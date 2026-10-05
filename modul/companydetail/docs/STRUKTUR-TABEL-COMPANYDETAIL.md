# Struktur Tabel — Company Detail

**Keputusan work owner 03/04-10-2026:** modul `companydetail` mengelola organisasi (Create dan ubah) di tabel datar
`POOLDATA.CLIENT`, `CLIENT_PICLIST`, dan `CLIENT_ADDRESS` — *"aku tidak mau ada json lagi"*. Layar SFAGIS Company
Detail **tidak ada di korpus XML**; acuannya screenshot Pega dari work owner dan katalog DEV (agregat, 03/04-10-2026).

Yang DIBUAT modul ini: tabel `M_ENUMERASI` (800-801), tabel datar `NATION` pengganti view-nya (802-804), dan
sequence `SEQ_CLIENT_ORG` nomor ORG (810, bab di bawah). Yang
DIUBAH: kolom tambahan `CLIENT` (805), `CLIENT_PICLIST` (806), `CLIENT_ADDRESS` (807), dan dua trigger `M_CLIENT`
dinonaktifkan (808). Bab tabel warisan di bawah hanya memuat kolom yang DIBUAT migrasi modul ini
(`TestKolomDDLCocokDenganStruktur`); kolom lamanya di sub-bab "kolom warisan".

## M_ENUMERASI

Pilihan dropdown Company Detail — *"jangan ada pake datapega, kalo mau kamu buat table baru ke pooldata"* (work owner
04-10-2026). Isinya disalin SEKALI dari enumerasi Pega DEV ke berkas migrasi 801 (213 baris); aplikasi dan migrasi
tidak membaca skema enumerasi Pega.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `JENIS` | teks | tidak | PK | nama daftar | tipe enumerasi Pega asal: `title`, `bidangusaha`, `posisi`, `jenisalamat`, `telfax`, `kodehp`, `gender` |
| `KODE` | teks | tidak | PK | nilai yang disimpan di tabel client | nama enumerasi Pega; gender `1`/`2` (data PIC lama) |
| `LABEL` | teks | ya | | teks layar | nilai enumerasi Pega; kosong di sumber = NULL (layar menampilkan kodenya) |
| `AKTIF` | teks | tidak | | `1` tampil di dropdown, `0` hanya menampilkan nilai lama | title status Pega; jenisalamat kelompok `cif`; telfax 3, 5 (screenshot Pega 04-10-2026; EMAIL 6 dinonaktifkan 809) |
| `URUTAN` | bilangan bulat | tidak | | urutan dropdown | urutan kode |

| JENIS | Baris | Aktif | Dipakai di |
| --- | --- | --- | --- |
| `title` | 7 | PT., CV., PD., UD. | `CLIENT.TITLE` (LABEL disimpan, seperti `RDBINSERTCLIENT`) |
| `bidangusaha` | 119 | semua | `CLIENT.BU_ID` |
| `posisi` | 20 | semua | saran Position (teks bebas: 123 dari 144 Position lama di luar daftar ini) |
| `jenisalamat` | 8 | 1 RUMAH, 2 KANTOR, 7 EMAIL, 8 KORESPONDENSI | `CLIENT_ADDRESS.ASMADDRESSTYPE` |
| `telfax` | 6 | 3 MOBILE PHONE, 5 OFFICE PHONE (6 EMAIL dinonaktifkan migrasi 809, work owner 04-10-2026) | `CLIENT_ADDRESS.TELFAX_TYPE` |
| `kodehp` | 51 | semua | `CLIENT_ADDRESS.TELFAX_CODE` |
| `gender` | 2 | 1 Male, 2 Female | `CLIENT_PICLIST.GENDER` |

## NATION

