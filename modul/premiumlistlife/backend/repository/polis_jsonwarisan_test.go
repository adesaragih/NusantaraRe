package repository

// JSON_OFFER_LIFE (Save Offer) dan JSON_POLIS (Confirm) - keputusan work owner
// 06-10-2026. Bentuk SQL saja; nama kolom = ALL_TAB_COLUMNS POOLDATA.

import (
	"database/sql"
	"strings"
	"testing"
)

// kolomJSONOfferLifeSah - seluruh kolom JSON_OFFER_LIFE (ALL_TAB_COLUMNS 06-10-2026).
var kolomJSONOfferLifeSah = map[string]bool{
	"ID": true, "TGL_INPUT": true, "IDPEGA": true, "BUSINESSCODE": true, "BUSINESSNAME": true, "CEDINGCO": true,
	"CEDINGCONAME": true, "POLICYHOLDER": true, "POLICYHOLDERNAME": true, "RECEIVED_DATE": true,
	"OFFERING_DATE": true, "RESPONSE_DATE": true, "CONFIRMATION_DATE": true, "MAX_TBC": true,
	"RECONFIRMATION_DATE": true, "REALIZATION_DATE": true, "BINDING_DATE": true, "STATUS": true,
	"STATUS_UPDATE": true, "STATUS_FINAL": true, "MARKETING_NOTE": true, "OPERATORNAME": true,
	"REINSURANCETYPE": true, "OLDID": true, "JSONDATA": true, "TYPECEDING": true, "UW_NOTE": true,
}

// lebarJSONOfferLifeSah - lebar VARCHAR2 ALL_TAB_COLUMNS.
var lebarJSONOfferLifeSah = map[string]int{
	"BUSINESSCODE": 10, "BUSINESSNAME": 100, "CEDINGCO": 100, "CEDINGCONAME": 100, "POLICYHOLDER": 100,
	"POLICYHOLDERNAME": 100, "STATUS": 10, "STATUS_UPDATE": 1000, "STATUS_FINAL": 1000, "MARKETING_NOTE": 1000,
	"OPERATORNAME": 100, "REINSURANCETYPE": 50, "TYPECEDING": 50, "UW_NOTE": 1000,
}

func TestKolomJSONOfferLife(t *testing.T) {
	isi := map[string]bool{}
	for _, k := range kolomJSONOfferLife {
		if !kolomJSONOfferLifeSah[k.Kolom] {
			t.Errorf("kolom %s bukan kolom JSON_OFFER_LIFE", k.Kolom)
		}
		if isi[k.Kolom] {
			t.Errorf("kolom %s ganda", k.Kolom)
		}
		isi[k.Kolom] = true
		if k.Lebar != lebarJSONOfferLifeSah[k.Kolom] {
			t.Errorf("lebar %s = %d, mau %d", k.Kolom, k.Lebar, lebarJSONOfferLifeSah[k.Kolom])
		}
		if k.Lebar > 0 && k.Label == "" {
			t.Errorf("kolom %s berlebar tanpa label pesan", k.Kolom)
		}
	}
	// Semua kolom terisi kecuali ID/IDPEGA (penampung) dan JSONDATA/OLDID (NULL).
	for k := range kolomJSONOfferLifeSah {
		switch k {
		case "ID", "IDPEGA":
		case "JSONDATA", "OLDID":
			if isi[k] {
				t.Errorf("%s tidak boleh diisi", k)
			}
		default:
			if !isi[k] {
				t.Errorf("kolom %s tidak diisi", k)
			}
		}
	}
}

