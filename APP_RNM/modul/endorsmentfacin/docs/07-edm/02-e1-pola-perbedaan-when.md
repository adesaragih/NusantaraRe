# Discovery Endorsement — E-1: pola perbedaan NB↔EDM

> **Sumber:** `Endorsment Fac In\` vs `NB FacIn\`. Perbandingan nama **peka huruf**
> (`HashSet[string]`), pengurutan **`[StringComparer]::Ordinal`**. Label mengikuti `CLAUDE.md` §3.
>
> 🔁 **Dokumen ini adalah REVISI.** Versi pertama melaporkan **1.202 identik / 505 berbeda
> (70,4 %)** dengan metode yang **tidak reproducible**. Angka itu **dicabut**. Riwayat pengukuran
> dan sebab tiap koreksi dicatat di §0 — bukan disembunyikan.

---

## 0. Riwayat angka — mengapa versi ini boleh dipercaya

⚠️ **Angka ini sudah bergerak empat kali.** Mencatatnya adalah bagian dari bukti: yang membuat angka
terakhir dapat dipercaya bukan besarannya, melainkan bahwa **daftar tag yang dibuang kini tertutup
dan teruji**, bukan ditentukan sambil jalan.

| # | Angka identik | Metode | Sebab gugur |
| ---: | ---: | --- | --- |
| 1 | **0** (0 %) | hash mentah byte-per-byte | Stempel waktu ekspor membuat setiap berkas berbeda |
| 2 | **1.202** (70,4 %) | buang `px*`/`pz*` **menyeluruh** + 2 tag `py` berstempel waktu | ⛔ **Membuang `pxRuleClassName` / `pxRuleFamilyName` — kelas rule ikut terbuang.** Beda kelas adalah beda fungsional, bukan artefak |
| 3 | **1.113** (65,2 %) | daftar volatile eksplisit (20 tag), kelas dipertahankan | Daftar **belum lengkap** — `pxHostId`, `pzIndexCount`, `pxWarningCreatedTime` masih ikut terhitung sebagai perbedaan |
| 4 | **1.195** (70,0 %) | daftar volatile **23 tag, tertutup** — lihat §1 | **berlaku** |

📌 Angka final **70,0 %** dekat dengan klaim lama **70,4 %**, tetapi **bukan pembenaran angka lama**.
Metode lama menghasilkannya dengan membuang kelas rule; metode ini mempertahankan kelas dan justru
menemukan **7 berkas tambahan yang berbeda** (512, bukan 505). Kedekatan angkanya kebetulan.

⚠️ **Yang sesungguhnya diperbaiki bukan angkanya, melainkan sifatnya:** angka lama tidak punya daftar
tag eksplisit, sehingga tidak dapat direproduksi work owner. Angka ini punya.

---

## 1. Kontrak metode — daftar tag yang dibuang

Ini **kontrak**. Tanpa daftar eksplisit, angka berapa pun tidak reproducible.

### 1.1 Dibuang seluruh barisnya (23 tag) — alasan per tag

| Tag | Alasan volatile |
| --- | --- |
| `pxCreateDateTime` · `pxUpdateDateTime` · `pxSaveDateTime` · `pxCommitDateTime` | Stempel waktu operasi penyimpanan |
| `pxMoveImportDateTime` · `pxOriginalCreateDateTime` | Stempel waktu impor/asal — berbeda per ekspor |
| `pyRuleFormStatusTime` | Stempel waktu status form. ⚠️ Berawalan `py` — **lolos** saringan `px*`/`pz*` |
| `pxWarningCreatedTime` | Stempel waktu peringatan rule |
| `pxCreateOperator` · `pxCreateOpName` · `pxUpdateOperator` · `pxUpdateOpName` | Identitas operator ekspor — **juga alasan privasi** (`CLAUDE.md` §3.5) |
| `pxMoveImportOperId` · `pxMoveImportOperName` · `pxOriginalCreateOperator` · `pxOriginalCreateOpName` | idem |
| `pxCreateSystemID` · `pxUpdateSystemID` · `pxOriginalCreateSystemID` | ID sistem Pega asal — korpus dirakit dari lebih dari satu server (`CLAUDE.md` §3.1) |
| `pxHostId` | Hostname mesin ekspor — memuat nama host DEV |
| `pzChecksum` | Checksum turunan isi, bukan isi |
| `pzIndexCount` | Cacah indeks turunan, dihitung ulang saat impor |
| `pyShowJavaWindowName` | Nama window Java ber-stempel waktu (`GeneratedJava<TS>`). ⚠️ Berawalan `py` |

### 1.2 Dipertahankan tagnya, **stempel waktu di dalamnya** diganti `<TS>` (6 tag)

`pzInsKey` · `pzIndexOwnerKey` · `pzOriginalInstanceKey` · `pzDocumentKey` · `pyJavaClassName` ·
`pxInsName`

⛔ **Tag ini memuat NAMA KELAS di dalam nilainya.** Membuangnya akan menyembunyikan perbedaan kelas.
Yang dibuang hanya substring stempel waktunya.

### 1.3 ⛔ TIDAK PERNAH dibuang — meskipun berawalan `px`/`pz`

`pxRuleClassName` · `pxRuleFamilyName` · `pxRuleObjClass` · `pxObjClass` · `pyDesignatedClass` ·
`pyClassName` · `pyRuleName` · `pyStepPageReference` · `pyActivityClass` · `pyStepsParentClass` ·
`rowdata` · `pyValue` · `pyLabel` · `pyName` · seluruh tag isi lainnya.

**Beda kelas = beda fungsional.** Inilah kesalahan yang menggugurkan angka #2.

### 1.4 Bukti daftar sudah tertutup

`[terverifikasi]` Setelah 23 tag itu dibuang, **tag apa yang masih berbeda** pada 512 pasangan:

| Tag sisa | Kemunculan | Sifat |
| --- | ---: | --- |
| `rowdata` | 26.226 | isi baris — **fungsional** |
| `pzIndexOwnerKey` | 4.525 | memuat kelas — **fungsional** |
| `pyStepPageReference` | 3.119 | rujukan halaman langkah — **fungsional** |
| `pxObjClass` · `pxRuleFamilyName` · `pxRuleClassName` · `pxRuleObjClass` | 2.270 · 2.047 · 1.610 · 372 | **kelas rule** |
| `pyRuleName` | 2.022 | nama rule yang dirujuk — **fungsional** |
| `pyCellId` · `pyAutomationID` · `pyWidth` · `pyHeight` · `pyUIElement` | 1.451 · 695 · 448 · 392 · 442 | tata letak UI |
| `pyValue` · `pyLabel` · `pyName` · `pyParametersParamName` | 718 · 700 · 559 · 638 | isi |

📌 **Tidak tersisa satu pun** tag stempel waktu, operator, host, atau checksum. Daftar tertutup.

### 1.5 Uji sensitivitas — apakah tata letak UI memengaruhi angka?

`[terverifikasi]` `pyCellId`/`pyAutomationID`/`pyWidth`/`pyHeight` adalah **kosmetik UI**, batas
penilaian satu-satunya yang tersisa. Diukur dua arah, bukan diputuskan diam-diam:

| Perlakuan | Berbeda | Identik |
| --- | ---: | ---: |
| Kosmetik UI **dihitung** | **512** (30,0 %) | 1.195 (70,0 %) |
| Kosmetik UI **diabaikan** | **512** (30,0 %) | 1.195 (70,0 %) |
| **Berkas yang bedanya HANYA kosmetik** | **0** | |

✅ **Angka tidak bergantung pada keputusan itu.** Nol berkas berbeda hanya karena tata letak.

### 1.6 Uji silang wajib

| Berkas | Harus | Hasil |
| --- | --- | --- |
| `When\FlagOldData` — **beda kelas** (`…Data-Accessory` ↔ `…Data-Coverage`) | **BEDA** | ✅ **BEDA** |
| `When\IsAdmin` · `IsBonding` · `IsAddButton` — beda hanya 2 tag stempel waktu | identik | ✅ identik |
| `DataPage\D_AnekaList` — beda 12 baris, seluruhnya volatile | identik | ✅ identik |

⚠️ `D_AnekaList` adalah uji yang **menjatuhkan angka #3**: di sana ia masih terklasifikasi "beda"
karena `pxHostId`/`pzIndexCount` belum masuk daftar.

### 1.7 Perintah audit — dapat dijalankan ulang

```powershell
$nbRoot="D:\migrasi\RNM\NB FacIn"; $edRoot="D:\migrasi\RNM\Endorsment Fac In"
$VOLATILE=@(
 'pxCreateDateTime','pxUpdateDateTime','pxSaveDateTime','pxCommitDateTime','pxMoveImportDateTime',
 'pxOriginalCreateDateTime','pyRuleFormStatusTime','pxWarningCreatedTime','pxCreateOperator',
 'pxCreateOpName','pxCreateSystemID','pxUpdateOperator','pxUpdateOpName','pxUpdateSystemID',
 'pxMoveImportOperId','pxMoveImportOperName','pxOriginalCreateOperator','pxOriginalCreateOpName',
 'pxOriginalCreateSystemID','pxHostId','pzChecksum','pzIndexCount','pyShowJavaWindowName')
