-- ⛔⛔ BERKAS INI MEMBUAT TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444`-`454`: rentang `treatyin` (`400-439`) PENUH.
--
-- ---------------------------------------------------------------------
-- Kurs MILIK KONTRAK — penghubung kontrak ↔ baris `TREATYEXCHANGEYEARLY`
-- ---------------------------------------------------------------------
--
-- Laporan pemakai 8 Oktober 2026 (kontrak HEALTH QUOTA SHARE 2026): grid
-- Rate of Exchange menampilkan TUJUH baris padahal yang diinput hanya DUA —
-- *"biarkan apa yg di input user yang tampil … jangan di tambah
-- tambahkan"*. Sebabnya: grid membaca SEMUA baris `TREATYEXCHANGEYEARLY`
-- tahun treaty-nya, dan tabel itu tidak tahu kontrak pemiliknya — setiap
-- kontrak 2026 melihat kurs kontrak 2026 lain.
--
-- ⭐ Kurs TETAP di `TREATYEXCHANGEYEARLY` (keputusan pemilik proses 4 & 6
-- Oktober 2026). Tabel ini hanya mencatat BARIS MANA milik KONTRAK MANA:
--
--   MASTERID  kontrak (`TREATY_IN.ID`)
--   IDKURS    `TREATYEXCHANGEYEARLY.ID`
--   URUTAN    urutan baris di grid kontrak itu
--
-- Save menulis ulang hubungan satu kontrak (Delete di grid melepas
-- hubungannya saja — baris kurs bersama TIDAK dihapus). Kontrak tanpa
-- hubungan (kontrak Pega lama) tetap membaca kurs tahun treaty-nya.
--
-- ⚠️ TANPA kunci asing — pola `448`-`454`.

CREATE TABLE {skema}.T_TREATY_KURS (
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  IDKURS    VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  CONSTRAINT PK_TT_KURS PRIMARY KEY (MASTERID, IDKURS)
)
/
CREATE INDEX {skema}.IX_TT_KURS_IDKURS ON {skema}.T_TREATY_KURS (IDKURS)
/
-- Isi awal — kontrak yang pemilik barisnya PASTI (Valid From = Commencement,
-- satu-satunya kontrak 2026 dengan periode itu):
--   1002305 TESTS                     08-10-2026  → 10124 IDR, 10125 USD
--   1002306 HEALTH QUOTA SHARE 2026   01-04-2026  → 10130 IDR, 10131 USD
-- Baris 01-07-2026 (10116, 10132, 10133) dimiliki salah satu dari TIGA
-- kontrak berperiode sama — TIDAK diisi di sini; pemiliknya ditentukan
-- pemilik proses.
--
-- ⚠️ `TREATYEXCHANGEYEARLY.ID` TIDAK unik: 10124, 10125, 10130, 10131
-- juga dipakai baris tahun 2025 (kurs JPY, KRW, PGK, PHP). Kuncinya
-- (ID, TREATYYEAR) — sama dengan `UPDATE` Save. Tanpa saringan tahun, isi
-- awal ini menyisipkan pasangan kembar dan jatuh di ORA-00001 PK_TT_KURS
-- (percobaan 9 Oktober 2026; tabel dan index-nya sudah berdiri, migrator
-- melewatinya saat diulang).
INSERT INTO {skema}.T_TREATY_KURS (MASTERID, IDKURS, URUTAN)
SELECT '1002305', ID, CASE CURRENCY WHEN 'IDR' THEN 1 ELSE 2 END FROM {skema}.TREATYEXCHANGEYEARLY WHERE TREATYYEAR = '2026' AND ID IN ('10124', '10125')
/
INSERT INTO {skema}.T_TREATY_KURS (MASTERID, IDKURS, URUTAN)
SELECT '1002306', ID, CASE CURRENCY WHEN 'IDR' THEN 1 ELSE 2 END FROM {skema}.TREATYEXCHANGEYEARLY WHERE TREATYYEAR = '2026' AND ID IN ('10130', '10131')
/
