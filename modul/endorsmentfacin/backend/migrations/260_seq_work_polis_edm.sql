-- 260 - SEQ_WORK_POLIS_EDM: penghitung nomor kasus Endorsment Fac In (irisan I3).
--
-- Keputusan work owner 11-10-2026: nomor pertama "EDM-00001" - sequence mulai 1. Bentuk teksnya (awalan `EDM-`, lima
-- digit berpadding nol) dirakit di Go, bukan di sini; sequence hanya ANGKA.
--
-- ⚠️ Kasus EDM Pega lama berbentuk `EDM-<n>` TANPA padding (`SetValueToEDMWork` langkah 7 `param.IDPrefix = "EDM-"`)
-- dan dipakai sebagai IDPEGA (JSON_POLIS PK IDPEGA, FACINPRODUCTION). Pengenal berpadding 5 digit tidak pernah sama
-- dengan pengenal Pega selama angkanya di bawah 10000 (`EDM-00001` ≠ `EDM-1`); mulai 10000 bentuknya sama. Bila nomor
-- awal diubah, DBA menyetel ulang lewat migrasi BARU (berkas ini tidak diubah).
-- Awalan lain di T_WORK_POLIS/tabel bersama: `NB-` (nbfacin), `EDMLF-` (endorsementlife), `EDMT-` (edmtreatyin).
--
-- NOCACHE: nomor tidak loncat saat instance dimulai ulang (pola 181 / 057). Ditulis, TIDAK dijalankan agent.
CREATE SEQUENCE {skema}.SEQ_WORK_POLIS_EDM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
