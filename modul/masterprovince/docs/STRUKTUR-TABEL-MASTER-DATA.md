# Struktur Tabel — Master Data

Modul `masterprovince` — pewaris migrasi modul `masterdata` yang dihapus saat Master Data dipecah menjadi delapan modul
menu (keputusan work owner 04-10-2026 "Dihapus, tabel pindah ke Province"; `issues/03-pemecahan-delapan-modul.md`). Tabel
di bawah ditulis kedelapan modul master lewat mesin bersama `inti/backend/master`. Modul ini **membuat tujuh tabel** —
tabel flat pengganti view warisan atas JSON `M_*`, migrasi **880** (dulu nbfacin 196, dipindah karena 196 ditahan
work owner; MD-1): `PROVINCE`, `CITYINPUT`, `DISTRICTINPUT`, `ACCUMULATEDTYPE`, `CZONE`, `ACCUMULATION` — dan **menulis**
tabel warisan yang tidak dibuatnya: `NATION` (status dan jejak ubahnya di `T_MASTER_STATUS`, migrasi 881 / 882 — tabel warisan tidak
diubah strukturnya, MD-2), `OBJECTITEMTYPE` (kolom status `ISACTIVE` sendiri; jejak ubah di `T_MASTER_STATUS`), dan membaca `BRANCH` (rujukan City).
Tabel ketujuh yang dibuatnya: `T_MASTER_STATUS`. View `CITY` / `DISTRICT` / `V_JN_OBJ_ITEM` TETAP view di atas tabel
ini (keputusan work owner "Flat-kan CITYINPUT & DISTRICTINPUT saja"), jadi perubahan master langsung terlihat oleh
pembacanya (popup / saran nbfacin).

Kolom view asal `[terverifikasi]` DDL `D:\migrasi\RNM\DDL\<NAMA>.txt`; tipe `VARCHAR2(4000)` = tipe hasil JSON
dot-notation Oracle `[dugaan]` (A180 nbfacin). Status: `STS_AKTIF` `'1'` aktif / `'0'` nonaktif (M-3); tanpa hapus.

## PROVINCE

Migrasi 880 — tabel flat bernama sama dengan view asalnya; isi disalin dari view.

| Kolom | Tipe | Tipe DDL (migrasi 880) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(4000) | kunci master Province (diisi pengguna, unik) |
| `NATIONID` | teks | VARCHAR2(4000) | rujukan `NATION.ID` |
| `NOTE` | teks | VARCHAR2(4000) | nama provinsi; saran nbfacin |
| `NATIONNAME` | teks | VARCHAR2(4000) | TURUNAN: `NATION.NOTE` ber-ID NATIONID (definisi view asal; M-5) |
| `STS_AKTIF` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | status master: `'1'` aktif (bawaan), `'0'` nonaktif |
| `CREATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pembuat (akun pelaku, = M_LOGIN_GO.LOGIN_ID); kosong untuk baris salinan 880 (MD-7) |
| `TGL_CREATE` | tanggal | DATE | migrasi 882 — tanggal buat (SYSDATE) |
| `UPDATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pengubah terakhir (ubah isian atau status) |
| `TGL_UPDATE` | tanggal | DATE | migrasi 882 — tanggal ubah terakhir (SYSDATE) |

## CITYINPUT

Migrasi 880 — tabel flat bernama sama dengan view asalnya; isi disalin dari view.

