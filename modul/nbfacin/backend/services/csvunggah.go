package services

// Unggah CSV objek TANPA BATAS baris (tiket 75 F1 + F2, keputusan work owner 10-10-2026 K75-1 / K75-2). Berkas dibaca MENGALIR
// (memori tidak bergantung jumlah baris) ke penampung T_NB_CSV_BARIS (migrasi 257); validasi dan simpan berjalan di LATAR per
// kelompok UkuranKelompokCsv baris dengan aturan per baris SAMA dengan Save Data tiket 66 (uraiBarisCsv, rencanakanBarisCsv,
// periksaObjek, generalCsv, CoverageAda). Semantik Save Data tetap MENGGANTI semua objek case: kelompok 1 = header General +
// hapus CurrencyList / cicilan + ganti objek (spreading dibuang); kelompok berikutnya MENAMBAH (SEQ_NO lokasi lanjut). Kemajuan
// pekerjaan ditulis DI TRANSAKSI YANG SAMA dengan objek kelompoknya, jadi Ulangi melanjutkan tepat dari kelompok yang belum COMMIT.
// Galat ketat seperti tiket 66 D7: satu saja galat baris -> gagal_validasi (berkas diperbaiki lalu diunggah ulang).
//
// Dialek berkas = pembaca layar tiket 66 (`frontend/components/csvObjek.ts` bacaCsvObjek): pemisah ";", kutip ganda per baris,
// BOM dibuang, baris kosong / berisi sel kosong semua dilewati, sel dipangkas, judul 61 kolom template atau 55 kolom Pega lama
// (tak peka huruf, sel kosong di ujung dibuang).

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// Status / fase pekerjaan unggah (kabel apa adanya).
const (
	StatusUnggahMemvalidasi   = "memvalidasi"
	StatusUnggahValid         = "valid"
	StatusUnggahGagalValidasi = "gagal_validasi"
	StatusUnggahMenyimpan     = "menyimpan"
	StatusUnggahSelesai       = "selesai"
	StatusUnggahGagal         = "gagal"
	FaseUnggahValidasi        = "validasi"
	FaseUnggahSimpan          = "simpan"
	jenisGalatUnggah          = "galat"
	jenisPesanUnggah          = "pesan"
	// UkuranKelompokCsv - baris per kelompok validasi / simpan (satu transaksi simpan per kelompok).
	UkuranKelompokCsv = 1000
	// BatasTampilUnggah - galat / pesan pertama yang dikirim status.
	BatasTampilUnggah = 100
	// kolomCsvPega - judul template Pega lama (55 kolom pertama KolomCsvObjek).
	kolomCsvPega = 55
	// batasPanjangBarisCsv - satu baris berkas lebih panjang dari ini = bukan CSV template (pengaman memori).
	batasPanjangBarisCsv = 1 << 20
)

var (
	// ErrUnggahBerjalan - case sedang memproses unggahan CSV (409).
	ErrUnggahBerjalan = errors.New("services: Unggahan CSV sedang diproses - tunggu sampai selesai")
	// ErrUnggahStatus - langkah tidak sah untuk status pekerjaan sekarang (409).
	ErrUnggahStatus = errors.New("services: status unggahan CSV tidak mengizinkan langkah ini")
	// ErrUnggahTidakAda - pekerjaan tidak ada / milik case lain (404).
	ErrUnggahTidakAda = errors.New("services: unggahan CSV tidak ada")
)

// StatusUnggahCsv - jawaban status pekerjaan: pekerjaan, galat / pesan pertama (paling banyak BatasTampilUnggah), case sudah berisi
// objek.
type StatusUnggahCsv struct {
	Unggah   models.UnggahCsv
	Galat    []models.PesanUnggahCsv
	Pesan    []models.PesanUnggahCsv
	AdaObjek bool
}

// DenganUnggahCsv memasang penyimpan pekerjaan unggah dan penambah objek (tiket 75).
func (s *Service) DenganUnggahCsv(u repository.PenyimpanUnggahCsv, p repository.PenambahObjek) *Service {
	s.unggahCsv, s.penambahObjek = u, p
	return s
}

// jalanLatar - pelaksana pekerjaan latar (nil = goroutine; uji: sinkron).
func (s *Service) jalanLatar(f func()) {
	if s.latar != nil {
		s.latar(f)
		return
	}
	go f()
}

// ---- pembaca mengalir ---------------------------------------------------------------------------------------------

