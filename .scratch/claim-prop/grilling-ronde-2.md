# Grilling — Claim Prop, Ronde 2

**Tanggal:** 2026-09-17 · **Modul:** `Claim Prop` (270 berkas) · **Status:** frontier fakta berjalan

> Korpus `D:\XML\RNM_BRD\` **READ-ONLY**. Tulisan hanya ke `OUTPUT_HASIL_RNM\`.
> ⛔ `D:\XML\nusantara-re\` di-blacklist.
> Identitas rule = `class/nama` dari `<pxInsName>`; tipe dari field pertama `<pzOriginalInstanceKey>`.
> ⚠️ **Bukti WAJIB memakai nomor step Pega** (`<pyStepPageReference>`), **bukan** nomor baris XML.
> Tanda: `[terverifikasi]` `[keputusan work owner]` `[data DBA]` `[terbuka]` `[dugaan]`.
> ⚠️ Tidak ada PII di berkas ini.

---

## 0. Aturan baca yang berlaku sejak ronde ini

Ronde 1 ditulis **tanpa** ketiga aturan ini. Itulah sebab banyak temuan `[terverifikasi]` di sana
yang berbunyi "digerbangi X" **terbalik atau tidak berlaku**.

### Aturan 1 — flag prakondisi

Teks syarat tersimpan di `pyStepsPreCondParamsWhen`, tetapi **hanya berlaku bila**
`pyStepsPreCondition = true`.

| `pyStepsPreCondition` | Arti |
| --- | --- |
| `true` | gerbang **berlaku** |
| `false` | teks syarat **TIDAK berlaku** — step **jalan tanpa saringan** |
| kosong | step memang tidak punya gerbang |

`[terverifikasi]` Dikonfirmasi work owner dari Pega Designer: step ber-flag `false` tampil **tanpa
when** di Designer dan selalu jalan.

`[terverifikasi]` **Sensus Claim Prop:** 535 step punya teks syarat — **451 aktif, 84 tidak berlaku
(15 %)**.

### Aturan 2 — blok di-remark

`pyStepsBlockName = //` → step **di-remark: tidak pernah jalan**, apa pun gerbangnya.

`[terverifikasi]` **Sensus:** 95 step di 43 rule se-korpus; **68 di antaranya di Claim Prop**.
Termasuk `TryMakePLA_Act`, `PrintDLATreatyIn`, `GetPayAttachmentAdj_Act`, `CountSpreading_Act`,
`GetMasterTreaty_Act`, `SetPayableTo_act`.

### Aturan 3 — kode arah pada tiap baris syarat

| `WhenTrue` / `WhenFalse` | Arah |
| --- | --- |
| BENAR=2, SALAH=3 | **normal** — step jalan bila syarat **terpenuhi** |
| BENAR=3, SALAH=2 | **TERBALIK** — step jalan bila syarat **TIDAK terpenuhi** |

Kode: `2` lanjut · `3` lewati step · `5` keluar iterasi `[terbuka — belum terbukti]` ·
`6` keluar activity `[terverifikasi]` dari keterangan step `CloseClaimProp` step 4
*"-- Kalau kena proteksi exit act"*.

`[terverifikasi]` **Sensus baris syarat pada gerbang AKTIF:** 380 normal · **44 TERBALIK** ·
89 memakai kode 5/6.

Rule dengan baris terbalik: `AttachmentProtect_ACT` (4) · `GetBase64Attachment` (3) ·
`PrintDLATreatyIn` (3) · `SaveOutstanding_Act` (3) · `TryMakePLA_Act` (3) · `AddAdjustment_Act` (2) ·
`GetUrlGoogleStorage_Act` (2) · `SendEmailKlaimRejectClose` (2) · `PrintFileAcceptance` (1) ·
`SetMOClaimTreaty_Act` (1) · `ProteksiData_act` (1) · `CloseClaimProp` (1).

**Contoh pembuktian aturan 3** — `SetDefNonCatastrope_Act`:

```
step 1   StsKatastrofe=="Catastrophe"   BENAR=2 SALAH=3  → jalan bila Catastrophe
step 2   StsKatastrofe=="Catastrophe"   BENAR=3 SALAH=2  → jalan bila BUKAN Catastrophe
```

Dibaca polos keduanya tampak bertentangan; dengan aturan 3 logikanya jelas.

⚠️ **Aturan baca untuk Section belum ditetapkan.** Gerbang di Section memakai *visible-when*, dan
mekanismenya **berbeda** dari Activity. Penetapannya sedang dikerjakan dan akan ditulis di §2
**sebelum** kesimpulan apa pun tentang gerbang UI diambil.

### Aturan penomoran — berlaku seterusnya

**Setiap identitas yang dirujuk — `P<n>`, `Q<n>`, `OQ-<n>`, `R<n>`, `AC <n>` — WAJIB punya definisi
utuh di dalam berkas yang sama, atau rujukan eksplisit ke berkas dan bab tempat definisinya berada.**

**Menyebut nomor tanpa definisi = cacat, sama beratnya dengan klaim tanpa bukti.**

Peta definisi berkas ini: **`P<n>`** → §19 · **`Q<n>`** → §21 · **`R<n>`** → §13 ·
**`OQ-<n>`** → `discovery/open-questions.md` (di luar berkas ini).

> ### ⚠️ RALAT 2026-09-17 — aturan ini lahir dari cacat nyata di berkas ini
>
> `Q1`–`Q5` disebut **tujuh kali** di berkas ini (§10c, §14a, §16c, §17b, §18 ×2, §19) **tanpa satu
> pun definisi**. Work owner karena itu **tidak dapat menjawab apa pun**. Diperbaiki di **§21**.
>
> Cacat yang sama ada di `grilling-ronde-1.md`: kalimat penutupnya menyebut `Q1–Q13` dan mengarahkan
> ke §21 berkas itu, padahal §21 di sana berisi **`P1`–`P14`** — **nol definisi `Q`**. Berkas ronde 1
> **tidak diubah**; cacatnya dilaporkan saja (§21 berkas ini, bagian pemeriksaan).

### Aturan arti kode — berlaku seterusnya

**Keterangan step (`pyStepsDescription`) menjelaskan APA YANG DIKERJAKAN step, bukan ARTI KODE yang
digerbanginya.**

**Arti sebuah kode hanya boleh disimpulkan dari data nyata (`[data DBA]`) atau dari tabel master —
tidak pernah dari keterangan step.**

Berkas ini sudah mencatat **lebih dari 20 nama rule dan kolom yang berbohong** (§8g, §14e, §16d).
**Keterangan step tunduk pada kecurigaan yang sama** — ia ditulis pengembang untuk dirinya sendiri,
tidak pernah diverifikasi, dan tidak pernah ikut diuji.

> ### ⚠️ RALAT 2026-09-17 — aturan ini lahir dari kesalahan nyata di berkas ini
>
> `TreatyType 10004` dicatat sebagai **"Fac Out"** di §2 (P14) semata-mata karena keterangan step
> `Activity/PrintDLATreatyIn.xml` step **15.11.1–3** berbunyi *"set fac out payment type 1, 2, 5"*.
> **Data produksi membantahnya**: `10004` = **`QS (R/I)`**. Ralat lengkap di §2 dan §10e.

### Aturan silsilah — berlaku seterusnya

**`pzOriginalInstanceKey` menyimpan asal *save-as*, BUKAN pemanggil.** Pencarian nama rule di dalam
berkas rule lain **akan menangkap field ini** dan tampak seperti rujukan.

**Sebelum menyatakan rule A merujuk rule B, pastikan kemunculannya bukan di `pzOriginalInstanceKey`,
`pzInsKey`, `pzDocumentKey`, atau `pzIndexOwnerKey`.**

`[terverifikasi]` Dua kasus nyata yang sempat salah dibaca sebagai pemanggilan:
`RDBList/GenerateNoCLMTreatyIn.xml` → `…GCNM!GENERATENOPLATREATYIN` dan
`RDBList/GetTreatyInMasterProp_SQL.xml` → `…ASM!GETTREATYINMASTER_SQL` (§26a).

### Aturan sensus nama rule — **wajib tidak peka huruf besar-kecil**

Pega menyimpan `pxInsName` dalam **huruf besar semua**, sementara **nama berkas mengikuti ejaan
penulisnya** dan rujukan di dalam rule mengikuti ejaan pengetiknya.

`[terverifikasi]` Kasus nyata: `When/IsAneka.xml` merujuk `IsAllRisk` dengan ejaan **`isAllRisk`**.
Sensus peka huruf besar-kecil melewatkannya dan **salah melabelinya yatim** (§27e).

### Aturan konflik cek versus aturan berdiri

**Bila cek numerik dalam sebuah perintah bertabrakan dengan aturan berdiri, aturan berdiri menang,
dan hasil ceknya dilaporkan apa adanya beserta alasannya.**

`[terverifikasi]` Kasus nyata: cek T7 mensyaratkan `grep -c "menunggu jawaban work owner"` → `0`,
sedangkan aturan berdiri melarang menghapus teks lama. Baris lama **dicoret, tidak dihapus**, cek
mengembalikan **1**, dan selisih itu dilaporkan — bukan disembunyikan.

---

## 1. Ralat Ronde 1 — gerbang yang ternyata tidak berlaku

`[terverifikasi]` Seluruh baris berikut **gerbangnya tidak berlaku**, jadi step-nya **jalan tanpa
saringan** — kebalikan dari yang tertulis di Ronde 1.

| Rule | Step | Syarat tersimpan yang TIDAK berlaku | Akibat sebenarnya |
| --- | --- | --- | --- |
| `CountSpreadingADJ_Act` | 3.1 | `.AcceptanceStatus != 2` | baris **reject ikut diproses** |
| | 6 | `paymenttype=="1"\|\|"2"\|\|"5"` | loop jalan untuk **semua** payment type |
| | 6.4 | tiga kode treaty | `SpreadForBreak` **ditimpa tiap putaran** → berisi baris **TERAKHIR** |
| | 7 | `CountSpreadQS != 0` | jalan **walau daftar QS kosong** |
| `InsertGoogleStorage_Act` | 6 | `IsPEGAPROD` | **jalan di luar produksi juga** |
| `HitServiceToKasir_Act` | 13.1.3 · 13.2.2.1.1.3 | `pyWorkIDPrefix=="CLMNP-"` | jalan untuk **semua prefix** |
| `SaveOutstanding_Act` | 39 | `AcceptStatus==1 && KomiteCount==KomiteLoop` | **Arasapas dipanggil tanpa menunggu komite selesai** |
| `SendCloseClaimToKomite` | 5.4 | `ParamCari.CARI2==1` | **child case komite dibuat tanpa syarat** |
| `GetPayAttachmentAdj_Act` | 3.2 | `.IsPayment==1` | REST dipanggil **walau bukan payment** |
| `GetDtlPaymentPremi_act` | 6 | satu nomor polis hardcode | REST jalan **untuk semua polis** |
| `PrintFileAcceptance` | 10 | `PaymentType!="4" \|\| !="3"` | gerbang **tautologi**, dan memang tidak berlaku |

⚠️ **`SendCloseClaimToKomite` step 5.4 menjelaskan §2d Ronde 1**: dua jalur komite bisa aktif
bersamaan **bukan** karena kelalaian desain yang rumit, melainkan karena **syaratnya memang tidak
dipasang**.

### 1a. Empat rule "hilang" ternyata dipanggil dari step DI-REMARK

`[terverifikasi]` Ronde 1 §17 mencatatnya sebagai "tidak ada di korpus" dan mengusulkan meminta ke
pemilik export. **Tidak perlu** — pemanggilnya mati:

| Rule yang dirujuk | Pemanggil | Step |
| --- | --- | --- |
| `GetBreakDownTreaty` | `SetTreatyNameSpreading_Act` | 10 |
| `InsertCatastrope_SQL` | `SaveCatasrtope_Act` | 2 |
| `SaveMasterAdjusterConsultant_sql` | `SaveAdjusterConsultant_Act` | 5 |
| `SetProtectionEstimation` | `ProteksiInitialandDate_Act` | 6 |

⚠️ Konsekuensi yang belum tuntas: bila step penyimpan master **di-remark**, **bagaimana master
adjuster/consultant dan master katastrofa sebenarnya tersimpan?** Sedang ditelusuri.

---

## 2. Pemblokir yang SUDAH TERTUTUP — tidak ditanyakan lagi

| # | Jawaban | Status |
| --- | --- | --- |
| P1 | Rumus deductible yang berlaku = **`CountListClaimAmountIDR`** (versi 2024): `.NetDeductibleValue = ClaimData.NetDeductibleValue`; `.Value = (share/100 × .ClaimAmount) − .NetDeductibleValue`. Rumus `AddListClaimAmount` (2022) **ditinggalkan**. | `[keputusan work owner]` |
| P3 | `.AcceptanceStatus` **nol penulisan** di seluruh Activity keempat modul. Penulisnya **di luar export**. | `[terverifikasi]` |
| P4 | Arti `PaymentType` 1–6 dan `.Type` sudah diberikan work owner; pemetaannya ada di catatan work owner. **Jangan tebak.** | `[keputusan work owner]` |
| P5 | `STSREJECT`: `0` belum diputus · `1` aksep · `2` reject · `4` saat `Adjustment.PaymentType` final ATAU klaim di-close. | korpus + work owner |
| P6 | 10 procedure + 2 fungsi + 1 sequence — lihat §3. | `[terverifikasi]` |
| P7 | `GENERATE_NOCLMTRTYINTEMP(KODE, TAHUN, TIPE)` versus `GENERATE_NOCLMTREATYIN(KODE, KODE_BIS, BULAN, TAHUN, TIPE)`. Nomor temp **tidak punya bulan dan kode bisnis**, jadi **tidak bisa dipromosikan** jadi final. Keduanya COMMIT sendiri. | `[terverifikasi]` |
| P8 | `TRT` = Treaty · `TNP` = Treaty Non Prop. **Bukan** kebocoran lintas modul — modul ini memang melayani tiga prefix, jadi punya satu rule simpan per jenis. | `[terverifikasi]` |
| P9 | **DLA = Definite Loss Advise.** | `[keputusan work owner]` |
| P10 | Huruf `S` pada `CLMS-`/`CLMPS-`/`CLMNPS-` = **server Syariah**. | `[keputusan work owner]` |
| P11 | Class `ASM-FW-GISFW-Int-ADJUSTERCONSULTANT`; kolom `ID NAME ADDRESS TELPNO USERNAME EDITDATE`; ID dari `ADJUSTERCONSULTANT_SEQ` + prefix `POOLDATA.M_SITE_DATABASE WHERE CURRENT_SITE='1'`. | `[terverifikasi]` |
| P12 | `StsKatastrofe = {Catastrophe, Non-Catastrophe}` · `NonKatastrofeType = {Big Claim, Claim}` · `@GCNM.GetPageJSONString()` / `@ASM.GetPageJSONString()` = pembungkus `ClipboardPage.getJSON(false)` — **dump seluruh step page**. | `[keputusan work owner]` |
| P13 | Kode enum arah — lihat §0 aturan 3. | `[terverifikasi]` |
| P14 | `InputData.CARI27` = kode jenis usaha, menentukan template teks Insured Interest: `10007` Property · `10016` Motor Vehicle · `10009` Marine Cargo · 14 kode lain Aneka. `TreatyType 10004` = **Fac Out** (keterangan step `PrintDLATreatyIn` 15.11.1–3). **Dua tabel kode BERBEDA** meski angka `10004` muncul di keduanya — jangan digabung jadi satu enum. Keputusan: **ikuti apa adanya**. | `[keputusan work owner]` |

⚠️ **Cacat yang ikut terbawa:** `10009` Marine Cargo memakai template **persis sama** dengan Aneka
(`"Project Name: Risk Location: Construction Period: Project Type:"`) — isinya field proyek
konstruksi, bukan kargo laut. Tampak **salah salin di Pega**. Keputusan "ikuti apa adanya" berarti
bug ini **ikut dimigrasikan**; bila kelak diperbaiki, catat sebagai **penyimpangan sadar**.
Pertanyaannya terdaftar sebagai **P17**.

