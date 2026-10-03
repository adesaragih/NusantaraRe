# Struktur Tabel — NB FacIn: PETA TABEL WARISAN yang dibaca + tabel yang dibuat

Modul ini **membuat tiga tabel** — `T_NB_OPPORTUNITY` (migrasi 180, tiket 29 / butir 76.3), dan `T_GENERAL_POLIS` /
`T_QUOTATIONDATA` **sebagian** (migrasi 182/183, tiket 31 / butir 78.4) — dan **menulis**
baris `T_WORK_POLIS` milik premiumlistlife (K-064; tidak dibuat, tidak dipetakan kolomnya di sini). Selebihnya ia hanya
**membaca** sembilan tabel yang sudah ada — enam tabel limit
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

## T_GENERAL_POLIS

Tiket 31 (blok General layar Inward Facultative), butir 78.4 — **dibuat sebagian** modul ini, migrasi
`182_t_general_polis.sql`: kolom sistem + kolom yang dipakai layar; kolom rancangan lain (termasuk kolom uang yang
menunggu keputusan presisi tim inti) ditambah tiket 23 lewat `ALTER`. Nama/tipe = rancangan (`loader/skema_gen.go`,
dijaga `TestMigrasiFlatSebagianCocokRancangan`). Berbagi PK dengan `T_WORK_POLIS` (K-064), tanpa FK lintas modul.

| Kolom | Tipe | Tipe DDL (migrasi 182) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(32) NOT NULL | PK = `T_WORK_POLIS.ID` (butir 76.1) |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader untuk data lama) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `START_DATE_TIME` | teks | VARCHAR2(30) | Begin date — `.PolicyData.StartDateTime`, teks Pega `YYYYMMDDTHHMMSS.mmm GMT`, ditulis 12:00 WIB = `T050000.000 GMT` (butir 78.1) |
| `OFFERING_DATE` | teks | VARCHAR2(30) | Offering date — `.PolicyData.OfferingDate`, teks Pega `YYYYMMDD` |
| `END_DATE_TIME` | teks | VARCHAR2(30) | End date — `.PolicyData.EndDateTime`, teks Pega |
| `FOLLOWING` | teks | VARCHAR2(50) | Old Policy Number — `.Following` (sel 72), tampil saja |

## T_QUOTATIONDATA

Tiket 31, butir 78.4 — **dibuat sebagian**, migrasi `183_t_quotationdata.sql` (+ `SEQ_T_QUOTATIONDATA`). Satu baris
per case (`UQ_T_QUOTATIONDATA_PARENT`, A87); `PARENT_ID` → `T_GENERAL_POLIS.ID` (FK tanpa `ON DELETE`).

| Kolom | Tipe | Tipe DDL (migrasi 183) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_QUOTATIONDATA` (A92) |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (butir 76.1) |
| `NO_OFFER_SLIP` | teks | VARCHAR2(50) | Reff. number — `.QuotationData.NoOfferSlip` (sel 9) |
| `QQ_NAME` | teks | VARCHAR2(500) | QQ name — `.QuotationData.QQName` (sel 20) |
| `POLICY_TYPE` | teks | VARCHAR2(50) | Policy Type — `.QuotationData.PolicyType` (sel 26); label layar apa adanya (butir 78.2) |
| `MOID` | teks | VARCHAR2(50) | Marketing Name — `.QuotationData.MOID` (sel 75) = `MARKETINGOFFICER.ID` |
| `EDM_DAY` | teks | VARCHAR2(50) | Day — `.QuotationData.EDMDay` (sel 78) |
| `TYPE_FACULTATIVE` | teks | VARCHAR2(50) | Type facultative — `.QuotationData.TypeFacultative` (sel 43); apa adanya (butir 78.3) |
| `SOB_NAME` | teks | VARCHAR2(500) | Source of business — `.QuotationData.SobName` (sel 48), tampil saja |
| `CEDING_CO_NAME` | teks | VARCHAR2(500) | Ceding co name — `.QuotationData.CedingCoName` (sel 49), tampil saja |
| `GROUP_NAME` | teks | VARCHAR2(500) | Group Name — `.QuotationData.GroupName` (sel 56), tampil saja |

## MARKETINGOFFICER

Tiket 31 (pilihan Marketing Name). Tabel warisan `POOLDATA`, **baca saja**. Sumber tipe `[terverifikasi]`: DDL
`D:\migrasi\RNM\DDL\MARKETINGOFFICER.txt` (15 kolom; `PEGA_MARKETINGOFFICER.txt` adalah prosedur penulis tabel
ini, bukan tabel). Hanya tiga kolom di bawah yang dibaca.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(100) | nilai pilihan (`.ID`) → disimpan ke `T_QUOTATIONDATA.MOID`; urutan DESC |
| `CLIENTNAME` | VARCHAR2(100) | teks pilihan (pyPrompt `.ClientName`) |
| `MOSTATUS` | VARCHAR2(100) | saringan RD `.MOStatus = "1"` |
