package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// csvTiruan - pencarian master Save Data CSV tiruan (nilai sintetis UJI-).
type csvTiruan struct {
	rw       map[string][]models.RWCsv // zip|kelurahan
	risk     map[string]string         // alamat|zip
	akum     map[string]string         // note
	negara   map[string]string
	prov     map[string]string
	mu       map[string]models.MataUangCsv
	item     map[string]string
	okupasi  map[string]string
	coverage map[string]bool
	hapus    []string
	panggil  []string       // kode pos tiap RWKelurahan
	n        map[string]int // jumlah panggilan per metode (query tiruan)
}

func (c *csvTiruan) hitung(m string) {
	if c.n == nil {
		c.n = map[string]int{}
	}
	c.n[m]++
}

func (c *csvTiruan) total() int {
	s := 0
	for _, v := range c.n {
		s += v
	}
	return s
}

func (c *csvTiruan) RWKelurahan(_ context.Context, zip, kel string) ([]models.RWCsv, error) {
	c.hitung("RWKelurahan")

	c.panggil = append(c.panggil, zip)
	return c.rw[zip+"|"+kel], nil
}
func (c *csvTiruan) IDRiskAddress(_ context.Context, a, z string) (string, error) {
	c.hitung("IDRiskAddress")

	return c.risk[a+"|"+z], nil
}
func (c *csvTiruan) IDAkumulasi(_ context.Context, n string) (string, error) {
	c.hitung("IDAkumulasi")
	return c.akum[n], nil
}
func (c *csvTiruan) InisialNegara(_ context.Context, n string) (string, error) {
	c.hitung("InisialNegara")

	return c.negara[n], nil
}
func (c *csvTiruan) IDProvinsi(_ context.Context, n string) (string, error) {
	c.hitung("IDProvinsi")
	return c.prov[n], nil
}
func (c *csvTiruan) MataUang(_ context.Context, n string) (models.MataUangCsv, bool, error) {
	c.hitung("MataUang")

	m, ok := c.mu[n]
	return m, ok, nil
}
func (c *csvTiruan) IDJenisItemFire(_ context.Context, n string) (string, error) {
	c.hitung("IDJenisItemFire")

	return c.item[n], nil
}
func (c *csvTiruan) NamaOkupasiFire(_ context.Context, k string) (string, error) {
	c.hitung("NamaOkupasiFire")

	return c.okupasi[k], nil
}
func (c *csvTiruan) CoverageAda(_ context.Context, id []string) (map[string]bool, error) {
	c.hitung("CoverageAda")

	m := map[string]bool{}
	for _, v := range id {
		m[v] = c.coverage[v]
	}
	return m, nil
}
func (c *csvTiruan) HapusCurrencyList(_ context.Context, _ *db.Tx, id string) error {
	c.hapus = append(c.hapus, id)
	return nil
}

// addCsvTiruan / penambahCsvTiruan - CZone menurut kode pos dan master akumulasi (port tiket 46).
type addCsvTiruan struct{}

func (addCsvTiruan) CZoneZip(_ context.Context, zip string) (models.CZoneZip, bool, error) {
	if zip == "12345" {
		return models.CZoneZip{ID: "UJI-CZ-ID", Code: "UJI-CZ"}, true, nil
	}
	return models.CZoneZip{}, false, nil
}
func (addCsvTiruan) CariZip(context.Context, string, string, int, int) ([]models.ZipAkumulasi, int, error) {
	return nil, 0, nil
}

type penambahCsvTiruan struct{ masukan []map[string]string }

func (p *penambahCsvTiruan) Tambah(_ context.Context, _ inti.Pelaku, _ string, m map[string]string) (string, error) {
	p.masukan = append(p.masukan, m)
	return "UJI-AKU-BARU", nil
}

func csvMasterUji() *csvTiruan {
	return &csvTiruan{
		rw:       map[string][]models.RWCsv{"12345|UJI KEL": {{Note: "UJI KEL", District: "UJI KEC", City: "UJI KOTA", Province: "UJI PROV", Nation: "UJI NEGARA"}}},
		risk:     map[string]string{},
		akum:     map[string]string{},
		negara:   map[string]string{"UJI NEGARA": "UJ"},
		prov:     map[string]string{"UJI PROV": "UJI-P1"},
		mu:       map[string]models.MataUangCsv{"IDR": {ID: "1", OldID: "01", Currency: "IDR"}},
		item:     map[string]string{"BUILDING": "UJI-OI1"},
		okupasi:  map[string]string{"2001": "UJI OKUPASI"},
		coverage: map[string]bool{"100815": true, "100828": true, "100825": true, "100829": true, "100840": true},
	}
}

