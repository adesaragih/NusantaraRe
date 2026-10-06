// Package tiruan - Gudang Treaty Group di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatygroup/backend/models"
	"nusantarare/modul/treatygroup/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung.
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// AnakTiruan - satu baris BUSINESSGROUP tiruan beserta TOPID-nya.
type AnakTiruan struct {
	TopID string
	models.BisnisGrup
}

// Gudang tiruan.
type Gudang struct {
	mu    sync.Mutex
	Baris map[string]models.Grup
	Ojk   []models.Ojk
	// BisnisGrup - BUSINESSGROUP tiruan (ID -> NOTE untuk nama COA; Anak untuk View).
	BisnisGrup []AnakTiruan
	// Seq - nomor sequence berikutnya; ID baru = "1" + 4 digit seperti situs aktif DEV.
	Seq int
}

// Contoh - tiga grup (dua se-OJK `01` ber-COAID `10013`), tiga OJK (OJK `03` belum punya grup), empat grup bisnis
// (satu berakhiran SYARIAH). Sequence berikutnya 9 - ID 10009 sudah terpakai (uji lompat).
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Grup{}, Seq: 9}
	g.Ojk = []models.Ojk{
		{ID: "01", Name: "UJI PROPERTY", NameIDN: "UJI HARTA BENDA", OrderNo: "1"},
		{ID: "09", Name: "UJI ENGINEERING", NameIDN: "UJI REKAYASA", OrderNo: "2"},
		{ID: "03", Name: "UJI CARGO", NameIDN: "UJI PENGANGKUTAN", OrderNo: "5"},
	}
	for _, r := range []models.Grup{
		{ID: "10007", OldID: "01", OjkID: "01", OjkName: "UJI PROPERTY", OjkNameIDN: "UJI HARTA BENDA", OrderNo: "1",
			Name: "UJI PROPERTY", SoaName: "UJI PROPERTY", TglUpdate: "20220105T135630.318 GMT", UserID: "UJI-PEGA", CoaID: "10013"},
		{ID: "10056", OjkID: "01", OjkName: "UJI PROPERTY", OjkNameIDN: "UJI HARTA BENDA", OrderNo: "1",
			Name: "UJI AGRICULTURE", CoaID: "10013"},
		{ID: "10009", OjkID: "09", OjkName: "UJI ENGINEERING LAMA", OjkNameIDN: "UJI REKAYASA", OrderNo: "2",
			Name: "UJI ENGINEERING", CoaID: "10019"},
	} {
		g.Baris[r.ID] = r
	}
	g.BisnisGrup = []AnakTiruan{
		{TopID: "10007", BisnisGrup: models.BisnisGrup{ID: "10013", Name: "UJI FIRE", Alias: "UJI FIRE"}},
		{TopID: "10007", BisnisGrup: models.BisnisGrup{ID: "10020", Name: "UJI FIRE SYARIAH", Alias: "UJI FIRE SYARIAH"}},
		{TopID: "10007", BisnisGrup: models.BisnisGrup{ID: "10014", Name: "UJI EARTHQUAKE", Alias: "UJI EQ"}},
		{TopID: "10009", BisnisGrup: models.BisnisGrup{ID: "10019", Name: "UJI CAR", Alias: "UJI CAR"}},
	}
	return g
}

func (g *Gudang) namaCoa(id string) string {
	for _, b := range g.BisnisGrup {
		if b.ID == id {
			return b.Name
		}
	}
	return ""
}

func (g *Gudang) dengan(r models.Grup) models.Grup {
	r.CoaName = g.namaCoa(r.CoaID)
	return r
}

// Daftar memenuhi services.Gudang - urut Order No (angka), nama, ID.
func (g *Gudang) Daftar(_ context.Context, kata, ojk string) ([]models.Grup, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := strings.ToUpper(strings.TrimSpace(kata))
	out := []models.Grup{}
	for _, r := range g.Baris {
		if ojk != "" && r.OjkID != ojk {
			continue
		}
		if k != "" && !strings.Contains(strings.ToUpper(r.ID+" "+r.Name+" "+r.SoaName+" "+r.OjkName), k) {
			continue
		}
		out = append(out, g.dengan(r))
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := fmt.Sprintf("%010s", out[i].OrderNo), fmt.Sprintf("%010s", out[j].OrderNo)
		if a != b {
			return a < b
		}
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Grup, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	r, ok := g.Baris[id]
	if !ok {
		return models.Grup{}, repository.ErrTidakAda
	}
	return g.dengan(r), nil
}

// Anak memenuhi services.Gudang - tanpa yang berakhiran SYARIAH, urut nama.
func (g *Gudang) Anak(_ context.Context, id string) ([]models.BisnisGrup, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []models.BisnisGrup{}
	for _, b := range g.BisnisGrup {
		if b.TopID == id && !strings.HasSuffix(strings.ToUpper(strings.TrimSpace(b.Name)), "SYARIAH") {
			out = append(out, b.BisnisGrup)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// DaftarOjk memenuhi services.Gudang - beserta COA setiap OJK (sama dengan CoaSeOjk).
func (g *Gudang) DaftarOjk(context.Context) ([]models.Ojk, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := append([]models.Ojk{}, g.Ojk...)
	for i := range out {
		out[i].CoaID = g.coaSeOjk(out[i].ID)
		out[i].CoaName = g.namaCoa(out[i].CoaID)
	}
	return out, nil
}

// AmbilOjk memenuhi services.Gudang.
func (g *Gudang) AmbilOjk(_ context.Context, _ *db.Tx, id string) (models.Ojk, error) {
	for _, o := range g.Ojk {
		if o.ID == id {
			return o, nil
		}
	}
	return models.Ojk{}, repository.ErrTidakAda
}

// PemakaiNama memenuhi services.Gudang.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []string
	for _, r := range g.Baris {
		if r.ID != kecualiID && strings.EqualFold(strings.TrimSpace(r.Name), strings.TrimSpace(nama)) {
			out = append(out, r.ID)
		}
	}
	return out, nil
}

// CoaSeOjk memenuhi services.Gudang - COAID paling sering di OJK itu, lalu terkecil.
func (g *Gudang) CoaSeOjk(_ context.Context, _ *db.Tx, ojk string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.coaSeOjk(ojk), nil
}

func (g *Gudang) coaSeOjk(ojk string) string {
	hitung := map[string]int{}
	for _, r := range g.Baris {
		if r.OjkID == ojk && r.CoaID != "" {
			hitung[r.CoaID]++
		}
	}
	terbaik := ""
	for c, n := range hitung {
		if terbaik == "" || n > hitung[terbaik] || (n == hitung[terbaik] && c < terbaik) {
			terbaik = c
		}
	}
	return terbaik
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
	_, ok := g.Baris[id]
	return ok, nil
}

// Sisip memenuhi services.Gudang.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, r models.Grup) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	r.CoaName, r.Diubah = "", ""
	g.Baris[r.ID] = r
	return nil
}

// Ubah memenuhi services.Gudang - ID dan OLDID tetap seperti SQL-nya.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, r models.Grup) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	lama, ok := g.Baris[r.ID]
	if !ok {
		return repository.ErrTidakAda
	}
	lama.OjkID, lama.OjkName, lama.OjkNameIDN, lama.OrderNo = r.OjkID, r.OjkName, r.OjkNameIDN, r.OrderNo
	lama.Name, lama.SoaName, lama.TglUpdate, lama.UserID, lama.CoaID = r.Name, r.SoaName, r.TglUpdate, r.UserID, r.CoaID
	g.Baris[r.ID] = lama
	return nil
}
