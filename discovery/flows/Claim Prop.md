# Telusur Flow — Claim Prop

STEP D2, Tahap 3 konteks #2. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `Claim Prop/Flow/Flow_TreatyIn.xml`
→ `RULE-OBJ-FLOW` / **`ASM-FW-GCNMFW-WORK-CLAIMTREATY`** / `FLOW_TREATYIN`,
`<pyStartActivity>Start1`.

> **Nama file sama, rule berbeda.** `Claim Non Prop` punya `Flow/Flow_TreatyIn.xml` dengan nama
> sama, class `ASM-FW-GCNMFW-WORK-CLAIMTREATY**NONPROP**`. **Dua rule berbeda**, ditelusur terpisah.
> Hash ternormalisasi: Prop = `9871faee8d`, Non Prop = `b899c224ab`.

---

## 1. Diagram alur

**2 Assignment, 1 Decision, 1 Event-Start, 4 connector** `[terverifikasi]`.

```
Start1 ─> Assignment2 "Outstanding Claim"  [WorkList, route=Current operator]
             │ FlowAction: OutstandingClaim
             v
          Assignment1 "Input Acceptation"  [WorkBasket, route=Custom]
             │ FlowAction: InputAcceptation
             v
          Decision3 "IsBack"
             ├─ When IsBackStage ─> Assignment2   (balik ke Outstanding)
             └─ Else ────────────> End1
```

`[terverifikasi]` `Claim Prop/When/IsBackStage.xml`: `.pyNote = "Back"` — nilai catatan dipakai
sebagai gerbang alur.

**Routing campuran** seperti Claim Life: tahap Outstanding ke pengguna sendiri (`WorkList` +
`Current operator`), tahap Akseptasi ke antrean bersama (`WorkBasket` + `Custom`) → OQ-028.

`[terverifikasi]` Graf ini **hanya memuat dua tahap**. Register, Estimasi, Adjustment, Close, dan
Reject **tidak muncul sebagai shape** — seluruhnya Activity di luar graf → OQ-039.

---

## 2. Status / state dan kode

| Kode | Nilai literal yang diuji | Arti |
| --- | --- | --- |
| `PaymentType` | **`1`, `2`, `3`, `4`, `5`, `6`** | **belum terverifikasi** |
| `TransferType` | `2` | **belum terverifikasi** |
| `.pyNote` | `"Back"` | gerbang `IsBackStage` |

Perintah audit:
```
grep -rhoE "(PaymentType|TransferType)[ ]*[=!]+[ ]*[\"']?[0-9]{1,2}" "Claim Prop" --include="*.xml" | sort | uniq -c
```

Frekuensi tertinggi: `PaymentType=1` (9×), `=4` (8×), `=2` (8×), `==1` (6×), `==3` (5×).

**Mapping `PaymentType` ke jenis pembayaran tidak ditebak.** Bukti Tahap 2 menunjukkan mapping
status ke setuju/tolak bisa berlawanan dengan dugaan nama
(`KomitePost_Reject` justru digerbangi `AcceptStatus=="1"`).

---

## 3. Jembatan ke Komite — dengan gerbang limit

`[terverifikasi]` `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml`
(`ASM-FW-GCNMFW-DATA-ADJUSTMENT` / `ADDKOMITETREATYCHILD_ACT`, 420.247 byte) membuat **case anak
Komite** dari sisi Claim.

Ini menjawab sebagian **OQ-039**: penyerahan ke tangga Komite dilakukan oleh Activity, bukan oleh
shape di flow.

### 3.1 Limit wewenang: ter-hardcode **dan** dari database

`[terverifikasi]` Modul ini punya **tiga rule Connect-SQL pengambil limit**:

| Rule | Path |
| --- | --- |
| `GetLimitDirekturUtama_SQL` | `RDBList/GetLimitDirekturUtama_SQL.xml` |
| `GetLimitPLATreatyin` | `RDBList/GetLimitPLATreatyin.xml` |
| `GetLimitsTreatyIn_SQL` | `RDBList/GetLimitsTreatyIn_SQL.xml` |

Jadi **limit diambil dari database** di sini — berbeda dari Claim Non Prop yang **meng-hardcode**
nilainya (lihat `Claim Non Prop.md` §3.1) dan dari `ApprovalKomite_Act` di Komite Claim FacIn yang
juga meng-hardcode pita nominal (OQ-037).

`[terverifikasi]` **Korpus memuat dua cara menentukan limit wewenang untuk proses yang setara.**
Nama rule `GetLimitDirekturUtama_SQL` `[dugaan]` menyiratkan limit khusus Direktur Utama;
**belum terverifikasi**. → **OQ-040**.

---

## 4. Objek Oracle yang disentuh

