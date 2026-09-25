# Telusur Flow — NB Treaty In

STEP D2, Tahap 1 konteks #1. Ditelusur 2026-09-13 dari korpus READ-ONLY `D:\XML\RNM_BRD\`.
Konvensi: `_METHOD.md`.

**Titik masuk:** `NB Treaty In/Flow/InputRealizationTreatyIn.xml`
→ `RULE-OBJ-FLOW` / `ASM-FW-GISFW-WORK` / `INPUTREALIZATIONTREATYIN`, `<pyStartActivity>Start1`.

**Catatan salinan:** Flow ini punya salinan **berisi identik** di
`NB FacIn/Flow/InputRealizationTreatyIn.xml` (identitas sama, hash ternormalisasi sama — lihat
`../inventory/_summary.md` OQ-006). Telusur ini dibaca dari **file NB Treaty In**.

**Catatan OQ-018:** seluruh pernyataan di bawah menyebut file yang dibaca. Korpus memuat hostname
DEV dan dirakit dari >1 server Pega, jadi ini **belum tentu cerminan production**.

---

## 1. Diagram alur

Graf dibaca dari `<pxSubscript>` (ID shape), `<pyShapeType>`, `<pyFrom>`/`<pyTo>`,
`<pyConditionType>`, dan `<pyExpression>`. Perintah audit ada di `_METHOD.md` §2.2.

Struktur: **6 Assignment, 12 Decision, 2 Utility, 30 connector** `[terverifikasi]`
(`grep -o "<pyShapeType>[^<]*" | sort | uniq -c`).

```
Start1
  └─(Always)─> Assignment2  "Input Realitation"   [WorkBasket]
                   │  FlowAction: InboxPolicyTreatyIn
                   v
               Decision3  "Is Correct?"  (When: isApproved)
                   ├─ Status=No  ──> End3   (selesai, tidak diproses lanjut)
                   └─ Status=YES ──> Decision6  "pxCreateOperator SPV?"
                                        ├─ When IsSPVCreate ──> Decision12 "IS SPV TREATY 1"
                                        │                          ├─ When IsSPVTreaty1 ─> Assignment6
                                        │                          └─ Else ─────────────> Assignment4
                                        └─ Else ─────────────────> Decision10 "is Treaty 1"
                                                                   ├─ When IsTreaty1 ───> Assignment4
                                                                   └─ Else ─────────────> Assignment6

  Assignment4 "Acceptance by Head. Treaty" [WorkBasket] ──FlowAction DeptHeadTreatyIn_UW──> Decision11
  Assignment6 "Acceptance by Head. Treaty" [WorkBasket] ──FlowAction DeptHeadTreatyIn_UW──> Decision4

  Decision11 "Is Correct?" ├─ Status=YES ─> Decision13
                           └─ Status=No  ─> Assignment2   (kembali ke input)
  Decision4  "Is Correct?" ├─ Status=YES ─> Decision13
                           └─ Status=No  ─> Assignment2   (kembali ke input)

  Decision13 "TO TREATY DEPT HEAD?"
        ├─ When ToTREATYDEPTHEAD ─> Assignment3 "Acceptance by Dept. Head" [WorkBasket]
        │                              └─FlowAction DeptHeadTreatyIn_UW─> Decision2 "Is Correct?"
        │                                     ├─ Status=YES ─> Decision8
        │                                     └─ Status=No  ─> Assignment2
        └─ Else ──────────────────> Decision8 "[Decision]"

  Decision8 ├─ When NopolisEmpty ─> Assignment3   (balik ke Dept. Head)
            └─ Else ─────────────> Utility1 "Save json policy"
                                       │ Activity: SaveJsonPolisTreatyIn_Act
                                       └─(ALWAYS)─> Utility2 "HIT SERVICE ARASAPAS"
                                                       │ Activity: serviceInsertArasapas_act
                                                       └─(ALWAYS)─> End3
```

### 1.1 Sub-graf tanpa connector masuk `[terverifikasi]`

