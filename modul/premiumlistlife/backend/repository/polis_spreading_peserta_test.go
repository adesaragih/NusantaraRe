package repository

// Popup spreading SATU peserta (keputusan work owner 07-10-2026). TANPA Oracle.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSqlSpreadingSatuPeserta(t *testing.T) {
	qi := sqlIdentitasPesertaSpreading("S.T_PREMIUM_LIST_DETAIL")
	qr := sqlRincianSpreadingPeserta("S.T_PREMIUM_LIST_SPREADING", "S.T_PREMIUM_LIST_SPREADING_RETRO", "S.TREATYYEAR_LIFE")
	for _, q := range []string{qi, qr} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(q, "SELECT *") {
			t.Errorf("SELECT * : %s", q)
		}
	}
	// ⛔ Kepemilikan: peserta dicari DI DALAM polis, bukan berdasarkan ID saja.
	for _, k := range []string{"WHERE d.PREMIUM_LIST_ID = :1 AND d.ID = :2", "d.RETROCEDED_SHARE"} {
		if !strings.Contains(qi, k) {
			t.Errorf("identitas tanpa %q: %s", k, qi)
		}
	}
	for _, k := range []string{
		"JOIN S.T_PREMIUM_LIST_SPREADING_RETRO r ON r.SPREADING_ID = s.ID",
		"LEFT JOIN S.TREATYYEAR_LIFE t ON TO_CHAR(t.ID) = s.TREATY_YEAR_LIFE",
		"WHERE s.DETAIL_ID = :1", "ORDER BY s.ID, r.ID",
	} {
		if !strings.Contains(qr, k) {
			t.Errorf("rincian tanpa %q: %s", k, qr)
		}
	}
}

// baris sqlRincianSpreadingPeserta: s.ID, type, tahunID, tahun, mulai, akhir, bagian, reinsurer, 7 angka.
func barisUjiSpreading(sID, tipe, bagian, re, rate string) []string {
	return []string{sID, tipe, "T2", "2026", "2026-01-01", "2026-12-31", bagian, re, "60", "1800000000", rate, "20", "10",
		"1", ".5"}
}

func TestSusunSpreadingPeserta(t *testing.T) {
	p := SpreadingSatuPeserta{ID: "UJI-D1"}
	SusunSpreadingPeserta(&p, [][]string{
		barisUjiSpreading("S01", "QS", "3000000000", "UJI-RE-A", "2.5"),
		barisUjiSpreading("S01", "QS", "3000000000", "UJI-RE-B", "2.5"),
		barisUjiSpreading("S02", "SURPLUS", "1500000000", "UJI-RE-C", ".72"),
	})
	if p.TreatyYear != "2026" || p.TreatyMulai != "2026-01-01" || p.TreatyAkhir != "2026-12-31" {
		t.Errorf("treaty year = %q %q %q", p.TreatyYear, p.TreatyMulai, p.TreatyAkhir)
	}
	if len(p.Lapisan) != 2 || p.Lapisan[0].ReinsType != "QS" || p.Lapisan[1].ReinsType != "SURPLUS" {
		t.Fatalf("lapisan = %+v", p.Lapisan)
	}
	if n := len(p.Lapisan[0].Reinsurer); n != 2 || p.Lapisan[0].Reinsurer[1].Nama != "UJI-RE-B" {
		t.Errorf("reinsurer lapisan 1 = %+v", p.Lapisan[0].Reinsurer)
	}
	// angkaUtuh: ".72" → "0.72" - teks desimal seragam untuk satu formatter layar.
	if r := p.Lapisan[1].Reinsurer[0]; r.Rate != "0.72" || r.Net != "0.5" {
		t.Errorf("angka = %+v", r)
	}
	// Tanpa label Transaction Year → ID treaty year.
	q := SpreadingSatuPeserta{}
	b := barisUjiSpreading("S01", "QS", "1", "UJI-RE", "1")
	b[3] = ""
	SusunSpreadingPeserta(&q, [][]string{b})
	if q.TreatyYear != "T2" {
		t.Errorf("treaty year cadangan = %q", q.TreatyYear)
	}
	// Nol baris → lapisan KOSONG (bukan nil) - jawaban sah, bukan galat.
	k := SpreadingSatuPeserta{}
	SusunSpreadingPeserta(&k, nil)
	if k.Lapisan == nil || len(k.Lapisan) != 0 || k.TreatyYear != "" {
		t.Errorf("kosong = %+v", k)
	}
}
