# Grilling Ronde 1 — Claim Fac In · SENSUS STRUKTURAL PENUH

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Fac In\` (READ-ONLY) · **482 berkas**
**Pembanding identitas rule:** `D:\XML\RNM_BRD\Claim Prop\` (READ-ONLY)

> ⛔ **Ronde ini MEMBACA, bukan merancang.** Nol kesimpulan untuk aplikasi Go.
> ⛔ Berkas lama NOL disunting · modul lain NOL · kode NOL · `CREATE TABLE` NOL · DDL NOL ·
> nilai rahasia NOL · nol nomor baris XML dikutip.
> ⛔ **Tidak satu pun butir `[terbuka]` dinyatakan tertutup.**

---

## §A — Titik mulai

### A1 — Aturan A–J dipakai sejak langkah pertama ✅

`ATURAN-BACA-KORPUS-PEGA.md` dipakai **sebagai alat, bukan sebagai temuan**. Yang paling menentukan
di ronde ini: **B** (pengurai XML sungguhan, bukan pembaca baris) · **C1** (gerbang hidup hanya bila
`pyStepsPreCondition = true`) · **C2** (arti keenam kode arah) · **C3** (kode 5 mengubah AND jadi OR) ·
**E2** (halaman berulang dibaca dari `pyStepsObjectName`) · **F2** (`//` mematikan langkah) ·
**I1** (penyaring tidak peka huruf besar-kecil) · **I2** (ketiadaan tidak boleh disimpulkan) ·
**I3** (hitung dua cara, laporkan selisihnya).

⭐ **Claim Prop butuh dua ronde untuk menemukan aturan itu. Ronde ini tidak mengulanginya** — dan
hasilnya, sensus penuh 482 berkas selesai dalam satu ronde.

### ⚠️ A1-RALAT — **DUA nilai muncul yang tidak ada di aturan baca**

Aturan A–J **tidak menampung keduanya**. ⛔ Tidak saya tebak artinya.

#### 1 · `pyStepsPreCondition` bernilai **`0`** — aturan C1 hanya menyebut `true` / `false` / kosong

`[terverifikasi]` **Enam langkah**, dan ⭐ **keenamnya bermetode `Property-Set-HTML`** — metode yang
**nol** di Claim Prop:

| Rule | Langkah |
| --- | --- |
| `CLaimFaceSheet_Act` | 39 |
| `DLAFacintoTreaty_Act` | 13.4 |
| `DraftGenerateDLAFacin_Act` | 13 |
| `GenerateDLAFacin_Act` | 16 |
| `GeneratePLATreaty_Act` | 18.23 |
| `PrintPDFAccep_MultiAksep` | 23 |

⚠️ Keenamnya bersyarat **kosong**, jadi akibatnya kemungkinan sama dengan flag kosong — tetapi
**itu dugaan, bukan bacaan**. `[terbuka]`.

#### 2 · `pyStepsRepeatDefHasRepeat` bernilai **`PROPERTYLIST`** — aturan E1 hanya menyebut `EMBEDDED` / `REPEAT`

`[terverifikasi]` **Satu langkah**: `SendEmail_ACT` langkah **19**, metode `Property-Set`,
`pyStepsObjectName` **kosong**. `[terbuka]` apa bedanya dengan `EMBEDDED`.

### A2 — ⭐ `DataPage`: **usulan aturan baca**, BUKAN aturan yang disahkan

⛔ `ATURAN-BACA-KORPUS-PEGA.md` **tidak disunting**. Yang di bawah adalah **usulan**, menunggu
pengesahan work owner.

> ### USULAN ATURAN K — `Rule-Declare-Pages` (DataPage)
>
> **K1. Kelas rule-nya `CODE-PEGA-LIST`**, bukan kelas kerja modul. Delapan berkas
> `DataPage/` di Claim Fac In seluruhnya berkelas itu.
>
> **K2. Bentuk halamannya di `pyStructure`** — `list` atau `page`. Kedelapan-delapannya `list`.
>
> **K3. Umur halamannya di `pyScope`** — `thread` · `requestor` · `node`. Kedelapan-delapannya
> `thread`, artinya **hidup sepanjang satu utas kerja**, bukan sepanjang sesi.
>
> **K4. ⭐ SUMBER DATANYA di `pyDeclarePagesDataSource`**, dan nama sumbernya di medan pasangannya:
>
> | Nilai `pyDeclarePagesDataSource` | Nama sumbernya dibaca di |
> | --- | --- |
> | `LoadActivity` | `pyLoadActivity` |
> | `DataTransform` | medan transform |
> | *(connector)* | `pyConnectorList` |
>
> **K5. Parameternya di `pyParametersParamName`** berikut `…ParamInOut` · `…ParamType` ·
> `…ParamReq`.
>
> **K6. `pyRefreshStrategy`** menentukan kapan ia dimuat ulang. Kedelapan-delapannya **`never`** —
> ⚠️ artinya **sekali dimuat, isinya tidak pernah disegarkan dalam utas itu**.
>
> **K7. Jumlah sumbernya di `pySourceCount`.** Semua yang diperiksa bernilai **1**.

### A3 — Kesimpulan Claim Prop **tidak disalin** ✅

Dua yang dilarang dibawa, dan keadaannya di sini:

| Yang dilarang dibawa | Keadaan di Claim Fac In |
| --- | --- |
| kelas kerja `WORK-CLAIMTREATY` | ⛔ **tidak dipakai.** Kelas kerjanya **`ASM-FW-GCNMFW-WORK-PNC` (89 rule)** |
| penggolongan `ObjectList`/`ObjectItemList` sebagai **bukan tabel** | ⛔ **tidak berlaku.** `DATA-OBJECTITEM` **56** + `DATA-OBJECT` **50** = **106 rule**, dan `.ObjectItemList` adalah **halaman yang PALING SERING diulang** (32 kali). ⚠️ Membawa penggolongan Claim Prop ke sini akan **membuang struktur terbesar kedua modul ini** |

