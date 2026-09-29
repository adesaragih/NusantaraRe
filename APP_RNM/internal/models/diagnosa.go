package models

// Diagnosa per peserta - butir bd.
//
// Untuk apa berkas ini: satu baris `.DiagnoseList`, daftar yang hidup DI
// DALAM halaman peserta. Aturan murni saja; yang menyentuh Oracle ada di
// `repository/diagnosa.go`.
//
// Dibaca sesudah: penyakit.go. Yang di sana adalah KATALOG penyakit (tabel
// warisan 97.586 baris, dibaca saja); yang di sini adalah diagnosa yang
// BENAR-BENAR dilekatkan pada seorang peserta klaim. Dua hal berbeda, dan
// satu-satunya yang menghubungkannya adalah tombol `Choose`.
//
// Istilah:
//   - diagnosa : satu baris pada peserta - ICD, nama penyakit, kelompok.
//   - penyakit : satu baris katalog `DISEASE_LIFE`. Sumber, bukan milik.

import (
	"strings"

	"nusantarare/inti/kontrak"
)

// Diagnosa adalah satu baris `.DiagnoseList` milik seorang peserta.
//
// `[terverifikasi]` `Section/ClaimLifeDetailGCNM.xml` b3923
// `pyPageListProperty .DiagnoseList`, kelas baris `ASM-FW-GISFW-Data-
// DiagnoseLife` b3915, disajikan `RepeatGrid` b3926.
//
// ⛔ BANYAK per peserta, dan buktinya tombolnya - bukan bentuk gridnya.
// Grid DAPAT berarti tampilan satu baris; grid ber-`Add` b4690 DAN
// ber-`Delete` b6160 tidak dapat. Ronde sebelumnya berhenti pada gridnya dan
// menyimpulkan "satu lawan banyak, tidak dapat diputuskan" (OQ-K.2) -
// jawabannya dua baris di bawah tempat pembacaan itu berhenti.
// ⛔ SETIAP medan bertag JSON eksplisit, termasuk yang "sudah benar"
// kalau ditulis huruf kecil. Tanpa tag, Go mengirim `KodeICD` sedangkan
// layar membaca `kodeIcd`, dan hasilnya `undefined` yang tampil sebagai sel
// kosong - tanpa satu pun galat, di kedua sisi. Itu bentuk cacat yang sudah
// LIMA kali terjadi di modul ini, dan tiap kali tiap sisi hijau sendirian.
// Dikunci `TestNamaJSONDiagnosaDikunci`.
type Diagnosa struct {
	// ID dari `SEQ_CLAIMLF_DIAGNOSE` (ADR-0006).
	//
	// ⚠️ Angka, bukan teks - berbeda dengan pengenal peserta dan pengenal
	// klaim. Ia milik KITA: tidak ada padanannya di sistem lama, sebab di
	// Pega baris ini hanya punya subscript di dalam halaman induknya.
	ID int64 `json:"id"`
	// PesertaID menunjuk `T_CLAIMLF_PREMIUMLIST_DETAIL.ID`.
	PesertaID string `json:"pesertaId"`
	// Urutan adalah posisi baris di grid, mulai 1.
	//
	// ⛔ Padanan `.pxListSubscript`. Ia bukan hiasan: `SetSTS_Reject.xml`
	// b241 memutar `.DiagnoseList` dan b345 menandainya `EMBEDDED`, jadi
	// urutan yang orang lihat adalah urutan yang rule itu tempuh.
	Urutan int `json:"urutan"`
	// KodeICD - kolom `ICD_CODE`, grid `.ICDCODE` b5616.
	//
	// ⛔ TEKS, dan read-only di layar (b5566 `pyEditOptions Read-only`):
	// nilainya datang dari baris katalog yang dipilih, tidak diketik.
	KodeICD string `json:"kodeIcd"`
	// Nama - kolom `DISEASE`, grid `.DISEASE` b5422, read-only b5374.
	Nama string `json:"nama"`
	// GroupDiagnose - kolom `GROUP_DIAGNOSE`, grid `.GROUPDIAGNOSE` b5860.
	//
	// ⛔ `[terbuka - OQ-L]` DAFTAR PILIHANNYA TIDAK ADA. Dropdown b5863
	// ber-`pyListSource associated`, artinya daftarnya hidup pada rule
	// properti `.GROUPDIAGNOSE` kelas `Data-DiagnoseLife` - dan rule itu
	// tidak ada di ekspor maupun di katalog DEV. Nilainya karena itu
	// diterima apa adanya sebagai teks dan TIDAK diperiksa terhadap daftar
	// mana pun: daftar karangan akan menolak nilai yang sah di sistem lama.
	GroupDiagnose string `json:"groupDiagnose"`
	// KodeStatus adalah CERMIN `STS_REJECT` PESERTA, bukan keputusan sendiri.
	//
	// `[terverifikasi]` `Activity/SetSTS_Reject.xml` b257-258:
	// `.STS_REJECT = Primary.STS_REJECT`, berulang `EMBEDDED` b345 atas
	// halaman langkah `.DiagnoseList` b241, tanpa prasyarat (b325 `WhenTrue`
	// 2 / `WhenFalse` 2 - keduanya LANJUT).
	KodeStatus string `json:"kodeStatus"`
}