// barisCsvUji - satu baris template sintetis (55 kunci; kosong kecuali yang diisi), ubah = pasangan kolom/nilai tambahan.
func barisCsvUji(ubah ...string) map[string]string {
	r := map[string]string{}
	for _, k := range KolomCsvObjek {
		r[k] = ""
	}
	for k, v := range map[string]string{"NO": "1", "NO_OFFER_SLIP": "UJI-SLIP-1", "BEGIN_DATE": "19/10/2026", "END_DATE": "19/10/2027",
		"OCCUPATION_CODE": "2001", "CLASS_OF_CONSTRUCTION": "UJI KELAS 1", "TYPE_ADDRESS": "jl.", "ADDRESS": "uji raya",
		"BUILDING_NO": "7", "POSTAL_CODE": "12345", "KELURAHAN": "uji kel", "OBJECT_TYPE": "dwelling house", "INTEREST": "UJI",
		"TOP_RISK": "Y", "ITEM_TYPE": "BUILDING", "TSI_OBJECT_ITEM": "1.000.000,50", "CURRENCY": "IDR", "LOL": "500.000",
		"RATE_FLEXAS": "0,25", "NET_RATE_FLEXAS": "0,2", "DISCOUNT_RATE_FLEXAS": "10", "ACCUMULATION_CODE": "UJI-AKU-CSV",
		"RATE_RSMDCC": "0"} {
		r[k] = v
	}
	for i := 0; i+1 < len(ubah); i += 2 {
		r[ubah[i]] = ubah[i+1]
	}
	return r
}

// layananCsvUji - case UJI-NB-1 dengan satu objek lama, General tersimpan ber-Ceding C1.
func layananCsvUji(c *csvTiruan, p *penambahCsvTiruan) (*Service, *models.Kasus, objekTiruan, *models.AlamatBaru) {
	kasus := &models.Kasus{CaseID: "UJI-NB-1", General: models.General{ReffNumber: "UJI-SLIP-LAMA", StartDateTime: "20260101T050000.000 GMT",
		EndDateTime: "20270101T050000.000 GMT", TypeFacultative: "Facultative In", CedingList: []models.Ceding{{ID: "C1", Name: "UJI CEDING"}},
		Tersimpan: true}}
	ob := objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {{ObjectNo: "LAMA", ObjectType: "Office"}}}}
	var zip string
	alamat := &models.AlamatBaru{}
	n := 0
	svc := Baru(nil).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": kasus}}).DenganObjek(ob).DenganTransaksi(tanpaTx).
		DenganPilihanItem(pilihanItemTiruan{&n}).DenganRW(rwTiruan{zip: &zip, alamat: alamat}).DenganCsvObjek(c).
		DenganAddAkumulasi(addCsvTiruan{}, p).DenganJam(jamUjiCsv)
	return svc, kasus, ob, alamat
}

// jamUjiCsv - 07-10-2026 22:04:05 WIB (OfferingDate 20261007).
func jamUjiCsv() time.Time { return time.Date(2026, 10, 7, 15, 4, 5, 0, time.UTC) }

var pelakuCsv = inti.Pelaku{AkunID: "UJI-MKT", Peran: []string{"ReasFacInMarketing"}}

