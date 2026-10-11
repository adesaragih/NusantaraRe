package services

// Tulisan produksi case MARINE CARGO (tiket 74 P5, keputusan work owner 10-10-2026 K74-9/10 "ikut Pega apa adanya", K74-11
// M_RISK_LOSS_PROFILE = Aneka / K73-48). Kolom Marine terbukti di korpus `[terverifikasi]` (DDL\FACINOFFER.txt L17-50,
// DDL\FACINPRODUCTION.txt L18-63; RDBList NB InsertOfferProduction_Sql / InsertTreatyProduction_Sql):
//   - barisFacInOfferCargo = SaveOfferProduction_Act langkah 16 (WHEN IsMarineCargo): per CargoList × CoverageList × SpreadingList;
//     16.1 OBJECTNO = .PackingID, OBJECT_NAME = .PackingNote, OBJECT_ITEM = .GoodNote, TRADING_NAME = .TradingNote,
//     CONVEYANCE_NAME = .ConveyanceNote, IDSHIP = .PolicyData.Ship.ID, FROMRUTE / TORUTE, SAILDATE = @FormatDateTime(SailDate,
//     "dd/MM/yyyy hh:mm:ss", Asia/Jakarta) (hh 12-jam -> 12:00 WIB, K74-9b); alamat / okupasi / risk location "" (NULL);
//     LICENSEPLATE / NOSERIAL / ENGINENO Pega "NULL" (konvensi FACINOFFER sistem baru: NULL); 16.2.1 COVERAGE / CURR_ID =
//     Currency.Name / COVERAGENOTE / ACCUMULATIONCODE (coverage cargo tidak menyimpannya: "") / RATE; 16.2.2.1 share, jenis treaty,
//     TSI / premi spreaded. TOP_RISK (CARI48) hanya di-SET cabang FIRE -> NULL.
//   - barisFacInProductionCargo = SaveTreatyProduction_Act langkah 22 (IsMarineCargo, "NB-"): kolom 22.1 sama dengan di atas +
//     RICOMM / BROKERAGE = PremiumSpreaded × % / 100, LOL_MENJADI / LOL_SELISIH = ClaimEstimation (kosong 0.0); langkah bersama
//     11 / 16: MOID case didahulukan, OBJECT_NAME_ID teks "NULL", LAYERNO NULL (pola Aneka tiket 73 P5); IDX_OBJECT_ITEM
//     (Datain1.CARI29) / TOP_RISK (CARI49) / TSI100 tidak di-SET -> NULL (Datain1.CARI30 = indeks coverage tidak masuk INSERT).
//     DEDUCTION2_* = PremiumSpreaded × Payment.PctDeduction2 / 100 di Pega; PctDeduction2 TIDAK tersimpan di sistem baru (tab
//     Premium Deduction menyimpan Deduction2 per mata uang) -> NULL seperti FIRE / Aneka (K73-45) - keputusan work owner 11-10-2026 K74-22.
//   - Satu baris per baris spreading (K74-9), sama dengan port FIRE / Aneka. Catatan korpus: setiap cabang SaveOfferProduction_Act /
//     SaveTreatyProduction_Act (FIRE juga) memberi prakondisi `.pxListSubscript>1` pada langkah INSERT; keputusan work owner
//     11-10-2026 K74-23: tetap tulis semua baris spreading (FIRE / Aneka tidak diubah).
//   - IDSHIP VARCHAR2(20) sedangkan POOLDATA.SHIP.ID VARCHAR2(100): ID kapal > 20 karakter -> ErrIDKapalFacIn (409), tidak dipotong.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// bentukSailDateFacIn - "dd/MM/yyyy hh:mm:ss" Pega (hh 12-jam) -> To_date HH24.
	bentukSailDateFacIn = "02/01/2006 03:04:05"
	// batasIDKapalFacIn - FACINOFFER / FACINPRODUCTION.IDSHIP VARCHAR2(20).
	batasIDKapalFacIn = 20
)

// ErrIDKapalFacIn - ID kapal cargo tidak muat kolom IDSHIP. 409.
var ErrIDKapalFacIn = errors.New("ID kapal cargo lebih dari 20 karakter - tidak muat kolom IDSHIP FACINOFFER / FACINPRODUCTION")

// tanggalSailFacIn - teks Pega DateTime SailDate -> "dd/MM/yyyy hh:mm:ss" WIB 12-jam; kosong / tak terurai -> "" (NULL).
func tanggalSailFacIn(teks string) string {
	t, err := time.Parse(BentukWaktuPega, strings.TrimSpace(teks))
	if err != nil {
		return ""
	}
	return t.In(WIB).Format(bentukSailDateFacIn)
}

