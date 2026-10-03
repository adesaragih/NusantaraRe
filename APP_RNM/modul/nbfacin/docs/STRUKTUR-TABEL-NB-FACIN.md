# Struktur Tabel — NB FacIn: PETA TABEL WARISAN yang dibaca + tabel yang dibuat

Modul ini **membuat satu tabel**, `T_NB_OPPORTUNITY` (migrasi 180, tiket 29 / butir 76.3 — bab terakhir), dan **menulis**
baris `T_WORK_POLIS` milik premiumlistlife (K-064; tidak dibuat, tidak dipetakan kolomnya di sini). Selebihnya ia hanya
**membaca** delapan tabel yang sudah ada — enam tabel limit
akseptasi (tiket 20), tabel akun `T_M_ACCOUNT` (tiket 27), dan tabel bisnis `BUSINESS` (tiket 28). Enam tabel limit akseptasi yang sudah ada
di `POOLDATA`, dengan nama tabel dan kolom **verbatim**. Berkas ini **peta**, bukan DDL: hanya kolom yang dibaca
repository (`backend/repository/limit.go`). Ke-enamnya dinyatakan di `MODUL.md` bab "Tabel warisan: dibaca, tidak
dibuat".

Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` (diberikan work owner, 01-10-2026; butir 42).
Kolom lain di DDL (`ID`, `MAX_LIMIT_*`, `BATAS_WAKTU`, `TGL_UPDATE`, `EFFECTIVE_DATE`, `WORKBASKET`, `JABATAN_ATASAN`,
`LIMITTRADE_BOTTOM`) tidak dibaca. ⛔ `NAMA` dan `LOGIN` **tidak pernah** dibaca (CLAUDE.md §4 butir 10, K-025).

## M_LIMIT_PROPERTYY

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga (`LetterNo`) |
| `TEAM_GROUP` | VARCHAR2(20) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_PROPERTY_NON_PREFERREDD

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(20) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_ENGINEERINGG

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_NONPROPANDENGG

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_FINANCIALINS

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga bentuk B (berspasi, mis. `DIREKTUR TEKNIK`) |
| `LIMITBOND_BOTTOM` | NUMBER(*,0) | limit Bond |
| `LIMITCREDITCL_BOTTOM` | NUMBER(*,0) | limit Kredit CL (juga Trade Credit, A26) |
| `LIMITCREDITNCL_BOTTOM` | NUMBER(*,0) | limit Kredit NCL |

## T_M_ACCOUNT

Tiket 27 (popup ChooseAccount). Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\T_M_ACCOUNT.txt` (ditambahkan work
owner 02-10-2026); kelima kolom DDL dibaca — tidak ada kolom lain.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(255 CHAR) | identitas baris akun |
| `GROUPBUSINESSID` | VARCHAR2(32 CHAR) | id group business |
| `GROUPBUSINESS` | VARCHAR2(64 CHAR) | kolom layar Group Business; dicari (A69) |
| `INSUREDID` | VARCHAR2(255 CHAR) | kolom layar Insured ID; dicari; urutan (A72) |
| `INSUREDNAME` | VARCHAR2(64 CHAR) | kolom layar Insured Name; dicari |

## BUSINESS

Tiket 28 (pilihan Class Of Business). Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\BUSINESS.txt`
(`POOLDATA.BUSINESS`, 18 kolom). **Hanya tiga kolom di bawah yang dibaca**; 15 kolom lain tidak disentuh dan tidak
didaftar di sini (tabel warisan — penjaga kolom tidak membandingkannya).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(4000 BYTE) | identitas pilihan; pemutus seri urutan (A75) |
| `NOTE` | VARCHAR2(4000 BYTE) | teks pilihan yang tampil (`pyDisplayProperty .Note` di section `InputLossRecord_Sec`); NULL dibuang (A78); urutan (A74) |
| `BUSINESSGROUPID` | VARCHAR2(4000 BYTE) | saringan `= :1` (RD filter C `.BusinessGroupID = Param.Group`) |

## T_NB_OPPORTUNITY

Tiket 29 (tombol Create opportunity), butir 76.3 — **dibuat** modul ini, migrasi `180_t_nb_opportunity.sql`. Satu baris
per case NB, **berbagi PK** dengan `T_WORK_POLIS.ID` (`NB-<n>`), tanpa constraint FK (pola `T_PREMIUM_LIST`
premiumlistlife). Dasar tabel sendiri `[terverifikasi]`: di Pega opportunity adalah kelas work tersendiri
(`NB FacIn\ReportDefinition\GetListOpportunity.xml`), dan rancangan tabel flat tidak punya tabel untuknya. Tipe =
keputusan agent A81.

| Kolom | Tipe | Tipe DDL (migrasi 180) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(32) NOT NULL | PK = `T_WORK_POLIS.ID` |
| `ESTIMATED_CLOSING_DATE` | DATE | DATE | Estimated Closing Date (kabel `DD-MM-YYYY`) |
| `BUSINESS_PROSPECT_NAME` | teks | VARCHAR2(255) | Business Prospect Name |
| `ACCOUNT_ID` | teks | VARCHAR2(255 CHAR) | `T_M_ACCOUNT.ID` akun terpilih |
| `INSURED_ID` | teks | VARCHAR2(255 CHAR) | `T_M_ACCOUNT.INSUREDID` |
| `GROUP_BUSINESS_ID` | teks | VARCHAR2(32 CHAR) | `T_M_ACCOUNT.GROUPBUSINESSID` |
| `GROUP_BUSINESS` | teks | VARCHAR2(64 CHAR) | `T_M_ACCOUNT.GROUPBUSINESS` |
| `CLASS_OF_BUSINESS` | teks | VARCHAR2(4000 BYTE) | Class Of Business (lebar `BUSINESS.NOTE`) |
| `TYPE_OF_INWARD` | teks | VARCHAR2(255) | Type Of Inward |
| `TYPE_OF_FACULTATIVE` | teks | VARCHAR2(255) | Type Of Facultative (kosong bila bukan Facultative) |
| `PHASE` | teks | VARCHAR2(255) | Phase |
| `STAGE` | teks | VARCHAR2(255) | Stage |
| `OPPORTUNITY_SOURCE` | teks | VARCHAR2(255) | Opportunity Source |
| `BUSINESS_STATUS` | teks | VARCHAR2(255) | Business Status |
| `DESCRIPTION` | teks | VARCHAR2(4000) | Description |
