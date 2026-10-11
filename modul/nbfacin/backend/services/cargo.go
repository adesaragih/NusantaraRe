package services

// Tab Object + Coverage MARINE CARGO (tiket 74 P1 + P2, keputusan work owner 10-10-2026 K74-1…K74-6). Bukti `[terverifikasi]`
// (`NB\` = `D:\migrasi\RNM\NB FacIn\`):
//   - gerbang When IsMarineCargo (`NB\When\IsMarineCargo.xml` L44);
//   - medan server objek = ProtectObjectMarine (`NB\Activity\ProtectObjectMarine.xml`): Obj-Browse master per ID -> `.Description`
//     ke TradingNote (L783), ConveyanceName + ConveyanceNote (L1470 / L1516), GoodNote (L2158), PackingNote (L2826); MaxAgeDetail =
//     `.Warranty` BrowseMaxAge_RD (MarineCargoDtl L14886). Nama tabel master `[pertanyaan terbuka]` (Q-DBA, cek DEV oleh work owner)
//     - selama PembacaMasterCargo belum terpasang, catatan lama dipertahankan bila ID tidak berubah; ID baru -> 503;
//   - CountAgeShip (`NB\Activity\CountAgeShip.xml` L352-439): AGE = @toDecimal(@substring(SailDate,0,4)) - @toDecimal(Y_MAKE1), hanya
//     bila Y_MAKE1 terisi (L531) - selain itu AGE tidak diubah. Pega menjalankannya saat Sail Date berubah; di sini setiap PUT
//     `[keputusan agent]` (hasil sama bila kapal tidak berubah);
//   - SelectShip_act L664 / L705 / L726: GRT / NRT / DWT `@replaceAll(…,',','.')` (lalu @toDecimal - teks tak terurai `belum
//     terverifikasi`, jadi disimpan hasil replaceAll apa adanya `[keputusan agent]`);
//   - CountCoverageMarine L998-1073: Premium = @Math.divide(TSI×Rate,100,4) (premium.Calculate lini MARINE CARGO), Discount =
//     @Math.divide(Premium×DiscountPercentage,100,4), TSILiability = TSI; CountTSIandRateMarineCargo_Act L245: Rate =
//     @Math.divide(Premium×100,TSI,4); SetCurrencyCoverage_Act L744 / L791: Currency.ID / OldID dari master CURRENCY;
//   - AppendCoverageMarineCargo_ACT L249 / L383 (prakondisi CoverageList kosong L462): ProRatePercent = min(100,
//     @Math.divide(hari(Start..End),365,8)×100) - disintesis di GET (TIDAK disimpan sampai PUT);
//   - CountPremiAndTSIRNMCargo_ACT langkah 7: TSINusantaraRe = PS×TSI/100; PremiNusantaraRe = PS×Premium/100 − PS×Discount/100
//     (skala 20), disimpan 8 desimal;
//   - SumTotalTSIPremiGross_Act langkah 4 (MARINE, prakondisi IsMarineCargo L7297): per cargo, per `.Currency.Name` ΣTSI, ΣPremium,
//     Rate = @if(TSI=0,0,@divide(Premium,TSI,20)×100) (L6890 / L7154).
// K74-3: MOP / Master / Declaration tidak diport (Premium=0 MOP, CurrencyMaster, PolicyMaster). Wajib isi (`pyRequiredNew always`,
// InputDtlCargo_FacIn): Trading, Goods, Conveyance, Packing, From Rute, To Rute.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services/premium"
)

var (
	// ErrBukanMarineCargo - tab Marine Cargo dipakai case yang bukan MARINE CARGO (409).
	ErrBukanMarineCargo = errors.New("services: case bukan MARINE CARGO (When IsMarineCargo)")
	// ErrMasterCargoBelumAda - master Trading / Goods / Conveyance / Packing / MaxAge belum terpasang (nama tabel menunggu
	// cek DEV, Q-DBA tiket 74) dan catatan tidak dapat dipertahankan (503).
	ErrMasterCargoBelumAda = errors.New("services: master Marine Cargo belum dipetakan (nama tabel Oracle menunggu cek DEV)")
)

// Jenis master ProtectObjectMarine (Param.Data).
const (
	MasterTrading    = "Trading"
	MasterGoods      = "Goods"
	MasterConveyance = "Conveyance"
	MasterPacking    = "Packing"
)

