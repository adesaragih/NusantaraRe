# 02 — Skema Oracle yang Terbaca dari Korpus Pega

> **Sumber**: `D:\migrasi\RNM\NB FacIn\`, `D:\migrasi\RNM\RNW Fac In\`, `D:\migrasi\RNM\Endorsment Fac In\` (READ-ONLY).
>
> **Batas pengetahuan yang mengikat seluruh dokumen ini**: korpus **tidak memuat satu baris DDL pun**.
> Tidak ada `CREATE TABLE`, tidak ada dump `ALL_TAB_COLUMNS`, tidak ada indeks, tidak ada constraint.
> Yang ada hanyalah **teks SQL di dalam rule `Rule-Connect-SQL`**. Akibatnya:
>
> - Nama tabel dan nama kolom **terverifikasi** bila dikutip dari SQL.
> - **Tipe kolom, panjang, nullability, primary key, foreign key, dan indeks seluruhnya `belum
>   terverifikasi`** — kecuali bila SQL sendiri membuktikannya (deklarasi variabel PL/SQL,
>   `TO_DATE(...)`, `TO_NUMBER(...)`, `DBMS_LOB.CREATETEMPORARY`).
> - Daftar kolom per tabel di bawah ini adalah **kolom yang dipakai aplikasi**, bukan kolom yang ada.
>   Tabel nyata hampir pasti memiliki kolom lain yang tidak pernah disentuh Pega.
>
> Dokumen ini **melengkapi** `01-model-domain.md` (model domain sisi Pega) dan tidak mengulanginya.
> Yang di sini adalah sisi Oracle: tabel, kolom, verba SQL, prosedur, dan konsekuensi migrasinya.

---

## 1. Angka dasar dan perintah auditnya

### 1.1 Berkas SQL

| Besaran | Nilai |
| --- | --- |
| Berkas `RDBList\*.xml` — NB FacIn | 217 |
| Berkas `RDBList\*.xml` — RNW Fac In | 200 |
| Berkas `RDBList\*.xml` — Endorsment Fac In | 188 |
| **Total berkas SQL** | **605** |

```powershell
foreach ($c in @('NB FacIn','RNW Fac In','Endorsment Fac In')) {
  $n = (Get-ChildItem "D:\migrasi\RNM\$c\RDBList" -File -Filter *.xml | Measure-Object).Count
  "{0,-20} {1}" -f $c, $n
}
```

`[terverifikasi]`

### 1.2 Blok SQL di dalam berkas

Setiap `Rule-Connect-SQL` punya empat kantong SQL: `pyBrowseSQL`, `pyOpenSQL`, `pySaveSQL`,
`pyDeleteSQL`. Yang terisi:

| Kantong | Blok terisi | Catatan |
| --- | --- | --- |
| `pyBrowseSQL` | **605** | 1:1 dengan berkas — setiap berkas punya tepat satu |
| `pyOpenSQL` | **0** | tidak dipakai sama sekali |
| `pySaveSQL` | **8** | seluruhnya boilerplate mati — lihat §7.2 |
| `pyDeleteSQL` | **6** | hanya 2 nama rule — lihat §6.3 |
| **Total** | **619** | |

```powershell
$tags = @('pyBrowseSQL','pyOpenSQL','pySaveSQL','pyDeleteSQL')
$rows = @()
foreach ($cyc in @('NB FacIn','RNW Fac In','Endorsment Fac In')) {
  Get-ChildItem "D:\migrasi\RNM\$cyc\RDBList" -File -Filter *.xml | ForEach-Object {
    $t = Get-Content $_.FullName -Raw -Encoding UTF8
    foreach ($tg in $tags) {
      $m = [regex]::Match($t, "<$tg>(.*?)</$tg>", 'Singleline')
      if ($m.Success -and $m.Groups[1].Value.Trim().Length -gt 0) {
        $rows += [pscustomobject]@{Cycle=$cyc; File=$_.Name; Tag=$tg; SQL=$m.Groups[1].Value}
      }
    }
  }
}
"TOTAL: $($rows.Count)"
$rows | Group-Object Tag | ForEach-Object { "{0,-14} {1}" -f $_.Name, $_.Count }
$rows | Export-Clixml "$env:TEMP\sqlblocks.xml"   # dipakai ulang oleh perintah audit berikutnya
```

`[terverifikasi]`

### 1.3 Klasifikasi verba

Dihitung atas 605 blok `pyBrowseSQL`. Satu blok dapat masuk lebih dari satu kelas (blok PL/SQL yang
berisi `INSERT` sekaligus memanggil prosedur).

| Kelas | Jumlah blok |
| --- | --- |
| `SELECT` murni (tidak menulis apa pun) | **445** |
| Blok PL/SQL anonim (`BEGIN` / `DECLARE`) | **120** |
| — di antaranya murni panggil prosedur, tanpa DML inline | **68** |
| Memuat `INSERT INTO` | **58** |
| Memuat `UPDATE … SET` | **19** |
| Memuat `DELETE FROM` di `pyBrowseSQL` | **0** |
| Memuat `MERGE INTO` | **0** |
| Memuat `COMMIT` | **117** |
| `DELETE FROM` di `pyDeleteSQL` | **6** |

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$b = $rows | Where-Object { $_.Tag -eq 'pyBrowseSQL' }
"SELECT murni      : " + ($b | Where-Object { $_.SQL -notmatch '(?is)\b(INSERT\s+INTO|UPDATE\s+"?[A-Za-z_][A-Za-z0-9_$.]*\s+(a\s+)?SET|DELETE\s+FROM)\b' -and $_.SQL -notmatch '(?i)\b(POOLDATA|FIRE|GENERAL|MBU|NEW_GENERAL)\.[A-Z_][A-Z0-9_]*\s*\(' }).Count
"blok PL/SQL       : " + ($b | Where-Object { $_.SQL -match '(?is)^\s*(DECLARE|BEGIN)\b' }).Count
"  murni prosedur  : " + ($b | Where-Object { $_.SQL -match '(?is)^\s*(DECLARE|BEGIN)\b' -and $_.SQL -notmatch '(?is)\b(INSERT\s+INTO|UPDATE\s+\w+\s+SET)\b' }).Count
"INSERT INTO       : " + ($b | Where-Object { $_.SQL -match '(?is)\bINSERT\s+INTO\b' }).Count
"UPDATE … SET      : " + ($b | Where-Object { $_.SQL -match '(?is)\bUPDATE\s+"?[A-Za-z_][A-Za-z0-9_$.]*\s+(a\s+)?SET\b' }).Count
"DELETE (browse)   : " + ($b | Where-Object { $_.SQL -match '(?is)\bDELETE\s+FROM\b' }).Count
"MERGE             : " + ($b | Where-Object { $_.SQL -match '(?is)\bMERGE\s+INTO\b' }).Count
"COMMIT            : " + ($b | Where-Object { $_.SQL -match '(?is)\bCOMMIT\b' }).Count
"DELETE (pyDelete) : " + ($rows | Where-Object { $_.Tag -eq 'pyDeleteSQL' }).Count
```

`[terverifikasi]`

**Konsekuensi arsitektur.** 445 dari 605 berkas (73,6 %) adalah pembacaan murni. Hanya **52 nama rule
unik** yang benar-benar menulis ke database (59 nama rule memuat pemanggilan berskema, tetapi 7 di
antaranya hanya memanggil **fungsi baca** di dalam `SELECT` — lihat §8.2). Permukaan tulis aplikasi
karena itu **sempit dan dapat diinventarisasi seluruhnya**; ini menguntungkan migrasi.

### 1.4 Tabel dan view unik

| Besaran | Nilai |
| --- | --- |
| **Tabel / view unik dirujuk** | **112** |
| Nama rule `RDBList` unik (lintas 3 siklus) | 239 |
| — ada di ketiga siklus | 163 |
| — ada di 2 siklus | 40 |
| — ada di 1 siklus | 36 |

Verifikasi jumlah: 163×3 + 40×2 + 36×1 = 605. ✓

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$cte = @('PREP_DATA','PREP_PARAM','WORDLIST')   # nama CTE, bukan tabel — lihat §1.5
$tabs = @{}
foreach ($r in $rows) {
  foreach ($m in [regex]::Matches($r.SQL,'(?is)\b(?:FROM|JOIN|INSERT\s+INTO|UPDATE)\s+"?([A-Za-z_][A-Za-z0-9_$]*)"?(?:\s*\.\s*"?([A-Za-z_][A-Za-z0-9_$]*)"?)?')) {
    if ($m.Groups[2].Success) { $t = $m.Groups[2].Value.ToUpper() } else { $t = $m.Groups[1].Value.ToUpper() }
    if ($t -in @('SELECT','DUAL','TABLE')) { continue }
    if ($t -in $cte) { continue }
    $tabs[$t] = 1
  }
}
"TABEL/VIEW UNIK: $($tabs.Count)"
```

`[terverifikasi]`

### 1.5 Tiga nama yang **bukan** tabel

Regex naif atas `FROM`/`JOIN` menangkap tiga nama CTE (`WITH … AS`). Nama-nama ini **tidak ada di
Oracle** dan harus dikeluarkan dari inventaris:

| Nama | Sebenarnya | Sumber |
| --- | --- | --- |
| `PREP_DATA`, `PREP_PARAM` | CTE di dalam `WITH` | `NB FacIn/RDBList/SearchAccumulationbypersetase_SQL.xml` |
| `WORDLIST` | CTE di dalam `WITH` | `NB FacIn/RDBList/CekDuplicateLocationOffer_SQL.xml` |

`[terverifikasi]` — dikutip: `WITH PREP_DATA AS ( SELECT ID, NOTE, ZIPCODE, … FROM ACCUMULATION )`.

### 1.6 Sebaran prefix skema

| Skema | Rujukan berprefix |
| --- | --- |
| `POOLDATA` | 233 |
| `NEW_UNDERWRITING` | 15 |
| `GENERAL` | 14 |
| `FIRE` | 12 |
| `DATAPEGA` | 7 |
| `ARASAPAS` | 1 |

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$sc = @{}
foreach ($r in $rows) {
  foreach ($m in [regex]::Matches($r.SQL,'(?i)\b(?:FROM|JOIN|INSERT\s+INTO|UPDATE)\s+"?([A-Z_][A-Z0-9_$]*)"?\s*\.\s*"?[A-Z_]')) {
    $n = $m.Groups[1].Value.ToUpper(); if (-not $sc.ContainsKey($n)) { $sc[$n]=0 }; $sc[$n]++
  }
}
$sc.GetEnumerator() | Sort-Object {-$_.Value} | ForEach-Object { "{0,-20} {1}" -f $_.Key,$_.Value }
```

`[terverifikasi]`

> **Jebakan.** Banyak tabel dirujuk **kadang dengan prefix `POOLDATA.`, kadang tanpa** — contoh
> `FACINPRODUCTION` (`InsertTreatyProduction_Sql` menulis ke `facinproduction` tanpa prefix,
> `BrowseDataViewPolis1` membaca `POOLDATA.FACINPRODUCTION`). Artinya *default schema* koneksi Pega
> adalah `POOLDATA` `[dugaan kuat]`. **Jangan menyimpulkan ada dua tabel berbeda.** Untuk migrasi:
> kualifikasi skema secara eksplisit di seluruh query Go, jangan mengandalkan default.

### 1.7 Tidak ada SQL inline di luar `RDBList`

Dari 88 rule `Activity` yang punya `pyStepsJavaSource` terisi, **0** memuat kata kunci SQL.

```powershell
$n=0; $hits=0
foreach ($c in @('NB FacIn','RNW Fac In','Endorsment Fac In')) {
  Get-ChildItem "D:\migrasi\RNM\$c\Activity" -File -Filter *.xml | ForEach-Object {
    $t = Get-Content $_.FullName -Raw -Encoding UTF8
    if ($t -match '(?i)<pyStepsJavaSource>\s*[^<\s]') {
      $n++
      $j = [regex]::Match($t,'(?is)<pyStepsJavaSource>(.*?)</pyStepsJavaSource>').Groups[1].Value
      if ($j -match '(?i)\b(SELECT|INSERT\s+INTO|UPDATE\s+\w+\s+SET|DELETE\s+FROM)\b') { $hits++ }
    }
  }
}
"Activity dgn Java source terisi: $n ; memuat kata kunci SQL: $hits"
```

`[terverifikasi]` — **seluruh akses SQL aplikasi terkonsentrasi di 605 berkas `RDBList`.** Ini temuan
yang menyenangkan: inventaris §1 adalah inventaris yang lengkap, bukan sampel.

---

## 2. Tabel produksi — inti aplikasi

### 2.1 `FACINPRODUCTION` (skema `POOLDATA`) — produksi Fac Inward

Tabel hasil akhir spreading. Ditulis baris-per-baris oleh `InsertTreatyProduction_Sql`.

**Asal**: `NB FacIn/RDBList/InsertTreatyProduction_Sql.xml` (rule `InsertTreatyProduction_Sql`).

Jumlah kolom pada `INSERT`: **79** (NB), **79** (RNW), **82** (EDM). Lihat §2.2 — perbedaannya
material.

Kolom terbaca (dikutip apa adanya, termasuk inkonsistensi kapitalisasi di korpus):

