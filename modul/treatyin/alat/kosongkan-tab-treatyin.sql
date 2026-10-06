-- Mengosongkan kesembilan tabel tab Treaty In - jalur BALIK pemuat
--
-- ⛔ BUKAN berkas migrasi, dan sengaja TIDAK di `backend/migrations/`.
-- Berkas di sana dijalankan sekali dan tercatat di `T_MIGRASI`; yang ini
-- dijalankan sesering pemuatan diulang.
--
-- ⚠️ `DELETE`, bukan `DROP` dan bukan `TRUNCATE`.
--   `DROP`    membongkar strukturnya - itu `430_..._down.sql`, dan ronde ini
--             melarang `DROP` terhadap POOLDATA.
--   `TRUNCATE` ber-DDL, jadi ia MENGIKAT transaksinya (commit implisit).
--             Pemuat yang gagal di tengah lalu `ROLLBACK` tidak akan
--             mendapatkan kembali baris yang `TRUNCATE` buang.
--   `DELETE`  dapat dibatalkan. Itu seluruh sebabnya.
--
-- ⛔ NOL `COMMIT` di berkas ini. Yang memanggilnya yang memutuskan kapan
-- perubahannya mengikat - dan di dalam uji, yang memanggilnya `ROLLBACK`.
--
-- Anak sebelum induknya. `FK_TT_INSTALLMENT_ITEM_1` memang
-- `ON DELETE CASCADE`, jadi baris pertama di bawah secara teknis mubazir -
-- ia ditulis apa adanya supaya berkas ini tetap benar bila kaskade itu suatu
-- hari dicabut, dan supaya cacah barisnya dapat dibaca satu per satu.
--
-- Dijalankan untuk SATU kontrak: tambahkan `WHERE MASTERID = :1` pada
-- kesembilan pernyataan. Pemuat melakukan persis itu; lihat
-- `backend/repository/pendaratan_muat.go`.
DELETE FROM {skema}.T_TREATY_INSTALLMENT_ITEM
/
DELETE FROM {skema}.T_TREATY_INSTALLMENT
/
DELETE FROM {skema}.T_VIEW_COMMENT
/
DELETE FROM {skema}.M_TREATYIN_COINSCALE
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
