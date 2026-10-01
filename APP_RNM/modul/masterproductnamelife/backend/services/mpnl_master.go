package services

// Tujuh pemilih master (paket 2, tiket 04, PARITAS §4).
//
//	medan `Search` (`SearchPolicyHolder.CARI1`) → Enter → `SearchPolicyHolder_act` 1 b236 `·`
//	`CARI1 = @toUpperCase(CARI1)` → grid RD berparam CARI1 → `Choose` → `set*_DT`
//
// ⛔ R/I Rate menunggu OQ-MPNL-03: 503 berkalimat, bukan daftar kosong.

import (
	"context"
	"errors"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

var (
	// ErrMasterTidakTerbaca - master rujukan tidak dapat dibaca (503, pesan menyebut objeknya).
	ErrMasterTidakTerbaca = repository.ErrMasterTidakTerbaca
	// ErrJenisMasterTidakDikenal - segmen jalur bukan pemilih (404).
	ErrJenisMasterTidakDikenal = repository.ErrJenisMasterTidakDikenal
	// ErrRIRateMenungguPersetujuan - sumber R/I Rate (`BrowseRateLifeSummary`) dan
	// `View Rate` (`BrowseRateLife_RD`) adalah view atas JSON rate; membacanya
	// menuntut persetujuan work owner (OQ-MPNL-03, preseden OQ-MCRL-13). 503:
	// keadaan server, bukan permintaan yang salah.
	ErrRIRateMenungguPersetujuan = errors.New("services: the R/I Rate tables are not read yet - reading the life " +
		"rate views requires work owner approval (OQ-MPNL-03)")
)

// GudangMaster - pembaca master (bagian Gudang).
type GudangMaster interface {
	CariMaster(ctx context.Context, jenis models.JenisMaster, kata string, batas int) ([]models.NilaiMaster, error)
	AmbilMaster(ctx context.Context, jenis models.JenisMaster, id string) (models.NilaiMaster, bool, error)
	CariPlan(ctx context.Context, kata string) ([]models.JenisPlan, error)
	AmbilPlan(ctx context.Context, id string) (models.JenisPlan, bool, error)
}

// CariMaster - grid pemilih / autocomplete; kata cari dihurufbesarkan seperti
// `SearchPolicyHolder_act` b236.
//
// `batas` > 0 membatasi baris yang dibaca (autocomplete); 0 = seluruh hasil RD (grid pemilih).
func (l *Layanan) CariMaster(ctx context.Context, p inti.Pelaku, jenis models.JenisMaster, kata string, batas int) ([]models.NilaiMaster, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if jenis == models.MasterRIRate {
		return nil, ErrRIRateMenungguPersetujuan
	}
	if batas < 0 {
		return nil, GalatValidasi{Pesan: []string{"batas must not be negative"}}
	}
	return l.gudang.CariMaster(ctx, jenis, strings.ToUpper(strings.TrimSpace(kata)), batas)
}

// CariPlan - autocomplete `Plan Name` (`.Plan` b33121, `BrowseProductTypeLife_RD`).
func (l *Layanan) CariPlan(ctx context.Context, p inti.Pelaku, kata string) ([]models.JenisPlan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return l.gudang.CariPlan(ctx, strings.ToUpper(strings.TrimSpace(kata)))
}

// DaftarRate - tombol `View Rate` b34113 (`SetParamRate` + `localAction ViewRate`,
// section `ViewRate` RD `BrowseRateLife_RD`): view atas JSON rate - menunggu
// OQ-MPNL-03, 503 berkalimat.
func (l *Layanan) DaftarRate(ctx context.Context, p inti.Pelaku, riRateID string) ([]models.NilaiMaster, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return nil, ErrRIRateMenungguPersetujuan
}
