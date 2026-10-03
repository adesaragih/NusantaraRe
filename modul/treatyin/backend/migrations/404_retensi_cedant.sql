-- RETENSI_CEDANT - retensi cedant per kelompok treaty per mata uang
--
-- Migrasi tiket 22 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/22-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/21_RETENSI_CEDANT.sql`.
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
-- INV-08: satu retensi per (versi, kelompok treaty, mata uang). Mata uang
-- IKUT di dalam kunci alaminya - dan itu bukan hiasan: UNIQUE tanpa mata uang
-- berarti "dilarang dua baris bermata uang berbeda", yang MENOLAK DATA SAH.
--
-- ⛔ KODE_MATA_UANG di sini VARCHAR2, bukan pengenal - ia tidak dapat
-- berkunci asing ke MATA_UANG. Lihat butir 6 `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
--
-- ⛔ INV-39 dan INV-40 BELUM DITULIS, dan sebabnya BUKAN penundaan.
--   Tingkat pencatatan NILAI_RETENSI sudah DIPUTUSKAN: `TREATY_100_PERSEN`,
--   `2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md` §2.2 - "keduanya membagi
--   risiko treaty, bukan mengambil bagian NuRe". §5 menyatakan INV-39 dan
--   INV-40 "kini dapat ditulis untuk 37 paket", dan paket ini salah satunya;
--   yang `ddl-usulan/` tinggalkan tanpa constraint hanya KEENAM paket yang
--   tingkatnya belum ditentukan, dan ini bukan salah satunya.
--   Jadi ketiadaannya di sini adalah LUBANG, bukan keputusan.
--   DITAGIH: tiket lapisan aplikasi, bersama INV-39/INV-40 paket lain.
CREATE TABLE {skema}.RETENSI_CEDANT (
  ID_RETENSI_CEDANT      NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK       NUMBER(19)          NOT NULL,
  ID_KELOMPOK_TREATY     NUMBER(19)          NOT NULL,
  NILAI_RETENSI          NUMBER(38,8)        NOT NULL,
  KODE_MATA_UANG         VARCHAR2(1000 CHAR) NOT NULL,
  CATATAN                VARCHAR2(1000 CHAR),
  PERSEN_BAGIAN_DIPAKAI  NUMBER(38,8),
  CONSTRAINT PK_RETENSI_CEDANT PRIMARY KEY (ID_RETENSI_CEDANT),
  CONSTRAINT UQ_RETENSI_CEDANT UNIQUE (ID_VERSI_KONTRAK, ID_KELOMPOK_TREATY, KODE_MATA_UANG),
  CONSTRAINT FK_RETENSI_CEDANT_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK),
  CONSTRAINT FK_RETENSI_CEDANT_2 FOREIGN KEY (ID_KELOMPOK_TREATY)
    REFERENCES {skema}.KELOMPOK_TREATY (ID_KELOMPOK_TREATY)
)
/
CREATE INDEX {skema}.IX_RETENSI_CEDANT_KLP ON {skema}.RETENSI_CEDANT (ID_KELOMPOK_TREATY)
/