// PembacaMasterCargo - Obj-Browse master ProtectObjectMarine / BrowseMaxAge_RD (ada = false bila ID tidak ada).
type PembacaMasterCargo interface {
	Deskripsi(ctx context.Context, jenis, id string) (string, bool, error)
	WarrantyMaxAge(ctx context.Context, id string) (string, bool, error)
	// IDMaster - unggah CSV (tiket 74 P6): ID master `jenis` (Trading / Goods / Packing / Conveyance) ber-Description PERSIS
	// `deskripsi` (BrowseTrading_RD / BrowseGoodsType / BrowsePacking_RD / BrowseConveyance_RD `.Description =`, peka huruf).
	IDMaster(ctx context.Context, jenis, deskripsi string) (string, bool, error)
	// Agen - BrowseLLOYDAGENT_RD `.ID =` (ID, Name, Address, City, Country -> ASMAgentID / pyFullName / pyAddress / pyCity / pyCountry).
	Agen(ctx context.Context, id string) (models.AgenSurvei, bool, error)
	// Daftar - dropdown form cargo master `jenis` (Trading / Goods / Conveyance / Packing).
	Daftar(ctx context.Context, jenis string) ([]models.BarisMasterCargo, error)
	// CariAgen / CariMaxAge / CariKondisiMarine - popup Select Agent, autocomplete Max Age, Select Coverage / Warranty Extension
	// (repository/mastercargo.go); kata kosong = tanpa saringan.
	CariAgen(ctx context.Context, kata string) ([]models.BarisAgenSurvei, error)
	CariMaxAge(ctx context.Context, kata string) ([]models.BarisMaxAge, error)
	CariKondisiMarine(ctx context.Context, kata string) ([]models.BarisKondisiMarine, error)
}

// CariAgenSurvei - GET /api/nbfacin/marine/agent?q= (master belum terpasang -> 503).
func (s *Service) CariAgenSurvei(ctx context.Context, q string) ([]models.BarisAgenSurvei, error) {
	if s.masterCargo == nil {
		return nil, fmt.Errorf("%w: agent", ErrMasterCargoBelumAda)
	}
	return s.masterCargo.CariAgen(ctx, q)
}

// CariMaxAge - GET /api/nbfacin/marine/max-age?q=.
func (s *Service) CariMaxAge(ctx context.Context, q string) ([]models.BarisMaxAge, error) {
	if s.masterCargo == nil {
		return nil, fmt.Errorf("%w: max-age", ErrMasterCargoBelumAda)
	}
	return s.masterCargo.CariMaxAge(ctx, q)
}

// CariKondisiMarine - GET /api/nbfacin/marine/coverage?q=.
func (s *Service) CariKondisiMarine(ctx context.Context, q string) ([]models.BarisKondisiMarine, error) {
	if s.masterCargo == nil {
		return nil, fmt.Errorf("%w: coverage", ErrMasterCargoBelumAda)
	}
	return s.masterCargo.CariKondisiMarine(ctx, q)
}

// jenisMasterKabel - segmen GET /api/nbfacin/marine/{jenis} -> jenis master.
var jenisMasterKabel = map[string]string{"trading": MasterTrading, "goods": MasterGoods, "conveyance": MasterConveyance,
	"packing": MasterPacking}

// DaftarMasterCargo - GET /api/nbfacin/marine/{trading|goods|conveyance|packing} (dropdown form cargo). Jenis lain (coverage /
// max-age / agent - DDL belum ada) -> ErrMasterCargoBelumAda (503); master belum terpasang -> 503.
func (s *Service) DaftarMasterCargo(ctx context.Context, jenisKabel string) ([]models.BarisMasterCargo, error) {
	jenis, ok := jenisMasterKabel[jenisKabel]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMasterCargoBelumAda, jenisKabel)
	}
	if s.masterCargo == nil {
		return nil, fmt.Errorf("%w: %s", ErrMasterCargoBelumAda, jenis)
	}
	return s.masterCargo.Daftar(ctx, jenis)
}

// DenganObjekCargo memasang penyimpan tab Marine Cargo (tiket 74).
func (s *Service) DenganObjekCargo(o repository.PenyimpanObjekCargo) *Service {
	s.objekCargo = o
	return s
}

// DenganMasterCargo memasang pembaca master Marine Cargo (tiket 74; menunggu nama tabel DEV).
func (s *Service) DenganMasterCargo(m PembacaMasterCargo) *Service { s.masterCargo = m; return s }

// CargoTampil - satu cargo GET / PUT dengan total gross per mata uang (baca saja).
type CargoTampil struct {
	Objek      models.ObjekCargo
	TotalGross []models.TotalGrossCargo
}

// BarisHitungCargo - satu coverage hitung-coverage-cargo dan sel pemicunya.
type BarisHitungCargo struct {
	Coverage models.CoverageCargo
	Pemicu   string
}

// Pemicu hitung-coverage-cargo (sel on-change InputCoverageCargo_FacIn).
const (
	PemicuCargoCurrency     = "currency"     // L2787 SetCurrencyCoverage_Act + CountCoverageMarine
	PemicuCargoRate         = "rate"         // CountCoverageMarine
	PemicuCargoTSI          = "tsi"          // CountCoverageMarine
	PemicuCargoDiskonPersen = "diskonPersen" // CountCoverageMarine
	PemicuCargoPremium      = "premium"      // L6083 CountTSIandRateMarineCargo_Act
)

var pemicuCargo = map[string]bool{PemicuCargoCurrency: true, PemicuCargoRate: true, PemicuCargoTSI: true,
	PemicuCargoDiskonPersen: true, PemicuCargoPremium: true}

