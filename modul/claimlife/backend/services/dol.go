package services

// Validasi Date of Loss per Type, dan penurunan jenis klaim - tiket 06.
//
// Untuk apa berkas ini: menolak tanggal kejadian yang berada di luar jendela
// valuasi polis, sehingga klaim yang tidak tertanggung tidak pernah masuk ke
// siklus; dan menurunkan jenis klaim dari kode produk.
//
// Dibaca sesudah: dokumen.go.
//
// ⛔ Seluruh isi berkas ini MURNI. Jendela valuasi sudah DISALIN ke peserta
// klaim saat pendaftaran (tiket 02), justru supaya validasi ini tidak perlu
// menyentuh tabel 66,8 juta baris sama sekali.
//
// Istilah:
//   - DOL            : Date of Loss, tanggal kejadian yang diklaim.
//   - jendela valuasi: rentang tanggal polis yang menanggung kejadian.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// Nilai Type yang punya cabang di `ValidasiDOL_Act`. Tidak ada nilai lain.
//
// ⚠️ `[terbuka - OQ-020]` Arti QP/QR belum dijawab siapa pun; TP = Payable dan
// TR = Receivable `[keputusan work owner, ADR-U-0012]`. Perilaku kedua cabang
// sudah terbaca penuh - hanya namanya yang belum, dan nama tidak dipakai
// menghitung apa pun di sini.
const (
	TypeQR = "QR"
	TypeQP = "QP"
	TypeTR = "TR"
	TypeTP = "TP"
)

// ⛔ JENDELA DOL — KEPUTUSAN WORK OWNER 11-10-2026, MENYIMPANG SADAR dari Pega: "tanggal awal valuation masih
// termasuk". Jendelanya INKLUSIF di kedua ujung dan SAMA untuk semua jalur (Save grid, Edit Date, Save to RNM 11.7-11.8)
// dan kedua cabang:
//
//	QR/QP : GROSS_VALUATION_BEGIN <= DOL <= GROSS_VALUATION_EXPIRED
//	TR/TP : RETRO_VALUATION_BEGIN <= DOL <= RETRO_VALUATION_EXPIRED
//
// dibandingkan per TANGGAL kalender (jam diabaikan). Frontend `galatDOL` (RincianPeserta.tsx) cermin aturan ini.
//
// Yang DIBUANG (catatan sejarah, lihat PARITAS 11-10-2026 lanjutan 9):
//   - ValidasiDOL_Act b698 menggeser DOL retro SATU HARI (`@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)`, posisi keempat =
//     hari - jangkar angka 7 = WIB di LoadDataPeserta_Act), sedang langkah 11.8 Save to RNM (b4464) tidak menggeser:
//     tanggal yang lolos Save grid dapat ditolak Save to RNM. Kini tanpa geser di mana pun.
//   - `@CompareDates` ketat (dugaan) membuat jendela gross (BEGIN, EXPIRED] dan retro [BEGIN, EXPIRED). Kini [BEGIN, EXPIRED].

var (
	// ErrDOLTidakSah - pesannya PERSIS `local.errmsg` di XML, tanpa tambahan.
	// Peserta mana yang bermasalah dibawa GalatDOL di medannya sendiri.
	ErrDOLTidakSah = errors.New("Invalid DOL")
	// ErrTypeTidakDikenal - Type di luar keempat nilai di atas.
	ErrTypeTidakDikenal = errors.New("services: Type klaim tidak dikenal")
	// ErrDOLKosong - tanggal kejadian tidak diisi sama sekali.
	ErrDOLKosong = errors.New("services: tanggal kejadian belum diisi")
	// ErrBusinessCodeTidakDikenal - kode produk di luar daftar acuan.
	ErrBusinessCodeTidakDikenal = errors.New("services: BusinessCode tidak dikenal")
)

// GalatDOL adalah DOL yang tidak sah, beserta peserta yang memilikinya.
//
// ⛔ `Error()` mengembalikan kalimat XML APA ADANYA. Pengenal pesertanya hidup
// di medan tersendiri, bukan di dalam kalimat: pesan galat adalah logika
// bisnis, dan menambahinya berarti sistem baru berbicara dengan kalimat yang
// tidak pernah ada di sistem lama. Pola yang sama dengan gerbang dokumen
// tiket 03.
type GalatDOL struct {
	PesertaID string
	// Nomor urut peserta di klaim (1..n), diisi pemanggil; 0 = tidak diketahui.
	Nomor int
	// Awal / Akhir - tanggal valuasi jendela yang dipakai (dd/mm/yyyy), untuk kalimat layar.
	Awal, Akhir string
}

