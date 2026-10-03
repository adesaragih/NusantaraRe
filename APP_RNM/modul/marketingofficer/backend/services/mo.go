package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/marketingofficer/backend/models"
	"nusantarare/modul/marketingofficer/backend/repository"
)

// Lebar kolom VARCHAR2 dalam BYTE (katalog DEV `POOLDATA.MARKETINGOFFICER` 03-10-2026).
const (
	// LebarAksesLogin - `AKSES_LOGIN` VARCHAR2(50); `M_LOGIN_GO.LOGIN_ID` bisa sampai 64.
	LebarAksesLogin = 50
	// LebarTeks - ID, CLIENTID, CLIENTNAME, CLIENTID2, MOLEADER, TEAMGROUP, USERUPDATE, BRANCHDETAILID/NAME.
	LebarTeks = 100
	// LebarBranchParent - `BRANCHPARENT` VARCHAR2(10).
	LebarBranchParent = 10
)

// Galat layanan. Pesan untuk layar (bahasa Inggris) mengikuti sesudah `: `.
var (
	// ErrTidakAda - baris MO tidak ada (404).
	ErrTidakAda = errors.New("marketing officer not found")
	// ErrMasukanTidakSah - isian ditolak (422).
	ErrMasukanTidakSah = errors.New("invalid input")
	// ErrSudahAktif - orang atau akun itu sudah punya baris AKTIF lain (409).
	ErrSudahAktif = errors.New("already active")
	// ErrTanpaPelaku - permintaan tanpa akun pelaku (401).
	ErrTanpaPelaku = errors.New("user account is required")
)

func tolak(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrMasukanTidakSah, fmt.Sprintf(format, a...))
}

// Transaksi menjalankan fn di dalam satu transaksi dan menutupnya.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang adalah semua yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	DaftarMO(ctx context.Context) ([]models.MarketingOfficer, error)
	AmbilMO(ctx context.Context, tx *dbTx, id string) (models.MarketingOfficer, error)
	SisipMO(ctx context.Context, tx *dbTx, m models.MarketingOfficer) (string, error)
	PerbaruiMO(ctx context.Context, tx *dbTx, m models.MarketingOfficer) error
	AktifDenganClientID(ctx context.Context, tx *dbTx, clientID, kecuali string) ([]string, error)
	AktifDenganAkses(ctx context.Context, tx *dbTx, login, kecuali string) ([]string, error)
	ClientIDDariAkses(ctx context.Context, tx *dbTx, login string) (string, bool, error)
	DaftarAkun(ctx context.Context) ([]models.Akun, error)
	AmbilAkun(ctx context.Context, tx *dbTx, login string) (models.Akun, error)
	DaftarCabang(ctx context.Context) ([]models.Cabang, error)
	AmbilCabang(ctx context.Context, tx *dbTx, id string) (models.Cabang, error)
	// LogMO - MARKETINGOFFICER_LOG satu MO, tertua lebih dulu.
	LogMO(ctx context.Context, id string) ([]models.BarisLog, error)
	// TandaiLog - AKSES_LOGIN lama + tanda UPDATE-GO pada baris log UPDATE yang baru terjadi (transaksi sama).
	TandaiLog(ctx context.Context, tx *dbTx, id, aksesLama string) error
}

// Layanan memegang seluruh aturan modul ini di atas satu Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
}

// BaruLayanan menyusun Layanan - dipakai uji dengan gudang tiruan dan transaksi tiruan (fn(nil)).
func BaruLayanan(g Gudang, tx Transaksi) *Layanan { return &Layanan{gudang: g, tx: tx} }

// Status akun sebuah baris di daftar.
const (
	AkunAktif    = "aktif"
	AkunNonaktif = "nonaktif"
	// AkunTidakAda - AKSES_LOGIN lama (Operator ID Pega) yang tidak ada di M_LOGIN_GO: tetap ditampilkan, ditandai.
	AkunTidakAda = "tidak-ada"
)

