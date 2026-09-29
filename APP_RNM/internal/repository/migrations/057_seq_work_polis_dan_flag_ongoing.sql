-- SEQ_WORK_POLIS dan T_WORK_POLIS.FLAG_ONGOING_POLICY - butir bn (GILIRAN-13).
--
-- `[DIPUTUSKAN; veto work owner]` bn, di atas pl3 (brief modul PremiumList
-- baris 43: pengenal work dirakit dari SEQ_WORK_POLIS, pola
-- PengenalWorkBerikut). Tanpa keduanya tombol portal `Input Offer` /
-- `Input Premium` tidak dapat membuat kasus.
--
-- ⛔ SEQ_WORK_POLIS hanya memberi ANGKA. Awalan `NBLF-` dan bentuknya (tanpa
-- nol di depan) dirakit di Go - `repository/polis_kasus.go` - sebab bentuk
-- itu dibaca dari data warisan, bukan dari DDL.
--
-- ⚠️ START WITH 1 AMAN HANYA SELAMA T_WORK_POLIS BELUM BERISI BARIS WARISAN.
-- Pengenal warisan `NBLF-<n>` sudah mencapai lima digit (sampel DEV
-- 29-09-2026). Sebelum data warisan dimigrasikan ke T_WORK_POLIS, urutan ini
-- wajib dimajukan melewati angka warisan terbesar - OQ-PL-15.
--
-- ⛔ FLAG_ONGOING_POLICY menyimpan nilai VERBATIM `FlagPolicy` tombol portal:
-- "0" (`Input Offer`, Section/PremiumList.xml b3310/b3597) dan "1"
-- (`Input Premium`, b3958/b4233), yang CreateInputLife b618 tuliskan ke
-- `curWorkPage.FlagOnGoingPolicy`. TEKS satu karakter, bukan bilangan
-- (ADR-U-0022): decision table `IsFlagOnGoingPolicy` membandingkannya sebagai
-- teks (b293/b294).
--
-- ⛔ Nullable: baris yang lahir sebelum migrasi ini tidak punya benderanya,
-- dan mengarang nilainya lebih buruk daripada mengosongkannya (ADR-U-0027).
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
CREATE SEQUENCE {skema}.SEQ_WORK_POLIS START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
ALTER TABLE {skema}.T_WORK_POLIS ADD (
  FLAG_ONGOING_POLICY VARCHAR2(1)
)
/
