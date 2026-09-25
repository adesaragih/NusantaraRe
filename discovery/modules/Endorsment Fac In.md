# Modul — Endorsment Fac In

STEP D3. Sintesis dari `../inventory/Endorsment Fac In.md` (D1) +
`../flows/Endorsment Fac In.md` (D2 Tahap 4) + `../flows/_SUMMARY-facultative.md`.
**2.061 file** — modul terbesar kedua di korpus.
Audit: `find "Endorsment Fac In" -name '*.xml' | wc -l`

> Ejaan nama modul (`Endorsment`, tanpa `e`) **ditulis apa adanya** sesuai korpus.

## 1. Peran modul

`[terverifikasi]` Menangani **endorsement / addendum facultative inward** — perubahan atas
penutupan yang sudah berjalan. Class kerja `ASM-FW-GISFW-WORK` (531 rule), sama dengan NB FacIn dan
RNW Fac In.

`[terverifikasi]` Bagian dari **satu ruleset facultative berdiskriminator siklus**: `When/IsNB.xml`,
`IsRenewal.xml`, `IsEDM.xml`, `IsNotEDM.xml` ada di ketiga modul → OQ-015.

## 2. Proses / fitur utama

Rincian: **`../flows/Endorsment Fac In.md`** (52 rule ditelusur — terbanyak di D2).
**Dua rule `Flow`** (paling sedikit dari trio facultative):

| File | Start | Hash | Catatan |
| --- | --- | --- | --- |
| `Flow/InputAddendumFacultativeIn.xml` | **`Start2`** | `f0b2f4ea86` | khas modul ini |
| `Flow/OfferFacRetro.xml` | `Start1` | **`b370146c63`** | **VARIAN BERBEDA** dari NB/RNW (`bfd6252070`) — OQ-011 #338 |

`[terverifikasi]` Modul ini **tidak punya** `InputInwardFacultativeRISlip`, `OfferFacOut`,
`InputQuotation`, `InputInwardFacultativeOffer`, maupun `InputRealizationTreatyIn`.

`[terverifikasi]` `InputAddendumFacultativeIn`: 15 Assignment, 37 Decision, 9 Utility,
3 SubProcess, 133 connector.

### 2.1 Gerbang khas endorsement

`[terverifikasi]` Enam shape Decision tanpa padanan di NB FacIn / RNW Fac In:
`Decision23`/`Decision29` "Is it EDM RISLIP?" (`When IsEDMRiSlip`),
`Decision30`/`Decision31` "Is EDM PPN PPH?" (`When IsEdmPPNPPH`),
`Decision28` "Is EDM Internal Retro" (`When IsEdmInternalRetro`),
`Decision19` "IS LIFE?".

`[terverifikasi]` Simpan lewat **`SaveEDMToJsonPolicy_Act`**, bukan `SaveJsonPolicyFacIn_Act`.

### 2.2 `OfferFacRetro` — varian berbeda, perbedaan perilaku terukur

`[terverifikasi]`

| | NB FacIn / RNW Fac In (`bfd6252070`) | **Endorsment Fac In** (`b370146c63`) |
| --- | --- | --- |
| Anak tangga persetujuan fac out | **2** (`ReasFacOutAdmin` → `ReasFacOutHead`) | **4** (+ `ReasFacOutGroupLeader` → `ReasFacOutTechnicalDirector`) |
| Decision `IsUWAccepted` | 2 | 4 |
| Tujuan `reject` di tangga ≥2 | `END52` / `Assignment1` | **seluruhnya kembali ke `Assignment1`** (Admin) |

Audit:
```
awk -f flow.awk "NB FacIn/Flow/OfferFacRetro.xml"          > a.tsv
awk -f flow.awk "Endorsment Fac In/Flow/OfferFacRetro.xml" > b.tsv
diff <(sort -u a.tsv) <(sort -u b.tsv)
```
→ 2 Assignment baru, 2 Decision baru, 6 connector baru, 1 connector berubah tujuan.

**Ini perbedaan perilaku, bukan kosmetik.** Alasannya **tidak ada di korpus**.

### 2.3 Taksonomi tipe endorsement — terbaca penuh

`[terverifikasi]` **14 rule `When`** menguji `QuotationData.StatusBusiness = 3` **dan**
`QuotationData.EdmType = 4` **dan** nilai `QuotationData.Type`, menghasilkan **12 nilai berbeda**:

