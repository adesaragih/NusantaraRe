//go:build db

package repository_test

// Uji seam repository lawan Oracle SUNGGUHAN untuk PEMBACAAN VIEW KONTRAK
// `TREATYINDETAILJOINEDM` (spec AC 17, 89; tiket 01). ⛔ Belum pernah
// dijalankan: skema uji K11 kosong (PERMINTAAN C4).
//
// View itu WARISAN `POOLDATA` (TREATYINDETAIL UNION ALL TREATYINDETAILEDM),
// tidak dibuat migrasi modul ini. Dua keadaan skema uji:
//
//   - view SUNGGUHAN disediakan DBA di skema uji -> AC 89 diperiksa lawan
//     katalognya, AC 17 lawan baris pertama yang ke-33 kolom RD-nya terisi;
//   - tidak ada -> AC 17 memakai TIRUAN berbentuk katalog (fakta
//     `ALL_TAB_COLUMNS` 03-10-2026, PROMPT putaran 2 bab 1: kolom nilai
//     `NUMBER` tanpa presisi/skala, `COMMENCEMENT`/`TERMINATION` `DATE`,
//     `INSTALLMENTNO` `VARCHAR2(1000)`, kolom lain `VARCHAR2(1000)`, `ID`
//     `VARCHAR2(100)`) yang dibuang lagi; AC 89 MELEWATI - ia menuntut view
//     sungguhan, dan tiruan buatan uji ini tidak membuktikan apa pun tentangnya.
//
// Pemanggil `skemauji.Buka()` tetap SATU untuk modul ini (`pasang`,
// polis_db_test.go) - penjaga claimlife `TestSetiapPemanggilBukaMemeriksaBolehDilewati`
// menghitungnya. Fixture berawalan UJI-.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

const viewKontrak = "TREATYINDETAILJOINEDM"

// kolomViewDibaca - setiap kolom view yang DIBACA modul: 33 kolom RD
// `BrowseTreatyJoinEDM` (`repository.KolomRDDetail`), `COMMENCEMENT` /
// `TERMINATION` (halaman master, `DetailKontrak`), dan `RIOGR` / `RIONR`
// (`KomisiKontrak`, `TreatyInputPctCommSpreading` 2.1.1.1) - 37 dari 39 kolom
// view; `RNM_SHARE` dan `BROKERAGE` dibaca NOL rule terjangkau (spec AC 26).
func kolomViewDibaca() []string {
	k := append([]string{}, repository.KolomRDDetail...)
	return append(k, "COMMENCEMENT", "TERMINATION", models.KolomRIOGR, models.KolomRIONR)
}

// kolomViewAngka - kolom bertipe NUMBER menurut katalog (PROMPT putaran 2 bab 1).
var kolomViewAngka = map[string]bool{
	"LIMITVALUE": true, "RETENTIONVALUE": true, "EPIVALUE": true, "NETPREMIVALUE": true,
	"SHAREVALUE": true, "MDPVALUE": true, "DEDUCTION1": true, "DEDUCTION2": true,
	"RIOGR": true, "RIONR": true, "RNM_SHARE": true, "BROKERAGE": true,
}

// jenisObjekView - jenis objek `TREATYINDETAILJOINEDM` di skema uji ("" = tidak ada).
func jenisObjekView(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) string {
	t.Helper()
	var jenis sql.NullString
	err := sqlDB.QueryRowContext(ctx, `SELECT MAX(OBJECT_TYPE) FROM ALL_OBJECTS WHERE OWNER = UPPER(:1) AND OBJECT_NAME = :2`,
		skema, viewKontrak).Scan(&jenis)
	if err != nil {
		t.Fatal(err)
	}
	return jenis.String
}

// pasangTiruanView membuat tiruan berbentuk katalog (39 kolom) dan membuangnya
// sesudah uji. Dipanggil HANYA bila objeknya tidak ada.
func pasangTiruanView(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) {
	t.Helper()
	var kolom []string
	for _, k := range append(kolomViewDibaca(), "RNM_SHARE", "BROKERAGE") {
		switch {
		case k == "ID":
			kolom = append(kolom, "ID VARCHAR2(100)")
		case k == "COMMENCEMENT" || k == "TERMINATION":
			kolom = append(kolom, k+" DATE")
		case kolomViewAngka[k]:
			kolom = append(kolom, k+" NUMBER")
		default:
			kolom = append(kolom, k+" VARCHAR2(1000)")
		}
	}
	if len(kolom) != 39 { // spec §5.1: view 39 kolom
		t.Fatalf("tiruan %d kolom, view 39", len(kolom))
	}
	q := fmt.Sprintf(`CREATE TABLE %s.%s (%s)`, skema, viewKontrak, strings.Join(kolom, ", "))
	if _, err := sqlDB.ExecContext(ctx, q); err != nil {
		t.Fatalf("tiruan %s: %v", viewKontrak, err)
	}
	t.Cleanup(func() { _, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`DROP TABLE %s.%s PURGE`, skema, viewKontrak)) })
}

