# ReportDefinition siklus New Business — 123 rule

**Ruang lingkup**: `D:\migrasi\RNM\NB FacIn\ReportDefinition\` (123 berkas `.xml`).
**Status area**: belum pernah dibaca pada discovery mana pun sebelumnya.
**Korpus**: READ-ONLY. Tidak ada berkas korpus yang diubah untuk menghasilkan dokumen ini.

Label yang dipakai: `[terverifikasi]` (didukung tag yang dikutip) · `[dugaan]` (dari pola/nama,
belum dikonfirmasi) · `[pertanyaan terbuka]` (tidak dapat dijawab dari korpus) · `[usulan]`
(rekomendasi migrasi, bukan temuan).

---

## 0. Ringkasan temuan

1. **123 ReportDefinition terverifikasi**, tetapi hanya **satu** yang membaca kelas *Work* siklus NB
   (`GetMasterPolicyLife` atas `ASM-FW-GISFW-Work-NB`). Sisanya membaca kelas integrasi
   `ASM-FW-GISFW-Int-*` — tabel referensi/master. **[terverifikasi]**
2. **62 dari 123 report dirujuk dari `Activity`, `DataPage`, atau `RDBList`** — artinya menggerakkan
   logika, bukan sekadar mengisi layar. 60 sisanya hanya dirujuk dari `Section`/`Harness`.
   **[terverifikasi]**
3. **Tidak satu pun ReportDefinition dirujuk dari `Flow` atau `FlowAction`.** Report tidak pernah
   menjadi gerbang transisi status; ia dipanggil dari dalam Activity/DataPage.
   **[terverifikasi]**
4. **Tidak satu pun ReportDefinition menyentuh `HISTORYAKSEPTASIPEGA`, `JSON_POLIS`, `PRODKE`,
   `FACINPRODUCTION`, atau `M_LINK_SERVICE`.** Tangga akseptasi dan produksi dibaca lewat `RDBList`
   (SQL mentah), bukan lewat ReportDefinition. **[terverifikasi]**
5. Mayoritas isi folder bukan Facultative Inward: **28 report** melayani objek pertanggungan lini
   langsung (kendaraan, kapal, kargo, aksesori) dan **3 report** melayani travel/PA ritel — warisan
   framework GIS/SFA. **[terverifikasi]** untuk pengelompokan; **[dugaan]** untuk kesimpulan
   "warisan framework".
6. Hanya **1 report yatim** (`pyDefaultReport`, rule bawaan Pega `@baseclass`).
7. Ditemukan beberapa anomali yang **kandidat perbaikan, bukan bagian migrasi** — lihat §7.

---

## 1. Verifikasi jumlah

```powershell
# Jumlah berkas ReportDefinition NB
(Get-ChildItem "D:\migrasi\RNM\NB FacIn\ReportDefinition" -File -Recurse -Filter *.xml).Count
# -> 123
```

Jumlah berkas per tipe rule di ketiga folder siklus (konteks):

```powershell
$f=@{NB="D:\migrasi\RNM\NB FacIn";RNW="D:\migrasi\RNM\RNW Fac In";EDM="D:\migrasi\RNM\Endorsment Fac In"}
foreach($k in 'NB','RNW','EDM'){ "{0} ReportDefinition = {1}" -f $k,
  (Get-ChildItem "$($f[$k])\ReportDefinition" -File -Filter *.xml).Count }
# -> NB ReportDefinition = 123 / RNW = 118 / EDM = 138
```

Seluruh 123 berkas dapat di-parse sebagai XML tanpa galat (0 `PARSEFAIL`) — verifikasi dilakukan
dengan `[xml]$x = Get-Content <file> -Raw -Encoding UTF8` atas setiap berkas.

**Struktur tag yang dibaca** (contoh lengkap terbaca di
`NB FacIn/ReportDefinition/BrowseCedingCo_RD.xml`):

| Informasi | Tag |
| --- | --- |
| Kelas sumber data | `pyContent/pyClassName` (mis. baris 705) |
| Kolom yang diambil | `pyContent/pyFields/pyListFields/rowdata/pyFieldName` (baris 719-745) |
| Urutan | `pySortType` + `pySortOrder` pada baris kolom yang sama |
| Filter | `pyContent/pyFilters/pyFilter/rowdata` → `pyFilterName`, `pyFilterOperation`, `pyFilterValue` (baris 815-872) |
| Perangkaian filter | `pyContent/pyFilters/pyFilterLogic` (baris 809: `A OR B AND C`) |
| Parameter | `pyContent/pyParameters/rowdata/pyParametersParamName` (baris 747-760) |
| Batas baris | `pyContent/pyMaxRecords` (baris 704) |
| Versi ruleset | `pyRuleSetVersion` (baris 22) |

> **Catatan metode.** Blok `pyUI/...` menduplikasi definisi filter dan kolom untuk keperluan
> tampilan designer. Yang dibaca di dokumen ini adalah blok **`pyContent`** (definisi eksekusi).
> Pada seluruh 123 berkas kedua blok konsisten; bila keduanya berbeda pada rule lain, `pyContent`
> yang dipakai. **[dugaan]** bahwa `pyContent` adalah yang dieksekusi — ini konvensi Pega, bukan
> sesuatu yang dibuktikan oleh korpus.

---

## 2. Kelompok per kelas sumber data

```powershell
# (skrip ekstraksi menulis rd_nb.csv berisi kolom File,Class,Fields,Filters,...)
Import-Csv rd_nb.csv | Group-Object Class | Sort-Object Count -Descending |
  ForEach-Object { "{0,-45} {1}" -f $_.Name,$_.Count }
