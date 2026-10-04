package acceptance

import (
	"errors"
	"strings"
	"testing"
)

func harusPanic(t *testing.T, nama string, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: tidak panic", nama)
		}
	}()
	f()
}

// TestDekodeDomainK021 - enam nilai sah dan hasil dekoder IsUWAccepted.
// Harapan dari K-021 dan DecisionTable IsUWAccepted (bukan dari kode).
func TestDekodeDomainK021(t *testing.T) {
	for _, u := range []struct {
		nilai string
		mau   Keputusan
		hasil string
	}{
		{"1", Accept, "confirm"}, {"2", Reject, "reject"}, {"3", Ask, "ask"},
		{"4", Banding, "banding"}, {"7", Decline, "decline"}, {"9", Revise, "revise"},
	} {
		k, err := PeriksaDomain(ParalelRun, u.nilai, nil)
		if err != nil || k != u.mau {
			t.Errorf("%s: %v %v, mau %v", u.nilai, k, err, u.mau)
		}
		if h := k.HasilIsUWAccepted(); h != u.hasil {
			t.Errorf("%s: hasil IsUWAccepted %q, mau %q", u.nilai, h, u.hasil)
		}
	}
}

// TestRejectBukanDecline - K-014: Reject dapat dibanding, Decline final. Efek
// samping mematikan konfirmasi binding + penerimaan R/I slip dimiliki Reject DAN
// Revise (temuan: SetReviseProposal L179/L209), tidak Decline.
func TestRejectBukanDecline(t *testing.T) {
	if Reject == Decline {
		t.Fatal("Reject dan Decline disatukan")
	}
	for k, mau := range map[Keputusan]bool{Reject: true, Revise: true, Decline: false, Accept: false, Ask: false, Banding: false} {
		if got := k.MatikanBindingDanRISlip(); got != mau {
			t.Errorf("%v: efek samping %v, mau %v", k, got, mau)
		}
	}
}

// TestFlagFase - satu flag, tiga situasi (K-008, tiket 13).
func TestFlagFase(t *testing.T) {
	var catatan []string
	catat := func(s string) { catatan = append(catatan, s) }

	// (a) dikenali, baris belum terverifikasi.
	harusPanic(t, "(a) paralel run", func() { JalurBelumTerverifikasi(ParalelRun, Accept, "uji", catat) })
	if k := JalurBelumTerverifikasi(Produksi, Accept, "uji", catat); k != Decline {
		t.Errorf("(a) produksi: %v, mau Decline", k)
	}
	if len(catatan) != 1 || !strings.Contains(catatan[0], "uji") {
		t.Errorf("(a) produksi tanpa catatan: %v", catatan)
	}

	// (b) tidak dikenali sama sekali - kedua fase.
	harusPanic(t, "(b) paralel run", func() { KondisiTakDikenali(ParalelRun, "uji") })
	harusPanic(t, "(b) produksi", func() { KondisiTakDikenali(Produksi, "uji") })

	// Di luar domain: selalu bercatatan; paralel run panic, produksi galat.
	catatan = nil
	for _, v := range []string{"5", "6", "8", "", "01", "Accept"} {
		harusPanic(t, "luar domain paralel "+v, func() { _, _ = PeriksaDomain(ParalelRun, v, catat) })
		if _, err := PeriksaDomain(Produksi, v, catat); !errors.Is(err, ErrDiLuarDomain) {
			t.Errorf("luar domain produksi %q: %v", v, err)
		}
	}
	if len(catatan) != 12 {
		t.Errorf("catatan luar domain %d, mau 12 (selalu dicatat)", len(catatan))
	}
}

// TestFaseNolDitolak - Fase nol (tidak diisi) bukan Produksi diam-diam.
func TestFaseNolDitolak(t *testing.T) {
	harusPanic(t, "PeriksaDomain fase nol", func() { _, _ = PeriksaDomain(Fase(0), "5", nil) })
	harusPanic(t, "JalurBelumTerverifikasi fase nol", func() { JalurBelumTerverifikasi(Fase(0), Accept, "uji", nil) })
}
