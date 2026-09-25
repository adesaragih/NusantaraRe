# Telusur Flow — Komite Claim FacIn

STEP D2, Tahap 2 konteks #2. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `Komite Claim FacIn/Flow/Komite_Flow.xml`
→ `RULE-OBJ-FLOW` / `ASM-FW-GCNMFW-WORK-KOMITE` / `KOMITE_FLOW`, `<pyStartActivity>Start1`.

---

## 1. Diagram alur

**1 Assignment, 1 Decision, 4 connector** `[terverifikasi]` — bentuk graf **sama persis** dengan
Komite Claim Life.

```
Start1 ─(Always)─> ASSIGNMENT63 "KomiteRouter"  [WorkList, pyRouteTo=Custom]
                        │ FlowAction: ViewTransferDtl
                        v
                    Decision1 "KomiteLoop"
                        ├─ When IsKomiteLoop ─> ASSIGNMENT63   ← LOOP
                        └─ Else ─────────────> END52
```

### 1.1 Bentuk sama, rule berbeda — dibuktikan, bukan diasumsikan

`[terverifikasi]` Keempat rule kunci **berbeda** dari Komite Claim Life, diuji dengan hash
ternormalisasi 18 tag (`_METHOD.md` §1.1):

| Rule | Komite Claim Life | Komite Claim FacIn | Hasil |
| --- | --- | --- | --- |
| Flow | `b69929c407` | `f9366db22d` | **BERBEDA** |
| `When/IsKomiteLoop.xml` | `ed34148ff5` | `597362839b` | **BERBEDA** |
| `Activity/KomiteRouter.xml` | `ee4d432924` | `83c9e74e7e` | **BERBEDA** |
| `FlowAction/ViewTransferDtl.xml` | `16bf7dbd52` | `cda2c47d22` | **BERBEDA** |

Class-nya juga berbeda: `ASM-FW-GCNMFW-WORK-KOMITE` vs `...-KOMITELIFE`. Jadi ini **rule yang
benar-benar berlainan**, bukan salinan.

**Namun kondisi `IsKomiteLoop` sama persis secara teks** `[terverifikasi]`:
`.AcceptStatus = "1"` DAN `.KomiteCount <= .KomiteLoop` — identik dengan Komite Claim Life.
Hash berbeda karena class dan metadata, bukan karena logikanya.

### 1.2 Routing: tangga **ter-hardcode empat tingkat**

`[terverifikasi]` `Komite Claim FacIn/Activity/KomiteRouter.xml`
(`ASM-FW-GCNMFW-WORK-KOMITE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`) — 7 step:

| Precondition | `param.AssignTo` |
| --- | --- |
| `.KomiteCount == 1` | `"komitepnc"` |
| `.KomiteCount == 2` | `"komitepnc2"` |
| `.KomiteCount == 3` | `"komitepnc3"` |
| `.KomiteCount == 4` | `"komitepnc4"` |
| `Primary.TransferType == '2'` | (cabang terpisah) |
| `Primary.TransferType != '2'` | (cabang terpisah) |

**Perbedaan mendasar dari Komite Claim Life:** di Life sasaran penugasan adalah **nilai data**
(`param.AssignTo = .KomiteID`); di FacIn sasaran adalah **empat string ter-hardcode**.

`[dugaan]` `komitepnc`…`komitepnc4` adalah workbasket/operator per tingkat komite; `PNC`
**kepanjangan belum terverifikasi**. Ini nilai konfigurasi yang **tidak boleh dipindahkan apa
adanya** ke sistem baru → **OQ-036**.

`[terverifikasi]` Berbeda dari Life, FacIn punya cabang eksplisit untuk `TransferType != '2'`.

---

## 2. Roster komite dibangun dari database

`[terverifikasi]` `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml`
(`ASM-FW-GCNMFW-WORK-KOMITE` / `APPROVALKOMITE_ACT` / `RULE-OBJ-ACTIVITY`, **146.689 byte**).

Langkah yang terbaca:

| Step | Method | Precondition |
| ---: | --- | --- |
| 3 | `Property-Set` | `.IsFacRetro == 1` |
| 4 | `Exit-Activity` | `Local.Retro == 1` |
| 5 | `Obj-Browse` | `pyWorkPage.KomiteCount == 1` |
| 6 | `Page-Remove` | — |
| 7 | `Obj-Browse` | **`Local.TotalAdj > 30000000.00 && Local.TotalAdj <= 57750000.00`** |
| 9 | `Property-Remove` | `pyWorkPage.KomiteList(1).KomiteID == .OPERATOR_ID` |
| 10 | `Property-Set` | `KomiteCount==1 && @LengthOfPageList(pyWorkPage.KomiteList) = 1` |

