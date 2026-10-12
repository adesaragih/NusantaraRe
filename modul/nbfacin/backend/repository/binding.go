package repository

// Accept tahap Underwriting ke jalur akhir (tiket 63): Flow\InputInwardFacultativeOffer.xml Decision23 Else -> Utility16
// "SAVE JSON_OFFER" (Activity\SaveJsonOfferFacIn_Act.xml) -> Decision16 ELSE -> Transition38 -> Assignment9 "MARKETING
// (BINDING)". Semua dalam transaksi Submit (keputusan agent D2: gagal -> Submit batal; Pega mengabaikan galatnya).
//
//   - JSON_OFFER: PORT Go prosedur POOLDATA.PEGA_M_JSON_OFFER (`[terverifikasi]` DDL\PEGA_M_JSON_OFFER.txt) TANPA CALL
//     (ADR-0043): ada baris ber-IDPEGA -> UPDATE JSONDATA; belum -> NOOFFER cabang P_BusinessFac 'F' = 'RNM-OFR.' ||
//     MM.YYYY || '.' || LPAD(RNM_JSON_OFFER_SEQ.NEXTVAL, 5, '0') (LPAD memotong ke 5 karakter - ditiru), ID = '1' ||
//     JSON_OFFER_SEQ.NEXTVAL, PRODKE 0, TGL_INPUT bawaan sysdate. Cabang 'F' `[belum terverifikasi]` (nilai
//     Quotation.BusinessFac NB FIRE, pertanyaan terbuka Product, issues/63). JSONDATA = dokumen minimal (keputusan work owner
//     D1; ADR-0028).
//   - FACINOFFER: Activity\SaveOfferProduction_Act.xml langkah 11 RDB-Delete (`DELETE FROM POOLDATA.FACINOFFER WHERE IDPEGA
//     ={pyWorkPage.pyID}`) lalu cabang FIRE langkah 12 satu INSERT per baris SpreadingList (RDBList\InsertOfferProduction_Sql,
//     50 kolom). NOOFFER = JSON_OFFER.NOOFFER case (GetNoOffer_Sql), SPREAD_DATE = JSON_OFFER.TGL_INPUT nomor itu
//     (GetTglInputOffer_sql). Perbaikan cacat (keputusan work owner D3): tanggal 24-jam (Pega "hh" 12-jam), literal "NULL"
//     -> SQL NULL, galat INSERT -> transaksi batal.
//   - T_WORK_POLIS.FLAG_ONGOING_POLICY = 2 (Transition38).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelJSONOffer / TabelFacInOffer - tabel lama POOLDATA, ditulis Accept jalur akhir (tiket 63).
	TabelJSONOffer  = "JSON_OFFER"
	TabelFacInOffer = "FACINOFFER"
	// sequenceJSONOffer / sequenceNoOffer - sequence prosedur PEGA_M_JSON_OFFER (DDL sequence tidak ada di korpus;
	// keberadaan / akses = pertanyaan terbuka DBA, issues/63).
	sequenceJSONOffer = "JSON_OFFER_SEQ"
	sequenceNoOffer   = "RNM_JSON_OFFER_SEQ"
)

// ErrSequenceJSONOffer - sequence JSON_OFFER tidak ada / tanpa akses (ORA-02289): galat jelas, bukan diam (D2).
var ErrSequenceJSONOffer = errors.New("repository: sequence JSON_OFFER_SEQ / RNM_JSON_OFFER_SEQ tidak ada atau tanpa akses " +
	"(pertanyaan terbuka DBA, tiket 63)")

func sqlFlagBinding(work string) string {
	return "UPDATE " + work + " SET FLAG_ONGOING_POLICY = :1 WHERE ID = :2 AND LINI = :3"
}

// sqlCariJSONOffer - COUNT prosedur + nomor offer yang sudah ada (GetNoOffer_Sql memungut nomor baris case itu).
func sqlCariJSONOffer(jo string) string {
	return "SELECT COUNT(*), MAX(NOOFFER) FROM " + jo + " WHERE IDPEGA IN (:1, :2)" // dua bentuk: case lama tidak disisip ganda
}

func sqlSisipJSONOffer(jo string) string {
	return "INSERT INTO " + jo + " (ID, IDPEGA, JSONDATA, PRODKE, NOOFFER) VALUES (:1, :2, :3, 0, :4)"
}

