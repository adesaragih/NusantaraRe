// Package tiruan adalah gudang di memori untuk uji services dan handlers
// modul Endorsement Life - tanpa Oracle.
//
// ⛔ Ia meniru PERILAKU repository (urutan, saringan, penyalinan), bukan
// SQL-nya. Bentuk SQL diuji di paket repository; perilaku terhadap Oracle di
// uji bertag `db`.
//
// ⛔ Fixture berawalan `UJI-`: nol nama orang, nol nomor polis sungguhan.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

// Polis adalah satu baris `T_PREMIUM_LIST` tiruan (versi NB, versi EDM, kasus).
type Polis struct {
	ID          string
	NoPolis     string // `NO_POLIS` - terisi pada versi EDM resmi
	OldPolicyNo string // `OLD_POLICY_NO` - kasus EDM
	ProdKe      int    // 0 = kosong (dibaca 1)
	EdmType     string
	EdmNote     string
	EdmDate     string
	PLNumberEDM string
	NoEndors    string
	Status      string // `STATUSS`
	TglInput    string
	Pembuat     string
	Kepala      map[string]string
}

// Peserta adalah satu baris `T_PREMIUM_LIST_DETAIL` tiruan.
type Peserta struct {
	ID        string
	PolisID   string
	ParentID  string
	PLNumber  string
	EdmStatus string
	Nilai     map[string]string
}

// PolisWarisan adalah satu baris `JSON_POLIS` tiruan.
type PolisWarisan struct {
	IDPega   string
	NoPolis  string
	ProdKe   int // 0 = NULL
	EdmType  string
	TglInput string
	Kepala   map[string]string // properti `DATA_JSON` → nilai
}

// Spreading adalah satu baris `T_PREMIUM_LIST_SPREADING` tiruan beserta retronya.
type Spreading struct {
	ID       string
	DetailID string
	Nama     string
	Share    string
	Retro    []models.SpreadingRetro
}

// PesertaWarisan adalah satu baris `M_LIFE_PREMIUM_DETAIL` tiruan.
type PesertaWarisan struct {
	ID       string
	PLNumber string
	IDPega   string
	Nilai    map[string]string
}

// Gudang adalah tiruan `repository.Gudang`.
type Gudang struct {
	mu        sync.Mutex
	Polis     map[string]*Polis
	Peserta   []*Peserta
	Spreading []*Spreading
	Rekap     map[string]int // PREMIUM_LIST_ID → cacah baris rekap
	// RekapData - baris rekap mata uang per kasus (kunci `models.KolomRekapKasus`).
	RekapData      map[string][]map[string]string
	PolisWarisan   []*PolisWarisan
	PesertaWarisan []*PesertaWarisan
	// RekapWarisan - baris `M_LIFE_PREMIUM_SUMMARY` (kunci kolom).
	RekapWarisan []map[string]string
	// Dibayar - nomor invoice Arasapas yang punya baris pelunasan.
	Dibayar map[string]bool
	// GalatArasapas - bila terisi, `SudahDibayar` mengembalikannya.
	GalatArasapas error
	// Galat - bila terisi, setiap metode mengembalikannya (uji jalur gagal).
	Galat error
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{Polis: map[string]*Polis{}, Rekap: map[string]int{}, Dibayar: map[string]bool{}}
}

