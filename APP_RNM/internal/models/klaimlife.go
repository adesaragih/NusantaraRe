// Package models memuat struct domain.
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/uang"
	"nusantarare/inti/utils"
)

// StatusBaris adalah keadaan satu baris AdjustmentList.
//
// Baris AdjustmentList - bukan klaim, bukan peserta - adalah unit keputusan
// mesin status (ADR-U-0011).
type StatusBaris int

const (
	// StatusTidakDiketahui dipakai untuk kolom kosong maupun kode di luar
	// ketiga nilai yang tertulis di spec. Artinya TIDAK ditebak.
	StatusTidakDiketahui StatusBaris = iota
	StatusOutstanding
	StatusAksep
	StatusDitolak
)

// String menulis status sebagai kata yang dibaca pengguna.
//
// Tiket 01 AC-3 / spec.md US-26: status ditampilkan sebagai kata, bukan sebagai
// nama field `STS_REJECT` dan bukan sebagai angka.
func (s StatusBaris) String() string {
	switch s {
	case StatusOutstanding:
		return "Outstanding"
	case StatusAksep:
		return "Aksep"
	case StatusDitolak:
		return "Ditolak"
	default:
		return "Tidak diketahui"
	}
}

// Diketahui membedakan status yang benar-benar tertulis di spec dari yang tidak.
func (s StatusBaris) Diketahui() bool { return s != StatusTidakDiketahui }

// Kode mentah kolom STS_REJECT, ditulis dan dibandingkan sebagai TEKS
// (ADR-U-0022). Menamainya di satu tempat membuat penulisnya dapat dicari:
// nol yang tersebar sebagai literal di banyak berkas tidak dapat ditelusuri.
//
// ⛔ Kode "4" SENGAJA tidak punya nama di sini. Ia ada di data warisan,
// artinya belum diputuskan work owner, dan sistem baru tidak pernah
// menulisnya. Nama akan membuatnya tampak seperti pilihan yang sah.
const (
	KodeOutstanding = "0"
	KodeAksep       = "1"
	KodeDitolak     = "2"
)

// StatusBarisDariKode menerjemahkan nilai kolom `STS_REJECT`.
//
// ⚠️ Namanya menyesatkan: nilai "1" berarti **DIAKSEP**, bukan ditolak.
// spec.md bab Problem butir 3, dan tabel "Penyimpangan sadar 4":
//
//	ReasLifeAdmin insert ke Outstanding -> "0"
//	ReasLifeAdmin reject langsung       -> "2"
//	ReasLifeSPV tambah baris Outstanding-> "0"
//
// dan `[terverifikasi]` nilai "1" = diaksep (spec.md US-26, CONTEXT.md).
//
// Perbandingan dilakukan atas TEKS, tidak pernah lewat bilangan (ADR-U-0022):
// "006" yang dibaca sebagai 6 lolos pulang-pergi dan memecahkan penggolong.
// Karena itu "00" dan "01" BUKAN "0" dan "1".
func StatusBarisDariKode(kode string) StatusBaris {
	switch kode {
	case KodeOutstanding:
		return StatusOutstanding
	case KodeAksep:
		return StatusAksep
	case KodeDitolak:
		return StatusDitolak
	default:
		return StatusTidakDiketahui
	}
}

