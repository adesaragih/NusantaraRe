package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// lookupAnekaCsvTiruan - BrowseOccupationANEKA_RD / BrowseCoverageFacIn_RD tiruan: Contains tak peka huruf, urutan larik = urutan
// RD (baris pertama dipakai importer); mencatat parameter BusinessName / BizCode.
type lookupAnekaCsvTiruan struct {
	okupasi  []models.PilihanOkupasiAneka
	coverage []models.BarisCoverage
	kelas    []string
	bisnis   []string
}

func (l *lookupAnekaCsvTiruan) CariKapal(context.Context, string) ([]models.BarisKapal, error) {
	return nil, nil
}
func (l *lookupAnekaCsvTiruan) DaftarTipeKendaraan(context.Context) ([]models.BarisTipeKendaraan, error) {
	return nil, nil
}
func (l *lookupAnekaCsvTiruan) IndemnityPct(context.Context, string, string) ([]*apd.Decimal, error) {
	return nil, nil
}
func (l *lookupAnekaCsvTiruan) CariOkupasiAneka(_ context.Context, kelas, kata string) ([]models.PilihanOkupasiAneka, error) {
	l.kelas = append(l.kelas, kelas)
	var hasil []models.PilihanOkupasiAneka
	for _, o := range l.okupasi {
		if strings.Contains(strings.ToUpper(o.Name), strings.ToUpper(kata)) {
			hasil = append(hasil, o)
		}
	}
	return hasil, nil
}
func (l *lookupAnekaCsvTiruan) CariCoverageAneka(_ context.Context, kode, kata string) ([]models.BarisCoverage, error) {
	l.bisnis = append(l.bisnis, kode)
	var hasil []models.BarisCoverage
	for _, c := range l.coverage {
		if strings.Contains(strings.ToUpper(c.Nama), strings.ToUpper(kata)) {
			hasil = append(hasil, c)
		}
	}
	return hasil, nil
}

// anekaCsvTiruan - PenulisAnekaCsv: ganti = isi lama dibuang; mencatat ganti / jumlah lokasi per kelompok.
type anekaCsvTiruan struct {
	ada    map[string][]models.LokasiAneka
	ganti  []bool
	ukuran []int
	// gagalKe - panggilan ke-n (1..) gagal dengan galat itu tanpa menulis (transaksi di-rollback).
	gagalKe int
	galat   error
	n       int
}

func (a *anekaCsvTiruan) TulisAnekaCsv(_ context.Context, _ *db.Tx, id string, lokasi []models.LokasiAneka, ganti bool) error {
	a.n++
	if a.n == a.gagalKe {
		return a.galat
	}
	if a.ada == nil {
		a.ada = map[string][]models.LokasiAneka{}
	}
	if ganti {
		a.ada[id] = nil
	}
	a.ada[id] = append(a.ada[id], lokasi...)
	a.ganti, a.ukuran = append(a.ganti, ganti), append(a.ukuran, len(lokasi))
	return nil
}

// layananUnggahAnekaUji - layananUnggahUji untuk case MARINE HULL (GroupPanel 003, OldID 05): okupasi "UJI OKUPASI KAPAL" (dan
// "UJI OKUPASI KAPAL LAIN" sesudahnya), coverage "UJI HULL" / "UJI HULL MACHINERY", RW kode pos 12345 tanpa kelurahan.
func layananUnggahAnekaUji() (*Service, *unggahTiruan, *anekaCsvTiruan, *lookupAnekaCsvTiruan, *csvTiruan) {
	svc, u, _, _, c := layananUnggahUji()
	svc.tableOfLimit = tolTiruan{bisnis: satuBisnis("003", "05")}
	c.rw["12345|"] = []models.RWCsv{{Note: "UJI RW", District: "UJI KEC", City: "UJI KOTA", Province: "UJI PROV"}}
	l := &lookupAnekaCsvTiruan{
		okupasi: []models.PilihanOkupasiAneka{{ID: "OK9", Name: "UJI OKUPASI KAPAL"}, {ID: "OK8", Name: "UJI OKUPASI KAPAL LAIN"}},
		coverage: []models.BarisCoverage{{ID: "C1", OldID: "1", Nama: "UJI HULL"}, {ID: "C2", OldID: "2", Nama: "UJI HULL MACHINERY"},
			{ID: "C3", OldID: "3", Nama: "UJI P&I"}}}
	a := &anekaCsvTiruan{}
	svc = svc.DenganLookupAneka(l).DenganUnggahCsvAneka(a)
	return svc, u, a, l, c
}

