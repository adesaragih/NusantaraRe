# Rancangan Tabel Flat — Pengganti `JSON_POLIS.DATA_JSON`

> ## 🔄 SEDANG DIREKONSILIASI — 24 September 2026 · **belum digantikan**
>
> Ada rancangan kedua yang dibangun dari **114 contoh** (`Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx`,
> keputusan V-1…V-50). Keduanya **bertentangan di sembilan titik**.
>
> ⛔ **Dokumen ini MASIH BERLAKU** — ia dipakai menyusun tiket, dan belum ada keputusan work owner
> yang menggantikannya. Jangan memakai rancangan baru sebagai sumber kebenaran sampai kesembilan
> butir diputuskan.
>
> **Hasil pemeriksaan kesembilan butir: `00-rekonsiliasi-rancangan-flat.md`.**
>
> Satu `[pertanyaan terbuka]` dokumen ini **sudah tertutup** oleh pemeriksaan itu: **kunci versi**
> (§3.1) — `[terverifikasi]` `JSON_POLIS` ber-`PRIMARY KEY (IDPEGA)` tunggal, sehingga dua versi
> dengan `IDPEGA` sama mustahil. Prinsip §1 *"semua tabel anak membawa `IDPEGA` + `PRODKE`"* ikut
> terpengaruh — lihat butir 5 di berkas rekonsiliasi.

Dokumen ini merancang tabel flat relasional untuk menggantikan CLOB JSON bersarang
`JSON_POLIS.DATA_JSON`. Rancangan dibangun **dari struktur nyata** lima berkas kasus produksi,
bukan dari asumsi.

> Kepatuhan CLAUDE.md §3: setiap tabel di bawah menyebut **sumber** (berkas + jalur kunci JSON +
> `pxObjClass` pembawa) dan berlabel `[terverifikasi]` bila strukturnya terbaca langsung di korpus.
> **Tidak ada nilai PII** (nama orang, nomor polis, alamat, email, nomor slip) yang disalin; hanya
> nama field dan tipe. Nilai contoh yang muncul hanya angka/kode struktural netral.

---

## 0. Sumber yang dipakai — [terverifikasi]

Lima berkas contoh di `D:\migrasi\RNM\DDL\` (isi = string `DATA_JSON` utuh dari `JSON_POLIS`):

| Berkas | COB | Ciri struktur bawah `LocationList`/root | Ukuran |
| --- | --- | --- | --- |
| `P-5 NB-176005 ( AS. KREDIT ).txt` | ANEKA (Asuransi Kredit) | `LocationList→Property→RiskLocation→OccupationList→AnekaList→CoverageList→SpreadingList` | 11 KB |
| `P-5 NB-181231 (FIRE).txt` | FIRE (Property All Risk) | `LocationList→PropertyItemList→CoverageList→{DeductibleList,LayerList,SpreadingList}` + `ScoringRisk` besar | 38 KB |
| `P-5 NB-184233 ( MARINE CARGO).txt` | MARINE CARGO | `CargoList→CoverageList→{DeductibleList,SpreadingList,AdditionalCoverage}` di **root**; `PolicyData.Ship` | 282 KB |
| `P-5 RNW-10579 (FIRE).txt` | FIRE (Renewal) | **struktur identik** dengan NB FIRE, 3× LocationList (multi-versi/riwayat dalam satu dokumen) | 105 KB |
| `P-5 EDM-13445 (FIRE).txt` | FIRE (Endorsement) | struktur identik NB FIRE + `Int-INWARDSCALE` + field `*Old`/`IsOldData` | 273 KB |

**Perintah audit** (menghasilkan inventaris kunci-list + `pxObjClass` di atas):
`OUTPUT\_tmp_list_keys.ps1` (regex `"([A-Za-z0-9_]+List)":\[` dan `"pxObjClass":"([^"]+)"`).

### 0.1 Tiga temuan yang mengarahkan rancangan

1. **[terverifikasi] RNW dan EDM memakai struktur yang SAMA dengan NB.** Himpunan `pxObjClass`
   RNW-10579 dan EDM-13445 identik dengan NB FIRE (hanya EDM menambah `Int-INWARDSCALE`). Ini
   membenarkan keputusan work owner: RNW/EDM menyalin data dari polis NB. **Konsekuensi rancangan:
   satu skema flat melayani ketiga siklus** — dibedakan oleh `PRODKE` (versi), bukan tabel terpisah.

2. **[terverifikasi] Bentuk anak `LocationList` berbeda per COB.** FIRE memakai `PropertyItemList`;
   ANEKA memakai `OccupationList→AnekaList`; MARINE tidak memakai `LocationList` sama sekali,
   melainkan `CargoList` di root. **Konsekuensi: `CoverageList` adalah simpul konvergen** — semua COB
   berujung pada `Coverage` + `SpreadingList`. Tabel `Coverage` dan `Spreading` menampung kolom FK
   yang menunjuk induk yang berbeda-beda (property-item / aneka / cargo), diisi sesuai COB.

