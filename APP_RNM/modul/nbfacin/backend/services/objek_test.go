package services

import (
	"context"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// objekTiruan - penyimpan berkeadaan: daftar per case, diganti utuh.
type objekTiruan struct{ ada map[string][]models.ObjekFire }

func (o objekTiruan) BacaObjek(_ context.Context, id string) ([]models.ObjekFire, error) {
	d, ok := o.ada[id]
	if !ok {
		return nil, repository.ErrKasusTidakAda
	}
	return append([]models.ObjekFire{}, d...), nil
}

func (o objekTiruan) GantiObjek(_ context.Context, _ *db.Tx, id string, baris []models.ObjekFire) error {
	if _, ok := o.ada[id]; !ok {
		return repository.ErrKasusTidakAda
	}
	o.ada[id] = append([]models.ObjekFire{}, baris...)
	return nil
}

// TestGantiObjek - tiket 35: urutan dipertahankan, ganti utuh, kosong mengosongkan, dibaca
// ulang; 401 / 400 (indeks baris) / 503 / 404.
func TestGantiObjek(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	svc := Baru(nil).DenganObjek(objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}).DenganTransaksi(tanpaTx)
	d, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", []models.ObjekFire{
		{ObjectNo: "2", ObjectType: "UJI TIPE B", IsTopRisk: true, NumberOfFloor: "0"},
		{ObjectNo: "1", ObjectType: "UJI TIPE A", NumberOfFloor: "12", RoofType: "14"},
	})
	if err != nil || len(d) != 2 || d[0].ObjectNo != "2" || !d[0].IsTopRisk || d[1].RoofType != "14" {
		t.Fatalf("%+v (%v)", d, err)
	}
	if d, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", nil); err != nil || len(d) != 0 {
		t.Errorf("kosong: %+v (%v)", d, err)
	}
	if _, err := svc.GantiObjek(ctx, inti.Pelaku{}, "UJI-NB-1", nil); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	for nama, u := range map[string]struct {
		baris []models.ObjekFire
		pesan string
	}{
		"tipe kosong":    {[]models.ObjekFire{{ObjectType: "UJI"}, {ObjectType: "  "}}, "baris[1].objectType wajib"},
		"lantai minus":   {[]models.ObjekFire{{ObjectType: "UJI", NumberOfFloor: "-1"}}, "baris[0].numberOfFloor"},
		"lantai desimal": {[]models.ObjekFire{{ObjectType: "UJI", NumberOfFloor: "1.5"}}, "baris[0].numberOfFloor"},
		"lantai teks":    {[]models.ObjekFire{{ObjectType: "UJI", NumberOfFloor: "dua"}}, "baris[0].numberOfFloor"},
		"nama 501 bita":  {[]models.ObjekFire{{ObjectType: "UJI", ObjectName: strings.Repeat("U", 501)}}, "baris[0].objectName paling banyak 500"},
		"lokasi 51":      {[]models.ObjekFire{{ObjectType: "UJI", RiskLocation: strings.Repeat("U", 51)}}, "baris[0].riskLocation paling banyak 50"},
		"ownership 51":   {[]models.ObjekFire{{ObjectType: "UJI", Ownership: strings.Repeat("2", 51)}}, "baris[0].ownership paling banyak 50"},
		"konstruksi 501": {[]models.ObjekFire{{ObjectType: "UJI", SurroundingRisk: models.SurroundingRisk{Back: models.SisiRisiko{Construction: strings.Repeat("U", 501)}}}}, "baris[0].surroundingRisk.back.construction paling banyak 500"},
		"remark 501":     {[]models.ObjekFire{{ObjectType: "UJI", SurroundingRisk: models.SurroundingRisk{HousekeepingRemark: strings.Repeat("U", 501)}}}, "baris[0].surroundingRisk.housekeepingRemark paling banyak 500"},
		"jarak minus":    {[]models.ObjekFire{{ObjectType: "UJI"}, {ObjectType: "UJI", SurroundingRisk: models.SurroundingRisk{Left: models.SisiRisiko{Distance: "-1"}}}}, "baris[1].surroundingRisk.left.distance harus angka 0..100000"},
	} {
		if _, err := Baru(nil).GantiObjek(ctx, akun, "UJI-NB-1", u.baris); !errors.Is(err, ErrMasukanObjek) || !strings.Contains(err.Error(), u.pesan) {
			t.Errorf("%s: %v, mau 400 %q (sebelum DB)", nama, err, u.pesan)
		}
	}
	if _, err := Baru(nil).GantiObjek(ctx, akun, "UJI-NB-1", nil); !errors.Is(err, ErrObjekTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-TIDAK-ADA", nil); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("case tidak ada: %v", err)
	}
	if _, err := svc.BacaObjek(ctx, "UJI-NB-TIDAK-ADA"); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("baca case tidak ada: %v", err)
	}
}

