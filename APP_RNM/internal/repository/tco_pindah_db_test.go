//go:build db

// Seam `repository` terhadap skema uji Oracle NYATA - tiket 01 Treaty
// Contract Out (tco2).
//
// Untuk apa berkas ini: membuktikan bahwa migrasi 300-306 berdiri, bahwa
// enam tabel warisan tiruan berpindah ke T_* dalam SATU transaksi tanpa
// mengubah satu digit pun, bahwa temuan yang memblokir membatalkan
// SELURUHNYA, dan bahwa sequence diselaraskan ke identitas terbesar + 1.
//
// Jalankan: make test-db   (perlu ORACLE_DSN, ORACLE_SCHEMA, ORACLE_SKEMA_UJI)
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan - bukan
// lulus diam-diam. Fixture: nol nama orang, seluruh nilai berawalan UJI.
package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
)

func siapkanTCO(t *testing.T) (*sql.DB, string, *repository.DB, func()) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	db, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	return sqlDB, skema, db, func() {
		_ = skemauji.Bongkar(ctx, sqlDB, skema)
		_ = db.Close()
		_ = sqlDB.Close()
	}
}

// isiWarisanContoh mengisi enam tiruan warisan dengan data yang memuat
// jebakan nyata: koma desimal, delapan angka di belakang koma, stempel Pega,
// CHAR bertabur spasi, TREATYYEARID NULL, kolom mati berisi.
func isiWarisanContoh(t *testing.T, sqlDB *sql.DB, skema string) {
	t.Helper()
	ctx := context.Background()
	isi := func(tabel string, baris ...map[string]string) {
		if err := skemauji.IsiWarisanTCO(ctx, sqlDB, skema, tabel, baris); err != nil {
			t.Fatalf("mengisi %s: %v", tabel, err)
		}
	}
	isi("TREATYYEAR", map[string]string{
		"ID": "1000001", "TREATYYEAR": "2026", "UNDERWRITINGYEAR": "2026", "TREATYGROUPID": "10001",
		"TREATYGROUPNAME": "UJI GRUP", "USERID": "UJI-OP", "TGLUPDATE": "15-SEP-26",
		"PROPORTION": "P", "STARTDATE": "20260101", "ENDDATE": "20261231T170000.000 GMT",
	}, map[string]string{
		"ID": "1000005", "TREATYYEAR": "2025", "TREATYGROUPID": "10001", "STARTDATE": "20250101", "ENDDATE": "20251231",
	})
	isi("TREATYCONTRACT", map[string]string{
		"ID": "1000007", "IDTREATYYEAR": "1000001", "REINSTYPEID": "10002", "REINSTYPENAME": "UJI QS",
		"TREATYSTARTDATE": "2026-01-01 00:00:00", "TREATYENDDATE": "2026-12-31 00:00:00",
		"USERID": "UJI-OP", "TGLUPDATE": "01/02/2026",
	})
	isi("TREATYREINSURER", map[string]string{
		"ID": "1000003", "TREATYYEAR": "2026", "TREATYGROUPID": "10001", "REINSTYPEID": "10002",
		"REINSURERID": "R1", "CLIENTID": "C1", "NAME": "UJI REINS A", "RICOMM": "12.5",
		"PCTSHARE": "33.33333333", "IUDATE": "UJI-APA-ADANYA", "STARTDATE": "20260101T170000.000 GMT",
		"OPERATORNAME": "UJI-OP",
	}, map[string]string{
		"ID": "1000004", "TREATYYEAR": "2026", "TREATYGROUPID": "10001", "REINSTYPEID": "10002",
		"REINSURERID": "R2", "NAME": "UJI REINS B", "PCTSHARE": "66.66666667",
	})
	isi("MTREATYSECURITY", map[string]string{
		"THN_TREATY": "2026", "REAS_ID": "1000003", "PCT_SHARE": "50,5", "REAS_SECURITY": "SEC-A",
	}, map[string]string{
		"THN_TREATY": "2026", "REAS_ID": "1000003", "PCT_SHARE": "49.5", "REAS_SECURITY": "SEC-B",
	}, map[string]string{
		"THN_TREATY": "2026", "REAS_ID": "1000004", "PCT_SHARE": "100", "REAS_SECURITY": "SEC-A",
	})
	isi("TREATYBUSINESS", map[string]string{
		"ID": "1000002", "ISACTIVE": "1", "TREATYYEAR": "2026", "TREATYYEARID": "",
		"TREATYGROUPID": "10001", "REINSTYPEID": "10002", "BIZCODE": "UJI-BIZ", "BIZNAME": "UJI BISNIS",
	})
	isi("PROPORTIONALARRG", map[string]string{
		"ID": "10000009", "TREATYYEAR": "2026", "TREATYYEARID": "1000001", "TREATYGROUPID": "10001",
		"TREATYDESCID": "10001", "TREATYDESCNAME": "UJI LIMIT", "REINSTYPEID": "10002",
		"PARENTREINSTYPEID": "00", "RP": "1000000000,12345678", "USD": "0.5", "PCT": "100",
		"TREATYLIMIT": "5000000", "COINS_MIN": "0.1", "COINS_MAX": "0.9", "MORERP": "1", "MOREUSD": "2",
		"TGLUPDATE": "2026-01-02 03:04:05", "OBJECT": "UJI-MATI",
	}, map[string]string{
		"ID": "10000010", "TREATYYEAR": "2026", "TREATYYEARID": "1000001", "TREATYGROUPID": "10001",
		"TREATYDESCID": "10001", "REINSTYPEID": "10002", "PARENTREINSTYPEID": "10002",
		"PCT": "12,5", "RP": "125000000.01543210", "USD": ".0625",
	})
}