// TestSimpanCsvObjek - InsertUploadFire_act: satu lokasi per baris, huruf besar, RW / risk address baru / akumulasi baru
// (Keyword menurut Object Type baku), okupasi dari master, item + coverage FLEXAS, header D1 (b), Ceding dipertahankan,
// CurrencyList dihapus, objek lama diganti, "Success Insert Data".
func TestSimpanCsvObjek(t *testing.T) {
	ctx := context.Background()
	c, p := csvMasterUji(), &penambahCsvTiruan{}
	svc, kasus, ob, alamat := layananCsvUji(c, p)
	lama := kasusSatuCoverage("50")
	lama.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "10007", SharePercentage: d("100")}}
	sp := &spreadingTiruan{baca: []models.KasusSpreading{lama}}
	svc.DenganSpreading(sp)
	h, err := svc.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{barisCsvUji()})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(h.Pesan, []string{"Success Insert Data"}) || len(h.Baris) != 1 || len(ob.ada["UJI-NB-1"]) != 1 {
		t.Fatalf("%+v", h)
	}
	o := h.Baris[0]
	if o.ObjectNo != "1" || o.RoadType != "JL." || o.RoadName != "UJI RAYA" || o.RiskLocation != "JL. UJI RAYA" || o.BuildingNo != "7" ||
		o.ObjectType != "Dwelling House" || o.ObjectName != "DWELLING HOUSE" || !o.IsTopRisk || o.ZipCode != "12345" ||
		o.City != "UJI KOTA" || o.District != "UJI KEC" || o.Territory != "UJI KEL" || o.Province != "UJI PROV" || o.Country != "UJI NEGARA" ||
		o.RiskAddressID != "UJI000000000001" {
		t.Errorf("lokasi %+v", o)
	}
	if *alamat != (models.AlamatBaru{NationName: "UJI NEGARA", ProvinceName: "UJI PROV", DistrictName: "UJI KEC", CityName: "UJI KOTA",
		TerritoryName: "UJI KEL", Title: "JL.", Address: "UJI RAYA", PostalCode: "12345"}) {
		t.Errorf("risk address %+v", *alamat)
	}
	if len(p.masukan) != 1 || p.masukan[0]["accumulation"] != "10000001" || p.masukan[0]["accumulationType"] != "FLEXAS" ||
		p.masukan[0]["note"] != "UJI RAYA" || p.masukan[0]["keyword"] != "KOMPLEKS/KAV.PERUMAHAN/ VILLA/PEMUKIMAN" ||
		p.masukan[0]["scopeArea"] != "AREA" || p.masukan[0]["cZone"] != "UJI-CZ" || p.masukan[0]["provinceId"] != "UJI-P1" ||
		p.masukan[0]["zipCode"] != "12345" {
		t.Errorf("akumulasi %+v", p.masukan)
	}
	if len(o.Occupations) != 1 || o.Occupations[0] != (models.OkupasiObjek{OccupationID: "2001", OccupationName: "UJI OKUPASI",
		ConstructionClass: "UJI KELAS 1"}) || len(o.LossRecords) != 1 || o.LossRecords[0].Detail != "-" {
		t.Errorf("okupasi / kerugian %+v %+v", o.Occupations, o.LossRecords)
	}
	if len(o.Items) != 1 {
		t.Fatalf("item %+v", o.Items)
	}
	it := o.Items[0]
	if it.ItemType != "BUILDING" || it.ItemTypeID != "UJI-OI1" || it.Currency != "IDR" || t8(it.TSI) != "1000000.5" || it.Note != "UJI" ||
		it.IsAdjustable || len(it.Coverages) != 1 {
		t.Fatalf("item %+v", it)
	}
	cv := it.Coverages[0]
	if cv.OldID != "FLEXAS" || cv.Coverage != "100815" || cv.CoverageNote != "FLEXAS (PSAKI)" || cv.CoverageBasis != "1" ||
		t8(cv.Rate) != "0.25" || t8(cv.NetRate) != "0.2" || t8(cv.DiscountPercentage) != "10" || t8(cv.IndemnityPercentage) != "100" ||
		t8(cv.LimitOfLiability) != "500000" || cv.AccumulationCode != "UJI-AKU-BARU" || cv.AccumulationDescription != "UJI RAYA" ||
		cv.Premium == nil {
		t.Errorf("coverage %+v", cv)
	}
	g := kasus.General
	if g.ReffNumber != "UJI-SLIP-1" || g.StartDateTime != "20261019T050000.000 GMT" || g.EndDateTime != "20271019T050000.000 GMT" ||
		g.TypeFacultative != "Facultative In" || g.OfferingDate != "20261007" || !reflect.DeepEqual(g.CedingIDs, []string{"C1"}) {
		t.Errorf("General %+v", g)
	}
	if !reflect.DeepEqual(c.hapus, []string{"UJI-NB-1"}) || len(sp.tulis) != 0 {
		t.Errorf("CurrencyList %v / spreading ditulis %d", c.hapus, len(sp.tulis))
	}
}

