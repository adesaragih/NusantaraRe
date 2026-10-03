-- 148 - M_PRODUCTNAME_LIFE_OUTWARD: empat kolom baris reasuradur outward + rate-nya.
--
-- Keputusan work owner 02-10-2026 OQ-FLAT-08 "pindahkan": uji kering DEV menemukan 189 objek `OutwardList` yang
-- keenam kunci OR-nya kosong tetapi berisi OUTWARDNAMEID / OUTWARDNAME (189) dan OUTWARDRATEID / OUTWARDRATE (187).
-- Kolom ditambah di berkas BARU, bukan di 146: 146 sudah dijalankan di DEV (T_MIGRASI 02-10-2026 15:39) dan
-- T_MIGRASI mencatat nama berkas - mengubah isi 146 tidak pernah sampai ke skema yang sudah memakainya.
-- Lebar = medan JSON yang sama di sisi umum (`OUTWARDNAMEID` 100, `OUTWARDNAME` 1000, `OUTWARDRATEID` 100,
-- `OUTWARDRATE` 500). NULLABLE: baris `On Retention` tidak mengisinya.
ALTER TABLE {skema}.M_PRODUCTNAME_LIFE_OUTWARD ADD (
  OUTWARDNAMEID    VARCHAR2(100),
  OUTWARDNAME      VARCHAR2(1000),
  OUTWARDRATEID    VARCHAR2(100),
  OUTWARDRATE      VARCHAR2(500)
)
/
