package services

// Zone dan 4.2 Construction form coverage FIRE (tiket 52). Rate dari tabel db-link (m_eqs_rate dan kawan-kawan) TIDAK
// diisi otomatis dan Coverage Basis 5 / Layering tidak dibangun (keputusan work owner, tiket 52 butir a/b).
//
// Saringan zona `[terverifikasi]` `DataTransform\SetPropertyCoverageDT_FacIn.xml`: .Coverage 100835 / 100828 / 100817 /
// 100816 -> Zone.CARI1..3 = "X" / "Y" / "Z"; 100829 / 100821 -> Zone.CARI1..6 = "1".."6"; selain itu keenamnya kosong
// (saringan C..H BrowseZones_RD dibuang). Coverage di luar daftar langkah induk DT (`.Coverage == 100829 || ...`)
// mungkin mewarisi Zone.CARI coverage sebelumnya di Pega `[dugaan]` - di sini dianggap kosong.
//
// Rate menurut zona (cabang "4.2 PRGBI" CountPremiRate_ACT): ratecoverage.go.

import (
	"context"
	"errors"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrZonaTanpaDatabase - ZONES / CONSTRUCTION / ZONEOJK tidak terbaca. 503.
var ErrZonaTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, zona / konstruksi tidak terbaca")

// zonaMenurutCoverage - SetPropertyCoverageDT_FacIn langkah zona: `rsmd` X/Y/Z (coverage RSMD, `[dugaan]` dari daftar
// kode tiket 52), `gempa` 1-6.
var zonaMenurutCoverage = func() map[string][]string {
	rsmd, gempa := []string{"X", "Y", "Z"}, []string{"1", "2", "3", "4", "5", "6"}
	return map[string][]string{"100835": rsmd, "100828": rsmd, "100817": rsmd, "100816": rsmd, "100829": gempa, "100821": gempa}
}()

// DenganZona memasang pembaca ZONES / CONSTRUCTION / ZONEOJK (tiket 52).
func (s *Service) DenganZona(z repository.PembacaZona) *Service { s.zona = z; return s }

// SaranZona - GET /api/nbfacin/zona: BrowseZones_RD dengan OLDID = oldId dan zona menurut kode coverage.
func (s *Service) SaranZona(ctx context.Context, coverage, oldID, kata string) ([]models.SaranAkumulasi, error) {
	if err := periksaParamKueri(ErrMasukanCoverage, paramAkumulasi{"coverage", &coverage}, paramAkumulasi{"oldId", &oldID},
		paramAkumulasi{"q", &kata}); err != nil {
		return nil, err
	}
	if s.zona == nil {
		return nil, ErrZonaTanpaDatabase
	}
	return s.zona.SaranZona(ctx, oldID, zonaMenurutCoverage[coverage], kata)
}

// SaranKonstruksiGempa - GET /api/nbfacin/konstruksi-gempa: SearchClassConstructionReas_SQL.
func (s *Service) SaranKonstruksiGempa(ctx context.Context, kata string) ([]models.SaranAkumulasi, error) {
	if err := periksaParamKueri(ErrMasukanCoverage, paramAkumulasi{"q", &kata}); err != nil {
		return nil, err
	}
	if s.zona == nil {
		return nil, ErrZonaTanpaDatabase
	}
	return s.zona.SaranKonstruksiGempa(ctx, kata)
}
