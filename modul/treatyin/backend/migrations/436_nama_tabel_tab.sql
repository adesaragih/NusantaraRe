-- Tabel tab DIGANTI NAMA mengikuti `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`,
-- `PXOBJCLASS` DICABUT, dan isinya DIKOSONGKAN.
--
-- Keputusan pemilik proses 5 Oktober 2026, tiga hal sekaligus.
--
-- ---------------------------------------------------------------------
-- ⚠️ DUA HAL YANG XLSX ITU MINTA DAN TIDAK DAPAT DIPENUHI
-- ---------------------------------------------------------------------
--  1. Tiap `T_TREATY_*` di sana membawa `FK TREATY_IN_ID -> TREATY_IN.ID
--     CASCADE`. ⛔ MUSTAHIL: `TREATY_IN` TIDAK punya kunci utama maupun
--     kunci unik atas `ID` -- terukur, dan tercatat sebagai §2 di
--     `KOREKSI-ERD-VERSUS-POOLDATA.md`. Kunci asing ke sana ditolak Oracle
--     dengan `ORA-02270`. Kolom penautnya karena itu tetap `MASTERID`
--     bertipe teks, tanpa kunci asing, persis seperti sebelum berkas ini.
--
--  2. `M_TREATYIN_COINSCALE` NOL PADANAN di xlsx itu -- disapu seluruh enam
--     lembarnya, nol kata `coins`, `scale`, maupun `skala`. Ia karena itu
--     TIDAK diganti nama. Mengarang `T_TREATY_COINS_SCALE` berarti menaruh
--     nama yang tidak pernah seorang pun putuskan.
--
-- ⚠️ Berkas itu sendiri menyatakan di baris pertamanya "POTRET SISTEM LAMA,
--    24 September 2026. BUKAN RANCANGAN", dan menunjuk `ERD-SKEMA-BARU.xlsx`
--    (25 September) sebagai gambar skema barunya -- yang justru memakai nama
--    Indonesia yang modul ini sudah pakai. Pemilik proses diberi tahu dan
--    tetap memilih nama `T_TREATY_*`. Dicatat di sini supaya yang membaca
--    nanti tahu pilihannya SADAR, bukan terlewat.
--
-- Pembalikan: `436_nama_tabel_tab_down.sql` mengembalikan nama lama beserta
-- kolom `PXOBJCLASS`. ⛔ Ia TIDAK mengembalikan isinya -- 27.238 baris yang
-- berkas ini hapus dimuat ulang dari `M_TREATY_IN.JSONDATA` lewat pemuat
-- yang sudah ada, bukan dari pembalikan ini.

-- --------------------------------------------------------------------
-- 1 · ISI DIKOSONGKAN. Anak lebih dulu: `INSTALLMENTITEM` menunjuk
--     `INSTALLMENT` dengan `ON DELETE CASCADE`, dan urutan terbalik
--     membuat baris anak terhapus diam-diam oleh kaskade alih-alih oleh
--     pernyataan yang tertulis di sini.
-- --------------------------------------------------------------------
DELETE FROM {skema}.M_TREATYIN_INSTALLMENTITEM
/
DELETE FROM {skema}.M_TREATYIN_INSTALLMENT
/
DELETE FROM {skema}.M_TREATYIN_REPORTINGPERIOD
/
DELETE FROM {skema}.M_TREATYIN_PORTFOLIO
/
DELETE FROM {skema}.M_TREATYIN_ACCUMULATION
/
DELETE FROM {skema}.M_TREATYIN_EGNPI
/
DELETE FROM {skema}.M_TREATYIN_RETENTION
/
DELETE FROM {skema}.M_TREATYIN_COMMENT
/
DELETE FROM {skema}.M_TREATYIN_COINSCALE
/

-- --------------------------------------------------------------------
-- 2 · `PXOBJCLASS` DICABUT dari kesembilannya.
--
--     Ia nama kelas Pega yang ikut terbawa saat pendaratan. Nol pembaca di
--     seluruh modul ini, dan menyimpan nama kelas sistem lama di dalam
--     tabel model baru adalah mengundang orang menurunkan perilaku darinya.
-- --------------------------------------------------------------------
ALTER TABLE {skema}.M_TREATYIN_REPORTINGPERIOD DROP COLUMN PXOBJCLASS
/
ALTER TABLE {skema}.M_TREATYIN_PORTFOLIO DROP COLUMN PXOBJCLASS
/
ALTER TABLE {skema}.M_TREATYIN_ACCUMULATION DROP COLUMN PXOBJCLASS
/
ALTER TABLE {skema}.M_TREATYIN_EGNPI DROP COLUMN PXOBJCLASS
/
ALTER TABLE {skema}.M_TREATYIN_RETENTION DROP COLUMN PXOBJCLASS
/
ALTER TABLE {skema}.M_TREATYIN_INSTALLMENT DROP COLUMN PXOBJCLASS
/
ALTER TABLE {skema}.M_TREATYIN_INSTALLMENTITEM DROP COLUMN PXOBJCLASS
/
ALTER TABLE {skema}.M_TREATYIN_COMMENT DROP COLUMN PXOBJCLASS
/
ALTER TABLE {skema}.M_TREATYIN_COINSCALE DROP COLUMN PXOBJCLASS
/

