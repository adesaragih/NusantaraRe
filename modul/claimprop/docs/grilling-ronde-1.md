# Grilling — Claim Prop, Ronde 1

**Tanggal:** 2026-09-17 · **Modul:** `Claim Prop` (270 berkas) · **Status:** frontier BELUM kosong

> Korpus `D:\XML\RNM_BRD\` **READ-ONLY**. Semua tulisan hanya ke `OUTPUT_HASIL_RNM\`.
> ⛔ `D:\XML\nusantara-re\` di-blacklist — tidak dibaca, tidak dirujuk.
> Identitas rule = `class/nama` dari `<pxInsName>`; tipe dari field pertama `<pzOriginalInstanceKey>`.
> Tanda: `[terverifikasi]` `[keputusan work owner]` `[data DBA]` `[terbuka]` `[dugaan]`.
> ⚠️ **Tidak ada PII di berkas ini.** Korpus memuat nama orang, alamat email, dan satu operator ID
> tertanam; seluruhnya dirujuk secara deskriptif, nilainya tidak pernah disalin.

---

## 0. Ralat terhadap artefak discovery — tiga klaim "tidak terbaca" ternyata SALAH

Discovery D2/D3 menyatakan kondisi `When` di modul ini **tidak terbaca dari tag**, karena
`<pyLabel>` hanya berisi template kosong `[first value][relation][second value]`. Itu keliru:
kondisinya tersimpan di **tag lain**. Metode yang benar: `<pyConditionString>` untuk kondisi
ekspresi, `<pyParametersParamValue>` untuk kondisi "simple".

| Rule | Kondisi sebenarnya `[terverifikasi]` | Bukti |
| --- | --- | --- |
| `IsPEGAPROD` | `pxProcess.pzProductionLevel = "5"` | `Claim Prop/When/IsPEGAPROD.xml` baris **174**; ditegaskan baris **325** `compareTwoValues(pxProcess.pzProductionLevel, "=", "5")` |
| `IsBackStage` | `.pyNote = "Back"` | `Claim Prop/When/IsBackStage.xml` — `<pyDesignatedProperty>pyNote`, `compareTwoValues(.pyNote, "=", "Back")` |
| `IsCLM` | `pyWorkCover.pyWorkIDPrefix = "CLM-"` **atau** `pyWorkPage.pyWorkIDPrefix = "CLM-"` | `Claim Prop/When/IsCLM.xml` — 2 kondisi (label A, B) |
| `IsCLMP` | idem, `"CLMP-"` | `Claim Prop/When/IsCLMP.xml` |
| `IsCLMNP` | idem, `"CLMNP-"` | `Claim Prop/When/IsCLMNP.xml` |
| `IsPEGASyariah` | `pxProcess.pxSystemNodeID = "jboss1074"` | `Claim Prop/When/IsPEGASyariah.xml` baris **299 / 306 / 352** |

⚠️ `[terverifikasi]` **Label berbohong lagi:** `IsBackStage` menampilkan label `Notes`, tetapi
properti yang benar-benar diuji adalah `pyNote`.

**Akibatnya:**

- `[terverifikasi]` **OQ-029 dapat ditutup untuk Claim Prop** — kondisinya identik dengan yang sudah
  terbukti di Claim Life (`pzProductionLevel = "5"`). Dugaan tidak lagi diperlukan.
- **OQ-041 separuh tertutup** — "apa yang diuji keempat rule ini" kini terjawab. Yang tersisa hanya
  "mengapa tiap modul punya versinya sendiri" (butuh perbandingan lintas modul).

`[terverifikasi]` Pembanding negatif yang menguatkan bahwa ini bukan kebetulan: `<pyConditionString>`
**kosong** untuk keempat `ISCLM*`/`IsPEGASyariah` tetapi **terisi** untuk `IsPEGAPROD` dan
`IsBackStage` — Pega menyimpan kondisi "simple" dan kondisi "expression" di tag yang berbeda.

---

## 1. Temuan struktural utama — modul ini melayani TIGA prefix klaim, bukan satu

`[terverifikasi]` Rule **di dalam** `Claim Prop` bercabang atas `pyWorkIDPrefix` milik modul lain.

| Prefix | Muncul | Pasangan berakhiran `S` |
| --- | ---: | --- |
| `"CLMNP-"` | 18× | `"CLMNPS-"` 2× |
| `"CLM-"` | 16× | `"CLMS-"` 4× |
| `"CLMP-"` | 12× | `"CLMPS-"` 2× |

Bukti percabangan lintas-prefix:

- `Claim Prop/Activity/HitServiceToKasir_Act.xml` baris **617**:
  `(pyWorkPage.pyWorkIDPrefix=="CLMP-" && .Type=="3") || (pyWorkPage.pyWorkIDPrefix=="CLM-" && .PaymentType=="3") || (pyWorkPage.pyWorkIDPrefix=="CLMNP-" && .PaymentType=="3")`
- `Claim Prop/Activity/HitServiceToKasir_Act.xml` baris **2360**:
  `@if(pyWorkPage.pyWorkIDPrefix="CLMNP-", pyWorkPage.ClaimData.AdjustmentList(1).CurencyAdjustment …)`
- `Claim Prop/Activity/CekPremiLunas_Act.xml` baris **529** (`=="CLMP-"`) dan **717** (`=="CLM-"`)
- `Claim Prop/Activity/SendEmailKlaimRejectClose.xml` baris **739, 930, 942, 945, 1394, 1644, 1915** —
  selalu berpasangan `X- || XS-`

⚠️ **Dua konsekuensi yang bertabrakan dengan peta konteks:**

1. Peta konteks (§2.5, §4) menyatakan trio Claim non-life adalah **tiga basis kode terpisah**
   (tumpang tindih identitas 12–40 %). Itu benar untuk **identitas rule**, tetapi **perilaku
   runtime**-nya tidak terpisah: beberapa rule di modul ini mengeksekusi cabang untuk `CLM-` dan
   `CLMNP-`. `[terbuka]` Apakah 102 identitas yang dibagi dengan Claim Non Prop justru rule-rule
   inilah — **belum diverifikasi**.
2. `[terverifikasi]` Ada **enam** prefix, bukan tiga. Tiga berakhiran `S`.
   `[dugaan — JANGAN dipakai sebagai fakta]` `S` = Syariah, karena `IsPEGASyariah` ada di modul ini.
   Korpus **tidak memuat** label atau kepanjangan mana pun → **`[terbuka]`**.

### 1a. `.Type` di Claim Prop **bukan** `.Type` di Claim Life

`[terverifikasi]` Sensus literal di `Claim Prop`: `.Type` dibandingkan dengan **`1`, `2`, `3`, `4`**
(numerik). Di Claim Life `.Type` bernilai `QR`/`QP`/`TP`/`TR`. **Nama properti sama, enum berbeda.**
Keputusan Life (ADR-0012, `TP`=Payable / `TR`=Receivable) **tidak boleh** dibawa ke sini.

`[terverifikasi]` Lebih tajam lagi, baris 617 memperlihatkan **asimetri**: untuk prefix `CLMP-`
gerbangnya `.Type=="3"`, sedangkan untuk `CLM-` dan `CLMNP-` gerbangnya `.PaymentType=="3"` —
**properti berbeda, literal sama**. Apakah ini sengaja (dua enum berbeda) atau cacat (satu enum,
dua nama) **tidak terbaca dari korpus** → `[terbuka]`, menyambung **OQ-020**.

`[terverifikasi]` Sensus `.PaymentType` di modul ini: `1, 2, 3, 4, 5, 6`.

### 1b. `IsPEGASyariah` digerbangi **identitas node aplikasi**

`[terverifikasi]` `pxProcess.pxSystemNodeID = "jboss1074"` — bukan data bisnis, melainkan **nama
node JBoss**. Satu-satunya perujuk di modul ini: `Claim Prop/Activity/HitServiceToKasir_Act.xml`.
Di Go tidak ada padanan identitas node. Ini kembaran ADR-0005 (`IsPEGAPROD` → flag lingkungan) dan
menyentuh **OQ-018**.

⚠️ Jangan keliru dengan `<pxHostId>` — itu metadata export (sebaran: `pega-nusre` 156×,
`jboss117` 50×, `jboss122117` 49×, dan empat hash lain). **Bukan** temuan.

### 1c. Pemakai keempat `When`

| `When` | Berkas perujuk `[terverifikasi]` |
| --- | --- |
| `IsPEGAPROD` | 6 — a.l. `getStatusKonversi_Act`, `HitServiceToKasir_Act`, `InsertGoogleStorage_Act`, `KonversiKlaim_Act` |
| `IsCLM` / `IsCLMP` / `IsCLMNP` | 2 masing-masing — `HitServiceToKasir_Act`, `SendEmailKlaim` |
| `IsPEGASyariah` | 1 — `HitServiceToKasir_Act` |

Artinya keempatnya menggerbangi **efek keluar**, bukan alur klaim inti.

---

## 2. Penyerahan ke Komite — DUA jalur, dan keduanya tidak saling mengunci

### 2a. F1 — kenapa muatan digemukkan: **kontrak balik-arah beralamat POSISI**

`[terverifikasi]` `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml`
(`ASM-FW-GCNMFW-DATA-ADJUSTMENT` / `ADDKOMITETREATYCHILD_ACT` / `RULE-OBJ-ACTIVITY`).

Penunjuk yang menyeberang:

| Properti | Nilai | Baris | Sifat |
| --- | --- | ---: | --- |
| `childPageKomite.CLMNO` | `pyWorkPage.pyID` | 3390 | **stabil** |
| `childPageKomite.IndexAdjustment` | `.pxListSubscript` | 3222 | **indeks posisi** |
| `childPageKomite.Adjustment.CoverageSubscript` | `.IndexCoverage` | 3201 | indeks posisi |

`[terverifikasi]` **Nol** `Property-Set` terhadap ID baris adjustment — tidak ada `AdjustmentID`,
GUID, atau kunci alternatif apa pun.

`[terverifikasi]` Sisi Komite **menulis balik lewat subscript itu**:
`Komite Claim Prop/Activity/KomitePostAdjustment.xml` baris **904** (`Local.AdjusmentID =
.IndexAdjustment`), lalu 12 penulisan ke `TempOpenPage.ClaimData.AdjustmentList(Local.AdjusmentID).*`
(a.l. baris 2888 `.IsApproved`, juga 1400, 1447, 4500, 5623, 6603). Induk dibuka dengan
`Obj-Open-By-Handle` atas `pyWorkPage.pxCoverInsKey` (baris 1027 / 1076).

**Jawaban F1:** muatan digemukkan **karena alamat baliknya rapuh**. Komite tidak dapat memverifikasi
ulang baris mana yang dimaksud, sehingga 24 field isi baris ikut disalin sebagai **snapshot**. Jadi
ini **bukan** salinan defensif yang disengaja, dan **bukan** karena Komite tak bisa membaca balik —
Komite **bisa** (dan memang melakukannya), tetapi hanya lewat indeks posisi.

### 2b. F3 — jalur kedua: `SendCloseClaimToKomite`

`[terverifikasi]` `ASM-FW-GCNMFW-WORK-CLAIMTREATY` / `SENDCLOSECLAIMTOKOMITE` / `RULE-OBJ-ACTIVITY`.
13 step. Membuat case anak **kelas dan flow yang identik** dengan jalur pertama —
`ChildClass = ASM-FW-GCNMFW-Work-KomiteTreaty`, `FlowName = KomiteTreaty_Flow` (baris 1803–1804;
bandingkan 6147–6149 di jalur pertama).

| | `AddKomiteTreatyChild_ACT` | `SendCloseClaimToKomite` |
| --- | --- | --- |
| `Property-Set` ke `childPageKomite.*` | **49** | **14** |
| Field `Adjustment.*` | **24** | **0** |
| `TransferType` | **`2`** | **`4`** |
| Roster | `FilterEmailKomiteWithLimit`, `Param.STS_KLAIM = "PROP"`, `Param.LIMIT_BOTTOM = Local.TotalAdjustment` | ⚠️ **hardcode satu orang**, `KomiteLoop = 1`, `KomiteCount = 1` |
| Dipanggil dari | `Harness/CommitteeTreaty.xml` (2×), `Section/ComiteeClaimTreaty.xml` (2×) | `Section/PreventRejectClaimProp.xml` baris 3077 & 3156, tombol "Yes" |

`[terverifikasi]` Muatan F3 adalah **himpunan bagian murni** dari muatan F1 — nol properti baru.

⚠️ `[terverifikasi]` **Jalur close TIDAK punya tangga roster sama sekali.** Ia tidak memanggil
`FilterEmailKomiteWithLimit`; satu jabatan, satu alamat email, dan satu operator ID **tertanam di
dalam rule** (baris 1607–1668; nilai tidak disalin ke sini — PII). Ditegaskan dari sisi Komite:
`Komite Claim Prop/Activity/SetKomiteList_Act.xml` baris 344 bergerbang `.TransferType==2`, dengan
memo baris 12 berbunyi `skip selain .TransferType =2`.

⚠️ `[terverifikasi]` Alamat email penerima di `SendEmailKlaimRejectClose.xml` juga **tertanam**
(baris 615, 806, 2670) — bukan lookup ke master.

### 2c. F2 — `TransferType` **HIDUP** di Claim Prop, dan ia adalah dispatcher

`[terverifikasi]` Di `Claim Prop`: **2 tulis, 1 baca.**

| Berkas | Baris | Peran |
| --- | ---: | --- |
| `Activity/AddKomiteTreatyChild_ACT.xml` | 3180 | tulis literal **`2`** |
| `Activity/SendCloseClaimToKomite.xml` | 1501 | tulis literal **`4`** |
| `Activity/SendEmailKlaimRejectClose.xml` | 507 | baca: `@if(childPageKomite.TransferType = "3", "Reject", "Close")` |

`[terverifikasi]` Di `Komite Claim Prop`: **20 pembacaan**, dan ia menggerbangi alur.
`Komite Claim Prop/Activity/KomitePost.xml` adalah switch murni:

| Nilai | Cabang | Baris |
| ---: | --- | ---: |
| `2` | `Call KomitePostAdjustment` | 281 / 225 |
| `3` | `Call KomitePost_Reject` | 379 / 319 |
| `4` | `Call KomitePost_Close` | 479 / 419 |

Routing assignment: `KomiteRouter.xml` baris **1190** (`=='2'`) dan **1343** (`!='2'`).

⚠️ **Konsekuensi mengikat:** kesimpulan Claim Life bahwa `TransferType` adalah kode mati
**TIDAK berlaku di sini**. Di Claim Prop ia kontrak lintas-case yang menentukan cabang mana yang
dijalankan Komite.

⚠️ `[terverifikasi]` **Cabang tak tercapai:** pembacaan lokal di `SendEmailKlaimRejectClose:507`
menguji `= "3"`, tetapi **tidak satu pun** dari kedua jalur di Claim Prop menghasilkan `3` — maka
subjek email "Reject" tak pernah terpilih lewat jalur-jalur ini. Siapa yang menghasilkan `3`
**tidak terbaca dari korpus `Claim Prop`** → `[terbuka]`.

### 2d. ⚠️ Tidak ada yang mencegah kedua jalur aktif pada klaim yang sama

`[terverifikasi]` Tiga bukti terpisah:

1. Tombol pemanggil di `Section/PreventRejectClaimProp.xml` punya
   `<pyActionConditions REPEATINGTYPE="PageList"/>` **kosong** (baris 3096, 3175) — tanpa kondisi UI.
2. `AddKomiteTreatyChild_ACT.xml` baris **926** men-set
   `pyWorkPage.ClaimData.FlagOnGoingCommitte = "Send Commite"` — tampak seperti flag "sedang di
   komite", tetapi grep seluruh `Claim Prop` menemukannya **hanya di baris itu**.
   **Flag tulis-saja, tidak pernah dibaca.**
3. Di `SendCloseClaimToKomite`, precondition `.AcceptanceStatus==""` (baris 722) memasang
   `Local.Error = "Can not close claim, there is adjustment in comitee!"`, tetapi
   `<pyStepsTransition/>` pada `Page-Set-Messages` (baris 807) **kosong** — eksekusi **tetap lanjut**
   ke `Page-New` (992) dan `Call pxAddChildWork` (1747). Yang tercegah hanyalah pengiriman email
   (gerbang baris 2638).

**Artinya perilaku existing memungkinkan dua case anak `Work-KomiteTreaty` di bawah satu induk,
dengan `TransferType` 2 dan 4 sekaligus** — pesan error terpasang tetapi tidak menghentikan apa pun.

⚠️ `[terverifikasi]` Gerbang `pxAddChildWork` di jalur close adalah `ParamCari.CARI2==1`
(baris 1961), dan `ParamCari.CARI2` **tidak pernah di-set di berkas itu** — asalnya **tidak terbaca
dari korpus** → `[terbuka]`.

### 2e. Batas transaksi jalur Komite

`[terverifikasi]` `AddKomiteTreatyChild_ACT`: **1** `Obj-Save` (baris 7317) dengan `WriteNow=true`,
`WithErrors=true`; **tidak ada** step `Commit`. `Call pxAddChildWork` membawa parameter
`Commit=true` (baris 6147). Pega sendiri memasang peringatan di berkas yang sama (baris 9636):
*"Obj-Save Write Now … interferes with the PRPC transaction model and may result in stale database
records"*.

`[terverifikasi]` `SendCloseClaimToKomite`: **1** `Obj-Save` (baris 2403, `WriteNow=true`,
`WithErrors=false`), juga tanpa `Commit`.

---

## 3. Silsilah rule — `pzOriginalInstanceKey` menyimpan **asal save-as**, bukan identitas

`[terverifikasi]` Pada 3 dari 4 rule penomoran, nama di `<pzOriginalInstanceKey>` **berbeda** dari
`<pxInsName>`, padahal tiap berkas hanya memuat **satu** instance (diverifikasi: `pxInsName` 1×,
`pzOriginalInstanceKey` 1× per berkas):

| Berkas | `pxInsName` (identitas) | nama di `pzOriginalInstanceKey` (asal) |
| --- | --- | --- |
| `RDBList/GenerateNoCLMTreatyIn.xml` | `…!GCNM!GENERATENOCLMTREATYIN` | `GCNM!GENERATENOPLATREATYIN` |
| `RDBList/GetSequenceNumber_SQL.xml` | `…!RNM!GETSEQUENCENUMBER_SQL` | `RNM!GETTANGGALCLOSING_SQL` |
| `RDBList/GetKodeProdNonLife_SQL.xml` | `…!RNM!GETKODEPRODNONLIFE_SQL` | `RNM!GETSEQUENCENUMBER_SQL` |
| `RDBList/GETTanggalClosing_SQL.xml` | `…!RNM!GETTANGGALCLOSING_SQL` | (cocok) |

Pola sama di dua activity Komite: `AddKomiteTreatyChild_ACT` berasal dari `ADDCHILDANDSAVE_ACT`;
`SendCloseClaimToKomite` berasal dari kelas **`ASM-FW-GCNMFW-Work-PNC`**, bukan `ClaimTreaty`.

**Konsekuensi metodologis:** konvensi yang berlaku sudah benar — identitas **hanya** dari
`<pxInsName>`, dan `<pzOriginalInstanceKey>` dipakai **hanya** untuk field tipe. Mencari rule
berdasarkan nama di `pzOriginalInstanceKey` akan menyesatkan.

---

## 4. Penomoran — dua prosedur, salah satunya "temp"

`[terverifikasi]` Lima rule Connect-SQL, isinya:

| Rule (`Claim Prop/RDBList/`) | Procedure / objek |
| --- | --- |
| `GenerateNoCLMTreatyIn.xml` | `BEGIN POOLDATA.GENERATE_NOCLMTREATYIN` |
| `GenerateNoTRTInTemp.xml` | `BEGIN POOLDATA.GENERATE_NOCLMTRTYINTEMP` |
| `GetSequenceNumber_SQL.xml` | `BEGIN POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` |
| `GetKodeProdNonLife_SQL.xml` | `SELECT` atas `POOLDATA.KODE_PRODUKSI` |
| `GETTanggalClosing_SQL.xml` | `SELECT` atas `POOLDATA.TANGGAL_CLOSING` |

`[terverifikasi]` `PROC_GENERATE_SEQUENCE_NUMBER` dan `TANGGAL_CLOSING` **sama** dengan yang dipakai
Komite Claim Life — mesin sequence dipakai bersama lintas lini. Prefix produksi di sini di-lookup
lewat rule bernama `…NonLife`, sejalan **ADR-0006** (jangan replikasi logika penomoran).

⚠️ `[terbuka]` **Mengapa ada nomor "temp"** (`GENERATE_NOCLMTRTYINTEMP`) di samping nomor final
(`GENERATE_NOCLMTREATYIN`)? Apakah klaim memperoleh nomor sementara lebih dulu lalu nomor tetap?
Kapan promosinya, dan apa yang terjadi pada nomor temp bila klaim ditolak? Body kedua procedure
**tidak ada di korpus** → butuh **dump DBA** (OQ-002).

---

## 5. OQ-021 — verifikasi NEGATIF untuk Claim Prop

`[terverifikasi]` Tidak ditemukan guard identitas literal: pola `pyUserIdentifier = "<literal>"`
menghasilkan **nol** hasil di seluruh `Claim Prop`. Delapan kemunculan
`OperatorID.pyUserIdentifier` seluruhnya adalah **nilai yang ditulis** (stempel siapa yang
membuat/mengubah): `Activity/AddAdjustment_Act.xml` (877), `Activity/SaveOutstanding_Act.xml` (4736),
`Harness/TambahCauseofLoss.xml`, `Harness/TambahMasterCauseOfLoss.xml`.

⚠️ Nama orang di korpus (`<pxCreateOpName>`, `<pxCreateOperator>`, `<pxUpdateOpName>`,
`<pxMoveImportOperName>`, `<pxWarningCreateOperator>`) adalah **metadata penulis rule**, bukan
gerbang otorisasi.

⚠️ **Tetapi** PII yang sebenarnya ada di tempat lain: alamat email dan operator ID **tertanam
sebagai data runtime** di `SendCloseClaimToKomite.xml` dan `SendEmailKlaimRejectClose.xml` (§2b).
Itu bukan metadata — itu konfigurasi yang tertanam di kode, dan **wajib** menjadi data master di
sistem baru.

---

## 6. Estimation — siklus hidup

### 6a. Bentuk data `[terverifikasi]`

**Page list**, bukti definisional `<pyPageListProperty>` + `<pyPageListPropertyClass>`:

| Properti | Class elemen | Bukti |
| --- | --- | --- |
| `.ClaimData.EstimationList` | `ASM-FW-GCNMFW-Data-Estimasi` | `Section/InputAcceptation_Est.xml` 26382/26373; `Section/OutstandingClaim_Est.xml` 26597/26588 |
| `.ClaimData.ListTotalEstimation` | `ASM-FW-GISFW-Data-TreatyInTotal` | `InputAcceptation_Est.xml` 32239/32230; `OutstandingClaim_Est.xml` 32261/32252 |

`[terverifikasi]` **1:N murni** — `EstimationList(<APPEND>)` sekali lalu `(<LAST>)` untuk kolom sisanya
(`Activity/AddEstimation_Act.xml` baris 2322–2468); subscript variabel `(Param.Index)` / `(Local.Index)`;
**nol** indeks angka literal.

`[terverifikasi]` **Induk struktural = header klaim** (`pyWorkPage.ClaimData.EstimationList`), sejajar
`AdjustmentList` dan `SpreadingRisk` — **tidak** menggantung pada baris adjustment maupun coverage.

⚠️ `[terverifikasi]` **Granularitas semantiknya per-TREATY, bukan per-klaim.** `AddEstimation_Act`
baris 1894/1902 beriterasi atas `LossAlloc.pxResults` (*"looping temp dan append data to estimation"*),
dan tiap iterasi meng-append satu baris Estimation yang identitasnya diambil dari hasil **loss
allocation**. Jadi **Loss Allocation adalah hulu Estimation** — dua konsep yang tidak ada di Life
ternyata berantai.

### 6b. ⚠️ Dua nama kolom berbohong

| Kolom | Namanya menyiratkan | `[terverifikasi]` isinya |
| --- | --- | --- |
| `.GrossEstimationPct` | persentase | **nilai uang** — diisi dari `.ClaimSpreaded` (`AddEstimation_Act:2368`) lalu dikalikan `KursValue` untuk IDR |
| `.TypeLoss` / `.TypeLossID` | jenis kerugian | **identitas treaty** — diisi dari `.TreatyName` / `.TreatyType` (`AddEstimation_Act:2448, 2468`) |

Kolom turunan (`Activity/CountEstimation_Act.xml`, step page class `Data-Estimasi`, baris 717):
`.PersenRNM` (727), `.EstimationValue = RNMShareP/100 × GrossEstimationPct` (773),
`.ConvertValue = EstimationValue × KursValue` (793),
`.ConvertGrossEstimasi = GrossEstimationPct × KursValue` (813).
⚠️ `.EstimastionReserve` (964) — **typo permanen** di nama properti korpus.

### 6c. ⚠️⚠️ Invariant "satu klaim satu mata uang" **TIDAK BERLAKU** di Claim Prop

Ini membalik asumsi yang berlaku di Claim Life (OQ-060). `[terverifikasi]` tiga bukti bebas:

1. `Activity/CurencyEstimation_Act.xml` menulis `CurrencyID` (296), `Currency` (764), dan `KursValue`
   (1114) ke **baris tertentu saja** lewat `(Param.Index)` — mata uang ditetapkan **per baris**.
2. Baris 1456 menyelaraskan kurs antar-baris **hanya bila `CurrencyID` sama**
   (`…EstimationList(Param.Index).CurrencyID = .CurrencyID`) — secara eksplisit **mengizinkan** baris
   bermata uang berbeda hidup berdampingan.
3. `Activity/CountEstimation_Act.xml` mengagregasi ke `Money.pxResults` per `CurrencyID` (1239–2014)
   lalu **bercabang eksplisit**: `Local.SizeMoney < 2` (baris 2790) versus `> 1` (baris 3959) —
   jalur "satu mata uang" versus "multi mata uang". Hasilnya ke `ListTotalEstimation`.

**Konsekuensi:** bentuk uang di Claim Prop wajib `(amount, currency, kurs)` **per baris**, dan
subtotal per mata uang adalah entitas tersendiri. Menyalin invariant Life ke sini akan salah.

### 6d. Gerbang tanggal `[terverifikasi]`

`Activity/CheckEstimateDate_Act.xml` menegakkan **DOL ≤ EstimationDate ≤ hari ini**, inklusif di
kedua ujung, zona `Asia/Jakarta`:

- baris 552: `@defaultCompareDates(pyWorkPage.ClaimData.DateOfLoss, .EstimationDate)`; pesan
  `"Estimation Date should not be less than Date of Loss"` (354) dipasang bila `ResultCompare=="1"` (768)
- baris 899: `@CompareDates(Param.CurrentDate, Param.EstimateDate)`; pesan
  `"Estimation Date should not be more than todays date"` (375) bila `ResultCompare=="false"` (1435),
  dengan koreksi "sama dengan hari ini → lolos" (1035/1285)

Pesan dipasang **per baris** (`Property-Set-Messages` pada `.EstimationDate` di class `Data-Estimasi`),
bukan di header.

### 6e. Dua plafon lain yang ditegakkan Estimation `[terverifikasi]`

| Plafon | Kondisi | Bukti |
| --- | --- | --- |
| TSI / interest insured | `Local.TotalEstimasiIDR > Local.TotalInsured` | `CountEstimation_Act.xml` 571 (pesan `"Value Estimation RNM more than TSI Interest Insured"`), dipasang di 5860–5918 |
| Cash call | `Local.TotalEstimasiIDR > Local.CashCall` | `DeleteEstimation_Act.xml` 2462; `CashCall` difilter `.TreatyGroupID` (2182) dan `.TreatyType` (2235) |

### 6f. Estimation ↔ Adjustment — **agregat, bukan FK**

`[terverifikasi]` **Tidak ada properti penghubung** dari satu baris Estimation ke satu baris
Adjustment. Yang menyeberang hanya **rollup**:

- `AddAdjustment_Act.xml` 1153: `AdjustmentList(<LAST>).TotalEstimasiValue ← .ClaimData.TotalEstimasiIDR`
- dipakai sebagai **plafon validasi** di `CountValueADJTreaty_Act.xml`: `.AdjustmentValue >
  .TotalEstimasiValue` (3431), `.ValueAdjustment > .TotalEstimasiValue` (3624),
  `Local.SumProposeAdjustment > .TotalEstimasiValue && .Type=="1"` (3647); pesan
  `"Adjustment RNM should not be greater than estimation RNM"` (493) dan varian total (577)

⚠️ `[terbuka]` **`EstimationValue` pada baris adjustment yang menyeberang ke Komite bukan berasal dari
Estimation.** Di `AddKomiteTreatyChild_ACT.xml` 3369 ia disalin dari `.EstimationValue` pada class
**`ASM-FW-GCNMFW-Data-Adjustment`** — properti bernama sama di class berbeda. Pencarian
`AdjustmentList(…).EstimationValue` di seluruh korpus: **nol hasil**. **Sumber pengisinya tidak
ditemukan di korpus.**

### 6g. Persistensi — tidak ada rule Estimation yang menulis ke Oracle

`[terverifikasi]` **Nol `RDB-Save` di seluruh korpus.** Kelima activity Estimation hanya melakukan
`RDB-List` (lookup kurs, read-only). Baris Estimation hidup di clipboard sampai klaim disimpan.

`[terverifikasi]` Penulis sesungguhnya `Activity/SaveOutstanding_Act.xml`: loop `EstimationList`
(4144), memetakan tiap baris ke `TempOSAkseptasi.*` (4390–4735) dan paralel ke
`pyWorkPage.AkseptasiList(<APPEND>)` (4755–4984), lalu `RDB-List` (5612) ke:

| Rule SQL (`Claim Prop/RDBList/`) | Procedure Oracle | Pemanggil |
| --- | --- | --- |
| `SaveOSClaim_SQL.xml` (`…!ASM!SAVEOSCLAIM_SQL`) | `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` | `SaveOutstanding_Act` 5612, 11338, 11818 |
| `SaveOSKlaimTreaty_SQL.xml` (`…!GCNM!SAVEOSKLAIMTREATY_SQL`) | `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT` | `SetOutstanding_Act` 832/1303; `UpdateTableOS` 1410/2303 |
| `SaveDataToOsAkseptasiNP.xml` (`…!GCNM!SAVEDATATOOSAKSEPTASINP`) | `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` | `CloseClaimProp` 1836, 3006 |

⚠️ Ini menjawab sebagian akhiran yang selama ini `[terbuka]`: `TRT` dipakai jalur **Treaty**
(`SetOutstanding`/`UpdateTableOS`), `TNP` dipakai jalur **Close**. Kepanjangannya tetap
**tidak terbaca dari korpus** → `[terbuka]`.

### 6h. ⚠️ Nama Section berbohong

`[terverifikasi]` `Section/InputAcceptation_Est.xml` dan `Section/OutstandingClaim_Est.xml` adalah
**Save-As dari section `_Intrs` (Interest)** yang di-rename: `<pzOriginalInstanceKey>` keduanya
berbunyi `…INPUTACCEPTATION_INTRS` dan `…OUTSTANDINGCLAIM_INTRS` (baris 74 dan 75). Keduanya
**section klaim penuh multi-grid**, bukan layar khusus Estimation — Estimation hanya satu dari 5–6
repeat grid di dalamnya. Menguatkan §3: `pzOriginalInstanceKey` menyimpan silsilah, bukan identitas.

`[terverifikasi]` Kedua section **tidak identik**: `.ClaimData.SpreadingClaim` (36504) dan
`.ClaimData.SpreadingBreakQS` (40479) hanya ada di `InputAcceptation_Est.xml`.

---

## 7. Bentuk penyimpanan — **tidak ada skema relasional untuk dimigrasikan**

### 7a. ⚠️⚠️ Temuan terpenting: data klaim inti TIDAK pernah lewat Connect-SQL

`[terverifikasi]` `AdjustmentList` dan `EstimationList` hidup sebagai **page list di dalam work
object Pega** berkelas `ASM-FW-GCNMFW-Work-ClaimTreaty`, dan dipersist lewat **`Obj-Save` ke tabel
work Pega (blob)** — bukan lewat Connect-SQL.

Bukti negatif (manipulasi page murni, nol penulisan DB):

| Rule | Bukti |
| --- | --- |
| `Activity/AddAdjustment_Act.xml` | `Property-Set` step page `.ClaimData.AdjustmentList` (2340–2347); `Property-Remove` (2433) |
| `Activity/DeleteAjsutment_Act.xml` | hanya `Property-Remove` (332); **nol `RDB-List`** di seluruh berkas |
| `Activity/AddEstimation_Act.xml` | step page `…EstimationList` (1004); satu-satunya `RDB-List` (2128) adalah lookup kurs `CurrencyStandard` |
| `Activity/DeleteEstimation_Act.xml` | `Property-Remove` (621, 789); tanpa `RDB-List` |

Bukti positif `Obj-Save` atas `pyWorkPage`: `SaveOutstanding_Act` 3884 & 9405 ·
`SetValueToClaim_Act` 3297 · `SetPayableTreaty_Act` 4193 · `PrintDLATreatyIn` 7185 ·
`AddKomiteTreatyChild_ACT` 7317 · `CheckNoPolicy` 1634 · `SendCloseClaimToKomite` 2403.

**Konsekuensi mengikat untuk `/to-spec`:** skema relasional Claim Prop **tidak dapat direkayasa-balik
dari korpus** — ia belum pernah ada. Sumber kebenaran sekarang adalah blob work object. Skema Oracle
untuk klaim, adjustment, estimation, spreading, interest, dan loss allocation **harus dirancang
baru**, dan itu menjadikannya **PREFACTOR** — pola yang sama seperti tiket 14 Claim Life dan tiket 00
PremiumList Life.

### 7b. Yang benar-benar ditulis ke Oracle hanyalah **proyeksi**

`[terverifikasi]` Dua tabel, keduanya **hibrida**: satu kolom dokumen `DATA_JSON` + beberapa kolom
kunci relasional.

| Peran | Tabel | Kolom relasional | Atribut `DATA_JSON` yang terbaca |
| --- | --- | --- | --- |
| Dokumen klaim | **`JSON_KLAIM`** | `NOPOLIS`, `IDPEGA`, `TGL_INPUT` | `DateOfLoss`, `ClaimNo`, `CauseOfLoss`, `ClaimEstimate` (`RDBList/CariHistoryClaim_SQL.xml` 78–84) |
| Baris outstanding / akseptasi | **`OS_AKSEPTASI_KLAIM`** | `NOCLAIM`, `NOPOLIS`, `TANGGAL` | `Currency`, `CurrencyID`, `GrossValue`, `Value`, `TotalGross`, `KursValue`, `PersenRNM`, `AcceptedNo`, `Type` (`RDBList/DataOutstandingTreatyin.xml` 60–71; `GetDataKlaimTreatyin.xml` 60–69) |
| Log | `pooldata.monitoring_klaim_log`, `POOLDATA.DIRECTTOKASIR_LOG` | — | `InsertLogServiceClaim.xml` 80; `InsertLOGDirectKasir_SQL.xml` 86 |
| Core eksternal | `reinsurance.trloss_detail_t` | **hanya dibaca** (`getStatusKonversi_SQL.xml` 79) | penulisannya lewat Connect-REST `KonversiKlaimNonLife` |

⚠️ Akses `a.DATA_JSON.NamaAtribut` adalah **dot-notation JSON Oracle** — jadi `DATA_JSON` memang
kolom dokumen, bukan CLOB buram. Nama `JSON_POLIS` juga **jujur**: `GetYearofQuartal.xml` 79
mengakses `a.DATA_JSON.YearOfQuartal`.

### 7c. ⚠️ `InsertJsonClaimTreaty_act` — di sini namanya **JUJUR** (kebalikan kasus Life)

`[terverifikasi]` `ASM-FW-GCNMFW-WORK-KOMITETREATY` / `INSERTJSONCLAIMTREATY_ACT` /
`RULE-OBJ-ACTIVITY` — perhatikan class-nya **KomiteTreaty**, bukan ClaimTreaty. Hanya **4 step**:

1. `Call GetBase64Attachment` (284)
2. `Property-Set` → `TempPNC.BUSINESS_CODE ← TempOpenPage.pzInsKey` (407), `.POLICY_NO` (454), `.No_Klaim` (475)
3. `Property-Set` → `TempPNC.TSI ← @GCNM.GetPageJSONString()` (583)
4. `RDB-List` → `InsertClaimPNC` (691)

⚠️ **Dua nama properti berbohong berat:** `BUSINESS_CODE` sebenarnya berisi **case key Pega**
(`pzInsKey`), dan **`TSI` berisi seluruh blob JSON** — bukan sum insured.

⚠️ `[terverifikasi]` Library-nya **`GCNM`**, bukan `ASM`: `@GCNM.GetPageJSONString()`. Varian
`@ASM.GetPageJSONString()` juga dipakai di modul ini tetapi di activity lain (`SaveOutstanding_Act`
5094, `CloseClaimProp` 1677). **Dua library berbeda untuk fungsi bernama sama** — definisi keduanya
tidak ada di korpus → `[terbuka]`.

SQL yang dijalankan (`RDBList/InsertClaimPNC.xml` 85–90):
`BEGIN POOLDATA.PEGA_JSON_KLAIM_PNC({No_Klaim},{POLICY_NO},{BUSINESS_CODE},{TSI},{BRANCH_NAME out},{BRANCH_CODE out}); COMMIT; END;`

### 7d. Dualitas JSON versi Claim Prop — **dua tabel, selalu berpasangan**

`[terverifikasi]` Bukan dualitas JSON-vs-relasional seperti di Life, melainkan **dua penulisan JSON
ke dua tabel berbeda**, dipanggil berpasangan dalam activity yang sama:

| Activity | Tulis `OS_AKSEPTASI_KLAIM` | lalu tulis `JSON_KLAIM` |
| --- | ---: | ---: |
| `SetOutstanding_Act.xml` | 774 (`_KLAIMTRT`) | 993 |
| `UpdateTableOS.xml` | 1353 (`_KLAIMTRT`) | 1572 |
| `SaveOutstanding_Act.xml` | 5556 (`_KLAIM`) | 9316 |
| `CloseClaimProp.xml` | 1782 (`_KLAIMTNP`) | 2247 |

### 7e. Tiga procedure `PEGA_JSON_OS_AKSEP_*` — payload **hibrida**

`[terverifikasi]` Ketiganya menerima **satu parameter JSON besar + beberapa kunci relasional**:

| Param | Isi `[terverifikasi]` |
| --- | --- |
| `InputData.CARI1` | **string JSON** — `@ASM.GetPageJSONString()` (`SaveOutstanding_Act` 5093; `CloseClaimProp` 1676) |
| `InputData.CARI2` | `pyWorkPage.pzInsKey` |
| `InputData.CARI3` | `…PolicyData.PolicyNo` |
| `InputData.CARI10` | `TempOSAkseptasi.Type` (di `CloseClaimProp` literal `"4"`) |
| `InputData.CARI18` | `ClaimData.IDMaster` / `TreatyInMaster.ID` |
| `InputData.CARI29` | `ClaimData.NoClaim` |

⚠️ `[terverifikasi]` **Parameter yang dikirim tetapi TIDAK PERNAH di-set** di seluruh `Claim Prop`:
`InputData.CARI16`, `CARI17`, `CARI20`, `InputParamOs.Currency`, `InputParamOs.TypeLoss`. Nilainya
yang sampai ke procedure **tidak terbaca dari korpus** → `[terbuka]`.

⚠️ `[terverifikasi]` `SetOutstanding_Act` dan `UpdateTableOS` memanggil `_KLAIMTRT` **tanpa men-set
`CARI1` (payload JSON) di activity itu sendiri** — mereka bergantung pada sisa state page `InputData`
di thread. Rapuh; wajib dimodelkan ulang secara eksplisit.

`[terverifikasi]` Pola pemanggilan lintas modul: `_KLAIMTRT` **hanya** di `Claim Prop` (konsisten
dengan `TRT` = Treaty); `_KLAIMTNP` muncul di `Claim Prop`, `Claim Non Prop`, `Komite Claim Prop`,
`Komite Claim Non Prop` — jadi **`TNP` BUKAN sekadar "Non Prop"**. Kepanjangannya **tidak terbaca
dari korpus** → `[terbuka]`.

### 7f. ⚠️ Titik potong transaksi — **13 rule COMMIT sendiri**

`[terverifikasi]` `COMMIT;` di dalam blok PL/SQL: `SaveOSClaim_SQL` (82), `SaveOSKlaimTreaty_SQL`
(88), `SaveDataToOsAkseptasiNP` (88), `InsertClaimPNC` (88), `InsertLogServiceClaim` (82),
`InsertLOGDirectKasir_SQL` (93), `Insert_T_Storage_SQL` (102), `UpdateMCauseOfLoss` (69),
`UpdateDCauseOfLoss` (94), `GenerateNoCLMTreatyIn` (98), `GenerateNoTRTInTemp` (70),
`GetSequenceNumber_SQL` (88), `GetTokenStorage_SQL` (87).

`[terverifikasi]` **Tanpa** COMMIT: `Update_T_Storage_SQL` (85–90), `DeleteDataTreatyGroup_SQL` (58).

`[terverifikasi]` Satu-satunya step method `Commit` di seluruh modul:
`Activity/GetPayAttachmentAdj_Act.xml` 2296. `Obj-Save` ber-`WriteNow=true` di 7 tempat
(`AddKomiteTreatyChild_ACT` 7369, `CheckNoPolicy` 1686, `GetNameCauseofLoss_Act` 414,
`GetPayAttachmentAdj_Act` 1540 & 2205, `HitServiceToKasir_Act` 9079, `InsertDocument_Act` 1232).

⚠️ **Perilaku existing TIDAK atomik.** Satu aksi pengguna di `SetOutstanding_Act` menghasilkan
**≥ 2 commit Oracle terpisah + 1 simpan work object**: step 4 commit, step 5 commit, lalu `Obj-Save`.
Bila step 5 gagal, step 4 sudah ter-commit.

**Peluang:** berbeda dari Claim Life (yang cut point-nya hilang begitu prosedur JSON dibuang), di
sini Go **bisa** membuatnya atomik — syaratnya `COMMIT;` di 13 rule itu tidak direplikasi dan seluruh
penulisan dipindah ke satu transaksi.

### 7g. Nama yang berbohong — daftar tambahan `[terverifikasi]`

| Rule | Kebohongan |
| --- | --- |
| Folder `RDBList/` | **nol** rule bertipe `RULE-RDB-*`; seluruh 54 berkas adalah `RULE-CONNECT-SQL` |
| `getStatusKonversi_SQL.xml` | class-nya `…INT-OS_AKSEPTASI_KLAIM`, SQL-nya `select COUNT(1) from reinsurance.trloss_detail_t` — **tidak menyentuh** tabel itu |
| `GetIDConsultanAdj_SQL.xml` | silsilah `…ASM!SAVEMASTERADJUSTERCONSULTANT_SQL`, isinya SELECT murni (`ADJUSTERCONSULTANT_SEQ.nextval`) |
| `Update_T_Storage_SQL`, `GetTokenStorage_SQL` | silsilah keduanya `…RNM!INSERT_T_STORAGE_SQL` |
| `Insert_T_Storage_SQL` | silsilah `…RNM!GENERATEIMAGEID_SQL` |
| `T_STORAGE_IMAGE` | kolom `TANGGAL_UPLOAD` ada di UPDATE (85–90) tetapi **tidak** di INSERT (85–103) |

---

## 8. Loss Allocation — nama properti aslinya `SpreadingRisk`

⚠️ `[terverifikasi]` **Nama berbohong lagi.** Di level klaim konsep ini **tidak** bernama
`LossAllocation`; ia bernama **`.ClaimData.SpreadingRisk`**. Nama `LossAllocation` hanya dipakai
untuk **salinannya di dalam baris Adjustment**.

| Induk | Properti daftar | Class elemen | Bukti |
| --- | --- | --- | --- |
| `Data-ClaimData` | `.ClaimData.SpreadingRisk` | `ASM-FW-GISFW-Data-SpreadingRisk` | `Section/InputAcceptation_Est.xml` 21975/21982 |
| `Data-Adjustment` | `.LossAllocation` | `ASM-FW-GISFW-Data-SpreadingRisk` | `Section/AdjustmentDetail.xml` 15620/15630 |

Bahwa grid `SpreadingRisk` itulah yang berlabel "Loss Allocation": `<pyTitle>Loss Allocation` di
`InputAcceptation_Est.xml` 21460 dan `OutstandingClaim_Est.xml` 20071.

`[terverifikasi]` **Induk = header klaim**; salinan ke baris adjustment adalah **snapshot**:
`AddAdjustment_Act.xml` 1069 → `AdjustmentList(<LAST>).LossAllocation ← ClaimData.SpreadingRisk`.
Treaty **bukan** induk — `TreatyName`/`TreatyType` hanya field datar di dalam baris.

⚠️ `[terverifikasi]` Class `ASM-FW-GISFW-Data-SpreadingRisk` dipakai bersama **8 page list berbeda**
(`SpreadingRisk`, `SpreadingClaim`, `SpreadingBreakQS`, `SpreadingAdjustment`,
`SpreadingAdjustmentQS`, `SpreadingQuotaShare`, `LossAllocation`) — satu class, banyak peran.

### 8a. Kolom dan pembagian nilai

| Kolom | Isi `[terverifikasi]` |
| --- | --- |
| `TreatyName`, `TreatyType` | identitas treaty (`AddLossAllocation_act.xml` 1279–1940) |
| `SharePercentage` | porsi pembagian |
| `CurrencyID`, `Currency` | dari `ListClaimAmount(<LAST>)` |
| `PremiumSpreaded` | ⚠️ **berisi KURS, bukan premi** — diisi dari `ListClaimAmount.IDR` yang sendirinya `← .KursObjectItem` (`CountTotalInsterest_Act.xml` 3788 → 4000), lalu dipakai sebagai pengali konversi IDR (4769) |
| `ClaimSpreaded`, `ClaimEstimation` | turunan |
| `IsOldData` | filter |

Perhitungan `[terverifikasi]` (`CountTotalInsterest_Act.xml`, loop `SpreadingRisk` 4663):

```
.ClaimSpreaded   = .SharePercentage * Local.Value / 100          (4723)
.ClaimEstimation = .ClaimSpreaded * .PremiumSpreaded             (4769)
precond: .CurrencyID == Local.CurencyID && .IsOldData != "Yes"   (4820)
```

`[terverifikasi]` Validasi **total share ≤ 100 % per mata uang** — `CountPersen_act.xml` 1174 akumulasi,
1393 `@greaterThan(Local.Total,100)`, pesan `"Total more than 100%"` (294), dipasang di 1433–1448.

`[terverifikasi]` Persentase loss allocation juga mengalikan nilai adjustment:
`CountGrossAdjTreaty_Act.xml` 810 dan `CountValueADJTreaty_Act.xml` 2574/2627 —
`… × PersenRNM/100 × PersenLossAllocation/100 × ShareCeding/100`.

**Dimensi pembagian: TREATY (`TreatyName` + `TreatyType`), disegmentasi per MATA UANG.**
Tidak ada dimensi tahun underwriting maupun coverage. Ini **berbeda** dari spreading Claim Life
(yang berdimensi treaty-year → reinsurer).

### 8b. ⚠️ `FilterLossAllocation_Act` — arah gerbangnya terbalik dari pola korpus

`[terverifikasi]` Kondisi di `FilterLossAllocation_Act.xml` 564:
`local.TreatyName == .TreatyName && local.SharePersen == .SharePercentage`, dengan
`<pyStepsPreCondParamsWhenTrue>3` dan `<pyStepsPreCondParamsWhenFalse>2` (558/560).

`[terverifikasi]` Pola standar di seluruh korpus adalah `WhenTrue=2 / WhenFalse=3`, dan step tanpa
kondisi selalu `2/2` (1.275 kemunculan berbanding 77).

⚠️ `[dugaan — enum Pega tidak didefinisikan di korpus]` Bila `2 = jalankan` dan `3 = lewati`, maka
efeknya: baris yang **cocok justru DIPERTAHANKAN**, yang **tidak cocok DIHAPUS** — jadi ini
"filter/keep", bukan "dedup/remove". **Arti enum ini tidak dapat diverifikasi dari korpus** dan
harus dikonfirmasi sebelum perilaku ini direplikasi.

### 8c. Persistensi

`[terverifikasi]` **Tidak ditemukan penulis Loss Allocation ke Oracle.** Konsisten dengan §7a —
hidup di work object, disimpan lewat `Obj-Save`.

---

## 9. ⚠️⚠️ "Interest" di Claim Prop **BUKAN bunga** — ia Insured Interest (objek pertanggungan)

`[terverifikasi]` Salah satu kebohongan nama terbesar di modul ini. Bukti bebas:

- `Section/OutstandingClaim_Intrs.xml` 626: `<pyTitle>Insured Interests 100 %`
- 1315: caption kolom `Insured Interest`; 1602: `Value In IDR`
- `AddInterest_act.xml` 294: `DataChronology.CARI1 = "Add Insured Insterest"`
- `DeleteInterest_act.xml` 365: `"Delete  Insured Interest"`

| Properti daftar | Class elemen | Bukti |
| --- | --- | --- |
| `.ClaimData.InterestList` | `ASM-FW-GCNMFW-Data-ObjectItem` | `OutstandingClaim_Intrs.xml` 1140/1149 |
| `.ClaimData.InterestListDtl` | `ASM-FW-GCNMFW-Data-ObjectItem` | `ViewDetailInterest.xml` 1118/1127 — salinan untuk layar detail (`DataTransform/CekInterestListDtl_DT.xml` 182) |
| `.ClaimData.TotalInterestInsured` | `ASM-FW-GISFW-Data-TreatyInTotal` | `OutstandingClaim_Intrs.xml` 5381/5390 — agregat **per mata uang** |

Kolom satu baris: `ObjectName`, `CurrencyID`, `Currency`, **`KursObjectItem`** (kurs per baris),
`TSIPerObject`, `TSIPerObjectIDR`, `IsAdjVal`. Induk = **header klaim**.

### 9a. `CountTotalInsterest_Act` — tidak ada rumus bunga sama sekali

`[terverifikasi]` **Nol** kemunculan rate, jumlah hari, atau basis 360/365 di seluruh 9.117 baris.
Yang ada: agregasi TSI per mata uang + konversi kurs + **re-kalkulasi turunan** (loss allocation,
estimation, spreading). Inti perhitungannya:

```
.TSIPerObjectIDR = @toDecimal(.KursObjectItem) * @toDecimal(.TSIPerObject)   (2216)
Local.TotalInIDR = Local.TotalInIDR + .TSIPerObjectIDR                        (2262)  precond .IsAdjVal==""
pyWorkPage.ClaimData.TotalSumInsuredIDR = Local.TotalInIDR                    (3405)
```

⚠️ `[terverifikasi]` **Blok rumus Estimation muncul DUA KALI secara literal** (5559–5726 dan
5841–6008) dengan precondition identik — salin-tempel. Mana yang menang **tidak terbaca**.

### 9b. Menguatkan §6c — mata uang per baris

`[terverifikasi]` `SetCurencyInterest_act.xml` menulis `CurrencyID` (267), `Currency` (735), dan
`KursObjectItem` (1085) **per `Param.Index`**; agregat `TotalInterestInsured` disusun per mata uang
dengan dedup `CurrencyID#Currency` (2539–2560). Multi-mata-uang bukan hanya di Estimation — ia pola
modul.