// TestCsvObjekPesanDanMaster - KodePos tidak ditemukan -> pesan verbatim (tanpa "Success Insert Data"), tetap tersimpan,
// tanpa risk address / akumulasi baru; master yang sudah ada dipakai; akumulasi baru tanpa provinsi -> kode CSV.
func TestCsvObjekPesanDanMaster(t *testing.T) {
	ctx := context.Background()
	c, p := csvMasterUji(), &penambahCsvTiruan{}
	svc, _, _, _ := layananCsvUji(c, p)
	h, err := svc.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{barisCsvUji(), barisCsvUji("KELURAHAN", "lain", "TOP_RISK", "y",
		"ADJUSTABLE_PCT", "100")})
	if err != nil {
		t.Fatal(err)
	}
	mau := []string{"KodePos 12345 and kelurahan LAIN pada lokasi ke - 2 tidak ditemukan!"}
	if !reflect.DeepEqual(h.Pesan, mau) || len(h.Baris) != 2 || h.Baris[1].ObjectNo != "2" || h.Baris[1].City != "" ||
		h.Baris[1].Items[0].Coverages[0].AccumulationCode != "UJI-AKU-CSV" || len(p.masukan) != 1 || h.Baris[1].IsTopRisk ||
		h.Baris[1].Items[0].IsAdjustable || t8(h.Baris[1].Items[0].PctAdjustOther) != "100" {
		t.Errorf("%+v / %d", h, len(p.masukan))
	}
	// Master sudah ada: tanpa tulisan baru.
	c2, p2 := csvMasterUji(), &penambahCsvTiruan{}
	c2.risk["UJI RAYA|12345"], c2.akum["UJI RAYA"] = "UJI-RISK-ADA", "UJI-AKU-ADA"
	svc2, _, _, alamat2 := layananCsvUji(c2, p2)
	h2, err := svc2.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{barisCsvUji()})
	if err != nil || h2.Baris[0].RiskAddressID != "UJI-RISK-ADA" || h2.Baris[0].Items[0].Coverages[0].AccumulationCode != "UJI-AKU-ADA" ||
		len(p2.masukan) != 0 || alamat2.Address != "" {
		t.Errorf("master ada: %v %+v", err, h2)
	}
	// Provinsi tidak ada di master: akumulasi tidak dibuat, kode CSV tetap.
	c3, p3 := csvMasterUji(), &penambahCsvTiruan{}
	c3.prov = map[string]string{}
	svc3, _, _, _ := layananCsvUji(c3, p3)
	h3, err := svc3.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{barisCsvUji()})
	if err != nil || h3.Baris[0].Items[0].Coverages[0].AccumulationCode != "UJI-AKU-CSV" || len(p3.masukan) != 0 {
		t.Errorf("tanpa provinsi: %v %+v", err, h3)
	}
}

// TestCsvObjekItemDanHeader - tujuh slot item (TSI_OBJECT_ITEM_6 ikut dinormalisasi), lima peril, ADJUSTABLE_PCT; header baris
// terakhir yang terisi menang, kolom kosong tidak menimpa.
func TestCsvObjekItemDanHeader(t *testing.T) {
	ctx := context.Background()
	isi := []string{"RATE_RSMDCC", "0,1", "RATE_FWTWD", "0,1", "RATE_EQ", "0,1", "RATE_OTHERS", "0,1", "ADJUSTABLE_PCT", "80"}
	for _, sl := range ItemCsvObjek[1:] {
		isi = append(isi, "ITEM_TYPE"+sl, "BUILDING", "TSI_OBJECT_ITEM"+sl, "1.000")
	}
	c := csvMasterUji()
	svc, kasus, _, _ := layananCsvUji(c, &penambahCsvTiruan{})
	h, err := svc.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{barisCsvUji(isi...),
		barisCsvUji("NO_OFFER_SLIP", "", "BEGIN_DATE", "", "END_DATE", "20/11/2027")})
	if err != nil {
		t.Fatal(err)
	}
	items := h.Baris[0].Items
	if len(items) != 7 || t8(items[6].TSI) != "1000" || len(items[6].Coverages) != 5 || !items[0].IsAdjustable || t8(items[0].PctAdjustOther) != "80" {
		t.Fatalf("%d item, %+v", len(items), items[len(items)-1])
	}
	var old []string
	for _, cv := range items[6].Coverages {
		old = append(old, cv.OldID)
	}
	if !reflect.DeepEqual(old, []string{"FLEXAS", "4.1A CC", "4.3", "4.2 PRGBI", "OTHERS"}) {
		t.Errorf("%v", old)
	}
	g := kasus.General
	if g.ReffNumber != "UJI-SLIP-1" || g.StartDateTime != "20261019T050000.000 GMT" || g.EndDateTime != "20271120T050000.000 GMT" {
		t.Errorf("%+v", g)
	}
	// Prorate memakai periode BARU (19-10-2026 s.d. 20-11-2027), bukan periode tersimpan (setahun).
	baru, err1 := Prorate(Periode{Mulai: g.StartDateTime, Akhir: g.EndDateTime})
	lamaPro, err2 := Prorate(Periode{Mulai: "20260101T050000.000 GMT", Akhir: "20270101T050000.000 GMT"})
	if err1 != nil || err2 != nil || t8(items[0].Coverages[0].ProRatePercent) != t8(baru) || t8(baru) == t8(lamaPro) {
		t.Errorf("prorate %s, baru %s, lama %s", t8(items[0].Coverages[0].ProRatePercent), t8(baru), t8(lamaPro))
	}
}

