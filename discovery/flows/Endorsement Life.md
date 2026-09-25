# Telusur Flow — Endorsement Life

STEP D2, Tahap 1 konteks #3. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `Endorsement Life/Flow/InputEDMLife.xml`
→ `RULE-OBJ-FLOW` / `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE`,
`<pyStartActivity>Start1`. Ukuran 68.374 byte. Modul ber-`Flow` terkecil di korpus (75 file).

---

## 1. Diagram alur

**2 Assignment, 1 Decision, 1 Utility, 4 connector** `[terverifikasi]` — flow paling sederhana
yang ditelusur sejauh ini.

```
Start1 ─(Always)─> ASSIGNMENT63 "InputEDMLife"  [WorkList]
                        │ FlowAction: InputEDMLife
                        v
                    Decision1 "Accept"  (DecisionTable: IsLifeAccepted)
                        ├─ Status=Decline ─> End1
                        └─ Status=Confirm ─> Utility1  Activity: InsertJsonPolisLife_Act
                                                 └─(…)─> END52
```

### 1.1 Shape tanpa connector masuk `[terverifikasi]`

```
Assignment1 "Input EDM Summary" [WorkList] ──FlowAction InputEDMLife_Summary──> END52
```

Tidak ada satu pun pasangan `<pyFrom>`/`<pyTo>` yang menuju `Assignment1`. Pasangan yang ada:
`Start1→ASSIGNMENT63`, `ASSIGNMENT63→Decision1`, `Decision1→End1`, `Decision1→Utility1`,
`Utility1→END52`, `Assignment1→END52`.

Ini **pola kedua** setelah NB Treaty In (OQ-023): shape yang hanya punya connector keluar.
Ditambahkan ke **OQ-023**.

### 1.2 Routing memakai `WorkList`, bukan `WorkBasket` `[terverifikasi]`

Kedua Assignment memakai `<pyImplementation>WorkList`. Kontras dengan NB Treaty In dan
EDM Treaty In yang seluruh Assignment-nya `WorkBasket`.

`[dugaan]` `WorkList` = antrean per-pengguna, `WorkBasket` = antrean bersama/peran — perbedaan ini
berdampak pada model penugasan. **Belum terverifikasi** → **OQ-028**.

---

## 2. Status / state yang berubah

