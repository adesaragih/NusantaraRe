package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// Tiket 57b - Send Email Banding: port Protection_Act yang sama dengan Accept, tetapi yang menahan hanya pesan ber-flag di
// syarat tombol Banding (EmailSection tahap marketing / EmailSectionCeding tahap Binding).

// TestPengumpulPesanFlag - kembar sekali dengan gabungan flag; penahan menjaga urutan asli; tanpa flag tidak menahan.
func TestPengumpulPesanFlag(t *testing.T) {
	var p pengumpulPesan
	p.tambah("A")
	p.tambahFlag("B", flagMarine2)
	p.tambahFlag("A", flagSpreading23)
	p.tambahFlag("C", flagZipCode1)
	p.tambahFlag("B", flagMarine2)
	p.tambahFlag("", flagSpreading1)
	if !reflect.DeepEqual(p.pesan, []string{"A", "B", "C"}) || len(p.flag[1]) != 1 {
		t.Fatalf("%q %v", p.pesan, p.flag)
	}
	if got := p.penahan(flagBandingMarketing); !reflect.DeepEqual(got, []string{"A", "B"}) {
		t.Errorf("marketing %q", got)
	}
	if got := p.penahan(flagBandingBinding); !reflect.DeepEqual(got, []string{"A"}) {
		t.Errorf("binding %q", got)
	}
	var kosong pengumpulPesan
	if got := kosong.penahan(flagBandingMarketing); got == nil || len(got) != 0 {
		t.Errorf("kosong %#v", got)
	}
}

// layananRusak57b - kasus TestProteksiAcceptGagal: setiap pesan lokasi / item / mata uang + MOID dan Scoring kosong.
func layananRusak57b(st *submitTiruan) *Service {
	rusak := lokasiSah()
	rusak.Province, rusak.City, rusak.ZipCode, rusak.District = "", "", "99999", "UJI KEC"
	rusak.Occupations = []models.OkupasiObjek{{OccupationID: ""}}
	rusak.LossRecords = nil
	it := rusak.Items[0]
	it.Currency, it.ItemTypeID = "", ""
	it.TotalNetRate = d("0.5")
	it.Coverages = []models.CoverageObjek{{Coverage: "FLEXAS", OldID: "FLEXAS", AccumulationCode: "UNKNOWNID"}}
	asing := itemSah()
	asing.Currency = "UJI-XYZ"
	rusak.Items = []models.ItemObjek{it, asing}
	svc := layananSubmit(st, lokasiSah(), rusak)
	k := svc.kasus.(kasusTiruan).ada["UJI-NB-1"]
	k.General.MarketingID, k.RiskScoring.FinalScore = "", nil
	return svc
}

// pesanSpreading1Rusak - pesan layananRusak57b ber-flag ProtectSpreading.CARI1 (menahan di kedua tahap).
var pesanSpreading1Rusak = []string{
	"Currency Pada Lokasi 2 Object Item 2 Bermasalah, Coba Pilih Ulang Currency",
	"Province on Location 2can't be null!",
	"City on Location 2can't be null!",
	"Loss record on location 2 can't be empty !",
	"Occupation ID on location2can't be null!",
	"Currency on location 2 object Item 1 must be choose!",
	"Choose Object Item on location 2 !",
}

