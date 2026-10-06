-- 924 - RICOMM_LIFE: tabel flat rincian R/I Comm Life (modul `ricommlife`, MASTER TREATY), menggantikan VIEW warisan
-- RICOMM_LIFE (`SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.CONTRACT, a.JSONDATA.YEAR,
-- a.JSONDATA.COMM FROM M_RICOMM_LIFE a`, dicek work owner 06-10-2026). Keputusan work owner 06-10-2026 butir 1-3:
-- tabel flat lewat migrasi inti (modul tanpa rentang sendiri, `—`); M_RICOMM_LIFE (JSON, 0 baris DEV) TIDAK disentuh,
-- isinya dipindah alat `modul/ricommlife/backend/alat/pindahflat`.
--
-- ⛔ PRASYARAT: VIEW POOLDATA.RICOMM_LIFE sudah dibuang DBA lewat
-- `modul/ricommlife/docs/DBA-LEPAS-VIEW-RICOMM_LIFE.sql` SEBELUM `-migrate`. Selama view itu ada, pra-terbang pelari
-- (`praTerbangBentuk`/`objekAda`) menemukan objek bernama sama dengan kolom yang sama, CREATE TABLE dijawab ORA-00955
-- dan DILEWATI, lalu langkah ini tercatat selesai tanpa tabel. Urutan langkah work owner: modul/ricommlife/MODUL.md.
--
-- Tipe ikut section Pega `InboxRIComm` (kelas ASM-FW-GISFW-Int-RI_COMM_LIFE); XML tidak memberi presisi, jadi
-- preseden K6 (TestNolNumberTanpaPresisi): CONTRACT (pxNumber b1907) dan YEAR (pyMax 4 b2206) bilangan bulat
-- NUMBER(5), COMM (grid pxNumber b9101) desimal NUMBER(38,8). Bukti per kolom: docs/STRUKTUR-TABEL-RICOMMLIFE.md.
-- Seluruh kolom NULLABLE kecuali PK; wajib-isi ditegakkan Go. ID: site || LPAD(M_RICOMM_LIFE_SEQ.NEXTVAL, 6, '0')
-- dari Go (sequence warisan, nol sequence baru). Batas transaksi milik Go (ADR-U-0029); NOL COMMIT.
CREATE TABLE {skema}.RICOMM_LIFE (
  ID        VARCHAR2(10) NOT NULL,
  IDUSEDBY  VARCHAR2(10),
  USEDBY    VARCHAR2(200),
  CONTRACT  NUMBER(5),
  YEAR      NUMBER(5),
  COMM      NUMBER(38,8),
  CONSTRAINT PK_RICOMM_LIFE PRIMARY KEY (ID)
)
/
CREATE INDEX {skema}.IX_RICOMM_LIFE_IDUSEDBY ON {skema}.RICOMM_LIFE (IDUSEDBY)
/