Dulu VIEW di atas `M_NATION` (`SELECT a.JSONDATA.ID, OLDID, a.JSONDATA.Note, a.JSONDATA.NationInitial`), kini tabel
datar — *"ambil dari select * from NATION; sebelum itu ubah view itu jadi flat table"* (work owner 04-10-2026). Nama
dan kolom SAMA dengan view-nya, jadi pembaca lama (Pega `BrowseNation_RD`, `ConvertNationality`, dan lainnya) tetap
jalan. ⚠️ Sesudah 804, negara yang ditambah atau diubah lewat form Pega (`InputNation`, `UpdateMasterNation`) masuk
`M_NATION` saja, tidak ke `NATION`.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | — | ID dokumen `M_NATION` (DEV: 57 baris, 6 digit, unik) |
| `OLDID` | teks | ya | | `CLIENT.COUNTRY` | `M_NATION.OLDID` (DEV: 42 terisi; 18.326/18.326 `COUNTRY` Org cocok) |
| `NOTE` | teks | ya | | `CLIENT.COUNTRYNAME`, teks dropdown COUNTRY | nama negara dokumen `M_NATION` |
| `NATIONINITIAL` | teks | ya | | — | inisial negara dokumen `M_NATION` |

## CLIENT

Tabel warisan Pega, satu baris per organisasi/kontak (DEV: 18.414 baris; `FLAG` `Org` 18.345). **Tanpa PK atau
unique** — satu index non-unik. Dibaca 5 view (`ASURADUR`, `B2B_CLIENT`, `EXCLUDEREINSSOA`, `HISTORY_B2B_CLIENT`,
`LIFE_PREMIUM_DETAIL`), prosedur `RDBINSERTCLIENT`/`RDBMCLIENTUPDATE`, rule NB FacIn/RNW/Endorsment
(`GetGroupName_SQL`, `GetInsuredID`, `GetNPWPbyCedingCo_SQL`, `GetNPWPbySOB_SQL`), dan modul Go `premiumlistlife`.
Kolom yang DIBUAT migrasi `805_client_kolom.sql`:

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `PARENT_ID` | teks | ya | | Parent organization | `CLIENT.ID` organisasi induk (pindah: dokumen `ParentID`) |
| `NOTE` | teks | ya | | Note | isian layar (pindah: dokumen `pyDescription`) |
| `CREATED_BY` | teks | ya | | jejak | `LOGIN_ID` pembuat lewat aplikasi Go |
| `CREATED_AT` | timestamp | ya | | jejak | `SYSTIMESTAMP` saat Create |
| `UPDATED_BY` | teks | ya | | jejak | `LOGIN_ID` pengubah terakhir lewat aplikasi Go |
| `UPDATED_AT` | timestamp | ya | | jejak | `SYSTIMESTAMP` setiap simpan |

### Kolom warisan CLIENT (tidak dibuat migrasi mana pun)

| Kolom | Isi | Ditulis modul ini |
| --- | --- | --- |
| `ID` | `ASM-SFAGIS-WORK-ORG ORG-n` | Create: n = `SEQ_CLIENT_ORG.NEXTVAL` (810); tidak berubah |
| `IDVIEW` | `ORG-n` | Create; tidak berubah |
| `FLAG` | `Org` | Create |
| `NAME` | Organization Name | ya |
| `TITLE` | LABEL title (`PT.`) | ya |
| `NPWP` | NPWP | ya |
| `COUNTRY` / `COUNTRYNAME` | `NATION.OLDID` / `NATION.NOTE` | ya |
| `BU_ID` | kode Business Field | ya |
| `GROUPNAME` | nama organisasi induk | ya, disalin dari induk |
| `BU_NOTE` | kategori bisnis (RE/LIFE/GENERAL INSURANCE, BROKER RE) — dipakai pencarian pemegang polis | **tidak** |
| `IDNUMBER`, `OLDID`, `MCLID`, `MCL_TGL_LAHIR` | — | **tidak** |

## CLIENT_PICLIST

PIC organisasi (DEV: 1 baris; 213 PIC di dokumen `M_CLIENT` — trigger lama langsung menghapus PIC tanpa
`ASMUserIdentifier`, 212 dari 213). Tanpa PK. Kolom yang DIBUAT migrasi `806_client_piclist_kolom.sql`:

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `GENDER` | teks | ya | | Gender | `M_ENUMERASI` gender: `1` Male, `2` Female |

### Kolom warisan CLIENT_PICLIST (tidak dibuat migrasi mana pun)

| Kolom | Isi |
| --- | --- |
| `CLIENTID` | `CLIENT.ID` |
| `USERIDENTIFIER` | kunci baris di organisasinya; PIC baru `PIC-n` (nomor tertinggi organisasi itu + 1) |
| `NICKNAME` | Name |
| `POSITION` | Position, teks bebas |
| `EMAIL` | Email |
| `DATEOFBIRTH` | Date of birth, `YYYYMMDD` |
| `PHONENUMBER` | Phone number |