// barisAnekaUji - satu baris Marine Hull sah (per nama kolom).
func barisAnekaUji() map[string]string {
	return map[string]string{"NO": "1", "TYPE_ADDRESS": "JL.", "ADDRESS": "UJI DERMAGA", "BUILDING_NO": "7", "POSTAL_CODE": "12345",
		"COUNTRY": "UJI NEGARA", "OCCUPATION_NAME": "okupasi kapal", "OBJECT_DESCRIPTION": "UJI KAPAL A", "SECTION": "SECTION I",
		"CURRENCY": "IDR", "TSI": "1.000.000,50", "YEAR": "2001", "PLACE_OF_MANUFACTURE": "UJI GALANGAN", "OBJECT_TYPE": "TUG",
		"NRT": "10", "GRT": "20", "DWT": "30", "FLAGNAME": "UJI BENDERA", "CONSTRUCTION": "STEEL", "CLASSIFICATION": "UJI KELAS",
		"COVERAGE_NAME": "UJI HULL", "CALCULATE_METHOD": "PRO RATA", "DAY": "", "PCT_LOADING": "5", "PCT_SCALE": "100",
		"PCT_RATE": "0,25", "PREMIUM": "2.500,75", "PCT_DISCOUNT": "10", "DISCOUNT": "250", "PCT_LOL": "50", "LOL": "500.000"}
}

// judulAnekaUji - 31 judul template Marine Hull urutan TERBALIK + satu kolom tak dikenal di depan.
func judulAnekaUji() []string {
	j := slices.Clone(models.KolomCsvAneka[:models.JumlahJudulMarineHull])
	slices.Reverse(j)
	return append([]string{"KETERANGAN"}, j...)
}

// csvAnekaUji - berkas n baris (judulAnekaUji); ubah(i, baris) mengubah baris ke-i (1..n).
func csvAnekaUji(n int, ubah func(int, map[string]string)) string {
	judul := judulAnekaUji()
	var b strings.Builder
	b.WriteString("\uFEFF" + strings.Join(judul, ";") + "\r\n")
	for i := 1; i <= n; i++ {
		r := barisAnekaUji()
		r["KETERANGAN"] = "x"
		if ubah != nil {
			ubah(i, r)
		}
		sel := make([]string, len(judul))
		for k, j := range judul {
			sel[k] = r[j]
		}
		b.WriteString(strings.Join(sel, ";") + "\r\n")
	}
	return b.String()
}

// isiAnekaUji - sel urut models.KolomCsvAneka.
func isiAnekaUji(r map[string]string) []string {
	isi := make([]string, len(models.KolomCsvAneka))
	for i, k := range models.KolomCsvAneka {
		isi[i] = r[k]
	}
	return isi
}

// TestPetaJudulAneka - judul per NAMA, urutan bebas; 31 judul template wajib; SERIALNO / CHASSISNUMBER / ENGINENUMBER /
// KELURAHAN opsional; ganda / kurang -> 400 berpesan.
func TestPetaJudulAneka(t *testing.T) {
	p, err := petaJudulAneka(strings.Join(judulAnekaUji(), ";"))
	if err != nil || p.n != 35 || p.idx[0] != -1 || p.kolom[0] != "LOL" || len(p.kolom) != 31 || p.idx[1] != 30 {
		t.Fatalf("%v %+v", err, p)
	}
	if p, err := petaJudulAneka(strings.Join(models.KolomCsvAneka, ";")); err != nil || len(p.kolom) != 35 {
		t.Errorf("35 judul: %v %+v", err, p)
	}
	if _, err := petaJudulAneka(" no ;" + strings.Join(models.KolomCsvAneka[:31], ";")); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "NO dua kali") {
		t.Errorf("ganda: %v", err)
	}
	if _, err := petaJudulAneka("NO;TSI"); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "template Marine Hull - kolom TYPE_ADDRESS, ADDRESS") {
		t.Errorf("kurang: %v", err)
	}
}

