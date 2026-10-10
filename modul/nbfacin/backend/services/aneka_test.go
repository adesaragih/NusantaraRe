package services

import (
	"context"

	"errors"
	"github.com/cockroachdb/apd/v3"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// anekaTiruan - penyimpan berkeadaan: lokasi per case, diganti utuh; `tulis` = salinan terakhir yang diterima.
type anekaTiruan struct {
	ada     map[string][]models.LokasiAneka
	tulis   *[]models.LokasiAneka
	share   *apd.Decimal
	proRate **apd.Decimal // nilai tersimpan T_GENERAL_POLIS.PRO_RATE_PERCENT (nil = tanpa penyimpanan)
	// shareTulis - % Share RNM terakhir yang ditulis UbahPercentShare (nil = tidak direkam).
	shareTulis **apd.Decimal
}

func (a anekaTiruan) ProRateKasus(context.Context, string) (*apd.Decimal, error) {
	if a.proRate == nil {
		return nil, nil
	}
	return *a.proRate, nil
}

func (a anekaTiruan) UbahPercentShare(_ context.Context, _ *db.Tx, _ string, ps *apd.Decimal) error {
	if a.shareTulis != nil {
		*a.shareTulis = ps
	}
	return nil
}

func (a anekaTiruan) PercentShareKasus(context.Context, string) (*apd.Decimal, error) {
	return a.share, nil
}

func (a anekaTiruan) BacaObjekAneka(_ context.Context, id string) ([]models.LokasiAneka, error) {
	d, ok := a.ada[id]
	if !ok {
		return nil, repository.ErrKasusTidakAda
	}
	return d, nil
}

func (a anekaTiruan) GantiObjekAneka(_ context.Context, _ *db.Tx, id string, lokasi []models.LokasiAneka, pr *apd.Decimal) error {
	if _, ok := a.ada[id]; !ok {
		return repository.ErrKasusTidakAda
	}
	if a.proRate != nil {
		*a.proRate = pr
	}
	a.ada[id] = lokasi
	if a.tulis != nil {
		*a.tulis = lokasi
	}
	return nil
}

// layananAneka - case "K" (Class of Business UJI, group 3) ber-BusinessCode `kode` (BUSINESS ber-GROUPPANEL 003).
func layananAneka(kode string, tulis *[]models.LokasiAneka) *Service {
	n := 0
	k := &models.Kasus{CaseID: "K", Opportunity: models.Opportunity{ClassOfBusiness: "UJI KELAS", GroupBusinessID: "3"}}
	return Baru(nil).DenganObjekAneka(anekaTiruan{ada: map[string][]models.LokasiAneka{"K": {}}, tulis: tulis}).
		DenganTransaksi(tanpaTx).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"K": k}}).
		DenganTableOfLimit(tolTiruan{bisnis: []models.BarisBisnisLini{{ID: kode, OldID: "36", GroupPanel: "003"}}}).
		DenganPilihanItem(pilihanItemTiruan{panggil: &n})
}

func lokasiAnekaUji(o ...models.ObjekAneka) []models.LokasiAneka {
	return []models.LokasiAneka{{ObjectNo: "1", Okupasi: []models.OkupasiAneka{{OccupationID: "OK1", Objek: o}},
		ObjekGolf: []models.ObjekAneka{{ObjectName: "UJI GOLF", Currency: "USD"}}}}
}

