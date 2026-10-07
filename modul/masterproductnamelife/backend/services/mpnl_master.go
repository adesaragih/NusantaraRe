package services

// Tujuh pemilih master (paket 2, tiket 04, PARITAS §4).
//
//	medan `Search` (`SearchPolicyHolder.CARI1`) → Enter → `SearchPolicyHolder_act` 1 b236 `·`
//	`CARI1 = @toUpperCase(CARI1)` → grid RD berparam CARI1 → `Choose` → `set*_DT`
//
// R/I Rate (`Choose R/I Rate`) dan `View Rate` membaca tabel `M_RATE_LIFE_SUMMARY` / `RATE_LIFE` baca saja
// sejak K1 keputusan work owner 01-10-2026 (OQ-MPNL-03). View tak terbaca = 503 yang menyebut objeknya.

import (
	"context"
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
)

// GudangMaster - pembaca master (bagian Gudang).
type GudangMaster interface {
	CariMaster(ctx context.Context, jenis models.JenisMaster, kata string, batas int) ([]models.NilaiMaster, error)
	AmbilMaster(ctx context.Context, jenis models.JenisMaster, id string) (models.NilaiMaster, bool, error)
	CariPlan(ctx context.Context, kata string) ([]models.JenisPlan, error)
	AmbilPlan(ctx context.Context, id string) (models.JenisPlan, bool, error)
	// DaftarRate - view `RATE_LIFE` satu IDUSEDBY; `terpotong` = lebih dari repository.BatasRate baris.
	DaftarRate(ctx context.Context, idUsedBy string) ([]models.BarisRate, bool, error)
}

// CariMaster - grid pemilih / autocomplete; kata cari dihurufbesarkan seperti
// `SearchPolicyHolder_act` b236.
//
// `batas` > 0 membatasi baris yang dibaca (autocomplete); 0 = seluruh hasil RD (grid pemilih).
func (l *Layanan) CariMaster(ctx context.Context, p inti.Pelaku, jenis models.JenisMaster, kata string, batas int) ([]models.NilaiMaster, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
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

// JawabanRate - dialog `View Rate`. `Terpotong` benar bila view memuat lebih dari repository.BatasRate
// baris (`BrowseRateLife_RD` b730 `pyMaxRecords` 500: Pega memotong diam-diam).
type JawabanRate struct {
	Daftar    []models.BarisRate `json:"daftar"`
	Total     int                `json:"total"`
	Terpotong bool               `json:"terpotong"`
}

// DaftarRate - tombol `View Rate` b34113 (`SetParamRate` b34310 + `localAction ViewRate` b34354, section
// `ViewRate` RD `BrowseRateLife_RD`), view `RATE_LIFE` baca saja (K1 keputusan work owner 01-10-2026, OQ-MPNL-03).
//
// ⚠️ Penyimpangan sadar: grid `ViewRate.xml` b1024 menyaring `idusedby = ParamID.OUTWARDRATEID`, padahal
// `SetParamRate` b259 hanya mengisi `ParamID.RIRATEID` - dari halaman `InputBusinessLife` milik Retro Life
// yang tidak ada di layar ini - dan nol rule korpus mengisi `ParamID.OUTWARDRATEID`. Di Pega dialog ini
// tidak pernah menampilkan rate baris plan. Di sini disaring `RIRATEID` baris plan, maksud tombol yang
// tampil bila `.RIRATE` terisi (b34113).
func (l *Layanan) DaftarRate(ctx context.Context, p inti.Pelaku, riRateID string) (JawabanRate, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return JawabanRate{}, err
	}
	if strings.TrimSpace(riRateID) == "" {
		return JawabanRate{}, GalatValidasi{Pesan: []string{"riRateId is required"}}
	}
	d, terpotong, err := l.gudang.DaftarRate(ctx, riRateID)
	if err != nil {
		return JawabanRate{}, err
	}
	if d == nil {
		d = []models.BarisRate{}
	}
	return JawabanRate{Daftar: d, Total: len(d), Terpotong: terpotong}, nil
}
