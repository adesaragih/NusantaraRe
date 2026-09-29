package services_test

// Baris adjustment PERTAMA - GILIRAN-13 butir bo. TANPA Oracle.
//
// Dibaca sesudah: adjustment.go (BarisPertama) dan hasilkomite.go (cabang
// baris pertama di Putaran.Tambah).

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

// TestBarisPertamaLahirKosong - `Add` b17937 pada grid kosong.
//
// ⛔ KOSONG, termasuk mata uang. `SetIndexAdjustmentList` langkah 3
// (b570-744) menyalin `.AdjustmentList(1)` ke `.AdjustmentList(<LAST>)`, dan
// pada grid kosong keduanya baris yang SAMA - nol yang diwarisi.
func TestBarisPertamaLahirKosong(t *testing.T) {
	p := models.Peserta{ID: "P-1", MataUang: "IDR"}
	baru, err := services.BarisPertama(p)
	if err != nil {
		t.Fatalf("BarisPertama: %v", err)
	}
	if baru.KodeStatus != "" {
		t.Errorf("baris pertama berkode %q; baris lahir TANPA status - "+
			"STS_REJECT ditulis Save Outstanding (22.1.3.2), bukan Add", baru.KodeStatus)
	}
	if baru.JumlahKlaim.Currency != "" || baru.CurrencyID != "" {
		t.Errorf("mata uang %q/%q diisi dari peserta; Add b17947 menambah baris KOSONG",
			baru.JumlahKlaim.Currency, baru.CurrencyID)
	}
	// Seluruhnya nilai nol - bukan hanya medan yang terpikir hari ini.
	if !reflect.DeepEqual(baru, models.BarisAdjustment{}) {
		t.Errorf("baris pertama membawa isian: %+v", baru)
	}
}

// TestBarisPertamaHanyaUntukGridKosong - baris kedua lahir lewat putaran.
func TestBarisPertamaHanyaUntukGridKosong(t *testing.T) {
	p := models.Peserta{ID: "P-1", Baris: []models.BarisAdjustment{{ID: "A-1"}}}
	if _, err := services.BarisPertama(p); !errors.Is(err, services.ErrBarisTidakSah) {
		t.Errorf("peserta berbaris: galat = %v, mau ErrBarisTidakSah", err)
	}
}

// TestBarisPertamaMenghormatiGerbangSTSReject - tujuh gerbang yang sudah ada.
//
// `pyDisabledWhen` `.STS_REJECT=='1' || .STS_REJECT=='2'` (b2628, b4682,
// b5059, b5870, b6152, b7335, b15234) membekukan isian layar Detail bagi
// peserta yang sudah DIPUTUS - diaksep ATAU ditolak. Butir bo memberlakukannya
// pada baris pertama.
func TestBarisPertamaMenghormatiGerbangSTSReject(t *testing.T) {
	for _, kode := range []string{models.KodeAksep, models.KodeDitolak, " 2 "} {
		p := models.Peserta{ID: "P-1", KodeStatus: kode}
		if _, err := services.BarisPertama(p); !errors.Is(err, services.ErrPesertaSudahDiputus) {
			t.Errorf("peserta berkode %q: galat = %v, mau ErrPesertaSudahDiputus", kode, err)
		}
	}
	for _, kode := range []string{"", models.KodeOutstanding} {
		p := models.Peserta{ID: "P-1", KodeStatus: kode}
		if _, err := services.BarisPertama(p); err != nil {
			t.Errorf("peserta berkode %q ditolak: %v", kode, err)
		}
	}
}

// TestPemegangClaimAnalisSamaDenganPeranPutaran - dua gerbang, satu peran.
//
// Cabang baris pertama menuntut PEMEGANG tahap Claim Analis; jalur putaran
// menuntut `PeranSimpanOutstanding`. Keduanya satu tombol (`Add` b17937,
// tampil hanya bila `pyWorkPage.pyPosition =='ReasLifeSPV'` b18160). Bila
// kelak keduanya berbeda, satu tombol akan menuntut dua peran.
func TestPemegangClaimAnalisSamaDenganPeranPutaran(t *testing.T) {
	peran, ada := models.PeranPemegangTahap(models.TahapClaimAnalis)
	if !ada || peran != services.PeranSimpanOutstanding {
		t.Errorf("pemegang Claim Analis %q (%v), peran putaran %q - satu tombol, dua peran",
			peran, ada, services.PeranSimpanOutstanding)
	}
}

// TestTombolAddAdjustmentVERBATIMDariKorpus - bentuk `Add` dan `Delete`.
//
// ⚠️ Melewati bila korpus tidak terjangkau - dan mengatakannya.
//
// ⚠️ b17991 (activity pewaris delapan kolom) SENGAJA tidak diperiksa di sini:
// namanya memuat pola indeks posisi yang dilarang penjaga AC 62
// (`komite_statik_test.go`) di luar komentar. Ia dikutip di komentar
// `BarisPertama`, bukan dikelabui dengan memecah teksnya.
func TestTombolAddAdjustmentVERBATIMDariKorpus(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\Claim Life\Section\ClaimLifeDetailGCNM.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); tombol tidak terperiksa", err)
	}
	baris := strings.Split(string(isi), "\n")
	for nomor, mau := range map[int]string{
		17126: `<pyPageListProperty>.AdjustmentList</pyPageListProperty>`,
		17937: `<pyLabel>Add</pyLabel>`,
		17947: `<pyAction>addRow</pyAction>`,
		18160: `<pyCondition>pyWorkPage.pyPosition =='ReasLifeSPV'</pyCondition>`,
		19120: `<pyLabel>Delete</pyLabel>`,
		19130: `<pyAction>deleteRow</pyAction>`,
		19399: `<pyCondition>.PrintFaceClaim == '' &amp;&amp; pyWorkPage.pyPosition =='ReasLifeSPV'</pyCondition>`,
	} {
		if got := strings.TrimSpace(baris[nomor-1]); got != mau {
			t.Errorf("b%d = %q, mau %q", nomor, got, mau)
		}
	}
}
