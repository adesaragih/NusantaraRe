package models

// RencanaSimpan - semua yang SATU penekanan tombol tulis (`Save`, `Submit`,
// `Actions`, `Decline offer`) tulis, dalam satu transaksi.
type RencanaSimpan struct {
	// ID kontrak; kosong = kontrak BARU (pengenal dibuat repository).
	ID string
	// Dokumen `TreatyIn` UTUH yang mendarat di `T_TREATY_*` dan kepala
	// `TREATY_IN`.
	Dokumen map[string]any
	// Kurs tahun kontrak — `nil` = grid Rate of Exchange tidak disimpan.
	Kurs *KursSimpan
	// Operator (`USERID` kurs) dan stempel waktu Pega (`DATEIU`/`DATEIN`).
	Operator string
	Stempel  string
}

// RencanaPenyesuaian - satu penekanan tombol tulis layar ADJUSTMENT (EDM):
// kepala `TREATY_IN_EDM` dan kedua sisi dokumennya di `T_TREATY_*`, dalam
// satu transaksi.
type RencanaPenyesuaian struct {
	// ID penyesuaian (`1000080/R01`). Wajib — dibuat layar saat `Choose`.
	ID string
	// Draf - penyesuaian BARU dari `Choose`, belum pernah disimpan: kepala
	// di-`INSERT`, dan pengenal yang sudah ada DITOLAK.
	Draf bool
	// Baru - halaman `TreatyIn` (sisi New), mendarat di `MASTERID = ID`.
	Baru map[string]any
	// Lama - halaman `OLDDATA` (sisi Old), mendarat di `ID + #LAMA`.
	// `nil` = tidak disentuh (penyesuaian tersimpan: sisi Old baca saja).
	Lama map[string]any
}
