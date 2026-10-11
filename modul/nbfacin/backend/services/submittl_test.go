package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

var pelakuTL = inti.Pelaku{AkunID: "UJI-TL", Peran: []string{"ReasFacInTeamLeader"}}

// layananSubmitTL - kasus UJI-NB-1 di posisi Team Leader; bedaMataUang = satu deductible USD pada item IDR (IsBedaCurr).
func layananSubmitTL(st *submitTiruan, bedaMataUang bool) *Service {
	k := st.konteks["UJI-NB-1"]
	k.Position = "ReasFacInTeamLeader"
	st.konteks["UJI-NB-1"] = k
	o := lokasiSah()
	if bedaMataUang {
		o.Items[0].Coverages[0].Deductibles = []models.Deductible{{Currency: "USD"}}
	} else {
		o.Items[0].Coverages[0].Deductibles = []models.Deductible{{Currency: "IDR"}, {Currency: "UNKNOWN"}, {}}
	}
	return layananSubmit(st, o)
}

// TestTampilanSubmitTL - tahap teamLeader, bolehSubmit untuk pemegang ReasFacInTeamLeader (marketing tidak), isian
// Deductible tersimpan, IsBedaCurr (Currency deductible beda dari item; kosong / UNKNOWN / sama bukan beda).
func TestTampilanSubmitTL(t *testing.T) {
	ctx := context.Background()
	st := submitUji()
	svc := layananSubmitTL(st, true)
	k := st.konteks["UJI-NB-1"]
	k.DeductibleAcceptance, k.CommentDeductible = "true", "UJI catatan"
	st.konteks["UJI-NB-1"] = k
	v, err := svc.TampilanSubmit(ctx, pelakuTL, "UJI-NB-1")
	if err != nil || v.Tahap != "teamLeader" || !v.BolehSubmit || !v.DeductibleBedaMataUang || !v.DeductibleAcceptance ||
		v.CommentDeductible != "UJI catatan" {
		t.Fatalf("%v %+v", err, v)
	}
	if v, _ := svc.TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1"); v.BolehSubmit || v.Alasan != "Akun tidak memegang workbasket ReasFacInTeamLeader." {
		t.Errorf("marketing di TL: %+v", v)
	}
	if v, err := layananSubmitTL(submitUji(), false).TampilanSubmit(ctx, pelakuTL, "UJI-NB-1"); err != nil || v.DeductibleBedaMataUang {
		t.Errorf("mata uang sama / kosong / UNKNOWN bukan beda: %v %+v", err, v)
	}
	if v, err := layananSubmit(submitUji()).TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1"); err != nil || v.Tahap != "marketing" || v.DeductibleBedaMataUang {
		t.Errorf("marketing: %v %+v", err, v)
	}
}

// TestProteksiDeductibleTL - Protection_Act 9-12 (verbatim): IsBedaCurr && tidak dicentang -> "Deductible Acceptance cannot be
// empty"; dicentang && Comment kosong -> "Comment Deductible cannot be empty"; marketing tidak diperiksa; pesan TL di depan.
func TestProteksiDeductibleTL(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		beda   bool
		isi    IsiPeriksa
		pelaku inti.Pelaku
		mau    []string
	}{
		{true, IsiPeriksa{EmailType: "1"}, pelakuTL, []string{"Deductible Acceptance cannot be empty"}},
		{true, IsiPeriksa{EmailType: "1", DeductibleAcceptance: true, CommentDeductible: "  "}, pelakuTL, []string{"Comment Deductible cannot be empty"}},
		{false, IsiPeriksa{EmailType: "1", DeductibleAcceptance: true}, pelakuTL, []string{"Comment Deductible cannot be empty"}},
		{true, IsiPeriksa{EmailType: "1", DeductibleAcceptance: true, CommentDeductible: "UJI"}, pelakuTL, []string{}},
		{false, IsiPeriksa{EmailType: "1"}, pelakuTL, []string{}},
		{true, IsiPeriksa{EmailType: "2"}, pelakuTL, []string{}},
		{true, IsiPeriksa{EmailType: "3"}, pelakuTL, []string{}},
	} {
		h, err := layananSubmitTL(submitUji(), c.beda).PeriksaSubmit(ctx, c.pelaku, "UJI-NB-1", c.isi)
		if err != nil || !reflect.DeepEqual(h.Pesan, c.mau) {
			t.Errorf("%+v: %v %q", c.isi, err, h.Pesan)
		}
	}
	// Urutan: langkah 10 sebelum 18 (Scoring).
	svc := layananSubmitTL(submitUji(), true)
	svc.kasus.(kasusTiruan).ada["UJI-NB-1"].RiskScoring.FinalScore = nil
	if h, err := svc.PeriksaSubmit(ctx, pelakuTL, "UJI-NB-1", IsiPeriksa{EmailType: "1"}); err != nil ||
		!reflect.DeepEqual(h.Pesan, []string{"Deductible Acceptance cannot be empty", "Scoring Risk Cannot Empty !!"}) {
		t.Errorf("urutan: %v %q", err, h.Pesan)
	}
	// Email Type per tahap: TL tanpa Banding, marketing tanpa Ask.
	if _, err := layananSubmitTL(submitUji(), false).PeriksaSubmit(ctx, pelakuTL, "UJI-NB-1", IsiPeriksa{EmailType: "4"}); !errors.Is(err, ErrMasukanSubmit) {
		t.Errorf("TL Banding: %v", err)
	}
	if _, err := layananSubmit(submitUji()).PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "3"}); !errors.Is(err, ErrMasukanSubmit) {
		t.Errorf("marketing Ask: %v", err)
	}
}