// KalimatDOLDiLuar - kalimat DOL di luar jendela valuasi yang DITAMPILKAN (keputusan work owner 11-10-2026: "date of loss
// harus ada di periode valuation date ... jika di luar tidak bisa lanjut dan ada pesan"). SAMA dengan frontend
// `galatDOL` (RincianPeserta.tsx). "Invalid DOL" (`local.errmsg` XML) tetap menjadi KODE (Error(), medan `kode`).
func KalimatDOLDiLuar(awal, akhir string) string {
	return fmt.Sprintf("Date of Loss must be within the valuation period (%s – %s).", awal, akhir)
}

// Kalimat - kalimat layar yang menyebut peserta: "Insured number N: Date of Loss must be within …".
func (e GalatDOL) Kalimat() string {
	k := KalimatDOLDiLuar(e.Awal, e.Akhir)
	if e.Nomor > 0 {
		return fmt.Sprintf("Insured number %d: %s", e.Nomor, k)
	}
	return k
}

// BeriNomorDOL - menempelkan nomor urut peserta pada GalatDOL; galat lain dikembalikan apa adanya.
func BeriNomorDOL(err error, nomor int) error {
	var g GalatDOL
	if errors.As(err, &g) {
		g.Nomor = nomor
		return g
	}
	return err
}

// tanggalLayar - dd/mm/yyyy, bentuk tanggal layar (tanggalTampil frontend).
func tanggalLayar(t time.Time) string { return t.Format("02/01/2006") }

// Error menulis kalimat XML, tanpa tambahan apa pun.
func (e GalatDOL) Error() string { return ErrDOLTidakSah.Error() }

// Unwrap membuat `errors.Is(err, ErrDOLTidakSah)` bernilai benar tanpa
// kalimatnya ikut berubah - Unwrap tidak menyentuh Error().
func (e GalatDOL) Unwrap() error { return ErrDOLTidakSah }

// ValidasiDOL menolak Date of Loss di luar jendela valuasi peserta.
//
// Meniru `Claim Life/Activity/ValidasiDOL_Act.xml` langkah 2, 3, dan 4:
//
//	langkah 2  Type QR atau QP -> jendela GROSS_VALUATION_*
//	langkah 3  Type TR atau TP -> jendela RETROCESSION_VALUATION_* (geser +1 hari b698 DIBUANG, keputusan WO 11-10-2026)
//	langkah 4  galat bila `local.Begin==false || local.Expired==true`
//
// ⛔ Type di luar keempat nilai itu adalah GALAT, bukan izin lewat. Di Pega
// kedua precondition tidak terpenuhi sehingga activity-nya diam saja dan
// tanggal apa pun lolos; diam seperti itu adalah lubang, bukan aturan.
func ValidasiDOL(tipe string, dol time.Time, p models.Peserta) error {
	var gross bool

	if !TypeDikenal(tipe) {
		return fmt.Errorf("%w: %q; yang dikenal hanya %s, %s, %s, dan %s",
			ErrTypeTidakDikenal, tipe, TypeQR, TypeQP, TypeTR, TypeTP)
	}
	switch tipe {
	case TypeQR, TypeQP:
		gross = true
	case TypeTR, TypeTP:
		gross = false
	}

	// ⛔ DOL kosong ditolak TERPISAH. Waktu nol adalah tahun 1 Masehi, yang
	// pasti di luar jendela mana pun - jadi tanpa penjaga ini pengguna
	// mendapat "Invalid DOL" untuk tanggal yang tidak pernah ia isi, dan
	// pesan itu menyesatkan ke arah yang salah.
	if dol.IsZero() {
		return fmt.Errorf("%w: peserta %s", ErrDOLKosong, p.ID)
	}

	awal, akhir, err := p.JendelaValuasi(gross)
	if err != nil {
		return err
	}

	// Keputusan work owner 11-10-2026: [BEGIN, EXPIRED] inklusif, tanpa geser (lihat catatan di atas konstanta Type).
	if !DOLDalamJendela(dol, awal, akhir) {
		return GalatDOL{PesertaID: p.ID, Awal: tanggalLayar(awal), Akhir: tanggalLayar(akhir)}
	}
	return nil
}

// DOLDalamJendela - SATU-SATUNYA aturan jendela DOL (Save grid, Edit Date, Save to RNM): awal <= DOL <= akhir, per
// tanggal kalender, inklusif di kedua ujung (keputusan work owner 11-10-2026).
func DOLDalamJendela(dol, awal, akhir time.Time) bool {
	d, a, z := hariSaja(dol), hariSaja(awal), hariSaja(akhir)
	return !d.Before(a) && !d.After(z)
}

