-- T_CLAIMLF_PREMIUMLIST_DETAIL.IDPEGA - keputusan work owner 10-10-2026 ("iya tambahkan"), migrasi disetujui.
--
-- IDPEGA = M_LIFE_PREMIUM_DETAIL.IDPEGA baris sumber peserta (pengenal kasus versi polis, NBLF-/EDML-): kunci header
-- data polis klaim (baris utama T_PREMIUM_LIST ber-ID = ID_PEGA = IDPEGA; keputusan work owner 10-10-2026 "ambil ke
-- tabel T_PREMIUM_LIST berdasarkan idpega peserta nya saja"). Lebar = sumbernya: VARCHAR2(50) (katalog
-- M_LIFE_PREMIUM_DETAIL kolom 80, `docs/KATALOG-TABEL-PESERTA-DAN-TREATY.md`). Nullable (ADR-U-0027).
--
-- Pendaftaran menulisnya bersama peserta (kolomSalin -> kolomPeserta); pembaca tahap lanjut memakai nilai tersimpan.
--
-- ⛔ Kolom ditambah dengan ALTER biasa (penjaga STRUKTUR membaca jalur maju; blok berpelindung untuk kolom baru
-- ditolak penjaga inti). Idempotensi dijaga pelari (T_MIGRASI).
--
-- Isi ulang (backfill) baris lama dari baris sumbernya, berkunci PL_NUMBER + CERTIFICATE_NO (index tabel 66,8 juta
-- baris) + ID = SOURCE_ID; aman diulang (`IDPEGA IS NULL`). Baris tanpa SOURCE_ID (pindahan warisan) atau tanpa baris
-- sumber tetap NULL - pembaca memakai cadangan lewat SOURCE_ID / melaporkannya.
--
-- ⛔ NOL COMMIT (ADR-U-0029). Dijalankan work owner (-migrate), bukan oleh giliran ini.

ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ADD (
  IDPEGA VARCHAR2(50)
)
/
UPDATE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL d
   SET d.IDPEGA = (SELECT m.IDPEGA FROM {skema}.M_LIFE_PREMIUM_DETAIL m
                    WHERE m.PL_NUMBER = d.PL_NUMBER AND m.CERTIFICATE_NO = d.CERTIFICATE_NO AND m.ID = d.SOURCE_ID
                      AND ROWNUM = 1)
 WHERE d.IDPEGA IS NULL AND d.SOURCE_ID IS NOT NULL
/
