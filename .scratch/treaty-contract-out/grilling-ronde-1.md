# Grilling Ronde 1 — Treaty Contract Out

Konteks: `treaty-contract-out` — **konteks non-Life pertama**. Modul `Treaty Contract Out` (303 berkas).
Tanggal: 2026-09-15
Skill: `/mattpocock-skills:grilling`
Sumber: korpus `D:\XML\RNM_BRD\Treaty Contract Out\` (READ-ONLY), `discovery/modules/Treaty Contract Out.md`,
`discovery/flows/Treaty Contract Out.md`, `discovery/context-map.md` §2.4 & §3, `CONTEXT.md`,
artefak `.scratch/master-contract-retro-life/`, `.scratch/master-product-name-life/`

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[fakta bisnis — work owner]`; `[data DBA]`; `[terbuka]` = OQ.
> **Identitas rule wajib menyertakan class** — nama sama di class berbeda = rule berbeda.

---

## 0. Yang sudah mengendap (tidak ditanyakan ulang)

Fondasi dari prompt work owner + D2/D3, tidak diuji ulang: peran modul (editor master
term/arrangement, **bukan** treaty outward), titik masuk `Harness/InboxTreatyContract.xml`
(`DATA-PORTAL!INBOXTREATYCONTRACT`), bentuk editor master tanpa tangga persetujuan (nol `Flow`,
nol `When`, nol `StatusAkseptasi`), distribusi rule, dan 7 objek Oracle inti.

---

## 1. Temuan baru Ronde 1 — hasil telusur isi rule

### 1.1 ⚠️ Temuan utama: **tabel kembar** — JSON `M_*` berdampingan dengan relasional non-`M_*`

`[terverifikasi]` Modul ini membaca **dua keluarga tabel yang menyimpan hal yang sama**.

| Entitas logis | Tabel **JSON** (dibaca `a.JSONDATA.<field>`) | Tabel **relasional** (kolom datar) |
| --- | --- | --- |
| Tahun treaty | `M_TREATYYEAR` | — (kelas Pega `ASM-FW-GISFW-INT-TREATYYEAR` memetakan tabel yang tak tampak di korpus) |
| Kontrak treaty | `M_TREATYCONTRACT` | `TREATYCONTRACT` (hanya dipakai saat **hapus**) |
| Business | `M_TREATYBUSINESS` (hanya saat **hapus**) | `TREATYBUSINESS` (baca + hapus) |
| Arrangement (klausul) | `M_PROPORTIONALARRG` (**10 rujukan**) | `PROPORTIONALARRG` (**1 rujukan**) |
| Reinsurer | — | `TREATYREINSURER` |
| Security reinsurer | — | `MTREATYSECURITY` (ejaan tanpa garis bawah) |
| Kurs | — | `TREATYEXCHANGE` |

Bukti paling tajam — **dua rule menjawab pertanyaan yang sama dari dua tabel berbeda**:

- `RDBList/GetMasterDescriptionLimitParentList.xml`
  (`ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL`):
  `SELECT a.TreatyYear, a.TreatyGroupID, a.TreatyDescID, a.ReinsTypeID, a.ParentReinsTypeID, a.Pct
  FROM PROPORTIONALARRG a WHERE a.TreatyYearID=… AND a.TreatyGroupID=… AND a.TreatyDescID=… AND
  a.ParentReinsTypeID=…` → **relasional**.
- `RDBList/GetMasterDescriptionEPIParentList.xml`: kueri **identik secara semantik**, tetapi
  `FROM m_PROPORTIONALARRG a WHERE a.JSONDATA.TreatyYear=… AND a.JSONDATA.TreatyGroupID=… AND
  a.JSONDATA.TreatyDescID=… AND a.JSONDATA.ParentReinsTypeID=…` → **JSON**.

Dan satu rule memegang **kedua-duanya sekaligus** — `RDBList/DeleteRowBusinessList.xml`
(`ASM-FW-GISFW-INT-TREATYBUSINESS` / `ASM!DELETEROWBUSINESSLIST` / `RULE-CONNECT-SQL`):

```
<pyBrowseSQL>  delete from treatybusiness   where id = {InputData.HASIL12}
<pyDeleteSQL>  delete from m_treatybusiness where id = {InputData.HASIL12}; commit;
```

