package services

// Jalur baca kelima grid dan master rujukan (paket 1).
//
// Setiap grid anak dibaca dengan KEDUA kunci induknya dari baris induk yang
// sungguh ada - persis parameter RD Pega (`PARITAS-LAYAR-DAN-AKSI.md` §3–§6):
// kontrak <- tahun; reinsurer, business <- (tahun, kontrak); security <-
// (tahun, kontrak, reinsurer). Induk yang tidak ada = 404, bukan daftar kosong.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/repository"
)

// Gudang adalah seluruh sentuhan basis data modul ini.
type Gudang interface {
	DaftarTahun(ctx context.Context) ([]models.TahunTreaty, error)
	AmbilTahun(ctx context.Context, tx *db.Tx, id string) (models.TahunTreaty, error)
	DaftarKontrak(ctx context.Context, tx *db.Tx, tahunID string) ([]models.Kontrak, error)
	AmbilKontrak(ctx context.Context, tx *db.Tx, id string) (models.Kontrak, error)
	DaftarReinsurer(ctx context.Context, tx *db.Tx, tahunID, kontrakID string) ([]models.Reinsurer, error)
	AmbilReinsurer(ctx context.Context, tx *db.Tx, id string) (models.Reinsurer, error)
	DaftarSecurity(ctx context.Context, tx *db.Tx, tahunID, kontrakID, reinsurerID string) ([]models.SecurityReinsurer, error)
	AmbilSecurity(ctx context.Context, tx *db.Tx, id string) (models.SecurityReinsurer, error)
	DaftarBusiness(ctx context.Context, tx *db.Tx, tahunID, kontrakID string) ([]models.Business, error)
	AmbilBusiness(ctx context.Context, tx *db.Tx, id string) (models.Business, error)

	JenisReasuransiLife(ctx context.Context) ([]models.JenisReasuransi, error)
	CariMasterReinsurer(ctx context.Context, kata string) ([]models.MasterReinsurer, error)
	AmbilMasterReinsurer(ctx context.Context, id string) (models.MasterReinsurer, bool, error)
	CariMasterBusiness(ctx context.Context, kata string) ([]models.MasterBusiness, error)
	AmbilMasterBusiness(ctx context.Context, id string) (models.MasterBusiness, bool, error)

	// Penulis tahun treaty (paket 2). `USERID` dari model, `TGLUPDATE` =
	// SYSDATE, ID baru dari sequence - ketiganya di dalam gudang.
	SisipTahun(ctx context.Context, tx *db.Tx, t models.TahunTreaty) (string, error)
	PerbaruiTahun(ctx context.Context, tx *db.Tx, t models.TahunTreaty) error
	// SalinTahunKeAnak menulis salinan tahun (K4) ke anaknya yang berbeda:
	// business `TREATYYEAR`, kontrak `TREATYSTARTDATE`/`TREATYENDDATE`; baris
	// yang berubah mendapat `USERID` t.UserID dan `TGLUPDATE` baru (K6).
	SalinTahunKeAnak(ctx context.Context, tx *db.Tx, t models.TahunTreaty) (int64, error)

	// Penulis kontrak (paket 3).
	SisipKontrak(ctx context.Context, tx *db.Tx, k models.Kontrak) (string, error)
	PerbaruiKontrak(ctx context.Context, tx *db.Tx, k models.Kontrak) error
	// SalinKontrakKeAnak menulis salinan jenis kontrak (K4) ke reinsurer dan
	// business kontrak itu yang berbeda; baris berubah mendapat USERID/TGLUPDATE.
	SalinKontrakKeAnak(ctx context.Context, tx *db.Tx, k models.Kontrak) (int64, error)

	// Penulis reinsurer (paket 4).
	SisipReinsurer(ctx context.Context, tx *db.Tx, r models.Reinsurer) (string, error)
	PerbaruiReinsurer(ctx context.Context, tx *db.Tx, r models.Reinsurer) error
	// TotalSharePerKontrak - total PCTSHARE tiap kontrak (tahunID kosong = semua),
	// kontrak tanpa reinsurer ikut dengan total nil (tiket 11).
	TotalSharePerKontrak(ctx context.Context, tahunID string) ([]models.TotalShareKontrak, error)

	// Penulis security (paket 5).
	SisipSecurity(ctx context.Context, tx *db.Tx, s models.SecurityReinsurer) (string, error)
	PerbaruiSecurity(ctx context.Context, tx *db.Tx, s models.SecurityReinsurer) error

	// Penulis business (paket 6) - juga dipakai salin-semua.
	SisipBusiness(ctx context.Context, tx *db.Tx, b models.Business) (string, error)
	PerbaruiBusiness(ctx context.Context, tx *db.Tx, b models.Business) error

	// Kaskade hapus (paket 7, K2): pencacah dan penghapus memakai predikat
	// yang SAMA; penghapus menghapus anak lebih dulu dan mengembalikan cacahnya.
	DampakHapusKontrak(ctx context.Context, tx *db.Tx, id string) (models.Dampak, error)
	HapusKontrak(ctx context.Context, tx *db.Tx, id string) (models.Dampak, error)
	DampakHapusReinsurer(ctx context.Context, tx *db.Tx, id string) (models.Dampak, error)
	HapusReinsurer(ctx context.Context, tx *db.Tx, id string) (models.Dampak, error)
	HapusSecurity(ctx context.Context, tx *db.Tx, id string) error
	HapusBusiness(ctx context.Context, tx *db.Tx, id string) error
}

