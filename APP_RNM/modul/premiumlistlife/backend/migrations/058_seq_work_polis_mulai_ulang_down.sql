-- Jalur mundur 058: SEQ_WORK_POLIS kembali ke bentuk 057 (START WITH 1).
--
-- ⚠️ Nomor yang sudah terbit di atas 22374 TIDAK ikut mundur - kasus yang
-- lahir sesudah 058 tetap bernomor itu, dan sequence baru mulai lagi dari 1.
-- Jalur mundur ini untuk skema uji; di skema berisi data nyata ia membuka
-- tabrakan nomor yang justru ditutup 058.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
DROP SEQUENCE {skema}.SEQ_WORK_POLIS
/
CREATE SEQUENCE {skema}.SEQ_WORK_POLIS START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
