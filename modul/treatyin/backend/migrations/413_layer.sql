-- LAYER beserta limit, deductible, MDP, dan pemulihan limitnya
--
-- Migrasi tiket 31 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/31-*.md`, `2-to-spec/KAMUS-KOLOM.md`,
-- `ddl-usulan/{15_LAYER,26_NILAI_MDP,27_NILAI_MDP_MINIMUM,28_PEMULIHAN_LIMIT}.sql`.
--
-- Empat tabel dalam satu langkah, dan itu disengaja: ketiga tabel anak tidak
-- dapat berdiri tanpa LAYER, dan tiket 31 menyebut keempatnya sebagai SATU
-- artefak ("LAYER, PEMULIHAN_LIMIT, dan kedua tabel anak paket uang").
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS DITETAPKAN SADAR. ⚠️ RALAT 2 Oktober 2026.
-- ---------------------------------------------------------------------
--   Berkas ini pernah menyebut INV-18 sambil MEMBIARKAN bawaan Oracle. Itu
--   membaca INV-18 terbalik: ia menuntut perilaku hapus "DITETAPKAN SADAR,
--   TIDAK DIBIARKAN BAWAAN", dan bawaan yang kebetulan cocok bukan keputusan.
--   Sumber keputusannya `4-erd-dan-tabel-datar/ERD.md` §2 - dokumen MENGIKAT
--   yang tidak pernah dibuka sampai hari ini.
--
--   Keempat kunci asing berkas ini **IKUT HAPUS**, dan keempatnya disebut
--   ERD.md: LAYER di §2.4, PEMULIHAN_LIMIT di §2.3b, NILAI_MDP dan
--   NILAI_MDP_MINIMUM di §2.3c - "sebuah nilai per mata uang tidak punya arti
--   tanpa induknya". Berkas ini sebelumnya menulis TANPA ON DELETE pada
--   keempatnya; itu keliru di keempatnya.
--
-- INV-05: satu layer per (versi, nomor layer, bagian layer).
--
-- ---------------------------------------------------------------------
-- PAKET UANG GOLONGAN C - kolom mata uang dipasang, BOLEH KOSONG
-- ---------------------------------------------------------------------
--   DEDUCTIBLE dan LIMIT_AGREGAT memperoleh kolom mata uang BERNAMA di
--   sebelahnya (MATA_UANG_DEDUCTIBLE, MATA_UANG_LIMIT_AGREGAT), bukan
--   "KODE_MATA_UANG" polos - baris ini memuat lebih dari satu paket uang, dan
--   satu nama untuk dua arti adalah cacat yang ditanam dengan tangan sendiri.
--   Sistem lama TIDAK merekam mata uang untuk besaran ini sama sekali, jadi
--   kolomnya akan KOSONG pada seluruh baris hasil migrasi. Kosong di sana
--   berarti "tidak pernah dicatat", bukan "belum diisi".
--   DILIHAT DI: KTV-C. DITAGIH: T-6 - ia MENYEMPITKAN, tidak memblokir.
--
-- ---------------------------------------------------------------------
-- LIMIT_AGREGAT - CONSTRAINT YANG SENGAJA TIDAK DIPASANG
-- ---------------------------------------------------------------------
--   APA        : tidak ada CHECK maupun NOT NULL yang memasangkan paket uang
--                ini dengan PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40).
--   KENAPA     : tingkat pencatatannya BELUM DITENTUKAN - ia salah satu dari
--                KEENAM paket yang `TINGKAT-PENCATATAN-14-PAKET-UANG.md` §5
--                tinggalkan tanpa constraint. Constraint yang dipasang
--                sekarang menuntut nilai yang migrasi tidak tahu cara
--                mengisinya.
--   AKIBAT     : kolom ini menerima angka tanpa menyatakan angkanya untuk
--                SELURUH TREATY atau untuk BAGIAN NuRe.
--   DITAGIH    : ketika jawaban T-3 kembali dari teknik treaty.
--
-- ---------------------------------------------------------------------
-- ✅ UQ_LAYER MENEGAKKAN INV-05 — TERVERIFIKASI DI ORACLE 02-10-2026
-- ---------------------------------------------------------------------
--   ⚠️ RALAT. Baris ini pernah berbunyi "UQ_LAYER TIDAK menegakkan INV-05
--   sepenuhnya", dengan alasan BAGIAN_LAYER boleh kosong dan Oracle
--   memperlakukan NULL sebagai tidak sama dengan NULL. ITU KELIRU, dan
--   kekeliruannya hanya ketahuan karena diuji terhadap Oracle sungguhan.
--
--   Aturan Oracle yang sebenarnya: sebuah entri DILEWATI indeks unik hanya
--   bila SELURUH kolom kuncinya NULL. Di sini ID_VERSI_KONTRAK dan
--   NOMOR_LAYER selalu terisi, sehingga barisnya DIINDEKS - dan dua baris
--   (versi, 1, NULL) ditolak.
--
--   Bukti: menyisipkan baris kedua (9000001, 1, NULL) menghasilkan
--     ORA-00001: unique constraint (UQ_LAYER) violated
--   Uji kendali, layer berbagian (versi, 2, 1) ganda, juga ditolak UQ_LAYER.
--   Keduanya dijalankan 02-10-2026, lalu ROLLBACK.
--
--   Pelajaran yang layak disimpan: penalaran tentang semantik NULL Oracle
--   GAGAL dua kali di berkas ini - sekali oleh penulisnya, sekali oleh
--   peninjaunya - dan yang memperbaikinya bukan argumen, melainkan satu
--   INSERT. Jangan menulis pernyataan keputusan tentang perilaku basis data
--   yang belum dijalankan.
--
-- ⛔ INV-49 BELUM TERTULIS. Daftar periksa tiket 31 menuntut "INV-49 tertulis
-- sebagai pengecualian bernama terhadap INV-47 - bukan dibiarkan terbaca
-- sebagai pelanggaran". INV-47 milik batch 2 (tiket 47, "paling banyak satu
-- versi tak-terminal per kontrak") dan belum berdiri, sehingga pengecualian
-- terhadapnya belum punya tempat berdiri. DITAGIH: bersama tiket 47.
--
-- ⚠️ PERSEN_ROL - "rate on line", tarif premi sebagai persen dari limit.
-- Dibawa apa adanya. `F-15` dapat MENCABUTNYA dan menambahkan tabel anak premi
-- diperoleh; bentuk LAYER sendiri tidak berubah karenanya.
CREATE TABLE {skema}.LAYER (
  ID_LAYER                      NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK              NUMBER(19)          NOT NULL,
  NOMOR_LAYER                   NUMBER(10)          NOT NULL,
  BAGIAN_LAYER                  NUMBER(10),
  JENIS_LAYER                   VARCHAR2(40 CHAR)   NOT NULL,
  JENIS_BAGIAN_LAYER            VARCHAR2(40 CHAR),
  CAKUPAN                       VARCHAR2(1000 CHAR),
  LIMIT                         NUMBER(38,8)        NOT NULL,
  KODE_MATA_UANG                VARCHAR2(1000 CHAR) NOT NULL,
  DEDUCTIBLE                    NUMBER(38,8)        NOT NULL,
  MATA_UANG_DEDUCTIBLE          VARCHAR2(1000 CHAR),
  PERSEN_PENYESUAIAN            NUMBER(38,8),
  PERSEN_MINIMUM_DEPOSIT        NUMBER(38,8),
  PORSI_PEMULIHAN_LIMIT         NUMBER(38,8),
  TARIF_PREMI_PEMULIHAN         NUMBER(38,8),
  LIMIT_AGREGAT                 NUMBER(38,8),
  MATA_UANG_LIMIT_AGREGAT       VARCHAR2(1000 CHAR),
  MDP                           NUMBER(38,8),
  MDP_MINIMUM                   NUMBER(38,8),
  PERSEN_MDP_MINIMUM            NUMBER(38,8),
  MDP_DIGABUNG                  VARCHAR2(40 CHAR)   NOT NULL,
  TANPA_HITUNG_PREMI_PEMULIHAN  VARCHAR2(40 CHAR)   NOT NULL,
  DEDUCTIBLE_KEDUA              NUMBER(38,8),
  PERSEN_ROL                    NUMBER(38,8),
  CONSTRAINT PK_LAYER PRIMARY KEY (ID_LAYER),
  CONSTRAINT UQ_LAYER UNIQUE (ID_VERSI_KONTRAK, NOMOR_LAYER, BAGIAN_LAYER),
  CONSTRAINT FK_LAYER_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/
CREATE TABLE {skema}.NILAI_MDP (
  ID_NILAI_MDP    NUMBER(19)          NOT NULL,
  ID_LAYER        NUMBER(19)          NOT NULL,
  KODE_MATA_UANG  VARCHAR2(1000 CHAR) NOT NULL,
  NILAI           NUMBER(38,8)        NOT NULL,
  CONSTRAINT PK_NILAI_MDP PRIMARY KEY (ID_NILAI_MDP),
  CONSTRAINT FK_NILAI_MDP_1 FOREIGN KEY (ID_LAYER)
    REFERENCES {skema}.LAYER (ID_LAYER)
)
/
CREATE TABLE {skema}.NILAI_MDP_MINIMUM (
  ID_NILAI_MDP_MINIMUM  NUMBER(19)          NOT NULL,
  ID_LAYER              NUMBER(19)          NOT NULL,
  KODE_MATA_UANG        VARCHAR2(1000 CHAR) NOT NULL,
  NILAI                 NUMBER(38,8)        NOT NULL,
  CONSTRAINT PK_NILAI_MDP_MINIMUM PRIMARY KEY (ID_NILAI_MDP_MINIMUM),
  CONSTRAINT FK_NILAI_MDP_MINIMUM_1 FOREIGN KEY (ID_LAYER)
    REFERENCES {skema}.LAYER (ID_LAYER)
)
/
CREATE TABLE {skema}.PEMULIHAN_LIMIT (
  ID_PEMULIHAN_LIMIT    NUMBER(19)   NOT NULL,
  ID_LAYER              NUMBER(19)   NOT NULL,
  NOMOR_URUT_PEMULIHAN  NUMBER(10)   NOT NULL,
  PERSEN_PEMULIHAN      NUMBER(38,8) NOT NULL,
  PERSEN_TAMBAHAN       NUMBER(38,8),
  CATATAN               VARCHAR2(1000 CHAR),
  CONSTRAINT PK_PEMULIHAN_LIMIT PRIMARY KEY (ID_PEMULIHAN_LIMIT),
  CONSTRAINT FK_PEMULIHAN_LIMIT_1 FOREIGN KEY (ID_LAYER)
    REFERENCES {skema}.LAYER (ID_LAYER)
)
/
CREATE INDEX {skema}.IX_NILAI_MDP_LAYER ON {skema}.NILAI_MDP (ID_LAYER)
/
CREATE INDEX {skema}.IX_NILAI_MDP_MIN_LAYER ON {skema}.NILAI_MDP_MINIMUM (ID_LAYER)
/
CREATE INDEX {skema}.IX_PEMULIHAN_LIMIT_LAYER ON {skema}.PEMULIHAN_LIMIT (ID_LAYER)
/
