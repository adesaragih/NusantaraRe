# Discovery Endorsement — E-4: 354 berkas EDM-only

> **Sumber:** `D:\migrasi\RNM\Endorsment Fac In\` (banding `NB FacIn\` dan `RNW Fac In\`). READ-ONLY.
> Korpus Treaty (`RNM_BRD\`) **tidak dibaca** (K-005). Label mengikuti `CLAUDE.md` §3.
>
> **EDM-only** = berkas di `Endorsment Fac In\` yang **tidak ada padanan namanya di `NB FacIn\`**.
> Orientasi mencatat 314 dari 354 tidak masuk tema nama apa pun — karena itu pengelompokan di
> dokumen ini dibangun dari **isi**, bukan nama.

---

## 0. Rekonsiliasi angka — mengapa 354, bukan 349

⚠️ Pengukuran pertama saya menghasilkan **349**. Selisihnya **metodologis**, dan angka resmi **354
tetap benar**. Dicatat karena perbedaannya adalah jebakan yang sama yang sudah menjatuhkan pengukuran
lain di proyek ini.

| Metode | Definisi | Hasil |
| --- | --- | ---: |
| **A — per folder** (dipakai orientasi & K-042) | nama tidak ada di folder NB **yang sepadan** | **354** ✅ |
| B — lintas folder | nama tidak ada di **folder NB mana pun** | 349 |

`[terverifikasi]` Kelima berkas selisihnya bernama sama di NB **tetapi tipe rule-nya berbeda**:

| Berkas | EDM | NB |
| --- | --- | --- |
| `ChangeUpKendaraan` | `DataTransform\` — `Rule-Obj-Model` | `Activity\` — `Rule-Obj-Activity` |
| `InputCoverageAneka_FacIn` | `FlowAction\` — `Rule-Obj-FlowAction` | `Section\` — `Rule-HTML-Section` |
| `InputDtlCoverage_FacIn_IsUW` | `FlowAction\` — `Rule-Obj-FlowAction` | `Section\` — `Rule-HTML-Section` |
| `ObjectDtlAneka_FacIn` | `FlowAction\` — `Rule-Obj-FlowAction` | `Section\` — `Rule-HTML-Section` |
| `SelClauseList` | `FlowAction\` — `Rule-Obj-FlowAction` | `Section\` — `Rule-HTML-Section` |

⛔ **Metode A benar** justru karena tipe rule mengikat: sebuah `Rule-Obj-FlowAction` bernama `X`
**bukan** rule yang sama dengan `Rule-HTML-Section` bernama `X`. Ini keluarga jebakan `GetInsuredID`
(nama sama, tipe beda) — kali ini dalam lima kejadian sekaligus.

---

## 1. Klasifikasi 354 menurut `pxObjClass`

`[terverifikasi]` Tipe rule dibaca dari `<pxObjClass>` tiap berkas, **bukan dari nama folder**:

| `pxObjClass` | Jumlah | Folder |
| --- | ---: | --- |
| `Rule-Obj-Activity` | **94** | `Activity\` |
| `Rule-HTML-Section` | **81** | `Section\` |
| `Rule-Obj-FlowAction` | **54** | `FlowAction\` |
| `Rule-Obj-Model` | **29** | `DataTransform\` |
| `Rule-Obj-Report-Definition` | **28** | `ReportDefinition\` |
| `Rule-Obj-When` | **27** | `When\` |
| **`Rule-Connect-SQL`** | **22** | `RDBList\` ⚠️ |
| `Rule-Declare-Pages` | **14** | `DataPage\` |
| `Rule-HTML-Harness` | **4** | `Harness\` |
| `Rule-Obj-Flow` | **1** | `Flow\` |
| **TOTAL** | **354** | |

✅ Cocok persis dengan angka orientasi.

📌 **Jebakan folder ditegaskan, dan sekarang terukur.** Seluruh **22** berkas di `RDBList\` ber-`pxObjClass`
**`Rule-Connect-SQL`**, bukan `Rule-Obj-RDBList`. Temuan E-2 (`InsertTreatyProduction_Sql`) dan E-3
(`GetEDMOldData_SQL`) **bukan pengecualian — itu aturannya**. Folder `RDBList\` di korpus ini memuat
rule Connect-SQL.

⚠️ Di **dalam** himpunan 354, folder konsisten dengan `pxObjClass` (0 anomali). Ketidaksepadanannya
muncul saat **membandingkan antar-korpus** — itulah sebab §0.

```powershell
# menghasilkan tabel di atas
$nbR="D:\migrasi\RNM\NB FacIn"; $edR="D:\migrasi\RNM\Endorsment Fac In"
function ObjClassOf($p){ $oc=''
  $rd=[System.Xml.XmlReader]::Create($p)
  while($rd.Read()){ if($rd.NodeType -eq 'Element' -and $rd.Name -eq 'pxObjClass'){
    $oc=$rd.ReadElementContentAsString(); break } }
  $rd.Close(); $oc }