func sqlUbahJSONOffer(jo string) string {
	return "UPDATE " + jo + " SET JSONDATA = :1 WHERE IDPEGA IN (:2, :3)"
}

func sqlHapusFacInOffer(fo string) string { return "DELETE FROM " + fo + " WHERE IDPEGA IN (:1, :2)" }

// kolomFacInOffer - urutan kolom InsertOfferProduction_Sql (50); GROUPNAME / TOP_RISK di akhir seperti sumbernya.
var kolomFacInOffer = []string{"IDPEGA", "NOOFFER", "BUSINESSCODE", "CEDINGCO", "SOB", "BEGINDATE", "ENDDATE", "PROVINCE",
	"KABUPATEN", "KECAMATAN", "KELURAHAN", "ROADNAME", "BUILDINGNO", "ZIP_CODE", "OBJECTNO", "OBJECT_NAME", "RISK_LOCATION",
	"OBJECT_ITEM", "COVERAGE", "COVERAGENOTE", "RATE", "PERCENT_RNM", "PERCENT_SHARE_SPREADED", "JN_REAS", "CURR_ID",
	"TSI_MENJADI", "PREMI_MENJADI", "TSI_SELISIH", "PREMI_SELISIH", "SPREAD_DATE", "ACCUMULATIONCODE", "OCCUPATIONCODE",
	"MARKETINGOFFICERCODE", "INSUREDNAME", "TRADING_NAME", "CONVEYANCE_NAME", "PERCENT_RI_COMM", "RISLIPCEDING", "LICENSEPLATE",
	"TYPE_NAME", "MODEL_NAME", "GROUP_TYPE", "NOSERIAL", "ENGINENO", "IDSHIP", "FROMRUTE", "TORUTE", "SAILDATE", "GROUPNAME",
	"TOP_RISK"}

// nilaiKolomFacInOffer - ekspresi VALUES per kolom: "?" = bind berikutnya; selain itu teks SQL tetap (NULL untuk kolom yang
// Pega isi "" / "NULL"). Tanggal To_date 24-jam; RATE = Replace(rate, '.', ',') seperti sumbernya; SPREAD_DATE = TGL_INPUT
// JSON_OFFER ber-NOOFFER itu.
func nilaiKolomFacInOffer(kolom, jo string) string {
	switch kolom {
	case "BEGINDATE", "ENDDATE", "SAILDATE":
		return "TO_DATE(?, 'DD/MM/YYYY HH24:MI:SS')"
	case "RATE":
		return "REPLACE(?, '.', ',')"
	case "PERCENT_RNM", "PERCENT_SHARE_SPREADED", "TSI_MENJADI", "PREMI_MENJADI", "TSI_SELISIH", "PREMI_SELISIH", "PERCENT_RI_COMM":
		return angkaMasuk("?")
	case "SPREAD_DATE":
		return "(SELECT MAX(j.TGL_INPUT) FROM " + jo + " j WHERE j.NOOFFER = ?)"
	case "LICENSEPLATE", "TYPE_NAME", "MODEL_NAME", "ENGINENO":
		return "NULL"
	}
	return "?"
}

// sqlSisipFacInOffer - INSERT 50 kolom; bind :1..:n berurutan menurut kemunculan (go-ora posisional).
func sqlSisipFacInOffer(fo, jo string) string {
	nilai := make([]string, len(kolomFacInOffer))
	n := 0
	for i, k := range kolomFacInOffer {
		e := nilaiKolomFacInOffer(k, jo)
		if strings.Contains(e, "?") {
			n++
			e = strings.Replace(e, "?", fmt.Sprintf(":%d", n), 1)
		}
		nilai[i] = e
	}
	return "INSERT INTO " + fo + " (" + strings.Join(kolomFacInOffer, ", ") + ") VALUES (" + strings.Join(nilai, ", ") + ")"
}