// BarisAdjustment adalah satu baris `AdjustmentList` milik seorang peserta.
//
// Sumber kolom: `T_CLAIMLF_ADJUSTMENT` (STRUKTUR-TABEL-CLAIM-LIFE.md).
type BarisAdjustment struct {
	ID string
	// KodeStatus adalah nilai mentah `STS_REJECT`, disimpan sebagai teks dan
	// tidak pernah dikonversi ke bilangan (ADR-U-0022).
	KodeStatus string
	// Kedelapan kolom WARISAN: baris kedua dan seterusnya menyalinnya dari
	// baris PERTAMA peserta yang sama (`SetIndexAdjustmentList` langkah 3).
	// Keenamnya uang meski namanya berakhiran _SHARE - nama kolom di korpus
	// ini terbukti menipu; klasifikasinya ADR-U-0003, bukan namanya.
	//
	// ⚠️ Kedelapan nama itu hidup di DUA tempat sekaligus: di sini sebagai
	// medan, dan di `services.KolomDiwarisi()` sebagai daftar nama kolom yang
	// dikunci test. Keduanya sengaja - daftar nama itulah yang membuktikan
	// STS_REJECT tidak ikut.
	ShareNusantaraRe uang.Money
	CedingRetention  uang.Money
	SumReasured      uang.Money
	SumInsured       uang.Money
	ShareRetro       uang.Money
	RetrocededShare  uang.Money
	CurrencyID       string
	// JumlahKlaim adalah `CLAIM_AMOUNT` beserta `CURRENCY`-nya. Uang tidak
	// pernah float (ADR-U-0003, ADR-U-0016). Ia BUKAN kolom warisan: tiap
	// baris punya jumlah klaimnya sendiri.
	JumlahKlaim      uang.Money
	NomorAkseptasi   string
	TanggalAkseptasi time.Time
	// KomiteID kosong bila baris belum pernah dikirim ke Komite. Roster dan
	// keputusan komite bukan milik Claim Life - ini rujukan, bukan salinan.
	// Isinya T_WORK_CLAIM.ID baris komite, berformat KMTLF-xxxxxx.
	KomiteID string
	// Ketiga medan bank adalah tujuan pembayaran klaim baris ini. Nama
	// warisannya berbeda: NAME_OF_BANK tetap, IDBANK menjadi ID_BANK, dan
	// ACCOUNTNO menjadi ACCOUNT_NO.
	//
	// ⛔ NomorRekening tetap TEKS dan tidak pernah menjadi bilangan. Nomor
	// rekening berawalan nol adalah hal biasa, dan mengubahnya menjadi angka
	// menghilangkan nol itu (ADR-U-0022).
	NamaBank      string
	IDBank        string
	NomorRekening string
	// Spreading adalah pecahan baris ini per treaty-year (tiket 14).
	// Induknya baris adjustment, bukan peserta dan bukan header klaim.
	Spreading []Spreading
}

// Status menerjemahkan kode mentah baris ini.
func (b BarisAdjustment) Status() StatusBaris { return StatusBarisDariKode(b.KodeStatus) }

// MarshalJSON menulis baris untuk kontrak API.
//
// Nama field warisan `STS_REJECT` tidak pernah muncul di kontrak; yang keluar
// adalah kata statusnya, ditemani kode mentah supaya nilai yang tidak dikenal
// tetap dapat dilaporkan tanpa ditebak.
func (b BarisAdjustment) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID               string     `json:"id"`
		Status           string     `json:"status"`
		KodeStatus       string     `json:"kodeStatus"`
		StatusDiketahui  bool       `json:"statusDiketahui"`
		JumlahKlaim      uang.Money `json:"jumlahKlaim"`
		NomorAkseptasi   string     `json:"nomorAkseptasi"`
		TanggalAkseptasi string     `json:"tanggalAkseptasi"`
		KomiteID         string     `json:"komiteId"`
	}{
		ID:               b.ID,
		Status:           b.Status().String(),
		KodeStatus:       b.KodeStatus,
		StatusDiketahui:  b.Status().Diketahui(),
		JumlahKlaim:      b.JumlahKlaim,
		NomorAkseptasi:   b.NomorAkseptasi,
		TanggalAkseptasi: utils.FormatTanggal(b.TanggalAkseptasi),
		KomiteID:         b.KomiteID,
	})
}

