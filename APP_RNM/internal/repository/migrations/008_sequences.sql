-- Sequence identitas
--
-- Migrasi tiket 14 Claim Life. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Aturan yang dijaga seluruh berkas di folder ini:
--   ADR-U-0027  seluruh kolom nullable kecuali kunci utama; wajib-isi di services
--   ADR-U-0003  ADR-U-0016  uang NUMBER(38,8), tidak pernah float
--   ADR-U-0029  nol COMMIT di teks SQL; transaksi dibuka-ditutup aplikasi
--   ADR-U-0006  identitas dari sequence, kecuali yang dinyatakan berformat
--
-- Sequence adalah pembangkit angka berurut milik Oracle. ADR-U-0006
-- menetapkan identitas berasal dari sequence, bukan dari cap waktu maupun teks.
--
-- ADR-U-0006 berlaku untuk T_CLAIMLF_* dan DOCUMENT_CLAIM. Ia TIDAK berlaku
-- untuk T_WORK_CLAIM dan kedua tabel ber-shared-PK: identitas ketiganya adalah
-- nomor bisnis berformat, dan itu penyimpangan sadar yang dicatat di tiket 14.
--
-- NOCYCLE berarti sequence tidak pernah berputar kembali ke awal, sehingga
-- satu nomor tidak pernah terpakai dua kali.
CREATE SEQUENCE {skema}.SEQ_CLAIMLF_PLD START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_CLAIMLF_ADJ START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_CLAIMLF_SPR START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_CLAIMLF_SPR_RETRO START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_T_CLAIMLF_DOCUMENT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