> ### ⚠️ RALAT 2026-09-17 — **`TreatyType 10004` BUKAN "Fac Out"**
>
> **Sumber bukti baru:** `[data DBA]` — **satu baris contoh** data produksi dari
> `os_akseptasi_klaim.DATA_JSON` dan `json_klaim.DATA_JSON`, klaim
> `RNM-K22.09.2026.T00718` / `CLMP-7134`. **Bukan korpus. Bukan `[terverifikasi]`.**
> Satu klaim — **jangan digeneralisasi**.
>
> **Yang salah:** `TreatyType 10004` dicatat sebagai **"Fac Out"**, bersumber **hanya** dari
> keterangan step `Activity/PrintDLATreatyIn.xml` step **15.11.1–3** (*"set fac out payment type
> 1, 2, 5"*).
>
> **Yang benar menurut data** `[data DBA]`:
>
> | Daftar | `TreatyType` | `TreatyName` | share |
> | --- | --- | --- | ---: |
> | `SpreadingBreakQS[0]` | `10028` | `QS (OR)` | 40 |
> | `SpreadingBreakQS[1]` | **`10004`** | **`QS (R/I)`** | 60 |
> | `SpreadingClaim[0]` | `10236` | `2023 QS 150M TRT` | 100 |
> | `SpreadingRisk[0]` | `10035` | `QUOTA SHARE` | 100 |
>
> **`10004` = `QS (R/I)`, bukan Fac Out.** Keterangan step itu menjelaskan **apa yang dikerjakan
> step**, bukan arti kode yang digerbanginya — persis pola yang kini ditetapkan sebagai
> **§0 Aturan arti kode**.
>
> ### ✅ P14 TERKONFIRMASI dari data — dan `CARI27` ternyata **`TreatyGroupID`**
>
> `[data DBA]` Satu baris contoh yang sama memuat:
>
> ```
> TreatyGroupID    "10007"
> TreatyGroupName  "PROPERTY"
> InsuredInterest  "Constrution Class : Occupation:  Risk Category: Risk Location: Coverage :"
> ```
>
> Teks itu **persis** template yang di-set `Activity/SetValueToClaim_Act.xml` step **9** untuk
> `CARI27 == "10007"` — **termasuk typo `Constrution`**.
>
> **Status pemetaan `CARI27` naik dari dugaan menjadi `[terverifikasi]` + `[data DBA]`**, dan
> `InputData.CARI27` kini diketahui sebagai **`TreatyGroupID`** — **bukan kode tak dikenal**.
>
> ⚠️ **Tetap berlaku, dan makin penting:** `CARI27 = 10004` (dalam keranjang **Aneka**) dan
> `TreatyType = 10004` (**QS (R/I)**) adalah **dua tabel kode yang berbeda**. **Jangan digabung.**
> Kebetulan angkanya sama justru yang membuat kesalahan ini mudah terulang.
>
> ### ✅ P17 TERTUTUP
>
> Pertanyaan "apakah `10009` Marine Cargo memang seharusnya memakai template Aneka" **tercakup oleh
> keputusan "ikuti apa adanya"** pada P14. Cacat salin-template **ikut dimigrasikan**; bila kelak
> diperbaiki, tetap dicatat sebagai **penyimpangan sadar**.

---

## 3. Fakta fondasi Ronde 2

### 3a. Komposisi korpus — tujuh jenis rule belum pernah disentuh

`[terverifikasi]` 270 berkas:

```
Activity 112 · RDBList 54 · Section 36 · ReportDefinition 19 · FlowAction 12
DataTransform 11 · Harness 11 · ConnectREST 6 · When 6 · Flow 1
DecisionTable 1 · SystemSettings 1
```

⚠️ **Ronde 1 hanya menyisir Activity.** Sisanya (158 berkas, 59 %) belum pernah dibaca. Penyisiran
sedang berjalan.

### 3b. Tiga schema Oracle, bukan satu

`[terverifikasi]`

| Schema | Dipakai oleh | Objek |
| --- | --- | --- |
| `pooldata` | mayoritas | — |
| `arasapas` | `CekLunasPremi_Sql` | `arasapas.invoice`, `arasapas.detail_invoice` |
| `reinsurance` | `getStatusKonversi_SQL` | `reinsurance.trloss_detail_t` |

### 3c. Inventaris stored procedure lengkap — **P6 tertutup**

`[terverifikasi]`

```
POOLDATA.GENERATE_NOCLMTREATYIN      (KODE, KODE_BIS, BULAN, TAHUN, TIPE → START_DATE, TSI)
POOLDATA.GENERATE_NOCLMTRTYINTEMP    (KODE, TAHUN, TIPE                  → START_DATE, TSI)
POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER
POOLDATA.GET_TOKEN_STORAGE
POOLDATA.PEGA_JSON_KLAIM_PNC
POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM       (11 param)
POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT    ( 7 param)
POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP    (11 param)
POOLDATA.PEGA_D_CAUSE_OF_LOSS
POOLDATA.PEGA_M_CAUSE_OF_LOSS

fungsi   : pooldata.getcurrencystandard · gl.f_get_email
sequence : ADJUSTERCONSULTANT_SEQ
```

### 3d. ⚠️ Nama orang di-hardcode di dalam SQL

`[terverifikasi]` `RDBList/GetLimitDirekturUtama_SQL.xml`:

```sql
select LIMIT_BOTTOM from POOLDATA.EMAILKOMITE where name = '<satu nama orang>'
```

Limit wewenang Direktur Utama dicari lewat **nama literal**. **Bila orangnya berganti, query ini
diam-diam mengembalikan kosong** — dan karena `LewatLimit` lalu tidak pernah bernilai 1, kewajiban
dokumen ADU hilang tanpa pesan apa pun. Belum tercatat di Ronde 1. (Nilai tidak disalin ke sini.)

### 3e. `.ClaimSpreaded` — 7 activity, 22 tempat, **4 perlakuan pembulatan**

`[terverifikasi]`

| Bentuk | Tempat |
| --- | --- |
| `/100` polos | `CountSpreadingADJ` 6.2, 7.2 · `CopyOldataCurr` 7.1, 9.1, 10.1 · `CountPersen` 5.2.1 · `CountTotalInsterest` 12.2.1, 16.2.1 |
| `@divide(...,100,4)` | `CountSpreadingADJ` 9.2.1, 9.3.1 · `CountTotalInsterest` 17.2.1 |
| `@divide(...,100,10)` | `CountPersen` 6.2.1 |
| `@divide(...,100,20)` | `CountGrossAdjTreaty` 5.2 |

⚠️ `CountTotalInsterest_Act` **sendirian memakai tiga cara berbeda**.

**Pembulatan bertingkat** — empat tempat memakai hasil berpembulatan sebagai masukan:

| Tempat | Basis |
| --- | --- |
| `CountSpreadingADJ` 7.2 | `SpreadForBreak` ← `.ClaimSpreaded` step 6.4 |
| `CountSpreading` 10.1 | `SpreadingClaim(...).ClaimSpreaded` |
| `CountSpreading` 11.2.1 | `SpreadingClaim(...).ClaimSpreaded` |
| `CopyOldataCurr` 10.1 | `Local.ClaimSpreaded` |

`[keputusan work owner]` `.ClaimSpreaded` adalah **detail/turunan, bukan master** — boleh dihitung
ulang dari `AdjustmentValue × SharePercentage`.

### 3f. ✅ **P16 TERTUTUP** — bentuk JSON lama terbukti dari data

**Sumber:** `[data DBA]` — **satu baris contoh** dari `os_akseptasi_klaim.DATA_JSON` dan
`json_klaim.DATA_JSON`, klaim `RNM-K22.09.2026.T00718` / `CLMP-7134`. **Bukan korpus.**

`@GCNM.GetPageJSONString()` = `stepPage.getJSON(false)`. Perilaku argumen `false` kini terbaca:

| Properti sistem | Hadir di dump? | Catatan |
| --- | --- | --- |
| `pxObjClass` | **ada** | di **setiap** halaman bersarang, termasuk **tiap baris page-list** |
| `pxListSubscript` | **ada** | di tiap baris page-list |
| `pxCreateOperator` | **ada** | di `os_akseptasi` |
| `pyExpanded` | **ada** — nilai `"true"` | ⚠️ **status buka-tutup grid di LAYAR** |
| `pzInsKey` · `pxCreateDateTime` · `pyID` | **tidak ada** | — |

**Kesimpulan `[data DBA]`:** `getJSON(false)` menuliskan **properti yang kebetulan ada di halaman,
apa pun awalannya**. Ia **bukan** daftar metadata tetap yang dapat diprediksi.

#### Konsekuensi mengikat untuk migrasi

1. **Parser Go membuang seluruh kunci berawalan `px`, `py`, `pz`** — dan **tidak boleh mengandalkan
   kunci tertentu selalu hadir**. Kehadiran `pxObjClass` di satu dokumen tidak menjamin kehadirannya
   di dokumen lain.
2. ⚠️ **Status tampilan ikut tersimpan ke database** (`pyExpanded` = status buka-tutup grid). Ini
   **kebocoran lapisan UI ke lapisan penyimpanan** — dicatat sebagai **cacat**, dan **tidak boleh
   ditiru**.
3. **Semua nilai bertipe string**, termasuk angka dan tanggal. Tanggal memakai **dua format yang
   berdampingan dalam satu dokumen**: `YYYYMMDD` (`"20250906"`) dan
   `YYYYMMDDTHHMMSS.mmm GMT` (`"20260902T170000.000 GMT"`).

⚠️ Butir 3 **senada dengan P22** (`RDBList/Update_T_Storage_SQL.xml`, dua format tanggal berlawanan
dalam satu statement) — pola dua-format bukan kejadian tunggal.

---

## 4. Pemblokir yang masih terbuka

| # | Pertanyaan | Pemilik | Status |
| --- | --- | --- | --- |
| P2 | `.ClaimSpreaded`: basis pengali mana yang benar — nilai satu baris adjustment (step 6.2) atau nilai baris `Hepeng.pxResults` (step 9.2.1)? Dan apakah pembulatan 4 desimal hanya di step 9 disengaja? | Product + UW | terbuka |
| P15 | **BARU.** Jumlah seluruh baris spreading wajib **sama persis** dengan nilai adjustment-nya, atau boleh ada selisih pembulatan? Pega **tidak punya langkah koreksi selisih**. | Finance | terbuka |
| P16 | **BARU.** Argumen `false` pada `stepPage.getJSON(false)` — apakah properti sistem (`px*`, `py*`) ikut terserialisasi? Menentukan apakah JSON lama di `json_klaim` / `os_akseptasi_klaim` berisi metadata Pega atau data murni. **Cukup satu baris dump JSON untuk menjawabnya.** | DBA | terbuka |
| P17 | **BARU.** Apakah `10009` Marine Cargo memang seharusnya memakai template Insured Interest yang sama dengan Aneka? | Product + UW | terbuka |

---

## 5. Penyempurnaan aturan baca — dua temuan dari audit

### 5a. ⚠️ Aturan 2 perlu dipertajam: `pyStepsBlockName` juga dipakai sebagai **label blok**

`[terverifikasi]` Nilai `pyStepsBlockName` di korpus ini **tidak selalu** berarti remark. Ia juga
dipakai sebagai **nama label blok** untuk lompatan: `FAIL`, `EXIT`, `END`, `TO`, `UP`.

**Hanya nilai `//` yang berarti di-remark.** Contoh: `SaveAdjusterConsultant_Act` step 8 ber-
`pyStepsBlockName = FAIL` — itu **label tujuan lompatan** dari step 4, **bukan** step mati.

### 5b. ⚠️ Aturan 3 punya kategori keempat: **kode arah KOSONG sebelah**

`[terverifikasi]` Selain normal (2/3) dan terbalik (3/2), ada baris syarat yang salah satu sisinya
**kosong**:

| Rule | Step | Syarat | Kode |
| --- | --- | --- | --- |
| `ProteksiData_act` | 3 | `@contains(TreatyInMaster.StatusAkseptasi,"Resolve Complete")` | BENAR=3, SALAH=**kosong** |
| `ProteksiData_act` | 6 | `Local.SizeEstimation>0` | BENAR=**kosong**, SALAH=3 |
| `CountSpreading_Act` | 6.2.3 | `@greaterThan(Local.TotalPersen,100)` | BENAR=**kosong**, SALAH=3 |

⚠️ Perilakunya **bergantung default engine Pega**, dan default itu **tidak terbaca dari korpus**
→ `[terbuka]`. Pada `ProteksiData_act` step 3 akibatnya pesan *"ID MASTER Sedang dalam tahap Revisi,
Tidak bisa Claim"* muncul justru saat status **BUKAN** "Resolve Complete" — berlawanan dengan
keterangan step-nya, meski hasil akhirnya kebetulan setara.

---

## 6. Ralat §14a dan §16 Ronde 1 — hasil audit

### 6a. ⚠️⚠️ `CountSpreading_Act` — **tidak ada percabangan sama sekali**

**KOREKSI besar terhadap Ronde 1 §14a.** Ronde 1 mencatat dua cabang rumus (`SizeEst==1` versus
`SizeEst>1`). `[terverifikasi]` **Keduanya tidak ada dalam perilaku runtime:**

| Step | Status | Akibat |
| --- | --- | --- |
| 7 | **DI-REMARK** (`//`) | cabang `SizeEst==1` **kode mati** |
| 10 | **DI-REMARK** (`//`) | loop BreakQS cabang `SizeEst==1` mati; anak 10.1 ikut mati |
| 6.3 | **DI-REMARK** (`//`) | duplikat pesan >100 % |
| 8, 9, 11, 11.2 | `pyStepsPreCondition = false` | gerbang `SizeEst>1` **MATI** → **selalu jalan** |
| **9.2.2** | `PRE=true`, 2/3 normal | **satu-satunya rumus `ClaimSpreaded` induk yang berjalan** |
| **11.2.1** | `PRE=true`, 2/3 normal | **satu-satunya rumus breakdown-QS yang berjalan** |

**Artinya sistem SELALU memakai jalur "multi-kurs"**, bahkan ketika hanya ada satu kurs estimasi.
Percabangan yang tampak di korpus adalah ilusi.

⚠️ `[terverifikasi]` Step 12 (`pyWorkPage.IsEditEstimation=="false"`) gerbangnya **mati** →
`DataChronology.CARI1` **selalu** diisi `"Add % spreading Claim"`, lalu step 13 (aktif) menimpanya
jadi `"Edit % spreading Claim"` bila `IsOutstanding==1`. **Jejak audit jadi tidak dapat dipercaya.**

### 6b. `CountSpreadingADJ_Act` — empat penulis, dan satu bug nyata

`[terverifikasi]` **Nol step di-remark** di rule ini.

| Step | Ekspresi persis | Basis | Status |
| --- | --- | --- | --- |
| 6.2 | `Local.Adjustment * Local.Sharepersen / 100` | `.AdjustmentValue` adjustment berjalan | aktif (induk gerbang mati) |
| 7.2 | `Local.SpreadForBreak * Local.Sharepersen / 100` | `SpreadForBreak` dari step 6.4 | aktif (induk gerbang mati) |
| 9.2.1 | `@divide(.SharePercentage,100,4) * Local.TotalAllADJ` | `Hepeng.pxResults` | **AKTIF**, gerbang normal |
| 9.3.1 | `@divide(.SharePercentage,100,4) * Local.TotalAllADJ` | idem | **AKTIF**, gerbang normal |

**Jawaban P2 sebagian — apa itu `Hepeng.pxResults`:** `[terverifikasi]` ia **page kerja clipboard,
bukan hasil query DB**. Class `Code-Pega-List`, elemen `ASM-FW-GCNMFW-Data-Adjustment`. Isinya
**daftar kurs unik + total adjustment per kurs**, disusun step 3 → dedup Java step 4 → akumulasi step
5 → dibaca balik step 9.1.

⚠️ `[terverifikasi]` **Dua cacat dari gerbang mati:**

1. **Step 3.1** (`.AcceptanceStatus != 2`, mati) → baris **reject ikut di-append** ke `Hepeng`.
   Tetapi step 5.2.1 yang menjumlah **aktif dan menyaring** `AcceptanceStatus != 2`. Jadi daftar kurs
   tercemar baris reject, sementara totalnya bersih — **tidak konsisten**.
2. **Step 6.4** (tiga kode treaty, mati) → `Local.SpreadForBreak` **ditimpa setiap iterasi**, sehingga
   saat step 7 berjalan ia berisi `ClaimSpreaded` baris **terakhir**, bukan baris treaty QS. **Basis
   pengali step 7.2 salah secara semantik.**

⚠️ `[terverifikasi]` Total di step 10 (`.TotalSharePersen`, `.TotalSpreadBreakQs`,
`.TotalSpreadAdjustment`) diambil dari akumulator step 6.3/7.3 — yaitu versi **`/100` polos**,
**bukan** versi `@divide(...,100,4)`. Total dan rincian karena itu berbeda presisi.

⚠️ `[terverifikasi]` Tidak ada `Page-Remove Hepeng` di awal activity (bandingkan `CountSpreading_Act`
step 1 yang meng-`Page-Remove Estimate`) → **risiko sisa data antar-pemanggilan**, sebagian diredam
dedup step 4.

### 6c. `SetTreatyNameSpreading_Act` — sumber pengganti memakai indeks hardcode

`[terverifikasi]` **TETAP BERLAKU**: step 10 (`GetBreakDownTreaty`) di-remark, dan itu **satu-satunya**
remark di rule ini.

⚠️ `[terverifikasi]` Penggantinya: step 11/11.1 mengisi `Databreak.pxResults` dari clipboard
`pyWorkPage.TreatyInMaster.Limits(1).Detail(1).SpreadingList` — **indeks `Limits(1).Detail(1)`
di-hardcode**. Bila master treaty punya lebih dari satu limit atau detail, sisanya **diabaikan diam-
diam** → `[terbuka]`.

### 6d. Sensus pembulatan `.ClaimSpreaded` — **20 titik tulis, bukan 22**

**KOREKSI:** angka 22 di §3e menghitung juga 2 titik **baca** (`AddKomiteTreatyChild_ACT` step 7.2;
`SendEmailKlaim` step 8.2.3). Titik **tulis** ada **20**, tiga di antaranya **mati**.

| Bentuk | Titik AKTIF | Di mana |
| --- | ---: | --- |
| `/100` polos | **9** | `CopyOldataCurr` 7.1/9.1/10.1 · `CountSpreadingADJ` 6.2/7.2 · `CountSpreading` 9.2.2/11.2.1 · `CountTotalInsterest` 12.2.1/16.2.1 |
| `@divide(...,100,4)` | **3** | `CountSpreadingADJ` 9.2.1/9.3.1 · `CountTotalInsterest` 17.2.1 |
| `@divide(...,100,10)` | **1** | `CountPersen` 6.2.1 — **nilai final yang bertahan** di `ClaimData.SpreadingRisk` |
| `@divide(...,100,20)` | **1** | `CountGrossAdjTreaty` 5.2 (dua faktor) |
| tanpa pembagian | **3** | `AddEstimation` 2.1 & 4.2.2.1.2 · `CountPersen` 4.4 (page sementara) |
| **mati** | **3** | `CountSpreading` 7 & 10.1 · `CountPersen` 5.2.1 (induk step 5 di-remark) |

⚠️⚠️ `[terverifikasi]` **Nilai `ClaimSpreaded` pada baris yang sama ditimpa berulang oleh rule
berbeda dengan presisi berbeda** — mis. `CountSpreadingADJ_Act` step 6.2 (`/100`) lalu step 9.2.1
(`,100,4`); `CountPersen_act` step 5.2.1 (mati) lalu 6.2.1 (`,100,10`). **Nilai akhir bergantung pada
urutan pemanggilan rule, bukan pada kebijakan pembulatan yang disengaja.**

### 6e. ⚠️ Tidak ada koreksi selisih pembulatan — **di mana pun**

`[terverifikasi]` Dicari dengan tiga cara dan nihil semua: (1) regex `selisih|pembulatan|rounding|
@round|sisa|balancing|adjustLast|LastRow|difference` → **0 kecocokan**; (2) ke-20 titik tulis
`.ClaimSpreaded` seluruhnya berbentuk `share × basis / 100` murni atau salin/akumulasi — **nol**
ekspresi koreksi; (3) sumber Java di lima tempat hanya dedup HashSet.

**Selisih pembulatan dibiarkan mengambang.** Ini menguatkan **P15** (Finance).

> ### ⚠️ RALAT 2026-09-17 — **P15 menyempit: penjumlahan ternyata TEPAT**
>
> **Sumber:** `[data DBA]` — **satu baris contoh**, klaim `RNM-K22.09.2026.T00718` / `CLMP-7134`.
>
> ```
> SpreadingClaim[0].ClaimSpreaded       19999999.999992      share 100
>   SpreadingBreakQS[0].ClaimSpreaded    7999999.9999968     share  40
>   SpreadingBreakQS[1].ClaimSpreaded   11999999.9999952     share  60
>                                       ────────────────
>                                       19999999.999992      ← pas, sampai digit terakhir
> ```
>
> **Jumlahnya tepat** — karena share berjumlah **100** dan **tidak ada pembulatan** yang diterapkan.
>
> `[keputusan work owner]` **Data disimpan penuh tanpa pembulatan; pembulatan hanya untuk tampilan.**
>
> **Maka alokasi sisa pembulatan TIDAK diperlukan** — dengan satu syarat: **total share selalu 100 %**.
> Temuan §6e (nol mekanisme koreksi) **tetap berlaku**, tetapi berubah arti: ketiadaan koreksi bukan
> cacat, melainkan **konsekuensi wajar** dari menyimpan nilai penuh.
>
> **P15 menyempit menjadi satu pertanyaan ke Finance:**
> **apakah total share reinsurer wajib tepat 100 %, atau boleh kurang?**
>
> ⚠️ `[terverifikasi]` **Bukti bahwa syarat itu TIDAK dijaga:** `Activity/CountSpreading_Act.xml`
> step **6.2.3** — gerbang **AKTIF**, kondisi `@greaterThan(Local.TotalPersen, 100)`. Ini
> **satu-satunya pemeriksaan total share di 270 dari 270 berkas**, dan ia **hanya menangkap lebih
> dari 100** — **tidak pernah kurang dari 100**. **Penjaga sebelah.**
>
> Bila Finance menjawab "wajib tepat 100 %", maka sistem baru perlu penjaga **dua sisi**, dan itu
> **penyimpangan sadar** dari perilaku lama.

### 6f. Master tetap tersimpan — lewat `Obj-Save`, bukan SQL yang di-remark

**Ini menjawab kekhawatiran §1a.** `[terverifikasi]`

| Master | Step SQL di-remark | **Jalur simpan sebenarnya** |
| --- | --- | --- |
| Adjuster/Consultant | step 5 (`SaveMasterAdjusterConsultant_sql`) | **`Obj-Save` step 4** → class `ASM-FW-GISFW-Int-ADJUSTERCONSULTANT`, page `MstConsultant`, `WriteNow=true` |
| Katastrofa | step 2 (`InsertCatastrope_SQL`) | **`Obj-Save` step 3** → class `ASM-FW-GCNMFW-Int-CATASTROPHE`, page `ParamCat`, `WriteNow=true` |

Keduanya **tanpa `Commit` eksplisit** → commit bergantung pemanggil.

⚠️ `[terverifikasi]` `SaveAdjusterConsultant_Act` step **6 juga di-remark** → `OutputParam.DATASHOW`
**tidak pernah di-set ke 2**, dan `MstConsultant.ID` tidak pernah di-refresh dari output SQL.

### 6g. ⚠️ KOREKSI: adjuster/consultant **hanya wajib pada adjustment PERTAMA**

**Ronde 1 §16 salah.** `[terverifikasi]` Gerbang di `AddAdjustment_Act` step 9 memang **aktif** dan
arahnya **normal** — sejauh itu benar. **Tetapi ia memuat DUA baris syarat yang di-AND-kan:**

```
baris 1: ConsultantID == "" || AppointedADJID == ""
baris 2: @LengthOfPageList(pyWorkPage.ClaimData.AdjustmentList) = 1
```

Karena step 4 sudah meng-`<APPEND>` satu baris, panjang = 1 **hanya pada adjustment pertama**.
**Adjustment ke-2 dan seterusnya dapat ditambahkan tanpa adjuster/consultant terisi.**

### 6h. ⚠️⚠️ Bug baru: baris adjustment yatim tertinggal saat validasi gagal

`[terverifikasi]` `AddAdjustment_Act` step 11 dan 12 keduanya **TERBALIK** (`Local.err==""`, 3/2) —
jalan bila error **ADA**. Step 11 memasang pesan, step 12 (+12.1) melakukan **rollback**
(`AktifButton=0` dan `Property-Remove AdjustmentList(<LAST>)`).

**Tetapi step 11 punya transisi `pyStepsTransParamsWhen = true` dengan `WhenTrue = 6` = keluar
activity.** Jadi begitu step 11 jalan, **activity langsung keluar dan step 12/12.1 TIDAK PERNAH
dieksekusi.**

**Akibatnya:** baris adjustment kosong yang sudah di-append di step 4 **tetap tertinggal** di
`AdjustmentList` meski pesan error muncul, dan `AktifButton` tetap `"1"`. Rollback yang ditulis
dengan sengaja **tidak pernah jalan**.

### 6i. `ProteksiInitialandDate_Act` — gerbang aktif, tetapi pesan tertempel di kolom yang salah

`[terverifikasi]` Step 6 di-remark (satu-satunya). Ketiga gerbang pesan **aktif dan normal**:
step 4.2 (`"Data Bank Account Can't NULL"` → field `.BranchOfBank`), step 5
(`"Please choose Payment Type first"` → field `.PayableTo`), step 7
(`"Please choose payable first"` → field `.Type`).

⚠️ **Pesan dan field tertukar:** step 5 dipicu `.Type==""` tetapi menempel di `.PayableTo`;
step 7 dipicu `Payable==""` tetapi menempel di `.Type`. Di UI pesan muncul di kolom yang salah.

⚠️ Step 5 dan 7 berada **di luar** loop step 4 → dievaluasi terhadap step page primer, bukan
per-baris `AdjustmentList`.

### 6j. `ProteksiData_act` — Cause of Loss **TETAP BERLAKU** sebagai gerbang wajib

`[terverifikasi]` Step 14 (`CauseOfLoss==""`) `PRE=true`, arah **normal**, tidak di-remark →
klaim Ronde 1 **benar**. Sembilan gerbang wajib lain juga aktif dan normal (Estimation list,
SpreadingList, `.Type` estimasi, share spreading, interest, `IDMaster`, `Location`, `PolicyNo`).

`[terverifikasi]` Step 18 dan 19 **di-remark** — duplikat pesan Cause of Loss.
Baris terbalik satu-satunya di step 9, dan efeknya **benign** (justru mencegah loop atas daftar
kosong).

---

## 8. RDBList — inventaris 54 rule (sebelumnya nol disentuh)

### 8a. ⚠️ KOREKSI §3b: **EMPAT schema Oracle, bukan tiga**

`[terverifikasi]` Schema keempat: **`gl`** — `RDBList/GetEmailCeding_SQL.xml`:
`select gl.f_get_email({TempCeding.CARI1}) AS HASIL1 from dual`.

⚠️⚠️ `[terverifikasi]` **Lebih penting: mayoritas tabel TIDAK diprefiks schema sama sekali** —
`bankaccount`, `agent`, `m_client`, `treatyinproduction`, `json_polis`, `os_akseptasi_klaim`,
`currency`, `treatyreinsurer`, `marketingofficer`, `treatybusiness`, `treatyyear`, `rw`, `city`,
`business`. Resolusinya bergantung pada **default schema koneksi**.

Kasus terburuk: `CekLunasPremi_Sql` menulis `from arasapas.invoice a, detail_invoice b` — tabel kedua
**tanpa prefix** berdampingan dengan yang berprefix. **Join itu bisa mengenai tabel yang salah bila
default schema berubah.** → `[terbuka]` untuk DBA.

### 8b. ⚠️ KOREKSI §7f Ronde 1: **16 penulis, bukan 15** — tetapi 13 COMMIT memang benar

`[terverifikasi]` Penulis tambahan yang terlewat: **`GetIDConsultanAdj_SQL`** — ia `SELECT`, tetapi
**memajukan sequence** (`ADJUSTERCONSULTANT_SEQ.nextval`). Efek samping di dalam rule bernama "Get".

⚠️ `[terverifikasi]` **Dua penulis DML asli TANPA `COMMIT`:**

| Rule | Operasi | Objek |
| --- | --- | --- |
| `DeleteDataTreatyGroup_SQL` | **DELETE** | `m_treatygroup` |
| `Update_T_Storage_SQL` | **UPDATE** | `T_STORAGE_IMAGE` |

Kontras langsung: `Insert_T_Storage_SQL` — **induk save-as** dari `Update_T_Storage_SQL` — **punya**
`commit;`. Jadi pola COMMIT **tidak konsisten** dan bergantung apakah penulis rule ingat menyalinnya.

**Catatan metode** `[terverifikasi]`: menghitung `commit` dengan grep polos ke XML memberi **54 hasil
palsu** karena tag metadata `<pxCommitDateTime>`. Hitungan yang benar hanya di dalam
`<pyBrowseSQL>` → ~~**tepat 13**~~.

> ### ⚠️ RALAT 2026-09-17 — **15 dari 58, bukan 13 dari 54**
>
> Angka **13** dihitung saat `RDBList/` masih **54 berkas**. Setelah empat rule ditambahkan (P21),
> sensus dijalankan ulang atas **58 dari 58 berkas** `Claim Prop/RDBList/`, mencari `COMMIT` **di
> dalam `<pyBrowseSQL>`**:
>
> | | Jumlah |
> | --- | ---: |
> | **Ber-`COMMIT` sendiri** | **15 dari 58** |
> | Tanpa `COMMIT` | **43 dari 58** |
>
> `[terverifikasi]` Kelima belas: `GenerateNoCLMTreatyIn` · **`GenerateNoDlaTreatyIn`** (BARU) ·
> **`GenerateNoPLATreatyIn`** (BARU) · `GenerateNoTRTInTemp` · `GetSequenceNumber_SQL` ·
> `GetTokenStorage_SQL` · `InsertClaimPNC` · `InsertLOGDirectKasir_SQL` · `InsertLogServiceClaim` ·
> `Insert_T_Storage_SQL` · `SaveDataToOsAkseptasiNP` · `SaveOSClaim_SQL` · `SaveOSKlaimTreaty_SQL` ·
> `UpdateDCauseOfLoss` · `UpdateMCauseOfLoss`.
>
> **Catatan metode tetap berlaku** — grep polos ke XML tetap memberi hasil palsu karena tag
> `<pxCommitDateTime>`; hitungan sah hanya di dalam `<pyBrowseSQL>`.
>
> ⚠️ **Rujukan angka 13 di luar berkas ini:** `grilling-ronde-1.md` **§7f** (judul *"Titik potong
> transaksi — 13 rule COMMIT sendiri"*) masih memuat angka lama. Berkas itu **di luar lingkup tugas
> ini dan tidak diubah** — dicatat di sini sebagai penunjuk. Di dalam berkas ini, satu-satunya
> rujukan lain (§27d) **sudah** memakai `13 dari 54 → 15 dari 58`.
>
> ### ⚠️⚠️ Yang lebih penting daripada selisih dua: **polanya**
>
> `[terverifikasi]` **Lima dari lima belas rule ber-`COMMIT` adalah GENERATOR NOMOR** —
> `GetSequenceNumber_SQL` · `GenerateNoCLMTreatyIn` · `GenerateNoPLATreatyIn` ·
> `GenerateNoDlaTreatyIn` · `GenerateNoTRTInTemp`. **Semuanya `COMMIT` sendiri, di dalam procedure,
> sebelum pemanggilnya selesai.**
>
> **Akibatnya nomor terbakar begitu di-generate:** bila langkah mana pun sesudahnya gagal, nomor itu
> **sudah ter-commit dan tidak dapat ditarik kembali**.
>
> ⚠️ **Ini bukan cacat satu rule — ini pola sistemik pada SELURUH jalur penomoran modul.**
>
> **Menguatkan §10e** (`Activity/TryMakePLA_Act.xml` **nol `Obj-Save`** setelah sequence diambil):
> di sana **gejalanya** terlihat — nomor diambil lalu hilang saat `Java` melempar; **di sini
> sebabnya** — procedure-nya memang commit sendiri. Dua temuan yang selama ini terpisah adalah satu
> masalah.
>
> **Bahan langsung untuk keputusan batas transaksi di spec:** di Go, **pengambilan nomor dan
> penyimpanan hasilnya harus berada dalam satu transaksi**, **atau** nomor harus **dapat
> dikembalikan**. Ini melengkapi §7 (butir PREFACTOR) dan menjawab sisi teknis Q12 yang sudah
> dijawab work owner di §21.

### 8c. Silsilah — ~~**21 dari 54 (39 %)**~~ **23 dari 58 (39 %) adalah hasil save-as**

> ⚠️ **RALAT 2026-09-17 — penyebut berubah, persentasenya bertahan.** Sensus ulang atas **58 dari 58
> berkas**: **23 save-as**, tetap **39 %**. Dua tambahan dari berkas baru —
> `RDBList/GenerateNoDlaTreatyIn.xml` (asal `GenerateNoDla`, induknya **tidak ada** di modul) dan
> `RDBList/GetTreatyInMaster_SQL.xml` (asal `GetLimitsTreatyIn_SQL`, induknya **ada**). Lihat §26a.
> ⚠️ Dihitung sesuai **Aturan silsilah** (§0): kemunculan nama di `pzOriginalInstanceKey` adalah
> **asal**, bukan rujukan.

`[terverifikasi]` Beberapa berbagi *timestamp* instance asal yang identik → lahir dari induk yang sama:

| Timestamp asal | Bersaudara |
| --- | --- |
| `#20180316T035719.205` | `GetAddressCeding`, `GetAddressTreatyIn`, `GetLeaderReport` |
| `#20250730T025520.617` | `GenerateImageID_SQL`, `Insert_T_Storage_SQL` |
| `#20250730T030543.459` | `GetTokenStorage_SQL`, `Update_T_Storage_SQL` |
| `#20180116T023908.482` | `SaveOSClaim_SQL`, `SaveOSKlaimTreaty_SQL` |
| `#20231208T040022.143` | `GETTanggalClosing_SQL`, `GetSequenceNumber_SQL` |
| `#20191114T022240.600` | `GetDataBankAccount_sql`, `GetDataBankAccount2_sql` |

⚠️ Satu rule membawa jejak **ruleset pra-GCNM**: `UpdateDCauseOfLoss` berasal dari class
`ASM-FW-**CNM**FW-…` dengan prefix `CNM!`, sedangkan sekarang `ASM-FW-**GCNM**FW` / `GCNM!`.

### 8d. ⚠️ **Nol rule yatim** — tetapi **7 `RequestType` menggantung**

`[terverifikasi]` Seluruh 54 rule punya minimal satu perujuk. **Yang berisiko justru kebalikannya:**
tujuh rule **dipanggil tetapi tidak ada berkasnya** di korpus.

| `RequestType` dipanggil | Status pemanggil |
| --- | --- |
| `GetBreakDownTreaty` | dipanggil dari step **DI-REMARK** (§1a) → tidak berdampak |
| `InsertCatastrope_SQL` | dipanggil dari step **DI-REMARK** (§1a) → tidak berdampak |
| `SaveMasterAdjusterConsultant_sql` | dipanggil dari step **DI-REMARK** (§1a) → tidak berdampak |
| **`GenerateNoDlaTreatyIn`** | ⚠️ perlu dicek status step pemanggilnya |
| **`GenerateNoPLATreatyIn`** | ⚠️ idem — dan ini **induk save-as** dari `GenerateNoCLMTreatyIn` |
| **`GetDataNopolisTreatyin`** | ⚠️ idem |
| **`GetTreatyInMaster_SQL`** | ⚠️ idem — **induk save-as** dari `GetTreatyInMasterProp_SQL` |

**Sintesis:** tiga yang pertama sudah terjelaskan oleh remark. **Empat sisanya adalah celah nyata** —
rule hidup yang tidak ikut diekspor. Tiga di antaranya adalah **induk save-as** yang masih dipanggil,
artinya induk-induk itu **masih hidup di produksi**, hanya tidak ikut paket ini. → `[terbuka]`,
pemilik export Pega.

### 8e. ⚠️ Lima rule **mati secara fungsional** — hanya terjangkau dari harness admin

`[terverifikasi]` Rantai rujukannya berhenti di harness mandiri, **bukan** di `Flow/Flow_TreatyIn.xml`:

| Rule | Ujung rantai |
| --- | --- |
| `UpdateMCauseOfLoss` | `Harness/TambahMasterCauseOfLoss` |
| `UpdateDCauseOfLoss` | `Harness/TambahCauseofLoss` |
| `GetLBUID_SQL` | klaster Cause-of-Loss yang sama |
| `DeleteDataTreatyGroup_SQL` | `Harness/MstAdjusterConsultant` |
| `GetIDConsultanAdj_SQL` | klaster master adjuster |

⚠️ **Kedua penulis paling berbahaya ada di klaster ini**: `DeleteDataTreatyGroup_SQL` (DELETE tanpa
COMMIT) dan `GetIDConsultanAdj_SQL` (memajukan sequence). Bila Claim Prop di-deploy tanpa portal
admin, kelimanya efektif mati. **Keputusan batas modul** → masuk frontier.

### 8f. Sensus nilai hardcode — **17 dari 54 rule**

`[terverifikasi]` Selain yang sudah diketahui (`name='<nama orang>'`, `treatydescid='10001'`,
`type='5' and STS_AKSEP='1'`, `CURRENT_SITE='1'`, `STORAGE='standard'`):

**Flag status:** `STS_CODE='1'` (3 rule bankaccount) · `MOSTATUS='1'` · `isactive='1'`

**Kode bisnis:** `TYPE='NONLIFE'` (`GetKodeProdNonLife_SQL`) ·
`PROPORTIONTYPE='Proportional'` (`GetTreatyInMasterProp_SQL`, dua cabang UNION) ·
`'ASMPP'` salt + algoritma `'MD5'` (`GenerateImageID_SQL`) ·
`'0' AS "Prodke"` (`SetPolicyTreatyProp` — **nilai keluaran yang dikarang di klausa SELECT**)

⚠️⚠️ `[terverifikasi]` **`Update_T_Storage_SQL` memakai DUA format tanggal berlawanan dalam satu
statement**: `EXPDATE` pakai `'DD/MM/YYYY HH24:MI:SS'`, `TANGGAL_UPLOAD` pakai
`'MM/DD/YYYY HH24:MI:SS'`. **Salah satu hampir pasti keliru** → `[terbuka]`.

⚠️ `[terverifikasi]` `GetIDConsultanAdj_SQL` mem-`lpad` sequence ke **4 digit** → **jebol setelah
9.999 konsultan**.

`[terverifikasi]` **Nomor polis yang ditanam: TIDAK ADA** di seluruh 54 rule. (Ini mengoreksi dugaan
di §1 tentang `GetDtlPaymentPremi_act` — polis hardcode-nya ada di Activity, bukan di SQL.)

### 8g. Nama yang berbohong — 15 tambahan

`[terverifikasi]` Yang paling material:

| Rule | Kebohongan |
| --- | --- |
| **`BrowseRW_SQL`** | class `Int-V_POLIS` tetapi membaca tabel **alamat** (`rw`, `city`), dan **seluruh aliasnya menyesatkan**: `DISTRICTID` disajikan sebagai `POLICY_NO`, `provincename` sebagai `CUSTOMER`, `note` kota sebagai `CURRENCY`; parameternya bernama `BUSINESS_CODE` padahal isinya **kode pos** |
| **`GenerateImageID_SQL`** | **bukan generator ID** — `STANDARD_HASH('ASMPP' \|\| TO_CHAR(SYSTIMESTAMP,…), 'MD5')`. Tidak menyentuh tabel maupun sequence; keunikan bergantung resolusi **milidetik**. **Dua unggahan dalam milidetik sama menghasilkan ImageID identik** |
| **`GetDataBankAccount2_SQL`** | angka "2" menyembunyikan perubahan semantik: induknya memfilter `clientid AND CURRENCYID`, versi ini **membuang filter `clientid`** → hasil mencakup **semua klien**. Dipakai berdampingan di `SetPayableTreaty_Act` |
| `GetTokenStorage_SQL` | rule bernama "Get" yang **men-COMMIT** |
| `GetLeaderReport` | bukan report — `select id from agent where clientname=…`; save-as dari rule alamat yang dikosongkan |
| `SetPolicyTreatyProp` | "Set" tetapi **tidak menulis apa pun** |
| `CekLunasPremi_Sql` | tidak mengembalikan status lunas — ia `sum(ivd_trans_sign * ivd_total)`; keputusan lunas diserahkan ke pemanggil |
| `CurrencyStandard`, `GetCurrency` | silsilahnya `UPDATEMASTER…` padahal keduanya read-only |
| `GetAppName_SQL`, `GETTanggalClosing_SQL` | **tanpa klausa WHERE** — menarik seluruh tabel lalu memakai baris pertama secara implisit |
| `GetStatusKasir_SQL` | mengambil kolom **`KET`** (keterangan bebas), bukan kode status terstruktur |

### 8h. `DATA_JSON` — jendela ke isi JSON lama, dan **dua konvensi berdampingan**

`[terverifikasi]` **Konvensi 1 — `DATA_JSON` pada tabel transaksi:**

| Tabel | Atribut yang terbaca lewat dot-notation |
| --- | --- |
| `json_klaim` | `DateOfLoss`, `ClaimNo`, `CauseOfLoss`, `ClaimEstimate` (4) |
| `os_akseptasi_klaim` | `Currency`, `CurrencyID`, `GrossValue`, `Value`, `AcceptedNo`, `Type`, `KursValue`, `TotalGross`, `PersenRNM` (9) |
| `json_polis` | `YearOfQuartal` (1) |

`[terverifikasi]` **Konvensi 2 — `JSONDATA` pada tabel master:** `m_client` (`AddressList.ASMAddress`,
`.RWName`, `.DistrictName`, `.CityName`, `.ASMZipCode`), `m_treatygroup` (`ID`),
`m_treaty_in`/`m_treaty_in_edm` (kolom utuh).

⚠️ `[terverifikasi]` `DeleteDataTreatyGroup_SQL` **menghapus dengan predikat atribut JSON**:
`DELETE FROM m_treatygroup a where a.JSONDATA.ID = {InputData.CARI1}`.

⚠️ `[terverifikasi]` **Bug alias:** `GetDataKlaimTreatyin` memakai alias **`CARI5` dua kali** —
`DATA_json.PersenRNM AS CARI5` dan `NOCLAIM AS CARI5`. Nilai kedua menimpa yang pertama →
**`PersenRNM` kemungkinan besar tidak pernah sampai ke Pega.**

⚠️ `[terverifikasi]` Pemetaan slot `CARIn` **tidak konsisten** antar dua pembaca tabel yang sama:
`GrossValue → CARI3` di `DataOutstandingTreatyin`, `GrossValue → CARI12` di `GetDataKlaimTreatyin`.

⚠️ `[terverifikasi]` `GetAddressTreatyIn` (save-as dari `GetAddressCeding`) **kehilangan atribut
`ASMZipCode`** — varian membaca 4 atribut, induknya 5.

**Nilai untuk P16:** daftar 19 nama atribut di atas adalah satu-satunya jendela korpus ke isi JSON
lama. Ia **tidak** menjawab apakah `px*`/`py*` ikut terserialisasi — itu tetap butuh satu baris dump
dari DBA.

---

## 10. Ralat §2, §18, §19 Ronde 1 — hasil audit gerbang

### 10a. ⚠️⚠️ KOREKSI besar §2d: `SendCloseClaimToKomite` **memang berhenti**

**Ronde 1 salah.** `[terverifikasi]` Benar bahwa step 4 (`Page-Set-Messages`) **tidak** menghentikan
activity — kode SALAH-nya `2` (lanjut), bukan `6`, dan tidak ada `Exit-Activity` di rule ini.

**Tetapi step 5 — blok "UNTUK CREATE KOMITE" — punya gerbang AKTIF** `Local.Error==""` (T=2 F=3,
normal). Jadi seluruh blok 5.1–5.9, **termasuk** `Call pxAddChildWork` (5.4), `Obj-Save` (5.8) dan
email (5.9), **dilewati bila ada error**.

**Perilaku sebenarnya:** activity berjalan sampai habis **tanpa efek apa pun**, pesan tersangkut di
`pyWorkPage`. Efek nettonya **menghentikan**.

⚠️ Nuansa yang tetap berlaku: gerbang step 5.4 sendiri (`ParamCari.CARI2==1`) memang **mati**, jadi
begitu blok 5 jalan, child case **selalu** dibuat. Gerbang yang menahan ada di **induknya**, bukan di
step itu. Dan gerbang induk hanya menguji `.AcceptanceStatus==""` — sehingga bila adjustment sudah
ber-status, **dua case komite tetap dapat lahir pada klaim yang sama**.

`[terverifikasi]` `SendCloseClaimToKomite` **nol step di-remark**. Roster hardcode satu orang
**aktif, tanpa gerbang**, di **dua** step (5.2 dan 5.3) — **TETAP BERLAKU**.

### 10b. `AddKomiteTreatyChild_ACT`

`[terverifikasi]` **Satu-satunya step di-remark: step 10** — blok iterasi `FacRetroList` beserta anak
10.1 (`.TFAllObj = Local.TFAllObj * .PctShareAllObj/100`). **Perhitungan fac-retro TFAllObj tidak
pernah jalan.**

`[terverifikasi]` **Lima gerbang mati:**

| Step | Syarat yang tidak berlaku | Akibat |
| --- | --- | --- |
| 7.2 | `ExGratia==1` | `Local.TotalEstimasireas` & `.TotalSpread` di-set untuk **semua** kasus |
| 9.3 | `ExGratia==1` | `Local.TFAllObj = .TSISpreaded` tanpa saringan |
| 22.2 | `Local.Subjectivity==true` (terbalik) | baris `ComiteeClaim(<APPEND>)` **selalu** ditambah; niat "kalau subjectivity jangan tambah komite" **mati** |
| 30 | `KomiteLoop==KomiteCount` | `pyWorkPage.AktifButton = 0` **selalu** |
| 32 | `.KomiteNo=="" && .IsKomite=1` | `.KomiteNo` **ditimpa** walau sudah terisi |

`[terverifikasi]` **TETAP BERLAKU**: step 13 (`TransferType=2`, `IndexAdjustment=.pxListSubscript`,
`CLMNO`), step 21 (roster lewat RD — **dinamis, bukan hardcode**), step 33 (`Obj-Save`), step 34
(`Call SendEmailKlaim`) — semuanya **tanpa gerbang**.

`[terverifikasi]` Step 25 (`Call pxAddChildWork`) bergerbang **aktif** `@hasMessages(pyWorkPage)`
dengan **T=6 (keluar activity)**; step 24 `Exit-Activity` sudah menjaga lebih dulu.

### 10c. Efek keluar — ralat

`[terverifikasi]` **KOREKSI:** `KonversiKlaim_Act` step 4, `InsertGoogleStorage_Act` step 13, dan
`GetUrlGoogleStorage_Act` step 8 ber-`pyStepsBlockName` **`END`/`EXIT`** — itu **label blok**, bukan
remark. Step-step itu **jalan normal**. (Lihat §5a.)

⚠️⚠️ `[terverifikasi]` **`HitServiceToKasir_Act` step 1 DI-REMARK** — blok proteksi
*"kalau salvage jangan lanjut"* (`(CLMP- && .Type=="3") || (CLM- && .PaymentType=="3") ||
(CLMNP- && .PaymentType=="3")`, T=6). **Proteksi salvage MATI TOTAL** meski gerbangnya sendiri aktif.

`[terverifikasi]` **TETAP BERLAKU:** kedua `Connect-REST` (step 9.7 dan 13.4) **memang** digerbangi
`IsPEGAPROD` yang **aktif dan normal**. Begitu pula `KonversiKlaim_Act` step 3.

⚠️ `[terverifikasi]` Tetapi **step log 9.9 dan 13.6 tanpa gerbang sama sekali**, dan step 9.8/13.5
(penulis `.StatusKasir`) **gerbangnya mati** → **di non-PROD, log tetap ditulis dan `.StatusKasir`
tetap ditimpa** dari respons yang kosong/basi. `Obj-Save` step 14 juga tanpa gerbang.

⚠️⚠️ `[terverifikasi]` **Google Storage menembak di SEMUA environment.** `Connect-REST` di
`InsertGoogleStorage_Act` step 11 dan `GetUrlGoogleStorage_Act` step 6.5 **tidak punya gerbang
`IsPEGAPROD` sama sekali** (`pyStepsPreCondition` kosong) — berbeda dari Kasir dan Konversi. Ditambah
step 6 `InsertGoogleStorage_Act` yang gerbang `IsPEGAPROD`-nya **mati**. **Ini memperkuat Q4 (§21).**

`[terverifikasi]` **TETAP BERLAKU** — `InsertDocument_Act` step 5 (`Obj-Save`): gerbang **aktif**,
arah **terbalik**, dan terbaliknya **benar secara desain** — baris dokumen **tidak** tersimpan saat
upload gagal. Rantainya utuh karena `InsertGoogleStorage_Act` step 12 juga aktif-terbalik.

⚠️ `[terverifikasi]` `GetUrlGoogleStorage_Act` step 6 **terbalik**: blok refresh signed-URL jalan
ketika `@CompareDates(UploadDoc.exp, now)` bernilai **FALSE**. Bila `exp` kosong atau tak terbaca,
compare menghasilkan false → blok **tetap** jalan dan menembak Google. **Rapuh.**

### 10d. Email — **KOREKSI**: yang hardcode tanpa gerbang adalah **BCC**, bukan CC

`[terverifikasi]` `SendEmailKlaim`:

| Step | Isi | Gerbang |
| --- | --- | --- |
| **2** | `Local.BCC` = **dua alamat tertanam** | **TANPA GERBANG → aktif di SEMUA environment** |
| 3 | `Local.EmailCC` = satu alamat tertanam | `IsPEGAPROD` **aktif** → PROD saja |
| 8.7 | `Call SendEmailWithAttachments` | `IsPEGAPROD` **aktif** → PROD saja |

⚠️⚠️ `[terverifikasi]` **Step 8.3 gerbangnya MATI**: syarat tertulis
`pxRequestor.pxReqContextURI == "<URL host dev internal>"` — pembatasan ke satu host dev yang **tidak
berlaku**. Akibatnya penerima **selalu** diambil dari `.ComiteeClaim(1).KomiteEmail` — **anggota
pertama saja**. Ini temuan lingkungan ketiga setelah `jboss1074` dan `pzProductionLevel`.

`[terverifikasi]` Empat gerbang mati lain di rule ini (8.1, 8.2, 8.2.1, 8.4) membuat blok badan email
**jalan tiap iterasi** dan nilainya **tertimpa baris terakhir**.

### 10e. Dokumen cetak — ralat besar

⚠️⚠️ `[terverifikasi]` **`TryMakePLA_Act` — jalur penomoran LAMA di-remark.** Step 5, 6, 7, **8**
(`RDB-List GenerateNoPLATreatyIn`) semuanya `//`. **KOREKSI Ronde 1** yang menyebut
`GenerateNoPLATreatyIn` masih aktif.

Yang **aktif** adalah `GetSequenceNumber_SQL` **step 11**, bergerbang **terbalik** pada `NoPla!=""` →
sequence ditarik **hanya bila PLA belum bernomor**.

⚠️⚠️ **Klaim "sequence terbakar" TETAP BERLAKU — dan lebih kuat.** `[terverifikasi]`
**`TryMakePLA_Act` tidak punya step `Obj-Save` sama sekali.** Urutannya: step 11 menembak
`PROC_GENERATE_SEQUENCE_NUMBER` → step 14 menaruh hasilnya **hanya di clipboard** → step 26 (`Java`,
**aktif**) dapat melempar `PRRuntimeException`. **Nomor sudah dikonsumsi di DB tetapi tidak pernah
dipersist → lubang sequence.** `PrintDLATreatyIn` punya `Obj-Save` di step 17, tetapi `Java` step
15.18 juga bisa melempar sebelum mencapainya → lubang yang sama.

⚠️⚠️ `[terverifikasi]` **`PrintDLATreatyIn` — enam step di-remark**, termasuk **15.11.2** (fac-out
payment type 3 / salvage) dan **15.11.3** (type 4 / adjuster fee). Hanya **15.11.1**
(`.TreatyType=="10004"`) yang hidup → `Local.Estimasi` **selalu** memakai `AdjustmentValue`.
**DLA untuk salvage dan adjuster fee mencetak angka dari cabang yang salah.**

> ### ⚠️ RALAT 2026-09-17 — label "fac-out" di kalimat di atas berasal dari keterangan step
>
> **Yang TETAP BERLAKU** (dari `pyStepsBlockName`, bukan dari keterangan): enam step di
> `Activity/PrintDLATreatyIn.xml` **di-remark**, termasuk **15.11.2** dan **15.11.3**; hanya
> **15.11.1** yang hidup.
>
> **Yang DIRALAT:** sebutan *"fac-out"* pada ketiga step itu saya ambil dari **keterangan step**
> (*"set fac out payment type …"*), dan `[data DBA]` **satu baris contoh** menunjukkan
> `TreatyType 10004` sebenarnya **`QS (R/I)`**, bukan Fac Out (§2). Keterangan step **tidak boleh**
> dipakai menetapkan arti kode — §0 Aturan arti kode.
>
> **Perumusan ulang yang aman:** step **15.11.1** bergerbang `.TreatyType == "10004"` — yaitu
> **satu jenis treaty tertentu**, yang pada satu baris contoh bernama `QS (R/I)`. Karena dua cabang
> tetangganya di-remark, `Local.Estimasi` **hanya terisi pada baris spreading ber-`TreatyType`
> `10004`**, dan baris ber-`TreatyType` lain (`10028`, `10035`, `10236` pada contoh yang sama)
> **tidak pernah mendapat nilai dari cabang mana pun**.
>
> ⚠️ Apakah akibatnya "mencetak angka dari cabang yang salah" atau "tidak mencetak angka sama
> sekali" **tidak dapat dipastikan dari korpus** — bergantung nilai awal `Local.Estimasi`, yang tidak
> terbaca → `[terbuka]`. Kalimat lama yang menyatakannya sebagai kepastian **dilunakkan**.

`[terverifikasi]` **KOREKSI:** `PrintFileAcceptance` **nol step di-remark**.
⚠️ Tetapi **step 11 gerbangnya MATI** (`IsFire`) → template Acceptance Note dipakai untuk **semua
lini bisnis**, bukan hanya Fire. Step 10 juga mati (syaratnya sendiri **tautologi**).

`[terverifikasi]` `Call InsertDocument_Act` **aktif tanpa gerbang** di ketiganya (PLA step 27,
DLA step 15.19, Acceptance step 15). Yang menahan di hulu adalah step `Java` ber-`param.InsHandle==""`
dengan **T=6**.

---

## 11. DataTransform (11) — jejak audit dan cacatnya

### 11a. `InsertChronology_DT` — satu-satunya mesin jejak audit

`[terverifikasi]` Dipanggil dari **27 Activity**; **nol** di antaranya di-remark. Ia meng-`APPEND`
satu halaman ke **`pyWorkPage.ClaimData.SuggestList`** — ⚠️ class barisnya
**`ASM-FW-GISFW-Data-OfferFacIn-SuggestList`**, yaitu **struktur milik modul OfferFacIn yang
ditumpangi**.

| Properti | Nilai |
| --- | --- |
| `.CommentSuggest` | `DataChronology.CARI1` — teks aksi |
| `.PICSuggest` | `OperatorID.pyUserName` |
| `.DateSuggest` | `@DateTime.CurrentDateTime()` |
| `.IsCedingConfirm` | `"Claim Admin"` default; tiga **nama orang hardcode** dipetakan ke jabatan |

⚠️⚠️ `[terverifikasi]` **Gerbang terluar: `OperatorID.pyPosition != "IT Developer"`.**
**Seluruh aksi operator berposisi "IT Developer" tidak meninggalkan jejak apa pun.** Ini lubang audit
struktural, bukan kelalaian per-kasus.

### 11b. ⚠️ Bug: riwayat penutupan klaim tercatat kosong atau basi

`[terverifikasi]` `CloseClaimProp` step 5 mengisi **`DataChronology.CARI12`**, sedangkan
`InsertChronology_DT` membaca **`DataChronology.CARI1`**. Berkas itu **tidak pernah menulis `CARI1`**.

Karena `DataChronology` adalah *named page* yang bertahan di clipboard, entri riwayat "Finish
Adjustment (Close Claim)" tercatat dengan komentar **kosong atau sisa nilai dari aksi sebelumnya**.
Residu ekspresi `@if(.pyNote == "","OK",Data.CARI12)` di step 1.1.1 memperkuat dugaan ini pernah
hendak diperbaiki lalu ditinggalkan.

### 11c. ⚠️ Lima gerbang mati membuat jejak audit salah label

`[terverifikasi]` Pada step pengisi teks aksi: `CountEstimation_Act`(23), `CountListClaimAmountIDR`(4),
`CountPersen_act`(8), `CountSpreading_Act`(12), `SaveOutstanding_Act`(23.1.2).

**Kasus paling merusak** — `CountListClaimAmountIDR`: step 3 (`"Edit Value Claim Amount"`, gerbang
**berlaku**) **selalu ditimpa** step 4 (`"Add Value Claim Amount"`, gerbang **mati**) sebelum Apply di
step 5. **Riwayat selalu tercatat "Add", tidak pernah "Edit".** Pola sama di `CountSpreading_Act`
(§6a) dan `CountPersen_act`.

⚠️ `[terverifikasi]` **Typo ikut tersimpan ke basis data riwayat**: `"Add Insured Insterest"`,
`"Send to Commite"`, `"Send Adjustment to Committe (Acceptation)"`, `"Edit  Value Estimation"` dan
`"Delete  Insured Interest"` (spasi ganda).

⚠️ `[terverifikasi]` **Jejak audit hanya mencatat JENIS aksi, bukan nilai sebelum/sesudah.** Nilai
uang tidak pernah ikut tersimpan di entri riwayat — batas nyata untuk rekonstruksi audit sistem lama.

### 11d. Temuan lain

`[terverifikasi]` `CekInterestListDtl_DT` bersifat **seed-once** (`@SizeOfPropertyList(InterestListDtl)
< 1`) → bila `InterestList` berubah setelah detail pernah terisi, **angka uang dan mata uang di
daftar detail tertinggal usang**. Dipanggil sebagai `pyPreProcessingTransformRule` FlowAction
`ViewDetailInterest` — **tanpa gerbang**, jalan tiap kali dibuka.

⚠️ `[terverifikasi]` `DisableEditRNMShare` menulis `.TreatyInMaster.RNMShareP` =
`@toDecimal(TreatyShareTemp.CARI1)` dari **field scratch bertipe teks**, **tanpa `When` maupun
validasi**, lalu **mengunci** `.IsEditRNMShare="false"`. Input non-numerik mendarat langsung di
properti share, dan **nilai salah tidak bisa diperbaiki lewat field yang sama**.

`[terverifikasi]` `SetDateOutstanding` dan `SetDateAcceptation` adalah stempel waktu tahapan, dan
karena **tanpa `When`**, membuka flow action yang sama dua kali akan **menimpa** stempel sebelumnya.

`[terverifikasi]` **Tidak ada DataTransform penetap nilai awal klaim baru** di korpus ini — bentuk
data awal harus dicari di Flow/Activity pembuat case.

`[terverifikasi]` **Nol DataTransform yatim.** Penulisan mati: `CloseStsSaveTreatyGroup` step 2
(`OutputData.HASILD7` tidak dikonsumsi siapa pun).

⚠️ `[terverifikasi]` Empat DataTransform Cause-of-Loss adalah **save-as dari transform Surveyors**
(nama asal `CNMREFRESHSURVEYORS_DT`, `CNMSHOWINSERTSUVERYORS_DT` — typo ikut terbawa), dan dua di
antaranya menyetel label `"Update"` pada aksi **"Add New"** — anomali logika.

⚠️ `[terverifikasi]` `CNMRefreshListDetailCauseOfLoss_dt` menunjuk class
**`ASM-FW-CNMFW-Int-V_D_CAUSE_OF_LOSS`** (tanpa `G`), sedangkan saudaranya memakai
`ASM-FW-GCNMFW-…`. Salah satu menunjuk class yang tidak konsisten.

---

### 11e. ⚠️⚠️ `IsCedingConfirm` adalah **TINGKAT WEWENANG** — dinaikkan dari §23c

> **Ditulis 2026-09-17.** Temuan ini semula dicatat di §23c sekadar sebagai "nama kolom yang
> berbohong". Ia lebih dari itu: **isinya tangga persetujuan**, jadi tempatnya di sini — di dalam bab
> mesin jejak audit — bukan di daftar nama yang menyesatkan.

`[terverifikasi]` **Dua rule menulis kolom ini**, dan tiga dari empat nilainya **identik persis**:

| Tingkat | `DataTransform/InsertChronology_DT.xml` | `Activity/SethistoryKlaimTreaty.xml` |
| --- | --- | --- |
| 1 (terbawah) | langkah **1.1.4** → `"Claim Admin"` | ⚠️ `"Admin Claim " + .PICSuggest` |
| 2 | langkah **1.1.5.1** → `"Claim Dept. Head"` | `"Claim Dept. Head"` |
| 3 | langkah **1.1.6.1** → `"Operational Director"` | `"Operational Director"` |
| 4 (teratas) | langkah **1.1.7.1** → `"Technical Director"` | `"Technical Director"` |

**Tangganya:** `Claim Admin` → `Claim Dept. Head` → `Operational Director` → `Technical Director`.

`[data DBA]` **Satu baris contoh** (klaim `RNM-K22.09.2026.T00718` / `CLMP-7134`) berisi
`"Admin Claim " + nama orang` — **persis pola penulis kedua**. Jadi entri audit pada klaim itu
ditulis oleh **`Activity/SethistoryKlaimTreaty.xml`**, bukan oleh `InsertChronology_DT`.
Nilai tidak disalin ke sini — PII.

#### 11e-1. ⚠️ Jejak auditnya **TIDAK KONSISTEN**

`[terverifikasi]` Tiga tingkat teratas ditulis **sama persis** oleh kedua penulis. **Tingkat terbawah
ditulis dua bentuk berbeda**, dan salah satunya **menempelkan nama orang ke dalam nilai yang sama**:

```
InsertChronology_DT      "Claim Admin"
SethistoryKlaimTreaty    "Admin Claim " + nama orang
```

**Tiga konsekuensi yang mengikat:**

1. ⚠️ **Menyaring atau mengelompokkan jejak audit berdasarkan kolom ini akan meleset** — ada **dua
   ejaan untuk tingkat yang sama**, dan salah satunya tidak pernah cocok dengan pencocokan persis.
   `[terverifikasi]` Korpus sendiri sudah mengakalinya: pembacaan
   `@contains(.IsCedingConfirm,"Admin") || .IsCedingConfirm==""` memakai **pencocokan sebagian**,
   bukan kesetaraan — pola yang hanya masuk akal bila nilainya teks bebas berimbuhan.
2. ⚠️ **Penyimpangan sadar:** di Go keempat tingkat menjadi **satu enum**, dan **nama pelaku dipisah
   ke kolomnya sendiri**. Menyatukan **peran** dan **identitas** dalam satu kolom adalah cacat yang
   **tidak ditiru**.
3. **Migrasi data lama wajib memetakan** `"Admin Claim <nama>"` → tingkat **`Claim Admin`** + nama ke
   **kolom pelaku**. Data lama memuat **kedua bentuk**, jadi pemetaan ini tidak opsional.

#### 11e-2. Sensus — **4 dari 270 berkas**, dan ia **tampil ke pengguna**

⚠️ **RALAT terhadap sensus saya sendiri di §23d**, yang hanya menyebut dua penulis.

`[terverifikasi]` `.IsCedingConfirm` disentuh **4 dari 270 berkas** `Claim Prop`:

| Berkas | Peran | Kemunculan |
| --- | --- | ---: |
| `Activity/SethistoryKlaimTreaty.xml` | **penulis** | 7 |
| `DataTransform/InsertChronology_DT.xml` | **penulis** | 5 |
| `Section/OutstandingClaim.xml` | **tampilan** | 2 |
| `Section/InputAcceptation.xml` | **tampilan** | 2 |

**Hasil pemeriksaan kedua Section** (diminta: apakah elemennya termasuk gerbang mati atau tersembunyi
permanen — daftar itu ada di **§14a**, bukan di bab 22; berkas ini tidak memiliki §22):

`[terverifikasi]` Pada **kedua** Section, kemunculannya berpasangan dan berperan berbeda:

- **Kemunculan pertama = sel grid yang benar-benar tampil.** `pxObjClass = Embed-Display-Table-Cell`
  (`pyCellId` 188 di `Section/OutstandingClaim.xml`, 167 di `Section/InputAcceptation.xml`),
  `pyValue = .IsCedingConfirm`, di dalam grid berclass
  `ASM-FW-GISFW-Data-OfferFacIn-SuggestList` — yaitu **grid jejak audit** yang sudah tercatat di §14b.
- **Kemunculan kedua = metadata indeks**, bukan kode: `Embed-Reference-Rule` dengan
  `pxRuleObjClass = Rule-Obj-Property`, `pxRuleFamilyName = ISCEDINGCONFIRM`.

`[terverifikasi]` **Status gerbangnya: `pyVisible = ALWAYS` dengan `pyCondition` KOSONG** di kedua
Section. Menurut **R1** (§13) itu berarti **tampil biasa** — **bukan** termasuk **44 dari 180** sel
bergerbang mati, dan **bukan** termasuk **12** sel tersembunyi permanen (§14a).

**Artinya tingkat wewenang ini memang sengaja ditampilkan ke pengguna di dua layar** —
Outstanding Claim dan Input Acceptation. ⚠️ Nama kolom yang tampil (caption) **tidak ditemukan pada
jendela yang saya periksa** → `[terbuka]` ringan; yang pasti adalah **isinya** yang tampil.

#### 11e-3. ⚠️ `[terbuka]` — **tiga daftar wewenang yang tidak saling merujuk** → **P24**

`[terverifikasi]` Modul ini memuat **tiga mekanisme wewenang yang hidup berdampingan tanpa satu pun
merujuk yang lain**:

| Mekanisme | Isi | Bukti |
| --- | --- | --- |
| **Tingkat di jejak audit** | 4 tingkat, **berhenti di `Technical Director`** | §11e |
| **Pencarian batas nilai** | menyebut **Direktur Utama** — yang **tidak muncul** di tangga jejak audit | `RDBList/GetLimitDirekturUtama_SQL.xml`; dipanggil `Activity/AttachmentProtect_ACT.xml` step **5** |
| **Roster komite** | `LIMIT_BOTTOM` · `LIMIT_TOP` · `TYPE_KOMITE` · `DEGREE` | `POOLDATA.EMAILKOMITE` lewat `ReportDefinition/FilterEmailKomiteWithLimit.xml` |

⚠️ **Tidak disimpulkan apakah ketiganya seharusnya satu tangga** — korpus tidak menyatakannya, dan
menebaknya akan menetapkan struktur wewenang dari sisa kode. → **P24** (§19).

⚠️ Ini **bukan bagian dari Q5** (§21). Q5 hanya soal **kunci pencarian** batas nilai (nama orang
versus jabatan); P24 soal **bentuk tangganya**.

> ### ✅ JAWABAN P24 — `[keputusan work owner]` 2026-09-17. **P24 TERTUTUP.**
>
> Ketiga daftar **bukan** tiga hal setara, dan **bukan** satu tangga tunggal. Pemetaannya:
>
> **Butir a — Direktur Utama berada DI LUAR tangga persetujuan klaim.**
> `RDBList/GetLimitDirekturUtama_SQL.xml` adalah **wewenang terpisah dengan urusannya sendiri**.
> Ia **bukan tingkat kelima**, dan **tidak berpadanan** dengan tingkat mana pun di jejak audit.
>
> **Butir b — roster komite dan jejak audit adalah TANGGA YANG SAMA.**
> `POOLDATA.EMAILKOMITE.DEGREE` **berpadanan** dengan tingkat di `.IsCedingConfirm`.
>
> **Butir c — `Claim Admin` adalah nilai DASAR, bukan anggota komite.**
> Ia **tidak ada** di `EMAILKOMITE.DEGREE`, dan **memang tidak seharusnya ada**.
>
> `[terverifikasi]` **Struktur rule membuktikannya** — `DataTransform/InsertChronology_DT.xml`:
>
> ```
> langkah 1.1.4                          .IsCedingConfirm = "Claim Admin"      ← TANPA SYARAT
> langkah 1.1.5  WHEN <nama orang 1>  →  "Claim Dept. Head"
> langkah 1.1.6  WHEN <nama orang 2>  →  "Operational Director"
> langkah 1.1.7  WHEN <nama orang 3>  →  "Technical Director"
> ```
>
> `"Claim Admin"` **di-set lebih dulu tanpa syarat, lalu ditimpa** oleh tiga cabang bernama. Jadi
> bentuknya **bukan empat tingkat sejajar**, melainkan:
>
> | | Bentuk sebenarnya |
> | --- | --- |
> | **`Claim Admin`** | **nilai dasar** — siapa pun yang **BUKAN** anggota komite |
> | Tiga tingkat di atasnya | **anggota komite**, bersumber dari `EMAILKOMITE.DEGREE` |
>
> ⚠️ Ini menjelaskan mengapa §11e-1 menemukan **tingkat terbawah ditulis dua bentuk berbeda** oleh
> dua penulis sementara tiga tingkat di atasnya identik: tiga yang atas berasal dari **satu sumber
> data**, sedangkan yang terbawah adalah **default yang diketik terpisah di tiap rule**.
>
> Aturan lengkapnya ditulis di **§28**.

---

### 11e-4. Empat `[terbuka]` yang ditutup — T9

> **Ditulis 2026-09-17.** Ketiganya sudah dapat diselesaikan dari korpus; satu lagi ternyata temuan
> baru. Ditulis sebagai ralat terhadap kesimpulan saya sendiri.

#### 9a. ✅ Caption kolom `.IsCedingConfirm` — **TIDAK ADA**

`[terverifikasi]` `Section/OutstandingClaim.xml` `pyCellId` **188** dan
`Section/InputAcceptation.xml` `pyCellId` **167**, keduanya ber-`pyCaption` = `"Label"` —
**placeholder bawaan Pega, bukan judul**.

**Tingkat wewenang tampil di grid tanpa nama kolom.** Tanda `[terbuka]` di §11e-2 **dicabut**,
diganti `[terverifikasi]`.

#### 9b. ⚠️ RALAT — lubang audit "IT Developer" **LEBIH LEBAR**, bukan lebih sempit

**Kesimpulan saya di §23d salah arah.** `[terverifikasi]` `DataTransform/InsertChronology_DT.xml`:

```
langkah 1      WHEN  OperatorID.pyPosition != "IT Developer"
  langkah 1.1    APPEND_AND_MAP_TO .ClaimData.SuggestList
```

**Seluruh penulisan jejak audit bersarang di bawah WHEN itu.** Bila posisi operator
`"IT Developer"`, **NOL entri ditulis** oleh rule ini — bukan sebagian bidang yang hilang.

**Yang benar dari koreksi saya:** `[terverifikasi]` `Activity/SethistoryKlaimTreaty.xml` **tidak**
punya pengecualian itu. Gerbangnya di step **1.1** berbunyi
`@contains(.IsCedingConfirm,"Admin") || .IsCedingConfirm==""` — soal **tingkat entri**, bukan
**posisi operator**. **Aksi yang lewat penulis kedua tetap terekam.**

**Jadi kedua sisi berlaku:** lubangnya **total** pada penulis pertama, dan **tidak ada** pada penulis
kedua.

#### 9c. ⚠️ TEMUAN BARU — tingkat wewenang ditentukan **nama orang yang di-hardcode**

`[terverifikasi]` `DataTransform/InsertChronology_DT.xml`:

```
langkah 1.1.4                                  .IsCedingConfirm = "Claim Admin"
langkah 1.1.5  WHEN .PICSuggest == <nama 1>  →  "Claim Dept. Head"
langkah 1.1.6  WHEN .PICSuggest == <nama 2>  →  "Operational Director"
langkah 1.1.7  WHEN .PICSuggest == <nama 3>  →  "Technical Director"
```

**Tiga nama orang**, ditambah **satu nama** di `RDBList/GetLimitDirekturUtama_SQL.xml` →
**empat orang di-hardcode** di modul ini. (Nilai tidak disalin — PII.)

⚠️ **Akibat:** orangnya pindah jabatan → **jejak audit salah label tanpa peringatan apa pun**.
Ditambahkan ke **P24** sebagai bukti penguat.

`[keputusan work owner 2026-09-17]` Butir ini **diikuti apa adanya** — nama tetap di kode (lihat P24
di §19, status 🟡 sebagian).

#### 9d. ✅ Bug `CARI12` — **PASTI**, bukan `[terbuka]`

`[terverifikasi]`

```
InsertChronology_DT  langkah 1.1.1   .CommentSuggest = DataChronology.CARI1
CloseClaimProp       step 5          DataChronology.CARI12 = "Finish Adjustment (Close Claim)"

sensus 270 dari 270 berkas:
  DataChronology.CARI1    34 kemunculan
  DataChronology.CARI12    1 kemunculan   ← hanya CloseClaimProp step 5
```

**Satu-satunya pemakai `CARI12` di seluruh modul adalah langkah yang salah itu.** Teks
`"Finish Adjustment (Close Claim)"` **tidak pernah sampai** ke jejak audit; entri close merekam
**apa pun yang tersisa di `CARI1`**.

`[data DBA]` Dikuatkan **satu baris contoh**: entri terakhir berbunyi `"Save Outstanding"`, **bukan**
teks close.

**Status dinaikkan dari `[terbuka]` menjadi `[terverifikasi]`** — meralat §23d dan §11b.

---

## 13. ⚠️ ATURAN BACA GERBANG DI SECTION — ditetapkan sebelum kesimpulan apa pun

Ini memenuhi syarat §F: mekanisme *visible-when* **berbeda** dari Activity, jadi aturannya ditetapkan
lebih dulu.

`[terverifikasi]` Blok gerbang di Section berbentuk `pyUserData` / `pyDefaultUserData` berkelas
`Embed-Harness-HeaderElements` atau `Embed-Harness-UserData`. Ada **tiga keluarga gerbang, masing-
masing dengan selektor modenya sendiri** — dan **satu keluarga tanpa selektor**.

| # | Aturan |
| --- | --- |
| **R1** | Gerbang tampil sel/kontrol: baca `pyCondition` **hanya bila `pyVisible = OTHER`**. Bila `ALWAYS` → selalu tampil, `pyCondition` **residu**. `NOTBLANK` → tampil bila properti terikat tidak kosong. `<pyVisible/>` kosong → tampil. |
| **R2** | Gerbang tampil container: baca `pyContainerVisibleWhen` **hanya bila `pyIsVisibilityOption = CONDITION`** (atau `ExpressionCondition`). Bila `ALWAYS` → selalu tampil. |
| **R3** | Gerbang disable: baca `pyDisabledWhen` **hanya bila `pyDisabledNew = true`**. `always` = selalu disable; `false` = tidak pernah. |
| **R4** | Gerbang wajib: baca `pyRequiredWhen` **hanya bila `pyRequiredNew = true`/`truewhn`**. `always` = wajib tanpa syarat. |
| **R5** | Gerbang read-only: `pyReadOnlyCondition` **dibaca langsung, tanpa flag pengaktif**. |
| **R6** | Semua ekspresi **dibaca lurus** — tidak ada kode arah/negasi seperti 2/3 di Activity. Negasi hanya lewat operator di dalam teks. |
| ~~**R7**~~ | ~~**Tidak ada gerbang pada aksi tombol** — `pyActionConditions` kosong di **seluruh 36 berkas**. Tombol yang tampil = tombol yang jalan.~~ **⚠️ DIBATALKAN — lihat R7-baru di bawah.** |

> ### ⚠️ RALAT 2026-09-17 — **R7 DIBATALKAN dan disusun ulang**
>
> R7 lama disimpulkan dari premis **"`pyActionConditions` kosong di mana-mana"**, dan premis itu
> **palsu**. Sensus penghitungan **anak dari setiap elemen `pyActionConditions` di 270 dari 270
> berkas** `Claim Prop`:
>
> | Berkas | Baris |
> | --- | ---: |
> | `Section/AdjustmentDetail.xml` | **10** |
> | `Section/OutstandingClaim.xml` | **6** |
> | `Section/InputAcceptation_Est.xml` | **1** |
> | `Section/AdjustmentDetail_Section.xml` | 0 |
> | `FlowAction/AdjustmentDetail.xml` | 0 |
> | (tidak ada berkas lain yang berisi) | — |
>
> **Jadi 3 dari 36 Section berisi, total 17 baris.**
>
> **R7-baru.** `pyActionConditions` **adalah** gerbang pada aksi tombol, dan ia **jarang tetapi
> nyata**. Tombol yang tampil **TIDAK** otomatis berarti aksinya jalan. Tiap baris adalah satu
> perbandingan (`pyOtherOperandLeft` + operator + nilai).
> ⚠️ `[terbuka]` **Cara menggabungkan beberapa baris (AND atau OR) tidak terbaca dari korpus** —
> jangan disimpulkan. Klaim sebelumnya bahwa dua baris di `Section/AdjustmentDetail.xml` digabung
> `AND` berasal dari penyisiran yang jumlah barisnya sendiri salah, jadi **tidak dapat dipercaya**.
>
> **Penyebab kesalahan: tabrakan nama.** Tiga berkas berbeda bernama `AdjustmentDetail` tersebar di
> `Section/`, `FlowAction/`, dan `Section/AdjustmentDetail_Section.xml`. Satu penyisiran menghitung
> nol karena menyisir berkas yang salah; satu lagi menyebutnya "Harness" padahal **tidak ada Harness
> bernama `AdjustmentDetail`** — berkas itu Section.

### R8 — `pyCondition` yang berisi label pemilih **bukan ekspresi**

**R8.** `pyCondition` yang berisi `"Other Property"`, `"When Rule"`, atau label pemilih lain
**bukan ekspresi** — itu status dropdown Designer. Sebelum menghitung sebuah `pyCondition` sebagai
gerbang, pastikan isinya ekspresi. Aturan sama berlaku untuk `pyStepsPreCondParamsWhen`.

### R9 — di rule `When`, `<pyConditionViewer>` adalah **sisa basi**

**R9.** Di berkas `When`, kondisi yang **benar-benar dieksekusi** ada di PageList `<pyCondition>`
(tiap `rowdata` ber-`<pyConditionLabel>` + `<pyConditionValue1>`), dan penggabungnya adalah
`<pyLogic>` **anak langsung `pagedata`**. Isi `<pyConditionViewer>` — termasuk `<pyConditionString>`
dan `<pyLogic>` **di dalamnya** — adalah **teks tampilan editor yang sering tidak sinkron** dan
**tidak boleh dipakai**.

⚠️ Bukti di §25a: `When/IsMarineCargo.xml` menyimpan enam baris kode bisnis di viewer yang **tidak
pernah dieksekusi**; `When/IsFire.xml` menyimpan `pyLogic` viewer dengan 21 label untuk satu baris
kosong. Tanpa R9, keduanya akan terbaca sebagai aturan aktif.
>
> ### ⚠️ ATURAN KERJA BARU — berlaku seterusnya
>
> 1. Rujuk berkas dengan **path lengkap dari akar modul** (`Section/AdjustmentDetail.xml`), **tidak
>    pernah dengan nama saja**.
> 2. Sebelum menyimpulkan "rule X begini", **pastikan dulu tidak ada berkas bernama sama di folder
>    lain**.
> 3. Setiap angka sensus **wajib menyebut penyebutnya** — "44 dari 180", bukan "44".
> 4. Sebelum menulis "tidak ada di korpus", **sebutkan berapa berkas yang dicari dan di folder mana**.

### 13a. Bukti yang mendasari

`[terverifikasi]` **Korelasi ko-kemunculan tanpa pengecualian di sisi flag-aktif:**

| Keluarga | Flag "pakai ekspresi" | Tanpa ekspresi | Dengan ekspresi |
| --- | --- | ---: | ---: |
| A (sel) | `OTHER` | **0** | 119 |
| A (sel) | `ALWAYS` | 1961 | **44 residu** |
| B (container) | `CONDITION` | **0** | 42 |
| B (container) | `ALWAYS` | 217 | **13 residu** |
| C (disable) | `pyDisabledNew=true` | 0 | **45/45** |
| C (wajib) | `pyRequiredNew=true/truewhn` | 0 | **9/9** |

⚠️ **Keluarga C justru yang membuktikan aturannya**: ketika flag dan ekspresi dikelola benar,
keduanya **selalu** berpasangan — nol yatim. Di keluarga A dan B ada yang tidak berpasangan; itulah
residunya.

`[terverifikasi]` **Bukti perilaku yang menentukan** — `Section/Catastrope_Sec.xml`
(`ASM-FW-GCNMFW-WORK!CATASTROPE_SEC`): sel `.ClaimData.StsKatastrofe` (`pxRadioButtons`) ber-
`pyVisible=ALWAYS` + `pyCondition=1=2`. Sel tetangganya `.ClaimData.NonKatastrofeType` ber-
`pyVisible=OTHER` + kondisi `StsKatastrofe='Non-Catastrophe'`. **Bila `1=2` berlaku**, satu-satunya
kontrol yang bisa mengisi `StsKatastrofe` tak pernah tampil, gerbang hidup tetangganya tak pernah
bisa benar, dan **seluruh layar Catastrophe mati**. Maka `ALWAYS` menang atas `pyCondition`.

`[terverifikasi]` **Sidik jari residu** di berkas yang sama: ekspresi nyaris identik hidup di level
**container** (`pyIsVisibilityOption=CONDITION`) sekaligus yatim di level **sel**
(`pyVisible=ALWAYS`). Pola "gerbang dipindah, salinan lama tertinggal".

`[terverifikasi]` **Bukti bantahan yang dicari tetapi tidak ditemukan:** nol blok `OTHER` berkondisi
kosong; nol `pyVisible` berupa ekspresi; nol tag alternatif (`pyVisibleWhen`, `pyWhenName`).

### 13b. Batas pengetahuan aturan ini

⚠️ `[terbuka]` R1–R4 disimpulkan dari **korelasi + penalaran perilaku**, **bukan** dari kode mesin
Pega. `pyJavaStream` dan `pySectionXML` **kosong di seluruh 36 berkas**, sehingga tidak ada verifikasi
langsung terhadap markup keluaran.

⚠️ `[terbuka]` Semantik persis `NOTBLANK` (8×) · beda `true` versus `truewhn` (2×) ·
`ExpressionCondition` (1 sampel) · beda `pyReadOnly` kosong versus `false` (15×).

---

## 14. Section (36) — sensus isi

### 14a. ⚠️⚠️ **57 gerbang UI MATI** — 44 sel + 13 container

Ini padanan 84 gerbang mati di Activity, dan **belum pernah tercatat sama sekali**.

> ### ⚠️ RALAT 2026-09-17 — angka 57 direkonsiliasi, dan batas atasnya belum diketahui
>
> **Dari mana 57 datang:** ia **menjumlahkan dua keluarga gerbang yang berbeda** tanpa mengatakannya —
> **44 dari keluarga A** (sel/kontrol: `pyVisible` + `pyCondition`) **+ 13 dari keluarga B**
> (container: `pyIsVisibilityOption` + `pyContainerVisibleWhen`). Sensus pembanding hanya mencakup
> keluarga A, sehingga angkanya tidak bertemu.
>
> **Rekonsiliasi, dengan penyebut:**
>
> | Keluarga | Total elemen berkondisi | Berlaku | **Pasti MATI** | **Tidak terbaca** |
> | --- | ---: | ---: | ---: | ---: |
> | **A** — `pyCondition` | **180** | 119 (`pyVisible=OTHER`) | **44** (`pyVisible=ALWAYS`) | **17** (`pyVisible` kosong) |
> | **B** — `pyContainerVisibleWhen` | **56** | 43 (`CONDITION` 42 + `ExpressionCondition` 1) | **13** (`pyIsVisibilityOption=ALWAYS`) | **0** |
> | **Jumlah** | **236** | **162** | **57** | **17** |
>
> `[terverifikasi]` Keluarga B **tidak punya kategori "tidak terbaca"**: sensus 36 dari 36 Section
> memberi `pyIsVisibilityOption` = `ALWAYS` 230 · `CONDITION` 42 · `ExpressionCondition` 1 · **kosong
> 33**, dan `pyContainerVisibleWhen` terisi **56**. Karena 42 + 1 + 13 = 56, **ke-33 yang kosong tidak
> satu pun membawa ekspresi** — jadi tidak ada yang menggantung.
>
> **Kesimpulan:** **57 dari 236 pasti mati; batas atasnya 74 dari 236** (57 + 17), bukan 61. Angka 61
> adalah batas atas **keluarga A saja**.
>
> ~~⚠️ **Daftar triase Q1 (§21) belum boleh dipakai** — selisih batas bawah dan batas atas adalah
> **17 elemen**, dan arah 17 itu menentukan apakah gerbang bisnisnya hidup atau mati. Lihat **P23**.~~

> ### ⚠️ RALAT 2026-09-17 — kolom "tidak terbaca" DIBATALKAN. 57 adalah angka pasti.
>
> Ke-17 elemen ber-`pyVisible` kosong **tidak membawa syarat sama sekali**: `pyCondition`-nya berisi
> `"Other Property"` (15) atau `"When Rule"` (2) — label pilihan dropdown Designer, bukan ekspresi.
>
> | Keluarga | Elemen ber-`pyCondition` | Bukan syarat | Syarat sungguhan | Berlaku | **Mati** |
> | --- | ---: | ---: | ---: | ---: | ---: |
> | A | 180 | 17 | 163 | 119 | **44** |
> | B |  56 |  0 |  56 |  43 | **13** |
> | Jumlah | 236 | 17 | 219 | 162 | **57** |
>
> **Batas atas 74 dibatalkan.** Kolom "tidak terbaca" dihapus — nilainya nol di kedua keluarga.
>
> Pembanding, isi `pyCondition` yang benar-benar syarat: `pyVisible=OTHER` → `.IsKomite != 1`,
> `.AcceptanceStatus = 1`; `pyVisible=ALWAYS` → `isGolfInsurance || IsAneka || IsFire`, `.Type==1`.
> Aturan bacanya kini **R8** (§13).
>
> **Kalimat "daftar triase Q1 belum boleh dipakai" DICABUT** — Q1 sudah dijawab (§21) dan angkanya
> pasti.

**(i) Instruksi sembunyi-permanen yang tidak berlaku — 12.** Ditulis "jangan pernah tampil"
(`1=2`, `1=3`, `NEVER`), nyatanya **selalu tampil**: radio `StsKatastrofe` di `Catastrope_Sec` dan
`CatastrofeList_Sec` · `.ClaimData.PolicyData.PolicyNo` di `InputAcceptation` · `.ClaimData.IDMaster`
di `OutstandingClaim` · `.ClaimData.InsuredInterest` di `OutstandingClaim_Intrs` · lima tombol baris
grid ber-`NEVER`.

**(ii) Gerbang bisnis nyata yang hilang — 32.** Yang paling material:

| Section | Kontrol | Syarat yang **tidak berlaku** |
| --- | --- | --- |
| `AdjustmentDetail` | `.DirectToKasir` (checkbox) | `.IsKomite = '1'` |
| `AdjustmentDetail` | 6 sel `IndividualRisk*` + `.Currency` | `.Type==1` |
| `InputAcceptation` / `OutstandingClaim` | `.TreatyInMaster.Ceding` | `IsMarineCargo` |
| `InputAcceptation` / `OutstandingClaim` | `LeadingReinsSource`, `Bordeaux`, `BordereauxNote`, `TeritorialScope` (4 sel × 2) | `isGolfInsurance \|\| IsAneka \|\| IsFire` |
| `InputAcceptation` / `OutstandingClaim` | `.TreatyInMaster.AccountingMode` | `ProportionalType=='NonProportional'` |
| `InputAcceptation` | blok Standard | `.ClaimData.IsCloseFile = true` |
| `*_Est` | `.ClaimData.TSIDeductible` (2 sel masing-masing) | `.ClaimData.TypeDeductible` |
| `InputAcceptation_Est`, `OutstandingClaim` | tombol PLA | `.ClaimData.IsPLA != 1` |
| `MstAdjusterConsultant` | `OutputParam.ERRMSG6`, `OutputParam.ERRMSG` | `OutputParam.ERRMSG != ''` — ⚠️ **pesan error selalu tampil, termasuk saat kosong** |

**(iii) Container mati — 13, dan 7 di antaranya panel "Spreading".** `[terverifikasi]` Seluruh skema
sembunyi-Spreading di sisi Adjustment/Estimation **tidak berfungsi** — panel `Spreading List`,
`Spreading In`, `Spreading Out` selalu tampil meski ditulis `NEVER` atau bersyarat
`IsEditEstimation=='true'`. Juga `MasterTreatyInList` (`Title` tampil walau hasil kosong) dan
`MstAdjusterConsultant` (`IsFire`).

### 14b. Hierarki data — **45 pasangan page-list, dan enam Spreading berbagi satu class**

`[terverifikasi]` Bukti definisional lengkap. Akar pertama `pyWorkPage.ClaimData`:

```
ClaimData ─┬─ EstimationList          ASM-FW-GCNMFW-Data-Estimasi
           ├─ InterestList            ASM-FW-GCNMFW-Data-ObjectItem
           ├─ InterestListDtl         ASM-FW-GCNMFW-Data-ObjectItem
           ├─ AdjustmentList          ASM-FW-GCNMFW-Data-Adjustment
           ├─ ListTotalEstimation     ASM-FW-GISFW-Data-TreatyInTotal
           ├─ ListClaimAmount         ASM-FW-GISFW-Data-TreatyInTotal
           ├─ TotalInterestInsured    ASM-FW-GISFW-Data-TreatyInTotal
           ├─ SuggestList             ASM-FW-GISFW-Data-OfferFacIn-SuggestList  ⬅ ditumpangi
           └─ SpreadingRisk · SpreadingClaim · SpreadingAdjustment ·
              SpreadingAdjustmentQS · SpreadingBreakQS · SpreadingQuotaShare
                                      ASM-FW-GISFW-Data-SpreadingRisk  ⬅ ENAM properti, SATU class
```

⚠️ `[terverifikasi]` **Akar KEDUA yang belum pernah tercatat: `pyWorkPage.PaymentData`** — hierarki
tiga tingkat, terbukti dari dua Section berbeda:

```
PaymentData ─┬─ AcceptationList   ASM-FW-GCNMFW-Data-Acceptance
             └─ CurrencyList      ASM-FW-GCNMFW-Data-Currency
                   └─ DetailPayment   ASM-FW-GCNMFW-Data-Payment
```

Ditambah `.ComiteeClaim` (`ASM-FW-GCNMFW-Data-Comitee`) dan `.LossAllocation` di `AdjustmentDetail`
(class `Data-SpreadingRisk` — menguatkan §8 Ronde 1).

### 14c. ⚠️ **Nol validasi di lapisan Section**

`[terverifikasi]` `pyValidate`, `pyEditValidate`, `pyMaxLength`, `pyMinValue`, `pyMaxValue`,
`pyValidationMsg`, `pyErrorMessage`, `pyCustomValidation` — **semuanya nol berisi** di seluruh 36
berkas. **Tidak ada validasi format, panjang, maupun rentang di UI.** Pemeriksaan nilai hanya terjadi
lewat activity yang dipicu `change` (`CheckDateDOL_Act`, `CheckEstimateDate_Act`, dll.) — di luar
Section.

`[terverifikasi]` Field wajib: **43 tanpa syarat** + **9 bersyarat** (semuanya hidup). Yang bersyarat
antara lain `.DataCommitteeTreaty.Occupation` (`.Type = 1`), `.AdjusterFee` (`.Type = 2 || 4`),
`.Salvage` (`.Type = 3`) — **selaras dengan bukti `PaymentType` di §12c/§16a Ronde 1**.

`[terverifikasi]` Read-only bersyarat **128, semuanya hidup**; pendorong terbesar
`pyWorkPage.IsOutstanding==1` (47×). Disable bersyarat **45, semuanya hidup**.
⚠️ Anomali: `pyReadOnlyCondition = 1=1` (2×) — read-only tanpa syarat ditulis sebagai ekspresi
selalu-benar; karena R5 tanpa flag pengaktif, ini **berlaku**.

### 14d. ✅ Nol pemeriksaan peran — **kesimpulan Ronde 1 TERVERIFIKASI**

`[terverifikasi]` `pyPrivilege` **0 berisi / 199 kosong** · `pyPrivilegeView` **0 / 392** ·
`pyPrivilegeUpdate` **0 / 392** · `pyWhenNotPrivilege` **0 / 199** · `pyAssociatedPrivileges` 278
**semuanya `false`** · tag `pyAccessGroup`/`pyRole`/`HasRole`/`pxAccess` **tidak ada di korpus**.

**Nol dari 418 ekspresi gerbang UI** menyebut operator, peran, access group, atau privilege. Seluruh
gerbang didasarkan pada **state data kasus** (`.IsKomite`, `.Type`, `IsOutstanding`,
`IsAnyAcceptation`, `IsPLA`, `StsKatastrofe`, `IsOldData`), **jenis bisnis** (`IsMarineCargo`,
`isGolfInsurance`, `IsAneka`, `IsFire`), dan **parameter layar**.

Yang **mirip tetapi bukan** pemeriksaan peran: `pxRequestor.pyIsMobile` (cek perangkat) ·
`OperatorID.pyUserName`/`pyUserIdentifier` sebagai **sel tampilan** · `D_pyUserWorkList` (filter
worklist bawaan Pega).

### 14e. Lain-lain

`[terverifikasi]` **11 dari 36 Section adalah save-as** — mis. `InputAcceptation` ← `INPUTESTIMASIADMIN`,
`OutstandingClaim` ← `VIEWPOLIS`, `Subjectivity` ← `DTLDATACOMMITTE`, `CauseofLoss_Section` ← class
`Work-PNC`. **Ini menjelaskan asal-usul residu gerbang mati.**

`[terverifikasi]` **133 pemicu activity** dari Section. Melengkapi peta OQ-039: ~~karena **R7**
(`pyActionConditions` kosong di seluruh 36), satu-satunya cara menonaktifkan tombol adalah lewat
`pyDisabledWhen` atau gerbang tampil.~~ Tombol `Send to Committe` disabled bila
`.IsKomite==1 || pyWorkPage.IsError>1`; tombol Kasir disabled bila
`.DirectToKasir='false' || .StatusKasir != ''`.

> ### ⚠️ RALAT 2026-09-17 — kalimat bersandar R7 yang dibatalkan
>
> Klaim "satu-satunya cara menonaktifkan tombol adalah `pyDisabledWhen` atau gerbang tampil"
> **gugur**. Menurut **R7-baru**, ada **cara ketiga**: `pyActionConditions` — terisi di **3 dari 36
> Section**, **17 baris** (`Section/AdjustmentDetail.xml` 10, `Section/OutstandingClaim.xml` 6,
> `Section/InputAcceptation_Est.xml` 1).
>
> Sisa kalimatnya **tetap berlaku**: kedua gerbang `pyDisabledWhen` yang dikutip memang hidup
> menurut **R3**.

---

## 16. ReportDefinition (19), DecisionTable (1), SystemSettings (1)

### 16a. ⚠️ Aturan baca tambahan: `pyRDName` versus `pyGridRDName`

`[terverifikasi]` Sumber data grid yang **benar-benar aktif** adalah `pyRDName` pada
`Embed-Display-Table-Restricted-Grid`. **`pyGridRDName` hanyalah stempel per-sel** yang bisa
tertinggal basi ketika sumber grid diganti. Sumber autocomplete adalah `pySourceName`. Pemanggilan
dari Activity hanya sah bila ada step `Call pxShowReport` / `Call pxRetrieveReportData`.

**Tanpa aturan ini, empat RD mati akan terbaca sebagai hidup.**

### 16b. Empat RD yatim

| RD | Bukti mati |
| --- | --- |
| **`BrowseRW_RD`** | satu-satunya rujukan adalah `Property-Set` di `GetAdders_Act` **step 2** — tetapi activity itu **tidak punya step pemanggil RD sama sekali**; data diambil `RDB-List` step 3. ⚠️ Berisiko bila "dihidupkan": `pyMaxRecords=10000` |
| **`BrowseTreatyInDetail`** | 15 stempel `pyGridRDName` basi di `MasterTreatyInList`; grid aktifnya `BrowseCLAIM_MASTER_TREATY`. **RD terbesar** (33 kolom, 12 kondisi, 13 param) tetapi tak terpanggil |
| **`BrowseTreatyGroup_RD`** | 1 stempel basi di `MstAdjusterConsultant`; grid aktifnya `BrowseAdjusterConsultant` |
| **`BrowseVMstUserTeknis_RD`** | 1 stempel basi; **tanpa filter apa pun** dan param `type_business` tak terpakai |

**Semi-yatim** `[terverifikasi]`: `RejectedClaim_RD` kondisi `B` (`.ID = Param.ClaimID`) — parameter
**tidak pernah diisi** pemanggil tunggalnya (`CheckDateDOL_Act` step 21 hanya mengisi `INSKEY`).
`BrowseReinsuranceType_RD` dipanggil dua activity (`SetNameTreaty_Act` step 2,
`SetTreatyNameSpreading_Act` step 3) **tanpa satu pun parameter filter diisi**.

### 16c. Koreksi dan penguatan

`[terverifikasi]` **`FilterEmailKomiteWithLimit` punya 18 kolom, bukan 17** — `.JABATAN` ikut diambil,
dan hasilnya **diurut ASC pada `.DEGREE`**. Fakta 2 mencatat 17 tanpa `JABATAN`.

⚠️ **Ini relevan untuk Q5 (§21):** kolom `JABATAN` **sudah diambil RD** — jadi mengganti kunci pencarian
limit dari nama orang ke jabatan **tidak menuntut perubahan skema**.

`[terverifikasi]` **`LIMIT_TOP` TERVERIFIKASI tidak pernah dipakai sebagai filter** — 11 kemunculan,
**seluruhnya di satu berkas**, dan hanya sebagai **kolom keluaran**. Nol kemunculan di seluruh jenis
rule lain. Batas atas limit komite **tidak pernah menjadi kriteria seleksi**.

`[terverifikasi]` **`BrowseVDCauseOfLoss_RD` memuat `INNER JOIN`** ke class
`ASM-FW-GCNMFW-Int-V_M_CAUSE_OF_LOSS` dengan kondisi `.M_COL_ID = MCause.M_COL_ID` — **bukti
independen** hubungan dua tingkat master (§15a Ronde 1).

⚠️ `[terverifikasi]` **`.STS_AKTIF` diambil sebagai kolom tetapi TIDAK difilter** di
`BrowseBankAccount`, `BrowseCouseOfLoss_Business`, dan `BrowseVDCauseOfLoss_RD` → **baris non-aktif
ikut terambil**. Bandingkan `BrowseRW_RD` dan `FilterEmailKomiteWithLimit` yang memfilter `= "1"`.

`[terverifikasi]` **Hardcode di dalam RD:** `STS_AKTIF="1"` (2 RD) · `KLAIMTYPE="NON-LIFE"`
(`GetCatastrope_RD`). **Disuntik pemanggil:** `"PROP"` · `"INDONESIA"` · `"Proportional"` · `4`.

`[terverifikasi]` **Nol filter berbasis peran/jabatan/operator** di seluruh 19 RD;
`pyPrivilegeList` **kosong di semuanya**. `.DEGREE` dan `.JABATAN` hanya kolom keluaran.
**Kesimpulan Ronde 1 terverifikasi di lapis ketiga.**

### 16d. DecisionTable — satu-satunya di modul: `GetMimeType`

`[terverifikasi]` `ASM-FW-GISFW-INT-T_STORAGE_IMAGE!GETMIMETYPE` / `RULE-DECLARE-DECISIONTABLE`.
Satu kolom masukan `UploadDoc.ext`, **42 baris** ekstensi → MIME, *otherwise*
`"application/octet-stream"`. `pyEvaluateAllRows = no`. Kolom property-set ada tetapi **seluruh 42
nilainya kosong** → tanpa efek samping.

**Pemanggil tunggal:** `InsertGoogleStorage_Act` **step 4** (`Property-Map-DecisionTable`,
deskripsi *"Get Mimetype"*, step page `UploadDoc`, **tanpa precondition**).

⚠️⚠️ **Peran bisnisnya menurut ISI, bukan nama:** ini **satu-satunya titik klasifikasi tipe berkas di
jalur unggah**, dan ia **tidak memblokir apa pun**. Ekstensi tak dikenal **tidak ditolak** — ia jatuh
ke `application/octet-stream` dan **tetap diunggah**. Tidak ada penyaringan ekstensi berbahaya, dan
`html`/`svg`/`url`/`lnk` justru dipetakan ke MIME yang dapat dirender browser. Keputusan sepenuhnya
bergantung pada **nama ekstensi yang dikirim klien**, bukan pemeriksaan isi berkas.

⚠️ Kopling: kolom masukan ditulis sebagai referensi page **absolut** `UploadDoc.ext`, bukan `.ext`
relatif → tabel hanya dapat dipakai ulang bila pemanggil kebetulan memakai nama page `UploadDoc`.

### 16e. ⚠️⚠️ SystemSettings `LinkService` — **tidak memuat URL sama sekali**

**Ini menggeser OQ-018 satu lapis.** `[terverifikasi]` `LINKSERVICE!LINKSERVICE` /
`RULE-ADMIN-SYSTEM-SETTINGS`, `pyClassName` **kosong** (rule tingkat sistem).

Nilainya **bukan literal** melainkan **ekspresi** `=ResponLink.URL`, dan **kelima production level
bernilai identik**:

| Level | Nilai |
| --- | --- |
| 1 Sandbox · 2 Development · 3 QA · 4 Staging · 5 Production | `=ResponLink.URL` (semua sama) |

⚠️ **Konsekuensi:** mekanisme pembedaan lingkungan bawaan Pega (*production level*) **tidak dipakai
sama sekali**. Pemisahan Dev/QA/Staging/Prod **sepenuhnya bergantung pada isi database**.

`[terverifikasi]` Sumber nilai sesungguhnya:
`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE` — step 1 `Page-New ResponLink`; step 2 `Obj-Browse`
ke `ASM-FW-GISFW-Int-M_LINK_SERVICE` dengan kriteria `.KATEGORI_1` dan `.KATEGORI_2`; step 3
`ResponLink.URL = linkService.pxResults(1).URL`; step 4 `Page-Remove`.

**Jadi base URL seluruh integrasi keluar ditentukan oleh baris tabel `M_LINK_SERVICE` di Oracle**,
dikunci sepasang kategori, dan **selalu memakai baris pertama tanpa penanganan bila kosong**.
Ini **menguatkan ADR-0013**.

⚠️⚠️ `[terverifikasi]` **Hanya 2 dari 6 Connect-REST yang terbukti memanggil `GetLinkService` lebih
dulu** (`GetDtlPaymentClaimTreatyin_Act`, `GetDtlPaymentPremi_act`). Untuk empat lainnya —
**termasuk `ServiceGoogle` yang dipakai unggah berkas** — **tidak ditemukan pemanggilan
`GetLinkService` di dalam `Claim Prop`**. Base URL berpotensi **kosong saat runtime**, atau diisi di
luar batas korpus ini → `[terbuka]`.

`[terverifikasi]` Nol setting lain di berkas ini; folder `SystemSettings/` hanya berisi berkas ini.
Nol `pxGetSystemSetting` / `@getSetting(...)` / `D_pxSystemSettings` di seluruh modul.

---

## 17. FlowAction (12), Harness (11), When (6), ConnectREST (6)

### 17a. FlowAction — **nol gerbang blocking di lapisan ini**

⚠️ `[terverifikasi]` **Tag `pyPostActivity` TIDAK ADA sama sekali** di schema ekspor ini. Slot hook
yang benar-benar ada:

| Slot | Terisi |
| --- | ---: |
| `pyPreProcessingActivity` | 2 |
| `pyPreProcessingTransformRule` | 3 |
| **`pyValidateActivity`** | **0 dari 12** |
| `pyActionTransformRule` | 0 |
| `pyLocalActionActivity` | 3 |

⚠️ `pyPreActivity` memang muncul 3×/berkas tetapi **selalu kosong** dan berada di sub-page
layout/grid — **bukan** hook flow action.

**Klaim Ronde 1 BENAR dan lebih kuat:** nol `pyValidateActivity`, nol Rule-Obj-Validate.
Gerbang transisi yang sesungguhnya adalah **pre-processing**, yang jalan **sebelum form tampil** —
jadi **tidak bisa menolak submit**. Satu-satunya penghalang adalah `pyClientValidation` di browser,
yang dapat dilewati. **Tidak ada gerbang blocking sisi server di lapisan FlowAction.**

`[terverifikasi]` Hook yang terisi: `InputAcceptation` → `GetPICAdjutment_Act` + `SetDateAcceptation` ·
`OutstandingClaim` → `CheeckNoRNM_Act` + `SetDateOutstanding` + local action `ProteksiData_act` ·
`GenerateDLATreaty` → `PrintDLATreatyIn` · `GeneratePLA` → `TryMakePLA_Act` ·
`ViewDetailInterest` → `CekInterestListDtl_DT`.

`[terverifikasi]` **Nol FlowAction yatim.** ⚠️ Asimetri: `AdjustmentDetail` hanya punya **satu** titik
rujuk (`pyEditAction` di `InputAcceptation_Adjs`) — **jalur Outstanding tidak memanggilnya sama
sekali**, padahal jalur Acceptance memanggilnya.

### 17b. Harness — gerbangnya bukan milik Harness

⚠️ `[terverifikasi]` **Dari 9 `pyContainerVisibleWhen` di seluruh 11 berkas, NOL berada di badan rule
Harness** — semuanya di rekaman **Section yang dibundel** dalam berkas yang sama. Harness di modul
ini praktis **tidak menggerbangi apa pun sendiri**.

`[terverifikasi]` Gerbang komite `Protect.CARI1 =1 && Protect.CARI2= 1` **terverifikasi**, dengan
`pyIsVisibilityOption=CONDITION` dan **`pyIsClientWhen=true`**.

⚠️ `[terverifikasi]` Tombol `Send Claim to Committee` punya **lapis ketiga**:
`pyDisabledWhen = pyWorkPage.ClaimData.Payable = ''` dengan `pyDisabledNew=true` (gerbang **hidup**
menurut R3).

⚠️⚠️ **`IsFire` dirujuk sebagai `pyContainerVisibleWhen` di `MstAdjusterConsultant`, tetapi rule
`IsFire` TIDAK ADA di folder `When/`** → dependensi luar korpus, tidak dapat diverifikasi →
`[terbuka]`. ~~(Ia juga dirujuk di `PrintFileAcceptance` step 11 — §10e.)~~

> ### ⚠️ RALAT 2026-09-17 — **P19: bukan dua rujukan, melainkan ENAM berkas / 12 kemunculan**
>
> `[terverifikasi]` **Klaim ketiadaan, dengan cakupan:** dicari di **6 dari 6 berkas** folder
> `Claim Prop/When/` — isinya `IsBackStage.xml`, `IsCLM.xml`, `IsCLMNP.xml`, `IsCLMP.xml`,
> `IsPEGAPROD.xml`, `IsPEGASyariah.xml`. **`When/IsFire.xml` tidak ada.**
>
> `[terverifikasi]` Perujuk, dicari di **270 dari 270 berkas** `Claim Prop`: **6 berkas, 12
> kemunculan**.
>
> | Berkas (path lengkap dari akar modul) | Kali | Yang digerbangi | Status gerbang |
> | --- | ---: | --- | --- |
> | `Activity/PrintDLATreatyIn.xml` | 1 | step **`RH_1.pySteps(15).pySteps(16)`** | **MATI** — `pyStepsPreCondition = 0`; arah `WhenTrue` **kosong** / `WhenFalse = 3` |
> | `Activity/PrintFileAcceptance.xml` | 1 | step **`RH_1.pySteps(12)`** — set `param.HTMLStream="FILEAcceptanceNote"`, `AttachmentCategory="AcceptanceNote"` | **MATI** — `pyStepsPreCondition = 0`; arah `WhenTrue` **kosong** / `WhenFalse = 3` |
> | `Section/InputAcceptation.xml` | **4** | 4 sel ber-`pyCondition = isGolfInsurance \|\| IsAneka \|\| IsFire` | **MATI** — keempatnya `pyVisible = ALWAYS` |
> | `Section/OutstandingClaim.xml` | **4** | idem, 4 sel | **MATI** — keempatnya `pyVisible = ALWAYS` |
> | `Section/MstAdjusterConsultant.xml` | 1 | container `pyContainerVisibleWhen = IsFire` | **MATI** menurut sensus keluarga B (`pyIsVisibilityOption = ALWAYS`) — ⚠️ atribusi tag opsi tidak saya konfirmasi ulang sendiri → `[terbuka]` ringan |
> | `Harness/MstAdjusterConsultant.xml` | 1 | salinan baris di atas, pada **rekaman Section yang dibundel** | idem |
>
> ⚠️ Ralat kecil yang ikut: rujukan di `Activity/PrintFileAcceptance.xml` **bukan step 11** melainkan
> **step `RH_1.pySteps(12)`** — §10e memakai nomor yang salah.
>
> **Akibat bila rule `IsFire` memang tidak pernah ada:** `[terverifikasi]` **nihil di runtime** —
> **kedua belas** kemunculan itu berada di gerbang yang **sudah mati lebih dulu**
> (`pyStepsPreCondition = 0` atau `pyVisible = ALWAYS`). Ekspresinya **tidak pernah dievaluasi**,
> jadi ada-tidaknya rule tidak mengubah perilaku.
>
> ⚠️ `[terbuka]` **Bagaimana Pega memperlakukan `when` yang hilang — benar atau salah — TIDAK terbaca
> dari korpus.** Tidak ada rule, log, atau pesan di 270 dari 270 berkas yang menunjukkannya.
> **Jangan ditebak.** Pertanyaan ini baru menggigit bila salah satu dari 12 gerbang itu dihidupkan.
>
> **Konsekuensi migrasi:** karena seluruh pemakaiannya mati, `IsFire` **tidak perlu dimigrasikan** —
> tetapi **niat** yang tertulis (membedakan lini Fire) ikut hilang, dan itu masuk **daftar triase
> Q1 (§21)**, bukan diterapkan diam-diam.

`[terverifikasi]` Nol pemeriksaan peran/privilege di seluruh 11 harness — **lapis keempat** yang
menguatkan §14d.

⚠️ `[terverifikasi]` Tiga activity kritikal (`SaveAdjusterConsultant_Act`, `CancelTreatyGroup`,
`NewAdjustConsult_Act`) dipanggil dengan `pyActivityClass=@baseclass` — resolusi rule melebar ke
seluruh hierarki.

`[terverifikasi]` `TambahCauseofLoss` dan `TambahMasterCauseOfLoss` membundel rekaman ketiga
`@BASECLASS!PYDASHBOARDMYWORKLIST` dan membawa `pyUsage` berbunyi *"thin wrapper around the My Cases
display in the Case Manager portal…"* → **save-as dari harness standar `pyDashboardMyWorklist`**
dengan `pyUsage` yang tidak pernah dibersihkan.

### 17c. ⚠️⚠️ KONTRADIKSI antar-penelusuran — `pyActionConditions`

**Dua penyisiran memberi hasil berlawanan, dan saya tidak memilih salah satunya:**

| Sumber | Klaim |
| --- | --- |
| Penyisiran Section (36 berkas) | `pyActionConditions` **kosong di seluruh 36 berkas** → dasar **R7** |
| Penyisiran Harness/FlowAction | `Section/AdjustmentDetail.xml` punya `pyActionConditions` **dua baris**: `Protect.CARI1 = "1"` **AND** `Protect.CARI2 = "1"`, menggerbangi pembukaan harness `CommitteeTreaty` |

⚠️ **Konsekuensi bila penyisiran kedua benar: R7 gugur**, dan gerbang tombol memang ada.
Konsekuensi bila yang pertama benar: gerbang komite hanya satu lapis, bukan dua.

**Ini harus dituntaskan sebelum `/to-spec`** — ia menentukan apakah ada gerbang di lapisan aksi
tombol sama sekali. → **P18 (BARU)**.

> ### ⚠️ RALAT 2026-09-17 — **P18 TERTUTUP. Kedua penyisiran saya SALAH.**
>
> Bukan salah satu yang benar — **keduanya keliru**:
>
> | Penyisiran | Klaimnya | Kenyataan |
> | --- | --- | --- |
> | Section (36 berkas) | `pyActionConditions` **kosong di seluruh 36** | **salah** — 3 Section berisi, **17 baris** |
> | Harness/FlowAction | "`AdjustmentDetail` punya **2 baris**" | **salah dua kali** — jumlahnya **10**, dan **tidak ada Harness bernama `AdjustmentDetail`**; berkas itu **`Section/AdjustmentDetail.xml`** |
>
> **Sensus yang benar** (menghitung anak tiap elemen `pyActionConditions`, 270 dari 270 berkas):
> `Section/AdjustmentDetail.xml` **10** · `Section/OutstandingClaim.xml` **6** ·
> `Section/InputAcceptation_Est.xml` **1** · `Section/AdjustmentDetail_Section.xml` **0** ·
> `FlowAction/AdjustmentDetail.xml` **0** · tidak ada berkas lain.
>
> **Penyebab: tabrakan nama** — tiga berkas berbeda bernama `AdjustmentDetail` di folder berbeda.
> Aturan kerja pencegahnya sudah ditulis di §13 (path lengkap dari akar modul; cek nama kembar
> sebelum menyimpulkan).
>
> **Akibat pada gerbang komite:** gerbangnya **tetap dua lapis** —
> `pyActionConditions` di `Section/AdjustmentDetail.xml` **dan** `pyContainerVisibleWhen` di Section
> yang dibundel `Harness/CommitteeTreaty.xml` — tetapi **isi persis 10 baris itu dan cara
> menggabungkannya belum diverifikasi ulang** → `[terbuka]`, bukan bagian dari P18 yang sudah tertutup.
>
> **Yang TIDAK berubah:** keduanya tetap **client-side** (`pyIsClientWhen=true`), dan **nol
> `pyValidateActivity` di 12 dari 12 FlowAction** tetap berlaku — jadi kesimpulan "tidak ada gerbang
> blocking sisi server" **tetap berdiri**.

### 17d. When — keenam **TEPAT**, dan `A OR B` terkonfirmasi

`[terverifikasi]` Keenam kondisi yang dicatat Ronde 1 **tidak ada yang meleset**.

⚠️ `[terverifikasi]` **Penggabungan `IsCLM*` adalah `OR`**, bukan `AND` — penanda otoritatif
`<pyLogic>A OR B</pyLogic>` di ketiga berkas.
**Jebakan:** sub-page `<pyConditionViewer>` menyimpan `<pyLogicalOperator>AND</pyLogicalOperator>`,
tetapi itu **nilai default sisa editor** — `IsPEGAPROD` dan `IsBackStage` yang hanya punya satu
kondisi pun menyimpannya.

**Akibat fungsional:** `IsCLM` true bila prefix ada di **cover ATAU work page** — inilah yang membuat
`IsCLM`/`IsCLMP`/`IsCLMNP` **bisa true bersamaan** pada struktur cover/child.

`[terverifikasi]` Rantai save-as yang belum pernah tercatat: **`IsCLM` → `IsCLMP` → `IsCLMNP`**.

`[terverifikasi]` `IsBackStage` dipakai sebagai **konektor flow `Transition12`** dari gateway
`Decision3`, bukan sebagai step.

### 17e. ⚠️ Empat gerbang `When` mati — **tiga belum pernah dilaporkan**

| Rule | Step | When | Efek |
| --- | --- | --- | --- |
| `InsertGoogleStorage_Act` | 6 | `IsPEGAPROD` (`false`) | sudah diketahui |
| **`HitServiceToKasir_Act`** | **9.9** | `IsCLMNP` (**tag kosong**) | *"Insert Log Direct to Kasir"* **jalan untuk SEMUA prefix** |
| **`HitServiceToKasir_Act`** | **13.6** | `IsCLMNP` (**tag kosong**) | idem |
| **`HitServiceToKasir_Act`** | **13.1.5** | `IsCLMP` (**tag kosong**) | *"Set Tanggal TglBolehBayar"* jalan tanpa syarat prefix |

`[terverifikasi]` Nol perujuk When yang di-remark; nol When yatim.

### 17f. ConnectREST — enam connector, satu setting

`[terverifikasi]` Keenam: **nol resource path** (`pyQueryStringParameters` dan
`pyResourcePathParameters` kosong), semua **POST**, semua lewat `LinkService!LinkService`.
Hanya **`SendAcceptationToKasir`** memakai autentikasi (profil autentikasi bernama — namanya
sengaja tidak dicatat); lima lainnya
`pyUseAuthentication=false`.

⚠️ `[terverifikasi]` **`pyResponseTimeout` tidak konsisten** untuk enam connector yang menunjuk
setting endpoint **yang sama**:

| Nilai | Connector |
| --- | --- |
| **`0`** (tanpa batas) | `GetDtlPaymentClaim`, `getPremiumPaidOnTreatyIn` |
| `30000` | `SendAcceptationToKasir`, `getPayAttachment` |
| `300000` | `KonversiKlaimNonLife`, `ServiceGoogle` |

⚠️ Semua call-site `ExecutionMode=Run` — **sinkron dan memblokir**. Dua connector ber-timeout `0` di
jalur UI = **panggilan sinkron tanpa batas waktu**.

`[terverifikasi]` **Lima dari enam responsnya DIBACA**; hanya `KonversiKlaimNonLife` yang
fire-and-forget — **konfirmasi Ronde 1**.

⚠️⚠️ `[terverifikasi]` **`SendAcceptationToKasir` membaca properti salah eja**:
`Primary.StatusServiceKasir.ReponseCode` (kurang huruf `s`). Bila properti sebenarnya `ResponseCode`,
ekspresi `@if(...=1, "Akseptasi Sudah Masuk ke Kasir", ...ResponseMsg)` **selalu jatuh ke cabang
else** → pesan sukses **tidak pernah muncul**, dan `.StatusKasir` selalu berisi `ResponseMsg`.
~~**De facto fire-and-forget.** Ini juga menjelaskan mengapa `CloseClaimProp` memblokir close dengan
*"direct to cashier that has not been successful"*.~~

> ### ⚠️ RALAT 2026-09-17 — **P20 diturunkan jadi cacat KOSMETIK, keluar dari daftar pemblokir**
>
> **Temuannya benar, akibatnya yang saya lebih-lebihkan.**
>
> `[terverifikasi]` Ejaan, dicari di **270 dari 270 berkas**: **`ResponseCode` nol kemunculan**;
> `ResponseMsg` dieja benar. **Satu ekspresi memuat dua ejaan sekaligus.**
>
> `[terverifikasi]` Pemakaiannya **hanya dua tempat, keduanya di berkas yang sama dan keduanya
> `Property-Set` tanpa `when`**:
>
> | Berkas | Step |
> | --- | --- |
> | `Activity/HitServiceToKasir_Act.xml` | **9.8** |
> | `Activity/HitServiceToKasir_Act.xml` | **13.5** |
>
> `[terverifikasi]` **Tidak ada satu pun step yang bercabang pada nilai itu.** Yang terdampak hanya
> **teks status yang tampil**, bukan alur. **Tidak ada transaksi yang gagal karenanya**, dan klaim
> "de facto fire-and-forget" **gugur** — respons tetap dipetakan dan `ResponseMsg` tetap terbaca.
>
> Klaim bahwa ini "menjelaskan mengapa `CloseClaimProp` memblokir close" juga **gugur**: gerbang itu
> menguji `.StatusKasir == "Akseptasi Sudah Masuk ke Kasir"`, dan apakah nilai itu pernah terbentuk
> bergantung pada nama properti di **respons JSON Kasir**, bukan pada ejaan di ekspresi.
>
> ⚠️ `[terbuka]` **Nama properti respons tidak terbaca dari korpus.** Diperiksa langsung di
> `ConnectREST/SendAcceptationToKasir.xml`: pemetaan responsnya adalah `pyMapTo = JSON` ke
> **`pyMapToKey = .StatusServiceKasir`** — yaitu **seluruh halaman**, tanpa menyebut sub-properti.
> **Baik `ReponseCode` maupun `ResponseCode` tidak muncul sama sekali di berkas connector itu.**
> Nama sub-properti lahir dari struktur JSON balasan Kasir saat runtime, **di luar korpus**.
> Jadi **tidak dapat ditunjukkan** bahwa pemetaan responsnya menyasar nama yang salah eja —
> **dan tidak saya simpulkan**.
>
> **Status: cacat kosmetik. P20 keluar dari daftar pemblokir `/to-spec`.**

⚠️ `[terverifikasi]` `GetDtlPaymentClaim` dan `getPremiumPaidOnTreatyIn` memakai **`.PaymentData`
sebagai sumber request DAN target respons** — halaman ditimpa oleh hasilnya sendiri.

### 17g. ⚠️ KOREKSI §8f: **ada nomor polis hardcoded** — di Activity, bukan SQL

`[terverifikasi]` `GetDtlPaymentPremi_act` **step 6**, precondition **mati** (`false`):
`pyWorkPage.OfferFacIn.PolicyData.PolicyNo == "<satu nomor polis>"` — **sisa debugging**. Karena
gerbangnya mati, connector `getPremiumPaidOnTreatyIn` **selalu jalan** untuk semua polis.

§8f benar bahwa **nol nomor polis di dalam SQL**; yang keliru adalah menganggapnya nol di seluruh
modul.

`[terverifikasi]` `getPayAttachment` — gerbang `.IsPayment==1` juga **mati** di
`GetPayAttachmentAdj_Act` step 3.2.

⚠️⚠️ `[terverifikasi]` Di `GetPayAttachmentAdj_Act`, step **Obj-Save** (3.3.8) dan **Commit** (3.3.9)
**DI-REMARK**, sementara pemanggilan connector-nya tetap jalan tanpa gerbang. **Lampiran pembayaran
diambil tetapi tidak pernah disimpan.**

⚠️ `[terverifikasi]` `ServiceGoogle` dipanggil dari dua tempat, **keduanya tanpa gerbang sama
sekali** (`InsertGoogleStorage_Act` step 11, `GetUrlGoogleStorage_Act` step 6.5) — menguatkan §10c
dan **Q4** (§21).

---

## 18. Status frontier Ronde 2

**Frontier fakta: KOSONG.** Seluruh 270 berkas dan 12 jenis rule sudah disisir; tujuh jenis yang
sebelumnya nol tersentuh kini tercatat.

~~**Frontier keputusan: Q1–Q7 menunggu jawaban work owner — definisinya di §21.**~~

> ### ⚠️ RALAT 2026-09-17 — **frontier keputusan KOSONG**
>
> **Q1–Q7 terjawab seluruhnya** pada 2026-09-17 `[keputusan work owner]` — jawaban tertulis di
> **§21**, di bawah masing-masing `### Q<n>`, dan tiga risiko yang diterima sadar di **§21c**.
>
> **Frontier keputusan: KOSONG.** Yang tersisa hanya **frontier fakta** — lihat §19.

> ### ⚠️ RALAT 2026-09-17 (kedua) — frontier fakta juga hampir kosong
>
> Setelah lima jawaban pemblokir (U1) dan korpus bertambah menjadi **328 berkas**:
>
> | Frontier | Status |
> | --- | --- |
> | **Keputusan** (`Q1`–`Q7`) | ✅ **KOSONG** — terjawab seluruhnya, §21 |
> | **Fakta** (`P1`–`P24`) | 🟡 **tersisa satu** — **P24 butir 2–3** saja |
>
> P15 · P19 · P21 · P23 ✅ tertutup (U1). P2 · P14 · P16 · P17 · P18 · P20 · P22 sudah tertutup lebih
> dulu. **Yang tersisa: apakah ketiga daftar wewenang (§11e-3) satu tangga yang sama.**
>
> Dua `[terbuka]` ringan yang **tidak** memblokir: arti niat step 17 `CountValueADJTreaty_Act`
> (§24b) dan salah ketik `IsCustomBonds` versus `IsCustomBond` (§27d).

## 19. Pemblokir — diperbarui

| # | Pertanyaan | Pemilik | Status |
| --- | --- | --- | --- |
| ~~P2~~ | ~~`.ClaimSpreaded`: basis pengali — baris adjustment (step 6.2) atau `Hepeng.pxResults` (step 9.2.1)?~~ | — | ✅ **GUGUR 2026-09-17** — bukan kontradiksi, melainkan **dua daftar berbeda**: step **6** melingkupi `.SpreadingAdjustment` pada **halaman baris `AdjustmentList`**, step **9.2** melingkupi `pyWorkPage.ClaimData.SpreadingAdjustment` (**tingkat halaman**). Terbukti dari step **2** yang menyalin `ClaimData.SpreadingClaim` ke tingkat halaman. Keduanya sah dan tidak saling meniadakan |
| **P15** | **MENYEMPIT 2026-09-17.** Pertanyaan tersisa: **apakah total share reinsurer wajib tepat 100 %, atau boleh kurang?** Penjumlahan sendiri terbukti **tepat** dari `[data DBA]` satu baris contoh; alokasi sisa pembulatan **tidak diperlukan** selama total share 100 % | Finance | **terbuka (menyempit)** — RALAT di **§6e**. ⚠️ Penjaga yang ada **sebelah**: `Activity/CountSpreading_Act.xml` step **6.2.3** hanya menangkap **lebih dari** 100 |
| ~~P16~~ | ~~`getJSON(false)` — apakah `px*`/`py*` ikut terserialisasi?~~ | — | ✅ **TERTUTUP 2026-09-17** `[data DBA]` — `pxObjClass`, `pxListSubscript`, `pxCreateOperator`, `pyExpanded` **ada**; `pzInsKey`, `pxCreateDateTime`, `pyID` **tidak**. Ia menulis **properti yang ada di halaman**, bukan daftar metadata tetap. Konsekuensi migrasi di **§3f** |
| ~~P17~~ | ~~`10009` Marine Cargo memakai template Aneka — disengaja?~~ | — | ✅ **TERTUTUP 2026-09-17** — tercakup keputusan **"ikuti apa adanya"** pada P14 (§2). Cacat salin-template ikut dimigrasikan; bila kelak diperbaiki, catat sebagai **penyimpangan sadar** |
| ~~P14~~ | ~~`TreatyType 10004` = Fac Out (dari keterangan step)~~ | — | ✅ **TERTUTUP + DIRALAT 2026-09-17** — `[data DBA]` **`10004` = `QS (R/I)`, BUKAN Fac Out**. `InputData.CARI27` **terkonfirmasi = `TreatyGroupID`**. Ralat lengkap di **§2**; aturan pencegahnya di **§0 Aturan arti kode** |
| ~~P18~~ | ~~`pyActionConditions` — kosong di seluruh 36 Section, atau terisi dua baris di `AdjustmentDetail`?~~ | — | ✅ **TERTUTUP 2026-09-17** — kedua penyisiran saya salah; sensus benar: **3 dari 36 Section, 17 baris**. **R7 dibatalkan dan disusun ulang** (§13). Lihat RALAT di **§17c** |
| **P19** | Rule `IsFire` **tidak ada** di `When/` (dicari di **6 dari 6 berkas** folder itu), tetapi dirujuk di **6 berkas / 12 kemunculan** (dicari di **270 dari 270 berkas**). ⚠️ Angka "dua tempat" **diralat**. Pertanyaan yang tersisa: **bagaimana Pega memperlakukan `when` yang hilang — benar atau salah?** | pemilik export Pega | terbuka — **§17b**. ⚠️ **Prioritas rendah**: kedua belas pemakaiannya **sudah mati lebih dulu**, jadi dampaknya nihil di runtime |
| ~~P20~~ | ~~`SendAcceptationToKasir` membaca `ReponseCode` (salah eja)~~ | — | ✅ **KELUAR dari daftar 2026-09-17** — **cacat kosmetik**. Hanya 2 step `Property-Set` tanpa `when`; **nol step bercabang** pada nilainya. Nama properti respons **tidak terbaca dari korpus** → `[terbuka]`, bukan pemblokir. Lihat RALAT di **§17f** |
| **P23** | **DIPERSEMPIT 2026-09-17.** Bagian Section **TERTUTUP** — ke-17 elemen ber-`pyVisible` kosong ternyata label dropdown, bukan syarat (§14a ralat, **R8**). **Sisa: 3 kasus kode arah sisi-kosong di Activity** (§5b). Cara menutup: buka satu step di Pega Designer. | Pega Designer | terbuka — sisa 3 kasus |
| **P21** | **BARU.** Empat `RequestType` hidup tetapi rule-nya tidak ikut diekspor: `GenerateNoDlaTreatyIn`, `GenerateNoPLATreatyIn`, `GetDataNopolisTreatyin`, `GetTreatyInMaster_SQL`. | pemilik export Pega | terbuka — **§8d** |
| ~~P22~~ | ~~`Update_T_Storage_SQL` memakai dua format tanggal berlawanan dalam satu statement. Mana yang benar?~~ | — | ✅ **TERTUTUP 2026-09-17** `[keputusan work owner]`: **ikuti apa adanya** — kedua format dipertahankan. ⚠️ Pola dua-format juga muncul di dump JSON (§3f butir 3), jadi ini **bukan kejadian tunggal** |

| **P24** | **BARU 2026-09-17.** Tiga daftar wewenang hidup berdampingan **tanpa saling merujuk**: tingkat di **jejak audit** (4 tingkat, berhenti di `Technical Director`), pencarian **batas nilai Direktur Utama**, dan **roster komite** `EMAILKOMITE`. Apakah ketiganya **satu tangga yang sama**, atau memang **tiga hal berbeda**? | work owner + Finance | terbuka — **§11e-3**. ⚠️ Bukan bagian **Q5** (§21): Q5 soal kunci pencarian, P24 soal bentuk tangganya |
| **P15** | ✅ **TERTUTUP 2026-09-17.** `[keputusan work owner]` Total share reinsurer **wajib tepat 100 %**. Maka: pasang **penjaga dua sisi** saat input (Pega hanya menjaga sebelah — `Activity/CountSpreading_Act.xml` step 6.2.3, `@greaterThan(Local.TotalPersen,100)`, satu-satunya pemeriksaan total share). Karena data disimpan penuh tanpa pembulatan, **alokasi sisa pembulatan tidak diperlukan**. | Finance | tertutup |
| **P19** | ✅ **TERTUTUP 2026-09-17.** Kelima rule `When` sudah ditambahkan ke korpus: `IsFire` · `isGolfInsurance` · `IsAneka` · `IsMarineCargo` · `IsSpreadingUW`. Folder `When/` kini **60 berkas**. | pemilik export Pega | tertutup |
| **P21** | ✅ **TERTUTUP 2026-09-17.** Keempat rule sudah ditambahkan: `GenerateNoDlaTreatyIn` · `GenerateNoPLATreatyIn` · `GetDataNopolisTreatyin` · `GetTreatyInMaster_SQL`. Folder `RDBList/` kini **58 berkas**. | pemilik export Pega | tertutup |
| **P23** | ✅ **TERTUTUP 2026-09-17.** `[keputusan work owner]` Kode arah yang **kosong** berarti **continue when** — setara nilai `2`, step dijalankan. Berlaku untuk **67 baris syarat di 19 rule**. ⚠️ Menggantikan status **DIPERSEMPIT** yang ditulis lebih dulu pada hari yang sama. | Pega Designer | tertutup |
| **P24** | 🟡 **SEBAGIAN 2026-09-17.** Butir 1 (tingkat wewenang di jejak audit ditentukan empat nama orang yang di-hardcode — §11e-4 butir 9c) → `[keputusan work owner]` **ikuti apa adanya**, nama tetap di kode. Butir 2 dan 3 — apakah ketiga daftar wewenang satu tangga yang sama — tetap `[terbuka]`. | work owner + Finance | sebagian |
| ~~P24~~ | ✅ **TERTUTUP 2026-09-17.** `[keputusan work owner]` **(a)** Direktur Utama **di LUAR** tangga persetujuan klaim — bukan tingkat kelima. **(b)** Roster komite dan jejak audit adalah **tangga yang SAMA**: `POOLDATA.EMAILKOMITE.DEGREE` berpadanan dengan `.IsCedingConfirm`. **(c)** `Claim Admin` adalah **nilai DASAR**, bukan anggota komite — tidak ada di `DEGREE` dan memang tidak seharusnya ada. ⚠️ **Butir 1 BERUBAH**: nama hardcode **tidak ikut ke Go**; tingkat dibaca dari **satu sumber**. Aturan lengkap di **§28**. | work owner + Finance | **tertutup** |

### 19a. Sisa pemblokir sebelum `/to-spec` — **lima**

| # | Pemilik | Pertanyaan |
| --- | --- | --- |
| ~~P15~~ | Finance | ✅ **TERTUTUP** — wajib tepat 100 %; pasang penjaga dua sisi |
| ~~P19~~ | pemilik export Pega | ✅ **TERTUTUP** — lima rule `When` sudah masuk korpus; `When/` kini **60 berkas** |
| ~~P21~~ | pemilik export Pega | ✅ **TERTUTUP** — empat rule sudah masuk korpus; `RDBList/` kini **58 berkas** |
| ~~P23~~ | Pega Designer | ✅ **TERTUTUP** — kode arah kosong = **continue when** (setara `2`) |
| ~~P24~~ | work owner + Finance | ✅ **TERTUTUP** — tangga tunggal dari `EMAILKOMITE`; Direktur Utama di luar tangga; nama hardcode dibuang (§28) |

> ### ⚠️ RALAT 2026-09-17 — ringkasan ini diperbarui setelah U1
>
> Judul lama berbunyi *"sisa pemblokir — **lima**"*. Setelah lima jawaban U1 masuk, **sisa sebenarnya
> satu**: **P24 butir 2–3**. Empat lainnya tertutup. Status **DIPERSEMPIT** pada P23 (ditulis pada
> T6 hari yang sama) **digantikan** oleh **TERTUTUP**.

> ### ⚠️ RALAT 2026-09-17 (kedua) — setelah V1–V4
>
> **Yang berubah:**
>
> | Butir | Status |
> | --- | --- |
> | `[terbuka]` §25e — dua mekanisme klasifikasi lini bisnis | ✅ **DICABUT** — korpus menjawabnya sendiri: rule `When` menguji model **Fac In**, tidak ada di objek kerja Claim Prop. Yang berlaku hanya **`TreatyGroupID`** |
> | `[terbuka]` niat gerbang tak-menyaring | **tetap** — dan lokasinya **diralat ke step 16** (§24b). **Tidak memblokir** |
> | P7 (§2) | ✅ **lengkap** — triplet CLM/PLA/DLA seragam, hanya TEMP yang berbeda (§26b) |
>
> **Frontier setelah V1–V4:**
>
> | Frontier | Status |
> | --- | --- |
> | **Keputusan** (`Q1`–`Q7`) | ✅ **KOSONG** |
> | **Fakta** (`P1`–`P24`) | 🟡 **tersisa P24 butir 2–3 saja** |
>
> Dua `[terbuka]` ringan yang **tidak** memblokir: niat step 16 (§24b) dan salah ketik
> `IsCustomBonds` versus `IsCustomBond` (§27d).

> ### ✅ RALAT 2026-09-17 (ketiga) — **KEDUA FRONTIER KOSONG**
>
> Setelah **P24 tertutup** (§28):
>
> | | Status |
> | --- | --- |
> | **Frontier keputusan** (`Q1`–`Q7`) | ✅ **KOSONG** — terjawab seluruhnya |
> | **Frontier fakta** (`P1`–`P24`) | ✅ **KOSONG** — **sebelas pemblokir tertutup** |
>
> **Tabel §19 akhir: `P1`–`P24` seluruhnya tertutup.**
>
> **Yang tersisa dan TIDAK memblokir:**
>
> | Jenis | Butir |
> | --- | --- |
> | `[terbuka]` | niat gerbang `Activity/CountValueADJTreaty_Act.xml` step **16** (§24b) |
> | `[terbuka]` | salah ketik `IsCustomBonds` versus `IsCustomBond` (§27d) |
> | `[data DBA]` ringan | kunci pencocokan `EMAILKOMITE` — `EMAIL` atau kolom `USER_ID` baru (§28f) |
>
> ⚠️ Yang `[data DBA]` **wajib dikonfirmasi sebelum tiket ditulis**, tetapi **tidak menahan
> `/to-spec`** — bentuk aturannya sudah pasti, hanya nama kolom kuncinya yang menunggu.

✅ **Tertutup/gugur pada ralat 2026-09-17:** P2 (gugur) · P14 · P16 · P17 · P18 · P20 (keluar) · P22.

> ### ⚠️ RALAT 2026-09-17 — §19a semula menyebut "empat"
>
> Angka itu benar **saat ditulis**. **P24** lahir kemudian, saat meralat §23d dan menaikkan
> `IsCedingConfirm` ke §11e — jadi sisanya **lima**, bukan empat. Baris lama dicoret di judul.

1. Audit gerbang §2 (jalur Komite), §18 (efek keluar), §19 (dokumen cetak) dengan tiga aturan §0
2. Audit §14a (spreading) dan §16 (adjuster/consultant) dengan tiga aturan §0
3. **Section (36)** — menetapkan aturan baca *visible-when* lebih dulu
4. **RDBList (54)** — inventaris lengkap
5. **ReportDefinition (19) + DecisionTable (1) + SystemSettings (1)**
6. **FlowAction (12) + Harness (11) + When (6) + ConnectREST (6)**
7. **DataTransform (11)**

---

## 20. Temuan sampingan saat ralat

Hanya yang muncul sendiri saat meralat — bukan penelusuran baru.

### 20a. ⚠️ **Lima rule `When` dirujuk tetapi tidak ada**, bukan satu

`[terverifikasi]` Saat memverifikasi P19 saya menghitung isi folder `Claim Prop/When/`: **6 dari 6
berkas** adalah `IsBackStage`, `IsCLM`, `IsCLMNP`, `IsCLMP`, `IsPEGAPROD`, `IsPEGASyariah`.
Nama `When` lain yang dipakai sebagai kondisi di **270 dari 270 berkas** modul:

| Nama `When` dirujuk | Berkas perujuk | `When/<nama>.xml` ada? |
| --- | ---: | --- |
| `IsFire` | 6 | **TIDAK ADA** |
| `isGolfInsurance` | 2 | **TIDAK ADA** |
| `IsAneka` | 2 | **TIDAK ADA** |
| `IsMarineCargo` | 2 | **TIDAK ADA** |
| `IsSpreadingUW` | 2 | **TIDAK ADA** |

Keempat nama selain `IsFire` muncul di ekspresi yang sama atau bersebelahan
(`isGolfInsurance || IsAneka || IsFire` di `Section/InputAcceptation.xml` dan
`Section/OutstandingClaim.xml`; `IsMarineCargo` menggerbangi `.TreatyInMaster.Ceding`;
`IsSpreadingUW` menggerbangi satu tombol di `Section/MstAdjusterConsultant.xml`).

⚠️ **Seluruhnya berada di gerbang yang sudah mati** (`pyVisible = ALWAYS`) menurut sensus §14a —
jadi dampaknya di runtime **nihil**, sama seperti `IsFire`. Tetapi ini memperluas **P19** dari satu
rule menjadi **lima**, dan memperkuat pola: **klasifikasi lini bisnis di UI seluruhnya bersandar pada
rule `When` yang tidak ikut diekspor.**

### 20b. Nomor step yang salah di §10e

`[terverifikasi]` Rujukan `IsFire` di `Activity/PrintFileAcceptance.xml` berada di step
**`RH_1.pySteps(12)`**, bukan step 11 seperti tertulis di §10e. Status gerbangnya tidak berubah —
tetap `pyStepsPreCondition = 0` (mati).

---

## 21. Pertanyaan keputusan Q1–Q7

Bab ini menutup cacat yang dijelaskan di **§0 Aturan penomoran**: `Q1`–`Q5` disebut tujuh kali di
berkas ini tanpa pernah didefinisikan.

**Semua Q1–Q5 DITULIS, tidak ada yang dicabut** — kelimanya dapat direkonstruksi utuh dari
pertanyaan yang memang diajukan saat frontier Ronde 2 dihitung. Dua di antaranya ternyata memuat
**dua keputusan dalam satu pertanyaan**, dan dipecah sesuai aturan: **Q4 → Q4 + Q6**, **Q5 → Q5 + Q7**.

**Pembeda `P` dan `Q`:** `P<n>` (§19) adalah **fakta yang belum diketahui** — dijawab dengan
menemukan bukti. `Q<n>` (bab ini) adalah **keputusan yang harus diambil** — tidak ada bukti yang bisa
menjawabnya, karena korpus hanya merekam apa yang *terjadi*, bukan apa yang *seharusnya*.

---

### Q1 — Gerbang yang tertulis tetapi tidak berjalan: ikuti kenyataan atau ikuti niat?

**Pertanyaan:** Untuk aturan yang tertulis di sistem lama tetapi ternyata **tidak pernah dijalankan**,
apakah sistem baru harus berperilaku **sama seperti sistem lama yang berjalan sekarang**, atau
**menegakkan aturan sebagaimana tertulis**?

**Pilihan jawaban:**
1. **Ikuti kenyataan** — sistem baru berperilaku persis seperti produksi hari ini; aturan yang tidak
   berjalan diperlakukan sebagai catatan niat, bukan spesifikasi.
2. **Ikuti niat** — pasang semua aturan yang tertulis; sistem baru menjadi lebih ketat daripada yang
   lama.
3. **Triase satu per satu** — putuskan per kasus, dengan daftar lengkap disiapkan lebih dulu.

**Bukti:** **84 dari 535** step ber-syarat di Activity tidak berlaku (§0 Aturan 1), dan **57 dari
236** gerbang UI di Section mati (§14a). Contoh konkret:
`Activity/CountSpreadingADJ_Act.xml` step **3.1** — syarat "abaikan baris yang ditolak" tidak
berlaku, sehingga baris yang sudah ditolak ikut dihitung;
`Activity/SaveOutstanding_Act.xml` step **39** — lihat Q3.

**Kenapa ini keputusan manusia:** korpus menunjukkan **dua kebenaran yang bertentangan** — apa yang
ditulis dan apa yang dijalankan — dan tidak memuat apa pun yang menyatakan mana yang dimaksudkan.
Data produksi selama bertahun-tahun lahir di bawah perilaku yang berjalan, bukan di bawah niat.

**Pemilik:** work owner (bersama Product + UW untuk kasus yang menyentuh angka)

**Yang terbuka bila dijawab:** bentuk aturan bisnis di seluruh spec. ~~⚠️ **Prasyarat: P23 (§19)
harus tertutup lebih dulu** — selama arah 17 dari 180 elemen belum diketahui, jumlah gerbang mati
hanya diketahui sebagai rentang **57–74 dari 236**, dan daftar triase (pilihan 3) belum dapat
disusun.~~

> ⚠️ **RALAT 2026-09-17 — prasyarat DICABUT.** Ke-17 elemen itu ternyata **tidak membawa syarat sama
> sekali** (label dropdown Designer — §14a ralat, **R8** §13). Angka gerbang UI mati **pasti 57 dari
> 236**; rentang 57–74 **dibatalkan**. Q1 dapat dijawab tanpa menunggu P23, dan **sudah dijawab**.


**Jawaban work owner 2026-09-17:** Pilihan 1 — ikuti kenyataan. `[keputusan work owner]`
Gerbang yang teks syaratnya tersimpan tetapi `pyStepsPreCondition` tidak `true` **tidak ditulis
di Go**. Yang tidak ditulis adalah **WHEN-nya**; **STEP-nya tetap ditulis** dan berjalan tanpa
saringan, persis seperti produksi hari ini. Alasan work owner: teks itu sisa yang tidak terhapus
bersih dari Pega, bukan aturan yang dimaksudkan berlaku. Cakupan: 84 di Activity, 44 di Section.
---

### Q2 — Langkah yang di-remark: tetapkan sebagai aturan berdiri "tidak dimigrasikan"?

**Pertanyaan:** Apakah langkah kerja yang sudah dinonaktifkan di sistem lama boleh **seluruhnya
diabaikan** dalam sistem baru, sebagai aturan yang berlaku untuk seluruh modul?

**Pilihan jawaban:**
1. **Ya, aturan berdiri** — langkah ter-remark tidak dimigrasikan, konsisten dengan **ADR-0006** yang
   sudah menetapkan hal sama untuk Claim Life.
2. **Ya, tetapi dengan kewajiban lapor** — tidak dimigrasikan, dan setiap remark yang memutus sebuah
   jalur kerja **wajib dilaporkan** sebelum diabaikan.
3. **Tidak** — periksa satu per satu sebelum memutuskan.

**Bukti:** **68 dari 270** berkas modul memuat step ter-remark (§0 Aturan 2). Yang memutus jalur
kerja: `Activity/TryMakePLA_Act.xml` step **8** (jalur penomoran PLA lama mati);
`Activity/PrintDLATreatyIn.xml` step **15.11.2** dan **15.11.3** (perhitungan salvage dan adjuster fee
tidak pernah jalan); `Activity/SaveAdjusterConsultant_Act.xml` step **5** (jalur simpan master lewat
SQL mati — ternyata digantikan `Obj-Save` step **4**, §6f).

**Kenapa ini keputusan manusia:** korpus menunjukkan langkah itu mati, tetapi **tidak menyatakan
apakah matinya disengaja atau kelalaian**. Tiga kasus di atas mematikan fungsi yang masih dipakai
bisnis.

**Pemilik:** work owner

**Yang terbuka bila dijawab:** §6f, §10e, dan daftar rule "hilang" di §8d dapat ditutup — tiga dari
tujuh `RequestType` menggantung ternyata dipanggil dari step ter-remark, jadi tidak perlu diminta ke
pemilik export.


**Jawaban work owner 2026-09-17:** Pilihan 1 — aturan berdiri. `[keputusan work owner]`
Langkah ber-`pyStepsBlockName` = `//` **tidak dimigrasikan sama sekali** — STEP-nya tidak ditulis,
bukan hanya syaratnya. Berlaku se-modul, 68 step di Claim Prop. Alasan: menjaga kode tetap bersih.
---

### Q3 — Pengiriman ke Arasapas mendahului keputusan komite: pertahankan atau perbaiki?

**Pertanyaan:** Sistem lama mengirim data akseptasi ke sistem luar **sebelum** komite selesai
memutuskan, padahal aturan yang tertulis mensyaratkan menunggu — apakah sistem baru harus menunggu?

**Pilihan jawaban:**
1. **Perbaiki** — tunggu komite selesai sebelum mengirim.
2. **Pertahankan** — sistem penerima memang mengharapkan data lebih awal; menunggu justru merusak.
3. **Kirim dua tahap** — kirim penanda awal, lalu kirim hasil final setelah komite memutuskan.

**Bukti:** `Activity/SaveOutstanding_Act.xml` step **39** — syarat tertulisnya menuntut keputusan
komite sudah final, tetapi syarat itu **tidak berlaku**, sehingga pengiriman terjadi tanpa menunggu
(§1).

**Kenapa ini keputusan manusia:** pihak penerima berada **di luar korpus**; apakah mereka bergantung
pada kiriman awal itu tidak dapat diketahui dari rule Pega mana pun.

**Pemilik:** work owner + pemilik sistem Arasapas

**Yang terbuka bila dijawab:** urutan efek keluar di §18g dapat difinalkan. Ini juga **kasus uji
paling tajam untuk Q1**: bila Q1 dijawab "ikuti kenyataan" tetapi Q3 dijawab "perbaiki", maka Q1
berlaku sebagai default dengan pengecualian untuk efek yang melintasi batas sistem.


**Jawaban work owner 2026-09-17:** Ikut Q1. `[keputusan work owner]`
Gerbang `AcceptStatus==1 && KomiteCount==KomiteLoop` pada `Activity/SaveOutstanding_Act.xml`
step 39 tidak ditulis. Arasapas dikirim tanpa menunggu komite, sama seperti produksi hari ini.
---

### Q4 — Unggahan berkas berjalan di lingkungan non-produksi: digerbangi atau tidak?

**Pertanyaan:** Apakah unggahan berkas ke penyimpanan awan harus **dibatasi hanya pada lingkungan
produksi**?

**Pilihan jawaban:**
1. **Ya, digerbangi** — hanya produksi yang boleh mengunggah ke penyimpanan sungguhan.
2. **Tidak** — semua lingkungan mengunggah ke penyimpanan yang sama seperti sekarang.

**Bukti:** `Activity/InsertGoogleStorage_Act.xml` step **11** dan
`Activity/GetUrlGoogleStorage_Act.xml` step **6.5** — keduanya memanggil layanan penyimpanan
**tanpa gerbang lingkungan sama sekali**; ditambah `Activity/InsertGoogleStorage_Act.xml` step **6**
yang gerbang lingkungannya **ada tetapi tidak berlaku** (§10c, §17e). Bandingkan
`Activity/HitServiceToKasir_Act.xml` step **9.7** dan **13.4**, serta
`Activity/KonversiKlaim_Act.xml` step **3** — ketiganya **digerbangi dan gerbangnya hidup**, dengan
keterangan step berbunyi "kalau diserver dev jangan dijalanin" (§17e).

**Kenapa ini keputusan manusia:** korpus memperlihatkan **dua pola berlawanan pada masalah yang
sama**, dan tidak menyatakan mana yang dimaksudkan. Ini juga menyangkut apakah berkas uji boleh
mendarat di penyimpanan produksi — pertanyaan tata kelola, bukan teknis.

**Pemilik:** work owner + pemilik infrastruktur

**Yang terbuka bila dijawab:** §18f dan §10c dapat difinalkan; **Q6** baru dapat dijawab.


**Jawaban work owner 2026-09-17:** Ikut Q1. `[keputusan work owner]`
Gerbang `IsPEGAPROD` pada `Activity/InsertGoogleStorage_Act.xml` step 6 tidak ditulis. Semua
lingkungan mengunggah ke penyimpanan yang sama.
---

### Q5 — Batas wewenang dicari lewat nama orang: diganti dengan apa?

**Pertanyaan:** Batas nilai wewenang saat ini dicari dengan **mencocokkan nama satu orang**; apa yang
harus dipakai sebagai gantinya?

**Pilihan jawaban:**
1. **Jabatan** — cari berdasarkan posisi, bukan orang.
2. **Peran atau tingkat komite** — cari berdasarkan tingkat wewenang.
3. **Tetap nama orang** — dengan kewajiban memutakhirkan setiap kali orangnya berganti.

**Bukti:** `RDBList/GetLimitDirekturUtama_SQL.xml` (rule Connect-SQL — tanpa step Pega) mencari baris
dengan mencocokkan **nama orang sebagai literal**; dipanggil dari
`Activity/AttachmentProtect_ACT.xml` step **5** (§13a, §13c). Kolom **`JABATAN`** dan **`DEGREE`**
sudah tersedia pada sumber data yang sama — `ReportDefinition/FilterEmailKomiteWithLimit.xml`
mengambil keduanya (§16c), jadi pilihan 1 dan 2 **tidak menuntut perubahan skema**.

**Kenapa ini keputusan manusia:** korpus tidak memuat isi tabel roster, sehingga tidak dapat
dipastikan jabatan mana yang setara dengan nama yang tertanam. Itu fakta organisasi, bukan fakta kode.

**Pemilik:** Finance + work owner

**Yang terbuka bila dijawab:** §13a dan §13c dapat difinalkan; **Q7** baru dapat dijawab.


**Jawaban work owner 2026-09-17:** Pilihan 3 — ikuti apa adanya. `[keputusan work owner]`
`RDBList/GetLimitDirekturUtama_SQL.xml` tetap mencari `where name = <nama orang>` (nilai tidak
disalin — PII; di korpus ia literal).
---

### Q6 — Pecahan dari Q4: apa yang dilakukan di lingkungan non-produksi?

**Pertanyaan:** Bila unggahan berkas digerbangi (Q4 dijawab "ya"), apa yang terjadi saat pengguna
mengunggah berkas di lingkungan pengujian?

**Pilihan jawaban:**
1. **Penyimpanan tiruan** — berkas disimpan di tempat lain, alur dokumen tetap dapat diuji
   ujung-ke-ujung.
2. **Tolak** — unggahan gagal dengan pesan jelas.
3. **Lewati diam-diam** — seperti perilaku lingkungan dev pada integrasi Kasir sekarang.

**Bukti:** `Activity/InsertDocument_Act.xml` step **5** — baris dokumen **tidak tersimpan** bila
unggahan gagal (§10c). Artinya pilihan 2 dan 3 membuat **seluruh alur dokumen tidak dapat diuji** di
luar produksi.

**Kenapa ini keputusan manusia:** ini keputusan cara kerja tim, bukan aturan bisnis; korpus tidak
memuat lingkungan pengujian sama sekali.

**Pemilik:** work owner + pemilik infrastruktur

**Yang terbuka bila dijawab:** rencana pengujian alur dokumen yang menyentuh §18f.


**Jawaban work owner 2026-09-17: — GUGUR.** Q6 hanya berlaku bila Q4 dijawab "digerbangi".
Q4 dijawab "tidak", jadi tidak ada perilaku khusus di lingkungan non-produksi.
---

### Q7 — Pecahan dari Q5: apa yang terjadi bila batas wewenang tidak ditemukan?

**Pertanyaan:** Bila pencarian batas nilai wewenang **tidak menemukan baris apa pun**, apakah operasi
harus **gagal terang-terangan** atau **lanjut seolah batas tidak terlampaui**?

**Pilihan jawaban:**
1. **Gagal terang-terangan** — operasi ditolak dengan pesan; tidak ada yang lolos tanpa pemeriksaan.
2. **Lanjut** — seperti perilaku sekarang.

**Bukti:** `RDBList/GetLimitDirekturUtama_SQL.xml` mengembalikan kosong bila nama yang ditanam tidak
lagi cocok; pemanggilnya `Activity/AttachmentProtect_ACT.xml` step **5** tidak memiliki penanganan
untuk hasil kosong, sehingga penanda "melewati batas" tidak pernah menyala dan **kewajiban dokumen
tambahan hilang tanpa pesan apa pun** (§13a).

**Kenapa ini keputusan manusia:** ini pilihan arah kegagalan — longgar atau ketat — dan menyangkut
risiko pengendalian, bukan fakta yang dapat ditemukan di korpus.

**Pemilik:** Finance + work owner

**Yang terbuka bila dijawab:** §13a dapat difinalkan bersama Q5.


**Jawaban work owner 2026-09-17:** Baris batas wewenang **dijamin selalu ada** secara bisnis.
`[keputusan work owner]` Perlakuan sama dengan roster komite di tiket 10 Claim Life: dijamin ada,
**tetapi tetap dipasang penjaga defensif**, ditandai ⚠️ penyimpangan sadar — bukan alur normal.
---

### 21c. ⚠️ Tiga risiko yang diterima sadar

Konsekuensi langsung dari jawaban Q1, Q3, Q5. Diterima work owner, bukan kelalaian migrasi.

1. **Uang.** `Activity/CountSpreadingADJ_Act.xml` step 6 — saringan payment type `1|2|5` tidak
   ditulis, sehingga perhitungan spreading di Go berjalan untuk **semua** payment type.
2. **Data di sistem luar.** Klaim yang akhirnya **ditolak** komite tetap memiliki data di Arasapas.
   `[terverifikasi]` **Nol rule pembatal di 270 dari 270 berkas.**
3. **Nama orang masuk ke kode.** Berbeda dari Q1–Q4: gerbang mati tidak meninggalkan jejak apa pun,
   sedangkan `GetLimitDirekturUtama_SQL` adalah query **hidup** — nama orang itu benar-benar ada di
   dalam kode Go.

---




### 21f. Pemeriksaan identitas — dijalankan ulang setelah W1–W3 (2026-09-17) — **frontier kosong**

Mekanis atas berkas, bukan dari ingatan.

```
Identitas dirujuk : P1 … P24                                  (24)  ← SELURUHNYA TERTUTUP
                    Q1 … Q7                                    (7)  ← SELURUHNYA TERJAWAB
                    R1 … R9                                    (9)
                    OQ-018 OQ-039                              (2, di badan berkas)

Punya definisi    : P1 P3 … P14                    — di bab 2   (pemblokir TERTUTUP)
                    P2 P14 P15 … P24               — di bab 19  (P14 di kedua bab)
                    Q1 … Q7                        — di bab 21
                    R1 … R9                        — di bab 13
                    OQ-018 OQ-039                  — discovery/open-questions.md (§0)
TIDAK ada definisi: nihil
```

**Bab §28 tidak melahirkan identitas `P`/`Q`/`R`/`OQ` baru** — ia menuliskan **jawaban** P24, bukan
pertanyaan baru. Satu butir `[data DBA]` di §28f sengaja **tidak** diberi nomor `P`: ia bukan fakta
korpus yang hilang, melainkan konfirmasi skema yang menunggu DBA, dan **tidak memblokir `/to-spec`**.

**Dikecualikan, sama seperti pemeriksaan sebelumnya:** `Q8`–`Q13` dan tujuh `OQ` lain muncul hanya di
dalam §21b — inventaris identitas berkas `grilling-ronde-1.md`.
### 21e. Pemeriksaan identitas — dijalankan ulang setelah V1–V5 (2026-09-17)

Mekanis atas berkas, bukan dari ingatan.

```
Identitas dirujuk : P1 … P24                                  (24)
                    Q1 … Q7                                    (7)
                    R1 … R9                                    (9)
                    OQ-018 OQ-039                              (2, di badan berkas)

Punya definisi    : P1 P3 … P14                    — di bab 2   (pemblokir TERTUTUP)
                    P2 P14 P15 … P24               — di bab 19  (P14 di kedua bab)
                    Q1 … Q7                        — di bab 21
                    R1 … R9                        — di bab 13
                    OQ-018 OQ-039                  — discovery/open-questions.md (§0)
TIDAK ada definisi: nihil
```

**Tidak ada identitas baru dari V1–V5.** Tiga aturan yang ditambahkan ke §0 pada V2 — **Aturan
silsilah**, **Aturan sensus nama rule**, **Aturan konflik cek versus aturan berdiri** — sengaja
**tidak diberi nomor `R`**, karena R1–R9 adalah aturan **membaca gerbang**, sedangkan ketiganya
aturan **kerja**. Keduanya tetap terdefinisi utuh di §0 sesuai Aturan penomoran.

**Dikecualikan, sama seperti pemeriksaan sebelumnya:** `Q8`–`Q13` dan tujuh `OQ` lain muncul hanya di
dalam §21b — inventaris identitas berkas `grilling-ronde-1.md`.
### 21d. Pemeriksaan identitas — dijalankan ulang setelah T1–T9 dan U1–U6 (2026-09-17)

Mekanis atas berkas, bukan dari ingatan.

```
Identitas dirujuk : P1 … P24                                  (24)
                    Q1 … Q7                                    (7)
                    R1 … R9                                    (9)   ← R8, R9 baru
                    OQ-018 OQ-039                              (2, di badan berkas)

Punya definisi    : P1 P3 P4 … P14                  — di bab 2   (pemblokir TERTUTUP)
                    P2 P14 P15 … P24                — di bab 19  (P14 di kedua bab)
                    Q1 … Q7                         — di bab 21
                    R1 … R7                         — di bab 13  (tabel aturan baca)
                    R8 R9                           — di bab 13  (### R8, ### R9)
                    OQ-018 OQ-039                   — discovery/open-questions.md
                                                      (rujukan eksplisit, §0)
TIDAK ada definisi: nihil
```

**Dikecualikan, sama seperti pemeriksaan sebelumnya:** `Q8`–`Q13` dan tujuh `OQ` lain muncul **hanya
di dalam §21b** — inventaris identitas berkas `grilling-ronde-1.md`, bukan rujukan dari badan berkas
ini.

**Identitas baru sejak pemeriksaan terakhir:** `R8` (§13, label dropdown bukan ekspresi) dan `R9`
(§13, `pyConditionViewer` sisa basi). Keduanya terdefinisi di bab yang sama dengan R1–R7.
Bab **§22**, **§24**, **§25**, **§26**, **§27** tidak melahirkan identitas `P`/`Q`/`R`/`OQ` baru.
### 21a. Pemeriksaan identitas — `grilling-ronde-2.md`

Dijalankan secara mekanis atas seluruh isi berkas, bukan dari ingatan.

```
Identitas dirujuk : P1 P2 P3 P4 P5 P6 P7 P8 P9 P10 P11 P12 P13 P14
                    P15 P16 P17 P18 P19 P20 P21 P22 P23        (23)
                    Q1 Q2 Q3 Q4 Q5 Q6 Q7                        (7)
                    R1 R2 R3 R4 R5 R6 R7                        (7)
                    OQ-018 OQ-039                               (2, di badan berkas)

Punya definisi    : P1 P3 P4 P5 P6 P7 P8 P9 P10 P11 P12 P13 P14 — di bab 2
                                                                  (pemblokir TERTUTUP)
                    P2 P14 P15 P16 P17 P18 P19 P20 P21 P22 P23  — di bab 19
                                                                  (status diperbarui;
                                                                   P14 muncul di kedua bab)
                    Q1 Q2 Q3 Q4 Q5 Q6 Q7                        — di bab 21
                    R1 R2 R3 R4 R5 R6 R7                        — di bab 13
                                                                  (R7 dibatalkan + R7-baru,
                                                                   keduanya di bab 13)
                    OQ-018 OQ-039                               — di luar berkas ini:
                                                                  discovery/open-questions.md
                                                                  (rujukan eksplisit, sesuai §0)
TIDAK ada definisi: nihil
```

> ### ⚠️ RALAT 2026-09-17 — pemeriksaan dijalankan ulang setelah ralat `[data DBA]`
>
> Setelah §2, §3f, §6e, §10e, §19, §19a dan §23 ditulis, pemeriksaan **dijalankan ulang secara
> mekanis**. Hasilnya **tetap `nihil`**, dengan dua perubahan:
>
> - **`P14` kini muncul di §19** (baris status `✅ TERTUTUP + DIRALAT`) selain definisinya di §2 —
>   keduanya sah, dan tidak menimbulkan identitas tanpa definisi.
> - **`P2`, `P16`, `P17`, `P22`** berpindah status di §19 dari *terbuka* menjadi *gugur/tertutup*;
>   definisinya **tetap ada**, jadi tidak ada rujukan yang menggantung.
>
> Tidak ada identitas baru yang lahir dari bab **§23** — bab itu merujuk nomor bab (§3e, §8g, §11,
> §14b), bukan identitas `P`/`Q`/`R`/`OQ`.

> ### ⚠️ RALAT 2026-09-17 (kedua) — dijalankan ulang setelah §11e dan P24
>
> Setelah `IsCedingConfirm` dinaikkan ke **§11e** dan **P24** dibuat, pemeriksaan **dijalankan ulang
> secara mekanis**. Hasilnya **tetap `nihil`**, dengan satu tambahan:
>
> ```
> Identitas dirujuk : P1 … P24   (24)   ← bertambah P24
>                     Q1 … Q7     (7)
>                     R1 … R7     (7)
>                     OQ-018 OQ-039 (2, di badan berkas)
>
> Punya definisi    : P1 P3 … P14                  — di bab 2
>                     P2 P14 P15 … P24             — di bab 19  (P24 baru; P14 di kedua bab)
>                     Q1 … Q7                      — di bab 21
>                     R1 … R7                      — di bab 13
>                     OQ-018 OQ-039                — discovery/open-questions.md
> TIDAK ada definisi: nihil
> ```
>
> **P24** didefinisikan penuh di **§19** (baris tabel) dan diuraikan di **§11e-3**. Bab **§11e** tidak
> melahirkan identitas baru selain itu.

**Dua kemunculan yang sengaja dikecualikan, dan alasannya:**

- **`Q8`–`Q13`** muncul di §0 dan §21b, tetapi **bukan identitas yang dipakai berkas ini** — keduanya
  adalah **laporan tentang cacat di `grilling-ronde-1.md`**, dan di sana justru dinyatakan **tidak
  punya definisi**.
- **`OQ-002`, `OQ-020`, `OQ-021`, `OQ-029`, `OQ-040`, `OQ-041`, `OQ-060`** muncul **hanya di dalam
  daftar §21b**, yaitu inventaris identitas berkas ronde 1 — bukan rujukan dari badan berkas ini.
  Seluruhnya didefinisikan di `discovery/open-questions.md`.

> ### ⚠️ RALAT 2026-09-17 — versi pertama pemeriksaan ini sendiri tidak akurat
>
> Tulisan pertama §21a menyebut hanya **`P2 P15…P23`** dan **`OQ-018 OQ-039`**, dan melewatkan
> **`P1`, `P3`–`P14`** yang didefinisikan di **§2**, serta tujuh `OQ` lain yang muncul di §21b.
> Kesalahan itu ketahuan karena pemeriksaannya **dijalankan mekanis atas berkas**, bukan disusun dari
> ingatan — dan itulah gunanya langkah ini. Angka di atas adalah hasil yang sudah dikoreksi.

### 21b. Pemeriksaan identitas — `grilling-ronde-1.md` (dilaporkan saja, berkas tidak diubah)

```
Identitas dirujuk : P1 … P14   (tabel §21 berkas itu)
                    Q1 … Q13   (kalimat penutup berkas itu)
                    OQ-002 OQ-018 OQ-020 OQ-021 OQ-029 OQ-039 OQ-040 OQ-041 OQ-060
Punya definisi    : P1 … P14                 — di bab 21 berkas itu
                    OQ-002 … OQ-060          — di luar berkas: discovery/open-questions.md
TIDAK ada definisi: Q1 Q2 Q3 Q4 Q5 Q6 Q7 Q8 Q9 Q10 Q11 Q12 Q13   ← 13 identitas
```

⚠️ **Cacat di `grilling-ronde-1.md`:** kalimat penutupnya menyebut **`Q1–Q13`** dan mengarahkan
pembaca ke **§21 berkas itu**, padahal §21 di sana berisi **`P1`–`P14`** — **nol definisi `Q`**.
Ketiga belas pertanyaan itu hanya pernah diucapkan, tidak pernah ditulis.

**Tidak diperbaiki** sesuai batasan tugas. Bila Ronde 1 akan dipakai untuk `/to-spec`, cacat ini
harus ditutup lebih dulu dengan cara yang sama seperti §21 di sini.

---

> ⚠️ **RALAT 2026-09-17:** catatan lama berbunyi *"Nomor §22 sengaja tidak dipakai"*. **Tidak lagi berlaku — §22 kini terisi** (aturan sisa versus niat).

---

## 22. Aturan turunan — sisa versus niat

### 22a. Keputusan: 75 elemen bersyarat selalu-salah tidak dibuat di Go

Sensus 36 dari 36 Section, syarat bernilai selalu salah (`NEVER`, `never`, `1=2`, `1=3`):

| Keluarga | Syarat berlaku | Syarat TIDAK berlaku | Jumlah |
| --- | ---: | ---: | ---: |
| A — `pyCondition` + `pyVisible` | 43 | 12 | 55 |
| B — `pyContainerVisibleWhen` | 12 | 8 | 20 |
| **Jumlah** | **55** | **20** | **75** |

Sebaran keluarga A: `OutstandingClaim` 10 · `InputAcceptation` 7 · `InputAcceptation_Est` 7 ·
`OutstandingClaim_Est` 7 · `InputAcceptation_Adjs` 3 · sisanya tersebar.

`[keputusan work owner]` **Ketujuh puluh lima elemen tidak dibuat di Go** — termasuk 20 yang hari
ini masih tampil karena sakelarnya tidak terpasang.

### 22b. Aturan: sisa versus niat

> Teks syarat yang tertinggal karena Pega tidak menghapusnya bersih **tidak membawa niat** — ia
> diabaikan, dan step-nya tetap ditulis tanpa saringan (Q1). Sedangkan `NEVER` atau `1=2` yang
> **diketik developer** adalah niat yang terbaca terang: elemen itu memang dimaksudkan tidak tampil,
> sehingga **elemennya tidak dibuat**. Keduanya berujung tidak-ditulis, tetapi alasannya berbeda dan
> dicatat berbeda.

### 22c. Peta nasib 236 ekspresi UI

| Nasib di Go | A | B | Jumlah |
| --- | ---: | ---: | ---: |
| Tidak dibuat — syarat selalu salah | 55 | 20 | **75** |
| Syarat ditulis jadi kondisi tampil | 76 | 31 | **107** |
| Syarat tidak ditulis, elemen selalu tampil | 32 | 5 | **37** |
| Diabaikan — label dropdown, bukan syarat | 17 | 0 | **17** |
| **Jumlah** | **180** | **56** | **236** |

⚠️ Dari **57** gerbang UI yang tidak berlaku, **20** menempel pada elemen yang toh tidak dibuat.
Yang bermakna "elemen tampil padahal ada syaratnya" hanya **37**.

---

## 23. Temuan dari dump JSON produksi

**Sumber seluruh bab ini:** `[data DBA]` — **satu baris contoh** dari
`os_akseptasi_klaim.DATA_JSON` dan `json_klaim.DATA_JSON`, klaim
`RNM-K22.09.2026.T00718` / `CLMP-7134`. **Bukan korpus, bukan `[terverifikasi]`.**
Satu klaim saja — **jangan digeneralisasi**. Setiap kali bab ini menyebut korpus, penyebutnya
disebutkan.

---

### 23a. Satu nilai, enam penulisan berbeda

`[data DBA]`

```
83333333.3333                          ClaimAmount
83333333.33330000000000                TotalListClaimAmount    · ListClaimAmount.Value
83333333.333300000000000               TotalListClaimAmountIDR · ListClaimAmount.USD
83333333.333300000000000000000000      TotalGrossEstimateIDR   · TotalGrossEstimateTreaty
83333333.3333000000000000000000000     SpreadingRisk.ClaimEstimation
```

**Nilainya identik; yang berbeda hanya panjang padding nol.** Ini **artefak tipe decimal Pega**,
**bukan presisi yang bermakna**.

**Konsekuensi migrasi:** di Go seluruhnya menjadi **satu nilai decimal tunggal**. Panjang padding
**tidak boleh** dibawa, dan **tidak boleh** dipakai menyimpulkan presisi kolom.

---

### 23b. ⚠️ Pemotongan 4 desimal terjadi di HULU, lalu merambat

`[data DBA]`

```
GrossValue  83333333.3333 000000…      ← hanya 4 desimal bermakna, sisanya nol
Value       19999999.999992 0000…      = 83333333.3333 × 24 %   (cocok sampai digit terakhir)
```

**`Value` dihitung dari `GrossValue` yang SUDAH terpotong**, bukan dari angka asalnya.

Ini **contoh nyata** dari pola "hasil dipakai sebagai masukan" yang sudah tercatat di **empat
tempat** di §3e (`CountSpreadingADJ` 7.2 · `CountSpreading` 10.1 dan 11.2.1 · `CopyOldataCurr` 10.1).
Sebelumnya pola itu hanya terbaca dari struktur rule; kini **terlihat akibatnya pada angka**.

⚠️ **Konsekuensi mengikat** `[keputusan work owner]`: **Go harus memutus rantai ini** — nilai turunan
dihitung dari **sumber asli**, bukan dari nilai yang sudah terpotong.

**Tanpa itu, menyimpan presisi penuh tidak ada gunanya** — angkanya sudah rusak sebelum disimpan.
Ini melengkapi keputusan §3e bahwa `.ClaimSpreaded` adalah turunan yang boleh dihitung ulang: bukan
hanya *boleh*, melainkan **harus**, dan **dari hulu**.

---

### 23c. Dua nama kolom yang berbohong — tambahan untuk §8g

`[data DBA]`

| Kolom | Nilai pada satu baris contoh | Kebohongannya |
| --- | --- | --- |
| `ListClaimAmount[].USD` | `83333333.3333` | **bukan nilai USD** — `Currency` pada baris yang sama berbunyi `"IDR"` |
| `ListClaimAmount[].IDR` | `"1.0"` | **ini KURS, bukan nilai** |
| `SuggestList[].IsCedingConfirm` | peran + **nama orang** (nilai tidak disalin — PII) | **bukan boolean**, meski namanya berawalan `Is` |

⚠️ Pasangan `USD`/`IDR` itu **saling bertukar peran**: kolom bernama mata uang berisi **nilai**, dan
kolom bernama mata uang lain berisi **kurs**. Menambah daftar §8g yang sudah memuat lebih dari 20
nama rule dan kolom yang berbohong.

⚠️ `IsCedingConfirm` melengkapi pola yang sama pada `.IsOldData` (bernilai `"Yes"`, bukan boolean)
dan `STS_REJECT` (nilai `1` berarti *diaksep*) — **awalan `Is`/`STS` di korpus ini bukan penanda
tipe boolean.**

> ### ⚠️ RALAT 2026-09-17 — `IsCedingConfirm` DINAIKKAN ke §11e; ia bukan sekadar nama yang berbohong
>
> Mencatatnya di sini **terlalu rendah**. Isinya bukan "peran + nama orang" yang kebetulan aneh —
> isinya **tangga persetujuan empat tingkat**:
> `Claim Admin` → `Claim Dept. Head` → `Operational Director` → `Technical Director`,
> dengan **tiga tingkat teratas identik persis** di kedua penulisnya.
>
> Karena itu ia **bagian dari mesin jejak audit**, dan dipindahkan ke **§11e** —
> lengkap dengan ketidakkonsistenan tingkat terbawah, tiga konsekuensi migrasinya, sensus 4 dari 270
> berkas, dan hasil pemeriksaan dua Section penampilnya.
>
> **Yang tetap berlaku di sini:** namanya memang berbohong — `Is…` **bukan** boolean. Dua baris lain
> di tabel §23c (`ListClaimAmount[].USD` dan `.IDR`) **tidak berubah**.

---

### 23d. `SuggestList` adalah jejak audit — dan urutannya **tidak kronologis**

`[data DBA]` **12 entri** pada klaim contoh. Bidangnya: `CommentSuggest` · `DateSuggest` ·
`PICSuggest` · `IsCedingConfirm`.

⚠️ **Entri terakhir bertanggal `20260915`, sementara sebelas lainnya `20260903`** — **daftar tidak
terurut waktu**.

**Konsekuensi migrasi:** di Go, urutan jejak audit **harus diambil dari kolom tanggal**, **bukan**
dari urutan baris. Menampilkan apa adanya akan menyajikan riwayat dengan urutan yang salah.

#### Hubungan dengan §11 — **satu jejak, tetapi DUA penulis**

Pertanyaannya: apakah `SuggestList` di dump adalah tujuan yang sama dengan `InsertChronology_DT`?

**Ya — keempat bidangnya cocok persis** dengan yang ditulis `DataTransform/InsertChronology_DT.xml`
menurut §11a (`.CommentSuggest`, `.PICSuggest`, `.DateSuggest`, `.IsCedingConfirm`). **Satu jejak,
bukan dua.**

⚠️ **Tetapi ada penulis KEDUA yang belum tercatat di §11.** Sensus `SuggestList` di **270 dari 270
berkas** `Claim Prop`:

| Berkas | Peran |
| --- | --- |
| `DataTransform/InsertChronology_DT.xml` | penulis — sudah tercatat §11a |
| **`Activity/SethistoryKlaimTreaty.xml`** (`ASM-FW-GCNMFW-WORK` / `SETHISTORYKLAIMTREATY`) | ⚠️ **penulis kedua** — menyentuh `.IsCedingConfirm` **7 kali**, di dalam blok berulang **`RH_1.pySteps(1)`** beserta sub-step **1–11** |
| `Section/InputAcceptation.xml` · `Section/OutstandingClaim.xml` | **pembaca** — grid `pyPageListProperty = pyWorkPage.ClaimData.SuggestList` (§14b) |

`[terverifikasi]` Sensus penulis `.IsCedingConfirm`: `Activity/SethistoryKlaimTreaty.xml` **7** ·
`DataTransform/InsertChronology_DT.xml` **5** · kedua Section **2 masing-masing**.

⚠️ **Bukti bahwa data produksi datang dari penulis kedua, bukan dari `InsertChronology_DT`:**
§11a mencatat `InsertChronology_DT` hanya pernah menulis empat nilai — `"Claim Admin"` (default) dan
tiga jabatan hasil pemetaan nama. ~~Nilai pada satu baris contoh **bukan salah satu dari
keempatnya**; ia berbentuk **peran + nama orang**.~~ `[terverifikasi]` Korpus juga memuat pembacaan
`@contains(.IsCedingConfirm,"Admin") || .IsCedingConfirm==""` — pola yang hanya masuk akal bila
nilainya **teks bebas berimbuhan**, bukan enum tertutup.

> ### ⚠️ RALAT 2026-09-17 — "bukan salah satu dari keempatnya" SALAH
>
> **Kesimpulannya benar, alasannya keliru.** Nilai pada satu baris contoh **memang cocok** dengan
> nilai yang ditulis rule — hanya **bukan dengan rule yang saya bandingkan**.
>
> `[terverifikasi]` `Activity/SethistoryKlaimTreaty.xml` menulis
> **`"Admin Claim " + .PICSuggest`** untuk tingkat terbawah, lalu `"Claim Dept. Head"`,
> `"Operational Director"`, `"Technical Director"` — tiga terakhir **identik** dengan
> `InsertChronology_DT`.
>
> `[data DBA]` Nilai pada satu baris contoh berbentuk **`"Admin Claim " + nama orang`** — itu
> **persis pola penulis kedua**, bukan nilai di luar enum. Jadi entri audit pada klaim contoh
> **ditulis `SethistoryKlaimTreaty`**, dan itulah yang membuktikan penulis kedua aktif di produksi.
>
> Tangga lengkap, ketidakkonsistenannya, dan konsekuensi migrasinya kini di **§11e**.

⚠️ **Akibatnya §11a perlu dibaca ulang:** kesimpulan bahwa jejak audit sepenuhnya dihasilkan
`InsertChronology_DT` — termasuk lubang "aksi IT Developer tidak terekam" — **hanya berlaku untuk
entri yang ditulis rule itu**. Entri dari `Activity/SethistoryKlaimTreaty.xml` **tidak melewati
gerbang `OperatorID.pyPosition != "IT Developer"`**, sehingga lubang audit itu **lebih sempit dari
yang saya tulis** — tetapi **seberapa sempit tidak terbaca dari korpus** → `[terbuka]`.

⚠️ Ini juga menjelaskan **bug `CloseClaimProp`** (§11b): rule itu mengisi `DataChronology.CARI12`
sementara `InsertChronology_DT` membaca `CARI1`. Bila jalur audit yang sebenarnya dipakai saat close
adalah `SethistoryKlaimTreaty`, maka dampaknya berbeda dari yang saya tulis — **tidak dapat
dipastikan dari korpus** → `[terbuka]`.

---

## 24. Akibat P23 — sebaran kode arah menjadi pasti

`[keputusan work owner]` Kode arah **kosong = continue when** (setara `2`).

Sensus **513 baris syarat pada gerbang AKTIF**, Claim Prop:

| BENAR | SALAH | Jumlah | Arti setelah P23 |
| --- | --- | ---: | --- |
| 2 | 3 | 380 | arah normal |
| (kosong) | 3 | **47** | **arah normal** — kosong = lanjut |
| 3 | 2 | 44 | arah TERBALIK |
| 3 | (kosong) | **17** | **arah TERBALIK** |
| 6 | 2 | 9 | benar → keluar activity |
| 5 | 2 | 7 | benar → keluar iterasi `[terbuka]` |
| 2 | 6 | 6 | salah → keluar activity |
| 6 | (kosong) | 2 | benar → keluar activity · salah → lanjut |
| (kosong) | 2 | 1 | **dua sisi lanjut — gerbang tidak menyaring apa pun** |

⚠️ **Gerbang berarah terbalik = 44 + 17 = 61**, bukan 44. Angka 44 di §14a dan mana pun ia dirujuk
harus diralat menjadi **61** untuk konteks Activity.

> ### ⚠️ RALAT 2026-09-17 — §5b `[terbuka]` ditutup oleh P23
>
> §5b mencatat **tiga kasus kode arah sisi-kosong** di Activity (`Activity/ProteksiData_act.xml`
> step 3 dan 6, `Activity/CountSpreading_Act.xml` step 6.2.3) dan menyatakan perilakunya
> **"bergantung default engine Pega, tidak terbaca dari korpus"** → `[terbuka]`.
>
> **P23 menutupnya:** kosong = **continue when**. Ketiga kasus itu kini terbaca pasti.
> Khusus `Activity/ProteksiData_act.xml` step 3 (BENAR=3, SALAH=kosong→lanjut): pembacaan lama
> — pesan *"ID MASTER Sedang dalam tahap Revisi"* muncul justru saat status **BUKAN**
> "Resolve Complete" — **TETAP BERLAKU**, dan sekarang berstatus `[terverifikasi]`, bukan dugaan.

### 24a. Tujuh belas gerbang terbalik yang baru ketahuan

Polanya seragam — `X != ""` dengan BENAR=3, artinya step jalan justru saat nomornya **KOSONG**.
Ini langkah **pembuatan nomor**: buat nomor bila belum ada.

| Berkas | Step | Syarat |
| --- | --- | --- |
| `Activity/SaveOutstanding_Act.xml` | 10 · 13 · 14 · 20 | `NoClaim` / `ClaimNo != ""` |
| `Activity/TryMakePLA_Act.xml` | 7 · 8 · 14 | `NoPla != ""` |
| `Activity/PrintDLATreatyIn.xml` | 6 · 7 | `DLA_No != ""` |
| `Activity/UpdateTableOS.xml` | 2 · 3 · 5 | `NoClaim != ""` |
| `Activity/DeleteEstimation_Act.xml` | 14 · 15 | `Local.CheckEst` |
| `Activity/CheckNopolicy_Act.xml` | 4 | `PolicyNo == ""` |
| `Activity/ProteksiData_act.xml` | 3 | `@contains(…,"Resolve Complete")` |
| `Activity/CountValueADJTreaty_Act.xml` | 13.1 | `.AcceptanceStatus == 2` |

⚠️ **Dibaca tanpa aturan P23, kesimpulannya TERBALIK TOTAL** — akan terbaca sebagai *"nomor dibuat
ulang padahal sudah ada"*.

#### Kesimpulan Ronde 1 dan Ronde 2 yang terdampak

| Rule / step | Status kesimpulan lama |
| --- | --- |
| `Activity/TryMakePLA_Act.xml` step **7 · 8** | **tidak terdampak** — kedua step **DI-REMARK** (§10e), jadi arahnya tak berarti |
| `Activity/PrintDLATreatyIn.xml` step **6 · 7** | **tidak terdampak** — keduanya **DI-REMARK** (§10e) |
| `Activity/TryMakePLA_Act.xml` step **14** | ✅ **TETAP BERLAKU** — §10e sudah membacanya sebagai "tulis hanya bila `NoPla` masih kosong"; itu pembacaan terbalik yang benar |
| `Activity/ProteksiData_act.xml` step **3** | ✅ **TETAP BERLAKU**, dan `[terbuka]`-nya ditutup (lihat ralat di atas) |
| `Activity/SaveOutstanding_Act.xml` step **10 · 13 · 14 · 20** | ⚠️ **belum pernah dianalisis** di ronde mana pun — tidak ada kesimpulan lama yang perlu diralat, tetapi **belum ada pula yang menutupinya** |
| `Activity/UpdateTableOS.xml` step **2 · 3 · 5** | ⚠️ idem — belum pernah dianalisis |
| `Activity/DeleteEstimation_Act.xml` step **14 · 15** | ⚠️ idem |
| `Activity/CheckNopolicy_Act.xml` step **4** | ⚠️ idem |
| `Activity/CountValueADJTreaty_Act.xml` step **13.1** | ⚠️ idem |

**Nol kesimpulan lama yang harus dibalik.** Yang ada adalah **lima rule yang gerbangnya kini terbaca
tetapi perilakunya belum pernah ditelusuri** — dicatat apa adanya, **tidak** ditebak.

### 24b. ⚠️ Satu gerbang yang tidak menyaring apa pun

> ⚠️ **RALAT nomor step 2026-09-17.** Perintah menyebut *step 16*; menurut
> `<pyStepPageReference>` gerbang ini milik **step `RH_1.pySteps(17)`**. Nomor step diambil dari
> korpus, bukan dari perintah.

```
Activity/CountValueADJTreaty_Act.xml   step 17
  pyStepsPreCondition          true
  pyStepsPreCondParamsWhen     .AdjustmentValue > .TotalEstimasiValue
  BENAR = (kosong)   SALAH = 2
```

`[terverifikasi]` **Apa yang dikerjakan step 17:** `Property-Set-Messages`, keterangan step
*"set error message adjustment lebih besar estimasi hanya untuk claim"* — yaitu memasang pesan
kesalahan **"nilai adjustment melebihi estimasi"**.

**Akibat mekanisnya:** BENAR kosong = lanjut (P23), SALAH = 2 = lanjut. **Kedua sisi lanjut.**
Gerbangnya ada, tetapi **step berjalan tanpa peduli hasil perbandingannya** — pesan kesalahan
dipasang **baik saat adjustment melebihi estimasi maupun saat tidak**.

⚠️ Perbandingan ini menyentuh **uang** — nilai adjustment terhadap estimasi.

⚠️ `[terbuka]` **Niatnya tidak disimpulkan.** Apakah pesan itu memang dimaksudkan selalu terpasang
(mis. ditampilkan bersyarat di UI), atau gerbangnya rusak, **tidak terbaca dari korpus**. Yang pasti
hanya pembacaan mekanisnya di atas. Pertanyaan ini **tidak** diberi nomor P baru karena bukan fakta
yang hilang, melainkan niat yang hanya work owner tahu.

> ### ⚠️ RALAT 2026-09-17 — **step 16, bukan 17. Saya tertukar, dan akibatnya lebih tajam.**
>
> **Saya salah dua kali:** menempatkan gerbang tak-menyaring di step 17, dan mengutip keterangan step
> 17 untuk menjelaskannya. **Keterangan itu milik step 17; pasangan kodenya milik step 16.**
>
> `[terverifikasi]` Urutan sebenarnya di `Activity/CountValueADJTreaty_Act.xml` — kedua step sama-sama
> ber-`pyStepsPreCondition = true`, dan **masing-masing punya DUA baris syarat**:
>
> ```
> step 16   Property-Set — ket "set iserror == 2 adjustment lebih besar estimasi"
>    baris 1   BENAR=(kosong)  SALAH=2    .AdjustmentValue > .TotalEstimasiValue
>    baris 2   BENAR=(kosong)  SALAH=3    Local.SumProposeAdjustment > .TotalEstimasiValue
>
> step 17   Property-Set-Messages — ket "set error message adjustment lebih besar estimasi hanya untuk claim"
>    baris 1   BENAR=(kosong)  SALAH=3    .ValueAdjustment > .TotalEstimasiValue
>    baris 2   BENAR=(kosong)  SALAH=3    Local.SumProposeAdjustment > .TotalEstimasiValue && .Type=="1"
> ```
>
> **Hanya SATU baris yang tidak menyaring: baris 1 step 16** (`kosong`/`2` — kedua sisi lanjut).
> Semua baris lain ber-`SALAH=3` dan **menyaring normal**. Tabel §24 tetap benar: **satu** baris
> berkategori "dua sisi lanjut", dan ia ada di **step 16**.
>
> ### ⚠️⚠️ Akibatnya — lebih tajam dari yang saya tulis
>
> Step 16 punya **dua** baris syarat: **baris pertama tidak menyaring, baris kedua menyaring**.
> Jadi step 16 **efektif digerbangi hanya oleh TOTAL**, sementara **pemeriksaan per baris tertulis
> tetapi tidak berpengaruh**.
>
> `[terverifikasi]` **Satu baris adjustment yang melebihi estimasi TIDAK ditandai error, selama total
> seluruh baris masih di bawah estimasi.** Pemeriksaan per baris `.AdjustmentValue >
> .TotalEstimasiValue` **ada di korpus tetapi tidak berpengaruh**. **Menyentuh uang.** Niat aslinya
> `[terbuka]`.
>
> **Catatan metode:** kesalahan saya berasal dari mengaitkan syarat ke `<pyStepPageReference>` yang
> muncul **sesudahnya**. Yang benar: deskripsi, `pyStepsPreCondition`, dan baris-baris syarat
> mengikuti **sesudah** referensi step pemiliknya.

---


---

## 25. Sisir 60 rule `When` — 54 di antaranya belum pernah dibaca

**Cakupan: 60 dari 60 berkas** `Claim Prop/When/` (dulu 6). Modul kini **328 berkas**.

### 25a. ⚠️ Aturan baca baru — **R9** (lihat juga §13)

`[terverifikasi]` Tiap berkas `When` punya **dua** tempat kondisi, dan keduanya **sering tidak
sinkron**:

| Lokasi | Isi | Status |
| --- | --- | --- |
| `<pyCondition>` PageList — tiap `rowdata` ber-`<pyConditionLabel>` (A, B, C…) + `<pyConditionValue1>` | ekspresi yang **benar-benar dieksekusi** | **otoritatif** |
| `<pyConditionViewer>` → `pyNestedConditions` → `<pyConditionString>` + `<pyLogic>` di dalamnya | teks tampilan editor | **sisa basi** |

⚠️ `<pyLogic>` yang berlaku adalah **anak langsung `pagedata`**, **bukan** yang di dalam
`pyConditionViewer`. Labelnya selalu cocok 1:1 dengan label di `pyCondition`.

**Bukti sisa-basi, bukan tafsiran:**

- `When/IsMarineCargo.xml` — viewer memuat **6 baris** `Kode Bisnis = "06"/"07"/"10"/"17"/"18"/"24"`
  dengan `pyLogic` viewer `A OR … OR F`, sedangkan `pyCondition` **hanya punya label G dan H**, dan
  `pyLogic` asli = **`G OR H`**. Label A–F sudah dihapus; sampahnya tertinggal di viewer.
- Enam baris `Kode Bisnis = 06/07/10/17/18/24` yang **identik** muncul di `IsHE`, `IsLiability`,
  `IsMarineCargo`, `IsOilGas`, `IsProductsLiability`, `IsProfessionalLiability` — **jejak
  salin-tempel, bukan logika**.
- `When/IsFire.xml` — viewer ber-`pyLogic` `A OR … OR U` (**21 label**) padahal isinya hanya satu
  baris `[Double click to add condition]`.

⚠️ **Tanpa R9, kesimpulan tentang kode lini bisnis akan salah total** — orang akan membaca
`Kode Bisnis 06/07/10/17/18/24` sebagai aturan aktif, padahal ia tidak pernah dieksekusi.

`[terverifikasi]` Menguatkan §3 (silsilah): `<pzOriginalInstanceKey>` **tidak dapat dipakai sebagai
identitas**. Jejak *save-as* di `When/`: `IsPA`←`ISTRAVEL` · `IsWorkmenCompensation`←`ISCIT` ·
`IsCLMP`←`ISCLM` · `IsEnvironmental`←`ISCRIME` · `IsYieldShortfall`←`ISGROWINGTREES` ·
`IsEdmAdjShareCedant`←`ISEDMADJREFNO` · `IsContractorsPlantMachinery`←`ISWORKMENCOMPENSATION`.
Seluruh **60 dari 60** bertipe `RULE-OBJ-WHEN`.

### 25b. Inventaris — ringkasan 60 dari 60

Mayoritas berpola sama: satu uji `Quotation.BusinessType = "<literal>"`. Yang menyimpang:

| Rule | Yang diuji |
| --- | --- |
| `When/IsCLM.xml` · `IsCLMP` · `IsCLMNP` | `pyWorkCover.pyWorkIDPrefix` **OR** `pyWorkPage.pyWorkIDPrefix` |
| `When/IsPEGAPROD.xml` | `pxProcess.pzProductionLevel = "5"` |
| `When/IsPEGASyariah.xml` | `pxProcess.pxSystemNodeID = "jboss1074"` |
| `When/IsBackStage.xml` | `.pyNote = "Back"` — class `ASM-FW-GCNMFW-WORK-PNC` |
| `When/IsUW.xml` | `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName` = 4 nama workbasket (`A OR B OR C OR D`) |
| `When/IsSpreadingUW.xml` | `pyWorkPage.IsSpreadingUW = "1"` |
| `When/IsEDM.xml` · `IsEdmAdjRefNo` · `IsEdmAdjShareCedant` | `StatusBusiness = 3`, `.EdmType`, `.Type` |
| `When/whenRiskType.xml` | `.IndividualRiskType = 1` **OR** `= 2` |
| `When/isPA_PNC.xml` | `pyWorkPage.Quotation.GroupPanel = "002"` |
| `When/IsNotTravelPA.xml` | `!= "Travel"` **AND** `!= "PA"` |

⚠️ **Dua anomali:**

- `When/IsBondingKBG.xml` — **dua kondisi persis sama** (`A OR B` redundan).
- `When/isElectronicEquipment.xml` — membandingkan dengan literal **ber-spasi**
  `" ElectronicEquipment "`, sehingga **tidak akan pernah cocok** dengan nilai bersih
  `"ElectronicEquipment"`.

⚠️ `When/IsUW.xml` adalah **satu-satunya** rule di modul ini yang menguji **workbasket operator** —
mendekati pemeriksaan peran. Tetapi ia **yatim** (§25d), jadi tidak mengubah kesimpulan §14d bahwa
nol gerbang berbasis peran benar-benar aktif.

### 25c. ⚠️ Kenapa `IsAneka` 118 KB — ia **simpul agregator**, bukan daftar kode

`[terverifikasi]` `When/IsAneka.xml` punya **26 kondisi**, dan **25 di antaranya `evaluateWhen(...)`
ke rule `When` lain**. Hanya **satu** yang menguji properti: label **U** =
`pyWorkPage.Quotation.BusinessType = "Aneka"`.

⚠️ **Label D hilang** dari `pyLogic` (`A OR B OR C OR E OR …`) — satu kondisi pernah dihapus, label
lain tidak dinomori ulang.

**Sebab ukurannya:** tiap `rowdata` membawa satu blok `<pyFunctionData>` lengkap plus `<pyParameters>`
berisi 5 baris `Embed-MethodParams` walau hanya 2–3 terisi ≈ **~4 KB boilerplate per kondisi**.
Garisnya konsisten di seluruh folder:

| Rule | Kondisi | Ukuran |
| --- | ---: | ---: |
| `isGolfInsurance` | 1 | 18 KB |
| `IsMarineCargo` | 2 | 30 KB |
| `IsFire` | 5 | 36 KB |
| `IsLiability` | 7 | 51 KB |
| **`IsAneka`** | **26** | **118 KB** |

**Jadi ia besar karena menggabung 25 rule lain dengan `OR`, bukan karena mendaftar kode satu per
satu.**

`When/IsFire.xml` (`A OR B OR C OR D OR E`): rantai ke `IsKPR`, `IsOilGas`, `IsFireStyle1`,
`IsFireStyle2`, ditambah satu uji `BusinessType = "Fire"`.
⚠️ Inkonsistensi jalur: kondisi E membaca `pyWorkPage.OfferFacIn.QuotationData.BusinessType`,
sedangkan keempat sub-rule-nya membaca `pyWorkPage.Quotation.BusinessType` — **dua jalur properti
berbeda dalam satu rule**.

### 25d. Yatim — **14 dari 60**

`[terverifikasi]` Dicari di **328 dari 328 berkas**, dengan batas kata, setelah membuang field
non-fungsional (`pzOriginalInstanceKey`, `pyMemo`, `pyNote`, `pyDescription`). Rujukan yang dihitung
hanya pemanggilan nyata: `evaluateWhen("…")`, `<pyCondition>`, `<pyStepsPreCondParamsWhen>`,
`<pyContainerVisibleWhen>`, `<pyTaskWhen>`.

**46 dari 60 dirujuk · 14 dari 60 yatim:**
`isAnalistorTransfer` · `IsBonding` · `IsClaim` · `IsCustomBond` · `IsEDM` · `IsEdmAdjRefNo` ·
`IsEdmAdjShareCedant` · `IsNotTravelPA` · `IsPA` · `isPA_PNC` · `IsSPK` · `IsTravel` · `IsUW` ·
`whenRiskType`.

Catatan penting per kasus:

- `When/IsBonding.xml` — `pyMemo` di `IsAneka` berbunyi *"ganti IsBonding jadi
  IsBondingAndCustomBonds"*. **Sengaja digantikan**, berkasnya tertinggal.
- `When/IsClaim.xml` — duplikat fungsional kondisi B `IsCLM`; yang dipakai `IsCLM`.
- ⚠️ `When/IsCustomBond.xml` — **rujukan menggantung**. `IsBondingAndCustomBonds` kondisi C memanggil
  `evaluateWhen("IsCustomBonds")` — **dengan huruf S**. Tidak ada rule bernama `IsCustomBonds` di 328
  berkas; yang ada `ISCUSTOMBOND` (tanpa S). **Rantai ini putus di dalam korpus** → `[terbuka]`.
- `When/IsUW.xml` — tampak dirujuk `Activity/GetDetailPolis_act.xml`, tetapi itu **false positive**:
  entri di sana ber-`pxRuleObjClass = Rule-Obj-Property`, bukan rule `When`.

⚠️ **Ketergantungan terpusat:** hanya **12 dari 60** dipanggil konsumen non-`When`. **34 dari 60**
hidup hanya lewat rantai `When`→`When`, dan **25 dari 34 hanya dirujuk `When/IsAneka.xml`**.
**Bila `IsAneka` dinonaktifkan, 25 rule `When` langsung menjadi kode mati.**

### 25e. ⚠️⚠️ Hubungan ke P14 — **DUA mekanisme klasifikasi lini bisnis, nol irisan**

**Jawaban langsung: TIDAK, keempat rule `When` itu tidak memakai kode yang sama dengan `CARI27`.**

`[terverifikasi]` **Nol dari 60 berkas** `When/` menyebut `TreatyGroupID` maupun `CARI27`.

| | **Mekanisme A** | **Mekanisme B** |
| --- | --- | --- |
| Kunci | `InputData.CARI27` = `ClaimData.TreatyGroupID` | `Quotation.BusinessType` (+ `Quotation.BusinessCode`) |
| Bentuk nilai | **numerik**, 19 kode | **literal teks**, ~37 literal (+8 kode) |
| Dipakai untuk | memilih template `InsuredInterest` | mengontrol visibilitas UI + precondition activity |
| Bukti | `Activity/SetValueToClaim_Act.xml` step **9, 10, 11, 12** | 60 rule `When` |

`[terverifikasi]` **Irisan antara 8 kode `BusinessCode` dan 19 kode `CARI27` = KOSONG.** Rentangnya
pun berbeda: `CARI27` seluruhnya `10002`–`10059`; `BusinessCode` di `When` `10021`–`10184` tanpa satu
pun beririsan. **Tidak ada berkas di korpus yang menerjemahkan antara keduanya.**
**Tidak saya simpulkan mana yang benar** — itu pertanyaan work owner.

⚠️ **Titik singgungnya nyata:** `Section/InputAcceptation.xml` dan `Section/OutstandingClaim.xml`
memakai `pyCondition` = `isGolfInsurance || IsAneka || IsFire` (**mekanisme B**) untuk mengatur
visibilitas, sementara `ClaimData.InsuredInterest` di **halaman yang sama** diisi lewat
**mekanisme A**. Keduanya bertemu di satu layar.

> ### ⚠️ RALAT 2026-09-17 — dua koreksi terhadap §2 (P14)
>
> 1. **Pemetaan `CARI27` bukan hanya di step 9.** Ia tersebar di
>    `Activity/SetValueToClaim_Act.xml` step **9, 10, 11, dan 12**, masing-masing dengan template
>    berbeda. Penautan `CARI27` → `TreatyGroupID` sendiri terjadi di **step 8**:
>    `TreatyGroupID ← @if(TreatyGroupID == "", InputData.CARI27, TreatyGroupID)`.
> 2. **Keranjang Aneka memuat 16 kode, bukan 14.** Step 10 mendaftar `10012` `10011` `10055` `10059`
>    `10002` `10003` `10004` `10006` `10008` `10010` `10013` `10026` `10051` `10057` `10058` `10054`.
>    **Total kode `CARI27` yang dikenal = 19** (16 + `10007` + `10016` + `10009`).
>
> ✅ **Menguatkan P17:** step **10** (Aneka) dan step **12** (Marine Cargo) memakai **template yang
> persis sama**. Jadi kesamaan template Marine Cargo ↔ Aneka bukan tafsiran — ia **dua step terpisah
> dengan isi identik**.

> ### ⚠️ RALAT 2026-09-17 — **bukan dua mekanisme yang harus dipilih. `[terbuka]` DICABUT.**
>
> Saya menyerahkan ini ke work owner sebagai "dua mekanisme berdampingan, tidak saya simpulkan mana
> yang benar". **Tidak perlu — korpus menjawabnya sendiri.**
>
> `[terverifikasi]` Kelima rule lini bisnis menguji **halaman model Fac In**, bukan model Claim Prop:
>
> ```
> When/IsFire.xml           pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Fire"
> When/IsMarineCargo.xml    .OfferFacIn.QuotationData.BusinessType           = "MarineCargo"
>                           pyWorkPage.Quotation.BusinessType                = "MarineCargo"
> When/isGolfInsurance.xml  pyWorkPage.Quotation.BusinessType = "GolfInsurance"
> When/IsAllRisk.xml        pyWorkPage.Quotation.BusinessType = "AllRisk"
> When/IsAneka.xml          pyWorkPage.Quotation.BusinessType = "Aneka"
> ```
>
> `[data DBA]` **Data klaim nyata — satu baris contoh — tidak punya properti itu:**
>
> ```
> QuotationData.BusinessName = "PROPERTY ALL RISK"
> QuotationData.BusinessCode = "10045"
> ```
>
> **`BusinessName` / `BusinessCode`, bukan `BusinessType`. Beda nama, beda halaman.**
>
> ### Kesimpulan yang menggantikan yang lama
>
> `[terverifikasi]` **Inilah sebab keenam belas gerbang lini bisnis dimatikan.** Rule `When` itu
> menguji properti yang **tidak ada pada objek kerja Claim Prop** — ia **sisa impor dari model
> Fac In**. **Dihidupkan pun semuanya bernilai salah.**
>
> **Bukan dua mekanisme yang harus dipilih salah satu: yang berlaku di Claim Prop hanya satu —
> `TreatyGroupID` (`InputData.CARI27`).** `[terbuka]` **dicabut — tidak perlu keputusan work owner.**
>
> ⚠️ `When/IsMarineCargo.xml` menguji **dua halaman berbeda dalam satu rule**
> (`.OfferFacIn.QuotationData` **dan** `pyWorkPage.Quotation`) — **dua generasi bertumpuk**.
>
> ### Seberapa luas sisa impor Fac In ini
>
> `[terverifikasi]` Sensus **60 dari 60 berkas** `Claim Prop/When/`:
>
> | | Jumlah |
> | --- | ---: |
> | Menyebut `OfferFacIn` atau `Quotation` — **sisa impor Fac In** | **51 dari 60** |
> | Tidak menyebut keduanya | **9 dari 60** |
>
> Kesembilan yang bukan sisa impor: `IsBackStage` · `IsCLM` · `IsCLMNP` · `IsCLMP` · `IsClaim` ·
> `IsPEGAPROD` · `IsPEGASyariah` · `IsSpreadingUW` · `whenRiskType` — persis kelompok "menyimpang"
> di §25b, ditambah `IsClaim`.
>
> Sensus properti di **60 dari 60**: `BusinessType` **44** berkas · `BusinessCode` **10** ·
> `BusinessName` **1** · **`TreatyGroupID` NOL**.
>
> ⚠️ **Jadi 85 % folder `When/` adalah sisa impor model Fac In**, dan **nol** di antaranya menyentuh
> kunci klasifikasi yang benar-benar dipakai Claim Prop. Ini menjelaskan sekaligus **mengapa 14 dari
> 60 yatim** (§25d) dan **mengapa gerbang yang memakainya mati** (§14a).
>
> ⚠️ *Catatan cakupan:* sensus 51/9 dihitung dari **kemunculan nama** di berkas. Karena **R9**,
> sebagian kemunculan bisa berada di `pyConditionViewer` yang basi; kondisi otoritatif kelima rule
> utama sudah dibaca langsung di §25c dan cocok.

---


---

## 26. Sisir 4 rule `RDBList` yang baru masuk korpus

**Cakupan: 58 dari 58 berkas** `Claim Prop/RDBList/` (dulu 54).

### 26a. Keempat rule

| | `RDBList/GenerateNoDlaTreatyIn.xml` | `RDBList/GenerateNoPLATreatyIn.xml` | `RDBList/GetDataNopolisTreatyin.xml` | `RDBList/GetTreatyInMaster_SQL.xml` |
| --- | --- | --- | --- | --- |
| Class | `ASM-FW-GCNMFW-Int-V_POLIS` | `ASM-FW-GCNMFW-Int-V_POLIS` | `ASM-FW-GISFW-Int-policyjson` | `ASM-FW-GISFW-Int-OFFERJSON` |
| Nama asal | ⚠️ **`GenerateNoDla`** (save-as; induk **tidak ada** di modul) | sama — **rule asli** | sama | ⚠️ **`GetLimitsTreatyIn_SQL`** (save-as; induk **ada**) |
| Objek | `POOLDATA.GENERATE_NODLATREATYIN` | `POOLDATA.GENERATE_NOPLATREATYIN` | tabel `policyjson` — **tanpa prefix schema** | `POOLDATA.TREATYINDETAIL` + `POOLDATA.TREATYINDETAILEDM` (`UNION ALL`) |
| Masuk | `TempDla.CARI30/32/33/34/31` | `TempCFS.CARI30/32/33/34/31` | `{InputData.CARI2}` | **nol parameter** |
| Keluar | `OutputData.START_DATE` · `.TSI` | `OutputData.START_DATE` · `.TSI` | nol `out`; 7 alias kolom `CARI1..CARI7` | nol `out`; 9 kolom per cabang |
| **`COMMIT;`** | **YA** | **YA** | **tidak** | **tidak** |
| Hardcode | — | — | — | **`'Proportional'`** ×2 |
| Pemanggil | `Activity/PrintDLATreatyIn.xml` step **7** | `Activity/TryMakePLA_Act.xml` step **8** | `Activity/CheckNoPolicy.xml` step **4** | `Activity/GetMasterTreaty_Act.xml` step **3** |

> ### ⚠️ RALAT 2026-09-17 — dua premis rujukan di perintah keliru
>
> Keduanya bukan pemanggilan runtime, melainkan **garis keturunan save-as**:
>
> 1. `RDBList/GenerateNoCLMTreatyIn.xml` **tidak memanggil** `GenerateNoPLATreatyIn`. Sebaliknya —
>    `pzOriginalInstanceKey` milik CLM **persis sama** dengan `pzInsKey` milik PLA. **CLM adalah
>    hasil save-as DARI PLA.** Itulah sebabnya CLM tetap memakai halaman `TempCFS`.
> 2. `RDBList/GetTreatyInMasterProp_SQL.xml` **tidak memanggil** `GetTreatyInMaster_SQL`. Rantai
>    keturunannya tiga tingkat:
>    `GetCopyNBForClaim` → `GetLimitsTreatyIn_SQL` → **`GetTreatyInMaster_SQL`** →
>    `GetTreatyInMasterProp_SQL`. Hanya `GetCopyNBForClaim` yang di luar modul.
>
> Ini contoh langsung dari **§3** (silsilah `pzOriginalInstanceKey`) dan **aturan kerja §13**
> (rujuk path lengkap, periksa nama kembar).

⚠️ `[terverifikasi]` **Tabrakan slot `CARI2`** di `GetDataNopolisTreatyin`: dipakai **masuk** sebagai
nomor polis, tetapi **keluar** sebagai `Nooffer`. Halamannya berbeda (`InputData` versus hasil),
tetapi rawan tertukar saat dibaca manusia.

### 26b. ⚠️ Penomoran — DLA dan PLA **sepola dengan klaim**; yang menyimpang justru *temp*

**Menjawab P7 (§2) selengkapnya.**

| | `GENERATE_NOCLMTREATYIN` | `GENERATE_NODLATREATYIN` | `GENERATE_NOPLATREATYIN` | `GENERATE_NOCLMTRTYINTEMP` |
| --- | --- | --- | --- | --- |
| Masuk | **5** — KODE, KODE_BIS, BULAN, TAHUN, TIPE | **5** — sama | **5** — sama | **3** — KODE, TAHUN, TIPE |
| Keluar | 2 — START_DATE, TSI | 2 — sama | 2 — sama | 2 — sama |
| Pemetaan CARI | 30, 32, 33, 34, 31 | 30, 32, 33, 34, 31 | 30, 32, 33, 34, 31 | 30, 34, 31 |
| Halaman input | `TempCFS` | ⚠️ **`TempDla`** | `TempCFS` | `TempCFS` |
| Halaman output | `OutputData` | `OutputData` | `OutputData` | ⚠️ **`OutptData`** (salah ketik) |
| `COMMIT;` | ya | ya | ya | ya |

**Ketiga procedure `GENERATE_NO{CLM,DLA,PLA}TREATYIN` adalah triplet identik** — tanda tangan sama
(5→2), pemetaan `CARI` sama, blok `DECLARE…BEGIN…END` sama baris demi baris. **Yang menyimpang justru
varian *temp***, yang membuang `KODE_BIS` dan `BULAN` — konsisten dengan P7 (§2): nomor temp **tidak
dapat dipromosikan** jadi final karena tidak punya bulan dan kode bisnis.

⚠️ Keidentikan CLM ↔ PLA **bukan kebetulan** — CLM di-save-as dari PLA (§26a ralat). **Perubahan pada
salah satu tidak akan merambat ke yang lain.**

⚠️ `[terverifikasi]` `RDBList/GenerateNoTRTInTemp.xml` menulis ke **`OutptData`** (kurang huruf `u`),
sedangkan tiga saudaranya ke `OutputData`. Salah ketik lama, tetapi menjadikan penomoran *temp*
satu-satunya yang menulis ke **halaman berbeda**.

> ### ✅ P7 LENGKAP 2026-09-17 — pola penomoran kini utuh
>
> `[terverifikasi]` Dari `<pyBrowseSQL>` keempat rule:
>
> ```
> POOLDATA.GENERATE_NOCLMTREATYIN     KODE, KODE_BIS, BULAN, TAHUN, TIPE  → START_DATE, TSI
> POOLDATA.GENERATE_NOPLATREATYIN     KODE, KODE_BIS, BULAN, TAHUN, TIPE  → START_DATE, TSI
> POOLDATA.GENERATE_NODLATREATYIN     KODE, KODE_BIS, BULAN, TAHUN, TIPE  → START_DATE, TSI
> POOLDATA.GENERATE_NOCLMTRTYINTEMP   KODE,           TAHUN, TIPE         → START_DATE, TSI
> ```
>
> **Nomor klaim, PLA, dan DLA memakai pola tanda tangan yang sama persis — lima parameter masuk, dua
> keluar. Hanya varian TEMP yang berbeda: tanpa `KODE_BIS` dan tanpa `BULAN`.**
>
> ✅ **Menguatkan kesimpulan P7** (§2): **nomor temp tidak dapat dipromosikan menjadi nomor final —
> ia kekurangan dua komponen pembentuk.** Sebelumnya itu disimpulkan dari dua procedure saja; kini
> terbukti dari **empat**, dengan triplet CLM/PLA/DLA sebagai pembanding yang seragam.
>
> **Dua SQL baru lainnya, ringkas:** `RDBList/GetDataNopolisTreatyin.xml` membaca tabel `policyjson`;
> `RDBList/GetTreatyInMaster_SQL.xml` membaca `POOLDATA.TREATYINDETAIL` dengan
> `PROPORTIONTYPE = 'Proportional'` lalu `UNION ALL` ke `POOLDATA.TREATYINDETAILEDM`.

### 26c. ⚠️ Inventaris stored procedure — ralat §3c: **12 procedure**, bukan 10

`[terverifikasi]` **Dua procedure baru**, keduanya dari berkas baru:

```
POOLDATA.GENERATE_NODLATREATYIN     (5 in / 2 out)   ← BARU
POOLDATA.GENERATE_NOPLATREATYIN     (5 in / 2 out)   ← BARU
```

**Inventaris terkini: 12 procedure + 2 fungsi + 1 sequence** (dulu 10 + 2 + 1).

`[terverifikasi]` Dua nama yang tampak seperti procedure sebenarnya **tabel** yang tertangkap regex
karena daftar kolom `INSERT INTO`: `POOLDATA.DIRECTTOKASIR_LOG` dan `pooldata.monitoring_klaim_log`.
`DBMS_LOB.CREATETEMPORARY` adalah built-in Oracle, bukan objek aplikasi.

### 26d. ⚠️ Nama yang berbohong — satu kasus baru

`[terverifikasi]` **`RDBList/GetTreatyInMaster_SQL.xml`** menjanjikan "**Master**" treaty-in, tetapi
**kedua cabang `UNION ALL` menyaring diam-diam** `WHERE a.PROPORTIONTYPE = 'Proportional'`.
**Seluruh treaty non-proporsional tidak akan pernah muncul**, tanpa apa pun di nama rule yang
mengisyaratkannya.

⚠️ Perbandingan yang menajamkan: tetangganya `GetTreatyInMasterProp_SQL` memakai filter **yang sama
persis** tetapi **jujur menyebut "Prop"** di namanya. **Nama yang lebih spesifik justru yang akurat;
nama yang lebih umum yang menyesatkan.**

⚠️ Tambahan: `pyClassName`-nya `ASM-FW-GISFW-Int-**OFFERJSON**` padahal **tidak menyentuh tabel offer
sama sekali** — warisan dari save-as `GetLimitsTreatyIn_SQL`.

`[terverifikasi]` **Bukan** kasus baru: nol dari empat berkas baru yang bernama `Get*` melakukan DML
atau `COMMIT`. Pelanggar kategori "Get yang menulis" tetap yang lama — `GetSequenceNumber_SQL` dan
`GetTokenStorage_SQL`.

---

## 27. Sensus ulang — penyebut berubah dari 270 menjadi 328

⚠️ **Korpus `Claim Prop` bertambah 58 berkas:** `When/` 6 → **60** (+54), `RDBList/` 54 → **58** (+4).
Total **270 → 328**. Setiap sensus lama yang berbunyi *"270 dari 270"* **tidak lagi sah** untuk kedua
folder itu.

`[terverifikasi]` Berkas ini memuat **11 klaim ber-penyebut "270 dari 270"**. Berikut hasil
penjalanan ulang atas **328 dari 328 berkas**.

### 27a. Sensus yang dijalankan ulang — **hasil tidak berubah, hanya penyebutnya**

| Sensus | Bab | Angka lama (270) | **Angka baru (328)** |
| --- | --- | --- | --- |
| `DataChronology.CARI1` | §11e-4 / 9d | 34 kemunculan | **34** — tidak berubah |
| `DataChronology.CARI12` | §11e-4 / 9d | 1 kemunculan | **1** — tidak berubah |
| `ResponseCode` (ejaan benar) | §17f | 0 kemunculan | **0** — tidak berubah |
| `ReponseCode` (salah eja) | §17f | — | **3** kemunculan |
| Penulis/pembaca `SuggestList` | §11e-2 / §23d | 4 berkas | **4** — tidak berubah |
| Pemeriksaan total share `@greaterThan(Local.TotalPersen,100)` | §6e / P15 | 1 kemunculan | **1** — tidak berubah, tetap **penjaga sebelah** |
| `pyActionConditions` berisi anak | §13 R7-baru / §17c | 3 Section, 17 baris | **3 Section, 17 baris** — tidak berubah (`Section/AdjustmentDetail.xml` 10 · `Section/OutstandingClaim.xml` 6 · `Section/InputAcceptation_Est.xml` 1; `Section/AdjustmentDetail_Section.xml` dan `FlowAction/AdjustmentDetail.xml` nol) |

**Alasan semuanya tetap:** ke-58 berkas baru seluruhnya di `When/` dan `RDBList/`, sedangkan ketujuh
sensus di atas menyasar `Activity/`, `Section/`, dan `FlowAction/`.

### 27b. ⚠️ Sensus yang **BERUBAH** — §20a

| Sensus | Angka lama (270) | **Angka baru (328)** |
| --- | --- | --- |
| Rule `When` dirujuk tetapi **tidak ada** di `When/` | **5** — `IsFire` · `isGolfInsurance` · `IsAneka` · `IsMarineCargo` · `IsSpreadingUW` | **0** |

`[terverifikasi]` Dicari di **60 dari 60 berkas** folder `Claim Prop/When/`: **kelimanya kini ADA**.
**§20a tertutup** — ini yang menutup **P19**.

### 27c. Sensus yang penyebutnya **tetap**

Sensus yang hanya menyisir `Activity/` tidak terpengaruh: penyebutnya tetap **112 dari 112**
(mis. 84 gerbang tidak berlaku, 68 step di-remark, 513 baris syarat aktif di §24). Sensus yang hanya
menyisir `Section/` tetap **36 dari 36** (57 gerbang UI mati, 75 elemen selalu-salah di §22).

---

### 27d. Empat sensus yang WAJIB dijalankan ulang — hasil lama versus baru

| Sensus | Bab | **Lama (270 / 54 / 6)** | **Baru (328 / 58 / 60)** |
| --- | --- | --- | --- |
| `RequestType` menggantung | §8d | **7** menggantung | **3** menggantung |
| Rule yatim di `RDBList/` | §8d | **0 dari 54** | **0 dari 58** — rekor bersih bertahan |
| Nilai hardcode di SQL | §8f | **17 dari 54** rule | **18 dari 58** rule |
| Nama rule yang berbohong | §8g | daftar lama | **+1 kasus baru** |
| Rule `When` dirujuk tetapi tidak ada | §20a | **5** | **1** |
| `COMMIT;` di dalam `<pyBrowseSQL>` | §8b | **13 dari 54** | **15 dari 58** |
| Silsilah *save-as* | §8c | **21 dari 54** (39 %) | **23 dari 58** (39 %) |
| Inventaris stored procedure | §3c · §26c | **10** procedure + 2 fungsi + 1 sequence | **12** procedure + **2** fungsi + **1** sequence |

> ### ⚠️ RALAT 2026-09-17 — X3: verifikasi ulang keempat sensus atas 58 dari 58
>
> Dijalankan mekanis, bukan diambil dari laporan penyisiran:
>
> | Sensus | Lama | **Baru** | Berubah? |
> | --- | --- | --- | --- |
> | `COMMIT;` di `<pyBrowseSQL>` (§8b) | 13 dari 54 | **15 dari 58** | ✅ **ya, +2** — kedua generator nomor baru |
> | Silsilah *save-as* (§8c) | 21 dari 54 · 39 % | **23 dari 58 · 39 %** | ✅ ya, +2 — persentase bertahan |
> | Nilai hardcode di SQL (§8f) | 17 dari 54 | **18 dari 58** | ✅ ya, +1 — `GetTreatyInMaster_SQL` (`'Proportional'` 2×) |
> | Nama yang berbohong (§8g) | daftar lama | **+1 kasus** | ✅ ya — `GetTreatyInMaster_SQL` (§26d) |
> | Inventaris procedure (§3c) | 10 + 2 + 1 | **12 + 2 + 1** | ✅ **dipastikan** — `GENERATE_NODLATREATYIN`, `GENERATE_NOPLATREATYIN` |
>
> `[terverifikasi]` Pemastian §3c: sensus procedure/fungsi atas **58 dari 58** memberi **12
> procedure** (`GENERATE_NOCLMTREATYIN` · `GENERATE_NOCLMTRTYINTEMP` · `GENERATE_NODLATREATYIN` ·
> `GENERATE_NOPLATREATYIN` · `GET_TOKEN_STORAGE` · `PEGA_D_CAUSE_OF_LOSS` · `PEGA_JSON_KLAIM_PNC` ·
> `PEGA_JSON_OS_AKSEP_KLAIM` · `_KLAIMTNP` · `_KLAIMTRT` · `PEGA_M_CAUSE_OF_LOSS` ·
> `PROC_GENERATE_SEQUENCE_NUMBER`), **2 fungsi** (`pooldata.getcurrencystandard` · `gl.f_get_email`),
> **1 sequence** (`ADJUSTERCONSULTANT_SEQ`).
>
> ⚠️ `POOLDATA.MONITORING_KLAIM_LOG` tertangkap regex sebagai procedure — ia **tabel**, terjaring
> karena daftar kolom `INSERT INTO`. **Tidak masuk inventaris.**

#### §8d — `RequestType` menggantung turun 7 → 3

`[terverifikasi]` Dari **61 nama `RequestType` unik** yang dirujuk di 328 berkas, tersisa **3** tanpa
berkas:

| Nama dirujuk | Pemanggil | Status |
| --- | --- | --- |
| `GetBreakDownTreaty` | `Activity/SetTreatyNameSpreading_Act.xml` step 10 | **DI-REMARK** (§1a) → tidak berdampak |
| `InsertCatastrope_SQL` | `Activity/SaveCatasrtope_Act.xml` step 2 | **DI-REMARK** (§1a) → tidak berdampak |
| `SaveMasterAdjusterConsultant_sql` | `Activity/SaveAdjusterConsultant_Act.xml` step 5 | **DI-REMARK** (§1a) → tidak berdampak |

⚠️ **Ketiganya sudah terjelaskan sejak §1a.** Empat yang hilang tepat sejumlah berkas baru, dan
**nol nama menggantung baru** muncul bersama 58 berkas itu. **§8d praktis tertutup** — yang tersisa
hanya rule yang pemanggilnya memang mati.

#### §8f — hardcode 17 dari 54 → **18 dari 58**

`[terverifikasi]` Hanya **satu** berkas baru menambah literal: `RDBList/GetTreatyInMaster_SQL.xml`
dengan **`'Proportional'`** (2×, keduanya di `WHERE`). Tiga berkas baru lainnya **bersih**.
Yang paling patut diawasi **tetap** `GetLimitDirekturUtama_SQL.xml` (nama orang) dan
`GetLimitPLATreatyin.xml` (`'10001'`).

#### §8g — satu kasus baru

`RDBList/GetTreatyInMaster_SQL.xml` — lihat **§26d**.

#### §20a — rule `When` hilang turun 5 → **1**, tetapi yang tersisa **baru**

`[terverifikasi]` Penyebut: **60 berkas** di `When/` · **47 nama `When` berbeda** dirujuk di 328
berkas. **Kelima nama lama kini ADA semua.** Namun muncul **satu nama menggantung baru**:

| Nama dirujuk | Dirujuk dari | Masalah |
| --- | --- | --- |
| **`IsCustomBonds`** | `When/IsBondingAndCustomBonds.xml` | folder hanya punya **`IsCustomBond`** — **tanpa huruf `s`** |

⚠️ **Hampir pasti salah ketik, bukan rule yang benar-benar hilang.** Bukti dari arah berlawanan:
**`When/IsCustomBond.xml` justru YATIM** — tidak pernah dirujuk siapa pun. Pasangan "satu nama
dirujuk tanpa berkas / satu berkas tanpa perujuk, beda satu huruf" praktis memastikannya. Resolusi
nama rule Pega **tidak peka huruf besar-kecil**, tetapi ini **beda ejaan**, bukan beda kapitalisasi
— jadi tidak tertolong. → `[terbuka]` ringan; perbaikannya sepele.

### 27e. ⚠️ Catatan metode — dua jebakan sensus yang terbukti

`[terverifikasi]` **Jebakan 1 — kapitalisasi.** Dua penyisiran memberi daftar yatim `When/` yang
berbeda pada satu nama. Sebabnya: `When/IsAneka.xml` merujuk `IsAllRisk` dengan ejaan
**`isAllRisk`** (huruf `i` kecil). Sensus **peka huruf besar-kecil** melewatkannya dan salah
melabelinya yatim.

**Verifikasi saya sendiri, case-insensitive:**

| Nama | Rujukan | Putusan |
| --- | ---: | --- |
| `isAnalistorTransfer` | **0** | **yatim** |
| `IsAllRisk` | **1** — dari `When/IsAneka.xml` | **BUKAN yatim** |
| `IsCustomBond` | **0** | **yatim** |

**Jumlah yatim `When/` yang benar tetap 14 dari 60** (daftar di §25d). Sensus nama rule **wajib
case-insensitive**.

`[terverifikasi]` **Jebakan 2 — SQL satu baris.** `RDBList/GetDataNopolisTreatyin.xml` menaruh
seluruh SQL dalam **satu baris**, sehingga ekstraksi berbasis baris menyeret sisa berkas dan ikut
menangkap `<pxCommitDateTime>` — persis jebakan yang sudah dicatat di **§8b**. Dengan ekstraksi
multiline yang benar: **nol `COMMIT`** di berkas itu.

---


---

## 28. Aturan wewenang Claim Prop

`[keputusan work owner]` 2026-09-17. Bab ini menutup **P24** dan mengubah **P24 butir 1**.
Ia menggantikan tiga catatan terpisah yang sebelumnya tersebar di §11e, §13a, dan §19.

### 28a. Tangga empat tingkat

```
Claim Admin              ← nilai DASAR
   ↓
Claim Dept. Head         ┐
Operational Director     ├─ anggota komite
Technical Director       ┘
```

`[terverifikasi]` Keempat nilai muncul di `DataTransform/InsertChronology_DT.xml`
(langkah **1.1.4** · **1.1.5** · **1.1.6** · **1.1.7**) dan di `Activity/SethistoryKlaimTreaty.xml`
(blok berulang `RH_1.pySteps(1)`).

### 28b. Sumber tunggal: `POOLDATA.EMAILKOMITE`

`[keputusan work owner]` Tingkat wewenang dibaca dari **satu sumber**:

```
cari pelaku di EMAILKOMITE
  ketemu        → tingkat = DEGREE
  tidak ketemu  → tingkat = Claim Admin
```

`POOLDATA.EMAILKOMITE.DEGREE` **berpadanan** dengan tingkat di `.IsCedingConfirm` (butir b jawaban
P24). Kolom `DEGREE` sudah tersedia dan sudah diambil
`ReportDefinition/FilterEmailKomiteWithLimit.xml` (§16c), jadi **tidak menuntut perubahan skema**
untuk bagian ini.

### 28c. Aturan nilai dasar

`Claim Admin` **bukan anggota komite** dan **tidak ada di `EMAILKOMITE.DEGREE`** — memang tidak
seharusnya ada. Ia adalah **default untuk siapa pun yang tidak ditemukan di roster**.

⚠️ Konsekuensi: pencarian yang **tidak menemukan baris** di sini **bukan kegagalan** — ia hasil yang
sah, berbeda dari pencarian batas nilai di §28d yang justru harus gagal terang-terangan (Q7, §21).

### 28d. Direktur Utama **di luar** tangga

`[keputusan work owner]` `RDBList/GetLimitDirekturUtama_SQL.xml` adalah **wewenang terpisah dengan
urusannya sendiri** — **bukan tingkat kelima**, dan **tidak berpadanan** dengan tingkat mana pun di
jejak audit. Ia menggerbangi **kewajiban dokumen ADU** lewat
`Activity/AttachmentProtect_ACT.xml` step **5** (§13a), bukan tangga persetujuan.

### 28e. ⚠️ Nama orang yang di-hardcode **TIDAK ikut ke Go**

> #### ⚠️ RALAT 2026-09-17 — **P24 butir 1 BERUBAH**
>
> Jawaban sebelumnya (§19, U1) adalah **"ikuti apa adanya — nama tetap di kode"**. Itu diambil
> **sebelum** diketahui bahwa roster dan jejak audit adalah **tangga yang sama** (butir b).
> Dengan informasi itu, `[keputusan work owner]` **mengubahnya**.

`[keputusan work owner]` ⚠️ **Penyimpangan sadar.** **Empat nama yang di-hardcode tidak ikut ke Go** —
tiga di `DataTransform/InsertChronology_DT.xml` (langkah 1.1.5–1.1.7) dan satu di
`RDBList/GetLimitDirekturUtama_SQL.xml`. Nilainya tidak disalin ke sini (PII).

**Perilakunya tidak berubah** — orang yang sama tetap mendapat tingkat yang sama — **tetapi saat
orangnya berganti, cukup ubah baris tabel.** Ini menutup cacat yang dicatat di §11e-4 butir 9c:
*"orangnya pindah jabatan → jejak audit salah label tanpa peringatan"*.

### 28f. ⚠️ Kunci pencocokan — butuh konfirmasi DBA

`[terverifikasi]` `.PICSuggest` diisi `OperatorID.pyUserName`, dan ketiga nilai hardcode
**bentuknya tidak seragam**: satu tanpa spasi huruf besar semua · satu kapital awal · satu huruf
besar semua. **Mencocokkan lewat `NAME` rapuh.**

`[keputusan work owner]` Pakai kolom baru **`USER_ID`** di `POOLDATA.EMAILKOMITE`, diisi user ID
Pega/akun. **Bila ternyata email akun sama dengan kolom `EMAIL` yang sudah ada, kolom `EMAIL` boleh
dipakai sebagai gantinya dan kolom baru tidak perlu.**

> **`[data DBA]` terbuka ringan — bukan pemblokir, tetapi WAJIB dikonfirmasi sebelum tiket ditulis:**
>
> **Apakah akun Pega memakai email yang sama dengan `EMAILKOMITE.EMAIL`?**
> Bila **ya** → pakai `EMAIL` sebagai kunci. Bila **tidak** → tambah kolom `USER_ID`.

---

