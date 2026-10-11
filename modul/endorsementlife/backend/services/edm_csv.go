package services

// Tiket 07 - unggah CSV peserta baru (`Upload CSV` b8973, `Add CSV Data`
// b10405 → `SaveCSVEDMLife`).
//
// ⛔ TANPA BATAS BARIS (spec penyimpangan 6, AC 36): berkas dibaca BERTAHAP -
// satu baris di memori untuk validasi, satu kelompok `ukuranKelompokCSV` untuk
// penyisipan. Gerbang `>50000` `SetPremi_EDM` b3373 tidak ditiru.
//
// ⛔ SELURUH BARIS DIPERIKSA DULU (R24). Pega menambahkan baris sampai baris
// pertama yang `PLAN`/`POLICY_HOLDER`-nya beda lalu berhenti; di sini satu
// pelanggaran menolak seluruh berkas - lintasan pertama hanya memeriksa,
// lintasan kedua (di dalam transaksi) memeriksa ulang dan menyisip.
//
// ⛔ Sesudah Process Endorsement (`Save`, rekap ada) unggahan DITOLAK 409
// `ErrSudahDiproses` - Upload CSV (periksa) maupun Add CSV Data, sebelum tulisan
// apa pun (keputusan work owner 08-10-2026; mengganti R31 yang mengikuti Pega:
// CSV boleh sesudah Save lalu rekap dihitung ulang).
//
// ⛔ Sesudah `Add CSV Data` unggahan terkunci (`.EditInput1 = 1` b2899 →
// `pyDisabledWhen` b8965/b10403): kasus ber-baris `New` menolak unggahan
// berikutnya (409). Pembuangan baris `New` langkah 2.1 tetap dijalankan.

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

const (
	// ukuranKelompokCSV - baris per kelompok penyisipan.
	ukuranKelompokCSV = 1000
	// batasPesanCSV - penolakan yang dikembalikan; sisanya hanya dicacah.
	batasPesanCSV = 200
)

var (
	// ErrCSVBukanPerubahanData - baris `New` hanya lahir pada `EdmType=1`
	// (4.3 b2791; wadah unggah b8698 `.EdmType=1`) - AC 13.
	ErrCSVBukanPerubahanData = errors.New("services: CSV upload is only available for EDM Type Perubahan Data.")
	// ErrSudahDiproses - Process Endorsement (`Save`) sudah dijalankan; CSV
	// terkunci (keputusan work owner 08-10-2026 mengganti R31). 409 seperti
	// `ErrSudahDisimpan`.
	ErrSudahDiproses = errors.New("services: " + models.PesanSudahDiproses)
	// ErrCSVTerkunci - `Add CSV Data` sudah dijalankan (`.EditInput1 = 1`).
	ErrCSVTerkunci = errors.New("services: CSV data has already been added to this endorsement, so another file cannot be uploaded.")
	// ErrCSVKosong - berkas tanpa baris data.
	ErrCSVKosong = errors.New("services: The file has no data rows below the header.")
	// ErrCSVRusak - berkas tidak dapat diurai sebagai CSV.
	ErrCSVRusak = errors.New("services: The file cannot be read as a CSV file.")
	// ErrCSVTanpaAcuan - kasus tanpa peserta lama: `PremiumListDetail(1)` tidak ada.
	ErrCSVTanpaAcuan = errors.New("services: This endorsement has no existing participants, so new participants cannot be added.")
)

// HasilPeriksaCSV - tinjauan unggahan: cacah baris, yang ditolak, dan pesannya.
type HasilPeriksaCSV struct {
	Total     int               `json:"total"`
	Ditolak   int               `json:"ditolak"`
	Pesan     []models.PesanCSV `json:"pesan"`
	Terpotong bool              `json:"terpotong"`
	// Diabaikan - judul kolom yang bukan properti 4.1: tidak disimpan, ditampilkan.
	Diabaikan []string `json:"diabaikan"`
}

// GalatCSV - berkas ditolak karena barisnya; nol baris tersimpan.
type GalatCSV struct{ HasilPeriksaCSV }

func (g GalatCSV) Error() string {
	return fmt.Sprintf("services: %d of %d rows have a problem; no row was saved.", g.Ditolak, g.Total)
}

// GalatBerkasCSV - penolakan berkas/judul CSV berkalimat layar BERSIH (10-10-2026): `PesanLayar` hanya kalimat untuk
// user (tanpa "invalid endorsement input:" / "models:"), `Is` menyamakannya dengan `jenis` (ErrMasukanTidakSah atau
// ErrCSVRusak) sehingga pemetaan HTTP (422) TIDAK berubah. Hanya jalur CSV yang memakainya; jalur lain tetap.
type GalatBerkasCSV struct {
	jenis error
	asal  error
	Pesan string
}