// TestProteksiBandingPemetaan - kasus rusak: marketing ditahan CARI1 + ProtectMarine.CARI2 (coverage / akumulasi) +
// ProtectNeRate.CARI1; Binding hanya CARI1. MOID (ProtectMarketing), Scoring (tanpa flag), Master RW (ProtectZipCode) tidak
// menahan. Accept tetap memuat semuanya (16 pesan, TestProteksiAcceptGagal).
func TestProteksiBandingPemetaan(t *testing.T) {
	ctx := context.Background()
	svc := layananRusak57b(submitUji())
	mkt, err := svc.proteksiBanding(ctx, "UJI-NB-1", flagBandingMarketing)
	mauMkt := append(append([]string{}, pesanSpreading1Rusak...),
		" Coverage 1 on location 2 object item 1 can't be empty!",
		" Accumulation on location 2 object item 1 can't be empty!",
		"coverage on location 2 object item 1 can't be Net Rate!")
	if err != nil || !reflect.DeepEqual(mkt, mauMkt) {
		t.Errorf("marketing %v\n%q\nmau\n%q", err, mkt, mauMkt)
	}
	bnd, err := svc.proteksiBanding(ctx, "UJI-NB-1", flagBandingBinding)
	if err != nil || !reflect.DeepEqual(bnd, pesanSpreading1Rusak) {
		t.Errorf("binding %v\n%q", err, bnd)
	}
	if all, err := svc.ProteksiAccept(ctx, "UJI-NB-1", nil); err != nil || len(all) != 16 {
		t.Errorf("accept %v %d", err, len(all))
	}
	// 9.2.7: ZipCode kosong / OccupationList kosong (CARI1) menahan Binding.
	o := lokasiSah()
	o.ZipCode, o.Occupations = "", nil
	svc = layananSubmit(submitUji(), o)
	bnd, err = svc.proteksiBanding(ctx, "UJI-NB-1", flagBandingBinding)
	if err != nil || !dalam(bnd, "ZipCode on Location 1can't be null!") || !dalam(bnd, "Occupation on location 1 can't be null!") {
		t.Errorf("9.2.7: %v %q", err, bnd)
	}
	// Payment: tanpa cicilan -> "List Installment" (5.8 CARI1) menahan; "Total premium RNM ..." (CARI28) tidak.
	svc = layananSubmit(submitUji())
	svc.payment.(*paymentTiruan).ada.Cicilan = map[string][]models.CicilanPayment{}
	for nama, tahan := range map[string][]flagProteksi{"marketing": flagBandingMarketing, "binding": flagBandingBinding} {
		if p, err := svc.proteksiBanding(ctx, "UJI-NB-1", tahan); err != nil || !reflect.DeepEqual(p, []string{"List Installment can't be null!"}) {
			t.Errorf("cicilan %s: %v %q", nama, err, p)
		}
	}
	// Tanpa mata uang payment -> "List Payment can't be null!" (ProtectPayment.CARI1): marketing saja.
	svc = layananSubmit(submitUji())
	svc.spreading.(*spreadingTiruan).baca = []models.KasusSpreading{{}}
	if p, err := svc.proteksiBanding(ctx, "UJI-NB-1", flagBandingMarketing); err != nil || !dalam(p, "List Payment can't be null!") {
		t.Errorf("list payment marketing: %v %q", err, p)
	}
	if p, err := svc.proteksiBanding(ctx, "UJI-NB-1", flagBandingBinding); err != nil || dalam(p, "List Payment can't be null!") {
		t.Errorf("list payment binding: %v %q", err, p)
	}
	// Kategori okupasi (ProtectOccupation.CARI1) tidak menahan.
	o = lokasiSah()
	o.Occupations = []models.OkupasiObjek{{OccupationID: "UJI-9"}}
	svc = layananSubmit(submitUji(), o)
	if a, _ := svc.ProteksiAccept(ctx, "UJI-NB-1", nil); len(a) != 1 {
		t.Fatalf("instrumen kategori: %q", a)
	}
	if p, err := svc.proteksiBanding(ctx, "UJI-NB-1", flagBandingMarketing); err != nil || len(p) != 0 {
		t.Errorf("kategori: %v %q", err, p)
	}
}

// TestProteksiPremiNol - Premi 0 (ProtectSpreading.CARI23) menahan di kedua tahap.
func TestProteksiPremiNol(t *testing.T) {
	ctx := context.Background()
	svc := layananSubmit(submitUji())
	svc.spreading.(*spreadingTiruan).baca = []models.KasusSpreading{kasusDuaItem("0", "USD")}
	a, err := svc.ProteksiAccept(ctx, "UJI-NB-1", nil)
	if err != nil || !dalam(a, "Premi on paymnet can't be 0") {
		t.Fatalf("instrumen: %v %q", err, a)
	}
	for nama, tahan := range map[string][]flagProteksi{"marketing": flagBandingMarketing, "binding": flagBandingBinding} {
		if p, err := svc.proteksiBanding(ctx, "UJI-NB-1", tahan); err != nil || !dalam(p, "Premi on paymnet can't be 0") {
			t.Errorf("%s: %v %q", nama, err, p)
		}
	}
}