```

**84 kelas berbeda** untuk 123 report. Kelas dengan lebih dari satu report:

| Kelas | Report | Report yang memakainya |
| --- | ---: | --- |
| `ASM-FW-GISFW-Int-AGENT` | 7 | `BrowseAgent_RD`, `BrowseAgentHierarkiList_RD`, `BrowseAgentNonLife_RD`, `BrowseAgentNusaRe_RD`, `BrowseAgentNusaReLife_RD`, `BrowseCedingCo_RD`, `BrowseSearchSobCeding_RD` |
| `ASM-FW-GISFW-Int-RW` | 6 | `BrowseCity_RD`, `BrowseDistrict_RD`, `BrowseProvince_RD`, `BrowseRW_RD`, `BrowseRWInput_RD`, `BrowseTeritory_RD` |
| `ASM-FW-GISFW-Int-OCCUPATION` | 5 | `BrowseOccupationANEKA_RD`, `BrowseOccupationFacInFIRE_RD`, `BrowseOccupationFIRE_RD`, `BrowseOccupationList_RD`, `BrowseOccupationMBU_RD` |
| `ASM-FW-GISFW-Int-BRANDDETAIL` | 4 | `BrowseBrandDetail_RD`, `BrowseModelList`, `BrowseType_RD`, `BrowseTypeList_RD` |
| `ASM-FW-GISFW-Int-CLAUSE` | 3 | `BrowseClauseAneka_RD`, `BrowseClauseMBU_RD`, `BrowseClausePA_RD` |
| `ASM-FW-SFAGISFW-Work-Opportunity` | 3 | `GetListOpportunity`, `GetListOpportunityF`, `GetListOpportunityLife` |
| `ASM-FW-GISFW-Int-RISKADDRESS` | 3 | `BrowseRiskAddress_RD`, `BrowseRiskAddressZipCode_RD`, `BrowseRisksAddress_RD` |
| `ASM-FW-GISFW-Int-V_DUPLICATE_VEHICLE` | 3 | `CheckMST_KENDARAANByChas`, `CheckMST_KENDARAANByEngine`, `CheckMST_KENDARAANByPlat` |
| `ASM-FW-GISFW-Int-marketingofficer` | 3 | `BrowseMarketingOfficer_RD`, `GetTeamGroup_RD`, `SelectLeader_RD` |
| 12 kelas lain | 2 masing-masing | `CURRENCY`, `SHIP`, `DISTRICT`, `V_JN_OBJ_ITEM`, `COVERAGE`, `Link-Attachment`, `Data-Admin-Operator-ID`, `BRAND`, `DEDUCTIBLE`, `Data-Enumeration`, `ACCUMULATION_LIFE`, `TREATYYEAR_LIFE` |
| 72 kelas sisanya | 1 masing-masing | — |

**Per ruleset** (`pyRuleSet`): `GISFW` 115 · `GISFWInt` 6 · `Pega-EndUserUI` 1 · `Pega-Reporting` 1.
Dua yang terakhir adalah rule produk Pega, bukan kode Nusantara Re. **[terverifikasi]**

### 2.1 Pemetaan kelas → tabel Oracle

Nama tabel **tidak ada** di dalam berkas ReportDefinition. Ia hanya dapat dipulihkan dengan
mencocokkan `pyClassName` terhadap klausa `FROM`/`JOIN` pada `RDBList` sekelas:

```powershell
# untuk setiap RDBList: ambil pyClassName + tabel pada FROM/JOIN di pyBrowseSQL,
# lalu cocokkan dengan daftar kelas yang dipakai ReportDefinition
```

| Kelas | Tabel yang terbukti dari SQL | Bukti |
| --- | --- | --- |
| `Int-PROPORTIONALARRG` | `pooldata.proportionalarrg` | `NB FacIn/RDBList/GetBreakDownSpread_SQL.xml` |
| `Int-ACCUMULATION` | `accumulation`, `czone`, `rw`, `facinproduction`, `pooldata.json_polis` | RDBList sekelas |
| `Int-AGENT` | `agent` | RDBList sekelas |
| `Int-CLIENT` | `pooldata.client` | `NB FacIn/RDBList/GetInsuredID.xml` |
| `Int-NATION` | `m_country` | RDBList sekelas |
| `Int-BRANCH` | `m_branch` | RDBList sekelas |
| `Int-REINSURANCETYPE` | `m_reinsurancetype` **dan** `reinsurancetype` | RDBList sekelas |
| `Int-COVERAGE_FACIN` | `coverage_facin` | RDBList sekelas |
| `Int-OCCUPATION` | `pooldata.occupation` | RDBList sekelas |
| `Int-V_JOB_PA` | `general.v_job_pa` | `NB FacIn/RDBList/SearchJobIDSQL.xml` |
| `Int-M_TREATY_IN`* | `pooldata.M_TREATY_IN` ∪ `pooldata.M_TREATY_IN_edm` | `NB FacIn/RDBList/BrowseTreatyIn.xml` |

\* kelas `Int-TREATY_IN`, tidak dipakai ReportDefinition; dicantumkan sebagai pembanding konvensi.

**Konvensi `ASM-FW-GISFW-Int-<X>` → tabel `<X>` berlaku untuk sebagian besar, tetapi bukan aturan:**
`Int-NATION` → `m_country` dan `Int-BRANCH` → `m_branch` melanggarnya. **[dugaan]** untuk konvensi;
**[pertanyaan terbuka]** untuk **64 kelas** yang tidak punya `RDBList` sekelas di folder NB — nama
tabelnya **belum terverifikasi** dan harus diambil dari definisi kelas Pega (`Rule-Obj-Class`) yang
**tidak ada di korpus**.

---

## 3. Kelompok per tema

```powershell
# pengelompokan berdasarkan regex atas pyClassName; lihat definisi tema di bawah
```

| Tema | Report | di antaranya penggerak alur |
| --- | ---: | ---: |
| A. Geografi & alamat risiko (`RW`,`CITY`,`DISTRICT`,`PROVINCE`,`NATION`,`RISKADDRESS`,`ZONES`,`CZONE`) | 16 | 6 |
| B. Mitra, cedant & organisasi (`AGENT`,`CLIENT`,`marketingofficer`,`lloydagent`,`MV_AGEN`,`BRANCH*`,`V_MARKETING_LEADER`,`Operator`,`V_USER_ASSIGNMENT`) | 19 | 7 |
| C. Treaty, limit, akumulasi & kurs (`PROPORTIONALARRG`,`TABLEOFLIMIT`,`INWARDSCALE`,`ACCUMULATION*`,`TREATYYEAR_LIFE`,`RETROCESSIONLIFE`,`RATE_LIFE_SUMMARY`,`CURRENCY`) | 12 | 7 |
| D. Produk, cover & klausul (`BUSINESS`,`COVERAGE*`,`CLAUSE*`,`DEDUCTIBLE`,`SUBJECTTO`,`WARRANTY`,`BENEFIT*`,`*PLAN`,`*WORDING`,`REINSURANCETYPE`) | 21 | 14 |
| E. Okupasi & pekerjaan (`OCCUPATION*`,`JOB`,`V_JOB_PA`,`MAXAGE`) | 9 | 4 |
| F. **Objek pertanggungan lini langsung** (`BRAND*`,`ACCESORY*`,`COLOR`,`*VEHICLE`,`SHIP`,`CONVEYANCE`,`PACKING`,`GOODSTYPE`,`TRADING`,`MARINE*`,`OBJECT*`,`JN_OBJ`,`KONSTRUKSI`,`OC`,`TYPESHIP`) | 28 | 13 |
| G. Travel / PA / Life ritel (`*TRAVEL`,`PACKAGE_AGE`) | 3 | 3 |
| H. Proses, dokumen & infrastruktur (sisanya) | 15 | 8 |
| **Total** | **123** | **62** |

Tema **F** dan **G** (31 report, 25 % dari folder) menggambarkan asuransi langsung — kendaraan
bermotor, rangka kapal, kargo laut, travel, personal accident. Fac Inward tidak menanggung objek;
ia menanggung porsi risiko cedant. **[dugaan]** bahwa 31 report ini adalah warisan framework
GIS/SFA yang ikut terbawa ke ruleset NB dan tidak dipakai jalur Fac In — namun 16 di antaranya
**tetap dirujuk** dari Activity/Section NB, jadi tidak bisa dibuang tanpa konfirmasi bisnis.

---

## 4. Katalog lengkap 123 ReportDefinition

Kolom **Peran**: `**ALUR** n` = dirujuk dari `Activity`/`DataPage`/`RDBList` (n = total rujukan di
seluruh folder NB) · `UI n` = hanya dirujuk dari `Section`/`Harness` · `— 0` = tanpa rujukan.
Kolom kelas memotong awalan `ASM-FW-GISFW-`. Filter ditulis `properti OPERASI nilai`, dipisah `·`;
daftar kolom dan filter dipotong bila terlalu panjang (isi penuh ada di berkas sumber).

| Report | Kelas | Kolom (`pyListFields`) | `pyFilterLogic` | Filter | Urutan | Peran |
| --- | --- | --- | --- | --- | --- | --- |
| `BrowseAccesorryBrand_RD` | `Int-ACCESORYBRAND` | .ID, .Note | A | .Note StartsWith Param.Note |  | UI 4 |
| `BrowseAccesory_RD` | `Int-ACCESORY` | .ID, .Note | A | .Note StartsWith Param.Note |  | UI 4 |
| `BrowseAccumulatedType_RD` | `Int-ACCUMULATEDTYPE` | .ID, .AccumulationType, .Keyword, .Note, .Type | A AND B | .AccumulationType = Param.AccType · .Note IS NOT NULL  |  | **ALUR** 3 |
| `BrowseAccumulationLife_RD` | `Int-ACCUMULATION_LIFE` | .ID, .ACCUMULATION, .ACCUMULATIONNAME, .NOTE, .KEYWORD, .PERSONNAME, .DOB, .IDCARD, .GENDER, .WEIGHT, .HEIGHT, .AGE, .LEFTHANDED, .OCCUPATION | A AND B | .ID = Param.ID · .ACCUMULATIONNAME = Param.ACCUMULATIONNAME | .ID:DESC:1 | UI 1 |
| `BrowseAgent_RD` | `Int-AGENT` | .ID, .CodeAgentOld, .ClientID, .ClientName, .CountryID, .CountryName, .BU_ID, .Leader0, .ChildCount, .StatusActive, .UpdateUser, .UpdateTime, .Cause, … | A | .ClientID IS NOT NULL  | .ID:DESC:1 | UI 4 |
| `BrowseAgentHierarkiList_RD` | `Int-AGENT` | .ID, .ClientName, .Leader0, .ChildCount, .ClientID | A AND B | .Leader0 = Param.Leader · .StatusActive IS NULL  | .ClientName:ASC:1 | **ALUR** 1 |
| `BrowseAgentNonLife_RD` | `Int-AGENT` | .ID, .ClientID, .ClientName, .CountryID, .CountryName, .BU_ID, .Leader0, .ChildCount, .StatusActive, .UpdateUser, .UpdateTime, .Cause, .OtherCause, .R… | B AND E AND A AND C | .AgentType2 != "LIFE INSURANCE" · .ClientName Contains Param.ClientName · .StatusActive = 1 · .ClientID IS NOT NULL  | .ID:DESC:1 | UI 2 |
| `BrowseAgentNusaRe_RD` | `Int-AGENT` | .ID, .ClientID, .ClientName, .CountryID, .CountryName, .BU_ID, .Leader0, .ChildCount, .StatusActive, .UpdateUser, .UpdateTime, .Cause, .OtherCause, .R… | A AND B AND D AND E AND G | .StatusActive = 1 · .AgentType2 != "LIFE INSURANCE" · .ChildCount = Param.ChildCount · .ClientName Contains Param.ClientName · .ClientID IS NOT NULL  | .ID:DESC:1 | UI 6 |
| `BrowseAgentNusaReLife_RD` | `Int-AGENT` | .ID, .ClientID, .ClientName, .CountryID, .CountryName, .BU_ID, .Leader0, .ChildCount, .StatusActive, .UpdateUser, .UpdateTime, .Cause, .OtherCause, .R… | A AND B AND D AND E AND G | .StatusActive = 1 · .AgentType2 = "LIFE INSURANCE" · .ChildCount = Param.ChildCount · .ClientName Contains Param.ClientName · .ClientID IS NOT NULL  | .ID:DESC:1 | UI 2 |
| `BrowseBankGroup` | `Int-LST_BANK_GROUP` | .LBG_ID, .BANK_GROUP |  |  |  | **ALUR** 1 |
| `BrowseBenefitClause` | `Int-VIEW_BENEFIT_PROPERTY` | .TYPE_PROPERTY, .NOTE, .VALUE, .VALUE_MANDATORY, .PROGRAM_ID, .PLAN_ID | A AND B AND C | .PROGRAM_ID = Param.ProgramId · .PLAN_ID = Param.PlanId · .BENEFIT_ID = Param.ID |  | **ALUR** 1 |
| `BrowseBenefitPkg` | `Int-V_PKG_BENEFIT` | .PCKG_ID, .PROGRAM_ID, .BENEFIT_NAME, .LIMIT, .PLAN_ID, .SORT_ID | A AND B | .PROGRAM_ID = Param.pProgId · .PLAN_ID = Param.pPlanId | .BENEFIT_NAME:ASC:2 / .SORT_ID:ASC:1 | **ALUR** 1 |
| `BrowseBranch_RD` | `Int-BRANCH` | .ID, .Name, .BASTerritory, .CompanyID, .JabodetabekStatus, .Telephone, .BranchParentID, .BranchName, .Status, .BranchStatus, .Kanwil, .BranchType, .HS… | A | .ID = Param.ID | .ID:ASC:1 | UI 1 |
| `BrowseBranchDetail_RD` | `Int-BRANCHDETAIL` | .ID, .BranchStatus, .BranchType, .Name, .CompanyID, .Status, .Kanwil, .KanwilGroup | A AND B AND C AND D | .ID = Param.ID · .Status = 1 · .BranchStatus = Param.BranchStatus · .BranchParentID = Param.ParentID | .ID:ASC:1 | UI 2 |
| `BrowseBrand_RD` | `Int-BRAND` | .ID, .Name, .Type | A AND B | .Name Contains Param.Name · .Type = "MBU" |  | UI 7 |
| `BrowseBrandDetail_RD` | `Int-BRANDDETAIL` | .ID, .BrandID, .Name, .CarTypeID, .BuiltYear, .Automatic, .Silinder, .BusinessID, .VehicleCodeID, .IsActive, .InputDate, .RPTVehicleCodeID, .GroupType… | B AND C AND D | .Type = "MBU" · .Name StartsWith Param.Name · .BrandID = Param.BrandID |  | UI 4 |
| `BrowseBusiness_RD` | `Int-BUSINESS` | .NoteINA, .Note, .MaxDisc, .StatusRandomTF, .ContentNote, .PolicyCost, .MinPremi, .GroupPanel, .ID, .OLDID, .BusinessGroupName, .EndorseCost, .Deducta… | (A OR B) AND C | .ID = Param.ID · .Note Contains Param.Note · .BusinessGroupID = Param.Group | .ID:DESC:1 | UI 1 |
| `BrowseCedingCo_RD` | `Int-AGENT` | .ID, .ClientName, .Leader0, .ChildCount | A OR B AND C | .Leader0 = Param.CedingCoLeader · .ID = Param.ID · .StatusActive IS NULL  | .ClientName:ASC:1 | **ALUR** 1 |
| `BrowseCity_RD` | `Int-RW` | .CITYNAME | A AND B | .PROVINCENAME = Param.Province · .NATIONNAME = Param.Nation |  | UI 4 |
| `BrowseCityInput_RD` | `Int-CITY` | .ID, .Note | A AND B AND C AND D | .ZipCode = Param.ZipCode · .ID = Param.CITYID · .Note = Param.Note · .ProvinceID = Param.PROVINCEID | .Note:ASC:1 | **ALUR** 5 |
| `BrowseClauseAneka_RD` | `Int-CLAUSE` | .ID, .Info, .SectionCode, .BusinessCode, .StatusTransCenter, .BrokerCode, .ClauseData, .StatusActive, .Type | B AND A | .Type = "ANEKA" · .BusinessCode = Param.BusinessCode |  | **ALUR** 1 |
| `BrowseClauseMBU_RD` | `Int-CLAUSE` | .ID, .Info, .CoverageID, .CoverageName, .Count, .Queue, .NoteENG, .ClauseData, .ClauseDataENG, .StatusActive, .StatusTransWorkshop, .Type, .SumOfArgum… | A | .Type = "MBU" |  | **ALUR** 1 |
| `BrowseClausePA_RD` | `Int-CLAUSE` | .ID, .Info, .ApplicantType, .ActiveStatus, .UserID, .InputDate, .Type | A | .Type = "PA" |  | **ALUR** 1 |
| `BrowseClausePlan` | `Int-M_KLAUSUL_PLAN` | .KLAUSUL_ID, .KLAUSUL_NOTE, .GROUP_PROG_ID | A AND B OR A AND C AND D | .STS_AKTIF = 1 · .GROUP_PROG_ID = Param.GroupProgId · .GROUP_PROG_ID = 0 · .CLIENT_TYPE = pyWorkPage.Quotation.BusinessCode |  | **ALUR** 1 |
| `BrowseClieent_RD` | `Int-CLIENT` | .ID, .Name, Agent.AgentType2, Agent.ID, .OldID, .MCLID, .MCL_TGL_Lahir, .BU_ID, .BU_Note, .IDVIEW, .IDNUMBER | A AND B AND C AND D AND E AND F AND (G OR H) | .ID Contains Param.ID · .Name Contains Param.Name · .Name IS NOT NULL  · .BU_Note Contains Param.BU_Note · .IDNUMBER = Param.IDNUMBER · Agent.AgentType2 = Param.AgentType · Agent.StatusActive != "0" · Agent.Sta… | .ID:DESC:1 / .Name:ASC:2 | UI 3 |
| `BrowseColor_RD` | `Int-COLOR` | .ID, .Note | A | .Note StartsWith Param.Note |  | **ALUR** 5 |
| `BrowseConveyance_RD` | `Int-CONVEYANCE` | .ID, .Description, .TypeConveyance, .Type | A AND B | .Description = Param.Description · .ID = Param.ID |  | **ALUR** 5 |
| `BrowseConveyanceType` | `Int-MCONVEYANCE` | .CONVEYANCE_ID, .CONVEYANCE_DESC | A AND B | .CONVEYANCE_DESC = Param.CONVEYANCE_DESC · .CONVEYANCE_ID = Param.CONVEYANCE_ID |  | UI 2 |
| `BrowseCoverage_RD` | `Int-COVERAGE` | .ID, .Name, .BusinessCode, .StatusMainCover, .ActiveStatus, .Type, .DiscountID, .Status, .ACCTypeCode, .Description, .Name2, .GroupCode, .Rate, .FlagR… | (A OR F) AND B AND C AND D AND E | .BusinessCode = Param.BusinessCode · .StatusMainCover <= Param.StatusMainCover · .Name Contains Param.Name · .Name NotContain "MACHINERY BREAK DOWN + 4.1ACC" · .ID = Param.ID · .BusinessOldID = Param.BusinessCo… |  | UI 1 |
| `BrowseCoverageAneka` | `Int-V_COVERAGE_ANEKA` | .KD_COVER, .NAMA_COVER | A AND B AND C AND D AND E | .KD_BISNIS = Param.kdbisnis · .STS_COVER_UTAMA <= Param.coverutama · .NAMA_COVER Contains Param.note · .NAMA_COVER NotContain "MACHINERY BREAK DOWN + 4.1ACC" · .KD_COVER = Param.kdCoverage |  | **ALUR** 1 |
| `BrowseCoverageFacIn_RD` | `Int-COVERAGE_FACIN` | .ID, .BizCode, .NamaCoverage, .Type, .OLDID, .ACTIVESTATUS | A AND B AND C AND (D OR E) | .Type = Param.Type · .BizCode = Param.BizCode · .NamaCoverage Contains Param.Nama · .ACTIVESTATUS = 1 · .ACTIVESTATUS IS NULL  | .NamaCoverage:ASC:1 / .OLDID:ASC:2 | **ALUR** 6 |
| `BrowseCoverageMBU` | `Int-COVERAGE` | .ID, .RateMBU, .OrderedNum, .MBUDeductible, .TransWorkShopStatus, .HeavyValueSharia, .HeavyValue, .Name | B AND C | .Type = "MBU" · .Name Contains Param.Name |  | UI 3 |
| `BrowseCurrency_RD` | `Int-CURRENCY` | .ID, .CountryID, .Currency, .CurrencySymbol, .CountryName, .ISOSymbol, .Note, .OLDID | A AND B AND C | .Currency = Param.Currency · .ID = Param.ID · .Currency != "ITL" |  | **ALUR** 26 |
| `BrowseCurrencyTreatyIn_RD` | `Int-CURRENCY` | .ID, .CountryID, .Currency, .CurrencySymbol, .CountryName, .ISOSymbol, .Note, .OLDID | A AND B | .Currency = Param.Currency · .Currency != "ITL" |  | UI 4 |
| `BrowseCZoneIsNotNull_RD` | `Int-CZONE` | .Description, .ID, .GroupOf, .Code, .GroupOfName | A AND B AND C | .GroupOf IS NOT NULL  · .GroupOfName Contains Param.GroupOfName · .Code = Param.Code | .Description:ASC:1 | UI 2 |
| `BrowseDeductibleFacIn_RD` | `Int-DEDUCTIBLE` | .OccupationID, .OccupationsName, .MinTSI, .TSIKursID, .Min, .Type, .MinPCT, .Description, .MaxKursID, .Score, .MaxTSI, .StatusActive, .MinDeductOF, .I… | A AND B AND C AND D | .Type = Param.Type · .BusinessCode = Param.BizCode · .CoverageCode = Param.CoverageCode · .DeductName Contains Param.DeductName | .Description:ASC:1 / .DeductName:ASC:2 | UI 2 |
| `BrowseDeductMBU_RD` | `Int-DEDUCTIBLE` | .ID, .Template, .Note, .Type, .Language | A | .Type = "MBU" |  | UI 3 |
| `BrowseDistrict_RD` | `Int-RW` | .DistrictName | A AND B AND C | .PROVINCENAME = Param.Province · .CITYNAME = Param.City · .NATIONNAME = Param.Nation |  | UI 3 |
| `BrowseDistrictInput_RD` | `Int-DISTRICT` | .ID, .DistrictName | A AND B | .ZipCode = Param.ZipCode · .CityID = Param.CityCode |  | **ALUR** 4 |
| `BrowseDistrictInputC_RD` | `Int-DISTRICT` | .ID, .CityID, .DistrictName, .CityName | A AND B | .CityID = Param.ID · .CityName = Param.CityName |  | UI 1 |
| `BrowseEnumData` | `Data-Enumeration` | .Type, .Value, .Value2, .Name | A AND B AND C | .Type = Param.TIPE · .Value = Param.VALUE · .Name2 = Param.NAME2 |  | **ALUR** 1 |
| `BrowseGoodsType` | `Int-GOODSTYPE` | .ID, .Description, .Type | A AND B | .Description = Param.Description · .ID = Param.ID |  | **ALUR** 5 |
| `BrowseJob_RD` | `Int-JOB` | .ID, .Describe, .ClassID, .ClassName, .EngDescribe, .Type | A AND B AND C | .ID = Param.ID · .Describe = Param.DESCRIBE · .Type = Param.TIPE |  | **ALUR** 9 |
| `BrowseJobPA` | `Int-V_JOB_PA` | .JOB_ID, .DESCRIBE, .KELAS_ID, .GOLONGAN |  |  |  | UI 2 |
| `BrowseJTypePropertyPlan` | `Int-VJ_M_TYPE_PROPERTY_PLAN` | .ID_PROPERTY, .ID_PROPERTY_OLD, .GROUP_PROG_ID, .USER_IU, .TGL_IU, .DESCRIPTION, .STS_AKTIF |  |  |  | **ALUR** 1 |
| `BrowseLimit` | `Int-PROPORTIONALARRG` | .ParentReinsTypeID, .ReinsTypeID, .Rp, .Usd, .TreatyYear, .TreatyGroupID, .TreatyDescID, .ID, .ReinsTypeName, .Line | (A OR G) AND B AND C AND D AND E AND F AND H | .Rp > Param.Rp · .Rp <= Param.Rupiah · .TreatyYear = Param.TreatyYear · .TreatyGroupID = Param.TreatyGroupID · .ID = Param.ID · .Usd = .Usd · .Rp IS NOT NULL  · .Rp > Param.RPFacout |  | UI 6 |
| `BrowseLLOYDAGENT_RD` | `Int-lloydagent` | .ID, .Contact, .Type, .Name, .Address, .City, .Country, .TypeAgentID | A AND B | .ID = Param.ID · .Name Contains Param.Name |  | **ALUR** 7 |
| `BrowseM_KONSTRUKSIList` | `Int-M_KONSTRUKSI` | .DESCRIPTION, .LANDSL_RATE, .FLOOD_RATE, .BISNIS_ID, .KONSTRUKSI, .TGL_AKTIF | A AND B | .BISNIS_ID = Param.BISNISID · .TGL_AKTIF = Param.TGLAKTIF |  | UI 5 |
| `BrowseM_OKUPASI_ANEKA` | `Int-M_OKUPASI_ANEKA` | .KD_OKUPASI, .NAMA_OKUPASI | A AND B AND C | .KD_BISNIS = Param.kdbisnis · .KD_OKUPASI = Param.kdokupasi · .NAMA_OKUPASI = Param.note |  | **ALUR** 9 |
| `BrowseMarineCondition_RD` | `Int-MARINECONDITION` | .ID, .Description, .ConditionAdd, .Type |  |  |  | UI 7 |
| `BrowseMarketingOfficer_RD` | `Int-marketingofficer` | .ID, .ClientID, .BranchDetailID, .BranchDetailName, .MOLeader, .MOStatus, .ClientName, .ClientID2, .BranchStatus, .TeamGroup, .Tanggal, .UserUpdate | A | .MOStatus = 1 | .ID:DESC:1 | UI 4 |
| `BrowseMaxAge_RD` | `Int-MAXAGE` | .ID, .Warranty, .Type |  |  |  | UI 4 |
| `BrowseMaxHariTravel_RD` | `Int-PREMITRAVEL` |  | A | .ID = Param.PLAN_ID |  | **ALUR** 1 |
| `BrowseMerk_RD` | `Int-BRAND` | .ID, .BrandName, .ActiveStatus, .Type | A AND B AND C | .Type = "ANEKA" · .ID = Param.ID · .BrandName Contains Param.BrandName |  | UI 3 |
| `BrowseModelList` | `Int-BRANDDETAIL` | .CarCategoryID, .CarCategoryName | F1 | .CarCategoryName Contains Param.CarCategoryName | .CarCategoryName:ASC:1 | **ALUR** 1 |
| `BrowseNation_RD` | `Int-NATION` | .ID, .Note, .NationInitial | A OR B | .ID = Param.ID · .Note = Param.Note |  | **ALUR** 5 |
| `BrowseObject_RD` | `Int-OBJECT` | .ID, .ObjectName, .BusinessCode, .ActiveStatus, .BusinessName | A AND B | .ID = Param.ID · .ObjectName Contains Param.Note |  | UI 3 |
| `BrowseObjectItemType_RD` | `Int-OBJECTITEMTYPE` | .ID, .Note, .PCTAdjustable2, .Type, .ObjectItemTypeINA, .PCTAdjustable1, .ObjectItemType, .Group |  |  |  | UI 1 |
| `BrowseOC_RD` | `Int-OC` | .ID, .Name, .IdentityType, .Identity, .Address, .City, .PostalCode, .Country, .PPHTypeID, .BankID, .NoAccount, .Institution, .CentralTransStatus, .Cli… |  |  |  | UI 2 |
| `BrowseOccupationANEKA_RD` | `Int-OCCUPATION` | .ID, .Name, .BusinessCode, .Type, .BusinessName, .OldID | A AND B AND C AND D | .ID = Param.ID · .Type = Param.TYPE · .BusinessName = Param.BUSINESSNAME · .Name Contains Param.NAME | .ID:DESC:1 | **ALUR** 5 |
| `BrowseOccupationFacInFIRE_RD` | `Int-OCCUPATION` | .IsNegativeList, .ParentID, .HazardLevel, .Notes, .KDRiskExposure, .Name, .TerrorismGroup, .GroupOccupation, .ID, .TimeExcessDays, .PCTLimitSyariah, .… | A AND B AND D AND E AND (F OR C) | .ID = Param.ID · .ParentID = Param.PARENTID · .Name Contains Param.DESCRIPTION · .Type = Param.TYPE · .KDRiskExposure = Param.KDRISKEXPOSURE · .OldID Contains Param.OLDID | .Name:ASC:1 / .OldID:ASC:2 | **ALUR** 9 |
| `BrowseOccupationFIRE_RD` | `Int-OCCUPATION` | .IsNegativeList, .OccupationCode, .ParentID, .HazardLevel, .Notes, .KDRiskExposure, .Name, .TerrorismGroup, .GroupOccupation, .ID, .TimeExcessDays, .P… | ((A AND B AND C AND D AND F) AND E) OR ((A OR C OR F) AND D) | .ID = Param.ID · .ParentID = Param.PARENTID · .Name Contains Param.DESCRIPTION · .Type = "FIRE" · .OccupationCode NotStartsWith Param.OccupationCode · .OccupationCode = Param.OccupationCodelama | .ID:DESC:1 | UI 3 |
| `BrowseOccupationList_RD` | `Int-OCCUPATION` | .ID, .Notes | A AND B AND C AND D | .Type = "MBU" · .STSActive = "1" · .ID = Param.Code · .Notes = Param.OccNotes |  | **ALUR** 1 |
| `BrowseOccupationMBU_RD` | `Int-OCCUPATION` | .ID, .Notes, .CarClassID, .CASCOJHP, .CASCOMINPREMI, .TLOJHP, .TLOMINPREMI, .MINOWNRISK, .MAXOWNRISK, .PassangerDestination, .DateSetting, .Function, … | A AND B | .Type = "MBU" · .STSActive = "1" |  | UI 4 |
| `BrowseOccupations_RD` | `Int-OCCUPATIONS` | .IsNegativeList, .ParentID, .HazardLevel, .Notes, .KDRiskExposure, .Description, .TerrorismGroup, .GroupOccupation, .ID, .TimeExcessDays, .PCtLimitSya… | A AND B AND C | .ID = Param.ID · .ParentID = Param.PARENTID · .Description Contains Param.DESCRIPTION |  | UI 1 |
| `BrowsePacking_RD` | `Int-PACKING` | .ID, .Description, .Type | A AND B | .Description = Param.Description · .ID = Param.ID |  | **ALUR** 5 |
| `BrowsePlanByAge` | `Int-V_PACKAGE_AGE_KLAUSUL` | .PROGRAM_ID, .PLAN_ID, .PREMIUM, .GENDER, .KLAUSUL_VALUE_FROM, .KLAUSUL_VALUE, .PLAN_NAME, .YEAR_LIMIT, .LIMIT_KEJADIAN, .GROUP_PROG_ID, .PROGRAM_NAME | A AND B AND C AND D | .KLAUSUL_VALUE >= Param.dAge · .KLAUSUL_VALUE_FROM < Param.dAge · .PCKG_ID = param.dPckgId · .GENDER = param.sSex |  | **ALUR** 2 |
| `BrowsePlanProperty` | `Int-M_TYPE_PROPERTY_PLAN` | .ID_PROPERTY, .DESCRIPTION, .GROUP_PROG_ID | A AND B OR A AND C | .STS_AKTIF = 1 · .GROUP_PROG_ID = Param.GroupProgId · .GROUP_PROG_ID = 0 |  | **ALUR** 3 |
| `BrowsePremiTravel` | `Int-VIEW_PREMI_TRAVEL` | .PREMI_LEBIH_PERMINGGU, .SAMPAI_HARI, .DARI_HARI, .PLAN_ID, .PREMI_LEBIH_PERMINGGU, .TSI, .PREMI, .PLAN_JENIS | A AND B AND C | .PLAN_ID = Param.PLAN_ID · .SAMPAI_HARI >= Param.HARI · .DARI_HARI <= Param.HARI |  | **ALUR** 1 |
| `BrowseProvince_RD` | `Int-RW` | .PROVINCENAME | F1 AND A | .NATIONNAME = Param.Nation · .PROVINCENAME Contains Param.PROVINCENAME |  | **ALUR** 8 |
| `BrowseProvince2_RD` | `Int-PROVINCE` | .ID, .NationID, .Note, .NationName | A AND B AND C AND D | .ID = param.ID · .NationID = param.NationID · .NationName = Param.NationName · .Note = Param.Note |  | UI 1 |
| `BrowserAccessoryType_RD` | `Int-ACCESORYTYPE` | .ID, .Note | A | .Note StartsWith Param.Note |  | UI 4 |
| `BrowseRateLifeSummary` | `Int-RATE_LIFE_SUMMARY` | .ID, .USEDBY, .OPERATORID, .MODIFIEDDATE, .TYPE | A AND B | .ID = param.id · .USEDBY Contains param.idusedby | .ID:ASC:1 | UI 2 |
| `BrowseReinsuranceType_RD` | `Int-REINSURANCETYPE` | .ID, .Note, .Type, .SOANote, .Code, .UserID, .Flag | (A OR B) AND C AND D | .ID = Param.ID · .Note = Param.Note · .Flag = Param.Flag · .Type = Param.Type | .ID:DESC:1 | **ALUR** 19 |
| `BrowseRetrocessionLife_RD` | `Int-RETROCESSIONLIFE` | .ID, .IDTREATYYEAR_LIFE, .PERCENTSHARE, .RATE, .COMMISION, .TREATYTYPEID, .TREATYTYPENAME, .REINSURERNAME, .TREATYSTARTDATE, .TREATYENDDATE, .USERID, … | F1 AND A | .IDTREATYYEAR_LIFE = Param.IDTreatyYear_Life · .ID = Param.IDRetro | .TGLUPDATE:DESC:1 | **ALUR** 3 |
| `BrowseRiskAddress_RD` | `Int-RISKADDRESS` | .ID, .IDCity, .CityName, .Address, .IDDistrict, .DistrictName, .IDNation, .NationName, .IDProvince, .ProvinceName, .IDTerritory, .TerritoryName, .Post… | A AND B AND C AND D AND E AND F AND G AND H | .Address Contains Param.JALAN · .PostalCode Contains Param.KODEPOS · .NationName Contains Param.COUNTRY · .ProvinceName Contains Param.PROVINCE · .CityName Contains Param.CITY · .DistrictName Contains Param.DIS… |  | **ALUR** 2 |
| `BrowseRiskAddressZipCode_RD` | `Int-RISKADDRESS` | .ID, .IDCity, .CityName, .Address, .IDDistrict, .DistrictName, .IDNation, .NationName, .IDProvince, .ProvinceName, .IDTerritory, .TerritoryName, .Post… | A | .ProvinceName Contains Param.ProvinceName |  | UI 2 |
| `BrowseRisksAddress_RD` | `Int-RISKADDRESS` | .ID, .CityName, .Address, .DistrictName, .NationName, .ProvinceName, .TerritoryName, .PostalCode, .Title | A AND B AND C AND D AND E AND F AND G | .Address Contains Param.JALAN · .PostalCode Contains Param.KODEPOS · .NationName Contains Param.COUNTRY · .ProvinceName Contains Param.PROVINCE · .CityName Contains Param.CITY · .DistrictName Contains Param.DIS… |  | UI 1 |
| `BrowseRW_RD` | `Int-RW` | .ZipCode, .CZONE, .Note, .DistrictName, .CITYNAME, .PROVINCENAME, .NATIONNAME | A AND B AND C AND D AND E AND F | .ZipCode = Param.ZipCode · .PROVINCENAME = Param.Province · .CITYNAME = Param.City · .DistrictName = Param.District · .Note = Param.Teritory · .STS_AKTIF = "1" |  | UI 3 |
| `BrowseRWInput_RD` | `Int-RW` | .ID, .Note, .ZipCode, .DistrictName, .DistrictID | A AND B AND C | .ZipCode = Param.ZipCode · .DistrictID = Param.DistrictCode · .CityID = Param.CityCode |  | **ALUR** 2 |
| `BrowseSameAccumulationLife_RD` | `Int-ACCUMULATION_LIFE` | .ID, .ACCUMULATION, .ACCUMULATIONNAME, .NOTE, .KEYWORD, .PERSONNAME, .DOB, .IDCARD, .GENDER, .WEIGHT, .HEIGHT, .AGE, .LEFTHANDED, .OCCUPATION | A AND B | .DOB = Param.DOB · .PERSONNAME = Param.PersonName | .ID:DESC:1 | **ALUR** 1 |
| `BrowseSearchSobCeding_RD` | `Int-AGENT` | .ID, .ClientName, .Leader0, .ChildCount | B AND A | .ID = Param.ID · .StatusActive = 1 | .ClientName:ASC:1 | **ALUR** 1 |
| `BrowseShip_RD` | `Int-SHIP` | .ID, .Name, .ShipTypeID, .ShipName, .FormerShipName, .Classification, .FlagName, .Construction, .GRT, .NRT, .DWT, .YearMake1, .YearMake2, .Rebuilt, .R… | A AND B AND C AND D AND E AND F AND (G OR H) | .Register = Param.Register · .Name Contains Param.Name · .ID = Param.ID · .FlagName Contains Param.FlagName · .Construction Contains Param.Konstruksi · .Remark Contains Param.Remark · .Remark != "Not Active" · … | .ID:DESC:1 | **ALUR** 9 |
| `BrowseShip2_RD` | `Int-SHIP` | .ID, .Name, .ShipTypeID, .ShipName, .FormerShipName, .Classification, .FlagName, .Construction, .GRT, .NRT, .DWT, .YearMake1, .YearMake2, .Rebuilt, .R… | A AND B AND C AND D AND E AND F AND (G OR H) | .Register = Param.Register · .Name = Param.Name · .ID = Param.ID · .FlagName Contains Param.FlagName · .Construction Contains Param.Konstruksi · .Remark Contains Param.Remark · .Remark != "Not Active" · .Remark… | .ID:DESC:1 | **ALUR** 1 |
| `BrowseSourceBizList` | `Int-MV_AGEN` | .NAMA, .LAG_AGEN_ID | A OR B OR C | .NAMA Contains Param.Nama · .LAG_AGEN_ID = Param.Kode · .LAG_LEADER = Param.Leader |  | UI 2 |
| `BrowseSubContract_RD` | `Int-SUBCONTRACT` | .ID, .Name |  |  |  | **ALUR** 1 |
| `BrowseSubjectTo_RD` | `Int-SUBJECTTO` | .ID, .BusinessID, .BusinessName, .OccupationID, .OccupationName, .Question, .PrintPolice, .YN, .StatusActive, .UserID, .InputData, .LossExperience, .T… | A AND B AND C | .BusinessID = Param.BusinessID · .OccupationID = Param.OccupationID · .StatusActive = "1" |  | **ALUR** 1 |
| `BrowseTableOfLimit_RD` | `Int-TABLEOFLIMIT` | .Bizcode, .Category, .Description, .PctLimit, .Note | A AND B AND C AND D | .Tahun = Param.Tahun · .Bizcode = Param.Bizcode · .Category = Param.Category · .ID = Param.ID | .Category:ASC:1 / .Description:ASC:2 | **ALUR** 7 |
| `BrowseTeritory_RD` | `Int-RW` | .Note, .ZipCode | A AND B AND C AND D | .PROVINCENAME = Param.Province · .CITYNAME = Param.City · .DistrictName = Param.District · .NATIONNAME = Param.Nation |  | UI 2 |
| `BrowseTrading_RD` | `Int-TRADING` | .ID, .Description, .Type | A AND B | .ID = Param.ID · .Description = Param.Description |  | **ALUR** 5 |
| `BrowseTreatyYear_Life_RD` | `Int-TREATYYEAR_LIFE` | .ID, .USERID, .TGLUPDATE, .STARTDATE, .ENDDATE, .TREATYYEAR, .UNDERWRITINGYEAR | A | .ID = Param.ID | .ID:ASC:1 | **ALUR** 7 |
| `BrowseTreatyYearOR_Life_RD` | `Int-TREATYYEAR_LIFE` | .ID, .TREATYTYPEID, .USERID, .TGLUPDATE, .TREATYYEAR_LIFE, .TREATYTYPENAME, .IDR, .USD, .B_IDR, .B_USD | A | .TREATYTYPENAME = Param.TName | .ID:ASC:1 | **ALUR** 1 |
| `BrowseType_RD` | `Int-BRANDDETAIL` | .ActiveStatus, .Type, .ID, .TypeName, .ObjectID, .ObjectName | A AND B AND C | .ID = Param.ID · .Type = "ANEKA" · .TypeName = Param.TypeName | .ID:DESC:1 | UI 3 |
| `BrowseTypeList_RD` | `Int-BRANDDETAIL` | .BrandID, .Name, .CarTypeName, .CarTypeID, .ID, .VehicleCodeID, .VehicleName, .BrandName | A AND B AND C AND D | .Name = Param.Note · .CarTypeName = Param.ModelName · .BrandName = Param.BrandName · .Type = "MBU" | .Name:ASC:1 | **ALUR** 1 |
| `BrowseTypeShip_RD` | `Int-TYPESHIP` | .ID, .Name, .Type |  |  |  | UI 4 |
| `BrowseUserAssignment` | `Int-V_USER_ASSIGNMENT` | .LUS_ID, .USER_ID, .MCL_NAME, .JABATAN |  |  |  | UI 4 |
| `BrowseV_ALM_RISIKOList` | `Int-V_ALM_RISIKO` | .TITLE, .PROPINSI_ID, .KOTA_NM, .KABUPATEN_ID, .NM_JALAN, .WILAYAH_NM, .NEGARA_NM, .PROPINSI_NM, .ALM_RISK_ID, .KOTA_ID, .NEGARA_ID, .WILAYAH_ID, .KD_… | A AND B AND C | .NM_JALAN Contains Param.JALAN · .KD_POS = Param.KODEPOS · .ALM_RISK_ID = Param.RISKID |  | UI 1 |
| `BrowseV_JN_OBJ_ITEM` | `Int-V_JN_OBJ_ITEM` | .MJOI_KODE, .JN_OBJ_ITEM, .KETERANGAN, .KELOMPOK | A AND B AND C | .JN_OBJ_ITEM = Param.JN_OBJ_ITEM · .ISACTIVE = 1 · .KELOMPOK = Param.KELOMPOK | .JN_OBJ_ITEM:ASC:1 | UI 2 |
| `BrowseVClauseArg` | `Int-V_CLAUSES_ARG` | .LKL_ID, .NO_ARG, .DESCRIPTION, .DEFAULT_VAL | A AND B | .LKL_AKTIF = "1" · .LKL_ID = Param.CLAUSE |  | **ALUR** 1 |
| `BrowseVDraftWordingList` | `Int-V_DRAFTWORDING` | .ID, .BIZCODE, .NAMA_BISNIS, .JENIS | A | .BIZCODE = Param.bizcode |  | UI 3 |
| `BrowseVMarketingLeaderList` | `Int-V_MARKETING_LEADER` | .LMO_ID, .LMO_LEADER, .LEADER_NAME, .LDC_ID | A AND B | .LMO_ID = Param.pLMO_ID · .LDC_ID = Param.pLDC_ID |  | **ALUR** 1 |
| `BrowseWarranty_RD` | `Int-WARRANTY` | .ID, .Description, .Type | A AND B | .ID = Param.ID · .Description = Param.DESCRIPTION |  | UI 1 |
| `BrowseZones_RD` | `Int-ZONES` | .OLDID, .Type, .Zone, .Description, .ID, .Description2, .ActiveDate, .Rate | A AND B AND (C OR D OR E OR F OR G OR H) AND I | .ID = Param.ID · .ActiveDate = Param.ACTIVEDATE · .Zone = Zone.CARI1 · .Zone = Zone.CARI2 · .Zone = Zone.CARI3 · .Zone = Zone.CARI4 · .Zone = Zone.CARI5 · .Zone = Zone.CARI6 · .OLDID = Param.COVERAGE |  | UI 2 |
| `CheckMST_KENDARAANByChas` | `Int-V_DUPLICATE_VEHICLE` |  | B | .MKE_NO_RANGKA = @@pxUpperCase(@@pxLeftTrim(@@pxRightTrim(@@pxUpperCase(Param.MKE_NO_RANGKA)))) |  | **ALUR** 1 |
| `CheckMST_KENDARAANByEngine` | `Int-V_DUPLICATE_VEHICLE` |  | C | .MKE_NO_MESIN = @@pxUpperCase(@@pxLeftTrim(@@pxRightTrim(@@pxUpperCase(Param.MKE_NO_MESIN)))) |  | **ALUR** 1 |
| `CheckMST_KENDARAANByPlat` | `Int-V_DUPLICATE_VEHICLE` |  | A | .MKE_NO_PLAT = @@pxUpperCase(@@pxLeftTrim(@@pxRightTrim(@@pxUpperCase(Param.MKE_NO_PLAT)))) |  | **ALUR** 1 |
| `DataTableEditorReport` | `Data-Enumeration` | .Type, .Name, .Language, .Status, .Value, .Description, .Name2, .Value2, .Name3, .Name4 | A AND B AND C AND D AND E | .Type = Param.Type · .Language = Param.Language · .Status = Param.Status · .Name2 = Param.Name2 · .Name3 = Param.Name3 | .Type:ASC:1 / .Language:ASC:2 / .Value:ASC:3 | **ALUR** 1 |
| `GetAllDocument` | `Int-DOCUMENT_POLIS` | .ID, .T_STORAGE_ID, .INSKEY_LINK, .MIME, .NAMAFILE, .KATEGORI_2, .TANGGAL | F1 AND F2 | .KATEGORI_2 = Param.Kategori · .IDPEGA = Param.IDPega |  | UI 2 |
| `GetListOpportunity` | `ASM-FW-SFAGISFW-Work-Opportunity` | .pyID, .Name, .pxCreateDateTime, .pxCreateOpName, .pzInsKey, .NBHandle, .TextNoQuotation, .SellingMode, A.NBStatus, A.NBStatusNew, A.Quotation.Marketi… | B AND D AND E AND F AND A AND (C OR G) AND F1 | .TextNoQuotation Contains "NB-" · A.pyStatusWork != "Resolved-Completed" · A.Quotation.BusinessFac = T · A.pxCreateOperator = Param.UserIdentifier · A.pyStatusWork != "Resolved-Rejected" · .Name Contains Param.… | .pxCreateDateTime:DESC:1 | UI 1 |
| `GetListOpportunityF` | `ASM-FW-SFAGISFW-Work-Opportunity` | .pyID, .Name, A.Quotation.InsuredName, .pxCreateDateTime, .pxCreateOpName, .pxCreateOperator, .pzInsKey, .NBHandle, .TextNoQuotation, .SellingMode, A.… | A AND B AND C AND E AND D AND F AND (G OR H) AND F1 | A.pyStatusWork != "Resolved-Completed" · A.pyStatusWork != "Resolved-Rejected" · A.Quotation.BusinessFac = F · A.Quotation.TeamGroup = Param.TeamGroup · .TextNoQuotation Contains "NB-" · A.pxCreateOperator = Pa… | .pxCreateDateTime:DESC:1 | UI 1 |
| `GetListOpportunityLife` | `ASM-FW-SFAGISFW-Work-Opportunity` | .pyID, .Name, A.Quotation.InsuredName, A.Quotation.BusinessName, .pxCreateDateTime, .pxCreateOpName, .pzInsKey, .TextNoQuotation, A.NBStatus, A.NBStat… | B AND A AND D AND G AND C AND E AND (F OR H) AND F1 | A.pxCreateOperator = Param.CreateOP · A.pyStatusWork != "Resolved-Completed" · A.pyStatusWork != "Resolved-Rejected" · A.Quotation.BusinessType = "Life" · .TextNoQuotation StartsWith "NB-" · A.Quotation.Busines… | .pxCreateDateTime:DESC:1 | UI 1 |
| `GetMasterPolicyLife` | `Work-NB` | .pyID, .pzInsKey, .OfferFacIn.QuotationData.PolicyType, .OfferFacIn.QuotationData.BusinessType, .OfferFacIn.PersonList(1).AccumulationName | A AND B AND C | .OfferFacIn.QuotationData.PolicyType = 1 · .OfferFacIn.QuotationData.BusinessType = "Life" · .pyStatusWork = "Resolved-Completed" |  | UI 1 |
| `GetObjectItem` | `Int-V_JN_OBJ_ITEM` | .MJOI_KODE, .JN_OBJ_ITEM, .KETERANGAN | F1 | .MJOI_KODE StartsWith Param.JN_OBJ_ITEM | .JN_OBJ_ITEM:ASC:1 | **ALUR** 5 |
| `GetOperatorIDbyWorkBasket_RD` | `Data-Admin-Operator-ID` | .pyUserName, .pyUserIdentifier | (F1 OR F2 OR F3) AND F4 | .pyWorkbasket(1) = Param.Workbasket · .pyWorkbasket(2) = Param.Workbasket · .pyWorkbasket(3) = Param.Workbasket · .pyOperatorIsDeactivated IsFalse  |  | **ALUR** 1 |
| `GetTeamGroup_RD` | `Int-marketingofficer` | .ID, .TeamGroup, .ClientName, .MOLeader | (A OR B) AND C | .ID = Param.IDMarketing · .ClientName = Param.MktName · .MOStatus = 1 |  | **ALUR** 1 |
| `InwardScale_RD` | `Int-INWARDSCALE` | .ID, .Tahun, .ShareMin, .ShareMax, .PctTreatyLimit, .UserID |  |  | .ID:DESC:1 | **ALUR** 2 |
| `pyDefaultReport` | `@baseclass` | .pyLabel, @@pxDay(.pxCreateDateTime) | A | .pxUpdateDateTime = Last 7 Days |  | — 0 |
| `pyGetAllAttachments` | `Link-Attachment` | .pxLinkedRefFrom, .pxLinkedRefTo, .pxCreateDateTime, .pyCategory, .pyMemo, .pxLinkedClassFrom, .asmCategory2, .ASMClaimCategory, .pxListSubscript, .px… | A AND B | .pxLinkedRefFrom = Param.LinkRefFrom · .pyCategory = Param.Category | .pxCreateDateTime:DESC:1 | **ALUR** 1 |
| `pyGetListOfOperators` | `Data-Admin-Operator-ID` | .pyUserIdentifier, .pyUserName, .pyLabel, .pzInsKey, .pyEmailAddress | A | .pyUserIdentifier != "Administrator@pega.com" |  | UI 1 |
| `ReasGetAllAttachments` | `Link-Attachment` | .pxLinkedRefFrom, .pyCategory, .pxLinkedRefTo, .pzInsKey, .pyMemo, .PNOTE, .pxCreateDateTime | A AND B AND C | .pyCategory = Param.category · .pxLinkedRefFrom = Param.inskey · .PNOTE = Param.GCNMCategory | .pxCreateDateTime:DESC:1 | **ALUR** 5 |
| `SearchDueDate_RD` | `Int-LST_BUSINESS_DUE_DATE` | .LBU_ID_DUE_DATE, .DUE_DATE | A | .LBU_ID_DUE_DATE = Param.BizCode |  | **ALUR** 1 |
| `SearchRiskAccumulation_RD` | `Int-ACCUMULATION` | .Accumulation, .CZone, .CZoneID, .Note, .Keyword, .ScopeArea, .ID, .AccumulationName, .PostalCode, .ProvinceName, .ProvinceID | A AND B AND C AND D AND E AND F AND G | .ID = Param.AccumulationCode · .Note Contains Param.AccumulationDescription · .AccumulationName = Param.Coverage · .CZone = Param.Czone · .Keyword = Param.Keyword · .PostalCode = Param.PostalCode · .ProvinceID … |  | UI 3 |
| `SelectLeader_RD` | `Int-marketingofficer` | .ID, .ClientID, .BranchDetailID, .MOLeader, .MOStatus, .ClientName, .ClientID2, .BranchStatus, .TeamGroup, .Tanggal, .UserUpdate, .BranchParent, .Bran… | F1 AND A | .ClientID2 = "LEADER" · .MOStatus = Param.sts_aktif | .ID:DESC:1 | UI 1 |

---

## 5. Report yang menggerakkan alur (62)

Rujukan dibuktikan dengan dua mekanisme yang terbaca eksplisit di korpus:

**(a) Activity memanggil `pxRetrieveReportData`** — nama report diisikan ke `Param.pyReportName`:

```xml
<!-- NB FacIn/Activity/SetNBStatus_Act.xml, baris 907-908 -->
<PropertiesName>Param.pyReportName</PropertiesName>
<PropertiesValue>"GetOperatorIDbyWorkBasket_RD"</PropertiesValue>
```

```xml
<!-- NB FacIn/Activity/SearchInwardScale_Act.xml, baris 308-309 -->
<PropertiesName>Param.pyReportName</PropertiesName>
<PropertiesValue>"InwardScale_RD"</PropertiesValue>
```

**(b) DataPage menyatakan report sebagai sumbernya** — `pyLoadReportDefinition` +
`pyLoadActivity = pxCallRetrieveReportData`:

```xml
<!-- NB FacIn/DataPage/D_BrowseTableOfLimit.xml, baris 328-344 -->
<pyLoadReportDefinition>BrowseTableOfLimit_RD</pyLoadReportDefinition>
<pyDeclarePagesDataSource>ReportDefinition</pyDeclarePagesDataSource>
<pyReportDefinitionClass>ASM-FW-GISFW-Int-TABLEOFLIMIT</pyReportDefinitionClass>
<pyLoadActivity>pxCallRetrieveReportData</pyLoadActivity>
```

DataPage yang sama juga mendaftarkan rujukan formal:

```xml
<!-- NB FacIn/DataPage/D_BrowseTableOfLimit.xml, baris 443-446 -->
<pyRuleName>BrowseTableOfLimit_RD</pyRuleName>
<pxRuleObjClass>Rule-Obj-Report-Definition</pxRuleObjClass>
```

**[terverifikasi]** untuk kedua mekanisme.

### 5.1 Daftar lengkap pemanggil

| Report | Dipanggil dari |
| --- | --- |
| `BrowseAccumulatedType_RD` | `Activity/SetAccumulationData_Act` |
| `BrowseAgentHierarkiList_RD` | `Activity/AgentSourceBiz_Act` |
| `BrowseBankGroup` | `DataPage/D_BankGroup` |
| `BrowseBenefitClause` | `DataPage/D_BENEFIT_CLAUSE` |
| `BrowseBenefitPkg` | `Activity/IndMedPrintLoading_PreAct` |
| `BrowseCedingCo_RD` | `Activity/AgentSourceBiz_Act` |
| `BrowseCityInput_RD` | `DataPage/D_CITY` |
| `BrowseClauseAneka_RD` · `BrowseClauseMBU_RD` · `BrowseClausePA_RD` | `Activity/SearchClauseSQL_PreAct` (ketiganya) |
| `BrowseClausePlan` | `DataPage/D_PlanClause` |
| `BrowseColor_RD` | `Activity/UploadCSVSupportVehicleColor` |
| `BrowseConveyance_RD` | `Activity/ProtectUploadMarineCargo_act`, `Activity/SetConveyence_Act` |
| `BrowseCoverageAneka` | `Activity/BrowseCoverageSupport` |
| `BrowseCoverageFacIn_RD` | `Activity/ProtectUploadAneka_Act` |
| `BrowseCurrency_RD` | `Activity/AddCurencyList_ACT`, `AddCurrencyListPA_ACT`, `AddPaymentCurency_ACT`, `CountCoverageMarine`, `GetCurrencyMaster`, `SetCurrencyCoverage_Act` |
| `BrowseDistrictInput_RD` | `DataPage/D_DISTRICT` |
| `BrowseEnumData` | `Activity/UploadCSVSupportCoverageT` |
| `BrowseGoodsType` | `Activity/ProtectUploadMarineCargo_act`, `Activity/SetGoods_Act` |
| `BrowseJob_RD` | `Activity/UploadCSVSupportJob` |
| `BrowseJTypePropertyPlan` | `DataPage/D_J_PlanProperty` |
| `BrowseLLOYDAGENT_RD` | `Activity/ProtectUploadMarineCargo_act` |
| `BrowseM_OKUPASI_ANEKA` | `Activity/BrowseCoverageSupport`, `Activity/BrowseM_OKUPASI_ANEKA`, `Activity/UploadCSVAneka_PostAct` |
| `BrowseMaxHariTravel_RD` | `Activity/fillActPremiTravel` |
| `BrowseModelList` | `Activity/BrowseActV_VEHICLE_MODEL` |
| `BrowseNation_RD` | `Activity/SetCountryID_Act` |
| `BrowseOccupationANEKA_RD` | `Activity/ProtectUploadAneka_Act`, `DataPage/D_Occupation` |
| `BrowseOccupationFacInFIRE_RD` | `Activity/InsertUploadFire_act`, `DataPage/D_BrowseOccupationFacInFIRE` |
| `BrowseOccupationList_RD` | `Activity/UploadCSVSupportVehicleOccu` |
| `BrowsePacking_RD` | `Activity/ProtectUploadMarineCargo_act`, `Activity/SetPackingNote_Act` |
| `BrowsePlanByAge` | `Activity/IndMedPrintLoading_PreAct`, `Activity/IndMedPrintProposal_PreAct` |
| `BrowsePlanProperty` | `DataPage/D_J_PlanProperty`, `DataPage/D_PlanProperty` |
| `BrowsePremiTravel` | `Activity/fillActPremiTravel` |
| `BrowseProvince_RD` | `Activity/SetProvinceID_Act` |
| `BrowseReinsuranceType_RD` | `Activity/GeneratePDFFacOutMemoPlacing` |
| `BrowseRetrocessionLife_RD` | `Activity/changepercentrnm_act`, `Activity/changepercentrnmLife_act`, `Activity/SetLayerSpreading_act` |
| `BrowseRiskAddress_RD` | `Activity/ProtectUploadAneka_Act` |
| `BrowseRWInput_RD` | `DataPage/D_RWNAME` |
| `BrowseSameAccumulationLife_RD` | `Activity/createNewAccumulation_act` |
| `BrowseSearchSobCeding_RD` | `Activity/SetDataSobCeding_Act` |
| `BrowseShip_RD` · `BrowseShip2_RD` | `Activity/ProtectUploadMarineCargo_act` (keduanya) |
| `BrowseSubContract_RD` | `Activity/SearchSubContractSQL_PreAct` |
| `BrowseSubjectTo_RD` | `Activity/SetSubjectToAct` |
| `BrowseTableOfLimit_RD` | `DataPage/D_BrowseTableOfLimit` |
| `BrowseTrading_RD` | `Activity/ProtectUploadMarineCargo_act`, `Activity/SetTrading_Act` |
| `BrowseTreatyYear_Life_RD` | `Activity/changepercentrnm_act`, `Activity/changepercentrnmLife_act`, `Activity/GenerateSpreadingLife_act` |
| `BrowseTreatyYearOR_Life_RD` | `Activity/GenerateSpreadingLifeOR_act` |
| `BrowseTypeList_RD` | `Activity/UploadCSVSupportVehicle` |
| `BrowseVClauseArg` | `Activity/InputClause_ViewArg_PreAct` |
| `BrowseVMarketingLeaderList` | `Activity/IndMedPrintLoading_PreAct` |
| `CheckMST_KENDARAANByChas` · `ByEngine` · `ByPlat` | `Activity/ValidateDuplicateVehicle` (ketiganya) |
| `DataTableEditorReport` | `DataPage/D_EnumerationList` |
| `GetObjectItem` | `Activity/ProtectFIREMBUPA_Act`, `SetDataFacOutFire_Act`, `SetObjItemType_Act`, `UploadDoc_AI_V2`, `RDBList/GetObjectItembyName_SQL` |
| `GetOperatorIDbyWorkBasket_RD` | `Activity/SetNBStatus_Act` |
| `GetTeamGroup_RD` | `Activity/GetTeamGroup_Act` |
| `InwardScale_RD` | `Activity/SearchInwardScale_Act` |
| `pyGetAllAttachments` | `Activity/LoadAttachmentData` |
| `ReasGetAllAttachments` | `Activity/ViewAttach_RD` |
| `SearchDueDate_RD` | `Activity/FillActInstallment` |

### 5.2 Yang paling penting untuk keputusan

Enam report ini dipakai mengambil keputusan, bukan menampilkan daftar:

#### `GetOperatorIDbyWorkBasket_RD` — routing penugasan
Kelas `Data-Admin-Operator-ID` (tabel operator internal Pega), `pyMemo` = `Copied from prod`,
v01-01-87. Dipanggil `Activity/SetNBStatus_Act` — aktivitas penetapan status NB.

- Kolom: `.pyUserName`, `.pyUserIdentifier`
- `pyFilterLogic`: `(F1 OR F2 OR F3) AND F4`
- Filter: `.pyWorkbasket(1) = Param.Workbasket` · `.pyWorkbasket(2) = …` · `.pyWorkbasket(3) = …` ·
  `.pyOperatorIsDeactivated IsFalse`

**Satu operator hanya dicari pada tiga slot workbasket pertama.** Bila operator memiliki workbasket
ke-4 dan seterusnya, ia tidak akan ditemukan. **[terverifikasi]** atas isi filter;
**[pertanyaan terbuka]** apakah batas 3 itu disengaja.
Tabel `Data-Admin-Operator-ID` adalah tabel internal Pega → menurut aturan proyek §4.3 tidak boleh
dimigrasikan apa adanya; harus diganti direktori operator sistem baru.

#### `BrowseTableOfLimit_RD` — tabel limit akseptasi per lini bisnis
Kelas `Int-TABLEOFLIMIT`, v01-01-80, `pyMemo` = `distinct`. Dibaca lewat `DataPage/D_BrowseTableOfLimit`.

- Kolom: `.Bizcode`, `.Category`, `.Description`, `.PctLimit`, `.Note`
- `pyFilterLogic`: `A AND B AND C AND D` → `.Tahun = Param.Tahun` · `.Bizcode = Param.Bizcode` ·
  `.Category = Param.Category` · `.ID = Param.ID`
- Urutan: `.Category ASC` lalu `.Description ASC`

`.PctLimit` **tidak** difilter, hanya dikembalikan → limit dinyatakan sebagai persentase, dan
pemilihan barisnya dilakukan oleh pemanggil. **[terverifikasi]** atas struktur;
**[pertanyaan terbuka]** basis persentase `.PctLimit` (dari TSI? dari kapasitas treaty?).

#### `BrowseLimit` — pita limit treaty proporsional
Kelas `Int-PROPORTIONALARRG` → tabel `pooldata.proportionalarrg` **[terverifikasi]**. v01-01-53.
Hanya dirujuk dari `Harness`(1) + `Section`(5) — tidak dari Activity.

- Kolom: `.ParentReinsTypeID`, `.ReinsTypeID`, `.Rp`, `.Usd`, `.TreatyYear`, `.TreatyGroupID`,
  `.TreatyDescID`, `.ID`, `.ReinsTypeName`, `.Line`
- `pyFilterLogic`: `(A OR G) AND B AND C AND D AND E AND F AND H`
- Filter: `.Rp > Param.Rp` (A) · `.Rp <= Param.Rupiah` (B) · `.TreatyYear = Param.TreatyYear` (C) ·
  `.TreatyGroupID = Param.TreatyGroupID` (D) · `.ID = Param.ID` (E) · **`.Usd = .Usd`** (F) ·
  `.Rp IS NOT NULL` (G) · `.Rp > Param.RPFacout` (H)

**Anomali**: filter F membandingkan kolom dengan dirinya sendiri (`.Usd = .Usd`) — selalu benar
kecuali bila `NULL`, sehingga efeknya setara `\.Usd IS NOT NULL`. **[terverifikasi]** atas isi tag;
**[dugaan]** bahwa ini tidak disengaja. Kandidat perbaikan, **bukan** bagian migrasi.
Nilai limit disimpan dalam dua kolom mata uang terpisah (`.Rp`, `.Usd`) — konsisten dengan aturan
proyek §4.1 bahwa mata uang wajib menyertai jumlah.

#### `InwardScale_RD` — skala share inward
Kelas `Int-INWARDSCALE`, v01-01-52, `pyMemo` = `remove condition`. Dipanggil
`Activity/SearchInwardScale_Act`.

- Kolom: `.ID`, `.Tahun`, `.ShareMin`, `.ShareMax`, `.PctTreatyLimit`, `.UserID`
- **`pyFilterLogic` kosong; `pyFilter` kosong** — tidak ada filter sama sekali
- Parameter yang *dideklarasikan*: `Tahun`, `ID`, `ShareMin`, `ShareMax` — **keempatnya tidak
  dipakai oleh filter mana pun**
- Urutan: `.ID DESC`; `pyMaxRecords` = 500

Report ini mengembalikan **seluruh** tabel inward scale (dibatasi 500 baris) dan menyerahkan
penyaringan kepada pemanggil. `pyMemo` "remove condition" menunjukkan filter pernah ada lalu
dihapus. **[terverifikasi]**. Konsekuensi migrasi: bila tabel melampaui 500 baris, hasil terpotong
diam-diam. **[pertanyaan terbuka]** apakah batas 500 itu disengaja.

#### `SearchRiskAccumulation_RD` — kontrol akumulasi risiko
Kelas `Int-ACCUMULATION`, v01-01-84, `pyMemo` = `remve duplicate`. Dirujuk dari `Harness`(1) +
`Section`(2) saja.

- Kolom: `.Accumulation`, `.CZone`, `.CZoneID`, `.Note`, `.Keyword`, `.ScopeArea`, `.ID`,
  `.AccumulationName`, `.PostalCode`, `.ProvinceName`, `.ProvinceID`
- `pyFilterLogic`: `A AND B AND C AND D AND E AND F AND G` — tujuh filter, semuanya `AND`:
  `.ID = Param.AccumulationCode` · `.Note Contains Param.AccumulationDescription` ·
  `.AccumulationName = Param.Coverage` · `.CZone = Param.Czone` · `.Keyword = Param.Keyword` ·
  `.PostalCode = Param.PostalCode` · `.ProvinceID = Param.Province`

`RDBList` sekelas (`Int-ACCUMULATION`) menyentuh `facinproduction` dan `pooldata.json_polis`
**[terverifikasi]**, tetapi **report ini tidak** — ia hanya membaca tabel akumulasi.

#### `GetMasterPolicyLife` — satu-satunya report atas kelas Work siklus NB
Kelas `ASM-FW-GISFW-Work-NB`, v01-01-61, `pyMemo` = `set time 200 sec`.
Tidak ditemukan pemanggil di folder NB (lihat §6).

- Kolom: `.pyID`, `.pzInsKey`, `.OfferFacIn.QuotationData.PolicyType`,
  `.OfferFacIn.QuotationData.BusinessType`, `.OfferFacIn.PersonList(1).AccumulationName`
- `pyFilterLogic`: `A AND B AND C` → `.OfferFacIn.QuotationData.PolicyType = 1` ·
  `.OfferFacIn.QuotationData.BusinessType = "Life"` · `.pyStatusWork = "Resolved-Completed"`

Ini bukti langsung tiga enumerasi di agregat `OfferFacIn`: `PolicyType` bernilai numerik (`1`),
`BusinessType` bernilai teks (`"Life"`), `pyStatusWork` memakai status Pega baku
(`"Resolved-Completed"`). **[terverifikasi]** atas nilai yang dibandingkan.
**[pertanyaan terbuka]**: arti `PolicyType = 1` — **belum terverifikasi**.
Perhatikan `PolicyType` dibandingkan tanpa tanda kutip sedangkan `BusinessType` dengan tanda kutip;
apakah `PolicyType` numerik atau teks **belum terverifikasi**.

---

## 6. Report yang menyentuh tabel produksi / limit / riwayat akseptasi

```powershell
$p="D:\migrasi\RNM\NB FacIn\ReportDefinition"
foreach($t in 'HISTORYAKSEPTASI','FACINPRODUCTION','JSON_POLIS','PRODKE','DATAPEGA','PC_WORK','M_LINK_SERVICE'){
  "{0,-18}: {1}" -f $t, @(Select-String -Path "$p\*.xml" -Pattern $t -SimpleMatch -List).Count }
