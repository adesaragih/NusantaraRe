-- T_WORK_POLIS - work object polis, LINTAS-LINI (tiket 00 PremiumList Life).
--
-- Satu baris = satu work object polis, untuk SELURUH lini (Life dan
-- Non-Life). Ia BUKAN anak `T_PREMIUM_LIST`.
--
-- ⛔ KENAPA TERPISAH (AC 45, penyimpangan sadar 4). Keadaan tangga kerja -
-- posisi dan status - BUKAN atribut polis. Menaruhnya sebagai kolom header
-- polis membuat SETIAP VERSI polis membawa keadaan tangganya sendiri,
-- padahal tangganya satu. Polanya sama persis dengan `T_WORK_CLAIM`.
--
-- ⛔ HUBUNGAN 1:1 LEWAT SHARED PK, TANPA CONSTRAINT FK. `T_PREMIUM_LIST.ID`
-- bernilai SAMA PERSIS dengan `T_WORK_POLIS.ID`, dan STRUKTUR menyatakannya
-- "tanpa kolom penyambung"; kolom `Kunci` untuk `T_PREMIUM_LIST.ID` berisi
-- `PK` saja, bukan `PK, FK`. Menambahkan constraint di antara keduanya
-- berarti memutuskan sesuatu yang dokumen acuan tidak putuskan - dan arah
-- ketergantungannya pun belum ditetapkan (mana yang lahir lebih dulu).
--
-- ⚠️ `POSITION` berisi nilai connector VERBATIM - `Confirm` / `Decline` /
-- `Reject` / `Offer` / `Premium` - bukan angka. Angka menuntut peta kedua di
-- suatu tempat, dan peta kedua adalah tempat kedua untuk salah (ADR-U-0022).
--
-- ⛔ KOLOM AUDIT SENGAJA TIDAK ADA. Tiket 00 menyebut "audit" tanpa
-- menamainya, dan STRUKTUR menolak menuliskannya dengan alasan yang benar:
-- menamainya sendiri berarti mengarang. Ia masuk lampiran, bukan DDL ini.
--
-- ⛔ TIPE FISIKNYA KEPUTUSAN KAMI, BUKAN BACAAN DARI KORPUS.
-- STRUKTUR menyebut KATEGORI logis (teks / angka desimal / bilangan bulat /
-- DATE) dan menyatakan presisi fisik `[data DBA]`. Pemetaannya:
--
--   angka desimal   -> NUMBER(38,8)   uang dan share (ADR-0003, revisi-penyimpanan)
--   bilangan bulat  -> NUMBER(5)     umur, periode, PROD_KE, nomor baris
--                     (konvensi yang SUDAH ada di repo ini; penjaga
--                      TestNolNumberTanpaPresisi hanya menerima
--                      NUMBER(19), NUMBER(38,8), dan NUMBER(5))
--   DATE            -> DATE
--   teks            -> VARCHAR2(255); ID dan kolom ber-akhiran _ID -> VARCHAR2(32)
--
-- ⚠️ Ketiganya menunggu pencocokan DBA, dan kedua arah salahnya TIDAK
-- setara: VARCHAR2 yang terlalu pendek MENOLAK data yang sah - kegagalan yang
-- terlihat - sedangkan NUMBER yang terlalu pendek MEMBULATKAN uang diam-diam.
--
-- ⛔ SELURUH kolom nullable kecuali PK. STRUKTUR menyatakannya, dan
-- wajib-isi ditegakkan di Go. NOT NULL yang ditambahkan di sini akan menolak
-- baris warisan yang memang kosong saat migrasi data (tiket 09).
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). Batas transaksi milik Go.
--
-- ⚠️ Nama kolom di bawah DIBANGKITKAN dari STRUKTUR, tidak diketik ulang.
--
-- ⛔ RALAT 01-10-2026 - `STATUS` MENJADI `STATUS_WORK` (disunting DI TEMPAT).
-- `[keputusan work owner 01-10-2026]` ikuti struktur DEV: di sana kolom ini
-- sudah bernama STATUS_WORK (pola `T_WORK_CLAIM`), diubah di luar repo, dan
-- kode yang membaca `STATUS` gagal ORA-00904. Disunting di tempat supaya skema
-- BARU langsung benar dan penjaga DDL (yang tidak membaca RENAME) melihat nama
-- yang benar; lingkungan yang sempat memasang 050 lama diselaraskan migrasi
-- 063 (RENAME berpelindung katalog). DEV memuat pula COVER_KEY, CREATE_OP,
-- CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE (nullable) yang TIDAK dibuat di sini.
CREATE TABLE {skema}.T_WORK_POLIS (
  ID       VARCHAR2(32) NOT NULL,
  LINI     VARCHAR2(255),
  POSITION VARCHAR2(255),
  STATUS_WORK VARCHAR2(255),
  CONSTRAINT PK_T_WORK_POLIS PRIMARY KEY (ID)
)
/