-- --------------------------------------------------------------------
-- 3 · NAMA TABEL. Delapan diganti; `COINSCALE` tetap (lihat kepala).
-- --------------------------------------------------------------------
ALTER TABLE {skema}.M_TREATYIN_REPORTINGPERIOD RENAME TO T_TREATY_REPORTING_PERIOD
/
ALTER TABLE {skema}.M_TREATYIN_PORTFOLIO RENAME TO T_TREATY_PORTFOLIO
/
ALTER TABLE {skema}.M_TREATYIN_ACCUMULATION RENAME TO T_TREATY_ACCUMULATION
/
ALTER TABLE {skema}.M_TREATYIN_EGNPI RENAME TO T_TREATY_EGNPI
/
ALTER TABLE {skema}.M_TREATYIN_RETENTION RENAME TO T_TREATY_RETENTION
/
ALTER TABLE {skema}.M_TREATYIN_INSTALLMENT RENAME TO T_TREATY_INSTALLMENT
/
ALTER TABLE {skema}.M_TREATYIN_INSTALLMENTITEM RENAME TO T_TREATY_INSTALLMENT_ITEM
/
ALTER TABLE {skema}.M_TREATYIN_COMMENT RENAME TO T_VIEW_COMMENT
/

-- --------------------------------------------------------------------
-- 4 · NAMA BATASAN, INDEX, DAN SEQUENCE ikut.
--
--     ⛔ Batasan bernama `PK_MTI_*` di atas tabel bernama `T_TREATY_*`
--     adalah dua kebenaran tentang satu benda. Yang membaca katalog akan
--     mengira tabelnya belum diganti.
-- --------------------------------------------------------------------
ALTER TABLE {skema}.T_TREATY_REPORTING_PERIOD RENAME CONSTRAINT PK_MTI_REPORTINGPERIOD TO PK_TT_REPORTING_PERIOD
/
ALTER TABLE {skema}.T_TREATY_REPORTING_PERIOD RENAME CONSTRAINT UQ_MTI_REPORTINGPERIOD TO UQ_TT_REPORTING_PERIOD
/
ALTER TABLE {skema}.T_TREATY_PORTFOLIO RENAME CONSTRAINT PK_MTI_PORTFOLIO TO PK_TT_PORTFOLIO
/
ALTER TABLE {skema}.T_TREATY_PORTFOLIO RENAME CONSTRAINT UQ_MTI_PORTFOLIO TO UQ_TT_PORTFOLIO
/
ALTER TABLE {skema}.T_TREATY_ACCUMULATION RENAME CONSTRAINT PK_MTI_ACCUMULATION TO PK_TT_ACCUMULATION
/
ALTER TABLE {skema}.T_TREATY_ACCUMULATION RENAME CONSTRAINT UQ_MTI_ACCUMULATION TO UQ_TT_ACCUMULATION
/
ALTER TABLE {skema}.T_TREATY_EGNPI RENAME CONSTRAINT PK_MTI_EGNPI TO PK_TT_EGNPI
/
ALTER TABLE {skema}.T_TREATY_EGNPI RENAME CONSTRAINT UQ_MTI_EGNPI TO UQ_TT_EGNPI
/
ALTER TABLE {skema}.T_TREATY_RETENTION RENAME CONSTRAINT PK_MTI_RETENTION TO PK_TT_RETENTION
/
ALTER TABLE {skema}.T_TREATY_RETENTION RENAME CONSTRAINT UQ_MTI_RETENTION TO UQ_TT_RETENTION
/
ALTER TABLE {skema}.T_TREATY_INSTALLMENT RENAME CONSTRAINT PK_MTI_INSTALLMENT TO PK_TT_INSTALLMENT
/
ALTER TABLE {skema}.T_TREATY_INSTALLMENT RENAME CONSTRAINT UQ_MTI_INSTALLMENT TO UQ_TT_INSTALLMENT
/
ALTER TABLE {skema}.T_TREATY_INSTALLMENT_ITEM RENAME CONSTRAINT PK_MTI_INSTALLMENTITEM TO PK_TT_INSTALLMENT_ITEM
/
ALTER TABLE {skema}.T_TREATY_INSTALLMENT_ITEM RENAME CONSTRAINT UQ_MTI_INSTALLMENTITEM TO UQ_TT_INSTALLMENT_ITEM
/
ALTER TABLE {skema}.T_TREATY_INSTALLMENT_ITEM RENAME CONSTRAINT FK_MTI_INSTALLMENTITEM_1 TO FK_TT_INSTALLMENT_ITEM_1
/
ALTER TABLE {skema}.T_VIEW_COMMENT RENAME CONSTRAINT PK_MTI_COMMENT TO PK_TV_COMMENT
/
ALTER INDEX {skema}.IX_MTI_INSTALLMENTITEM_MST RENAME TO IX_TT_INSTALLMENT_ITEM_MST
/
RENAME SEQ_MTI_REPORTINGPERIOD TO SEQ_TT_REPORTING_PERIOD
/
RENAME SEQ_MTI_PORTFOLIO TO SEQ_TT_PORTFOLIO
/
RENAME SEQ_MTI_ACCUMULATION TO SEQ_TT_ACCUMULATION
/
RENAME SEQ_MTI_EGNPI TO SEQ_TT_EGNPI
/
RENAME SEQ_MTI_RETENTION TO SEQ_TT_RETENTION
/
RENAME SEQ_MTI_INSTALLMENT TO SEQ_TT_INSTALLMENT
/
RENAME SEQ_MTI_INSTALLMENTITEM TO SEQ_TT_INSTALLMENT_ITEM
/
RENAME SEQ_MTI_COMMENT TO SEQ_TV_COMMENT
/