// AC 89: "Modul ini membaca 39 kolom view; nol medan yang dipakai tidak
// tersedia. Test yang menemukan medan hilang gagal."
func TestViewKontrakMemuatSetiapKolomYangDibaca(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	if jenisObjekView(t, ctx, sqlDB, skema) == "" {
		t.Skip("lewati: view warisan TREATYINDETAILJOINEDM tidak ada di skema uji (PERMINTAAN C4) - " +
			"AC 89 menuntut view SUNGGUHAN, bukan tiruan")
	}
	ada := map[string]bool{}
	rows, err := sqlDB.QueryContext(ctx, `SELECT COLUMN_NAME FROM ALL_TAB_COLUMNS WHERE OWNER = UPPER(:1) AND TABLE_NAME = :2`, skema, viewKontrak)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			t.Fatal(err)
		}
		ada[k] = true
	}
	rows.Close()
	for _, k := range kolomViewDibaca() {
		if !ada[k] {
			t.Errorf("kolom %s dibaca modul, tidak ada di view", k)
		}
	}
	// Jalur baca yang sebenarnya: `pilihKolom` menolak kolom yang hilang
	// (ErrDataKontrakTidakAda "kolom ... tidak ada"); ID yang tidak ada
	// menghasilkan galat ID, BUKAN galat kolom.
	g := repository.Baru(d)
	_, err = g.DetailKontrak(ctx, "UJI-KONTRAK-TIDAK-ADA")
	if !errors.Is(err, repository.ErrDataKontrakTidakAda) || strings.Contains(err.Error(), "kolom") {
		t.Fatalf("DetailKontrak atas view sungguhan: %v", err)
	}
}

// AC 17: "Ke-33 medan yang dipakai laporan tersedia dari view. Test yang
// menemukan medan laporan tidak terisi gagal."
func TestDetailKontrakMengisiKe33MedanRD(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	if jenisObjekView(t, ctx, sqlDB, skema) != "" {
		// view sungguhan: baris pertama yang ke-33 kolom RD-nya terisi.
		var syarat []string
		for _, k := range repository.KolomRDDetail {
			syarat = append(syarat, k+" IS NOT NULL")
		}
		var id string
		err := sqlDB.QueryRowContext(ctx, fmt.Sprintf(`SELECT ID FROM %s.%s WHERE %s FETCH FIRST 1 ROWS ONLY`,
			skema, viewKontrak, strings.Join(syarat, " AND "))).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			t.Skip("lewati: view sungguhan tanpa baris yang ke-33 kolom RD-nya terisi")
		}
		if err != nil {
			t.Fatal(err)
		}
		b, err := g.DetailKontrak(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		for _, k := range repository.KolomRDDetail {
			if b[k] == "" {
				t.Errorf("%s terisi di view, kosong di hasil baca", k)
			}
		}
		return
	}
	pasangTiruanView(t, ctx, sqlDB, skema)
	// Satu baris UJI- dengan nilai yang dihitung tangan. Angka ditulis sebagai
	// literal SQL (titik desimal tidak bergantung NLS sesi), tanggal sebagai
	// literal DATE.
	harap := map[string]string{}
	var kolom, nilai []string
	for i, k := range append(kolomViewDibaca(), "RNM_SHARE", "BROKERAGE") {
		kolom = append(kolom, k)
		switch {
		case k == "ID":
			harap[k] = "UJI-KONTRAK-1"
			nilai = append(nilai, "'UJI-KONTRAK-1'")
		case k == "COMMENCEMENT":
			harap[k] = "2026-01-01 00:00:00"
			nilai = append(nilai, "DATE '2026-01-01'")
		case k == "TERMINATION":
			harap[k] = "2026-12-31 00:00:00"
			nilai = append(nilai, "DATE '2026-12-31'")
		case k == "INSTALLMENTNO":
			harap[k] = "4"
			nilai = append(nilai, "'4'")
		case kolomViewAngka[k]:
			// 1500000000.5, 1500000001.5, ... - setengah supaya pemisah desimal
			// ikut diperiksa; DEDUCTION2 0.5 memeriksa ".5" -> "0.5".
			v := fmt.Sprintf("%d.5", 1500000000+i)
			if k == "DEDUCTION2" {
				v = "0.5"
			}
			harap[k] = v
			nilai = append(nilai, v)
		default:
			harap[k] = "UJI-" + k
			nilai = append(nilai, "'UJI-"+k+"'")
		}
	}
	if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s.%s (%s) VALUES (%s)`,
		skema, viewKontrak, strings.Join(kolom, ", "), strings.Join(nilai, ", "))); err != nil {
		t.Fatal(err)
	}
	b, err := g.DetailKontrak(ctx, "UJI-KONTRAK-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range append(append([]string{}, repository.KolomRDDetail...), "COMMENCEMENT", "TERMINATION") {
		if b[k] != harap[k] {
			t.Errorf("%s = %q, harap %q", k, b[k], harap[k])
		}
	}
	if _, ada := b["RNM_SHARE"]; ada {
		t.Error("RNM_SHARE dibaca padahal nol rule terjangkau memakainya (AC 26)")
	}
	// RIOGR/RIONR lewat pembaca komisi (TREATYID = NoOffer).
	k, err := g.KomisiKontrak(ctx, "UJI-TREATYID")
	if err != nil {
		t.Fatal(err)
	}
	if len(k) != 1 || k[0][models.KolomRIOGR] != harap["RIOGR"] || k[0][models.KolomRIONR] != harap["RIONR"] {
		t.Fatalf("komisi kontrak %+v", k)
	}
}
