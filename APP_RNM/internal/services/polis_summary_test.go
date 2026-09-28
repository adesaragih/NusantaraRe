package services

// Rekap dan submit premium list - tiket 05a bagian 2. TANPA Oracle.
//
// ⛔ Alasan penjaga statik sama dengan polis_nomor_test.go: jalur ini hanya
// berjalan dengan Oracle, dan uji yang menuntut Oracle SKIP diam-diam.
// ⚠️ Bukti tingkat Oracle (injeksi kegagalan di peserta ke-N → nol baris
// tersimpan) BELUM ADA: skema uji belum memasang tabel 050-056. Dilaporkan,
// tidak diklaim.

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/repository"
)

func sumberSummary(t *testing.T) string {
	t.Helper()
	isi, err := os.ReadFile("polis_summary.go")
	if err != nil {
		t.Fatalf("membaca polis_summary.go: %v", err)
	}
	return string(isi)
}

// badanFungsi memotong satu fungsi dari sumbernya.
func badanFungsi(t *testing.T, teks, kepala string) string {
	t.Helper()
	i := strings.Index(teks, kepala)
	if i < 0 {
		t.Fatalf("fungsi %q tidak ditemukan", kepala)
	}
	j := strings.Index(teks[i:], "\n}\n")
	if j < 0 {
		t.Fatalf("ujung fungsi %q tidak ditemukan", kepala)
	}
	return teks[i : i+j]
}

// TestSubmitSummaryUrutanTerkunci - AC 24 spec.
//
// ⛔ Urutannya bagian dari kebenaran: nomor lebih dahulu (salinan warisan
// berkunci `PL_NUMBER`), rekap sesudah nomor, warisan paling akhir.
func TestSubmitSummaryUrutanTerkunci(t *testing.T) {
	badan := badanFungsi(t, sumberSummary(t), "func (s *SummaryPremiumList) Submit(")
	urut := []string{
		"s.siapkan(",
		"DalamTransaksi(ctx",
		"s.nomor.terbitkanDalam(",
		"s.rekapDalam(",
		"ringkas.GantiRekap(",
		"ringkas.SumberWarisan(",
		"warisan.Ganti(",
	}
	lalu := -1
	for _, jejak := range urut {
		i := strings.Index(badan, jejak)
		if i < 0 {
			t.Fatalf("Submit tidak memanggil %s", jejak)
		}
		if i < lalu {
			t.Errorf("%s dipanggil sebelum langkah sebelumnya; urutan terkunci: %v", jejak, urut)
		}
		lalu = i
	}
	// ⛔ SATU transaksi - bukan dua. Transaksi kedua menghidupkan kembali
	// titik potong Pega.
	if n := strings.Count(badan, "DalamTransaksi("); n != 1 {
		t.Errorf("Submit membuka %d transaksi, mau 1", n)
	}
}

// TestPenomoranSatuFungsiUntukDuaJalur - gerbang lahir-sekali tidak bercabang.
func TestPenomoranSatuFungsiUntukDuaJalur(t *testing.T) {
	nomor := sumberNomorPL(t)
	terbit := badanFungsi(t, nomor, "func (n *NomorPremiumList) Terbitkan(")
	if !strings.Contains(terbit, "n.terbitkanDalam(") {
		t.Error("Terbitkan tidak memakai terbitkanDalam; dua salinan rantai penomoran")
	}
	if strings.Contains(terbit, "penghitung.UrutNomorBerikut(") {
		t.Error("Terbitkan menaikkan penghitung sendiri, di luar terbitkanDalam")
	}
}

// TestLangkah15TidakDitiru - pl7.
//
// ⛔ `@replaceAll(PL_NUMBER, .MM.YYYY., …)` tidak pernah cocok di produksi.
// Kode yang menggeser periode nomor saat submit adalah perilaku baru, bukan
// replikasi - dan keputusannya milik pemilik Pega.
func TestLangkah15TidakDitiru(t *testing.T) {
	for _, berkas := range []string{"polis_summary.go", "polis_nomor.go"} {
		isi, err := os.ReadFile(berkas)
		if err != nil {
			t.Fatal(err)
		}
		var kode []string
		for _, b := range strings.Split(string(isi), "\n") {
			if !strings.HasPrefix(strings.TrimSpace(b), "//") {
				kode = append(kode, b)
			}
		}
		k := strings.Join(kode, "\n")
		for _, jejak := range []string{"ReplaceAll(", "NextMMYY", "CurrentMMYY"} {
			if strings.Contains(k, jejak) {
				t.Errorf("%s memuat %q - langkah 15 ditiru, padahal pl7 menolaknya", berkas, jejak)
			}
		}
	}
}

// TestSummaryTanpaOracleDitolakTerang - bukan panik, bukan hijau palsu.
func TestSummaryTanpaOracleDitolakTerang(t *testing.T) {
	s := New(nil).SummaryPremiumList()
	pelaku := Pelaku{AkunID: "UJI-OPR"}
	if _, err := s.Submit(context.Background(), pelaku, "P1"); !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("Submit tanpa Oracle: %v", err)
	}
	if _, err := s.Lihat(context.Background(), pelaku, "P1"); !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("Lihat tanpa Oracle: %v", err)
	}
	if _, err := s.Submit(context.Background(), Pelaku{}, "P1"); err == nil {
		t.Error("Submit tanpa identitas diterima")
	}
}
