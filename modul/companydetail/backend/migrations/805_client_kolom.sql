-- 805 - kolom tambahan CLIENT untuk layar Company Detail (keputusan work owner 03/04-10-2026: tanpa dokumen Pega,
-- tabel datar saja). Tabel warisan; kolom lamanya tidak disentuh.
--   PARENT_ID  - Parent organization: ID baris CLIENT induk (namanya tetap di GROUPNAME, dibaca NB FacIn)
--   NOTE       - Note
--   CREATED_*  - jejak pembuatan lewat aplikasi Go (LOGIN_ID, waktu)
--   UPDATED_*  - jejak perubahan terakhir lewat aplikasi Go
-- Established Date, Number of Employees, Owner, dan Client Status TIDAK disimpan (keputusan work owner 04-10-2026).
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.CLIENT ADD (
  PARENT_ID  VARCHAR2(100),
  NOTE       VARCHAR2(4000),
  CREATED_BY VARCHAR2(64),
  CREATED_AT TIMESTAMP(6),
  UPDATED_BY VARCHAR2(64),
  UPDATED_AT TIMESTAMP(6)
)
/