// hariSaja - tanggal kalender (jam, menit, detik dibuang) pada zona waktu nilainya sendiri.
func hariSaja(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// ContentNoteDari menurunkan jenis klaim dari kode produk.
//
// Aturannya di sini, datanya di `models` - kode yang tidak dikenal adalah
// keputusan lapisan ini, bukan sifat tabelnya.
//
// ⛔ Kode asing menjadi GALAT, bukan jenis kosong. Jenis klaim menggerbangi
// jalur Medical Check dan Claim Analis di hilir; menebaknya berarti mengarahkan
// klaim ke jalur yang belum tentu miliknya.
func ContentNoteDari(kodeBisnis string) (string, error) {
	note, ada := models.ContentNoteUntuk(kodeBisnis)
	if !ada {
		return "", fmt.Errorf("%w: %q", ErrBusinessCodeTidakDikenal, kodeBisnis)
	}
	return note, nil
}

// TanggalKejadian adalah layanan pengisian Date of Loss seorang peserta.
//
// ⭐ Inilah jalur yang membuat ValidasiDOL benar-benar menolak sesuatu. Di Pega
// validasinya dipicu dua Section saat tanggal diubah di layar
// (`EditDateClaimLife_Section` dan `ClaimLifeDetailGCNM`, lewat
// `<pyActivity>ValidasiDOL_Act</pyActivity>`); di sini pemicunya satu pintu
// HTTP yang menyetel tanggal itu.
type TanggalKejadian struct {
	svc *Service
}

// TanggalKejadian menyusun layanan itu.
func (s *Service) TanggalKejadian() *TanggalKejadian { return &TanggalKejadian{svc: s} }

// Set memvalidasi lalu menyimpan Date of Loss seorang peserta.
//
// Urutannya disengaja: Type dan peserta dibaca dari basis data lebih dulu,
// validasi berjalan atas nilai yang TERSIMPAN, dan penulisan hanya terjadi
// bila validasinya lolos. Memvalidasi terhadap nilai yang dikirim klien berarti
// mempercayai klien untuk menyatakan jendela valuasinya sendiri.
func (t *TanggalKejadian) Set(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string, dol time.Time) error {

	// ⛔ FAIL-CLOSED atas pelaku anonim, alasan yang sama dengan pendaftaran:
	// tanpa identitas, perubahan tanggal kejadian tidak dapat ditelusuri.
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(pesertaID) == "" {
		return fmt.Errorf("%w: pengenal klaim dan peserta wajib diisi", galat.ErrPermintaanTidakSah)
	}
	// ⛔ BUTIR bj `[DIPUTUSKAN 28-09-2026, veto work owner]` - PERUBAHAN AUTHZ,
	// dicatat bertanggal di tiket 07. Gerbang DOL kini SAMA dengan tiga tanggal
	// lainnya menurut XML: isian DOL `EditDateClaimLife_Section` baca-saja bila
	// `pyPosition!='ReasLifeAdmin'` (b1000). Sebelum ini siapa pun yang
	// beridentitas dapat mengubah DOL pada kasus terbuka di tahap mana pun.
	if err := inti.WajibPeran(pelaku, inti.PeranAdmin); err != nil {
		return err
	}
	if !t.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}

	// ⛔ BUTIR bb: kasus yang sudah ditutup tidak dapat diubah lagi.
	// Satu pintu untuk seluruh rute pengubah - lihat
	// services.PastikanKasusTerbuka, yang pemanggilannya ditagih penjaga
	// statik. Diperiksa SESUDAH wewenang: pemanggil yang tidak berhak tidak
	// berhak pula tahu keadaan kasusnya.
	if err := t.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return err
	}

	baca := repository.NewKlaimLife(t.svc.DB())
	if err := gerbangTahapDialogTanggal(ctx, baca, klaimID); err != nil {
		return err
	}
	tipe, err := baca.TypeKlaim(ctx, klaimID)
	if err != nil {
		return err
	}
	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	var target *models.Peserta
	nomor := 0
	for i := range peserta {
		if peserta[i].ID == pesertaID {
			target = &peserta[i]
			nomor = i + 1
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%w: peserta %q bukan milik klaim %q",
			galat.ErrPermintaanTidakSah, pesertaID, klaimID)
	}
	if err := ValidasiDOL(tipe, dol, *target); err != nil {
		return BeriNomorDOL(err, nomor)
	}
	return t.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		return baca.PerbaruiTanggalKejadian(ctx, tx, pesertaID, dol)
	})
}

