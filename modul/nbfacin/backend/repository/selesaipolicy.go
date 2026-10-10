package repository

// Selesai Policy (tiket 69a-2): port Activity\SaveJsonPolicyFacIn_Act.xml jalur NB FIRE dalam transaksi Submit - nol CALL
// (ADR-0043), nol COMMIT (ADR-U-0029). Keputusan work owner K69a2-1..9 (docs/issues/69-jalur-policy.md). Urutan:
//   - nomor polis (langkah 3-12): JSON_POLIS case sudah ada -> NOPOLIS / TGL_PROD-nya dipakai ulang (Pega: PolicyNo terisi ->
//     langkah 11-12 dilewati); belum -> TANGGAL_CLOSING (kosong -> 25, langkah 4), ProdDateTime (langkah 5), periode +
//     penghitung GENERATE_SEQUENCE_NUMBER lewat inti/backend/penomor (port PROC_GENERATE_SEQUENCE_NUMBER), awalan
//     KODE_PRODUKSI TYPE 'NONLIFE' (GetKodeProdNonLife_SQL);
//   - JSON_POLIS (langkah 15, port INSERTJSONPOLIS): 9 kolom, DATA_JSON minimal, PRODKE '0', NOENDORS NULL;
//   - CEDING_FACINPRODUCTION (langkah 16 InsertCedingProduction): DELETE per IDPEGA lalu INSERT per ceding (port
//     InsertUpdateCedingProduction; baris kembar IDPEGA + CEDINGID + PCT_SHARE dilewati). PRODKE / NOPOLIS_ARAS
//     dihitung untuk baris itu saja - UPDATE tanpa WHERE prosedur atas seluruh tabel TIDAK ditiru (K69a2-1);
//   - M_RISK_LOSS_PROFILE (langkah 21): sekali per NOPOLIS (K69a2-8);
//   - FACINPRODUCTION (langkah 23 -> SaveTreatyProduction_Act -> SaveFacinProdFireNB_Act): hanya bila IDPEGA belum ada (CekFacin,
//     K69a2-7). Tanggal 24 jam (K69a2-2), teks "NULL" ditiru (K69a2-4), MARKETINGOFFICERCODE = GetMOID_Sql (K69a2-3).
// Tidak diport: JSON_POLIS_ERROR (langkah 17) dan UpdateFacProdObjItemID_SQL (langkah 27) - K69a2-1.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// Tabel lama POOLDATA yang ditulis selesai Policy (JSON_POLIS: TabelJSONPolis, FACINPRODUCTION: TabelFacInProduction).
	TabelCedingProduksi     = "CEDING_FACINPRODUCTION"
	TabelProfilRisikoLoss   = "M_RISK_LOSS_PROFILE"
	tipeKodeProduksiNonLife = "NONLIFE" // GetKodeProdNonLife_SQL `TYPE ='NONLIFE'`
	// bentukTanggalPolis - "DD/MM/YYYY HH24:MI:SS" (TO_DATE kolom tanggal).
	bentukTanggalPolis = "02/01/2006 15:04:05"
)

type tabelPolicy struct{ jp, ceding, risiko, prod, mo string }

func (r *SubmitOracle) tabelPolicy() (tabelPolicy, error) {
	var t tabelPolicy
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelJSONPolis, &t.jp}, {TabelCedingProduksi, &t.ceding}, {TabelProfilRisikoLoss, &t.risiko},
		{TabelFacInProduction, &t.prod}, {TabelMarketingOfficer, &t.mo}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

// ---- SQL ---------------------------------------------------------------------------------------------------------

// sqlCariJSONPolis - JSON_POLIS case (PK IDPEGA): NOPOLIS + TGL_PROD (jam dinding).
func sqlCariJSONPolis(jp string) string {
	return "SELECT NOPOLIS, TO_CHAR(TGL_PROD, 'DD/MM/YYYY HH24:MI:SS') FROM " + jp + " WHERE IDPEGA = :1"
}

