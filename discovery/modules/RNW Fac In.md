# Modul — RNW Fac In

STEP D3. Sintesis dari `../inventory/RNW Fac In.md` (D1) + `../flows/RNW Fac In.md` (D2 Tahap 4) +
`../flows/_SUMMARY-facultative.md`. **1.927 file**.
Audit: `find "RNW Fac In" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **renewal facultative inward**. Class kerja `ASM-FW-GISFW-WORK`
(510 rule) — class yang sama dengan NB FacIn dan Endorsment Fac In.

`[terverifikasi]` **Bukan basis kode terpisah**: 1.906 dari 1.926 identitas rule modul ini
(**99,0 %**) juga ada di `NB FacIn` (`../inventory/_summary.md`, OQ-015). Pembedanya dibaca saat
runtime — lihat §5.

## 2. Proses / fitur utama

Rincian: **`../flows/RNW Fac In.md`** (28 rule ditelusur). **Empat rule `Flow`**:

| File | Start | Hash | Telusur |
| --- | --- | --- | --- |
| `Flow/InputRenewalFacultativeIn.xml` | **`Start2`** | `f3dff3940e` | **khas modul ini** |
| `Flow/InputInwardFacultativeRISlip.xml` | `Start1` | `c099bf4ebb` | **identik** dgn NB FacIn → tidak ditelusur ulang |
| `Flow/OfferFacOut.xml` | `Start62` | `832fb12b9d` | **identik** dgn NB FacIn |
| `Flow/OfferFacRetro.xml` | `Start1` | `bfd6252070` | **identik** dgn NB FacIn (varian NB/RNW) |

`[terverifikasi]` Modul ini **tidak punya** `InputInwardFacultativeOffer`, `InputQuotation`, maupun
`InputRealizationTreatyIn`.

`[terverifikasi]` `InputRenewalFacultativeIn`: 17 Assignment, 36 Decision, 11 Utility,
5 SubProcess, 140 connector. Tahapan: Marketing → Group?/PKS?/UW Financial? → tangga akseptasi
berbasis limit → binding → R/I Slip → fac out/retro → produksi.

### 2.1 Selisih terukur dari siklus new business

| Aspek | NB FacIn | **RNW Fac In** |
| --- | --- | --- |
| Titik masuk | `Start1` | **`Start2`** |
| Assignment / Decision | 21 / 42 | 17 / 36 |
| **Cabang Life** | ada lengkap (Medical Life, UW Life, DeptHead UW Life) | **tidak ada sama sekali** |
| Activity limit | `GetLimitAkseptasi_ActFlow` + `GetLimitAkseptasiLife_Act` | `GetLimitAkseptasi_Act`, **tanpa varian Life** |
| Shape Arasapas | `serviceInsertArasapas_act` | **`serviceInsertArasapasRNW_act`** (eksklusif) |
| Simpan | `SaveJsonOfferFacIn_Act` + `SaveJsonPolicyFacIn_Act` | `SaveJsonPolicyFacIn_Act`; **tanpa SAVE JSON_OFFER** |

`[terverifikasi]` **Tidak ada jalur Life di siklus renewal** — selisih perilaku terbesar.

`[terverifikasi]` Activity yang **eksklusif** RNW terhadap NB FacIn hanya **tiga**:
`GetBusinessGroup_Act`, `SetDataInsuredEDM_Act`, `serviceInsertArasapasRNW_act`. Sebaliknya
**57 activity NB FacIn tidak ada di sini** — sebagian besar terkait treaty inward
(`SaveJsonPolisTreatyIn_Act`, `SetTreatyIn_Act`, `TreatyIn*`, `GeneratePolicyNoTreaty_Act`) dan
perhitungan premi/komisi (`CalculatePremi_Act`, `CountResult*`, `CountRiComm*`, `SetPPNPPH`).

Audit:
```
comm -23 <(ls "RNW Fac In/Activity" | sort) <(ls "NB FacIn/Activity" | sort) | wc -l   # 3
comm -13 <(ls "RNW Fac In/Activity" | sort) <(ls "NB FacIn/Activity" | sort) | wc -l   # 57
```

### 2.2 Penanda renewal yang terbaca

`[terverifikasi]` `RNW Fac In/Activity/ProtectRenewal_Act.xml`
(`ASM-FW-GISFW-WORK!PROTECTRENEWAL_ACT`, 46.258 B, 3 langkah): precondition `When IsRenewal`,
ekspresi `@contains(pyWorkPage.pxInsName,"NB-")` dan `Local.DateDif >= 0`, lalu
`Page-Set-Messages` (`Local.ErrMsg`). **Prefix ID case membawa penanda siklus.** Isi pesan dan arti
`Local.DateDif` **belum terverifikasi**.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 556 Activity, 396 Section, 239 FlowAction,
200 Connect-SQL, 188 When, 138 DataTransform, 118 ReportDefinition, 41 Harness, 30 DataPage,
10 DecisionTable, 4 Flow, 3 ConnectREST.

| Objek | Rule perujuk |
| --- | ---: |
| **`JSON_POLIS`** | 18 |
| `FACINPRODUCTION` | 7 |
| `RW`, `ACCUMULATION` | 6 masing-masing |
| `POOLDATA.REINSURANCETYPE`, `POOLDATA.JSON_POLIS`, `OCCUPATION`, `CURRENCY` | 5 masing-masing |
| `TREATYBUSINESS`, `POOLDATA.M_LIMIT_PROPERTYY`, `MARKETINGOFFICER` | 4 masing-masing |

`[terverifikasi]` **Selisih dari NB FacIn:** `POOLDATA.CLIENT` tidak masuk daftar teratas di sini,
`TREATYBUSINESS` turun 5→4. Sebaran keseluruhan **hampir sama persis** — konsisten dengan §5.

## 4. Integrasi eksternal

`[terverifikasi]` Tiga `RULE-CONNECT-REST`, sama persis dengan NB FacIn: `ServiceGoogle`,
`convertJsonNusareToProduction`, `getPremiumPaidOn` — seluruhnya `SETTING` →
`LinkService!LinkService`, **tanpa URL literal**.

`[terverifikasi]` **`serviceInsertArasapasRNW_act` eksklusif modul ini**
(`ASM-FW-GISFW-WORK!SERVICEINSERTARASAPASRNW_ACT`, 89.947 B, hash `cb70714c1a`), 10 langkah:
`Call serviceInsertArasapas_act` → `Property-Set` → `RDB-List` → `Property-Set` →
**`Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService`** → **`Connect-REST`** → `Property-Set` →
`Call InsertLogServiceProd` → `Page-Remove` → `RDB-List`.

`[terverifikasi]` **Rantai resolusi endpoint terbaca di sini**: alamat layanan berada di tabel
Oracle **`M_LINK_SERVICE`** (kunci `KATEGORI_1`/`KATEGORI_2`), bukan hanya di SystemSettings
→ **OQ-047**. Sembilan rule modul ini memakainya
(`grep -rl "M_LINK_SERVICE" "RNW Fac In" --include="*.xml" | wc -l`).

`[terverifikasi]` **`GeminiAIGoogle_Act.xml` ada** — pemanggilan model AI pihak ketiga; isinya
**belum dibaca** → OQ-048.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **NB FacIn** | 99,0 % identitas dibagi; 3 dari 4 Flow **identik via hash**; `SetOldData.xml` identik (`df291257aa`); `CountPremiNusareRetro_Act` identik (`8f1268a660`) | `[terverifikasi]` |
| **Endorsment Fac In** | `When/IsNB.xml`, `IsRenewal.xml`, `IsEDM.xml`, `IsNotEDM.xml` ada di ketiga modul → satu ruleset berdiskriminator siklus | `[terverifikasi]` — OQ-015 |
| **Master rating eksternal** | **7** objek db-link (bukan 8) — `LST_KURS_STANDARD` **tidak** dirujuk dari modul renewal | `[terverifikasi]`; `[dugaan]` renewal memakai kurs tersimpan pada data lama — **belum terverifikasi** |
| Class `ASM-FW-GISFW-Work-Renewal` | hanya satu rule lain: `ReportDefinition/RenewalList_RD.xml` | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **26 stored procedure** — daftar NB FacIn **dikurangi tiga**:
`POOLDATA.PEGA_JSON_POLIS_TREATYIN`, `POOLDATA.PEGA_TREATY_IN`, `POOLDATA.PEGA_M_JSON_OFFER`
**tidak dipanggil** dari sini. Konsisten dengan §2.1 (tanpa jalur treaty-in & tanpa pembuatan JSON
offer baru). Isinya tidak ada di korpus → OQ-002.

`[terverifikasi]` **Blocker `serviceInsertArasapas_act` berlaku di sini juga** — file modul ini
adalah pembungkus 18 KB (`ASM-FW-GISFW-DATA-POLICYTREATYIN!…`, hash `b013027e0c`); implementasi
232 KB berclass `ASM-FW-GISFW-WORK` **tidak ada di modul ini** → OQ-025.

`[terverifikasi]` **`RNW Fac In` tidak punya `When/IsUWAccepted.xml`** — hanya
`DecisionTable/IsUWAccepted.xml` (hash `f99bc43c45`, sama dengan NB FacIn). Ini **mempersempit
OQ-026**, tetapi baris tabelnya tetap tidak ada di ekspor → OQ-043.

`[terverifikasi]` **34 activity spreading/scoring belum habis dibaca**; terbesar
`SumTSIPremiSpreadedRNM_FIRE_Act.xml` 996.052 B. Yang tidak ada di sini tetapi ada di NB FacIn:
`BreakDownSpreading_Act`, `CountSpreading_Act`, `SpreadingAdditionalProtection`,
`SumTreatyCapacity_Act`. Batas **cakupan telusur**.

`[terverifikasi]` `Activity/SetOldData.xml` (450.829 B, 59 langkah `Property-Set`) — mesin salinan
data lama, **identik di ketiga modul facultative**; **belum habis dibaca**.

`[terverifikasi]` **Prorate:** `ProRateType` diuji dengan nilai `3` dan `4` (9 kemunculan
berkondisi, vs 23 di Endorsment Fac In). Procedure `FIRE.CEK_PRORATA_TANGGAL` **tidak dipanggil**
dari modul ini. Arti `ProRateType` **belum terverifikasi** → OQ-020.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **36 file**,
`OperatorID.pyTelephone` **2 file**. Tiga rule guard identik dengan NB FacIn → OQ-021.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-011**, **OQ-012**, **OQ-015**, **OQ-017**, **OQ-018**, **OQ-020**, **OQ-021**,
**OQ-025**, **OQ-026**, **OQ-028**, **OQ-043**, **OQ-045**, **OQ-046**, **OQ-047**, **OQ-048**.
