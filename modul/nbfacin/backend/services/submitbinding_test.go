package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// lokasiBinding - lokasiSah + satu deductible terisi per coverage (tiket 64b: deductible 9.2.30.6.6 / 13 berjalan di flag 2).
func lokasiBinding() models.ObjekFire {
	o := lokasiSah()
	for i := range o.Items {
		for j := range o.Items[i].Coverages {
			o.Items[i].Coverages[j].Deductibles = []models.Deductible{{TypeDeductible: "1", Condition: "1", MinMax: "1", Amount: d("1000")}}
		}
	}
	return o
}

// layananBinding - kasus UJI-NB-1 di tahap MARKETING (BINDING): POSITION ReasFacInMarketing, flag 2, ConfirmBinding `confirm`;
// General ber-Ceding Co dan SOB (proteksi langkah 33 / 34 lolos); deductible terisi (lokasiBinding).
func layananBinding(st *submitTiruan, confirm string, ubahKasus func(*models.Kasus)) *Service {
	k := st.konteks["UJI-NB-1"]
	k.Position, k.FlagOnGoingPolicy, k.ConfirmBinding, k.ReceivedRiSlip = repository.WorkbasketMarketing, "2", confirm, "false"
	st.konteks["UJI-NB-1"] = k
	svc := layananSubmit(st, lokasiBinding())
	kasus := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{InsuredID: "UJI-ORG"},
		General:     models.General{MarketingID: "UJI-MO", CedingCoName: "UJI CEDING", SourceOfBusinessID: "UJI-SOB"},
		RiskScoring: models.RingkasanScoring{FinalScore: d("61")}}
	if ubahKasus != nil {
		ubahKasus(kasus)
	}
	// K73-29: CheckSpreadingProtect berjalan di tahap Binding - spreading sah (satu baris ORS 100 % per coverage, premi 0 = premi
	// RNM kosong) agar uji tiket 64 tetap menguji proteksinya sendiri.
	sp := kasusDuaItem("50", "USD")
	setiapCoverage(&sp, func(_ *models.ItemSpreading, c *models.CoverageSpreading) {
		c.Spreading = []models.BarisSpreading{{TreatyType: treatyORS, TreatyName: namaORS, SharePercentage: d("100")}}
	})
	return svc.DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": kasus}}).
		DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}})
}

