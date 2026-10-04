package repository

// Dua aturan pemilihan nilai peserta, ditiru dari `SaveInsuredClaim_Act`.
//
// Untuk apa berkas ini: layar Pega `Select Insured` memindahkan baris polis ke
// daftar klaim, dan saat memindahkannya ia MEMILIH nilai untuk dua medan -
// umur dan share Nusantara Re. Kedua pilihan itu aturan bisnis, bukan hiasan
// layar, sehingga rumahnya di sini, pada jalur yang membaca sumbernya.
//
// ⛔ Kenapa di repository dan BUKAN di layar: server membaca ulang baris polis
// saat POST justru supaya nilai uang tidak pernah datang dari klien
// (`AmbilUntukKlaim`). Menaruh aturan ini di React berarti angka uang
// ditentukan browser - persis yang ADR-U-0003 larang. Di Pega ia kebetulan
// berada di layar; yang ditiru MAKSUDnya.
//
// Sumber, dibaca sebagai pohon langkah utuh (bukan jendela maju):
//
//	`Activity/SaveInsuredClaim_Act.xml`
//	  langkah 2 `Property-Set` halaman `TempDetail1.pxResults`
//	    sublangkah 1 b508 "Set property PremiumListDetail"
//	      b601  SHARE_NUSANTARA_RE
//	      b661  AGE
//	    sublangkah 2 b1339 "Set property AdjustmentList"
//	      b1412 SHARE_NUSANTARA_RE  (ungkapan yang sama persis)
//
// Dibaca sesudah: pesertapolis.go.

// UmurPeserta memilih umur peserta dari tiga kolom sumber.
//
//	b661: @if(.AGE!="",.AGE,@if(.ENTRY_AGE!="",.ENTRY_AGE,.CURRENT_AGE))
//
// ⛔ Hanya teks KOSONG yang membuatnya turun ke kolom berikutnya. "0" adalah
// umur yang sah - bayi - dan XML tidak pernah membandingkannya dengan "0".
// Bandingkan dengan ShareNusantaraReTeks, yang justru memperlakukan "0"
// sebagai kosong: kedua aturan ini SENGAJA berbeda, dan ujinya mengunci
// perbedaan itu supaya tidak ada yang menyeragamkannya.
func UmurPeserta(age, entryAge, currentAge string) string {
	if age != "" {
		return age
	}
	if entryAge != "" {
		return entryAge
	}
	return currentAge
}

// ShareNusantaraReTeks memilih share Nusantara Re, dengan cadangan gross.
//
//	b601/b1412: @if(.SHARE_NUSANTARA_RE=="0", GROSS,
//	                @if(.SHARE_NUSANTARA_RE=="", GROSS, SHARE))
//
// ⛔ Di sini "0" MEMANG berarti kosong. Memo penghapusan rule itu sendiri
// berbunyi "fix share nusantara re": baris polis yang share-nya belum dibagi
// tersimpan sebagai nol, dan nilai yang benar ada di kolom gross-nya.
//
// ⚠️ Perbandingannya dengan teks "0" PERSIS, bukan "bernilai nol". Itu
// mengikuti XML. Pembaca kita memakai TO_CHAR(..,'TM9'), yang menuliskan nol
// sebagai "0" tepat - jadi jalur produksinya bertemu; "0.00" hanya mungkin
// datang dari sumber lain, dan untuk itu XML memilih share apa adanya.
func ShareNusantaraReTeks(share, gross string) string {
	if share == "0" || share == "" {
		return gross
	}
	return share
}
