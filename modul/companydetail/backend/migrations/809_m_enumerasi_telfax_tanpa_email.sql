-- 809 - Phone and Fax tanpa pilihan EMAIL (perintah work owner 04-10-2026: "Phone and Fax ada opsi EMAIL, hapus
-- aja"). Dropdown kini hanya 3 MOBILE PHONE dan 5 OFFICE PHONE.
--
-- Barisnya TIDAK dibuang, hanya AKTIF = 0: nomor lama berjenis EMAIL (DEV: 39 di dokumen Pega) tetap terbaca di
-- layar sebagai nilai lama. Berkas 801 tidak disunting - ia sudah dijalankan.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
UPDATE {skema}.M_ENUMERASI SET AKTIF = '0'
WHERE JENIS = 'telfax' AND KODE = '6'
/
