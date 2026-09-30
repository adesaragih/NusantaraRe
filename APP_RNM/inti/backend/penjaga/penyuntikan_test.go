package penjaga

// Penjaga penyuntikan implementasi nyata - perapian A2.
//
// Untuk apa berkas ini: MENGGANTIKAN cabang 501 yang dibuang.
//
// ⛔ Cabang `…BelumDiputuskan` → 501 di kelima handler kini MATI, sebab tiap
// handler menyuntikkan implementasi Oracle-nya sebelum memanggil. Membuang
// cabang mati itu benar; membuangnya TANPA menggantikan apa pun tidak. Bila
// kelak seseorang menghapus satu baris `Dengan…`, stub bawaannya hidup lagi,
// layanan gagal dengan pesan "belum diputuskan work owner" yang menyesatkan,
// dan tidak ada satu pun test yang jatuh.
//
// Berkas ini test itu: bawaan layanan WAJIB gagal terang (dijaga di
// services/), dan handler WAJIB menggantinya di sini.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// petaSuntikan (`modulmd_test.go`) memasangkan tiap berkas handler dengan
// kewajibannya - dari bab "Penyuntikan wajib di handler" `MODUL.md` setiap
// modul.
//
// ⚠️ Didaftar per berkas, bukan dicari otomatis. Daftar yang dirakit sendiri
// oleh test dari kode yang diujinya akan selalu cocok dengan kode itu - dan
// tidak menjaga apa pun. Struktur tim satu folder per modul (30-09-2026):
// daftarnya tetap ditulis tangan, oleh MODUL pemilik handler-nya, bukan di sini.

// bacaHandler membaca satu berkas handler; `rel` jalur relatif akar aplikasi
// (`modul/<nama>/backend/handlers/<berkas>`).
func bacaHandler(rel string) ([]byte, error) {
	return os.ReadFile(filepath.Join(akarAplikasi, filepath.FromSlash(rel)))
}

func TestHandlerMenyuntikkanImplementasiNyata(t *testing.T) {
	diperiksa := 0
	peta := petaSuntikan(t)
	for berkas, daftar := range peta {
		isi, err := bacaHandler(berkas)
		if err != nil {
			t.Errorf("membaca %s: %v", berkas, err)
			continue
		}
		teks := string(isi)
		for _, s := range daftar {
			if !strings.Contains(teks, s.penyusun) {
				t.Errorf("%s: tidak lagi memanggil %s - peta penjaga ini usang, "+
					"dan penjaga usang tidak menjaga apa pun", berkas, s.penyusun)
				continue
			}
			for _, w := range s.wajib {
				diperiksa++
				if !strings.Contains(teks, w) {
					t.Errorf("%s: memanggil %s tetapi tidak menyuntikkan %s; "+
						"stub bawaannya akan hidup lagi dan gagal dengan pesan "+
						"\"belum diputuskan work owner\" yang menyesatkan",
						berkas, s.penyusun, w)
				}
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol penyuntikan diperiksa; pembacanya yang rusak")
	}
	t.Logf("%d penyuntikan wajib diperiksa di %d berkas", diperiksa, len(peta))
}

// Dan cabang 501-nya memang sudah tidak ada lagi.
//
// ⛔ Penjaga arah-balik: tanpa ini, seseorang dapat "memperbaiki" kegagalan
// penjaga di atas dengan menambahkan kembali cabang 501 alih-alih
// menyuntikkan implementasinya.
func TestNolCabang501StubDiHandler(t *testing.T) {
	stub := []string{
		"ErrPenomorBelumDiputuskan",
		"ErrJejakBelumDiputuskan",
		"ErrRosterBelumDiputuskan",
		"ErrKasusKomiteBelumDiputuskan",
		"ErrResolverBelumDiputuskan",
		"ErrAntreanBelumDiputuskan",
	}
	for berkas := range petaSuntikan(t) {
		isi, err := bacaHandler(berkas)
		if err != nil {
			t.Errorf("membaca %s: %v", berkas, err)
			continue
		}
		for _, s := range stub {
			if strings.Contains(string(isi), s) {
				t.Errorf("%s masih memetakan %s ke kode HTTP; stubnya sudah "+
					"diganti, jadi cabang itu mati - dan cabang mati "+
					"menyembunyikan bahwa penyuntikannya hilang", berkas, s)
			}
		}
	}
}