Seluruhnya dari `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 22 step)
kecuali disebut lain.

| Properti / kondisi | Nilai literal | Peran | Arti |
| --- | --- | --- | --- |
| `Local.currentdate` | diuji `@toDecimal(...) > 25` | precondition step 2 | **ambang tanggal 25** — sama dengan default `TglProd=25` di EDM Treaty In |
| `pyWorkPage.PremiumListSummary.PL_NUMBER_EDM` | `== ""` | precondition step 5 & 6 | nomor premium list kosong |
| `pyWorkPage.Type` | **`"QR"`, `"QP"`** | precondition step 12 | **arti belum terverifikasi** (OQ-020) |
| `pyWorkPage.Type` | **`"TR"`, `"TP"`** | precondition step 13 | **arti belum terverifikasi** (OQ-020) |
| `.EDMStatus` | `"Old"` / `!= "Old"` | precondition step 14 & 15 | **arti belum terverifikasi** |
| `TempError.CARIDESC` | `== 1` | precondition step 11 | penanda error |
| `@SizeOfPropertyList(TempWorkPage.ListLifePremiumDetailUpload)` | `> 50000` | precondition step 16 | **ambang 50.000 baris** |
| `@SizeOfPropertyList(pyWorkPage.Addendum.OldData.CurrencyList)` | `<= 0` | precondition step 3 | daftar mata uang kosong |
| `OutDataLife.pxResults(1).PL_NUMBER` | `== ""` | precondition `SendEmailNotification` | memicu notifikasi email |
| `IsPEGAPROD` | When guard | precondition `serviceInsertArasapasLife_act` | **guard lingkungan** — lihat §5.1 |

`[terverifikasi]` Kode `QR`/`QP`/`TR`/`TP` **tidak eksklusif milik PremiumList Life** — dipakai juga
di sini sebagai precondition step. Memperluas **OQ-020**.

### 2.1 Bukti untuk `EdmType` (OQ-020)

`[terverifikasi]` `EdmType` dipakai di **13 rule** di modul ini:

```
grep -rlE "EdmType" "Endorsement Life" --include="*.xml"
```

| Tipe rule | File |
| --- | --- |
| Activity | `CreateCaseEMDL.xml`, `MappingEDMLife.xml`, `SaveCSVEDMLife.xml`, `SetErrorBatalEndorsement_Act.xml`, `SetPremi_EDM.xml` |
| Harness | `EndorsmentLife_harnes.xml`, `InboxEndorsementLife.xml` |
| Section | `EndorsmentLife_Section.xml`, `InboxEndorsementLife.xml`, `InputEDMLife.xml`, `ShowLifePremiumSummary_EDM.xml` |
| RDBList | `GetEdmTypeLife.xml` |
| ReportDefinition | `InboxEDMLife.xml` |

**Nilai literal yang diuji: hanya `1` dan `3`** `[terverifikasi]` (`EdmType=1` 8×, `EdmType==3` 6×,
`EdmType=3` 5×, plus varian berkutip `'1'`/`'3'`). **Tidak ditemukan nilai 2** di modul ini.

```
grep -rhoE "EdmType[^<\"]{0,3}[=!]+[^<\"]{0,6}" "Endorsement Life" --include="*.xml" | sort | uniq -c
```

**Arti nilai 1 dan 3 belum terverifikasi.** Ketiadaan nilai 2 juga tidak dapat ditafsirkan —
mungkin dipakai modul lain, mungkin memang tidak ada.

**Sumber nilainya** `[terverifikasi]` — `Endorsement Life/RDBList/GetEdmTypeLife.xml`:

```sql
SELECT A.DATA_JSON.EdmType AS CARI1 FROM POOLDATA.JSON_POLIS A
WHERE NOPOLIS = {TempWork.PolicyNo} AND PRODKE IS NOT NULL ORDER BY PRODKE DESC
```

`EdmType` **disimpan di dalam kolom JSON** `DATA_JSON` pada tabel `POOLDATA.JSON_POLIS`, diambil
dari baris `PRODKE` terbesar. Lihat §5.3.

---

## 3. Objek Oracle yang disentuh

Tujuh pemanggilan `RDB-List` dari `InsertJsonPolisLife_Act`:

| Class | RequestType | Precondition |
| --- | --- | --- |
| `ASM-FW-GISFW-Int-OFFERJSON` | `GetProdKeOldData_SQL` | `CurrencyList` kosong |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` | `Generate_NoEndorsmentLife` | `PL_NUMBER_EDM == ""` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_SUMMARY` | `InsertJsonPolisEDM` | — |
| `ASM-FW-GISFW-Work-EndorsementLife` | `SaveLifeinProduction_SQL` | — |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` | `SaveMasterLPDet` | > 50.000 baris |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_SUMMARY` | `InsertPLSummary` | — |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_SUMMARY` | `GetNopolisByIDPega` | — |

Ditambah `Endorsement Life/RDBList/GetEdmTypeLife.xml` → `POOLDATA.JSON_POLIS`.

Objek fisik yang disentuh tiap rule ada di `../inventory/Endorsement Life.md` §6. Tipe kolom tetap
tidak diketahui (OQ-001).

---

## 4. Integrasi eksternal

`[terverifikasi]` Langkah terakhir `InsertJsonPolisLife_Act`:

```
call serviceInsertArasapasLife_act   pre: IsPEGAPROD
```

Panggilan service eksternal **digerbangi oleh When `IsPEGAPROD`** — lihat §5.1. Modul ini punya
1 rule `ConnectREST` (`../inventory/Endorsement Life.md` §7); apakah `serviceInsertArasapasLife_act`
memanggilnya **belum ditelusur** (di luar jalur inti).

Notifikasi email: `call @baseclass.SendEmailNotification` bila
`OutDataLife.pxResults(1).PL_NUMBER == ""`.

---

## 5. Batas pengetahuan

### 5.1 `IsPEGAPROD` — guard lingkungan, dan ia **berkonflik** (OQ-011)

`[terverifikasi]` `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` ada di **14 modul** dan
**terdaftar di register OQ-011 entri #305 — 14 varian, 2 isi berbeda**.

Pengelompokan hash ternormalisasi `[terverifikasi]`:

| Grup | Modul |
| --- | --- |
| A | Claim Fac In, Claim Non Prop, Claim Prop, EDM Treaty In, Endorsment Fac In, Komite Claim Non Prop, Komite Claim Prop |
| B | Claim Life, **Endorsement Life**, Komite Claim FacIn, Komite Claim Life, NB FacIn, PremiumList Life, RNW Fac In |

Telusur ini membaca **varian grup B** (`Endorsement Life/When/IsPEGAPROD.xml`). Varian grup A
**berbeda isi** dan tidak dibaca.

