-- T_PREMIUM_LIST.STATUS_PENAWARAN - radio "Status" layar Input Offer.
--
-- `[keputusan work owner 01-10-2026]` status disimpan di header, dan header
-- menyimpan SATU BARIS PER STATUS untuk satu kasus:
--
--   baris utama   ID = ID_PEGA = nomor kasus  -> status TERAKHIR disimpan;
--                 inilah baris yang dibaca kotak masuk, Premium List Detail,
--                 penomoran PL, summary, dan Claim Life (kunci tidak berubah)
--   baris status  ID = 32 heksa (nomor kasus + status), ID_PEGA = nomor kasus
--                 -> salinan isian penawaran saat status itu terakhir disimpan
--
-- Simpan dengan status yang sama -> baris itu diperbarui; status berbeda ->
-- isian baris utama dipindah ke baris status lamanya, baris utama mengambil
-- status baru (repository/polis_penawaran.go `PindahkanBarisStatus`).
--
-- ⚠️ Di Pega radio ini (`pyWorkPage.Status` / `EmailTypePL`) tidak punya
-- kolom; nilainya hanya mengalir ke riwayat dan `SetStatusAkseptasi`.
-- Nilai yang disimpan: KODE radio (Accept/Pending/Bind/Decline/Closed untuk
-- bendera "0"; 1/2/3 untuk bendera "1").
--
-- Nullable; nol `COMMIT` (ADR-U-0029).
ALTER TABLE {skema}.T_PREMIUM_LIST ADD (
  STATUS_PENAWARAN VARCHAR2(255)
)
/
