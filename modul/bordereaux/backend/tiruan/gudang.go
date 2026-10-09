// Package tiruan - Gudang Bordereaux di memori untuk uji layanan dan handler. Data UJI-* buatan.
package tiruan

import (
	"context"
	"sort"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung; galat tidak dibatalkan (cukup untuk uji aturan).
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// Gudang tiruan.
type Gudang struct {
	mu         sync.Mutex
	Header     map[string]models.Header
	Rinci      map[string]map[string][]models.Baris // tabel -> BDX_ID -> baris (ber-ID)
	Riw        map[string][]models.Riwayat
	Treaty     []models.MasterTreaty
	Cedant     []models.Cedant
	IDTerpakai map[string]bool
	// JSON - isi M_BORDEREAUX.DATA_JSON per ID (Copy Old Data).
	JSON map[string]string
	// Jam - TANGGAL berkas baru.
	Jam string
	// Kategori - master M_KATEGORIBORDEREAUX (ID, Nama); Lampiran - M_ATTACHMENTBORDEREAUX per ID.
	Kategori []models.KategoriLampiran
	Lampiran map[string]models.Lampiran
}

// Contoh - satu cedant, dua treaty, satu berkas lama di tangan Checker.
func Contoh() *Gudang {
	g := &Gudang{
		Header: map[string]models.Header{}, Rinci: map[string]map[string][]models.Baris{}, Riw: map[string][]models.Riwayat{},
		IDTerpakai: map[string]bool{}, JSON: map[string]string{}, Jam: "04-10-2026 09:30", Lampiran: map[string]models.Lampiran{},
		Kategori: []models.KategoriLampiran{{ID: "00000", Nama: "Others"}, {ID: "00001", Nama: "UJI SOA"}},
		Cedant:   []models.Cedant{{ID: "UJI-AG-1", Nama: "UJI CEDANT SATU"}, {ID: "UJI-AG-2", Nama: "UJI CEDANT DUA"}},
		Treaty: []models.MasterTreaty{
			{ID: "UJI-TI-1", ContractName: "UJI KONTRAK SATU", ReinsType: "Proportional", SobID: "UJI-AG-9", SobName: "UJI SOB", CedingID: "UJI-AG-1", CedingName: "UJI CEDANT SATU"},
			{ID: "UJI-TI-2", ContractName: "UJI KONTRAK DUA", ReinsType: "NonProportional", SobID: "UJI-AG-9", SobName: "UJI SOB", CedingID: "UJI-AG-1", CedingName: "UJI CEDANT SATU"},
		},
	}
	g.Header["BDX-UJI.1"] = models.Header{BdxID: "BDX-UJI.1", Tanggal: "01-09-2026 10:00", UserInput: "UJI-MAKER", Type: "PREMIUM",
		TypeBusiness: "FIRE", MasterID: "UJI-TI-1", CedingName: "UJI CEDANT SATU", ReportStart: "01-07-2026", ReportEnd: "30-09-2026",
		Position: models.PosisiChecker, StatusAksep: models.StatusAccept}
	return g
}

// Daftar memenuhi services.Gudang (saringan tiruan: BdxID, Type, Ceding).
func (g *Gudang) Daftar(_ context.Context, f models.Filter, offset, ukuran int) ([]models.Header, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var semua []models.Header
	for _, h := range g.Header {
		if f.BdxID != "" && !strings.Contains(strings.ToUpper(h.BdxID), strings.ToUpper(f.BdxID)) {
			continue
		}
		if f.Type != "" && !strings.EqualFold(h.Type, f.Type) {
			continue
		}
		if f.Ceding != "" && !strings.Contains(strings.ToUpper(h.CedingName), strings.ToUpper(f.Ceding)) {
			continue
		}
		semua = append(semua, h)
	}
	sort.Slice(semua, func(i, j int) bool { return semua[i].BdxID > semua[j].BdxID })
	total := len(semua)
	if offset > total {
		offset = total
	}
	akhir := offset + ukuran
	if akhir > total {
		akhir = total
	}
	return semua[offset:akhir], total, nil
}

// AmbilHeader memenuhi services.Gudang.
func (g *Gudang) AmbilHeader(_ context.Context, _ *db.Tx, id string) (models.Header, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	h, ok := g.Header[id]
	if !ok {
		return models.Header{}, repository.ErrTidakAda
	}
	return h, nil
}

// AdaBerkas memenuhi services.Gudang.
func (g *Gudang) AdaBerkas(_ context.Context, _ *db.Tx, id string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, ok := g.Header[id]
	return ok, nil
}

// SisipHeader memenuhi services.Gudang.
func (g *Gudang) SisipHeader(_ context.Context, _ *db.Tx, h models.Header) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if h.Tanggal == "" {
		h.Tanggal = g.Jam
	}
	g.Header[h.BdxID] = h
	return nil
}

