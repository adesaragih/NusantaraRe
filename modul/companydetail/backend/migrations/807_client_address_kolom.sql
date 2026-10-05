-- 807 - Phone and Fax di CLIENT_ADDRESS (keputusan work owner 04-10-2026: "CLIENT_ADDRESS itu kan list bisa banyak
-- row"): SATU BARIS PER NOMOR. Alamat dengan tiga nomor = tiga baris beralamat sama; alamat tanpa nomor = satu baris
-- dengan ketiga kolom ini kosong. Tanpa tabel baru.
--   TELFAX_TYPE - kode M_ENUMERASI jenis telfax (3 MOBILE PHONE, 5 OFFICE PHONE, 6 EMAIL)
--   TELFAX_CODE - kode area, M_ENUMERASI jenis kodehp (maks. 4 karakter di data lama)
--   TELFAX_NO   - nomornya (maks. 37 karakter di data lama)
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.CLIENT_ADDRESS ADD (
  TELFAX_TYPE VARCHAR2(10),
  TELFAX_CODE VARCHAR2(10),
  TELFAX_NO   VARCHAR2(50)
)
/
