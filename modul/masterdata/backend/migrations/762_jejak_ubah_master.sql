-- 762 - Master Data: jejak ubah (MD-7, keputusan work owner 04-10-2026 "Kolom diubah oleh/tanggal": pembuat, tanggal
-- buat, pengubah terakhir, tanggal ubah terakhir - hanya perubahan terakhir yang tercatat).
--
-- Nama / tipe mengikuti pola T_WORK_* (claimlife 001, premiumlistlife 059): CREATE_OP VARCHAR2(64) (= lebar
-- M_LOGIN_GO.LOGIN_ID), TGL_CREATE DATE, ditambah pasangan UPDATE_OP VARCHAR2(64), TGL_UPDATE DATE (MD-10). Kolom boleh
-- kosong: baris yang disalin 760 belum punya pembuat. Tabel flat milik modul ini mendapat kolomnya sendiri; jejak NATION
-- dan OBJECTITEMTYPE (warisan, tidak di-ALTER - MD-2) disimpan di T_MASTER_STATUS, yang ID_BARIS-nya dilebarkan ke
-- VARCHAR2(4000) = lebar OBJECTITEMTYPE.ID. Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.PROVINCE ADD (
  CREATE_OP  VARCHAR2(64),
  TGL_CREATE DATE,
  UPDATE_OP  VARCHAR2(64),
  TGL_UPDATE DATE
)
/
ALTER TABLE {skema}.CITYINPUT ADD (
  CREATE_OP  VARCHAR2(64),
  TGL_CREATE DATE,
  UPDATE_OP  VARCHAR2(64),
  TGL_UPDATE DATE
)
/
ALTER TABLE {skema}.DISTRICTINPUT ADD (
  CREATE_OP  VARCHAR2(64),
  TGL_CREATE DATE,
  UPDATE_OP  VARCHAR2(64),
  TGL_UPDATE DATE
)
/
ALTER TABLE {skema}.ACCUMULATEDTYPE ADD (
  CREATE_OP  VARCHAR2(64),
  TGL_CREATE DATE,
  UPDATE_OP  VARCHAR2(64),
  TGL_UPDATE DATE
)
/
ALTER TABLE {skema}.CZONE ADD (
  CREATE_OP  VARCHAR2(64),
  TGL_CREATE DATE,
  UPDATE_OP  VARCHAR2(64),
  TGL_UPDATE DATE
)
/
ALTER TABLE {skema}.ACCUMULATION ADD (
  CREATE_OP  VARCHAR2(64),
  TGL_CREATE DATE,
  UPDATE_OP  VARCHAR2(64),
  TGL_UPDATE DATE
)
/
ALTER TABLE {skema}.T_MASTER_STATUS ADD (
  CREATE_OP  VARCHAR2(64),
  TGL_CREATE DATE,
  UPDATE_OP  VARCHAR2(64),
  TGL_UPDATE DATE
)
/
ALTER TABLE {skema}.T_MASTER_STATUS MODIFY (
  ID_BARIS VARCHAR2(4000)
)
/
