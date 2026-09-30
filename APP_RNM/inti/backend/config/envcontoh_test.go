package config

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Setiap env yang DIBACA kode harus ada di .env.example.
//
// ⛔ Kenapa test ini ada, dengan sebab yang sungguh terjadi: `AUTH_STUB`
// dibaca config sejak awal tetapi tidak pernah tertulis di berkas contoh.
// Akibatnya backend berjalan dengan gerbang identitas MATI, mengabaikan
// header yang dikirim React, dan menjawab 401 - sehingga Inbox tampak KOSONG
// di layar. Kosong terbaca "tidak ada pekerjaan", bukan "permintaan ditolak",
// jadi tidak ada yang mencurigai konfigurasi.
//
// Knob yang ada di kode tetapi tidak ada di berkas contoh adalah knob yang
// tidak diketahui siapa pun - dan yang tidak diketahui tidak akan disetel.
func TestSetiapEnvTerdokumentasi(t *testing.T) {
	akar := filepath.Join("..", "..", "..") // APP_RNM, dari inti/backend/config
	contoh, err := os.ReadFile(filepath.Join(akar, ".env.example"))
	if err != nil {
		t.Fatalf("tidak dapat membaca .env.example: %v", err)
	}
	teks := string(contoh)

	pola := regexp.MustCompile(`os\.Getenv\("([A-Z_][A-Z0-9_]*)"\)`)
	terbaca := map[string]string{}
	err = filepath.Walk(akar, func(jalur string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(jalur, ".go") {
			return nil
		}
		if strings.HasSuffix(jalur, "_test.go") {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return nil
		}
		for _, m := range pola.FindAllStringSubmatch(string(isi), -1) {
			terbaca[m[1]] = jalur
		}
		return nil
	})
	if err != nil {
		t.Fatalf("menelusuri sumber: %v", err)
	}
	if len(terbaca) == 0 {
		t.Fatal("nol os.Getenv terbaca; penelusurnya yang rusak, bukan kodenya")
	}

	for nama, jalur := range terbaca {
		// Dicari sebagai KUNCI di awal baris, bukan sekadar disebut:
		// menyebutkannya di dalam kalimat komentar bukan mendokumentasikan
		// cara menyetelnya.
		if !strings.Contains(teks, "\n"+nama+"=") && !strings.HasPrefix(teks, nama+"=") {
			t.Errorf("%s dibaca di %s tetapi tidak ada di .env.example sebagai %s=; "+
				"knob yang tidak tertulis tidak akan disetel siapa pun", nama, jalur, nama)
		}
	}
	t.Logf("%d env diperiksa", len(terbaca))
}