Lima shape membentuk jalur yang **tidak dituju connector mana pun**:

```
Decision5 "Claim" ─When isClaimTreaty─> Assignment5 "Claim Manager" [WorkBasket]
      └─FlowAction DeptHeadTreatyIn_UW─> Decision7 "Is Correct?"
             ├─ Status=YES ─> Assignment1 "Acceptance by Dir." [WorkBasket]
             │                    └─FlowAction DeptHeadTreatyIn_UW─> Decision1 "Is Correct?"
             │                           ├─ Status=Yes ─> Decision8
             │                           └─ Status=No  ─> Assignment2
             └─ Status=No  ─> Assignment2
```

Perintah audit: dari daftar pasangan `<pyFrom>`/`<pyTo>`, tidak ada baris yang `pyTo = Decision5`.

`[pertanyaan terbuka]` Apakah jalur komite klaim + persetujuan Direktur ini **mati** (tidak pernah
dieksekusi), ataukah connector masuknya hilang saat ekspor? → **OQ-023**.

### 1.2 Routing

Keenam Assignment memakai `<pyImplementation>WorkBasket` dan `<pyRouteTo>Custom`
`[terverifikasi]`. Nama workbasket yang muncul di `<pyRuleName>` flow ini: `ReasTreatyInAdmin`,
`ReasTreatyInSecHead`, `ReasTreatyInGroupLeader`, `ReasTreatyInDeptHead`, `ReasTreatyInDirector`.

`[dugaan]` Kelimanya mewakili **peran** dalam alur akseptasi treaty inward. Pemetaan shape → nama
workbasket tidak dapat dipastikan dari tag yang terbaca (`<pyWorkBasket>` kosong, routing `Custom`)
→ **OQ-024**. Ini bukti paling konkret sejauh ini untuk **OQ-007** (tidak ada rule identitas/otorisasi
di korpus).

---

## 2. Status / state yang berubah

| Properti | Nilai literal | Diubah/diuji oleh | Arti |
| --- | --- | --- | --- |
| `pyWorkPage.PolicyTreatyIn.IsApproved` | `1` = lolos | `When isApproved` (`NB Treaty In/When/isApproved.xml`) | `[terverifikasi]` sebagai kondisi; arti nilai selain 1 **belum terverifikasi** |
| `pyWorkPage.PolicyTreatyIn.PolicyNo` | `""` (kosong) | `When NopolisEmpty` (`NB Treaty In/When/NopolisEmpty.xml`) | kosong → kembali ke Dept. Head |
| `pyWorkPage.LetterNo` | `"TREATYINDEPTHEAD"` | `When ToTREATYDEPTHEAD` (`NB Treaty In/When/ToTREATYDEPTHEAD.xml`) | `[dugaan]` token routing di field nomor surat; **arti belum terverifikasi** |
| `.PolicyTreatyIn.ClaimPaymentType` | `"Claim"`, `"Salvage"` | `When isClaimTreaty` (`NB Treaty In/When/isClaimTreaty.xml`) | **arti belum terverifikasi** |
| `.PolicyTreatyIn.ClaimType` | `"XOL"` | `When isClaimTreaty` | `XOL` **kepanjangan belum terverifikasi** |
| `.Quotation.BusinessCode` | `!= "40"` | precondition step 1 `SaveJsonPolisTreatyIn_Act` | **arti kode 40 belum terverifikasi** |
| `param.isFOR` | `"EDM"`, `"POLICY"` | precondition step 3 & 5 `SaveJsonPolisTreatyIn_Act` | membedakan dua mode simpan; `EDM` **kepanjangan belum terverifikasi** (OQ-014) |
| `OperatorID.pyTelephone` | `"TREATY1"`, `"SPVTREATY1"` | `When IsTreaty1`, `When IsSPVTreaty1` | **field telepon dipakai menyimpan kode peran** — lihat §7 |
| `pyWorkPage.pxCreateOperator` | dua identifier orang | `When IsSPVCreate` | guard identitas orang — lihat §7 |