**Kondisi rule ini tidak terbaca dari tag.** `<pyLabel>` hanya berisi template kosong
`[first value][relation][second value]` — kondisi sebenarnya tersimpan dalam bentuk yang tidak
terekstraksi oleh pembacaan tag sederhana. **Apa yang diuji `IsPEGAPROD` tidak diketahui.**

**Mengapa ini penting:** namanya menyiratkan pemeriksaan apakah sistem berjalan di lingkungan
production, dan ia **menggerbangi pemanggilan service eksternal**. Bila benar, maka perilaku
integrasi **berbeda antar lingkungan**, dan dua isi berbeda di 14 modul berarti gerbang itu **tidak
konsisten**. Ini menyentuh langsung **OQ-018**. → **OQ-029**.

**Tidak boleh disimpulkan dari namanya** (`_METHOD.md` §2.2).

### 5.2 Fungsi dan rule yang belum ditelusur

`serviceInsertArasapasLife_act`, `SendEmailNotification` (`@baseclass`),
`CreateCaseEMDL`, `MappingEDMLife`, `SaveCSVEDMLife`, `SetPremi_EDM` — di luar jalur inti.

### 5.3 Struktur JSON **sebagian dapat dibaca** — kemajuan untuk OQ-012

`[terverifikasi]` SQL di korpus menanyakan **ke dalam** kolom JSON memakai notasi path Oracle,
mis. `A.DATA_JSON.EdmType`. Sapuan korpus menemukan **49 path JSON unik**:

```
grep -rhoE "(DATA_JSON|JSONDATA)\.[A-Za-z_]+" . --include="*.xml" | sort -u | wc -l
```

Yang paling sering: `DATA_JSON.osAkseptasi` (78×), `DATA_JSON.CurrencyList` (33×),
`DATA_JSON.Value` (23×), `DATA_JSON.AcceptedNo` (23×), `JSONDATA.TreatyYear` (22×),
`JSONDATA.TreatyGroupID` (22×), `DATA_JSON.CoverageName` (21×), `JSONDATA.TreatyDescID` (18×).

**Artinya:** struktur JSON **tidak sepenuhnya opaque**. Nama field yang benar-benar dipakai dapat
direkonstruksi dari SQL. Yang tetap **tidak diketahui**: tipe, kardinalitas, field yang ditulis
tetapi tidak pernah dibaca, dan field opsional. → memperbarui **OQ-012**.

Dua nama kolom JSON berbeda muncul: `DATA_JSON` dan `JSONDATA` — `[pertanyaan terbuka]` apakah dua
kolom berbeda di tabel berbeda, atau penamaan tidak konsisten.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Endorsement Life/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-FLOW` | `Flow/InputEDMLife.xml` | tidak |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/InsertJsonPolisLife_Act.xml` | tidak |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `ISLIFEACCEPTED` / `RULE-DECLARE-DECISIONTABLE` | `DecisionTable/IsLifeAccepted.xml` | tidak |
| (FlowAction) `INPUTEDMLIFE` | `FlowAction/InputEDMLife.xml` | tidak |
| (FlowAction) `INPUTEDMLIFE_SUMMARY` | `FlowAction/InputEDMLife_Summary.xml` | tidak |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `When/IsPEGAPROD.xml` | **YA — entri #305** |
| (RDBList) `GETEDMTYPELIFE` | `RDBList/GetEdmTypeLife.xml` | tidak |
| 7 rule `RULE-CONNECT-SQL` (§3) | `RDBList/` | tidak |

**14 rule ditelusur.**

`[terverifikasi]` **`InputEDMLife` adalah nama yang dipakai tiga tipe rule sekaligus** di modul ini:
`Flow/InputEDMLife.xml`, `FlowAction/InputEDMLife.xml`, `Section/InputEDMLife.xml`. Contoh konkret
aturan `../README.md` §3.2 — rule wajib ditulis `class / nama / tipe`.

---

## 7. Pertanyaan terbuka baru

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-028** | `WorkList` vs `WorkBasket` — dua model penugasan berbeda antar modul; apa bedanya secara proses? | Product+UW + IAM |
| **OQ-029** | `IsPEGAPROD` menggerbangi panggilan service eksternal, kondisinya tidak terbaca, dan ia punya **2 isi berbeda** di 14 modul | IAM + DBA + Product+UW |

OQ yang **dikuatkan**: OQ-011 (§5.1), OQ-012 (§5.3 — kini sebagian terbaca), OQ-018 (§5.1),
OQ-020 (§2, §2.1), OQ-023 (§1.1), OQ-026 (`IsLifeAccepted` juga `DecisionTable`).
