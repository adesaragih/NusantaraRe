# Rule yang hanya ada di folder `NB FacIn\` — 139 rule

**Ruang lingkup**: perbandingan **nama berkas** per tipe rule antara `D:\migrasi\RNM\NB FacIn\`,
`D:\migrasi\RNM\RNW Fac In\`, dan `D:\migrasi\RNM\Endorsment Fac In\`.
**Isi yang dibaca dan didokumentasikan hanya dari `NB FacIn\`**; dua folder lain dipakai semata
untuk daftar nama berkas.
**Korpus READ-ONLY** — tidak ada berkas korpus yang diubah.

Label: `[terverifikasi]` · `[dugaan]` · `[pertanyaan terbuka]`.

> ## ⚠️ Dokumen ini ditulis **sebelum** keputusan lingkup. Baca ini dulu.
>
> **Treaty Inward TERMASUK lingkup proyek** — keputusan work owner **K-004**, 15 September 2026
> (`00-KEPUTUSAN-WORK-OWNER.md`). Karena itu:
>
> - **80 rule kelompok G1 (Treaty Inward) TIDAK dikeluarkan dari lingkup**, termasuk seluruh
>   22 activity perhitungan premi Treaty. Setiap kalimat di bawah yang menyebut G1 "bukan perilaku
>   New Business Fac In" tetap **benar sebagai deskripsi** — Treaty memang lini produk berbeda —
>   tetapi **tidak lagi berarti "jangan dimigrasikan"**.
> - Pertanyaan "apakah Treaty Inward masuk lingkup?" di §6.2 dan §7 **sudah terjawab: ya.**
>   Ia bukan lagi pertanyaan terbuka dan bukan blocker.
> - Pertanyaan tentang `When/IsOfferFacIn` di §7 **sudah terjawab K-002**: yang berlaku adalah
>   ekspresi tersimpan `.Quotation.BusinessFac = "F"`.
>
> ⚠️ **Batas cakupan (K-005):** batas kerja tetap `NB FacIn\`. Korpus Treaty Inward tersendiri
> (1.149 berkas di luar `D:\migrasi\RNM\`) **tidak dibaca**. Jadi 80 rule G1 di dokumen ini adalah
> jejak Treaty **di dalam folder Fac In**, bukan gambaran modul Treaty Inward yang utuh.
>
> Isi dokumen di bawah **tidak diubah** selain catatan ini dan tiga penanda di §6.2 dan §7 —
> temuan korpusnya tetap sahih; yang berubah hanya keputusan lingkup di atasnya.

> **Batas metode.** Perbandingan ini berbasis **nama berkas**, sesuai penugasan. Rule yang namanya
> ada di ketiga folder tetap bisa **berbeda isinya** antar folder (brief menyebut 86 rule berbeda
> `pyRuleSetVersion`). Dokumen ini tidak menjawab pertanyaan itu.

---

## 0. Ringkasan temuan

1. **139 rule eksklusif NB**, tersebar di 14 tipe rule. Jumlah `When` eksklusif = **18**,
   **memverifikasi ulang angka pass sebelumnya**. **[terverifikasi]**
2. **Hipotesis "beda waktu ekspor" tidak berlaku di sini.** Ketiga folder diekspor dalam rentang
   satu pekan (commit terakhir: NB 2026-09-01, RNW 2026-09-07, EDM 2026-08-30) dan mencapai patch
   ruleset yang hampir sama (`01-01-96` / `01-01-96` / `01-01-95`). **Seluruh 139 rule eksklusif NB
   di-commit sebelum tanggal ekspor EDM dan RNW** — jadi ketiadaannya di folder lain adalah
   **perbedaan cakupan, bukan perbedaan waktu**. **[terverifikasi]** (perintah audit §2.1)
3. **81 dari 139 (58 %) adalah Treaty Inward, bukan Facultative Inward** — lini produk yang berbeda.
   Termasuk **seluruh 22 aktivitas perhitungan premi** yang berkelas
   `ASM-FW-GISFW-Data-PolicyTreatyIn`. **[terverifikasi]** atas kelas; **[dugaan]** atas kesimpulan
   "lini produk berbeda".
4. **15 rule adalah produk Pega / framework SFA-CRM**, bukan kode Nusantara Re (ruleset
   `PegaCRM-SFA`, `SFAGIS`, `SFAGISFW`, `Pega-UIEngine`, `Pega-API`).
5. **Hanya 28 rule adalah perilaku Fac In New Business sejati** — dan di antaranya terdapat
   **mesin tangga akseptasi NB** (`GetLimitAkseptasi_ActFlow`) yang **tidak ada di RNW maupun EDM**.
   Ini rule yang **memblokir** implementasi NB.
6. **3 rule tanpa pemanggil** di seluruh folder NB: `Activity/InputDtlAdditional_PreAct`,
   `Activity/ValidateCoverage`, `Harness/SFAPortalOpportunities`. Dua yang pertama kandidat kode
   mati; yang ketiga **bukan** (lihat §5.1).
7. Satu rule berada di **ruleset personal seorang pengembang** (`Flow/InputQuotation`,
   ruleset `<nama>@`, v01-01-01) — artinya *checkout* pribadi ikut terekspor. Nama tidak disalin.

---

## 1. Jumlah per tipe rule

### 1.1 Perintah audit

```powershell
$nb="D:\migrasi\RNM\NB FacIn"; $rnw="D:\migrasi\RNM\RNW Fac In"; $edm="D:\migrasi\RNM\Endorsment Fac In"
foreach ($t in (Get-ChildItem $nb -Directory).Name | Sort-Object) {
  $a = @(Get-ChildItem "$nb\$t"  -Filter *.xml -File -EA SilentlyContinue).Name
  $b = @(Get-ChildItem "$rnw\$t" -Filter *.xml -File -EA SilentlyContinue).Name
  $c = @(Get-ChildItem "$edm\$t" -Filter *.xml -File -EA SilentlyContinue).Name
  $other = New-Object 'System.Collections.Generic.HashSet[string]' ([string[]]($b+$c), [StringComparer]::OrdinalIgnoreCase)
  $only = @($a | Where-Object { -not $other.Contains($_) })
  "{0,-18} NB={1,4}  eksklusif={2,3}" -f $t, $a.Count, $only.Count
}
```

### 1.2 Hasil

| Tipe rule | NB | RNW | EDM | **Eksklusif NB** |
| --- | ---: | ---: | ---: | ---: |
| `Activity` | 609 | 556 | 582 | **50** |
| `Section` | 432 | 397 | 434 | **27** |
| `When` | 210 | 189 | 202 | **18** |
| `RDBList` | 217 | 200 | 188 | **15** |
| `DataTransform` | 148 | 138 | 157 | **8** |
| `FlowAction` | 250 | 239 | 265 | **8** |
| `ReportDefinition` | 123 | 118 | 138 | **5** |
| `Flow` | 6 | 4 | 2 | **3** |
| `DecisionTable` | 12 | 10 | 10 | **2** |
| `DecisionTree` | 1 | 1 | 1 | **1** |
| `Harness` | 41 | 41 | 37 | **1** |
| `SystemSettings` | 1 | 1 | 1 | **1** |
| `ConnectREST` | 3 | 3 | 3 | 0 |
| `DataPage` | 30 | 30 | 41 | 0 |
| **Total** | **2.083** | **1.927** | **2.061** | **139** |

Catatan: `DecisionTree` dan `SystemSettings` masing-masing hanya berisi **satu** rule di tiap
folder, dan namanya **berbeda** di NB — karena itu terhitung eksklusif meski jumlahnya sama.

**Verifikasi angka 18 `When`**: pass sebelumnya melaporkan 18 rule `When` eksklusif NB. Perhitungan
ulang di atas menghasilkan **18** — **cocok**. **[terverifikasi]**

---

## 2. "Eksklusif NB" atau "artefak ekspor"?

Ini pertanyaan yang menentukan apakah 139 rule ini wajib diimplementasikan.

### 2.1 Bukti: ketiga ekspor hampir sezaman

```powershell
foreach ($k in 'NB FacIn','RNW Fac In','Endorsment Fac In') {
  $vers=@(); $commits=@()
  foreach ($f in (Get-ChildItem "D:\migrasi\RNM\$k" -Recurse -Filter *.xml -File)) {
    $t = Get-Content $f.FullName -Raw
    if ($t -match '<pyRuleSetVersion>([^<]+)') { $vers += $Matches[1] }
    if ($t -match '<pxCommitDateTime>([^<]+)') { $commits += $Matches[1] } }
  $gis = @($vers | Where-Object { $_ -like '01-01-*' } | ForEach-Object { [int]($_ -split '-')[2] })
  "{0}: patch GISFW tertinggi={1}  commit min={2} max={3}" -f $k,
     ($gis|Measure-Object -Maximum).Maximum, ($commits|Sort-Object)[0], ($commits|Sort-Object)[-1] }