// lebarTeksCargo - lebar kolom migrasi 256 (byte) medan teks yang ditulis.
var lebarTeksCargo = []struct {
	jalur string
	nilai func(*models.ObjekCargo) string
	n     int
}{
	{"tradingId", func(c *models.ObjekCargo) string { return c.TradingID }, 50},
	{"goodId", func(c *models.ObjekCargo) string { return c.GoodID }, 50},
	{"conveyanceId", func(c *models.ObjekCargo) string { return c.ConveyanceID }, 50},
	{"packingId", func(c *models.ObjekCargo) string { return c.PackingID }, 50},
	{"fromRute", func(c *models.ObjekCargo) string { return c.FromRute }, 50},
	{"toRute", func(c *models.ObjekCargo) string { return c.ToRute }, 50},
	{"blNumber", func(c *models.ObjekCargo) string { return c.BLNumber }, 4000},
	{"invoiceNumber", func(c *models.ObjekCargo) string { return c.InvoiceNumber }, 50},
	{"maxAge", func(c *models.ObjekCargo) string { return c.MaxAge }, 50},
	{"surveyAgent.agentId", func(c *models.ObjekCargo) string { return c.SurveyAgent.AgentID }, 100},
	{"surveyAgent.fullName", func(c *models.ObjekCargo) string { return c.SurveyAgent.FullName }, 4000},
	{"surveyAgent.address", func(c *models.ObjekCargo) string { return c.SurveyAgent.Address }, 500},
	{"surveyAgent.city", func(c *models.ObjekCargo) string { return c.SurveyAgent.City }, 500},
	{"surveyAgent.country", func(c *models.ObjekCargo) string { return c.SurveyAgent.Country }, 500},
	{"lc.bank", func(c *models.ObjekCargo) string { return c.LC.Bank }, 500},
	{"lc.number", func(c *models.ObjekCargo) string { return c.LC.Number }, 500},
	{"lc.condition", func(c *models.ObjekCargo) string { return c.LC.Condition }, 500},
	{"lc.remark", func(c *models.ObjekCargo) string { return c.LC.Remark }, 500},
	{"ship.id", func(c *models.ObjekCargo) string { return c.Ship.ID }, 100},
	{"ship.nmShip", func(c *models.ObjekCargo) string { return c.Ship.NmShip }, 1000},
	{"ship.shipTypeId", func(c *models.ObjekCargo) string { return c.Ship.ShipTypeID }, 1000},
	{"ship.shipTypeName", func(c *models.ObjekCargo) string { return c.Ship.ShipTypeName }, 1000},
	{"ship.grt", func(c *models.ObjekCargo) string { return c.Ship.GRT }, 1000},
	{"ship.nrt", func(c *models.ObjekCargo) string { return c.Ship.NRT }, 1000},
	{"ship.dwt", func(c *models.ObjekCargo) string { return c.Ship.DWT }, 1000},
	{"ship.flagName", func(c *models.ObjekCargo) string { return c.Ship.FlagName }, 1000},
	{"ship.construction", func(c *models.ObjekCargo) string { return c.Ship.Construction }, 1000},
	{"ship.yMake1", func(c *models.ObjekCargo) string { return c.Ship.YMake1 }, 1000},
	{"ship.yMake2", func(c *models.ObjekCargo) string { return c.Ship.YMake2 }, 1000},
	{"ship.remark", func(c *models.ObjekCargo) string { return c.Ship.Remark }, 1000},
	{"ship.register", func(c *models.ObjekCargo) string { return c.Ship.Register }, 1000},
	{"ship.imo", func(c *models.ObjekCargo) string { return c.Ship.Imo }, 1000},
	{"ship.formerShipName", func(c *models.ObjekCargo) string { return c.Ship.FormerShipName }, 1000},
	{"ship.classification", func(c *models.ObjekCargo) string { return c.Ship.Classification }, 1000},
	{"ship.rebuilt", func(c *models.ObjekCargo) string { return c.Ship.Rebuilt }, 1000},
}

// lebarKapalTambahan - lebar T_ADDITIONALSHIP.
var lebarKapalTambahan = []struct {
	jalur string
	nilai func(models.KapalTambahan) string
	n     int
}{
	{"id", func(k models.KapalTambahan) string { return k.ID }, 100},
	{"voyageNo", func(k models.KapalTambahan) string { return k.VoyageNo }, 500},
	{"register", func(k models.KapalTambahan) string { return k.Register }, 1000},
	{"nmShip", func(k models.KapalTambahan) string { return k.NmShip }, 1000},
	{"shipTypeId", func(k models.KapalTambahan) string { return k.ShipTypeID }, 1000},
	{"flagName", func(k models.KapalTambahan) string { return k.FlagName }, 1000},
	{"nrt", func(k models.KapalTambahan) string { return k.NRT }, 1000},
	{"dwt", func(k models.KapalTambahan) string { return k.DWT }, 1000},
	{"grt", func(k models.KapalTambahan) string { return k.GRT }, 1000},
	{"yMake1", func(k models.KapalTambahan) string { return k.YMake1 }, 1000},
	{"construction", func(k models.KapalTambahan) string { return k.Construction }, 1000},
}