### A4 — Identitas rule dari `pxInsName` ✅ — dan itu **langsung terbukti berguna**, lihat §B3.

---

## §B — Inventaris dan identitas

### B1 — Dihitung ulang: **seluruh angka COCOK**

| Yang diperiksa | Brief | Hitungan saya | Cocok? |
| --- | --- | --- | --- |
| berkas | 482 | **482** | ✅ |
| Activity · RDBList · When · Section | 179 · 63 · 60 · 57 | **179 · 63 · 60 · 57** | ✅ |
| FlowAction · DataTransform · ReportDefinition | 33 · 28 · 27 | **33 · 28 · 27** | ✅ |
| Harness · ConnectREST · DataPage | 16 · 8 · 8 | **16 · 8 · 8** | ✅ |
| DecisionTable · Flow · SystemSettings | 1 · 1 · 1 | **1 · 1 · 1** | ✅ |
| langkah activity | 2 509 | **2 509** | ✅ |
| langkah ber-`//` | 112 | **112** | ✅ |
| langkah Java | 36 | **36** | ✅ |
| `pyMemo` berisi | 459 di 449 berkas | **459 di 449** | ✅ |
| `pxInsName` unik | 470 dari 482 | **470 dari 482** | ✅ |

⛔ **Nol angka berbeda.** Tidak ada yang perlu diralat.

### B3 — Selisih 12 `pxInsName`: **TERJELASKAN**

`[terverifikasi]` **Nol berkas tanpa `pxInsName`** — keempat ratus delapan puluh dua punya.

⭐ **Selisihnya karena 11 identitas rule muncul di LEBIH DARI SATU berkas ekspor**:

| `pxInsName` | Berkasnya |
| --- | --- |
| `ASM-FW-GISFW-DATA-QUOTATION!CEDINGCEDANT` | ⭐ **TIGA** — `FlowAction` + `Harness` + `Section` |
| `…DATA-ADJUSTMENT!INPUTADJUSTMENT` | FlowAction + Section |
| `…DATA-OBJECTITEM!ESTIMASI` | FlowAction + Section |
| `…DATA-OBJECTITEM!ESTIMASIPA` | FlowAction + Section |
| `…DATA-OBJECTITEM!VIEWATTACHMENT` | Harness + Section |
| `…WORK-PNC!INPUTESTIMASI` | FlowAction + Section |
| `…WORK-PNC!INPUTREGISTER` | FlowAction + Section |
| `…WORK-PNC!PREVENTREJECTCLAIM` | FlowAction + Section |
| `…WORK-PNC!PROTECTDOL` | FlowAction + Section |
| `…DATA-FACOFFER!SHOWSECURITYREINSURER` | FlowAction + Section |
| `DATA-PORTAL!MSTADJUSTERCONSULTANT` | Harness + Section |

**Hitungannya tepat:** 10 pasangan × 1 lebihan + 1 rangkap-tiga × 2 lebihan = **12**. ✅

⭐ **Ini menjawab §D4 sekaligus**: kesebelas itu **rule yang SAMA**, bukan kebetulan nama —
dibuktikan `pxInsName`, bukan nama berkas.

### B2 — Peta kelas: **62 kelas berbeda**

