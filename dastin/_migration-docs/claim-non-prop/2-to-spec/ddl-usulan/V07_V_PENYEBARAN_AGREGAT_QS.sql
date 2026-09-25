-- =====================================================================
-- V07_V_PENYEBARAN_AGREGAT_QS.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 31. Pengganti SpreadingAdjustmentQS. Hanya jenis reasuransi quota share.
-- =====================================================================

-- ---------------------------------------------------------------------
-- DUA NILAI PENYARING, DAN DAFTAR PENUHNYA TIDAK DIKETAHUI
-- ---------------------------------------------------------------------
-- '10028' dan '10004' EVIDENCED dari CountLossAllocation_act Langkah 10.
-- Apakah HANYA kedua nilai itu yang berarti quota share TIDAK DITEMUKAN.
--
-- Daftar penuhnya ada di POOLDATA.REINSURANCETYPE, dan tabel itu TIDAK
-- PUNYA SATU PUN CONSTRAINT — jadi bahkan di sana tidak ada yang menjamin
-- daftarnya tertutup.
--
-- Dituliskan sebagai literal DI SATU TEMPAT, bukan disebar, supaya bila
-- kelak nilai ketiga ditemukan ada tepat satu baris yang perlu diubah.
-- ---------------------------------------------------------------------

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

CREATE OR REPLACE VIEW KLAIMNP.V_PENYEBARAN_AGREGAT_QS AS
SELECT ID_KLAIM, JENIS_REASURANSI, MATA_UANG,
       SUM(NILAI_TERSEBAR)   AS NILAI_TERSEBAR,
       SUM(CASE WHEN KEADAAN_BARIS <> 'LENGKAP' THEN 1 ELSE 0 END) AS CACAH_BELUM_LENGKAP
  FROM KLAIMNP.PENYEBARAN
 WHERE JENIS_REASURANSI_ID IN ('10028', '10004')
 GROUP BY ID_KLAIM, JENIS_REASURANSI, MATA_UANG;

-- ADR-0017/0028: hilir membaca HANYA lewat view, tanpa hak atas tabel kanonik.
GRANT SELECT ON KLAIMNP.V_PENYEBARAN_AGREGAT_QS TO POOLDATA;
GRANT SELECT ON KLAIMNP.V_PENYEBARAN_AGREGAT_QS TO KLAIMNP_HILIR;