$set=@()
foreach($dir in (Get-ChildItem $edR -Directory)){
  $nbF=New-Object 'System.Collections.Generic.HashSet[string]'   # peka huruf
  Get-ChildItem "$nbR\$($dir.Name)" -Filter *.xml -File -EA SilentlyContinue|%{[void]$nbF.Add($_.BaseName)}
  foreach($f in (Get-ChildItem $dir.FullName -Filter *.xml -File)){
    if(-not $nbF.Contains($f.BaseName)){ $set += [pscustomobject]@{Nama=$f.BaseName;ObjClass=(ObjClassOf $f.FullName)} } } }
$set.Count; $set|Group-Object ObjClass|Sort-Object Count -Descending|Select Count,Name
```

---

## 2. Kelompok fungsi — 94 Activity EDM-only

`[terverifikasi]` Dikelompokkan dari **isi** (cacah token properti + metode langkah + kelas), bukan
nama. Ⓟ menandai kelompok yang **menyentuh perhitungan/premi** — prioritas.

| Kelompok | Jml | Dasar pengelompokan (tag/token yang dihitung) |
| --- | ---: | --- |
| **1. Before-image / OldData** Ⓟ | 9 | `OldData` / `IsOldData` / `FlagOldData` ≥ 20 kemunculan |
| **2. Spreading & kapasitas treaty** Ⓟ | 7 | token `Spread` ≥ 20 |
| **3. Perhitungan premi / TSI / rate** Ⓟ | 8 | `Premi`+`TSI`+`.Rate` ≥ 60 |
| **4. Jalur Life** Ⓟ | 7 | token `Life` ≥ 10 |
| **5. Jalur retro / fac out** | 2 | token `Retro` ≥ 10 |
| **6. Integrasi / notifikasi** | 1 | `Connect-REST`/`Email`/`Correspond` ≥ 5 |
| **7. Manipulasi coverage / objek & kait layar** | 60 | sisanya — lihat §2.1 |

### Kelompok 1 — Before-image / OldData Ⓟ `[terverifikasi]`

| Berkas | Langkah | Sinyal `OldData` |
| --- | ---: | ---: |
| `Activity\SetValueToEDMWork` | 46 | 85 |
| `Activity\GetdataOldEDMError` | 146 | 115 |
| `Activity\SetOLDValueToEDMWork_{FIRE,MC,Aneka,MBU,LIFE,PA,GOLF}` | 8–24 | 11–19 |

📌 Menyatu dengan E-3: kesembilannya **adalah** lapis A dan lapis C. Seluruh mesin before-image
adalah kode **EDM-only** — tidak ada padanannya di NB. Itu memperkuat K-043: before-image bukan
"NB dengan tambahan", melainkan subsistem tersendiri.

### Kelompok 2 — Spreading & kapasitas treaty Ⓟ `[terverifikasi]`

| Berkas | Langkah | premi / TSI / spread |
| --- | ---: | --- |
| `Activity\CopySpreading_Act` | 95 | 52 / 54 / **247** |
| `Activity\SpreadingProtection` | 20 | 0 / 0 / **96** |
| `Activity\CheckLimitTreatyType_Act` | 43 | 5 / **88** / 20 |
| `Activity\CountEndorsementData` | 50 | **77** / 35 / 25 |
| `Activity\CountDataEDMElse` | 42 | **118** / 28 / 27 |
| `Activity\CountEdmExtPeriode_Act` | 46 | 51 / 43 / 21 |
| `Activity\ReCountPremiLifeEDM` | 22 | 63 / 51 / 24 |

⛔ **`CountEndorsementData`, `CountDataEDMElse`, dan `ReCountPremiLifeEDM` adalah EDM-only** — ketiga
berkas yang seluruh tabel rumus `EDMPremiMenjadi` (E-3 §5) bersandar padanya. Perhitungan premi
endorsement **tidak mewarisi apa pun dari NB**.

### Kelompok 3 — Perhitungan premi / TSI / rate Ⓟ `[terverifikasi]`

`Activity\SetTSIAllCoverageLife_ACT` (TSI 174, rate 66) · `Activity\SetLocalNonMbuProrate`
(premi 50, TSI 96, rate 97) · `Activity\SetSumTSISpreading_Act` (TSI 92) ·
`Activity\CopyTemplateCoverageSQL_PostActEDM` (77 langkah) · `Activity\CalculatePremiumTravel` ·
`Activity\CopyTemplateCoverage` · `Activity\calculatePremiPA` · `Activity\CalculatePremiFire`

⚠️ **Empat aktivitas ber-nama `Calculate*`/`calculate*` adalah EDM-only** — perhitungan premi
per-lini di endorsement punya implementasi sendiri. `[pertanyaan terbuka]` apakah rumusnya setara
dengan jalur NB; **belum dibandingkan** (bukan lingkup E-4).

### Kelompok 4 — Jalur Life Ⓟ

`Activity\CalculateScorLife_Act` (**210 langkah** — terpanjang di seluruh himpunan 354) ·
`setRateLife_act` · `SetTSIAllCoverage_ACT` · `AddCurrencyListLife_ACT` ·
`CountCoveragePropertyLife_Act` · `copyCoveragePeriode_act` · `AddAgeLife_ACT` — lihat §5.

### 2.1 Kelompok 7 — 60 berkas, dipilah lebih lanjut `[dugaan]`

Sisanya bukan "tak terklasifikasi" melainkan **berbobot perhitungan rendah**. Pola yang terlihat dari
isi:

| Sub-pola | Contoh | Sifat |
| --- | --- | --- |
| Salin/terapkan ke seluruh objek | `CopyDiscountAllObject_act` (premi 46) · `CopyCoverageToALL_HIO_Act` · `CopyDeductibleToAllCargo_Act` | Ⓟ ringan — menyalin angka |
| Kait layar (`_PreAct`/`_PostAct`) | `InputCoverageAneka_PreAct` · `CopyTemplateCoverageSQL_PostAct` · `CoverageList_FacIn_PreAct` | tampilan |
| Gerbang/validasi masuk | `SetErrorBatalEndorsement_Act` (45 langkah) · `CheckEDMPolisDate` · `CekDoubleCoverage` | kontrol alur |
| Medis / benefit (Life) | `PreInputMedial_Act` · `SaveMedical` · `SetParamLab_Act` · `CalculatePhysicalExam` | lihat §5 |
| Master/lookup | `BrowseActV_VEHICLE_BRAND` · `GetBusinessGroup_Act` · `GetInwardScale_ACT` | data |

⚠️ Dilabeli `[dugaan]` karena pemilahan ini dari pola isi agregat, **belum** dibaca langkah-per-langkah.

---

## 3. Kelompok fungsi — 81 Section EDM-only

`[terverifikasi]` Dikelompokkan dari cacah token isi (`Coverage`, `Premi`, `TSI`, `Life`/`Benefit`,
`Old`, `pyActivity`/`pyLocalAction`):

| Kelompok | Jml | Contoh + sinyal |
| --- | ---: | --- |
| **C. Coverage / objek pertanggungan** | **21** | `Section\InputPerCoverageFire_IsUW` (cov 424) · `InputCoverageTravel` (cov 179) |
| **A. Life (benefit/medis)** Ⓟ | **14** | `Section\Medical_Sec` (life 1592) · `InputMedical2_Sec` (1227) · `InputCoverageLife_FacIn` (cov 293) |
| **E. Layar beraksi (tombol/activity)** | 7 | `Section\InputPackageMedi` · `InpMedPackage` · `PeriodeEndorsement_IsUW` |
| **B. Menampilkan/mengubah angka premi-TSI-rate** Ⓟ | 6 | `Section\InputCoverageMBU_FacIn` (premi 47, cov 664) · `InputCoverageGridMBU_FacIn` (old 32) |
| **D. Menampilkan nilai lama (before-image)** Ⓟ | 2 | `Section\periode` (old 16) · `Section\InputEndorsement_IsUW` (old 10) |
| **F. Layar/komponen umum** | 31 | `Section\ViewRiCommFacOut1` · `PlanSection` · `ObjectDtlAneka_FacIn_IsUW` |

⛔ **Hanya 2 section menampilkan nilai lama.** Lapis B mengisi 52 properti `*Old` (E-3 §3.3), tetapi
yang menampilkannya hanya `Section\periode` dan `Section\InputEndorsement_IsUW`.
`[pertanyaan terbuka]` apakah `*Old` lain dipakai untuk **perhitungan saja**, bukan tampilan — belum
ditelusuri.

📌 `Section\periode` adalah salah satu dari **12 varian kapitalisasi** K-042 (`periode` ↔ `Periode`).

---

## 4. Peta pemanggilan — dan mengapa "tak dirujuk" hampir salah total

### 4.1 Dua mekanisme, bukan satu

`[terverifikasi]` Korpus Pega memuat indeks rujukan bawaan `<pxRuleReferences REPEATINGTYPE="PageList">`
pada **1.872 dari 2.061** berkas EDM — **46.254 sisi rujukan**. (Divalidasi: `SetValueToEDMWork`
menghasilkan **136**, cocok dengan cacah `rowdata`-nya.)

| Tahap | Hasil |
| --- | ---: |
| Dirujuk menurut indeks `pxRuleReferences` | **203** dari 354 |
| Tampak **tidak** dirujuk setelah tahap 1 | 151 |
| Dari 151 itu, **ternyata disebut** berkas lain (mekanisme kedua) | **146** |
| **Benar-benar tidak disebut di mana pun** | **5** |

⛔ **Kalau berhenti di tahap 1, saya akan melaporkan 151 rujukan menggantung — salah 146 di
antaranya.** `pxRuleReferences` **tidak** mengindeks pemanggilan berbasis **string nama saat
runtime**. Yang terlewat, dengan bukti tag:

| Tipe | Terlewat | Mekanisme sebenarnya (kutipan konteks) |
| --- | ---: | --- |
| `Rule-Connect-SQL` | **20 dari 22** | `Activity\GetBusinessGroup_Act` → `<RequestType>CariBusinessGID</RequestType>` |
| `Rule-Obj-Activity` | 52 | `Section\EmailSectionEDM` → `<pyActivity>AcceptNotificationEdm_Post</pyActivity>` |
| `Rule-Obj-FlowAction` | 24 | `Section\InputPerCoverageFire_IsUW` → `<pyLocalAction>ChooseAccumulation</pyLocalAction>` |
| `Rule-Obj-Model` | 16 | `Section\InputPerCoverageFire_IsUW` → `<pyName>BlankChangeSublimit_DT</pyName>` dalam `Embed-Invoke-DataTransform` |
| `Rule-Obj-Report-Definition` | 13 | `Section\ViewCoverageCargo` → `<pySourceName>BrowseCoverageCargoList</pySourceName>` |

📌 Ini penerapan langsung **K-006**: ketiadaan dalam satu indeks **bukan bukti** tak terpakai.

### 4.2 Siapa memanggil apa

`[terverifikasi]` Sebaran tipe **pemanggil** untuk 203 yang terindeks:

| Pemanggil | Sisi rujukan |
| --- | ---: |
| `Section\` | 186 |
| `Activity\` | 69 |
| `FlowAction\` | 58 |
| `DataPage\` | 11 |
| `When\` | 10 |
| `Flow\` | 3 |
| `Harness\` | 2 |

⛔ **Section adalah pemanggil terbesar.** Delta EDM digerakkan dari **lapisan tampilan**, bukan dari
flow. Konsisten dengan orientasi (EDM menambah DataPage +11, ReportDefinition +15, FlowAction +15
dibanding NB) — dan menjelaskan mengapa hanya ada 2 Flow.

### 4.3 Lima berkas yang benar-benar tidak disebut `[pertanyaan terbuka]`

| Folder | Berkas | `pxObjClass` |
| --- | --- | --- |
| `RDBList\` | `GetMasterKlausulAge(1)` | `Rule-Connect-SQL` |
| `RDBList\` | `SearchClobClauseSQL(1)` | `Rule-Connect-SQL` |
| `When\` | `IsCustomBond` | `Rule-Obj-When` |
| `Activity\` | `SaveClausePA_Act` | `Rule-Obj-Activity` |
| `Harness\` | `SFAPortalEndorsement` | `Rule-HTML-Harness` |

⛔ **Tidak dinyatakan usang** (`CLAUDE.md` §4.5, K-006). Dua di antaranya bernama berakhiran **`(1)`**
— `[dugaan]` pola salinan "save as" Pega, tetapi **nama bukan bukti**; belum diperiksa isinya.
`SFAPortalEndorsement` adalah **Harness** (titik masuk layar) yang tidak dipanggil dari korpus EDM —
`[dugaan]` dipanggil dari konfigurasi portal di luar ekspor. → **E-Q24**.

```powershell
# validasi indeks rujukan: menghasilkan 136
$rxB=[regex]'(?s)<pxRuleReferences[^>]*>(.*?)</pxRuleReferences>'
$rxP=[regex]'(?s)<pyRuleName>([^<]*)</pyRuleName>.*?<pxRuleObjClass>([^<]*)</pxRuleObjClass>'
$t=[IO.File]::ReadAllText("D:\migrasi\RNM\Endorsment Fac In\Activity\SetValueToEDMWork.xml")
$rxP.Matches($rxB.Match($t).Groups[1].Value).Count
```

---

## 5. Jalur Life — subsistem tersendiri

`[terverifikasi]` **52 berkas EDM-only** ber-sinyal Life, tersebar di **9 tipe rule**:

| Tipe | Jml | Contoh |
| --- | ---: | --- |
| `Activity` | 17 | `CalculateScorLife_Act` (**210 langkah**) · `setRateLife_act` · `ReCountPremiLifeEDM` |
| `Section` | 15 | `Medical_Sec` · `InputMedical1_Sec` · `InputMedical2_Sec` · `BenefitList` |
| `FlowAction` | 8 | `HistoryofDisease` · `InpBenefit` · `UploadCSVEDM_BenefitLimit` |
| `ReportDefinition` | 5 | `BrowseJBenefit` · `BrowseBenefitLife_RD` · `BrowseBenefitClause_RD` |
| `DataPage` | 3 | `D_BenefitList` · `D_J_BenefitList` · `D_BenefitPkg` |
| `Harness` | 1 | **`Medical_Harnes`** |
| `RDBList` | 1 | `BrowseLifeRisk_SQL` (`Rule-Connect-SQL`) |
| `DataTransform` | 1 | `DelSubHistoryDisease` |

⚠️ Penyaring saya juga menangkap `Flow\InputAddendumFacultativeIn` (sinyal 74) — itu **flow utama**,
bukan milik Life. Dikecualikan dari kelompok; dicatat agar ambang tidak disalahbaca.

### Ia benar-benar membentuk sub-graf

`[terverifikasi]` **31 sisi rujukan Life→Life**, contoh:

```
Section\CreateListBenefit      -> Section\BenefitList , Section\SelectBenefitList
Section\BenefitList            -> DataPage\D_BenefitList , DataPage\D_J_BenefitList
Section\BenefitClauseList      -> ReportDefinition\BrowseBenefitClause_RD
Section\Medical_Sec            -> FlowAction\HistoryofDisease
Section\InputCoverageLife_FacIn-> FlowAction\InputDtlCoverageLife_FacIn
Activity\SetTSIAllCoverageLife_ACT -> Activity\CountCoveragePropertyLife_Act
```

⛔ **Life punya Harness sendiri (`Medical_Harnes`), DataPage sendiri, ReportDefinition sendiri, dan
rantai Section→FlowAction→Activity sendiri.** Ini bukan variasi kecil — ini **aplikasi di dalam
aplikasi**.

### Menyambung tiga temuan sebelumnya

| Temuan | Asal | Menyatu jadi |
| --- | --- | --- |
| `SetOldData` **keluar** bila `IsLife` (gerbang 1) | E-3 §3.1 | Life tidak memakai lapis B |
| **Hanya Life** menerima `EdmType==1` pada cabang perhitungan | E-3 §5.3 (E-Q21) | Life punya aturan pembatalan sendiri |
| Life **melompati seluruh tangga akseptasi**; `GetLimitAkseptasiLife_Act` tidak ada di korpus | arsip §4.5 (#6, MEMBLOKIR) | Life tidak memakai mesin akseptasi |
| 52 berkas EDM-only ber-sub-graf lengkap | **E-4 (ini)** | Life punya lapisan tampilan & data sendiri |

⛔ **Keempatnya searah:** endorsement Life adalah **jalur terpisah**, bukan varian lini bisnis.
`[pertanyaan terbuka]` apakah Life masuk lingkup migrasi tahap ini → **E-Q26**. Ini keputusan work
owner, bukan kesimpulan teknis.

---

## 6. Enam berkas EDM+RNW (bukan NB) — ⛔ orientasi perlu dikoreksi

`[terverifikasi]` Diuji dengan **skrip yang sama persis** dengan §7.6 — **dua varian**, satu standar,
bukan dua:

| Berkas | Tipe | Varian A: **23-tag** | Varian B: **+ buang `pzIndexes`** | Beda hanya indeks? |
| --- | --- | --- | --- | --- |
| `Activity\GetBusinessGroup_Act` | `Rule-Obj-Activity` | ✅ **IDENTIK** | ✅ IDENTIK | — |
| `RDBList\CariBusinessGID` | `Rule-Connect-SQL` | ✅ **IDENTIK** | ✅ IDENTIK | — |
| `ReportDefinition\BrowseAccountInsuredEDM` | `Rule-Obj-Report-Definition` | ✅ **IDENTIK** | ✅ IDENTIK | — |
| `Activity\SetDataInsuredEDM_Act` | `Rule-Obj-Activity` | ✅ **IDENTIK** | ✅ IDENTIK | — |
| `Harness\ChooseInsured` | `Rule-HTML-Harness` | BEDA | ✅ IDENTIK | **YA** |
| `Section\ChooseInsuredDtl` | `Rule-HTML-Section` | BEDA | ✅ IDENTIK | **YA** |

⛔ **Penting untuk konsistensi:** keempat berkas pertama **identik tanpa perlu membuang `pzIndexes`
sama sekali** — di bawah kontrak 23-tag yang sedang berlaku. Perluasan §7 hanya dibutuhkan untuk dua
berkas terakhir. Jadi klaim "4 dari 6 identik" **tidak bergantung** pada E-Q23.

> 🔁 **Koreksi orientasi.** `01-orientasi-discovery-edm.md` §3 menyatakan *"Keenamnya **berbeda
> isinya** antara EDM dan RNW"*. Itu diperoleh dari **hash mentah**, yang selalu berbeda karena
> stempel waktu **dan urutan tag ekspor** (E-3 §3.6). Dengan kontrak 23-tag yang berlaku:
> **4 dari 6 identik**.

### Dan dua yang "berbeda" pun 100 % penomoran indeks

`[terverifikasi]` Cacah baris pembeda di bawah kontrak 23-tag, dipilah bentuknya:

| Berkas | Total baris beda | Berbentuk `rowdata` indeks | **Bukan** indeks |
| --- | ---: | ---: | ---: |
| `Harness\ChooseInsured` | 37 | **37** | **0** |
| `Section\ChooseInsuredDtl` | 58 | **58** | **0** |

Seluruhnya berbentuk `<rowdata REPEATINGINDEX="RuleReference|NamedPageReference|Warnings">N</rowdata>`
di dalam blok `<pzIndexes REPEATINGTYPE="PropertyGroup">` (§7.1–§7.3). **Nol** baris berupa properti,
kondisi, kelas, ekspresi, atau nilai.

⛔ **Konsekuensi untuk E-Q3 (nasib pilih-tertanggung EDM):** dasar orientasi untuk *"K-038 tidak
otomatis berlaku untuk EDM karena isinya beda"* **tidak bertahan**. Keempat berkas pilih-tertanggung
di EDM **setara fungsional** dengan salinan RNW. Tetapi **apakah K-038 berlaku untuk EDM tetap
keputusan work owner** — E-4 hanya mencabut premis teknisnya. → **E-Q25**.

---

## 7. ⛔ E-Q23 — blok indeks turunan `pzIndexes` belum ada di kontrak 23-tag

> 🔁 **Koreksi istilah & metode.** Versi pertama menarasikan "`pzIndexes`" tetapi **skripnya**
> menyaring berdasarkan **atribut `REPEATINGINDEX` pada `rowdata`** — dua definisi berbeda untuk satu
> klaim. Bagian ini menggantinya dengan **satu definisi struktural**, dibuktikan pada berkas yang
> ditunjuk work owner, dan menghasilkan angka yang sama.

### 7.1 Nama tag XML persis — dan mengapa pencariannya bisa menghasilkan 0

`[terverifikasi]` Nama elemennya memang **`pzIndexes`**. Pencarian literal `<pzIndexes>` menghasilkan
**0** karena elemen ini **tidak pernah diserialisasi tanpa atribut**:

```
<pzIndexes REPEATINGTYPE="PropertyGroup">
```

`[terverifikasi]` Di `Activity\GetLimitAkseptasi_Act.xml`:

| Pola dicari | NB | EDM |
| --- | ---: | ---: |
| `<pzIndexes>` (literal, berpenutup) | **0** | **0** |
| `<pzIndexes` (tanpa penutup `>`) | **82** | **82** |
| `REPEATINGINDEX` | 820 | 820 |

⚠️ Ini **jebakan yang sama** yang menjatuhkan regex saya di E-4 §4.1 (`<pxRuleReferences
REPEATINGTYPE="PageList">`). Tag ber-atribut tidak tertangkap pencarian literal berpenutup. Pencarian
wajib memakai `<pzIndexes` **tanpa** `>`.

### 7.2 Satu contoh berdampingan — NB vs EDM

`[terverifikasi]` `Activity\GetLimitAkseptasi_Act.xml`, rantai induk dan **nomor baris identik** di
kedua korpus:

```
        NB                                     EDM