Perintah audit kondisi `When`:
`grep -oE "<pyLabel>[^<]{1,160}" "NB Treaty In/When/<nama>.xml" | sort -u`

---

## 3. Objek Oracle yang disentuh

Hanya objek yang benar-benar tersentuh **oleh jalur yang ditelusur**.

| Objek | Operasi | Dipanggil dari | Path rule |
| --- | --- | --- | --- |
| `dual` + sequence `POOLDATA.JSON_POLIS_TREATYIN_SEQ` | SELECT / `NEXTVAL` | `RDB-List` step 4 `SaveJsonPolisTreatyIn_Act` | `NB Treaty In/RDBList/GenerateNoPolicy.xml` (`ASM-FW-GISFW-INT-POLISTREATYIN` / `ASM!GENERATENOPOLICY` / `RULE-CONNECT-SQL`) |
| — (blok PL/SQL, bukan tabel langsung) | CALL + `COMMIT` | `RDB-List` step 7 `SaveJsonPolisTreatyIn_Act` | `NB Treaty In/RDBList/SavePolisTreatyIn_SQL.xml` (`ASM-FW-GISFW-INT-POLISTREATYIN` / `ASM!SAVEPOLISTREATYIN_SQL` / `RULE-CONNECT-SQL`) |

### 3.1 Format nomor polis `[terverifikasi]`

`GenerateNoPolicy` membentuk nomor polis dari `dual`:

```sql
SELECT 'RNM-' || {InputData.CARI20} || '.T' || {pyWorkPage.PolicyTreatyIn.QuotationData.BusinessOldId}
       || '.' || to_char(sysdate,'MM.yyyy') || '.'
       || LPAD (POOLDATA.JSON_POLIS_TREATYIN_SEQ.NEXTVAL, 5, '0') AS "HASIL"
FROM dual
```

Komponen: prefix `RNM-`, `CARI20`, literal `.T`, `BusinessOldId`, bulan-tahun `MM.yyyy`, dan nomor
urut 5 digit dari sequence. **Isi `CARI20` tidak terlihat di activity ini** — ia diisi di tempat lain
`[pertanyaan terbuka]`.

---

## 4. Integrasi eksternal

| Shape | Activity | Status |
| --- | --- | --- |
| `Utility2` "HIT SERVICE ARASAPAS" | `serviceInsertArasapas_act` | **terblokir** — lihat §5.2 |

Tidak ada rule `ConnectREST` yang tersentuh **langsung** oleh jalur flow ini. Modul NB Treaty In
tidak memiliki rule `ConnectREST` sama sekali (`../inventory/NB Treaty In.md` §1).

`[dugaan]` Nama shape menyebut "ARASAPAS", dan D1 mencatat skema Oracle `ARASAPAS` dengan objek
`ARASAPAS.INVOICE` / `ARASAPAS.DETAIL_INVOICE` yang dirujuk modul lain. Keterkaitannya **belum
terverifikasi** — label shape bukan bukti perilaku (`_METHOD.md` §2.2).

---

## 5. Batas pengetahuan

### 5.1 Stored procedure — isinya tidak diketahui

Rule `ASM-FW-GISFW-INT-POLISTREATYIN` / `ASM!SAVEPOLISTREATYIN_SQL` / `RULE-CONNECT-SQL`
(`NB Treaty In/RDBList/SavePolisTreatyIn_SQL.xml`) memanggil `POOLDATA.PEGA_JSON_POLIS_TREATYIN`
dengan 8 parameter:

| # | Parameter | Sumber |
| ---: | --- | --- |
| 1 | `{pyWorkPage.pzInsKey}` | kunci case Pega |
| 2 | `{pyWorkPage.PolicyTreatyIn.PolicyNo}` | nomor polis |
| 3 | `NULL` | — |
| 4 | `'0'` | literal; **arti belum terverifikasi** |
| 5 | `{InputData.CARI21}` | **belum terlihat diisi di jalur ini** |
| 6 | `{OperatorID.pyUserIdentifier}` | identitas pengguna |
| 7 | `{InputData.CARI3}` | **JSON**, diisi `@ASM.GetPageJSONString()` (step 6) |
| 8 | `{OutputData.HASIL1 out}` | parameter keluaran |

