-- T_VIEW_SUGGEST - riwayat penawaran (tiket 00 PremiumList Life).
--
-- RALAT 28-09-2026 (pl6) - KOLOM INITIAL DIGANTI NAMA.
--
-- Migrasi ini GAGAL di DEV saat work owner menjalankan -migrate 09.33:
-- INITIAL adalah KATA CADANGAN Oracle. Bukti yang dijalankan work owner:
--
--   SELECT 1 AS INITIAL   FROM DUAL  -> ORA-00923
--   SELECT 1 AS "INITIAL" FROM DUAL  -> lolos (berkutip)
--   SELECT 1 AS NO        FROM DUAL  -> lolos
--
-- Sebabnya BUKAN penanda {skema}: sebelas langkah lain di jalannya yang
-- sama lolos dengan penanda itu. Nama baru INITIAL_SUGGEST mengikuti pola
-- saudaranya DATE_SUGGEST, PIC_SUGGEST, COMMENT_SUGGEST.
--
-- Berkas ini disunting DI TEMPAT, bukan ditambah migrasi 057: ia belum
-- pernah terpasang di mana pun - T_MIGRASI DEV tidak memuat 056, dan skema
-- uji belum pernah dibuat. Menambah ALTER RENAME untuk kolom yang belum
-- pernah ada berarti mewariskan riwayat yang tidak terjadi.
--
-- Penjaga TestNolKataCadanganOracleSebagaiKolom lahir bersama ralat ini:
-- yang membuat INITIAL lolos bukan kecerobohan, melainkan ketiadaan
-- pemeriksa.
--
-- Satu baris = satu langkah riwayat penawaran polis.
--
-- ⚠️ Sumber warisannya `JSON_OFFER_LIFE` (CLOB + 24 kolom flat), ditulis
-- `SaveOfferJsonLife_SQL`. CLOB-nya DIBUANG; yang pindah adalah kolomnya.
--
-- ⛔ `[terverifikasi]` `GetOfferLife_sql` menyaring
-- `WHERE STATUS = 'Bind' AND OLDID IS NULL`. Adanya `OLDID` menandakan
-- VERSIONING penawaran. Apakah versi lama ikut dibawa adalah keputusan
-- migrasi data (tiket 09) - dilaporkan, tidak ditebak. Karena itu TIDAK ADA
-- kolom `OLD_ID` yang dikarang di sini.
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
CREATE TABLE {skema}.T_VIEW_SUGGEST (
  ID                VARCHAR2(32) NOT NULL,
  PREMIUM_LIST_ID   VARCHAR2(32),
  NO                NUMBER(5),
  DATE_SUGGEST      DATE,
  PIC_SUGGEST       VARCHAR2(255),
  IS_CEDING_CONFIRM VARCHAR2(255),
  COMMENT_SUGGEST   VARCHAR2(255),
  INITIAL_SUGGEST   VARCHAR2(255),
  CONSTRAINT PK_T_VIEW_SUGGEST PRIMARY KEY (ID),
  CONSTRAINT FK_VS_PL FOREIGN KEY (PREMIUM_LIST_ID)
    REFERENCES {skema}.T_PREMIUM_LIST (ID) ON DELETE CASCADE
)
/

CREATE INDEX {skema}.IDX_VS_PL ON {skema}.T_VIEW_SUGGEST (PREMIUM_LIST_ID)
/
