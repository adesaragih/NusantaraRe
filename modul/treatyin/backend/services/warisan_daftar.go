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
	s := TanggalDelapanDigit(yyyymmdd)
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

// TanggalDelapanDigit memotong cap waktu Pega menjadi tanggalnya saja.
//
// ⛔ LAHIR DARI CACAT NYATA, 6 Oktober 2026. Grid `Rate of Exchange`
// menampilkan `20250701T075400.000 GMT` apa adanya, sebab `TanggalTampil`
// menuntut PERSIS delapan digit dan mengembalikan apa adanya untuk yang lain.
// `TREATYEXCHANGEYEARLY.STARTDATE` menyimpan cap waktu penuh, bukan tanggal
// polos — dan tabel pendaratan menyimpan tanggal polos, sehingga bedanya
// tidak terlihat sampai sumbernya berganti.
//
// ⚠️ SESEMPIT bentuknya: delapan digit yang DIIKUTI `T`. Teks lain lewat
// apa adanya, sehingga nilai yang aneh tetap terlihat aneh alih-alih
// dikarang menjadi tanggal.
func TanggalDelapanDigit(nilai string) string {
	s := strings.TrimSpace(nilai)
	// `YYYYMMDD` + `T` + `hhmmss` = 15 aksara, dan keenam aksara jam WAJIB
	// angka. Cap waktu yang TERPOTONG (`20250101T00`) tidak diterima:
	// `TestTanggalPanjangTidakMengarang` menuntutnya lewat apa adanya, dan
	// tuntutan itu benar — nilai yang aneh harus tetap terlihat aneh.
	if len(s) < 15 || s[8] != 'T' {
		return s
	}
	for _, r := range s[9:15] {
		if r < '0' || r > '9' {
			return s
		}
	}
	return s[:8]
}

// TanggalTampilPanjang mengubah `YYYYMMDD` menjadi `dd/mm/yyyy`.
//
// ⛔ BENTUK KEDUA, dan ia perlu — bukan kerapian. Dokumen desain 5 Oktober
// 2026 memperlihatkan DUA bentuk tanggal di layar yang sama, dan bedanya
// disengaja:
//
//	gambar 01  medan kepala `Commencement`   `01/01/2025`  dd/mm/yyyy
//	gambar 01  grid Rate of Exchange          `01/01/25`    dd/mm/yy
//	gambar 29  grid EGNPI `As Date`           `18/01/25`    dd/mm/yy
//	gambar 29  rincian EGNPI `As At`          `18/01/2025`  dd/mm/yyyy
//	gambar 38  grid Installment `Due Date`    `18/01/2024`  dd/mm/yyyy
//
// ⭐ Aturannya terbaca dari kelima contoh itu: BARIS GRID YANG TERLIPAT
// memakai dua digit tahun; MEDAN — kepala maupun rincian yang terbuka —
// memakai empat. Installment tidak melanggarnya: kolom `Due Date`-nya medan
// isian di dalam baris yang sudah terbuka, bukan teks baris terlipat.
//
// ⚠️ Satu bentuk dipakai di dua tempat yang di Pega berbeda akan terbaca
// sebagai layar yang bukan layar yang sama — itu sebabnya keduanya ada.
//
// ⛔ Yang BUKAN delapan angka dikembalikan APA ADANYA, aturan yang sama
// dengan `TanggalTampil`.
func TanggalTampilPanjang(yyyymmdd string) string {
	s := TanggalDelapanDigit(yyyymmdd)
	if len(s) != 8 {
		return yyyymmdd
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return yyyymmdd
		}
	}
	return fmt.Sprintf("%s/%s/%s", s[6:8], s[4:6], s[0:4])
}

