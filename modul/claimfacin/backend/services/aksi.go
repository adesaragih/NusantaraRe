package services

// Untuk apa berkas ini: AKSI LAYAR. Setiap aksi (refresh ber-activity sebuah sel, tombol, tombol Choose pop-up) berjalan
// dalam SATU transaksi:
//
//	kunci baris kasus (tahap belum berubah, belum tertutup) -> pelaku pemegang assignment
//	-> halaman TERSIMPAN dimuat (+ halaman polis) + pra-proses + turunan -> tata dievaluasi
//	-> aksi harus tampil dan aktif di tata KONTEKS-nya (layar utama, panel baris, atau modal) -> kiriman layar digabung
//	   HANYA untuk jalur terbuka -> port activity dijalankan -> turunan dihitung ulang -> halaman disimpan
//	   -> InsertProgressClaim (pra-proses assignment, idempoten)
//
// ⚠️ PENYIMPANGAN SADAR (PARITAS, pola Claim Prop): Pega menyimpan clipboard hanya pada Save / Obj-Save; di sini setiap
// aksi menyimpan halaman (write-through). Aksi yang keluar karena pesan validasi (`GalatValidasi`) dibatalkan
// seluruhnya; layar tetap menampilkan isian dan pesannya.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// PermintaanAksi - satu aksi layar.
type PermintaanAksi struct {
	// Aksi - nama aksi di tata (nama activity tanpa akhiran, atau ID tombol); ":param" = parameter activity.
	Aksi string `json:"aksi"`
	// Konteks - kunci panel baris (`models.KunciPanel`) atau modal ("pla:1") tempat aksi berada; "" = layar utama.
	Konteks string `json:"konteks,omitempty"`
	// Indeks - baris grid (1..n) di dalam konteks untuk sel / tombol baris; 0 = bukan sel grid.
	Indeks int `json:"indeks"`
	// Param - pilihan pop-up / autocomplete (dibaca ULANG di server) atau parameter activity.
	Param string `json:"param"`
	// Tahap - tahap kasus yang dilihat layar (penjaga balapan: tahap berubah = 409).
	Tahap string `json:"tahap"`
	// Masukan - nilai medan terbuka layar (jalur -> teks).
	Masukan map[string]string `json:"masukan"`
	// Mode - penanda mode layar terakhir (`Layar.Mode`, models.ModeLayar).
	Mode map[string]string `json:"mode,omitempty"`
}

// jalanAksi - konteks satu aksi untuk penanganannya.
type jalanAksi struct {
	l     *Layanan
	ctx   context.Context
	tx    *db.Tx
	k     *models.Konteks
	kasus models.Kasus
	p     inti.Pelaku
	h     *models.Halaman
	r     PermintaanAksi
	// o, i, e - objek / item / baris (estimasi / adjustment) yang dituju aksi (1..n; 0 = tidak ada).
	o, i, e int
	// info / bukaModal - pemberitahuan dan modal yang dibuka layar sesudah aksi.
	info, bukaModal string
	// selesai - kasus sudah ditutup / dipindah oleh aksi ini (halaman disimpan sebelum pemindahan).
	selesai bool
	// tanpaProgres - InsertProgressClaim tidak dijalankan sesudah aksi (pemindahan tahap menuliskannya sendiri).
	tanpaProgres bool
	// dokumen - bahan PDF akseptasi yang diproses SESUDAH aksi tersimpan (SaveAcceptation 12).
	dokumen *dokumenTunda
}

type penanganAksi func(j *jalanAksi) error