// BatasGroupDiagnose adalah lebar kolom `GROUP_DIAGNOSE` (migrasi 018).
//
// ⚠️ `[terbuka - OQ-L]` Angkanya lebar KOLOM, bukan aturan dagang: tanpa
// daftar pilihan, satu-satunya hal yang dapat kami katakan tentang nilainya
// adalah bahwa ia harus muat. Ketika OQ-L dijawab, pemeriksaan daftar
// menggantikan pemeriksaan panjang ini - dan panjangnya tetap berlaku.
const BatasGroupDiagnose = 255

// BatasKodeICDDiagnosa dan BatasNamaDiagnosa mengikuti tabel sumbernya.
//
// `[data DBA - katalog DEV 27-09-2026]` `POOLDATA.DISEASE_LIFE`:
// `ICD_CODE VARCHAR2(100)`, `DISEASE VARCHAR2(1000)`. Lebar sumber dipakai
// apa adanya - menyempitkannya "karena isi terpanjangnya 7 dan 290" berarti
// menebak bahwa katalog tidak akan pernah tumbuh.
const (
	BatasKodeICDDiagnosa = 100
	BatasNamaDiagnosa    = 1000
)

// DiagnosaTerkunci menjawab apakah baris diagnosa tidak boleh diubah lagi.
//
// `[terverifikasi]` `Section/ClaimLifeDetailGCNM.xml`, SATU kalimat di EMPAT
// tempat - keempatnya `pyDisabledWhen`:
//
//	b4682  `Add`
//	b5059  `Find Disease`
//	b5870  `GROUPDIAGNOSE`
//	b6152  `Delete`
//
//	.STS_REJECT=='1' || .STS_REJECT=='2'
//
// ⛔ YANG DIUJI ADALAH `STS_REJECT` PESERTA, bukan milik baris diagnosa.
// Keempat kontrol itu berdiri di dalam grid, tetapi `.` di sana menunjuk
// halaman PESERTA yang memuat gridnya - dan kalaupun ia menunjuk baris
// diagnosa, hasilnya sama: `SetSTS_Reject` menyalin nilai peserta ke setiap
// barisnya, jadi keduanya selalu sama nilainya. Pemanggil karena itu
// menyerahkan `STS_REJECT` PESERTA, dan itulah yang namanya sebut.
//
// ⚠️ Memakai KodeAksep/KodeDitolak, bukan literal "1"/"2", supaya kalimat
// ini dan mesin status tidak dapat bergeser sendiri-sendiri. Keduanya
// dikunci `TestKodeStatusLiteralHanyaDiModels`.
func DiagnosaTerkunci(stsRejectPeserta string) bool {
	k := strings.TrimSpace(stsRejectPeserta)
	return k == kontrak.KodeAksep || k == kontrak.KodeDitolak
}

// UrutanBerikutnya mengembalikan nomor urut baris baru.
//
// `Add` b4700 menyisipkan `After` (b4715 `pyPosition`), jadi baris baru
// selalu di EKOR. Mulai dari 1: nol adalah nilai yang tidak dapat dibedakan
// dari "belum diisi" (ADR-U-0027).
func UrutanBerikutnya(daftar []Diagnosa) int {
	tertinggi := 0
	for _, d := range daftar {
		if d.Urutan > tertinggi {
			tertinggi = d.Urutan
		}
	}
	return tertinggi + 1
}