3. **[terverifikasi] Field EDM (`IsOldData`, `PremiumOld`, `PremiNusantaraReOld`, `TSIOld`,
   `ProRatePercent`) hadir di level `Coverage`.** Terlihat di MARINE CARGO dan konsisten dengan
   before-image EDM (Seam 4). Kolom-kolom `*Old` ini dipertahankan berpasangan dengan nilai berjalan,
   seukuran pasangan `_MENJADI`/`_SELISIH` yang diminta CLAUDE.md §4.3.

---

## 1. Prinsip perancangan yang mengikat

- **[K-010/K-012] Uang bukan `float`.** Setiap nilai uang jadi `NUMBER` di Oracle + kolom mata uang
  pendamping (dari `Currency.Name`/`OldID` pada simpul terdekat). Di JSON API: string desimal.
- **[K-027] Nilai masuk sebagai string berkoma/berdesimal panjang** (mis.
  `"41744600.00000000000000000000"`). Kolom `NUMBER` Oracle menyimpan presisinya; konversi string↔
  desimal terjadi di lapisan repository, tidak menebak pembulatan.
- **[§4.3] `PRODKE` tetap penomoran versi; riwayat = baris bertambah.** Semua tabel anak membawa
  `IDPEGA` + `PRODKE` (mewarisi dari header) agar satu versi polis dapat dipugar utuh dan versi lama
  tetap ada. **Tidak ada update in-place.**
- **[§4.3] Pasangan `_MENJADI`/`_SELISIH` dan `*Old` dipertahankan berpasangan**, tidak diratakan
  jadi satu nilai.
- **Kunci alami dari JSON** (`IndexProperty`, `IndexPropertyItem`, `IndexCoverage`, `IndexCargo`,
  `IdxLocation`, `IdxOccupation`, `IdxAneka`, `pxListSubscript`) **dipertahankan sebagai kolom** —
  itulah yang menautkan Coverage ke induknya di data asli. Jangan mengganti dengan surrogate saja.
- **[§4.5 gagal keras] Arti enumerasi belum dipetakan** (`TreatyType`, `CalculateMethod_FacIn`,
  `CoverageBasis`, `TypeDeductible`, `ConveyanceID`, dll.) — kolom disimpan **apa adanya sebagai kode**,
  penerjemahan menunggu keputusan bisnis, bukan ditebak.
- **PII tidak dipetakan sebagai kolom yang wajib diisi dari korpus** — nama orang/insured/alamat tetap
  ada sebagai kolom (bentuk dipertahankan) tetapi **tidak diisi contoh** di dokumen ini.

---

## 2. Peta relasi (ringkas)

```
JSON_POLIS_HEADER (1 baris per IDPEGA+PRODKE)  ← menggantikan baris JSON_POLIS
├── FLAT_QUOTATION            (1:1)   QuotationData
├── FLAT_POLICY               (1:1)   PolicyData (+ Ship utk MARINE)
├── FLAT_PAYMENT_HEADER       (1:1)   PolicyData.Payment
├── FLAT_CEDING_CEDANT        (1:N)   CedingCedantList (+ CurrencyList anak)
├── FLAT_CURRENCY_TOP         (1:N)   CurrencyList (root) → Policy.Payment → ListInstallment
│   └── FLAT_INSTALLMENT      (1:N)   ListInstallment
├── FLAT_TOTAL_TSI_PREMI      (1:N)   TotalTSIPremiGrossList / TotalTSIList
├── FLAT_SPREAD_RNM           (1:N)   TotalTSIPremiSpreadRNM (level polis)
├── FLAT_VIEW_SUGGEST         (1:N)   ViewSuggest (jejak komentar/approval)
├── FLAT_LOCATION             (1:N)   LocationList            [FIRE/ANEKA]
│   ├── FLAT_PROPERTY_ITEM    (1:N)   PropertyItemList        [FIRE]
│   ├── FLAT_OCCUPATION       (1:N)   OccupationList          [ANEKA/FIRE]
│   │   └── FLAT_ANEKA        (1:N)   AnekaList               [ANEKA]
│   ├── FLAT_CAUSE_OF_LOSS    (1:N)   ListCauseOfLoss         [FIRE]
│   └── FLAT_SCORING_*        (…)     ScoringRisk (FIRE)      [lihat §5]
├── FLAT_CARGO                (1:N)   CargoList (root)        [MARINE]
├── FLAT_COVERAGE             (1:N)   CoverageList  ← simpul konvergen semua COB
│   ├── FLAT_DEDUCTIBLE       (1:N)   DeductibleList
│   ├── FLAT_LAYER            (1:N)   LayerList
│   ├── FLAT_ADDITIONAL_COV   (1:N)   AdditionalCoverage      [MARINE]
│   └── FLAT_SPREADING        (1:N)   SpreadingList  ← baris spreading (≈ FACINPRODUCTION)
└── FLAT_PERSON               (1:N)   PersonList (Data-Party-Person)  [jika ada]
```

