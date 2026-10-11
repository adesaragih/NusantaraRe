-- 206 - tab utama Loss Record kasus FIRE (tiket 55): daftar kerugian TINGKAT CASE `.OfferFacIn.ListCauseOfLoss`, Premi
-- PAR / MB / PL / Other, dan loss ratio 1 / 3 / 5 tahun (gross dan net).
--
-- TABEL BARU (keputusan work owner 06-10-2026 butir 4 "setuju tabel baru"): workbook rancangan (`loader/skema_gen.go`) hanya
-- punya ListCauseOfLoss per LOKASI (T_LISTCAUSEOFLOSS, tiket 42) dan loss ratio per lokasi (T_LOCATIONLIST) - jalur case
-- tidak ada, jadi bukan amandemen rancangan. Nama mengikuti properti Pega (`Activity\SaveLossRecord.xml`):
--   T_LOSSRECORD     satu per kasus (UNIQUE PARENT_ID = T_GENERAL_POLIS.ID VARCHAR2(32), butir 76.1) - PremiPAR / PremiMB /
--                    PremiPL / PremiOther, LossRatio{1,3,5}YearPercent, NetLossRatio{1,3,5}Year;
--   T_LOSSRECORDLIST baris kerugian urut SEQ_NO (DOL terbaru dulu) - UWYEAR, TahunKe, ClassOfBusiness, DOL (teks Pega
--                    DateTime, pola T_LISTCAUSEOFLOSS.DATE_OF_LOSS), CauseOfLoss, Currency, LOSS_AMOUNT, LOSS_AMOUNT_IDR,
--                    KursIDR.
-- Uang / persen NUMBER(38,8) (ADR-0016); kurs NUMBER(38,10) (TOIDR TREATYEXCHANGEYEARLY = teks VARCHAR2(255)). Lebar teks =
-- sumbernya: CLASS_OF_BUSINESS = BUSINESS.NOTE 4000, CAUSE_OF_LOSS 500 (pola T_LISTCAUSEOFLOSS), CURRENCY 10. ID NUMBER(19) dari
-- sequence (A92). Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_LOSSRECORD (
  ID                       NUMBER(19) NOT NULL,
  PARENT_ID                VARCHAR2(32) NOT NULL,
  PREMI_PAR                NUMBER(38,8),
  PREMI_MB                 NUMBER(38,8),
  PREMI_PL                 NUMBER(38,8),
  PREMI_OTHER              NUMBER(38,8),
  LOSS_RATIO1_YEAR_PERCENT NUMBER(38,8),
  LOSS_RATIO3_YEAR_PERCENT NUMBER(38,8),
  LOSS_RATIO5_YEAR_PERCENT NUMBER(38,8),
  NET_LOSS_RATIO1_YEAR     NUMBER(38,8),
  NET_LOSS_RATIO3_YEAR     NUMBER(38,8),
  NET_LOSS_RATIO5_YEAR     NUMBER(38,8),
  CONSTRAINT PK_T_LOSSRECORD PRIMARY KEY (ID),
  CONSTRAINT UQ_T_LOSSRECORD_PARENT UNIQUE (PARENT_ID),
  CONSTRAINT FK_LOSSRECORD_GENERAL FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID)
)
/
CREATE SEQUENCE {skema}.SEQ_T_LOSSRECORD START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.T_LOSSRECORDLIST (
  ID                NUMBER(19) NOT NULL,
  PARENT_ID         NUMBER(19) NOT NULL,
  SEQ_NO            NUMBER(5) NOT NULL,
  ROW_UID           VARCHAR2(36) NOT NULL,
  UW_YEAR           VARCHAR2(10),
  TAHUN_KE          NUMBER(5),
  CLASS_OF_BUSINESS VARCHAR2(4000),
  DATE_OF_LOSS      VARCHAR2(30),
  CAUSE_OF_LOSS     VARCHAR2(500),
  CURRENCY          VARCHAR2(10),
  LOSS_AMOUNT       NUMBER(38,8),
  LOSS_AMOUNT_IDR   NUMBER(38,8),
  KURS_IDR          NUMBER(38,10),
  CONSTRAINT PK_T_LOSSRECORDLIST PRIMARY KEY (ID),
  CONSTRAINT FK_LOSSRECORDLIST_KEPALA FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_LOSSRECORD (ID)
)
/
CREATE INDEX {skema}.IX_LOSSRECORDLIST_PARENT ON {skema}.T_LOSSRECORDLIST (PARENT_ID, SEQ_NO)
/
CREATE SEQUENCE {skema}.SEQ_T_LOSSRECORDLIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