// sqlSisipJSONPolis - INSERT prosedur INSERTJSONPOLIS (TGL_INPUT bawaan sysdate; PRODKE '0'; NOENDORS / OLDNOPOLIS NULL - NB).
func sqlSisipJSONPolis(jp string) string {
	return "INSERT INTO " + jp + " (IDPEGA, DATA_JSON, TGL_INPUT, NOPOLIS, PRODKE, OLDNOPOLIS, TGL_PROD, USERNAME, NOENDORS)" +
		" VALUES (:1, :2, SYSDATE, :3, '0', NULL, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, NULL)"
}

func sqlHapusCedingProduksi(ceding string) string {
	return "DELETE FROM " + ceding + " WHERE IDPEGA = :1"
}

// sqlSisipCedingProduksi - INSERT prosedur InsertUpdateCedingProduction + nilai yang UPDATE-nya isi untuk baris itu (PRODKE dari
// JSON_POLIS, NOPOLIS_ARAS = NOPOLIS tanpa titik). NOENDORS_ARAS TIDAK ditulis: kolom itu ada di DDL korpus
// (DDL\CEDING_FACINPRODUCTION.txt baris 11) tetapi TIDAK di skema DEV (ORA-00904, Accept TL Policy NB-151029, 09-10-2026);
// nilainya selalu NULL (NOENDORS NULL -> REPLACE(NULL) = NULL), jadi tanpa kolom itu setara Pega `[keputusan agent 0f]`.
func sqlSisipCedingProduksi(ceding, jp string) string {
	return "INSERT INTO " + ceding + " (IDPEGA, NOPOLIS, NOENDORS, CEDINGID, CEDINGNAME, PCT_SHARE, PRODKE, NOPOLIS_ARAS)" +
		" VALUES (:1, :2, NULL, :3, :4, " + angkaMasuk(":5") + ", (SELECT MAX(j.PRODKE) FROM " + jp +
		" j WHERE j.IDPEGA = :6), REPLACE(:7, '.', ''))"
}

func sqlAdaProfilRisiko(risiko string) string {
	return "SELECT COUNT(*) FROM " + risiko + " WHERE NOPOLIS = :1"
}

// sqlSisipProfilRisiko - InsertRiskAndLossProfile_SQL (TANGGALSPREADING = SYSDATE).
func sqlSisipProfilRisiko(risiko string) string {
	return "INSERT INTO " + risiko + " (NOPOLIS, TANGGALSPREADING, TSINUSANTARARE, PREMINUSANTARARE, PREMIUMOR, PREMIUMQS," +
		" PREMIUMFACOUT, COB) VALUES (:1, SYSDATE, " + angkaMasuk(":2") + ", " + angkaMasuk(":3") + ", " + angkaMasuk(":4") + ", " +
		angkaMasuk(":5") + ", " + angkaMasuk(":6") + ", :7)"
}

// sqlAdaFacInProduction - CekFacin_Sql (`select distinct(IDPEGA) ... WHERE IDPEGA = pzInsKey`).
func sqlAdaFacInProduction(prod string) string {
	return "SELECT COUNT(*) FROM " + prod + " WHERE IDPEGA = :1"
}

