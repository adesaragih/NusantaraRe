-- ⛔⛔ BERKAS INI MEMBUAT TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444` dan `445`, dengan sebab yang sama: rentang `treatyin` (`400-439`)
-- PENUH, dan `446` jatuh di rentang modul ini (`440-479`). Penjaga DDL
-- `treatyin` (`migrasi_pendaratan_test.go`) membaca folder ini juga.
--
-- ---------------------------------------------------------------------
-- `T_TREATY_HAZARD_LIMIT` — dari `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`
-- ---------------------------------------------------------------------
-- Permintaan pemakai 7 Oktober 2026: "cek di file ini kemana tablenya
-- disimpan … kalau blm ada tambahkan saja". Keempat lembar diagram v2
-- memuat kotak:
--
--   T_TREATY_HAZARD_LIMIT   1:1 · BERSAMA → Prop, Non Prop, EDM Prop, EDM Non Prop
--   FK TREATY_IN_ID → TREATY_IN.ID   CASCADE
--   ← TreatyIn   [BATAS_BAHAYA]   bukti: JALUR-PENUH
--
-- dan `peta-nama-tabel-treatyin.tsv` merinci isinya: "Earthquake FloodJab
-- FloodNation RSMDLimit MaxCoGroup MaxCoNonGroup + Currency… — 10 skalar".
-- Kesepuluhnya skalar AKAR `TreatyIn`, dan sebelum tabel ini nol kolom
-- pendaratan menampungnya:
--
--   RSMDLIMIT / CURRENCYRSMD               tab Event Limits (Non-Prop)
--   EARTHQUAKE / CURRENCYEARTHQUAKE        ""
--   FLOODJAB / CURRENCYFLOODJAB            ""
--   FLOODNATION / CURRENCYFLOODNAT         ""
--   MAXCOGROUP / MAXCONONGROUP             tab Co-Ins Scale (Prop), sel 277/278
--
-- ⛔ RALAT atas kepala `439`, yang menolak tabel ini karena "nilainya sudah
-- berumah di `T_TREATY_LIMIT_DETAIL`". Itu benar untuk nilai TINGKAT DETAIL
-- (`Limits[].Detail[]`, tab Limits Prop). Tab Event Limits NON-PROP
-- mengikat properti AKAR — terukur 6 Oktober 2026: akar berisi di 49 dari
-- 772 kontrak Non-Prop, `Detail[]` di 1 — dan Max Co terisi di 173/87 dari
-- 1.855 dokumen. Keduanya tidak punya rumah lain.
--
-- ⚠️ BENTUK 1:1, bukan "diputar jadi daftar": diagram menyatakan
-- kardinalitas 1:1 dengan `TREATY_IN`, dan pola tabel akar yang sudah ada
-- (`T_TREATY_REVISION`) — satu baris per dokumen, `URUTAN` 0 — dipakai apa
-- adanya. Memutarnya menjadi satu baris per bahaya adalah bentuk model
-- sasaran (`BATAS_PER_BAHAYA`, `KAMUS-KOLOM.md` §10.18), bukan pendaratan.
--
-- ⛔ NOL kunci asing ke `TREATY_IN`: tabel warisan itu tidak punya PK/UNIQUE
-- (`ORA-02270`), sama seperti seluruh tabel pendaratan. Kaitannya lewat
-- `MASTERID`; sisi `OLDDATA` penyesuaian memakai akhiran `#LAMA`.
--
-- ⚠️ `VARCHAR2(4000 CHAR)` seperti seluruh kolom pendaratan: nilainya teks
-- apa adanya dari dokumen, ditafsirkan services.

CREATE TABLE {skema}.T_TREATY_HAZARD_LIMIT (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  RSMDLIMIT                        VARCHAR2(4000 CHAR),
  CURRENCYRSMD                     VARCHAR2(4000 CHAR),
  EARTHQUAKE                       VARCHAR2(4000 CHAR),
  CURRENCYEARTHQUAKE               VARCHAR2(4000 CHAR),
  FLOODJAB                         VARCHAR2(4000 CHAR),
  CURRENCYFLOODJAB                 VARCHAR2(4000 CHAR),
  FLOODNATION                      VARCHAR2(4000 CHAR),
  CURRENCYFLOODNAT                 VARCHAR2(4000 CHAR),
  MAXCOGROUP                       VARCHAR2(4000 CHAR),
  MAXCONONGROUP                    VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_HAZARD_LIMIT PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_HAZARD_LIMIT UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_HAZARD_LIMIT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