// TestUraiBarisAnekaCsv - angka (ribuan "." dibuang lalu "," -> "."; persen "," -> "."), Section / Calculate Method / Day, lebar.
func TestUraiBarisAnekaCsv(t *testing.T) {
	r := barisAnekaUji()
	b, m := uraiBarisAnekaCsv(4, isiAnekaUji(r))
	if len(m) != 0 || utils.FormatDecimal(b.tsi) != "1000000.5" || utils.FormatDecimal(b.premium) != "2500.75" ||
		utils.FormatDecimal(b.pctRate) != "0.25" || utils.FormatDecimal(b.lol) != "500000" || b.section != "1" ||
		b.calculateMethod != "1" || b.day != "365" {
		t.Fatalf("%v %+v", m, b)
	}
	for masuk, mau := range map[[3]string][3]string{
		{"SECTION II", "SHORT PERIOD", "360"}: {"2", "2", "360"}, {"WARRANTY", "", "366"}: {"3", "3", "366"},
		{"section i", "pro rata", ""}: {"4", "3", "365"}, {"", "FIX RATE", "1"}: {"4", "3", "1"},
	} {
		r["SECTION"], r["CALCULATE_METHOD"], r["DAY"] = masuk[0], masuk[1], masuk[2]
		b, _ := uraiBarisAnekaCsv(1, isiAnekaUji(r))
		if [3]string{b.section, b.calculateMethod, b.day} != mau {
			t.Errorf("%v: %s %s %s", masuk, b.section, b.calculateMethod, b.day)
		}
	}
	// Pega: TSI "1.5" -> "15" (titik = ribuan) apa adanya; persen "1.5" tetap 1.5.
	r = barisAnekaUji()
	r["TSI"], r["PCT_RATE"] = "1.5", "1.5"
	if b, _ := uraiBarisAnekaCsv(1, isiAnekaUji(r)); utils.FormatDecimal(b.tsi) != "15" || utils.FormatDecimal(b.pctRate) != "1.5" {
		t.Errorf("titik: %v %v", b.tsi, b.pctRate)
	}
	r["PREMIUM"], r["YEAR"], r["ADDRESS"] = "satu", strings.Repeat("9", 51), strings.Repeat("A", lebarRiskAddress-3)
	_, m = uraiBarisAnekaCsv(7, isiAnekaUji(r))
	gabung := strings.Join(m, "|")
	if len(m) != 3 || !strings.Contains(gabung, `baris 7 kolom PREMIUM: angka tidak sah "satu"`) ||
		!strings.Contains(gabung, "baris 7 kolom YEAR: paling banyak 50 byte") || !strings.Contains(gabung, "Risk Location") {
		t.Errorf("%q", m)
	}
}

// TestProtectAnekaCsv - ProtectUploadAneka_Act: pesan verbatim, kosong menimpa, okupasi tidak ada -> keluar (coverage tidak
// diperiksa), Contains + baris pertama, BusinessName = Class of Business, BizCode = BUSINESS.ID; mata uang tidak ada = galat.
func TestProtectAnekaCsv(t *testing.T) {
	ctx := context.Background()
	svc, _, _, l, _ := layananUnggahAnekaUji()
	c, err := svc.cacheAnekaCsv(ctx, "UJI-NB-1")
	if err != nil {
		t.Fatal(err)
	}
	periksa := func(ubah map[string]string) (barisAnekaCsv, []string) {
		r := barisAnekaUji()
		for k, v := range ubah {
			r[k] = v
		}
		b, _ := uraiBarisAnekaCsv(3, isiAnekaUji(r))
		m, err := svc.periksaBarisAnekaCsv(ctx, &b, c)
		if err != nil {
			t.Fatal(err)
		}
		return b, m
	}
	b, m := periksa(nil)
	if len(m) != 0 || b.okupasi.ID != "OK9" || b.coverage.ID != "C1" || b.coverage.Nama != "UJI HULL" || b.mu.ID != "1" ||
		l.kelas[0] != "UJI KELAS" || l.bisnis[0] != "10999" {
		t.Fatalf("%q %+v kelas %v bisnis %v", m, b, l.kelas, l.bisnis)
	}
	for nama, c := range map[string]struct {
		ubah map[string]string
		mau  []string
	}{
		"okupasi kosong":   {map[string]string{"OCCUPATION_NAME": ""}, []string{"Occupation Cannot NULL"}},
		"coverage kosong":  {map[string]string{"COVERAGE_NAME": ""}, []string{"Coverage Cannot NULL"}},
		"keduanya kosong":  {map[string]string{"OCCUPATION_NAME": "", "COVERAGE_NAME": ""}, []string{"Coverage Cannot NULL"}},
		"okupasi tak ada":  {map[string]string{"OCCUPATION_NAME": "XX", "COVERAGE_NAME": "YY"}, []string{"* Occupation XX Not Available in Master"}},
		"coverage tak ada": {map[string]string{"COVERAGE_NAME": "YY"}, []string{"* Coverage YY Not Available in Master"}},
		"okupasi kosong, coverage tak ada": {map[string]string{"OCCUPATION_NAME": "", "COVERAGE_NAME": "YY"},
			[]string{"Occupation Cannot NULL", "* Coverage YY Not Available in Master"}},
		"mata uang tak ada": {map[string]string{"CURRENCY": "XYZ"},
			[]string{`baris 3 kolom CURRENCY: mata uang "XYZ" tidak ada di master CURRENCY`}},
		"mata uang kosong": {map[string]string{"CURRENCY": ""}, []string{`baris 3 kolom CURRENCY: mata uang "" tidak ada di master CURRENCY`}},
	} {
		if _, m := periksa(c.ubah); !slices.Equal(m, c.mau) {
			t.Errorf("%s: %q", nama, m)
		}
	}
	// Contains + baris pertama (K73-59): "machinery" -> C2; "hull" -> C1 (baris pertama dari dua).
	if b, _ := periksa(map[string]string{"COVERAGE_NAME": "machinery"}); b.coverage.ID != "C2" {
		t.Errorf("machinery: %+v", b.coverage)
	}
	if b, _ := periksa(map[string]string{"COVERAGE_NAME": "hull", "OCCUPATION_NAME": "LAIN"}); b.coverage.ID != "C1" || b.okupasi.ID != "OK8" {
		t.Errorf("hull: %+v %+v", b.coverage, b.okupasi)
	}
}

