-- M_LOGIN_GO: email, nomor HP, NIK, dan jabatan akun (Kelola User, permintaan work owner 03-10-2026:
-- "tambahkan email, no hp, nik dan jabatan; buat dalam bahasa inggris").
--
-- Nama kolom berbahasa Inggris. SEMUANYA opsional: akun yang sudah ada mendapat NULL, dan isian kosong ditulis
-- NULL. Format diperiksa aplikasi saat menyimpan (`login.PeriksaKontak`): EMAIL berbentuk alamat email,
-- PHONE_NUMBER 8-15 digit (boleh diawali +, boleh spasi dan tanda hubung), EMPLOYEE_ID = NIK (Nomor Induk
-- Karyawan, bukan NIK kependudukan) huruf/angka/titik/garis miring/tanda hubung, JOB_POSITION teks jabatan.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
ALTER TABLE {skema}.M_LOGIN_GO ADD (
  EMAIL        VARCHAR2(254),
  PHONE_NUMBER VARCHAR2(30),
  EMPLOYEE_ID  VARCHAR2(30),
  JOB_POSITION VARCHAR2(150)
)
/
