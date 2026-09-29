-- Kode bisnis pada header klaim - TEMUAN AUDIT A0.
--
-- Pemilik: A1. Bukan butir §1; ia lahir dari audit.
--
-- ⛔ Nomor akseptasi MEMUATNYA: `[terverifikasi]` `Claim Life/RDBList/
-- Generate_NoAccept_Life.xml` baris 85
-- `SELECT 'RNML-A'||{pyWorkPage.BusinessCode}||'.'||…`. Tetapi model
-- relasional TIDAK menyimpannya - `T_GENERAL_CLAIM` hanya punya
-- `BUSINESS_NAME`, dan `BUSINESSID` bukan salah satu dari 18 kolom datar
-- warisan yang `Simpan` tulis. Ia dipakai sekali saat pendaftaran lalu hilang.
--
-- Tanpa kolom ini jalur akseptasi gagal terang di
-- `ErrKodeBisnisBelumTersimpan` (HTTP 501).
--
-- ⛔ TEKS. Kode produk berawalan huruf (`L1`, `L2`, …) dan tidak pernah
-- menjadi bilangan (ADR-U-0022).
ALTER TABLE {skema}.T_GENERAL_CLAIM ADD (
  BUSINESS_CODE VARCHAR2(64)
)
/
