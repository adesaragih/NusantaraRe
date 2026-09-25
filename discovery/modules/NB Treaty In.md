# Modul — NB Treaty In

STEP D3. Sintesis dari `../inventory/NB Treaty In.md` (D1) + `../flows/NB Treaty In.md`
(D2 Tahap 1). **278 file**.
Audit: `find "NB Treaty In" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **realisasi penutupan baru treaty inward** — tahap setelah penawaran,
dengan tangga akseptasi berbasis workbasket. Class kerja `ASM-FW-GISFW-WORK` (90 rule).
Titik masuk `NB Treaty In/Flow/InputRealizationTreatyIn.xml` → `RULE-OBJ-FLOW` /
`ASM-FW-GISFW-WORK` / `INPUTREALIZATIONTREATYIN`, `<pyStartActivity>Start1`.

`[terverifikasi]` Rule `Flow` ini punya salinan **berisi identik** di
`NB FacIn/Flow/InputRealizationTreatyIn.xml` (hash `1306f68d56`) — satu rule, dua modul.

## 2. Proses / fitur utama

Rincian: **`../flows/NB Treaty In.md`** (15 rule ditelusur).

`[terverifikasi]` Graf: 6 Assignment, 12 Decision, 2 Utility, 30 connector.
Tahapan: Input Realitation → "Is Correct?" (`When isApproved`) → percabangan SPV/Treaty →
Acceptance by Head. Treaty → (opsional) Acceptance by Dept. Head → cek nomor polis →
simpan JSON policy → HIT SERVICE ARASAPAS.

`[terverifikasi]` **Keenam Assignment memakai `WorkBasket`/`Custom`** — tidak satu pun `WorkList`
→ OQ-028. Nama workbasket yang muncul di `<pyRuleName>`: `ReasTreatyInAdmin`,
`ReasTreatyInSecHead`, `ReasTreatyInGroupLeader`, `ReasTreatyInDeptHead`, `ReasTreatyInDirector`.
Pemetaan shape → nama workbasket **tidak terbaca** (`<pyWorkBasket>` kosong, routing `Custom`)
→ OQ-024.

`[pertanyaan terbuka]` **Sub-graf lima shape tanpa connector masuk** (`Decision5` "Claim" →
`Assignment5` "Claim Manager" → … → `Assignment1` "Acceptance by Dir.") — jalur mati atau ekspor
tidak lengkap? → **OQ-023**.

### 2.1 Format nomor polis — terbaca penuh

`[terverifikasi]` `NB Treaty In/RDBList/GenerateNoPolicy.xml`
(`ASM-FW-GISFW-INT-POLISTREATYIN` / `ASM!GENERATENOPOLICY` / `RULE-CONNECT-SQL`):

```sql
SELECT 'RNM-' || {InputData.CARI20} || '.T' || {pyWorkPage.PolicyTreatyIn.QuotationData.BusinessOldId}
       || '.' || to_char(sysdate,'MM.yyyy') || '.'
       || LPAD (POOLDATA.JSON_POLIS_TREATYIN_SEQ.NEXTVAL, 5, '0') AS "HASIL"
FROM dual
```

`[pertanyaan terbuka]` Isi `CARI20` **tidak terlihat diisi di activity ini** → OQ-059.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 92 Activity, **75 When**, 41 Connect-SQL, 25 Section,
15 ReportDefinition, 12 DataTransform, 9 FlowAction, 6 Harness, 2 DecisionTable, 1 Flow.
**Tidak ada `ConnectREST` sama sekali.**

| Objek | Rule perujuk |
| --- | ---: |
| `CURRENCY` | 3 |
| `TREATYBUSINESS`, `REINSURANCETYPE`, `PROPORTIONALARRG`, `POOLDATA.REINSURANCETYPE`, **`POOLDATA.M_TREATY_OUT`**, `POOLDATA.M_TREATY_IN_EDM`, `POOLDATA.M_TREATY_IN`, `HISTORYAKSEPTASIPEGA`, `CATEGORY_ATTACH_REAS` | 2 masing-masing |

`[terverifikasi]` **Modul ini adalah penyentuh objek treaty OUTWARD terbanyak di korpus — 8 file**
(`grep -rli "M_TREATY_OUT\|TREATY_OUT\|TREATYOUTDETAIL" "NB Treaty In" --include='*.xml' | wc -l`).
Tiga rule berclass `ASM-FW-GISFW-INT-TREATYOUTDETAIL`: `Activity/SetValueRetro_Act.xml`,
`RDBList/BrowseTreatyOutDetail.xml`, `ReportDefinition/BrowseTreatyOutDetail.xml`.
Bandingkan: modul bernama `Treaty Contract Out` menyentuhnya **0 file** → OQ-022.

## 4. Integrasi eksternal

`[terverifikasi]` **Tidak ada rule `RULE-CONNECT-REST` di modul ini.** Integrasi tersentuh sebagai
**shape**, bukan ConnectREST: `Utility2` "HIT SERVICE ARASAPAS" → `serviceInsertArasapas_act`.

`[dugaan]` Nama shape menyebut "ARASAPAS", dan D1 mencatat skema Oracle `ARASAPAS` dengan objek
`ARASAPAS.INVOICE` / `ARASAPAS.DETAIL_INVOICE` yang dirujuk modul lain. Keterkaitannya **belum
terverifikasi** — label shape bukan bukti perilaku.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **NB FacIn** | `Flow/InputRealizationTreatyIn.xml` **identik** (`1306f68d56`); `InputQuotation` di NB FacIn bercabang ke flow ini | `[terverifikasi]` |
| **Treaty In (master)** | `RDBList/SaveTreatyIn.xml` (`ASM!SAVETREATYIN`, `e371c194`) & `BrowseTreatyIn.xml` (`196d49b5`) — identik di 6 modul; procedure `POOLDATA.PEGA_TREATY_IN` | `[terverifikasi]` |
| **EDM Treaty In / Endorsment Fac In** | memanggil `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` yang **hanya ada di dua modul itu** | `[terverifikasi]` — OQ-025 |
| **Treaty outward** | `POOLDATA.M_TREATY_OUT`, class `…INT-TREATYOUTDETAIL` | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Tiga stored procedure**, isinya tidak ada di korpus (OQ-002):
`POOLDATA.PEGA_JSON_POLIS_TREATYIN`, `POOLDATA.PEGA_TREATY_IN`,
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` `RDBList/SavePolisTreatyIn_SQL.xml` memanggil
`POOLDATA.PEGA_JSON_POLIS_TREATYIN` dengan **8 parameter** dan melakukan **`COMMIT` di dalam blok
PL/SQL** — batas transaksi berada di sisi database → **OQ-013**.