// TestSubmitBindingMatriks - tiket 64: Accept (TL / UW Financial bonding, flag 1, ConfirmBinding + Received RI Slip,
// FACINPERFORMANCE BINDING/ACCEPT sekali), Banding (SetBanding, flag 0), Revise (Offer flag 0, reset ConfirmBinding), Decline
// (End3, BINDING/DECLINE, NB "WAS DECLINED BY"); riwayat PUTARAN "2", POSISI tahap tujuan, APPROVAL pilihan.
func TestSubmitBindingMatriks(t *testing.T) {
	ctx := context.Background()
	decline := models.IsiSubmit{EmailType: "7", SebabNonTc: "Accumulation / Capacity"}
	for nama, c := range map[string]struct {
		isi                         models.IsiSubmit
		grup                        KonfigGrup
		teamGroup                   string
		mauPosisi, mauNB, mauFlag   string
		mauStatus                   *models.StatusBinding
		mauRiwayatPosisi, mauPega   string
		mauPerf, mauRISlip, mauWork string
	}{
		"Accept TL": {models.IsiSubmit{EmailType: "1", Comment: "UJI", ReceivedRiSlip: true}, KonfigGrup{}, "1", posisiTeamLeader,
			"NB IS IN TEAM LEADER'S INBOX", "1", &models.StatusBinding{ConfirmBinding: "1", ReceivedRiSlip: "true"}, "Policy", "ACCEPT",
			"ACCEPT", "true", ""},
		"Accept bonding team 5": {models.IsiSubmit{EmailType: "1", Comment: "UJI"}, KonfigGrup{}, "5", posisiUWFinancial,
			"NB IS IN UNDERWRITER FINANCIAL'S INBOX", "1", &models.StatusBinding{ConfirmBinding: "1", ReceivedRiSlip: "false"}, "Policy",
			"ACCEPT", "ACCEPT", "false", ""},
		"Accept bonding nama": {models.IsiSubmit{EmailType: "1", Comment: "UJI"}, KonfigGrup{BondingMarketingName: []string{"uji - pemasar"}},
			"1", posisiUWFinancial, "NB IS IN UNDERWRITER FINANCIAL'S INBOX", "1",
			&models.StatusBinding{ConfirmBinding: "1", ReceivedRiSlip: "false"}, "Policy", "ACCEPT", "ACCEPT", "false", ""},
		"Banding": {models.IsiSubmit{EmailType: "4", Comment: "UJI", BandingTo: "JUW_B"}, KonfigGrup{}, "1", posisiJuniorUWB,
			"NB IS IN JUNIOR UNDERWRITER (B)'S INBOX", "0", nil, "Offer", "BANDING", "", "", ""},
		"Revise": {models.IsiSubmit{EmailType: "9", Comment: "UJI", ReceivedRiSlip: true}, KonfigGrup{}, "1", repository.WorkbasketMarketing,
			"NB IS IN Uji - Pemasar'S INBOX", "0", &models.StatusBinding{ConfirmBinding: "0", ReceivedRiSlip: "false"}, "Offer", "REVISE",
			"", "", ""},
		"Decline": {decline, KonfigGrup{}, "1", posisiDecline, "NB WAS DECLINED BY  UJI PENGGUNA", "", nil, "Binding", "DECLINE",
			"DECLINE", "false", statusDecline},
	} {
		st := submitUji()
		k := st.konteks["UJI-NB-1"]
		k.TeamGroup = c.teamGroup
		st.konteks["UJI-NB-1"] = k
		svc := layananBinding(st, "1", nil).DenganGrup(c.grup, false)
		h, err := svc.AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", c.isi)
		if err != nil || len(st.tulis) != 1 {
			t.Fatalf("%s: %v (%d tulis)", nama, err, len(st.tulis))
		}
		w := st.tulis[0]
		if h.Position != c.mauPosisi || h.NBStatus != c.mauNB || w.FlagBaru != c.mauFlag || !reflect.DeepEqual(w.StatusBinding, c.mauStatus) ||
			w.Riwayat.Posisi != c.mauRiwayatPosisi || w.Riwayat.Putaran != "2" || w.Pega.Status != c.mauPega ||
			w.Pega.Workbasket != repository.WorkbasketMarketing || w.StatusWorkBaru != c.mauWork || w.PerformanceLanjut != nil ||
			w.Binding != nil {
			t.Errorf("%s: %+v / %+v", nama, h, w)
		}
		if c.mauPerf == "" {
			if w.Performance != nil {
				t.Errorf("%s: performance %+v", nama, w.Performance)
			}
		} else if f := w.Performance; f == nil || f.PositionDoc != "BINDING" || f.StatusDoc != c.mauPerf || f.RISlipReceived != c.mauRISlip {
			t.Errorf("%s: performance %+v", nama, f)
		}
		if w.Riwayat.Approval != labelApproval[c.isi.EmailType].riwayat {
			t.Errorf("%s: approval %q", nama, w.Riwayat.Approval)
		}
	}
}

