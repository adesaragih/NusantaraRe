-- Dua celah bentuk pada migrasi 013 - tiket 00 Komite Claim Life.
--
-- ⛔ BUKAN pembuatan ulang. Butir km1: migrasi 030+ hanya untuk yang KURANG
-- dari 013; tabelnya sudah ada dan tidak disentuh isinya.
--
-- ⛔ CELAH 1 - kaskade yang dijanjikan dokumen tetapi tidak ada di DDL.
--
-- `STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md` menyebut relasi
-- `T_GENERAL_KOMITE` -> `T_KOMITE_KOMITELIST` lewat `DATA_KOMITE_ID` sebagai
-- `ON DELETE CASCADE` di TIGA tempat (baris 159, 238, 253), dan tiket 00
-- mengulanginya. `013_tabel_komite.sql` membuat `FK_KOMITELIST_KOMITE`
-- TANPA `ON DELETE`.
--
-- ⚠️ Akibatnya nyata, bukan kosmetik: menghapus satu kasus komite akan
-- DITOLAK Oracle (ORA-02292) selama masih ada baris roster yang
-- menunjuknya - dan penghapusan kasus komite adalah jalur yang tiket 05
-- perlukan saat jalur balik dijalankan.
--
-- ⛔ CELAH 2 - lebar kolom yang tidak sama dengan induknya.
--
--   T_WORK_CLAIM.ID              VARCHAR2(32)   (migrasi 001)
--   T_CLAIMLF_ADJUSTMENT.ID      VARCHAR2(32)   (migrasi 004)
--   T_GENERAL_KOMITE.ID          VARCHAR2(40)   <- shared PK, seharusnya 32
--   T_GENERAL_KOMITE.ADJUSTMENT_ID VARCHAR2(40) <- menunjuk kolom 32
--   T_KOMITE_KOMITELIST.ID       VARCHAR2(40)
--   T_KOMITE_KOMITELIST.DATA_KOMITE_ID VARCHAR2(40)
--
-- ⚠️ Oracle MENERIMA kunci tamu antarlebar berbeda, jadi ini tidak pernah
-- gagal - ia hanya berbohong. Kolom 40 karakter yang menunjuk kolom 32
-- karakter menjanjikan ruang yang tidak dapat dipakai: nilai ke-33 sampai
-- ke-40 tidak akan pernah punya induk. Shared PK yang lebarnya berbeda dari
-- induknya bukan shared PK, ia kebetulan yang sedang cocok.
--
-- ⚠️ `MODIFY` menyempitkan kolom. Ia AMAN hanya bila tidak ada nilai yang
-- lebih panjang; migrasi ini dijalankan sebelum ada data komite, dan bila
-- kelak ada, Oracle sendiri yang menolak (ORA-01441) - gagal terang.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).

ALTER TABLE {skema}.T_KOMITE_KOMITELIST DROP CONSTRAINT FK_KOMITELIST_KOMITE
/

ALTER TABLE {skema}.T_KOMITE_KOMITELIST MODIFY (
  ID             VARCHAR2(32),
  DATA_KOMITE_ID VARCHAR2(32)
)
/

ALTER TABLE {skema}.T_GENERAL_KOMITE MODIFY (
  ID            VARCHAR2(32),
  ADJUSTMENT_ID VARCHAR2(32)
)
/

ALTER TABLE {skema}.T_KOMITE_KOMITELIST ADD CONSTRAINT FK_KOMITELIST_KOMITE
  FOREIGN KEY (DATA_KOMITE_ID)
  REFERENCES {skema}.T_GENERAL_KOMITE (ID) ON DELETE CASCADE
/
