-- T_PREMIUM_LIST_SPREADING - spreading peserta (tiket 00 PremiumList Life).
--
-- Satu baris = satu treaty-year dari satu peserta. Nilainya DIBEKUKAN saat
-- polis disimpan.
--
-- ⛔ FK-nya menunjuk PESERTA (`DETAIL_ID`), bukan header (AC 40). Spreading
-- yang menggantung pada header kehilangan peserta pemiliknya, dan jumlahnya
-- tidak lagi dapat dicocokkan per peserta saat rekonsiliasi migrasi.
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
CREATE TABLE {skema}.T_PREMIUM_LIST_SPREADING (
  ID               VARCHAR2(32) NOT NULL,
  DETAIL_ID        VARCHAR2(32),
  TREATY_TYPE_ID   VARCHAR2(32),
  TREATY_TYPE_NAME VARCHAR2(255),
  TREATY_YEAR_LIFE VARCHAR2(255),
  RETROCADED_SHARE NUMBER(38,8),
  IDR              NUMBER(38,8),
  USD              NUMBER(38,8),
  IDR_SELISIH      NUMBER(38,8),
  USD_SELISIH      NUMBER(38,8),
  B_IDR            NUMBER(38,8),
  B_USD            NUMBER(38,8),
  TGL_UPDATE       DATE,
  USER_ID          VARCHAR2(32),
  CONSTRAINT PK_T_PL_SPREADING PRIMARY KEY (ID),
  CONSTRAINT FK_PLS_DETAIL FOREIGN KEY (DETAIL_ID)
    REFERENCES {skema}.T_PREMIUM_LIST_DETAIL (ID) ON DELETE CASCADE
)
/

CREATE INDEX {skema}.IDX_PLS_DETAIL ON {skema}.T_PREMIUM_LIST_SPREADING (DETAIL_ID)
/