```
idpega, nopolis, noendors, businesscode, cedingco, sob, begindate, enddate, province, kabupaten,
kecamatan, kelurahan, roadname, buildingno, zip_code, objectno, object_name, risk_location,
object_item, coverage, coveragenote, rate, percent_rnm, percent_share_spreaded, jn_reas, curr_id,
tsi_menjadi, premi_menjadi, tsi_selisih, premi_selisih, spread_date, accumulationcode, occupationcode,
marketingofficercode, insuredname, trading_name, conveyance_name, percent_ri_comm, rislipceding,
licenseplate, type_name, model_name, group_type, noserial, engineno, qqname, prorate,
pct_ri_comm_selisih, ricomm, ricomm_selisih, object_name_id, layerno, IDSHIP, FROMRUTE, TORUTE,
SAILDATE, POLICYTYPE, LOL_MENJADI, LOL_SELISIH, PCT_BROKERGARE_FEE, BROKERAGE_FEE_MENJADI,
BROKERAGE_FEE_SELISIH, PCT_BROKERAGE_FEE_SELISIH, CEDINGID, SOBID, INSUREDID, IDX_OBJECT_ITEM,
TSI100_MENJADI, TSI100_SELISIH, BINDER, ISB2B, ASKRED, STATUS_BUSINESS, GROUPNAME, TOP_RISK,
NO_MASTER_POLIS, TYPE_EDM, DEDUCTION2_MENJADI, DEDUCTION2_SELISIH
```

`[terverifikasi]`

Tipe kolom yang **dibuktikan SQL** (sisanya `belum terverifikasi`):

| Kolom | Bukti | Kesimpulan |
| --- | --- | --- |
| `begindate`, `enddate`, `spread_date`, `SAILDATE` | `To_date({…}, 'DD/MM/YYYY HH24:MI:SS')` | bertipe tanggal/waktu |
| 21 kolom numerik (`rate`, `percent_rnm`, `tsi_menjadi`, `premi_menjadi`, `tsi_selisih`, `premi_selisih`, …) | dibungkus `To_number(Replace({…},',','.'))` | bertipe numerik; **nilai masuk sebagai string berkoma desimal** |
| sisanya | tidak ada konversi | `belum terverifikasi` |

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$p = ($rows | Where-Object { $_.File -eq 'InsertTreatyProduction_Sql.xml' -and $_.Cycle -eq 'NB FacIn' -and $_.Tag -eq 'pyBrowseSQL' }).SQL
"To_number(Replace(',','.')) : " + ([regex]::Matches($p,"(?i)To_number\(Replace\([^)]*',','\.'").Count)
"To_date                     : " + ([regex]::Matches($p,'(?i)To_date').Count)
```

`[terverifikasi]`

> **PCT_BROKERGARE_FEE — salah ketik yang sudah jadi skema.** Kolom "persentase brokerage fee" dieja
> `PCT_BROKER**GARE**_FEE`, sementara pasangan selisihnya dieja benar `PCT_BROKERAGE_FEE_SELISIH`.
> Salah ketik ini ada di ketiga siklus dan identik. Ini **bukan** kandidat perbaikan dalam migrasi —
> skema Oracle tidak berubah; kode Go harus memakai ejaan yang salah itu, dengan komentar.
> `[terverifikasi]`

> **`nopolis LIKE 'RNM-F%'` sebagai penanda lini bisnis.** `GetTotalAccumulation_sql` menyaring
> `where nopolis like 'RNM-F%'` — nomor polis Fac Inward berawalan literal. Konvensi ini **tidak
> tercatat di kolom mana pun**; ia hanya ada di dalam string SQL.
> Asal: `NB FacIn/RDBList/GetTotalAccumulation_sql.xml`. `[terverifikasi]`

### 2.2 `FACINPRODUCTION` menulis kolom **berbeda** per siklus — MEMBLOKIR

`InsertTreatyProduction_Sql` adalah salah satu dari **3 rule** yang SQL-nya berbeda antar siklus
(dari 163 yang hadir di ketiganya; 160 sisanya identik persis).

| Siklus | 2 kolom terakhir setelah `TYPE_EDM` | Jumlah kolom |
| --- | --- | --- |
| NB FacIn | `DEDUCTION2_MENJADI, DEDUCTION2_SELISIH` | 79 |
| RNW Fac In | `DEDUCTION2_MENJADI, DEDUCTION2_SELISIH` | 79 |
| **Endorsment Fac In** | `NO_KONTRAK, NO_NPWP, NO_KTP, STARTDATE_DEBITUR, ENDDATE_DEBITUR` | **82** |

`[terverifikasi]` — dikutip langsung dari daftar kolom masing-masing berkas.

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$b = $rows | Where-Object { $_.Tag -eq 'pyBrowseSQL' }
$b | Group-Object File | Where-Object { $_.Count -eq 3 -and (($_.Group | ForEach-Object { $_.SQL.Trim() } | Select-Object -Unique).Count -gt 1) } | ForEach-Object { $_.Name }
# -> GenerateImageID_SQL.xml, GetMasterKlausulAge.xml, InsertTreatyProduction_Sql.xml
foreach ($c in @('NB FacIn','RNW Fac In','Endorsment Fac In')) {
  $r = $b | Where-Object { $_.File -eq 'InsertTreatyProduction_Sql.xml' -and $_.Cycle -eq $c }
  $m = [regex]::Match($r.SQL,'(?is)INSERT\s+INTO\s+facinproduction\s*\((.*?)\)\s*VALUES')
  "{0,-20} {1} kolom" -f $c, ($m.Groups[1].Value -split ',').Count
}
```

Konsekuensi:

1. **`FACINPRODUCTION` memuat kolom PII**: `NO_NPWP` (nomor pokok wajib pajak), `NO_KTP` (nomor induk
   kependudukan). Keduanya **hanya ditulis oleh siklus endorsement**. Penanganan data pribadi di
   aplikasi baru harus mencakup dua kolom ini. `[terverifikasi]` bahwa kolomnya ada dan ditulis;
   `[pertanyaan terbuka]` apakah ada kewajiban retensi/masking.
2. **`DEDUCTION2_MENJADI` / `DEDUCTION2_SELISIH` tidak pernah ditulis oleh siklus endorsement.** Baris
   endorsement karena itu selalu `NULL` pada pasangan kolom ini — atau memuat sisa nilai lama bila
   tabel diperbarui alih-alih disisipkan. **Ini harus dikonfirmasi bisnis sebelum implementasi.**
3. Satu tabel, tiga bentuk penulisan. Repository Go tidak boleh punya satu fungsi
   `InsertFacinProduction` generik; perlu tiga varian yang jelas, atau satu varian dengan daftar
   kolom eksplisit per siklus.

### 2.3 `FACINPRODUCTION` — varian penyesuaian kurs

**Asal**: `NB FacIn/RDBList/ForInputCurrencyAdj_Sql.xml` (rule `ForInputCurrencyAdj_Sql`).

`INSERT` ketiga ke tabel yang sama, **67 kolom** (bukan 79). Yang membedakannya: nilai diambil dari
**dua halaman berbeda** — `Datain`/`Datain1` (nilai baru) dan `DatainOld`/`DatainOld1` (nilai lama) —
bercampur kolom per kolom. Contoh dikutip:

```
percent_rnm         <- To_number(Replace({Datain.CARI20},',','.'))       -- BARU
percent_share_spreaded <- To_number(Replace({DatainOld.CARI30},',','.')) -- LAMA
rate                <- To_number(Replace({DatainOld.CARI34},',','.'))    -- LAMA
```

`[terverifikasi]`. **Aturan pemilihan baru-vs-lama per kolom tidak dijelaskan di mana pun** — ia hanya
tersirat dari nama halaman. `[pertanyaan terbuka]`, lihat §12.

### 2.4 `FACOUTPRODUCTION` — produksi Fac Outward / retro

**Asal**: `NB FacIn/RDBList/InsertTreatyProd_Sql.xml` (rule `InsertTreatyProd_Sql`). **65 kolom**,
identik di ketiga siklus.

> Perhatikan jebakan penamaan: rule bernama `InsertTreaty**Prod**_Sql` menulis ke **FACOUTPRODUCTION**,
> sedangkan `InsertTreatyProd**uction**_Sql` menulis ke **FACINPRODUCTION**. Dua nama yang nyaris sama
> untuk dua tabel berbeda. `[terverifikasi]`

```
IDPEGA, POLICYNO, GROUPPANEL, REINSURER_ID, REINSURER_NAME, TGL_PRINT, RISTARTPERIOD, RIENDPERIOD,
START_DATE, END_DATE, RISLIPNO, NOENDORS, PACKINGID, PACKINGNOTE, GOODNOTE, TRADINGNOTE, SHIPID,
FROMRUTE, TORUTE, SAILDATE, CONVEYANCENOTE, OBJECTNO, ZIPCODE, PROVINCE, CITY, DISTRICT, RW, ADDRESS,
BUILDINGNO, ROADNAME, OBJECTNAME, OCCUPATION, OBJECTITEM, OBJECTITEMID, BRANDNAME, LICENSEPLATE,
TYPENAME, MODELNAME, ENGINENUMBER, CLASSPA, DOB, CURRENCY, CURRENCYID, TSIRNM, TSISPREADED,
PCTOFFERED, SHAREOFFERED, OBJECTPREMI, RICOMM, COMMISION, RATE, SHAREOFFERED_SELISIH,
OBJECTPREMI_SELISIH, COMMISION_SELISIH, TYPEFACULTATIVE, PRORATE, RATE_COVERAGE,
PREMI_COVERAGE_MENJADI, PREMI_COVERAGE_SELISIH, COVERAGE_NAME, COVERAGE_ID, COMMISION_COVERAGE_PCT,
COMMISION_COVERAGE_MENJADI, COMMISION_COVERAGE_SELISIH, OBJECTNO_FACIN
```

`[terverifikasi]`

Catatan: kolom `DOB` (tanggal lahir) dan `RW` (nama satuan wilayah) ada di tabel ini. `CURRENCY` dan
`CURRENCYID` berdampingan — dua representasi mata uang, `[pertanyaan terbuka]` mana yang otoritatif.

### 2.5 `FACINOFFER` (skema `POOLDATA`) — penawaran, sebelum jadi produksi

**Asal**: `NB FacIn/RDBList/InsertOfferProduction_Sql.xml` (rule `InsertOfferProduction_Sql`).
**50 kolom**, identik di ketiga siklus.

```
idpega, nooffer, businesscode, cedingco, sob, begindate, enddate, province, kabupaten, kecamatan,
kelurahan, roadname, buildingno, zip_code, objectno, object_name, risk_location, object_item,
coverage, coveragenote, rate, percent_rnm, percent_share_spreaded, jn_reas, curr_id, tsi_menjadi,
premi_menjadi, tsi_selisih, premi_selisih, spread_date, accumulationcode, occupationcode,
marketingofficercode, insuredname, trading_name, conveyance_name, percent_ri_comm, rislipceding,
licenseplate, type_name, model_name, group_type, noserial, ENGINENO, IDSHIP, FROMRUTE, TORUTE,
SAILDATE, GROUPNAME, TOP_RISK
```

`[terverifikasi]`

#### Asimetri konversi angka antara `FACINOFFER` dan `FACINPRODUCTION` — MEMBLOKIR

Dua tabel dengan nama kolom yang sama (`rate`, `tsi_menjadi`, `premi_menjadi`, `percent_rnm`, …)
diisi dengan **dua konvensi yang berlawanan**:

| | `FACINOFFER` (`InsertOfferProduction_Sql`) | `FACINPRODUCTION` (`InsertTreatyProduction_Sql`) |
| --- | --- | --- |
| Jumlah `To_number(...)` | **0** | **21** |
| Perlakuan `rate` | `Replace({Datain.CARI34},'.',',')` — titik **menjadi** koma | `To_number(Replace({Datain.CARI34},',','.'))` — koma **menjadi** titik, lalu ke angka |
| Kolom numerik lain | disisipkan mentah, tanpa konversi | seluruhnya `To_number(Replace(…))` |

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$o = ($rows | Where-Object { $_.File -eq 'InsertOfferProduction_Sql.xml'   -and $_.Cycle -eq 'NB FacIn' -and $_.Tag -eq 'pyBrowseSQL' }).SQL
$p = ($rows | Where-Object { $_.File -eq 'InsertTreatyProduction_Sql.xml' -and $_.Cycle -eq 'NB FacIn' -and $_.Tag -eq 'pyBrowseSQL' }).SQL
"FACINOFFER      To_number = " + ([regex]::Matches($o,'(?i)To_number').Count)
"FACINPRODUCTION To_number = " + ([regex]::Matches($p,'(?i)To_number').Count)
```

`[terverifikasi]` bahwa asimetrinya ada.

Pembacaannya yang paling wajar: **`FACINOFFER` menyimpan angka sebagai teks berkoma desimal**,
`FACINPRODUCTION` menyimpan sebagai `NUMBER`. Tetapi ini `[dugaan]` — **tipe kolom
`FACINOFFER.RATE`, `FACINOFFER.TSI_MENJADI` dan seterusnya belum terverifikasi.** Korpus tidak
memuat DDL, dan `INSERT` tanpa konversi bisa juga berarti Oracle melakukan konversi implisit dengan
`NLS_NUMERIC_CHARACTERS` tertentu — yang akan **gagal senyap atau salah** bila sesi punya setelan
berbeda.

**Ini memblokir.** Sebelum menulis satu baris kode repository untuk `FACINOFFER`, jalankan atas
database nyata:

```sql
SELECT table_name, column_name, data_type, data_length, data_precision, data_scale
FROM   all_tab_columns
WHERE  owner = 'POOLDATA'
  AND  table_name IN ('FACINOFFER','FACINPRODUCTION','FACOUTPRODUCTION')