// pecahBarisCsv - padanan pecahBarisCsv layar (`csvLossRecord.ts`): pemisah di luar kutip; "" di dalam kutip = ".
func pecahBarisCsv(baris string, pemisah rune) []string {
	var sel []string
	var kini strings.Builder
	dalam := false
	r := []rune(baris)
	for i := 0; i < len(r); i++ {
		c := r[i]
		switch {
		case dalam && c == '"' && i+1 < len(r) && r[i+1] == '"':
			kini.WriteRune('"')
			i++
		case dalam && c == '"':
			dalam = false
		case dalam:
			kini.WriteRune(c)
		case c == '"':
			dalam = true
		case c == pemisah:
			sel = append(sel, kini.String())
			kini.Reset()
		default:
			kini.WriteRune(c)
		}
	}
	return append(sel, kini.String())
}

// judulCsvSah - judul 61 kolom template atau 55 kolom Pega lama; tak peka huruf / spasi tepi; sel kosong ujung dibuang.
func judulCsvSah(baris string) bool {
	sel := pecahBarisCsv(baris, ';')
	for i := range sel {
		sel[i] = strings.ToUpper(strings.TrimSpace(sel[i]))
	}
	for len(sel) > 0 && sel[len(sel)-1] == "" {
		sel = sel[:len(sel)-1]
	}
	cocok := func(n int) bool {
		if len(sel) != n {
			return false
		}
		for i := range sel {
			if sel[i] != models.KolomCsvObjek[i] {
				return false
			}
		}
		return true
	}
	return cocok(len(models.KolomCsvObjek)) || cocok(kolomCsvPega)
}

// pesanJudulCsv - pesan judul salah (sama maksudnya dengan PESAN_CSV_OBJEK.judul layar).
const pesanJudulCsv = "judul berkas CSV tidak sesuai template objek (61 kolom, atau 55 kolom template lama, pemisah \";\")"

// BacaCsvObjekMengalir - judul lalu setiap baris data berisi dari aliran: onJudul sekali sesudah judul sah, onBaris(no, isi) per
// baris data (no 1..n, isi = 61 sel urut KolomCsvObjek, kekurangan = ""). Jawab jumlah baris data. Judul salah / tidak ada ->
// ErrMasukanObjek. Memori sebatas satu baris (paling panjang batasPanjangBarisCsv).
func BacaCsvObjekMengalir(r io.Reader, onJudul func() error, onBaris func(no int, isi []string) error) (int, error) {
	return bacaCsvMengalir(r, petaJudulFire, func(petaJudulCsv) error {
		if onJudul != nil {
			return onJudul()
		}
		return nil
	}, onBaris)
}

