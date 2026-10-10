-- NILAI_SELISIH dan NILAI_SEBELUM_PRO_RATE - PINDAH ke modul ini.
--
-- Migrasi tiket 76 dan 77, papan Treaty In Adjustment. Satu berkas = satu
-- langkah migrasi. Pernyataan dipisahkan oleh baris yang hanya berisi tanda
-- garis miring.
--
-- ⛔ BENTUKNYA DISALIN APA ADANYA dari `treatyin/.../426_nilai_selisih.sql`
-- - presisi `NUMBER(38,8)`, nama sequence (`SEQ_TRIA_*`), nama kunci asing,
-- nama index, urutan kolom. Migrasi yang membangun bentuk berbeda bukan
-- pemindahan, dan selisih bentuk antara yang dicabut dan yang dibangun
-- adalah selisih yang tidak akan pernah ada yang mencarinya.
--
-- Keputusan pemilik proses 4 Oktober 2026:
-- `modul/treatyin/docs/KEPUTUSAN-PENYELARASAN-REPO.md` §19. Pasangannya
-- `treatyin/.../435_cabut_nilai_selisih.sql`, dan keduanya WAJIB berjalan
-- bersama.
--
-- ---------------------------------------------------------------------
-- ⚠️ KETERGANTUNGAN URUTAN - `VERSI_KONTRAK` WAJIB SUDAH BERDIRI
-- ---------------------------------------------------------------------
--   Kedua tabel berkunci asing keluar ke `{skema}.VERSI_KONTRAK`, yang
--   dibuat `treatyin/.../401_kontrak_dan_versi_kontrak.sql` - MODUL LAIN.
--
--   Yang menjaganya nomor rentang: modul `treatyin` memakai `400`-`439`,
--   modul ini `440`+. Pelari migrasi berjalan menurut nomor, jadi `401`
--   selalu mendahului `442`. ⛔ Ketergantungan itu TIDAK dijaga penjaga
--   mana pun - ia dijaga oleh pilihan rentang, dan karena itu ditulis di
--   sini: siapa pun yang mengusulkan menomori ulang rentang modul ini
--   memutus kunci asing di bawah.
--
-- ---------------------------------------------------------------------
-- ALASAN RANCANGANNYA ADA DI `426`, DAN TIDAK DISALIN KE SINI
-- ---------------------------------------------------------------------
--   Kepala `426` memuat empat hal yang masih berlaku penuh, dan menyalinnya
--   ke sini hanya membuat salinannya membeku pada hari yang asli berubah:
--
--     * kardinalitas 1:N, dan kenapa ERD HTML menulis 1:1;
--     * penghalang kunci asing kedua - `BESARAN_DAPAT_DISESUAIKAN` tidak
--       ada di mana pun, jadi `KODE_BESARAN` berdiri sebagai TEKS;
--     * kenapa `NILAI_SEBELUM_PRO_RATE` berinduk `VERSI_KONTRAK` dan bukan
--       `KONTRAK` - kolom `MEMAKAI_PRORATA` ada di versi;
--     * INV-58 - tabel bernama `NILAI_SELISIH` TIDAK punya kolom selisih.
--
--   Keempatnya tetap mengikat di sini. `426` tidak dihapus dari riwayat
--   justru supaya rujukan ini tidak menggantung.
--
-- ---------------------------------------------------------------------
-- ⛔ PENJAGA `TestNolTabelBaru` MODUL INI DIUBAH, BUKAN DICABUT
-- ---------------------------------------------------------------------
--   Penjaga itu menyatakan *"modul ini hanya mengubah tabel milik
--   `treatyin`"*, dan ia MENOLAK berkas ini ketika pemindahan pertama kali
--   dicoba - riwayatnya tercatat di kepala `426`.
--
--   Keputusan §19 menggantikan aturan itu UNTUK DUA TABEL INI SAJA.
--   Penjaganya kini berupa daftar-izin: kedua tabel ini boleh, apa pun
--   selain keduanya tetap ditolak. Melonggarkannya menjadi "modul ini boleh
--   membuat tabel" akan membuang penjaga yang masih diperlukan.
--
--   ⚠️ DITAGIH: `treaty-in/SPEC-MODEL-DATA.md` masih menyatakan model
--   datanya SATU, dan penjaga lama mengutipnya. Pernyataan itu kini tidak
--   lagi benar seluruhnya, dan yang dapat memperbaikinya pemilik spec -
--   bukan berkas migrasi ini.
--
-- Pembalikan: `442_nilai_selisih_down.sql`, dan ia WAJIB berjalan bersama
-- `treatyin/.../435_cabut_nilai_selisih_down.sql`.
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
