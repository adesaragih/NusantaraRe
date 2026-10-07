// Package tiruan adalah gudang EDM Treaty In di memori - untuk uji seam HTTP (handler) tanpa Oracle. Ia memenuhi
// `services.Gudang` dengan perilaku yang PENTING bagi uji: transaksi yang benar-benar batal (potret dipulihkan bila
// fn gagal), tahap lama di syarat pindah tahap, UNIQUE OLD_POLIS_ID dan (NOPOLIS, PRODKE) bernomor, generasi
// tertutup tidak dapat disunting, proyeksi selisih SUMBER 'PEGA' beku. Pola: salinan
// `modul/nbtreatyin/backend/tiruan/gudang.go` (06-10-2026), isinya menurut rantai generasi EDM.
//
// ⛔ Data uji berawalan UJI-: nol data polis atau nama orang sungguhan.
package tiruan

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/repository"
)

// Generasi - satu baris T_GENERAL_POLIS_TREATY beserta anaknya (halaman = proyeksi katalog).
type Generasi struct {
	NoPolis    string
	ProdKe     int
	EDMNo      string
	OldPolisID string
	EDMType    string
	Halaman    *models.Halaman
}

// Selisih - satu baris T_POLIS_DIFFERENCE beserta anaknya (halaman berjalur PolicyTreatyIn.TreatyDifference.* dan
// TreatyXOLDifferenceList).
type Selisih struct {
	Kunci   repository.KunciSelisih
	Halaman *models.Halaman
}

// Gudang adalah penyimpanan di memori.
type Gudang struct {
	mu sync.Mutex

	urut     int
	Kasus    map[string]models.Kasus
	Generasi map[string]*Generasi
	Selisih  map[string]*Selisih
	Riwayat  []models.Riwayat
	Usulan   []BarisRiwayatProduksi
	// Simpanan - Utility1 (`SimpanPolisProduksi`), urutan tulis.
	Simpanan []models.SimpananPolis
	// JSONPolis - cacah generasi berbeda per nomor polis di POOLDATA.JSON_POLIS (CheckNopolisAvailability = cacah > 0;
	// penjaga Create = cacah tidak melebihi generasi terakhir + 1). Bertambah per generasi ditanam / selesai.
	JSONPolis map[string]int
	// NoMaster - NOOFFER TREATYINPRODUCTION per nomor polis (FetchNoOfferFromNoPolis).
	NoMaster map[string]string

	Nama        map[string]string // login -> nama tampilan
	NamaPembuat map[string]string
	Bisnis      map[string]models.BarisBisnis
	PKPAgen     map[string]string
	StsPKP      string
	Closing     int
	// Popup - baris grid BusinessAndSOBListEDM yang dijawab DaftarBisnisEDM; SaringanPopup - masukan terakhir.
	Popup         []models.Baris
	SaringanPopup repository.SaringanPopupEDM
	// Master - master kontrak per kunci: "ID:<id>", "OLDID:<id>", "EDM:<id>", "OUT:<id>".
	Master     map[string][]models.MasterXOL
	MataUang   map[string]string // nama -> ID
	KotakMasuk map[string]models.PemegangKotakMasuk
	GagalDi    string
	Panggil    []string
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{
		Kasus: map[string]models.Kasus{}, Generasi: map[string]*Generasi{}, Selisih: map[string]*Selisih{},
		JSONPolis: map[string]int{}, NoMaster: map[string]string{}, Nama: map[string]string{},
		NamaPembuat: map[string]string{}, Bisnis: map[string]models.BarisBisnis{}, Closing: 25,
		Master: map[string][]models.MasterXOL{}, MataUang: map[string]string{},
	}
}