$KEYLIKE=@('pzInsKey','pzIndexOwnerKey','pzOriginalInstanceKey','pzDocumentKey','pyJavaClassName','pxInsName')
$rxVol=[regex]('^</?(' + ($VOLATILE -join '|') + ')[ />]|^<(' + ($VOLATILE -join '|') + ')>')
$rxKey=[regex]('^<(' + ($KEYLIKE -join '|') + ')>')
$rxTs=[regex]'\d{8}T\d{6}(\.\d+)?( GMT)?|_\d{8}T\d{6}_\d+'
$sha=[Security.Cryptography.SHA256]::Create()
function Get-Sig([string]$path){
  $keep=New-Object Collections.Generic.List[string]
  foreach($ln in [IO.File]::ReadAllLines($path)){ $tr=$ln.Trim(); if(-not $tr){continue}
    if($rxVol.IsMatch($tr)){continue}
    if($rxKey.IsMatch($tr)){ $tr=$rxTs.Replace($tr,'<TS>') }   # buang stempel waktu, PERTAHANKAN kelas
    $keep.Add($tr) }
  $a=$keep.ToArray(); [Array]::Sort($a,[StringComparer]::Ordinal)
  [BitConverter]::ToString($sha.ComputeHash([Text.Encoding]::UTF8.GetBytes([string]::Join("`n",$a)))).Replace('-','') }
