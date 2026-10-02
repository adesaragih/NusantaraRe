# Telusur Flow — Claim Fac In

STEP D2, Tahap 3 konteks #4. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.
Modul klaim terbesar: **482 file**.

**Titik masuk:** `Claim Fac In/Flow/Register_Flow.xml`
→ `RULE-OBJ-FLOW` / **`ASM-FW-GCNMFW-WORK-PNC`** / `REGISTER_FLOW`, `<pyStartActivity>Start1`.

> **Nama file sama, rule berbeda.** `Claim Life` punya `Flow/Register_Flow.xml` dengan nama sama,
> class `ASM-FW-GCNMFW-WORK-CLAIMLIFE`, dan `<pyStartActivity>Start2` (bukan `Start1`).
> **Dua rule berbeda**, ditelusur terpisah.

---

## 1. Diagram alur

**3 Assignment, 3 Decision, 1 Event-Start, 9 connector** `[terverifikasi]`.

```
Start1 ─> Decision4 "B2B"
             ├─ When IsSPK ─> Assignment3 "Choose Surveyor"  [WorkBasket, Custom]
             └─ Else ──────> Assignment1 "Input Register"    [WorkList, Current operator]
                                │ FlowAction: InputRegister
                                v
                             Assignment7 "Input Estimasi"    [WorkList, Custom]
                                │ FlowAction: InputEstimasi
                                v
                             Decision5 "IsBack"
                                ├─ When IsBackStage ─> Assignment1  (balik ke Register)
                                └─ Else ────────────> Assignment3 "Choose Surveyor"
                                                         │ FlowAction: InputSurveyor
                                                         v
                                                      Decision8 "IsBack"
                                                         ├─ When IsBackStage ─> Assignment7 (balik ke Estimasi)
                                                         └─ Else ────────────> END52
```

**Tahapan siklus `[terverifikasi]`:** (percabangan B2B) → Register → Estimasi → Choose Surveyor →
selesai, dengan dua jalur mundur lewat `IsBackStage`.

### 1.1 Percabangan di titik masuk — khas modul ini

`[terverifikasi]` Satu-satunya modul Claim yang **bercabang sebelum tahap pertama**. `Decision4`
berlabel **"B2B"**, digerbangi `When IsSPK`: bila benar, alur **melewati Register dan Estimasi**
dan langsung ke Choose Surveyor.

`[dugaan]` jalur business-to-business (klaim masuk dari mitra) tidak memerlukan input manual.
**Belum terverifikasi** — label shape bukan bukti perilaku (`_METHOD.md` §2.2).

`SPK` **kepanjangan belum terverifikasi**.

### 1.2 Routing campuran

| Shape | `pyImplementation` | `pyRouteTo` |
| --- | --- | --- |
| `Assignment1` "Input Register" | `WorkList` | `Current operator` |
| `Assignment7` "Input Estimasi" | `WorkList` | `Custom` |
| `Assignment3` "Choose Surveyor" | **`WorkBasket`** | `Custom` |

Pola sama dengan Claim Life, Claim Prop, dan Claim Non Prop → OQ-028.

---

## 2. Status / state dan kode

| Kode | Nilai literal | Arti |
| --- | --- | --- |
| `PaymentType` | **`1`–`7`** | **belum terverifikasi** |
| `TransferType` | `2` | **belum terverifikasi** |
| `.pyNote` | `"Back"` | gerbang `IsBackStage` |

`[terverifikasi]` **Pemakaian `PaymentType` paling padat dari keempat modul Claim**:
`=4` 23×, `=6` 21×, `==3` 19×, `=2` 18×, `=="3"` 17×, `=1` 17×, `!=3` 17×, `==4` 15×.

Perintah audit:
```
grep -rhoE "(PaymentType|TransferType)[ ]*[=!]+[ ]*[\"']?[0-9]{1,2}" "Claim Fac In" --include="*.xml" | sed 's/ //g' | sort | uniq -c | sort -rn
```