// TestCsvObjekDitolak - 400 terkumpul (angka, tanggal, kunci, mata uang, coverage master), "Please Upload CSV", batas 500,
// penjaga PUT …/objek (uw 409, tanpa identitas 401); tanpa tulisan apa pun.
func TestCsvObjekDitolak(t *testing.T) {
	ctx := context.Background()
	banyak := make([]map[string]string, BatasBarisCsvObjek+1)
	for i := range banyak {
		banyak[i] = barisCsvUji()
	}
	tanpaCov := csvMasterUji()
	tanpaCov.coverage["100815"] = false
	for nama, c := range map[string]struct {
		baris  []map[string]string
		master *csvTiruan
		pelaku inti.Pelaku
		mau    []string
	}{
		"kosong":         {nil, nil, pelakuCsv, []string{"Please Upload CSV"}},
		"terlalu banyak": {banyak, nil, pelakuCsv, []string{"paling banyak 3000 baris CSV"}},
		"isian": {[]map[string]string{barisCsvUji("TSI_OBJECT_ITEM", "satu juta", "BEGIN_DATE", "2026-10-19"),
			barisCsvUji("KOLOM_ASING", "x", "TSI_OBJECT_ITEM_6", "1,2,3")}, nil, pelakuCsv,
			[]string{"baris 1 kolom TSI_OBJECT_ITEM", "baris 1 kolom BEGIN_DATE", "baris 2: kolom KOLOM_ASING tidak dikenal",
				"baris 2 kolom TSI_OBJECT_ITEM_6"}},
		"mata uang": {[]map[string]string{barisCsvUji("CURRENCY", "XYZ")}, nil, pelakuCsv,
			[]string{`baris 1 kolom CURRENCY: mata uang "XYZ" tidak ada di master CURRENCY`}},
		"coverage":  {[]map[string]string{barisCsvUji()}, tanpaCov, pelakuCsv, []string{"coverage 100815 (FLEXAS) tidak ada di master COVERAGE"}},
		"identitas": {[]map[string]string{barisCsvUji()}, nil, inti.Pelaku{}, nil},
	} {
		m := c.master
		if m == nil {
			m = csvMasterUji()
		}
		p := &penambahCsvTiruan{}
		svc, kasus, ob, _ := layananCsvUji(m, p)
		_, err := svc.SimpanCsvObjek(ctx, c.pelaku, "UJI-NB-1", c.baris)
		if c.mau == nil {
			if !errors.Is(err, inti.ErrTanpaIdentitas) {
				t.Errorf("%s: %v", nama, err)
			}
		} else if !errors.Is(err, ErrMasukanObjek) {
			t.Errorf("%s: %v", nama, err)
		}
		for _, s := range c.mau {
			if err == nil || !strings.Contains(err.Error(), s) {
				t.Errorf("%s: tanpa %q: %v", nama, s, err)
			}
		}
		if ob.ada["UJI-NB-1"][0].ObjectNo != "LAMA" || kasus.General.ReffNumber != "UJI-SLIP-LAMA" || len(m.hapus) != 0 || len(p.masukan) != 0 {
			t.Errorf("%s: ada tulisan", nama)
		}
	}
	// Penjaga: case di tahap Underwriting -> sama dengan PUT …/objek (bukan tahap ubah).
	svc, _, ob, _ := layananCsvUji(csvMasterUji(), &penambahCsvTiruan{})
	svc.DenganTahap(tahapTiruan{"UJI-NB-1": {Position: "ReasFacInUnderwriting", FlagOnGoingPolicy: "0"}})
	if _, err := svc.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{barisCsvUji()}); err == nil ||
		ob.ada["UJI-NB-1"][0].ObjectNo != "LAMA" {
		t.Errorf("penjaga uw: %v", err)
	}
}

// TestJenisObjekBaku - D2 (b): tak peka huruf ke 9 nilai baku, selain itu "Others"; Keyword akumulasi per jenis.
func TestJenisObjekBaku(t *testing.T) {
	for in, mau := range map[string]string{"DWELLING HOUSE": "Dwelling House", "shop houses": "Shop Houses", " Factory ": "Factory",
		"PUBLIC WAREHOUSE": "Public Warehouse", "RUMAH": "Others", "": "Others"} {
		if got := jenisObjekBaku(in); got != mau {
			t.Errorf("%q: %q", in, got)
		}
	}
	if keywordAkumulasi("Private Warehouse") != "KAWASAN PERGUDANGAN" || keywordAkumulasi("Others") != "ZONA LAIN-LAIN" {
		t.Error("keyword")
	}
}