$sama=0; $beda=0
foreach($tipe in (Get-ChildItem $edRoot -Directory).Name){
  $nbSet=New-Object 'System.Collections.Generic.HashSet[string]'   # peka huruf — JANGAN Test-Path
  Get-ChildItem "$nbRoot\$tipe" -Filter *.xml -File -EA SilentlyContinue|%{[void]$nbSet.Add($_.BaseName)}
  foreach($f in (Get-ChildItem "$edRoot\$tipe" -Filter *.xml -File -EA SilentlyContinue)){
    if($nbSet.Contains($f.BaseName)){
      if((Get-Sig "$nbRoot\$tipe\$($f.BaseName).xml") -eq (Get-Sig $f.FullName)){$sama++}else{$beda++} } } }
"identik=$sama beda=$beda total=$($sama+$beda)"   # → identik=1195 beda=512 total=1707
```

---

## 2. Angka terverifikasi

`[terverifikasi]` **1.707** berkas bernama sama (peka huruf):

| Ukuran | Jumlah | % |
| --- | ---: | ---: |
| **Isi identik** (setelah normalisasi §1) | **1.195** | **70,0 %** |
| **Benar-benar berbeda** | **512** | **30,0 %** |

### Per tipe rule

| Tipe | Identik | **Berbeda** | Bersama | % berbeda |
| --- | ---: | ---: | ---: | ---: |
| Activity | 294 | **194** | 488 | 39,8 % |
| Section | 252 | **101** | 353 | 28,6 % |
| When | 112 | **63** | 175 | 36,0 % |
| FlowAction | 166 | 45 | 211 | 21,3 % |
| DataTransform | 93 | 35 | 128 | 27,3 % |
| ReportDefinition | 83 | 27 | 110 | 24,5 % |
| RDBList | 148 | 18 | 166 | 10,8 % |
| Harness | 20 | 13 | 33 | 39,4 % |
| **DecisionTable** | **0** | **10** | 10 | **100 %** |
| **ConnectREST** | **0** | **3** | 3 | **100 %** |
| **Flow** | **0** | **1** | 1 | **100 %** |
| **DecisionTree** | **0** | **1** | 1 | **100 %** |
| DataPage | 26 | 1 | 27 | 3,7 % |
| SystemSettings | 1 | 0 | 1 | 0 % |
| **TOTAL** | **1.195** | **512** | **1.707** | **30,0 %** |

### ⛔ Enam belas berkas yang **seluruh tipenya** berbeda — daftar penuh

`[terverifikasi]` Bukan sampel. Keenam belasnya:

| Tipe | Berkas |
| --- | --- |
| **DecisionTable** (10/10) | `GetMimeType` · **`IsUWAccepted`** · `LicensePlatRegion_DeT` · `MappingAdditionalCoverageIndex` · `MappingCoverage` · `MappingOutgoIndex` · `MappingOutgoIndex2` · `SetUploadHubAW1` · `SetUploadHubAW2` · `SetUploadHubAW3` |
| **ConnectREST** (3/3) | `convertJsonNusareToProduction` · `getPremiumPaidOn` · `ServiceGoogle` |
| **Flow** (1/1) | `OfferFacRetro` |
| **DecisionTree** (1/1) | `Tree_ShortPeriod` |
| **DataPage** (1/27) | `D_EnumerationList` |

⚠️ **`DecisionTable\IsUWAccepted`** dan **`DecisionTree\Tree_ShortPeriod`** menyentuh langsung dua
modul yang sudah terspec untuk NB: **tangga akseptasi** dan **periode pendek/premi**. Keduanya
berbeda di EDM. `[dugaan]` belum dibaca isinya — E-2/E-3.

📌 **`RDBList` 10,8 % dan `DataPage` 3,7 %** — lapisan query dan pengambilan data justru **paling
banyak sama**. Perbedaan EDM terkonsentrasi pada **logika dan tampilan**, bukan pada akses data.

---

## 3. Lingkup `When` yang diukur

`[terverifikasi]` Perbandingan **peka huruf**:

| | Jumlah |
| --- | ---: |
| `When` di NB | 210 |
| `When` di EDM | 202 |
| **Bernama sama** | **175** |
| EDM-only | 27 |
| NB-only | 35 |

⚠️ 11 dari 12 varian kapitalisasi K-042 adalah rule `When`. Perbandingan **tidak peka huruf** memberi
angka berbeda.

✅ **Angka `When` tidak berubah** oleh koreksi metode: 112 identik / 63 berbeda, sama pada metode #2
dan #4. Taksonomi di §4–§6 karena itu tetap berlaku.

---

## 4. Distribusi perbedaan pada 175 rule `When`

`[terverifikasi]`

| Kategori | Jumlah | % |
| --- | ---: | ---: |
| **(c) Identik seluruhnya** | **112** | 64,0 % |
| **(d) Kondisi identik, berbeda di luar kondisi** | **32** | 18,3 % |
| **(a) Beda JALUR properti saja** — logika sama | **19** | 10,9 % |
| **(b) Beda LOGIKA / SUMBER sungguhan** | **12** | 6,9 % |
| **total** | **175** | |

⚠️ Kategori **(d) — 32 berkas** — kondisinya identik tetapi isinya berbeda di tempat lain (label,
kelas, tag non-kondisi). **Belum ditelusuri**; bukan nol, jadi tidak boleh diabaikan.

---

## 5. ⛔ Pola P1 — **terbantah sebagai aturan**, NB sendiri sudah campuran

Warisan arsip berbunyi: *"EDM baca agregat tersimpan `.OfferFacIn.QuotationData`, NB baca halaman
aktif `pyWorkPage.Quotation`."*

`[terverifikasi]` **Premis tentang NB itu salah.** Dibuktikan dari isi tag, bukan dari nama:

| Rule | Jalur di **NB** (`pyConditionValue1String`) |
| --- | --- |
| `NB FacIn\When\IsFire.xml` | `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Fire"` |
| `NB FacIn\When\IsMBU.xml` | `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "MBUCar"` / `"MBUMotorCycle"` |
| `NB FacIn\When\IsMarineCargo.xml` | `.OfferFacIn.QuotationData.BusinessType = "MarineCargo"` **dan** `pyWorkPage.Quotation.BusinessType = "MarineCargo"` |
| `NB FacIn\When\IsAneka.xml` · `IsPA.xml` · `IsGolfInsurance.xml` · `IsTravel.xml` | `pyWorkPage.Quotation.*` |