// BarisMO adalah satu baris daftar: baris MO dan status akunnya.
type BarisMO struct {
	models.MarketingOfficer
	// StatusAkun - kosong (tanpa akun), `aktif`, `nonaktif`, atau `tidak-ada`.
	StatusAkun string `json:"statusAkun"`
	// EmailAkun - `M_LOGIN_GO.EMAIL` akun itu; kosong bila tidak ada.
	EmailAkun string `json:"emailAkun"`
}

// OpsiLeader adalah satu pilihan dropdown "Leader" (`SelectLeader_RD`: `ClientID2 = LEADER`, aktif).
type OpsiLeader struct {
	ID        string `json:"id"`
	Nama      string `json:"nama"`
	SubBranch string `json:"subBranch"`
}

// Pilihan adalah isi pilihan form.
type Pilihan struct {
	// Akun - akun M_LOGIN_GO yang AKTIF.
	Akun   []models.Akun `json:"akun"`
	Leader []OpsiLeader  `json:"leader"`
	// Branch - cabang yang menjadi induk (ID-nya muncul sebagai BRANCHPARENTID).
	Branch []models.Cabang `json:"branch"`
	// SubBranch - cabang aktif; `induk` menyaringnya per Branch.
	SubBranch []models.Cabang `json:"subBranch"`
}

func rapikan(isi models.Isian) models.Isian {
	isi.AksesLogin = strings.TrimSpace(isi.AksesLogin)
	isi.LeaderID = strings.TrimSpace(isi.LeaderID)
	isi.BranchParent = strings.TrimSpace(isi.BranchParent)
	isi.BranchDetailID = strings.TrimSpace(isi.BranchDetailID)
	return isi
}

func status(aktif bool) string {
	if aktif {
		return models.StatusAktif
	}
	return models.StatusNonaktif
}

// periksaPelaku - akun pelaku ditulis ke `USERUPDATE` VARCHAR2(100).
func periksaPelaku(p inti.Pelaku) error {
	if strings.TrimSpace(p.AkunID) == "" {
		return ErrTanpaPelaku
	}
	if len(p.AkunID) > LebarTeks {
		return tolak("your user account is longer than %d bytes", LebarTeks)
	}
	return nil
}

// Daftar membaca seluruh baris MO beserta status akunnya.
func (l *Layanan) Daftar(ctx context.Context) ([]BarisMO, error) {
	baris, err := l.gudang.DaftarMO(ctx)
	if err != nil {
		return nil, err
	}
	akun, err := l.gudang.DaftarAkun(ctx)
	if err != nil {
		return nil, err
	}
	peta := map[string]models.Akun{}
	for _, a := range akun {
		peta[strings.ToUpper(a.LoginID)] = a
	}
	out := make([]BarisMO, 0, len(baris))
	for _, m := range baris {
		b := BarisMO{MarketingOfficer: m}
		if m.AksesLogin != "" {
			if a, ada := peta[strings.ToUpper(m.AksesLogin)]; !ada {
				b.StatusAkun = AkunTidakAda
			} else {
				b.StatusAkun, b.EmailAkun = AkunNonaktif, a.Email
				if a.Aktif {
					b.StatusAkun = AkunAktif
				}
			}
		}
		out = append(out, b)
	}
	return out, nil
}