L11644  <pxRuleReferences>                     <pxRuleReferences>
L11645    <rowdata>                              <rowdata>
L11653      <pzIndexes REPEATINGTYPE="PropertyGroup">   <pzIndexes REPEATINGTYPE="PropertyGroup">
L11654        <rowdata REPEATINGINDEX="RuleReference">86</rowdata>
                                                  <rowdata REPEATINGINDEX="RuleReference">6</rowdata>
```

Setelah normalisasi 23-tag, berkas ini berbeda **82 baris di tiap sisi**, dan **seluruh 82-nya
berbentuk sama**:

| | NB memakai indeks | EDM memakai indeks | Cacah |
| --- | --- | --- | ---: |
| `REPEATINGINDEX="RuleReference"` | 86 … 159 | 6 … 79 | **74 : 74** |
| `REPEATINGINDEX="NamedPageReference"` | 82 … 85 | 2 … 5 | **4 : 4** |
| `REPEATINGINDEX="Warnings"` | 80, 81 | 0, 1 | **4 : 4** |

⛔ **Cacahnya sama persis; hanya titik awal penomoran bergeser +80.** Tidak ada satu pun baris berupa
properti, kondisi, kelas, ekspresi, atau nilai. Ukuran berkas beda **78 byte dari 544.882 (0,014 %)**
— konsisten dengan selisih lebar digit, bukan isi.

### 7.3 Bukti bahwa `pzIndexes` tidak pernah memuat isi fungsional

`[terverifikasi]` Dipindai di **seluruh** korpus EDM:

| Ukuran | Nilai |
| --- | ---: |
| Berkas ber-blok `pzIndexes` | **1.872** dari 2.061 |
| Total blok `pzIndexes` | **59.143** |
| Elemen di dalamnya yang **bukan** `rowdata` | **0** |
| Nilai di dalamnya yang **bukan bilangan bulat** | **0** |

Nama `REPEATINGINDEX` yang muncul: `RuleReference` (51.251) · `Warnings` (4.163) ·
`NamedPageReference` (3.510) · `pzParameters` (72) · `pzCSS` (41) · `pzDataSources` (41) ·
`pzDataReference` (20) · `WorkBasket` (16) · `Privileges` (3) · `IndexCustomFields` (2) ·
`AddWorkProcesses` (2).

✅ **Tidak ada tag pembawa kelas di dalam `pzIndexes`** — `pxRuleClassName`, `pxRuleFamilyName`, dan
`pyDesignatedClass` **tidak tersentuh** oleh penyaring ini. Aturan sejak E-1 tetap utuh.

### 7.4 Hitung menyeluruh 1.707 pasangan — angka apa adanya

`[terverifikasi]`

| Perlakuan | Berbeda | Identik |
| --- | ---: | ---: |
| **Kontrak 23-tag sekarang** (dasar K-042/K-043) | **512** | **1.195** |
| 23-tag **+ buang isi blok `pzIndexes`** | **354** | **1.353** |
| **Beda HANYA karena blok `pzIndexes`** | **158** | |

**Uji silang wajib — seluruhnya tetap lolos di bawah metode baru:**

| Berkas | Harus | Hasil |
| --- | --- | --- |
| `When\FlagOldData` (beda kelas) | **BEDA** | ✅ BEDA |
| `When\IsAdmin` · `IsBonding` · `IsAddButton` | identik | ✅ identik |
| `DataPage\D_AnekaList` | identik | ✅ identik |

⚠️ **158 berkas** berbeda semata karena penomoran indeks. Di antaranya yang menentukan perilaku inti:

```
Activity\GetLimitAkseptasi_Act          ← mesin tangga akseptasi
Activity\GetLimitAkseptasi_JUW_UW       ← mesin tangga akseptasi
Activity\CopyAllObj_ACT
Activity\CountPremiAndTSINusantaraRe_ACT
Activity\CheckSpreadingProtect_ACT
Activity\GetAksepBanding
Activity\GetKapasitasTreaty
```

⚠️ **Kebetulan yang berbahaya:** angka beda-setelah-koreksi (**354**) **sama persis** dengan cacah
EDM-only (**354**). Keduanya **tidak berhubungan**. Dicatat agar tidak tertukar.

### 7.5 Usulan (bukan keputusan)

⛔ **Saya tidak mengubah kontrak.** Kontrak 23-tag adalah dasar **K-042 dan K-043**.

**Usulan:** amandemen kontrak menjadi **"23 tag volatile + isi blok `pzIndexes` diabaikan"**, dengan
alasan: isinya **terbukti** hanya penanda posisi bilangan bulat (§7.3), tidak memuat tag pembawa
kelas, dan kelima uji silang tetap lolos (§7.4). Konsekuensinya lingkup K-043 **866 → 708**
(354 beda + 354 EDM-only).

**Keputusan milik work owner.**

### 7.6 Skrip reproducible — dapat dijalankan ulang

```powershell
$nbR="D:\migrasi\RNM\NB FacIn"; $edR="D:\migrasi\RNM\Endorsment Fac In"
$VOL=@('pxCreateDateTime','pxUpdateDateTime','pxSaveDateTime','pxCommitDateTime','pxMoveImportDateTime',
 'pxOriginalCreateDateTime','pyRuleFormStatusTime','pxWarningCreatedTime','pxCreateOperator','pxCreateOpName',
 'pxCreateSystemID','pxUpdateOperator','pxUpdateOpName','pxUpdateSystemID','pxMoveImportOperId',
 'pxMoveImportOperName','pxOriginalCreateOperator','pxOriginalCreateOpName','pxOriginalCreateSystemID',
 'pxHostId','pzChecksum','pzIndexCount','pyShowJavaWindowName')
