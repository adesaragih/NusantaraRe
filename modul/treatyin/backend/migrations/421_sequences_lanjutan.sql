-- Sequence identitas untuk entitas rekonsiliasi ERD
--
-- Migrasi tiket 65, 67, 68, 69, 70, 74 Treaty In. Satu berkas = satu langkah
-- migrasi. Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis
-- miring.
--
-- Berdiri TERPISAH dari tabelnya, mengikuti 402: `ddl-usulan/` memuat nol
-- `CREATE SEQUENCE`, dan menambalnya ke dalam berkas tabel membuat ketiadaan
-- itu berhenti terbaca sebagai temuan.
--
-- INV-02: pengenal datang dari SEQUENCE, tidak pernah dari cap waktu maupun
-- teks. INV-03: `NOCYCLE` - satu nomor tidak pernah terpakai dua kali.
--
-- Nama sequence dibatasi 30 bita. Yang terpanjang di bawah,
-- SEQ_TRIN_KELAS_BISNIS_LAYER, 27 bita. `KELAS_BISNIS_KELOMPOK` disingkat
-- menjadi `KELAS_BISNIS_KLP` sebab bentuk penuhnya tepat 30 bita - tepat di
-- batas, dan tepat di batas adalah tempat yang buruk untuk berdiri.
-- Singkatan `KLP` sudah dipakai `IX_RETENSI_CEDANT_KLP` di migrasi 404.
CREATE SEQUENCE {skema}.SEQ_TRIN_KELAS_BISNIS_LAYER START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_KELOMPOK_LAYER START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_KELAS_BISNIS_KLP START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_PENCAPAIAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_RINCIAN_ANGSURAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_ARSIP_MUATAN_KELUAR START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