### 2.1 `IsSPK` — kondisi tidak terbaca

`[terverifikasi]` `Claim Fac In/When/IsSPK.xml` (`ASM-FW-GCNMFW-WORK` / `ISSPK`) — `<pyLabel>` hanya
template kosong `[first value][relation][second value]`.

**Ini kasus ketiga** pola When tak terbaca di D2 (setelah `IsPEGAPROD` dan
`Claim Life/When/IsSendtoMedical.xml`), dan yang paling berdampak: ia menggerbangi **percabangan di
titik masuk flow**. Apa yang memisahkan jalur B2B dari jalur normal **tidak diketahui** → OQ-029.

Perhatikan class-nya `ASM-FW-GCNMFW-WORK` (class induk), bukan `...WORK-PNC` seperti `IsBackStage`.

---

## 3. Objek Oracle yang disentuh

| Objek | Rule perujuk |
| --- | ---: |
| **`OS_AKSEPTASI_KLAIM`** | 6 |
| `JSON_POLIS` | 5 |
| `FACINPRODUCTION` | 4 |
| `T_STORAGE_IMAGE` | 3 |
| **`REINSURANCE.TRLOSS_DETAIL_T`** | 3 |
| `PROPORTIONALARRG` | 3 |
| `JSON_KLAIM` | 3 |
| `AGENT` | 3 |
| `TREATYCONTRACT`, `TREATYBUSINESS`, `POOLDATA.DIRECTTOKASIR_LOG`, `M_CLIENT` | 2 masing-masing |

`[terverifikasi]` **Rujukan terbanyak ke `OS_AKSEPTASI_KLAIM`** dari keempat modul Claim (6 rule) —
tabel yang sama yang disentuh `Komite Claim FacIn` lewat `SaveOSClaim_SQL` (Tahap 2). Titik temu
Claim ↔ Komite; **tidak disimpulkan ulang**.

### 3.1 Kebocoran batas ke domain Treaty

`[terverifikasi]` Modul **klaim fakultatif** membaca objek domain **treaty**: `TREATYCONTRACT`,
`TREATYBUSINESS`, `PROPORTIONALARRG`. Ditambah activity
`Claim Fac In/Activity/DLAFacintoTreaty_Act.xml`
(`ASM-FW-GCNMFW-DATA-OBJECT` / `DLAFACINTOTREATY_ACT`, 644.505 byte) yang namanya menyebut
jembatan Fac → Treaty.

Pola sejajar dengan `Claim Non Prop` yang membaca `M_TREATY_OUT` (OQ-042). **Dicatat untuk D4;
penetapan konteks bukan di sini.** `DLA` **kepanjangan belum terverifikasi**.

`[terverifikasi]` Skema **`REINSURANCE`** dirujuk (`REINSURANCE.TRLOSS_DETAIL_T`) — salah satu dari
11 skema di OQ-016.

---

## 4. Integrasi eksternal

`[terverifikasi]` Modul dengan `ConnectREST` terbanyak di korpus: **8 rule**.

| `pyServiceName` | Sumber base URL | URL literal? |
| --- | --- | --- |
| `HitDLAClaimFacin` | `SETTING` | tidak |
| `KonversiKlaimNonLife` | `SETTING` | tidak |
| `SendAcceptationToKasir` | `SETTING` | tidak |
| `ServiceGoogle` | `SETTING` | tidak |
| `getPayAttachment` | `SETTING` | tidak |
| `getPaymentClaim` | `SETTING` | tidak |
| `getPremiumPaidOn` | `SETTING` | tidak |
| **`getPremiumPaidOnMarine`** | **`URL`** | **YA** |

Tujuh dari delapan memakai konfigurasi `LinkService!LinkService`. Yang kedelapan memuat **URL
literal** — sudah tercatat di D1 sebagai satu-satunya penyimpangan `ConnectREST` batch 2, dan URL-nya
memuat contoh nilai data. Nilainya **tidak disalin** ke artefak ini → OQ-018.

