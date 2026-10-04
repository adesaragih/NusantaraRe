-- Tabel Komite Claim Life - butir af.
--
-- Pemilik: A1/A2 (brief lanjutan 4 §1, `[DIPUTUSKAN 27-09-2026]`). Bentuknya
-- dari `.scratch/komite-claim-life/STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md`, yang
-- §1 sebut sebagai sumbernya; tiket 00 Komite kelak MEMVERIFIKASI, bukan
-- membuat ulang.
--
-- ⛔ RONDE PERTAMA BERKAS INI SALAH DI LIMA TEMPAT, dan seluruhnya ketahuan
-- ketika dokumen sumbernya benar-benar dibaca - bukan ketika ia dikira sudah
-- diketahui:
--
--   1. awalan `ID` `KMT-`        -> `KMTLF-` (lintas-lini: Life vs Prop)
--   2. `ADJUSTMENT_ID` nullable  -> NOT NULL, dan ber-index UNIK
--   3. `REFERENCES` ke adjustment -> DIBUANG (dua tabel tujuan menurut LINI)
--   4. FK ke induk `KOMITE_ID`   -> `DATA_KOMITE_ID`
--   5. `KOMITE_ID`/`ID_KOMITE`   -> `KOMITE_OPERATORID`/`KOMITE_JABATAN`
--
-- Berkas ini belum pernah dijalankan di Oracle mana pun (G1), jadi
-- memperbaikinya di tempat adalah yang benar - bukan menambah langkah ralat.

-- ⛔ SHARED PRIMARY KEY. `ID` = `T_WORK_CLAIM.ID` baris komite, teks
-- `KMTLF-xxxxxx`. Nol kolom `WORK_CLAIM_ID` (REVISI 2026-09-18) - dua jalan
-- menuju satu baris akhirnya berbeda.
--
-- ⛔ `ADJUSTMENT_ID` TANPA `REFERENCES`, dan itu disengaja serta dinyatakan:
-- `T_GENERAL_KOMITE` lintas-lini, sehingga satu kolom ini punya DUA tabel
-- tujuan menurut `T_WORK_CLAIM.LINI` - LIFE ke `T_CLAIMLF_ADJUSTMENT`, PROP ke
-- `T_CLAIMP_ADJUSTMENT`. Oracle hanya dapat menunjuk satu, jadi klausanya
-- tidak dipasang sama sekali dan KEUTUHANNYA DIJAGA KODE GO.
--
-- ⚠️ `[keputusan work owner]` satu baris `AdjustmentList` = TEPAT satu kasus
-- komite, dan sebaliknya. Karena itu `ADJUSTMENT_ID` NOT NULL dan ber-index
-- UNIK: dua kasus komite tidak boleh menunjuk baris adjustment yang sama.
--
-- `[terverifikasi]` kosakata kolom dari `Komite Claim Life/Activity/
-- KomitePostAdjustment.xml`: `KomiteLoop` (1322), `KomiteCount` (1398),
-- `AcceptStatus` (gerbang 5695, 8119, 8648, 8887).
CREATE TABLE {skema}.T_GENERAL_KOMITE (
  ID             VARCHAR2(40) NOT NULL,
  ADJUSTMENT_ID  VARCHAR2(40) NOT NULL,
  KOMITE_LOOP    NUMBER(5),
  KOMITE_COUNT   NUMBER(5),
  ACCEPT_STATUS  VARCHAR2(8),
  CONSTRAINT PK_GENERAL_KOMITE PRIMARY KEY (ID),
  CONSTRAINT FK_GENERAL_KOMITE_WORK FOREIGN KEY (ID)
    REFERENCES {skema}.T_WORK_CLAIM (ID)
)
/
CREATE UNIQUE INDEX {skema}.UX_GENERAL_KOMITE_ADJ ON {skema}.T_GENERAL_KOMITE (ADJUSTMENT_ID)
/
-- Roster dan keputusan komite. Satu baris = satu anggota pada satu jenjang
-- tangga dari satu kasus komite.
--
-- ⛔ TIGA NAMA KOLOM DIBETULKAN `[keputusan work owner 2026-09-18 sore]`,
-- sebab nama korpusnya MENYESATKAN ke dua arah sekaligus:
--
--   `KomiteID`  isinya akun OPERATOR -> `KOMITE_OPERATORID`
--   `IDKomite`  isinya JABATAN       -> `KOMITE_JABATAN`
--   `KomiteAproval` (satu P)         -> `KOMITE_APPROVAL` (ejaan dibetulkan)
--
-- `[terverifikasi]` sensus 555 berkas korpus: `.KomiteID := .OPERATOR_ID`
-- (8 penulisan, 5 berkas) dan `.IDKomite := .JABATAN` (12 penulisan, 7
-- berkas). Dua properti bernama nyaris sama, isinya berbeda.
--
-- ⚠️ `KOMITE_OPERATORID` dan `KOMITE_EMAIL` memuat DATA ORANG. Dibaca saat
-- jalan dari `EMAILKOMITE`; nol baris disalin ke fixture, tiket, maupun log.
--
-- ⛔ `KOMITE_APPROVAL` bertipe TEKS bernilai `"0"`/`"1"`/`"2"` - kode, bukan
-- bilangan (ADR-U-0022).
CREATE TABLE {skema}.T_KOMITE_KOMITELIST (
  ID                VARCHAR2(40) NOT NULL,
  DATA_KOMITE_ID    VARCHAR2(40),
  KOMITE_URUT       NUMBER(5),
  KOMITE_OPERATORID VARCHAR2(64),
  KOMITE_JABATAN    VARCHAR2(140),
  KOMITE_EMAIL      VARCHAR2(255),
  KOMITE_APPROVAL   VARCHAR2(8),
  KOMITE_COMMENT    VARCHAR2(4000),
  DATE_APPROVE      DATE,
  CONSTRAINT PK_KOMITE_KOMITELIST PRIMARY KEY (ID),
  CONSTRAINT FK_KOMITELIST_KOMITE FOREIGN KEY (DATA_KOMITE_ID)
    REFERENCES {skema}.T_GENERAL_KOMITE (ID)
)
/
CREATE INDEX {skema}.IX_KOMITELIST_KOMITE ON {skema}.T_KOMITE_KOMITELIST (DATA_KOMITE_ID)
/
CREATE SEQUENCE {skema}.SEQ_KOMITE_KOMITELIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