ORDER  BY table_name, column_id;
```

Bila `FACINOFFER.RATE` ternyata `VARCHAR2`, maka seluruh perbandingan angka atasnya adalah
perbandingan string — persis kelas cacat yang sudah dicatat di `CLAUDE.md` §4.1, dan **harus
dipertahankan apa adanya saat migrasi**, bukan diperbaiki diam-diam.

### 2.6 `TREATYPRODUCTION_BACKUP` (skema `POOLDATA`) — cadangan produksi

**Asal**: `NB FacIn/RDBList/InsertFacinProductionBackup_Sql.xml`, dan varian kurs
`ForInputCurrencyAdjBackup_Sql.xml`. **70 kolom**.

Bentuknya paralel dengan `FACINPRODUCTION` (nama kolom sama, termasuk pasangan `_MENJADI`/`_SELISIH`).
`[terverifikasi]` ia ditulis; **kapan dan mengapa dibuat cadangan `belum terverifikasi`** — rule
`When` pemicunya tidak terbaca dari sisi SQL.

### 2.7 `CEDING_FACINPRODUCTION` (skema `POOLDATA`)

Ditulis lewat prosedur `POOLDATA.InsertUpdateCedingProduction`, dihapus lewat `pyDeleteSQL`.

**Asal**: `NB FacIn/RDBList/InsertCedingProduction_SQL.xml`.

```sql
-- pyDeleteSQL
DELETE FROM POOLDATA.CEDING_FACINPRODUCTION WHERE IDPEGA ={pyWorkPage.pzInsKey}
```

```sql
-- pyBrowseSQL
BEGIN
pooldata.InsertUpdateCedingProduction(
  {pyWorkPage.pzInsKey},
  {pyWorkPage.OfferFacIn.PolicyData.PolicyNo},
  {pyWorkPage.OfferFacIn.PolicyData.EndorsementNo},
  {TempInput.CARI3}, {TempInput.CARI4}, {TempInput.CARI5},
  {OutputData.HASIL1 out}, {OutputData.HASIL2 out});
COMMIT;
END;
```

Kolom yang terbaca hanya `IDPEGA`. Enam parameter prosedur → **kolom sisanya `belum terverifikasi`**.
`[terverifikasi]` untuk `IDPEGA`, sisanya tidak.

### 2.8 Tabel produksi jiwa (life)

| Tabel | Rule penulis | Catatan |
| --- | --- | --- |
| `POOLDATA.FACINLIFE` | `InsertFacinLife_Sql`, `InsertFacinLifeMonthly_Sql` | dua varian: sekali-bayar dan bulanan |
| `POOLDATA.FACINOFFERLIFE` | `InsertFacinOfferLife_Sql` | hanya ada di 2 siklus |
| `POOLDATA.FACINSPREADLIFE` | `InsertIntoFacinSPreadLife_Sql`, `InsertIntoFacinSPreadLifeMonthly_Sql` | idem |
| `RATE_LIFE` | dibaca `BrowseLifeRateRetro_SQL` | kolom: `ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE, TYPE` `[terverifikasi]` |
| `POOLDATA.M_ACCUMULATION_LIFE` | prosedur `PEGA_M_ACCUMULATION_LIFE` | menerima data pribadi tertanggung — lihat §9 |

Kolom lengkap `FACINLIFE` / `FACINSPREADLIFE` ada di berkas masing-masing dan belum ditabulasi di
sini karena lini jiwa berada di luar jalur akseptasi utama; **bila lini jiwa masuk lingkup migrasi,
tabulasi ini wajib dilengkapi.** `[pertanyaan terbuka]`

---

## 3. `JSON_POLIS` — versi polis dan riwayat endorsement

### 3.1 Kolom

**Asal**: 62 blok SQL merujuk `JSON_POLIS`; 47 di antaranya hanya menyentuh tabel ini.

| Kolom | Bukti |
| --- | --- |
| `IDPEGA` | `SELECT IDPEGA AS HASIL FROM JSON_POLIS WHERE NOPOLIS=… and PRODKE=…` (`GetIDPega_SQL`) |
| `NOPOLIS` | idem |
| `NOENDORS` | `select noendors as HASIL1 from json_polis where nopolis =…` (`GetNoEndors_SQL`) |
| `PRODKE` | lihat §3.2 |
| `DATA_JSON` | `Update Json_Polis set DATA_JSON = {InputData.CARI2}, …` (`UpdatePolisEndorsement_SQL`) |
| `USERNAME` | idem |
| `TGL_INPUT` | `order by TGL_INPUT desc` (`GetProdKeOldData_SQL`) |
| `TGL_PROD` | `InsertIntoJsonError_SQL` (tabel error yang sebangun) |
| `OLDNOPOLIS` | idem |
| `STS_KONVERSI` | terbaca di predikat |
| `STS_KONVERSI_RETRO` | `update json_polis set STS_KONVERSI_RETRO = '8' where idpega = …` (`UpdateSTSKonversi_SQL`) |

`[terverifikasi]`. **Tipe seluruh kolom `belum terverifikasi`**, dengan satu pengecualian di §3.3.

> `STS_KONVERSI_RETRO` di-set ke literal **string** `'8'`, bukan angka `8`. Arti kode `8` **tidak
> dijelaskan korpus** — `[pertanyaan terbuka]`. Jangan menebak.

### 3.2 `PRODKE` — penomoran versi polis

Idiom "ambil versi terakhir" muncul dalam **tiga bentuk yang tidak setara**:

**Bentuk A — `COUNT(...) - 1`** (dominan; `GetProdKeEDM_SQL`, `GetNoEndors_SQL`, `GetEDMStatus_SQL`,
`GetBusinessType_Sql`, `GetEDMOldData_SQL`):

```sql
-- Asal: Endorsment Fac In/RDBList/GetProdKeEDM_SQL.xml
select PRODKE as HASIL2 from json_polis
where NOPOLIS= {InputData.CARI17}
  and PRODKE=(SELECT COUNT(NOPOLIS)-1 FROM JSON_POLIS b WHERE NOPOLIS={InputData.CARI17})
```

**Bentuk B — `ORDER BY PRODKE DESC` + `ROWNUM<=1`** (`GetLastPPNCheckEDM`, `GetEDMOldIDPEGA`):

```sql
-- Asal: NB FacIn/RDBList/GetLastPPNCheckEDM.xml
select * from (
  select a.data_json.PPnCheck as cari1 from pooldata.json_polis a
  where nopolis ={pyWorkPage.OfferFacIn.QuotationData.OldPolicyNo} order by prodke desc
) where rownum <=1
```

**Bentuk C — `ORDER BY TGL_INPUT DESC`** (`GetProdKeOldData_SQL`):

```sql
-- Asal: NB FacIn/RDBList/GetProdKeOldData_SQL.xml
select PRODKE as HASIL2 from pooldata.json_polis where NOPOLIS= {TempPolis.CARI4} order by TGL_INPUT desc
```

`[terverifikasi]` ketiganya ada.

Apa yang dibuktikan bentuk A: **`PRODKE` berbasis 0 dan rapat** — versi pertama `0`, dan jumlah baris
untuk satu `NOPOLIS` selalu `PRODKE_maks + 1`. Ini mengonfirmasi model "riwayat endorsement =
baris bertambah" di `CLAUDE.md` §4.3. `[terverifikasi]`

Apa yang **tidak** dibuktikan, dan berbahaya:

1. Bentuk A **pecah bila satu baris pernah dihapus** — `COUNT-1` lalu menunjuk `PRODKE` yang salah atau
   tidak ada, dan query mengembalikan **nol baris secara senyap**. Tidak ada `MAX(PRODKE)` di mana pun.
2. Bentuk A dan bentuk C **tidak setara** bila `TGL_INPUT` tidak monoton terhadap `PRODKE`.
3. Tidak ada bukti korpus tentang **siapa yang menetapkan nilai `PRODKE`** saat penyisipan. `INSERT`
   ke `JSON_POLIS` tidak pernah eksplisit — semuanya lewat prosedur `POOLDATA.INSERTJSONPOLIS`
   (§8.1), yang isinya `belum terverifikasi`. **Argumen ke-7 prosedur itu adalah literal `'0'` pada
   jalur NB dan `{InputData.CARI4}` pada jalur EDM** — pola yang konsisten dengan "prodke", tetapi
   itu `[dugaan]`, bukan bukti.

**Untuk migrasi**: pilih **satu** definisi "versi terakhir" — `MAX(PRODKE)` adalah yang benar — tetapi
**catat bahwa ini mengubah perilaku** terhadap data yang pernah dihapus, dan bawa ke bisnis sebelum
mengubahnya. Lihat §12.

### 3.3 `DATA_JSON` — kolom dokumen, diakses dengan notasi titik

Korpus mengakses isi `DATA_JSON` langsung di SQL, memakai alias tabel:

```sql
-- Asal: Endorsment Fac In/RDBList/GetBusinessType_Sql.xml
select a.Data_json.QuotationData.BusinessType as CARI2 from json_polis a where nopolis =…
```

Jalur yang terbaca di korpus, lengkap:

```
data_json.PPnCheck
data_json.PolicyData.StartDateTime      data_json.PolicyData.EndDateTime
data_json.QuotationData.BusinessType    data_json.QuotationData.CedingCoName
data_json.QuotationData.EdmDate         data_json.QuotationData.EdmType
data_json.QuotationData.InsuredName     data_json.QuotationData.NoOfferSlip
data_json.QuotationData.QQName          data_json.QuotationData.SobName
```

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
($rows | ForEach-Object { [regex]::Matches($_.SQL,'(?i)DATA_JSON\s*\.\s*[A-Za-z0-9_.]+') | ForEach-Object { $_.Value } }) | Sort-Object -Unique
```

`[terverifikasi]` bahwa jalur-jalur itu diakses.

**Tipe `DATA_JSON` `belum terverifikasi`.** Notasi titik Oracle bekerja untuk kolom `JSON`, untuk
`VARCHAR2`/`CLOB`/`BLOB` yang punya constraint `IS JSON`, **dan** untuk kolom bertipe objek. Tiga
kemungkinan berbeda dengan perilaku berbeda. Yang dapat dipastikan: prosedur saudaranya
`INSERTJSONPOLISMONITORING` menerima payload yang **dideklarasikan `CLOB`** —

```sql
-- Asal: NB FacIn/RDBList/INSERTJSON_JSONPOLISMONITORING_FACIN.xml
DECLARE
  P_JSONDATA CLOB;
BEGIN
  DBMS_LOB.CREATETEMPORARY(P_JSONDATA, true);
  P_JSONDATA := {InputData.CARI3};
  …
```

`[terverifikasi]` untuk parameter prosedur monitoring; `[dugaan]` bahwa `JSON_POLIS.DATA_JSON` juga
`CLOB`. **Jangan menulis kode yang mengandaikan tipe apa pun sebelum `all_tab_columns` diperiksa.**

Konsekuensi migrasi yang tidak boleh dilewatkan: **skema `DATA_JSON` tidak terdokumentasi di mana
pun.** Sebelas jalur di atas adalah yang *dibaca oleh SQL*; dokumen yang sesungguhnya jauh lebih besar
(ia adalah serialisasi agregat `OfferFacIn` — lihat `01-model-domain.md`). Struktur JSON-nya
**`belum terverifikasi`** dan **tidak boleh ditebak**.

### 3.4 Tabel turunan `JSON_*`

| Tabel | Kolom terbaca | Asal |
| --- | --- | --- |
| `JSON_POLIS_ERROR` | `TGL_INPUT, DATA_JSON, NOPOLIS, NOENDORS, PRODKE, IDPEGA, OLDNOPOLIS, TGL_PROD` | `NB FacIn/RDBList/InsertIntoJsonError_SQL.xml` `[terverifikasi]` |
| `POOLDATA.JSON_POLIS_MONITORING` | `NOPOLIS, ERR_NOTE` (+ parameter prosedur) | `NB FacIn/RDBList/UpdateErrorNoteJsonPolisMonitoring.xml` `[terverifikasi]` |
| `POOLDATA.JSON_OFFER` | `ID, JSONDATA, TGL_BIND`, dan jalur `JSONDATA.IDNewBisnis` | `UpdateJsonOffer_SQL.xml`, `UpdateTgl_BindJsonOffer_SQL.xml` `[terverifikasi]` |
| `POOLDATA.JSON_FOLLOWING` | 21 kolom, lihat di bawah | `NB FacIn/RDBList/InsertJsonFollowing_SQL.xml` `[terverifikasi]` |

`JSON_FOLLOWING`:

```
TGL_INPUT, DATA_JSON, IDPEGA, CEDINGCO, CEDINGCONAME, MARKETINGNAME, MARKETINGCODE, MOID, TEAMGROUP,
BUSINESSNAME, BUSINESSCODE, SPECIALACCEPTANCE, CEDINGCONFIRM, RNMSHARE, TOTALPREMINUSARE,
TOTALTSINUSARE, TOTALTSINUSARESPREADING, TOTALTSITOPRISK, INSUREDNAME, INSUREDID, STARTDATE
```

Perhatikan: `JSON_OFFER` memakai nama kolom `JSONDATA` (tanpa garis bawah), sedangkan `JSON_POLIS`,
`JSON_POLIS_ERROR` dan `JSON_FOLLOWING` memakai `DATA_JSON`. **Dua konvensi di satu skema.**
`[terverifikasi]`

---

## 4. Pasangan kolom nilai-sesudah / selisih — daftar lengkap

Ini pola sentral endorsement: setiap besaran yang dapat berubah disimpan sebagai **nilai sesudah
(`_MENJADI`)** dan **delta (`_SELISIH`)**. `CLAUDE.md` §4.3 mewajibkan pasangan ini dipertahankan
berpasangan.

### 4.1 Pasangan lengkap (8)

