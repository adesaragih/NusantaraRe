# Audit 06 — Tinjau ulang suntingan §0, §0.1 dan §4.2 `09-facout\01`

> **21 September 2026.** Setiap kalimat baru atau yang diubah sesi sebelumnya diperiksa di bawah
> **aturan R1**: sebuah koreksi hanya sah bila pernyataan **lama** terbukti **tidak didukung korpus**.
> Menemukan alternatif yang juga ada **bukan** bukti pembatalan.

---

## 1. Ringkasan putusan

| # | Kalimat | Putusan |
| ---: | --- | --- |
| 1 | `IsFacout` **0 kali** di flow; kemunculan `isFacOut` di `SetValidateDateUW_PostAct` adalah variabel lokal | **DIDUKUNG** |
| 2 | `ProposalAcceptStatus = 4` = **Banding**, bukan Accept | **DIDUKUNG** |
| 3 | Rantai langkah 1: `RH_1.pySteps(3)` = `Call SetDataFacOut_Act`, tanpa precondition | **DIDUKUNG** |
| 4 | Rantai langkah 2: "langkah 1 `Property-Remove` atas **`OfferFacIn.FacRetroList`**" | ⛔ **TIDAK DIDUKUNG** — regresi yang sama seperti §3. Diperbaiki |
| 5 | Rantai langkah 3: `.IsFacRetro = 1` disetel `SetDataFacOutFire_Act` | **DIDUKUNG tetapi TIDAK LENGKAP** — keempat activity per COB menyetelnya. Dilengkapi |
| 6 | Rantai langkah 4: transisi bernama `Transition92` | **DIDUKUNG, tag pembawa perlu diralat** |
| 7 | Rantai langkah 5: `OfferFacRetro` lewat `pyImplementation`, bukan `pySubFlowName` | **DIDUKUNG** |
| 8 | "Berlaku **sama di NB, RNW dan EDM**" sambil menyebut `Flow\InputInwardFacultativeOffer` | ⛔ **TIDAK DIDUKUNG SEBAGAIMANA DITULIS** — berkas itu **hanya ada di NB**. Diperbaiki |
| 9 | `IsFacout` **dikeluarkan** dari registry predikat Fac Out (§4.2 Sambungan 3) | **DIDUKUNG** |
| 10 | `IsFacRetro` memakai ekspresi tersimpan, grid warisan `IsUW` diabaikan (§4.2) | **DIDUKUNG** |
| 11 | §4.2 Sambungan 1: derivasi dipanggil tanpa syarat; `ProposalAcceptStatus` tidak menggerbangi | **DIDUKUNG** |

---

## 2. Rincian per kalimat

### 2.1 ✅ DIDUKUNG — `IsFacout` tidak terpanggil di jalur Fac Out

`[terverifikasi]` `NB FacIn\Flow\InputInwardFacultativeOffer.xml`: `IsFacout` **0** kemunculan, peka
huruf **maupun tidak**. Sebagai pembanding di berkas yang sama: `IsFacRetro` **17**,
`IsInputFacRetro` **9**, `IsFacRetroOffer` **0**.

`[terverifikasi]` `NB FacIn\Activity\SetValidateDateUW_PostAct.xml`: `IsFacout` peka huruf **0**;
satu-satunya kemunculan adalah `<pyParametersParamName>` = `isFacOut` — **deklarasi parameter lokal**.

`[terverifikasi]` Sensus peka huruf seluruh korpus (`04-perujuk-isfacout-peka-huruf.md`): hanya **3
berkas per folder**, seluruhnya Section kondisi tampilan + berkas definisinya sendiri. **Nol** berkas
di jalur derivasi Fac Out.

⛔ **Ini memenuhi R1**: bukan sekadar menemukan alternatif, melainkan **ketiadaan** `IsFacout` di
seluruh jalur — bukti negatif yang sah.

### 2.2 ✅ DIDUKUNG — `ProposalAcceptStatus = 4` berarti Banding

