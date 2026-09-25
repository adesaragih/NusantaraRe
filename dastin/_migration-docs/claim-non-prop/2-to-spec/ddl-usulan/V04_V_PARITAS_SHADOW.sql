-- =====================================================================
-- V04_V_PARITAS_SHADOW.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 23. Perbandingan baris per baris terhadap sistem lama (ADR-0005).
-- Toleransi pembandingnya sendiri: AK-3. Pembandingnya BUKAN bagian view ini.
-- =====================================================================


-- ---------------------------------------------------------------------
-- Kenapa view ini menyebut KESEMBILAN BELAS tabel sasaran, bukan satu
-- ---------------------------------------------------------------------
-- Versi sebelumnya hanya menggabungkan AKSEPTASI. Akibatnya dua, dan
-- keduanya diam:
--
--   1. Baris korelasi untuk 18 tabel lain tetap muncul, tetapi
--      KEADAAN_BARIS-nya NULL — bukan karena barisnya belum lengkap,
--      melainkan karena view-nya tidak melihat tabelnya.
--   2. Penyaring "KEADAAN_BARIS = 'LENGKAP'" yang dijanjikan tiket 23
--      butir 2 akan MEMBUANG seluruh 18 tabel itu tanpa sepatah pesan.
--
-- Itu persis cacat yang ADR-0019 larang: NULL diperlakukan sebagai nilai.
-- Tiga keadaan dibedakan sekarang, dan tidak satu pun NULL:
--
--   KEADAAN_BARIS = 'LENGKAP' / 'MENUNGGU_KURS' / 'GAGAL_URAI'
--       tabelnya memang punya kolom keadaan, dan inilah isinya.
--   KEADAAN_BARIS = 'TIDAK_BERLAKU'
--       tabelnya TIDAK punya kolom keadaan — delapan tabel acuan dan
--       catatan. Berbeda dari "belum lengkap", dan sekarang terbaca begitu.
--   KEADAAN_BARIS = 'TABEL_TIDAK_DIKENALI'
--       baris korelasi menunjuk tabel yang tidak ada di daftar mana pun.
--       Tanpa cabang terakhir ini, baris seperti itu HILANG DARI VIEW —
--       dan baris yang hilang dari alat paritas adalah selisih yang tidak
--       pernah terhitung.
--
-- KEBERADAAN memisahkan pertanyaan kedua dari yang pertama: 'HILANG'
-- berarti baris kanoniknya tidak ada, apa pun keadaannya. Itu kegagalan
-- migrasi, bukan baris yang belum lengkap.
--
-- SETIAP BARIS KORELASI MUNCUL TEPAT SEKALI (tiket 23 butir 1): cabangnya
-- disaring TABEL_TUJUAN, dan ke-19 nilai itu saling lepas; cabang terakhir
-- menampung sisanya.
--
-- TIGA TABEL SENGAJA TIDAK ADA DI SINI: MIGRASI_KORELASI, MIGRASI_PENDARATAN,
-- dan MIGRASI_NILAI_DITOLAK. Ketiganya jembatan, bukan sasaran migrasi — sebuah
-- baris korelasi tidak pernah menunjuk tabel korelasi. Bila ternyata ada yang
-- menunjuk ke sana, cabang terakhir menangkapnya sebagai TABEL_TIDAK_DIKENALI,
-- dan itu memang yang seharusnya terjadi.


