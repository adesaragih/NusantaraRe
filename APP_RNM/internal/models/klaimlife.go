package models

import (
	"encoding/json"
	"time"

	"nusantarare/pkg/utils"
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
	case "0":
		return StatusOutstanding
	case "1":
		return StatusAksep
	case "2":
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
	// JumlahKlaim adalah `CLAIM_AMOUNT` beserta `CURRENCY`-nya. Uang tidak
	// pernah float (ADR-U-0003, ADR-U-0016).
	JumlahKlaim      Money
	NomorAkseptasi   string
	TanggalAkseptasi time.Time
	// KomiteID kosong bila baris belum pernah dikirim ke Komite. Roster dan
	// keputusan komite bukan milik Claim Life - ini rujukan, bukan salinan.
	// Isinya T_WORK_CLAIM.ID baris komite, berformat KMT-xxxxxx.
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
		ID               string `json:"id"`
		Status           string `json:"status"`
		KodeStatus       string `json:"kodeStatus"`
		StatusDiketahui  bool   `json:"statusDiketahui"`
		JumlahKlaim      Money  `json:"jumlahKlaim"`
		NomorAkseptasi   string `json:"nomorAkseptasi"`
		TanggalAkseptasi string `json:"tanggalAkseptasi"`
		KomiteID         string `json:"komiteId"`
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
	IsCheck string

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

	// Uang polis. Seluruhnya Money kecuali EMPercent, yang perbandingan -
	// keduanya sengaja bertipe berbeda supaya tidak pernah terjumlahkan
	// (ADR-F-0004).
	SumInsured       Money
	SumReasured      Money
	GrossPremium     Money
	NetPremium       Money
	CedingRetention  Money
	ShareNusantaraRe Money
	ShareRetro       Money
	RetrocededShare  Money
	EMPercent        Ratio

	// ⛔ NAMA ORANG SENGAJA TIDAK ADA DI SINI. NAME_OF_INSURED dan
	// POLICY_HOLDER punya kolomnya di DDL, tetapi tidak disalin: nama
	// tertanggung hanya diperlukan LAYAR saat memilih, dan itu dilayani
	// CalonPeserta. Menyalinnya ke tabel klaim berarti menduplikasi data
	// pribadi tanpa satu pun AC yang memintanya. Kolomnya tetap NULL sampai
	// ada keputusan work owner. [terbuka]

	Baris []BarisAdjustment
}

// MarshalJSON menulis peserta untuk kontrak API.
func (p Peserta) MarshalJSON() ([]byte, error) {
	baris := p.Baris
	if baris == nil {
		baris = []BarisAdjustment{}
	}
	return json.Marshal(struct {
		ID              string            `json:"id"`
		NomorPremiList  string            `json:"nomorPremiList"`
		NomorPolis      string            `json:"nomorPolis"`
		NomorSertifikat string            `json:"nomorSertifikat"`
		MataUang        string            `json:"mataUang"`
		Baris           []BarisAdjustment `json:"baris"`
	}{p.ID, p.NomorPremiList, p.NomorPolis, p.NomorSertifikat, p.MataUang, baris})
}

// Klaim adalah satu klaim Life beserta seluruh pesertanya.
//
// Sumber kolom: `T_GENERAL_CLAIM` (STRUKTUR-TABEL-CLAIM-LIFE.md).
type Klaim struct {
	ID         string
	NomorKlaim string
	NomorPolis string
	NamaBisnis string
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
	ClaimRetro Money

	Peserta []Peserta
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
	// ⛔ TIDAK ada kata status di tingkat klaim. spec.md, "Aturan yang
	// mengikat": "**`STS_REJECT` adalah status baris**, bukan status klaim" dan
	// "Nilai `2` selalu berarti 'baris ini ditolak', tidak pernah 'klaim
	// selesai'". Header hanyalah cerminan baris terakhir - turunan, bukan unit
	// keputusan - sehingga menerjemahkannya menjadi kata akan mengarang status
	// klaim yang tidak pernah ditetapkan siapa pun. Kode mentahnya dibawa apa
	// adanya, tanpa tafsir.
	return json.Marshal(struct {
		ID         string    `json:"id"`
		NomorKlaim string    `json:"nomorKlaim"`
		NomorPolis string    `json:"nomorPolis"`
		NamaBisnis string    `json:"namaBisnis"`
		KodeStatus string    `json:"kodeStatus"`
		ClaimRetro Money     `json:"claimRetro"`
		Peserta    []Peserta `json:"peserta"`
		CacahBaris int       `json:"cacahBaris"`
	}{
		ID:         k.ID,
		NomorKlaim: k.NomorKlaim,
		NomorPolis: k.NomorPolis,
		NamaBisnis: k.NamaBisnis,
		KodeStatus: k.KodeStatus,
		// Money punya MarshalJSON sendiri: jumlahnya TEKS, tidak pernah angka
		// JSON (ADR-U-0003). Tanpa baris ini nilai yang sudah dibaca dari
		// Oracle dibuang diam-diam saat serialisasi - separuh butir w2 mati.
		ClaimRetro: k.ClaimRetro,
		Peserta:    peserta,
		CacahBaris: k.CacahBaris(),
	})
}