| Kelas | Jumlah | Mewakili apa | Ada di Claim Prop? |
| --- | --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-PNC` | **89** | ⭐ **objek kerja klaim Fac In** | ⛔ **TIDAK** — Prop memakai `WORK-CLAIMTREATY` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM` | **56** | rincian item di dalam objek pertanggungan | ⛔ tidak sebagai kelas |
| `ASM-FW-GCNMFW-DATA-OBJECT` | **50** | objek pertanggungan | ⛔ tidak sebagai kelas |
| `@BASECLASS` | 44 | rule lintas-aplikasi | ✅ ya |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | 43 | baris penyesuaian | ✅ **ya** — kelas yang sama |
| `ASM-FW-GCNMFW-WORK` | 25 | objek kerja induk | ✅ ya |
| `ASM-FW-GISFW-INT-POLICYJSON` | 18 | antarmuka JSON polis | ✅ ya |
| `ASM-FW-GISFW-DATA` | 13 | data bersama | ✅ ya |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 10 | penyimpanan berkas | ✅ ya |
| `ASM-FW-GISFW-DATA-QUOTATION` | 10 | data penawaran | ⚠️ ada, jauh lebih kecil |
| `ASM-FW-GISFW-WORK` | 9 | objek kerja GISFW | ✅ ya |
| **`CODE-PEGA-LIST`** | **8** | ⭐ **kelas DataPage** | ⛔ **TIDAK ADA** |
| `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | 7 | antarmuka tabel akseptasi | ✅ ya |
| `ASM-FW-GCNMFW-DATA-CLAIMDATA` · `DATA-PORTAL` · `INT-V_POLIS` | 6 tiap | data klaim · portal · view polis | sebagian |

⚠️ **Dua kelas terbesar kedua dan ketiga — `DATA-OBJECT` dan `DATA-OBJECTITEM` — tidak punya
padanan kelas di Claim Prop.** Di sana konsep itu hidup sebagai halaman di dalam kasus
(`ObjectList`/`ObjectItemList`); di sini ia **kelas rule tersendiri dengan 106 rule**.

### B4 — Rule bersama vs khas: **dihitung ulang, COCOK**

| | Brief | Hitungan saya |
| --- | --- | --- |
| beridentitas sama dengan Claim Prop | 162 | **162** ✅ |
| `pxUpdateDateTime` **IDENTIK** | 158 | **158** ✅ |
| **VERSI BEDA** | 4 | **4** ✅ |
| khas Claim Fac In | 308 | **308** ✅ |

#### ⚠️ Keempat yang beda versi — apa bedanya, mana yang lebih baru

| `pxInsName` | Claim Fac In | Claim Prop | Yang lebih baru |
| --- | --- | --- | --- |
| `ASM-FW-GISFW-INT-CURRENCY!BROWSECURRENCY_RD` | 2020-03-02 03:52:58 | 2019-05-21 15:12:42 | ⭐ **Fac In**, lebih baru **10 bulan** |
| `ASM-FW-GISFW-INT-CURRENCYSTANDARD!ASM!CURRENCYSTANDARD` | 2022-04-14 06:15:17 | 2018-01-12 06:08:01 | ⭐ **Fac In**, lebih baru **4 tahun 3 bulan** |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE!GETMIMETYPE` | 2026-02-19 01:50:45 | 2025-09-04 15:16:37 | ⭐ **Fac In**, lebih baru **5 bulan** |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE!RNM!GENERATEIMAGEID_SQL` | 2025-09-10 14:53:38 | 2025-09-04 14:39:33 | ⭐ **Fac In**, lebih baru **6 hari** |

⭐ **Keempatnya lebih baru di Claim Fac In.** ⛔ **Tidak saya simpulkan mana yang benar** — versi
lebih baru belum tentu versi yang berlaku di produksi untuk kedua modul. `[terbuka]`.

⚠️ **Ketiga wilayah yang disentuh keempatnya adalah wilayah angka dan dokumen:** daftar mata uang ·
kurs standar · jenis berkas dan pembuat pengenal berkas. ⚠️ **Kesimpulan Claim Prop tentang kurs
standar dan penyimpanan berkas berdiri di atas versi yang LEBIH TUA.**

### B5 — Dari 158 rule bersama: **150 menyentuh uang atau efek keluar**

`[terverifikasi]` menyentuh **uang** 141 · menyentuh **efek keluar** 71 · salah satu/keduanya
**150** dari 158 *(95%)*.

⚠️ **Artinya hampir seluruh rule bersama itu memegang angka atau menembak keluar.** ⛔ Merujuk
kesimpulan Claim Prop untuk rule-rule itu **boleh**, tetapi **keempat yang beda versi harus dibaca
sendiri** — dan keempatnya justru di wilayah itu.

---

## §C — Sensus langkah activity: 179 rule · 2 509 langkah

⛔ Dibaca dengan **pengurai XML sungguhan** (aturan B).

### C1 — Empat keluarga gerbang

| Keluarga | Baris |
| --- | --- |
| **1** `pyStepsPreCondParams` | **2 631** — ⚠️ **lebih banyak dari jumlah langkah**, karena satu langkah bisa punya beberapa baris syarat |
| **2** `pyStepsTransParams` | **2 508** — ⚠️ **satu kurang dari 2 509 langkah** |
| **3** `pyStepsRepeatDef` | terisi **521** |
| **4** `pyStepsCallParams` / `pyParamArray` | menyertai tiap langkah bermetode |

⚠️ `[terbuka]` **Satu langkah tidak punya baris keluarga-2**, padahal aturan **D1** menyatakan satu
baris per langkah *(terbukti 552/552 di Komite, 1 365 di Claim Prop)*. **Tidak saya kejar di ronde
ini.**

**Sebaran kode arah keluarga-1** *(12 teratas dari 2 631)*:

```
2/2  1424    2/3   740    _/3   261    3/2    95    3/_    44    6/2    14
2/6    11    1/2    11    _/_     7    5/2     5    1/_     5    _/1     4
```

**Keluarga-2** *(2 508)*: `2/2` **2 478** · `_/_` 12 · `6/_` 5 · `1/2` 4 · `6/2` 3 · sisanya 1–1.

**Kode khusus:** **1 = 23** · **4 = 1** · **5 = 6** · **6 = 30**.

> ### ⭐ PERBEDAAN STRUKTURAL TERBESAR DARI CLAIM PROP
>
> | Kode | Claim Prop | Claim Fac In |
> | --- | --- | --- |
> | **1** *(Jump to Later Step)* | ⛔ **NOL** | ⭐ **23** |
> | **4** *(Exit Iteration)* | ⛔ **NOL** | **1** |
> | **5** *(Skip Whens)* | 7 | **6** |
> | **6** *(Exit Activity)* | 16 | **30** |
>
> **Claim Prop tidak memakai lompatan sama sekali. Claim Fac In memakainya, dan memakainya untuk
> BERCABANG MENURUT LINI PRODUK.**

**Kedua puluh tiga lompatan dan sasarannya** — 22 baris berparameter:

| Rule | Langkah | Kode | Sasaran |
| --- | --- | --- | --- |
| `CLaimFaceSheet_Act` | 4 · 5 | `1/_` · `_/1` | `NoEx` |
| `CloseClaim` | 6 | `1/_` | `END` |
| `DeleteLocation_Act` | 2 | `2/1` | `SAVE` |
| **`GetAllData_Act`** | **11** | `1/_` · `1/2` ×2 | ⭐ **`MBU` · `TRAVEL` · `PA`** |
| **`GetAllData_Act`** | **16.1.1.1.1** | `1/2` ×3 | ⭐ **`TRAVEL2` · `MBU2` · `PA2`** |
| **`GetAllData_Act`** | **16.1.2.1.1** | `1/2` ×3 | ⭐ **`TRAVEL3` · `MBU3` · `PA3`** |
| `GettsiAneka_Act` | 1.1.1.1.2.1 · 1.1.1.1.2.2 · 4 | `1/2` | `MBD` ×2 · `MBD1` |
| `SetNilaiResikoSendiri` | 18 | `_/1` | `TO` |
| `SetProtectionEstimation` | 12 · 13 · 19.1 · 21 | `1/_` · `1/3` | `TO` · `B` · `TO` · `Er` |
| `TravelDocument_act` | 3 | `_/1` | `TO` |

⭐ **`GetAllData_Act` sendiri memuat sembilan lompatan bernama lini produk** — `MBU` *(kendaraan
bermotor)*, `TRAVEL`, `PA` *(personal accident)*, masing-masing tiga kali dengan akhiran `2` dan
`3`. **Modul ini bercabang per produk dengan lompatan, bukan dengan gerbang.**

⚠️ **Lompatan menggantung: SATU** *(aturan F3)* — `SetProtectionEstimation` langkah **26** menunjuk
label `END` yang **tidak ada** di activity itu. ⭐ Dan flag prakondisinya **`false`**, jadi
gerbangnya mati — **persis pola yang aturan F3 katakan: seluruh lompatan menggantung sejauh ini
ditemukan pada langkah bergerbang mati.** Pola itu **bertahan**.

### C2 — Flag prakondisi

| Nilai | Langkah |
| --- | --- |
| `true` | **942** |
| `false` | **130** |
| kosong | **1 431** |
| ⚠️ **`0`** | **6** — lihat A1-RALAT |

⭐ **121 baris bersyarat berada pada langkah berflag MATI** — syaratnya tertulis, **tetapi tidak
berlaku**. Itu 121 gerbang yang akan salah dibaca oleh siapa pun yang tidak memeriksa flag-nya.

### C3 — Arah gerbang: **dihitung DUA CARA** *(aturan I3)*

Basis: **1 061** baris bersyarat pada langkah berflag `true`.

| | NORMAL | TERBALIK | lain |
| --- | --- | --- | --- |
| **Cara A** — `_` dihitung *lain* | 683 | **115** | 263 |
| **Cara B** — `_` = *tanpa tindakan* = lanjut = NORMAL | **930** | **115** | 16 |

⭐ **Selisihnya 247 baris, seluruhnya berkode `_` di sisi benar.** ⛔ **Tidak saya pilih salah
satu** — keduanya ditulis, sesuai aturan I3. **Angka TERBALIK sama di kedua cara: 115.**

**Bandingkan Claim Prop:** 44 terbalik dari 482 baris aktif *(9,1%)*. **Claim Fac In: 115 dari
1 061 *(10,8%)*.** ⭐ **Proporsinya hampir sama** — gaya penulisannya sama, hanya skalanya beda.

### C4 — Langkah ber-remark: **112**

⭐ `[data work owner 2026-09-19]` `//` **memang mematikan langkah** — dipakai apa adanya, tidak
diuji ulang.

