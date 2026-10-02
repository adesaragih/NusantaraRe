# Telusur Flow — Claim Non Prop

STEP D2, Tahap 3 konteks #3. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `Claim Non Prop/Flow/Flow_TreatyIn.xml`
→ `RULE-OBJ-FLOW` / **`ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP`** / `FLOW_TREATYIN`,
`<pyStartActivity>Start1`.

> **Nama file sama, rule berbeda.** `Claim Prop` punya `Flow/Flow_TreatyIn.xml` dengan nama sama,
> class `ASM-FW-GCNMFW-WORK-CLAIMTREATY` (tanpa `NONPROP`). **Dua rule berbeda**, ditelusur
> terpisah. Hash ternormalisasi: Non Prop = `b899c224ab`, Prop = `9871faee8d`.

---

## 1. Diagram alur

**2 Assignment, 1 Decision, 1 Event-Start, 4 connector** `[terverifikasi]` — bentuk graf **sama
persis** dengan Claim Prop, tetapi rule-nya berbeda.

```
Start1 ─> Assignment2 "Outstanding Claim"  [WorkList, route=Current operator]
             │ FlowAction: OutstandingClaim
             v
          Assignment1 "Input Acceptation"  [WorkBasket, route=Custom]
             │ FlowAction: InputAcceptation
             v
          Decision3 "IsBack"
             ├─ When IsBackStage ─> Assignment2
             └─ Else ────────────> End1
```

`Claim Non Prop/When/IsBackStage.xml`: `.pyNote = "Back"` — **teks kondisi sama** dengan Claim Prop.

Seperti Claim Prop, graf hanya memuat dua tahap; sisanya Activity di luar graf → OQ-039.

---

## 2. Status / state dan kode

| Kode | Nilai literal | Arti |
| --- | --- | --- |
| `PaymentType` | **`1`, `2`, `3`, `4`, `5`, `6`, `7`** | **belum terverifikasi** |
| `TransferType` | `2` | **belum terverifikasi** |
| `.pyNote` | `"Back"` | gerbang `IsBackStage` |

`[terverifikasi]` Modul ini memakai **nilai `7`**, yang tidak muncul di Claim Prop. Frekuensi:
`PaymentType=="3"` (5×), `!=7` (3×), `=6`/`=5`/`=4`/`=2` (2× masing-masing).

Perintah audit sama dengan Claim Prop, path diganti.

---

## 3. Jembatan ke Komite — **limit ter-hardcode**

`[terverifikasi]` `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml`
(`ASM-FW-GCNMFW-DATA-ADJUSTMENT` / `CREATECHILDKOMITECNP_ACT`, **756.836 byte** — activity terbesar
yang ditemui D2 sejauh ini) membuat **case anak Komite**.

Langkah awal yang terbaca: `Page-Remove`, `Property-Set`, `Obj-Refresh-And-Lock`, `Page-New` ×3,
`Page-Remove`, lalu serangkaian `Property-Set`, dengan gerbang limit.

### 3.1 Nilai limit ter-hardcode `[terverifikasi]`

| Variabel | Nilai ter-hardcode |
| --- | ---: |
| `Local.LimitMax` | **30000000.00** |
| `Local.LimitMaxDivHead` | **50000000.00** |
| `Local.LimitPersenMax` | **30.00** |

Gerbang yang memakainya:

```
Local.TotalValueAdjust <= Local.LimitMax
Local.TotalValueAdjust >  Local.LimitMax && Local.TotalValueAdjust <= Local.LimitMa…
Local.TotalValueAdjust <= Local.LimitMaxDivHead
```

`[terverifikasi]` Batas bawah **30.000.000,00 sama persis** dengan ambang di
`Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` (OQ-037), tetapi batas atasnya berbeda
(50.000.000,00 di sini vs 57.750.000,00 di sana).

`[terverifikasi]` Modul ini juga punya rule pengambil limit dari DB
(`RDBList/GetLimitsTreatyIn_SQL.xml`, `RDBList/GetLimitTONPPLA.xml`) — jadi **dua mekanisme limit
hidup berdampingan di modul yang sama**.

Mata uang **tidak disebut**; arti `LimitPersenMax = 30.00` (persen dari apa) **belum
terverifikasi**. `DivHead` **kepanjangan belum terverifikasi**. → **OQ-040**.

---

## 4. Objek Oracle — termasuk **kebocoran batas ke Treaty Out**

| Objek | Rule perujuk |
| --- | ---: |
| `BANKACCOUNT` | 5 |
| `T_STORAGE_IMAGE`, `TREATYINPRODUCTION`, `AGENT` | 3 masing-masing |
| `POOLDATA.TREATYINDETAILEDM` | 3 |
| **`OS_AKSEPTASI_KLAIM`** | 3 |
| `POOLDATA.TREATYINDETAIL`, `M_CLIENT`, `JSON_POLIS`, `JSON_KLAIM` | 2 masing-masing |
| **`M_TREATY_OUT`**, **`M_TREATY_OUT_DETAIL`**, **`TREATY_OUT`** | 1 masing-masing |
| `TREATYCONTRACT`, `TREATYBUSINESS`, `V_D_CAUSE_OF_LOSS_BUSINESS` | 1 masing-masing |

### 4.1 Recovery lewat master Treaty **Out** `[terverifikasi]`

Tiga rule Connect-SQL modul klaim ini membaca objek **treaty outward**:

| Rule | Objek |
| --- | --- |
| `RDBList/GetDataMasterTOutNP.xml` | `M_TREATY_OUT` |
| `RDBList/BrowseDtlTreatyOutNP.xml` | `M_TREATY_OUT_DETAIL` |
| `RDBList/GetLimitTONPPLA.xml` | `TREATY_OUT` |

