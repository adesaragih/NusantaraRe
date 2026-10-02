# Modul — Claim Non Prop

STEP D3. Sintesis dari `../inventory/Claim Non Prop.md` (D1) + `../flows/Claim Non Prop.md`
(D2 Tahap 3). **279 file**.
Audit: `find "Claim Non Prop" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **siklus kerugian untuk treaty inward non-proporsional**. Class kerja
`ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` (84 rule). Titik masuk
`Claim Non Prop/Flow/Flow_TreatyIn.xml` → `RULE-OBJ-FLOW` /
`ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` / `FLOW_TREATYIN`, `<pyStartActivity>Start1`,
hash **`b899c224ab`**.

`[terverifikasi]` **Bukan salinan** `Claim Prop` — nama file Flow sama, class dan hash berbeda
(`9871faee8d`). Bentuk grafnya sama persis (2 Assignment, 1 Decision, 4 connector), tetapi rule-nya
berbeda.

## 2. Proses / fitur utama

Rincian: **`../flows/Claim Non Prop.md`** (15 rule ditelusur).

`[terverifikasi]` Tahapan graf: Outstanding Claim → Input Acceptation, jalur mundur
`When IsBackStage` (`.pyNote = "Back"` — **teks kondisi sama** dengan Claim Prop, rule berbeda).

`[terverifikasi]` Adjustment, Close, Reject **di luar graf** → OQ-039.

### 2.1 Jembatan ke Komite dengan limit ter-hardcode

`[terverifikasi]` `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml`
(`ASM-FW-GCNMFW-DATA-ADJUSTMENT` / `CREATECHILDKOMITECNP_ACT`, **756.836 byte**) membuat case anak
Komite dan memuat nilai ter-hardcode:

| Variabel | Nilai |
| --- | ---: |
| `Local.LimitMax` | **30000000.00** |
| `Local.LimitMaxDivHead` | **50000000.00** |
| `Local.LimitPersenMax` | **30.00** |

`[terverifikasi]` Batas bawah **30.000.000,00 sama persis** dengan ambang di
`Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` (OQ-037), tetapi batas atasnya berbeda
(50 jt di sini vs 57,75 jt di sana). **Mata uang tidak disebut**; arti `LimitPersenMax = 30.00`
**belum terverifikasi**. `DivHead` kepanjangan belum terverifikasi. → OQ-040.

`[terverifikasi]` Modul ini **sekaligus** punya rule pengambil limit dari DB
(`RDBList/GetLimitsTreatyIn_SQL.xml`, `RDBList/GetLimitTONPPLA.xml`) — **dua mekanisme limit hidup
berdampingan dalam satu modul**.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 114 Activity, 50 Connect-SQL, 36 Section,
22 ReportDefinition, 16 FlowAction, 14 Harness, 10 DataTransform, 7 When, 7 ConnectREST, 1 Flow,
1 DecisionTable, 1 SystemSettings.

| Objek | Rule perujuk |
| --- | ---: |
| `BANKACCOUNT` | 5 |
| `T_STORAGE_IMAGE`, `TREATYINPRODUCTION`, `POOLDATA.TREATYINDETAILEDM`, `OS_AKSEPTASI_KLAIM`, `AGENT` | 3 masing-masing |
| `POOLDATA.TREATYINDETAIL`, `M_CLIENT`, `JSON_POLIS`, `JSON_KLAIM` | 2 masing-masing |
| **`M_TREATY_OUT`**, **`M_TREATY_OUT_DETAIL`**, **`TREATY_OUT`**, `TREATYCONTRACT`, `TREATYBUSINESS`, `V_D_CAUSE_OF_LOSS_BUSINESS` | 1 masing-masing |

## 4. Integrasi eksternal

`[terverifikasi]` **7 rule `RULE-CONNECT-REST`**: `InsertClaimOutstanding_NP`,
`KonversiKlaimNonLife`, `SendAcceptationToKasir`, `ServiceGoogle`, `getPayAttachment`,
`getPremiumPaidOnTreatyIn`, `insertClaimFinalOrClosed_NP`. Seluruhnya `SETTING` →
`LinkService!LinkService`, tanpa URL literal.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Komite Claim Non Prop** | `OS_AKSEPTASI_KLAIM` juga disentuh lewat `InsertOSKlaimCNP`; jembatan `CreateChildKomiteCNP_Act` | `[terverifikasi]` |
| **Treaty inward (master)** | `TREATYINPRODUCTION`, `POOLDATA.TREATYINDETAIL`, `TREATYINDETAILEDM` | `[terverifikasi]` |
| **Treaty OUTWARD** — kebocoran batas | 3 rule Connect-SQL membaca objek outward: `RDBList/GetDataMasterTOutNP.xml` → `M_TREATY_OUT`; `RDBList/BrowseDtlTreatyOutNP.xml` → `M_TREATY_OUT_DETAIL`; `RDBList/GetLimitTONPPLA.xml` → `TREATY_OUT` | `[terverifikasi]` keberadaan / `[dugaan]` jalur recovery/retrosesi — **OQ-042** |
| **Claim Prop** | berbagi 102 identitas (38 %) | `[terverifikasi]` |

`[terverifikasi]` **Catatan silang penting:** modul bernama `Treaty Contract Out` justru **tidak**
menyentuh satu pun objek treaty outward (0 file) — yang menyentuhnya adalah modul ini (4 file),
`NB Treaty In` (8), `EDM Treaty In` (4), `Treaty In Adjustment` (2), `Claim Prop` (1) → OQ-022,
OQ-042.

Audit:
```
for m in "Claim Non Prop" "Treaty Contract Out" "NB Treaty In"; do
  echo "$m: $(grep -rli 'M_TREATY_OUT\|TREATY_OUT\|TREATYOUTDETAIL' "$m" --include='*.xml' | wc -l)"