// tanggalCargo - tiga medan tanggal (kabel DD-MM-YYYY, disimpan teks Pega DateTime pola Date of Loss, sama dengan Aneka).
var tanggalCargo = []struct {
	jalur string
	ke    func(*models.ObjekCargo) *string
}{
	{"sailDate", func(c *models.ObjekCargo) *string { return &c.SailDate }},
	{"invoiceDate", func(c *models.ObjekCargo) *string { return &c.InvoiceDate }},
	{"lc.endDate", func(c *models.ObjekCargo) *string { return &c.LC.EndDate }},
}

// desimalCoverageCargo - medan angka coverage yang diisi layar (dibulatkan 8 desimal saat simpan, wajib muat NUMBER(38,8)).
var desimalCoverageCargo = []struct {
	jalur string
	ke    func(*models.CoverageCargo) **apd.Decimal
}{
	{"rate", func(c *models.CoverageCargo) **apd.Decimal { return &c.Rate }},
	{"limitOfLiability", func(c *models.CoverageCargo) **apd.Decimal { return &c.LimitOfLiability }},
	{"discountPercentage", func(c *models.CoverageCargo) **apd.Decimal { return &c.DiscountPercentage }},
	{"tsi", func(c *models.CoverageCargo) **apd.Decimal { return &c.TSI }},
	{"premium", func(c *models.CoverageCargo) **apd.Decimal { return &c.Premium }},
	{"minPremium", func(c *models.CoverageCargo) **apd.Decimal { return &c.MinPremium }},
}

// wajibIsiCargo - medan `pyRequiredNew always` form cargo.
var wajibIsiCargo = []struct {
	jalur string
	nilai func(*models.ObjekCargo) string
}{
	{"tradingId", func(c *models.ObjekCargo) string { return c.TradingID }},
	{"goodId", func(c *models.ObjekCargo) string { return c.GoodID }},
	{"conveyanceId", func(c *models.ObjekCargo) string { return c.ConveyanceID }},
	{"packingId", func(c *models.ObjekCargo) string { return c.PackingID }},
	{"fromRute", func(c *models.ObjekCargo) string { return c.FromRute }},
	{"toRute", func(c *models.ObjekCargo) string { return c.ToRute }},
}

