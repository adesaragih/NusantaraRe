package services

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// layananPolicyAneka - layananPolicy (TL Policy, PolicyType 2 = deklarasi) untuk case Aneka BUSINESS (003, old) dengan objek
// Aneka `lokasi` (satu lokasi = baris T_LOCATIONLIST lokasiSah).
func layananPolicyAneka(st *submitTiruan, old string, lokasi []models.LokasiAneka) *Service {
	return layananPolicy(st, posisiTeamLeader, models.KasusSpreading{PercentShare: d("12.5"), StartDateTime: "20261019T170000.000 GMT",
		Lokasi: []models.LokasiSpreading{{ObjectNo: "1"}}}, 0).
		DenganTableOfLimit(tolTiruan{oldID: []string{oldIDUji}, bisnis: satuBisnis("003", old)}).
		DenganObjekAneka(anekaTiruan{ada: map[string][]models.LokasiAneka{"UJI-NB-1": lokasi}})
}

// TestIsFlagUWAneka - K73-37 / K73-41: keluarga Aneka (IsAneka / Golf) = 1 kecuali polis deklarasi (PolicyType 2) = 0; FIRE
// (lini tak terbaca) tetap K69-5.
func TestIsFlagUWAneka(t *testing.T) {
	ctx := context.Background()
	for nama, c := range map[string]struct {
		old, policyType string
		mau             bool
	}{
		"Aneka individual": {"30", "0", true}, "Aneka deklarasi": {"30", "2", false}, "Golf individual": {"27", "0", true},
		"Golf deklarasi": {"27", "2", false}, "Bonding master": {"36", "1", true},
	} {
		svc := layananPolicyAneka(submitUji(), c.old, lokasiFacInOfferAneka())
		svc.kasus.(kasusTiruan).ada["UJI-NB-1"].General.PolicyType = c.policyType
		got, err := svc.isFlagUW(ctx, "UJI-NB-1")
		if err != nil || got != c.mau {
			t.Errorf("%s: %v %v", nama, got, err)
		}
	}
}

// TestBarisFacInProductionAneka - SaveTreatyProduction_Act 19 / 21: kolom beda FIRE (OBJECT_ITEM, CURR_ID / MBD, OCCUPATIONCODE nama
// per okupasi, NOSERIAL, OBJECT_NAME_ID "NULL", MOID case, TOP_RISK / IDX / TSI100 kosong, Golf LoL 0, alamat ASM).
func TestBarisFacInProductionAneka(t *testing.T) {
	st := submitUji()
	k := st.konteks["UJI-NB-1"]
	k.MOID = "MO-UJI"
	dp := dasarProduksi{id: "UJI-NB-1", kodeBisnis: "10048", groupType: "2", k: k}
	objek := []models.ObjekFire{{ObjectNo: "1", ObjectName: "LOKASI UJI", City: "UJI KOTA", District: "UJI KEC", Territory: "UJI RW"}}
	lokasi := lokasiFacInOfferAneka()
	lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading[0].ClaimEstimation = d("7")
	lokasi[0].ObjekGolf[0].Coverages[0].Spreading[0].ClaimEstimation = d("9")
	lokasi[0].ObjekGolf[0].Coverages[0].Currency = "IDR" // Golf CURR_ID = objek saja (L11343), juga bila IsMBD
	teks := func(b models.BarisFacInProduction) string {
		return strings.Join([]string{b.OccupationCode, b.ObjectItem, b.CurrID, b.NoSerial, b.ObjectNameID, b.MOIDKasus, b.TopRisk,
			b.IdxObjectItem, teksPolis(b.TSI100), teksPolis(b.LoL), b.Kabupaten, b.Kelurahan, b.ObjectName, b.LayerNo, b.JnReas}, "|")
	}
	for nama, u := range map[string]struct {
		l   liniFacInOffer
		mau []string
	}{
		"Aneka": {liniFacInOffer{aneka: true}, []string{"OKU1|A|IDR|SN-1|NULL|MO-UJI|||<nil>|7|UJI KOTA|UJI RW|LOKASI UJI|NULL|10007",
			"OKU2|B|IDR||NULL|MO-UJI|||<nil>|0|UJI KOTA|UJI RW|LOKASI UJI|NULL|Q1", "OKU2|B|IDR||NULL|MO-UJI|||<nil>|0|UJI KOTA|UJI RW|LOKASI UJI|NULL|S1"}},
		"MBD": {liniFacInOffer{aneka: true, mbd: true}, []string{"OKU1|A|IDR|SN-1|NULL|MO-UJI|||<nil>|7|UJI KOTA|UJI RW|LOKASI UJI|NULL|10007",
			"OKU2|B|USD||NULL|MO-UJI|||<nil>|0|UJI KOTA|UJI RW|LOKASI UJI|NULL|Q1", "OKU2|B|USD||NULL|MO-UJI|||<nil>|0|UJI KOTA|UJI RW|LOKASI UJI|NULL|S1"}},
		"Golf": {liniFacInOffer{golf: true, mbd: true}, []string{"|G|EUR||NULL|MO-UJI|||<nil>|0|UJI KOTA|UJI RW|LOKASI UJI|NULL|10007"}},
	} {
		b, err := barisFacInProductionAneka(dp, objek, lokasi, u.l)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, x := range b {
			got = append(got, teks(x))
		}
		if !reflect.DeepEqual(got, u.mau) {
			t.Errorf("%s:\n%q", nama, got)
		}
		if b[0].IDPega != "UJI-NB-1" || b[0].BusinessCode != "10048" || teksPolis(b[0].TSI) != "10" ||
			teksPolis(b[0].RIComm) != "0.1" {
			t.Errorf("%s dasar: %+v", nama, b[0])
		}
	}
	if _, err := barisFacInProductionAneka(dp, nil, lokasi, liniFacInOffer{aneka: true}); err == nil {
		t.Error("tidak sejalan")
	}
}

// TestSelesaiPolicyAneka - TL Accept polis deklarasi Aneka (IsFlagUW 0, K73-41) -> selesai Policy: FACINPRODUCTION cabang 19,
// M_RISK_LOSS_PROFILE TSI NULL + premi 0 (K73-48).
func TestSelesaiPolicyAneka(t *testing.T) {
	st := submitUji()
	lokasi := lokasiFacInOfferAneka() // proteksi P4 berlaku di TL Policy: coverage berketerangan, premi terisi
	setiapCoverageAneka(lokasi, true, true, func(_ int, _ *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		c.CoverageNote, c.Premium = "UJI", d("10")
	})
	svc := layananPolicyAneka(st, "30", lokasi)
	h, err := svc.AjukanSubmit(context.Background(), pelakuTL, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "UJI"})
	if err != nil || len(st.tulis) != 1 || st.tulis[0].Policy == nil || h.NomorPolis == "" {
		t.Fatalf("%v %+v", err, h)
	}
	p := st.tulis[0].Policy
	if len(p.Produksi) != 3 || p.Produksi[0].ObjectItem != "A" || p.Produksi[1].OccupationCode != "OKU2" {
		t.Errorf("produksi %+v", p.Produksi)
	}
	r := p.RisikoLoss
	if r.TSINusantaraRe != nil || teksPolis(r.PremiNusantaraRe) != "0" || teksPolis(r.PremiumOR) != "0" || teksPolis(r.PremiumQS) != "0" ||
		teksPolis(r.PremiumFacout) != "0" || r.COB != "UJI COB" {
		t.Errorf("risiko %+v", r)
	}
}