// argFacInOffer - nilai bind sqlSisipFacInOffer, urutan kolomFacInOffer (kolom NULL tetap dilewati).
func argFacInOffer(b models.BarisFacInOffer, noOffer string) []any {
	k := db.KosongJadiNil
	nilai := map[string]any{"IDPEGA": k(b.IDPega), "NOOFFER": k(noOffer), "BUSINESSCODE": k(b.BusinessCode),
		"CEDINGCO": k(b.CedingCo), "SOB": k(b.SOB), "BEGINDATE": k(b.BeginDate), "ENDDATE": k(b.EndDate),
		"PROVINCE": k(b.Province), "KABUPATEN": k(b.Kabupaten), "KECAMATAN": k(b.Kecamatan), "KELURAHAN": k(b.Kelurahan),
		"ROADNAME": k(b.RoadName), "BUILDINGNO": k(b.BuildingNo), "ZIP_CODE": k(b.ZipCode), "OBJECTNO": k(b.ObjectNo),
		"OBJECT_NAME": k(b.ObjectName), "RISK_LOCATION": k(b.RiskLocation), "OBJECT_ITEM": k(b.ObjectItem),
		"COVERAGE": k(b.Coverage), "COVERAGENOTE": k(b.CoverageNote), "RATE": k(b.Rate), "PERCENT_RNM": ikatDesimal(b.PercentRNM),
		"PERCENT_SHARE_SPREADED": ikatDesimal(b.PercentShareSpreaded), "JN_REAS": k(b.JnReas), "CURR_ID": k(b.CurrID),
		"TSI_MENJADI": ikatDesimal(b.TSI), "PREMI_MENJADI": ikatDesimal(b.Premi), "TSI_SELISIH": ikatDesimal(b.TSI),
		"PREMI_SELISIH": ikatDesimal(b.Premi), "SPREAD_DATE": k(noOffer), "ACCUMULATIONCODE": k(b.AccumulationCode),
		"OCCUPATIONCODE": k(b.OccupationCode), "MARKETINGOFFICERCODE": k(b.MarketingOfficerCode),
		"INSUREDNAME": k(b.InsuredName), "PERCENT_RI_COMM": ikatDesimal(b.PercentRIComm), "RISLIPCEDING": k(b.RISlipCeding),
		"GROUP_TYPE": k(b.GroupType), "GROUPNAME": k(b.GroupName), "TOP_RISK": k(b.TopRisk),
		"NOSERIAL": k(b.NoSerial), // tiket 73 P4: Aneka cabang 13 `.SerialNo`; "" = NULL (FIRE / Bonding / Golf "NULL")
		// tiket 74 P5: cabang Marine Cargo langkah 16; "" = NULL (cabang lain, sama dengan NULL tetap sebelumnya)
		"TRADING_NAME": k(b.TradingName), "CONVEYANCE_NAME": k(b.ConveyanceName), "IDSHIP": k(b.IDShip),
		"FROMRUTE": k(b.FromRute), "TORUTE": k(b.ToRute), "SAILDATE": k(b.SailDate)}
	var arg []any
	for _, kol := range kolomFacInOffer {
		if strings.Contains(nilaiKolomFacInOffer(kol, ""), "?") {
			arg = append(arg, nilai[kol])
		}
	}
	return arg
}

// noOfferBaru - 'RNM-OFR.' || bulan || '.' || LPAD(n, 5, '0'); LPAD Oracle MEMOTONG teks yang lebih panjang dari 5.
func noOfferBaru(bulan, n string) string {
	if len(n) >= 5 {
		n = n[:5]
	} else {
		n = strings.Repeat("0", 5-len(n)) + n
	}
	return "RNM-OFR." + bulan + "." + n
}

// galatSequence - ORA-02289 (sequence does not exist) -> ErrSequenceJSONOffer.
func galatSequence(err error) error {
	if err != nil && strings.Contains(err.Error(), "ORA-02289") {
		return fmt.Errorf("%w: %v", ErrSequenceJSONOffer, err)
	}
	return err
}