// Aksi menjalankan satu aksi dan mengembalikan layar sesudahnya.
func (l *Layanan) Aksi(ctx context.Context, p inti.Pelaku, id string, r PermintaanAksi) (*Layar, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	nama, param := r.Aksi, ""
	if i := strings.Index(nama, ":"); i >= 0 { // "SetEditCatastrope:Edit"
		nama, param = nama[:i], nama[i+1:]
	}
	f, ada := penangan[nama]
	if !ada {
		return nil, fmt.Errorf("%w: aksi %q", ErrPermintaanTidakSah, r.Aksi)
	}
	var hasil *Layar
	var gagal *GalatValidasi
	var dokumen *dokumenTunda
	err := l.g.Transaksi(ctx, func(tx *db.Tx) error {
		kasus, err := l.g.Keadaan(ctx, tx, id)
		if err != nil {
			return err
		}
		if kasus.Tertutup() {
			return ErrKasusTertutup
		}
		if r.Tahap != "" && r.Tahap != kasus.Tahap {
			return ErrTahapBerubah
		}
		if !Pemegang(p, kasus) {
			return ErrBukanPemegang
		}
		if kasus, err = l.g.KunciKasus(ctx, tx, id, kasus.Tahap); err != nil {
			return err
		}
		_, h, err := l.muat(ctx, tx, id)
		if err != nil {
			return err
		}
		models.PasangMode(h, r.Mode)
		kt, err := l.konteks(ctx, p, kasus, l.jam())
		if err != nil {
			return err
		}
		if err := l.siapkan(ctx, kt, kasus, h); err != nil {
			return err
		}
		if err := l.turunkan(ctx, kasus, h); err != nil {
			return err
		}
		h.BersihkanPesan() // pesan pra-proses tidak dibawa ke hasil aksi
		j := &jalanAksi{l: l, ctx: ctx, tx: tx, k: kt, kasus: kasus, p: p, h: h, r: r}
		if f, ada := praTata[nama]; ada {
			if err := f(j); err != nil {
				return err
			}
		}
		utama, panel, modal := tataKasus(kasus, h, false)
		konteks := utama
		if r.Konteks != "" {
			t, ada := panel[r.Konteks]
			if !ada {
				t, ada = modal[r.Konteks]
			}
			if !ada {
				return fmt.Errorf("%w: konteks %q", ErrPermintaanTidakSah, r.Konteks)
			}
			konteks = t
		}
		if !aksiPascaGabung[nama] && !aksiTerbuka(konteks, utama, nama, r.Aksi, r.Indeks) {
			return fmt.Errorf("%w: %q", ErrAksiTertutup, r.Aksi)
		}
		models.GabungMasukan(h, r.Masukan, models.MedanTerbuka(semuaTata(utama, panel, modal)))
		if param != "" && r.Param == "" {
			r.Param = param
		}
		j.r = r
		if err := j.sasaran(); err != nil {
			return err
		}
		// Aksi berpesan validasi dibatalkan seluruhnya: layar validasinya dibangun dari halaman SESUDAH isian digabung
		// dan SEBELUM penangan mengubahnya (mis. baris yang dihapus), ditambah pesannya - layar = keadaan tersimpan +
		// isian layar (temuan review 10-10-2026: dahulu dari halaman yang sudah diubah penangan).
		sebelum := h.Salin()
		if err := f(j); err != nil {
			var gv *GalatValidasi
			if errors.As(err, &gv) {
				gagal = gv
				sebelum.Pesan = h.Pesan
				models.HitungTurunan(sebelum)
				gv.Layar = l.layar(kasus, sebelum, true)
			}
			return err
		}
		if !j.selesai {
			if err := l.g.SimpanHalaman(ctx, tx, id, h); err != nil {
				return err
			}
			if err := l.g.SentuhKasus(ctx, tx, id, kt.Sekarang); err != nil {
				return err
			}
			if !j.tanpaProgres {
				if err := l.tulisProgres(ctx, tx, kt, id, h, models.ParamProgres{}); err != nil {
					return err
				}
			}
		}
		pesan := h.Pesan
		k2, h2, err := l.muat(ctx, tx, id)
		if err != nil {
			return err
		}
		if Pemegang(p, k2) {
			kt2, err := l.konteks(ctx, p, k2, kt.Sekarang)
			if err != nil {
				return err
			}
			if err := l.siapkan(ctx, kt2, k2, h2); err != nil {
				return err
			}
		}
		models.BawaSementara(h, h2)
		if err := l.turunkan(ctx, k2, h2); err != nil {
			return err
		}
		h2.BersihkanPesan()
		h2.Pesan = pesan
		hasil = l.layar(k2, h2, Pemegang(p, k2))
		hasil.Info, hasil.BukaModal = j.info, j.bukaModal
		dokumen = j.dokumen
		return nil
	})
	if gagal != nil {
		return nil, gagal
	}
	if errors.Is(err, repository.ErrAdjustmentBerkomite) { // jagaKomite: konflik, bukan galat server
		return nil, fmt.Errorf("%w: %v", ErrAdjustmentDiKomite, err)
	}
	if err != nil {
		return nil, err
	}
	if dokumen != nil { // akseptasi sudah tersimpan: pemutusan klien tidak membatalkan dokumennya
		if err := l.simpanDokumenAkseptasi(context.WithoutCancel(ctx), *dokumen); err != nil {
			log.Printf("claimfacin: dokumen akseptasi kasus %s: %v", id, err)
			hasil.Info = models.PesanDokumenGagal
		}
	}
	return hasil, nil
}

