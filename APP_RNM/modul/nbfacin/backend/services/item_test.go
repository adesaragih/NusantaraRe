package services

import (
	"context"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

type pilihanItemTiruan struct{ panggil *int }

func (p pilihanItemTiruan) DaftarJenisItem(context.Context) ([]models.JenisItem, error) {
	return []models.JenisItem{{Kode: "UJI01", Nama: "UJI MESIN", Keterangan: "UJI KET"}}, nil
}

func (p pilihanItemTiruan) DaftarMataUang(context.Context) ([]string, error) {
	*p.panggil++
	return []string{"IDR", "USD"}, nil
}

func itemSah() models.ItemObjek {
	return models.ItemObjek{ItemTypeID: "UJI01", ItemType: "UJI MESIN", Currency: "IDR", TSI: ds("1500000000.12345678"), Unit: "2"}
}

// TestPeriksaItem - tiket 39: currency wajib (A133), tsi / pct desimal <= 8 (NUMBER(38,8)), unit bulat > 0
// (SetErrorMessageUnit_Act), Adjustable -> pctAdjustOther 60..100 (ValidateAdjustPct), lebar kolom; pesan
// menyebut indeks baris dan item.
func TestPeriksaItem(t *testing.T) {
	if m := periksaItem(0, []models.ItemObjek{itemSah(), {Currency: "USD"}}); len(m) != 0 {
		t.Fatalf("sah ditolak: %v", m)
	}
	ubah := func(f func(*models.ItemObjek)) []models.ItemObjek {
		it := itemSah()
		f(&it)
		return []models.ItemObjek{itemSah(), it}
	}
	for nama, u := range map[string]struct {
		item  []models.ItemObjek
		pesan string
	}{
		"tanpa mata uang": {ubah(func(i *models.ItemObjek) { i.Currency = "" }), "baris[3].items[1].currency wajib"},
		"mata uang spasi": {ubah(func(i *models.ItemObjek) { i.Currency = "  " }), "baris[3].items[1].currency wajib"},
		"tsi minus":       {ubah(func(i *models.ItemObjek) { i.TSI = ds("-1") }), "baris[3].items[1].tsi harus"},
		"tsi 9 desimal":   {ubah(func(i *models.ItemObjek) { i.TSI = ds("1.123456789") }), "baris[3].items[1].tsi"},
		"tsi 31 digit":    {ubah(func(i *models.ItemObjek) { i.TSI = ds(strings.Repeat("9", 31)) }), "baris[3].items[1].tsi"},
		"unit nol":        {ubah(func(i *models.ItemObjek) { i.Unit = "00" }), "baris[3].items[1].unit harus bilangan bulat > 0"},
		"unit desimal":    {ubah(func(i *models.ItemObjek) { i.Unit = "1.5" }), "baris[3].items[1].unit"},
		"pct2 minus":      {ubah(func(i *models.ItemObjek) { i.PctAdjust2 = ds("-7") }), "baris[3].items[1].pctAdjust2"},
		"adjust 59.99":    {ubah(func(i *models.ItemObjek) { i.IsAdjustable, i.PctAdjustOther = true, ds("59.99") }), "less than 60%"},
		"adjust 100.01":   {ubah(func(i *models.ItemObjek) { i.IsAdjustable, i.PctAdjustOther = true, ds("100.01") }), "more than 100%"},
		"adjust kosong":   {ubah(func(i *models.ItemObjek) { i.IsAdjustable, i.PctAdjustOther = true, nil }), "baris[3].items[1].pctAdjustOther"},
		"catatan 501":     {ubah(func(i *models.ItemObjek) { i.Note = strings.Repeat("U", 501) }), "baris[3].items[1].note paling banyak 500"},
		"kondisi 501":     {ubah(func(i *models.ItemObjek) { i.Condition = strings.Repeat("U", 501) }), "baris[3].items[1].condition paling banyak 500"},
	} {
		if m := strings.Join(periksaItem(3, u.item), "; "); !strings.Contains(m, u.pesan) {
			t.Errorf("%s: %q, mau %q", nama, m, u.pesan)
		}
	}
	for _, pct := range []string{"60", "100", "75.5", "100.00000000"} {
		it := itemSah()
		it.IsAdjustable, it.PctAdjustOther = true, ds(pct)
		if m := periksaItem(0, []models.ItemObjek{it}); len(m) != 0 {
			t.Errorf("adjust %s ditolak: %v", pct, m)
		}
	}
	// Tidak Adjustable: pctAdjustOther tidak dibatasi 60..100 (ResetPct_Adjustment mengisinya 100 di layar).
	it := itemSah()
	it.PctAdjustOther = ds("10")
	if m := periksaItem(0, []models.ItemObjek{it}); len(m) != 0 {
		t.Errorf("tidak adjustable: %v", m)
	}
}

// TestGantiObjekMataUangItem - tiket 39: mata uang item diperiksa ke CURRENCY (sekali per simpan, hanya bila ada
// item); tidak ada -> 400 ber-indeks; tanpa pembaca -> 503; item ikut tersimpan dan terbaca ulang.
func TestGantiObjekMataUangItem(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var panggil int
	svc := Baru(nil).DenganObjek(objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}).DenganTransaksi(tanpaTx).
		DenganPilihanItem(pilihanItemTiruan{&panggil})
	obj := []models.ObjekFire{{ObjectType: "UJI", Items: []models.ItemObjek{itemSah(), {Currency: "USD"}}}, {ObjectType: "UJI"}}
	d, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", obj)
	if err != nil || len(d) != 2 || len(d[0].Items) != 2 || utils.FormatDecimal(d[0].Items[0].TSI) != "1500000000.12345678" || panggil != 1 {
		t.Fatalf("%+v (%v) panggil=%d", d, err, panggil)
	}
	obj[1].Items = []models.ItemObjek{{Currency: "ITL"}}
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", obj); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), `baris[1].items[0].currency "ITL" tidak ada`) {
		t.Errorf("ITL: %v", err)
	}
	tanpaPilihan := Baru(nil).DenganObjek(objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}).DenganTransaksi(tanpaTx)
	if _, err := tanpaPilihan.GantiObjek(ctx, akun, "UJI-NB-1", obj); !errors.Is(err, ErrPilihanItemTanpaDatabase) {
		t.Errorf("tanpa pembaca mata uang: %v", err)
	}
	panggil = 0
	if _, err := tanpaPilihan.GantiObjek(ctx, akun, "UJI-NB-1", []models.ObjekFire{{ObjectType: "UJI"}}); err != nil {
		t.Errorf("tanpa item tidak butuh CURRENCY: %v", err)
	}
}

