package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// akumulasiTiruan - PenyimpanAkumulasiCoverage dalam memori.
type akumulasiTiruan struct{ tulis [][]models.AkumulasiCoverage }

func (a *akumulasiTiruan) UbahAkumulasiCoverage(_ context.Context, _ *db.Tx, _ string, baris []models.AkumulasiCoverage) error {
	a.tulis = append(a.tulis, baris)
	return nil
}

// kasusTigaCoverage - kasusDuaItem dengan ID baris coverage C1 11, C2 12, C3 13.
func kasusTigaCoverage() models.KasusSpreading {
	k := kasusDuaItem("50", "USD")
	k.Lokasi[0].Items[0].Coverages[1].ID = "12"
	k.Lokasi[0].Items[1].Coverages[0].ID = "13"
	return k
}

// badanAkumulasi - badan bentuk kasusTigaCoverage; akumulasi baru + medan lain yang "diubah" (harus diabaikan).
func badanAkumulasi() []models.ObjekFire {
	c := func(old, kode string) models.CoverageObjek {
		return models.CoverageObjek{OldID: old, AccumulationCode: kode, AccumulationDescription: "ALAMAT " + kode,
			Premium: d("999"), Rate: d("9"), TSI: d("1"), Coverage: "UBAH"}
	}
	return []models.ObjekFire{{Items: []models.ItemObjek{{Coverages: []models.CoverageObjek{c("C1", "AK1"), c("C2", "")}},
		{Coverages: []models.CoverageObjek{c("C3", "AK3")}}}}}
}

func layananAkumulasi(tahap tahapTiruan) (*Service, *akumulasiTiruan, *spreadingTiruan, objekTiruan) {
	a := &akumulasiTiruan{}
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusTigaCoverage()}}
	ob := objekTiruan{ada: map[string][]models.ObjekFire{}}
	for id := range tahap {
		ob.ada[id] = []models.ObjekFire{}
	}
	ob.ada["NB-1"] = []models.ObjekFire{}
	svc := Baru(nil).DenganObjek(ob).DenganSpreading(src).DenganAkumulasiCoverage(a).DenganTransaksi(tanpaTx).DenganTahap(tahap)
	return svc, a, src, ob
}

// TestMatriksAkumulasiUW - tiket 60 + RALAT K60-2: tahap × workbasket untuk PUT …/akumulasi dan ubahAkumulasi GET kasus;
// enam posisi IsUWSpread (termasuk DepHead / UnderwritingFinancial) lolos untuk pemegang.
func TestMatriksAkumulasiUW(t *testing.T) {
	ctx := context.Background()
	posisi := map[string]error{"ReasFacInUnderwriting": nil, "ReasFacInSeniorUnderwriting": nil, "ReasFacInJuniorUnderwritingB": nil,
		"ReasFacInJuniorUnderwritingA": nil, "ReasFacInDepHeadUnderwriting": nil,
		"ReasFacInUnderwritingFinancial": nil, "ReasFacInJuniorUnderwriting": ErrAkumulasiHanyaUW,
		"Offer": ErrAkumulasiHanyaUW, "ReasFacInMarketing": ErrAkumulasiHanyaUW, "ReasFacInTeamLeader": ErrAkumulasiHanyaUW,
		"Decline": ErrAkumulasiHanyaUW}
	tahap := tahapTiruan{}
	for p := range posisi {
		tahap["A-"+p] = models.TahapKasus{Position: p, FlagOnGoingPolicy: "0"}
	}
	for p, mau := range posisi {
		for _, pegang := range []bool{true, false} {
			peran := "ReasFacInLainUji"
			if pegang {
				peran = p
			}
			if p == "Offer" && pegang {
				peran = "ReasFacInMarketing"
			}
			svc, a, _, _ := layananAkumulasi(tahap)
			pl := inti.Pelaku{AkunID: "UJI", Peran: []string{peran}}
			harap := mau
			if mau == nil && !pegang {
				harap = inti.ErrTanpaWewenang
			}
			_, err := svc.SimpanAkumulasi(ctx, pl, "A-"+p, badanAkumulasi())
			if !galatSama(err, harap) || (harap != nil && len(a.tulis) != 0) {
				t.Errorf("%s pegang=%v: %v, mau %v (tulis %d)", p, pegang, err, harap, len(a.tulis))
			}
			if got := svc.AksesAkumulasi(ctx, pl, "A-"+p); got != (harap == nil) {
				t.Errorf("akses %s pegang=%v: %v", p, pegang, got)
			}
		}
	}
	// 401 lebih dulu; tanpa pembaca tahap: akses false, endpoint tidak dijaga (penjaga tidak aktif).
	svc, _, _, _ := layananAkumulasi(tahap)
	if _, err := svc.SimpanAkumulasi(ctx, inti.Pelaku{}, "A-ReasFacInUnderwriting", badanAkumulasi()); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	svc, _, _, _ = layananAkumulasi(nil)
	svc.DenganTahap(nil)
	if svc.AksesAkumulasi(ctx, pelakuUji, "NB-1") {
		t.Errorf("tanpa pembaca tahap: akses harus false")
	}
	if _, err := svc.SimpanAkumulasi(ctx, pelakuUji, "NB-1", badanAkumulasi()); err != nil {
		t.Errorf("tanpa pembaca tahap: %v", err)
	}
	if svc.AksesAkumulasi(ctx, pelakuUji, "UJI-GALAT") {
		t.Errorf("galat baca tahap harus false")
	}
}