```

| Folder | Patch `01-01-xx` tertinggi | `pxCommitDateTime` terawal | `pxCommitDateTime` terakhir |
| --- | ---: | --- | --- |
| `NB FacIn` | **96** | 2025-05-13T20:53:19 | **2026-09-01T07:03:04** |
| `RNW Fac In` | **96** | 2025-05-13T20:54:00 | **2026-09-07T09:05:08** |
| `Endorsment Fac In` | **95** | 2025-05-13T15:08:50 | **2026-08-30T14:02:15** |

Ketiganya bertitik-awal sama (13 Mei 2025) dan berakhir dalam rentang **8 hari** (30 Agu – 7 Sep
2026). **[terverifikasi]**

### 2.2 Bukti: seluruh 139 rule eksklusif sudah ada sebelum ekspor terakhir folder lain

```powershell
# excl_full.csv memuat kolom Name, Type, RuleSet, Ver, Commit untuk 139 rule eksklusif NB
$r = Import-Csv excl_full.csv
"commit SETELAH ekspor EDM (>20260830T140215): " + @($r | Where-Object { $_.Commit -gt '20260830T140215' }).Count
"commit SEBELUM/ sama dengan             : " + @($r | Where-Object { $_.Commit -le '20260830T140215' }).Count
# -> 0  dan  139
```

**Hasil: 0 dan 139.** Tidak satu pun rule eksklusif NB lebih baru daripada ekspor EDM.
Karena itu **hipotesis "rule baru yang belum masuk ekspor folder lain" gugur untuk seluruh 139
rule**. **[terverifikasi]**

### 2.3 Bukti positif: tiap siklus punya Flow utamanya sendiri

```powershell
foreach ($d in 'NB FacIn','RNW Fac In','Endorsment Fac In') {
  "{0,-20}: {1}" -f $d, (((Get-ChildItem "D:\migrasi\RNM\$d\Flow").BaseName | Sort-Object) -join ', ') }
```

| Folder | Isi folder `Flow\` |
| --- | --- |
| `NB FacIn` | `InputInwardFacultativeOffer`, `InputInwardFacultativeRISlip`, **`InputQuotation`**, **`InputRealizationTreatyIn`**, `OfferFacOut`, `OfferFacRetro` |
| `RNW Fac In` | `InputInwardFacultativeRISlip`, `InputRenewalFacultativeIn`, `OfferFacOut`, `OfferFacRetro` |
| `Endorsment Fac In` | `InputAddendumFacultativeIn`, `OfferFacRetro` |

Tiga flow entri berbeda untuk tiga siklus — `InputInwardFacultativeOffer` (NB),
`InputRenewalFacultativeIn` (RNW), `InputAddendumFacultativeIn` (EDM). Ini **perbedaan cakupan yang
disengaja**, bukan artefak. **[terverifikasi]**

### 2.4 Kesimpulan

| Pernyataan | Status |
| --- | --- |
| "Hanya ada di NB" berarti "rule baru yang belum masuk ekspor folder lain" | **gugur** untuk 139/139 rule **[terverifikasi]** |
| "Hanya ada di NB" berarti "cakupan folder NB memang lebih luas" | **didukung** — NB memuat Treaty Inward + portal SFA yang tidak diekspor ke folder lain **[terverifikasi]** |
| Setiap rule eksklusif NB otomatis = perilaku khas New Business | **salah** — hanya 28 dari 139 **[terverifikasi]** |

---

## 3. Pengelompokan menurut fungsi

```powershell
# klasifikasi: ruleset (PegaCRM/SFAGIS/Pega-*) -> G2 ; nama atau pyClassName bermuatan
# Treaty/PolicyTreatyIn/Data-Portal/POLISTREATYIN -> G1 ; Warranty|Vehicle|Cargo|Coverage -> G4 ;
# util (GetSQLDate|GetCurrentDate|SearchJobID|LinkService|pz*|py*) -> G5 ; sisanya -> G3
```

| Grup | Rule | Perilaku khas NB? |
| --- | ---: | --- |
| **G1 — Treaty Inward** | **80** | **Tidak.** Lini produk berbeda (Treaty), kebetulan diekspor bersama NB. |
| **G2 — SFA/CRM & produk Pega** | **15** | **Tidak.** Rule produk pihak ketiga; hilang bersama Pega. |
| **G3 — Facultative Inward New Business** | **28** | **Ya.** Wajib ada di sistem baru. |
| **G4 — Lini langsung / cover ritel** | **11** | **Tidak** (dugaan). Warisan framework GIS. |
| **G5 — Infrastruktur & util** | **5** | Sebagian; lihat §3.5. |
| **Total** | **139** | |

> **Koreksi atas klasifikasi otomatis.** Skrip menempatkan `Section/FormulaTreatyCapacityDesc` di G1
> karena namanya memuat "Treaty". Pemeriksaan isi menunjukkan kelasnya
> `ASM-FW-GISFW-Data-OfferFacIn` dan pemanggilnya `Section/InwardFacIn`,
> `Section/OfferFacIn_NusaRe`, `Section/FormulaTreatyCapacityDesc_IsUW` — jadi ia **milik jalur
> Fac In** (menampilkan kapasitas treaty di dalam layar offer Fac In). Angka di tabel sudah
> dikoreksi: G1 = 80, G3 = 28. Contoh telanjang dari aturan **nama bukan bukti**.

### 3.1 G1 — Treaty Inward (80 rule) — **bukan perilaku New Business Fac In**

Bukti terkuat berupa **kelas**, bukan nama:

```powershell
# kelas dari 50 Activity eksklusif NB
Import-Csv excl_grp.csv | Where-Object {$_.Type -eq 'Activity'} | Group-Object Cls |
  Sort-Object Count -Descending | ForEach-Object { "{0,-44} {1}" -f $_.Name,$_.Count }