// Transaksi menjalankan fn tanpa transaksi sungguhan (uji services).
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// Inbox meniru `sqlInbox`: kasus EDM terbuka, terbaru dahulu, pemutus seri ID.
func (g *Gudang) Inbox(_ context.Context, halaman, ukuran int) ([]models.BarisInbox, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return nil, 0, g.Galat
	}
	var semua []*Polis
	for _, p := range g.Polis {
		if models.KasusEDM(p.ID) && p.Status == "" {
			semua = append(semua, p)
		}
	}
	sort.Slice(semua, func(i, j int) bool {
		if semua[i].TglInput != semua[j].TglInput {
			return semua[i].TglInput > semua[j].TglInput
		}
		return semua[i].ID < semua[j].ID
	})
	total := len(semua)
	awal := (halaman - 1) * ukuran
	if awal > total {
		awal = total
	}
	akhir := awal + ukuran
	if akhir > total {
		akhir = total
	}
	var hasil []models.BarisInbox
	for _, p := range semua[awal:akhir] {
		hasil = append(hasil, models.BarisInbox{
			CaseID: p.ID, EndorsementNo: p.OldPolicyNo, PolicyNo: p.OldPolicyNo, Tipe: p.Kepala["TYPE"],
			EdmType: p.EdmType, Sob: p.Kepala["SOB_NAME"], Ceding: p.Kepala["CEDING_CO_NAME"],
			PolicyHolder: p.Kepala["POLICY_HOLDER_NAME"], MarketingName: p.Kepala["MARKETING_NAME"],
			CreateDate: p.TglInput, CreateOperator: p.Pembuat, Status: p.Status,
		})
	}
	return hasil, total, nil
}

// AmbilKasus meniru `sqlKasus` (kunci diabaikan).
func (g *Gudang) AmbilKasus(_ context.Context, _ *db.Tx, id string, _ bool) (models.Kasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return models.Kasus{}, g.Galat
	}
	p, ada := g.Polis[id]
	if !ada || !models.KasusEDM(id) {
		return models.Kasus{}, repository.ErrTidakAda
	}
	kepala := map[string]string{}
	for _, k := range models.KolomKepalaSalin {
		kepala[k.Kolom] = p.Kepala[k.Kolom]
	}
	prod := p.ProdKe
	if prod == 0 {
		prod = 1
	}
	return models.Kasus{
		ID: p.ID, NomorPolis: p.OldPolicyNo, PLNumber: p.OldPolicyNo, PLNumberEDM: p.PLNumberEDM,
		ProdKe: prod, EdmType: p.EdmType, EdmNote: p.EdmNote, EdmDate: p.EdmDate, Status: p.Status,
		TglInput: p.TglInput, Pembuat: p.Pembuat, Kepala: kepala,
	}, nil
}

// pesertaKasus - peserta satu polis, urut seperti `sqlPeserta`.
func (g *Gudang) pesertaKasus(polisID string) []*Peserta {
	var hasil []*Peserta
	for _, d := range g.Peserta {
		if d.PolisID == polisID {
			hasil = append(hasil, d)
		}
	}
	sort.SliceStable(hasil, func(i, j int) bool {
		a, b := hasil[i].Nilai, hasil[j].Nilai
		if a["CERTIFICATE_NO"] != b["CERTIFICATE_NO"] {
			return a["CERTIFICATE_NO"] < b["CERTIFICATE_NO"]
		}
		if a["NAME_OF_INSURED"] != b["NAME_OF_INSURED"] {
			return a["NAME_OF_INSURED"] < b["NAME_OF_INSURED"]
		}
		return hasil[i].ID < hasil[j].ID
	})
	return hasil
}

// DaftarPeserta meniru `sqlPeserta`.
func (g *Gudang) DaftarPeserta(_ context.Context, kasusID string, halaman, ukuran int) ([]models.Peserta, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return nil, g.Galat
	}
	semua := g.pesertaKasus(kasusID)
	awal := (halaman - 1) * ukuran
	if awal > len(semua) {
		awal = len(semua)
	}
	akhir := awal + ukuran
	if akhir > len(semua) {
		akhir = len(semua)
	}
	var hasil []models.Peserta
	for _, d := range semua[awal:akhir] {
		nilai := map[string]string{}
		for _, k := range models.KolomPesertaGrid {
			nilai[k.Nama] = d.Nilai[k.Nama]
		}
		hasil = append(hasil, models.Peserta{
			ID: d.ID, ParentID: d.ParentID, EdmStatus: d.EdmStatus, Nilai: nilai,
			Terkunci: d.EdmStatus == models.StatusDelete || d.EdmStatus == models.StatusBatal,
		})
	}
	return hasil, nil
}