// TestSusunAnekaCsv - Insert 4.7 / 4.9 / 4.12 / 4.13: lokasi per NO, okupasi tunggal DITIMPA, objek per (deskripsi, section) dengan
// medan dari baris PERTAMA, coverage per nama (baris coverage sama berurutan diabaikan), coverage ber-mata uang objek.
func TestSusunAnekaCsv(t *testing.T) {
	baris := func(ubah map[string]string) barisAnekaCsv {
		r := barisAnekaUji()
		for k, v := range ubah {
			r[k] = v
		}
		b, m := uraiBarisAnekaCsv(1, isiAnekaUji(r))
		if len(m) != 0 {
			t.Fatal(m)
		}
		b.okupasi = models.PilihanOkupasiAneka{ID: "OK-" + r["OCCUPATION_NAME"], Name: strings.ToUpper(r["OCCUPATION_NAME"])}
		b.coverage = models.BarisCoverage{ID: "C-" + r["COVERAGE_NAME"], Nama: r["COVERAGE_NAME"]}
		return b
	}
	z := &susunAnekaCsv{}
	rw := &models.RWCsv{Note: "UJI RW", District: "UJI KEC", City: "UJI KOTA", Province: "UJI PROV"}
	z.terapkan(baris(nil), rw)
	z.terapkan(baris(map[string]string{"COVERAGE_NAME": "UJI MESIN", "TSI": "7", "PREMIUM": "5", "YEAR": "1999"}), nil)
	z.terapkan(baris(map[string]string{"COVERAGE_NAME": "UJI MESIN", "PREMIUM": "9"}), nil) // coverage sama: diabaikan
	z.terapkan(baris(map[string]string{"SECTION": "SECTION II", "CURRENCY": "USD", "OCCUPATION_NAME": "lain"}), nil)
	z.terapkan(baris(map[string]string{"NO": "2", "TYPE_ADDRESS": "", "ADDRESS": "UJI DUA", "OBJECT_DESCRIPTION": "UJI KAPAL B"}), nil)
	if len(z.lokasi) != 2 {
		t.Fatalf("%d lokasi", len(z.lokasi))
	}
	l := z.lokasi[0]
	if l.ObjectNo != "1" || l.RoadType != "JL." || l.RoadName != "UJI DERMAGA" || l.RiskLocation != "JL. UJI DERMAGA" ||
		l.City != "UJI KOTA" || l.Province != "UJI PROV" || l.District != "UJI KEC" || l.Territory != "UJI RW" ||
		l.ZipCode != "12345" || l.Country != "UJI NEGARA" || len(l.Okupasi) != 1 || l.Okupasi[0].OccupationID != "OK-lain" {
		t.Errorf("lokasi 1: %+v", l)
	}
	ob := l.Okupasi[0].Objek
	if len(ob) != 2 || ob[0].Year != "2001" || ob[0].Section != "1" || utils.FormatDecimal(ob[0].TSI) != "1000000.5" ||
		ob[0].MarineHull.TypeOfVessel != "TUG" || ob[0].MarineHull.NRTRT != "10" || ob[0].Flag != "UJI BENDERA" ||
		ob[0].BuiltIn != "UJI GALANGAN" || ob[1].Section != "2" || ob[1].Currency != "USD" {
		t.Errorf("objek: %+v", ob)
	}
	cov := ob[0].Coverages
	if len(cov) != 2 || cov[0].Coverage != "C-UJI HULL" || cov[0].CoverageNote != "UJI HULL" || cov[1].Coverage != "C-UJI MESIN" ||
		utils.FormatDecimal(cov[1].TSI) != "7" || utils.FormatDecimal(cov[1].Premium) != "5" || cov[0].Currency != "IDR" ||
		cov[0].CalculateMethodFacIn != "1" || cov[0].Day != "365" || utils.FormatDecimal(cov[0].Rate) != "0.25" ||
		utils.FormatDecimal(cov[0].LostLimit) != "100" || utils.FormatDecimal(cov[0].Loading) != "5" ||
		utils.FormatDecimal(cov[0].DiscountPercentage) != "10" || utils.FormatDecimal(cov[0].Discount) != "250" ||
		utils.FormatDecimal(cov[0].PctLoL) != "50" || utils.FormatDecimal(cov[0].LimitOfLiability) != "500000" ||
		cov[0].TSINusantaraRe != nil || cov[0].PremiNusantaraRe != nil || cov[0].MinPremium != nil {
		t.Errorf("coverage: %+v", cov)
	}
	if c := ob[1].Coverages; len(c) != 1 || c[0].Currency != "USD" {
		t.Errorf("coverage objek 2: %+v", c)
	}
	if l2 := z.lokasi[1]; l2.RiskLocation != " UJI DUA" || l2.City != "" || l2.Okupasi[0].OccupationID != "OK-okupasi kapal" ||
		len(l2.Okupasi[0].Objek) != 1 {
		t.Errorf("lokasi 2: %+v", l2)
	}
}

