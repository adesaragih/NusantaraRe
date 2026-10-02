# Modul — Claim Life

STEP D3. Sintesis dari `../inventory/Claim Life.md` (D1) + `../flows/Claim Life.md` (D2 Tahap 3).
**136 file** — modul Claim terkecil.
Audit: `find "Claim Life" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **siklus kerugian untuk lini life**. Class kerja
`ASM-FW-GCNMFW-WORK-CLAIMLIFE` (44 rule). Titik masuk `Claim Life/Flow/Register_Flow.xml`
→ `RULE-OBJ-FLOW` / `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW`, `<pyStartActivity>Start2`.

`[terverifikasi]` **Modul Claim paling terpisah dari yang lain**: overlap identitas hanya
**12–14 %** terhadap Claim Prop / Non Prop / Fac In (`../flows/_SUMMARY-claim.md` §9).

## 2. Proses / fitur utama

Rincian: **`../flows/Claim Life.md`** (16 rule ditelusur).

`[terverifikasi]` Tahapan: Register → Outstanding → **Medical Check** → Claim Analis.
**Satu-satunya modul Claim yang punya tahap Medical Check.**

`[pertanyaan terbuka]` Gerbang `Claim Life/When/IsSendtoMedical.xml` — kondisinya **tidak terbaca**
dari tag → OQ-029.

`[terverifikasi]` Routing campuran (input `WorkList`, tahap penilaian `WorkBasket`) → OQ-028.

`[terverifikasi]` **Tidak memakai `PaymentType` sama sekali** — berbeda dari ketiga modul Claim lain
yang memakainya (`../flows/_SUMMARY-claim.md` §4).

`[terverifikasi]` Memakai kode lini bisnis **`L1`–`L11`** dalam satu precondition → OQ-038.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 51 Activity, 29 Connect-SQL, 20 Section, 16 FlowAction,
8 ReportDefinition, 3 When, 3 Harness, 2 ConnectREST, 1 DataTransform, 1 Flow, 1 DecisionTable,
1 SystemSettings.

Objek Oracle:

| Objek | Rule perujuk |
| --- | ---: |
| **`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`** | 6 |
| `T_STORAGE_IMAGE`, `POOLDATA.M_LIFE_PREMIUM_DETAIL` | 3 masing-masing |
| `M_PRODUCT_LIFE` | 2 |
| `TREATYYEAR_LIFE`, `RETROCESSIONLIFE`, `PRODUCT_LIFE`, `POOLDATA.TANGGAL_CLOSING`, `POOLDATA.RATE_LIFE`, `POOLDATA.T_FOLDER_IMAGE` | 1 masing-masing |

Audit:
```
awk -F'\t' '$1 ~ "^Claim Life/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
 | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

`[terverifikasi]` Seluruh objeknya bersufiks/berkonteks **`_LIFE`** — tidak menyentuh tabel klaim
non-life (`OS_AKSEPTASI_KLAIM` tanpa sufiks, `JSON_KLAIM`).

## 4. Integrasi eksternal

`[terverifikasi]` Dua `RULE-CONNECT-REST`: `ServiceGoogle` dan
**`convertJsonNusareToProductionClaimLife`** (varian life dari `convertJsonNusareToProduction`;
hanya di modul ini). Keduanya `pyBaseURLSelectionType = SETTING` → `LinkService!LinkService`,
tanpa URL literal.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Komite Claim Life** | `UpdateOsAkseptasiClaimLife_sql` adalah **satu rule dengan identitas dan hash ternormalisasi sama (`c50bfd9a12`)** di kedua modul, dan **tidak terdaftar di register OQ-011** → satu Connect-SQL bersama, bukan dua penulis independen | `[terverifikasi]` |
| **Master Product Name Life** | membaca `M_PRODUCT_LIFE`, `PRODUCT_LIFE` | `[terverifikasi]` |
| **Master Contract Retro Life** | membaca `RETROCESSIONLIFE`, `TREATYYEAR_LIFE` | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Tiga stored procedure** (paling sedikit di domain Claim), isinya tidak ada di
korpus (OQ-002): `POOLDATA.GET_TOKEN_STORAGE`, `POOLDATA.PEGA_JSON_KLAIM_PNC`,
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` Activity besar yang **belum habis dibaca**: `SaveOutStandingLife_Act.xml`
(± 643 KB) — batas **cakupan telusur**.

`[terverifikasi]` **Kode tanpa arti terverifikasi**: `STS_REJECT` bernilai `0`/`1`/`2`;
`L1`–`L11` → OQ-020, OQ-038.

`[terverifikasi]` **`ISCLM*` tidak menyentuh modul ini** — keempat rule berkonflik itu ada di Claim
Prop, Non Prop, dan Fac In, bukan di sini (OQ-041).

`[terverifikasi]` `Claim Life/When/IsPEGAPROD.xml` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **2 file**,
`OperatorID.pyTelephone` **0**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-018**, **OQ-019**, **OQ-020**, **OQ-021**, **OQ-028**, **OQ-029**, **OQ-038**,
**OQ-039**.
