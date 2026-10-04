-- Sequence identitas tiket 33, 34, dan 37
--
-- Migrasi tiket 33, 34, dan 37 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- INV-02: pengenal dari SEQUENCE. INV-03: NOCYCLE.
--
-- ⛔ `ddl-usulan/` tidak memuat satu pun CREATE SEQUENCE (butir 4
-- `docs/KEPUTUSAN-PENYELARASAN-REPO.md`).
--
-- ⚠️ `SEQ_TRIN_NILAI_PB_MINIMUM` DISINGKAT. Bentuk penuhnya,
-- `SEQ_TRIN_NILAI_PREMI_BRUTO_MINIMUM`, 34 bita - melewati batas 30 (§16).
-- Singkatan dipilih di awalan `NILAI_PREMI_BRUTO` -> `NILAI_PB`, bukan di
-- `MINIMUM`, sebab "MINIMUM" yang hilang membuat kedua sequence tidak dapat
-- dibedakan. Dijaga `TestNamaObjekDiBawahTigaPuluhBita`.
CREATE SEQUENCE {skema}.SEQ_TRIN_BAGIAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_NILAI_PREMI_BRUTO START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_NILAI_PB_MINIMUM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_DETAIL_PROPORSIONAL START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_NILAI_CADANGAN_PREMI START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_POTONGAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