// Peserta adalah satu peserta yang diklaim, beserta barisnya sendiri.
//
// `[terverifikasi]` `Claim Life/Activity/SavePesertaClaim.xml`
// (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEPESERTACLAIM` / `RULE-OBJ-ACTIVITY`)
// menulis ke `…PremiumListDetail(<LAST>).AdjustmentList(<LAST>)` - **setiap
// peserta punya daftar adjustment sendiri**. Menggantungkan baris ke header
// menghapus informasi peserta pemilik dan mematahkan mesin status (ADR-U-0011).
//
// ⚠️ Tracer ini sengaja TIDAK membawa medan bernama orang (`NAME_OF_INSURED`,
// `POLICY_HOLDER`). Tiket 01 tidak memerlukannya, dan menahannya menjauhkan
// fixture dari data pribadi.
type Peserta struct {
	ID              string
	NomorPremiList  string
	NomorPolis      string
	NomorSertifikat string
	MataUang        string

	// SumberID adalah kolom ID baris asalnya di M_LIFE_PREMIUM_DETAIL,
	// disimpan di SOURCE_ID. Ia jejak: dari baris polis mana peserta klaim ini
	// disalin.
	SumberID string
	// IsCheck menandai peserta yang DIPILIH untuk diklaim (AC 8 tiket 03).
	//
	// ⛔ Nilainya SELALU PenandaDipilih, tidak pernah "1" atau "Y".
	IsCheck string
	// KodeStatus adalah CERMIN `STS_REJECT` baris adjustment yang terakhir
	// diputus, bukan keputusan tersendiri (tiket 04).
	//
	// ⛔ Ditambahkan 26-09-2026 karena kolomnya DITULIS sejak tiket 04 tetapi
	// tidak pernah DIBACA - sehingga AC pencerminan ke peserta tidak dapat
	// dibuktikan, dan separuh pencerminan dapat mati tanpa satu pun test
	// gagal. Cacat yang persis sama pernah terjadi pada CLAIM_RETRO.
	KodeStatus string

	// Empat tanggal valuasi dan WPC. ⭐ Disalin saat pendaftaran justru supaya
	// tiket 06 dapat memvalidasi DOL TANPA query ulang ke tabel 66,8 juta
	// baris - satu-satunya alasan kolom ini ada di sini.
	ValuasiGrossMulai   string
	ValuasiGrossSelesai string
	ValuasiRetroMulai   string
	ValuasiRetroSelesai string
	WPC                 string

	// Tanggal polis.
	TanggalMulai   string
	TanggalEfektif string
	TanggalLapse   string
	TanggalExpired string

	// STNC treaty, dibawa apa adanya sebagai teks.
	STNC string

	// Umur adalah umur tertanggung saat baris polisnya disalin, dipilih
	// dari tiga kolom sumber oleh repository.UmurPeserta
	// (`SaveInsuredClaim_Act.xml:661`).
	//
	// ⛔ TEKS, bukan int. Kolomnya boleh kosong, dan kosong BUKAN nol
	// (ADR-U-0027): bayi berumur nol tahun dan peserta yang umurnya tidak
	// tercatat adalah dua keadaan berbeda, dan int tidak dapat
	// membedakannya tanpa penunjuk.
	Umur string

	// TanggalKejadian adalah DATE_OF_LOSS - tanggal kejadian yang diklaim,
	// dan ia milik PESERTA, bukan klaim. `[terverifikasi]` ValidasiDOL_Act
	// berkelas Int-LIFE_PREMIUM_DETAIL dan menempelkan pesan galatnya pada
	// `.DATE_OF_LOSS` peserta (parameter Field langkah 4). Kolomnya sudah
	// disediakan migrasi 003; tiket 06 yang memakainya.
	TanggalKejadian string

	// Tiga tanggal klaim lain dialog Edit Date, milik PESERTA seperti DOL
	// (`EditDateClaimLife_Section` berkelas Int-LIFE_PREMIUM_DETAIL). Kolomnya
	// sudah disediakan migrasi 003. Kosong berarti belum diisi (ADR-U-0027).
	//
	//	TanggalTerimaKlaim    CLAIM_RECEIVED_DATE  b1076
	//	TanggalDokumenLengkap COMPLETE_DATE        b1387 "DOCUMENT COMPLETE DATE"
	//	TanggalKonfirmasi     CONFIRMATION_DATE    b1626
	//
	// ⛔ Satu tipe, bukan tiga medan lepas (GILIRAN-11 paket 4, data clump):
	// ketiganya selalu berjalan bersama - dibaca, ditampilkan, dan disimpan
	// satu tombol. Medannya tetap diakses langsung (`p.TanggalTerimaKlaim`).
	TanggalKlaimTeks

	// PenandaTerimaKlaim adalah `.MAXCLAIM_RECEIVED` (butir bk), DIHITUNG saat
	// baca - kosong berarti sah. PenandaTerimaKlaimAlasan terisi bila ia TIDAK
	// dapat dihitung (polis, produk, atau ambangnya belum terbaca): kosong di
	// sana bukan "sah".
	PenandaTerimaKlaim       string
	PenandaTerimaKlaimAlasan string

	// Diagnosa adalah daftar `.DiagnoseList` peserta ini - butir bd.
	//
	// ⛔ Daftar, bukan sepasang kolom. `ClaimLifeDetailGCNM.xml` b3923
	// menyajikannya `RepeatGrid` dengan tombol `Add` b4690 dan `Delete`
	// b6160; kolom tunggal `DISEASE`/`ICD_CODE` pada peserta (migrasi 003)
	// adalah warisan yang hanya sanggup menyimpan satu.
	Diagnosa []Diagnosa

	// Uang polis. Seluruhnya Money kecuali EMPercent, yang perbandingan -
	// keduanya sengaja bertipe berbeda supaya tidak pernah terjumlahkan
	// (ADR-F-0004).
	SumInsured       uang.Money
	SumReasured      uang.Money
	GrossPremium     uang.Money
	NetPremium       uang.Money
	CedingRetention  uang.Money
	ShareNusantaraRe uang.Money
	ShareRetro       uang.Money
	RetrocededShare  uang.Money
	EMPercent        uang.Ratio
	// JumlahKlaim adalah `CLAIM_AMOUNT` peserta - GILIRAN-14 butir bp.
	//
	// `[terverifikasi]` `SavePesertaClaim` 7.7 b3280 menyalinnya ke peserta,
	// dan 7.8 b3828 ke baris adjustment pertama - keduanya dari kolom yang
	// sama di sumber (`M_LIFE_PREMIUM_DETAIL.CLAIM_AMOUNT`). Sebelum bp ia
	// tidak dibaca sama sekali, sehingga baris pertama tidak dapat lahir
	// bernilai.
	JumlahKlaim uang.Money

	// ⛔ NAMA ORANG SENGAJA TIDAK ADA DI SINI. NAME_OF_INSURED dan
	// POLICY_HOLDER punya kolomnya di DDL, tetapi tidak disalin: nama
	// tertanggung hanya diperlukan LAYAR saat memilih, dan itu dilayani
	// CalonPeserta. Menyalinnya ke tabel klaim berarti menduplikasi data
	// pribadi tanpa satu pun AC yang memintanya. Kolomnya tetap NULL sampai
	// ada keputusan work owner. [terbuka]

	Baris []BarisAdjustment

	// Total adalah keenam total uang peserta ini, DIHITUNG saat dibaca dan
	// tidak pernah disimpan - lihat totalpeserta.go. Ia diisi oleh
	// services.KlaimLife.Ambil; repository tidak pernah menyentuhnya, dan
	// ada penjaga statik yang menagihnya.
	Total TotalPeserta

	// Dokumen pendukung peserta ini. Ia GERBANG simpan ke Outstanding, bukan
	// pelengkap: SaveOutStandingLife_Act menolak menyimpan bila peserta yang
	// dipilih belum mengunggah dokumen.
	Dokumen []Dokumen
}

