package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nusantarare/internal/models"
)

// Kontrak penanda DIPILIH - satu nilai, dua sisi.
//
// ⛔ Sebab uji ini ada, dan itu cacat yang sungguh terjadi: jalur PENDAFTARAN
// menulis `IsCheck = "1"` (`repository/pesertapolis.go`), sedangkan gerbang
// AKSEPTASI menuntut `"true"` (`services/akseptasi.go`). Akibatnya setiap
// peserta yang baru didaftarkan DITOLAK saat hendak diaksep, dengan pesan
// yang menyebut nilainya - dan seluruh jalur akseptasi tidak dapat dijalankan
// untuk klaim baru.
//
// Cacatnya tidak berbunyi di satu sisi mana pun: pendaftaran benar menurut
// komentarnya sendiri, akseptasi benar menurut XML, dan hanya PERTEMUANNYA
// yang salah. Karena itu yang diuji di sini adalah pertemuannya.
//
// XML menjawab tegas nilainya:
//
//	`Activity/SetIndexAdjustmentList.xml` b328  `.IsCheck` = `true`
//	`Activity/SavePesertaClaim.xml` b3631, b3919  `.IsCheck=="true"`
//	                                (langkah 7.7 dan 7.8 - keduanya HIDUP)
//
// RALAT GILIRAN-11: dulu dikutip b1812, b1997, b2413 - WHEN langkah 7.1
// (precondition false, mati) serta 7.2 dan 7.4 (`//`, ter-remark). Nilainya
// sama; hanya buktinya yang dipindah ke baris yang benar-benar berjalan.

func TestPenandaDipilihSatuNilaiSaja(t *testing.T) {
	// Peserta sebagaimana DIBUAT jalur pendaftaran harus lolos gerbang
	// akseptasi. Bila nilainya bergeser di salah satu sisi, uji ini merah.
	p := models.Peserta{ID: "P-1", IsCheck: models.PenandaDipilih}
	if !PesertaDipilih(p) {
		t.Fatalf("peserta ber-IsCheck %q ditolak gerbang akseptasi", p.IsCheck)
	}
	for _, bukan := range []string{"", "1", "0", "false", "ya"} {
		if PesertaDipilih(models.Peserta{IsCheck: bukan}) {
			t.Errorf("IsCheck %q diterima sebagai dipilih; hanya %q yang berarti dipilih",
				bukan, models.PenandaDipilih)
		}
	}
}

func TestNilaiPenandaDipilihDariXML(t *testing.T) {
	// ⛔ Nilainya VERBATIM dari rule, bukan pilihan gaya. Pega menuliskan
	// `true` dan mengujinya `=="true"`; menyimpan "1" atau "Y" berarti data
	// kita tidak dapat dibaca aturan yang sama.
	if models.PenandaDipilih != "true" {
		// Sumbernya disebut di komentar kepala uji ini, bukan di dalam
		// pesan: nama rule di string yang dapat mendarat di log dituduh
		// penjaga rujukan Komite, dan penjaga itu benar untuk berhati-hati.
		t.Errorf("PenandaDipilih = %q, mau %q", models.PenandaDipilih, "true")
	}
}

// Tidak ada satu pun berkas sumber yang menulis penanda itu dengan nilai lain.
//
// ⛔ Penjaga statik, sebab yang dijaga bukan hasil satu jalur melainkan
// KESERAGAMAN seluruh penulis. Jalur kedua yang menulis "1" tidak akan
// membuat uji perilaku mana pun merah sampai seseorang mencoba mengaksep
// peserta yang lahir dari jalur itu.
func TestNolPenulisPenandaDipilihBernilaiLain(t *testing.T) {
	diperiksa := 0
	err := filepath.Walk(filepath.Join("..", ".."), func(jalur string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			switch info.Name() {
			case "frontend", "node_modules", "bin", ".git", "dist", "unggahan":
				return filepath.SkipDir
			}
		}
		if err != nil || info.IsDir() || !strings.HasSuffix(jalur, ".go") ||
			strings.HasSuffix(jalur, "_test.go") {
			return err
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		diperiksa++
		for n, baris := range strings.Split(string(isi), "\n") {
			b := strings.TrimSpace(baris)
			if strings.HasPrefix(b, "//") || !strings.Contains(b, "IsCheck") {
				continue
			}
			// Penulisan: `IsCheck:` atau `IsCheck =`, bukan `IsCheck ==`.
			tulis := strings.Contains(b, "IsCheck:") ||
				(strings.Contains(b, "IsCheck =") && !strings.Contains(b, "IsCheck =="))
			if !tulis {
				continue
			}
			// Nilai yang sah: konstanta itu sendiri, atau penerusan nilai
			// yang sudah ada (mis. dari kolom basis data).
			if strings.Contains(b, "PenandaDipilih") || strings.Contains(b, "p.IsCheck") ||
				strings.Contains(b, "teks(") || strings.Contains(b, "IsCheck string") {
				continue
			}
			t.Errorf("%s:%d menulis penanda dipilih dengan nilai harfiah: %s\n"+
				"\tpakai models.PenandaDipilih - dua nilai berarti dua sisi yang "+
				"tidak pernah bertemu", filepath.ToSlash(jalur), n+1, b)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 30 {
		t.Fatalf("hanya %d berkas terbaca; penelusurnya yang rusak", diperiksa)
	}
}
