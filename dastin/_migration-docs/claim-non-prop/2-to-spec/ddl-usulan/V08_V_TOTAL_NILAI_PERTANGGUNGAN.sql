-- =====================================================================
-- V08_V_TOTAL_NILAI_PERTANGGUNGAN.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 31. Pengganti TotalInterestInsured. TANPA tabel penyimpan.
-- =====================================================================

-- Di sistem lama, salah satu rutin dedup Java MEMBUANG BARIS dari
-- .ClaimData.TotalInterestInsured berdasarkan (CurrencyID, Currency) — data
-- klaim, tanpa pencatatan, tanpa penanda, tanpa jejak (D40). Sebagai view,
-- tidak ada baris yang dapat dibuang: yang tampak selalu seluruh isinya.

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

CREATE OR REPLACE VIEW KLAIMNP.V_TOTAL_NILAI_PERTANGGUNGAN AS
SELECT ID_KLAIM, MATA_UANG,
       SUM(NILAI_PERTANGGUNGAN) AS NILAI_PERTANGGUNGAN,
       SUM(CASE WHEN KEADAAN_BARIS <> 'LENGKAP' THEN 1 ELSE 0 END) AS CACAH_BELUM_LENGKAP
  FROM KLAIMNP.OBJEK_PERTANGGUNGAN
 GROUP BY ID_KLAIM, MATA_UANG;

-- ADR-0017/0028: hilir membaca HANYA lewat view, tanpa hak atas tabel kanonik.
GRANT SELECT ON KLAIMNP.V_TOTAL_NILAI_PERTANGGUNGAN TO POOLDATA;
GRANT SELECT ON KLAIMNP.V_TOTAL_NILAI_PERTANGGUNGAN TO KLAIMNP_HILIR;
