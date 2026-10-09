package kontrak

// Kontrak Claim Non Prop untuk Komite Claim Non Prop - perintah work owner 09-10-2026 (komite Non Prop ikut pola Claim
// Prop: tanpa menu komite, kasus komite dibuka dari inbox Claim Non Prop, tangga workbasket).
//
// Komite Claim Non Prop memutus SATU baris akseptasi klaim treaty inward non proporsional (XoL). Kasus komite `KMTNP-`
// dilahirkan Claim Non Prop (`CreateChildKomiteCNP_Act`); Komite membaca klaim induknya dan, lewat `KomitePostAdjustment`
// (korpus `Komite Claim Non Prop`), menulis kembali ke kasus klaim itu (`TempMainWork`, Obj-Open-By-Handle
// `pxCoverInsKey`). Tulisan ke tabel klaim Non Prop HANYA lewat kontrak ini, di dalam transaksi Komite sendiri.
//
// Bentuk data sama dengan `KlaimTreatyKomite` (`KlaimTreaty`, `UbahanKlaimTreaty`); daftar putih jalurnya berbeda dan
// `UbahanKlaimTreaty.FacRetro` tidak dipakai (XML Non Prop tanpa Fac Retro). Disediakan `claimnonprop`
// (`services.KlaimUntukKomite`), dipakai `komiteclaimnonprop`.

import (
	"context"

	"nusantarare/inti/backend/db"
)

// KlaimTreatyNonPropKomite adalah yang Komite Claim Non Prop butuhkan dari Claim Non Prop.
type KlaimTreatyNonPropKomite interface {
	// BacaKlaimTreaty membaca halaman kasus klaim induk seperti dibuka Claim Non Prop beserta posisi baris akseptasi
	// `adjID` (ID baris T_CLAIM_ADJUSTMENT). `tx` nil = baca di luar transaksi (layar).
	BacaKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, adjID string) (KlaimTreaty, error)
	// KunciKlaimTreaty = `KomitePostAdjustment` S5 (Obj-Open-By-Handle `pyWorkPage.pxCoverInsKey`, `Lock=true`).
	KunciKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID string) error
	// TulisBalikKlaimTreaty menerapkan ubahan Komite atas kasus klaim induk lalu menyimpannya; jalur di luar
	// `JalurHeaderKomiteNonProp` / `PropAdjustmentKomiteNonProp` -> `ErrUbahanKlaimTreatyTidakSah`.
	TulisBalikKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, adjID string, u UbahanKlaimTreaty) error
}

// JalurHeaderKomiteNonProp - jalur header kasus klaim yang boleh ditulis Komite Claim Non Prop, beserta langkah XML
// penulisnya (`KomitePostAdjustment`). Tidak dimuat: `ClaimData.IsFInalAccXOL` (S14.17), `stsReject` (S14.19), dan
// `ClaimData.IsSubjectivity` (S18) - tanpa pembaca di korpus Non Prop dan tanpa kolom di katalog Claim Non Prop
// (`KonversiKlaim_Act` menerima STSREJECT "1" tertulis mati).
var JalurHeaderKomiteNonProp = map[string]string{
	"IsCloseFile":   "S14.12 (<- Adjustment.IsProposeClose)",
	"CNPStatusCase": "S14.17 (CLAIM ACCEPTED) dan S20 (CLAIM REJECTED)",
}

// PropAdjustmentKomiteNonProp - properti baris akseptasi yang boleh ditulis Komite Claim Non Prop, beserta langkah XML
// penulisnya.
var PropAdjustmentKomiteNonProp = map[string]string{
	"AcceptedNo":       "KomitePostAdjustment S14.13",
	"AcceptedDate":     "KomitePostAdjustment S14.13",
	"AcceptanceStatus": "KomitePostAdjustment S14.13 (1) dan S11.4 (2)",
	"IsKomite":         "KomitePostAdjustment S16 (subjectivity -> 0)",
	"IsSubjectivity":   "KomitePostAdjustment S18",
	"SubjectivityNote": "KomitePostAdjustment S18",
	"NoAccount":        "HitServiceToKasirKMT_Act S7 (KomitePostAdjustment S19.3; angka saja)",
	"IDOfBank":         "HitServiceToKasirKMT_Act S12-S13 (KomitePostAdjustment S19.3)",
}
