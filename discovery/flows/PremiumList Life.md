# Telusur Flow — PremiumList Life

STEP D2, Tahap 1 konteks #4. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `PremiumList Life/InputPolicyHolder.xml`
→ `RULE-OBJ-FLOW` / `ASM-FW-GISFW-WORK-LIFE` / `INPUTPOLICYHOLDER`, `<pyStartActivity>Start1`.

**Catatan letak file:** rule `Flow` ini berada **langsung di root modul**, bukan di folder `Flow/`.
Tipenya sudah diverifikasi di D1 batch 4 dari `<pzOriginalInstanceKey>` **dan** `<pxObjClass>`
(`Rule-Obj-Flow`) — lihat **OQ-004 (terjawab)** di `../open-questions.md`. Telusur ini memakai tipe
yang sudah terverifikasi itu, bukan menebak dari nama.

---

## 1. Diagram alur

**3 Assignment, 3 Decision, 1 Utility, 11 connector** `[terverifikasi]`.

```
Start1 ─(Always)─> Assignment2 "Input Offer"  [WorkList]
                        │ FlowAction: InputDataOfferLife
                        v
                    Decision1 "Accept"  (DecisionTable: IsLifeAccepted)
                        ├─ Status=Decline ─> End1
                        └─ Status=Confirm ─> Decision3 "FlagOnGoingPolicy"
                                              (DecisionTable: IsFlagOnGoingPolicy)
                              ├─ Status=Offer   ─> END52          (selesai di tahap penawaran)
                              └─ Status=Premium ─> ASSIGNMENT63 "Input Premium List Detail" [WorkList]
                                                      │ FlowAction: ShowLifePremiumDetail
                                                      v
                                                  Decision2 "Accept" (IsLifeAccepted)
                                                      ├─ Status=Decline ─> End1
                                                      ├─ Status=Reject  ─> Assignment2  (balik ke Input Offer)
                                                      └─ Status=Confirm ─> Utility1
                                                                             │ Activity: InsertJsonPolisLife_Act
                                                                             └─> END52
```

`[terverifikasi]` Nilai status connector yang muncul: `Confirm`, `Decline`, `Reject`, `Offer`,
`Premium`. **Arti masing-masing belum terverifikasi** selain dari posisinya di graf.

### 1.1 Shape tanpa connector masuk `[terverifikasi]`

```
Assignment1 "Input Premium List Summary" [WorkList] ──FlowAction ShowLifePremiumSummary──> END52
```

Tidak ada pasangan `<pyFrom>`/`<pyTo>` yang menuju `Assignment1`. Ini **kejadian ketiga** pola yang
sama (NB Treaty In, Endorsement Life, PremiumList Life) → memperkuat **OQ-023**.

Layar "Premium List Summary" jelas bagian dari proses (ada Section dan FlowAction-nya), tetapi
**tidak terhubung ke graf flow**. `[pertanyaan terbuka]` apakah dicapai lewat jalan lain (menu,
harness) atau memang mati.

### 1.2 Routing

Ketiga Assignment memakai `<pyImplementation>WorkList` — sama dengan Endorsement Life, berbeda
dari kedua konteks treaty yang memakai `WorkBasket` (OQ-028).

---

## 2. Status / state yang berubah

Dari `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml`
(`ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY`, **21 step**):

| Properti / kondisi | Nilai literal | Peran | Arti |
| --- | --- | --- | --- |
| `Local.currentdate` | `@toDecimal(...) > 25` | precondition step 4 | **ambang tanggal 25** — muncul di 3 konteks (lihat §2.2) |
| `.ID` | `@contains(.ID,"1000032")`, `"1000033"`, `"1000034"`, `"1000035"` | precondition step 7–10 | **empat ID literal** menggerbangi `Property-Set` berbeda; **arti belum terverifikasi** |
| `IsPEGAPROD` | When guard | precondition step 18, 19, 21 | menggerbangi **email, service eksternal, dan Connect-REST** — §5.1 |
| `.Type` | `QR`, `QP`, `TP`, `TR` | lihat §2.1 | **arti belum terverifikasi** (OQ-020) |

### 2.1 Kode `.Type` — QR / QP / TP / TR `[terverifikasi]`

Dipakai di **11 rule** di modul ini:

```
grep -rlE "[\"'](QR|QP|TP|TR)[\"']" "PremiumList Life" --include="*.xml"
```

| Tipe rule | File |
| --- | --- |
| Activity | `Calculate1_Act.xml`, `GetPLNumber_Act.xml`, `ProtectAccept.xml`, `SavePremiumList_Act.xml`, `SubmitPremiumList_Act.xml`, `ValidasiUploadPL_act.xml`, `WPCLife_Act.xml` |
| DataTransform | `AppendCurrencySummary_DT.xml` |
| Section | `PL_Detail_Sec.xml`, `ShowLifePremiumDetail.xml`, `ShowLifePremiumSummary.xml` |

