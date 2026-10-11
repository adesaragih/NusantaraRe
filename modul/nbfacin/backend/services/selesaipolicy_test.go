package services

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// Tiket 69a-2 - selesai Policy (SaveJsonPolicyFacIn_Act jalur NB FIRE).

// teksPolis - desimal -> teks (nil = "<nil>").
func teksPolis(x *apd.Decimal) string {
	if x == nil {
		return "<nil>"
	}
	return utils.FormatDecimal(x)
}

// TestSelesaiPolicyTulisan - TL Accept IsFlagUW 0: POSITION / NB_STATUS tetap, Resolved-Completed; FACINPRODUCTION satu baris per
// spreading (FireNB langkah 4-13), CEDING_FACINPRODUCTION per cedant, M_RISK_LOSS_PROFILE, DATA_JSON minimal.
func TestSelesaiPolicyTulisan(t *testing.T) {
	ctx := context.Background()
	sp := spreadingPolicy("1000", "10001", "QS UJI")
	sp.Lokasi[0].Items[0].Coverages[0].Spreading = append(sp.Lokasi[0].Items[0].Coverages[0].Spreading,
		models.BarisSpreading{TreatyType: "10003", TreatyName: "QS UJI 2", SharePercentage: d("10"), TSISpreaded: d("50"),
			PremiumSpreaded: d("3"), ClaimEstimation: d("7"), TSIGrossSpreaded: d("400")})
	st := submitUji()
	svc := layananPolicy(st, posisiTeamLeader, sp, 0)
	h, err := svc.AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "UJI"})
	if err != nil || len(st.tulis) != 1 {
		t.Fatalf("%v", err)
	}
	// nomorPolis = nomor yang diisi repository di transaksi (tiruan: "UJI-F" + OldId + …).
	if h.NomorPolis != "UJI-F"+oldIDUji+".10.2026.00001" || h.Position != posisiTeamLeader {
		t.Errorf("hasil %+v", h)
	}
	w := st.tulis[0]
	p := w.Policy
	if p == nil || w.PosisiBaru != posisiTeamLeader || w.NBStatus != "" || !w.NBStatusTetap || w.StatusWorkBaru != "Resolved-Completed" ||
		w.FlagBaru != "" || w.Binding != nil || w.Riwayat.Posisi != "Policy" {
		t.Fatalf("tulisan: %+v", w)
	}
	if p.BusinessOldID != oldIDUji || p.Username != "UJI-TL" || !p.Sekarang.Equal(svc.sekarang()) ||
		p.MulaiPolis.UTC().Format(BentukWaktuPega) != "20261019T170000.000 GMT" {
		t.Errorf("dasar nomor: %+v", p)
	}
	if len(p.Produksi) != 2 {
		t.Fatalf("produksi %d", len(p.Produksi))
	}
	b := p.Produksi[0]
	teks := []string{b.IDPega, b.BusinessCode, b.CedingCo, b.SOB, b.BeginDate, b.EndDate, b.Province, b.Kabupaten, b.Kecamatan,
		b.ObjectItem, b.CurrID, b.ObjectNameID, b.IdxObjectItem, b.Coverage, b.CoverageNote, b.AccumulationCode, b.OccupationCode,
		b.MarketingCode, b.InsuredName, b.RISlipCeding, b.GroupType, b.QQName, b.LayerNo, b.PolicyType, b.CedingID, b.SOBID,
		b.InsuredID, b.StatusBusiness, b.GroupName, b.TopRisk, b.JnReas, b.RoadName, b.BuildingNo}
	mau := []string{"UJI-NB-1", "", "UJI CEDING A;UJI CEDING B", "UJI SOB", "20/10/2026 00:00:00",
		"20/10/2027 00:00:00", "UJI PROV", "UJI KOTA", "UJI KEC", "UJI MESIN", "IDR", "UJI01", "1", "FLEXAS", "UJI", "UJI-AKU",
		"2001", "UJI-MC", "UJI TERTANGGUNG", "UJI SLIP", "2", "UJI QQ", "NULL", "2", "C1;C2", "UJI-SOB", "UJI-ORG", "1", "UJI GRUP",
		"false", "10001", "", "NO. "}
	if !reflect.DeepEqual(teks, mau) {
		t.Errorf("teks baris:\n%q\n%q", teks, mau)
	}
	f := p.Produksi[1]
	angka := []string{teksPolis(f.PercentRNM), teksPolis(f.PercentShareSpreaded), teksPolis(f.TSI), teksPolis(f.Premi), teksPolis(f.TSI100),
		teksPolis(f.LoL), teksPolis(f.PercentRIComm), teksPolis(f.PctRICommSelisih), teksPolis(f.RIComm), teksPolis(f.PctBrokerageFee),
		teksPolis(f.BrokerageFee), f.JnReas}
	// Premi 3 × RICommision 10 % = 0.3; × % brokerage 5 = 0.15; ClaimEstimation 7 -> LoL.
	if !reflect.DeepEqual(angka, []string{"12.5", "10", "50", "3", "400", "7", "10", "10", "0.3", "5", "0.15", "10003"}) {
		t.Errorf("angka baris kedua: %q", angka)
	}
	if teksPolis(b.LoL) != teksPolis(simpan8(sp.Lokasi[0].Items[0].Coverages[0].Spreading[0].ClaimEstimation)) {
		t.Errorf("LoL baris 1: %s", teksPolis(b.LoL))
	}
	if len(p.Ceding) != 2 || p.Ceding[0].CedingID != "C1" || p.Ceding[0].CedingName != "UJI CEDING A" || teksPolis(p.Ceding[0].Share) != "60" ||
		p.Ceding[1].CedingID != "C2" || teksPolis(p.Ceding[1].Share) != "40" {
		t.Errorf("ceding %+v", p.Ceding)
	}
	// Baris 1 OR (10001, PremiumQs belum terisi = 0) lalu baris 2 QS 3 × 0.8.
	if p.RisikoLoss.COB != "UJI COB" || teksPolis(p.RisikoLoss.PremiumQS) != "2.4" || p.RisikoLoss.PremiumFacout != nil ||
		teksPolis(p.RisikoLoss.PremiumOR) != teksPolis(simpan8(sp.Lokasi[0].Items[0].Coverages[0].Spreading[0].PremiumSpreaded)) {
		t.Errorf("risiko %+v", p.RisikoLoss)
	}
	// DATA_JSON minimal - work owner 09-10-2026 "sistem baru tidak pakai data_json lgi" (isi polis tidak disalin).
	var doc map[string]any
	if err := json.Unmarshal([]byte(p.JSONData), &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc) != 1 || doc["IDNewBisnis"] != "UJI-NB-1" || p.JSONData != `{"IDNewBisnis":"UJI-NB-1"}` {
		t.Errorf("doc %v", doc)
	}
}

