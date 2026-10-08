package models_test

import (
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

// Panel AdjustmentDetail mengikuti layout Section `AdjustmentDetail_Section` / `AdjustmentDetail` (work owner 08-10-2026
// "perbaiki tampilan", screenshot Pega "ikuti dan rapihkan"): tiga tab, blok atas Inline grid double bersarang
// [Type, Currency | Payment Type, DLA] | [Value In IDR, RNM Share, Allocation | Share], checkbox "Transfer Direct to
// Kasir", tabel bebas 6 x 4, Payable dan bank berdampingan, Spreading In tidak dibangun (ContainerVisibleWhen NEVER).
func TestLayarAdjustmentIkutLayoutXML(t *testing.T) {
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarAdjustment, []models.Baris{{"Type": "1", "IsSubjectivity": "true"}})
	ts := models.Evaluasi(h, models.LayarAdjustment(1), false)

	var cari func(xs []models.Tata, ok func(models.Tata) bool) *models.Tata
	cari = func(xs []models.Tata, ok func(models.Tata) bool) *models.Tata {
		for i := range xs {
			if ok(xs[i]) {
				return &xs[i]
			}
			if x := cari(xs[i].Anak, ok); x != nil {
				return x
			}
		}
		return nil
	}
	jalur := func(xs []models.Tata) []string {
		var out []string
		for _, x := range xs {
			out = append(out, x.Jalur)
		}
		return out
	}

	if len(ts) != 1 || ts[0].Letak != models.LetakTab {
		t.Fatalf("panel bukan layout group Tab: %+v", ts)
	}
	var tab []string
	for _, a := range ts[0].Anak {
		tab = append(tab, a.Label)
	}
	if strings.Join(tab, "|") != "Detail Chronology to Committe|Adjusment List Detail|Subjectivity" {
		t.Fatalf("tab panel %v", tab)
	}

	atas := cari(ts, func(x models.Tata) bool {
		return x.Letak == models.LetakDua && cari(x.Anak, func(y models.Tata) bool { return y.Label == "Value In IDR" }) != nil
	})
	if atas == nil || len(atas.Anak) != 2 || atas.Anak[0].Letak != models.LetakDua || len(atas.Anak[0].Anak) != 2 ||
		len(atas.Anak[0].Anak[0].Anak) != 2 || len(atas.Anak[0].Anak[1].Anak) != 3 || len(atas.Anak[1].Anak) != 3 ||
		atas.Anak[1].Anak[2].Letak != models.LetakDua {
		t.Fatalf("blok atas bukan [[Type, Currency | Payment Type, DLA x2] | Value, RNM Share, Allocation|Share]: %+v", atas)
	}
	kasir := cari(ts, func(x models.Tata) bool {
		return x.Letak == models.LetakSebaris && len(x.Anak) > 1 &&
			x.Anak[0].Jalur == models.JalurAdj(1, "DirectToKasir")
	})
	if kasir == nil || kasir.Anak[1].Label != "Transfer Direct to Kasir" {
		t.Fatalf("checkbox Transfer Direct to Kasir: %+v", kasir)
	}

	tabel := cari(ts, func(x models.Tata) bool { return x.Letak == models.LetakTabel })
	if tabel == nil || len(tabel.Anak) != 4 {
		t.Fatalf("tabel Deductible / Treaty Gross bukan 4 baris: %+v", tabel)
	}
	var judul []string
	for _, c := range tabel.Anak[0].Anak {
		judul = append(judul, c.Label)
	}
	if got := judul; len(got) != 6 || got[0] != "Dedutible Type" || got[3] != "Treaty Gross (100%)" ||
		got[5] != "RNM Share Claim" {
		t.Fatalf("judul kolom tabel %v", got)
	}
	for i, b := range tabel.Anak {
		if len(b.Anak) != 6 {
			t.Fatalf("baris tabel %d punya %d sel, mau 6", i, len(b.Anak))
		}
	}
	j := models.JalurAdj
	if got := jalur(tabel.Anak[2].Anak); got[0] != j(1, "IndividualRiskType") || got[5] != j(1, "IndividualRiskRNM") {
		t.Fatalf("baris Deductible %v", got)
	}
	if b := tabel.Anak[3].Anak; b[0].Label != "Total" || b[3].Jalur != j(1, "ProposeAdjustmentValue") ||
		b[5].Jalur != j(1, "AdjustmentValue") {
		t.Fatalf("baris Total %+v", b)
	}

	bayar := cari(ts, func(x models.Tata) bool {
		return x.Letak == models.LetakDua && cari(x.Anak, func(y models.Tata) bool { return y.Label == "Payable To" }) != nil
	})
	if bayar == nil || len(bayar.Anak) != 2 {
		t.Fatalf("Payable To | Name of Bank tidak berdampingan: %+v", bayar)
	}

	if cari(ts, func(x models.Tata) bool { return x.Label == "Spreading In" }) != nil {
		t.Errorf("Spreading In tampil, padahal XML ContainerVisibleWhen NEVER")
	}
	if cari(ts, func(x models.Tata) bool { return x.Label == "Spreading Out" }) == nil {
		t.Errorf("Spreading Out hilang")
	}
}