> `FLAT_SPREADING` adalah granularitas terhalus dan setara baris `FACINPRODUCTION` yang sudah ada
> (satu baris per spreading). `FACINPRODUCTION` tetap dipakai untuk pelaporan; `FLAT_SPREADING` adalah
> sumber transaksional yang memberinya makan.

---

## 3. Tabel level polis

Konvensi kolom umum di **semua** tabel anak (mewarisi dari header):
`IDPEGA VARCHAR2(50)`, `PRODKE VARCHAR2(5)`, `NOPOLIS VARCHAR2(100)`. Ini FK komposit ke
`JSON_POLIS_HEADER`. Kolom `*_CCY VARCHAR2(10)` menyertai tiap kolom uang (mata uang dari `Currency`
simpul terdekat). Semua nilai uang = `NUMBER` (bukan `BINARY_DOUBLE`).

### 3.1 `JSON_POLIS_HEADER` — pengganti langsung baris `JSON_POLIS`
Sumber: root object `OfferFacIn` (`pxObjClass=ASM-FW-GISFW-Data-OfferFacIn`) + kolom relasional
existing `JSON_POLIS`. [terverifikasi] field root dari kelima berkas.

| Kolom | Tipe | Asal (JSON key) | Catatan |
| --- | --- | --- | --- |
| IDPEGA | VARCHAR2(50) PK-bagian | (dari `JSON_POLIS.IDPEGA`) | PK = (IDPEGA) existing; lihat catatan versi |
| NOPOLIS | VARCHAR2(100) | `JSON_POLIS.NOPOLIS` | |
| NOENDORS | VARCHAR2(100) | `JSON_POLIS.NOENDORS` | |
| PRODKE | VARCHAR2(5) | `JSON_POLIS.PRODKE` | **versi polis** — riwayat = baris bertambah |
| OLDNOPOLIS | VARCHAR2(100) | `JSON_POLIS.OLDNOPOLIS` | |
| ID_WORK | VARCHAR2(60) | `ID` | mis. `ASM-FW-GISFW-WORK NB-xxxxx` |
| ID_NEW_BISNIS | VARCHAR2(30) | `IDNewBisnis` | |
| PROR_ATE_PERCENT | NUMBER | `ProRatePercent` | pecahan desimal (K-048), diport apa adanya |
| PRO_RATE_TYPE | VARCHAR2(5) | `ProRateType` | kode, arti belum dipetakan → gagal keras jika dipakai logika |
| IS_PRO_RATE | VARCHAR2(20) | `IsProRate` | `"Prorate"` dll |
| PERCENT_SHARE | NUMBER | `PercentShare` | share RNM (%) |
| CEDING_RETENTION | NUMBER | `CedingRetention` | |
| SHARE_CEDANT_TYPE | VARCHAR2(5) | `ShareCedantType` | kode |
| IS_FAC_RETRO | VARCHAR2(5) | `IsFacRetro` | |
| IS_FLAG_REJECT | VARCHAR2(5) | `IsFlagReject` | dipakai keputusan (bukan audit) |
| FLAG_SPECIAL_ACCEPT | VARCHAR2(5) | `FlagSpecialAcceptance` | |
| IS_DEDUCTIBLE_ACCEPT | VARCHAR2(5) | `IsDeductibleAcceptance` | |
| IS_OCCUP_EXCEPTION | VARCHAR2(5) | `IsOccupException` | |
| PARAM_CONDITION | VARCHAR2(30) | `ParamCondition` | `FIRE`/`ANEKA`/`PROPERTY`/… penanda COB di root |
| CONDITION | VARCHAR2(30) | `Condition` | (FIRE) |
| COMMENT_TXT | VARCHAR2(2000) | `Comment` | |
| DATE_OFFER | VARCHAR2(30) | `Date` | string GMT Pega — diport apa adanya |
| DATE_VALIDITY | VARCHAR2(30) | `DateValidity` | |
| DAYS_VALIDITY | NUMBER | `DaysValidity` | |
| CURRENT_YEAR | VARCHAR2(5) | `CurrentYear` | |
| TOTAL_PREMI_NUSARE | NUMBER | `TotalPremiNusaRe` | + TOTAL_PREMI_NUSARE_CCY |
| TOTAL_TSI_NUSARE | NUMBER | `TotalTSINusaRe` | + _CCY |
| TOTAL_TSI_NUSARE_SPREAD | NUMBER | `TotalTSINusaReSpreading` | + _CCY |
| TOTAL_TSI_TOP_RISK | NUMBER | `TotalTSITopRisk` | + _CCY |
| TREATY_CAPACITY | NUMBER | `TreatyCapacity` | (tak selalu ada) |
| IS_B2B | VARCHAR2(10) | `IsB2B` | (MARINE) |
| POLICY_MASTER_IDPEGA | VARCHAR2(60) | `PolicyMasterIDPega` | (MARINE, deklarasi) |
| POLICY_MASTER_NUMBER | VARCHAR2(60) | `PolicyMasterNumber` | (MARINE) |
| FOLLOWING | VARCHAR2(60) | `Following` | (FIRE banding) |
| FOLLOWING_NB | VARCHAR2(30) | `FollowingNB` | |
| IS_BANDING | VARCHAR2(5) | `IsBanding` | |
| POLICY_STATUS | VARCHAR2(30) | `PolicyStatus` | |
| TGL_INPUT / TGL_PROD / TGL_KONVERSI | DATE | `JSON_POLIS.*` | dipertahankan dari tabel existing |
| STS_KONVERSI | NUMBER | `JSON_POLIS.STS_KONVERSI` | |
| USERNAME | VARCHAR2(50) | `JSON_POLIS.USERNAME` | |

