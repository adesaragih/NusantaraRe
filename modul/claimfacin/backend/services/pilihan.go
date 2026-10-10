package services

// Untuk apa berkas ini: PILIHAN POP-UP DAN AUTOCOMPLETE - isi grid ViewPolis (SearchPolis_act), CauseofLoss_Harness,
// CatastrofeList, pop-up Outstanding (GetAllData_Act), dan sumber sel autocomplete (adjuster, wilayah, item / coverage /
// aneka / okupasi dari halaman polis). Baca-saja; pilihan yang dipilih dikirim balik sebagai aksi dan DIBACA ULANG di
// server.

import (
	"context"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/claimfacin/backend/models"
)

// Jenis pilihan pop-up.
const (
	PilihanPolis       = "polis"       // ViewPolis (SearchPolis_act + GetPolisForClaim_SQL)
	PilihanSebab       = "sebab"       // CauseofLoss_Harness
	PilihanKatastrofe  = "katastrofe"  // CatastrofeList
	PilihanOutstanding = "outstanding" // harness Outstanding (GetAllData_Act), n = objek
)

// PermintaanPilihan - satu permintaan daftar pilihan.
type PermintaanPilihan struct {
	Jenis string
	// Konteks - panel baris tempat sel berada (sumber tingkat item); "" = layar utama.
	Konteks string
	// N - baris grid di konteks (item untuk sumber tingkat item; objek untuk Outstanding).
	N int
	// Cari - teks saringan (autocomplete / kotak cari pop-up).
	Cari string
	// JenisCari - Search Type pop-up ViewPolis.
	JenisCari string
}

