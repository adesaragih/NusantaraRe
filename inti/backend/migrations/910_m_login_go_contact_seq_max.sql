-- 910 - M_LOGIN_GO_CONTACT_SEQ mengikuti nomor CON tertinggi (perintah work owner 04-10-2026: "tambahan untuk seq
-- con kelola user, script-nya ambil dari max con+1").
--
-- 905 membuatnya dengan angka tetap START WITH 1001. Langkah ini membuangnya lalu membuatnya lagi dengan nilai awal
-- = nomor CON tertinggi di M_LOGIN_GO.CONTACT_ID + 1, DIHITUNG di basis data tempat migrasi berjalan (bentuk
-- `migrasi.BacaSequenceDariKueri`), jadi DEV dan PROD masing-masing mulai sesudah CON terbesarnya sendiri.
-- Tidak pernah di bawah 1001: jaminan 905 - nomor CON aplikasi di atas nomor kontak SFAGIS (CON-108 di DEV) - tetap.
-- M_LOGIN_GO satu-satunya sumber CON baru (MARKETINGOFFICER.CLIENTID baris baru menyalinnya).
--
-- ⚠️ Bila langkah ini gagal SESUDAH DROP (jarang: blok hanya membaca M_LOGIN_GO), pengulangan berhenti di
-- ORA-02289 - jalankan blok kedua secara manual, lalu catat `910_m_login_go_contact_seq_max` di T_MIGRASI.
-- Selama sequence belum ada, membuat akun di Kelola User gagal.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
DROP SEQUENCE {skema}.M_LOGIN_GO_CONTACT_SEQ
/
DECLARE
  n    NUMBER;
  awal NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES
   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = 'M_LOGIN_GO_CONTACT_SEQ';
  IF n = 0 THEN
    SELECT GREATEST(NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(CONTACT_ID, '^CON-([0-9]+)$', 1, 1, NULL, 1))), 0), 1000) + 1
      INTO awal FROM {skema}.M_LOGIN_GO;
    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}.M_LOGIN_GO_CONTACT_SEQ START WITH ' || awal || ' INCREMENT BY 1 NOCACHE';
  END IF;
END;
/