`[terverifikasi]` Parameter ke-7 diisi **`@ASM.GetPageJSONString()`** — fungsi Pega kustom yang
**source-nya tidak ada di korpus**; ini sumber langsung isi kolom JSON → **OQ-012**.

### 6.1 BLOCKER — `serviceInsertArasapas_act`

`[terverifikasi]` `NB Treaty In/Activity/serviceInsertArasapas_act.xml` beridentitas
`ASM-FW-GISFW-DATA-POLICYTREATYIN!SERVICEINSERTARASAPAS_ACT` (18.013 B, hash `b013027e0c`) —
pembungkus satu langkah `Call serviceInsertArasapas_act` dengan
**`<pyStepsClassName>ASM-FW-GISFW-Work`**.

Implementasi berclass `ASM-FW-GISFW-WORK` (232 KB) **tidak ada di modul ini** — hanya di
`EDM Treaty In` (`ab3ae3c9ba`) dan `Endorsment Fac In` (`7314b6c49c`), dan keduanya terdaftar
berkonflik (OQ-011 #415). Telusur **dihentikan**; varian modul lain **tidak dipinjam** → **OQ-025**.

`[terverifikasi]` Risiko salah tafsirnya **lebih kecil dari perkiraan D1**: kedua varian di korpus
berbeda **hanya pada 3 tag metadata Pega**; 20 langkah dan rule yang dirujuk identik. Tetapi
blocker prosedural tetap berdiri — modul ini tidak memuat implementasinya.

### 6.2 Kode dan guard

`[terverifikasi]` Kode tanpa arti terverifikasi: `IsApproved = 1`; `LetterNo = "TREATYINDEPTHEAD"`;
`.ClaimPaymentType` = `"Claim"`/`"Salvage"`; `.ClaimType` = `"XOL"`; `.BusinessCode != "40"`;
`param.isFOR` = `"EDM"`/`"POLICY"` → OQ-020, OQ-014.

`[terverifikasi]` **`OperatorID.pyTelephone` menyimpan kode peran** — `"TREATY1"`, `"SPVTREATY1"`
di modul ini → **OQ-027** (Tahap 5 menemukan dua nilai lagi: `TREATY2`, `SPVTREATY2`).

`[terverifikasi]` `isApproved` ada sebagai **dua tipe rule** dengan identitas sama
(`ASM-FW-GISFW-WORK!ISAPPROVED`): `When/isApproved.xml` dan `DecisionTable/isApproved.xml`.
Mana yang dipakai flow **tidak dapat ditentukan** → **OQ-026**.

`[terverifikasi]` Folder **`Claude outputs`** berisi `Struktur_MenuNBTreatyIn.xlsx` (275.860 B) —
berkas non-Pega di dalam korpus; **tidak dibaca** → OQ-054.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **12 file**,
`OperatorID.pyTelephone` **2 file**; `pyWorkPage.pxCreateOperator` memuat identifier orang di
`When/IsSPVCreate.xml`. Nilai nama orang **tidak disalin** → OQ-021, OQ-027.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-009**, **OQ-011**, **OQ-012**, **OQ-013**, **OQ-018**, **OQ-020**, **OQ-021**,
**OQ-022**, **OQ-023**, **OQ-024**, **OQ-025**, **OQ-026**, **OQ-027**, **OQ-028**, **OQ-054**,
**OQ-059**.