// CacahPeserta meniru `sqlCacahPeserta`.
func (g *Gudang) CacahPeserta(_ context.Context, _ *db.Tx, kasusID string) (map[string]int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return nil, g.Galat
	}
	hasil := map[string]int{}
	for _, d := range g.Peserta {
		if d.PolisID == kasusID {
			hasil[d.EdmStatus]++
		}
	}
	return hasil, nil
}

// AdaRekap meniru `IsJsonPolis`.
func (g *Gudang) AdaRekap(_ context.Context, _ *db.Tx, kasusID string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return false, g.Galat
	}
	return g.Rekap[kasusID] > 0, nil
}

// VersiBerjalan meniru ketiga jalur `repository.Gudang.VersiBerjalan`.
func (g *Gudang) VersiBerjalan(_ context.Context, _ *db.Tx, nomorPolis string, sebelum int) (models.Versi, bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return models.Versi{}, false, g.Galat
	}
	lolos := func(prod int) bool { return sebelum <= 0 || prod < sebelum }
	var terbaik models.Versi
	for _, p := range g.Polis {
		prod := p.ProdKe
		if prod == 0 {
			prod = 1
		}
		if !lolos(prod) {
			continue
		}
		resmi := p.NoPolis == nomorPolis && p.Status == models.StatusKasusSelesai
		nb := p.EdmType == "" && g.punyaPLNumber(p.ID, nomorPolis)
		if resmi || nb {
			terbaik = lebihBaruTiruan(terbaik, models.Versi{Jenis: models.SumberAplikasi, ID: p.ID, ProdKe: prod, EdmType: p.EdmType})
		}
	}
	for _, w := range g.PolisWarisan {
		prod := w.ProdKe
		if prod == 0 {
			prod = 1
		}
		if w.NoPolis == nomorPolis && lolos(prod) {
			terbaik = lebihBaruTiruan(terbaik, models.Versi{Jenis: models.SumberWarisan, ID: w.IDPega, ProdKe: prod, EdmType: w.EdmType})
		}
	}
	return terbaik, terbaik.ID != "", nil
}

// lebihBaruTiruan - `models.LebihBaru`, ditambah pemutus seri `ID` seperti
// `ORDER BY …, p.ID` supaya peta tanpa urutan tetap deterministik.
func lebihBaruTiruan(a, b models.Versi) models.Versi {
	if a.ID != "" && b.ProdKe == a.ProdKe && b.Jenis == a.Jenis && b.ID > a.ID {
		return a
	}
	if a.ID != "" && b.ProdKe == a.ProdKe && b.Jenis == a.Jenis {
		return b
	}
	return models.LebihBaru(a, b)
}

func (g *Gudang) punyaPLNumber(polisID, nomor string) bool {
	for _, d := range g.Peserta {
		if d.PolisID == polisID && d.PLNumber == nomor {
			return true
		}
	}
	return false
}

// UrutanBerikut - pengganti `SEQ_WORK_EDM_LIFE`.
func (g *Gudang) urutanBerikut() string {
	n := 0
	for id := range g.Polis {
		if strings.HasPrefix(id, models.AwalanKasus) {
			if v, err := strconv.Atoi(strings.TrimPrefix(id, models.AwalanKasus)); err == nil && v > n {
				n = v
			}
		}
	}
	return strconv.Itoa(n + 1)
}

// AdaKasusTerbuka meniru `sqlKasusTerbuka`.
func (g *Gudang) AdaKasusTerbuka(_ context.Context, _ *db.Tx, nomorPolis string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return false, g.Galat
	}
	for _, p := range g.Polis {
		if models.KasusEDM(p.ID) && p.Status == "" && p.OldPolicyNo == nomorPolis {
			return true, nil
		}
	}
	return false, nil
}

// SudahDibayar meniru `sqlSudahDibayar`.
func (g *Gudang) SudahDibayar(_ context.Context, nomorInvoice string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.GalatArasapas != nil {
		return false, g.GalatArasapas
	}
	return g.Dibayar[nomorInvoice], nil
}

