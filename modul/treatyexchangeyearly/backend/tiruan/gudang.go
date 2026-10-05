// Package tiruan - Gudang Treaty Exchange Yearly di memori untuk uji layanan dan handler. Data UJI buatan; Kunci tiruan
// menggantikan ROWID.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyexchangeyearly/backend/models"
	"nusantarare/modul/treatyexchangeyearly/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung.
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// Gudang tiruan - Baris per Kunci.
type Gudang struct {
	mu       sync.Mutex
	Baris    map[string]models.Kurs
	MataUang []models.MataUang
	// Seq - nomor sequence berikutnya; ID baru = "1" + 4 digit seperti situs aktif DEV.
	Seq   int
	kunci int
}

// Contoh - lima baris: dua berbagi ID 10114 (ID warisan tidak unik), satu bertanggal warisan rusak (`T00000`), satu
// ber-jam (2025). Sequence berikutnya 114 - ID 10114 terpakai (uji lompat).
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Kurs{}, Seq: 114}
	g.MataUang = []models.MataUang{{ID: "10001", Kode: "USD", Nama: "UJI DOLLAR"}, {ID: "10002", Kode: "SGD", Nama: "UJI SGD"},
		{ID: "10025", Kode: "EUR", Nama: "UJI EURO"}}
	for _, k := range []models.Kurs{
		{Kunci: "AAA1", ID: "10010", TreatyYear: "2019", IDCurrency: "10001", Currency: "USD", StartDate: "20190801T00000.000 GMT",
			EndDate: "20200630T000000.000 GMT", ToIDR: "14500.00", ToUSD: "0", Quarter: "0", UserID: "UJI-PEGA",
			DateIU: "20230713T103900.000 GMT", DateIn: "20230713T103923.330 GMT"},
		{Kunci: "AAA2", ID: "10114", TreatyYear: "2025", IDCurrency: "10001", Currency: "USD", StartDate: "20250701T075400.000 GMT",
			EndDate: "20260630T152800.000 GMT", ToIDR: "16300", ToUSD: "1", Quarter: "0"},
		{Kunci: "AAA3", ID: "10114", TreatyYear: "2025", IDCurrency: "10025", Currency: "EUR", StartDate: "20250701T000000.000 GMT",
			EndDate: "20260630T000000.000 GMT", ToIDR: "18900", Quarter: "0"},
		{Kunci: "AAA4", ID: "10050", TreatyYear: "2024", IDCurrency: "10002", Currency: "SGD", StartDate: "20240701T000000.000 GMT",
			EndDate: "20250630T000000.000 GMT", ToIDR: "12000", ToUSD: "16000", Quarter: "0"},
		{Kunci: "AAA5", ID: "10051", TreatyYear: "2024", IDCurrency: "10002", Currency: "SGD", StartDate: "20240701T000000.000 GMT",
			EndDate: "20240930T000000.000 GMT", ToIDR: "11900", Quarter: "1"},
	} {
		g.Baris[k.Kunci] = k
	}
	return g
}

func (g *Gudang) dengan(k models.Kurs) models.Kurs {
	for _, m := range g.MataUang {
		if m.ID == k.IDCurrency {
			k.CurrencyName = m.Nama
		}
	}
	return k
}

// Daftar memenuhi services.Gudang - urut tahun terbaru, kode, quarter, ID.
func (g *Gudang) Daftar(_ context.Context, kata, tahun string) ([]models.Kurs, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	kt := strings.ToUpper(strings.TrimSpace(kata))
	out := []models.Kurs{}
	for _, k := range g.Baris {
		k = g.dengan(k)
		if tahun != "" && k.TreatyYear != tahun {
			continue
		}
		if kt != "" && !strings.Contains(strings.ToUpper(k.ID+" "+k.Currency+" "+k.CurrencyName), kt) {
			continue
		}
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.TreatyYear != b.TreatyYear {
			return a.TreatyYear > b.TreatyYear
		}
		if a.Currency != b.Currency {
			return a.Currency < b.Currency
		}
		if a.Quarter != b.Quarter {
			return a.Quarter < b.Quarter
		}
		return a.ID < b.ID
	})
	return out, nil
}

// AmbilKunci memenuhi services.Gudang.
func (g *Gudang) AmbilKunci(_ context.Context, _ *db.Tx, kunci string) (models.Kurs, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ok := g.Baris[kunci]
	if !ok {
		return models.Kurs{}, repository.ErrTidakAda
	}
	return g.dengan(k), nil
}

// AmbilID memenuhi services.Gudang - baris pertama menurut Kunci.
func (g *Gudang) AmbilID(_ context.Context, _ *db.Tx, id string) (models.Kurs, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var kunci []string
	for k, b := range g.Baris {
		if b.ID == id {
			kunci = append(kunci, k)
		}
	}
	if len(kunci) == 0 {
		return models.Kurs{}, repository.ErrTidakAda
	}
	sort.Strings(kunci)
	return g.dengan(g.Baris[kunci[0]]), nil
}

// Kembar memenuhi services.Gudang.
func (g *Gudang) Kembar(_ context.Context, _ *db.Tx, tahun, idCurrency, quarter, kecualiKunci string) ([]string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []string
	for k, b := range g.Baris {
		if k != kecualiKunci && b.TreatyYear == tahun && b.IDCurrency == idCurrency && b.Quarter == quarter {
			out = append(out, b.ID)
		}
	}
	sort.Strings(out)
	return out, nil
}

// DaftarMataUang memenuhi services.Gudang.
func (g *Gudang) DaftarMataUang(context.Context) ([]models.MataUang, error) {
	return append([]models.MataUang{}, g.MataUang...), nil
}

// AmbilMataUang memenuhi services.Gudang.
func (g *Gudang) AmbilMataUang(_ context.Context, _ *db.Tx, id string) (models.MataUang, error) {
	for _, m := range g.MataUang {
		if m.ID == id {
			return m, nil
		}
	}
	return models.MataUang{}, repository.ErrTidakAda
}

// DaftarTahun memenuhi services.Gudang.
func (g *Gudang) DaftarTahun(context.Context) ([]string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	ada := map[string]bool{}
	var out []string
	for _, k := range g.Baris {
		if !ada[k.TreatyYear] {
			ada[k.TreatyYear] = true
			out = append(out, k.TreatyYear)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(out)))
	return out, nil
}

// IDBaru memenuhi services.Gudang.
func (g *Gudang) IDBaru(context.Context, *db.Tx) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	id := fmt.Sprintf("1%04d", g.Seq)
	g.Seq++
	return id, nil
}

// AdaID memenuhi services.Gudang.
func (g *Gudang) AdaID(_ context.Context, _ *db.Tx, id string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, b := range g.Baris {
		if b.ID == id {
			return true, nil
		}
	}
	return false, nil
}

// Sisip memenuhi services.Gudang - Kunci tiruan baru.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, k models.Kurs) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.kunci++
	k.Kunci, k.CurrencyName = fmt.Sprintf("BARU%d", g.kunci), ""
	g.Baris[k.Kunci] = k
	return nil
}

// Ubah memenuhi services.Gudang - ID dan DATEIN tetap seperti SQL-nya.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, k models.Kurs) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	lama, ok := g.Baris[k.Kunci]
	if !ok {
		return repository.ErrTidakAda
	}
	k.ID, k.DateIn, k.CurrencyName = lama.ID, lama.DateIn, ""
	k.Mulai, k.Akhir, k.Diubah = "", "", ""
	g.Baris[k.Kunci] = k
	return nil
}
