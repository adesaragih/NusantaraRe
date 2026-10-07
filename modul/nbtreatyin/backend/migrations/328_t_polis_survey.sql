-- 328 - T_POLIS_SURVEY <- PolicyTreatyIn.QuotationData.SurveyReportList: popup Historical
-- Survey Report (Section/HistoricalSurveyReportDtl). Keputusan work owner 06-10-2026,
-- membatalkan K7 - tabel kesembilan, di luar diagram grilling.
-- NOURUT = nomor urut baris di dalam induknya, unik per induk (ID-11, AC 8, 10);
-- kunci pasangan antar generasi, bukan kunci dagang (ID-13, AC 11).
CREATE TABLE {skema}.T_POLIS_SURVEY (
  ID               VARCHAR2(32) NOT NULL,
  POLIS_ID         VARCHAR2(32) NOT NULL,
  NOURUT           NUMBER(5) NOT NULL,
  DATE_OF_SURVEY   DATE,
  SURVEYED_BY      VARCHAR2(255),
  LOSS_PREVENTION  NUMBER(38,10),
  REMARKS          VARCHAR2(64),
  CONSTRAINT PK_POLIS_SURVEY PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_SURVEY_POLIS FOREIGN KEY (POLIS_ID) REFERENCES {skema}.T_GENERAL_POLIS_TREATY (ID),
  CONSTRAINT UQ_POLIS_SURVEY_NOURUT UNIQUE (POLIS_ID, NOURUT)
)
/