func (g GalatBerkasCSV) Error() string        { return g.jenis.Error() + ": " + g.Pesan }
func (g GalatBerkasCSV) PesanLayar() string   { return g.Pesan }
func (g GalatBerkasCSV) Is(target error) bool { return target == g.jenis }
func (g GalatBerkasCSV) Unwrap() error        { return g.asal }

// galatBerkasCSV - galat models berkalimat layar (`models.PesanGalatCSV`) menjadi GalatBerkasCSV 422.
func galatBerkasCSV(err error) error {
	pesan := err.Error()
	var l pesanLayar
	if errors.As(err, &l) {
		pesan = l.PesanLayar()
	}
	return GalatBerkasCSV{jenis: ErrMasukanTidakSah, asal: err, Pesan: pesan}
}

// HasilTambahCSV - hasil `Add CSV Data`.
type HasilTambahCSV struct {
	Disimpan int                 `json:"disimpan"`
	Dibuang  int                 `json:"dibuang"`
	Rekap    []map[string]string `json:"rekap"`
}

// periksaKasusCSV - gerbang kasus untuk unggahan; `k.SudahSimpan` = rekap ada.
func periksaKasusCSV(k models.Kasus) error {
	switch {
	case !k.Terbuka():
		return ErrKasusTertutup
	case k.SudahSimpan:
		return ErrSudahDiproses
	case k.EdmType != models.EdmTypePerubahanData:
		return ErrCSVBukanPerubahanData
	case k.Cacah[models.StatusNew] > 0:
		return ErrCSVTerkunci
	}
	return nil
}

// pemisahCSV - pemisah unggahan (keputusan work owner 10-10-2026, opsi a): hanya baris JUDUL yang dibaca (bufio,
// tetap bertahap) - memuat ";" dan tidak memuat "," → ";", selain itu ",". BOM UTF-8 diabaikan untuk deteksi (tetap
// di aliran; `JudulCSV` membuangnya). Baris judul dikembalikan ke aliran (MultiReader), jadi penomoran baris (judul
// = 1) tidak bergeser. Baris petunjuk Excel "sep=" TIDAK didukung (tetap dibaca sebagai judul).
func pemisahCSV(berkas io.Reader) (io.Reader, rune, error) {
	br := bufio.NewReader(berkas)
	judul, err := br.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, 0, galatBaca(err)
	}
	periksa := strings.TrimPrefix(judul, string(rune(0xFEFF)))
	pemisah := ','
	if strings.Contains(periksa, ";") && !strings.Contains(periksa, ",") {
		pemisah = ';'
	}
	return io.MultiReader(strings.NewReader(judul), br), pemisah, nil
}

// telusuriCSV membaca berkas baris demi baris; `sisip` menerima kelompok baris
// lolos SELAMA belum ada penolakan (nil = periksa saja).
func telusuriCSV(berkas io.Reader, acuan models.AcuanCSV, sisip func([]models.BarisCSV) error) (HasilPeriksaCSV, error) {
	aliran, pemisah, err := pemisahCSV(berkas)
	if err != nil {
		return HasilPeriksaCSV{}, err
	}
	c := csv.NewReader(aliran)
	c.Comma = pemisah
	c.FieldsPerRecord = -1
	c.TrimLeadingSpace = true
	c.ReuseRecord = true
	mentah, err := c.Read()
	if errors.Is(err, io.EOF) {
		return HasilPeriksaCSV{}, ErrCSVKosong
	}
	if err != nil {
		return HasilPeriksaCSV{}, galatBaca(err)
	}
	judul, asing, err := models.JudulCSV(append([]string{}, mentah...), acuan.Tipe)
	if err != nil {
		return HasilPeriksaCSV{}, galatBerkasCSV(err)
	}
	h := HasilPeriksaCSV{Pesan: []models.PesanCSV{}, Diabaikan: asing}
	if h.Diabaikan == nil {
		h.Diabaikan = []string{}
	}
	tolak := func(p ...models.PesanCSV) {
		for _, x := range p {
			if len(h.Pesan) < batasPesanCSV {
				h.Pesan = append(h.Pesan, x)
			} else {
				h.Terpotong = true
			}
		}
	}
	// CERTIFICATE_NO ganda DALAM berkas (10-10-2026): baris pertama per kunci (`models.KunciSertifikat`).
	pertama := map[string]int{}
	var kelompok []models.BarisCSV
	kirim := func() error {
		if sisip == nil || h.Ditolak > 0 || len(kelompok) == 0 {
			kelompok = kelompok[:0]
			return nil
		}
		err := sisip(kelompok)
		kelompok = nil
		return err
	}
	for {
		rec, err := c.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return h, galatBaca(err)
		}
		if strings.TrimSpace(strings.Join(rec, "")) == "" {
			continue // baris kosong - bukan peserta
		}
		h.Total++
		// Nomor baris = baris BERKAS tempat rekamannya mulai (judul = 1), seperti
		// nomor baris spreadsheet - baris kosong dan medan berkutip-banyak-baris
		// tidak menggeser laporan penolakan (AC 37).
		nomor, _ := c.FieldPos(0)
		baris := make(map[string]string, len(judul))
		for i, v := range rec {
			if i < len(judul) {
				baris[judul[i]] = v
			}
		}
		nilai, pesan := models.RapikanBarisCSV(nomor, baris, acuan)
		if cn := models.KunciSertifikat(baris["CERTIFICATE_NO"]); cn != "" {
			if awal, ada := pertama[cn]; ada {
				pesan = append(pesan, models.PesanCSV{Baris: nomor, Kolom: "CERTIFICATE_NO", Pesan: fmt.Sprintf(models.PesanCSVSertifikatGanda, strings.TrimSpace(baris["CERTIFICATE_NO"]), awal)})
			} else {
				pertama[cn] = nomor
			}
		}
		if len(rec) > len(judul) {
			pesan = append(pesan, models.PesanCSV{Baris: nomor, Pesan: fmt.Sprintf(models.PesanCSVJumlahMedan, len(rec), len(judul))})
		}
		if len(pesan) > 0 {
			h.Ditolak++
			tolak(pesan...)
			continue
		}
		if sisip != nil && h.Ditolak == 0 {
			kelompok = append(kelompok, models.BarisCSV{Nomor: nomor, Nilai: nilai})
			if len(kelompok) >= ukuranKelompokCSV {
				if err := kirim(); err != nil {
					return h, err
				}
			}
		}
	}
	if h.Total == 0 {
		return h, ErrCSVKosong
	}
	return h, kirim()
}