// TestSelesaiPolicyDitolak - BusinessOldId tidak diketahui -> 409 tanpa tulisan; objek / spreading tidak sejalan -> galat.
func TestSelesaiPolicyDitolak(t *testing.T) {
	ctx := context.Background()
	qs := spreadingPolicy("1000", "10001", "QS UJI")
	accept := models.IsiSubmit{EmailType: "1", Comment: "UJI"}
	for nama, old := range map[string][]string{"0 baris": nil, "ganda": {"A", "B"}, "kosong": {" "}} {
		st := submitUji()
		svc := denganOldID(layananPolicy(st, posisiTeamLeader, qs, 0), old)
		if _, err := svc.AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", accept); !errors.Is(err, ErrBusinessOldIDTidakAda) || len(st.tulis) != 0 {
			t.Errorf("%s: %v", nama, err)
		}
	}
	st := submitUji()
	sp := spreadingPolicy("1000", "10001", "QS UJI")
	sp.Lokasi[0].Items[0].Coverages[0].OldID = "LAIN"
	if _, err := layananPolicy(st, posisiTeamLeader, sp, 0).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", accept); !errors.Is(err, ErrBindingTidakSejalan) ||
		len(st.tulis) != 0 {
		t.Errorf("tidak sejalan: %v", err)
	}
}

