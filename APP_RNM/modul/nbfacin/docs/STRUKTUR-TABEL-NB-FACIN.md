# Struktur Tabel — NB FacIn: PETA TABEL WARISAN yang dibaca + tabel yang dibuat

Modul ini **membuat enam belas tabel** — `T_NB_OPPORTUNITY` (migrasi 180, tiket 29 / butir 76.3), `T_GENERAL_POLIS` /
`T_QUOTATIONDATA` **sebagian** (migrasi 182/183, tiket 31 / butir 78.4), `T_CEDINGCOLIST` utuh (185, tiket 34), dan tabel tab Object FIRE (186, tiket 35: `T_LOCATIONLIST` /
`T_PROPERTY` sebagian, `T_RISKLOCATION` / `T_BUILDINGCONSTRUCTION` utuh; 187, tiket 38: `T_SURROUNDINGRISK` utuh; 188, tiket 39: `T_PROPERTYITEMLIST` sebagian; 189, tiket 40: `T_OCCUPATIONLIST` sebagian, `T_TABLEOFLIMIT` utuh; 190, tiket 41: `T_FEALIST` baru; 191, tiket 42: `T_LISTCAUSEOFLOSS`, `T_COINSDATA` utuh; 193, tiket 43: `T_COVERAGELIST` sebagian) — dan **menulis**
baris `T_WORK_POLIS` milik premiumlistlife (K-064; tidak dibuat, tidak dipetakan kolomnya di sini). Selebihnya ia hanya
**membaca** delapan belas tabel / view yang sudah ada — enam tabel limit
akseptasi (tiket 20), tabel akun `T_M_ACCOUNT` (tiket 27), tabel bisnis `BUSINESS` (tiket 28), `MARKETINGOFFICER` (31), `AGENT`
(33), `RISKADDRESS` (36; juga DISISIPI tiket 37), `RW` (36/37), `OCCUPATION` (38), `V_JN_OBJ_ITEM` dan `CURRENCY` (39), `TABLEOFLIMIT` (40), view `COVERAGE_FACIN` dan `COVERAGE` (43) — sama dengan bab "Tabel warisan" `MODUL.md`. Enam tabel limit akseptasi yang sudah ada
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
| `SOURCE_OF_BUSINESS` | teks | VARCHAR2(50) | kode SOB — `.QuotationData.SourceOfBusiness`; migrasi **184** (`ALTER … ADD`, tiket 33) |
| `SOB_NAME` | teks | VARCHAR2(500) | Source of business — `.QuotationData.SobName` (sel 48); ditulis server dari `AGENT` menurut kode (tiket 33, E-4) |
| `CEDING_CO` | teks | VARCHAR2(1000) | gabungan `;` kode Ceding Co — `.QuotationData.CedingCo`; migrasi **185** (`ALTER … ADD`, butir 80) |
| `CEDING_CO_NAME` | teks | VARCHAR2(4000) | Ceding co name — `.QuotationData.CedingCoName` (sel 49), gabungan `;` nama dari `AGENT` (tiket 34); 183 VARCHAR2(500), dilebarkan migrasi **185** (butir 80) |
| `GROUP_NAME` | teks | VARCHAR2(500) | Group Name — `.QuotationData.GroupName` (sel 56), tampil saja |

## T_CEDINGCOLIST

Tiket 34 (daftar Ceding Co), butir 80 — **dibuat utuh** modul ini, migrasi `185_t_cedingcolist.sql` (+ `SEQ_T_CEDINGCOLIST`,
indeks `IX_CEDINGCOLIST_PARENT`). Tabel RANCANGAN (jalur `QuotationData/CedingCoList`, induk `T_QUOTATIONDATA`); nama/tipe =
rancangan, ID/PARENT_ID NUMBER(19) (A92). Satu baris per ceding, urut pilih; diganti utuh tiap *Save for later*.

