package services

// Nomor akseptasi Komite - tiket 04a. TANPA Oracle.

import (
	"os"
	"strings"
	"testing"
)

func sumberKomiteAkseptasi(t *testing.T) string {
	t.Helper()
	isi, err := os.ReadFile("komite_akseptasi.go")
	if err != nil {
		t.Fatal(err)
	}
	var kode []string
	for _, b := range strings.Split(string(isi), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(b), "//") {
			kode = append(kode, b)
		}
	}
	return strings.Join(kode, "\n")
}

// TestPenghitungKomiteBukanSequenceClaimLife - AC 14 spec, keputusan o.
func TestPenghitungKomiteBukanSequenceClaimLife(t *testing.T) {
	kode := sumberKomiteAkseptasi(t)
	for _, wajib := range []string{"penghitung.AwalanProduksi(", "penghitung.UrutNomorBerikut(",
		"models.ClassPenghitungKomiteLife", "models.JenisPenghitungKomite(", "models.NomorAkseptasiKomite("} {
		if !strings.Contains(kode, wajib) {
			t.Errorf("jalur nomor Komite tidak memakai %s", wajib)
		}
	}
	for _, terlarang := range []string{"UrutAkseptasiBerikut(", "ACCEPTATIONNOLIFE_SEQ",
		"Generate_NoAccept_KMT", "RakitNomorAkseptasi(", `"RNML-`} {
		if strings.Contains(kode, terlarang) {
			t.Errorf("jalur nomor Komite memuat %q - penghitung/awalan milik jalur lain atau literal", terlarang)
		}
	}
}

// TestKeunikanDiperiksaSebelumStempel - OQ-K-04a, tabrakan gagal terang.
func TestKeunikanDiperiksaSebelumStempel(t *testing.T) {
	kode := sumberKomiteAkseptasi(t)
	iLama := strings.Index(kode, ".NomorAkseptasiDipakai(ctx, tx, nomor)")
	iAdj := strings.Index(kode, ".NomorAkseptasiDipakaiDiAdjustment(ctx, tx, nomor)")
	iStempel := strings.Index(kode, "baca.PerbaruiStatusBaris(")
	if iLama < 0 || iAdj < 0 || iStempel < 0 {
		t.Fatal("pemeriksa keunikan atau stempel hilang")
	}
	// Pemeriksa ada di terbitkanNomor, yang dipanggil SEBELUM stempel.
	iPanggil := strings.Index(kode, "p.terbitkanNomor(ctx, tx,")
	if iPanggil < 0 || iPanggil > iStempel {
		t.Error("nomor tidak diterbitkan (dan diperiksa) sebelum baris distempel")
	}
}

// TestRekamAkhirSatuJalurDalamAkseptasi - tiket 04b, AC 16 spec.
func TestRekamAkhirSatuJalurDalamAkseptasi(t *testing.T) {
	kode := sumberKomiteAkseptasi(t)
	iHeader := strings.Index(kode, "baca.CerminkanHeader(")
	iRekam := strings.Index(kode, "p.rekamAkhir(ctx, tx, kasus, models.KodeAksep, nomor, saat)")
	iJejak := strings.Index(kode, "p.jejak.Rekam(")
	if iHeader < 0 || iRekam < 0 || iJejak < 0 || !(iHeader < iRekam && iRekam < iJejak) {
		t.Error("urutan akseptasi harus: stempel → rekam akhir → jejak")
	}
	if n := strings.Count(kode, "RekamAkhirWarisan("); n != 1 {
		t.Errorf("rekam akhir dipanggil di %d tempat, mau 1 (satu jalur berparameter status)", n)
	}
}