// aksiInduk - aksi pop-up (pilihan grid / tombol harness) yang terbuka bila salah satu tombol PEMBUKA pop-up-nya terbuka
// di layar utama (grid pop-up diisi pilihan server, bukan halaman kasus).
var aksiInduk = map[string][]string{
	"SearchPolis":        {"BukaPolis"},
	"CopyNB":             {"BukaPolis"},
	"SetInputParam":      {"BukaPolis"},
	"GetNameCauseofLoss": {"BukaSebab"},
	"SetCatastrope":      {"CatastrofeList"},
	"SaveCatasrtope":     {"CatastrofeList"},
	"InputCatastrope":    {"CatastrofeList"},
}

// aksiTerbuka - aksi (nama, atau nama berparameter apa adanya) terbuka di konteksnya, atau salah satu pembuka
// pop-up-nya terbuka di layar utama.
func aksiTerbuka(konteks, utama []models.Tata, nama, utuh string, indeks int) bool {
	if models.AksiTerbuka(konteks, utuh, indeks) || models.AksiTerbuka(konteks, nama, indeks) {
		return true
	}
	for _, induk := range aksiInduk[nama] {
		if models.AksiTerbuka(utama, induk, 0) {
			return true
		}
	}
	return false
}

// pembukaTerbuka - tombol PEMBUKA pop-up (aksi `pembuka` di tata `konteks` - "" layar utama, selainnya kunci panel - baris
// `indeks`) masih terbuka pada halaman sekarang. Submit pop-up memeriksanya karena modal pop-up dievaluasi untuk setiap
// baris, bukan hanya baris yang tombol pembukanya aktif (temuan review 10-10-2026).
func pembukaTerbuka(j *jalanAksi, konteks, pembuka string, indeks int) bool {
	utama, panel, _ := tataKasus(j.kasus, j.h, false)
	ts := utama
	if konteks != "" {
		ts = panel[konteks]
	}
	return models.AksiTerbuka(ts, pembuka, indeks)
}

// sasaran membaca objek / item / baris yang dituju aksi dari konteks dan indeks permintaan.
//
//	""               layar utama: Indeks = baris grid objek (tombol CFS / Print PLA / Outstanding)
//	"est:..(o)"      panel objek: Indeks = baris grid item
//	"estitem:..(o)..(i)" panel item: Indeks = baris grid estimasi
//	"pla:o"          modal Print PLA objek o
//	"dla:o"          modal Print DLA objek o
//	"adj..." / "cedant..."  panel / modal layar Input Adjustment (adjustment.go `konteksAdj`)
func (j *jalanAksi) sasaran() error {
	k := j.r.Konteks
	if n, ok := modalDLA(k); ok {
		j.o = n
		return nil
	}
	if o, i, a, ok := konteksAdj(k); ok {
		switch {
		case i == 0: // panel objek: Indeks = baris item
			j.o, j.i = o, j.r.Indeks
		case a == 0: // panel item: Indeks = baris Adjustment
			j.o, j.i, j.e = o, i, j.r.Indeks
		default: // panel adjustment / pop-up Cedant: Indeks = baris grid di dalamnya
			j.o, j.i, j.e = o, i, a
		}
		return nil
	}
	switch {
	case k == "":
		j.o = j.r.Indeks
	case strings.HasPrefix(k, AwalanModalPLA+":"):
		n, err := strconv.Atoi(strings.TrimPrefix(k, AwalanModalPLA+":"))
		if err != nil {
			return fmt.Errorf("%w: konteks %q", ErrPermintaanTidakSah, k)
		}
		j.o = n
	case k == ModalPilihPolis || k == ModalProtectDOL || k == ModalTutup || k == ModalTolak:
	default:
		prefiks, idx, ok := models.UraiPanel(k)
		if !ok {
			return fmt.Errorf("%w: konteks %q", ErrPermintaanTidakSah, k)
		}
		switch prefiks {
		case models.PanelObjekEst:
			j.o, j.i = idx[0], j.r.Indeks
		case models.PanelItemEst:
			if len(idx) < 2 {
				return fmt.Errorf("%w: konteks %q", ErrPermintaanTidakSah, k)
			}
			j.o, j.i, j.e = idx[0], idx[1], j.r.Indeks
		}
	}
	return nil
}