// Pilihan membaca satu daftar pilihan kasus.
func (l *Layanan) Pilihan(ctx context.Context, p inti.Pelaku, id string, r PermintaanPilihan) (any, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	switch r.Jenis { // pilihan tanpa halaman kasus
	case PilihanPolis:
		return l.a.CariPolis(ctx, r.JenisCari, r.Cari)
	case PilihanSebab:
		return l.a.DaftarSebab(ctx, r.Cari)
	case PilihanKatastrofe:
		return l.a.DaftarKatastrofe(ctx, r.Cari)
	case models.SumberAdjuster:
		return l.a.DaftarAdjuster(ctx, r.Cari)
	case models.SumberNegara:
		return l.a.DaftarNegara(ctx, r.Cari)
	case models.SumberMataUang:
		return l.a.DaftarMataUang(ctx)
	case models.SumberJenisReas:
		return l.a.DaftarJenisReas(ctx)
	}
	k, h, err := l.muat(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	switch r.Jenis {
	case models.SumberProvinsi:
		return l.a.DaftarProvinsi(ctx, h.Ambil(models.CD+"Country"), r.Cari)
	case models.SumberKota:
		return l.a.DaftarKota(ctx, h.Ambil(models.CD+"Province"), r.Cari)
	case models.SumberDistrik:
		return l.a.DaftarDistrik(ctx, h.Ambil(models.CD+"CityID"), r.Cari)
	case models.SumberRW:
		return l.a.DaftarRW(ctx, h.Ambil(models.CD+"DistrictID"), r.Cari)
	case PilihanOutstanding:
		if r.N < 1 || r.N > len(h.AmbilDaftar(models.DaftarObjek)) {
			return nil, fmt.Errorf("%w: objek %d", ErrPermintaanTidakSah, r.N)
		}
		rows, err := l.g.BacaOS(ctx, k.ID, h.Ambil(models.CD+"NoClaim"))
		if err != nil {
			return nil, err
		}
		kt, err := l.konteks(ctx, p, k, l.jam())
		if err != nil {
			return nil, err
		}
		return models.RingkasOutstanding(kt, h, rows)
	case models.SumberMataUangAdj, models.SumberRekening:
		models.HitungTurunan(h)
		return l.pilihanAdj(ctx, h, r)
	}
	models.HitungTurunan(h)
	it, err := itemKonteks(h, r)
	if err != nil {
		return nil, err
	}
	var out []models.Pilihan
	saring := func(nilai, label string, tambahan map[string]string) {
		c := strings.ToUpper(strings.TrimSpace(r.Cari))
		if c == "" || strings.Contains(strings.ToUpper(label), c) || strings.Contains(strings.ToUpper(nilai), c) {
			out = append(out, models.Pilihan{Nilai: nilai, Label: label, Tambahan: tambahan})
		}
	}
	switch r.Jenis {
	case models.SumberItemProp: // D_FilteredPropertyItemList: nilai .ItemType, isi .IndexPropertyItem
		for _, b := range models.OpsiPropertyItem(h, it["IndexObject"]) {
			saring(b["IndexPropertyItem"], b["ItemType"], nil)
		}
	case models.SumberCovFire: // D_FilteredCoverageList: nilai .OLDID
		for _, b := range models.OpsiCoverageFire(h, it["IndexObject"], it["ObjectItemID"]) {
			saring(b["OLDID"], b["OLDID"], map[string]string{"CoverageNote": b["CoverageNote"]})
		}
	case models.SumberAneka: // D_AnekaList: nilai .IdxAneka
		for _, b := range models.OpsiAneka(h, it["IndexObject"]) {
			saring(b["IdxAneka"], b["IdxAneka"], map[string]string{"ObjectName": b["ObjectName"]})
		}
	case models.SumberCovAneka: // D_FilteredCoverageAnekaList: nilai .CoverageNote
		for _, b := range models.OpsiCoverageAneka(h, it) {
			saring(b["CoverageNote"], b["CoverageNote"], nil)
		}
	case models.SumberCovObjek: // D_Coverage{MBU,PA,Travel}ClaimList: nilai .Coverage
		for _, b := range models.OpsiCoverageObjek(h, it) {
			saring(b["Coverage"], b["Coverage"], map[string]string{"CoverageNote": b["CoverageNote"]})
		}
	case models.SumberOkupasi: // D_OccupationList: nilai .OccupationName, isi .OccupationId
		for _, b := range models.OpsiOkupasi(h, it["IndexObject"]) {
			saring(b["OccupationId"], b["OccupationName"], nil)
		}
	default:
		return nil, fmt.Errorf("%w: pilihan %q", ErrPermintaanTidakSah, r.Jenis)
	}
	if out == nil {
		out = []models.Pilihan{}
	}
	return out, nil
}

// pilihanAdj - sumber sel panel adjustment ("adjdtl:...(o)...(i)...(a)"): Choose Currency (`.CurencyAdjustment`
// baris) dan Name of Bank (`Result.pxResults` SetPayable_Act 6 / 8.2 - nilai "AccountNo|NameOfBank").
func (l *Layanan) pilihanAdj(ctx context.Context, h *models.Halaman, r PermintaanPilihan) ([]models.Pilihan, error) {
	prefiks, idx, ok := models.UraiPanel(r.Konteks)
	if !ok || prefiks != models.PanelAdj || len(idx) < 3 {
		return nil, fmt.Errorf("%w: konteks pilihan %q", ErrPermintaanTidakSah, r.Konteks)
	}
	o, i, a := idx[0], idx[1], idx[2]
	if _, err := models.Adj(h, o, i, a); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPermintaanTidakSah, err)
	}
	out := []models.Pilihan{}
	if r.Jenis == models.SumberMataUangAdj {
		for _, m := range h.AmbilDaftar(models.DaftarDiAdj(o, i, a, models.DaftarMataUangAdj)) {
			out = append(out, models.Pilihan{Nilai: m["CurrencyID"], Label: m["Currency"]})
		}
		return out, nil
	}
	rek, err := l.rekeningAdj(ctx, h, o, i, a)
	if err != nil {
		return nil, err
	}
	c := strings.ToUpper(strings.TrimSpace(r.Cari))
	for _, x := range rek {
		if c != "" && !strings.Contains(strings.ToUpper(x.NameOfBank), c) {
			continue
		}
		out = append(out, models.Pilihan{Nilai: models.KunciRekening(x), Label: x.NameOfBank,
			Tambahan: map[string]string{"BranchOfBank": x.BranchOfBank, "NoAccount": x.AccountNo,
				"SwiftCode": x.SwiftCode, "IDOfBank": x.IDOfBank}})
	}
	return out, nil
}

// rekeningAdj - SetPayable_Act 6 / 8.2: rekening klien payable + mata uang adjustment (Payable 3: mata uang saja).
// Halaman dipakai salinan - RencanaPayable menulis payable baris.
func (l *Layanan) rekeningAdj(ctx context.Context, h *models.Halaman, o, i, a int) ([]models.RekeningBank, error) {
	r, err := models.RencanaPayable(h.Salin(), o, i, a)
	if err != nil {
		return nil, err
	}
	if r.SemuaKlien {
		return l.a.RekeningBankMataUang(ctx, r.MataUang)
	}
	if r.Klien == "" {
		return nil, nil
	}
	return l.a.RekeningBank(ctx, r.Klien, r.MataUang)
}

// itemKonteks - baris item sumber autocomplete tingkat item: konteks panel objek "est:...(o)" + N = item.
func itemKonteks(h *models.Halaman, r PermintaanPilihan) (models.Baris, error) {
	prefiks, idx, ok := models.UraiPanel(r.Konteks)
	if !ok || prefiks != models.PanelObjekEst {
		return nil, fmt.Errorf("%w: konteks pilihan %q", ErrPermintaanTidakSah, r.Konteks)
	}
	it, err := models.Item(h, idx[0], r.N)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPermintaanTidakSah, err)
	}
	return it, nil
}
