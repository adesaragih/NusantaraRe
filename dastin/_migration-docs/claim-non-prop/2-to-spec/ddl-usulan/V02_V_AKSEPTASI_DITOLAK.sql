-- =====================================================================
-- V02_V_AKSEPTASI_DITOLAK.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 29. Kelompok yang DITOLAK V_AKSEPTASI_KOMPATIBEL, beserta sebabnya.
-- =====================================================================


-- ---------------------------------------------------------------------
-- Kenapa pengelompokan dan penyaringnya harus SAMA PERSIS dengan V01
-- ---------------------------------------------------------------------
-- Kedua view ini satu janji, bukan dua: "tidak ada yang hilang diam-diam".
-- Janji itu hanya berlaku bila setiap kelompok yang JATUH dari V01 MUNCUL
-- di sini. Begitu keduanya mengelompokkan dengan cara berbeda, ada kelompok
-- yang jatuh dari keduanya — dan itu keadaan terburuk: hilang, dan tidak
-- terlihat hilang.
--
-- Versi sebelumnya TIDAK sama, di dua tempat:
--   * penyambungannya melewatkan LAYER_BAGIAN dan LAYER_BAGIAN_JENIS,
--     sehingga ia melebur baris yang V01 pisahkan;
--   * ia tidak menyaring KEADAAN_BARIS, sementara V01 menyaring 'LENGKAP',
--     sehingga kelompok yang jatuh dari V01 karena barisnya belum lengkap
--     tidak pernah muncul di sini.
--
-- Diperbaiki 19 September 2026: penyambungan, penyaring, dan pengelompokan
-- kini identik. Yang berbeda hanya HAVING-nya — di sana sepakat, di sini
-- tidak sepakat — dan keduanya saling melengkapi tepat.


CREATE OR REPLACE VIEW KLAIMNP.V_AKSEPTASI_DITOLAK AS
SELECT a.ID_KLAIM,
       a.JENIS_REASURANSI,
       a.JENIS_REASURANSI_ID,
       a.MATA_UANG,
       COUNT(*) AS CACAH_BARIS,
       CASE WHEN MIN(al.LIMIT_LAYER)            <> MAX(al.LIMIT_LAYER)            THEN 'LIMIT_LAYER'
            WHEN MIN(al.PREMI_DEPOSIT)          <> MAX(al.PREMI_DEPOSIT)          THEN 'PREMI_DEPOSIT'
            WHEN MIN(al.PERSEN_PREMI_PEMULIHAN) <> MAX(al.PERSEN_PREMI_PEMULIHAN) THEN 'PERSEN_PREMI_PEMULIHAN'
            WHEN MIN(a.KURS)                    <> MAX(a.KURS)                    THEN 'KURS'
            ELSE 'PORSI_REASURADUR' END AS SEBAB
  FROM KLAIMNP.AKSEPTASI a
  LEFT JOIN KLAIMNP.ALOKASI_LAYER al
         ON al.ID_KLAIM = a.ID_KLAIM
        AND al.LAYER = a.LAYER
        AND al.LAYER_JENIS = a.LAYER_JENIS
        AND al.LAYER_BAGIAN = a.LAYER_BAGIAN
        AND al.LAYER_BAGIAN_JENIS = a.LAYER_BAGIAN_JENIS
        AND al.MATA_UANG = a.MATA_UANG
 WHERE a.KEADAAN_BARIS = 'LENGKAP'
   AND (al.ID_ALOKASI IS NULL OR al.KEADAAN_BARIS = 'LENGKAP')
 GROUP BY a.ID_KLAIM, a.JENIS_REASURANSI, a.JENIS_REASURANSI_ID, a.MATA_UANG
HAVING MIN(al.LIMIT_LAYER)            <> MAX(al.LIMIT_LAYER)
    OR MIN(al.PREMI_DEPOSIT)          <> MAX(al.PREMI_DEPOSIT)
    OR MIN(al.PERSEN_PREMI_PEMULIHAN) <> MAX(al.PERSEN_PREMI_PEMULIHAN)
    OR MIN(a.KURS)                    <> MAX(a.KURS)
    OR MIN(al.PORSI_REASURADUR)       <> MAX(al.PORSI_REASURADUR);


-- ---------------------------------------------------------------------
-- Satu keadaan yang TIDAK tertangkap view ini, dan dicatat terbuka
-- ---------------------------------------------------------------------
-- Kelompok yang seluruh barisnya ber-KEADAAN_BARIS bukan 'LENGKAP' tidak
-- muncul di V01 MAUPUN di sini — penyaringnya membuang barisnya lebih dulu,
-- jadi kelompoknya tidak pernah terbentuk.
--
-- Itu BUKAN kelalaian: kelompok begitu belum siap dikirim, dan alasannya
-- bukan "tidak sepakat" melainkan "belum lengkap" — sebab yang berbeda,
-- yang tempatnya di V_PARITAS_SHADOW dan di uji tiket 37, bukan di view
-- penolakan kontrak.
--
-- Dicatat di sini supaya yang membaca kedua view ini tahu pasangannya tidak
-- mencakup keadaan ketiga itu, alih-alih menyimpulkan sendiri bahwa
-- cakupannya lengkap.
-- ---------------------------------------------------------------------

-- ADR-0017/0028: hilir membaca HANYA lewat view, tanpa hak atas tabel kanonik.
GRANT SELECT ON KLAIMNP.V_AKSEPTASI_DITOLAK TO POOLDATA;
GRANT SELECT ON KLAIMNP.V_AKSEPTASI_DITOLAK TO KLAIMNP_HILIR;