**Dua puluh dua di antaranya bermetode berat** — jalur simpan atau efek keluar yang **MATI**:

| Metode | Jumlah | Rule · langkah |
| --- | --- | --- |
| **`RDB-List`** | **16** | `BackToRegister_act` 3 · `CLaimFaceSheet_Act` 8 · 11 · `DLAFacintoTreaty_Act` 4 · `DraftGenerateDLAFacin_Act` 4 · 7 · `GenerateDLAFacin_Act` 3 · 6 · `GeneratePLATreaty_Act` 4 · `SaveAcceptation` 10 · `SaveAdjusterConsultant_Act` 5 · `SaveCatasrtope_Act` 2 · `SearchPolis_act` 8 · 10 · 11 · `SetPayableTo_act` 2 |
| **`Obj-Save`** | **5** | `GetPayAttachmentAdj_Act` 3.3.4 · 3.3.8 · `SetKomiteList_ACT` 4 · `SetPayableTo_act` 4 · `SetProtectionEstimation` 31 |
| **`Commit`** | **1** | `GetPayAttachmentAdj_Act` 3.3.9 |

⚠️ ⭐ **Bandingkan Claim Prop: di sana hanya 4 langkah berat yang mati. Di sini 22 — lima kali
lipat.** Tiga di antaranya *(`GetPayAttachmentAdj_Act` 3.3.4 · 3.3.8 · 3.3.9)* **rule yang SAMA**
dengan yang mati di Claim Prop — jadi ia **satu temuan bersama**, bukan dua.

### C5 — Langkah Java: **36**

| | |
| --- | --- |
| Menyentuh basis data | ⛔ **0** — tidak ketemu di pola `connect` · `executeRDB` · `Database` · `PreparedStatement` · `SQL` |
| Memanggil layanan luar | **1** |
| Menangani kegagalan | **13** |

⛔ Kode Java **tidak disalin**. Identifikasi dari `pyStepsJavaSource` **terisi**, bukan dari nama
metode — ⚠️ dan itu perlu, karena ejaan metode di modul ini **tidak konsisten**: `page-copy`
huruf kecil muncul **12 kali** berdampingan dengan `Page-Copy` huruf besar **54 kali** *(aturan I1)*.

### C6 — Perulangan: **521 langkah**

| Jenis | Jumlah |
| --- | --- |
| `EMBEDDED` | **409** |
| `REPEAT` | **111** |
| ⚠️ **`PROPERTYLIST`** | **1** — lihat A1-RALAT |

⭐ **Step Page TERISI 412 · kosong 109 · halaman berbeda 124.** Dibaca dari `pyStepsObjectName`
*(aturan E2)* — bukan dari `pyStepPage` maupun `pyStepsRepeatType`.

**Dua belas halaman yang paling sering diulang:**

| Halaman | Kali |
| --- | --- |
| ⭐ **`.ObjectItemList`** | **32** |
| `.SpreadingList` | 31 |
| `.EstimationList` | 26 |
| `.Adjustment` | 23 |
| ⭐ **`pyWorkPage.ClaimData.ObjectList`** | **22** |
| `pyWorkPage.OfferFacIn.FacRetroList` | 19 |
| `pyWorkPage.OfferFacIn.LocationList` | 11 |
| `pyReportContentPage.pxResults` | 10 |
| `TempData.ClaimData.ObjectList(1).ObjectItemList` | 10 |
| `.CoverageList` | 9 |
| `Spreading.pxResults` | 9 |
| `OutputData.pxResults` | 9 |

⭐ **Halaman yang paling sering diulang di modul ini adalah `ObjectItemList` dan `ObjectList`** —
tepat dua konsep yang di Claim Prop digolongkan **bukan tabel**. ⚠️ Ditulis di sini sebagai fakta;
⛔ **nol kesimpulan untuk Go**.

⚠️ `pyWorkPage.OfferFacIn.*` muncul dua kali di daftar — ⭐ **halaman penawaran fakultatif
(`OfferFacIn`) adalah struktur yang hidup di modul ini**, sedangkan di Claim Prop ia hanya
disinggung sekilas.

