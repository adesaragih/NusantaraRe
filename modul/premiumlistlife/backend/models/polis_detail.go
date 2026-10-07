package models

// Kolom grid Premium List Detail - tiket 03 PremiumList Life.
//
// Untuk apa berkas ini: SATU daftar kolom, dari mana SQL, JSON, dan layar
// sama-sama diturunkan. Tiga daftar yang diketik terpisah akan berselisih,
// dan yang berselisih di sini adalah kolom uang yang tampil di bawah judul
// kolom lain.
//
// Sumbernya `Section/PL_Detail_Sec.xml`, dibaca 28-09-2026 dalam URUTAN
// DOKUMEN - bukan urutan abjad.
//
// `[terverifikasi]` 40 medan `<pyValue>` ber-titik yang UNIK tampil di sana;
// 38 punya kolom di migrasi 052, dan 2 tidak.
//
// Jendela sensus: SATU berkas, `PremiumList Life/Section/PL_Detail_Sec.xml`,
// hanya elemen `<pyValue>` yang isinya diawali titik (rujukan properti).
// Dihitung DUA cara, dan keduanya menjawab 40:
//
//	cara A - cacah nilai UNIK:
//	  grep -o "<pyValue>\.[A-Za-z0-9_.]*</pyValue>" Section/PL_Detail_Sec.xml \
//	    | sed 's/<[^>]*>//g' | sort -u | wc -l
//	cara B - urai urut dokumen, buang duplikat tanpa mengurutkan:
//	  grep -o "<pyValue>\.[A-Za-z0-9_.]*</pyValue>" Section/PL_Detail_Sec.xml \
//	    | sed 's/<[^>]*>//g' | awk '!s[$0]++' | wc -l
//
// ⚠️ JEBAKANNYA DISEBUT: kemunculan MENTAHnya 109, bukan 40 - satu medan
// dapat muncul berkali-kali (mode baca, mode sunting, kolom tersembunyi).
// Sensus yang mencacah kemunculan alih-alih nilai unik akan menjawab 109, dan
// menyimpulkan gridnya hampir tiga kali lebih lebar daripada yang ada.
//
// Dua yang tidak punya kolom:
//
//	REINSTYPENAME    nol kolom di migrasi 050-056 mana pun, dan nol rule lain
//	                 di korpus PremiumList Life yang menyebutnya - hanya
//	                 section ini
//	RetrocadedShare  BUKAN salah ketik dari RETROCEDED_SHARE: KEDUANYA tampil
//	                 di grid yang sama, jadi keduanya medan yang berbeda.
//	                 Mana yang mana `[terbuka - pemilik kerja]`
//
// ⛔ KEDUANYA TIDAK DITAMPILKAN, dan ketiadaannya DINYATAKAN - sama dengan
// `KetentuanUnderwriting` di kotak masuk. Sel kosong di layar terbaca "memang
// kosong", bukan "kami tidak punya datanya", dan kolom yang ditebak isinya
// lebih buruk lagi.
//
// ⛔ NAME_OF_INSURED TIDAK ADA di `KolomGridPeserta` (tiruan PL_Detail_Sec):
// Pega tidak menampilkannya di grid ini. Sejak 02-10-2026 ia tampil lewat
// `KolomGridTambahan` atas keputusan work owner. Ia nama orang: tidak pernah
// masuk fixture, tiket, atau log.
//
// ⛔ SELURUH KOLOM UANG KELUAR SEBAGAI TEKS (ADR-U-0003, ADR-U-0016). Tidak
// satu pun melewati `float64`, termasuk saat hanya ditampilkan - pembulatan
// yang terjadi di jalan pulang tidak kalah salah dari yang terjadi saat
// menyimpan.
//
// Dibaca sesudah: polis_nomor.go.

// JenisKolomPeserta menyatakan cara satu kolom dibaca dari Oracle.
//
// ⚠️ Ini BUKAN tipe fisiknya, melainkan cara membacanya. Uang dan tanggal
// keduanya keluar sebagai teks, tetapi lewat pembungkus yang berbeda -
// `TM9` untuk yang pertama, pola tanggal untuk yang kedua.
type JenisKolomPeserta int

