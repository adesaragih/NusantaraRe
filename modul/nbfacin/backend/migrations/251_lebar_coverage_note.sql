-- 251 - lebar T_COVERAGELIST.COVERAGE_NOTE = lebar sumbernya (bug DEV 11-10-2026: unggah CSV Marine Hull case NB-151034 gagal
-- simpan ORA-12899 T_COVERAGELIST.COVERAGE_NOTE actual 1276 / maximum 500). Rentang KETIGA nbfacin 250-254. Ditulis, TIDAK
-- dijalankan agent; work owner / DBA yang menjalankan. Nol COMMIT (ADR-U-0029). SATU pernyataan per blok `/` (ORA-02158 /
-- ORA-00922 bila digabung).
--
-- T_COVERAGELIST.COVERAGE_NOTE VARCHAR2(500) (migrasi 193) adalah lebar RANCANGAN (Tabel-Flat xlsx), bukan batas Pega - di Pega
-- coverage tinggal di work object. Sumbernya POOLDATA.COVERAGE.NAME VARCHAR2(4000 BYTE) (`D:\migrasi\RNM\DDL\COVERAGE.txt` L4;
-- BrowseCoverageFacIn_RD NamaCoverage -> .CoverageNote: tab Coverage FIRE / Aneka / cargo dan CSV) - pola butir 80 (lebar
-- sumber). 193 TIDAK diubah (sudah jalan di DEV). Hanya MELEBARKAN; data yang ada tidak berubah.
-- ⛔ T_ADDITIONALCOVERAGE.COVERAGE_NOTE (Warranty Extension cargo, 500) SENGAJA tidak di sini: tabelnya lahir di 256, dan pelari
-- migrasi mengurutkan menurut NAMA berkas - 251 berjalan sebelum 256, sehingga di basis data baru ALTER-nya gagal ORA-00942.
ALTER TABLE {skema}.T_COVERAGELIST MODIFY (COVERAGE_NOTE VARCHAR2(4000))
/