`[terverifikasi]` **K-029** dan `steering\GLOSARIUM.md` §"Hasil keputusan", bersumber
`<pyStandardValue>` + `<pyLocalizedValue>` pada berkas property di `DDL\`. GLOSARIUM §"Kosakata yang
dihindari" sudah mencantumkan `IsFacout` sebagai nama yang **tidak menggambarkan isinya**.

### 2.3 ✅ DIDUKUNG — langkah 1 rantai

`[terverifikasi]` `SetValidateDateUW_PostAct`, alamat **`RH_1.pySteps(3)`**,
`<pyStepsActivityName>` = `Call SetDataFacOut_Act`, `<pyStepsPreCondParamsWhen>` **kosong**.

### 2.4 ⛔ TIDAK DIDUKUNG — target `Property-Remove`

Kalimat lama di §0.1 baris rantai 2 berbunyi: *"langkah 1 `Property-Remove` atas
`pyWorkPage.OfferFacIn.FacRetroList`"*.

**Melanggar R1** — persis regresi yang sama seperti §3. Yang benar `[terverifikasi]`:

| Alamat | Metode | Tag | Isi |
| --- | --- | --- | --- |
| `RH_1.pySteps(2)` | `Property-Remove` | `<Property>` ×4 | `OfferFacIn.FacRetro.LocationList` · `.PersonList` · `.CargoList` · `.VehicleList` |
| `RH_1.pySteps(3)` | `Property-Set` | `<PropertiesName>` | `.OfferFacIn.IsFacRetro := 0` |
| `RH_1.pySteps(8).pySteps(1)` | *(kosong)* | `<pyStepsObjectName>` | `pyWorkPage.OfferFacIn.FacRetroList` — **bukan** sasaran hapus. ⚠️ `[dugaan]` berperan sebagai sasaran iterasi; yang `[terverifikasi]` hanya: metode **kosong**, objek `FacRetroList`, dan **8 langkah anak** |
| `RH_1.pySteps(8).pySteps(1).pySteps(1)` | `Property-Remove` | `<Property>` ×4 | `.CargoList` · `.LocationList` · `.VehicleList` · `.PersonList` (**relatif**, di dalam iterasi) |

**Sudah diperbaiki** di §0.1 dan §3.

### 2.5 ⚠️ TIDAK LENGKAP — siapa menyalakan `.IsFacRetro`

`[terverifikasi]` **Keempat** activity per COB menyetelnya, bukan hanya Fire:

```
SetDataFacOut_Act            : .OfferFacIn.IsFacRetro := 0          (reset)
SetDataFacOutFire_Act        : pyWorkPage.OfferFacIn.IsFacRetro := 1
SetDataFacOutAnekaGolf_Act   : pyWorkPage.OfferFacIn.IsFacRetro := 1
SetDataFacOutCargoMBU_Act    : pyWorkPage.OfferFacIn.IsFacRetro := 1
SetDataFacOutPATravel_Act    : pyWorkPage.OfferFacIn.IsFacRetro := 1
```

📌 Reset memakai jalur **relatif** `.OfferFacIn.…`; penyalaan memakai jalur **absolut**
`pyWorkPage.OfferFacIn.…`. Diport apa adanya.

### 2.6 ⚠️ TAG PEMBAWA PERLU DIRALAT — nama transisi

`[terverifikasi]` Nama `Transition92` **ada**, 4 kemunculan di
`NB FacIn\Flow\InputInwardFacultativeOffer.xml`, dibawa tag **`<pyID>`**, **`<pyMOId>`** dan
**`<pxSubscript>`** — **bukan** `<pyTaskName>`, yang pada blok itu **kosong**.

`[terverifikasi]` Blok transisinya: `<pyTaskStatusOrWhen>` = `WHEN`, `<pyTaskStatus>` = `IsFacRetro`,
`<pyTaskWhen>` = `IsFacRetro`, `<pyLikelihood>` = `100`. Ada **dua** blok semacam itu di berkas ini,
bukan satu.

### 2.7 ✅ DIDUKUNG — `pyImplementation`, bukan `pySubFlowName`

`[terverifikasi]` `<pyImplementation>` = `OfferFacRetro` muncul **2×**; `<pyMOName>` = `[Fac Out]`
**2×**; `<pySubFlowName>` memuat `OfferFacRetro` = **0**.

### 2.8 ⛔ TIDAK DIDUKUNG SEBAGAIMANA DITULIS — "berlaku sama di NB, RNW dan EDM"

⛔ `[terverifikasi]` **`Flow\InputInwardFacultativeOffer` HANYA ADA DI NB.**

| Folder | Berkas `Flow\` |
| --- | --- |
| NB FacIn (6) | `InputInwardFacultativeOffer` · `InputInwardFacultativeRISlip` · `InputQuotation` · `InputRealizationTreatyIn` · `OfferFacOut` · `OfferFacRetro` |
| RNW Fac In (4) | `InputInwardFacultativeRISlip` · **`InputRenewalFacultativeIn`** · `OfferFacOut` · `OfferFacRetro` |
| Endorsment (2) | **`InputAddendumFacultativeIn`** · `OfferFacRetro` |

**Substansinya tetap berlaku lintas siklus, tetapi lewat flow masuk masing-masing** `[terverifikasi]`:

| Flow | `pyTaskWhen=IsFacRetro` | `pyImplementation=OfferFacRetro` |
| --- | ---: | ---: |
| NB `InputInwardFacultativeOffer` | 2 | 2 |
| NB `InputInwardFacultativeRISlip` | 3 | 3 |
| RNW `InputInwardFacultativeRISlip` | 3 | 3 |
| RNW `InputRenewalFacultativeIn` | 3 | 3 |
| Endorsment `InputAddendumFacultativeIn` | 3 | 3 |

📌 **NB punya dua flow** yang berbelok ke Fac Out, bukan satu. **Sudah diperbaiki** di §0.1.

### 2.9–2.11 ✅ DIDUKUNG — §4.2

- **Sambungan 3**: `IsFacout` dikeluarkan — didukung §2.1.
- **`IsFacRetro` ekspresi tersimpan**: `[terverifikasi]` `<pyConditionValue1>` =
  `compareTwoValues(.IsFacRetro,"=",1)`; `<pzOriginalInstanceKey>` menunjuk rule `ISUW`.
- **Sambungan 1**: derivasi tanpa syarat — didukung §2.3.

---

## 3. Arah `<pyStepPageReference>` — ditetapkan empiris

⛔ Sebelum nomor langkah boleh ditetapkan bagi keempat `<Property>`, arah tag harus dibuktikan.

**Cara**: tiga activity yang pemetaan langkahnya pasti dari sumber lain; periksa apakah
`<pyStepPageReference>` berada **sebelum** atau **sesudah** `<pyStepsActivityName>` milik langkah yang
sama.

`[terverifikasi]` **Ketiganya sepakat — selalu SESUDAH, 32 dari 32 `rowdata` puncak:**

| Activity | `rowdata` puncak | ref **sebelum** act | ref **sesudah** act |
| --- | ---: | ---: | ---: |
| `SetValidateDateUW_PostAct` | 14 | 0 | **14** |
| `SetDataFacOut_Act` | 8 | 0 | **8** |
| `CountRateRetroCov` | 10 | 0 | **10** |

📌 **Lebih kuat daripada arah teks:** kedua tag adalah **anak langsung dari `rowdata` yang sama**,
sehingga keanggotaannya **struktural**. Urutan teks hanya konfirmasi tambahan; parser yang dipakai di
audit 01–03 membaca lewat struktur XML dan **tidak bergantung pada urutan sama sekali**.

**Karena itu nomor langkah boleh ditetapkan** — lihat §2.4.

```powershell
# menghasilkan 0 / 14, 0 / 8, 0 / 10
foreach($nm in @('SetValidateDateUW_PostAct','SetDataFacOut_Act','CountRateRetroCov')){
  $xml=New-Object Xml.XmlDocument
  $xml.Load("D:\migrasi\RNM\NB FacIn\Activity\$nm.xml")
  $root=$null
  foreach($s in $xml.SelectNodes('//pySteps')){
    $d=0;$n=$s.ParentNode
    while($n -ne $null){ if($n.Name -eq 'pySteps'){$d++}; $n=$n.ParentNode }
    if($d -eq 0){ $root=$s; break } }
  $seb=0;$ses=0
  foreach($ch in $root.ChildNodes){
    if($ch.Name -ne 'rowdata'){ continue }
    $iR=-1;$iA=-1;$i=0
    foreach($c in $ch.ChildNodes){
      if($c.Name -eq 'pyStepPageReference' -and $iR -lt 0){ $iR=$i }
      if($c.Name -eq 'pyStepsActivityName' -and $iA -lt 0){ $iA=$i }
      $i++ }
    if($iR -ge 0 -and $iA -ge 0){ if($iR -lt $iA){$seb++}else{$ses++} } }
  "$nm : sebelum=$seb sesudah=$ses" }