// TestKasusSelesaiBacaSaja - STATUS_WORK Resolved-*: POSITION TL / UW flag 1 tidak lagi membuka tahap (Submit 409, tab baca-saja,
// spreading baca-saja).
func TestKasusSelesaiBacaSaja(t *testing.T) {
	ctx := context.Background()
	for _, pos := range []string{posisiTeamLeader, posisiUWUtama, posisiSeniorUW} {
		tk := models.TahapKasus{Position: pos, FlagOnGoingPolicy: "1", StatusWork: "Resolved-Completed"}
		k := models.KonteksSubmit{Position: pos, FlagOnGoingPolicy: "1", StatusWork: "Resolved-Completed"}
		alasan, err := alasanTidakBoleh(k, pelakuPosisi(pos))
		if tahapDari(tk) != "" || tahapKonteks(k) != "" || tahapUW(tk) || !errors.Is(err, ErrKasusBacaSaja) ||
			alasan != "Case baca-saja: case sudah selesai (Resolved-Completed)." {
			t.Errorf("%s: %q %v", pos, alasan, err)
		}
		svc := Baru(nil).DenganTahap(tahapTiruan{"K": tk})
		if boleh, _ := svc.AksesUbahKasus(ctx, pelakuPosisi(pos), "K"); boleh {
			t.Errorf("%s boleh ubah", pos)
		}
		if a, b := svc.AksesSpreading(ctx, pelakuPosisi(pos), "K"); a || b {
			t.Errorf("%s spreading", pos)
		}
	}
	// Instrumen: tanpa Resolved tahapnya terbuka.
	if tahapKonteks(models.KonteksSubmit{Position: posisiSeniorUW, FlagOnGoingPolicy: "1", StatusWork: "Pending-Policy"}) != TahapUWPolicy ||
		tahapKonteks(models.KonteksSubmit{Position: posisiTeamLeader, FlagOnGoingPolicy: "1"}) != TahapTeamLeaderPolicy {
		t.Error("instrumen")
	}
	st := submitUji()
	svc := layananPolicy(st, posisiTeamLeader, spreadingPolicy("1000", "10001", "QS UJI"), 0)
	k := st.konteks["UJI-NB-1"]
	k.StatusWork = "Resolved-Completed"
	st.konteks["UJI-NB-1"] = k
	if _, err := svc.AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "UJI"}); !errors.Is(err, ErrKasusBacaSaja) ||
		len(st.tulis) != 0 {
		t.Errorf("submit ulang: %v", err)
	}
}

// TestBarisCedingProduksi - InsertCedingProduction 2.1-2.5 (berurutan, langkah sesudahnya menimpa).
func TestBarisCedingProduksi(t *testing.T) {
	dua := []models.BarisCedant{{CedingCo: "A;", CedingCoName: "N;A", ShareCeding: d("40")}, {CedingCo: "B", CedingCoName: "NB"}}
	satu := []models.BarisCedant{{CedingCo: "A", ShareCeding: d("40")}}
	tiga := []models.BarisCedant{{CedingCo: "A"}, {CedingCo: "B"}, {CedingCo: "C"}}
	for nama, c := range map[string]struct {
		k    models.KasusCedant
		ps   string
		mau  []string
		kode string
	}{
		"gross":           {models.KasusCedant{ShareCedantType: "0", Cedant: dua}, "12.5", []string{"320", "<nil>"}, "A"},
		"gross PS kosong": {models.KasusCedant{ShareCedantType: "0", Cedant: satu}, "", []string{"<nil>"}, "A"},
		"gross 1/3":       {models.KasusCedant{ShareCedantType: "0", Cedant: []models.BarisCedant{{CedingCo: "A", ShareCeding: d("10")}}}, "30", []string{"33.33333333"}, "A"},
		"net":             {models.KasusCedant{ShareCedantType: "1", Cedant: dua}, "12.5", []string{"40", "<nil>"}, "A"},
		"kosong":          {models.KasusCedant{ShareCedantType: "", Cedant: dua}, "12.5", []string{"40", "50"}, "A"},
		"kosong satu":     {models.KasusCedant{ShareCedantType: "", Cedant: satu}, "12.5", []string{"100"}, "A"},
		"kosong tiga":     {models.KasusCedant{ShareCedantType: "", Cedant: tiga}, "", []string{"33.3333", "33.3333", "33.3333"}, "A"},
		"tipe lain":       {models.KasusCedant{ShareCedantType: "2", Cedant: satu}, "12.5", []string{"<nil>"}, "A"},
	} {
		var ps = d("0")
		if c.ps != "" {
			ps = d(c.ps)
		}
		h := barisCedingProduksi(c.k, ps)
		got := []string{}
		for _, x := range h {
			got = append(got, teksPolis(x.Share))
		}
		if !reflect.DeepEqual(got, c.mau) || h[0].CedingID != c.kode {
			t.Errorf("%s: %q %+v", nama, got, h)
		}
	}
	if h := barisCedingProduksi(models.KasusCedant{ShareCedantType: "1", Cedant: dua}, d("1")); h[0].CedingName != "NA" {
		t.Errorf("nama tanpa ';': %q", h[0].CedingName)
	}
}