// kolomFacInProduction - urutan kolom InsertTreatyProduction_Sql (79).
var kolomFacInProduction = []string{"IDPEGA", "NOPOLIS", "NOENDORS", "BUSINESSCODE", "CEDINGCO", "SOB", "BEGINDATE", "ENDDATE",
	"PROVINCE", "KABUPATEN", "KECAMATAN", "KELURAHAN", "ROADNAME", "BUILDINGNO", "ZIP_CODE", "OBJECTNO", "OBJECT_NAME",
	"RISK_LOCATION", "OBJECT_ITEM", "COVERAGE", "COVERAGENOTE", "RATE", "PERCENT_RNM", "PERCENT_SHARE_SPREADED", "JN_REAS",
	"CURR_ID", "TSI_MENJADI", "PREMI_MENJADI", "TSI_SELISIH", "PREMI_SELISIH", "SPREAD_DATE", "ACCUMULATIONCODE",
	"OCCUPATIONCODE", "MARKETINGOFFICERCODE", "INSUREDNAME", "TRADING_NAME", "CONVEYANCE_NAME", "PERCENT_RI_COMM",
	"RISLIPCEDING", "LICENSEPLATE", "TYPE_NAME", "MODEL_NAME", "GROUP_TYPE", "NOSERIAL", "ENGINENO", "QQNAME", "PRORATE",
	"PCT_RI_COMM_SELISIH", "RICOMM", "RICOMM_SELISIH", "OBJECT_NAME_ID", "LAYERNO", "IDSHIP", "FROMRUTE", "TORUTE", "SAILDATE",
	"POLICYTYPE", "LOL_MENJADI", "LOL_SELISIH", "PCT_BROKERGARE_FEE", "BROKERAGE_FEE_MENJADI", "BROKERAGE_FEE_SELISIH",
	"PCT_BROKERAGE_FEE_SELISIH", "CEDINGID", "SOBID", "INSUREDID", "IDX_OBJECT_ITEM", "TSI100_MENJADI", "TSI100_SELISIH",
	"BINDER", "ISB2B", "ASKRED", "STATUS_BUSINESS", "GROUPNAME", "TOP_RISK", "NO_MASTER_POLIS", "TYPE_EDM",
	"DEDUCTION2_MENJADI", "DEDUCTION2_SELISIH"}

// nilaiKolomFacInProduction - ekspresi VALUES per kolom: "?" = bind berikutnya; selain itu teks SQL tetap. Pega
// `To_number(Replace(x, ',', '.'))` -> angkaMasuk (RATE VARCHAR2: konversi implisit seperti Pega). Teks tetap: Pega mengisi
// "" (NULL) atau "NULL" (K69a2-4) / 1 (PRORATE) untuk FIRE; halaman Pega yang tidak dikenal sistem baru (BinderRNM, IsB2B,
// NPLStatus, PolicyMasterNumber) dan yang tidak di-set di jalur FIRE (SAILDATE, TYPE_EDM, DEDUCTION2_*; NoEndors tak
// diteruskan) = NULL.
func nilaiKolomFacInProduction(kolom, mo string) string {
	switch kolom {
	case "BEGINDATE", "ENDDATE", "SPREAD_DATE":
		return "TO_DATE(?, 'DD/MM/YYYY HH24:MI:SS')"
	case "RATE", "PERCENT_RNM", "PERCENT_SHARE_SPREADED", "TSI_MENJADI", "PREMI_MENJADI", "TSI_SELISIH", "PREMI_SELISIH",
		"PERCENT_RI_COMM", "PCT_RI_COMM_SELISIH", "RICOMM", "RICOMM_SELISIH", "LOL_MENJADI", "LOL_SELISIH", "PCT_BROKERGARE_FEE",
		"BROKERAGE_FEE_MENJADI", "BROKERAGE_FEE_SELISIH", "PCT_BROKERAGE_FEE_SELISIH", "TSI100_MENJADI", "TSI100_SELISIH":
		return angkaMasuk("?")
	case "MARKETINGOFFICERCODE":
		// GetMOID_Sql: `select id ... from marketingofficer where (clientid = Quotation.MarketingCode or clientid =
		// QuotationData.MarketingCode) and mostatus = '1'`, pxResults(1) - urutan Pega tak tentu; MIN(ID) `[keputusan agent]`.
		// Tiket 73 P5: bind pertama = MOID case cabang Aneka / Golf (S11 L2492 `@if(QuotationData.MOID=="",GetMOID,MOID)`, K73-46);
		// FIRE mengirim NULL -> NVL jatuh ke GetMOID (FireNB F6 menimpa S11) - hasil FIRE tidak berubah.
		return "NVL(?, (SELECT MIN(m.ID) FROM " + mo + " m WHERE m.CLIENTID = ? AND m.MOSTATUS = '1'))"
	case "NOSERIAL":
		return "?" // tiket 73 P5: `.SerialNo` objek Aneka; selain itu teks "NULL" (argFacInProduction)
	case "LICENSEPLATE", "ENGINENO":
		return "'NULL'"
	case "PRORATE":
		return "1"
	case "NOENDORS", "TRADING_NAME", "CONVEYANCE_NAME", "TYPE_NAME", "MODEL_NAME", "IDSHIP", "FROMRUTE", "TORUTE", "SAILDATE",
		"BINDER", "ISB2B", "ASKRED", "NO_MASTER_POLIS", "TYPE_EDM", "DEDUCTION2_MENJADI", "DEDUCTION2_SELISIH":
		return "NULL"
	}
	return "?"
}

