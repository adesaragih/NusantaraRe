-- 263 - konversi kasus Endorsment Fac In ke bentuk POSITION / POSITION_NOTE (keputusan work owner W-EDM-18, mengganti
-- W-EDM-11; pola NB tiket 77): POSITION = tahap 'Policy', POSITION_NOTE = workbasket 'ReasFacInMarketing'.
--
-- ⚠️ URUTAN: kolom T_WORK_POLIS.POSITION_NOTE VARCHAR2(255) dibuat migrasi nbfacin 254. Pelari mengurutkan menurut NAMA
-- berkas lintas modul, jadi 254 berjalan lebih dulu di -migrate yang sama. Bila 254 belum ada, pernyataan ini gagal
-- ORA-00904 dan -migrate berhenti di sini - jalankan sesudah 254.
-- Biner sendiri TOLERAN (repository/posisinote.go): tanpa kolom ia menulis bentuk lama (POSITION = 'ReasFacInMarketing');
-- baris semacam itu yang dikonversi di sini. Hanya baris LINI 'FACEDM' berbentuk lama; idempoten.
-- Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
UPDATE {skema}.T_WORK_POLIS SET POSITION_NOTE = POSITION, POSITION = 'Policy', TGL_UPDATE = SYSDATE
WHERE LINI = 'FACEDM' AND POSITION = 'ReasFacInMarketing'
/
