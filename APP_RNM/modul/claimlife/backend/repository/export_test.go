package repository

// Ekspor KHUSUS UJI untuk paket `repository_test` (uji `db` tidak dapat berpaket `repository`: `uji/skemauji`
// mengimpor paket ini - siklus impor).

// VersiTerakhirDibangun, KasusVersiTerakhir, PilihAcuan - tabel kasus dan acuan aturan versi terakhir
// (`versiterakhir_test.go`) bagi `versiterakhir_db_test.go`.
const VersiTerakhirDibangun = versiTerakhirDibangun

var (
	KasusVersiTerakhir = kasusVersiTerakhir
	PilihAcuan         = pilihAcuan
)
