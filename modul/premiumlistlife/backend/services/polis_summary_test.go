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
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
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

// TestSimpanDalamUrutanTerkunci - AC 24 spec.
//
// ⛔ Urutannya bagian dari kebenaran: rekap TERSIMPAN dibaca lebih dahulu
// (tanpa rekap, Confirm ditolak sebelum penghitung nomor dikunci), lalu
// nomor (salinan warisan berkunci `PL_NUMBER`), warisan paling akhir.
//
// ⛔ [keputusan work owner 03-10-2026] Confirm TIDAK menghitung rekap: ia
// dihitung saat Save Data / simpan peserta CSV (`perbaruiRekapDalam`).
func TestSimpanDalamUrutanTerkunci(t *testing.T) {
	badan := badanFungsi(t, sumberSummary(t), "func (s *SummaryPremiumList) simpanDalam(")
	for _, larangan := range []string{"s.rekapDalam(", "ringkas.GantiRekap(", "perbaruiRekapDalam("} {
		if strings.Contains(badan, larangan) {
			t.Errorf("simpanDalam (Confirm) masih menghitung rekap lewat %s", larangan)
		}
	}
	urut := []string{
		"ringkas.BacaRekap(",
		"ErrSummaryBelumAda",
		"s.nomor.terbitkanDalam(",
		"ringkas.TulisNomorRekap(",
		"ringkas.KepalaSummaryWarisan(",
		"summaryWarisan.Ganti(",
		"ringkas.SumberWarisan(",
		"warisan.Ganti(",
	}
	lalu := -1
	for _, jejak := range urut {
		i := strings.Index(badan, jejak)
		if i < 0 {
			t.Fatalf("simpanDalam tidak memanggil %s", jejak)
		}
		if i < lalu {
			t.Errorf("%s dipanggil sebelum langkah sebelumnya; urutan terkunci: %v", jejak, urut)
		}
		lalu = i
	}
	// ⛔ Ia TIDAK membuka transaksi sendiri - transaksinya milik pemanggil,
	// supaya simpan dan penutupan kasus satu commit.
	if strings.Contains(badan, "DalamTransaksi(") {
		t.Error("simpanDalam membuka transaksi sendiri; simpan dan penutupan akan terpisah")
	}
}