| Besaran | Kolom nilai-sesudah | Kolom selisih | Tabel |
| --- | --- | --- | --- |
| TSI | `TSI_MENJADI` | `TSI_SELISIH` | `FACINPRODUCTION`, `FACINOFFER`, `TREATYPRODUCTION_BACKUP` |
| TSI 100 % | `TSI100_MENJADI` | `TSI100_SELISIH` | `FACINPRODUCTION`, `TREATYPRODUCTION_BACKUP` |
| Premi | `PREMI_MENJADI` | `PREMI_SELISIH` | `FACINPRODUCTION`, `FACINOFFER`, `TREATYPRODUCTION_BACKUP` |
| Premi per coverage | `PREMI_COVERAGE_MENJADI` | `PREMI_COVERAGE_SELISIH` | `FACOUTPRODUCTION` |
| Limit of Liability | `LOL_MENJADI` | `LOL_SELISIH` | `FACINPRODUCTION`, `TREATYPRODUCTION_BACKUP` |
| Brokerage fee | `BROKERAGE_FEE_MENJADI` | `BROKERAGE_FEE_SELISIH` | `FACINPRODUCTION`, `TREATYPRODUCTION_BACKUP` |
| Komisi per coverage | `COMMISION_COVERAGE_MENJADI` | `COMMISION_COVERAGE_SELISIH` | `FACOUTPRODUCTION` |
| Deduction 2 | `DEDUCTION2_MENJADI` | `DEDUCTION2_SELISIH` | `FACINPRODUCTION` (**NB/RNW saja** — lihat §2.2) |

### 4.2 Kolom `_SELISIH` **tanpa** pasangan `_MENJADI` (6)

Di sini nilai-sesudah disimpan pada kolom bernama polos:

| Besaran | Kolom nilai-sesudah | Kolom selisih | Tabel |
| --- | --- | --- | --- |
| Share offered | `SHAREOFFERED` | `SHAREOFFERED_SELISIH` | `FACOUTPRODUCTION` |
| Premi objek | `OBJECTPREMI` | `OBJECTPREMI_SELISIH` | `FACOUTPRODUCTION` |
| Komisi | `COMMISION` | `COMMISION_SELISIH` | `FACOUTPRODUCTION` |
| RI commission | `RICOMM` | `RICOMM_SELISIH` | `FACINPRODUCTION` |
| Persentase RI comm | `PERCENT_RI_COMM` | `PCT_RI_COMM_SELISIH` | `FACINPRODUCTION` |
| Persentase brokerage | `PCT_BROKERGARE_FEE` *(sic)* | `PCT_BROKERAGE_FEE_SELISIH` | `FACINPRODUCTION` |

**Dua konvensi penamaan untuk satu konsep.** Enam pasangan terakhir tidak mengikuti pola `_MENJADI`,
dan satu di antaranya punya ejaan berbeda antara pasangannya (`BROKERGARE` vs `BROKERAGE`). Setiap
mapper Go yang menurunkan nama kolom dari nama field secara mekanis **akan salah pada keenamnya**.

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
"kolom *_MENJADI unik:"
($rows | ForEach-Object { [regex]::Matches($_.SQL,'(?i)\b([A-Z0-9_]+_MENJADI)\b') | ForEach-Object { $_.Groups[1].Value.ToUpper() } }) | Sort-Object -Unique
"kolom *SELISIH* unik:"
($rows | ForEach-Object { [regex]::Matches($_.SQL,'(?i)\b([A-Z0-9_]*SELISIH[A-Z0-9_]*)\b') | ForEach-Object { $_.Groups[1].Value.ToUpper() } }) | Sort-Object -Unique
```

Hasil: **8** kolom `_MENJADI`, **14** kolom `*SELISIH*`. 8 berpasangan + 6 tak berpasangan = 14. ✓
`[terverifikasi]`

### 4.3 Yang tidak dijelaskan korpus

- Apakah `_SELISIH` = nilai_baru − nilai_lama, atau sebaliknya? **`[pertanyaan terbuka]`** — tidak ada
  satu pun SQL yang menghitungnya; nilainya selalu datang sudah jadi dari halaman Pega.
- Pada baris **New Business** (bukan endorsement), apakah `_SELISIH` diisi nol atau sama dengan
  `_MENJADI`? **`[pertanyaan terbuka]`**. Ini menentukan apakah agregasi seperti
  `sum(Premi_selisih)` di `GetSummaryRiskAccumulation_Sql` menghitung eksposur total atau hanya
  perubahan.
- Mata uang tidak pernah berdampingan dengan pasangan ini di tabel yang sama selain kolom `curr_id`
  tunggal per baris — konsisten dengan peringatan `CLAUDE.md` §4.1 bahwa korpus kehilangan informasi
  mata uang pada tingkat nilai.

---

## 5. Tabel limit akseptasi per lini bisnis

Inilah tabel yang menggerakkan tangga persetujuan. **7 tabel**, satu per lini bisnis.

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$lim = @{}
foreach ($r in $rows) {
  foreach ($m in [regex]::Matches($r.SQL,'(?i)\b(?:FROM|JOIN)\s+"?(?:POOLDATA\.)?"?(M_LIMIT[A-Z0-9_]*)\b')) {
    $n = $m.Groups[1].Value.ToUpper(); if (-not $lim.ContainsKey($n)) { $lim[$n] = 0 }; $lim[$n]++
  }
}
"TABEL M_LIMIT_* UNIK: $($lim.Count)"
$lim.GetEnumerator() | Sort-Object Name | ForEach-Object { "{0,-42} {1}" -f $_.Key, $_.Value }
```

| Tabel | Rujukan | Lini bisnis | Rule pembaca |
| --- | --- | --- | --- |
| `POOLDATA.M_LIMIT_PROPERTYY` | 21 | Property (utama) | `GetLimitAkseptasi_SQL`, `GetLimitAkseptasiJUWA_SQL`, `GetLimitAkseptasiBanding_SQL`, `GetAksepBanding_SQL` |
| `M_LIMIT_ENGINEERINGG` | 18 | Engineering | `GetLimitAccEngineeringUW_SQL`, `…JUWA_SQL`, `…Banding_SQL` |
| `POOLDATA.M_LIMIT_NONPROPANDENGG` | 18 | Non-property & non-engineering | `GetLimitAkseptasiNonFire_SQL`, `…JUWA_SQL`, `…Banding_SQL` |
| `M_LIMIT_PROPERTY_NON_PREFERREDD` | 18 | Property non-preferred | `GetLimitAkseptasiNonPrefer_SQL`, `…JUWA_SQL`, `…Banding_SQL` |
| `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` | 12 | Property preferred commercial | `GetLimitAkseptasiPreferedComm_SQL`, `…JUWA_SQL` |
| `M_LIMIT_FINANCIALINS` | 9 | Financial lines (bond & kredit) | `GetLimitAkseptasiBond_SQL`, `…KreditCL_SQL`, `…KreditNCL_SQL` |
| `M_LIMIT_LIFE` | 1 | Jiwa | `GetLimitAkseptasiLife_SQL` |

`[terverifikasi]`

> **Tiga nama tabel berakhir huruf ganda** — `PROPERTYY`, `ENGINEERINGG`, `NON_PREFERREDD`,
> `COMMERCIALL`. Ini bukan salah salin: keempatnya konsisten di ketiga siklus. **Pertahankan ejaannya.**
> `[terverifikasi]`

### 5.1 Kolom

Lima tabel "bertangga" (`M_LIMIT_PROPERTYY`, `M_LIMIT_ENGINEERINGG`, `M_LIMIT_NONPROPANDENGG`,
`M_LIMIT_PROPERTY_NON_PREFERREDD`, `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL`) berbagi bentuk yang sama:

| Kolom | Peran terbaca | Bukti |
| --- | --- | --- |
| `TEAM_GROUP` | penyekat tangga per tim | `WHERE team_group = {…QuotationData.TeamGroup}` |
| `LOGIN` | identitas operator (kunci "siapa saya sekarang") | `AND LOGIN ={OperatorID.pyUserIdentifier}` |
| `JABATAN` | jabatan / anak tangga | `SELECT JABATAN AS CARI1` |
| `JABATAN_ATASAN` | jabatan atasan (jalur banding) | `GetAksepBanding_SQL` |
| `NAMA` | **nama orang** — lihat peringatan di bawah | `NAMA AS CARI5` |
| `LIMIT_BOTTOM` | ambang bawah tangga jalur normal | `LIMIT_BOTTOM > (SELECT LIMIT_BOTTOM FROM … WHERE LOGIN = …)` |
| `LIMIT_BOTTOM2` | ambang bawah tangga **jalur banding** | `GetLimitAkseptasiBanding_SQL` |
| `MAX_LIMIT_IDR` | batas nilai IDR | `MAX_LIMIT_IDR AS CARI2` |
| `MAX_LIMIT_USD` | batas nilai USD | `MAX_LIMIT_USD AS CARI3` |
| `BATAS_WAKTU` | batas waktu (SLA) | `BATAS_WAKTU AS CARI4` |
| `WORKBASKET` | keranjang kerja Pega | `WHERE WORKBASKET IN (…)` — §6.2 |

`M_LIMIT_FINANCIALINS` berbeda: `JABATAN`, `NAMA`, dan **tiga** ambang terpisah
`LIMITBOND_BOTTOM`, `LIMITCREDITCL_BOTTOM`, `LIMITCREDITNCL_BOTTOM`.

`M_LIMIT_LIFE`: `ID`, `JABATAN`, `NAMA`, `LIMIT_BOTTOM`.

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$lb = $rows | Where-Object { $_.SQL -match '(?i)M_LIMIT' }
$ids=@{}
foreach ($r in $lb) {
  $s = $r.SQL -replace '(?s)\{[^}]*\}','?'
  foreach ($m in [regex]::Matches($s,'\b([A-Za-z_][A-Za-z0-9_]{2,})\b')) {
    $v=$m.Groups[1].Value.ToUpper()
    if ($v -match '^(SELECT|FROM|WHERE|AND|OR|ORDER|GROUP|BY|ASC|DESC|AS|IN|NOT|IS|NULL|POOLDATA|CARI[0-9]*|CARID[0-9]*|ROWNUM|M_LIMIT[A-Z_]*|HISTORYAKSEPTASIPEGA)$') { continue }
    if (-not $ids.ContainsKey($v)) { $ids[$v]=0 }; $ids[$v]++
  }
}
$ids.GetEnumerator() | Sort-Object {-$_.Value} | ForEach-Object { "{0,-30} {1}" -f $_.Key,$_.Value }
```

`[terverifikasi]`

> ### Peringatan data pribadi — `NAMA` dan `LOGIN`
>
> Kolom `NAMA` pada ketujuh tabel limit memuat **nama orang**; `LOGIN` memuat **identitas operator**.
> Query `GetLimitAkseptasi_SQL` memilih penyetuju berikutnya dengan membandingkan `LIMIT_BOTTOM`
> pemakai yang sedang login terhadap tangga — yaitu **otorisasi berbasis baris data, bukan berbasis
> peran**. Sesuai `CLAUDE.md` §3.5, **tidak satu pun nilai dari kolom ini boleh masuk ke kode, tes,
> fixture, komentar, atau tiket.** Yang boleh dicatat hanya mekanismenya, seperti di atas.
>
> Untuk migrasi: pertahankan tabelnya apa adanya (skema Oracle tidak berubah), tetapi **tarik nilai
> `NAMA`/`LOGIN` hanya pada saat runtime**, jangan pernah disalin ke artefak repo.

### 5.2 Bentuk query tangga

```sql
-- Asal: NB FacIn/RDBList/GetLimitAkseptasi_SQL.xml (rule GetLimitAkseptasi_SQL)
SELECT JABATAN AS CARI1, MAX_LIMIT_IDR AS CARI2, MAX_LIMIT_USD AS CARI3, BATAS_WAKTU AS CARI4, NAMA AS CARI5
  FROM POOLDATA.M_LIMIT_PROPERTYY
 WHERE team_group = {pyWorkPage.OfferFacIn.QuotationData.TeamGroup}
   AND LIMIT_BOTTOM > (SELECT LIMIT_BOTTOM
                         FROM POOLDATA.M_LIMIT_PROPERTYY
                        WHERE team_group = {…TeamGroup} AND LOGIN ={OperatorID.pyUserIdentifier})
   AND LIMIT_BOTTOM < {DataSearch.CARID2}
 ORDER BY LIMIT_BOTTOM ASC
```

Tiga varian per lini, dengan perbedaan yang **halus dan material**:

| Varian | Perbedaan | Berkas |
| --- | --- | --- |
| UW (normal) | titik awal dari `LOGIN = operator yang login` | `GetLimitAkseptasi_SQL` |
| JUWA | titik awal dari `JABATAN = {DataSearch.CARI5}` — **jabatan, bukan operator** | `GetLimitAkseptasiJUWA_SQL` |
| Banding | memakai `LIMIT_BOTTOM2`, **bukan** `LIMIT_BOTTOM`, pada kedua sisi perbandingan | `GetLimitAkseptasiBanding_SQL` |

`[terverifikasi]`

> **Cacat yang terbawa.** Varian Banding menyaring dengan `LIMIT_BOTTOM2` tetapi mengurutkan dengan
> `ORDER BY LIMIT_BOTTOM ASC` — kolom yang berbeda. Bila kedua kolom tidak searah, **urutan anak
> tangga banding salah.** Ini **kandidat perbaikan, bukan bagian migrasi** (`CLAUDE.md` §1):
> pertahankan `ORDER BY LIMIT_BOTTOM` dan catat sebagai temuan untuk bisnis.
> `[terverifikasi]` bahwa kolomnya berbeda; `[pertanyaan terbuka]` apakah itu disengaja.

---

## 6. `HISTORYAKSEPTASIPEGA` — riwayat akseptasi yang **dibaca untuk mengambil keputusan**

Pertanyaan yang diminta dijawab: apakah tabel ini hanya di-`INSERT` (audit), atau juga di-`SELECT`
untuk mengambil keputusan? **Jawabannya: juga di-`SELECT`, dan hasilnya menggerakkan alur.**
`[terverifikasi]` — buktinya dikutip di bawah.

Tabel ini dirujuk oleh **9 blok SQL / 3 nama rule**, masing-masing hadir di ketiga siklus.

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$h = $rows | Where-Object { $_.SQL -match '(?i)HISTORYAKSEPTASIPEGA' }
"blok: $($h.Count) ; nama rule unik: " + (($h | ForEach-Object { $_.File } | Sort-Object -Unique).Count)
$h | ForEach-Object { $_.File } | Sort-Object -Unique
```

