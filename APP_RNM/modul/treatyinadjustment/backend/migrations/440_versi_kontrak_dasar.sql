-- VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR - rujukan ke versi yang menjadi dasarnya
--
-- Migrasi tiket 01 Treaty In Adjustment. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `D:\XML_NURE\_migration-docs\treaty-in-adjustment\5-tiket\issues\01-*.md`,
-- `treaty-in/SPEC-MODEL-DATA.md` §10.2, `GRL-10`, bahan to-spec `B-3` dan `B-4`.
--
-- ⛔ TABEL INI MILIK MODUL `treatyin` (migrasi 401, tiket 14). Modul ini
-- MENAMBAHKAN satu kolom padanya, dan itu sah: papan Adjustment dipisahkan dari
-- papan induk 25-09-2026 sementara MODEL DATANYA tetap satu - modul ini tidak
-- punya §10 sendiri. Kolom VERSI_KONTRAK karena itu digambarkan di
-- `modul/treatyin/docs/STRUKTUR-TABEL-TREATY-IN.md`, SATU dokumen, bukan dua
-- (dua dokumen atas tabel yang sama wajib identik - dan salinan yang wajib
-- identik adalah salinan yang akan menyimpang).
--
-- Urutan hulu ke hilir: 440 > 401, jadi tabelnya sudah berdiri saat berkas ini
-- berjalan.
--
-- Apa yang dasarnya: versi BERLAKU TERAKHIR pada saat versi baru dibuat.
-- BUKAN baris yang dipilih di layar, dan BUKAN versi `DISETUJUI` mana pun.
-- Sistem lama memakai `OLDID` = ID baris yang dipilih di picker, tanpa pembeda
-- jenis dan tanpa saringan keadaan (`TDA-11`); selisih yang dihitung terhadap
-- dasar yang salah TIDAK menghasilkan galat di sana.
--
-- ---------------------------------------------------------------------
-- CONSTRAINT YANG SENGAJA TIDAK DIPASANG - ini PERNYATAAN KEPUTUSAN
-- ---------------------------------------------------------------------
--   APA        : kolomnya NULLABLE. Tidak ada NOT NULL, tidak ada CHECK yang
--                memasangkan keterisiannya dengan NOMOR_URUT_VERSI, dan tidak
--                ada yang melarang menunjuk versi berkeadaan DITOLAK maupun
--                DIBATALKAN.
--   KENAPA     : ketiganya aturan bisnis, dan ADR-0056 (K-4) menahan aturan
--                bisnis di lapisan services - bukan di basis data. Bahan to-spec
--                `B-3` (wajib terisi pada versi penyesuaian, wajib kosong pada
--                versi pertama) dan `B-4` (yang ditunjuk harus versi berlaku
--                terakhir) menuntut mengetahui KEADAAN baris lain, dan constraint
--                Oracle tidak dapat menyatakannya tanpa trigger.
--                Lagi pula tiket 10 menyatakan baris warisan BOLEH berdasar
--                kosong sampai penomoran ulangnya selesai - NOT NULL di sini
--                membuat setiap baris warisan gagal dimuat.
--   AKIBAT     : sampai jalur simpan berdiri, kolom ini menerima pengenal versi
--                mana pun, termasuk yang keadaannya tidak sah sebagai dasar.
--   DITAGIH    : tiket lapisan aplikasi, yang menulis jalur simpannya; daftar
--                periksa tiket 01 menuntut keempat penolakan itu diuji di sana.
--
-- INV-18 - perilaku hapus DITETAPKAN SADAR. TANPA "ON DELETE": aturan bawaan
-- Oracle MENOLAK, dan menolak memang yang dikehendaki - kolom ini penunjuk ke
-- baris SEJARAH, bukan kepemilikan. Menghapus versi yang masih menjadi dasar
-- versi lain akan gagal dengan ORA-02292, bukan diam-diam memutus rantainya.
ALTER TABLE {skema}.VERSI_KONTRAK ADD (
  ID_VERSI_KONTRAK_DASAR  NUMBER(19)
)
/
ALTER TABLE {skema}.VERSI_KONTRAK ADD CONSTRAINT FK_VERSI_KONTRAK_DASAR
  FOREIGN KEY (ID_VERSI_KONTRAK_DASAR) REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
/
CREATE INDEX {skema}.IX_VERSI_KONTRAK_DASAR ON {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK_DASAR)
/
