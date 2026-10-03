// Package kontrak memuat antarmuka lintas modul - TANPA implementasi.
//
// Refactor bentuk B (30-09-2026): modul di bawah `modul/<nama>/` tidak pernah
// mengimpor modul lain. Bila satu modul memerlukan modul lain, yang ia kenal
// hanyalah antarmuka di paket ini; implementasinya disambung `cmd/api`.
// Tipe data yang menyeberang antarmuka ikut tinggal di sini.
package kontrak

import (
	"context"
	"errors"
)

// PembacaPolis membaca polis ringkas PremiumList Life - butir pl4/av.
//
// Disediakan PremiumList (`modul/premiumlistlife/backend/services.PembacaPolis`), dipakai
// Claim Life (layar `PolicyDataLife`, `Save to RNM`, ambang klaim). Galat untuk
// nomor yang tidak ada: `ErrPolisNomorTakDitemukan`.
type PembacaPolis interface {
	Ringkas(ctx context.Context, nomorPolis string) (PolisRingkas, error)
}

// ErrPolisNomorTakDitemukan - tidak ada polis dengan nomor itu.
var ErrPolisNomorTakDitemukan = errors.New(
	"repository: nomor polis tidak ditemukan di PremiumList Life")

// PolisRingkas adalah sepuluh medan `PolicyDataLife` yang PUNYA kolom.
//
// ⚠️ TIGA MEDAN LAYAR LAMA TIDAK ADA DI SINI, dan ketiadaannya dinyatakan
// alih-alih diisi teks kosong: `TanggalRespon`, `TanggalKonfirmasi`, dan
// `TanggalRealisasi` tidak punya kolom di migrasi 050-056 mana pun. Ketiganya
// properti halaman kerja Pega; di mana nilainya tinggal `[terbuka]`.
type PolisRingkas struct {
	NomorPolis       string
	Type             string
	MarketingName    string
	CedingCoName     string
	PolicyHolderName string
	BusinessName     string
	// DateReceived adalah `Date Received Email` di layar Claim Life.
	DateReceived *string
	Status       string
	StatusUpdate string
	// ProductNameID adalah KUNCI ambang batas hari (butir ba) - ia yang
	// dipakai membaca `PRODUCTINWARD_LIFE`.
	ProductNameID string
	ProductName   string
	// ProdKe adalah versi yang terbaca, supaya pemanggil dapat menyatakannya.
	ProdKe int

	// Tujuh medan layar lagi - butir av-2 (GILIRAN-11 paket 2), urut section.
	//
	//	TypeCeding/TypeCedingName `.TypeCeding`        "System Reinsurance" b9301
	//	ProRateType               `.ProRateType`       "Premium Method"     b9694
	//	WPC                       `.WPC`               "WPC"                b10149
	//	RetroName                 `.RetroName`         "Retro Name"         b10335 (TP/TR)
	//	SecurityReinsurer         `.SecurityReinsurer` "Security Reinsurer" b10617 (TP/TR)
	//	SobName                   `.SobName`           "SOB"                b11930
	TypeCeding        string
	TypeCedingName    string
	ProRateType       string
	WPC               *string
	RetroName         string
	SecurityReinsurer string
	SobName           string

	// Empat KUNCI yang `Save to RNM` pakai - bukan medan layar, sehingga
	// TIDAK menyeberang ke JSON `PolicyDataLife` (kontraknya dikunci dua sisi).
	//
	//	BusinessCode        `PolicyDataLife.BusinessCode`  langkah 10 (ContentNote), 13-20
	//	CedingCo            `PolicyDataLife.CedingCo`      langkah 11.1 `CARI1` (klaim ganda)
	//	RetroID             `PolicyDataLife.RetroID`       langkah 27
	//	SecurityReinsurerID `PolicyDataLife.SecurityReinsurerID` langkah 27
	BusinessCode        string
	CedingCo            string
	RetroID             string
	SecurityReinsurerID string
}