// periksaCargo - batas baris, lebar teks, tanggal DD-MM-YYYY, wajib isi, angka coverage muat NUMBER(38,8) sesudah 8 desimal.
func periksaCargo(cargo []models.ObjekCargo) error {
	if len(cargo) > batasBarisObjek {
		return fmt.Errorf("%w: paling banyak %d cargo", ErrMasukanObjek, batasBarisObjek)
	}
	var masalah []string
	for i := range cargo {
		c := &cargo[i]
		j := fmt.Sprintf("cargo[%d]", i)
		for _, w := range wajibIsiCargo {
			if strings.TrimSpace(w.nilai(c)) == "" {
				masalah = append(masalah, j+"."+w.jalur+" wajib diisi")
			}
		}
		for _, x := range lebarTeksCargo {
			if len(x.nilai(c)) > x.n {
				masalah = append(masalah, fmt.Sprintf("%s.%s paling banyak %d byte", j, x.jalur, x.n))
			}
		}
		for _, d := range tanggalCargo {
			uraiKabel(j+"."+d.jalur, *d.ke(c), &masalah)
		}
		if len(c.Ship.AdditionalShip) > batasBarisObjek || len(c.Coverages) > batasBarisObjek {
			masalah = append(masalah, fmt.Sprintf("%s: paling banyak %d baris", j, batasBarisObjek))
		}
		for k, x := range c.Ship.AdditionalShip {
			for _, l := range lebarKapalTambahan {
				if len(l.nilai(x)) > l.n {
					masalah = append(masalah, fmt.Sprintf("%s.ship.additionalShip[%d].%s paling banyak %d byte", j, k, l.jalur, l.n))
				}
			}
		}
		for k := range c.Coverages {
			v := &c.Coverages[k]
			jv := fmt.Sprintf("%s.coverages[%d]", j, k)
			for _, x := range []struct {
				jalur, nilai string
				n            int
			}{{"coverage", v.Coverage, 50}, {"coverageNote", v.CoverageNote, 4000}, {"coverageInitial", v.CoverageInitial, 500}, // 251
				{"currency", v.Currency, 10}} {
				if len(x.nilai) > x.n {
					masalah = append(masalah, fmt.Sprintf("%s.%s paling banyak %d byte", jv, x.jalur, x.n))
				}
			}
			for _, d := range desimalCoverageCargo {
				if x := simpan8(*d.ke(v)); x != nil && !repository.PolaAngkaMasuk.MatchString(utils.FormatDecimal(x)) {
					masalah = append(masalah, fmt.Sprintf("%s.%s melebihi 30 digit bulat", jv, d.jalur))
				}
			}
			for m, w := range v.WarrantyExtension {
				if len(w.Coverage) > 50 || len(w.CoverageNote) > 500 {
					masalah = append(masalah, fmt.Sprintf("%s.warrantyExtension[%d] coverage / coverageNote paling banyak 50 / 500 byte", jv, m))
				}
			}
			if len(v.Deductibles) > batasBarisObjek || len(v.WarrantyExtension) > batasBarisObjek {
				masalah = append(masalah, fmt.Sprintf("%s: paling banyak %d baris", jv, batasBarisObjek))
			}
			for m, de := range v.Deductibles {
				if len(de.Descriptions) > 4000 {
					masalah = append(masalah, fmt.Sprintf("%s.deductibles[%d].descriptions paling banyak 4000 byte", jv, m))
				}
			}
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	return nil
}

// adalahMarineCargo - When IsMarineCargo case `id` (BUSINESS tidak terbaca -> salah).
func (s *Service) adalahMarineCargo(ctx context.Context, id string) bool {
	kb, ada := s.kasusBisnisDari(ctx, id)
	return ada && salahSatu(id, kb, "IsMarineCargo")
}

// BacaObjekCargo - GET objek-cargo: tanggal ke kabel, coverage kosong disintesis (AppendCoverageMarineCargo_ACT, tidak disimpan),
// Currency ID / OldID dari master, total gross per cargo.
func (s *Service) BacaObjekCargo(ctx context.Context, id string) ([]CargoTampil, error) {
	if s.objekCargo == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	cargo, err := s.objekCargo.BacaObjekCargo(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	var proRate *apd.Decimal
	for i := range cargo {
		if len(cargo[i].Coverages) == 0 {
			if proRate == nil {
				proRate = s.proRateCargo(ctx, id)
			}
			cargo[i].Coverages = []models.CoverageCargo{coverageCargoKosong(proRate)}
		}
	}
	return s.tampilanCargo(ctx, cargo)
}

// coverageCargoKosong - AppendCoverageMarineCargo_ACT: satu coverage dengan ProRatePercent.
func coverageCargoKosong(proRate *apd.Decimal) models.CoverageCargo {
	return models.CoverageCargo{ProRatePercent: proRate, WarrantyExtension: []models.WarrantyExtension{},
		Deductibles: []models.DeductibleCargo{}, Spreading: []models.BarisSpreading{}}
}

// proRateCargo - AppendCoverageMarineCargo_ACT L249 / L383: min(100, @Math.divide(hari(Start..End),365,8)×100); periode tidak
// terbaca -> nil.
func (s *Service) proRateCargo(ctx context.Context, id string) *apd.Decimal {
	periode, err := s.periodeKasus(ctx, id)
	if err != nil {
		return nil
	}
	mulai, ok1 := tanggalJakarta(periode.Mulai)
	akhir, ok2 := tanggalJakarta(periode.Akhir)
	if !ok1 || !ok2 {
		return nil
	}
	h, err := bagi(apd.New(hari(mulai, akhir), 0), tigaRatusLima, 8)
	if err != nil {
		return nil
	}
	pr := kali(h, seratus)
	if pr.Cmp(seratus) > 0 {
		pr = apd.New(100, 0)
	}
	return simpan8(pr)
}

// tampilanCargo - tanggal ke kabel, mata uang ID / OldID, total gross.
func (s *Service) tampilanCargo(ctx context.Context, cargo []models.ObjekCargo) ([]CargoTampil, error) {
	hasil := make([]CargoTampil, 0, len(cargo))
	mu := map[string]models.MataUangCsv{}
	for _, c := range cargo {
		for _, d := range tanggalCargo {
			*d.ke(&c) = tanggalKerugianKeKabel(*d.ke(&c))
		}
		for k := range c.Coverages {
			v := &c.Coverages[k]
			if v.Currency == "" {
				continue
			}
			m, ada := mu[v.Currency]
			if !ada {
				var err error
				if m, _, err = s.objekCargo.MataUangCargo(ctx, v.Currency); err != nil {
					return nil, err
				}
				mu[v.Currency] = m
			}
			v.CurrencyID, v.CurrencyOldID = m.ID, m.OldID
		}
		hasil = append(hasil, CargoTampil{Objek: c, TotalGross: totalGrossCargo(c.Coverages)})
	}
	return hasil, nil
}

// totalGrossCargo - SumTotalTSIPremiGross_Act langkah 4: per `.Currency.Name` (urutan kemunculan) ΣTSI, ΣPremium, Rate =
// @if(TSI=0,0,@divide(Premium,TSI,20)×100).
func totalGrossCargo(cov []models.CoverageCargo) []models.TotalGrossCargo {
	hasil := []models.TotalGrossCargo{}
	for _, c := range cov {
		n := -1
		for i := range hasil {
			if hasil[i].Currency == c.Currency {
				n = i
				break
			}
		}
		if n < 0 {
			hasil = append(hasil, models.TotalGrossCargo{Currency: c.Currency, TSI: nol(), Premium: nol()})
			n = len(hasil) - 1
		}
		hasil[n].TSI = tambah(hasil[n].TSI, nolBila(c.TSI))
		hasil[n].Premium = tambah(hasil[n].Premium, nolBila(c.Premium))
	}
	for i := range hasil {
		hasil[i].Rate = nol()
		if !hasil[i].TSI.IsZero() {
			if r, err := bagi(hasil[i].Premium, hasil[i].TSI, skalaBagi); err == nil {
				hasil[i].Rate = simpan8(kali(r, seratus)) // tampil 8 desimal (bentuk tersimpan)
			}
		}
	}
	return hasil
}

// GantiObjekCargo - PUT objek-cargo: daftar diganti utuh. Cargo ber-ID lama: coverage tidak dikirim (nil) -> coverage lama
// dipertahankan; dikirim -> spreading coverage lama di posisi yang sama dipertahankan (P3). Medan server dihitung ulang.
func (s *Service) GantiObjekCargo(ctx context.Context, pelaku inti.Pelaku, id string, masuk []models.ObjekCargo) ([]CargoTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return nil, err
	}
	if s.objekCargo == nil || s.transaksi == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	if !s.adalahMarineCargo(ctx, id) {
		return nil, ErrBukanMarineCargo
	}
	if err := s.cekKunciUnggah(ctx, id); err != nil { // tiket 74 P6: unggah CSV cargo berjalan (K75-2)
		return nil, err
	}
	cargo := salinCargo(masuk)
	if err := periksaCargo(cargo); err != nil {
		return nil, err
	}
	lama, err := s.objekCargo.BacaObjekCargo(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	lamaPerID := map[string]models.ObjekCargo{}
	for _, c := range lama {
		lamaPerID[c.ID] = c
	}
	share, err := s.objekCargo.PercentShareKasus(ctx, id)
	if err != nil {
		return nil, err
	}
	var proRate *apd.Decimal
	for i := range cargo {
		c := &cargo[i]
		l, adaLama := lamaPerID[c.ID]
		if err := s.medanServerObjekCargo(ctx, c, l, adaLama); err != nil {
			return nil, err
		}
		switch {
		case c.Coverages == nil && adaLama:
			c.Coverages = l.Coverages
		case c.Coverages == nil:
			c.Coverages = []models.CoverageCargo{}
		}
		for k := range c.Coverages {
			v := &c.Coverages[k]
			if v.Spreading == nil && adaLama && k < len(l.Coverages) {
				v.Spreading = l.Coverages[k].Spreading
			}
			if v.Spreading == nil {
				v.Spreading = []models.BarisSpreading{}
			}
			// ProRatePercent medan server: coverage lama di posisi yang sama, selain itu AppendCoverageMarineCargo_ACT.
			v.ProRatePercent = nil
			if adaLama && k < len(l.Coverages) {
				v.ProRatePercent = l.Coverages[k].ProRatePercent
			}
			if v.ProRatePercent == nil {
				if proRate == nil {
					proRate = s.proRateCargo(ctx, id)
				}
				v.ProRatePercent = proRate
			}
			if err := s.medanServerCoverageCargo(ctx, v, share); err != nil {
				return nil, err
			}
		}
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error { return s.objekCargo.GantiObjekCargo(ctx, tx, id, cargo) })
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	return s.BacaObjekCargo(ctx, id)
}

// salinCargo - salinan dalam (masukan pemanggil tidak diubah); larik nil (Coverages) tetap nil.
func salinCargo(masuk []models.ObjekCargo) []models.ObjekCargo {
	hasil := make([]models.ObjekCargo, len(masuk))
	for i, c := range masuk {
		c.Ship.AdditionalShip = append([]models.KapalTambahan{}, c.Ship.AdditionalShip...)
		if c.Coverages != nil {
			cov := make([]models.CoverageCargo, len(c.Coverages))
			for k, v := range c.Coverages {
				v.WarrantyExtension = append([]models.WarrantyExtension{}, v.WarrantyExtension...)
				v.Deductibles = append([]models.DeductibleCargo{}, v.Deductibles...)
				if v.Spreading != nil {
					v.Spreading = append([]models.BarisSpreading{}, v.Spreading...)
				}
				cov[k] = v
			}
			c.Coverages = cov
		}
		hasil[i] = c
	}
	return hasil
}

// medanServerObjekCargo - catatan master (ProtectObjectMarine), MaxAgeDetail, tanggal ke Pega, GRT / NRT / DWT koma -> titik,
// AGE (CountAgeShip).
func (s *Service) medanServerObjekCargo(ctx context.Context, c *models.ObjekCargo, lama models.ObjekCargo, adaLama bool) error {
	for _, m := range []struct {
		jenis       string
		id, idLama  string
		catatan     []*string
		catatanLama []string
	}{
		{MasterTrading, c.TradingID, lama.TradingID, []*string{&c.TradingNote}, []string{lama.TradingNote}},
		{MasterGoods, c.GoodID, lama.GoodID, []*string{&c.GoodNote}, []string{lama.GoodNote}},
		{MasterConveyance, c.ConveyanceID, lama.ConveyanceID, []*string{&c.ConveyanceName, &c.ConveyanceNote},
			[]string{lama.ConveyanceName, lama.ConveyanceNote}},
		{MasterPacking, c.PackingID, lama.PackingID, []*string{&c.PackingNote}, []string{lama.PackingNote}},
	} {
		isi, err := s.catatanMaster(c.ID, m.id, adaLama && m.idLama == m.id, m.catatanLama[0], func() (string, bool, error) {
			return s.masterCargo.Deskripsi(ctx, m.jenis, m.id)
		}, m.jenis)
		if err != nil {
			return err
		}
		for k, ke := range m.catatan {
			*ke = isi
			if k < len(m.catatanLama) && adaLama && m.idLama == m.id && s.masterCargo == nil {
				*ke = m.catatanLama[k]
			}
		}
	}
	detail, err := s.catatanMaster(c.ID, c.MaxAge, adaLama && lama.MaxAge == c.MaxAge, lama.MaxAgeDetail, func() (string, bool, error) {
		return s.masterCargo.WarrantyMaxAge(ctx, c.MaxAge)
	}, "MaxAge")
	if err != nil {
		return err
	}
	c.MaxAgeDetail = detail
	for _, d := range tanggalCargo {
		*d.ke(c) = tanggalKerugianKePega(*d.ke(c))
	}
	for _, x := range []*string{&c.Ship.GRT, &c.Ship.NRT, &c.Ship.DWT} {
		*x = strings.ReplaceAll(*x, ",", ".")
	}
	// AGE / FlagDelete / FlagOldData medan server (kabel tidak membawanya): CountAgeShip; Y_MAKE1 kosong -> AGE tersimpan cargo ini
	// tidak diubah (prakondisi L531). FlagDelete (hapus lunak EDM) / FlagOldData tidak diubah layar NB.
	c.Ship.Age, c.FlagDelete, c.FlagOldData = "", "", ""
	if adaLama {
		c.Ship.Age, c.FlagDelete, c.FlagOldData = lama.Ship.Age, lama.FlagDelete, lama.FlagOldData
	}
	if age, ok := umurKapal(c.SailDate, c.Ship.YMake1); ok {
		c.Ship.Age = age
	}
	return nil
}

// catatanMaster - catatan satu ID master: kosong -> ""; master terpasang -> Obj-Browse (tidak ada -> 400); belum terpasang ->
// catatan lama bila ID tidak berubah, selain itu ErrMasterCargoBelumAda.
func (s *Service) catatanMaster(idCargo, id string, sama bool, lama string, cari func() (string, bool, error), jenis string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", nil
	}
	if s.masterCargo == nil {
		if sama {
			return lama, nil
		}
		return "", fmt.Errorf("%w: %s", ErrMasterCargoBelumAda, jenis)
	}
	isi, ada, err := cari()
	if errors.Is(err, repository.ErrMasterBelumAda) { // master terpasang sebagian (MaxAge belum ada DDL): pola belum terpasang
		if sama {
			return lama, nil
		}
		return "", fmt.Errorf("%w: %s", ErrMasterCargoBelumAda, jenis)
	}
	if err != nil {
		return "", err
	}
	if !ada {
		return "", fmt.Errorf("%w: %s %q tidak ada di master", ErrMasukanObjek, jenis, id)
	}
	return isi, nil
}

// umurKapal - CountAgeShip: @toDecimal(@substring(SailDate,0,4)) − @toDecimal(Y_MAKE1); hanya bila Y_MAKE1 terisi (L531).
// Teks Pega SailDate "YYYYMMDDT…"; salah satu tidak terurai -> tidak diubah `[keputusan agent]` (perilaku @toDecimal atas teks
// tak terurai belum terverifikasi).
func umurKapal(sailPega, yMake1 string) (string, bool) {
	y := strings.TrimSpace(yMake1)
	if y == "" || len(sailPega) < 4 {
		return "", false
	}
	tahun, err1 := utils.ParseDecimal(sailPega[:4])
	buat, err2 := utils.ParseDecimal(y)
	if err1 != nil || err2 != nil {
		return "", false
	}
	return utils.FormatDecimal(kurang(tahun, buat)), true
}

// medanServerCoverageCargo - SetCurrencyCoverage_Act (ID / OldID), CountCoverageMarine (Discount, TSILiability),
// CountPremiAndTSIRNMCargo_ACT langkah 7 (RNM), pembulatan 8 desimal semua angka tersimpan.
func (s *Service) medanServerCoverageCargo(ctx context.Context, v *models.CoverageCargo, share *apd.Decimal) error {
	if err := s.mataUangCoverageCargo(ctx, v); err != nil {
		return err
	}
	if err := countCoverageMarineDiskon(v); err != nil {
		return err
	}
	rnmCargo(v, share)
	bulatkanCoverageCargo(v)
	return nil
}

// bulatkanCoverageCargo - bentuk tersimpan NUMBER(38,8) semua angka coverage (simpan8); dipakai PUT dan jawaban hitung (layar
// menampilkan nilai yang sama dengan yang akan disimpan) `[keputusan agent]`.
func bulatkanCoverageCargo(v *models.CoverageCargo) {
	for _, d := range desimalCoverageCargo {
		*d.ke(v) = simpan8(*d.ke(v))
	}
	v.Discount, v.TSILiability = simpan8(v.Discount), simpan8(v.TSILiability)
	v.TSINusantaraRe, v.PremiNusantaraRe = simpan8(v.TSINusantaraRe), simpan8(v.PremiNusantaraRe)
	v.ProRatePercent = simpan8(v.ProRatePercent)
	for k := range v.Deductibles {
		de := &v.Deductibles[k]
		de.PctDeductible, de.PctDeductible2, de.Amount = simpan8(de.PctDeductible), simpan8(de.PctDeductible2), simpan8(de.Amount)
	}
}

// mataUangCoverageCargo - SetCurrencyCoverage_Act: nama -> ID / OldID (CURRENCY); kosong -> kosong; tidak ada -> 400.
func (s *Service) mataUangCoverageCargo(ctx context.Context, v *models.CoverageCargo) error {
	v.CurrencyID, v.CurrencyOldID = "", ""
	if v.Currency == "" {
		return nil
	}
	m, ada, err := s.objekCargo.MataUangCargo(ctx, v.Currency)
	if err != nil {
		return err
	}
	if !ada {
		return fmt.Errorf("%w: currency %q tidak ada di master CURRENCY", ErrMasukanObjek, v.Currency)
	}
	v.CurrencyID, v.CurrencyOldID = m.ID, m.OldID
	return nil
}

// countCoverageMarineDiskon - CountCoverageMarine L1051-1073: Discount = @Math.divide(Premium×DiscountPercentage,100,4) (kosong
// = 0, pola @Math Pega), TSILiability = TSI.
func countCoverageMarineDiskon(v *models.CoverageCargo) error {
	v.TSILiability = v.TSI
	d, err := bagi(kali(nolBila(v.Premium), nolBila(v.DiscountPercentage)), seratus, 4)
	if err != nil {
		return err
	}
	v.Discount = d
	return nil
}

// rnmCargo - CountPremiAndTSIRNMCargo_ACT langkah 7: TSINusantaraRe = PS×TSI/100; PremiNusantaraRe = PS×Premium/100 −
// PS×Discount/100 (skala 20). PS kosong -> keduanya kosong.
func rnmCargo(v *models.CoverageCargo, share *apd.Decimal) {
	if share == nil {
		v.TSINusantaraRe, v.PremiNusantaraRe = nil, nil
		return
	}
	v.TSINusantaraRe = persenDari(nolBila(v.TSI), share)
	v.PremiNusantaraRe = kurang(persenDari(nolBila(v.Premium), share), persenDari(nolBila(v.Discount), share))
}

// countCoverageMarine - CountCoverageMarine L998-1073: Premium (premium.Calculate lini MARINE CARGO = @Math.divide(TSI×Rate,100,4)),
// lalu Discount + TSILiability.
func countCoverageMarine(v *models.CoverageCargo) error {
	hasil, err := premium.Calculate(premium.Input{LiniBisnis: premium.LiniMarineCargo, MataUang: v.Currency,
		TSI: teksAtauNol(v.TSI), Rate: teksAtauNol(v.Rate)})
	if err != nil {
		return err
	}
	v.Premium = hasil.Amount
	return countCoverageMarineDiskon(v)
}

// HitungCoverageCargo - POST hitung-coverage-cargo (tanpa menyimpan): rantai on-change sel pemicu + medan RNM (PS case).
func (s *Service) HitungCoverageCargo(ctx context.Context, id string, baris []BarisHitungCargo) ([]models.CoverageCargo, error) {
	var masalah []string
	for i, b := range baris {
		if !pemicuCargo[b.Pemicu] {
			masalah = append(masalah, "baris["+strconv.Itoa(i)+"].pemicu tidak dikenal")
		}
	}
	if len(masalah) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	if s.objekCargo == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	share, err := s.objekCargo.PercentShareKasus(ctx, id)
	if err != nil {
		return nil, err
	}
	hasil := make([]models.CoverageCargo, 0, len(baris))
	for i, b := range baris {
		v := b.Coverage
		switch b.Pemicu {
		case PemicuCargoCurrency:
			if err := s.mataUangCoverageCargo(ctx, &v); err != nil {
				return nil, fmt.Errorf("baris[%d]: %w", i, err)
			}
			err = countCoverageMarine(&v)
		case PemicuCargoPremium:
			// CountTSIandRateMarineCargo_Act L245: Rate = @Math.divide(Premium×100,TSI,4); TSI kosong / 0 -> Rate tidak diubah
			// (pembagian nol di Pega `belum terverifikasi`) `[keputusan agent]`.
			if v.TSI != nil && !v.TSI.IsZero() {
				v.Rate, err = bagi(kali(nolBila(v.Premium), seratus), v.TSI, 4)
			}
			if err == nil {
				err = countCoverageMarineDiskon(&v)
			}
		default:
			err = countCoverageMarine(&v)
		}
		if err != nil {
			return nil, fmt.Errorf("baris[%d]: %w", i, err)
		}
		rnmCargo(&v, share)
		bulatkanCoverageCargo(&v)
		hasil = append(hasil, v)
	}
	return hasil, nil
}