// TestSubmitBindingDitolak - Accept tanpa Binding -> 422 (pesan sistem); proteksi langkah 33 / 34; Email Type 2 / 3 -> 400;
// workbasket bukan marketing -> 403; tanpa tulisan.
func TestSubmitBindingDitolak(t *testing.T) {
	ctx := context.Background()
	terima := models.IsiSubmit{EmailType: "1", Comment: "UJI"}
	for nama, c := range map[string]struct {
		confirm  string
		ubah     func(*models.Kasus)
		pelaku   inti.Pelaku
		isi      models.IsiSubmit
		mau      error
		mauPesan []string
	}{
		"belum binding": {"", nil, pelakuMkt, terima, ErrProteksiSubmit, []string{pesanBelumBinding}},
		"ceding / sob kosong": {"1", func(k *models.Kasus) { k.General.CedingCoName, k.General.SourceOfBusinessID = "", "" }, pelakuMkt,
			terima, ErrProteksiSubmit, []string{pesanCedingCoKosong, pesanSOBKosong}},
		"reject":     {"1", nil, pelakuMkt, models.IsiSubmit{EmailType: "2", Comment: "UJI"}, ErrMasukanSubmit, nil},
		"ask":        {"1", nil, pelakuMkt, models.IsiSubmit{EmailType: "3", Comment: "UJI"}, ErrMasukanSubmit, nil},
		"workbasket": {"1", nil, inti.Pelaku{AkunID: "UJI-TL", Peran: []string{posisiTeamLeader}}, terima, inti.ErrTanpaWewenang, nil},
	} {
		st := submitUji()
		_, err := layananBinding(st, c.confirm, c.ubah).AjukanSubmit(ctx, c.pelaku, "UJI-NB-1", c.isi)
		var g GalatProteksi
		if !errors.Is(err, c.mau) || len(st.tulis) != 0 || (c.mauPesan != nil && (!errors.As(err, &g) || !reflect.DeepEqual(g.Pesan, c.mauPesan))) {
			t.Errorf("%s: %v (%d tulis)", nama, err, len(st.tulis))
		}
	}
	// Scoring kosong: tidak menahan Accept (tiket 64b: penahan Accept Binding = ProtectSpreading.CARI1 / CARI23 /
	// ProtectZipCode.CARI1 saja) maupun Banding (tiket 57b); Revise tidak diproteksi.
	for e, mau := range map[string]bool{"1": false, "4": false, "9": false} {
		st := submitUji()
		svc := layananBinding(st, "1", func(k *models.Kasus) { k.RiskScoring = models.RingkasanScoring{} })
		_, err := svc.AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: e, Comment: "UJI", BandingTo: "JUW_B"})
		if errors.Is(err, ErrProteksiSubmit) != mau {
			t.Errorf("proteksi %s: %v", e, err)
		}
	}
	// Periksa: Accept tanpa Binding memuat pesan sistem; Banding diproteksi (lolos); Revise tanpa proteksi.
	st := submitUji()
	svc := layananBinding(st, "", nil)
	if h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "1"}); err != nil || !reflect.DeepEqual(h.Pesan, []string{pesanBelumBinding}) {
		t.Errorf("periksa accept: %v %+v", err, h)
	}
	for _, e := range []string{"4", "9", "7"} {
		if h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: e}); err != nil || len(h.Pesan) != 0 {
			t.Errorf("periksa %s: %v %+v", e, err, h)
		}
	}
}

// TestTampilanDanTabBinding - tampilan tahap "binding" (bolehSubmit, ConfirmBinding / ReceivedRiSlip); tab case tahap Binding
// DAPAT DIUBAH sebelum tombol Binding, baca-saja sesudahnya (tiket 64 lanjutan - RALAT "tab case TETAP baca-saja");
// di luar Binding kedua bendera false.
func TestTampilanDanTabBinding(t *testing.T) {
	ctx := context.Background()
	st := submitUji()
	svc := layananBinding(st, "1", nil)
	k := st.konteks["UJI-NB-1"]
	k.ReceivedRiSlip = "true"
	st.konteks["UJI-NB-1"] = k
	tp, err := svc.TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1")
	if err != nil || tp.Tahap != TahapBinding || !tp.BolehSubmit || !tp.ConfirmBinding || !tp.ReceivedRiSlip {
		t.Errorf("%v %+v", err, tp)
	}
	for confirm, mau := range map[string]string{"": TahapBinding, "0": TahapBinding, "1": ""} {
		if got := tahapDari(models.TahapKasus{Position: repository.WorkbasketMarketing, FlagOnGoingPolicy: "2", ConfirmBinding: confirm}); got != mau {
			t.Errorf("tab case tahap Binding ConfirmBinding %q: %q mau %q", confirm, got, mau)
		}
	}
	st2 := submitUji()
	k2 := st2.konteks["UJI-NB-1"]
	k2.Position, k2.ConfirmBinding, k2.ReceivedRiSlip = repository.WorkbasketMarketing, "1", "true"
	st2.konteks["UJI-NB-1"] = k2
	if tp, err := layananSubmit(st2).TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1"); err != nil || tp.Tahap != TahapMarketing || tp.ConfirmBinding || tp.ReceivedRiSlip {
		t.Errorf("marketing: %v %+v", err, tp)
	}
}

