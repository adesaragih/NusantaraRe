package services

// Baris adjustment - TANPA Oracle.
//
// Pemilik: tiket 03.
//
// Dibaca sesudah: adjustment.go. Perhitungan spreading diuji terpisah di
// spreading_test.go, sebab ia tumbuh menjadi seam tersendiri sesudah
// pembacaan ulang XML 26 September 2026.

import (
	"testing"

	"nusantarare/inti/kontrak"
	"nusantarare/inti/uang"
	"nusantarare/inti/utils"
	"nusantarare/modul/claimlife/models"
)

func TestKolomDiwarisiDikunci(t *testing.T) {
	mau := []string{
		"SHARE_NUSANTARA_RE", "CEDING_RETENTION", "SUM_REASURED", "SUM_INSURED",
		"SHARE_RETRO", "RETROCEDED_SHARE", "CURRENCYID", "CURRENCY",
	}
	got := KolomDiwarisi()
	if len(got) != len(mau) {
		t.Fatalf("kolom diwarisi %d, mau %d", len(got), len(mau))
	}
	for i := range mau {
		if got[i] != mau[i] {
			t.Errorf("posisi %d = %q, mau %q", i, got[i], mau[i])
		}
	}
	// ⛔ Inti AC 5: status TIDAK pernah diwarisi.
	for _, k := range got {
		if k == "STS_REJECT" {
			t.Error("STS_REJECT ikut diwarisi; status adalah keputusan per baris (ADR-U-0011)")
		}
	}
	// Salinan, bukan daftar aslinya: pemanggil tidak boleh dapat mengubahnya.
	got[0] = "DIUBAH"
	if KolomDiwarisi()[0] == "DIUBAH" {
		t.Error("KolomDiwarisi mengembalikan daftar aslinya, bukan salinan")
	}
}

// Baris baru TIDAK mewarisi status maupun tanggal akseptasi.
func TestBarisBaruTidakMewarisiStatus(t *testing.T) {
	pertama := models.BarisAdjustment{JumlahKlaim: uang.Money{Currency: "IDR"}}
	baru := models.BarisAdjustment{KodeStatus: "1"}
	WarisiKolom(pertama, &baru)
	if baru.KodeStatus != "" {
		t.Errorf("KodeStatus = %q; baris yang belum disimpan ke Outstanding "+
			"tidak punya status (AC 11)", baru.KodeStatus)
	}
	if !baru.TanggalAkseptasi.IsZero() {
		t.Error("TanggalAkseptasi distempel saat insert; ia milik aksi akseptasi (AC 12)")
	}
	if baru.JumlahKlaim.Currency != "IDR" {
		t.Errorf("mata uang tidak diwarisi: %q", baru.JumlahKlaim.Currency)
	}
}

// TestBarisKeduaMewarisiDelapanKolomTanpaStatus mengunci SetIndexAdjustmentList
// langkah 3: baris terakhir menyalin baris PERTAMA, dan status tidak ikut.
func TestBarisKeduaMewarisiDelapanKolomTanpaStatus(t *testing.T) {
	// Tiap kolom warisan diberi nilai BERBEDA, supaya satu kolom yang lupa
	// disalin tidak tertutup nilai kolom tetangganya.
	uangUji := func(s string) uang.Money {
		d, err := utils.ParseDecimal(s)
		if err != nil {
			t.Fatalf("uang %q: %v", s, err)
		}
		return uang.Money{Amount: d, Currency: "IDR"}
	}
	pertama := models.BarisAdjustment{
		ShareNusantaraRe: uangUji("11"),
		CedingRetention:  uangUji("22"),
		SumReasured:      uangUji("33"),
		SumInsured:       uangUji("44"),
		ShareRetro:       uangUji("55"),
		RetrocededShare:  uangUji("66"),
		CurrencyID:       "UJI-CUR-1",
		JumlahKlaim:      uangUji("77"),
	}
	p := models.Peserta{MataUang: "IDR"}
	if err := TambahBaris(&p, pertama); err != nil {
		t.Fatalf("baris pertama: %v", err)
	}
	// Baris pertama sudah disimpan ke Outstanding.
	p.Baris[0].KodeStatus = kontrak.KodeOutstanding

	if err := TambahBaris(&p, models.BarisAdjustment{KodeStatus: kontrak.KodeAksep}); err != nil {
		t.Fatalf("baris kedua: %v", err)
	}
	if len(p.Baris) != 2 {
		t.Fatalf("baris = %d, mau 2", len(p.Baris))
	}
	kedua := p.Baris[1]

	// Kedelapan kolom, satu per satu. Daftar nama kolomnya dikunci terpisah
	// oleh TestKolomDiwarisiDikunci; di sini yang diuji NILAInya benar-benar
	// pindah.
	for _, k := range []struct{ nama, mau, dapat string }{
		{"SHARE_NUSANTARA_RE", "11 IDR", kedua.ShareNusantaraRe.String()},
		{"CEDING_RETENTION", "22 IDR", kedua.CedingRetention.String()},
		{"SUM_REASURED", "33 IDR", kedua.SumReasured.String()},
		{"SUM_INSURED", "44 IDR", kedua.SumInsured.String()},
		{"SHARE_RETRO", "55 IDR", kedua.ShareRetro.String()},
		{"RETROCEDED_SHARE", "66 IDR", kedua.RetrocededShare.String()},
		{"CURRENCYID", "UJI-CUR-1", kedua.CurrencyID},
		{"CURRENCY", "IDR", kedua.JumlahKlaim.Currency},
	} {
		if k.dapat != k.mau {
			t.Errorf("%s tidak diwarisi: %q, mau %q", k.nama, k.dapat, k.mau)
		}
	}

	// ⛔ CLAIM_AMOUNT BUKAN kolom warisan - tiap baris punya jumlahnya sendiri.
	if !kedua.JumlahKlaim.Kosong() {
		t.Errorf("CLAIM_AMOUNT ikut diwarisi (%q); ia bukan salah satu dari kedelapan kolom",
			kedua.JumlahKlaim.String())
	}
	if kedua.KodeStatus != "" {
		t.Errorf("baris kedua berstatus %q; status adalah keputusan per baris (ADR-U-0011)",
			kedua.KodeStatus)
	}
	// Baris pertama tidak ikut berubah.
	if p.Baris[0].KodeStatus != kontrak.KodeOutstanding {
		t.Errorf("baris pertama berubah menjadi %q", p.Baris[0].KodeStatus)
	}
}

