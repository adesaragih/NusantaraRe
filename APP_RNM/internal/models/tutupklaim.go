package models

// Gerbang `Close Claim` — kelompok Detail & Tutup.
//
// Untuk apa berkas ini: sebelum sebuah klaim boleh ditutup, Pega memeriksa
// seluruh pesertanya satu per satu. Bila ada yang belum diaksep, penugasannya
// TIDAK diselesaikan dan pemakai diberi tahu peserta mana yang menahan.
//
// Sumber, dibaca sebagai pohon langkah utuh:
//
//	`Activity/ProtectCloseClaim_act.xml`
//	  b236 `Property-Set`          `ProtectLife.CARI1 = ""`
//	  b370 `Property-Set` halaman `pyWorkPage...PremiumListDetail`
//	       b812 ULANG(EMBEDDED) atas seluruh peserta
//	       b398 `local.idx  = .pxListSubscript`
//	       b445 `local.name = .NAME_OF_INSURED`
//	    b484 `Property-Set`  prasyarat b608 `.STS_REJECT!=1`
//	                         `WhenTrue=2` lanjut / `WhenFalse=3` lewati
//	         b511 `ProtectLife.CARI1 = 1`
//	         b557 `local.errmsg` = nama + " is not approved yet, on list " + idx
//	    b646 `Page-Set-Messages`  prasyarat b756 `ProtectLife.CARI1==1`
//	  b836 `Call FinishAssignment` prasyarat b992 `ProtectLife.CARI1==""`
//
// ⚠️ `STS_REJECT == 1` berarti DIAKSEP, bukan ditolak - lihat KodeAksep.
// Nama kolomnya menyesatkan dan itu warisan; maknanya diambil dari kodenya.
//
// Dibaca sesudah: statusklaim.go.

import "fmt"

// BarisTutup adalah satu peserta sebagaimana dilihat gerbang tutup.
//
// Hanya tiga hal yang diperlukan, dan tidak satu pun di antaranya nama orang.
type BarisTutup struct {
	// Urutan adalah nomor barisnya di daftar, mulai 1 - padanan
	// `.pxListSubscript` (b398). Ia muncul di pesan, jadi ia harus nomor
	// yang DILIHAT pemakai di layar, bukan indeks slice.
	Urutan int
	// NomorSertifikat menunjuk peserta itu di pesan.
	NomorSertifikat string
	// KodeStatus adalah cermin `STS_REJECT` baris itu.
	KodeStatus string
}

// Penghalang adalah satu peserta yang menahan penutupan.
type Penghalang struct {
	Urutan          int
	NomorSertifikat string
}

// bentukPesanBelumAksep adalah kalimat Pega, dengan SATU penggantian.
//
// ⛔ PENYIMPANGAN SADAR, dan sebabnya bukan gaya penulisan. Pega menyusun
// kalimatnya dari `.NAME_OF_INSURED` (b445, b557). Kita SENGAJA tidak
// menyimpan nama tertanggung: `kolomSalin` tiket 02 meninggalkan
// `NAME_OF_INSURED` dan `POLICY_HOLDER` justru supaya data pribadi tidak
// berganda ke tabel klaim. Nama itu karena itu TIDAK ADA untuk dipakai di
// sini, dan mengambilnya kembali dari tabel warisan hanya demi sebuah pesan
// berarti membatalkan keputusan itu.
//
// Penggantinya nomor sertifikat: ia menunjuk baris yang sama persis, memang
// kita simpan, dan tidak memuat nama siapa pun. Sisa kalimatnya VERBATIM,
// supaya orang dapat mencari teks yang sama di kedua sistem.
const bentukPesanBelumAksep = "%s is not approved yet, on list %d"

// Pesan menyusun kalimat yang dilihat pemakai.
func (p Penghalang) Pesan() string {
	return fmt.Sprintf(bentukPesanBelumAksep, p.NomorSertifikat, p.Urutan)
}

// PenghalangTutupKlaim mengembalikan SELURUH peserta yang menahan penutupan.
//
// ⛔ Seluruhnya, bukan yang pertama. Pega memasang pesannya DI DALAM loop
// (b646 adalah sublangkah b370), sekali untuk tiap iterasi yang tertandai.
// Melaporkan satu saja memaksa pemakai menutup klaim berulang kali dan
// menemukan satu penghalang baru setiap kali - perbaikan yang terasa seperti
// hukuman.
func PenghalangTutupKlaim(baris []BarisTutup) []Penghalang {
	var out []Penghalang
	for _, b := range baris {
		// b608 `.STS_REJECT != 1` apa adanya. ⛔ BUKAN "== 0": status
		// DITOLAK ("2") juga bukan 1, jadi ia menahan pula. Memperlakukan
		// "ditolak" sebagai "selesai" akan menutup klaim yang barisnya belum
		// diputus ulang, dan XML tidak pernah mengatakannya.
		// ⛔ Dibandingkan PERSIS, tanpa TrimSpace. Ronde pertama memakai
		// TrimSpace "untuk aman", dan itu justru MELONGGARKAN gerbang uang:
		// " 1 " akan menutup klaim di sini padahal Pega - yang membandingkan
		// `.STS_REJECT != 1` apa adanya - menahannya. Paket ini pun sudah
		// menyatakan aturannya di StatusBarisDariKode: perbandingan atas
		// TEKS, dan "00" BUKAN "0" (ADR-U-0022).
		if b.KodeStatus == KodeAksep {
			continue
		}
		out = append(out, Penghalang{Urutan: b.Urutan, NomorSertifikat: b.NomorSertifikat})
	}
	return out
}

// BolehTutupKlaim menjawab padanan prasyarat `ProtectLife.CARI1==""` (b992).
//
// ⚠️ Klaim TANPA peserta boleh ditutup, dan itu ditiru apa adanya: loopnya
// tidak pernah berjalan, penandanya tetap kosong, dan `Call FinishAssignment`
// berprasyarat kosong. Kita tidak menambahkan larangan yang XML tidak punya.
func BolehTutupKlaim(baris []BarisTutup) bool {
	return len(PenghalangTutupKlaim(baris)) == 0
}
