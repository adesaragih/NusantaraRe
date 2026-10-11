-- 259 - indeks kolom FK anak objek FIRE yang belum berindeks (keputusan work owner 11-10-2026, temuan skrip kosongkan NB-151032).
-- Rentang KEDUA nbfacin 255-259: 259 = nomor TERAKHIR (rentang habis). Ditulis, TIDAK dijalankan agent; work owner / DBA yang
-- menjalankan. Nol COMMIT (ADR-U-0029).
--
-- [terverifikasi migrasi 186-198]: anak objek FIRE lain sudah berindeks - UNIQUE (PARENT_ID) T_PROPERTY / T_RISKLOCATION /
-- T_BUILDINGCONSTRUCTION (186), T_SURROUNDINGRISK (187), T_TABLEOFLIMIT (189), T_COINSDATA (191); indeks T_LOCATIONLIST (186),
-- T_OCCUPATIONLIST / T_COVERAGELIST (189 / 193), T_DEDUCTIBLELIST (194), T_SPREADINGLIST (198). Tiga FK anak di bawah ini TANPA
-- indeks: menghapus baris induk (GantiObjek / Save Data CSV / PATCH objek / skrip DBA) memindai penuh tabel anak per baris induk
-- (sampai high-water mark), dan baca pohon objek case besar (sqlBacaItem / sqlBacaKerugian / sqlBacaFEA) menggabung lewat kolom
-- ini.
--   T_PROPERTYITEMLIST.PARENT_ID -> T_PROPERTY (188), T_LISTCAUSEOFLOSS.PARENT_ID -> T_PROPERTY (191),
--   T_FEALIST.PARENT_ID -> T_LOCATIONLIST (190).
-- Tabel di DEV bisa berisi ±1 juta baris (NB-151032). `ONLINE` (DML tetap jalan selama pembuatan indeks) memerlukan Oracle
-- Enterprise Edition - edisi DEV `belum terverifikasi`, jadi TIDAK dipakai di berkas ini; DBA boleh menambah `ONLINE` di akhir tiap
-- CREATE INDEX bila edisinya mendukung.
CREATE INDEX {skema}.IX_T_PROPERTYITEMLIST_PARENT ON {skema}.T_PROPERTYITEMLIST (PARENT_ID)
/
CREATE INDEX {skema}.IX_T_LISTCAUSEOFLOSS_PARENT ON {skema}.T_LISTCAUSEOFLOSS (PARENT_ID)
/
CREATE INDEX {skema}.IX_T_FEALIST_PARENT ON {skema}.T_FEALIST (PARENT_ID)
/