⚠️ `[terverifikasi]` `SetCurencyInterest_act` adalah **Save-As dari `CURENCYESTIMATION_ACT`**
(`<pzOriginalInstanceKey>` baris 60).

### 9c. Kapan dihitung

`[terverifikasi]` **On-demand di UI, bukan saat simpan maupun akseptasi.** Hanya dua pemicu:
event `change` pada `.TSIPerObject` (`OutstandingClaim_Intrs.xml` 3036/3086/3200) dan
`Call CountTotalInsterest_Act` dari `SetCurencyInterest_act` 1235. **Nol** pemanggilan dari
`SaveOutstanding_Act` atau `SaveAcceptationTreaty_Act`.

`[terverifikasi]` Tidak ditemukan penulis Interest ke Oracle.

---

## 10. Deductible — dan ⚠️⚠️ **kontradiksi urutan yang memblokir**

### 10a. Rumus `[terverifikasi]`

`Activity/CountDeductible_Act.xml`, seluruhnya bergerbang `.ClaimData.FormType==2` (1188).

```
cabang TypeDeductible==1  (basis = TSI dari TotalInterestInsured, dicocokkan CurrencyDeductible)
  local.DeductibleClaim = @divide(ClaimData.Amount * Local.UangNilai, 100, 8)          (690)
  ClaimData.NetDeductibleValue = @if(local.DeductibleClaim > ClaimData.DeductibleValue,
                                     local.DeductibleClaim, ClaimData.DeductibleValue) (742)

cabang TypeDeductible==2  (basis = field TSIDeductible)
  Local.DeductibleTSI = @divide(ClaimData.Amount * ClaimData.TSIDeductible, 100, 8)    (969)
  ClaimData.NetDeductibleValue = @if(Local.DeductibleTSI > ClaimData.DeductibleValue,
                                     Local.DeductibleTSI, ClaimData.DeductibleValue)   (1021)
```