// bacaCsvMengalir - BacaCsvObjekMengalir dengan pemeta judul importer (tiket 74 P6): isi = sel urut kolom penampung peta.
func bacaCsvMengalir(r io.Reader, judulKe func(string) (petaJudulCsv, error), onJudul func(petaJudulCsv) error,
	onBaris func(no int, isi []string) error) (int, error) {
	br := bufio.NewReaderSize(r, 64<<10)
	judul, n := false, 0
	var peta petaJudulCsv
	pertama := true
	for {
		baris, err := bacaSatuBaris(br)
		if err != nil && !errors.Is(err, io.EOF) {
			return n, err
		}
		if baris != "" || err == nil {
			if pertama {
				baris = strings.TrimPrefix(baris, "\uFEFF")
				pertama = false
			}
			if strings.TrimSpace(baris) != "" {
				if !judul {
					p, e := judulKe(baris)
					if e != nil {
						return 0, e
					}
					judul, peta = true, p
					if e := onJudul(peta); e != nil {
						return 0, e
					}
				} else if isi, ada := selDataPeta(baris, peta); ada {
					n++
					if e := onBaris(n, isi); e != nil {
						return n, e
					}
				}
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	if !judul {
		return 0, fmt.Errorf("%w: %s", ErrMasukanObjek, pesanJudulCsv)
	}
	return n, nil
}

// bacaSatuBaris - satu baris tanpa \r\n; panjang > batasPanjangBarisCsv -> galat. io.EOF dikembalikan bersama baris terakhir.
func bacaSatuBaris(br *bufio.Reader) (string, error) {
	var b strings.Builder
	for {
		potong, err := br.ReadSlice('\n')
		b.Write(potong)
		if b.Len() > batasPanjangBarisCsv {
			return "", fmt.Errorf("%w: satu baris berkas lebih dari %d byte - bukan CSV template objek", ErrMasukanObjek, batasPanjangBarisCsv)
		}
		if errors.Is(err, bufio.ErrBufferFull) {
			continue
		}
		return strings.TrimRight(b.String(), "\r\n"), err
	}
}

// selDataCsv - sel baris data FIRE dipangkas dan dipetakan menurut posisi; semua kosong -> ada = false (sel berlebih ikut menandai
// baris berisi - layar: sel.every kosong).
func selDataCsv(baris string) ([]string, bool) {
	p, _ := petaJudulFire(strings.Join(models.KolomCsvObjek, ";"))
	return selDataPeta(baris, p)
}

// petaBarisCsv - isi penampung -> peta kolom (masukan uraiBarisCsv).
func petaBarisCsv(isi []string) map[string]string {
	m := make(map[string]string, len(models.KolomCsvObjek))
	for i, k := range models.KolomCsvObjek {
		if i < len(isi) {
			m[k] = isi[i]
		}
	}
	return m
}

// ---- kunci & status ----------------------------------------------------------------------------------------------

// berjalan - pekerjaan sedang memvalidasi / menyimpan dan masih berdetak.
func berjalan(u models.UnggahCsv) bool {
	return (u.Status == StatusUnggahMemvalidasi || u.Status == StatusUnggahMenyimpan) && !u.Basi
}

// cekKunciUnggah - K75-2: selama pekerjaan case berjalan, Save objek lain ditolak (409). Tanpa penyimpan unggah -> bebas.
func (s *Service) cekKunciUnggah(ctx context.Context, id string) error {
	if s.unggahCsv == nil {
		return nil
	}
	u, ada, err := s.unggahCsv.UnggahTerakhir(ctx, id)
	if err != nil {
		return err
	}
	if ada && berjalan(u) {
		return fmt.Errorf("%w (%s, %d/%d baris)", ErrUnggahBerjalan, u.Status, maks(u.BarisDiproses, u.BarisTersimpan), u.TotalBaris)
	}
	return nil
}

func maks(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// tulisUnggah - satu transaksi pendek: ubah pekerjaan (+ fn tambahan di transaksi yang sama).
func (s *Service) tulisUnggah(ctx context.Context, u models.UnggahCsv, fn func(tx *db.Tx) error) error {
	return s.transaksi(ctx, func(tx *db.Tx) error {
		if fn != nil {
			if err := fn(tx); err != nil {
				return err
			}
		}
		return s.unggahCsv.UbahUnggah(ctx, tx, u)
	})
}

// gagalUnggah - tandai gagal (fase tetap) dengan galat proses dan baris awal kelompok; galat penandaan hanya dicatat.
func (s *Service) gagalUnggah(ctx context.Context, u models.UnggahCsv, diBaris int, sebab error) {
	u.Status, u.GalatProses, u.GalatDiBaris = StatusUnggahGagal, sebab.Error(), diBaris
	if diBaris > 0 {
		u.GalatProses = fmt.Sprintf("gagal di baris %d: %v", diBaris, sebab)
	}
	if err := s.tulisUnggah(ctx, u, nil); err != nil {
		log.Printf("nbfacin: unggah CSV %s: menandai gagal: %v (sebab: %v)", u.ID, err, sebab)
	}
	log.Printf("nbfacin: unggah CSV %s case %s gagal (%s): %v", u.ID, u.CaseID, u.Fase, sebab)
}

// prasyaratUnggah - identitas, tahap boleh ubah, dependensi, ID case, case ada, importer menurut Class of Business (jenisCsvKasus)
// dan dependensi importer itu. Jawab jenis importer.
func (s *Service) prasyaratUnggah(ctx context.Context, pelaku inti.Pelaku, id string) (string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return "", err
	}
	if s.unggahCsv == nil || s.transaksi == nil || s.csv == nil || s.kasus == nil {
		return "", ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return "", ErrKasusTidakAda
	}
	if _, err := s.kasus.BacaKasus(ctx, id); errors.Is(err, repository.ErrKasusTidakAda) {
		return "", ErrKasusTidakAda
	} else if err != nil {
		return "", err
	}
	jenis, err := s.jenisCsvKasus(ctx, id)
	if err != nil {
		return "", err
	}
	if jenis == models.JenisCsvCargo {
		if s.objekCargo == nil || s.penulisCargoCsv == nil || s.kapalMaster == nil {
			return "", ErrObjekTanpaDatabase
		}
		if s.masterCargo == nil {
			return "", fmt.Errorf("%w: Trading / Goods / Packing / Conveyance / Agent", ErrMasterCargoBelumAda)
		}
		return jenis, nil
	}
	if jenis == models.JenisCsvAneka { // tiket 73 P6
		if s.penulisAnekaCsv == nil || s.lookupAneka == nil {
			return "", ErrObjekTanpaDatabase
		}
		return jenis, nil
	}
	if s.penambahObjek == nil || s.objek == nil || s.rw == nil {
		return "", ErrObjekTanpaDatabase
	}
	return jenis, nil
}

// wajibJenisUnggah - pekerjaan `u` hanya dilanjutkan untuk case yang importernya masih sama (Class of Business bisa berubah).
func (s *Service) wajibJenisUnggah(jenis string, u models.UnggahCsv) error {
	if u.Jenis != jenis {
		return fmt.Errorf("%w: pekerjaan unggah CSV %s tidak sesuai Class of Business case (%s)", ErrUnggahStatus, u.Jenis, jenis)
	}
	return nil
}

// bacaUnggahKasus - pekerjaan `idUnggah` milik case `id` (404 bila tidak).
func (s *Service) bacaUnggahKasus(ctx context.Context, id, idUnggah string) (models.UnggahCsv, error) {
	if s.unggahCsv == nil {
		return models.UnggahCsv{}, ErrObjekTanpaDatabase
	}
	u, err := s.unggahCsv.BacaUnggah(ctx, idUnggah)
	if errors.Is(err, repository.ErrUnggahTidakAda) || (err == nil && u.CaseID != id) {
		return models.UnggahCsv{}, ErrUnggahTidakAda
	}
	return u, err
}

// statusUnggah - status lengkap; pekerjaan berjalan yang basi ditandai gagal lebih dulu (dapat diulang).
func (s *Service) statusUnggah(ctx context.Context, u models.UnggahCsv) (StatusUnggahCsv, error) {
	if (u.Status == StatusUnggahMemvalidasi || u.Status == StatusUnggahMenyimpan) && u.Basi {
		s.gagalUnggah(ctx, u, 0, fmt.Errorf("pekerjaan terputus (tanpa kemajuan lebih dari %d menit) - gunakan Ulangi",
			models.BatasDetakUnggahMenit))
		var err error
		if u, err = s.unggahCsv.BacaUnggah(ctx, u.ID); err != nil {
			return StatusUnggahCsv{}, err
		}
	}
	galat, err := s.unggahCsv.BacaPesanUnggah(ctx, u.ID, jenisGalatUnggah, BatasTampilUnggah)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	pesan, err := s.unggahCsv.BacaPesanUnggah(ctx, u.ID, jenisPesanUnggah, BatasTampilUnggah)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	if u.Status == StatusUnggahSelesai && len(pesan) == 0 {
		pesan = []models.PesanUnggahCsv{{Teks: pesanSuksesCsv}}
	}
	ada, err := s.unggahCsv.AdaObjek(ctx, u.CaseID, u.Jenis)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	return StatusUnggahCsv{Unggah: u, Galat: galat, Pesan: pesan, AdaObjek: ada}, nil
}

// StatusUnggahCsvKasus - GET …/objek/csv/{idUnggah}.
func (s *Service) StatusUnggahCsvKasus(ctx context.Context, id, idUnggah string) (StatusUnggahCsv, error) {
	u, err := s.bacaUnggahKasus(ctx, id, idUnggah)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	return s.statusUnggah(ctx, u)
}

// UnggahCsvAktif - GET …/objek/csv/aktif: pekerjaan terakhir case (ada = false bila belum pernah).
func (s *Service) UnggahCsvAktif(ctx context.Context, id string) (StatusUnggahCsv, bool, error) {
	if s.unggahCsv == nil {
		return StatusUnggahCsv{}, false, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return StatusUnggahCsv{}, false, ErrKasusTidakAda
	}
	u, ada, err := s.unggahCsv.UnggahTerakhir(ctx, id)
	if err != nil || !ada {
		return StatusUnggahCsv{}, false, err
	}
	st, err := s.statusUnggah(ctx, u)
	return st, err == nil, err
}

// HalamanBarisUnggah - satu halaman View Upload: Baris (Isi urut models.KolomCsvJenis(Jenis)), Total baris data, Kolom = judul
// tampil urut berkas (tiket 74 P6; pekerjaan FIRE = urutan template).
type HalamanBarisUnggah struct {
	Baris []models.BarisUnggahCsv
	Total int
	Jenis string
	Kolom []string
}

// BarisUnggahCsv - GET …/objek/csv/{idUnggah}/baris: View Upload per halaman (halaman mulai 1).
func (s *Service) BarisUnggahCsv(ctx context.Context, id, idUnggah string, halaman, ukuran int) (HalamanBarisUnggah, error) {
	u, err := s.bacaUnggahKasus(ctx, id, idUnggah)
	if err != nil {
		return HalamanBarisUnggah{}, err
	}
	if halaman < 1 || ukuran < 1 || ukuran > 1000 {
		return HalamanBarisUnggah{}, fmt.Errorf("%w: halaman >= 1, ukuran 1..1000", ErrMasukanObjek)
	}
	h := HalamanBarisUnggah{Baris: []models.BarisUnggahCsv{}, Total: u.TotalBaris, Jenis: u.Jenis, Kolom: u.Kolom}
	if h.Kolom == nil {
		h.Kolom = models.KolomCsvJenis(u.Jenis)
	}
	dari := (halaman-1)*ukuran + 1
	if dari > u.TotalBaris {
		return h, nil
	}
	h.Baris, err = s.unggahCsv.BacaBarisUnggah(ctx, u.ID, u.Jenis, dari, min(dari+ukuran-1, u.TotalBaris))
	return h, err
}

// ---- F1: unggah + validasi latar ---------------------------------------------------------------------------------

// MulaiUnggahCsv - POST …/objek/csv/berkas: berkas dibaca mengalir ke penampung (kelompok UkuranKelompokCsv baris), lalu
// validasi latar dimulai. Judul salah -> 400 tanpa pekerjaan. Aliran putus di tengah -> pekerjaan gagal + galat.
func (s *Service) MulaiUnggahCsv(ctx context.Context, pelaku inti.Pelaku, id, nama string, r io.Reader) (StatusUnggahCsv, error) {
	jenis, err := s.prasyaratUnggah(ctx, pelaku, id)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	if err := s.cekKunciUnggah(ctx, id); err != nil {
		return StatusUnggahCsv{}, err
	}
	hitung := &penghitungByte{r: r}
	u := models.UnggahCsv{CaseID: id, Status: StatusUnggahMemvalidasi, Fase: FaseUnggahValidasi, NamaBerkas: nama, DibuatOleh: pelaku.AkunID,
		Jenis: jenis}
	var tampung []models.BarisUnggahCsv
	salur := func() error {
		if len(tampung) == 0 {
			return nil
		}
		err := s.transaksi(ctx, func(tx *db.Tx) error { return s.unggahCsv.SisipBarisUnggah(ctx, tx, u.ID, u.Jenis, tampung) })
		tampung = tampung[:0]
		return err
	}
	total, err := bacaCsvMengalir(hitung, petaJudulJenis(jenis), func(p petaJudulCsv) error {
		u.Kolom = p.kolom
		return s.transaksi(ctx, func(tx *db.Tx) error {
			var err error
			u.ID, err = s.unggahCsv.BuatUnggah(ctx, tx, u)
			return err
		})
	}, func(no int, isi []string) error {
		tampung = append(tampung, models.BarisUnggahCsv{No: no, Isi: isi})
		if len(tampung) >= UkuranKelompokCsv {
			return salur()
		}
		return nil
	})
	if err == nil {
		err = salur()
	}
	u.TotalBaris, u.Ukuran = total, hitung.n
	if err != nil {
		if u.ID == "" {
			return StatusUnggahCsv{}, err // judul salah / tidak ada pekerjaan
		}
		s.gagalUnggah(ctx, u, 0, fmt.Errorf("unggahan berkas terputus sesudah %d baris: %w", total, err))
		return StatusUnggahCsv{}, err
	}
	if total == 0 {
		u.Status, u.JumlahGalat = StatusUnggahGagalValidasi, 1
		if err := s.tulisUnggah(ctx, u, func(tx *db.Tx) error {
			return s.unggahCsv.SisipPesanUnggah(ctx, tx, u.ID, jenisGalatUnggah, 1, []models.PesanUnggahCsv{{Teks: pesanUnggahCsv}})
		}); err != nil {
			return StatusUnggahCsv{}, err
		}
	} else {
		if err := s.tulisUnggah(ctx, u, nil); err != nil {
			return StatusUnggahCsv{}, err
		}
		latar := context.WithoutCancel(ctx)
		s.jalanLatar(func() { s.validasiUnggah(latar, u) })
	}
	return s.StatusUnggahCsvKasus(ctx, id, u.ID)
}

// penghitungByte - ukuran berkas yang diterima.
type penghitungByte struct {
	r io.Reader
	n int64
}

func (p *penghitungByte) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.n += int64(n)
	return n, err
}

// headerTeks - langkah 5.3 D1 (b) atas teks CSV (baris terakhir yang terisi menang).
type headerTeks struct{ noOfferSlip, mulai, akhir string }

func (h *headerTeks) ambil(b barisCsv) {
	for _, x := range []struct {
		v  string
		ke *string
	}{{b.t("NO_OFFER_SLIP"), &h.noOfferSlip}, {b.t("BEGIN_DATE"), &h.mulai}, {b.t("END_DATE"), &h.akhir}} {
		if x.v != "" {
			*x.ke = x.v
		}
	}
}

// headerCsvDari - header tersimpan pekerjaan -> headerCsv (tanggal sudah lolos validasi).
func headerCsvDari(u models.UnggahCsv) headerCsv {
	h := headerCsv{noOfferSlip: u.NoOfferSlip}
	if t, ok := tanggalCsv(u.BeginDate); ok {
		h.mulai = t
	}
	if t, ok := tanggalCsv(u.EndDate); ok {
		h.akhir = t
	}
	return h
}

// validasiKelompok - satu kelompok baris: rencana objek, galat (bernomor baris), pesan KodePos. Galat basis data = error.
func (s *Service) validasiKelompok(ctx context.Context, baris []models.BarisUnggahCsv, pakaiCoverage map[string]string,
	h *headerTeks) ([]rencanaCsv, []models.PesanUnggahCsv, []models.PesanUnggahCsv, error) {
	var rencana []rencanaCsv
	var galat, pesan []models.PesanUnggahCsv
	for _, x := range baris {
		b, masalah := uraiBarisCsv(x.No, petaBarisCsv(x.Isi))
		for _, m := range masalah {
			galat = append(galat, models.PesanUnggahCsv{NoBaris: x.No, Teks: m})
		}
		if len(masalah) > 0 {
			continue
		}
		if h != nil {
			h.ambil(b)
		}
		r, p, m, err := s.rencanakanBarisCsv(ctx, b, pakaiCoverage)
		if err != nil {
			return nil, nil, nil, err
		}
		for _, t := range p {
			pesan = append(pesan, models.PesanUnggahCsv{NoBaris: x.No, Teks: t})
		}
		if len(m) > 0 { // urutan jalur lama: galat master (D6) menghentikan sebelum periksaObjek
			for _, t := range m {
				galat = append(galat, models.PesanUnggahCsv{NoBaris: x.No, Teks: t})
			}
			continue
		}
		if err := periksaObjek([]models.ObjekFire{r.objek}); err != nil {
			galat = append(galat, models.PesanUnggahCsv{NoBaris: x.No, Teks: fmt.Sprintf("baris %d: %v", x.No, err)})
			continue
		}
		rencana = append(rencana, r)
	}
	return rencana, galat, pesan, nil
}

// validasiUnggah - pekerjaan latar fase validasi: semua kelompok, lalu coverage master + header General. Galat basis data ->
// gagal (dapat diulang).
func (s *Service) validasiUnggah(ctx context.Context, u models.UnggahCsv) {
	sc := s.denganCacheCsv()
	u.Status, u.Fase, u.BarisDiproses, u.JumlahGalat, u.JumlahPesan, u.GalatProses, u.GalatDiBaris =
		StatusUnggahMemvalidasi, FaseUnggahValidasi, 0, 0, 0, "", 0
	if err := s.tulisUnggah(ctx, u, func(tx *db.Tx) error { return s.unggahCsv.HapusPesanUnggah(ctx, tx, u.ID) }); err != nil {
		s.gagalUnggah(ctx, u, 0, err)
		return
	}
	if u.Jenis == models.JenisCsvCargo {
		s.validasiUnggahCargo(ctx, u)
		return
	}
	if u.Jenis == models.JenisCsvAneka { // tiket 73 P6
		sc.validasiUnggahAneka(ctx, u)
		return
	}
	pakaiCoverage := map[string]string{}
	var h headerTeks
	for dari := 1; dari <= u.TotalBaris; dari += UkuranKelompokCsv {
		baris, err := s.unggahCsv.BacaBarisUnggah(ctx, u.ID, u.Jenis, dari, min(dari+UkuranKelompokCsv-1, u.TotalBaris))
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		_, galat, pesan, err := sc.validasiKelompok(ctx, baris, pakaiCoverage, &h)
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		urutGalat, urutPesan := u.JumlahGalat+1, u.JumlahPesan+1
		u.BarisDiproses, u.JumlahGalat, u.JumlahPesan = dari+len(baris)-1, u.JumlahGalat+len(galat), u.JumlahPesan+len(pesan)
		if err := s.tulisUnggah(ctx, u, func(tx *db.Tx) error {
			if err := s.unggahCsv.SisipPesanUnggah(ctx, tx, u.ID, jenisGalatUnggah, urutGalat, galat); err != nil {
				return err
			}
			return s.unggahCsv.SisipPesanUnggah(ctx, tx, u.ID, jenisPesanUnggah, urutPesan, pesan)
		}); err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
	}
	var akhir []models.PesanUnggahCsv
	if err := s.periksaCoverageMasterCsv(ctx, pakaiCoverage); err != nil {
		if !errors.Is(err, ErrMasukanObjek) {
			s.gagalUnggah(ctx, u, 0, err)
			return
		}
		akhir = append(akhir, models.PesanUnggahCsv{Teks: err.Error()})
	}
	u.NoOfferSlip, u.BeginDate, u.EndDate = h.noOfferSlip, h.mulai, h.akhir
	if u.JumlahGalat == 0 {
		kasus, err := s.kasus.BacaKasus(ctx, u.CaseID)
		if err != nil {
			s.gagalUnggah(ctx, u, 0, err)
			return
		}
		if _, err := s.generalCsv(kasus, headerCsvDari(u)); err != nil {
			akhir = append(akhir, models.PesanUnggahCsv{Teks: err.Error()})
		}
	}
	urut := u.JumlahGalat + 1
	u.JumlahGalat += len(akhir)
	u.Status = StatusUnggahValid
	if u.JumlahGalat > 0 {
		u.Status = StatusUnggahGagalValidasi
	}
	if err := s.tulisUnggah(ctx, u, func(tx *db.Tx) error {
		return s.unggahCsv.SisipPesanUnggah(ctx, tx, u.ID, jenisGalatUnggah, urut, akhir)
	}); err != nil {
		s.gagalUnggah(ctx, u, 0, err)
	}
}

// periksaCoverageMasterCsv - D8: ID coverage ter-hardcode yang dipakai wajib ada di master (ErrMasukanObjek bila tidak).
func (s *Service) periksaCoverageMasterCsv(ctx context.Context, pakai map[string]string) error {
	ids := make([]string, 0, len(pakai))
	for id := range pakai {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil
	}
	ada, err := s.csv.CoverageAda(ctx, ids)
	if err != nil {
		return err
	}
	var masalah []string
	for _, id := range ids {
		if !ada[id] {
			masalah = append(masalah, fmt.Sprintf("coverage %s (%s) tidak ada di master COVERAGE", id, pakai[id]))
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	return nil
}

// ---- F2: simpan latar per kelompok + ulangi ----------------------------------------------------------------------

// SimpanUnggahCsv - POST …/objek/csv/{idUnggah}/simpan: hanya dari 'valid'; mulai simpan latar dari kelompok 1.
func (s *Service) SimpanUnggahCsv(ctx context.Context, pelaku inti.Pelaku, id, idUnggah string) (StatusUnggahCsv, error) {
	jenis, err := s.prasyaratUnggah(ctx, pelaku, id)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	u, err := s.bacaUnggahKasus(ctx, id, idUnggah)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	if err := s.wajibJenisUnggah(jenis, u); err != nil {
		return StatusUnggahCsv{}, err
	}
	if u.Status != StatusUnggahValid {
		return StatusUnggahCsv{}, fmt.Errorf("%w: simpan hanya dari status valid (sekarang %s)", ErrUnggahStatus, u.Status)
	}
	u.Status, u.Fase, u.BarisTersimpan, u.KelompokTerakhir, u.GalatProses, u.GalatDiBaris =
		StatusUnggahMenyimpan, FaseUnggahSimpan, 0, 0, "", 0
	if err := s.tulisUnggah(ctx, u, nil); err != nil {
		return StatusUnggahCsv{}, err
	}
	latar := context.WithoutCancel(ctx)
	s.jalanLatar(func() { s.simpanUnggah(latar, pelaku, u) })
	return s.StatusUnggahCsvKasus(ctx, id, u.ID)
}

// UlangiUnggahCsv - POST …/objek/csv/{idUnggah}/ulangi: dari 'gagal'; fase validasi -> validasi ulang; fase simpan -> lanjut dari
// kelompok sesudah KELOMPOK_TERAKHIR.
func (s *Service) UlangiUnggahCsv(ctx context.Context, pelaku inti.Pelaku, id, idUnggah string) (StatusUnggahCsv, error) {
	jenis, err := s.prasyaratUnggah(ctx, pelaku, id)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	u, err := s.bacaUnggahKasus(ctx, id, idUnggah)
	if err != nil {
		return StatusUnggahCsv{}, err
	}
	if err := s.wajibJenisUnggah(jenis, u); err != nil {
		return StatusUnggahCsv{}, err
	}
	if u.Basi && (u.Status == StatusUnggahMemvalidasi || u.Status == StatusUnggahMenyimpan) {
		u.Status = StatusUnggahGagal // terputus: boleh diulang langsung
	}
	if u.Status != StatusUnggahGagal {
		return StatusUnggahCsv{}, fmt.Errorf("%w: ulangi hanya dari status gagal (sekarang %s)", ErrUnggahStatus, u.Status)
	}
	latar := context.WithoutCancel(ctx)
	if u.Fase == FaseUnggahSimpan {
		u.Status, u.GalatProses, u.GalatDiBaris = StatusUnggahMenyimpan, "", 0
		if err := s.tulisUnggah(ctx, u, nil); err != nil {
			return StatusUnggahCsv{}, err
		}
		s.jalanLatar(func() { s.simpanUnggah(latar, pelaku, u) })
	} else {
		u.Status, u.GalatProses, u.GalatDiBaris = StatusUnggahMemvalidasi, "", 0
		if err := s.tulisUnggah(ctx, u, nil); err != nil {
			return StatusUnggahCsv{}, err
		}
		s.jalanLatar(func() { s.validasiUnggah(latar, u) })
	}
	return s.StatusUnggahCsvKasus(ctx, id, u.ID)
}

// simpanUnggah - pekerjaan latar fase simpan dari kelompok KELOMPOK_TERAKHIR+1. Kemajuan ditulis di transaksi objek kelompoknya.
func (s *Service) simpanUnggah(ctx context.Context, pelaku inti.Pelaku, u models.UnggahCsv) {
	// bug DEV 11-10-2026: pengaman jalur latar - importer pekerjaan wajib sama dengan Class of Business case sekarang.
	jenis, err := s.jenisCsvKasus(ctx, u.CaseID)
	if err == nil {
		err = s.wajibJenisUnggah(jenis, u)
	}
	if err != nil {
		s.gagalUnggah(ctx, u, 0, err)
		return
	}
	if u.Jenis == models.JenisCsvCargo {
		s.simpanUnggahCargo(ctx, pelaku, u)
		return
	}
	if u.Jenis == models.JenisCsvAneka { // tiket 73 P6
		s.denganCacheCsv().simpanUnggahAneka(ctx, u) // master ber-cache per pekerjaan
		return
	}
	sc := s.denganCacheCsv()
	kasus, err := s.kasus.BacaKasus(ctx, u.CaseID)
	if err != nil {
		s.gagalUnggah(ctx, u, 0, err)
		return
	}
	g, err := s.generalCsv(kasus, headerCsvDari(u))
	if err != nil {
		s.gagalUnggah(ctx, u, 0, err)
		return
	}
	periode := &Periode{Mulai: g.StartDateTime, Akhir: g.EndDateTime, EDMDay: g.Day}
	for k := u.KelompokTerakhir + 1; (k-1)*UkuranKelompokCsv < u.TotalBaris; k++ {
		dari := (k-1)*UkuranKelompokCsv + 1
		baris, err := s.unggahCsv.BacaBarisUnggah(ctx, u.ID, u.Jenis, dari, min(dari+UkuranKelompokCsv-1, u.TotalBaris))
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		rencana, galat, _, err := sc.validasiKelompok(ctx, baris, map[string]string{}, nil)
		if err == nil && len(galat) > 0 {
			err = fmt.Errorf("%w: data master berubah sejak validasi: %s", ErrMasukanObjek, galat[0].Teks)
		}
		if err == nil {
			err = sc.tulisMasterCsv(ctx, pelaku, u.CaseID, rencana)
		}
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		objek := make([]models.ObjekFire, len(rencana))
		for i, r := range rencana {
			objek[i] = r.objek
		}
		maju := u
		maju.KelompokTerakhir, maju.BarisTersimpan = k, dari+len(baris)-1
		if k == 1 {
			err = sc.gantiObjek(ctx, u.CaseID, objek, opsiGantiObjek{periode: periode, buangSpreading: true,
				sebelum: func(tx *db.Tx) error {
					if err := s.kasus.SimpanGeneral(ctx, tx, u.CaseID, g); err != nil {
						return err
					}
					if err := s.csv.HapusCurrencyList(ctx, tx, u.CaseID); err != nil {
						return err
					}
					return s.unggahCsv.UbahUnggah(ctx, tx, maju)
				}})
		} else {
			err = sc.tambahObjek(ctx, u.CaseID, objek, dari-1, periode, func(tx *db.Tx) error {
				return s.unggahCsv.UbahUnggah(ctx, tx, maju)
			})
		}
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		u = maju
	}
	u.Status = StatusUnggahSelesai
	if err := s.tulisUnggah(ctx, u, nil); err != nil {
		s.gagalUnggah(ctx, u, 0, err)
		return
	}
	log.Printf("nbfacin: unggah CSV %s case %s oleh %s: %d lokasi tersimpan", u.ID, u.CaseID, pelaku.AkunID, u.BarisTersimpan)
}

// tambahObjek - padanan gantiObjek untuk kelompok berikutnya: mata uang, kerugian, coverage (periode header), lalu TambahObjek
// (SEQ_NO mulai awal+1) di satu transaksi bersama `sebelum`.
func (s *Service) tambahObjek(ctx context.Context, id string, baris []models.ObjekFire, awal int, periode *Periode,
	sebelum func(tx *db.Tx) error) error {
	if err := s.periksaMataUang(ctx, baris); err != nil {
		return err
	}
	simpan, err := s.siapkanKerugian(ctx, id, baris)
	if err != nil {
		return err
	}
	if err := s.siapkanCoverageDengan(ctx, id, simpan, periode); err != nil {
		return err
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		if sebelum != nil {
			if err := sebelum(tx); err != nil {
				return err
			}
		}
		return s.penambahObjek.TambahObjek(ctx, tx, id, simpan, awal)
	})
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return ErrKasusTidakAda
	}
	return err
}