// PengenalKasusBaru meniru `SEQ_WORK_EDM_LIFE`.
func (g *Gudang) PengenalKasusBaru(_ context.Context, _ *db.Tx) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return "", g.Galat
	}
	return models.RakitPengenalKasus(g.urutanBerikut())
}

// KepalaSumber meniru `sqlKepalaAplikasi` / `sqlKepalaWarisan`.
func (g *Gudang) KepalaSumber(_ context.Context, _ *db.Tx, v models.Versi, nomorPolis string) (map[string]string, []string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return nil, nil, g.Galat
	}
	kepala := map[string]string{}
	switch v.Jenis {
	case models.SumberAplikasi:
		p, ada := g.Polis[v.ID]
		if !ada {
			return nil, nil, repository.ErrTidakAda
		}
		for _, k := range models.KolomKepalaSalin {
			kepala[k.Kolom] = p.Kepala[k.Kolom]
		}
		return kepala, nil, nil
	case models.SumberWarisan:
		for _, w := range g.PolisWarisan {
			if w.IDPega == v.ID && w.NoPolis == nomorPolis {
				var rusak []string
				for _, k := range models.KolomKepalaSalin {
					s := w.Kepala[k.Properti]
					if k.Tanggal {
						iso, ok := repository.TanggalPega(s)
						if !ok {
							rusak = append(rusak, k.Kolom)
						}
						s = iso
					}
					kepala[k.Kolom] = s
				}
				return kepala, rusak, nil
			}
		}
		return nil, nil, repository.ErrTidakAda
	}
	return nil, nil, fmt.Errorf("tiruan: sumber %q", v.Jenis)
}

// SisipKasus meniru `sqlSisipKasus` + index unik `UX_PL_EDM_TERBUKA` (482).
func (g *Gudang) SisipKasus(_ context.Context, _ *db.Tx, k repository.KasusTulis) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return g.Galat
	}
	for _, p := range g.Polis {
		if models.KasusEDM(p.ID) && p.Status == "" && p.OldPolicyNo == k.NomorPolis {
			return fmt.Errorf("%w: %s", repository.ErrKasusTerbukaGanda, k.NomorPolis)
		}
	}
	kepala := map[string]string{}
	for kol, v := range k.Kepala {
		kepala[kol] = v
	}
	g.Polis[k.ID] = &Polis{
		ID: k.ID, OldPolicyNo: k.NomorPolis, EdmType: k.EdmType, EdmDate: k.EdmDate, EdmNote: k.EdmNote,
		ProdKe: k.ProdKe, Pembuat: k.Pembuat, TglInput: "2026-10-01 12:00:00", Kepala: kepala,
	}
	return nil
}

// SalinVersi meniru `sqlSalinPesertaAplikasi` (+ spreading) dan
// `sqlSalinPesertaWarisan`: baris `Delete` tidak disalin, sisanya `Old`.
func (g *Gudang) SalinVersi(_ context.Context, _ *db.Tx, kasusID string, v models.Versi, nomorPolis string) (repository.Salinan, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return repository.Salinan{}, g.Galat
	}
	var s repository.Salinan
	switch v.Jenis {
	case models.SumberAplikasi:
		var lama []*Peserta
		for _, d := range g.Peserta {
			if d.PolisID == v.ID && d.EdmStatus != models.StatusDelete {
				lama = append(lama, d)
			}
		}
		for _, d := range lama {
			nilai := map[string]string{}
			for kol, x := range d.Nilai {
				nilai[kol] = x
			}
			baru := &Peserta{ID: kasusID + "/D/" + d.ID, PolisID: kasusID, ParentID: d.ID, PLNumber: d.PLNumber,
				EdmStatus: models.StatusOld, Nilai: nilai}
			g.Peserta = append(g.Peserta, baru)
			s.Peserta++
			var salinan []*Spreading
			for _, sp := range g.Spreading {
				if sp.DetailID == d.ID {
					salinan = append(salinan, &Spreading{ID: kasusID + "/S/" + sp.ID, DetailID: baru.ID, Nama: sp.Nama,
						Share: sp.Share, Retro: append([]models.SpreadingRetro(nil), sp.Retro...)})
					s.Spreading++
					s.SpreadingRetro += len(sp.Retro)
				}
			}
			g.Spreading = append(g.Spreading, salinan...)
		}
	case models.SumberWarisan:
		if v.EdmType != "" {
			return repository.Salinan{}, repository.ErrSumberWarisanEDM
		}
		for _, m := range g.PesertaWarisan {
			if m.PLNumber == nomorPolis && m.IDPega == v.ID {
				nilai := map[string]string{}
				for kol, x := range m.Nilai {
					nilai[kol] = x
				}
				g.Peserta = append(g.Peserta, &Peserta{ID: kasusID + "/W/" + m.ID, PolisID: kasusID, PLNumber: m.PLNumber,
					EdmStatus: models.StatusOld, Nilai: nilai})
				s.Peserta++
			}
		}
	}
	return s, nil
}

