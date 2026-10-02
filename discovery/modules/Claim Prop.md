# Modul — Claim Prop

STEP D3. Sintesis dari `../inventory/Claim Prop.md` (D1) + `../flows/Claim Prop.md` (D2 Tahap 3).
**270 file**.
Audit: `find "Claim Prop" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **siklus kerugian untuk treaty inward proporsional**. Class kerja
`ASM-FW-GCNMFW-WORK-CLAIMTREATY` (103 rule). Titik masuk `Claim Prop/Flow/Flow_TreatyIn.xml`
→ `RULE-OBJ-FLOW` / `ASM-FW-GCNMFW-WORK-CLAIMTREATY` / `FLOW_TREATYIN`,
`<pyStartActivity>Start1`, hash ternormalisasi **`9871faee8d`**.

`[terverifikasi]` **Bukan salinan** `Claim Non Prop`: nama file Flow sama persis, tetapi class
berbeda (`…CLAIMTREATYNONPROP`) dan hash berbeda (`b899c224ab`) → dua rule berbeda,
ditelusur terpisah.

## 2. Proses / fitur utama

Rincian: **`../flows/Claim Prop.md`** (14 rule ditelusur).

`[terverifikasi]` Tahapan graf hanya **dua**: Outstanding Claim → Input Acceptation, dengan jalur
mundur `When IsBackStage` (`.pyNote = "Back"`).

`[terverifikasi]` Adjustment, Close, dan Reject **tidak muncul sebagai shape** — seluruhnya
Activity di luar graf → OQ-039.

`[terverifikasi]` Routing campuran (`Assignment2` Outstanding = `WorkList`/`Current operator`;
`Assignment1` Input Acceptation = `WorkBasket`/`Custom`) → OQ-028.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 112 Activity, 54 Connect-SQL, 36 Section,
19 ReportDefinition, 12 FlowAction, 11 DataTransform, 11 Harness, 6 When, 6 ConnectREST, 1 Flow,
1 DecisionTable, 1 SystemSettings.

Objek Oracle:

| Objek | Rule perujuk |
| --- | ---: |
| `T_STORAGE_IMAGE`, `TREATYINPRODUCTION`, `JSON_POLIS`, `BANKACCOUNT`, `AGENT` | 3 masing-masing |
| `POOLDATA.DIRECTTOKASIR_LOG`, `OS_AKSEPTASI_KLAIM`, `M_CLIENT` | 2 masing-masing |
| `V_D_CAUSE_OF_LOSS_BUSINESS`, `TREATYYEAR` | 1 masing-masing |

Audit:
```
awk -F'\t' '$1 ~ "^Claim Prop/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
 | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

## 4. Integrasi eksternal

`[terverifikasi]` **6 rule `RULE-CONNECT-REST`**: `GetDtlPaymentClaim`, `KonversiKlaimNonLife`,
`SendAcceptationToKasir`, `ServiceGoogle`, `getPayAttachment`, `getPremiumPaidOnTreatyIn`.
Seluruhnya `pyBaseURLSelectionType = SETTING` → `LinkService!LinkService`, tanpa URL literal.

`[terverifikasi]` `POOLDATA.DIRECTTOKASIR_LOG` (2 rule) — jejak integrasi Kasir di sisi database.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Komite Claim Prop** | `OS_AKSEPTASI_KLAIM` (2 rule) juga disentuh modul Komite pasangannya; jembatan lewat `AddKomiteTreatyChild_ACT` (`ASM-FW-GCNMFW-DATA-ADJUSTMENT`, 420.247 byte) | `[terverifikasi]` keberadaan |
| **Treaty In (master)** | membaca `TREATYINPRODUCTION` (3 rule), `TREATYYEAR` | `[terverifikasi]` |
| **Claim Non Prop** | berbagi **102 identitas** (38 % dari 270) — basis kode sebagian, bukan ruleset bersama | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Delapan stored procedure — terbanyak di domain Claim**, isinya tidak ada di
korpus (OQ-002): `POOLDATA.GENERATE_NOCLMTREATYIN`, `POOLDATA.GENERATE_NOCLMTRTYINTEMP`,
`POOLDATA.PEGA_D_CAUSE_OF_LOSS`, `POOLDATA.PEGA_M_CAUSE_OF_LOSS`,
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`, `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`,
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT`, `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.
Ditambah `GL.F_GET_EMAIL`, `POOLDATA.GETCURRENCYSTANDARD`, `POOLDATA.GET_TOKEN_STORAGE`,
`POOLDATA.PEGA_JSON_KLAIM_PNC`, `DBMS_LOB.CREATETEMPORARY`.

`[terverifikasi]` **Limit wewenang diambil dari database** di modul ini —
`GetLimitDirekturUtama_SQL`, `GetLimitPLATreatyin`, `GetLimitsTreatyIn_SQL` — berbeda dari
`Claim Non Prop` yang meng-hardcode nilainya → OQ-040.

`[terverifikasi]` Lima activity > 380 KB **belum habis dibaca** — batas **cakupan telusur**.

`[terverifikasi]` **Kode tanpa arti terverifikasi**: `PaymentType` `1`–`6`, `TransferType` `2`
→ OQ-020.

`[terverifikasi]` `ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` — varian modul ini berbeda hash dari
varian Claim Non Prop dan Claim Fac In, **dan kondisinya tidak terbaca** → OQ-041.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **7 file**,
`OperatorID.pyTelephone` **0**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-011**, **OQ-018**, **OQ-019**, **OQ-020**, **OQ-021**, **OQ-028**, **OQ-029**,
**OQ-039**, **OQ-040**, **OQ-041**.
