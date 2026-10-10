package models

// Baris daftar kontrak WARISAN — dibaca dari `POOLDATA.TREATY_IN`.
//
// ⛔ ENTITAS YANG BERBEDA dari `BarisDaftarKontrak`, dan itu bukan duplikasi.
// Keduanya menjawab pertanyaan yang berbeda:
//
//	BarisDaftarKontrak   KONTRAK + VERSI_KONTRAK   kontrak yang sudah DIPINDAHKAN
//	BarisDaftarWarisan   TREATY_IN                 kontrak yang ADA DI SISTEM LAMA
//
// Hari ini yang pertama nol baris dan yang kedua 1.854. Ketika tiket 59
// berjalan keduanya akan berisi, dan yang membacanya harus tahu mana yang
// mana - menyatukan keduanya ke dalam satu bentuk menghapus perbedaan itu
// tepat ketika ia mulai penting.
//
// ⛔ `ID` bertipe TEKS, bukan bilangan. Kolom `TREATY_IN.ID` adalah
// `VARCHAR2(100)`. Sapuan 3 Oktober 2026 menemukan ke-1.854 nilainya memang
// berupa angka tujuh digit (1000001-1001856, nol ganda, nol bukan-angka) -
// tetapi KEBETULAN BUKAN JAMINAN: kolomnya teks, dan baris ke-1.855 tidak
// terikat apa pun.

// BarisDaftarWarisan adalah satu baris layar daftar, dari tabel warisan.
//
// ⚠️ Tiap medan terjemahan membawa nilai ASLINYA di sebelahnya. Itu bukan
// kemewahan: layar menampilkan bentuk yang dibaca orang, dan yang menyelidiki
// selisih pemindahan butuh bentuk yang tersimpan. Membuang salah satunya
// membuat salah satu pertanyaan itu tidak terjawab.
type BarisDaftarWarisan struct {
	// `TREATY_IN.ID` apa adanya.
	ID string `json:"id"`
	// `TREATYCONTRACTNAME`.
	NamaKontrak string `json:"namaKontrak"`

	// `PROPORTIONTYPE` apa adanya - `Proportional` atau `NonProportional`.
	SifatProporsiAsli string `json:"sifatProporsiAsli"`
	// Bentuk layar: `NonProportional` -> `Non Proportional`.
	SifatProporsi string `json:"sifatProporsi"`

	// `LEADINGREINSSOURCE` - nama broker, bukan pengenal.
	AsalBisnis string `json:"asalBisnis"`
	// `CEDING` - nama cedant.
	Cedant string `json:"cedant"`

	// `COMMENCEMENT`/`TERMINATION` apa adanya - `YYYYMMDD`.
	TanggalMulaiAsli    string `json:"tanggalMulaiAsli"`
	TanggalBerakhirAsli string `json:"tanggalBerakhirAsli"`
	// Bentuk layar `dd/mm/yy`; sama dengan aslinya bila tidak berbentuk
	// delapan angka - nol tanggal dikarang.
	TanggalMulai    string `json:"tanggalMulai"`
	TanggalBerakhir string `json:"tanggalBerakhir"`

	// `POSITIONUSERNAME`. ⚠️ 1.822 dari 1.854 baris NULL, bukan teks kosong.
	PosisiKe string `json:"posisiKe"`
	// `POSITION` - workbasket tempat berkas menunggu; syarat tampil tombol
	// `Revision` (`.Position = ''`).
	Posisi string `json:"posisi"`
	// `STATUSAKSEPTASI`. Empat nilai di data nyata: `Resolve Complete`
	// (1.820), `Accept` (12), `Decline` (11), dan NULL (11).
	StatusAkseptasi string `json:"statusAkseptasi"`
}

// HalamanDaftarWarisan adalah satu halaman beserta cacah seluruhnya.
//
// ⛔ `Total` dibawa supaya penomoran halaman dapat menghitung halaman
// terakhir tanpa menarik seluruh baris. Dengan 1.854 baris, menarik semuanya
// ke peramban untuk menghitung sepuluh tombol adalah ongkos yang tidak perlu
// dibayar siapa pun.
type HalamanDaftarWarisan struct {
	Baris   []BarisDaftarWarisan `json:"baris"`
	Halaman int                  `json:"halaman"`
	Ukuran  int                  `json:"ukuran"`
	Total   int                  `json:"total"`
}
