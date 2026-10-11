package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// Tiket 70 - Following Policy Number.

// followingTiruan - PenyimpanFollowing dalam memori: polis = "case|nomor".
type followingTiruan struct {
	polis map[string]bool
	cari  []string // "jenis|nilai|offset|batas"
	baris []models.BarisFollowing
	tulis []tulisFollowingUji
}

type tulisFollowingUji struct {
	caseID, noPolis string
	nilai           *models.NilaiFollowing
}

func (f *followingTiruan) CariFollowing(_ context.Context, jenis, nilai string, offset, batas int) ([]models.BarisFollowing, int, error) {
	f.cari = append(f.cari, jenis+"|"+nilai+"|"+itoa(offset)+"|"+itoa(batas))
	return f.baris, 41, nil
}

func (f *followingTiruan) AdaPolisFollowing(_ context.Context, caseID, noPolis string) (bool, error) {
	return f.polis[caseID+"|"+noPolis], nil
}

func (f *followingTiruan) TulisFollowing(_ context.Context, _ *db.Tx, caseID, noPolis string, nilai *models.NilaiFollowing) error {
	f.tulis = append(f.tulis, tulisFollowingUji{caseID, noPolis, nilai})
	return nil
}

func itoa(n int) string { return string(rune('0'+n/10)) + string(rune('0'+n%10)) }

// TestJenisCariFollowing - GetListViewFollowing_act 4.1-4.9 BERURUTAN (langkah sesudahnya menimpa); opsi 3 / lain 400; filter kosong
// tanpa kueri.
func TestJenisCariFollowing(t *testing.T) {
	a := "" // IDPEGA polos dari services; dua bentuk dicocokkan repository (bindFollowing)
	for nama, c := range map[string]struct{ opsi, filter, jenis, nilai string }{
		"nopolis":              {"2", "RNM-F1.10.2026.00001", "nopolis", "RNM-F1.10.2026.00001"},
		"id NB":                {"1", "NB-12", "idpega", a + "NB-12"},
		"id RNW":               {"1", "RNW-3", "idpega", a + "RNW-3"},
		"id lain":              {"1", "X-3", "", ""},
		"id EDM":               {"1", "EDM-7", "edm", a + "EDM-7"},
		"id NB + EDM (4.5)":    {"1", "NB-1EDM-2", "edm", a + "NB-1EDM-2"},
		"RNM-E timpa opsi 1":   {"1", "NB-RNM-E1", "noendors", "NB-RNM-E1"},
		"RNM-E timpa opsi 2":   {"2", "RNM-E.10.2026.1", "noendors", "RNM-E.10.2026.1"},
		"opsi 4 timpa RNM-E":   {"4", "RNM-E1", "reff", "%RNM-E1%"},
		"reff":                 {"4", "SLIP", "reff", "%SLIP%"},
		"insured":              {"5", "PT A", "insured", "%PT A%"},
		"sob":                  {"6", "BRK", "sob", "%BRK%"},
		"qq":                   {"7", "QQ", "qq", "%QQ%"},
		"filter kosong":        {"5", "", "", ""},
		"EDM opsi 2 bukan 4.5": {"2", "EDM-7", "nopolis", "EDM-7"},
	} {
		j, n, err := jenisCariFollowing(c.opsi, c.filter)
		if err != nil || j != c.jenis || n != c.nilai {
			t.Errorf("%s: %q %q %v", nama, j, n, err)
		}
	}
	for _, o := range []string{"3", "", "8", "0"} {
		if _, _, err := jenisCariFollowing(o, "X"); !errors.Is(err, ErrMasukanFollowing) {
			t.Errorf("opsi %q: %v", o, err)
		}
	}
}