Jadi **deductible = MAX(persentase × basis, nilai flat)**.

⚠️ `[terverifikasi]` `ClaimData.Amount` adalah **persentase, bukan uang** — label `%` di
`InputAcceptation_Est.xml` 9740, tepat sebelum field di 9769.

⚠️ `[terverifikasi]` **Cache expression-builder di berkas yang sama BASI dan tidak sinkron** dengan
rumus aktif: baris 703/706 dan 982/985 menyimpan teks berbeda dari `<PropertiesValue>` aktif di
691/970. **Yang otoritatif adalah `<PropertiesValue>`** — `<pyExpression>` dan
`<rowdata REPEATINGINDEX="PropertiesValue">` adalah cache UI. Jebakan salin-kutip.

`[terverifikasi]` Deductible melekat pada **header klaim** (properti skalar `Data-ClaimData`), lalu
disalin ke tiap baris adjustment saat dibuat (`AddAdjustment_Act.xml` 1325/1372/1393, tiga cabang di
1445, 1632, 1819). Dipicu **hanya** dari event `change` di UI pada 5 field — **nol** pemanggilan dari
alur simpan/akseptasi.

### 10b. ⚠️⚠️ PEMBLOKIR — dua rumus berlawanan untuk nilai yang sama

`[terverifikasi]` Dua activity menulis `ListClaimAmount.Value` dengan urutan **berbeda** terhadap
`ShareCeding`:

| Rule | Rumus | Dibuat |
| --- | --- | --- |
| `Activity/AddListClaimAmount.xml` 2006/2052 | `Value = Share% × (ClaimAmount − Deductible)` | 2022-06-15 |
| `Activity/CountListClaimAmountIDR.xml` 338/386 | `Value = (Share% × ClaimAmount) − Deductible` | 2024-04-19 |

⚠️ Lebih buruk: properti `.NetDeductibleValue` **berarti berbeda** di keduanya — di (a) ia
"klaim setelah dikurangi deductible", di (b) ia "nilai deductible mentah". **Nama sama, semantik
berlawanan.**

`[terverifikasi]` Keduanya dipicu event UI berbeda (`AddListClaimAmount` dari tombol click;
`CountListClaimAmountIDR` dari `change` pada `.ClaimAmount`) dan **keduanya menulis field yang sama**
— sehingga **hasil akhir bergantung pada urutan interaksi pengguna**. Korpus **tidak** menyediakan
urutan deterministik.

**Ini memblokir bentuk perhitungan nilai klaim.** Setara dengan kasus `PREMIUM_SPREADED_NET` di
Claim Life: dua cabang rumus untuk satu nilai, dan mana yang berlaku tidak terbaca. **Jangan tebak.**

Konsekuensi berantai: `SpreadingRisk.ClaimSpreaded` dihitung dari `ListClaimAmount.Value`
(`CountTotalInsterest_Act.xml` 4578 → 4724), jadi seluruh loss allocation dan estimation ikut
bergantung pada rumus mana yang menang.