// galatBaca - rusak bentuk CSV menjadi `ErrCSVRusak` berbaris/berkolom (422);
// galat membaca badan (putus, terlalu besar, berkas sementara) TIDAK - ia
// bukan kesalahan berkas pengguna, dan teksnya (jalur sementara) tidak ke layar.
func galatBaca(err error) error {
	var pe *csv.ParseError
	if errors.As(err, &pe) {
		return GalatBerkasCSV{jenis: ErrCSVRusak, asal: err,
			Pesan: fmt.Sprintf("The file cannot be read as a CSV file (problem in row %d). Check for an unclosed quote mark (\").", pe.Line)}
	}
	return fmt.Errorf("services: reading the uploaded file: %w", err)
}

// acuanCSV membaca `PremiumListDetail(1)`.
func (l *Layanan) acuanCSV(ctx context.Context, tx *db.Tx, id string) (models.AcuanCSV, error) {
	a, ada, err := l.gudang.AcuanCSV(ctx, tx, id)
	if err != nil {
		return models.AcuanCSV{}, err
	}
	if !ada {
		return models.AcuanCSV{}, ErrCSVTanpaAcuan
	}
	return a, nil
}

// SiapCSV - gerbang kasus unggahan (identitas, kasus terbuka Perubahan Data
// belum terkunci, ada peserta acuan). Handler memanggilnya SEBELUM berkas
// diterima, supaya berkas peserta tidak pernah ditampung untuk kasus yang menolaknya.
func (l *Layanan) SiapCSV(ctx context.Context, p inti.Pelaku, id string) (models.AcuanCSV, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.AcuanCSV{}, err
	}
	k, err := l.muatKasus(ctx, id)
	if err != nil {
		return models.AcuanCSV{}, err
	}
	if err := periksaKasusCSV(k); err != nil {
		return models.AcuanCSV{}, err
	}
	a, err := l.acuanCSV(ctx, nil, id)
	if err != nil {
		return models.AcuanCSV{}, err
	}
	a.Tipe = k.Kepala["TYPE"]
	return acuanPolisCSV(a, k)
}

// acuanPolisCSV - PLAN / POLICY_HOLDER peserta baru dari polis (keputusan work owner 10-10-2026, seperti PremiumList
// Life): kepala kasus, lalu peserta Old pertama; keduanya kosong → `ErrMasukanTidakSah` berpesan (422).
func acuanPolisCSV(a models.AcuanCSV, k models.Kasus) (models.AcuanCSV, error) {
	a, err := models.NilaiPolisCSV(a, k.Kepala["BUSINESS_NAME"], k.Kepala["POLICY_HOLDER_NAME"])
	if err != nil {
		return models.AcuanCSV{}, galatBerkasCSV(err)
	}
	return a, nil
}

