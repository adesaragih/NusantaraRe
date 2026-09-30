-- T_TREATYCO_JEJAK - jejak audit modul Treaty Contract Out (tiket 01).
--
-- ADR-0007: setiap penyimpanan dan penghapusan merekam SIAPA dan KAPAN
-- (AC 41 spec; tiket 10 menuntut pula "berapa baris tiap jenis" pada kaskade).
--
-- ⛔ Tabel BARU tanpa padanan warisan. USERID/TGLUPDATE pada tiap baris master
-- hanya menyimpan penulis TERAKHIR, dan baris yang dihapus tidak dapat
-- menyimpan siapa yang menghapusnya. Jejak karena itu tabel sendiri.
--
-- ⛔ WAKTU bertipe TIMESTAMP, bukan DATE: dua tindakan dalam detik yang sama
-- harus tetap terbedakan (pola T_CLAIMLF_JEJAK Claim Life).
--
-- ⛔ AKUN_ID adalah pengenal akun, bukan nama orang (ADR-U-0030).
--
-- ⛔ NOL COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_TREATYCO_JEJAK (
  ID         VARCHAR2(32) NOT NULL,
  WAKTU      TIMESTAMP,
  AKUN_ID    VARCHAR2(255),
  TABEL      VARCHAR2(255),
  BARIS_ID   VARCHAR2(32),
  AKSI       VARCHAR2(255),
  KETERANGAN VARCHAR2(1000),
  CONSTRAINT PK_T_TREATYCO_JEJAK PRIMARY KEY (ID)
)
/

CREATE INDEX {skema}.IDX_T_TREATYCO_JEJAK_BARIS ON {skema}.T_TREATYCO_JEJAK (TABEL, BARIS_ID)
/

CREATE SEQUENCE {skema}.SEQ_T_TREATYCO_JEJAK START WITH 1 INCREMENT BY 1 NOCACHE
/
