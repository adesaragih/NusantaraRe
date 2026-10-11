package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// Tiket 64b - penahan tahap Binding mengikuti pyDisabledWhen EmailSectionCeding: tombol Send Email Accept (CARI1 / CARI23 /
// ProtectZipCode.CARI1 + ConfirmBinding), tombol Binding (CARI1 / CARI23 / ProtectDeductible.CARI2 / CARI4), Banding
// (CARI1 / CARI23, 57b).

const (
	pesanDedTambah = " Deductible on location 1 object item 1 Coverage 1 should be added!"
	pesanDedKosong = " Deductible on location 1 object item 1 Coverage 1  DeductibleList 2 can't be empty!"
	pesanZipMaster = "ZipCode on Location 1 doesn’t exist in Master!"
	pesanProvinsi  = "Province on Location 1can't be null!"
)

// TestProteksiDeductibleBinding - 9.2.30.6.6 / 13.2 hanya bila deductibleBinding; per coverage sesudah pesan coverage, baris
// kosong SELURUHNYA saja (satu medan terisi = bukan kosong).
func TestProteksiDeductibleBinding(t *testing.T) {
	var p pengumpulPesan
	it := itemSah()
	it.Coverages = []models.CoverageObjek{
		{Coverage: "FLEXAS", CoverageNote: "UJI", OldID: "FLEXAS", AccumulationCode: "UJI"},
		{Coverage: "", CoverageNote: "", OldID: "4.3", AccumulationCode: "UJI", Deductibles: []models.Deductible{
			{TypeDeductible: "1"}, {}, {Amount: d("5")}, {TimeExcess: d("3")}, {MinMax: "2"}, {Condition: "1"}}},
	}
	proteksiItem(&p, "UJI-NB-1", 2, 3, it, true)
	mau := []string{" Deductible on location 2 object item 3 Coverage 1 should be added!",
		" Coverage 2 on location 2 object item 3 can't be empty!",
		" Deductible on location 2 object item 3 Coverage 2  DeductibleList 2 can't be empty!"}
	if !reflect.DeepEqual(p.pesan, mau) {
		t.Fatalf("%q", p.pesan)
	}
	if !reflect.DeepEqual(p.flag[0], []flagProteksi{flagDeductible4}) || !reflect.DeepEqual(p.flag[2], []flagProteksi{flagDeductible2}) {
		t.Errorf("flag %v", p.flag)
	}
	var q pengumpulPesan
	proteksiItem(&q, "UJI-NB-1", 2, 3, it, false)
	if len(q.pesan) != 1 {
		t.Errorf("tanpa binding %q", q.pesan)
	}
}

// layananBinding64b - kasus Binding dengan satu lokasi yang dapat dirusak; ConfirmBinding `confirm`.
func layananBinding64b(st *submitTiruan, confirm string, ubah func(*models.ObjekFire), kasus func(*models.Kasus)) *Service {
	svc := layananBinding(st, confirm, kasus)
	o := lokasiBinding()
	if ubah != nil {
		ubah(&o)
	}
	svc.objek.(objekTiruan).ada["UJI-NB-1"] = []models.ObjekFire{o}
	return svc
}