// TestSubmitBandingMarketing - tahap marketing Email Type 4: pesan penahan -> 422 tanpa tulisan (hanya pesan penahan);
// hanya pesan bukan penahan (MOID / Scoring) -> lolos; periksa memuat pesan penahan yang sama. Accept tetap 16 pesan.
func TestSubmitBandingMarketing(t *testing.T) {
	ctx := context.Background()
	banding := models.IsiSubmit{EmailType: "4", Comment: "UJI", BandingTo: "JUW_B"}
	st := submitUji()
	svc := layananRusak57b(st)
	_, err := svc.AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", banding)
	var g GalatProteksi
	if !errors.Is(err, ErrProteksiSubmit) || !errors.As(err, &g) || len(g.Pesan) != 10 || g.Pesan[0] != pesanSpreading1Rusak[0] || len(st.tulis) != 0 {
		t.Errorf("rusak: %v %q (%d tulis)", err, g.Pesan, len(st.tulis))
	}
	h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "4"})
	if err != nil || !reflect.DeepEqual(h.Pesan, g.Pesan) {
		t.Errorf("periksa rusak: %v %q", err, h.Pesan)
	}
	if h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "1"}); err != nil || len(h.Pesan) != 16 {
		t.Errorf("periksa accept: %v %d", err, len(h.Pesan))
	}
	// Hanya MOID + Scoring kosong (bukan penahan): Banding lolos, Accept 422.
	st = submitUji()
	svc = layananSubmit(st)
	k := svc.kasus.(kasusTiruan).ada["UJI-NB-1"]
	k.General.MarketingID, k.RiskScoring.FinalScore = "", nil
	if h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "4"}); err != nil || len(h.Pesan) != 0 {
		t.Errorf("periksa bukan penahan: %v %q", err, h.Pesan)
	}
	if _, err := svc.AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "UJI"}); !errors.Is(err, ErrProteksiSubmit) {
		t.Errorf("accept bukan penahan: %v", err)
	}
	if _, err := svc.AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", banding); err != nil || len(st.tulis) != 1 {
		t.Errorf("banding bukan penahan: %v (%d tulis)", err, len(st.tulis))
	}
}

// TestSubmitBandingBinding - tahap Binding Email Type 4: Province kosong (CARI1) -> 422; coverage kosong (ProtectMarine.CARI2,
// penahan EmailSection tetapi BUKAN EmailSectionCeding) -> lolos; periksa sama.
func TestSubmitBandingBinding(t *testing.T) {
	ctx := context.Background()
	banding := models.IsiSubmit{EmailType: "4", Comment: "UJI", BandingTo: "JUW_B"}
	for nama, c := range map[string]struct {
		ubah  func(*models.ObjekFire)
		tahan bool
	}{
		"province": {func(o *models.ObjekFire) { o.Province = "" }, true},
		"coverage": {func(o *models.ObjekFire) { o.Items[0].Coverages[0].CoverageNote = "" }, false},
	} {
		st := submitUji()
		svc := layananBinding(st, "1", nil)
		o := lokasiSah()
		c.ubah(&o)
		svc.objek.(objekTiruan).ada["UJI-NB-1"] = []models.ObjekFire{o}
		if a, _ := svc.ProteksiAccept(ctx, "UJI-NB-1", nil); len(a) != 1 {
			t.Fatalf("%s instrumen: %q", nama, a)
		}
		h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "4"})
		if err != nil || (len(h.Pesan) == 1) != c.tahan {
			t.Errorf("%s periksa: %v %q", nama, err, h.Pesan)
		}
		_, err = svc.AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", banding)
		if errors.Is(err, ErrProteksiSubmit) != c.tahan || (len(st.tulis) == 0) != c.tahan {
			t.Errorf("%s submit: %v (%d tulis)", nama, err, len(st.tulis))
		}
	}
}
