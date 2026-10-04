# Discovery Endorsement — E-6: penutup

> **Sumber:** `D:\migrasi\RNM\Endorsment Fac In\` (banding `NB FacIn\`, `RNW Fac In\`, `DDL\`).
> READ-ONLY. Korpus Treaty (`RNM_BRD\`) **tidak dibaca** (K-005). Label mengikuti `CLAUDE.md` §3.
>
> **Ini konsolidasi, bukan temuan baru.** Tujuannya satu daftar tuntas yang siap jadi (a) bahan paket
> DBA/IT dan (b) titik masuk spec EDM.
>
> Kontrak volatile yang berlaku: **23 tag + isi blok `pzIndexes` diabaikan** (amandemen K-043).
> Lingkup EDM: **708 berkas** = 354 berbeda + 354 EDM-only.

---

## 1. Angka akhir discovery EDM

`[terverifikasi]` Di bawah kontrak K-043 yang sudah diamandemen:

| Ukuran | Jumlah |
| --- | ---: |
| Total `.xml` di `Endorsment Fac In\` | **2.061** |
| Bernama sama dengan NB (peka huruf, per folder) | 1.707 |
| — **identik** | **1.353** (79,3 %) |
| — **berbeda** | **354** (20,7 %) |
| **EDM-only** | **354** |
| **Lingkup kerja EDM** | **708** |

### Perbedaan per tipe rule — dan satu butir yang ikut tertutup

| Tipe | Identik | Berbeda | % beda | Perubahan dari E-1 |
| --- | ---: | ---: | ---: | --- |
| Activity | 373 | **115** | 23,6 % | dari 194 |
| Section | 288 | 65 | 18,4 % | dari 101 |
| When | 119 | **56** | 32,0 % | dari 63 |
| FlowAction | 170 | 41 | 19,4 % | dari 45 |
| DataTransform | 105 | 23 | 18,0 % | dari 35 |
| ReportDefinition | 92 | 18 | 16,4 % | dari 27 |
| RDBList | 148 | 18 | 10,8 % | tetap |
| **DecisionTable** | 0 | **10** | **100 %** | tetap |
| Harness | 28 | 5 | 15,2 % | dari 13 |
| **Flow** | 0 | **1** | **100 %** | tetap |
| **DecisionTree** | 0 | **1** | **100 %** | tetap |
| DataPage | 26 | 1 | 3,7 % | tetap |
| **ConnectREST** | **3** | **0** | **0 %** | ⛔ **dari 3/3 berbeda** |
| SystemSettings | 1 | 0 | 0 % | tetap |

⛔ **`ConnectREST` tidak lagi berbeda.** Ketiga berkas integrasi yang di E-1 tampak "100 % berbeda"
ternyata **beda penomoran indeks saja**. Ini menutup separuh **E-Q11**. Yang **tetap** 100 % berbeda:
`DecisionTable` (10), `Flow` (1), `DecisionTree` (1).

---

## 2. Konsolidasi E-Q1 … E-Q32

Satu entri per nomor. Nomor yang muncul di beberapa berkas digabung — **tidak dihitung ganda**.

### 2.1 TERTUTUP (7)

| # | Isi ringkas | Ditutup oleh |
| ---: | --- | --- |
| **E-Q1** | Resolusi rule Pega peka huruf? Berdampak pada 12 berkas varian kapitalisasi | **K-042** menetapkan angka resmi **peka huruf** (1.707/354). Pertanyaan teknis ke IT tidak lagi memblokir |
| **E-Q2** | Mengapa EDM hanya punya 2 Flow? | **E-4 §4.2** — delta EDM digerakkan dari **Section** (186 sisi rujukan vs Flow 3). Alur endorsement memang terpusat |
| **E-Q3** | Apakah K-038 berlaku untuk EDM? | **Digabung ke E-Q25** — premisnya gugur di E-4 §6 (bukan duplikat baru) |
| **E-Q4** | Amandemen K-042 — lingkup 866 | **K-043** |
| **E-Q14** | `EDMPremiMenjadi` banyak bentuk | **E-3 §5** — dituntaskan: **22 penugasan · 12 rumus literal · 11 semantik**, tabel per (lini × `EdmType` × jalur). Sisanya keputusan rancangan, bukan pertanyaan korpus |
| **E-Q20** | 47 gerbang lini bisnis dinonaktifkan | **Digabung ke E-Q16** — bagian dari pertanyaan semantik yang sama |
| **E-Q23** | `pzIndexes` masuk kontrak volatile? | **Amandemen K-043** — 158 tereproduksi, lingkup 866→708 |

### 2.2 TERVERIFIKASI faktanya — menunggu keputusan (4)

Fakta korpusnya **sudah terbukti dengan bukti baris**; yang tersisa keputusan bisnis.

| # | Fakta terbukti | Pemilik keputusan | Memblokir? |
| ---: | --- | --- | --- |
| **E-Q18** | `RateOld` guard self-referential `@If(.RateOld!="",.Rate,0)` — **6 dari 6**; `PremiumOld` **1 dari 18** | work owner + Aktuaria | tidak |
| **E-Q19** | `PremiNusantaraReOld` fallback `"0"` berkutip — **2** (L6294, L7527) vs **50** numerik; cache `pyExpressionGadget` bahkan berbeda dari yang dieksekusi | IT + Aktuaria | tidak |
| **E-Q21** | Dari 12 cabang perhitungan, **hanya LIFE `x.2`** (@L4513) menerima `EdmType==1`; 6 langkah penihilan menerima `==1` di semua lini | work owner + Underwriting | tidak |
| **E-Q11** *(separuh)* | `ConnectREST` **tidak lagi berbeda**; `DecisionTable` 10/10, `Flow` 1/1, `DecisionTree` 1/1 **tetap** berbeda | internal — tinggal dibaca isinya | tidak |

### 2.3 TERBUKA (21)

| # | Isi ringkas | Pemilik | MEMBLOKIR |
| ---: | --- | --- | :---: |
| **E-Q5** | Bagaimana `OutData.pxResults(1).CARI2` diisi; apa isi kolom `CARI2` | IT + DBA | ⛔ **YA** — 6 predikat COB tak dapat diimplementasikan |
| **E-Q6** | `IsUW` — workbasket marketing dihitung sebagai underwriter di EDM | Underwriting | ⛔ **YA** — menentukan siapa boleh menyetujui |
| **E-Q7** | `IsClaim` — NB uji prefiks ID, EDM uji keberadaan halaman; kelas beda | work owner | tidak |
| **E-Q8** | `IsTravel` — EDM dua sumber, NB satu | work owner | tidak |
| **E-Q9** | 32 rule `When` kategori (d) — kondisi identik, isi beda di luar kondisi | internal | tidak |
| **E-Q10** | `DecisionTable\IsUWAccepted` + `DecisionTree\Tree_ShortPeriod` berbeda — menyentuh modul NB terspec | internal → work owner | ⛔ **YA** — tangga akseptasi & periode pendek |
| **E-Q12** | Registry predikat COB EDM: registry kedua atau satu registry bersumber-ganti | rancangan | menunggu E-Q5 |
| **E-Q13** | `IsNotEDM` **bukan** negasi `IsEDM` (jalur properti berbeda) | work owner + IT | tidak |
| **E-Q15** | `pyStepsRepeatDef` / `pyStepsPreCondParams` memuat nama orang + email — perluasan K-025 | work owner | tidak |
| **E-Q16** | `pyStepsPreCondition=false` — **617** kemunculan di 173 berkas, **47** di antaranya gerbang lini bisnis | **IT** | ⛔ **YA** — menentukan validitas seluruh pemetaan |
| **E-Q17** | Hanya `_FIRE` menyetel `IsProRate="Prorate"` | work owner | tidak |
| **E-Q22** | Himpunan lini lapis B ≠ lapis C (B: Travel tanpa Life; C: Life tanpa Travel); `_LIFE` pakai dua penanda | work owner | tidak |
| **E-Q24** | 5 berkas tidak disebut di mana pun: `GetMasterKlausulAge(1)` · `SearchClobClauseSQL(1)` · `IsCustomBond` · `SaveClausePA_Act` · `SFAPortalEndorsement` | work owner + IT | tidak |
| **E-Q25** | Apakah K-038 (pilih-tertanggung dibuang) berlaku untuk EDM — **premis teknis sudah gugur** | **work owner** | tidak |
| **E-Q26** | Apakah endorsement **Life** masuk lingkup tahap ini — 52 berkas ber-sub-graf lengkap | **work owner** | ⛔ **YA** — menentukan lingkup |
| **E-Q27** | Hanya 2 Section menampilkan `*Old`, padahal lapis B mengisi 52 properti | work owner + Underwriting | tidak |
| **E-Q28** | `SaveFacinProdEDMBonding_Act` ada di `DDL\`, tidak di korpus — dilipat masuk (pola K-006)? | **work owner** | ⛔ **YA** — lini Bonding tanpa jalur produksi |
| **E-Q29** | Kualifikasi skema: `facinproduction` vs `pooldata.facinproduction` | **DBA** | ⛔ **YA** — guard idempotensi bisa tak berfungsi |
| **E-Q30** | Dua semantik "versi terakhir `PRODKE`": `COUNT-1` vs `TGL_INPUT desc` | **DBA** + work owner | ⛔ **YA** — before-image bisa salah versi |
| **E-Q31** | Tidak ada cabang Life **maupun Travel** di `SaveFacinProdAllEDM_Act` | work owner | tidak |
| **E-Q32** | Guard idempotensi fac out hanya dipakai keluar bila `Type=="7"` | work owner | tidak |

**Rekap:** 7 tertutup · 4 terverifikasi menunggu keputusan · 21 terbuka — **8 di antaranya MEMBLOKIR**.

---

## 3. Blocker per pemilik

### 3.1 DBA — 4 blocker

| # | Yang dibutuhkan | Bukti lokasi panggilan |
| --- | --- | --- |
| **D-1** | Isi **4 stored procedure** | `POOLDATA.INSERTJSONPOLIS` ← `RDBList\INSERTJSON_JSONPOLISEDM_FACIN` · `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` ← `RDBList\DeleteDataProduction` · `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` ← `RDBList\GetSequenceNumber_SQL` · blok PL/SQL anonim ← `RDBList\MachingDataFacin_Sql` |
| **D-2** (E-Q29) | Apakah `facinproduction` (tanpa skema) dan `pooldata.facinproduction` objek yang sama | `RDBList\InsertTreatyProduction_Sql` vs `RDBList\CekFacin_Sql` |
| **D-3** (E-Q30) | Adakah constraint yang menjamin `PRODKE` rapat dari 0 tanpa lubang | `GetProdKeEDM_SQL` & `GetEDMOldData_SQL` (`COUNT-1`) vs `GetProdKeOldData_SQL` (`TGL_INPUT desc`) |
| **D-4** (arsip §10 #10) | Isi tabel `ASM-FW-GISFW-Int-OPENPROTEKSI_EDM` + siapa mengisinya | klep pembatal 4 gerbang penolakan, `Param.Type` 1–4 |

⚠️ Tambahan yang sudah diketahui menunggu: tabel flat + `ALL_SOURCE` untuk jalur produksi NB
(pembanding jalur produksi EDM).

### 3.2 IT — 2 blocker

| # | Yang dibutuhkan | Skala |
| --- | --- | --- |
| **I-1** (E-Q16/E-Q20) | Arti `<pyStepsPreCondition>false</pyStepsPreCondition>` — prakondisi tidak dievaluasi, atau langkah dinonaktifkan? | **617** kemunculan di **173** berkas; **47** di antaranya gerbang lini bisnis — termasuk langkah 4 `SetOldData` (seluruh lapis B) dan langkah 20 `SaveEDMToJsonPolicy_Act` (pengambilan `PRODKE`) |
| **I-2** (arsip §10 #1) | Versi rule mana yang berlaku di produksi — `SetToInbox_ACT` dan `SetBanding_ACT` | 86 rule bernama sama berbeda `pyRuleSetVersion`; rekonsiliasi paralel run mustahil tanpa satu ekspor produksi tunggal |

### 3.3 Work owner — 4 keputusan

| # | Keputusan | Dampak lingkup |
| --- | --- | --- |
| **W-1** (E-Q26) | Apakah endorsement **Life** masuk lingkup tahap ini | **52 berkas** ber-sub-graf lengkap — bisa memangkas atau menambah subsistem penuh |
| **W-2** (E-Q25) | Apakah **K-038** (pilih-tertanggung dibuang) berlaku untuk EDM | 4 berkas; premis teknis sudah gugur |
| **W-3** (E-Q28) | Apakah `SaveFacinProdEDMBonding_Act` dari `DDL\` dilipat masuk | menentukan apakah lini **Bonding** punya jalur produksi |
| **W-4** (E-Q31) | Apakah endorsement **Travel** memang tidak masuk produksi | Travel ditangani `SetOldData` tetapi tidak punya cabang produksi |

### 3.4 Underwriting / Aktuaria — 3 keputusan

**E-Q6** (marketing sebagai UW) ⛔ memblokir · **E-Q18** (`RateOld`) · **E-Q21** (`EdmType==1` hanya Life).

---

## 4. Audit rujukan menggantung — menyeluruh

`[terverifikasi]` Enam mekanisme pemanggilan dipindai di seluruh **2.061** berkas EDM; nama yang
dirujuk dicocokkan ke **tiga folder korpus** (peka huruf, per tipe rule).

| Ukuran | Jumlah |
| --- | ---: |
| Rujukan unik (mekanisme + nama) diperiksa | **989** |
| Tidak ditemukan di 3 korpus | **94** |
| — A. bawaan platform Pega (prefiks `px`/`pz`/`py`) | 13 |
| — B. metode bawaan / bukan rule aplikasi | 16 |
| — **C. ADA di `DDL\` → pola K-006** | **4** |
| — **D. HILANG dari semua lokasi → `panic` (§4.5)** | **61** |

⚠️ **Dua kesalahan klasifikasi saya sendiri diperbaiki sebelum angka ini keluar:** `pySourceName`
dicocokkan hanya ke `ReportDefinition` (padahal bisa **DataPage** — `D_EnumerationList` jelas ada di
korpus), dan `Obj-Save`/`ObjSave` dihitung sebagai rule (padahal **metode**). Tanpa koreksi itu,
angkanya 108, bukan 94.

### 4.1 Kategori C — ada di `DDL\`, **bukan** `panic`

`[terverifikasi]` Empat rujukan. Sesuai pola amandemen **K-006**, cabang pemanggilnya
**tidak dihapus dan tidak divonis usang**:

| Mekanisme | Nama | Dirujuk dari |
| --- | --- | --- |
| `Call` | **`SaveFacinProdEDMBonding_Act`** | `Activity\SaveFacinProdAllEDM_Act` ← **E-Q28** |
| `Call` | `GetLimitAkseptasi_Act2` | `Activity\SetValidateDate_PostAct` |
| `RequestType` | `GetLimitAkseptasi1SA_Act` | `Activity\GetLimitAkseptasi_Act` |
| `RequestType` | `GetLimitAkseptasi2SA_Act` | `Activity\GetLimitAkseptasi_Act` |

⛔ **Tiga dari empat menyentuh tangga akseptasi.** `GetLimitAkseptasi_Act` merujuk dua rule SQL yang
hanya ada di `DDL\`. Ini memperluas E-Q28 dari satu berkas menjadi **empat**.

### 4.2 Kategori D — hilang dari semua lokasi → `panic`

`[terverifikasi]` **61** rujukan. Yang menyentuh jalur EDM inti:

| Mekanisme | Nama | Dirujuk dari | Catatan |
| --- | --- | --- | --- |
| `RequestType` | **`GenerateEndorsementNo`** | `SaveEDMToJsonPolicy_Act` · `GenerateNopolis_Act` | penomoran endorsement **dan** polis |
| `RequestType` | **`GenerateEDMNoLife`** | `SaveEDMToJsonPolicy_Act` | penomoran endorsement Life |
| `Call` | `SumTSIPremiSpreadRNMMultiCob_Act` | `SumTSIPremiSpreadedRNM_Act` | spreading multi-COB |
| `Call` | `EDMPayment` | `InputDtlPayment_PreAct` | pembayaran EDM |
| `Call` | `AveragePremium_Act` | `FillActInstallment` | premi rata-rata |
| `Call` | `SetCurrencyServis_act` | `InputAddendumFacIn_PreAct` | mata uang, di pra-aktivitas layar EDM |
| `pyActivity` | `CountNetPremiumAsm_Act` · `SetRateErrorMessageASMRate_Act` · `SetErrorMessageRIComm_Act` | `Section\OfferFacIn_NusaRe` (+`_IsUW`) | masing-masing 4× |
| `pyActivity` | `SetPctProrateAct` · `CalculateDiffDate_Act` · `StartDateLog_Act` · `EndDateLog_Act` · `CheckMaxRange` | `Section\PeriodePolicy` | periode & prorata |
| `RequestType` | `MachingDataOffer_Sql` · `INSERTERROROFFER_Sql` · `InsertOfferProductionBackup_Sql` | `SaveOfferProduction_Act` | masing-masing 7× |

Sisanya berupa lampiran/korespondensi/cetak (`FacOfferCorrespondenceDownload`, `GeneratePDFRISliping`,
`DeleteImage_Act`, …), DataPage master (`D_TypeList`, `D_CityList`, `D_DistrictList`, …), dan
FlowAction lokal (`ConfirmDateChange`, `PreventDeleteAttachment_LA`, …).

⛔ **Semua kategori D gagal keras (`panic`)** per `CLAUDE.md` §4.5 — **tidak** ditebak, **tidak**
di-stub, **tidak** dihapus cabangnya. ⚠️ Per **K-006**: "hilang dari ekspor" **bukan** bukti usang.
Sebagian besar kemungkinan ada di ruleset yang tidak ikut terekspor — itu pertanyaan untuk IT
bersama **I-2**, bukan vonis.

```powershell
# audit rujukan menggantung — enam mekanisme, cocokkan ke 3 korpus lalu DDL\
$edR="D:\migrasi\RNM\Endorsment Fac In"
$ada=@{}; foreach($a in @('NB FacIn','RNW Fac In','Endorsment Fac In')){
  Get-ChildItem "D:\migrasi\RNM\$a" -Recurse -Filter *.xml -File | % { $ada["$($_.Directory.Name)|$($_.BaseName)"]=1 } }