// Galat "tidak ada" per entitas (404) - handler tidak mengimpor repository.
var (
	ErrTahunTidakAda     = errors.New("services: treaty year not found")
	ErrKontrakTidakAda   = errors.New("services: treaty contract not found")
	ErrReinsurerTidakAda = errors.New("services: reinsurer not found")
	ErrSecurityTidakAda  = errors.New("services: security reinsurer not found")
	ErrBusinessTidakAda  = errors.New("services: business not found")
	// ErrMasterTidakTerbaca - master rujukan tidak dapat dibaca (503, pesan menyebut objeknya).
	ErrMasterTidakTerbaca = repository.ErrMasterTidakTerbaca
	// ErrParameterWajib - parameter kueri wajib kosong (400).
	ErrParameterWajib = errors.New("services: a required query parameter is empty")
	// ErrRateMenungguPersetujuan - sumber tabel rate (autocomplete `R/I RATE`,
	// section `Rate List`) adalah view atas JSON produk rate; membacanya
	// menuntut persetujuan manusia (OQ-MCRL-13). 503: keadaan server, bukan
	// permintaan yang salah.
	ErrRateMenungguPersetujuan = errors.New("services: the life rate tables are not read yet - reading them " +
		"requires work owner approval (OQ-MCRL-13)")
)

// tidakAda menerjemahkan ErrTidakAda repository menjadi galat entitasnya.
func tidakAda(err, jadi error, id string) error {
	if errors.Is(err, repository.ErrTidakAda) {
		return fmt.Errorf("%w: %s", jadi, id)
	}
	return err
}

// ambilTahun, ambilKontrak, ambilReinsurer - induk yang WAJIB ada.
func (l *Layanan) ambilTahun(ctx context.Context, tx *db.Tx, id string) (models.TahunTreaty, error) {
	t, err := l.gudang.AmbilTahun(ctx, tx, id)
	return t, tidakAda(err, ErrTahunTidakAda, id)
}

func (l *Layanan) ambilKontrak(ctx context.Context, tx *db.Tx, id string) (models.Kontrak, error) {
	k, err := l.gudang.AmbilKontrak(ctx, tx, id)
	return k, tidakAda(err, ErrKontrakTidakAda, id)
}

