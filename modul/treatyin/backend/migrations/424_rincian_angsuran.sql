-- RINCIAN_ANGSURAN - anak TERMIN yang tidak pernah ikut dibangun
--
-- Migrasi tiket 70 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal struktur: `ERD-TREATY-IN-DAN-EDM.html` baris relasi 27 -
-- `T_TREATY_INSTALLMENT 1:N T_TREATY_INSTALLMENT_ITEM` lewat `INSTALLMENT_ID`,
-- `CASCADE`. Nama kolom: `2-to-spec/KAMUS-KOLOM.md` §10.28.
--
-- Tiket 27 membangun `TERMIN` dan TIDAK membangun anaknya - bukan karena
-- diputuskan begitu, melainkan karena bentuk anaknya hanya terbaca lewat
-- salinan: ERD menandai buktinya DAUN-RELATIF dengan jalur
-- `TreatyIn.ValueDifference.Installment.InstallmentList`. Jadi daftar yang di
-- layar terlihat sebagai "termin" sebenarnya DUA tingkat, dan satu tingkat
-- penuh hilang dari rancangan tanpa pernah dinyatakan hilang.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS: IKUT HAPUS
-- ---------------------------------------------------------------------
--   `ERD-TREATY-IN-DAN-EDM.html` baris 27 menulis `CASCADE`, dan `ERD.md` §2.3
--   memberlakukan ikut hapus pada seluruh anak langsung baris versi. Keduanya
--   sepakat. Rincian angsuran tanpa terminnya tidak berarti apa pun.
--
-- ---------------------------------------------------------------------
-- EMPAT KOLOM, DAN KEEMPATNYA BELUM PUNYA PENULIS HIDUP
-- ---------------------------------------------------------------------
--   `DueDate`, `InstallmentPct`, `PaymentDate`, `WPC` hanya terbaca lewat
--   salinan `OLDDATA`/`ValueDifference`; nol aktivitas di ekspor 2026-09
--   menulisnya langsung. Bila ternyata ada kolom kelima, tabel ini bertambah
--   kolom - dan batas waktunya tiket 44, sebab menambah kolom sesudah data
--   masuk jauh lebih mahal. Penghalangnya tercatat di tiket 70.
--
--   ⚠️ `WPC` dibuat dengan NAMANYA APA ADANYA. Kepanjangannya tidak ada di
--   satu berkas pun - korpus sudah disapu habis (`CONTEXT.md` §2.10) - dan
--   menamainya dengan tebakan lebih buruk daripada menyimpannya apa adanya.
--
-- PERNYATAAN KEPUTUSAN - kunci alami (ID_TERMIN, NOMOR_URUT_RINCIAN) TIDAK
-- dipasang sebagai UNIQUE: `SPEC-INVARIAN.md` berhenti di INV-71. Aturan
-- `Z00_KUNCI_ALAMI.sql` berlaku. Tagihan di `SPEC-INVARIAN.md` §13.
--
-- INV-41: persentase di rentang 0-100 ditegakkan di services, bukan CHECK -
-- ADR-0056, basis data tidak membawa aturan bisnis.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
CREATE TABLE {skema}.RINCIAN_ANGSURAN (
  ID_RINCIAN_ANGSURAN  NUMBER(19)          NOT NULL,
  ID_TERMIN            NUMBER(19)          NOT NULL,
  NOMOR_URUT_RINCIAN   NUMBER(10)          NOT NULL,
  TANGGAL_JATUH_TEMPO  DATE,
  PERSEN_ANGSURAN      NUMBER(38,8),
  TANGGAL_BAYAR        DATE,
  WPC                  VARCHAR2(1000 CHAR),
  CONSTRAINT PK_RINCIAN_ANGSURAN PRIMARY KEY (ID_RINCIAN_ANGSURAN),
  CONSTRAINT FK_RINCIAN_ANGSURAN_1 FOREIGN KEY (ID_TERMIN)
    REFERENCES {skema}.TERMIN (ID_TERMIN) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_RINCIAN_ANGSURAN_TERMIN ON {skema}.RINCIAN_ANGSURAN (ID_TERMIN)
/