| Kolom | Tipe | Tipe DDL (migrasi 880) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(4000) | kunci master City |
| `PROVINCEID` | teks | VARCHAR2(4000) | rujukan `PROVINCE.ID` |
| `PROVINCENAME` | teks | VARCHAR2(4000) | migrasi 883 — nama provinsi, isi awal dari `RW.PROVINCENAME` ber-`CITYNAME = NOTE` (hanya bila tepat satu nilai; MD-11); diubah menu City, bukan turunan |
| `NOTE` | teks | VARCHAR2(4000) | nama kota (view CITY) |
| `BRANCHID` | teks | VARCHAR2(4000) | rujukan `BRANCH.ID` |
| `EMAIL` | teks | VARCHAR2(4000) | disalin dari view; sejak 04-10-2026 TIDAK tampil / diisi menu City ("di menu city kolom email dan mo id di hapus saja") — tambah: NULL, ubah: tidak disentuh |
| `MOID` | teks | VARCHAR2(4000) | disalin dari view; sejak 04-10-2026 TIDAK tampil / diisi menu City — tambah: NULL, ubah: tidak disentuh |
| `JABODETABEKSTATUS` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `STS_AKTIF` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | status master: `'1'` aktif (bawaan), `'0'` nonaktif |
| `CREATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pembuat (akun pelaku, = M_LOGIN_GO.LOGIN_ID); kosong untuk baris salinan 880 (MD-7) |
| `TGL_CREATE` | tanggal | DATE | migrasi 882 — tanggal buat (SYSDATE) |
| `UPDATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pengubah terakhir (ubah isian atau status) |
| `TGL_UPDATE` | tanggal | DATE | migrasi 882 — tanggal ubah terakhir (SYSDATE) |

## DISTRICTINPUT

Migrasi 880 — tabel flat bernama sama dengan view asalnya; isi disalin dari view.

| Kolom | Tipe | Tipe DDL (migrasi 880) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(4000) | kunci master District |
| `CITYID` | teks | VARCHAR2(4000) | rujukan `CITYINPUT.ID` |
| `DISTRICTNAME` | teks | VARCHAR2(4000) | nama kecamatan (view DISTRICT) |
| `STS_AKTIF` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | status master: `'1'` aktif (bawaan), `'0'` nonaktif |
| `CREATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pembuat (akun pelaku, = M_LOGIN_GO.LOGIN_ID); kosong untuk baris salinan 880 (MD-7) |
| `TGL_CREATE` | tanggal | DATE | migrasi 882 — tanggal buat (SYSDATE) |
| `UPDATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pengubah terakhir (ubah isian atau status) |
| `TGL_UPDATE` | tanggal | DATE | migrasi 882 — tanggal ubah terakhir (SYSDATE) |

## ACCUMULATEDTYPE

Migrasi 880 — tabel flat bernama sama dengan view asalnya; isi disalin dari view.

| Kolom | Tipe | Tipe DDL (migrasi 880) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(4000) | kunci master Accumulated Type |
| `ACCUMULATIONTYPE` | teks | VARCHAR2(4000) | nama tipe; saran nbfacin |
| `KEYWORD` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `NOTE` | teks | VARCHAR2(4000) | saringan saran nbfacin `IS NOT NULL` |
| `TYPE` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `STS_AKTIF` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | status master: `'1'` aktif (bawaan), `'0'` nonaktif |
| `CREATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pembuat (akun pelaku, = M_LOGIN_GO.LOGIN_ID); kosong untuk baris salinan 880 (MD-7) |
| `TGL_CREATE` | tanggal | DATE | migrasi 882 — tanggal buat (SYSDATE) |
| `UPDATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pengubah terakhir (ubah isian atau status) |
| `TGL_UPDATE` | tanggal | DATE | migrasi 882 — tanggal ubah terakhir (SYSDATE) |

## CZONE

Migrasi 880 — tabel flat bernama sama dengan view asalnya; isi disalin dari view.

| Kolom | Tipe | Tipe DDL (migrasi 880) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(4000) | kunci master CZone |
| `CODE` | teks | VARCHAR2(4000) | kode; saran nbfacin |
| `DESCRIPTION` | teks | VARCHAR2(4000) | urutan saran nbfacin |
| `GROUPOF` | teks | VARCHAR2(4000) | rujukan `CZONE.CODE` |
| `GROUPOFNAME` | teks | VARCHAR2(4000) | TURUNAN: `CZONE.DESCRIPTION` ber-CODE GROUPOF (definisi view asal) |
| `TGLUPDATE` | teks | VARCHAR2(4000) | disalin dari view; tidak ditulis menu (format belum terverifikasi) |
| `USERID` | teks | VARCHAR2(4000) | disalin dari view; tidak ditulis menu |
| `STS_AKTIF` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | status master: `'1'` aktif (bawaan), `'0'` nonaktif |
| `CREATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pembuat (akun pelaku, = M_LOGIN_GO.LOGIN_ID); kosong untuk baris salinan 880 (MD-7) |
| `TGL_CREATE` | tanggal | DATE | migrasi 882 — tanggal buat (SYSDATE) |
| `UPDATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pengubah terakhir (ubah isian atau status) |
| `TGL_UPDATE` | tanggal | DATE | migrasi 882 — tanggal ubah terakhir (SYSDATE) |