func (l *Layanan) ambilReinsurer(ctx context.Context, tx *db.Tx, id string) (models.Reinsurer, error) {
	r, err := l.gudang.AmbilReinsurer(ctx, tx, id)
	return r, tidakAda(err, ErrReinsurerTidakAda, id)
}

// DaftarTahun - grid tahun treaty (halaman awal).
func (l *Layanan) DaftarTahun(ctx context.Context, p inti.Pelaku) ([]models.TahunTreaty, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return l.gudang.DaftarTahun(ctx)
}

// JawabanKontrak - grid kontrak satu tahun, beserta tahunnya (kepala panel).
type JawabanKontrak struct {
	Tahun  models.TahunTreaty `json:"tahun"`
	Daftar []models.Kontrak   `json:"daftar"`
}

// DaftarKontrak - tombol `ReinsType` (`InputRetrocessionLife.xml` b11988).
func (l *Layanan) DaftarKontrak(ctx context.Context, p inti.Pelaku, tahunID string) (JawabanKontrak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return JawabanKontrak{}, err
	}
	t, err := l.ambilTahun(ctx, nil, tahunID)
	if err != nil {
		return JawabanKontrak{}, err
	}
	d, err := l.gudang.DaftarKontrak(ctx, nil, t.ID)
	return JawabanKontrak{Tahun: t, Daftar: kosongBukanNil(d)}, err
}

// DaftarBusiness - tombol `Business List` (`InputRetroLimitReinsurers.xml` b12999).
func (l *Layanan) DaftarBusiness(ctx context.Context, p inti.Pelaku, kontrakID string) (JawabanBusiness, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return JawabanBusiness{}, err
	}
	k, err := l.ambilKontrak(ctx, nil, kontrakID)
	if err != nil {
		return JawabanBusiness{}, err
	}
	d, err := l.gudang.DaftarBusiness(ctx, nil, k.IDTreatyYear, k.ID)
	return JawabanBusiness{Kontrak: k, Daftar: kosongBukanNil(d)}, err
}

// JawabanReinsurer - grid reinsurer satu kontrak, beserta kontraknya dan
// `Total Share -->>` (b14339; `CountingPercentShare_Act`). `TotalBukan100` =
// penanda mencolok (tiket 05); penyimpanan tidak pernah diblokir karenanya.
type JawabanReinsurer struct {
	Kontrak       models.Kontrak     `json:"kontrak"`
	Daftar        []models.Reinsurer `json:"daftar"`
	TotalShare    string             `json:"totalShare"`
	TotalBukan100 bool               `json:"totalBukan100"`
}

// DaftarReinsurer - tombol `Reinsurer List` (`InputRetroLimitReinsurers.xml` b14034).
func (l *Layanan) DaftarReinsurer(ctx context.Context, p inti.Pelaku, kontrakID string) (JawabanReinsurer, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return JawabanReinsurer{}, err
	}
	k, err := l.ambilKontrak(ctx, nil, kontrakID)
	if err != nil {
		return JawabanReinsurer{}, err
	}
	d, err := l.gudang.DaftarReinsurer(ctx, nil, k.IDTreatyYear, k.ID)
	if err != nil {
		return JawabanReinsurer{}, err
	}
	total, err := jumlahShare(sharesReinsurer(d))
	if err != nil {
		return JawabanReinsurer{}, err
	}
	return JawabanReinsurer{Kontrak: k, Daftar: kosongBukanNil(d), TotalShare: total.Text('f'),
		TotalBukan100: total.Cmp(seratus) != 0}, nil
}

// JawabanSecurity - grid security satu reinsurer, beserta induknya (kepala
// panel `ID Reinsurer` / `Reinsurer Name` / `PCT Share (%)`) dan eksposur
// efektif tiap baris, DIKUNCI ID baris (bukan urutan) - tiket 06, OQ-MCRL-08.
type JawabanSecurity struct {
	Induk    models.Reinsurer           `json:"induk"`
	Daftar   []models.SecurityReinsurer `json:"daftar"`
	Eksposur map[string]string          `json:"eksposur"`
}