> **[pertanyaan terbuka] Kunci utama versi.** `JSON_POLIS` existing PK = `IDPEGA` saja, dengan UNIQUE
> `(NOPOLIS, IDPEGA)`. Karena PRODKE menandai versi dan riwayat = baris bertambah, tabel flat butuh
> PK efektif `(IDPEGA, PRODKE)` **atau** satu IDPEGA baru per versi. Mana yang dipakai bergantung pada
> bagaimana `INSERTJSONPOLIS` dipanggil untuk endorsement (satu IDPEGA di-append PRODKE baru, atau
> IDPEGA baru). **Belum terverifikasi dari korpus** — perlu contoh dua versi dengan IDPEGA sama.
> Jangan menebak. Ditandai untuk paket DBA.

### 3.2 `FLAT_QUOTATION` (1:1) — sumber `QuotationData` (`Data-Quotation`)
[terverifikasi] semua field dari NB/RNW/EDM/MARINE.

Kolom: BRANCH_CODE, BRANCH_NAME, BUSINESS_CODE, BUSINESS_FAC, BUSINESS_NAME, BUSINESS_OLD_ID,
BUSINESS_TYPE, BUSINESS_TYPE2, CEDING_CO, CEDING_CO_NAME, EDM_DAY (NUMBER), GROUP_NAME, GROUP_PANEL,
INSURED_ID, INSURED_NAME *(PII — kolom ada, tak diisi contoh)*, IS_GROUP, MARKETING_CODE,
MARKETING_NAME *(PII)*, MOID, NO_OFFER_SLIP, OPERATOR_ID *(PII)*, POLICY_TYPE, SHARE_OF_CEDING,
SOB_LEADER0, SOB_LSG, SOB_NAME, SOURCE_OF_BUSINESS, STATUS_BUSINESS, STATUS_SYARIAH, TEAM_GROUP,
TYPE_FACULTATIVE, EMAIL *(PII)*, BTN_QUOTATION.
> `CedingCoList` (anak QuotationData) → gunakan `FLAT_CEDING_CEDANT` (§3.4) atau tabel kecil
> `FLAT_QUOTATION_CEDINGCO` bila perlu dibedakan dari CedingCedantList level root.

### 3.3 `FLAT_POLICY` (1:1) — sumber `PolicyData` (`Data-Policy`)
Kolom: START_DATETIME, END_DATETIME, OFFERING_DATE, PROD_DATETIME, POLICY_NO *(PII: nomor polis —
kolom ada, tak diisi)*.
Sub-objek `PolicyData.Payment` → `FLAT_PAYMENT_HEADER`: INSTALLMENT (NUMBER), PCT_BROKERAGE_FEE,
RI_COMMISION.
Sub-objek `PolicyData.Ship` **[MARINE, terverifikasi]** → tabel `FLAT_POLICY_SHIP` (1:1 opsional):
INVOICE_NUMBER, SAIL_DATE, SHIP (nama kapal), DWT, GRT, NM_SHIP, NRT, REMARK, ADDITIONAL_SHIP,
SURVEY_AGENT. `pxObjClass=ASM-FW-GISFW-Int-MSHIP` untuk lookup kapal.