### 6.1 Kolom

```sql
-- Asal: NB FacIn/RDBList/InsertHistoryAkseptasiPega_Sql.xml (rule InsertHistoryAkseptasiPega_Sql)
BEGIN
INSERT INTO HISTORYAKSEPTASIPEGA
(ID_PEGA, Tgl_Transfer, Status, Username, Workbasket, ID_KOMITE)
VALUES
( {InsertHistory.CARI1}, sysdate, {InsertHistory.CARI5}, {InsertHistory.CARI4},
  {InsertHistory.CARI2}, {InsertHistory.CARI6} );
COMMIT;
END;
```

| Kolom | Catatan |
| --- | --- |
| `ID_PEGA` | kunci case — diisi `pyWorkPage.pzInsKey` pada sisi baca |
| `TGL_TRANSFER` | `sysdate` saat penyisipan → bertipe tanggal/waktu `[terverifikasi]` |
| `STATUS` | enumerasi; satu nilai terbaca: `'REJECT'` (§6.3) |
| `USERNAME` | identitas operator |
| `WORKBASKET` | keranjang kerja; nilai terbaca: `'ReasFacInMarketing'`, `'ReasFacInTeamLeader'` |
| `ID_KOMITE` | arti `belum terverifikasi` |

`[terverifikasi]` untuk keberadaan kolom; tipe seluruhnya `belum terverifikasi` kecuali `TGL_TRANSFER`.

### 6.2 Dibaca untuk menentukan **jalur banding** — bukti

```sql
-- Asal: NB FacIn/RDBList/GetAksepBanding_SQL.xml (rule GetAksepBanding_SQL)
select JABATAN AS CARI1, JABATAN_ATASAN AS CARI2
from POOLDATA.M_LIMIT_PROPERTYY
WHERE WORKBASKET IN (
  SELECT WORKBASKET FROM (
    SELECT WORKBASKET
    FROM HISTORYAKSEPTASIPEGA
    WHERE ID_PEGA ={pyWorkPage.pzInsKey}
      AND WORKBASKET !='ReasFacInMarketing'
    ORDER BY TGL_TRANSFER DESC
  ) WHERE ROWNUM = 1)
AND TEAM_GROUP = {pyWorkPage.OfferFacIn.QuotationData.TeamGroup}
```

Yang terjadi: **workbasket terakhir yang bukan Marketing** diambil dari riwayat, lalu dipakai
menjoin ke tabel limit untuk memperoleh `JABATAN` dan `JABATAN_ATASAN`. **Riwayat adalah masukan
keputusan, bukan jejak audit.** `[terverifikasi]`

> Perhatikan cacat Oracle klasik di sini: `ORDER BY` ada **di dalam** subquery inline dan `ROWNUM = 1`
> diterapkan pada query pembungkusnya — bentuk ini benar. Tetapi §6.3 memakai bentuk yang **tidak**
> benar.

### 6.3 Dibaca untuk menentukan **flag reject** — bukti

```sql
-- Asal: NB FacIn/RDBList/GetFlagReject_SQL.xml (rule GetFlagReject_SQL)
SELECT COUNT(*) AS CARI1 FROM HISTORYAKSEPTASIPEGA
WHERE ID_PEGA = {DataSearch.CARI20}
  AND WORKBASKET != 'ReasFacInMarketing'
  AND WORKBASKET != 'ReasFacInTeamLeader'
  AND STATUS ='REJECT'
ORDER BY TGL_TRANSFER DESC
```

Menghitung berapa kali case pernah ditolak di luar Marketing dan TeamLeader. `[terverifikasi]`

> **`ORDER BY` pada query agregat tanpa `GROUP BY` adalah tak bermakna** — `COUNT(*)` menghasilkan satu
> baris; mengurutkannya tidak melakukan apa-apa. Ini **sisa dari query yang pernah mengembalikan
> baris detail**. Dalam migrasi: tulis ulang tanpa `ORDER BY`, **tetapi jangan mengubah predikatnya** —
> hasilnya identik, jadi ini bukan perubahan perilaku. `[terverifikasi]`

> **`!=` dua kali, bukan `NOT IN`.** Setara secara logika **kecuali bila `WORKBASKET` bisa `NULL`** —
> dalam hal itu barisnya tersaring habis pada kedua bentuk. Tidak ada bukti korpus soal nullability.
> `belum terverifikasi`.

### 6.4 Konsekuensi migrasi

`CLAUDE.md` §4.3 sudah menyatakan tabel ini "dibaca untuk mengambil keputusan". Dokumen ini
**memasok buktinya**: dua query di atas. Implikasinya keras:

- `HISTORYAKSEPTASIPEGA` **bukan** tabel audit yang boleh dipindah ke skema log, diarsipkan, atau
  dipangkas retensinya. Memangkas riwayat **mengubah keputusan** yang diambil `GetAksepBanding_SQL`.
- Urutan `TGL_TRANSFER` adalah **semantik**, bukan kosmetik. Dua baris dengan `sysdate` identik
  (penyisipan dalam detik yang sama) membuat "workbasket terakhir" **tidak deterministik**.
  `[pertanyaan terbuka]` — apakah ada kolom urut lain, atau presisi `TGL_TRANSFER` cukup?
  Ini **memblokir** bila jawabannya "tidak".

### 6.5 Tabel riwayat yang **lain** — jangan tertukar

| Tabel | Ditulis oleh | Dibaca untuk keputusan? |
| --- | --- | --- |
| `HISTORYAKSEPTASIPEGA` | `InsertHistoryAkseptasiPega_Sql` | **Ya** — §6.2, §6.3 |
| `POOLDATA.HISTORYAKSEPTASIPRODUCTION` | `InsertViewSuggest_SQL` | Tidak ada `SELECT` di korpus `[terverifikasi]` |
| `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK` | (internal Pega) | **Ya** — §7.1 |

`POOLDATA.HISTORYAKSEPTASIPRODUCTION`, kolom dari `NB FacIn/RDBList/InsertViewSuggest_SQL.xml`:

```
IDPEGA, TYPE_POLIS, NOURUT, POSISI, PIC, TGL_INP, DIV, TYPE, PUTARAN, APPROVAL, KETERANGAN,
AKSES_LOGIN, B2B, BUSINESS_CODE, PERCENT_RNM
```

`[terverifikasi]`. `KETERANGAN` dipotong `substr({InputData.CARI10},0,3990)` → kolom **paling banyak
4000 karakter** `[terverifikasi]`; tipe pastinya `belum terverifikasi` (kemungkinan `VARCHAR2(4000)`).
`PIC` dan `AKSES_LOGIN` memuat identitas orang — berlaku peringatan §5.1.

---

## 7. Tabel internal Pega — **hilang bersama Pega**

### 7.1 `DATAPEGA.PC_*`

Tujuh rujukan, dua tabel, tiga rule.

| Tabel | Rule | Peran |
| --- | --- | --- |
| `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK` | `GetHistoryAccPega_SQL`, `GetOPFacOut_Sql` | riwayat internal work object Pega |
| `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | `GetCountClaim` | work object aplikasi **klaim** (GCNMFW), bukan Fac In |

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$d = $rows | Where-Object { $_.SQL -match '(?i)DATAPEGA\.' }
"blok: $($d.Count) ; rule unik: " + (($d | ForEach-Object { $_.File } | Sort-Object -Unique).Count)
$d | ForEach-Object { $_.File } | Sort-Object -Unique
```

Kolom terbaca (seluruhnya kolom internal Pega, berawalan `PX`/`PY`):

```sql
-- Asal: NB FacIn/RDBList/GetOPFacOut_Sql.xml (rule GetOPFacOut_Sql)
select PYPERFORMER as "pyOwnerUserID" from datapega.pc_history_asm_fw_gisfw_work
  where (PYMESSAGEKEY like '%OfferFacOut%' or PYMESSAGEKEY like '%OfferRetro%')
    and PYHISTORYTYPE='F'
    and PXHISTORYFORREFERENCE={pyWorkPage.pzInsKey}
  order by PXCOMMITDATETIME
```

```sql
-- Asal: NB FacIn/RDBList/GetHistoryAccPega_SQL.xml (rule GetHistoryAccPega_SQL)
select PXSAVEDATETIME AS TGL_TRANSFER from DATAPEGA.pc_History_ASM_FW_GISFW_Work
 where PXHISTORYFORREFERENCE ={pyWorkPage.pzInsKey}
   and PYASSIGNEDTO is not null and pylabel is not null
   and (PXADDEDBYID not in (…) or PXADDEDBYID is null)
 order by PXSAVEDATETIME
```

Kolom: `PXHISTORYFORREFERENCE`, `PXSAVEDATETIME`, `PXCOMMITDATETIME`, `PXADDEDBYID`, `PYASSIGNEDTO`,
`PYLABEL`, `PYPERFORMER`, `PYMESSAGEKEY`, `PYHISTORYTYPE`, `MASTERID`. `[terverifikasi]`

> **Guard berbasis identitas.** Predikat `PXADDEDBYID not in (…)` menyaring atas **1 identitas
> operator + 1 literal sistem**. Sesuai `CLAUDE.md` §3.5, nilainya tidak disalin ke sini. Mekanismenya
> saja: penyaringan riwayat dengan daftar-hitam identitas, ditulis langsung di dalam SQL. Ini satu-
> satunya guard identitas berbentuk literal SQL di seluruh 619 blok.
>
> ```powershell
> $rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
> ($rows | Where-Object { $_.SQL -match "(?i)(PXADDEDBYID|USERNAME|LOGIN|PYPERFORMER)\s*(NOT\s+)?IN\s*\(\s*'" } | ForEach-Object { "$($_.Cycle)/$($_.File)" })
> ```
>
> `[terverifikasi]` — tepat 3 blok, satu nama rule, hadir di ketiga siklus.

**Konsekuensi migrasi — MEMBLOKIR.** Ketiga query ini membaca tabel internal engine Pega. **Tabel itu
lenyap saat Pega dimatikan.** Yang harus diputuskan bisnis:

1. `GetOPFacOut_Sql` mengembalikan **siapa pemilik langkah Fac Out / Retro** — dipakai untuk routing.
   Di aplikasi baru informasi ini harus berasal dari tabel milik aplikasi sendiri. **Tabel itu belum
   ada.**
2. `GetHistoryAccPega_SQL` mengembalikan **daftar waktu transfer** — sebagian tumpang tindih dengan
   `HISTORYAKSEPTASIPEGA.TGL_TRANSFER`, tetapi **tidak identik** (yang ini menyaring `PYASSIGNEDTO`
   dan `PYLABEL` non-null). `[pertanyaan terbuka]`: dapatkah `HISTORYAKSEPTASIPEGA` menggantikannya
   sepenuhnya? Bila tidak, data historis **tidak dapat direkonstruksi setelah Pega mati** —
   yang berarti **ekstraksi harus dilakukan sebelum dekomisioning.**
3. `GetCountClaim` membaca work object **aplikasi klaim (GCNMFW)** — sistem yang berbeda, di luar
   lingkup migrasi ini. `[pertanyaan terbuka]`: apakah aplikasi klaim juga bermigrasi? Bila tidak,
   hitungan klaim harus datang lewat integrasi, bukan query lintas-tabel.

### 7.2 Delapan blok `pySaveSQL` — seluruhnya **kode mati**

Kedelapan blok `pySaveSQL` di seluruh korpus **identik byte-per-byte** dan tidak ada hubungannya
dengan nama rule yang memuatnya:

| Rule pemuat | Nama menjanjikan | Isi `pySaveSQL` sebenarnya |
| --- | --- | --- |
| `CariBusinessGID` (RNW, EDM) | pencarian ID bisnis | `POOLDATA.PROSESCOPY('2007','10018','1000134', errmsg)` |
| `ConvertBusinessField` (NB, RNW, EDM) | konversi field bisnis | idem, literal identik |
| `ConvertNationality` (NB, RNW, EDM) | konversi kewarganegaraan | idem, literal identik |

```sql
-- Asal: NB FacIn/RDBList/ConvertBusinessField.xml, tab pySaveSQL
DECLARE
   vTHN_TREATY VARCHAR2(10);
   vTOP_ID VARCHAR2(10);
   vIDTreatyYear VARCHAR2(10);
   errmsg VARCHAR2(4000);
BEGIN
   vTHN_TREATY := '2007';
   vTOP_ID := '10018';
   vIDTreatyYear:='1000134';
   POOLDATA.PROSESCOPY(vTHN_TREATY,vTOP_ID,vIDTreatyYear,errmsg);
   dbms_output.put_line(errmsg);
END;
```

`[terverifikasi]`

Ini **template Pega yang tidak pernah dibersihkan**: nilai keras `'2007'`, `'10018'`, `'1000134'`,
dan `dbms_output.put_line` (keluaran konsol yang tak dibaca siapa pun di runtime aplikasi).

Tiga hal yang perlu dicatat:

1. **Tidak ada bukti blok ini pernah dieksekusi.** Pega memanggil `pySaveSQL` hanya lewat metode
   `RDB-Save`; korpus memakai `RDB-Save` **nol kali** (hanya `RDB-List` 614×, `RDB-Delete` 2× di
   `NB FacIn/Activity`). `[terverifikasi]`
