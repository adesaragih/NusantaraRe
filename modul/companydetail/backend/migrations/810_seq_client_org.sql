-- 810 - sequence nomor ORG Company Detail (perintah work owner 04-10-2026: "buat seq aja, start-nya dari id
-- max+1" dan "max dari org yang ada di client").
--
-- Create organisasi mengambil nomor dari SEQ_CLIENT_ORG.NEXTVAL, menggantikan "nomor tertinggi + 1 di bawah
-- LOCK TABLE CLIENT". Nilai awal = nomor ORG tertinggi yang sudah terpakai + 1, DIHITUNG di basis data tempat
-- migrasi ini berjalan - Oracle hanya menerima angka pada START WITH, jadi angka DEV tidak dibekukan ke PROD:
--   - CLIENT.IDVIEW `ORG-n` (organisasi di tabel datar), dan
--   - M_CLIENT.ID `ASM-SFAGIS-WORK-ORG ORG-n` (dokumen Pega yang belum pindah ke CLIENT - Copy Old memindahkannya
--     dengan nomor aslinya, jadi nomor itu pun tidak boleh dipakai ulang).
-- Polanya = models.AwalanID dan models.AwalanIDView (dijaga TestMigrasi810SequenceMulaiSesudahNomorTertinggi).
--
-- Bentuk blok `migrasi.BacaSequenceDariKueri`: tanya ALL_SEQUENCES dulu, supaya pengulangan sesudah gagal di
-- tengah tidak mati di ORA-00955. NOCACHE: nomor tidak melompat saat instance dimulai ulang. Nomor yang sudah
-- diambil transaksi yang gagal tidak kembali (sifat sequence) - nomor ORG boleh berlubang, tidak pernah kembar.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n    NUMBER;
  awal NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES
   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = 'SEQ_CLIENT_ORG';
  IF n = 0 THEN
    SELECT NVL(MAX(N), 0) + 1 INTO awal FROM (
      SELECT TO_NUMBER(REGEXP_SUBSTR(IDVIEW, '[0-9]+$')) N FROM {skema}.CLIENT
       WHERE REGEXP_LIKE(IDVIEW, '^ORG-[0-9]+$')
      UNION ALL
      SELECT TO_NUMBER(REGEXP_SUBSTR(ID, '[0-9]+$')) N FROM {skema}.M_CLIENT
       WHERE REGEXP_LIKE(ID, '^ASM-SFAGIS-WORK-ORG ORG-[0-9]+$'));
    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}.SEQ_CLIENT_ORG START WITH ' || awal || ' INCREMENT BY 1 NOCACHE NOCYCLE';
  END IF;
END;
/
