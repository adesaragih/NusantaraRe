-- 197 - Clauses kasus FIRE (tiket 47): T_CLAUSELIST = `pyWorkPage.ClauseList` (kelas Data-Clause; `SaveJsonOfferFacIn_Act`
-- `OfferFacIn.ClauseList = pyWorkPage.ClauseList` - di AKAR kasus, induk T_GENERAL_POLIS seperti LocationList) BESERTA
-- `.ClauseList(n).ArgumentList` dalam SATU tabel + sequence.
--
-- ClauseList TIDAK ada di workbook rancangan (`loader/skema_gen.go`, 78 tabel) dan kosong di seluruh 108 contoh kasus
-- `DDL\CONTOH` `[terverifikasi]` - jadi tabel ini rancangan agent (K47-1, menunggu konfirmasi), berpola tabel rancangan:
-- ID NUMBER(19) dari sequence, IDPEGA / COB_GROUP / PARENT_ID / SEQ_NO / ROW_UID. Kolom = properti yang ditulis
-- `Activity\SearchClauseFireSQL_PostAct.xml` (ClauseCode, ClauseTitle, ClauseDescription, ClauseLanguage, ClauseLanguageID,
-- ClauseContent, ClauseContentTemp, ArgumentCount) dan grid "Isi Klasula" (ArgumentNumber, ArgumentDescription,
-- ArgumentValue).
--
-- SATU TABEL (K47-4; work owner 05-10-2026: "T_CLAUSELIST dan T_CLAUSEARGUMENTLIST bukannya bisa jadi 1 tabel aja ?
-- T_CLAUSEARGUMENTLIST disimpan di T_CLAUSELIST juga"): SATU BARIS PER ARGUMEN - kolom klausa berulang di setiap baris
-- argumennya, `ARGUMENT_SEQ_NO` 1..n; klausa tanpa argumen = satu baris ber-`ARGUMENT_SEQ_NO` dan kolom argumen kosong.
-- SEQ_NO = urutan klausa (sama di semua baris satu klausa). Daftar argumen sebagai teks terstruktur di satu kolom TIDAK
-- dipakai: penjaga migrasi `TestKolomUangDesimalDanNolJSON` melarang atribut di dalam dokumen (keputusan tim inti).
--
-- Tipe (K47-2): `M_CLAUSE` = TABEL (ID VARCHAR2(7), OLDID VARCHAR2(7), JSONDATA berkas besar; `DDL\M_CLAUSE.txt`
-- 05-10-2026), tetapi ClauseCode = `.ID` view CLAUSE = JSON `ID` (dot-notation), bukan kolom tabel M_CLAUSE.ID - lebarnya
-- tidak dijamin 7 -> CLAUSE_CODE VARCHAR2(4000) seperti kolom dot-notation lain. Isi klausa juga VARCHAR2(4000):
-- kolom dokumen besar dilarang penjaga migrasi, dan isi asal dibaca lewat dot-notation (maks 4000) - isi sesudah argumen
-- diganti yang melebihi 4000 bita DITOLAK aplikasi (400), tidak dipotong. ClauseLanguage kode "0" / "1" / "2" ->
-- VARCHAR2(10); ArgumentCount teks angka (hasil COUNT, frontend teks) -> VARCHAR2(50) seperti angka di tabel rancangan.
-- Diganti utuh tiap Save. Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_CLAUSELIST (
  ID                   NUMBER(19) NOT NULL,
  IDPEGA               VARCHAR2(50),
  COB_GROUP            VARCHAR2(20),
  PARENT_ID            VARCHAR2(32) NOT NULL,
  SEQ_NO               NUMBER(5) NOT NULL,
  ARGUMENT_SEQ_NO      NUMBER(5),
  ROW_UID              VARCHAR2(36) NOT NULL,
  CLAUSE_CODE          VARCHAR2(4000),
  CLAUSE_TITLE         VARCHAR2(4000),
  CLAUSE_DESCRIPTION   VARCHAR2(4000),
  CLAUSE_LANGUAGE      VARCHAR2(10),
  CLAUSE_LANGUAGE_ID   VARCHAR2(4000),
  CLAUSE_CONTENT       VARCHAR2(4000),
  CLAUSE_CONTENT_TEMP  VARCHAR2(4000),
  ARGUMENT_COUNT       VARCHAR2(50),
  ARGUMENT_NUMBER      VARCHAR2(4000),
  ARGUMENT_DESCRIPTION VARCHAR2(4000),
  ARGUMENT_VALUE       VARCHAR2(4000),
  CONSTRAINT PK_T_CLAUSELIST PRIMARY KEY (ID),
  CONSTRAINT FK_CLAUSELIST_GENERAL FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID)
)
/
CREATE INDEX {skema}.IX_CLAUSELIST_PARENT ON {skema}.T_CLAUSELIST (PARENT_ID, SEQ_NO, ARGUMENT_SEQ_NO)
/
CREATE SEQUENCE {skema}.SEQ_T_CLAUSELIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