---

## 11. OQ-039 — ⚠️ **di Claim Prop tidak ada aksi "Reject" sama sekali**

`[terverifikasi]` Keempat activity dipicu **dari Section**, tidak satu pun dari flow:

| Activity | Pemanggil | Pemicu |
| --- | --- | --- |
| `ADDADJUSTMENT_ACT` | `Section/InputAcceptation_Adjs.xml` 30137, 30284 | tombol, `pyEvent=click` |
| `DELETEAJSUTMENT_ACT` | idem 31529, 31616 | tombol, `click` |
| `CHANGEDATA_ACT` | idem 1556, 1678 | **bukan tombol** — kontrol `.IsEditEstimation`, `pyEvent=change` |
| `CLOSECLAIMPROP` | `Section/PreventRejectClaimProp.xml` 4688, 4784 | tombol "Yes" di dialog |

Dialog Close dibuka sebagai **local action** dari tombol `Close Claim`
(`Section/InputAcceptation.xml` 5598 → 5626 `<pyLocalAction>PreventRejectClaimProp`).

### 11a. ⚠️ `PreventRejectClaimProp` — namanya berbohong total

`[terverifikasi]` Ia **bukan** pencegah reject; ia **dialog konfirmasi CLOSE** dengan dua cabang:

| Cabang | Teks | Tombol "Yes" memanggil |
| --- | --- | --- |
| `AllocationShareSalvage == 'true'` (2364) | *"Are you sure want close this claim without payment?"* (2586) | `SendCloseClaimToKomite` (3077, 3156) |
| `== 'false'` (3975) | *"Are you sure want close this claim?"* (4197) | `CloseClaimProp` (4688, 4784) |
| `Komite.CARI1 != ''` (5757) | *"Success Create Request to Committee"* (6002) | — hanya tombol tutup |

⚠️ `[terverifikasi]` Properti pengendali bernama **`TempCommiteClaim.AllocationShareSalvage`** tetapi
caption-nya **`Close Without Payment ?`** (7184). **Nama properti tidak ada hubungannya dengan
fungsinya.**

`[terverifikasi]` `FlowAction/PreventRejectClaimProp.xml` hanya membungkus Section bernama sama;
`pyLocalActionActivity`, `pyValidateActivity`, `pyPreActivity` **semuanya kosong**.

`[terverifikasi]` `ReportDefinition/RejectedClaim_RD.xml` juga **bukan** soal reject klaim — ia
dipakai `CheckDateDOL_Act` sebagai **pengecualian deteksi duplikat**: bila klaim lama ada di tabel
`CLAIMREJECTED`, duplikat Date of Loss **diizinkan** (loop 4094, `Local.Flagerr=0` di 4161).
Pesan duplikatnya `"There is a similar Date of Loss to this policy number"` (833).

`[terverifikasi]` Satu-satunya "reject" lain adalah parameter keluar `STS_REJECT` ke sistem inti
non-life lewat `KonversiKlaim_Act`: `"0"` dari `SaveOutstanding_Act` 9584, `"1"` dari
`SaveAcceptationTreaty_Act` 4203, `"4"` dari `CloseClaimProp` 2043.

### 11b. ⚠️ Gerbang Close — ada, tapi bocor di banyak tempat

`[terverifikasi]` `CloseClaimProp` punya gerbang nyata (step 2–4):

```
.AcceptanceStatus=="" || =="0"                        → "Can not close claim, there is adjustment in comitee!"
.AcceptanceStatus=="1" && .DirectToKasir
   && bukan .StatusKasir=="Akseptasi Sudah Masuk ke Kasir"
                                                      → "Cannot close claim, there is a direct to
                                                         cashier that has not been successful."
Local.FlagError=="1" → WhenTrue=6 = EXIT ACTIVITY     (1048/1087/1093)
```

⚠️ Berbeda dari `SendCloseClaimToKomite` (§2d), di sini gerbangnya **benar-benar menghentikan**
(`exit activity`), bukan sekadar memasang pesan. Jadi **dua jalur close berperilaku berbeda terhadap
error yang sama.**

Tetapi `[terverifikasi]` empat kebocoran:

1. **`DeleteAjsutment_Act` nol precondition** — semua `<pyStepsPreCondParamsWhen/>` kosong. Tidak ada
   cek status, tidak ada cek komite. Baris adjustment dapat dihapus kapan saja.
2. **`ChangeData_Act` nol precondition** juga.
3. **Tombol `Close Claim` tanpa gerbang tampil.** Selnya ber-`pyVisible=ALWAYS` (5839) sehingga
   `pyCondition .ClaimData.IsCloseFile = true` (5842) **tidak aktif** — dan `IsCloseFile` muncul
   **hanya di baris itu** di seluruh korpus, tidak pernah di-set.
4. **Tidak ada gerbang "sudah pernah di-close"** — `CloseClaimProp` dapat dijalankan berulang.

### 11c. ⚠️⚠️ `.AcceptanceStatus` dibaca 26×, **tidak pernah ditulis**

`[terverifikasi]` Nol `<PropertiesName>…AcceptanceStatus</PropertiesName>` di seluruh `Claim Prop`.
Ia gerbang paling banyak dipakai di modul ini (Close, penyerahan komite, proteksi lampiran,
spreading, `CheckAnyAcceptationProp`) tetapi **sumber nilainya di luar korpus** — kemungkinan ditulis
sisi Komite. → **`[terbuka]` pemblokir**: mesin status klaim tidak dapat digambarkan utuh.

`[terverifikasi]` **Nol** rujukan ke `pyStatusWork` dan `.StatusKlaim` di seluruh modul — status kerja
Pega tidak pernah dipakai sebagai gerbang, meski `CloseClaimProp` step 10 men-set
`Resolved-Completed` lewat `Call ASMForceCaseClose`.

---

## 12. Proteksi — tiga rule, **tidak satu pun mengunci**

### 12a. `ProteksiData_act` — validasi field wajib, bukan penguncian

`[terverifikasi]` Semua yang ditulis hanya `Local.*` + `Page-Set-Messages`. **Nol properti kerja
yang ditulis.** Sebelas pemeriksaan: status master treaty, Estimation list kosong, SpreadingClaim
kosong, `EstimationList().Type`, `SpreadingClaim().SharePercentage`, tiga field Interest,
`CauseOfLoss`, `IDMaster`, `Location`, `PolicyNo`.

⚠️ `[terverifikasi]` Step 18 dan 19 (pesan `"Please Cause of loss"`) **tanpa precondition efektif**
(3586/3738 keduanya `2`/`2`) — selalu jalan.

Dipanggil dari tiga tempat, **semuanya tanpa precondition**: `FlowAction/OutstandingClaim.xml` 145
(`pyLocalActionActivity`), `CheeckNoRNM_Act.xml` 1824, dan `SaveOutstanding_Act.xml` **step 1** (441).

### 12b. `ProteksiInitialandDate_Act` — namanya berbohong

`[terverifikasi]` Isinya salin-remark dari baris adjustment sebelumnya + validasi bank/pembayaran —
bukan "initial and date". Pesan: `"Data Bank Account Can't NULL"` (707), `"Please choose Payment Type
first"` (686), `"Please choose payable first"` (665).

⚠️ `[terverifikasi]` Step 6 memanggil `Call SetProtectionEstimation` (1492) — **rule itu tidak ada di
korpus** → `[terbuka]`.

### 12c. `AttachmentProtect_ACT` — **ini gerbang yang sesungguhnya**

`[terverifikasi]` Ia menghitung `Protect.CARI1` dan `Protect.CARI2`, dan keduanya menggerbangi
tampilnya tombol penyerahan komite: `Harness/CommitteeTreaty.xml` 6941
`pyContainerVisibleWhen = Protect.CARI1 =1 && Protect.CARI2= 1`; juga
`Section/ComiteeClaimTreaty.xml` 5958 dan `Section/AdjustmentDetail.xml` 34901/34914.

⚠️ **Dokumen wajib ditentukan oleh `PaymentType`** `[terverifikasi]`:

| `PaymentType` | Dokumen wajib | Deskripsi step di korpus |
| --- | --- | --- |
| `1` | LOD + DLA + SPGR (+ **ADU** bila lewat limit) | *"Set Protek agar tidak send komite **Claim**"* (4581) |
| `2` atau `4` | Invoice + DLA | *"… **Expense**"* (4761) |
| `3` | Salvage + DLA | *"… **Salvage**"* (4941) |

⚠️ Ini **bukti korpus terkuat sejauh ini** tentang arti `PaymentType` — deskripsi step-nya sendiri
menyebut **Claim / Expense / Salvage**. Tetapi itu komentar pengembang, **bukan** definisi bisnis;
arti resmi tetap **OQ-020** dan **jangan ditetapkan dari sini**.

Gerbang tambahan: `Protect.CARI1` juga menuntut data bank lengkap (5132) dan spreading terisi (5273);
`Protect.CARI2=1` (5571) menuntut `.IsKomite==""` **dan** `ParamData.HASIL1 != "BELUM LUNAS"`.

