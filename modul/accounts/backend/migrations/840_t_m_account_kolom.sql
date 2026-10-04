-- 840 - T_M_ACCOUNT: kolom Owner, Create Date, dan Description layar Accounts (keputusan work owner 04-10-2026:
-- "tambahkan create date dan create op", "Owner = CREATEOP", "Description tambah"). Tabel warisan POOLDATA
-- (Pega SFAGIS); baris lama (18.658 di DEV) kosong pada ketiga kolom - nol default, nol isi ulang.
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
ALTER TABLE {skema}.T_M_ACCOUNT ADD (
  CREATEDATE  DATE,
  CREATEOP    VARCHAR2(100),
  DESCRIPTION VARCHAR2(4000)
)
/
