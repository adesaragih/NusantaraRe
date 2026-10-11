package repository

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// 11-10-2026: kotak masuk membawa PL Number (peserta aktif) dan Update Operator Name (jejak terakhir, cadangan pembuat)
// lewat subkueri skalar berindeks - satu kueri, tanpa N+1.
func TestSQLInboxPLNumberDanOperatorPengubah(t *testing.T) {
	q := strings.Join(strings.Fields(sqlInbox("S.W", "S.H", "S.P", "S.J", true)), " ")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{
		"(SELECT MIN(TRIM(p.PL_NUMBER)) FROM S.P p WHERE p.CLAIM_ID = w.ID AND p.STS_HAPUS IS NULL)",
		"NVL((SELECT MAX(j.AKUN_ID) KEEP (DENSE_RANK LAST ORDER BY j.WAKTU, j.ID) FROM S.J j WHERE j.KLAIM_ID = w.ID), w.CREATE_OP_NAME)",
		"ORDER BY w.TGL_CREATE DESC, w.ID DESC",
		"AND w.CREATE_OP = :akun",
	} {
		if !strings.Contains(q, mau) {
			t.Errorf("sqlInbox tanpa %q:\n%s", mau, q)
		}
	}
}

// 11-10-2026: efek unggah dokumen yang masih antre dibaca per rujukan (tabel outbox bersama), terkunci untuk dibatalkan.
func TestSQLEfekUnggahAntre(t *testing.T) {
	q := strings.Join(strings.Fields(sqlEfekUnggahAntre("S.O")), " ")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if q != "SELECT ID FROM S.O WHERE MODUL = :1 AND JENIS_EFEK = :2 AND RUJUKAN = :3 AND STATUS = :4 FOR UPDATE" {
		t.Errorf("efek unggah antre: %s", q)
	}
	if strings.Join(strings.Fields(sqlDokumenAda("S.D")), " ") != "SELECT COUNT(*) FROM S.D WHERE ID = :1" {
		t.Error("dokumen ada")
	}
}

// 11-10-2026: kategori dokumen = Document Claim produk (Master Product Name Life), urutan master; efek unggah satu
// dokumen dipungut per rujukan untuk dikirim segera sesudah unggah.
func TestSQLDokumenKlaimProdukDanPungutUnggah(t *testing.T) {
	q := strings.Join(strings.Fields(sqlDokumenKlaimProduk("S.DC")), " ")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if q != "SELECT d.DOCUMENT FROM S.DC d WHERE d.PRODUCTID = :1 ORDER BY d.URUT" || TabelDokumenKlaimProduk != "M_PRODUCTNAME_LIFE_DOCCLAIM" {
		t.Errorf("Document Claim produk: %s", q)
	}
	p := strings.Join(strings.Fields(sqlPungutEfekUnggahDokumen("S.O")), " ")
	if err := db.PeriksaSQL(p); err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"WHERE STATUS = :1 AND JADWAL_BERIKUT <= :2 AND MODUL = :3 AND JENIS_EFEK = :4 AND RUJUKAN = :5",
		"FETCH FIRST 1 ROWS ONLY FOR UPDATE SKIP LOCKED"} {
		if !strings.Contains(p, mau) {
			t.Errorf("pungut efek unggah tanpa %q: %s", mau, p)
		}
	}
}

// 11-10-2026: objek dokumen klaim di T_STORAGE_IMAGE - bentuk SQL sama dengan master / rule warisan (tanpa COMMIT).
func TestSQLObjekPenyimpanan(t *testing.T) {
	for q, mau := range map[string]string{
		sqlSisipObjek("S.SI"):    "INSERT INTO S.SI (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE) VALUES (:1, :2, :3, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, :6, 'standard')",
		sqlAmbilObjek("S.SI"):    "SELECT URLPUBLIC, APPFOLDER, TO_CHAR(EXPDATE, 'DD/MM/YYYY HH24:MI:SS'), APPNAME, FILENAME FROM S.SI WHERE IMAGEID = :1 FETCH FIRST 1 ROWS ONLY",
		sqlPerbaruiObjek("S.SI"): "UPDATE S.SI SET URLPUBLIC = :1, APPFOLDER = :2, EXPDATE = TO_DATE(:3, 'DD/MM/YYYY HH24:MI:SS'), TANGGAL_UPLOAD = TO_DATE(:4, 'MM/DD/YYYY HH24:MI:SS') WHERE IMAGEID = :5",
		sqlHapusObjek("S.SI"):    "DELETE FROM S.SI WHERE IMAGEID = :1",
		sqlNamaAplikasi("S.FI"):  "SELECT APPNAME FROM S.FI FETCH FIRST 1 ROWS ONLY",
	} {
		satu := strings.Join(strings.Fields(q), " ")
		if err := db.PeriksaSQL(satu); err != nil {
			t.Error(err)
		}
		if satu != mau {
			t.Errorf("SQL:\n%s\nmau\n%s", satu, mau)
		}
	}
	if TabelObjekPenyimpanan != "T_STORAGE_IMAGE" || TabelFolderPenyimpanan != "T_FOLDER_IMAGE" {
		t.Error("nama tabel bersama")
	}
}

// 11-10-2026: jumlah dokumen per peserta Accept AKTIF - satu kueri berkelompok (peserta dicabut diabaikan).
func TestSQLCacahDokumenPesertaAccept(t *testing.T) {
	q := strings.Join(strings.Fields(sqlCacahDokumenPesertaAccept("S.P", "S.D")), " ")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if q != "SELECT p.ID, COUNT(d.ID) FROM S.P p LEFT JOIN S.D d ON d.PREMIUM_LIST_DETAIL_ID = p.ID WHERE p.CLAIM_ID = :1 AND p.STS_HAPUS IS NULL AND p.IS_ACCEPT = 'true' GROUP BY p.ID" {
		t.Errorf("cacah dokumen: %s", q)
	}
}
