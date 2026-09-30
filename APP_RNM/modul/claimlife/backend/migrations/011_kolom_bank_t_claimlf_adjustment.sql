-- Tiga kolom bank sisa pada T_CLAIMLF_ADJUSTMENT - butir aj.
--
-- Pemilik: A1 (brief lanjutan 4 §1, `[DIPUTUSKAN 27-09-2026]`).
--
-- `[terverifikasi]` `Claim Life/Section/AdjustmentDetail_Section.xml` berkas
-- pecahan menampilkan ENAM medan bank pada satu baris adjustment:
--
--   5929  PayableTo
--   6122  NameOfBank        <- sudah ada sebagai NAME_OF_BANK
--   6245  NAMEOFBANK        <- ejaan kedua medan yang sama
--   6665  SwiftCode
--   6870  BranchOfBank
--   7049  NoAccount         <- sudah ada sebagai ACCOUNT_NO
--
-- Tiga yang sudah ada sejak tiket 14 (NAME_OF_BANK, ID_BANK, ACCOUNT_NO)
-- ditambah tiga di bawah menutup keenamnya.
--
-- ⛔ Ketiganya TEKS. Nomor dan kode bank tidak pernah menjadi bilangan:
-- SWIFT_CODE memuat huruf, dan kode cabang berawalan nol adalah hal biasa
-- (ADR-U-0022).
ALTER TABLE {skema}.T_CLAIMLF_ADJUSTMENT ADD (
  BRANCH_OF_BANK  VARCHAR2(140),
  SWIFT_CODE      VARCHAR2(35),
  PAYABLE_TO      VARCHAR2(140)
)
/
