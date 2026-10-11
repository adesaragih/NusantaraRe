-- Jalur mundur untuk 435_cabut_nilai_selisih.sql
--
-- Membangun kembali `NILAI_SELISIH` dan `NILAI_SEBELUM_PRO_RATE` DI MODUL
-- INI, bentuknya persis seperti `426_nilai_selisih.sql` membuatnya.
--
-- ⛔ WAJIB dijalankan BERSAMA pembalikan
-- `modul/treatyinadjustment/.../442_nilai_selisih_down.sql`. Dua modul yang
-- sama-sama membuat tabel yang sama akan bertabrakan: yang kedua gagal
-- dengan ORA-00955, dan pemasangan berhenti di tengah.
--
-- ⚠️ Alasan RANCANGANNYA tidak diulang di sini - ia ada di kepala `426`,
-- lengkap: kardinalitas 1:N lawan 1:1 ERD, penghalang
-- `BESARAN_DAPAT_DISESUAIKAN`, induk `VERSI_KONTRAK` lawan `KONTRAK`, dan
-- INV-58 (nol kolom selisih). Menyalinnya ke sini membuat salinannya
-- membeku pada hari yang asli berubah.
--
-- ⚠️ Sequence yang dibangun kembali memulai dari 1 dan menawarkan ulang
-- pengenal yang sudah terpakai (INV-03). Aman selama keduanya nol baris.
CREATE SEQUENCE {skema}.SEQ_TRIA_NILAI_SELISIH START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIA_NILAI_SBL_PRORATA START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.NILAI_SELISIH (
  ID_NILAI_SELISIH  NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK  NUMBER(19)          NOT NULL,
  KODE_BESARAN      VARCHAR2(1000 CHAR) NOT NULL,
  KODE_MATA_UANG    VARCHAR2(1000 CHAR) NOT NULL,
  NILAI_LAMA        NUMBER(38,8),
  NILAI_BARU        NUMBER(38,8),
  CONSTRAINT PK_NILAI_SELISIH PRIMARY KEY (ID_NILAI_SELISIH),
  CONSTRAINT FK_NILAI_SELISIH_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK) ON DELETE CASCADE
)
/
CREATE TABLE {skema}.NILAI_SEBELUM_PRO_RATE (
  ID_NILAI_SEBELUM_PRO_RATE  NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK           NUMBER(19)          NOT NULL,
  KODE_BESARAN               VARCHAR2(1000 CHAR) NOT NULL,
  KODE_MATA_UANG             VARCHAR2(1000 CHAR) NOT NULL,
  NILAI                      NUMBER(38,8),
  CONSTRAINT PK_NILAI_SEBELUM_PRO_RATE PRIMARY KEY (ID_NILAI_SEBELUM_PRO_RATE),
  CONSTRAINT FK_NILAI_SEBELUM_PRO_RATE_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_NILAI_SELISIH_VERSI ON {skema}.NILAI_SELISIH (ID_VERSI_KONTRAK)
/
CREATE INDEX {skema}.IX_NILAI_SBL_PRORATA_VERSI ON {skema}.NILAI_SEBELUM_PRO_RATE (ID_VERSI_KONTRAK)
/