// Pilihan menyusun isi dropdown form.
func (l *Layanan) Pilihan(ctx context.Context) (Pilihan, error) {
	p := Pilihan{Akun: []models.Akun{}, Leader: []OpsiLeader{}, Branch: []models.Cabang{}, SubBranch: []models.Cabang{}}
	akun, err := l.gudang.DaftarAkun(ctx)
	if err != nil {
		return p, err
	}
	for _, a := range akun {
		if a.Aktif {
			p.Akun = append(p.Akun, a)
		}
	}
	baris, err := l.gudang.DaftarMO(ctx)
	if err != nil {
		return p, err
	}
	for _, m := range baris {
		if m.Leader() && m.Aktif() {
			p.Leader = append(p.Leader, OpsiLeader{ID: m.ID, Nama: m.ClientName, SubBranch: m.BranchDetailName})
		}
	}
	sort.SliceStable(p.Leader, func(i, j int) bool { return strings.ToUpper(p.Leader[i].Nama) < strings.ToUpper(p.Leader[j].Nama) })
	cabang, err := l.gudang.DaftarCabang(ctx)
	if err != nil {
		return p, err
	}
	p.SubBranch = cabang
	induk := map[string]bool{}
	for _, c := range cabang {
		if c.Induk != "" {
			induk[c.Induk] = true
		}
	}
	for _, c := range cabang {
		if induk[c.ID] {
			p.Branch = append(p.Branch, c)
			delete(induk, c.ID)
		}
	}
	sisa := make([]string, 0, len(induk))
	for id := range induk {
		sisa = append(sisa, id)
	}
	sort.Strings(sisa)
	for _, id := range sisa {
		p.Branch = append(p.Branch, models.Cabang{ID: id, Nama: id})
	}
	return p, nil
}

// akunSah - akun M_LOGIN_GO yang ada, aktif, dan muat di AKSES_LOGIN.
func (l *Layanan) akunSah(ctx context.Context, tx *dbTx, login string) (models.Akun, error) {
	a, err := l.gudang.AmbilAkun(ctx, tx, login)
	if errors.Is(err, repository.ErrTidakAda) {
		return a, tolak("Login Account %s is not in the user list (Kelola User)", login)
	}
	if err != nil {
		return a, err
	}
	if !a.Aktif {
		return a, tolak("Login Account %s is not active", login)
	}
	if len(a.LoginID) > LebarAksesLogin {
		return a, tolak("Login Account %s is longer than %d characters and does not fit the marketing officer table",
			login, LebarAksesLogin)
	}
	return a, nil
}

// isiLeader mengisi CLIENTID2 dan MOLEADER. `lama` nil = baris baru; leader yang TIDAK diubah tidak diperiksa ulang
// (leader yang belakangan dinonaktifkan tidak boleh mengunci perubahan lain).
func (l *Layanan) isiLeader(ctx context.Context, tx *dbTx, m *models.MarketingOfficer, isi models.Isian, lama *models.MarketingOfficer) error {
	if isi.Leader {
		if lama != nil && lama.Leader() {
			m.ClientID2, m.MOLeader = lama.ClientID2, lama.MOLeader
			return nil
		}
		m.ClientID2, m.MOLeader = models.NilaiLeader, m.ClientName
		return nil
	}
	if isi.LeaderID == "" {
		return tolak("Leader is required, or tick Set as a leader")
	}
	if lama != nil && !lama.Leader() && isi.LeaderID == lama.ClientID2 {
		m.ClientID2, m.MOLeader = lama.ClientID2, lama.MOLeader
		return nil
	}
	if lama != nil && isi.LeaderID == lama.ID {
		return tolak("A marketing officer cannot be their own leader")
	}
	ld, err := l.gudang.AmbilMO(ctx, tx, isi.LeaderID)
	if errors.Is(err, repository.ErrTidakAda) {
		return tolak("Leader %s not found", isi.LeaderID)
	}
	if err != nil {
		return err
	}
	if !ld.Leader() || !ld.Aktif() {
		return tolak("Leader %s is not an active leader", isi.LeaderID)
	}
	m.ClientID2, m.MOLeader = ld.ID, ld.ClientName
	return nil
}