// TestCariFollowing - halaman 20 (offset), filter kosong tanpa kueri, 400 / 401 / 503.
func TestCariFollowing(t *testing.T) {
	ctx := context.Background()
	f := &followingTiruan{baris: []models.BarisFollowing{{CaseID: "NB-1"}}}
	svc := Baru(nil).DenganFollowing(f)
	h, err := svc.CariFollowing(ctx, pelakuMkt, "UJI-NB-1", "5", "  PT A ", 3)
	if err != nil || h.Total != 41 || h.Halaman != 3 || h.Ukuran != 20 || len(h.Baris) != 1 ||
		!reflect.DeepEqual(f.cari, []string{"insured|%PT A%|40|20"}) {
		t.Errorf("%+v %v %v", h, err, f.cari)
	}
	if h, err := svc.CariFollowing(ctx, pelakuMkt, "UJI-NB-1", "1", "", 1); err != nil || h.Baris == nil || len(h.Baris) != 0 || len(f.cari) != 1 {
		t.Errorf("kosong: %+v %v", h, err)
	}
	if _, err := svc.CariFollowing(ctx, pelakuMkt, "UJI-NB-1", "3", "X", 1); !errors.Is(err, ErrMasukanFollowing) {
		t.Errorf("opsi 3: %v", err)
	}
	if _, err := svc.CariFollowing(ctx, pelakuMkt, "UJI-NB-1", "2", "X", 0); !errors.Is(err, ErrMasukanFollowing) {
		t.Errorf("halaman 0: %v", err)
	}
	if _, err := svc.CariFollowing(ctx, inti.Pelaku{}, "UJI-NB-1", "2", "X", 1); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	if _, err := Baru(nil).CariFollowing(ctx, pelakuMkt, "UJI-NB-1", "2", "X", 1); !errors.Is(err, ErrFollowingTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
}

// TestSalinNopolisFollowing - FOLLOWING saja; penjaga: ConfirmBinding 1 -> 409 khusus, baca-saja 409, workbasket 403, polis tidak
// ada 404, badan 400; bolehFollowing.
func TestSalinNopolisFollowing(t *testing.T) {
	ctx := context.Background()
	baru := func(tk models.TahapKasus) (*Service, *followingTiruan) {
		f := &followingTiruan{polis: map[string]bool{"NB-9|RNM-F9": true}}
		return Baru(nil).DenganFollowing(f).DenganTransaksi(tanpaTx).DenganTahap(tahapTiruan{"UJI-NB-1": tk}), f
	}
	mkt := models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0"}
	svc, f := baru(mkt)
	if err := svc.SalinNopolisFollowing(ctx, pelakuMkt, "UJI-NB-1", "NB-9", "RNM-F9"); err != nil ||
		!reflect.DeepEqual(f.tulis, []tulisFollowingUji{{"UJI-NB-1", "RNM-F9", nil}}) || !svc.AksesFollowing(ctx, pelakuMkt, "UJI-NB-1") {
		t.Errorf("salin: %v %+v", err, f.tulis)
	}
	for nama, c := range map[string]struct {
		tk            models.TahapKasus
		pelaku        inti.Pelaku
		sumber, polis string
		mau           error
	}{
		"binding ditekan": {models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "2", ConfirmBinding: "1"}, pelakuMkt, "NB-9", "RNM-F9", ErrFollowingSesudahBinding},
		"TL Policy":       {models.TahapKasus{Position: posisiTeamLeader, FlagOnGoingPolicy: "1", ConfirmBinding: "1"}, pelakuTL, "NB-9", "RNM-F9", ErrFollowingSesudahBinding},
		"baca-saja":       {models.TahapKasus{Position: posisiUWUtama, FlagOnGoingPolicy: "0"}, pelakuMkt, "NB-9", "RNM-F9", ErrKasusBacaSaja},
		"selesai":         {models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0", StatusWork: "Resolved-Completed"}, pelakuMkt, "NB-9", "RNM-F9", ErrKasusBacaSaja},
		"workbasket":      {mkt, pelakuTL, "NB-9", "RNM-F9", inti.ErrTanpaWewenang},
		"polis tidak ada": {mkt, pelakuMkt, "NB-9", "RNM-F8", ErrPolisFollowingTidakAda},
		"tanpa nomor":     {mkt, pelakuMkt, "NB-9", " ", ErrMasukanFollowing},
		"sumber tak sah":  {mkt, pelakuMkt, "", "RNM-F9", ErrMasukanFollowing},
	} {
		svc, f := baru(c.tk)
		if err := svc.SalinNopolisFollowing(ctx, c.pelaku, "UJI-NB-1", c.sumber, c.polis); !errors.Is(err, c.mau) || len(f.tulis) != 0 {
			t.Errorf("%s: %v", nama, err)
		}
		if c.mau != ErrPolisFollowingTidakAda && c.mau != ErrMasukanFollowing && svc.AksesFollowing(ctx, c.pelaku, "UJI-NB-1") {
			t.Errorf("%s bolehFollowing", nama)
		}
	}
	if err := Baru(nil).SalinNopolisFollowing(ctx, pelakuMkt, "UJI-NB-1", "NB-9", "RNM-F9"); !errors.Is(err, ErrFollowingTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
}

// TestSalinSemuaFollowing - Copy All: lokasi case sumber menggantikan milik case (gantiObjek), spreading case sumber dipasang,
// FOLLOWING + PercentShare + RICommision sumber; case sumber tanpa tabel -> 404.
func TestSalinSemuaFollowing(t *testing.T) {
	ctx := context.Background()
	sumber := []models.ObjekFire{{ObjectNo: "7", ObjectType: "UJI SUMBER"}}
	ob := objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {{ObjectNo: "1", ObjectType: "UJI LAMA"}}, "NB-9": sumber}}
	spSumber := models.KasusSpreading{PercentShare: d("30"), Lokasi: []models.LokasiSpreading{{ObjectNo: "7"}}}
	src := &spreadingTiruan{baca: []models.KasusSpreading{spSumber}}
	f := &followingTiruan{polis: map[string]bool{"NB-9|RNM-F9": true, "NB-8|RNM-F8": true}}
	dt := &deductionTiruan{ada: models.KasusDeduction{RICommision: d("12.5")}}
	svc := Baru(nil).DenganObjek(ob).DenganSpreading(src).DenganDeduction(dt).DenganFollowing(f).DenganTransaksi(tanpaTx).
		DenganTahap(tahapTiruan{"UJI-NB-1": {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0"}})
	if err := svc.SalinSemuaFollowing(ctx, pelakuMkt, "UJI-NB-1", "NB-9", "RNM-F9"); err != nil {
		t.Fatal(err)
	}
	if got := ob.ada["UJI-NB-1"]; len(got) != 1 || got[0].ObjectNo != "7" || got[0].ObjectType != "UJI SUMBER" {
		t.Errorf("objek %+v", got)
	}
	if len(f.tulis) != 1 || f.tulis[0].noPolis != "RNM-F9" || f.tulis[0].nilai == nil || f.tulis[0].nilai.PercentShare.String() != "30" ||
		f.tulis[0].nilai.RICommision.String() != "12.5" {
		t.Errorf("following %+v", f.tulis)
	}
	// Spreading case sumber dipasang ke coverage case (posisi + oldId sama), dihitung ulang dengan % Share sumber.
	spIsi := spreadingPolicy("1000", "10001", "QS UJI")
	spIsi.PercentShare = d("30")
	target := spreadingPolicy("1000", "10001", "QS UJI")
	target.Lokasi[0].Items[0].Coverages[0].Spreading = nil
	src2 := &spreadingTiruan{baca: []models.KasusSpreading{spIsi, target}}
	f2 := &followingTiruan{polis: map[string]bool{"NB-9|RNM-F9": true}}
	svc2 := Baru(nil).DenganObjek(ob).DenganSpreading(src2).DenganDeduction(dt).DenganFollowing(f2).DenganTransaksi(tanpaTx)
	if err := svc2.SalinSemuaFollowing(ctx, pelakuMkt, "UJI-NB-1", "NB-9", "RNM-F9"); err != nil || len(src2.tulis) != 1 {
		t.Fatalf("spreading: %v %d", err, len(src2.tulis))
	}
	if w := src2.tulis[0]; w.PercentShare.String() != "30" || len(w.Lokasi[0].Items[0].Coverages[0].Spreading) != 1 ||
		w.Lokasi[0].Items[0].Coverages[0].Spreading[0].TreatyType != "10001" {
		t.Errorf("spreading tertulis %+v", w)
	}
	if err := svc.SalinSemuaFollowing(ctx, pelakuMkt, "UJI-NB-1", "NB-8", "RNM-F8"); !errors.Is(err, ErrCaseSumberTidakAda) {
		t.Errorf("sumber tanpa tabel: %v", err)
	}
	if err := svc.SalinSemuaFollowing(ctx, pelakuMkt, "UJI-NB-1", "NB-9", "RNM-F8"); !errors.Is(err, ErrPolisFollowingTidakAda) {
		t.Errorf("polis: %v", err)
	}
}
