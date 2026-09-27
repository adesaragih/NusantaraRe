-- T_PREMIUM_LIST_SPREADING_RETRO - retrosesi per baris spreading (tiket 00).
--
-- ⛔ FK-nya menunjuk BARIS SPREADING (`SPREADING_ID`) - bukan header, bukan
-- peserta (AC 40). Tingkat KEEMPAT pohon, dan tingkatnya penting: retrosesi
-- milik satu treaty-year tertentu, bukan milik pesertanya secara keseluruhan.
--
-- ⚠️ Polis TANPA retrosesi pindah dengan NOL baris di sini. Itu keadaan
-- yang sah, bukan kegagalan migrasi (AC 42).
--
-- ⚠️ Nama tabelnya 30 byte PERSIS - batas Oracle di bawah 12.2. Nama
-- constraint dan index-nya karena itu DIPENDEKKAN (`PK_T_PL_SPR_RETRO`,
-- `FK_PLSR_SPR`); memakai nama tabel sebagai awalan akan melewati batas dan
-- ditolak ORA-00972.
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
CREATE TABLE {skema}.T_PREMIUM_LIST_SPREADING_RETRO (
  ID                     VARCHAR2(32) NOT NULL,
  SPREADING_ID           VARCHAR2(32),
  REINSURER_NAME         VARCHAR2(255),
  PERCENT_SHARE          NUMBER(38,8),
  AMOUNT                 NUMBER(38,8),
  COMMISION              NUMBER(38,8),
  OVR_COMM               NUMBER(38,8),
  RATE                   NUMBER(38,8),
  PREMIUM_SPREADED_GROSS NUMBER(38,8),
  PREMIUM_SPREADED_NET   NUMBER(38,8),
  TREATY_TYPE_ID         VARCHAR2(32),
  TREATY_TYPE_NAME       VARCHAR2(255),
  TREATY_START_DATE      DATE,
  TREATY_END_DATE        DATE,
  TGL_UPDATE             DATE,
  USER_ID                VARCHAR2(32),
  CONSTRAINT PK_T_PL_SPR_RETRO PRIMARY KEY (ID),
  CONSTRAINT FK_PLSR_SPR FOREIGN KEY (SPREADING_ID)
    REFERENCES {skema}.T_PREMIUM_LIST_SPREADING (ID) ON DELETE CASCADE
)
/

CREATE INDEX {skema}.IDX_PLSR_SPR ON {skema}.T_PREMIUM_LIST_SPREADING_RETRO (SPREADING_ID)
/