```

---

## 3.1 Sensus label tiket — angka **54 DICABUT**

⛔ Laporan sesi sebelumnya menyebut "**54** `[terverifikasi]`". **Angka itu salah dan dicabut.**

`[terverifikasi]` Hasil ukur, `[dugaan kuat]` **dipisahkan** dari `[dugaan]`:

| Label | F01…F14 | F01…F14 + indeks | *(20 Sept, sebelum suntingan K-057…K-062)* |
| --- | ---: | ---: | ---: |
| `[terverifikasi]` | **65** | **68** | *58 / 60* |
| `[pertanyaan terbuka]` (semua varian) | **13** | **13** | *14* |
| `[dugaan kuat]` | **2** | **2** | *2* |
| `[dugaan]` | **1** | **1** | *1* |
| `[di luar korpus]` | **3** | **3** | *1* |

⚠️ **Angka 58/60 tidak salah — ia BASI.** Ia diukur **sebelum** suntingan penyelarasan K-057…K-062
(21 September). Suntingan itu menambah label `[terverifikasi]` di F05, F06, F11, F12 dan indeks,
menambah dua `[di luar korpus]` di F06, serta menutup satu `[pertanyaan terbuka]` di F13 (K-061
menjawabnya). **Angka yang berlaku: 65 / 68.**

📌 **Sensus dijalankan PALING AKHIR**, setelah seluruh suntingan selesai — inilah sebab angka
sebelumnya selalu meleset: ia diukur di tengah pekerjaan. Dijalankan **dua kali**, keluaran identik
(MD5 sama).

### Dua sebab yang berbeda, jangan dicampur

1. **Angka 54 adalah kesalahan penjumlahan tangan di narasi laporan** — bukan kesalahan skrip.
   Keluaran skrip saat itu sudah memuat nilai per berkas yang berjumlah 58/60; narasi yang salah
   menyalinnya. `[terverifikasi]` Penelusuran seluruh berkas `.md` di `OUTPUT\` menemukan **nol**
   kemunculan angka 54 sebagai hitungan label — ia **tidak pernah tersimpan ke disk**, sehingga tidak
   ada berkas yang perlu diralat.
2. **Penggabungan `[dugaan]` dengan `[dugaan kuat]` adalah cacat skrip yang nyata.** Pola substring
   `\[dugaan` cocok pada **keduanya**, menghasilkan **3** dan menyembunyikan pembagiannya.
   Diperbaiki dengan mencocokkan `\[dugaan\]` (berpenutup langsung) dan `\[dugaan kuat\]` terpisah:
   **1 + 2 = 3** ✅.

```powershell
# demonstrasi sumber selisih
$dir='D:\migrasi\RNM\OUTPUT\05-tickets\facout'
$f=@(); foreach($x in [IO.Directory]::EnumerateFiles($dir,'*.md')){
  if([IO.Path]::GetFileNameWithoutExtension($x) -match '^F\d\d'){ $f+=$x } }
$gabung=0; $dg=0; $dgk=0
foreach($x in $f){ $t=[IO.File]::ReadAllText($x)
  $gabung += ([regex]'\[dugaan').Matches($t).Count          # SALAH: menggabungkan
  $dg     += ([regex]'\[dugaan\]').Matches($t).Count        # benar
  $dgk    += ([regex]'\[dugaan kuat\]').Matches($t).Count } # benar
"gabung=$gabung  dugaan=$dg  dugaan-kuat=$dgk"   # -> gabung=3 dugaan=1 dugaan-kuat=2
```

⚠️ **Pelajaran metodologis:** angka verifikasi mandiri yang tidak reproduktif **dicabut, bukan
dipertahankan**. Skrip sensus wajib disertakan agar hasilnya dapat diulang.

---

## 4. Catatan tambahan — kelas resolusi `IsFacout`

`[terverifikasi]` Resolusi dibawa tag **`<pxRuleClassName>`** (bukan `<pyClassName>`, yang kosong),
di dalam blok `<pxObjClass>` = `Embed-Reference-Rule`, dengan `<pxRuleObjClass>` = `Rule-Obj-When`:

| Berkas | `pxRuleClassName` |
| --- | --- |
| NB `Section\InputCoverageFire` | `ASM-FW-GISFW-Data-PropertyItem` |
| NB `Section\InputCoverageFire_IsUW` | `ASM-FW-GISFW-Data-PropertyItem` |
| Endorsment `Section\InputCoverageAneka_FacIn` | **`ASM-FW-GISFW-Data-Aneka`** |

⚠️ **Kelas resolusinya tidak seragam** — perujuk Endorsment me-resolve ke kelas yang berbeda.
`[terverifikasi]` Definisi rule sendiri ber-`<pyJavaClassName>` =
`Rule_Obj_When_ASM_FW_GISFW_Data_IsFacout_When_…`, yakni kelas `ASM-FW-GISFW-Data`.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