// ErrValuasiKosong menandai kolom tanggal valuasi yang kosong.
//
// ⛔ Kosong DITOLAK, tidak diperlakukan sebagai waktu nol (ADR-U-0022 Akibat 2,
// ADR-U-0027). Waktu nol adalah tahun 1 Masehi, sehingga setiap tanggal akan
// tampak sesudah tanggal mulai dan setiap klaim lolos - kegagalan yang arahnya
// paling berbahaya: diam dan meloloskan.
var ErrValuasiKosong = errors.New("models: tanggal valuasi peserta kosong")

// JendelaValuasi mengurai sepasang tanggal valuasi peserta menjadi waktu.
//
// ⚠️ [terbuka - tiket 14] Tempat SEBENARNYA konversi ini adalah pemuat di
// `repository` (ADR-U-0022 Akibat 1: "konversi adalah tanggung jawab pemuat di
// repository, bukan tersebar di lapisan layanan"). Ia ada di sini - pada tipe
// yang memiliki teksnya, satu tempat dan tidak tersebar - karena kesembilan
// medan tanggal peserta bertipe TEKS, dan mengubahnya menjadi waktu adalah
// keputusan sekali untuk seluruh model beserta daftar kolomnya. Itu milik
// tiket 14, bukan tiket 06.
//
// Nama kolom yang dilaporkan adalah nama kolom SEBENARNYA di migrasi 003 -
// bukan nama medan Go, dan bukan nama ketiga yang dikarang.
func (p Peserta) JendelaValuasi(gross bool) (mulai, selesai time.Time, err error) {
	namaMulai, namaSelesai := "RETRO_VALUATION_BEGIN_DATE", "RETRO_VALUATION_EXPIRED_DATE"
	teksMulai, teksSelesai := p.ValuasiRetroMulai, p.ValuasiRetroSelesai
	if gross {
		namaMulai, namaSelesai = "GROSS_VALUATION_BEGIN_DATE", "GROSS_VALUATION_EXPIRED_DATE"
		teksMulai, teksSelesai = p.ValuasiGrossMulai, p.ValuasiGrossSelesai
	}
	if mulai, err = p.uraiTanggal(namaMulai, teksMulai); err != nil {
		return time.Time{}, time.Time{}, err
	}
	if selesai, err = p.uraiTanggal(namaSelesai, teksSelesai); err != nil {
		return time.Time{}, time.Time{}, err
	}
	return mulai, selesai, nil
}