// TestUnggahCsvAnekaAlur - unggah Marine Hull 2500 baris (dua baris per lokasi, berbeda coverage): jenis aneka, validasi, simpan
// berkelompok tanpa memotong lokasi (kelompok 1 ganti + CurrencyList dihapus sekali), baris NO kosong dilewati.
func TestUnggahCsvAnekaAlur(t *testing.T) {
	ctx := context.Background()
	svc, _, a, _, c := layananUnggahAnekaUji()
	berkas := csvAnekaUji(2500, func(i int, r map[string]string) {
		r["NO"] = fmt.Sprint((i + 1) / 2)
		if i%2 == 0 {
			r["COVERAGE_NAME"] = "MACHINERY"
		}
		if i == 2500 {
			r["NO"] = ""
		}
	})
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
	if err != nil {
		t.Fatal(err)
	}
	x := st.Unggah
	if x.Status != StatusUnggahValid || x.Jenis != models.JenisCsvAneka || x.TotalBaris != 2500 || x.JumlahGalat != 0 ||
		!slices.Equal(x.Kolom, judulAnekaUji()[1:]) {
		t.Fatalf("validasi: %+v %v", x, st.Galat)
	}
	st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", x.ID)
	if err != nil || st.Unggah.Status != StatusUnggahSelesai || st.Unggah.BarisTersimpan != 2500 {
		t.Fatalf("simpan: %v %+v", err, st.Unggah)
	}
	lokasi := a.ada["UJI-NB-1"]
	if len(lokasi) != 1250 || !slices.Equal(a.ganti, []bool{true, false, false}) || !slices.Equal(a.ukuran, []int{500, 500, 250}) ||
		len(c.hapus) != 1 {
		t.Fatalf("lokasi %d ganti %v ukuran %v hapus %d", len(lokasi), a.ganti, a.ukuran, len(c.hapus))
	}
	if l := lokasi[1249]; l.ObjectNo != "1250" || len(l.Okupasi[0].Objek[0].Coverages) != 1 {
		t.Errorf("lokasi terakhir (baris 2500 NO kosong dilewati): %+v", l)
	}
	if cov := lokasi[0].Okupasi[0].Objek[0].Coverages; len(cov) != 2 || cov[1].Coverage != "C2" || lokasi[0].City != "UJI KOTA" {
		t.Errorf("lokasi 1: %+v", lokasi[0])
	}
}