// TestGantiObjekAneka - tiket 73 P1: SuretyType diisi server dari BusinessCode (InputDtlObjectAneka_PreAct; isian klien
// diabaikan, kode di luar tujuh -> kosong), tanggal kabel -> teks Pega dan kembali DD-MM-YYYY, masukan pemanggil tidak diubah;
// 401 / 400 (lebar, tanggal, mata uang) / 503 / 404.
func TestGantiObjekAneka(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var tulis []models.LokasiAneka
	masuk := lokasiAnekaUji(models.ObjekAneka{ObjectName: "UJI BOND", Currency: "IDR", StartDate: "01-10-2026",
		Briguna: models.BrigunaAneka{DateOfBirth: "31-12-1990"}, Memorandum: models.MemorandumAccept{SuretyType: "KLIEN"}})
	d, err := layananAneka("10076", &tulis).GantiObjekAneka(ctx, akun, "K", masuk, nil)
	if err != nil {
		t.Fatal(err)
	}
	o := tulis[0].Okupasi[0].Objek[0]
	if o.Memorandum.SuretyType != "SB-BB" || tulis[0].ObjekGolf[0].Memorandum.SuretyType != "SB-BB" {
		t.Errorf("SuretyType: %+v", o.Memorandum)
	}
	if o.StartDate != "20261001T050000.000 GMT" || o.Briguna.DateOfBirth != "19901231T050000.000 GMT" || o.EndDate != "" {
		t.Errorf("tanggal Pega: %q %q %q", o.StartDate, o.Briguna.DateOfBirth, o.EndDate)
	}
	if b := d[0].Okupasi[0].Objek[0]; b.StartDate != "01-10-2026" || b.Briguna.DateOfBirth != "31-12-1990" {
		t.Errorf("baca ulang: %q %q", b.StartDate, b.Briguna.DateOfBirth)
	}
	if masuk[0].Okupasi[0].Objek[0].StartDate != "01-10-2026" || masuk[0].Okupasi[0].Objek[0].Memorandum.SuretyType != "KLIEN" {
		t.Error("masukan pemanggil diubah")
	}
	for kode, mau := range map[string]string{"10007": "SB-PB", "10077": "SB-APB", "10083": "SB-MB", "10011": "SB-PayB",
		"10075": "CB", "10126": "KBG", "10021": "", "": ""} {
		if _, err := layananAneka(kode, &tulis).GantiObjekAneka(ctx, akun, "K", lokasiAnekaUji(models.ObjekAneka{
			Memorandum: models.MemorandumAccept{SuretyType: "KLIEN"}}), nil); err != nil || tulis[0].Okupasi[0].Objek[0].Memorandum.SuretyType != mau {
			t.Errorf("kode %q: %q (%v), mau %q", kode, tulis[0].Okupasi[0].Objek[0].Memorandum.SuretyType, err, mau)
		}
	}
	if _, err := layananAneka("10076", nil).GantiObjekAneka(ctx, inti.Pelaku{}, "K", nil, nil); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	for nama, u := range map[string]struct {
		objek models.ObjekAneka
		pesan string
	}{
		"nama 4001":    {models.ObjekAneka{ObjectName: strings.Repeat("U", 4001)}, "lokasi[0].okupasi[0].objek[0].objectName paling banyak 4000"},
		"memo 4001":    {models.ObjekAneka{Memorandum: models.MemorandumAccept{PrincipalsData: strings.Repeat("U", 4001)}}, "objek[0].memorandum.principalsData paling banyak 4000"},
		"kapal 1001":   {models.ObjekAneka{MarineHull: models.MarineHull{GRTRT: strings.Repeat("1", 1001)}}, "objek[0].marineHull.grtrt paling banyak 1000"},
		"briguna 51":   {models.ObjekAneka{Briguna: models.BrigunaAneka{IDCardNo: strings.Repeat("1", 51)}}, "objek[0].briguna.idCardNo paling banyak 50"},
		"tanggal ISO":  {models.ObjekAneka{EndDate: "2026-10-01"}, "objek[0].endDate bukan tanggal DD-MM-YYYY"},
		"lahir 31-02":  {models.ObjekAneka{Briguna: models.BrigunaAneka{DateOfBirth: "31-02-1990"}}, "objek[0].briguna.dateOfBirth bukan tanggal"},
		"mata uang 11": {models.ObjekAneka{Currency: strings.Repeat("I", 11)}, "objek[0].currency paling banyak 10"},
	} {
		if _, err := Baru(nil).GantiObjekAneka(ctx, akun, "K", lokasiAnekaUji(u.objek), nil); !errors.Is(err, ErrMasukanObjek) ||
			!strings.Contains(err.Error(), u.pesan) {
			t.Errorf("%s: %v, mau 400 %q (sebelum DB)", nama, err, u.pesan)
		}
	}
	if _, err := Baru(nil).GantiObjekAneka(ctx, akun, "K", []models.LokasiAneka{{RiskLocation: strings.Repeat("U", 4001)}}, nil); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "lokasi[0].riskLocation paling banyak 4000") {
		t.Errorf("lokasi: %v", err)
	}
	if _, err := layananAneka("10076", nil).GantiObjekAneka(ctx, akun, "K", lokasiAnekaUji(models.ObjekAneka{Currency: "XYZ"}), nil); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), `lokasi[0].okupasi[0].objek[0].currency "XYZ" tidak ada`) {
		t.Errorf("mata uang: %v", err)
	}
	// Aviation Hull tanpa Currency: sah, CURRENCY tidak dibaca.
	n := 0
	svc := layananAneka("10076", nil).DenganPilihanItem(pilihanItemTiruan{panggil: &n})
	if _, err := svc.GantiObjekAneka(ctx, akun, "K", []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{
		{AviationHull: models.AviationHull{Pilots: "2"}}}}}}}, nil); err != nil || n != 0 {
		t.Errorf("tanpa mata uang: %v, CURRENCY dibaca %d kali", err, n)
	}
	if _, err := Baru(nil).GantiObjekAneka(ctx, akun, "K", nil, nil); !errors.Is(err, ErrObjekTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	if _, err := layananAneka("10076", nil).GantiObjekAneka(ctx, akun, "K-404", nil, nil); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("case tidak ada: %v", err)
	}
	if _, err := layananAneka("10076", nil).BacaObjekAneka(ctx, "K-404"); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("baca 404: %v", err)
	}
	if _, err := Baru(nil).BacaObjekAneka(ctx, "K"); !errors.Is(err, ErrObjekTanpaDatabase) {
		t.Errorf("baca tanpa DB: %v", err)
	}
}

