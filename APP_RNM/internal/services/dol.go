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

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
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

// ⛔ DUA ASUMSI PUSTAKA PEGA, bukan fakta korpus. Keduanya diberi nama supaya
// terlihat, dapat dicari, dan dapat diubah dalam satu baris disertai testnya.
//
// ⛔ RALAT A0 — 27-09-2026. Ronde pertama menebak tanda tangannya dari CACAH
// ARGUMEN saja: `(tanggal, tahun, bulan, hari, jam, menit, detik, milidetik)`,
// yang menaruh JAM di posisi keempat - dan karena itu membaca
// `(.DATE_OF_LOSS,0,0,0,1,0,0,0)` sebagai satu JAM.
//
// Korpus sendiri menentukannya, dan itu terlewat: `[terverifikasi]`
// `Claim Life/Activity/LoadDataPeserta_Act.xml` pecahan baris 1097, 1143,
// 1163, 1183, 1203, 1223, 1243, dan 1263 memuat kedelapan tanggal peserta
// lewat `@addCalendar(<tanggal>,0,0,0,0,7,0,0)` - angka **7** di posisi
// KELIMA. Tujuh menit tidak menggeser apa pun yang berarti pada tanggal lahir,
// dan tujuh minggu memindahkannya ke bulan lain; tujuh JAM adalah WIB (UTC+7).
// Maka posisi kelima JAM, posisi keempat HARI - tanda tangan baku
// `(tanggal, tahun, bulan, minggu, hari, jam, menit, detik)`.
//
// `[terverifikasi - turunan]` `pergeseranDOLRetro`: `ValidasiDOL_Act` baris
// 698 menggeser DOL **SATU HARI** sebelum dibandingkan dengan jendela
// RETROCESSION (baris 719, 746). Cabang jendela GROSS di baris 458 memakai
// seluruh argumen nol - tidak bergeser sama sekali (baris 486, 526).
//
// `[dugaan - Product+UW]` `bandingKetat`: `@CompareDates(a,b)` bernilai benar
// bila `a` SESUDAH `b`, ketat. Korpus tidak memuat definisinya.
//
// Akibat kedua asumsi itu, jendelanya ASIMETRIS - dan asimetri itu memang
// berasal dari argumen yang BERBEDA di kedua cabang, bukan dari kami:
//
//	QP/QR : (BEGIN, EXPIRED]
//	TP/TR : [BEGIN, EXPIRED)
//
// Daftar LENGKAP kasus yang berbalik bila salah satu asumsi diputuskan lain
// ada di bab "Pembacaan ulang XML" tiket 06 - satu tempat, supaya tiga salinan
// tidak berjalan sendiri-sendiri. Ringkasnya: mengubah `bandingKetat` membalik
// batas di KEDUA ujung dan pada KEDUA cabang, bukan hanya satu kasus.
const (
	// Satu HARI, bukan satu jam - lihat ralat di atas.
	pergeseranDOLRetro = 24 * time.Hour
	bandingKetat       = true
)

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
}

// Error menulis kalimat XML, tanpa tambahan apa pun.
func (e GalatDOL) Error() string { return ErrDOLTidakSah.Error() }

// Unwrap membuat `errors.Is(err, ErrDOLTidakSah)` bernilai benar tanpa
// kalimatnya ikut berubah - Unwrap tidak menyentuh Error().
func (e GalatDOL) Unwrap() error { return ErrDOLTidakSah }

// ValidasiDOL menolak Date of Loss di luar jendela valuasi peserta.
//
// Meniru `Claim Life/Activity/ValidasiDOL_Act.xml` langkah 2, 3, dan 4:
//
//	langkah 2  Type QR atau QP -> jendela GROSS_VALUATION_*, tanpa pergeseran
//	langkah 3  Type TR atau TP -> jendela RETROCESSION_VALUATION_*, digeser
//	langkah 4  galat bila `local.Begin==false || local.Expired==true`
//
// ⛔ Type di luar keempat nilai itu adalah GALAT, bukan izin lewat. Di Pega
// kedua precondition tidak terpenuhi sehingga activity-nya diam saja dan
// tanggal apa pun lolos; diam seperti itu adalah lubang, bukan aturan.
func ValidasiDOL(tipe string, dol time.Time, p models.Peserta) error {
	var gross bool
	var geser time.Duration

	if !TypeDikenal(tipe) {
		return fmt.Errorf("%w: %q; yang dikenal hanya %s, %s, %s, dan %s",
			ErrTypeTidakDikenal, tipe, TypeQR, TypeQP, TypeTR, TypeTP)
	}
	switch tipe {
	case TypeQR, TypeQP:
		gross = true
	case TypeTR, TypeTP:
		geser = pergeseranDOLRetro
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

	// `local.DOL = @addCalendar(...)` - satu-satunya beda kedua cabang.
	dolGeser := dol.Add(geser)

	// `local.Begin` dan `local.Expired`, lalu gerbang langkah 4:
	// galat bila `Begin==false || Expired==true`.
	begin := sesudah(dolGeser, awal)
	expired := sesudah(dolGeser, akhir)
	if !begin || expired {
		return GalatDOL{PesertaID: p.ID}
	}
	return nil
}

// sesudah meniru `@CompareDates(a, b)`.
func sesudah(a, b time.Time) bool {
	if bandingKetat {
		return a.After(b)
	}
	return a.After(b) || a.Equal(b)
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
	for i := range peserta {
		if peserta[i].ID == pesertaID {
			target = &peserta[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%w: peserta %q bukan milik klaim %q",
			galat.ErrPermintaanTidakSah, pesertaID, klaimID)
	}
	if err := ValidasiDOL(tipe, dol, *target); err != nil {
		return err
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

// PergeseranDOLRetro membuka pergeseran itu untuk diuji.
//
// ⚠️ Fungsi, bukan konstanta yang diekspor: nilainya hasil pembacaan korpus
// dan boleh berubah bila bacaan itu diralat lagi - dan bila berubah, testnya
// yang menyebut jangkarnya ikut merah.
func PergeseranDOLRetro() time.Duration { return pergeseranDOLRetro }
