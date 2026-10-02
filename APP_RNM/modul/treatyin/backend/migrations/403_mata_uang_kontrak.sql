-- MATA_UANG_KONTRAK - beberapa mata uang berlaku, masing-masing berkurs dan berperiode
--
-- Migrasi tiket 20 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/20-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/16_MATA_UANG_KONTRAK.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS DITETAPKAN SADAR. ⚠️ RALAT 2 Oktober 2026.
-- ---------------------------------------------------------------------
--   Berkas ini pernah berbunyi: "INV-18 - TANPA ON DELETE: bawaan Oracle
--   MENOLAK, dan menolak yang dikehendaki." ITU MEMBACA INV-18 TERBALIK.
--
--   INV-18 berbunyi "perilaku hapus setiap kunci asing DITETAPKAN SADAR,
--   TIDAK DIBIARKAN BAWAAN". Membiarkan bawaan lalu menyebut nama invarian
--   yang melarangnya bukan pemenuhan - bahkan ketika bawaannya kebetulan
--   cocok. Yang dituntut keputusan yang DINYATAKAN, bukan yang kebetulan benar.
--
--   Keputusan sadarnya sudah ada dan MENGIKAT: `4-erd-dan-tabel-datar/ERD.md`
--   §2, yang menyatakan `ikut hapus`, `tolak`, atau `putus` per relasi.
--   Dokumen itu tidak pernah dibuka sampai 2 Oktober 2026 - tiga puluh tiga
--   kunci asing ditulis tanpanya.
--
--   Perilaku hapus tiap kunci asing di berkas ini kini disebut di barisnya
--   sendiri. Daftar lengkap seluruh modul: bab "Kaskade"
--   di `MODUL.md`, yang juga menyalakan `TestKaskadeHanyaPadaRelasiTerdaftar`.
--
-- INV-07: satu mata uang hanya sekali per versi kontrak.
-- INV-44: KODE_MATA_UANG di sini BERTIPE NUMBER(19) - ia pengenal baris
-- MATA_UANG, bukan teks - sehingga kunci asingnya dapat dipasang. Kolom
-- bernama sama di RETENSI_CEDANT, EGNPI, TERMIN, dan BATAS_PER_BAHAYA
-- BERTIPE VARCHAR2(1000) di `ddl-usulan/`, dan di sana kunci asing yang sama
-- TIDAK dapat dipasang. Ketidakseragaman itu ada di spec, bukan dibuat di
-- sini; dicatat di `docs/KEPUTUSAN-PENYELARASAN-REPO.md` butir 6.
-- ---------------------------------------------------------------------
-- INV-38 dan INV-57 SENGAJA BUKAN CONSTRAINT - pernyataan keputusan.
-- ---------------------------------------------------------------------
--   APA     : tidak ada CHECK yang menolak KURS nol, tidak ada yang menolak
--             kurs yang dipaksakan menjadi satu (INV-38), dan tidak ada yang
--             menolak tanggal kurs lebih akhir daripada transaksi yang
--             memakainya (INV-57).
--   KENAPA  : INV-38 MEMANG dapat ditulis sebagai CHECK - ia membandingkan
--             kolom pada BARIS YANG SAMA, tidak seperti INV-53/55/56. Yang
--             menahannya hanya ADR-0056 (K-4), dan keputusan itu berlaku
--             seragam di seluruh modul ini: nol aturan bisnis di basis data.
--             Memasangnya di sini saja membuat satu invarian ditegakkan di
--             tempat yang berbeda dari sembilan lainnya.
--             INV-57 tidak dapat menjadi CHECK sama sekali - "transaksi yang
--             memakainya" ada di tabel lain.
--   AKIBAT  : kurs nol dan kurs bertanggal mundur diterima sampai jalur
--             simpan berdiri. Kurs nol MEMBAGI DENGAN NOL di setiap
--             perhitungan hilir - ini bukan cacat kosmetik.
--   DITAGIH : tiket lapisan aplikasi. Daftar periksa tiket 20 menuntut DUA
--             uji untuk INV-38 - "menolak kurs nol DAN menolak kurs yang
--             dipaksakan menjadi satu" - bukan satu.
--   ⛔ Bila pemilik proses memutuskan CHECK boleh untuk invarian SEBARIS,
--      INV-38 adalah calon pertamanya, bersama INV-29 (berkas 401).

CREATE TABLE {skema}.MATA_UANG_KONTRAK (
  ID_MATA_UANG_KONTRAK   NUMBER(19)   NOT NULL,
  ID_VERSI_KONTRAK       NUMBER(19)   NOT NULL,
  KODE_MATA_UANG         NUMBER(19)   NOT NULL,
  KURS                   NUMBER(38,8) NOT NULL,
  TANGGAL_MULAI_BERLAKU  DATE         NOT NULL,
  TANGGAL_AKHIR_BERLAKU  DATE         NOT NULL,
  CONSTRAINT PK_MATA_UANG_KONTRAK PRIMARY KEY (ID_MATA_UANG_KONTRAK),
  CONSTRAINT UQ_MATA_UANG_KONTRAK UNIQUE (ID_VERSI_KONTRAK, KODE_MATA_UANG),
  CONSTRAINT FK_MATA_UANG_KONTRAK_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK),
  CONSTRAINT FK_MATA_UANG_KONTRAK_2 FOREIGN KEY (KODE_MATA_UANG)
    REFERENCES {skema}.MATA_UANG (ID_MATA_UANG)
)
/
CREATE INDEX {skema}.IX_MATA_UANG_KONTRAK_MU ON {skema}.MATA_UANG_KONTRAK (KODE_MATA_UANG)
/
