# Audit 09 — Peta pemanggil `CountRateRetroCov`

> **21 September 2026.** Menjawab: berapa jalur COB yang benar-benar memanggil rumus premi retro, dan
> kelas apa saja yang dirujuk. Sumber bagi **K-060 (amandemen)** dan tiket
> `05-tickets\facout\F06-premi-retro.md`.

---

## 1. Cakupan pencarian

`[terverifikasi]` **Peka huruf**, seluruh korpus + `DDL\`:

| Ukuran | Nilai |
| --- | ---: |
| Berkas dipindai (14 subfolder × 3 folder + `DDL\*.xml`) | **6.105** |
| Berkas menyebut `CountRateRetroCov` | **15** |
| — di antaranya **definisi rule** | **5** |
| — **pemanggil** | **10** |

---

## 2. Definisi rule — **dua**, di kelas berbeda

`[terverifikasi]` Identitas dibaca dari `pzInsKey` + `pyClassName`, **bukan** nama berkas.

| Berkas | `pyClassName` | Basis `pzInsKey` |
| --- | --- | --- |
| `NB FacIn\Activity\CountRateRetroCov.xml` | `ASM-FW-GISFW-Data-PropertyItem` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-PROPERTYITEM COUNTRATERETROCOV` |
| `RNW Fac In\Activity\CountRateRetroCov.xml` | `ASM-FW-GISFW-Data-PropertyItem` | idem |
| `DDL\CountRateRetroCov.xml` | `ASM-FW-GISFW-Data-PropertyItem` | idem ✅ **berlaku** |
| `Endorsment Fac In\Activity\CountRateRetroCov.xml` | `ASM-FW-GISFW-Data-Aneka` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-ANEKA COUNTRATERETROCOV` |
| `DDL\CountRateRetroCov(ANEKA).xml` | `ASM-FW-GISFW-Data-Aneka` | idem ✅ **berlaku** |

⛔ **Tidak ada berkas rule berkelas `ASM-FW-GISFW-Data-Cargo`** di korpus maupun `DDL\`.

---

## 3. Pemanggil — 10 berkas, 3 nama activity

`[terverifikasi]` Langkah ber-`<pyStepsActivityName>` = `Call CountRateRetroCov`, alamat dari
`<pyStepPageReference>`:

### 3.1 `InsertFacoutProduction` — 3 panggilan, semua bergerbang

Ada di **NB · RNW · Endorsment · DDL**. Alamat (korpus, `RH_1`; `DDL\` memakai `RH_2`):

| Alamat | Gerbang |
| --- | --- |
| `RH_1.pySteps(8).pySteps(1)` | `.CoverageList(1).PremiumRetro==""` |
| `RH_1.pySteps(9).pySteps(2).pySteps(2).pySteps(1)` | `.CoverageList(1).PremiumRetro==""` |
| `RH_1.pySteps(10).pySteps(3).pySteps(1)` | `.CoverageList(1).PremiumRetro==""` |

**Kelas di indeks `Embed-Reference-Rule`:** `Data-Aneka` · **`Data-Cargo`** · `Data-PropertyItem`.

### 3.2 `InsertFacoutProductionEDM` — 3 panggilan, semua bergerbang

Ada di **NB · RNW · Endorsment**.

| Alamat | Gerbang |
| --- | --- |
| `RH_1.pySteps(10).pySteps(1)` | `.CoverageList(1).PremiumRetro==""` |
| `RH_1.pySteps(11).pySteps(2).pySteps(2).pySteps(1)` | `.CoverageList(1).PremiumRetro==""` |
| `RH_1.pySteps(12).pySteps(3).pySteps(1)` | `.CoverageList(1).PremiumRetro==""` |

**Kelas di indeks:** `Data-Aneka` · **`Data-Cargo`** · `Data-PropertyItem`.

### 3.3 `CopyAllObjFacOutFireAneka_ACT` — 2 panggilan, ⛔ **TANPA gerbang**

Ada di **NB · RNW · Endorsment**.

| Alamat | Gerbang |
| --- | --- |
| `RH_1.pySteps(4).pySteps(2).pySteps(1).pySteps(9)` | **(tidak ada)** |
| `RH_1.pySteps(5).pySteps(2).pySteps(1).pySteps(1).pySteps(5)` | **(tidak ada)** |

⚠️ **Nol entri `Embed-Reference-Rule`** untuk `CountRateRetroCov` di berkas ini — indeksnya tidak
merekam resolusi apa pun, meski memanggilnya dua kali.

📌 Namanya memuat **"FireAneka"**, dan kedua panggilannya berada di **dua cabang bersarang berbeda** —
`[dugaan]` satu cabang per lini. **Tidak dapat dipastikan dari korpus** cabang mana memanggil varian
mana.

---

## 4. Kelas yang dirujuk indeks — **tiga**, merata

`[terverifikasi]` Kemunculan `Embed-Reference-Rule` ber-`pyRuleName` = `CountRateRetroCov`:

| Kelas | Berkas |
| --- | ---: |
| `ASM-FW-GISFW-Data-PropertyItem` | **7** |
| `ASM-FW-GISFW-Data-Aneka` | **7** |
| **`ASM-FW-GISFW-Data-Cargo`** | **7** |

**Ketujuh berkas itu sama untuk ketiga kelas** — tiap berkas yang mengindeks satu, mengindeks semua:

`NB\InsertFacoutProduction` · `NB\InsertFacoutProductionEDM` · `RNW\InsertFacoutProduction` ·
`RNW\InsertFacoutProductionEDM` · `Endorsment\InsertFacoutProduction` ·
`Endorsment\InsertFacoutProductionEDM` · `DDL\InsertFacoutProduction`

> ⛔ **Menjawab pertanyaan yang diajukan:** kelas `Data-Cargo` **muncul juga di berkas lain selain
> `InsertFacoutProduction`** — yakni di **`InsertFacoutProductionEDM`**, di ketiga folder. Jadi
> rujukan ke kelas ketiga itu **bukan kejadian tunggal**.

---

## 5. Berapa jalur COB yang memanggil rumus ini

| Jalur | Bukti | Status |
| --- | --- | :-: |
| **FIRE / PropertyItem** | rule ada; dirujuk indeks 7× | ✅ pasti |
| **ANEKA** | rule ada; dirujuk indeks 7× | ✅ pasti |
| **CARGO** | ⛔ rule **tidak ada**; dirujuk indeks 7×; tiga langkah `Call` di tiap pemanggil produksi | ❓ **`[pertanyaan terbuka]`** |

### ⚠️ Jangan disimpulkan salah satu

Dua pembacaan sama-sama konsisten dengan korpus, dan **korpus tidak dapat memutuskan**:

1. **Rule `Data-Cargo` pernah ada lalu dihapus.** Indeks rujukan merekam resolusi **saat berkas
   terakhir disimpan**, jadi jejaknya bertahan setelah rule-nya hilang. Ini pola yang sama dengan
   K-006 dan K-019.
2. **Rule `Data-Cargo` tidak pernah ada.** Indeks dapat mencatat kelas yang **dicoba di-resolve**
   Pega saat menyusun daftar, tanpa menjamin keberadaannya.

**Work owner menyatakan hanya ada DUA rule di Pega.** Keterangan itu **mengalahkan** indeks rujukan
(`PANDUAN-KERJA` §1), tetapi **tidak menghapus** pertanyaan mengapa tiap pemanggil produksi memuat
**tiga** langkah `Call` yang seragam.

📌 **Yang perlu dijawab UI Pega:** apakah langkah `Call` ketiga di
`InsertFacoutProduction`/`InsertFacoutProductionEDM` masih dapat tercapai, dan bila tercapai, ke rule
mana ia resolve.

---

## 6. Perintah audit

```powershell
# 15 berkas menyebut, peka huruf, korpus + DDL
$rx=[regex]'CountRateRetroCov'
$subs=@('Activity','ConnectREST','DataPage','DataTransform','DecisionTable','DecisionTree','Flow',
        'FlowAction','Harness','RDBList','ReportDefinition','Section','SystemSettings','When')