// sqlSisipFacInProduction - INSERT 79 kolom; bind :1..:n berurutan menurut kemunculan (go-ora posisional).
func sqlSisipFacInProduction(prod, mo string) string {
	nilai := make([]string, len(kolomFacInProduction))
	n := 0
	for i, k := range kolomFacInProduction {
		e := nilaiKolomFacInProduction(k, mo)
		for strings.Contains(e, "?") { // MARKETINGOFFICERCODE memuat dua bind (tiket 73 P5)
			n++
			e = strings.Replace(e, "?", fmt.Sprintf(":%d", n), 1)
		}
		nilai[i] = e
	}
	return "INSERT INTO " + prod + " (" + strings.Join(kolomFacInProduction, ", ") + ") VALUES (" + strings.Join(nilai, ", ") + ")"
}

// argFacInProduction - nilai bind sqlSisipFacInProduction, urutan kolomFacInProduction (kolom tetap dilewati).
func argFacInProduction(b models.BarisFacInProduction, noPolis, spreadDate string) []any {
	k := db.KosongJadiNil
	d := ikatDesimal
	nilai := map[string]any{"IDPEGA": k(b.IDPega), "NOPOLIS": k(noPolis), "BUSINESSCODE": k(b.BusinessCode),
		"CEDINGCO": k(b.CedingCo), "SOB": k(b.SOB), "BEGINDATE": k(b.BeginDate), "ENDDATE": k(b.EndDate),
		"PROVINCE": k(b.Province), "KABUPATEN": k(b.Kabupaten), "KECAMATAN": k(b.Kecamatan), "KELURAHAN": k(b.Kelurahan),
		"ROADNAME": k(b.RoadName), "BUILDINGNO": k(b.BuildingNo), "ZIP_CODE": k(b.ZipCode), "OBJECTNO": k(b.ObjectNo),
		"OBJECT_NAME": k(b.ObjectName), "RISK_LOCATION": k(b.RiskLocation), "OBJECT_ITEM": k(b.ObjectItem),
		"COVERAGE": k(b.Coverage), "COVERAGENOTE": k(b.CoverageNote), "RATE": d(b.Rate), "PERCENT_RNM": d(b.PercentRNM),
		"PERCENT_SHARE_SPREADED": d(b.PercentShareSpreaded), "JN_REAS": k(b.JnReas), "CURR_ID": k(b.CurrID),
		"TSI_MENJADI": d(b.TSI), "PREMI_MENJADI": d(b.Premi), "TSI_SELISIH": d(b.TSI), "PREMI_SELISIH": d(b.Premi),
		"SPREAD_DATE": k(spreadDate), "ACCUMULATIONCODE": k(b.AccumulationCode), "OCCUPATIONCODE": k(b.OccupationCode),
		"MARKETINGOFFICERCODE": k(b.MarketingCode), "INSUREDNAME": k(b.InsuredName), "PERCENT_RI_COMM": d(b.PercentRIComm),
		"RISLIPCEDING": k(b.RISlipCeding), "GROUP_TYPE": k(b.GroupType), "QQNAME": k(b.QQName),
		"PCT_RI_COMM_SELISIH": d(b.PctRICommSelisih), "RICOMM": d(b.RIComm), "RICOMM_SELISIH": d(b.RIComm),
		"OBJECT_NAME_ID": k(b.ObjectNameID), "LAYERNO": k(b.LayerNo), "POLICYTYPE": k(b.PolicyType), "LOL_MENJADI": d(b.LoL),
		"LOL_SELISIH": d(b.LoL), "PCT_BROKERGARE_FEE": d(b.PctBrokerageFee), "BROKERAGE_FEE_MENJADI": d(b.BrokerageFee),
		"BROKERAGE_FEE_SELISIH": d(b.BrokerageFee), "PCT_BROKERAGE_FEE_SELISIH": d(b.PctBrokerageFee),
		"CEDINGID": k(b.CedingID), "SOBID": k(b.SOBID), "INSUREDID": k(b.InsuredID), "IDX_OBJECT_ITEM": k(b.IdxObjectItem),
		"TSI100_MENJADI": d(b.TSI100), "TSI100_SELISIH": d(b.TSI100), "STATUS_BUSINESS": k(b.StatusBusiness),
		"GROUPNAME": k(b.GroupName), "TOP_RISK": k(b.TopRisk), "NOSERIAL": noSerialProduksi(b.NoSerial)}
	var arg []any
	for _, kol := range kolomFacInProduction {
		if kol == "MARKETINGOFFICERCODE" { // NVL(MOID case, GetMOID(clientid))
			arg = append(arg, k(b.MOIDKasus), nilai[kol])
			continue
		}
		if strings.Contains(nilaiKolomFacInProduction(kol, ""), "?") {
			arg = append(arg, nilai[kol])
		}
	}
	return arg
}

