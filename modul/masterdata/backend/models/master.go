// Package models memuat bentuk data modul Master Data: daftar tabel master, kolomnya, dan rujukannya.
package models

// Kolom - satu kolom tabel master. Semua kolom master bertipe teks (VARCHAR2) - DDL `D:\migrasi\RNM\DDL\`.
type Kolom struct {
	// Nama - nama kolom Oracle apa adanya (dikutip bila kata cadangan, mis. `"GROUP"`).
	Nama string
	// JSON - kunci kabel (camelCase).
	JSON string
	// Lebar - batas byte kolom (DDL / migrasi 760).
	Lebar int
	// Wajib - harus terisi saat simpan.
	Wajib bool
	// Turunan - diisi backend dari tabel rujukannya (M-5), tidak diterima dari masukan.
	Turunan bool
	// HurufBesar - disimpan huruf besar (ACCUMULATION.NOTE: view asal `UPPER (b.JSONDATA.Note)`).
	HurufBesar bool
	// Tanggal - kolom DATE (jejak ubah): dibaca `TO_CHAR(..., 'YYYY-MM-DD HH24:MI:SS')`, ditulis SYSDATE.
	Tanggal bool
}

// Rujukan - kolom yang menunjuk baris tabel lain: nilainya wajib ada di `Tabel.KolomKunci` (bila diisi), dan bila
// `Turunan` terisi, kolom turunan itu diisi `KolomNilai` baris rujukan.
type Rujukan struct {
	Sumber     string // kolom Oracle di tabel ini
	Tabel      string // tabel rujukan
	KolomKunci string
	KolomNilai string // kosong = hanya pemeriksaan keberadaan
	Turunan    string // kolom Oracle tabel ini yang diisi
}

// TabelMaster - satu master yang dapat diubah lewat menu Master Data.
type TabelMaster struct {
	// Kunci - nama di rute `/api/masterdata/{tabel}`.
	Kunci string
	// Judul - nama tampil (= nama master di rencana).
	Judul string
	// Nama - tabel Oracle.
	Nama string
	// Kolom - kolom data; yang pertama = kunci baris `ID`.
	Kolom []Kolom
	// KolomCari - kolom yang dicocokkan `q` (Contains tidak peka huruf).
	KolomCari []string
	// KolomStatus - kolom aktif / nonaktif: '1' = aktif. Kosong = tabel warisan tanpa kolom status; statusnya di
	// T_MASTER_STATUS (tanpa baris = aktif).
	KolomStatus string
	// IDOtomatis - ID dibuat backend (ACCUMULATION, M-4); selainnya diisi pengguna dan unik.
	IDOtomatis bool
	// JejakTerpisah - tabel warisan yang tidak di-ALTER (MD-2): jejak ubah (dan status bila KolomStatus kosong) di
	// T_MASTER_STATUS; selainnya kolom jejak tabel itu sendiri (migrasi 762).
	JejakTerpisah bool
	Rujukan       []Rujukan
}

// Baris - satu baris master di kabel: kunci JSON kolom -> nilai teks.
type Baris map[string]string

// Halaman - hasil daftar satu master.
type Halaman struct {
	Baris  []Baris
	Aktif  []bool
	Total  int
	Nomor  int
	Ukuran int
}

// SeluruhKolom - kolom data master lalu keempat kolom jejak ubah (KolomAudit): urutan kabel GET dan pindai baris.
func (t TabelMaster) SeluruhKolom() []Kolom {
	return append(append([]Kolom{}, t.Kolom...), KolomAudit...)
}