// TestJarakSah - tiket 38 (A129/A131): kosong, atau angka >= 0, <= 2 desimal, <= 100000.
func TestJarakSah(t *testing.T) {
	for s, sah := range map[string]bool{"": true, "0": true, "12": true, "12.5": true, "12.50": true, "100000": true,
		"100000.00": true, "007": true, "100000.01": false, "100001": false, "1.234": false, "-1": false, "1,5": false,
		".5": false, "5.": false, " 5": false, "1e3": false, "dua": false} {
		if jarakSah(s) != sah {
			t.Errorf("%q: sah=%v, mau %v", s, !sah, sah)
		}
	}
	// lebar 50 tetap diperiksa: angka sah dengan nol di depan bisa melebihi kolom.
	baris := []models.ObjekFire{{ObjectType: "UJI", SurroundingRisk: models.SurroundingRisk{Right: models.SisiRisiko{Distance: strings.Repeat("0", 51)}}}}
	if err := periksaObjek(baris); !errors.Is(err, ErrMasukanObjek) || !strings.Contains(err.Error(), "baris[0].surroundingRisk.right.distance paling banyak 50") {
		t.Errorf("jarak 51 bita: %v", err)
	}
	for _, s := range sisiObjek {
		var o models.ObjekFire
		o.ObjectType = "UJI"
		o.SurroundingRisk = models.SurroundingRisk{Front: models.SisiRisiko{Distance: "x"}, Left: models.SisiRisiko{Distance: "x"},
			Back: models.SisiRisiko{Distance: "x"}, Right: models.SisiRisiko{Distance: "x"}}
		if err := periksaObjek([]models.ObjekFire{o}); err == nil || !strings.Contains(err.Error(), "surroundingRisk."+s.nama+".distance harus") {
			t.Errorf("sisi %s tidak diperiksa: %v", s.nama, err)
		}
	}
}

// TestLebarSekitarSamaDenganMigrasi - tiket 38: lebar validasi surroundingRisk / ownership = lebar
// kolom migrasi 187 (dipasangkan lewat NAMA kolom, bukan urutan).
func TestLebarSekitarSamaDenganMigrasi(t *testing.T) {
	b, err := os.ReadFile("../migrations/187_t_surroundingrisk.sql")
	if err != nil {
		t.Fatal(err)
	}
	lebar := map[string]int{}
	for _, m := range regexp.MustCompile(`(?m)^\s+([A-Z_]+)\s+VARCHAR2\((\d+)\)`).FindAllStringSubmatch(string(b), -1) {
		lebar[m[1]], _ = strconv.Atoi(m[2])
	}
	ular := regexp.MustCompile(`([a-z])([A-Z])`)
	diperiksa := 0
	for _, l := range lebarObjek {
		nama := strings.TrimPrefix(l.nama, "surroundingRisk.")
		if nama == l.nama && l.nama != "ownership" {
			continue // medan tiket 35 (migrasi 186)
		}
		kolom := strings.ToUpper(ular.ReplaceAllString(strings.ReplaceAll(nama, ".", "_"), "${1}_${2}"))
		if n, ada := lebar[kolom]; !ada || n != l.n {
			t.Errorf("%s: lebar %d, kolom %s migrasi 187 = %d (ada=%v)", l.nama, l.n, kolom, n, ada)
		}
		diperiksa++
	}
	if diperiksa != 21 {
		t.Errorf("%d medan diperiksa, mau 21 (ownership + 16 sisi + 4)", diperiksa)
	}
}