// UbahHeader memenuhi services.Gudang.
func (g *Gudang) UbahHeader(_ context.Context, _ *db.Tx, h models.Header) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	lama, ok := g.Header[h.BdxID]
	if !ok {
		return repository.ErrTidakAda
	}
	h.Tanggal, h.UserInput, h.Position, h.StatusAksep = lama.Tanggal, lama.UserInput, lama.Position, lama.StatusAksep
	g.Header[h.BdxID] = h
	return nil
}

// UbahStatus memenuhi services.Gudang.
func (g *Gudang) UbahStatus(_ context.Context, _ *db.Tx, id, posisiLama, statusLama, posisiBaru, statusBaru string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	h, ok := g.Header[id]
	if !ok || h.Position != posisiLama || h.StatusAksep != statusLama {
		return false, nil
	}
	h.Position, h.StatusAksep = posisiBaru, statusBaru
	g.Header[id] = h
	return true, nil
}

// Detail memenuhi services.Gudang.
func (g *Gudang) Detail(_ context.Context, _ *db.Tx, k models.KombinasiBdx, id string) ([]models.Baris, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]models.Baris{}, g.Rinci[k.Tabel][id]...), nil
}

// HapusDetail memenuhi services.Gudang.
func (g *Gudang) HapusDetail(_ context.Context, _ *db.Tx, k models.KombinasiBdx, id string) (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := len(g.Rinci[k.Tabel][id])
	for _, b := range g.Rinci[k.Tabel][id] {
		delete(g.IDTerpakai, b[models.KolomID])
	}
	delete(g.Rinci[k.Tabel], id)
	return int64(n), nil
}

// SisipDetail memenuhi services.Gudang; ID terpakai = repository.ErrIDTerpakai.
func (g *Gudang) SisipDetail(_ context.Context, _ *db.Tx, k models.KombinasiBdx, id, idDetail string, b models.Baris) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.IDTerpakai[idDetail] {
		return repository.ErrIDTerpakai
	}
	g.IDTerpakai[idDetail] = true
	salin := models.Baris{models.KolomID: idDetail}
	for kk, v := range b {
		if kk != models.KolomID {
			salin[kk] = v
		}
	}
	if g.Rinci[k.Tabel] == nil {
		g.Rinci[k.Tabel] = map[string][]models.Baris{}
	}
	g.Rinci[k.Tabel][id] = append(g.Rinci[k.Tabel][id], salin)
	return nil
}

// HapusBerkas memenuhi services.Gudang.
func (g *Gudang) HapusBerkas(_ context.Context, _ *db.Tx, id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.Header[id]; !ok {
		return repository.ErrTidakAda
	}
	for t := range g.Rinci {
		delete(g.Rinci[t], id)
	}
	delete(g.Riw, id)
	delete(g.JSON, id)
	delete(g.Header, id)
	return nil
}

// Riwayat memenuhi services.Gudang.
func (g *Gudang) Riwayat(_ context.Context, id string) ([]models.Riwayat, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]models.Riwayat{}, g.Riw[id]...), nil
}

// SisipRiwayat memenuhi services.Gudang.
func (g *Gudang) SisipRiwayat(_ context.Context, _ *db.Tx, id, pic string, setuju bool, komentar string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Riw[id] = append(g.Riw[id], models.Riwayat{Tanggal: g.Jam, PIC: pic, IsApproved: setuju, Komentar: komentar})
	return nil
}

// CariCedant memenuhi services.Gudang.
func (g *Gudang) CariCedant(_ context.Context, kata string, batas int) ([]models.Cedant, error) {
	var out []models.Cedant
	for _, c := range g.Cedant {
		if strings.Contains(strings.ToUpper(c.Nama), strings.ToUpper(strings.TrimSpace(kata))) && len(out) < batas {
			out = append(out, c)
		}
	}
	return out, nil
}

// CariTreaty memenuhi services.Gudang.
func (g *Gudang) CariTreaty(_ context.Context, cedingID string, batas int) ([]models.MasterTreaty, error) {
	var out []models.MasterTreaty
	for _, t := range g.Treaty {
		if t.CedingID == cedingID && len(out) < batas {
			out = append(out, t)
		}
	}
	return out, nil
}

// AmbilTreaty memenuhi services.Gudang.
func (g *Gudang) AmbilTreaty(_ context.Context, id string) (models.MasterTreaty, error) {
	for _, t := range g.Treaty {
		if t.ID == id {
			return t, nil
		}
	}
	return models.MasterTreaty{}, repository.ErrTidakAda
}

