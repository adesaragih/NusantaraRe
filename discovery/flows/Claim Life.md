# Telusur Flow — Claim Life

STEP D2, Tahap 3 konteks #1. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `Claim Life/Flow/Register_Flow.xml`
→ `RULE-OBJ-FLOW` / **`ASM-FW-GCNMFW-WORK-CLAIMLIFE`** / `REGISTER_FLOW`,
`<pyStartActivity>Start2`.

> **Nama file sama, rule berbeda.** `Claim Fac In` punya `Flow/Register_Flow.xml` dengan nama sama,
> class `ASM-FW-GCNMFW-WORK-PNC`. **Dua rule berbeda**, ditelusur terpisah. Dokumen ini hanya
> varian Claim Life.

---

## 1. Diagram alur — siklus kerugian

**4 Assignment, 3 Decision, 1 Event-Start, 11 connector** `[terverifikasi]`.

```
Start2 ─> Assignment2 "Input Register"     [WorkList, route=Current operator]
             │ FlowAction: InputRegisterClaimLife
             v
          Assignment1 "Outstanding Claim"  [WorkList, route=Custom]
             │ FlowAction: OSClaimLife
             v
          Decision3
             ├─ When IsSendtoAdmin ─> Assignment2   (balik ke Register)
             └─ Else ──────────────> Assignment3 "Medical Check" [WorkBasket, Custom]
                                        │ FlowAction: MedicalCheck
                                        v
                                     Decision1
                                        ├─ When IsSendtoAdmin ─> Assignment1 (Outstanding)
                                        └─ Else ──────────────> Assignment4 "Claim Analis" [WorkBasket, Custom]
                                                                   │ FlowAction: AkseptasiClaimLife
                                                                   v
                                                                Decision2
                                                                   ├─ When IsSendtoAdmin   ─> Assignment1
                                                                   ├─ When IsSendtoMedical ─> Assignment3
                                                                   └─ Else ────────────────> End1
```

**Tahapan siklus `[terverifikasi]`:** Register → Outstanding → Medical Check → Claim Analis
(Akseptasi) → selesai. Dua jalur kembali: `IsSendtoAdmin` (dari tiga titik) dan `IsSendtoMedical`
(dari tahap akseptasi).

`[terverifikasi]` **Tidak ada tahap Adjustment maupun Close/Reject sebagai shape** di flow ini —
keduanya ada sebagai Activity (`SaveAdjustment_Act`, `ProtectCloseClaim_act`,
`RejectOSClaimLife_Act`) tetapi **di luar graf flow**. `[pertanyaan terbuka]` bagaimana dipicu.

### 1.1 Routing campuran — pertama kali di D2

`[terverifikasi]` Modul ini memakai **dua model penugasan sekaligus**:

| Shape | `pyImplementation` | `pyRouteTo` |
| --- | --- | --- |
| `Assignment2` "Input Register" | `WorkList` | **`Current operator`** |
| `Assignment1` "Outstanding Claim" | `WorkList` | `Custom` |
| `Assignment3` "Medical Check" | **`WorkBasket`** | `Custom` |
| `Assignment4` "Claim Analis" | **`WorkBasket`** | `Custom` |

Sebelumnya D2 menemukan modul **seragam**: treaty inward semuanya `WorkBasket`; Life, PremiumList,
dan seluruh Komite semuanya `WorkList`. **Claim Life campuran** — tahap awal ke pengguna sendiri,
tahap penilaian (medis, analis) ke antrean bersama. → memperkaya **OQ-028**.

---

## 2. Status / state yang berubah

| Properti | Nilai literal | Peran | Arti |
| --- | --- | --- | --- |
| `pyWorkPage.SendtoAdmin` | `1` | guard `IsSendtoAdmin` — memicu balik ke tahap sebelumnya | **belum terverifikasi** |
| `STS_REJECT` | `'0'`, `'1'`, `'2'` | dipakai 20+ kali sebagai gate | **belum terverifikasi** |
| `AcceptStatus` | `2` | 1 kemunculan | **belum terverifikasi** |
| `.IsAccept` | `"true"` | precondition `Property-Set` | — |
| `.Protect` | `@contains(.Protect,"1")` | gate `Page-Set-Messages` | **belum terverifikasi** |
| `BusinessCode` | **`L1`…`L11`** | satu precondition menguji **sebelas kode** berderet | **belum terverifikasi** |
| `ContentNote` | `"DEATH"` | gate `RDB-List` + `Property-Set` | `[dugaan]` sebab klaim |