// validasi mengubah pesan halaman menjadi GalatValidasi (aksi dibatalkan).
func validasi(h *models.Halaman) error {
	if !h.AdaPesan() {
		return nil
	}
	return &GalatValidasi{Pesan: h.SemuaPesan()}
}

// butuh - aksi yang menuntut objek / item / baris sasaran.
func (j *jalanAksi) butuh(o, i, e bool) error {
	if (o && j.o < 1) || (i && j.i < 1) || (e && j.e < 1) {
		return fmt.Errorf("%w: baris sasaran aksi %q", ErrPermintaanTidakSah, j.r.Aksi)
	}
	return nil
}

// ---------------------------------------------------------------- penangan sederhana

func halamanSaja(f func(j *jalanAksi)) penanganAksi {
	return func(j *jalanAksi) error { f(j); return nil }
}

func dariModel(f func(k *models.Konteks, h *models.Halaman) error) penanganAksi {
	return func(j *jalanAksi) error { return f(j.k, j.h) }
}

// diItem - aksi panel item (o, i).
func diItem(f func(j *jalanAksi) error) penanganAksi {
	return func(j *jalanAksi) error {
		if err := j.butuh(true, true, false); err != nil {
			return err
		}
		return f(j)
	}
}

// diBaris - aksi sel / tombol grid estimasi (o, i, e).
func diBaris(f func(j *jalanAksi) error) penanganAksi {
	return func(j *jalanAksi) error {
		if err := j.butuh(true, true, true); err != nil {
			return err
		}
		return f(j)
	}
}

// bukaModal - tombol yang hanya membuka pop-up.
func bukaModal(m func(j *jalanAksi) string) penanganAksi {
	return func(j *jalanAksi) error { j.bukaModal = m(j); return nil }
}

// penangan - aksi -> port activity. Nama = nama activity tanpa akhiran `_Act`/`_act`/`_ACT`.
var penangan map[string]penanganAksi