// TestUnggahCsvAnekaBerhentiGalatPertama - K73-62: validasi berhenti di baris PERTAMA bermasalah (pesan baris itu saja); simpan
// ditolak; K73-53: NO berulang tidak berurutan = galat.
func TestUnggahCsvAnekaBerhentiGalatPertama(t *testing.T) {
	ctx := context.Background()
	svc, u, a, _, _ := layananUnggahAnekaUji()
	berkas := csvAnekaUji(10, func(i int, r map[string]string) {
		r["NO"] = fmt.Sprint(i)
		switch i {
		case 4:
			r["OCCUPATION_NAME"], r["COVERAGE_NAME"] = "", ""
		case 7:
			r["CURRENCY"] = "XYZ"
		}
	})
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
	if err != nil {
		t.Fatal(err)
	}
	x := st.Unggah
	if x.Status != StatusUnggahGagalValidasi || x.JumlahGalat != 1 || x.BarisDiproses != 4 || len(st.Galat) != 1 ||
		st.Galat[0].NoBaris != 4 || st.Galat[0].Teks != "Coverage Cannot NULL" {
		t.Fatalf("%+v %+v", x, st.Galat)
	}
	if _, err := svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", x.ID); !errors.Is(err, ErrUnggahStatus) || len(a.ganti) != 0 {
		t.Errorf("simpan: %v %v", err, a.ganti)
	}
	berkas = csvAnekaUji(4, func(i int, r map[string]string) { r["NO"] = []string{"1", "2", "1", "3"}[i-1] })
	st, err = svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
	if err != nil || st.Unggah.Status != StatusUnggahGagalValidasi || len(st.Galat) != 1 || st.Galat[0].NoBaris != 3 ||
		!strings.Contains(st.Galat[0].Teks, `lokasi NO "1" sudah muncul sebelumnya`) {
		t.Errorf("NO berulang: %v %+v %+v", err, st.Unggah, st.Galat)
	}
	_ = u
}

// TestUnggahCsvAnekaKunci - K75-2: selama unggah Aneka berjalan, PUT objek-aneka 409.
func TestUnggahCsvAnekaKunci(t *testing.T) {
	ctx := context.Background()
	svc, u, _, _, _ := layananUnggahAnekaUji()
	svc = svc.DenganObjekAneka(anekaTiruan{ada: map[string][]models.LokasiAneka{"UJI-NB-1": nil}})
	u.BuatUnggah(ctx, nil, models.UnggahCsv{CaseID: "UJI-NB-1", Status: StatusUnggahMenyimpan, Fase: FaseUnggahSimpan, TotalBaris: 3,
		Jenis: models.JenisCsvAneka})
	if _, err := svc.GantiObjekAneka(ctx, pelakuCsv, "UJI-NB-1", nil, nil); !errors.Is(err, ErrUnggahBerjalan) {
		t.Errorf("PUT objek-aneka: %v", err)
	}
}

// TestKelompokAnekaTidakMemotongLokasi - lokasi tiga baris (batas kelompok 1000 jatuh di tengah lokasi 334): kelompok diperpanjang
// sampai lokasi selesai, tidak ada lokasi terbelah (ObjectNo unik, tiap lokasi 3 coverage).
func TestKelompokAnekaTidakMemotongLokasi(t *testing.T) {
	ctx := context.Background()
	svc, _, a, _, _ := layananUnggahAnekaUji()
	svc.lookupAneka.(*lookupAnekaCsvTiruan).coverage = append(svc.lookupAneka.(*lookupAnekaCsvTiruan).coverage,
		models.BarisCoverage{ID: "C4", Nama: "UJI TIGA"})
	berkas := csvAnekaUji(2500, func(i int, r map[string]string) {
		r["NO"] = fmt.Sprint((i + 2) / 3)
		r["COVERAGE_NAME"] = []string{"UJI HULL", "MACHINERY", "UJI TIGA"}[(i-1)%3]
	})
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
	if err != nil || st.Unggah.Status != StatusUnggahValid {
		t.Fatalf("%v %+v %v", err, st.Unggah, st.Galat)
	}
	if st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", st.Unggah.ID); err != nil || st.Unggah.Status != StatusUnggahSelesai {
		t.Fatalf("simpan: %v %+v", err, st.Unggah)
	}
	lokasi := a.ada["UJI-NB-1"]
	lihat := map[string]bool{}
	for i, l := range lokasi {
		n := len(l.Okupasi[0].Objek[0].Coverages)
		if lihat[l.ObjectNo] || (i < len(lokasi)-1 && n != 3) {
			t.Fatalf("lokasi %s terbelah / %d coverage", l.ObjectNo, n)
		}
		lihat[l.ObjectNo] = true
	}
	// kelompok = baris 1-1002, 1003-2004, 2005-2500 (diperpanjang sampai lokasi selesai).
	if len(lokasi) != 834 || !slices.Equal(a.ukuran, []int{334, 334, 166}) {
		t.Errorf("%d lokasi, kelompok %v", len(lokasi), a.ukuran)
	}
}