Perintah audit:
```
grep -rhoE "(AcceptStatus|STS_REJECT|SendtoAdmin)[ ]*[=!]+[ ]*[\"']?[0-9A-Za-z]{1,10}" "Claim Life" --include="*.xml" | sort | uniq -c
```

`[terverifikasi]` **`PaymentType` tidak dipakai sama sekali di Claim Life** — berbeda dari ketiga
modul Claim lain yang memakainya (§`_SUMMARY-claim.md`).

### 2.1 Sebelas kode lini bisnis ter-hardcode

`Claim Life/Activity/SaveOutStandingLife_Act.xml` memuat satu precondition yang menguji
`pyWorkPage.PolicyDataLife.BusinessCode` terhadap **`"L1"`, `"L2"`, … `"L11"`** berderet dengan `||`.

`[dugaan]` daftar lini bisnis Life. **Arti tiap kode belum terverifikasi**, dan daftar ter-hardcode
seperti ini berarti penambahan lini memerlukan perubahan rule → **OQ-038**.

### 2.2 `IsSendtoMedical` — kondisi tidak terbaca

`[terverifikasi]` `Claim Life/When/IsSendtoMedical.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL`) — `<pyLabel>` hanya berisi template kosong
`[first value][relation][second value]`, sama seperti `IsPEGAPROD` (OQ-029).

**Apa yang diuji tidak diketahui**, padahal ia menentukan apakah kasus dikembalikan ke tahap medis.
Ini **kasus kedua** pola When tak terbaca → ditambahkan ke **OQ-029**.

Sebaliknya `When/IsSendtoAdmin.xml` terbaca: `pyWorkPage.SendtoAdmin = 1`.

---

## 3. Objek Oracle yang disentuh

Dari `Claim Life/Activity/SaveOutStandingLife_Act.xml` (**642.787 byte**):

| Class | RequestType |
| --- | --- |
| `ASM-FW-GISFW-Int-policyjson` | `GETTanggalClosing_SQL` |
| `ASM-FW-GCNMFW-Work-ClaimLife` | `CountPesertaAkseptasiLife_SQL` |
| `ASM-FW-GCNMFW-Work-ClaimLife` | `CountPesertaAkseptasiLifeHealth_SQL` |
| `ASM-FW-GCNMFW-Work-ClaimLife` | `GetCategoryLife_SQL` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` | `Generate_NoKlaim_Life` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` | `Generate_NoKlaim_LifeRetro` |
| `ASM-FW-GISFW-Int-policyjson` | `GetKodeProdLife_SQL` |
| `ASM-FW-GISFW-Int-policyjson` | `GetSequenceNumber_SQL` |

`[terverifikasi]` Penomoran klaim **bercabang retro / non-retro** (`Generate_NoKlaim_Life` vs
`..._LifeRetro`) — pola sejajar dengan penomoran **akseptasi** di Komite Claim Life
(`Generate_NoAccept_KMT_Life` vs `..._LifeRetro`).

### 3.1 Operasi tulis yang terbaca langsung

| Rule | Operasi |
| --- | --- |
| `RDBList/UpdateDateClaimLife_SQL.xml` | `UPDATE POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` |
| `RDBList/InsertLogServiceClaim.xml` | `INSERT INTO pooldata.monitoring_klaim_log` |
| `RDBList/Insert_T_Storage_SQL.xml` | `insert into t_storage_image` |
| `RDBList/Update_T_Storage_SQL.xml` | `update T_STORAGE_IMAGE` |
| `RDBList/GetSequenceNumber_SQL.xml` | memanggil `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER(` → **batas pengetahuan** |

### 3.2 Tabel akseptasi dibagi dengan Komite — **satu rule, dua pemanggil**

`[terverifikasi]` Ini klarifikasi penting terhadap temuan Tahap 2.

`RDBList/UpdateOsAkseptasiClaimLife_sql.xml` ada di **Claim Life** *dan* **Komite Claim Life**
dengan:

| | Identitas | Hash ternormalisasi |
| --- | --- | --- |
| Claim Life | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` | `c50bfd9a12` |
| Komite Claim Life | identitas **sama** | `c50bfd9a12` — **identik** |

dan **tidak terdaftar di register OQ-011**.

**Jadi bukan dua penulis independen ke `OS_AKSEPTASI_KLAIM_LIFE`, melainkan satu rule Connect-SQL
yang sama dipanggil dari dua sisi** (siklus klaim dan tangga komite). 55 kolom yang ditulis sudah
terdaftar di `Komite Claim Life.md` §3.1 — **tidak diulang di sini** untuk menghindari kesimpulan
ganda.

---

## 4. Integrasi eksternal

Modul ini punya `Activity/serviceInsertArasapasClaimLife_act.xml` — **satu-satunya salinan di
korpus**, dan justru inilah yang dipanggil `Komite Claim Life` (OQ-035). Juga
`Activity/SendEmailKlaimLF.xml`, `Activity/GetLinkService.xml`, `Activity/InsertGoogleStorage_Act.xml`,
`Activity/GetUrlGoogleStorage_Act.xml`, dan 2 rule `ConnectREST`.

**Belum ditelusur** — di luar jalur inti flow.

`Claim Life/When/IsPEGAPROD.xml` ada (varian grup yang sama dengan Endorsement Life). Kondisinya
**tidak terbaca**; cabang aktif **tidak ditebak** → OQ-029.

---

## 5. Batas pengetahuan

### 5.1 Konflik OQ-011 — `ISCLM*` **tidak** menyentuh modul ini

`[terverifikasi]` Keempat `When` berkonflik 6-versi (`ISCLM`, `ISCLMP`, `ISCLMNP`,
`ISPEGASYARIAH`) berada di **6 modul**: Claim Fac In, Claim Non Prop, Claim Prop, Komite Claim
FacIn, Komite Claim Non Prop, Komite Claim Prop — **bukan Claim Life**.

```
find . -iname "IsCLM.xml" -not -path "./OUTPUT_HASIL_RNM/*"
```

Jadi telusur Claim Life **tidak tersentuh** konflik itu. Modul lain di Tahap 3 akan tersentuh.

### 5.2 Activity besar yang belum habis dibaca

| Activity | Ukuran | Status |
| --- | ---: | --- |
| `SaveOutStandingLife_Act.xml` | **642.787 byte** | 16 langkah pertama terbaca; **sisanya belum** |

Pembacaan tag menampilkan sebagian; jumlah langkah seluruhnya dan cabang lain **belum terukur**.

### 5.3 Stored procedure

`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` dipanggil dari `GetSequenceNumber_SQL`. **Isinya tidak
diketahui** (OQ-002).

### 5.4 Tahap yang tidak muncul di graf

`SaveAdjustment_Act`, `ProtectCloseClaim_act`, `RejectOSClaimLife_Act`, `CreateKMTLife_Act`,
`GetListKomiteLife` ada sebagai Activity tetapi **tidak dirujuk** flow. Bagaimana Adjustment,
Close, Reject, dan penyerahan ke Komite dipicu **tidak diketahui dari flow ini** → **OQ-039**.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Claim Life/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `Flow/Register_Flow.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SaveOutStandingLife_Act.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN` | `When/IsSendtoAdmin.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN` | `When/IsSendtoMedical.xml` | tidak (kondisi tak terbaca) |
| (FlowAction) `INPUTREGISTERCLAIMLIFE`, `OSCLAIMLIFE`, `MEDICALCHECK`, `AKSEPTASICLAIMLIFE` | `FlowAction/` | tidak |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | tidak — **identik** dengan salinan Komite |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `When/IsPEGAPROD.xml` | **YA — #305** |
| 8+ rule `RULE-CONNECT-SQL` (§3) | `RDBList/` | sebagian |

**16 rule ditelusur.**

---

## 7. Pertanyaan terbuka baru

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-038** | Sebelas kode lini bisnis `L1`…`L11` ter-hardcode dalam satu precondition | Product+UW |
| **OQ-039** | Adjustment, Close, Reject, dan penyerahan ke Komite tidak muncul di graf flow — bagaimana dipicu? | Product+UW |

**Guard identitas** `[terverifikasi]`: `pyUserIdentifier` **0 file**, `pyTelephone` **0 file**,
**`pyPosition` 8 file** — jumlah `pyPosition` tertinggi dari seluruh modul yang ditelusur D2 sejauh
ini. Nilainya tercatat di D1 (`ReasLifeAdmin`, `ReasLifeSPV`, `ReasLifeMedicalAdvisor`,
`IT Developer`). → OQ-021.

OQ dikuatkan: OQ-020 (§2), OQ-028 (§1.1 — routing campuran), OQ-029 (§2.2 — When kedua tak terbaca),
OQ-002 (§5.3), OQ-035 (§4).