// TestSimpanSebelumTutupDalamSatuTransaksi - tiket 05b.
//
// ⛔ `Utility1` berdiri di ANTARA `Confirm` dan `END52`. Menutup lebih dahulu
// lalu menyimpan di transaksi lain meninggalkan kasus Resolved-Completed
// tanpa rekap bila simpannya gagal.
func TestSimpanSebelumTutupDalamSatuTransaksi(t *testing.T) {
	isi, err := os.ReadFile("polis_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	badan := badanFungsi(t, string(isi), "func (p *Penawaran) terapkan(")
	iTx := strings.Index(badan, "DalamTransaksi(ctx")
	iGerbang := strings.Index(badan, "if akibat.SimpanPolis")
	iSimpan := strings.Index(badan, ".simpanDalam(ctx, tx,")
	iTutup := strings.Index(badan, "kerja.TutupKasus(")
	if iTx < 0 || iGerbang < 0 || iSimpan < 0 || iTutup < 0 {
		t.Fatalf("terapkan kehilangan salah satu dari: transaksi, gerbang SimpanPolis, simpanDalam, TutupKasus")
	}
	if !(iTx < iGerbang && iGerbang < iSimpan && iSimpan < iTutup) {
		t.Error("urutan terapkan harus: transaksi → gerbang SimpanPolis → simpanDalam → TutupKasus")
	}
	if n := strings.Count(badan, "DalamTransaksi("); n != 1 {
		t.Errorf("terapkan membuka %d transaksi, mau 1", n)
	}
}

// TestSubmitSummaryHanyaDariTahapSummary - `finishAssignment` milik Assignment1.
func TestSubmitSummaryHanyaDariTahapSummary(t *testing.T) {
	badan := badanFungsi(t, sumberSummary(t), "func (s *SummaryPremiumList) Submit(")
	iGerbang := strings.Index(badan, "models.TahapPolisSummary")
	iTerap := strings.Index(badan, ".terapkan(")
	if iGerbang < 0 || iTerap < 0 || iGerbang > iTerap {
		t.Error("Submit tidak menolak tahap selain Input Premium Summary sebelum menutup kasus")
	}
	if !strings.Contains(badan, "models.PenyelesaianSummary()") {
		t.Error("Submit tidak memakai models.PenyelesaianSummary (Transition2 -> END52)")
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
	pelaku := inti.Pelaku{AkunID: "UJI-OPR"}
	if _, err := s.Submit(context.Background(), pelaku, "P1", time.Now()); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("Submit tanpa Oracle: %v", err)
	}
	if _, err := s.Lihat(context.Background(), pelaku, "P1"); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("Lihat tanpa Oracle: %v", err)
	}
	if _, err := s.Submit(context.Background(), inti.Pelaku{}, "P1", time.Now()); err == nil {
		t.Error("Submit tanpa identitas diterima")
	}
}

// TestRekapDihitungSaatSave - keputusan work owner 03-10-2026: rekap dihitung
// dan disimpan di transaksi Save Data DAN simpan peserta CSV, dan rekap yang
// tidak dapat dihitung DIHAPUS (bukan dibiarkan usang).
func TestRekapDihitungSaatSave(t *testing.T) {
	for berkas, kepala := range map[string]string{
		"polis_datapolis.go": "func (f *FormDataPolis) Simpan(",
		"polis_unggah.go":    "func (u *UnggahPremiumList) Simpan(",
	} {
		isi, err := os.ReadFile(berkas)
		if err != nil {
			t.Fatal(err)
		}
		badan := badanFungsi(t, string(isi), kepala)
		iTx := strings.Index(badan, "DalamTransaksi(")
		iRekap := strings.Index(badan, ".perbaruiRekapDalam(ctx, tx, polisID)")
		if iTx < 0 || iRekap < iTx {
			t.Errorf("%s: rekap tidak diperbarui di dalam transaksi simpan", berkas)
		}
	}
	badan := badanFungsi(t, sumberSummary(t), "func (s *SummaryPremiumList) perbaruiRekapDalam(")
	for _, jejak := range []string{"ringkas.HapusRekap(", "ringkas.GantiRekap(",
		"repository.ErrPolisTanpaPeserta", "penomor.ErrTipePLTanpaCabang"} {
		if !strings.Contains(badan, jejak) {
			t.Errorf("perbaruiRekapDalam tanpa %s", jejak)
		}
	}
}

// TestWPCDitulisSaatConfirm - WPC dihitung di Go (tanpa POOLDATA.GETQUARTER)
// dan ditulis SESUDAH nomor terbit, di transaksi Confirm (03-10-2026).
func TestWPCDitulisSaatConfirm(t *testing.T) {
	badan := badanFungsi(t, sumberSummary(t), "func (s *SummaryPremiumList) simpanDalam(")
	iNomor := strings.Index(badan, "s.nomor.terbitkanDalam(")
	iWPC := strings.Index(badan, "models.WPCPolis(identitas.Tipe, nomor.Periode)")
	iTulis := strings.Index(badan, ".TulisWPC(ctx, tx, polisID, wpc)")
	if iNomor < 0 || iWPC < iNomor || iTulis < iWPC {
		t.Errorf("WPC tidak dihitung lalu ditulis sesudah nomor terbit (nomor=%d wpc=%d tulis=%d)", iNomor, iWPC, iTulis)
	}
	if strings.Contains(badan, "GETQUARTER") {
		t.Error("simpanDalam masih memanggil POOLDATA.GETQUARTER")
	}
}