done
```

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Empat stored procedure** khas + pendukung, isinya tidak ada di korpus (OQ-002):
`POOLDATA.PEGA_D_CAUSE_OF_LOSS`, `POOLDATA.PEGA_M_CAUSE_OF_LOSS`,
**`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`**, `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`;
ditambah `GL.F_GET_EMAIL`, `POOLDATA.GETCURRENCYSTANDARD`, `POOLDATA.GET_TOKEN_STORAGE`,
`POOLDATA.PEGA_JSON_KLAIM_PNC`, `DBMS_LOB.CREATETEMPORARY`.

`[terverifikasi]` **Delapan activity > 390 KB belum habis dibaca** — termasuk
`CreateChildKomiteCNP_Act.xml` (756.836 B, hanya 12 langkah awal + nilai limit yang terbaca),
`CountLossAllocation_act.xml` (635.029), `SaveDataToOSAksep_Act.xml` (611.417),
`GenerateCACNP_Act.xml` (533.564), `HitServiceToKasir_Act.xml` (482.194).

**Perhitungan alokasi kerugian non-proporsional (XOL) berada di dalam activity ini dan belum
ditelusur** — batas **cakupan telusur**, bukan batas korpus.

Akhiran `CNP`, `CA`, `CFS`, `TNP` **kepanjangan belum terverifikasi**.

`[terverifikasi]` **Kode tanpa arti terverifikasi**: `PaymentType` `1`–`7` (memakai nilai `7` yang
tidak muncul di Claim Prop), `TransferType` `2` → OQ-020.

`[terverifikasi]` `ISCLM` (`541179a0`), `ISCLMP` (`9b10dc97`), `ISCLMNP` (`d3255d2c`),
`ISPEGASYARIAH` (`1bbf8269`) — seluruhnya berbeda dari varian modul lain, **dan kondisinya tidak
terbaca** → OQ-041.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **7 file**,
`OperatorID.pyTelephone` **0**, `pyPosition` **2 file**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-011**, **OQ-018**, **OQ-019**, **OQ-020**, **OQ-021**, **OQ-022**, **OQ-028**,
**OQ-029**, **OQ-039**, **OQ-040**, **OQ-041**, **OQ-042**.