// RincianPeserta meniru `sqlRincianPeserta` + `sqlSpreadingPeserta`.
func (g *Gudang) RincianPeserta(_ context.Context, kasusID, pesertaID string) (models.RincianPeserta, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return models.RincianPeserta{}, g.Galat
	}
	for _, d := range g.Peserta {
		if d.ID != pesertaID || d.PolisID != kasusID {
			continue
		}
		nilai := map[string]string{}
		for _, k := range models.KolomPesertaRinci {
			nilai[k.Nama] = d.Nilai[k.Nama]
		}
		r := models.RincianPeserta{
			Peserta: models.Peserta{ID: d.ID, ParentID: d.ParentID, EdmStatus: d.EdmStatus, Nilai: nilai,
				Terkunci: d.EdmStatus == models.StatusDelete || d.EdmStatus == models.StatusBatal},
			Spreading: []models.Spreading{},
		}
		for _, sp := range g.Spreading {
			if sp.DetailID == d.ID {
				r.Spreading = append(r.Spreading, models.Spreading{ID: sp.ID, TreatyTypeName: sp.Nama,
					RetrocadedShare: sp.Share, Retro: append([]models.SpreadingRetro{}, sp.Retro...)})
			}
		}
		return r, nil
	}
	return models.RincianPeserta{}, repository.ErrTidakAda
}

// PesertaVersi meniru `PesertaVersi` repository.
func (g *Gudang) PesertaVersi(_ context.Context, v models.Versi, nomorPolis string, halaman, ukuran int) ([]models.Peserta, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return nil, 0, g.Galat
	}
	var semua []models.Peserta
	switch v.Jenis {
	case models.SumberAplikasi:
		for _, d := range g.pesertaKasus(v.ID) {
			semua = append(semua, models.Peserta{ID: d.ID, ParentID: d.ParentID, EdmStatus: d.EdmStatus, Nilai: d.Nilai})
		}
	case models.SumberWarisan:
		for _, m := range g.PesertaWarisan {
			if m.PLNumber == nomorPolis && m.IDPega == v.ID {
				semua = append(semua, models.Peserta{ID: m.ID, Nilai: m.Nilai})
			}
		}
	}
	awal := (halaman - 1) * ukuran
	if awal > len(semua) {
		awal = len(semua)
	}
	akhir := awal + ukuran
	if akhir > len(semua) {
		akhir = len(semua)
	}
	return semua[awal:akhir], len(semua), nil
}

// RekapPolis meniru `sqlRekapPolis` (seluruh versi bernomor PL itu).
func (g *Gudang) RekapPolis(_ context.Context, nomorPolis string) ([]map[string]string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return nil, g.Galat
	}
	var hasil []map[string]string
	for _, r := range g.RekapWarisan {
		if r["PL_NUMBER"] == nomorPolis {
			hasil = append(hasil, r)
		}
	}
	return hasil, nil
}