// tulisBinding - lihat kepala berkas; dipanggil TulisSubmit sesudah posisi berhasil diubah.
func (r *SubmitOracle) tulisBinding(ctx context.Context, tx *db.Tx, t tabelSubmit, caseID, idPega string, b models.TulisanBinding) error {
	if _, err := jalankan(ctx, tx, sqlFlagBinding(t.work), "FLAG_ONGOING_POLICY binding", b.FlagBaru, caseID, LiniFacIn); err != nil {
		return err
	}
	q := sqlCariJSONOffer(t.jo)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	var ada int
	var noOffer sql.NullString
	polos, awalan := DuaBentukIDPega(idPega)
	if err := tx.QueryRowContext(ctx, q, polos, awalan).Scan(&ada, &noOffer); err != nil {
		return fmt.Errorf("repository: baca %s: %w", TabelJSONOffer, err)
	}
	if ada == 0 {
		nomor, err := r.db.NomorBerikut(ctx, tx, sequenceNoOffer)
		if err != nil {
			return galatSequence(err)
		}
		urut, err := r.db.NomorBerikut(ctx, tx, sequenceJSONOffer)
		if err != nil {
			return galatSequence(err)
		}
		noOffer = sql.NullString{String: noOfferBaru(b.BulanNoOffer, nomor), Valid: true}
		if _, err := jalankan(ctx, tx, sqlSisipJSONOffer(t.jo), TabelJSONOffer, "1"+urut, idPega, b.JSONData, noOffer.String); err != nil {
			return err
		}
	} else if _, err := jalankan(ctx, tx, sqlUbahJSONOffer(t.jo), TabelJSONOffer, b.JSONData, polos, awalan); err != nil {
		return err
	}
	fp, fa := DuaBentukIDPega(caseID)
	if _, err := jalankan(ctx, tx, sqlHapusFacInOffer(t.fo), "hapus "+TabelFacInOffer, fp, fa); err != nil {
		return err
	}
	sisip := sqlSisipFacInOffer(t.fo, t.jo)
	for _, o := range b.Offer {
		if _, err := jalankan(ctx, tx, sisip, TabelFacInOffer, argFacInOffer(o, noOffer.String)...); err != nil {
			return err
		}
	}
	return nil
}

// ---- penerima email "Waiting for Bind" (SendEmailPolicy) --------------------------------------------------------------

// sqlEmailAkun - email akun ber-LOGIN_ID itu (Obj-Browse Data-Admin-Operator-ID .pyEmailAddress; padanan sistem baru
// M_LOGIN_GO, seperti korespondensi tiket 56).
func sqlEmailAkun(login string) string {
	return "SELECT MAX(EMAIL) FROM " + login + " WHERE UPPER(TRIM(LOGIN_ID)) = UPPER(TRIM(:1))"
}

// sqlEmailMarketing - GetMKTandLeader_SQL (`ID = MOID OR ID = (SELECT CLIENTID2 … WHERE ID = MOID)`; MAX mencegah ORA-01427)
// -> AKSES_LOGIN -> email akun; urut ID.
func sqlEmailMarketing(mo, login string) string {
	return "SELECT (SELECT MAX(l.EMAIL) FROM " + login + " l WHERE UPPER(TRIM(l.LOGIN_ID)) = UPPER(TRIM(m.AKSES_LOGIN))) FROM " + mo +
		" m WHERE m.ID = :1 OR m.ID = (SELECT MAX(m2.CLIENTID2) FROM " + mo + " m2 WHERE m2.ID = :2) ORDER BY m.ID"
}

// PenerimaBinding - email pembuat case dan email marketing + atasannya (kosong dibuang).
func (r *SubmitOracle) PenerimaBinding(ctx context.Context, createOp, moid string) (string, []string, error) {
	t, err := r.tabel()
	if err != nil {
		return "", nil, err
	}
	var pembuat sql.NullString
	if err := r.db.QueryRowContext(ctx, sqlEmailAkun(t.login), createOp).Scan(&pembuat); err != nil {
		return "", nil, fmt.Errorf("repository: email pembuat case: %w", err)
	}
	baris, err := r.db.QueryContext(ctx, sqlEmailMarketing(t.mo, t.login), moid, moid)
	if err != nil {
		return "", nil, fmt.Errorf("repository: email marketing: %w", err)
	}
	defer baris.Close()
	var cc []string
	for baris.Next() {
		var e sql.NullString
		if err := baris.Scan(&e); err != nil {
			return "", nil, fmt.Errorf("repository: email marketing: %w", err)
		}
		if s := strings.TrimSpace(e.String); s != "" {
			cc = append(cc, s)
		}
	}
	if err := baris.Err(); err != nil {
		return "", nil, fmt.Errorf("repository: email marketing: %w", err)
	}
	return strings.TrimSpace(pembuat.String), cc, nil
}

// ---- form tahap MARKETING (BINDING) (tiket 64) -----------------------------------------------------------------------

var (
	// ErrSudahBinding - CONFIRM_BINDING sudah "1" (tombol Binding disabled `.OfferFacIn.ConfirmBinding = 1`).
	ErrSudahBinding = errors.New("repository: penawaran sudah di-Binding")
	// ErrPenawaranTidakAda - JSON_OFFER case belum ada (Accept jalur akhir tiket 63 belum terjadi).
	ErrPenawaranTidakAda = errors.New("repository: JSON_OFFER case belum tersimpan")
)