### 3.4 `FLAT_CEDING_CEDANT` (1:N) — `CedingCedantList` (`Data-Quotation`)
Kolom: SUBSCRIPT (dari `pxListSubscript`), CEDING_CO, CEDING_CO_NAME, SHARE_CEDING (NUMBER).
Anak `CurrencyList` di dalam tiap cedant → `FLAT_CEDING_CEDANT_CCY`: NAME, PREMIUM (NUMBER)+_CCY,
TSI (NUMBER)+_CCY.

### 3.5 `FLAT_CURRENCY_TOP` (1:N) — `CurrencyList` root (`OfferFacIn-Currency`)
Kolom: SUBSCRIPT, NAME (mata uang), OLD_ID, SUM_TOTAL_PAYMENT (NUMBER).
Sub `Policy.Payment` → kolom di baris ini: BROKERAGE_FEE, COMMISION, NET_PREMIUM, PPH, PPN, PREMIUM,
TSI_TOTAL (semua NUMBER).
Anak `ListInstallment` → `FLAT_INSTALLMENT` (§3.6).

### 3.6 `FLAT_INSTALLMENT` (1:N) — `ListInstallment` (`Data-Installment`)
FK ke FLAT_CURRENCY_TOP (SUBSCRIPT induk). Kolom: INSTALLMENT_NO, INSTALLMENT_PERCENTAGE (NUMBER),
DUE_DATE, PREMIUM, PAYMENT_TOTAL, PPH, PPN, BROKERAGE_FEE, RI_COMMISION, DISCOUNT, DEDUCTION2,
STAMP (semua NUMBER kecuali tanggal/urut).

### 3.7 `FLAT_TOTAL_TSI_PREMI` (1:N) — `TotalTSIPremiGrossList` + `TotalTSIList`
(`OfferFacIn-Currency`). Kolom: LIST_KIND ('GROSS'|'TSI'), NAME, PREMIUM (NUMBER)+_CCY, RATE (NUMBER),
TSI (NUMBER)+_CCY.

### 3.8 `FLAT_SPREAD_RNM` (1:N) — `TotalTSIPremiSpreadRNM` (`Data-SpreadingRisk`, level polis)
Kolom: CLAIM_ESTIMATION, PREMIUM_SPREADED, SHARE_PERCENTAGE, TREATY_NAME, TREATY_TYPE (kode),
TSI_SPREADED, CLAIM_AMOUNT_IDR, CLAIM_SPREADED (NUMBER). Bedakan dari `FLAT_SPREADING` (§4.7) yang
per-coverage.

### 3.9 `FLAT_VIEW_SUGGEST` (1:N) — `ViewSuggest` (`OfferFacIn-SuggestList`)
Jejak komentar/approval underwriting. Kolom: NO, SUBSCRIPT, APPROVAL, COMMENT_SUGGEST *(bisa memuat
nama — sanitasi)*, DATE_SUGGEST, DATE_TRANSFER, IS_CEDING_CONFIRM, IS_SAVE, PIC_SUGGEST *(PII)*.
> Ini bahan `HISTORYAKSEPTASIPEGA`-like; simpan bentuknya, dipakai untuk jejak keputusan.

---

## 4. Tabel risiko — per COB, konvergen di Coverage

### 4.1 `FLAT_LOCATION` (1:N) — `LocationList` (`OfferFacIn-LocationReinsurance`) [FIRE/ANEKA]
Kolom kunci alami: LOC_SUBSCRIPT (`pxListSubscript`). Kolom: PY_CITY, PY_COUNTRY_NAME,
LOSS_RATIO_1YEAR_AMOUNT/PERCENT, LOSS_RATIO_35YEAR_AMOUNT/PERCENT (NUMBER).
Sub-objek `Property` (`Data-Property`) diratakan ke kolom di baris ini: PROP_ALM_RISK_ID,
PROP_COUNTRY, PROP_PROVINCE, PROP_OBJECT_NO, PROP_OBJECT_NAME, PROP_OBJECT_TYPE, PROP_ROAD_NAME
*(alamat — PII)*, PROP_ROAD_TYPE, PROP_IS_TOP_RISK, PROP_OWNERSHIP, PROP_TOTAL_TSI (NUMBER),
PROP_IS_MATERIAL_DAMAGE, PROP_IS_FLAMMABLE, PROP_IS_HOTWORK, PROP_IS_PRODUCTION_PROCESS.
Sub `Property.BuildingConstruction` (`Data-BuildingConstruction`) → kolom: BC_FLOOR_TYPE,
BC_ROOF_TYPE (kode), BC_WALL_TYPE (kode).
Sub `RiskLocation` (`Data-Address`) → kolom: ADDR_ASM_ADDRESS *(PII)*, ADDR_ASM_CITY, ADDR_ASM_DISTRICT,
ADDR_ASM_RW, ADDR_ASM_ZIPCODE.
Sub `SurroundingRisk` (`Data-SurroundingRisk`) → kolom: SR_FLOOD_AREA_STATUS (kode),
SR_HOUSEKEEPING_STATUS (kode).

