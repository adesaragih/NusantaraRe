package skemauji

// Tiruan objek warisan yang DIUBAH migrasi modul `companydetail` (04-10-2026): `CLIENT`, `CLIENT_PICLIST`,
// `CLIENT_ADDRESS` (805-807 menambah kolom), `M_CLIENT` beserta dua triggernya (808 menonaktifkannya), dan `M_NATION`
// beserta view `NATION` (802-804 mengganti view itu dengan tabel datar). Kolom dari katalog DEV 03/04-10-2026.
//
// Seperti tiruan Marketing Officer, semuanya dibuat SEBELUM migrasi. Berbeda darinya, semuanya dibongkar SESUDAH
// migrasi mundur: jalur mundur 808 menghidupkan trigger dan jalur mundur 802 membuat ulang view `NATION` dari
// `M_NATION` - keduanya butuh objek tiruannya masih ada.
//
// ⛔ Nol nama orang, nol data produksi - hanya bentuk tabel; trigger tiruan tidak berbuat apa pun.

import "fmt"

// namaTabelTiruanCompanyDetail - tabel tiruan, urut dibongkar (view NATION ikut lewat namaViewTiruanCompanyDetail).
var namaTabelTiruanCompanyDetail = []string{"CLIENT_ADDRESS", "CLIENT_PICLIST", "CLIENT", "M_CLIENT", "M_NATION"}

// namaViewTiruanCompanyDetail - view tiruan; sesudah migrasi maju 802 ia sudah tidak ada (ORA-00942 ditoleransi).
var namaViewTiruanCompanyDetail = []string{"NATION"}

// ddlTiruanCompanyDetail - DDL seluruh tiruan, urut dibuat.
func ddlTiruanCompanyDetail(skema string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE %s.CLIENT (ID VARCHAR2(100), OLDID VARCHAR2(1000), NAME VARCHAR2(4000),
  IDNUMBER VARCHAR2(1000), MCLID VARCHAR2(1000), MCL_TGL_LAHIR VARCHAR2(1000), FLAG VARCHAR2(1000), BU_ID VARCHAR2(1000),
  BU_NOTE VARCHAR2(1000), IDVIEW VARCHAR2(1000), NPWP VARCHAR2(1000), GROUPNAME VARCHAR2(1000), TITLE VARCHAR2(10),
  COUNTRY VARCHAR2(20), COUNTRYNAME VARCHAR2(50))`, skema),
		fmt.Sprintf(`CREATE TABLE %s.CLIENT_PICLIST (CLIENTID VARCHAR2(100), USERIDENTIFIER VARCHAR2(500),
  NICKNAME VARCHAR2(500), EMAIL VARCHAR2(500), POSITION VARCHAR2(100), DATEOFBIRTH VARCHAR2(100), PHONENUMBER VARCHAR2(100))`, skema),
		fmt.Sprintf(`CREATE TABLE %s.CLIENT_ADDRESS (ASMADDRESS VARCHAR2(500), ASMADDRESSTYPE VARCHAR2(10),
  ASMCITY VARCHAR2(10), CITYNAME VARCHAR2(100), PXCREATEOPERATOR VARCHAR2(100), CLIENTID VARCHAR2(50),
  PXCREATEDATETIME VARCHAR2(50), ASMZIPCODE VARCHAR2(50), DISTRICTNAME VARCHAR2(100), PROVINCENAME VARCHAR2(100),
  RWNAME VARCHAR2(100))`, skema),
		fmt.Sprintf(`CREATE TABLE %s.M_CLIENT (ID VARCHAR2(100), OLDID VARCHAR2(100), JSONDATA CLOB,
  CONSTRAINT ENSURE_UJI_M_CLIENT_JSON CHECK (JSONDATA IS JSON))`, skema),
		fmt.Sprintf(`CREATE OR REPLACE TRIGGER %s.TRG_M_CLIENT AFTER INSERT ON %s.M_CLIENT FOR EACH ROW
BEGIN
  NULL;
END;`, skema, skema),
		fmt.Sprintf(`CREATE OR REPLACE TRIGGER %s.TRG_M_CLIENT_PIC AFTER INSERT ON %s.M_CLIENT FOR EACH ROW
BEGIN
  NULL;
END;`, skema, skema),
		fmt.Sprintf(`CREATE TABLE %s.M_NATION (ID VARCHAR2(6), OLDID VARCHAR2(6), JSONDATA CLOB,
  CONSTRAINT ENSURE_UJI_M_NATION_JSON CHECK (JSONDATA IS JSON))`, skema),
		fmt.Sprintf(`CREATE VIEW %s.NATION AS SELECT a.JSONDATA.ID, a.OLDID, a.JSONDATA.Note, a.JSONDATA.NationInitial
  FROM %s.M_NATION a`, skema, skema),
	}
}
