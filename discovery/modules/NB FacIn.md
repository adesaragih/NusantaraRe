# Modul — NB FacIn

STEP D3. Sintesis dari `../inventory/NB FacIn.md` (D1) + `../flows/NB FacIn.md` (D2 Tahap 4) +
`../flows/_SUMMARY-facultative.md`. **2.083 file — modul terbesar di korpus**.
Audit: `find "NB FacIn" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **penutupan baru (new business) facultative inward**: dari quotation →
penawaran → akseptasi berjenjang → binding → R/I slip → fac out/retro → produksi. Class kerja
`ASM-FW-GISFW-WORK` (551 rule).

`[terverifikasi]` Modul ini juga **simpul yang menyatukan lini facultative dan treaty inward**:
`NB FacIn/Flow/InputQuotation.xml` (`RULE-OBJ-FLOW` / **`ASM-FW-GISFW-WORK-NB`** / `INPUTQUOTATION`)
bercabang ke `InputInwardFacultativeOffer` atau `InputRealizationTreatyIn`.

## 2. Proses / fitur utama

Rincian: **`../flows/NB FacIn.md`** (45 rule ditelusur). **Enam rule `Flow`**:

| File | `pxInsName` | Start | Hash | Telusur |
| --- | --- | --- | --- | --- |
| `Flow/InputQuotation.xml` | `ASM-FW-GISFW-WORK-NB!INPUTQUOTATION` | `Start1` | `d30995fece` | ya |
| `Flow/InputInwardFacultativeOffer.xml` | `ASM-FW-GISFW-WORK!INPUTINWARDFACULTATIVEOFFER` | `Start1` | `39cfbbec4f` | ya |
| `Flow/InputInwardFacultativeRISlip.xml` | `…!INPUTINWARDFACULTATIVERISLIP` | `Start1` | `c099bf4ebb` | ya (identik dgn RNW) |
| `Flow/OfferFacOut.xml` | `…!OFFERFACOUT` | `Start62` | `832fb12b9d` | ya (identik dgn RNW) |
| `Flow/OfferFacRetro.xml` | `…!OFFERFACRETRO` | `Start1` | `bfd6252070` | ya — varian NB/RNW |
| `Flow/InputRealizationTreatyIn.xml` | `…!INPUTREALIZATIONTREATYIN` | `Start1` | `1306f68d56` | **dipakai ulang** dari `NB Treaty In.md` |

`[terverifikasi]` `InputInwardFacultativeOffer` adalah **graf terbesar korpus**: 21 Assignment,
42 Decision, 12 Utility, 3 SubProcess, **170 connector**. Sebagai pembanding, flow terbesar domain
Claim punya 7 shape.

`[terverifikasi]` Tahapan: Marketing → (Life? → Medical Life) → simpan JSON offer →
Group?/PKS?/UW Financial? → tangga akseptasi berbasis limit → binding → R/I Slip → fac out/retro →
produksi.

`[terverifikasi]` Seluruh 21 Assignment memakai `WorkBasket`/`Custom` — **tidak satu pun
`WorkList`**, kecuali `OfferFacOut` yang berupa screen flow (`Data-MO-Event-Start-StartScreenFlow`)
dan memakai `WorkList` → OQ-028.

`[pertanyaan terbuka]` **16 shape yatim** di flow terbesar: 5 tanpa connector masuk
(`Decision14/16/21/38/42`), 11 tanpa connector keluar → OQ-023 (skalanya jauh di atas 5 shape yang
ditemukan di NB Treaty In).

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 609 Activity, 431 Section, 250 FlowAction,
217 Connect-SQL, 209 When, 148 DataTransform, 123 ReportDefinition, 41 Harness, 30 DataPage,
12 DecisionTable, 6 Flow, 3 ConnectREST.

| Objek | Rule perujuk |
| --- | ---: |
| **`JSON_POLIS`** | 18 — objek paling dirujuk dari modul mana pun yang ditelusur |
| `FACINPRODUCTION` | 7 |
| `RW`, `ACCUMULATION` | 6 masing-masing |
| `TREATYBUSINESS`, `POOLDATA.REINSURANCETYPE`, `POOLDATA.JSON_POLIS`, `OCCUPATION`, `CURRENCY` | 5 masing-masing |
| `PROPORTIONALARRG`, **`POOLDATA.M_LIMIT_PROPERTYY`**, `POOLDATA.CLIENT`, `MARKETINGOFFICER` | 4 masing-masing |
| `T_STORAGE_IMAGE`, **`POOLDATA.M_LIMIT_NONPROPANDENGG`** | 3 masing-masing |

`[terverifikasi]` **Limit akseptasi diambil dari database**: `Activity/GetLimitAkseptasi_ActFlow.xml`
(539.687 B) memanggil **12 RDB-List** di class `ASM-FW-GISFW-Int-policyjson`, terpisah per lini
bisnis (Prefered Comm, Non Prefer, Non Fire, Kredit NCL, Kredit CL, Bond, Engineering) dengan
**varian "Banding" tersendiri** → OQ-040.

## 4. Integrasi eksternal

`[terverifikasi]` Tiga `RULE-CONNECT-REST`, seluruhnya `pyBaseURLSelectionType = SETTING` →
`LinkService!LinkService`, **tanpa URL literal endpoint**: `ServiceGoogle`,
`convertJsonNusareToProduction`, `getPremiumPaidOn`.

`[terverifikasi]` Alamat sesungguhnya diambil dari **tabel Oracle `M_LINK_SERVICE`** lewat
`Activity/GetLinkService.xml` (`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE`, `Obj-Browse`
dengan kunci `KATEGORI_1`/`KATEGORI_2`) → **OQ-047**.

`[terverifikasi]` **`GeminiAIGoogle_Act.xml` ada di modul ini** (dan di RNW Fac In, tidak di
Endorsment Fac In) — pemanggilan model AI pihak ketiga di dalam alur underwriting.
**Isinya belum dibaca** → **OQ-048**.

`[terverifikasi]` Shape `Utility2` "HIT SERVICE ARASAPAS" → `serviceInsertArasapas_act` —
**terblokir**, lihat §6.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **RNW Fac In / Endorsment Fac In** | satu ruleset: `IsNB`/`IsRenewal`/`IsEDM` ada di ketiganya; 1.906 dari 1.926 identitas RNW (99,0 %) juga ada di sini | `[terverifikasi]` — OQ-015 |
| **NB Treaty In** | `Flow/InputRealizationTreatyIn.xml` **identik** (`1306f68d56`); procedure `POOLDATA.PEGA_TREATY_IN` dan `POOLDATA.PEGA_JSON_POLIS_TREATYIN` dipanggil dari sini | `[terverifikasi]` |
| **Treaty In (master)** | `RDBList/SaveTreatyIn.xml` (`ASM!SAVETREATYIN`, hash `e371c194`) dan `BrowseTreatyIn.xml` (`196d49b5`) — **identik di 6 modul**, tidak di register OQ-011 | `[terverifikasi]` |
| **Master rating eksternal** | 8 objek lewat db-link `@ASMD.SINARMAS.CO.ID` | `[terverifikasi]` — OQ-017 |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **29 stored procedure** dipanggil, isinya tidak ada di korpus (OQ-002). Yang khas
facultative: `POOLDATA.GENERATE_FACRETRO_NO`, `POOLDATA.PEGA_M_JSON_OFFER`. Dua skema baru muncul
di sini: `GENERAL` (`F_GET_NM_ASURADUR`), `MBU` (`F_CEK_HURUF`) → OQ-016.

`[terverifikasi]` **8 objek db-link** `@ASMD.SINARMAS.CO.ID`: `M_EQS_RATE`, `M_TERORISME_RATE`,
`M_RSMD_RATE`, `M_FLEXAS_RATE`, `M_FLOOD_AREA`, `FIRE.M_FLOOD_RATE`, `FIRE.M_BI_INDEMNITY`,
`LST_KURS_STANDARD` — **input rating berada di luar basis data yang dimigrasikan** → OQ-017.

`[terverifikasi]` **Baris tabel keputusan tidak ada di ekspor.**
`DecisionTable/IsUWAccepted.xml` menggerbangi **22 shape Decision** dan mendeklarasikan enam hasil
(`confirm`, `reject`, `ask`, `banding`, `revise`, `decline`) lewat `<pyTaskStatusXml>` — tetapi
baris kondisi→hasilnya **tidak ikut terekspor**, seperti seluruh 49 `DecisionTable` korpus
→ **OQ-043**.

`[terverifikasi]` **`ProposalAcceptStatus`**: nilai literal yang ditulis terbaca
(`1`, `2`, `3`, `4`, `7`, `9`, dari 8 Data Transform `Set*Proposal*`), tetapi **pemetaannya ke enam
hasil tidak dapat dinyatakan** → OQ-020, OQ-043. Nilai `4` **ambigu**: ditulis
`SetBandingProposal_DT` tetapi diuji `When/IsFacout.xml` → **OQ-044**.

`[terverifikasi]` **`pyWorkPage.LetterNo` dipakai sebagai token routing persetujuan** dengan 10 kode
peran literal (`KADIVTEKNIK`, `DIREKTURTEKNIK`, `DEPHEADUNDERWRITER`, …) → **OQ-045**.

`[terverifikasi]` **38 activity spreading/capacity/scoring belum habis dibaca** — terbesar
`SumTSIPremiSpreadedRNM_FIRE_Act.xml` **996.052 byte**, terbesar di seluruh D2. Ini batas
**cakupan telusur**, bukan batas korpus.

`[terverifikasi]` **Satu rumus terbaca penuh**: `Activity/CountPremiNusareRetro_Act.xml`
(hash `8f1268a660`, identik di ketiga modul facultative) memuat pangsa **0,45** / **0,05** dan
ambang **3.000.000.000** yang dibandingkan **sebagai string berkutip**; mata uang tidak disebut
→ **OQ-046**.

`[terverifikasi]` **Blocker `serviceInsertArasapas_act`**: file di modul ini beridentitas
`ASM-FW-GISFW-DATA-POLICYTREATYIN!SERVICEINSERTARASAPAS_ACT` (18.013 B) — hanya pembungkus satu
langkah `Call serviceInsertArasapas_act` dengan **`<pyStepsClassName>ASM-FW-GISFW-Work`**.
Implementasi 232 KB berclass `ASM-FW-GISFW-WORK` **tidak ada di modul ini** → telusur dihentikan,
varian modul lain **tidak dipinjam** → **OQ-025**.

`[terverifikasi]` **Empat `When` kondisinya tidak terbaca**: `IsPKSASM`, `LetterNoNull`, `ToJUW_A`,
`ToUW` → OQ-029. `PKS`, `ASM` kepanjangan belum terverifikasi.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **43 file** (terbanyak di
korpus), `OperatorID.pyTelephone` **4 file**. Tiga rule guard — `When/IsGroup.xml` (3 identitas),
`When/IsGroupCreate.xml` (2), `When/IsTBonding.xml` (1 lewat `.MarketingName`, **field data bisnis**,
pola baru) — **identik isinya** di ketiga modul facultative (`IsGroup` hash `2cff7b9a67`).
Nilai nama orang **tidak disalin** → OQ-021.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-011**, **OQ-012**, **OQ-015**, **OQ-016**, **OQ-017**, **OQ-018**, **OQ-020**,
**OQ-021**, **OQ-023**, **OQ-024**, **OQ-025**, **OQ-026**, **OQ-028**, **OQ-029**, **OQ-038**,
**OQ-040**, **OQ-043**, **OQ-044**, **OQ-045**, **OQ-046**, **OQ-047**, **OQ-048**.