📌 **NB memakai kedua jalur, kadang keduanya dalam satu rule** (`IsMarineCargo`). "Halaman aktif" dan
"agregat tersimpan" **bukan pembeda siklus** — keduanya hidup berdampingan di dalam NB sendiri.

Arah perpindahan pada 19 rule kategori (a) tetap tercatat, tetapi **sebagai sebaran, bukan aturan**:

| NB → EDM | Jumlah | Rule |
| --- | ---: | --- |
| `pyWorkPage.Quotation` → **`.OfferFacIn.QuotationData`** | **10** | `IsBoiler` · `IsCIS` · `IsCIT` · `IsContractorsPlantMachinery` · `IsEar` · `IsGlass` · `IsGrowingTrees` · `IsHE` · `IsLandRig` · `IsMBD` |
| `pyWorkPage.Quotation` → `.Quotation` | 3 | `IsBondingKBG` · `IsCustomBonds` · `IsMBUCar` |
| `pyWorkPage.Quotation` → `pyWorkPage.OfferFacIn.QuotationData` | 2 | `IsFireStyle1` · `IsFireStyle2` |
| **(relatif)** → `pyWorkPage.Policy` | 1 | `IsCorporate` |
| **(relatif)** → `pyWorkPage.OfferFacIn` | 1 | `IsFacRetro` |
| `pyWorkPage` → **(relatif)** | 1 | `IsProposalTransfer` |
| (tanpa prefiks) → (tanpa prefiks) | 1 | `IsPASSG` |

