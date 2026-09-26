-- Kolom isi T_CLAIMLF_DOCUMENT - dokumen pendukung per peserta.
--
-- Pemilik: tiket 03 (AC 16-18).
--
-- `[keputusan work owner 26-09-2026, butir ad]` — dan hasil sensusnya BUKAN
-- yang dikira brief. Sensus yang diminta adalah `.DocumentList` di
-- `SaveOutStandingLife_Act`; yang ditemukan di sana hanya SATU rujukan, dan
-- kelasnya `Link-Attachment` — lampiran bawaan Pega, bukan tabel karangan.
--
-- Kolom tabel dokumen datang dari tempat lain: `InsertDocument_Act` milik kelas
-- `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`. Ia memakai `Obj-Save`, yang menulis
-- SELURUH properti kelas itu, jadi tidak ada daftar kolom pada pernyataan SQL
-- mana pun — tetapi ADA daftar Property-Set yang mengisinya, dan itulah yang
-- disensus.
--
-- ⭐ Cacahnya DUA KALI, dari sumber yang sungguh berbeda (CLAUDE.md §4a):
--   1. katalog instance pengembangan: tabel warisan `DOCUMENT_CLAIM` punya 14 kolom
--      (brief ronde 4 §9);
--   2. seluruh Property-Set `InsertDocument_Act` sebelum `Obj-Save`
--      `[terverifikasi]` baris 647-940 berkas pecahan: .ID, .TANGGAL, .IDPEGA,
--      .NAMAFILE, .MIME, .KATEGORI_1, .KATEGORI_2, .NOAKSEP, .NOPREKAS,
--      .PAYMENTDATE, .pxCreateOperator, ditambah .T_STORAGE_ID yang diisi
--      activity storage dan menjadi precondition Obj-Save-nya — 12 properti.
-- Kedua belas properti itu seluruhnya termuat di keempat belas kolom katalog;
-- dua sisanya (INSKEY_LINK, INSKEY_DATA) memang tidak pernah disentuh rule
-- ini. Kedua cacah COCOK, dan itulah yang membuat daftar di bawah bukan
-- tebakan.
--
-- ⚠️ Yang TIDAK ikut dan sebabnya:
--   - IDPEGA, INSKEY_LINK, INSKEY_DATA, PXCREATEOPERATOR : milik Pega, bukan
--     milik data dokumen. Sistem baru bukan Pega, dan kolom yang isinya kunci
--     internal mesin lama tidak berarti apa-apa di sini.
--   - NOAKSEP, NOPREKAS : nomor akseptasi dan prekas hidup di baris
--     adjustment, bukan di dokumen. Menyalinnya ke sini membuat satu nilai
--     punya dua rumah.
--
-- ⭐ T_STORAGE_ID adalah PENUNJUK ke berkasnya di Google Storage (ADR-U-0010) -
-- isi berkas tidak pernah masuk basis data.
ALTER TABLE {skema}.T_CLAIMLF_DOCUMENT ADD (
  NAMA_FILE     VARCHAR2(255),
  MIME          VARCHAR2(128),
  KATEGORI_1    VARCHAR2(64),
  KATEGORI_2    VARCHAR2(64),
  TANGGAL       DATE,
  T_STORAGE_ID  VARCHAR2(64),
  PAYMENT_DATE  DATE
)
/
