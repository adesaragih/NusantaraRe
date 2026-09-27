package handlers

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
	"strings"
	"testing"
)

// suntikan adalah satu pasangan: penyusun layanan, dan yang wajib menyertainya.
type suntikan struct {
	// penyusun adalah pemanggilan yang mengembalikan layanan ber-stub.
	penyusun string
	// wajib adalah penyuntikan yang harus muncul di fungsi yang sama.
	wajib []string
}

// petaSuntikan memasangkan tiap berkas handler dengan kewajibannya.
//
// ⚠️ Didaftar per berkas, bukan dicari otomatis. Daftar yang dirakit sendiri
// oleh test dari kode yang diujinya akan selalu cocok dengan kode itu - dan
// tidak menjaga apa pun.
var petaSuntikan = map[string][]suntikan{
	"akseptasi.go": {{
		penyusun: "svc.Akseptasi()",
		wajib: []string{
			"DenganJejak(services.PerekamJejakOracle(svc))",
			"DenganPenerbit(services.PenerbitAkseptasiOracle(svc))",
		},
	}},
	"komite.go": {{
		penyusun: "svc.Komite()",
		wajib: []string{
			"DenganJejak(services.PerekamJejakOracle(svc))",
			"DenganRoster(services.RosterKomiteOracle(svc))",
			"DenganKasus(services.KasusKomiteOracle(svc))",
			"DenganPenyalur(services.PenyalurClaimLifeOracle(svc))",
		},
	}},
	"putaran.go": {{
		penyusun: "svc.Putaran()",
		wajib:    []string{"DenganJejak(services.PerekamJejakOracle(svc))"},
	}},
	"register.go": {{
		penyusun: "svc.Pendaftaran()",
		wajib:    []string{"DenganPenomor(services.PenomorCounterOracle(svc))"},
	}},
	"tolak.go": {{
		penyusun: "svc.Status()",
		wajib:    []string{"DenganJejak(services.PerekamJejakOracle(svc))"},
	}},
}

func TestHandlerMenyuntikkanImplementasiNyata(t *testing.T) {
	diperiksa := 0
	for berkas, daftar := range petaSuntikan {
		isi, err := os.ReadFile(berkas)
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
	t.Logf("%d penyuntikan wajib diperiksa di %d berkas", diperiksa, len(petaSuntikan))
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
	for berkas := range petaSuntikan {
		isi, err := os.ReadFile(berkas)
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