## ACCUMULATION

Migrasi 880 — tabel flat bernama sama dengan view asalnya; isi disalin dari view (baris view `ACCUMULATION` hanya yang `IsActive IS NULL`).

| Kolom | Tipe | Tipe DDL (migrasi 880) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(4000) | DIBUAT backend: `negara-ZIPCODE-lpad(ACCUMULATION_SEQ,6)` (prosedur RDBMASTERACCUMULATION, M-4) |
| `ACCUMULATION` | teks | VARCHAR2(4000) | rujukan `ACCUMULATEDTYPE.ID` |
| `ACCUMULATIONNAME` | teks | VARCHAR2(4000) | TURUNAN: `ACCUMULATEDTYPE.ACCUMULATIONTYPE` ber-ID ACCUMULATION (definisi view asal) |
| `NOTE` | teks | VARCHAR2(4000) | disimpan HURUF BESAR (view asal `UPPER`); Note + ZIPCODE unik (prosedur) |
| `KEYWORD` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `SCOPEAREA` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `CZONE` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `CZONEID` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `PROVINCE` | teks | VARCHAR2(4000) | TURUNAN: `PROVINCE.NOTE` ber-ID PROVINCEID `[dugaan]` (MD-3) |
| `PROVINCEID` | teks | VARCHAR2(4000) | rujukan `PROVINCE.ID` |
| `ZIPCODE` | teks | VARCHAR2(4000) | wajib; JOIN RW popup akumulasi nbfacin |
| `ACCUMULATIONTYPE` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `SYARIAHSTATUS` | teks | VARCHAR2(4000) | disalin dari view; diubah menu master |
| `STS_AKTIF` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | status master: `'1'` aktif (bawaan), `'0'` nonaktif |
| `CREATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pembuat (akun pelaku, = M_LOGIN_GO.LOGIN_ID); kosong untuk baris salinan 880 (MD-7) |
| `TGL_CREATE` | tanggal | DATE | migrasi 882 — tanggal buat (SYSDATE) |
| `UPDATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pengubah terakhir (ubah isian atau status) |
| `TGL_UPDATE` | tanggal | DATE | migrasi 882 — tanggal ubah terakhir (SYSDATE) |

## NATION

Tabel warisan `POOLDATA.NATION`, **tidak dibuat** modul ini (bab "Tabel warisan" `MODUL.md`). Sumber tipe
`[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\NATION.txt`. Strukturnya TIDAK diubah; status dan jejak ubah di `T_MASTER_STATUS`. Ditulis menu master Nation;
dibaca saran Country popup akumulasi nbfacin. Kelas Int-NATION → tabel NATION `[dugaan]` (nama dan kolom sama persis).

| Kolom | Tipe | Tipe DDL (companydetail 803) | Dipakai |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(10) NOT NULL | kunci master (diisi pengguna, unik) |
| `OLDID` | teks | VARCHAR2(6) | ditulis menu |
| `NOTE` | teks | VARCHAR2(100) | nama negara (wajib) |
| `NATIONINITIAL` | teks | VARCHAR2(20) | ditulis menu; ekstra saran nbfacin |

Sejak merge `origin/dev` 04-10-2026, `NATION` bukan lagi view/tabel warisan: `companydetail` 802 / 803 membuatnya tabel
datar (`modul/companydetail/docs/STRUKTUR-TABEL-COMPANYDETAIL.md`). Tipe di atas = DDL 803; kolomnya sama persis.