⛔ **Konsekuensi mengikat:** resolver jalur properti **tidak dapat dibuat generik**, dan **tidak dapat
dibuat per-siklus**. Ia harus **per-rule**, dibaca dari XML.

---

## 6. ✅ Pola P2 — terkonfirmasi, dan lebih dalam dari perbedaan jalur

`[terverifikasi]` Enam predikat lini bisnis (COB) di EDM **tidak membaca properti kasus sama sekali**
— mereka membaca **baris hasil query**. Diverifikasi per path + isi tag, bukan sampel:

> 🔁 **Bagian ini dikoreksi pada E-2** (lihat `03-e2-verifikasi-alur.md` §4.1–§4.2). Versi pertama
> menulis "beda kelas deklarasi rule" — **itu salah**. Teks di bawah sudah diperbaiki.

| Rule | `pxObjClass` NB → EDM | `pyDesignatedClass` EDM (kelas **sumber data**) | Sumber di EDM |
| --- | --- | --- | --- |
| `IsAneka` | `Rule-Obj-When` → `Rule-Obj-When` | **`Code-Pega-List`** | `OutData.pxResults(1).CARI2` |
| `IsFire` | `Rule-Obj-When` → `Rule-Obj-When` | **`Code-Pega-List`** | `OutData.pxResults(1).CARI2` |
| `IsGolfInsurance` | `Rule-Obj-When` → `Rule-Obj-When` | **`Code-Pega-List`** | `OutData.pxResults(1).CARI2` |
| `IsMarineCargo` | `Rule-Obj-When` → `Rule-Obj-When` | **`Code-Pega-List`** | `OutData.pxResults(1).CARI2` |
| `IsMBU` | `Rule-Obj-When` → `Rule-Obj-When` | **`Code-Pega-List`** | `OutData.pxResults(1).CARI2` |
| `IsPA` | `Rule-Obj-When` → `Rule-Obj-When` | **`ASM-FW-GISFW-Data-Search`** | `OutData.pxResults(1).CARI2` |

⛔ **Yang berbeda adalah KELAS SUMBER DATA yang dibaca kondisi, BUKAN kelas deklarasi rule.**
`pxObjClass` keenamnya tetap **`Rule-Obj-When` di kedua sisi** — tipe dan deklarasinya sama.
Di EDM kondisinya membaca **halaman hasil query** (`Code-Pega-List` / `…Data-Search`), di NB membaca
properti kasus. Rule yang sama namanya **membaca sumber yang berbeda**.

📌 Ini **jebakan keenam** keluarga "nama bukan bukti": nama sama, tipe sama, **sumber data berbeda**.

### 6.1 Delegasi vs literal — **tidak seragam**

> 🔁 Versi pertama menyimpulkan "NB mendelegasikan, EDM membandingkan literal" dari dua contoh.
> **Generalisasi itu salah.** Hitungan penuh keenamnya:

`[terverifikasi]`

| Rule | NB delegasi | NB literal/properti | EDM delegasi | EDM literal |
| --- | ---: | ---: | ---: | ---: |
| `IsAneka` | **25** | 1 | 0 | 25 |
| `IsFire` | **4** | 1 | 0 | 5 |
| `IsPA` | 0 | 1 | 0 | 1 |
| `IsGolfInsurance` | 0 | 1 | 0 | 1 |
| `IsMarineCargo` | 0 | 2 | 0 | 1 |
| `IsMBU` | 0 | 2 | 0 | 4 |

