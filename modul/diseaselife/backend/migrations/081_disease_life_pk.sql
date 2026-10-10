-- 081 - Disease Life (diseaselife): PK_DISEASE_LIFE pada DISEASE_LIFE.ID (keputusan work owner 08-10-2026 D1.2).
--
-- DISEASE_LIFE (97.586 baris DEV) TANPA PK dan TANPA indeks; ID NULLABLE. PRIMARY KEY sekaligus menjadikan ID NOT NULL
-- (Oracle menambah NOT NULL pada kolom PK) dan membuat indeks unik bernama sama. Kolom dan nama tabel TIDAK berubah:
-- pencarian diagnosa claimlife dan prosedur Pega PEGA_DISEASE_LIFE tetap jalan.
--
-- BERHENTI dengan galat yang jelas bila masih ada ID NULL / kembar - TANPA menghapus data (D1.2):
--   - ID kembar  -> ORA-02437: cannot validate (<skema>.PK_DISEASE_LIFE) - primary key violated
--   - ID NULL    -> ORA-01449: column contains NULL values; cannot alter to NOT NULL
--   - PK lain sudah ada (nama lain) -> ORA-02260: table can have only one primary key
-- Galat itu dari ALTER itu sendiri: satu pernyataan DDL, berhasil seluruhnya atau tidak sama sekali (nol PK, nol
-- indeks, nol baris berubah), pelari berhenti, 081 tidak tercatat. Pemaksa UPDATE terpisah (bentuk 944 / 948) TIDAK
-- dipakai: `SET ID = NULL` tidak gagal pada kolom yang masih NULLABLE, `RPAD(ID, ...)` tidak menggigit ID NULL, dan
-- keduanya menulis ke tabel 97.586 baris - sedangkan galat ALTER sudah menyebut constraint dan sebabnya.
-- Bukti baca-saja SEBELUM -migrate (baris yang akan menghentikan langkah ini): docs/sql/dis_bukti.sql kueri C / D.
-- DEV: satu-satunya ID kembar = 102051 (baris uji TEST123 / Sakit) - DIHAPUS WO lewat docs/sql/disease_hapus_baris_uji.sql
-- SESUDAH cadangan P3 dan SEBELUM -migrate (D2); migrasi ini tidak pernah menghapus.
--
-- Blok berpelindung katalog (ALL_CONSTRAINTS, n = 0): aman DIULANG sesudah PK berdiri. Tanpa tanda kutip di teks
-- EXECUTE IMMEDIATE (polaPerintahKatalog). K0: migrasi MODUL (rentang 080-084). NOL COMMIT (ADR-U-0029). -migrate
-- dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'DISEASE_LIFE' AND CONSTRAINT_NAME = 'PK_DISEASE_LIFE';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.DISEASE_LIFE ADD CONSTRAINT PK_DISEASE_LIFE PRIMARY KEY (ID)';
  END IF;
END;
/
