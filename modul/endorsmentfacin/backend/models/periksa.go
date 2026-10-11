package models

// Kabel JSON `POST /api/endorsmentfacin/periksa` - tombol periksa sebelum Create (irisan I2, kontrak sesi utama
// 11-10-2026). Nol tulis. Di Pega dua pemeriksaan TERPISAH yang dipicu layar `Section/WorkPrimaryDetails.xml`:
// `SetErrorBatalEndorsement_Act` (PolicyNo / EdmType / EdmTypeNew / Internal Retro berubah) dan `CheckEDMPolisDate`
// (EndorsementDate berubah); tombol OK aktif bila keduanya `CARIGROUP=="true"` (`Section/crmNewHarnessButtons.xml`).
// Di sini satu panggilan membawa seluruh isian.

// MasukanPeriksa - badan permintaan. Tanggal teks `DD-MM-YYYY`; kosong = pemeriksaan tanggal belum jalan.
type MasukanPeriksa struct {
	PolicyNo                 string `json:"policyNo"`
	EdmType                  string `json:"edmType"`
	EdmTypeNew               string `json:"edmTypeNew"`
	EndorsementDate          string `json:"endorsementDate"`
	EndorsementInternalRetro string `json:"endorsementInternalRetro"`
}

// HasilGerbangPeriksa - nilai `hasil` satu gerbang.
type HasilGerbangPeriksa string

const (
	GerbangLolos    HasilGerbangPeriksa = "lolos"
	GerbangTolak    HasilGerbangPeriksa = "tolak"
	GerbangDiklep   HasilGerbangPeriksa = "diklep"
	GerbangDilewati HasilGerbangPeriksa = "dilewati"
	// GerbangBelumDiperiksa - data pemeriksanya tidak tersedia (rule tanpa body di korpus, atau skema lain tidak
	// terbaca). BUKAN lolos.
	GerbangBelumDiperiksa HasilGerbangPeriksa = "belum-diperiksa"
)

// PesanPeriksa - satu `Property-Set-Messages`, teks verbatim korpus.
type PesanPeriksa struct {
	Medan string `json:"medan"`
	Teks  string `json:"teks"`
}

// GerbangPeriksa - satu dari enam gerbang `SetErrorBatalEndorsement_Act` langkah 5.
type GerbangPeriksa struct {
	No      int                 `json:"no"`
	Hasil   HasilGerbangPeriksa `json:"hasil"`
	Catatan string              `json:"catatan"`
}

// BisnisPeriksa - langkah 7-18: jenis bisnis polis lama dan predikat lini.
type BisnisPeriksa struct {
	Hasil HasilGerbangPeriksa `json:"hasil"`
	// JenisBisnis - SESUDAH langkah 8 (polis berpola RNML = "Life").
	JenisBisnis string `json:"jenisBisnis"`
}

// TanggalPeriksa - `CheckEDMPolisDate`; Mulai/Akhir = tanggal kalender periode yang dibandingkan (`DD-MM-YYYY`,
// kosong bila tidak terbaca).
type TanggalPeriksa struct {
	Hasil HasilGerbangPeriksa `json:"hasil"`
	Mulai string              `json:"mulai"`
	Akhir string              `json:"akhir"`
	// Catatan - mis. tabel pengecualian belum dibuat DBA (daftar dianggap kosong). Kosong = tidak dikirim.
	Catatan string `json:"catatan,omitempty"`
}

// HasilPeriksa - jawaban 200.
type HasilPeriksa struct {
	// Ditolak - ada pesan dari gerbang, bisnis, atau tanggal.
	Ditolak bool `json:"ditolak"`
	// Lengkap - false bila ada gerbang "belum-diperiksa".
	Lengkap bool             `json:"lengkap"`
	Pesan   []PesanPeriksa   `json:"pesan"`
	Gerbang []GerbangPeriksa `json:"gerbang"`
	Bisnis  BisnisPeriksa    `json:"bisnis"`
	// Tanggal - nil bila endorsementDate kosong.
	Tanggal *TanggalPeriksa `json:"tanggal"`
}
