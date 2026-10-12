-- T_CLAIMLF_PREMIUMLIST_DETAIL - indeks pemeriksa peserta aktif di klaim lain (aturan in-app Register,
-- repository/pesertainapp.go `PesertaAktifDiKlaimLain`). Keputusan work owner 10-10-2026: "ya buat saja indeks nya".
--
-- ⛔ INDEKS FUNGSI, ekspresinya PERSIS kueri: `WHERE TRIM(p.PL_NUMBER) = :3 AND TRIM(p.CERTIFICATE_NO) IN (...)`.
-- Indeks biasa (PL_NUMBER, CERTIFICATE_NO) tidak terpakai oleh kolom yang dibungkus TRIM, dan TRIM tidak dapat
-- dibuang dari kueri: penulisnya tidak menjamin nilai ter-trim - `pohonklaim.go` `Simpan` menyalin PL_NUMBER dan
-- CERTIFICATE_NO apa adanya dari M_LIFE_PREMIUM_DETAIL (`kolomSalin`), dan pindahan warisan menyalin
-- OS_AKSEPTASI_KLAIM_LIFE apa adanya. Penjaga `TestIndeksPesertaAktifCocokDenganKueri` mengunci kesamaan ekspresinya.
--
-- Subkueri baris adjustment (`a.PREMIUM_LIST_DETAIL_ID = p.ID`) sudah ber-indeks IX_ADJ_PLD_ID (migrasi 004); join
-- header `h.ID = p.CLAIM_ID` memakai PK T_GENERAL_CLAIM - tidak ada indeks lain yang perlu.
--
-- Idempoten: blok berpelindung ALL_INDEXES (pola 100_relasi_lima_tabel_life). NOL COMMIT (ADR-U-0029). Dijalankan
-- work owner (-migrate), bukan oleh giliran ini.

DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_CLAIMLF_PREMIUMLIST_DETAIL' AND INDEX_NAME = 'IX_PLD_TRIM_PL_SERT';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'CREATE INDEX {skema}.IX_PLD_TRIM_PL_SERT ON {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL (TRIM(PL_NUMBER), TRIM(CERTIFICATE_NO))';
  END IF;
END;
/