// TestObjekAnekaBacaSaja - hak ubah = tab Object FIRE: tahap Policy (FlagOnGoingPolicy 1) -> 409 sebelum isian / DB.
func TestObjekAnekaBacaSaja(t *testing.T) {
	svc := layananAneka("10076", nil).DenganTahap(tahapUji)
	_, err := svc.GantiObjekAneka(context.Background(), inti.Pelaku{AkunID: "UJI"}, "UJI-TL1", lokasiAnekaUji(models.ObjekAneka{
		ObjectName: strings.Repeat("U", 5000)}), nil)
	if !errors.Is(err, ErrTahapPolicyBacaSaja) {
		t.Errorf("tahap Policy: %v", err)
	}
}

// TestProRateKasusAneka - K73-12 (migrasi 255): GET memakai T_GENERAL_POLIS.PRO_RATE_PERCENT tersimpan; NULL -> dari tanggal
// polis (CPCA langkah 4, Day 365); PUT menulis nilai badan (simpan8; nil = NULL).
func TestProRateKasusAneka(t *testing.T) {
	ctx := context.Background()
	var tersimpan *apd.Decimal
	k := &models.Kasus{CaseID: "K", General: models.General{StartDateTime: "20260101T050000.000 GMT", EndDateTime: "20260702T050000.000 GMT"}}
	svc := Baru(nil).DenganObjekAneka(anekaTiruan{ada: map[string][]models.LokasiAneka{"K": {}}, proRate: &tersimpan}).
		DenganTransaksi(tanpaTx).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"K": k}})
	samaD(t, "NULL -> tanggal", svc.ProRateKasusAneka(ctx, "K"), "49.8630137")
	if _, err := svc.GantiObjekAneka(ctx, inti.Pelaku{AkunID: "UJI"}, "K", nil, ds("12.123456789")); err != nil {
		t.Fatal(err)
	}
	samaD(t, "PUT menulis (simpan8)", tersimpan, "12.12345679")
	samaD(t, "GET tersimpan", svc.ProRateKasusAneka(ctx, "K"), "12.12345679")
	if _, err := svc.GantiObjekAneka(ctx, inti.Pelaku{AkunID: "UJI"}, "K", nil, nil); err != nil || tersimpan != nil {
		t.Errorf("PUT kosong = NULL: %v %v", tersimpan, err)
	}
	if pr := Baru(nil).ProRateKasusAneka(ctx, "K"); pr != nil {
		t.Errorf("tanpa DB: %v", pr)
	}
}

// TestPertahankanSpreadingAneka - tiket 73 P3 (pola K48-7): coverage dikirim ulang mewarisi spreading tersimpan di posisi (jalur
// objek, indeks) yang oldId-nya sama; oldId beda / posisi baru -> kosong; spreading kiriman klien diabaikan.
func TestPertahankanSpreadingAneka(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	or := []models.BarisSpreading{{TreatyType: "OR"}}
	qs := []models.BarisSpreading{{TreatyType: "QS"}}
	var tulis []models.LokasiAneka
	svc := layananAneka("10076", &tulis)
	tiruan := svc.objekAneka.(anekaTiruan)
	tiruan.ada["K"] = lokasiAnekaUji(models.ObjekAneka{Coverages: []models.CoverageAneka{{OldID: "1", Spreading: or},
		{OldID: "2", Spreading: qs}}})
	tiruan.ada["K"][0].ObjekGolf[0].Coverages = []models.CoverageAneka{{OldID: "9", Spreading: qs}}
	masuk := lokasiAnekaUji(models.ObjekAneka{Currency: "IDR", Coverages: []models.CoverageAneka{{OldID: "1"}, {OldID: "3"},
		{OldID: "2", Spreading: or}}})
	masuk[0].ObjekGolf[0].Coverages = []models.CoverageAneka{{OldID: "9", Spreading: or}}
	if _, err := svc.GantiObjekAneka(ctx, akun, "K", masuk, nil); err != nil {
		t.Fatal(err)
	}
	c := tulis[0].Okupasi[0].Objek[0].Coverages
	if len(c[0].Spreading) != 1 || c[0].Spreading[0].TreatyType != "OR" || len(c[1].Spreading) != 0 || c[1].Spreading == nil ||
		len(c[2].Spreading) != 0 {
		t.Errorf("objek okupasi: %+v", c)
	}
	if g := tulis[0].ObjekGolf[0].Coverages; len(g[0].Spreading) != 1 || g[0].Spreading[0].TreatyType != "QS" {
		t.Errorf("Golf: %+v", g)
	}
	if masuk[0].Okupasi[0].Objek[0].Coverages[0].Spreading != nil {
		t.Error("masukan pemanggil diubah")
	}
	// tanpa coverage dikirim: tetap tersimpan (coverage lama dipindah repository bersama spreading-nya).
	tiruan.ada["K"] = nil
	if _, err := svc.GantiObjekAneka(ctx, akun, "K", lokasiAnekaUji(models.ObjekAneka{Currency: "IDR"}), nil); err != nil {
		t.Errorf("tanpa coverage: %v", err)
	}
}