| `Type` | Rule `When` |
| ---: | --- |
| `0` | `IsEdmAdjRefNo`, `IsEDMRiSlip` (via `EdmTypeNew = 4`) |
| `1` | `IsEdmExtendPeriod` |
| `2` | `IsEdmAdjTSI` |
| `3` | `IsEdmAdjRate` |
| `4` | `IsEdmAdjSpreading` |
| `5` | `IsEdmAddObject` |
| `6` | `IsEdmAdjPeriod` |
| `7` | `IsEdmAdjCurrency` |
| `8` | `IsEdmAdjRIC` |
| `9` | `IsEdmAdjInsured` |
| `11` | `IsEdmAdjShareCedant` |
| `12` | `IsEdmAdjCeding`, `IsEdmAdjRIC`, `IsEdmPPNPPH` |

`[terverifikasi]` Nilai **`10` tidak dipakai**; nilai **`12` dipakai tiga rule**.
**Arti setiap nilai belum terverifikasi** (OQ-020) — nama rule bukan bukti.

`[terverifikasi]` `When/IsEdmPerubahan.xml` adalah **komposit**: kondisinya rujukan
"Rule[X] evaluates to true" atas tujuh dari dua belas tipe.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 582 Activity, 433 Section, 265 FlowAction, 202 When,
188 Connect-SQL, 157 DataTransform, 138 ReportDefinition, 41 DataPage, 37 Harness,
10 DecisionTable, 3 ConnectREST, 2 Flow.

| Objek | Rule perujuk |
| --- | ---: |
| **`JSON_POLIS`** | 17 |
| `FACINPRODUCTION` | 6 |
| `RW`, `POOLDATA.REINSURANCETYPE`, `OCCUPATION`, `CURRENCY` | 5 masing-masing |
| `TREATYBUSINESS`, `PROPORTIONALARRG`, `POOLDATA.M_LIMIT_PROPERTYY`, `POOLDATA.JSON_POLIS`, `MARKETINGOFFICER`, **`FACOUTPRODUCTION`**, `ACCUMULATION` | 4 masing-masing |

`[terverifikasi]` **`FACOUTPRODUCTION` (4 rule) hanya menonjol di modul ini** — sejalan dengan
tangga fac out yang lebih panjang (§2.2).

## 4. Integrasi eksternal

`[terverifikasi]` Tiga `RULE-CONNECT-REST`: `ServiceGoogle`, `convertJsonNusareToProduction`,
`getPremiumPaidOn` — seluruhnya `SETTING` → `LinkService!LinkService`, **tanpa URL literal**.

`[terverifikasi]` **`GeminiAIGoogle_Act` TIDAK ada di modul ini** — berbeda dari NB FacIn dan
RNW Fac In (OQ-048). Tiga activity Google Storage lainnya ada.

### 4.1 `serviceInsertArasapas_act` — varian modul ini ADA, telusur tidak terhenti

`[terverifikasi]` Berbeda dari NB Treaty In, NB FacIn, dan RNW Fac In, modul ini **memiliki
implementasinya**: `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT`, **231.975 byte**,
hash `7314b6c49c`. Dua puluh langkahnya terbaca:

`Obj-Refresh-And-Lock` → `Property-Set` → `RDB-List` → `Property-Set` →
**`Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService`** → **`Connect-REST` ×2** → `Property-Set` →
`Obj-Save` → `Page-Set-Messages` → `RDB-List` ×2 → `Property-Set` ×2 → `RDB-List` ×2 →
`Page-Remove` → `Call ASMForceCaseClose` → `Property-Set` → `Call SendEmailWithAttachments`.

Rule yang dirujuk: `ASM GetPolicyNoByCaseId`, `ASM INSERTJSON_JSONPOLISMONITORING_FACIN`,
`ASM UpdateErrorNoteJsonPolisMonitoring`, **`RNM CekSTSKonversiJson`**,
**`RNM DeleteDataProduction`**, dan fungsi kustom **`@ASM.GetPageJSONString()`**.

`[terverifikasi]` Dua rule memakai prefix **`RNM`**, sisanya `ASM`, **di class yang sama** — bukti
konteks baru untuk OQ-008.