// ⚠️ Ketiganya SENGAJA sejajar dengan `golonganKolom` di
// `repository/kolompeserta.go` - teks, tanggal, angka. Golongan keempat untuk
// "bilangan bulat" sempat ada di sini dan DIBUANG: ia dibaca dengan
// pembungkus yang persis sama dengan angka, jadi ia perbedaan yang tidak
// membedakan apa pun - dan perbedaan semacam itu membuat orang mengira ada
// perlakuan khusus yang harus dijaga.
const (
	// KolomPesertaTeks dibaca apa adanya.
	KolomPesertaTeks JenisKolomPeserta = iota
	// KolomPesertaAngka dibungkus TO_CHAR ber-TM9 - uang, share, rate,
	// faktor, DAN umur serta periode.
	KolomPesertaAngka
	// KolomPesertaTanggal dibungkus TO_CHAR berpola tanggal.
	KolomPesertaTanggal
)

// KolomPeserta adalah satu kolom grid.
type KolomPeserta struct {
	// Nama adalah nama kolom di migrasi 052, HURUF BESAR.
	Nama string
	// Jenis menentukan pembungkusnya saat dibaca.
	Jenis JenisKolomPeserta
}

// KolomGridPeserta adalah kolom grid Premium List Detail, urut seperti
// `PL_Detail_Sec.xml` menampilkannya.
//
// ⚠️ URUTANNYA BAGIAN DARI TIRUAN. Orang yang sudah bertahun-tahun memakai
// layar lama membaca grid ini dengan mata, bukan dengan judul kolom; menyusun
// ulang kolomnya "supaya lebih rapi" membuat mereka salah baca baris yang
// benar.
var KolomGridPeserta = []KolomPeserta{
	{"BEGIN_DATE", KolomPesertaTanggal},
	{"ENTRY_AGE", KolomPesertaAngka},
	{"PLAN", KolomPesertaTeks},
	{"STNC", KolomPesertaTeks},
	{"EXPIRED_DATE", KolomPesertaTanggal},
	{"CURRENT_AGE", KolomPesertaAngka},
	{"PERIOD_MM", KolomPesertaAngka},
	{"POLICY_HOLDER", KolomPesertaTeks},
	{"EFFECTIVE_DATE", KolomPesertaTanggal},
	{"WPC", KolomPesertaTeks},
	{"CERTIFICATE_NO", KolomPesertaTeks},
	{"CURRENCY", KolomPesertaTeks},
	{"SUM_REASURED", KolomPesertaAngka},
	{"DEDUCTION", KolomPesertaAngka},
	{"BROKERAGE_FEE", KolomPesertaAngka},
	{"SUM_AT_RISK_RETRO", KolomPesertaAngka},
	{"RETROCEDED_SHARE", KolomPesertaAngka},
	{"SUM_INSURED", KolomPesertaAngka},
	{"SHARE_NUSANTARA_RE_GROSS", KolomPesertaAngka},
	{"EM_PERCENT", KolomPesertaAngka},
	{"SUM_AT_RISK_GROSS", KolomPesertaAngka},
	{"GROSS_PREMIUM", KolomPesertaAngka},
	{"CEDING_RETENTION", KolomPesertaAngka},
	{"RATE", KolomPesertaAngka},
	// ⛔ FACTOR adalah DESIMAL, bukan bilangan bulat (AC 39 spec). Kolomnya
	// `NUMBER(38,8)`, dan membacanya sebagai bilangan bulat membuang tujuh
	// angka di belakang koma - tepat angka yang membedakan dua faktor.
	{"FACTOR", KolomPesertaAngka},
	{"RISK", KolomPesertaTeks},
	{"NET_PREMIUM", KolomPesertaAngka},
	{"DEDUCTION_REFUND", KolomPesertaAngka},
	{"BROKERAGE_FEE_REFUND", KolomPesertaAngka},
	{"GROSS_PREMIUM_REFUND", KolomPesertaAngka},
	{"NET_PREMIUM_REFUND", KolomPesertaAngka},
	{"RI_ADMIN_FEE_REFUND_RETRO", KolomPesertaAngka},
	{"SHARE_RETRO", KolomPesertaAngka},
	{"GROSS_PREMIUM_REFUND_RETRO", KolomPesertaAngka},
	{"NET_PREMIUM_REFUND_RETRO", KolomPesertaAngka},
	{"RI_ADMIN_FEE_RETRO", KolomPesertaAngka},
	{"GROSS_PREMIUM_RETRO", KolomPesertaAngka},
	{"NET_PREMIUM_RETRO", KolomPesertaAngka},
}