// uraiTanggal mengurai satu kolom tanggal peserta, menyebut kolom dan peserta.
func (p Peserta) uraiTanggal(kolom, nilai string) (time.Time, error) {
	if strings.TrimSpace(nilai) == "" {
		return time.Time{}, fmt.Errorf("%w: peserta %s kolom %s",
			ErrValuasiKosong, p.ID, kolom)
	}
	w, err := utils.ParseTanggal(nilai)
	if err != nil {
		return time.Time{}, fmt.Errorf("models: peserta %s kolom %s: %w", p.ID, kolom, err)
	}
	return w, nil
}

// MarshalJSON menulis peserta untuk kontrak API.
func (p Peserta) MarshalJSON() ([]byte, error) {
	baris := p.Baris
	if baris == nil {
		baris = []BarisAdjustment{}
	}
	// Daftar KOSONG, bukan nil: `encoding/json` menulis nil sebagai `null`,
	// dan layar yang menerima `null` harus menjaganya sendiri.
	dokumen := p.Dokumen
	if dokumen == nil {
		dokumen = []Dokumen{}
	}
	diagnosa := p.Diagnosa
	if diagnosa == nil {
		diagnosa = []Diagnosa{}
	}
	return json.Marshal(struct {
		ID              string `json:"id"`
		NomorPremiList  string `json:"nomorPremiList"`
		NomorPolis      string `json:"nomorPolis"`
		NomorSertifikat string `json:"nomorSertifikat"`
		MataUang        string `json:"mataUang"`
		// ⭐ IsCheck menyeberang sejak audit A0: layar memerlukannya untuk
		// memutuskan apakah kontrol "Save Adjustment" tampil - prasyarat XML
		// `SaveAdjustment_Act` menuntut peserta DIPILIH.
		IsCheck string `json:"isCheck"`
		// ⭐ TanggalKejadian menyeberang sejak kelompok Detail & Tutup:
		// tombol `Edit Date` (`ClaimLifeDetailGCNM.xml` b14115) harus
		// menampilkan tanggal yang SEDANG berlaku. Tanpa itu kotaknya selalu
		// terbuka kosong, dan pemakai tidak dapat membedakan "belum diisi"
		// dari "sudah diisi tetapi tidak terlihat" - dan yang kedua membuat
		// orang mengetik ulang tanggal yang sudah benar.
		//
		// ⚠️ TEKS apa adanya, bukan tanggal yang diurai ulang di layar.
		TanggalKejadian string `json:"tanggalKejadian"`
		// ⭐ Tiga tanggal klaim menyeberang sejak 28-09-2026 dengan sebab
		// yang sama: dialog Edit Date membuka nilai yang SEDANG berlaku.
		TanggalTerimaKlaim    string `json:"tanggalTerimaKlaim"`
		TanggalDokumenLengkap string `json:"tanggalDokumenLengkap"`
		TanggalKonfirmasi     string `json:"tanggalKonfirmasi"`
		// ⭐ Butir bk - sel `MAX CLAIM RECEIVED` b12131, dihitung saat baca.
		PenandaTerimaKlaim       string            `json:"penandaTerimaKlaim"`
		PenandaTerimaKlaimAlasan string            `json:"penandaTerimaKlaimAlasan"`
		Baris                    []BarisAdjustment `json:"baris"`
		// ⭐ Total menyeberang sejak ralat 27-09-2026: enam total, bukan
		// lima, dan bukan penanda "belum tersedia". Ia turunan, bukan kolom.
		Total TotalPeserta `json:"total"`
		// ⭐ Dokumen menyeberang sejak kelompok 1 giliran 12. Tanpa ini
		// layar Detail memanggil rute kedua hanya untuk menampilkan daftar
		// yang sudah ada di tangan - dan rute kedua itu tidak pernah dibuat,
		// sehingga daftarnya tidak pernah tampil.
		//
		// ⛔ NAMA BERKAS ikut; nama ORANG tidak. `NAMA_FILE` adalah nama
		// berkas unggahan, bukan nama tertanggung.
		Dokumen []Dokumen `json:"dokumen"`
		// ⭐ Diagnosa menyeberang sejak butir bd. Tanpa ini grid diagnosa
		// harus memanggil rute kedua untuk setiap peserta yang dibuka -
		// dan klaim grup berpeserta ratusan membuat itu terasa rusak.
		Diagnosa []Diagnosa `json:"diagnosa"`
		// ⭐ KodeStatus menyeberang sejak butir bd, dan ia MEDAN KEENAM
		// yang hampir menjadi cacat lintas-lapis keenam. Layar memerlukannya
		// untuk meniru gerbang `pyDisabledWhen` b4682 - yang diuji di sana
		// `.STS_REJECT` PESERTA. Tanpa medan ini `peserta.kodeStatus` di
		// TypeScript adalah `undefined`, pembandingnya selalu tidak sama,
		// dan gerbangnya tidak pernah menutup: tombol tetap hidup, backend
		// menjawab 409, dan pemakai belajar mengabaikan galat.
		//
		// ⚠️ Kolomnya DITULIS sejak tiket 04; yang kurang hanya
		// penyeberangannya. Cacat serupa pernah nyata pada CLAIM_RETRO.
		KodeStatus string `json:"kodeStatus"`
	}{
		// ⛔ Diisi BERNAMA, bukan berposisi. Tujuh medan berurutan yang
		// enam di antaranya bertipe string: dua yang tertukar tetap
		// dikompilasi, tetap lolos uji bentuk, dan baru terlihat ketika
		// nomor sertifikat muncul di kolom tanggal.
		ID:                       p.ID,
		NomorPremiList:           p.NomorPremiList,
		NomorPolis:               p.NomorPolis,
		NomorSertifikat:          p.NomorSertifikat,
		MataUang:                 p.MataUang,
		IsCheck:                  p.IsCheck,
		TanggalKejadian:          p.TanggalKejadian,
		TanggalTerimaKlaim:       p.TanggalTerimaKlaim,
		TanggalDokumenLengkap:    p.TanggalDokumenLengkap,
		TanggalKonfirmasi:        p.TanggalKonfirmasi,
		PenandaTerimaKlaim:       p.PenandaTerimaKlaim,
		PenandaTerimaKlaimAlasan: p.PenandaTerimaKlaimAlasan,
		Baris:                    baris,
		Total:                    p.Total,
		Dokumen:                  dokumen,
		Diagnosa:                 diagnosa,
		KodeStatus:               p.KodeStatus,
	})
}

