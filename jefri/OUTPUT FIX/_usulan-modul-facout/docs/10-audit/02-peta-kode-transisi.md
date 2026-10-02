# Audit 02 — Peta kode transisi `pyStepsPreCondParamsWhenTrue` / `WhenFalse`

> **Diukur ulang 21 September 2026.** Sumber angka bagi `05-tickets\facout\F10-produksi-facoutproduction.md`.

---

## 1. Sensus lengkap

`[terverifikasi]` Ketiga folder korpus, seluruh 14 subfolder:

| Ukuran | Nilai |
| --- | ---: |
| Berkas dipindai | **6.071** |
| Berkas memuat tag gerbang | **1.747** |
| Blok precondition (`pyStepsPreCondParams/rowdata`) | **36.265** |
| Pasangan (T, F) berbeda | **31** |

| T | F | Total | Bergerbang | Tanpa gerbang |
| --- | --- | ---: | ---: | ---: |
| 2 | 2 | 17.115 | 74 | **17.041** |
| 2 | 3 | 10.020 | 9.731 | 289 |
| *(kosong)* | 3 | 5.103 | 5.017 | 86 |
| 3 | 2 | 1.759 | 1.750 | 9 |
| **5** | 2 | 908 | 908 | 0 |
| 3 | *(kosong)* | 298 | 293 | 5 |
| *(kosong)* | *(kosong)* | 203 | 3 | 200 |
| **1** | 2 | 156 | 156 | 0 |
| **6** | 2 | 142 | 139 | 3 |
| 2 | 1 | 101 | 98 | 3 |
| 2 | **4** | 93 | 91 | 2 |
| 2 | 5 | 73 | 73 | 0 |
| 2 | 6 | 57 | 54 | 3 |

*(18 pasangan sisanya masing-masing < 40; seluruhnya tercakup perintah audit §4.)*

📌 **17.041 dari 17.115** pasangan `2/2` **tidak punya gerbang sama sekali** — itu bentuk bawaan
langkah tanpa precondition.

---

## 2. Penambatan — apa yang terbukti dan apa yang tidak

### 2.1 Keanggotaan gerbang dan T/F pada langkah yang sama — `[terverifikasi]`

⛔ Syarat yang dituntut: **buktikan nilai T/F dan gerbangnya milik langkah yang SAMA lewat
`pyStepPageReference`.** Hasilnya:

| Jangkar | Alamat | Metode | Gerbang | T | F |
| --- | --- | --- | --- | ---: | ---: |
| `SaveEDMToJsonPolicy_Act` | `RH_1.pySteps(10)` | *(kosong)* | `@LengthOfPageList(EndorsementFacIn.pxResults)>=1` | **6** | 3 |
| `SaveTreatyProduction_Act` | `RH_1.pySteps(4)` | *(kosong)* | `CekFacin.pxResults(1).CARI1==""` | 2 | **6** |
| `SaveEDMToJsonPolicy_Act` | `RH_1.pySteps(12)` | RDB-List | `IsLife` | **3** | 2 |
| `SaveEDMToJsonPolicy_Act` | `RH_1.pySteps(13)` | RDB-List | `IsLife` | 2 | **3** |
| `InsertFacoutProduction` | `RH_1.pySteps(7)` | *(kosong)* | `@SizeOfPropertyList(FacoutList.pxResults)>0` | **6** | 2 |

`[terverifikasi]` Pada kelimanya, `pyStepsPreCondParamsWhen`, `…WhenTrue` dan `…WhenFalse` berada di
**`rowdata` `pyStepsPreCondParams` yang sama**, dan blok itu adalah **anak langsung** dari `rowdata`
langkah yang membawa `pyStepPageReference` di atas. Keanggotaan ditetapkan **secara struktural**, bukan
dari kedekatan teks.

### 2.2 Arti kodenya — `[dugaan kuat]`, **bukan** `[terverifikasi]`

| Kode | Arti yang diusulkan | Dasar |
| ---: | --- | --- |
| **2** | jalankan langkah ini, lanjut | 17.041 pasangan `2/2` **tanpa gerbang** = bentuk bawaan · jangkar `SaveTreatyProduction_Act` T=2 saat baris belum ada |
| **3** | lewati langkah ini | pasangan terbalik `pySteps(12)` T=3/F=2 vs `pySteps(13)` T=2/F=3 **atas gerbang `IsLife` yang sama** — hanya masuk akal bila satu berjalan saat benar dan satunya saat salah |
| **6** | keluar dari activity | **dua jangkar independen di dua berkas berbeda**, keduanya mengaitkan `6` dengan keadaan "baris/dokumen SUDAH ada" |