**Di tag apa kode itu muncul** `[terverifikasi]` — menentukan apa yang digerbanginya:

| Tag | Contoh | Yang digerbangi |
| --- | --- | --- |
| `<pyStepsPreCondParamsWhen>` | `.Type == "TR"` | **eksekusi step** di Activity |
| `<pyCondition>` | `.Type = 'QR'`, `.Type == 'TR' \|\| .Type == 'TP'` | kondisi rule |
| `<pyContainerVisibleWhen>` | `.Type = 'QP'`, `.Type = 'QR'` | **visibilitas bagian layar** |
| `<PropertiesValue>` | `"QR"`, `"QP"`, `"TP"`, `"TR"` | nilai yang **di-set** ke properti |

`[terverifikasi]` Pengelompokan yang muncul berulang: **`TR` dan `TP` diperlakukan berpasangan**
(`.Type=='TR' || .Type=='TP'`, tiga variasi penulisan), sedangkan `QR` dan `QP` muncul sebagai
kondisi tunggal terpisah.

`[dugaan]` Ada dua keluarga kode: huruf pertama `Q` vs `T`, huruf kedua `R` vs `P`.
**Arti keempat kode belum terverifikasi** — korpus tidak memuat tabel kode. → **OQ-020**.

Perbandingan lintas modul `[terverifikasi]`: kode yang sama juga menggerbangi step di
`Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` (`pyWorkPage.Type=="QR" || =="QP"` dan
`=="TR" || =="TP"`). **Bukan kode lokal PremiumList Life.**

### 2.2 Ambang tanggal 25 muncul di tiga konteks `[terverifikasi]`

| Konteks | Bentuk | Path |
| --- | --- | --- |
| EDM Treaty In | `@if(Local.TglProd=="",25,Local.TglProd)` | `Activity/SaveJsonPolisTreatyInEDM_Act.xml` |
| Endorsement Life | `@toDecimal(Local.currentdate) > 25` | `Activity/InsertJsonPolisLife_Act.xml` |
| PremiumList Life | `@toDecimal(Local.currentdate) > 25` | `Activity/InsertJsonPolisLife_Act.xml` |

`[dugaan]` angka 25 adalah **hari tutup buku bulanan**; di EDM Treaty In ia bahkan menjadi
*default* ketika query `GETTanggalClosing_SQL` tidak mengembalikan nilai. **Belum terverifikasi**,
dan **tidak boleh ditebak** → **OQ-030**.

---

## 3. Objek Oracle yang disentuh

Lima pemanggilan `RDB-List` dari `InsertJsonPolisLife_Act`:

