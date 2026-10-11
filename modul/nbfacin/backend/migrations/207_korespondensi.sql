-- 207 - tab Correspondence kasus Inward Facultative (tiket 56): `.OfferFacIn.CorrespondenceList` dan `.RecipientList`.
--
-- TABEL BARU (keputusan work owner 06-10-2026 "saya setuju dengan rekomendasi kamu"): tidak ada di rancangan loader dan
-- `DDL\` (Pega menyimpannya di halaman kerja + JSON_OFFER; sampel CONTOH 113 CorrespondenceList, SEMUA kosong). Nama
-- mengikuti properti Pega (kelas Data-Correspondence):
--   T_CORRESPONDENCELIST induk T_GENERAL_POLIS (PARENT_ID = ID kasus VARCHAR2(32), butir 76.1), urut SEQ_NO = urutan Add;
--                        TANGGAL / SENT_AT = teks Pega DateTime (pola START_DATE_TIME); SENDER / SENDER_MAIL = akun login
--                        (M_LOGIN_GO.NAME 150 / EMAIL 254, inti 902 / 904); MESSAGE = HTML yang sudah dibersihkan layanan,
--                        VARCHAR2(4000) - penjaga inti melarang CLOB di migrasi (pengecualian hanya lewat tim inti), jadi
--                        pesan bersih > 4000 bita ditolak 400; IS_EDITABLE "1" draf / "0" terkirim (`.IsEditable`);
--                        CREATE_OP = LOGIN_ID pembuat.
--   T_RECIPIENTLIST      induk T_CORRESPONDENCELIST, urut SEQ_NO; RECIPIENT 150 / RECIPIENT_MAIL 254; IS_EDITABLE ikut induk.
-- ID NUMBER(19) dari sequence (A92). Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_CORRESPONDENCELIST (
  ID          NUMBER(19) NOT NULL,
  PARENT_ID   VARCHAR2(32) NOT NULL,
  SEQ_NO      NUMBER(5) NOT NULL,
  ROW_UID     VARCHAR2(36) NOT NULL,
  TANGGAL     VARCHAR2(30),
  SENDER      VARCHAR2(150),
  SENDER_MAIL VARCHAR2(254),
  SUBJECT     VARCHAR2(1000),
  MESSAGE     VARCHAR2(4000),
  IS_EDITABLE VARCHAR2(1) DEFAULT '1' NOT NULL,
  SENT_AT     VARCHAR2(30),
  CREATE_OP   VARCHAR2(64),
  TGL_CREATE  DATE,
  CONSTRAINT PK_T_CORRESPONDENCELIST PRIMARY KEY (ID),
  CONSTRAINT FK_CORRESPONDENCE_GENERAL FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID)
)
/
CREATE INDEX {skema}.IX_CORRESPONDENCE_PARENT ON {skema}.T_CORRESPONDENCELIST (PARENT_ID, SEQ_NO)
/
CREATE SEQUENCE {skema}.SEQ_T_CORRESPONDENCELIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.T_RECIPIENTLIST (
  ID             NUMBER(19) NOT NULL,
  PARENT_ID      NUMBER(19) NOT NULL,
  SEQ_NO         NUMBER(5) NOT NULL,
  ROW_UID        VARCHAR2(36) NOT NULL,
  RECIPIENT      VARCHAR2(150),
  RECIPIENT_MAIL VARCHAR2(254),
  IS_EDITABLE    VARCHAR2(1) DEFAULT '1' NOT NULL,
  CONSTRAINT PK_T_RECIPIENTLIST PRIMARY KEY (ID),
  CONSTRAINT FK_RECIPIENT_CORRESPONDENCE FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_CORRESPONDENCELIST (ID)
)
/
CREATE INDEX {skema}.IX_RECIPIENTLIST_PARENT ON {skema}.T_RECIPIENTLIST (PARENT_ID, SEQ_NO)
/
CREATE SEQUENCE {skema}.SEQ_T_RECIPIENTLIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