// TestKolomCsvObjekSamaDenganFrontend - 61 kolom, urutan sama dengan KOLOM_CSV_OBJEK (frontend/components/csvObjek.ts); saat
// ditulis juga sama dengan judul D:\migrasi\RNM\DDL\templatecsvfire.csv (61 = 61 = 61 sesudah 6 kolom wilayah 08-10-2026, selisih dua arah kosong).
func TestKolomCsvObjekSamaDenganFrontend(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "..", "frontend", "components", "csvObjek.ts"))
	if err != nil {
		t.Fatal(err)
	}
	blok := regexp.MustCompile(`(?s)KOLOM_CSV_OBJEK = \[(.*?)\]`).FindStringSubmatch(string(isi))
	if blok == nil {
		t.Fatal("KOLOM_CSV_OBJEK tidak ditemukan")
	}
	var ts []string
	for _, m := range regexp.MustCompile(`'([A-Z_0-9]+)'`).FindAllStringSubmatch(blok[1], -1) {
		ts = append(ts, m[1])
	}
	if len(KolomCsvObjek) != 61 || !reflect.DeepEqual(ts, KolomCsvObjek) {
		t.Errorf("go %d, ts %d", len(KolomCsvObjek), len(ts))
	}
}

// TestCsvObjekKodePosDanWilayah - tiket 66 lanjutan (keputusan work owner 08-10-2026): POSTAL_CODE dulu; gagal -> KODE_POS;
// keduanya gagal -> pesan verbatim ber-POSTAL_CODE; kode pos yang menemukan RW dipakai lokasi / risk address / akumulasi; > 1 RW
// -> dipilih menurut KECAMATAN / KABUPATEN_KOTA / PROVINSI (tak peka huruf), tidak cocok -> baris pertama; SUMBER_KODE_POS /
// TINGKAT_KEYAKINAN tidak memengaruhi.
func TestCsvObjekKodePosDanWilayah(t *testing.T) {
	ctx := context.Background()
	jalan := func(c *csvTiruan, p *penambahCsvTiruan, baris map[string]string) (HasilCsvObjek, *models.AlamatBaru) {
		t.Helper()
		svc, _, _, alamat := layananCsvUji(c, p)
		h, err := svc.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{baris})
		if err != nil {
			t.Fatal(err)
		}
		return h, alamat
	}
	// 1. POSTAL_CODE ketemu: KODE_POS tidak dicoba; kolom tampil-saja diabaikan.
	c, p := csvMasterUji(), &penambahCsvTiruan{}
	h, _ := jalan(c, p, barisCsvUji("KODE_POS", "88888", "SUMBER_KODE_POS", "ALAMAT", "TINGKAT_KEYAKINAN", "TINGGI"))
	if !reflect.DeepEqual(h.Pesan, []string{"Success Insert Data"}) || h.Baris[0].ZipCode != "12345" ||
		!reflect.DeepEqual(c.panggil, []string{"12345"}) {
		t.Errorf("1: %+v %v", h, c.panggil)
	}
	// 2. POSTAL_CODE gagal, KODE_POS ketemu: lokasi / risk address / akumulasi memakai KODE_POS.
	c, p = csvMasterUji(), &penambahCsvTiruan{}
	h, alamat := jalan(c, p, barisCsvUji("POSTAL_CODE", "99999", "KODE_POS", "12345"))
	if !reflect.DeepEqual(h.Pesan, []string{"Success Insert Data"}) || h.Baris[0].ZipCode != "12345" || h.Baris[0].City != "UJI KOTA" ||
		alamat.PostalCode != "12345" || len(p.masukan) != 1 || p.masukan[0]["zipCode"] != "12345" ||
		!reflect.DeepEqual(c.panggil, []string{"99999", "12345"}) {
		t.Errorf("2: %+v %+v %v %v", h, alamat, p.masukan, c.panggil)
	}
	// 2b. Risk address yang sudah ada dicari dengan kode pos yang menemukan RW (KODE_POS).
	c = csvMasterUji()
	c.risk["UJI RAYA|12345"], c.risk["UJI RAYA|99999"] = "UJI-RISK-KP", "UJI-RISK-POS"
	h, _ = jalan(c, &penambahCsvTiruan{}, barisCsvUji("POSTAL_CODE", "99999", "KODE_POS", "12345"))
	if h.Baris[0].RiskAddressID != "UJI-RISK-KP" {
		t.Errorf("2b: %q", h.Baris[0].RiskAddressID)
	}
	// 3. Keduanya gagal: pesan verbatim memakai POSTAL_CODE, tanpa master baru.
	c, p = csvMasterUji(), &penambahCsvTiruan{}
	h, alamat = jalan(c, p, barisCsvUji("POSTAL_CODE", "99999", "KODE_POS", "88888"))
	if !reflect.DeepEqual(h.Pesan, []string{"KodePos 99999 and kelurahan UJI KEL pada lokasi ke - 1 tidak ditemukan!"}) ||
		h.Baris[0].ZipCode != "99999" || alamat.PostalCode != "" || len(p.masukan) != 0 {
		t.Errorf("3: %+v", h)
	}
	// 3b. KODE_POS sama dengan POSTAL_CODE: satu pencarian saja.
	c = csvMasterUji()
	jalan(c, &penambahCsvTiruan{}, barisCsvUji("POSTAL_CODE", "99999", "KODE_POS", "99999"))
	if !reflect.DeepEqual(c.panggil, []string{"99999"}) {
		t.Errorf("3b: %v", c.panggil)
	}
	// 4. Banyak RW: dipilih menurut wilayah terisi; tidak ada yang cocok -> baris pertama.
	dua := []models.RWCsv{{Note: "UJI KEL", District: "LAIN KEC", City: "LAIN KOTA", Province: "UJI PROV", Nation: "UJI NEGARA"},
		{Note: "UJI KEL", District: "UJI KEC", City: "UJI KOTA", Province: "UJI PROV", Nation: "UJI NEGARA"}}
	for kec, mau := range map[string]string{" uji kec ": "UJI KEC", "TIDAK ADA": "LAIN KEC", "": "LAIN KEC"} {
		c = csvMasterUji()
		c.rw["12345|UJI KEL"] = dua
		h, _ = jalan(c, &penambahCsvTiruan{}, barisCsvUji("KECAMATAN", kec, "KABUPATEN_KOTA", "", "PROVINSI", "uji prov"))
		if h.Baris[0].District != mau {
			t.Errorf("4 %q: %q", kec, h.Baris[0].District)
		}
	}
	c = csvMasterUji()
	c.rw["12345|UJI KEL"] = dua
	h, _ = jalan(c, &penambahCsvTiruan{}, barisCsvUji("KECAMATAN", "UJI KEC", "KABUPATEN_KOTA", "LAIN KOTA"))
	if h.Baris[0].District != "LAIN KEC" { // semua kolom terisi wajib cocok: tidak ada -> baris pertama
		t.Errorf("4 semua: %q", h.Baris[0].District)
	}
}