`[terverifikasi]` Sumber roster: `Obj-Browse` terhadap `<ObjClass>ASM-FW-GCNMFW-Int-EMAILKOMITE`
ke halaman `GetKomite`. Jadi **daftar anggota komite berasal dari tabel database**, bukan
ter-hardcode — berbeda dari sasaran routing di §1.2 yang justru ter-hardcode.

### 2.1 Ambang nominal ter-hardcode `[terverifikasi]`

```
Local.TotalAdj > 30000000.00 && Local.TotalAdj <= 57750000.00
```

Sebuah **pita nilai** (30.000.000,00 sampai 57.750.000,00) tertanam langsung di rule, dipakai
memilih `Obj-Browse` mana yang dijalankan — artinya **komposisi roster komite bergantung pada
besaran nilai**.

`Local.TotalAdj` diakumulasi dari `Local.TotalAdj + .ValueAdjustment`.

`[dugaan]` ini **batas wewenang persetujuan** berdasarkan nominal. **Belum terverifikasi** — mata
uang tidak disebut, dan hanya satu pita yang terbaca dari pembacaan tag ini. → **OQ-037**.

**Tidak disimpulkan** berapa tingkat wewenang seluruhnya atau siapa yang berwenang di tiap pita.

---

## 3. Status / state yang berubah

Dari `Komite Claim FacIn/Activity/KomitePost_Adjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITE` / `KOMITEPOST_ADJUSTMENT` / `RULE-OBJ-ACTIVITY`, **626.150 byte** —
rule terbesar yang ditelusur sejauh ini di D2):

| Properti / kondisi | Nilai literal | Peran |
| --- | --- | --- |
| `pyWorkPage.TransferType` | `"2"` | gerbang `Obj-Open-By-Handle` |
| `pyWorkPage.AcceptStatus` | `"1"`, `"2"` | dua `Property-Set` berbeda |
| `pyWorkPage.KomiteCount` | `!= local.operator` | gerbang step 1 |
| terminal tangga | `KomiteCount == Local.TotalKomite && AcceptStatus == "1"` | — |
| `OutputData.START_DATE` | `== ""` | gerbang **enam** RDB-List |
| `.AcceptanceStatus` | `""`, `"0"` | gerbang RDB-List + Property-Set |
| `CekPLATreaty.START_DATE` | `== ""` | gerbang RDB-List |
| `IsFire` | When | memilih `GenerateNoAccept` vs `GenerateNoAcceptNonFire` |
| `@LengthOfPageList(.ComiteeClaim) < @LengthOfPageList(pyWorkPage.KomiteList)` | — | gerbang Property-Set |

`[terverifikasi]` **Terminal tangga memakai `Local.TotalKomite`**, bukan `.KomiteLoop` seperti di
Komite Claim Life — dua modul memakai properti berbeda untuk konsep yang sama.

### 3.1 Kode PaymentType / TransferType (OQ-020)

`[terverifikasi]` Nilai literal yang diuji di modul ini:

| Kode | Nilai yang muncul |
| --- | --- |
| `PaymentType` | `1`, `2`, `3`, `4`, `5`, `6` (terbanyak: `"4"` 32×, `"6"` 22×) |
| `TransferType` | `2` |

```
grep -rhoE "(PaymentType|TransferType)[ ]*[=!]+[ ]*[\"']?[0-9]{1,2}" "Komite Claim FacIn" --include="*.xml" | sort | uniq -c
```

**Arti seluruh nilai belum terverifikasi.** Dicatat literal apa adanya.

### 3.2 Guard identitas orang `[terverifikasi]`

**Empat activity** memakai identitas operator sebagai literal:

| File | Peran dalam tangga |
| --- | --- |
| `Activity/KomitePost_Adjustment.xml` | pasca-keputusan penyesuaian |
| `Activity/KomitePost_CloseClaim.xml` | pasca-keputusan tutup klaim |
| `Activity/KomitePost_Reject.xml` | pasca-keputusan tolak |
| `Activity/SetProteksiSubmiteKomite.xml` | proteksi submit |

Guard berada di `<pyExpression>` (7×), `<pyMemo>` (1×), `<pyPropertiesValue>` (1×).
**Nilai nama orang tidak disalin** ke artefak ini (OQ-021).

Guard jabatan: `pyPosition != "IT Developer"` (2×), `== "IT Developer"` (2×), `== "SPV B"` (1×).

**Ini modul dengan guard identitas terpadat di domain Komite** — dan guard-nya berada tepat di
tiga activity keputusan komite. → OQ-021.

---

## 4. Objek Oracle yang disentuh