## CLIENT_ADDRESS

Alamat organisasi (DEV: 233 baris = 233 alamat dokumen; kunci `CLIENTID` + `ASMADDRESS` tanpa ganda; tanpa PK).
**Satu baris per nomor Phone and Fax** (work owner 04-10-2026): alamat bernomor tiga = tiga baris beralamat sama;
alamat tanpa nomor = satu baris dengan kolom telfax kosong. Kolom yang DIBUAT migrasi `807_client_address_kolom.sql`:

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `TELFAX_TYPE` | teks | ya | | jenis nomor | `M_ENUMERASI` telfax: 3 MOBILE PHONE, 5 OFFICE PHONE (nomor lama: semua jenis) |
| `TELFAX_CODE` | teks | ya | | kode area | `M_ENUMERASI` kodehp |
| `TELFAX_NO` | teks | ya | | nomor | isian layar (maks. 37 karakter di data lama) |

### Kolom warisan CLIENT_ADDRESS (tidak dibuat migrasi mana pun)

| Kolom | Isi |
| --- | --- |
| `CLIENTID` | `CLIENT.ID` |
| `ASMADDRESSTYPE` | Type, `M_ENUMERASI` jenisalamat |
| `ASMADDRESS` | Address; unik di satu organisasi |
| `ASMCITY`, `CITYNAME`, `ASMZIPCODE`, `DISTRICTNAME`, `PROVINCENAME`, `RWNAME` | tidak di layar; dibawa dari baris lama alamat itu saat disimpan ulang |
| `PXCREATEOPERATOR`, `PXCREATEDATETIME` | jejak Pega; alamat baru: `LOGIN_ID` pelaku dan waktu bentuk Pega (`YYYYMMDDTHHMMSS.mmm GMT`) |

## M_CLIENT

Dokumen organisasi Pega SFAGIS (DEV: 27.132 baris), terdaftar `Tabel warisan: dibaca, tidak dibuat`. Kolom `ID`-nya
dibaca SEKALI oleh migrasi 810 (nilai awal `SEQ_CLIENT_ORG`), dan dua trigger `AFTER INSERT OR UPDATE`-nya
dinonaktifkan (808): keduanya
menulis ulang `CLIENT_ADDRESS` dan `CLIENT_PICLIST` dari dokumen setiap kali Pega menyimpannya (MERGE lalu DELETE).
Isinya dibaca SEKALI oleh alat pindah (`backend/alat`).

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(100) | `ASM-SFAGIS-WORK-ORG ORG-n` untuk organisasi |
| `OLDID` | VARCHAR2(100) | — |

## SEQ_CLIENT_ORG

Sequence nomor ORG organisasi baru (migrasi 810; perintah work owner 04-10-2026: *"buat seq aja, start-nya dari id
max+1"*). Create memakai `NEXTVAL` sebagai `n` di `IDVIEW` `ORG-n` dan `ID` `ASM-SFAGIS-WORK-ORG ORG-n`; `LOCK TABLE
CLIENT` tidak lagi dipakai.

| Sifat | Isi |
| --- | --- |
| Nilai awal | nomor ORG tertinggi di `CLIENT.IDVIEW` dan `M_CLIENT.ID` + 1, DIHITUNG di basis data tempat 810 dijalankan (DEV dan PROD masing-masing) |
| Opsi | `INCREMENT BY 1 NOCACHE NOCYCLE` |
| Lubang nomor | nomor yang diambil Create yang gagal tidak kembali; nomor tidak pernah kembar |
| Mundur | `DROP SEQUENCE`; Create menjawab "belum dimigrasi" (ORA-02289) sampai 810 dijalankan lagi |

## M_NATION

Dokumen negara Pega, terdaftar `Tabel warisan: dibaca, tidak dibuat`; sumber salinan sekali `NATION` (804). Form Pega
`InputNation` tetap menulisnya.

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(6) | — |
| `OLDID` | VARCHAR2(6) | kode negara lama = `CLIENT.COUNTRY` |
