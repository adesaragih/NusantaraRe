# Modul — EDM Treaty In

STEP D3. Sintesis dari `../inventory/EDM Treaty In.md` (D1) + `../flows/EDM Treaty In.md`
(D2 Tahap 1). **163 file**.
Audit: `find "EDM Treaty In" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **endorsement / addendum treaty inward**. Titik masuk
`EDM Treaty In/Flow/InputAddendumTreatyIn.xml` → `RULE-OBJ-FLOW` /
**`ASM-FW-GISFW-WORK-ENDORSEMENTTREATY`** / `INPUTADDENDUMTREATYIN`, `<pyStartActivity>Start1`,
ukuran 202.269 byte.

`[terverifikasi]` Class terbanyak justru `ASM-FW-GISFW-DATA-POLICYTREATYIN` (46 rule), lalu
`ASM-FW-GISFW-WORK` (33) — modul ini berbagi banyak class dengan `NB Treaty In`.

## 2. Proses / fitur utama

Rincian: **`../flows/EDM Treaty In.md`** (20 rule ditelusur — terbanyak di Tahap 1).

`[terverifikasi]` Alur bertulang sama dengan realisasi treaty inward (`NB Treaty In`), dengan
FlowAction berakhiran **`Addendum`**: `FlowAction/InboxPolicyTreatyInAddendum.xml`,
`FlowAction/DeptHeadTreatyIn_UWAddendum.xml`.

### 2.1 Bukti untuk OQ-014 (arti `EDM`)

`[terverifikasi]` Telusur ini mengumpulkan **bukti penamaan, bukan kepanjangan singkatannya**:

| Bukti | Path |
| --- | --- |
| Flow bernama **`InputAddendumTreatyIn`** berada di modul bernama **`EDM Treaty In`** | `Flow/InputAddendumTreatyIn.xml` |
| Class flow adalah `ASM-FW-GISFW-WORK-**ENDORSEMENTTREATY**` | `<pxInsName>` flow |
| FlowAction berakhiran **`Addendum`** | `FlowAction/InboxPolicyTreatyInAddendum.xml` |
| Activity simpan berakhiran **`EDM`** untuk proses yang sama | `Activity/SaveJsonPolisTreatyInEDM_Act.xml` |
| Di NB Treaty In, `param.isFOR` bernilai `"EDM"` atau `"POLICY"` sebagai **dua mode simpan** | `NB Treaty In/Activity/SaveJsonPolisTreatyIn_Act.xml` |

`[dugaan]` **EDM**, **Addendum**, dan **Endorsement** dipakai bergantian untuk konsep yang sama:
perubahan atas polis yang sudah ada, dibedakan dari penutupan baru (`POLICY`).
**Kepanjangan `EDM` tetap belum terverifikasi** → **OQ-014**.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 66 Activity, 36 Connect-SQL, 18 Section, 11 When,
10 DataTransform, 9 ReportDefinition, 6 FlowAction, 3 Harness, 1 Flow, 1 DecisionTable,
1 ConnectREST, 1 SystemSettings.

| Objek | Rule perujuk |
| --- | ---: |
| `POOLDATA.JSON_POLIS` | 4 |
| **`POOLDATA.M_TREATY_IN_EDM`**, `JSON_POLIS` | 3 masing-masing |
| `TREATYINPRODUCTION`, `REINSURANCETYPE`, `CURRENCY` | 2 masing-masing |
| **`POOLDATA.TREATY_OUT2`**, `POOLDATA.TREATY_IN_EDM`, `POOLDATA.TREATY_IN`, `POOLDATA.TREATYINPRODUCTION` | 1 masing-masing |

`[terverifikasi]` Modul ini menyentuh objek treaty **outward** di **4 file**
(`POOLDATA.TREATY_OUT2`, `RDBList/BrowseTreatyOutDetailEDM.xml` berclass
`ASM-FW-GISFW-INT-TREATYOUTDETAIL`) — salah satu dari lima modul yang menyentuhnya, dan
**`Treaty Contract Out` bukan salah satunya** → OQ-022.

## 4. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `convertJsonNusareToProduction`
(`SETTING` → `LinkService!LinkService`, tanpa URL literal).

`[terverifikasi]` **Modul ini memuat salah satu dari dua implementasi `SERVICEINSERTARASAPAS_ACT`**
(`ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT`, 232.000 B, hash **`ab3ae3c9ba`**). Varian kedua ada
di `Endorsment Fac In` (`7314b6c49c`).

`[terverifikasi]` Kedua varian berbeda **hanya pada 3 tag metadata Pega** (`<pyDelete>`,
`<pyVersionSecure>`, `<pzIsPrivateCheckOut>`); **urutan 20 langkah dan daftar rule yang dirujuk
identik** → konflik OQ-011 #415 adalah artefak check-out, bukan percabangan perilaku → OQ-025.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **NB Treaty In** | berbagi class `ASM-FW-GISFW-DATA-POLICYTREATYIN` (46 rule) dan `ASM-FW-GISFW-WORK`; memanggil `POOLDATA.PEGA_JSON_POLIS_TREATYIN` dan `POOLDATA.PEGA_TREATY_IN` yang sama | `[terverifikasi]` |
| **NB Treaty In + NB FacIn + RNW Fac In** | ketiganya **memanggil** implementasi Arasapas yang ada di modul ini, tetapi **tidak memuatnya** | `[terverifikasi]` — OQ-025 |
| **Treaty In (master)** | `RDBList/SaveTreatyIn.xml` (`ASM!SAVETREATYIN`, `e371c194`) & `BrowseTreatyIn.xml` (`196d49b5`) — identik di 6 modul | `[terverifikasi]` |
| **Treaty In Adjustment** | `RDBList/BrowseTreatyOutDetailEDM.xml` ada di kedua modul | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Tujuh stored procedure**, isinya tidak ada di korpus (OQ-002):
`POOLDATA.INSERTJSONPOLISMONITORING`, **`POOLDATA.INSERTUPDATEACHIEVMENT`**,
`POOLDATA.PEGA_DELETE_ERROR_KONVERSI`, `POOLDATA.PEGA_JSON_POLIS_TREATYIN`,
`POOLDATA.PEGA_TREATY_IN`, `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`, `DBMS_LOB.CREATETEMPORARY`.

`[terverifikasi]` Struktur `JSONDATA` pada `M_TREATY_IN_EDM` dan `JSON_POLIS` **tidak ada di
korpus** → OQ-012.

`[terverifikasi]` `DecisionTable/isApproved.xml` ada di modul ini — baris keputusannya **tidak ikut
terekspor**, seperti seluruh 49 `DecisionTable` korpus → **OQ-043**, dan bersama
`When/isApproved.xml` di NB Treaty In → **OQ-026**.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **5 file**,
`OperatorID.pyTelephone` **2 file**. Nilai nama orang **tidak disalin** → OQ-021, OQ-027.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-011**, **OQ-012**, **OQ-014**, **OQ-018**, **OQ-021**, **OQ-022**, **OQ-024**,
**OQ-025**, **OQ-026**, **OQ-027**, **OQ-028**, **OQ-029**, **OQ-043**.