⚠️ `[terverifikasi]` Page `AttachCategory` — sumber seluruh flag kelengkapan dokumen —
**tidak pernah diisi di dalam activity ini**; hanya dideklarasikan di Pages & Classes. Sumbernya
**di luar korpus** → `[terbuka]`.

### 12d. ⚠️ `CekProteksiKlaim` — semantiknya **kebalikan** dari mengunci

`[terverifikasi]` SQL (`RDBList/CekProteksiKlaim.xml` 79):

```sql
select policy_no as "PolicyNo" from pooldata.openproteksi_edm
where type = '5' and STS_AKSEP = '1' and policy_no = {ParamData.CARI3}
```

`[terverifikasi]` Di `CekPremiLunas_Act`, blokir default adalah
`"Akseptasi tidak dapat dilanjutkan dikarenakan Premi belum Lunas"` (1094) bila sisa premi > 0.
Step 7.2 lalu men-set `ParamData.HASIL1="Lunas"` dengan deskripsi **"set lunas jika proteksi dibuka"**
(1471) bila ada baris di `openproteksi_edm`. **Jadi "proteksi" di sini adalah WAIVER** — adanya baris
justru **membatalkan** blokir premi belum lunas.

### 12e. Bukan lock Pega, dan **tidak ada pemeriksaan peran sama sekali**

`[terverifikasi]` Nol `Obj-Refresh-And-Lock` di keempat rule proteksi. Lock Pega dipakai di rule
lain (`SaveOutstanding_Act` 4016, `GetPayAttachmentAdj_Act` 2795, `SetPayableTreaty_Act` 4325), dan
Pega sendiri memperingatkan hasil lock **tidak dicek** (`SaveOutstanding_Act` 12430).

⚠️⚠️ `[terverifikasi]` **Tidak ditemukan satu pun pemeriksaan peran, privilege, atau access group**
di seluruh `Claim Prop`: `<pyWhenNotPrivilege/>` kosong; nol kondisi visible-when berbasis
operator/role. Yang ada hanya **identitas orang yang ditanam sebagai literal** (§13b). Ini
menegaskan **ADR-0014** — wewenang harus ditegakkan di lapisan layanan, karena di sistem lama ia
praktis tidak ditegakkan sama sekali.

---

## 13. OQ-040 — limit wewenang: jawabannya **bukan seperti yang tertulis di register**

### 13a. Ketiga rule "limit" ternyata bukan tangga wewenang

| Rule | SQL `[terverifikasi]` | Perannya sebenarnya |
| --- | --- | --- |
| `GetLimitDirekturUtama_SQL` | `select LIMIT_BOTTOM from POOLDATA.EMAILKOMITE where name = '<satu nama orang>'` (60) | ambang tunggal untuk **mewajibkan dokumen ADU** — bukan wewenang menyetujui |
| `GetLimitPLATreatyin` | `SELECT … FROM PROPORTIONALARRG WHERE treatydescid = '10001' and …` (60) | ambil share/PCT-RP-USD treaty; `treatydescid` **hardcode** |
| `GetLimitsTreatyIn_SQL` | `select a.JSONDATA from m_treaty_in a where id={…} union all … m_treaty_in_edm` (86–90) | ⚠️ **bukan limit sama sekali** — ini pengambil JSON master treaty. Silsilahnya `…ASM!GETCOPYNBFORCLAIM` |

⚠️⚠️ `[terverifikasi]` `GetLimitDirekturUtama_SQL` memilih barisnya dengan **nama orang sebagai
literal di `WHERE`** — bukan jabatan, bukan `DEGREE`, bukan `JABATAN`, padahal kedua kolom itu ada di
class yang sama. Satu-satunya petunjuk bahwa ini Direktur Utama adalah **komentar** pada step
pemanggil (`AttachmentProtect_ACT` 2877). **Tidak dapat diverifikasi dari korpus** → `[terbuka]`.

`[terverifikasi]` Akibat melewati limit: `local.LewatLimit = 1` (2888) bila
`local.TotalAdjustment >= LIMIT.pxResults(1).LIMIT_BOTTOM` (2967) → **dokumen ADU menjadi wajib**
untuk `PaymentType==1`.

### 13b. Tangga wewenang yang sesungguhnya — lewat roster, dan `LIMIT_TOP` tidak dipakai

`[terverifikasi]` `FilterEmailKomiteWithLimit` difilter `A AND C AND B` (662):
`.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM` · `.STS_KLAIM = Param.STS_KLAIM` · `.STS_AKTIF = "1"`.
Dipanggil dua tempat, keduanya dengan `Param.LIMIT_BOTTOM = Local.TotalAdjustment` dan
`Param.STS_KLAIM = "PROP"` (hardcode): `SetKomiteTreaty_ACT` 630/697/718 dan
`AddKomiteTreatyChild_ACT` 4971/5038/5060.

⚠️ `[terverifikasi]` Arah perbandingannya `<=`: **makin besar nilai adjustment, makin BANYAK anggota
komite** yang masuk daftar. Kolom **`LIMIT_TOP` ada di class tetapi tidak pernah dipakai sebagai
filter** — jadi pita atas tidak ditegakkan. Ini konsisten dengan Fakta 2 (yang boleh diwarisi dari
Life).

### 13c. Nilai hardcode yang relevan

`[terverifikasi]` **Tidak ada ambang nominal (rupiah/USD) yang di-hardcode di `Claim Prop`** —
berbeda dari Claim Non Prop. Yang hardcode justru **kunci dan identitas**:

| Nilai | Lokasi | Peran |
| --- | --- | --- |
| satu **nama orang** | `GetLimitDirekturUtama_SQL` 60 | penentu baris limit |
| `"PROP"` | `SetKomiteTreaty_ACT` 719; `AddKomiteTreatyChild_ACT` 5060 | filter roster |
| `"1"` (`STS_AKTIF`) | `FilterEmailKomiteWithLimit` 705 | filter roster |
| `'10001'` (`treatydescid`) | `GetLimitPLATreatyin` 60 | filter treaty |
| `'5'`, `'1'` | `CekProteksiKlaim` 79 | penanda waiver |
| satu **operator ID + email + jabatan** | `SendCloseClaimToKomite` 1103/1169/1189 | approver tunggal jalur close — **melewati mekanisme limit sepenuhnya** |
| tiga **nama orang** | `SethistoryKlaimTreaty` 538/680/822 | pemetaan nama→jabatan di riwayat |
| satu **operator ID** | `SendEmailKlaimRejectClose` 693/884 | cabang email khusus satu orang |

⚠️ Nilai PII di atas **sengaja tidak disalin**. Seluruhnya **wajib** menjadi data master di sistem
baru. Ini menyentuh **OQ-021** — dan melengkapi §5: guard identitas literal memang nihil, tetapi
**identitas sebagai kunci data** ada di tujuh tempat.

---

## 14. Spreading di Prop — **bentuknya berbeda dari Claim Life**

`[terverifikasi]` Bukti definisional di `Section/OutstandingClaim_Sprd.xml`:

| Baris | `pyPageListProperty` | Class elemen |
| --- | --- | --- |
| 1149/1158 | `.ClaimData.SpreadingClaim` | `ASM-FW-GISFW-Data-SpreadingRisk` |
| 4991/5000 | `.ClaimData.SpreadingBreakQS` | `ASM-FW-GISFW-Data-SpreadingRisk` |

⚠️ `[terverifikasi]` **Keduanya pada kedalaman XML yang sama (14) → dua repeat BERSAUDARA, bukan
bersarang.** Keterkaitannya lewat **subscript**, bukan penyarangan: `CountSpreading_Act.xml` 3749/4180
menghitung baris BreakQS dari `SpreadingClaim(Param.Index).ClaimSpreaded`.

`[terverifikasi]` Di level adjustment barulah **benar-benar bersarang** —
`AddAdjustment_Act.xml` 1048/1090:
`AdjustmentList(<LAST>).SpreadingAdjustment ← ClaimData.SpreadingClaim` dan
`AdjustmentList(<LAST>).SpreadingQuotaShare ← ClaimData.SpreadingBreakQS`.

**Peta tingkat:**

```
klaim ─┬─ SpreadingClaim        (1 tingkat)   ⬅ dapat diedit
       ├─ SpreadingBreakQS      (1 tingkat, SEJAJAR, seluruh sel pyReadOnly=true)
       └─ AdjustmentList ──┬─ SpreadingAdjustment   (2 tingkat)
                           └─ SpreadingQuotaShare   (2 tingkat)
```

**Bandingkan Claim Life:** `SpreadingList` (per treaty-year) → `RetroLifeList` (per reinsurer),
dua tingkat bersarang, dimensi treaty-year → reinsurer. **Prop tidak begitu.** Jangan salin bentuknya.

### 14a. ⚠️ Rumus — **enam cabang, dan dua di antaranya berbasis pengali berbeda**

`[terverifikasi]` `CountSpreading_Act` (class `Work-ClaimTreaty`, basis **estimasi**):

```
SpreadingClaim.ClaimSpreaded, cabang Local.SizeEst==1  (2261) :
  = SharePercentage × ClaimData.TotalEstimasi / 100                        (2176)
SpreadingClaim.ClaimSpreaded, cabang Local.SizeEst>1   (3634) + CurrID cocok (3541):
  = SharePercentage × Local.Estimation / 100                               (3456)
SpreadingBreakQS.ClaimSpreaded, dua cabang serupa:
  = (SpreadingClaim(Index).ClaimSpreaded × SharePercentage) / 100          (3748 / 4179)
```

`[terverifikasi]` `CountSpreadingADJ_Act` (class `Data-Adjustment`, basis **adjustment**) — **empat
cabang**:

| Cabang | Rumus | Gerbang |
| --- | --- | --- |
| A | `.ClaimSpreaded = Local.Adjustment × Local.Sharepersen / 100` (1945) | `paymenttype == 1 \|\| 2 \|\| 5` (2376) |
| B | `.ClaimSpreaded = Local.SpreadForBreak × Local.Sharepersen / 100` (2621) | `CountSpreadQS != 0` (2910); `SpreadForBreak` diisi hanya bila `.TreatyType == "10196" \|\| "10184" \|\| "10085"` (2324) |
| C | `.ClaimSpreaded = @divide(.SharePercentage,100,4) × Local.TotalAllADJ` (3350) | `.CurrencyID == Local.CurrID` (3433) |
| D | `.ClaimSpreaded = @divide(.SharePercentage,100,4) × Local.TotalAllADJ` (3580) | `.CurrencyID == Local.CurrID` (3663) |

⚠️⚠️ `[terverifikasi]` **Cabang B dan D menulis kolom yang sama dengan basis pengali BERBEDA**:
B memakai `ClaimSpreaded` baris induk, D memakai `Local.TotalAllADJ` (total adjustment per mata
uang). **Mana yang dimaksud tidak terbaca** → `[terbuka]`.

⚠️ `[terverifikasi]` `Local.TotalAllADJ` punya **dua makna** dalam satu activity: akumulator
(1411, syarat `.CurrencyID=Local.CurrID && .AcceptanceStatus!=2`) lalu penugasan langsung
(`= .AdjustmentValue`, 3205).

⚠️ `[terverifikasi]` **Tiga kode treaty type di-hardcode**: `"10196"`, `"10184"`, `"10085"` (2324).
Artinya **tidak terbaca dari korpus** → `[terbuka]`.

⚠️ `[terverifikasi]` `Local.Adjuster ← .AdjusterFeeValue` (473) **tidak pernah dipakai** — variabel mati.

`[terverifikasi]` Validasi share: `"Persen Share tidak boleh lebih besar dari 100"` (397), dipicu
`@greaterThan(Local.TotalPersen,100)` (1708).

### 14b. Pemicu dan persistensi

`[terverifikasi]` `CountSpreading_Act` dipicu **onChange grid** di banyak Section, plus dipanggil dari
`CountEstimation_Act` 7502 dan `DeleteEstimation_Act` 2870. `CountSpreadingADJ_Act` **hanya satu
pemanggil**: `CountValueADJTreaty_Act.xml` 4328.

`[terverifikasi]` **Tidak ditemukan penulis spreading ke Oracle** — nol di `RDBList/`, `ConnectREST/`,
`DataTransform/`, `When/`, `Flow/`, `FlowAction/`, `ReportDefinition/`. `InsertJsonClaimTreaty_act`
**tidak menyebut spreading sama sekali**. Nilainya **direkomputasi** tiap perubahan. Konsisten §7a.

---

## 15. Catastrophe + Cause of Loss

### 15a. Dua tingkat master, dihubungkan `M_COL_ID`

`[terverifikasi]` Master `M_CAUSE_OF_LOSS` berkolom `.M_COL_ID`, `.OLD_M_COL_ID`, `.COL_DESC` —
**hanya tiga** (`ReportDefinition/BrowseVMCauseOfLoss_RD.xml` 577–609). Detail
`D_CAUSE_OF_LOSS` berkolom `.D_COL_ID`, `.OLD_D_COL_ID`, `.DESCRIPTION`, `.LOSS_CODE`,
**`.M_COL_ID`** (FK), `.STS_AKTIF` (`SelectVDCauseOfLoss_RD.xml` 590–671).