CREATE OR REPLACE VIEW KLAIMNP.V_PARITAS_SHADOW AS
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_KLAIM               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_KLAIM IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.KLAIM t ON t.ID_KLAIM = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'KLAIM'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_NILAI_KLAIM               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_NILAI_KLAIM IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.NILAI_KLAIM_MATA_UANG t ON t.ID_NILAI_KLAIM = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'NILAI_KLAIM_MATA_UANG'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_ALOKASI               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_ALOKASI IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.ALOKASI_LAYER t ON t.ID_ALOKASI = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'ALOKASI_LAYER'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_RETENSI               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_RETENSI IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.RETENSI_CEDANT t ON t.ID_RETENSI = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'RETENSI_CEDANT'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_AKSEPTASI               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_AKSEPTASI IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.AKSEPTASI t ON t.ID_AKSEPTASI = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'AKSEPTASI'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_ADJUSTMENT               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_ADJUSTMENT IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.ADJUSTMENT t ON t.ID_ADJUSTMENT = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'ADJUSTMENT'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_PREMI_PEMULIHAN               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_PREMI_PEMULIHAN IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.PREMI_PEMULIHAN t ON t.ID_PREMI_PEMULIHAN = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'PREMI_PEMULIHAN'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_OBJEK               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_OBJEK IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.OBJEK_PERTANGGUNGAN t ON t.ID_OBJEK = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'OBJEK_PERTANGGUNGAN'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_PEMBAGIAN               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_PEMBAGIAN IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.PEMBAGIAN_KERUGIAN t ON t.ID_PEMBAGIAN = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'PEMBAGIAN_KERUGIAN'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_PENYEBARAN               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_PENYEBARAN IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.PENYEBARAN t ON t.ID_PENYEBARAN = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'PENYEBARAN'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_ESTIMASI               AS ID_BARIS_KANONIK,
       t.KEADAAN_BARIS        AS KEADAAN_BARIS,
       CASE WHEN t.ID_ESTIMASI IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.ESTIMASI_AWAL t ON t.ID_ESTIMASI = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'ESTIMASI_AWAL'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_REKENING               AS ID_BARIS_KANONIK,
       'TIDAK_BERLAKU'        AS KEADAAN_BARIS,
       CASE WHEN t.ID_REKENING IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.REKENING_PENERIMA t ON t.ID_REKENING = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'REKENING_PENERIMA'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_KRONOLOGI               AS ID_BARIS_KANONIK,
       'TIDAK_BERLAKU'        AS KEADAAN_BARIS,
       CASE WHEN t.ID_KRONOLOGI IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.KRONOLOGI_KLAIM t ON t.ID_KRONOLOGI = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'KRONOLOGI_KLAIM'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_DOKUMEN               AS ID_BARIS_KANONIK,
       'TIDAK_BERLAKU'        AS KEADAAN_BARIS,
       CASE WHEN t.ID_DOKUMEN IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.DOKUMEN_KLAIM t ON t.ID_DOKUMEN = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'DOKUMEN_KLAIM'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_TUTUP_BUKU               AS ID_BARIS_KANONIK,
       'TIDAK_BERLAKU'        AS KEADAAN_BARIS,
       CASE WHEN t.ID_TUTUP_BUKU IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.TUTUP_BUKU t ON t.ID_TUTUP_BUKU = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'TUTUP_BUKU'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_TARIF               AS ID_BARIS_KANONIK,
       'TIDAK_BERLAKU'        AS KEADAAN_BARIS,
       CASE WHEN t.ID_TARIF IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.TARIF_BERLAKU t ON t.ID_TARIF = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'TARIF_BERLAKU'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_KOREKSI               AS ID_BARIS_KANONIK,
       'TIDAK_BERLAKU'        AS KEADAAN_BARIS,
       CASE WHEN t.ID_KOREKSI IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.KOREKSI_NILAI t ON t.ID_KOREKSI = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'KOREKSI_NILAI'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_PENJAGA_TANGGAL               AS ID_BARIS_KANONIK,
       'TIDAK_BERLAKU'        AS KEADAAN_BARIS,
       CASE WHEN t.ID_PENJAGA_TANGGAL IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.KLAIM_PENJAGA_TANGGAL t ON t.ID_PENJAGA_TANGGAL = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'KLAIM_PENJAGA_TANGGAL'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       t.ID_ARSIP               AS ID_BARIS_KANONIK,
       'TIDAK_BERLAKU'        AS KEADAAN_BARIS,
       CASE WHEN t.ID_ARSIP IS NULL THEN 'HILANG' ELSE 'ADA' END AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
  LEFT JOIN KLAIMNP.ARSIP_MUATAN_KELUAR t ON t.ID_ARSIP = mk.ID_TUJUAN
 WHERE mk.TABEL_TUJUAN = 'ARSIP_MUATAN_KELUAR'
UNION ALL
SELECT mk.ID_KORELASI, mk.TABEL_TUJUAN, mk.ID_TUJUAN,
       mk.PZINSKEY_LAMA, mk.PYID_LAMA, mk.CASEID_LAMA,
       mk.INDEX_OBJECT_LAMA, mk.URUTAN_BARIS_LAMA,
       mk.SUMBER_LAMA, mk.DIMIGRASI_PADA,
       CAST(NULL AS NUMBER(19)) AS ID_BARIS_KANONIK,
       'TABEL_TIDAK_DIKENALI'   AS KEADAAN_BARIS,
       'HILANG'                 AS KEBERADAAN
  FROM KLAIMNP.MIGRASI_KORELASI mk
 WHERE mk.TABEL_TUJUAN NOT IN ('KLAIM',
                               'NILAI_KLAIM_MATA_UANG',
                               'ALOKASI_LAYER',
                               'RETENSI_CEDANT',
                               'AKSEPTASI',
                               'ADJUSTMENT',
                               'PREMI_PEMULIHAN',
                               'OBJEK_PERTANGGUNGAN',
                               'PEMBAGIAN_KERUGIAN',
                               'PENYEBARAN',
                               'ESTIMASI_AWAL',
                               'REKENING_PENERIMA',
                               'KRONOLOGI_KLAIM',
                               'DOKUMEN_KLAIM',
                               'TUTUP_BUKU',
                               'TARIF_BERLAKU',
                               'KOREKSI_NILAI',
                               'KLAIM_PENJAGA_TANGGAL',
                               'ARSIP_MUATAN_KELUAR');


-- ---------------------------------------------------------------------
-- Apa yang view ini TIDAK lakukan
-- ---------------------------------------------------------------------
--   1. Ia TIDAK membandingkan nilai. Ia memasangkan baris. Pembanding dan
--      toleransinya AK-3, dan keduanya di luar lapisan data.
--   2. Ia TIDAK menjamin baseline sistem lama benar. AK-2b mencatat tiga
--      pengecualian bernama — TotalUR selalu nol, dan dua rumus premi
--      pemulihan bekerja atas nilai yang berbeda. Paritas terhadap angka
--      yang salah tetap paritas; itu sebabnya ketiganya bernama.
--   3. Ia TIDAK bertahan. Ia hidup selama MIGRASI_KORELASI hidup, dan
--      keduanya dijatuhkan setelah paritas diterima (ADR-0005).
-- ---------------------------------------------------------------------

-- ADR-0017/0028: hilir membaca HANYA lewat view, tanpa hak atas tabel kanonik.
GRANT SELECT ON KLAIMNP.V_PARITAS_SHADOW TO POOLDATA;
GRANT SELECT ON KLAIMNP.V_PARITAS_SHADOW TO KLAIMNP_HILIR;