// CaraPembukuanTampil menerjemahkan nilai tersimpan `AccountingMode` dan
// `AccountingModeNonProp` menjadi label layar.
//
// ⛔ Ketiga padanan ini DIBACA DARI TANGKAPAN LAYAR SISTEM BERJALAN, bukan
// dari ekspor — rule `Rule-Obj-Property` keduanya tidak ikut diekspor
// (`pyListSource = associated`, dan korpus nol folder `Property`):
//
//	gambar 01  kontrak 1001846 prop     tersimpan `underwriting` -> "Underwriting Year"
//	gambar 26  kontrak 1001841 non-prop tersimpan `loss`         -> "Loss Occuring"
//
// `accounting` -> "Accounting Year" adalah pasangan yang pemilik proses
// nyatakan sendiri dalam briefing 5 Oktober 2026.
//
// ⚠️ EJAANNYA "Loss Occuring" dengan SATU `r`, persis seperti di layar.
// Memperbaikinya menjadi "Occurring" membuat layar ini berbeda dari layar
// yang orang hafal.
//
// ⛔ `risk` TIDAK diterjemahkan, dan itu keputusan: labelnya tidak ada di
// tangkapan layar mana pun dari ke-43, dan tidak ada di korpus. Satu-satunya
// kemunculan "Risk Attaching" di `D:\XML_NURE` ada di MEMO berbahasa
// Indonesia pada `Claim Non Prop\Activity\SetEndDate_Act.xml` — catatan
// pengembang di modul lain, bukan label properti ini. Menebaknya akan
// menuliskan karangan ke layar yang orang percayai; nilainya tampil apa
// adanya, dan pertanyaannya di `PERTANYAAN-TERBUKA-LAYAR-PEGA.md`.
func CaraPembukuanTampil(tersimpan string) string {
	switch strings.TrimSpace(tersimpan) {
	case "underwriting":
		return "Underwriting Year"
	case "accounting":
		return "Accounting Year"
	case "loss":
		return "Loss Occuring"
	}
	// Termasuk `risk` — 21 dari 1.854 kontrak. Lihat komentar di atas.
	return tersimpan
}

// BordereauxTampil menerjemahkan nilai tersimpan `Bordeaux` menjadi labelnya.
//
// ⭐ SATU padanan yang terukur: gambar 01 memperlihatkan kontrak 1001846
// (tersimpan `reporting`) tampil sebagai "Reporting".
//
// ⛔ `nonreporting` TIDAK diterjemahkan. Ia tidak muncul di satu pun dari 43
// tangkapan layar — kedua kontrak contoh bernilai `reporting` atau tidak
// punya medannya. "Non Reporting" terdengar jelas dan tetap tebakan; 687
// kontrak memakainya, dan 687 layar yang salah lebih buruk daripada 687
// layar yang jujur.
func BordereauxTampil(tersimpan string) string {
	if strings.TrimSpace(tersimpan) == "reporting" {
		return "Reporting"
	}
	return tersimpan
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

// Domain TERSIMPAN dropdown kepala — terukur atas 1.854 dokumen
// (`frontend/labels.ts`, `bordereauxNilai` · `caraPembukuanNilai` ·
// `caraPembukuanNonPropNilai`): `reporting` 1.164 · `nonreporting` 687;
// `underwriting` 1.110 · `accounting` 741; `loss` 1.830 · `risk` 21.
var (
	domainBordereaux           = []string{"reporting", "nonreporting"}
	domainCaraPembukuan        = []string{"underwriting", "accounting"}
	domainCaraPembukuanNonProp = []string{"loss", "risk"}
	// Nilai yang `Activity/TreatyInSetReport.xml` kenal (langkah 9–12);
	// tersimpan: quarter 1.848 · month 2 · other 1 · kosong 4.
	domainPeriodePelaporan = []string{"quarter", "half", "month", "other"}
)

// OpsiKepalaKontrak menyusun pilihan dropdown kepala: NILAI tersimpan
// berpasangan LABEL tampil.
//
// ⛔ Labelnya dari penerjemah yang SAMA (`CaraPembukuanTampil`,
// `BordereauxTampil`) — nol penerjemah kedua. Nilai yang penerjemahnya tidak
// kenal (`risk`, `nonreporting`) berlabel nilainya sendiri, persis seperti
// penerjemahnya mengembalikannya.
//
// ⚠️ Sebab bug "Underwriting Year (tidak ada di daftar referensi)": form
// mengisi dropdown dengan LABEL hasil terjemahan, sementara pilihannya NILAI
// tersimpan. Dengan pasangan ini dropdown memegang nilai, dan menampilkan label.
func OpsiKepalaKontrak() models.OpsiKepala {
	susun := func(domain []string, tampil func(string) string) []models.Opsi {
		out := make([]models.Opsi, 0, len(domain))
		for _, v := range domain {
			out = append(out, models.Opsi{Nilai: v, Label: tampil(v)})
		}
		return out
	}
	return models.OpsiKepala{
		Bordereaux:           susun(domainBordereaux, BordereauxTampil),
		CaraPembukuan:        susun(domainCaraPembukuan, CaraPembukuanTampil),
		CaraPembukuanNonProp: susun(domainCaraPembukuanNonProp, CaraPembukuanTampil),
		PeriodePelaporan:     susun(domainPeriodePelaporan, PeriodePelaporanTampil),
	}
}

// OpsiKepala - rute `GET /warisan/opsi-kepala`: pilihan dropdown kepala.
func (l *Layanan) OpsiKepala(p inti.Pelaku) (models.OpsiKepala, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.OpsiKepala{}, err
	}
	return OpsiKepalaKontrak(), nil
}
