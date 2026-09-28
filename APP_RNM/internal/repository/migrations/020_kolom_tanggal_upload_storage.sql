-- T_CLAIMLF_STORAGE.TANGGAL_UPLOAD - butir be, ralat 28-09-2026.
--
-- ⛔ KOLOM YANG TERLEWAT DI 019, dan buktinya rule yang menulisnya:
--
--   RDBList/Update_T_Storage_SQL.xml b86-b90
--     update T_STORAGE_IMAGE set
--       URLPUBLIC      = {UpdateDoc.URLImage},
--       APPFOLDER      = {UpdateDoc.appfolder},
--       EXPDATE        = To_date({UpdateDoc.exp},      'DD/MM/YYYY HH24:MI:SS'),
--       TANGGAL_UPLOAD = To_date({UpdateDoc.DateTime}, 'MM/DD/YYYY HH24:MI:SS')
--     where imageid = {UpdateDoc.ImageID}
--
-- Migrasi 019 dibuat dari `Insert_T_Storage_SQL.xml` saja, yang memang tidak
-- menyebut kolom ini - ia ditulis oleh UPDATE, bukan oleh INSERT. Membaca
-- satu dari dua rule penulis lalu menyimpulkan tentang tabelnya: bentuk
-- kekeliruan yang sama untuk kelima kalinya di modul ini.
--
-- ⚠️ DUA BENTUK TANGGAL YANG BERBEDA DI SATU PERNYATAAN, dan itu ada di
-- rule aslinya - bukan salah salin kami:
--
--   EXPDATE        'DD/MM/YYYY HH24:MI:SS'   <- hari dulu
--   TANGGAL_UPLOAD 'MM/DD/YYYY HH24:MI:SS'   <- bulan dulu
--
-- Akibatnya nilai warisan pada kedua kolom itu TIDAK dapat dibedakan untuk
-- tanggal 1-12 tiap bulan. Kami tidak mewarisi masalahnya - Go mengikat
-- `time.Time`, bukan teks - tetapi migrasi data (A4) harus tahu, dan itu
-- sebab catatan ini berdiri di sini.
--
-- ⛔ Sumbernya BUKAN jam kami: `GetUrlGoogleStorage_Act.xml` b2273-2274
-- menyetel `UpdateDoc.DateTime = UploadDoc.Response.DateTime`, yaitu cap
-- waktu yang DIKEMBALIKAN layanan penyimpanan. Selama pelaksananya stub,
-- layanan itu adalah proses kita sendiri, jadi jamnya jam kita - penyimpangan
-- yang dinyatakan, bukan disembunyikan.
--
-- ⛔ NOL COMMIT (ADR-U-0029).

ALTER TABLE {skema}.T_CLAIMLF_STORAGE ADD (
  TANGGAL_UPLOAD DATE
)
/