// MedanGridTanpaKolom adalah medan `PL_Detail_Sec` yang TIDAK punya kolom.
//
// ⛔ Didaftar, bukan dilupakan. Daftar ini yang membuat selisih antara empat
// puluh medan di layar lama dan tiga puluh delapan kolom di grid kami dapat
// dijawab - oleh siapa pun yang kelak menghitungnya dan bertanya.
var MedanGridTanpaKolom = map[string]string{
	"REINSTYPENAME": "nol kolom di migrasi 050-056, dan nol rule lain di " +
		"korpus PremiumList Life yang menyebutnya",
	"RetrocadedShare": "berdampingan dengan RETROCEDED_SHARE di grid yang " +
		"sama, jadi medan yang berbeda; mana yang mana [terbuka]",
}

// BarisPeserta adalah satu baris grid.
//
// ⚠️ Nilainya PETA, bukan empat puluh medan bernama. Gridnya memang
// digerakkan daftar kolom - satu daftar yang merakit SQL-nya, kunci petanya,
// dan judul kolomnya sekaligus. Empat puluh medan bernama berarti empat
// puluh tempat yang harus ikut berubah setiap kali satu kolom bergeser, dan
// yang tertinggal satu di antaranya tidak akan ketahuan.
type BarisPeserta struct {
	// ID adalah kunci baris - dipakai layar sebagai `key`, bukan ditampilkan.
	ID string `json:"id"`
	// Nilai berkunci nama kolom, seluruhnya TEKS.
	Nilai map[string]string `json:"nilai"`
}

// NamaKolomGridPeserta mengembalikan nama kolom grid, urut.
func NamaKolomGridPeserta() []string {
	nama := make([]string, 0, len(KolomGridPeserta))
	for _, k := range KolomGridPeserta {
		nama = append(nama, k.Nama)
	}
	return nama
}

// KolomGridTambahan adalah kolom grid di LUAR `PL_Detail_Sec` - ditambahkan
// atas keputusan work owner 02-10-2026 ("tampilkan NAME_OF_INSURED, DOB,
// GROSS_VALUATION_BEGIN_DATE, GROSS_VALUATION_EXPIRED_DATE").
//
// ⛔ DAFTAR TERPISAH, bukan disisipkan ke `KolomGridPeserta`: daftar itu
// tiruan `PL_Detail_Sec` (cacah dan urutannya dijaga uji paritas), dan
// mencampurnya membuat paritas itu tidak lagi dapat diperiksa.
//
// ⚠️ NAME_OF_INSURED adalah nama orang. Ia tampil di layar karena keputusan
// work owner, tetapi aturan lainnya TETAP: tidak pernah masuk fixture,
// tiket, atau log.
var KolomGridTambahan = []KolomPeserta{
	// POLICY_NO dari CSV, kolom paling depan grid (permintaan work owner 05-10-2026).
	{"POLICY_NO", KolomPesertaTeks},
	{"NAME_OF_INSURED", KolomPesertaTeks},
	{"DOB", KolomPesertaTanggal},
	{"GROSS_VALUATION_BEGIN_DATE", KolomPesertaTanggal},
	{"GROSS_VALUATION_EXPIRED_DATE", KolomPesertaTanggal},
	// RI Admin Fee - ikut disimpan unggahan sejak 03-10-2026 dan dipakai rekap.
	{"RI_ADMIN_FEE", KolomPesertaAngka},
}

// KolomGridTampil mengembalikan kolom yang benar-benar dibaca dan dikirim ke
// layar: `KolomGridPeserta` lalu `KolomGridTambahan`. Urutan tampilnya
// disusun layar (`susunKolom`).
func KolomGridTampil() []KolomPeserta {
	k := make([]KolomPeserta, 0, len(KolomGridPeserta)+len(KolomGridTambahan))
	k = append(k, KolomGridPeserta...)
	return append(k, KolomGridTambahan...)
}

// NamaKolomGridTampil - nama kolom `KolomGridTampil`, urut.
func NamaKolomGridTampil() []string {
	k := KolomGridTampil()
	nama := make([]string, 0, len(k))
	for _, x := range k {
		nama = append(nama, x.Nama)
	}
	return nama
}

// NamaKolomGridAngka - nama kolom `KolomGridTampil` yang ber-jenis ANGKA,
// urut. Layar memberinya pemisah ribuan (03-10-2026).
func NamaKolomGridAngka() []string {
	var nama []string
	for _, k := range KolomGridTampil() {
		if k.Jenis == KolomPesertaAngka {
			nama = append(nama, k.Nama)
		}
	}
	return nama
}
