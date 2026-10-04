package repository

// SQL Marketing Officer TANPA Oracle: kolom sama dengan prosedur Pega, urutan bind, skema, dan nol COMMIT.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/marketingofficer/backend/models"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// INSERT menulis kolom yang SAMA, berurutan sama, dengan `INSERT` prosedur `PEGA_MARKETINGOFFICER`; TANGGAL SYSDATE.
func TestSisipSamaDenganProsedurPega(t *testing.T) {
	q := satuBaris(sqlSisipMO("S.MARKETINGOFFICER"))
	mau := "INSERT INTO S.MARKETINGOFFICER (ID, CLIENTID, BRANCHDETAILID, MOLEADER, MOSTATUS, CLIENTID2, TEAMGROUP, " +
		"CLIENTNAME, TANGGAL, USERUPDATE, BRANCHPARENT, BRANCHDETAILNAME, AKSES_LOGIN) " +
		"VALUES (:1, :2, :3, :4, :5, :6, :7, :8, SYSDATE, :9, :10, :11, :12)"
	if q != mau {
		t.Errorf("sisip:\n dapat %s\n mau   %s", q, mau)
	}
	m := models.MarketingOfficer{ID: "i", ClientID: "c", BranchDetailID: "bd", MOLeader: "ml", MOStatus: "1", ClientID2: "c2",
		TeamGroup: "tg", ClientName: "cn", UserUpdate: "u", BranchParent: "bp", BranchDetailName: "bn", AksesLogin: "a"}
	n := NilaiSisip(m)
	mauNilai := []any{"i", "c", "bd", "ml", "1", "c2", "tg", "cn", "u", "bp", "bn", "a"}
	if len(n) != len(mauNilai) {
		t.Fatalf("nilai sisip %d, mau %d", len(n), len(mauNilai))
	}
	for i := range n {
		if n[i] != mauNilai[i] {
			t.Errorf("bind :%d = %v, mau %v", i+1, n[i], mauNilai[i])
		}
	}
	if kosong := NilaiSisip(models.MarketingOfficer{ID: "i"}); kosong[1] != nil || kosong[11] != nil {
		t.Errorf("kosong = NULL: %v", kosong)
	}
}

// UPDATE tidak pernah menyentuh CLIENTID (Marketing Code) dan CLIENTNAME; ID di bind terakhir.
func TestPerbaruiTanpaMarketingCodeDanNama(t *testing.T) {
	q := satuBaris(sqlPerbaruiMO("S.MARKETINGOFFICER"))
	if strings.Contains(q, " CLIENTID =") || strings.Contains(q, "CLIENTNAME") || strings.Contains(q, "BRANCHSTATUS") {
		t.Errorf("perbarui menyentuh CLIENTID/CLIENTNAME/BRANCHSTATUS: %s", q)
	}
	if !strings.Contains(q, "TANGGAL = SYSDATE") || !strings.HasSuffix(q, "WHERE ID = :10") {
		t.Errorf("perbarui: %s", q)
	}
	n := NilaiPerbarui(models.MarketingOfficer{ID: "10000201", AksesLogin: "UJI"})
	if len(n) != 10 || n[9] != "10000201" || n[8] != "UJI" {
		t.Errorf("nilai perbarui: %v", n)
	}
}

// ID = 1 + 7 digit seperti `1||lpad(currency_seq.nextval,7,'0')`; lebih dari 7 digit ditolak, bukan dipotong.
func TestFormatIDMO(t *testing.T) {
	for n, mau := range map[int64]string{133: "10000133", 9999999: "19999999", 0: "10000000"} {
		if id, err := FormatIDMO(n); err != nil || id != mau {
			t.Errorf("%d: %q %v, mau %q", n, id, err, mau)
		}
	}
	if _, err := FormatIDMO(10000000); err != ErrNomorMelampaui {
		t.Errorf("8 digit: %v", err)
	}
}

// Setiap SQL berskema, tanpa COMMIT; kecuali-ID NULL tetap membandingkan (Oracle: teks kosong = NULL).
func TestSemuaSQLBerskemaTanpaCommit(t *testing.T) {
	for nama, q := range map[string]string{
		"daftar":         sqlDaftarMO("S.MARKETINGOFFICER"),
		"ambil":          sqlAmbilMO("S.MARKETINGOFFICER"),
		"nomor":          sqlNomorMO("S.CURRENCY_SEQ"),
		"sisip":          sqlSisipMO("S.MARKETINGOFFICER"),
		"perbarui":       sqlPerbaruiMO("S.MARKETINGOFFICER"),
		"aktif clientid": sqlAktifDenganClientID("S.MARKETINGOFFICER"),
		"aktif akses":    sqlAktifDenganAkses("S.MARKETINGOFFICER"),
		"clientid lama":  sqlClientIDDariAkses("S.MARKETINGOFFICER"),
		"daftar akun":    sqlDaftarAkun("S.M_LOGIN_GO"),
		"ambil akun":     sqlAmbilAkun("S.M_LOGIN_GO"),
		"daftar cabang":  sqlDaftarCabang("S.BRANCH"),
		"ambil cabang":   sqlAmbilCabang("S.BRANCH"),
	} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		if !strings.Contains(q, " S.") {
			t.Errorf("%s tanpa skema: %s", nama, q)
		}
	}
	for _, q := range []string{sqlAktifDenganClientID("T"), sqlAktifDenganAkses("T")} {
		if !strings.Contains(q, "ID <> NVL(:3, '-')") {
			t.Errorf("kecuali-ID harus NVL: %s", q)
		}
	}
	if q := satuBaris(sqlDaftarCabang("T")); !strings.Contains(q, "WHERE STATUS = :1") {
		t.Errorf("cabang aktif saja (BrowseBranchDetail_RD .Status = 1): %s", q)
	}
}

// Log: tertua dulu (LOG_TIME kosong lebih dulu, lalu TANGGAL, lalu ROWID); jalur sebelum migrasi 760 tanpa kedua
// kolom barunya; penanda hanya baris terakhir yang belum bertanda.
func TestSQLLogMO(t *testing.T) {
	if q := satuBaris(sqlLogMO("S.MARKETINGOFFICER_LOG")); !strings.HasSuffix(q, "WHERE ID = :1 ORDER BY LOG_TIME NULLS FIRST, TANGGAL NULLS FIRST, ROWID") ||
		!strings.Contains(q, "ACTION, TO_CHAR(LOG_TIME, 'YYYY-MM-DD HH24:MI:SS') FROM") {
		t.Errorf("log: %s", q)
	}
	if q := satuBaris(sqlLogMOLama("T")); strings.Contains(q, "LOG_TIME") || strings.Contains(q, "AKSES_LOGIN") {
		t.Errorf("log lama menyebut kolom 760: %s", q)
	}
	q := satuBaris(sqlTandaiLog("T"))
	if q != "UPDATE T SET ACTION = :1, AKSES_LOGIN = :2 WHERE ID = :3 AND ACTION IS NULL AND LOG_TIME = (SELECT MAX(LOG_TIME) FROM T WHERE ID = :4)" {
		t.Errorf("tandai log: %s", q)
	}
	for _, s := range []string{sqlLogMO("S.T"), sqlLogMOLama("S.T"), sqlTandaiLog("S.T")} {
		if err := db.PeriksaSQL(s); err != nil {
			t.Error(err)
		}
	}
}