// DaftarSecurity - tombol `Security Reinsurer` (`InputSecurityLifeReinsurers.xml` b12441).
func (l *Layanan) DaftarSecurity(ctx context.Context, p inti.Pelaku, reinsurerID string) (JawabanSecurity, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return JawabanSecurity{}, err
	}
	r, err := l.ambilReinsurer(ctx, nil, reinsurerID)
	if err != nil {
		return JawabanSecurity{}, err
	}
	d, err := l.gudang.DaftarSecurity(ctx, nil, r.TreatyYearID, r.TreatyContractID, r.ID)
	if err != nil {
		return JawabanSecurity{}, err
	}
	eks := map[string]string{}
	for _, s := range d {
		e, err := Eksposur(s.PctShare, r.PctShare)
		if err != nil {
			return JawabanSecurity{}, err
		}
		eks[s.ID] = utils.FormatDecimal(e)
	}
	return JawabanSecurity{Induk: r, Daftar: kosongBukanNil(d), Eksposur: eks}, nil
}

// JawabanBusiness - grid business satu kontrak, beserta kontraknya.
type JawabanBusiness struct {
	Kontrak models.Kontrak    `json:"kontrak"`
	Daftar  []models.Business `json:"daftar"`
}

// JenisReasuransi - dropdown `REINS TYPE`. Master kosong = 503 berkata-kata.
func (l *Layanan) JenisReasuransi(ctx context.Context, p inti.Pelaku) ([]models.JenisReasuransi, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	d, err := l.gudang.JenisReasuransiLife(ctx)
	if err != nil {
		return nil, err
	}
	if len(d) == 0 {
		return nil, fmt.Errorf("%w: %s has no row with FLAG = 1 (life)", ErrMasterTidakTerbaca,
			repository.MasterJenisReasuransi)
	}
	return d, nil
}

// CariMasterReinsurer - autocomplete `REINSURER NAME`.
func (l *Layanan) CariMasterReinsurer(ctx context.Context, p inti.Pelaku, kata string) ([]models.MasterReinsurer, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	d, err := l.gudang.CariMasterReinsurer(ctx, kata)
	return kosongBukanNil(d), err
}

// CariMasterBusiness - autocomplete `BUSINESS NAME`.
func (l *Layanan) CariMasterBusiness(ctx context.Context, p inti.Pelaku, kata string) ([]models.MasterBusiness, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	d, err := l.gudang.CariMasterBusiness(ctx, kata)
	return kosongBukanNil(d), err
}

// CariRingkasanRate - autocomplete `R/I RATE` (`BrowseRateLifeSummary`).
//
// ⛔ Menunggu persetujuan (OQ-MCRL-13): galat terang, bukan daftar kosong
// yang membuat pengguna mengira tidak ada tabel rate.
func (l *Layanan) CariRingkasanRate(_ context.Context, p inti.Pelaku, _ string) ([]struct{}, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return nil, ErrRateMenungguPersetujuan
}

// DaftarRate - section `ViewRate` (`Rate List`) untuk satu `RIRATEID`
// (`ViewRate.xml` b1061 `idusedby = ParamID.RIRATEID`).
func (l *Layanan) DaftarRate(_ context.Context, p inti.Pelaku, idUsedBy string) ([]struct{}, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if strings.TrimSpace(idUsedBy) == "" {
		return nil, fmt.Errorf("%w: idusedby", ErrParameterWajib)
	}
	return nil, ErrRateMenungguPersetujuan
}

// kosongBukanNil - JSON `[]`, bukan `null`, untuk daftar kosong.
func kosongBukanNil[T any](d []T) []T {
	if d == nil {
		return []T{}
	}
	return d
}
