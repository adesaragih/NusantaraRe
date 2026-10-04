-- 182 - T_GENERAL_POLIS SEBAGIAN: kolom yang dipakai blok General layar Inward Facultative (tiket 31), butir 78.
--
-- Keputusan work owner 03-10-2026 (butir 78.4): tabel flat dibuat SEKARANG hanya dengan kolom sistem +
-- kolom yang dipakai layar; kolom lain (termasuk ~30 kolom uang NUMBER yang menunggu keputusan presisi tim inti)
-- ditambah tiket 23 lewat ALTER. Nama dan tipe kolom = rancangan (`loader/skema_gen.go` = DDL draf
-- `08-flat`); ID VARCHAR2(32) = butir 76.1 (berbagi PK dengan T_WORK_POLIS, K-064 relasi 52), tanpa
-- constraint FK lintas modul (pola A79). OLD_POLIS_ID (NUMBER, K-071: NULL di fase 1) belum dibuat.
--
-- Tiga kolom tanggal menyimpan TEKS BENTUK PEGA (butir 78.1): OFFERING_DATE 'YYYYMMDD'; START/END_DATE_TIME
-- 'YYYYMMDDTHHMMSS.mmm GMT', ditulis 12:00 WIB = 05:00 GMT (78.1 diralat). Seluruh kolom nullable kecuali PK. Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_GENERAL_POLIS (
  ID              VARCHAR2(32) NOT NULL,
  IDPEGA          VARCHAR2(50),
  COB_GROUP       VARCHAR2(20),
  START_DATE_TIME VARCHAR2(30),
  OFFERING_DATE   VARCHAR2(30),
  END_DATE_TIME   VARCHAR2(30),
  FOLLOWING       VARCHAR2(50),
  CONSTRAINT PK_T_GENERAL_POLIS PRIMARY KEY (ID)
)
/
