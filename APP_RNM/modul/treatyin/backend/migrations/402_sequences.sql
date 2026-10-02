-- Sequence identitas Treaty In
--
-- Migrasi tiket 14 dan 15 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- INV-02: pengenal datang dari SEQUENCE, TIDAK PERNAH dari cap waktu maupun
-- dari teks. Sistem lama membentuk `M_TREATY_IN.ID` dari pola bernomor revisi
-- (TDA-11) - teks yang membawa arti, dan karena itu teks yang dapat bentrok.
--
-- INV-03: pengenal tidak dipakai ulang. NOCYCLE berarti sequence tidak pernah
-- berputar kembali ke awal, sehingga satu nomor tidak pernah terpakai dua kali.
--
-- ⛔ `ddl-usulan/` TIDAK memuat satu pun CREATE SEQUENCE - keenam tabel acuan
-- dan kedua tabel kontrak di sana berdiri tanpa pembangkit pengenalnya,
-- sementara INV-02 dan daftar periksa tiket 14 menuntutnya ("pengenal keduanya
-- datang dari SEQUENCE tanpa CYCLE; tidak ada jalur lain yang dapat memberi
-- pengenal"). Berkas ini menutup lubang itu, dan menutupnya di sini - bukan
-- menambalnya diam-diam ke dalam berkas tabel - supaya ketiadaannya di
-- `ddl-usulan/` tetap terbaca sebagai temuan.
--
-- Nama sequence dibatasi 30 bita (§16). Yang terpanjang di bawah,
-- SEQ_TRIN_JENIS_REASURANSI, 25 bita.
CREATE SEQUENCE {skema}.SEQ_TRIN_KONTRAK START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_VERSI_KONTRAK START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_MATA_UANG START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_JENIS_POTONGAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_JENIS_REASURANSI START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_BAHAYA START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_KELOMPOK_TREATY START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_KELAS_BISNIS START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
