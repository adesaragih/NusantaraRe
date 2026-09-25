-- =====================================================================
-- V03_V_REKAP_KLAIM_MATA_UANG.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 30. Rekap per klaim per mata uang. Retensi Cedant TIDAK PERNAH tercampur.
-- =====================================================================

-- ---------------------------------------------------------------------
-- CACAH_BELUM_LENGKAP — kenapa setiap view agregat di sini membawanya
-- ---------------------------------------------------------------------
-- Sebuah penjumlahan tidak menyatakan berapa banyak baris yang IKUT
-- terjumlah, dan tidak menyatakan apakah baris-baris itu sudah selesai.
-- Dua kelompok dengan angka yang sama persis dapat berarti hal yang sangat
-- berbeda: yang satu lengkap, yang lain separuhnya masih menunggu kurs.
--
-- ADR-0019 melarang nilai yang belum dihitung terbaca sebagai nol. Di
-- tingkat baris, larangan itu ditegakkan kolom KEADAAN_BARIS. Di tingkat
-- AGREGAT, larangan itu hilang begitu SUM dijalankan — kecuali agregatnya
-- ikut membawa keterangan itu keluar.
--
-- CACAH_BELUM_LENGKAP = 0 berarti seluruh baris kelompok ini 'LENGKAP'.
-- Lebih dari nol berarti angkanya BELUM BOLEH DIPAKAI sebagai final, dan
-- pembacanya dapat memutuskan itu sendiri tanpa membuka tabel sumbernya.
--
-- Baris yang belum lengkap TETAP IKUT DIJUMLAH, dan itu disengaja: nilai
-- dalam mata uang aslinya sudah benar sejak awal; yang belum ada hanya
-- padanan rupiahnya. Membuangnya akan menghasilkan angka yang salah
-- diam-diam — persis cacat yang sama dari arah berlawanan.
--
-- Tidak ada nilai IDR yang dijumlahkan di view mana pun di berkas ini.
-- Bila kelak ada yang menambahkannya, CACAH_BELUM_LENGKAP adalah penjaga
-- yang membuat kesalahan itu terbaca alih-alih tersembunyi.
-- ---------------------------------------------------------------------

CREATE OR REPLACE VIEW KLAIMNP.V_REKAP_KLAIM_MATA_UANG AS
SELECT kl.ID_KLAIM,
       n.MATA_UANG,
       (SELECT SUM(x.NILAI_TERALOKASI_DIPAKAI) FROM KLAIMNP.ALOKASI_LAYER x
         WHERE x.ID_KLAIM = kl.ID_KLAIM AND x.MATA_UANG = n.MATA_UANG) AS NILAI_TERALOKASI_LAYER,
       (SELECT SUM(r.NILAI_RETENSI_DIPAKAI) FROM KLAIMNP.RETENSI_CEDANT r
         WHERE r.ID_KLAIM = kl.ID_KLAIM AND r.MATA_UANG = n.MATA_UANG) AS NILAI_RETENSI_CEDANT,
       (SELECT SUM(p.PREMI_PEMULIHAN) FROM KLAIMNP.PREMI_PEMULIHAN p
         WHERE p.ID_KLAIM = kl.ID_KLAIM AND p.MATA_UANG = n.MATA_UANG) AS PREMI_PEMULIHAN,
       (SELECT COUNT(*) FROM KLAIMNP.ALOKASI_LAYER x
         WHERE x.ID_KLAIM = kl.ID_KLAIM AND x.MATA_UANG = n.MATA_UANG
           AND x.KEADAAN_BARIS <> 'LENGKAP')
     + (SELECT COUNT(*) FROM KLAIMNP.RETENSI_CEDANT r
         WHERE r.ID_KLAIM = kl.ID_KLAIM AND r.MATA_UANG = n.MATA_UANG
           AND r.KEADAAN_BARIS <> 'LENGKAP')                           AS CACAH_BELUM_LENGKAP
  FROM KLAIMNP.KLAIM kl
  JOIN KLAIMNP.NILAI_KLAIM_MATA_UANG n ON n.ID_KLAIM = kl.ID_KLAIM;

-- ---------------------------------------------------------------------
-- Kenapa TIGA subkueri terpisah, bukan satu penjumlahan
-- ---------------------------------------------------------------------
-- Inilah seluruh isi ADR-0010 sebagai SQL. Di sistem lama, Retensi Cedant
-- adalah BARIS di dalam daftar alokasi yang sama — sehingga setiap loop atas
-- SpreadingRisk ikut melihatnya KECUALI menyaringnya sendiri, dan setiap
-- penulis loop harus ingat melakukannya (BLUEPRINT 2.4).
--
-- Di sini ia TABEL TERSENDIRI dan KOLOM TERSENDIRI. Menjumlahkan alokasi
-- layer tidak dapat memuat Retensi Cedant — bukan karena ada penyaring yang
-- harus diingat, melainkan karena barisnya tidak ada di sana.
--
-- SUM atas nol baris menghasilkan NULL, bukan 0, dan itu DIPERTAHANKAN:
-- NULL berarti "tidak ada baris untuk dijumlah", dan itu berbeda dari nol
-- yang berarti "sudah dihitung, hasilnya nol" (ADR-0019).
--
-- Baris Retensi Cedant lama membawa TotalClaim yang SELALU NOL pada cabang
-- mata uang sama (D7). Sistem baru menghitungnya seperti Retensi Cedant
-- biasa (AK-2b), jadi selisih terhadap sistem lama DIHARAPKAN dan masuk
-- pengecualian bernama shadow-run — bukan kegagalan.
-- ---------------------------------------------------------------------

-- ADR-0017/0028: hilir membaca HANYA lewat view, tanpa hak atas tabel kanonik.
GRANT SELECT ON KLAIMNP.V_REKAP_KLAIM_MATA_UANG TO POOLDATA;
GRANT SELECT ON KLAIMNP.V_REKAP_KLAIM_MATA_UANG TO KLAIMNP_HILIR;