### C7 — `Obj-*` dan `Commit`: **52 langkah**

| Metode | Jumlah |
| --- | --- |
| `Obj-Save` | **27** |
| `Obj-Browse` | 12 |
| `Obj-Refresh-And-Lock` | 4 |
| `Commit` | 4 |
| `Obj-Open-By-Handle` | 3 |
| `Obj-Sort` | 2 |

> ## ⚠️ ANGKA YANG PALING KERAS DI RONDE INI
>
> **Dari 52 langkah penyimpanan dan pembacaan objek, hanya DUA yang punya jalur kegagalan**
> *(baris keluarga-2 bersyarat)*. **Enam lainnya ber-remark, jadi mati.**
>
> ⛔ **Nol kesimpulan untuk Go ditarik dari sini** — tetapi angkanya dicatat apa adanya.

### C8 — Parameter activity

`pyXMLSignature` **tidak ketemu** sebagai anak langsung elemen akar pada activity mana pun di
medan yang saya sisir. ⚠️ Sesuai aturan **I2**, kalimatnya dibatasi: *tidak ketemu di medan
`pyXMLSignature` pada tingkat akar*. ⛔ **Belum saya sisir sampai ke dalam** — `[terbuka]`,
dan ini **kekurangan ronde ini**, bukan pernyataan ketiadaan.

---

## §D — Layar: Section 57 · Harness 16 · FlowAction 33

### D1 — Lima keluarga gerbang

| Keluarga | Muncul | HIDUP | MATI *(`NEVER` / `1=2`)* |
| --- | --- | --- | --- |
| **tampil-sel** `pyUserData` / `pyDefaultUserData` | **3 608** blok · `pyVisible=OTHER` **214** | **247** | **50** |
| **tampil-layout** `pyIsVisibilityOption` | **135** | **126** | **17** |
| **disable** `pyDisabledNew` | **883** | **72** | 0 |
| **wajib** `pyRequiredNew` | **443** | **9** | 0 |
| **hanya-baca** `pyReadOnlyCondition` | **101** muncul | **101** | 0 |

⭐ **Total gerbang layar HIDUP: 555 · MATI: 67.**

⚠️ `pyVisible` dibaca dari dalam **sub-halaman** `pyUserData`, bukan anak langsung — aturan **I2**
pernah menggigit persis di titik ini di modul Komite, dan **tidak terulang**.

⚠️ **`pyRequiredNew` muncul 443 kali tetapi hanya 9 yang bersyarat hidup** — sisanya wajib-isi
tetap atau tidak wajib sama sekali. Angka yang mencolok, dicatat tanpa kesimpulan.

### D2 — Gerbang aksi tombol: **`pyActionConditions` TIDAK KETEMU**

⚠️ Sesuai aturan **I2**, kalimat yang sah: **tidak ketemu di medan `pyActionConditions`,
`pyActions`, dan `pyActionSetName`** pada Section, Harness, maupun FlowAction.

Yang **ada** adalah **`pyActionName` — 33 baris di 33 berkas**, tepat satu per FlowAction, jadi
ia **nama flow action itu sendiri**, bukan gerbang tombol.

⚠️ **Bandingkan Claim Prop: `pyActionConditions` ada 17 baris di 3 berkas.** ⭐ Di sini **nol**.
`[terbuka]` apakah modul ini memang tidak memakai gerbang aksi, atau memakainya lewat medan lain
yang belum saya sisir.

### D3 — Peta layar: **21 sambungan Harness → Section**

| Harness | Section yang dimuat |
| --- | --- |
| `CauseofLoss_Harness` | `CauseofLoss_Section` |
| `CedingCedant` | `CedingCedant` · `CedingCedantHierarki` |
| `ChoosePolis` | `InputInwardFacultativeDtl` · `ViewPolis` |
| `Comittee` | `ClaimComite` |
| `ListPaymentClaim_Harness` | `ListPaymentClaim_SC` |
| `ListPaymentPremi_Harness` | `PaymentPremiList_SC` |
| `ListPayment_Harness` | `ListPayment_SC` |
| `MstAdjusterConsultant` | `MstAdjusterConsultant` |
| `Outstanding` | `Outstanding_SC` |
| `Pla_Dtl_Harness` | `Pla_Dtl` |
| `PrintDLA_dtl_Harness` | `PrintDLA_dtl` |
| `RetroList_Harnness` | `RetroList_SC` |
| `TambahCauseofLoss` | `BrowseDetailCauseOfLoss` · `ListDetailCauseOfLoss` |
| `TambahMasterCauseOfLoss` | `BrowseCauseOfLoss` · `GridCauseOfLoss` |
| `ViewAttachment` | `ViewAttachment` |
| `ViewCedantpanels` | `CedingCedant` · `ViewCedantPanel` |

⚠️ **Medannya `pyStreamName`, bukan `pySectionName`.** Sisiran pertama saya memakai `pySectionName`
dan menghasilkan **nol untuk keenam belas Harness** — ⛔ angka itu **salah**, dan aturan **I2**
yang mencegahnya menjadi kesimpulan. Dicatat sebagai kejadian, bukan disembunyikan.

### D4 — Nama kembar: ⭐ **11 pasang, SELURUHNYA rule yang SAMA**

⛔ Berbeda dari Claim Prop, yang hasilnya **3 sepasang : 3 kebetulan**. Di sini **11 : 0** —
dibuktikan `pxInsName` di §B3.

⚠️ **Asumsi "sama seperti Claim Prop" akan salah di sini.** Penyebabnya terbaca: di modul ini
satu rule layar sering diekspor **sekaligus** sebagai Section, FlowAction, **dan** Harness.

---

## §E — DataPage · RDBList · ConnectREST

### E1 — DataPage: kedelapan dibaca **satu per satu** ✅