$src=New-Object Collections.Generic.List[string]
foreach($f in @('NB FacIn','RNW Fac In','Endorsment Fac In')){ foreach($s in $subs){
  $p="D:\migrasi\RNM\$f\$s"; if(-not [IO.Directory]::Exists($p)){continue}
  foreach($x in [IO.Directory]::EnumerateFiles($p,'*.xml')){ $src.Add($x) } } }
foreach($x in [IO.Directory]::EnumerateFiles('D:\migrasi\RNM\DDL','*.xml')){ $src.Add($x) }
$n=0; foreach($x in $src){ if($rx.IsMatch([IO.File]::ReadAllText($x))){ $n++ } }
"dipindai=$($src.Count)  menyebut=$n"      # -> 6105 / 15
```

```powershell
# kelas yang dirujuk indeks Embed-Reference-Rule -> PropertyItem 7, Aneka 7, Cargo 7
$k=@{}
foreach($x in $src){
  $t=[IO.File]::ReadAllText($x); if(-not $rx.IsMatch($t)){ continue }
  $d=New-Object Xml.XmlDocument; try{ $d.LoadXml($t) }catch{ continue }
  foreach($nd in $d.SelectNodes('//rowdata')){
    $rn='';$cls='';$obj=''
    foreach($c in $nd.ChildNodes){
      if($c.Name -eq 'pyRuleName'){ $rn=$c.InnerText.Trim() }
      if($c.Name -eq 'pxRuleClassName'){ $cls=$c.InnerText.Trim() }
      if($c.Name -eq 'pxObjClass'){ $obj=$c.InnerText.Trim() } }
    if($obj -eq 'Embed-Reference-Rule' -and $rn -ceq 'CountRateRetroCov' -and $cls -ne ''){
      if($k.ContainsKey($cls)){$k[$cls]++}else{$k[$cls]=1} } } }
$k.GetEnumerator() | Sort-Object Name
```

---

## 7. Yang TIDAK dapat ditentukan

1. `[pertanyaan terbuka]` Apakah varian berkelas `Data-Cargo` pernah ada. Korpus memuat **jejak
   rujukan** tanpa rule-nya; keterangan work owner menyatakan hanya dua. **Butuh UI Pega.**
2. `[pertanyaan terbuka]` **Apa yang memilih varian saat runtime.** Tidak ada gerbang pemilih di
   korpus; pemilihan terjadi lewat **kelas halaman** tempat rule dipanggil. Pemetaan COB → kelas
   **tidak terbaca** dari korpus.
3. `[dugaan]` Dua panggilan tanpa gerbang di `CopyAllObjFacOutFireAneka_ACT` diduga satu cabang per
   lini (Fire dan Aneka), sesuai namanya. ⚠️ **Nama bukan bukti** — belum dibedah.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