| Class | RequestType |
| --- | --- |
| `ASM-FW-GISFW-Int-TREATYCONTRACT_LIFE` | `GetJsonProductLife` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_SUMMARY` | `InsertPLSummary` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_SUMMARY` | `InsertJsonPolis` |
| `ASM-FW-GISFW-Work-LIFE` | `SaveLifeinProduction_SQL` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_SUMMARY` | `GetNopolisByIDPega` |

`[terverifikasi]` Activity ini memakai **`Obj-Refresh-And-Lock`** di langkah 1 dan **`Commit`**
eksplisit di langkah 20 — jadi ada pengelolaan transaksi **di sisi Pega**, terpisah dari `COMMIT`
yang berada di dalam stored procedure (OQ-013).

Objek fisik per rule ada di `../inventory/PremiumList Life.md` §6. Tipe kolom tidak diketahui (OQ-001).

---

## 4. Integrasi eksternal

`[terverifikasi]` Tiga langkah terakhir yang menyentuh luar sistem, **seluruhnya digerbangi
`IsPEGAPROD`**:

| Step | Method | Precondition |
| ---: | --- | --- |
| 18 | `call @baseclass.SendEmailNotification` | `IsPEGAPROD` |
| 19 | `call serviceInsertArasapasLife_act` | `IsPEGAPROD` |
| 21 | `Connect-REST` | `IsPEGAPROD` |

Modul ini punya 1 rule `ConnectREST` (`../inventory/PremiumList Life.md` §7) yang memakai
`<pyBaseURLSelectionType>SETTING` → `LinkService!LinkService`, jadi alamatnya **konfigurasi**,
bukan literal.

**Ini bukti terkuat sejauh ini bahwa `IsPEGAPROD` adalah gerbang integrasi keluar** — bukan sekadar
guard lokal. Lihat §5.1.

---

## 5. Batas pengetahuan

### 5.1 `IsPEGAPROD` — kondisinya tidak diketahui, isinya berkonflik

`@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` (`PremiumList Life/When/IsPEGAPROD.xml`) terdaftar di
register **OQ-011 entri #305** — 14 modul, **2 isi berbeda**. Varian PremiumList Life berada di
grup yang sama dengan Endorsement Life, Claim Life, Komite Claim FacIn, Komite Claim Life,
NB FacIn, RNW Fac In.

**Kondisi yang diuji tidak terbaca** dari tag — `<pyLabel>` hanya berisi template kosong.
Karena rule ini menggerbangi **tiga** jalur integrasi keluar di modul ini, ketidaktahuan ini
material → **OQ-029**.

### 5.2 Empat ID literal yang menggerbangi logika

`[terverifikasi]` Step 7–10 memakai `@contains(.ID,"1000032")` … `"1000035"` sebagai precondition
untuk empat `Property-Set` berbeda. **Apa yang diwakili keempat ID itu tidak diketahui** — korpus
tidak memuat tabel referensi. → **OQ-031**.

### 5.3 Rule yang belum ditelusur

`serviceInsertArasapasLife_act`, `SendEmailNotification` (`@baseclass`), `Calculate1_Act`,
`GetPLNumber_Act`, `SavePremiumList_Act`, `SubmitPremiumList_Act`, `ValidasiUploadPL_act`,
`WPCLife_Act`, `ProtectAccept` — di luar jalur inti flow.

### 5.4 `IsFlagOnGoingPolicy` — tipe berbeda antar modul

`[terverifikasi]` Identitas dengan nama sama ada sebagai **dua tipe rule berbeda** di modul berbeda:

| Path | Tipe |
| --- | --- |
| `PremiumList Life/DecisionTable/IsFlagOnGoingPolicy.xml` | `DecisionTable` |
| `NB FacIn/When/IsFlagOnGoingPolicy.xml` | `When` |
| `RNW Fac In/When/IsFlagOnGoingPolicy.xml` | `When` |
| `Endorsment Fac In/When/IsFlagOnGoingPolicy.xml` | `When` |

Ini **varian baru** dari pola OQ-026: sebelumnya ditemukan dua tipe dalam **satu** modul
(`isApproved` di NB Treaty In); di sini tipenya berbeda **antar modul**. Isi `DecisionTable`
tidak terekstraksi oleh pembacaan tag sederhana. → memperluas **OQ-026**.

### 5.5 Nama activity sama, identitas berbeda

`[terverifikasi]` `InsertJsonPolisLife_Act` ada di dua modul dengan **class berbeda**, sehingga
**dua rule berbeda**, bukan konflik OQ-011:

| Path | Identitas | Step |
| --- | --- | ---: |
| `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` | `ASM-FW-GISFW-WORK-LIFE!INSERTJSONPOLISLIFE_ACT` | 21 |
| `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE!INSERTJSONPOLISLIFE_ACT` | 22 |

Hal yang sama berlaku untuk `IsLifeAccepted` (`...WORK-LIFE` vs `...WORK-ENDORSEMENTLIFE`).
Contoh konkret mengapa identitas rule **wajib** menyertakan class.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `PremiumList Life/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GISFW-WORK-LIFE` / `INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW` | `InputPolicyHolder.xml` (root modul) | tidak |
| `ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/InsertJsonPolisLife_Act.xml` | tidak |
| `ASM-FW-GISFW-WORK-LIFE` / `ISLIFEACCEPTED` / `RULE-DECLARE-DECISIONTABLE` | `DecisionTable/IsLifeAccepted.xml` | tidak |
| (DecisionTable) `ISFLAGONGOINGPOLICY` | `DecisionTable/IsFlagOnGoingPolicy.xml` | tidak |
| (FlowAction) `INPUTDATAOFFERLIFE` | `FlowAction/InputDataOfferLife.xml` | tidak |
| (FlowAction) `SHOWLIFEPREMIUMDETAIL` | `FlowAction/ShowLifePremiumDetail.xml` | tidak |
| (FlowAction) `SHOWLIFEPREMIUMSUMMARY` | `FlowAction/ShowLifePremiumSummary.xml` | tidak |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `When/IsPEGAPROD.xml` | **YA — entri #305** |
| 5 rule `RULE-CONNECT-SQL` (§3) | `RDBList/` | tidak |

**13 rule ditelusur.**

`ShowLifePremiumDetail` dan `ShowLifePremiumSummary` masing-masing ada sebagai **FlowAction dan
Section** — tabrakan nama lintas tipe, konsisten dengan temuan D1.

---

## 7. Pertanyaan terbuka baru

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-030** | Ambang **tanggal 25** muncul di 3 konteks (sebagai default dan sebagai perbandingan) — hari tutup buku? | Finance + Product+UW |
| **OQ-031** | Empat ID literal `1000032`–`1000035` menggerbangi `Property-Set` berbeda — mewakili apa? | Product+UW |

OQ yang **dikuatkan**: OQ-011 & OQ-029 (§5.1), OQ-018 (§4), OQ-020 (§2.1), OQ-023 (§1.1),
OQ-026 (§5.4 — varian baru: tipe berbeda antar modul), OQ-028 (§1.2).