// ErrTahapTidakBolehUbahTanggal - dialog Edit Date hanya terbuka di tahap
// yang dipegang Admin DAN membuka grid peserta (models.TahapBolehUbahTanggalKlaim).
var ErrTahapTidakBolehUbahTanggal = errors.New(
	"services: tanggal klaim hanya dapat diubah di tahap Outstanding Claim")

// ErrTanggalTerkunciSesudahSaveRNM - OQ-M1 (GILIRAN-17): tanggal klaim
// terkunci sesudah Save to RNM pertama berhasil (b1000 `CLAIM_NO!=”`).
var ErrTanggalTerkunciSesudahSaveRNM = errors.New(
	"services: tanggal klaim terkunci sesudah Save to RNM pertama berhasil")

// SetTanggalKlaim menyimpan tiga tanggal klaim lain seorang peserta -
// `UpdateDateClaimLife_Act`, tombol `Save` `EditDateClaimLife_Section` b1910.
//
// Urutan gerbangnya: identitas, pengenal, peran, kasus terbuka, tahap, lalu
// kepemilikan peserta. ⚠️ BEDA dengan `DiagnosaPeserta.pagari`: di sana
// perannya DITURUNKAN dari tahap dan diperiksa sesudahnya; di sini tahap yang
// sah hanya satu (Outstanding, dipegang Admin), jadi perannya tetap dan
// diperiksa lebih dulu - pemanggil tak berhak tidak berhak tahu keadaan kasus.
//
// ⛔ Peran DAN tahap, bukan salah satunya. `pyPosition=='ReasLifeAdmin'`
// adalah posisi KASUS (Register_Flow b582-b731 menyetelnya per tahap). Di
// Pega hanya pemegang penugasan yang dapat membuka kasus di posisi itu -
// penugasannya dirutekan ke Admin - jadi peran pelaku adalah padanan
// perutean itu, pola yang sama dengan `TahapLayanan.Pindah` ("orang
// mengubah pekerjaan yang sedang ia pegang"). Rute BARU ini lahir dengan
// gerbangnya; rute DOL lama tidak diubah (OQ-M8).
//
// ⚠️ Yang TIDAK dilakukan, dan sebabnya ada di kepala `models/tanggalklaim.go`:
// `ValidasiClaimReceived_Act` (penandanya tanpa kolom - OQ-M9) dan cermin
// warisan `UpdateDateClaimLife_SQL` (OQ-M2). Separuh "CLAIM_NO tidak kosong"
// ditiru maknanya sejak GILIRAN-17 (OQ-M1) di `gerbangTahapDialogTanggal`.
func (t *TanggalKejadian) SetTanggalKlaim(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string, tgl models.TanggalKlaim) error {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(pesertaID) == "" {
		return fmt.Errorf("%w: pengenal klaim dan peserta wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if err := inti.WajibPeran(pelaku, inti.PeranAdmin); err != nil {
		return err
	}
	if !t.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	// ⛔ BUTIR bb - kasus tertutup tidak dapat diubah lagi.
	if err := t.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return err
	}

	baca := repository.NewKlaimLife(t.svc.DB())
	if err := gerbangTahapDialogTanggal(ctx, baca, klaimID); err != nil {
		return err
	}

	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	milik := false
	for _, p := range peserta {
		if p.ID == pesertaID {
			milik = true
			break
		}
	}
	if !milik {
		return fmt.Errorf("%w: peserta %q bukan milik klaim %q",
			galat.ErrPermintaanTidakSah, pesertaID, klaimID)
	}
	return t.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		return baca.PerbaruiTanggalKlaim(ctx, tx, klaimID, pesertaID, tgl)
	})
}

// gerbangTahapDialogTanggal adalah gerbang SELURUH isian dialog Edit Date -
// DOL dan tiga tanggal lainnya (b1000, b1313, b1550, b1788), satu tempat:
// tahapnya, lalu penanda "sudah Save to RNM" (OQ-M1, GILIRAN-17).
func gerbangTahapDialogTanggal(ctx context.Context, baca *repository.KlaimLife, klaimID string) error {
	tahap, err := tahapKasus(ctx, baca, klaimID)
	if err != nil {
		return err
	}
	if !models.TahapBolehUbahTanggalKlaim(tahap) {
		return fmt.Errorf("%w: tahap %s", ErrTahapTidakBolehUbahTanggal, tahap)
	}
	sudah, err := baca.SudahSaveRNM(ctx, klaimID)
	if err != nil {
		return err
	}
	if !models.BolehUbahTanggalKlaim(tahap, sudah) {
		return ErrTanggalTerkunciSesudahSaveRNM
	}
	return nil
}