⚠️ Satu rule, dua tabel, dua ejaan — **`treatybusiness` dan `m_treatybusiness`**.

`[terbuka]` Yang **tidak** dapat dijawab dari korpus: apakah tabel relasional adalah **turunan**
yang di-flatten oleh stored procedure dari JSON, atau sebaliknya, atau keduanya ditulis
berdampingan. Seluruh jalur tulis melewati procedure yang bodinya tidak ada → **OQ-002**.

Catatan pembanding: tanda tangan procedure di sini **berbeda** dari pola Life.
`PEGA_TREATYCONTRACT` menerima **parameter skalar berposisi**, bukan satu CLOB JSON —
berbeda dari `PEGA_M_PRODUCT_LIFE`/`PEGA_M_PRODUCT_INWARD_LIFE` (Master Product Name Life) yang
menerima `DATAPEGA CLOB`. Satu-satunya procedure di modul ini yang menerima CLOB adalah
`PEGA_M_ATTACHMENT` (lampiran).

### 1.2 ⚠️ Temuan utama: **25 jenis klausul menulis ke DUA procedure saja**

`[terverifikasi]` Sensus penuh — setiap `SaveTreatyArr*` dipetakan ke rule Connect-SQL yang
dipanggilnya (`grep` atas 25 berkas di `Treaty Contract Out/Activity/`):

| Menulis lewat | Jumlah | Jenis klausul |
| --- | ---: | --- |
| `SaveMasterProportionalArrg` → `POOLDATA.PEGA_PROPORTIONALARRG` (**induk**) | **18** | BordereAux, CashLossLimit, ClaimCoorp, CoinsPanel, EPI, ExGratia, ExclutionTreaty, FacIn, LimitMB, MaxCoinsPanel, MinLOL, MinLOLMB, PLA, Portfolio, ProfitComm, Ricomm, TerrLimit, TreatyLimit |
| `SaveMasterProportionalArrgChild` → `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD` (**anak**) | **7** | CashLossLimitList, ClaimCoorpChild, EpiList, ExGratiaChildList, FacInList, PLAList, TreatyLimitChild |

**Mekanismenya** `[terverifikasi]`: tiap klausul punya halaman masukannya sendiri
(`InputTreatyArrEpi`, `InputTreatyArrPLA`, `InputTreatyArrProfit`, `InputTreatyArrBord`, …), lalu
langkah **`Page-Copy`** menyalinnya ke halaman bersama **`InputTreatyArrTreatyLimit`** (induk) atau
**`InputTreatyArrTreatyLimitChild`** (anak) sebelum langkah `RDB-List` memanggil procedure.
Contoh: `Activity/SaveTreatyArrEPI_Act.xml`
(`ASM-FW-GISFW-INT-PROPORTIONALARRG` / `SAVETREATYARREPI_ACT` / `RULE-OBJ-ACTIVITY`, 111.418 byte,
10 langkah) — 31 rujukan `InputTreatyArrEpi`, 3 rujukan `InputTreatyArrTreatyLimit`, langkah 7 =
`Page-Copy`, langkah 8 = `RDB-List`.

**Akibatnya: `M_PROPORTIONALARRG` adalah SATU tabel generik** yang memuat seluruh jenis klausul,
dibedakan oleh kolom **`TreatyDescID` / `TreatyDescName`**. Bukan 25 tabel.

⚠️ Koreksi terhadap D2/D3: dokumen itu menulis **27** `SaveTreatyArr*_Act`; hitungan ulang
(`ls | grep -c "^SaveTreatyArr"`) = **25**. Dan daftar 16 `CancelActivity*` **bukan** daftar jenis
klausul yang lengkap — ia hanya 16 dari 25.

### 1.3 `[terverifikasi]` Bentuk baris arrangement — 36 parameter induk, 26 anak

`RDBList/SaveMasterProportionalArrg.xml` → `POOLDATA.PEGA_PROPORTIONALARRG` (36 in + 2 out):

```
ID, TreatyYear, TreatyYearID, TreatyGroupID, TreatyGroupName, TreatyDescID, TreatyDescName,
ReinsTypeID, ReinsTypeName, Layer, LayerPart, LayerPartType, LayerType, Kurs, TglUpdate, UserID,
Line, Pct, PctMe, Ydcf, Method, TerritorialLimit, ParentReinsTypeID, SpreadingOrder, Rp, Usd,
ID_Occupation, Occupation, ID_Clause, Clause, TreatyLimit, CoIns_Min, CoIns_Max, MoreRp, MoreUsd
→ OUT: HASIL1, HASIL2 ;  COMMIT internal
```

