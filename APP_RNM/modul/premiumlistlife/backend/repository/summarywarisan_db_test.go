//go:build db

package repository_test

// PL-09 (GILIRAN-18) terhadap skema uji Oracle: rekap ke `M_LIFE_PREMIUM_SUMMARY`
// seperti `PEGA_M_LIFE_PREMIUM_SUMMARY` - satu baris per mata uang, ID dari
// sequence, uang kosong = 0, idempoten per nomor + work, di transaksi pemanggil.
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan.

import (
	"context"
	"testing"

	"github.com/cockroachdb/apd/v3"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/repository"
	"nusantarare/uji/skemauji"
)

// rekapUjiDB - IDR lengkap turunannya, USD tanpa satu pun jumlah (nol).
func rekapUjiDB() []models.RekapMataUang {
	return []models.RekapMataUang{
		{Currency: "IDR", Premium: apd.New(9750000, -4), Commission: apd.New(1250, -4),
			Balance: apd.New(-8970000, -4),
			Jumlah:  map[string]*apd.Decimal{"CLAIM": apd.New(10, 0), "DEDUCTION": apd.New(333, -2)}},
		{Currency: "USD", Premium: apd.New(1, 0), Commission: apd.New(0, 0), Balance: apd.New(1, 0),
			Jumlah: map[string]*apd.Decimal{}},
	}
}

func TestSummaryWarisanDitulisSepertiProsedur(t *testing.T) {
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	defer func() { _ = skemauji.Bongkar(ctx, sqlDB, skema) }()
	for _, q := range []string{
		`INSERT INTO ` + skema + `.T_WORK_POLIS (ID) VALUES ('UJI-POLIS-1801')`,
		`INSERT INTO ` + skema + `.T_PREMIUM_LIST (ID, BUSINESS_NAME) VALUES ('UJI-POLIS-1801', 'UJI-COB-1')`,
		// Baris endorsemen ber-PL_NUMBER sama milik work LAIN: tidak boleh terhapus.
		`INSERT INTO ` + skema + `.M_LIFE_PREMIUM_SUMMARY (ID, PL_NUMBER, IDPEGA, CURRENCY)
		   VALUES ('UJI-LAIN', 'UJI-PL-1801', 'UJI-WORK-LAIN', 'IDR')`,
	} {
		if _, err := sqlDB.ExecContext(ctx, q); err != nil {
			t.Fatalf("polis tiruan: %v", err)
		}
	}
	db, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	ringkas := repository.NewSummaryPolis(db)
	tulis := repository.NewSummaryWarisan(db)

	tulisSekali := func(nomor string, commit bool) (int, int) {
		t.Helper()
		tx, err := db.Mulai(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer func(tx *intidb.Tx) { _ = tx.Rollback() }(tx)
		kepala, err := ringkas.KepalaSummaryWarisan(ctx, tx, "UJI-POLIS-1801", nomor)
		if err != nil {
			t.Fatalf("kepala: %v", err)
		}
		h, s, err := tulis.Ganti(ctx, tx, kepala, rekapUjiDB())
		if err != nil {
			t.Fatalf("ganti summary warisan: %v", err)
		}
		if commit {
			if err := tx.Commit(); err != nil {
				t.Fatal(err)
			}
		}
		return h, s
	}
	hitung := func(where string, arg ...any) int {
		t.Helper()
		var n int
		if err := sqlDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+
			`.M_LIFE_PREMIUM_SUMMARY WHERE `+where, arg...).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}

	if h, s := tulisSekali("UJI-PL-1801", true); h != 0 || s != 2 {
		t.Errorf("putaran 1: terhapus %d, tersisip %d; mau 0, 2", h, s)
	}
	// Simpan ulang: idempoten - rekap lama work ini diganti, bukan ditumpuk.
	if h, s := tulisSekali("UJI-PL-1801", true); h != 2 || s != 2 {
		t.Errorf("putaran 2: terhapus %d, tersisip %d; mau 2, 2", h, s)
	}
	if n := hitung(`PL_NUMBER = :1 AND IDPEGA = :2`, "UJI-PL-1801", "UJI-POLIS-1801"); n != 2 {
		t.Errorf("baris work ini %d, mau 2 (satu per mata uang)", n)
	}
	if n := hitung(`ID = 'UJI-LAIN'`); n != 1 {
		t.Errorf("baris work lain ber-PL_NUMBER sama ikut terhapus")
	}
	if n := hitung(`IDPEGA = :1 AND CURRENCY = 'IDR' AND PREMIUM = 975 AND COMMISSION = 0.125
		AND BALANCE = -897 AND CLAIM = 10 AND DEDUCTION = 3.33 AND BROKERAGE_FEE = 0
		AND COB = 'UJI-COB-1' AND PL_NUMBER_EDM IS NULL AND ID IS NOT NULL`, "UJI-POLIS-1801"); n != 1 {
		t.Errorf("rekap IDR tidak tersimpan seperti prosedur: %d baris cocok, mau 1", n)
	}
	if n := hitung(`IDPEGA = :1 AND CURRENCY = 'USD' AND RI_ADMIN_FEE = 0 AND TAX = 0
		AND NET_PREMIUM_REFUND_RETRO = 0`, "UJI-POLIS-1801"); n != 1 {
		t.Errorf("uang kosong bukan 0 (PL-10): %d baris cocok, mau 1", n)
	}
	var beda int
	if err := sqlDB.QueryRowContext(ctx, `SELECT COUNT(DISTINCT ID) FROM `+skema+
		`.M_LIFE_PREMIUM_SUMMARY WHERE IDPEGA = 'UJI-POLIS-1801'`).Scan(&beda); err != nil {
		t.Fatal(err)
	}
	if beda != 2 {
		t.Errorf("ID dari sequence tidak unik per baris: %d", beda)
	}
	// Transaksi pemanggil yang dibatalkan membatalkan tulisannya: nol COMMIT.
	tulisSekali("UJI-PL-1802", false)
	if n := hitung(`PL_NUMBER = 'UJI-PL-1802'`); n != 0 {
		t.Errorf("tulisan summary warisan bertahan sesudah rollback: %d baris", n)
	}
}
