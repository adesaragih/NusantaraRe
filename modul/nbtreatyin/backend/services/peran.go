package services

// Untuk apa berkas ini: TEMPAT BERPERAN - tiket 05 (spec §5.4; AC 12, 13, 81,
// 82, 91). Mekanisme dan pemetaannya tinggal di `models/peran_tempat.go`
// (konstanta kode, keputusan work owner K16 03-10-2026 - BUKAN tabel); di sini
// hanya disambungkan ke peran pelaku `inti.Pelaku.Peran` (workbasket akun).

import (
	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// TempatTanggalProduksi - `Section/ListSuggest` medan `.ProductionDate`.
const TempatTanggalProduksi = models.TempatTanggalProduksi

// tempat menghitung tampil/tidaknya setiap tempat bagi pelaku menurut
// `models.PemetaanPeranTempat` (kosong sampai IAM menjawab = semua tertunda).
func (l *Layanan) tempat(p inti.Pelaku) map[string]bool {
	return models.TempatTampil(models.PemetaanPeranTempat, p.PunyaPeran)
}
