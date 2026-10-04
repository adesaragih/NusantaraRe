package models

// Copy Old - permintaan work owner 03-10-2026: tombol di samping `Add` membuka popup berisi SEMUA produk tabel JSON
// warisan (`M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE`) yang belum ada di tabel flat; yang dicentang disalin ke tabel
// flat lewat `Process Copy`. Tabel JSON hanya DIBACA.

// ProdukLama - satu baris popup `Copy Old`. Kolomnya sama dengan grid produk. `Alasan` (kenapa tidak dapat disalin)
// dan `Catatan` (yang berubah saat disalin) menyebut tabel dan kolom, tidak pernah nilainya.
type ProdukLama struct {
	ID           string   `json:"id"`
	ProductName  string   `json:"productName"`
	Ceding       string   `json:"ceding"`
	TreatyNumber string   `json:"treatyNumber"`
	InwardName   string   `json:"inwardName"`
	CreateOp     string   `json:"createOp"`
	UpdateOp     string   `json:"updateOp"`
	BolehDisalin bool     `json:"bolehDisalin"`
	Alasan       []string `json:"alasan"`
	Catatan      []string `json:"catatan"`
}

// Status satu produk sesudah `Process Copy`.
const (
	// SalinDisalin - induk + ketujuh anak tertulis di tabel flat dan terbaca ulang sama.
	SalinDisalin = "disalin"
	// SalinSudahAda - ID sudah ada di tabel flat (disalin sebelumnya, oleh pemakai lain, atau alat pindah).
	SalinSudahAda = "sudahAda"
	// SalinDitolak - rekonsiliasi menolak produk ini (`Pesan` = alasannya); tidak ada yang ditulis.
	SalinDitolak = "ditolak"
	// SalinGagal - galat basis data; transaksi produk ini dibatalkan, produk lain tidak terpengaruh.
	SalinGagal = "gagal"
)

// HasilSalinLama - hasil satu ID yang diminta.
type HasilSalinLama struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Pesan  []string `json:"pesan"`
}

// JawabanSalinLama - jawaban `Process Copy`: hasil per ID, urutan permintaan.
type JawabanSalinLama struct {
	Hasil   []HasilSalinLama `json:"hasil"`
	Disalin int              `json:"disalin"`
}
