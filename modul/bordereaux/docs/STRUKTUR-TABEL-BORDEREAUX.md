# Struktur Tabel — Bordereaux

**Keputusan work owner 04-10-2026:** modul `bordereaux` mengikuti XML Pega folder korpus `Bordereaux` dengan bug Pega
diperbaiki di Go, **tanpa JSON**: header ke `BORDEREAUX`, detail ke 29 tabel detail warisan, riwayat persetujuan ke
tabel baru `BORDEREAUX_HISTORY` (satu-satunya tabel baru, migrasi 890). Peran baru: tiga baris `M_WORKBASKET`
(migrasi 891).

Katalog dan profil DEV 04-10-2026 (agregat): `BORDEREAUX` 16 berkas (ID `BDX-yyyy.MM.dd.ssSSS`), `M_BORDEREAUX` 16 JSON
(1,6-71,6 KB), detail berisi hanya Premi Engineering 21, Klaim Engineering 21, Klaim Fire 22; tiga prosedur Pega
(`PEGA_INSERTUPDATE_BORDEREAUX`, `PEGA_DELETE_BORDEREAUX`, `PEGA_REPORTBORDERAUX_LOG`) TIDAK dipanggil - isinya
ditulis ulang di Go.

## BORDEREAUX

Tabel warisan Pega, terdaftar `Tabel warisan: dibaca, tidak dibuat` di `MODUL.md`.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `BDX_ID` | teks | tidak | PK | daftar, form | `BDX-` + `yyyy.MM.dd.ssSSS` Asia/Jakarta (XML `InputNew` langkah 3, format diikuti apa adanya); bentrok = milidetik berikutnya |
| `TANGGAL` | date | ya | | daftar | waktu buat berkas (Pega menyimpan `hh` 12 jam sebagai HH24 - diperbaiki) |
| `USER_INPUT` | teks | ya | | hak Edit/Delete | pembuat (`pxCreateOperator`) = `LOGIN_ID` |
| `TYPE` | teks | ya | | daftar, form | PREMIUM / CLAIM / SUBROGATION |
| `TYPE_BUSINESS` | teks | ya | | daftar, form | 14 Business; SUBROGATION = BONDING |
| `MASTERID` | teks | ya | | form | `TREATY_IN.ID` dari popup Choose Master Treaty |
| `CEDINGID` | teks | ya | | form | `TREATY_IN.CEDINGID` |
| `CEDINGNAME` | teks | ya | | daftar, form | `TREATY_IN.CEDING` |
| `SOBID` | teks | ya | | form | `TREATY_IN.LEADINGREINSSOURCEID` |
| `SOBNAME` | teks | ya | | form | `TREATY_IN.LEADINGREINSSOURCE` |
| `TREATYNAME` | teks | ya | | daftar, form | `TREATY_IN.TREATYCONTRACTNAME` |
| `BDXREPORT_START` | date | ya | | daftar, form | wajib; tidak sesudah End |
| `BDXREPORT_END` | date | ya | | daftar, form | wajib |
| `REFFNO_OF_SOA` | teks | ya | | daftar, form | ketikan, maks. 1000 |
| `REFFNO_OF_BDX` | teks | ya | | daftar, form | ketikan, maks. 1000 |
| `POSITION` | teks | ya | | daftar, hak | `LOGIN_ID` pembuat (di tangan pembuat), `Checker`, `Supervisor`, kosong (selesai) |
| `STATUSAKSEP` | teks | ya | | daftar, hak | kosong (draf), `Accept`, `Rejected`, `Resolve-Complete` (`AkseptasiBdx_DT`) |

## BORDEREAUX_HISTORY

Tabel baru (migrasi 890). Padanan page list Pega `CommentList`, yang di Pega hanya ada di dalam JSON.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | bilangan bulat | tidak | PK | — | `SEQ_BORDEREAUX_HISTORY` |
| `BDX_ID` | teks | tidak | | riwayat form | `BORDEREAUX.BDX_ID` - tanpa FK; dihapus bersama berkasnya |
| `TANGGAL` | DATE | tidak | | riwayat form | `SYSDATE` saat Submit / Approve / Reject (`CommentList.Date`) |
| `PIC` | teks | tidak | | riwayat form | `LOGIN_ID` pelaku (`CommentList.OperatorName`) |
| `IS_APPROVED` | teks | tidak | CHECK | riwayat form | `'1'` Approve / Submit, `'0'` Reject (`CommentList.IsApproved`) |
| `KOMENTAR` | teks | ya | | riwayat form | komentar; wajib saat Reject (`CommentList.Suggest`) |