// rwCsvTiruan / penambahDbTiruan - penulis master yang MENYIMPAN ke csvTiruan seperti basis data (risk address / akumulasi baru
// ditemukan pencarian berikutnya), untuk uji kesamaan hasil dengan / tanpa cache.
type rwCsvTiruan struct {
	c     *csvTiruan
	sisip *int
}

func (r rwCsvTiruan) CariRW(context.Context, string) ([]models.BarisRW, error) { return nil, nil }
func (r rwCsvTiruan) SisipRiskAddress(_ context.Context, _ *db.Tx, a models.AlamatBaru) (string, error) {
	*r.sisip++
	id := fmt.Sprintf("UJI-RA-%d", *r.sisip)
	r.c.risk[a.Address+"|"+a.PostalCode] = id
	return id, nil
}

type penambahDbTiruan struct {
	c     *csvTiruan
	tulis int
}

func (p *penambahDbTiruan) Tambah(_ context.Context, _ inti.Pelaku, _ string, m map[string]string) (string, error) {
	p.tulis++
	id := fmt.Sprintf("UJI-AK-%d", p.tulis)
	p.c.akum[strings.ToUpper(m["note"])] = id
	return id, nil
}

// barisCsvBanyak - n baris; alamat ke-i = i % alamat (alamat berbeda, kode pos / kelurahan sama - RW ditemukan semua).
func barisCsvBanyak(n, alamat int) []map[string]string {
	out := make([]map[string]string, n)
	for i := range out {
		out[i] = barisCsvUji("NO", fmt.Sprint(i+1), "ADDRESS", fmt.Sprintf("uji raya %d", i%alamat))
	}
	return out
}

// jalankanCsvBanyak - satu Save Data atas master "basis data" tiruan segar; kembalikan hasil, objek tersimpan, penghitung.
func jalankanCsvBanyak(t *testing.T, baris []map[string]string, cache bool) (HasilCsvObjek, []models.ObjekFire, *csvTiruan, int, int) {
	t.Helper()
	c := csvMasterUji()
	sisip, p := 0, &penambahDbTiruan{c: c}
	svc, _, ob, _ := layananCsvUji(c, nil)
	svc.DenganRW(rwCsvTiruan{c: c, sisip: &sisip}).DenganAddAkumulasi(addCsvTiruan{}, p)
	h, err := svc.simpanCsvObjek(context.Background(), pelakuCsv, "UJI-NB-1", baris, cache)
	if err != nil {
		t.Fatal(err)
	}
	return h, ob.ada["UJI-NB-1"], c, sisip, p.tulis
}

