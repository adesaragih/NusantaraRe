# 02-layar / 01 — Skema field layar New Business (NB FacIn)

**Sumber tunggal:** `D:\migrasi\RNM\NB FacIn\` — ekspor rule Pega, dibaca read-only.
Tidak ada bahan dari `RNW Fac In\`, `Endorsment Fac In\`, maupun `RNM_BRD\`.
**Tanggal audit:** 2026-09-15. Setiap angka disertai perintah PowerShell yang menghasilkannya.

**Hubungan dengan dokumen sebelumnya.** `OUTPUT/_ARSIP-lintas-siklus/03-ui/01-inventaris-layar.md`
menjawab *layar apa saja yang ada* (lintas tiga siklus). Dokumen ini menjawab pertanyaan berikutnya
dan hanya untuk NB: **field apa yang diikat di tiap layar, kontrolnya apa, kapan tampil, kapan
dapat disunting.** Angka di sini dihitung ulang khusus NB dan **tidak boleh** dibandingkan langsung
dengan angka lintas-siklus di dokumen arsip.

---

## Daftar isi

| § | Isi |
| --- | --- |
| 0 | Metode: tag mana yang dipakai sebagai bukti pengikatan, dan mengapa |
| 1 | Angka pokok — 432 Section, 6.314 field terikat |
| 2 | Skema field per Section (ringkasan; tabel penuh di Lampiran A) |
| 3 | Katalog properti unik — tulang punggung model data |
| 4 | Rumpun lini bisnis dari properti terikat |
| 5 | Repeating grid |
| 6 | Field read-only dan field tersembunyi permanen |
| 7 | Validasi yang terbaca — dan yang tidak ada sama sekali |
| 8 | Tipe data: apa yang terbukti dan apa yang tidak |
| 9 | [usulan] pemetaan ke React + Go |
| 10 | Pertanyaan terbuka (dengan penanda MEMBLOKIR) |
| Lampiran A | Tabel penuh: Section → properti → kontrol → read-only → kondisi tampil (432 Section) |
| Lampiran B | Katalog 1.805 ekspresi properti unik, diurutkan menurun |
| Lampiran C | Seluruh RepeatGrid dan properti PageList-nya |

---

## 0. Metode — tag mana yang menjadi bukti pengikatan field

### 0.1 Berkas sampel yang dibaca penuh sebelum menetapkan tag

Empat berkas dibaca utuh untuk menentukan tag mana yang sahih; sisanya (428 berkas) diproses massal
dengan parser XML.

| Berkas | Mengapa dibaca |
| --- | --- |
| `D:\migrasi\RNM\NB FacIn\Section\DetailLocation.xml` (24 KB, terkecil) | melihat kerangka Section tanpa field sama sekali |
| `D:\migrasi\RNM\NB FacIn\Section\InputDtlOtherSchedule.xml` | sel field pertama yang lengkap: pengikatan + kontrol + kelas |
| `D:\migrasi\RNM\NB FacIn\Section\InputCoverageFire.xml` | sel numerik dengan format desimal & mata uang; sel `SUB_SECTION` |
| `D:\migrasi\RNM\NB FacIn\Section\InputObjectDtlFire.xml` (38 KB) | Section yang hanya membungkus sub-section (0 field) |

### 0.2 Tag pengikat yang dipakai — [terverifikasi]

Struktur UI Pega di ekspor ini adalah **sel tabel**: setiap elemen layar adalah satu `<rowdata>`
ber-`<pxObjClass>Embed-Display-Table-Cell</pxObjClass>`. Properti yang diikat ada di **`<pyValue>`
milik sel itu**, bukan di `pyPropertyName`, `pyReference`, atau `pyFieldName` — **ketiga tag itu
tidak ada sama sekali** di korpus Section NB.

Bukti langsung, `D:\migrasi\RNM\NB FacIn\Section\InputDtlOtherSchedule.xml`, satu sel utuh
(potongan berurutan dari berkas):

```xml
<pyType>FIELD</pyType>
<pyReadOnly>true</pyReadOnly>
<pyControlDisplayTitle>Text input</pyControlDisplayTitle>
<pyEditOptions>Read-only</pyEditOptions>
<pyRequired>false</pyRequired>
<pySmartPromptParentClass>ASM-FW-GISFW-Data-AnekaOtherSchedule</pySmartPromptParentClass>
<pyLabelPreview>Header</pyLabelPreview>
<pySmartPromptClass>ASM-FW-GISFW-Data-AnekaOtherSchedule</pySmartPromptClass>
<pyTempValue>.pyTemplateInputBox</pyTempValue>
<pxObjClass>Embed-Display-Table-Cell</pxObjClass>
<pyValue>.Header</pyValue>
<pyFormat>pxTextInput</pyFormat>
```

Tag yang dipakai dokumen ini dan artinya:

| Tag (anak langsung sel) | Arti | Label |
| --- | --- | --- |
| `<pyValue>` | **properti yang diikat** (`.Header`) — atau, pada sel `SUB_SECTION`, nama section yang disisipkan | [terverifikasi] |
| `<pyType>` | jenis sel: `FIELD` / `LABEL` / `LAYOUT` / `SUB_SECTION` / `PARAGRAPH` / `HIERARCHY` | [terverifikasi] |
| `<pyFormat>` | **jenis kontrol** (`pxTextInput`, `pxNumber`, `pxCheckbox`, `pxDateTime`, …) | [terverifikasi] |
| `<pyControlDisplayTitle>` | label kontrol versi perancang (`Text input`) | [terverifikasi] |
| `<pyReadOnly>` | read-only keras di sel | [terverifikasi] |
| `<pyEditOptions>` | `Read-only` / `Auto` / `Editable` / `Read-only-always` | [terverifikasi] |
| `<pyRequired>` | wajib isi | [terverifikasi] |
| `<pySmartPromptClass>` | **kelas Pega tempat properti itu hidup** (`ASM-FW-GISFW-Data-Coverage`) | [terverifikasi] |
| `pyUserData/<pyVisible>` + `pyUserData/<pyCondition>` | visibilitas: `ALWAYS` / `NOTBLANK` / `NOTZERO` / `OTHER`+ekspresi | [terverifikasi] |
| `pyUserData/<pyReadOnlyCondition>` | ekspresi read-only bersyarat | [terverifikasi] |
| `pyModes/rowdata/<pyDecimalPlaces>`, `<pyCurrencyType>` | pembulatan tampilan dan penanda mata uang | [terverifikasi] |
| `SectionBody/<pyRepeatDirection>` + `<pyPageListProperty>` | grid berulang dan properti daftar sumbernya | [terverifikasi] |

### 0.3 Tiga jebakan yang ditemukan dan dihindari — [terverifikasi]

**(a) `<pyValue>` tidak selalu properti bisnis.** 997 sel mengikat **placeholder template Pega**
(`.pyTemplateInputBox`, `.pyTemplateButton`, `.pyTemplateCalendar`, `.pyTemplateDisplayText`,
`.pyTemplateGeneric`). Bukti — `D:\migrasi\RNM\NB FacIn\Section\AddFacOfferList.xml`, satu sel:

```xml
<pyLabelPreview>Total Offered</pyLabelPreview>
<pySmartPromptClass>ASM-FW-GISFW-Data-FacOffer</pySmartPromptClass>
<pyValue>.pyTemplateInputBox</pyValue>
<pyFormat>pxButton</pyFormat>
```

Label manusianya "Total Offered", tetapi **tidak ada properti yang diikat** — itu tombol. Seluruh
997 sel ini **dikeluarkan** dari hitungan field. Menghitungnya akan menaikkan jumlah field 15,8 %
secara palsu.

**(b) Nama Section tidak membuktikan isi.** `InputObjectDtlFire.xml` (38 KB) **tidak mengikat satu
properti pun**; ia hanya menyisipkan `InputDtlObjFire`. 39 dari 432 Section NB seperti ini —
pembungkus murni. Daftar lengkapnya di §2.4.

**(c) `pySectionBody` sub-section tidak berisi isi section yang disisipkan.** Di ekspor, badan
section yang disisipkan selalu kosong (`<pySectionBody REPEATINGTYPE="PageList"><rowdata>
<pxObjClass>Embed-Harness-SectionBody</pxObjClass><pyTable>…` tanpa baris). Akibatnya **kedalaman
grid-dalam-grid di dalam satu berkas selalu 0** (§5.3) — penyarangan hanya terjadi lewat
penyisipan antar-Section.

### 0.4 Skrip induk — semua angka di dokumen ini berasal dari dua CSV ini

```powershell
# Skrip induk: mengekstrak setiap sel tabel dan setiap badan-grid dari 432 Section NB.
# Keluaran: nb-cells.csv (19.509 baris) dan nb-grids.csv (1.897 baris).
$src = 'D:\migrasi\RNM\NB FacIn\Section'
$cells = New-Object System.Collections.Generic.List[object]
$grids = New-Object System.Collections.Generic.List[object]
function Get-Txt($n,$t){ if($null -eq $n){return ''}; $c=$n.SelectSingleNode($t); if($null -eq $c){return ''}; $c.InnerText }
function Get-GridDepth($n){ $d=0; $p=$n.ParentNode
  while($null -ne $p){ if($p.Name -eq 'rowdata'){
      $oc=$p.SelectSingleNode('pxObjClass')
      if($null -ne $oc -and $oc.InnerText -eq 'Embed-Harness-SectionBody'){
        $rd=$p.SelectSingleNode('pyRepeatDirection')
        if($null -ne $rd -and $rd.InnerText -eq 'RepeatGrid'){$d++} } }
    $p=$p.ParentNode }
  $d }
foreach($f in (Get-ChildItem -LiteralPath $src -File -Filter *.xml | Sort-Object Name)){
  $doc=New-Object System.Xml.XmlDocument; $doc.Load($f.FullName); $sec=$f.BaseName
  foreach($n in $doc.SelectNodes("//rowdata[pxObjClass='Embed-Display-Table-Cell']")){
    $val=Get-Txt $n 'pyValue'; $typ=Get-Txt $n 'pyType'
    if($val -eq '' -and $typ -eq ''){continue}
    $ud=$n.SelectSingleNode('pyUserData'); $dec='';$cur=''
    foreach($m in $n.SelectNodes('pyModes/rowdata')){
      $dp=Get-Txt $m 'pyDecimalPlaces'; if($dp -ne ''){$dec=$dp}
      $ct=Get-Txt $m 'pyCurrencyType'; if($ct -ne ''){$cur=$ct} }
    $cells.Add([pscustomobject]@{Section=$sec;CellType=$typ;Value=$val;
      Format=(Get-Txt $n 'pyFormat');CtrlTitle=(Get-Txt $n 'pyControlDisplayTitle');
      ReadOnly=(Get-Txt $n 'pyReadOnly');EditOptions=(Get-Txt $n 'pyEditOptions');
      Required=(Get-Txt $n 'pyRequired');Visible=(Get-Txt $ud 'pyVisible');
      Cond=(Get-Txt $ud 'pyCondition');ROCond=(Get-Txt $ud 'pyReadOnlyCondition');
      SmartClass=(Get-Txt $n 'pySmartPromptClass');Decimals=$dec;CurrType=$cur;
      DateFormat=(Get-Txt $n 'pyDateFormat');GridDepth=(Get-GridDepth $n)})|Out-Null }
  foreach($b in $doc.SelectNodes("//rowdata[pxObjClass='Embed-Harness-SectionBody']")){
    $rd=Get-Txt $b 'pyRepeatDirection'; if($rd -eq ''){continue}
    $grids.Add([pscustomobject]@{Section=$sec;Direction=$rd;PageList=(Get-Txt $b 'pyPageListProperty');
      ListClass=(Get-Txt $b 'pyPageListPropertyClass');SrcType=(Get-Txt $b 'pySourceType');
      DPName=(Get-Txt $b 'pyDPName');Depth=(Get-GridDepth $b)})|Out-Null }
  $doc=$null }
$cells | Export-Csv -NoTypeInformation -Encoding UTF8 .\nb-cells.csv
$grids | Export-Csv -NoTypeInformation -Encoding UTF8 .\nb-grids.csv
"FILES=432 CELLS=$($cells.Count) GRIDROWS=$($grids.Count)"
```

Sepanjang dokumen, **`$real`** berarti himpunan field terikat nyata:

```powershell
$c = Import-Csv .\nb-cells.csv
$real = $c | Where-Object { $_.CellType -eq 'FIELD' -and
  $_.Value -match '^(\.[A-Za-z]|[A-Za-z][A-Za-z0-9_]*\.[A-Za-z])' -and
  $_.Value -notmatch '(?i)^\.pyTemplate' }
$real.Count      # 6314
```

---

## 1. Angka pokok

### 1.1 Berkas NB — 432 Section terverifikasi

```powershell
$nb='D:\migrasi\RNM\NB FacIn'
Get-ChildItem -LiteralPath $nb -Directory | ForEach-Object {
  "{0,-20} {1}" -f $_.Name, (Get-ChildItem -LiteralPath $_.FullName -File -Recurse -Filter *.xml).Count }
```

| Tipe rule (NB) | Berkas `.xml` |
| --- | ---: |
| Activity | 609 |
| **Section** | **432** |
| FlowAction | 250 |
| RDBList | 217 |
| When | 210 |
| DataTransform | 148 |
| ReportDefinition | 123 |
| Harness | 41 |
| DataPage | 30 |
| DecisionTable | 12 |
| Flow | 6 |
| ConnectREST | 3 |
| DecisionTree | 1 |
| SystemSettings | 1 |
| **Jumlah `.xml`** | **2.083** (+1 berkas `.xlsx` = 2.084 berkas) |

**432 Section terverifikasi.** Ukuran totalnya **159.743.907 byte (152,3 MB)** —
`ScoringRisk.xml` sendiri 8.048.648 byte.

```powershell
(Get-ChildItem -LiteralPath 'D:\migrasi\RNM\NB FacIn\Section' -File -Filter *.xml |
  Measure-Object Length -Sum).Sum
```

### 1.2 Sel dan field

```powershell
# sel mentah
$n=0; Get-ChildItem -LiteralPath 'D:\migrasi\RNM\NB FacIn\Section' -File -Filter *.xml |
  ForEach-Object { $n += [regex]::Matches([IO.File]::ReadAllText($_.FullName),
    '<pxObjClass>Embed-Display-Table-Cell</pxObjClass>').Count }; $n     # 21938
# jenis sel yang terekam
$c | Group-Object CellType | Sort-Object Count -Descending
```

| Metrik | Nilai |
| --- | ---: |
| Sel `Embed-Display-Table-Cell` mentah | **21.938** |
| Sel dengan `pyType` **atau** `pyValue` terisi (yang direkam) | 19.509 |
| — `LABEL` | 9.270 |
| — **`FIELD`** | **7.382** |
| — `LAYOUT` | 1.947 |
| — `SUB_SECTION` | 892 |
| — `PARAGRAPH` | 15 |
| — `HIERARCHY` | 3 |
| Sel `FIELD` yang mengikat ekspresi properti | 7.311 |
| — placeholder `.pyTemplate*` (**dibuang**) | 997 |
| — **field terikat nyata** | **6.314** |
| Ekspresi properti unik | **1.805** |
| Nama properti daun (segmen terakhir) unik | **1.096** |
| Kelas Pega berbeda (`pySmartPromptClass`) | **91** |
| Section dengan ≥1 field terikat nyata | **393** |
| Section tanpa field terikat sama sekali | **39** |

```powershell
"real=$($real.Count) unikEkspresi=$((($real.Value)|Sort-Object -Unique).Count)"
"unikDaun=$((($real|ForEach-Object{$_.Value -replace '^.*\.',''})|Sort-Object -Unique).Count)"
"kelas=$((($real|Where-Object{$_.SmartClass -ne ''}|Group-Object SmartClass)).Count)"
"sectionBerfield=$((($real|Group-Object Section)).Count)"
```

### 1.3 Lingkup properti — di dalam agregat case atau di luar

```powershell
$rel = $real | Where-Object { $_.Value -like '.*' }
$abs = $real | Where-Object { $_.Value -notlike '.*' }
"relatif=$($rel.Count) berkualifikasi=$($abs.Count)"
$abs | Group-Object { ($_.Value -split '\.')[0] } | Sort-Object Count -Descending
```

| Lingkup | Sel | Ekspresi unik | Arti untuk model Go |
| --- | ---: | ---: | --- |
| Relatif `.X` — relatif terhadap step page section | **5.990** | 1.597 | berada di bawah agregat case (`pyWorkPage`) **kecuali** bila section dirender di atas halaman lain; lihat §10-Q4 |
| `pyWorkPage.X` — eksplisit agregat case | 105 | — | pasti di agregat case |
| Halaman lain (219 sel, 35 nama halaman) | 219 | — | **di luar agregat case**: hasil pencarian, parameter, cache tampilan |

35 nama halaman non-case yang muncul, diurut jumlah sel:
`ViewOfferStatus` 19 · `InputParam` 17 · `InputBranch` 15 · `A` 14 · `InputRiskAddress` 14 ·
`InputData` 14 · `SearchAccumulation` 12 · `InputRW` 9 · `InputAccumulation` 9 ·
`AccumulationRisk` 7 · `TempPolis` 7 · `InputCity` 7 · `InputMarketingOfficer` 7 ·
`FacOfferList` 7 · `AdjustmentRisk` 7 · `InputAccumulatedType` 6 · `RISlipContent` 6 ·
`SearchSOB` 5 · `SFAResults` 4 · `ProtectSpreading` 3 · `InputProvince` 3 · `InputDistrict` 3 ·
`OutputParam` 3 · `CountAccumulation` 3 · `InputNation` 3 · `CARI` 3 · `SearchLife` 2 ·
`TempError` 2 · `Clause` 2 · `pyPortal` 1 · `InputAccumulationCode` 1 · `TempIndex` 1 ·
`InputBranchDetail` 1 · `InputPar` 1 · `test` 1.

Halaman `test` **[pertanyaan terbuka]** — sisa uji coba yang tertinggal di produksi, atau dipakai?

### 1.4 Keterjangkauan Section dari layar

```powershell
$nb='D:\migrasi\RNM\NB FacIn'; $fa=@{}
Get-ChildItem -LiteralPath "$nb\FlowAction" -File -Filter *.xml | ForEach-Object {
  $v=[regex]::Match([IO.File]::ReadAllText($_.FullName),'<pySectionReference>([^<]*)</pySectionReference>').Groups[1].Value
  if($v){$fa[$v]=1+$fa[$v]} }
$hs=@{}; Get-ChildItem -LiteralPath "$nb\Harness" -File -Filter *.xml | ForEach-Object {
  foreach($m in [regex]::Matches([IO.File]::ReadAllText($_.FullName),'<pyStreamName>([^<]+)</pyStreamName>')){$hs[$m.Groups[1].Value]=1} }
$secNames=@{}; (Get-ChildItem -LiteralPath "$nb\Section" -File -Filter *.xml).BaseName|ForEach-Object{$secNames[$_]=$true}
(@($fa.Keys+$hs.Keys)|Sort-Object -Unique|Where-Object{$secNames.ContainsKey($_)}).Count
```

| Metrik | Nilai |
| --- | ---: |
| FlowAction NB dengan `pySectionReference` terisi | 248 dari 250 |
| Section berbeda yang dirujuk FlowAction | 229 |
| Nama stream unik pada 41 Harness NB | 88 |
| **Section yang dirujuk langsung Harness/FlowAction** | **272 dari 432** |

160 Section sisanya hanya dicapai lewat penyisipan dari Section lain (`SUB_SECTION`). Untuk React
ini berarti: **272 kandidat "layar/dialog", 160 kandidat "komponen"** — [dugaan], karena rujukan
langsung bukan bukti bahwa sesuatu dirender sebagai halaman.

---

## 2. Skema field per Section

**Tabel penuh ada di Lampiran A** — 432 Section, masing-masing dengan tabel
`Properti | Kontrol | n | Read-only | Kondisi tampil`. Bagian ini hanya memberi peta.

### 2.1 Distribusi ukuran layar

```powershell
$sum = Import-Csv .\nb-section-summary.csv   # dihasilkan §4.2
$sum | Group-Object { switch([int]$_.Fields){ {$_ -eq 0}{'0'} {$_ -le 5}{'1-5'}
  {$_ -le 20}{'6-20'} {$_ -le 50}{'21-50'} default{'>50'} } }
```

| Field terikat | Section |
| --- | ---: |
| 0 | 39 |
| 1–5 | 146 |
| 6–20 | 176 |
| 21–50 | 46 |
| > 50 | 25 |
| **Jumlah** | **432** |

### 2.2 30 Section terbesar — [terverifikasi]

| Section | Field | Unik | Read-only | RepeatGrid | Rumpun |
| --- | ---: | ---: | ---: | ---: | --- |
| `ScoringRisk` | 438 | 436 | 233 | 0 | FIRE |
| `ScoringRiskForm2` | 158 | 158 | 78 | 0 | FIRE |
| `ViewDtlOtherObjectAneka` | 119 | 63 | 118 | 0 | ANEKA |
| `InputDtlOtherObjectAneka` | 119 | 63 | 119 | 0 | ANEKA |
| `InputDtlOtherObjectAneka_FacIn` | 119 | 65 | 119 | 0 | INTI+ANEKA |
| `InputDtlOtherObjectAneka_FacIn_ISUW` | 119 | 65 | 119 | 0 | INTI+ANEKA |
| `InputDtlObject_FacIn` | 103 | 56 | 59 | 22 | PA-LIFE |
| `InputDtlObject_FacIn_IsUW` | 100 | 55 | 59 | 22 | PA-LIFE |
| `ScoringRiskForm1` | 96 | 96 | 48 | 0 | FIRE |
| `ScoringRiskForm4` | 96 | 96 | 48 | 0 | FIRE |
| `GeneralPolicyTreatyIn` | 87 | 76 | 49 | 3 | INTI |
| `GeneralDeptHeadTreatyIn_UW` | 87 | 76 | **84** | 3 | INTI |
| `DetailDeptHeadTreatyIn_UW` | 87 | 76 | **84** | 3 | INTI |
| `DetailPolicyTreatyIn` | 87 | 76 | 49 | 3 | INTI |
| `InputDtlObjectAneka_FacIn` | 83 | 53 | 83 | 0 | ANEKA |
| `InputObjectDtlFire_GCNM` | 80 | 75 | 70 | 6 | FIRE |
| `InputDtlObjectAneka_FacIn_IsUW` | 73 | 48 | 73 | 0 | ANEKA |
| `InputDtlObjFire` | 72 | 72 | 67 | 4 | FIRE |
| `InputDtlObjFire_IsUW` | 72 | 72 | 71 | 4 | FIRE |
| `ViewInwardFacultativeDtl_IsUW` | 69 | 30 | 57 | 12 | INTI |
| `DetailPolicyTreatyInNonProportional` | 64 | 23 | 18 | 17 | INTI |
| `InputInwardFacultativeDtl` | 63 | 36 | 37 | 10 | INTI |
| `ScoringRiskForm3` | 58 | 58 | 30 | 0 | FIRE |
| `ViewDtlObjectAneka` | 52 | 37 | 52 | 0 | ANEKA |
| `ViewCSVMarine` | 52 | 52 | 52 | 1 | INTI |
| `CoverageItem` | 48 | 42 | 44 | 3 | INTI |
| `InputEndorsementDtl` | 46 | 30 | 28 | 7 | INTI |
| `InputInwardFacultativeDtl_IsUW` | 46 | 31 | 36 | 8 | INTI |
| `PropertyItemListCoverageSpreading` | 46 | 17 | 40 | 12 | FIRE |
| `PersonCoverage` | 45 | 29 | 12 | 5 | PA-LIFE |

**Temuan #1 — `ScoringRisk` adalah anomali.** 438 field terikat dengan 436 ekspresi berbeda: hampir
tidak ada pengulangan, 22 nama properti daun saja (`ChechBox1`…`ChechBox10`, `Score1`…`Score7`,
`Score`, `PercentLossRatio`, `InsuredName`, `StartDateTime`, `EndDateTime`) yang dipakai lewat
**indeks larik keras di dalam ekspresi**, misalnya
`.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox1`. Itu berarti jumlah baris skoring
**dikodekan di layar**, bukan data. Salah eja `ChechBox` dipertahankan apa adanya. [terverifikasi]

### 2.3 Duplikasi `_IsUW` — [terverifikasi]

```powershell
$names=(Get-ChildItem -LiteralPath 'D:\migrasi\RNM\NB FacIn\Section' -File -Filter *.xml).BaseName
$uw=$names|Where-Object{$_ -match '(?i)(_IsUW|_ISUW|_UW)$'}
"uw=$($uw.Count) berpasangan=$((($uw|Where-Object{$names -contains ($_ -replace '(?i)(_IsUW|_ISUW|_UW)$','')})).Count)"
```

**87** Section NB bernama `*_IsUW` / `*_ISUW` / `*_UW`; **81** di antaranya punya kembaran non-UW.
Contoh pasangan dengan skema field praktis identik: `InputDtlObjFire` (72 field) vs
`InputDtlObjFire_IsUW` (72 field) — bedanya read-only 67 vs 71.

**Konsekuensi migrasi:** 81 pasang layar kembar = **satu komponen React dengan prop `mode`**, bukan
dua komponen. Tetapi **perbedaan read-only-nya harus dipetakan satu per satu** dari Lampiran A —
kembar bukan berarti sama.

### 2.4 39 Section pembungkus (0 field terikat) — [terverifikasi]

`additemButton`, `AllSummarySection`, `CallDualScoringRisk`, `ChooseAccumulation_FacIn`,
`ComfirmPolis_Life`, `ComfirmPolis_NonLife`, `ConfirmChangeOutgo`, `CreateListClause`,
`deleteItemButton`, `DetailLocation`, `DetailPoliciesNonProportional`, `FacOutOffer`,
`FinalPolicyRO`, `InputAnekaParticipant`, `InputObjectDtlFire`, `InputObjectDtlFire_IsUW`,
`InputOtherObjectAneka_FacIn`, `InputOtherObjectAneka_FacIn_ISUW`, `InputOtherObjectAneka_UW`,
`OfferFacinLifeNotification`, `PolicyTreatyInDeclineConfirm`, `Property`, `Property_IsUW`,
`ProtectCurrency`, `ProtectCurrency_IsUW`, `ProtectCurrencyCargo`, `ProtectCurrencyCargo_IsUW`,
`pyAttachmentScreen`, `RejectNotificationSection`, `ReviseNotification`, `SFAPortal_Opportunities`,
`ViewCoverage`, `ViewCoveragePA`, `ViewCoverageTravel`, `ViewDeliveryAddressFacOut`,
`ViewDtlDeductible`, `ViewDtlOutgo`, `ViewObjectSavior`, `ViewPerilsFacOut`.

Empat di antaranya berukuran besar (`FinalPolicyRO` 46 KB, `ComfirmPolis_Life` 50 KB,
`ConfirmChangeOutgo` 49 KB, `DetailPoliciesNonProportional`) — **[pertanyaan terbuka]** apakah
isinya benar-benar kosong atau seluruh isinya berupa teks statis/HTML non-autogenerated yang tidak
terbaca sebagai sel.

---

## 3. Katalog properti unik — tulang punggung model data

**Lampiran B memuat seluruh 1.805 ekspresi**, diurutkan menurun, dengan kolom
`n | Section | Kelas Pega | Kontrol | Read-only | Lingkup`.

```powershell
$real | Group-Object Value | Sort-Object Count -Descending | Select-Object -First 40
```

### 3.1 40 properti paling sering diikat — [terverifikasi]

| n | Section | Ekspresi | Kelas Pega dominan | Lingkup |
| ---: | ---: | --- | --- | --- |
| 107 | 81 | `.Premium` | `…Data-Coverage` | agregat |
| 103 | 66 | `.TSI` | `…Data-OfferFacIn-Currency` | agregat |
| 92 | 47 | `.Name` | `…Data-OfferFacIn-Currency` | agregat |
| 85 | 48 | `.Currency` | `…Data-PropertyItem` | agregat |
| 76 | 39 | `.ObjectName` | `…Data-Aneka` | agregat |
| 60 | 25 | `.pyFullName` | `Data-Party-Person` | agregat |
| 57 | 45 | `.CoverageNote` | `…Data-Coverage` | agregat |
| 52 | 41 | `.Rate` | `…Data-Coverage` | agregat |
| 51 | 24 | `.ASMDateOfBirth` | `Data-Party-Person` | agregat |
| 43 | 17 | `.TreatyType` | `…Data-SpreadingRisk` | agregat |
| 42 | 21 | `.ASMGender` | `Data-Party-Person` | agregat |
| 40 | 34 | `.OccupationName` | `…Data-Occupation` | agregat |
| 36 | 17 | `.SharePercentage` | `…Data-SpreadingRisk` | agregat |
| 34 | 17 | `.TSISpreaded` | `…Data-SpreadingRisk` | agregat |
| 32 | 8 | `.Value` | `…Data-TreatyInTotal` | agregat |
| 32 | 21 | `.Currency.Name` | `…Data-Aneka` | agregat |
| 30 | 21 | `.Property.RiskLocation.ASMAddress` | `…Data-OfferFacIn-LocationReinsurance` | agregat |
| 30 | 14 | `.PremiumSpreaded` | `…Data-SpreadingRisk` | agregat |
| 29 | 18 | `.ItemType` | `…Data-PropertyItem` | agregat |
| 25 | 15 | `.TSIObjectItem` | `…Data-PropertyItem` | agregat |
| 25 | 19 | `.BrandName` | `…Data-Vehicle` | agregat |
| 24 | 20 | `.PremiNusantaraRe` | `…Data-Coverage` | agregat |
| 23 | 17 | `.ASMLeftHanded` | `Data-Party-Person` | agregat |
| 23 | 23 | `.OccupationID` | `…Data-Occupation` | agregat |
| 23 | 8 | `.Year` | `…Data-Aneka` | agregat |
| 22 | 18 | `.TypeName` | `…Data-Vehicle` | agregat |
| 22 | 17 | `.ASMJobName` | `Data-Party-Person` | agregat |
| 20 | 20 | `.TSINusantaraRe` | `…Data-Coverage` | agregat |
| 20 | 20 | `.ID` | `…Int-SUBCONTRACT` | agregat |
| 19 | 12 | `.SumTotalPayment` | `…Data-OfferFacIn-Currency` | agregat |
| 19 | 15 | `.Descriptions` | `…Data-Deductible` | agregat |
| 18 | 4 | `.Aneka.Obligee.UserAssignment` | `ASM-FW-GISFW-Work` | agregat |
| 18 | 4 | `.Aneka.Obligee.PublishedIn` | `ASM-FW-GISFW-Work` | agregat |
| 18 | 16 | `.ConveyanceNote` | `…Data-Cargo` | agregat |
| 18 | 13 | `.ASMRelation` | `Data-Party-Person` | agregat |
| 18 | 11 | `.ClaimSpreaded` | `…Data-SpreadingRisk` | agregat |
| 18 | 13 | `.Property.ObjectNo` | `…Data-OfferFacIn-LocationReinsurance` | agregat |
| 17 | 17 | `.ASMIDCard` | `Data-Party-Person` | agregat |
| 17 | 17 | `.LicensePlate` | `…Data-Vehicle` | agregat |
| 17 | 17 | `.ModelName` | `…Data-Vehicle` | agregat |

Perhatikan: **satu nama daun dipakai oleh banyak kelas.** `.Premium` muncul di 81 Section dan
diatribusikan ke sekurangnya 4 kelas berbeda. Di Go, `.Premium` **bukan satu field** — ia field
pada `Coverage`, pada `PropertyItem`, pada `PolicyPayment`, dan pada `CurrencyLine`. Menyatukannya
akan menghasilkan struct yang salah.

### 3.2 Bentuk agregat case yang terbaca dari jalur properti — [terverifikasi]

```powershell
$deep = $real | Where-Object { $_.Value -match '^\.[A-Za-z0-9_]+\.' }
$deep | Group-Object { ($_.Value -split '\.')[1] } | Sort-Object Count -Descending
```

**1.426** field terikat memakai jalur berjenjang. Segmen pertama, diurut jumlah:

| Segmen | Sel | Segmen | Sel |
| --- | ---: | --- | ---: |
| `.OfferFacIn.` | 323 | `.ScoringRisk.` | 31 |
| `.Aneka.` | 280 | `.VehicleHE.` | 31 |
| `.Property.` | 184 | `.MarineHull.` | 27 |
| `.PolicyData.` | 120 | `.RiskLocation.` | 24 |
| `.Policy.` | 109 | `.BuildingConstruction.` | 21 |
| `.QuotationData.` | 78 | `.AviationHull.` | 18 |
| `.SurroundingRisk.` | 60 | `.Customer_P.` | 17 |
| `.Currency.` | 35 | `.AnekaPolicySchedule.` 9 · `.AnekaMaintenance.` 9 · `.TableOfLimit.` 6 · `.CoinsData.` 6 · `.Occupation.` 6 · `.DataFEA.` 6 · `.Customer_C.` 4 · `.ASMBankData.` 4 · `.InwardScale.` 4 · `.AnekaAnotherSchedule.` 4 | |

Ini adalah **peta struktur agregat case yang terbukti dari UI** — bukan dari nama rule. Tiap segmen
adalah kandidat struct bersarang di Go.

---

## 4. Rumpun lini bisnis — dari properti terikat, bukan nama berkas

### 4.1 Pembeda lini bisnis yang sesungguhnya — [terverifikasi]

Lini bisnis **tidak** ditentukan nama Section. Ia ditentukan satu properti, terbaca dari rule
`When` yang menggerbangi visibilitas field:

```powershell
foreach($n in @('IsFire','IsMBU','IsMarineCargo','IsPA','IsTravel','IsAneka','IsHE','IsLife')){
  $c=[IO.File]::ReadAllText("D:\migrasi\RNM\NB FacIn\When\$n.xml")
  "{0,-16} -> {1}" -f $n, (([regex]::Matches($c,'<pyUnmodifiedPath>([^<]*)</pyUnmodifiedPath>')|
     ForEach-Object{$_.Groups[1].Value}|Sort-Object -Unique) -join ' ; ') }
```

| Rule `When` (NB) | Properti yang dibandingkan |
| --- | --- |
| `IsFire`, `IsMBU` | `pyWorkPage.OfferFacIn.QuotationData.BusinessType` |
| `IsMarineCargo` | `.OfferFacIn.QuotationData.BusinessType` **dan** `pyWorkPage.Quotation.BusinessType` |
| `IsPA`, `IsTravel`, `IsAneka`, `IsHE`, `IsMarineHull`, `IsGolfInsurance`, `IsMBUCar` | `pyWorkPage.Quotation.BusinessType` |
| `IsEngineering`, `IsLiability` | `pyWorkPage.Quotation.BusinessCode` **dan** `…BusinessType` |
| `IsCreditBriguna` | `pyWorkPage.Quotation.BusinessCode` |
| `IsLife` | `pyWorkPage.Quotation.BusinessOldId` |
| `IsCustomBonds` | `…OfferFacIn.QuotationData.BusinessType`, `…Policy.Quotation.BusinessType`, `pyWorkPage.Quotation.BusinessType` |

**Temuan #2 — [terverifikasi]: dua jalur berbeda untuk satu pembeda.** `IsMarineCargo` dan
`IsCustomBonds` membaca `BusinessType` dari **dua/tiga lokasi berbeda dalam agregat**
(`.Quotation` vs `.OfferFacIn.QuotationData` vs `.Policy.Quotation`). Sebelum model Go dikunci,
bisnis harus menjawab apakah ketiganya selalu bernilai sama. Bila tidak, satu rule dapat memberi
hasil berbeda tergantung jalur. **MEMBLOKIR** untuk desain `Quotation` (§10-Q1).

Nilai-nilai `BusinessType`, `BusinessCode`, `BusinessOldId` **tidak terbaca** dari ekspor Section
(isi kondisi `When` ada di rule `When`, di luar cakupan dokumen ini) — **belum terverifikasi**.

### 4.2 Pengelompokan Section — dua sinyal, bukan nama berkas

Aturan yang dipakai (dinyatakan agar dapat dibantah):

1. **Sinyal kelas** — `pySmartPromptClass` sel field dipetakan ke rumpun.
2. **Sinyal jalur** — segmen jalur properti (`.Property.`, `.Aneka.`, `.MarineHull.`, …) dipetakan
   ke rumpun.
3. Rumpun dengan skor tertinggi menang. Bila skor tertinggi × 3 < jumlah field Section itu, label
   diberi awalan `INTI+` — **ambang 3 ini pilihan pelaporan, bukan temuan korpus.**
4. Section tanpa sinyal apa pun = `INTI` (lintas lini).

```powershell
# skrip lengkap ada di §0.4 + peta kelas/jalur di bawah; keluarannya nb-section-summary.csv
$rows | Group-Object Lob | Sort-Object Count -Descending
```

| Rumpun | Section | Sel field (sinyal kelas saja) |
| --- | ---: | ---: |
| **INTI** (lintas lini) | 230 | 4.062 |
| **FIRE / PROPERTY** | 67 | 1.085 |
| *(tanpa field)* | 39 | — |
| **PA / LIFE** | 33 | 411 |
| **ANEKA** (bonding, kredit, golf, travel-objek, alat berat) | 25 | 330 |
| **MARINE CARGO** | 17 | 282 |
| `INTI+FIRE` | 11 | — |
| **MBU / KENDARAAN** | 6 | 144 |
| `INTI+ANEKA` | 2 | — |
| `INTI+PA-LIFE` | 1 | — |
| `INTI+MBU` | 1 | — |

Kelas Pega yang dipetakan ke tiap rumpun (dasar sinyal 1):

| Rumpun | Kelas `pySmartPromptClass` |
| --- | --- |
| FIRE | `…Data-Property`, `…Data-PropertyItem`, `…Data-Occupation`, `…Int-OCCUPATION`, `…Int-RISKADDRESS`, `…Data-FEA`, `…Data-ScoringRisk` |
| CARGO | `…Data-Cargo`, `…Data-Conveyance`, `…Data-Good`, `…Data-Packing`, `…Data-Trading`, `…Int-SHIP`, `…Int-MSHIP` |
| MBU | `…Data-Vehicle`, `…Data-Accessory` |
| PA-LIFE | `Data-Party-Person`, `…Data-Benefit`, `…Data-Plan`, `…Data-Plan-PremiumList`, `…Int-RETROCESSIONLIFE` |
| ANEKA | `…Data-Aneka`, `…Data-AnekaParticipant`, `…Data-AnekaOtherSchedule` |

### 4.3 Daftar Section per rumpun (selain INTI) — [terverifikasi]

**FIRE / PROPERTY (67)** — angka dalam kurung = jumlah field terikat:
`ScoringRisk`(438), `ScoringRiskForm2`(158), `ScoringRiskForm1`(96), `ScoringRiskForm4`(96),
`InputObjectDtlFire_GCNM`(80), `InputDtlObjFire_IsUW`(72), `InputDtlObjFire`(72),
`ScoringRiskForm3`(58), `PropertyItemListCoverageSpreading`(46),
`PropertyItemListCoverageSpreading_IsUW`(26), `RiskAround_IsUW`(24), `RiskAround`(24),
`InputDtlObjectLocation_FacIn`(23), `ObjectDetails_IsUW`(23), `ObjectDetails`(23),
`PropertyItemListCoverage`(23), `PropertyItemListCoverage_IsUW`(21), `PropertyItemList`(18),
`PropertyItemList_IsUW`(18), `PropertyItemFacIn_Section`(16),
`PropertyItemListCoverageCommision`(15), `InputDtlObjectLocation_FacIn_IsUW`(15),
`InputObjItem`(12), `PropertyItemFacIn_Section_IsUW`(12), `ShowCoverageFacOut`(12),
`ViewCheckListOfferFacOut`(12), `InputFEA`(11), `LocationDetail`(9), `PropertyItemListFacOut`(9),
`InputFEA_IsUW`(9), `SummaryCoverage_Section`(8), `InputDtlObjectLocation_GISFW`(8),
`ChooseRiskAddress_ResultList`(8), `OccupationFacOut`(7), `OccupationFacOut_IsUW`(7),
`ChooseRiskAddress`(7), `AnekaListFacOut_IsUW`(6), `AnekaListFacOut`(6),
`PropertyItemListFacOut_IsUW`(6), `ViewObjectFireFacOut`(5), `InputOkupasi`(5),
`InputCoverageFire`(5), `InputOkupasiHull_GCNM`(5), `InputCoverageFire_IsUW`(5),
`ChooseZipCodeDtl`(4), `ObjectList`(4), `ObjectList_IsUW`(4),
`InputObjectSubContract_FacIn_IsUW`(3), `CoverageList`(3), `ObjectItemSummary`(3),
`CoverageList_IsUW`(3), `ChooseOccupation_Dtl`(3), `ChooseOccupation`(3), `OccupationList`(3),
`OccupationList_IsUW`(3), `InputObjectSubContract_FacIn`(3), `OccupationDetail`(3),
`OccupationItemFacIn_Section`(3), `OccupationItemFacIn_Section_IsUW`(3), `FireSummarySection`(2),
`GrowingTreesOccupationSummary`(2), `CopyCoverageFrom`(2), `ViewObjectOccupation`(1),
`OccupationPropertyFacOut`(1), `LocationPropertyFacOut`(1), `GolfSummarySection`(1),
`GrowingTreesSummarySection`(1).

**PA / LIFE (33):** `InputDtlObject_FacIn`(103), `InputDtlObject_FacIn_IsUW`(100),
`PersonCoverage`(45), `ShowObjectFacOut`(38), `CoverageCommisionList`(28),
`CoverageSpreadingList`(28), `CoverageSpreadingList_IsUW`(28), `ShowObjectFacOut_IsUW`(28),
`ViewDtlCoverageFacOut`(24), `CoverageCommisionList_IsUW`(24), `InputDtlParticipantDM`(18),
`ShowCoveragePAFacOut_Dtl`(18), `InputDtlParticipantLife_FacIn`(17),
`InputDtlParticipantPA_FacIn`(17), `ViewDtlParticipantPA`(16),
`InputDtlParticipantPA_FacIn_IsUW`(16), `InputDtlParticipantLife_FacIn_IsUW`(15),
`InputDtlParticipantTravel_IsUW`(12), `InputDtlParticipantTravel`(12),
`InputDtlParticipantTravel_GISFW`(12), `InputPerson`(10), `ShowCoveragePAFacOutIsUW_Dtl`(9),
`ShowCoverageTravelFacOut_Dtl_IsUW`(9), `ShowCoverageTravelFacOut_Dtl`(9),
`InputDtlLayerSpreading_FacIn`(8), `PersonListDetail`(8), `PlanContains_UWRO`(8),
`ViewObjectPAFacOut`(5), `ViewObjectTravelFacOut`(5), `ViewPolicyListRO`(5),
`SelectPlanListRO`(2), `confirmCreateAccumulation`(1), `PASummarySection`(1).

**Travel tidak punya rumpun sendiri** — sepenuhnya memakai kelas `Data-Party-Person`
(`InputDtlParticipantTravel*`, `ShowCoverageTravelFacOut_Dtl*`). [terverifikasi]

**ANEKA (25):** `InputDtlOtherObjectAneka`(119), `ViewDtlOtherObjectAneka`(119),
`InputDtlObjectAneka_FacIn`(83), `InputDtlObjectAneka_FacIn_IsUW`(73), `ViewDtlObjectAneka`(52),
`ObjectOccupation_FacIn`(20), `InputOkupasiAneka_FacIn`(19), `InputOkupasiAneka_FacIn_IsUW`(19),
`ViewOkupasiAneka`(16), `InputDtlAnekaPolicySchedule`(9), `ObjectDtlAneka_FacIn`(9),
`InputDtlMaintenance`(9), `InputOkupasiAneka_GCNM`(7), `GrowingTreesLocationDetail`(7),
`InputDtlParticipant`(6), `InputCoverageDeductibleHull_GCNM`(5),
`InputCoverageDeductibleAneka_GCNM`(5), `GolfLocationDetail`(5), `CommisionObjectAneka`(4),
`SpreadingObjectAneka`(4), `InputOtherSchedule`(4), `SpreadingObjectAneka_IsUW`(3),
`GolfObjectItemSummary`(3), `InputDtlOtherSchedule`(2), `ViewObjectDtlAneka`(1).
Ditambah `INTI+ANEKA`: `InputDtlOtherObjectAneka_FacIn`(119),
`InputDtlOtherObjectAneka_FacIn_ISUW`(119).

**MARINE CARGO (17):** `InputDtlCargo_FacIn_IsUW`(38), `InputDtlCargo_FacIn`(37),
`MarineCargoDtl_IsUW`(33), `MarineCargoDtl`(32), `SelectShip`(16), `ShowCoverageFacOut_IsUW`(12),
`InputDtlTrading`(10), `InputDtlTrading_IsUW`(9), `ChooseObject_Ship`(5), `InputDtlGoods`(5),
`InputDtlCargo`(4), `ViewDetailCargo`(4), `ViewObjectMarineCargoFacOut`(4),
`InputDtlGoods_IsUW`(4), `SelectCoverage`(3), `InputDtlConveyance`(3),
`InputDtlConveyance_IsUW`(3).

**MBU / KENDARAAN (6):** `ViewVehicleGridFacOut`(17), `VehicleGrid_FacIn_IsUW`(17),
`VehicleGrid_FacIn`(17), `VehicleGrid`(17), `OfferStatusFacOut`(8),
`ViewObjectVehicleFacOut`(6). Ditambah `INTI+MBU`: `OldOfferStatusFacOut`(20).

`INTI+FIRE` (11): `Periode`(44), `Periode_IsUW`(37), `ViewCheckListOffer`(33),
`SummaryLossRecord_Section`(10), `UploadMember`(10), `ViewDtlObjectLocation`(7),
`InputDtlObjectLocationHull_GISFW`(7), `FEAList`(6), `FEAList_IsUW`(6),
`ViewCoverageFacOutShow`(5), `ViewCoverageFacOutShow_isUW`(5).
`INTI+PA-LIFE` (1): `ViewTabGroupPolisFacOutSavior`(32).

### 4.4 Field bersama vs field khas — jawaban atas "satu struct atau banyak struct"

```powershell
$leafLob=@{}
foreach($r in $real){ $l = <rumpun dari kelas>; $k=($r.Value -replace '^.*\.','')
  if(-not $leafLob.ContainsKey($k)){$leafLob[$k]=@{}}; $leafLob[$k][$l]=$true }
"total=$($leafLob.Count) bersama=$((@($leafLob.Keys|Where-Object{$leafLob[$_].Keys.Count -gt 1})).Count)"
```

| Kategori | Nama properti daun |
| --- | ---: |
| Total nama daun unik | **1.096** |
| Muncul di **> 1** rumpun (**BERSAMA**) | **133** |
| Muncul di tepat **1** rumpun (**KHAS**) | **963** |
| — khas INTI | 790 |
| — khas MARINE-CARGO | 51 |
| — khas FIRE/PROPERTY | 45 |
| — khas ANEKA | 43 |
| — khas PA/LIFE | 25 |
| — khas MBU/KENDARAAN | 9 |

26 nama daun yang muncul di **≥ 3 rumpun** — kandidat terkuat untuk field bersama:
`BrandName`, `Commision`, `Construction`, `Description`, `EngineNumber`, `HASIL12`, `ID`,
`IsTopRisk`, `Name`, `ObjectName`, `ObjectPremi`, `OccupationName`, `PercentOffered`, `Premium`,
`pyCity`, `pyFullName`, `pyName`, `pySelected`, `Rate`, `Remark`, `RIComm`, `ShareOffered`,
`TSI`, `TSISpreaded`, `TypeName`, `Year`.

**Field khas per lini, dari properti terikat — [terverifikasi]:**

| Rumpun | Contoh field yang hanya muncul di rumpun ini |
| --- | --- |
| FIRE | `EML`, `PML`, `FloodArea`, `FloodAreaStatus`, `BackConstruction`/`FrontConstruction`/`LeftConstruction`/`RightConstruction`, `BackDistance`/`FrontDistance`/…, `RoofType`, `WallType`, `FloorType`, `PartitionType`, `SupportWallType`, `NumberOfFloor`, `HousekeepingStatus`, `IsHotWorkProcessFlag`, `IsFlammableItemFlag`, `IsProductionProcessFlag`, `FirstLossLimitFac`, `RiskCategory`, `AlmRiskID`, `FEAInside`/`FEAOutside`, `PrivateFireBrigade`, `TeamSOPSafety` |
| MARINE CARGO | `BLNumber`, `SailDate`, `FromRute`, `ToRute`, `GRT`, `NRT`, `DWT`, `SHIPTYPE_NAME`, `NM_SHIP`, `Y_MAKE1`, `MaxAge`, `MaxAgeDetail`, `LCNumber`, `LCBANK`, `LCEndDate`, `LCCondition`, `InvoiceNumber`, `InvoiceDate`, `PackingID`, `GoodID`, `ConveyanceID`, `TradingID` |
| MBU | `LicensePlate`, `CC`, `ColorName`, `TransmissionType`, `ManufactureYear`, `ModelName`, `CarParking` — perhatikan `ChassisNumber` dan `EngineNumber` **bukan** khas MBU: masing-masing 11 sel di `Data-Vehicle` **dan** 11 sel di `Data-Aneka` |
| PA / LIFE | `ASMDateOfBirth`, `ASMGender`, `ASMLeftHanded`, `ASMJobName`, `ASMIDCard`, `ASMRelation`, `ASMWeight`, `ASMHeight`, `Age` |
| ANEKA | `Obligee.*` (`ObligeeName`, `ProjectName`, `ContractNo`, `ContractDate`, `AuctionDate`, `PublishedIn`, `UserAssignment`, `Jabatan`, `PaymentClaim`), `LoanAccountNumber`, `LoanName`, `LoanType`, `DepositAccountNo`, `LimitOfLiability`, `GeographicalLimits`, `Pilots`, `RegistrationMarks`, `TypeOfVessel`, `LENGTHM`/`BREADTHM`/`DEPTHM`, `SpecialUses`/`StandardUses`/`SpecialRentalUses` |

**Kesimpulan untuk model Go — [usulan], lihat §9.2:** rasio 963 khas : 133 bersama **menolak**
satu struct tunggal berisi semua field. Yang didukung bukti adalah **satu inti + lima bagian
opsional**.

---

## 5. Repeating grid

### 5.1 Jumlah — [terverifikasi]

```powershell
$g = Import-Csv .\nb-grids.csv
$g | Group-Object Direction | Sort-Object Count -Descending
$rg = $g | Where-Object { $_.Direction -eq 'RepeatGrid' }
"RepeatGrid=$($rg.Count) sectionUnik=$((($rg|Group-Object Section)).Count)"
"pagelistUnik=$((($rg|Where-Object{$_.PageList -ne ''}).PageList|Sort-Object -Unique).Count)"
```

| `pyRepeatDirection` | Badan section |
| --- | ---: |
| `Bottom` | 632 |
| `Top` | 632 |
| **`RepeatGrid`** | **629** |
| `TreeGrid` | 2 |
| `HORIZONTAL` | 1 |
| `Tree` | 1 |

| Metrik grid NB | Nilai |
| --- | ---: |
| Layout `RepeatGrid` | **629** |
| Section yang memuatnya | **281 dari 432 (65,0 %)** |
| Properti `pyPageListProperty` unik | **181** |
| `pySourceType` = `Property` | **629 / 629** (tidak ada grid bersumber Report Definition/Data Page) |
| `pyExpandable=true` | 424 kemunculan |
| `pyLoadDeferred=true` | 22 kemunculan |

**Seluruh grid bersumber properti PageList** — tidak satu pun dari Report Definition. Untuk Go ini
berarti: **daftar dikirim bersama agregat case**, bukan diambil lewat endpoint daftar terpisah.
[terverifikasi]

### 5.2 30 properti daftar paling sering menjadi grid — [terverifikasi]

| n | Properti PageList | n | Properti PageList |
| ---: | --- | ---: | --- |
| 40 | `.CoverageList` | 10 | `.CurrencyList` |
| 32 | `.PersonList` | 10 | `.Property.RiskLocation.AnekaList` |
| 21 | `.LocationList` | 10 | `.ASMHeir` |
| 20 | `.Property.PropertyItemList` | 9 | `.CargoList` |
| 20 | `.AnekaList` | 9 | `.OfferFacIn.FacRetroList` |
| 17 | `.DeductibleList` | 8 | `SpreadingList.pxResults` |
| 16 | `.VehicleList` | 8 | `.OccupationList` |
| 15 | `pyWorkPage.OfferFacIn.CurrencyList` | 7 | `.SpreadingList` |
| 15 | `TotalSpreadAll.pxResults` | 6 | `TempViewSuggest.pxResults` |
| 14 | `.PropertyList` | 6 | `.OfferFacIn.PersonList` |
| 14 | `.Property.RiskLocation.OccupationList` | 6 | `.OfferFacIn.LocationList` |
| 13 | `.ASMCoverage` | 6 | `.Property.TotalTSIPremiGrossList` |
| 12 | `.ClauseList` | 5 | `.SpreadingRiskList` |
| 12 | `.OfferFacIn.CurrencyList` | 5 | `.SubContractList` |
| 11 | `TotalSpreadingCurrency.pxResults` | 5 | `.FEAList` |

Perhatikan **`.CurrencyList` diikat dengan tiga ejaan jalur berbeda** —
`.CurrencyList`, `.OfferFacIn.CurrencyList`, `pyWorkPage.OfferFacIn.CurrencyList` — total 37 grid.
Sama halnya `.PersonList` vs `.OfferFacIn.PersonList`, dan `.LocationList` vs
`.OfferFacIn.LocationList`. **[dugaan]** ketiganya menunjuk daftar yang sama, dilihat dari step page
berbeda; **belum terverifikasi** dari korpus.

Grid bersumber `*.pxResults` (`TotalSpreadAll`, `TotalSpreadingCurrency`, `SpreadingList`,
`TempViewSuggest`) **bukan bagian agregat case** — itu halaman hasil query yang diisi Activity/
RDBList. Di Go ini menjadi **endpoint baca terpisah**, bukan field pada model case. [terverifikasi]

### 5.3 Kedalaman penyarangan — [terverifikasi]

```powershell
$rg | Group-Object Depth   # depth = jumlah induk RepeatGrid dalam berkas yang sama
```

**Kedalaman struktural di dalam satu berkas = 0 untuk seluruh 629 grid.** Sebabnya §0.3(c):
ekspor tidak memuat isi section yang disisipkan. Penyarangan karena itu **hanya dapat diukur
lewat graf penyisipan**:

```powershell
$sub = $c | Where-Object { $_.CellType -eq 'SUB_SECTION' -and $_.Value -ne '' }
$nbNames=@{}; (Get-ChildItem -LiteralPath 'D:\migrasi\RNM\NB FacIn\Section' -File -Filter *.xml).BaseName|ForEach-Object{$nbNames[$_]=$true}
$inNB = $sub | Where-Object { $nbNames.ContainsKey($_.Value) }
$grid=@{}; $rg | ForEach-Object { $grid[$_.Section]=$true }
$nest = $inNB | Where-Object { $grid.ContainsKey($_.Section) -and $grid.ContainsKey($_.Value) }
"pasangan=$((($nest|ForEach-Object{$_.Section+'>'+$_.Value})|Sort-Object -Unique).Count) induk=$((($nest.Section)|Sort-Object -Unique).Count)"
```

| Metrik graf penyisipan (NB) | Nilai |
| --- | ---: |
| Sel `SUB_SECTION` dengan nama | 892 |
| Nama target unik | 144 |
| — ada sebagai berkas Section NB | 138 |
| — komponen bawaan Pega / tidak ada di NB | 6 |
| Sisi Section→Section di dalam NB | **270** |
| **Pasangan grid-dalam-grid** | **121** |
| Section induk grid-dalam-grid | **33** |

Enam target di luar korpus NB: `pzPegaDefaultGridIcons` (535 sel — ikon baris grid bawaan Pega),
`pyGridPaginator` (82), `pzAttachFilesScreen` (2), `RatePropertyFacOut` (1), `EmailSectionEDM` (1),
`ErrorList` (1). Tiga terakhir **[pertanyaan terbuka]**: dirujuk dari NB tetapi berkasnya tidak ada
di folder NB.

Induk grid-dalam-grid terdalam:

| Section induk | Section-grid yang disisipkan |
| --- | ---: |
| `InputInwardFacultativeDtl` | 15 |
| `InputInwardFacultativeDtl_IsUW` | 13 |
| `InputEndorsementDtl` | 12 |
| `InputEndorsementDtl_IsUW` | 11 |
| `ViewTabGroupPolisFacOutSavior` | 10 |
| `OldDataEndorsementDtl` | 9 |
| `ViewInwardFacultativeDtl_IsUW` | 9 |

**Lampiran C** memuat seluruh 629 grid: Section, properti PageList, kelas baris.

---

## 6. Field read-only dan field tersembunyi permanen

### 6.1 Read-only keras di sel — [terverifikasi]

```powershell
$real | Group-Object EditOptions | Sort-Object Count -Descending
"RO=$(($real|Where-Object{$_.ReadOnly -eq 'true'}).Count)"
```

| `pyEditOptions` | Sel | % dari 6.314 |
| --- | ---: | ---: |
| **`Read-only`** | **4.834** | **76,6 %** |
| `Auto` | 1.276 | 20,2 % |
| `Editable` | 202 | 3,2 % |
| `Read-only-always` | 2 | 0,03 % |

`pyReadOnly=true` pada **4.836** sel, `false` pada 1.478.

**Temuan #3 — [terverifikasi]: layar NB pada dasarnya adalah layar baca.** Hanya **1.478 dari
6.314** field (23,4 %) yang dapat disunting di titik mana pun. Dari 393 Section yang berfield:

| Kategori | Section |
| --- | ---: |
| **Seluruh** field read-only (layar tampilan murni) | **183** |
| **Tidak ada** field read-only (layar isian murni) | 63 |
| Campuran | 147 |

### 6.2 Read-only bersyarat — [terverifikasi]

**1.141** sel field punya `pyReadOnlyCondition` terisi; **93** ekspresi berbeda.

| Kemunculan | Ekspresi |
| ---: | --- |
| 374 | `IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac` |
| **115** | **`1!=2`** *(selalu benar)* |
| 95 | `IsSpreadingUW` |
| 74 | `.FlagDelete==1` |
| 64 | `IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW` |
| 56 | `IsUW` |
| **32** | **`1==1`** *(selalu benar)* |
| 31 | `.FlagDelete = 1` |
| 30 | `.FlagDelete = '1'` |
| 21 | `IsGroupFac \|\| IsGroupUWFac` |
| **19** | **`1=1`** *(selalu benar)* |
| **15** | **`ALWAYS`** |
| 15 | `pyWorkPage.IsFacRetroOffer=='1'` |
| 13 | `IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW` |
| 12 | `IsUW \|\| IsSpreadingUW` |
| 11 | `.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsShowOpenPolicyMarine == true` |

**Tiga ejaan untuk satu maksud** juga ada di NB: `.FlagDelete==1` (74) + `.FlagDelete = 1` (31) +
`.FlagDelete = '1'` (30) = **135** kemunculan. Yang ketiga membandingkan **sebagai string**.
Sesuai CLAUDE.md §1 ini **kandidat perbaikan, bukan keputusan migrasi** — jangan diseragamkan
diam-diam.

**Temuan #4 — [dugaan]: 1.043 dari 1.141 kondisi read-only itu tidak berpengaruh.** 1.043 sel
memiliki `pyReadOnlyCondition` terisi **sekaligus** `pyReadOnly=true`; hanya **98** sel yang
`pyReadOnly=false` sehingga kondisinya benar-benar menentukan. Label [dugaan] karena semantik
presedensi Pega (mana yang menang antara flag sel dan kondisi) tidak dapat dibuktikan dari ekspor.

```powershell
"ROCond&RO=true : $(($real|Where-Object{$_.ROCond -ne '' -and $_.ReadOnly -eq 'true'}).Count)"
"ROCond&RO=false: $(($real|Where-Object{$_.ROCond -ne '' -and $_.ReadOnly -eq 'false'}).Count)"
```

### 6.3 Kondisi selalu-benar / selalu-salah — hitungan khusus NB

```powershell
$s='D:\migrasi\RNM\NB FacIn\Section'; $dead=0;$df=@{}; $ro=0;$rf=@{}
Get-ChildItem -LiteralPath $s -File -Filter *.xml | ForEach-Object {
  $x=[IO.File]::ReadAllText($_.FullName)
  $m=[regex]::Matches($x,'(?i)<pyCondition>\s*(1\s*[=]{1,2}\s*2|never)\s*</pyCondition>').Count
  if($m){$dead+=$m;$df[$_.BaseName]=$m}
  $n=[regex]::Matches($x,'(?i)<pyReadOnlyCondition>\s*(1\s*!=\s*2|1\s*[=]{1,2}\s*1|ALWAYS)\s*</pyReadOnlyCondition>').Count
  if($n){$ro+=$n;$rf[$_.BaseName]=$n} }
"MATI=$dead berkas=$($df.Count) ; SELALU-RO=$ro berkas=$($rf.Count)"
```

| Pola | Kemunculan (semua sel NB) | Berkas Section |
| --- | ---: | ---: |
| `pyCondition` = `1=2` / `1==2` / `never` → **tidak pernah tampil** | **407** | **105** |
| `pyReadOnlyCondition` = `1!=2` / `1=1` / `1==1` / `ALWAYS` → **selalu read-only** | **183** | **20** |

**Pemindaian wajib `(?i)`.** Ejaan huruf tidak seragam di korpus dan menyembunyikan puluhan
kemunculan bila diabaikan:

| Pola mati | Ejaan persis | n |
| --- | --- | ---: |
| `1=2` | `1=2` | 321 |
| `never` | **`never` huruf kecil 40** + `Never` 5 | 45 |
| `1 = 2` | `1 = 2` | 25 |
| `1==2` | `1==2` | 14 |
| `1 == 2` | `1 == 2` | 2 |
| **Jumlah** | | **407** |

| Pola selalu-read-only | Ejaan persis | n |
| --- | --- | ---: |
| `1!=2` | `1!=2` | 115 |
| `1==1` | `1==1` | 32 |
| `1=1` | `1=1` | 19 |
| `ALWAYS` | **`Always` 13** + `ALWAYS` 2 | 15 |
| `1 = 1` | `1 = 1` | 2 |
| **Jumlah** | | **183** |

Regex peka huruf memberi 367/94 dan 170/18 — **40 + 13 kemunculan luput**. Ini bukan sekadar
kebersihan data: mesin aturan Go yang mengevaluasi ekspresi ini harus memutuskan apakah
pembandingan namanya peka huruf atau tidak, dan korpus tidak menjawabnya. **[pertanyaan terbuka]**

> Angka lintas-siklus di dokumen arsip (936 dan 344) **tidak sebanding**: itu tiga folder digabung
> dan memakai pola regex berbeda. Angka NB di atas adalah yang berlaku untuk pass ini.

18 berkas dengan elemen mati terbanyak (dari 105): `InputInwardFacultativeDtl` (33),
`InputDtlObject_FacIn_IsUW` (25), `InputEndorsementDtl` (24), `DetailPolicyTreatyIn` (17),
`DetailDeptHeadTreatyIn_UW` (17), `GeneralDeptHeadTreatyIn_UW` (17), `GeneralPolicyTreatyIn` (17),
`InputDtlObject_FacIn` (14), `OldDataEndorsementDtl` (13), `InputOtherObjectAneka_FacIn` (12),
`InputOtherObjectAneka_FacIn_ISUW` (12), `InputEndorsementDtl_IsUW` (12), `ScoringRisk` (10),
`InputRiskAddress` (8), `Periode` (8), `ReasViewAttachment` (6), `PrintRISlip` (5),
`Periode_IsUW` (5).

Seluruh **20** berkas selalu-read-only, jumlah kemunculan **183**:
`InputDtlOtherObjectAneka_FacIn_ISUW` (**98**), `DetailDeptHeadTreatyIn_UW` (18),
`GeneralDeptHeadTreatyIn_UW` (18), `InputDtlObjectAneka_FacIn_IsUW` (14),
`InputDtlObjectAneka_FacIn` (12), `SummaryLossRecord_Section` (6), `InputAccumulationCov` (2),
`PersonListDetail` (2), `InputDtlParticipantPA_FacIn_IsUW` (2), `AddFacOfferList2` (1),
`InputEndorsementDtl_IsUW` (1), `DetailPolicyTreatyIn` (1), `FacOutPrintRISlipSectionInside` (1),
`InputFacOffer` (1), `Periode_IsUW` (1), `OccupationItemFacIn_Section` (1), `Periode` (1),
`AddFacOfferList_IsUW` (1), `ShowCedingCoList` (1), `GeneralPolicyTreatyIn` (1).

**Seluruh 183 kemunculan itu menempel pada sel yang mengikat properti** — angka ini identik dengan
hitungan tingkat-field di §6.4, jadi tidak ada kondisi selalu-read-only yang terbuang pada sel
label atau layout.

### 6.4 Terbatas pada field terikat saja

```powershell
$dead = $real | Where-Object { $_.Cond -match '^\s*(1\s*[=]{1,2}\s*2|Never)\s*$' }
"selMati=$($dead.Count) propUnik=$((($dead.Value)|Sort-Object -Unique).Count) section=$((($dead|Group-Object Section)).Count)"
$aro = $real | Where-Object { $_.ROCond -match '^\s*(1\s*!=\s*2|1\s*[=]{1,2}\s*1|ALWAYS)\s*$' }
"selSelaluRO=$($aro.Count) section=$((($aro|Group-Object Section)).Count)"
```

| Metrik (hanya field terikat nyata) | Nilai |
| --- | ---: |
| Field **tidak pernah tampil** | **106 sel**, 64 properti unik, **56 Section** |
| Field **selalu read-only lewat kondisi** | **183 sel**, **20 Section** |

Section dengan field mati terbanyak: `InputInwardFacultativeDtl` (15), `InputEndorsementDtl` (9),
`Periode` (6), `Periode_IsUW` (4), `ViewDataPolicy` (4), `InputRiskAddress` (3), `InwardFacIn` (3),
`PeriodeEndorsement` (3).

**Ini pertanyaan bisnis, bukan keputusan migrasi.** Membawa 106 field mati ke React berarti
membawa mati-suri; membuangnya berarti memutuskan sendiri bahwa bisnis tidak akan menyalakannya
lagi. Lihat §10-Q6.

### 6.5 Visibilitas bersyarat — [terverifikasi]

```powershell
$real | Group-Object Visible | Sort-Object Count -Descending
$real | Where-Object{$_.Cond -ne ''} | Group-Object Cond | Sort-Object Count -Descending | Select-Object -First 25
```

| `pyVisible` | Sel field |
| --- | ---: |
| `ALWAYS` | 5.203 |
| `OTHER` (ekspresi di `pyCondition`) | 909 |
| `NOTBLANK` | 199 |
| `NOTZERO` | 3 |

**1.051** sel punya `pyCondition` terisi (378 ekspresi berbeda), tetapi hanya 909 yang
`pyVisible=OTHER`. Selisih **142 sel** menyimpan ekspresi kondisi sementara `pyVisible=ALWAYS` —
**ekspresi itu tidak dievaluasi**. [dugaan] sisa penyuntingan; **jangan** diterjemahkan sebagai
kondisi aktif di React.

25 ekspresi visibilitas teratas pada field terikat:

| n | Ekspresi | n | Ekspresi |
| ---: | --- | ---: | --- |
| 251 | `1=2` | 14 | `1==2` |
| 28 | `!IsSpreadingUW` | 14 | `pyWorkPage.ProposalPosition = 'H1' \|\| … = '1' \|\| … = '1A'` |
| 25 | `IsEDM` | 13 | `IsLife` |
| 22 | `!IsClaim && .FlagDelete != '1'` | 13 | `IsHE` |
| 22 | `1 = 2` | 12 | `IsObjectLimiliabilityAneka` |
| 21 | `!IsUW` | 12 | `isGolfInsurance` |
| 21 | `IsCreditBriguna` | 12 | `IsObjectWithQuantityYear \|\| IsHE` |
| 20 | `.CoverageBasis==2` | 11 | `IsFac` |
| 18 | `.MinMax = 3` | 10 | `IsAneka OR IsBonding OR IsFire OR IsCustomBonds` |
| 17 | `IsObjectSectionAneka` | 10 | `OutputParam.ERRMSG6==''` |
| 16 | `Never` | 10 | `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy'` |
| 16 | `.IsNewPolicyNonProp != 1` | 10 | `.BalanceDueTo >= 0` |
| 14 | `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` | | |

Ekspresi `IsAneka OR IsBonding OR IsFire OR IsCustomBonds` memakai **`OR` kata**, sementara sisanya
memakai `||`. Dua dialek ekspresi hidup berdampingan. [terverifikasi]

---

## 7. Validasi yang terbaca — dan yang tidak ada sama sekali

```powershell
$s='D:\migrasi\RNM\NB FacIn\Section'
$tags=@('pyMaxLength','pyValidateAs','pyEditValidate','pyMaxChars','pyMinChars','pyRequired',
        'pyRequiredNew','pyPattern','pyMask','pyTextAreaMaxLength','pyCharacterLimit','pyMinValue','pyMaxValue')
$h=@{}; Get-ChildItem -LiteralPath $s -File -Filter *.xml | ForEach-Object { $x=[IO.File]::ReadAllText($_.FullName)
  foreach($t in $tags){ $n=([regex]::Matches($x,"<$t>([^<]*)</$t>")|Where-Object{$_.Groups[1].Value -notin @('','false')}).Count
    if($n){$h[$t]=$n+$h[$t]} } }
$h.GetEnumerator()|Sort-Object Value -Descending
```

| Tag validasi | Kemunculan non-kosong di 432 Section NB |
| --- | ---: |
| `pyRequired` | 928 |
| `pyRequiredNew` (mode kontrol) | 457 |
| `pyEditValidate` | **2** (keduanya `true` di `pyAttachmentScreen` — bukan nama rule) |
| `pyMaxLength` | **0** |
| `pyValidateAs` | **0** |
| `pyMaxChars` / `pyMinChars` / `pyCharacterLimit` / `pyTextAreaMaxLength` | **0** |
| `pyPattern` / `pyMask` | **0** |
| `pyMinValue` / `pyMaxValue` | **0** |

**Temuan #5 — [terverifikasi] dan MEMBLOKIR: lapisan Section tidak memuat satu pun aturan panjang,
format, rentang, atau pola.** Satu-satunya validasi yang terbaca adalah **wajib-isi**. Semua
batasan lain hidup di rule `Property`, `Edit Validate`, atau `Validate` — **dan tidak satu pun
tipe rule itu ada di ekspor NB** (folder yang ada hanya 14 tipe, §1.1).

### 7.1 Field wajib isi — [terverifikasi]

```powershell
$req = $real | Where-Object { $_.Required -eq 'true' }
"sel=$($req.Count) propUnik=$((($req.Value)|Sort-Object -Unique).Count) section=$((($req|Group-Object Section)).Count)"
```

**477 sel**, **191 properti unik**, tersebar di **73 Section**. Yang terbanyak:

| n | Properti | n | Properti |
| ---: | --- | ---: | --- |
| 14 | `.OfferFacIn.Obligee.ObligeeName` | 8 | `.OfferFacIn.Obligee.ContractDate` |
| 14 | `.OfferFacIn.Obligee.PublishedIn` | 8 | `.OfferFacIn.Obligee.ContractNo` |
| 14 | `.OfferFacIn.Obligee.UserAssignment` | 7 | `.Aneka.Obligee.Jabatan` |
| 14 | `.ASMDateOfBirth` | 7 | `.Aneka.Obligee.ObligeeAddress` |
| 14 | `.OfferFacIn.Obligee.Jabatan` | 7 | `.Aneka.Obligee.ProjectName` |
| 12 | `.OfferFacIn.Obligee.ObligeeAddress` | 7 | `.Aneka.Obligee.ObligeeName` |
| 12 | `.OfferFacIn.Obligee.PaymentClaim` | 6 | `.Aneka.Obligee.PaymentClaim` |
| 12 | `.OfferFacIn.Obligee.ProjectName` | 5 | `.pyFullName` |
| 10 | `.Aneka.Obligee.PublishedIn` | 4 | `.OfferFacIn.PercentShare` |
| 10 | `.Aneka.Obligee.UserAssignment` | 4 | `.StartDate` / `.EndDate` |

Blok `Obligee.*` (bonding/surety) adalah **satu-satunya kelompok field yang benar-benar wajib** di
NB. Sisanya tersebar tipis: `.Claim`, `.PremiOnp`, `.PremiOgp`, `.RiCommOnp`, `.RiCommOgp`,
`.OveriddingCommOnp`, `.OveriddingCommOgp`, `.Deduction1`, `.Deduction2`, `.SalvageValue`,
`.ExcessLoss`, `.OutstandingClaim`, `.StatementDate` (masing-masing 4) — semuanya pada layar
pernyataan hasil (statement) yang sama.

**Yang perlu diingat:** wajib-isi di Pega dievaluasi **saat submit flow action**, bukan saat
mengetik. Memindahkannya ke validasi klien React **mengubah perilaku**.

---

## 8. Tipe data: apa yang terbukti dan apa yang tidak

### 8.1 Korpus NB tidak memuat rule `Property` — ini temuan, bukan kekurangan pembacaan

Folder `D:\migrasi\RNM\NB FacIn\` memuat 14 tipe rule (§1.1). **`Property` bukan salah satunya.**
Karena itu **tipe deklaratif setiap properti adalah `belum terverifikasi`.** Satu-satunya bukti
tipe yang tersedia adalah **kontrol UI** yang mengikatnya.

### 8.2 Bukti tipe dari kontrol — [terverifikasi]

```powershell
$real | Group-Object Format | Sort-Object Count -Descending
```

| `pyFormat` | Sel | Bukti tipe yang diberikan |
| --- | ---: | --- |
| `pxNumber` | 1.770 | numerik desimal — **bukan bukti presisi** |
| `pxTextInput` | 1.370 | teks — bisa juga angka sebagai string |
| `pxDropdown` | 539 | enumerasi — **nilai enum tidak ada di korpus Section** |
| **`pxCheckbox`** | **534** | **boolean** (bukti terkuat) |
| `pxDisplayText` | 519 | hanya tampil; tipe tidak terbukti |
| *(tanpa `pyFormat`)* | 466 | **tidak ada bukti tipe** |
| `pxTextArea` | 370 | teks panjang |
| `pxAutoComplete` | 260 | teks + lookup |
| **`pxDateTime`** | **239** | **tanggal/waktu** |
| **`pxInteger`** | **81** | **bilangan bulat** |
| `pxRadioButtons` | 65 | enumerasi kecil |
| **`pxCurrency`** | **65** | **uang** |
| `pxLink` | 16 | navigasi |
| `pxRichTextEditor` | 11 | HTML |
| `pxHidden` | 4 | tersembunyi |
| `ReasAttachmentDescription_name` | 2 | kontrol kustom — isinya tidak terbaca |
| `ASMThumbnail` / `pxSlider` / `pxPassword` | 1 masing-masing | — |

| Metrik bukti tipe (atas 1.805 ekspresi unik) | Nilai |
| --- | ---: |
| Punya ≥1 kontrol pembuktian tipe (checkbox/date/int/number/currency/dropdown/radio) | **1.034** |
| **Tidak punya `pyFormat` sama sekali** | **66** |
| **Diikat dengan >1 jenis kontrol berbeda** (bukti bertentangan) | **150** |

### 8.3 Temuan #6 — [terverifikasi]: bukti tipe saling bertentangan

150 ekspresi diikat dengan lebih dari satu jenis kontrol. Yang paling sering:

| Ekspresi | Kontrol yang dipakai |
| --- | --- |
| `.Currency` | `pxAutoComplete`, `pxDisplayText`, `pxDropdown`, **`pxNumber`**, `pxTextInput` |
| `.TSI` | `pxDisplayText`, **`pxInteger`**, `pxNumber`, `pxTextInput` |
| `.Premium` | `pxDisplayText`, `pxNumber`, `pxTextInput` |
| **`.ASMDateOfBirth`** | **`pxDateTime` dan `pxInteger`** |
| `.CoverageNote` | `pxAutoComplete`, `pxDropdown`, `pxTextArea`, `pxTextInput` |
| `.ObjectName` | `pxAutoComplete`, `pxDisplayText`, `pxTextArea`, `pxTextInput` |
| `.Rate` | `pxDisplayText`, `pxNumber`, `pxTextInput` |
| `.ASMGender` | `pxDisplayText`, `pxDropdown` |
| `.TreatyType` | `pxDropdown`, `pxTextInput` |

Bukti `.ASMDateOfBirth` — **dua pengikatan dalam satu berkas yang sama**,
`D:\migrasi\RNM\NB FacIn\Section\CoverageSpreadingList_IsUW.xml`:

```powershell
$c=[IO.File]::ReadAllText('D:\migrasi\RNM\NB FacIn\Section\CoverageSpreadingList_IsUW.xml')
foreach($m in [regex]::Matches($c,'<pyValue>\.ASMDateOfBirth</pyValue>')){
  $seg=$c.Substring($m.Index,300)
  "idx=$($m.Index) format=$([regex]::Match($seg,'<pyFormat>([^<]*)</pyFormat>').Groups[1].Value)" }
# idx=214676 format=pxDateTime
# idx=418857 format=pxInteger
```

**`.Currency` diikat sebagai `pxNumber`** di sebagian layar. Bila `.Currency` menyimpan kode mata
uang, itu keliru; bila menyimpan kurs, namanya keliru. **Tidak dapat diputuskan dari korpus.**
**MEMBLOKIR** untuk tipe `Money` (§10-Q2).

### 8.4 Uang: presisi dan mata uang — [terverifikasi]

```powershell
$real | Where-Object{$_.Decimals -ne ''} | Group-Object Decimals | Sort-Object Name
$real | Where-Object{$_.CurrType -ne ''} | Group-Object CurrType
```

| `pyDecimalPlaces` | Sel field |
| ---: | ---: |
| **4** | **1.025** |
| **2** | **538** |
| 0 | 35 |
| **-999** | 8 |
| 3 | 2 |
| 5 | 2 |
| 6 | 1 |
| 10 | 1 |
| **-1** | 1 |
| **Total** | **1.613** |

`pyCurrencyType`: `local` **1.993** sel, `other` 4. Nilai `-999` dan `-1` **belum terverifikasi**
artinya.

Bukti format uang, `D:\migrasi\RNM\NB FacIn\Section\InputCoverageFire.xml`, sel `.RateOJK`:

```xml
<pySmartPromptClass>ASM-FW-GISFW-Data-Coverage</pySmartPromptClass>
<pyValue>.RateOJK</pyValue>
<pyFormat>pxNumber</pyFormat>
...
<pyDecimalPlaces>4</pyDecimalPlaces>
<pyCurrencyPosition>left</pyCurrencyPosition>
<pyCurrencyType>local</pyCurrencyType>
<pyNegativeFormat>none</pyNegativeFormat>
```

**Temuan #7 — [terverifikasi] dan MEMBLOKIR: presisi satu properti berbeda antar layar.**
**20 dari 439** properti bernilai desimal punya `pyDecimalPlaces` **tidak konsisten**:

| Properti | Desimal yang dipakai |
| --- | --- |
| `.Premium` | **0 / 2 / 4** |
| `.TSI` | 0 / 4 |
| `.TSIObjectItem` | 0 / 2 / 4 |
| `.SharePercentage` | **-1 / 2 / 4** |
| `.Rate` | 4 / 6 |
| `.RateLife` | **10 / 4** |
| `.LimitofLiability` | 4 / **-999** |
| `.PctDeductible2` | 2 / **-999** |
| `.PremiumSpreaded`, `.ClaimSpreaded`, `.Amount`, `.PaymentTotal`, `.ClaimPercentage`, `.GrossPremium` | 2 / 4 |
| `pyWorkPage.PolicyTreatyIn.TotalPremium`, `…TotalClaim`, `…TotalSharePercentagePremium`, `…TotalSharePercentageClaim` | 2 / 4 |
| `.SumTotalTSI`, `.MinPremium` | 0 / 4 |

Contoh terverifikasi: `.Premium` dirender `pxNumber` dengan `pyDecimalPlaces=2` di
`PropertyItemListFacOut`, tetapi dengan 4 desimal di mayoritas layar lain, dan dengan 0 desimal di
`PersonPlanAndDeduct` serta `InputDeductible_GCNM`. **Apakah ini hanya format tampilan atau
pembulatan yang benar-benar tersimpan tidak dapat dijawab dari korpus Section.**

**Temuan #8 — [terverifikasi]: mata uang tidak menyertai nilai uang.** Hanya **93** pengikatan
properti bernama mata uang di seluruh NB — `.Currency` 85, `.Currency.Currency` 3,
`pyWorkPage.Policy.Currency` 2, `.CurrencyMaster` 2, `.CurrencyID` 1 — berhadapan dengan **2.007**
sel numerik (543 ekspresi unik).

```powershell
$mre='(?i)(TSI|Premi|Premium|Amount|Claim|Comm|Komisi|Payment|Deduction|Fee|Salvage|Discount|Loading|Capital|Price)'
$money=$real|Where-Object{($_.Format -in @('pxCurrency','pxNumber') -or $_.Decimals -ne '') -and ($_.Value -replace '^.*\.','') -match $mre}
$secMoney=($money|Group-Object Section).Name
$curr=$real|Where-Object{($_.Value -replace '^.*\.','') -match '(?i)^(Currency|CurrencyID|CurrencyMaster|CurrencyName)$' -or $_.Value -match '(?i)\.Currency\.'}
$secCurr=($curr|Group-Object Section).Name
"uang=$($secMoney.Count) mataUang=$($secCurr.Count) keduanya=$((($secMoney|Where-Object{$secCurr -contains $_})).Count)"
```

| Metrik | Nilai |
| --- | ---: |
| Sel bernilai uang (kontrol numerik + nama seperti uang) | **899** (140 ekspresi unik) |
| Section yang memuat field uang | **170** |
| Section yang memuat field mata uang | 71 |
| Section yang memuat **keduanya** | 58 |
| **Section yang memuat uang TANPA mata uang** | **112** |

Ini adalah konfirmasi independen atas CLAUDE.md §4.1: **di 112 layar NB, angka uang dirender tanpa
mata uang mana pun di layar yang sama.** Model Go **wajib** membawa `Currency` bersama `Amount`
walaupun korpus tidak menyebutkannya, dan asal mata uang untuk 112 layar itu adalah
**[pertanyaan terbuka]**.

15 field uang paling sering, dengan desimalnya:

| n | Properti | Desimal | n | Properti | Desimal |
| ---: | --- | --- | ---: | --- | --- |
| 96 | `.Premium` | 0/2/4 | 16 | `.TSIOld` | 4 |
| 86 | `.TSI` | 0/4 | 16 | `.Amount` | 2/4 |
| 30 | `.PremiumSpreaded` | 2/4 | 15 | `.TSILiability` | 4 |
| 29 | `.TSISpreaded` | 4 | 14 | `.Policy.Payment.EDMPremiMenjadi` | 4 |
| 24 | `.PremiNusantaraRe` | 4 | 14 | `.TSIAddCap` | 4 |
| 23 | `.TSIObjectItem` | 0/2/4 | 14 | `.Policy.Payment.EDMOldPayment` | 4 |
| 20 | `.TSINusantaraRe` | 4 | 14 | `.OldTSI` | 4 |
| 19 | `.SumTotalPayment` | 4 | 14 | `.ClaimEstimation` | 4 |
| 18 | `.ClaimSpreaded` | 2/4 | 13 | `.Policy.Payment.Premium` | 4 |

`pxCurrency` hanya dipakai pada 17 properti, semuanya di layar pernyataan hasil:
`.ResultOnp1`, `.ResultOnp2`, `.NetPremium`, `.PremiOnp`, `.PPHValue`, `.PPNValue`, `.Deduction1`,
`.Deduction2`, `.ResultOgp1`, `.ResultOgp2`, `.ShareValue`, `.PremiOgp`, `.SalvageValue`,
`.ExcessLoss`, `.Claim`, `.OutstandingClaim` (4 sel masing-masing) dan
`SFAResults.pxResults(1).pySummaryValue(1)` (1). **Seluruh field uang lain memakai `pxNumber`** —
artinya "ini uang" **tidak dinyatakan di UI**, hanya tersirat dari nama.

### 8.5 Boolean dan tanggal yang terbukti

**270 ekspresi** diikat `pxCheckbox` → boolean [terverifikasi]. Terbanyak: `.IsTopRisk` (12),
`.IsMaterialDamage` (8), `.Property.IsTopRisk` (7), `.OfferFacIn.PPnCheck` (5), `.IsUsedFlag` (4),
`.PolicyData.LC.LC` (4), `.FlagPPH` (4), `.IsFlammableItemFlag` (3), `.IsProductionProcessFlag` (3),
`.IsAdjustableFlag` (3), `.IsShowDetail` (3), `.IsHotWorkProcessFlag` (3),
`.IsSpecialAcceptance` (3).

**71 ekspresi** diikat `pxDateTime` → tanggal [terverifikasi]. Terbanyak: `.ASMDateOfBirth` (50 —
tetapi lihat §8.3), `.DateOfLoss` (10), `.DueDate` (9), `.Aneka.Obligee.ContractDate` (8),
`.OfferFacIn.Obligee.ContractDate` (8), `.StartPeriod` (7), `.EndPeriod` (7), `.DateTransfer` (6),
`.PolicyData.SailDate` (5), `.DateofSurvey` (4), `.PolicyData.LC.LCEndDate` (4),
`.OfferFacIn.Obligee.AuctionDate` (4), `.PolicyData.InvoiceDate` (4),
`.Aneka.Obligee.AuctionDate` (4), `.UpdateDate` (4).

**539 sel `pxDropdown` + 65 `pxRadioButtons` = 604 sel enumerasi, dan tidak satu pun nilai enum
terbaca dari Section.** Sumber pilihan ada di `pyListDataSource` yang menunjuk
`pyListSource=associated` (Property rule — tidak ada) atau ke Data Page / Report Definition.
**MEMBLOKIR** (§10-Q3).

---

## 9. [usulan] Pemetaan ke React + Go

> **Seluruh §9 adalah USULAN.** Korpus tidak memuat struktur React atau Go mana pun. Dasarnya
> angka §2–§8; tiap baris menyebut angkanya agar dapat dibantah.

### 9.1 [usulan] Komponen React bersama

| Komponen | Dasar bukti | Catatan |
| --- | --- | --- |
| `<DataGrid>` (virtualisasi + baris expandable) | 629 RepeatGrid di 281 Section; 424 `pyExpandable=true`; 121 pasangan grid-dalam-grid | wajib membatasi kedalaman render secara eksplisit |
| `<MoneyInput>` / `<MoneyText>` | 899 sel uang; `pyDecimalPlaces` 4 (1.025) dan 2 (538) | presisi = **prop**, bukan konstanta; lihat §8.4 |
| `<PercentInput>` | `.SharePercentage`, `.PercentOffered`, `.ClaimPercentage`, `.PctLoL`, `.PercentageAdjustment` | desimal 2 vs 4 berbeda antar layar |
| `<DateField>` | 239 sel `pxDateTime`, 71 ekspresi | |
| `<CheckboxField>` | 534 sel, 270 ekspresi | |
| `<LookupSelect>` | 539 `pxDropdown` + 260 `pxAutoComplete` + 65 `pxRadioButtons` | sumber pilihan **belum terverifikasi** |
| `<ReadOnlyGuard mode={…}>` | 4.834 sel `Read-only`; 81 pasang layar `_IsUW` | satu komponen + prop `mode`, bukan komponen kembar |
| `<ConditionalField when={…}>` | 909 sel `pyVisible=OTHER`, 378 ekspresi | evaluator ekspresi, bukan `if` tersebar |
| `<PersonPanel>` | `Data-Party-Person` 442 sel di 33 Section PA/LIFE | dipakai PA, Life, Travel |
| `<ObligeePanel>` | blok `Obligee.*` — satu-satunya kelompok wajib-isi (§7.1) | |

### 9.2 [usulan] Bentuk model Go — satu inti + lima bagian opsional

Didukung rasio **963 field khas : 133 field bersama** (§4.4) dan peta segmen jalur (§3.2).

```go
// Asal: NB FacIn/Section/* -- pySmartPromptClass + segmen jalur pyValue
type OfferFacIn struct {
    // inti, lintas lini bisnis (790 nama daun khas INTI)
    Quotation   Quotation      // .Quotation. / .QuotationData.  (78 sel)
    Currencies  []CurrencyLine // .CurrencyList (37 grid, 3 ejaan jalur)
    Coverages   []Coverage     // .CoverageList (40 grid)
    Deductibles []Deductible   // .DeductibleList (17 grid)
    Spreading   []SpreadingRisk// .SpreadingList / .SpreadingRiskList (12 grid)
    Clauses     []Clause       // .ClauseList (12 grid)

    // bagian opsional per lini bisnis -- dipilih oleh Quotation.BusinessType (Sec 4.1)
    Property *PropertyRisk // .Property. (184 sel) + .SurroundingRisk. (60) + .ScoringRisk. (31)
    Cargo    *CargoRisk    // kelas Data-Cargo/Conveyance/Good/Packing/Trading (282 sel)
    Vehicles []Vehicle     // .VehicleList (16 grid), kelas Data-Vehicle (144 sel)
    Persons  []InsuredPerson // .PersonList (32 grid), kelas Data-Party-Person (411 sel)
    Aneka    *AnekaRisk    // .Aneka. (280 sel) + Obligee, kredit, golf, alat berat
}

// WAJIB bertipe uang (decimal + mata uang) -- Sec 8.4.
// Daftar lengkap: Lampiran B, kolom Kontrol = pxCurrency/pxNumber dengan pyDecimalPlaces terisi.
type Coverage struct {
    TSI              money.Money // .TSI        -- 86 sel, desimal 0/4 (TIDAK konsisten)
    Premium          money.Money // .Premium    -- 96 sel, desimal 0/2/4 (TIDAK konsisten)
    PremiNusantaraRe money.Money // .PremiNusantaraRe -- 24 sel, desimal 4
    TSINusantaraRe   money.Money // .TSINusantaraRe   -- 20 sel, desimal 4
    Rate             decimal.Decimal // .Rate    -- desimal 4/6 (TIDAK konsisten)
    // ... lihat Lampiran A: CoverageItem (48 field, 42 unik)
}
```

**Field yang WAJIB bertipe uang** (kriteria: kontrol `pxCurrency`, atau `pxNumber`/`pyDecimalPlaces`
dengan nama uang) — **140 ekspresi unik**, seluruhnya terdaftar di Lampiran B. Yang teratas:
`.Premium`, `.TSI`, `.PremiumSpreaded`, `.TSISpreaded`, `.PremiNusantaraRe`, `.TSIObjectItem`,
`.TSINusantaraRe`, `.SumTotalPayment`, `.ClaimSpreaded`, `.TSIOld`, `.Amount`, `.TSILiability`,
`.TSIAddCap`, `.OldTSI`, `.ClaimEstimation`, `.TotalGrossPremi`, `.Discount`, `.PremiumOld`,
`.NetPremium`, `.PremiOnp`, `.PremiOgp`, `.Claim`, `.OutstandingClaim`, `.SalvageValue`,
`.ExcessLoss`, `.PPHValue`, `.PPNValue`, `.Deduction1`, `.Deduction2`, `.ShareValue`.

### 9.3 [usulan] Lima keputusan yang harus diambil sebelum kode ditulis

1. **Presisi uang** — simpan presisi penuh dan format di UI, atau bulatkan saat simpan? Korpus
   menunjukkan 0/2/4 desimal untuk properti yang sama (§8.4). Memilih salah satunya secara diam-diam
   akan memunculkan selisih saat rekonsiliasi paralel run.
2. **Mata uang** — 112 layar menampilkan uang tanpa mata uang. Dari mana `Currency` diisi?
3. **Enumerasi** — 604 sel dropdown/radio tanpa satu pun nilai yang terbaca.
4. **Kedalaman grid** — 121 pasangan grid-dalam-grid perlu batas render eksplisit di React;
   di Pega dipotong kondisi visibilitas sisi-server.
5. **Field mati** — 106 field yang tidak pernah tampil dan 183 yang selalu read-only: dibawa atau
   tidak? Keputusan bisnis.

---

## 10. Pertanyaan terbuka

| # | Pertanyaan | Status | Dampak |
| --- | --- | --- | --- |
| **Q1** | `BusinessType` dibaca dari **tiga jalur berbeda** (`.Quotation.`, `.OfferFacIn.QuotationData.`, `.Policy.Quotation.`) oleh `IsMarineCargo` dan `IsCustomBonds` (§4.1). Apakah ketiganya selalu sama? | **MEMBLOKIR** | menentukan bentuk `Quotation` dan pemilihan bagian opsional lini bisnis |
| **Q2** | `.Premium` dirender dengan **0, 2, dan 4 desimal** pada layar berbeda; `.SharePercentage` dengan `-1/2/4`; `.RateLife` dengan `10/4` (§8.4). Mana pembulatan yang tersimpan? Apa arti `-999` dan `-1`? | **MEMBLOKIR** | tipe `money.Money`; selisih angka saat paralel run |
| **Q3** | **604 sel** `pxDropdown`/`pxRadioButtons` tanpa satu pun nilai enumerasi di korpus (§8.5). | **MEMBLOKIR** | tidak ada layar isian yang dapat dibangun tanpa daftar pilihan |
| **Q4** | Korpus **tidak memuat rule `Property`** (§8.1) — tipe 1.805 ekspresi properti `belum terverifikasi`; 66 di antaranya bahkan tanpa `pyFormat`. | **MEMBLOKIR** | model Go dan skema JSON tidak dapat dikunci |
| **Q5** | Tidak ada satu pun aturan **panjang, format, rentang, pola** di lapisan Section (§7). Di mana validasi sesungguhnya? | **MEMBLOKIR** | form React akan menerima masukan yang ditolak backend |
| **Q6** | **106 field tidak pernah tampil** (`1=2`/`Never`) dan **183 field selalu read-only** (§6.3, §6.4). Disengaja atau sisa? | penting | dibawa = mati-suri; dibuang = keputusan bisnis diambil diam-diam |
| Q7 | `.Currency` diikat sebagai **`pxNumber`** di sebagian layar (§8.3). Kode mata uang atau kurs? | penting | penamaan field `Money.Currency` |
| Q8 | **112 Section** menampilkan uang tanpa field mata uang mana pun (§8.4). | penting | dari mana UI mengambil mata uang untuk ditampilkan? |
| Q9 | `.CurrencyList` / `.PersonList` / `.LocationList` muncul dengan **2–3 ejaan jalur** (§5.2). Daftar yang sama? | penting | menentukan apakah ada duplikasi daftar di agregat |
| Q10 | Empat Section besar tanpa satu pun field terikat (`FinalPolicyRO` 46 KB, `ComfirmPolis_Life` 50 KB, `ConfirmChangeOutgo` 49 KB, `DetailPoliciesNonProportional`) (§2.4). | sedang | mungkin memuat HTML non-autogenerated yang belum terbaca |
| Q11 | Tiga section dirujuk dari NB tetapi berkasnya tidak ada di folder NB: `RatePropertyFacOut`, `EmailSectionEDM`, `ErrorList` (§5.3). | sedang | celah korpus |
| Q12 | Halaman `test` mengikat 1 field (§1.3). Sisa uji coba di produksi? | rendah | kebersihan |
| Q13 | **142 sel** menyimpan `pyCondition` sementara `pyVisible=ALWAYS` — ekspresi tidak dievaluasi (§6.5). | rendah | jangan diterjemahkan sebagai kondisi aktif |
| Q14 | `ScoringRisk` mengkodekan **indeks baris di dalam ekspresi properti** (`.DataScoringRiskList(2)…`) (§2.2). Jumlah baris skoring memang tetap? | penting | menentukan apakah skoring adalah slice atau struct tetap |
| Q15 | 1.043 dari 1.141 `pyReadOnlyCondition` menyertai `pyReadOnly=true` (§6.2) — presedensi Pega tidak terbukti dari ekspor. | sedang | 1.043 kondisi mungkin dapat dibuang, mungkin tidak |
| Q16 | Ekspresi kondisi ditulis dengan ejaan huruf tidak seragam: `never`/`Never`, `Always`/`ALWAYS`, `OR` vs `\|\|`, `1=2`/`1 = 2`/`1==2`/`1 == 2` (§6.3, §6.5). Apakah evaluator Pega peka huruf? | penting | mesin aturan Go: satu keputusan salah menghidupkan/mematikan puluhan elemen |

---

## Lampiran A — Skema field per Section (432 Section)

Kolom:
`Properti` = ekspresi `pyValue` sel ·
`Kontrol` = `pyFormat` (digabung bila satu properti dipakai dengan beberapa kontrol) ·
`n` = jumlah sel yang mengikat properti itu di Section ini ·
`Read-only` = `pyReadOnly` (`campur` bila keduanya muncul); dalam kurung siku = `pyReadOnlyCondition` ·
`Kondisi tampil` = `pyCondition`, atau `pyVisible` bila bukan `ALWAYS`.

Placeholder `.pyTemplate*` **tidak** ditampilkan (§0.3a). Ekspresi dan kondisi dipotong bila
terlalu panjang (ditandai `...`); teks utuhnya ada di berkas XML yang bersangkutan.

#### AccumulationRisk

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Search` - field terikat 21 (unik 20) - read-only 20 - RepeatGrid 2 - grid: AccumulationRiskDtl.pxResults, AccumulationRiskPolicyDtl.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CARI1` | pxDisplayText | 2 | ya | ALWAYS |
| `.CARI10` | pxDisplayText | 1 | ya | ALWAYS |
| `.CARI11` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CARI12` | pxNumber | 1 | ya | ALWAYS |
| `.CARI13` | pxDisplayText | 1 | ya | ALWAYS |
| `.CARI14` | pxDisplayText | 1 | ya | ALWAYS |
| `.CARI2` | pxNumber | 1 | ya | ALWAYS |
| `.CARI3` | pxNumber | 1 | ya | ALWAYS |
| `.CARI4` | pxNumber | 1 | ya | ALWAYS |
| `.CARI5` | pxNumber | 1 | ya | ALWAYS |
| `AccumulationRisk.AccumulationType` | pxDisplayText | 1 | ya | ALWAYS |
| `AccumulationRisk.CenterTransStatus` | pxDisplayText | 1 | ya | ALWAYS |
| `AccumulationRisk.CZone` | pxDisplayText | 1 | ya | ALWAYS |
| `AccumulationRisk.ID` | pxDisplayText | 1 | ya | ALWAYS |
| `AccumulationRisk.Keyword` | pxDisplayText | 1 | ya | ALWAYS |
| `AccumulationRisk.Note` | pxTextInput | 1 | ya | ALWAYS |
| `AccumulationRisk.SyariahStatus` | pxTextInput | 1 | tidak | ALWAYS |
| `CountAccumulation.CARI2` | pxNumber | 1 | ya | ALWAYS |
| `CountAccumulation.CARI3` | pxNumber | 1 | ya | ALWAYS |
| `CountAccumulation.CARI4` | pxNumber | 1 | ya | ALWAYS |

#### addDeductible

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 11 (unik 11) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | tidak | .TypeDeductible != 7 |
| `.Condition` | pxDropdown | 1 | tidak | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [.TypeDeductible = 0] | ALWAYS |
| `.Descriptions` | pxTextArea | 1 | tidak | .Descriptions != '' && pyWorkPage.OfferFacIn.IsB2B !='' |
| `.InputCondition` | pxTextInput | 1 | tidak | .Condition = 5 |
| `.MinMax` | pxDropdown | 1 | ya [.TypeDeductible = 0 \|\| .TypeDeductible = 7] | .TypeDeductible != 7 |
| `.PctDeductible` | pxNumber | 1 | tidak | .TypeDeductible != 7 |
| `.PctDeductible2` | pxNumber | 1 | tidak | .MinMax = 3 |
| `.TimeExcess` | pxNumber | 1 | tidak | ALWAYS |
| `.TypeDeductible` | pxDropdown | 1 | tidak | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | tidak | .MinMax = 3 |

#### addDeductible_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 11 (unik 11) - read-only 11 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | ya [pyWorkPage.OfferFacIn.ConfirmBinding = 1] | .TypeDeductible != 7 |
| `.Condition` | pxDropdown | 1 | ya [pyWorkPage.OfferFacIn.ConfirmBinding = 1] | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [.TypeDeductible = 0] | ALWAYS |
| `.Descriptions` | pxDisplayText | 1 | ya | .Descriptions != '' && pyWorkPage.OfferFacIn.IsB2B !='' |
| `.InputCondition` | pxTextInput | 1 | ya [pyWorkPage.OfferFacIn.ConfirmBinding = 1] | .Condition = 5 |
| `.MinMax` | pxDropdown | 1 | ya [.TypeDeductible = 0 \|\| .TypeDeductible = 7] | .TypeDeductible != 7 |
| `.PctDeductible` | pxNumber | 1 | ya [pyWorkPage.OfferFacIn.ConfirmBinding = 1] | .TypeDeductible != 7 |
| `.PctDeductible2` | pxNumber | 1 | ya [pyWorkPage.OfferFacIn.ConfirmBinding = 1] | .MinMax = 3 |
| `.TimeExcess` | pxNumber | 1 | ya [pyWorkPage.OfferFacIn.ConfirmBinding = 1] | ALWAYS |
| `.TypeDeductible` | pxDropdown | 1 | ya [pyWorkPage.OfferFacIn.ConfirmBinding = 1] | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | ya [pyWorkPage.OfferFacIn.ConfirmBinding = 1] | .MinMax = 3 |

#### addDeductibleCargo

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 10 (unik 10) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | tidak | ALWAYS |
| `.Condition` | pxDropdown | 1 | tidak | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [.TypeDeductible = 0] | ALWAYS |
| `.InputCondition` | pxTextInput | 1 | tidak | .Condition = 5 |
| `.MinMax` | pxDropdown | 1 | ya [.TypeDeductible =0 \|\| .TypeDeductible = 7] | .TypeDeductible !=7 |
| `.PctDeductible` | pxNumber | 1 | tidak | .PctDeductible !=7 |
| `.PctDeductible2` | pxNumber | 1 | tidak | .MinMax = 3 |
| `.TimeExcess` | pxNumber | 1 | tidak | ALWAYS |
| `.TypeDeductible` | pxDropdown | 1 | tidak | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | tidak | .MinMax = 3 |

#### addDeductibleCargo_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 10 (unik 10) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | tidak | ALWAYS |
| `.Condition` | pxDropdown | 1 | tidak | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [.TypeDeductible = 0] | ALWAYS |
| `.InputCondition` | pxTextInput | 1 | tidak | .Condition = 5 |
| `.MinMax` | pxDropdown | 1 | ya [.TypeDeductible = 0 \|\| .TypeDeductible = 7] | .TypeDeductible != 7 |
| `.PctDeductible` | pxNumber | 1 | tidak | .TypeDeductible != 7 |
| `.PctDeductible2` | pxNumber | 1 | tidak | .MinMax = 3 |
| `.TimeExcess` | pxNumber | 1 | tidak | ALWAYS |
| `.TypeDeductible` | pxDropdown | 1 | tidak | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | tidak | .MinMax = 3 |

#### AddFacOfferList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-FacOffer` - field terikat 18 (unik 17) - read-only 9 - RepeatGrid 2 - grid: .CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Attention` | pxTextArea | 1 | tidak | ALWAYS |
| `.EndPeriod` | pxDateTime | 1 | tidak | ALWAYS |
| `.Name` | pxTextInput | 2 | ya | ALWAYS |
| `.NoOfferSlip` | pxTextInput | 1 | tidak | ALWAYS |
| `.PctPremiAllObj` | pxNumber | 1 | tidak | ALWAYS |
| `.PctPremiAllObjUSD` | pxNumber | 1 | tidak | ALWAYS |
| `.PctShareAllObj` | pxNumber | 1 | tidak | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.NetPremium` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.ReinsurerName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.RiCommAllObj` | pxNumber | 1 | tidak | ALWAYS |
| `.StartPeriod` | pxDateTime | 1 | tidak | ALWAYS |
| `.SumTotalPayment` | pxNumber | 1 | ya | ALWAYS |
| `.TotalOffered` | pxNumber | 1 | ya | NOTZERO |

#### AddFacOfferList_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 13 (unik 12) - read-only 13 - RepeatGrid 2 - grid: .CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Attention` | pxTextArea | 1 | ya | ALWAYS |
| `.Name` | pxTextInput | 2 | ya | ALWAYS |
| `.PctShareAllObj` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.NetPremium` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.ReinsurerName` | pxAutoComplete | 1 | ya [1=1] | ALWAYS |
| `.RiCommAllObj` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 1 | ya | ALWAYS |
| `.TotalOffered` | pxNumber | 1 | ya | NOTZERO |

#### AddFacOfferList2

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 13 (unik 12) - read-only 13 - RepeatGrid 2 - grid: .CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Attention` | pxTextArea | 1 | ya | ALWAYS |
| `.Name` | pxTextInput | 2 | ya | ALWAYS |
| `.PctShareAllObj` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.NetPremium` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.ReinsurerName` | pxAutoComplete | 1 | ya [1=1] | ALWAYS |
| `.RiCommAllObj` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 1 | ya | ALWAYS |
| `.TotalOffered` | pxNumber | 1 | ya | NOTZERO |

#### additemButton

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### AdjustmentRiskAccumulation

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Search` - field terikat 26 (unik 14) - read-only 26 - RepeatGrid 3 - grid: RiskAccumAdj.pxResults, RiskAccumCurYear.pxResults, RiskAccumPrevYear.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CARI1` | pxDateTime/pxDisplayText | 3 | ya | ALWAYS |
| `.CARI10` | pxDisplayText | 3 | ya | ALWAYS |
| `.CARI11` | pxDisplayText/pxNumber | 3 | ya | ALWAYS |
| `.CARI12` | pxNumber | 3 | ya | ALWAYS |
| `.CARI13` | pxNumber | 3 | ya | ALWAYS |
| `.CARI14` | pxDisplayText/pxNumber | 3 | ya | ALWAYS |
| `.CARI15` | pxDisplayText | 1 | ya | ALWAYS |
| `AdjustmentRisk.CARI18` | pxDisplayText | 1 | ya | ALWAYS |
| `AdjustmentRisk.CARI30` | pxNumber | 1 | ya | ALWAYS |
| `AdjustmentRisk.CARI31` | pxNumber | 1 | ya | ALWAYS |
| `AdjustmentRisk.CARI32` | pxNumber | 1 | ya | ALWAYS |
| `AdjustmentRisk.CARI33` | pxNumber | 1 | ya | ALWAYS |
| `AdjustmentRisk.CARI34` | pxNumber | 1 | ya | ALWAYS |
| `AdjustmentRisk.CARI35` | pxNumber | 1 | ya | ALWAYS |

#### AllSummarySection

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### AnalysLocationbyAI_Sec

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 2 (unik 2) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.HasilAI` | pxDisplayText | 1 | ya | NOTBLANK |
| `.pyScore` | pxSlider | 1 | tidak | ALWAYS |

#### AnekaListFacOut

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 6 (unik 6) - read-only 4 - RepeatGrid 1 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | tidak | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 1 | ya | ALWAYS |

#### AnekaListFacOut_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 1 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 1 | ya | ALWAYS |

#### AttachmentGridReas

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-OFFERJSON` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 1 - grid: Attachment.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CountAttach` | pxTextInput | 1 | tidak | ALWAYS |
| `.NOTE` | pxTextInput | 1 | tidak | ALWAYS |

#### BenefitClauseRO

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: .ClauseList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClauseContent` | pxDisplayText | 1 | ya | ALWAYS |
| `.ClauseParam` | pxDisplayText | 1 | ya | .ClauseType == '1' |

#### BenefitListRO

rumpun **INTI** - kelas dominan `(kosong)` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 1 - grid: .BenefitList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Limit` | pxNumber | 1 | tidak | ALWAYS |
| `.Name` | pxTextInput | 1 | tidak | ALWAYS |

#### CallDualScoringRisk

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### CauseOfLoss_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-CauseOfLoss` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 1 - grid: .Property.ListCauseOfLoss

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | ya | ALWAYS |
| `.Claim` | pxNumber | 1 | ya | ALWAYS |
| `.CoinsData.CoinsName` | pxTextInput | 1 | ya | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.DateOfLoss` | pxDateTime | 1 | ya | ALWAYS |
| `.LossObject` | pxTextInput | 1 | ya | ALWAYS |

#### CauseOfLoss_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-CauseOfLoss` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 1 - grid: .Property.ListCauseOfLoss

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | ya | ALWAYS |
| `.Claim` | pxNumber | 1 | ya | ALWAYS |
| `.CoinsData.CoinsName` | pxTextInput | 1 | ya | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.DateOfLoss` | pxDateTime | 1 | ya | ALWAYS |
| `.LossObject` | pxTextInput | 1 | ya | ALWAYS |

#### CauseOfLossClaim_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-CauseOfLossClaim` - field terikat 11 (unik 10) - read-only 11 - RepeatGrid 1 - grid: .Property.ListCauseOfLossClaim

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AkseptasiKlaim` | pxNumber | 1 | ya | ALWAYS |
| `.Currency` | (tanpa) | 1 | ya | ALWAYS |
| `.DateOfLoss` | pxDateTime | 2 | ya | ALWAYS |
| `.EstiamsiKlaim` | pxNumber | 1 | ya | ALWAYS |
| `.InccuredKlaim` | pxNumber | 1 | ya | ALWAYS |
| `.Location` | (tanpa) | 1 | ya | ALWAYS |
| `.LocationNo` | (tanpa) | 1 | ya | ALWAYS |
| `.LossRatio` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Remark` | (tanpa) | 1 | ya | ALWAYS |

#### CedingCedant

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 14 (unik 8) - read-only 12 - RepeatGrid 2 - grid: .CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CedingCoName` | pxDisplayText | 2 | ya | ALWAYS |
| `.Name` | (tanpa) | 2 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.ShareCeding` | pxNumber | 2 | tidak [IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSp...] | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.QuotationData.ShareOfCeding` | pxTextInput | 2 | ya | ALWAYS |

#### CedingCedant_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 11 (unik 8) - read-only 11 - RepeatGrid 2 - grid: .CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CedingCoName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Name` | (tanpa) | 2 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.ShareCeding` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.QuotationData.ShareOfCeding` | pxTextInput | 1 | ya | ALWAYS |

#### CedingCedantHierarki

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Agent` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClientName` | pxTextInput | 1 | tidak | ALWAYS |
| `SearchSOB.CARI1` | pxAutoComplete | 1 | tidak | ALWAYS |

#### CedingCoHierarki

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-AGENT` - field terikat 6 (unik 4) - read-only 3 - RepeatGrid 1 - grid: pgRepPgSubSectionCedingCoHierarkiBBBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClientID` | (tanpa) | 1 | ya | ALWAYS |
| `.ClientName` | pxTextInput | 2 | campur | ALWAYS |
| `.ID` | (tanpa) | 1 | ya | ALWAYS |
| `SearchSOB.CARI1` | pxAutoComplete | 2 | tidak | ALWAYS |

#### ChooseAccumulation_FacIn

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### ChooseClassofContraction

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-TABLEOFLIMIT` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: pgRepPgSubSectionChooseClassofContractionBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Description` | (tanpa) | 1 | ya | ALWAYS |
| `.PctLimit` | pxNumber | 1 | ya | Never |

#### ChooseClause

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-CLAUSE` - field terikat 3 (unik 3) - read-only 0 - RepeatGrid 1 - grid: clausePage.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ID` | pxTextInput | 1 | tidak | ALWAYS |
| `.Info` | pxTextInput | 1 | tidak | ALWAYS |
| `.IsSelected` | pxCheckbox | 1 | tidak | ALWAYS |

#### ChooseClauseFire

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-CLAUSE` - field terikat 6 (unik 6) - read-only 0 - RepeatGrid 1 - grid: SearchClauseListOutput.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ID` | pxTextInput | 1 | tidak | ALWAYS |
| `.Info` | pxTextInput | 1 | tidak | ALWAYS |
| `.IsSelected` | pxCheckbox | 1 | tidak | ALWAYS |
| `.Title` | pxTextInput | 1 | tidak | ALWAYS |
| `Clause.ClauseLanguageID` | pxDropdown | 1 | tidak [IsSpreadingUW] | ALWAYS |
| `Clause.pyNote` | pxTextInput | 1 | tidak [IsSpreadingUW] | .pxListSubscript = 1 |

#### ChooseCoverage

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-COVERAGE_FACIN` - field terikat 4 (unik 4) - read-only 3 - RepeatGrid 1 - grid: pgRepPgSubSectionChooseCoverageBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ID` | (tanpa) | 1 | ya | ALWAYS |
| `.NamaCoverage` | (tanpa) | 1 | ya | ALWAYS |
| `.OLDID` | (tanpa) | 1 | ya | ALWAYS |
| `CARI.pyName` | pxTextInput | 1 | tidak | ALWAYS |

#### ChooseDeductible

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-DEDUCTIBLE` - field terikat 3 (unik 3) - read-only 2 - RepeatGrid 1 - grid: pgRepPgSubSectionChooseDeductibleBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DeductName` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ID` | pxDisplayText | 1 | ya | NOTBLANK |
| `TempIndex.BizCode` | pxTextInput | 1 | tidak | ALWAYS |

#### ChooseObject_Ship

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Int-SHIP` - field terikat 5 (unik 5) - read-only 4 - RepeatGrid 1 - grid: pgRepPgSubSectionChooseObject_ShipBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.FlagName` | (tanpa) | 1 | ya | ALWAYS |
| `.ID` | (tanpa) | 1 | ya | ALWAYS |
| `.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.YearMake1` | (tanpa) | 1 | ya | ALWAYS |
| `CARI.pyName` | pxTextInput | 1 | tidak | ALWAYS |

#### ChooseOccupation

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Int-OCCUPATION` - field terikat 3 (unik 3) - read-only 2 - RepeatGrid 1 - grid: pgRepPgSubSectionChooseOccupationBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.OldID` | (tanpa) | 1 | ya | ALWAYS |
| `CARI.pyName` | pxTextInput | 1 | tidak | ALWAYS |

#### ChooseOccupation_Dtl

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Int-OCCUPATION` - field terikat 3 (unik 3) - read-only 0 - RepeatGrid 1 - grid: D_Occupation.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ID` | (tanpa) | 1 | tidak | ALWAYS |
| `.Name` | (tanpa) | 1 | tidak | ALWAYS |
| `InputParam.CARI5` | pxTextInput | 1 | tidak | ALWAYS |

#### ChooseRiskAddress

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 7 (unik 7) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.Country` | pxTextInput | 1 | tidak | ALWAYS |
| `.Property.Province` | pxTextInput | 1 | tidak | InputRiskAddress.IDNation!='' |
| `.Property.RiskLocation.ASMCity` | pxTextInput | 1 | tidak | InputRiskAddress.IDProvince!='' |
| `.Property.RiskLocation.ASMDistrict` | pxTextInput | 1 | tidak | InputRiskAddress.IDCity!='' |
| `.Property.RiskLocation.ASMRW` | pxTextInput | 1 | tidak | InputRiskAddress.IDDistrict!='' |
| `.Property.RiskLocation.ASMZipCode` | pxTextInput | 1 | tidak | ALWAYS |
| `.Property.RoadName` | pxTextInput | 1 | tidak | ALWAYS |

#### ChooseRiskAddress_ResultList

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Int-RISKADDRESS` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 1 - grid: pgRepPgSubSectionChooseRiskAddress_ResultListB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Address` | pxDisplayText | 1 | ya | ALWAYS |
| `.CityName` | pxDisplayText | 1 | ya | ALWAYS |
| `.DistrictName` | pxDisplayText | 1 | ya | ALWAYS |
| `.NationName` | pxDisplayText | 1 | ya | ALWAYS |
| `.PostalCode` | pxDisplayText | 1 | ya | ALWAYS |
| `.ProvinceName` | pxDisplayText | 1 | ya | ALWAYS |
| `.TerritoryName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Title` | pxDisplayText | 1 | ya | ALWAYS |

#### ChooseSubContract

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-SUBCONTRACT` - field terikat 3 (unik 3) - read-only 2 - RepeatGrid 1 - grid: pyReportContentPage.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ID` | pxTextInput | 1 | ya | ALWAYS |
| `.IS_SELECTED` | pxCheckbox | 1 | tidak | ALWAYS |
| `.Name` | pxTextInput | 1 | ya | ALWAYS |

#### ChooseSubContract_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-SUBCONTRACT` - field terikat 3 (unik 3) - read-only 2 - RepeatGrid 1 - grid: pyReportContentPage.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ID` | pxTextInput | 1 | ya | ALWAYS |
| `.IS_SELECTED` | pxCheckbox | 1 | tidak | ALWAYS |
| `.Name` | pxTextInput | 1 | ya | ALWAYS |

#### ChooseZipCodeDtl

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Int-RISKADDRESS` - field terikat 4 (unik 4) - read-only 0 - RepeatGrid 1 - grid: pgRepPgSubSectionChooseZipCodeDtlBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CityName` | (tanpa) | 1 | tidak | ALWAYS |
| `.NationName` | (tanpa) | 1 | tidak | ALWAYS |
| `.PostalCode` | (tanpa) | 1 | tidak | ALWAYS |
| `.ProvinceName` | (tanpa) | 1 | tidak | ALWAYS |

#### ClauseList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-M_KLAUSUL_PLAN` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: D_PlanClause.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.KLAUSUL_NOTE` | pxDisplayText | 1 | ya | ALWAYS |

#### ClauseValue

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 12 (unik 2) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClauseDescription` | pxAutoComplete | 9 | tidak | .ClauseCode == 100014 ; .ClauseCode == 100017 ; .ClauseCode == 100080 ... |
| `.ClauseParam` | pxRadioButtons/pxTextInput | 3 | tidak [pyWorkPage.ProposalPosition = 2 \|\| pyWorkPage....] | .ClauseCode != 100083 &&.ClauseCode != 100101 && .ClauseCode != 100047... |

#### ComfirmPolis_Life

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### ComfirmPolis_NonLife

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### CommisionCoverageAneka

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |

#### CommisionItem

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 10 (unik 8) - read-only 5 - RepeatGrid 1 - grid: .LayerList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.LayerNo` | (tanpa) | 1 | tidak | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.PremiRp` | pxNumber | 1 | ya | IsMBU |
| `.Premium` | pxNumber/pxTextInput | 2 | campur | !IsMBU |
| `.Rate` | pxTextInput | 1 | tidak | ALWAYS |
| `.TSI` | pxTextInput | 1 | tidak | ALWAYS |
| `.TSILiability` | pxNumber/pxTextInput | 2 | campur | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### CommisionItemAneka

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### CommisionObjectAneka

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 4 (unik 4) - read-only 3 - RepeatGrid 1 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | tidak | ALWAYS |
| `.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |

#### ConfirmChangeOutgo

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### confirmCreateAccumulation

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 1 (unik 1) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `OutputParam.ERRMSG3` | pxTextInput | 1 | ya | NOTBLANK |

#### CopyCoverageFrom

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.SelectedLocationAddress` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.SelectedObjectItem` | pxAutoComplete | 1 | tidak | .SelectedLocationAddress != '' |

#### Correspondence

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Correspondence` - field terikat 5 (unik 5) - read-only  - RepeatGrid 1 - grid: .OfferFacIn.CorrespondenceList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Date` | (tanpa) | 1 | tidak | ALWAYS |
| `.IsEditable` | pxDropdown | 1 | ya | ALWAYS |
| `.Recipient` | (tanpa) | 1 | tidak | ALWAYS |
| `.Sender` | (tanpa) | 1 | tidak | ALWAYS |
| `.Subject` | (tanpa) | 1 | tidak | ALWAYS |

#### CorrespondenceContent

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Correspondence` - field terikat 10 (unik 8) - read-only 2 - RepeatGrid 3 - grid: .AttachmentList, .RecipientList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Message` | pxRichTextEditor | 1 | ya [.IsEditable != 1] | ALWAYS |
| `.pyCategory` | (tanpa) | 2 | tidak | ALWAYS |
| `.pyFileName` | pxLink | 2 | tidak | ALWAYS |
| `.Recipient` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.RecipientMail` | (tanpa) | 1 | ya | ALWAYS |
| `.Sender` | pxTextInput | 1 | tidak | ALWAYS |
| `.SenderMail` | pxTextInput | 1 | tidak | ALWAYS |
| `.Subject` | pxTextInput | 1 | tidak | ALWAYS |

#### CoverageCommisionList

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 28 (unik 23) - read-only 27 - RepeatGrid 7 - grid: .CargoList, .LocationList, .PersonList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Age` | pxInteger | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya | ALWAYS |
| `.ASMGender` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMHeight` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | pxInteger | 1 | ya | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMParticipantStatus` | pxInteger | 1 | ya | ALWAYS |
| `.ASMWeight` | pxInteger | 1 | ya | ALWAYS |
| `.BrandName` | pxDisplayText | 1 | ya | ALWAYS |
| `.ChassisNumber` | (tanpa) | 1 | ya | ALWAYS |
| `.ConveyanceNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.EngineNumber` | (tanpa) | 1 | ya | ALWAYS |
| `.GoodNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.LicensePlate` | pxDropdown | 1 | ya | ALWAYS |
| `.ModelName` | pxDisplayText | 1 | ya | ALWAYS |
| `.PackingNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxDisplayText | 2 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText | 2 | ya | ALWAYS |
| `.pxListSubscript` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyFullName` | pxDisplayText | 3 | ya | ALWAYS |
| `.TradingNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.TypeName` | pxDisplayText | 1 | ya | ALWAYS |

#### CoverageCommisionList_IsUW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 24 (unik 21) - read-only 6 - RepeatGrid 6 - grid: .CargoList, .LocationList, .PersonList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Age` | pxInteger | 1 | tidak | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 1 | ya | ALWAYS |
| `.ASMGender` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMHeight` | pxInteger | 1 | tidak | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMWeight` | pxInteger | 1 | tidak | ALWAYS |
| `.BrandName` | (tanpa) | 1 | tidak | ALWAYS |
| `.ChassisNumber` | (tanpa) | 1 | tidak | ALWAYS |
| `.ConveyanceNote` | pxTextInput | 1 | tidak | ALWAYS |
| `.EngineNumber` | (tanpa) | 1 | tidak | ALWAYS |
| `.GoodNote` | pxTextInput | 1 | tidak | ALWAYS |
| `.LicensePlate` | pxDropdown | 1 | tidak | ALWAYS |
| `.ModelName` | pxInteger | 1 | tidak | ALWAYS |
| `.PackingNote` | pxTextInput | 1 | tidak | ALWAYS |
| `.Property.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxDisplayText | 2 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText/pxTextInput | 2 | ya | ALWAYS |
| `.pxListSubscript` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyFullName` | (tanpa) | 2 | tidak | ALWAYS |
| `.TradingNote` | pxTextInput | 1 | tidak | ALWAYS |
| `.TypeName` | pxInteger | 1 | tidak | ALWAYS |

#### CoverageFormula

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.IsShowFormula` | pxCheckbox | 1 | tidak | ALWAYS |

#### CoverageItem

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 48 (unik 42) - read-only 44 - RepeatGrid 3 - grid: .DeductibleList, .LayerList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AccumulationCode` | pxTextInput | 1 | ya | ALWAYS |
| `.AccumulationDescription` | pxDisplayText | 1 | ya | ALWAYS |
| `.Amount` | pxNumber | 1 | ya | ALWAYS |
| `.ASMPremium` | pxNumber | 1 | ya | ALWAYS |
| `.ASMTSI` | pxNumber | 1 | ya | ALWAYS |
| `.Condition` | pxDropdown | 1 | ya | .Condition != 5 |
| `.Conditions` | pxTextArea | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.ConstructionEarthquake` | pxAutoComplete | 1 | ya [.FlagDelete = 1] | .Coverage=='100821' \|\| .Coverage=='100829' |
| `.CoverageBasis` | pxDropdown | 1 | ya [.FlagDelete = 1] | .CoverageBasis==2 |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.Day` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.Descriptions` | pxDropdown | 1 | ya | ALWAYS |
| `.Discount` | pxTextInput | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.DiscountPercentage` | pxNumber | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.EmlPml` | pxNumber | 1 | ya [.FlagDelete = 1] | .CoverageBasis==3 |
| `.FirstLoss` | pxNumber | 1 | ya [.FlagDelete = 1] | .CoverageBasis==2 |
| `.FirstScale` | pxNumber | 1 | ya | .CoverageBasis==2 |
| `.Indemnity` | pxNumber | 1 | ya [.FlagDelete = 1] | .CoverageBasis==2 |
| `.IndemnityPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.InputCondition` | pxTextInput | 1 | ya | .Condition = 5 |
| `.Layer` | pxNumber | 1 | ya [.FlagDelete = 1] | .CoverageBasis == 5 |
| `.LayerNo` | (tanpa) | 1 | tidak | ALWAYS |
| `.LimitofLiability` | pxNumber | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.LostLimit` | pxNumber | 1 | ya [.FlagDelete = 1] | .CoverageBasis==4 \|\| .CoverageBasis==5 |
| `.MinMax` | pxDropdown | 1 | ya | ALWAYS |
| `.NetRate` | pxNumber | 1 | ya [.FlagDelete = 1] | 1=2 |
| `.OLDID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PctDeductible` | pxNumber | 1 | ya | ALWAYS |
| `.PctDeductible2` | pxNumber | 1 | ya | .MinMax = 3 |
| `.PctLoL` | pxNumber | 1 | tidak [.FlagDelete = 1] | ALWAYS |
| `.Premium` | pxNumber/pxTextInput | 2 | ya [.FlagDelete = 1] | ALWAYS |
| `.ProRatePercent` | pxNumber | 1 | ya [.FlagDelete = 1 \|\| pyWorkPage.OfferFacIn.IsPro...] | .ProRatePercent>0 |
| `.Rate` | pxNumber/pxTextInput | 2 | ya [.FlagDelete = 1] | ALWAYS |
| `.Sublimit` | pxNumber | 1 | ya [.FlagDelete = 1] | .CoverageBasis==4 \|\| .CoverageBasis==5 |
| `.TimeExcess` | pxNumber | 2 | ya | ALWAYS |
| `.TSI` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.TSILiability` | pxNumber/pxTextInput | 2 | campur | .CoverageBasis==2 |
| `.TypeDeductible` | pxDropdown | 1 | ya | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | ya | .MinMax = 3 |
| `.Unit` | pxDropdown | 1 | ya [.FlagDelete = 1] | .CoverageBasis==2 |
| `.Zone` | pxAutoComplete/pxTextArea | 2 | ya [.FlagDelete = 1] | (.Coverage=='100835' \|\| .Coverage=='100821' \|\| .Coverage=='100829'... |
| `pyWorkPage.OfferFacIn.ProRatePercent` | pxNumber | 1 | ya [.FlagDelete = 1 \|\| pyWorkPage.OfferFacIn.IsPro...] | .ProRatePercent<1 |

#### CoverageItem_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 44 (unik 38) - read-only 43 - RepeatGrid 3 - grid: .DeductibleList, .LayerList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AccumulationCode` | pxTextInput | 1 | ya | ALWAYS |
| `.AccumulationDescription` | pxDisplayText | 1 | ya | ALWAYS |
| `.Amount` | pxNumber | 1 | ya | ALWAYS |
| `.Condition` | pxDropdown | 1 | ya | .Condition != 5 |
| `.Conditions` | pxTextArea | 1 | ya | ALWAYS |
| `.ConstructionEarthquake` | pxAutoComplete | 1 | ya | .Coverage=='100821' \|\| .Coverage=='100829' |
| `.CoverageBasis` | pxDropdown | 1 | ya | .CoverageBasis==2 |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.Descriptions` | pxDropdown | 1 | ya | ALWAYS |
| `.Discount` | pxTextInput | 1 | ya | ALWAYS |
| `.DiscountPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.EmlPml` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsUW \|\| pyPortal.I...] | .CoverageBasis==3 |
| `.FirstLoss` | pxNumber | 1 | ya | .CoverageBasis==2 |
| `.FirstScale` | pxNumber | 1 | ya | .CoverageBasis==2 |
| `.Indemnity` | pxNumber | 1 | ya | .CoverageBasis==2 |
| `.IndemnityPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.InputCondition` | pxTextInput | 1 | ya | .Condition = 5 |
| `.Layer` | pxNumber | 1 | ya | .CoverageBasis == 5 |
| `.LayerNo` | (tanpa) | 1 | tidak | ALWAYS |
| `.LimitofLiability` | pxNumber | 1 | ya | ALWAYS |
| `.LostLimit` | pxNumber | 1 | ya | .CoverageBasis==4 \|\| .CoverageBasis==5 |
| `.MinMax` | pxDropdown | 1 | ya | ALWAYS |
| `.NetRate` | pxNumber | 1 | ya | NOTBLANK |
| `.OLDID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PctDeductible` | pxNumber | 1 | ya | ALWAYS |
| `.PctDeductible2` | pxNumber | 1 | ya | .MinMax = 3 |
| `.PctLoL` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.Rate` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.Sublimit` | pxNumber | 1 | ya | .CoverageBasis==4 \|\| .CoverageBasis==5 |
| `.TimeExcess` | pxNumber | 2 | ya | ALWAYS |
| `.TSI` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.TSILiability` | pxTextInput | 2 | ya | .CoverageBasis==2 |
| `.TypeDeductible` | pxDropdown | 1 | ya | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | ya | .MinMax = 3 |
| `.Unit` | pxDropdown | 1 | ya | .CoverageBasis==2 |
| `.Zone` | pxAutoComplete/pxTextArea | 2 | ya | (.Coverage=='100835' \|\| .Coverage=='100821' \|\| .Coverage=='100829'... |
| `pyWorkPage.OfferFacIn.ProRatePercent` | pxNumber | 1 | ya | ALWAYS |

#### CoverageList

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: .LocationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText | 1 | ya | ALWAYS |

#### CoverageList_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: .LocationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText | 1 | ya | ALWAYS |

#### CoveragePropertyFacOut

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | tidak | ALWAYS |

#### CoverageSpreadingList

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 28 (unik 23) - read-only 16 - RepeatGrid 7 - grid: .CargoList, .LocationList, .PersonList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Age` | pxInteger | 1 | tidak | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya | ALWAYS |
| `.ASMGender` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMHeight` | pxInteger | 1 | tidak | ALWAYS |
| `.ASMIDCard` | (tanpa) | 1 | tidak | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMParticipantStatus` | (tanpa) | 1 | tidak | ALWAYS |
| `.ASMWeight` | pxInteger | 1 | tidak | ALWAYS |
| `.BrandName` | pxDisplayText | 1 | ya | ALWAYS |
| `.ChassisNumber` | (tanpa) | 1 | tidak | ALWAYS |
| `.ConveyanceNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.EngineNumber` | (tanpa) | 1 | tidak | ALWAYS |
| `.GoodNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.LicensePlate` | pxDropdown | 1 | tidak | ALWAYS |
| `.ModelName` | pxDisplayText | 1 | ya | ALWAYS |
| `.PackingNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxDisplayText | 2 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText | 2 | ya | NOTBLANK |
| `.pxListSubscript` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyFullName` | pxDisplayText | 3 | ya | ALWAYS |
| `.TradingNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.TypeName` | pxInteger | 1 | tidak | ALWAYS |

#### CoverageSpreadingList_IsUW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 28 (unik 23) - read-only 24 - RepeatGrid 7 - grid: .CargoList, .LocationList, .PersonList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Age` | pxInteger | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime/pxInteger | 2 | campur | ALWAYS |
| `.ASMGender` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMHeight` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | (tanpa) | 1 | tidak | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMParticipantStatus` | (tanpa) | 1 | tidak | ALWAYS |
| `.ASMWeight` | pxInteger | 1 | ya | ALWAYS |
| `.BrandName` | pxDisplayText | 1 | ya | ALWAYS |
| `.ChassisNumber` | (tanpa) | 1 | ya | ALWAYS |
| `.ConveyanceNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.EngineNumber` | (tanpa) | 1 | ya | ALWAYS |
| `.GoodNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.LicensePlate` | pxDisplayText | 1 | ya | ALWAYS |
| `.ModelName` | pxDisplayText | 1 | ya | ALWAYS |
| `.PackingNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxDisplayText | 2 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText | 2 | ya | ALWAYS |
| `.pxListSubscript` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyFullName` | pxDisplayText | 3 | ya | ALWAYS |
| `.TradingNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.TypeName` | pxDisplayText | 1 | ya | ALWAYS |

#### CoverageSummary

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: D_CoverageSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.OLDID` | (tanpa) | 1 | ya | ALWAYS |
| `.Rate` | (tanpa) | 1 | ya | ALWAYS |

#### CreateListClause

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### CreateListWarranty

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work-NB` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.pySearchCase` | pxTextInput | 1 | tidak | ALWAYS |
| `.pySearchString` | pxTextInput | 1 | tidak | ALWAYS |

#### DeductibleDetail

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DeductibleType` | pxTextInput | 1 | ya | ALWAYS |
| `.Descriptions` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationID` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.TimeExcess` | pxTextInput | 1 | ya | ALWAYS |
| `.Value` | pxTextInput | 1 | ya | ALWAYS |

#### deleteItemButton

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### DetailCoverageListFire_GCNM

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 3 (unik 3) - read-only 2 - RepeatGrid 1 - grid: .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 1 | ya | ALWAYS |
| `.DeductibleNote` | (tanpa) | 1 | tidak | ALWAYS |
| `.TSISublimit` | pxTextInput | 1 | ya | ALWAYS |

#### DetailCoverageTravelGrid_GCNM

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 1 - grid: .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Coverage` | pxTextInput | 1 | ya | ALWAYS |
| `.CoverageNote` | pxDropdown | 1 | ya | ALWAYS |
| `.DeductibleNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Discount` | pxTextInput | 1 | ya | ALWAYS |
| `.DiscountPercentage` | pxTextInput | 1 | ya | ALWAYS |
| `.Premium` | pxTextInput | 1 | ya | ALWAYS |
| `.TempJaminan` | pxTextInput | 1 | ya | ALWAYS |
| `.TSI` | pxTextInput | 1 | ya | ALWAYS |

#### DetailDeptHeadTreatyIn_UW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-PolicyTreatyIn` - field terikat 87 (unik 76) - read-only 84 - RepeatGrid 3 - grid: .BreakDownSpreadList, .ListInstallment, .SpreadingRiskList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BalanceBeforePPH` | pxTextInput | 1 | ya | .BalanceDueTo >= 0 |
| `.BalanceBeforeTax` | pxTextInput | 1 | ya | .BalanceDueTo >= 0 |
| `.BalanceDueTo` | pxTextInput | 2 | ya | .BalanceDueTo < 0 ; .BalanceDueTo >=0 |
| `.BizName` | pxDisplayText | 1 | ya | .IsNewPolicyNonProp != 1 && NEVER |
| `.CedingCoName` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Claim` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.ClaimPaymentType` | pxDropdown | 1 | ya | ALWAYS |
| `.ClaimPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.ClaimType` | pxDropdown | 1 | ya | ALWAYS |
| `.Currency` | pxAutoComplete | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.Deduction1` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.Deduction2` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.DueDate` | pxDateTime | 1 | ya | ALWAYS |
| `.DueTo` | pxRadioButtons | 1 | tidak | 1=2 |
| `.EndDate` | pxDateTime | 1 | ya [1==1] | ALWAYS |
| `.ExcessLoss` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.FlagPPH` | pxCheckbox | 1 | tidak | ALWAYS |
| `.GrossPremium` | pxNumber | 1 | ya | ALWAYS |
| `.Installment` | pxTextInput | 1 | ya | ALWAYS |
| `.InstallmentNo` | pxInteger | 1 | ya | ALWAYS |
| `.InstallmentPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.Layer` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerPart` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerPartType` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerType` | pxTextInput | 1 | ya | ALWAYS |
| `.MarketingOfficer` | pxDisplayText | 1 | ya | NOTBLANK |
| `.NetPremium` | pxCurrency | 1 | ya | ALWAYS |
| `.NoOffer` | pxDisplayText | 1 | ya | ALWAYS |
| `.OutstandingClaim` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.OveriddingCommOgp` | pxNumber | 1 | ya [1==1] | ALWAYS |
| `.OveriddingCommOnp` | pxNumber | 1 | ya [1==1] | ALWAYS |
| `.PaymentTotal` | pxNumber | 1 | ya | ALWAYS |
| `.PPHValue` | pxCurrency | 1 | ya | ALWAYS |
| `.PPNValue` | pxCurrency | 1 | ya | ALWAYS |
| `.PremiOgp` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.PremiOnp` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.Quartal` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | 1 | ya [1=1] | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.QuotationData.NoOfferSlip` | pxTextArea | 1 | tidak | ALWAYS |
| `.QuotationData.ProportionalType` | pxDropdown | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.Remark` | pxTextArea | 1 | ya | ALWAYS |
| `.ResultOgp1` | pxCurrency | 1 | ya [1 = 1] | 1==2 |
| `.ResultOgp2` | pxCurrency | 1 | ya | ALWAYS |
| `.ResultOnp1` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.ResultOnp2` | pxCurrency | 1 | ya | ALWAYS |
| `.RiCommOgp` | pxNumber | 1 | ya [1==1] | ALWAYS |
| `.RiCommOnp` | pxNumber | 1 | ya [1==1] | ALWAYS |
| `.SalvageValue` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.ShareCurrency` | pxTextInput | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | ya | ALWAYS |
| `.ShareValue` | pxCurrency | 1 | ya | ALWAYS |
| `.SOBName` | pxDisplayText | 1 | ya | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.StartDate` | pxDateTime | 1 | ya [1==1] | ALWAYS |
| `.StatementDate` | pxDateTime | 1 | ya [1==1] | ALWAYS |
| `.StatementType` | pxDropdown | 1 | ya | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.TotalClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalPremium` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentageClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentagePremium` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyGroupName` | pxDisplayText | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.TreatyType` | pxDropdown/pxTextInput | 3 | ya | ALWAYS |
| `.TreatyYear` | pxTextInput | 2 | ya | .IsNewPolicyNonProp != 1 |
| `.TypeTax` | pxDisplayText | 1 | ya | NOTBLANK |
| `.YearOfQuartal` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.PolicyNo` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.ProductionDate` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.Commencement` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.Termination` | pxDateTime | 1 | ya | ALWAYS |

#### DetailLocation

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### DetailPoliciesNonProportional

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### DetailPolicyTreatyIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-PolicyTreatyIn` - field terikat 87 (unik 76) - read-only 49 - RepeatGrid 3 - grid: .BreakDownSpreadList, .ListInstallment, .SpreadingRiskList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BalanceBeforePPH` | pxTextInput | 1 | ya | .BalanceDueTo >= 0 |
| `.BalanceBeforeTax` | pxTextInput | 1 | ya | .BalanceDueTo >= 0 |
| `.BalanceDueTo` | pxTextInput | 2 | ya | .BalanceDueTo < 0 ; .BalanceDueTo >= 0 |
| `.BizName` | pxTextInput | 1 | ya | 1=2 |
| `.CedingCoName` | pxTextInput | 1 | ya | ALWAYS |
| `.Claim` | pxCurrency | 1 | tidak | ALWAYS |
| `.ClaimPaymentType` | pxDropdown | 1 | tidak | ALWAYS |
| `.ClaimPercentage` | pxNumber | 1 | tidak | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.ClaimType` | pxDropdown | 1 | tidak | ALWAYS |
| `.Deduction1` | pxCurrency | 1 | tidak | ALWAYS |
| `.Deduction2` | pxCurrency | 1 | tidak | ALWAYS |
| `.DueDate` | pxDateTime | 1 | ya [IsUW \|\| IsClaim \|\| pyPortal.IsShowOpenPolicy...] | ALWAYS |
| `.DueTo` | pxRadioButtons | 1 | tidak | NEVER |
| `.EndDate` | pxDateTime | 1 | tidak | ALWAYS |
| `.ExcessLoss` | pxCurrency | 1 | tidak | ALWAYS |
| `.FlagPPH` | pxCheckbox | 1 | tidak | ALWAYS |
| `.FlagRetroTreaty` | pxCheckbox | 1 | tidak | .ClaimType != 'XOL Retro' |
| `.GrossClaim` | pxNumber | 1 | tidak | ALWAYS |
| `.GrossPremium` | pxNumber | 1 | tidak | ALWAYS |
| `.IDCurrency` | pxDropdown | 1 | tidak | .IsNewPolicyNonProp != 1 |
| `.Installment` | pxTextInput | 1 | tidak | ALWAYS |
| `.InstallmentNo` | pxInteger | 1 | ya | ALWAYS |
| `.InstallmentPercentage` | pxNumber | 1 | ya [IsUW \|\| pyPortal.IsShowOpenPolicyMarine == tru...] | ALWAYS |
| `.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.Layer` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerPart` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerPartType` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerType` | pxTextInput | 1 | ya | ALWAYS |
| `.NetPremium` | pxCurrency | 1 | ya | ALWAYS |
| `.NoOffer` | pxDisplayText | 1 | ya | ALWAYS |
| `.OutstandingClaim` | pxCurrency | 1 | tidak | ALWAYS |
| `.OveriddingCommOgp` | pxNumber | 1 | tidak | ALWAYS |
| `.OveriddingCommOnp` | pxNumber | 1 | tidak | ALWAYS |
| `.PaymentTotal` | pxNumber | 1 | ya | ALWAYS |
| `.PPHValue` | pxCurrency | 1 | ya | ALWAYS |
| `.PPNValue` | pxCurrency | 1 | ya | ALWAYS |
| `.PremiOgp` | pxCurrency | 1 | tidak | ALWAYS |
| `.PremiOnp` | pxCurrency | 1 | tidak | ALWAYS |
| `.Premium` | pxNumber | 1 | ya [IsUW \|\| pyPortal.IsShowOpenPolicyMarine == tru...] | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.Quartal` | pxTextInput | 1 | tidak | ALWAYS |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | 1 | tidak [IsUW] | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.QuotationData.MOID` | pxDropdown | 1 | tidak | ALWAYS |
| `.QuotationData.NoOfferSlip` | pxTextArea | 1 | tidak | ALWAYS |
| `.QuotationData.ProportionalType` | pxDropdown | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.Remark` | pxTextArea | 1 | tidak | ALWAYS |
| `.ResultOgp1` | pxCurrency | 1 | tidak | ALWAYS |
| `.ResultOgp2` | pxCurrency | 1 | tidak | ALWAYS |
| `.ResultOnp1` | pxCurrency | 1 | tidak | ALWAYS |
| `.ResultOnp2` | pxCurrency | 1 | tidak | ALWAYS |
| `.RiCommOgp` | pxNumber | 1 | tidak | ALWAYS |
| `.RiCommOnp` | pxNumber | 1 | tidak | ALWAYS |
| `.SalvageValue` | pxCurrency | 1 | tidak | ALWAYS |
| `.ShareCurrency` | pxTextInput | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | campur | ALWAYS |
| `.ShareValue` | pxCurrency | 1 | ya | ALWAYS |
| `.StartDate` | pxDateTime | 1 | tidak | ALWAYS |
| `.StatementDate` | pxDateTime | 1 | ya [ALWAYS] | ALWAYS |
| `.StatementType` | pxDropdown | 1 | tidak [IsUW] | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.TotalClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalPremium` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentageClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentagePremium` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyGroupName` | pxTextInput | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.TreatyType` | pxDropdown/pxTextInput | 3 | campur | ALWAYS |
| `.TreatyYear` | pxTextInput | 2 | ya | .IsNewPolicyNonProp != 1 |
| `.TypeTax` | pxRadioButtons | 1 | tidak | .FlagPPH = true |
| `.YearOfQuartal` | pxTextInput | 1 | tidak | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.SOBName` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.Commencement` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.Termination` | pxDateTime | 1 | ya | ALWAYS |

#### DetailPolicyTreatyInNonProportional

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-LimitSummaryList` - field terikat 64 (unik 23) - read-only 18 - RepeatGrid 17 - grid: pyWorkPage.PolicyTreatyIn.ListInstallment, pyWorkPage.TreatyIn.LimitFacShareSummaryList, pyWorkPage.TreatyIn.LimitShareSummaryList, pyWorkPage.TreatyIn.LimitSummaryList, pyWorkPage.TreatyIn.TotalFacShareDeductionNP, pyWorkPage.TreatyIn.TotalFacShareGrossNP, pyWorkPage.TreatyIn.TotalFacShareNetNP, pyWorkPage.TreatyIn.TotalFacShareRnmNP, pyWorkPage.TreatyIn.TotalLimitDeductblNP, pyWorkPage.TreatyIn.TotalLimitIOONP, pyWorkPage.TreatyIn.TotalLimitMDPNP, pyWorkPage.TreatyIn.TotalShareDeductionNP, pyWorkPage.TreatyIn.TotalShareGrossNP, pyWorkPage.TreatyIn.TotalShareNetNP, pyWorkPage.TreatyIn.TotalShareRnmNP, pyWorkPage.TreatyIn.TotalSpreadedNetPremi, pyWorkPage.TreatyIn.TotalSpreadedNetPremiRI

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxNumber | 14 | tidak | ALWAYS |
| `.Deductible` | pxNumber | 3 | campur | ALWAYS |
| `.Deductible2` | pxNumber | 3 | campur | ALWAYS |
| `.Limit` | pxNumber | 3 | campur | ALWAYS |
| `.Limit2` | pxNumber | 3 | campur | ALWAYS |
| `.MDP` | pxNumber | 3 | campur | ALWAYS |
| `.MDP2` | pxNumber | 3 | campur | ALWAYS |
| `.NetPremi` | pxNumber | 2 | campur | ALWAYS |
| `.NetPremi2` | pxNumber | 2 | campur | ALWAYS |
| `.NetPremiAfterPPH` | pxNumber | 1 | tidak | ALWAYS |
| `.NetPremiAfterPPH2` | pxNumber | 1 | tidak | ALWAYS |
| `.NetPremiAfterPPN` | pxNumber | 1 | tidak | ALWAYS |
| `.NetPremiAfterPPN2` | pxNumber | 1 | tidak | ALWAYS |
| `.Note` | (tanpa) | 3 | tidak | ALWAYS |
| `.TotalNetPremiAfterPPN` | pxNumber | 1 | tidak | ALWAYS |
| `.TotalNetPremiAfterTax` | pxNumber | 1 | tidak | ALWAYS |
| `.TotalPPHValue` | pxNumber | 1 | tidak | ALWAYS |
| `.TotalPPNValue` | pxNumber | 1 | tidak | ALWAYS |
| `.Value` | pxNumber | 13 | tidak | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.Installment` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.FacultativeShare` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.RNMShare` | pxTextInput | 1 | ya | pyWorkPage.TreatyIn.FacultativeShare = 0 |
| `pyWorkPage.TreatyIn.RnmShareDeducted` | pxTextInput | 1 | ya | pyWorkPage.TreatyIn.FacultativeShare != 0 |

#### DetailPolicyTreatyInNonProportionalEDM

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-LimitSummaryList` - field terikat 29 (unik 11) - read-only 2 - RepeatGrid 10 - grid: pyWorkPage.TreatyIn.Installment, pyWorkPage.TreatyIn.ValueDifference.LimitShareSummaryList, pyWorkPage.TreatyIn.ValueDifference.LimitSummaryList, pyWorkPage.TreatyIn.ValueDifference.TotalLimitDeductblNP, pyWorkPage.TreatyIn.ValueDifference.TotalLimitMDPNP, pyWorkPage.TreatyIn.ValueDifference.TotalShareDeductionNP, pyWorkPage.TreatyIn.ValueDifference.TotalShareGrossNP, pyWorkPage.TreatyIn.ValueDifference.TotalShareNetNP, pyWorkPage.TreatyIn.ValueDifference.TotalSpreadedNetPremi, pyWorkPage.TreatyIn.ValueDifference.TotalSpreadedNetPremiRI

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxNumber | 8 | tidak | ALWAYS |
| `.Deductible` | (tanpa) | 2 | tidak | ALWAYS |
| `.Deductible2` | (tanpa) | 2 | tidak | ALWAYS |
| `.MDP` | (tanpa) | 2 | tidak | ALWAYS |
| `.MDP2` | (tanpa) | 2 | tidak | ALWAYS |
| `.NetPremi` | (tanpa) | 1 | tidak | ALWAYS |
| `.NetPremi2` | (tanpa) | 1 | tidak | ALWAYS |
| `.Note` | (tanpa) | 2 | tidak | ALWAYS |
| `.Value` | pxNumber | 7 | tidak | ALWAYS |
| `pyWorkPage.TreatyIn.InstallmentNo` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.RNMShare` | pxTextInput | 1 | ya | ALWAYS |

#### DetailPolicyTreatyOutNonProportional

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-LimitSummaryList` - field terikat 34 (unik 14) - read-only 17 - RepeatGrid 10 - grid: pyWorkPage.PolicyTreatyIn.ListInstallment, pyWorkPage.TreatyIn.LimitShareSummaryList, pyWorkPage.TreatyIn.LimitSummaryList, pyWorkPage.TreatyIn.TotalLimitDeductblNP, pyWorkPage.TreatyIn.TotalLimitIOONP, pyWorkPage.TreatyIn.TotalLimitMDPNP, pyWorkPage.TreatyIn.TotalShareDeductionNP, pyWorkPage.TreatyIn.TotalShareGrossNP, pyWorkPage.TreatyIn.TotalShareNetNP, pyWorkPage.TreatyIn.TotalShareRnmNP

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxNumber | 8 | tidak | ALWAYS |
| `.Deductible` | pxNumber | 2 | ya | ALWAYS |
| `.Deductible2` | pxNumber | 2 | ya | ALWAYS |
| `.Limit` | pxNumber | 2 | ya | ALWAYS |
| `.Limit2` | pxNumber | 2 | ya | ALWAYS |
| `.MDP` | pxNumber | 2 | ya | ALWAYS |
| `.MDP2` | pxNumber | 2 | ya | ALWAYS |
| `.NetPremi` | pxNumber | 1 | ya | ALWAYS |
| `.NetPremi2` | pxNumber | 1 | ya | ALWAYS |
| `.Note` | (tanpa) | 2 | tidak | ALWAYS |
| `.Value` | pxNumber | 7 | tidak | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.Installment` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.ReinsurerShare` | pxTextInput | 1 | ya | pyWorkPage.TreatyIn.FacultativeShare = 0 |
| `pyWorkPage.TreatyIn.RnmShareDeducted` | pxTextInput | 1 | ya | pyWorkPage.TreatyIn.FacultativeShare != 0 |

#### EmailSection

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-SuggestList` - field terikat 10 (unik 10) - read-only 6 - RepeatGrid 1 - grid: TempViewSuggest.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Approval` | pxDropdown | 1 | ya | ALWAYS |
| `.BandingTo` | pxDropdown | 1 | tidak | .EmailType = 4 |
| `.CommentSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.DateTransfer` | pxDateTime | 1 | ya | ALWAYS |
| `.EmailType` | pxRadioButtons | 1 | tidak | .FlagOnGoingPolicy != 1 && pyWorkPage.PositionNote != 'ReasFacInTeamLe... |
| `.EmailTypeTL` | pxRadioButtons | 1 | tidak | pyWorkPage.PositionNote == 'ReasFacInTeamLeader' |
| `.EmailTypeUWPolicy` | pxRadioButtons | 1 | tidak | .FlagOnGoingPolicy == 1 && pyWorkPage.PositionNote != 'ReasFacInTeamLe... |
| `.IsCedingConfirm` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.Comment` | pxTextArea | 1 | ya [.OfferFacIn.QuotationData.BusinessType = 'Life'] | .OfferFacIn.QuotationData.StatusBusiness!='3' |
| `.PICSuggest` | pxDisplayText | 1 | ya | ALWAYS |

#### EmailSection_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 16 (unik 16) - read-only 5 - RepeatGrid 2 - grid: pyWorkPage.QuotationList, TempViewSuggest.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Approval` | pxDropdown | 1 | ya | ALWAYS |
| `.CommentSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.DateTransfer` | pxDateTime | 1 | ya | ALWAYS |
| `.Email` | pxTextInput | 1 | tidak | ALWAYS |
| `.EmailText` | pxRichTextEditor | 1 | tidak | 1=2 |
| `.EmailTypeUW` | pxRadioButtons | 1 | tidak | pyWorkPage.OfferFacIn.QuotationData.StatusBusiness!=3&& .FlagOnGoingPo... |
| `.EmailTypeUWEDM` | pxRadioButtons | 1 | tidak | pyWorkPage.OfferFacIn.QuotationData.StatusBusiness!=3&& .FlagOnGoingPo... |
| `.EmailTypeUWPolicy` | pxRadioButtons | 1 | tidak | 1=2 |
| `.IsCedingConfirm` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.Comment` | pxTextArea | 1 | tidak | ALWAYS |
| `.OfferFacIn.DateValidity` | pxDateTime | 1 | tidak | ALWAYS |
| `.OfferFacIn.DaysValidity` | pxTextInput | 1 | tidak | ALWAYS |
| `.OfferFacIn.IsBanding` | pxCheckbox | 1 | tidak | .EmailType = 1 \|\| .EmailTypeUWPolicy = 1 |
| `.OfferFacIn.PolicyData.ProdDateTime` | pxDateTime | 1 | tidak | ProtectSpreading.CARI31==0 |
| `.OfferFacIn.WaitingBindDate` | pxDateTime | 1 | tidak | ALWAYS |
| `.PICSuggest` | pxDisplayText | 1 | ya | ALWAYS |

#### EmailSection_Life

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-SuggestList` - field terikat 7 (unik 7) - read-only 6 - RepeatGrid 1 - grid: TempViewSuggest.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Approval` | pxDropdown | 1 | ya | ALWAYS |
| `.CommentSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.DateTransfer` | pxDateTime | 1 | ya | ALWAYS |
| `.EmailTypeLife` | pxRadioButtons | 1 | tidak | .FlagOnGoingPolicy != 1 \|\| pyWorkPage.PositionNote == 'ReasFacInTeam... |
| `.IsCedingConfirm` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.Comment` | pxTextArea | 1 | ya [.OfferFacIn.QuotationData.BusinessType != 'Life'] | .OfferFacIn.QuotationData.StatusBusiness!='3' |
| `.PICSuggest` | pxDisplayText | 1 | ya | ALWAYS |

#### EmailSectionCeding

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-SuggestList` - field terikat 10 (unik 10) - read-only 6 - RepeatGrid 1 - grid: TempViewSuggest.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Approval` | pxDropdown | 1 | ya | ALWAYS |
| `.BandingTo` | pxDropdown | 1 | tidak | .EmailTypeBinding= 4 |
| `.CommentSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.DateTransfer` | pxDateTime | 1 | ya | ALWAYS |
| `.EmailTypeBinding` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.EmailTypeCeding` | pxRadioButtons | 1 | tidak | 1 = 2 |
| `.IsCedingConfirm` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.CommentCeding` | pxTextArea | 1 | ya [.OfferFacIn.QuotationData.BusinessType = 'Life'] | ALWAYS |
| `.OfferFacIn.ReceivedRiSlip` | pxCheckbox | 1 | tidak | ALWAYS |
| `.PICSuggest` | pxDisplayText | 1 | ya | ALWAYS |

#### FacOutOffer

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### FacOutPrintRISlipSectionInside

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 16 (unik 15) - read-only 16 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoinsList(1).TSIShare` | pxNumber | 1 | ya | ALWAYS |
| `.OfferFacIn.FacRetroDetails.AdditionalInfo` | pxTextArea | 1 | ya [IsRISlip] | 1=2 |
| `.OfferFacIn.FacRetroDetails.BackUpStatus` | pxDropdown | 1 | ya [1=1] | 1=2 |
| `.OfferFacIn.FacRetroDetails.DocumentPosition` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.OfficerFacOut` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.PercentShare` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.PolicyData.EndDateTime` | pxDateTime | 1 | ya [.pyID != ""] | ALWAYS |
| `.OfferFacIn.PolicyData.PolicyNo` | pxTextInput | 1 | ya | NOTBLANK |
| `.OfferFacIn.PolicyData.StartDateTime` | pxDateTime | 1 | ya [.pyID != ""] | ALWAYS |
| `.OfferFacIn.QuotationData.BusinessName` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.QuotationData.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.pyID` | pxTextInput | 1 | ya | ALWAYS |
| `FacOfferList.FacNo` | pxTextInput | 1 | ya | ALWAYS |
| `InputData.CARI11` | pxNumber | 1 | ya | IsFire |
| `InputData.CARI12` | pxTextInput | 2 | ya | ALWAYS |

#### FacultativeLetter

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-FacOffer` - field terikat 5 (unik 5) - read-only 3 - RepeatGrid 1 - grid: .OfferFacIn.FacRetroList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Attention` | (tanpa) | 1 | tidak | ALWAYS |
| `.OfferFacIn.FacRetroDetails.BackUpStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.OfferFacIn.FacRetroDetails.DocumentPosition` | pxTextInput | 1 | ya | ALWAYS |
| `.OurRef` | (tanpa) | 1 | ya | ALWAYS |
| `.ReinsurerName` | (tanpa) | 1 | tidak | ALWAYS |

#### FEAList

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 1 - grid: .FEAList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.APAR` | pxDisplayText | 1 | ya | ALWAYS |
| `.DataFEA.PrivateTruckBrigade` | pxDisplayText | 1 | ya | ALWAYS |
| `.Hydrant` | pxDisplayText | 1 | ya | ALWAYS |
| `.InfoFEA` | pxDisplayText | 1 | ya | ALWAYS |
| `.SmokeDetector` | pxDisplayText | 1 | ya | ALWAYS |
| `.Sprinkler` | pxDisplayText | 1 | ya | ALWAYS |

#### FEAList_IsUW

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 1 - grid: .FEAList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.APAR` | pxDisplayText | 1 | ya | ALWAYS |
| `.DataFEA.PrivateTruckBrigade` | pxDisplayText | 1 | ya | ALWAYS |
| `.Hydrant` | pxDisplayText | 1 | ya | ALWAYS |
| `.InfoFEA` | pxDisplayText | 1 | ya | ALWAYS |
| `.SmokeDetector` | pxDisplayText | 1 | ya | ALWAYS |
| `.Sprinkler` | pxDisplayText | 1 | ya | ALWAYS |

#### FinalPolicyRO

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### FireSummarySection

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: D_LocationSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.ObjectName` | (tanpa) | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | (tanpa) | 1 | ya | ALWAYS |

#### FormulaDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AdditionalCapital` | pxNumber | 1 | ya [IsUW] | ALWAYS |
| `.MaxPctTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.MaxTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.ShareInTSI` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | 1 = 2 |

#### FormulaTreatyCapacityDesc

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 8 (unik 8) - read-only 7 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.InwardScale.PctTreatyLimit` | pxNumber | 1 | ya | ALWAYS |
| `.IsGrossNetShow` | pxCheckbox | 1 | tidak | ALWAYS |
| `.MaxPctTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.MaxTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.Parameters.LowestPctLimit` | pxNumber | 1 | ya | ALWAYS |
| `.ShareInTSI` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |

#### FormulaTreatyCapacityDesc_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 8 (unik 8) - read-only 7 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.InwardScale.PctTreatyLimit` | pxNumber | 1 | ya | ALWAYS |
| `.IsGrossNetShow` | pxCheckbox | 1 | tidak | ALWAYS |
| `.MaxPctTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.MaxTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.Parameters.LowestPctLimit` | pxNumber | 1 | ya | ALWAYS |
| `.ShareInTSI` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |

#### GeneralDeptHeadTreatyIn_UW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-PolicyTreatyIn` - field terikat 87 (unik 76) - read-only 84 - RepeatGrid 3 - grid: .BreakDownSpreadList, .ListInstallment, .SpreadingRiskList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BalanceBeforePPH` | pxTextInput | 1 | ya | .BalanceDueTo >= 0 |
| `.BalanceBeforeTax` | pxTextInput | 1 | ya | .BalanceDueTo >= 0 |
| `.BalanceDueTo` | pxTextInput | 2 | ya | .BalanceDueTo < 0 ; .BalanceDueTo >=0 |
| `.BizName` | pxDisplayText | 1 | ya | .IsNewPolicyNonProp != 1 && NEVER |
| `.CedingCoName` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Claim` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.ClaimPaymentType` | pxDropdown | 1 | ya | ALWAYS |
| `.ClaimPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.ClaimType` | pxDropdown | 1 | ya | ALWAYS |
| `.Currency` | pxAutoComplete | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.Deduction1` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.Deduction2` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.DueDate` | pxDateTime | 1 | ya | ALWAYS |
| `.DueTo` | pxRadioButtons | 1 | tidak | 1=2 |
| `.EndDate` | pxDateTime | 1 | ya [1==1] | ALWAYS |
| `.ExcessLoss` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.FlagPPH` | pxCheckbox | 1 | tidak | ALWAYS |
| `.GrossPremium` | pxNumber | 1 | ya | ALWAYS |
| `.Installment` | pxTextInput | 1 | ya | ALWAYS |
| `.InstallmentNo` | pxInteger | 1 | ya | ALWAYS |
| `.InstallmentPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.Layer` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerPart` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerPartType` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerType` | pxTextInput | 1 | ya | ALWAYS |
| `.MarketingOfficer` | pxDisplayText | 1 | ya | NOTBLANK |
| `.NetPremium` | pxCurrency | 1 | ya | ALWAYS |
| `.NoOffer` | pxDisplayText | 1 | ya | ALWAYS |
| `.OutstandingClaim` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.OveriddingCommOgp` | pxNumber | 1 | ya [1==1] | ALWAYS |
| `.OveriddingCommOnp` | pxNumber | 1 | ya [1==1] | ALWAYS |
| `.PaymentTotal` | pxNumber | 1 | ya | ALWAYS |
| `.PPHValue` | pxCurrency | 1 | ya | ALWAYS |
| `.PPNValue` | pxCurrency | 1 | ya | ALWAYS |
| `.PremiOgp` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.PremiOnp` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.Quartal` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | 1 | ya [1=1] | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.QuotationData.NoOfferSlip` | pxTextArea | 1 | tidak | ALWAYS |
| `.QuotationData.ProportionalType` | pxDropdown | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.Remark` | pxTextArea | 1 | ya | ALWAYS |
| `.ResultOgp1` | pxCurrency | 1 | ya [1 = 1] | 1==2 |
| `.ResultOgp2` | pxCurrency | 1 | ya | ALWAYS |
| `.ResultOnp1` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.ResultOnp2` | pxCurrency | 1 | ya | ALWAYS |
| `.RiCommOgp` | pxNumber | 1 | ya [1==1] | ALWAYS |
| `.RiCommOnp` | pxNumber | 1 | ya [1==1] | ALWAYS |
| `.SalvageValue` | pxCurrency | 1 | ya [1==1] | ALWAYS |
| `.ShareCurrency` | pxTextInput | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | ya | ALWAYS |
| `.ShareValue` | pxCurrency | 1 | ya | ALWAYS |
| `.SOBName` | pxDisplayText | 1 | ya | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.StartDate` | pxDateTime | 1 | ya [1==1] | ALWAYS |
| `.StatementDate` | pxDateTime | 1 | ya [1==1] | ALWAYS |
| `.StatementType` | pxDropdown | 1 | ya | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.TotalClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalPremium` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentageClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentagePremium` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyGroupName` | pxDisplayText | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.TreatyType` | pxDropdown/pxTextInput | 3 | ya | ALWAYS |
| `.TreatyYear` | pxTextInput | 2 | ya | .IsNewPolicyNonProp != 1 |
| `.TypeTax` | pxDisplayText | 1 | ya | NOTBLANK |
| `.YearOfQuartal` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.PolicyNo` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.ProductionDate` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.Commencement` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.Termination` | pxDateTime | 1 | ya | ALWAYS |

#### GeneralPolicyTreatyIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-PolicyTreatyIn` - field terikat 87 (unik 76) - read-only 49 - RepeatGrid 3 - grid: .BreakDownSpreadList, .ListInstallment, .SpreadingRiskList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BalanceBeforePPH` | pxTextInput | 1 | ya | .BalanceDueTo >= 0 |
| `.BalanceBeforeTax` | pxTextInput | 1 | ya | .BalanceDueTo >= 0 |
| `.BalanceDueTo` | pxTextInput | 2 | ya | .BalanceDueTo < 0 ; .BalanceDueTo >= 0 |
| `.BizName` | pxTextInput | 1 | ya | 1=2 |
| `.CedingCoName` | pxTextInput | 1 | ya | ALWAYS |
| `.Claim` | pxCurrency | 1 | tidak | ALWAYS |
| `.ClaimPaymentType` | pxDropdown | 1 | tidak | ALWAYS |
| `.ClaimPercentage` | pxNumber | 1 | tidak | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.ClaimType` | pxDropdown | 1 | tidak | ALWAYS |
| `.Deduction1` | pxCurrency | 1 | tidak | ALWAYS |
| `.Deduction2` | pxCurrency | 1 | tidak | ALWAYS |
| `.DueDate` | pxDateTime | 1 | ya [IsUW \|\| IsClaim \|\| pyPortal.IsShowOpenPolicy...] | ALWAYS |
| `.DueTo` | pxRadioButtons | 1 | tidak | NEVER |
| `.EndDate` | pxDateTime | 1 | tidak | ALWAYS |
| `.ExcessLoss` | pxCurrency | 1 | tidak | ALWAYS |
| `.FlagPPH` | pxCheckbox | 1 | tidak | ALWAYS |
| `.FlagRetroTreaty` | pxCheckbox | 1 | tidak | .ClaimType != 'XOL Retro' |
| `.GrossClaim` | pxNumber | 1 | tidak | ALWAYS |
| `.GrossPremium` | pxNumber | 1 | tidak | ALWAYS |
| `.IDCurrency` | pxDropdown | 1 | tidak | .IsNewPolicyNonProp != 1 |
| `.Installment` | pxTextInput | 1 | tidak | ALWAYS |
| `.InstallmentNo` | pxInteger | 1 | ya | ALWAYS |
| `.InstallmentPercentage` | pxNumber | 1 | ya [IsUW \|\| pyPortal.IsShowOpenPolicyMarine == tru...] | ALWAYS |
| `.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.Layer` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerPart` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerPartType` | pxTextInput | 1 | ya | ALWAYS |
| `.LayerType` | pxTextInput | 1 | ya | ALWAYS |
| `.NetPremium` | pxCurrency | 1 | ya | ALWAYS |
| `.NoOffer` | pxDisplayText | 1 | ya | ALWAYS |
| `.OutstandingClaim` | pxCurrency | 1 | tidak | ALWAYS |
| `.OveriddingCommOgp` | pxNumber | 1 | tidak | ALWAYS |
| `.OveriddingCommOnp` | pxNumber | 1 | tidak | ALWAYS |
| `.PaymentTotal` | pxNumber | 1 | ya | ALWAYS |
| `.PPHValue` | pxCurrency | 1 | ya | ALWAYS |
| `.PPNValue` | pxCurrency | 1 | ya | ALWAYS |
| `.PremiOgp` | pxCurrency | 1 | tidak | ALWAYS |
| `.PremiOnp` | pxCurrency | 1 | tidak | ALWAYS |
| `.Premium` | pxNumber | 1 | ya [IsUW \|\| pyPortal.IsShowOpenPolicyMarine == tru...] | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.Quartal` | pxTextInput | 1 | tidak | ALWAYS |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | 1 | tidak [IsUW] | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.QuotationData.MOID` | pxDropdown | 1 | tidak | ALWAYS |
| `.QuotationData.NoOfferSlip` | pxTextArea | 1 | tidak | ALWAYS |
| `.QuotationData.ProportionalType` | pxDropdown | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.Remark` | pxTextArea | 1 | tidak | ALWAYS |
| `.ResultOgp1` | pxCurrency | 1 | tidak | ALWAYS |
| `.ResultOgp2` | pxCurrency | 1 | tidak | ALWAYS |
| `.ResultOnp1` | pxCurrency | 1 | tidak | ALWAYS |
| `.ResultOnp2` | pxCurrency | 1 | tidak | ALWAYS |
| `.RiCommOgp` | pxNumber | 1 | tidak | ALWAYS |
| `.RiCommOnp` | pxNumber | 1 | tidak | ALWAYS |
| `.SalvageValue` | pxCurrency | 1 | tidak | ALWAYS |
| `.ShareCurrency` | pxTextInput | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | campur | ALWAYS |
| `.ShareValue` | pxCurrency | 1 | ya | ALWAYS |
| `.StartDate` | pxDateTime | 1 | tidak | ALWAYS |
| `.StatementDate` | pxDateTime | 1 | ya [ALWAYS] | ALWAYS |
| `.StatementType` | pxDropdown | 1 | tidak [IsUW] | pyWorkPage.Quotation.ProportionalType != 'NonProportional' |
| `.TotalClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalPremium` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentageClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentagePremium` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyGroupName` | pxTextInput | 1 | ya | .IsNewPolicyNonProp != 1 |
| `.TreatyType` | pxDropdown/pxTextInput | 3 | campur | ALWAYS |
| `.TreatyYear` | pxTextInput | 2 | ya | .IsNewPolicyNonProp != 1 |
| `.TypeTax` | pxRadioButtons | 1 | tidak | .FlagPPH = true |
| `.YearOfQuartal` | pxTextInput | 1 | tidak | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.SOBName` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber | 2 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.Commencement` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.TreatyIn.Termination` | pxDateTime | 1 | ya | ALWAYS |

#### GolfCoverageSummary

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: D_GrowingTreesCoverageSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | (tanpa) | 1 | ya | ALWAYS |

#### GolfLocationDetail

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 2 - grid: D_AnekaList.pxResults, D_GrowingTreesCoverageSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Currency.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectName` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | (tanpa) | 1 | ya | ALWAYS |
| `.TSI` | (tanpa) | 1 | ya | ALWAYS |

#### GolfObjectItemSummary

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: D_AnekaList.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectName` | (tanpa) | 1 | ya | ALWAYS |
| `.TSI` | (tanpa) | 1 | ya | ALWAYS |

#### GolfSummarySection

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: D_GolfLocationSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.RiskLocation.ASMAddress` | (tanpa) | 1 | ya | ALWAYS |

#### GridViewFollowingNB

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-policyjson` - field terikat 8 (unik 8) - read-only 3 - RepeatGrid 1 - grid: TempListFollowing.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CaseID` | pxTextInput | 1 | tidak | ALWAYS |
| `.CustomerName` | pxTextInput | 1 | tidak | ALWAYS |
| `.EndDateTime` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyNo` | pxTextInput | 1 | tidak | ALWAYS |
| `.QQ` | (tanpa) | 1 | tidak | ALWAYS |
| `.RiSlipCeding` | pxDisplayText | 1 | ya | ALWAYS |
| `.SourceOfBusinessName` | pxTextInput | 1 | tidak | ALWAYS |
| `.StartDateTime` | pxDateTime | 1 | ya | ALWAYS |

#### GrowingTreesLocationDetail

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 7 (unik 7) - read-only 7 - RepeatGrid 3 - grid: .Property.RiskLocation.OccupationList, D_AnekaList.pxResults, D_GrowingTreesCoverageSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Currency.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectName` | (tanpa) | 1 | ya | ALWAYS |
| `.OccupationId` | (tanpa) | 1 | ya | ALWAYS |
| `.OccupationName` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | (tanpa) | 1 | ya | ALWAYS |
| `.TSI` | (tanpa) | 1 | ya | ALWAYS |

#### GrowingTreesOccupationSummary

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: .Property.RiskLocation.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OccupationId` | (tanpa) | 1 | ya | ALWAYS |
| `.OccupationName` | (tanpa) | 1 | ya | ALWAYS |

#### GrowingTreesSummarySection

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: D_GrowingTreesLocationSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText | 1 | ya | ALWAYS |

#### HistoricalSurveyReportDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Quotation` - field terikat 4 (unik 4) - read-only  - RepeatGrid 1 - grid: .QuotationData.SurveyReportList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DateofSurvey` | pxDateTime | 1 | tidak | ALWAYS |
| `.QuotationData.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.Remarks` | pxDropdown | 1 | tidak | ALWAYS |
| `.SurveyedBy` | (tanpa) | 1 | tidak | ALWAYS |

#### HistoricalSurveyReportDtlUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Quotation` - field terikat 4 (unik 4) - read-only  - RepeatGrid 1 - grid: .QuotationData.SurveyReportList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DateofSurvey` | pxDateTime | 1 | tidak | ALWAYS |
| `.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.Remarks` | pxDropdown | 1 | tidak | ALWAYS |
| `.SurveyedBy` | (tanpa) | 1 | tidak | ALWAYS |

#### InputAccumulatedType_FacIn

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 6 (unik 5) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputAccumulatedType.AccumulationType` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputAccumulatedType.ID` | pxTextInput | 1 | ya | 1=2 |
| `InputAccumulatedType.Keyword` | pxTextInput | 1 | tidak | ALWAYS |
| `InputAccumulatedType.Note` | pxTextInput | 1 | tidak | ALWAYS |
| `InputAccumulatedType.Type` | pxDisplayText/pxHidden | 2 | campur | 1==2 |

#### InputAccumulationCov

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 9 (unik 9) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputAccumulation.AccumulationName` | pxDisplayText | 1 | ya | NOTBLANK |
| `InputAccumulation.AccumulationType` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputAccumulation.CZone` | pxAutoComplete | 1 | tidak [1=1] | ALWAYS |
| `InputAccumulation.ID` | pxTextInput | 1 | ya | 1=2 |
| `InputAccumulation.Keyword` | pxDropdown | 1 | tidak | ALWAYS |
| `InputAccumulation.Note` | pxTextArea | 1 | tidak | ALWAYS |
| `InputAccumulation.PostalCode` | pxAutoComplete | 1 | ya [1!=2] | ALWAYS |
| `InputAccumulation.ProvinceName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputAccumulation.ScopeArea` | pxDropdown | 1 | tidak | ALWAYS |

#### InputAnekaParticipant

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 1 - grid: .AnekaParticipantList

_Tidak ada properti terikat._

#### InputBranch

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 16 (unik 16) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputBranch.BASTerritory` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.BranchName` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.BranchStatus` | pxRadioButtons | 1 | tidak | ALWAYS |
| `InputBranch.BranchType` | pxRadioButtons | 1 | tidak | ALWAYS |
| `InputBranch.CompanyID` | pxHidden | 1 | tidak | ALWAYS |
| `InputBranch.HSGradSalary` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.ID` | pxTextInput | 1 | ya | ALWAYS |
| `InputBranch.JabodetabekStatus` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.JHSGradSalary` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.Kanwil` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.KanwilGroup` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.Name` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.Status` | pxRadioButtons | 1 | tidak | ALWAYS |
| `InputBranch.Telephone` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranch.UndergradSalary` | pxTextInput | 1 | tidak | ALWAYS |
| `InputBranchDetail.CompanyID` | pxHidden | 1 | ya | ALWAYS |

#### InputCauseOfDecline

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-CauseOfDecline` - field terikat 9 (unik 7) - read-only 2 - RepeatGrid 2 - grid: CauseOfDeclineList.pxResults, CauseOfDeclineList2.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CheckBox1` | pxCheckbox | 2 | tidak [.FlagDisable != true] | ALWAYS |
| `.CheckBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.Description1` | pxAutoComplete | 1 | tidak [.TERMCONDITION != 'Period Of Insurance'] | ALWAYS |
| `.Description2` | pxTextInput | 1 | tidak | ALWAYS |
| `.PctRNM` | pxNumber | 1 | tidak | ALWAYS |
| `.PctSOB` | pxNumber | 1 | tidak | ALWAYS |
| `.TERMCONDITION` | pxTextInput | 2 | ya | ALWAYS |

#### InputCauseOfLoss_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-CauseOfLoss` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CauseOfLoss` | pxTextInput | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `.Claim` | pxNumber | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | .Quotation.BusinessCode='04' |
| `.DateOfLoss` | pxDateTime | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `.Detail` | pxTextArea | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `.LossObject` | pxTextInput | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `.PreventionOfLoss` | pxNumber | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `.Remarks` | pxDropdown | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `pyWorkPage.OfferFacIn.QuotationData.InsuredName` | (tanpa) | 1 | ya | ALWAYS |

#### InputCauseOfLoss_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-CauseOfLoss` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CauseOfLoss` | pxTextInput | 1 | ya | ALWAYS |
| `.Claim` | pxTextInput | 1 | ya | ALWAYS |
| `.Currency` | pxAutoComplete | 1 | ya | .Quotation.BusinessCode='04' |
| `.DateOfLoss` | pxDateTime | 1 | ya | ALWAYS |
| `.Detail` | pxTextArea | 1 | ya | ALWAYS |
| `.LossObject` | pxTextInput | 1 | ya | ALWAYS |
| `.PreventionOfLoss` | pxTextInput | 1 | ya | ALWAYS |
| `.Remarks` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.QuotationData.InsuredName` | (tanpa) | 1 | ya | ALWAYS |

#### InputCity

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 7 (unik 7) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputCity.BranchName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputCity.Email` | pxTextInput | 1 | tidak | ALWAYS |
| `InputCity.ID` | pxTextInput | 1 | ya | 1==2 |
| `InputCity.JABODETABEKStatus` | pxRadioButtons | 1 | tidak | ALWAYS |
| `InputCity.MOName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputCity.Note` | pxTextArea | 1 | tidak | ALWAYS |
| `InputCity.ProvinceName` | pxAutoComplete | 1 | tidak | ALWAYS |

#### InputClause

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ArgumentCount` | pxTextInput | 1 | ya | ALWAYS |
| `.ClauseCode` | pxTextInput | 1 | ya | ALWAYS |
| `.ClauseDescription` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ClauseLanguage` | pxTextInput | 1 | ya | ALWAYS |
| `.KodeSection` | (tanpa) | 1 | ya [IsSpreadingUW] | pyWorkPage.Quotation.BusinessCode = 38 \|\| pyWorkPage.Quotation.Busin... |

#### InputClause_ViewDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Argument` - field terikat 4 (unik 4) - read-only 3 - RepeatGrid 1 - grid: .ArgumentList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ArgumentDescription` | pxTextInput | 1 | ya | ALWAYS |
| `.ArgumentNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.ArgumentValue` | pxTextInput | 1 | tidak | ALWAYS |
| `.ClauseContent` | pxRichTextEditor | 1 | ya | ALWAYS |

#### InputClauseFire_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ArgumentCount` | pxTextInput | 1 | ya | .ArgumentCount > 0 |
| `.ClauseCode` | pxTextInput | 1 | ya | ALWAYS |
| `.ClauseDescription` | pxTextArea | 1 | ya [IsSpreadingUW\|\| IsUW] | ALWAYS |
| `.ClauseTitle` | pxTextArea | 1 | ya | ALWAYS |

#### InputClauseFire_FacInUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ArgumentCount` | pxTextInput | 1 | ya | .ArgumentCount > 0 |
| `.ClauseCode` | pxTextInput | 1 | ya | ALWAYS |
| `.ClauseDescription` | pxTextArea | 1 | ya | ALWAYS |
| `.ClauseTitle` | pxTextArea | 1 | ya | ALWAYS |

#### InputClauseFire_ViewDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Argument` - field terikat 4 (unik 4) - read-only 3 - RepeatGrid 1 - grid: .ArgumentList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ArgumentDescription` | pxTextInput | 1 | ya | ALWAYS |
| `.ArgumentNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.ArgumentValue` | pxTextInput | 1 | tidak | ALWAYS |
| `.ClauseContent` | pxRichTextEditor | 1 | ya | ALWAYS |

#### InputCommentDeductible

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.InputCommentDeductible` | pxTextArea | 1 | tidak | .IsDeductibleAcceptance = true |
| `.Password` | pxPassword | 1 | tidak | ALWAYS |

#### InputCommissionLife_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only  - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 1 | tidak | ALWAYS |
| `.IndexCoverage` | (tanpa) | 1 | tidak | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxTextInput | 1 | tidak | ALWAYS |

#### InputCommissionMBU_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 10 (unik 6) - read-only 10 - RepeatGrid 2 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 2 | ya | ALWAYS |
| `.PremiNusantaraReOld` | pxNumber | 1 | ya | ALWAYS |
| `.PremiRp` | pxNumber | 2 | ya | ALWAYS |
| `.PremiRpOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 2 | ya | ALWAYS |

#### InputCommissionMC_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 8 (unik 5) - read-only 8 - RepeatGrid 2 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |

#### InputCommissionPA_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 8 (unik 5) - read-only 6 - RepeatGrid 2 - grid: .ASMCoverage

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | tidak | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |

#### InputCoverageAneka_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 31 (unik 30) - read-only 24 - RepeatGrid 2 - grid: .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | ya | ALWAYS |
| `.CalculateMethod_FacIn` | pxDropdown | 1 | ya [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.Condition` | pxDropdown | 1 | ya | .Condition != 5 |
| `.CoverageNote` | pxAutoComplete | 1 | ya | !IsYieldShortfall |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.Currency.Name` | pxDropdown | 1 | ya | IsPortRisk \|\| IsBuilderRisk\|\|IsObjectWithQuantityYear \|\| isAllRi... |
| `.Day` | pxRadioButtons | 1 | ya [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.Descriptions` | pxTextArea | 1 | ya | ALWAYS |
| `.Discount` | pxNumber | 1 | ya [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.DiscountPercentage` | pxTextInput | 1 | ya [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.Indemnity` | pxNumber | 1 | tidak [IsClaim \|\| IsRenewal \|\| IsUW \|\| pyPortal.I...] | IsMBD |
| `.IndemnityPercentage` | pxNumber | 1 | ya | IsMBD |
| `.InputCondition` | pxTextInput | 1 | ya | .Condition = 5 |
| `.LimitofLiability` | pxNumber | 1 | ya [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.Loading` | pxTextInput | 1 | ya [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.LostLimit` | pxTextInput | 1 | tidak [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.MinMax` | pxDropdown | 1 | ya | ALWAYS |
| `.PctDeductible` | pxNumber | 1 | ya | ALWAYS |
| `.PctDeductible2` | pxNumber | 1 | ya | .MinMax = 3 |
| `.PctLoL` | pxNumber | 1 | ya [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.PctShortPeriod` | pxTextInput | 2 | ya [IsClaim \|\| IsUW \|\| IsGroupUWFac \|\| IsSprea...] | .CalculateMethod_FacIn==2 ; isEDM && .CalculateMethod==2 |
| `.Premium` | pxNumber | 1 | tidak [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1] | ALWAYS |
| `.Rate` | pxTextInput | 1 | tidak [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1 ...] | ALWAYS |
| `.SubLimitNote` | pxTextArea | 1 | tidak | ALWAYS |
| `.TimeExcess` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TypeDeductible` | pxDropdown | 1 | ya | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | ya | .MinMax = 3 |
| `.Unit` | pxDropdown | 1 | tidak [IsClaim \|\| IsRenewal \|\| IsUW \|\| pyPortal.I...] | IsMBD |
| `pyWorkPage.OfferFacIn.ProRatePercent` | pxNumber | 1 | tidak [IsClaim \|\| IsUW \|\| pyWorkPage.IsOldData = 1] | .CalculateMethod_FacIn==1 |

#### InputCoverageCargo_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 25 (unik 23) - read-only 24 - RepeatGrid 3 - grid: .AdditionalCoverage, .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | ya | ALWAYS |
| `.Condition` | pxDropdown | 1 | ya | .Condition != 5 |
| `.Coverage` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.CoverageInitial` | pxTextInput | 1 | ya [(.IsOldData = 'old' && (IsEdmPerubahan \|\| IsED...] | ALWAYS |
| `.CoverageNote` | pxTextArea | 2 | ya | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.Currency.Name` | pxDropdown | 1 | ya [IsRenewal \|\| IsUW \|\| pyPortal.IsShowOpenPoli...] | ALWAYS |
| `.CurrencyMaster` | pxTextInput | 1 | ya | ALWAYS |
| `.Descriptions` | pxTextArea | 1 | ya | ALWAYS |
| `.Discount` | pxNumber | 1 | ya | ALWAYS |
| `.DiscountPercentage` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsUW] | ALWAYS |
| `.InputCondition` | pxTextInput | 1 | ya | .Condition = 5 |
| `.LimitofLiability` | pxNumber | 2 | ya [IsUW] | pyWorkPage.OfferFacIn.PolicyMasterNumber !='RNM-F04.10.2022.04666' ; p... |
| `.MinMax` | pxDropdown | 1 | ya | ALWAYS |
| `.MinPremium` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsUW] | ALWAYS |
| `.PctDeductible` | pxNumber | 1 | ya | ALWAYS |
| `.PctDeductible2` | pxNumber | 1 | ya | .MinMax = 3 |
| `.Premium` | pxNumber | 1 | ya [IsUW \|\| pyPortal.IsShowOpenPolicyMarine == tru...] | ALWAYS |
| `.Rate` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsUW] | ALWAYS |
| `.TimeExcess` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya [IsUW] | ALWAYS |
| `.TypeDeductible` | pxDropdown | 1 | ya | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | ya | .MinMax = 3 |

#### InputCoverageCargo_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 31 (unik 25) - read-only 31 - RepeatGrid 3 - grid: .AdditionalCoverage, .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | ya | ALWAYS |
| `.Condition` | pxDropdown | 1 | ya | .Condition != 5 |
| `.Coverage` | pxAutoComplete | 1 | ya | ALWAYS |
| `.CoverageInitial` | pxTextArea | 1 | ya | ALWAYS |
| `.CoverageNote` | pxTextArea | 2 | ya | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.Currency.Name` | pxAutoComplete | 1 | ya | ALWAYS |
| `.CurrencyMaster` | pxTextInput | 1 | ya | pyWorkPage.OfferFacIn.PolicyMasterNumber ='RNM-F04.10.2022.04666' |
| `.Descriptions` | pxTextArea | 1 | ya | ALWAYS |
| `.Discount` | pxNumber | 2 | ya | ALWAYS |
| `.DiscountPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.InputCondition` | pxTextInput | 1 | ya | .Condition = 5 |
| `.LimitofLiability` | pxNumber | 2 | ya | pyWorkPage.OfferFacIn.PolicyMasterNumber !='RNM-F04.10.2022.04666' ; p... |
| `.MinMax` | pxDropdown | 1 | ya | ALWAYS |
| `.MinPremium` | pxNumber | 2 | ya | ALWAYS |
| `.PctDeductible` | pxNumber | 1 | ya | ALWAYS |
| `.PctDeductible2` | pxNumber | 1 | ya | .MinMax = 3 |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |
| `.TimeExcess` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |
| `.TypeDeductible` | pxDropdown | 1 | ya | ALWAYS |
| `.TypeDeductible2` | pxDropdown | 1 | ya | .MinMax = 3 |

#### InputCoverageCommisionFire

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OLDID` | pxTextInput | 1 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |

#### InputCoverageDeductibleAneka_GCNM

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Quantity` | pxTextInput | 1 | ya | ALWAYS |
| `.SumDiscountObjectAneka` | pxTextInput | 1 | ya | ALWAYS |
| `.SumPremiObjectAneka` | pxTextInput | 1 | ya | ALWAYS |

#### InputCoverageDeductibleHull_GCNM

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Quantity` | pxTextInput | 1 | ya | ALWAYS |
| `.SumDiscountObjectAneka` | pxTextInput | 1 | ya | ALWAYS |
| `.SumPremiObjectAneka` | pxTextInput | 1 | ya | ALWAYS |

#### InputCoverageFire

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 5 (unik 5) - read-only 3 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OLDID` | pxTextInput | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.RateOJK` | pxNumber | 1 | ya | ALWAYS |
| `.SelectedLocationAddress` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.SelectedObjectItem` | pxAutoComplete | 1 | tidak | .SelectedLocationAddress != '' |

#### InputCoverageFire_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 5 (unik 5) - read-only 3 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OLDID` | pxTextInput | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.RateOJK` | pxNumber | 1 | ya | ALWAYS |
| `.SelectedLocationAddress` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.SelectedObjectItem` | pxAutoComplete | 1 | tidak | .SelectedLocationAddress != '' |

#### InputCoverageSpreadingAneka

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 12 (unik 7) - read-only 12 - RepeatGrid 2 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 2 | ya | ALWAYS |
| `.Loading` | pxNumber | 2 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 2 | ya | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |

#### InputCoverageSpreadingAneka_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Loading` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |

#### InputCoverageSpreadingFire

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OLDID` | pxTextInput | 1 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |

#### InputCoverageSpreadingFire_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OLDID` | pxTextInput | 1 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |

#### InputCoverageSpreadingMarineCargo

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 10 (unik 6) - read-only 10 - RepeatGrid 2 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 2 | ya | ALWAYS |
| `.PremiNusantaraReOld` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 2 | ya | ALWAYS |

#### InputCoverageSpreadingMarineCargoISUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 1 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |

#### InputCoverageSpreadingMBU

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 10 (unik 6) - read-only 10 - RepeatGrid 2 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 2 | ya | ALWAYS |
| `.PremiNusantaraReOld` | pxNumber | 1 | ya | ALWAYS |
| `.PremiRp` | pxNumber | 2 | ya | ALWAYS |
| `.PremiRpOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 2 | ya | ALWAYS |

#### InputCoverageSpreadingMBUISUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 1 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |

#### InputDeductible_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 11 (unik 8) - read-only 10 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BasisType` | pxDropdown | 1 | ya [IsSpreadingUW] | .DeductibleType = 4 \|\| .DeductibleType = 7 \|\| .DeductibleType = 9 |
| `.ClaimCategory` | pxDropdown | 1 | ya [IsSpreadingUW] | !IsFire |
| `.Coverage` | pxAutoComplete/pxDropdown | 2 | ya [IsSpreadingUW ; pyPortal.IsShowOpenPolicyMarine ...] | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [IsSpreadingUW] | IsDeducType |
| `.DeductibleType` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.Descriptions` | pxTextArea | 3 | campur [IsSpreadingUW ; pyPortal.IsShowOpenPolicyMarine ...] | IsDescDeductType |
| `.Type` | (tanpa) | 1 | ya [IsSpreadingUW \|\| IsFlagDelete] | IsTypeDeductType |
| `.Value` | pxNumber | 1 | ya [IsSpreadingUW] | IsValueDeductType |

#### InputDeductible_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 9 (unik 8) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BasisType` | pxDropdown | 1 | ya [IsSpreadingUW] | .DeductibleType = 4 \|\| .DeductibleType = 7 \|\| .DeductibleType = 9 |
| `.ClaimCategory` | pxDropdown | 1 | ya | !IsFire |
| `.Coverage` | pxDropdown | 1 | ya | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | IsDeducType |
| `.DeductibleType` | pxDropdown | 1 | ya | ALWAYS |
| `.Descriptions` | pxTextArea | 2 | ya | IsDescDeductType |
| `.Type` | (tanpa) | 1 | ya | IsTypeDeductType |
| `.Value` | pxNumber | 1 | ya | IsValueDeductType |

#### InputDeductible_GCNM

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 1 - grid: .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BeginDate` | pxTextInput | 1 | ya | ALWAYS |
| `.CoverageNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Descriptions` | (tanpa) | 1 | ya | ALWAYS |
| `.EndDate` | pxTextInput | 1 | ya | ALWAYS |
| `.MinPremium` | pxTextInput | 1 | ya | ALWAYS |
| `.Premium` | pxTextInput | 1 | ya | ALWAYS |
| `.ProRatePercent` | pxTextInput | 1 | ya | ALWAYS |
| `.Rate` | pxTextInput | 1 | ya | ALWAYS |
| `.TSI` | pxTextInput | 1 | ya | ALWAYS |

#### InputDeductibleHull_GCNM

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 1 - grid: .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BeginDate` | pxTextInput | 1 | ya | ALWAYS |
| `.CoverageNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Descriptions` | (tanpa) | 1 | ya | ALWAYS |
| `.EndDate` | pxTextInput | 1 | ya | ALWAYS |
| `.MinPremium` | pxTextInput | 1 | ya | ALWAYS |
| `.Premium` | pxTextInput | 1 | ya | ALWAYS |
| `.ProRatePercent` | pxTextInput | 1 | ya | ALWAYS |
| `.Rate` | pxTextInput | 1 | ya | ALWAYS |
| `.TSI` | pxTextInput | 1 | ya | ALWAYS |

#### InputDistrict

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 3 (unik 3) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputDistrict.CityName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputDistrict.DistrictName` | pxTextInput | 1 | tidak | ALWAYS |
| `InputDistrict.ID` | pxTextInput | 1 | ya | 1==2 |

#### InputDtlAnekaPolicySchedule

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AnekaPolicySchedule.DeductibleInsured1` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaPolicySchedule.DeductibleInsured2` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaPolicySchedule.DeductibleRisk` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaPolicySchedule.InsuredItemSection1` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaPolicySchedule.InsuredItemSection2` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaPolicySchedule.LimitOfIndemnitySection1` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaPolicySchedule.LimitOfIndemnitySection2` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaPolicySchedule.Risk` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaPolicySchedule.SumInsured` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlCargo

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Cargo` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ConveyanceNote` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.GoodNote` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.PackingNote` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TradingNote` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlCargo_FacIn

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Cargo` - field terikat 37 (unik 37) - read-only 34 - RepeatGrid 1 - grid: .PolicyData.Ship.AdditionalShip

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ConveyanceID` | pxDropdown | 1 | ya [pyPortal.IsShowOpenPolicyMarine == true \|\| IsS...] | ALWAYS |
| `.FromRute` | pxTextInput | 1 | ya [IsUW] | ALWAYS |
| `.GoodID` | pxDropdown | 1 | ya [pyPortal.IsShowOpenPolicyMarine == true \|\| IsS...] | ALWAYS |
| `.ID` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 1 | ya [pyPortal.IsShowOpenPolicyMarine == true \|\| IsS...] | ALWAYS |
| `.NM_SHIP` | pxTextInput | 1 | tidak | ALWAYS |
| `.PackingID` | pxDropdown | 1 | ya [pyPortal.IsShowOpenPolicyMarine == true \|\| IsS...] | ALWAYS |
| `.PolicyData.BLNumber` | pxTextArea | 1 | ya [IsUW] | ALWAYS |
| `.PolicyData.InvoiceDate` | pxDateTime | 1 | ya [IsUW] | ALWAYS |
| `.PolicyData.InvoiceNumber` | pxTextInput | 1 | ya [IsUW] | ALWAYS |
| `.PolicyData.LC.LC` | pxCheckbox | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCBANK` | pxTextInput | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCCondition` | pxTextInput | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCEndDate` | pxDateTime | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCNumber` | pxTextInput | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCRemark` | pxTextInput | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.MaxAge` | pxAutoComplete | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.MaxAgeDetail` | pxTextArea | 1 | ya | ALWAYS |
| `.PolicyData.SailDate` | pxDateTime | 1 | ya [IsUW] | IsUW |
| `.PolicyData.Ship.AGE` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.CONST` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.DWT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.FLAG_NAME` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.GRT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NM_SHIP` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NRT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.REMARK` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.SHIPTYPE_NAME` | pxDropdown | 1 | ya | ALWAYS |
| `.PolicyData.Ship.Y_MAKE1` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.ASMAgentID` | pxAutoComplete | 1 | ya [IsUW] | 1=2 |
| `.PolicyData.SurveyAgent.pyAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyCity` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyCountry` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyFullName` | pxTextArea | 1 | ya | ALWAYS |
| `.ToRute` | pxTextInput | 1 | ya [IsUW] | ALWAYS |
| `.TradingID` | pxDropdown | 1 | ya [pyPortal.IsShowOpenPolicyMarine == true \|\| IsS...] | ALWAYS |
| `.VOYAGE_NO` | pxTextInput | 1 | tidak | ALWAYS |

#### InputDtlCargo_FacIn_IsUW

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Cargo` - field terikat 38 (unik 38) - read-only 36 - RepeatGrid 1 - grid: .PolicyData.Ship.AdditionalShip

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ConveyanceNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.FromRute` | pxTextInput | 1 | ya | ALWAYS |
| `.GoodNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 1 | tidak | ALWAYS |
| `.NM_SHIP` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.PolicyData.SurveyAgent.pyCity` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.PolicyData.SurveyAgent.pyCountry` | pxTextInput | 1 | ya | ALWAYS |
| `.PackingNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PolicyData.BLNumber` | pxTextArea | 1 | ya | ALWAYS |
| `.PolicyData.InvoiceDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyData.InvoiceNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.LC.LC` | pxCheckbox | 1 | tidak | ALWAYS |
| `.PolicyData.LC.LCBANK` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.LC.LCCondition` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.LC.LCEndDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyData.LC.LCNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.LC.LCRemark` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.MaxAge` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PolicyData.MaxAgeDetail` | pxTextArea | 1 | ya | ALWAYS |
| `.PolicyData.SailDate` | pxDateTime | 1 | ya | IsUW |
| `.PolicyData.Ship.AGE` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.CONST` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.DWT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.FLAG_NAME` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.GRT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.ID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NM_SHIP` | pxDropdown | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NRT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.REMARK` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.SHIPTYPE_NAME` | pxDropdown | 1 | ya | ALWAYS |
| `.PolicyData.Ship.Y_MAKE1` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.ASMAgentID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyFullName` | pxTextInput | 1 | ya | ALWAYS |
| `.REGISTER_ID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ToRute` | pxTextInput | 1 | ya | ALWAYS |
| `.TradingNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.VOYAGE_NO` | pxTextInput | 1 | ya | ALWAYS |

#### InputDtlClause_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 10 (unik 9) - read-only 6 - RepeatGrid 7 - grid: .ClauseList, .Policy.ClauseList, .WarrantyList, AttachmentClauses.pxResults, pgRepPgSubSectionInputDtlClause_FacInBBBBBBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClauseCode` | pxDisplayText | 2 | ya | ALWAYS |
| `.ClauseContent` | pxDisplayText | 1 | ya | ALWAYS |
| `.ClauseTitle` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | ALWAYS |
| `.CountAttach` | pxTextInput | 1 | tidak | ALWAYS |
| `.NOTE` | pxTextInput | 1 | tidak | ALWAYS |
| `.PNOTE` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyMemo` | ReasAttachmentDescription_name | 1 | tidak | ALWAYS |
| `.WarrantyDesc` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | ALWAYS |
| `.WarrantyID` | pxDisplayText | 1 | ya | ALWAYS |

#### InputDtlClause_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 10 (unik 9) - read-only 6 - RepeatGrid 7 - grid: .ClauseList, .Policy.ClauseList, .WarrantyList, AttachmentClauses.pxResults, pgRepPgSubSectionInputDtlClause_FacIn_IsUWBBBBBBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClauseCode` | pxDisplayText | 2 | ya | ALWAYS |
| `.ClauseContent` | pxDisplayText | 1 | ya | ALWAYS |
| `.ClauseTitle` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | ALWAYS |
| `.CountAttach` | pxTextInput | 1 | tidak | ALWAYS |
| `.NOTE` | pxTextInput | 1 | tidak | ALWAYS |
| `.PNOTE` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyMemo` | ReasAttachmentDescription_name | 1 | tidak | ALWAYS |
| `.WarrantyDesc` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | ALWAYS |
| `.WarrantyID` | pxDisplayText | 1 | ya | ALWAYS |

#### InputDtlConveyance

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Conveyance` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ConveyanceDutyRange` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ConveyanceID` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ConveyanceNote` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | never |

#### InputDtlConveyance_IsUW

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Conveyance` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ConveyanceDutyRange` | pxTextInput | 1 | ya | ALWAYS |
| `.ConveyanceID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ConveyanceNote` | pxTextArea | 1 | ya | ALWAYS |

#### InputDtlCoverage_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 3 - grid: .CoverageList, .TotalTSIPremiGrossList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageInitial` | (tanpa) | 1 | ya | ALWAYS |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | (tanpa) | 1 | ya | ALWAYS |
| `.Rate` | (tanpa) | 1 | ya | ALWAYS |
| `.TSI` | (tanpa) | 1 | ya | ALWAYS |

#### InputDtlCoverage_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 3 - grid: .CoverageList, .TotalTSIPremiGrossList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageInitial` | (tanpa) | 1 | ya | ALWAYS |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | (tanpa) | 1 | ya | ALWAYS |
| `.Rate` | (tanpa) | 1 | ya | ALWAYS |
| `.TSI` | (tanpa) | 1 | ya | ALWAYS |

#### InputDtlGoods

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Packing` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 1 - grid: .PackingList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.GoodID` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac ...] | ALWAYS |
| `.GoodName2` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac ...] | .GoodID ='10022' |
| `.PackingID` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac ...] | ALWAYS |
| `.PackingNote` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac ...] | never |
| `.pxListSubscript` | (tanpa) | 1 | ya | ALWAYS |

#### InputDtlGoods_IsUW

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Packing` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .PackingList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.GoodID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.GoodName` | pxTextArea | 1 | ya | ALWAYS |
| `.PackingID` | pxDropdown | 1 | ya | ALWAYS |
| `.PackingNote` | pxTextInput | 1 | ya | ALWAYS |

#### InputDtlLayerSpreading_FacIn

rumpun **PA-LIFE** - kelas dominan `ASM-FW-GISFW-Int-RETROCESSIONLIFE` - field terikat 8 (unik 8) - read-only 6 - RepeatGrid 1 - grid: .RetroList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.COMMISION` | pxNumber | 1 | ya [.REINSURERNAME ="PT. Reasuransi Nusantara Makmur...] | NOTBLANK |
| `.COMMISION_AMOUNT` | pxNumber | 1 | ya | ALWAYS |
| `.OVR_COMM` | pxNumber | 1 | ya [.REINSURERNAME ="PT. Reasuransi Nusantara Makmur...] | ALWAYS |
| `.OVR_COMM_AMOUNT` | pxNumber | 1 | ya | NOTBLANK |
| `.PERCENTSHARE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.PREMIUM_SPREADED_GROSS` | pxNumber | 1 | tidak | NOTBLANK |
| `.PREMIUM_SPREADED_NET` | pxNumber | 1 | tidak | NOTBLANK |
| `.REINSURERNAME` | pxDisplayText | 1 | ya | NOTBLANK |

#### InputDtlMaintenance

rumpun **ANEKA** - kelas dominan `Data-Address` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AnekaMaintenance.EndPeriod` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaMaintenance.ExtendedPeriodDay` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaMaintenance.ExtendedPeriodMonth` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaMaintenance.ExtendedPeriodWeek` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaMaintenance.MaintenancePeriodDay` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaMaintenance.MaintenancePeriodMonth` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaMaintenance.MaintenancePeriodWeek` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaMaintenance.StartPeriod` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaMaintenance.TitleOfContract` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlMemorandumAccept

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-MemorandumAccept` - field terikat 15 (unik 15) - read-only 15 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Accumulation` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ApplicationNumber` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.BuildingStatus` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.Collateral` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ConclusionsProposals` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.InsuredSiup` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.KeyPerson` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.OriginalRate` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PrincipalsBackground` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PrincipalsData` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PrincipalsGuarantee` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PrincipalsOutGo` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.SummaryFinancial` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.SupportingDocument` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.TypeOfWording` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |

#### InputDtlMemorandumAccept_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-MemorandumAccept` - field terikat 15 (unik 15) - read-only 15 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Accumulation` | pxTextInput | 1 | ya | ALWAYS |
| `.ApplicationNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.BuildingStatus` | pxTextArea | 1 | ya | ALWAYS |
| `.Collateral` | pxTextInput | 1 | ya | ALWAYS |
| `.ConclusionsProposals` | pxTextArea | 1 | ya | ALWAYS |
| `.InsuredSiup` | pxTextInput | 1 | ya | ALWAYS |
| `.KeyPerson` | pxTextInput | 1 | ya | ALWAYS |
| `.OriginalRate` | pxTextInput | 1 | ya | ALWAYS |
| `.PrincipalsBackground` | pxTextArea | 1 | ya | ALWAYS |
| `.PrincipalsData` | pxTextArea | 1 | ya | ALWAYS |
| `.PrincipalsGuarantee` | pxTextArea | 1 | ya | ALWAYS |
| `.PrincipalsOutGo` | pxTextInput | 1 | ya | ALWAYS |
| `.SummaryFinancial` | pxTextArea | 1 | ya | ALWAYS |
| `.SupportingDocument` | pxTextArea | 1 | ya | ALWAYS |
| `.TypeOfWording` | pxTextInput | 1 | ya | ALWAYS |

#### InputDtlObject_FacIn

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 103 (unik 56) - read-only 59 - RepeatGrid 22 - grid: .FinalOccupationList, .OfferFacIn.CargoList, .OfferFacIn.ConveyanceList, .OfferFacIn.GoodsList, .OfferFacIn.LocationList, .OfferFacIn.PersonList, .OfferFacIn.TradingList, .OfferFacIn.VehicleList, .PersonList, .PolicyList, .PropertyList, .SubContractList, pyWorkPage.BatchMemberList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Age` | pxInteger | 1 | ya | ALWAYS |
| `.ASMCoverage(1).Loading` | (tanpa) | 3 | ya | NOTBLANK |
| `.ASMCoverage(1).PremiRp` | (tanpa) | 3 | tidak | ALWAYS |
| `.ASMCoverage(1).Premium` | (tanpa) | 3 | tidak | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 6 | campur | ALWAYS |
| `.ASMGender` | pxDropdown | 5 | campur | ALWAYS |
| `.ASMHeight` | pxInteger | 1 | tidak | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMJobName` | pxDropdown | 2 | campur | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 4 | tidak | ALWAYS |
| `.ASMMaritalStatus` | pxDropdown | 3 | tidak | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMRelation` | pxDropdown | 3 | tidak | ALWAYS |
| `.ASMStatusCoverage` | pxDisplayText | 3 | ya | pyWorkPage.ProposalPosition = '2' \|\| pyWorkPage.ProposalPosition = '... |
| `.ASMWeight` | pxInteger | 1 | tidak | ALWAYS |
| `.BookNo` | pxLink | 2 | tidak | ALWAYS |
| `.BrandName` | pxDropdown | 1 | tidak | ALWAYS |
| `.ChassisNumber` | pxTextInput | 1 | tidak | ALWAYS |
| `.ConveyanceDutyRange` | (tanpa) | 1 | tidak | ALWAYS |
| `.ConveyanceID` | pxDropdown | 1 | tidak | ALWAYS |
| `.ConveyanceNote` | pxDisplayText | 2 | ya | ALWAYS |
| `.EmployeeName` | pxDisplayText | 2 | ya | ALWAYS |
| `.EngineNumber` | pxTextInput | 1 | tidak | ALWAYS |
| `.FromRute` | pxDisplayText | 1 | ya | ALWAYS |
| `.GoodID` | pxDisplayText | 1 | ya | ALWAYS |
| `.GoodName` | pxDisplayText | 1 | ya | ALWAYS |
| `.GoodNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.IsMaterialDamage` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 2 | campur | IsFac |
| `.LicensePlate` | pxTextInput | 1 | tidak | ALWAYS |
| `.ModelName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Name` | pxDisplayText | 2 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxDisplayText | 2 | ya | ALWAYS |
| `.PlanIP` | pxDisplayText | 2 | ya | ALWAYS |
| `.PlanOP` | pxDisplayText | 2 | ya | ALWAYS |
| `.PolicyData.SailDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NM_SHIP` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.IsTopRisk` | pxCheckbox | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText/pxTextInput | 2 | ya | ALWAYS |
| `.pxCreateDateTime` | pxDisplayText | 2 | ya | ALWAYS |
| `.pxCreateOpName` | pxDisplayText | 2 | ya | ALWAYS |
| `.pxListSubscript` | (tanpa) | 4 | campur | ALWAYS |
| `.pyFullName` | pxDisplayText | 7 | campur | ALWAYS |
| `.pyLabel` | pxDisplayText | 2 | ya | ALWAYS |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.SPPANo` | pxLink | 2 | tidak | ALWAYS |
| `.SubContractID` | pxTextInput | 1 | ya | ALWAYS |
| `.SubContractName` | pxTextInput | 1 | ya | ALWAYS |
| `.ToRute` | pxDisplayText | 1 | ya | ALWAYS |
| `.TradingNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.Tranship` | pxDisplayText | 1 | ya | ALWAYS |
| `.TypeName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Warehouse` | pxDisplayText | 1 | ya | ALWAYS |

#### InputDtlObject_FacIn_IsUW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 100 (unik 55) - read-only 59 - RepeatGrid 22 - grid: .FinalOccupationList, .OfferFacIn.CargoList, .OfferFacIn.ConveyanceList, .OfferFacIn.GoodsList, .OfferFacIn.LocationList, .OfferFacIn.PersonList, .OfferFacIn.TradingList, .OfferFacIn.VehicleList, .PersonList, .PolicyList, .PropertyList, .SubContractList, pyWorkPage.BatchMemberList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Age` | pxInteger | 1 | ya | ALWAYS |
| `.ASMCoverage(1).Loading` | (tanpa) | 3 | ya | NOTBLANK |
| `.ASMCoverage(1).PremiRp` | (tanpa) | 3 | tidak | ALWAYS |
| `.ASMCoverage(1).Premium` | (tanpa) | 3 | tidak | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 6 | campur | ALWAYS |
| `.ASMGender` | pxDropdown | 5 | campur | ALWAYS |
| `.ASMHeight` | pxInteger | 1 | tidak | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMJobName` | pxDropdown | 2 | campur | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 4 | tidak | ALWAYS |
| `.ASMMaritalStatus` | pxDropdown | 3 | tidak | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMRelation` | pxDropdown | 3 | tidak | ALWAYS |
| `.ASMStatusCoverage` | pxDisplayText | 3 | ya | pyWorkPage.ProposalPosition = '2' \|\| pyWorkPage.ProposalPosition = '... |
| `.ASMWeight` | pxInteger | 1 | tidak | ALWAYS |
| `.BookNo` | pxLink | 2 | tidak | ALWAYS |
| `.BrandName` | pxDropdown | 1 | tidak | ALWAYS |
| `.ChassisNumber` | pxDisplayText | 1 | ya | ALWAYS |
| `.ConveyanceDutyRange` | pxDisplayText | 1 | ya | ALWAYS |
| `.ConveyanceID` | pxDropdown | 1 | tidak | ALWAYS |
| `.ConveyanceNote` | pxDisplayText | 2 | ya | ALWAYS |
| `.EmployeeName` | pxDisplayText | 2 | ya | ALWAYS |
| `.EngineNumber` | pxDisplayText | 1 | ya | ALWAYS |
| `.FromRute` | pxDisplayText | 1 | ya | ALWAYS |
| `.GoodID` | pxDisplayText | 1 | ya | ALWAYS |
| `.GoodName` | pxDisplayText | 1 | ya | ALWAYS |
| `.GoodNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.IsMaterialDamage` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 2 | campur | IsFac |
| `.LicensePlate` | pxDisplayText | 1 | ya | ALWAYS |
| `.ModelName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Name` | pxDisplayText | 2 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxDisplayText | 2 | ya | ALWAYS |
| `.PackingNote` | pxDisplayText | 1 | ya | ALWAYS |
| `.PlanIP` | pxDisplayText | 2 | ya | ALWAYS |
| `.PlanOP` | pxDisplayText | 2 | ya | ALWAYS |
| `.Property.IsTopRisk` | pxCheckbox | 1 | tidak | IsMBD \|\| IsContractorsPlantMachinery \|\| IsHE\|\|IsElectronicEquipm... |
| `.Property.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText/pxTextInput | 2 | ya | ALWAYS |
| `.pxCreateDateTime` | pxDisplayText | 2 | ya | ALWAYS |
| `.pxCreateOpName` | pxDisplayText | 2 | ya | ALWAYS |
| `.pxListSubscript` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyFullName` | pxDisplayText | 7 | campur | ALWAYS |
| `.pyLabel` | pxDisplayText | 2 | ya | ALWAYS |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.SPPANo` | pxLink | 2 | tidak | ALWAYS |
| `.SubContractID` | pxTextInput | 1 | ya | ALWAYS |
| `.SubContractName` | pxTextInput | 1 | ya | ALWAYS |
| `.ToRute` | pxDisplayText | 1 | ya | ALWAYS |
| `.TradingNote` | pxDisplayText | 2 | ya | ALWAYS |
| `.Tranship` | pxDisplayText | 1 | ya | ALWAYS |
| `.TypeName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Warehouse` | pxDisplayText | 1 | ya | ALWAYS |

#### InputDtlObjectAneka_FacIn

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 83 (unik 53) - read-only 83 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AreaHectar` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.AviationHull.MaxPassengers` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.AviationHull.Pilots` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.AviationHull.RegistrationMarks` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.AviationHull.SpecialRentalUses` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.AviationHull.SpecialUses` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.AviationHull.StandardUses` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.BuiltIn` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.Currency.Name` | pxDropdown | 7 | ya [.FlagDelete==1] | ALWAYS |
| `.DateOfBirth` | pxDateTime | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.DepositAccountNo` | pxTextInput | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.DraftWordingID` | pxDropdown | 1 | ya [.FlagDelete==1] | IsBonding |
| `.END_DATE` | pxDateTime | 1 | ya [.FlagDelete==1] | IsObjectLimiliabilityAneka |
| `.Flag` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.GeographicalLimits` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.IDCardNo` | pxTextInput | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.Institution` | pxTextInput | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.InvoiceNo` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.Job` | pxTextInput | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.LimitIdemnityDay` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.LimitIdemnityMonth` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.LoanAccountNumber` | pxTextInput | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.LoanName` | pxTextInput | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.LoanType` | pxTextInput | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.MaintenancePeriod` | pxNumber | 2 | ya [.FlagDelete==1] | IsEngineering |
| `.MarineHull.BREADTHM` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.MarineHull.Classification` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.MarineHull.Construction` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.MarineHull.DEPTHM` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.MarineHull.DWTRT` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.MarineHull.GRTRT` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.MarineHull.LENGTHM` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.MarineHull.NRTRT` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.MarineHull.TypeOfVessel` | pxTextInput | 1 | ya [Always] | ALWAYS |
| `.NO_KONTRAK` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.NO_KTP` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.NO_NPWP` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.NoOfTree` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.ObjectName` | pxAutoComplete/pxTextArea/pxTextInput | 8 | ya [.FlagDelete==1] | ALWAYS |
| `.Others` | pxTextArea | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.Quantity` | pxTextInput | 3 | ya [.FlagDelete==1] | pyWorkPage.Quotation.BusinessCode!=20 |
| `.Salary` | pxNumber | 1 | ya [.FlagDelete==1] | IsCreditBriguna |
| `.Section` | pxDropdown | 3 | ya [.FlagDelete==1] | IsObjectSectionAneka |
| `.SerialNo` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.START_DATE` | pxDateTime | 1 | ya [.FlagDelete==1] | IsObjectLimiliabilityAneka |
| `.TSI` | pxNumber/pxTextInput | 7 | ya [.FlagDelete==1] | IsObjectLimiliabilityAneka |
| `.UnitQuantity` | pxTextInput | 1 | ya [.FlagDelete==1] | pyWorkPage.Quotation.BusinessCode!=20 |
| `.VehicleHE.BrandName` | pxAutoComplete | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.VehicleHE.ChassisNumber` | pxTextInput | 2 | ya [.FlagDelete==1] | ALWAYS |
| `.VehicleHE.EngineNumber` | pxTextInput | 2 | ya [.FlagDelete==1] | ALWAYS |
| `.VehicleHE.ObjectNameHE` | pxAutoComplete | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.VehicleHE.TypeName` | pxAutoComplete | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.Year` | pxTextInput | 5 | ya [.FlagDelete==1 ; Always] | ALWAYS |

#### InputDtlObjectAneka_FacIn_IsUW

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 73 (unik 48) - read-only 73 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AreaHectar` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.MaxPassengers` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.Pilots` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.RegistrationMarks` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.SpecialRentalUses` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.SpecialUses` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.StandardUses` | pxTextInput | 1 | ya | ALWAYS |
| `.BuiltIn` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.Currency.Name` | pxAutoComplete/pxDropdown | 4 | ya | ALWAYS |
| `.DateOfBirth` | pxDateTime | 1 | ya | IsCreditBriguna |
| `.DepositAccountNo` | pxTextInput | 1 | ya | IsCreditBriguna |
| `.DraftWordingID` | pxDropdown | 1 | ya [1!=2] | IsBonding |
| `.Flag` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.GeographicalLimits` | pxTextInput | 1 | ya | ALWAYS |
| `.IDCardNo` | pxTextInput | 1 | ya | IsCreditBriguna |
| `.Institution` | pxTextInput | 1 | ya | IsCreditBriguna |
| `.InvoiceNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Job` | pxTextInput | 1 | ya | IsCreditBriguna |
| `.LimitIdemnityDay` | pxTextInput | 1 | ya | ALWAYS |
| `.LimitIdemnityMonth` | pxTextInput | 1 | ya | ALWAYS |
| `.LoanAccountNumber` | pxTextInput | 1 | ya | IsCreditBriguna |
| `.LoanName` | pxTextInput | 1 | ya | IsCreditBriguna |
| `.LoanType` | pxTextInput | 1 | ya | IsCreditBriguna |
| `.MaintenancePeriod` | pxNumber | 2 | ya | IsEngineering |
| `.MarineHull.BREADTHM` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.MarineHull.Classification` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.MarineHull.Construction` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.MarineHull.DEPTHM` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.MarineHull.DWTRT` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.MarineHull.GRTRT` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.MarineHull.LENGTHM` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.MarineHull.NRTRT` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.MarineHull.TypeOfVessel` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.NoOfTree` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextArea/pxTextInput | 8 | ya [1!=2] | ALWAYS |
| `.Others` | pxTextArea | 1 | ya | ALWAYS |
| `.Quantity` | pxTextInput | 3 | ya | pyWorkPage.Quotation.BusinessCode!=20 |
| `.Salary` | pxNumber | 1 | ya | IsCreditBriguna |
| `.Section` | pxDropdown | 3 | ya | IsObjectSectionAneka |
| `.SerialNo` | pxTextInput | 1 | ya | 1=2 |
| `.TSI` | pxNumber | 5 | ya | IsObjectLimiliabilityAneka |
| `.UnitQuantity` | pxTextInput | 1 | ya | pyWorkPage.Quotation.BusinessCode!=20 |
| `.VehicleHE.BrandName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.VehicleHE.ChassisNumber` | pxTextInput | 2 | ya | ALWAYS |
| `.VehicleHE.EngineNumber` | pxTextInput | 2 | ya | ALWAYS |
| `.VehicleHE.ObjectNameHE` | pxAutoComplete | 1 | ya | ALWAYS |
| `.VehicleHE.TypeName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Year` | pxTextInput | 5 | ya [1!=2] | ALWAYS |

#### InputDtlObjectLocation_FacIn

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 23 (unik 21) - read-only 18 - RepeatGrid 2 - grid: .Property.RiskLocation.AnekaList, .Property.RiskLocation.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextArea | 1 | tidak [IsUW] | ALWAYS |
| `.OccupationId` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Property.AlmRiskID` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.BuildingNo` | pxNumber | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.Property.Country` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.IsTopRisk` | pxCheckbox | 1 | tidak | IsMBD \|\| IsHE \|\| IsContractorsPlantMachinery\|\|IsElectronicEquipm... |
| `.Property.Province` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxTextArea | 2 | ya [IsUW \|\| pyWorkPage.IsOldData = 1 \|\| IsSpread...] | ALWAYS |
| `.Property.RiskLocation.ASMCity` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMDistrict` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMRW` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMZipCode` | pxTextInput | 2 | ya [IsUW \|\| pyWorkPage.IsOldData = 1 \|\| IsSpread...] | ALWAYS |
| `.Property.RiskLocation.CityName` | pxAutoComplete | 1 | ya [IsUW \|\| pyWorkPage.IsOldData = 1] | InputRiskAddress.IDProvince!='' |
| `.Property.RiskLocation.DistrictName` | pxAutoComplete | 1 | ya [IsUW \|\| pyWorkPage.IsOldData = 1] | .Property.RiskLocation.CityName!='' |
| `.Property.RiskLocation.RWName` | pxAutoComplete | 1 | ya [IsUW \|\| pyWorkPage.IsOldData = 1] | .Property.RiskLocation.DistrictName!='' |
| `.Property.RoadName` | pxTextArea | 1 | ya | .Property.AlmRiskID != '' |
| `.Property.RoadType` | pxDropdown | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | tidak | ALWAYS |

#### InputDtlObjectLocation_FacIn_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 15 (unik 15) - read-only 15 - RepeatGrid 2 - grid: .Property.RiskLocation.AnekaList, .Property.RiskLocation.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextArea | 1 | ya | ALWAYS |
| `.OccupationId` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.AlmRiskID` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.BuildingNo` | pxNumber | 1 | ya | ALWAYS |
| `.Property.Country` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.Province` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMCity` | pxTextInput | 1 | ya | InputRiskAddress.IDProvince!='' |
| `.Property.RiskLocation.ASMDistrict` | pxTextInput | 1 | ya | .Property.RiskLocation.CityName!='' |
| `.Property.RiskLocation.ASMRW` | pxTextInput | 1 | ya | .Property.RiskLocation.DistrictName!='' |
| `.Property.RiskLocation.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RoadType` | pxDropdown | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |

#### InputDtlObjectLocation_GISFW

rumpun **FIRE** - kelas dominan `Data-Address` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 1 - grid: .OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.ASMCity` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMDistrict` | pxDropdown | 1 | ya | .ASMCity != '' |
| `.ASMRW` | pxDropdown | 1 | ya | .ASMDistrict != '' |
| `.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationId` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.SumTotalTSI` | (tanpa) | 1 | ya | ALWAYS |

#### InputDtlObjectLocationHull_GISFW

rumpun **INTI+FIRE** - kelas dominan `Data-Address` - field terikat 7 (unik 7) - read-only 5 - RepeatGrid 1 - grid: .OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.ASMCity` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMDistrict` | pxDropdown | 1 | ya | .ASMCity != '' |
| `.ASMRW` | pxDropdown | 1 | ya | .ASMDistrict != '' |
| `.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationId` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlObjFire

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Property` - field terikat 72 (unik 72) - read-only 67 - RepeatGrid 4 - grid: .FEAList, .ListCauseOfLoss, .OccupationList, .PropertyItemList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AlmRiskID` | pxTextInput | 1 | ya | ALWAYS |
| `.Amount` | pxNumber | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.BuildingConstruction.FloorType` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.BuildingConstruction.NumberOfFloor` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.BuildingConstruction.OthersType` | pxTextInput | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.BuildingConstruction.PartitionType` | pxTextInput | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.BuildingConstruction.RoofType` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.BuildingConstruction.SupportWallType` | pxTextInput | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.BuildingConstruction.WallType` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.BuildingNo` | pxTextInput | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.CoinsData.CoinsName` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.Construction` | pxAutoComplete | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.DateOfLoss` | pxDateTime | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.Detail` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.EML` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.FEAID` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.FEAInside` | pxInteger | 1 | ya | ALWAYS |
| `.FEAOutside` | pxInteger | 1 | ya | ALWAYS |
| `.FirstLossLimitFac` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.IsFlammableItemFlag` | pxCheckbox | 1 | tidak | ALWAYS |
| `.IsHotWorkProcessFlag` | pxCheckbox | 1 | tidak | ALWAYS |
| `.IsMaterialDamage` | pxCheckbox | 1 | tidak | ALWAYS |
| `.IsProductionProcessFlag` | pxCheckbox | 1 | tidak | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 1 | tidak | IsFac |
| `.IsUsedFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 1 | ya | ALWAYS |
| `.LossObject` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | .ObjectType = 'Lainnya' |
| `.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectType` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.OccupationId` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.Ownership` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.PctOfContribution` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PML` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.PropertyYear` | pxInteger | 1 | ya | ALWAYS |
| `.pyCity` | pxTextInput | 1 | ya | ALWAYS |
| `.pyCountryName` | pxTextInput | 1 | ya | ALWAYS |
| `.RIComm` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.RiskCategory` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | IsFac |
| `.RiskLocation.ASMAddress` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.RiskLocation.ASMCity` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMDistrict` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMRW` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.RoadName` | pxTextInput | 1 | ya [.AlmRiskID != ''] | ALWAYS |
| `.RoadType` | pxDropdown | 1 | ya [.AlmRiskID != ""] | ALWAYS |
| `.SurroundingRisk.BackConstruction` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.BackDistance` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.BackNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.BackOccupation` | pxAutoComplete | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.FloodArea` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | .SurroundingRisk.FloodAreaStatus==0 |
| `.SurroundingRisk.FloodAreaStatus` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.FrontConstruction` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.FrontDistance` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.FrontNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.FrontOccupation` | pxAutoComplete | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.HousekeepingRemark` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.HousekeepingStatus` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.LeftConstruction` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.LeftDistance` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.LeftNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.LeftOccupation` | pxAutoComplete | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.RightConstruction` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.RightDistance` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.RightNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.RightOccupation` | pxAutoComplete | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.TreatyGroup` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | IsOilGas |
| `.TSIObjectItem` | pxNumber | 1 | ya | ALWAYS |
| `.Unit` | pxInteger | 1 | ya | ALWAYS |
| `.UpdateDate` | pxDateTime | 1 | ya | ALWAYS |

#### InputDtlObjFire_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Property` - field terikat 72 (unik 72) - read-only 71 - RepeatGrid 4 - grid: .FEAList, .ListCauseOfLoss, .OccupationList, .PropertyItemList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AlmRiskID` | pxTextInput | 1 | ya | ALWAYS |
| `.Amount` | pxNumber | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.BuildingConstruction.FloorType` | pxDropdown | 1 | ya | ALWAYS |
| `.BuildingConstruction.NumberOfFloor` | pxNumber | 1 | ya | ALWAYS |
| `.BuildingConstruction.OthersType` | pxTextInput | 1 | ya | ALWAYS |
| `.BuildingConstruction.PartitionType` | pxTextInput | 1 | ya | ALWAYS |
| `.BuildingConstruction.RoofType` | pxDropdown | 1 | ya | ALWAYS |
| `.BuildingConstruction.SupportWallType` | pxTextInput | 1 | ya | ALWAYS |
| `.BuildingConstruction.WallType` | pxDropdown | 1 | ya | ALWAYS |
| `.BuildingNo` | pxTextInput | 1 | ya | ALWAYS |
| `.CoinsData.CoinsName` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.Construction` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.DateOfLoss` | pxDateTime | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.Detail` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.EML` | pxNumber | 1 | ya | ALWAYS |
| `.FEAID` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.FEAInside` | pxInteger | 1 | ya | ALWAYS |
| `.FEAOutside` | pxInteger | 1 | ya | ALWAYS |
| `.FirstLossLimitFac` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.IsFlammableItemFlag` | pxCheckbox | 1 | tidak | ALWAYS |
| `.IsHotWorkProcessFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsMaterialDamage` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsProductionProcessFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 1 | ya | IsFac |
| `.IsUsedFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 1 | ya | ALWAYS |
| `.LossObject` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | .ObjectType = 'Lainnya' |
| `.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectType` | pxDropdown | 1 | ya | ALWAYS |
| `.OccupationId` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.Ownership` | pxDropdown | 1 | ya | ALWAYS |
| `.PctOfContribution` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PML` | pxNumber | 1 | ya | ALWAYS |
| `.PropertyYear` | pxInteger | 1 | ya | ALWAYS |
| `.pyCity` | pxTextInput | 1 | ya | ALWAYS |
| `.pyCountryName` | pxTextInput | 1 | ya | ALWAYS |
| `.RIComm` | pxNumber | 1 | ya | ALWAYS |
| `.RiskCategory` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | IsFac |
| `.RiskLocation.ASMAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.RiskLocation.ASMCity` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMDistrict` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMRW` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.RoadName` | pxTextInput | 1 | ya | ALWAYS |
| `.RoadType` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.BackConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.BackDistance` | pxNumber | 1 | ya | ALWAYS |
| `.SurroundingRisk.BackNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.BackOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.SurroundingRisk.FloodArea` | pxDropdown | 1 | ya | .SurroundingRisk.FloodAreaStatus==0 |
| `.SurroundingRisk.FloodAreaStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.FrontConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.FrontDistance` | pxNumber | 1 | ya | ALWAYS |
| `.SurroundingRisk.FrontNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.FrontOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.SurroundingRisk.HousekeepingRemark` | pxTextArea | 1 | ya | ALWAYS |
| `.SurroundingRisk.HousekeepingStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.LeftConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.LeftDistance` | pxNumber | 1 | ya | ALWAYS |
| `.SurroundingRisk.LeftNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.LeftOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.SurroundingRisk.RightConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.RightDistance` | pxNumber | 1 | ya | ALWAYS |
| `.SurroundingRisk.RightNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.SurroundingRisk.RightOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.TreatyGroup` | pxDropdown | 1 | ya | IsOilGas |
| `.TSIObjectItem` | pxNumber | 1 | ya | ALWAYS |
| `.Unit` | pxInteger | 1 | ya | ALWAYS |
| `.UpdateDate` | pxDateTime | 1 | ya | ALWAYS |

#### InputDtlOtherObjectAneka

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 119 (unik 63) - read-only 119 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Aneka.Business` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isBisnisOtherObjAneka |
| `.Aneka.BusinessGolf` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.BusinessTerritorialLimit` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.ConditionExtentions` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Contestant` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.Coverage` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isCoverageOtherObjAneka |
| `.Aneka.Event` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.Exclusion` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isExclusion |
| `.Aneka.InterestInsureds` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isInterestInsuredOtherObjAneka |
| `.Aneka.LimitOfGuarantee` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.LocationInterestInsured` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.AdministrationFee` | pxInteger | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.AgreementDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.AgreementNo` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.AuctionDate` | pxDateTime | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.AuctionNo` | pxTextInput | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.BankAddress` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.BankGuaranteeNo` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.BankName` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ContractDate` | pxDateTime | 4 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ContractNo` | pxTextInput | 5 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.CustomsActivity` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.DateOfGuarantee` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.FacilityDecreeDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.FacilityDecreeNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.FacilityStatus` | (tanpa) | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.InsuredNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.Jabatan` | pxTextInput | 7 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.JenisPenjaminan` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.KBGDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.KBGNo` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.LengthPeriod` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.Location` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.MaxOfLimitOfLiability` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.NIPER` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ObligeeAddress` | pxTextArea | 6 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ObligeeName` | pxTextInput | 7 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ObligeePhone` | pxTextInput | 5 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.OCDescription` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.OtherFacilityNo` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.Pasal` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PaymentClaim` | pxTextInput | 6 | ya [IsSpreadingUW] | ALWAYS |
| `.Aneka.Obligee.PIBDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PIBNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PKDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PKDate2` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PKNo1` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PKNo2` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PlaceOfEvent` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ProjectName` | pxTextArea | 6 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ProjectValue` | pxTextInput | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PublishedIn` | pxTextArea/pxTextInput | 8 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode!=35 &&  pyWorkPage.Quotation.Busines... |
| `.Aneka.Obligee.RegisterNo` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ScheduleType` | pxRadioButtons | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.UserAssignment` | pxAutoComplete | 8 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode!=35 &&  pyWorkPage.Quotation.Busines... |
| `.Aneka.Occupation` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.OtherLOLDesc` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.PersonInsured` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Remark` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode == 18 \|\| pyWorkPage.Quotation.Busi... |
| `.Aneka.ScopeOfCover` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode == 17 \|\| pyWorkPage.Quotation.Busi... |
| `.Aneka.SpecialCondition` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.TerritorialLimits` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.TradingWarranties` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlOtherObjectAneka_FacIn

rumpun **INTI+ANEKA** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 119 (unik 65) - read-only 119 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Aneka.Business` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isBisnisOtherObjAneka |
| `.Aneka.BusinessGolf` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.BusinessTerritorialLimit` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.ConditionExtentions` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Contestant` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.Coverage` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isCoverageOtherObjAneka |
| `.Aneka.Event` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.Exclusion` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isExclusion |
| `.Aneka.InterestInsureds` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isInterestInsuredOtherObjAneka |
| `.Aneka.LimitOfGuarantee` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.LocationInterestInsured` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PublishedIn` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode!=35 &&  pyWorkPage.Quotation.Busines... |
| `.Aneka.Obligee.UserAssignment` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode!=35 &&  pyWorkPage.Quotation.Busines... |
| `.Aneka.Occupation` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.OtherLOLDesc` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.PersonInsured` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Remark` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode == 18 \|\| pyWorkPage.Quotation.Busi... |
| `.Aneka.ScopeOfCover` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode == 17 \|\| pyWorkPage.Quotation.Busi... |
| `.Aneka.SpecialCondition` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.TerritorialLimits` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.TradingWarranties` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.AdministrationFee` | pxInteger | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.AgreementDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.AgreementNo` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.AuctionDate` | pxDateTime | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.AuctionNo` | pxTextInput | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.BankAddress` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.BankGuaranteeNo` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.BankName` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.ContractDate` | pxDateTime | 4 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.ContractNo` | pxTextInput | 5 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.CustomsActivity` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.DateOfGuarantee` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.FacilityDecreeDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.FacilityDecreeNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.FacilityStatus` | (tanpa) | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.InsuredNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.Jabatan` | pxTextInput | 7 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.JenisPenjaminan` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.KBGDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.KBGNo` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.LengthPeriod` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.Location` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.MaxOfLimitOfLiability` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.NIPER` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.ObligeeAddress` | pxTextArea | 6 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.ObligeeName` | pxTextInput | 7 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.ObligeePhone` | pxTextInput | 5 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.OCDescription` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.OtherFacilityNo` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.Pasal` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.PaymentClaim` | pxTextInput | 6 | ya [IsSpreadingUW] | ALWAYS |
| `.OfferFacIn.Obligee.PIBDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.PIBNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.PKDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.PKDate2` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.PKNo1` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.PKNo2` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.PlaceOfEvent` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.ProjectName` | pxTextArea | 6 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.ProjectValue` | pxTextInput | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.PublishedIn` | pxTextArea | 7 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.RegisterNo` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.ScheduleType` | pxRadioButtons | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.UserAssignment` | pxAutoComplete | 7 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlOtherObjectAneka_FacIn_ISUW

rumpun **INTI+ANEKA** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 119 (unik 65) - read-only 119 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Aneka.Business` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isBisnisOtherObjAneka |
| `.Aneka.BusinessGolf` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.BusinessTerritorialLimit` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.ConditionExtentions` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Contestant` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.Coverage` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isCoverageOtherObjAneka |
| `.Aneka.Event` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isGolfInsurance |
| `.Aneka.Exclusion` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isExclusion |
| `.Aneka.InterestInsureds` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | isInterestInsuredOtherObjAneka |
| `.Aneka.LimitOfGuarantee` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.LocationInterestInsured` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.PublishedIn` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode!=35 &&  pyWorkPage.Quotation.Busines... |
| `.Aneka.Obligee.UserAssignment` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode!=35 &&  pyWorkPage.Quotation.Busines... |
| `.Aneka.Occupation` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.OtherLOLDesc` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.PersonInsured` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Remark` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode == 18 \|\| pyWorkPage.Quotation.Busi... |
| `.Aneka.ScopeOfCover` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode == 17 \|\| pyWorkPage.Quotation.Busi... |
| `.Aneka.SpecialCondition` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.TerritorialLimits` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.TradingWarranties` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OfferFacIn.Obligee.AdministrationFee` | pxInteger | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.AgreementDate` | pxDateTime | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.AgreementNo` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.AuctionDate` | pxDateTime | 2 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.AuctionNo` | pxTextInput | 2 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.BankAddress` | pxTextArea | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.BankGuaranteeNo` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.BankName` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.ContractDate` | pxDateTime | 4 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.ContractNo` | pxTextInput | 5 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.CustomsActivity` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.DateOfGuarantee` | pxDateTime | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.FacilityDecreeDate` | pxDateTime | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.FacilityDecreeNumber` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.FacilityStatus` | (tanpa) | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.InsuredNumber` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.Jabatan` | pxTextInput | 7 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.JenisPenjaminan` | pxDropdown | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.KBGDate` | pxDateTime | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.KBGNo` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.LengthPeriod` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.Location` | pxTextArea | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.MaxOfLimitOfLiability` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.NIPER` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.ObligeeAddress` | pxTextArea | 6 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.ObligeeName` | pxTextInput | 7 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.ObligeePhone` | pxTextInput | 5 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.OCDescription` | pxTextArea | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.OtherFacilityNo` | pxTextArea | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.Pasal` | pxDropdown | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PaymentClaim` | pxTextInput | 6 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PIBDate` | pxDateTime | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PIBNumber` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PKDate` | pxDateTime | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PKDate2` | pxDateTime | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PKNo1` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PKNo2` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PlaceOfEvent` | pxTextArea | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.ProjectName` | pxTextArea | 6 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.ProjectValue` | pxTextInput | 2 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.PublishedIn` | pxTextArea | 7 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.RegisterNo` | pxTextInput | 1 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.ScheduleType` | pxRadioButtons | 2 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.Obligee.UserAssignment` | pxAutoComplete | 7 | ya [1!=2] | ALWAYS |

#### InputDtlOtherSchedule

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-AnekaOtherSchedule` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Description` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Header` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlParticipant

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-AnekaParticipant` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Branch` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Name` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Occupation` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Person` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Sex` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ValueLoading` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlParticipantDM

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 18 (unik 15) - read-only 18 - RepeatGrid 1 - grid: .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMBMI` | pxNumber | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMHeight` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMHeirPercentage` | pxInteger | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMJoinDate` | pxDateTime | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMLeftHanded` | pxRadioButtons | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.ASMWeight` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.pyAddress` | (tanpa) | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.pyFullName` | pxTextInput | 2 | ya [IsSpreadingUW] | ALWAYS |
| `.pyPhoneNumber` | (tanpa) | 1 | ya [IsSpreadingUW] | ALWAYS |

#### InputDtlParticipantLife_FacIn

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 17 (unik 14) - read-only 14 - RepeatGrid 1 - grid: .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AccumulationName` | pxAutoComplete | 1 | tidak | 1=2 |
| `.Age` | pxInteger | 1 | tidak | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMHeight` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMHeirPercentage` | pxInteger | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | 1=2 |
| `.ASMLeftHanded` | pxRadioButtons | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMWeight` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.pyFullName` | pxTextInput | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Remark` | pxTextArea | 1 | tidak | ALWAYS |
| `OutputParam.ERRMSG2` | pxDisplayText | 1 | ya | NOTBLANK |

#### InputDtlParticipantLife_FacIn_IsUW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 15 (unik 12) - read-only 15 - RepeatGrid 1 - grid: .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Age` | pxInteger | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | ya | ALWAYS |
| `.ASMHeight` | pxNumber | 1 | ya | ALWAYS |
| `.ASMHeirPercentage` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya | 1=2 |
| `.ASMLeftHanded` | pxRadioButtons | 1 | ya | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMWeight` | pxNumber | 1 | ya | ALWAYS |
| `.pyFullName` | pxTextInput | 2 | ya | ALWAYS |
| `.Remark` | pxTextArea | 1 | ya | ALWAYS |

#### InputDtlParticipantPA_FacIn

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 17 (unik 14) - read-only 15 - RepeatGrid 1 - grid: .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMCCAmount` | pxInteger | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMClass` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | campur [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMHeight` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMHeirPercentage` | pxInteger | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMLeftHanded` | pxRadioButtons | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMWeight` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.pyFullName` | pxTextInput | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Remark` | pxTextArea | 1 | tidak | ALWAYS |

#### InputDtlParticipantPA_FacIn_IsUW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 16 (unik 13) - read-only 16 - RepeatGrid 1 - grid: .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMCCAmount` | pxInteger | 1 | ya [1=1] | ALWAYS |
| `.ASMClass` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya [1=1] | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | ya | ALWAYS |
| `.ASMHeight` | pxNumber | 1 | ya | ALWAYS |
| `.ASMHeirPercentage` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ASMLeftHanded` | pxRadioButtons | 1 | ya | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMWeight` | pxNumber | 1 | ya | ALWAYS |
| `.pyFullName` | pxTextInput | 2 | ya | ALWAYS |

#### InputDtlParticipantTravel

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 12 (unik 9) - read-only 11 - RepeatGrid 1 - grid: .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | 1=2 |
| `.ASMHeirPercentage` | pxInteger | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.pyFullName` | pxTextInput | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Remark` | pxTextArea | 1 | tidak | ALWAYS |

#### InputDtlParticipantTravel_GISFW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 12 (unik 9) - read-only 11 - RepeatGrid 2 - grid: .ASMCoverage, .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | ya | 1=2 |
| `.ASMHeirPercentage` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya | ALWAYS |
| `.CoverageNote` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyFullName` | pxDropdown/pxTextInput | 2 | ya | ALWAYS |

#### InputDtlParticipantTravel_IsUW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 12 (unik 9) - read-only 12 - RepeatGrid 1 - grid: .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | ya | 1=2 |
| `.ASMHeirPercentage` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya | ALWAYS |
| `.pyFullName` | pxTextInput | 2 | ya | ALWAYS |
| `.Remark` | pxTextArea | 1 | ya | ALWAYS |

#### InputDtlPayment_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Installment` - field terikat 15 (unik 15) - read-only 15 - RepeatGrid 1 - grid: .Policy.Payment.ListInstallment

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AdminFee` | pxNumber | 1 | ya | IsHealthIndividu |
| `.BrokerageFee` | pxNumber | 1 | ya | ALWAYS |
| `.Deduction2` | pxNumber | 1 | ya | ALWAYS |
| `.Discount` | pxNumber | 1 | ya | pyWorkPage.Quotation.BusinessCode=='10101' |
| `.DueDate` | pxDateTime | 1 | ya [1 = 2] | ALWAYS |
| `.InstallmentNo` | pxInteger | 1 | ya | ALWAYS |
| `.InstallmentPercentage` | pxNumber | 1 | ya [pyPortal.IsShowOpenPolicyMarine == true \|\| pyW...] | ALWAYS |
| `.PaymentTotal` | pxNumber | 1 | ya | ALWAYS |
| `.PPh` | pxNumber | 1 | ya | ALWAYS |
| `.PPN` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.RICommision` | pxTextInput | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |
| `pyWorkPage.PPH_Note` | (tanpa) | 1 | ya | ALWAYS |
| `pyWorkPage.PPN_Note` | (tanpa) | 1 | ya | ALWAYS |

#### InputDtlPayment_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Installment` - field terikat 13 (unik 13) - read-only 13 - RepeatGrid 1 - grid: .Policy.Payment.ListInstallment

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BrokerageFee` | pxNumber | 1 | ya | ALWAYS |
| `.DueDate` | pxDateTime | 1 | ya [OpenWpc.CARI9!=1\|\|.FlagInsPayed==1] | ALWAYS |
| `.InstallmentNo` | pxInteger | 1 | ya | ALWAYS |
| `.InstallmentPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.PaymentTotal` | pxNumber | 1 | ya | ALWAYS |
| `.PPh` | pxNumber | 1 | ya | ALWAYS |
| `.PPN` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.RICommision` | pxTextInput | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |
| `pyWorkPage.PaymentList(1).ListInstallment(1).Premium` | pxNumber | 1 | ya | pyWorkPage.Quotation.BusinessCode=='10101' |
| `pyWorkPage.PPH_Note` | (tanpa) | 1 | ya | ALWAYS |
| `pyWorkPage.PPN_Note` | (tanpa) | 1 | ya | ALWAYS |

#### InputDtlPayment_FacInLife

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Installment` - field terikat 8 (unik 8) - read-only 7 - RepeatGrid 1 - grid: .Policy.Payment.ListInstallment

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Commission` | pxNumber | 1 | ya | pyWorkPage.Quotation.BusinessCode=='10101' |
| `.DueDate` | pxDateTime | 1 | tidak | NOTBLANK |
| `.InstallmentNo` | pxInteger | 1 | ya | ALWAYS |
| `.PaymentTotal` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.RateEM` | pxTextInput | 1 | ya | ALWAYS |
| `.RateLife` | pxTextInput | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |

#### InputDtlPayment_FacInLife_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Installment` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 1 - grid: .Policy.Payment.ListInstallment

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Commission` | pxNumber | 1 | ya | pyWorkPage.Quotation.BusinessCode=='10101' |
| `.DueDate` | pxDateTime | 1 | ya | NOTBLANK |
| `.InstallmentNo` | pxInteger | 1 | ya | ALWAYS |
| `.PaymentTotal` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.RateEM` | pxTextInput | 1 | ya | ALWAYS |
| `.RateLife` | pxTextInput | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |

#### InputDtlPaymentFacRetro

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Policy.Payment.Commision` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3 |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3 |
| `.Policy.Payment.NetPremium` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 1 | ya | pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3 |

#### InputDtlPaymentFacRetro_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Policy.Payment.Commision` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3 |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3 |
| `.Policy.Payment.NetPremium` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 1 | ya | pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3 |

#### InputDtlSpreadingCommisionCoverage_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-SpreadingRisk` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 1 - grid: .SpreadingList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClaimEstimation` | pxNumber | 1 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyType` | pxDropdown | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 1 | ya | ALWAYS |

#### InputDtlSpreadingCoverage_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-SpreadingRisk` - field terikat 15 (unik 7) - read-only 15 - RepeatGrid 3 - grid: .SpreadingList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClaimEstimation` | pxNumber | 1 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 3 | ya [.FlagDelete = 1 \|\| (OperatorID.pyUserIdentifie...] | ALWAYS |
| `.SharePercentage` | pxNumber | 3 | ya [.FlagDelete = 1] | ALWAYS |
| `.TreatyType` | pxDropdown | 3 | ya [.FlagDelete = 1] | ALWAYS |
| `.TSIGrossSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 3 | ya [.FlagDelete = 1 \|\| (OperatorID.pyUserIdentifie...] | ALWAYS |
| `OutputParam.ERRMSG6` | pxDisplayText | 1 | ya | OutputParam.ERRMSG6!='' && pyWorkPage.OfferFacIn.QuotationData.IsGroup... |

#### InputDtlSpreadingCoverage_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-SpreadingRisk` - field terikat 14 (unik 6) - read-only 14 - RepeatGrid 3 - grid: .SpreadingList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClaimEstimation` | pxNumber | 1 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 3 | ya [.FlagDelete = 1 \|\| OperatorID.pyUserIdentifier...] | ALWAYS |
| `.SharePercentage` | pxNumber | 3 | ya [.FlagDelete = 1] | ALWAYS |
| `.TreatyType` | pxDropdown | 3 | ya [.FlagDelete = 1] | ALWAYS |
| `.TSIGrossSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 3 | ya [.FlagDelete = 1 \|\| OperatorID.pyUserIdentifier...] | ALWAYS |

#### InputDtlSpreadingPA_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSICeding` | pxNumber | 1 | ya | 1==2 |
| `.TSILiability` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### InputDtlSpreadingPA_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSILiability` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### InputDtlSpreadingTravel_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSICeding` | pxNumber | 1 | ya | 1==2 |
| `.TSILiability` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### InputDtlTrading

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Trading` - field terikat 10 (unik 10) - read-only 10 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.FromRute` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.PartshipmentNote` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.PartshipmentStatus` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ToRute` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TradingID` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TradingNote2` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | .TradingID ='10007' |
| `.Tranship` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TranshipmentNote` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TranshipmentStatus` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Warehouse` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputDtlTrading_IsUW

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Trading` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.FromRute` | pxTextInput | 1 | ya | ALWAYS |
| `.PartshipmentNote` | pxTextInput | 1 | ya | ALWAYS |
| `.PartshipmentStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.ToRute` | pxTextInput | 1 | ya | ALWAYS |
| `.TradingNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Tranship` | pxTextInput | 1 | ya | ALWAYS |
| `.TranshipmentNote` | pxTextInput | 1 | ya | ALWAYS |
| `.TranshipmentStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.Warehouse` | pxTextInput | 1 | ya | ALWAYS |

#### InputEndorsement

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.IsShowDetail` | pxCheckbox | 1 | tidak | ALWAYS |

#### InputEndorsementDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 46 (unik 30) - read-only 28 - RepeatGrid 7 - grid: .OfferFacIn.CedingCedantList, .OfferFacIn.CurrencyList, SpreadingList.pxResults, TotalSpreadAll.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CedingCoName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Name` | (tanpa) | 4 | campur | ALWAYS |
| `.OfferFacIn.CedingRetention` | pxNumber | 2 | campur [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | 1=2 ; IsLife |
| `.OfferFacIn.CedingRetentionNominal` | pxNumber | 1 | tidak [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | 1=2 |
| `.OfferFacIn.PercentShare` | pxNumber | 3 | campur [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | 1=2 |
| `.OfferFacIn.PPnCheck` | pxCheckbox | 1 | tidak | ALWAYS |
| `.OfferFacIn.QuotationData.SobName` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.ShareCedantType` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.OfferFacIn.ShareRNML` | pxNumber | 1 | ya [IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSp...] | 1=2 |
| `.OfferFacIn.TotalPremiNusaRe` | pxNumber/pxTextInput | 3 | ya | 1=2 |
| `.OfferFacIn.TotalTSINusaRe` | pxNumber/pxTextInput | 2 | campur | 1 = 2 ; 1==2 |
| `.OfferFacIn.TotalTSINusaReSpreading` | pxNumber/pxTextInput | 3 | ya | 1=2 |
| `.OfferFacIn.TotalTSITopRisk` | pxNumber | 1 | ya | IsFire |
| `.OldTSI` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.BrokerageFee` | pxNumber | 1 | tidak | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 1 | tidak | ALWAYS |
| `.Policy.Payment.Deduction2` | pxNumber | 1 | tidak | ALWAYS |
| `.Policy.Payment.PctBrokerageFee` | pxNumber | 1 | tidak | .CalcBrokerFee = true |
| `.Policy.Payment.PPh` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.PPN` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.RICommision` | pxNumber | 1 | tidak [IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSp...] | !IsUW |
| `.PPH_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.PPN_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.ShareCeding` | pxNumber | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | campur [IsUW] | ALWAYS |
| `.TreatyType` | pxDropdown | 2 | tidak [IsUW] | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |
| `SearchLife.CARI7` | pxAutoComplete | 1 | tidak | 1=2 |

#### InputEndorsementDtl_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 34 (unik 27) - read-only 26 - RepeatGrid 6 - grid: .OfferFacIn.CedingCedantList, .OfferFacIn.CurrencyList, SpreadingList.pxResults, TotalSpreadAll.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CedingCoName` | pxDisplayText | 1 | ya | ALWAYS |
| `.FlagSpreading` | pxDropdown | 1 | tidak [IsUW] | ALWAYS |
| `.Name` | (tanpa) | 4 | tidak | ALWAYS |
| `.OfferFacIn.CedingRetention` | pxNumber | 1 | ya | IsLife |
| `.OfferFacIn.FlagSpreading` | pxTextInput | 1 | ya | NOTBLANK |
| `.OfferFacIn.PercentShare` | pxNumber | 2 | ya [1!=2] | ALWAYS |
| `.OfferFacIn.PPnCheck` | pxCheckbox | 1 | tidak | ALWAYS |
| `.OfferFacIn.QuotationData.SobName` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.ShareCedantType` | pxRadioButtons | 1 | ya | ALWAYS |
| `.OfferFacIn.TotalPremiNusaRe` | pxNumber | 1 | ya | ALWAYS |
| `.OfferFacIn.TotalTSINusaReSpreading` | pxNumber | 1 | ya | ALWAYS |
| `.OldTSI` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.BrokerageFee` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Deduction2` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.PctBrokerageFee` | pxNumber | 1 | ya | !IsUW |
| `.Policy.Payment.PPh` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.PPN` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.RICommision` | pxNumber | 1 | ya | !IsUW |
| `.PPH_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.PPN_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.ShareCeding` | pxNumber | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 1 | tidak [IsUW] | ALWAYS |
| `.TreatyType` | pxDropdown | 1 | tidak [IsUW] | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |

#### InputFacOffer

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 29 (unik 20) - read-only 22 - RepeatGrid 6 - grid: .OfferFacIn.FacRetroList, CurrencyListRetroAll.pxResults, GrandTotalListRetroAll.pxResults, TempViewSuggest.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Approval` | pxDropdown | 1 | ya | ALWAYS |
| `.CommentSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.CurrentYear` | (tanpa) | 1 | tidak | ALWAYS |
| `.DateTransfer` | pxDateTime | 1 | ya | ALWAYS |
| `.EmailTypeRetro` | pxRadioButtons | 1 | tidak | .FlagOnGoingPolicy!=1\|\|(.FlagOnGoingPolicy==1&&pyWorkPage.OfferFacIn... |
| `.EndPeriod` | pxDateTime | 2 | ya | ALWAYS |
| `.IsCedingConfirm` | pxDisplayText | 1 | ya | ALWAYS |
| `.Name` | pxTextInput | 3 | campur | ALWAYS |
| `.OfferFacIn.Comment` | pxTextArea | 1 | tidak | ALWAYS |
| `.OurRef` | pxTextInput | 2 | ya | ALWAYS |
| `.PICSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 2 | campur | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.NetPremium` | pxNumber | 2 | campur | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 2 | campur | ALWAYS |
| `.ReinsurerName` | (tanpa) | 2 | ya | ALWAYS |
| `.StartPeriod` | pxDateTime | 2 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 1 | ya | ALWAYS |
| `pyWorkPage.FlagErrorKonversi` | (tanpa) | 1 | ya [1=1] | NOTBLANK |

#### InputFacOffer_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 21 (unik 12) - read-only 21 - RepeatGrid 5 - grid: .OfferFacIn.FacRetroList, CurrencyListRetroAll.pxResults, GrandTotalListRetroAll.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CurrentYear` | (tanpa) | 1 | ya | ALWAYS |
| `.EndPeriod` | pxDateTime | 2 | ya | ALWAYS |
| `.Name` | pxTextInput | 3 | ya | ALWAYS |
| `.OurRef` | pxTextInput | 2 | ya | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.NetPremium` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.ReinsurerName` | (tanpa) | 2 | ya | ALWAYS |
| `.StartPeriod` | pxDateTime | 2 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 1 | ya | ALWAYS |

#### InputFEA

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-FEA` - field terikat 11 (unik 11) - read-only 11 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DistributionStatus` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.FEAID` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.FEAInside` | pxInteger | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.FEAOutside` | pxInteger | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.MaintenanceStatus` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.PrivateFireBrigade` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PrivateTruckBrigade` | pxTextInput | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.RemarkFEA` | pxTextArea | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.TeamSOPRiskManagement` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.TeamSOPSafety` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.UpdateDate` | pxDateTime | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |

#### InputFEA_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.APAR` | pxNumber | 1 | ya | ALWAYS |
| `.DataFEA.PrivateFireBrigade` | pxDropdown | 1 | ya | ALWAYS |
| `.DataFEA.PrivateTruckBrigade` | pxNumber | 1 | ya | ALWAYS |
| `.DataFEA.TeamSOPRiskManagement` | pxDropdown | 1 | ya | ALWAYS |
| `.DataFEA.TeamSOPSafety` | pxDropdown | 1 | ya | ALWAYS |
| `.Hydrant` | pxNumber | 1 | ya | ALWAYS |
| `.InfoFEA` | pxTextArea | 1 | ya | ALWAYS |
| `.SmokeDetector` | pxNumber | 1 | ya | ALWAYS |
| `.Sprinkler` | pxNumber | 1 | ya | ALWAYS |

#### InputHistoricalSurveyReportDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Quotation` - field terikat 3 (unik 3) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DateofSurvey` | pxDateTime | 1 | tidak | ALWAYS |
| `.Remarks` | pxDropdown | 1 | tidak | ALWAYS |
| `.SurveyedBy` | pxTextInput | 1 | tidak | ALWAYS |

#### InputHistoricalSurveyReportDtlUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Quotation` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DateofSurvey` | pxDateTime | 1 | ya | ALWAYS |
| `.Remarks` | pxDropdown | 1 | ya | ALWAYS |
| `.SurveyedBy` | pxTextInput | 1 | ya | ALWAYS |

#### InputInwardFacultative

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.IsShowDetail` | pxCheckbox | 1 | tidak | ALWAYS |

#### InputInwardFacultative_2

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.IsShowDetail` | pxCheckbox | 1 | tidak | ALWAYS |
| `.IsShowFacRetro` | pxCheckbox | 1 | tidak | ALWAYS |

#### InputInwardFacultativeDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 63 (unik 36) - read-only 37 - RepeatGrid 10 - grid: .OfferFacIn.CedingCedantList, .OfferFacIn.CurrencyList, SpreadingList.pxResults, TotalSpreadAll.pxResults, TotalSpreadingCurrency.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CedingCoName` | pxDisplayText | 1 | ya | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.Name` | (tanpa) | 4 | tidak | ALWAYS |
| `.OfferFacIn.CedingRetention` | pxNumber | 2 | campur [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | 1=2 ; IsLife |
| `.OfferFacIn.CedingRetentionNominal` | pxNumber | 1 | tidak [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | 1=2 |
| `.OfferFacIn.PercentShare` | pxNumber | 3 | campur [IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSp...] | 1=2 |
| `.OfferFacIn.PPnCheck` | pxCheckbox | 1 | tidak | ALWAYS |
| `.OfferFacIn.QuotationData.SobName` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.ShareCedantType` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.OfferFacIn.ShareRNML` | pxNumber | 1 | ya [IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSp...] | 1=2 |
| `.OfferFacIn.TotalPremiNusaRe` | pxNumber/pxTextInput | 4 | ya | 1=2 |
| `.OfferFacIn.TotalRIComNusaRe` | pxTextInput | 1 | ya | 1=2 |
| `.OfferFacIn.TotalTSINusaRe` | pxNumber | 2 | ya | 1 = 2 |
| `.OfferFacIn.TotalTSINusaReSpreading` | pxNumber/pxTextInput | 3 | ya | 1=2 |
| `.OfferFacIn.TotalTSITopRisk` | pxNumber | 2 | ya | 1=2 ; IsFire |
| `.OldTSI` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.BrokerageFee` | pxNumber | 1 | tidak | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 1 | tidak | ALWAYS |
| `.Policy.Payment.Deduction2` | pxNumber | 1 | tidak | ALWAYS |
| `.Policy.Payment.PctBrokerageFee` | pxNumber | 1 | tidak [IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSp...] | !IsUW |
| `.Policy.Payment.PctDeduction2` | pxNumber | 1 | tidak [IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSp...] | never |
| `.Policy.Payment.PPh` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.PPN` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.RICommision` | pxNumber | 1 | tidak [IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSp...] | !IsUW |
| `.PPH_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.PPN_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 3 | ya | ALWAYS |
| `.ShareCeding` | pxNumber | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | tidak [IsUW] | ALWAYS |
| `.TreatyName` | (tanpa) | 3 | tidak | ALWAYS |
| `.TreatyType` | pxDropdown | 5 | tidak [IsUW] | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 3 | ya | ALWAYS |
| `SearchLife.CARI7` | pxAutoComplete | 1 | tidak | 1=2 |

#### InputInwardFacultativeDtl_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 46 (unik 31) - read-only 36 - RepeatGrid 8 - grid: .OfferFacIn.CedingCedantList, .OfferFacIn.CurrencyList, SpreadingList.pxResults, TotalSpreadAll.pxResults, TotalSpreadingCurrency.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CedingCoName` | pxDisplayText | 1 | ya | ALWAYS |
| `.ClaimEstimation` | pxNumber | 2 | ya | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.FlagSpreading` | pxDropdown | 1 | tidak [IsUW] | ALWAYS |
| `.Name` | (tanpa) | 4 | tidak | ALWAYS |
| `.OfferFacIn.FlagSpreading` | pxTextInput | 1 | ya | NOTBLANK |
| `.OfferFacIn.PercentShare` | pxNumber | 2 | ya [!IsUWSpread] | ALWAYS |
| `.OfferFacIn.PPnCheck` | pxCheckbox | 1 | tidak | ALWAYS |
| `.OfferFacIn.QuotationData.SobName` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.ShareCedantType` | pxRadioButtons | 1 | ya | ALWAYS |
| `.OfferFacIn.TotalPremiNusaRe` | pxNumber | 1 | ya | ALWAYS |
| `.OfferFacIn.TotalTSINusaReSpreading` | pxNumber | 1 | ya | ALWAYS |
| `.OldTSI` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.BrokerageFee` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Deduction2` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.PctBrokerageFee` | pxNumber | 1 | ya | !IsUW |
| `.Policy.Payment.PPh` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.PPN` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.RICommision` | pxNumber | 1 | ya | !IsUW |
| `.PPH_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.PPN_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.ShareCeding` | pxNumber | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 3 | campur [IsUW] | OutputParam.ERRMSG6=='' |
| `.TreatyName` | (tanpa) | 2 | tidak | ALWAYS |
| `.TreatyType` | pxDropdown | 3 | campur [IsUW] | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 2 | ya | ALWAYS |

#### InputInwardFacultativeSuggest

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-SuggestList` - field terikat 7 (unik 7) - read-only 6 - RepeatGrid 1 - grid: .ViewSuggest

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Approval` | pxDropdown | 1 | ya | ALWAYS |
| `.Comment` | pxTextArea | 1 | tidak | 1==2 |
| `.CommentSuggest` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Date` | pxDateTime | 1 | ya | ALWAYS |
| `.DateSuggest` | pxDateTime | 1 | ya | NOTBLANK |
| `.PIC` | pxDisplayText | 1 | ya | ALWAYS |
| `.PICSuggest` | pxDisplayText | 1 | ya | NOTBLANK |

#### InputLossExperiance

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-CauseOfLoss` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.CoinsData.CoinsName` | pxAutoComplete | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [IsGroupFac \|\| isGroupUwFac \|\| IsSpreadingUW] | .Quotation.BusinessCode='04' |
| `.DateOfLoss` | pxDateTime | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.Detail` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.LossObject` | pxTextInput | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |

#### InputLossRecord_Sec

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-CauseOfLoss` - field terikat 8 (unik 8) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CauseOfLoss` | pxTextInput | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `.ClassOfBusiness` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | .Quotation.BusinessCode='04' |
| `.DOL` | pxDateTime | 1 | tidak | ALWAYS |
| `.LOSS_AMOUNT` | pxNumber | 1 | tidak | ALWAYS |
| `.LOSS_AMOUNT_IDR` | pxNumber | 1 | tidak | ALWAYS |
| `.UWYEAR` | pxTextInput | 1 | tidak [.FlagDelete = 1 \|\| IsClaim \|\| pyPortal.IsSho...] | ALWAYS |
| `pyWorkPage.OfferFacIn.QuotationData.InsuredName` | (tanpa) | 1 | ya | ALWAYS |

#### InputMarketingOfficer

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 7 (unik 7) - read-only 5 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputMarketingOfficer.A_Leader` | pxCheckbox | 1 | tidak | ALWAYS |
| `InputMarketingOfficer.BranchDetailID` | pxDropdown | 1 | ya [InputMarketingOfficer.A_Leader = false \|\| Inpu...] | InputMarketingOfficer.BranchParent != '' |
| `InputMarketingOfficer.BranchParent` | pxDropdown | 1 | ya [InputMarketingOfficer.A_Leader = false \|\| Inpu...] | ALWAYS |
| `InputMarketingOfficer.ClientName` | pxAutoComplete | 1 | ya [InputMarketingOfficer.ID != ""] | ALWAYS |
| `InputMarketingOfficer.ID` | pxTextInput | 1 | ya | ALWAYS |
| `InputMarketingOfficer.MOLeader` | pxAutoComplete | 1 | ya [InputMarketingOfficer.ID != ""] | InputMarketingOfficer.A_Leader == false |
| `InputMarketingOfficer.MOStatus` | pxRadioButtons | 1 | tidak | ALWAYS |

#### InputNation

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 3 (unik 3) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputNation.ID` | pxTextInput | 1 | ya | 1==2 |
| `InputNation.NationInitial` | pxTextArea | 1 | tidak | ALWAYS |
| `InputNation.Note` | pxTextArea | 1 | tidak | ALWAYS |

#### InputObjectDtlFire

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### InputObjectDtlFire_GCNM

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Property` - field terikat 80 (unik 75) - read-only 70 - RepeatGrid 6 - grid: .CoverageList, .FEAList, .ListCauseOfLoss, .OccupationList, .PropertyItemList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AlmRiskID` | pxTextInput | 1 | ya | ALWAYS |
| `.Amount` | pxNumber | 1 | tidak | ALWAYS |
| `.BuildingConstruction.FloorType` | pxDropdown | 1 | ya | ALWAYS |
| `.BuildingConstruction.NumberOfFloor` | pxNumber | 1 | ya | ALWAYS |
| `.BuildingConstruction.OthersType` | pxTextInput | 1 | ya | ALWAYS |
| `.BuildingConstruction.PartitionType` | pxTextInput | 1 | ya | ALWAYS |
| `.BuildingConstruction.RoofType` | pxDropdown | 1 | ya | ALWAYS |
| `.BuildingConstruction.SupportWallType` | pxTextInput | 1 | ya | ALWAYS |
| `.BuildingConstruction.WallType` | pxDropdown | 1 | ya | ALWAYS |
| `.BuildingNo` | pxTextInput | 1 | ya | ALWAYS |
| `.CoinsData.CoinsName` | pxTextInput | 1 | tidak | ALWAYS |
| `.Construction` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Coverage` | (tanpa) | 1 | tidak | ALWAYS |
| `.CoverageNote` | (tanpa) | 1 | tidak | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.DateOfLoss` | pxDateTime | 1 | tidak | ALWAYS |
| `.Detail` | pxTextInput | 1 | tidak | ALWAYS |
| `.EML` | pxNumber | 1 | ya | ALWAYS |
| `.FEAID` | pxDropdown | 1 | tidak | ALWAYS |
| `.FEAInside` | pxInteger | 1 | ya | ALWAYS |
| `.FEAOutside` | pxInteger | 1 | ya | ALWAYS |
| `.FirstLossLimitFac` | pxTextInput | 1 | ya | ALWAYS |
| `.IsFlammableItemFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsHotWorkProcessFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsMaterialDamage` | pxCheckbox | 1 | tidak | ALWAYS |
| `.IsProductionProcessFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 1 | tidak | IsFac |
| `.IsUsedFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 2 | ya | ALWAYS |
| `.LossObject` | pxTextInput | 1 | tidak | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | .ObjectType = 'Lainnya' |
| `.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectType` | pxDropdown | 1 | ya | ALWAYS |
| `.OccupationId` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.Ownership` | pxDropdown | 1 | ya | ALWAYS |
| `.PctOfContribution` | pxTextInput | 1 | ya | ALWAYS |
| `.PML` | pxNumber | 1 | ya | ALWAYS |
| `.PropertyYear` | pxInteger | 2 | ya | ALWAYS |
| `.pyCity` | pxTextInput | 1 | ya | ALWAYS |
| `.pyCountryName` | pxTextInput | 1 | ya | ALWAYS |
| `.RIComm` | pxNumber | 1 | ya | ALWAYS |
| `.RiskCategory` | pxDropdown | 1 | ya | IsFac |
| `.RiskLocation.ASMAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.RiskLocation.ASMCity` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMDistrict` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMRW` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.RoadName` | pxTextInput | 1 | ya | ALWAYS |
| `.RoadType` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.BackConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.BackDistance` | pxNumber | 1 | ya | ALWAYS |
| `.SurroundingRisk.BackNote` | pxTextInput | 1 | ya | ALWAYS |
| `.SurroundingRisk.BackOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.SurroundingRisk.FloodArea` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.FloodAreaStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.FrontConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.FrontDistance` | pxNumber | 1 | ya | ALWAYS |
| `.SurroundingRisk.FrontNote` | pxTextInput | 1 | ya | ALWAYS |
| `.SurroundingRisk.FrontOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.SurroundingRisk.HousekeepingRemark` | pxTextArea | 1 | ya | ALWAYS |
| `.SurroundingRisk.HousekeepingStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.LeftConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.LeftDistance` | pxNumber | 1 | ya | ALWAYS |
| `.SurroundingRisk.LeftNote` | pxTextInput | 1 | ya | ALWAYS |
| `.SurroundingRisk.LeftOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.SurroundingRisk.RightConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.SurroundingRisk.RightDistance` | pxNumber | 1 | ya | ALWAYS |
| `.SurroundingRisk.RightNote` | pxTextInput | 1 | ya | ALWAYS |
| `.SurroundingRisk.RightOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.TreatyGroup` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac] | IsOilGas |
| `.TSIObjectItem` | pxNumber | 2 | ya | ALWAYS |
| `.Unit` | pxInteger | 2 | ya | ALWAYS |
| `.UpdateDate` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.Policy.Currency` | (tanpa) | 2 | ya | ALWAYS |

#### InputObjectDtlFire_IsUW

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### InputObjectSubContract_FacIn

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-SubContract` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 2 - grid: .OfferFacIn.LocationList, .SubContractList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.SubContractID` | pxTextInput | 1 | ya | ALWAYS |
| `.SubContractName` | pxTextInput | 1 | ya | ALWAYS |

#### InputObjectSubContract_FacIn_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-SubContract` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 2 - grid: .OfferFacIn.LocationList, .SubContractList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.SubContractID` | pxTextInput | 1 | ya | ALWAYS |
| `.SubContractName` | pxTextInput | 1 | ya | ALWAYS |

#### InputObjItem

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 12 (unik 12) - read-only 7 - RepeatGrid 1 - grid: .StockAdjustmentList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Condition` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.IsAdjustableFlag` | pxCheckbox | 1 | tidak | IsTypeStock |
| `.ItemType` | pxAutoComplete | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.MonthName` | pxDropdown | 1 | tidak | ALWAYS |
| `.PercentageAdjustment` | pxDropdown | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | IsAdjustableFlag |
| `.PropertiItemNote` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.PropertyYear` | pxInteger | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.TSI` | pxInteger | 1 | tidak | ALWAYS |
| `.TSIAdjustment` | pxNumber | 1 | tidak | ALWAYS |
| `.TSIObjectItem` | pxNumber | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.Unit` | pxInteger | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.Year` | pxTextInput | 1 | tidak | ALWAYS |

#### InputOfferFacInLossRatio

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.LossRatio1YearAmount` | pxNumber | 1 | ya | ALWAYS |
| `.LossRatio1YearPercent` | pxNumber | 1 | ya | ALWAYS |
| `.LossRatio35YearAmount` | pxNumber | 1 | ya | ALWAYS |
| `.LossRatio35YearPercent` | pxNumber | 1 | ya | ALWAYS |

#### InputOkupasi

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 5 (unik 5) - read-only 4 - RepeatGrid 1 - grid: .SubjectToList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Answer` | pxRadioButtons | 1 | ya [IsFlagDelete] | ALWAYS |
| `.IsUsedFlag` | pxCheckbox | 1 | tidak | .IsUsedFlag=='true' |
| `.OccupationId` | pxAutoComplete | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.OccupationName` | pxTextArea | 1 | ya [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW] | ALWAYS |
| `.Question` | pxTextInput | 1 | ya [IsFlagDelete \|\| IsSpreadingUW] | !IsNotFlagDelete |

#### InputOkupasiAneka_FacIn

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 19 (unik 12) - read-only 5 - RepeatGrid 3 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | tidak | ALWAYS |
| `.LimitOfLiability` | pxTextInput | 2 | tidak | IsObjectSectionAneka |
| `.ObjectName` | pxDisplayText | 3 | ya | ALWAYS |
| `.OccupationId` | pxAutoComplete | 1 | ya [pyWorkPage.IsOldData = 1 \|\| IsSpreadingUW] | 1=2 |
| `.OccupationName` | pxTextArea | 1 | ya | 1=2 |
| `.Quantity` | pxNumber | 1 | tidak | IsObjectWithQuantityYear \|\| IsHE |
| `.Section` | pxDropdown | 2 | tidak | IsObjectSectionAneka |
| `.TSI` | (tanpa) | 1 | tidak | ALWAYS |
| `.UnitQuantity` | pxTextInput | 1 | tidak | IsObjectWithQuantityYear \|\| IsHE |
| `.VehicleHE.ChassisNumber` | pxNumber | 2 | tidak | IsHE |
| `.VehicleHE.EngineNumber` | pxNumber/pxTextInput | 2 | tidak | IsHE |
| `.Year` | pxTextInput | 2 | tidak | IsObjectWithQuantityYear \|\| IsHE |

#### InputOkupasiAneka_FacIn_IsUW

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 19 (unik 12) - read-only 19 - RepeatGrid 3 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.LimitOfLiability` | pxTextInput | 2 | ya | IsObjectSectionAneka |
| `.ObjectName` | pxTextInput | 3 | ya | ALWAYS |
| `.OccupationId` | pxAutoComplete | 1 | ya | 1=2 |
| `.OccupationName` | pxTextArea | 1 | ya | 1=2 |
| `.Quantity` | pxNumber | 1 | ya | IsObjectWithQuantityYear \|\| IsHE |
| `.Section` | pxDropdown | 2 | ya | IsObjectSectionAneka |
| `.TSI` | (tanpa) | 1 | ya | ALWAYS |
| `.UnitQuantity` | pxTextInput | 1 | ya | IsObjectWithQuantityYear \|\| IsHE |
| `.VehicleHE.ChassisNumber` | pxNumber | 2 | ya | IsHE |
| `.VehicleHE.EngineNumber` | pxNumber/pxTextInput | 2 | ya | IsHE |
| `.Year` | pxTextInput | 2 | ya | IsObjectWithQuantityYear \|\| IsHE |

#### InputOkupasiAneka_GCNM

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 7 (unik 7) - read-only 7 - RepeatGrid 1 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationId` | pxAutoComplete | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.Quantity` | (tanpa) | 1 | ya | ALWAYS |
| `.SumDiscountObjectAneka` | (tanpa) | 1 | ya | ALWAYS |
| `.SumPremiObjectAneka` | (tanpa) | 1 | ya | ALWAYS |
| `.SumTotalTSI` | pxTextInput | 1 | ya | ALWAYS |

#### InputOkupasiHull_GCNM

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 5 (unik 5) - read-only 3 - RepeatGrid 1 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ObjectName` | pxTextInput | 1 | tidak | ALWAYS |
| `.OccupationId` | pxAutoComplete | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.pySelected` | pxCheckbox | 1 | tidak | ALWAYS |
| `.SumTotalTSI` | pxTextInput | 1 | ya | ALWAYS |

#### InputOtherObjectAneka_FacIn

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### InputOtherObjectAneka_FacIn_ISUW

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### InputOtherObjectAneka_UW

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### InputOtherSchedule

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .AnekaAnotherScheduleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AnekaAnotherSchedule.OtherAnnualPremiumDesc` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaAnotherSchedule.RetroactiveDate` | pxDateTime | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaAnotherSchedule.TerritorialLimit` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.AnekaAnotherSchedule.Yurisdiction` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### InputOutgo

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Outgo` - field terikat 7 (unik 6) - read-only 6 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OutgoAmount` | pxNumber | 1 | ya [IsSpreadingUW] | .TypeOfNominal='false' |
| `.PercentOutgo` | pxTextInput | 1 | ya [IsSpreadingUW] | .TypeOfNominal='true' |
| `.ReceiverOutgo` | pxAutoComplete | 1 | ya [IsSpreadingUW] | .TypeOfOutgo != '1' |
| `.SourceBizName` | pxAutoComplete | 2 | ya [IsSpreadingUW] | .TypeOfOutgo == '1' ; .TypeOfOutgo == '4' |
| `.TypeOfNominal` | pxCheckbox | 1 | tidak | ALWAYS |
| `.TypeOfOutgo` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |

#### InputPerson

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 10 (unik 10) - read-only 2 - RepeatGrid 1 - grid: .PersonList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMCoverage(1).Loading` | (tanpa) | 1 | ya | NOTBLANK |
| `.ASMCoverage(1).PremiRp` | (tanpa) | 1 | tidak | ALWAYS |
| `.ASMCoverage(1).Premium` | (tanpa) | 1 | tidak | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 1 | tidak | ALWAYS |
| `.ASMGender` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMMaritalStatus` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMStatusCoverage` | pxDisplayText | 1 | ya | pyWorkPage.ProposalPosition = '2' \|\| pyWorkPage.ProposalPosition = '... |
| `.pyFullName` | (tanpa) | 1 | tidak | ALWAYS |

#### InputProvince

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 3 (unik 3) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputProvince.ID` | pxTextInput | 1 | ya | 1==2 |
| `InputProvince.NationName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputProvince.Note` | pxTextArea | 1 | tidak | ALWAYS |

#### InputRiskAddress

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 11 (unik 11) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputRiskAddress.Address` | pxTextArea | 1 | tidak | ALWAYS |
| `InputRiskAddress.CityName` | pxAutoComplete | 1 | tidak | InputRiskAddress.ProvinceName != '' |
| `InputRiskAddress.DistrictName` | pxAutoComplete | 1 | tidak | InputRiskAddress.CityName != '' |
| `InputRiskAddress.ID` | pxTextInput | 1 | ya | never |
| `InputRiskAddress.NationName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputRiskAddress.PostalCode` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputRiskAddress.ProvinceName` | pxAutoComplete | 1 | tidak | InputRiskAddress.NationName != '' |
| `InputRiskAddress.StatusTransCenter` | pxRadioButtons | 1 | tidak | never |
| `InputRiskAddress.TerritoryName` | pxAutoComplete | 1 | tidak | InputRiskAddress.DistrictName != '' |
| `InputRiskAddress.Title` | pxDropdown | 1 | tidak | ALWAYS |
| `InputRiskAddress.Type` | pxTextInput | 1 | tidak | never |

#### InputRW

rumpun **INTI** - kelas dominan `@baseclass` - field terikat 9 (unik 9) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputRW.AssessmentZone` | pxTextInput | 1 | tidak | ALWAYS |
| `InputRW.CITYNAME` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputRW.CZONE` | pxTextInput | 1 | tidak | ALWAYS |
| `InputRW.DistrictName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `InputRW.ID` | pxTextInput | 1 | ya | 1==2 |
| `InputRW.Note` | pxTextInput | 1 | tidak | ALWAYS |
| `InputRW.PROVINCENAME` | pxAutoComplete | 1 | tidak | InputRiskAddress.NationName != '' |
| `InputRW.STS_AKTIF` | pxRadioButtons | 1 | tidak | ALWAYS |
| `InputRW.ZipCode` | pxTextInput | 1 | tidak | ALWAYS |

#### InputSpreadingLife_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 10 (unik 5) - read-only 0 - RepeatGrid 2 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | tidak | ALWAYS |
| `.IndexCoverage` | (tanpa) | 2 | tidak | ALWAYS |
| `.PremiLifeNusantaraRe` | pxNumber | 2 | tidak | ALWAYS |
| `.PremiumRetro` | pxNumber | 2 | tidak | ALWAYS |
| `.RateLifeRetro` | pxTextInput | 2 | tidak | ALWAYS |

#### InputSpreadingLife_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 10 (unik 5) - read-only 0 - RepeatGrid 2 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | tidak | ALWAYS |
| `.IndexCoverage` | (tanpa) | 2 | tidak | ALWAYS |
| `.PremiLifeNusantaraRe` | pxNumber | 2 | tidak | ALWAYS |
| `.PremiumRetro` | pxNumber | 2 | tidak | ALWAYS |
| `.RateLifeRetro` | pxTextInput | 2 | tidak | ALWAYS |

#### InputSpreadingPA_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 7 (unik 4) - read-only 0 - RepeatGrid 2 - grid: .ASMCoverage

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | tidak | ALWAYS |
| `.Premium` | pxNumber | 2 | tidak | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | tidak | ALWAYS |
| `.Rate` | pxTextInput | 2 | tidak | ALWAYS |

#### InputSpreadingPA_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 7 (unik 4) - read-only 7 - RepeatGrid 2 - grid: .ASMCoverage

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 2 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxTextInput | 2 | ya | ALWAYS |

#### InputSpreadingTravel_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 3 (unik 3) - read-only 0 - RepeatGrid 1 - grid: .ASMCoverage

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 1 | tidak | ALWAYS |
| `.Premium` | pxNumber | 1 | tidak | ALWAYS |
| `.TSI` | pxTextInput | 1 | tidak | ALWAYS |

#### InputSpreadingTravel_FacIn_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 3 (unik 3) - read-only 0 - RepeatGrid 1 - grid: .ASMCoverage

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | pxTextInput | 1 | tidak | ALWAYS |
| `.Premium` | pxNumber | 1 | tidak | ALWAYS |
| `.TSI` | pxTextInput | 1 | tidak | ALWAYS |

#### InstallmentList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Installment` - field terikat 6 (unik 6) - read-only 0 - RepeatGrid 1 - grid: .InstallmentList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxTextInput | 1 | tidak | ALWAYS |
| `.DueDate` | (tanpa) | 1 | tidak | ALWAYS |
| `.InstallmentPercentage` | (tanpa) | 1 | tidak | ALWAYS |
| `.Premium` | pxTextInput | 1 | tidak | ALWAYS |
| `.PremiumAfterPPN` | pxTextInput | 1 | tidak | ALWAYS |
| `.PremiumAfterTax` | pxTextInput | 1 | tidak | ALWAYS |

#### Installments_ReadOnly

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-TreatyInInstallment` - field terikat 8 (unik 8) - read-only 6 - RepeatGrid 1 - grid: .InstallmentList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Amount` | pxNumber | 1 | ya | ALWAYS |
| `.AmountTotal` | pxNumber | 1 | tidak | ALWAYS |
| `.DueDate` | pxDateTime | 1 | ya | ALWAYS |
| `.Installment` | pxTextInput | 1 | ya | ALWAYS |
| `.InstallmentPct` | pxNumber | 1 | ya | ALWAYS |
| `.PaymentDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PctTotal` | pxNumber | 1 | tidak | ALWAYS |
| `.WPC` | pxTextInput | 1 | ya | ALWAYS |

#### InsuredNameTable

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-CLIENT` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 1 - grid: pgRepPgSubSectionInsuredNameTableBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ID` | (tanpa) | 1 | tidak | ALWAYS |
| `.Name` | (tanpa) | 1 | tidak | ALWAYS |

#### InwardFacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 25 (unik 24) - read-only 18 - RepeatGrid 1 - grid: pyWorkPage.OfferFacIn.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AdditionalCapital` | pxNumber | 1 | tidak [IsClaim \|\| IsRenewal \|\| IsUW] | ALWAYS |
| `.ASMRate` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMRIComm` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMShareTotal` | pxNumber | 1 | ya | ALWAYS |
| `.Currency.Currency` | pxAutoComplete | 1 | tidak [IsClaim \|\| IsRenewal \|\| IsUW] | ALWAYS |
| `.FacInGrossPremiumASM` | pxNumber | 1 | ya | ALWAYS |
| `.FacInNetPremiumASM` | pxNumber | 1 | ya | ALWAYS |
| `.InwardScale.Desc` | pxDisplayText | 1 | ya | 1 = 2 |
| `.InwardScale.ShareMin` | pxNumber | 1 | ya | 1 = 2 |
| `.LocationOfTopRisk` | pxAutoComplete/pxDisplayText | 2 | campur | .LocationOfTopRisk != '' ; .LocationOfTopRisk = '' |
| `.MaxPctTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.MaxTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.Name` | pxTextInput | 1 | ya | ALWAYS |
| `.PercentShare` | pxNumber | 1 | ya [IsUW] | ALWAYS |
| `.PolicyData.CedantRetention` | pxNumber | 1 | ya [IsUW] | ALWAYS |
| `.Remarks` | pxTextArea | 1 | tidak | ALWAYS |
| `.ShareInTSI` | pxNumber | 1 | ya | ALWAYS |
| `.ShareOffered` | pxTextInput | 1 | ya [IsUW] | ALWAYS |
| `.Status` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.SumTotalTSI` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | 1 = 2 |

#### LayerListDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 15 (unik 15) - read-only 14 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Day` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.Discount` | pxTextInput | 1 | ya [IsClaim \|\| IsRenewal \|\| IsGroupFac \|\| IsGr...] | ALWAYS |
| `.DiscountPercentage` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsGroupFac \|\| IsGr...] | ALWAYS |
| `.EmlPml` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsUW \|\| pyPortal.I...] | .CoverageBasis==3 |
| `.FirstLoss` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsUW \|\| pyPortal.I...] | .CoverageBasis==2 |
| `.FirstScale` | pxNumber | 1 | ya | .CoverageBasis==2 |
| `.Indemnity` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsUW \|\| pyPortal.I...] | .CoverageBasis==2 |
| `.IndemnityPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.LostLimit` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsUW \|\| pyPortal.I...] | .CoverageBasis==4 \|\| .CoverageBasis==5 |
| `.Premium` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsGroupFac \|\| IsGr...] | ALWAYS |
| `.Rate` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsGroupFac \|\| IsGr...] | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSILiability` | pxTextInput | 1 | ya | ALWAYS |
| `.Unit` | pxDropdown | 1 | ya [IsClaim \|\| IsUW \|\| pyPortal.IsShowOpenPolicy...] | .CoverageBasis==2 |
| `pyWorkPage.OfferFacIn.ProRatePercent` | pxNumber | 1 | ya | ALWAYS |

#### LayerListDtl_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 14 (unik 14) - read-only 14 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Discount` | pxTextInput | 1 | ya | ALWAYS |
| `.DiscountPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.EmlPml` | pxNumber | 1 | ya | .CoverageBasis==3 |
| `.FirstLoss` | pxNumber | 1 | ya | .CoverageBasis==2 |
| `.FirstScale` | pxNumber | 1 | ya | .CoverageBasis==2 |
| `.Indemnity` | pxNumber | 1 | ya | .CoverageBasis==2 |
| `.IndemnityPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.LostLimit` | pxNumber | 1 | ya | .CoverageBasis==4 \|\| .CoverageBasis==5 |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsGroupFac \|\| IsGr...] | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSILiability` | pxTextInput | 1 | ya | ALWAYS |
| `.Unit` | pxDropdown | 1 | ya | .CoverageBasis==2 |
| `pyWorkPage.OfferFacIn.ProRatePercent` | pxNumber | 1 | ya | ALWAYS |

#### LimitTreaty

rumpun **INTI** - kelas dominan `(kosong)` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 1 - grid: OfferJson.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.REINSTYPENAME` | pxTextInput | 1 | tidak | ALWAYS |
| `.RP` | pxNumber | 1 | tidak | ALWAYS |

#### ListPayment_SC

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: pyWorkPage.PaymentData.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AgingAmount` | pxNumber | 1 | ya | ALWAYS |
| `.CurrencyID` | (tanpa) | 1 | ya | ALWAYS |
| `.TotalPayment` | pxNumber | 1 | ya | ALWAYS |
| `pyWorkPage.PaymentData.AgingStatus` | pxTextInput | 1 | ya | ALWAYS |

#### ListSuggest

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-SuggestList` - field terikat 7 (unik 5) - read-only 4 - RepeatGrid 1 - grid: .SuggestList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Date` | pxDateTime | 1 | ya | NOTBLANK |
| `.IsApproved` | pxDropdown/pxRadioButtons | 2 | campur | pyWorkPage.PositionNote !='ReasTreatyInSecHead' |
| `.OperatorName` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ProductionDate` | pxDateTime | 1 | tidak | .IsApproved == 1 && (OperatorID.pyUserIdentifier=='NANDINA' \|\| Opera... |
| `.Suggest` | pxDisplayText/pxTextArea | 2 | campur | NOTBLANK |

#### LocationDetail

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 3 - grid: .Property.OccupationList, .Property.PropertyItemList, D_CoverageSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Currency` | (tanpa) | 1 | ya | ALWAYS |
| `.ItemType` | (tanpa) | 1 | ya | ALWAYS |
| `.OccupationId` | (tanpa) | 1 | ya | ALWAYS |
| `.OccupationName` | (tanpa) | 1 | ya | ALWAYS |
| `.OLDID` | (tanpa) | 1 | ya | ALWAYS |
| `.Rate` | (tanpa) | 1 | ya | ALWAYS |
| `.TableOfLimit.Description` | (tanpa) | 1 | ya | ALWAYS |
| `.TSIObjectItem` | (tanpa) | 1 | ya | ALWAYS |

#### LocationPropertyFacOut

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PrintMemoPlacingData` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 1 - grid: ListLocation.Object

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | tidak | ALWAYS |

#### MarineCargoDtl

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Cargo` - field terikat 32 (unik 32) - read-only 29 - RepeatGrid 1 - grid: .PolicyData.Ship.AdditionalShip

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.FromRute` | pxTextInput | 1 | ya [IsUW] | ALWAYS |
| `.ID` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.NM_SHIP` | pxTextInput | 1 | tidak | ALWAYS |
| `.PolicyData.BLNumber` | pxTextArea | 1 | ya [IsUW] | ALWAYS |
| `.PolicyData.InvoiceDate` | pxDateTime | 1 | ya [IsUW] | ALWAYS |
| `.PolicyData.InvoiceNumber` | pxTextInput | 1 | ya [IsUW] | ALWAYS |
| `.PolicyData.LC.LC` | pxCheckbox | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCBANK` | pxTextInput | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCCondition` | pxTextInput | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCEndDate` | pxDateTime | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCNumber` | pxTextInput | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.LC.LCRemark` | pxTextInput | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.MaxAge` | pxAutoComplete | 1 | ya [IsUW \|\| IsSpreadingUW] | ALWAYS |
| `.PolicyData.MaxAgeDetail` | pxTextArea | 1 | ya | ALWAYS |
| `.PolicyData.SailDate` | pxDateTime | 1 | ya [IsUW] | IsUW |
| `.PolicyData.Ship.AGE` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.CONST` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.DWT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.FLAG_NAME` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.GRT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NM_SHIP` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NRT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.REMARK` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.SHIPTYPE_NAME` | pxDropdown | 1 | ya | ALWAYS |
| `.PolicyData.Ship.Y_MAKE1` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.ASMAgentID` | pxAutoComplete | 1 | ya [IsUW] | 1=2 |
| `.PolicyData.SurveyAgent.pyAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyCity` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyCountry` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyFullName` | pxTextArea | 1 | ya | ALWAYS |
| `.ToRute` | pxTextInput | 1 | ya [IsUW] | ALWAYS |
| `.VOYAGE_NO` | pxTextInput | 1 | tidak | ALWAYS |

#### MarineCargoDtl_IsUW

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Cargo` - field terikat 33 (unik 33) - read-only 32 - RepeatGrid 1 - grid: .PolicyData.Ship.AdditionalShip

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.FromRute` | pxTextInput | 1 | ya | ALWAYS |
| `.NM_SHIP` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.PolicyData.SurveyAgent.pyCity` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.PolicyData.SurveyAgent.pyCountry` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.BLNumber` | pxTextArea | 1 | ya | ALWAYS |
| `.PolicyData.InvoiceDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyData.InvoiceNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.LC.LC` | pxCheckbox | 1 | tidak | ALWAYS |
| `.PolicyData.LC.LCBANK` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.LC.LCCondition` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.LC.LCEndDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyData.LC.LCNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.LC.LCRemark` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.MaxAge` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PolicyData.MaxAgeDetail` | pxTextArea | 1 | ya | ALWAYS |
| `.PolicyData.SailDate` | pxDateTime | 1 | ya | IsUW |
| `.PolicyData.Ship.AGE` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.CONST` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.DWT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.FLAG_NAME` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.GRT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.ID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NM_SHIP` | pxDropdown | 1 | ya | ALWAYS |
| `.PolicyData.Ship.NRT` | pxNumber | 1 | ya | ALWAYS |
| `.PolicyData.Ship.REMARK` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.Ship.SHIPTYPE_NAME` | pxDropdown | 1 | ya | ALWAYS |
| `.PolicyData.Ship.Y_MAKE1` | pxDisplayText | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.ASMAgentID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.PolicyData.SurveyAgent.pyFullName` | pxTextInput | 1 | ya | ALWAYS |
| `.REGISTER_ID` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ToRute` | pxTextInput | 1 | ya | ALWAYS |
| `.VOYAGE_NO` | pxTextInput | 1 | ya | ALWAYS |

#### ObjectDetails

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 23 (unik 23) - read-only 23 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.AlmRiskID` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.BuildingConstruction.FloorType` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.BuildingConstruction.NumberOfFloor` | pxNumber | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.BuildingConstruction.OthersType` | pxTextInput | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.BuildingConstruction.PartitionType` | pxTextInput | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.BuildingConstruction.RoofType` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.BuildingConstruction.SupportWallType` | pxTextInput | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.BuildingConstruction.WallType` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.BuildingNo` | pxNumber | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.Country` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.IsMaterialDamage` | pxCheckbox | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.IsTopRisk` | pxCheckbox | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.ObjectName` | pxTextInput | 1 | ya [.FlagDelete = '1'] | .Property.ObjectType = 'Others' |
| `.Property.ObjectNo` | pxTextInput | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.ObjectType` | pxDropdown | 1 | ya [.FlagDelete = '1'] | .Property.ObjectType!='Rumah Tinggal' |
| `.Property.Province` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMCity` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMDistrict` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMRW` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RoadName` | pxTextArea | 1 | ya | .Property.AlmRiskID != '' |
| `.Property.RoadType` | pxDropdown | 1 | ya | ALWAYS |

#### ObjectDetails_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 23 (unik 23) - read-only 23 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.AlmRiskID` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.BuildingConstruction.FloorType` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.BuildingConstruction.NumberOfFloor` | pxNumber | 1 | ya | ALWAYS |
| `.Property.BuildingConstruction.OthersType` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.BuildingConstruction.PartitionType` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.BuildingConstruction.RoofType` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.BuildingConstruction.SupportWallType` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.BuildingConstruction.WallType` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.BuildingNo` | pxNumber | 1 | ya | ALWAYS |
| `.Property.Country` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.IsMaterialDamage` | pxCheckbox | 1 | ya | ALWAYS |
| `.Property.IsTopRisk` | pxCheckbox | 1 | ya | ALWAYS |
| `.Property.ObjectName` | pxTextInput | 1 | ya | .Property.ObjectType = 'Others' |
| `.Property.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.ObjectType` | pxDropdown | 1 | ya | .Property.ObjectType!='Rumah Tinggal' |
| `.Property.Province` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMCity` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMDistrict` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMRW` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RoadName` | pxTextArea | 1 | ya | .Property.AlmRiskID != '' |
| `.Property.RoadType` | pxDropdown | 1 | ya | ALWAYS |

#### ObjectDtlAneka_FacIn

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 9 (unik 4) - read-only 5 - RepeatGrid 4 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 2 | tidak | ALWAYS |
| `.ObjectName` | pxDisplayText/pxTextInput | 4 | campur [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |

#### ObjectItemSummary

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: .Property.PropertyItemList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | (tanpa) | 1 | ya | ALWAYS |
| `.ItemType` | (tanpa) | 1 | ya | ALWAYS |
| `.TSIObjectItem` | (tanpa) | 1 | ya | ALWAYS |

#### ObjectList

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .LocationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.IsTopRisk` | pxCheckbox | 1 | ya | ALWAYS |
| `.Property.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText | 1 | ya | ALWAYS |

#### ObjectList_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .LocationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.IsTopRisk` | pxCheckbox | 1 | ya | ALWAYS |
| `.Property.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxDisplayText | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxDisplayText | 1 | ya | ALWAYS |

#### ObjectList_LossRecord

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 12 (unik 12) - read-only 8 - RepeatGrid 1 - grid: .OfferFacIn.ListCauseOfLoss

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CauseOfLoss` | pxTextInput | 1 | ya | ALWAYS |
| `.ClassOfBusiness` | pxTextInput | 1 | ya | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | ALWAYS |
| `.DOL` | pxDateTime | 1 | ya | ALWAYS |
| `.LOSS_AMOUNT` | pxNumber | 1 | ya | ALWAYS |
| `.LOSS_AMOUNT_IDR` | pxNumber | 1 | ya | ALWAYS |
| `.OfferFacIn.PremiMB` | pxNumber | 1 | tidak | ALWAYS |
| `.OfferFacIn.PremiOther` | pxNumber | 1 | tidak | ALWAYS |
| `.OfferFacIn.PremiPAR` | pxNumber | 1 | tidak | ALWAYS |
| `.OfferFacIn.PremiPL` | pxNumber | 1 | tidak | ALWAYS |
| `.UWYEAR` | pxTextInput | 1 | ya | ALWAYS |
| `TempError.ERRMSG` | pxTextArea | 1 | ya | NOTBLANK |

#### ObjectOccupation_FacIn

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 20 (unik 10) - read-only 14 - RepeatGrid 6 - grid: .Property.RiskLocation.AnekaList, .Property.RiskLocation.OccupationList, .Property.TotalTSIPremiGrossList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 2 | tidak | ALWAYS |
| `.Name` | pxDisplayText | 2 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 2 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OccupationName` | pxTextInput | 2 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 2 | ya | ALWAYS |
| `.RateOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 4 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 2 | ya | ALWAYS |

#### OccupationDetail

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: .Property.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OccupationId` | (tanpa) | 1 | ya | ALWAYS |
| `.OccupationName` | (tanpa) | 1 | ya | ALWAYS |
| `.TableOfLimit.Description` | (tanpa) | 1 | ya | ALWAYS |

#### OccupationFacOut

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 7 (unik 7) - read-only 4 - RepeatGrid 2 - grid: .Property.RiskLocation.AnekaList, .Property.RiskLocation.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | tidak | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OccupationName` | (tanpa) | 1 | tidak | ALWAYS |
| `.PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 1 | ya | ALWAYS |

#### OccupationFacOut_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 7 (unik 7) - read-only 7 - RepeatGrid 2 - grid: .Property.RiskLocation.AnekaList, .Property.RiskLocation.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | (tanpa) | 1 | ya | ALWAYS |
| `.PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 1 | ya | ALWAYS |

#### OccupationItemFacIn_Section

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OccupationId` | pxAutoComplete | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextArea | 1 | ya | ALWAYS |
| `.TableOfLimit.Description` | pxAutoComplete | 1 | ya [Always] | ALWAYS |

#### OccupationItemFacIn_Section_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OccupationId` | pxAutoComplete | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextArea | 1 | ya | ALWAYS |
| `.TableOfLimit.Description` | pxAutoComplete | 1 | ya | ALWAYS |

#### OccupationList

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: .Property.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OccupationId` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.TableOfLimit.Description` | pxAutoComplete | 1 | ya | ALWAYS |

#### OccupationList_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: .Property.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OccupationId` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.TableOfLimit.Description` | pxAutoComplete | 1 | ya | ALWAYS |

#### OccupationPropertyFacOut

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: .OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OccupationName` | pxTextArea | 1 | ya | ALWAYS |

#### OfferFacIn_NusaRe

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 18 (unik 18) - read-only 12 - RepeatGrid 1 - grid: pyWorkPage.OfferFacIn.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AdditionalCapital` | pxNumber | 1 | tidak [IsClaim \|\| IsRenewal \|\| IsUW] | ALWAYS |
| `.ASMRate` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMRIComm` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMShareTotal` | pxNumber | 1 | ya | ALWAYS |
| `.Currency.Currency` | pxAutoComplete | 1 | tidak [IsClaim \|\| IsRenewal \|\| IsUW] | ALWAYS |
| `.FacInGrossPremiumASM` | pxNumber | 1 | ya | ALWAYS |
| `.FacInNetPremiumASM` | pxNumber | 1 | ya | ALWAYS |
| `.MaxPctTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.MaxTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.Name` | pxTextInput | 1 | ya | ALWAYS |
| `.Remarks` | pxTextArea | 1 | tidak | ALWAYS |
| `.ShareInTSI` | pxNumber | 1 | ya | ALWAYS |
| `.Status` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.SumTotalTSI` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | 1 = 2 |

#### OfferFacIn_NusaRe_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 18 (unik 18) - read-only 14 - RepeatGrid 1 - grid: pyWorkPage.OfferFacIn.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AdditionalCapital` | pxNumber | 1 | ya [IsClaim \|\| IsRenewal \|\| IsUW] | ALWAYS |
| `.ASMRate` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMRIComm` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMShareTotal` | pxNumber | 1 | ya | ALWAYS |
| `.Currency.Currency` | pxAutoComplete | 1 | ya [IsClaim \|\| IsRenewal \|\| IsUW] | ALWAYS |
| `.FacInGrossPremiumASM` | pxNumber | 1 | ya | ALWAYS |
| `.FacInNetPremiumASM` | pxNumber | 1 | ya | ALWAYS |
| `.MaxPctTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.MaxTreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.Name` | pxTextInput | 1 | ya | ALWAYS |
| `.Remarks` | pxTextArea | 1 | tidak | ALWAYS |
| `.ShareInTSI` | pxNumber | 1 | ya | ALWAYS |
| `.Status` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.SumTotalTSI` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyCapacity` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | 1 = 2 |

#### OfferFacinLifeNotification

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### OfferStatusFacOut

rumpun **MBU** - kelas dominan `ASM-FW-GISFW-Data-Vehicle` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 2 - grid: .OfferFacIn.FacRetroList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Attention` | pxTextArea | 1 | ya | ALWAYS |
| `.BrandName` | pxDropdown | 1 | ya | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya | ALWAYS |
| `.ModelName` | pxTextInput | 1 | ya | ALWAYS |
| `.ReinsurerName` | (tanpa) | 1 | ya | ALWAYS |
| `.TypeName` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.FacRetroDetails.BackUpStatus` | pxDropdown | 1 | ya | ALWAYS |

#### OldDataEndorsementDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 12 (unik 11) - read-only 7 - RepeatGrid 1 - grid: SpreadingList.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OfferFacIn.CedingRetention` | pxNumber | 1 | tidak [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | IsLife |
| `.OfferFacIn.PercentShare` | pxNumber | 1 | tidak [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | ALWAYS |
| `.OfferFacIn.PolicyData.Payment.RICommision` | pxNumber | 1 | tidak [IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW ...] | ALWAYS |
| `.OfferFacIn.TotalPremiNusaRe` | pxTextInput | 2 | ya | ALWAYS |
| `.OfferFacIn.TotalRIComNusaRe` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.TotalTSINusaRe` | pxNumber | 1 | tidak | 1=2 |
| `.OfferFacIn.TotalTSINusaReSpreading` | pxTextInput | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 1 | ya [IsUW] | ALWAYS |
| `.TreatyType` | pxDropdown | 1 | ya [IsUW] | ALWAYS |
| `TempPolis.CARI12` | pxDropdown | 1 | tidak | ALWAYS |
| `TempPolis.CARI16` | pxDisplayText | 1 | ya | NOTBLANK |

#### OldOfferStatusFacOut

rumpun **INTI+MBU** - kelas dominan `ASM-FW-GISFW-Data-FacOffer` - field terikat 20 (unik 12) - read-only 19 - RepeatGrid 4 - grid: .OfferFacIn.OldData.OldData.FacRetroList, .VehicleList, TempOldData.OfferFacIn.FacRetroList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Attention` | pxTextArea | 3 | ya | ALWAYS |
| `.BrandName` | pxDropdown | 1 | ya | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya | ALWAYS |
| `.ModelName` | pxTextInput | 1 | ya | ALWAYS |
| `.NoOfferSlip` | (tanpa) | 2 | ya | ALWAYS |
| `.PrintRISlip.FACRETROSLIPNUMBER` | (tanpa) | 2 | ya | ALWAYS |
| `.ReinsurerName` | (tanpa) | 3 | ya | ALWAYS |
| `.TypeName` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.FacRetroDetails.BackUpStatus` | pxDropdown | 3 | ya | ALWAYS |
| `TempPolis.CARI12` | pxDropdown | 1 | tidak | TempPolis.CARI18=='1' \|\| TempPolis.CARI18=='2' |
| `TempPolis.CARI16` | pxDisplayText | 1 | ya | ALWAYS |

#### PASummarySection

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: D_PersonListSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.pyFullName` | (tanpa) | 1 | ya | ALWAYS |

#### PaymentCurrencyList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 9 (unik 6) - read-only 9 - RepeatGrid 3 - grid: pyWorkPage.OfferFacIn.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Name` | pxTextInput | 3 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Installment` | pxTextInput | 1 | ya [pyWorkPage.OfferFacIn.QuotationData.StatusBusine...] | !IsLife |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 2 | ya | ALWAYS |

#### PaymentCurrencyList_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 9 (unik 6) - read-only 9 - RepeatGrid 3 - grid: pyWorkPage.OfferFacIn.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Name` | pxTextInput | 3 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Installment` | pxTextInput | 1 | ya | !IsClaim |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 2 | ya | ALWAYS |

#### PaymentCurrencyListLife

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 11 (unik 5) - read-only 11 - RepeatGrid 3 - grid: pyWorkPage.OfferFacIn.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Name` | pxTextInput | 3 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 3 | ya | ALWAYS |

#### PaymentCurrencyListLife_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 11 (unik 5) - read-only 11 - RepeatGrid 3 - grid: pyWorkPage.OfferFacIn.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Name` | pxTextInput | 3 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 2 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 3 | ya | ALWAYS |

#### Periode

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 44 (unik 41) - read-only 23 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CommentDeductible` | pxTextArea | 1 | tidak | pyWorkPage.PositionNote = 'ReasFacInTeamLeader' &&.IsDeductibleAccepta... |
| `.EndorsmentReason` | pxTextInput | 1 | ya | NOTBLANK |
| `.Following` | pxTextInput | 1 | ya | !IsLife |
| `.FollowingNB` | pxDropdown | 2 | tidak | 1=2 |
| `.IsDeductibleAcceptance` | pxCheckbox | 1 | tidak | pyWorkPage.PositionNote = 'ReasFacInTeamLeader' && pyWorkPage.IsBedaCu... |
| `.IsProRate` | pxDropdown | 1 | tidak [IsUW] | Never |
| `.IsSpecialAcceptance` | pxCheckbox | 1 | tidak | 1=2 |
| `.PolicyData.EndDateTime` | pxDateTime | 1 | ya [.ConfirmBinding = 1\|\|pyWorkPage.FlagOnGoingPol...] | ALWAYS |
| `.PolicyData.EndorsementNo` | pxTextInput | 1 | ya | NOTBLANK |
| `.PolicyData.OfferingDate` | pxDateTime | 1 | ya [.ConfirmBinding = 1] | ALWAYS |
| `.PolicyData.PolicyNo` | pxTextInput | 1 | ya | NOTBLANK |
| `.PolicyData.StartDateTime` | pxDateTime | 1 | ya [.ConfirmBinding = 1\|\|pyWorkPage.FlagOnGoingPol...] | ALWAYS |
| `.PolicyMasterNumber` | pxTextInput | 1 | ya | IsMarineCargo |
| `.ProductBriguna` | pxDropdown | 1 | tidak [.FlagDelete==1] | IsCreditBriguna |
| `.ProRateType` | pxDropdown | 1 | tidak [IsUW] | IsLife |
| `.QuotationData.BusinessName` | pxDisplayText | 1 | ya | ALWAYS |
| `.QuotationData.CedingCoName` | pxDisplayText | 2 | ya | pyWorkPage.IsCedingConfirm!='Offer' && pyWorkPage.IsCedingConfirm!='Bi... |
| `.QuotationData.CommentOldPolicy` | pxTextArea | 1 | tidak | 1=2 |
| `.QuotationData.EdmDate` | pxDateTime | 1 | ya | IsEDM |
| `.QuotationData.EDMDay` | pxRadioButtons | 1 | tidak [IsUW] | IsEdmAdjCurrency\|\| IsEdmAdjTSI \|\| IsEdmExtendPeriod \|\| IsEdmAdjR... |
| `.QuotationData.GroupName` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.InputMonth` | pxTextInput | 1 | tidak | IsLife |
| `.QuotationData.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.IsMOP` | pxRadioButtons | 1 | ya | IsMarineCargo |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | 1 | tidak [IsUW] | 1=2 |
| `.QuotationData.IsTemplate` | pxDropdown | 1 | ya | .QuotationData.IsTemplate != '' && .QuotationData.StatusBusiness = 1 |
| `.QuotationData.MOID` | pxDropdown | 1 | tidak | IsLife |
| `.QuotationData.NoOfferSlip` | pxTextInput | 1 | tidak [IsUW] | ALWAYS |
| `.QuotationData.NPLStatus` | pxDropdown | 1 | tidak | .QuotationData.BusinessCode ='10053' |
| `.QuotationData.Period` | pxTextInput | 1 | tidak | IsLife |
| `.QuotationData.PeriodMM` | pxTextInput | 1 | tidak | IsLife |
| `.QuotationData.PolicyType` | pxRadioButtons | 1 | tidak | IsMarineCargo |
| `.QuotationData.QQName` | pxTextInput | 1 | tidak | pyWorkPage.IsCedingConfirm!='Offer' |
| `.QuotationData.SobName` | pxDisplayText | 2 | ya | pyWorkPage.IsCedingConfirm!='Offer' && pyWorkPage.IsCedingConfirm!='Bi... |
| `.QuotationData.StatusBusiness` | pxDropdown | 1 | ya | .QuotationData.IsTemplate = '' \|\| .QuotationData.StatusBusiness != 1 |
| `.QuotationData.TypeFacultative` | pxDropdown | 1 | tidak | ALWAYS |
| `.ScoringRisk.FinalScore` | pxInteger | 1 | ya | ALWAYS |
| `.ScoringRisk.NoteFinalScore` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.FlagErrorKonversi` | (tanpa) | 1 | ya [1=1] | NOTBLANK |
| `pyWorkPage.IsSpecialCase` | pxCheckbox | 1 | tidak | IsMarineCargo && .QuotationData.PolicyType ='2' && .QuotationData.Grou... |
| `TempError.ERRMSG` | pxTextArea | 1 | ya | NOTBLANK |

#### Periode_IsUW

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 37 (unik 37) - read-only 35 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BinderRNM` | pxTextInput | 1 | ya | NOTBLANK |
| `.EndorsmentReason` | pxTextInput | 1 | ya | NOTBLANK |
| `.Following` | pxTextInput | 1 | ya | IsFire \|\| IsMBD \|\| IsLiability |
| `.FollowingNB` | pxDropdown | 1 | ya | 1=2 |
| `.IsProRate` | pxDropdown | 1 | ya | ALWAYS |
| `.IsSpecialAcceptance` | pxCheckbox | 1 | tidak | 1=2 |
| `.PolicyData.EndDateTime` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyData.OfferingDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyData.StartDateTime` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyMasterNumber` | pxTextInput | 1 | ya | IsMarineCargo |
| `.ProductBriguna` | pxDropdown | 1 | ya | IsCreditBriguna |
| `.QuotationData.BusinessName` | pxDisplayText | 1 | ya | ALWAYS |
| `.QuotationData.CedingCoName` | pxDisplayText | 1 | ya | ALWAYS |
| `.QuotationData.CommentOldPolicy` | pxTextArea | 1 | ya | 1=2 |
| `.QuotationData.EdmDate` | pxDateTime | 1 | ya | IsEDM |
| `.QuotationData.EdmType` | pxDropdown | 1 | ya | NOTBLANK |
| `.QuotationData.GroupName` | pxDisplayText | 1 | ya | ALWAYS |
| `.QuotationData.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.IsGroup` | pxDropdown | 1 | ya | 1=2 |
| `.QuotationData.IsMOP` | pxRadioButtons | 1 | ya | IsMarineCargo |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | 1 | ya | IsMarineCargo |
| `.QuotationData.IsTemplate` | pxDropdown | 1 | ya | .QuotationData.IsTemplate != '' && .QuotationData.StatusBusiness = 1 |
| `.QuotationData.MarketingName` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.NoOfferSlip` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.NPLStatus` | pxDropdown | 1 | tidak | .QuotationData.BusinessCode ='10053' |
| `.QuotationData.OldPolicyNo` | pxTextInput | 1 | ya | IsRenewal |
| `.QuotationData.Period` | pxTextInput | 1 | ya | NOTBLANK |
| `.QuotationData.PolicyType` | pxRadioButtons | 1 | ya | NOTBLANK |
| `.QuotationData.QQName` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.SobName` | pxDisplayText | 1 | ya | ALWAYS |
| `.QuotationData.StatusBusiness` | pxDropdown | 1 | ya | .QuotationData.IsTemplate = '' \|\| .QuotationData.StatusBusiness != 1 |
| `.QuotationData.TypeFacultative` | pxTextInput | 1 | ya | ALWAYS |
| `.ScoringRisk.FinalScore` | pxInteger | 1 | ya | ALWAYS |
| `.ScoringRisk.NoteFinalScore` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.FlagErrorKonversi` | (tanpa) | 1 | ya [1=1] | NOTBLANK |
| `TempPolis.CARI24` | pxTextInput | 1 | ya | NOTBLANK |
| `TempPolis.CARI4` | pxTextInput | 1 | ya | NOTBLANK |

#### PeriodeEndorsement

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 26 (unik 25) - read-only 17 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.EndorsmentReason` | pxDisplayText | 1 | ya | ALWAYS |
| `.IsEditReffNumber` | pxCheckbox | 1 | tidak | !IsEDMRiSlip |
| `.IsSpecialAcceptance` | pxCheckbox | 1 | tidak | 1=2 |
| `.PolicyData.EndDateTime` | pxDateTime | 1 | tidak | ALWAYS |
| `.PolicyData.OfferingDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PolicyData.StartDateTime` | pxDateTime | 1 | ya | IsEdmAdjTSI \|\| IsEdmExtendPeriod \|\| IsEdmAdjRate \|\| IsEdmAdjSpre... |
| `.ProRateType` | pxDropdown | 1 | tidak [IsUW] | IsLife |
| `.QuotationData.BusinessName` | pxDisplayText | 1 | ya | ALWAYS |
| `.QuotationData.CedingCoName` | pxDisplayText | 2 | ya | !IsEdmAdjShareCedant&&!IsEdmAdjRefNo&&!IsEdmAdjCeding |
| `.QuotationData.EdmDate` | pxDateTime | 1 | ya | ALWAYS |
| `.QuotationData.EDMDay` | pxRadioButtons | 1 | tidak [IsUW] | IsEdmAdjCurrency\|\| IsEdmAdjTSI \|\| IsEdmExtendPeriod \|\| IsEdmAdjR... |
| `.QuotationData.EdmType` | pxDropdown | 1 | ya | 1=2 |
| `.QuotationData.EdmTypeNew` | pxDropdown | 1 | ya | .QuotationData.EdmType = 4 |
| `.QuotationData.GroupName` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | 1 | tidak [IsUW] | 1=2 |
| `.QuotationData.MarketingName` | pxTextInput | 1 | ya | ALWAYS |
| `.QuotationData.NoOfferSlip` | pxTextInput | 1 | tidak [!IsEdmAdjShareCedant&&!IsEdmAdjRefNo] | ALWAYS |
| `.QuotationData.OldPolicyNo` | pxDisplayText | 1 | ya | ALWAYS |
| `.QuotationData.Period` | pxTextInput | 1 | tidak | IsLife |
| `.QuotationData.PeriodMM` | pxTextInput | 1 | tidak | IsLife |
| `.QuotationData.QQName` | pxTextInput | 1 | ya [!IsEdmAdjInsured] | ALWAYS |
| `.QuotationData.SobName` | pxDisplayText | 1 | ya | ALWAYS |
| `.QuotationData.StatusBusiness` | pxDropdown | 1 | ya | ALWAYS |
| `.QuotationData.Type` | pxDropdown | 1 | ya | ALWAYS |

#### PeriodePolicy

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Policy` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.EndDateTime` | pxDateTime | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.StartDateTime` | pxDateTime | 1 | ya [IsSpreadingUW] | ALWAYS |

#### PersonCoverage

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 45 (unik 29) - read-only 12 - RepeatGrid 5 - grid: .ASMClause, .ASMCoverage, .ASMHeir, .ASMNotes, .ASMNotesDokter

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAge` | pxDisplayText | 2 | ya | ALWAYS |
| `.ASMBankData.ASMAccountNo` | pxTextInput | 2 | tidak | ALWAYS |
| `.ASMBankData.ASMBankName` | pxAutoComplete | 2 | tidak | ALWAYS |
| `.ASMBMI` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 3 | tidak [IsSpreadingUW] | ALWAYS |
| `.ASMDateTimeChronology` | pxDateTime | 2 | ya | NOTBLANK |
| `.ASMGender` | pxDropdown | 3 | tidak [!IsSpreadingUW] | ALWAYS |
| `.ASMHeight` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMHeirPercentage` | pxInteger | 1 | tidak [IsSpreadingUW] | ALWAYS |
| `.ASMInsuredRelationship` | pxDropdown | 1 | tidak | pyWorkPage.Quotation.BusinessCode=10057 |
| `.ASMLeftHanded` | pxDropdown | 1 | tidak | pyWorkPage.Quotation.BusinessCode == '10135' |
| `.ASMMaritalStatus` | pxDropdown | 2 | tidak | pyWorkPage.Quotation.BusinessCode!=10057 |
| `.ASMRelation` | pxDropdown | 2 | tidak [IsSpreadingUW] | pyWorkPage.Quotation.BusinessCode!=10057 |
| `.ASMUser` | pxDisplayText | 2 | ya | NOTBLANK |
| `.ASMWeight` | pxTextInput | 1 | tidak | ALWAYS |
| `.ClauseContent` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.ClauseDescription` | pxTextArea | 1 | tidak | ALWAYS |
| `.Coverage` | pxDropdown | 1 | tidak | ALWAYS |
| `.Loading` | (tanpa) | 1 | tidak | ALWAYS |
| `.LoadingAmount` | pxDisplayText | 1 | ya | ALWAYS |
| `.NoteDokter` | pxTextArea | 1 | tidak | IsNotEDM |
| `.PlanList(1).Name` | (tanpa) | 1 | tidak | ALWAYS |
| `.PremiRp` | (tanpa) | 1 | tidak | ALWAYS |
| `.Premium` | (tanpa) | 1 | tidak | ALWAYS |
| `.pyAddress` | (tanpa) | 1 | tidak [IsSpreadingUW] | ALWAYS |
| `.pyFullName` | pxTextInput | 3 | tidak [IsSpreadingUW] | ALWAYS |
| `.pyNote` | pxDisplayText/pxTextArea | 3 | campur | pyWorkPage.ProposalPosition != 'H3' && IsNotEDM |
| `.pyOrg` | pxDisplayText | 2 | ya | NOTBLANK |
| `.pyPhoneNumber` | (tanpa) | 1 | tidak [IsSpreadingUW] | ALWAYS |

#### PersonCoverageRO

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 2 (unik 2) - read-only  - RepeatGrid 1 - grid: .ASMCoverage

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Coverage` | pxDropdown | 1 | tidak | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |

#### PersonListDetail

rumpun **PA-LIFE** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 1 - grid: D_PersonListCoverageSummary.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMClass` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 1 | ya [1=1] | ALWAYS |
| `.ASMGender` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya [1=1] | ALWAYS |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | (tanpa) | 1 | ya | ALWAYS |
| `.Rate` | (tanpa) | 1 | ya | ALWAYS |
| `.TSI` | (tanpa) | 1 | ya | ALWAYS |

#### PersonPlanAndDeduct

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 7 (unik 7) - read-only 4 - RepeatGrid 1 - grid: .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Descriptions` | pxTextInput | 1 | tidak | ALWAYS |
| `.Loading` | pxTextInput | 1 | tidak | ALWAYS |
| `.LoadingAmount` | pxDisplayText | 1 | ya | ALWAYS |
| `.PlanList(1).Limit` | pxDisplayText | 1 | ya | ALWAYS |
| `.PlanList(1).Name` | pxDisplayText | 1 | ya | ALWAYS |
| `.Premium` | pxDisplayText | 1 | ya | ALWAYS |
| `.Value` | (tanpa) | 1 | tidak | ALWAYS |

#### PlanContains_UWRO

rumpun **PA-LIFE** - kelas dominan `ASM-FW-GISFW-Data-Benefit` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 3 - grid: .BenefitList, .ClauseList, .PremiumList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClauseContent` | (tanpa) | 1 | ya | ALWAYS |
| `.ClauseParam` | pxDisplayText | 1 | ya | ALWAYS |
| `.ID` | (tanpa) | 1 | ya | ALWAYS |
| `.Limit` | pxNumber | 1 | ya | ALWAYS |
| `.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.NumberOfMember` | pxDisplayText | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TypeName` | pxDropdown | 1 | ya | ALWAYS |

#### PlanListRO

rumpun **INTI** - kelas dominan `(kosong)` - field terikat 3 (unik 3) - read-only 0 - RepeatGrid 1 - grid: .PlanList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Limit` | pxNumber | 1 | tidak | ALWAYS |
| `.LimitCase` | pxNumber | 1 | tidak | ALWAYS |
| `.Name` | pxTextInput | 1 | tidak | ALWAYS |

#### PolicyTreatyInDeclineConfirm

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### PrintRISlip

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-FacOffer` - field terikat 16 (unik 11) - read-only 14 - RepeatGrid 4 - grid: .OfferFacIn.FacRetroList, TempViewSuggest.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Approval` | pxDropdown | 1 | ya | ALWAYS |
| `.Commentform` | pxTextArea | 1 | tidak | ALWAYS |
| `.CommentSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.Confirceding` | pxDropdown | 1 | tidak | ALWAYS |
| `.DateTransfer` | pxDateTime | 1 | ya | ALWAYS |
| `.EndPeriod` | pxDateTime | 2 | ya | ALWAYS |
| `.IsCedingConfirm` | pxDisplayText | 1 | ya | ALWAYS |
| `.OurRef` | pxTextInput | 2 | ya | ALWAYS |
| `.PICSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.ReinsurerName` | pxTextInput | 3 | ya | ALWAYS |
| `.StartPeriod` | pxDateTime | 2 | ya | ALWAYS |

#### PrintRISlip_Endorsement

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-FacOffer` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `RISlipContent.CARI1` | pxRichTextEditor | 1 | ya | ALWAYS |
| `RISlipContent.CARI3` | pxRichTextEditor | 1 | ya [RISlipContent.CARI2 = 0] | ALWAYS |
| `RISlipContent.CARI4` | pxRichTextEditor | 1 | ya | ALWAYS |

#### PrintRISlips_Endorsement

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-FacOffer` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `RISlipContent.CARI1` | pxRichTextEditor | 1 | ya | ALWAYS |
| `RISlipContent.CARI3` | pxRichTextEditor | 1 | ya [RISlipContent.CARI2 = 0] | ALWAYS |
| `RISlipContent.CARI4` | pxRichTextEditor | 1 | ya | ALWAYS |

#### Property

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### Property_IsUW

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### PropertyItemFacIn_Section

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 16 (unik 16) - read-only 15 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AreaHectar` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.Condition` | pxDropdown | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.IsAdjustableFlag` | pxCheckbox | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.ItemType` | pxDisplayText | 1 | ya | ALWAYS |
| `.ItemTypeID` | pxDropdown | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.NoOfTree` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |
| `.PctAdjust1` | pxDropdown | 1 | ya [IsClaim \|\| .FlagDelete = 1] | 1=2 |
| `.PctAdjust2` | pxDropdown | 1 | ya [IsClaim \|\| .FlagDelete = 1] | !IsAdjustableFlag |
| `.PctAdjustOther` | pxNumber | 1 | tidak | IsAdjustableFlag |
| `.PropertiItemNote` | pxTextArea | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.PropertyYear` | pxNumber | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.Remark` | pxTextArea | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.TSIObjectItem` | pxNumber | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.Unit` | pxNumber | 1 | ya [.FlagDelete = 1] | ALWAYS |
| `.Year` | pxTextInput | 1 | ya [.FlagDelete==1] | ALWAYS |

#### PropertyItemFacIn_Section_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 12 (unik 11) - read-only 11 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Condition` | pxDropdown | 1 | ya | ALWAYS |
| `.Currency` | pxAutoComplete | 1 | ya | ALWAYS |
| `.IsAdjustableFlag` | pxCheckbox | 1 | tidak | 1 = 2 |
| `.ItemType` | pxAutoComplete/pxDisplayText | 2 | ya | ALWAYS |
| `.PctAdjust2` | pxDropdown | 1 | ya | !IsAdjustableFlag |
| `.PctAdjustOther` | pxDropdown | 1 | ya | IsAdjustableFlag |
| `.PropertiItemNote` | pxTextArea | 1 | ya | ALWAYS |
| `.PropertyYear` | pxNumber | 1 | ya | ALWAYS |
| `.Remark` | pxTextArea | 1 | ya | ALWAYS |
| `.TSIObjectItem` | pxNumber | 1 | ya | ALWAYS |
| `.Unit` | pxNumber | 1 | ya | ALWAYS |

#### PropertyItemList

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 18 (unik 10) - read-only 16 - RepeatGrid 4 - grid: .Property.PropertyItemList, .Property.TotalTSIList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Condition` | pxDropdown | 2 | ya | ALWAYS |
| `.Currency` | pxDisplayText | 2 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 2 | ya | ALWAYS |
| `.Name` | (tanpa) | 2 | tidak | ALWAYS |
| `.PropertyYear` | pxNumber | 2 | ya | ALWAYS |
| `.TSI` | pxTextInput | 2 | ya | ALWAYS |
| `.TSIObjectItem` | pxNumber | 2 | ya | ALWAYS |
| `.TSIObjectItemOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSIOld` | pxTextInput | 1 | ya | ALWAYS |
| `.Unit` | pxNumber | 2 | ya | ALWAYS |

#### PropertyItemList_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 18 (unik 10) - read-only 16 - RepeatGrid 4 - grid: .Property.PropertyItemList, .Property.TotalTSIList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Condition` | pxDropdown | 2 | ya | ALWAYS |
| `.Currency` | pxDisplayText | 2 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 2 | ya | ALWAYS |
| `.Name` | (tanpa) | 2 | tidak | ALWAYS |
| `.PropertyYear` | pxNumber | 2 | ya | ALWAYS |
| `.TSI` | pxTextInput | 2 | ya | ALWAYS |
| `.TSIObjectItem` | pxNumber | 2 | ya | ALWAYS |
| `.TSIObjectItemOld` | pxNumber | 1 | ya | IsEDM |
| `.TSIOld` | pxTextInput | 1 | ya | ALWAYS |
| `.Unit` | pxNumber | 2 | ya | ALWAYS |

#### PropertyItemListCoverage

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 23 (unik 14) - read-only 23 - RepeatGrid 4 - grid: .Property.PropertyItemList, .Property.TotalTSIPremiGrossList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxDisplayText | 2 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 2 | ya | ALWAYS |
| `.Name` | pxDisplayText | 2 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 2 | ya | ALWAYS |
| `.RateOld` | pxNumber | 1 | ya | ALWAYS |
| `.TotalGrossPremi` | pxNumber | 2 | ya | ALWAYS |
| `.TotalGrossPremiOld` | pxNumber | 1 | ya | ALWAYS |
| `.TotalNetRate` | pxNumber | 2 | ya [.FlagNetRate = false] | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIObjectItem` | pxNumber | 2 | ya | ALWAYS |
| `.TSIObjectItemOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |

#### PropertyItemListCoverage_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 21 (unik 13) - read-only 21 - RepeatGrid 4 - grid: .Property.PropertyItemList, .Property.TotalTSIPremiGrossList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxDisplayText | 2 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 2 | ya | ALWAYS |
| `.Name` | pxDisplayText | 2 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumOld` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 2 | ya | ALWAYS |
| `.RateOld` | pxNumber | 1 | ya | ALWAYS |
| `.TotalGrossPremi` | pxNumber | 2 | ya | ALWAYS |
| `.TotalGrossPremiOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIObjectItem` | pxNumber | 2 | ya | ALWAYS |
| `.TSIObjectItemOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |

#### PropertyItemListCoverageCommision

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 15 (unik 10) - read-only 14 - RepeatGrid 4 - grid: .Property.PropertyItemList, .Property.RiskLocation.AnekaList, .Property.RiskLocation.OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxDisplayText | 2 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 2 | ya | ALWAYS |
| `.ObjectName` | (tanpa) | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TotalGrossPremi` | pxNumber | 2 | ya | ALWAYS |
| `.TotalGrossPremiOld` | pxNumber | 1 | ya | ALWAYS |
| `.TotalPremiumNusantaraRe` | pxNumber | 2 | ya | ALWAYS |
| `.TotalPremiumNusantaraReOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSIObjectItem` | pxNumber | 2 | ya | ALWAYS |
| `.TSIObjectItemOld` | pxNumber | 1 | ya | ALWAYS |

#### PropertyItemListCoverageSpreading

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 46 (unik 17) - read-only 40 - RepeatGrid 12 - grid: .Property.PropertyItemList, .Property.RiskLocation.AnekaList, .Property.RiskLocation.OccupationList, .Property.TotalTSIPremiSpreadRNM

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClaimEstimation` | pxNumber | 2 | ya | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.Currency` | pxDisplayText | 4 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 4 | ya | ALWAYS |
| `.ObjectName` | (tanpa) | 2 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 4 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | ya | OutputParam.ERRMSG6=='' |
| `.TotalGrossPremi` | pxNumber | 4 | ya | ALWAYS |
| `.TotalGrossPremiOld` | pxNumber | 2 | ya | ALWAYS |
| `.TotalPremiumNusantaraRe` | pxNumber | 4 | ya | ALWAYS |
| `.TotalPremiumNusantaraReOld` | pxNumber | 2 | ya | ALWAYS |
| `.TreatyName` | (tanpa) | 2 | tidak | ALWAYS |
| `.TreatyType` | pxDropdown | 2 | ya | ALWAYS |
| `.TSIObjectItem` | pxNumber | 4 | ya | ALWAYS |
| `.TSIObjectItemOld` | pxNumber | 2 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 2 | ya | ALWAYS |

#### PropertyItemListCoverageSpreading_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 26 (unik 19) - read-only 23 - RepeatGrid 6 - grid: .Property.PropertyItemList, .Property.RiskLocation.AnekaList, .Property.RiskLocation.OccupationList, .Property.TotalTSIPremiSpreadRNM, SpreadingListLoc.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClaimAmountIDR` | pxNumber | 1 | ya | ALWAYS |
| `.ClaimEstimation` | pxNumber | 1 | ya | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.Currency` | pxDisplayText | 2 | ya | ALWAYS |
| `.FlagSpreading` | pxDropdown | 1 | tidak [IsUW] | ALWAYS |
| `.ItemType` | pxTextInput | 2 | ya | ALWAYS |
| `.ObjectName` | (tanpa) | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | campur [IsUW] | ALWAYS |
| `.TotalGrossPremi` | pxNumber | 2 | ya | ALWAYS |
| `.TotalGrossPremiOld` | pxNumber | 1 | ya | ALWAYS |
| `.TotalPremiumNusantaraRe` | pxNumber | 2 | ya | ALWAYS |
| `.TotalPremiumNusantaraReOld` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyName` | pxTextInput | 1 | ya | ALWAYS |
| `.TreatyType` | pxDropdown | 2 | campur [IsUW] | ALWAYS |
| `.TSIObjectItem` | pxNumber | 2 | ya | ALWAYS |
| `.TSIObjectItemOld` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 1 | ya | ALWAYS |

#### PropertyItemListFacOut

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 9 (unik 9) - read-only 8 - RepeatGrid 2 - grid: .Property.PropertyItemList, .Property.TotalTSIPremiRetro

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxDisplayText | 1 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 1 | ya | ALWAYS |
| `.Name` | pxDisplayText | 1 | ya | ALWAYS |
| `.PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | tidak | ALWAYS |
| `.ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | (tanpa) | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 1 | ya | ALWAYS |

#### PropertyItemListFacOut_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 1 - grid: .Property.PropertyItemList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxDisplayText | 1 | ya | ALWAYS |
| `.ItemType` | pxTextInput | 1 | ya | ALWAYS |
| `.PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 1 | ya | ALWAYS |

#### ProtectCurrency

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### ProtectCurrency_IsUW

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### ProtectCurrencyCargo

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### ProtectCurrencyCargo_IsUW

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### pyAttachmentScreen

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### ReasViewAttachment

rumpun **INTI** - kelas dominan `Data-WorkAttach-File` - field terikat 7 (unik 7) - read-only 4 - RepeatGrid 2 - grid: AttachShowList.pxResults, pgRepPgSubSectionReasViewAttachmentBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.KATEGORI_2` | (tanpa) | 1 | ya | ALWAYS |
| `.NAMAFILE` | pxLink | 1 | tidak | ALWAYS |
| `.PNOTE` | (tanpa) | 1 | tidak | ALWAYS |
| `.pxCommitDateTime` | pxDateTime | 1 | ya | ALWAYS |
| `.pyAttachStream` | ASMThumbnail | 1 | tidak | 1 = 2 |
| `.pyNote` | pxTextInput | 1 | ya | ALWAYS |
| `.TANGGAL` | (tanpa) | 1 | ya | ALWAYS |

#### RejectNotificationSection

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### RemarksFacLetter

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-FacLetter` - field terikat 3 (unik 3) - read-only 0 - RepeatGrid 1 - grid: .FacLetterList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DateCreated` | pxDateTime | 1 | tidak | ALWAYS |
| `.Message` | pxTextArea | 1 | tidak | ALWAYS |
| `.Remarks` | pxTextInput | 1 | tidak | ALWAYS |

#### ReviseNotification

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### RiskAccumulationReport

rumpun **INTI** - kelas dominan `Data-Portal` - field terikat 1 (unik 1) - read-only  - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `InputParam.CARIDESC` | pxRichTextEditor | 1 | ya | ALWAYS |

#### RiskAround

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 24 (unik 24) - read-only 21 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.IsFlammableItemFlag` | pxCheckbox | 1 | tidak [IsUW] | ALWAYS |
| `.Property.IsHotWorkProcessFlag` | pxCheckbox | 1 | tidak [IsUW] | ALWAYS |
| `.Property.IsProductionProcessFlag` | pxCheckbox | 1 | tidak [IsUW] | ALWAYS |
| `.Property.Ownership` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.BackConstruction` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.BackDistance` | pxNumber | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.BackNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.BackOccupation` | pxAutoComplete | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.FloodArea` | pxDropdown | 1 | ya [.FlagDelete = '1'] | .Property.SurroundingRisk.FloodAreaStatus==0 |
| `.Property.SurroundingRisk.FloodAreaStatus` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.FrontConstruction` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.FrontDistance` | pxNumber | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.FrontNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.FrontOccupation` | pxAutoComplete | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.HousekeepingRemark` | pxTextArea | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.HousekeepingStatus` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.LeftConstruction` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.LeftDistance` | pxNumber | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.LeftNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.LeftOccupation` | pxAutoComplete | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.RightConstruction` | pxDropdown | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.RightDistance` | pxNumber | 1 | ya [.FlagDelete = '1'] | ALWAYS |
| `.Property.SurroundingRisk.RightNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.RightOccupation` | pxAutoComplete | 1 | ya [.FlagDelete = '1'] | ALWAYS |

#### RiskAround_IsUW

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` - field terikat 24 (unik 24) - read-only 24 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Property.IsFlammableItemFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.Property.IsHotWorkProcessFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.Property.IsProductionProcessFlag` | pxCheckbox | 1 | ya | ALWAYS |
| `.Property.Ownership` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.BackConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.BackDistance` | pxNumber | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.BackNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.BackOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.FloodArea` | pxDropdown | 1 | ya | .Property.SurroundingRisk.FloodAreaStatus==0 |
| `.Property.SurroundingRisk.FloodAreaStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.FrontConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.FrontDistance` | pxNumber | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.FrontNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.FrontOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.HousekeepingRemark` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.HousekeepingStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.LeftConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.LeftDistance` | pxNumber | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.LeftNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.LeftOccupation` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.RightConstruction` | pxDropdown | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.RightDistance` | pxNumber | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.RightNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Property.SurroundingRisk.RightOccupation` | pxAutoComplete | 1 | ya | ALWAYS |

#### ScoringRisk

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-ScoringRisk` - field terikat 438 (unik 436) - read-only 233 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.FireAlarmSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.FireAlarmSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox3= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox4= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox5= true |
| `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.HydrantSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.HydrantSystem.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.HydrantSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.HydrantSystem.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.HydrantSystem.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.HydrantSystem.ChechBox3= true |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox3= true |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox4= true |
| `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SecurityGuard.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SecurityGuard.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.SecurityGuard.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SecurityGuard.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.SecurityGuard.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SecurityGuard.ChechBox3= true |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox3= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox10` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox7` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox8` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox9` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.PercentLossRatio` | pxNumber | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox1= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score10` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox10= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox2= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox3= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox4= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox5= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox6= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score7` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox7= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score8` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox8= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score9` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox9= true |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox1 = true |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox2 = true |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox3 = true |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox4 = true |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox5 = true |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox1 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox2 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox3 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox4 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox5 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox6 =... |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox1 = true |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox2 = true |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox3 = true |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox4 = true |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox5 = true |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox1= true |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox2= true |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox3= true |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox4= true |
| `.DataScoringRiskList(1).Occupation.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox1 = true |
| `.DataScoringRiskList(1).Occupation.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox2 = true |
| `.DataScoringRiskList(1).Occupation.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox3 = true |
| `.DataScoringRiskList(1).Occupation.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox4 = true |
| `.DataScoringRiskList(1).Occupation.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox5 = true |
| `.DataScoringRiskList(1).Others.Clauses.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Clauses.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Clauses.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Clauses.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Clauses.ChechBox1= true |
| `.DataScoringRiskList(1).Others.Clauses.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Clauses.ChechBox2= true |
| `.DataScoringRiskList(1).Others.Clauses.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Clauses.ChechBox3= true |
| `.DataScoringRiskList(1).Others.Deductible.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Deductible.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Deductible.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Deductible.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Deductible.ChechBox1= true |
| `.DataScoringRiskList(1).Others.Deductible.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Deductible.ChechBox2= true |
| `.DataScoringRiskList(1).Others.Deductible.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Deductible.ChechBox3= true |
| `.DataScoringRiskList(1).Others.OtherInformation.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.OtherInformation.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.OtherInformation.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.OtherInformation.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.OtherInformation.ChechBox1= true |
| `.DataScoringRiskList(1).Others.OtherInformation.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.OtherInformation.ChechBox2= true |
| `.DataScoringRiskList(1).Others.OtherInformation.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.OtherInformation.ChechBox3= true |
| `.DataScoringRiskList(1).Others.Rate.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Rate.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Rate.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Rate.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Rate.ChechBox1= true |
| `.DataScoringRiskList(1).Others.Rate.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Rate.ChechBox2= true |
| `.DataScoringRiskList(1).Others.Rate.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Rate.ChechBox3= true |
| `.DataScoringRiskList(1).Others.TotalSumInsured.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.TotalSumInsured.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.TotalSumInsured.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.TotalSumInsured.ChechBox1= true |
| `.DataScoringRiskList(1).Others.TotalSumInsured.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.TotalSumInsured.ChechBox2= true |
| `.DataScoringRiskList(1).RiskImprovement.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskImprovement.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskImprovement.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskImprovement.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskImprovement.ChechBox1= true |
| `.DataScoringRiskList(1).RiskImprovement.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskImprovement.ChechBox2= true |
| `.DataScoringRiskList(1).RiskImprovement.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskImprovement.ChechBox3= true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox7` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox1 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox2= true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox3 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox4 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox5 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox6 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score7` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox7= true |
| `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Flood.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Flood.ChechBox1 = true |
| `.DataScoringRiskList(1).RiskLocation.Flood.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Flood.ChechBox2 = true |
| `.DataScoringRiskList(1).RiskLocation.Flood.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Flood.ChechBox3 = true |
| `.DataScoringRiskList(1).RiskLocation.Flood.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Flood.ChechBox4 = true |
| `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Riot.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Riot.ChechBox1 = true |
| `.DataScoringRiskList(1).RiskLocation.Riot.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Riot.ChechBox2 = true |
| `.DataScoringRiskList(1).RiskLocation.Riot.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Riot.ChechBox3 = true |
| `.DataScoringRiskList(1).RiskLocation.Riot.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Riot.ChechBox4 = true |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox1 = true |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox2 = true |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox3 = true |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox4 = true |
| `.DataScoringRiskList(1).Score` | pxNumber | 1 | ya | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox1= true |
| `.DataScoringRiskList(1).SurveyReport.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox2= true |
| `.DataScoringRiskList(1).SurveyReport.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox3= true |
| `.DataScoringRiskList(1).SurveyReport.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox4= true |
| `.DataScoringRiskList(1).SurveyReport.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox5= true |
| `.DataScoringRiskList(1).SurveyReport.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox6= true |
| `.DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.FireAlarmSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.FireAlarmSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox3= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox4= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox5= true |
| `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.HydrantSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.HydrantSystem.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.HydrantSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.HydrantSystem.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.HydrantSystem.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.HydrantSystem.ChechBox3= true |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox3= true |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox4= true |
| `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SecurityGuard.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SecurityGuard.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.SecurityGuard.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SecurityGuard.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.SecurityGuard.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SecurityGuard.ChechBox3= true |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox3= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox10` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox7` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox8` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox9` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.PercentLossRatio` | pxNumber | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox1= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score10` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox10= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox2= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox3= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox4= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox5= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox6= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score7` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox7= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score8` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox8= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score9` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox9= true |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox1 = true |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox2 = true |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox3 = true |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox4 = true |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox5 = true |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox1 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox2 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox3 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox4 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox5 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox6 =... |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox1 = true |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox2 = true |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox3 = true |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox4 = true |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox5 = true |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox1= true |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox2= true |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox3= true |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox4= true |
| `.DataScoringRiskList(2).Occupation.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox1 = true |
| `.DataScoringRiskList(2).Occupation.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox2 = true |
| `.DataScoringRiskList(2).Occupation.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox3 = true |
| `.DataScoringRiskList(2).Occupation.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox4 = true |
| `.DataScoringRiskList(2).Occupation.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox5 = true |
| `.DataScoringRiskList(2).Others.Clauses.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Clauses.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Clauses.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Clauses.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Clauses.ChechBox1= true |
| `.DataScoringRiskList(2).Others.Clauses.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Clauses.ChechBox2= true |
| `.DataScoringRiskList(2).Others.Clauses.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Clauses.ChechBox3= true |
| `.DataScoringRiskList(2).Others.Deductible.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Deductible.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Deductible.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Deductible.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Deductible.ChechBox1= true |
| `.DataScoringRiskList(2).Others.Deductible.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Deductible.ChechBox2= true |
| `.DataScoringRiskList(2).Others.Deductible.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Deductible.ChechBox3= true |
| `.DataScoringRiskList(2).Others.OtherInformation.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.OtherInformation.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.OtherInformation.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.OtherInformation.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.OtherInformation.ChechBox1= true |
| `.DataScoringRiskList(2).Others.OtherInformation.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.OtherInformation.ChechBox2= true |
| `.DataScoringRiskList(2).Others.OtherInformation.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.OtherInformation.ChechBox3= true |
| `.DataScoringRiskList(2).Others.Rate.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Rate.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Rate.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Rate.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Rate.ChechBox1= true |
| `.DataScoringRiskList(2).Others.Rate.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Rate.ChechBox2= true |
| `.DataScoringRiskList(2).Others.Rate.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Rate.ChechBox3= true |
| `.DataScoringRiskList(2).Others.TotalSumInsured.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.TotalSumInsured.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.TotalSumInsured.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.TotalSumInsured.ChechBox1= true |
| `.DataScoringRiskList(2).Others.TotalSumInsured.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.TotalSumInsured.ChechBox2= true |
| `.DataScoringRiskList(2).RiskImprovement.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskImprovement.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskImprovement.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskImprovement.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskImprovement.ChechBox1= true |
| `.DataScoringRiskList(2).RiskImprovement.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskImprovement.ChechBox2= true |
| `.DataScoringRiskList(2).RiskImprovement.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskImprovement.ChechBox3= true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox7` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox1 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox2= true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox3 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox4 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox5 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox6 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score7` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox7= true |
| `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Flood.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Flood.ChechBox1 = true |
| `.DataScoringRiskList(2).RiskLocation.Flood.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Flood.ChechBox2 = true |
| `.DataScoringRiskList(2).RiskLocation.Flood.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Flood.ChechBox3 = true |
| `.DataScoringRiskList(2).RiskLocation.Flood.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Flood.ChechBox4 = true |
| `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Riot.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Riot.ChechBox1 = true |
| `.DataScoringRiskList(2).RiskLocation.Riot.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Riot.ChechBox2 = true |
| `.DataScoringRiskList(2).RiskLocation.Riot.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Riot.ChechBox3 = true |
| `.DataScoringRiskList(2).RiskLocation.Riot.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Riot.ChechBox4 = true |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox1 = true |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox2 = true |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox3 = true |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox4 = true |
| `.DataScoringRiskList(2).Score` | pxNumber | 1 | ya | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox1= true |
| `.DataScoringRiskList(2).SurveyReport.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox2= true |
| `.DataScoringRiskList(2).SurveyReport.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox3= true |
| `.DataScoringRiskList(2).SurveyReport.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox4= true |
| `.DataScoringRiskList(2).SurveyReport.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox5= true |
| `.DataScoringRiskList(2).SurveyReport.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox6= true |
| `.ScoringRisk.DataScoringRiskList(1).OccupationNote` | pxTextInput | 1 | ya | ALWAYS |
| `.ScoringRisk.DataScoringRiskList(1).TopRiskLocation` | pxTextInput | 1 | ya | ALWAYS |
| `.ScoringRisk.DataScoringRiskList(2).OccupationNote` | pxAutoComplete | 1 | tidak | ALWAYS |
| `.ScoringRisk.DataScoringRiskList(2).TopRiskLocation` | pxTextInput | 1 | ya | ALWAYS |
| `.ScoringRisk.FinalScore` | pxNumber | 2 | ya | ALWAYS |
| `.ScoringRisk.NoteFinalScore` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.AboveAverage` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.AcceptableNotApproval` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.AcceptableWithApproval` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.Average` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.BelowAverage` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.Good` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.LossRatio` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.NotAcceptable` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.ObjectConditions` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.Occupation` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.OperationalDirector` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.OperationalDivHead` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.Others` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.Poor` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.PresidentDirector` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.TechnicalDirector` | pxCheckbox | 1 | ya | ALWAYS |
| `.ScoringRisk.ScoringResult.TotalSumInsuredR1` | pxTextInput | 1 | ya | 1=2 |
| `.ScoringRisk.ScoringResult.TotalSumInsuredR2` | pxNumber | 1 | ya | 1=2 |
| `.ScoringRisk.Status` | pxTextInput | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.PolicyData.EndDateTime` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.PolicyData.StartDateTime` | pxDateTime | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.QuotationData.InsuredName` | pxTextInput | 1 | ya | ALWAYS |

#### ScoringRiskForm1

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-ScoringRisk` - field terikat 96 (unik 96) - read-only 48 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DataScoringRiskList(1).Occupation.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Occupation.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox1 = true |
| `.DataScoringRiskList(1).Occupation.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox2 = true |
| `.DataScoringRiskList(1).Occupation.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox3 = true |
| `.DataScoringRiskList(1).Occupation.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox4 = true |
| `.DataScoringRiskList(1).Occupation.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).Occupation.ChechBox5 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox7` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox1 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox2= true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox3 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox4 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox5 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox6 = true |
| `.DataScoringRiskList(1).RiskLocation.Earthquake.Score7` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox7= true |
| `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Flood.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Flood.ChechBox1 = true |
| `.DataScoringRiskList(1).RiskLocation.Flood.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Flood.ChechBox2 = true |
| `.DataScoringRiskList(1).RiskLocation.Flood.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Flood.ChechBox3 = true |
| `.DataScoringRiskList(1).RiskLocation.Flood.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Flood.ChechBox4 = true |
| `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Riot.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Riot.ChechBox1 = true |
| `.DataScoringRiskList(1).RiskLocation.Riot.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Riot.ChechBox2 = true |
| `.DataScoringRiskList(1).RiskLocation.Riot.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Riot.ChechBox3 = true |
| `.DataScoringRiskList(1).RiskLocation.Riot.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Riot.ChechBox4 = true |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox1 = true |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox2 = true |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox3 = true |
| `.DataScoringRiskList(1).RiskLocation.Tsunami.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox4 = true |
| `.DataScoringRiskList(2).Occupation.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Occupation.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox1 = true |
| `.DataScoringRiskList(2).Occupation.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox2 = true |
| `.DataScoringRiskList(2).Occupation.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox3 = true |
| `.DataScoringRiskList(2).Occupation.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox4 = true |
| `.DataScoringRiskList(2).Occupation.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).Occupation.ChechBox5 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox7` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox1 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox2= true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox3 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox4 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox5 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox6 = true |
| `.DataScoringRiskList(2).RiskLocation.Earthquake.Score7` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox7= true |
| `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Flood.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Flood.ChechBox1 = true |
| `.DataScoringRiskList(2).RiskLocation.Flood.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Flood.ChechBox2 = true |
| `.DataScoringRiskList(2).RiskLocation.Flood.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Flood.ChechBox3 = true |
| `.DataScoringRiskList(2).RiskLocation.Flood.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Flood.ChechBox4 = true |
| `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Riot.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Riot.ChechBox1 = true |
| `.DataScoringRiskList(2).RiskLocation.Riot.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Riot.ChechBox2 = true |
| `.DataScoringRiskList(2).RiskLocation.Riot.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Riot.ChechBox3 = true |
| `.DataScoringRiskList(2).RiskLocation.Riot.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Riot.ChechBox4 = true |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox1 = true |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox2 = true |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox3 = true |
| `.DataScoringRiskList(2).RiskLocation.Tsunami.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox4 = true |

#### ScoringRiskForm2

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-ScoringRisk` - field terikat 158 (unik 158) - read-only 78 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.FireAlarmSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.FireAlarmSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox3= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox4= true |
| `.DataScoringRiskList(1).FEA.Firebrigade.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.Firebrigade.ChechBox5= true |
| `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.HydrantSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.HydrantSystem.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.HydrantSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.HydrantSystem.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.HydrantSystem.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.HydrantSystem.ChechBox3= true |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox3= true |
| `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox4= true |
| `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SecurityGuard.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SecurityGuard.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.SecurityGuard.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SecurityGuard.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.SecurityGuard.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SecurityGuard.ChechBox3= true |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox1= true |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox2= true |
| `.DataScoringRiskList(1).FEA.SprinklerSystem.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox3= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox10` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox7` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox8` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox9` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.PercentLossRatio` | pxNumber | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox1= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score10` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox10= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox2= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox3= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox4= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox5= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox6= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score7` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox7= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score8` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox8= true |
| `.DataScoringRiskList(1).LossRatio.LossRatio.Score9` | pxNumber | 1 | ya | .DataScoringRiskList(1).LossRatio.LossRatio.ChechBox9= true |
| `.DataScoringRiskList(1).RiskImprovement.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskImprovement.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskImprovement.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).RiskImprovement.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskImprovement.ChechBox1= true |
| `.DataScoringRiskList(1).RiskImprovement.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskImprovement.ChechBox2= true |
| `.DataScoringRiskList(1).RiskImprovement.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).RiskImprovement.ChechBox3= true |
| `.DataScoringRiskList(1).SurveyReport.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).SurveyReport.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox1= true |
| `.DataScoringRiskList(1).SurveyReport.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox2= true |
| `.DataScoringRiskList(1).SurveyReport.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox3= true |
| `.DataScoringRiskList(1).SurveyReport.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox4= true |
| `.DataScoringRiskList(1).SurveyReport.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox5= true |
| `.DataScoringRiskList(1).SurveyReport.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(1).SurveyReport.ChechBox6= true |
| `.DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.FireAlarmSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.FireAlarmSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox3= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox4= true |
| `.DataScoringRiskList(2).FEA.Firebrigade.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.Firebrigade.ChechBox5= true |
| `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.HydrantSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.HydrantSystem.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.HydrantSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.HydrantSystem.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.HydrantSystem.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.HydrantSystem.ChechBox3= true |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox3= true |
| `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox4= true |
| `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SecurityGuard.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SecurityGuard.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.SecurityGuard.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SecurityGuard.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.SecurityGuard.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SecurityGuard.ChechBox3= true |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox1= true |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox2= true |
| `.DataScoringRiskList(2).FEA.SprinklerSystem.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox3= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox10` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox7` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox8` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox9` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.PercentLossRatio` | pxNumber | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox1= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score10` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox10= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox2= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox3= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox4= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox5= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox6= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score7` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox7= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score8` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox8= true |
| `.DataScoringRiskList(2).LossRatio.LossRatio.Score9` | pxNumber | 1 | ya | .DataScoringRiskList(2).LossRatio.LossRatio.ChechBox9= true |
| `.DataScoringRiskList(2).RiskImprovement.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskImprovement.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskImprovement.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).RiskImprovement.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskImprovement.ChechBox1= true |
| `.DataScoringRiskList(2).RiskImprovement.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskImprovement.ChechBox2= true |
| `.DataScoringRiskList(2).RiskImprovement.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).RiskImprovement.ChechBox3= true |
| `.DataScoringRiskList(2).SurveyReport.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).SurveyReport.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox1= true |
| `.DataScoringRiskList(2).SurveyReport.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox2= true |
| `.DataScoringRiskList(2).SurveyReport.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox3= true |
| `.DataScoringRiskList(2).SurveyReport.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox4= true |
| `.DataScoringRiskList(2).SurveyReport.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox5= true |
| `.DataScoringRiskList(2).SurveyReport.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(2).SurveyReport.ChechBox6= true |

#### ScoringRiskForm3

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-ScoringRisk` - field terikat 58 (unik 58) - read-only 30 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DataScoringRiskList(1).Others.Clauses.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Clauses.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Clauses.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Clauses.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Clauses.ChechBox1= true |
| `.DataScoringRiskList(1).Others.Clauses.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Clauses.ChechBox2= true |
| `.DataScoringRiskList(1).Others.Clauses.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Clauses.ChechBox3= true |
| `.DataScoringRiskList(1).Others.Deductible.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Deductible.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Deductible.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Deductible.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Deductible.ChechBox1= true |
| `.DataScoringRiskList(1).Others.Deductible.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Deductible.ChechBox2= true |
| `.DataScoringRiskList(1).Others.Deductible.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Deductible.ChechBox3= true |
| `.DataScoringRiskList(1).Others.OtherInformation.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.OtherInformation.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.OtherInformation.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.OtherInformation.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.OtherInformation.ChechBox1= true |
| `.DataScoringRiskList(1).Others.OtherInformation.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.OtherInformation.ChechBox2= true |
| `.DataScoringRiskList(1).Others.OtherInformation.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.OtherInformation.ChechBox3= true |
| `.DataScoringRiskList(1).Others.Rate.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Rate.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Rate.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.Rate.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Rate.ChechBox1= true |
| `.DataScoringRiskList(1).Others.Rate.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Rate.ChechBox2= true |
| `.DataScoringRiskList(1).Others.Rate.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.Rate.ChechBox3= true |
| `.DataScoringRiskList(1).Others.TotalSumInsured.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.TotalSumInsured.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).Others.TotalSumInsured.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.TotalSumInsured.ChechBox1= true |
| `.DataScoringRiskList(1).Others.TotalSumInsured.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).Others.TotalSumInsured.ChechBox2= true |
| `.DataScoringRiskList(1).Score` | pxNumber | 1 | ya | ALWAYS |
| `.DataScoringRiskList(2).Others.Clauses.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Clauses.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Clauses.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Clauses.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Clauses.ChechBox1= true |
| `.DataScoringRiskList(2).Others.Clauses.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Clauses.ChechBox2= true |
| `.DataScoringRiskList(2).Others.Clauses.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Clauses.ChechBox3= true |
| `.DataScoringRiskList(2).Others.Deductible.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Deductible.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Deductible.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Deductible.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Deductible.ChechBox1= true |
| `.DataScoringRiskList(2).Others.Deductible.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Deductible.ChechBox2= true |
| `.DataScoringRiskList(2).Others.Deductible.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Deductible.ChechBox3= true |
| `.DataScoringRiskList(2).Others.OtherInformation.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.OtherInformation.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.OtherInformation.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.OtherInformation.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.OtherInformation.ChechBox1= true |
| `.DataScoringRiskList(2).Others.OtherInformation.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.OtherInformation.ChechBox2= true |
| `.DataScoringRiskList(2).Others.OtherInformation.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.OtherInformation.ChechBox3= true |
| `.DataScoringRiskList(2).Others.Rate.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Rate.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Rate.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.Rate.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Rate.ChechBox1= true |
| `.DataScoringRiskList(2).Others.Rate.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Rate.ChechBox2= true |
| `.DataScoringRiskList(2).Others.Rate.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.Rate.ChechBox3= true |
| `.DataScoringRiskList(2).Others.TotalSumInsured.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.TotalSumInsured.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).Others.TotalSumInsured.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.TotalSumInsured.ChechBox1= true |
| `.DataScoringRiskList(2).Others.TotalSumInsured.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).Others.TotalSumInsured.ChechBox2= true |
| `.DataScoringRiskList(2).Score` | pxNumber | 1 | ya | ALWAYS |

#### ScoringRiskForm4

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-ScoringRisk` - field terikat 96 (unik 96) - read-only 48 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox1 = true |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox2 = true |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox3 = true |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox4 = true |
| `.DataScoringRiskList(1).ObjectConditions.Building.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Building.ChechBox5 = true |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox1 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox2 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox3 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox4 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox5 =... |
| `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox6 =... |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox1 = true |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox2 = true |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox3 = true |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox4 = true |
| `.DataScoringRiskList(1).ObjectConditions.Machinery.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox5 = true |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox1= true |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox2= true |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox3= true |
| `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox4= true |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox1 = true |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox2 = true |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox3 = true |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox4 = true |
| `.DataScoringRiskList(2).ObjectConditions.Building.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Building.ChechBox5 = true |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox6` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox1 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox2 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox3 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox4 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox5 =... |
| `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score6` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox6 =... |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox5` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox1 = true |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox2 = true |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox3 = true |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox4 = true |
| `.DataScoringRiskList(2).ObjectConditions.Machinery.Score5` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox5 = true |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBo... |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox1` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox2` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox3` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox4` | pxCheckbox | 1 | tidak | ALWAYS |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score1` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox1= true |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score2` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox2= true |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score3` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox3= true |
| `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score4` | pxNumber | 1 | ya | .DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox4= true |

#### SearchRiskAccumCov

rumpun **INTI** - kelas dominan `Data-Portal` - field terikat 18 (unik 18) - read-only 6 - RepeatGrid 2 - grid: pgRepPgSubSectionSearchRiskAccumCovBBBBBBBBBBBBBB.pxResults, ResultsAccumulation.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AccumulationName` | pxDisplayText | 1 | ya | ALWAYS |
| `.CARI1` | pxDisplayText | 1 | ya | ALWAYS |
| `.CARI2` | pxDisplayText | 1 | ya | ALWAYS |
| `.CARI3` | pxDisplayText | 1 | ya | ALWAYS |
| `.ID` | pxDisplayText | 1 | ya | ALWAYS |
| `.Note` | pxDisplayText | 1 | ya | ALWAYS |
| `SearchAccumulation.Area` | pxAutoComplete | 1 | tidak | ALWAYS |
| `SearchAccumulation.CenterTransStatus` | pxTextInput | 1 | tidak | ALWAYS |
| `SearchAccumulation.City` | pxAutoComplete | 1 | tidak | ALWAYS |
| `SearchAccumulation.CZone` | pxAutoComplete | 1 | tidak | ALWAYS |
| `SearchAccumulation.District` | pxAutoComplete | 1 | tidak | ALWAYS |
| `SearchAccumulation.ID` | pxTextInput | 1 | tidak | ALWAYS |
| `SearchAccumulation.Keyword` | pxDropdown | 1 | tidak | ALWAYS |
| `SearchAccumulation.Nation` | pxAutoComplete | 1 | tidak | ALWAYS |
| `SearchAccumulation.Note` | pxTextInput | 1 | tidak | ALWAYS |
| `SearchAccumulation.PostalCode` | pxTextInput | 1 | tidak | ALWAYS |
| `SearchAccumulation.ProvinceName` | pxAutoComplete | 1 | tidak | ALWAYS |
| `SearchAccumulation.Type` | pxAutoComplete | 1 | tidak | ALWAYS |

#### SelClauseList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 2 (unik 2) - read-only 0 - RepeatGrid 1 - grid: .ClauseList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClauseContent` | (tanpa) | 1 | tidak | ALWAYS |
| `.ClauseParam` | (tanpa) | 1 | tidak | ALWAYS |

#### SelectAgent

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-lloydagent` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 1 - grid: pgRepPgSubSectionSelectAgentBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Address` | pxDisplayText | 1 | ya | NOTBLANK |
| `.City` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Country` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ID` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Name` | pxDisplayText | 1 | ya | NOTBLANK |

#### SelectCoverage

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Int-CONDITION` - field terikat 3 (unik 3) - read-only  - RepeatGrid 2 - grid: pgRepPgSubSectionSelectCoverageBB.pxResults, TempMasterPolis.OfferFacIn.CargoList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ConditionName` | (tanpa) | 1 | tidak | ALWAYS |
| `.ID` | (tanpa) | 1 | tidak | ALWAYS |
| `.pxListSubscript` | pxTextInput | 1 | ya | ALWAYS |

#### SelectPlanListRO

rumpun **PA-LIFE** - kelas dominan `ASM-FW-GISFW-Data-Plan` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: .PlanList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Limit` | (tanpa) | 1 | ya | ALWAYS |
| `.Name` | (tanpa) | 1 | ya | ALWAYS |

#### SelectShip

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Int-SHIP` - field terikat 16 (unik 16) - read-only 16 - RepeatGrid 1 - grid: pgRepPgSubSectionSelectShipBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Classification` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Construction` | pxDisplayText | 1 | ya | NOTBLANK |
| `.DWT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.FlagName` | pxDisplayText | 1 | ya | NOTBLANK |
| `.FormerShipName` | pxDisplayText | 1 | ya | NOTBLANK |
| `.GRT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ID` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Imo` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Name` | pxDisplayText | 1 | ya | NOTBLANK |
| `.NRT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Rebuilt` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Register` | pxDisplayText | 1 | ya | NOTBLANK |
| `.Remark` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ShipName` | pxDisplayText | 1 | ya | NOTBLANK |
| `.YearMake1` | pxDisplayText | 1 | ya | NOTBLANK |
| `.YearMake2` | pxDisplayText | 1 | ya | NOTBLANK |

#### SelectWarrantyList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Warranty` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 1 - grid: .WarrantyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.WarrantyDesc` | (tanpa) | 1 | tidak | ALWAYS |

#### SFAPortal_Opportunities

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### SFAPortal_OpportunitiesList

rumpun **INTI** - kelas dominan `ASM-FW-SFAGISFW-Work-Opportunity` - field terikat 21 (unik 8) - read-only 14 - RepeatGrid 3 - grid: pgRepPgSubSectionSFAPortal_OpportunitiesListBBBBBBBBBBBB.pxResults, pgRepPgSubSectionSFAPortal_OpportunitiesListBBBBBBBBBBBBBB.pxResults, pgRepPgSubSectionSFAPortal_OpportunitiesListBBBBBBBBBBBBBBBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.FilterTermForOpportunity` | pxTextInput | 1 | tidak | ALWAYS |
| `.Name` | pxLink | 3 | tidak | ALWAYS |
| `.TextNoQuotation` | pxDisplayText | 3 | ya | ALWAYS |
| `A.NBStatus` | pxDisplayText | 3 | ya | .pxPages(A).NBStatus != .pxPages(A).NBStatusNew |
| `A.NBStatusNew` | pxDisplayText | 2 | ya | ALWAYS |
| `A.Quotation.BusinessName` | pxDisplayText | 3 | ya | ALWAYS |
| `A.Quotation.InsuredName` | pxDisplayText | 3 | ya | ALWAYS |
| `A.Quotation.MarketingName` | (tanpa) | 3 | tidak | ALWAYS |

#### SFAPortal_OpportunitiesList_Header

rumpun **INTI** - kelas dominan `Data-Portal` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `SFAResults.pxResults(1).pySummaryCount(1)` | pxNumber | 1 | ya | ALWAYS |
| `SFAResults.pxResults(1).pySummaryDateTime(1)` | pxDateTime | 1 | ya | crmIsReview |
| `SFAResults.pxResults(1).pySummaryDateTime(2)` | pxDateTime | 1 | ya | crmIsReview |
| `SFAResults.pxResults(1).pySummaryValue(1)` | pxCurrency | 1 | ya | NOTBLANK |

#### SFAPortalOpportunitiesHeader

rumpun **INTI** - kelas dominan `Data-Portal` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CurrentMode` | pxHidden | 1 | tidak | ALWAYS |

#### ShowCedingCoList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Quotation` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: pyWorkPage.Quotation.CedingCoList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CedingCoName` | pxTextInput | 1 | ya [1!=2] | ALWAYS |

#### ShowCoverageFacOut

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 12 (unik 11) - read-only 11 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageList(1).FacOutObjectList(1).commision` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).ObjectPremi` | pxNumber | 1 | tidak | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).PercentOffered` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).Rate` | pxNumber | 2 | ya | IsAneka \|\| IsMarineCargo \|\| IsMBU \|\| IsTravel ; IsFire \|\| IsPA |
| `.CoverageList(1).FacOutObjectList(1).RiComm` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.CoverageList(1).FacOutTSI` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).NominalShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).PercentFacOut` | pxNumber | 1 | ya | ALWAYS |
| `.Currency` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |

#### ShowCoverageFacOut_IsUW

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Cargo` - field terikat 12 (unik 11) - read-only 12 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageList(1).Currency.Name` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).commision` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).ObjectPremi` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).Rate` | pxNumber | 2 | ya | IsAneka \|\| IsMarineCargo \|\| IsMBU \|\| IsTravel ; IsFire \|\| IsPA |
| `.CoverageList(1).FacOutObjectList(1).RiComm` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutTSI` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).NominalShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).PercentFacOut` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |

#### ShowCoveragePAFacOut_Dtl

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 18 (unik 17) - read-only 16 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMCoverage(1).FacOutObjectList(1).ObjectPremi` | pxNumber | 1 | tidak | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).PercentOffered` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).Rate` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).RiComm` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.ASMCoverage(1).FacOutTSI` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).NominalShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).PercentFacOut` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).ObjectPremi` | pxNumber | 1 | tidak | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).PercentOffered` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).Rate` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).RiComm` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.CoverageList(1).FacOutTSI` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).NominalShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageList(1).PercentFacOut` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 2 | ya | NOTBLANK |

#### ShowCoveragePAFacOutIsUW_Dtl

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMCoverage(1).FacOutObjectList(1).ObjectPremi` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).Rate` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).RiComm` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).FacOutTSI` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).NominalShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).PercentFacOut` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |

#### ShowCoverageTravelFacOut_Dtl

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 9 (unik 9) - read-only 8 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMCoverage(1).FacOutObjectList(1).ObjectPremi` | pxNumber | 1 | tidak | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).PercentOffered` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).Rate` | pxNumber | 1 | ya | 1=2 |
| `.ASMCoverage(1).FacOutObjectList(1).RiComm` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.ASMCoverage(1).FacOutTSI` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).NominalShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).PercentFacOut` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |

#### ShowCoverageTravelFacOut_Dtl_IsUW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMCoverage(1).FacOutObjectList(1).ObjectPremi` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).PercentOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).Rate` | pxNumber | 1 | ya | 1=2 |
| `.ASMCoverage(1).FacOutObjectList(1).RiComm` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=='1'] | ALWAYS |
| `.ASMCoverage(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).FacOutTSI` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).NominalShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.ASMCoverage(1).PercentFacOut` | pxNumber | 1 | ya | ALWAYS |
| `InputParam.HASIL12` | pxDisplayText | 1 | ya | NOTBLANK |

#### ShowObjectFacOut

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 38 (unik 23) - read-only 16 - RepeatGrid 8 - grid: .CargoList, .LocationList, .PersonList, .PropertyList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMDateOfBirth` | pxDateTime | 3 | tidak | ALWAYS |
| `.ASMGender` | pxDropdown | 3 | tidak | ALWAYS |
| `.ASMJobName` | pxDropdown | 3 | tidak | ALWAYS |
| `.BrandName` | pxDropdown | 1 | ya | ALWAYS |
| `.ConveyanceNote` | pxDropdown | 1 | tidak | ALWAYS |
| `.GoodNote` | pxInteger | 1 | tidak | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya | ALWAYS |
| `.ModelName` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectName` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectNo` | (tanpa) | 1 | ya | ALWAYS |
| `.PackingNote` | pxInteger | 1 | tidak | ALWAYS |
| `.Property.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.ObjectNo` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxTextInput | 2 | ya | ALWAYS |
| `.pyFullName` | (tanpa) | 3 | tidak | ALWAYS |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.SumTotalTSI` | pxNumber | 1 | ya | ALWAYS |
| `.TradingNote` | (tanpa) | 1 | tidak | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusaRe` | (tanpa) | 3 | tidak | ALWAYS |
| `.TSISpreaded` | pxNumber | 4 | campur | ALWAYS |
| `.TypeName` | pxTextInput | 1 | ya | ALWAYS |

#### ShowObjectFacOut_IsUW

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 28 (unik 19) - read-only 16 - RepeatGrid 6 - grid: .CargoList, .LocationList, .PersonList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMDateOfBirth` | pxDateTime | 2 | tidak | ALWAYS |
| `.ASMGender` | pxDropdown | 2 | tidak | ALWAYS |
| `.ASMJobName` | pxDropdown | 2 | tidak | ALWAYS |
| `.BrandName` | pxDropdown | 1 | ya | ALWAYS |
| `.ConveyanceNote` | pxDropdown | 1 | ya | ALWAYS |
| `.GoodNote` | pxInteger | 1 | ya | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya | ALWAYS |
| `.ModelName` | pxTextInput | 1 | ya | ALWAYS |
| `.PackingNote` | pxInteger | 1 | ya | ALWAYS |
| `.Property.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.Property.RiskLocation.ASMAddress` | pxTextInput | 2 | ya | ALWAYS |
| `.pxListSubscript` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.pyFullName` | (tanpa) | 2 | tidak | ALWAYS |
| `.TradingNote` | (tanpa) | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusaRe` | (tanpa) | 2 | tidak | ALWAYS |
| `.TSISpreaded` | pxNumber | 3 | campur | ALWAYS |
| `.TypeName` | pxTextInput | 1 | ya | ALWAYS |

#### ShowPolicyNoTreaty_SC

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-PolicyTreatyIn` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `pyWorkPage.PolicyTreatyIn.PolicyNo` | (tanpa) | 1 | ya | ALWAYS |
| `pyWorkPage.pyID` | (tanpa) | 1 | ya | ALWAYS |

#### ShowPolis_sc

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `pyWorkPage.pyID` | (tanpa) | 1 | ya | ALWAYS |
| `test.CARI40` | (tanpa) | 1 | ya | ALWAYS |

#### SourceHierarki

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-AGENT` - field terikat 5 (unik 4) - read-only 3 - RepeatGrid 1 - grid: pgRepPgSubSectionSourceHierarkiBBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClientID` | (tanpa) | 1 | ya | ALWAYS |
| `.ClientName` | (tanpa) | 1 | ya | ALWAYS |
| `.ID` | (tanpa) | 1 | ya | ALWAYS |
| `SearchSOB.CARI1` | pxAutoComplete | 2 | tidak | ALWAYS |

#### SpreadingCoverageAneka

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |

#### SpreadingCoverageAneka_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageNote` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |

#### SpreadingItem

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 10 (unik 8) - read-only 9 - RepeatGrid 1 - grid: .LayerList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.LayerNo` | (tanpa) | 1 | tidak | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 2 | ya | ALWAYS |
| `.PremiRp` | pxNumber | 1 | ya | IsMBU |
| `.Premium` | pxNumber | 2 | ya | !IsMBU |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSILiability` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### SpreadingItem_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 10 (unik 7) - read-only 10 - RepeatGrid 1 - grid: .LayerList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.LayerNo` | (tanpa) | 1 | ya | ALWAYS |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.Rate` | pxTextInput | 1 | ya | ALWAYS |
| `.TSI` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.TSILiability` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### SpreadingItemAneka

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### SpreadingItemAneka_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxNumber | 1 | ya | ALWAYS |

#### SpreadingItemLayer

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxDisplayText | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSILiability` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxDisplayText | 1 | ya | ALWAYS |

#### SpreadingItemLayer_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxDisplayText | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSILiability` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxDisplayText | 1 | ya | ALWAYS |

#### SpreadingItemLayerCommision

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PremiNusantaraRe` | pxDisplayText | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSILiability` | pxNumber | 1 | ya | ALWAYS |
| `.TSINusantaraRe` | pxDisplayText | 1 | ya | ALWAYS |

#### SpreadingObjectAneka

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 4 (unik 4) - read-only 3 - RepeatGrid 1 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | tidak | ALWAYS |
| `.ObjectName` | pxDisplayText | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TSIOld` | pxNumber | 1 | ya | ALWAYS |

#### SpreadingObjectAneka_IsUW

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |

#### SpreadingRiskList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-PolicyTreatyIn` - field terikat 13 (unik 13) - read-only 13 - RepeatGrid 1 - grid: .SpreadingRiskList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClaimPercentage` | pxNumber | 1 | ya [pyWorkPage.TreatyIn.FacultativeShare >0] | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 1 | ya [pyWorkPage.TreatyIn.FacultativeShare >0] | ALWAYS |
| `.TotalClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalPremium` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentageClaim` | pxNumber | 1 | ya | ALWAYS |
| `.TotalSharePercentagePremium` | pxNumber | 1 | ya | ALWAYS |
| `.TreatyType` | pxDropdown | 1 | ya [pyWorkPage.TreatyIn.FacultativeShare >0] | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber | 1 | ya | ALWAYS |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber | 1 | ya | ALWAYS |

#### SummaryCoverage_Section

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-PropertyItem` - field terikat 8 (unik 5) - read-only 8 - RepeatGrid 2 - grid: TempTotal.pxResults, TempTotalItem.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Currency` | pxNumber | 2 | ya | ALWAYS |
| `.ItemType` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.TotalNetRate` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |

#### SummaryLossRecord_Section

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 10 (unik 10) - read-only 10 - RepeatGrid 1 - grid: TotalLossCOB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClassOfBusiness` | pxNumber | 1 | ya | ALWAYS |
| `.Currency` | pxNumber | 1 | ya | ALWAYS |
| `.LOSS_AMOUNT` | pxNumber | 1 | ya | ALWAYS |
| `.LOSS_AMOUNT_IDR` | pxNumber | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.LossRatio1YearPercent` | pxTextInput | 1 | ya [1=1] | ALWAYS |
| `pyWorkPage.OfferFacIn.LossRatio3YearPercent` | pxTextInput | 1 | ya [1=1] | ALWAYS |
| `pyWorkPage.OfferFacIn.LossRatio5YearPercent` | pxTextInput | 1 | ya [1=1] | ALWAYS |
| `pyWorkPage.OfferFacIn.NetLossRatio1Year` | pxTextInput | 1 | ya [1=1] | ALWAYS |
| `pyWorkPage.OfferFacIn.NetLossRatio3Year` | pxTextInput | 1 | ya [1=1] | ALWAYS |
| `pyWorkPage.OfferFacIn.NetLossRatio5Year` | pxTextInput | 1 | ya [1=1] | ALWAYS |

#### SummaryRiskAccumulation

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Search` - field terikat 10 (unik 9) - read-only 9 - RepeatGrid 3 - grid: pgRepPgSubSectionSummaryRiskAccumulationBBBB.pxResults, ResultsAccumulation.pxResults, SummaryRiskAccum.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AccumulationName` | pxDisplayText | 1 | ya | ALWAYS |
| `.CARI1` | pxDisplayText | 2 | ya | ALWAYS |
| `.CARI2` | pxDisplayText | 1 | ya | ALWAYS |
| `.CARI3` | pxDisplayText | 1 | ya | ALWAYS |
| `.HASILD1` | pxNumber | 1 | ya | ALWAYS |
| `.HASILD2` | pxNumber | 1 | ya | ALWAYS |
| `.ID` | pxDisplayText | 1 | ya | ALWAYS |
| `.Note` | pxDisplayText | 1 | ya | ALWAYS |
| `InputPar.CARIYEAR` | pxTextInput | 1 | tidak | ALWAYS |

#### SummarySpreading_Section

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-SpreadingRisk` - field terikat 28 (unik 14) - read-only 24 - RepeatGrid 4 - grid: TotalSpreadAll.pxResults, TotalSpreadingCurrency.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ClaimAmountIDR` | pxNumber | 1 | ya | ALWAYS |
| `.ClaimEstimation` | pxNumber | 2 | ya | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.Name` | (tanpa) | 2 | tidak | ALWAYS |
| `.OldTSI` | pxNumber | 2 | ya | ALWAYS |
| `.Premium` | pxNumber | 2 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 2 | ya | OutputParam.ERRMSG6=='' |
| `.TreatyName` | (tanpa) | 2 | tidak | ALWAYS |
| `.TreatyType` | pxDropdown | 2 | ya | ALWAYS |
| `.TSI` | pxNumber | 2 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 2 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 2 | ya | ALWAYS |
| `ProtectSpreading.CARI14` | pxDisplayText | 3 | ya | ALWAYS |

#### TableInwardScale

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn` - field terikat 8 (unik 8) - read-only 4 - RepeatGrid 1 - grid: D_InwardScale.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.PctTreatyLimit` | pxNumber | 1 | tidak | ALWAYS |
| `.ShareMax` | pxNumber | 1 | tidak | ALWAYS |
| `.ShareMin` | pxNumber | 1 | tidak | ALWAYS |
| `.Tahun` | pxTextInput | 1 | tidak | ALWAYS |
| `pyWorkPage.OfferFacIn.CurrentYear` | (tanpa) | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.Parameters.CariScaleMax` | (tanpa) | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.Parameters.CariScaleMin` | (tanpa) | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.Parameters.CariTreatyLimit` | (tanpa) | 1 | ya | ALWAYS |

#### TableOfLimit

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-TABLEOFLIMIT` - field terikat 6 (unik 6) - read-only 3 - RepeatGrid 1 - grid: pgRepPgSubSectionTableOfLimitBBB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Category` | pxTextInput | 1 | tidak | ALWAYS |
| `.Description` | pxTextInput | 1 | tidak | ALWAYS |
| `.PctLimit` | pxTextInput | 1 | tidak | ALWAYS |
| `pyWorkPage.OfferFacIn.Parameters.CariOccupationConstruction` | (tanpa) | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.Parameters.CariOccupationRiskCategory` | (tanpa) | 1 | ya | ALWAYS |
| `pyWorkPage.OfferFacIn.Parameters.LowestPctLimit` | (tanpa) | 1 | ya | ALWAYS |

#### TotalAccumulation_FacIn

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Search` - field terikat 5 (unik 5) - read-only 2 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CARI10` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI15` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI16` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARID1` | pxNumber | 1 | ya | ALWAYS |
| `.CARID2` | pxNumber | 1 | ya | ALWAYS |

#### TotalAccumulationDtl

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Search` - field terikat 9 (unik 9) - read-only 3 - RepeatGrid 1 - grid: TotalAccumulationDtl.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CARI11` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI15` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI16` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI17` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI18` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI19` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARID1` | pxNumber | 1 | ya | ALWAYS |
| `.CARID2` | pxNumber | 1 | ya | ALWAYS |
| `InputAccumulationCode.CARI10` | pxTextInput | 1 | ya | ALWAYS |

#### TSIForFacOut

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |

#### TSIPropertyFacOut

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |

#### UploadMember

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Data-Policy` - field terikat 10 (unik 10) - read-only 8 - RepeatGrid 3 - grid: .BatchMemberList, .FinalOccupationList, .PolicyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BookNo` | pxLink | 1 | tidak | ALWAYS |
| `.EmployeeName` | pxDisplayText | 1 | ya | ALWAYS |
| `.Name` | pxDisplayText | 1 | ya | ALWAYS |
| `.OccupationName` | pxDisplayText | 1 | ya | ALWAYS |
| `.PlanIP` | pxDisplayText | 1 | ya | ALWAYS |
| `.PlanOP` | pxDisplayText | 1 | ya | ALWAYS |
| `.pxCreateDateTime` | pxDisplayText | 1 | ya | ALWAYS |
| `.pxCreateOpName` | pxDisplayText | 1 | ya | ALWAYS |
| `.pyLabel` | pxDisplayText | 1 | ya | ALWAYS |
| `.SPPANo` | pxLink | 1 | tidak | ALWAYS |

#### VehicleGrid

rumpun **MBU** - kelas dominan `ASM-FW-GISFW-Data-Vehicle` - field terikat 17 (unik 15) - read-only 17 - RepeatGrid 1 - grid: .AccessoryList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BrandName` | pxAutoComplete | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.CarParking` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.CC` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ChassisNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ColorName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.EngineNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ModelName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Occupation.OccupationName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Price` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TransmissionType` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Type` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | 1=2 |
| `.TypeName` | pxAutoComplete | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ValueName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### VehicleGrid_FacIn

rumpun **MBU** - kelas dominan `ASM-FW-GISFW-Data-Vehicle` - field terikat 17 (unik 15) - read-only 17 - RepeatGrid 1 - grid: .AccessoryList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BrandName` | pxAutoComplete | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.CarParking` | pxDropdown | 1 | ya [IsSpreadingUW] | ALWAYS |
| `.CC` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ChassisNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ColorName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.EngineNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ModelName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Occupation.OccupationName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Price` | pxNumber | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.TransmissionType` | pxDropdown | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Type` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | 1=2 |
| `.TypeName` | pxAutoComplete | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.ValueName` | pxAutoComplete | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### VehicleGrid_FacIn_IsUW

rumpun **MBU** - kelas dominan `ASM-FW-GISFW-Data-Vehicle` - field terikat 17 (unik 15) - read-only 17 - RepeatGrid 1 - grid: .AccessoryList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BrandName` | pxAutoComplete | 2 | ya | ALWAYS |
| `.CarParking` | pxDropdown | 1 | ya | ALWAYS |
| `.CC` | pxTextInput | 1 | ya | ALWAYS |
| `.ChassisNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.ColorName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.EngineNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya | ALWAYS |
| `.ModelName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Occupation.OccupationName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Price` | pxNumber | 1 | ya | ALWAYS |
| `.TransmissionType` | pxDropdown | 1 | ya | ALWAYS |
| `.Type` | pxTextInput | 1 | ya | 1=2 |
| `.TypeName` | pxAutoComplete | 2 | ya | ALWAYS |
| `.ValueName` | pxAutoComplete | 1 | ya | ALWAYS |

#### ViewAlamatKirim

rumpun **INTI** - kelas dominan `Data-Address` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMAddressType` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMContactInfo` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.pyPhoneNumber` | pxTextInput | 1 | ya | ALWAYS |

#### ViewCheckListOffer

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 33 (unik 27) - read-only 14 - RepeatGrid 7 - grid: .Policy.FacOfferList, .PropertyList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BrandName` | pxTextInput | 2 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 4 | ya | ALWAYS |
| `.Occupation.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `FacOfferList.Attention` | pxTextInput | 1 | ya | ALWAYS |
| `FacOfferList.OurRef` | pxTextInput | 1 | ya | ALWAYS |
| `FacOfferList.ReinsurerName` | pxTextInput | 1 | ya | ALWAYS |
| `InputData.CARI11` | pxNumber | 2 | ya | IsFire ; IsMBU |
| `InputData.CARI34` | pxDateTime | 1 | ya | ALWAYS |
| `ViewOfferStatus.AdditionalInfo` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.ASMShare` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.ClassOfConstruction` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.Clauses` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.Deductibles` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.Information` | pxTextArea | 1 | tidak | InputData.CARI5!=5 && InputData.CARI6!=6 && InputData.CARI17!=17 && In... |
| `ViewOfferStatus.Insured` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.LineBusiness` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.Location` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.LossProtection` | pxCheckbox | 2 | tidak | ALWAYS |
| `ViewOfferStatus.Occupation` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.OurRetention` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.PerilsCovered` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.PeriodOfInsurance` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.PeriodOfReinsurance` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.ShareOffered` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.SumInsuredTSI` | pxCheckbox | 1 | tidak | ALWAYS |
| `ViewOfferStatus.SurveyReport` | pxCheckbox | 1 | tidak | ALWAYS |

#### ViewCheckListOfferFacOut

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Property` - field terikat 12 (unik 7) - read-only 4 - RepeatGrid 7 - grid: .Policy.FacOfferList, .PropertyList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BrandName` | pxTextInput | 2 | tidak | ALWAYS |
| `.ObjectName` | pxTextInput | 4 | campur | ALWAYS |
| `.Occupation.OccupationName` | pxTextInput | 1 | tidak | ALWAYS |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | tidak | ALWAYS |
| `.Vehicle.FacOut` | pxDropdown | 1 | tidak | ALWAYS |
| `InputData.CARI11` | pxNumber | 2 | ya | IsFire ; IsMBU |
| `InputData.CARI33` | pxTextArea | 1 | tidak | ALWAYS |

#### ViewClaimList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Search` - field terikat 9 (unik 9) - read-only 9 - RepeatGrid 1 - grid: TempListClaim.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CARI11` | (tanpa) | 1 | ya | ALWAYS |
| `.CARI12` | (tanpa) | 1 | ya | ALWAYS |
| `.CARI13` | (tanpa) | 1 | ya | ALWAYS |
| `.CARI14` | (tanpa) | 1 | ya | ALWAYS |
| `.CARI15` | pxNumber | 1 | ya | ALWAYS |
| `.CARI16` | pxTextInput | 1 | ya | ALWAYS |
| `.CARI17` | pxTextInput | 1 | ya | ALWAYS |
| `.CARI18` | (tanpa) | 1 | ya | ALWAYS |
| `.CARI20` | (tanpa) | 1 | ya | ALWAYS |

#### ViewClauseFire

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Clause` - field terikat 8 (unik 8) - read-only 7 - RepeatGrid 1 - grid: .ArgumentList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ArgumentCount` | pxTextInput | 1 | ya | .ArgumentCount > 0 |
| `.ArgumentDescription` | pxTextInput | 1 | ya | ALWAYS |
| `.ArgumentNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.ArgumentValue` | pxTextInput | 1 | tidak | ALWAYS |
| `.ClauseCode` | pxTextInput | 1 | ya | ALWAYS |
| `.ClauseDescription` | pxTextInput | 1 | ya | ALWAYS |
| `.ClauseLanguageID` | pxDropdown | 1 | ya | ALWAYS |
| `.ClauseTitle` | pxTextInput | 1 | ya | ALWAYS |

#### ViewCoins

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coins` - field terikat 15 (unik 15) - read-only 15 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Brokerage` | pxNumber | 1 | ya | ALWAYS |
| `.CoinsName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.DiscountShare` | pxNumber | 1 | ya | ALWAYS |
| `.HandingFee` | pxNumber | 1 | ya | ALWAYS |
| `.Leader` | pxCheckbox | 1 | ya | ALWAYS |
| `.LeaderPolicyNo` | pxTextInput | 1 | ya | ALWAYS |
| `.PercentBrokerage` | pxNumber | 1 | ya | ALWAYS |
| `.PercentHandlingFee` | pxNumber | 1 | ya | ALWAYS |
| `.PercentPPH` | pxNumber | 1 | ya | ALWAYS |
| `.PercentPPN` | pxNumber | 1 | ya | ALWAYS |
| `.PercentShare` | pxNumber | 1 | ya | ALWAYS |
| `.PPH` | pxNumber | 1 | ya | ALWAYS |
| `.PPN` | pxNumber | 1 | ya | ALWAYS |
| `.PremiShare` | pxNumber | 1 | ya | ALWAYS |
| `.TSIShare` | pxNumber | 1 | ya | ALWAYS |

#### ViewCompanyInformationFacOut

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 7 (unik 7) - read-only 7 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Customer_C.ASMClientID` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_C.ASMDirectorPerson.pyFullName` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_C.ASMNPWP` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_C.pyCompany` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.pyTitle` | pxDropdown | 1 | ya | ALWAYS |
| `InputData.CARI10` | pxTextInput | 1 | ya | ALWAYS |
| `InputData.CARI9` | pxTextInput | 1 | ya | ALWAYS |

#### ViewCoverage

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 1 - grid: .CoverageList

_Tidak ada properti terikat._

#### ViewCoverageAneka

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 13 (unik 13) - read-only 13 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CalculateMethod` | pxDropdown | 1 | ya | ALWAYS |
| `.CoverageNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Day` | pxRadioButtons | 1 | ya | ALWAYS |
| `.Discount` | pxNumber | 1 | ya | ALWAYS |
| `.DiscountPercentage` | pxTextInput | 1 | ya | ALWAYS |
| `.Loading` | pxTextInput | 1 | ya | ALWAYS |
| `.MinPremium` | pxNumber | 1 | ya | ALWAYS |
| `.PctShortPeriod` | pxTextInput | 1 | ya [IsGroupUWFac \|\| IsSpreadingUW] | isEDM && .CalculateMethod==2 |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.ProRatePercent` | pxTextInput | 1 | ya | ALWAYS |
| `.Rate` | pxTextInput | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `.TypeOfDiscount` | pxCheckbox | 1 | ya | ALWAYS |

#### ViewCoverageCargoFacOut

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 8 (unik 8) - read-only 8 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Coverage` | pxAutoComplete | 1 | ya | ALWAYS |
| `.CoverageNote` | pxTextArea | 1 | ya | ALWAYS |
| `.Discount` | pxNumber | 1 | ya | ALWAYS |
| `.DiscountPercentage` | pxNumber | 1 | ya | ALWAYS |
| `.MinPremium` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |

#### ViewCoverageFacOutShow

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageBasis` | pxDropdown | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.CoverageNote` | pxTextInput | 1 | ya | ALWAYS |
| `.PremiumRetro` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=="1"] | ALWAYS |
| `.Rate` | pxNumber | 1 | ya [pyWorkPage.IsFacRetroOffer=="1"] | ALWAYS |

#### ViewCoverageFacOutShow_isUW

rumpun **INTI+FIRE** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 5 (unik 5) - read-only 5 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CoverageBasis` | pxDropdown | 1 | ya | ALWAYS |
| `.CoverageList(1).FacOutObjectList(1).ShareOffered` | pxNumber | 1 | ya | ALWAYS |
| `.OLDID` | pxTextInput | 1 | ya | ALWAYS |
| `.PremiumRetro` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |

#### ViewCoverageFire

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .CoverageList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Coverage` | pxTextInput | 1 | ya | ALWAYS |
| `.CoverageNote` | pxTextInput | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |

#### ViewCoveragePA

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 1 - grid: .ASMCoverage

_Tidak ada properti terikat._

#### ViewCoverageTravel

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 1 - grid: .ASMCoverage

_Tidak ada properti terikat._

#### ViewCSVAneka

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-TABLE_B2B` - field terikat 22 (unik 22) - read-only 22 - RepeatGrid 1 - grid: TempWorkPage.ListB2BHostUpload

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ADDRESS` | pxDisplayText | 1 | ya | NOTBLANK |
| `.BUILDING_NO` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CALCULATE_METHOD` | pxDisplayText | 1 | ya | NOTBLANK |
| `.COUNTRY` | pxDisplayText | 1 | ya | NOTBLANK |
| `.COVERAGE_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CURRENCY` | pxDisplayText | 1 | ya | NOTBLANK |
| `.DAY` | (tanpa) | 1 | ya | ALWAYS |
| `.DISCOUNT` | pxDropdown | 1 | ya | ALWAYS |
| `.LOL` | pxDropdown | 1 | ya | ALWAYS |
| `.OBJECT_DESCRIPTION` | pxDisplayText | 1 | ya | NOTBLANK |
| `.OCCUPATION_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.PCT_DISCOUNT` | pxNumber | 1 | ya | ALWAYS |
| `.PCT_LOADING` | pxNumber | 1 | ya | ALWAYS |
| `.PCT_LOL` | pxNumber | 1 | ya | ALWAYS |
| `.PCT_RATE` | pxNumber | 1 | ya | ALWAYS |
| `.PCT_SCALE` | pxNumber | 1 | ya | ALWAYS |
| `.POSTAL_CODE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.PREMIUM` | pxNumber | 1 | ya | ALWAYS |
| `.pxListSubscript` | pxDisplayText | 1 | ya | NOTBLANK |
| `.SECTION` | (tanpa) | 1 | ya | NOTBLANK |
| `.TSI` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TYPE_ADDRESS` | pxDisplayText | 1 | ya | NOTBLANK |

#### ViewCSVCredit

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-TABLE_B2B` - field terikat 26 (unik 26) - read-only 26 - RepeatGrid 1 - grid: TempWorkPage.ListB2BHostUpload

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CEDING_CO_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CLASS_OF_BUSINESS` | pxDisplayText | 1 | ya | NOTBLANK |
| `.COVERAGE_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CURRENCY` | pxDisplayText | 1 | ya | NOTBLANK |
| `.DATE_END` | pxDisplayText | 1 | ya | NOTBLANK |
| `.FACULTATIVE_TYPE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.INSURED_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.JOB` | pxDisplayText | 1 | ya | NOTBLANK |
| `.LOAN_ACCOUNT_NO` | pxDisplayText | 1 | ya | NOTBLANK |
| `.LOAN_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.LOAN_TYPE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.MARKETING_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.NO` | pxDisplayText | 1 | ya | NOTBLANK |
| `.NO_OFFER_SLIP` | pxDisplayText | 1 | ya | NOTBLANK |
| `.OBJECT_DESCRIPTION` | pxDisplayText | 1 | ya | NOTBLANK |
| `.OCCUPATION_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.PCT_COMM` | pxNumber | 1 | ya | NOTBLANK |
| `.PCT_SHARE_RNM` | pxNumber | 1 | ya | NOTBLANK |
| `.PREMIUM_COVERAGE` | pxNumber | 1 | ya | NOTBLANK |
| `.PREMIUM_RNM` | pxNumber | 1 | ya | NOTBLANK |
| `.PRODUCT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.RATE_COVERAGE` | pxNumber | 1 | ya | NOTBLANK |
| `.SHARE_RNM` | pxNumber | 1 | ya | NOTBLANK |
| `.SOB_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.START_DATE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TSI_OBJECT` | pxNumber | 1 | ya | NOTBLANK |

#### ViewCSVMarine

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-TABLE_B2B` - field terikat 52 (unik 52) - read-only 52 - RepeatGrid 1 - grid: TempWorkPage.ListB2BHostUpload

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AMOUNT1` | (tanpa) | 1 | ya | ALWAYS |
| `.AMOUNT2` | (tanpa) | 1 | ya | ALWAYS |
| `.AMOUNT3` | (tanpa) | 1 | ya | ALWAYS |
| `.AMOUNT4` | (tanpa) | 1 | ya | ALWAYS |
| `.AMOUNT5` | (tanpa) | 1 | ya | ALWAYS |
| `.CERTNO` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CONDITION_NAME1` | (tanpa) | 1 | ya | ALWAYS |
| `.CONDITION_NAME2` | (tanpa) | 1 | ya | ALWAYS |
| `.CONDITION_NAME3` | (tanpa) | 1 | ya | ALWAYS |
| `.CONDITION_NAME4` | (tanpa) | 1 | ya | ALWAYS |
| `.CONDITION_NAME5` | (tanpa) | 1 | ya | ALWAYS |
| `.CONDITION1` | pxDropdown | 1 | ya | ALWAYS |
| `.CONDITION2` | pxDropdown | 1 | ya | ALWAYS |
| `.CONDITION3` | pxDropdown | 1 | ya | ALWAYS |
| `.CONDITION4` | pxDropdown | 1 | ya | ALWAYS |
| `.CONDITION5` | pxDropdown | 1 | ya | ALWAYS |
| `.CONVEYENCE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CONVEYENCE_NAME` | (tanpa) | 1 | ya | NOTBLANK |
| `.CURRENCY1` | (tanpa) | 1 | ya | ALWAYS |
| `.CURRENCY2` | (tanpa) | 1 | ya | ALWAYS |
| `.CURRENCY3` | (tanpa) | 1 | ya | ALWAYS |
| `.CURRENCY4` | (tanpa) | 1 | ya | ALWAYS |
| `.CURRENCY5` | (tanpa) | 1 | ya | ALWAYS |
| `.DEPARTURE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.DESTINATION` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ETA` | pxDisplayText | 1 | ya | NOTBLANK |
| `.GOODS` | pxDisplayText | 1 | ya | NOTBLANK |
| `.GOODS_NAME` | (tanpa) | 1 | ya | ALWAYS |
| `.GROSSPREMIUM` | pxNumber | 1 | ya | ALWAYS |
| `.MINMAX1` | pxDropdown | 1 | ya | ALWAYS |
| `.MINMAX2` | pxDropdown | 1 | ya | ALWAYS |
| `.MINMAX3` | pxDropdown | 1 | ya | ALWAYS |
| `.MINMAX4` | pxDropdown | 1 | ya | ALWAYS |
| `.MINMAX5` | pxDropdown | 1 | ya | ALWAYS |
| `.PCTDEDUCTIBLE1` | (tanpa) | 1 | ya | ALWAYS |
| `.PCTDEDUCTIBLE2` | (tanpa) | 1 | ya | ALWAYS |
| `.PCTDEDUCTIBLE3` | (tanpa) | 1 | ya | ALWAYS |
| `.PCTDEDUCTIBLE4` | (tanpa) | 1 | ya | ALWAYS |
| `.PCTDEDUCTIBLE5` | (tanpa) | 1 | ya | ALWAYS |
| `.pxListSubscript` | pxDisplayText | 1 | ya | NOTBLANK |
| `.RATE` | pxNumber | 1 | ya | ALWAYS |
| `.STATUS` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TOTALGOODS` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TOTALPRICE` | (tanpa) | 1 | ya | ALWAYS |
| `.TRADING` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TYPEDEDUCTIBLE1` | pxDropdown | 1 | ya | ALWAYS |
| `.TYPEDEDUCTIBLE2` | pxDropdown | 1 | ya | ALWAYS |
| `.TYPEDEDUCTIBLE3` | pxDropdown | 1 | ya | ALWAYS |
| `.TYPEDEDUCTIBLE4` | pxDropdown | 1 | ya | ALWAYS |
| `.TYPEDEDUCTIBLE5` | pxDropdown | 1 | ya | ALWAYS |
| `.UNITY` | pxDisplayText | 1 | ya | NOTBLANK |
| `.VESSELNAME` | pxDisplayText | 1 | ya | NOTBLANK |

#### ViewCSVResult_TableB2B

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-TABLE_B2B` - field terikat 43 (unik 43) - read-only 43 - RepeatGrid 1 - grid: TempWorkPage.ListB2BHostUpload

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ACCUMULATION_CODE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ACCUMULATION_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.BEGIN_DATE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CEDING_CO_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CLASS_OF_BUSINESS` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CLASS_OF_CONSTRUCTION` | pxDisplayText | 1 | ya | NOTBLANK |
| `.COMMISION` | pxDisplayText | 1 | ya | NOTBLANK |
| `.COVERAGE_BASIS` | pxDisplayText | 1 | ya | NOTBLANK |
| `.COVERAGE_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.COVERAGE_PREMIUM` | pxDisplayText | 1 | ya | NOTBLANK |
| `.COVERAGE_RATE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.CURRENCY` | pxDisplayText | 1 | ya | NOTBLANK |
| `.DEDUCTIBLE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.END_DATE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.FACULTATIVE_TYPE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.INDEMNITY` | pxDisplayText | 1 | ya | NOTBLANK |
| `.INSURED_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.INTEREST` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ITEM_TYPE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.LAYER` | pxDisplayText | 1 | ya | NOTBLANK |
| `.LIMIT_OF_LIABILITY` | pxDisplayText | 1 | ya | NOTBLANK |
| `.LOSSLIMIT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.MARKETING_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.NO_OFFER_SLIP` | pxDisplayText | 1 | ya | NOTBLANK |
| `.NOTES` | pxDisplayText | 1 | ya | NOTBLANK |
| `.OBJECT_TYPE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.OCCUPATION_CODE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.OCCUPATION_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.POSTAL_CODE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.RATE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.RISK_LOCATION` | pxDisplayText | 1 | ya | NOTBLANK |
| `.RNM_PREMIUM` | pxDisplayText | 1 | ya | NOTBLANK |
| `.RNM_SHARE` | pxDisplayText | 1 | ya | NOTBLANK |
| `.RNM_SHARE_PCT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.SOB_NAME` | pxDisplayText | 1 | ya | NOTBLANK |
| `.SUBLIMIT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.SUM_INSURED_BI` | pxDisplayText | 1 | ya | NOTBLANK |
| `.SUM_INSURED_MD` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TSI_LOSSLIMIT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TSI_MD_BI` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TSI_OBJECT_ITEM` | pxDisplayText | 1 | ya | NOTBLANK |
| `.TSI_SUBLIMIT` | pxDisplayText | 1 | ya | NOTBLANK |
| `.ZONE` | pxDisplayText | 1 | ya | NOTBLANK |

#### ViewDataPolicy

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 16 (unik 16) - read-only 14 - RepeatGrid 1 - grid: TotalSpreadAll.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Name` | (tanpa) | 1 | ya | ALWAYS |
| `.OfferFacIn.FacRetroDetails.BackUpStatus` | pxDropdown | 1 | tidak | 1=2 |
| `.OfferFacIn.OfficerFacOut` | pxTextInput | 1 | ya | 1=2 |
| `.OfferFacIn.PercentShare` | pxNumber | 1 | ya | ALWAYS |
| `.OfferFacIn.PolicyData.EndDateTime` | pxDateTime | 1 | ya [.pyID != ""] | ALWAYS |
| `.OfferFacIn.PolicyData.StartDateTime` | pxDateTime | 1 | ya [.pyID != ""] | ALWAYS |
| `.OfferFacIn.QuotationData.BusinessName` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.QuotationData.CedingCoName` | pxDisplayText | 1 | ya | pyWorkPage.IsCedingConfirm!='Offer' |
| `.OfferFacIn.QuotationData.EdmDate` | pxDateTime | 1 | ya [.pyID != ""] | ALWAYS |
| `.OfferFacIn.QuotationData.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.QuotationData.OldPolicyNo` | pxTextInput | 1 | ya | IsEDM |
| `.OfferFacIn.QuotationData.SobName` | pxDisplayText | 1 | ya | pyWorkPage.IsCedingConfirm!='Offer' |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |
| `FacOfferList.AdditionalInfo` | pxTextArea | 1 | tidak | 1=2 |
| `InputData.CARI7` | pxTextInput | 1 | ya | 1=2 |

#### ViewDataPolicy_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 11 (unik 11) - read-only 11 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OfferFacIn.FacRetroDetails.BackUpStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.OfferFacIn.OfficerFacOut` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.PercentShare` | pxNumber | 1 | ya | ALWAYS |
| `.OfferFacIn.PolicyData.EndDateTime` | pxDateTime | 1 | ya [.pyID != ""] | ALWAYS |
| `.OfferFacIn.PolicyData.StartDateTime` | pxDateTime | 1 | ya [.pyID != ""] | ALWAYS |
| `.OfferFacIn.QuotationData.BusinessName` | pxTextInput | 1 | ya | ALWAYS |
| `.OfferFacIn.QuotationData.InsuredName` | pxTextInput | 1 | ya | ALWAYS |
| `.pyID` | pxTextInput | 1 | ya | ALWAYS |
| `FacOfferList.AdditionalInfo` | pxTextArea | 1 | ya | 1=2 |
| `FacOfferList.FacNo` | pxTextInput | 1 | ya | ALWAYS |
| `InputData.CARI7` | pxTextInput | 1 | ya | ALWAYS |

#### ViewDeductible

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 10 (unik 8) - read-only 10 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BasisType` | pxDropdown | 1 | ya | .DeductibleType = 4 \|\| .DeductibleType = 7 \|\| .DeductibleType = 9 |
| `.ClaimCategory` | pxDropdown | 1 | ya | !IsFire |
| `.Coverage` | pxDropdown | 2 | ya | ALWAYS |
| `.Currency` | pxDropdown | 1 | ya | IsDeducType |
| `.DeductibleType` | pxDropdown | 1 | ya | ALWAYS |
| `.Descriptions` | pxTextArea | 2 | ya | IsDescDeductType |
| `.Type` | (tanpa) | 1 | ya | IsTypeDeductType |
| `.Value` | pxTextInput | 1 | ya | IsValueDeductType |

#### ViewDeliveryAddressFacOut

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### ViewDetailCargo

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Cargo` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ConveyanceNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.GoodNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.PackingNote` | pxAutoComplete | 1 | ya | ALWAYS |
| `.TradingNote` | pxAutoComplete | 1 | ya | ALWAYS |

#### ViewDetailPayment_SC

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Payment` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .DetailPayment

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.GLDate` | pxDateTime | 1 | ya | ALWAYS |
| `.PaymentAmount` | pxNumber | 1 | ya | ALWAYS |
| `.PaymentDate` | (tanpa) | 1 | ya | ALWAYS |
| `.PaymentId` | (tanpa) | 1 | ya | ALWAYS |

#### ViewDtlAddCargo

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: .AdditionalCoverage

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Coverage` | pxAutoComplete | 1 | ya | ALWAYS |
| `.CoverageNote` | pxTextArea | 1 | ya | ALWAYS |

#### ViewDtlAdditional

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Coverage` - field terikat 6 (unik 6) - read-only 6 - RepeatGrid 1 - grid: .AdditionalCoverage

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CalculateMethod` | pxDropdown | 1 | ya | ALWAYS |
| `.CoverageNote` | pxDropdown | 1 | ya | ALWAYS |
| `.Discount` | pxNumber | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.Rate` | pxNumber | 1 | ya | ALWAYS |
| `.TSI` | pxNumber | 1 | ya | ALWAYS |

#### ViewDtlCoverageFacOut

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 24 (unik 23) - read-only 23 - RepeatGrid 6 - grid: .CargoList, .LocationList, .PersonList, .PropertyList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 1 | ya | ALWAYS |
| `.ASMHeight` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMJobName` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMWeight` | pxInteger | 1 | ya | ALWAYS |
| `.BrandName` | pxDropdown | 1 | ya | ALWAYS |
| `.ConveyanceNote` | pxTextInput | 1 | ya | ALWAYS |
| `.GoodNote` | pxTextInput | 1 | ya | ALWAYS |
| `.IsMaterialDamage` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 1 | ya | IsFac |
| `.LicensePlate` | pxTextInput | 1 | ya | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya | ALWAYS |
| `.ModelName` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.PackingNote` | pxTextInput | 1 | ya | ALWAYS |
| `.pyFullName` | (tanpa) | 2 | ya | ALWAYS |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.TradingNote` | pxTextInput | 1 | ya | ALWAYS |
| `.TypeName` | pxTextInput | 1 | ya | ALWAYS |

#### ViewDtlDeductible

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 1 - grid: .DeductibleList

_Tidak ada properti terikat._

#### ViewDtlMemorandumAccept

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-MemorandumAccept` - field terikat 15 (unik 15) - read-only 15 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Accumulation` | pxTextInput | 1 | ya | ALWAYS |
| `.ApplicationNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.BuildingStatus` | pxTextArea | 1 | ya | ALWAYS |
| `.Collateral` | pxTextInput | 1 | ya | ALWAYS |
| `.ConclusionsProposals` | pxTextArea | 1 | ya | ALWAYS |
| `.InsuredSiup` | pxTextInput | 1 | ya | ALWAYS |
| `.KeyPerson` | pxTextInput | 1 | ya | ALWAYS |
| `.OriginalRate` | pxTextInput | 1 | ya | ALWAYS |
| `.PrincipalsBackground` | pxTextArea | 1 | ya | ALWAYS |
| `.PrincipalsData` | pxTextArea | 1 | ya | ALWAYS |
| `.PrincipalsGuarantee` | pxTextArea | 1 | ya | ALWAYS |
| `.PrincipalsOutGo` | pxTextInput | 1 | ya | ALWAYS |
| `.SummaryFinancial` | pxTextArea | 1 | ya | ALWAYS |
| `.SupportingDocument` | pxTextArea | 1 | ya | ALWAYS |
| `.TypeOfWording` | pxTextInput | 1 | ya | ALWAYS |

#### ViewDtlObjectAneka

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 52 (unik 37) - read-only 52 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.AreaHectar` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.MaxPassengers` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.Pilots` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.RegistrationMarks` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.SpecialRentalUses` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.SpecialUses` | pxTextInput | 1 | ya | ALWAYS |
| `.AviationHull.StandardUses` | pxTextInput | 1 | ya | ALWAYS |
| `.BuiltIn` | pxTextInput | 1 | ya | ALWAYS |
| `.DraftWordingID` | pxDropdown | 1 | ya | IsBonding |
| `.Flag` | pxTextInput | 1 | ya | ALWAYS |
| `.GeographicalLimits` | pxTextInput | 1 | ya | ALWAYS |
| `.InvoiceNo` | pxTextInput | 1 | ya | ALWAYS |
| `.LimitIdemnityDay` | pxTextInput | 1 | ya | ALWAYS |
| `.LimitIdemnityMonth` | pxTextInput | 1 | ya | ALWAYS |
| `.LimitOfLiability` | pxTextInput | 1 | ya | IsObjectLimiliabilityAneka |
| `.MarineHull.BREADTHM` | pxTextInput | 1 | ya | ALWAYS |
| `.MarineHull.Classification` | pxTextInput | 1 | ya | ALWAYS |
| `.MarineHull.Construction` | pxTextInput | 1 | ya | ALWAYS |
| `.MarineHull.DEPTHM` | pxTextInput | 1 | ya | ALWAYS |
| `.MarineHull.DWTRT` | pxTextInput | 1 | ya | ALWAYS |
| `.MarineHull.GRTRT` | pxTextInput | 1 | ya | ALWAYS |
| `.MarineHull.LENGTHM` | pxTextInput | 1 | ya | ALWAYS |
| `.MarineHull.NRTRT` | pxTextInput | 1 | ya | ALWAYS |
| `.MarineHull.TypeOfVessel` | pxTextInput | 1 | ya | ALWAYS |
| `.NoOfTree` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextArea/pxTextInput | 8 | ya | ALWAYS |
| `.Quantity` | pxTextInput | 3 | ya | pyWorkPage.Quotation.BusinessCode!=20 |
| `.Section` | pxDropdown | 2 | ya [IsSpreadingUW] | IsObjectSectionAneka |
| `.SerialNo` | pxTextInput | 1 | ya | ALWAYS |
| `.TSI` | pxTextInput | 2 | ya | ALWAYS |
| `.UnitQuantity` | pxTextInput | 1 | ya | pyWorkPage.Quotation.BusinessCode!=20 |
| `.VehicleHE.BrandName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.VehicleHE.ChassisNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.VehicleHE.EngineNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.VehicleHE.ObjectNameHE` | pxAutoComplete | 1 | ya | ALWAYS |
| `.VehicleHE.TypeName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.Year` | pxTextInput | 5 | ya | ALWAYS |

#### ViewDtlObjectLocation

rumpun **INTI+FIRE** - kelas dominan `Data-Address` - field terikat 7 (unik 7) - read-only 5 - RepeatGrid 1 - grid: .OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.ASMCity` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMDistrict` | pxDropdown | 1 | ya | .ASMCity != '' |
| `.ASMRW` | pxDropdown | 1 | ya | .ASMDistrict != '' |
| `.ASMZipCode` | pxTextInput | 1 | ya | ALWAYS |
| `.OccupationId` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### ViewDtlOtherObjectAneka

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 119 (unik 63) - read-only 118 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Aneka.Business` | pxTextArea | 1 | ya | isBisnisOtherObjAneka |
| `.Aneka.BusinessGolf` | pxDropdown | 1 | ya | isGolfInsurance |
| `.Aneka.BusinessTerritorialLimit` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.ConditionExtentions` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Contestant` | pxTextInput | 1 | ya | isGolfInsurance |
| `.Aneka.Coverage` | pxTextArea | 1 | ya | isCoverageOtherObjAneka |
| `.Aneka.Event` | pxTextArea | 1 | ya | isGolfInsurance |
| `.Aneka.Exclusion` | pxTextArea | 1 | ya | isExclusion |
| `.Aneka.InterestInsureds` | pxTextArea | 1 | ya | isInterestInsuredOtherObjAneka |
| `.Aneka.LimitOfGuarantee` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.LocationInterestInsured` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Obligee.AdministrationFee` | pxInteger | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.AgreementDate` | pxDateTime | 1 | ya | ALWAYS |
| `.Aneka.Obligee.AgreementNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.AuctionDate` | pxDateTime | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.AuctionNo` | pxTextInput | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.BankAddress` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Obligee.BankGuaranteeNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.BankName` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.ContractDate` | pxDateTime | 4 | ya | ALWAYS |
| `.Aneka.Obligee.ContractNo` | pxTextInput | 5 | ya | ALWAYS |
| `.Aneka.Obligee.CustomsActivity` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.DateOfGuarantee` | pxDateTime | 1 | ya | ALWAYS |
| `.Aneka.Obligee.FacilityDecreeDate` | pxDateTime | 1 | ya | ALWAYS |
| `.Aneka.Obligee.FacilityDecreeNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.FacilityStatus` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Obligee.InsuredNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.Jabatan` | pxTextInput | 7 | ya | ALWAYS |
| `.Aneka.Obligee.JenisPenjaminan` | pxDropdown | 1 | ya | ALWAYS |
| `.Aneka.Obligee.KBGDate` | pxDateTime | 1 | ya | ALWAYS |
| `.Aneka.Obligee.KBGNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.LengthPeriod` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.Location` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Obligee.MaxOfLimitOfLiability` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.NIPER` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.ObligeeAddress` | pxTextArea | 6 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ObligeeName` | pxTextInput | 7 | campur [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ObligeePhone` | pxTextInput | 5 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.OCDescription` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Obligee.OtherFacilityNo` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Obligee.Pasal` | pxDropdown | 1 | ya | ALWAYS |
| `.Aneka.Obligee.PaymentClaim` | pxTextInput | 6 | ya | ALWAYS |
| `.Aneka.Obligee.PIBDate` | pxDateTime | 1 | ya | ALWAYS |
| `.Aneka.Obligee.PIBNumber` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.PKDate` | pxDateTime | 1 | ya | ALWAYS |
| `.Aneka.Obligee.PKDate2` | pxDateTime | 1 | ya | ALWAYS |
| `.Aneka.Obligee.PKNo1` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.PKNo2` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.PlaceOfEvent` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Obligee.ProjectName` | pxTextArea | 6 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | ALWAYS |
| `.Aneka.Obligee.ProjectValue` | pxTextInput | 2 | ya | ALWAYS |
| `.Aneka.Obligee.PublishedIn` | pxTextArea/pxTextInput | 8 | ya | pyWorkPage.Quotation.BusinessCode!=35 &&  pyWorkPage.Quotation.Busines... |
| `.Aneka.Obligee.RegisterNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Aneka.Obligee.ScheduleType` | pxRadioButtons | 2 | ya | ALWAYS |
| `.Aneka.Obligee.UserAssignment` | pxAutoComplete | 8 | ya | pyWorkPage.Quotation.BusinessCode!=35 &&  pyWorkPage.Quotation.Busines... |
| `.Aneka.Occupation` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.OtherLOLDesc` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.PersonInsured` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.Remark` | pxTextArea | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac] | pyWorkPage.Quotation.BusinessCode == 18 \|\| pyWorkPage.Quotation.Busi... |
| `.Aneka.ScopeOfCover` | pxTextArea | 1 | ya | pyWorkPage.Quotation.BusinessCode == 17 \|\| pyWorkPage.Quotation.Busi... |
| `.Aneka.SpecialCondition` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.TerritorialLimits` | pxTextArea | 1 | ya | ALWAYS |
| `.Aneka.TradingWarranties` | pxTextArea | 1 | ya | ALWAYS |

#### ViewDtlOutgo

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 1 - grid: .OutgoList

_Tidak ada properti terikat._

#### ViewDtlParticipantPA

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 16 (unik 14) - read-only 16 - RepeatGrid 1 - grid: .ASMHeir

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMCCAmount` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMClass` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 2 | ya | ALWAYS |
| `.ASMGender` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMHeight` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMHeirPercentage` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMJobName` | pxAutoComplete | 1 | ya | ALWAYS |
| `.ASMLeftHanded` | pxRadioButtons | 1 | ya | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMRelation` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMWeight` | pxTextInput | 1 | ya | ALWAYS |
| `.pyFullName` | pxTextInput | 2 | ya | ALWAYS |
| `.Remark` | pxTextArea | 1 | ya | ALWAYS |

#### ViewGeneralPolisFacOut

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 15 (unik 15) - read-only 15 - RepeatGrid 1 - grid: .Policy.CIFData.AddressList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | (tanpa) | 1 | ya | ALWAYS |
| `.ASMAddressType` | pxDropdown | 1 | ya | ALWAYS |
| `.Policy.EndDateTime` | pxDateTime | 1 | ya | ALWAYS |
| `.Policy.Note` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.PolicyNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.QQName` | pxDropdown | 1 | ya | ALWAYS |
| `.Policy.RefNo` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.StartDateTime` | pxDateTime | 1 | ya | ALWAYS |
| `.Policy.SumOfTSI` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.TheInsured` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.TypeOfCoins` | pxDropdown | 1 | ya | ALWAYS |
| `.Quotation.BranchName` | pxTextInput | 1 | ya | ALWAYS |
| `.Quotation.BusinessName` | pxTextInput | 1 | ya | ALWAYS |
| `.Quotation.MarketingName` | pxTextInput | 1 | ya | ALWAYS |
| `InputData.CARI12` | pxTextInput | 1 | ya | ALWAYS |

#### ViewIndemnity_Section

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Search` - field terikat 3 (unik 3) - read-only 0 - RepeatGrid 1 - grid: ViewIndemnity.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CARI1` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI11` | (tanpa) | 1 | tidak | ALWAYS |
| `.CARI12` | (tanpa) | 1 | tidak | ALWAYS |

#### ViewIndividualInformationFacOut

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 20 (unik 20) - read-only 20 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CIFData.Customer_P.pyCity` | pxTextInput | 1 | ya | ALWAYS |
| `.CIFData.Customer_P.pyCompany` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMAnotherIncome` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMClientID` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMCompanyInfo.ASMFoundedIn` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMDateOfBirth` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMDependentCount` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMGender` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMIDCard` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMJobDesc` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMMaritalStatus` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMNationality` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMNickName` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMNPWP` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMPosition` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMSourceOfIncome` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.ASMYearlyIncome` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.pyFirstName` | pxTextInput | 1 | ya | ALWAYS |
| `.Customer_P.pyLastName` | pxTextInput | 1 | ya | ALWAYS |
| `.pyTitle` | pxTextInput | 1 | ya | ALWAYS |

#### ViewInwardFacultativeDtl_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-SpreadingRisk` - field terikat 69 (unik 30) - read-only 57 - RepeatGrid 12 - grid: .OfferFacIn.CedingCedantList, .OfferFacIn.CurrencyList, SpreadingList.pxResults, TotalSpreadAll.pxResults, TotalSpreadingCurrency.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.CedingCoName` | pxDisplayText | 1 | ya | ALWAYS |
| `.ClaimEstimation` | pxNumber | 4 | ya | ALWAYS |
| `.ClaimSpreaded` | pxNumber | 2 | ya | ALWAYS |
| `.Name` | (tanpa) | 6 | tidak | ALWAYS |
| `.OfferFacIn.CedingRetention` | pxNumber | 1 | ya | IsLife |
| `.OfferFacIn.PercentShare` | pxNumber | 2 | ya | ALWAYS |
| `.OfferFacIn.PPnCheck` | pxCheckbox | 1 | tidak | ALWAYS |
| `.OfferFacIn.QuotationData.SobName` | pxDisplayText | 1 | ya | ALWAYS |
| `.OfferFacIn.ShareCedantType` | pxRadioButtons | 1 | tidak | ALWAYS |
| `.OfferFacIn.TotalPremiNusaRe` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.OfferFacIn.TotalTSINusaRe` | pxTextInput | 1 | ya | 1 = 2 |
| `.OfferFacIn.TotalTSINusaReSpreading` | pxNumber/pxTextInput | 2 | ya | ALWAYS |
| `.OldTSI` | pxNumber | 4 | ya | ALWAYS |
| `.Policy.Payment.BrokerageFee` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Commision` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.PctBrokerageFee` | pxNumber | 1 | ya | !IsUW |
| `.Policy.Payment.PPh` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.PPN` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.RICommision` | pxNumber | 1 | ya | !IsUW |
| `.PPH_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.PPN_Note` | (tanpa) | 1 | ya | ALWAYS |
| `.Premium` | pxNumber | 4 | ya | ALWAYS |
| `.PremiumSpreaded` | pxNumber | 4 | ya | ALWAYS |
| `.ShareCeding` | pxNumber | 1 | ya | ALWAYS |
| `.SharePercentage` | pxNumber | 5 | ya | OutputParam.ERRMSG6=='' |
| `.TreatyName` | (tanpa) | 4 | tidak | ALWAYS |
| `.TreatyType` | pxDropdown | 5 | ya | ALWAYS |
| `.TSI` | pxNumber | 4 | ya | ALWAYS |
| `.TSIAddCap` | pxNumber | 2 | ya | ALWAYS |
| `.TSISpreaded` | pxNumber | 4 | ya | ALWAYS |

#### ViewObjectAnekaFacOut

rumpun **INTI** - kelas dominan `Data-Address` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 1 - grid: .LocationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMCity` | (tanpa) | 1 | ya | ALWAYS |
| `.SumOfTSI` | (tanpa) | 1 | ya | ALWAYS |

#### ViewObjectDtlAneka

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 1 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ObjectName` | pxTextInput | 1 | tidak [IsGroupFac \|\| IsGroupUWFac] | ALWAYS |

#### ViewObjectFireFacOut

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Property` - field terikat 5 (unik 5) - read-only 4 - RepeatGrid 1 - grid: .PropertyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.IsMaterialDamage` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 1 | ya | IsFac |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.pySelected` | pxCheckbox | 1 | tidak | ALWAYS |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |

#### ViewObjectMarineCargoFacOut

rumpun **CARGO** - kelas dominan `ASM-FW-GISFW-Data-Cargo` - field terikat 4 (unik 4) - read-only 4 - RepeatGrid 1 - grid: .CargoList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ConveyanceNote` | (tanpa) | 1 | ya | ALWAYS |
| `.GoodNote` | (tanpa) | 1 | ya | ALWAYS |
| `.PackingNote` | (tanpa) | 1 | ya | ALWAYS |
| `.TradingNote` | (tanpa) | 1 | ya | ALWAYS |

#### ViewObjectMarineHullFacOut

rumpun **INTI** - kelas dominan `Data-Address` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 1 - grid: .LocationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | pxTextInput | 1 | tidak | ALWAYS |

#### ViewObjectOccupation

rumpun **FIRE** - kelas dominan `ASM-FW-GISFW-Data-Occupation` - field terikat 1 (unik 1) - read-only 0 - RepeatGrid 1 - grid: .OccupationList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OccupationName` | pxTextInput | 1 | tidak | ALWAYS |

#### ViewObjectPAFacOut

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 5 (unik 5) - read-only 0 - RepeatGrid 1 - grid: .PersonList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMHeight` | pxInteger | 1 | tidak | ALWAYS |
| `.ASMJobName` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMWeight` | pxInteger | 1 | tidak | ALWAYS |
| `.pyFullName` | (tanpa) | 1 | tidak | ALWAYS |

#### ViewObjectSavior

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 0

_Tidak ada properti terikat._

#### ViewObjectSubContract

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-SubContract` - field terikat 3 (unik 3) - read-only 3 - RepeatGrid 2 - grid: .LocationList, .SubContractList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.SubContractID` | pxTextInput | 1 | ya | ALWAYS |
| `.SubContractName` | pxTextInput | 1 | ya | ALWAYS |

#### ViewObjectTravelFacOut

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 5 (unik 5) - read-only 0 - RepeatGrid 1 - grid: .PersonList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMDateOfBirth` | pxDateTime | 1 | tidak | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | tidak | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | tidak | ALWAYS |
| `.pyFullName` | (tanpa) | 1 | tidak | ALWAYS |
| `.pySelected` | pxCheckbox | 1 | tidak | ALWAYS |

#### ViewObjectVehicleFacOut

rumpun **MBU** - kelas dominan `ASM-FW-GISFW-Data-Vehicle` - field terikat 6 (unik 6) - read-only 0 - RepeatGrid 1 - grid: .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BrandName` | pxDropdown | 1 | tidak | ALWAYS |
| `.ChassisNumber` | pxTextInput | 1 | tidak | ALWAYS |
| `.EngineNumber` | (tanpa) | 1 | tidak | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | tidak | ALWAYS |
| `.ModelName` | pxTextInput | 1 | tidak | ALWAYS |
| `.TypeName` | (tanpa) | 1 | tidak | ALWAYS |

#### ViewOkupasiAneka

rumpun **ANEKA** - kelas dominan `ASM-FW-GISFW-Data-Aneka` - field terikat 16 (unik 10) - read-only 2 - RepeatGrid 2 - grid: .AnekaList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.LimitOfLiability` | pxTextInput | 2 | tidak | IsObjectSectionAneka |
| `.ObjectName` | pxTextInput | 2 | tidak | ALWAYS |
| `.OccupationId` | pxAutoComplete | 1 | ya | ALWAYS |
| `.OccupationName` | pxTextInput | 1 | ya | ALWAYS |
| `.Quantity` | pxTextInput | 1 | tidak | IsObjectWithQuantityYear \|\| IsHE |
| `.Section` | pxDropdown | 2 | tidak | IsObjectSectionAneka |
| `.UnitQuantity` | pxTextInput | 1 | tidak | IsObjectWithQuantityYear \|\| IsHE |
| `.VehicleHE.ChassisNumber` | pxTextInput | 2 | tidak | IsHE |
| `.VehicleHE.EngineNumber` | pxTextInput | 2 | tidak | IsHE |
| `.Year` | pxTextInput | 2 | tidak | IsObjectWithQuantityYear \|\| IsHE |

#### ViewOldDataPayment_IsUW

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-OfferFacIn-Currency` - field terikat 8 (unik 6) - read-only 8 - RepeatGrid 2 - grid: .OfferFacIn.CurrencyList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Name` | pxTextInput | 2 | ya | ALWAYS |
| `.OfferFacIn.PolicyData.Payment.Installment` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMOldPayment` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.EDMPremiMenjadi` | pxNumber | 1 | ya | ALWAYS |
| `.Policy.Payment.Premium` | pxNumber | 1 | ya | ALWAYS |
| `.SumTotalPayment` | pxNumber | 2 | ya | ALWAYS |

#### ViewOldDeductible

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Deductible` - field terikat 1 (unik 1) - read-only  - RepeatGrid 1 - grid: .DeductibleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Descriptions` | pxTextArea | 1 | ya | ALWAYS |

#### ViewOpenFollowingPolicy

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Work` - field terikat 14 (unik 13) - read-only 5 - RepeatGrid 2 - grid: Attachment.pxResults, tempWorkPageView.OfferFacIn.ViewSuggest

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Approval` | pxDropdown | 1 | ya | ALWAYS |
| `.CommentSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `.CountAttach` | pxTextInput | 1 | tidak | ALWAYS |
| `.DateSuggest` | pxDateTime | 1 | ya | ALWAYS |
| `.FilterFollowingNB` | pxTextInput | 2 | tidak | 1=2 |
| `.IsCedingConfirm` | pxDisplayText | 1 | ya | ALWAYS |
| `.NOTE` | pxTextInput | 1 | tidak | ALWAYS |
| `.PICSuggest` | pxDisplayText | 1 | ya | ALWAYS |
| `InputRiskAddress.CityName` | pxAutoComplete | 1 | tidak | InputRiskAddress.IDProvince!='' |
| `InputRiskAddress.DistrictName` | pxAutoComplete | 1 | tidak | InputRiskAddress.IDCity!='' |
| `InputRiskAddress.ProvinceName` | pxAutoComplete | 1 | tidak | ViewPolis.CARI40!='' |
| `pyPortal.FilterForViewRetro` | pxDropdown | 1 | tidak | ALWAYS |
| `TempPolis.CARI12` | pxDropdown | 1 | tidak | TempPolis.CARI26==1 |

#### ViewOutgo

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Data-Outgo` - field terikat 7 (unik 6) - read-only 7 - RepeatGrid 0

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.OutgoAmount` | pxNumber | 1 | ya | .TypeOfNominal='false' |
| `.PercentOutgo` | pxTextInput | 1 | ya | .TypeOfNominal='true' |
| `.ReceiverOutgo` | pxAutoComplete | 1 | ya | .TypeOfOutgo != '1' |
| `.SourceBizName` | pxAutoComplete | 2 | ya | .TypeOfOutgo == '1' ; .TypeOfOutgo == '4' |
| `.TypeOfNominal` | pxCheckbox | 1 | ya | ALWAYS |
| `.TypeOfOutgo` | pxDropdown | 1 | ya | ALWAYS |

#### ViewPerilsFacOut

rumpun **(tanpa field)** - kelas dominan `(kosong)` - field terikat 0 (unik 0) - read-only 0 - RepeatGrid 1 - grid: .VehicleList

_Tidak ada properti terikat._

#### ViewPolicyListRO

rumpun **PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 5 (unik 5) - read-only 4 - RepeatGrid 1 - grid: .PersonList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ASMDateOfBirth` | pxDateTime | 1 | tidak | ALWAYS |
| `.ASMGender` | pxDisplayText | 1 | ya | ALWAYS |
| `.ASMRegNo` | pxDisplayText | 1 | ya | ALWAYS |
| `.BookNo` | pxDisplayText | 1 | ya | ALWAYS |
| `.pyFullName` | pxDisplayText | 1 | ya | ALWAYS |

#### ViewTabGroupPolisFacOutSavior

rumpun **INTI+PA-LIFE** - kelas dominan `Data-Party-Person` - field terikat 32 (unik 31) - read-only 31 - RepeatGrid 10 - grid: .CargoList, .ClauseList, .CoinsList, .LocationList, .PersonList, .PropertyList, .VehicleList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.ArgumentCount` | (tanpa) | 1 | ya | ALWAYS |
| `.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMDateOfBirth` | pxDateTime | 1 | ya | ALWAYS |
| `.ASMHeight` | pxInteger | 1 | ya | ALWAYS |
| `.ASMIDCard` | pxTextInput | 1 | ya | ALWAYS |
| `.ASMJobName` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMLeftHanded` | pxDropdown | 1 | ya | ALWAYS |
| `.ASMParticipantStatus` | pxDropdown | 1 | tidak | ALWAYS |
| `.ASMWeight` | pxInteger | 1 | ya | ALWAYS |
| `.BrandName` | pxDropdown | 1 | ya | ALWAYS |
| `.ClauseCode` | (tanpa) | 1 | ya | ALWAYS |
| `.ClauseDescription` | (tanpa) | 1 | ya | ALWAYS |
| `.ClauseLanguage` | (tanpa) | 1 | ya | ALWAYS |
| `.ClauseTitle` | (tanpa) | 1 | ya | ALWAYS |
| `.ConveyanceNote` | pxTextInput | 1 | ya | ALWAYS |
| `.GoodNote` | pxTextInput | 1 | ya | ALWAYS |
| `.IsMaterialDamage` | pxCheckbox | 1 | ya | ALWAYS |
| `.IsTopRisk` | pxCheckbox | 1 | ya | IsFac |
| `.KodeSection` | (tanpa) | 1 | ya | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya | ALWAYS |
| `.ModelName` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectName` | pxTextInput | 1 | ya | ALWAYS |
| `.ObjectNo` | pxTextInput | 1 | ya | ALWAYS |
| `.PackingNote` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.STNC` | pxTextInput | 1 | ya | ALWAYS |
| `.Policy.STNCDate` | pxTextInput | 1 | ya | ALWAYS |
| `.pyFullName` | (tanpa) | 2 | ya | ALWAYS |
| `.RiskLocation.ASMAddress` | pxTextInput | 1 | ya | ALWAYS |
| `.TradingNote` | pxTextInput | 1 | ya | ALWAYS |
| `.TypeName` | pxTextInput | 1 | ya | ALWAYS |

#### ViewVehicleGridFacOut

rumpun **MBU** - kelas dominan `ASM-FW-GISFW-Data-Vehicle` - field terikat 17 (unik 15) - read-only 17 - RepeatGrid 1 - grid: .AccessoryList

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.BrandName` | pxAutoComplete | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac ...] | ALWAYS |
| `.CarParking` | pxDropdown | 1 | ya [.CarParking!=''\|\|IsSpreadingUW] | ALWAYS |
| `.CC` | pxTextInput | 1 | ya [.CC != '' \|\|IsSpreadingUW \|\| IsGroupFac \|\|...] | ALWAYS |
| `.ChassisNumber` | pxTextInput | 1 | ya [.ChassisNumber!=''\|\|IsSpreadingUW \|\| IsGroup...] | ALWAYS |
| `.ColorName` | pxAutoComplete | 1 | ya [.ColorName!=''\|\|IsSpreadingUW \|\| IsGroupFac ...] | ALWAYS |
| `.EngineNumber` | pxTextInput | 1 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac ...] | ALWAYS |
| `.LicensePlate` | pxTextInput | 1 | ya [.LicensePlate!=''\|\|IsSpreadingUW \|\| IsGroupF...] | ALWAYS |
| `.ManufactureYear` | pxTextInput | 1 | ya [.ManufactureYear!=''\|\|IsSpreadingUW \|\| IsGro...] | ALWAYS |
| `.ModelName` | pxAutoComplete | 1 | ya [.ModelName!= 0 \|\| IsSpreadingUW \|\| IsGroupFa...] | ALWAYS |
| `.Occupation.OccupationName` | pxAutoComplete | 1 | ya [.Occupation.OccupationName!=''\|\|IsSpreadingUW ...] | ALWAYS |
| `.Price` | pxNumber | 1 | ya | ALWAYS |
| `.TransmissionType` | pxDropdown | 1 | ya [.TransmissionType != '' \|\| IsSpreadingUW \|\| ...] | ALWAYS |
| `.Type` | pxTextInput | 1 | ya [.Type!= '' \|\| IsSpreadingUW \|\| IsGroupFac \|...] | 1=2 |
| `.TypeName` | pxAutoComplete | 2 | ya [IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac ...] | ALWAYS |
| `.ValueName` | pxAutoComplete | 1 | ya | ALWAYS |

#### WarrantyList

rumpun **INTI** - kelas dominan `ASM-FW-GISFW-Int-WARRANTY` - field terikat 2 (unik 2) - read-only 2 - RepeatGrid 1 - grid: pgRepPgSubSectionWarrantyListB.pxResults

| Properti | Kontrol | n | Read-only | Kondisi tampil |
| --- | --- | ---: | --- | --- |
| `.Description` | pxDisplayText | 1 | ya | ALWAYS |
| `.ID` | pxDisplayText | 1 | ya | ALWAYS |

---

## Lampiran B - Katalog 1.805 ekspresi properti unik

Diurutkan menurun menurut jumlah pengikatan. Kolom:
`n` = jumlah sel yang mengikatnya di seluruh NB ·
`Section` = berapa Section berbeda memakainya ·
`Kelas Pega dominan` = `pySmartPromptClass` terbanyak ·
`Kontrol` = seluruh `pyFormat` yang pernah dipakai (lebih dari satu = bukti tipe bertentangan, Sec 8.3) ·
`RO` = sel read-only / total sel ·
`Lingkup` = relatif (di bawah agregat case), agregat case eksplisit (`pyWorkPage.`), atau nama halaman di luar agregat.

Placeholder `.pyTemplate*` tidak termasuk.
| # | Ekspresi properti | n | Section | Kelas Pega dominan | Kontrol | RO | Lingkup |
| ---: | --- | ---: | ---: | --- | --- | ---: | --- |
| 1 | `.Premium` | 107 | 81 | ASM-FW-GISFW-Data-Coverage | pxDisplayText/pxNumber/pxTextInput | 98/107 | relatif (agregat case) |
| 2 | `.TSI` | 103 | 66 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxDisplayText/pxInteger/pxNumber/pxTextInput | 98/103 | relatif (agregat case) |
| 3 | `.Name` | 92 | 47 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxDisplayText/pxLink/pxTextInput | 57/92 | relatif (agregat case) |
| 4 | `.Currency` | 85 | 48 | ASM-FW-GISFW-Data-PropertyItem | pxAutoComplete/pxDisplayText/pxDropdown/pxNumber/pxTextInput | 54/85 | relatif (agregat case) |
| 5 | `.ObjectName` | 76 | 39 | ASM-FW-GISFW-Data-Aneka | pxAutoComplete/pxDisplayText/pxTextArea/pxTextInput | 63/76 | relatif (agregat case) |
| 6 | `.pyFullName` | 60 | 25 | Data-Party-Person | pxDisplayText/pxDropdown/pxTextInput | 39/60 | relatif (agregat case) |
| 7 | `.CoverageNote` | 57 | 45 | ASM-FW-GISFW-Data-Coverage | pxAutoComplete/pxDropdown/pxTextArea/pxTextInput | 43/57 | relatif (agregat case) |
| 8 | `.Rate` | 52 | 41 | ASM-FW-GISFW-Data-Coverage | pxDisplayText/pxNumber/pxTextInput | 47/52 | relatif (agregat case) |
| 9 | `.ASMDateOfBirth` | 51 | 24 | Data-Party-Person | pxDateTime/pxInteger | 31/51 | relatif (agregat case) |
| 10 | `.TreatyType` | 43 | 17 | ASM-FW-GISFW-Data-SpreadingRisk | pxDropdown/pxTextInput | 31/43 | relatif (agregat case) |
| 11 | `.ASMGender` | 42 | 21 | Data-Party-Person | pxDisplayText/pxDropdown | 24/42 | relatif (agregat case) |
| 12 | `.OccupationName` | 40 | 34 | ASM-FW-GISFW-Data-Occupation | pxAutoComplete/pxDisplayText/pxTextArea/pxTextInput | 28/40 | relatif (agregat case) |
| 13 | `.SharePercentage` | 36 | 17 | ASM-FW-GISFW-Data-SpreadingRisk | pxNumber | 28/36 | relatif (agregat case) |
| 14 | `.TSISpreaded` | 34 | 17 | ASM-FW-GISFW-Data-SpreadingRisk | pxNumber | 29/34 | relatif (agregat case) |
| 15 | `.Currency.Name` | 32 | 21 | ASM-FW-GISFW-Data-Aneka | pxAutoComplete/pxDropdown/pxTextInput | 23/32 | relatif (agregat case) |
| 16 | `.Value` | 32 | 8 | ASM-FW-GISFW-Data-TreatyInTotal | pxNumber/pxTextInput | 4/32 | relatif (agregat case) |
| 17 | `.PremiumSpreaded` | 30 | 14 | ASM-FW-GISFW-Data-SpreadingRisk | pxNumber | 30/30 | relatif (agregat case) |
| 18 | `.Property.RiskLocation.ASMAddress` | 30 | 21 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDisplayText/pxTextArea/pxTextInput | 30/30 | relatif (agregat case) |
| 19 | `.ItemType` | 29 | 18 | ASM-FW-GISFW-Data-PropertyItem | pxAutoComplete/pxDisplayText/pxNumber/pxTextInput | 29/29 | relatif (agregat case) |
| 20 | `.BrandName` | 25 | 19 | ASM-FW-GISFW-Data-Vehicle | pxAutoComplete/pxDisplayText/pxDropdown/pxTextInput | 19/25 | relatif (agregat case) |
| 21 | `.TSIObjectItem` | 25 | 15 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 25/25 | relatif (agregat case) |
| 22 | `.PremiNusantaraRe` | 24 | 20 | ASM-FW-GISFW-Data-Coverage | pxDisplayText/pxNumber | 24/24 | relatif (agregat case) |
| 23 | `.ASMLeftHanded` | 23 | 17 | Data-Party-Person | pxDropdown/pxRadioButtons | 10/23 | relatif (agregat case) |
| 24 | `.OccupationID` | 23 | 23 | ASM-FW-GISFW-Data-Occupation | pxAutoComplete/pxTextInput | 20/23 | relatif (agregat case) |
| 25 | `.Year` | 23 | 8 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 18/23 | relatif (agregat case) |
| 26 | `.ASMJobName` | 22 | 17 | Data-Party-Person | pxAutoComplete/pxDropdown | 14/22 | relatif (agregat case) |
| 27 | `.TypeName` | 22 | 18 | ASM-FW-GISFW-Data-Vehicle | pxAutoComplete/pxDisplayText/pxDropdown/pxInteger/pxTextInput | 19/22 | relatif (agregat case) |
| 28 | `.ID` | 20 | 20 | ASM-FW-GISFW-Int-SUBCONTRACT | pxAutoComplete/pxDisplayText/pxTextInput | 13/20 | relatif (agregat case) |
| 29 | `.TSINusantaraRe` | 20 | 20 | ASM-FW-GISFW-Data-Coverage | pxDisplayText/pxNumber | 20/20 | relatif (agregat case) |
| 30 | `.Descriptions` | 19 | 15 | ASM-FW-GISFW-Data-Deductible | pxDisplayText/pxDropdown/pxTextArea/pxTextInput | 16/19 | relatif (agregat case) |
| 31 | `.SumTotalPayment` | 19 | 12 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 19/19 | relatif (agregat case) |
| 32 | `.Aneka.Obligee.PublishedIn` | 18 | 4 | ASM-FW-GISFW-Work | pxTextArea/pxTextInput | 18/18 | relatif (agregat case) |
| 33 | `.Aneka.Obligee.UserAssignment` | 18 | 4 | ASM-FW-GISFW-Work | pxAutoComplete | 18/18 | relatif (agregat case) |
| 34 | `.ASMRelation` | 18 | 13 | Data-Party-Person | pxDropdown | 9/18 | relatif (agregat case) |
| 35 | `.ClaimSpreaded` | 18 | 11 | ASM-FW-GISFW-Data-SpreadingRisk | pxNumber | 18/18 | relatif (agregat case) |
| 36 | `.ConveyanceNote` | 18 | 16 | ASM-FW-GISFW-Data-Cargo | pxAutoComplete/pxDisplayText/pxDropdown/pxTextArea/pxTextInput | 16/18 | relatif (agregat case) |
| 37 | `.Property.ObjectNo` | 18 | 13 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDisplayText/pxNumber/pxTextInput | 18/18 | relatif (agregat case) |
| 38 | `.ASMIDCard` | 17 | 17 | Data-Party-Person | pxInteger/pxTextInput | 12/17 | relatif (agregat case) |
| 39 | `.LicensePlate` | 17 | 17 | ASM-FW-GISFW-Data-Vehicle | pxDisplayText/pxDropdown/pxTextInput | 13/17 | relatif (agregat case) |
| 40 | `.ModelName` | 17 | 17 | ASM-FW-GISFW-Data-Vehicle | pxAutoComplete/pxDisplayText/pxInteger/pxTextInput | 15/17 | relatif (agregat case) |
| 41 | `.TSILiability` | 17 | 13 | ASM-FW-GISFW-Data-Coverage | pxNumber/pxTextInput | 15/17 | relatif (agregat case) |
| 42 | `.Amount` | 16 | 16 | ASM-FW-GISFW-Data-Deductible | pxNumber | 12/16 | relatif (agregat case) |
| 43 | `.ASMHeight` | 16 | 16 | Data-Party-Person | pxInteger/pxNumber/pxTextInput | 10/16 | relatif (agregat case) |
| 44 | `.ASMWeight` | 16 | 16 | Data-Party-Person | pxInteger/pxNumber/pxTextInput | 10/16 | relatif (agregat case) |
| 45 | `.Condition` | 16 | 14 | ASM-FW-GISFW-Data-Deductible | pxDropdown | 13/16 | relatif (agregat case) |
| 46 | `.OfferFacIn.PercentShare` | 16 | 9 | ASM-FW-GISFW-Work | pxNumber/pxTextInput | 13/16 | relatif (agregat case) |
| 47 | `.ReinsurerName` | 16 | 10 | ASM-FW-GISFW-Data-FacOffer | pxAutoComplete/pxDisplayText/pxTextInput | 14/16 | relatif (agregat case) |
| 48 | `.TradingNote` | 16 | 15 | ASM-FW-GISFW-Data-Cargo | pxAutoComplete/pxDisplayText/pxTextInput | 14/16 | relatif (agregat case) |
| 49 | `.TSIOld` | 16 | 15 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber/pxTextInput | 15/16 | relatif (agregat case) |
| 50 | `.Unit` | 16 | 13 | ASM-FW-GISFW-Data-PropertyItem | pxDropdown/pxInteger/pxNumber | 15/16 | relatif (agregat case) |
| 51 | `.ASMParticipantStatus` | 15 | 15 | Data-Party-Person | pxDropdown/pxInteger | 8/15 | relatif (agregat case) |
| 52 | `.ClauseDescription` | 15 | 7 | ASM-FW-GISFW-Data-Clause | pxAutoComplete/pxTextArea/pxTextInput | 5/15 | relatif (agregat case) |
| 53 | `.PackingNote` | 15 | 15 | ASM-FW-GISFW-Data-Cargo | pxAutoComplete/pxDisplayText/pxInteger/pxTextInput | 13/15 | relatif (agregat case) |
| 54 | `.pxListSubscript` | 15 | 11 | ASM-FW-GISFW-Data-Cargo | pxDisplayText/pxNumber/pxTextInput | 9/15 | relatif (agregat case) |
| 55 | `.Quantity` | 15 | 9 | ASM-FW-GISFW-Data-Aneka | pxNumber/pxTextInput | 13/15 | relatif (agregat case) |
| 56 | `.Section` | 15 | 7 | ASM-FW-GISFW-Data-Aneka | pxDropdown | 11/15 | relatif (agregat case) |
| 57 | `InputParam.HASIL12` | 15 | 14 | ASM-FW-GISFW-Data-PropertyItem | pxDisplayText | 15/15 | LUAR agregat: InputParam |
| 58 | `.Aneka.Obligee.Jabatan` | 14 | 2 | ASM-FW-GISFW-Work | pxTextInput | 14/14 | relatif (agregat case) |
| 59 | `.Aneka.Obligee.ObligeeName` | 14 | 2 | ASM-FW-GISFW-Work | pxTextInput | 13/14 | relatif (agregat case) |
| 60 | `.ClaimEstimation` | 14 | 8 | ASM-FW-GISFW-Data-SpreadingRisk | pxNumber | 14/14 | relatif (agregat case) |
| 61 | `.Coverage` | 14 | 12 | ASM-FW-GISFW-Data-Coverage | pxAutoComplete/pxDropdown/pxTextInput | 10/14 | relatif (agregat case) |
| 62 | `.Discount` | 14 | 13 | ASM-FW-GISFW-Data-Coverage | pxDropdown/pxNumber/pxTextInput | 14/14 | relatif (agregat case) |
| 63 | `.GoodNote` | 14 | 14 | ASM-FW-GISFW-Data-Cargo | pxAutoComplete/pxDisplayText/pxInteger/pxTextInput | 12/14 | relatif (agregat case) |
| 64 | `.LimitofLiability` | 14 | 9 | ASM-FW-GISFW-Data-Aneka | pxNumber/pxTextInput | 10/14 | relatif (agregat case) |
| 65 | `.OfferFacIn.Obligee.Jabatan` | 14 | 2 | ASM-FW-GISFW-Work | pxTextInput | 14/14 | relatif (agregat case) |
| 66 | `.OfferFacIn.Obligee.ObligeeName` | 14 | 2 | ASM-FW-GISFW-Work | pxTextInput | 14/14 | relatif (agregat case) |
| 67 | `.OfferFacIn.Obligee.PublishedIn` | 14 | 2 | ASM-FW-GISFW-Work | pxTextArea | 14/14 | relatif (agregat case) |
| 68 | `.OfferFacIn.Obligee.UserAssignment` | 14 | 2 | ASM-FW-GISFW-Work | pxAutoComplete | 14/14 | relatif (agregat case) |
| 69 | `.OldTSI` | 14 | 6 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 14/14 | relatif (agregat case) |
| 70 | `.Policy.Payment.Commision` | 14 | 12 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 11/14 | relatif (agregat case) |
| 71 | `.Policy.Payment.EDMOldPayment` | 14 | 12 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 14/14 | relatif (agregat case) |
| 72 | `.Policy.Payment.EDMPremiMenjadi` | 14 | 12 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 14/14 | relatif (agregat case) |
| 73 | `.Policy.Payment.Premium` | 14 | 12 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 13/14 | relatif (agregat case) |
| 74 | `.Remark` | 14 | 14 | Data-Party-Person | pxDisplayText/pxTextArea | 9/14 | relatif (agregat case) |
| 75 | `.TreatyName` | 14 | 6 | ASM-FW-GISFW-Data-SpreadingRisk | pxTextInput | /14 | relatif (agregat case) |
| 76 | `.TSIAddCap` | 14 | 12 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 14/14 | relatif (agregat case) |
| 77 | `.CedingCoName` | 13 | 12 | ASM-FW-GISFW-Data-Quotation | pxDisplayText/pxTextInput | 13/13 | relatif (agregat case) |
| 78 | `.NOTE` | 13 | 9 | ASM-FW-GISFW-Data-LimitSummaryList | pxDisplayText/pxTextInput | 2/13 | relatif (agregat case) |
| 79 | `.OfferFacIn.TotalPremiNusaRe` | 13 | 6 | ASM-FW-GISFW-Work | pxNumber/pxTextInput | 13/13 | relatif (agregat case) |
| 80 | `.Property.ObjectName` | 13 | 13 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDisplayText/pxTextInput | 13/13 | relatif (agregat case) |
| 81 | `.Aneka.Obligee.ObligeeAddress` | 12 | 2 | ASM-FW-GISFW-Work | pxTextArea | 12/12 | relatif (agregat case) |
| 82 | `.Aneka.Obligee.PaymentClaim` | 12 | 2 | ASM-FW-GISFW-Work | pxTextInput | 12/12 | relatif (agregat case) |
| 83 | `.Aneka.Obligee.ProjectName` | 12 | 2 | ASM-FW-GISFW-Work | pxTextArea | 12/12 | relatif (agregat case) |
| 84 | `.IsTopRisk` | 12 | 10 | ASM-FW-GISFW-Data-Property | pxCheckbox | 7/12 | relatif (agregat case) |
| 85 | `.OfferFacIn.Obligee.ObligeeAddress` | 12 | 2 | ASM-FW-GISFW-Work | pxTextArea | 12/12 | relatif (agregat case) |
| 86 | `.OfferFacIn.Obligee.PaymentClaim` | 12 | 2 | ASM-FW-GISFW-Work | pxTextInput | 12/12 | relatif (agregat case) |
| 87 | `.OfferFacIn.Obligee.ProjectName` | 12 | 2 | ASM-FW-GISFW-Work | pxTextArea | 12/12 | relatif (agregat case) |
| 88 | `.OLDID` | 12 | 12 | ASM-FW-GISFW-Data-Coverage | pxAutoComplete/pxTextInput | 12/12 | relatif (agregat case) |
| 89 | `.PremiumOld` | 12 | 12 | ASM-FW-GISFW-Data-Coverage | pxNumber | 11/12 | relatif (agregat case) |
| 90 | `.RiskLocation.ASMAddress` | 12 | 12 | ASM-FW-GISFW-Data-Property | pxTextArea/pxTextInput | 10/12 | relatif (agregat case) |
| 91 | `.TimeExcess` | 12 | 10 | ASM-FW-GISFW-Data-Deductible | pxNumber/pxTextInput | 9/12 | relatif (agregat case) |
| 92 | `.TotalGrossPremi` | 12 | 5 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 12/12 | relatif (agregat case) |
| 93 | `.ChassisNumber` | 11 | 11 | ASM-FW-GISFW-Data-Vehicle | pxDisplayText/pxTextInput | 7/11 | relatif (agregat case) |
| 94 | `.EngineNumber` | 11 | 11 | ASM-FW-GISFW-Data-Vehicle | pxDisplayText/pxTextInput | 7/11 | relatif (agregat case) |
| 95 | `.OfferFacIn.TotalTSINusaReSpreading` | 11 | 6 | ASM-FW-GISFW-Work | pxNumber/pxTextInput | 11/11 | relatif (agregat case) |
| 96 | `.PropertyYear` | 11 | 8 | ASM-FW-GISFW-Data-PropertyItem | pxInteger/pxNumber | 11/11 | relatif (agregat case) |
| 97 | `.VehicleHE.ChassisNumber` | 11 | 6 | ASM-FW-GISFW-Data-Aneka | pxNumber/pxTextInput | 7/11 | relatif (agregat case) |
| 98 | `.VehicleHE.EngineNumber` | 11 | 6 | ASM-FW-GISFW-Data-Aneka | pxNumber/pxTextInput | 7/11 | relatif (agregat case) |
| 99 | `.Aneka.Obligee.ContractNo` | 10 | 2 | ASM-FW-GISFW-Work | pxTextInput | 10/10 | relatif (agregat case) |
| 100 | `.Aneka.Obligee.ObligeePhone` | 10 | 2 | ASM-FW-GISFW-Work | pxTextInput | 10/10 | relatif (agregat case) |
| 101 | `.ASMAddress` | 10 | 10 | Data-Address | pxTextArea/pxTextInput | 9/10 | relatif (agregat case) |
| 102 | `.ASMHeirPercentage` | 10 | 10 | Data-Party-Person | pxInteger | 9/10 | relatif (agregat case) |
| 103 | `.DateOfLoss` | 10 | 9 | ASM-FW-GISFW-Data-CauseOfLoss | pxDateTime | 9/10 | relatif (agregat case) |
| 104 | `.DiscountPercentage` | 10 | 10 | ASM-FW-GISFW-Data-Coverage | pxNumber/pxTextInput | 10/10 | relatif (agregat case) |
| 105 | `.DueDate` | 10 | 10 | ASM-FW-GISFW-Data-Installment | pxDateTime | 8/10 | relatif (agregat case) |
| 106 | `.ManufactureYear` | 10 | 10 | ASM-FW-GISFW-Data-Vehicle | pxTextInput | 10/10 | relatif (agregat case) |
| 107 | `.OfferFacIn.Obligee.ContractNo` | 10 | 2 | ASM-FW-GISFW-Work | pxTextInput | 10/10 | relatif (agregat case) |
| 108 | `.OfferFacIn.Obligee.ObligeePhone` | 10 | 2 | ASM-FW-GISFW-Work | pxTextInput | 10/10 | relatif (agregat case) |
| 109 | `.PctDeductible2` | 10 | 10 | ASM-FW-GISFW-Data-Deductible | pxNumber | 7/10 | relatif (agregat case) |
| 110 | `.Remarks` | 10 | 10 | ASM-FW-GISFW-Data-Quotation | pxDropdown/pxTextArea/pxTextInput | 3/10 | relatif (agregat case) |
| 111 | `.TypeDeductible2` | 10 | 10 | ASM-FW-GISFW-Data-Deductible | pxDropdown | 7/10 | relatif (agregat case) |
| 112 | `.ASMMaritalStatus` | 9 | 4 | Data-Party-Person | pxDropdown | 0/9 | relatif (agregat case) |
| 113 | `.CARI1` | 9 | 5 | ASM-FW-GISFW-Data-Search | pxDateTime/pxDisplayText | 8/9 | relatif (agregat case) |
| 114 | `.ClauseCode` | 9 | 7 | ASM-FW-GISFW-Data-Clause | pxDisplayText/pxTextInput | 9/9 | relatif (agregat case) |
| 115 | `.InputCondition` | 9 | 9 | ASM-FW-GISFW-Data-Deductible | pxTextInput | 6/9 | relatif (agregat case) |
| 116 | `.Limit` | 9 | 6 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 6/9 | relatif (agregat case) |
| 117 | `.MinMax` | 9 | 9 | ASM-FW-GISFW-Data-Deductible | pxDropdown | 9/9 | relatif (agregat case) |
| 118 | `.PctDeductible` | 9 | 9 | ASM-FW-GISFW-Data-Deductible | pxNumber | 6/9 | relatif (agregat case) |
| 119 | `.Policy.Payment.NetPremium` | 9 | 7 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 8/9 | relatif (agregat case) |
| 120 | `.TypeDeductible` | 9 | 9 | ASM-FW-GISFW-Data-Deductible | pxDropdown | 6/9 | relatif (agregat case) |
| 121 | `pyWorkPage.PolicyTreatyIn.TotalClaim` | 9 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 9/9 | agregat case (eksplisit) |
| 122 | `pyWorkPage.PolicyTreatyIn.TotalPremium` | 9 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 9/9 | agregat case (eksplisit) |
| 123 | `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | 9 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 9/9 | agregat case (eksplisit) |
| 124 | `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | 9 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 9/9 | agregat case (eksplisit) |
| 125 | `.Age` | 8 | 8 | Data-Party-Person | pxInteger | 5/8 | relatif (agregat case) |
| 126 | `.Aneka.Obligee.ContractDate` | 8 | 2 | ASM-FW-GISFW-Work | pxDateTime | 8/8 | relatif (agregat case) |
| 127 | `.Approval` | 8 | 8 | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | pxDropdown | 8/8 | relatif (agregat case) |
| 128 | `.Attention` | 8 | 6 | ASM-FW-GISFW-Data-FacOffer | pxTextArea | 6/8 | relatif (agregat case) |
| 129 | `.BalanceDueTo` | 8 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 8/8 | relatif (agregat case) |
| 130 | `.Claim` | 8 | 8 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency/pxNumber/pxTextInput | 6/8 | relatif (agregat case) |
| 131 | `.ClauseContent` | 8 | 8 | ASM-FW-GISFW-Data-Clause | pxAutoComplete/pxDisplayText/pxRichTextEditor | 6/8 | relatif (agregat case) |
| 132 | `.CommentSuggest` | 8 | 8 | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | pxDisplayText | 8/8 | relatif (agregat case) |
| 133 | `.Deductible` | 8 | 4 | ASM-FW-GISFW-Data-LimitSummaryList | pxDisplayText/pxNumber | 5/8 | relatif (agregat case) |
| 134 | `.FromRute` | 8 | 8 | ASM-FW-GISFW-Data-Trading | pxDisplayText/pxTextInput | 8/8 | relatif (agregat case) |
| 135 | `.InstallmentNo` | 8 | 8 | ASM-FW-GISFW-Data-Installment | pxInteger | 8/8 | relatif (agregat case) |
| 136 | `.IsMaterialDamage` | 8 | 8 | ASM-FW-GISFW-Data-Property | pxCheckbox | 6/8 | relatif (agregat case) |
| 137 | `.LossObject` | 8 | 8 | ASM-FW-GISFW-Data-CauseOfLoss | pxTextInput | 7/8 | relatif (agregat case) |
| 138 | `.ObjectNo` | 8 | 8 | ASM-FW-GISFW-Data-Property | pxTextInput | 8/8 | relatif (agregat case) |
| 139 | `.OfferFacIn.Obligee.ContractDate` | 8 | 2 | ASM-FW-GISFW-Work | pxDateTime | 8/8 | relatif (agregat case) |
| 140 | `.PaymentTotal` | 8 | 8 | ASM-FW-GISFW-Data-Installment | pxNumber | 8/8 | relatif (agregat case) |
| 141 | `.PICSuggest` | 8 | 8 | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | pxDisplayText | 8/8 | relatif (agregat case) |
| 142 | `.ShareCeding` | 8 | 7 | ASM-FW-GISFW-Data-Quotation | pxNumber | 6/8 | relatif (agregat case) |
| 143 | `.ToRute` | 8 | 8 | ASM-FW-GISFW-Data-Trading | pxDisplayText/pxTextInput | 8/8 | relatif (agregat case) |
| 144 | `.TotalPremiumNusantaraRe` | 8 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 8/8 | relatif (agregat case) |
| 145 | `.TreatyYear` | 8 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 8/8 | relatif (agregat case) |
| 146 | `.TSIObjectItemOld` | 8 | 7 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 8/8 | relatif (agregat case) |
| 147 | `.ASMCoverage(1).Loading` | 7 | 3 | Data-Party-Person | (tanpa) | 7/7 | relatif (agregat case) |
| 148 | `.ASMCoverage(1).PremiRp` | 7 | 3 | Data-Party-Person | (tanpa) | 0/7 | relatif (agregat case) |
| 149 | `.ASMCoverage(1).Premium` | 7 | 3 | Data-Party-Person | (tanpa) | 0/7 | relatif (agregat case) |
| 150 | `.ASMStatusCoverage` | 7 | 3 | Data-Party-Person | pxDisplayText | 7/7 | relatif (agregat case) |
| 151 | `.CARI11` | 7 | 5 | ASM-FW-GISFW-Data-Search | pxDisplayText/pxNumber | 5/7 | relatif (agregat case) |
| 152 | `.Deductible2` | 7 | 3 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 4/7 | relatif (agregat case) |
| 153 | `.EndPeriod` | 7 | 4 | ASM-FW-GISFW-Data-FacOffer | pxDateTime | 6/7 | relatif (agregat case) |
| 154 | `.InstallmentPercentage` | 7 | 7 | ASM-FW-GISFW-Data-Installment | pxNumber | 6/7 | relatif (agregat case) |
| 155 | `.IsCedingConfirm` | 7 | 7 | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | pxDisplayText | 7/7 | relatif (agregat case) |
| 156 | `.Layer` | 7 | 7 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDisplayText/pxNumber/pxTextInput | 7/7 | relatif (agregat case) |
| 157 | `.Loading` | 7 | 6 | ASM-FW-GISFW-Data-Coverage | pxNumber/pxTextInput | 5/7 | relatif (agregat case) |
| 158 | `.MDP` | 7 | 3 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 4/7 | relatif (agregat case) |
| 159 | `.MDP2` | 7 | 3 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 4/7 | relatif (agregat case) |
| 160 | `.MinPremium` | 7 | 6 | ASM-FW-GISFW-Data-Coverage | pxNumber/pxTextInput | 7/7 | relatif (agregat case) |
| 161 | `.OfferFacIn.CedingRetention` | 7 | 5 | ASM-FW-GISFW-Work | pxNumber | 4/7 | relatif (agregat case) |
| 162 | `.OurRef` | 7 | 4 | ASM-FW-GISFW-Data-FacOffer | pxTextInput | 7/7 | relatif (agregat case) |
| 163 | `.PremiRp` | 7 | 5 | ASM-FW-GISFW-Data-Coverage | pxNumber | 6/7 | relatif (agregat case) |
| 164 | `.Property.IsTopRisk` | 7 | 7 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxCheckbox | 5/7 | relatif (agregat case) |
| 165 | `.QuotationData.IsSurveyReport` | 7 | 7 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxRadioButtons | 3/7 | relatif (agregat case) |
| 166 | `.QuotationData.NoOfferSlip` | 7 | 7 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextArea/pxTextInput | /7 | relatif (agregat case) |
| 167 | `.ShareOffered` | 7 | 7 | ASM-FW-GISFW-Data-Occupation | pxNumber/pxTextInput | 7/7 | relatif (agregat case) |
| 168 | `.StartPeriod` | 7 | 4 | ASM-FW-GISFW-Data-FacOffer | pxDateTime | 6/7 | relatif (agregat case) |
| 169 | `.SumTotalTSI` | 7 | 7 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber/pxTextInput | 7/7 | relatif (agregat case) |
| 170 | `.Type` | 7 | 7 | ASM-FW-GISFW-Data-Vehicle | pxTextInput | 7/7 | relatif (agregat case) |
| 171 | `.BookNo` | 6 | 4 | ASM-FW-GISFW-Data-Policy | pxDisplayText/pxLink | /6 | relatif (agregat case) |
| 172 | `.CARI12` | 6 | 4 | ASM-FW-GISFW-Data-Search | pxNumber | 5/6 | relatif (agregat case) |
| 173 | `.ClauseParam` | 6 | 4 | ASM-FW-GISFW-Data-Clause | pxDisplayText/pxRadioButtons/pxTextInput | 2/6 | relatif (agregat case) |
| 174 | `.ClauseTitle` | 6 | 6 | ASM-FW-GISFW-Data-Clause | pxTextArea/pxTextInput | 6/6 | relatif (agregat case) |
| 175 | `.CoinsData.CoinsName` | 6 | 6 | ASM-FW-GISFW-Data-CauseOfLoss | pxAutoComplete/pxTextInput | 5/6 | relatif (agregat case) |
| 176 | `.DateTransfer` | 6 | 6 | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | pxDateTime | 6/6 | relatif (agregat case) |
| 177 | `.Detail` | 6 | 6 | ASM-FW-GISFW-Data-CauseOfLoss | pxTextArea/pxTextInput | 5/6 | relatif (agregat case) |
| 178 | `.EndDate` | 6 | 6 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDateTime/pxTextInput | 4/6 | relatif (agregat case) |
| 179 | `.Indemnity` | 6 | 6 | ASM-FW-GISFW-Data-Coverage | pxDisplayText/pxNumber | 5/6 | relatif (agregat case) |
| 180 | `.MaxPctTreatyCapacity` | 6 | 6 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 6/6 | relatif (agregat case) |
| 181 | `.MaxTreatyCapacity` | 6 | 6 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 6/6 | relatif (agregat case) |
| 182 | `.Occupation.OccupationName` | 6 | 6 | ASM-FW-GISFW-Data-Vehicle | pxAutoComplete/pxTextInput | 5/6 | relatif (agregat case) |
| 183 | `.OfferFacIn.QuotationData.SobName` | 6 | 6 | ASM-FW-GISFW-Work | pxDisplayText | 6/6 | relatif (agregat case) |
| 184 | `.OfferFacIn.TotalTSINusaRe` | 6 | 4 | ASM-FW-GISFW-Work | pxNumber/pxTextInput | 4/6 | relatif (agregat case) |
| 185 | `.PercentOffered` | 6 | 6 | ASM-FW-GISFW-Data-Occupation | pxNumber | 6/6 | relatif (agregat case) |
| 186 | `.PremiumRetro` | 6 | 4 | ASM-FW-GISFW-Data-Coverage | pxNumber | 2/6 | relatif (agregat case) |
| 187 | `.Property.RiskLocation.ASMZipCode` | 6 | 5 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 5/6 | relatif (agregat case) |
| 188 | `.ShareInTSI` | 6 | 6 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 6/6 | relatif (agregat case) |
| 189 | `.TableOfLimit.Description` | 6 | 6 | ASM-FW-GISFW-Data-Occupation | pxAutoComplete | 6/6 | relatif (agregat case) |
| 190 | `.TotalGrossPremiOld` | 6 | 5 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 6/6 | relatif (agregat case) |
| 191 | `.TreatyCapacity` | 6 | 6 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 6/6 | relatif (agregat case) |
| 192 | `.UnitQuantity` | 6 | 6 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 4/6 | relatif (agregat case) |
| 193 | `.ArgumentCount` | 5 | 5 | ASM-FW-GISFW-Data-Clause | pxTextInput | 5/5 | relatif (agregat case) |
| 194 | `.CARI10` | 5 | 3 | ASM-FW-GISFW-Data-Search | pxDisplayText | 4/5 | relatif (agregat case) |
| 195 | `.CARI13` | 5 | 3 | ASM-FW-GISFW-Data-Search | pxDisplayText/pxNumber | 5/5 | relatif (agregat case) |
| 196 | `.CARI14` | 5 | 3 | ASM-FW-GISFW-Data-Search | pxDisplayText/pxNumber | 5/5 | relatif (agregat case) |
| 197 | `.ClaimPercentage` | 5 | 5 | ASM-FW-GISFW-Data-SpreadingRisk | pxNumber | 3/5 | relatif (agregat case) |
| 198 | `.ConveyanceID` | 5 | 5 | ASM-FW-GISFW-Data-Conveyance | pxAutoComplete/pxDropdown | 3/5 | relatif (agregat case) |
| 199 | `.CoverageList(1).FacOutObjectList(1).Rate` | 5 | 3 | ASM-FW-GISFW-Data-Cargo | pxNumber | 5/5 | relatif (agregat case) |
| 200 | `.CoverageList(1).FacOutObjectList(1).ShareOffered` | 5 | 5 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 5/5 | relatif (agregat case) |
| 201 | `.Day` | 5 | 5 | ASM-FW-GISFW-Data-Coverage | pxRadioButtons | 3/5 | relatif (agregat case) |
| 202 | `.Deduction2` | 5 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency/pxNumber | 3/5 | relatif (agregat case) |
| 203 | `.EmployeeName` | 5 | 3 | ASM-FW-GISFW-Data-BatchMember | pxDisplayText | 5/5 | relatif (agregat case) |
| 204 | `.GoodID` | 5 | 5 | ASM-FW-GISFW-Data-Good | pxAutoComplete/pxDisplayText/pxDropdown | 5/5 | relatif (agregat case) |
| 205 | `.GrossPremium` | 5 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 3/5 | relatif (agregat case) |
| 206 | `.IndemnityPercentage` | 5 | 5 | ASM-FW-GISFW-Data-Coverage | pxNumber | 5/5 | relatif (agregat case) |
| 207 | `.IndexCoverage` | 5 | 3 | ASM-FW-GISFW-Data-Coverage | (tanpa) | 0/5 | relatif (agregat case) |
| 208 | `.Installment` | 5 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 3/5 | relatif (agregat case) |
| 209 | `.InsuredName` | 5 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 5/5 | relatif (agregat case) |
| 210 | `.LayerNo` | 5 | 5 | ASM-FW-GISFW-Data-Coverage | (tanpa) | /5 | relatif (agregat case) |
| 211 | `.Limit2` | 5 | 2 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 4/5 | relatif (agregat case) |
| 212 | `.LostLimit` | 5 | 5 | ASM-FW-GISFW-Data-Coverage | pxNumber/pxTextInput | 4/5 | relatif (agregat case) |
| 213 | `.OfferFacIn.PPnCheck` | 5 | 5 | ASM-FW-GISFW-Work | pxCheckbox | 0/5 | relatif (agregat case) |
| 214 | `.OfferFacIn.ShareCedantType` | 5 | 5 | ASM-FW-GISFW-Work | pxRadioButtons | 2/5 | relatif (agregat case) |
| 215 | `.PlanIP` | 5 | 3 | ASM-FW-GISFW-Data-BatchMember | pxDisplayText | 5/5 | relatif (agregat case) |
| 216 | `.PlanOP` | 5 | 3 | ASM-FW-GISFW-Data-BatchMember | pxDisplayText | 5/5 | relatif (agregat case) |
| 217 | `.Policy.Payment.BrokerageFee` | 5 | 5 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 3/5 | relatif (agregat case) |
| 218 | `.Policy.Payment.PctBrokerageFee` | 5 | 5 | ASM-FW-GISFW-Work | pxNumber | 3/5 | relatif (agregat case) |
| 219 | `.Policy.Payment.PPh` | 5 | 5 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 5/5 | relatif (agregat case) |
| 220 | `.Policy.Payment.PPN` | 5 | 5 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 5/5 | relatif (agregat case) |
| 221 | `.Policy.Payment.RICommision` | 5 | 5 | ASM-FW-GISFW-Work | pxNumber | 3/5 | relatif (agregat case) |
| 222 | `.PolicyData.SailDate` | 5 | 5 | ASM-FW-GISFW-Data-Cargo | pxDateTime | 5/5 | relatif (agregat case) |
| 223 | `.PolicyData.Ship.NM_SHIP` | 5 | 5 | ASM-FW-GISFW-Data-Cargo | pxDisplayText/pxDropdown | 5/5 | relatif (agregat case) |
| 224 | `.PPH_Note` | 5 | 5 | ASM-FW-GISFW-Work | (tanpa) | 5/5 | relatif (agregat case) |
| 225 | `.PPN_Note` | 5 | 5 | ASM-FW-GISFW-Work | (tanpa) | 5/5 | relatif (agregat case) |
| 226 | `.Property.Country` | 5 | 5 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 4/5 | relatif (agregat case) |
| 227 | `.Property.Province` | 5 | 5 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 4/5 | relatif (agregat case) |
| 228 | `.Property.RiskLocation.ASMCity` | 5 | 5 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 4/5 | relatif (agregat case) |
| 229 | `.Property.RiskLocation.ASMDistrict` | 5 | 5 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 4/5 | relatif (agregat case) |
| 230 | `.Property.RiskLocation.ASMRW` | 5 | 5 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 4/5 | relatif (agregat case) |
| 231 | `.pxCreateDateTime` | 5 | 3 | ASM-FW-GISFW-Data-Policy | pxDisplayText | 5/5 | relatif (agregat case) |
| 232 | `.pxCreateOpName` | 5 | 3 | ASM-FW-GISFW-Data-Policy | pxDisplayText | 5/5 | relatif (agregat case) |
| 233 | `.pyLabel` | 5 | 3 | ASM-FW-GISFW-Data-Policy | pxDisplayText | 5/5 | relatif (agregat case) |
| 234 | `.QuotationData.CedingCoName` | 5 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText | 5/5 | relatif (agregat case) |
| 235 | `.SPPANo` | 5 | 3 | ASM-FW-GISFW-Data-Policy | pxLink | 0/5 | relatif (agregat case) |
| 236 | `.SubContractID` | 5 | 5 | ASM-FW-GISFW-Data-SubContract | pxTextInput | 5/5 | relatif (agregat case) |
| 237 | `.SubContractName` | 5 | 5 | ASM-FW-GISFW-Data-SubContract | pxTextInput | 5/5 | relatif (agregat case) |
| 238 | `.TotalClaim` | 5 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 5/5 | relatif (agregat case) |
| 239 | `.TotalPremium` | 5 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 5/5 | relatif (agregat case) |
| 240 | `.TotalSharePercentageClaim` | 5 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 5/5 | relatif (agregat case) |
| 241 | `.TotalSharePercentagePremium` | 5 | 5 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 5/5 | relatif (agregat case) |
| 242 | `.TSINusaRe` | 5 | 2 | Data-Party-Person | (tanpa) | 0/5 | relatif (agregat case) |
| 243 | `.Zone` | 5 | 3 | ASM-FW-GISFW-Data-Coverage | pxAutoComplete/pxDisplayText/pxTextArea | 5/5 | relatif (agregat case) |
| 244 | `InputData.CARI11` | 5 | 3 | ASM-FW-GISFW-Work | pxNumber | 5/5 | LUAR agregat: InputData |
| 245 | `pyWorkPage.OfferFacIn.ProRatePercent` | 5 | 5 | ASM-FW-GISFW-Data-Coverage | pxNumber | 4/5 | agregat case (eksplisit) |
| 246 | `SearchSOB.CARI1` | 5 | 3 | ASM-FW-GISFW-Data-Quotation | pxAutoComplete | 0/5 | LUAR agregat: SearchSOB |
| 247 | `.AdditionalCapital` | 4 | 4 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | 2/4 | relatif (agregat case) |
| 248 | `.Aneka.Business` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 249 | `.Aneka.BusinessGolf` | 4 | 4 | ASM-FW-GISFW-Work | pxDropdown | 4/4 | relatif (agregat case) |
| 250 | `.Aneka.BusinessTerritorialLimit` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 251 | `.Aneka.ConditionExtentions` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 252 | `.Aneka.Contestant` | 4 | 4 | ASM-FW-GISFW-Work | pxTextInput | 4/4 | relatif (agregat case) |
| 253 | `.Aneka.Coverage` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 254 | `.Aneka.Event` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 255 | `.Aneka.Exclusion` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 256 | `.Aneka.InterestInsureds` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 257 | `.Aneka.LimitOfGuarantee` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 258 | `.Aneka.LocationInterestInsured` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 259 | `.Aneka.Obligee.AuctionDate` | 4 | 2 | ASM-FW-GISFW-Work | pxDateTime | 4/4 | relatif (agregat case) |
| 260 | `.Aneka.Obligee.AuctionNo` | 4 | 2 | ASM-FW-GISFW-Work | pxTextInput | 4/4 | relatif (agregat case) |
| 261 | `.Aneka.Obligee.ProjectValue` | 4 | 2 | ASM-FW-GISFW-Work | pxTextInput | 4/4 | relatif (agregat case) |
| 262 | `.Aneka.Obligee.ScheduleType` | 4 | 2 | ASM-FW-GISFW-Work | pxRadioButtons | 4/4 | relatif (agregat case) |
| 263 | `.Aneka.Occupation` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 264 | `.Aneka.OtherLOLDesc` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 265 | `.Aneka.PersonInsured` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 266 | `.Aneka.Remark` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 267 | `.Aneka.ScopeOfCover` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 268 | `.Aneka.SpecialCondition` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 269 | `.Aneka.TerritorialLimits` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 270 | `.Aneka.TradingWarranties` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 4/4 | relatif (agregat case) |
| 271 | `.AreaHectar` | 4 | 4 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 4/4 | relatif (agregat case) |
| 272 | `.ASMCity` | 4 | 4 | Data-Address | pxDropdown | 4/4 | relatif (agregat case) |
| 273 | `.ASMClass` | 4 | 4 | Data-Party-Person | pxAutoComplete | 4/4 | relatif (agregat case) |
| 274 | `.ASMCoverage(1).FacOutObjectList(1).ObjectPremi` | 4 | 4 | Data-Party-Person | pxNumber | 2/4 | relatif (agregat case) |
| 275 | `.ASMCoverage(1).FacOutObjectList(1).PercentOffered` | 4 | 4 | Data-Party-Person | pxNumber | 4/4 | relatif (agregat case) |
| 276 | `.ASMCoverage(1).FacOutObjectList(1).Rate` | 4 | 4 | Data-Party-Person | pxNumber | 4/4 | relatif (agregat case) |
| 277 | `.ASMCoverage(1).FacOutObjectList(1).RiComm` | 4 | 4 | Data-Party-Person | pxNumber | 4/4 | relatif (agregat case) |
| 278 | `.ASMCoverage(1).FacOutObjectList(1).ShareOffered` | 4 | 4 | Data-Party-Person | pxNumber | 4/4 | relatif (agregat case) |
| 279 | `.ASMCoverage(1).FacOutTSI` | 4 | 4 | Data-Party-Person | pxNumber | 4/4 | relatif (agregat case) |
| 280 | `.ASMCoverage(1).NominalShareOffered` | 4 | 4 | Data-Party-Person | pxNumber | 4/4 | relatif (agregat case) |
| 281 | `.ASMCoverage(1).PercentFacOut` | 4 | 4 | Data-Party-Person | pxNumber | 4/4 | relatif (agregat case) |
| 282 | `.ASMZipCode` | 4 | 4 | Data-Address | pxTextInput | 4/4 | relatif (agregat case) |
| 283 | `.BalanceBeforePPH` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 4/4 | relatif (agregat case) |
| 284 | `.BalanceBeforeTax` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 4/4 | relatif (agregat case) |
| 285 | `.BizName` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDisplayText/pxTextInput | 4/4 | relatif (agregat case) |
| 286 | `.CARI15` | 4 | 4 | ASM-FW-GISFW-Data-Search | pxDisplayText/pxNumber | 2/4 | relatif (agregat case) |
| 287 | `.CarParking` | 4 | 4 | ASM-FW-GISFW-Data-Vehicle | pxDropdown | 4/4 | relatif (agregat case) |
| 288 | `.CauseOfLoss` | 4 | 4 | ASM-FW-GISFW-Data-CauseOfLoss | pxTextInput | 4/4 | relatif (agregat case) |
| 289 | `.CC` | 4 | 4 | ASM-FW-GISFW-Data-Vehicle | pxTextInput | 4/4 | relatif (agregat case) |
| 290 | `.ClaimPaymentType` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDropdown | 2/4 | relatif (agregat case) |
| 291 | `.ClaimType` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDropdown | 2/4 | relatif (agregat case) |
| 292 | `.ClientName` | 4 | 3 | ASM-FW-GISFW-Int-AGENT | pxTextInput | 2/4 | relatif (agregat case) |
| 293 | `.ColorName` | 4 | 4 | ASM-FW-GISFW-Data-Vehicle | pxAutoComplete | 4/4 | relatif (agregat case) |
| 294 | `.Construction` | 4 | 4 | ASM-FW-GISFW-Data-Property | pxAutoComplete/pxDisplayText | 4/4 | relatif (agregat case) |
| 295 | `.ConveyanceDutyRange` | 4 | 4 | ASM-FW-GISFW-Data-Conveyance | pxDisplayText/pxTextInput | 3/4 | relatif (agregat case) |
| 296 | `.CountAttach` | 4 | 4 | ASM-FW-GISFW-Int-OFFERJSON | pxTextInput | 0/4 | relatif (agregat case) |
| 297 | `.CoverageBasis` | 4 | 4 | ASM-FW-GISFW-Data-Coverage | pxDropdown | 4/4 | relatif (agregat case) |
| 298 | `.CoverageInitial` | 4 | 4 | ASM-FW-GISFW-Data-Coverage | pxTextArea/pxTextInput | 4/4 | relatif (agregat case) |
| 299 | `.DateofSurvey` | 4 | 4 | ASM-FW-GISFW-Data-Quotation | pxDateTime | /4 | relatif (agregat case) |
| 300 | `.DeductibleType` | 4 | 4 | ASM-FW-GISFW-Data-Deductible | pxDropdown/pxTextInput | 4/4 | relatif (agregat case) |
| 301 | `.Deduction1` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 302 | `.Description` | 4 | 4 | ASM-FW-GISFW-Int-TABLEOFLIMIT | pxDisplayText/pxTextArea/pxTextInput | 3/4 | relatif (agregat case) |
| 303 | `.DueTo` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxRadioButtons | 0/4 | relatif (agregat case) |
| 304 | `.EmlPml` | 4 | 4 | ASM-FW-GISFW-Data-Coverage | pxNumber | 4/4 | relatif (agregat case) |
| 305 | `.ExcessLoss` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 306 | `.FEAID` | 4 | 4 | ASM-FW-GISFW-Data-FEA | pxDropdown | 3/4 | relatif (agregat case) |
| 307 | `.FEAInside` | 4 | 4 | ASM-FW-GISFW-Data-FEA | pxInteger | 4/4 | relatif (agregat case) |
| 308 | `.FEAOutside` | 4 | 4 | ASM-FW-GISFW-Data-FEA | pxInteger | 4/4 | relatif (agregat case) |
| 309 | `.FirstLoss` | 4 | 4 | ASM-FW-GISFW-Data-Coverage | pxNumber | 4/4 | relatif (agregat case) |
| 310 | `.FirstScale` | 4 | 4 | ASM-FW-GISFW-Data-Coverage | pxNumber | 4/4 | relatif (agregat case) |
| 311 | `.FlagPPH` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCheckbox | 0/4 | relatif (agregat case) |
| 312 | `.IsUsedFlag` | 4 | 4 | ASM-FW-GISFW-Data-Occupation | pxCheckbox | 3/4 | relatif (agregat case) |
| 313 | `.LayerPart` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 4/4 | relatif (agregat case) |
| 314 | `.LayerPartType` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 4/4 | relatif (agregat case) |
| 315 | `.LayerType` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 4/4 | relatif (agregat case) |
| 316 | `.MaintenancePeriod` | 4 | 2 | ASM-FW-GISFW-Data-Aneka | pxNumber | 4/4 | relatif (agregat case) |
| 317 | `.NetPremi` | 4 | 3 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 2/4 | relatif (agregat case) |
| 318 | `.NetPremi2` | 4 | 3 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 2/4 | relatif (agregat case) |
| 319 | `.NetPremium` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 4/4 | relatif (agregat case) |
| 320 | `.NM_SHIP` | 4 | 4 | ASM-FW-GISFW-Int-MSHIP | pxTextInput | 2/4 | relatif (agregat case) |
| 321 | `.NoOffer` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDisplayText | 4/4 | relatif (agregat case) |
| 322 | `.NoOfTree` | 4 | 4 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 4/4 | relatif (agregat case) |
| 323 | `.OfferFacIn.Comment` | 4 | 4 | ASM-FW-GISFW-Work | pxTextArea | 2/4 | relatif (agregat case) |
| 324 | `.OfferFacIn.FacRetroDetails.BackUpStatus` | 4 | 4 | ASM-FW-GISFW-Work | pxDropdown | 3/4 | relatif (agregat case) |
| 325 | `.OfferFacIn.Obligee.AuctionDate` | 4 | 2 | ASM-FW-GISFW-Work | pxDateTime | 4/4 | relatif (agregat case) |
| 326 | `.OfferFacIn.Obligee.AuctionNo` | 4 | 2 | ASM-FW-GISFW-Work | pxTextInput | 4/4 | relatif (agregat case) |
| 327 | `.OfferFacIn.Obligee.ProjectValue` | 4 | 2 | ASM-FW-GISFW-Work | pxTextInput | 4/4 | relatif (agregat case) |
| 328 | `.OfferFacIn.Obligee.ScheduleType` | 4 | 2 | ASM-FW-GISFW-Work | pxRadioButtons | 4/4 | relatif (agregat case) |
| 329 | `.OutstandingClaim` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 330 | `.OveriddingCommOgp` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 2/4 | relatif (agregat case) |
| 331 | `.OveriddingCommOnp` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 2/4 | relatif (agregat case) |
| 332 | `.Policy.Payment.Deduction2` | 4 | 4 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 2/4 | relatif (agregat case) |
| 333 | `.PolicyData.BLNumber` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextArea | 4/4 | relatif (agregat case) |
| 334 | `.PolicyData.InvoiceDate` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxDateTime | 4/4 | relatif (agregat case) |
| 335 | `.PolicyData.InvoiceNumber` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextInput | 4/4 | relatif (agregat case) |
| 336 | `.PolicyData.LC.LC` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxCheckbox | 2/4 | relatif (agregat case) |
| 337 | `.PolicyData.LC.LCBANK` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextInput | 4/4 | relatif (agregat case) |
| 338 | `.PolicyData.LC.LCCondition` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextInput | 4/4 | relatif (agregat case) |
| 339 | `.PolicyData.LC.LCEndDate` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxDateTime | 4/4 | relatif (agregat case) |
| 340 | `.PolicyData.LC.LCNumber` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextInput | 4/4 | relatif (agregat case) |
| 341 | `.PolicyData.LC.LCRemark` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextInput | 4/4 | relatif (agregat case) |
| 342 | `.PolicyData.MaxAge` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxAutoComplete | 4/4 | relatif (agregat case) |
| 343 | `.PolicyData.MaxAgeDetail` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextArea | 4/4 | relatif (agregat case) |
| 344 | `.PolicyData.Ship.AGE` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxDisplayText | 4/4 | relatif (agregat case) |
| 345 | `.PolicyData.Ship.CONST` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxDisplayText | 4/4 | relatif (agregat case) |
| 346 | `.PolicyData.Ship.DWT` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxNumber | 4/4 | relatif (agregat case) |
| 347 | `.PolicyData.Ship.FLAG_NAME` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxDisplayText | 4/4 | relatif (agregat case) |
| 348 | `.PolicyData.Ship.GRT` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxNumber | 4/4 | relatif (agregat case) |
| 349 | `.PolicyData.Ship.NRT` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxNumber | 4/4 | relatif (agregat case) |
| 350 | `.PolicyData.Ship.REMARK` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxDisplayText | 4/4 | relatif (agregat case) |
| 351 | `.PolicyData.Ship.SHIPTYPE_NAME` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxDropdown | 4/4 | relatif (agregat case) |
| 352 | `.PolicyData.Ship.Y_MAKE1` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxDisplayText | 4/4 | relatif (agregat case) |
| 353 | `.PolicyData.SurveyAgent.ASMAgentID` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxAutoComplete | 4/4 | relatif (agregat case) |
| 354 | `.PolicyData.SurveyAgent.pyAddress` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextInput | 4/4 | relatif (agregat case) |
| 355 | `.PolicyData.SurveyAgent.pyFullName` | 4 | 4 | ASM-FW-GISFW-Data-Cargo | pxTextArea/pxTextInput | 4/4 | relatif (agregat case) |
| 356 | `.PPHValue` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 4/4 | relatif (agregat case) |
| 357 | `.PPNValue` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 4/4 | relatif (agregat case) |
| 358 | `.PremiLifeNusantaraRe` | 4 | 2 | ASM-FW-GISFW-Data-Coverage | pxNumber | 0/4 | relatif (agregat case) |
| 359 | `.PremiOgp` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 360 | `.PremiOnp` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 361 | `.Price` | 4 | 4 | ASM-FW-GISFW-Data-Accessory | pxNumber | 4/4 | relatif (agregat case) |
| 362 | `.Property.AlmRiskID` | 4 | 4 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 4/4 | relatif (agregat case) |
| 363 | `.Property.BuildingNo` | 4 | 4 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | 4/4 | relatif (agregat case) |
| 364 | `.Property.RoadName` | 4 | 4 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextArea/pxTextInput | 3/4 | relatif (agregat case) |
| 365 | `.Property.RoadType` | 4 | 4 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 4/4 | relatif (agregat case) |
| 366 | `.ProRatePercent` | 4 | 4 | ASM-FW-GISFW-Data-Coverage | pxNumber/pxTextInput | 4/4 | relatif (agregat case) |
| 367 | `.pyNote` | 4 | 2 | Data-WorkAttach-Note | pxDisplayText/pxTextArea/pxTextInput | 3/4 | relatif (agregat case) |
| 368 | `.Quartal` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 2/4 | relatif (agregat case) |
| 369 | `.QuotationData.InsuredName` | 4 | 4 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 4/4 | relatif (agregat case) |
| 370 | `.QuotationData.ProportionalType` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDropdown | 4/4 | relatif (agregat case) |
| 371 | `.QuotationData.SobName` | 4 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText | 4/4 | relatif (agregat case) |
| 372 | `.RateLifeRetro` | 4 | 2 | ASM-FW-GISFW-Data-Coverage | pxTextInput | 0/4 | relatif (agregat case) |
| 373 | `.ResultOgp1` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 374 | `.ResultOgp2` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 375 | `.ResultOnp1` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 376 | `.ResultOnp2` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 377 | `.RiCommOgp` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 2/4 | relatif (agregat case) |
| 378 | `.RiCommOnp` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 2/4 | relatif (agregat case) |
| 379 | `.SalvageValue` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 2/4 | relatif (agregat case) |
| 380 | `.ScoringRisk.FinalScore` | 4 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxInteger/pxNumber | 4/4 | relatif (agregat case) |
| 381 | `.ScoringRisk.NoteFinalScore` | 4 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber/pxTextInput | 4/4 | relatif (agregat case) |
| 382 | `.ShareCurrency` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 4/4 | relatif (agregat case) |
| 383 | `.ShareValue` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCurrency | 4/4 | relatif (agregat case) |
| 384 | `.SourceBizName` | 4 | 2 | ASM-FW-GISFW-Data-Outgo | pxAutoComplete | 4/4 | relatif (agregat case) |
| 385 | `.StartDate` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDateTime | 2/4 | relatif (agregat case) |
| 386 | `.StatementDate` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDateTime | 4/4 | relatif (agregat case) |
| 387 | `.StatementType` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDropdown | 2/4 | relatif (agregat case) |
| 388 | `.Status` | 4 | 4 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText/pxRadioButtons | /4 | relatif (agregat case) |
| 389 | `.SurveyedBy` | 4 | 4 | ASM-FW-GISFW-Data-Quotation | pxTextInput | /4 | relatif (agregat case) |
| 390 | `.TotalPremiumNusantaraReOld` | 4 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 4/4 | relatif (agregat case) |
| 391 | `.Tranship` | 4 | 4 | ASM-FW-GISFW-Data-Trading | pxDisplayText/pxTextInput | 4/4 | relatif (agregat case) |
| 392 | `.TransmissionType` | 4 | 4 | ASM-FW-GISFW-Data-Vehicle | pxDropdown | 4/4 | relatif (agregat case) |
| 393 | `.TreatyGroupName` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDisplayText/pxTextInput | 4/4 | relatif (agregat case) |
| 394 | `.TypeTax` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDisplayText/pxRadioButtons | 2/4 | relatif (agregat case) |
| 395 | `.UpdateDate` | 4 | 4 | ASM-FW-GISFW-Data-FEA | pxDateTime | 4/4 | relatif (agregat case) |
| 396 | `.ValueName` | 4 | 4 | ASM-FW-GISFW-Data-Accessory | pxAutoComplete | 4/4 | relatif (agregat case) |
| 397 | `.VOYAGE_NO` | 4 | 4 | ASM-FW-GISFW-Int-MSHIP | pxTextInput | 2/4 | relatif (agregat case) |
| 398 | `.Warehouse` | 4 | 4 | ASM-FW-GISFW-Data-Trading | pxDisplayText/pxTextInput | 4/4 | relatif (agregat case) |
| 399 | `.YearOfQuartal` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 2/4 | relatif (agregat case) |
| 400 | `pyWorkPage.OfferFacIn.FacRetroDetails.BackUpStatus` | 4 | 2 | ASM-FW-GISFW-Data-FacOffer | pxDropdown | 4/4 | agregat case (eksplisit) |
| 401 | `pyWorkPage.OfferFacIn.QuotationData.InsuredName` | 4 | 4 | ASM-FW-GISFW-Data-CauseOfLoss | pxTextInput | 4/4 | agregat case (eksplisit) |
| 402 | `pyWorkPage.TreatyIn.Commencement` | 4 | 4 | (kosong) | pxDateTime | 4/4 | agregat case (eksplisit) |
| 403 | `pyWorkPage.TreatyIn.Termination` | 4 | 4 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDateTime | 4/4 | agregat case (eksplisit) |
| 404 | `.Accumulation` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextInput | 3/3 | relatif (agregat case) |
| 405 | `.AccumulationName` | 3 | 3 | ASM-FW-GISFW-Int-ACCUMULATION | pxAutoComplete/pxDisplayText | 2/3 | relatif (agregat case) |
| 406 | `.Address` | 3 | 3 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 3/3 | relatif (agregat case) |
| 407 | `.AlmRiskID` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 408 | `.APAR` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 409 | `.ApplicationNumber` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextInput | 3/3 | relatif (agregat case) |
| 410 | `.ArgumentDescription` | 3 | 3 | ASM-FW-GISFW-Data-Argument | pxTextInput | 3/3 | relatif (agregat case) |
| 411 | `.ArgumentNumber` | 3 | 3 | ASM-FW-GISFW-Data-Argument | pxTextInput | 3/3 | relatif (agregat case) |
| 412 | `.ArgumentValue` | 3 | 3 | ASM-FW-GISFW-Data-Argument | pxTextInput | 0/3 | relatif (agregat case) |
| 413 | `.ASMCCAmount` | 3 | 3 | Data-Party-Person | pxInteger/pxTextInput | 3/3 | relatif (agregat case) |
| 414 | `.ASMDistrict` | 3 | 3 | Data-Address | pxDropdown | 3/3 | relatif (agregat case) |
| 415 | `.ASMRate` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 0/3 | relatif (agregat case) |
| 416 | `.ASMRIComm` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 0/3 | relatif (agregat case) |
| 417 | `.ASMRW` | 3 | 3 | Data-Address | pxDropdown | 3/3 | relatif (agregat case) |
| 418 | `.ASMShareTotal` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | 3/3 | relatif (agregat case) |
| 419 | `.AviationHull.MaxPassengers` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 420 | `.AviationHull.Pilots` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 421 | `.AviationHull.RegistrationMarks` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 422 | `.AviationHull.SpecialRentalUses` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 423 | `.AviationHull.SpecialUses` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 424 | `.AviationHull.StandardUses` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 425 | `.BasisType` | 3 | 3 | ASM-FW-GISFW-Data-Deductible | pxDropdown | 3/3 | relatif (agregat case) |
| 426 | `.BuildingConstruction.FloorType` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 427 | `.BuildingConstruction.NumberOfFloor` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxNumber | 3/3 | relatif (agregat case) |
| 428 | `.BuildingConstruction.OthersType` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 429 | `.BuildingConstruction.PartitionType` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 430 | `.BuildingConstruction.RoofType` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 431 | `.BuildingConstruction.SupportWallType` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 432 | `.BuildingConstruction.WallType` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 433 | `.BuildingNo` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 434 | `.BuildingStatus` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextArea | 3/3 | relatif (agregat case) |
| 435 | `.BuiltIn` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 436 | `.CARI16` | 3 | 3 | ASM-FW-GISFW-Data-Search | pxTextInput | /3 | relatif (agregat case) |
| 437 | `.CARI2` | 3 | 3 | ASM-FW-GISFW-Data-Search | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 438 | `.CARI3` | 3 | 3 | ASM-FW-GISFW-Data-Search | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 439 | `.ClaimCategory` | 3 | 3 | ASM-FW-GISFW-Data-Deductible | pxDropdown | 3/3 | relatif (agregat case) |
| 440 | `.ClassOfBusiness` | 3 | 3 | ASM-FW-GISFW-Data-CauseOfLoss | pxAutoComplete/pxNumber/pxTextInput | 2/3 | relatif (agregat case) |
| 441 | `.Collateral` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextInput | 3/3 | relatif (agregat case) |
| 442 | `.ConclusionsProposals` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextArea | 3/3 | relatif (agregat case) |
| 443 | `.COVERAGE_NAME` | 3 | 3 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 3/3 | relatif (agregat case) |
| 444 | `.CoverageList(1).FacOutObjectList(1).ObjectPremi` | 3 | 3 | Data-Party-Person | pxNumber | /3 | relatif (agregat case) |
| 445 | `.CoverageList(1).FacOutObjectList(1).PercentOffered` | 3 | 3 | Data-Party-Person | pxNumber | 3/3 | relatif (agregat case) |
| 446 | `.CoverageList(1).FacOutObjectList(1).RiComm` | 3 | 3 | Data-Party-Person | pxNumber | 3/3 | relatif (agregat case) |
| 447 | `.CoverageList(1).FacOutTSI` | 3 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 3/3 | relatif (agregat case) |
| 448 | `.CoverageList(1).NominalShareOffered` | 3 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 3/3 | relatif (agregat case) |
| 449 | `.CoverageList(1).PercentFacOut` | 3 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 3/3 | relatif (agregat case) |
| 450 | `.Currency.Currency` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxAutoComplete | /3 | relatif (agregat case) |
| 451 | `.DataFEA.PrivateTruckBrigade` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 452 | `.Date` | 3 | 3 | ASM-FW-GISFW-Data-SuggestList | pxDateTime | 2/3 | relatif (agregat case) |
| 453 | `.DraftWordingID` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxDropdown | 3/3 | relatif (agregat case) |
| 454 | `.EML` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxNumber | 3/3 | relatif (agregat case) |
| 455 | `.EndorsmentReason` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText/pxTextInput | 3/3 | relatif (agregat case) |
| 456 | `.FacInGrossPremiumASM` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | 3/3 | relatif (agregat case) |
| 457 | `.FacInNetPremiumASM` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | 3/3 | relatif (agregat case) |
| 458 | `.FirstLossLimitFac` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 459 | `.Flag` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 460 | `.FlagSpreading` | 3 | 3 | ASM-FW-GISFW-Data-SpreadingRisk | pxDropdown | 0/3 | relatif (agregat case) |
| 461 | `.FollowingNB` | 3 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | /3 | relatif (agregat case) |
| 462 | `.GeographicalLimits` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 463 | `.GoodName` | 3 | 3 | ASM-FW-GISFW-Data-Good | pxDisplayText/pxTextArea | 3/3 | relatif (agregat case) |
| 464 | `.Hydrant` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 465 | `.InfoFEA` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDisplayText/pxTextArea | 3/3 | relatif (agregat case) |
| 466 | `.InsuredSiup` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextInput | 3/3 | relatif (agregat case) |
| 467 | `.InvoiceNo` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 468 | `.IsAdjustableFlag` | 3 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxCheckbox | /3 | relatif (agregat case) |
| 469 | `.IsFlammableItemFlag` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxCheckbox | /3 | relatif (agregat case) |
| 470 | `.IsHotWorkProcessFlag` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxCheckbox | 2/3 | relatif (agregat case) |
| 471 | `.IsProductionProcessFlag` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxCheckbox | 2/3 | relatif (agregat case) |
| 472 | `.IsShowDetail` | 3 | 3 | ASM-FW-GISFW-Work | pxCheckbox | 0/3 | relatif (agregat case) |
| 473 | `.IsSpecialAcceptance` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | 0/3 | relatif (agregat case) |
| 474 | `.Job` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxDisplayText/pxTextInput | 3/3 | relatif (agregat case) |
| 475 | `.KeyPerson` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextInput | 3/3 | relatif (agregat case) |
| 476 | `.LimitIdemnityDay` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 477 | `.LimitIdemnityMonth` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 478 | `.LOSS_AMOUNT` | 3 | 3 | ASM-FW-GISFW-Data-CauseOfLoss | pxNumber | 2/3 | relatif (agregat case) |
| 479 | `.LOSS_AMOUNT_IDR` | 3 | 3 | ASM-FW-GISFW-Data-CauseOfLoss | pxNumber | 2/3 | relatif (agregat case) |
| 480 | `.MarineHull.BREADTHM` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 481 | `.MarineHull.Classification` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 482 | `.MarineHull.Construction` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 483 | `.MarineHull.DEPTHM` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 484 | `.MarineHull.DWTRT` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 485 | `.MarineHull.GRTRT` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 486 | `.MarineHull.LENGTHM` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 487 | `.MarineHull.NRTRT` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 488 | `.MarineHull.TypeOfVessel` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 489 | `.NoOfferSlip` | 3 | 2 | ASM-FW-GISFW-Data-FacOffer | pxTextInput | 2/3 | relatif (agregat case) |
| 490 | `.ObjectType` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 491 | `.OCCUPATION_NAME` | 3 | 3 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 3/3 | relatif (agregat case) |
| 492 | `.OfferFacIn.OfficerFacOut` | 3 | 3 | ASM-FW-GISFW-Work | pxTextInput | 3/3 | relatif (agregat case) |
| 493 | `.OfferFacIn.PolicyData.EndDateTime` | 3 | 3 | ASM-FW-GISFW-Work | pxDateTime | 3/3 | relatif (agregat case) |
| 494 | `.OfferFacIn.PolicyData.StartDateTime` | 3 | 3 | ASM-FW-GISFW-Work | pxDateTime | 3/3 | relatif (agregat case) |
| 495 | `.OfferFacIn.QuotationData.BusinessName` | 3 | 3 | ASM-FW-GISFW-Work | pxTextInput | 3/3 | relatif (agregat case) |
| 496 | `.OfferFacIn.QuotationData.InsuredName` | 3 | 3 | ASM-FW-GISFW-Work | pxTextInput | 3/3 | relatif (agregat case) |
| 497 | `.OfferFacIn.TotalTSITopRisk` | 3 | 2 | ASM-FW-GISFW-Work | pxNumber | 3/3 | relatif (agregat case) |
| 498 | `.OriginalRate` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextInput | 3/3 | relatif (agregat case) |
| 499 | `.Ownership` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 500 | `.PackingID` | 3 | 3 | ASM-FW-GISFW-Data-Packing | pxDropdown | 3/3 | relatif (agregat case) |
| 501 | `.PctLoL` | 3 | 3 | ASM-FW-GISFW-Data-Coverage | pxNumber | 2/3 | relatif (agregat case) |
| 502 | `.PctOfContribution` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 503 | `.PctShareAllObj` | 3 | 3 | ASM-FW-GISFW-Data-FacOffer | pxNumber | 2/3 | relatif (agregat case) |
| 504 | `.PctShortPeriod` | 3 | 2 | ASM-FW-GISFW-Data-Coverage | pxTextInput | 3/3 | relatif (agregat case) |
| 505 | `.PERCENTSHARE` | 3 | 3 | ASM-FW-GISFW-Data-Coins | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 506 | `.PML` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxNumber | 3/3 | relatif (agregat case) |
| 507 | `.PNOTE` | 3 | 3 | Link-Attachment | (tanpa) | 0/3 | relatif (agregat case) |
| 508 | `.PolicyData.EndDateTime` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDateTime | 2/3 | relatif (agregat case) |
| 509 | `.PolicyData.OfferingDate` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDateTime | 3/3 | relatif (agregat case) |
| 510 | `.PolicyData.StartDateTime` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDateTime | 3/3 | relatif (agregat case) |
| 511 | `.PPh` | 3 | 3 | ASM-FW-GISFW-Data-Installment | pxNumber | 3/3 | relatif (agregat case) |
| 512 | `.PPN` | 3 | 3 | ASM-FW-GISFW-Data-Installment | pxNumber | 3/3 | relatif (agregat case) |
| 513 | `.PremiNusantaraReOld` | 3 | 3 | ASM-FW-GISFW-Data-Coverage | pxNumber | 3/3 | relatif (agregat case) |
| 514 | `.PrincipalsBackground` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextArea | 3/3 | relatif (agregat case) |
| 515 | `.PrincipalsData` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextArea | 3/3 | relatif (agregat case) |
| 516 | `.PrincipalsGuarantee` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextArea | 3/3 | relatif (agregat case) |
| 517 | `.PrincipalsOutGo` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextInput | 3/3 | relatif (agregat case) |
| 518 | `.PropertiItemNote` | 3 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxTextArea | 3/3 | relatif (agregat case) |
| 519 | `.pyCity` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 520 | `.pyCountryName` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 521 | `.pyPhoneNumber` | 3 | 3 | Data-Party-Person | pxTextInput | 2/3 | relatif (agregat case) |
| 522 | `.pySelected` | 3 | 3 | Data-Party-Person | pxCheckbox | 0/3 | relatif (agregat case) |
| 523 | `.QuotationData.BusinessName` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText | 3/3 | relatif (agregat case) |
| 524 | `.QuotationData.EdmDate` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDateTime | 3/3 | relatif (agregat case) |
| 525 | `.QuotationData.GroupName` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText/pxTextInput | 3/3 | relatif (agregat case) |
| 526 | `.QuotationData.MOID` | 3 | 3 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDropdown | 0/3 | relatif (agregat case) |
| 527 | `.QuotationData.Period` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /3 | relatif (agregat case) |
| 528 | `.QuotationData.QQName` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 2/3 | relatif (agregat case) |
| 529 | `.QuotationData.StatusBusiness` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | 3/3 | relatif (agregat case) |
| 530 | `.RateOld` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | 3/3 | relatif (agregat case) |
| 531 | `.RIComm` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxNumber | 3/3 | relatif (agregat case) |
| 532 | `.RiCommAllObj` | 3 | 3 | ASM-FW-GISFW-Data-FacOffer | pxNumber | 2/3 | relatif (agregat case) |
| 533 | `.RiskCategory` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 534 | `.RiskLocation.ASMCity` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 535 | `.RiskLocation.ASMDistrict` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 536 | `.RiskLocation.ASMRW` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 537 | `.RiskLocation.ASMZipCode` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 538 | `.RoadName` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextInput | 3/3 | relatif (agregat case) |
| 539 | `.RoadType` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 540 | `.SelectedLocationAddress` | 3 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxAutoComplete | 0/3 | relatif (agregat case) |
| 541 | `.SelectedObjectItem` | 3 | 3 | ASM-FW-GISFW-Data-PropertyItem | pxAutoComplete | 0/3 | relatif (agregat case) |
| 542 | `.SerialNo` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 543 | `.SmokeDetector` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 544 | `.Sprinkler` | 3 | 3 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 545 | `.Sublimit` | 3 | 3 | ASM-FW-GISFW-Data-Coverage | pxDisplayText/pxNumber | 3/3 | relatif (agregat case) |
| 546 | `.SumDiscountObjectAneka` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 547 | `.SummaryFinancial` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextArea | 3/3 | relatif (agregat case) |
| 548 | `.SumPremiObjectAneka` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 3/3 | relatif (agregat case) |
| 549 | `.SupportingDocument` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextArea | 3/3 | relatif (agregat case) |
| 550 | `.SurroundingRisk.BackConstruction` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 551 | `.SurroundingRisk.BackDistance` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxNumber | 3/3 | relatif (agregat case) |
| 552 | `.SurroundingRisk.BackNote` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextArea/pxTextInput | 3/3 | relatif (agregat case) |
| 553 | `.SurroundingRisk.BackOccupation` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxAutoComplete | 3/3 | relatif (agregat case) |
| 554 | `.SurroundingRisk.FloodArea` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 555 | `.SurroundingRisk.FloodAreaStatus` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 556 | `.SurroundingRisk.FrontConstruction` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 557 | `.SurroundingRisk.FrontDistance` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxNumber | 3/3 | relatif (agregat case) |
| 558 | `.SurroundingRisk.FrontNote` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextArea/pxTextInput | 3/3 | relatif (agregat case) |
| 559 | `.SurroundingRisk.FrontOccupation` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxAutoComplete | 3/3 | relatif (agregat case) |
| 560 | `.SurroundingRisk.HousekeepingRemark` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextArea | 3/3 | relatif (agregat case) |
| 561 | `.SurroundingRisk.HousekeepingStatus` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 562 | `.SurroundingRisk.LeftConstruction` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 563 | `.SurroundingRisk.LeftDistance` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxNumber | 3/3 | relatif (agregat case) |
| 564 | `.SurroundingRisk.LeftNote` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextArea/pxTextInput | 3/3 | relatif (agregat case) |
| 565 | `.SurroundingRisk.LeftOccupation` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxAutoComplete | 3/3 | relatif (agregat case) |
| 566 | `.SurroundingRisk.RightConstruction` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 567 | `.SurroundingRisk.RightDistance` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxNumber | 3/3 | relatif (agregat case) |
| 568 | `.SurroundingRisk.RightNote` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxTextArea/pxTextInput | 3/3 | relatif (agregat case) |
| 569 | `.SurroundingRisk.RightOccupation` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxAutoComplete | 3/3 | relatif (agregat case) |
| 570 | `.TextNoQuotation` | 3 | 3 | ASM-FW-SFAGISFW-Work-Opportunity | pxDisplayText | 3/3 | relatif (agregat case) |
| 571 | `.TotalNetRate` | 3 | 2 | ASM-FW-GISFW-Data-PropertyItem | pxNumber | 3/3 | relatif (agregat case) |
| 572 | `.TotalOffered` | 3 | 3 | ASM-FW-GISFW-Data-FacOffer | pxNumber | 3/3 | relatif (agregat case) |
| 573 | `.TreatyGroup` | 3 | 3 | ASM-FW-GISFW-Data-Property | pxDropdown | 3/3 | relatif (agregat case) |
| 574 | `.TypeOfWording` | 3 | 3 | ASM-FW-GISFW-Data-MemorandumAccept | pxTextInput | 3/3 | relatif (agregat case) |
| 575 | `.VehicleHE.BrandName` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxAutoComplete | 3/3 | relatif (agregat case) |
| 576 | `.VehicleHE.ObjectNameHE` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxAutoComplete | 3/3 | relatif (agregat case) |
| 577 | `.VehicleHE.TypeName` | 3 | 3 | ASM-FW-GISFW-Data-Aneka | pxAutoComplete | 3/3 | relatif (agregat case) |
| 578 | `.WarrantyDesc` | 3 | 3 | ASM-FW-GISFW-Data-Warranty | pxTextArea | 2/3 | relatif (agregat case) |
| 579 | `A.NBStatus` | 3 | 3 | ASM-FW-SFAGISFW-Work-Opportunity | pxDisplayText | 3/3 | LUAR agregat: A |
| 580 | `A.Quotation.BusinessName` | 3 | 3 | ASM-FW-SFAGISFW-Work-Opportunity | pxDisplayText | 3/3 | LUAR agregat: A |
| 581 | `A.Quotation.InsuredName` | 3 | 3 | ASM-FW-SFAGISFW-Work-Opportunity | pxDisplayText | 3/3 | LUAR agregat: A |
| 582 | `A.Quotation.MarketingName` | 3 | 3 | ASM-FW-SFAGISFW-Work-Opportunity | (tanpa) | 0/3 | LUAR agregat: A |
| 583 | `CARI.pyName` | 3 | 3 | ASM-FW-GISFW-Data-Occupation | pxTextInput | 0/3 | LUAR agregat: CARI |
| 584 | `InputData.CARI12` | 3 | 2 | ASM-FW-GISFW-Work | pxTextInput | 3/3 | LUAR agregat: InputData |
| 585 | `ProtectSpreading.CARI14` | 3 | 3 | ASM-FW-GISFW-Work | pxDisplayText | 3/3 | LUAR agregat: ProtectSpreading |
| 586 | `pyWorkPage.FlagErrorKonversi` | 3 | 3 | ASM-FW-GISFW-Work | (tanpa) | 3/3 | agregat case (eksplisit) |
| 587 | `pyWorkPage.OfferFacIn.QuotationData.ShareOfCeding` | 3 | 2 | ASM-FW-GISFW-Data-Quotation | pxTextInput | 3/3 | agregat case (eksplisit) |
| 588 | `pyWorkPage.PolicyTreatyIn.PolicyNo` | 3 | 3 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 3/3 | agregat case (eksplisit) |
| 589 | `TempPolis.CARI12` | 3 | 3 | ASM-FW-GISFW-Work | pxDropdown | 0/3 | LUAR agregat: TempPolis |
| 590 | `.AccumulationCode` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxTextInput | 2/2 | relatif (agregat case) |
| 591 | `.AccumulationDescription` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxDisplayText | 2/2 | relatif (agregat case) |
| 592 | `.Aneka.Obligee.AdministrationFee` | 2 | 2 | ASM-FW-GISFW-Work | pxInteger | 2/2 | relatif (agregat case) |
| 593 | `.Aneka.Obligee.AgreementDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 594 | `.Aneka.Obligee.AgreementNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 595 | `.Aneka.Obligee.BankAddress` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 596 | `.Aneka.Obligee.BankGuaranteeNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 597 | `.Aneka.Obligee.BankName` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 598 | `.Aneka.Obligee.CustomsActivity` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 599 | `.Aneka.Obligee.DateOfGuarantee` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 600 | `.Aneka.Obligee.FacilityDecreeDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 601 | `.Aneka.Obligee.FacilityDecreeNumber` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 602 | `.Aneka.Obligee.FacilityStatus` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 603 | `.Aneka.Obligee.InsuredNumber` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 604 | `.Aneka.Obligee.JenisPenjaminan` | 2 | 2 | ASM-FW-GISFW-Work | pxDropdown | 2/2 | relatif (agregat case) |
| 605 | `.Aneka.Obligee.KBGDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 606 | `.Aneka.Obligee.KBGNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 607 | `.Aneka.Obligee.LengthPeriod` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 608 | `.Aneka.Obligee.Location` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 609 | `.Aneka.Obligee.MaxOfLimitOfLiability` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 610 | `.Aneka.Obligee.NIPER` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 611 | `.Aneka.Obligee.OCDescription` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 612 | `.Aneka.Obligee.OtherFacilityNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 613 | `.Aneka.Obligee.Pasal` | 2 | 2 | ASM-FW-GISFW-Work | pxDropdown | 2/2 | relatif (agregat case) |
| 614 | `.Aneka.Obligee.PIBDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 615 | `.Aneka.Obligee.PIBNumber` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 616 | `.Aneka.Obligee.PKDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 617 | `.Aneka.Obligee.PKDate2` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 618 | `.Aneka.Obligee.PKNo1` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 619 | `.Aneka.Obligee.PKNo2` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 620 | `.Aneka.Obligee.PlaceOfEvent` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 621 | `.Aneka.Obligee.RegisterNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 622 | `.ASMAddressType` | 2 | 2 | Data-Address | pxDropdown | 2/2 | relatif (agregat case) |
| 623 | `.ASMAge` | 2 | 2 | Data-Party-Person | pxDisplayText | 2/2 | relatif (agregat case) |
| 624 | `.ASMBankData.ASMAccountNo` | 2 | 2 | Data-Party-Person | pxTextInput | 0/2 | relatif (agregat case) |
| 625 | `.ASMBankData.ASMBankName` | 2 | 2 | Data-Party-Person | pxAutoComplete | 0/2 | relatif (agregat case) |
| 626 | `.ASMBMI` | 2 | 2 | Data-Party-Person | pxNumber/pxTextInput | 2/2 | relatif (agregat case) |
| 627 | `.ASMDateTimeChronology` | 2 | 2 | Data-WorkAttach-Note | pxDateTime | 2/2 | relatif (agregat case) |
| 628 | `.ASMUser` | 2 | 2 | Data-WorkAttach-Note | pxDisplayText | 2/2 | relatif (agregat case) |
| 629 | `.BandingTo` | 2 | 2 | ASM-FW-GISFW-Work | pxDropdown | 0/2 | relatif (agregat case) |
| 630 | `.BeginDate` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxTextInput | 2/2 | relatif (agregat case) |
| 631 | `.BrokerageFee` | 2 | 2 | ASM-FW-GISFW-Data-Installment | pxNumber | 2/2 | relatif (agregat case) |
| 632 | `.CalculateMethod` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxDropdown | 2/2 | relatif (agregat case) |
| 633 | `.CARI17` | 2 | 2 | ASM-FW-GISFW-Data-Search | pxTextInput | /2 | relatif (agregat case) |
| 634 | `.CARI18` | 2 | 2 | ASM-FW-GISFW-Data-Search | (tanpa) | /2 | relatif (agregat case) |
| 635 | `.CARID1` | 2 | 2 | ASM-FW-GISFW-Data-Search | pxNumber | 2/2 | relatif (agregat case) |
| 636 | `.CARID2` | 2 | 2 | ASM-FW-GISFW-Data-Search | pxNumber | 2/2 | relatif (agregat case) |
| 637 | `.CEDING_CO_NAME` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 638 | `.CheckBox1` | 2 | 2 | ASM-FW-GISFW-Data-CauseOfDecline | pxCheckbox | 0/2 | relatif (agregat case) |
| 639 | `.CityName` | 2 | 2 | ASM-FW-GISFW-Int-RISKADDRESS | pxDisplayText | /2 | relatif (agregat case) |
| 640 | `.ClaimAmountIDR` | 2 | 2 | ASM-FW-GISFW-Data-SpreadingRisk | pxNumber | 2/2 | relatif (agregat case) |
| 641 | `.CLASS_OF_BUSINESS` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 642 | `.ClauseLanguage` | 2 | 2 | ASM-FW-GISFW-Data-Clause | pxTextInput | 2/2 | relatif (agregat case) |
| 643 | `.ClientID` | 2 | 2 | ASM-FW-GISFW-Int-AGENT | (tanpa) | 2/2 | relatif (agregat case) |
| 644 | `.COMMISION` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText/pxNumber | 2/2 | relatif (agregat case) |
| 645 | `.Commission` | 2 | 2 | ASM-FW-GISFW-Data-Installment | pxNumber | 2/2 | relatif (agregat case) |
| 646 | `.Conditions` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxTextArea | 2/2 | relatif (agregat case) |
| 647 | `.ConstructionEarthquake` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxAutoComplete | 2/2 | relatif (agregat case) |
| 648 | `.Country` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 649 | `.CoverageList(1).FacOutObjectList(1).commision` | 2 | 2 | ASM-FW-GISFW-Data-Cargo | pxNumber | 2/2 | relatif (agregat case) |
| 650 | `.CurrencyMaster` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxTextInput | 2/2 | relatif (agregat case) |
| 651 | `.CurrentYear` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-Currency | (tanpa) | /2 | relatif (agregat case) |
| 652 | `.DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 653 | `.DataScoringRiskList(1).FEA.FireAlarmSystem.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 654 | `.DataScoringRiskList(1).FEA.FireAlarmSystem.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 655 | `.DataScoringRiskList(1).FEA.FireAlarmSystem.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 656 | `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 657 | `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 658 | `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 659 | `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 660 | `.DataScoringRiskList(1).FEA.Firebrigade.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 661 | `.DataScoringRiskList(1).FEA.Firebrigade.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 662 | `.DataScoringRiskList(1).FEA.Firebrigade.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 663 | `.DataScoringRiskList(1).FEA.Firebrigade.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 664 | `.DataScoringRiskList(1).FEA.Firebrigade.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 665 | `.DataScoringRiskList(1).FEA.Firebrigade.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 666 | `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 667 | `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 668 | `.DataScoringRiskList(1).FEA.HydrantSystem.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 669 | `.DataScoringRiskList(1).FEA.HydrantSystem.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 670 | `.DataScoringRiskList(1).FEA.HydrantSystem.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 671 | `.DataScoringRiskList(1).FEA.HydrantSystem.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 672 | `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 673 | `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 674 | `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 675 | `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 676 | `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 677 | `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 678 | `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 679 | `.DataScoringRiskList(1).FEA.PortableFireExtinguisher.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 680 | `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 681 | `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 682 | `.DataScoringRiskList(1).FEA.SecurityGuard.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 683 | `.DataScoringRiskList(1).FEA.SecurityGuard.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 684 | `.DataScoringRiskList(1).FEA.SecurityGuard.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 685 | `.DataScoringRiskList(1).FEA.SecurityGuard.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 686 | `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 687 | `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 688 | `.DataScoringRiskList(1).FEA.SprinklerSystem.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 689 | `.DataScoringRiskList(1).FEA.SprinklerSystem.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 690 | `.DataScoringRiskList(1).FEA.SprinklerSystem.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 691 | `.DataScoringRiskList(1).FEA.SprinklerSystem.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 692 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 693 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox10` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 694 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 695 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 696 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 697 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 698 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 699 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox7` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 700 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox8` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 701 | `.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox9` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 702 | `.DataScoringRiskList(1).LossRatio.LossRatio.PercentLossRatio` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 0/2 | relatif (agregat case) |
| 703 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 704 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score10` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 705 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 706 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 707 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 708 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 709 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score6` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 710 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score7` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 711 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score8` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 712 | `.DataScoringRiskList(1).LossRatio.LossRatio.Score9` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 713 | `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 714 | `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 715 | `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 716 | `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 717 | `.DataScoringRiskList(1).ObjectConditions.Building.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 718 | `.DataScoringRiskList(1).ObjectConditions.Building.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 719 | `.DataScoringRiskList(1).ObjectConditions.Building.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 720 | `.DataScoringRiskList(1).ObjectConditions.Building.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 721 | `.DataScoringRiskList(1).ObjectConditions.Building.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 722 | `.DataScoringRiskList(1).ObjectConditions.Building.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 723 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 724 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 725 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 726 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 727 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 728 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.ChechBox6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 729 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 730 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 731 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 732 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 733 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 734 | `.DataScoringRiskList(1).ObjectConditions.ClassConstruction.Score6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 735 | `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 736 | `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 737 | `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 738 | `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 739 | `.DataScoringRiskList(1).ObjectConditions.Machinery.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 740 | `.DataScoringRiskList(1).ObjectConditions.Machinery.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 741 | `.DataScoringRiskList(1).ObjectConditions.Machinery.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 742 | `.DataScoringRiskList(1).ObjectConditions.Machinery.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 743 | `.DataScoringRiskList(1).ObjectConditions.Machinery.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 744 | `.DataScoringRiskList(1).ObjectConditions.Machinery.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 745 | `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 746 | `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 747 | `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 748 | `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 749 | `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 750 | `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 751 | `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 752 | `.DataScoringRiskList(1).ObjectConditions.Stockinthepolicycover.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 753 | `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 754 | `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 755 | `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 756 | `.DataScoringRiskList(1).ObjectConditions.TypeofStock.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 757 | `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 758 | `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 759 | `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 760 | `.DataScoringRiskList(1).ObjectConditions.TypeofStock.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 761 | `.DataScoringRiskList(1).Occupation.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 762 | `.DataScoringRiskList(1).Occupation.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 763 | `.DataScoringRiskList(1).Occupation.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 764 | `.DataScoringRiskList(1).Occupation.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 765 | `.DataScoringRiskList(1).Occupation.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 766 | `.DataScoringRiskList(1).Occupation.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 767 | `.DataScoringRiskList(1).Occupation.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 768 | `.DataScoringRiskList(1).Occupation.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 769 | `.DataScoringRiskList(1).Occupation.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 770 | `.DataScoringRiskList(1).Occupation.Score5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 771 | `.DataScoringRiskList(1).Others.Clauses.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 772 | `.DataScoringRiskList(1).Others.Clauses.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 773 | `.DataScoringRiskList(1).Others.Clauses.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 774 | `.DataScoringRiskList(1).Others.Clauses.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 775 | `.DataScoringRiskList(1).Others.Clauses.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 776 | `.DataScoringRiskList(1).Others.Clauses.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 777 | `.DataScoringRiskList(1).Others.Deductible.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 778 | `.DataScoringRiskList(1).Others.Deductible.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 779 | `.DataScoringRiskList(1).Others.Deductible.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 780 | `.DataScoringRiskList(1).Others.Deductible.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 781 | `.DataScoringRiskList(1).Others.Deductible.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 782 | `.DataScoringRiskList(1).Others.Deductible.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 783 | `.DataScoringRiskList(1).Others.OtherInformation.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 784 | `.DataScoringRiskList(1).Others.OtherInformation.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 785 | `.DataScoringRiskList(1).Others.OtherInformation.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 786 | `.DataScoringRiskList(1).Others.OtherInformation.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 787 | `.DataScoringRiskList(1).Others.OtherInformation.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 788 | `.DataScoringRiskList(1).Others.OtherInformation.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 789 | `.DataScoringRiskList(1).Others.Rate.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 790 | `.DataScoringRiskList(1).Others.Rate.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 791 | `.DataScoringRiskList(1).Others.Rate.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 792 | `.DataScoringRiskList(1).Others.Rate.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 793 | `.DataScoringRiskList(1).Others.Rate.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 794 | `.DataScoringRiskList(1).Others.Rate.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 795 | `.DataScoringRiskList(1).Others.TotalSumInsured.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 796 | `.DataScoringRiskList(1).Others.TotalSumInsured.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 797 | `.DataScoringRiskList(1).Others.TotalSumInsured.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 798 | `.DataScoringRiskList(1).Others.TotalSumInsured.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 799 | `.DataScoringRiskList(1).RiskImprovement.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 800 | `.DataScoringRiskList(1).RiskImprovement.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 801 | `.DataScoringRiskList(1).RiskImprovement.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 802 | `.DataScoringRiskList(1).RiskImprovement.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 803 | `.DataScoringRiskList(1).RiskImprovement.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 804 | `.DataScoringRiskList(1).RiskImprovement.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 805 | `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 806 | `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 807 | `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 808 | `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 809 | `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 810 | `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 811 | `.DataScoringRiskList(1).RiskLocation.Earthquake.ChechBox7` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 812 | `.DataScoringRiskList(1).RiskLocation.Earthquake.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 813 | `.DataScoringRiskList(1).RiskLocation.Earthquake.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 814 | `.DataScoringRiskList(1).RiskLocation.Earthquake.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 815 | `.DataScoringRiskList(1).RiskLocation.Earthquake.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 816 | `.DataScoringRiskList(1).RiskLocation.Earthquake.Score5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 817 | `.DataScoringRiskList(1).RiskLocation.Earthquake.Score6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 818 | `.DataScoringRiskList(1).RiskLocation.Earthquake.Score7` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 819 | `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 820 | `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 821 | `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 822 | `.DataScoringRiskList(1).RiskLocation.Flood.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 823 | `.DataScoringRiskList(1).RiskLocation.Flood.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 824 | `.DataScoringRiskList(1).RiskLocation.Flood.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 825 | `.DataScoringRiskList(1).RiskLocation.Flood.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 826 | `.DataScoringRiskList(1).RiskLocation.Flood.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 827 | `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 828 | `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 829 | `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 830 | `.DataScoringRiskList(1).RiskLocation.Riot.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 831 | `.DataScoringRiskList(1).RiskLocation.Riot.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 832 | `.DataScoringRiskList(1).RiskLocation.Riot.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 833 | `.DataScoringRiskList(1).RiskLocation.Riot.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 834 | `.DataScoringRiskList(1).RiskLocation.Riot.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 835 | `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 836 | `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 837 | `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 838 | `.DataScoringRiskList(1).RiskLocation.Tsunami.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 839 | `.DataScoringRiskList(1).RiskLocation.Tsunami.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 840 | `.DataScoringRiskList(1).RiskLocation.Tsunami.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 841 | `.DataScoringRiskList(1).RiskLocation.Tsunami.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 842 | `.DataScoringRiskList(1).RiskLocation.Tsunami.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 843 | `.DataScoringRiskList(1).Score` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 844 | `.DataScoringRiskList(1).SurveyReport.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 845 | `.DataScoringRiskList(1).SurveyReport.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 846 | `.DataScoringRiskList(1).SurveyReport.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 847 | `.DataScoringRiskList(1).SurveyReport.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 848 | `.DataScoringRiskList(1).SurveyReport.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 849 | `.DataScoringRiskList(1).SurveyReport.ChechBox6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 850 | `.DataScoringRiskList(1).SurveyReport.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 851 | `.DataScoringRiskList(1).SurveyReport.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 852 | `.DataScoringRiskList(1).SurveyReport.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 853 | `.DataScoringRiskList(1).SurveyReport.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 854 | `.DataScoringRiskList(1).SurveyReport.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 855 | `.DataScoringRiskList(1).SurveyReport.Score6` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 856 | `.DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 857 | `.DataScoringRiskList(2).FEA.FireAlarmSystem.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 858 | `.DataScoringRiskList(2).FEA.FireAlarmSystem.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 859 | `.DataScoringRiskList(2).FEA.FireAlarmSystem.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 860 | `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 861 | `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 862 | `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 863 | `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 864 | `.DataScoringRiskList(2).FEA.Firebrigade.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 865 | `.DataScoringRiskList(2).FEA.Firebrigade.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 866 | `.DataScoringRiskList(2).FEA.Firebrigade.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 867 | `.DataScoringRiskList(2).FEA.Firebrigade.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 868 | `.DataScoringRiskList(2).FEA.Firebrigade.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 869 | `.DataScoringRiskList(2).FEA.Firebrigade.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 870 | `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 871 | `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 872 | `.DataScoringRiskList(2).FEA.HydrantSystem.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 873 | `.DataScoringRiskList(2).FEA.HydrantSystem.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 874 | `.DataScoringRiskList(2).FEA.HydrantSystem.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 875 | `.DataScoringRiskList(2).FEA.HydrantSystem.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 876 | `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 877 | `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 878 | `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 879 | `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 880 | `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 881 | `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 882 | `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 883 | `.DataScoringRiskList(2).FEA.PortableFireExtinguisher.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 884 | `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 885 | `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 886 | `.DataScoringRiskList(2).FEA.SecurityGuard.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 887 | `.DataScoringRiskList(2).FEA.SecurityGuard.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 888 | `.DataScoringRiskList(2).FEA.SecurityGuard.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 889 | `.DataScoringRiskList(2).FEA.SecurityGuard.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 890 | `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 891 | `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 892 | `.DataScoringRiskList(2).FEA.SprinklerSystem.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 893 | `.DataScoringRiskList(2).FEA.SprinklerSystem.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 894 | `.DataScoringRiskList(2).FEA.SprinklerSystem.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 895 | `.DataScoringRiskList(2).FEA.SprinklerSystem.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 896 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 897 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox10` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 898 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 899 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 900 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 901 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 902 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 903 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox7` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 904 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox8` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 905 | `.DataScoringRiskList(2).LossRatio.LossRatio.ChechBox9` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 906 | `.DataScoringRiskList(2).LossRatio.LossRatio.PercentLossRatio` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 0/2 | relatif (agregat case) |
| 907 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 908 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score10` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 909 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 910 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 911 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 912 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 913 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score6` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 914 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score7` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 915 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score8` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 916 | `.DataScoringRiskList(2).LossRatio.LossRatio.Score9` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 917 | `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 918 | `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 919 | `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 920 | `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 921 | `.DataScoringRiskList(2).ObjectConditions.Building.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 922 | `.DataScoringRiskList(2).ObjectConditions.Building.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 923 | `.DataScoringRiskList(2).ObjectConditions.Building.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 924 | `.DataScoringRiskList(2).ObjectConditions.Building.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 925 | `.DataScoringRiskList(2).ObjectConditions.Building.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 926 | `.DataScoringRiskList(2).ObjectConditions.Building.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 927 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 928 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 929 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 930 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 931 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 932 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.ChechBox6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 933 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 934 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 935 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 936 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 937 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 938 | `.DataScoringRiskList(2).ObjectConditions.ClassConstruction.Score6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 939 | `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 940 | `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 941 | `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 942 | `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 943 | `.DataScoringRiskList(2).ObjectConditions.Machinery.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 944 | `.DataScoringRiskList(2).ObjectConditions.Machinery.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 945 | `.DataScoringRiskList(2).ObjectConditions.Machinery.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 946 | `.DataScoringRiskList(2).ObjectConditions.Machinery.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 947 | `.DataScoringRiskList(2).ObjectConditions.Machinery.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 948 | `.DataScoringRiskList(2).ObjectConditions.Machinery.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 949 | `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 950 | `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 951 | `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 952 | `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 953 | `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 954 | `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 955 | `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 956 | `.DataScoringRiskList(2).ObjectConditions.Stockinthepolicycover.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 957 | `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 958 | `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 959 | `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 960 | `.DataScoringRiskList(2).ObjectConditions.TypeofStock.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 961 | `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 962 | `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 963 | `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 964 | `.DataScoringRiskList(2).ObjectConditions.TypeofStock.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 965 | `.DataScoringRiskList(2).Occupation.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 966 | `.DataScoringRiskList(2).Occupation.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 967 | `.DataScoringRiskList(2).Occupation.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 968 | `.DataScoringRiskList(2).Occupation.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 969 | `.DataScoringRiskList(2).Occupation.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 970 | `.DataScoringRiskList(2).Occupation.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 971 | `.DataScoringRiskList(2).Occupation.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 972 | `.DataScoringRiskList(2).Occupation.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 973 | `.DataScoringRiskList(2).Occupation.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 974 | `.DataScoringRiskList(2).Occupation.Score5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 975 | `.DataScoringRiskList(2).Others.Clauses.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 976 | `.DataScoringRiskList(2).Others.Clauses.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 977 | `.DataScoringRiskList(2).Others.Clauses.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 978 | `.DataScoringRiskList(2).Others.Clauses.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 979 | `.DataScoringRiskList(2).Others.Clauses.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 980 | `.DataScoringRiskList(2).Others.Clauses.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 981 | `.DataScoringRiskList(2).Others.Deductible.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 982 | `.DataScoringRiskList(2).Others.Deductible.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 983 | `.DataScoringRiskList(2).Others.Deductible.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 984 | `.DataScoringRiskList(2).Others.Deductible.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 985 | `.DataScoringRiskList(2).Others.Deductible.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 986 | `.DataScoringRiskList(2).Others.Deductible.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 987 | `.DataScoringRiskList(2).Others.OtherInformation.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 988 | `.DataScoringRiskList(2).Others.OtherInformation.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 989 | `.DataScoringRiskList(2).Others.OtherInformation.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 990 | `.DataScoringRiskList(2).Others.OtherInformation.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 991 | `.DataScoringRiskList(2).Others.OtherInformation.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 992 | `.DataScoringRiskList(2).Others.OtherInformation.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 993 | `.DataScoringRiskList(2).Others.Rate.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 994 | `.DataScoringRiskList(2).Others.Rate.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 995 | `.DataScoringRiskList(2).Others.Rate.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 996 | `.DataScoringRiskList(2).Others.Rate.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 997 | `.DataScoringRiskList(2).Others.Rate.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 998 | `.DataScoringRiskList(2).Others.Rate.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 999 | `.DataScoringRiskList(2).Others.TotalSumInsured.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1000 | `.DataScoringRiskList(2).Others.TotalSumInsured.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1001 | `.DataScoringRiskList(2).Others.TotalSumInsured.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1002 | `.DataScoringRiskList(2).Others.TotalSumInsured.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1003 | `.DataScoringRiskList(2).RiskImprovement.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1004 | `.DataScoringRiskList(2).RiskImprovement.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1005 | `.DataScoringRiskList(2).RiskImprovement.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1006 | `.DataScoringRiskList(2).RiskImprovement.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1007 | `.DataScoringRiskList(2).RiskImprovement.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1008 | `.DataScoringRiskList(2).RiskImprovement.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1009 | `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1010 | `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1011 | `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1012 | `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1013 | `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1014 | `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1015 | `.DataScoringRiskList(2).RiskLocation.Earthquake.ChechBox7` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1016 | `.DataScoringRiskList(2).RiskLocation.Earthquake.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1017 | `.DataScoringRiskList(2).RiskLocation.Earthquake.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1018 | `.DataScoringRiskList(2).RiskLocation.Earthquake.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1019 | `.DataScoringRiskList(2).RiskLocation.Earthquake.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1020 | `.DataScoringRiskList(2).RiskLocation.Earthquake.Score5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1021 | `.DataScoringRiskList(2).RiskLocation.Earthquake.Score6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1022 | `.DataScoringRiskList(2).RiskLocation.Earthquake.Score7` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1023 | `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1024 | `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1025 | `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1026 | `.DataScoringRiskList(2).RiskLocation.Flood.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1027 | `.DataScoringRiskList(2).RiskLocation.Flood.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1028 | `.DataScoringRiskList(2).RiskLocation.Flood.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1029 | `.DataScoringRiskList(2).RiskLocation.Flood.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1030 | `.DataScoringRiskList(2).RiskLocation.Flood.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1031 | `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1032 | `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1033 | `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1034 | `.DataScoringRiskList(2).RiskLocation.Riot.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1035 | `.DataScoringRiskList(2).RiskLocation.Riot.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1036 | `.DataScoringRiskList(2).RiskLocation.Riot.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1037 | `.DataScoringRiskList(2).RiskLocation.Riot.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1038 | `.DataScoringRiskList(2).RiskLocation.Riot.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1039 | `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1040 | `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1041 | `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1042 | `.DataScoringRiskList(2).RiskLocation.Tsunami.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1043 | `.DataScoringRiskList(2).RiskLocation.Tsunami.Score1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1044 | `.DataScoringRiskList(2).RiskLocation.Tsunami.Score2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1045 | `.DataScoringRiskList(2).RiskLocation.Tsunami.Score3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1046 | `.DataScoringRiskList(2).RiskLocation.Tsunami.Score4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1047 | `.DataScoringRiskList(2).Score` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1048 | `.DataScoringRiskList(2).SurveyReport.ChechBox1` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1049 | `.DataScoringRiskList(2).SurveyReport.ChechBox2` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1050 | `.DataScoringRiskList(2).SurveyReport.ChechBox3` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1051 | `.DataScoringRiskList(2).SurveyReport.ChechBox4` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1052 | `.DataScoringRiskList(2).SurveyReport.ChechBox5` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1053 | `.DataScoringRiskList(2).SurveyReport.ChechBox6` | 2 | 2 | ASM-FW-GISFW-Data-ScoringRisk | pxCheckbox | 0/2 | relatif (agregat case) |
| 1054 | `.DataScoringRiskList(2).SurveyReport.Score1` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1055 | `.DataScoringRiskList(2).SurveyReport.Score2` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1056 | `.DataScoringRiskList(2).SurveyReport.Score3` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1057 | `.DataScoringRiskList(2).SurveyReport.Score4` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1058 | `.DataScoringRiskList(2).SurveyReport.Score5` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1059 | `.DataScoringRiskList(2).SurveyReport.Score6` | 2 | 2 | (kosong) | pxNumber | 2/2 | relatif (agregat case) |
| 1060 | `.DateOfBirth` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxDateTime | 2/2 | relatif (agregat case) |
| 1061 | `.DateSuggest` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | pxDateTime | 2/2 | relatif (agregat case) |
| 1062 | `.DeductibleNote` | 2 | 2 | ASM-FW-GISFW-Data-Deductible | pxTextArea | /2 | relatif (agregat case) |
| 1063 | `.DepositAccountNo` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 2/2 | relatif (agregat case) |
| 1064 | `.DOL` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | /2 | relatif (agregat case) |
| 1065 | `.EmailTypeUWPolicy` | 2 | 2 | ASM-FW-GISFW-Work | pxRadioButtons | 0/2 | relatif (agregat case) |
| 1066 | `.END_DATE` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDateTime/pxDisplayText | 2/2 | relatif (agregat case) |
| 1067 | `.EndDateTime` | 2 | 2 | ASM-FW-GISFW-Data-Policy | pxDateTime | 2/2 | relatif (agregat case) |
| 1068 | `.FACULTATIVE_TYPE` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 1069 | `.FilterFollowingNB` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 0/2 | relatif (agregat case) |
| 1070 | `.FlagName` | 2 | 2 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | 2/2 | relatif (agregat case) |
| 1071 | `.FlagRetroTreaty` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxCheckbox | 0/2 | relatif (agregat case) |
| 1072 | `.Following` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 2/2 | relatif (agregat case) |
| 1073 | `.GrossClaim` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxNumber | 0/2 | relatif (agregat case) |
| 1074 | `.IDCardNo` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 2/2 | relatif (agregat case) |
| 1075 | `.IDCurrency` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDropdown | 0/2 | relatif (agregat case) |
| 1076 | `.Info` | 2 | 2 | ASM-FW-GISFW-Int-CLAUSE | pxTextInput | 0/2 | relatif (agregat case) |
| 1077 | `.Institution` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 2/2 | relatif (agregat case) |
| 1078 | `.INSURED_NAME` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 1079 | `.InwardScale.PctTreatyLimit` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | 2/2 | relatif (agregat case) |
| 1080 | `.IS_SELECTED` | 2 | 2 | ASM-FW-GISFW-Int-SUBCONTRACT | pxCheckbox | 0/2 | relatif (agregat case) |
| 1081 | `.IsApproved` | 2 | 2 | ASM-FW-GISFW-Data-SuggestList | pxDropdown/pxRadioButtons | /2 | relatif (agregat case) |
| 1082 | `.IsGrossNetShow` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | 0/2 | relatif (agregat case) |
| 1083 | `.IsProRate` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | /2 | relatif (agregat case) |
| 1084 | `.IsSelected` | 2 | 2 | ASM-FW-GISFW-Int-CLAUSE | pxCheckbox | 0/2 | relatif (agregat case) |
| 1085 | `.KodeSection` | 2 | 2 | ASM-FW-GISFW-Data-Clause | (tanpa) | 2/2 | relatif (agregat case) |
| 1086 | `.LoadingAmount` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxDisplayText | 2/2 | relatif (agregat case) |
| 1087 | `.LoanAccountNumber` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 2/2 | relatif (agregat case) |
| 1088 | `.LoanName` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 2/2 | relatif (agregat case) |
| 1089 | `.LoanType` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxTextInput | 2/2 | relatif (agregat case) |
| 1090 | `.LocationOfTopRisk` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxAutoComplete/pxDisplayText | /2 | relatif (agregat case) |
| 1091 | `.MARKETING_NAME` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 1092 | `.MarketingOfficer` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDisplayText | 2/2 | relatif (agregat case) |
| 1093 | `.Message` | 2 | 2 | ASM-FW-GISFW-Data-FacLetter | pxRichTextEditor/pxTextArea | /2 | relatif (agregat case) |
| 1094 | `.NationName` | 2 | 2 | ASM-FW-GISFW-Int-RISKADDRESS | pxDisplayText | /2 | relatif (agregat case) |
| 1095 | `.NetRate` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxNumber | 2/2 | relatif (agregat case) |
| 1096 | `.NO_OFFER_SLIP` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 1097 | `.OBJECT_DESCRIPTION` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 1098 | `.OfferFacIn.CedingRetentionNominal` | 2 | 2 | ASM-FW-GISFW-Work | pxNumber | 0/2 | relatif (agregat case) |
| 1099 | `.OfferFacIn.FacRetroDetails.DocumentPosition` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1100 | `.OfferFacIn.FlagSpreading` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1101 | `.OfferFacIn.Obligee.AdministrationFee` | 2 | 2 | ASM-FW-GISFW-Work | pxInteger | 2/2 | relatif (agregat case) |
| 1102 | `.OfferFacIn.Obligee.AgreementDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 1103 | `.OfferFacIn.Obligee.AgreementNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1104 | `.OfferFacIn.Obligee.BankAddress` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 1105 | `.OfferFacIn.Obligee.BankGuaranteeNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1106 | `.OfferFacIn.Obligee.BankName` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1107 | `.OfferFacIn.Obligee.CustomsActivity` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1108 | `.OfferFacIn.Obligee.DateOfGuarantee` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 1109 | `.OfferFacIn.Obligee.FacilityDecreeDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 1110 | `.OfferFacIn.Obligee.FacilityDecreeNumber` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1111 | `.OfferFacIn.Obligee.FacilityStatus` | 2 | 2 | ASM-FW-GISFW-Work | (tanpa) | 2/2 | relatif (agregat case) |
| 1112 | `.OfferFacIn.Obligee.InsuredNumber` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1113 | `.OfferFacIn.Obligee.JenisPenjaminan` | 2 | 2 | ASM-FW-GISFW-Work | pxDropdown | 2/2 | relatif (agregat case) |
| 1114 | `.OfferFacIn.Obligee.KBGDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 1115 | `.OfferFacIn.Obligee.KBGNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1116 | `.OfferFacIn.Obligee.LengthPeriod` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1117 | `.OfferFacIn.Obligee.Location` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 1118 | `.OfferFacIn.Obligee.MaxOfLimitOfLiability` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1119 | `.OfferFacIn.Obligee.NIPER` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1120 | `.OfferFacIn.Obligee.OCDescription` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 1121 | `.OfferFacIn.Obligee.OtherFacilityNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 1122 | `.OfferFacIn.Obligee.Pasal` | 2 | 2 | ASM-FW-GISFW-Work | pxDropdown | 2/2 | relatif (agregat case) |
| 1123 | `.OfferFacIn.Obligee.PIBDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 1124 | `.OfferFacIn.Obligee.PIBNumber` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1125 | `.OfferFacIn.Obligee.PKDate` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 1126 | `.OfferFacIn.Obligee.PKDate2` | 2 | 2 | ASM-FW-GISFW-Work | pxDateTime | 2/2 | relatif (agregat case) |
| 1127 | `.OfferFacIn.Obligee.PKNo1` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1128 | `.OfferFacIn.Obligee.PKNo2` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1129 | `.OfferFacIn.Obligee.PlaceOfEvent` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | 2/2 | relatif (agregat case) |
| 1130 | `.OfferFacIn.Obligee.RegisterNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1131 | `.OfferFacIn.PolicyData.SurveyAgent.pyCity` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1132 | `.OfferFacIn.PolicyData.SurveyAgent.pyCountry` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1133 | `.OfferFacIn.ShareRNML` | 2 | 2 | ASM-FW-GISFW-Work | pxNumber | 2/2 | relatif (agregat case) |
| 1134 | `.OfferFacIn.TotalRIComNusaRe` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1135 | `.Others` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxTextArea | 2/2 | relatif (agregat case) |
| 1136 | `.OutgoAmount` | 2 | 2 | ASM-FW-GISFW-Data-Outgo | pxNumber | 2/2 | relatif (agregat case) |
| 1137 | `.Parameters.LowestPctLimit` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | 2/2 | relatif (agregat case) |
| 1138 | `.PartshipmentNote` | 2 | 2 | ASM-FW-GISFW-Data-Trading | pxTextInput | 2/2 | relatif (agregat case) |
| 1139 | `.PartshipmentStatus` | 2 | 2 | ASM-FW-GISFW-Data-Trading | pxDropdown | 2/2 | relatif (agregat case) |
| 1140 | `.PaymentDate` | 2 | 2 | ASM-FW-GISFW-Data-Payment | pxDateTime | 2/2 | relatif (agregat case) |
| 1141 | `.PctAdjust2` | 2 | 2 | ASM-FW-GISFW-Data-PropertyItem | pxDropdown | 2/2 | relatif (agregat case) |
| 1142 | `.PctAdjustOther` | 2 | 2 | ASM-FW-GISFW-Data-PropertyItem | pxDropdown/pxNumber | /2 | relatif (agregat case) |
| 1143 | `.PctLimit` | 2 | 2 | ASM-FW-GISFW-Int-TABLEOFLIMIT | pxNumber/pxTextInput | /2 | relatif (agregat case) |
| 1144 | `.PercentOutgo` | 2 | 2 | ASM-FW-GISFW-Data-Outgo | pxTextInput | 2/2 | relatif (agregat case) |
| 1145 | `.PlanList(1).Name` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxDisplayText | /2 | relatif (agregat case) |
| 1146 | `.Policy.Payment.Installment` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1147 | `.PolicyData.Ship.ID` | 2 | 2 | ASM-FW-GISFW-Data-Cargo | pxAutoComplete | 2/2 | relatif (agregat case) |
| 1148 | `.PolicyData.SurveyAgent.pyCity` | 2 | 2 | ASM-FW-GISFW-Data-Cargo | pxTextInput | 2/2 | relatif (agregat case) |
| 1149 | `.PolicyData.SurveyAgent.pyCountry` | 2 | 2 | ASM-FW-GISFW-Data-Cargo | pxTextInput | 2/2 | relatif (agregat case) |
| 1150 | `.PolicyMasterNumber` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 2/2 | relatif (agregat case) |
| 1151 | `.POSTAL_CODE` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 1152 | `.PostalCode` | 2 | 2 | ASM-FW-GISFW-Int-RISKADDRESS | pxDisplayText | /2 | relatif (agregat case) |
| 1153 | `.PremiRpOld` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxNumber | 2/2 | relatif (agregat case) |
| 1154 | `.PreventionOfLoss` | 2 | 2 | ASM-FW-GISFW-Data-CauseOfLoss | pxNumber/pxTextInput | 2/2 | relatif (agregat case) |
| 1155 | `.PrintRISlip.FACRETROSLIPNUMBER` | 2 | 2 | ASM-FW-GISFW-Data-FacOffer | (tanpa) | 2/2 | relatif (agregat case) |
| 1156 | `.ProductBriguna` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | /2 | relatif (agregat case) |
| 1157 | `.Property.BuildingConstruction.FloorType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1158 | `.Property.BuildingConstruction.NumberOfFloor` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | 2/2 | relatif (agregat case) |
| 1159 | `.Property.BuildingConstruction.OthersType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 2/2 | relatif (agregat case) |
| 1160 | `.Property.BuildingConstruction.PartitionType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 2/2 | relatif (agregat case) |
| 1161 | `.Property.BuildingConstruction.RoofType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1162 | `.Property.BuildingConstruction.SupportWallType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextInput | 2/2 | relatif (agregat case) |
| 1163 | `.Property.BuildingConstruction.WallType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1164 | `.Property.IsFlammableItemFlag` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxCheckbox | /2 | relatif (agregat case) |
| 1165 | `.Property.IsHotWorkProcessFlag` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxCheckbox | /2 | relatif (agregat case) |
| 1166 | `.Property.IsMaterialDamage` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxCheckbox | 2/2 | relatif (agregat case) |
| 1167 | `.Property.IsProductionProcessFlag` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxCheckbox | /2 | relatif (agregat case) |
| 1168 | `.Property.ObjectType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1169 | `.Property.Ownership` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1170 | `.Property.SurroundingRisk.BackConstruction` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1171 | `.Property.SurroundingRisk.BackDistance` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | 2/2 | relatif (agregat case) |
| 1172 | `.Property.SurroundingRisk.BackNote` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextArea | 2/2 | relatif (agregat case) |
| 1173 | `.Property.SurroundingRisk.BackOccupation` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxAutoComplete | 2/2 | relatif (agregat case) |
| 1174 | `.Property.SurroundingRisk.FloodArea` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1175 | `.Property.SurroundingRisk.FloodAreaStatus` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1176 | `.Property.SurroundingRisk.FrontConstruction` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1177 | `.Property.SurroundingRisk.FrontDistance` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | 2/2 | relatif (agregat case) |
| 1178 | `.Property.SurroundingRisk.FrontNote` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextArea | 2/2 | relatif (agregat case) |
| 1179 | `.Property.SurroundingRisk.FrontOccupation` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxAutoComplete | 2/2 | relatif (agregat case) |
| 1180 | `.Property.SurroundingRisk.HousekeepingRemark` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextArea | 2/2 | relatif (agregat case) |
| 1181 | `.Property.SurroundingRisk.HousekeepingStatus` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1182 | `.Property.SurroundingRisk.LeftConstruction` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1183 | `.Property.SurroundingRisk.LeftDistance` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | 2/2 | relatif (agregat case) |
| 1184 | `.Property.SurroundingRisk.LeftNote` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextArea | 2/2 | relatif (agregat case) |
| 1185 | `.Property.SurroundingRisk.LeftOccupation` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxAutoComplete | 2/2 | relatif (agregat case) |
| 1186 | `.Property.SurroundingRisk.RightConstruction` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxDropdown | 2/2 | relatif (agregat case) |
| 1187 | `.Property.SurroundingRisk.RightDistance` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | 2/2 | relatif (agregat case) |
| 1188 | `.Property.SurroundingRisk.RightNote` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxTextArea | 2/2 | relatif (agregat case) |
| 1189 | `.Property.SurroundingRisk.RightOccupation` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxAutoComplete | 2/2 | relatif (agregat case) |
| 1190 | `.ProRateType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | 0/2 | relatif (agregat case) |
| 1191 | `.ProvinceName` | 2 | 2 | ASM-FW-GISFW-Int-RISKADDRESS | pxDisplayText | /2 | relatif (agregat case) |
| 1192 | `.pyAddress` | 2 | 2 | Data-Party-Person | (tanpa) | /2 | relatif (agregat case) |
| 1193 | `.pyCategory` | 2 | 2 | Embed-DragDropFile | (tanpa) | 0/2 | relatif (agregat case) |
| 1194 | `.pyFileName` | 2 | 2 | Embed-DragDropFile | pxLink | 0/2 | relatif (agregat case) |
| 1195 | `.pyID` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | relatif (agregat case) |
| 1196 | `.pyMemo` | 2 | 2 | Link-Attachment | ReasAttachmentDescription_name | 0/2 | relatif (agregat case) |
| 1197 | `.pyOrg` | 2 | 2 | Data-WorkAttach-Note | pxDisplayText | 2/2 | relatif (agregat case) |
| 1198 | `.QuotationData.CommentOldPolicy` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxTextArea | /2 | relatif (agregat case) |
| 1199 | `.QuotationData.EDMDay` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxRadioButtons | 0/2 | relatif (agregat case) |
| 1200 | `.QuotationData.EdmType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | 2/2 | relatif (agregat case) |
| 1201 | `.QuotationData.IsMOP` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxRadioButtons | 2/2 | relatif (agregat case) |
| 1202 | `.QuotationData.IsTemplate` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | 2/2 | relatif (agregat case) |
| 1203 | `.QuotationData.MarketingName` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 2/2 | relatif (agregat case) |
| 1204 | `.QuotationData.NPLStatus` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | 0/2 | relatif (agregat case) |
| 1205 | `.QuotationData.OldPolicyNo` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText/pxTextInput | 2/2 | relatif (agregat case) |
| 1206 | `.QuotationData.PeriodMM` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 0/2 | relatif (agregat case) |
| 1207 | `.QuotationData.PolicyType` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxRadioButtons | /2 | relatif (agregat case) |
| 1208 | `.QuotationData.TypeFacultative` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown/pxTextInput | /2 | relatif (agregat case) |
| 1209 | `.RateEM` | 2 | 2 | ASM-FW-GISFW-Data-Installment | pxTextInput | 2/2 | relatif (agregat case) |
| 1210 | `.RateLife` | 2 | 2 | ASM-FW-GISFW-Data-Installment | pxTextInput | 2/2 | relatif (agregat case) |
| 1211 | `.RateOJK` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxNumber | 2/2 | relatif (agregat case) |
| 1212 | `.ReceiverOutgo` | 2 | 2 | ASM-FW-GISFW-Data-Outgo | pxAutoComplete | 2/2 | relatif (agregat case) |
| 1213 | `.Recipient` | 2 | 2 | ASM-FW-GISFW-Data-Correspondence | pxAutoComplete | 0/2 | relatif (agregat case) |
| 1214 | `.REGISTER_ID` | 2 | 2 | ASM-FW-GISFW-Int-MSHIP | pxAutoComplete | 2/2 | relatif (agregat case) |
| 1215 | `.RICommision` | 2 | 2 | ASM-FW-GISFW-Data-Installment | pxTextInput | 2/2 | relatif (agregat case) |
| 1216 | `.Salary` | 2 | 2 | ASM-FW-GISFW-Data-Aneka | pxNumber | 2/2 | relatif (agregat case) |
| 1217 | `.Sender` | 2 | 2 | ASM-FW-GISFW-Data-Correspondence | pxTextInput | 0/2 | relatif (agregat case) |
| 1218 | `.SOB_NAME` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | 2/2 | relatif (agregat case) |
| 1219 | `.SOBName` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDisplayText | 2/2 | relatif (agregat case) |
| 1220 | `.START_DATE` | 2 | 2 | ASM-FW-GISFW-Int-TABLE_B2B | pxDateTime/pxDisplayText | 2/2 | relatif (agregat case) |
| 1221 | `.StartDateTime` | 2 | 2 | ASM-FW-GISFW-Data-Policy | pxDateTime | 2/2 | relatif (agregat case) |
| 1222 | `.Subject` | 2 | 2 | ASM-FW-GISFW-Data-Correspondence | pxTextInput | 0/2 | relatif (agregat case) |
| 1223 | `.Suggest` | 2 | 2 | ASM-FW-GISFW-Data-SuggestList | pxDisplayText/pxTextArea | /2 | relatif (agregat case) |
| 1224 | `.TERMCONDITION` | 2 | 2 | ASM-FW-GISFW-Data-CauseOfDecline | pxTextInput | 2/2 | relatif (agregat case) |
| 1225 | `.Title` | 2 | 2 | ASM-FW-GISFW-Int-RISKADDRESS | pxDisplayText/pxTextInput | /2 | relatif (agregat case) |
| 1226 | `.TradingID` | 2 | 2 | ASM-FW-GISFW-Data-Trading | pxDropdown | 2/2 | relatif (agregat case) |
| 1227 | `.TranshipmentNote` | 2 | 2 | ASM-FW-GISFW-Data-Trading | pxTextInput | 2/2 | relatif (agregat case) |
| 1228 | `.TranshipmentStatus` | 2 | 2 | ASM-FW-GISFW-Data-Trading | pxDropdown | 2/2 | relatif (agregat case) |
| 1229 | `.TSICeding` | 2 | 2 | ASM-FW-GISFW-Data-Coverage | pxNumber | 2/2 | relatif (agregat case) |
| 1230 | `.TSIGrossSpreaded` | 2 | 2 | ASM-FW-GISFW-Data-SpreadingRisk | pxNumber | 2/2 | relatif (agregat case) |
| 1231 | `.TypeOfNominal` | 2 | 2 | ASM-FW-GISFW-Data-Outgo | pxCheckbox | /2 | relatif (agregat case) |
| 1232 | `.TypeOfOutgo` | 2 | 2 | ASM-FW-GISFW-Data-Outgo | pxDropdown | 2/2 | relatif (agregat case) |
| 1233 | `.UWYEAR` | 2 | 2 | ASM-FW-GISFW-Data-CauseOfLoss | pxTextInput | /2 | relatif (agregat case) |
| 1234 | `.WarrantyID` | 2 | 2 | ASM-FW-GISFW-Data-Warranty | pxDisplayText | 2/2 | relatif (agregat case) |
| 1235 | `.YearMake1` | 2 | 2 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | 2/2 | relatif (agregat case) |
| 1236 | `A.NBStatusNew` | 2 | 2 | ASM-FW-SFAGISFW-Work-Opportunity | pxDisplayText | 2/2 | LUAR agregat: A |
| 1237 | `FacOfferList.AdditionalInfo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextArea | /2 | LUAR agregat: FacOfferList |
| 1238 | `FacOfferList.FacNo` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | LUAR agregat: FacOfferList |
| 1239 | `InputAccumulatedType.Type` | 2 | 2 | @baseclass | pxDisplayText/pxHidden | /2 | LUAR agregat: InputAccumulatedType |
| 1240 | `InputData.CARI7` | 2 | 2 | ASM-FW-GISFW-Work | pxTextInput | 2/2 | LUAR agregat: InputData |
| 1241 | `InputRiskAddress.CityName` | 2 | 2 | ASM-FW-GISFW-Work | pxAutoComplete | 0/2 | LUAR agregat: InputRiskAddress |
| 1242 | `InputRiskAddress.DistrictName` | 2 | 2 | ASM-FW-GISFW-Work | pxAutoComplete | 0/2 | LUAR agregat: InputRiskAddress |
| 1243 | `InputRiskAddress.ProvinceName` | 2 | 2 | ASM-FW-GISFW-Work | pxAutoComplete | 0/2 | LUAR agregat: InputRiskAddress |
| 1244 | `pyWorkPage.Policy.Currency` | 2 | 2 | ASM-FW-GISFW-Data-PropertyItem | (tanpa) | 2/2 | agregat case (eksplisit) |
| 1245 | `pyWorkPage.PolicyTreatyIn.Installment` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 2/2 | agregat case (eksplisit) |
| 1246 | `pyWorkPage.PolicyTreatyIn.ProductionDate` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDateTime | 2/2 | agregat case (eksplisit) |
| 1247 | `pyWorkPage.PolicyTreatyIn.SOBName` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 2/2 | agregat case (eksplisit) |
| 1248 | `pyWorkPage.PPH_Note` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-Currency | (tanpa) | 2/2 | agregat case (eksplisit) |
| 1249 | `pyWorkPage.PPN_Note` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn-Currency | (tanpa) | 2/2 | agregat case (eksplisit) |
| 1250 | `pyWorkPage.pyID` | 2 | 2 | ASM-FW-GISFW-Work | (tanpa) | 2/2 | agregat case (eksplisit) |
| 1251 | `pyWorkPage.TreatyIn.RNMShare` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 2/2 | agregat case (eksplisit) |
| 1252 | `pyWorkPage.TreatyIn.RnmShareDeducted` | 2 | 2 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | 2/2 | agregat case (eksplisit) |
| 1253 | `RISlipContent.CARI1` | 2 | 2 | ASM-FW-GISFW-Data-FacOffer | pxRichTextEditor | 2/2 | LUAR agregat: RISlipContent |
| 1254 | `RISlipContent.CARI3` | 2 | 2 | ASM-FW-GISFW-Data-FacOffer | pxRichTextEditor | 2/2 | LUAR agregat: RISlipContent |
| 1255 | `RISlipContent.CARI4` | 2 | 2 | ASM-FW-GISFW-Data-FacOffer | pxRichTextEditor | 2/2 | LUAR agregat: RISlipContent |
| 1256 | `SearchLife.CARI7` | 2 | 2 | ASM-FW-GISFW-Work | pxAutoComplete | 0/2 | LUAR agregat: SearchLife |
| 1257 | `TempError.ERRMSG` | 2 | 2 | ASM-FW-GISFW-Data-OfferFacIn | pxTextArea | 2/2 | LUAR agregat: TempError |
| 1258 | `TempPolis.CARI16` | 2 | 2 | ASM-FW-GISFW-Work | pxDisplayText | 2/2 | LUAR agregat: TempPolis |
| 1259 | `ViewOfferStatus.LossProtection` | 2 | 2 | ASM-FW-GISFW-Work | pxCheckbox | 0/2 | LUAR agregat: ViewOfferStatus |
| 1260 | `.ACCUMULATION_CODE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1261 | `.ACCUMULATION_NAME` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1262 | `.AdminFee` | 1 | 1 | ASM-FW-GISFW-Data-Installment | pxNumber | /1 | relatif (agregat case) |
| 1263 | `.AgingAmount` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | /1 | relatif (agregat case) |
| 1264 | `.AkseptasiKlaim` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfLossClaim | pxNumber | /1 | relatif (agregat case) |
| 1265 | `.AMOUNT1` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1266 | `.AMOUNT2` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1267 | `.AMOUNT3` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1268 | `.AMOUNT4` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1269 | `.AMOUNT5` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1270 | `.AmountTotal` | 1 | 1 | ASM-FW-GISFW-Data-TreatyInInstallment | pxNumber | 0/1 | relatif (agregat case) |
| 1271 | `.AnekaAnotherSchedule.OtherAnnualPremiumDesc` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1272 | `.AnekaAnotherSchedule.RetroactiveDate` | 1 | 1 | ASM-FW-GISFW-Work | pxDateTime | /1 | relatif (agregat case) |
| 1273 | `.AnekaAnotherSchedule.TerritorialLimit` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1274 | `.AnekaAnotherSchedule.Yurisdiction` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1275 | `.AnekaMaintenance.EndPeriod` | 1 | 1 | Data-Address | pxDateTime | /1 | relatif (agregat case) |
| 1276 | `.AnekaMaintenance.ExtendedPeriodDay` | 1 | 1 | Data-Address | pxTextInput | /1 | relatif (agregat case) |
| 1277 | `.AnekaMaintenance.ExtendedPeriodMonth` | 1 | 1 | Data-Address | pxTextInput | /1 | relatif (agregat case) |
| 1278 | `.AnekaMaintenance.ExtendedPeriodWeek` | 1 | 1 | Data-Address | pxTextInput | /1 | relatif (agregat case) |
| 1279 | `.AnekaMaintenance.MaintenancePeriodDay` | 1 | 1 | Data-Address | pxTextInput | /1 | relatif (agregat case) |
| 1280 | `.AnekaMaintenance.MaintenancePeriodMonth` | 1 | 1 | Data-Address | pxTextInput | /1 | relatif (agregat case) |
| 1281 | `.AnekaMaintenance.MaintenancePeriodWeek` | 1 | 1 | Data-Address | pxTextInput | /1 | relatif (agregat case) |
| 1282 | `.AnekaMaintenance.StartPeriod` | 1 | 1 | Data-Address | pxDateTime | /1 | relatif (agregat case) |
| 1283 | `.AnekaMaintenance.TitleOfContract` | 1 | 1 | Data-Address | pxTextArea | /1 | relatif (agregat case) |
| 1284 | `.AnekaPolicySchedule.DeductibleInsured1` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1285 | `.AnekaPolicySchedule.DeductibleInsured2` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1286 | `.AnekaPolicySchedule.DeductibleRisk` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1287 | `.AnekaPolicySchedule.InsuredItemSection1` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1288 | `.AnekaPolicySchedule.InsuredItemSection2` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1289 | `.AnekaPolicySchedule.LimitOfIndemnitySection1` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1290 | `.AnekaPolicySchedule.LimitOfIndemnitySection2` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1291 | `.AnekaPolicySchedule.Risk` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1292 | `.AnekaPolicySchedule.SumInsured` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1293 | `.Answer` | 1 | 1 | ASM-FW-GISFW-Data-SubjectTo | pxRadioButtons | /1 | relatif (agregat case) |
| 1294 | `.ASMContactInfo` | 1 | 1 | Data-Address | pxTextInput | /1 | relatif (agregat case) |
| 1295 | `.ASMInsuredRelationship` | 1 | 1 | Data-Party-Person | pxDropdown | 0/1 | relatif (agregat case) |
| 1296 | `.ASMJoinDate` | 1 | 1 | Data-Party-Person | pxDateTime | /1 | relatif (agregat case) |
| 1297 | `.ASMPremium` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxNumber | /1 | relatif (agregat case) |
| 1298 | `.ASMRegNo` | 1 | 1 | Data-Party-Person | pxDisplayText | /1 | relatif (agregat case) |
| 1299 | `.ASMTSI` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxNumber | /1 | relatif (agregat case) |
| 1300 | `.BEGIN_DATE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1301 | `.BinderRNM` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | relatif (agregat case) |
| 1302 | `.Branch` | 1 | 1 | ASM-FW-GISFW-Data-AnekaParticipant | pxTextInput | /1 | relatif (agregat case) |
| 1303 | `.Brokerage` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1304 | `.BUILDING_NO` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1305 | `.CALCULATE_METHOD` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1306 | `.CalculateMethod_FacIn` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxDropdown | /1 | relatif (agregat case) |
| 1307 | `.CARI19` | 1 | 1 | ASM-FW-GISFW-Data-Search | (tanpa) | 0/1 | relatif (agregat case) |
| 1308 | `.CARI20` | 1 | 1 | ASM-FW-GISFW-Data-Search | (tanpa) | /1 | relatif (agregat case) |
| 1309 | `.CARI4` | 1 | 1 | ASM-FW-GISFW-Data-Search | pxNumber | /1 | relatif (agregat case) |
| 1310 | `.CARI5` | 1 | 1 | ASM-FW-GISFW-Data-Search | pxNumber | /1 | relatif (agregat case) |
| 1311 | `.CaseID` | 1 | 1 | ASM-FW-GISFW-Int-policyjson | pxTextInput | 0/1 | relatif (agregat case) |
| 1312 | `.Category` | 1 | 1 | ASM-FW-GISFW-Int-TABLEOFLIMIT | pxTextInput | 0/1 | relatif (agregat case) |
| 1313 | `.CERTNO` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1314 | `.CheckBox2` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfDecline | pxCheckbox | 0/1 | relatif (agregat case) |
| 1315 | `.CIFData.Customer_P.pyCity` | 1 | 1 | ASM-FW-GISFW-Data-Policy | pxTextInput | /1 | relatif (agregat case) |
| 1316 | `.CIFData.Customer_P.pyCompany` | 1 | 1 | ASM-FW-GISFW-Data-Policy | pxTextInput | /1 | relatif (agregat case) |
| 1317 | `.City` | 1 | 1 | ASM-FW-GISFW-Int-lloydagent | pxDisplayText | /1 | relatif (agregat case) |
| 1318 | `.CLASS_OF_CONSTRUCTION` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1319 | `.Classification` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1320 | `.ClauseLanguageID` | 1 | 1 | ASM-FW-GISFW-Data-Clause | pxDropdown | /1 | relatif (agregat case) |
| 1321 | `.CoinsList(1).TSIShare` | 1 | 1 | ASM-FW-GISFW-Work | pxNumber | /1 | relatif (agregat case) |
| 1322 | `.CoinsName` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxAutoComplete | /1 | relatif (agregat case) |
| 1323 | `.Comment` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextArea | 0/1 | relatif (agregat case) |
| 1324 | `.CommentDeductible` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextArea | 0/1 | relatif (agregat case) |
| 1325 | `.Commentform` | 1 | 1 | ASM-FW-GISFW-Data-FacOffer | pxTextArea | 0/1 | relatif (agregat case) |
| 1326 | `.COMMISION_AMOUNT` | 1 | 1 | ASM-FW-GISFW-Int-RETROCESSIONLIFE | pxNumber | /1 | relatif (agregat case) |
| 1327 | `.CONDITION_NAME1` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1328 | `.CONDITION_NAME2` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1329 | `.CONDITION_NAME3` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1330 | `.CONDITION_NAME4` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1331 | `.CONDITION_NAME5` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1332 | `.CONDITION1` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1333 | `.CONDITION2` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1334 | `.CONDITION3` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1335 | `.CONDITION4` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1336 | `.CONDITION5` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1337 | `.ConditionName` | 1 | 1 | ASM-FW-GISFW-Int-CONDITION | (tanpa) | 0/1 | relatif (agregat case) |
| 1338 | `.Confirceding` | 1 | 1 | ASM-FW-GISFW-Data-FacOffer | pxDropdown | 0/1 | relatif (agregat case) |
| 1339 | `.CONVEYENCE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1340 | `.CONVEYENCE_NAME` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1341 | `.COVERAGE_BASIS` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1342 | `.COVERAGE_PREMIUM` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1343 | `.COVERAGE_RATE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1344 | `.CoverageList(1).Currency.Name` | 1 | 1 | ASM-FW-GISFW-Data-Cargo | pxNumber | /1 | relatif (agregat case) |
| 1345 | `.CURRENCY1` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1346 | `.CURRENCY2` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1347 | `.CURRENCY3` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1348 | `.CURRENCY4` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1349 | `.CURRENCY5` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1350 | `.CurrencyID` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-Currency | (tanpa) | /1 | relatif (agregat case) |
| 1351 | `.CurrentMode` | 1 | 1 | Data-Portal | pxHidden | 0/1 | relatif (agregat case) |
| 1352 | `.Customer_C.ASMClientID` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1353 | `.Customer_C.ASMDirectorPerson.pyFullName` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1354 | `.Customer_C.ASMNPWP` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1355 | `.Customer_C.pyCompany` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1356 | `.Customer_P.ASMAnotherIncome` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1357 | `.Customer_P.ASMClientID` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1358 | `.Customer_P.ASMCompanyInfo.ASMFoundedIn` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1359 | `.Customer_P.ASMDateOfBirth` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1360 | `.Customer_P.ASMDependentCount` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1361 | `.Customer_P.ASMGender` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1362 | `.Customer_P.ASMIDCard` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1363 | `.Customer_P.ASMJobDesc` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1364 | `.Customer_P.ASMMaritalStatus` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1365 | `.Customer_P.ASMNationality` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1366 | `.Customer_P.ASMNickName` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1367 | `.Customer_P.ASMNPWP` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1368 | `.Customer_P.ASMPosition` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1369 | `.Customer_P.ASMSourceOfIncome` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1370 | `.Customer_P.ASMYearlyIncome` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1371 | `.Customer_P.pyFirstName` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1372 | `.Customer_P.pyLastName` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1373 | `.CustomerName` | 1 | 1 | ASM-FW-GISFW-Int-policyjson | pxTextInput | 0/1 | relatif (agregat case) |
| 1374 | `.DataFEA.PrivateFireBrigade` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDropdown | /1 | relatif (agregat case) |
| 1375 | `.DataFEA.TeamSOPRiskManagement` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDropdown | /1 | relatif (agregat case) |
| 1376 | `.DataFEA.TeamSOPSafety` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | pxDropdown | /1 | relatif (agregat case) |
| 1377 | `.DATE_END` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1378 | `.DateCreated` | 1 | 1 | ASM-FW-GISFW-Data-FacLetter | pxDateTime | 0/1 | relatif (agregat case) |
| 1379 | `.DeductName` | 1 | 1 | ASM-FW-GISFW-Int-DEDUCTIBLE | pxDisplayText | /1 | relatif (agregat case) |
| 1380 | `.DEPARTURE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1381 | `.Description1` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfDecline | pxAutoComplete | 0/1 | relatif (agregat case) |
| 1382 | `.Description2` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfDecline | pxTextInput | 0/1 | relatif (agregat case) |
| 1383 | `.DESTINATION` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1384 | `.DiscountShare` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1385 | `.DistributionStatus` | 1 | 1 | ASM-FW-GISFW-Data-FEA | pxDropdown | /1 | relatif (agregat case) |
| 1386 | `.DistrictName` | 1 | 1 | ASM-FW-GISFW-Int-RISKADDRESS | pxDisplayText | /1 | relatif (agregat case) |
| 1387 | `.DWT` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1388 | `.Email` | 1 | 1 | ASM-FW-GISFW-Data-Quotation | pxTextInput | 0/1 | relatif (agregat case) |
| 1389 | `.EmailText` | 1 | 1 | ASM-FW-GISFW-Work | pxRichTextEditor | 0/1 | relatif (agregat case) |
| 1390 | `.EmailType` | 1 | 1 | ASM-FW-GISFW-Work | pxRadioButtons | 0/1 | relatif (agregat case) |
| 1391 | `.EmailTypeBinding` | 1 | 1 | ASM-FW-GISFW-Work | pxRadioButtons | 0/1 | relatif (agregat case) |
| 1392 | `.EmailTypeCeding` | 1 | 1 | ASM-FW-GISFW-Work | pxRadioButtons | 0/1 | relatif (agregat case) |
| 1393 | `.EmailTypeLife` | 1 | 1 | ASM-FW-GISFW-Work | pxRadioButtons | 0/1 | relatif (agregat case) |
| 1394 | `.EmailTypeRetro` | 1 | 1 | ASM-FW-GISFW-Work | pxRadioButtons | 0/1 | relatif (agregat case) |
| 1395 | `.EmailTypeTL` | 1 | 1 | ASM-FW-GISFW-Work | pxRadioButtons | 0/1 | relatif (agregat case) |
| 1396 | `.EmailTypeUW` | 1 | 1 | ASM-FW-GISFW-Work | pxRadioButtons | 0/1 | relatif (agregat case) |
| 1397 | `.EmailTypeUWEDM` | 1 | 1 | ASM-FW-GISFW-Work | pxRadioButtons | 0/1 | relatif (agregat case) |
| 1398 | `.EstiamsiKlaim` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfLossClaim | pxNumber | /1 | relatif (agregat case) |
| 1399 | `.ETA` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1400 | `.FilterTermForOpportunity` | 1 | 1 | Data-Portal | pxTextInput | 0/1 | relatif (agregat case) |
| 1401 | `.FormerShipName` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1402 | `.GLDate` | 1 | 1 | ASM-FW-GISFW-Data-Payment | pxDateTime | /1 | relatif (agregat case) |
| 1403 | `.GoodName2` | 1 | 1 | ASM-FW-GISFW-Data-Good | pxTextArea | /1 | relatif (agregat case) |
| 1404 | `.GOODS` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1405 | `.GOODS_NAME` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1406 | `.GRT` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1407 | `.HandingFee` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1408 | `.HasilAI` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText | /1 | relatif (agregat case) |
| 1409 | `.HASILD1` | 1 | 1 | ASM-FW-GISFW-Data-Search | pxNumber | /1 | relatif (agregat case) |
| 1410 | `.HASILD2` | 1 | 1 | ASM-FW-GISFW-Data-Search | pxNumber | /1 | relatif (agregat case) |
| 1411 | `.Header` | 1 | 1 | ASM-FW-GISFW-Data-AnekaOtherSchedule | pxTextInput | /1 | relatif (agregat case) |
| 1412 | `.Imo` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1413 | `.InccuredKlaim` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfLossClaim | pxNumber | /1 | relatif (agregat case) |
| 1414 | `.InputCommentDeductible` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextArea | 0/1 | relatif (agregat case) |
| 1415 | `.InstallmentPct` | 1 | 1 | ASM-FW-GISFW-Data-TreatyInInstallment | pxNumber | /1 | relatif (agregat case) |
| 1416 | `.INTEREST` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1417 | `.InwardScale.Desc` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText | /1 | relatif (agregat case) |
| 1418 | `.InwardScale.ShareMin` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | /1 | relatif (agregat case) |
| 1419 | `.IsDeductibleAcceptance` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | 0/1 | relatif (agregat case) |
| 1420 | `.IsEditable` | 1 | 1 | ASM-FW-GISFW-Data-Correspondence | pxDropdown | /1 | relatif (agregat case) |
| 1421 | `.IsEditReffNumber` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | 0/1 | relatif (agregat case) |
| 1422 | `.IsShowFacRetro` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | relatif (agregat case) |
| 1423 | `.IsShowFormula` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxCheckbox | 0/1 | relatif (agregat case) |
| 1424 | `.ITEM_TYPE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1425 | `.ItemTypeID` | 1 | 1 | ASM-FW-GISFW-Data-PropertyItem | pxDropdown | /1 | relatif (agregat case) |
| 1426 | `.KATEGORI_2` | 1 | 1 | ASM-FW-GISFW-Int-DOCUMENT_POLIS | (tanpa) | /1 | relatif (agregat case) |
| 1427 | `.KLAUSUL_NOTE` | 1 | 1 | ASM-FW-GISFW-Int-M_KLAUSUL_PLAN | pxDisplayText | /1 | relatif (agregat case) |
| 1428 | `.Leader` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxCheckbox | /1 | relatif (agregat case) |
| 1429 | `.LeaderPolicyNo` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxTextInput | /1 | relatif (agregat case) |
| 1430 | `.LIMIT_OF_LIABILITY` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1431 | `.LimitCase` | 1 | 1 | (kosong) | pxNumber | 0/1 | relatif (agregat case) |
| 1432 | `.LOAN_ACCOUNT_NO` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1433 | `.LOAN_NAME` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1434 | `.LOAN_TYPE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1435 | `.Location` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfLossClaim | (tanpa) | /1 | relatif (agregat case) |
| 1436 | `.LocationNo` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfLossClaim | (tanpa) | /1 | relatif (agregat case) |
| 1437 | `.LOL` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1438 | `.LOSSLIMIT` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1439 | `.LossRatio` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfLossClaim | pxNumber | /1 | relatif (agregat case) |
| 1440 | `.LossRatio1YearAmount` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | /1 | relatif (agregat case) |
| 1441 | `.LossRatio1YearPercent` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | /1 | relatif (agregat case) |
| 1442 | `.LossRatio35YearAmount` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | /1 | relatif (agregat case) |
| 1443 | `.LossRatio35YearPercent` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxNumber | /1 | relatif (agregat case) |
| 1444 | `.MaintenanceStatus` | 1 | 1 | ASM-FW-GISFW-Data-FEA | pxDropdown | /1 | relatif (agregat case) |
| 1445 | `.MINMAX1` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1446 | `.MINMAX2` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1447 | `.MINMAX3` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1448 | `.MINMAX4` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1449 | `.MINMAX5` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1450 | `.MonthName` | 1 | 1 | ASM-FW-GISFW-Data-StockDeclaration | pxDropdown | 0/1 | relatif (agregat case) |
| 1451 | `.NamaCoverage` | 1 | 1 | ASM-FW-GISFW-Int-COVERAGE_FACIN | (tanpa) | /1 | relatif (agregat case) |
| 1452 | `.NAMAFILE` | 1 | 1 | ASM-FW-GISFW-Int-DOCUMENT_POLIS | pxLink | 0/1 | relatif (agregat case) |
| 1453 | `.NetPremiAfterPPH` | 1 | 1 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 0/1 | relatif (agregat case) |
| 1454 | `.NetPremiAfterPPH2` | 1 | 1 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 0/1 | relatif (agregat case) |
| 1455 | `.NetPremiAfterPPN` | 1 | 1 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 0/1 | relatif (agregat case) |
| 1456 | `.NetPremiAfterPPN2` | 1 | 1 | ASM-FW-GISFW-Data-LimitSummaryList | pxNumber | 0/1 | relatif (agregat case) |
| 1457 | `.NO` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1458 | `.NO_KONTRAK` | 1 | 1 | ASM-FW-GISFW-Data-Aneka | pxTextInput | /1 | relatif (agregat case) |
| 1459 | `.NO_KTP` | 1 | 1 | ASM-FW-GISFW-Data-Aneka | pxTextInput | /1 | relatif (agregat case) |
| 1460 | `.NO_NPWP` | 1 | 1 | ASM-FW-GISFW-Data-Aneka | pxTextInput | /1 | relatif (agregat case) |
| 1461 | `.NoteDokter` | 1 | 1 | Data-Party-Person | pxTextArea | 0/1 | relatif (agregat case) |
| 1462 | `.NOTES` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1463 | `.NRT` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1464 | `.NumberOfMember` | 1 | 1 | ASM-FW-GISFW-Data-Plan-PremiumList | pxDisplayText | /1 | relatif (agregat case) |
| 1465 | `.OBJECT_TYPE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1466 | `.Occupation` | 1 | 1 | ASM-FW-GISFW-Data-AnekaParticipant | pxTextInput | /1 | relatif (agregat case) |
| 1467 | `.OCCUPATION_CODE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1468 | `.OfferFacIn.CommentCeding` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1469 | `.OfferFacIn.DateValidity` | 1 | 1 | ASM-FW-GISFW-Work | pxDateTime | 0/1 | relatif (agregat case) |
| 1470 | `.OfferFacIn.DaysValidity` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | 0/1 | relatif (agregat case) |
| 1471 | `.OfferFacIn.FacRetroDetails.AdditionalInfo` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | /1 | relatif (agregat case) |
| 1472 | `.OfferFacIn.IsBanding` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | relatif (agregat case) |
| 1473 | `.OfferFacIn.PolicyData.Payment.Installment` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1474 | `.OfferFacIn.PolicyData.Payment.RICommision` | 1 | 1 | ASM-FW-GISFW-Work | pxNumber | 0/1 | relatif (agregat case) |
| 1475 | `.OfferFacIn.PolicyData.PolicyNo` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1476 | `.OfferFacIn.PolicyData.ProdDateTime` | 1 | 1 | ASM-FW-GISFW-Work | pxDateTime | 0/1 | relatif (agregat case) |
| 1477 | `.OfferFacIn.PremiMB` | 1 | 1 | ASM-FW-GISFW-Work | pxNumber | 0/1 | relatif (agregat case) |
| 1478 | `.OfferFacIn.PremiOther` | 1 | 1 | ASM-FW-GISFW-Work | pxNumber | 0/1 | relatif (agregat case) |
| 1479 | `.OfferFacIn.PremiPAR` | 1 | 1 | ASM-FW-GISFW-Work | pxNumber | 0/1 | relatif (agregat case) |
| 1480 | `.OfferFacIn.PremiPL` | 1 | 1 | ASM-FW-GISFW-Work | pxNumber | 0/1 | relatif (agregat case) |
| 1481 | `.OfferFacIn.QuotationData.CedingCoName` | 1 | 1 | ASM-FW-GISFW-Work | pxDisplayText | /1 | relatif (agregat case) |
| 1482 | `.OfferFacIn.QuotationData.EdmDate` | 1 | 1 | ASM-FW-GISFW-Work | pxDateTime | /1 | relatif (agregat case) |
| 1483 | `.OfferFacIn.QuotationData.OldPolicyNo` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1484 | `.OfferFacIn.ReceivedRiSlip` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | relatif (agregat case) |
| 1485 | `.OfferFacIn.WaitingBindDate` | 1 | 1 | ASM-FW-GISFW-Work | pxDateTime | 0/1 | relatif (agregat case) |
| 1486 | `.OperatorName` | 1 | 1 | ASM-FW-GISFW-Data-SuggestList | pxDisplayText | /1 | relatif (agregat case) |
| 1487 | `.OVR_COMM` | 1 | 1 | ASM-FW-GISFW-Int-RETROCESSIONLIFE | pxNumber | /1 | relatif (agregat case) |
| 1488 | `.OVR_COMM_AMOUNT` | 1 | 1 | ASM-FW-GISFW-Int-RETROCESSIONLIFE | pxNumber | /1 | relatif (agregat case) |
| 1489 | `.Password` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxPassword | 0/1 | relatif (agregat case) |
| 1490 | `.PaymentAmount` | 1 | 1 | ASM-FW-GISFW-Data-Payment | pxNumber | /1 | relatif (agregat case) |
| 1491 | `.PaymentId` | 1 | 1 | ASM-FW-GISFW-Data-Payment | (tanpa) | /1 | relatif (agregat case) |
| 1492 | `.PCT_COMM` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1493 | `.PCT_DISCOUNT` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1494 | `.PCT_LOADING` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1495 | `.PCT_LOL` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1496 | `.PCT_RATE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1497 | `.PCT_SCALE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1498 | `.PCT_SHARE_RNM` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1499 | `.PctAdjust1` | 1 | 1 | ASM-FW-GISFW-Data-PropertyItem | pxDropdown | /1 | relatif (agregat case) |
| 1500 | `.PCTDEDUCTIBLE1` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1501 | `.PCTDEDUCTIBLE3` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1502 | `.PCTDEDUCTIBLE4` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1503 | `.PCTDEDUCTIBLE5` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1504 | `.PctPremiAllObj` | 1 | 1 | ASM-FW-GISFW-Data-FacOffer | pxNumber | 0/1 | relatif (agregat case) |
| 1505 | `.PctPremiAllObjUSD` | 1 | 1 | ASM-FW-GISFW-Data-FacOffer | pxNumber | 0/1 | relatif (agregat case) |
| 1506 | `.PctRNM` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfDecline | pxNumber | 0/1 | relatif (agregat case) |
| 1507 | `.PctSOB` | 1 | 1 | ASM-FW-GISFW-Data-CauseOfDecline | pxNumber | 0/1 | relatif (agregat case) |
| 1508 | `.PctTotal` | 1 | 1 | ASM-FW-GISFW-Data-TreatyInInstallment | pxNumber | 0/1 | relatif (agregat case) |
| 1509 | `.PctTreatyLimit` | 1 | 1 | (kosong) | pxNumber | 0/1 | relatif (agregat case) |
| 1510 | `.PercentageAdjustment` | 1 | 1 | ASM-FW-GISFW-Data-PropertyItem | pxDropdown | /1 | relatif (agregat case) |
| 1511 | `.PercentBrokerage` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1512 | `.PercentHandlingFee` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1513 | `.PercentPPH` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1514 | `.PercentPPN` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1515 | `.Person` | 1 | 1 | ASM-FW-GISFW-Data-AnekaParticipant | pxTextInput | /1 | relatif (agregat case) |
| 1516 | `.PIC` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxDisplayText | /1 | relatif (agregat case) |
| 1517 | `.PlanList(1).Limit` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxDisplayText | /1 | relatif (agregat case) |
| 1518 | `.Policy.EndDateTime` | 1 | 1 | ASM-FW-GISFW-Work | pxDateTime | /1 | relatif (agregat case) |
| 1519 | `.Policy.Note` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1520 | `.Policy.Payment.PctDeduction2` | 1 | 1 | ASM-FW-GISFW-Work | pxNumber | 0/1 | relatif (agregat case) |
| 1521 | `.Policy.PolicyNo` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1522 | `.Policy.pyTitle` | 1 | 1 | ASM-FW-GISFW-Work | pxDropdown | /1 | relatif (agregat case) |
| 1523 | `.Policy.QQName` | 1 | 1 | ASM-FW-GISFW-Work | pxDropdown | /1 | relatif (agregat case) |
| 1524 | `.Policy.RefNo` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1525 | `.Policy.StartDateTime` | 1 | 1 | ASM-FW-GISFW-Work | pxDateTime | /1 | relatif (agregat case) |
| 1526 | `.Policy.STNC` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1527 | `.Policy.STNCDate` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1528 | `.Policy.SumOfTSI` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1529 | `.Policy.TheInsured` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | relatif (agregat case) |
| 1530 | `.Policy.TypeOfCoins` | 1 | 1 | ASM-FW-GISFW-Work | pxDropdown | /1 | relatif (agregat case) |
| 1531 | `.PolicyData.CedantRetention` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | /1 | relatif (agregat case) |
| 1532 | `.PolicyData.EndorsementNo` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | relatif (agregat case) |
| 1533 | `.PolicyData.PolicyNo` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | relatif (agregat case) |
| 1534 | `.PolicyNo` | 1 | 1 | ASM-FW-GISFW-Int-policyjson | pxTextInput | 0/1 | relatif (agregat case) |
| 1535 | `.PremiShare` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1536 | `.PREMIUM_COVERAGE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1537 | `.PREMIUM_RNM` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1538 | `.PREMIUM_SPREADED_GROSS` | 1 | 1 | ASM-FW-GISFW-Int-RETROCESSIONLIFE | pxNumber | 0/1 | relatif (agregat case) |
| 1539 | `.PREMIUM_SPREADED_NET` | 1 | 1 | ASM-FW-GISFW-Int-RETROCESSIONLIFE | pxNumber | 0/1 | relatif (agregat case) |
| 1540 | `.PremiumAfterPPN` | 1 | 1 | ASM-FW-GISFW-Data-Installment | pxTextInput | 0/1 | relatif (agregat case) |
| 1541 | `.PremiumAfterTax` | 1 | 1 | ASM-FW-GISFW-Data-Installment | pxTextInput | 0/1 | relatif (agregat case) |
| 1542 | `.PrivateFireBrigade` | 1 | 1 | ASM-FW-GISFW-Data-FEA | pxDropdown | /1 | relatif (agregat case) |
| 1543 | `.PrivateTruckBrigade` | 1 | 1 | ASM-FW-GISFW-Data-FEA | pxTextInput | /1 | relatif (agregat case) |
| 1544 | `.PRODUCT` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1545 | `.ProductionDate` | 1 | 1 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxDateTime | 0/1 | relatif (agregat case) |
| 1546 | `.Property.RiskLocation.CityName` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxAutoComplete | /1 | relatif (agregat case) |
| 1547 | `.Property.RiskLocation.DistrictName` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxAutoComplete | /1 | relatif (agregat case) |
| 1548 | `.Property.RiskLocation.RWName` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | pxAutoComplete | /1 | relatif (agregat case) |
| 1549 | `.pxCommitDateTime` | 1 | 1 | Data-WorkAttach-File | pxDateTime | /1 | relatif (agregat case) |
| 1550 | `.pyAttachStream` | 1 | 1 | Data-WorkAttach-File | ASMThumbnail | 0/1 | relatif (agregat case) |
| 1551 | `.pyScore` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxSlider | 0/1 | relatif (agregat case) |
| 1552 | `.pySearchCase` | 1 | 1 | ASM-FW-GISFW-Work-NB | pxTextInput | 0/1 | relatif (agregat case) |
| 1553 | `.pySearchString` | 1 | 1 | ASM-FW-GISFW-Work-NB | pxTextInput | 0/1 | relatif (agregat case) |
| 1554 | `.pyTitle` | 1 | 1 | ASM-FW-GISFW-Data-Policy | pxTextInput | /1 | relatif (agregat case) |
| 1555 | `.QQ` | 1 | 1 | ASM-FW-GISFW-Int-policyjson | (tanpa) | 0/1 | relatif (agregat case) |
| 1556 | `.Question` | 1 | 1 | ASM-FW-GISFW-Data-SubjectTo | pxTextInput | /1 | relatif (agregat case) |
| 1557 | `.Quotation.BranchName` | 1 | 1 | ASM-FW-GISFW-Data-Policy | pxTextInput | /1 | relatif (agregat case) |
| 1558 | `.Quotation.BusinessName` | 1 | 1 | ASM-FW-GISFW-Data-Policy | pxTextInput | /1 | relatif (agregat case) |
| 1559 | `.Quotation.MarketingName` | 1 | 1 | ASM-FW-GISFW-Data-Policy | pxTextInput | /1 | relatif (agregat case) |
| 1560 | `.QuotationData.EdmTypeNew` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | /1 | relatif (agregat case) |
| 1561 | `.QuotationData.InputMonth` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | 0/1 | relatif (agregat case) |
| 1562 | `.QuotationData.IsGroup` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | /1 | relatif (agregat case) |
| 1563 | `.QuotationData.Type` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxDropdown | /1 | relatif (agregat case) |
| 1564 | `.RATE_COVERAGE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1565 | `.Rebuilt` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1566 | `.RecipientMail` | 1 | 1 | ASM-FW-GISFW-Data-Correspondence | (tanpa) | /1 | relatif (agregat case) |
| 1567 | `.Register` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1568 | `.REINSTYPENAME` | 1 | 1 | (kosong) | pxTextInput | 0/1 | relatif (agregat case) |
| 1569 | `.RemarkFEA` | 1 | 1 | ASM-FW-GISFW-Data-FEA | pxTextArea | /1 | relatif (agregat case) |
| 1570 | `.RISK_LOCATION` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1571 | `.RiSlipCeding` | 1 | 1 | ASM-FW-GISFW-Int-policyjson | pxDisplayText | /1 | relatif (agregat case) |
| 1572 | `.RNM_PREMIUM` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1573 | `.RNM_SHARE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1574 | `.RNM_SHARE_PCT` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1575 | `.RP` | 1 | 1 | (kosong) | pxNumber | 0/1 | relatif (agregat case) |
| 1576 | `.ScoringRisk.DataScoringRiskList(1).OccupationNote` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | relatif (agregat case) |
| 1577 | `.ScoringRisk.DataScoringRiskList(1).TopRiskLocation` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | relatif (agregat case) |
| 1578 | `.ScoringRisk.DataScoringRiskList(2).OccupationNote` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxAutoComplete | 0/1 | relatif (agregat case) |
| 1579 | `.ScoringRisk.DataScoringRiskList(2).TopRiskLocation` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | relatif (agregat case) |
| 1580 | `.ScoringRisk.ScoringResult.AboveAverage` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1581 | `.ScoringRisk.ScoringResult.AcceptableNotApproval` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1582 | `.ScoringRisk.ScoringResult.AcceptableWithApproval` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1583 | `.ScoringRisk.ScoringResult.Average` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1584 | `.ScoringRisk.ScoringResult.BelowAverage` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1585 | `.ScoringRisk.ScoringResult.Good` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1586 | `.ScoringRisk.ScoringResult.LossRatio` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1587 | `.ScoringRisk.ScoringResult.NotAcceptable` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1588 | `.ScoringRisk.ScoringResult.ObjectConditions` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1589 | `.ScoringRisk.ScoringResult.Occupation` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1590 | `.ScoringRisk.ScoringResult.OperationalDirector` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1591 | `.ScoringRisk.ScoringResult.OperationalDivHead` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1592 | `.ScoringRisk.ScoringResult.Others` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1593 | `.ScoringRisk.ScoringResult.Poor` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1594 | `.ScoringRisk.ScoringResult.PresidentDirector` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1595 | `.ScoringRisk.ScoringResult.TechnicalDirector` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | /1 | relatif (agregat case) |
| 1596 | `.ScoringRisk.ScoringResult.TotalSumInsuredR1` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | relatif (agregat case) |
| 1597 | `.ScoringRisk.ScoringResult.TotalSumInsuredR2` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxNumber | /1 | relatif (agregat case) |
| 1598 | `.ScoringRisk.Status` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | relatif (agregat case) |
| 1599 | `.SenderMail` | 1 | 1 | ASM-FW-GISFW-Data-Correspondence | pxTextInput | 0/1 | relatif (agregat case) |
| 1600 | `.Sex` | 1 | 1 | ASM-FW-GISFW-Data-AnekaParticipant | pxTextInput | /1 | relatif (agregat case) |
| 1601 | `.SHARE_RNM` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1602 | `.ShareMax` | 1 | 1 | (kosong) | pxNumber | 0/1 | relatif (agregat case) |
| 1603 | `.ShareMin` | 1 | 1 | (kosong) | pxNumber | 0/1 | relatif (agregat case) |
| 1604 | `.ShipName` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1605 | `.SourceOfBusinessName` | 1 | 1 | ASM-FW-GISFW-Int-policyjson | pxTextInput | 0/1 | relatif (agregat case) |
| 1606 | `.SubLimitNote` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxTextArea | 0/1 | relatif (agregat case) |
| 1607 | `.SUM_INSURED_BI` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1608 | `.SUM_INSURED_MD` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1609 | `.SumOfTSI` | 1 | 1 | Data-Address | (tanpa) | /1 | relatif (agregat case) |
| 1610 | `.Tahun` | 1 | 1 | ASM-FW-GISFW-Int-INWARDSCALE | pxTextInput | 0/1 | relatif (agregat case) |
| 1611 | `.TANGGAL` | 1 | 1 | ASM-FW-GISFW-Int-DOCUMENT_POLIS | (tanpa) | /1 | relatif (agregat case) |
| 1612 | `.TeamSOPRiskManagement` | 1 | 1 | ASM-FW-GISFW-Data-FEA | pxDropdown | /1 | relatif (agregat case) |
| 1613 | `.TeamSOPSafety` | 1 | 1 | ASM-FW-GISFW-Data-FEA | pxDropdown | /1 | relatif (agregat case) |
| 1614 | `.TempJaminan` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxTextInput | /1 | relatif (agregat case) |
| 1615 | `.TerritoryName` | 1 | 1 | ASM-FW-GISFW-Int-RISKADDRESS | pxDisplayText | /1 | relatif (agregat case) |
| 1616 | `.TOTALGOODS` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1617 | `.TotalNetPremiAfterPPN` | 1 | 1 | ASM-FW-GISFW-Data-TreatyInTotal | pxNumber | 0/1 | relatif (agregat case) |
| 1618 | `.TotalNetPremiAfterTax` | 1 | 1 | ASM-FW-GISFW-Data-TreatyInTotal | pxNumber | 0/1 | relatif (agregat case) |
| 1619 | `.TotalPayment` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn-Currency | pxNumber | /1 | relatif (agregat case) |
| 1620 | `.TotalPPHValue` | 1 | 1 | ASM-FW-GISFW-Data-TreatyInTotal | pxNumber | 0/1 | relatif (agregat case) |
| 1621 | `.TotalPPNValue` | 1 | 1 | ASM-FW-GISFW-Data-TreatyInTotal | pxNumber | 0/1 | relatif (agregat case) |
| 1622 | `.TOTALPRICE` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | (tanpa) | /1 | relatif (agregat case) |
| 1623 | `.TRADING` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1624 | `.TradingNote2` | 1 | 1 | ASM-FW-GISFW-Data-Trading | pxTextArea | /1 | relatif (agregat case) |
| 1625 | `.TSI_LOSSLIMIT` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1626 | `.TSI_MD_BI` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1627 | `.TSI_OBJECT` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxNumber | /1 | relatif (agregat case) |
| 1628 | `.TSI_OBJECT_ITEM` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1629 | `.TSI_SUBLIMIT` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1630 | `.TSIAdjustment` | 1 | 1 | ASM-FW-GISFW-Data-StockDeclaration | pxNumber | 0/1 | relatif (agregat case) |
| 1631 | `.TSIShare` | 1 | 1 | ASM-FW-GISFW-Data-Coins | pxNumber | /1 | relatif (agregat case) |
| 1632 | `.TSISublimit` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxTextInput | /1 | relatif (agregat case) |
| 1633 | `.TYPE_ADDRESS` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1634 | `.TYPEDEDUCTIBLE1` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1635 | `.TYPEDEDUCTIBLE3` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1636 | `.TYPEDEDUCTIBLE4` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1637 | `.TYPEDEDUCTIBLE5` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDropdown | /1 | relatif (agregat case) |
| 1638 | `.TypeOfDiscount` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxCheckbox | /1 | relatif (agregat case) |
| 1639 | `.UNITY` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1640 | `.ValueLoading` | 1 | 1 | ASM-FW-GISFW-Data-AnekaParticipant | pxTextInput | /1 | relatif (agregat case) |
| 1641 | `.Vehicle.FacOut` | 1 | 1 | ASM-FW-GISFW-Work | pxDropdown | 0/1 | relatif (agregat case) |
| 1642 | `.VESSELNAME` | 1 | 1 | ASM-FW-GISFW-Int-TABLE_B2B | pxDisplayText | /1 | relatif (agregat case) |
| 1643 | `.WPC` | 1 | 1 | ASM-FW-GISFW-Data-TreatyInInstallment | pxTextInput | /1 | relatif (agregat case) |
| 1644 | `.YearMake2` | 1 | 1 | ASM-FW-GISFW-Int-SHIP | pxDisplayText | /1 | relatif (agregat case) |
| 1645 | `AccumulationRisk.AccumulationType` | 1 | 1 | Data-Portal | pxDisplayText | /1 | LUAR agregat: AccumulationRisk |
| 1646 | `AccumulationRisk.CenterTransStatus` | 1 | 1 | Data-Portal | pxDisplayText | /1 | LUAR agregat: AccumulationRisk |
| 1647 | `AccumulationRisk.CZone` | 1 | 1 | Data-Portal | pxDisplayText | /1 | LUAR agregat: AccumulationRisk |
| 1648 | `AccumulationRisk.ID` | 1 | 1 | Data-Portal | pxDisplayText | /1 | LUAR agregat: AccumulationRisk |
| 1649 | `AccumulationRisk.Keyword` | 1 | 1 | Data-Portal | pxDisplayText | /1 | LUAR agregat: AccumulationRisk |
| 1650 | `AccumulationRisk.Note` | 1 | 1 | Data-Portal | pxTextInput | /1 | LUAR agregat: AccumulationRisk |
| 1651 | `AccumulationRisk.SyariahStatus` | 1 | 1 | Data-Portal | pxTextInput | 0/1 | LUAR agregat: AccumulationRisk |
| 1652 | `AdjustmentRisk.CARI18` | 1 | 1 | Data-Portal | pxDisplayText | /1 | LUAR agregat: AdjustmentRisk |
| 1653 | `AdjustmentRisk.CARI30` | 1 | 1 | Data-Portal | pxNumber | /1 | LUAR agregat: AdjustmentRisk |
| 1654 | `AdjustmentRisk.CARI31` | 1 | 1 | Data-Portal | pxNumber | /1 | LUAR agregat: AdjustmentRisk |
| 1655 | `AdjustmentRisk.CARI32` | 1 | 1 | Data-Portal | pxNumber | /1 | LUAR agregat: AdjustmentRisk |
| 1656 | `AdjustmentRisk.CARI33` | 1 | 1 | Data-Portal | pxNumber | /1 | LUAR agregat: AdjustmentRisk |
| 1657 | `AdjustmentRisk.CARI34` | 1 | 1 | Data-Portal | pxNumber | /1 | LUAR agregat: AdjustmentRisk |
| 1658 | `AdjustmentRisk.CARI35` | 1 | 1 | Data-Portal | pxNumber | /1 | LUAR agregat: AdjustmentRisk |
| 1659 | `Clause.ClauseLanguageID` | 1 | 1 | ASM-FW-GISFW-Work | pxDropdown | 0/1 | LUAR agregat: Clause |
| 1660 | `Clause.pyNote` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | 0/1 | LUAR agregat: Clause |
| 1661 | `CountAccumulation.CARI2` | 1 | 1 | (kosong) | pxNumber | /1 | LUAR agregat: CountAccumulation |
| 1662 | `CountAccumulation.CARI3` | 1 | 1 | (kosong) | pxNumber | /1 | LUAR agregat: CountAccumulation |
| 1663 | `CountAccumulation.CARI4` | 1 | 1 | (kosong) | pxNumber | /1 | LUAR agregat: CountAccumulation |
| 1664 | `FacOfferList.Attention` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | LUAR agregat: FacOfferList |
| 1665 | `FacOfferList.OurRef` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | LUAR agregat: FacOfferList |
| 1666 | `FacOfferList.ReinsurerName` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | LUAR agregat: FacOfferList |
| 1667 | `InputAccumulatedType.AccumulationType` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputAccumulatedType |
| 1668 | `InputAccumulatedType.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputAccumulatedType |
| 1669 | `InputAccumulatedType.Keyword` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputAccumulatedType |
| 1670 | `InputAccumulatedType.Note` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputAccumulatedType |
| 1671 | `InputAccumulation.AccumulationName` | 1 | 1 | @baseclass | pxDisplayText | /1 | LUAR agregat: InputAccumulation |
| 1672 | `InputAccumulation.AccumulationType` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputAccumulation |
| 1673 | `InputAccumulation.CZone` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputAccumulation |
| 1674 | `InputAccumulation.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputAccumulation |
| 1675 | `InputAccumulation.Keyword` | 1 | 1 | @baseclass | pxDropdown | 0/1 | LUAR agregat: InputAccumulation |
| 1676 | `InputAccumulation.Note` | 1 | 1 | @baseclass | pxTextArea | 0/1 | LUAR agregat: InputAccumulation |
| 1677 | `InputAccumulation.PostalCode` | 1 | 1 | @baseclass | pxAutoComplete | /1 | LUAR agregat: InputAccumulation |
| 1678 | `InputAccumulation.ProvinceName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputAccumulation |
| 1679 | `InputAccumulation.ScopeArea` | 1 | 1 | @baseclass | pxDropdown | 0/1 | LUAR agregat: InputAccumulation |
| 1680 | `InputAccumulationCode.CARI10` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxTextInput | /1 | LUAR agregat: InputAccumulationCode |
| 1681 | `InputBranch.BASTerritory` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1682 | `InputBranch.BranchName` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1683 | `InputBranch.BranchStatus` | 1 | 1 | @baseclass | pxRadioButtons | 0/1 | LUAR agregat: InputBranch |
| 1684 | `InputBranch.BranchType` | 1 | 1 | @baseclass | pxRadioButtons | 0/1 | LUAR agregat: InputBranch |
| 1685 | `InputBranch.CompanyID` | 1 | 1 | @baseclass | pxHidden | 0/1 | LUAR agregat: InputBranch |
| 1686 | `InputBranch.HSGradSalary` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1687 | `InputBranch.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputBranch |
| 1688 | `InputBranch.JabodetabekStatus` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1689 | `InputBranch.JHSGradSalary` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1690 | `InputBranch.Kanwil` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1691 | `InputBranch.KanwilGroup` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1692 | `InputBranch.Name` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1693 | `InputBranch.Status` | 1 | 1 | @baseclass | pxRadioButtons | 0/1 | LUAR agregat: InputBranch |
| 1694 | `InputBranch.Telephone` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1695 | `InputBranch.UndergradSalary` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputBranch |
| 1696 | `InputBranchDetail.CompanyID` | 1 | 1 | @baseclass | pxHidden | /1 | LUAR agregat: InputBranchDetail |
| 1697 | `InputCity.BranchName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputCity |
| 1698 | `InputCity.Email` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputCity |
| 1699 | `InputCity.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputCity |
| 1700 | `InputCity.JABODETABEKStatus` | 1 | 1 | @baseclass | pxRadioButtons | 0/1 | LUAR agregat: InputCity |
| 1701 | `InputCity.MOName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputCity |
| 1702 | `InputCity.Note` | 1 | 1 | @baseclass | pxTextArea | 0/1 | LUAR agregat: InputCity |
| 1703 | `InputCity.ProvinceName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputCity |
| 1704 | `InputData.CARI10` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | LUAR agregat: InputData |
| 1705 | `InputData.CARI33` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | 0/1 | LUAR agregat: InputData |
| 1706 | `InputData.CARI34` | 1 | 1 | ASM-FW-GISFW-Work | pxDateTime | /1 | LUAR agregat: InputData |
| 1707 | `InputData.CARI9` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | LUAR agregat: InputData |
| 1708 | `InputDistrict.CityName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputDistrict |
| 1709 | `InputDistrict.DistrictName` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputDistrict |
| 1710 | `InputDistrict.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputDistrict |
| 1711 | `InputMarketingOfficer.A_Leader` | 1 | 1 | @baseclass | pxCheckbox | 0/1 | LUAR agregat: InputMarketingOfficer |
| 1712 | `InputMarketingOfficer.BranchDetailID` | 1 | 1 | @baseclass | pxDropdown | /1 | LUAR agregat: InputMarketingOfficer |
| 1713 | `InputMarketingOfficer.BranchParent` | 1 | 1 | @baseclass | pxDropdown | /1 | LUAR agregat: InputMarketingOfficer |
| 1714 | `InputMarketingOfficer.ClientName` | 1 | 1 | @baseclass | pxAutoComplete | /1 | LUAR agregat: InputMarketingOfficer |
| 1715 | `InputMarketingOfficer.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputMarketingOfficer |
| 1716 | `InputMarketingOfficer.MOLeader` | 1 | 1 | @baseclass | pxAutoComplete | /1 | LUAR agregat: InputMarketingOfficer |
| 1717 | `InputMarketingOfficer.MOStatus` | 1 | 1 | @baseclass | pxRadioButtons | 0/1 | LUAR agregat: InputMarketingOfficer |
| 1718 | `InputNation.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputNation |
| 1719 | `InputNation.NationInitial` | 1 | 1 | @baseclass | pxTextArea | 0/1 | LUAR agregat: InputNation |
| 1720 | `InputNation.Note` | 1 | 1 | @baseclass | pxTextArea | 0/1 | LUAR agregat: InputNation |
| 1721 | `InputPar.CARIYEAR` | 1 | 1 | Data-Portal | pxTextInput | 0/1 | LUAR agregat: InputPar |
| 1722 | `InputParam.CARI5` | 1 | 1 | ASM-FW-GISFW-Data-Occupation | pxTextInput | 0/1 | LUAR agregat: InputParam |
| 1723 | `InputParam.CARIDESC` | 1 | 1 | Data-Portal | pxRichTextEditor | /1 | LUAR agregat: InputParam |
| 1724 | `InputProvince.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputProvince |
| 1725 | `InputProvince.NationName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputProvince |
| 1726 | `InputProvince.Note` | 1 | 1 | @baseclass | pxTextArea | 0/1 | LUAR agregat: InputProvince |
| 1727 | `InputRiskAddress.Address` | 1 | 1 | @baseclass | pxTextArea | 0/1 | LUAR agregat: InputRiskAddress |
| 1728 | `InputRiskAddress.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputRiskAddress |
| 1729 | `InputRiskAddress.NationName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputRiskAddress |
| 1730 | `InputRiskAddress.PostalCode` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputRiskAddress |
| 1731 | `InputRiskAddress.StatusTransCenter` | 1 | 1 | @baseclass | pxRadioButtons | 0/1 | LUAR agregat: InputRiskAddress |
| 1732 | `InputRiskAddress.TerritoryName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputRiskAddress |
| 1733 | `InputRiskAddress.Title` | 1 | 1 | @baseclass | pxDropdown | 0/1 | LUAR agregat: InputRiskAddress |
| 1734 | `InputRiskAddress.Type` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputRiskAddress |
| 1735 | `InputRW.AssessmentZone` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputRW |
| 1736 | `InputRW.CITYNAME` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputRW |
| 1737 | `InputRW.CZONE` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputRW |
| 1738 | `InputRW.DistrictName` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputRW |
| 1739 | `InputRW.ID` | 1 | 1 | @baseclass | pxTextInput | /1 | LUAR agregat: InputRW |
| 1740 | `InputRW.Note` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputRW |
| 1741 | `InputRW.PROVINCENAME` | 1 | 1 | @baseclass | pxAutoComplete | 0/1 | LUAR agregat: InputRW |
| 1742 | `InputRW.STS_AKTIF` | 1 | 1 | @baseclass | pxRadioButtons | 0/1 | LUAR agregat: InputRW |
| 1743 | `InputRW.ZipCode` | 1 | 1 | @baseclass | pxTextInput | 0/1 | LUAR agregat: InputRW |
| 1744 | `OutputParam.ERRMSG2` | 1 | 1 | Data-Party-Person | pxDisplayText | /1 | LUAR agregat: OutputParam |
| 1745 | `OutputParam.ERRMSG3` | 1 | 1 | Data-Party-Person | pxTextInput | /1 | LUAR agregat: OutputParam |
| 1746 | `OutputParam.ERRMSG6` | 1 | 1 | ASM-FW-GISFW-Data-Coverage | pxDisplayText | /1 | LUAR agregat: OutputParam |
| 1747 | `pyPortal.FilterForViewRetro` | 1 | 1 | ASM-FW-GISFW-Work | pxDropdown | 0/1 | LUAR agregat: pyPortal |
| 1748 | `pyWorkPage.IsSpecialCase` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxCheckbox | 0/1 | agregat case (eksplisit) |
| 1749 | `pyWorkPage.OfferFacIn.CurrentYear` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | (tanpa) | /1 | agregat case (eksplisit) |
| 1750 | `pyWorkPage.OfferFacIn.LossRatio1YearPercent` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | agregat case (eksplisit) |
| 1751 | `pyWorkPage.OfferFacIn.LossRatio3YearPercent` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | agregat case (eksplisit) |
| 1752 | `pyWorkPage.OfferFacIn.LossRatio5YearPercent` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | agregat case (eksplisit) |
| 1753 | `pyWorkPage.OfferFacIn.NetLossRatio1Year` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | agregat case (eksplisit) |
| 1754 | `pyWorkPage.OfferFacIn.NetLossRatio3Year` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | agregat case (eksplisit) |
| 1755 | `pyWorkPage.OfferFacIn.NetLossRatio5Year` | 1 | 1 | ASM-FW-GISFW-Work | pxTextInput | /1 | agregat case (eksplisit) |
| 1756 | `pyWorkPage.OfferFacIn.Parameters.CariOccupationConstruction` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | (tanpa) | /1 | agregat case (eksplisit) |
| 1757 | `pyWorkPage.OfferFacIn.Parameters.CariOccupationRiskCategory` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | (tanpa) | /1 | agregat case (eksplisit) |
| 1758 | `pyWorkPage.OfferFacIn.Parameters.CariScaleMax` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | (tanpa) | /1 | agregat case (eksplisit) |
| 1759 | `pyWorkPage.OfferFacIn.Parameters.CariScaleMin` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | (tanpa) | /1 | agregat case (eksplisit) |
| 1760 | `pyWorkPage.OfferFacIn.Parameters.CariTreatyLimit` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | (tanpa) | /1 | agregat case (eksplisit) |
| 1761 | `pyWorkPage.OfferFacIn.Parameters.LowestPctLimit` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | (tanpa) | /1 | agregat case (eksplisit) |
| 1762 | `pyWorkPage.OfferFacIn.PolicyData.EndDateTime` | 1 | 1 | ASM-FW-GISFW-Data-ScoringRisk | pxDateTime | /1 | agregat case (eksplisit) |
| 1763 | `pyWorkPage.OfferFacIn.PolicyData.StartDateTime` | 1 | 1 | ASM-FW-GISFW-Data-ScoringRisk | pxDateTime | /1 | agregat case (eksplisit) |
| 1764 | `pyWorkPage.PaymentData.AgingStatus` | 1 | 1 | ASM-FW-GISFW-Work-NB | pxTextInput | /1 | agregat case (eksplisit) |
| 1765 | `pyWorkPage.PaymentList(1).ListInstallment(1).Premium` | 1 | 1 | ASM-FW-GISFW-Data-Installment | pxNumber | /1 | agregat case (eksplisit) |
| 1766 | `pyWorkPage.TreatyIn.FacultativeShare` | 1 | 1 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | /1 | agregat case (eksplisit) |
| 1767 | `pyWorkPage.TreatyIn.InstallmentNo` | 1 | 1 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | /1 | agregat case (eksplisit) |
| 1768 | `pyWorkPage.TreatyIn.ReinsurerShare` | 1 | 1 | ASM-FW-GISFW-Data-PolicyTreatyIn | pxTextInput | /1 | agregat case (eksplisit) |
| 1769 | `SearchAccumulation.Area` | 1 | 1 | Data-Portal | pxAutoComplete | 0/1 | LUAR agregat: SearchAccumulation |
| 1770 | `SearchAccumulation.CenterTransStatus` | 1 | 1 | Data-Portal | pxTextInput | 0/1 | LUAR agregat: SearchAccumulation |
| 1771 | `SearchAccumulation.City` | 1 | 1 | Data-Portal | pxAutoComplete | 0/1 | LUAR agregat: SearchAccumulation |
| 1772 | `SearchAccumulation.CZone` | 1 | 1 | Data-Portal | pxAutoComplete | 0/1 | LUAR agregat: SearchAccumulation |
| 1773 | `SearchAccumulation.District` | 1 | 1 | Data-Portal | pxAutoComplete | 0/1 | LUAR agregat: SearchAccumulation |
| 1774 | `SearchAccumulation.ID` | 1 | 1 | Data-Portal | pxTextInput | 0/1 | LUAR agregat: SearchAccumulation |
| 1775 | `SearchAccumulation.Keyword` | 1 | 1 | Data-Portal | pxDropdown | 0/1 | LUAR agregat: SearchAccumulation |
| 1776 | `SearchAccumulation.Nation` | 1 | 1 | Data-Portal | pxAutoComplete | 0/1 | LUAR agregat: SearchAccumulation |
| 1777 | `SearchAccumulation.Note` | 1 | 1 | Data-Portal | pxTextInput | 0/1 | LUAR agregat: SearchAccumulation |
| 1778 | `SearchAccumulation.PostalCode` | 1 | 1 | Data-Portal | pxTextInput | 0/1 | LUAR agregat: SearchAccumulation |
| 1779 | `SearchAccumulation.ProvinceName` | 1 | 1 | Data-Portal | pxAutoComplete | 0/1 | LUAR agregat: SearchAccumulation |
| 1780 | `SearchAccumulation.Type` | 1 | 1 | Data-Portal | pxAutoComplete | 0/1 | LUAR agregat: SearchAccumulation |
| 1781 | `SFAResults.pxResults(1).pySummaryCount(1)` | 1 | 1 | Data-Portal | pxNumber | /1 | LUAR agregat: SFAResults |
| 1782 | `SFAResults.pxResults(1).pySummaryDateTime(1)` | 1 | 1 | Data-Portal | pxDateTime | /1 | LUAR agregat: SFAResults |
| 1783 | `SFAResults.pxResults(1).pySummaryDateTime(2)` | 1 | 1 | Data-Portal | pxDateTime | /1 | LUAR agregat: SFAResults |
| 1784 | `SFAResults.pxResults(1).pySummaryValue(1)` | 1 | 1 | Data-Portal | pxCurrency | /1 | LUAR agregat: SFAResults |
| 1785 | `TempIndex.BizCode` | 1 | 1 | ASM-FW-GISFW-Data-Deductible | pxTextInput | 0/1 | LUAR agregat: TempIndex |
| 1786 | `TempPolis.CARI24` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | LUAR agregat: TempPolis |
| 1787 | `TempPolis.CARI4` | 1 | 1 | ASM-FW-GISFW-Data-OfferFacIn | pxTextInput | /1 | LUAR agregat: TempPolis |
| 1788 | `test.CARI40` | 1 | 1 | ASM-FW-GISFW-Work | (tanpa) | /1 | LUAR agregat: test |
| 1789 | `ViewOfferStatus.AdditionalInfo` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1790 | `ViewOfferStatus.ASMShare` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1791 | `ViewOfferStatus.ClassOfConstruction` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1792 | `ViewOfferStatus.Clauses` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1793 | `ViewOfferStatus.Deductibles` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1794 | `ViewOfferStatus.Information` | 1 | 1 | ASM-FW-GISFW-Work | pxTextArea | 0/1 | LUAR agregat: ViewOfferStatus |
| 1795 | `ViewOfferStatus.Insured` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1796 | `ViewOfferStatus.LineBusiness` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1797 | `ViewOfferStatus.Location` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1798 | `ViewOfferStatus.Occupation` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1799 | `ViewOfferStatus.OurRetention` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1800 | `ViewOfferStatus.PerilsCovered` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1801 | `ViewOfferStatus.PeriodOfInsurance` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1802 | `ViewOfferStatus.PeriodOfReinsurance` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1803 | `ViewOfferStatus.ShareOffered` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1804 | `ViewOfferStatus.SumInsuredTSI` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |
| 1805 | `ViewOfferStatus.SurveyReport` | 1 | 1 | ASM-FW-GISFW-Work | pxCheckbox | 0/1 | LUAR agregat: ViewOfferStatus |

---

## Lampiran C - Seluruh RepeatGrid NB dan properti PageList-nya

629 layout `RepeatGrid` di 281 Section. `Kelas baris` = `pyPageListPropertyClass`; `(kosong)`
berarti ekspor tidak mencantumkannya - kelas baris **belum terverifikasi** untuk baris itu.
| Section | Properti PageList | Kelas baris | n grid |
| --- | --- | --- | ---: |
| AccumulationRisk | `AccumulationRiskDtl.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| AccumulationRisk | `AccumulationRiskPolicyDtl.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| AddFacOfferList_IsUW | `.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| AddFacOfferList | `.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| AddFacOfferList2 | `.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| AdjustmentRiskAccumulation | `RiskAccumAdj.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| AdjustmentRiskAccumulation | `RiskAccumCurYear.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| AdjustmentRiskAccumulation | `RiskAccumPrevYear.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| AnekaListFacOut_IsUW | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| AnekaListFacOut | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| AttachmentGridReas | `Attachment.pxResults` | ASM-FW-GISFW-Int-OFFERJSON | 1 |
| BenefitClauseRO | `.ClauseList` | ASM-FW-GISFW-Data-Clause | 1 |
| BenefitListRO | `.BenefitList` | ASM-FW-GISFW-Data-Benefit | 1 |
| CauseOfLoss_FacIn_IsUW | `.Property.ListCauseOfLoss` | ASM-FW-GISFW-Data-CauseOfLoss | 1 |
| CauseOfLoss_FacIn | `.Property.ListCauseOfLoss` | ASM-FW-GISFW-Data-CauseOfLoss | 1 |
| CauseOfLossClaim_FacIn | `.Property.ListCauseOfLossClaim` | ASM-FW-GISFW-Data-CauseOfLossClaim | 1 |
| CedingCedant_IsUW | `.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| CedingCedant | `.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| CedingCoHierarki | `pgRepPgSubSectionCedingCoHierarkiBBBB.pxResults` | ASM-FW-GISFW-Int-AGENT | 1 |
| ChooseClassofContraction | `pgRepPgSubSectionChooseClassofContractionBB.pxResults` | ASM-FW-GISFW-Int-TABLEOFLIMIT | 1 |
| ChooseClause | `clausePage.pxResults` | ASM-FW-GISFW-Int-CLAUSE | 1 |
| ChooseClauseFire | `SearchClauseListOutput.pxResults` | ASM-FW-GISFW-Int-CLAUSE | 1 |
| ChooseCoverage | `pgRepPgSubSectionChooseCoverageBB.pxResults` | ASM-FW-GISFW-Int-COVERAGE_FACIN | 1 |
| ChooseDeductible | `pgRepPgSubSectionChooseDeductibleBB.pxResults` | ASM-FW-GISFW-Int-DEDUCTIBLE | 1 |
| ChooseObject_Ship | `pgRepPgSubSectionChooseObject_ShipBB.pxResults` | ASM-FW-GISFW-Int-SHIP | 1 |
| ChooseOccupation_Dtl | `D_Occupation.pxResults` | ASM-FW-GISFW-Int-OCCUPATION | 1 |
| ChooseOccupation | `pgRepPgSubSectionChooseOccupationBB.pxResults` | ASM-FW-GISFW-Int-OCCUPATION | 1 |
| ChooseRiskAddress_ResultList | `pgRepPgSubSectionChooseRiskAddress_ResultListB.pxResults` | ASM-FW-GISFW-Int-RISKADDRESS | 1 |
| ChooseSubContract_FacIn | `pyReportContentPage.pxResults` | ASM-FW-GISFW-Int-SUBCONTRACT | 1 |
| ChooseSubContract | `pyReportContentPage.pxResults` | ASM-FW-GISFW-Int-SUBCONTRACT | 1 |
| ChooseZipCodeDtl | `pgRepPgSubSectionChooseZipCodeDtlBB.pxResults` | ASM-FW-GISFW-Int-RISKADDRESS | 1 |
| ClauseList | `D_PlanClause.pxResults` | ASM-FW-GISFW-Int-M_KLAUSUL_PLAN | 1 |
| CommisionCoverageAneka | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| CommisionItem | `.LayerList` | ASM-FW-GISFW-Data-Coverage | 1 |
| CommisionObjectAneka | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| Correspondence | `.OfferFacIn.CorrespondenceList` | ASM-FW-GISFW-Data-Correspondence | 1 |
| CorrespondenceContent | `.AttachmentList` | Embed-DragDropFile | 2 |
| CorrespondenceContent | `.RecipientList` | ASM-FW-GISFW-Data-Correspondence | 1 |
| CoverageCommisionList_IsUW | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| CoverageCommisionList_IsUW | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 2 |
| CoverageCommisionList_IsUW | `.PersonList` | Data-Party-Person | 2 |
| CoverageCommisionList_IsUW | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| CoverageCommisionList | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| CoverageCommisionList | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 2 |
| CoverageCommisionList | `.PersonList` | Data-Party-Person | 3 |
| CoverageCommisionList | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| CoverageItem_IsUW | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 2 |
| CoverageItem_IsUW | `.LayerList` | ASM-FW-GISFW-Data-Coverage | 1 |
| CoverageItem | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 2 |
| CoverageItem | `.LayerList` | ASM-FW-GISFW-Data-Coverage | 1 |
| CoverageList_IsUW | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| CoverageList | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| CoveragePropertyFacOut | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| CoverageSpreadingList_IsUW | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| CoverageSpreadingList_IsUW | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 2 |
| CoverageSpreadingList_IsUW | `.PersonList` | Data-Party-Person | 3 |
| CoverageSpreadingList_IsUW | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| CoverageSpreadingList | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| CoverageSpreadingList | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 2 |
| CoverageSpreadingList | `.PersonList` | Data-Party-Person | 3 |
| CoverageSpreadingList | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| CoverageSummary | `D_CoverageSummary.pxResults` | ASM-FW-GISFW-Data-Coverage | 1 |
| DetailCoverageListFire_GCNM | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 1 |
| DetailCoverageTravelGrid_GCNM | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 1 |
| DetailDeptHeadTreatyIn_UW | `.BreakDownSpreadList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| DetailDeptHeadTreatyIn_UW | `.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| DetailDeptHeadTreatyIn_UW | `.SpreadingRiskList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| DetailPolicyTreatyIn | `.BreakDownSpreadList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| DetailPolicyTreatyIn | `.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| DetailPolicyTreatyIn | `.SpreadingRiskList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.PolicyTreatyIn.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.LimitFacShareSummaryList` | ASM-FW-GISFW-Data-LimitSummaryList | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.LimitShareSummaryList` | ASM-FW-GISFW-Data-LimitSummaryList | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.LimitSummaryList` | ASM-FW-GISFW-Data-LimitSummaryList | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalFacShareDeductionNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalFacShareGrossNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalFacShareNetNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalFacShareRnmNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalLimitDeductblNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalLimitIOONP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalLimitMDPNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalShareDeductionNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalShareGrossNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalShareNetNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalShareRnmNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalSpreadedNetPremi` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportional | `pyWorkPage.TreatyIn.TotalSpreadedNetPremiRI` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.Installment` | ASM-FW-GISFW-Data-TreatyInInstallment | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.LimitShareSummaryList` | ASM-FW-GISFW-Data-LimitSummaryList | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.LimitSummaryList` | ASM-FW-GISFW-Data-LimitSummaryList | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.TotalLimitDeductblNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.TotalLimitMDPNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.TotalShareDeductionNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.TotalShareGrossNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.TotalShareNetNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.TotalSpreadedNetPremi` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyInNonProportionalEDM | `pyWorkPage.TreatyIn.ValueDifference.TotalSpreadedNetPremiRI` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.PolicyTreatyIn.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.LimitShareSummaryList` | ASM-FW-GISFW-Data-LimitSummaryList | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.LimitSummaryList` | ASM-FW-GISFW-Data-LimitSummaryList | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.TotalLimitDeductblNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.TotalLimitIOONP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.TotalLimitMDPNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.TotalShareDeductionNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.TotalShareGrossNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.TotalShareNetNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| DetailPolicyTreatyOutNonProportional | `pyWorkPage.TreatyIn.TotalShareRnmNP` | ASM-FW-GISFW-Data-TreatyInTotal | 1 |
| EmailSection_IsUW | `pyWorkPage.QuotationList` | ASM-FW-GISFW-Data-Quotation | 1 |
| EmailSection_IsUW | `TempViewSuggest.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | 1 |
| EmailSection_Life | `TempViewSuggest.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | 1 |
| EmailSection | `TempViewSuggest.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | 1 |
| EmailSectionCeding | `TempViewSuggest.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | 1 |
| FacultativeLetter | `.OfferFacIn.FacRetroList` | ASM-FW-GISFW-Data-FacOffer | 1 |
| FEAList_IsUW | `.FEAList` | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | 1 |
| FEAList | `.FEAList` | ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList | 1 |
| FireSummarySection | `D_LocationSummary.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| GeneralDeptHeadTreatyIn_UW | `.BreakDownSpreadList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| GeneralDeptHeadTreatyIn_UW | `.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| GeneralDeptHeadTreatyIn_UW | `.SpreadingRiskList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| GeneralPolicyTreatyIn | `.BreakDownSpreadList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| GeneralPolicyTreatyIn | `.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| GeneralPolicyTreatyIn | `.SpreadingRiskList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| GolfCoverageSummary | `D_GrowingTreesCoverageSummary.pxResults` | ASM-FW-GISFW-Data-Coverage | 1 |
| GolfLocationDetail | `D_AnekaList.pxResults` | ASM-FW-GISFW-Data-Aneka | 1 |
| GolfLocationDetail | `D_GrowingTreesCoverageSummary.pxResults` | ASM-FW-GISFW-Data-Coverage | 1 |
| GolfObjectItemSummary | `D_AnekaList.pxResults` | ASM-FW-GISFW-Data-Aneka | 1 |
| GolfSummarySection | `D_GolfLocationSummary.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| GridViewFollowingNB | `TempListFollowing.pxResults` | ASM-FW-GISFW-Int-policyjson | 1 |
| GrowingTreesLocationDetail | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| GrowingTreesLocationDetail | `D_AnekaList.pxResults` | ASM-FW-GISFW-Data-Aneka | 1 |
| GrowingTreesLocationDetail | `D_GrowingTreesCoverageSummary.pxResults` | ASM-FW-GISFW-Data-Coverage | 1 |
| GrowingTreesOccupationSummary | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| GrowingTreesSummarySection | `D_GrowingTreesLocationSummary.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| HistoricalSurveyReportDtl | `.QuotationData.SurveyReportList` | ASM-FW-GISFW-Data-Quotation | 1 |
| HistoricalSurveyReportDtlUW | `.QuotationData.SurveyReportList` | ASM-FW-GISFW-Data-Quotation | 1 |
| InputAnekaParticipant | `.AnekaParticipantList` | ASM-FW-GISFW-Data-AnekaParticipant | 1 |
| InputCauseOfDecline | `CauseOfDeclineList.pxResults` | ASM-FW-GISFW-Data-CauseOfDecline | 1 |
| InputCauseOfDecline | `CauseOfDeclineList2.pxResults` | ASM-FW-GISFW-Data-CauseOfDecline | 1 |
| InputClause_ViewDtl | `.ArgumentList` | ASM-FW-GISFW-Data-Argument | 1 |
| InputClauseFire_ViewDtl | `.ArgumentList` | ASM-FW-GISFW-Data-Argument | 1 |
| InputCommissionLife_FacIn | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCommissionMBU_FacIn | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputCommissionMC_FacIn | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputCommissionPA_FacIn | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputCoverageAneka_FacIn | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 2 |
| InputCoverageCargo_FacIn_IsUW | `.AdditionalCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageCargo_FacIn_IsUW | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 2 |
| InputCoverageCargo_FacIn | `.AdditionalCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageCargo_FacIn | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 2 |
| InputCoverageCommisionFire | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageDeductibleAneka_GCNM | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageDeductibleHull_GCNM | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageFire_IsUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageFire | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageSpreadingAneka_IsUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageSpreadingAneka | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputCoverageSpreadingFire_IsUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageSpreadingFire | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageSpreadingMarineCargo | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputCoverageSpreadingMarineCargoISUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputCoverageSpreadingMBU | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputCoverageSpreadingMBUISUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputDeductible_GCNM | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 1 |
| InputDeductibleHull_GCNM | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 1 |
| InputDtlCargo_FacIn_IsUW | `.PolicyData.Ship.AdditionalShip` | ASM-FW-GISFW-Int-MSHIP | 1 |
| InputDtlCargo_FacIn | `.PolicyData.Ship.AdditionalShip` | ASM-FW-GISFW-Int-MSHIP | 1 |
| InputDtlClause_FacIn_IsUW | `.ClauseList` | ASM-FW-GISFW-Data-Clause | 3 |
| InputDtlClause_FacIn_IsUW | `.Policy.ClauseList` | ASM-FW-GISFW-Data-Clause | 1 |
| InputDtlClause_FacIn_IsUW | `.WarrantyList` | ASM-FW-GISFW-Data-Warranty | 1 |
| InputDtlClause_FacIn_IsUW | `AttachmentClauses.pxResults` | ASM-FW-GISFW-Int-OFFERJSON | 1 |
| InputDtlClause_FacIn_IsUW | `pgRepPgSubSectionInputDtlClause_FacIn_IsUWBBBBBBB.pxResults` | Link-Attachment | 1 |
| InputDtlClause_FacIn | `.ClauseList` | ASM-FW-GISFW-Data-Clause | 3 |
| InputDtlClause_FacIn | `.Policy.ClauseList` | ASM-FW-GISFW-Data-Clause | 1 |
| InputDtlClause_FacIn | `.WarrantyList` | ASM-FW-GISFW-Data-Warranty | 1 |
| InputDtlClause_FacIn | `AttachmentClauses.pxResults` | ASM-FW-GISFW-Int-OFFERJSON | 1 |
| InputDtlClause_FacIn | `pgRepPgSubSectionInputDtlClause_FacInBBBBBBB.pxResults` | Link-Attachment | 1 |
| InputDtlCoverage_FacIn_IsUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputDtlCoverage_FacIn_IsUW | `.TotalTSIPremiGrossList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| InputDtlCoverage_FacIn | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputDtlCoverage_FacIn | `.TotalTSIPremiGrossList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| InputDtlGoods_IsUW | `.PackingList` | ASM-FW-GISFW-Data-Packing | 1 |
| InputDtlGoods | `.PackingList` | ASM-FW-GISFW-Data-Packing | 1 |
| InputDtlLayerSpreading_FacIn | `.RetroList` | ASM-FW-GISFW-Int-RETROCESSIONLIFE | 1 |
| InputDtlObject_FacIn_IsUW | `.FinalOccupationList` | ASM-FW-GISFW-Data-Occupation | 2 |
| InputDtlObject_FacIn_IsUW | `.OfferFacIn.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| InputDtlObject_FacIn_IsUW | `.OfferFacIn.ConveyanceList` | ASM-FW-GISFW-Data-Conveyance | 1 |
| InputDtlObject_FacIn_IsUW | `.OfferFacIn.GoodsList` | ASM-FW-GISFW-Data-Good | 1 |
| InputDtlObject_FacIn_IsUW | `.OfferFacIn.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 2 |
| InputDtlObject_FacIn_IsUW | `.OfferFacIn.PersonList` | Data-Party-Person | 3 |
| InputDtlObject_FacIn_IsUW | `.OfferFacIn.TradingList` | ASM-FW-GISFW-Data-Trading | 1 |
| InputDtlObject_FacIn_IsUW | `.OfferFacIn.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| InputDtlObject_FacIn_IsUW | `.PersonList` | Data-Party-Person | 4 |
| InputDtlObject_FacIn_IsUW | `.PolicyList` | ASM-FW-GISFW-Data-Policy | 2 |
| InputDtlObject_FacIn_IsUW | `.PropertyList` | ASM-FW-GISFW-Data-Property | 1 |
| InputDtlObject_FacIn_IsUW | `.SubContractList` | ASM-FW-GISFW-Data-SubContract | 1 |
| InputDtlObject_FacIn_IsUW | `pyWorkPage.BatchMemberList` | ASM-FW-GISFW-Data-BatchMember | 2 |
| InputDtlObject_FacIn | `.FinalOccupationList` | ASM-FW-GISFW-Data-Occupation | 2 |
| InputDtlObject_FacIn | `.OfferFacIn.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| InputDtlObject_FacIn | `.OfferFacIn.ConveyanceList` | ASM-FW-GISFW-Data-Conveyance | 1 |
| InputDtlObject_FacIn | `.OfferFacIn.GoodsList` | ASM-FW-GISFW-Data-Good | 1 |
| InputDtlObject_FacIn | `.OfferFacIn.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 2 |
| InputDtlObject_FacIn | `.OfferFacIn.PersonList` | Data-Party-Person | 3 |
| InputDtlObject_FacIn | `.OfferFacIn.TradingList` | ASM-FW-GISFW-Data-Trading | 1 |
| InputDtlObject_FacIn | `.OfferFacIn.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| InputDtlObject_FacIn | `.PersonList` | Data-Party-Person | 4 |
| InputDtlObject_FacIn | `.PolicyList` | ASM-FW-GISFW-Data-Policy | 2 |
| InputDtlObject_FacIn | `.PropertyList` | ASM-FW-GISFW-Data-Property | 1 |
| InputDtlObject_FacIn | `.SubContractList` | ASM-FW-GISFW-Data-SubContract | 1 |
| InputDtlObject_FacIn | `pyWorkPage.BatchMemberList` | ASM-FW-GISFW-Data-BatchMember | 2 |
| InputDtlObjectLocation_FacIn_IsUW | `.Property.RiskLocation.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| InputDtlObjectLocation_FacIn_IsUW | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| InputDtlObjectLocation_FacIn | `.Property.RiskLocation.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| InputDtlObjectLocation_FacIn | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| InputDtlObjectLocation_GISFW | `.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| InputDtlObjectLocationHull_GISFW | `.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| InputDtlObjFire_IsUW | `.FEAList` | ASM-FW-GISFW-Data-FEA | 1 |
| InputDtlObjFire_IsUW | `.ListCauseOfLoss` | ASM-FW-GISFW-Data-CauseOfLoss | 1 |
| InputDtlObjFire_IsUW | `.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| InputDtlObjFire_IsUW | `.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 1 |
| InputDtlObjFire | `.FEAList` | ASM-FW-GISFW-Data-FEA | 1 |
| InputDtlObjFire | `.ListCauseOfLoss` | ASM-FW-GISFW-Data-CauseOfLoss | 1 |
| InputDtlObjFire | `.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| InputDtlObjFire | `.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 1 |
| InputDtlParticipantDM | `.ASMHeir` | Data-Party-Person | 1 |
| InputDtlParticipantLife_FacIn_IsUW | `.ASMHeir` | Data-Party-Person | 1 |
| InputDtlParticipantLife_FacIn | `.ASMHeir` | Data-Party-Person | 1 |
| InputDtlParticipantPA_FacIn_IsUW | `.ASMHeir` | Data-Party-Person | 1 |
| InputDtlParticipantPA_FacIn | `.ASMHeir` | Data-Party-Person | 1 |
| InputDtlParticipantTravel_GISFW | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputDtlParticipantTravel_GISFW | `.ASMHeir` | Data-Party-Person | 1 |
| InputDtlParticipantTravel_IsUW | `.ASMHeir` | Data-Party-Person | 1 |
| InputDtlParticipantTravel | `.ASMHeir` | Data-Party-Person | 1 |
| InputDtlPayment_FacIn_IsUW | `.Policy.Payment.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| InputDtlPayment_FacIn | `.Policy.Payment.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| InputDtlPayment_FacInLife_IsUW | `.Policy.Payment.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| InputDtlPayment_FacInLife | `.Policy.Payment.ListInstallment` | ASM-FW-GISFW-Data-Installment | 1 |
| InputDtlSpreadingCommisionCoverage_FacIn | `.SpreadingList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| InputDtlSpreadingCoverage_FacIn_IsUW | `.SpreadingList` | ASM-FW-GISFW-Data-SpreadingRisk | 3 |
| InputDtlSpreadingCoverage_FacIn | `.SpreadingList` | ASM-FW-GISFW-Data-SpreadingRisk | 3 |
| InputEndorsementDtl_IsUW | `.OfferFacIn.CedingCedantList` | ASM-FW-GISFW-Data-Quotation | 1 |
| InputEndorsementDtl_IsUW | `.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputEndorsementDtl_IsUW | `SpreadingList.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| InputEndorsementDtl_IsUW | `TotalSpreadAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputEndorsementDtl | `.OfferFacIn.CedingCedantList` | ASM-FW-GISFW-Data-Quotation | 1 |
| InputEndorsementDtl | `.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputEndorsementDtl | `SpreadingList.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 2 |
| InputEndorsementDtl | `TotalSpreadAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputFacOffer_IsUW | `.OfferFacIn.FacRetroList` | ASM-FW-GISFW-Data-FacOffer | 2 |
| InputFacOffer_IsUW | `CurrencyListRetroAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputFacOffer_IsUW | `GrandTotalListRetroAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| InputFacOffer | `.OfferFacIn.FacRetroList` | ASM-FW-GISFW-Data-FacOffer | 2 |
| InputFacOffer | `CurrencyListRetroAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputFacOffer | `GrandTotalListRetroAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| InputFacOffer | `TempViewSuggest.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | 1 |
| InputInwardFacultativeDtl_IsUW | `.OfferFacIn.CedingCedantList` | ASM-FW-GISFW-Data-Quotation | 1 |
| InputInwardFacultativeDtl_IsUW | `.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputInwardFacultativeDtl_IsUW | `SpreadingList.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| InputInwardFacultativeDtl_IsUW | `TotalSpreadAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputInwardFacultativeDtl_IsUW | `TotalSpreadingCurrency.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 2 |
| InputInwardFacultativeDtl | `.OfferFacIn.CedingCedantList` | ASM-FW-GISFW-Data-Quotation | 1 |
| InputInwardFacultativeDtl | `.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputInwardFacultativeDtl | `SpreadingList.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 2 |
| InputInwardFacultativeDtl | `TotalSpreadAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| InputInwardFacultativeDtl | `TotalSpreadingCurrency.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 3 |
| InputInwardFacultativeSuggest | `.ViewSuggest` | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | 1 |
| InputObjectDtlFire_GCNM | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputObjectDtlFire_GCNM | `.FEAList` | ASM-FW-GISFW-Data-FEA | 1 |
| InputObjectDtlFire_GCNM | `.ListCauseOfLoss` | ASM-FW-GISFW-Data-CauseOfLoss | 1 |
| InputObjectDtlFire_GCNM | `.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| InputObjectDtlFire_GCNM | `.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 2 |
| InputObjectSubContract_FacIn_IsUW | `.OfferFacIn.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| InputObjectSubContract_FacIn_IsUW | `.SubContractList` | ASM-FW-GISFW-Data-SubContract | 1 |
| InputObjectSubContract_FacIn | `.OfferFacIn.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| InputObjectSubContract_FacIn | `.SubContractList` | ASM-FW-GISFW-Data-SubContract | 1 |
| InputObjItem | `.StockAdjustmentList` | ASM-FW-GISFW-Data-StockDeclaration | 1 |
| InputOkupasi | `.SubjectToList` | ASM-FW-GISFW-Data-SubjectTo | 1 |
| InputOkupasiAneka_FacIn_IsUW | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 3 |
| InputOkupasiAneka_FacIn | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 3 |
| InputOkupasiAneka_GCNM | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| InputOkupasiHull_GCNM | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| InputOtherSchedule | `.AnekaAnotherScheduleList` | ASM-FW-GISFW-Data-AnekaOtherSchedule | 1 |
| InputPerson | `.PersonList` | Data-Party-Person | 1 |
| InputSpreadingLife_FacIn_IsUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputSpreadingLife_FacIn | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputSpreadingPA_FacIn_IsUW | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputSpreadingPA_FacIn | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 2 |
| InputSpreadingTravel_FacIn_IsUW | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| InputSpreadingTravel_FacIn | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| InstallmentList | `.InstallmentList` | ASM-FW-GISFW-Data-Installment | 1 |
| Installments_ReadOnly | `.InstallmentList` | ASM-FW-GISFW-Data-TreatyInInstallment | 1 |
| InsuredNameTable | `pgRepPgSubSectionInsuredNameTableBB.pxResults` | ASM-FW-GISFW-Int-CLIENT | 1 |
| InwardFacIn | `pyWorkPage.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| LimitTreaty | `OfferJson.pxResults` | ASM-FW-GISFW-Int-OFFERJSON | 1 |
| ListPayment_SC | `pyWorkPage.PaymentData.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| ListSuggest | `.SuggestList` | ASM-FW-GISFW-Data-SuggestList | 1 |
| LocationDetail | `.Property.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| LocationDetail | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 1 |
| LocationDetail | `D_CoverageSummary.pxResults` | ASM-FW-GISFW-Data-Coverage | 1 |
| LocationPropertyFacOut | `ListLocation.Object` | ASM-FW-GISFW-Data-PrintMemoPlacingData | 1 |
| MarineCargoDtl_IsUW | `.PolicyData.Ship.AdditionalShip` | ASM-FW-GISFW-Int-MSHIP | 1 |
| MarineCargoDtl | `.PolicyData.Ship.AdditionalShip` | ASM-FW-GISFW-Int-MSHIP | 1 |
| ObjectDtlAneka_FacIn | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 4 |
| ObjectItemSummary | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 1 |
| ObjectList_IsUW | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| ObjectList_LossRecord | `.OfferFacIn.ListCauseOfLoss` | ASM-FW-GISFW-Data-CauseOfLoss | 1 |
| ObjectList | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 1 |
| ObjectOccupation_FacIn | `.Property.RiskLocation.AnekaList` | ASM-FW-GISFW-Data-Aneka | 2 |
| ObjectOccupation_FacIn | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 2 |
| ObjectOccupation_FacIn | `.Property.TotalTSIPremiGrossList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| OccupationDetail | `.Property.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| OccupationFacOut_IsUW | `.Property.RiskLocation.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| OccupationFacOut_IsUW | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| OccupationFacOut | `.Property.RiskLocation.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| OccupationFacOut | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| OccupationList_IsUW | `.Property.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| OccupationList | `.Property.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| OccupationPropertyFacOut | `.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| OfferFacIn_NusaRe_IsUW | `pyWorkPage.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| OfferFacIn_NusaRe | `pyWorkPage.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| OfferStatusFacOut | `.OfferFacIn.FacRetroList` | ASM-FW-GISFW-Data-FacOffer | 1 |
| OfferStatusFacOut | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| OldDataEndorsementDtl | `SpreadingList.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| OldOfferStatusFacOut | `.OfferFacIn.OldData.OldData.FacRetroList` | ASM-FW-GISFW-Data-FacOffer | 1 |
| OldOfferStatusFacOut | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| OldOfferStatusFacOut | `TempOldData.OfferFacIn.FacRetroList` | ASM-FW-GISFW-Data-FacOffer | 2 |
| PASummarySection | `D_PersonListSummary.pxResults` | Data-Party-Person | 1 |
| PaymentCurrencyList_IsUW | `pyWorkPage.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 3 |
| PaymentCurrencyList | `pyWorkPage.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 3 |
| PaymentCurrencyListLife_IsUW | `pyWorkPage.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 3 |
| PaymentCurrencyListLife | `pyWorkPage.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 3 |
| PersonCoverage | `.ASMClause` | ASM-FW-GISFW-Data-Clause | 1 |
| PersonCoverage | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| PersonCoverage | `.ASMHeir` | Data-Party-Person | 1 |
| PersonCoverage | `.ASMNotes` | Data-WorkAttach-Note | 1 |
| PersonCoverage | `.ASMNotesDokter` | Data-WorkAttach-Note | 1 |
| PersonCoverageRO | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| PersonListDetail | `D_PersonListCoverageSummary.pxResults` | ASM-FW-GISFW-Data-Coverage | 1 |
| PersonPlanAndDeduct | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 1 |
| PlanContains_UWRO | `.BenefitList` | ASM-FW-GISFW-Data-Benefit | 1 |
| PlanContains_UWRO | `.ClauseList` | ASM-FW-GISFW-Data-Clause | 1 |
| PlanContains_UWRO | `.PremiumList` | ASM-FW-GISFW-Data-Plan-PremiumList | 1 |
| PlanListRO | `.PlanList` | ASM-FW-GISFW-Data-Plan | 1 |
| PrintRISlip | `.OfferFacIn.FacRetroList` | ASM-FW-GISFW-Data-FacOffer | 3 |
| PrintRISlip | `TempViewSuggest.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | 1 |
| PropertyItemList_IsUW | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 2 |
| PropertyItemList_IsUW | `.Property.TotalTSIList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| PropertyItemList | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 2 |
| PropertyItemList | `.Property.TotalTSIList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| PropertyItemListCoverage_IsUW | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 2 |
| PropertyItemListCoverage_IsUW | `.Property.TotalTSIPremiGrossList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| PropertyItemListCoverage | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 2 |
| PropertyItemListCoverage | `.Property.TotalTSIPremiGrossList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| PropertyItemListCoverageCommision | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 2 |
| PropertyItemListCoverageCommision | `.Property.RiskLocation.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| PropertyItemListCoverageCommision | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| PropertyItemListCoverageSpreading_IsUW | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 2 |
| PropertyItemListCoverageSpreading_IsUW | `.Property.RiskLocation.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| PropertyItemListCoverageSpreading_IsUW | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| PropertyItemListCoverageSpreading_IsUW | `.Property.TotalTSIPremiSpreadRNM` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| PropertyItemListCoverageSpreading_IsUW | `SpreadingListLoc.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| PropertyItemListCoverageSpreading | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 4 |
| PropertyItemListCoverageSpreading | `.Property.RiskLocation.AnekaList` | ASM-FW-GISFW-Data-Aneka | 2 |
| PropertyItemListCoverageSpreading | `.Property.RiskLocation.OccupationList` | ASM-FW-GISFW-Data-Occupation | 4 |
| PropertyItemListCoverageSpreading | `.Property.TotalTSIPremiSpreadRNM` | ASM-FW-GISFW-Data-SpreadingRisk | 2 |
| PropertyItemListFacOut_IsUW | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 1 |
| PropertyItemListFacOut | `.Property.PropertyItemList` | ASM-FW-GISFW-Data-PropertyItem | 1 |
| PropertyItemListFacOut | `.Property.TotalTSIPremiRetro` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| ReasViewAttachment | `AttachShowList.pxResults` | Data-WorkAttach-File | 1 |
| ReasViewAttachment | `pgRepPgSubSectionReasViewAttachmentBB.pxResults` | ASM-FW-GISFW-Int-DOCUMENT_POLIS | 1 |
| RemarksFacLetter | `.FacLetterList` | ASM-FW-GISFW-Data-FacLetter | 1 |
| SearchRiskAccumCov | `pgRepPgSubSectionSearchRiskAccumCovBBBBBBBBBBBBBB.pxResults` | ASM-FW-GISFW-Int-ACCUMULATION | 1 |
| SearchRiskAccumCov | `ResultsAccumulation.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| SelClauseList | `.ClauseList` | ASM-FW-GISFW-Data-Clause | 1 |
| SelectAgent | `pgRepPgSubSectionSelectAgentBB.pxResults` | ASM-FW-GISFW-Int-lloydagent | 1 |
| SelectCoverage | `pgRepPgSubSectionSelectCoverageBB.pxResults` | ASM-FW-GISFW-Int-CONDITION | 1 |
| SelectCoverage | `TempMasterPolis.OfferFacIn.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| SelectPlanListRO | `.PlanList` | ASM-FW-GISFW-Data-Plan | 1 |
| SelectShip | `pgRepPgSubSectionSelectShipBB.pxResults` | ASM-FW-GISFW-Int-SHIP | 1 |
| SelectWarrantyList | `.WarrantyList` | ASM-FW-GISFW-Data-Warranty | 1 |
| SFAPortal_OpportunitiesList | `pgRepPgSubSectionSFAPortal_OpportunitiesListBBBBBBBBBBBB.pxResults` | ASM-FW-SFAGISFW-Work-Opportunity | 1 |
| SFAPortal_OpportunitiesList | `pgRepPgSubSectionSFAPortal_OpportunitiesListBBBBBBBBBBBBBB.pxResults` | ASM-FW-SFAGISFW-Work-Opportunity | 1 |
| SFAPortal_OpportunitiesList | `pgRepPgSubSectionSFAPortal_OpportunitiesListBBBBBBBBBBBBBBBB.pxResults` | ASM-FW-SFAGISFW-Work-Opportunity | 1 |
| ShowCedingCoList | `pyWorkPage.Quotation.CedingCoList` | ASM-FW-GISFW-Data-Quotation | 1 |
| ShowObjectFacOut_IsUW | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| ShowObjectFacOut_IsUW | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 2 |
| ShowObjectFacOut_IsUW | `.PersonList` | Data-Party-Person | 2 |
| ShowObjectFacOut_IsUW | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| ShowObjectFacOut | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| ShowObjectFacOut | `.LocationList` | ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance | 2 |
| ShowObjectFacOut | `.PersonList` | Data-Party-Person | 3 |
| ShowObjectFacOut | `.PropertyList` | ASM-FW-GISFW-Data-Property | 1 |
| ShowObjectFacOut | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| SourceHierarki | `pgRepPgSubSectionSourceHierarkiBBB.pxResults` | ASM-FW-GISFW-Int-AGENT | 1 |
| SpreadingCoverageAneka_IsUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| SpreadingCoverageAneka | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| SpreadingItem_IsUW | `.LayerList` | ASM-FW-GISFW-Data-Coverage | 1 |
| SpreadingItem | `.LayerList` | ASM-FW-GISFW-Data-Coverage | 1 |
| SpreadingObjectAneka_IsUW | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| SpreadingObjectAneka | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| SpreadingRiskList | `.SpreadingRiskList` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| SummaryCoverage_Section | `TempTotal.pxResults` | ASM-FW-GISFW-Data-PropertyItem | 1 |
| SummaryCoverage_Section | `TempTotalItem.pxResults` | ASM-FW-GISFW-Data-PropertyItem | 1 |
| SummaryLossRecord_Section | `TotalLossCOB.pxResults` | ASM-FW-GISFW-Data-CauseOfLoss | 1 |
| SummaryRiskAccumulation | `pgRepPgSubSectionSummaryRiskAccumulationBBBB.pxResults` | ASM-FW-GISFW-Int-ACCUMULATION | 1 |
| SummaryRiskAccumulation | `ResultsAccumulation.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| SummaryRiskAccumulation | `SummaryRiskAccum.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| SummarySpreading_Section | `TotalSpreadAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| SummarySpreading_Section | `TotalSpreadingCurrency.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 2 |
| TableInwardScale | `D_InwardScale.pxResults` | ASM-FW-GISFW-Int-INWARDSCALE | 1 |
| TableOfLimit | `pgRepPgSubSectionTableOfLimitBBB.pxResults` | ASM-FW-GISFW-Int-TABLEOFLIMIT | 1 |
| TotalAccumulationDtl | `TotalAccumulationDtl.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| TSIForFacOut | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| TSIPropertyFacOut | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| UploadMember | `.BatchMemberList` | ASM-FW-GISFW-Data-BatchMember | 1 |
| UploadMember | `.FinalOccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| UploadMember | `.PolicyList` | ASM-FW-GISFW-Data-Policy | 1 |
| VehicleGrid_FacIn_IsUW | `.AccessoryList` | ASM-FW-GISFW-Data-Accessory | 1 |
| VehicleGrid_FacIn | `.AccessoryList` | ASM-FW-GISFW-Data-Accessory | 1 |
| VehicleGrid | `.AccessoryList` | ASM-FW-GISFW-Data-Accessory | 1 |
| ViewCheckListOffer | `.Policy.FacOfferList` | ASM-FW-GISFW-Data-FacOffer | 1 |
| ViewCheckListOffer | `.PropertyList` | ASM-FW-GISFW-Data-Property | 4 |
| ViewCheckListOffer | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 2 |
| ViewCheckListOfferFacOut | `.Policy.FacOfferList` | ASM-FW-GISFW-Data-FacOffer | 1 |
| ViewCheckListOfferFacOut | `.PropertyList` | ASM-FW-GISFW-Data-Property | 4 |
| ViewCheckListOfferFacOut | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 2 |
| ViewClaimList | `TempListClaim.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| ViewClauseFire | `.ArgumentList` | ASM-FW-GISFW-Data-Argument | 1 |
| ViewCoverage | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| ViewCoverageFacOutShow_isUW | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| ViewCoverageFacOutShow | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| ViewCoverageFire | `.CoverageList` | ASM-FW-GISFW-Data-Coverage | 1 |
| ViewCoveragePA | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| ViewCoverageTravel | `.ASMCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| ViewCSVAneka | `TempWorkPage.ListB2BHostUpload` | ASM-FW-GISFW-Int-TABLE_B2B | 1 |
| ViewCSVCredit | `TempWorkPage.ListB2BHostUpload` | ASM-FW-GISFW-Int-TABLE_B2B | 1 |
| ViewCSVMarine | `TempWorkPage.ListB2BHostUpload` | ASM-FW-GISFW-Int-TABLE_B2B | 1 |
| ViewCSVResult_TableB2B | `TempWorkPage.ListB2BHostUpload` | ASM-FW-GISFW-Int-TABLE_B2B | 1 |
| ViewDataPolicy | `TotalSpreadAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 1 |
| ViewDetailPayment_SC | `.DetailPayment` | ASM-FW-GISFW-Data-Payment | 1 |
| ViewDtlAddCargo | `.AdditionalCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| ViewDtlAdditional | `.AdditionalCoverage` | ASM-FW-GISFW-Data-Coverage | 1 |
| ViewDtlCoverageFacOut | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| ViewDtlCoverageFacOut | `.LocationList` | Data-Address | 1 |
| ViewDtlCoverageFacOut | `.PersonList` | Data-Party-Person | 2 |
| ViewDtlCoverageFacOut | `.PropertyList` | ASM-FW-GISFW-Data-Property | 1 |
| ViewDtlCoverageFacOut | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| ViewDtlDeductible | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 1 |
| ViewDtlObjectLocation | `.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| ViewDtlOutgo | `.OutgoList` | ASM-FW-GISFW-Data-Outgo | 1 |
| ViewDtlParticipantPA | `.ASMHeir` | Data-Party-Person | 1 |
| ViewGeneralPolisFacOut | `.Policy.CIFData.AddressList` | Data-Address | 1 |
| ViewIndemnity_Section | `ViewIndemnity.pxResults` | ASM-FW-GISFW-Data-Search | 1 |
| ViewInwardFacultativeDtl_IsUW | `.OfferFacIn.CedingCedantList` | ASM-FW-GISFW-Data-Quotation | 1 |
| ViewInwardFacultativeDtl_IsUW | `.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| ViewInwardFacultativeDtl_IsUW | `SpreadingList.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 1 |
| ViewInwardFacultativeDtl_IsUW | `TotalSpreadAll.pxResults` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 4 |
| ViewInwardFacultativeDtl_IsUW | `TotalSpreadingCurrency.pxResults` | ASM-FW-GISFW-Data-SpreadingRisk | 4 |
| ViewObjectAnekaFacOut | `.LocationList` | Data-Address | 1 |
| ViewObjectDtlAneka | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 1 |
| ViewObjectFireFacOut | `.PropertyList` | ASM-FW-GISFW-Data-Property | 1 |
| ViewObjectMarineCargoFacOut | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| ViewObjectMarineHullFacOut | `.LocationList` | Data-Address | 1 |
| ViewObjectOccupation | `.OccupationList` | ASM-FW-GISFW-Data-Occupation | 1 |
| ViewObjectPAFacOut | `.PersonList` | Data-Party-Person | 1 |
| ViewObjectSubContract | `.LocationList` | Data-Address | 1 |
| ViewObjectSubContract | `.SubContractList` | ASM-FW-GISFW-Data-SubContract | 1 |
| ViewObjectTravelFacOut | `.PersonList` | Data-Party-Person | 1 |
| ViewObjectVehicleFacOut | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| ViewOkupasiAneka | `.AnekaList` | ASM-FW-GISFW-Data-Aneka | 2 |
| ViewOldDataPayment_IsUW | `.OfferFacIn.CurrencyList` | ASM-FW-GISFW-Data-OfferFacIn-Currency | 2 |
| ViewOldDeductible | `.DeductibleList` | ASM-FW-GISFW-Data-Deductible | 1 |
| ViewOpenFollowingPolicy | `Attachment.pxResults` | ASM-FW-GISFW-Int-OFFERJSON | 1 |
| ViewOpenFollowingPolicy | `tempWorkPageView.OfferFacIn.ViewSuggest` | ASM-FW-GISFW-Data-OfferFacIn-SuggestList | 1 |
| ViewPerilsFacOut | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| ViewPolicyListRO | `.PersonList` | Data-Party-Person | 1 |
| ViewTabGroupPolisFacOutSavior | `.CargoList` | ASM-FW-GISFW-Data-Cargo | 1 |
| ViewTabGroupPolisFacOutSavior | `.ClauseList` | ASM-FW-GISFW-Data-Clause | 3 |
| ViewTabGroupPolisFacOutSavior | `.CoinsList` | ASM-FW-GISFW-Data-Coins | 1 |
| ViewTabGroupPolisFacOutSavior | `.LocationList` | Data-Address | 1 |
| ViewTabGroupPolisFacOutSavior | `.PersonList` | Data-Party-Person | 2 |
| ViewTabGroupPolisFacOutSavior | `.PropertyList` | ASM-FW-GISFW-Data-Property | 1 |
| ViewTabGroupPolisFacOutSavior | `.VehicleList` | ASM-FW-GISFW-Data-Vehicle | 1 |
| ViewVehicleGridFacOut | `.AccessoryList` | ASM-FW-GISFW-Data-Accessory | 1 |
| WarrantyList | `pgRepPgSubSectionWarrantyListB.pxResults` | ASM-FW-GISFW-Int-WARRANTY | 1 |

---

*Akhir dokumen. Seluruh angka dihasilkan pada 2026-09-15 dari `D:\migrasi\RNM\NB FacIn\` saja.
Korpus tidak diubah.*