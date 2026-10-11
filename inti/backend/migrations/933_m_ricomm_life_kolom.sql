-- 933 - M_RICOMM_LIFE mendapat kolom rincian (langkah 1 dari 2; langkah 2 = 934).
--
-- Keputusan work owner 08-10-2026: rincian R/I Comm Life = tabel flat M_RICOMM_LIFE berkolom PERSIS seperti tabel
-- RICOMM_LIFE (924): ID (PK warisan) + IDUSEDBY VARCHAR2(10), USEDBY VARCHAR2(200), CONTRACT NUMBER(5), YEAR NUMBER(5),
-- COMM NUMBER(38,8) - tipe dan lebar 924 (bukti per kolom: modul/ricommlife/docs/STRUKTUR-TABEL-RICOMMLIFE.md); tabel
-- RICOMM_LIFE dibuang 934; JSONDATA dibuang 934; nol tabel baru.
-- Berdiri sendiri (alasan sama dengan 931). NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.M_RICOMM_LIFE ADD (
  IDUSEDBY  VARCHAR2(10),
  USEDBY    VARCHAR2(200),
  CONTRACT  NUMBER(5),
  YEAR      NUMBER(5),
  COMM      NUMBER(38,8)
)
/