// BarisRiwayatProduksi - satu baris HISTORYAKSEPTASIPRODUCTION beserta IDPEGA.
type BarisRiwayatProduksi struct {
	IDPega string
	models.UsulanProduksi
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

// TanamPolis menanam satu generasi polis yang SUDAH selesai (NOPOLIS terisi) - polis NB (prodKe 0) atau endorsemen
// lama; `noMaster` = NOOFFER TREATYINPRODUCTION. Jawab ID generasinya.
func (g *Gudang) TanamPolis(id, nopolis string, prodKe int, oldPolisID string, h *models.Halaman, noMaster string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	h = models.ProyeksiKatalog(h)
	h.Setel(models.HalamanPolis+".PolicyNo", nopolis)
	g.Generasi[id] = &Generasi{NoPolis: nopolis, ProdKe: prodKe, OldPolisID: oldPolisID, Halaman: h}
	g.JSONPolis[nopolis]++
	if noMaster != "" {
		g.NoMaster[nopolis] = noMaster
	}
	return id
}

type potret struct {
	urut     int
	kasus    map[string]models.Kasus
	generasi map[string]Generasi
	selisih  map[string]Selisih
	riwayat  []models.Riwayat
	usulan   []BarisRiwayatProduksi
	simpanan []models.SimpananPolis
}

func (g *Gudang) potret() potret {
	p := potret{urut: g.urut, kasus: map[string]models.Kasus{}, generasi: map[string]Generasi{}, selisih: map[string]Selisih{}}
	for k, v := range g.Kasus {
		p.kasus[k] = v
	}
	for k, v := range g.Generasi {
		s := *v
		s.Halaman = v.Halaman.Salin()
		p.generasi[k] = s
	}
	for k, v := range g.Selisih {
		s := *v
		s.Halaman = v.Halaman.Salin()
		p.selisih[k] = s
	}
	p.riwayat = append(p.riwayat, g.Riwayat...)
	p.usulan = append(p.usulan, g.Usulan...)
	p.simpanan = append(p.simpanan, g.Simpanan...)
	return p
}

func (g *Gudang) pulihkan(p potret) {
	g.urut, g.Kasus = p.urut, p.kasus
	g.Generasi = map[string]*Generasi{}
	for k, v := range p.generasi {
		s := v
		g.Generasi[k] = &s
	}
	g.Selisih = map[string]*Selisih{}
	for k, v := range p.selisih {
		s := v
		g.Selisih[k] = &s
	}
	g.Riwayat, g.Usulan, g.Simpanan = p.riwayat, p.usulan, p.simpanan
}

// Transaksi menjalankan fn; galat memulihkan seluruh keadaan.
func (g *Gudang) Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error {
	g.mu.Lock()
	sebelum := g.potret()
	g.mu.Unlock()
	err := fn(&db.Tx{})
	g.mu.Lock()
	defer g.mu.Unlock()
	if err != nil {
		g.pulihkan(sebelum)
	}
	return err
}

// ------------------------------------------------------------------ kasus

func (g *Gudang) IDKasusBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.urut++
	return models.RakitIDKasus(fmt.Sprint(9000 + g.urut)), nil
}

func (g *Gudang) SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string, gb repository.GenerasiBaru) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.gagal("SisipKasus"); err != nil {
		return err
	}
	if gb.OldPolisID != "" {
		for _, x := range g.Generasi {
			if x.OldPolisID == gb.OldPolisID {
				return repository.ErrGenerasiSudahDiendorse
			}
		}
	}
	g.Kasus[id] = models.Kasus{ID: id, Position: models.PositionAdmin, StatusWork: models.AssignmentAdmin,
		PositionNote: models.PosisiAdmin, ProdKe: gb.ProdKe, OldPolisID: gb.OldPolisID, CreateOp: pembuat,
		TglCreate: "2026-10-06 10:00:00"}
	g.NamaPembuat[id] = namaPembuat
	g.Generasi[id] = &Generasi{ProdKe: gb.ProdKe, EDMNo: gb.EDMNo, OldPolisID: gb.OldPolisID, EDMType: gb.EDMType,
		Halaman: models.HalamanBaru()}
	return nil
}

// tertutup - generasi `id` punya penerus (ID-10).
func (g *Gudang) tertutup(id string) bool {
	for _, x := range g.Generasi {
		if x.OldPolisID == id {
			return true
		}
	}
	return false
}

func (g *Gudang) keadaan(id string) (models.Kasus, error) {
	k, ada := g.Kasus[id]
	if !ada {
		return models.Kasus{}, repository.ErrKasusTidakAda
	}
	gen := g.Generasi[id]
	k.NoPolis, k.OldPolisID, k.ProdKe = gen.NoPolis, gen.OldPolisID, gen.ProdKe
	k.GenerasiTertutup = g.tertutup(id)
	return k, nil
}

func (g *Gudang) Keadaan(ctx context.Context, tx *db.Tx, id string) (models.Kasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.keadaan(id)
}