`[terverifikasi]` **Tingkat ketiga**: `V_D_CAUSE_OF_LOSS_BUSINESS` menghubungkan detail ke **lini
bisnis** — `RDBList/GetLBUID_SQL.xml` 85:
`select ID, NOTE from V_D_CAUSE_OF_LOSS_BUSINESS a, BUSINESS b where a.BISNISID=b.ID and D_COL_ID={…}`.

### 15b. ⚠️ Dua procedure master menerima **CLOB JSON** yang menyamar sebagai kolom teks

`[terverifikasi]` `RDBList/UpdateMCauseOfLoss.xml` 60–71 dan `UpdateDCauseOfLoss.xml` 85–96 keduanya:

```sql
DECLARE IDPega VARCHAR2(32767); Datapega CLOB;
BEGIN
  IDPega := {InputData.M_COL_ID};            -- resp. D_COL_ID
  DBMS_LOB.CREATETEMPORARY(Datapega, true);
  Datapega := {InputData.COL_DESC};          -- resp. OLD_D_COL_ID
  POOLDATA.PEGA_M_CAUSE_OF_LOSS(Datapega, IDPega, {OutputData.COL_DESC out});
  COMMIT;
END;
```

⚠️⚠️ `[terverifikasi]` **`InputData.COL_DESC` dan `InputData.OLD_D_COL_ID` BUKAN deskripsi maupun ID
lama** — keduanya diisi `@GCNM.GetPageJSONString()` (`CNMInsertCauseOfLoss_act.xml` 424;
`CNMInsertDetailCauseOfLoss_act.xml` 426). Parameter CLOB itu **payload JSON**. Nama parameter
berbohong.

`[terverifikasi]` Parameter OUT dipakai sebagai **pesan status**, bukan data — hasilnya masuk ke
`TempCauseOfLoss.pyNote` / `TempDcol.pyNote`.

### 15c. Katastrofa punya **identitas event sendiri** yang mengelompokkan banyak klaim

`[terverifikasi]` Flag di klaim: `ClaimData.StsKatastrofe` bernilai `'Catastrophe'` /
`'Non-Catastrophe'`; pendamping `NonKatastrofeType` (`'Claim'` / `'Big Claim'`), `KatastrofeID`,
`KatastrofeNote`, `EditCatastrope`.

`[terverifikasi]` **Master event** berkelas `ASM-FW-GCNMFW-Int-CATASTROPHE` dengan kolom `.ID`,
`.NOTE`, `.STSKATASTROFE`, `.NONKATASTROFETYPE`, `KLAIMTYPE`, `USER_INPUT`, `TGL_INPUT`.
Nomor event dibentuk di `SaveCatasrtope_Act.xml` 331:

```
ParamCat.ID = "CTS-" + @pxReplaceAllViaRegex(@DateTime.CurrentDateTime(), "[GMT.]", "")
```

`[terverifikasi]` Penempelan ke klaim (`SetCatastrope_act.xml` 267–315) hanya men-set
`ClaimData.KatastrofeNote ← .NOTE` dan `ClaimData.KatastrofeID ← .ID`. Karena beberapa klaim dapat
memilih baris master yang sama, **satu `KatastrofeID` mengelompokkan banyak klaim** — ini **entitas
lintas-klaim**, dan tidak punya padanan di Claim Life.

⚠️ `[terverifikasi]` `SetDefNonCatastrope_Act` — **namanya JUJUR** (langka di modul ini): bila
katastrofa → kosongkan `NonKatastrofeType`; bila bukan → default `"Claim"` dan bersihkan
`KatastrofeID`/`KatastrofeNote`. Yang salah hanya **ejaannya** (`Catastrope`), konsisten di seluruh
nama rule.

⚠️ `[terverifikasi]` `InsertCatastrope_SQL` (dipakai `SaveCatasrtope_Act` 601) **tidak ada di
korpus** → `[terbuka]`. Daftar nilai dropdown `StsKatastrofe`/`NonKatastrofeType` bersumber
`pyListSource=associated` — rule properti itu **tidak diekspor**, jadi daftar lengkap opsinya
**tidak terbaca**.

### 15d. Cause of Loss **adalah gerbang wajib**