📌 Pernyataan yang benar: **NB mendelegasikan HANYA pada `IsFire` dan `IsAneka`**; empat lainnya
membaca properti langsung. EDM memakai literal pada keenamnya.

⛔ **Konsekuensi yang tetap berlaku:** di endorsement, klasifikasi COB **tidak melewati sub-rule sama
sekali**. Registry predikat NB tidak dapat dipakai ulang apa adanya untuk EDM — cabangnya dibaca dari
nilai literal, dan literal itu **harus cocok dengan isi kolom query**, bukan dengan properti kasus.

### 6.2 Dua kejanggalan yang tercatat apa adanya

| Temuan | Bukti | Sikap |
| --- | --- | --- |
| `Endorsment…\When\IsMBU.xml` memuat **cabang ganda** — `"MBUCar"` dan `"MBUMotorCycle"` masing-masing **dua kali** (4 entri `pyConditionValue1String`) | terbaca langsung dari tag | Redundan, tidak mengubah hasil. **Diport apa adanya** (`CLAUDE.md` §1) |
| `Endorsment…\When\IsAneka.xml` memuat 6 `pyConditionString` berlabel `Kode Bisnis = "24"/"18"/…` yang **tidak berkaitan** dengan 25 `pyConditionValue1String`-nya | label tampilan ≠ ekspresi tersimpan | ⚠️ Penegasan `CLAUDE.md` §4.5 — **yang dieksekusi adalah `pyConditionValue1String`**. Label tampilan **tidak boleh** dijadikan dasar perilaku |

⚠️ `[pertanyaan terbuka]` **Mekanisme pengisian `OutData` belum ditelusuri** — query mana yang
mengisinya, dan apa isi kolom `CARI2`. Tanpa itu, predikat COB endorsement tidak dapat
diimplementasikan.

---

## 7. Tiga perbedaan perilaku — diverifikasi ulang per path + isi tag

### `IsUW` — EDM mengakui satu workbasket tambahan `[terverifikasi]`

Kelas **sama** di kedua sisi (`Code-Pega-Requestor`), jadi ini murni perbedaan himpunan kondisi.

| Berkas | `pyConditionValue1String` |
| --- | --- |
| `NB FacIn\When\IsUW.xml` | `…pyWorkBasketName = ReasFacInDirector` · `= ReasFacInGroupLeader` · `= ReasFacInUnderwriting` |
| `Endorsment Fac In\When\IsUW.xml` | `= ReasFacInDirector` · **`= ReasFacInMarketing`** · `= ReasFacInGroupLeader` · `= ReasFacInUnderwriting` |

⛔ **Di endorsement, operator ber-workbasket marketing dihitung sebagai underwriter.** Di NB dan RNW
tidak. Ini menentukan **siapa yang boleh menyetujui** — belum diputuskan.

### `IsClaim` — uji berbeda **dan kelas berbeda** `[terverifikasi]`

| Berkas | Kelas | Kondisi |
| --- | --- | --- |
| `NB FacIn\When\IsClaim.xml` | `pyDesignatedClass` = `ASM-FW-GCNMFW-Work-PNC`, `pxRuleClassName` = `ASM-FW-GISFW-Data` | `pyWorkPage.pyWorkIDPrefix = "CLM-"` — menguji **prefiks ID kasus** |
| `Endorsment Fac In\When\IsClaim.xml` | keduanya = **`ASM-FW-GISFW-Work-Claim`** | `pyWorkPage.ClaimData has a value` — menguji **keberadaan halaman data** |

Bukan jalur berbeda, bukan nilai berbeda: **dua pengujian berbeda sifatnya, pada kelas berbeda.**
Keduanya bisa memberi jawaban berlainan pada kasus yang sama.

### `IsTravel` — dua sumber, dan label tampilan yang tidak sejalan `[terverifikasi]`

| Berkas | Kelas | `pyConditionValue1String` (dieksekusi) | `pyConditionString` (label) |
| --- | --- | --- | --- |
| `NB FacIn\When\IsTravel.xml` | `ASM-FW-GISFW-Work-NB` / `…Data-Quotation` | `pyWorkPage.Quotation.BusinessType = "Travel"` (1) | `BusinessCode = "02"/"58"/"SB"/"SG"` |
| `Endorsment Fac In\When\IsTravel.xml` | `ASM-FW-GISFW-Work` / **`ASM-FW-GCNMFW-Work-PNC`** | `.Quotation.BusinessType = "Travel"` **OR** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Travel"` (2) | `Kode Bisnis = "77"` |

