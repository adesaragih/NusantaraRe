-- Tabel Komite Claim Life - butir af.
--
-- Pemilik: A1 (brief lanjutan 4 §1, `[DIPUTUSKAN 27-09-2026]`). Bentuknya dari
-- tiket 00 Komite Claim Life §"Bentuk yang dibangun"; tiket 00 kelak
-- MEMVERIFIKASI tabel ini, bukan membuatnya ulang.
--
-- ⛔ SHARED PRIMARY KEY. `T_GENERAL_KOMITE.ID` = `T_WORK_CLAIM.ID` baris
-- komite, bertipe TEKS `KMT-xxxxxx`. Tidak ada kolom `WORK_CLAIM_ID`
-- (REVISI 2026-09-18) - menambahkannya berarti dua jalan menuju satu baris,
-- dan dua jalan akhirnya berbeda.
--
-- ⛔ `ADJUSTMENT_ID` berada DI SINI, bukan di `T_WORK_CLAIM` - keputusan tiket
-- 00 Komite. Bersama `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` ia membentuk penunjuk
-- DUA ARAH, dan keduanya diisi dalam SATU transaksi saat kirim komite.
--
-- `[terverifikasi]` kosakata kolom dari `Komite Claim Life/Activity/
-- KomitePostAdjustment.xml`: `KomiteLoop` (baris 1322), `KomiteCount` (1398),
-- `KomiteAproval` (912, 993, 5899, 5985), `KomiteComment` (893, 5945),
-- `DateApprove` (935, 5965), `AcceptStatus` (gerbang 5695, 8119, 8648, 8887).
CREATE TABLE {skema}.T_GENERAL_KOMITE (
  ID             VARCHAR2(40) NOT NULL,
  ADJUSTMENT_ID  VARCHAR2(40),
  KOMITE_LOOP    NUMBER(5),
  KOMITE_COUNT   NUMBER(5),
  ACCEPT_STATUS  VARCHAR2(8),
  CONSTRAINT PK_GENERAL_KOMITE PRIMARY KEY (ID),
  CONSTRAINT FK_GENERAL_KOMITE_WORK FOREIGN KEY (ID)
    REFERENCES {skema}.T_WORK_CLAIM (ID),
  CONSTRAINT FK_GENERAL_KOMITE_ADJ FOREIGN KEY (ADJUSTMENT_ID)
    REFERENCES {skema}.T_CLAIMLF_ADJUSTMENT (ID)
)
/
CREATE INDEX {skema}.IX_GENERAL_KOMITE_ADJ ON {skema}.T_GENERAL_KOMITE (ADJUSTMENT_ID)
/
-- Keputusan per anggota komite, satu baris per tingkat tangga.
--
-- ⛔ `KOMITE_APROVAL` bertipe TEKS dan bernilai `"0"`, `"1"`, atau `"2"` -
-- kode, bukan bilangan (ADR-U-0022). Ejaan `APROVAL` (sic) dipertahankan
-- seperti korpus.
--
-- ⚠️ `KOMITE_EMAIL` dan `ID_KOMITE` memuat DATA ORANG. Ia dibaca saat jalan
-- dari `EMAILKOMITE`; nol baris disalin ke fixture, tiket, maupun log.
CREATE TABLE {skema}.T_KOMITE_KOMITELIST (
  ID             VARCHAR2(40) NOT NULL,
  KOMITE_ID      VARCHAR2(40) NOT NULL,
  KOMITE_URUT    NUMBER(5),
  ID_KOMITE      VARCHAR2(64),
  KOMITE_EMAIL   VARCHAR2(255),
  KOMITE_APROVAL VARCHAR2(8),
  KOMITE_COMMENT VARCHAR2(4000),
  DATE_APPROVE   DATE,
  CONSTRAINT PK_KOMITE_KOMITELIST PRIMARY KEY (ID),
  CONSTRAINT FK_KOMITELIST_KOMITE FOREIGN KEY (KOMITE_ID)
    REFERENCES {skema}.T_GENERAL_KOMITE (ID)
)
/
CREATE INDEX {skema}.IX_KOMITELIST_KOMITE ON {skema}.T_KOMITE_KOMITELIST (KOMITE_ID)
/
CREATE SEQUENCE {skema}.SEQ_KOMITE_KOMITELIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