| Kolom | Tipe | Tipe DDL (migrasi 185) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_CEDINGCOLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_QUOTATIONDATA.ID` (FK tanpa `ON DELETE`) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan pilih, mulai 1 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID — baris aplikasi v4 acak (A106), baris loader v5 |
| `CEDING_CO` | teks | VARCHAR2(50) | kode `AGENT.ID` — `.CedingCoList(n).CedingCo` |
| `CEDING_CO_NAME` | teks | VARCHAR2(500) | nama dari `AGENT.CLIENTNAME` — `.CedingCoList(n).CedingCoName` |

## T_LOCATIONLIST

Tiket 35 (tab Object FIRE) — migrasi `186_t_objek_fire.sql`, **sebagian** (pola butir 78.4): kolom sistem saja; kolom
uang `LOSS_RATIO*_AMOUNT` (presisi tim inti) dan kolom lain rancangan ditambah tiket 23. Jalur rancangan `LocationList`,
induk `T_GENERAL_POLIS`. Satu baris per objek, urut `SEQ_NO`; diganti utuh tiap Save.

Tiket 42 — migrasi 191 menambah empat kolom loss ratio rancangan (`ALTER TABLE … ADD`); diisi **hasil hitung server** saat
Save (W-4 / `SetLossRatio_Act`), tidak dari layar.

| Kolom | Tipe | Tipe DDL (migrasi 186/191) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_LOCATIONLIST` (A92) |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (butir 76.1; FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan baris grid, mulai 1 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID (aplikasi v4, A106) |
| `LOSS_RATIO1_YEAR_AMOUNT` | angka desimal | NUMBER(38,8) | LR 1 Year = ΣClaim/ΣAmount (≤ 365 hari), 8 desimal (A148) — **uang** |
| `LOSS_RATIO1_YEAR_PERCENT` | teks | VARCHAR2(50) | %LR 1 Years = ΣClaim·100/ΣAmount (8 desimal, A148), teks desimal |
| `LOSS_RATIO35_YEAR_AMOUNT` | angka desimal | NUMBER(38,8) | LR 3 - 5 Years (≤ 1825 hari, kumulatif) — **uang** |
| `LOSS_RATIO35_YEAR_PERCENT` | teks | VARCHAR2(50) | %LR 3 - 5 Years, teks desimal |

## T_PROPERTY

Tiket 35 — migrasi 186, **sebagian**: tanpa `TOTAL_TSI` (uang) — ditambah tiket 23. Tiket 38 — migrasi 187 menambah
empat kolom rancangan `OWNERSHIP`, `IS_PRODUCTION_PROCESS_FLAG`, `IS_HOT_WORK_PROCESS_FLAG`, `IS_FLAMMABLE_ITEM_FLAG`
(`ALTER TABLE … ADD`, tipe rancangan).
Satu baris per lokasi (`UQ_T_PROPERTY_PARENT`). Boolean disimpan teks `true`/`false` (bentuk data fixture).

| Kolom | Tipe | Tipe DDL (migrasi 186/187/192) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_PROPERTY` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_LOCATIONLIST.ID` |
| `ALM_RISK_ID` | teks | VARCHAR2(100) | Risk Address ID — `.Property.AlmRiskID` (sel 39) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `BUILDING_NO` | teks | VARCHAR2(50) | Building No. (sel 30) |
| `COUNTRY` | teks | VARCHAR2(100) | Country (sel 32) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `IS_MATERIAL_DAMAGE` | teks | VARCHAR2(10) | Material Damage (sel 11) `true`/`false` |
| `IS_TOP_RISK` | teks | VARCHAR2(10) | Top Risk (sel 12) `true`/`false` |
| `OBJECT_NAME` | teks | VARCHAR2(500) | Object Name (sel 13) |
| `OBJECT_NO` | teks | VARCHAR2(50) | Object No. (sel 5) |
| `OBJECT_TYPE` | teks | VARCHAR2(50) | Object Type (sel 6, wajib) |
| `PROVINCE` | teks | VARCHAR2(100) | Province (sel 38) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `ROAD_NAME` | teks | VARCHAR2(4000) | Address — `.Property.RoadName` (sel 29) — lebar 4000 sejak migrasi 192 (butir 87/88; semula 500) |
| `ROAD_TYPE` | teks | VARCHAR2(100) | Type — `.Property.RoadType` (sel 28) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `CURRENCY_CODE` | teks | VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL | K-069 (bawaan; tidak ditulis layar) |
| `OWNERSHIP` | teks | VARCHAR2(50) | Ownership — `.Property.Ownership` (`RiskAround` sel 65); kode apa adanya (187) |
| `IS_PRODUCTION_PROCESS_FLAG` | teks | VARCHAR2(10) | `.Property.IsProductionProcessFlag` (sel 77) `true`/`false` (187) |
| `IS_HOT_WORK_PROCESS_FLAG` | teks | VARCHAR2(10) | `.Property.IsHotWorkProcessFlag` (sel 78) `true`/`false` (187) |
| `IS_FLAMMABLE_ITEM_FLAG` | teks | VARCHAR2(10) | `.Property.IsFlammableItemFlag` (sel 79) `true`/`false` (187) |

## T_RISKLOCATION

Tiket 35 — migrasi 186, **utuh**. Satu baris per property (`UQ_T_RISKLOCATION_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 186/192) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_RISKLOCATION` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `ASM_ADDRESS` | teks | VARCHAR2(4000) | Risk Location — `.Property.RiskLocation.ASMAddress` (sel 40; kolom grid "Location") — lebar 4000 sejak migrasi 192 (butir 87/88; semula 50) |
| `ASM_CITY` | teks | VARCHAR2(100) | City (sel 36) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `ASM_DISTRICT` | teks | VARCHAR2(100) | District (sel 37) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `ASMRW` | teks | VARCHAR2(100) | Territory — `.Property.RiskLocation.ASMRW` (sel 35) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `ASM_ZIP_CODE` | teks | VARCHAR2(100) | Zip Code (sel 31) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |

## T_BUILDINGCONSTRUCTION

Tiket 35 — migrasi 186, **utuh** + tiga kolom **baru** (bukan di rancangan; A110, amandemen loader `amandemenBangunan`).
Satu baris per property (`UQ_T_BUILDINGCONSTR_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 186) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_BUILDINGCONSTRUCTION` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `FLOOR_TYPE` | teks | VARCHAR2(50) | Floor Type (sel 52) — kode apa adanya (G-4) |
| `NUMBER_OF_FLOOR` | teks | VARCHAR2(50) | Number of Floor (sel 49) — rancangan teks; diperiksa bilangan bulat ≥ 0 |
| `ROOF_TYPE` | teks | VARCHAR2(50) | Roof Type (sel 50) — kode apa adanya |
| `WALL_TYPE` | teks | VARCHAR2(50) | Wall Type (sel 51) — kode apa adanya |
| `PARTITION_TYPE` | teks | VARCHAR2(50) | **baru** — `.Property.BuildingConstruction.PartitionType` (sel 56) |
| `SUPPORT_WALL_TYPE` | teks | VARCHAR2(50) | **baru** — `.SupportWallType` (sel 57) |
| `OTHERS_TYPE` | teks | VARCHAR2(50) | **baru** — `.OthersType` (sel 58) |

## T_SURROUNDINGRISK

Tiket 38 — migrasi 187, **utuh** (rancangan: jalur `LocationList/Property/SurroundingRisk`, induk `T_PROPERTY`) + 18 kolom
**baru** (bukan di rancangan; A130, amandemen loader `amandemenSekitar`) dari `NB FacIn\Section\RiskAround.xml`.
Satu baris per property (`UQ_T_SURROUNDINGRISK_PARENT`). Kode disimpan apa adanya — server tidak mencocokkan ke daftar
pilihan (aturan properti `DDL\FrontConstruction.xml` / `Ownership.xml` / `FloodAreaStatus.xml` /
`HousekeepingStatus.xml` / `FloodArea.xml` dipakai frontend).

| Kolom | Tipe | Tipe DDL (migrasi 187) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SURROUNDINGRISK` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `FLOOD_AREA_STATUS` | teks | VARCHAR2(50) | Flood Area Status (sel 67) — kode apa adanya |
| `HOUSEKEEPING_STATUS` | teks | VARCHAR2(50) | Housekeeping Status (sel 66) — kode apa adanya |
| `FRONT_OCCUPATION` | teks | VARCHAR2(1000) | **baru** — `.Property.SurroundingRisk.FrontOccupation` (sel 9) = `OCCUPATION.OLDID` (lebar = sumber, butir 80) |
| `FRONT_CONSTRUCTION` | teks | VARCHAR2(500) | **baru** — `.FrontConstruction`; nilai standar terpanjang 215 bita |
| `FRONT_DISTANCE` | teks | VARCHAR2(50) | **baru** — `.FrontDistance`; teks angka 0..100000, ≤ 2 desimal (A129) |
| `FRONT_NOTE` | teks | VARCHAR2(1000) | **baru** — `.FrontNote` (diisi `OCCUPATION.NAME`; lebar = sumber, butir 80) |
| `LEFT_OCCUPATION` | teks | VARCHAR2(1000) | **baru** — sisi Left, pola sama (sel 23–28) |
| `LEFT_CONSTRUCTION` | teks | VARCHAR2(500) | **baru** |
| `LEFT_DISTANCE` | teks | VARCHAR2(50) | **baru** |
| `LEFT_NOTE` | teks | VARCHAR2(1000) | **baru** |
| `BACK_OCCUPATION` | teks | VARCHAR2(1000) | **baru** — sisi Back (sel 37–42) |
| `BACK_CONSTRUCTION` | teks | VARCHAR2(500) | **baru** |
| `BACK_DISTANCE` | teks | VARCHAR2(50) | **baru** |
| `BACK_NOTE` | teks | VARCHAR2(1000) | **baru** |
| `RIGHT_OCCUPATION` | teks | VARCHAR2(1000) | **baru** — sisi Right (sel 51–56) |
| `RIGHT_CONSTRUCTION` | teks | VARCHAR2(500) | **baru** |
| `RIGHT_DISTANCE` | teks | VARCHAR2(50) | **baru** |
| `RIGHT_NOTE` | teks | VARCHAR2(1000) | **baru** |
| `FLOOD_AREA` | teks | VARCHAR2(50) | **baru** — `.FloodArea` (sel 68); kode apa adanya (aturan properti `DDL\FloodArea.xml`: 4 nilai, 1 bita) |
| `HOUSEKEEPING_REMARK` | teks | VARCHAR2(500) | **baru** — `.HousekeepingRemark` (sel 69) |

## T_PROPERTYITEMLIST

Tiket 39 — migrasi 188, **sebagian** (pola A109): rancangan jalur `LocationList/Property/PropertyItemList`, induk
`T_PROPERTY`, **banyak** baris per property urut `SEQ_NO` (= `PROPERTY_ITEM_NO` 1..n). Tanpa kolom yang tidak dipakai layar
ini (`CURRENCY_ID`, `CURRENCY_OLD_ID`, `FLAG_NET_RATE`, `IS_OLD_DATA`, `PERCENTAGE_ADJUSTMENT`, `PROPERTY_ID`,
`SELECTED_LOCATION_ADDRESS`, `SELECTED_OBJECT_ITEM`, `TOTAL_PREMIUM_NUSANTARA_RE`) — ditambah tiket 23; `TOTAL_GROSS_PREMI` /
`TOTAL_NET_RATE` ditambah tiket 43 (migrasi 193 `ALTER ADD`). Enam kolom **baru** (A132, amandemen loader `amandemenItem`). Uang dan
persen `NUMBER(38,8)` (rancangan `NUMBER` polos; ADR-0016), ditulis/dibaca sebagai teks desimal bertitik — nol float.

| Kolom | Tipe | Tipe DDL (migrasi 188) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_PROPERTYITEMLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan item 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `CURRENCY` | teks | VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL | Currency (`CURRENCY.CURRENCY`, tanpa ITL); **wajib** diisi layar (A133; aplikasi tidak menulis `UNKNOWN`, K-012) |
| `IS_ADJUSTABLE_FLAG` | teks | VARCHAR2(10) | Adjustable `true`/`false` |
| `ITEM_TYPE` | teks | VARCHAR2(50) | Object Item Type — `V_JN_OBJ_ITEM.JN_OBJ_ITEM` |
| `ITEM_TYPE_ID` | teks | VARCHAR2(50) | `V_JN_OBJ_ITEM.MJOI_KODE` (nilai dropdown) |
| `PCT_ADJUST2` | angka desimal | NUMBER(38,8) | Adjustment Pct. (dropdown, tidak Adjustable) — persen |
| `PCT_ADJUST_OTHER` | angka desimal | NUMBER(38,8) | Adjustment Pct. (Adjustable: 60..100, `ValidateAdjustPct`) — persen |
| `PROPERTI_ITEM_NOTE` | teks | VARCHAR2(500) | Object Item Note — `V_JN_OBJ_ITEM.KETERANGAN` |
| `PROPERTY_ITEM_NO` | teks | VARCHAR2(50) | nomor item = `SEQ_NO` (K-6) |
| `REMARK` | teks | VARCHAR2(500) | Remark of `<ItemType>` |
| `TSI_OBJECT_ITEM` | angka desimal | NUMBER(38,8) | **UANG** — TSI Object Item (All Unit), ≥ 0, ≤ 8 desimal |
| `PROPERTY_YEAR` | teks | VARCHAR2(50) | **baru** — Year (`.PropertyYear`) |
| `UNIT` | teks | VARCHAR2(50) | **baru** — Unit(s) (`.Unit`), bilangan bulat > 0 |
| `CONDITION` | teks | VARCHAR2(500) | **baru** — Condition (`.Condition`), kode apa adanya |
| `YEAR` | teks | VARCHAR2(50) | **baru** — Year of Planting (`.Year`) |
| `NO_OF_TREE` | teks | VARCHAR2(50) | **baru** — No of Trees (`.NoOfTree`) |
| `AREA_HECTAR` | teks | VARCHAR2(50) | **baru** — Area ( Hectar ) (`.AreaHectar`) |
| `TOTAL_GROSS_PREMI` | angka desimal | NUMBER(38,8) | tiket 43 (193) — **uang**, Σ Premium coverage item, dihitung server saat PUT (`CountPremi_ACT` langkah 55-56); item tanpa coverage = kosong |
| `TOTAL_NET_RATE` | angka desimal | NUMBER(38,8) | tiket 43 (193) — ‰ Total Net Rate, disimpan apa adanya dari PUT (A159; rumusnya tahap C2) |

## T_COVERAGELIST

Tiket 43 — migrasi 193, **sebagian** (pola A109). Rancangan memberi tabel ini **induk jamak** (CargoList / PropertyItemList /
AnekaList / PersonList / VehicleList), dibedakan `PARENT_TABLE` / `SRC_PATH`; NB menulis hanya jalur
`LocationList/Property/PropertyItemList/CoverageList` (`PARENT_TABLE` `T_PROPERTYITEMLIST`). `PARENT_ID` **tanpa FK**
(pola A140), indeks `IX_T_COVERAGELIST_PARENT (PARENT_TABLE, PARENT_ID)`. Banyak baris per item, urut `SEQ_NO`. Seluruh
medan kontrak ADA di rancangan — tanpa kolom baru. Medan uang / rate / persen dihitung server (`CountPremi_ACT`, basis
1–4) dan disimpan setengah-ke-atas 8 desimal (A155).

| Kolom | Tipe | Tipe DDL (migrasi 193) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_COVERAGELIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTYITEMLIST.ID` |
| `PARENT_TABLE` | teks | VARCHAR2(30) | `T_PROPERTYITEMLIST` |
| `SRC_PATH` | teks | VARCHAR2(200) | `LocationList/Property/PropertyItemList/CoverageList` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan coverage 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `COVERAGE` | teks | VARCHAR2(50) | `.Coverage` — ID coverage (`COVERAGE_FACIN.ID`) |
| `OLDID` | teks | VARCHAR2(50) | `.OLDID` — kode tampil |
| `COVERAGE_NOTE` | teks | VARCHAR2(500) | `.CoverageNote` — nama coverage |
| `COVERAGE_BASIS` | teks | VARCHAR2(50) | `.CoverageBasis` — 1..4 (5 Layering ditolak 400, tahap C2–C4) |
| `DAY` | teks | VARCHAR2(50) | `.Day` |
| `TSI` | angka desimal | NUMBER(38,8) | **uang** — server: = `TSI_OBJECT_ITEM` item |
| `INDEMNITY` | teks | VARCHAR2(50) | `.Indemnity` |
| `RATE` | angka desimal | NUMBER(38,8) | ‰ Gross Rate (dihitung balik pada mode amount) |
| `RATE_OJK` | angka desimal | NUMBER(38,8) | ‰ Standard Rate (lookup = tahap C4) |
| `FIRST_LOSS` | teks | VARCHAR2(50) | % First Loss — teks desimal (rancangan VARCHAR2) |
| `DISCOUNT_PERCENTAGE` | angka desimal | NUMBER(38,8) | % Discount |
| `TSI_LIABILITY` | angka desimal | NUMBER(38,8) | **uang** — server, per basis |
| `NET_RATE` | angka desimal | NUMBER(38,8) | ‰ Net Rate (terisi → menggantikan Rate dalam rumus) |
| `LIMITOF_LIABILITY` | angka desimal | NUMBER(38,8) | **uang** — Limit of Liability |
| `PCT_LO_L` | angka desimal | NUMBER(38,8) | % LoL |
| `PRO_RATE_PERCENT` | angka desimal | NUMBER(38,8) | server — Prorate periode polis case |
| `INDEMNITY_PERCENTAGE` | angka desimal | NUMBER(38,8) | % Indemnity |
| `FIRST_SCALE` | teks | VARCHAR2(50) | % First Scale — teks desimal |
| `SUBLIMIT` | angka desimal | NUMBER(38,8) | % Sub Limit |
| `LOST_LIMIT` | teks | VARCHAR2(50) | % Loss Limit — teks desimal (ejaan rancangan / Pega `.LostLimit`) |
| `EML_PML` | teks | VARCHAR2(50) | % EML / PML — teks desimal |
| `DISCOUNT` | angka desimal | NUMBER(38,8) | **uang** — Discount |
| `PREMIUM` | angka desimal | NUMBER(38,8) | **uang** — Gross Premium (server pada mode percent) |
| `CONDITIONS` | teks | VARCHAR2(500) | `.Conditions` |
| `PCT_ADJUSTMENT` | angka desimal | NUMBER(38,8) | server — `CountPremi_ACT` langkah 11-15 (A156); tidak dikirim JSON |
| `CURRENCY_CODE` | teks | VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL | = `CURRENCY` item (K-069 / K-012); item ber-coverage dengan mata uang > 10 byte → 400 |

## T_OCCUPATIONLIST

Tiket 40 — migrasi 189, **sebagian** (pola A109; tanpa `CATEGORY` okupasi, `IS_OLD_DATA`, `LOCATION_ID` — tiket 23).
Rancangan memberi tabel ini **tiga induk** (`LocationList/Property/OccupationList`, `.../RiskLocation/OccupationList`,
`VehicleList/OccupationList`); NB menulis jalur pertama saja: `PARENT_TABLE = 'T_PROPERTY'`, `SRC_PATH` = jalur itu (isi sama
dengan loader). Karena berinduk jamak, `PARENT_ID` **tanpa FK** — indeks `IX_T_OCCUPATIONLIST_PARENT (PARENT_TABLE, PARENT_ID)`.

| Kolom | Tipe | Tipe DDL (migrasi 189) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_OCCUPATIONLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` (bila `PARENT_TABLE = 'T_PROPERTY'`) |
| `PARENT_TABLE` | teks | VARCHAR2(30) | tabel induk — `T_PROPERTY` |
| `SRC_PATH` | teks | VARCHAR2(200) | jalur sumber — `LocationList/Property/OccupationList` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan okupasi 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `OCCUPATION_ID` | teks | VARCHAR2(1000) | Occupation ID (`.OccupationId` = `OCCUPATION.OLDID`); lebar = sumber (A138, rancangan 50) |
| `OCCUPATION_NAME` | teks | VARCHAR2(1000) | Occupation Name (`OCCUPATION.NAME`); lebar = sumber (A138, rancangan 500) |

## T_TABLEOFLIMIT

Tiket 40 — migrasi 189, **utuh**. Halaman `.TableOfLimit` okupasi: satu baris per okupasi (`UQ_T_TABLEOFLIMIT_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 189) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_TABLEOFLIMIT` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_OCCUPATIONLIST.ID` |
| `CATEGORY` | teks | VARCHAR2(50) | Category I/II/III (`SetDataOccupation` dari KDRiskExposure 01/02/03) |
| `DESCRIPTION` | teks | VARCHAR2(500) | Class of Construction |
| `PCT_LIMIT` | teks | VARCHAR2(50) | PctLimit — **teks apa adanya** (butir 68.1: rancangan NUMBER; koma desimal dan spasi ujung dipertahankan); lebar A139 |

## T_FEALIST

Tiket 41 — migrasi 190, tabel **BARU** (rancangan flat tidak punya tabel FEA; A142, amandemen loader `amandemenFEA` +
jalur `LocationList/FEAList` + lipatan `.DataFEA`). `.FEAList` milik **baris lokasi** (kelas
`Data-OfferFacIn-LocationReinsurance`), jadi induknya `T_LOCATIONLIST`; banyak baris per lokasi, urut `SEQ_NO`. Kolom
sistem = pola tabel berulang rancangan (T_ADDITIONALSHIP).

| Kolom | Tipe | Tipe DDL (migrasi 190) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_FEALIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_LOCATIONLIST.ID` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan FEA 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `APAR` | teks | VARCHAR2(50) | APAR (`.APAR`) — jumlah unit, bilangan bulat ≥ 0 (M-2) |
| `SPRINKLER` | teks | VARCHAR2(50) | Sprinkler (`.Sprinkler`) — jumlah unit |
| `SMOKE_DETECTOR` | teks | VARCHAR2(50) | Smoke Detector & Alarm (`.SmokeDetector`) — jumlah unit |
| `HYDRANT` | teks | VARCHAR2(50) | Hydrant (`.Hydrant`) — jumlah unit |
| `PRIVATE_TRUCK_BRIGADE` | teks | VARCHAR2(50) | Private Truck Brigade (Unit) (`.DataFEA.PrivateTruckBrigade`) — jumlah unit |
| `PRIVATE_FIRE_BRIGADE` | teks | VARCHAR2(50) | Private Team Fire Brigade (`.DataFEA.PrivateFireBrigade`) — kode apa adanya |
| `TEAM_SOP_SAFETY` | teks | VARCHAR2(50) | Team & SOP Safety (`.DataFEA.TeamSOPSafety`) — kode apa adanya |
| `TEAM_SOP_RISK_MANAGEMENT` | teks | VARCHAR2(50) | Team & SOP Risk Management (`.DataFEA.TeamSOPRiskManagement`) — kode apa adanya |
| `INFO_FEA` | teks | VARCHAR2(500) | Others Info (`.InfoFEA`) |

## T_LISTCAUSEOFLOSS

Tiket 42 — migrasi 191, **utuh** (rancangan jalur `LocationList/Property/ListCauseOfLoss`, induk `T_PROPERTY`) + lima kolom
**baru** (A145, amandemen loader `amandemenKerugian`). Banyak baris per property, urut `SEQ_NO`.

| Kolom | Tipe | Tipe DDL (migrasi 191) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_LISTCAUSEOFLOSS` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan catatan 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `CLAIM` | angka desimal | NUMBER(38,8) | Total Claim (100%) — **uang** |
| `CURRENCY` | teks | VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL | Currency (`CURRENCY`, ≠ ITL); **wajib** (K-069 + K-012, pola A133) |
| `DETAIL` | teks | VARCHAR2(500) | Loss Detail — dilebarkan dari rancangan 50 (A146) |
| `REMARKS` | teks | VARCHAR2(500) | Remarks — apa adanya (`DDL\Remarks.xml`: Settled / Ex Gratia Payment / Withdraw / Others / --, dipakai frontend) |
| `DATE_OF_LOSS` | teks | VARCHAR2(30) | **baru** — Date of Loss, teks Pega `YYYYMMDDTHHMMSS.mmm GMT` pukul 12:00 WIB (pola Begin date) |
| `LOSS_OBJECT` | teks | VARCHAR2(500) | **baru** — Loss Object |
| `AMOUNT` | angka desimal | NUMBER(38,8) | **baru** — Total of Loss (baca-saja di layar, N-2) — **uang** |
| `PREVENTION_OF_LOSS` | angka desimal | NUMBER(38,8) | **baru** — Prevention Of Loss — **uang** |
| `CAUSE_OF_LOSS` | teks | VARCHAR2(500) | **baru** — Cause of Loss |

## T_COINSDATA

Tiket 42 — migrasi 191, **utuh** (rancangan `.../ListCauseOfLoss/CoinsData`, satu halaman per catatan:
`UQ_T_COINSDATA_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 191) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_COINSDATA` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_LISTCAUSEOFLOSS.ID` |
| `COINS_NAME` | teks | VARCHAR2(500) | Insured Name grid — **selalu** nama tertanggung case (`SetLossRatio_Act` langkah 3.1) |

## MARKETINGOFFICER

Tiket 31 (pilihan Marketing Name). Tabel warisan `POOLDATA`, **baca saja**. Sumber tipe `[terverifikasi]`: DDL
`D:\migrasi\RNM\DDL\MARKETINGOFFICER.txt` (15 kolom; `PEGA_MARKETINGOFFICER.txt` adalah prosedur penulis tabel
ini, bukan tabel). Hanya tiga kolom di bawah yang dibaca.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(100) | nilai pilihan (`.ID`) → disimpan ke `T_QUOTATIONDATA.MOID`; urutan DESC |
| `CLIENTNAME` | VARCHAR2(100) | teks pilihan (pyPrompt `.ClientName`) |
| `MOSTATUS` | VARCHAR2(100) | saringan RD `.MOStatus = "1"` |

## AGENT

Tiket 33 (popup Change SOB). Tabel warisan `POOLDATA.AGENT`, **baca saja**. Sumber tipe `[terverifikasi]`: DDL
`D:\migrasi\RNM\DDL\AGENT.txt` (ditambahkan work owner 03-10-2026; 28 kolom dihitung dua cara; tanpa PK, dua indeks).
⚠️ Kolom tipe-2 dieja **`AGENTTPYE2`** di DDL (salah eja di basis data) — dipetakan dari properti `.AgentType2` `[dugaan]`.
Hanya lima kolom di bawah yang dibaca.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(1000 BYTE) | kode SOB (`.ID`) → `T_QUOTATIONDATA.SOURCE_OF_BUSINESS` (≤ 50, butir 79.2); dicari; urut DESC |
| `CLIENTID` | VARCHAR2(1000 BYTE) | Client ID; syarat `IS NOT NULL`; dicari |
| `CLIENTNAME` | VARCHAR2(1000 BYTE) | nama SOB (`.ClientName`) → `T_QUOTATIONDATA.SOB_NAME` (500); dicari |
| `STATUSACTIVE` | VARCHAR2(1000 BYTE) | syarat `= '1'` (teks) |
| `AGENTTPYE2` | VARCHAR2(1000 BYTE) | syarat `IS NULL OR <> 'LIFE INSURANCE'` (butir 79.1) |

## RISKADDRESS

Tiket 36 (popup Choose Risk Address) dan 37 (Add). Tabel warisan `POOLDATA.RISKADDRESS`, **dibaca dan DISISIPI** (tiket 37:
alamat baru, port Go prosedur `InsertUpdateRISKADDRESS`, butir 81 — ID `GETCURRENTSITE || LPAD(TO_CHAR(RISKADDRESS_SEQ.NEXTVAL),12,'0')`, ≤ 15 karakter
memakai fungsi dan sequence POOLDATA yang ada; tidak dibuat modul ini). Sumber tipe `[terverifikasi]`:
DDL `D:\migrasi\RNM\DDL\RISKADDRESS.txt` (16 kolom, semua VARCHAR2(4000), tanpa PK). Dibaca sembilan kolom (kolom hasil RD
`BrowseRisksAddress_RD`); kesembilan kolom yang sama **ditulis** saat alamat baru disisipkan (tiket 37).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(4000) | nilai Pilih → `.Property.AlmRiskID`; urutan (A118) |
| `TITLE` | VARCHAR2(4000) | Type → `.Property.RoadType` |
| `ADDRESS` | VARCHAR2(4000) | Address → `.Property.RoadName`; saringan A (JALAN) |
| `NATIONNAME` | VARCHAR2(4000) | Country; saringan C |
| `PROVINCENAME` | VARCHAR2(4000) | Province; saringan D |
| `CITYNAME` | VARCHAR2(4000) | City; saringan E |
| `DISTRICTNAME` | VARCHAR2(4000) | District; saringan F |
| `TERRITORYNAME` | VARCHAR2(4000) | Territory; saringan G |
| `POSTALCODE` | VARCHAR2(4000) | Zip Code; saringan B; kunci JOIN ke `RW.ZIPCODE` |

## RW

Tiket 36/37. Tabel warisan `POOLDATA.RW`, **baca saja** — JOIN RD tiket 36 dan saran Zip Code tiket 37; JOIN RD (`ASM-FW-GISFW-Int-RW`, prefix `RW`, INNER JOIN
`.PostalCode = RW.ZipCode`; kelas → tabel `[dugaan]`). Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\RW.txt` (19 kolom;
`RDBMASTERRW.txt` adalah prosedur penulisnya).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ZIPCODE` | VARCHAR2(4000) | syarat JOIN `= RISKADDRESS.POSTALCODE` (tiket 36); saran Zip Code "diawali" (tiket 37) |
| `NOTE` | VARCHAR2(4000) | Territory saran Zip Code (RD `.Note`, tiket 37) |
| `DISTRICTNAME` | VARCHAR2(4000) | District saran Zip Code |
| `CITYNAME` | VARCHAR2(4000) | City saran Zip Code |
| `PROVINCENAME` | VARCHAR2(4000) | Province saran Zip Code |
| `NATION` | VARCHAR2(100) | Country saran Zip Code — RD `.NATIONNAME`; alias `NATION as "NATIONNAME"` `[terverifikasi]` `RDBList\BrowseRW2_SQL.xml` |
| `STS_AKTIF` | VARCHAR2(10) | saringan RD `= "1"` (tiket 37) |

## OCCUPATION

Tiket 38 (saran Occupation Surrounding Risk). Tabel warisan `POOLDATA.OCCUPATION`, **baca saja**. Sumber tipe
`[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\OCCUPATION.txt`; kelas → tabel `[terverifikasi]`
`RDBList\SearchOccupationIDSQL.xml` (kelas `ASM-FW-GISFW-INT-OCCUPATION`, `FROM OCCUPATION`). RD
`BrowseOccupationFacInFIRE_RD`. Hanya empat kolom di bawah yang dibaca.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `OLDID` | VARCHAR2(1000) | nilai saran (`.OldID`) → `T_SURROUNDINGRISK.*_OCCUPATION`; dicari; urutan kedua |
| `NAME` | VARCHAR2(1000) | nama (`.Name`) → `*_NOTE` di layar; dicari; urutan pertama |
| `TYPE` | VARCHAR2(1000) | saringan RD `.Type = Param.TYPE` = `'FIRE'` |
| `KDRISKEXPOSURE` | VARCHAR2(1000) | `kdRiskExposure` → Category okupasi (`SetDataOccupation`, tiket 40) |

## V_JN_OBJ_ITEM

Tiket 39 (pilihan Object Item Type). View warisan `POOLDATA`, **baca saja**. Kelas → tabel `[terverifikasi]`
`RDBList\GetObjectItembyName_SQL.xml` (kelas `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM`, `FROM V_JN_OBJ_ITEM ... AND ISACTIVE ='1'`).
RD `BrowseV_JN_OBJ_ITEM` (DISTINCT, urut JN_OBJ_ITEM, maks 10000) dan `GetObjectItem` (KETERANGAN). ⚠️ DDL view **tidak ada**
di `DDL\` — tipe kolom `belum terverifikasi`; MJOI_KODE dibaca lewat `TO_CHAR`.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `MJOI_KODE` | belum terverifikasi | `kode` → `T_PROPERTYITEMLIST.ITEM_TYPE_ID` |
| `JN_OBJ_ITEM` | belum terverifikasi | `nama` → `ITEM_TYPE`; urutan |
| `KETERANGAN` | belum terverifikasi | `keterangan` → `PROPERTI_ITEM_NOTE` (layar) |
| `ISACTIVE` | belum terverifikasi | saringan RD `.ISACTIVE = 1` (dibandingkan `'1'` seperti SQL Pega) |

## CURRENCY

Tiket 39 (pilihan Currency + pemeriksaan mata uang item). Tabel warisan `POOLDATA`, **baca saja**. Kelas → tabel
`[terverifikasi]` `RDBList\GetAllCurrency.xml` (kelas `ASM-FW-GISFW-INT-CURRENCY`, `select Currency as CURR, OldID as OldID
from currency`). RD `BrowseCurrency_RD` (`.Currency != "ITL"`, maks 500, tanpa urutan). ⚠️ DDL tabel **tidak ada** di
`DDL\` — tipe kolom `belum terverifikasi`.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `CURRENCY` | belum terverifikasi | pilihan Currency (≠ `ITL`), urut CURRENCY (A135); anggota sah `T_PROPERTYITEMLIST.CURRENCY` |

## TABLEOFLIMIT

Tiket 40 (popup Choose Class of Construction). Tabel warisan `POOLDATA.TABLEOFLIMIT`, **baca saja**. Sumber tipe
`[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\TABLEOFLIMIT.txt` (03-10-2026, seluruh kolom VARCHAR2(4000 BYTE)). RD
`BrowseTableOfLimit_RD`. Hanya lima kolom di bawah yang dibaca. `TAHUN` **tidak** disaring (A161, menggantikan A153 —
popup tombol Pega mengirim Tahun kosong).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `BIZCODE` | VARCHAR2(4000 BYTE) | saringan = `BUSINESS.ID` Class of Business case (butir 89); DISTINCT |
| `CATEGORY` | VARCHAR2(4000 BYTE) | saringan `category` (bila diisi); urutan pertama |
| `DESCRIPTION` | VARCHAR2(4000 BYTE) | `description` → Class of Construction; urutan kedua |
| `PCTLIMIT` | VARCHAR2(4000 BYTE) | `pctLimit` teks apa adanya |
| `NOTE` | VARCHAR2(4000 BYTE) | kolom laporan RD (ikut DISTINCT), tidak dikirim |

## COVERAGE_FACIN

Tiket 43 (popup Choose Coverage). **View** warisan `POOLDATA.COVERAGE_FACIN`, **baca saja**. Sumber `[terverifikasi]`: DDL
`D:\migrasi\RNM\DDL\COVERAGE_FACIN.txt` (03-10-2026): `CREATE OR REPLACE FORCE VIEW` lima kolom atas `JSON_VALUE` tabel
`M_COVERAGE` — **tanpa** `ACTIVESTATUS` (A160). Kelas → tabel `[terverifikasi]`: `RDBList\GetCoverageFacin.xml` ("FROM
COVERAGE_FACIN"). RD `BrowseCoverageFacIn_RD` (Type FIRE, NamaCoverage Contains, DISTINCT, maks 500, urut NamaCoverage lalu
OLDID). Tipe kolom view tidak tertulis di DDL (turunan `M_COVERAGE`, DDL-nya tidak ada) — `belum terverifikasi`.

| Kolom | Tipe | Dibaca untuk |
| --- | --- | --- |
| `ID` | belum terverifikasi (dibaca lewat `TO_CHAR`) | `id` → `T_COVERAGELIST.COVERAGE`; ikut DISTINCT |
| `BIZCODE` | belum terverifikasi | kolom laporan RD (ikut DISTINCT), tidak dikirim; saringan RD dikirim kosong |
| `NAMACOVERAGE` | belum terverifikasi | `nama`; saringan `cari` (UPPER LIKE); urutan pertama |
| `TYPE` | belum terverifikasi | saringan = `FIRE` |
| `OLDID` | belum terverifikasi | `oldId`; urutan kedua |

## COVERAGE

Tiket 43 (lima coverage otomatis). Tabel warisan `POOLDATA.COVERAGE`, **baca saja**. Sumber tipe `[terverifikasi]`: DDL
`D:\migrasi\RNM\DDL\COVERAGE.txt`. Kelas → tabel `[terverifikasi]`: `Activity\AddCoverageAutoFire.xml` Obj-Browse
`ASM-FW-GISFW-Int-COVERAGE`, `RDBList\SearchCoverageIDSQL.xml` ("FROM COVERAGE"). Hanya tiga kolom yang dibaca.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(4000 BYTE) | saringan = lima kode `AddCoverageAutoFire` (100815, 100828, 100829, 100825, 100840; urut korpus) |
| `NAME` | VARCHAR2(4000 BYTE) | `nama` → `.CoverageNote` |
| `OLDID` | VARCHAR2(23 BYTE) | `oldId` → `.OLDID` |
