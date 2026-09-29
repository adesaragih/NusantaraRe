-- T_CLAIMLF_PREMIUMLIST_DETAIL.STS_HAPUS - OQ-M6 DITUTUP 29-09-2026 (GILIRAN-17),
-- keputusan work owner: mencabut peserta = PENANDA, layar menyembunyikannya.
--
-- `[terverifikasi]` `Section/InputOSClaimLife.xml`: tombol `DELETE` b17865 (judul
-- kolom b16839) di grid `PremiumListDetail` b15924 menjalankan `deleteRow` b17874
-- tanpa konfirmasi (b18021), lalu `DeletePesertaClaimLife` b17909; tampil hanya
-- bila `CLAIM_NO == ''` (b18082). Activity itu (langkah 1 b227, 1.1, 1.2 `Obj-Save`
-- b443; `pyStepsBlockName` kosong b238/b326/b453) nol SQL - di Pega peserta
-- DILEPAS dari halaman kasus. ADR-U-0031: hapus = penanda, bukan hapus fisik.
--
-- ⛔ Tidak ada kolom penanda di 003 (`IS_CHECK` = dipilih untuk diklaim,
-- `STATUS`/`STS_REJECT` = status). Bentuknya mengikuti pola `STS_*` VARCHAR2(1)
-- teks kode (ADR-U-0022): NULL = aktif, '1' = dicabut. Nullable (ADR-U-0027),
-- sehingga seluruh baris yang ada tetap aktif tanpa pengisian ulang.
--
-- Setiap pembaca dan penulis tabel ini menyaring `STS_HAPUS IS NULL`; penjaganya
-- `TestSetiapPenyentuhTabelPesertaMenyaringPenandaCabut`.
--
-- ⛔ NOL COMMIT (ADR-U-0029). Dijalankan work owner sesudah giliran selesai.

ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ADD (
  STS_HAPUS VARCHAR2(1)
)
/
