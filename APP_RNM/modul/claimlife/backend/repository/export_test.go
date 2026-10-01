package repository

// Ekspor KHUSUS UJI untuk paket `repository_test` (uji `db` tidak dapat berpaket `repository`: `uji/skemauji`
// mengimpor paket ini - siklus impor).

// KasusVersiTerakhir, PilihAcuan - tabel kasus dan acuan aturan versi terakhir (`versiterakhir_test.go`) bagi
// `versiterakhir_db_test.go`.
var (
	KasusVersiTerakhir = kasusVersiTerakhir
	PilihAcuan         = pilihAcuan
)