// TestAjukanSubmitTL - Accept -> ReasFacInJuniorUnderwritingB + IsJuniorUW_B (NB status, DateToUW / SubmitToUW); Reject / Ask ->
// kembali ke ReasFacInMarketing (Transition170, MarketingName apa adanya); Decline -> Decline + Resolved-Decline tanpa
// FACINPERFORMANCE; riwayat Pega ber-workbasket ReasFacInTeamLeader; isian Deductible disimpan setiap Submit TL.
func TestAjukanSubmitTL(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		isi                   models.IsiSubmit
		posisi, status, label string
	}{
		{models.IsiSubmit{EmailType: "1", Comment: "ok", DeductibleAcceptance: true, CommentDeductible: " UJI terima "},
			"ReasFacInJuniorUnderwritingB", "NB IS IN JUNIOR UNDERWRITER (B)'S INBOX", "Accept"},
		{models.IsiSubmit{EmailType: "2", Comment: "tolak"}, "ReasFacInMarketing", "NB IS IN Uji - Pemasar'S INBOX", "Reject"},
		{models.IsiSubmit{EmailType: "3", Comment: "tanya"}, "ReasFacInMarketing", "NB IS IN Uji - Pemasar'S INBOX", "Ask"},
		{models.IsiSubmit{EmailType: "7", SebabNonTc: "Accumulation / Capacity"}, "Decline", "NB WAS DECLINED BY  UJI PENGGUNA", "Decline"},
	} {
		st := submitUji()
		h, err := layananSubmitTL(st, true).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", c.isi)
		if err != nil || h.Position != c.posisi || h.NBStatus != c.status || len(st.tulis) != 1 {
			t.Fatalf("%s: %v %+v", c.isi.EmailType, err, h)
		}
		w := st.tulis[0]
		pega := map[string]string{"Accept": "ACCEPT", "Reject": "REJECT", "Ask": "ASK", "Decline": "DECLINE"}[c.label]
		if w.PosisiLama != "ReasFacInTeamLeader" || w.Riwayat.Approval != c.label || w.Pega != (models.RiwayatPega{Status: pega,
			Username: "Uji Pengguna", Workbasket: "ReasFacInTeamLeader"}) || w.Performance != nil || w.Deductible == nil {
			t.Errorf("%s tulisan: %+v", c.label, w)
		}
		switch c.label {
		case "Accept":
			if *w.Deductible != (models.DeductibleTL{Acceptance: "true", Comment: "UJI terima"}) || w.KeUW != "20261007T150405.000 GMT" || w.StatusWorkBaru != "" {
				t.Errorf("Accept: %+v %q", w.Deductible, w.KeUW)
			}
		case "Decline":
			if w.StatusWorkBaru != "Resolved-Decline" || w.KeUW != "" {
				t.Errorf("Decline: %+v", w)
			}
		default:
			if w.KeUW != "" || w.StatusWorkBaru != "" || *w.Deductible != (models.DeductibleTL{Acceptance: "false"}) {
				t.Errorf("%s: %+v", c.label, w)
			}
		}
	}
	// 422 Deductible Acceptance tanpa tulisan; Banding TL 400; marketing di posisi TL 403.
	st := submitUji()
	if _, err := layananSubmitTL(st, true).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "ok"}); !errors.Is(err, ErrProteksiSubmit) ||
		len(st.tulis) != 0 {
		t.Errorf("422: %v", err)
	}
	if _, err := layananSubmitTL(submitUji(), false).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", models.IsiSubmit{EmailType: "4", Comment: "x", BandingTo: "JUW_B"}); !errors.Is(err, ErrMasukanSubmit) {
		t.Errorf("Banding TL: %v", err)
	}
	if _, err := layananSubmitTL(submitUji(), false).AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"}); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("marketing di TL: %v", err)
	}
}