// TestPenahanTahapBinding - matriks pesan × tombol: periksa Accept (pesan / pesanBinding), Submit Accept (sesudah Binding),
// tombol Binding (sebelum Binding), Banding; Accept tahap marketing tidak berubah (deductible tidak diperiksa).
func TestPenahanTahapBinding(t *testing.T) {
	ctx := context.Background()
	dedTambah := func(o *models.ObjekFire) { o.Items[0].Coverages[0].Deductibles = nil }
	dedKosong := func(o *models.ObjekFire) {
		o.Items[0].Coverages[0].Deductibles = append(o.Items[0].Coverages[0].Deductibles, models.Deductible{})
	}
	zip := func(o *models.ObjekFire) { o.ZipCode = "99999" }
	prov := func(o *models.ObjekFire) { o.Province = "" }
	for nama, c := range map[string]struct {
		ubah                    func(*models.ObjekFire)
		kasus                   func(*models.Kasus)
		pesan                   string
		accept, tombol, banding bool
	}{
		"deductible ditambah": {dedTambah, nil, pesanDedTambah, false, true, false},
		"deductible kosong":   {dedKosong, nil, pesanDedKosong, false, true, false},
		"zip master":          {zip, nil, pesanZipMaster, true, false, false},
		"provinsi":            {prov, nil, pesanProvinsi, true, true, true},
		"ceding":              {nil, func(k *models.Kasus) { k.General.CedingCoName = "" }, pesanCedingCoKosong, true, true, false},
		"sob":                 {nil, func(k *models.Kasus) { k.General.SourceOfBusinessID = "" }, pesanSOBKosong, true, true, false},
		"scoring":             {nil, func(k *models.Kasus) { k.RiskScoring = models.RingkasanScoring{} }, pesanScoringKosong, false, false, false},
	} {
		// Periksa Accept sebelum Binding: pesan = belum Binding + penahan Accept; pesanBinding = penahan tombol Binding.
		st := submitUji()
		svc := layananBinding64b(st, "", c.ubah, c.kasus)
		h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "1"})
		if err != nil || h.Pesan[0] != pesanBelumBinding || dalam(h.Pesan, c.pesan) != c.accept || dalam(h.PesanBinding, c.pesan) != c.tombol {
			t.Errorf("%s periksa: %v %q / %q", nama, err, h.Pesan, h.PesanBinding)
		}
		// Tombol Binding: 422 hanya oleh penahan tombol.
		_, err = svc.TekanBinding(ctx, pelakuMkt, "UJI-NB-1")
		var g GalatProteksi
		if errors.As(err, &g) != c.tombol || (c.tombol && !dalam(g.Pesan, c.pesan)) || (!c.tombol && err != nil) {
			t.Errorf("%s tombol: %v", nama, err)
		}
		// Submit Accept sesudah Binding: 422 hanya oleh penahan Accept.
		st = submitUji()
		svc = layananBinding64b(st, "1", c.ubah, c.kasus)
		_, err = svc.AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "UJI"})
		if errors.Is(err, ErrProteksiSubmit) != c.accept || (!c.accept && err != nil) {
			t.Errorf("%s accept: %v", nama, err)
		}
		// Banding: penahan EmailSectionCeding; pesanBinding [] (bukan Accept).
		h, err = svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "4"})
		if err != nil || dalam(h.Pesan, c.pesan) != c.banding || h.PesanBinding == nil || len(h.PesanBinding) != 0 {
			t.Errorf("%s banding: %v %q %#v", nama, err, h.Pesan, h.PesanBinding)
		}
	}
	// Premi 0 (ProtectSpreading.CARI23): menahan Accept DAN tombol Binding.
	st0 := submitUji()
	svc0 := layananBinding64b(st0, "", nil, nil)
	svc0.spreading.(*spreadingTiruan).baca = []models.KasusSpreading{kasusDuaItem("0", "USD")}
	const premiNol = "Premi on paymnet can't be 0"
	if h, err := svc0.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "1"}); err != nil || !dalam(h.Pesan, premiNol) ||
		!dalam(h.PesanBinding, premiNol) {
		t.Errorf("premi 0 periksa: %v %#v", err, h)
	}
	var g0 GalatProteksi
	if _, err := svc0.TekanBinding(ctx, pelakuMkt, "UJI-NB-1"); !errors.As(err, &g0) || !dalam(g0.Pesan, premiNol) {
		t.Errorf("premi 0 tombol: %v", err)
	}
	// Bersih: periksa Accept sesudah Binding - pesan dan pesanBinding kosong (bukan nil).
	st := submitUji()
	h, err := layananBinding64b(st, "1", nil, nil).PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "1"})
	if err != nil || len(h.Pesan) != 0 || h.PesanBinding == nil || len(h.PesanBinding) != 0 {
		t.Errorf("bersih: %v %#v", err, h)
	}
	// Accept tahap marketing: deductible TIDAK diperiksa, pesanBinding nil (handler menulis []).
	st = submitUji()
	svc := layananSubmit(st)
	if h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "1"}); err != nil || len(h.Pesan) != 0 || len(h.PesanBinding) != 0 {
		t.Errorf("marketing: %v %#v", err, h)
	}
}
