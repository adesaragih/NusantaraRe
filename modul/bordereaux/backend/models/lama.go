package models

// Copy Old Data (keputusan work owner 04-10-2026: "buat copy old data", hanya superadmin): berkas Pega yang isinya
// masih tertinggal di JSON `M_BORDEREAUX.DATA_JSON` disalin ke tabel. JSON hanya DIBACA.

// JSONLama - satu baris M_BORDEREAUX beserta header BORDEREAUX-nya bila ada.
type JSONLama struct {
	ID  string
	Isi string
	// AdaHeader - baris BORDEREAUX ber-BDX_ID sama ada; Type dan Business lalu diambil dari header itu.
	AdaHeader      bool
	Type, Business string
}

// BerkasLama - satu baris popup Copy Old Data: berkas yang masih punya isi JSON yang belum ada di tabel.
type BerkasLama struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Business string `json:"business"`
	Ceding   string `json:"ceding"`
	Treaty   string `json:"treaty"`
	Status   string `json:"status"`
	// BarisJSON - baris detail kombinasi berkas di JSON; BarisTabel - baris di tabel detailnya.
	BarisJSON  int `json:"barisJson"`
	BarisTabel int `json:"barisTabel"`
	// Komentar - baris CommentList JSON; Riwayat - baris BORDEREAUX_HISTORY.
	Komentar int `json:"komentar"`
	Riwayat  int `json:"riwayat"`
	// TanpaHeader - baris BORDEREAUX belum ada; dibuat dari JSON.
	TanpaHeader bool `json:"tanpaHeader"`
}

// Status hasil salin satu berkas.
const (
	SalinDisalin  = "disalin"
	SalinSudahAda = "sudahAda"
	SalinDitolak  = "ditolak"
	SalinGagal    = "gagal"
)

// HasilSalinLama - hasil Process Copy satu berkas.
type HasilSalinLama struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Pesan  []string `json:"pesan"`
}

// JawabanSalinLama - hasil satu permintaan Process Copy.
type JawabanSalinLama struct {
	Hasil   []HasilSalinLama `json:"hasil"`
	Disalin int              `json:"disalin"`
}