| Rule | Kelas data | Bentuk · umur · segar | Sumber | Parameter | Dipakai |
| --- | --- | --- | --- | --- | --- |
| `D_ANEKALIST` | `ASM-FW-GISFW-Data-Aneka` | list · thread · **never** | — | — | **5** berkas |
| `D_COVERAGEMBUCLAIMLIST` | `…Data-Coverage` | list · thread · never | — | — | 2 |
| `D_COVERAGEPACLAIMLIST` | `…Data-Coverage` | list · thread · never | — | — | 2 |
| `D_COVERAGETRAVELCLAIMLIST` | `…Data-Coverage` | list · thread · never | — | — | 2 |
| `D_FILTEREDCOVERAGEANEKALIST` | `…Data-Coverage` | list · thread · never | — | — | 2 |
| **`D_FILTEREDCOVERAGELIST`** | `…Data-Coverage` | list · thread · never | ⭐ **`LoadActivity` → `FilterCoverage_Act`** | `ObjectID` · `IndexPropertyItem` *(IN, STRING, tidak wajib)* | **4** |
| `D_FILTEREDPROPERTYITEMLIST` | `…Data-PropertyItem` | list · thread · never | — | — | **6** berkas |
| `D_OCCUPATIONLIST` | `…Data-Occupation` | list · thread · never | — | — | — |

⭐ **Ketiga nama produk muncul lagi di sini**: `MBU` · `PA` · `TRAVEL` · `Aneka` — **sama dengan
sasaran lompatan di `GetAllData_Act`**. Dua bukti bebas yang menunjuk hal yang sama: **modul ini
bercabang per lini produk**.

⚠️ `pyRefreshStrategy = never` pada **kedelapan-delapannya** — dicatat sebagai fakta.
⛔ **Nol kesimpulan untuk Go.**

### E2 — RDBList: **63 rule, 34 objek Oracle unik**

`[terverifikasi]` Keenam puluh tiga punya `pyBrowseSQL` terbaca — **nol yang kosong**.

| Aksi | Rule |
| --- | --- |
| `SELECT` | **50** |
| `BEGIN … END` *(blok anonim)* | **9** |
| `UPDATE` | 3 |
| `INSERT` | 2 |
| `DELETE` | 1 |

⭐ **9 rule membawa `COMMIT` tertanam di dalam teks SQL-nya.**

**Sepuluh objek Oracle terbanyak:**

| Objek | Sentuhan |
| --- | --- |
| ⭐ **`os_akseptasi_klaim`** | **27** |
| `json_polis` | 6 |
| `facinproduction` | 4 |
| `json_klaim` | 3 |
| `dual` | 3 |
| `agent` | 3 |
| **`reinsurance.trloss_detail_t`** | 3 |
| `proportionalarrg` | 3 |
| `t_storage_image` | 3 |
| `m_client` | 2 |

⚠️ ⭐ **`os_akseptasi_klaim` disentuh 27 kali — jauh mendominasi.** Di Claim Prop tabel yang sama
juga muncul, tetapi tidak sedominan ini.

⚠️ **Prefiks schema hampir seluruhnya TIDAK DITULIS.** Hanya `reinsurance.` yang muncul eksplisit.
⛔ Sama dengan cacat yang tercatat di Claim Prop — dicatat, tidak disimpulkan.

⚠️ **Alias SQL belum diperiksa satu per satu.** Di Claim Prop lima dari enam alias satu rule
**berbohong**. Di sini **belum saya periksa** — `[terbuka]`, dan ini **kekurangan ronde ini**.

### E3 — ConnectREST: **8 rule, 1 berautentikasi**

⛔ **NILAI KREDENSIAL TIDAK DICETAK, TIDAK DISALIN.** Dikonfirmasi — perintah yang dipakai hanya
membaca medan pengendali, bukan isinya.

| Rule (`pxInsName`) | Auth | Sumber endpoint |
| --- | --- | --- |
| `…WORK-PNC!HITDLACLAIMFACIN` | false | pengaturan |
| `…WORK!KONVERSIKLAIMNONLIFE` | false | pengaturan |
| ⭐ `…DATA-ADJUSTMENT!SENDACCEPTATIONTOKASIR` | **true** | pengaturan |
| `…INT-T_STORAGE_IMAGE!SERVICEGOOGLE` | false | pengaturan |
| `…DATA-ADJUSTMENT!GETPAYATTACHMENT` | false | pengaturan |
| `…WORK-PNC!GETPAYMENTCLAIM` | false | pengaturan |
| `…WORK-PNC!GETPREMIUMPAIDON` | false | pengaturan |
| ⚠️ `…WORK-PNC!GETPREMIUMPAIDONMARINE` | false | ⭐ **URL langsung**, bukan pengaturan |

⭐ **Satu yang berautentikasi adalah rule Kasir yang sama** dengan yang sudah tercatat di Claim
Prop — **rule bersama, satu sumber**.

⚠️ ⭐ **`GETPREMIUMPAIDONMARINE` mengambil endpoint dari URL langsung, bukan dari tabel
pengaturan** — satu-satunya dari delapan. `[terbuka]`.

⚠️ **Peringatan keamanan dicatat sekali lagi, tidak ditindaklanjuti sendiri:** kredensial Kasir
berada di dalam ekspor yang beredar dan **sebaiknya diganti sesudah migrasi**. Urusan tim Kasir.

### E4 — Jenis sisa

| Jenis | Berkas | Kelas berbeda | Kelas terbanyak |
| --- | --- | --- | --- |
| DataTransform | 28 | 9 | `WORK-PNC` 11 · `@BASECLASS` 4 |
| ReportDefinition | 27 | **20** | `INT-AGENT` 4 · `INT-RW` 4 |
| DecisionTable | 1 | 1 | `INT-T_STORAGE_IMAGE` |
| Flow | 1 | 1 | `WORK-PNC` |
| SystemSettings | 1 | 1 | `LINKSERVICE` |

---

## §F — Catatan pengembang dan apa lagi

### F1 — `pyMemo`: **459 catatan di 449 berkas**, **457 pasangan berbeda**