func init() {
	penangan = map[string]penanganAksi{
		// registrasi (InputRegisterDetail - juga tab View Registration / Registration layar sesudahnya)
		"BukaPolis":                      bukaModal(func(*jalanAksi) string { return ModalPilihPolis }),
		"SearchPolis":                    halamanSaja(func(*jalanAksi) {}), // hasil = pilihan "polis" (SearchPolis_act)
		"CopyNB":                         aksiPilihPolis,
		"SetInputParam":                  halamanSaja(func(j *jalanAksi) { models.SubmitPilihPolis(j.k, j.h) }),
		"CheckDate":                      halamanSaja(func(j *jalanAksi) { models.CheckDate(j.k, j.h) }),
		"CheckDateReport":                halamanSaja(func(j *jalanAksi) { models.CheckDateReport(j.k, j.h) }),
		"CheckDateReceived":              halamanSaja(func(j *jalanAksi) { models.CheckDateReceived(j.k, j.h) }),
		"GetReportStatus":                dariModel(models.GetReportStatus),
		"GetCeding":                      aksiPilihCeding,
		"SetConsultant":                  dariModel(models.SetConsultant),
		"SetAdjsuter":                    dariModel(models.SetAdjuster),
		"SetTempLocation":                halamanSaja(func(j *jalanAksi) { models.SetTempLocation(j.h) }),
		"InputKodePos1":                  dariModel(models.InputKodePos),
		"PilihNegara":                    aksiPilihWilayah(models.SumberNegara),
		"PilihProvinsi":                  aksiPilihWilayah(models.SumberProvinsi),
		"PilihKota":                      aksiPilihWilayah(models.SumberKota),
		"PilihDistrik":                   aksiPilihWilayah(models.SumberDistrik),
		"PilihRW":                        aksiPilihWilayah(models.SumberRW),
		"InputCurrencyValueAct_Register": dariModel(models.InputCurrencyValue),
		"CheckDoubleClaim":               aksiCekKembar,
		"CheckListEstimasi":              halamanSaja(func(j *jalanAksi) { models.CheckListEstimasi(j.h) }),
		"SetDefNonCatastrope":            halamanSaja(func(j *jalanAksi) { models.SetDefNonCatastrope(j.h) }),
		"SetEditCatastrope":              halamanSaja(func(j *jalanAksi) { models.SetEditCatastrope(j.h, j.r.Param) }),
		"BukaKatastrofe":                 halamanSaja(func(*jalanAksi) {}),
		"InputCatastrope":                halamanSaja(func(j *jalanAksi) { models.InputCatastrope(j.h, j.r.Param) }),
		"SetCatastrope":                  aksiPilihKatastrofe,
		"SaveCatasrtope":                 aksiSimpanKatastrofe,
		"BukaSebab":                      halamanSaja(func(*jalanAksi) {}),
		"GetNameCauseofLoss":             aksiPilihSebab,
		"Simpan":                         halamanSaja(func(*jalanAksi) {}),
		"BukaProtectDOL":                 bukaModal(func(*jalanAksi) string { return ModalProtectDOL }),
		"Submit":                         aksiSubmitRegister,
		// estimasi - panel objek (grid item)
		"TambahItem": func(j *jalanAksi) error {
			if err := j.butuh(true, false, false); err != nil {
				return err
			}
			_, err := models.TambahItem(j.h, j.o)
			return err
		},
		"HapusItem":          diItem(func(j *jalanAksi) error { return models.HapusItem(j.k, j.h, j.o, j.i) }),
		"SetObjectItem":      diItem(func(j *jalanAksi) error { return models.PilihItemProperti(j.k, j.h, j.o, j.i, j.r.Param) }),
		"PilihOkupasi":       diItem(func(j *jalanAksi) error { return models.PilihOkupasi(j.h, j.o, j.i, j.r.Param) }),
		"PilihCoverageFire":  diItem(func(j *jalanAksi) error { return models.PilihCoverageFire(j.h, j.o, j.i, j.r.Param) }),
		"PilihAneka":         diItem(func(j *jalanAksi) error { return models.PilihAneka(j.k, j.h, j.o, j.i, j.r.Param) }),
		"PilihCoverageAneka": diItem(func(j *jalanAksi) error { return models.PilihCoverageAneka(j.k, j.h, j.o, j.i, j.r.Param) }),
		"PilihCoverageObjek": diItem(func(j *jalanAksi) error { return models.PilihCoverageObjek(j.k, j.h, j.o, j.i, j.r.Param) }),
		// estimasi - panel item (grid estimasi, deductible)
		"ValidateInputEstimate": diItem(func(j *jalanAksi) error { return models.TambahEstimasi(j.k, j.h, j.o, j.i) }),
		"CheckEstimateValue":    diItem(func(j *jalanAksi) error { return models.CheckEstimateValue(j.k, j.h, j.o, j.i) }),
		"CountTSI":              diItem(func(j *jalanAksi) error { return models.CountTSI(j.h, j.o, j.i) }),
		"SetConvertValueKurs_Estimation": diItem(func(j *jalanAksi) error {
			return models.UbahMataUangEstimasi(j.k, j.h, j.o, j.i)
		}),
		"ProtectionDate": diBaris(func(j *jalanAksi) error {
			models.ProtectionDate(j.k, j.h, j.o, j.i, j.e)
			return nil
		}),
		"DeleteValueEstimation": diBaris(func(j *jalanAksi) error { return models.HapusEstimasi(j.k, j.h, j.o, j.i, j.e) }),
		"BukaRetro":             halamanSaja(func(*jalanAksi) {}),
		"BukaRetroList":         halamanSaja(func(*jalanAksi) {}),
		// estimasi - layar utama
		"CLaimFaceSheet":  aksiCFS,
		"BukaPLA":         aksiBukaPLA,
		"GeneratePLA":     aksiPLA,
		"BukaOutstanding": halamanSaja(func(*jalanAksi) {}), // isi = pilihan "outstanding" (GetAllData_Act)
		"SetDisable":      aksiKirimPIC,
		"BackToRegister":  aksiKembaliRegister,
	}
	for nama, f := range penanganAdjustment() {
		penangan[nama] = f
	}
}