// PeriksaCSV - `Upload CSV`: mengurai dan memvalidasi, nol penulisan (tinjauan
// sebelum `Add CSV Data`).
func (l *Layanan) PeriksaCSV(ctx context.Context, p inti.Pelaku, id string, berkas io.Reader) (HasilPeriksaCSV, error) {
	acuan, err := l.SiapCSV(ctx, p, id)
	if err != nil {
		return HasilPeriksaCSV{}, err
	}
	return telusuriCSV(berkas, acuan, nil)
}

// TambahCSV - `Add CSV Data` (`SaveCSVEDMLife`), satu transaksi: baris `New`
// lama dibuang (2.1), seluruh baris diperiksa ulang lalu disisip (4.1/4.3),
// jejak. Kasus yang sudah diproses ditolak (`ErrSudahDiproses`) - rekap tidak
// pernah dihitung di sini (08-10-2026; dulu dihitung ulang, R31).
func (l *Layanan) TambahCSV(ctx context.Context, p inti.Pelaku, id string, berkas io.ReadSeeker) (HasilTambahCSV, error) {
	cek, err := l.PeriksaCSV(ctx, p, id, berkas)
	if err != nil {
		return HasilTambahCSV{}, err
	}
	if cek.Ditolak > 0 {
		return HasilTambahCSV{}, GalatCSV{cek}
	}
	if _, err := berkas.Seek(0, io.SeekStart); err != nil {
		return HasilTambahCSV{}, fmt.Errorf("services: rewinding the CSV file: %w", err)
	}
	var hasil HasilTambahCSV
	err = l.tx(ctx, func(tx *db.Tx) error {
		k, err := l.gudang.AmbilKasus(ctx, tx, id, true)
		if errors.Is(err, repository.ErrTidakAda) {
			return fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
		}
		if err != nil {
			return err
		}
		if k.Cacah, err = l.gudang.CacahPeserta(ctx, tx, id); err != nil {
			return err
		}
		if k.SudahSimpan, err = l.gudang.AdaRekap(ctx, tx, id); err != nil {
			return err
		}
		if err := periksaKasusCSV(k); err != nil {
			return err
		}
		acuan, err := l.acuanCSV(ctx, tx, id)
		if err != nil {
			return err
		}
		acuan.Tipe = k.Kepala["TYPE"]
		if acuan, err = acuanPolisCSV(acuan, k); err != nil {
			return err
		}
		if hasil.Dibuang, err = l.gudang.HapusPesertaBaru(ctx, tx, id); err != nil {
			return err
		}
		h, err := telusuriCSV(berkas, acuan, func(b []models.BarisCSV) error {
			n, err := l.gudang.SisipPesertaCSV(ctx, tx, id, k.PLNumber, b)
			hasil.Disimpan += n
			return err
		})
		if err != nil {
			return err
		}
		if h.Ditolak > 0 {
			return GalatCSV{h} // berkas berubah di antara kedua lintasan
		}
		if hasil.Rekap, err = l.gudang.RekapKasus(ctx, tx, id); err != nil {
			return err
		}
		return l.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			KlaimID: id, Dari: models.TahapInputEDMLife, Ke: models.TahapInputEDMLife, AkunID: p.AkunID, Waktu: l.jam(),
			Komentar: models.AksiTambahCSV,
		})
	})
	if err != nil {
		return HasilTambahCSV{}, err
	}
	if hasil.Rekap == nil {
		hasil.Rekap = []map[string]string{}
	}
	return hasil, nil
}

// BerkasTemplat - templat CSV unduhan.
type BerkasTemplat struct {
	Nama string
	Isi  []byte
}

// TemplatCSV - `Download Template` (keputusan work owner 10-10-2026): judul 72 properti (kolom wajib Type kasus
// paling depan) + satu baris contoh dari peserta Old pertama (`models.BarisContohCSV`), UTF-8 berpemisah koma.
// Gerbang sama dengan unggahan (`SiapCSV`): kasus tertutup / sudah diproses / bukan Perubahan Data / terkunci ditolak.
func (l *Layanan) TemplatCSV(ctx context.Context, p inti.Pelaku, id string) (BerkasTemplat, error) {
	acuan, err := l.SiapCSV(ctx, p, id)
	if err != nil {
		return BerkasTemplat{}, err
	}
	contoh, ada, err := l.gudang.ContohPesertaCSV(ctx, nil, id)
	if err != nil {
		return BerkasTemplat{}, err
	}
	if !ada {
		return BerkasTemplat{}, ErrCSVTanpaAcuan
	}
	judul := models.JudulTemplatCSV(acuan.Tipe)
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	if err := w.Write(judul); err != nil {
		return BerkasTemplat{}, err
	}
	if err := w.Write(models.BarisContohCSV(judul, contoh)); err != nil {
		return BerkasTemplat{}, err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return BerkasTemplat{}, err
	}
	return BerkasTemplat{Nama: "Template_EDM_" + id + ".csv", Isi: b.Bytes()}, nil
}