// TestPilihanItem - tiket 39: dua lookup tanpa identitas; 503 tanpa basis data.
func TestPilihanItem(t *testing.T) {
	ctx := context.Background()
	var panggil int
	svc := Baru(nil).DenganPilihanItem(pilihanItemTiruan{&panggil})
	if j, err := svc.DaftarJenisItem(ctx); err != nil || len(j) != 1 || j[0].Keterangan != "UJI KET" {
		t.Errorf("jenis item: %v %v", j, err)
	}
	if m, err := svc.DaftarMataUang(ctx); err != nil || len(m) != 2 {
		t.Errorf("mata uang: %v %v", m, err)
	}
	if _, err := Baru(nil).DaftarJenisItem(ctx); !errors.Is(err, ErrPilihanItemTanpaDatabase) {
		t.Errorf("jenis item tanpa DB: %v", err)
	}
	if _, err := Baru(nil).DaftarMataUang(ctx); !errors.Is(err, ErrPilihanItemTanpaDatabase) {
		t.Errorf("mata uang tanpa DB: %v", err)
	}
}

// TestLebarItemSamaDenganMigrasi - lebar validasi item = lebar kolom migrasi 188 (dipasangkan lewat NAMA kolom).
func TestLebarItemSamaDenganMigrasi(t *testing.T) {
	b, err := os.ReadFile("../migrations/188_t_propertyitemlist.sql")
	if err != nil {
		t.Fatal(err)
	}
	lebar := map[string]int{}
	for _, m := range regexp.MustCompile(`(?m)^\s+([A-Z_]+)\s+VARCHAR2\((\d+)\)`).FindAllStringSubmatch(string(b), -1) {
		lebar[m[1]], _ = strconv.Atoi(m[2])
	}
	kolom := map[string]string{"itemTypeId": "ITEM_TYPE_ID", "itemType": "ITEM_TYPE", "note": "PROPERTI_ITEM_NOTE",
		"propertyYear": "PROPERTY_YEAR", "unit": "UNIT", "condition": "CONDITION", "currency": "CURRENCY",
		"yearOfPlanting": "YEAR", "noOfTree": "NO_OF_TREE", "areaHectar": "AREA_HECTAR", "remark": "REMARK"}
	for _, l := range lebarItem {
		if n, ada := lebar[kolom[l.nama]]; !ada || n != l.n {
			t.Errorf("%s: lebar %d, kolom %s migrasi 188 = %d (ada=%v)", l.nama, l.n, kolom[l.nama], n, ada)
		}
	}
	if len(lebarItem) != len(kolom) {
		t.Errorf("%d medan diperiksa, peta %d", len(lebarItem), len(kolom))
	}
}

// ds - teks desimal uji -> *apd.Decimal ("" = nil; panik bila salah tulis di uji). Boleh di luar aturan polaDesimal
// (mis. "-1") untuk menguji pemeriksaan ulang services.
func ds(s string) *apd.Decimal {
	if s == "" {
		return nil
	}
	d, err := utils.ParseDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

// TestUraiDesimalIsian - ADR-0034 (butir 94): SATU fungsi teks JSON -> desimal; kosong = nil; "0.5", 8 desimal, 30 digit
// bulat sah; koma / huruf / minus / 9 desimal / 31 digit -> pesan ber-jalur (sama dengan sebelum refaktor).
func TestUraiDesimalIsian(t *testing.T) {
	for _, s := range []string{"0.5", "0.12345678", "123456789012345678901234567890", "123456789012345678901234567890.12345678", "100"} {
		var m []string
		if d := UraiDesimalIsian("baris[0].items[0].tsi", s, &m); d == nil || len(m) != 0 || utils.FormatDecimal(d) != s {
			t.Errorf("%q -> %v %v", s, d, m)
		}
	}
	var m []string
	if d := UraiDesimalIsian("x", "", &m); d != nil || len(m) != 0 {
		t.Errorf("kosong -> %v %v", d, m)
	}
	for _, s := range []string{"1,5", "tujuh", "-1", "1.123456789", strings.Repeat("9", 31), " 1", "1e3", ".5"} {
		var m []string
		if d := UraiDesimalIsian("baris[2].lossRecords[1].claim", s, &m); d != nil || len(m) != 1 ||
			m[0] != "baris[2].lossRecords[1].claim harus angka >= 0 dengan paling banyak 8 desimal" {
			t.Errorf("%q -> %v %v", s, d, m)
		}
	}
	if GalatIsianObjek(nil) != nil {
		t.Error("tanpa masalah harus nil")
	}
	if err := GalatIsianObjek([]string{"a", "b"}); !errors.Is(err, ErrMasukanObjek) || !strings.Contains(err.Error(), "a; b") {
		t.Errorf("galat: %v", err)
	}
}