// TestSimpanAkumulasiHanyaDuaKolom - badan yang juga mengubah premi / rate / TSI / coverage hanya menulis ID + dua kolom
// akumulasi (menurut posisi, oldId dicocokkan); objek tidak diganti, spreading tidak ditulis.
func TestSimpanAkumulasiHanyaDuaKolom(t *testing.T) {
	tahap := tahapTiruan{"A-UW": {Position: "ReasFacInUnderwriting", FlagOnGoingPolicy: "0"}}
	svc, a, src, ob := layananAkumulasi(tahap)
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInUnderwriting"}}
	if _, err := svc.SimpanAkumulasi(context.Background(), uw, "A-UW", badanAkumulasi()); err != nil {
		t.Fatal(err)
	}
	mau := []models.AkumulasiCoverage{{ID: "11", AccumulationCode: "AK1", AccumulationDescription: "ALAMAT AK1"},
		{ID: "12", AccumulationCode: "", AccumulationDescription: "ALAMAT "}, {ID: "13", AccumulationCode: "AK3", AccumulationDescription: "ALAMAT AK3"}}
	if len(a.tulis) != 1 || !reflect.DeepEqual(a.tulis[0], mau) || len(src.tulis) != 0 || len(ob.ada["A-UW"]) != 0 {
		t.Errorf("tulis %+v spreading %d objek %d", a.tulis, len(src.tulis), len(ob.ada["A-UW"]))
	}
}

// TestSimpanAkumulasiBentukDanLebar - bentuk / oldId beda -> 409 ErrObjekBerubah tanpa tulisan; lebar > 50 / 500 byte ->
// 400 (aturan lebarCoverage PUT …/objek).
func TestSimpanAkumulasiBentukDanLebar(t *testing.T) {
	ctx := context.Background()
	tahap := tahapTiruan{"A-UW": {Position: "ReasFacInSeniorUnderwriting", FlagOnGoingPolicy: "0"}}
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInSeniorUnderwriting"}}
	for nama, ubah := range map[string]func([]models.ObjekFire) []models.ObjekFire{
		"kurang coverage": func(b []models.ObjekFire) []models.ObjekFire {
			b[0].Items[0].Coverages = b[0].Items[0].Coverages[:1]
			return b
		},
		"oldId beda": func(b []models.ObjekFire) []models.ObjekFire { b[0].Items[1].Coverages[0].OldID = "C9"; return b },
		"item lebih": func(b []models.ObjekFire) []models.ObjekFire {
			b[0].Items = append(b[0].Items, models.ItemObjek{})
			return b
		},
		"objek kosong": func([]models.ObjekFire) []models.ObjekFire { return nil },
	} {
		svc, a, _, _ := layananAkumulasi(tahap)
		if _, err := svc.SimpanAkumulasi(ctx, uw, "A-UW", ubah(badanAkumulasi())); !errors.Is(err, ErrObjekBerubah) || len(a.tulis) != 0 {
			t.Errorf("%s: %v tulis %d", nama, err, len(a.tulis))
		}
	}
	svc, a, _, _ := layananAkumulasi(tahap)
	b := badanAkumulasi()
	b[0].Items[0].Coverages[0].AccumulationCode = strings.Repeat("A", 51)
	b[0].Items[1].Coverages[0].AccumulationDescription = strings.Repeat("A", 501)
	_, err := svc.SimpanAkumulasi(ctx, uw, "A-UW", b)
	if !errors.Is(err, ErrMasukanObjek) || !strings.Contains(err.Error(), "accumulationCode paling banyak 50") ||
		!strings.Contains(err.Error(), "accumulationDescription paling banyak 500") || len(a.tulis) != 0 {
		t.Errorf("lebar: %v", err)
	}
	b = badanAkumulasi()
	b[0].Items[0].Coverages[0].AccumulationCode = strings.Repeat("A", 50)
	if _, err := svc.SimpanAkumulasi(ctx, uw, "A-UW", b); err != nil {
		t.Errorf("tepat 50: %v", err)
	}
}