| Class | RequestType | Precondition |
| --- | --- | --- |
| `ASM-FW-GISFW-Int-policyjson` | `GETTanggalClosing_SQL` | — |
| `ASM-FW-GISFW-Int-policyjson` | `GetKodeProdNonLife_SQL` | — |
| `ASM-FW-GISFW-Int-policyjson` | `GetSequenceNumber_SQL` | — |
| `ASM-FW-GCNMFW-Int-V_POLIS` | `GenerateNoAccept` | `IsFire` |
| `ASM-FW-GCNMFW-Int-V_POLIS` | `GenerateNoAcceptNonFire` | `OutputData.START_DATE==""` (6×) |
| `ASM-FW-GCNMFW-Int-V_POLIS` | `SaveOSClaim_SQL` | `.AcceptanceStatus=="" \|\| =="0"` |
| `ASM-FW-GISFW-int-policyjson` | `InsertHistoryAkseptasiPega_Sql` | — |
| `ASM-FW-GCNMFW-Work-Komite` | `UpdateSubProgresKlaim` | — |

`[terverifikasi]` Penomoran akseptasi **bercabang menurut lini bisnis**: `GenerateNoAccept` bila
`IsFire`, `GenerateNoAcceptNonFire` selain itu. Pola setara dengan Komite Claim Life yang bercabang
menurut kode `Type` (`QP/QR` vs `TP/TR`) — **dua modul memakai kriteria pencabangan berbeda untuk
keputusan yang sama.**

`[terverifikasi]` `ASM-FW-GISFW-Int-policyjson` dan `ASM-FW-GISFW-int-policyjson` muncul dengan
**beda kapitalisasi** di rule yang sama — `[pertanyaan terbuka]` apakah Pega memperlakukannya sebagai
class yang sama.

---

## 5. Batas pengetahuan

### 5.1 Rule berkonflik yang tersentuh

`[terverifikasi]` Dari register OQ-011, identitas berikut menyentuh modul ini dan berada di jalur
atau lingkungan telusur:

| Identitas | Entri |
| --- | --- |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | #305 — 14 modul, 2 isi |
| `ASM-FW-GCNMFW-WORK` / `KONVERSIKLAIM_ACT` / `RULE-OBJ-ACTIVITY` | berkonflik |
| `ASM-FW-GCNMFW-WORK` / `GETBASE64ATTACHMENT` / `RULE-OBJ-ACTIVITY` | berkonflik |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | berkonflik |
| keluarga `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` (7 rule) | berkonflik |

Rule-rule ini **tidak ditelusur isinya** di sini karena berada di luar jalur inti. Bila telusur
lanjutan menyentuhnya, **wajib per-varian** (`_METHOD.md` §1.1).

### 5.2 Isi `ApprovalKomite_Act` belum habis dibaca

Activity ini 146 KB; pembacaan tag menampilkan 11 langkah pertama. **Jumlah pita nominal
seluruhnya, dan aturan penyusunan roster lengkapnya, belum terbaca** → bagian dari **OQ-037**.

### 5.3 `IsPEGAPROD`

Tidak muncul di jalur inti yang ditelusur di modul ini, tetapi rule-nya ada
(`Komite Claim FacIn/When/IsPEGAPROD.xml`, varian grup yang sama dengan Komite Claim Life).
Kondisinya tetap **tidak terbaca** → OQ-029.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Komite Claim FacIn/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-KOMITE` / `KOMITE_FLOW` / `RULE-OBJ-FLOW` | `Flow/Komite_Flow.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY` | `Activity/KomiteRouter.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITE` / `APPROVALKOMITE_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/ApprovalKomite_Act.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITE` / `KOMITEPOST_ADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `Activity/KomitePost_Adjustment.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITE` / `ISKOMITELOOP` / `RULE-OBJ-WHEN` | `When/IsKomiteLoop.xml` | tidak |
| (FlowAction) `VIEWTRANSFERDTL` | `FlowAction/ViewTransferDtl.xml` | tidak |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `When/IsPEGAPROD.xml` | **YA — #305** |
| 8 rule `RULE-CONNECT-SQL` (§4) | `RDBList/` | sebagian |

**15 rule ditelusur.** Disebut tetapi belum ditelusur: `KomitePost_CloseClaim`, `KomitePost_Reject`,
`KomitePostAct`, `SaveAccept_ACT`, `SaveReject_ACT_KMT`, `SetProteksiSubmiteKomite`,
`HitServiceToKasirKMT_Act`, `SendEmailKlaim_KMT`, `KonversiKlaim_Act`.

---

## 7. Pertanyaan terbuka baru

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-036** | Empat sasaran routing ter-hardcode (`komitepnc`…`komitepnc4`); `PNC` kepanjangan belum terverifikasi | IAM + Product+UW |
| **OQ-037** | Ambang nominal ter-hardcode (`> 30.000.000,00 && <= 57.750.000,00`) menentukan komposisi roster komite; mata uang tidak disebut; pita lain belum terbaca | Finance + Product+UW |

OQ dikuatkan: OQ-011 & OQ-029 (§5), OQ-020 (§3.1), OQ-021 (§3.2), OQ-024 (§1.2 — sasaran routing
kini terbaca sebagian), OQ-028 (`WorkList`), OQ-032 (`Local.TotalKomite` vs `.KomiteLoop`).