| Objek | Rule perujuk |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 |
| `TREATYINPRODUCTION` | 3 |
| `JSON_POLIS` | 3 |
| `BANKACCOUNT` | 3 |
| `AGENT` | 3 |
| `POOLDATA.DIRECTTOKASIR_LOG` | 2 |
| **`OS_AKSEPTASI_KLAIM`** | 2 |
| `M_CLIENT` | 2 |
| `V_D_CAUSE_OF_LOSS_BUSINESS`, `TREATYYEAR` | 1 masing-masing |

`[terverifikasi]` **`OS_AKSEPTASI_KLAIM` juga disentuh Komite Claim Prop** (lewat `SaveOSClaim_SQL`,
Tahap 2 §5). Seperti pada pasangan Claim Life ↔ Komite Claim Life, tabel akseptasi adalah **titik
temu kedua sisi** — dicatat di sini, **tidak disimpulkan ulang**.

Perintah audit: filter dataset ekstraksi D1 pada `tables=` untuk path `Claim Prop/`.

---

## 5. Batas pengetahuan

### 5.1 Stored procedure — isinya tidak diketahui

`[terverifikasi]` Delapan procedure dipanggil dari rule Connect-SQL modul ini:

```
POOLDATA.GENERATE_NOCLMTREATYIN     POOLDATA.GENERATE_NOCLMTRTYINTEMP
POOLDATA.PEGA_D_CAUSE_OF_LOSS       POOLDATA.PEGA_M_CAUSE_OF_LOSS
POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM   POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP
POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT  POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER
```

**Isi seluruhnya tidak ada di korpus** (OQ-002). Akhiran `TNP` dan `TRT` **kepanjangan belum
terverifikasi**.

### 5.2 `ISCLMP` dan kerabatnya — berkonflik **dan** tidak terbaca

`[terverifikasi]` `@BASECLASS` / `ISCLMP` / `RULE-OBJ-WHEN` terdaftar di register **OQ-011 entri
#293 — 6 varian, 6 isi berbeda**. Varian modul ini (`Claim Prop/When/IsCLMP.xml`, hash `4919199b`)
yang relevan; varian 5 modul lain **tidak dibaca**.

**Masalah berlapis:** kondisinya **juga tidak terbaca** — `<pyLabel>` hanya template kosong
`[first value][relation][second value]`. Jadi kita tahu keenam varian berbeda, tetapi **tidak bisa
melihat bedanya di mana**.

Hal yang sama berlaku untuk `ISCLM` (`2245b132`), `ISCLMNP` (`50c584af`), `ISPEGASYARIAH`
(`0a7077fd`) di modul ini. → **OQ-041**.

### 5.3 Activity besar yang belum habis dibaca

| Activity | Ukuran |
| --- | ---: |
| `SaveOutstanding_Act.xml` | 535.235 byte |
| `HitServiceToKasir_Act.xml` | 482.194 byte |
| `AddKomiteTreatyChild_ACT.xml` | 420.247 byte |
| `CountEstimation_Act.xml` | 392.613 byte |
| `PrintDLATreatyIn.xml` | 380.551 byte |

**Seluruhnya belum ditelusur isinya** — hanya identitas dan ukuran yang dicatat.

### 5.4 `IsPEGAPROD`

Ada di modul ini; kondisinya **tidak terbaca**; cabang aktif **tidak ditebak** → OQ-029.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Claim Prop/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY` / `FLOW_TREATYIN` / `RULE-OBJ-FLOW` | `Flow/Flow_TreatyIn.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY` / `ISBACKSTAGE` / `RULE-OBJ-WHEN` | `When/IsBackStage.xml` | tidak |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` / `ADDKOMITETREATYCHILD_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/AddKomiteTreatyChild_ACT.xml` | tidak |
| (FlowAction) `OUTSTANDINGCLAIM`, `INPUTACCEPTATION` | `FlowAction/` | tidak |
| `@BASECLASS` / `ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` / `RULE-OBJ-WHEN` | `When/` | **YA — #291, #293, #292, #306** |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `When/IsPEGAPROD.xml` | **YA — #305** |
| 3 rule limit + 8 procedure (§3.1, §5.1) | `RDBList/` | sebagian |

**14 rule ditelusur.**

---

## 7. Pertanyaan terbuka baru

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-040** | Limit wewenang diambil dari DB di Claim Prop, tetapi ter-hardcode di Claim Non Prop dan Komite FacIn — mana yang berlaku? | Finance + Product+UW |
| **OQ-041** | `ISCLM`/`ISCLMP`/`ISCLMNP`/`ISPEGASYARIAH` berkonflik 6-versi **dan** kondisinya tidak terbaca dari tag | Product+UW + pemilik export |

OQ dikuatkan: OQ-002 (§5.1), OQ-011 (§5.2), OQ-020 (§2), OQ-028 (§1), OQ-029 (§5.4), OQ-039 (§3).
