-- SEQ_WORK_CLAIM - pembangkit urutan identitas work object.
--
-- Pemilik: tiket 02 (register klaim). Dibuat terpisah dari 008 supaya berkas
-- yang sudah ada tidak disunting ulang.
--
-- [keputusan work owner 26-09-2026, butir aa] Identitas T_WORK_CLAIM adalah
-- TEKS BERFORMAT - CLM-xxxxxx untuk baris klaim, KMT-xxxxxx untuk baris komite
-- (tiket 14 AC 34, penyimpangan sadar dari ADR-U-0006). Yang datang dari
-- sequence hanyalah ANGKA urutannya; awalan dan LPAD(6) dirakit di Go, sebab
-- awalan bergantung jenis baris dan itu keputusan aturan dagang, bukan DDL.
--
-- ⛔ TANPA reset tahunan. Korpus tidak memuat satu pun bukti bahwa urutan work
-- object di-reset per tahun, dan mengarang reset akan membuat dua baris
-- bernomor sama pada tahun berbeda. Bila kelak DBA atau work owner menunjukkan
-- bahwa Pega me-reset-nya, itu keputusan baru beserta jalur migrasinya.
--
-- NOCACHE: nomor yang hilang saat instance mati jauh lebih mahal daripada
-- kecepatan, sama dengan seluruh sequence di 008.
CREATE SEQUENCE {skema}.SEQ_WORK_CLAIM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
