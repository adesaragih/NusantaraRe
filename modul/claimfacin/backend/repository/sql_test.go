package repository

// Uji teks SQL tanpa Oracle: teks cari Choose Polis DIIKAT (perbaikan prompt §5 butir 1 - XML menyisipkan
// `{ASIS:InputData.CARI1}`), nol COMMIT dan nol `%[n]s` di setiap SQL modul ini (prompt §3, ADR-U-0033), UPDATE OS DLA
// menurut nomor akseptasi (InsertDLA_OS_SQL).

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
)

func TestCariPolisTeksTerikat(t *testing.T) {
	for _, persis := range []bool{true, false} {
		q := sqlCariPolis("UJI.FP", "UJI.BIS", "UJI.JP", "NOPOLIS", persis)
		if !strings.Contains(q, ":1") || strings.Contains(q, "ASIS") || strings.Contains(q, "{") {
			t.Fatalf("teks cari tidak terikat: %s", q)
		}
		if !strings.Contains(q, "FETCH FIRST 500 ROWS ONLY") {
			t.Fatalf("tanpa batas baris: %s", q)
		}
	}
	if q := sqlCariPolis("UJI.FP", "UJI.BIS", "UJI.JP", "CEDINGCO", false); !strings.Contains(q, "a.CEDINGCO LIKE :1 || '%'") {
		t.Fatalf("awalan nama ceding: %s", q)
	}
}

func TestSQLModulTanpaCommitDanIndeksFormat(t *testing.T) {
	tabel := &models.TabelAdjustment
	baca, _ := sqlBacaSimpul(tabel, "UJI.T", "UJI.ADJ")
	semua := []string{
		sqlCariPolis("A", "B", "C", "NOPOLIS", true), sqlSisipKasus("W"), sqlSisipInduk("G"), sqlKeadaan("W", "G", true),
		sqlPindahTahap("W"), sqlTutupKasus("W"), sqlSentuhKasus("W"), baca, sqlSisipSimpul("T", tabel),
		sqlUbahSimpul("T", tabel), sqlNomorKronologi("V"), sqlSisipKronologi("V"), sqlBacaKronologi("V"),
		sqlSisipKasusKomite("W"), sqlSisipKepalaKomite("G"), sqlSisipAnggotaKomite("L"), sqlTanggaKomite("L", "G", "W"),
		sqlSetelKomiteAdjustment("A"), sqlKategoriLampiran("K", "D"), sqlDaftarLampiran("D"), sqlSisipDokumenKlaim("D"),
		sqlHapusDokumenKlaim("D"), sqlPindahKategoriDokumen("D"), sqlSisipOS("O"), sqlTandaiDLAOS("O"), sqlBacaOS("O"),
		sqlAdaJSONKlaim("J"), sqlSisipJSONKlaim("J"), sqlSisipKatastrofe("C"), sqlLogLayanan("M"), sqlSisipProgres("P"),
		sqlSelesaiProgres("P"), sqlSisipSubProgres("S"), sqlGeserNourut("N"),
	}
	for _, q := range semua {
		u := strings.ToUpper(q)
		if strings.Contains(u, "COMMIT") || strings.Contains(q, "%[") || strings.Contains(u, "BEGIN ") {
			t.Errorf("SQL terlarang: %s", q)
		}
	}
}

func TestTandaiDLAOSMenurutNomorAkseptasi(t *testing.T) {
	q := sqlTandaiDLAOS("UJI.OS")
	if !strings.Contains(q, "SET DLA_NO = :1, DLA_DATE = :2") ||
		!strings.Contains(q, "JSON_VALUE(DATA_JSON, '$.AcceptedNo') = :3") {
		t.Fatalf("%s", q)
	}
}

// Setiap tabel stabil ber-UNIQUE (CLAIM_ID, NOURUT) digeser dua fase sebelum ditulis (temuan review: item baru di objek
// pertama bertabrakan dengan item objek kedua).
func TestGeserNourutSemuaTabelStabil(t *testing.T) {
	ada := map[string]bool{}
	for _, n := range tabelGeser() {
		ada[n] = true
	}
	for _, n := range []string{"T_CLAIM_OBJECT", "T_CLAIM_OBJECT_ITEM", "T_CLAIM_ADJUSTMENT"} {
		if !ada[n] {
			t.Errorf("%s tidak digeser: %v", n, tabelGeser())
		}
	}
	if q := sqlGeserNourut("UJI.T"); !strings.Contains(q, "SET NOURUT = -NOURUT WHERE CLAIM_ID = :1 AND NOURUT > 0") {
		t.Fatal(q)
	}
}

// Kolom milik komite (AcceptedNo / AcceptedDate / AcceptanceStatus / IsApproved) tidak ikut UPDATE simpan halaman klaim:
// simpan klaim yang membaca halaman sebelum komite memutuskan tidak menimpa keputusannya (temuan review 10-10-2026).
func TestUbahAdjustmentTanpaKolomKomite(t *testing.T) {
	q := sqlUbahSimpul("UJI.ADJ", &models.TabelAdjustment)
	for _, k := range []string{"ACCEPTED_NO", "ACCEPTED_DATE", "ACCEPTANCE_STATUS", "IS_APPROVED", "NOTES"} {
		if strings.Contains(q, k+" =") {
			t.Errorf("%s ikut UPDATE: %s", k, q)
		}
	}
	if !strings.Contains(q, "IS_PRINT_ACCEPT =") && !strings.Contains(q, "STATUS_KASIR =") {
		t.Fatalf("kolom klaim hilang: %s", q)
	}
	b := models.Baris{}
	semua, err := argsKatalog(&models.TabelAdjustment, b, false)
	if err != nil {
		t.Fatal(err)
	}
	ubah, err := argsKatalog(&models.TabelAdjustment, b, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(semua)-len(ubah) < 4 || len(regexp.MustCompile(`:\d+`).FindAllString(q, -1)) != len(ubah)+4 {
		t.Fatalf("jumlah bind: semua %d ubah %d, teks %s", len(semua), len(ubah), q)
	}
}

// TestUbahAdjustmentKomiteHanyaKolomKomite - penulis keputusan komite (kontrak KlaimFacInKomite) hanya menyentuh kolom
// milik komite, terikat ID + CLAIM_ID, bind urut kemunculan.
func TestUbahAdjustmentKomiteHanyaKolomKomite(t *testing.T) {
	var kolom []models.Kolom
	for _, p := range []string{"AcceptanceStatus", "AcceptedDate", "AcceptedNo", "IsApproved", "Notes"} {
		k, ok := kolomKomiteAdjustment(p)
		if !ok {
			t.Fatalf("%s bukan kolom milik komite", p)
		}
		kolom = append(kolom, k)
	}
	for _, p := range []string{"PaymentType", "IsPrintAccept", "StatusKasir", "IDOfBank", "IsFacRetro"} {
		if _, ok := kolomKomiteAdjustment(p); ok {
			t.Fatalf("%s dianggap kolom milik komite (harus lewat simpan halaman)", p)
		}
	}
	q := sqlUbahAdjustmentKomite("UJI.ADJ", kolom)
	want := "UPDATE UJI.ADJ SET ACCEPTANCE_STATUS = :1, ACCEPTED_DATE = TO_DATE(:2, 'YYYY-MM-DD HH24:MI:SS'), " +
		"ACCEPTED_NO = :3, IS_APPROVED = :4, NOTES = :5 WHERE ID = :6 AND CLAIM_ID = :7"
	if q != want {
		t.Fatalf("SQL got %s, want %s", q, want)
	}
}