> **Penting [terverifikasi]:** di ANEKA (176005) `OccupationList` berada **di dalam**
> `Property.RiskLocation`, sedangkan di FIRE (181231) `OccupationList` berada langsung di bawah
> `Property`. Kolom FK `FLAT_OCCUPATION` cukup menunjuk `FLAT_LOCATION` (LOC_SUBSCRIPT) untuk keduanya;
> jenjang antara tidak perlu tabel sendiri karena tidak membawa data selain wadah.

### 4.2 `FLAT_PROPERTY_ITEM` (1:N) — `PropertyItemList` (`Data-PropertyItem`) [FIRE]
FK: LOC_SUBSCRIPT. Kunci alami: INDEX_PROPERTY, INDEX_PROPERTY_ITEM, PROPERTY_ITEM_NO,
ITEM_SUBSCRIPT. Kolom: ITEM_TYPE, ITEM_TYPE_ID (kode), CURRENCY, CURRENCY_ID, CURRENCY_OLD_ID,
FLAG_NET_RATE, IS_ADJUSTABLE_FLAG, PCT_ADJUST2, PCT_ADJUST_OTHER (NUMBER), TSI_OBJECT_ITEM (NUMBER),
TOTAL_GROSS_PREMI (NUMBER), TOTAL_NET_RATE (NUMBER), TOTAL_PREMIUM_NUSANTARA_RE (NUMBER).

### 4.3 `FLAT_OCCUPATION` (1:N) — `OccupationList` (`Data-Occupation`) [ANEKA/FIRE]
FK: LOC_SUBSCRIPT. Kunci alami: IDX_LOCATION, IDX_OCCUPATION, OCC_SUBSCRIPT.
Kolom: OCCUPATION_ID (kode), OCCUPATION_NAME, CATEGORY.
Sub `TableOfLimit` (`Int-TABLEOFLIMIT`) → kolom: TOL_CATEGORY, TOL_DESCRIPTION.

### 4.4 `FLAT_ANEKA` (1:N) — `AnekaList` (`Data-Aneka`) [ANEKA]
FK: LOC_SUBSCRIPT, IDX_OCCUPATION. Kunci alami: IDX_ANEKA, IDX_LOCATION, ANEKA_SUBSCRIPT.
Kolom: OBJECT_NAME, SECTION, SELECTED_LOCATION_ADDRESS *(PII)*, SELECTED_OBJECT_ITEM, TSI (NUMBER)+_CCY.

### 4.5 `FLAT_CARGO` (1:N) — `CargoList` root (`Data-Cargo`) [MARINE]
FK: header saja (CargoList di root, tidak di LocationList). Kunci alami: CARGO_SUBSCRIPT, INDEX_CARGO
(dipakai Coverage MARINE). Kolom: CONVEYANCE_ID (kode), CONVEYANCE_NOTE, GOOD_ID (kode), GOOD_NOTE,
PACKING_ID (kode), PACKING_NOTE, TRADING_ID (kode), TRADING_NOTE, FROM_RUTE, TO_RUTE, IS_TOP_RISK.

### 4.6 `FLAT_COVERAGE` (1:N) — `CoverageList` (`Data-Coverage`) — **simpul konvergen semua COB**
Kolom FK **kondisional per COB** (isi yang relevan, sisanya NULL):
- FIRE: LOC_SUBSCRIPT + INDEX_PROPERTY + INDEX_PROPERTY_ITEM
- ANEKA: LOC_SUBSCRIPT + IDX_OCCUPATION + IDX_ANEKA
- MARINE: INDEX_CARGO
Kunci alami umum: INDEX_COVERAGE, COV_SUBSCRIPT.
Kolom nilai [terverifikasi]: COVERAGE (kode), COVERAGE_NOTE, COVERAGE_BASIS (kode),
CALCULATE_METHOD_FACIN (kode, FIRE/ANEKA), DAY (NUMBER), RATE (NUMBER), RATE_OJK (NUMBER),
NET_RATE (NUMBER), FLAG_NET_RATE, PREMIUM (NUMBER)+_CCY, PREMI_NUSANTARA_RE (NUMBER)+_CCY,
TSI (NUMBER)+_CCY, TSI_LIABILITY (NUMBER), TSI_NUSANTARA_RE (NUMBER)+_CCY,
DISCOUNT (NUMBER), DISCOUNT_PERCENTAGE (NUMBER), INDEMNITY, INDEMNITY_PERCENTAGE, LOST_LIMIT,
LIMIT_OF_LIABILITY, SUBLIMIT, TSI_SUBLIMIT (NUMBER), PCT_ADJUSTMENT, PCT_LOL, PCT_SHORT_PERIOD,
LOADING, FIRST_LOSS, FIRST_LOSS_SCALE, FIRST_SCALE, EML_PML, MIN_PREMIUM (NUMBER),
MAX_STANDARD_RATE, MIN_STANDARD_RATE, UNIT, ZONE, ACCUMULATION_CODE, ACCUMULATION_DESCRIPTION,
RISK_ZIP_CODE, OLD_ID.
Sub `Currency` (`OfferFacIn-Currency`) → kolom: COV_CCY_ID, COV_CCY_NAME, COV_CCY_OLD_ID.