// idKapalFacIn - Ship.ID cargo ke-n (1..), wajib muat IDSHIP.
func idKapalFacIn(n int, c models.ObjekCargo) (string, error) {
	id := strings.TrimSpace(c.Ship.ID)
	if utf8.RuneCountInString(id) > batasIDKapalFacIn || len(id) > batasIDKapalFacIn {
		return "", fmt.Errorf("%w (cargo %d)", ErrIDKapalFacIn, n)
	}
	return id, nil
}

// barisFacInOfferCargo - lihat kepala berkas; dasar = kolom tingkat case (dasarFacInOffer).
func barisFacInOfferCargo(dasar models.BarisFacInOffer, cargo []models.ObjekCargo) ([]models.BarisFacInOffer, error) {
	hasil := []models.BarisFacInOffer{}
	for n, c := range cargo {
		idKapal, err := idKapalFacIn(n+1, c)
		if err != nil {
			return nil, err
		}
		lok := dasar
		lok.ObjectNo, lok.ObjectName, lok.ObjectItem = c.PackingID, c.PackingNote, c.GoodNote
		lok.TradingName, lok.ConveyanceName, lok.IDShip = c.TradingNote, c.ConveyanceNote, idKapal
		lok.FromRute, lok.ToRute, lok.SailDate = c.FromRute, c.ToRute, tanggalSailFacIn(c.SailDate)
		for _, v := range c.Coverages {
			for _, sb := range v.Spreading {
				b := lok
				b.Coverage, b.CoverageNote, b.CurrID = v.Coverage, v.CoverageNote, v.Currency
				if v.Rate != nil {
					b.Rate = utils.FormatDecimal(v.Rate)
				}
				b.PercentShareSpreaded, b.JnReas, b.TSI, b.Premi = sb.SharePercentage, sb.TreatyType, sb.TSISpreaded, sb.PremiumSpreaded
				hasil = append(hasil, b)
			}
		}
	}
	return hasil, nil
}

// barisFacInProductionCargo - lihat kepala berkas; d = dasar FACINPRODUCTION (dasarBarisProduksi).
func barisFacInProductionCargo(d dasarProduksi, cargo []models.ObjekCargo) ([]models.BarisFacInProduction, error) {
	ri := nolBila(d.k.RICommision)
	dasar := dasarBarisProduksi(d)
	dasar.MOIDKasus = strings.TrimSpace(d.k.MOID)
	hasil := []models.BarisFacInProduction{}
	for n, c := range cargo {
		idKapal, err := idKapalFacIn(n+1, c)
		if err != nil {
			return nil, err
		}
		lok := dasar
		lok.ObjectNo, lok.ObjectName, lok.ObjectItem = c.PackingID, c.PackingNote, c.GoodNote
		lok.TradingName, lok.ConveyanceName, lok.IDShip = c.TradingNote, c.ConveyanceNote, idKapal
		lok.FromRute, lok.ToRute, lok.SailDate = c.FromRute, c.ToRute, tanggalSailFacIn(c.SailDate)
		lok.ObjectNameID, lok.LayerNo = "NULL", layerNoNull
		for _, v := range c.Coverages {
			for _, sb := range v.Spreading {
				b := lok
				b.Coverage, b.CoverageNote, b.CurrID, b.Rate = v.Coverage, v.CoverageNote, v.Currency, simpan8(v.Rate)
				b.PercentShareSpreaded, b.JnReas = simpan8(sb.SharePercentage), sb.TreatyType
				b.TSI, b.Premi = simpan8(sb.TSISpreaded), simpan8(sb.PremiumSpreaded)
				b.LoL = simpan8(nolBila(sb.ClaimEstimation))
				b.RIComm = simpan8(persenDari(sb.PremiumSpreaded, ri))
				b.BrokerageFee = simpan8(persenDari(sb.PremiumSpreaded, d.pctBrokerage))
				hasil = append(hasil, b)
			}
		}
	}
	return hasil, nil
}

// cargoProduksi - pohon cargo case MARINE CARGO untuk Binding / selesai Policy (ada = false bila bukan cargo atau penyimpan cargo
// tidak terpasang - cabang lain berjalan apa adanya).
func (s *Service) cargoProduksi(ctx context.Context, id string) ([]models.ObjekCargo, bool, error) {
	if s.objekCargo == nil || !s.adalahMarineCargo(ctx, id) {
		return nil, false, nil
	}
	cargo, err := s.objekCargo.BacaObjekCargo(ctx, id)
	if err != nil {
		return nil, false, galatKasus(err)
	}
	return cargo, true, nil
}