// sqlStatusBinding - CONFIRM_BINDING / RECEIVED_RI_SLIP (migrasi 213).
func sqlStatusBinding(general string) string {
	return "UPDATE " + general + " SET CONFIRM_BINDING = :1, RECEIVED_RI_SLIP = :2 WHERE ID = :3"
}

// sqlSentuhTahap - penjaga serentak tombol Binding: kasus masih di posisi + flag ini.
// sqlSentuhTahapCatatan - sqlSentuhTahap dengan posisi efektif (tiket 77).
func sqlSentuhTahapCatatan(work string) string {
	return "UPDATE " + work + " w SET TGL_UPDATE = SYSDATE WHERE w.ID = :1 AND w.LINI = :2 AND " + sqlPosisiEfektif("w") +
		" = :3 AND w.FLAG_ONGOING_POLICY = :4"
}

func sqlSentuhTahap(work string) string {
	return "UPDATE " + work + " SET TGL_UPDATE = SYSDATE WHERE ID = :1 AND LINI = :2 AND POSITION = :3 AND FLAG_ONGOING_POLICY = :4"
}

// sqlKonfirmasiBinding - ConfirmBinding = 1 (UpdateTglBindJsonOffer_Act langkah 3) hanya bila belum.
func sqlKonfirmasiBinding(general string) string {
	return "UPDATE " + general + " SET CONFIRM_BINDING = '1' WHERE ID = :1 AND (CONFIRM_BINDING IS NULL OR CONFIRM_BINDING <> '1')"
}

// sqlTglBind - UpdateTgl_BindJsonOffer_SQL (`UPDATE JSON_OFFER a SET a.TGL_BIND = sysdate, jsondata = … WHERE
// a.JSONDATA.IDNewBisnis = pyLabel`) dikunci IDPEGA seperti tiket 63 (keputusan agent D6); JSONDATA minimal (tiket 63 D1).
func sqlTglBind(jo string) string {
	return "UPDATE " + jo + " SET TGL_BIND = SYSDATE, JSONDATA = :1 WHERE IDPEGA IN (:2, :3)"
}

// TulisTekanBinding - lihat PenyimpanSubmit. FACINPERFORMANCE TIDAK ditulis di sini (keputusan work owner D4: BINDING/ACCEPT
// sekali saat Submit Accept).
func (r *SubmitOracle) TulisTekanBinding(ctx context.Context, tx *db.Tx, b models.TekanBinding) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	satu := func(q, aksi string, galat error, arg ...any) error {
		h, err := jalankan(ctx, tx, q, aksi, arg...)
		if err != nil {
			return err
		}
		n, err := h.RowsAffected()
		if err != nil {
			return fmt.Errorf("repository: %s: %w", aksi, err)
		}
		if n == 0 {
			return galat
		}
		return nil
	}
	qSentuh := sqlSentuhTahap(t.work)
	if adaCatatanPosisi(ctx, r.db, t.work) {
		qSentuh = sqlSentuhTahapCatatan(t.work) // tiket 77
	}
	if err := satu(qSentuh, "posisi Binding", ErrPosisiBerubah, b.CaseID, LiniFacIn, b.Position, b.Flag); err != nil {
		return err
	}
	if err := satu(sqlKonfirmasiBinding(t.general), "CONFIRM_BINDING", ErrSudahBinding, b.CaseID); err != nil {
		return err
	}
	if err := satu(sqlTglBind(t.jo), "TGL_BIND "+TabelJSONOffer, ErrPenawaranTidakAda, b.JSONData, b.IDPega,
		AwalanIDPega+IDPegaProduksi(b.IDPega)); err != nil {
		return err
	}
	k := db.KosongJadiNil
	p := b.Riwayat
	_, err = jalankan(ctx, tx, sqlSisipRiwayatProduksi(t.produksi), "riwayat "+TabelRiwayatProduksi, b.IDPega, k(p.TypePolis),
		k(p.Posisi), k(p.PIC), p.TglInp, k(p.Div), k(p.Type), k(p.Putaran), k(p.Approval), k(p.Keterangan), k(p.AksesLogin),
		ikatDesimal(p.PercentRNM), b.IDPega, AwalanIDPega+IDPegaProduksi(b.IDPega))
	return err
}
