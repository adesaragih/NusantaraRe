package services

// Layar daftar kontrak dari tabel WARISAN — keputusan pemilik proses
// 3 Oktober 2026.
//
// ⛔ Terjemahan bentuk tampil hidup DI SINI, bukan di repository dan bukan di
// peramban. Sebabnya dapat diperiksa: `gudangTiruan` membuatnya teruji tanpa
// satu pun koneksi Oracle, dan aturan yang teruji tanpa basis data adalah
// aturan yang masih teruji ketika basis datanya tidak terjangkau.
//
// ⛔ NILAI ASLINYA TIDAK DIUBAH. Tiap medan terjemahan berdampingan dengan
// `…Asli`-nya; yang menyelidiki selisih pemindahan membaca yang asli, yang
// membaca layar membaca terjemahannya.

import (
	"context"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// UkuranHalamanWarisan - baris per halaman layar daftar.
//
// ⚠️ Angka ini KEPUTUSAN KITA. Layar lama memakai penomoran Pega
// (`pyGridPaginator`) yang ukurannya tidak tertulis di Section mana pun, dan
// mengarangnya sebagai "angka dari sistem lama" akan keliru. 25 dipilih
// supaya 1.854 baris menjadi 75 halaman - cukup untuk menguji penomoran, dan
// cukup kecil agar satu halaman tidak pernah menjadi beban.
const UkuranHalamanWarisan = 25

// DaftarKontrakWarisan membaca satu halaman tabel warisan, sudah
// diterjemahkan untuk layar.
func (l *Layanan) DaftarKontrakWarisan(ctx context.Context, p inti.Pelaku, halaman int) (models.HalamanDaftarWarisan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.HalamanDaftarWarisan{}, err
	}
	// Halaman nol atau negatif BUKAN galat - ia permintaan yang dibulatkan
	// ke halaman pertama. Menolaknya membuat setiap pemanggil mengulang
	// pembulatan yang sama.
	if halaman < 1 {
		halaman = 1
	}

	total, err := l.gudang.CacahKontrakWarisan(ctx)
	if err != nil {
		return models.HalamanDaftarWarisan{}, err
	}
	hasil := models.HalamanDaftarWarisan{
		Baris:   []models.BarisDaftarWarisan{},
		Halaman: halaman,
		Ukuran:  UkuranHalamanWarisan,
		Total:   total,
	}
	offset := (halaman - 1) * UkuranHalamanWarisan
	if total > 0 && offset >= total {
		// Halaman di luar jangkauan menjawab halaman KOSONG, bukan galat:
		// pemakai yang menekan "berikutnya" sekali terlalu banyak tidak
		// sedang melakukan kesalahan.
		return hasil, nil
	}

	baris, err := l.gudang.DaftarKontrakWarisan(ctx, offset, UkuranHalamanWarisan)
	if err != nil {
		return models.HalamanDaftarWarisan{}, err
	}
	for i := range baris {
		baris[i].SifatProporsi = SifatProporsiTampil(baris[i].SifatProporsiAsli)
		baris[i].TanggalMulai = TanggalTampil(baris[i].TanggalMulaiAsli)
		baris[i].TanggalBerakhir = TanggalTampil(baris[i].TanggalBerakhirAsli)
	}
	hasil.Baris = baris
	return hasil, nil
}

// TanggalTampil mengubah `YYYYMMDD` menjadi `dd/mm/yy` — bentuk layar lama.
//
// ⛔ Yang BUKAN delapan angka dikembalikan APA ADANYA. Itu bukan kelonggaran,
// itu syarat: kolomnya `VARCHAR2` dan nullable, dan satu dari 1.854 baris
// memang kosong. Menebak tanggal untuk nilai yang tidak berbentuk tanggal
// berarti menampilkan hari yang tidak pernah ada di baris mana pun — dan
// pembacanya tidak punya cara tahu ia karangan.
func TanggalTampil(yyyymmdd string) string {
	s := strings.TrimSpace(yyyymmdd)
	if len(s) != 8 {
		return yyyymmdd
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return yyyymmdd
		}
	}
	// dd/mm/yy — dua digit tahun, seperti rujukan ("01/01/19").
	return fmt.Sprintf("%s/%s/%s", s[6:8], s[4:6], s[2:4])
}

// SifatProporsiTampil memberi spasi pada `NonProportional`.
//
// ⛔ HANYA nilai itu yang disentuh. `Proportional` tidak berubah, dan nilai
// lain apa pun dikembalikan apa adanya — termasuk kosong. Sapuan 3 Oktober
// 2026 menemukan tepat dua nilai di 1.854 baris (`Proportional` 1.079,
// `NonProportional` 775, nol NULL), tetapi kolomnya nullable dan teks bebas:
// baris ke-1.855 tidak terikat apa pun.
func SifatProporsiTampil(asli string) string {
	if strings.TrimSpace(asli) == "NonProportional" {
		return "Non Proportional"
	}
	return asli
}
