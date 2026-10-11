package services

import (
	"context"
	"errors"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// anekaSimulasi - penyimpan objek Aneka yang GantiObjekAneka-nya menjalankan lintasan tulis repository tanpa basis data
// (repository.SimulasiTulisAneka: traversal, pengelompokan, periksaBindBaris, array-bind / per baris yang sama dengan Oracle).
type anekaSimulasi struct {
	anekaTiruan
	eks *[]repository.EksekusiSimulasi
}

func (a anekaSimulasi) GantiObjekAneka(ctx context.Context, tx *db.Tx, id string, lokasi []models.LokasiAneka, pr *apd.Decimal) error {
	e, err := repository.SimulasiTulisAneka(ctx, id, lokasi)
	*a.eks = e
	if err != nil {
		return err
	}
	return a.anekaTiruan.GantiObjekAneka(ctx, tx, id, lokasi, pr)
}

// lokasiCaseDEV - bentuk case work owner bug DEV 10-10-2026 (dari layar, rekonstruksi): 2 lokasi, objek "TES LIABILITY" di
// lokasi 1 dan satu objek di lokasi 2, coverage PUBLIC LIABILITY; TSI 500.000.000.000.000 + 50.000.000.000.000 dan premi 1 %
// `[dugaan: dipilih agar ringkasan layar cocok - TSI spreaded 137.500.000.000.000, premi 1.375.000.000.000 pada % Share 25]`.
func lokasiCaseDEV() []models.LokasiAneka {
	objek := func(nama, tsi, premi string) []models.ObjekAneka {
		return []models.ObjekAneka{{ObjectName: nama, Currency: "IDR", TSI: d(tsi), Coverages: []models.CoverageAneka{{
			Coverage: "PUBLIC LIABILITY", OldID: "PL", TSI: d(tsi), TSILiability: d(tsi), Premium: d(premi)}}}}
	}
	return []models.LokasiAneka{
		{ObjectNo: "1", RiskLocation: "LOKASI 1", Okupasi: []models.OkupasiAneka{{OccupationName: "OKU", Objek: objek("TES LIABILITY",
			"500000000000000", "5000000000000")}}},
		{ObjectNo: "2", RiskLocation: "LOKASI 2", Okupasi: []models.OkupasiAneka{{OccupationName: "OKU", Objek: objek("OBJEK 2",
			"50000000000000", "500000000000")}}},
	}
}

// TestSimpanSpreadingAnekaBindCaseDEV - bug DEV 10-10-2026 (ORA-01722 T_SPREADINGLIST): jalur PUT spreading-aneka utuh (service
// -> lintasan tulis repository) atas bentuk case work owner, % Share RNM 25, satu baris Copy To All "2024 QS 155M FAC" 100 %
// per coverage. Setiap bind angka lolos pola mask (periksaBindBaris di dalam lintasan), T_SPREADINGLIST dikirim per baris, dan
// nilainya = ringkasan layar. Instrumen diuji dulu: nilai 9 desimal di lintasan yang sama -> ErrBindTidakSah.
func TestSimpanSpreadingAnekaBindCaseDEV(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	salah := lokasiCaseDEV()
	salah[0].Okupasi[0].Objek[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "QS", SharePercentage: d("33.333333333")}}
	if _, err := repository.SimulasiTulisAneka(ctx, "K", salah); !errors.Is(err, repository.ErrBindTidakSah) {
		t.Fatalf("instrumen: 9 desimal lolos: %v", err)
	}

	var eks []repository.EksekusiSimulasi
	var tulis []models.LokasiAneka
	var share *apd.Decimal
	pr := d("100")
	a := anekaSimulasi{anekaTiruan: anekaTiruan{ada: map[string][]models.LokasiAneka{"K": lokasiCaseDEV()}, tulis: &tulis,
		share: d("25"), shareTulis: &share, proRate: &pr}, eks: &eks}
	k := &models.Kasus{CaseID: "K", Opportunity: models.Opportunity{ClassOfBusiness: "UJI KELAS", GroupBusinessID: "3"}}
	svc := Baru(nil).DenganObjekAneka(a).DenganTransaksi(tanpaTx).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"K": k}}).
		DenganTableOfLimit(tolTiruan{bisnis: satuBisnis("003", "30")}).DenganPilihanItem(pilihanItemTiruan{panggil: new(int)}).
		DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{{PercentShare: d("25"), StartDateTime: "20250101T170000.000 GMT",
			Lokasi: []models.LokasiSpreading{{ObjectNo: "1"}, {ObjectNo: "2"}}}}, idUang: map[string]string{"IDR": "10026"}})
	badan := lokasiCaseDEV()
	for i := range badan {
		badan[i].Okupasi[0].Objek[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "QS", TreatyName: "2024 QS 155M FAC",
			SharePercentage: d("100")}}
	}
	if _, err := svc.SimpanSpreadingAneka(ctx, akun, "K", d("25"), badan); err != nil {
		t.Fatal(err)
	}
	var spr [][]any
	for _, e := range eks {
		if e.Tabel != repository.TabelSpreadingList {
			continue
		}
		if len(e.Baris) != 1 {
			t.Errorf("T_SPREADINGLIST %d baris per eksekusi, mau per baris", len(e.Baris))
		}
		spr = append(spr, e.Baris...)
	}
	if len(spr) != 2 {
		t.Fatalf("T_SPREADINGLIST %d baris, mau 2", len(spr))
	}
	jumlahTSI, jumlahPremi := new(apd.Decimal), new(apd.Decimal)
	for i, b := range spr {
		for c := 8; c <= 12; c++ {
			if s, ok := b[c].(string); ok && !repository.PolaAngkaMasuk.MatchString(s) {
				t.Errorf("baris %d bind :%d di luar mask", i+1, c+1)
			}
		}
		if b[6] != "QS" || b[7] != "2024 QS 155M FAC" || b[8] != "100" || b[13] != "IDR" {
			t.Errorf("baris %d: %v", i+1, b[6:9])
		}
		tsi, premi := d(b[9].(string)), d(b[11].(string))
		jumlahTSI, jumlahPremi = tambah(jumlahTSI, tsi), tambah(jumlahPremi, premi)
	}
	if spr[0][9] != "125000000000000" || t8(jumlahTSI) != "137500000000000" || t8(jumlahPremi) != "1375000000000" {
		t.Errorf("TSI spreaded %v / jumlah %s, premi %s", spr[0][9], t8(jumlahTSI), t8(jumlahPremi))
	}
}