**Index:** PK *(`ID`)*; `IX_BORDEREAUX_HISTORY_BDX` *(`BDX_ID`)*.

## Tabel detail (29 tabel warisan)

`BORDEREAUX_PREMI_<LOB>` (14), `BORDEREAUX_CLAIM_<LOB>` (14), `BORDEREAUX_SUBROGATION_BONDING`. Kolom per tabel dan
nomor kolom CSV-nya ada di `backend/models/kombinasi_gen.go` (dibangkitkan dari 29 activity `Mapping*Bdx*` + katalog
DEV + header templat). Setiap tabel punya `ID` (`BDX_DTL-` + `yyyyMMddhhmmssSSS`, format XML `SaveDetailData`) dan
penunjuk induk `BDX_ID` atau `ID_BDX` (tidak seragam - 11 tabel memakai `ID_BDX`; rule hapus Pega yang selalu memakai
`BDX_ID` karena itu gagal di tabel tersebut, diperbaiki di sini). `ZIP_CODE_AI` (FIRE dan ENGINEERING) diisi pekerja
kode pos AI.

Migrasi 892 melebarkan kolom spread dua tabel Aviation dari `NUMBER(10,4)` warisan ke `NUMBER(38,8)`:
`BORDEREAUX_CLAIM_AVIATION.SPREAD_OF_CLAIM_{OR,QS,OTHERS}` dan `BORDEREAUX_PREMI_AVIATION.SPREAD_OF_RISK_{OR,QS,SPL,OTHERS}`.
Isinya nominal, bukan persen (27 tabel lain: NUMBER bebas, sampai 11 digit di DEV), sehingga `NUMBER(10,4)` menolak
setiap nilai sejuta ke atas - Copy Old Data berkas klaim Aviation DEV ditolak karenanya (04-10-2026).

## M_ATTACHMENTBORDEREAUX

Tabel warisan Pega (lampiran berkas), terdaftar `Tabel warisan: dibaca, tidak dibuat` di `MODUL.md`. Katalog DEV
08-10-2026: PK `ID`; 1 baris. Ditulis `AttachDocBdx_Post` / `AttachDocumentBdx_SQL`, dibaca `GetKategoryDocBDX_SQL` dan
`getAttcachmentList`, dihapus `DeleteAttachmentBdx` dan bersama berkasnya.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | popup View File | `yyyyMMddHHmmssSSS` Asia/Jakarta (`AttachDocBdx_Post` b392); bentrok = milidetik berikutnya |
| `BDX_ID` | teks | ya | | semua | `BORDEREAUX.BDX_ID` - tanpa FK |
| `CATEGORY_ID` | teks (20) | ya | | grid, popup | `M_KATEGORIBORDEREAUX.ID` baris grid (`Primary.MASTERID`) |
| `CATEGORY` | teks (100) | ya | | kolom Type popup | `M_KATEGORIBORDEREAUX.NOTE` saat diunggah (`Primary.TYPE`) |
| `FILENAME` | teks (1000) | ya | | File Name | nama asli berkas, tanpa jalur folder peramban |
| `FILEMIMETYPE` | teks (100) | ya | | View Office Online | EKSTENSI huruf kecil (`.pyFileMimeType`; DEV: `csv`) |
| `DATEINPUT` | date | ya | | — | tidak ditulis Pega maupun modul ini |
| `USERNAME` | teks (100) | ya | | — | pengunggah (`pxRequestor.pxUserIdentifier`) |
| `T_STORAGE_ID` | teks (50) | ya | | unduh, hapus | `T_STORAGE_IMAGE.IMAGEID` objek berkas (`inti/backend/penyimpanan`) |

## M_KATEGORIBORDEREAUX

Master kategori lampiran, hanya dibaca (`ID`, `NOTE`; DEV 08-10-2026: satu baris `00000 Others`).

## M_WORKBASKET

Master peran warisan. Migrasi 891 menambah `ReasBordereauxAdmin`, `ReasBordereauxChecker`, `ReasBordereauxSupervisor`;
migrasi 893 membuang `ReasBordereauxAdmin` lagi (Input Data kini hak menu PENUH, keputusan work owner 04-10-2026)
(`IS_ACTIVE` 1). Pemegangnya di `M_LOGIN_GO_WORKBASKET`, diatur Kelola User.
