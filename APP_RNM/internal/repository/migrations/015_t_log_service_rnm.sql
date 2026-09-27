-- Outbox efek keluar lintas modul - butir aq.
-- [keputusan work owner 27-09-2026, butir ax] Nama tabel T_LOG_SERVICE_RNM
-- (semula T_EFEK_KELUAR): ia LOG sekaligus ANTREAN panggilan layanan luar. Nama
-- index, constraint, dan sequence mengikuti. Migrasi ini belum pernah dijalankan
-- di Oracle mana pun saat diganti nama, sehingga disunting di tempat.
--
-- Pemilik: A2 (brief lanjutan 5 §2, `[DIPUTUSKAN]` turunan tiket 12 AC
-- "antre-ulang, kegagalan tercatat" + ADR-U-0008; bentuknya milik executor).
--
-- ⛔ Ditulis DALAM TRANSAKSI YANG SAMA dengan aksi bisnisnya. Itulah inti pola
-- outbox: bila klaim tersimpan, efek keluarnya PASTI terantre; bila
-- transaksinya batal, antreannya ikut batal. Dua transaksi berarti ada saat
-- ketika salah satunya ada tanpa yang lain.
--
-- ⚠️ LINTAS MODUL - `LINI` dan `MODUL` memisahkannya. Komite tiket 06-07
-- memakai tabel yang sama; menulis tabel kembar per modul berarti dua worker,
-- dua kebijakan backoff, dan dua tempat kegagalan bersembunyi.
--
-- ⛔ `STATUS` bertipe TEKS dengan himpunan tertutup, bukan angka: `antre`,
-- `jalan`, `selesai`, `gagal-permanen`. Kode bukan bilangan (ADR-U-0022), dan
-- himpunan tertutup membuat nilai asing terlihat.
--
-- ⚠️ `MUATAN` CLOB berisi JSON. Ia TIDAK boleh memuat nama orang, kredensial,
-- maupun alamat - hanya pengenal dan angka. Penjaga statik menegakkannya.
CREATE TABLE {skema}.T_LOG_SERVICE_RNM (
  ID              VARCHAR2(40)  NOT NULL,
  LINI            VARCHAR2(16)  NOT NULL,
  MODUL           VARCHAR2(32)  NOT NULL,
  JENIS_EFEK      VARCHAR2(32)  NOT NULL,
  RUJUKAN         VARCHAR2(40)  NOT NULL,
  MUATAN          CLOB,
  STATUS          VARCHAR2(16)  NOT NULL,
  PERCOBAAN       NUMBER(5)     DEFAULT 0 NOT NULL,
  JADWAL_BERIKUT  TIMESTAMP,
  GALAT_TERAKHIR  VARCHAR2(4000),
  DIBUAT          TIMESTAMP     NOT NULL,
  DIPERBARUI      TIMESTAMP,
  CONSTRAINT PK_LOG_SERVICE_RNM PRIMARY KEY (ID)
)
/
-- Worker memilih lewat (STATUS, JADWAL_BERIKUT); tanpa index ia memindai
-- seluruh outbox setiap putaran.
CREATE INDEX {skema}.IX_LOG_SERVICE_RNM_JADWAL ON {skema}.T_LOG_SERVICE_RNM (STATUS, JADWAL_BERIKUT)
/
CREATE INDEX {skema}.IX_LOG_SERVICE_RNM_RUJUKAN ON {skema}.T_LOG_SERVICE_RNM (RUJUKAN)
/
CREATE SEQUENCE {skema}.SEQ_LOG_SERVICE_RNM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