`[terverifikasi]` **Kedua varian implementasi di korpus berbeda hanya pada 3 tag metadata Pega**
(`<pyDelete>`, `<pyVersionSecure>`, `<pzIsPrivateCheckOut>`); **urutan 20 langkah dan daftar rule
yang dirujuk identik**. Jadi konflik OQ-011 #415 adalah **artefak check-out/versi, bukan
percabangan perilaku** → OQ-025 diperbarui.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **NB FacIn / RNW Fac In** | satu ruleset (OQ-015); `SetOldData` (`df291257aa`), `CountPremiNusareRetro_Act` (`8f1268a660`), 3 rule guard identitas — seluruhnya **identik** | `[terverifikasi]` |
| **EDM Treaty In** | keduanya memuat varian `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT`; satu-satunya dua modul yang punya implementasinya | `[terverifikasi]` |
| **Master rating eksternal** | 8 objek db-link `@ASMD.SINARMAS.CO.ID`, termasuk `LST_KURS_STANDARD` | `[terverifikasi]` — OQ-017 |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **28 stored procedure**, isinya tidak ada di korpus (OQ-002). Yang **eksklusif
modul ini**: **`FIRE.CEK_PRORATA_TANGGAL`** (perhitungan prorata tanggal, dipanggil dari
`RDBList/SearchSQLRateKPR.xml` dan `SearchSQLRateNonKPR.xml`), **`FIRE.PEGA_FIRE_SET_RATE`**,
**`NEW_GENERAL.CEK_PLAT_NO`**. Skema `NEW_GENERAL` baru → OQ-016.

`[terverifikasi]` **Perhitungan prorata tanggal berada di database**, bukan di rule.
`ProRateType` dipakai paling padat di sini (**23 kemunculan berkondisi** vs 12 di NB FacIn, 9 di
RNW Fac In). Activity prorate eksklusif: `SetLocalNProrate` (76.187 B),
`SetLocalNonMbuProrate` (297.476 B), `ChangeProRatePreActPA_FacIn` — **isinya belum dibaca**.

### 6.1 Class baru bagi D2 — `ASM-SFAGIS-`

`[terverifikasi]` Mesin salinan nilai lama endorsement tinggal di class **`ASM-SFAGIS-WORK-ENDORSEMENT`**
(31 rule di korpus) — prefix yang **tidak** mengikuti pola `ASM-FW-GISFW-`:

| Activity | Ukuran | Langkah |
| --- | ---: | ---: |
| `SetValueToEDMWork.xml` | 445.703 | 41 |
| `SetErrorBatalEndorsement_Act.xml` | 425.872 | 39 |
| `CountEndorsementData.xml` | 357.791 | 25 |
| `SetOLDValueToEDMWork_{Aneka,FIRE,GOLF,LIFE,MBU,MC,PA}.xml` | — | 7 varian per lini |
| `SetEdmType.xml` | 73.041 | 4 |

Ditambah `EDMRetro_Act.xml` di class `ASM-FW-GISFW-WORK-ENDORSEMENT`. **Seluruhnya belum dibaca**
— batas **cakupan telusur**. Empat class bersufiks `ENDORSEMENT` hidup berdampingan di korpus
→ **OQ-049**.

`[terverifikasi]` Rule modul ini menyebut `OldData` di **133 file** — terbanyak dari trio
facultative (NB FacIn 129, RNW Fac In 114).

`[terverifikasi]` `DecisionTable/IsUWAccepted.xml` **berbeda isinya** di sini (`9226a7db77`)
dibanding NB FacIn/RNW Fac In (`f99bc43c45`) — aturan persetujuan **memang bercabang antar siklus**,
tetapi **apa yang bercabang tidak dapat dinyatakan** karena baris tabelnya tidak terekspor
→ gabungan OQ-011 × **OQ-043**.

`[terverifikasi]` Modul ini punya **6** Data Transform ber-`ProposalAcceptStatus`, bukan 8 —
`SetAkseptasiCeding_DT` dan `SetReviseProposal` **tidak ada**.

`[terverifikasi]` **38 activity spreading/scoring belum habis dibaca**. `CekLimitSpreading_Act`
**tidak ada** di sini; sebagai gantinya ada `CopySpreading_Act` (627.836 B),
`SetSpreadingCoverage_ACT`, `SetSumTSISpreading_Act`, `SpreadingProtection` — eksklusif modul ini.

`[terverifikasi]` `When/IsEdmInternalRetro.xml` — kondisinya **tidak terbaca** → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **33 file**,
`OperatorID.pyTelephone` **2 file**; tiga rule guard identik dengan NB FacIn → OQ-021.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-008**, **OQ-011**, **OQ-012**, **OQ-014**, **OQ-015**, **OQ-016**, **OQ-017**,
**OQ-018**, **OQ-020**, **OQ-021**, **OQ-023**, **OQ-025**, **OQ-026**, **OQ-028**, **OQ-029**,
**OQ-043**, **OQ-045**, **OQ-046**, **OQ-047**, **OQ-048**, **OQ-049**.
