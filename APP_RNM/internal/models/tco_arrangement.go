package models

// Entitas master arrangement kontrak treaty non-life - tiket 01 Treaty
// Contract Out.
//
// Untuk apa berkas ini: bentuk enam tabel `TREATY*` warisan (STRUKTUR-TABEL-
// TREATY-CONTRACT-OUT.md) di sisi Go, dengan tipe yang SUDAH dirapikan
// (penyimpangan sadar 6): uang dan persen `*apd.Decimal`, tanggal `time.Time`.
//
// ⛔ Nama medan mengikuti nama kolom warisan (VERBATIM, spec §4 - istilah
// Pega dipakai apa adanya, tidak diterjemahkan): `Ydcf`, `PctMe`,
// `LayerPartType`, `SpreadingOrder`, `Method`, `Line`. Artinya OQ Product + UW.
//
// ⛔ Nol float64 di jalur uang mana pun (ADR-U-0003). `*apd.Decimal` ber-
// MarshalText, sehingga di JSON ia TEKS - bukan angka JSON yang dibaca
// float64 di React.
//
// ⚠️ Waktu NOL (`time.Time{}`) berarti kolom KOSONG (ADR-U-0027), bukan
// tanggal 1 Januari tahun 1.
//
// Hierarki (spec §2): TahunTreaty -> KontrakTreaty; ReinsurerTreaty
// (-> SecurityReinsurer), BusinessTreaty, dan KlausulTreaty menggantung pada
// KUNCI GABUNGAN (TreatyYear, TreatyGroupID, ReinsTypeID) - BUKAN pada ID
// kontrak.

import (
	"time"

	"github.com/cockroachdb/apd/v3"
)

// TahunTreaty adalah satu baris `TREATYYEAR`.
type TahunTreaty struct {
	ID               string
	TreatyYear       string
	UnderwritingYear string
	TreatyGroupID    string
	TreatyGroupName  string
	// Proportion tetap TEKS: artinya `[terbuka]`; layar mengisinya dari
	// pilihan "Reinsurance Type" (`InputDtlTreatyContact.xml` b6829).
	Proportion string
	StartDate  time.Time
	EndDate    time.Time
	UserID     string
	TglUpdate  time.Time
}

// KontrakTreaty adalah satu baris `TREATYCONTRACT` - satu jenis reasuransi
// yang dibuka pada sebuah tahun treaty.
type KontrakTreaty struct {
	ID              string
	IDTreatyYear    string
	ReinsTypeID     string
	ReinsTypeName   string
	TreatyStartDate time.Time
	TreatyEndDate   time.Time
	UserID          string
	TglUpdate       time.Time
}

// ReinsurerTreaty adalah satu baris `TREATYREINSURER`.
//
// `Name` adalah nama PERUSAHAAN reinsurer (kolom warisan `NAME`), bukan nama
// orang. `OperatorName` diisi `OperatorID.pyUserName` - pengenal akun.
type ReinsurerTreaty struct {
	ID              string
	TreatyYear      string
	TreatyGroupID   string
	TreatyGroupName string
	ReinsTypeID     string
	ReinsTypeName   string
	ReinsurerID     string
	ClientID        string
	Name            string
	Ricomm          *apd.Decimal
	PctShare        *apd.Decimal
	// IUDate tetap TEKS: bentuk dan artinya `[terbuka]`, tidak ada di daftar
	// tanggal AC 53 - dibawa apa adanya.
	IUDate       string
	UserID       string
	StartDate    time.Time
	EndDate      time.Time
	StatusOn     string
	StdRating    string
	OperatorName string
	TglUpdate    time.Time
}

// SecurityReinsurer adalah satu baris `MTREATYSECURITY` warisan (tco4): tabel
// TANPA identitas - `ID` membawa `REAS_SECURITY` terpangkas, kunci baris
// (`REAS_ID`, `TRIM(REAS_SECURITY)`) seperti `UpdateMTreatySecurity`. RALAT
// penyimpangan sadar 5 (PK surrogate) - dibatalkan tco4.
type SecurityReinsurer struct {
	ID        string
	ThnTreaty string
	TopID     string
	TpTreaty  string
	// ReasID menunjuk ReinsurerTreaty.ID (tanpa FK - kaskade oleh layanan).
	ReasID       string
	PctShare     *apd.Decimal
	UserID       string
	ReasSecurity string
}

// BusinessTreaty adalah satu baris `TREATYBUSINESS`.
type BusinessTreaty struct {
	ID       string
	IsActive string
	// TreatyYearID boleh KOSONG pada data lama (`DeleteFromTREATYCONTRACT_SQL`:
	// `OR TREATYYEARID IS NULL`).
	TreatyYear      string
	TreatyYearID    string
	TreatyGroupID   string
	TreatyGroupName string
	ReinsTypeID     string
	ReinsTypeName   string
	BizCode         string
	BizName         string
	UserID          string
	TglUpdate       time.Time
}

// KlausulTreaty adalah satu baris `PROPORTIONALARRG` - SATU bentuk untuk
// dua puluh lima jenis klausul (penyimpangan sadar 2), dibedakan
// `TreatyDescID`.
//
// Sembilan medan terakhir - `IDOccupation` sampai `MoreUsd` - hanya diisi
// jenis INDUK (18 jenis); jenis ANAK (7 jenis) menyimpannya KOSONG (AC 25).
//
// `ParentReinsTypeID` bernilai sentinel `"00"` untuk baris tanpa induk
// (`Activity/GetPeriode.xml` b889) - dibawa apa adanya.
type KlausulTreaty struct {
	ID                string
	TreatyYear        string
	TreatyYearID      string
	TreatyGroupID     string
	TreatyGroupName   string
	TreatyDescID      string
	TreatyDescName    string
	ReinsTypeID       string
	ReinsTypeName     string
	Layer             string
	LayerPart         string
	LayerPartType     string
	LayerType         string
	Kurs              *apd.Decimal
	TglUpdate         time.Time
	UserID            string
	Line              string
	Pct               *apd.Decimal
	PctMe             *apd.Decimal
	Ydcf              string
	Method            string
	TerritorialLimit  string
	ParentReinsTypeID string
	SpreadingOrder    string
	// Rp dan Usd DUA nilai terpisah (AC 49), bukan satu nilai + kode mata uang.
	Rp           *apd.Decimal
	Usd          *apd.Decimal
	IDOccupation string
	Occupation   string
	IDClause     string
	Clause       string
	TreatyLimit  *apd.Decimal
	CoinsMin     *apd.Decimal
	CoinsMax     *apd.Decimal
	MoreRp       *apd.Decimal
	MoreUsd      *apd.Decimal
}

// ParentReinsTypeTanpaInduk adalah sentinel warisan untuk klausul tanpa
// induk. `[terverifikasi]` `Treaty Contract Out/Activity/GetPeriode.xml` b889-890:
// `InputTreatyExclutionTreaty.ParentReinsTypeID = "00"`.
const ParentReinsTypeTanpaInduk = "00"
