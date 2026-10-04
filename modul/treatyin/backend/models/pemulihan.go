package models

// Bentuk pemulihan limit — tiket 32.
//
// ⛔ PersenPemulihan dan PersenTambahan DUA ruas yang berbeda, dan itu POKOK
// tiket 32. Sistem lama menulis keduanya tetap `"100"` di dalam kode
// (`SetReinstatementPct.xml`), sehingga ketentuan pasar yang biasa - pemulihan
// pertama gratis, kedua berbayar penuh - tidak punya tempat untuk dinyatakan.
// Selama keduanya selalu 100, satu baris dan sepuluh baris menghasilkan angka
// yang sama, dan tidak ada yang pernah melihat batasnya.

import "github.com/cockroachdb/apd/v3"

// PemulihanLimit adalah satu pemulihan pada sebuah layer.
//
// PersenTambahan penunjuk sebab kolomnya boleh kosong (`413_layer.sql`):
// kosong berarti tarif premi pemulihannya BELUM dinyatakan - bukan nol, dan
// bukan seratus.
type PemulihanLimit struct {
	ID        int64 `json:"id,omitempty"`
	IDLayer   int64 `json:"idLayer,omitempty"`
	NomorUrut int64 `json:"nomorUrut"`
	// Porsi limit yang dipulihkan, persen. Wajib isi di DDL.
	PersenPemulihan *apd.Decimal `json:"persenPemulihan"`
	// Tarif premi pemulihan, persen. Boleh kosong.
	PersenTambahan *apd.Decimal `json:"persenTambahan,omitempty"`
	Catatan        string       `json:"catatan,omitempty"`
}
