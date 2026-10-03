package repository

// Uji SQL form penawaran dan pencarian master - tiket 01 bagian 3. TANPA Oracle.

import (
	"database/sql"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

func TestSimpanPenawaranHanyaMenyentuhKolomLayar(t *testing.T) {
	q := sqlSimpanPenawaran("SKEMAUJI.T_PREMIUM_LIST")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"CEDING_CO = :1", "CEDING_CO_NAME = :2", "POLICY_HOLDER = :3",
		"POLICY_HOLDER_NAME = :4", "TYPE_CEDING = :5", "TYPE_CEDING_NAME = :6", "BUSINESS_CODE = :7",
		"BUSINESS_NAME = :8", "DATE_RECEIVED = :9", "DESCRIPTION = :10",
		"BATAS_USIA_PESERTA = :11", "PERIODE_PERTANGGUNGAN = :12", "SUM_INSURED = :13", "TANGGAL_PENAWARAN = :14",
		"TANGGAL_RESPON = :15", "TANGGAL_KONFIRMASI = :16", "TBC = :17", "TANGGAL_TBC = :18", "STATUS_UPDATE = :19",
		"KETERANGAN_MARKETING = :20", "QQ_NAME = :21", "STATUS_FINAL = :27", "JENIS_ASURANSI = :28", "STATUS_PENAWARAN = :29", "WHERE ID = :30"} {
		if !strings.Contains(q, k) {
			t.Errorf("bentuk %q tidak ada:\n%s", k, q)
		}
	}
	// Kolom milik tahap lain tidak boleh tertimpa penawaran yang disimpan ulang.
	for _, k := range []string{" TYPE =", "PRODUCT_NAME", "SOB", "MARKETING_CODE", "MARKETING_NAME", "PRO_RATE_TYPE", "NO_OFFER"} {
		if strings.Contains(q, k) {
			t.Errorf("kolom %s ikut ditulis:\n%s", k, q)
		}
	}
	arg := argSimpanPenawaran("NBLF-1", models.PenawaranTersimpan{})
	if len(arg) != 30 || arg[29] != "NBLF-1" || arg[0] != nil || arg[8] != nil || arg[12] != nil {
		t.Errorf("argumen %v - kosong harus NULL, pengenal di :30", arg)
	}
}

func TestSuggestBernomorPerPolisDanBerpengenal32(t *testing.T) {
	q := sqlSisipSuggest("SKEMAUJI.T_VIEW_SUGGEST")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if got := urutanPenampung(q); got != "12345678" {
		t.Errorf("penampung %q", got)
	}
	if !strings.Contains(sqlNomorSuggestBerikut("S.T"), "NVL(MAX(NO), 0) + 1") ||
		!strings.Contains(sqlNomorSuggestBerikut("S.T"), "PREMIUM_LIST_ID = :1") {
		t.Error("nomor riwayat bukan MAX+1 per polis")
	}
	a, b := PengenalSuggest("NBLF-1", 1), PengenalSuggest("NBLF-1", 2)
	if len(a) != 32 || a == b || a != PengenalSuggest("NBLF-1", 1) {
		t.Errorf("pengenal %q / %q - harus 32 heksa, unik per nomor, deterministik", a, b)
	}
	if !strings.Contains(sqlRiwayatSuggest("S.T"), "ORDER BY s.NO DESC") {
		t.Error("riwayat tidak terbaru dahulu (Obj-Sort .No Descending)")
	}
}

// Saringan kedua RD, VERBATIM - dan teks cari lewat penanda, tidak disambung.
func TestPencarianMasterMengikutiReportDefinition(t *testing.T) {
	ceding := sqlCariCeding("POOLDATA.AGENT")
	for _, k := range []string{"STATUSACTIVE = :1", "ID LIKE :2", "UPPER(CLIENTNAME) LIKE :3", "FETCH FIRST"} {
		if !strings.Contains(ceding, k) {
			t.Errorf("ceding tanpa %q:\n%s", k, ceding)
		}
	}
	pemegang := sqlCariPemegangPolis("POOLDATA.CLIENT")
	for _, k := range []string{"NAME IS NOT NULL", "NAME <> :1", "BU_NOTE = :2", "UPPER(NAME) LIKE :3"} {
		if !strings.Contains(pemegang, k) {
			t.Errorf("pemegang polis tanpa %q:\n%s", k, pemegang)
		}
	}
	for _, q := range []string{ceding, pemegang} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
	}
	if got := polaLikeRujukan(" a%b_c "); got != `%A\%B\_C%` {
		t.Errorf("pola = %q - dibesarkan (SearchPolicyHolder_act) dan wildcard di-escape", got)
	}
}

// Satu daftar kolom merakit SELECT dan pengurainya; angka/tanggal keluar sebagai teks.
func TestBacaPenawaranDariSatuDaftarKolom(t *testing.T) {
	q := sqlBacaPenawaran("SKEMAUJI.T_PREMIUM_LIST")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"TO_CHAR(p.SUM_INSURED, 'TM9'", "TO_CHAR(p.TANGGAL_KONFIRMASI, 'YYYY-MM-DD",
		"TO_CHAR(p.TBC, 'TM9'", "p.KETERANGAN_MARKETING", "WHERE p.ID = :1"} {
		if !strings.Contains(q, k) {
			t.Errorf("bentuk %q tidak ada:\n%s", k, q)
		}
	}
	n := make([]sql.NullString, len(kolomBacaPenawaran))
	for i, k := range kolomBacaPenawaran {
		switch k.nama {
		case "SUM_INSURED":
			n[i] = sql.NullString{String: "6000000000.5", Valid: true}
		case "TBC":
			n[i] = sql.NullString{String: "25", Valid: true}
		case "TANGGAL_KONFIRMASI":
			n[i] = sql.NullString{String: "2026-09-29 00:00:00", Valid: true}
		case "STATUS_UPDATE":
			n[i] = sql.NullString{String: "UJI-STATUS", Valid: true}
		}
	}
	h, err := uraiPenawaran(n)
	if err != nil {
		t.Fatal(err)
	}
	if h.SumInsured == nil || h.SumInsured.Text('f') != "6000000000.5" || h.TBC == nil || *h.TBC != 25 ||
		h.TanggalKonfirmasi == nil || h.TanggalKonfirmasi.Day() != 29 || h.StatusUpdate != "UJI-STATUS" {
		t.Errorf("hasil = %+v", h)
	}
	if h.BatasUsiaPeserta != nil || h.TanggalPenawaran != nil {
		t.Error("NULL menjadi nilai, bukan nil")
	}
}
