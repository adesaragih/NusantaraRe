# Struktur Tabel — Treaty Description

**Keputusan work owner 05-10-2026:** modul `treatydescription` mengelola master jenis klausul treaty di tabel warisan
`POOLDATA.TREATYDESC` — tambah dan ubah, tanpa hapus. Nol tabel baru, nol DDL.

Katalog DEV 05-10-2026 (dibaca langsung, SELECT saja): 13 baris; tanpa PK, constraint, indeks, maupun trigger; satu
objek bergantung (prosedur `PEGA_TREATYDESC`).

## TREATYDESC

| Kolom | Tipe | Null | Isi |
| --- | --- | --- | --- |
| `ID` | VARCHAR2(100) | ya | `'1' || LPAD(TREATY_DESCRIPTION_SEQ.NEXTVAL, 4, '0')` (`PEGA_TREATYDESC`); DEV 10001-10018, celah 10002, 10005-10007, 10010 |
| `DESCNAME` | VARCHAR2(100) | ya | nama jenis klausul, huruf besar |
| `ISXOL` | VARCHAR2(100) | ya | `0` Non XOL / `1` XOL — saringan grid Treaty Contract Out; DEV: seluruhnya `0` |
| `STATUSAKTIF` | VARCHAR2(100) | ya | `1` aktif / `0` nonaktif; DEV: 7 baris `1`, 6 baris NULL (dibaca aktif), nol `0` |

Sequence `TREATY_DESCRIPTION_SEQ`: MIN 0, INCREMENT 1, CACHE 0, LAST_NUMBER 19 (ID berikutnya `10019`). Dipakai juga
`PEGA_M_TREATYDESC` (dengan awalan situs `M_SITE_DATABASE`, menulis `M_TREATYDESC` JSON) — modul ini tidak memanggil
keduanya dan tidak menyentuh `M_TREATYDESC`.

## Pemakai tabel ini (catatan riset - modul ini tidak membaca maupun mengubah tabel mereka)

| Pemakai | Cara pakai |
| --- | --- |
| `PROPORTIONALARRG` (Treaty Contract Out) | `TREATYDESCID` + salinan `TREATYDESCNAME`; DEV: ke-13 ID dipakai (10001: 1.501 baris); modul ini TIDAK membacanya |
| Treaty Contract Out (kode Go) | membaca seluruh baris urut ID; aturan wajib-isi per ID ditanam di kodenya — jenis baru tampil tetapi belum bisa disimpan di sana |
| Rule Pega klaim (Claim Fac In, Claim Prop, Komite Claim Prop) | `TREATYDESCID = '10001'` literal |
| Rule Pega NB FacIn `GetDataByID_SQL` | `TREATYDESCNAME = 'TREATY LIMIT'` (lewat salinan nama di PROPORTIONALARRG) |