```

| Kelas | Activity eksklusif |
| --- | ---: |
| `ASM-FW-GISFW-Data-PolicyTreatyIn` | **22** |
| `ASM-FW-GISFW-Work` | 14 |
| `Data-Portal` | 7 |
| `ASM-FW-GISFW-Data-Coverage` | 3 |
| `ASM-FW-GISFW-Data-TreatyInLimits` | 1 |
| `ASM-FW-GISFW-Data-Vehicle` | 1 |
| `ASM-FW-GISFW-Work-NB` | 1 |
| `Pega-API-CaseManagement-Case` | 1 |

**Seluruh mesin perhitungan premi yang "eksklusif NB" sebenarnya milik Treaty Inward**
(kelas `Data-PolicyTreatyIn`) — 22 aktivitas: `CalculatePremi_Act`, `CountNetPremi_act`,
`CountOGPONP_Act`, `CountOverridingCommOgp_Act`, `CountOverridingCommOnp_Act`,
`CountPctInstallment_Act`, `CountResult1_Act`, `CountResult1Onp_Act`, `CountResult2Ogp_act`,
`CountResult2Onp_act`, `CountRiCommOgp_act`, `CountRiCommOnp_act`, `CountSpreading_Act`,
`SetPPNPPH`, `RemoveTypeTax_ACT`, `SetCurrency_act`, `SetDueTo_act`, `ProtectDate`,
`CheckDataMkt`, `CekLimitTreatyAcc_Act`, `GeneratePolicyNoTreaty_Act`, `SaveJsonPolisTreatyIn_Act`.
**[terverifikasi]** dari tag `<pyClassName>` masing-masing berkas.

Isi yang terbaca pada wakil-wakilnya:

| Rule | Isi (dikutip) | Pemanggil |
| --- | --- | --- |
| `RDBList/SavePolisTreatyIn_SQL` | blok PL/SQL memanggil `POOLDATA.PEGA_JSON_POLIS_TREATYIN(...)` lalu `COMMIT` | `Activity/SaveJsonPolisTreatyIn_Act` |
| `RDBList/SaveTreatyIn` | `POOLDATA.PEGA_TREATY_IN(...)` dengan 21 argumen masuk + 3 keluar (`ERRMSG`, `IDPEGAOUT`, `S…`) | `Activity/SetTreatyIn_Act` |
| `RDBList/BrowseTreatyIn` | `select JSONDATA … from pooldata.M_TREATY_IN where ID={TreatyIn.ID}` **UNION ALL** idem dari `pooldata.M_TREATY_IN_edm` | `Activity/SetTreatyIn_Act` |
| `RDBList/FetchTreatyGroupOLDID` | `select oldid as CARI1, treatygroupname as CARI2 from pooldata.treatygroup where ID = {InputData.CARI1}` | `Activity/GeneratePolicyNoTreaty_Act` |
| `RDBList/GetBreakDownSpread_SQL` | `select DISTINCT Reinstypeid AS "TreatyType", PCT AS "SharePercentage" from pooldata.proportionalarrg a where a.treatygroupid={ParamData.CARI10} and a.PARENTREINSTYPEID={ParamData.CARI11} and TreatyDescID='10001'` | `Activity/BreakDownSpreading_Act` |
| `RDBList/TreatyRealizationCheckDuplicate` | `select nopolis as CARI1 from pooldata.treatyinproduction where NOOFFER=… and BEGINDATE=TO_DATE(…, 'YYYYMMDD') and … and BALANCE_DUE_TO = REPLACE({InputData.Totaltsi},',','.')` | `Activity/InputPolicyTreatyInPost_Act` |
| `RDBList/GenerateNoPolicy` | `SELECT 'RNM-' \|\| {InputData.CARI20} \|\| '.T' \|\| …BusinessOldId \|\| '.' \|\| to_char(sysdate,'MM.yyyy') \|\| '.' \|\| LPAD(POOLDATA.JSON_POLIS_TREATYIN_SEQ.NEXTVAL, 5, '0')` | `Activity/GenerateNopolis_Act`, `SaveJsonPolicyFacIn_Act`, `SaveJsonPolisTreatyIn_Act` |
| `RDBList/GetCountClaim` | `select count(1) as "CARI1" from DATAPEGA.PC_ASM_FW_GCNMFW_WORK where masterid={TreatyWarning.CAIREINSFACIN}` | `Activity/CheckDuplicateOffer` |

Tiga temuan dari kutipan di atas:

- **`TreatyRealizationCheckDuplicate` membandingkan jumlah uang sebagai string.**
  `BALANCE_DUE_TO = REPLACE({InputData.Totaltsi}, ',', '.')` — nilai TSI dikirim sebagai teks
  berkoma desimal, komanya diganti titik, lalu diadu dengan kolom. Tepat pola yang aturan proyek
  §4.1 peringatkan. **[terverifikasi]**
- **`GetCountClaim` membaca tabel internal Pega `DATAPEGA.PC_ASM_FW_GCNMFW_WORK`.** Menurut aturan
  proyek §4.3 tabel `PC_*` **hilang bersama Pega** dan tidak boleh dimigrasikan apa adanya.
  **[terverifikasi]**
- **`GenerateNoPolicy` ambigu.** Berkasnya berkelas `ASM-FW-GISFW-Int-POLISTREATYIN` dan memakai
  sekuens `JSON_POLIS_TREATYIN_SEQ` serta menyisipkan literal `'.T'`, namun ia **juga dipanggil dari
  `Activity/SaveJsonPolicyFacIn_Act`** dengan `BrowsePage = PolicyFacIn`. Lihat §6 butir MEMBLOKIR.

### 3.2 G2 — SFA/CRM & produk Pega (15 rule) — **bukan kode Nusantara Re**

| Ruleset | Rule |
| --- | --- |
| `PegaCRM-SFA` v07-22-01 | `When/crmCreateOpportunity`, `When/isSellingModeB2B`, `When/isSellingModeB2BB2C`, `When/isSellingModeB2C`, `When/pyIsIpadOrDesktop`, `Harness/SFAPortalOpportunities` |
| `SFAGISFW` | `Section/SFAPortal_Opportunities` (v01-01-39), `Section/SFAPortal_OpportunitiesList_Header` (v01-01-04), `Section/SFAPortalOpportunitiesHeader` (v01-01-39) |
| `SFAGIS` | `Section/SFAPortal_OpportunitiesList` (v01-01-32) |
| `Pega-UIEngine` | `When/pyIsMobile` (v08-02-01) |
| `Pega-API` | `Activity/pzChangeStageWrapper` (v08-05-01) |
| `GISFW` (tetapi melayani portal SFA) | `ReportDefinition/GetListOpportunity`, `…F`, `…Life` (v01-01-90) |

Isi kondisi `When` yang terbaca:

```
isSellingModeB2B      : @(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") EQUALS "B2B"
isSellingModeB2BB2C   : … equals "B2B_B2C"
isSellingModeB2C      : … equals "B2C"
pyIsMobile            : @(Pega-RulesEngine:Utilities).pzIsMobile(tools)
pyIsIpadOrDesktop     : Rule pyIsIPad evaluates to true  ||  pxRequestor.pxDeviceType = desktop
crmCreateOpportunity  : true = evaluateWhen("crmBypassOperatorAccessChecks")
                        || "true" equals Declare_crmOperatorAccess.canCreateOpportunity
                        || true = evaluateWhen("crmIsOpen")
```
(sumber: `NB FacIn/When/<nama>.xml`, tag `pyConditionValue1String`). **[terverifikasi]**

Ketiganya bergantung pada `Data-Admin-System-Settings` bernama `PegaCRM-`/`SellingMode`, yang
**tidak ada di korpus** → nilai aktif **belum terverifikasi**.

**Konsekuensi migrasi**: seluruh 15 rule ini **tidak dimigrasikan**. Portal opportunity harus
didesain ulang; `pyIsMobile`/`pyIsIpadOrDesktop` digantikan CSS responsif React.

### 3.3 G3 — Facultative Inward New Business (28 rule) — **wajib ada di sistem baru**

| Tipe | Rule | Kelas | Ver | Pemanggil |
| --- | --- | --- | --- | --- |
| `Flow` | **`InputInwardFacultativeOffer`** | `ASM-FW-GISFW-Work` | 01-01-94 | `Flow/InputQuotation` |
| `Flow` | `InputQuotation` | `ASM-FW-GISFW-Work-NB` | 01-01-01 (ruleset personal) | `Section/CedingCoHierarki` |
| `Activity` | **`GetLimitAkseptasi_ActFlow`** | `ASM-FW-GISFW-Work` | 01-01-87 | `Activity/GetLimitAkseptasi_JUW_UW`, `Flow/InputInwardFacultativeOffer` |
| `Activity` | **`GetLimitAkseptasiLife_Act`** | `ASM-FW-GISFW-Work` | 01-01-81 | `Flow/InputInwardFacultativeOffer` |
| `Activity` | `SaveJsonOfferFacIn_Act` | `ASM-FW-GISFW-Work` | 01-01-81 | 4 pemanggil |
| `Activity` | `SetToJsonOffer_ACT` | `ASM-FW-GISFW-Work` | 01-01-79 | `Activity/…`, `Flow/…` |
| `Activity` | `SetBusinessType_Act` | `ASM-FW-GISFW-Work` | 01-01-81 | `Flow/InputInwardFacultativeOffer` |
| `Activity` | `BreakDownSpreading_Act` | `ASM-FW-GISFW-Work` | 01-01-91 | 1 Activity |
| `Activity` | `CalcultePersentageSpeadingLife_Act` | `ASM-FW-GISFW-Work` | 01-01-53 | 2 Section |
| `Activity` | `ShowViewCheckListFacOut` | `ASM-FW-GISFW-Work` | 01-01-52 | 1 Harness + 1 Section |
| `DataTransform` | `InputQuotation_PreDT` | `ASM-FW-GISFW-Work` | 01-01-91 | `Flow/InputQuotation` |
| `DecisionTable` | `BusinessType_DeT` | `ASM-FW-GISFW-Work` | 01-01-87 | `Activity/SetBusinessType_Act` |
| `DecisionTree` | `Tree_ShortPeriod` | `ASM-FW-GISFW-Data-Coverage` | 01-01-52 | `Activity/CountProrateExtension_Act` |
| `RDBList` | **`GetLimitAkseptasiLife_SQL`** | `Int-policyjson` | 01-01-81 | `Activity/GetLimitAkseptasiLife_Act` |
| `RDBList` | `SaveOfferJson_SQL` | `Int-OFFERJSON` | 01-01-90 | `Activity/SaveJsonOfferFacIn_Act` |
| `RDBList` | `GetInsuredID` | `Int-CLIENT` | 01-01-56 | 3 Activity |
| `RDBList` | `GetBreakDownSpread_SQL` | `ASM-FW-GISFW-Work` | 01-01-95 | `Activity/BreakDownSpreading_Act` |
| `RDBList` | `GetCountClaim` | `Assign-Worklist` | 01-01-56 | `Activity/CheckDuplicateOffer` |
| `FlowAction` | `InstallmentList` | `Data-Installment` | 01-01-56 | 11 pemanggil |
| `FlowAction` | `ReasViewAttachment` | `ASM-FW-GISFW-Int` | 01-01-52 | 3 pemanggil |
| `Section` | `InwardFacIn` | `Data-OfferFacIn` | 01-01-53 | 10 pemanggil |
| `Section` | `InstallmentList` | `Data-Installment` | 01-01-80 | 11 pemanggil |
| `Section` | `InputInwardFacultativeSuggest` | `Data-OfferFacIn` | 01-01-53 | `Section/InputEndorsement`, `Section/InputInwardFacultative_2` |
| `Section` | `FormulaTreatyCapacityDesc` | `Data-OfferFacIn` | 01-01-53 | `Section/InwardFacIn`, `Section/OfferFacIn_NusaRe`, `…_IsUW` |
| `When` | **`ToDeptHeadUWLife`** | `ASM-FW-GISFW-Work` | 01-01-81 | `Flow/InputInwardFacultativeOffer` |
| `When` | `IsOfferFacIn` | `ASM-FW-GISFW-Work` | 01-01-55 | `Flow/InputQuotation` |
| `When` | `IsSPVCreate` | `ASM-FW-GISFW-Work` | 01-01-81 | `Flow/InputRealizationTreatyIn` + disebut di 2 rule `When` lain |
| `When` | `IsNotAdmin` | `@baseclass` | 01-01-55 | `Harness/SFAPortalOpportunities`, `Section/SFAPortalOpportunitiesHeader` |

#### 3.3.1 `GetLimitAkseptasi_ActFlow` — mesin tangga akseptasi NB

Ini rule tunggal terpenting dalam daftar. **Tidak ada di RNW maupun EDM**:

```powershell
foreach ($d in 'NB FacIn','RNW Fac In','Endorsment Fac In') {
  "{0,-20}: {1}" -f $d, (((Get-ChildItem "D:\migrasi\RNM\$d\Activity" -Filter 'GetLimitAkseptasi*').BaseName) -join ', ') }
# NB FacIn          : GetLimitAkseptasiLife_Act, GetLimitAkseptasi_Act, GetLimitAkseptasi_ActFlow, GetLimitAkseptasi_JUW_UW
# RNW Fac In        : GetLimitAkseptasi_Act, GetLimitAkseptasi_JUW_UW
# Endorsment Fac In : GetLimitAkseptasi_Act, GetLimitAkseptasi_JUW_UW
```

`GetLimitAkseptasi_Act` dan `GetLimitAkseptasi_JUW_UW` ada di ketiga folder; **varian `_ActFlow` dan
`Life_Act` hanya di NB**. **[terverifikasi]**

Yang dirujuknya — 12 `RDBList` tangga limit, **semuanya ada juga di RNW/EDM** (jadi bukan eksklusif):
`GetLimitAkseptasi_SQL`, `GetLimitAkseptasiBanding_SQL`, `GetLimitAkseptasiBond_SQL`,
`GetLimitAkseptasiKreditCL_SQL`, `GetLimitAkseptasiKreditNCL_SQL`, `GetLimitAkseptasiNonFire_SQL`,
`GetLimitAkseptasiNonFireBanding_SQL`, `GetLimitAkseptasiNonPrefer_SQL`,
`GetLimitAkseptasiNonPreferBanding_SQL`, `GetLimitAkseptasiPreferedComm_SQL`,
`GetLimitAccEngineeringUW_SQL`, `GetLimitAccEngineeringBanding_SQL`. **[terverifikasi]**

**Artinya: SQL tangga limitnya dipakai bersama tiga siklus; orkestrasinya khas NB.**

Percabangan yang terbaca pada `pyStepsPreCondParamsWhen` (kutipan terpilih):

```
IsGroup
pyWorkPage.IsAdaTopRisk=="true" || pyWorkPage.OfferFacIn.TotalTSITopRisk>0
pyWorkPage.OfferFacIn.QuotationData.StatusBusiness=="3"
pyWorkPage.OfferFacIn.IsPreferredRisk=="Preferred Risk"
pyWorkPage.OfferFacIn.IsPreferredRisk=="Non-Preferred Risk"
pyWorkPage.OfferFacIn.IsPreferredRisk=="Preferred Risk Commercial"
pyWorkPage.OfferFacIn.IsBanding=="true" && pyWorkPage.OfferFacIn.IsFlagReject=="true"
pyWorkPage.OfferFacIn.IsOccupException=="1" && pyWorkPage.OfferFacIn.IsSpecialAcceptance!="true"
pyWorkPage.IsB2B=="ASM"
pyWorkPage.PositionNote=="ReasFacInUnderwriting" | "…SeniorUnderwriting" | "…DepHeadUnderwriting"
   | "…ManagerTeknik" | "…GroupLeader" | "…FacultativeDivHead" | "…MarketingDirector"
   | "…UnderwritingFinancial" | "…FinDivHead"
.CARI1=="SENIOR UNDERWRITER" | "DEP.HEAD UNDERWRITER" | "MANAGER TEKNIK" | "MANAGERTEKNIK"
   | "KADIV FACULTATIVE" | "KADIVFACULTATIVE" | "KADIV TEKNIK" | "KADIV KEUANGAN"
   | "DIREKTUR MARKETING" | "DIREKTUR TEKNIK"
```

**[terverifikasi]** semuanya dari `NB FacIn/Activity/GetLimitAkseptasi_ActFlow.xml`.

Empat catatan penting:

1. **Tautologi 15 kali.** Ekspresi `Local.TotalTSI>=@toDecimal(Local.Limit) || Local.TotalTSI<=@toDecimal(Local.Limit)`
   muncul **15 kali**; `>=` sendirian muncul **18 kali**. Bentuk `x>=y || x<=y` **selalu benar**
   untuk nilai yang dapat dibandingkan. **Kandidat perbaikan, bukan bagian migrasi** — implementasi
   baru harus meniru perilakunya agar rekonsiliasi paralel run cocok.
   ```powershell
   $t = Get-Content "D:\migrasi\RNM\NB FacIn\Activity\GetLimitAkseptasi_ActFlow.xml" -Raw
   ([regex]::Matches($t,'Local\.TotalTSI&gt;=@toDecimal\(Local\.Limit\) \|\| Local\.TotalTSI&lt;=@toDecimal\(Local\.Limit\)')).Count  # -> 15
   ([regex]::Matches($t,'Local\.TotalTSI&gt;=@toDecimal\(Local\.Limit\)')).Count                                                     # -> 18
   ```
2. **Nomor polis literal sebagai gerbang.** Aktivitas ini memuat **4 nomor polis unik** (8 kemunculan)
   yang dibandingkan langsung dengan `pyWorkPage.OfferFacIn.QuotationData.OldPolicyNo` untuk
   membelokkan alur. Nilainya **tidak disalin** ke dokumen ini sesuai aturan §7; berkasnya
   `NB FacIn/Activity/GetLimitAkseptasi_ActFlow.xml`.
   ```powershell
   ([regex]::Matches($t,'RNM-F[0-9.]+') | ForEach-Object {$_.Value} | Select-Object -Unique).Count  # -> 4
   ```
3. **Dua konvensi penamaan jabatan hidup berdampingan** pada `.CARI1` — dengan spasi
   (`"MANAGER TEKNIK"`, `"KADIV FACULTATIVE"`) dan tanpa spasi (`"MANAGERTEKNIK"`,
   `"KADIVFACULTATIVE"`). Nilai ini datang dari kolom `JABATAN` hasil `GetLimitAkseptasi*_SQL`.
   **[pertanyaan terbuka]** apakah keduanya benar-benar ada di data, atau salah satu cabang mati.
4. **`IsGroup`** dipakai sebagai kondisi pertama. Catatan proyek menyebut `IsGroup` bercabang atas
   3 identitas operator + 1 kode kontak; dokumen ini tidak menyalin nilainya.

#### 3.3.2 `GetLimitAkseptasiLife_Act` + `GetLimitAkseptasiLife_SQL` — tangga akseptasi Life

SQL-nya utuh dan pendek:

```sql
-- NB FacIn/RDBList/GetLimitAkseptasiLife_SQL.xml, kelas ASM-FW-GISFW-Int-policyjson
SELECT JABATAN AS CARI1, NAMA AS CARI5
FROM M_LIMIT_LIFE
WHERE LIMIT_BOTTOM <= {DataSearch.CARID2}
ORDER BY ID ASC
```

**[terverifikasi]**. Tiga catatan:

- Tabel **`M_LIMIT_LIFE`** terpisah dari tabel tangga Fac In biasa. Skemanya tidak disebut (tanpa
  awalan `POOLDATA.`) — **belum terverifikasi**.
- **Hanya batas bawah** (`LIMIT_BOTTOM <=`). Tidak ada batas atas, sehingga query mengembalikan
  **semua** jabatan yang batas bawahnya di bawah nilai yang dicari; pemilihan barisnya ditentukan
  `ORDER BY ID ASC` + pengambilan baris pertama oleh pemanggil. **[dugaan]** untuk mekanisme
  pengambilan baris — belum dibuktikan dari langkah aktivitas.
- Parameter bernama **`CARID2`** (bukan `CARI2`) — penamaan menyimpang dari konvensi `CARIn` yang
  dipakai di seluruh korpus. Nama yang sama muncul di `GetLimitAkseptasi_ActFlow`, jadi konsisten
  antar kedua rule.

Percabangan `GetLimitAkseptasiLife_Act`:

```
IsFire
pyWorkPage.PositionNote=="ReasFacInUnderwriting"
IsT1T3
pyWorkPage.PositionNote=="ReasFacInUnderwritingLife"     -> .CARI1=="DEP.HEAD UNDERWRITER"
pyWorkPage.PositionNote=="ReasFacInDepHeadUnderwritingLife" -> .CARI1=="DIREKTUR TEKNIK"
```
**[terverifikasi]**. Tangga Life hanya punya **dua** anak tangga yang terbaca, jauh lebih pendek
daripada tangga non-Life.

#### 3.3.3 Persistensi offer NB

```sql
-- NB FacIn/RDBList/SaveOfferJson_SQL.xml, kelas ASM-FW-GISFW-Int-OFFERJSON
BEGIN
  POOLDATA.PEGA_M_JSON_OFFER(
    {InputParam.DATAPEGA},
    {pyWorkPage.pzInsKey},
    {pyWorkPage.OfferFacIn.QuotationData.BusinessOldId},
    {InputParam.FINDDATA1},
    {OutputParam.ERRMSG OUT},
    {OutputParam.IDPEGAOUT OUT},
    {OutputParam.STSSAVE OUT} );
  COMMIT;
END;
```

**[terverifikasi]**. Seluruh agregat `OfferFacIn` diserialkan ke satu argumen `DATAPEGA` dan
disimpan oleh stored procedure `POOLDATA.PEGA_M_JSON_OFFER`. **Isi prosedur tidak ada di korpus →
belum terverifikasi.** Pemanggilnya `Activity/SaveJsonOfferFacIn_Act`, yang juga memakai
`RDBList/GetInsuredID` (`select ID as CARI1 from pooldata.client where name = {InputData.CARI1}` —
**pencarian tertanggung berdasarkan nama, bukan ID**; berpotensi ambigu) dan
`RDBList/GetOldIDBusiness_SQL`.

Kondisi yang menggerakkan `SaveJsonOfferFacIn_Act` (tag `pyStepsPreCondParamsWhen`):

```
@LengthOfPageList(IDPegaOffer.pxResults)>=1
pyWorkPage.OfferFacIn.QuotationData.InsuredID==""
pyWorkPage.OfferFacIn.QuotationData.MarketingCode==""
@Utilities.PropertyHasValue(pyWorkPage.OfferFacIn.ID)
IsLife
```
**[terverifikasi]**

#### 3.3.4 Rule `When` khas NB

| Rule | Kondisi (`pyConditionValue1String`) | Catatan |
| --- | --- | --- |
| `IsOfferFacIn` | `.Quotation.BusinessFac = "F"` | **`pyConditionString` berbunyi lain**: `Kode Bisnis = "02" \|\| "58" \|\| "SB" \|\| "SG"`. Lihat §4. |
| `ToDeptHeadUWLife` | `pyWorkPage.LetterNo = "DEPTHEADUWLIFE"` | `LetterNo` dipakai sebagai **token routing**, bukan nomor surat. |
| `IsSPVCreate` | bercabang atas **2 identitas operator** (`pyWorkPage.pxCreateOperator`) | Nilai tidak disalin. Guard berbasis identitas. |
| `IsNotAdmin` | `OperatorID.pyPosition != "Admin"` | kelas `@baseclass`; dipakai portal SFA. |

### 3.4 G4 — Lini langsung / cover ritel (11 rule) — **dugaan: warisan framework GIS**

`Activity/InputDtlAdditional_PreAct` (0 pemanggil), `Activity/ValidateCoverage` (0 pemanggil),
`Activity/SpreadingAdditionalProtection`, `Section/VehicleGrid`, `Section/InputDtlCargo`,
`Section/ViewDtlAdditional`, `Section/WarrantyList`, `Section/SelectWarrantyList`,
`Section/CreateListWarranty`, `FlowAction/InputDtlWarranty`, `ReportDefinition/BrowseWarranty_RD`.

Bukti bahwa ini bukan Fac Inward — kelasnya berorientasi objek yang ditanggung langsung:
`ASM-FW-GISFW-Data-Vehicle`, `ASM-FW-GISFW-Data-Cargo`, `ASM-FW-GISFW-Data-Coverage`.
`Activity/ValidateCoverage` beroperasi atas `pyWorkPage.VehicleList(i).CoverageList(j).AdditionalCoverage(k)`
**[terverifikasi]**.

**Tetapi 9 dari 11 tetap dirujuk dari kode NB** (mis. `Section/WarrantyList` punya 8 pemanggil,
`Section/VehicleGrid` 4, `Section/InputDtlCargo` 4 — di antaranya `FlowAction/VehicleGrid_FacIn`
dan `FlowAction/InputDtlCargo_FacIn`, yang namanya justru berakhiran `_FacIn`).
**[pertanyaan terbuka]**: apakah layar kendaraan/kargo benar-benar dipakai pada offer Fac In, atau
hanya tersambung tanpa pernah dibuka.

### 3.5 G5 — Infrastruktur & util (5 rule)

| Rule | Isi | Pemanggil | Nasib |
| --- | --- | --- | --- |
| `RDBList/GetSQLDate` | `select sysdate as CARI1 from dual` | `GeneratePolicyNoTreaty_Act`, `InputPolicyTreatyInPre_Act` | ganti dengan waktu server aplikasi |
| `RDBList/SearchJobID` / `SearchJobIDSQL` | `select KELAS_ID, GOLONGAN from GENERAL.V_JOB_PA where DESCRIBE={SearchJobIDInput.DESCRIBE}` (identik; `SearchJobID` memanggil `SearchJobIDSQL`) | `Activity/UploadCSVPerson_PostAct` | duplikasi; ruleset `GISFWInt` v01-01-34 |
| `Activity/CopyTemplateCov_PostAct` | kelas `Data-Vehicle` | 1 Activity | ikut G4 secara substansi |
| **`SystemSettings/LinkService`** | `pySetting` = `=ResponLink.URL` ; `pyPurpose` = `LinkService` ; v01-01-88 | `Activity/GetLinkService`, `ConnectREST/convertJsonNusareToProduction`, `ConnectREST/getPremiumPaidOn`, **`ConnectREST/ServiceGoogle`** | lihat di bawah |

**`SystemSettings/LinkService` penting untuk aturan proyek §4.4.** Nilainya **bukan URL literal**
melainkan ekspresi `=ResponLink.URL` — rujukan ke halaman clipboard `ResponLink` yang diisi dari
hasil pencarian tabel `M_LINK_SERVICE`. **Tidak ada satu pun endpoint literal di berkas ini.**
**[terverifikasi]**. Ia adalah satu-satunya `SystemSettings` di folder NB dan namanya berbeda dari
padanan di RNW/EDM — **[pertanyaan terbuka]** nama setting di dua folder lain.

Salah satu pemanggilnya adalah `ConnectREST/ServiceGoogle`. Aturan proyek §6 mensyaratkan
persetujuan manusia sebelum memindahkan integrasi model AI pihak ketiga di alur underwriting;
`LinkService` adalah jalur konfigurasi endpoint-nya. **Isi `ServiceGoogle` tidak dibaca dalam pass
ini.**

---

## 4. `pyConditionString` vs `pyConditionValue1String` — keduanya wajib dibaca, dan bisa bertentangan

Brief menyarankan memeriksa kedua tag. Pada 18 `When` eksklusif NB, keduanya **identik pada 8 rule**,
`pyConditionString` berisi placeholder `[Double click to add condition]` pada **8 rule**, dan
**bertentangan pada 2 rule** (8 + 8 + 2 = 18):

| Rule | `pyConditionString` | `pyConditionValue1String` |
| --- | --- | --- |
| `IsOfferFacIn` | `Kode Bisnis = "02" \|\| "58" \|\| "SB" \|\| "SG"` | `.Quotation.BusinessFac = "F"` |
| `IsOperatorLife` | `OperatorID.pyAccessGroup = "GISFW:AuctionUsers"` | `OperatorID.pyWorkGroup = "ReasLife"` |

**[terverifikasi]** dari `NB FacIn/When/IsOfferFacIn.xml` dan `…/IsOperatorLife.xml`.

`pyConditionString` di Pega adalah **teks tampilan** di rule form; `pyConditionValue1String` adalah
ekspresi yang tersimpan. **[dugaan]** bahwa `pyConditionValue1String` yang dieksekusi — ini konvensi
Pega, **tidak dibuktikan oleh korpus**.

**Konsekuensi**: teks `Kode Bisnis = "02"/"58"/"SB"/"SG"` **tidak boleh** dijadikan spesifikasi
`IsOfferFacIn`. Empat kode bisnis itu kemungkinan sisa versi lama. **[pertanyaan terbuka]** — dan
ini **MEMBLOKIR**, karena `IsOfferFacIn` menggerbangi `Flow/InputQuotation`.

### 4.1 Kondisi lengkap 18 `When` eksklusif NB

| Rule | Ver | Kondisi tersimpan | Pemanggil |
| --- | --- | --- | --- |
| `crmCreateOpportunity` | 07-22-01 | 3 klausa `OR` berbasis `Declare_crmOperatorAccess` / `evaluateWhen` | `Harness/SFAPortalOpportunities`, `Section/SFAPortalOpportunitiesHeader` |
| `isApproved` | 01-01-52 | `pyWorkPage.PolicyTreatyIn.IsApproved = 1` | **13 pemanggil** (3 Activity, 4 DataTransform, 1 Flow, 5 Section) |
| `isClaimTreaty` | 01-01-56 | 7 klausa `OR`: `.PolicyTreatyIn.ClaimType = "XOL"` · `.Claim != 0` · `.SalvageValue != 0` · `.ExcessLoss != 0` · `.ClaimPaymentType = "Salvage"` · `= "Claim"` · `.Quotation.ProportionalType != "Proportional"` | `Flow/InputRealizationTreatyIn` |
| `IsNotAdmin` | 01-01-55 | `OperatorID.pyPosition != "Admin"` | 2 (portal SFA) |
| `IsOfferFacIn` | 01-01-55 | `.Quotation.BusinessFac = "F"` | `Flow/InputQuotation` |
| `IsOperatorLife` | 01-01-55 | `OperatorID.pyWorkGroup = "ReasLife"` | `Section/SFAPortal_OpportunitiesList` |
| `isSellingModeB2B` / `B2BB2C` / `B2C` | 07-22-01 | `getDataSystemSetting("PegaCRM-","SellingMode")` = `"B2B"` / `"B2B_B2C"` / `"B2C"` | 3 masing-masing (portal SFA) |
| `IsSPVCreate` | 01-01-81 | 2 klausa `OR` atas `pyWorkPage.pxCreateOperator` (**2 identitas operator**; nilai tidak disalin) | `Flow/InputRealizationTreatyIn` + disebut di `When/IsSPVTreaty1`, `When/IsTreaty1` |
| `IsSPVTreaty1` | 01-01-81 | `OperatorID.pyTelephone = "SPVTREATY1"` | `Flow/InputRealizationTreatyIn` |
| `IsTreaty1` | 01-01-81 | `OperatorID.pyTelephone = "TREATY1"` | `Flow/InputRealizationTreatyIn` |
| `NopolisEmpty` | 01-01-56 | `pyWorkPage.PolicyTreatyIn.PolicyNo = ""` | `Flow/InputRealizationTreatyIn` |
| `pyIsIpadOrDesktop` | 07-22-01 | `pyIsIPad` benar `\|\|` `pxRequestor.pxDeviceType = desktop` | 2 (portal SFA) |
| `pyIsMobile` | 08-02-01 | `@(Pega-RulesEngine:Utilities).pzIsMobile(tools)` | `Section/pyAttachmentScreen`, `Section/SFAPortal_OpportunitiesList` |
| `ToDeptHeadUWLife` | 01-01-81 | `pyWorkPage.LetterNo = "DEPTHEADUWLIFE"` | **`Flow/InputInwardFacultativeOffer`** |
| `ToTREATYDEPTHEAD` | 01-01-81 | `pyWorkPage.LetterNo = "TREATYINDEPTHEAD"` | `Flow/InputRealizationTreatyIn` |
| `TreatyMasterInEDM` | 01-01-56 | `pyWorkPage.TreatyIn.EDMState = "1" \|\| "2" \|\| "3"` | `Section/DetailPoliciesNonProportional` |

Dua pola yang layak dicatat:

- **`OperatorID.pyTelephone` dipakai sebagai kode peran.** `IsTreaty1` dan `IsSPVTreaty1` membaca
  kolom nomor telepon operator dan membandingkannya dengan `"TREATY1"` / `"SPVTREATY1"`.
  **[terverifikasi]**. Kandidat perbaikan; migrasi harus memindahkannya ke atribut peran yang
  sebenarnya, **tetapi keputusannya milik bisnis**.
- **`pyWorkPage.LetterNo` dipakai sebagai token routing.** `ToDeptHeadUWLife` dan `ToTREATYDEPTHEAD`
  membandingkannya dengan `"DEPTHEADUWLIFE"` / `"TREATYINDEPTHEAD"` — bukan nomor surat.
  **[terverifikasi]**. Ini menjelaskan mengapa rule `LetterNoNull` (salah satu dari lima `When`
  berkondisi kosong menurut aturan proyek §4.5) penting.

---

## 5. Rule eksklusif tanpa pemanggil

### 5.1 Perintah audit

```powershell
$nb="D:\migrasi\RNM\NB FacIn"
$corpus = foreach ($d in (Get-ChildItem $nb -Directory)) { foreach ($f in (Get-ChildItem $d.FullName -Filter *.xml -File)) {
  [pscustomobject]@{ Dir=$d.Name; Name=$f.BaseName; Text=(Get-Content $f.FullName -Raw) } } }
# 2.083 berkas; untuk tiap rule eksklusif cari kemunculan namanya (batas kata) di berkas LAIN
# -> TANPA PEMANGGIL: 3 / 139
```

### 5.2 Hasil: 3 rule

| Rule | Kelas | Isi | Penilaian |
| --- | --- | --- | --- |
| `Activity/InputDtlAdditional_PreAct` | `ASM-FW-GISFW-Data-Coverage` v01-01-52 | kondisi `\.Coverage=="21" \|\| "51" \|\| "52" \|\| "53"`; memakai `Apply-DataTransform AddCoverageProRate`, properti `.TSI`, `.TSITJH`, `.Day` | **kandidat kode mati.** Arti kode cover 21/51/52/53 **belum terverifikasi**. |
| `Activity/ValidateCoverage` | `ASM-FW-GISFW-Data-Coverage` v01-01-52 | validasi cover ganda pada `pyWorkPage.VehicleList(i).CoverageList(j).AdditionalCoverage(k)`; memakai `@Utilities.countInPageList(...)>1`, rujuk `IsEDM`, `IsFlagOldData`, `IsProposalTransfer` | **kandidat kode mati** untuk jalur Fac In (berorientasi kendaraan). |
| `Harness/SFAPortalOpportunities` | `PegaCRM-Portal`, ruleset `PegaCRM-SFA` v07-22-01 | harness portal; merujuk `SFAPortal_Opportunities`, `SFAPortalOpportunitiesHeader`, `crmCreateOpportunity`, `isSellingModeB2C`, `isSellingModeB2BB2C` | **BUKAN kode mati.** Harness portal dipasang lewat `Data-Admin-Operator-AccessGroup`, tipe rule yang **tidak diekspor** ke folder mana pun. Ini "belum ditemukan rujukannya". |

### 5.3 Pembedaan tegas

| Kategori | Definisi | Jumlah |
| --- | --- | --- |
| **Kode mati (kandidat)** | tidak dirujuk dari mana pun **dan** mekanisme rujukannya ada di dalam korpus | **2** (`InputDtlAdditional_PreAct`, `ValidateCoverage`) |
| **Belum ditemukan rujukannya** | tidak dirujuk, tetapi mekanisme rujukannya (`Rule-Access-*`, konfigurasi portal, siklus lain) **tidak ada di korpus** | **1** (`SFAPortalOpportunities`) + seluruh G2 lain |

Daftar kode mati **tidak dapat ditutup** tanpa ekspor `Rule-Access-Role-Obj`,
`Data-Admin-Operator-AccessGroup`, dan definisi portal. **[pertanyaan terbuka]**

---

## 6. Daftar prioritas — rule eksklusif NB yang MEMBLOKIR implementasi

| # | Rule | Mengapa memblokir | Yang dibutuhkan |
| --- | --- | --- | --- |
| **1** | `Activity/GetLimitAkseptasi_ActFlow` | Satu-satunya orkestrasi tangga akseptasi NB; tidak ada padanannya di RNW/EDM. Tanpa ini tidak ada jalur persetujuan New Business sama sekali. | Konfirmasi bisnis atas: (a) tautologi `>= \|\| <=` yang muncul 15 kali; (b) 4 nomor polis literal sebagai gerbang; (c) dua konvensi ejaan jabatan pada `.CARI1`. |
| **2** | `Activity/GetLimitAkseptasiLife_Act` + `RDBList/GetLimitAkseptasiLife_SQL` | Tangga akseptasi Life terpisah, atas tabel `M_LIMIT_LIFE` yang tidak disebut di tempat lain mana pun. | Skema + isi `M_LIMIT_LIFE`; aturan pemilihan baris (query hanya memberi batas bawah). |
| **3** | `When/IsOfferFacIn` | Menggerbangi `Flow/InputQuotation` — titik masuk pembuatan offer. `pyConditionString` dan `pyConditionValue1String` **bertentangan** (§4). | Konfirmasi mana yang berlaku, dan arti kode `"F"` pada `.Quotation.BusinessFac`. |
| **4** | `RDBList/SaveOfferJson_SQL` → `POOLDATA.PEGA_M_JSON_OFFER` | Satu-satunya jalur persistensi offer NB. Isi prosedur tidak ada di korpus. | Sumber stored procedure `PEGA_M_JSON_OFFER` (7 argumen, 3 di antaranya `OUT`). |
| **5** | `RDBList/GenerateNoPolicy` | Berkelas `Int-POLISTREATYIN` dan memakai `JSON_POLIS_TREATYIN_SEQ` + literal `'.T'`, **tetapi juga dipanggil `Activity/SaveJsonPolicyFacIn_Act`** (`BrowsePage = PolicyFacIn`). Nomor polis Fac In bisa salah format. | Konfirmasi apakah ada rule `GenerateNoPolicy` kedua pada kelas Fac In yang **hilang akibat penamaan ekspor** (lihat catatan di bawah). |
| **6** | `Flow/InputInwardFacultativeOffer` | Flow entri siklus NB. | — (isinya perlu dibaca utuh pada pass berikutnya) |
| **7** | `Flow/InputQuotation` | Flow yang memanggil `InputInwardFacultativeOffer` dan `InputRealizationTreatyIn`. Berada di **ruleset personal** v01-01-01 — versi yang terekspor mungkin *checkout* yang belum di-check-in. | Konfirmasi versi yang berlaku di produksi. |
| **8** | `SystemSettings/LinkService` | Jalur konfigurasi endpoint (`=ResponLink.URL`) yang dipakai 3 `ConnectREST`, termasuk integrasi AI pihak ketiga. | Isi tabel `M_LINK_SERVICE` (`KATEGORI_1` + `KATEGORI_2`) — tidak ada di korpus; dan persetujuan manusia per aturan proyek §6. |

**Catatan penamaan ekspor (butir 5).** `pyRuleName` sebuah `RDBList` adalah triplet
`<Kelas> <Access> <RequestType>` — contoh terbaca:
`ASM-FW-GISFW-Int-POLISTREATYIN ASM GenerateNoPolicy`. Namun **nama berkas ekspor hanya memakai
`RequestType`** (`GenerateNoPolicy.xml`). Folder NB memuat **217 berkas `RDBList` dari 55 kelas
berbeda**; bila dua kelas memakai `RequestType` yang sama, kedua rule akan bertabrakan pada satu
nama berkas dan **hanya satu yang tersisa di ekspor**. **[dugaan]** — belum dibuktikan bahwa
tabrakan benar-benar terjadi, tetapi risikonya nyata dan menjelaskan ambiguitas butir 5.

```powershell
$nb="D:\migrasi\RNM\NB FacIn\RDBList"
$cl=@(); foreach ($f in Get-ChildItem $nb -Filter *.xml) { $cl += [regex]::Match((Get-Content $f.FullName -Raw),'<pyClassName>([^<]+)').Groups[1].Value }
"kelas unik: {0} ; berkas: {1}" -f ($cl|Select-Object -Unique).Count, $cl.Count   # -> 55 ; 217
```

### 6.1 Tidak memblokir, tetapi wajib direplikasi apa adanya

- `When/ToDeptHeadUWLife` — routing via `LetterNo`.
- `When/IsSPVCreate` — guard atas 2 identitas operator.
- `DecisionTable/BusinessType_DeT` (kelas `ASM-FW-GISFW-Work`, v01-01-87) dan
  `DecisionTree/Tree_ShortPeriod` (kelas `ASM-FW-GISFW-Data-Coverage`, v01-01-52) — isi baris/cabang
  **belum dibaca** pada pass ini.
- `Section/InwardFacIn` (10 pemanggil), `Section/InstallmentList` + `FlowAction/InstallmentList`
  (11 pemanggil masing-masing), `Section/FormulaTreatyCapacityDesc` (3 pemanggil).

### 6.2 Jangan dimigrasikan

- **15 rule G2** (produk Pega / framework SFA-CRM).
- ~~**80 rule G1** — Treaty Inward sebagai lini produk terpisah.~~
  ✅ **DICABUT oleh K-004: Treaty Inward termasuk lingkup — 80 rule G1 DIMIGRASIKAN.**
  Konsekuensinya memang menambah ±80 rule ke lingkup NB, dan itu sudah diterima work owner.
  Cakupannya terbatas pada jejak Treaty di dalam `NB FacIn\` (K-005).
- `RDBList/GetCountClaim` sebagaimana adanya (membaca `DATAPEGA.PC_*`).
- `RDBList/GetSQLDate` dan `GetCurrentDate` (`select sysdate from dual`).

---

## 7. Pertanyaan terbuka

**MEMBLOKIR**

1. ~~**Apakah Treaty Inward masuk lingkup proyek ini?**~~ ✅ **TERJAWAB — K-004: ya, masuk.**
   80 dari 139 rule eksklusif NB adalah Treaty Inward (kelas `Data-PolicyTreatyIn`, `Data-Portal`,
   `Data-TreatyInLimits`, `Int-POLISTREATYIN`, `Int-TREATY_IN`) dan seluruhnya dimigrasikan.
   Cakupan terbatas pada jejaknya di dalam `NB FacIn\` (K-005).
2. ~~**Mana yang berlaku pada `When/IsOfferFacIn`**~~ ✅ **TERJAWAB — K-002: ekspresi tersimpan
   `pyConditionValue1String` (`.Quotation.BusinessFac = "F"`).** Teks `pyConditionString`
   (`Kode Bisnis = "02"/"58"/"SB"/"SG"`) dicatat sebagai kandidat perbaikan, tidak diimplementasikan.
3. **Skema, isi, dan aturan pemilihan baris tabel `M_LIMIT_LIFE`** (`GetLimitAkseptasiLife_SQL`).
   Query hanya memberi `LIMIT_BOTTOM <=` tanpa batas atas.
4. **Isi stored procedure `POOLDATA.PEGA_M_JSON_OFFER`** (persistensi offer NB), **dan** —
   karena Treaty Inward masuk lingkup (K-004) — `POOLDATA.PEGA_JSON_POLIS_TREATYIN` serta
   `POOLDATA.PEGA_TREATY_IN`. Ketiganya kini sama-sama memblokir, bukan bersyarat.
5. **Apakah `RDBList/GenerateNoPolicy` yang terekspor memang yang dipakai jalur Fac In?**
   Kelasnya `Int-POLISTREATYIN`, sekuensnya `JSON_POLIS_TREATYIN_SEQ`, dan ia menyisipkan literal
   `'.T'` ke nomor polis — padahal dipanggil juga dari `Activity/SaveJsonPolicyFacIn_Act`.
6. **Apakah 15 tautologi `TotalTSI>=Limit || TotalTSI<=Limit` di `GetLimitAkseptasi_ActFlow`
   disengaja?** Bila tidak, tangga akseptasi sistem lama tidak pernah benar-benar membandingkan TSI
   dengan limit pada 15 titik itu — dan sistem baru yang "memperbaikinya" akan menghasilkan jalur
   persetujuan yang berbeda.
7. **Apakah 4 nomor polis literal di `GetLimitAkseptasi_ActFlow` masih relevan?** Bila ya, nilainya
   harus dipindah ke data, bukan kode. Keputusan milik bisnis.
8. **Isi tabel `M_LINK_SERVICE`** (`KATEGORI_1` + `KATEGORI_2`) yang dirujuk
   `SystemSettings/LinkService` → `=ResponLink.URL`. Termasuk endpoint `ConnectREST/ServiceGoogle`
   yang menurut aturan proyek §6 memerlukan tinjauan keamanan.

**TIDAK MEMBLOKIR**

9. Apakah `OperatorID.pyTelephone` sebagai kode peran (`"TREATY1"`, `"SPVTREATY1"`) harus
   dipertahankan pada sistem baru, atau dipindahkan ke atribut peran yang sesungguhnya?
10. Dua konvensi ejaan jabatan pada `.CARI1` (`"MANAGER TEKNIK"` vs `"MANAGERTEKNIK"`,
    `"KADIV FACULTATIVE"` vs `"KADIVFACULTATIVE"`) — keduanya ada di data, atau salah satu cabang
    mati?
11. Arti kode cover `"21"`, `"51"`, `"52"`, `"53"` pada `Activity/InputDtlAdditional_PreAct`
    — belum terverifikasi.
12. Apakah 11 rule G4 (kendaraan, kargo, warranty) benar-benar dipakai pada offer Fac In?
    9 di antaranya masih dirujuk, dua di antaranya lewat FlowAction berakhiran `_FacIn`.
13. `Flow/InputQuotation` berada di ruleset personal seorang pengembang, v01-01-01. Versi mana yang
    aktif di produksi?
14. Nama `SystemSettings` di folder RNW dan EDM (folder NB memakai `LinkService`; ketiga folder
    hanya punya satu `SystemSettings` masing-masing, dan namanya berbeda).
15. Apakah `RDBList/SearchJobID` dan `RDBList/SearchJobIDSQL` (SQL identik, ruleset `GISFWInt`
    v01-01-34) memang perlu dua rule?
16. Isi `DecisionTable/BusinessType_DeT`, `DecisionTable/isApproved`, dan
    `DecisionTree/Tree_ShortPeriod` — belum dibaca pada pass ini.
17. Daftar kode mati tidak dapat ditutup tanpa ekspor `Rule-Access-*` dan konfigurasi portal (§5.3).
