-- Mengosongkan SELURUH tabel pendaratan tab Treaty In - jalur BALIK pemuat
--
-- ⛔ BUKAN berkas migrasi, dan sengaja TIDAK di `backend/migrations/`.
-- Berkas di sana dijalankan sekali dan tercatat di `T_MIGRASI`; yang ini
-- dijalankan sesering pemuatan diulang.
--
-- ---------------------------------------------------------------------
-- ⛔ DIPERBARUI 7 Oktober 2026 — DULU HANYA SEMBILAN TABEL
-- ---------------------------------------------------------------------
-- Bentuk sebelumnya ditulis ketika tabel pendaratan baru sembilan. Migrasi
-- 437, 438, 439, dan 446 menambah 21 lagi, dan berkas ini tidak ikut
-- tumbuh. Mengosongkan dengan bentuk lama menyisakan 21 tabel TERISI,
-- dan pemuatan berikutnya akan menumpuk di atasnya.
--
-- ⭐ DAFTARNYA DIAMBIL DARI `repository.PetaPendaratan` — sumber kebenaran
-- yang sama yang dipakai pemuat. `TestKosongkanMemuatSeluruhPeta` mengadu
-- berkas ini terhadap peta itu; menambah tabel ke peta tanpa menambahkannya
-- di sini membuat uji itu merah.
--
-- ---------------------------------------------------------------------
-- ⚠️ `DELETE`, bukan `DROP` dan bukan `TRUNCATE`
-- ---------------------------------------------------------------------
--   `DROP`     membongkar strukturnya - itu `430_..._down.sql`, dan ronde
--              ini melarang `DROP` terhadap POOLDATA.
--   `TRUNCATE` ber-DDL, jadi ia MENGIKAT transaksinya (commit implisit).
--              Pemuat yang gagal di tengah lalu `ROLLBACK` tidak akan
--              mendapatkan kembali baris yang `TRUNCATE` buang.
--   `DELETE`   dapat dibatalkan. Itu seluruh sebabnya.
--
-- ⛔ NOL `COMMIT` di berkas ini. Yang memanggilnya yang memutuskan kapan
-- perubahannya mengikat - dan di dalam uji, yang memanggilnya `ROLLBACK`.
--
-- ---------------------------------------------------------------------
-- ⛔ YANG TIDAK BOLEH IKUT DIKOSONGKAN
-- ---------------------------------------------------------------------
-- Tabel WARISAN hanya dibaca, dan sebagian DIPAKAI BERSAMA modul lain.
-- Mengosongkannya merusak modul yang tidak ada hubungannya dengan ronde ini:
--
--   TREATYEXCHANGEYEARLY   kurs - dipakai `aggregate`, `treatycontractout`,
--                          `nbfacin`, dan grid Rate of Exchange
--   REINSURANCETYPE        dropdown Treaty Type
--   ACHIEVEMENT            sumber `GetAchievement`
--   PROPORTIONALARRG       susunan treaty master (spreading tab Share)
--   M_TREATY_IN*           DILARANG KERAS dipakai aplikasi - jangan disentuh
--                          sama sekali, termasuk untuk dikosongkan
--
-- ---------------------------------------------------------------------
-- Urutan: ANAK sebelum INDUK — kebalikan urutan muat `PetaPendaratan`.
-- `FK_TT_INSTALLMENT_ITEM_1` memang `ON DELETE CASCADE`, jadi barisnya
-- secara teknis mubazir - ia ditulis apa adanya supaya berkas ini tetap
-- benar bila kaskade itu suatu hari dicabut, dan supaya cacah barisnya
-- dapat dibaca satu per satu.
--
-- Dijalankan untuk SATU kontrak: tambahkan `WHERE MASTERID = :1` pada
-- setiap pernyataan. Pemuat melakukan persis itu; lihat
-- `backend/repository/pendaratan_muat.go`.

-- [31..26] larik akar dan ringkasan
DELETE FROM {skema}.T_TREATY_SHARE_SUMMARY
/
DELETE FROM {skema}.T_TREATY_TOTAL
/
DELETE FROM {skema}.T_TREATY_LIMIT_SUMMARY
/
DELETE FROM {skema}.T_TREATY_LIMIT_MEASURE
/
DELETE FROM {skema}.T_TREATY_HAZARD_LIMIT
/
DELETE FROM {skema}.T_TREATY_REVISION
/

-- [25..20] cucu Limits dan Share
DELETE FROM {skema}.T_TREATY_FAC_SHARE_AMOUNT
/
DELETE FROM {skema}.T_TREATY_SHARE_AMOUNT
/
DELETE FROM {skema}.T_TREATY_LIMIT_AMOUNT
/
DELETE FROM {skema}.T_TREATY_LIMIT_GROUP_COB
/
DELETE FROM {skema}.T_TREATY_LIMIT_REINSTATEMENT
/
DELETE FROM {skema}.T_TREATY_LIMIT_ACH_PARAM
/
DELETE FROM {skema}.T_TREATY_LIMIT_DEDUCTION
/
DELETE FROM {skema}.T_TREATY_LIMIT_SPREADING
/
DELETE FROM {skema}.T_TREATY_LIMIT_ACHIEVEMENT
/
DELETE FROM {skema}.T_TREATY_LIMIT_COB
/

-- [19..15] anak Share dan Limits
DELETE FROM {skema}.T_TREATY_FAC_SHARE_DEDUCTION
/
DELETE FROM {skema}.T_TREATY_SHARE_DEDUCTION
/
DELETE FROM {skema}.T_TREATY_SHARE_SPREADING
/
DELETE FROM {skema}.T_TREATY_LIMIT_GROUP
/
DELETE FROM {skema}.T_TREATY_LIMIT_DETAIL
/

-- [14..10] induk Share dan Limits
DELETE FROM {skema}.T_TREATY_FAC_REINSURER
/
DELETE FROM {skema}.T_TREATY_FAC_SHARE
/
DELETE FROM {skema}.T_TREATY_RETRO_SHARE
/
DELETE FROM {skema}.T_TREATY_SHARE
/
DELETE FROM {skema}.T_TREATY_LIMITS
/

-- [9..1] tab berdiri sendiri
DELETE FROM {skema}.T_VIEW_COMMENT
/
DELETE FROM {skema}.M_TREATYIN_COINSCALE
/
DELETE FROM {skema}.T_TREATY_INSTALLMENT_ITEM
/
DELETE FROM {skema}.T_TREATY_INSTALLMENT
/
DELETE FROM {skema}.T_TREATY_RETENTION
/
DELETE FROM {skema}.T_TREATY_EGNPI
/
DELETE FROM {skema}.T_TREATY_ACCUMULATION
/
DELETE FROM {skema}.T_TREATY_PORTFOLIO
/
DELETE FROM {skema}.T_TREATY_REPORTING_PERIOD
/