⚠️ **Penggolongan ini KASAR** — berdasarkan kata kunci, **bukan** pembacaan isi rule-nya satu per
satu. ⛔ Di Claim Prop cara itu **terbukti meleset lima kali lipat** *(ronde 3 menduga 2, sensus
tuntas ronde 6 menemukan 11)*.

| Golongan | Jumlah kasar |
| --- | --- |
| **(a)** administratif | **136** |
| **(b)** menerangkan | **191** |
| **(c)** kandidat **MENGUBAH** pembacaan | **118** |
| **(d)** kandidat **LARANGAN** | **12** |

> ⛔ **Angka (c) dan (d) di atas adalah KANDIDAT, bukan vonis.** Menguji keduanya menuntut membuka
> tiap rule-nya — **itu pekerjaan ronde 2**, dan ia yang paling menahan *(lihat F6)*.

#### Kedua belas kandidat larangan — seluruhnya

| Rule | Catatan |
| --- | --- |
| `Activity/GetCurencyCoverage_Act.xml` | *"Hapus jika ada spreadingClaim yg sama"* |
| `Activity/InsertJsonClaimNonMBU_act.xml` | *"hapus yg hapus attachment list"* |
| `Activity/SaveAdjustmenttoDB_ACT.xml` | *"hapus yg hapus attachment list"* |
| `Activity/SendEmailDLA_ACT.xml` | *"hapus when step 6"* |
| `Activity/SetCatastrope_act.xml` | *"hapus save"* |
| `Activity/SetDefNonCatastrope_Act.xml` | *"hapus save"* |
| `Activity/SetSalvageValue.xml` | *"hapus yg set cronologu"* |
| `Activity/SetValueAdjusterFee.xml` | *"hapus yg set chronologi"* |
| `FlowAction/InputAdjustment.xml` | *"hapus pre"* |
| ⭐ `RDBList/GetLimitPLADLA_Sql.xml` | *"ganti biar jgn ambil dari DLA tpai dari treatyLimit"* |
| `ReportDefinition/BrowseAgentNusaRe_RD.xml` | *"hapus param yang ga dipake"* |
| `When/isPA_PNC.xml` | *"tambah pyWorkPage dan hapus policy"* |

⚠️ **Dua yang paling menuntut pemeriksaan:** `GetLimitPLADLA_Sql` *(menyebut sumber batas nilai —
DLA versus batas treaty; itu angka uang)* dan kedua *"hapus yg set chronologi"* *(menyentuh jejak
audit)*.

⚠️ **Empat catatan berbunyi sama persis dengan yang ada di Claim Prop** — *"hapus yg hapus
attachment list"* · *"hapus save"* ×2 · *"tambah pyWorkPage dan hapus policy"*. Di Claim Prop
keempatnya terbukti **sudah dipatuhi**. ⛔ **Tidak saya asumsikan sama di sini.**

### F2 — `pyUsage`: **1 267 terisi, hanya 6 nilai berbeda**

`FLOW` **1 068** · `java` **187** · sisanya 4 kalimat bawaan Pega, salah satunya justru bermakna:
*"Copy Currency from data-object item into estimation"*.

⭐ **Satu nilai bermakna dari 1 267** — pola yang sama dengan Claim Prop *(870 terisi, 4 nilai)*.

`pyXMLSignature` — lihat **C8**, belum tersisir penuh, `[terbuka]`.

### F3 — Berkas yang **TIDAK tersentuh sensus mana pun**: ⭐ **60 — seluruhnya `When`**

`[terverifikasi]` 482 berkas · tersentuh **422** · **tidak tersentuh 60**.

> ## ⚠️ INI CELAH PADA RANCANGAN RONDE INI SENDIRI
>
> **Tidak satu pun §A–§F menyebut `When`.** Keenam puluh rule itu lolos bukan karena tidak penting,
> melainkan karena **brief ronde 1 tidak memberi mereka tempat**. Dicatat sebagai kekurangan
> rancangan, bukan temuan korpus.

**Sisiran cepat — halaman yang diuji keenam puluhnya:**

| Yang diuji | Kemunculan |
| --- | --- |
| ⭐ **`pyWorkPage.Quotation`** | **38** |
| `.BusinessType` | 25 |
| `.BusinessCode` | 10 |
| `.BusinessName` · `pyWorkPage.pyWorkIDPrefix` · `.StatusBusiness` · `.pzProductionLevel` · `IsSpreadingUW` · `.OperatorID` · `.pyWorkBasketList` | 1 tiap |

⚠️ ⭐ **38 dari 60 menguji `pyWorkPage.Quotation`** — halaman **penawaran**. Di Claim Prop, 52 dari
61 rule `When` menguji halaman yang **tidak ada** pada objek kerjanya dan karena itu **selalu
bernilai salah**. ⛔ **Apakah `Quotation` ADA pada objek kerja Claim Fac In belum saya periksa.**
`[terbuka]` — dan ini pertanyaan yang **berbeda hasilnya** antara kedua modul.

### F4 — Kesimpulan yang PALING RAWAN salah

⛔ **Ditunjuk, tidak diperbaiki.**

> ⚠️ **Paling rawan: §F1 — penggolongan 459 catatan menjadi (a) 136 · (b) 191 · (c) 118 · (d) 12.**

Ia disusun dari **kata kunci pada teks catatannya**, bukan dari membuka rule-nya. ⭐ **Di Claim Prop
cara persis ini meleset lima kali lipat**: ronde 3 menduga 2 catatan mengubah pembacaan, sensus
tuntas di ronde 6 menemukan **11**. Di sini bahkan penyaringnya lebih longgar, sehingga angka **118**
hampir pasti **terlalu besar** dan angka **12** mungkin **terlalu kecil**.

**Paling rawan kedua:** §C3 *"TERBALIK 115"*. Angkanya stabil di kedua cara hitung, tetapi ia
bersandar pada satu definisi — *kode `3` di sisi benar* — dan **tidak memeriksa** apakah baris itu
bagian dari rantai yang kode 5-nya mengubah maknanya menjadi OR.

