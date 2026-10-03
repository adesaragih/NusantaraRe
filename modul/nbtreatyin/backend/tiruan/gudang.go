// Package tiruan adalah gudang NB Treaty In di memori - untuk uji seam 1
// (HTTP handler, spec §6.2) tanpa Oracle. Ia memenuhi
// `services.Gudang` dengan perilaku yang SAMA PENTINGNYA dengan Oracle untuk
// uji: transaksi yang benar-benar batal (snapshot dipulihkan bila fn gagal,
// AC 29, 83), tahap lama di syarat pindah tahap, nomor polis sekali.
//
// ⛔ Data uji berawalan UJI- (pagar keamanan brief): nol data polis atau nama
// orang sungguhan.
package tiruan

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// Gudang adalah penyimpanan di memori.
type Gudang struct {
	mu sync.Mutex

	urut    int
	Kasus   map[string]models.Kasus
	Halaman map[string]*models.Halaman
	Riwayat []models.Riwayat
	Nama    map[string]string // login -> nama tampilan
	Kontrak map[string]models.BarisKontrak
	Bisnis  map[string]models.BarisBisnis
	Serupa  []string
	Agen    []models.BarisAgen // RD BrowseAgentHierarkiList_RD (pemilih SOB)
	PKPAgen map[string]string  // STS_PKP per ID agen; tak terdaftar = StsPKP
	Tempat  []repository.PeranTempat
	StsPKP  string
	OJK     string
	urutPol int
	Closing int
	GagalDi string // nama operasi yang dipaksa gagal (uji pembatalan transaksi)
	dalamTx bool
	Panggil []string
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{
		Kasus:   map[string]models.Kasus{},
		Halaman: map[string]*models.Halaman{},
		Nama:    map[string]string{},
		Kontrak: map[string]models.BarisKontrak{},
		Bisnis:  map[string]models.BarisBisnis{},
		Closing: 25,
		OJK:     "UJI-OJK",
	}
}

// ErrDipaksa - galat yang disuntikkan uji.
var ErrDipaksa = errors.New("tiruan: galat dipaksa")

func (g *Gudang) gagal(op string) error {
	g.Panggil = append(g.Panggil, op)
	if g.GagalDi == op {
		return ErrDipaksa
	}
	return nil
}

type potret struct {
	urut, urutPol int
	kasus         map[string]models.Kasus
	halaman       map[string]*models.Halaman
	riwayat       []models.Riwayat
}

func (g *Gudang) potret() potret {
	p := potret{urut: g.urut, urutPol: g.urutPol, kasus: map[string]models.Kasus{}, halaman: map[string]*models.Halaman{}}
	for k, v := range g.Kasus {
		p.kasus[k] = v
	}
	for k, v := range g.Halaman {
		p.halaman[k] = v.Salin()
	}
	p.riwayat = append(p.riwayat, g.Riwayat...)
	return p
}

// Transaksi menjalankan fn; galat memulihkan seluruh keadaan.
func (g *Gudang) Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error {
	g.mu.Lock()
	sebelum := g.potret()
	g.dalamTx = true
	g.mu.Unlock()
	err := fn(&db.Tx{})
	g.mu.Lock()
	defer g.mu.Unlock()
	g.dalamTx = false
	if err != nil {
		g.urut, g.urutPol = sebelum.urut, sebelum.urutPol
		g.Kasus, g.Halaman, g.Riwayat = sebelum.kasus, sebelum.halaman, sebelum.riwayat
	}
	return err
}

func (g *Gudang) IDKasusBerikut(context.Context, *db.Tx) (string, error) {
	g.urut++
	return models.RakitIDKasus(fmt.Sprint(g.urut)), g.gagal("IDKasusBerikut")
}

func (g *Gudang) SisipKasus(_ context.Context, _ *db.Tx, id, pembuat, _ string) error {
	if err := g.gagal("SisipKasus"); err != nil {
		return err
	}
	g.Kasus[id] = models.Kasus{ID: id, Position: models.PositionAdmin, StatusWork: models.AssignmentAdmin,
		PositionNote: models.PosisiAdmin, CreateOp: pembuat, TglCreate: "2026-10-03 09:00:00"}
	g.Halaman[id] = models.HalamanBaru()
	return nil
}

func (g *Gudang) Keadaan(_ context.Context, _ *db.Tx, id string) (models.Kasus, error) {
	k, ada := g.Kasus[id]
	if !ada {
		return models.Kasus{}, repository.ErrKasusTidakAda
	}
	return k, nil
}