$KEY=@('pzInsKey','pzIndexOwnerKey','pzOriginalInstanceKey','pzDocumentKey','pyJavaClassName','pxInsName')
$rxV=[regex]('^</?(' + ($VOL -join '|') + ')[ />]|^<(' + ($VOL -join '|') + ')>')
$rxK=[regex]('^<(' + ($KEY -join '|') + ')>')
$rxT=[regex]'\d{8}T\d{6}(\.\d+)?( GMT)?|_\d{8}T\d{6}_\d+'
# --- penyaring STRUKTURAL: buang seluruh isi blok pzIndexes (tag selalu ber-atribut) ---
$rxO=[regex]'^<pzIndexes(\s[^>]*)?>$'; $rxS=[regex]'^<pzIndexes(\s[^>]*)?/>$'; $rxC=[regex]'^</pzIndexes>$'
function Sig([string]$p,[bool]$buangIdx){
  $l=New-Object Collections.Generic.List[string]; $d=0
  foreach($ln in [IO.File]::ReadAllLines($p)){ $tr=$ln.Trim(); if(-not $tr){continue}
    if($buangIdx){
      if($rxS.IsMatch($tr)){continue}
      if($rxO.IsMatch($tr)){$d++;continue}
      if($rxC.IsMatch($tr)){if($d -gt 0){$d--};continue}
      if($d -gt 0){continue} }
    if($rxV.IsMatch($tr)){continue}
    if($rxK.IsMatch($tr)){$tr=$rxT.Replace($tr,'<TS>')}   # buang stempel waktu, PERTAHANKAN kelas
    $l.Add($tr) }
  $a=$l.ToArray(); [Array]::Sort($a,[StringComparer]::Ordinal); [string]::Join("`n",$a) }