EDM menerima **dua sumber**, NB satu. `[dugaan]` label tampilan (4 kode vs 1 kode) **tidak** dieksekusi
— ia tidak boleh dijadikan bukti perilaku (§4.5).

📌 `IsTravel` **tidak** memakai pola P2 — ia membaca properti, bukan hasil query. P2 berlaku pada 6
rule COB saja, bukan pada seluruh predikat EDM.

---

## 8. Apa yang ini berarti bagi pendekatan EDM

| Pertanyaan E-1 | Jawaban |
| --- | --- |
| Apakah satu pola menjelaskan mayoritas perbedaan? | **Tidak.** P1 **terbantah sebagai aturan** (NB sendiri campuran); P2 nyata tetapi terbatas pada 6 rule COB; sisanya susunan logika + perbedaan perilaku + perbedaan kelas |
| Bisakah 1.707 jadi sedikit keputusan rancangan? | **Sebagian.** 1.195 (70,0 %) identik → dapat diperlakukan bersama NB. **512 berbeda** butuh pemeriksaan |
| Berapa lingkup kerja EDM sebenarnya? | **866 berkas** = 512 berbeda + 354 eksklusif, dari 2.061 — bukan 2.061, bukan 354 |

**Strategi yang muncul dari data:** perlakukan EDM sebagai **NB dengan 866 berkas menyimpang** — bukan
siklus baru (terlalu mahal), bukan varian tipis seperti RNW (terbukti salah).

⚠️ **Angka 866 belum diformalkan.** Ia menunggu keputusan work owner atas amandemen K-042 (E-Q4).

---

## 9. Pertanyaan terbuka

| # | Pertanyaan | Dampak |
| ---: | --- | --- |
| E-Q4 | **Amandemen K-042** — "0 identik" adalah artefak normalisasi; angka terverifikasi **1.195 identik / 512 berbeda**, lingkup kerja **866 berkas** | Menentukan lingkup EDM · **work owner** — **belum diformalkan** |
| E-Q5 | **Bagaimana `OutData.pxResults(1).CARI2` diisi**, dan apa isi kolom `CARI2`? | Tanpa ini, 6 predikat COB endorsement tidak dapat diimplementasikan |
| E-Q6 | **`IsUW`** — di endorsement, workbasket **marketing dihitung sebagai underwriter**. Disengaja? | Menentukan siapa yang boleh menyetujui · **Underwriting** |
| E-Q7 | **`IsClaim`** — NB menguji prefiks ID kasus, EDM menguji keberadaan halaman data; **kelasnya pun berbeda** | Dua uji bisa berbeda hasil pada kasus sama |
| E-Q8 | **`IsTravel`** — EDM menerima dua sumber, NB satu. Disengaja? | Klasifikasi produk travel di endorsement |
| E-Q9 | **32 rule `When` kategori (d)** — kondisi identik tetapi isinya berbeda di luar kondisi | Belum diketahui berdampak atau tidak |
| E-Q10 | **`DecisionTable\IsUWAccepted` dan `DecisionTree\Tree_ShortPeriod` berbeda** — keduanya menyentuh modul NB yang sudah terspec (tangga akseptasi, periode pendek) | **Prioritas E-2/E-3** |
| E-Q11 | **`ConnectREST` 3/3, `Flow` 1/1, `DecisionTable` 10/10 berbeda seluruhnya** — mengapa integrasi, alur, dan tabel keputusan **selalu** menyimpang? | Prioritas E-2 |
| E-Q12 | **Registry predikat NB tidak dapat dipakai ulang untuk COB EDM** (§6.1) — apakah dibuat registry kedua, atau satu registry dengan sumber data yang dapat diganti? | Keputusan rancangan · menunggu E-Q5 |

---

## 10. Belum dikerjakan

- **512 berkas berbeda** di luar `When` — Activity (194), Section (101), FlowAction (45), dst.
- **32 rule `When` kategori (d)** — sebab perbedaannya di luar kondisi
- **27 rule `When` EDM-only** dan **35 NB-only**
- Isi 16 berkas tipe yang 100 % berbeda (§2)
- Mekanisme pengisian `OutData` (E-Q5)

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