// PenandaDipilih adalah SATU-SATUNYA nilai yang berarti "peserta ini diklaim".
//
// ⛔ VERBATIM dari rule, bukan pilihan gaya:
//
//	`Activity/SetIndexAdjustmentList.xml` b328  `.IsCheck` = `true`
//	`Activity/SavePesertaClaim.xml` b3631, b3919  `.IsCheck=="true"`
//	                                (langkah 7.7 dan 7.8 - keduanya HIDUP)
//
// RALAT GILIRAN-11: dulu dikutip b1812, b1997, b2413 - WHEN langkah 7.1
// (precondition false, mati) serta 7.2 dan 7.4 (`//`, ter-remark). Nilainya
// sama; hanya buktinya yang dipindah ke baris yang benar-benar berjalan.
//
// ⛔ CACAT YANG PERNAH ADA, dan sebab konstanta ini lahir: jalur
// pendaftaran menuliskan "1" sedangkan gerbang akseptasi menuntut "true",
// sehingga setiap peserta yang baru didaftarkan DITOLAK saat hendak diaksep -
// dan seluruh jalur akseptasi tidak dapat dijalankan untuk klaim baru.
// Kedua sisi benar menurut dirinya sendiri; hanya pertemuannya yang salah,
// dan tidak ada satu pun uji yang memaksa keduanya bertemu.
const PenandaDipilih = "true"

