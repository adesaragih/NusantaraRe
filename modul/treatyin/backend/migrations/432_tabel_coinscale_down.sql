-- Jalur mundur untuk 432_tabel_coinscale.sql
--
-- ⚠️ Membongkar ini membuang 702 baris hasil pemuatan. Itu tidak kehilangan
-- data: setiap barisnya salinan dari `M_TREATY_IN.JSONDATA`, yang tidak
-- disentuh dan tetap menjadi sumbernya.
--
-- ⛔ Untuk MENGOSONGKAN tanpa membongkar strukturnya, pakai
-- `modul/treatyin/alat/kosongkan-tab-treatyin.sql` atau pemuat berbendera
-- `-kosongkan` — keduanya `DELETE`, dan karena itu dapat dibatalkan.
DROP TABLE {skema}.M_TREATYIN_COINSCALE
/