Diikuti `COMMIT` **di dalam blok PL/SQL**. **Isi `POOLDATA.PEGA_JSON_POLIS_TREATYIN` tidak diketahui**
(OQ-002); batas transaksi berada di sisi database (OQ-013).

### 5.2 Rule berkonflik yang variannya tidak ada di modul ini — **BLOCKER**

`NB Treaty In/Activity/serviceInsertArasapas_act.xml`
(`ASM-FW-GISFW-DATA-POLICYTREATYIN` / `SERVICEINSERTARASAPAS_ACT` / `RULE-OBJ-ACTIVITY`) berisi
**satu langkah**: `Call serviceInsertArasapas_act`, dengan `<pyStepsObjectName>pyWorkPage`.

`pyWorkPage` di flow ini berclass `ASM-FW-GISFW-WORK`, sehingga `[dugaan]` target panggilan adalah
`ASM-FW-GISFW-WORK / SERVICEINSERTARASAPAS_ACT`. Identitas itu **terdaftar di register OQ-011**
(entri #415, 2 varian, **2 isi berbeda**) dan variannya hanya ada di
`EDM Treaty In/Activity/serviceInsertArasapas_act.xml` dan
`Endorsment Fac In/Activity/serviceInsertArasapas_act.xml` — **tidak ada di NB Treaty In**.

Sesuai `_METHOD.md` §1.1 aturan 4: telusur **dihentikan di sini**. Varian dari modul lain
**tidak dipinjam**. Apa yang dilakukan langkah "HIT SERVICE ARASAPAS" **tidak dapat dinyatakan**
sebelum OQ-011 dijawab. → **OQ-025**.

### 5.3 Fungsi Pega kustom

`@ASM.GetPageJSONString()` menghasilkan JSON yang ditulis ke database (parameter 7 di §5.1).
**Source fungsi ini tidak ada di korpus** — tidak ada tipe rule fungsi/library di 17 tipe yang
diinventarisasi (`../README.md` §6.1). Ini sumber langsung isi kolom JSON → memperkuat **OQ-012**.

### 5.4 `isApproved` ada sebagai dua tipe rule

`[terverifikasi]` Di NB Treaty In terdapat **dua** rule beridentitas class+nama sama:

| Path | Tipe |
| --- | --- |
| `NB Treaty In/When/isApproved.xml` | `RULE-OBJ-WHEN` |
| `NB Treaty In/DecisionTable/isApproved.xml` | `RULE-DECLARE-DECISIONTABLE` |

Keduanya `ASM-FW-GISFW-WORK!ISAPPROVED`. Enam shape Decision merujuk `isApproved`; **mana dari dua
tipe itu yang dipakai tidak dapat ditentukan dari flow** → **OQ-026**. Kondisi yang dicatat di §2
berasal dari varian `When`.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GISFW-WORK` / `INPUTREALIZATIONTREATYIN` / `RULE-OBJ-FLOW` | `Flow/InputRealizationTreatyIn.xml` | tidak |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` / `SAVEJSONPOLISTREATYIN_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SaveJsonPolisTreatyIn_Act.xml` | tidak |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` / `SERVICEINSERTARASAPAS_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/serviceInsertArasapas_act.xml` | tidak (tapi **memanggil** identitas yang berkonflik — §5.2) |
| `ASM-FW-GISFW-WORK` / `INBOXPOLICYTREATYIN` / `RULE-OBJ-FLOWACTION` | `FlowAction/InboxPolicyTreatyIn.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `DEPTHEADTREATYIN_UW` / `RULE-OBJ-FLOWACTION` | `FlowAction/DeptHeadTreatyIn_UW.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISAPPROVED` / `RULE-OBJ-WHEN` | `When/isApproved.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISAPPROVED` / `RULE-DECLARE-DECISIONTABLE` | `DecisionTable/isApproved.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `NOPOLISEMPTY` / `RULE-OBJ-WHEN` | `When/NopolisEmpty.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISSPVCREATE` / `RULE-OBJ-WHEN` | `When/IsSPVCreate.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISSPVTREATY1` / `RULE-OBJ-WHEN` | `When/IsSPVTreaty1.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISTREATY1` / `RULE-OBJ-WHEN` | `When/IsTreaty1.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `TOTREATYDEPTHEAD` / `RULE-OBJ-WHEN` | `When/ToTREATYDEPTHEAD.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISCLAIMTREATY` / `RULE-OBJ-WHEN` | `When/isClaimTreaty.xml` | tidak |
| `ASM-FW-GISFW-INT-POLISTREATYIN` / `ASM!GENERATENOPOLICY` / `RULE-CONNECT-SQL` | `RDBList/GenerateNoPolicy.xml` | tidak |
| `ASM-FW-GISFW-INT-POLISTREATYIN` / `ASM!SAVEPOLISTREATYIN_SQL` / `RULE-CONNECT-SQL` | `RDBList/SavePolisTreatyIn_SQL.xml` | tidak |

**15 rule ditelusur.** Seluruh path relatif terhadap `D:\XML\RNM_BRD\NB Treaty In\`.

Rule yang **dirujuk tapi belum ditelusur** (di luar jalur inti; kandidat telusur lanjutan):
`InboxPolicyTreatyIn_postDT`, `InputPolicyTreatyInPost_Act`, `DeptHeadTreatyInUW_preDT`,
`DeptHeadTreatyIn_UW_postDT`, Section `GeneralPolicyTreatyIn`, `GeneralDeptHeadTreatyIn_UW`,
Activity `CopyToPolicy`.

---

## 7. Pertanyaan terbuka baru

Seluruhnya ditambahkan ke `../open-questions.md`.

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-023** | Sub-graf komite klaim + Direktur tanpa connector masuk — jalur mati atau ekspor tidak lengkap? | Product+Underwriting |
| **OQ-024** | Pemetaan Assignment → workbasket (`ReasTreatyIn*`) tidak terbaca; routing `Custom` | IAM + Product+UW |
| **OQ-025** | `SERVICEINSERTARASAPAS_ACT` berkonflik dan variannya tidak ada di NB Treaty In (blocker) | Product+UW + pemilik export |
| **OQ-026** | `isApproved` ada sebagai `When` **dan** `DecisionTable` di class yang sama — mana yang dipakai flow? | Product+UW |
| **OQ-027** | `OperatorID.pyTelephone` dipakai menyimpan **kode peran** (`TREATY1`, `SPVTREATY1`) — 14 file di 7 modul | IAM |

### 7.1 Perluasan OQ-021 — dua pola guard identitas yang belum pernah disapu

`[terverifikasi]` Telusur ini menemukan **dua pola guard yang luput dari sapuan OQ-021 di D1**
(yang hanya mencari `OperatorID.pyUserIdentifier` / `pyUserName` / `pyPosition`):

| Pola | File | Modul | Nilai |
| --- | ---: | ---: | --- |
| `OperatorID.pyTelephone = "<kode>"` | 14 | 7 | `TREATY1` (6×), `SPVTREATY1` (6×) |
| `pyWorkPage.pxCreateOperator = "<orang>"` | 6 | 5 | 4 identifier orang (nilai tidak disalin) |

Perintah audit:

```
grep -rl "pyTelephone" "<modul>" --include="*.xml"
grep -rlE "pxCreateOperator\]\[&amp;#61;\]" "<modul>" --include="*.xml"
```

Yang pertama **bukan** sekadar hardcode identitas — ia memakai **field telepon operator sebagai
penyimpan kode peran**. `[dugaan]` ini workaround karena tidak ada model peran; **belum
terverifikasi**. Dampaknya ke desain RBAC berbeda dari OQ-021 biasa, sehingga dicatat sebagai
**OQ-027** tersendiri.