// Klaim adalah satu klaim Life beserta seluruh pesertanya.
//
// Sumber kolom: `T_GENERAL_CLAIM` (STRUKTUR-TABEL-CLAIM-LIFE.md).
type Klaim struct {
	ID         string
	NomorKlaim string
	NomorPolis string
	NamaBisnis string
	// KodeBisnis adalah `BUSINESSID` - kode produk, bukan namanya.
	//
	// ⛔ TEMUAN AUDIT A0, 27-09-2026: nomor akseptasi memuatnya
	// (`'RNML-A'||{pyWorkPage.BusinessCode}||…`, `Generate_NoAccept_Life.xml`
	// baris 85), tetapi model relasional TIDAK menyimpannya - `T_GENERAL_CLAIM`
	// hanya punya `BUSINESS_NAME`, dan `BUSINESSID` bukan salah satu dari 18
	// kolom datar warisan yang `Simpan` tulis. Ia dipakai sekali saat
	// pendaftaran lalu hilang.
	//
	// ⚠️ Medannya ada di sini SEJAK SEKARANG supaya jalur akseptasi dapat
	// menyebutnya dan gagal terang; kolomnya lahir di A1, dan sampai itu
	// nilainya selalu kosong.
	KodeBisnis string
	// KodeStatus adalah `STS_REJECT` di tingkat header. Ia dibawa apa adanya
	// dan TIDAK pernah diterjemahkan menjadi kata: spec.md menyatakan ia
	// cerminan baris adjustment terakhir - turunan, bukan unit keputusan.
	// Unit keputusan tetap baris (ADR-U-0011).
	KodeStatus string

	// ClaimRetro adalah bagian retro dari jumlah klaim, di tingkat header.
	//
	// `[keputusan work owner 26-09-2026, butir w2]` ia UANG, bukan
	// perbandingan. Buktinya agregat instance pengembangan: nol nilai di
	// rentang 0-1 maupun 1-100, 4.778 nilai di atas 100, 475 berdesimal - dan
	// label Pega berbunyi "Claim Retro" di samping label uang lain. Karena itu
	// Money, bukan Ratio: keduanya sengaja bertipe berbeda supaya tidak pernah
	// terjumlahkan (ADR-F-0004).
	//
	// ⚠️ [terbuka] T_GENERAL_CLAIM tidak punya kolom mata uang, sehingga
	// Currency-nya datang dari baris adjustment saat klaim dibongkar dari tabel
	// warisan, dan KOSONG saat header dibaca sendirian. Menambah kolom mata
	// uang ke header adalah keputusan tersendiri, dan executor tidak
	// mengarangnya.
	ClaimRetro uang.Money

	Peserta []Peserta

	// Tahap adalah `T_WORK_CLAIM.TAHAP` - nama assignment VERBATIM. Ia tidak
	// datang dari `AmbilHeader` (yang membaca T_GENERAL_CLAIM) melainkan
	// dirakit `services.KlaimLife.Ambil` dari baris work.
	//
	// ⚠️ Boleh KOSONG. Aplikasi ini tidak pernah menyisipkan ke
	// T_WORK_CLAIM - baris itu lahir di sistem lama - jadi klaim tanpa baris
	// work adalah keadaan nyata, bukan kerusakan. Kosong berarti tahapnya
	// tidak diketahui, dan layar karena itu tidak menawarkan Close Claim.
	Tahap string

	// StatusWork adalah `T_WORK_CLAIM.STATUS_WORK` - butir bb. Kosong berarti
	// kasusnya BELUM ditutup (ADR-U-0027); satu-satunya nilai lain yang
	// pernah ditulis adalah models.StatusWorkSelesai.
	StatusWork string
}