// KunciKasus - tahap kasus masih tahap yang dibaca.
func (g *Gudang) KunciKasus(_ context.Context, _ *db.Tx, id, statusHarap string) error {
	k, ada := g.Kasus[id]
	if !ada {
		return repository.ErrKasusTidakAda
	}
	if k.StatusWork != statusHarap {
		return repository.ErrTahapBerubah
	}
	return nil
}

func (g *Gudang) PindahPosisi(_ context.Context, _ *db.Tx, id, statusLama, posisiBaru string) error {
	if err := g.gagal("PindahPosisi"); err != nil {
		return err
	}
	k := g.Kasus[id]
	if k.StatusWork != statusLama {
		return repository.ErrTahapBerubah
	}
	k.Position, k.StatusWork, k.PositionNote = models.PositionPosisi(posisiBaru), models.AssignmentPosisi(posisiBaru), posisiBaru
	g.Kasus[id] = k
	return nil
}

func (g *Gudang) TutupKasus(_ context.Context, _ *db.Tx, id, statusLama, statusAkhir string) error {
	if err := g.gagal("TutupKasus"); err != nil {
		return err
	}
	k := g.Kasus[id]
	if k.StatusWork != statusLama {
		return repository.ErrTahapBerubah
	}
	k.StatusWork, k.Position = statusAkhir, ""
	g.Kasus[id] = k
	return nil
}

func (g *Gudang) DaftarKasus(_ context.Context, s models.SaringanKasus) ([]models.RingkasanKasus, error) {
	var out []models.RingkasanKasus
	for id, k := range g.Kasus {
		if k.Tertutup() || (s.Posisi != "" && k.PositionNote != s.Posisi) {
			continue
		}
		if s.Cari != "" && !strings.Contains(strings.ToUpper(id), strings.ToUpper(s.Cari)) {
			continue
		}
		out = append(out, models.RingkasanKasus{ID: id, StatusWork: k.StatusWork, PositionNote: k.PositionNote, NoPolis: k.NoPolis})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (g *Gudang) SimpanHalaman(_ context.Context, _ *db.Tx, id string, h *models.Halaman) error {
	if err := g.gagal("SimpanHalaman"); err != nil {
		return err
	}
	if err := models.PeriksaBentukSimpan(h); err != nil { // sama dengan repository
		return fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, err)
	}
	if _, ada := g.Kasus[id]; !ada {
		return repository.ErrKasusTidakAda
	}
	if g.Kasus[id].GenerasiTertutup {
		return repository.ErrGenerasiTertutup
	}
	s := h.Salin()
	s.BersihkanPesan()
	// PolicyNo milik SetelNomorPolis, persis repository (NOPOLIS di luar katalog).
	s.Hapus(models.HalamanPolis + ".PolicyNo")
	s.Hapus(models.JalurStsPKP)
	if pn := h.Ambil("PositionNote"); pn != "" {
		k := g.Kasus[id]
		k.PositionNote = pn
		g.Kasus[id] = k
	}
	g.Halaman[id] = s
	return nil
}

func (g *Gudang) BacaHalaman(_ context.Context, _ *db.Tx, id string) (*models.Halaman, error) {
	h, ada := g.Halaman[id]
	if !ada {
		return nil, repository.ErrKasusTidakAda
	}
	s := h.Salin()
	s.Setel(models.HalamanPolis+".PolicyNo", g.Kasus[id].NoPolis)
	s.Setel("pyID", id)
	return s, nil
}

func (g *Gudang) SetelNomorPolis(_ context.Context, _ *db.Tx, id, nopol string) error {
	k := g.Kasus[id]
	if k.NoPolis != "" {
		return repository.ErrNomorPolisSudahAda
	}
	k.NoPolis = nopol
	g.Kasus[id] = k
	return nil
}

func (g *Gudang) TerbitkanNomorPolis(_ context.Context, _ *db.Tx, h *models.Halaman, sekarang time.Time) (repository.BahanNomor, error) {
	tipe := models.TipeNomorPolis(h)
	if tipe == "" {
		return repository.BahanNomor{}, repository.ErrTipeNomorKosong
	}
	g.urutPol++
	return repository.BahanNomor{
		NoPolis:        models.RakitNomorPolis("UJI-", tipe, h.Ambil(models.HalamanPolis+".OJKBusinessID"), sekarang.Format("01.2006"), g.urutPol),
		ProductionDate: sekarang,
	}, nil
}

func (g *Gudang) HariClosing(context.Context, *db.Tx) (int, error) { return g.Closing, nil }

func (g *Gudang) DetailKontrak(_ context.Context, id string) (models.BarisKontrak, error) {
	b, ada := g.Kontrak[id]
	if !ada {
		return nil, repository.ErrDataKontrakTidakAda
	}
	return b, nil
}

// KomisiKontrak - baris `Kontrak` ber-TREATYID itu, berurut ID (sama dengan
// repository).
func (g *Gudang) KomisiKontrak(_ context.Context, treatyID string) ([]models.BarisKontrak, error) {
	var out []models.BarisKontrak
	for _, b := range g.Kontrak {
		if treatyID != "" && b["TREATYID"] == treatyID {
			out = append(out, b)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i]["ID"] < out[j]["ID"] })
	return out, nil
}

