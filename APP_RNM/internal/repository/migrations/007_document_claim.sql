-- DOCUMENT_CLAIM - dokumen pendukung per peserta
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
-- Tabel LINTAS-LINI. Induknya untuk Life adalah T_CLAIMLF_PREMIUMLIST_DETAIL.
--
-- Penghapusannya ditangani DI GO, bukan ON DELETE CASCADE basis data: karena
-- tabel ini lintas-lini, induknya berbeda tabel per lini sehingga satu aturan
-- cascade tidak dapat seragam. Karena itu kunci tamunya TANPA ON DELETE.
--
-- [terbuka] DAFTAR KOLOMNYA TIDAK DAPAT DITURUNKAN. [data DBA] Kelas Pega-nya
-- ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM; SQL-nya dibuat Pega sendiri dan nol
-- kemunculan di rule SQL mana pun. Menuliskan kolom isinya berarti mengarang.
-- Yang dibuat di sini hanya kunci utama dan kunci tamu, yang keduanya memang
-- sudah ditetapkan. Kolom isinya menunggu DBA.
CREATE TABLE {skema}.DOCUMENT_CLAIM (
  ID                      NUMBER(19) NOT NULL,
  PREMIUM_LIST_DETAIL_ID  VARCHAR2(32),
  CONSTRAINT PK_DOCUMENT_CLAIM PRIMARY KEY (ID),
  CONSTRAINT FK_DOC_PLD FOREIGN KEY (PREMIUM_LIST_DETAIL_ID)
    REFERENCES {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL (ID)
)
/
CREATE INDEX {skema}.IX_DOC_PLD_ID ON {skema}.DOCUMENT_CLAIM (PREMIUM_LIST_DETAIL_ID)
/