// TestCsvObjekCacheSamaHasil - cache per permintaan: hasil, objek tersimpan, dan master baru SAMA dengan tanpa cache; alamat
// berulang menambah risk address / akumulasi SEKALI (baris berikutnya memakai yang baru - seperti Pega); query jauh lebih sedikit.
func TestCsvObjekCacheSamaHasil(t *testing.T) {
	baris := barisCsvBanyak(40, 5)
	h1, o1, c1, rs1, ak1 := jalankanCsvBanyak(t, baris, true)
	h2, o2, c2, rs2, ak2 := jalankanCsvBanyak(t, baris, false)
	if !reflect.DeepEqual(h1, h2) || !reflect.DeepEqual(o1, o2) || rs1 != rs2 || ak1 != ak2 {
		t.Fatalf("hasil berbeda: risk %d/%d akumulasi %d/%d", rs1, rs2, ak1, ak2)
	}
	if rs1 != 5 || ak1 != 5 || o1[0].RiskAddressID != o1[5].RiskAddressID || o1[0].RiskAddressID == o1[1].RiskAddressID ||
		o1[5].Items[0].Coverages[0].AccumulationCode != o1[0].Items[0].Coverages[0].AccumulationCode {
		t.Errorf("master ganda: risk %d akumulasi %d", rs1, ak1)
	}
	// Alamat sama, kode pos berbeda -> dua risk address (kunci cache = alamat + kode pos).
	beda := []map[string]string{barisCsvUji(), barisCsvUji("NO", "2", "POSTAL_CODE", "54321")}
	jalankanDua := func(cache bool) (HasilCsvObjek, int) {
		c := csvMasterUji()
		c.rw["54321|UJI KEL"] = c.rw["12345|UJI KEL"]
		sisip := 0
		svc, _, _, _ := layananCsvUji(c, nil)
		svc.DenganRW(rwCsvTiruan{c: c, sisip: &sisip}).DenganAddAkumulasi(addCsvTiruan{}, &penambahDbTiruan{c: c})
		h, err := svc.simpanCsvObjek(context.Background(), pelakuCsv, "UJI-NB-1", beda, cache)
		if err != nil {
			t.Fatal(err)
		}
		return h, sisip
	}
	hb1, sb1 := jalankanDua(true)
	hb2, sb2 := jalankanDua(false)
	if sb1 != 2 || sb2 != 2 || !reflect.DeepEqual(hb1, hb2) || hb1.Baris[0].RiskAddressID == hb1.Baris[1].RiskAddressID {
		t.Errorf("kode pos beda: risk %d / %d", sb1, sb2)
	}
	if c1.total() >= c2.total() || c1.n["RWKelurahan"] != 1 || c1.n["MataUang"] != 1 || c1.n["IDJenisItemFire"] != 1 ||
		c1.n["NamaOkupasiFire"] != 1 || c1.n["IDRiskAddress"] != 5 || c1.n["IDAkumulasi"] != 5 {
		t.Errorf("query cache %v, tanpa cache %v", c1.n, c2.n)
	}
}

// TestCsvObjekBatas3000 - RALAT D7: 3000 baris lolos (3001 ditolak di TestCsvObjekDitolak); jumlah query pencarian master untuk
// 1128 dan 3000 baris dicatat (log uji) - dengan cache vs tanpa cache, alamat berbeda semua (terburuk) dan 200 alamat.
func TestCsvObjekBatas3000(t *testing.T) {
	if testing.Short() {
		t.Skip("3000 baris")
	}
	for _, n := range []int{1128, BatasBarisCsvObjek} {
		for _, alamat := range []int{n, 200} {
			h, o, c, rs, ak := jalankanCsvBanyak(t, barisCsvBanyak(n, alamat), true)
			_, _, c2, _, _ := jalankanCsvBanyak(t, barisCsvBanyak(n, alamat), false)
			if len(h.Baris) != n || len(o) != n {
				t.Fatalf("%d baris: %d", n, len(o))
			}
			t.Logf("%d baris, %d alamat: query master cache %d (tanpa cache %d); master baru risk %d + akumulasi %d", n, alamat,
				c.total(), c2.total(), rs, ak)
		}
	}
}
