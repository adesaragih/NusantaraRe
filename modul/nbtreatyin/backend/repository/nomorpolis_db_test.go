//go:build db

package repository_test

// Spec AC 31: "Nomor polis terbentuk dari deret, tanpa bentrok. Test yang
// menemukan dua berkas bernomor sama gagal." - lawan Oracle SUNGGUHAN, lewat
// `TerbitkanNomorPolis` (= `GeneratePolicyNoTreaty_Act` langkah 5-28) dan
// penomor bersama `inti/backend/penomor` (padanan
// `PROC_GENERATE_SEQUENCE_NUMBER`, keputusan WO F8). ⛔ Belum pernah
// dijalankan: K11 kosong.
//
// Tiga tabel penomoran adalah WARISAN `POOLDATA`, tidak dibuat migrasi mana
// pun. Bila skema uji tidak memuatnya, dibuat TIRUAN lewat `buatTiruan` (U1),
// berbentuk katalog yang dicatat `modul/claimlife/docs/SUMBER-PENOMORAN-DBA.md`
// (dibaca dari `ALL_TAB_COLUMNS`). Bila DBA menyediakannya, baris penghitung
// kelas `KelasDeret` dipotret sebelum uji dan dipulihkan sesudahnya.
//
// Nilai harapan dari XML `GeneratePolicyNoTreaty_Act` langkah 28:
// awalan (`KODE_PRODUKSI`, TYPE NONLIFE) + tipe (DueTo "0" = QP, langkah 7-9)
// + ".T" + OJKBusinessID + "." + MM.YYYY + "." + urut lima angka.

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

const (
	// tiruan uji, bukan tabel aplikasi (bila DBA tidak menyediakannya):
	// penghitung deret warisan `POOLDATA.GENERATE_SEQUENCE_NUMBER`.
	tabelDeretUji = "GENERATE_SEQUENCE_NUMBER"
	// tiruan uji, bukan tabel aplikasi (bila DBA tidak menyediakannya): awalan
	// produksi warisan `POOLDATA.KODE_PRODUKSI`.
	tabelKodeProduksiUji = "KODE_PRODUKSI"
	// tiruan uji, bukan tabel aplikasi (bila DBA tidak menyediakannya): hari
	// tutup buku warisan `POOLDATA.TANGGAL_CLOSING`.
	tabelClosingUji = "TANGGAL_CLOSING"
)

// siapkanPenomoran memastikan ketiga tabel penomoran ada; mengembalikan true
// bila SELURUHNYA tiruan (nilai harapan dapat ditulis harfiah).
func siapkanPenomoran(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) bool {
	t.Helper()
	deret := buatTiruan(t, ctx, sqlDB, skema, tabelDeretUji, []string{"CLASS VARCHAR2(200) NOT NULL",
		"JENIS VARCHAR2(50) NOT NULL", "TAHUN VARCHAR2(5) NOT NULL", "NO_SEQ NUMBER", "TANGGAL DATE", "MM_YYYY VARCHAR2(10)"})
	kode := buatTiruan(t, ctx, sqlDB, skema, tabelKodeProduksiUji, []string{"KODE VARCHAR2(5)", "TYPE VARCHAR2(10)"})
	closing := buatTiruan(t, ctx, sqlDB, skema, tabelClosingUji, []string{"TANGGAL VARCHAR2(10)"})
	if kode {
		if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s.%s (KODE, TYPE) VALUES ('UJI-', :1)`, skema, tabelKodeProduksiUji), models.LiniKasus); err != nil {
			t.Fatal(err)
		}
	}
	if closing {
		if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s.%s (TANGGAL) VALUES ('25')`, skema, tabelClosingUji)); err != nil {
			t.Fatal(err)
		}
	}
	if !deret {
		potretDeret(t, ctx, sqlDB, skema)
	}
	return deret && kode && closing
}

// potretDeret menyimpan baris penghitung kelas NB Treaty In di tabel DBA dan
// memulihkannya sesudah uji.
func potretDeret(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) {
	t.Helper()
	rows, err := sqlDB.QueryContext(ctx, fmt.Sprintf(`SELECT JENIS, TAHUN, TO_CHAR(NO_SEQ), TO_CHAR(TANGGAL, 'YYYY-MM-DD HH24:MI:SS'), MM_YYYY
		FROM %s.%s WHERE CLASS = :1`, skema, tabelDeretUji), models.KelasDeret)
	if err != nil {
		t.Fatal(err)
	}
	var simpan [][5]sql.NullString
	for rows.Next() {
		var r [5]sql.NullString
		if err := rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4]); err != nil {
			t.Fatal(err)
		}
		simpan = append(simpan, r)
	}
	rows.Close()
	t.Cleanup(func() {
		_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s.%s WHERE CLASS = :1`, skema, tabelDeretUji), models.KelasDeret)
		for _, r := range simpan {
			_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s.%s (CLASS, JENIS, TAHUN, NO_SEQ, TANGGAL, MM_YYYY)
				VALUES (:1, :2, :3, TO_NUMBER(:4), TO_DATE(:5, 'YYYY-MM-DD HH24:MI:SS'), :6)`, skema, tabelDeretUji),
				models.KelasDeret, r[0], r[1], r[2], r[3], r[4])
		}
	})
}