2. Contoh paling murni dari aturan **"nama rule bukan bukti perilaku"**: tiga nama rule yang berbeda
   sama sekali, satu isi identik yang tak berkaitan dengan ketiganya.
3. **Jangan memigrasikan blok ini.** Ia tidak melakukan apa pun yang dibutuhkan aplikasi. Prosedur
   `POOLDATA.PROSESCOPY` sendiri **tidak boleh dipanggil** dari aplikasi baru sampai bisnis
   menjelaskan apa fungsinya — isinya `belum terverifikasi`.

```powershell
foreach ($c in @('NB FacIn','RNW Fac In','Endorsment Fac In')) {
  $n = (Get-ChildItem "D:\migrasi\RNM\$c\Activity" -File -Filter *.xml | ForEach-Object {
    ([regex]::Matches((Get-Content $_.FullName -Raw -Encoding UTF8),'<pyStepsActivityName>RDB-Save</pyStepsActivityName>')).Count
  } | Measure-Object -Sum).Sum
  "{0,-20} RDB-Save = {1}" -f $c, $n
}
```

### 7.3 Rule bernama "Insert…" yang **menghapus lalu menyisip ulang** — ditemukan

Diminta secara eksplisit. **Ditemukan, dua kejadian.**

`Rule-Connect-SQL` punya tab `Browse` dan tab `Delete` yang terpisah. Dua rule mengisi **keduanya**:

| Rule | `pyBrowseSQL` | `pyDeleteSQL` |
| --- | --- | --- |
| `InsertOfferProduction_Sql` | `INSERT INTO POOLDATA.FACINOFFER (… 50 kolom …)` + `COMMIT` | `DELETE FROM POOLDATA.FACINOFFER WHERE IDPEGA ={pyWorkPage.pyID}` |
| `InsertCedingProduction_SQL` | `BEGIN pooldata.InsertUpdateCedingProduction(…); COMMIT; END;` | `DELETE FROM POOLDATA.CEDING_FACINPRODUCTION WHERE IDPEGA ={pyWorkPage.pzInsKey}` |

`[terverifikasi]` — kedua tab dikutip di §2.5 dan §2.7.

Urutan eksekusinya dibuktikan dari activity pemanggil:

```powershell
$t = Get-Content "D:\migrasi\RNM\NB FacIn\Activity\SaveOfferProduction_Act.xml" -Raw -Encoding UTF8
$evts = @()
foreach ($m in [regex]::Matches($t,'RDB-Delete|RDB-List|InsertOfferProduction_Sql')) {
  $evts += [pscustomobject]@{Pos=$m.Index; Tok=$m.Value}
}
$evts | Sort-Object Pos | Select-Object -First 12 | ForEach-Object { "{0,8}  {1}" -f $_.Pos, $_.Tok }
```

Keluaran memperlihatkan `RDB-Delete` pada posisi 98806, diikuti rujukan `InsertOfferProduction_Sql`
pada posisi 100886 (**langkah yang sama**), lalu rujukan `InsertOfferProduction_Sql` berulang-ulang
pada langkah-langkah `RDB-List` berikutnya — sekali per baris spreading.

**Pola sebenarnya**: *hapus seluruh baris milik case ini → sisipkan ulang semuanya, satu `INSERT` +
`COMMIT` per baris.* `[terverifikasi]`

> **Dan label langkahnya berbohong.** Deskripsi langkah `RDB-Delete` di dalam
> `NB FacIn/Activity/SaveOfferProduction_Act.xml` berbunyi **`insert to db`**.
>
> ```powershell
> $t = Get-Content "D:\migrasi\RNM\NB FacIn\Activity\SaveOfferProduction_Act.xml" -Raw -Encoding UTF8
> $m = [regex]::Match($t,'(?s)<pyStepsActivityName>RDB-Delete</pyStepsActivityName>(.{0,3500})')
> [regex]::Match($m.Groups[1].Value,'<pyStepsDescription>(.*?)</pyStepsDescription>').Groups[1].Value
> ```
>
> `[terverifikasi]`. Contoh kedua, dari sisi activity, untuk aturan "label bukan bukti".

Konsekuensi migrasi:

1. **Tidak idempoten dan tidak atomik.** `COMMIT` dijalankan **per baris** (117 blok memuat `COMMIT`).
   Kegagalan di tengah meninggalkan case dengan sebagian baris terhapus dan sebagian tersisip —
   keadaan yang **tidak dapat dibedakan** dari case yang memang punya sedikit baris.
2. Di Go: bungkus hapus-dan-sisip-ulang dalam **satu transaksi**, `COMMIT` sekali di akhir. Ini
   **mengubah perilaku kegagalan** (menjadi lebih baik) tanpa mengubah hasil pada jalur sukses —
   catat sebagai keputusan sadar, bukan perbaikan diam-diam.
3. Semua kunci alami hilang: karena baris dihapus dan disisip ulang, **tidak ada identitas baris yang
   stabil** di `FACINOFFER`. Jangan merancang foreign key ke sana.

---

## 8. Stored procedure dan fungsi yang dipanggil — **nama saja**