**Paling rawan ketiga:** §D2 *"gerbang aksi tombol nol"*. Saya menyisir tiga medan. ⚠️ Claim Prop
punya 17 baris di medan `pyActionConditions`; kalau modul ini memakai medan keempat yang belum saya
kenal, pernyataan itu runtuh — **pola yang sudah empat kali menggigit proyek ini**.

### F5 — Pertanyaan untuk work owner — **3**

#### Pertanyaan 1 — Empat rule beda versi: yang mana yang berlaku di produksi?

**Apa yang ditanyakan.** Empat rule punya identitas sama di Claim Prop dan Claim Fac In tetapi
**versi berbeda**, dan **keempatnya lebih baru di Fac In** — selisihnya dari 6 hari sampai
**4 tahun 3 bulan** *(kurs standar)*. Keempatnya menyentuh **daftar mata uang, kurs standar, jenis
berkas, dan pembuat pengenal berkas**.

**Kenapa muncul.** Ekspor kedua modul diambil terpisah, dan rule bersama ikut terbawa di keduanya.

**Bedanya kalau A atau B.** **A — versi Fac In yang berlaku:** kesimpulan Claim Prop tentang **kurs
standar** dibuat di atas rule berumur 4 tahun lebih tua, dan perlu dibaca ulang. **B — tiap modul
memang memakai versinya sendiri:** berarti sistem berjalan punya **dua perilaku berbeda** untuk satu
identitas rule, dan aplikasi Go harus memutuskan mana yang ditiru.

**Apa yang tertahan.** Pembacaan kurs standar di **kedua** modul, dan apakah kesimpulan Claim Prop
tentang penyimpanan berkas perlu ditinjau.

#### Pertanyaan 2 — `pyStepsPreCondition = 0` dan `pyStepsRepeatDefHasRepeat = PROPERTYLIST`

**Apa yang ditanyakan.** Dua nilai muncul yang **tidak ada di aturan baca yang sudah disahkan**:
flag prakondisi bernilai `0` pada **6 langkah** *(semuanya `Property-Set-HTML`)*, dan penanda
perulangan bernilai `PROPERTYLIST` pada **1 langkah**.

**Kenapa muncul.** Aturan C1 dan E1 disusun dari dua modul yang **tidak memuat kedua nilai itu**.

**Bedanya kalau A atau B.** **A — keduanya setara dengan nilai yang sudah dikenal** (`0` = kosong,
`PROPERTYLIST` = `EMBEDDED`): aturan bacanya cukup ditambahi catatan. **B — artinya berbeda:**
tujuh langkah itu sedang dibaca **salah** oleh sensus ini, dan aturan A–J perlu direvisi sebelum
modul-modul berikutnya disensus dengan aturan yang sama.

**Apa yang tertahan.** Ketepatan §C2 dan §C6 di ronde ini, dan keandalan aturan A–J untuk **seluruh
modul yang belum disensus**.

#### Pertanyaan 3 — `pyWorkPage.Quotation`: ada atau tidak pada objek kerja Claim Fac In?

**Apa yang ditanyakan.** **38 dari 60** rule `When` menguji `pyWorkPage.Quotation`. Di Claim Prop,
52 dari 61 rule `When` menguji halaman yang **tidak ada** pada objek kerjanya, sehingga
**seluruhnya bernilai salah** dan diputuskan **tidak dimigrasikan**.

**Kenapa muncul.** Polanya terlihat sama; hasilnya mungkin **kebalikannya**, karena Fac In adalah
modul fakultatif yang memang lahir dari penawaran.

**Bedanya kalau A atau B.** **A — `Quotation` ADA:** ke-38 rule itu **hidup**, dan klasifikasi lini
bisnis di modul ini bekerja dengan mekanisme yang di Claim Prop mati. **B — TIDAK ADA:** ke-38
rule itu sisa impor seperti di Claim Prop, dan perlakuannya mengikuti keputusan yang sudah ada.

**Apa yang tertahan.** Seluruh pembacaan klasifikasi lini bisnis modul ini — dan **60 rule `When`
yang ronde ini belum sentuh sama sekali**.

### F6 — Yang paling menahan untuk ronde 2

| # | Yang dikerjakan | Kenapa menahan |
| --- | --- | --- |
| **1** | ⭐ **60 rule `When`** — ronde ini **tidak menyentuhnya sama sekali** | celah terbesar, dan 38 di antaranya bergantung pada Pertanyaan 3 |
| **2** | ⭐ **Uji 118 + 12 kandidat `pyMemo` dengan membuka rule-nya** | penggolongan kasar ronde ini **paling rawan salah**, dan di Claim Prop cara yang sama meleset lima kali lipat |
| **3** | **Baca `GetAllData_Act`** — 9 lompatan bernama lini produk | ia **tulang punggung percabangan modul ini**, dan bentuknya tidak ada padanannya di Claim Prop |
| **4** | **Sisir `pyXMLSignature` sampai ke dalam** *(C8)* | ronde ini hanya memeriksa tingkat akar; pernyataannya sengaja dibatasi |
| **5** | **Periksa alias SQL 63 RDBList satu per satu** | di Claim Prop lima dari enam alias satu rule berbohong; di sini belum diperiksa |
| **6** | **Baca keempat rule beda versi** *(Pertanyaan 1)* | menyentuh kurs dan penyimpanan berkas di **dua modul sekaligus** |
| **7** | **Sisir medan keempat untuk gerbang aksi tombol** *(D2)* | pernyataan "nol" bersandar pada tiga medan saja |

---

## Lampiran — bukti berkas lain tidak disentuh

Sidik jari MD5 atas **38 berkas** milik `claim-prop`, `komite-claim-prop`, dan
`ATURAN-BACA-KORPUS-PEGA.md` diambil sebelum ronde ini dan dibandingkan sesudahnya.

**Satu-satunya berkas baru: `grilling-ronde-1.md` di folder baru `.scratch/claim-facin/`.**