func TestSqlJSONOfferLife(t *testing.T) {
	const tujuan, polis, urutan = "SKEMAUJI.JSON_OFFER_LIFE", "SKEMAUJI.T_PREMIUM_LIST", "SKEMAUJI.JSON_OFFER_SEQ"
	if q := sqlHapusJSONOffer(tujuan); q != "DELETE FROM SKEMAUJI.JSON_OFFER_LIFE WHERE IDPEGA = :1" {
		t.Errorf("hapus = %q", q)
	}
	if q := sqlIDJSONOfferLama(tujuan); !strings.Contains(q, "SELECT ID FROM SKEMAUJI.JSON_OFFER_LIFE WHERE IDPEGA = :1") {
		t.Errorf("id lama = %q", q)
	}
	lama := sqlSisipJSONOffer(tujuan, polis, "")
	baru := sqlSisipJSONOffer(tujuan, polis, urutan)
	for _, k := range []string{
		"INSERT INTO SKEMAUJI.JSON_OFFER_LIFE (ID, IDPEGA, BUSINESSCODE, BUSINESSNAME",
		"SELECT :1, :2, p.BUSINESS_CODE", "WHERE p.ID = :3 AND p.ID = p.ID_PEGA",
		"p.STATUS_PENAWARAN", "p.JENIS_ASURANSI", "p.KETERANGAN_MARKETING", "p.CREATE_OP_NAME",
		"p.KETENTUAN_UNDERWRITING", "TRUNC(p.DATE_RECEIVED)", "TRUNC(p.TANGGAL_BIND)", "TRUNC(p.TANGGAL_TBC)",
		"WHEN '4' THEN 'XOL' END", "UW_NOTE, TGL_INPUT)", "SYSDATE FROM",
	} {
		if !strings.Contains(lama, k) {
			t.Errorf("sisip (ID lama) tanpa %q:\n%s", k, lama)
		}
	}
	for _, k := range []string{"SELECT TO_CHAR(SKEMAUJI.JSON_OFFER_SEQ.NEXTVAL), :1, p.BUSINESS_CODE",
		"WHERE p.ID = :2 AND p.ID = p.ID_PEGA"} {
		if !strings.Contains(baru, k) {
			t.Errorf("sisip (baru) tanpa %q:\n%s", k, baru)
		}
	}
	if strings.Contains(lama, "NEXTVAL") {
		t.Error("ID lama tidak boleh menggerakkan sequence")
	}
	if q := sqlLebarJSONOffer(polis); !strings.HasPrefix(q, "SELECT LENGTHB(p.BUSINESS_CODE), ") ||
		strings.Count(q, "LENGTHB(") != len(lebarJSONOfferLifeSah) {
		t.Errorf("lebar = %q", q)
	}
	for _, q := range []string{lama, baru, sqlHapusJSONOffer(tujuan), sqlIDJSONOfferLama(tujuan), sqlLebarJSONOffer(polis)} {
		for _, larangan := range []string{"SELECT *", "COMMIT", "JSONDATA", "OLDID", "INSERTJSONOFFERLIFE"} {
			if strings.Contains(strings.ToUpper(q), larangan) {
				t.Errorf("SQL memuat %q: %s", larangan, q)
			}
		}
	}
}

// TestLebarJSONOfferDitolakBerpesan - STATUS/BUSINESSCODE VARCHAR2(10): ditolak di Go.
func TestLebarJSONOfferDitolakBerpesan(t *testing.T) {
	var panjang []sql.NullInt64
	for _, k := range kolomJSONOfferLife {
		if k.Lebar == 0 {
			continue
		}
		n := sql.NullInt64{Int64: int64(k.Lebar), Valid: true}
		switch k.Kolom {
		case "BUSINESSCODE":
			n.Int64 = 11
		case "STATUS":
			n.Int64 = 12
		case "UW_NOTE":
			n.Valid = false // NULL muat
		}
		panjang = append(panjang, n)
	}
	kurang := periksaLebarJSONOffer(panjang)
	if len(kurang) != 2 || kurang[0] != "Business Code max 10 characters (now 11)" ||
		kurang[1] != "Status max 10 characters (now 12)" {
		t.Errorf("kurang = %q", kurang)
	}
}

func TestSqlJSONPolis(t *testing.T) {
	const tujuan, polis = "SKEMAUJI.JSON_POLIS", "SKEMAUJI.T_PREMIUM_LIST"
	if q := sqlHapusJSONPolis(tujuan); q != "DELETE FROM SKEMAUJI.JSON_POLIS WHERE IDPEGA = :1" {
		t.Errorf("hapus = %q", q)
	}
	q := sqlSisipJSONPolis(tujuan, polis)
	for _, k := range []string{
		"INSERT INTO SKEMAUJI.JSON_POLIS (IDPEGA, NOPOLIS, TGL_PROD, USERNAME, TGL_INPUT)",
		"SELECT :1, :2, :3, p.CREATE_OP_NAME, SYSDATE FROM SKEMAUJI.T_PREMIUM_LIST p",
		"WHERE p.ID = :4 AND p.ID = p.ID_PEGA",
	} {
		if !strings.Contains(q, k) {
			t.Errorf("sisip tanpa %q:\n%s", k, q)
		}
	}
	// USERNAME = CREATEOPNAME LIFEINPRODUCTION.
	if nilaiProduksiLife["CREATEOPNAME"] != "p.CREATE_OP_NAME" {
		t.Error("USERNAME JSON_POLIS harus sama dengan CREATEOPNAME LIFEINPRODUCTION")
	}
	for _, s := range []string{q, sqlHapusJSONPolis(tujuan), sqlTglProdLama(tujuan)} {
		for _, larangan := range []string{"SELECT *", "COMMIT", "DATA_JSON", "INSERTJSONPOLISLIFE", "PRODKE",
			"NOENDORS", "ERR_NOTE", "KONVERSI", "OLDNOPOLIS"} {
			if strings.Contains(strings.ToUpper(s), larangan) {
				t.Errorf("SQL memuat %q: %s", larangan, s)
			}
		}
	}
}