// TestSimpanSpreadingAnekaBindCaseDEVTahapUW - bug DEV 10-10-2026 11:31 (`T_SPREADINGLIST kolom CLAIM_ESTIMATION baris 1:
// panjang bulat 1, desimal 20`): cabang uw SimpanSpreadingAneka memanggil totalAneka untuk periksaJenisTreaty SESUDAH
// pembulatan; totalLokasiAneka menimpa ClaimEstimation baris (LoL kosong -> 0E-20) dan pohon itu yang ditulis. Bentuk case DEV
// (LoL kosong) di tahap uw lewat lintasan tulis repository: setiap bind lolos mask.
func TestSimpanSpreadingAnekaBindCaseDEVTahapUW(t *testing.T) {
	ctx := context.Background()
	const id = "UJI-UW"
	var eks []repository.EksekusiSimulasi
	var tulis []models.LokasiAneka
	var share *apd.Decimal
	pr := d("100")
	a := anekaSimulasi{anekaTiruan: anekaTiruan{ada: map[string][]models.LokasiAneka{id: lokasiCaseDEV()}, tulis: &tulis,
		share: d("25"), shareTulis: &share, proRate: &pr}, eks: &eks}
	k := &models.Kasus{CaseID: id, Opportunity: models.Opportunity{ClassOfBusiness: "UJI KELAS", GroupBusinessID: "3"}}
	svc := Baru(nil).DenganObjekAneka(a).DenganTransaksi(tanpaTx).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{id: k}}).
		DenganTableOfLimit(tolTiruan{bisnis: satuBisnis("003", "30")}).DenganPilihanItem(pilihanItemTiruan{panggil: new(int)}).
		DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{{PercentShare: d("25"), StartDateTime: "20250101T170000.000 GMT",
			Lokasi: []models.LokasiSpreading{{ObjectNo: "1"}, {ObjectNo: "2"}}}}, kap: &kapasitasUji, idUang: map[string]string{"IDR": "10026"}}).
		DenganTahap(tahapUji)
	badan := lokasiCaseDEV()
	for i := range badan {
		badan[i].Okupasi[0].Objek[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "10015", TreatyName: "2024 QS 155M FAC",
			SharePercentage: d("100")}}
	}
	p := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInJuniorUnderwritingB"}}
	if _, uw, err := svc.penjagaPutSpreading(ctx, p, id); err != nil || !uw {
		t.Fatalf("bukan tahap uw: %v %v", uw, err)
	}
	if _, err := svc.SimpanSpreadingAneka(ctx, p, id, d("25"), badan); err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, e := range eks {
		if e.Tabel == repository.TabelSpreadingList {
			n += len(e.Baris)
		}
	}
	if n != 2 {
		t.Fatalf("T_SPREADINGLIST %d baris, mau 2", n)
	}
	for _, c := range tulis[0].Okupasi[0].Objek[0].Coverages[0].Spreading {
		if c.ClaimEstimation != nil && -c.ClaimEstimation.Exponent > 8 {
			t.Errorf("ClaimEstimation tertulis berskala %d", -c.ClaimEstimation.Exponent)
		}
	}
}