// isiCabang mengisi BRANCHPARENT, BRANCHDETAILID, BRANCHDETAILNAME, dan TEAMGROUP (`SaveMarketingOfficer_Act`:
// TeamGroup = KanwilGroup, BranchDetailName = Name cabang). Cabang yang TIDAK diubah tidak diperiksa ulang.
func (l *Layanan) isiCabang(ctx context.Context, tx *dbTx, m *models.MarketingOfficer, isi models.Isian, lama *models.MarketingOfficer) error {
	if lama != nil && isi.BranchDetailID == lama.BranchDetailID && isi.BranchParent == lama.BranchParent {
		m.BranchParent, m.BranchDetailID, m.BranchDetailName, m.TeamGroup =
			lama.BranchParent, lama.BranchDetailID, lama.BranchDetailName, lama.TeamGroup
		return nil
	}
	if len(isi.BranchParent) > LebarBranchParent {
		return tolak("Branch %s is longer than %d characters", isi.BranchParent, LebarBranchParent)
	}
	if isi.BranchDetailID == "" {
		if isi.BranchParent != "" {
			if _, err := l.gudang.AmbilCabang(ctx, tx, isi.BranchParent); errors.Is(err, repository.ErrTidakAda) {
				return tolak("Branch %s not found", isi.BranchParent)
			} else if err != nil {
				return err
			}
		}
		m.BranchParent, m.BranchDetailID, m.BranchDetailName, m.TeamGroup = isi.BranchParent, "", "", ""
		return nil
	}
	c, err := l.gudang.AmbilCabang(ctx, tx, isi.BranchDetailID)
	if errors.Is(err, repository.ErrTidakAda) {
		return tolak("Sub Branch %s not found", isi.BranchDetailID)
	}
	if err != nil {
		return err
	}
	if !c.Aktif {
		return tolak("Sub Branch %s is not active", isi.BranchDetailID)
	}
	if isi.BranchParent != "" && c.Induk != isi.BranchParent {
		return tolak("Sub Branch %s does not belong to Branch %s", isi.BranchDetailID, isi.BranchParent)
	}
	if len(c.Induk) > LebarBranchParent || len(c.ID) > LebarTeks || len(c.Nama) > LebarTeks || len(c.KanwilGroup) > LebarTeks {
		return tolak("Sub Branch %s does not fit the marketing officer table", isi.BranchDetailID)
	}
	m.BranchParent, m.BranchDetailID, m.BranchDetailName, m.TeamGroup = c.Induk, c.ID, c.Nama, c.KanwilGroup
	return nil
}

// periksaSatuAktif - satu baris AKTIF per Marketing Code dan per akun. `kecuali` = baris yang sedang diubah.
func (l *Layanan) periksaSatuAktif(ctx context.Context, tx *dbTx, m models.MarketingOfficer, kecuali string) error {
	if !m.Aktif() {
		return nil
	}
	if m.ClientID != "" {
		ids, err := l.gudang.AktifDenganClientID(ctx, tx, m.ClientID, kecuali)
		if err != nil {
			return err
		}
		if len(ids) > 0 {
			return fmt.Errorf("%w: this person already has an active marketing officer row (ID %s); deactivate it first",
				ErrSudahAktif, ids[0])
		}
	}
	if m.AksesLogin != "" {
		ids, err := l.gudang.AktifDenganAkses(ctx, tx, m.AksesLogin, kecuali)
		if err != nil {
			return err
		}
		if len(ids) > 0 {
			return fmt.Errorf("%w: Login Account %s is already used by active marketing officer ID %s",
				ErrSudahAktif, m.AksesLogin, ids[0])
		}
	}
	return nil
}