**Kolom before-image EDM [terverifikasi di MARINE, Seam 4]** — dipertahankan berpasangan:
IS_OLD_DATA (`IsOldData` = 'old'/…), PREMIUM_OLD (NUMBER), PREMI_NUSANTARA_RE_OLD (NUMBER),
TSI_OLD (NUMBER), PRO_RATE_PERCENT (NUMBER, pecahan — K-048).
> Ini pasangan nilai-berjalan vs nilai-lama, sejajar prinsip `_MENJADI`/`_SELISIH` §4.3. Untuk NB
> murni kolom `*_OLD` bernilai 0/NULL; untuk EDM terisi. Jangan diratakan.

### 4.7 `FLAT_DEDUCTIBLE` (1:N) — `DeductibleList` (`Data-Deductible`)
FK: COV_SUBSCRIPT (+ index coverage/property utk keunikan). Kunci alami: INDEX_DEDUCTIBLE +
index induk. Kolom: TYPE_DEDUCTIBLE (kode), TYPE_DEDUCTIBLE2 (kode), MIN_MAX (kode),
PCT_DEDUCTIBLE (NUMBER), PCT_DEDUCTIBLE2 (NUMBER), AMOUNT (NUMBER)+_CCY (`Currency` field),
CONDITION (kode), INPUT_CONDITION, FLAG_CURRENCY.

### 4.8 `FLAT_LAYER` (1:N) — `LayerList` (kosong di semua contoh)
[terverifikasi struktur wadahnya ada, **isi kosong** di kelima kasus] → **[pertanyaan terbuka]**
skema kolom Layer belum dapat dipetakan dari korpus contoh ini. Tabel disiapkan (FK COV_SUBSCRIPT +
LAYER_SUBSCRIPT) tetapi definisi kolom **menunggu contoh kasus berlayer**. Jangan menebak kolom.

### 4.9 `FLAT_ADDITIONAL_COV` (1:N) — `AdditionalCoverage` (kosong di contoh) [MARINE]
[pertanyaan terbuka] sama seperti Layer — wadah ada, isi kosong. Siapkan FK COV_SUBSCRIPT, kolom
menunggu contoh.

### 4.10 `FLAT_SPREADING` (1:N) — `SpreadingList` (`Data-SpreadingRisk`) — granularitas terhalus
FK: COV_SUBSCRIPT (+ index induk). Kunci alami: SPREAD_SUBSCRIPT. Kolom: TREATY_NAME, TREATY_TYPE
(kode), SHARE_PERCENTAGE (NUMBER), PREMIUM_SPREADED (NUMBER)+_CCY, TSI_SPREADED (NUMBER),
TSI_GROSS_SPREADED (NUMBER), CLAIM_ESTIMATION (NUMBER).
> Setara baris `FACINPRODUCTION` existing (satu baris per spreading). `FACINPRODUCTION` tetap sebagai
> tabel pelaporan/produksi; `FLAT_SPREADING` sumber transaksionalnya.

### 4.11 `FLAT_CAUSE_OF_LOSS` (1:N) — `ListCauseOfLoss` (`Data-CauseOfLoss`) [FIRE]
FK: LOC_SUBSCRIPT. Kolom: CLAIM (NUMBER)+CURRENCY, DETAIL, REMARKS. Sub `CoinsData` (`Data-Coins`)
→ COINS_NAME.

### 4.12 `FLAT_PERSON` (1:N) — `PersonList` (`Data-Party-Person`) [kosong di contoh]
[pertanyaan terbuka] wadah ada, kosong di kelima kasus → kolom person (nama, dsb — **PII**) menunggu
contoh kasus Life/PA. Jangan menebak. Untuk Life lihat §6 (`InputEDMLife`, di luar cakupan contoh ini).

---

## 5. Scoring risiko (FIRE) — struktur dalam, banyak sub-objek seragam