// Chart memenuhi services.Gudang.
func (g *Gudang) Chart(context.Context) ([]models.IrisanChart, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	per := map[models.IrisanChart]int{}
	for _, h := range g.Header {
		per[models.IrisanChart{Business: h.TypeBusiness, Type: h.Type, CedingID: h.CedingID, CedingName: h.CedingName}]++
	}
	out := []models.IrisanChart{}
	for k, n := range per {
		k.Jumlah = n
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Business+out[i].Type+out[i].CedingName < out[j].Business+out[j].Type+out[j].CedingName
	})
	return out, nil
}

// JSONLama memenuhi services.Gudang (urut ID turun).
func (g *Gudang) JSONLama(_ context.Context, _ *db.Tx, id string) ([]models.JSONLama, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []models.JSONLama{}
	for k, isi := range g.JSON {
		if id != "" && k != id {
			continue
		}
		j := models.JSONLama{ID: k, Isi: isi}
		if h, ok := g.Header[k]; ok {
			j.AdaHeader, j.Type, j.Business = true, h.Type, h.TypeBusiness
		}
		out = append(out, j)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out, nil
}

// CacahDetail memenuhi services.Gudang.
func (g *Gudang) CacahDetail(_ context.Context, k models.KombinasiBdx) (map[string]int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := map[string]int{}
	for id, b := range g.Rinci[k.Tabel] {
		if len(b) > 0 {
			out[id] = len(b)
		}
	}
	return out, nil
}

// CacahRiwayat memenuhi services.Gudang.
func (g *Gudang) CacahRiwayat(context.Context) (map[string]int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := map[string]int{}
	for id, r := range g.Riw {
		if len(r) > 0 {
			out[id] = len(r)
		}
	}
	return out, nil
}

// SisipRiwayatLama memenuhi services.Gudang; Tanggal disimpan tanpa detik seperti bacaan Oracle.
func (g *Gudang) SisipRiwayatLama(_ context.Context, _ *db.Tx, id, tanggal, pic string, setuju bool, komentar string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(tanggal) >= 16 {
		tanggal = tanggal[:16]
	}
	g.Riw[id] = append(g.Riw[id], models.Riwayat{Tanggal: tanggal, PIC: pic, IsApproved: setuju, Komentar: komentar})
	return nil
}

// KategoriLampiran memenuhi services.Gudang - setiap kategori master dan jumlah lampiran berkas itu.
func (g *Gudang) KategoriLampiran(_ context.Context, bdxID string) ([]models.KategoriLampiran, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]models.KategoriLampiran, 0, len(g.Kategori))
	for _, k := range g.Kategori {
		k.Cacah = 0
		for _, l := range g.Lampiran {
			if l.BdxID == bdxID && l.KategoriID == k.ID {
				k.Cacah++
			}
		}
		out = append(out, k)
	}
	return out, nil
}

// NamaKategoriLampiran memenuhi services.Gudang.
func (g *Gudang) NamaKategoriLampiran(_ context.Context, id string) (string, bool, error) {
	for _, k := range g.Kategori {
		if k.ID == id {
			return k.Nama, true, nil
		}
	}
	return "", false, nil
}

// DaftarLampiran memenuhi services.Gudang (urut ID).
func (g *Gudang) DaftarLampiran(_ context.Context, bdxID, kategoriID string) ([]models.Lampiran, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []models.Lampiran{}
	for _, l := range g.Lampiran {
		if l.BdxID == bdxID && l.KategoriID == kategoriID {
			out = append(out, l)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// AmbilLampiran memenuhi services.Gudang.
func (g *Gudang) AmbilLampiran(_ context.Context, bdxID, kategoriID, id string) (models.Lampiran, bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	l, ada := g.Lampiran[id]
	if !ada || l.BdxID != bdxID || l.KategoriID != kategoriID {
		return models.Lampiran{}, false, nil
	}
	return l, true, nil
}

// SisipLampiran memenuhi services.Gudang; ID terpakai (atau di IDTerpakai) = ErrIDTerpakai.
func (g *Gudang) SisipLampiran(_ context.Context, _ *db.Tx, a models.Lampiran) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Lampiran[a.ID]; ada || g.IDTerpakai[a.ID] {
		return repository.ErrIDTerpakai
	}
	g.Lampiran[a.ID] = a
	return nil
}

// HapusLampiran memenuhi services.Gudang.
func (g *Gudang) HapusLampiran(_ context.Context, _ *db.Tx, bdxID, id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if l, ada := g.Lampiran[id]; ada && l.BdxID == bdxID {
		delete(g.Lampiran, id)
	}
	return nil
}
