-- M_LOGIN_GO dan M_LOGIN_GO_WORKBASKET - login sungguhan (keputusan work owner
-- 01-10-2026).
--
-- Permintaan work owner: "nama table nya M_LOGIN_GO, didalam situ ada kolom
-- untuk menampung data dari M_UNIT, M_DIVISION, M_ORGANIZATION; tapi untuk
-- M_WORKBASKET 1 orang bisa beberapa workbasket". Lalu: "simpan aja code nya,
-- jangan ID" dan "M_LOGIN_GO_WORKBASKET aja, yang lain ga usah".
--
-- ⛔ SATU BARIS = SATU ORANG. `LOGIN_ID` adalah akun yang diketik di layar
-- login dan yang tercatat sebagai `CREATE_OP`; `NAME` nama tampilannya
-- (`CREATE_OP_NAME`).
--
-- ⛔ SANDI TIDAK PERNAH DISIMPAN. `PASSWORD_HASH` = hash bcrypt (cost 12).
--
-- ⛔ ORGANISASI DISIMPAN SEBAGAI CODE (`M_ORGANIZATION.CODE`, `M_DIVISION.CODE`,
-- `M_UNIT.CODE` - masing-masing UNIQUE, VARCHAR2(20)), TANPA FK: ketiga master
-- itu tidak dibuat migrasi aplikasi ini. Aplikasi memeriksa saat menyimpan
-- bahwa setiap CODE ada dan aktif, dan bahwa jenjangnya cocok (unit milik
-- divisinya, divisi milik organisasinya). Login TIDAK memeriksanya.
--
-- ⛔ WORKBASKET = PERAN. `M_WORKBASKET.WORKBASKET_ID` berisi nama peran yang
-- sudah dipakai aplikasi (`ReasLifeAdmin`, `ReasLifeSPV`, ...), jadi
-- `M_LOGIN_GO_WORKBASKET` sekaligus tabel peran - tidak ada tabel peran
-- terpisah. Tanpa FK ke `M_WORKBASKET` (alasan yang sama); hanya workbasket
-- yang `IS_ACTIVE = 1` di master yang berlaku.
--
-- ⛔ TANPA TABEL SESI. Cookie sesi ditandatangani HMAC dan memuat
-- `SESSION_VERSION`; menaikkan kolom itu (logout, ganti sandi, nonaktif)
-- mencabut SEMUA cookie lama akun itu.
--
-- Penguncian: `FAILED_COUNT` naik tiap sandi salah; sesudah 5, `LOCKED_UNTIL`
-- = 15 menit ke depan. `MUST_CHANGE_PASSWORD = '1'` untuk akun baru dan sandi
-- yang direset.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
CREATE TABLE {skema}.M_LOGIN_GO (
  LOGIN_ID             VARCHAR2(64) NOT NULL,
  NAME                 VARCHAR2(150) NOT NULL,
  PASSWORD_HASH        VARCHAR2(255) NOT NULL,
  ORGANIZATION_CODE    VARCHAR2(20),
  DIVISION_CODE        VARCHAR2(20),
  UNIT_CODE            VARCHAR2(20),
  IS_ACTIVE            VARCHAR2(1) DEFAULT '1' NOT NULL,
  FAILED_COUNT         NUMBER(5) DEFAULT 0 NOT NULL,
  LOCKED_UNTIL         DATE,
  MUST_CHANGE_PASSWORD VARCHAR2(1) DEFAULT '1' NOT NULL,
  SESSION_VERSION      NUMBER(10) DEFAULT 1 NOT NULL,
  LAST_LOGIN           DATE,
  TGL_CREATE           DATE DEFAULT SYSDATE NOT NULL,
  TGL_UPDATE           DATE,
  CONSTRAINT PK_M_LOGIN_GO PRIMARY KEY (LOGIN_ID),
  CONSTRAINT CK_M_LOGIN_GO_IS_ACTIVE CHECK (IS_ACTIVE IN ('0', '1')),
  CONSTRAINT CK_M_LOGIN_GO_MUST_CHANGE CHECK (MUST_CHANGE_PASSWORD IN ('0', '1'))
)
/
-- Satu orang, banyak workbasket. PK (LOGIN_ID, WORKBASKET_ID) sekaligus index
-- untuk "workbasket milik akun ini"; index kedua untuk "siapa di workbasket ini".
-- FK ke M_LOGIN_GO tanpa ON DELETE: akun dinonaktifkan, tidak dihapus.
CREATE TABLE {skema}.M_LOGIN_GO_WORKBASKET (
  LOGIN_ID      VARCHAR2(64) NOT NULL,
  WORKBASKET_ID VARCHAR2(64) NOT NULL,
  TGL_CREATE    DATE DEFAULT SYSDATE NOT NULL,
  CONSTRAINT PK_M_LOGIN_GO_WORKBASKET PRIMARY KEY (LOGIN_ID, WORKBASKET_ID),
  CONSTRAINT FK_M_LOGIN_GO_WB_LOGIN FOREIGN KEY (LOGIN_ID) REFERENCES {skema}.M_LOGIN_GO (LOGIN_ID)
)
/
CREATE INDEX {skema}.IX_M_LOGIN_GO_WB_WORKBASKET ON {skema}.M_LOGIN_GO_WORKBASKET (WORKBASKET_ID)
/