func (g *Gudang) KunciKasus(ctx context.Context, tx *db.Tx, id, statusHarap string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ada := g.Kasus[id]
	if !ada {
		return repository.ErrKasusTidakAda
	}
	if k.StatusWork != statusHarap {
		return repository.ErrTahapBerubah
	}
	return nil
}

func (g *Gudang) PindahPosisi(ctx context.Context, tx *db.Tx, id, statusLama, posisiBaru, position string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.gagal("PindahPosisi"); err != nil {
		return err
	}
	k := g.Kasus[id]
	if k.StatusWork != statusLama {
		return repository.ErrTahapBerubah
	}
	if g.tertutup(id) {
		return repository.ErrGenerasiTertutup
	}
	if position != "" {
		k.Position = position
	}
	k.StatusWork, k.PositionNote = models.AssignmentPosisi(posisiBaru), posisiBaru
	g.Kasus[id] = k
	g.Generasi[id].Halaman.Setel("PositionNote", posisiBaru)
	return nil
}

func (g *Gudang) TutupKasus(ctx context.Context, tx *db.Tx, id, statusLama, statusAkhir string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
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

func (g *Gudang) DaftarKasus(ctx context.Context, s models.SaringanKasus) ([]models.RingkasanKasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []models.RingkasanKasus{}
	for id, k := range g.Kasus {
		h := g.Generasi[id].Halaman
		selesai := k.StatusWork == models.StatusDitolak || k.StatusWork == models.StatusSelesai
		if selesai != s.Selesai {
			continue
		}
		pt := func(m string) string { return h.Ambil(models.HalamanPolis + "." + m) }
		if !models.CocokCari(s.Cari, id, pt("NoOffer"), g.Generasi[id].NoPolis, models.NilaiQuotation(h, "OldPolicyNo"),
			g.Generasi[id].EDMNo, pt("InsuredName"), models.NilaiQuotation(h, "InsuredName"), models.NilaiQuotation(h, "BusinessName"),
			pt("SOBName"), pt("CedingCoName"), models.NilaiQuotation(h, "MarketingName"), pt("TreatyGroupName"), pt("BizName"),
			g.NamaPembuat[id]) {
			continue
		}
		if s.Posisi != "" && k.PositionNote != s.Posisi {
			continue
		}
		pembuat := s.Pembuat == "" || (k.CreateOp == s.Pembuat && (s.PembuatPosisi == "" || k.PositionNote == s.PembuatPosisi))
		antrean := false
		for _, a := range s.Antrean {
			antrean = antrean || a == k.PositionNote
		}
		switch {
		case s.Pembuat != "" && len(s.Antrean) > 0:
			if !pembuat && !antrean {
				continue
			}
		case s.Pembuat != "":
			if !pembuat {
				continue
			}
		case len(s.Antrean) > 0:
			if !antrean {
				continue
			}
		}
		np := g.Generasi[id].NoPolis
		if np == "" {
			np = models.NilaiQuotation(h, "OldPolicyNo")
		}
		out = append(out, models.RingkasanKasus{ID: id, NoOffer: h.Ambil(models.HalamanPolis + ".NoOffer"), NoPolis: np,
			EDMNo: g.Generasi[id].EDMNo, SOBName: h.Ambil(models.HalamanPolis + ".SOBName"),
			CedingCoName: h.Ambil(models.HalamanPolis + ".CedingCoName"), EDMType: g.Generasi[id].EDMType,
			ProportionalType: models.NilaiQuotation(h, "ProportionalType"), MarketingName: models.NilaiQuotation(h, "MarketingName"),
			NBStatus: h.Ambil("NBStatus"), StatusWork: k.StatusWork, PositionNote: k.PositionNote, TglCreate: k.TglCreate,
			TglProd: h.Ambil(models.HalamanPolis + ".ProductionDate")})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (g *Gudang) HitungKotakMasuk(ctx context.Context, akun string, admin bool, atasan []string) (map[string]int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := map[string]int{}
	for _, k := range g.Kasus {
		if k.StatusWork == models.StatusDitolak || k.StatusWork == models.StatusSelesai {
			continue
		}
		if admin && k.PositionNote == models.PosisiAdmin && k.CreateOp == akun {
			out[k.PositionNote]++
			continue
		}
		for _, a := range atasan {
			if a == k.PositionNote {
				out[k.PositionNote]++
			}
		}
	}
	return out, nil
}

// ------------------------------------------------------------------ generasi

func (g *Gudang) CacahGenerasiJSONPolis(ctx context.Context, nopolis string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.JSONPolis[nopolis], nil
}

func (g *Gudang) NoMasterDariNoPolis(ctx context.Context, nopolis string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.NoMaster[nopolis], nil
}

func (g *Gudang) AdaEDMBerjalan(ctx context.Context, nopolis string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for id, k := range g.Kasus {
		if k.StatusWork == models.StatusDitolak || k.StatusWork == models.StatusSelesai {
			continue
		}
		if models.NilaiQuotation(g.Generasi[id].Halaman, "OldPolicyNo") == nopolis {
			return true, nil
		}
	}
	return false, nil
}

func (g *Gudang) GenerasiTerakhir(ctx context.Context, tx *db.Tx, nopolis string) (repository.GenerasiPolis, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var r repository.GenerasiPolis
	ada := false
	for id, x := range g.Generasi {
		if x.NoPolis == nopolis && (!ada || x.ProdKe > r.ProdKe) {
			r, ada = repository.GenerasiPolis{ID: id, ProdKe: x.ProdKe, EDMNo: x.EDMNo}, true
		}
	}
	if !ada {
		return r, repository.ErrGenerasiPolisTidakAda
	}
	return r, nil
}

func (g *Gudang) SetelNomorPolisSelesai(ctx context.Context, tx *db.Tx, id, nopolis string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	gen := g.Generasi[id]
	if gen == nil || gen.NoPolis != "" {
		return repository.ErrNomorPolisSudahAda
	}
	for _, x := range g.Generasi {
		if x.NoPolis == nopolis && x.ProdKe == gen.ProdKe {
			return fmt.Errorf("tiruan: ORA-00001 UQ_GP_TREATY_NOPOLIS (%s, %d)", nopolis, gen.ProdKe)
		}
	}
	gen.NoPolis = nopolis
	g.JSONPolis[nopolis]++
	return nil
}

func (g *Gudang) LepasGenerasi(ctx context.Context, tx *db.Tx, id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if gen := g.Generasi[id]; gen != nil && gen.NoPolis == "" {
		gen.OldPolisID = ""
	}
	return nil
}

// ------------------------------------------------------------------ halaman

func (g *Gudang) SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.gagal("SimpanHalaman"); err != nil {
		return err
	}
	if err := models.PeriksaBentukSimpan(h); err != nil {
		return err
	}
	gen := g.Generasi[id]
	if gen == nil {
		return repository.ErrKasusTidakAda
	}
	if g.tertutup(id) {
		return repository.ErrGenerasiTertutup
	}
	gen.Halaman = models.ProyeksiKatalog(h)
	if s := h.AmbilDaftar(models.DaftarUsulan); s != nil {
		gen.Halaman.SetelDaftar(models.DaftarUsulan, nil)
	}
	return nil
}

// bacaGenerasi = repository.BacaGenerasi: proyeksi katalog + kunci generasi.
func (g *Gudang) bacaGenerasi(id string) (*models.Halaman, error) {
	gen := g.Generasi[id]
	if gen == nil {
		return nil, repository.ErrKasusTidakAda
	}
	h := gen.Halaman.Salin()
	h.Setel(models.HalamanPolis+".PolicyNo", gen.NoPolis)
	if gen.EDMNo != "" {
		h.Setel(models.HalamanPolis+".EDMNo", gen.EDMNo)
	}
	h.Setel(models.HalamanPolis+".ProdKe", fmt.Sprint(gen.ProdKe))
	if gen.EDMType != "" && h.Ambil(models.HalamanPolis+".EDMType") == "" {
		h.Setel(models.HalamanPolis+".EDMType", gen.EDMType)
	}
	return h, nil
}

func (g *Gudang) BacaGenerasi(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.bacaGenerasi(id)
}

// bacaSelisih = repository.BacaSelisih.
func (g *Gudang) bacaSelisih(polisID string, h *models.Halaman, awalan string) error {
	s := g.Selisih[polisID]
	if s == nil {
		return nil
	}
	sem := s.Halaman.Salin()
	if lapisan := models.DatarSelisihLapisan(sem); len(lapisan) > 0 {
		if err := models.BangunIndukSelisihXOL(sem, lapisan); err != nil {
			return err
		}
	}
	models.PindahAwalan(sem, h, models.HalamanPolis+".", models.HalamanPolis+"."+awalan)
	return nil
}

// BacaHalaman = repository.BacaHalaman (halaman_edm.go).
func (g *Gudang) BacaHalaman(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, err := g.keadaan(id)
	if err != nil {
		return nil, err
	}
	h, err := g.bacaGenerasi(id)
	if err != nil {
		return nil, err
	}
	nopolis := h.Ambil(models.HalamanPolis + ".PolicyNo")
	if nopolis == "" {
		nopolis = models.NilaiQuotation(h, "OldPolicyNo")
	}
	lama := k.OldPolisID
	if lama == "" {
		for gid, x := range g.Generasi {
			if nopolis != "" && x.NoPolis == nopolis && x.ProdKe == k.ProdKe-1 {
				lama = gid
			}
		}
	}
	if lama != "" {
		hl, err := g.bacaGenerasi(lama)
		if err != nil {
			return nil, err
		}
		if err := g.bacaSelisih(lama, hl, ""); err != nil {
			return nil, err
		}
		models.PasangOldData(h, hl)
	}
	if h.Ambil(models.HalamanPolis+".PolicyNo") == "" {
		h.Setel(models.HalamanPolis+".PolicyNo", h.Ambil(models.HalamanPolis+".OldData.PolicyNo"))
	}
	if err := g.bacaSelisih(id, h, ""); err != nil {
		return nil, err
	}
	var catatan []models.Baris
	for _, u := range g.Usulan {
		if u.IDPega == models.KunciInstans(id) {
			catatan = append(catatan, models.BarisCatatan(u.UsulanProduksi))
		}
	}
	if len(catatan) > 0 {
		h.SetelDaftar(models.DaftarUsulan, catatan)
	}
	h.Setel("pyID", id)
	return h, nil
}

// BuangSelisihGenerasi = repository.BuangSelisihGenerasi: hanya baris GO.
func (g *Gudang) BuangSelisihGenerasi(ctx context.Context, tx *db.Tx, polisID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	s := g.Selisih[polisID]
	if s == nil {
		return nil
	}
	if s.Kunci.Sumber != models.SumberGo {
		return repository.ErrSelisihBeku
	}
	delete(g.Selisih, polisID)
	return nil
}

// SimpanSelisih = repository.SimpanSelisih: hanya kolom katalog selisih; baris PEGA beku.
func (g *Gudang) SimpanSelisih(ctx context.Context, tx *db.Tx, polisID string, h *models.Halaman, k repository.KunciSelisih) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.gagal("SimpanSelisih"); err != nil {
		return err
	}
	if lama := g.Selisih[polisID]; lama != nil && lama.Kunci.Sumber == models.SumberPega {
		return repository.ErrSelisihBeku
	}
	s := models.HalamanBaru()
	for _, kol := range models.TabelSelisih.Kolom {
		if v := h.Ambil(kol.Properti); v != "" {
			s.Setel(kol.Properti, v)
		}
	}
	saring := func(t models.Tabel, b []models.Baris) []models.Baris {
		var out []models.Baris
		for _, x := range b {
			nb := models.Baris{}
			for _, kol := range t.Kolom {
				if v := x[kol.Properti]; v != "" {
					nb[kol.Properti] = v
				}
			}
			out = append(out, nb)
		}
		return out
	}
	s.SetelDaftar(models.DaftarSelisihSpreading, saring(models.TabelSelisihSpreading, h.AmbilDaftar(models.DaftarSelisihSpreading)))
	s.SetelDaftar(models.DaftarSelisihAngsuran, saring(models.TabelSelisihAngsuran, h.AmbilDaftar(models.DaftarSelisihAngsuran)))
	lapisan := saring(models.TabelSelisihLapisan, models.DatarSelisihLapisan(h))
	if len(lapisan) > 0 {
		if err := models.BangunIndukSelisihXOL(s, lapisan); err != nil {
			return err
		}
	}
	g.Selisih[polisID] = &Selisih{Kunci: k, Halaman: s}
	return nil
}

func (g *Gudang) HariClosing(ctx context.Context, tx *db.Tx) (int, error) { return g.Closing, nil }

// ------------------------------------------------------------------ acuan

func (g *Gudang) DaftarBisnisEDM(ctx context.Context, s repository.SaringanPopupEDM) ([]models.Baris, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.SaringanPopup = s
	return append([]models.Baris{}, g.Popup...), nil
}

func (g *Gudang) IDMataUangDariNama(ctx context.Context, nama string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.MataUang[nama], nil
}

func (g *Gudang) NamaMataUang(ctx context.Context, id string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for n, i := range g.MataUang {
		if i == id {
			return n, nil
		}
	}
	return "", nil
}

func (g *Gudang) StsPKPAgen(ctx context.Context, sobID string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if v, ada := g.PKPAgen[sobID]; ada {
		return v, nil
	}
	return g.StsPKP, nil
}

func (g *Gudang) BisnisDariKunci(ctx context.Context, kunci string) (models.BarisBisnis, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.Bisnis[kunci], nil
}

func (g *Gudang) MO(ctx context.Context, id string) (models.BarisMO, error) {
	return models.BarisMO{ID: id, ClientName: "UJI-MO-" + id}, nil
}

func pilihanUji(n string) []models.Pilihan {
	return []models.Pilihan{{Nilai: "UJI-" + n, Label: "UJI " + n}}
}

func (g *Gudang) DaftarMataUang(ctx context.Context) ([]models.Pilihan, error) {
	return pilihanUji("IDR"), nil
}
func (g *Gudang) DaftarMO(ctx context.Context) ([]models.Pilihan, error) {
	return pilihanUji("MO"), nil
}
func (g *Gudang) DaftarJenisReas(ctx context.Context) ([]models.Pilihan, error) {
	return pilihanUji("QS"), nil
}

// ------------------------------------------------------------------ master

func (g *Gudang) MasterXOL(ctx context.Context, noKontrak string) (models.MasterXOL, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	m := g.Master["ID:"+noKontrak]
	if len(m) == 0 {
		return models.MasterXOL{}, repository.ErrMasterXOLTidakAda
	}
	return m[len(m)-1], nil
}

// MasterEDM - tiruan `models.PembacaMasterEDM`.
func (g *Gudang) MasterEDM(ctx context.Context) models.PembacaMasterEDM { return pembacaMaster{g} }

type pembacaMaster struct{ g *Gudang }

func (p pembacaMaster) ambil(k string) ([]models.MasterXOL, error) {
	p.g.mu.Lock()
	defer p.g.mu.Unlock()
	return append([]models.MasterXOL{}, p.g.Master[k]...), nil
}

func (p pembacaMaster) MasterMenurutOldID(id string) ([]models.MasterXOL, error) {
	return p.ambil("OLDID:" + id)
}
func (p pembacaMaster) MasterMenurutID(id string) ([]models.MasterXOL, error) {
	return p.ambil("ID:" + id)
}
func (p pembacaMaster) MasterEDMMenurutID(id string) ([]models.MasterXOL, error) {
	return p.ambil("EDM:" + id)
}
func (p pembacaMaster) MasterOutMenurutID(id string) ([]models.MasterXOL, error) {
	return p.ambil("OUT:" + id)
}
func (p pembacaMaster) IDMataUang(nama string) (string, error) {
	p.g.mu.Lock()
	defer p.g.mu.Unlock()
	return p.g.MataUang[nama], nil
}

// ------------------------------------------------------------------ riwayat dan produksi

func (g *Gudang) SimpanPolisProduksi(ctx context.Context, tx *db.Tx, s models.SimpananPolis) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.gagal("SimpanPolisProduksi"); err != nil {
		return err
	}
	g.Simpanan = append(g.Simpanan, s)
	return nil
}

func (g *Gudang) CatatRiwayat(ctx context.Context, tx *db.Tx, r models.Riwayat) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.gagal("CatatRiwayat"); err != nil {
		return err
	}
	g.Riwayat = append(g.Riwayat, r)
	return nil
}

func (g *Gudang) CatatUsulan(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := 0
	for _, u := range g.Usulan {
		if u.IDPega == idPega {
			n++
		}
	}
	for _, b := range baris {
		n++
		b.NoUrut = n
		g.Usulan = append(g.Usulan, BarisRiwayatProduksi{IDPega: idPega, UsulanProduksi: b})
	}
	return nil
}

func (g *Gudang) NamaTampilan(ctx context.Context, loginID string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if n, ada := g.Nama[loginID]; ada {
		return n, nil
	}
	return loginID, nil
}

func (g *Gudang) PemegangKotakMasuk(ctx context.Context, workbasket string) (models.PemegangKotakMasuk, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.KotakMasuk[workbasket], nil
}