// RapatkanUrutan menomori ulang daftar 1..n menurut urutannya sekarang.
//
// ⛔ Dipakai SESUDAH penghapusan. `Delete` b6170 `deleteRow` di Pega
// menghapus anggota PageList, dan subscript anggota sesudahnya bergeser
// sendiri - tidak ada lubang. Kolom `URUTAN` tidak bergeser sendiri, jadi
// yang di Pega gratis harus di sini dikerjakan.
//
// ⚠️ Daftar diandaikan SUDAH terurut oleh pemanggilnya (repository membaca
// `ORDER BY URUTAN, ID`). Fungsi ini tidak mengurutkan ulang: mengurutkan di
// dua tempat berarti dua definisi "urutan", dan yang kedua akan diam-diam
// menang.
func RapatkanUrutan(daftar []Diagnosa) []Diagnosa {
	rapat := make([]Diagnosa, len(daftar))
	copy(rapat, daftar)
	for i := range rapat {
		rapat[i].Urutan = i + 1
	}
	return rapat
}

// PotongDiagnosa memangkas nilai yang melebihi lebar kolomnya - dan TIDAK
// memotong: ia melaporkan.
//
// ⛔ Pemotongan diam-diam adalah kehilangan data yang terlihat seperti
// keberhasilan. Nama diagnosa yang terpenggal di karakter ke-1000 tetap
// tersimpan, tetap tampil, dan tidak ada satu pun yang berbunyi - sampai
// seseorang membandingkannya dengan katalog. Yang dikembalikan karena itu
// nama kolom yang kepanjangan, bukan nilai yang sudah dipendekkan.
// ⚠️ `len()` menghitung BYTE, bukan rune - dan itu disengaja. Oracle
// `VARCHAR2(n)` bawaannya pun berukuran byte, jadi ini justru ukuran yang
// sama dengan kolomnya. Bila kelak skema memakai semantik CHAR, pemeriksaan
// ini menolak sedikit LEBIH AWAL daripada Oracle - arah yang benar: menolak
// nilai yang sebenarnya muat hanya merepotkan, sedangkan menerima nilai yang
// tidak muat menggagalkan transaksi di tempat yang jauh dari sebabnya.
func PotongDiagnosa(kodeICD, nama, grup string) (string, bool) {
	switch {
	case len(kodeICD) > BatasKodeICDDiagnosa:
		return "ICD_CODE", false
	case len(nama) > BatasNamaDiagnosa:
		return "DISEASE", false
	case len(grup) > BatasGroupDiagnose:
		return "GROUP_DIAGNOSE", false
	}
	return "", true
}

// tahapBergridPeserta adalah tahap yang layarnya memuat grid peserta - dan
// karena itu satu-satunya tempat grid diagnosa dapat dibuka.
//
// `[terverifikasi]` `ViewClaimDetailLifeGCNM` adalah `pyEditAction` grid
// peserta di TIGA section, ketiganya atas daftar yang sama
// (`pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail`):
//
//	Section/InputOSClaimLife.xml        b18252   grid b15924   Outstanding
//	Section/MedicalCheckClaimLife.xml   b17416   grid b15764   Medical Check
//	Section/InputAkseptasiClaimLife.xml b17387   grid b15735   Claim Analis
//
// ⛔ TIGA, dan `MedicalCheckClaimLife` TERMASUK - walau ia tidak muncul
// di daftar rule yang MEMUAT `ClaimLifeDetailGCNM`. Ia tidak memuatnya; ia
// MEMBUKANYA. Berhenti pada daftar pemuat berarti menyimpulkan bahwa Medical
// Advisor tidak dapat menyunting diagnosa - padahal dialah yang paling masuk
// akal mengisinya. Bentuk kekeliruan yang sama dengan `Close Claim`, OQ-H,
// dan OQ-K.2: berhenti pada X lalu menyimpulkan tentang Y.
//
// ⛔ Input Register TIDAK termasuk: nol grid peserta di sana. Peserta
// baru dipilih di tahap itu, belum disunting.
var tahapBergridPeserta = map[Tahap]bool{
	TahapOutstanding:  true,
	TahapMedicalCheck: true,
	TahapClaimAnalis:  true,
}

// TahapBergridPeserta menjawab apakah sebuah tahap membuka grid peserta.
func TahapBergridPeserta(t Tahap) bool { return tahapBergridPeserta[t] }

// TahapPembukaGridPeserta menyebut ketiganya, untuk pesan galat.
//
// Urutannya tangga kerja, bukan urutan peta - peta di Go tidak berurutan,
// dan pesan galat yang berganti susunan setiap kali sulit dipercaya.
func TahapPembukaGridPeserta() []string {
	return []string{
		TahapOutstanding.String(),
		TahapMedicalCheck.String(),
		TahapClaimAnalis.String(),
	}
}