# -> semuanya 0
```

**Hasil: nol.** Tidak satu pun dari 123 ReportDefinition menyebut `HISTORYAKSEPTASIPEGA`,
`FACINPRODUCTION`, `JSON_POLIS`, `PRODKE`, `DATAPEGA.PC_*`, atau `M_LINK_SERVICE`.
**[terverifikasi]**

> Peringatan pola: pencarian kata `LIMIT` mengembalikan 123/123 — karena setiap berkas memuat tag
> metadata `<pxLimitedAccess>`. Angka itu **bukan** bukti apa pun tentang tabel limit. Kesimpulan di
> atas memakai nama tabel yang spesifik, bukan kata `LIMIT`.

Yang tersisa, yang menyentuh domain limit/treaty/produksi lewat **kelas**-nya:

| Report | Kelas | Tabel (bila terbukti) | Peran |
| --- | --- | --- | --- |
| `BrowseLimit` | `Int-PROPORTIONALARRG` | `pooldata.proportionalarrg` **[terverifikasi]** | pita limit treaty |
| `BrowseTableOfLimit_RD` | `Int-TABLEOFLIMIT` | belum terverifikasi | limit % per lini bisnis |
| `InwardScale_RD` | `Int-INWARDSCALE` | belum terverifikasi | skala share inward |
| `SearchRiskAccumulation_RD` | `Int-ACCUMULATION` | `accumulation` **[terverifikasi]** | kontrol akumulasi |
| `BrowseRetrocessionLife_RD` | `Int-RETROCESSIONLIFE` | belum terverifikasi | share retro per treaty-year Life |
| `BrowseTreatyYear_Life_RD`, `BrowseTreatyYearOR_Life_RD` | `Int-TREATYYEAR_LIFE` | belum terverifikasi | tahun treaty Life |
| `BrowseReinsuranceType_RD` | `Int-REINSURANCETYPE` | `m_reinsurancetype` / `reinsurancetype` **[terverifikasi]** | jenis reasuransi |
| `GetMasterPolicyLife` | `Work-NB` | tabel work Pega (`PC_*`) **[dugaan]** | polis master Life selesai |
| `GetListOpportunity*` (3) | `SFAGISFW-Work-Opportunity` | tabel work Pega (`PC_*`) **[dugaan]** | antrian opportunity |
| `GetOperatorIDbyWorkBasket_RD`, `pyGetListOfOperators` | `Data-Admin-Operator-ID` | tabel operator Pega | routing & daftar operator |
| `pyGetAllAttachments`, `ReasGetAllAttachments` | `Link-Attachment` | tabel link Pega | lampiran |
| `GetAllDocument` | `Int-DOCUMENT_POLIS` | belum terverifikasi | dokumen polis (`.T_STORAGE_ID`, `.INSKEY_LINK`, `.MIME`, `.NAMAFILE`, `.KATEGORI_2`) |

**Konsekuensi migrasi**: tujuh report (`GetMasterPolicyLife`, 3 × `GetListOpportunity*`, 2 ×
operator, 2 × attachment) membaca tabel internal Pega. Menurut aturan proyek §4.3 tabel `PC_*`
**hilang bersama Pega** — report-report ini harus ditulis ulang di atas model data baru, bukan
dipetakan satu-satu.

---

## 7. Report yatim — dan yang bukan yatim

```powershell
# untuk setiap nama report, cari kemunculannya di seluruh berkas NB selain folder ReportDefinition
# -> YATIM(0 rujukan) = 1 ; DIRUJUK = 122
```

**Satu report yatim: `pyDefaultReport`** — kelas `@baseclass`, ruleset `Pega-Reporting`,
v08-01-01, `pyMemo` = `US-106351`. Isinya: kolom `.pyLabel` dan `@@pxDay(.pxCreateDateTime)`,
filter `.pxUpdateDateTime = Last 7 Days`. Ini **rule bawaan produk Pega** (laporan drill-down
default), bukan kode Nusantara Re. Ia muncul di setiap berkas ReportDefinition lain sebagai
`pyDrillDownReportName`, tetapi tidak dipanggil dari kode. **Bukan kandidat migrasi.**
**[terverifikasi]**

### 7.1 "Belum ditemukan rujukannya" — dibedakan tegas dari yatim

Deteksi rujukan di atas adalah pencocokan teks nama rule di seluruh 1.960 berkas non-ReportDefinition
folder NB. Ada dua sumber rujukan yang **tidak ada di korpus** dan karenanya tidak terdeteksi:

1. **Portal & Access Group.** Report dapat dipasang langsung pada portal/landing page yang
   konfigurasinya adalah `Data-Admin-Operator-AccessGroup` / `Rule-Access-Role-Obj` — tipe rule
   yang **tidak diekspor** ke folder mana pun.
2. **Siklus lain.** Report yang hanya dipanggil dari RNW/EDM akan tampak tanpa rujukan di NB.

Karena itu, report berikut — yang di §5.1 **tidak** muncul dan hanya dirujuk dari `Section` —
tetap harus dicatat sebagai **belum ditemukan rujukan yang menggerakkan logika**, bukan kode mati:

`GetMasterPolicyLife` (0 pemanggil di NB), `GetListOpportunity`, `GetListOpportunityF`,
`GetListOpportunityLife`, `pyGetListOfOperators`, `SelectLeader_RD`, `BrowseAgent_RD`,
`BrowseAgentNonLife_RD`, `BrowseAgentNusaRe_RD`, `BrowseAgentNusaReLife_RD`, `BrowseBranch_RD`,
`BrowseBusiness_RD`, `BrowseCoverage_RD`, `BrowseDistrictInputC_RD`, `BrowseObjectItemType_RD`,
`BrowseOccupations_RD`, `BrowseProvince2_RD`, `BrowseRisksAddress_RD`, `BrowseSubjectTo_RD`,
`BrowseV_ALM_RISIKOList`, `BrowseWarranty_RD`, dan seterusnya.

**[pertanyaan terbuka]** — daftar report yang benar-benar mati hanya dapat ditutup setelah ekspor
`Rule-Access-*` dan konfigurasi portal tersedia.

---

## 8. Anomali yang ditemukan (kandidat perbaikan, bukan bagian migrasi)

Semua butir di bawah **wajib direplikasi apa adanya** pada migrasi agar rekonsiliasi paralel run
menghasilkan angka yang sama. Perbaikannya adalah keputusan bisnis.

| # | Temuan | Berkas | Label |
| --- | --- | --- | --- |
| 1 | **Dua konvensi "cedant aktif" yang bertentangan.** `BrowseCedingCo_RD` (`pyMemo`: `.StatusActive is null`) menyaring `\.StatusActive IS NULL`; `BrowseSearchSobCeding_RD` (`pyMemo`: `ubah status aktif`) menyaring `\.StatusActive = 1`. Keduanya kelas `Int-AGENT`, keduanya dipanggil dari Activity (`AgentSourceBiz_Act` dan `SetDataSobCeding_Act`). | `ReportDefinition/BrowseCedingCo_RD.xml`, `…/BrowseSearchSobCeding_RD.xml` | [terverifikasi] |
| 2 | **Perbandingan angka sebagai string.** `BrowseClieent_RD` menyaring `Agent.StatusActive != "0"` — status numerik dibandingkan dengan literal string. | `ReportDefinition/BrowseClieent_RD.xml` | [terverifikasi] |
| 3 | **Filter tautologi.** `BrowseLimit` memuat `\.Usd = .Usd`. | `ReportDefinition/BrowseLimit.xml` | [terverifikasi] |
| 4 | **Pengecualian data bisnis ditanam di kode.** `BrowseCoverage_RD` menyaring `\.Name NotContain "MACHINERY BREAK DOWN + 4.1ACC"` — satu nama cover tertentu di-blacklist di dalam definisi report. | `ReportDefinition/BrowseCoverage_RD.xml` | [terverifikasi] |
| 5 | **Literal tanpa/dengan tanda kutip tidak konsisten.** `GetListOpportunity` memakai `A.Quotation.BusinessFac = T` dan `GetListOpportunityF` memakai `= F` (tanpa kutip), sedangkan `GetListOpportunityLife` memakai `= "F"` (dengan kutip). Arti `T`/`F` **belum terverifikasi** — jangan diasumsikan true/false. | `ReportDefinition/GetListOpportunity.xml`, `…F.xml`, `…Life.xml` | [terverifikasi] atas isi tag |
| 6 | **Report menyaring terhadap halaman clipboard, bukan parameter.** `BrowseZones_RD` menyaring `\.Zone = Zone.CARI1` … `Zone.CARI6` — enam slot pencarian pada halaman bernama `Zone`. Report jadi bergantung pada state clipboard pemanggil. | `ReportDefinition/BrowseZones_RD.xml` | [terverifikasi] |
| 7 | **Parameter dideklarasikan tetapi tidak dipakai.** `InwardScale_RD` mendeklarasikan 4 parameter dan tidak punya filter sama sekali. Pola serupa pada `BrowseAgent_RD` (5 parameter, 1 filter `\.ClientID IS NOT NULL`). | `ReportDefinition/InwardScale_RD.xml`, `…/BrowseAgent_RD.xml` | [terverifikasi] |
| 8 | **`pyFilterLogic` memakai presedensi campur tanpa kurung.** `BrowseCedingCo_RD` = `A OR B AND C`. Hasilnya bergantung pada presedensi `AND` di atas `OR` yang **tidak dinyatakan** di korpus. | `ReportDefinition/BrowseCedingCo_RD.xml` | [terverifikasi] atas tag; presedensi **[pertanyaan terbuka]** |
| 9 | **Batas baris seragam 500 pada hampir semua report**, dengan pengecualian 10.000 pada `BrowseCedingCo_RD`, `BrowseSearchSobCeding_RD`, `BrowseAgentHierarkiList_RD`. Pemotongan terjadi diam-diam. | seluruh folder | [terverifikasi] |

---

## 9. Usulan pemetaan ke target Go + React **[usulan]**

Seluruh isi bagian ini adalah **usulan**, bukan temuan korpus.

### 9.1 Prinsip

1. **ReportDefinition bukan satu konsep tunggal.** Korpus memakainya untuk tiga hal berbeda, dan
   ketiganya harus dipetakan berbeda:

   | Pemakaian di Pega | Ciri | Target |
   | --- | --- | --- |
   | **Lookup / autocomplete** (mayoritas: ~95 report) | dipanggil dari `Section` sebagai sumber grid/dropdown | endpoint lookup generik, cache di sisi klien |
   | **Sumber DataPage** (7 report) | `pyLoadReportDefinition` pada DataPage | lapisan `repository/lookup/`, dipanggil service |
   | **Query keputusan** (mis. `GetOperatorIDbyWorkBasket_RD`, `BrowseLimit`) | dipanggil Activity untuk mengambil keputusan | query bernama di `repository/oracle/`, **bukan** endpoint publik |

2. **Jangan membuat 123 endpoint.** 84 kelas sumber → satu keluarga endpoint lookup berparameter
   lebih tepat daripada 123 rute.

3. **Ketertelusuran wajib** (aturan proyek §4.6). Setiap query menyebut report asalnya:

   ```go
   // Asal: NB FacIn/ReportDefinition/BrowseTableOfLimit_RD.xml
   //       (dipanggil lewat NB FacIn/DataPage/D_BrowseTableOfLimit.xml)
   // Filter asli: Tahun = :tahun AND Bizcode = :bizcode AND Category = :kategori AND ID = :id
   // Urutan asli: Category ASC, Description ASC ; pyMaxRecords = 500
   const qTableOfLimit = `SELECT … `
   ```

### 9.2 Usulan endpoint Go

| Kelompok | Endpoint usulan | Report sumber |
| --- | --- | --- |
| Lookup generik | `GET /api/v1/lookup/{katalog}?q=&page=` | ~95 report bertema A, D, E, F, G |
| Limit & treaty | `GET /api/v1/limit/table-of-limit?tahun=&bizcode=&kategori=` | `BrowseTableOfLimit_RD` |
| | `GET /api/v1/limit/inward-scale?tahun=` | `InwardScale_RD` (catatan: aslinya tanpa filter) |
| | `GET /api/v1/limit/treaty-band?treaty_year=&group_id=&amount=` | `BrowseLimit` |
| Akumulasi | `GET /api/v1/akumulasi/risiko?kode=&czone=&provinsi=&kodepos=` | `SearchRiskAccumulation_RD` |
| Cedant | `GET /api/v1/cedant?leader=&id=` | `BrowseCedingCo_RD`, `BrowseSearchSobCeding_RD`, `BrowseAgentHierarkiList_RD` |
| Dokumen | `GET /api/v1/offer/{id}/dokumen?kategori=` | `GetAllDocument` |
| Lampiran | `GET /api/v1/offer/{id}/lampiran?kategori=` | `ReasGetAllAttachments`, `pyGetAllAttachments` |
| Antrian NB | `GET /api/v1/nb/antrian?search=&team=` | `GetListOpportunity*` — **tulis ulang**, sumber lama tabel Pega |
| Routing (internal) | tanpa endpoint — fungsi `repository` | `GetOperatorIDbyWorkBasket_RD` |

Aturan §4.1 berlaku pada semua endpoint di atas: nilai uang (`.Rp`, `.Usd`, `.MinTSI`, `.MaxTSI`,
`.PctLimit`, `.PctTreatyLimit`) dikirim sebagai **string desimal** disertai kode mata uang, bukan
`number` JSON.

### 9.3 Usulan komponen React

| Komponen | Untuk | Catatan |
| --- | --- | --- |
| `<LookupSelect catalog=… />` | ~95 report lookup | satu komponen, katalog sebagai prop; menggantikan ~95 grid Pega |
| `<LookupModalGrid catalog=… columns=… />` | report lookup berkolom banyak (`BrowseDeductibleFacIn_RD` 25 kolom, `BrowseOC_RD` 25, `BrowseAgent_RD` 24, `BrowseBrandDetail_RD` 24) | pencarian di server, bukan di klien |
| `<CascadingAddress />` | 16 report tema A | `BrowseNation_RD` → `BrowseProvince_RD` → `BrowseCity_RD` → `BrowseDistrict_RD` → `BrowseTeritory_RD` / `BrowseRW_RD` — rantai berjenjang yang terbaca dari parameter tiap report |
| `<LimitTable />` | `BrowseTableOfLimit_RD`, `BrowseLimit`, `InwardScale_RD` | read-only; angka sebagai string desimal |
| `<AccumulationSearch />` | `SearchRiskAccumulation_RD` | 7 filter `AND` |
| `<AttachmentList />` | `ReasGetAllAttachments` | urut `.pxCreateDateTime DESC` |
| `<NBWorklist />` | `GetListOpportunity*` | **desain ulang**, bukan port |

Rantai berjenjang pada baris ketiga **[terverifikasi]** dari parameter masing-masing report:
`BrowseProvince_RD(Nation, PROVINCENAME)` → `BrowseCity_RD(Nation, Province)` →
`BrowseDistrict_RD(Nation, Province, City)` → `BrowseTeritory_RD(Nation, Province, City, District)`.

### 9.4 Yang jangan dipetakan

- `pyDefaultReport` — rule produk Pega.
- `pyGetListOfOperators`, `GetOperatorIDbyWorkBasket_RD` — tabel operator Pega; ganti dengan
  direktori identitas sistem baru.
- `pyGetAllAttachments` — versi lama; `ReasGetAllAttachments` adalah versi yang dipakai jalur Reas
  (v01-01-81 vs v01-01-52, dan hanya `ReasGetAllAttachments` yang punya parameter `GCNMCategory`).
  **[dugaan]** — perlu konfirmasi sebelum salah satunya dibuang.
- 31 report tema F dan G, **bila** bisnis mengonfirmasi jalur Fac In tidak memakainya. Sampai itu
  terjadi, 16 di antaranya masih dirujuk dari kode NB dan tidak boleh dihapus.

---

## 10. Pertanyaan terbuka

**MEMBLOKIR** — implementasi NB tidak dapat diselesaikan dengan benar tanpa jawaban:

1. **Nama tabel Oracle untuk 64 kelas `ASM-FW-GISFW-Int-*`** yang tidak punya `RDBList` sekelas.
   Definisi kelas Pega tidak ada di korpus. Tanpa ini, tidak ada satu pun endpoint lookup yang dapat
   ditulis. Yang paling mendesak: `Int-TABLEOFLIMIT`, `Int-INWARDSCALE`, `Int-RETROCESSIONLIFE`,
   `Int-TREATYYEAR_LIFE`, `Int-DOCUMENT_POLIS`, `Int-COVERAGE_FACIN`, `Int-DEDUCTIBLE`.
2. **Basis persentase `.PctLimit` (`BrowseTableOfLimit_RD`) dan `.PctTreatyLimit`
   (`InwardScale_RD`)** — persen dari apa? Tanpa ini tangga limit tidak dapat dihitung ulang.
3. **Mana konvensi "cedant aktif" yang benar** — `StatusActive IS NULL` atau `StatusActive = 1`?
   Keduanya hidup berdampingan dan dipanggil dari Activity yang berbeda (anomali §8 #1). Pilihan
   yang salah mengubah daftar cedant yang boleh dipilih pada pembuatan offer.
4. **Arti `PolicyType = 1`** pada `GetMasterPolicyLife` — belum terverifikasi.
5. **Arti kode `T` dan `F`** pada `A.Quotation.BusinessFac` (`GetListOpportunity*`) — belum
   terverifikasi; jangan diasumsikan boolean.

**TIDAK MEMBLOKIR** — perlu dijawab sebelum go-live, tidak menghalangi implementasi:

6. Apakah batas 3 slot workbasket pada `GetOperatorIDbyWorkBasket_RD` disengaja?
7. Apakah batas `pyMaxRecords` 500 (dan 10.000 pada tiga report cedant) adalah aturan bisnis atau
   nilai default yang tak pernah ditinjau?
8. Presedensi operator pada `pyFilterLogic` tanpa kurung (mis. `A OR B AND C` pada
   `BrowseCedingCo_RD`) — konfirmasi ke dokumentasi Pega, bukan ke korpus.
9. Apakah blacklist literal satu nama cover pada `BrowseCoverage_RD` (§8 #4) masih relevan?
10. Apakah 31 report tema F (objek lini langsung) dan G (travel/PA) memang tidak dipakai jalur
    Fac Inward, sehingga boleh ditinggalkan?
11. `pyGetAllAttachments` vs `ReasGetAllAttachments` — mana yang dipensiunkan?
12. Daftar report yang benar-benar mati tidak dapat ditutup tanpa ekspor `Rule-Access-*` dan
    konfigurasi portal (lihat §7.1).