// noSerialProduksi - NOSERIAL: SerialNo objek Aneka (SaveTreatyProduction_Act L5004); kosong -> teks "NULL" (L5109; FIRE F L4149,
// Golf L11363) - K69a2-4.
func noSerialProduksi(s string) string {
	if strings.TrimSpace(s) == "" {
		return "NULL"
	}
	return s
}

// ---- tulis ---------------------------------------------------------------------------------------------------------

// nomorPolis - (NOPOLIS, ProdDateTime "DD/MM/YYYY HH24:MI:SS" WIB, baru). JSON_POLIS case sudah ada -> dipakai ulang.
func (r *SubmitOracle) nomorPolis(ctx context.Context, tx *db.Tx, jp, idPega string, p models.TulisanPolicy) (string, string, bool, error) {
	q := sqlCariJSONPolis(jp)
	if err := db.PeriksaSQL(q); err != nil {
		return "", "", false, err
	}
	var no, prod sql.NullString
	err := tx.QueryRowContext(ctx, q, idPega).Scan(&no, &prod)
	if err == nil {
		return no.String, prod.String, false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", "", false, fmt.Errorf("repository: baca %s: %w", TabelJSONPolis, err)
	}
	pen := penomor.NewPenomor(r.db)
	closing, err := pen.HariClosing(ctx, tx)
	if errors.Is(err, penomor.ErrTanggalClosingKosong) {
		closing = models.HariClosingBawaan
	} else if err != nil {
		return "", "", false, err
	}
	jakarta := penomor.DiJakarta(p.Sekarang).Location()
	saat := models.ProdDateTime(p.Sekarang, p.MulaiPolis, closing, jakarta)
	periode, err := penomor.HitungPeriodeNomor(saat, closing)
	if err != nil {
		return "", "", false, err
	}
	kode, err := pen.AwalanProduksi(ctx, tx, tipeKodeProduksiNonLife)
	if err != nil {
		return "", "", false, err
	}
	jenis := models.JenisNomorPolis(kode)
	// TANGGAL penghitung = p_proddate prosedur (TO_DATE ProdDateTime "dd/MM/yyyy").
	lokal := saat.In(jakarta)
	hari := time.Date(lokal.Year(), lokal.Month(), lokal.Day(), 0, 0, 0, 0, jakarta)
	urut, err := pen.UrutNomorBerikut(ctx, tx, models.KelasNomorPolis, jenis, periode, hari)
	if err != nil {
		return "", "", false, err
	}
	return models.RakitNomorPolis(jenis, p.BusinessOldID, periode.MMYYYY, urut), lokal.Format(bentukTanggalPolis), true, nil
}