var polaNomorUji = regexp.MustCompile(`^(.+QP\.TUJI1\.\d{2}\.\d{4}\.)(\d{5})$`)

func TestNomorPolisDariDeretTanpaBentrok(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	semuaTiruan := siapkanPenomoran(t, ctx, sqlDB, skema)
	g := repository.Baru(d)
	h := models.HalamanBaru()
	h.Setel(models.HalamanPolis+".DueTo", "0") // langkah 7-9: QP
	h.Setel(models.HalamanPolis+".OJKBusinessID", "UJI1")
	sekarang := time.Date(2026, 10, 5, 10, 0, 0, 0, time.FixedZone("WIB", 7*3600)) // hari 5 <= tutup buku 25

	terbitkan := func(tx *intidb.Tx) string {
		b, err := g.TerbitkanNomorPolis(ctx, tx, h, sekarang)
		if err != nil {
			t.Error(err)
			return ""
		}
		return b.NoPolis
	}
	var n1 string
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { n1 = terbitkan(tx); return nil }); err != nil {
		t.Fatal(err)
	}
	// Dua transaksi SERENTAK: B harus MENUNGGU kunci baris penghitung yang
	// dipegang A (`FOR UPDATE`), lalu mendapat urut berikutnya - bukan urut
	// yang sama.
	txA, err := d.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = txA.Rollback() }()
	nA := terbitkan(txA)
	hasilB := make(chan string, 1)
	go func() {
		txB, err := d.Mulai(ctx)
		if err != nil {
			hasilB <- "galat: " + err.Error()
			return
		}
		defer func() { _ = txB.Rollback() }()
		b, err := g.TerbitkanNomorPolis(ctx, txB, h, sekarang)
		if err != nil {
			hasilB <- "galat: " + err.Error()
			return
		}
		if err := txB.Commit(); err != nil {
			hasilB <- "galat: " + err.Error()
			return
		}
		hasilB <- b.NoPolis
	}()
	select {
	case nB := <-hasilB:
		t.Fatalf("transaksi B selesai (%s) tanpa menunggu kunci penghitung transaksi A (%s) - nomor dapat bentrok", nB, nA)
	case <-time.After(3 * time.Second):
	}
	if err := txA.Commit(); err != nil {
		t.Fatal(err)
	}
	nB := <-hasilB

	urut := func(n string) (string, int) {
		m := polaNomorUji.FindStringSubmatch(n)
		if m == nil {
			t.Fatalf("nomor %q tidak berpola langkah 28 (awalan+QP.T<OJK>.MM.YYYY.urut5)", n)
		}
		u, _ := strconv.Atoi(m[2])
		return m[1], u
	}
	b1, u1 := urut(n1)
	bA, uA := urut(nA)
	bB, uB := urut(nB)
	if b1 != bA || bA != bB || uA != u1+1 || uB != uA+1 {
		t.Fatalf("deret tidak berurut tanpa celah/kembar: %s, %s, %s", n1, nA, nB)
	}
	if semuaTiruan {
		for i, n := range []string{n1, nA, nB} {
			if harap := fmt.Sprintf("UJI-QP.TUJI1.10.2026.%05d", i+1); n != harap {
				t.Errorf("nomor ke-%d = %s, harap %s", i+1, n, harap)
			}
		}
		if seq := kolomTeks(t, ctx, sqlDB, fmt.Sprintf(`SELECT TO_CHAR(NO_SEQ) FROM %s.%s WHERE CLASS = :1 AND JENIS = :2 AND TAHUN = '2026'`,
			skema, tabelDeretUji), models.KelasDeret, models.JenisDeret("UJI-")); seq != "3" {
			t.Errorf("NO_SEQ penghitung = %s, harap 3", seq)
		}
	}
	// Ketiga nomor tersimpan di tiga berkas berbeda - tidak satu pun ditolak.
	for i, n := range []string{n1, nA, nB} {
		id := fmt.Sprintf("UJI-NB-DERET-%d", i+1)
		if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
			if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI"); err != nil {
				return err
			}
			return g.SetelNomorPolis(ctx, tx, id, n)
		}); err != nil {
			t.Fatalf("berkas %s nomor %s: %v", id, n, err)
		}
	}
}