func (g *Gudang) DaftarDetailKontrak(context.Context, repository.SaringanDetail) ([]models.BarisKontrak, error) {
	var out []models.BarisKontrak
	for _, b := range g.Kontrak {
		out = append(out, b)
	}
	return out, nil
}

func (g *Gudang) IDMataUangDariNama(_ context.Context, nama string) (string, error) {
	if nama == "" {
		return "", nil
	}
	return "UJI-ID-" + nama, nil
}

func (g *Gudang) NamaMataUang(_ context.Context, id string) (string, error) {
	return strings.TrimPrefix(id, "UJI-ID-"), nil
}

func (g *Gudang) OJKGrupTreaty(context.Context, string) (string, error)   { return g.OJK, nil }
func (g *Gudang) OldIDGrupTreaty(context.Context, string) (string, error) { return "UJI-OLD", nil }
func (g *Gudang) KlienDariNama(context.Context, string) (string, error)   { return "", nil }
func (g *Gudang) StsPKPAgen(_ context.Context, sobID string) (string, error) {
	if v, ada := g.PKPAgen[sobID]; ada {
		return v, nil
	}
	return g.StsPKP, nil
}

func (g *Gudang) BisnisDariKunci(_ context.Context, kunci string) (models.BarisBisnis, error) {
	return g.Bisnis[kunci], nil
}

func (g *Gudang) MO(_ context.Context, id string) (models.BarisMO, error) {
	return models.BarisMO{ID: id, ClientID: "UJI-MKT", ClientName: "UJI-MO", TeamGroup: "1"}, nil
}

func (g *Gudang) DaftarMataUang(context.Context) ([]models.Pilihan, error) {
	return []models.Pilihan{{Nilai: "UJI-ID-IDR", Label: "IDR"}}, nil
}
func (g *Gudang) DaftarMO(context.Context) ([]models.Pilihan, error)             { return nil, nil }
func (g *Gudang) DaftarJenisSpreading(context.Context) ([]models.Pilihan, error) { return nil, nil }
func (g *Gudang) DaftarJenisReas(context.Context) ([]models.Pilihan, error)      { return nil, nil }
func (g *Gudang) PolisSerupa(context.Context, *models.Halaman) ([]string, error) {
	return g.Serupa, nil
}

func (g *Gudang) DaftarAgenHierarki(context.Context) ([]models.BarisAgen, error) {
	return append([]models.BarisAgen(nil), g.Agen...), nil
}

func (g *Gudang) AgenHierarki(_ context.Context, id string) (models.BarisAgen, bool, error) {
	for _, b := range g.Agen {
		if b.ID == id {
			return b, true, nil
		}
	}
	return models.BarisAgen{}, false, nil
}

func (g *Gudang) CatatRiwayat(_ context.Context, _ *db.Tx, r models.Riwayat) error {
	if err := g.gagal("CatatRiwayat"); err != nil {
		return err
	}
	g.Riwayat = append(g.Riwayat, r)
	return nil
}

func (g *Gudang) DaftarRiwayat(_ context.Context, idPega string) ([]models.Riwayat, error) {
	var out []models.Riwayat
	for _, r := range g.Riwayat {
		if r.IDPega == idPega {
			out = append(out, r)
		}
	}
	return out, nil
}

func (g *Gudang) DaftarPeranTempat(context.Context) ([]repository.PeranTempat, error) {
	return g.Tempat, nil
}

func (g *Gudang) NamaTampilan(_ context.Context, login string) (string, error) {
	return g.Nama[login], nil // tanpa jatuh-balik ke login ID (AC 40, 42)
}