`Claim Fac In/When/IsPEGAPROD.xml` ada; kondisinya tidak terbaca → OQ-029.

---

## 5. Batas pengetahuan

### 5.1 Stored procedure

`[terverifikasi]` Lima procedure dipanggil:

```
POOLDATA.PEGA_D_CAUSE_OF_LOSS   POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM
POOLDATA.PEGA_PROGRESSCLAIM     POOLDATA.PEGA_SUBPROGRESSCLAIM
POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER
```

**Isinya tidak ada di korpus** (OQ-002). `PEGA_PROGRESSCLAIM` dan `PEGA_SUBPROGRESSCLAIM` hanya
muncul di modul ini `[terverifikasi]`.

### 5.2 Activity besar yang **belum habis dibaca**

| Activity | Ukuran |
| --- | ---: |
| `CLaimFaceSheet_Act.xml` | **804.401** — terbesar yang ditemui di seluruh D2 |
| `GenerateDLAFacin_Act.xml` | 690.099 |
| `DLAFacintoTreaty_Act.xml` | 644.505 |
| `DraftGenerateDLAFacin_Act.xml` | 630.950 |
| `CheckLimitSpreadingTreaty_Act.xml` | 602.168 |

**Seluruhnya belum ditelusur isinya.** `CheckLimitSpreadingTreaty_Act` `[dugaan]` memuat
pemeriksaan limit spreading — relevan ke OQ-040, **belum terverifikasi**.

### 5.3 `ISCLM*` — berkonflik dan tidak terbaca

Varian modul ini: `ISCLM` (`4ee647f1`), `ISCLMP` (`a44d0c98`), `ISCLMNP` (`91c9b037`),
`ISPEGASYARIAH` (`cac64576`) — berbeda dari varian Claim Prop dan Claim Non Prop. Kondisi
**tidak terbaca** → OQ-041. Varian modul lain **tidak dibaca**.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Claim Fac In/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-PNC` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `Flow/Register_Flow.xml` | tidak |
| `ASM-FW-GCNMFW-WORK` / `ISSPK` / `RULE-OBJ-WHEN` | `When/IsSPK.xml` | tidak (kondisi tak terbaca) |
| `ASM-FW-GCNMFW-WORK-PNC` / `ISBACKSTAGE` / `RULE-OBJ-WHEN` | `When/IsBackStage.xml` | tidak |
| `ASM-FW-GCNMFW-DATA-OBJECT` / `DLAFACINTOTREATY_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/DLAFacintoTreaty_Act.xml` | tidak |
| (FlowAction) `INPUTREGISTER`, `INPUTESTIMASI`, `INPUTSURVEYOR` | `FlowAction/` | tidak |
| `@BASECLASS` / `ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` | `When/` | **YA — #291, #293, #292, #306** |
| `@BASECLASS` / `ISPEGAPROD` | `When/IsPEGAPROD.xml` | **YA — #305** |
| 8 rule `RULE-CONNECT-REST` (§4) | `ConnectREST/` | sebagian (`SERVICEGOOGLE` berkonflik) |

**17 rule ditelusur.**

---

## 7. Pertanyaan terbuka

Tidak ada OQ baru eksklusif. Dikuatkan: OQ-002 (§5.1), OQ-011 (§5.3), OQ-016 (§3.1 — skema
`REINSURANCE`), OQ-018 (§4), OQ-020 (§2), OQ-028 (§1.2), OQ-029 (§2.1 — `IsSPK` kasus ketiga),
OQ-039, OQ-042 (§3.1 — kebocoran batas Fac → Treaty).

**Guard identitas** `[terverifikasi]`: `pyUserIdentifier` **0 file**, `pyTelephone` **0 file**,
`pyPosition` **2 file**.

`[pertanyaan terbuka]` baru dicatat ke OQ-029: `IsSPK` menggerbangi percabangan **di titik masuk
flow** dan kondisinya tidak terbaca — dampaknya lebih besar dari dua kasus sebelumnya.