> Korpus memuat **pemanggilannya**, bukan **isinya**. Untuk setiap entri di bawah ini: parameter
> terbaca dari teks pemanggilan; **badan prosedur, efek samping, transaksi internal, dan trigger yang
> mungkin ikut terpicu seluruhnya `belum terverifikasi`.** Jangan menebak.

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
$proc = @{}
foreach ($r in $rows) {
  foreach ($m in [regex]::Matches($r.SQL,'(?i)(\bINSERT\s+INTO\s+|\bFROM\s+|\bJOIN\s+|\bUPDATE\s+)?\b([A-Z_][A-Z0-9_$]*\.[A-Z_][A-Z0-9_$]*)\s*\(')) {
    if ($m.Groups[1].Success -and $m.Groups[1].Value.Trim().Length -gt 0) { continue }   # buang target INSERT INTO / FROM
    $n = $m.Groups[2].Value.ToUpper()
    if ($n -match '^(DBMS_|UTL_)') { continue }                                          # buang paket bawaan Oracle
    if (-not $proc.ContainsKey($n)) { $proc[$n] = 0 }; $proc[$n]++
  }
}
"PROSEDUR/FUNGSI NON-BAWAAN: $($proc.Count)"
$proc.GetEnumerator() | Sort-Object Name | ForEach-Object { "{0,-42} {1}" -f $_.Key,$_.Value }
```

Hasil: **31** nama non-bawaan. `[terverifikasi]`

### 8.1 Prosedur `POOLDATA` (24) — dipanggil dalam blok `BEGIN … END`

| Prosedur | Rule pemanggil | Peran terbaca dari konteks |
| --- | --- | --- |
| `INSERTJSONPOLIS` | `INSERTJSON_JSONPOLIS_FACIN`, `INSERTJSON_JSONPOLISEDM_FACIN` | **satu-satunya jalur tulis ke `JSON_POLIS`** (§3) — 9 argumen |
| `INSERTJSONPOLISMONITORING` | `INSERTJSON_JSONPOLISMONITORING_FACIN` | payload `CLOB`, + `pxProcess.pzProductionLevel` |
| `PEGA_JSON_POLIS_TREATYIN` | `SavePolisTreatyIn_SQL` | varian Treaty In dari yang di atas — 8 argumen |
| `PEGA_M_JSON_OFFER` | `SaveOfferJson_SQL` | tulis `JSON_OFFER`; keluaran `ERRMSG`, `IDPEGAOUT`, `STSSAVE` |
| `PEGA_TREATY_IN` | `SaveTreatyIn` | 24 argumen, termasuk `StatusAkseptasi` |
| `INSERTUPDATECEDINGPRODUCTION` | `InsertCedingProduction_SQL` | §2.7 |
| `INSERTUPDATERISKADDRESS` | `UpdateMasterRiskAddress_SQL` | 9 argumen masuk, 2 keluar |
| `PEGA_M_ACCUMULATION_LIFE` | `SaveNewAccumulation_SQL`, `Update_sql` | menerima data pribadi — §9 |
| `PEGA_MARKETINGOFFICER` | `UpdateMasterMarketingOfficer` | |
| `PEGA_DELETE_ERROR_KONVERSI` | `DeleteDataProduction` | **tidak ada `COMMIT`** di blok pemanggil |
| `FACINFORBACKUP` | `MachingDataFacin_Sql` | 8 argumen kunci baris + 1 keluar |
| `PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL` | 3 masuk (termasuk tanggal), 2 keluar |
| `GET_TOKEN_STORAGE` | `GetTokenStorage_SQL` | **kredensial** — §10 |
| `RDBINSERTCLIENT` | `INSERTCONORGJSON_MCLIENT` | |
| `RDBMASTERACCUMULATEDTYPE` | `UpdateMasterAccumulatedType` | |
| `RDBMASTERACCUMULATION` | `UpdateMasterAccumulation` | |
| `RDBMASTERBRANCH` | `UpdateMasterBranch` | |
| `RDBMASTERCITY` | `UpdateMasterCity` | |
| `RDBMASTERDISTRICT` | `UpdateMasterDistrict` | |
| `RDBMASTERNATION` | `UpdateMasterNation` | |
| `RDBMASTERPROVINCE` | `UpdateMasterProvince` | |
| `RDBMASTERRW` | `UpdateMasterRW` | |
| `RDBMASTERSHIP` | `UpdateMasterShip` | hanya di 2 siklus |
| `PROSESCOPY` | `ConvertBusinessField`, `ConvertNationality`, `CariBusinessGID` | **kode mati** — §7.2 |

Sembilan prosedur `RDBMASTER*` + `RDBINSERTCLIENT` + `INSERTUPDATERISKADDRESS` + `PEGA_MARKETINGOFFICER`
membentuk satu kelompok: **sinkronisasi data master dari Pega ke Oracle**. Arahnya
(Pega → Oracle, atau dua arah) `belum terverifikasi`.

### 8.2 Fungsi yang dipanggil di dalam `SELECT` (7) — baca saja

| Fungsi | Rule pemanggil | Bentuk |
| --- | --- | --- |
| `POOLDATA.GETCURRENCYSTANDARD` | `CurrencyStandard` | `select pooldata.getcurrencystandard({…Currency}, sysdate) as nilaiKurs from dual` |
| `POOLDATA.GENERATE_FACRETRO_NO` | `GenerateRISlipNumber` | `SELECT POOLDATA.GENERATE_FACRETRO_NO(…4 arg…) AS FacRetroSlipNumber FROM DUAL` |
| `GENERAL.F_GET_NM_ASURADUR` | `SearchCoinsSQL` | 21 rujukan dalam satu query |
| `FIRE.CEK_PRORATA_TANGGAL` | `SearchSQLRateKPR` | |
| `FIRE.PEGA_FIRE_SET_RATE` | `SearchSQLRateNonKPR` | |
| `MBU.F_CEK_HURUF` | `GetNoRangkaMesinDiff` | skema `MBU` — hanya di sini |
| `NEW_GENERAL.CEK_PLAT_NO` | `SearchTemplateMainCoverageSQL` | skema `NEW_GENERAL` — hanya di sini |

`[terverifikasi]`

> Dua skema (`MBU`, `NEW_GENERAL`) **muncul hanya lewat pemanggilan fungsi ini** dan tidak pernah
> sebagai pemilik tabel. Keduanya **tidak ada di daftar §1.6** karena §1.6 hanya menghitung prefix pada
> `FROM`/`JOIN`/`INSERT INTO`/`UPDATE`. Dependensi lintas-skema aplikasi karena itu **lebih luas dari
> enam skema di §1.6** — totalnya **delapan**. `[terverifikasi]`

> `POOLDATA.GETCURRENCYSTANDARD(currency, sysdate)` adalah **satu-satunya sumber kurs** yang terbaca.
> Ia mengambil `sysdate`, bukan tanggal efektif polis — artinya **konversi mata uang memakai kurs hari
> ini, bukan kurs pada tanggal transaksi**. `[terverifikasi]` bahwa argumennya `sysdate`;
> `[pertanyaan terbuka]` apakah ini disengaja. Ini material untuk rekonsiliasi paralel run: **dua
> sistem yang dijalankan pada hari berbeda akan menghasilkan angka berbeda.**

### 8.3 Paket bawaan Oracle yang dipakai

`DBMS_LOB.CREATETEMPORARY` (30 rujukan), `DBMS_OUTPUT.PUT_LINE` (8 — seluruhnya di kode mati §7.2),
`UTL_MATCH.EDIT_DISTANCE_SIMILARITY` (4 — pencarian akumulasi fuzzy, §11.2).

---

## 9. Tabel lookup / master

112 tabel dikurangi yang sudah dibahas. Dikelompokkan menurut fungsi; **kolom hanya dicantumkan bila
SQL mengutipnya**.

### 9.1 Master geografi dan akumulasi risiko

| Tabel | Kolom terbaca | Asal |
| --- | --- | --- |
| `POOLDATA.RW` | `NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION, ZIPCODE, CZONE` | `BrowseRW_SQL`, `BrowseRW2_SQL`, `BrowseZipCodeAndCzone` |
| `CZONE` | `ID, CODE` | `BrowseZipCodeAndCzone` |
| `ACCUMULATION` | `ID, ACCUMULATIONTYPE, NOTE, ZIPCODE, ACCUMULATIONNAME, KEYWORD, CZONE` | `GetAccumulationByNote_SQL`, `GetSummaryRiskAccumulation_Sql` |
| `ZONES`, `POOLDATA.ZONEOJK` | `belum terverifikasi` | |
| `M_FLOOD_AREA`, `M_BRANCH`, `M_COUNTRY` | `belum terverifikasi` | |
| `GENERAL.LST_DET_CABANG` | `belum terverifikasi` | |

`ACCUMULATION` adalah salah satu tabel terpanas (34 rujukan). `[terverifikasi]`

### 9.2 Master mitra dan produk

| Tabel | Catatan |
| --- | --- |
| `POOLDATA.CLIENT`, `POOLDATA.AGENT` | pihak; ditulis lewat `RDBINSERTCLIENT` |
| `GENERAL.LST_ASURADUR` | daftar perusahaan asuransi (9 rujukan) |
| `MARKETINGOFFICER` | 15 rujukan. Kolom: `ID, CLIENTID, CLIENTID2, CLIENTNAME, MOLEADER, TEAMGROUP, MOSTATUS, BRANCHSTATUS` `[terverifikasi]` (`GetDataMarketing_SQL`). Predikat aktif: `MOSTATUS = '1' and BRANCHSTATUS is null AND CLIENTID2 <> 'LEADER'` — kode `'1'` dan `'LEADER'` **belum terverifikasi** artinya. Memuat identitas orang — peringatan §5.1 berlaku |
| `POOLDATA.BUSINESS`, `M_BUSINESSFIELD`, `M_TREATYBUSINESS` | klasifikasi lini bisnis |
| `POOLDATA.REINSURANCETYPE`, `M_REINSURANCETYPE` | kolom `ID`, `NOTE` `[terverifikasi]` — `(select note from pooldata.reinsurancetype where id = jn_reas)` |
| `POOLDATA.OCCUPATION`, `NEW_UNDERWRITING.M_OCCUPATION`, `POOLDATA.OCUPATIONHANNOVER` | **tiga** tabel okupasi *(perhatikan ejaan `OCUPATION…`, satu huruf C)* |
| `POOLDATA.CURRENCY`, `M_CURRENCY`, `LST_KURS_STANDARD` | **tiga** tabel mata uang |
| `COVERAGE`, `COVERAGE_FACIN`, `NEW_UNDERWRITING.M_COVERAGE` | **tiga** tabel coverage |
| `M_CLAUSE`, `CLAUSE`, `LST_KLAUSULA`, `M_ARGCLAUSEFIRE` | **empat** tabel klausula |
| `M_KONSTRUKSI`, `CONSTRUCTION` | **dua** tabel konstruksi |
| `OBJECTITEMTYPE`, `V_JN_OBJ_ITEM` | jenis objek; `V_JN_OBJ_ITEM` punya `MJOI_KODE`, `JN_OBJ_ITEM` `[terverifikasi]` |
| `POOLDATA.CAUSEOFDECLINE` | alasan penolakan |
| `POOLDATA.MASTERCARGO`, `POOLDATA.M_PLANTRAVEL`, `POOLDATA.M_RISK` | |
| `LOADINGUSIAKENDARAAN`, `NEW_UNDERWRITING.M_MAPPING_PLAT_KEND` | kendaraan bermotor |
| `GENERAL.V_JOB_PA`, `POOLDATA.VJ_PKG_BENEFIT`, `POOLDATA.VJ_BENEFIT_PROPERTY`, `POOLDATA.VJ_PACKAGE_AGE_KLAUSUL`, `POOLDATA.VJ_PACKAGE_AGE_KLAUSUL_DM` | view (`V_`/`VJ_`) |
| `NEW_UNDERWRITING.M_DEDUCTIBLE_TYPE`, `M_RECEIVER_OUTGO`, `T_TEMPLATE_*` (5 tabel) | template underwriting |

> **Duplikasi master adalah risiko migrasi tersendiri.** Tiga tabel okupasi, tiga mata uang, tiga
> coverage, empat klausula. Mana yang otoritatif — dan apakah isinya konsisten — **`belum
> terverifikasi` dan tidak dapat dijawab dari korpus.** `[pertanyaan terbuka]`, lihat §12.

### 9.3 Master tarif

`M_EQS_RATE`, `POOLDATA.M_EQS_MULTIPLIER`, `M_FLEXAS_RATE`, `M_RSMD_RATE`, `M_TERORISME_RATE`,
`FIRE.M_FLOOD_RATE`, `FIRE.M_KPR_RATE`, `FIRE.T_KOMISI_KPR`, `FIRE.M_BI_INDEMNITY`, `RATE_LIFE`,
`SPREADSYARIAH`, `M_RISK_LOSS_PROFILE`, `RIRISK_LIFE`.

Kolom untuk sebagian besar **`belum terverifikasi`** — query berbentuk `SELECT *`.

`M_RISK_LOSS_PROFILE` terverifikasi lewat `INSERT`-nya
(`NB FacIn/RDBList/InsertRiskAndLossProfile_SQL.xml`):

```
NOPOLIS, TANGGALSPREADING, TSINUSANTARARE, PREMINUSANTARARE, PREMIUMOR, PREMIUMQS, PREMIUMFACOUT, COB
```

`[terverifikasi]`

### 9.4 Treaty

`TREATYBUSINESS`, `TREATYCONTRACT`, `TREATYYEAR`, `POOLDATA.TREATYGROUP`,
`POOLDATA.TREATYEXCHANGEYEARLY` (21 rujukan), `POOLDATA.TREATYINPRODUCTION`,
`POOLDATA.PROPORTIONALARRG`, `POOLDATA.M_TREATY_IN`, `POOLDATA.M_TREATY_IN_EDM`.

### 9.5 Operasional dan monitoring

| Tabel | Kolom terbaca | Asal |
| --- | --- | --- |
| `POOLDATA.ERRORFACINPROD` | `IDPEGA, NOPOLIS, NOENDORS, ERRORMESSAGE, TGL_INPUT, STATUS` | `INSERTERRORFACIN_Sql` `[terverifikasi]` |
| `POOLDATA.MONITORING_PROD_LOG` | `IDPEGA, NOPOLIS, PARAMETER, JN_SERVICE, STS_MESSAGE, RESPON_MESSAGE` | `InsertLogServiceProd` `[terverifikasi]` |
| `POOLDATA.FACINPERFORMANCE` | `IDPEGA, POSITIONDOC, STATUSDOC, DATETIME, USERNAME, INSUREDNAME, SOB, MOLEADER, COB, RISLIPRECEIVED` | `InsertDataFacinPerformance_SQL` `[terverifikasi]` |
| `POOLDATA.TANGGAL_CLOSING` | `SELECT *` → `belum terverifikasi` | `GETTanggalClosing_SQL` |
| `POOLDATA.KODE_PRODUKSI` | `KODE, TYPE` | `GetKodeProdLife_SQL`: `WHERE TYPE ='LIFE'` `[terverifikasi]` |
| `POOLDATA.DATAKLAIM` | `belum terverifikasi` (9 rujukan) | |
| `OS_AKSEPTASI_KLAIM` | `belum terverifikasi` | |
| `ARASAPAS.DETAIL_INVOICE` / `DETAIL_INVOICE` | `belum terverifikasi` | skema `ARASAPAS` hanya muncul sekali |

### 9.6 Dokumen dan lampiran

| Tabel | Kolom terbaca | Asal |
| --- | --- | --- |
| `T_STORAGE_IMAGE` | `IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE, TANGGAL_UPLOAD` | `Insert_T_Storage_SQL`, `Update_T_Storage_SQL`, `GetLinkStorage_SQL` `[terverifikasi]` |
| `POOLDATA.T_FOLDER_IMAGE` | `APPNAME` | `GetAppName_SQL` `[terverifikasi]` |
| `CATEGORY_ATTACH_REAS` | `NOTE` | `AttachmentLife`: `where note in (…13 nilai…)` `[terverifikasi]` |

`T_STORAGE_IMAGE.STORAGE` diisi literal `'standard'` saat penyisipan `[terverifikasi]` — menyiratkan
ada nilai lain, yang **`belum terverifikasi`**.

`T_STORAGE_IMAGE.IMAGEID` dibangkitkan di database:

```sql
-- Asal: NB FacIn/RDBList/GenerateImageID_SQL.xml
SELECT STANDARD_HASH('ASMPP' || TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF9') || SYS_GUID(), 'MD5')
  AS "InsertDoc.ImageID" FROM DUAL
```

> **Siklus endorsement memakai rumus yang berbeda** — salah satu dari 3 rule yang menyimpang (§2.2):
> presisi `'DD/MM/YYYY HH24:MI:SS.FF3'` dan **tanpa `SYS_GUID()`**. Tanpa GUID, dua unggahan dalam
> milidetik yang sama menghasilkan `IMAGEID` yang **sama** → tabrakan primary key atau penimpaan
> berkas. **Kandidat perbaikan, bukan bagian migrasi.** `[terverifikasi]`

---

## 10. Tabel lookup endpoint / link service

### 10.1 `M_LINK_SERVICE`

**Tidak dirujuk lewat SQL sama sekali** — ia diakses lewat `Obj-Browse` Pega, sehingga tidak muncul
dalam 112 tabel di §1.4.

**Asal**: `NB FacIn/Activity/GetLinkService.xml`, rule `GetLinkService`, kelas
`ASM-FW-GISFW-Int-M_LINK_SERVICE`.

```powershell
$t = Get-Content "D:\migrasi\RNM\NB FacIn\Activity\GetLinkService.xml" -Raw -Encoding UTF8
"Field : " + (([regex]::Matches($t,'<Field>(.*?)</Field>') | ForEach-Object { $_.Groups[1].Value }) -join ', ')
"Cond  : " + (([regex]::Matches($t,'<Condition>(.*?)</Condition>') | ForEach-Object { $_.Groups[1].Value }) -join ', ')
"Steps : " + (([regex]::Matches($t,'<pyStepsActivityName>(.*?)</pyStepsActivityName>') | ForEach-Object { $_.Groups[1].Value }) -join ' -> ')
```

| Aspek | Nilai |
| --- | --- |
| Tabel | `M_LINK_SERVICE` |
| **Kolom kunci** | `KATEGORI_1`, `KATEGORI_2` — keduanya dibandingkan `=` |
| **Kolom hasil** | `URL` |
| Langkah activity | `Page-New` → `Obj-Browse` → `Property-Set` → `Page-Remove` |

`[terverifikasi]`. Ini **mengonfirmasi** pernyataan `CLAUDE.md` §4.4 dari sisi korpus: kunci memang
`KATEGORI_1` + `KATEGORI_2`.

**Isi tabel tidak ada di korpus** dan tidak dicatat di sini. Aktivitasnya membuang halaman hasil
(`Page-Remove`) setelah menyalin `URL` ke satu properti.

**Untuk migrasi**: seluruh endpoint aplikasi baru dibaca dari tabel ini dengan kunci ganda yang sama,
atau dipindah ke konfigurasi/env var. Keduanya sah menurut `CLAUDE.md` §4.4; yang **tidak** sah adalah
menuliskan URL sebagai literal di kode.

### 10.2 `M_PROMPT_AI` — bentuk kunci yang sama, isi yang jauh lebih sensitif

```sql
-- Asal: NB FacIn/RDBList/GetPromptAI_SQL.xml (rule GetPromptAI_SQL)
SELECT PROMPT_AI AS "ParamAI.PromptAI" FROM M_PROMPT_AI
WHERE KATEGORI_1 = {ParamAI.CARI1}
  AND KATEGORI_2 = {ParamAI.CARI2}
```

| Kolom | Peran |
| --- | --- |
| `KATEGORI_1`, `KATEGORI_2` | kunci ganda — **pola identik dengan `M_LINK_SERVICE`** |
| `PROMPT_AI` | teks prompt yang dikirim ke model AI |

`[terverifikasi]`

**Ini material dan berkaitan langsung dengan `CLAUDE.md` §6.** Prompt yang menggerakkan
`GeminiAIGoogle_Act` di alur underwriting **tidak ada di kode** — ia **data di tabel Oracle**, dapat
diubah tanpa deployment dan tanpa jejak di version control. Konsekuensinya:

- Perilaku model AI dalam alur underwriting **tidak dapat direview dari korpus**; isinya
  `belum terverifikasi`.
- Siapa pun dengan akses tulis ke `M_PROMPT_AI` dapat mengubah instruksi yang memengaruhi keputusan
  underwriting. **Tinjauan keamanan yang diwajibkan `CLAUDE.md` §6 harus mencakup tabel ini**, bukan
  hanya rule `GeminiAIGoogle_Act`.
- `[pertanyaan terbuka]` — apakah ada audit trail perubahan `M_PROMPT_AI`? Tidak ada buktinya di
  korpus.

### 10.3 Kredensial: `POOLDATA.GET_TOKEN_STORAGE`

```sql
-- Asal: NB FacIn/RDBList/GetTokenStorage_SQL.xml (rule GetTokenStorage_SQL)
BEGIN
  pooldata.GET_TOKEN_STORAGE ( {UploadDoc.App}, {OperatorID.pyUserIdentifier},
                               {UploadDoc.Kodestring OUT}, {DocAPI.ResponseMsg OUT});
  COMMIT;
END;
```

Prosedur ini mengembalikan **token** (`Kodestring OUT`) untuk penyimpanan dokumen, berdasarkan nama
aplikasi dan identitas operator. **Isinya `belum terverifikasi`** — termasuk apakah token dibangkitkan,
diambil dari tabel, atau dipanggil ke layanan luar.

`[terverifikasi]` bahwa prosedurnya dipanggil dan mengembalikan token. **Setiap penanganan hasil
prosedur ini di aplikasi baru memerlukan persetujuan manusia** (`CLAUDE.md` §6: akses credential).

### 10.4 Tidak ada literal URL di dalam SQL

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
($rows | Where-Object { $_.SQL -match '(?i)https?://' }).Count    # -> 0
```

**0 blok dari 619.** `[terverifikasi]` — temuan yang menyenangkan: lapisan SQL bersih dari endpoint
tertanam. (Rule `SystemSettings\LinkService.xml` memang ada di ketiga siklus dan memuat jenjang
lingkungan Sandbox/Development/QA/Staging/Production, tetapi `pyValue`-nya dalam ekspor berisi satu
karakter dan bukan URL. Keberadaannya dicatat; nilainya tidak.)

---

## 11. Risiko teknis pada lapisan data

### 11.1 Interpolasi mentah `{ASIS:…}` — injeksi SQL

Pega mensubstitusi `{Page.Prop}` sebagai **bind variable**, tetapi `{ASIS:Page.Prop}` sebagai
**teks mentah**. Korpus memakai bentuk kedua di **9 blok / 4 nama rule**:

```powershell
$rows = Import-Clixml "$env:TEMP\sqlblocks.xml"
"blok  : " + ($rows | Where-Object { $_.SQL -match '\{ASIS:' }).Count
"rule  : " + (($rows | Where-Object { $_.SQL -match '\{ASIS:' } | ForEach-Object { $_.File } | Sort-Object -Unique) -join ', ')
```

| Rule | Bentuk | Tingkat bahaya |
| --- | --- | --- |
| `BrowseDataViewPolis1` | `WHERE ({ASIS:ParamViewIn.CARI5})` | **seluruh klausa `WHERE`** datang dari halaman |
| `GetDataMarketing_SQL` | `CLIENTNAME LIKE '{ASIS:Marketing.CARI1}%'` — di dalam literal `LIKE` | tinggi |
| `SearchAccumulationbypersetase_SQL` | `UPPER('{ASIS:ParamSearch.CARI1}')` — di dalam literal berkutip | tinggi |
| `CekDuplicateLocationOffer_SQL` | `UPPER('{ASIS:ParamCari.CARI3}')` ×2 | tinggi |

`[terverifikasi]`

**Untuk migrasi**: keempatnya **wajib** ditulis ulang sebagai query berparameter di Go. Untuk
`BrowseDataViewPolis1`, klausa `WHERE` dinamis harus dibangun dari **daftar-putih kolom dan operator**,
bukan konkatenasi string. Ini perbaikan keamanan, bukan perubahan perilaku — hasil pada masukan sah
tetap sama.

### 11.2 Pencocokan fuzzy sebagai logika bisnis

```sql
-- Asal: NB FacIn/RDBList/SearchAccumulationbypersetase_SQL.xml
WHERE UTL_MATCH.EDIT_DISTANCE_SIMILARITY(A.CLEAN_NOTE, P.CLEAN_PARAM) >= 70
```

Ambang **70** dan daftar kata yang dibuang
(`DESA|DUSUN|GANG|GEDUNG|JL\.|KOMPLEK|PERUMAHAN|OTHERS`) tertanam di dalam SQL. `[terverifikasi]`

`CekDuplicateLocationOffer_SQL` memakai ambang yang **berbeda**, datang dari parameter:
`HAVING COUNT(DISTINCT W.WORD) >= CEIL( MAX(W.TOTAL_WORD) * to_number({ParamCari.CARI4}) )`.
`[terverifikasi]`

Keduanya mendeteksi akumulasi risiko pada lokasi yang sama. **`UTL_MATCH` adalah paket Oracle**; jika
logika ini dipindahkan ke Go, hasilnya **tidak akan identik** kecuali algoritma jarak sunting yang
sama direplikasi persis. Rekomendasi: **pertahankan di SQL**.

### 11.3 `COMMIT` per operasi

117 dari 619 blok memuat `COMMIT` di dalam teks SQL. Aplikasi **tidak mengelola transaksi**; setiap
blok menutup transaksinya sendiri. Digabung dengan pola hapus-sisip-ulang (§7.3), tidak ada satu pun
operasi tulis multi-baris yang atomik. `[terverifikasi]`

### 11.4 Angka masuk sebagai string berkoma desimal

`To_number(Replace({…},',','.'))` muncul 21× dalam satu `INSERT` saja. Ini mengonfirmasi
`CLAUDE.md` §4.1 dari sisi database. Di Go: `decimal.Decimal`, **tidak pernah `float64`**; pada
batas API, string desimal. `[terverifikasi]`

---

## 12. Catatan migrasi per tabel

### 12.1 Dibawa apa adanya — skema tidak berubah

| Tabel | Alasan |
| --- | --- |
| `POOLDATA.FACINPRODUCTION` | tabel produksi inti; pertahankan 82 kolom, termasuk salah ketik `PCT_BROKERGARE_FEE` |
| `POOLDATA.FACINOFFER` | idem; **periksa tipe kolom dulu** (§2.5) |
| `FACOUTPRODUCTION` | produksi Fac Out |
| `POOLDATA.JSON_POLIS` | riwayat versi polis; `PRODKE` tetap penomoran versi, baris bertambah |
| `HISTORYAKSEPTASIPEGA` | **dibaca untuk keputusan** (§6) — bentuk dipertahankan, **jangan diarsipkan** |
| Seluruh 7 tabel `M_LIMIT_*` | tangga persetujuan; ejaan huruf ganda dipertahankan |
| `M_LINK_SERVICE`, `M_PROMPT_AI` | kunci `KATEGORI_1`+`KATEGORI_2`; isinya konfigurasi runtime |
| Seluruh master §9 | lookup; tidak tersentuh migrasi |
| `POOLDATA.TREATYPRODUCTION_BACKUP` | cadangan; pemicunya perlu dikonfirmasi |
| `T_STORAGE_IMAGE`, `POOLDATA.T_FOLDER_IMAGE` | metadata dokumen |

### 12.2 Hilang bersama Pega — perlu pengganti

| Tabel | Dipakai untuk | Yang harus disiapkan |
| --- | --- | --- |
| `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK` | pemilik langkah Fac Out/Retro (`GetOPFacOut_Sql`); daftar waktu transfer (`GetHistoryAccPega_SQL`) | tabel riwayat langkah milik aplikasi baru — **belum ada**; **ekstrak data historis sebelum dekomisioning** |
| `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | jumlah klaim (`GetCountClaim`) | integrasi ke sistem klaim, bukan query lintas-tabel |

**Jangan memigrasikan kedua tabel ini apa adanya** (`CLAUDE.md` §4.3).

### 12.3 Jangan dimigrasikan

| Objek | Alasan |
| --- | --- |
| 8 blok `pySaveSQL` (`PROSESCOPY` dengan literal keras) | kode mati, tidak pernah dieksekusi (§7.2) |
| `ORDER BY` pada `GetFlagReject_SQL` | tak bermakna pada agregat; hasil identik tanpanya (§6.3) |

### 12.4 Perlu keputusan bisnis sebelum ditulis kodenya

| Hal | Rujukan |
| --- | --- |
| Tipe kolom `FACINOFFER` (`NUMBER` vs `VARCHAR2`) | §2.5 — **MEMBLOKIR** |
| `PRODKE`: `COUNT-1` vs `MAX(PRODKE)` | §3.2 |
| Kolom PII `NO_KTP` / `NO_NPWP` di `FACINPRODUCTION` | §2.2 |
| Arah tanda `_SELISIH` | §4.3 |
| Kurs `sysdate` vs tanggal efektif | §8.2 |
| Master ganda (okupasi ×3, mata uang ×3, coverage ×3, klausula ×4) | §9.2 |

---

## 13. Pertanyaan terbuka

Ditandai **[MEMBLOKIR]** bila kode tidak dapat ditulis dengan benar tanpa jawabannya.

### Memblokir

1. **[MEMBLOKIR] Tipe kolom seluruh tabel.** Korpus tidak memuat DDL. Khususnya: apakah
   `FACINOFFER.RATE`, `FACINOFFER.TSI_MENJADI`, `FACINOFFER.PREMI_MENJADI` bertipe `NUMBER` atau
   `VARCHAR2`? `InsertOfferProduction_Sql` menyisipkan **tanpa satu pun `To_number`**, sedangkan
   `InsertTreatyProduction_Sql` memakai 21. Bila `VARCHAR2`, seluruh perbandingan angka atasnya adalah
   perbandingan string dan **harus dipertahankan** (§2.5). *Jalankan `all_tab_columns` atas
   `POOLDATA.FACINOFFER`, `FACINPRODUCTION`, `FACOUTPRODUCTION`.*

2. **[MEMBLOKIR] Tipe dan skema `JSON_POLIS.DATA_JSON`.** `JSON`, `CLOB IS JSON`, `BLOB`, atau tipe
   objek? Perilaku notasi titik berbeda untuk masing-masing. Dan: **apa struktur lengkap dokumennya?**
   Korpus hanya memperlihatkan 11 jalur yang dibaca SQL; dokumen sesungguhnya adalah serialisasi
   agregat `OfferFacIn` yang jauh lebih besar dan **tidak terdokumentasi** (§3.3).

3. **[MEMBLOKIR] Isi 24 stored procedure `POOLDATA`.** Seluruh jalur tulis ke `JSON_POLIS`,
   `JSON_OFFER`, `M_ACCUMULATION_LIFE`, `CEDING_FACINPRODUCTION` dan seluruh master melewati prosedur
   yang badannya tidak ada di korpus. **Tidak mungkin mereplikasi perilaku tulis tanpa membacanya.**
   Yang paling mendesak: `INSERTJSONPOLIS` (siapa menetapkan `PRODKE`?),
   `PROC_GENERATE_SEQUENCE_NUMBER`, `GENERATE_FACRETRO_NO`, `GETCURRENCYSTANDARD` (§8.1).
   *Minta `DBMS_METADATA.GET_DDL` atau akses baca `ALL_SOURCE`.*

4. **[MEMBLOKIR] Pengganti `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK`.** Dua query membaca tabel internal
   Pega untuk **mengambil keputusan routing**. Tabel itu lenyap bersama Pega. Apakah
   `HISTORYAKSEPTASIPEGA` memuat informasi yang setara? Bila tidak, **riwayat harus diekstraksi
   sebelum dekomisioning** — dan itu tidak dapat dibatalkan (§7.1).

5. **[MEMBLOKIR] Determinisme `TGL_TRANSFER`.** `GetAksepBanding_SQL` mengambil "workbasket terakhir"
   dengan `ORDER BY TGL_TRANSFER DESC` + `ROWNUM = 1`, sementara penyisipan memakai `sysdate`. Dua
   baris dalam detik yang sama → jalur banding **tidak deterministik**. Adakah kolom urut lain, atau
   presisi `TGL_TRANSFER` di bawah detik? (§6.4)

6. **[MEMBLOKIR] `DEDUCTION2_MENJADI` / `_SELISIH` pada baris endorsement.** Siklus endorsement tidak
   pernah menulis pasangan kolom ini. Apakah baris endorsement memang selalu `NULL` di sana, dan
   apakah ada agregasi hilir yang mengandaikan sebaliknya? (§2.2)

### Tidak memblokir, tetapi harus dijawab sebelum rekonsiliasi paralel run

7. **Arah tanda `_SELISIH`.** Baru − lama, atau lama − baru? Dan berapa nilainya pada baris New
   Business — nol, atau sama dengan `_MENJADI`? Menentukan apakah `sum(Premi_selisih)` bermakna
   eksposur total atau hanya perubahan (§4.3).

8. **`GETCURRENCYSTANDARD(currency, sysdate)`.** Kurs hari ini, bukan kurs tanggal efektif. Disengaja?
   Bila ya, **paralel run pada hari berbeda akan menghasilkan angka berbeda** dan selisihnya tidak
   akan dapat dijelaskan (§8.2).

9. **Arti kode `STS_KONVERSI_RETRO = '8'`** dan enumerasi lengkap `STS_KONVERSI` (§3.1).

10. **Arti kolom `ID_KOMITE`** pada `HISTORYAKSEPTASIPEGA` (§6.1).

11. **Definisi "versi terakhir" `JSON_POLIS`.** Tiga idiom yang tidak setara hidup berdampingan.
    Adakah baris `JSON_POLIS` yang pernah dihapus? Bila ya, `COUNT-1` sudah salah **sekarang** (§3.2).

12. **`ORDER BY LIMIT_BOTTOM` pada query banding yang menyaring `LIMIT_BOTTOM2`.** Disengaja, atau
    salin-tempel? Memengaruhi urutan anak tangga banding (§5.2).

13. **Master ganda.** Tiga tabel okupasi, tiga mata uang, tiga coverage, empat klausula, dua
    konstruksi. Mana yang otoritatif, dan apakah isinya konsisten? (§9.2)

14. **Pemicu `TREATYPRODUCTION_BACKUP`.** Kapan cadangan dibuat, dan apa kontrak retensinya? (§2.6)

15. **`M_PROMPT_AI`.** Adakah audit trail perubahan prompt yang memengaruhi keputusan underwriting?
    Tinjauan keamanan `CLAUDE.md` §6 harus mencakup tabel ini, bukan hanya rule Pega-nya (§10.2).

16. **`GET_TOKEN_STORAGE`.** Token apa, masa berlaku berapa, dan disimpan di mana? Penanganannya di
    aplikasi baru memerlukan persetujuan manusia (§10.3).

17. **Nullability `WORKBASKET`** pada `HISTORYAKSEPTASIPEGA` — menentukan apakah `!=` ganda setara
    dengan `NOT IN` (§6.3).

18. **Lini bisnis jiwa.** `FACINLIFE`, `FACINSPREADLIFE`, `FACINOFFERLIFE`, `M_LIMIT_LIFE`,
    `M_ACCUMULATION_LIFE` ada di korpus tetapi belum ditabulasi penuh di sini. **Apakah lini jiwa
    masuk lingkup migrasi?** Bila ya, dokumen ini perlu dilengkapi (§2.8).

19. **`nopolis LIKE 'RNM-F%'`.** Konvensi penomoran polis yang hanya hidup di dalam string SQL, tidak
    di kolom mana pun. Apakah stabil? (§2.1)