`[dugaan]` ini jalur **recovery / retrosesi**: klaim non-proporsional memulihkan sebagian nilai
lewat treaty outward. **Belum terverifikasi.**

**Ini kebocoran batas konteks yang terukur**: modul domain **klaim** membaca master domain
**treaty outward**. Perlu dicatat untuk D4 — **tetapi penetapan konteks bukan di sini**.

Catatan silang: D1 §21.3 menemukan modul bernama `Treaty Contract Out` justru **tidak** merujuk
objek treaty outward sama sekali (OQ-022), sementara `Claim Non Prop` **merujuk**. → memperkuat
OQ-022.

### 4.2 Titik temu dengan Komite

`OS_AKSEPTASI_KLAIM` (3 rule) juga disentuh `Komite Claim Non Prop` lewat `InsertOSKlaimCNP`
(Tahap 2). Dicatat sebagai titik temu; **tidak disimpulkan ulang**.

---

## 5. Batas pengetahuan

### 5.1 Activity besar yang **belum habis dibaca**

`[terverifikasi]` Modul ini memuat activity terbesar di domain klaim. Yang dicatat hanya identitas,
ukuran, dan (untuk satu activity) langkah awal:

| Activity | Ukuran | Terbaca |
| --- | ---: | --- |
| `CreateChildKomiteCNP_Act.xml` | 756.836 | 12 langkah awal + nilai limit |
| `CountLossAllocation_act.xml` | 635.029 | **belum** |
| `SaveDataToOSAksep_Act.xml` | 611.417 | **belum** |
| `GenerateCACNP_Act.xml` | 533.564 | **belum** |
| `HitServiceToKasir_Act.xml` | 482.194 | **belum** |
| `GenerateCFS_act.xml` | 419.291 | **belum** |
| `SetPayableTreatyNP_Act.xml` | 396.656 | **belum** |
| `AdjClaimCNP_Act.xml` | 394.446 | **belum** |

Akhiran `CNP`, `CA`, `CFS`, `TNP` **kepanjangan belum terverifikasi**.

**Perhitungan non-proporsional (alokasi kerugian, XOL) berada di dalam activity ini dan belum
ditelusur** — bukan karena lewat stored procedure, melainkan karena volumenya. Ini **batas
cakupan telusur**, bukan batas pengetahuan korpus.

### 5.2 Stored procedure — batas pengetahuan sesungguhnya

`[terverifikasi]` Empat procedure dipanggil:

```
POOLDATA.PEGA_D_CAUSE_OF_LOSS        POOLDATA.PEGA_M_CAUSE_OF_LOSS
POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP  POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER
```

**Isinya tidak ada di korpus** (OQ-002). Lebih sedikit dari Claim Prop (8 procedure).

Activity khusus XOL/CNP di modul Komite pasangannya (`InsertXOLKlaimCNP`, `InsertOSKlaimCNP`,
`GenerateAccCNP_act`) sudah dicatat di Tahap 2 sebagai belum ditelusur.

### 5.3 `ISCLMNP` dan kerabatnya — berkonflik dan tidak terbaca

`[terverifikasi]` Varian modul ini: `ISCLM` (`541179a0`), `ISCLMP` (`9b10dc97`),
`ISCLMNP` (`d3255d2c`), `ISPEGASYARIAH` (`1bbf8269`) — seluruhnya **berbeda** dari varian Claim Prop
dan Claim Fac In. Kondisinya **tidak terbaca** dari tag → OQ-041.

Varian modul lain **tidak dibaca**.

### 5.4 `IsPEGAPROD`

Ada; kondisinya tidak terbaca; cabang aktif **tidak ditebak** → OQ-029.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Claim Non Prop/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` / `FLOW_TREATYIN` / `RULE-OBJ-FLOW` | `Flow/Flow_TreatyIn.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` / `ISBACKSTAGE` / `RULE-OBJ-WHEN` | `When/IsBackStage.xml` | tidak |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` / `CREATECHILDKOMITECNP_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/CreateChildKomiteCNP_Act.xml` | tidak |
| (FlowAction) `OUTSTANDINGCLAIM`, `INPUTACCEPTATION` | `FlowAction/` | tidak |
| `@BASECLASS` / `ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` | `When/` | **YA — #291, #293, #292, #306** |
| `@BASECLASS` / `ISPEGAPROD` | `When/IsPEGAPROD.xml` | **YA — #305** |
| `GetDataMasterTOutNP`, `BrowseDtlTreatyOutNP`, `GetLimitTONPPLA`, `GetLimitsTreatyIn_SQL` | `RDBList/` | tidak |

**15 rule ditelusur.**

---

## 7. Pertanyaan terbuka baru

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-042** | Klaim non-proporsional membaca master **treaty outward** (`M_TREATY_OUT`, `M_TREATY_OUT_DETAIL`, `TREATY_OUT`) — jalur recovery/retrosesi? | Product+UW |

OQ dikuatkan: **OQ-040** (§3.1 — limit ter-hardcode vs DB, dalam satu modul), OQ-041 (§5.3),
OQ-002 (§5.2), OQ-011, OQ-020 (§2 — nilai `7`), OQ-022 (§4.1), OQ-028, OQ-029, OQ-037 (§3.1 —
batas bawah 30 jt sama), OQ-039 (§3).

**Guard identitas** `[terverifikasi]`: `pyUserIdentifier` **0 file**, `pyPosition` **2 file**.