// TestBarisNoKosongAnekaDilewati - K73-54: baris NO kosong tidak diperiksa (isinya rusak pun) dan tidak memutus perbandingan
// baris sebelumnya: baris 3 (NO 1, sesudah baris kosong) tetap lokasi 1.
func TestBarisNoKosongAnekaDilewati(t *testing.T) {
	ctx := context.Background()
	svc, _, a, _, _ := layananUnggahAnekaUji()
	berkas := csvAnekaUji(4, func(i int, r map[string]string) {
		r["NO"] = []string{"1", "", "1", "2"}[i-1]
		if i == 2 {
			r["CURRENCY"], r["TSI"], r["OCCUPATION_NAME"] = "XYZ", "rusak", ""
		}
		if i == 3 {
			r["COVERAGE_NAME"] = "MACHINERY"
		}
	})
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
	if err != nil || st.Unggah.Status != StatusUnggahValid {
		t.Fatalf("%v %+v %v", err, st.Unggah, st.Galat)
	}
	if st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", st.Unggah.ID); err != nil || st.Unggah.Status != StatusUnggahSelesai {
		t.Fatalf("simpan: %v %+v", err, st.Unggah)
	}
	if l := a.ada["UJI-NB-1"]; len(l) != 2 || len(l[0].Okupasi[0].Objek[0].Coverages) != 2 {
		t.Errorf("%+v", l)
	}
}

// TestUlangiSimpanAnekaSesudahGagal - bug DEV 11-10-2026 (ORA-12899 di kelompok 1): pekerjaan 'gagal' fase simpan, BARIS_TERSIMPAN
// tetap 0 (kemajuan ikut rollback); Ulangi melanjutkan fase simpan dari baris 1 dengan GANTI lagi, tanpa lokasi ganda. Gagal di
// kelompok 2 -> Ulangi mulai di kelompok 2 tanpa ganti.
func TestUlangiSimpanAnekaSesudahGagal(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		gagalKe, tersimpan int
		ganti              []bool
	}{{1, 0, []bool{true, false, false}}, {2, 1000, []bool{true, false, false}}} {
		svc, _, a, _, _ := layananUnggahAnekaUji()
		a.gagalKe, a.galat = c.gagalKe, errors.New("ORA-12899: value too large for column (uji)")
		berkas := csvAnekaUji(2500, func(i int, r map[string]string) { r["NO"] = fmt.Sprint((i + 1) / 2) })
		st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
		if err != nil {
			t.Fatal(err)
		}
		st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", st.Unggah.ID)
		x := st.Unggah
		if err != nil || x.Status != StatusUnggahGagal || x.Fase != FaseUnggahSimpan || x.BarisTersimpan != c.tersimpan ||
			!strings.Contains(x.GalatProses, "ORA-12899") {
			t.Fatalf("gagal ke-%d: %v %+v", c.gagalKe, err, x)
		}
		st, err = svc.UlangiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", x.ID)
		if err != nil || st.Unggah.Status != StatusUnggahSelesai || st.Unggah.BarisTersimpan != 2500 {
			t.Fatalf("ulangi ke-%d: %v %+v", c.gagalKe, err, st.Unggah)
		}
		if l := a.ada["UJI-NB-1"]; len(l) != 1250 || l[0].ObjectNo != "1" || l[1249].ObjectNo != "1250" || !slices.Equal(a.ganti, c.ganti) {
			t.Errorf("gagal ke-%d: %d lokasi, ganti %v", c.gagalKe, len(l), a.ganti)
		}
	}
}

