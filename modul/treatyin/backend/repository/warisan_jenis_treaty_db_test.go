//go:build db

package repository_test

import "testing"

// ⭐ Kode yang tersimpan di kontrak terbaca sebagai teksnya — dan pilihannya
// berurut `ID` menurun, seperti RD Pega.
func TestDaftarJenisTreatyDariReinsuranceType(t *testing.T) {
	g, ctx := gudangBaca(t)
	d, err := g.BacaDaftarJenisTreaty(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(d) == 0 || len(d) > 500 {
		t.Fatalf("%d pilihan", len(d))
	}
	nama := map[string]string{}
	for i, p := range d {
		nama[p.ID] = p.Nama
		if i > 0 && d[i-1].ID < p.ID {
			t.Errorf("urutan bukan ID menurun: %s lalu %s", d[i-1].ID, p.ID)
		}
	}
	// ⭐ SOA Name ikut terbaca — `2020 QS 89M FAC` (10210) ber-SOA `QUOTA SHARE 2020`.
	adaSOA := false
	for _, p := range d {
		if p.ID == "10210" {
			adaSOA = p.NamaSOA == "QUOTA SHARE 2020"
		}
	}
	if !adaSOA {
		t.Error("10210 tidak membawa SOA Name `QUOTA SHARE 2020`")
	}
	for id, mau := range map[string]string{"10042": "SURPLUS", "10035": "QUOTA SHARE", "10037": "2ND SURPLUS"} {
		if nama[id] != mau {
			t.Errorf("%s -> %q, mau %q", id, nama[id], mau)
		}
	}
}

// ⭐ Treaty Group dan mata uang — RD `BrowseTreatyGroup_RD` dan
// `BrowseCurrencyTreatyIn_RD`.
func TestDaftarKelompokDanMataUangLimit(t *testing.T) {
	g, ctx := gudangBaca(t)
	kel, err := g.BacaDaftarKelompokTreaty(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(kel) < 30 {
		t.Errorf("%d kelompok, terukur 33", len(kel))
	}
	mu, err := g.BacaDaftarMataUangLimit(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ada := map[string]string{}
	for _, m := range mu {
		ada[m.Nama] = m.ID
	}
	if ada["USD"] != "10001" || ada["IDR"] != "10026" {
		t.Errorf("USD %q IDR %q, mau 10001/10026", ada["USD"], ada["IDR"])
	}
	if _, itl := ada["ITL"]; itl {
		t.Error("ITL ikut — RD menyaringnya")
	}
}