func cacahBaris(t *testing.T, sqlDB *sql.DB, skema, tabel string) int {
	t.Helper()
	var n int
	if err := sqlDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", skema, tabel)).Scan(&n); err != nil {
		t.Fatalf("mencacah %s: %v", tabel, err)
	}
	return n
}

func TestPindahkanTCOSatuTransaksiTanpaKehilanganDigit(t *testing.T) {
	sqlDB, skema, db, bersihkan := siapkanTCO(t)
	defer bersihkan()
	ctx := context.Background()
	isiWarisanContoh(t, sqlDB, skema)

	lap, err := repository.NewMigrasiTCO(db).Pindahkan(ctx)
	if err != nil {
		t.Fatalf("pindahkan: %v\n%s", err, lap)
	}
	if lap.Memblokir() || len(lap.Selisih) != 0 {
		t.Fatalf("laporan tidak bersih:\n%s", lap)
	}
	mau := map[string]int{
		"T_TREATYYEAR": 2, "T_TREATYCONTRACT": 1, "T_TREATYREINSURER": 2,
		"T_MTREATYSECURITY": 3, "T_TREATYBUSINESS": 1, "T_PROPORTIONALARRG": 2,
	}
	for tabel, n := range mau {
		if lap.Ditulis[tabel] != n || lap.DibacaKembali[tabel] != n {
			t.Errorf("%s: ditulis %d, dibaca kembali %d, mau %d", tabel, lap.Ditulis[tabel], lap.DibacaKembali[tabel], n)
		}
		if dapat := cacahBaris(t, sqlDB, skema, tabel); dapat != n {
			t.Errorf("%s memuat %d baris sesudah commit, mau %d", tabel, dapat, n)
		}
	}
	// Kolom mati berisi dilaporkan, tidak memblokir (AC 70).
	if lap.KolomMatiBerisi["OBJECT"] != 1 || lap.KolomMatiBerisi["PROPORTIONALLIST"] != 0 {
		t.Errorf("kolom mati: %v", lap.KolomMatiBerisi)
	}

	// Digit demi digit lewat kontrak hilir (kolom VERBATIM).
	hilir := repository.NewKontrakHilirTCO(db)
	klausul, err := hilir.KlausulUntukHilir(ctx, "10001", "2026", "10001", "10002")
	if err != nil {
		t.Fatal(err)
	}
	if len(klausul) != 2 {
		t.Fatalf("klausul hilir %d, mau 2", len(klausul))
	}
	if klausul[0].Rp != "1000000000.12345678" || klausul[0].Usd != ".5" && klausul[0].Usd != "0.5" {
		t.Errorf("RP/USD baris induk: %q %q", klausul[0].Rp, klausul[0].Usd)
	}
	if klausul[1].Rp != "125000000.0154321" && klausul[1].Rp != "125000000.01543210" {
		t.Errorf("RP baris anak berubah: %q", klausul[1].Rp)
	}
	if klausul[1].Pct != "12.5" {
		t.Errorf("PCT koma desimal: %q", klausul[1].Pct)
	}
	reins, err := hilir.ReinsurerUntukHilir(ctx, "10002", "2026", "10001")
	if err != nil {
		t.Fatal(err)
	}
	if len(reins) != 2 || reins[0].PctShare != "33.33333333" || reins[0].Name != "UJI REINS A" {
		t.Errorf("reinsurer hilir: %+v", reins)
	}
	grup, err := hilir.GrupTreatyAktifUntukHilir(ctx, "UJI-BIZ", "2026")
	if err != nil || len(grup) != 1 || grup[0] != "10001" {
		t.Errorf("grup aktif: %v %v", grup, err)
	}

	// Security: PK surrogat, spasi ekor CHAR dibuang, FK utuh.
	var reasSec string
	var pctShare string
	if err := sqlDB.QueryRow(fmt.Sprintf(`SELECT REAS_SECURITY, TO_CHAR(PCT_SHARE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')
		FROM %s.T_MTREATYSECURITY WHERE REAS_ID = '1000003' AND PCT_SHARE = 50.5`, skema)).Scan(&reasSec, &pctShare); err != nil {
		t.Fatalf("membaca security: %v", err)
	}
	if reasSec != "SEC-A" || pctShare != "50.5" {
		t.Errorf("security: %q %q", reasSec, pctShare)
	}
	var idKosong int
	_ = sqlDB.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s.T_MTREATYSECURITY WHERE ID IS NULL OR LENGTH(ID) <> 7`, skema)).Scan(&idKosong)
	if idKosong != 0 {
		t.Errorf("%d baris security tanpa surrogate '1'+6 digit", idKosong)
	}

	// Tanggal: stempel Pega dan YYYYMMDD tiba sebagai DATE tanpa pergeseran.
	var mulai, akhir string
	if err := sqlDB.QueryRow(fmt.Sprintf(`SELECT TO_CHAR(STARTDATE,'YYYY-MM-DD HH24:MI:SS'), TO_CHAR(ENDDATE,'YYYY-MM-DD HH24:MI:SS')
		FROM %s.T_TREATYYEAR WHERE ID = '1000001'`, skema)).Scan(&mulai, &akhir); err != nil {
		t.Fatal(err)
	}
	if mulai != "2026-01-01 00:00:00" || akhir != "2026-12-31 00:00:00" {
		t.Errorf("tanggal tahun: %q %q", mulai, akhir)
	}

	// Sequence diselaraskan: tahun -> 6 (ekor terbesar 5), klausul -> 11.
	if lap.SequenceBerikut["SEQ_T_TREATYYEAR"] != 6 || lap.SequenceBerikut["SEQ_T_PROPORTIONALARRG"] != 11 {
		t.Errorf("sequence berikutnya: %v", lap.SequenceBerikut)
	}
	tx, err := db.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	idBaru, err := db.IdentitasBerikutTCO(ctx, tx, "SEQ_T_TREATYYEAR")
	if err != nil {
		t.Fatal(err)
	}
	if idBaru != "1000006" {
		t.Errorf("identitas tahun berikutnya %q, mau 1000006 - sequence tidak diselaraskan", idBaru)
	}

	// Diulang: ditolak, tidak menggandakan.
	_, err = repository.NewMigrasiTCO(db).Pindahkan(ctx)
	if !errors.Is(err, repository.ErrTabelBaruTCOSudahBerisi) {
		t.Errorf("pemindahan kedua harus ditolak ErrTabelBaruTCOSudahBerisi, dapat %v", err)
	}
	if cacahBaris(t, sqlDB, skema, "T_TREATYYEAR") != 2 {
		t.Error("pemindahan kedua menggandakan baris")
	}
}

func TestPindahkanTCOTemuanMemblokirMembatalkanSeluruhnya(t *testing.T) {
	sqlDB, skema, db, bersihkan := siapkanTCO(t)
	defer bersihkan()
	ctx := context.Background()
	isiWarisanContoh(t, sqlDB, skema)
	// Satu nilai rusak di tabel TERAKHIR yang ditulis: bila transaksinya tidak
	// satu, lima tabel sebelumnya sudah terlanjur berisi.
	if err := skemauji.IsiWarisanTCO(ctx, sqlDB, skema, "PROPORTIONALARRG", []map[string]string{{
		"ID": "10000011", "TREATYYEAR": "2026", "TREATYGROUPID": "10001", "REINSTYPEID": "10002",
		"RP": "1.000.000", "PCT": "0.123456789",
	}}); err != nil {
		t.Fatal(err)
	}

	lap, err := repository.NewMigrasiTCO(db).Pindahkan(ctx)
	if !errors.Is(err, repository.ErrMigrasiTCOTidakBersih) {
		t.Fatalf("harus dibatalkan ErrMigrasiTCOTidakBersih, dapat %v\n%s", err, lap)
	}
	jenis := map[string]int{}
	for _, tm := range lap.Temuan {
		jenis[tm.Jenis]++
	}
	if jenis[repository.TemuanUangTakTerurai] != 1 || jenis[repository.TemuanPresisiMelampaui] != 1 {
		t.Errorf("temuan: %v", jenis)
	}
	for _, tabel := range []string{"T_TREATYYEAR", "T_TREATYCONTRACT", "T_TREATYREINSURER",
		"T_MTREATYSECURITY", "T_TREATYBUSINESS", "T_PROPORTIONALARRG"} {
		if n := cacahBaris(t, sqlDB, skema, tabel); n != 0 {
			t.Errorf("%s memuat %d baris padahal pemindahan dibatalkan - bukan satu transaksi", tabel, n)
		}
	}
}

func TestPindahkanTCORujukanYatimDilaporkanBukanORA(t *testing.T) {
	sqlDB, skema, db, bersihkan := siapkanTCO(t)
	defer bersihkan()
	ctx := context.Background()
	isiWarisanContoh(t, sqlDB, skema)
	if err := skemauji.IsiWarisanTCO(ctx, sqlDB, skema, "MTREATYSECURITY", []map[string]string{{
		"THN_TREATY": "2026", "REAS_ID": "1999999", "PCT_SHARE": "1", "REAS_SECURITY": "YATIM",
	}}); err != nil {
		t.Fatal(err)
	}
	lap, err := repository.NewMigrasiTCO(db).Pindahkan(ctx)
	if !errors.Is(err, repository.ErrMigrasiTCOTidakBersih) {
		t.Fatalf("security yatim harus dilaporkan sebelum Oracle menolaknya, dapat %v\n%s", err, lap)
	}
	ada := false
	for _, tm := range lap.Temuan {
		if tm.Jenis == repository.TemuanRujukanYatim && tm.Nilai == "1999999" {
			ada = true
		}
	}
	if !ada {
		t.Errorf("temuan rujukan yatim tidak menyebut 1999999:\n%s", lap)
	}
}
