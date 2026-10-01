package models

// Keputusan kasus (tiket 08) - `ConfirmSection` (wadah `InputEDMLife.xml`
// b35518 `.IsJsonPolis=1`) → dua tombol `Submit` (b37494 Confirm, b38109
// Decline) → `AddHistorySuggest` → `IsLifeAccepted`.

// BarisRiwayat - satu baris riwayat keputusan (`T_VIEW_SUGGEST`), grid
// `ConfirmSection.xml` b1856: `Date` b1984 · `PIC` b2125 · `Status` b2263 ·
// `Comment` b2399.
type BarisRiwayat struct {
	No       int    `json:"no"`
	Tanggal  string `json:"date"`
	PIC      string `json:"pic"`
	Status   string `json:"status"`
	Komentar string `json:"comment"`
}

// RiwayatTulis - `AddHistorySuggest` 1 b233: `DateSuggest` b256
// `@CurrentDateTime()`, `PICSuggest` b332 (akun pelaku), `IsCedingConfirm`
// b353, `CommentSuggest` b374 `.Description`.
type RiwayatTulis struct {
	KasusID  string
	Waktu    string // `YYYY-MM-DD HH:MM:SS`
	PIC      string
	Status   string // `LabelKeputusanRiwayat`
	Komentar string
}

// ResmiKasus - kepala versi resmi sesudah `Confirm` (`GenerateNoEDM_Life` +
// `InsertJsonPolisLife_Act` 11.2-11.5): `NO_POLIS`, `PROD_KE`, `NO_ENDORS` =
// `PL_NUMBER_EDM`, `STATUSS` `Resolved-Completed`; peserta ber-`PL_NUMBER_EDM`,
// `STATUS_OLD` (`StatusLama`), `STATUS` (`StatusJenis`).
type ResmiKasus struct {
	ID          string
	NomorPolis  string
	ProdKe      int
	Nomor       string
	StatusJenis string
}

// BatasKomentar - `T_VIEW_SUGGEST.COMMENT_SUGGEST VARCHAR2(255)`.
const BatasKomentar = 255