⛔ **Mengapa tetap `[dugaan kuat]`:** yang terbukti dari korpus adalah **keanggotaan struktural** dan
**pola distribusinya**. Arti "keluar dari activity" tetap **inferensi** atas apa yang gerbang itu
harus lakukan; korpus tidak memuat berkas yang menyatakan pemetaan kode→perilaku secara eksplisit.
Kenaikan ke `[terverifikasi]` memerlukan konfirmasi UI Pega work owner.

⚠️ Deskripsi langkah pada jangkar pertama berbunyi seperti perintah keluar. **Deskripsi adalah label**
(`PANDUAN-KERJA` §3: nama dan label bukan bukti) — dicatat sebagai penguat, **tidak** dijadikan dasar.

### 2.3 Kode yang TIDAK tertambat

| Kode | Kemunculan | Status |
| ---: | ---: | --- |
| **1** | 156 (T,F=2) + 38 + 38 + 7 + 7 | `[pertanyaan terbuka — butuh UI Pega]` |
| **4** | 93 (T=2) + 14 + 12 + 3 | `[pertanyaan terbuka — butuh UI Pega]` |
| **5** | **908** (T=5,F=2) + 73 + 23 + 3 | `[pertanyaan terbuka — butuh UI Pega]` |

⚠️ Kode **5** adalah yang **paling sering** di antara yang tak tertambat (908 kemunculan pada satu
pasangan saja). Hipotesis "kendali perulangan" **sudah diuji dan gagal**: langkah pembawa kode 5
beriterasi hanya **9,8 %**, lebih rendah daripada kode 2 (30,8 %).

📌 Kode **1** adalah satu-satunya yang membawa sasaran di tag `…Prms` pasangannya — nilai contoh
`Close`, `Exit`, `save`. Strukturnya khas; **artinya tetap tidak tertambat**.

---

## 3. Arah gerbang idempotensi `InsertFacoutProduction`

`[terverifikasi]` Alamat **`RH_1.pySteps(7)`**, gerbang
`@SizeOfPropertyList(FacoutList.pxResults)>0`, **T=6 / F=2**.

`[terverifikasi]` Gerbang itu memakai hasil `RDBList\CekFacoutProd_Sql` (`Rule-Connect-SQL`, tag
`pyBrowseSQL`):

```sql
select DISTINCT IDPEGA from FACoutPRODUCTION
 where RISLIPNO = {DataIN.CARI7} and IDPEGA = {pyWorkPage.pzInsKey}
```

`[dugaan kuat]` Dengan peta §2.2: **baris sudah ada → `6` → keluar activity, tidak menyisipkan;
baris belum ada → `2` → lanjut dan sisipkan.** Rantai **aman diulang**.

⚠️ Bila peta §2.2 kelak terbantah UI Pega, **arah ini ikut gugur** dan F10 harus ditinjau ulang.

---

## 4. Perintah audit

```powershell
# sensus 31 pasangan; menghasilkan 6071 / 1747 / 36265
$folders=@('NB FacIn','RNW Fac In','Endorsment Fac In')
$subs=@('Activity','ConnectREST','DataPage','DataTransform','DecisionTable','DecisionTree','Flow',
        'FlowAction','Harness','RDBList','ReportDefinition','Section','SystemSettings','When')
$pair=@{}; $fTag=0; $blok=0; $fAll=0
foreach($f in $folders){ foreach($s in $subs){
  $p="D:\migrasi\RNM\$f\$s"; if(-not [IO.Directory]::Exists($p)){continue}
  foreach($x in [IO.Directory]::EnumerateFiles($p,'*.xml')){
    $fAll++; $t=[IO.File]::ReadAllText($x)
    if($t.IndexOf('pyStepsPreCondParamsWhenTrue') -lt 0){ continue }
    $fTag++; $d=New-Object Xml.XmlDocument
    try{ $d.LoadXml($t) }catch{ continue }
    foreach($row in $d.SelectNodes('//pyStepsPreCondParams/rowdata')){
      $blok++; $a='';$b=''
      foreach($c in $row.ChildNodes){
        if($c.Name -eq 'pyStepsPreCondParamsWhenTrue'){ $a=$c.InnerText.Trim() }
        if($c.Name -eq 'pyStepsPreCondParamsWhenFalse'){ $b=$c.InnerText.Trim() } }
      if($a -eq ''){$a='(kosong)'}; if($b -eq ''){$b='(kosong)'}
      $k="$a|$b"; if($pair.ContainsKey($k)){$pair[$k]++}else{$pair[$k]=1} } } } }
"berkas=$fAll tag=$fTag blok=$blok pasangan=$($pair.Count)"
```

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
