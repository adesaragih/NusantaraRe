package repository

// Filter A RD `GetListOpportunity` (`A.pxCreateOperator = Param.UserIdentifier`; keputusan work owner 06-10-2026
// "isi inbox ini muncul hanya untuk akun dia saja"): berkas buatan pelaku, ATAU berkas di antrean atasan yang ia
// pegang. Nilai selalu terikat.

import (
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestSqlDaftarKasusFilterPembuat(t *testing.T) {
	q, args := sqlDaftarKasus("S.W", "S.G", "S.Q", models.SaringanKasus{Pembuat: "UJI-A"})
	if !strings.Contains(q, "AND w.CREATE_OP = :4") || args[3] != "UJI-A" {
		t.Fatalf("hanya pembuat:\n%s\n%v", q, args)
	}
	q, args = sqlDaftarKasus("S.W", "S.G", "S.Q", models.SaringanKasus{Pembuat: "UJI-A",
		Antrean: []string{models.PosisiSecHead, models.PosisiDeptHead}})
	if !strings.Contains(q, "AND (w.CREATE_OP = :4 OR g.POSITION_NOTE IN (:5, :6))") || len(args) != 6 {
		t.Fatalf("pembuat ATAU antrean:\n%s\n%v", q, args)
	}
}

// Daftar kotak masuk Beranda: pembuat dibatasi posisi (buatan akun yang MASIH di Admin), tetap ATAU antrean.
func TestSqlDaftarKasusPembuatDiPosisi(t *testing.T) {
	q, args := sqlDaftarKasus("S.W", "S.G", "S.Q", models.SaringanKasus{Pembuat: "UJI-A", PembuatPosisi: models.PosisiAdmin,
		Antrean: []string{models.PosisiSecHead}})
	if !strings.Contains(q, "AND ((w.CREATE_OP = :4 AND g.POSITION_NOTE = :5) OR g.POSITION_NOTE IN (:6))") || len(args) != 6 {
		t.Fatalf("pembuat di posisi:\n%s\n%v", q, args)
	}
	for _, k := range []string{"g.CEDING_CO_NAME", "TO_CHAR(g.START_DATE", "TO_CHAR(w.TGL_UPDATE"} {
		if !strings.Contains(q, k) {
			t.Errorf("kolom Beranda %s tidak dibaca:\n%s", k, q)
		}
	}
}

// Switch portal Proses / Resolved (keputusan work owner 06-10-2026, bawaan Proses): status ditukar, penampung tetap.
func TestSqlDaftarKasusSwitchSelesai(t *testing.T) {
	q, args := sqlDaftarKasus("S.W", "S.G", "S.Q", models.SaringanKasus{Pembuat: "UJI-A"})
	if !strings.Contains(q, "AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT IN (:2, :3))") || strings.Contains(q, "AND w.STATUS_WORK IN (") {
		t.Fatalf("bawaan Proses:\n%s", q)
	}
	q, args2 := sqlDaftarKasus("S.W", "S.G", "S.Q", models.SaringanKasus{Pembuat: "UJI-A", Selesai: true})
	if !strings.Contains(q, "AND w.STATUS_WORK IN (:2, :3)") || strings.Contains(q, "NOT IN (:2, :3)") ||
		!strings.Contains(q, "AND w.CREATE_OP = :4") || len(args2) != len(args) ||
		args2[1] != models.StatusDitolak || args2[2] != models.StatusSelesai {
		t.Fatalf("Resolved:\n%s\n%v", q, args2)
	}
}