// TestTambahBarisMenandaiPesertaDipilih mengunci SetIndexAdjustmentList
// langkah 1: `.IsCheck = true` lahir bersama barisnya (AC 39).
func TestTambahBarisMenandaiPesertaDipilih(t *testing.T) {
	p := models.Peserta{MataUang: "IDR"}
	if p.IsCheck != "" {
		t.Fatalf("prasyarat: IsCheck sudah terisi")
	}
	if err := TambahBaris(&p, models.BarisAdjustment{}); err != nil {
		t.Fatalf("TambahBaris: %v", err)
	}
	if p.IsCheck != "true" {
		t.Errorf("IsCheck = %q, mau true", p.IsCheck)
	}
}

// TestTandaiOutstandingHanyaMenyentuhBarisTanpaStatus mengunci gerbang
// `.PrintFaceClaim==""` pada langkah 22.1.3.1.
func TestTandaiOutstandingHanyaMenyentuhBarisTanpaStatus(t *testing.T) {
	pohon := models.PohonKlaim{}
	pohon.Klaim.Peserta = []models.Peserta{{
		MataUang: "IDR",
		Baris: []models.BarisAdjustment{
			{ID: "A", KodeStatus: ""},
			{ID: "B", KodeStatus: kontrak.KodeAksep},
			{ID: "C", KodeStatus: ""},
		},
	}}
	n := TandaiOutstanding(&pohon)
	if n != 2 {
		t.Fatalf("baris tersentuh = %d, mau 2", n)
	}
	baris := pohon.Klaim.Peserta[0].Baris
	if baris[0].KodeStatus != kontrak.KodeOutstanding || baris[2].KodeStatus != kontrak.KodeOutstanding {
		t.Errorf("baris tanpa status tidak menjadi Outstanding: %v", baris)
	}
	if baris[1].KodeStatus != kontrak.KodeAksep {
		t.Errorf("baris yang sudah diaksep ditimpa menjadi %q; nol ditulis MENURUT AKSI, "+
			"dan baris yang sudah disimpan tidak disentuh ulang", baris[1].KodeStatus)
	}
}

// TestSimpanKeOutstandingTidakMenstempelTanggalAkseptasi - AC 43.
func TestSimpanKeOutstandingTidakMenstempelTanggalAkseptasi(t *testing.T) {
	pohon := models.PohonKlaim{}
	pohon.Klaim.Peserta = []models.Peserta{{
		Baris: []models.BarisAdjustment{{ID: "A"}},
	}}
	TandaiOutstanding(&pohon)
	if !pohon.Klaim.Peserta[0].Baris[0].TanggalAkseptasi.IsZero() {
		t.Error("ACCEPTATION_DATE distempel saat simpan ke Outstanding; " +
			"menyimpan bukan mengaksep (AC 43)")
	}
}