`RDBList/SaveMasterProportionalArrgChild.xml` → `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD` (26 in + 2 out):
sama persis **minus** `ID_Occupation, Occupation, ID_Clause, Clause, TreatyLimit, CoIns_Min,
CoIns_Max, MoreRp, MoreUsd` → OUT: `HASIL1`, **`BRANCH_CODE`**; `COMMIT` internal.

⚠️ Penamaan parameter keluaran **tidak konsisten antar jalur tulis**:
`HASIL1`/`HASIL2` (contract, year, arrangement induk) · `HASIL5`/`HASIL6` (business) ·
`ERRMSG`/`STSSAVE` (reinsurer, copy, attachment) · `HASIL1`/**`BRANCH_CODE`** (arrangement anak).
Arti tiap keluaran tidak ada di korpus → **OQ-002**.

`[terverifikasi]` JSON `M_PROPORTIONALARRG` memuat **daftar bersarang** — dibaca dengan `JSON_TABLE`
di `RDBList/GetMasterPortfolioListDetail.xml`:
`'$.ProportionalList[*]'` dengan kolom `PortfolioType`, `PremiLost`, `PortfolioValue`,
`PortfolioDesc`. Jadi Portfolio adalah **daftar di dalam dokumen**, bukan tabel anak.

### 1.4 `[terverifikasi]` Validasi wajib-isi **tidak seragam** antar klausul

Dari prasyarat langkah (`<pyStepsPreCondParamsWhen> … ==""`) pada 25 activity simpan:

| Klausul | Field wajib |
| --- | --- |
| CashLossLimit, ClaimCoorp, EPI, ExGratia, FacIn, PLA, TreatyLimit | `ReinsTypeID`, `Rp`, `Usd` |
| CashLossLimitList, ClaimCoorpChild, EpiList, ExGratiaChildList, FacInList, PLAList, TreatyLimitChild | `Pct`, `ReinsTypeID`, `Rp`, `Usd` |
| ProfitComm | `Pct`, `PctMe`, `ReinsTypeID` |
| Ricomm | `Method`, `Pct`, `ReinsTypeID` |
| CoinsPanel | `TerritorialLimit`, `TreatyLimit` |
| MaxCoinsPanel | `CoIns_Max`, `TerritorialLimit` |
| MinLOL, MinLOLMB | `Pct`, `TerritorialLimit` |
| TerrLimit | `TerritorialLimit` |
| ExclutionTreaty | `ID_Occupation`, `TerritorialLimit` |
| BordereAux | `Method` |
| **LimitMB** | **nol validasi** |
| **Portfolio** | **nol validasi** |

`[terverifikasi]` Penanda baris baru: `InputTreatyArr….ID = "UnknownId"`; anti-dobel memakai pesan
`OutputParam.ERRMSG4 = "Data sudah pernah di Input"` (`Activity/SaveTreatyArrEPI_Act.xml`).

### 1.5 ⚠️ Temuan: rule bernama **`_Old`** justru yang **terbaru dan terpakai**

`[terverifikasi]` Dua Report Definition di class yang sama:

| Berkas | `pxInsName` | RuleSet | Commit | Dipakai |
| --- | --- | --- | --- | ---: |
| `ReportDefinition/BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE!BROWSEREINSURANCETYPE_RD` | 01-01-83 | 2025-06-02 | **1** section |
| `ReportDefinition/BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE!BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` | **01-01-91** | **2026-02-05** | **11** section |

Nama berkata "Old"; ruleset, tanggal commit, dan jumlah pemakaian berkata **yang terbaru**.
(`pzOriginalInstanceKey`-nya menunjuk `BROWSEREINSURANCETYPE_RD #20170420…` — ia hasil *save-as*
dari versi 2017, lalu terus dirawat sendiri.) **Pelajaran OQ-066 terulang di modul ini.**

Perbedaan perilaku keduanya:

- **Yang dipakai 11 grid** (`_Old_Ljt_id_isnotnull`): `A AND B AND C` =
  `.ID NotStartsWith ("10004","10011","10012","10021","10022","10025","10026","10028","10248",
  "10249","10018","10217")` **AND** `.Flag = "active"` **AND** `.Type = ("1","2","3")`.
  → **12 ID jenis reasuransi di-blacklist secara literal.**
- **Yang dipakai 1 layar** (`BrowseReinsuranceType_RD`, dipakai `Section/InputTreatyContractReinsType.xml`
  dengan `Param.Flag="active"`): `(A OR B) AND C AND D` =
  `(.ID = Param.ID OR .Note = Param.Note) AND .Flag = Param.Flag AND .Type = Param.Type` —
  **tanpa blacklist**, seluruhnya berparameter.

`[terverifikasi]` Class `ASM-FW-GISFW-INT-REINSURANCETYPE` adalah **class yang sama** dengan master
jenis reasuransi Life (`Master Contract Retro Life/ReportDefinition/BrowseReinsuranceTypeLimit_RD.xml`,
nilai 10196–10200, filter `.Flag = 1`, nama dari `.Note`). Jadi **jenis reasuransi adalah master
bersama life & non-life**, dibedakan oleh `.Type` dan `.Flag`.
⚠️ Nilai `.Flag` pun berbeda ejaan antar konteks: **`1`** (Life) vs **`"active"`** (di sini).

### 1.6 `[terverifikasi]` Master jenis klausul: class `ASM-FW-GISFW-INT-TREATYDESC`

`ReportDefinition/BrowseTreatyDesc_RD.xml` (`…-INT-TREATYDESC` / `BROWSETREATYDESC_RD`):
field `.ID`, `.DescName`, **`.IsXOL`**, **`.StatusAktif`**; filter tunggal `.IsXOL = Param.IsXOL`.

`TreatyDescID` **tidak pernah literal** di modul ini — selalu datang dari parameter
(`TempArrg.TreatyDescID = Param.ID`, `TempArrg.TreatyDescName = Param.DescName`,
`InputData.CARIDESC = Param.TreatyDescID`). Jadi daftar jenis klausul adalah **data**, bukan kode —
walaupun Pega menyediakan 25 activity yang **di-hardcode per jenis**.

### 1.7 `[terverifikasi]` Hierarki kunci — **beda dari Retro Life**

Dari tanda tangan procedure penulis:

```
M_TREATYYEAR      (ID, TreatyYear, UnderwritingYear, TreatyGroupID, TreatyGroupName,
                   Proportion, StartDate, EndDate, UserID, TglUpdate)
   └─ M_TREATYCONTRACT (ID, IDTreatyYear → M_TREATYYEAR.ID, ReinsTypeID, ReinsTypeName,
                        TreatyStartDate, TreatyEndDate)

TREATYREINSURER   kunci alami = (TreatyYear, TreatyGroupID, ReinsTypeID)  ← BUKAN ID kontrak
   └─ MTREATYSECURITY  (REAS_ID → TREATYREINSURER.ID, REAS_SECURITY, PCT_SHARE, THN_TREATY)
TREATYBUSINESS    kunci alami = (TreatyYear, TreatyYearID, TreatyGroupID, ReinsTypeID) + BizCode
M_PROPORTIONALARRG kunci alami = (TreatyYear, TreatyYearID, TreatyGroupID, TreatyDescID, ReinsTypeID)
   └─ M_PROPORTIONALARRG_CHILD  (kunci sama, + ParentReinsTypeID)
```

⚠️ **Perbedaan mendasar dari Master Contract Retro Life**: di Life, anak-anak menggantung pada
**ID kontrak**. Di sini mereka menggantung pada **kunci gabungan (tahun + grup + jenis reasuransi)** —
**tanpa foreign key ke `ID` kontrak**. Itulah sebabnya `DeleteFromTREATYCONTRACT_SQL` harus menghapus
dengan mengulang seluruh kunci gabungan, bukan dengan satu `ID`.

### 1.8 ⚠️ `[terverifikasi]` Kaskade hapus ditulis tangan — dan **arrangement tidak ikut terhapus**

`RDBList/DeleteFromTREATYCONTRACT_SQL.xml` — SQL mentah, empat `DELETE` berurut, `COMMIT` di akhir:

1. `DELETE FROM treatycontract WHERE id=… AND IDTREATYYEAR=…`
2. `DELETE FROM treatybusiness WHERE TREATYYEAR=… AND (TREATYYEARID=… OR TREATYYEARID IS NULL) AND TREATYGROUPID=… AND REINSTYPEID=…`
3. `DELETE FROM MTREATYSECURITY WHERE REAS_ID IN (SELECT id FROM TREATYREINSURER a WHERE …)`
4. `DELETE FROM TREATYREINSURER WHERE TreatyYear=… AND TreatyGroupID=… AND ReinsTypeID=…`

⚠️ **`M_PROPORTIONALARRG` dan `M_PROPORTIONALARRG_CHILD` TIDAK dihapus** — seluruh klausul
kontrak yang dihapus **tertinggal sebagai baris yatim**.
⚠️ Juga: langkah 1 menghapus **`treatycontract`** (relasional) tetapi **bukan `M_TREATYCONTRACT`**
(JSON) — sedangkan `DeleteRowBusinessList` menghapus **keduanya**. Perlakuan tidak konsisten.

`RDBList/DeleteFromTreatyReinsurer_Act.xml` menghapus `MTREATYSECURITY` lalu `TREATYREINSURER`
**tanpa `COMMIT`**.

### 1.9 ⚠️ `[terverifikasi]` `MTREATYSECURITY` — satu-satunya master yang ditulis **SQL mentah**

```
RDBList/InsertToMTreatySecurity.xml   (ASM-FW-GISFW-INT-MTREATYSECURITY / ASM!INSERTTOMTREATYSECURITY)
  insert into mtreatysecurity values ({InputData.CARI8}, '', '', {InputData.CARI9},
                                      {InputTreatySecurity.PCT_SHARE}, '', {InputTreatySecurity.REAS_SECURITY})
RDBList/UpdateMTreatySecurity.xml
  update mtreatysecurity set THN_TREATY=…, PCT_SHARE=…, REAS_SECURITY=…
  where REAS_ID = … and trim(REAS_SECURITY) = trim(…)
RDBList/DeleteSecurityReinsurer.xml
  delete from mtreatysecurity where REAS_ID = … and trim(REAS_SECURITY) = trim(…)
```

⚠️ `INSERT` **tanpa daftar kolom** — 7 nilai berposisi, **3 di antaranya string kosong `''`**.
Kolom apa yang dikosongkan **tidak dapat diketahui dari korpus** → butuh DDL `[data DBA]`.
⚠️ Ketiganya **tanpa `COMMIT`**.
⚠️ Kunci pembaruan/penghapusan memakai `trim(REAS_SECURITY)` — artinya **nama security dipakai
sebagai bagian kunci**, dan datanya diketahui mengandung spasi berlebih.

### 1.10 `[terverifikasi]` `PROSESCOPY` — salin satu tahun treaty ke tahun lain

`RDBList/SaveMasterCopyData_SQL.xml`:

```
vTHN_TREATY := {InputData.CARI17}   -- tahun ASAL
vTOP_ID     := {InputData.CARI18}   -- TreatyGroupID
vTHN_TUJUAN := {InputData.CARI16}   -- tahun TUJUAN
vUserId     := {InputData.CARI19}
POOLDATA.PROSESCOPY(vTHN_TREATY, vTOP_ID, vTHN_TUJUAN, {OutputParam.ERRMSG out}, {OutputParam.STSSAVE out}, vUserId);
dbms_output.put_line(errmsg);
```

Pemanggil `[terverifikasi]`: `Activity/BrowseCopyData.xml`, `Activity/SaveTreatyYearMultiple_Act.xml`,
`Activity/BrowseDeleteRowTreatyInContract.xml`.

⚠️ **Tanpa `COMMIT`** — satu-satunya jalur tulis master yang tidak commit sendiri.
⚠️ `dbms_output.put_line(errmsg)` memakai variabel lokal `errmsg` yang **dideklarasikan tetapi
tidak pernah diisi** (keluaran procedure masuk ke `{OutputParam.ERRMSG out}`, bukan ke `errmsg`).
Jadi baris itu **selalu mencetak NULL** — kode mati.
⚠️ **Slot `InputData.CARI19` dipakai untuk dua arti berbeda**: `userId` pada jalur salin,
`ReinsTypeID` pada `DeleteFromTREATYCONTRACT_SQL`. Jebakan bagi migrasi yang meniru nama slot.

`[terbuka]` **Apa persisnya yang disalin** (kontrak saja? beserta reinsurer, business, seluruh
klausul?), dan **apa yang terjadi bila tahun tujuan sudah berisi** (timpa? tolak? gandakan?) —
seluruhnya di dalam body procedure → **OQ-002**.

### 1.11 `[terverifikasi]` Kurs — satu mata uang, di-hardcode

`RDBList/GetMasterKursList.xml`:

```
select TOIDR as HASIL1 from treatyexchange
 where Quarter='0'
   and to_date({InputData.CARI1},'YYYYMMDD') BETWEEN trunc(TO_TIMESTAMP_TZ(STARTDATE,…))
                                                 AND trunc(TO_TIMESTAMP_TZ(ENDDATE,…))
   and IDCURRENCY = '10001'
```

⚠️ **`IDCURRENCY='10001'` dan `Quarter='0'` literal.** Hasilnya masuk ke
`InputTreatyArrangement.Kurs` (`Activity/SetTreatyArrangementDesc_Act.xml`).
⚠️ `STARTDATE`/`ENDDATE` disimpan sebagai **teks** ber-format Pega
(`YYYYMMDD"T"HH24MISS.FF3 TZR`), lalu di-parse tiap kueri.
⚠️ Baris arrangement menyimpan **`Rp` dan `Usd` sebagai dua kolom terpisah** (+ `MoreRp`,
`MoreUsd`) — bukan satu nilai + kode mata uang.

### 1.12 `[terverifikasi]` Lampiran — pola sama dengan konteks lain, **tetapi kuncinya pinjaman**

Rantai berkas identik pola Master Product Name Life: `ConnectREST/ServiceGoogle.xml`
(`ASM-FW-GISFW-INT-T_STORAGE_IMAGE!SERVICEGOOGLE`, tanpa URL literal; alamat dari
`SystemSettings/LinkService.xml` → `LINKSERVICE!LINKSERVICE`), `GET_TOKEN_STORAGE`,
`T_STORAGE_IMAGE` (`GetLinkStorage_SQL`, `Update_T_Storage_SQL`, `DeleteStorage_SQL`),
`PEGA_M_ATTACHMENT` (CLOB JSON), kategori dari `CATEGORY_ATTACH_REAS`.

⚠️ `pzOriginalInstanceKey` `ServiceGoogle` = `…GOOGLESTORAGE_UPLOAD #20250729…` — sekali lagi
hasil *save-as*; nama berkas ≠ nama asal.

⚠️ **Kebocoran batas:** `RDBList/GetAllAttachment2_Sql.xml`
(class **`ASM-FW-GISFW-INT-TREATY_IN`**) membaca `M_ATTACHMENTTREATY_2 WHERE treatyid = {TreatyIn.ID}`.
Lampiran modul ini dikunci ke **ID treaty inward**, bukan ke ID kontrak/arrangement.
Empat rule lain berclass sama: `DELETE_ACT`, `DELETEATTACHMENT2_SQL`, `GETATTACHMENT2_SQL`,
`INSERTATATCHMENT_SQL`.

### 1.13 `[terverifikasi]` Kode mati — cukup banyak untuk perlu disensus

| Penanda | Berkas |
| --- | ---: |
| `1=2` | 33 |
| `1==2` | 1 |
| visibilitas `never` | 41 |
| langkah di-remark (`<pyStepsBlockName>//`) | 46 langkah |
| `<pyStepsBlockName>EXIT` | 2 langkah |
| memo "not used" | 1 |

⚠️ Sesuai pelajaran OQ-066 — **berkas ada ≠ dipakai**, dan **memo ≠ kebenaran**. Sensus hidup/mati
harus dikonfirmasi ke work owner sebelum ditiketkan.

---

## 2. Frontier Ronde 1

Dua belas pertanyaan di bawah ini prasyaratnya sudah terpenuhi. Yang **tidak** ditanyakan sekarang
karena bergantung pada jawaban di ronde ini:

- **jalur tulis** (panggil 9 procedure apa adanya vs tulis sendiri) — bergantung Q1
- **batas transaksi & atomisitas** — bergantung Q1 + Q2
- **bentuk migrasi data** — bergantung Q1 + Q2 + Q4
- **irisan tiket** — bergantung seluruh ronde ini

(Isi pertanyaan lengkap ada di pesan grilling; berkas ini menyimpan buktinya.)