$ddl=@{}; Get-ChildItem "D:\migrasi\RNM\DDL" -Recurse -File -EA SilentlyContinue | % { $ddl[$_.BaseName]=1 }
$mek=@(
 @('Call','<pyStepsActivityName>Call\s+([A-Za-z0-9_.\-]+)</pyStepsActivityName>',@('Activity')),
 @('RequestType','<RequestType>([^<]+)</RequestType>',@('RDBList')),
 @('pyActivity','<pyActivity>([A-Za-z0-9_\-]+)</pyActivity>',@('Activity')),
 @('pyLocalAction','<pyLocalAction>([A-Za-z0-9_\-]+)</pyLocalAction>',@('FlowAction')),
 @('pySourceName','<pySourceName>([A-Za-z0-9_\-]+)</pySourceName>',@('ReportDefinition','DataPage')),
 @('pySectionName','<pySectionName>([A-Za-z0-9_\-]+)</pySectionName>',@('Section')))
# => 989 rujukan unik / 94 tak ketemu / 13 platform / 16 metode / 4 di DDL / 61 hilang
```

---

## 5. Kode usang & kejanggalan — diport apa adanya

Sesuai `CLAUDE.md` §1: perbaikan dipisahkan dari migrasi; memutuskannya milik bisnis.
⛔ **Tidak satu pun diperbaiki.** Memperbaikinya diam-diam menghasilkan selisih angka yang tidak
dapat dijelaskan saat rekonsiliasi paralel run.

| # | Kejanggalan | Bukti | Menyentuh |
| ---: | --- | --- | --- |
| 1 | **`RateOld` guard self-referential** — `@If(.RateOld!="",.Rate,0)` **6/6**; `PremiumOld` **1/18** | E-3 §3.4 | SELISIH rate · E-Q18 |
| 2 | **`PremiNusantaraReOld` fallback `"0"` berkutip** — 2 penugasan (L6294, L7527) vs 50 numerik | E-3 §3.4 | tipe kolom · E-Q19 |
| 3 | **Blok `EdmType==1` tidak dihitung ulang** — 11 dari 12 cabang perhitungan `==2` saja; hanya LIFE `x.2` menerima `==1` | E-3 §5.3 | premi pembatalan · E-Q21 |
| 4 | **Label MBU bergerbang `IsMarineCargo`** — `CountDataEDMElse` [2] deskripsi *"…if MBU business"*, gerbang `IsMarineCargo` | E-2 §2, E-3 §8 | MBU tak dihitung; MC dihitung 2× |
| 5 | **`EDMPremiMenjadi` tak seragam** — 22 penugasan, **11 rumus semantik** di 7 berkas | E-3 §5 | modul premi · E-Q14 |
| 6 | **`PCT_BROKERGARE_FEE` salah ejaan di skema** (`BROKERGARE`, bukan `BROKERAGE`) | E-2 §1.5, E-5 §3.1 | skema tidak berubah (§4.3) |
| 7 | **Ketidaksetangkupan `*_MENJADI`/`*_SELISIH`** — `facinproduction` 5 vs 8; `FACOUTPRODUCTION` 2 vs 5 | E-5 §3.1, §5.2 | tiga pasangan pakai kolom tanpa sufiks |
| 8 | **Dua semantik `PRODKE`** dalam satu rantai | E-5 §4.2 | before-image · E-Q30 |
| 9 | **Guard idempotensi fac out** hanya dipakai keluar bila `Type=="7"` | E-5 §5.2 | baris ganda · E-Q32 |
| 10 | **`JSON_POLIS_MONITORING` hanya ditulis untuk prefiks `EDMT-`** | E-5 §6 | arsip §10 #23 |
| 11 | **Paradoks label *"batal sejak semula"*** — per K-029 = `EdmType==1`, tetapi bergerbang `==2` di 5 dari 6 lini | E-3 §5.3 | label ≠ gerbang (§3.3) |
| 12 | **Cabang ganda `IsMBU`** — `"MBUCar"`/`"MBUMotorCycle"` masing-masing 2× | E-2 §6.2 | redundan, hasil tak berubah |

---

## 6. Kesiapan spec — apa yang bisa dikerjakan sekarang

### 6.1 ✅ BISA dispec sekarang (tidak menunggu pihak luar)

| Area | Dasar | Catatan |
| --- | --- | --- |
| **Before-image tiga lapis** | E-3 §1–§4 | Lapis A (query + 54 salinan), lapis B (`SetOldData`, 52 properti `*Old`, peta per lini), lapis C (7 varian). ⚠️ Validitasnya bergantung **I-1** untuk langkah 4 `SetOldData` |
| **Rumus SELISIH** | E-3 §6, E-5 §3.3 | `SELISIH = (nilai_baru × ProrateEDMEnd) − nilai_lama_TreatyType_sama`; 57 penugasan pengurangan terverifikasi |
| **Prorata EDM** | E-3 §1.4 | `@Math.divide(…, 20)` + guard pembagian-nol; **Ratio**, bukan Money |
| **Tabel rumus `EDMPremiMenjadi`** | E-3 §5.2 | 22 penugasan, 11 rumus, per (lini × `EdmType` × jalur) |
| **Alur endorsement + gerbang masuk** | E-2 §1.2 | `Start2 → Assignment7`, 6 properti, workbasket `ReasFacInMarketing`, tiket `AdminPolicy` |
| **Predikat `When` EDM** | E-1 §5–§7 | 56 rule `When` berbeda; P1 terbantah, P2 terbatas 6 COB. ⚠️ 6 predikat COB menunggu **E-Q5** |
| **Klasifikasi & pemanggilan 354 EDM-only** | E-4 | per `pxObjClass`, kelompok fungsi, peta pemanggilan |

### 6.2 ⛔ MENUNGGU pihak luar

| Area | Menunggu |
| --- | --- |
| **Jalur produksi endorsement** (`facinproduction`, `FACOUTPRODUCTION`, `JSON_POLIS`) | **D-1** (4 stored procedure) · **D-2** (skema) · **D-3** (`PRODKE`) · tabel flat + `ALL_SOURCE` |
| **Gerbang penolakan pembuatan case EDM** | **D-4** (`OPENPROTEKSI_EDM`) |
| **Tangga akseptasi EDM** | **E-Q10** (`IsUWAccepted` belum dibaca) · **E-Q6** (marketing sebagai UW) · **I-2** (`SetToInbox_ACT`/`SetBanding_ACT`) · **E-Q28** (2 rule SQL akseptasi hanya di `DDL\`) |
| **Validitas seluruh pemetaan alur** | **I-1** — 617 prakondisi nonaktif |
| **Lingkup** | **W-1** (Life) · **W-2** (pilih-tertanggung) · **W-3** (Bonding) · **W-4** (Travel) |

### 6.3 Urutan kerja yang disarankan pasca-discovery

1. **Kirim paket ke DBA (D-1…D-4) dan IT (I-1, I-2)** — jalur terpanjang, mulai lebih dulu.
2. **Putuskan W-1…W-4** — menentukan lingkup sebelum spec ditulis; W-1 (Life) berdampak paling besar.
3. **Spec bagian §6.1 yang tidak memblokir** — before-image, SELISIH, prorata, alur masuk.
4. **Baca 12 berkas tipe yang masih 100 % berbeda** (`DecisionTable` 10, `Flow` 1, `DecisionTree` 1) —
   pekerjaan internal, tidak menunggu siapa pun, dan menutup **E-Q10**.
5. Spec jalur produksi **setelah** DBA menjawab.

⚠️ **Tidak ada spec yang ditulis di putaran ini.** `/to-spec` dan `/to-tickets` tidak dijalankan —
ini penutup discovery, bukan spec.

---

## 7. Yang tetap di luar lingkup

- Korpus Treaty `RNM_BRD\` — **tidak dibaca** (K-005)
- Isi `SaveFacinProdEDMBonding_Act` dan 3 rule akseptasi di `DDL\` — menunggu **E-Q28/W-3**
- Jalur produksi NB sebagai pembanding — menunggu DBA
- Arsip §4 (7 perbedaan gerbang akseptasi EDM) termasuk klaim *"nilai dasar akseptasi = SELISIH TSI"*
  — **belum diuji**; masuk pekerjaan pasca-discovery butir 4
- Arsip §10 **#12–#23** — belum ditinjau

---

*Discovery EDM ditutup. Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
