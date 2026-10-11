-- Jalur mundur untuk 430_tabel_tab_treatyin.sql
--
-- ⚠️ Membongkar ini MEMBUANG hasil pemuatan - 26.536 baris bila kedelapan
-- tabel sudah terisi penuh. Itu tidak kehilangan data: setiap barisnya
-- salinan dari `M_TREATY_IN.JSONDATA`, yang tidak disentuh dan tetap menjadi
-- sumbernya sampai Langkah 3 selesai. Sesudah Langkah 3, kalimat ini berhenti
-- benar dan berkas ini harus diberi peringatan yang lebih keras.
--
-- ⛔ Untuk MENGOSONGKAN tabelnya tanpa membongkar strukturnya - yang
-- diperlukan saat memuat ulang - jangan pakai berkas ini. Pakai
-- `modul/treatyin/alat/kosongkan-tab-treatyin.sql`, yang `DELETE`, bukan
-- `DROP`, dan karenanya dapat dibatalkan.
--
-- Urutan terbalik dari urutan pembuatan: anak sebelum induknya, supaya
-- kunci asing `FK_MTI_INSTALLMENTITEM_1` tidak menahan `DROP`.
DROP TABLE {skema}.M_TREATYIN_COMMENT
/
DROP TABLE {skema}.M_TREATYIN_INSTALLMENTITEM
/
DROP TABLE {skema}.M_TREATYIN_INSTALLMENT
/
DROP TABLE {skema}.M_TREATYIN_RETENTION
/
DROP TABLE {skema}.M_TREATYIN_EGNPI
/
DROP TABLE {skema}.M_TREATYIN_ACCUMULATION
/
DROP TABLE {skema}.M_TREATYIN_PORTFOLIO
/
DROP TABLE {skema}.M_TREATYIN_REPORTINGPERIOD
/