// TestTekanBinding - tombol Binding: satu tulisan (posisi + flag penjaga, JSONDATA minimal, riwayat "Offer has been bind."
// PUTARAN 2 POSISI Binding APPROVAL Accept), TampilanSubmit terbaru; sudah Binding / bukan tahap Binding -> 409; proteksi ->
// 422; galat repository dipetakan.
func TestTekanBinding(t *testing.T) {
	ctx := context.Background()
	st := submitUji()
	tp, err := layananBinding(st, "", nil).TekanBinding(ctx, pelakuMkt, "UJI-NB-1")
	if err != nil || len(st.tekan) != 1 || !tp.ConfirmBinding {
		t.Fatalf("%v %+v", err, tp)
	}
	b := st.tekan[0]
	if b.CaseID != "UJI-NB-1" || b.IDPega != "UJI-NB-1" || b.Position != repository.WorkbasketMarketing || b.Flag != "2" ||
		b.JSONData != `{"IDNewBisnis":"UJI-NB-1"}` || b.Riwayat.Keterangan != "Offer has been bind." || b.Riwayat.Posisi != "Binding" ||
		b.Riwayat.Putaran != "2" || b.Riwayat.Approval != "Accept" || b.Riwayat.PIC != "Uji Pengguna" || b.Riwayat.AksesLogin != "UJI-MKT" {
		t.Errorf("%+v", b)
	}
	for nama, c := range map[string]struct {
		confirm, posisi string
		ubah            func(*models.Kasus)
		galatRepo       error
		mau             error
	}{
		"sudah binding":  {"1", "", nil, nil, ErrSudahBinding},
		"bukan binding":  {"", "ReasFacInUnderwriting", nil, nil, ErrKasusBacaSaja},
		"proteksi":       {"", "", func(k *models.Kasus) { k.General.CedingCoName = "" }, nil, ErrProteksiSubmit},
		"repo sudah":     {"", "", nil, repository.ErrSudahBinding, ErrSudahBinding},
		"repo penawaran": {"", "", nil, repository.ErrPenawaranTidakAda, ErrPenawaranBelumAda},
		"repo posisi":    {"", "", nil, repository.ErrPosisiBerubah, ErrKasusBacaSaja},
	} {
		st := submitUji()
		svc := layananBinding(st, c.confirm, c.ubah)
		st.galatTekan = c.galatRepo
		pelaku := pelakuMkt
		if c.posisi != "" { // tahap uw (flag 0): lolos penjaga workbasket, ditolak cek tahap Binding
			k := st.konteks["UJI-NB-1"]
			k.Position, k.FlagOnGoingPolicy = c.posisi, "0"
			st.konteks["UJI-NB-1"] = k
			pelaku = inti.Pelaku{AkunID: "UJI-UW", Peran: []string{c.posisi}}
		}
		if _, err := svc.TekanBinding(ctx, pelaku, "UJI-NB-1"); !errors.Is(err, c.mau) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	if _, err := layananBinding(submitUji(), "", nil).TekanBinding(ctx, inti.Pelaku{AkunID: "UJI-TL", Peran: []string{posisiTeamLeader}},
		"UJI-NB-1"); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("403: %v", err)
	}
}