Sumber: `ScoringRisk` (`Data-ScoringRisk`) → `DataScoringRiskList` (`Data-ScoringRisk`) →
`FEA`/`LossRatio`/`ObjectConditions`/`Occupation`/`Others`/`RiskImprovement`/`RiskLocation`/
`SurveyReport`, masing-masing memuat banyak sub-objek `Data-RiskDetail` dengan pola seragam
(`ChechBox1..N`, `Score1..N`, `ScorePerFactor`, `Checked`). [terverifikasi di FIRE 181231].

**Usulan [dugaan, perlu konfirmasi bisnis]:** karena polanya seragam (faktor → skor), ratakan jadi
tabel EAV alih-alih ratusan kolom:

- `FLAT_SCORING_HEADER` (1:1 per DataScoringRisk): FK header, SCORING_SUBSCRIPT, FINAL_SCORE (NUMBER),
  NOTE_FINAL_SCORE, STATUS, OCCUPATION_CODE, OCCUPATION_NOTE, TOP_RISK_LOCATION *(PII: alamat)*.
- `FLAT_SCORING_FACTOR` (1:N): FK SCORING_SUBSCRIPT, GROUP_NAME (mis. 'FEA','LossRatio','Others'),
  FACTOR_NAME (mis. 'FireAlarmSystem'), CHECKED, CHECKBOX_IDX (NUMBER), SCORE_VALUE (NUMBER),
  SCORE_PER_FACTOR (NUMBER), REMARKS.

> **[pertanyaan terbuka]** Apakah bisnis ingin skema skoring dipertahankan per-checkbox (audit penuh)
> atau cukup skor akhir per faktor. Bentuk EAV mempertahankan segala nilai tanpa menebak arti tiap
> checkbox — aman untuk migrasi. Arti tiap `ChechBox`/`Score` **tidak ditebak**. Keputusan bentuk
> final milik bisnis + berkaitan dengan Seam 6 (ScoreMedical) untuk jalur Life.

---

## 6. Batas — apa yang TIDAK dirancang di sini

- **Jalur Life / `PersonList` terisi / medical** — tidak ada di kelima contoh (semua NONLIFE, PersonList
  kosong). Skema tabel person/medis menunggu contoh kasus Life. Flow `InputEDMLife` (K-052) di luar
  cakupan dokumen ini. **[pertanyaan terbuka]**
- **`LayerList` dan `AdditionalCoverage`** — wadah terverifikasi ada tetapi kosong; kolom menunggu
  contoh. **[pertanyaan terbuka]**
- **Arti semua kode enumerasi** (`TreatyType`, `CalculateMethod_FacIn`, `CoverageBasis`,
  `TypeDeductible`, `ConveyanceID`, `GoodID`, `PackingID`, `TradingID`, `SR_FLOOD_AREA_STATUS`, dst.)
  — disimpan apa adanya, penerjemahan menunggu tabel lookup/keputusan bisnis. **[§4.5 gagal keras
  bila dipakai untuk logika tanpa arti].**
- **Prosedur transformasi write (32 prosedur ALL_SOURCE)** — bagaimana `INSERTJSONPOLIS` diganti
  menjadi insert ke tabel-tabel flat ini, dan bagaimana `SaveJsonPolicyFacIn_Act`/`SaveEDMToJsonPolicy_Act`
  memetakan agregat → baris flat, **belum dapat diverifikasi** tanpa dump ALL_SOURCE dari DBA.
  Ini blocker E21 (jalur produksi). **[pertanyaan terbuka — paket DBA]**
- **Kunci versi (IDPEGA vs IDPEGA+PRODKE)** — lihat §3.1, perlu contoh dua versi. **[paket DBA]**

---

## 7. Ringkasan untuk keputusan work owner

**Yang [terverifikasi] dan siap:**
- 1 header + ~20 tabel anak melayani ketiga siklus (NB/RNW/EDM) karena strukturnya sama.
- `CoverageList` + `SpreadingList` adalah simpul konvergen semua COB → mesin premi/spreading yang
  sudah dirancang (Seam 3) memakannya langsung.
- Field before-image EDM (`*Old`) sudah berada di level Coverage → cocok dengan Seam 4.
- Uang selalu berpasangan dengan mata uang; PRODKE tetap versi; riwayat = baris bertambah.

**Yang menunggu (jangan ditebak):**
1. Kunci versi IDPEGA/PRODKE — butuh 2 versi contoh.
2. Skema `LayerList`, `AdditionalCoverage`, `PersonList`/Life — butuh contoh berisi.
3. Bentuk final tabel scoring (per-checkbox vs ringkas) — keputusan bisnis.
4. Pemetaan write (32 prosedur) → butuh dump ALL_SOURCE dari DBA (blocker E21).
5. Arti kode enumerasi → tabel lookup dari DBA/bisnis.

Poin 1, 4, 5 → masuk `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md`.