`[terverifikasi]` `ProteksiData_act.xml` step 2841 (deskripsi *"make sure cause of loss has been
choose"*): pesan `"Please fill Cause of loss"` bila `ClaimData.CauseOfLoss==""` (2963). Karena
`ProteksiData_act` dipanggil sebagai step 1 `SaveOutstanding_Act`, ini **memblokir simpan**.

⚠️ `[terverifikasi]` Tiga keganjilan pada `GetNameCauseofLoss_Act` (penulis `ClaimData.CauseOfLoss`):
class-nya **`ASM-FW-GISFW-Int-CLAUSE`** (bukan cause-of-loss); `Obj-Save`-nya memakai
`pyStepsClassName = ASM-FW-GCNMFW-Work-PNC` (378) — rujukan class usang; dan properti sumber `.Info`
**tidak ada** dalam daftar kolom RD yang memicunya.

⚠️ `[terverifikasi]` **`ClaimData.CauseOfLossID` dibaca** (`SaveOutstanding_Act` 4466/4756;
`CloseClaimProp` 1432) **tetapi tidak ada satu pun rule di korpus yang menulisnya** → `[terbuka]`.

⚠️ `[terverifikasi]` `Harness/TambahMasterCauseOfLoss.xml` — namanya berbohong: `pxInsName` =
`DATA-PORTAL!TAMBAHMASTERCAUSEOFLOSS`, **harness portal** pada class `Data-Portal`, bukan form tambah
master. Form sebenarnya `Section/BrowseCauseOfLoss.xml`.

---

## 16. Adjuster / Consultant — **satu master, dua peran**

⚠️ `[terverifikasi]` **SATU entitas master tanpa flag pembeda.** Master `ASM-FW-GISFW-Int-ADJUSTERCONSULTANT`
berkolom `.ID`, `.NAME`, `.EDITDATE`, `.USERNAME`, `.ADDRESS`, `.TELPNO` — **tidak ada kolom
`TYPE`/`ROLE`/`IS_ADJUSTER`**.

`[terverifikasi]` `SetAdjsuter_act` dan `SetConsultant_Act` adalah **klon**: RD sama
(`BrowseAdjusterConsultant`), class sama, satu-satunya beda adalah properti tujuan —
`AppointedADJID`/`AppointedADJ` versus `ConsultantID`/`ConsultantName`. Keduanya juga berbagi satu
sequence ID (`ADJUSTERCONSULTANT_SEQ`).

`[terverifikasi]` Keempat penunjuk di klaim adalah **skalar, bukan page list** → **1:1 per peran**.
Satu klaim = maksimal satu adjuster + satu consultant; satu master dapat dipakai banyak klaim.

`[terverifikasi]` **Wajib sebelum baris adjustment dapat ditambahkan** —
`AddAdjustment_Act.xml` 1859, pesan `"Please select the Adjuster and Consultant again."` (1885) bila
`ConsultantID == "" || AppointedADJID == ""` (1964). ⚠️ **Gerbangnya berbasis TAHAP, bukan ambang
nilai** — berlaku untuk semua klaim. Pencarian gerbang berbasis nilai: **nihil**.

⚠️ `[terverifikasi]` `GetIDConsultanAdj_SQL` **namanya berbohong** — ia bukan pembaca master,
melainkan **pembangkit ID**:
`SELECT ID||lpad(to_Char(ADJUSTERCONSULTANT_SEQ.nextval),4,'0') FROM POOLDATA.M_SITE_DATABASE WHERE CURRENT_SITE='1'`.

⚠️ `[terverifikasi]` **Nama tabel master adjuster/consultant TIDAK DIKETAHUI** — penyimpanannya
lewat `RequestType = SaveMasterAdjusterConsultant_sql` (`SaveAdjusterConsultant_Act` 1014), dan rule
itu **tidak ada di korpus** → `[terbuka]`.

### 16a. ⚠️ `AdjusterFee` **tidak masuk nilai klaim** — ia jenis pembayaran tersendiri

`[terverifikasi]` Rumus `.AdjustmentValue` (`CountValueADJTreaty_Act` 2573) **tidak memuat** suku
adjuster fee. Yang ada: akumulator kirim-ke-kasir di `HitServiceToKasir_Act`:

| Baris | Ekspresi | Gerbang |
| --- | --- | --- |
| 7677 | `TempKasir.CARI9 += .AdjustmentValue` | tanpa prasyarat |
| 7819 | `TempKasir.CARI9 += .AdjusterFeeValue` | `.PaymentType = 4 \|\| .PaymentType = 6` (7936) |
| 7961 | `TempKasir.CARI9 += .SalvageValue` | `.PaymentType == 3` |

⚠️ Ini **bukti korpus tambahan untuk OQ-020** yang selaras dengan §12c: `3` berkaitan salvage,
`4`/`6` berkaitan adjuster fee, `1` berkaitan klaim. Tetap **bukan definisi bisnis** — OQ-020 tetap
terbuka.

⚠️ `[terverifikasi]` `.AdjusterFeeValue` **tidak punya penulis maupun kontrol input di korpus**;
properti bernama mirip `.DataCommitteeTreaty.AdjusterFee` adalah field **di sisi komite**, dan tidak
ada rule yang menurunkan satu dari yang lain → `[terbuka]`.

---

## 17. Rule yang dirujuk tetapi **tidak ada di korpus** (kumulatif)

`[terverifikasi]` `GetBreakDownTreaty` (`SetTreatyNameSpreading_Act` 1901) · `InsertCatastrope_SQL`
(`SaveCatasrtope_Act` 601) · `SaveMasterAdjusterConsultant_sql` (`SaveAdjusterConsultant_Act` 1014) ·
`SetProtectionEstimation` (`ProteksiInitialandDate_Act` 1492) · rule properti daftar nilai
`StsKatastrofe`/`NonKatastrofeType` · definisi `@GCNM.GetPageJSONString()` dan
`@ASM.GetPageJSONString()`.

Penulis properti yang **tidak ditemukan**: `.AcceptanceStatus` (§11c) · `ClaimData.CauseOfLossID` ·
`.AdjusterFeeValue` · `InputData.CARI16/17/20` · `InputParamOs.Currency/TypeLoss` · isi page
`AttachCategory` · `ParamCari.CARI2` (§2d).

---

## 18. Efek keluar — **tanpa penanganan kegagalan sama sekali**

### 18a. `HitServiceToKasir_Act` — 14 langkah, nol jaring pengaman

`[terverifikasi]` **Tidak ada `Obj-Save` maupun `Commit` SEBELUM panggilan REST.** Satu-satunya
`Obj-Save` adalah langkah 14 (baris 9028, `WriteNow=true`, `WithErrors=false`) — **sesudah** kedua
`Connect-REST` (3150 dan 8431). Nol step `Commit` di seluruh berkas.

⚠️⚠️ `[terverifikasi]` **Penanganan kegagalan praktis tidak ada:**

| Mekanisme | Jumlah |
| --- | ---: |
| `Exit-Activity` | **0** |
| `Page-Set-Messages` | **0** |
| `Rollback` | **0** |
| `<pyOnException>` terisi | **0 dari 39** (termasuk pada kedua langkah `Connect-REST`) |
| retry | **0** |

`[terverifikasi]` Satu-satunya pemeriksaan hasil adalah precondition `Primary.pyStatusMessage=="OK"`
(3897), dan efeknya hanya **melewati** penandaan `.IsPrintAccept`. Blok `Java` (3070) menelan
exception: `catch(Exception e){ oLog.error(...); }`.

`[terverifikasi]` Bila REST gagal, langkah 9.8/13.5 **tetap jalan tanpa precondition** dan menulis
pesan error servis ke `.StatusKasir`; log tetap ditulis; `Obj-Save` tetap terjadi. **Fire-and-forget.**

`[terverifikasi]` Di server non-PROD kedua `Connect-REST` **dilewati total** (precond `IsPEGAPROD`,
3273/8554, deskripsi *"kalau diserver dev jangan dijalanin"*) — tetapi penulisan status dan log
**tetap berjalan**, sehingga data lingkungan dev tercatat seolah terkirim.

⚠️ **Tetapi ada gerbang belakangan:** `CloseClaimProp` memblokir penutupan klaim bila Kasir belum
sukses (§11b) dengan `exit activity`. Jadi kegagalan Kasir **tidak** membatalkan simpan/akseptasi,
tetapi **menghalangi Close**.

### 18b. `DIRECTTOKASIR_LOG` adalah **jejak**, bukan antrean

`[terverifikasi]` `RDBList/InsertLOGDirectKasir_SQL.xml` 85–94:
`INSERT INTO POOLDATA.DIRECTTOKASIR_LOG (DATA_JSON, IDPEGA, NOAKSEPTASI, KET) VALUES (…); COMMIT;`

Pembeda jejak-vs-antrean `[terverifikasi]`: **nol kolom status/retry**; seluruh korpus hanya punya
**1 INSERT + 1 SELECT** terhadap tabel ini (`GetStatusKasir_SQL.xml` 59, murni untuk tampilan);
**nol UPDATE, nol DELETE, nol agent/job** yang memprosesnya ulang. Bersifat *append-only* — kirim
ulang menambah baris baru.

⚠️ **Konsekuensi:** pola outbox transaksional (ADR-0015) **belum ada** di sistem lama; ia harus
dibangun, bukan ditiru.

### 18c. Bukan polling

`[terverifikasi]` `GetStatusKasir_Act` adalah **defer-load seksi UI**
(`Section/AdjustmentDetail.xml` 8271 `<pyDeferLoadRetrievalActivity>`), dengan
`pyRefreshRunActivity=false` dan `pyRefreshWhenActive=false` → **tidak ada interval auto-refresh**.
Hanya dua cabang nilai: `"Success"` → dipetakan ke teks `"Akseptasi Sudah Masuk ke Kasir"`; nilai
lain disalin apa adanya sebagai pesan error.

`[terverifikasi]` `getStatusKonversi_Act` adalah **pre-check satu kali**, dipanggil hanya dari
`HitServiceToKasir_Act` 751. SQL-nya menghitung baris di `reinsurance.trloss_detail_t`; bila nol,
`HitServiceToKasir_Act` **keluar** (transisi 784, false=exit). **Klaim hanya boleh dikirim ke Kasir
bila baris akseptasinya sudah ada di sistem inti.**

### 18d. `KonversiKlaim_Act` — **bukan konversi nilai**

`[terverifikasi]` Yang dikirim hanya tiga field: `caseId ← pzInsKey`, `noPolis ← PolicyNo`,
`stsReject`. Tidak ada konversi mata uang maupun format. Dipadankan dengan `getStatusKonversi_SQL`
yang memeriksa `reinsurance.trloss_detail_t`, isinya adalah **dorong/sinkron case klaim Pega ke
sistem inti reinsurance non-life**.

`[terverifikasi]` `.StatusService` **tidak pernah dibaca** → fire-and-forget.

| Pemanggil | `STSREJECT` |
| --- | --- |
| `SaveOutstanding_Act` 9584 | `"0"` |
| `SaveAcceptationTreaty_Act` 4203 | `"1"` |
| `CloseClaimProp` 2043 | `"4"` |

⚠️ Arti `0`/`1`/`4` **tidak terbaca dari korpus** — tidak ada decision table maupun field value yang
memetakannya → `[terbuka]`.

### 18e. Email

`[terverifikasi]` `SendEmailKlaim` dipicu **hanya** dari `AddKomiteTreatyChild_ACT` **langkah 34
(terakhir)**, tanpa precondition, **sesudah** `Obj-Save` langkah 33. Penerima `To` diambil dari
**roster komite pada case** (`.ComiteeClaim(1).KomiteEmail`) — **bukan** `GL.F_GET_EMAIL`, bukan
master. CC adalah **satu alamat tertanam** (baris 630, hanya bila `IsPEGAPROD`; nilai tidak disalin).

⚠️ `[terverifikasi]` `GL.F_GET_EMAIL` **memang ada** tetapi untuk hal lain: `GetEmailCeding_SQL`
mengisi field `Email` pada payload JSON ke Kasir.

`[terverifikasi]` Kegagalan email **tidak** menggagalkan apa pun: nol `Exit-Activity`, nol
`Page-Set-Messages`, `pyOnException` kosong, dan di non-PROD email **tidak dikirim sama sekali**
(precond `IsPEGAPROD`, 3848).

⚠️ `[terverifikasi]` **Satu alamat IP internal tertanam sebagai precondition** di
`SendEmailKlaim.xml` 2816 (`pxRequestor.pxReqContextURI`). Ini kembaran temuan `jboss1074` (§1b) —
konfigurasi lingkungan yang menggerbangi logika bisnis. Menyentuh **OQ-018**.

### 18f. Google Storage — dan satu kopling yang ketat

`[terverifikasi]` Alur: **token → upload → simpan metadata.**
`pooldata.GET_TOKEN_STORAGE({App}, {OperatorID.pyUserIdentifier}, {Kodestring OUT}, {ResponseMsg OUT})`
→ `Connect-REST ServiceGoogle POST` → `Insert_T_Storage_SQL`.

`[terverifikasi]` Berkas fisik di Google, metadata di Oracle (`T_STORAGE_IMAGE`:
`IMAGEID`, `URLPUBLIC`, `APPFOLDER`, `EXPDATE`, `FILENAME`, `APPNAME`, `STORAGE='standard'`).
`IMAGEID` = `STANDARD_HASH('ASMPP'||TO_CHAR(SYSTIMESTAMP,…),'MD5')`. URL **berbatas waktu**
(`Durasi=1800`), dengan jalur refresh bila `EXPDATE` lewat (`GetUrlGoogleStorage_Act` langkah 6).

⚠️ `[terverifikasi]` **Kopling ketat:** `InsertDocument_Act` **melewati `Obj-Save`** bila
`NewDocument.T_STORAGE_ID==""` (precond 1283) — artinya **baris dokumen tidak tersimpan bila upload
Google gagal**. Sejalan ADR-0010 (berkas tetap di Google Storage), tetapi di sistem baru ini harus
menjadi keputusan sadar, bukan efek samping.

### 18g. Urutan efek keluar per tahap `[terverifikasi]`

| Tahap | Urutan |
| --- | --- |
| **Simpan Outstanding** | `SaveOutstanding_Act` → `GetLimitPLATreatyin` → set `IsPLA` → `KonversiKlaim_Act` (`STSREJECT="0"`) → modal *"Please Print Pla"* |
| **Akseptasi** | `SaveAcceptationTreaty_Act` → `PrintFileAcceptance` → `InsertJsonClaimTreaty_act` → `SaveAcceptation_Act` → `KonversiKlaim_Act` (`"1"`) · lalu aksi terpisah `HitServiceToKasir_Act` → REST Kasir + log + `Obj-Save` · lalu modal *"Please Print DLA"* |
| **Close** | `CloseClaimProp` → gerbang FlagError → `SaveDataToOsAkseptasiNP` → `KonversiKlaim_Act` (`"4"`) → history → `InsertJsonClaimTreaty_act` → `ASMForceCaseClose` |
| **Close tanpa bayar** | `SendCloseClaimToKomite` → `pxAddChildWork` → `Obj-Save` → `SendEmailKlaimRejectClose` |

---

## 19. Dokumen PLA / DLA / Acceptance

### 19a. Kepanjangan — PLA lemah, **DLA tidak terbaca**

⚠️ `[terverifikasi]` String `PRELIMINARYLOSSADVICECOINS` memang ada di korpus
(`TryMakePLA_Act` 5768/5771; `PrintDLATreatyIn` 6187/6190; `CreatClaimAnalysis_Act` 1223/1226) —
**tetapi seluruhnya di dalam `<pyExpression>` / cache expression-builder yang sudah tidak aktif**,
bukan `<PropertiesValue>` yang dieksekusi. Nilai aktif PLA adalah
`"PLA " + ReinsurerName + " " + NoPla + ".pdf"` (5753).

Jadi: **PLA = "Preliminary Loss Advice"** didukung **residu korpus**, bukan label aktif — `[dugaan]`
kuat, bukan `[terverifikasi]`.

⚠️⚠️ **DLA `[terbuka]`.** String yang sama di `PrintDLATreatyIn` adalah **salin-tempel identik
karakter demi karakter** dari rule PLA → **tidak sah** sebagai bukti. Tidak ada label, komentar,
judul dokumen, atau teks Section yang mengejanya. **Jangan tebak.**

Pembanding yang jelas: dokumen ketiga berkategori `"AcceptanceNote"`, nama berkas
`"Persetujuan Klaim AcceptNo <no>.pdf"`.

### 19b. Pemicu dan tahap

| Dokumen | Rule | Pemicu | Tahap |
| --- | --- | --- | --- |
| PLA | `TryMakePLA_Act` | `FlowAction/GeneratePLA.xml` 137 `pyLocalActionActivity` — **local action**, tombol "PRINT PLA" | Outstanding |
| DLA | `PrintDLATreatyIn` | `FlowAction/GenerateDLATreaty.xml` 136 — local action, tombol "Generate DLA" | Adjustment / Akseptasi |
| Acceptance | `PrintFileAcceptance` | **bukan** local action — **langkah 11 `SaveAcceptationTreaty_Act`** | otomatis saat akseptasi disimpan |

### 19c. Pencetakan **bukan gerbang**

`[terverifikasi]` `<pyValidateActivity/>` **kosong** pada `GeneratePLA`, `GenerateDLATreaty`,
`OutstandingClaim`, `InputAcceptation`, `AdjustmentDetail`, `PreventRejectClaimProp`.
`CloseClaimProp` **tidak memeriksa** `NoPla`, `IsPLA`, `DLA_No`, maupun `IsPrintAccept`.

Yang ada hanya gerbang UI **arah sebaliknya** — mencegah cetak ganda: tombol PLA disabled bila
`.IsPLA != 1`, dan `TryMakePLA_Act` **mereset `IsPLA` ke `0`** setelah cetak (6689); tombol DLA
disabled bila `.IsFacRetro != 1 || .DLA_No != ''`. Modal *"Please Print Pla"* / *"Please Print DLA"*
hanya **pesan**, tidak memblokir transisi.

### 19d. ⚠️ `TryMakePLA_Act` — "Try" menyesatkan, dan sequence bisa terbakar

`[terverifikasi]` Nol `Exit-Activity`, nol `Page-Set-Messages`, nol `pyOnException` terisi, nol
`Obj-Save`. Langkah 26 (`Java`, 6356) justru **melempar exception keras**:

```java
catch (Exception e) { throw new PRRuntimeException("Can't attach the file to the Work Object"); }
```

⚠️⚠️ `[terverifikasi]` **Nomor PLA sudah dibentuk sebelum PDF dibuat** —
`GenerateNoPLATreatyIn` (1714) dan `GetSequenceNumber_SQL` (2245), keduanya ber-`COMMIT;` sendiri
(§7f) — dan **tidak ada rollback** bila langkah 24–27 gagal. **Sequence terbakar tanpa dokumen.**
Ini kasus nyata mengapa cut point transaksi harus dipindahkan (§7f).

### 19e. Dokumen **disimpan**, bukan sekadar di-stream

`[terverifikasi]` Ketiganya memanggil `InsertDocument_Act` dengan `KATEGORI_1` = `"PLA"` / `"DLA"` /
`"AcceptanceNote"`, yang meneruskan ke Google Storage + `T_STORAGE_IMAGE` + instance dokumen Pega
(§18f). Streaming ke pengguna juga terjadi, tetapi bukan satu-satunya tujuan.

⚠️ `[terverifikasi]` Stream HTML `PLAHTMLEksternal`, `DLAHTMLEksternal`, `FILEAcceptanceNote`
**tidak ada di korpus** → isi dokumen tidak dapat direkonstruksi → `[terbuka]`.

### 19f. `GetLimitPLATreatyin` **bukan** gerbang limit PLA

`[terverifikasi]` `TryMakePLA_Act` **tidak memanggilnya sama sekali**. Pemakainya hanya
`SaveAcceptationTreaty_Act` 2141 dan `SaveOutstanding_Act` 7943 — keduanya di jalur **simpan**.
Kaitan ke `IsPLA` hanya **kedekatan urutan langkah** (7943 lalu 9144); ekspresi yang menghubungkan
hasil limit ke nilai `IsPLA` **tidak terbaca dari korpus** → `[terbuka]`.

---

## 20. Status Ronde 1

**Seluruh penelusuran korpus selesai.** Frontier fakta kosong; frontier **keputusan** belum —
Q1–Q13 menunggu jawaban work owner. Daftar pemblokir dan pertanyaan ada di §21.

## 21. Pemblokir terkumpul

| # | Pemblokir | Pemilik |
| --- | --- | --- |
| P1 | Dua rumus berlawanan `ListClaimAmount.Value` — deductible sebelum vs sesudah share (§10b) | Product + UW |
| P2 | `CountSpreadingADJ_Act` cabang B vs D — basis pengali berbeda (§14a) | Product + UW |
| P3 | `.AcceptanceStatus` dibaca 26×, penulisnya di luar korpus (§11c) | pemilik Komite Claim Prop |
| P4 | Arti `PaymentType` 1–6 dan `.Type` 1–4, serta asimetri keduanya (§1a, §12c, §16a) — **OQ-020** | Finance + Product/UW |
| P5 | Arti `STSREJECT` `0`/`1`/`4` (§18d) | Product + UW |
| P6 | Delapan stored procedure + `GET_TOKEN_STORAGE`, `GL.F_GET_EMAIL`, `GETCURRENCYSTANDARD` — body tidak ada (**OQ-002**) | DBA |
| P7 | Nomor "temp" vs final — aturan promosi (§4) | Product/UW + DBA |
| P8 | Arti akhiran `TNP` / `TRT`; `TNP` dipakai lintas modul (§7e) | DBA + Product |
| P9 | Kepanjangan **DLA** dan isi stream dokumen (§19a, §19e) | Product + UW |
| P10 | Arti `S` pada `CLMS-`/`CLMPS-`/`CLMNPS-` (§1) | Product + UW |
| P11 | Tabel master adjuster/consultant — namanya tidak diketahui (§16) | DBA |
| P12 | Enam rule dirujuk tapi tidak ada di korpus (§17) | pemilik export Pega |
| P13 | Arti enum `pyStepsPreCondParamsWhenTrue/False` (§8b) | pemilik export Pega |
| P14 | Tiga kode treaty type hardcode `10196`/`10184`/`10085`, dan `10001`, `10004` (§14a, §13a) | Product + UW |