// hitung - SELECT COUNT(*) satu bind.
func hitung(ctx context.Context, tx *db.Tx, q, aksi string, arg any) (int, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := tx.QueryRowContext(ctx, q, arg).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: %s: %w", aksi, err)
	}
	return n, nil
}

// tulisPolicy - lihat kepala berkas; dipanggil TulisSubmit sesudah posisi berhasil diubah. p.NomorPolis diisi.
func (r *SubmitOracle) tulisPolicy(ctx context.Context, tx *db.Tx, idPega string, p *models.TulisanPolicy) error {
	t, err := r.tabelPolicy()
	if err != nil {
		return err
	}
	noPolis, prod, baru, err := r.nomorPolis(ctx, tx, t.jp, idPega, *p)
	if err != nil {
		return err
	}
	p.NomorPolis = noPolis
	k := db.KosongJadiNil
	if baru {
		if _, err := jalankan(ctx, tx, sqlSisipJSONPolis(t.jp), TabelJSONPolis, idPega, p.JSONData, noPolis, prod,
			k(p.Username)); err != nil {
			return err
		}
	}
	if _, err := jalankan(ctx, tx, sqlHapusCedingProduksi(t.ceding), "hapus "+TabelCedingProduksi, idPega); err != nil {
		return err
	}
	sudah := map[string]bool{}
	for _, c := range p.Ceding {
		kunci := c.CedingID + "\x00" + fmt.Sprint(ikatDesimal(c.Share))
		if sudah[kunci] {
			continue // prosedur: COUNT IDPEGA + CEDINGID + PCT_SHARE > 0 -> "Data sudah di simpan"
		}
		sudah[kunci] = true
		if _, err := jalankan(ctx, tx, sqlSisipCedingProduksi(t.ceding, t.jp), TabelCedingProduksi, idPega, noPolis,
			k(c.CedingID), k(c.CedingName), ikatDesimal(c.Share), idPega, noPolis); err != nil {
			return err
		}
	}
	n, err := hitung(ctx, tx, sqlAdaProfilRisiko(t.risiko), "baca "+TabelProfilRisikoLoss, noPolis)
	if err != nil {
		return err
	}
	if rl := p.RisikoLoss; n == 0 {
		if _, err := jalankan(ctx, tx, sqlSisipProfilRisiko(t.risiko), TabelProfilRisikoLoss, noPolis, ikatDesimal(rl.TSINusantaraRe),
			ikatDesimal(rl.PremiNusantaraRe), ikatDesimal(rl.PremiumOR), ikatDesimal(rl.PremiumQS), ikatDesimal(rl.PremiumFacout),
			k(rl.COB)); err != nil {
			return err
		}
	}
	n, err = hitung(ctx, tx, sqlAdaFacInProduction(t.prod), "CekFacin "+TabelFacInProduction, idPega)
	if err != nil || n > 0 {
		return err
	}
	sisip := sqlSisipFacInProduction(t.prod, t.mo)
	for _, b := range p.Produksi {
		if _, err := jalankan(ctx, tx, sisip, TabelFacInProduction, argFacInProduction(b, noPolis, prod)...); err != nil {
			return err
		}
	}
	return nil
}