// CacahBaris menghitung seluruh baris adjustment di seluruh peserta.
func (k Klaim) CacahBaris() int {
	n := 0
	for _, p := range k.Peserta {
		n += len(p.Baris)
	}
	return n
}

// MarshalJSON menulis klaim untuk kontrak API.
func (k Klaim) MarshalJSON() ([]byte, error) {
	peserta := k.Peserta
	if peserta == nil {
		peserta = []Peserta{}
	}
	// ⛔ KodeStatus header dibawa MENTAH, tanpa tafsir. spec.md, "Aturan yang
	// mengikat": "**`STS_REJECT` adalah status baris**, bukan status klaim" dan
	// "Nilai `2` selalu berarti 'baris ini ditolak', tidak pernah 'klaim
	// selesai'". Kolom itu cerminan baris terakhir, bukan unit keputusan, dan
	// menerjemahkannya menjadi kata akan mengarang status yang tidak pernah
	// ditetapkan siapa pun.
	//
	// ⭐ Diperluas tiket 04: yang berhak menjadi KATA adalah StatusTurunan,
	// yang DIHITUNG dari seluruh baris - bukan dibaca dari kolom cermin itu.
	// Keduanya sengaja berdampingan di kontrak: yang satu fakta mentah dari
	// basis data, yang lain kesimpulan yang dapat dipertanggungjawabkan.
	return json.Marshal(struct {
		ID            string     `json:"id"`
		NomorKlaim    string     `json:"nomorKlaim"`
		NomorPolis    string     `json:"nomorPolis"`
		NamaBisnis    string     `json:"namaBisnis"`
		KodeStatus    string     `json:"kodeStatus"`
		StatusTurunan string     `json:"statusTurunan"`
		ClaimRetro    uang.Money `json:"claimRetro"`
		Peserta       []Peserta  `json:"peserta"`
		CacahBaris    int        `json:"cacahBaris"`
		// ⭐ Tahap dan StatusWork menyeberang sejak butir bb. Layar Detail
		// memerlukan KEDUANYA untuk memutuskan apakah tombol `Close Claim`
		// pantas ditawarkan: `pyLocalAction>CloseClaim` hanya ada di dua
		// section (Outstanding, Akseptasi), dan kasus yang sudah tertutup
		// tidak menawarkannya lagi.
		//
		// ⚠️ Keduanya boleh KOSONG: aplikasi ini tidak pernah menyisipkan
		// ke T_WORK_CLAIM, sehingga klaim tanpa baris work adalah keadaan
		// nyata. Kosong = tahap tidak diketahui = tombolnya tidak ditawarkan.
		Tahap      string `json:"tahap"`
		StatusWork string `json:"statusWork"`
	}{
		StatusTurunan: k.StatusTurunan().String(),
		ID:            k.ID,
		NomorKlaim:    k.NomorKlaim,
		NomorPolis:    k.NomorPolis,
		NamaBisnis:    k.NamaBisnis,
		KodeStatus:    k.KodeStatus,
		// Money punya MarshalJSON sendiri: jumlahnya TEKS, tidak pernah angka
		// JSON (ADR-U-0003). Tanpa baris ini nilai yang sudah dibaca dari
		// Oracle dibuang diam-diam saat serialisasi - separuh butir w2 mati.
		ClaimRetro: k.ClaimRetro,
		Peserta:    peserta,
		CacahBaris: k.CacahBaris(),
		Tahap:      k.Tahap,
		StatusWork: k.StatusWork,
	})
}