// TestLebarRWAnekaCsv - audit lebar 11-10-2026 (keputusan work owner): RW master melebihi kolom lokasi 100 byte (BYTE = len UTF-8)
// -> validasi berhenti di baris lokasi itu dengan galat jelas; tepat 100 byte sah; RW berubah memanjang sesudah validasi -> simpan
// gagal "data master berubah".
func TestLebarRWAnekaCsv(t *testing.T) {
	ctx := context.Background()
	if m := periksaRWAnekaCsv(1, "12345", &models.RWCsv{Note: strings.Repeat("é", 50), District: strings.Repeat("D", 100)}); len(m) != 0 {
		t.Fatalf("tepat 100 byte: %q", m)
	}
	for medan, rw := range map[string]models.RWCsv{
		"NOTE 51 rune 102 byte melebihi kolom T_RISKLOCATION.ASMRW":        {Note: strings.Repeat("é", 51)},
		"DISTRICTNAME 101 byte melebihi kolom T_RISKLOCATION.ASM_DISTRICT": {District: strings.Repeat("D", 101)},
		"CITYNAME 101 byte melebihi kolom T_RISKLOCATION.ASM_CITY":         {City: strings.Repeat("K", 101)},
		"PROVINCENAME 101 byte melebihi kolom T_PROPERTY.PROVINCE":         {Province: strings.Repeat("P", 101)},
	} {
		mau := strings.Replace(medan, " 51 rune", "", 1)
		if m := periksaRWAnekaCsv(4, "12345", &rw); len(m) != 1 || m[0] != `baris 4: master RW kode pos "12345" - `+mau+" (paling banyak 100 byte)" {
			t.Errorf("%s: %q", medan, m)
		}
	}
	svc, _, a, _, c := layananUnggahAnekaUji()
	c.rw["99999|"] = []models.RWCsv{{City: strings.Repeat("é", 51)}}
	berkas := csvAnekaUji(4, func(i int, r map[string]string) {
		r["NO"] = fmt.Sprint((i + 1) / 2)
		if i >= 3 {
			r["POSTAL_CODE"] = "99999"
		}
	})
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
	if err != nil || st.Unggah.Status != StatusUnggahGagalValidasi || len(st.Galat) != 1 || st.Galat[0].NoBaris != 3 ||
		!strings.Contains(st.Galat[0].Teks, "CITYNAME 102 byte melebihi kolom T_RISKLOCATION.ASM_CITY") {
		t.Fatalf("validasi: %v %+v %+v", err, st.Unggah, st.Galat)
	}
	// RW sah saat validasi, memanjang sebelum simpan -> simpan gagal, objek tidak ditulis.
	c.rw["99999|"] = []models.RWCsv{{City: "UJI KOTA"}}
	st, err = svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
	if err != nil || st.Unggah.Status != StatusUnggahValid {
		t.Fatalf("validasi 2: %v %+v %v", err, st.Unggah, st.Galat)
	}
	c.rw["99999|"] = []models.RWCsv{{Province: strings.Repeat("P", 101)}}
	// cache RW per pekerjaan dibangun ulang saat simpan.
	st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", st.Unggah.ID)
	if err != nil || st.Unggah.Status != StatusUnggahGagal || !strings.Contains(st.Unggah.GalatProses, "data master berubah sejak validasi") ||
		!strings.Contains(st.Unggah.GalatProses, "PROVINCENAME 101 byte") || len(a.ganti) != 0 {
		t.Errorf("simpan: %v %+v ganti %v", err, st.Unggah, a.ganti)
	}
}

// TestPremiMelebihiTSIAneka - K73-63 (penyimpangan sadar dari Pega): PREMIUM > TSI baris yang sama = galat bernomor baris; sama =
// sah; selisih 0,01 = galat; bentuk CSV "43.200.000.002,00" vs "100.000.001,00" (baris contoh template) = galat; validasi berhenti
// di baris itu (K73-62); salah satu kosong = tanpa cek.
func TestPremiMelebihiTSIAneka(t *testing.T) {
	for nama, c := range map[string]struct{ premi, tsi, mau string }{
		"sama":            {"1.000.000,50", "1.000.000,50", ""},
		"lebih 0,01":      {"1.000.000,51", "1.000.000,50", "baris 5: PREMIUM 1000000.51 melebihi TSI 1000000.5"},
		"contoh template": {"43.200.000.002,00", "100.000.001,00", "baris 5: PREMIUM 43200000002 melebihi TSI 100000001"},
		"kurang":          {"2.500,75", "1.000.000,50", ""},
		"premi kosong":    {"", "1.000", ""},
		"TSI kosong":      {"1.000", "", ""},
	} {
		r := barisAnekaUji()
		r["PREMIUM"], r["TSI"] = c.premi, c.tsi
		_, m := uraiBarisAnekaCsv(5, isiAnekaUji(r))
		if (c.mau == "" && len(m) != 0) || (c.mau != "" && (len(m) != 1 || m[0] != c.mau)) {
			t.Errorf("%s: %q", nama, m)
		}
	}
	ctx := context.Background()
	svc, _, a, _, _ := layananUnggahAnekaUji()
	berkas := csvAnekaUji(4, func(i int, r map[string]string) {
		r["NO"] = fmt.Sprint(i)
		if i == 3 {
			r["PREMIUM"], r["TSI"] = "43.200.000.002,00", "100.000.001,00"
		}
	})
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "hull.csv", strings.NewReader(berkas))
	if err != nil || st.Unggah.Status != StatusUnggahGagalValidasi || len(st.Galat) != 1 || st.Galat[0].NoBaris != 3 ||
		st.Galat[0].Teks != "baris 3: PREMIUM 43200000002 melebihi TSI 100000001" || len(a.ganti) != 0 {
		t.Errorf("validasi: %v %+v %+v", err, st.Unggah, st.Galat)
	}
}
