-- 806 - kolom tambahan CLIENT_PICLIST: GENDER (layar PIC). Isinya kode M_ENUMERASI jenis gender: 1 Male,
-- 2 Female - sama dengan data PIC lama (keputusan work owner 04-10-2026).
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.CLIENT_PICLIST ADD (
  GENDER VARCHAR2(10)
)
/
