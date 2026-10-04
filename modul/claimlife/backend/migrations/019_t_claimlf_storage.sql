-- T_CLAIMLF_STORAGE - kartu berkas unggahan (butir be, 27-09-2026).
--
-- ⛔ TABEL KAMI SENDIRI, dan itu keputusan yang dinyatakan. Padanannya di
-- sistem lama adalah `T_STORAGE_IMAGE`, yang dibaca `GetLinkStorage_SQL` b85
-- dan ditulis `Insert_T_Storage_SQL` b85-103 / `Update_T_Storage_SQL` b86:
--
--   IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE
--
-- Kolomnya ditiru NAMA DEMI NAMA supaya migrasi data kelak mekanis. Yang
-- TIDAK ditiru: tempatnya. `T_STORAGE_IMAGE` tabel BERSAMA lintas modul dan
-- lintas aplikasi (`APPNAME` ada justru karena itu); menulis ke sana berarti
-- aplikasi ini menaruh baris di meja orang lain, dan brief menuntut
-- persetujuan manusia sebelum penyambungan penyimpanan nyata. Tabel sendiri
-- membuat jalur DEV berjalan penuh tanpa menyentuh milik siapa pun.
--
-- ⛔ DAN `Insert_T_Storage_SQL` b102 BER-`commit;`. ADR-U-0029 melarangnya:
-- transaksi milik pemanggil, bukan milik pernyataan. Itu sebab kedua tabel
-- ini tidak ditulis lewat rule warisan.
--
-- ⚠️ NOL kunci tamu ke `T_CLAIMLF_DOCUMENT`, dan itu meniru aslinya:
-- `T_STORAGE_ID` di tabel dokumen menunjuk `IMAGEID` TANPA constraint, sebab
-- di sistem lama barisnya lahir di layanan luar dan boleh mendahului atau
-- menyusul barisnya. Menambahkan FK di sini akan menolak urutan yang sah.
--
-- ⚠️ `URLPUBLIC` 2000: alamat berkas, bukan kalimat. Di DEV isinya
-- `/api/dokumen/{id}/isi`; di produksi kelak URL penyimpanan luar, dan URL
-- bertanda tangan bisa panjang.
--
-- ⛔ NOL COMMIT (ADR-U-0029). Identitas dari `IMAGEID` yang diturunkan
-- `models.IDDokumenBaru`, bukan dari sequence: ia HARUS sama dengan nilai
-- yang `T_CLAIMLF_DOCUMENT.T_STORAGE_ID` simpan, dan di sistem lama nilai itu
-- cap waktu - bukan nomor kita.

CREATE TABLE {skema}.T_CLAIMLF_STORAGE (
  IMAGEID    VARCHAR2(64) NOT NULL,
  URLPUBLIC  VARCHAR2(2000),
  APPFOLDER  VARCHAR2(255),
  EXPDATE    DATE,
  FILENAME   VARCHAR2(255),
  APPNAME    VARCHAR2(20),
  STORAGE    VARCHAR2(32),
  CONSTRAINT PK_T_CLAIMLF_STORAGE PRIMARY KEY (IMAGEID)
)
/
