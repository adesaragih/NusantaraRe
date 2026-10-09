-- Tiga anak cabang LAYER yang tidak pernah ikut ternormalisasi
--
-- Migrasi tiket 65, 67, 68 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal struktur: `4-erd-dan-tabel-datar/ERD-TREATY-IN-DAN-EDM.html` - ACUAN
-- struktur sistem lama sejak keputusan pemilik proses 2 Oktober 2026
-- ("fokus ini aja, yang lainnya ada yang salah itu"). Baris relasi 4, 9, 10.
-- Nama kolom: `2-to-spec/KAMUS-KOLOM.md` §10.24, §10.25, §10.26.
--
-- Tiga tabel dalam satu langkah, dan itu disengaja: `KELAS_BISNIS_KELOMPOK`
-- tidak dapat berdiri tanpa `KELOMPOK_LAYER`, dan memisahkannya menjadi dua
-- langkah membuat ada keadaan sah di mana induknya ada tanpa anaknya - keadaan
-- yang tidak seorang pun inginkan dan yang harus diurus setiap pemasangan.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS DITETAPKAN SADAR
-- ---------------------------------------------------------------------
--   Ketiganya **IKUT HAPUS** terhadap induk di dalam cabang, dan **TOLAK**
--   terhadap tabel acuan. Sumbernya DUA dokumen yang sepakat:
--
--     - `ERD-TREATY-IN-DAN-EDM.html` baris 4, 9, 10 - ketiganya `CASCADE`;
--     - `ERD.md` §2.5 (anak `DETAIL_PROPORSIONAL` ikut hapus) dan §2.7
--       (sepuluh rujukan tabel acuan, seluruhnya tolak).
--
--   ⚠️ Kolom `ON DELETE` di ERD HTML menyatakan dirinya **"USULAN rancangan
--   mengikuti konvensi contooh.xlsx, bukan perilaku sistem lama"**. Ia karena
--   itu TIDAK dipakai sendirian: yang mengikat tetap `ERD.md` §2, dan ERD HTML
--   dipakai untuk STRUKTUR - induk, nama kolom kunci asing, kardinalitas.
--   Keduanya sepakat di sini, dan itu disebut supaya pembaca berikutnya tahu
--   mana yang menang bila kelak tidak sepakat.
--
-- ---------------------------------------------------------------------
-- KENAPA TIGA TABEL, BUKAN SATU
-- ---------------------------------------------------------------------
--   `KELAS_BISNIS_LAYER` dan `KELAS_BISNIS_KELOMPOK` tampak kembar dan BUKAN.
--   Yang pertama menjawab "kelas apa saja yang ditanggung rincian proporsional
--   ini"; yang kedua "kelas apa saja yang ditanggung kelompok treaty X DI
--   DALAM layer ini". ERD menempatkannya pada kedalaman yang berbeda - baris 4
--   berinduk `LIMIT_DETAIL`, baris 10 berinduk `LIMIT_GROUP`. Menggabungkannya
--   melahirkan satu kolom induk yang separuh waktu kosong, dan `ADR-0041`
--   menolak bentuk itu.
--
-- ---------------------------------------------------------------------
-- NOL KOLOM NAMA - INV-59
-- ---------------------------------------------------------------------
--   Sistem lama menyalin `ClassOfBusiness` (nama) BERDAMPINGAN dengan
--   `ClassOfBusinessID`, dan `TreatyGroup` berdampingan dengan `TreatyGroupID`.
--   Salinan itu membuat nama yang diperbaiki di master tidak pernah merambat.
--   Di sini hanya kunci asingnya yang disimpan; namanya dibaca lewat join.
--
-- ---------------------------------------------------------------------
-- PERNYATAAN KEPUTUSAN - KUNCI ALAMI SENGAJA TIDAK DIPASANG
-- ---------------------------------------------------------------------
--   Ketiga tabel ini PUNYA kunci alami yang jelas:
--
--     KELAS_BISNIS_LAYER       (ID_DETAIL_PROPORSIONAL, ID_KELAS_BISNIS)
--     KELOMPOK_LAYER           (ID_LAYER, ID_KELOMPOK_TREATY)
--     KELAS_BISNIS_KELOMPOK    (ID_KELOMPOK_LAYER, ID_KELAS_BISNIS)
--
--   Ketiganya TIDAK dipasang sebagai UNIQUE, dan itu keputusan - bukan
--   kelalaian. `ddl-usulan/Z00_KUNCI_ALAMI.sql` menuliskan aturannya sendiri:
--   *"Setiap baris menyebut NOMOR INVARIANNYA. Constraint tanpa invarian tidak
--   ditulis di sini: ia tidak punya tempat untuk gagal."* `SPEC-INVARIAN.md`
--   berhenti di INV-71 dan belum menomori ketiganya.
--
--   ⚠️ Ongkosnya dibayar dan dinyatakan: baris kembar persis dapat masuk lewat
--   jalur mana pun. Yang menahannya untuk sementara lapisan services, dan
--   lapisan itu tidak menjaga pemuatan data maupun perbaikan manual.
--
--   SYARAT PEMBALIKAN: begitu `SPEC-INVARIAN.md` memberi nomornya, ketiganya
--   dipasang lewat migrasi korektif tersendiri - dan HARUS sebelum tiket 44
--   memuat data, sebab UNIQUE yang dipasang di atas data kembar akan gagal.
--   Tagihannya tercatat di `SPEC-INVARIAN.md` §13 dan di tiket 65, 67, 68.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
CREATE TABLE {skema}.KELAS_BISNIS_LAYER (
  ID_KELAS_BISNIS_LAYER   NUMBER(19) NOT NULL,
  ID_DETAIL_PROPORSIONAL  NUMBER(19) NOT NULL,
  ID_KELAS_BISNIS         NUMBER(19) NOT NULL,
  CONSTRAINT PK_KELAS_BISNIS_LAYER PRIMARY KEY (ID_KELAS_BISNIS_LAYER),
  CONSTRAINT FK_KELAS_BISNIS_LAYER_1 FOREIGN KEY (ID_DETAIL_PROPORSIONAL)
    REFERENCES {skema}.DETAIL_PROPORSIONAL (ID_DETAIL_PROPORSIONAL) ON DELETE CASCADE,
  CONSTRAINT FK_KELAS_BISNIS_LAYER_2 FOREIGN KEY (ID_KELAS_BISNIS)
    REFERENCES {skema}.KELAS_BISNIS (ID_KELAS_BISNIS)
)
/
CREATE TABLE {skema}.KELOMPOK_LAYER (
  ID_KELOMPOK_LAYER   NUMBER(19) NOT NULL,
  ID_LAYER            NUMBER(19) NOT NULL,
  ID_KELOMPOK_TREATY  NUMBER(19) NOT NULL,
  CONSTRAINT PK_KELOMPOK_LAYER PRIMARY KEY (ID_KELOMPOK_LAYER),
  CONSTRAINT FK_KELOMPOK_LAYER_1 FOREIGN KEY (ID_LAYER)
    REFERENCES {skema}.LAYER (ID_LAYER) ON DELETE CASCADE,
  CONSTRAINT FK_KELOMPOK_LAYER_2 FOREIGN KEY (ID_KELOMPOK_TREATY)
    REFERENCES {skema}.KELOMPOK_TREATY (ID_KELOMPOK_TREATY)
)
/
CREATE TABLE {skema}.KELAS_BISNIS_KELOMPOK (
  ID_KELAS_BISNIS_KELOMPOK  NUMBER(19) NOT NULL,
  ID_KELOMPOK_LAYER         NUMBER(19) NOT NULL,
  ID_KELAS_BISNIS           NUMBER(19) NOT NULL,
  CONSTRAINT PK_KELAS_BISNIS_KELOMPOK PRIMARY KEY (ID_KELAS_BISNIS_KELOMPOK),
  CONSTRAINT FK_KELAS_BISNIS_KELOMPOK_1 FOREIGN KEY (ID_KELOMPOK_LAYER)
    REFERENCES {skema}.KELOMPOK_LAYER (ID_KELOMPOK_LAYER) ON DELETE CASCADE,
  CONSTRAINT FK_KELAS_BISNIS_KELOMPOK_2 FOREIGN KEY (ID_KELAS_BISNIS)
    REFERENCES {skema}.KELAS_BISNIS (ID_KELAS_BISNIS)
)
/
-- ⚠️ ENAM index, bukan tiga. Tanpa UNIQUE di atas, kunci asing ke INDUK pun
-- tidak lagi terlayani index mana pun - biasanya ia terlayani karena kolomnya
-- MEMIMPIN kunci alaminya. Begitu kunci alaminya dipasang nanti, ketiga index
-- induk di bawah menjadi mubazir dan dibuang bersama migrasi korektifnya.
CREATE INDEX {skema}.IX_KELAS_BISNIS_LAYER_DP ON {skema}.KELAS_BISNIS_LAYER (ID_DETAIL_PROPORSIONAL)
/
CREATE INDEX {skema}.IX_KELOMPOK_LAYER_LAYER ON {skema}.KELOMPOK_LAYER (ID_LAYER)
/
CREATE INDEX {skema}.IX_KELAS_BISNIS_KLP_KL ON {skema}.KELAS_BISNIS_KELOMPOK (ID_KELOMPOK_LAYER)
/
CREATE INDEX {skema}.IX_KELAS_BISNIS_LAYER_KB ON {skema}.KELAS_BISNIS_LAYER (ID_KELAS_BISNIS)
/
CREATE INDEX {skema}.IX_KELOMPOK_LAYER_KLP ON {skema}.KELOMPOK_LAYER (ID_KELOMPOK_TREATY)
/
CREATE INDEX {skema}.IX_KELAS_BISNIS_KLP_KB ON {skema}.KELAS_BISNIS_KELOMPOK (ID_KELAS_BISNIS)
/
