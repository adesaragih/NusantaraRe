-- 947 - PRODUCT_TYPE_LIFE (tabel Pega M_PRODUCT_TYPE_LIFE, berganti nama di 946) mendapat lima kolom view lama
-- (langkah 2 dari 3; 948 isi + pemeriksaan + PK + buang JSONDATA).
--
-- Keputusan work owner 08-10-2026 K1: kolom akhir = kolom view lama, urutan sama (ID warisan VARCHAR2(6) di depan).
-- Lebar: nama-nama 200 (pola benefitlife 943; XML tanpa batas - pxTextInput b844 / pxAutoComplete b1107 / b1464)
-- [data DEV 08-10-2026: CoverName 34, Business 29, Benefit 48 byte]; ID master 10 = BUSINESS.ID / BENEFIT_LIFE.ID
-- VARCHAR2(10) (bukti plan_bukti_k1.sql kueri E). Tanpa FK ke BUSINESS / BENEFIT_LIFE (Pega tidak punya; K1).
-- Berdiri sendiri: ADD diulang mati di ORA-01430 dan penjaga inti menolak kolom baru lewat blok di jalur maju.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.PRODUCT_TYPE_LIFE ADD (
  COVERNAME     VARCHAR2(200),
  BUSINESS      VARCHAR2(200),
  BUSINESSID    VARCHAR2(10),
  BENEFIT       VARCHAR2(200),
  BENEFITID     VARCHAR2(10)
)
/
