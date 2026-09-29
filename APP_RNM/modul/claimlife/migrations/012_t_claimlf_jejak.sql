-- Jejak audit setiap transisi dan jalur balik - butir am.
--
-- Pemilik: A1 (brief lanjutan 4 §1, `[DIPUTUSKAN 27-09-2026]`); bentuknya
-- milik executor, AC-nya milik tiket 09.
--
-- ⛔ ADR-U-0007 menuntut SIAPA dan KAPAN untuk setiap transisi status klaim
-- DAN setiap jalur balik (`SendtoAdmin`, `SendtoMedical`). Tanpa tabel ini
-- KELIMA jalur tulis modul menjawab HTTP 501 - menolak baris (05), memindah
-- tahap (08), menyerahkan ke Komite (10), membuka putaran (11), dan mengaksep
-- (audit A0).
--
-- ⚠️ `[keputusan work owner 2026-09-14]` Jejak lama TIDAK dapat direkonstruksi
-- ke belakang: data sebelum cutover hanya punya CREATEOPNAME + empat tanggal.
-- Tabel ini karena itu mulai KOSONG, dan itu keadaan yang benar.
--
-- Bentuknya:
--   - ID dari sequence, seperti seluruh tabel T_CLAIMLF_* (ADR-U-0006).
--   - ADJUSTMENT_ID menunjuk baris; unit keputusan adalah BARIS (ADR-U-0011).
--     Ia TIDAK ber-FK: jejak harus selamat dari penghapusan apa pun, dan
--     jejak yang ikut terhapus bersama yang dijejakinya bukan jejak.
--   - DARI dan KE bertipe TEKS, bukan angka: keduanya memuat kode status
--     (`"0"`, `"1"`, `"2"`) DAN nama peran pada jalur balik tahap
--     (`"ReasLifeAdmin"`). Satu kolom, dua kosakata - dan itu disengaja,
--     sebab keduanya adalah "keadaan sebelum" dan "keadaan sesudah".
--   - AKUN_ID adalah identitas AKUN, bukan nama orang (ADR-U-0002).
CREATE TABLE {skema}.T_CLAIMLF_JEJAK (
  ID             VARCHAR2(40) NOT NULL,
  ADJUSTMENT_ID  VARCHAR2(40),
  KLAIM_ID       VARCHAR2(40),
  DARI           VARCHAR2(64),
  KE             VARCHAR2(64),
  AKUN_ID        VARCHAR2(64) NOT NULL,
  WAKTU          TIMESTAMP    NOT NULL,
  CONSTRAINT PK_CLAIMLF_JEJAK PRIMARY KEY (ID)
)
/
CREATE INDEX {skema}.IX_CLAIMLF_JEJAK_ADJ ON {skema}.T_CLAIMLF_JEJAK (ADJUSTMENT_ID)
/
CREATE INDEX {skema}.IX_CLAIMLF_JEJAK_KLAIM ON {skema}.T_CLAIMLF_JEJAK (KLAIM_ID)
/
CREATE SEQUENCE {skema}.SEQ_CLAIMLF_JEJAK START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
