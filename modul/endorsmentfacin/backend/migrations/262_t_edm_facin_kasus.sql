-- 262 - T_EDM_FACIN_KASUS: kepala kasus Endorsment Fac In (irisan I3, kontrak sesi utama 11-10-2026).
--
-- Satu baris per kasus `EDM-<5 digit>`; BERBAGI PK dengan T_WORK_POLIS.ID (LINI 'FACEDM'), TANPA constraint FK -
-- pola T_NB_OPPORTUNITY (nbfacin 180). T_WORK_POLIS milik premiumlistlife (K-064); constraint lintas modul tidak
-- diputuskan siapa pun.
--
-- Isi = medan yang `SetValueToEDMWork` / `DataToEDM` salin dari kasus portal ke kasus baru (services/bukakasus.go):
-- OLD_POLICY_NO (langkah 12, `OfferFacIn.QuotationData.OldPolicyNo`), EDM_TYPE / EDM_TYPE_NEW / DETAIL_TYPE
-- (`.Quotation.EdmType` / `EdmTypeNew` / `Type`), ENDORSEMENT_DATE (tanggal kalender WIB), ENDORSEMENT_REASON
-- (langkah 13, `.Note` → `OfferFacIn.EndorsmentReason`), sumber / status / admin charge / survey / internal retro,
-- NB_STATUS ("NEW EDM") dan IS_CEDING_CONFIRM ("Policy") dari konektor Start2 flow `InputAddendumFacultativeIn`.
-- Lebar = keputusan agent (lebar Pega tidak tertulis; properti nol rule Property di korpus); panjang diperiksa di Go
-- dalam BYTE sebelum menulis. Nol kolom NUMBER. Seluruh kolom nullable kecuali PK; wajib-isi ditegakkan di Go.
-- Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_EDM_FACIN_KASUS (
  ID                      VARCHAR2(32) NOT NULL,
  OLD_POLICY_NO           VARCHAR2(100),
  EDM_TYPE                VARCHAR2(10),
  EDM_TYPE_NEW            VARCHAR2(10),
  DETAIL_TYPE             VARCHAR2(10),
  ENDORSEMENT_DATE        DATE,
  ENDORSEMENT_REASON      VARCHAR2(4000),
  ENDORSEMENT_SOURCE      VARCHAR2(255),
  ENDORSEMENT_SOURCE_NOTE VARCHAR2(1000),
  ENDORSEMENT_STATUS      VARCHAR2(255),
  ADMIN_CHARGE            VARCHAR2(5),
  SURVEY                  VARCHAR2(5),
  INTERNAL_RETRO          VARCHAR2(10),
  NB_STATUS               VARCHAR2(255),
  IS_CEDING_CONFIRM       VARCHAR2(20),
  CONSTRAINT PK_T_EDM_FACIN_KASUS PRIMARY KEY (ID)
)
/
-- Gerbang 2 (/periksa) mencari kasus EDM yang belum selesai menurut nomor polis lama.
CREATE INDEX {skema}.IX_T_EDM_FACIN_KASUS_POLIS ON {skema}.T_EDM_FACIN_KASUS (OLD_POLICY_NO)
/