## T_MASTER_STATUS

Migrasi 881 (MD-2) — status aktif / nonaktif master yang tabelnya warisan TANPA kolom status (kini hanya `NATION`), dan
sejak 882 jejak ubah (MD-7) kedua tabel warisan yang tidak di-ALTER: `NATION` dan `OBJECTITEMTYPE`.
Satu baris per (tabel, ID) yang ditambah, diubah, atau diubah statusnya lewat menu (sejak 882; 881 saja: hanya
ubah status); **tanpa baris = aktif** (baris NATION yang disalin dari sistem lama). Pembaca (saran nbfacin) menyaring
`NAMA_TABEL = 'NATION'` berstatus selain `'1'`.

| Kolom | Tipe | Tipe DDL (migrasi 881, diubah 882) | Isi |
| --- | --- | --- | --- |
| `NAMA_TABEL` | teks | VARCHAR2(30) NOT NULL | nama tabel master (`NATION` / `OBJECTITEMTYPE`); PK bersama ID_BARIS |
| `ID_BARIS` | teks | VARCHAR2(4000) NOT NULL | ID baris master (NATION.ID / OBJECTITEMTYPE.ID); dilebarkan 882 dari VARCHAR2(100) |
| `STS_AKTIF` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | `'1'` aktif, `'0'` nonaktif (NATION; diabaikan untuk OBJECTITEMTYPE yang ber-ISACTIVE) |
| `CREATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pembuat (akun pelaku, = M_LOGIN_GO.LOGIN_ID); tabel warisan NATION / OBJECTITEMTYPE (MD-7) |
| `TGL_CREATE` | tanggal | DATE | migrasi 882 — tanggal buat (SYSDATE) |
| `UPDATE_OP` | teks | VARCHAR2(64) | migrasi 882 — pengubah terakhir (ubah isian atau status) |
| `TGL_UPDATE` | tanggal | DATE | migrasi 882 — tanggal ubah terakhir (SYSDATE) |

## OBJECTITEMTYPE

Tabel warisan `POOLDATA.OBJECTITEMTYPE`, **tidak dibuat** modul ini. Sumber tipe `[terverifikasi]`: DDL
`D:\migrasi\RNM\DDL\OBJECTITEMTYPE.txt` (seluruh kolom VARCHAR2(4000 BYTE) kecuali ISACTIVE). View `V_JN_OBJ_ITEM`
(dibaca nbfacin tiket 39) TETAP view di atasnya. Status = kolomnya sendiri `ISACTIVE`: `'1'` aktif `[terverifikasi]`
(`RDBList\GetObjectItembyName_SQL.xml`, `GetDataObjectItem.xml`: `ISACTIVE = '1'`); `'0'` untuk nonaktif `[dugaan]`.

| Kolom | Tipe DDL | Dipakai |
| --- | --- | --- |
| `ID` | VARCHAR2(4000 BYTE) | kunci master |
| `OBJECTITEMTYPE` | VARCHAR2(4000 BYTE) | nama (wajib) |
| `NOTE`, `PCTADJUSTABLE1`, `PCTADJUSTABLE2`, `OBJECTITEMTYPEINA`, `"GROUP"`, `TYPE` | VARCHAR2(4000 BYTE) | ditulis menu |
| `ISACTIVE` | VARCHAR2(10 BYTE) | status master |

## BRANCH

Tabel warisan `POOLDATA.BRANCH`, **baca saja**. Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\BRANCH.txt`.
Hanya `ID` yang dibaca: rujukan `CITYINPUT.BRANCHID` (master City).

Struktur lengkapnya digambarkan `modul/marketingofficer/docs/STRUKTUR-TABEL-MARKETINGOFFICER.md` (modul yang
menyatakannya tabel warisan sejak merge `origin/dev` 04-10-2026; satu tabel bersama = satu bentuk). Yang dipakai di sini:

- `ID` — VARCHAR2(4000 BYTE); pemeriksaan rujukan BRANCHID.
