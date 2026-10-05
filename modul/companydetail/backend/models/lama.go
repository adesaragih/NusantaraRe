package models

// Copy Old - perintah work owner 04-10-2026: "buatkan fungsinya copy old seperti pada productname; copy old hanya muncul
// untuk superadmin". Tombol di samping `Create` membuka popup berisi organisasi dokumen Pega (`M_CLIENT`) yang datanya
// BELUM pindah ke tabel datar; yang dicentang disalin lewat `Process Copy` dengan aturan alat pindah. Dokumen hanya
// DIBACA. Superadmin = akun pemegang menu Kelola User (admin aplikasi).

// OrgLama - satu baris popup `Copy Old`: apa yang AKAN ditulis untuk organisasi ini. Tanpa nilai PIC atau nomor.
type OrgLama struct {
	ID     string `json:"id"`
	IDView string `json:"idView"`
	// Nama - nama organisasi di dokumen.
	Nama string `json:"nama"`
	// Baru - organisasi belum punya baris CLIENT; Copy Old membuatnya.
	Baru bool `json:"baru"`
	// PIC - PIC dokumen yang akan ditambah ke CLIENT_PICLIST.
	PIC int `json:"pic"`
	// Nomor - nomor Phone and Fax dokumen yang akan masuk CLIENT_ADDRESS.
	Nomor int `json:"nomor"`
	// Isi - kolom CLIENT kosong yang akan diisi (`PARENT_ID`, `NOTE`, `TITLE`).
	Isi []string `json:"isi"`
}

// Status satu organisasi sesudah `Process Copy`.
const (
	// SalinDisalin - ada yang ditulis ke tabel datar.
	SalinDisalin = "disalin"
	// SalinSudahAda - tidak ada lagi yang perlu ditulis (disalin sebelumnya, oleh pemakai lain, atau alat pindah).
	SalinSudahAda = "sudahAda"
	// SalinDitolak - ID bukan organisasi dokumen Pega.
	SalinDitolak = "ditolak"
	// SalinGagal - galat basis data; transaksi organisasi ini dibatalkan, yang lain tidak terpengaruh.
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
