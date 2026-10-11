-- 484 - pengenal kasus Endorsement Life `EDML-<n>` (sama dengan Pega), sequence mulai 828
-- (keputusan work owner 07-10-2026). Dulu `EDMLF-<n>` (481/482) - tebakan analogi `NBLF-`.
--
-- a. Index unik `UX_PL_EDM_TERBUKA` (482, "satu endorsement terbuka per polis") dibuat ulang
--    dengan pola `EDML-%`. Nama index sama; ekspresi selain awalan tidak berubah. `LIKE 'EDML-%'`
--    tidak menangkap `NBLF-...` maupun `EDMLF-...` (karakter ke-5 `-` harfiah).
-- b. `SEQ_WORK_EDM_LIFE` (481) dimulai ulang pada 828: nomor `EDML-` tertinggi di Pega = 827
--    (diberikan work owner). Pola persis PremiumList 058.
--
-- ⚠️ Label:
-- [sementara - DBA memastikan nomor EDML- tertinggi Pega sebelum data nyata]
-- 827 datang dari work owner, bukan dari audit penghitung Pega (`PC_DATA_UNIQUEID`).
--
-- ⚠️ Kasus uji lama `EDMLF-<n>` di DEV TIDAK disentuh (nol isi mundur, nol hapus) - work owner
-- menghapusnya sendiri. Selama masih ada, baris itu tidak lagi terlindung index terbuka-ganda dan
-- tidak dikenali kode sebagai kasus EDM.
--
-- ⛔ DROP lalu CREATE, bukan ALTER (`ALTER SEQUENCE … RESTART` baru ada sejak 18c) - sama dengan 058.
-- ⚠️ BILA GAGAL DI TENGAH (DROP sudah jalan, CREATE belum): DDL Oracle menetapkan dirinya sendiri;
-- jalankan pernyataan CREATE yang tertinggal secara manual (DBA), lalu ulangi `-migrate`.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
DROP INDEX {skema}.UX_PL_EDM_TERBUKA
/
CREATE UNIQUE INDEX {skema}.UX_PL_EDM_TERBUKA ON {skema}.T_PREMIUM_LIST (
  CASE WHEN STATUSS IS NULL AND ID LIKE 'EDML-%' THEN OLD_POLICY_NO END
)
/
DROP SEQUENCE {skema}.SEQ_WORK_EDM_LIFE
/
CREATE SEQUENCE {skema}.SEQ_WORK_EDM_LIFE START WITH 828 INCREMENT BY 1 NOCACHE NOCYCLE
/