$tot=0;$bedaA=0;$bedaB=0;$hanyaIdx=@()
foreach($dir in (Get-ChildItem $edR -Directory)){
  $nbF=New-Object 'System.Collections.Generic.HashSet[string]'          # peka huruf
  Get-ChildItem "$nbR\$($dir.Name)" -Filter *.xml -File -EA SilentlyContinue|%{[void]$nbF.Add($_.BaseName)}
  foreach($f in (Get-ChildItem $dir.FullName -Filter *.xml -File)){
    if(-not $nbF.Contains($f.BaseName)){continue}
    $tot++; $pn="$nbR\$($dir.Name)\$($f.BaseName).xml"
    if((Sig $pn $false) -eq (Sig $f.FullName $false)){continue}
    $bedaA++
    if((Sig $pn $true) -ne (Sig $f.FullName $true)){$bedaB++} else {$hanyaIdx+="$($dir.Name)\$($f.BaseName)"} } }
"total=$tot  beda23=$bedaA  beda23+idx=$bedaB  hanyaIndeks=$($hanyaIdx.Count)"
# => total=1707  beda23=512  beda23+idx=354  hanyaIndeks=158
```

---

## 8. Pertanyaan terbuka baru dari E-4

| # | Pertanyaan | Pemilik | Dampak |
| ---: | --- | --- | --- |
| **E-Q23** ✅ **terbukti reproducible** | **Isi blok `pzIndexes` masuk kontrak volatile?** Elemen bernama **`pzIndexes`**, selalu ber-atribut `REPEATINGTYPE` (pencarian `<pzIndexes>` literal → 0; pakai `<pzIndexes` tanpa `>`). 59.143 blok di korpus EDM, isinya **0 %** non-`rowdata` dan **0 %** non-bilangan-bulat. **158 dari 512** beda hanya karenanya → lingkup K-043 **866 → 708**. Kelima uji silang tetap lolos. **Usulan** §7.5; keputusan milik work owner | **work owner** | Menggeser keputusan yang sudah diformalkan (K-042, K-043) |
| **E-Q24** | **5 berkas tidak disebut di mana pun** di korpus EDM: `GetMasterKlausulAge(1)` · `SearchClobClauseSQL(1)` · `IsCustomBond` · `SaveClausePA_Act` · `SFAPortalEndorsement`. Per K-006 **tidak** dinyatakan usang | work owner + IT | Kandidat rujukan menggantung; `SFAPortalEndorsement` adalah Harness — mungkin dipanggil dari konfigurasi portal di luar ekspor |
| **E-Q25** | **Apakah K-038 (pilih-tertanggung dibuang) berlaku untuk EDM?** Premis orientasi (*"isinya beda dari RNW"*) **tidak bertahan** — 2 dari 4 identik, 2 sisanya beda hanya nomor indeks | **work owner** | Menutup E-Q3 lama; menentukan 4 berkas masuk lingkup atau tidak |
| **E-Q26** | **Apakah endorsement Life masuk lingkup migrasi tahap ini?** 52 berkas EDM-only ber-sub-graf lengkap (Harness, DataPage, ReportDefinition sendiri), melompati tangga akseptasi, tidak memakai lapis B, dan satu-satunya yang jalan di `EdmType==1` | **work owner** | Bisa memangkas lingkup EDM secara berarti, atau menambah subsistem penuh |
| **E-Q27** | **Hanya 2 Section menampilkan nilai `*Old`**, padahal lapis B mengisi 52 properti. Apakah sisanya hanya untuk perhitungan? | work owner + Underwriting | Menentukan kebutuhan tampilan before-image |

---

## 9. Belum dikerjakan (bukan untuk prompt ini)

- **E-5** — jalur produksi endorsement (`_MENJADI`/`_SELISIH`, `PRODKE`). **Tidak dimulai tanpa perintah.**
- Isi langkah-per-langkah 60 Activity kelompok 7 (§2.1 masih `[dugaan]`)
- 54 FlowAction, 29 DataTransform, 28 ReportDefinition, 27 When, 22 Connect-SQL, 14 DataPage
  EDM-only — **belum dikelompokkan per fungsi**
- Banding rumus `Calculate*` EDM-only vs jalur NB (§2 kelompok 3)
- Isi 5 berkas tak-disebut (§4.3) — belum dibaca
- Arsip §4 (gerbang akseptasi EDM) — termasuk klaim "nilai dasar akseptasi = SELISIH TSI"

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