// Tambah menyisipkan baris MO baru.
func (l *Layanan) Tambah(ctx context.Context, p inti.Pelaku, isi models.Isian) (models.MarketingOfficer, error) {
	if err := periksaPelaku(p); err != nil {
		return models.MarketingOfficer{}, err
	}
	isi = rapikan(isi)
	if isi.AksesLogin == "" {
		return models.MarketingOfficer{}, tolak("Login Account is required")
	}
	var hasil models.MarketingOfficer
	err := l.tx(ctx, func(tx *dbTx) error {
		akun, err := l.akunSah(ctx, tx, isi.AksesLogin)
		if err != nil {
			return err
		}
		if strings.TrimSpace(akun.Nama) == "" || len(akun.Nama) > LebarTeks {
			return tolak("the name of Login Account %s must be 1 to %d bytes; fix it in Kelola User first", akun.LoginID, LebarTeks)
		}
		m := models.MarketingOfficer{ClientName: akun.Nama, AksesLogin: akun.LoginID, UserUpdate: p.AkunID,
			MOStatus: status(isi.Aktif)}
		// Marketing Code: milik orang itu bila akunnya sudah pernah dipakai baris MO, selain itu CONTACT_ID akunnya.
		lamaID, ada, err := l.gudang.ClientIDDariAkses(ctx, tx, akun.LoginID)
		if err != nil {
			return err
		}
		m.ClientID = lamaID
		if !ada {
			m.ClientID = akun.ContactID
		}
		if m.ClientID == "" {
			return tolak("Login Account %s has no Contact ID yet (migration 905 has not run)", akun.LoginID)
		}
		if len(m.ClientID) > LebarTeks {
			return tolak("the Marketing Code of Login Account %s is longer than %d bytes", akun.LoginID, LebarTeks)
		}
		if err := l.isiLeader(ctx, tx, &m, isi, nil); err != nil {
			return err
		}
		if err := l.isiCabang(ctx, tx, &m, isi, nil); err != nil {
			return err
		}
		if err := l.periksaSatuAktif(ctx, tx, m, ""); err != nil {
			return err
		}
		id, err := l.gudang.SisipMO(ctx, tx, m)
		if err != nil {
			return err
		}
		hasil, err = l.gudang.AmbilMO(ctx, tx, id)
		return err
	})
	return hasil, err
}

// Ubah menulis ulang baris MO. CLIENTID dan CLIENTNAME tidak pernah berubah.
func (l *Layanan) Ubah(ctx context.Context, p inti.Pelaku, id string, isi models.Isian) (models.MarketingOfficer, error) {
	if err := periksaPelaku(p); err != nil {
		return models.MarketingOfficer{}, err
	}
	isi = rapikan(isi)
	var hasil models.MarketingOfficer
	err := l.tx(ctx, func(tx *dbTx) error {
		lama, err := l.gudang.AmbilMO(ctx, tx, id)
		if errors.Is(err, repository.ErrTidakAda) {
			return ErrTidakAda
		}
		if err != nil {
			return err
		}
		m := lama
		m.UserUpdate, m.MOStatus = p.AkunID, status(isi.Aktif)
		// Akun: sama = tidak diperiksa ulang (AKSES_LOGIN lama Pega boleh tetap); kosong = dilepas; lain = akun sah.
		if isi.AksesLogin != lama.AksesLogin {
			m.AksesLogin = ""
			if isi.AksesLogin != "" {
				akun, err := l.akunSah(ctx, tx, isi.AksesLogin)
				if err != nil {
					return err
				}
				m.AksesLogin = akun.LoginID
			}
		}
		if err := l.isiLeader(ctx, tx, &m, isi, &lama); err != nil {
			return err
		}
		if err := l.isiCabang(ctx, tx, &m, isi, &lama); err != nil {
			return err
		}
		if err := l.periksaSatuAktif(ctx, tx, m, id); err != nil {
			return err
		}
		if err := l.gudang.PerbaruiMO(ctx, tx, m); errors.Is(err, repository.ErrTidakAda) {
			return ErrTidakAda
		} else if err != nil {
			return err
		}
		// Log perubahan (migrasi 760): trigger warisan sudah menyisipkan baris lama; AKSES_LOGIN lamanya diisi di
		// sini. Sebelum 760 dijalankan simpan tetap berhasil - log saja yang tanpa Login Account.
		if err := l.gudang.TandaiLog(ctx, tx, id, lama.AksesLogin); err != nil && !errors.Is(err, repository.ErrLogBelumDimigrasi) {
			return err
		}
		hasil, err = l.gudang.AmbilMO(ctx, tx, id)
		return err
	})
	return hasil, err
}