// TestProfilRisikoLoss - SaveJsonPolicy 19-20: Local ditimpa per iterasi, PremiumQs dikosongkan per coverage.
func TestProfilRisikoLoss(t *testing.T) {
	cov := func(tsi, premi string, b ...models.BarisSpreading) models.CoverageSpreading {
		return models.CoverageSpreading{TSINusantaraRe: d(tsi), PremiNusantaraRe: d(premi), Spreading: b}
	}
	sb := func(tipe, premi string) models.BarisSpreading {
		return models.BarisSpreading{TreatyType: tipe, PremiumSpreaded: d(premi)}
	}
	sp := models.KasusSpreading{Lokasi: []models.LokasiSpreading{{Items: []models.ItemSpreading{{Coverages: []models.CoverageSpreading{
		cov("100", "10", sb("10003", "100"), sb("10001", "50")),
		cov("200", "20", sb("10007", "30"), sb("10015", "20")),
	}}}}}}
	p := profilRisikoLoss(sp, "COB")
	got := []string{teksPolis(p.TSINusantaraRe), teksPolis(p.PremiNusantaraRe), teksPolis(p.PremiumOR), teksPolis(p.PremiumQS), teksPolis(p.PremiumFacout), p.COB}
	if !reflect.DeepEqual(got, []string{"200", "20", "30", "80", "20", "COB"}) {
		t.Errorf("dua coverage: %q", got)
	}
	sp.Lokasi[0].Items[0].Coverages = sp.Lokasi[0].Items[0].Coverages[:1]
	p = profilRisikoLoss(sp, "")
	if teksPolis(p.PremiumOR) != "70" || p.PremiumFacout != nil {
		t.Errorf("QS sebelum OR: %s %s", teksPolis(p.PremiumOR), teksPolis(p.PremiumFacout))
	}
	if p := profilRisikoLoss(models.KasusSpreading{}, "X"); p.TSINusantaraRe != nil || p.PremiumOR != nil {
		t.Errorf("kosong %+v", p)
	}
}

// TestKodeOkupasiProduksi - FireNB 13.2: yang baru di depan, pemisah ";".
func TestKodeOkupasiProduksi(t *testing.T) {
	if got := kodeOkupasiProduksi([]models.OkupasiObjek{{OccupationID: "A"}, {OccupationID: "B"}, {OccupationID: "C"}}); got != "C;B;A" {
		t.Errorf("%q", got)
	}
	if got := kodeOkupasiProduksi(nil); got != "" {
		t.Errorf("%q", got)
	}
}

// TestBarisFacInProductionBasisLayer - coverage basis "5" tanpa LayerList -> tanpa baris (coverage lain tetap).
func TestBarisFacInProductionBasisLayer(t *testing.T) {
	o := lokasiSah()
	o.Items[0].Coverages = append(o.Items[0].Coverages, models.CoverageObjek{Coverage: "EQ", OldID: "EQ", CoverageBasis: "5"})
	sp := spreadingPolicy("1000", "10001", "QS UJI")
	sp.Lokasi[0].Items[0].Coverages[0].Spreading[0].ClaimEstimation = nil
	c2 := sp.Lokasi[0].Items[0].Coverages[0]
	c2.OldID = "EQ"
	sp.Lokasi[0].Items[0].Coverages = append(sp.Lokasi[0].Items[0].Coverages, c2)
	h, err := barisFacInProduction(dasarProduksi{id: "K"}, []models.ObjekFire{o}, sp)
	if err != nil || len(h) != 1 || h[0].Coverage != "FLEXAS" {
		t.Errorf("%v %+v", err, h)
	}
	if utils.FormatDecimal(h[0].RIComm) != "0" || h[0].PercentRIComm.String() != "0" || h[0].PctBrokerageFee != nil || teksPolis(h[0].LoL) != "0" ||
		utils.FormatDecimal(h[0].BrokerageFee) != "0" {
		t.Errorf("RICommision / brokerage kosong: %+v", h[0])
	}
}
