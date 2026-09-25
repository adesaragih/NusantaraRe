# PENGETAHUAN — Modul Treaty In Adjustment (AS-IS)

> **Sifat berkas ini:** dokumen **pemahaman sistem lama**, bukan rancangan.
> Tidak ada DDL, tidak ada skema Golang, tidak ada komponen React di sini.
> Sumber tunggal: ekspor XML di `D:\XML_NURE\Treaty In Adjustment` (379 berkas).
> **Tanggal analisis:** 23 September 2026.

## Status pengisian

| Bagian | Isi | Keadaan |
|---|---|---|
| §1 | Inventaris ekspor dan hubungannya dengan ekspor Treaty In | terisi |
| §2 | Model data: empat pohon di dalam satu `JSONDATA` | terisi |
| §3 | Pintu masuk, picker, susunan layar, hak akses | terisi |
| §4 | Daur hidup addendum: jenis, penomoran, potret nilai lama | terisi |
| §5 | Mesin selisih (`TreatyEDM*`) | terisi |
| §6 | Penyimpanan dan tabel yang disentuh | terisi |
| §7 | Persetujuan dan penolakan | terisi |
| §8 | Layar `OldData` | terisi |
| §9 | Utang teknis, dan pemisahan "masih bisa terjadi" dari "pernah terjadi" | terisi |
| §10 | Akibat untuk migrasi React/Golang/Oracle | terisi |

**Tidak dikerjakan di sesi ini, dan disebut supaya tidak dikira sudah:** pemetaan atribut
`Int-treaty_in_edm` satu per satu ke model baru, peta telusur JSON milik modul ini, DDL, dan
verifikasi apa pun atas data produksi (§9.2 mendaftar delapan uji yang belum dijalankan).

Berkas ini ditulis bertahap: tiap bagian masuk ke disk begitu langkahnya selesai.

---

## 1. Inventaris ekspor

### 1.1 Angka pokok

| | Berkas XML |
|---|---|
| `Treaty In Adjustment` | **379** |
| `Treaty In` | **329** |
| Irisan (jalur relatif sama) | **323** |
| Hanya di Adjustment | **56** |
| Hanya di Treaty In | **6** |

Sebaran 379 berkas menurut tipe aturan:

| Tipe | Jumlah | | Tipe | Jumlah |
|---|---|---|---|---|
| Activity | 155 | | RDBList | 44 |
| Section | 68 | | FlowAction | 43 |
| DataTransform | 38 | | ReportDefinition | 18 |
| Harness | 6 | | When | 3 |
| ConnectREST | 1 | | DecisionTable | 1 |
| Menu | 1 | | SystemSettings | 1 |

### 1.2 Koreksi atas sesi sebelumnya: irisannya **bukan** byte-identik

`KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` menyebut 323 berkas itu "byte-identik". Diperiksa ulang dengan
membandingkan berkas demi berkas: **nol** pasang yang identik sebagai bita. Yang benar:

* **Seluruh 323 pasang berbeda pada dua tag ekspor** — `pyRuleFormStatusTime` dan
  `pyShowJavaWindowName` — yang dicap saat ekspor dijalankan, bukan saat aturan ditulis.
  Kedua ekspor diambil terpisah sekitar sembilan menit (berkas Treaty In bercap 10:52,
  Adjustment 11:01 pada 3 September).
* **Urutan elemen XML-nya pun berbeda**, sehingga perbandingan bita tidak pernah bisa sama.
* Setelah kedua tag itu dan artefak kode-ter-generate dibuang, **300 dari 323 pasang benar-benar
  sama isinya**, dan 23 sisanya berbeda hanya pada hal yang tidak mengubah perilaku:
  `pyStepPageReference`, daftar versi aturan yang ikut terbungkus, dan pengaturan toolbar editor.

**Kesimpulan sesi lalu tetap berdiri untuk sebagian besar berkas, tetapi TIDAK untuk semuanya.**
Yang menentukan bukan bita melainkan `pzInsKey` — dan diperiksa satu per satu, `pzInsKey` **tidak**
identik di enam jalur:

| Jalur | Versi di ekspor Adjustment | Versi di ekspor Treaty In |
|---|---|---|
| `Activity/AddSpreadingXOL.xml` | 01-01-96 | 01-01-95 |
| `Activity/LoadAttachmentData.xml` | 01-01-54 | 01-01-62 |
| `Activity/SetCurrName_Act.xml` | 01-01-56 | 01-01-94 |
| `Activity/TreatyInSubmit.xml` | 01-01-90 | 01-01-96 |
| `Activity/AddDeduction.xml` | 01-01-54 | 01-01-54 (`pzInsKey` berbeda, versi sama) |
| `FlowAction/CoBList.xml` | 01-01-55 | 01-01-55 (idem) |

Maka angka yang benar: **323 jalur bersama, 317 di antaranya ber-`pzInsKey` sama.** Dua ekspor yang
diambil berselang sembilan menit memuat **versi aturan yang berbeda** untuk keenam jalur itu, dan
yang menentukan versi mana yang berjalan adalah konteks penggunanya (`pzIndexOwnerKey`), bukan nama
berkas — `METODE` §2.0.

**Yang dibaca berkas ini adalah versi dari ekspor Adjustment**, untuk keenam jalur itu maupun
selebihnya. Diperiksa terhadap seluruh berkas yang menjadi bukti TDA-01 sampai TDA-15: **tidak satu
pun** berada di antara keenamnya, sehingga tidak ada temuan yang bergantung pada pilihan versi ini.

**Dan ini membatalkan `METODE` §8.1 dua kali:** irisannya bukan byte-identik, dan juga bukan
seluruhnya ber-`pzInsKey` sama.

Contoh jalur yang memang identik, diperiksa langsung — `Section/InputTreatyInOffer.xml`:

```
Treaty In            : RULE-HTML-SECTION DATA-PORTAL INPUTTREATYINOFFER #20251014T093249.589 GMT · GISFW 01-01-90
Treaty In Adjustment : RULE-HTML-SECTION DATA-PORTAL INPUTTREATYINOFFER #20251014T093249.589 GMT · GISFW 01-01-90
```

> **Aturan kerja yang lahir dari sini:** membandingkan dua ekspor Pega dilakukan atas `pzInsKey`,
> bukan atas isi berkas. Ekspor Pega membungkus **daftar versi lama aturan yang sama** di dalam
> `pyIncludedRuleXML` dan `pyRuleVersionsList`, dan bungkusan itu berbeda antar ekspor meskipun
> aturan aktifnya persis sama. Membandingkan isi berkas menghasilkan selisih palsu.

### 1.3 Maka: Adjustment bukan aplikasi terpisah

85 % permukaan ekspor Adjustment **adalah Treaty In**. Permukaan khasnya adalah **56 berkas**, dan
bentuk ke-56 itu langsung memberi tahu apa modul ini:

| Kelompok | Jumlah | Isi |
|---|---|---|
| Layar `OldData` | 20 | 9 Section + 8 Flow Action + 1 Harness + 2 Section pembungkus, semua ber-akhiran `OldData` |
| Pembuatan & penomoran addendum | 6 | `TreatyCreateEDM`, `TreatyInRevisi_post`, `TreatyInEDMSetValue`, `SetTreatyInEDM_Act`, `TreatyInSetEditPre`, `TreatyInSetAddendumToHistory` |
| Pemilih (picker) kontrak/addendum | 6 | `PickerTreatyInMaster`, `PickerTreatyInMasterRevisi` (Section + Flow Action), `TreatyLoadMasterJoinEdm`, `TreatyLoadMasterJoinEdmXOL` |
| Akses data addendum | 5 | `GetTreatyRevisionID`, `BrowseTreatyInEDM`, `BrowseTREATY_IN_EDM`, `BrowseTreatyOutDetailEDM`, `BrowseOffer(convert)` |
| Pintu masuk | 3 | Harness + Section `InputTreatyInAdjustment`, Menu `MasterNavTreaty` |
| Lampiran | 3 | `TreatyRevisionCopyAttachment`, `CopyAllAttachment2_Sql`, `pyGetAllAttachments` |
| Lain-lain GISFW | 7 | `RefreshAchievement`, `SetCurrencyID`, `SetTreatyGroupID`, `Installments_ReadOnly` (×2), `IsTreatyUser`, `TreatyInFacultativeShareCalculationOldData` (Harness) |
| Aturan platform Pega, bukan bisnis | 3 | `ActivityStatusSuccess` (Pega-ProCom), `pyAttachmentScreen` (Pega-EndUserUI), `recordEvent` (Pega-UIEngine) |

### 1.4 Enam berkas yang hanya ada di ekspor Treaty In

`AttachmentLife`, `BrowseOffer`, `CategoryAttach_SQL`, `InputParamUploadReas_act`,
`SetCategoryAttach`, `setCategoryAttachment_DT` — seluruhnya seputar **kategori lampiran**.
Di sisi Adjustment padanannya adalah `BrowseOffer(convert)` dan `pyGetAllAttachments`. Ini
perbedaan **jalur lampiran**, bukan perbedaan model kontrak.

### 1.5 Satu kelas yang benar-benar khas Adjustment

`ASM-FW-GISFW-Int-treaty_in_edm` — kelas integrasi tabel `M_TREATY_IN_EDM`. Rinciannya di §2.

---

## 2. Model data: satu halaman clipboard, empat pohon

Ini temuan terpenting sesi ini, dan ia **membalik anggapan sesi sebelumnya**.

Sebuah baris addendum menyimpan seluruh halaman `TreatyIn` sebagai satu `JSONDATA`
(`RDBList/SaveTreatyInEDM.xml`, parameter pertama `{InputParam.DATAPEGA}`, diisi di
`SaveTreatyIn_EDM_Act` langkah 7 dengan `@ASM.GetPageJSONString()` atas step page `TreatyIn`).
Di dalam satu halaman itu hidup **empat pohon paralel**:

| Pohon | Isi | Ditulis oleh |
|---|---|---|
| `TreatyIn.*` | nilai **baru** yang sedang disunting | layar addendum |
| `TreatyIn.OLDDATA.*` | **potret penuh kontrak lama**, dibekukan saat addendum dibuat | `TreatyInSetAddendumToHistory` |
| `TreatyIn.ActualValue.*` | salinan `TreatyIn` tanpa `OLDDATA` dan tanpa `ValueDifference`, dibuat **setiap kali disimpan** | `SaveTreatyIn_EDM_Act` langkah 2-4 |
| `TreatyIn.ValueDifference.*` | selisih baru minus lama, **165 jalur daun di bawah 33 akar** (§5.6) | `TreatyEDMDifference*` |

### 2.1 `OLDDATA` adalah potret, bukan rujukan - dan itu mengoreksi G2

`Activity/TreatyInSetAddendumToHistory.xml` seluruhnya tiga langkah:

```
1  Page-Copy   TreatyIn      -> TreatyIntemp
2  Page-Copy   TreatyIntemp  -> TreatyIn.OLDDATA
3  Page-Remove TreatyIntemp
```

Dipanggil sebagai langkah **pertama** `TreatyInRevisi_post`, yaitu sebelum pengenal addendum
dibentuk dan sebelum apa pun disunting. Jadi `OLDDATA` memuat kontrak **sebagaimana adanya pada
detik addendum dibuat**, dan ia ikut tersimpan ke `JSONDATA` addendum itu.

> **Koreksi atas `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` §G2.**
> Sesi sebelumnya menyimpulkan *"nilai selisih historis tidak dapat direproduksi"* karena `OLDID`
> menunjuk baris `treaty_in` yang disunting di tempat. Bagian pertamanya benar - `OLDID` memang
> penunjuk ke sasaran bergerak - tetapi **kesimpulannya tidak berlaku**, karena selisih tidak
> pernah dibaca lewat `OLDID`. Seluruh mesin selisih membaca `TreatyIn.OLDDATA`, yang **beku di
> dalam baris addendum itu sendiri**.
>
> Yang benar: **nilai lama sebuah addendum dapat direproduksi hari ini, dari `JSONDATA`-nya
> sendiri.** `OLDID` dipakai untuk hal lain - penyalinan lampiran (`CopyAllAttachment2_Sql`
> menyeleksi `where treatyid = {TreatyIn.OLDID}`) dan penelusuran rantai.
>
> Bukti berkas: `Activity/TreatyInSetAddendumToHistory.xml`, `Activity/TreatyEDMDifference*.xml`
> (seluruh ruas kanannya `TreatyIn.OLDDATA....`), `Activity/SaveTreatyIn_EDM_Act.xml` langkah 7.

Akibatnya untuk migrasi berlawanan arah dengan dugaan sebelumnya: nilai lama **tidak hilang**, tapi
ia **berganda**. Tiap addendum membawa salinan penuh kontrak, dan tiap penyimpanan menambah satu
salinan lagi ke `ActualValue`.

### 2.2 `ActualValue` ditimpa saat menyimpan, kecuali pada addendum premi

`SaveTreatyIn_EDM_Act` langkah 2-4:

```
2  Page-Copy   TreatyIn                  -> TreatyInTemp     [precondition]
3  Page-Remove TreatyInTemp.OLDDATA, TreatyInTemp.ValueDifference
4  Page-Copy   TreatyInTemp              -> TreatyIn.ActualValue
```

Precondition langkah 2, dibaca dari `pyStepsPreCondParamsWhen`:

```
WHEN TreatyIn.EDMState == "3"  ->  benar: LOMPAT ke blok "jmp" (langkah 5)
                                   salah: lanjut
```

Jadi penyalinan ke `ActualValue` berjalan untuk `EDMState` **1 dan 2**, dan **dilewati** untuk
`EDMState` 3. Itu masuk akal: pada addendum premi, `ActualValue` justru pohon yang disunting
pengguna (disemai `TreatyInSetEditPre` langkah 3.5, `ActualValue.EGNPI = EGNPI`), sehingga menimpanya
akan menghapus masukannya.

Deskripsi langkahnya berbunyi *"when edmtype= 2 copy data actual to treatyin.actualvalue"*.
Kondisinya bukan `= 2` melainkan `!= 3`, sehingga langkah itu **juga** berjalan untuk `EDMState` 1.
Labelnya kurang, kondisinya yang menang - tetapi arah dasarnya benar dan ini **bukan cacat**.

### 2.3 Kelas khas addendum

`ASM-FW-GISFW-Int-treaty_in_edm` - kelas integrasi ke `M_TREATY_IN_EDM` / `TREATY_IN_EDM`.
Kolom kepala yang diserahkan ke prosedur `POOLDATA.PEGA_M_TREATY_IN_EDM` (25 masukan, 3 keluaran):

```
JSONDATA(DATAPEGA), ID, OLDID, STSCARI1(datacount), ProportionType, TreatyContractName,
TeritorialScope, Commencement, Termination, ClassofBusiness, LeadingReinsSource,
LeadingReinsSourceID, Ceding, CedingID, LeadingReinsID, NusareSharePct, BrokeragePct,
Information, PositionUsername, Position, StatusAkseptasi, ChooseStatusAkseptasi, TreatyYear,
EDMState, EDMMaterialType
  -> ERRMSG out, IDPEGAOUT out, STSSAVE out
```

Tiga di antaranya hanya ada di jalur addendum: **`OLDID`**, **`EDMState`**, **`EDMMaterialType`**.
Sisanya lapisan kepala yang sama dengan kontrak.

---

## 3. Pintu masuk dan susunan layar

### 3.1 Satu layar, tiga kolom nilai

Modul ini punya satu harness dan satu section bernama sama: **`InputTreatyInAdjustment`**
(GISFW 01-01-56 / 01-01-85). Dari sanalah seluruh modul digantung:

```
InputTreatyInAdjustment  (Harness + Section, Data-Portal)
├── TreatyInNONProportional              -> NILAI BARU  (aturan bersama dengan Treaty In)
├── TreatyInNONProportionalOldData       -> NILAI LAMA  (khas Adjustment)
│   ├── TreatyInTabsProportionalOldData
│   │   ├── LimitProportionalOldData -> DetailLimitsOldData
│   │   └── TotalLimitsOldData
│   └── TreatyInTabsNonProportionalOldData
│       ├── LayersOldData -> CoBListOldData
│       └── TreatyInTabsNonProportionalOldDataShare -> ShareOldData -> DetailShareOldData
└── TreatyInTabsNPValueDifferenceProRate / _NoProRate  -> SELISIH (aturan bersama)
```

Jadi layarnya membandingkan tiga pohon yang sudah diuraikan di §2: `OLDDATA`, `TreatyIn`, dan
`ValueDifference`. **Seluruh 20 berkas ber-akhiran `OldData` adalah cermin baca dari `OLDDATA`** —
tidak satu pun membawa logika hitung; yang ada hanya `pyDisabledWhen` dan bentuk tampilan.

### 3.2 Dua tombol yang menetapkan jenis addendum, dan dua picker

Pada `InputTreatyInAdjustment` ada dua tombol yang memanggil data transform `TreatyCreateEDM`:

| Tombol | Parameter | Akibat |
|---|---|---|
| Revisi | `type = revision` | `TreatyIn.EDMState = "1"` |
| Penyesuaian | `type = adjustment` | `TreatyIn.EDMState = "3"` |

`DataTransform/TreatyCreateEDM.xml` seluruhnya empat baris: `SET TreatyIn = ""`, lalu dua cabang
`WHEN` itu. **Ia tidak pernah menetapkan `EDMState = "2"`** — lihat §4.2.

Sesudah itu salah satu picker terbuka, dan pickernya sendiri yang menentukan parameter pembuatan:

| Picker | Tampil bila | Memanggil | Parameter |
|---|---|---|---|
| `PickerTreatyInMaster` | `TreatyIn.EDMState = '3'` | `TreatyInEDMSetValue` | `ID=.CARI1`, **`InternalType=3`**, **`MaterialType=1`** (keduanya tetap) |
| `PickerTreatyInMasterRevisi` | `TreatyIn.EDMState = '1'` atau `'2'` | `TreatyInEDMSetValue` | `ID=.CARI1`, `InternalType=TreatyIn.EDMState`, `MaterialType=TreatyIn.EDMMaterialType` |

Artinya penyesuaian premi **selalu** material; hanya jalur revisi yang boleh memilih sifat
materialnya.

### KOREKSI 23 September 2026 (grilling, putaran 1) — jenis dan materialitas DIPILIH DI PICKER

Ditemukan saat grilling cabang A, dan ia mengubah §4.2, §4.3, dan TDA-13.

**Yang terbaca.** `PickerTreatyInMasterRevisi` memuat **dua radio group yang dapat disunting**,
terikat langsung ke properti penanda jenis dan materialitas:

| Kontrol | Terikat ke | `pyReadOnly` | `pyDisabledNew` | Kondisi tampil |
|---|---|---|---|---|
| `pxRadioButtons` | `TreatyIn.EDMState` | `false` | `false` | **tidak ada** |
| `pxRadioButtons` | `TreatyIn.EDMMaterialType` | `false` | `false` | **tidak ada** |

Kondisi tampil diperiksa sampai ke akar seksi — sel, baris, tabel, `Embed-Harness-SectionBody`,
dan `Embed-Harness-Section` di atasnya: **tidak satu pun membawa `pyVisible`, `pyCondition`,
maupun `pyVisibleWhen`.** Jadi kedua radio itu tampil dan dapat disunting setiap kali picker Revisi
dirender.

`PickerTreatyInMaster` (jalur penyesuaian premi) **tidak memuat keduanya** — sejalan dengan
parameternya yang tetap (`InternalType=3`, `MaterialType=1`).

**Akibatnya** `TreatyCreateEDM` hanya **menyemai** nilai awal; yang benar-benar terkirim ke
`TreatyInEDMSetValue` adalah nilai radio pada saat baris dipilih.

**Kadar klaim — dibatasi, dan batasnya penting.** Kedua radio memakai `pyListSource = associated`,
sehingga himpunan nilai yang ditawarkannya berasal dari aturan **Field Value** properti itu, dan
aturan itu **tidak ikut ter-ekspor**. Maka:

| Pernyataan | Kadar |
|---|---|
| mekanisme untuk mengubah jenis dan materialitas di layar **ada, tampil, dan tidak dikunci** | **terbaca** |
| `EDMState = "2"` dapat dibuat lewat jalur ter-ekspor | **masih bisa terjadi — bila Field Value-nya menawarkan `2`**; belum terbukti |
| radio Revisi juga menawarkan `3` | **tidak terbaca** |
| penyesuaian premi selalu material | **terbaca hanya untuk jalur tombol Penyesuaian**; untuk jalur radio, tidak terbaca |

Dibutuhkan ekspor tambahan: aturan Field Value untuk `EDMState` dan `EDMMaterialType`
(lihat `KEPUTUSAN-GRILLING-ADJUSTMENT.md`, daftar ekspor tambahan).

### 3.3 Daftar yang dipilih menyatukan kontrak dan addendum

`RDBList/TreatyLoadMasterJoinEdm.xml`:

```sql
select a.ID as CARI1, a.TREATYCONTRACTNAME as CARI2, a.PROPORTIONTYPE as CARI3,
       a.LEADINGREINSSOURCE as CARI4, a.CEDING as CARI5,
       a.COMMENCEMENT as CARI6, a.TERMINATION as CARI7
from pooldata.treaty_in a
UNION
select b.ID … from pooldata.treaty_in_edm b
```

Tujuh kolom yang sama dari dua tabel, **tanpa kolom pembeda jenis**, dan **tanpa penyaringan
keadaan**: kontrak yang belum disetujui dan addendum yang sudah kedaluwarsa sama-sama muncul.

Variannya `TreatyLoadMasterJoinEdmXOL` menambahkan `WHERE PROPORTIONTYPE = 'NonProportional'` pada
kedua sisi. Varian yang non-XOL **tidak** menyaring sebaliknya — daftar "proporsional" tetap memuat
baris non-proporsional. Ini terbaca dari kedua SQL berdampingan, bukan disimpulkan.

### 3.4 Siapa yang boleh masuk

`When/IsTreatyUser.xml` — logika `A OR B OR C OR D OR F OR E`:

```
OperatorID.pyWorkBasketList(2).pyWorkBasketName = "ReasTreatyInAdmin"        (A)
                                               = "ReasTreatyInGroupLeader"   (B)
                                               = "ReasTreatyInDirector"      (C)
                                               = "ReasTreatyInDeptHead"      (D)
                                               = "ReasTreatyInSecHead"       (F)
OperatorID.pyPosition = "IT Developer"                                       (E)
```

Tiga hal terbaca langsung: peran disimpan sebagai **nama workbasket**, diambil dari **indeks
kedua yang tetap** (`pyWorkBasketList(2)` — bukan pencarian di seluruh daftar), dan ada **jalur
lolos untuk "IT Developer"**.

`When/TreatyMasterInEDM.xml` — `EDMState = "1" OR "2" OR "3"`; inilah satu-satunya penanda
"halaman ini addendum, bukan kontrak".

---

## 4. Daur hidup addendum

### 4.1 Urutan pembuatan, langkah demi langkah

`Activity/TreatyInEDMSetValue.xml` (param: `ID`, `InternalType`, `type`, `MaterialType`):

| # | Langkah | Precondition sebenarnya |
|---|---|---|
| 1 | `Call SetTreatyInEDM_Act` — muat baris **addendum** dari `M_TREATY_IN_EDM` | `@length(Param.ID) == 7` benar -> LEWATI |
| 2 | `Call SetTreatyIn_Act` — muat baris **kontrak** dari `TREATY_IN` | `@length(Param.ID) == 7` benar -> lanjut |
| 3 | `Property-Set`: `OLDID <- ID`, `CommentList <- ""`, `EDMState <- Param.InternalType`, `EDMMaterialType <- Param.MaterialType`, `RevisionDate <- @CurrentDateTime()` | — |
| 4 | `Apply-DataTransform TreatyInSetEditPre` | — |
| 5 | `Apply-DataTransform TreatyInSetEdit` | — |
| 6 | `Call TreatyInRevisi_post` — potret `OLDDATA` + pembentukan pengenal | — |
| 7 | `Call TreatyInEdmCheckDuplicate` | **MATI (blok `//`)** |
| 8 | `Call TreatyRevisionCopyAttachment` | — |
| 9 | `Call SaveTreatyIn_EDM_Act` | — |

**Panjang pengenal adalah pembeda jenis baris.** Tujuh karakter berarti kontrak, selain itu
addendum. Tidak ada kolom jenis; bentuk teksnya sendiri yang dibaca.

### 4.2 `EDMState` — tiga nilai; aturan menulis `1` dan `3`, sedangkan `2` hanya lahir di layar

`DataTransform/TreatyInSetEditPre.xml` memberi nama ketiganya lewat teks komentar yang ia
tambahkan ke `CommentList`:

| `EDMState` | Komentar yang ditulis | Tambahan |
|---|---|---|
| `1` | *"Had Created Internal Edit"* | `EDMEffective <- Commencement` |
| `2` | *"Had Created External Addendum"* | — |
| `3` | *"Had Created Addendum Premium"* | `ActualValue.EGNPI <- EGNPI`, `AddendumPremi <- "1"` |

**Tidak ada *aturan* yang menetapkan `EDMState = "2"`** — `TreatyCreateEDM` hanya menghasilkan 1
dan 3. Tetapi nilainya **tidak hanya berasal dari aturan**: `PickerTreatyInMasterRevisi` memuat
radio group yang dapat disunting dan terikat langsung ke `TreatyIn.EDMState` (§3.2), dan picker itu
tampil untuk `EDMState = '1'` **maupun** `'2'`. Jalurnya karena itu terbaca utuh:

```
tombol Revisi -> TreatyCreateEDM(type=revision) -> EDMState = "1"
              -> PickerTreatyInMasterRevisi tampil (syaratnya '1' atau '2')
              -> pengguna memindahkan radio ke "2"          <- di sinilah 2 lahir
              -> pilih baris -> TreatyInEDMSetValue(InternalType = TreatyIn.EDMState = "2")
```

> **KOREKSI 26 September 2026 — `EXP-1` datang, dan bentuknya bukan Field Value.**
>
> Kedua aturan ada di `ekspor-tambahan/`, dan **bukan** Field Value di kelas `Data-Portal` seperti
> dugaan saya. Keduanya **`Rule-Obj-Property`** di kelas **`ASM-FW-GISFW-Int-TREATY_IN`**, ruleset
> `GISFW` **01-01-56**, dengan `pyTableOption = PromptList`:
>
> | Aturan | Daftar pilihan (badan) | `pyEditValidate` | `pyMemo` | Dibuat |
> |---|---|---|---|---|
> | `EDMState` | `1` = Internal, `2` = External | **kosong** | *"display only not for validation"* | 16 Jun 2020 |
> | `EDMMaterialType` | `1` = Material, `2` = Non Material | **kosong** | *"added prompt list"* | 26 Okt 2020 |
>
> Nilai `3` muncul di kedua berkas **hanya di dalam `pyRuleVersionsList`** — salinan versi lama,
> yang menurut aturan bukti ronde ini (`00-LINGKUP` §5a) tidak dihitung dan tidak jadi bukti.
>
> **Tiga pemeriksaan sebelum daftar ini dipakai, dan ketiganya dijalankan:**
>
> 1. **Sumber daftar radio.** Kedua sel radio di `PickerTreatyInMasterRevisi` membawa
>    `pyListSource = associated`, `pyListLoadMode = auto`, `pyIsValidDataPage = false`, dan
>    **tidak memuat satu pun daftar lokal** — nol `pyPromptTableList`, nol `pyOptionsList`, nol
>    data page. Jadi daftarnya memang datang dari definisi propertinya. Satu tanda yang mudah
>    disalahbaca dan karena itu disebut: `pySmartPromptClass = Data-Portal` pada sel itu adalah
>    kelas **bantuan pengetikan saat merancang**, bukan sumber daftar saat berjalan.
> 2. **Resolusi kelas.** Halaman `TreatyIn` dideklarasikan berkelas `ASM-FW-GISFW-Int-TREATY_IN`
>    pada **190 kemunculan `pyPagesAndClasses` di badan**, dan **tidak ada kelas lain**. Kelas
>    `ASM-FW-GISFW-Int-treaty_in_edm` memang ada (9 kemunculan), tetapi sebagai kelas *applies-to*
>    aturan RDB List yang menghadap basis data — bukan kelas halaman formulir. Jadi properti inilah
>    yang terpakai.
> 3. **Versi ruleset — dan di sinilah kadarnya dibatasi.** Ekspor ini memuat **01-01-56**.
>    Di badan ekspor Adjustment, versi ruleset tersebar dari 01-01-53 sampai **01-01-95**, dan
>    01-01-56 kebetulan yang paling sering (103 kemunculan). **Tidak terbaca** apakah ada versi
>    properti yang lebih tinggi di sistem yang berjalan. Maka pernyataan di bawah berlaku **untuk
>    versi 01-01-56**; bila ada versi lebih tinggi yang tidak ikut ter-ekspor, daftarnya dapat
>    berbeda. **Perlu dicek** ke pemilik sistem.
>
> **Akibatnya, dan ketiganya menaikkan kadar:**
>
> * **`EDMState = 2` ("External") DITAWARKAN layar.** TDA-13 naik dari *masih bisa terjadi*
>   menjadi **terbaca**. `EXP-1` ditutup.
> * **Radio tidak menawarkan `3`.** Maka "penyesuaian premi selalu material" **terbaca untuk kedua
>   jalur**, bukan hanya jalur tombol Penyesuaian — sebab jenis 3 tidak dapat dipilih dari picker.
> * **Kombinasi yang dapat dibuat lewat layar hanya lima:** (1,1), (1,2), (2,1), (2,2) dari radio,
>   ditambah (3,1) dari tombol Penyesuaian. Ditulis di kepala `UA-2`.
> * **TDA-10 menguat.** `pyEditValidate` **kosong** pada kedua properti, dan memo pengembangnya
>   sendiri berbunyi *"display only not for validation"*. Daftar itu tampilan, bukan penjaga —
>   persis bentuk TDA-10.
> * **Bahan C1, sebagai fakta sejarah dan bukan bukti niat:** sumbu materialitas lahir **empat
>   bulan sesudah** sumbu jenis, sebagai properti tersendiri.

**Jadi jenis 2 dapat dibuat lewat antarmuka yang ter-ekspor**, dan sejak 26 September 2026 itu
**terbaca**, bukan lagi dugaan. Sebarannya di produksi dihitung **UA-2**.

### 4.3 `EDMMaterialType` — dua nilai, penentu apa yang boleh disunting

Hanya `1` dan `2` yang muncul di seluruh ekspor.

> **KOREKSI 24 September 2026 (audit `TA-04`).** Rumusan lama berbunyi *"satu-satunya aturan yang
> menulisnya"*. Disapu ulang dengan perkakas terkalibrasi atas **lima jenis penulis**, penulisnya
> **tiga**, bukan satu — dan yang dua adalah kontrol layar, bukan aturan:
>
> | Jenis | Penulis | Catatan |
> |---|---|---|
> | Property-Set | `TreatyInEDMSetValue` langkah 3 | dari `Param.MaterialType`, yaitu parameter picker |
> | pengikatan kontrol | `PickerTreatyInMasterRevisi`, `pxRadioButtons` | **dapat disunting** |
> | pengikatan kontrol | `InputTreatyInAdjustment`, `pxDropdown` | **hanya-baca** - tidak menulis |
>
> Kesimpulannya tidak berubah: satu-satunya asal nilai adalah radio group di picker. Yang berubah
> adalah kalimatnya - "satu-satunya aturan" benar, "satu-satunya penulis" tidak.

Untuk **jalur tombol Penyesuaian** nilainya tetap `1`, tanpa kontrol apa pun: di situ penyesuaian
premi memang selalu material.

> **KOREKSI 26 September 2026 (`EXP-1`).** Pertanyaan *"apakah jenis penyesuaian premi dapat dicapai
> lewat radio"* kini **terjawab: tidak.** Daftar `PromptList` properti `EDMState` hanya memuat `1`
> dan `2`; nilai `3` hanya ada di dalam `pyRuleVersionsList` dan tidak dihitung. Maka jenis 3 tidak
> dapat dipasangkan dengan materialitas `2` lewat layar, dan **"penyesuaian premi selalu material"
> berlaku untuk kedua jalur** — bukan hanya jalur tombol.

Nilainya dibaca oleh **ratusan kondisi `pyDisabledWhen`** di layar penyuntingan.

> **KOREKSI 24 September 2026 (grilling, putaran 5) — hitungan lama ikut menghitung salinan.**
> Tabel sebelumnya menghitung seluruh kemunculan dalam berkas, termasuk yang berada di dalam
> `pyIncludedRuleXML` dan `pyRuleVersionsList` — yaitu **salinan aturan lain yang ikut terbungkus
> saat ekspor**, bukan isi aturan itu sendiri. Tiga baris teratasnya karena itu keliru, dan yang
> paling menyesatkan: `Section/InputTreatyInOffer.xml` dilaporkan 70, padahal **badan aturannya
> memuat nol** — seluruh 76 kemunculannya salinan.
>
> Hitungan di bawah **hanya badan aturan**.
>
> **KOREKSI 24 September 2026 (`MA-13`).** Baris ini semula berbunyi *"**220 kondisi** tersebar di
> **18 seksi**, ditambah 215 kemunculan salinan yang tidak dihitung"*. Angka itu **tidak dapat
> direproduksi** dan **dicabut** — bukan diganti angka baru di sini. Yang berlaku:
> kondisi penguncian yang menyebut `EDMMaterialType` (hitungannya `GRILL-D/01-TEMUAN.md` `TD-02`).

| Berkas (badan aturan) | Kondisi | | Berkas | Kondisi |
|---|---|---|---|---|
| `Section/DetailLimits.xml` | 60 | | `Section/DetailEGNPI.xml` | 4 |
| `Section/Layers.xml` | 38 | | `Section/CoBList.xml` | 4 |
| `Section/TreatyInTabsNonProportional.xml` | 32 | | `Section/TreatyInTabsNPValueDifference_NoProRate.xml` | 3 |
| `Section/TreatyInTabsProportional.xml` | 31 | | `Section/TreatyInTabsNPValueDifferenceProRate.xml` | 3 |
| `Section/TreatyInNONProportional.xml` | 13 | | `Section/LayersOldData.xml` | 2 |
| `Section/Share.xml` | 8 | | `Section/LayersEDM.xml` | 2 |
| `Section/ShareRetro.xml` | 7 | | `Section/DetailShareOldData.xml` | 2 |
| `Section/MaxRetention.xml` | 4 | | `Section/DetailShare.xml` | 2 |
| `Section/Installments.xml` | 4 | | `Section/DetailShareRetro.xml` | 1 |
| | | | **Total badan** | **220** |
| | | | `Section`/`Harness` `InputTreatyInOffer` | **0 badan**, 76 salinan masing-masing |

Bentuknya `pyDisabledWhen = TreatyIn.EDMMaterialType = 2`, sering digabung
`TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2`. Jadi **`2` mematikan field nilai** dan
`1` membukanya. Satu kondisi menyebut keduanya sekaligus:
`pyContainerVisibleWhen = TreatyIn.EDMState != 3 && TreatyIn.EDMMaterialType == 1`.

Kesimpulan yang dapat ditarik tanpa bertanya ke bisnis: **material** berarti nilai uangnya boleh
berubah, **non-material** berarti hanya keterangan yang boleh berubah. Sifat ini **tidak** ditegakkan
di penyimpanan — ia hanya kondisi tampilan. Tidak ada pemeriksaan di sisi penyimpanan yang menolak
perubahan nilai pada addendum non-material.

### 4.4 Pembentukan pengenal: `‹kontrak›/Rnn`, dihitung dari baris yang DIPILIH

`Activity/TreatyInRevisi_post.xml`:

| # | Langkah | Precondition |
|---|---|---|
| 1 | `Call TreatyInSetAddendumToHistory` (potret `OLDDATA`) | — |
| 2 | `InputData.CARI1 <- TreatyIn.ID + "%"` | — |
| 3 | `RDB-List GetTreatyRevisionID` | — |
| 4 | `TreatyIn.ID <- TreatyIn.ID + "/R01"` | `OutputData.pxResults(1).HASIL1 = ""` benar -> lanjut |
| 5 | `TreatyIn.OLDID <- TreatyIn.ID` lalu `TreatyIn.ID <- @If(...)` | `… = ""` benar -> LEWATI |

Ekspresi langkah 5, utuh:

```
@If(@substring(TreatyIn.ID,10,12) < 10,
    @substring(TreatyIn.ID,0,7) + "/R0" + (@toInt(@substring(TreatyIn.ID,10,12)) + 1),
    @substring(TreatyIn.ID,0,7) + "/R"  + (@toInt(@substring(TreatyIn.ID,10,12)) + 1))
```

Dan kuerinya, `RDBList/GetTreatyRevisionID.xml`:

```sql
select TO_NUMBER(SUBSTR(ID,10,2))+1 as HASIL1
from m_treaty_in_edm
where ID like {InputData.CARI1}
and ROWNUM = 1
order by ID desc
```

Tiga hal terbaca dari kedua potongan itu berdampingan:

1. **Nilai `HASIL1` tidak pernah dipakai.** Ia hanya diuji kosong/tidak-kosong. Nomor revisinya
   dihitung ulang di langkah 5 dari `TreatyIn.ID` sendiri. Karena itu cacat `ROWNUM = 1` yang
   mendahului `ORDER BY` — yang sudah diparkir sebagai butir A-2 — **tidak berakibat pada nomor
   yang dihasilkan**: baris mana pun yang kembali, hasil ujinya sama. Butir A-2 boleh diturunkan
   bobotnya.
2. **Nomor berikutnya dihitung dari baris yang DIPILIH pengguna, bukan dari revisi terakhir.**
   Memilih `1234567/R01` ketika `1234567/R02` sudah ada menghasilkan `1234567/R02` lagi.
3. **Penjaga duplikatnya mati.** `TreatyInEDMSetValue` langkah 7 (`TreatyInEdmCheckDuplicate`)
   ber-blok `//`.

Apa yang terjadi pada tabrakan itu terbaca dari prosedurnya — lihat §6.2. Ini **"masih bisa
terjadi"**; **"pernah terjadi"** hanya dapat dijawab data.

Satu hal yang **tidak terbaca** dan tidak saya tebak: perilaku `@substring` ketika posisi awalnya
melampaui panjang teks — yaitu ketika yang dipilih adalah kontrak tujuh karakter yang sudah punya
revisi, sehingga langkah 5-lah yang berjalan. Ekspor tidak memuat jawabannya.

#### Akibat konkret TDA-12: penomorannya patah pada revisi kesepuluh — dengan dua kemungkinan akhir

Klaim yang beredar bahwa penomorannya "hanya sampai R99" **tidak benar**. Ia patah jauh lebih awal.

**Sebabnya satu, bukan dua.** Cacatnya sepenuhnya ada pada ekspresi
`Activity/TreatyInRevisi_post.xml` **baris 848**, yang membaca mulai **indeks 10** — digit satuan
saja — pada pengenal yang digit puluhannya ada di indeks 9:

```
@If(@substring(TreatyIn.ID,10,12) < 10, …+"/R0"+(n+1), …+"/R"+(n+1))
```

Kueri `RDBList/GetTreatyRevisionID.xml` memakai `SUBSTR(ID,10,2)` yang — karena Oracle 1-based —
justru membaca **kedua digit dengan benar**. Tetapi nilainya **dibuang**: `HASIL1` hanya diuji
kosong atau tidak (§4.4 butir 1). Jadi pengurai yang benar tidak pernah memengaruhi nomor yang
dihasilkan, dan menyebut "dua pengurai" sebagai sebab adalah salah — saya cabut.

**Kadarnya dibatasi satu hal yang tidak terbaca.** Pengenal `1234567/Rnn` panjangnya **11
karakter**, sehingga `@substring(ID,10,12)` selalu meminta batas akhir **di luar** panjang teks.
Perilaku `@substring` dalam keadaan itu tidak terbaca dari ekspor — hal yang sudah ditandai di
akhir §4.4. Akibatnya bercabang dua, dan arah dampaknya berbeda (`METODE` §2.2):

| Bila `@substring` … | Yang terjadi | Sifat kerusakan |
|---|---|---|
| **memotong dengan longgar** (mengembalikan sisa teks) | `/R09` -> `/R010` (12 karakter, bentuk rusak) -> `/R11` (melompati `R10`) -> **`/R02`**, menimpa revisi kedua lewat cabang `UPDATE` yang melapor berhasil (§6.2) | **diam** |
| **gagal keras** | membuat revisi **dari sebuah revisi** selalu galat; hanya revisi pertama yang pernah berhasil | **berisik** |

Keduanya rusak. Yang memutuskan mana yang terjadi adalah **data, bukan kode** — **UA-1**, yang
karena itu diperluas menjadi tiga kueri. Satu saja baris yang `OLDID`-nya menunjuk baris
`TREATY_IN_EDM` lain sudah **menyingkirkan** kemungkinan gagal keras.

### 4.5 KOREKSI 24 September 2026 — membuka kontrak untuk disesuaikan TIDAK mengubah kontraknya

**Versi pertama bagian ini salah, dan dicabut seluruhnya.** Ia menyatakan bahwa memilih sebuah
kontrak di picker Adjustment menulis kembali ke baris kontraknya. Itu keliru: saya membaca daftar
langkah `SetTreatyIn_Act` tanpa memeriksa preconditionnya, dengan perkakas yang saat itu belum
membaca `pyStepsPreCondParams` yang tersarang.

**Yang benar.** Keenam langkah yang menulis ke baris kontrak **bersyarat pada sebuah parameter**,
dan syaratnya **aktif** — `pyStepsPreCondition = 'true'` pada keenamnya, tanpa blok mati:

| # | Baris | Syarat | `pyStepsPreCondition` | Bila benar | Bila salah | Blok |
|---|---|---|---|---|---|---|
| 6 | 1188 | `param.viewstate==1` | `true` — **aktif** | lanjut (kosong) | **3 = LEWATI** | — |
| 7 | 1392 | `param.revisionstate==1` | `true` — **aktif** | `2` = lanjut | **3 = LEWATI** | — |
| 8 | 1580 | `param.revisionstate==1` | `true` — **aktif** | `2` = lanjut | **3 = LEWATI** | — |
| 9 | 1767 | `param.revisionstate==1` | `true` — **aktif** | `2` = lanjut | **3 = LEWATI** | — |
| 10 | 1933 | `param.revisionstate==1` | `true` — **aktif** | `2` = lanjut | **3 = LEWATI** | — |
| 11 | 2117 | `param.revisionstate==1` | `true` — **aktif** | `2` = lanjut | **3 = LEWATI** | — |

Nomor baris merujuk `Activity/SetTreatyIn_Act.xml` pada ekspor Adjustment. Langkah 6 disertakan
karena ia memakai parameter yang **berbeda** (`viewstate`), sehingga sebagian pemanggil menyalakan
6 tanpa menyalakan 7-11.

**Penyisiran pemanggil yang tertutup.** Setiap kemunculan nama `SetTreatyIn_Act` di **kedua**
ekspor dihitung dan digolongkan menurut tag pembawanya, sehingga tidak ada jalur yang lolos
(`METODE` §2.0a, §3.1):

| Tag pembawa | Treaty In | Adjustment | Artinya |
|---|---|---|---|
| `pyActivity` | 8 badan + 8 salinan | 8 badan + 8 salinan | **action set di layar** — **delapan kontrol, di satu tempat**: `Section/InputTreatyInOffer.xml` (badan 8, salinan 0). Kedelapan kemunculan di harness-nya **seluruhnya di dalam `pyIncludedRuleXML`** (badan 0, salinan 8) — salinan susunan yang sama, bukan kontrol tambahan |
| `pyStepsActivityName` | 2 | 3 | **langkah `Call` di Activity** |
| `pxStepDefaultDescription` | 2 | 3 | teks bawaan langkah yang sama — bukan pemanggilan tersendiri |
| `pyRuleName`, `pyLabel`, `pyLabelOld`, `pxTabLabel`, `pyActivityName` | 7 | 8 | metadata aturan itu sendiri |

**Nol** kemunculan pada pra/pasca-proses Flow Action, data page, atau jalur lain. Jadi pemanggilnya
hanya dua jenis: langkah `Call` di tiga aktivitas, dan **delapan kontrol di satu seksi**.

> **Koreksi 24 September 2026:** angka `16/16` pada versi pertama tabel ini ikut menghitung salinan
> terbungkus. Kesimpulannya tidak berubah — hanya kontrol 7 dan 8 yang mengirim `revisionstate=1` —
> tetapi jumlah tempatnya **satu**, bukan dua.

| Pemanggil | Ada di | Parameter | Langkah 7-11 |
|---|---|---|---|
| `TreatyInEDMSetValue` langkah 2 — **jalur Adjustment** | hanya Adjustment | **tidak ada** | **DILEWATI** |
| `TreatyInDownloadAll` langkah 1 | kedua ekspor | tidak ada | DILEWATI |
| `TreatyInSaveROL` langkah 5 | kedua ekspor | `ID`, `viewstate=1` | DILEWATI |
| `InputTreatyInOffer` kontrol 1-6 | kedua ekspor | `ID`, sebagian `viewstate=1` | DILEWATI |
| **`InputTreatyInOffer` kontrol 7 dan 8** | kedua ekspor | `ID=.ID`, `viewstate=1`, **`revisionstate=1`** | **BERJALAN** |

**Apakah kontrol 7 dan 8 dapat dicapai dari layar addendum?** Tidak.
`Section/InputTreatyInAdjustment.xml` dan harness-nya **tidak menyertakan** `InputTreatyInOffer`
sama sekali — diperiksa dengan menyapu seluruh rujukan seksi. Picker Adjustment menargetkan
`InputTreatyInAdjustment` (`pySection` pada action set-nya), dan addendum dimuat lewat
`SetTreatyInEDM_Act`. Kedua kontrol itu tinggal di layar penawaran kontrak.

**Sifat kedua kontrol itu sendiri**, diperiksa sampai akar seksi: keduanya `pxButton`,
`pyReadOnly = false`, dan **nol** `pyVisible`, `pyCondition`, `pyVisibleWhen`, `pyDisabledWhen`,
maupun `pyPrivilege` pada sel, baris, tabel, `Embed-Harness-SectionBody`, dan
`Embed-Harness-Section` di atasnya. Parameter `ID` yang dikirim adalah `.ID`, yaitu pengenal baris
yang sedang dibuka di layar itu. Label tombolnya **tidak terbaca** dari ekspor.

> **Jadi penulisan ke baris kontrak hanya terjadi lewat dua kontrol di layar `InputTreatyInOffer`,
> yaitu layar Treaty In — bukan lewat jalur Adjustment sama sekali.** Kedua kontrol itu identik di
> kedua ekspor. Ia tampak sebagai tindakan "mulai revisi" yang disengaja, bukan efek samping
> membuka sesuatu.

**Apa yang tersisa sebagai temuan, dan apa yang tidak.**

* **Tidak tersisa:** "membuka kontrak mengubahnya". Tidak terjadi.
* **Tersisa, dan belum diadili:** dua kontrol itu **mengosongkan `StatusAkseptasi`** sebuah kontrak
  yang mungkin sudah disetujui, lalu menyimpannya — tanpa versi baru, tanpa jejak selain satu
  komentar. Itu temuan **Treaty In**, bukan Adjustment, dan berkasnya ada di irisan.
* **Tidak terbaca:** nilai `RevisionState` yang dibawa sebuah addendum. Jalur Adjustment tidak
  pernah menuliskannya, dan `TreatyInSetEdit` tidak mengosongkannya — jadi nilainya **diwarisi dari
  baris yang dimuat**. Berapa banyak baris membawa `1` adalah pertanyaan data (**UA-9**).

### 4.6 Menyunting nama cedant tidak selalu memperbarui ID-nya

> **KOREKSI 24 September 2026 (audit `TA-04`).** Dua angka di bagian ini keliru, dan keduanya
> berasal dari penghitungan yang tidak terkalibrasi.
>
> * **Kontrol `Ceding` di `Section/InputTreatyInAdjustment` ada SATU, bukan dua.** Disapu atas
>   pengikatan sel (`Embed-Display-Table-Cell` ber-`pyValue`), badan aturan saja: satu sel,
>   `pyValue = .Ceding`, `pyReadOnly = false`.
> * **Penulis `TreatyIn.CedingID` ada EMPAT, bukan dua** - lihat daftar di bawah.
>
> Kesimpulan pokoknya **tidak berubah**: kontrol `Ceding` pada layar addendum tetap tidak menulis
> `CedingID`.

Kontrol `Ceding` yang dapat disunting **tidak sama perilakunya**:

| Kontrol | Berkas | Bentuk | Menulis `CedingID`? |
|---|---|---|---|
| autocomplete | `Section/TreatyInNONProportional.xml` | `pxAutoComplete` | **YA** — selnya memuat `pyPropertyTarget = TreatyIn.CedingID`, dan action set-nya memanggil `TreatyInCheckCedingBlacklist` |
| kontrol layar addendum | `Section/InputTreatyInAdjustment.xml` | tanpa `pyFormat`, terikat `.Ceding` | **TIDAK** — nol rujukan `CedingID` di seluruh berkas itu |

**Penulis `TreatyIn.CedingID` ada EMPAT.** Disapu 24 September 2026 dengan perkakas terkalibrasi
atas lima jenis penulis, badan aturan saja:

| # | Jenis | Penulis | Menulis `Ceding` juga? | Terjangkau dari layar addendum? |
|---|---|---|---|---|
| 1 | kontrol Section (`pyPropertyTarget`) | autocomplete `Ceding` di `TreatyInNONProportional` | **ya** - selnya terikat `Ceding` | tidak |
| 2 | Property-Set | `TreatyInMappingDataconvert` (`<- .InsuredID`, `<- DataClient.pxResults(1).CARI1`) | **ya** (`Ceding <- .InsuredName`) | tidak - jalur konversi data, dipanggil `TreatyInConvertCallData_act` |
| 3 | Data Transform | **`TreatyInSetReinsured` langkah 1.2** (`<- param.id`) | **ya** - langkah 1.1 menulis `Ceding <- param.name` | tidak - dirujuk hanya `TreatyInSearchReinsured` / `…SoB`, keduanya hanya dari `TreatyInNONProportional` |
| 4 | pengikatan kontrol | **`Section/ShowSummary.xml`, `pxTextInput` dapat disunting** | ruas namanya kontrol terpisah | tidak - `ShowSummary` hanya dirujuk `InputTreatyInOffer` |

Tiga dari empat menulis **nama dan pengenal bersama-sama**, jadi ketidaksinkronan tidak berasal dari
mereka. Yang keempat, `ShowSummary`, justru bentuk sebaliknya: di sana **pengenalnya** dapat disunting
sebagai teks bebas, terpisah dari namanya - dan itu di layar penawaran, bukan layar addendum.

**Yang tetap berdiri:** satu-satunya kontrol `Ceding` pada layar addendum tidak menulis `CedingID`
sama sekali.

#### Temuan irisan yang lahir dari sapuan ini: `ShowSummary` membuka kelima field lapisan beku sebagai teks bebas

`Section/ShowSummary.xml` memuat `pxTextInput` **yang dapat disunting** untuk **kelima** field
lapisan beku ADR-0040 sekaligus:

| Field | Kontrol di `ShowSummary` | Di layar addendum |
|---|---|---|
| `Ceding` | `pxTextInput` dapat disunting (3 sel) | satu kontrol, dapat disunting |
| `CedingID` | **`pxTextInput` dapat disunting** | tidak ada kontrol |
| `Commencement` | `pxTextInput` dapat disunting | `pxDateTime` dapat disunting |
| `Termination` | `pxTextInput` dapat disunting | `pxDateTime` dapat disunting |
| `ProportionType` | **`pxTextInput` dapat disunting** | `pxDropdown` dan `pxRadioButtons`, keduanya **hanya-baca** |

Dua di antaranya tidak dapat dicapai dari mana pun selain di sini. `ShowSummary` dirujuk **hanya**
oleh `InputTreatyInOffer` beserta harness dan flow action-nya - jadi ini **layar penawaran**, yaitu
**cacat Treaty In yang berjalan hari ini**, bukan cacat jalur addendum.

Yang membuatnya berbeda dari kontrol lain: `pxTextInput` **tidak memvalidasi apa pun**. `CedingID`
di sana bukan hasil memilih cedant dari daftar - ia teks bebas, terpisah dari namanya. Inilah bentuk
paling langsung dari ketidaksinkronan yang diukur `UA-10`, dan ia ada di sisi kontrak.

Dicatat sebagai **usulan untuk modul induk** (GRL-01 butir 5). Diperiksa lebih dulu (`MA-06`):
**tidak satu pun dokumen di `_migration-docs/treaty-in/` menyebut `ShowSummary`** - bukan
`PENGETAHUAN.md` induk, bukan daftar eskalasi, bukan ADR mana pun. Ini titik buta induk, bukan
pengulangan.

Kedua kontrol `Ceding` di `Section/InputTreatyInAdjustment.xml` **tidak termasuk keduanya**: nol
`pyPropertyTarget` ke `CedingID` di seluruh berkas itu, badan maupun salinan.

**Akibatnya, terbaca dari bentuk kontrolnya:** mengubah nama cedant lewat kedua kontrol di layar
addendum meninggalkan `CedingID` **pada nilai lamanya**. Nama dan pengenal cedant pada baris
addendum karena itu dapat **tidak sinkron**, dan tidak ada apa pun yang memeriksanya.

Ketiadaan kontrol untuk `CedingID` karena itu **bukan** bukti bahwa `CedingID` tidak pernah
menyimpang — ia justru bentuk penyimpangannya (`METODE` §2.0). Diukur **UA-10**.

**Akibatnya ditarik sampai habis** (`METODE` §4.5). Prosedur `PEGA_M_TREATY_IN_EDM` menerima
`{TreatyIn.Ceding}` **dan** `{TreatyIn.CedingID}` sebagai dua parameter terpisah (§2.3), dan
menuliskan keduanya ke kolomnya masing-masing pada `TREATY_IN_EDM`. Jadi ketidaksinkronan itu
**tidak berhenti di clipboard** — ia tersimpan ke tabel datar, dan setiap pembaca hilir yang
memakai salah satunya akan melihat cedant yang berbeda dari pembaca yang memakai yang lain.

Tiga hal yang mengikutinya, dan tidak satu pun diputuskan di sini:

| Ke mana | Pertanyaannya |
|---|---|
| **cabang I** | bila nama dan pengenal cedant tidak sinkron pada baris warisan, **mana yang dimigrasikan sebagai kebenaran** — dan siapa di bisnis yang dapat memastikannya? |
| **cabang F** | model baru menyimpan **hanya rujukan** ke cedant, namanya diturunkan dari rujukan itu. Ditandai **PERUBAHAN**, berubah dari: nama dan pengenal disimpan terpisah tanpa pemeriksaan apa pun |
| **ke orang** | pembaca hilir — laporan, akuntansi, pihak lawan — membaca **nama** atau **pengenal**? Tidak terbaca dari ekspor (`METODE` §4.6); masuk daftar untuk dibantah |

---

## 5. Mesin selisih - `ValueDifference` = baru minus `OLDDATA`

### 5.1 Susunan pemanggilan

`Activity/TreatyEDMCalculateDifference.xml`, 11 langkah, **sepuluh di antaranya tanpa syarat**
(hanya langkah 11 bersyarat `TreatyIn.IsProRate == true`):

| # | Panggilan | Isi |
|---|---|---|
| 1 | `Property-Remove` | mengosongkan pohon `ValueDifference` lebih dulu |
| 2 | `TreatyEDMDifferencePremium` | EGNPI dan total premi |
| 3 | `TreatyEDMDifferenceLimits` | `Limits`, `LimitSummaryList`, lima `TotalLimit*` |
| 4 | `TreatyEDMDifferenceShare` | `Share` dan seluruh daftar bersarangnya |
| 5 | `TreatyEDMDifferenceDeduction` | potongan (dua langkah terakhirnya **mati**) |
| 6-8 | `TreatyInDifferenceShareSumary`, `...ShareTotal`, `TreatyInSummaryLimitShareActual` | aturan **milik Treaty In**, dipakai ulang |
| 9-10 | `TreatyInSetValueDifferenceInstallment` | selisih angsuran |
| 11 | `TreatyEDMProRateCalculation` | pro rata, **bila `TreatyIn.IsProRate == true`** (lima langkahnya **mati**) |

### 5.2 Bentuk rumusnya seragam

Setiap besaran mengikuti satu bentuk:

```
TreatyIn.ValueDifference.<jalur> = TreatyIn.<jalur> - TreatyIn.OLDDATA.<jalur>
```

Tidak ada pembulatan, tidak ada penjagaan nilai kosong, tidak ada pemeriksaan mata uang: dua nilai
dengan mata uang berbeda tetap dikurangkan. Satu-satunya pengecualian bentuk adalah `Deductible2`,
yang dibungkus `@toDecimal(...)` di kedua sisi - bukti bahwa sebagian nilai memang tersimpan sebagai
teks (penyakit yang sudah dicatat di modul induk).

### 5.3 Pemadanan barisnya menurut POSISI, bukan kunci bisnis

Untuk setiap daftar, baris ke-*n* yang baru dikurangi baris ke-*n* yang lama:

```
TreatyIn.ValueDifference.EGNPI(<LAST>).Amount
        = .Amount - TreatyIn.OLDDATA.EGNPI(.pxListSubscript).Amount
TreatyIn.ValueDifference.Limits(<CURRENT>).Limit
        = .Limit - TreatyIn.OLDDATA.Limits(<CURRENT>).Limit
TreatyIn.ValueDifference.Share(local.subscript).RnmLimitList(<CURRENT>).Value
        = .Value - TreatyIn.OLDDATA.Share(local.subscript).RnmLimitList(<CURRENT>).Value
```

`pxListSubscript`, `<CURRENT>`, dan `local.subscript` semuanya **nomor urut dalam daftar**. Tidak
ada satu pun pemadanan lewat nomor layer, pengenal pihak, atau kunci bisnis lain.

**Akibatnya, dan ini terbaca dari rumusnya sendiri, bukan disimpulkan:** bila sebuah baris
disisipkan atau dihapus di tengah daftar, seluruh baris sesudahnya dipadankan dengan baris lama
yang salah, dan selisihnya salah tanpa ada galat yang muncul. Bila daftar baru lebih panjang,
baris tambahannya dikurangi baris lama yang tidak ada. Panjang daftar tidak pernah dibandingkan.

Ini persis alasan `SEAM-ADJUSTMENT.md` §3 menuntut `KUNCI_PADANAN`. Sekarang tuntutan itu punya
bukti sistem lama, bukan hanya alasan rancangan.

### 5.4 Satu-satunya percabangan di mesin selisih: `EDMState == "3"`

`TreatyEDMDifferenceShare` langkah 1 dan 2 adalah pasangan saling-eksklusif yang **benar**:

| # | Step page | Precondition | Berjalan untuk |
|---|---|---|---|
| 1 | `TreatyIn.ActualValue` | `TreatyIn.EDMState=="3"` benar -> lanjut | addendum premi |
| 2 | `TreatyIn` | `TreatyIn.EDMState=="3"` benar -> LEWATI | revisi internal dan addendum eksternal |

Keduanya menulis `ValueDifference.RNMShare` dan `ValueDifference.BrokeragePercent`. Maka pada
addendum premi selisih share NuRe dan brokerage dihitung dari pohon `ActualValue`, pada dua jenis
lainnya dari pohon utama. Ini satu-satunya tempat di seluruh mesin selisih yang membedakan jenis
addendum.

> **TEMUAN 27 September 2026 — dan ia mencabut kelengkapan kalimat di atas.**
>
> Percabangan `EDMState == "3"` di `TreatyEDMDifferenceShare` memang satu-satunya **percabangan**.
> Tetapi akibatnya baru terbaca utuh ketika `TreatyEDMDifferencePremium` diperiksa dengan pertanyaan
> yang sama: **dari pohon mana ia membaca?**
>
> | Aktivitas | `pyStepsObjectName` langkah pengurangnya | Membaca `ActualValue`? |
> |---|---|---|
> | `TreatyEDMDifferenceShare` langkah 1 | `TreatyIn.ActualValue` | **ya**, untuk `EDMState == "3"` |
> | **`TreatyEDMDifferencePremium` langkah 1** | **`TreatyIn.EGNPI`** | **tidak** — `ActualValue` muncul di berkas itu **satu kali saja**, sebagai deklarasi `pyPagesAndClasses`, tidak pernah sebagai objek langkah maupun nilai |
>
> Rumusnya: `ValueDifference.EGNPI(<APPEND>).Amount = .Amount - TreatyIn.OLDDATA.EGNPI(.pxListSubscript).Amount`,
> beriterasi atas `TreatyIn.EGNPI`.
>
> **Maka pada addendum premi, selisih EGNPI selalu nol.** Pengguna menyunting
> `TreatyIn.ActualValue.EGNPI` — pohon yang disemai `TreatyInSetEditPre` langkah 3.5 justru untuk
> itu — sedangkan mesin selisih mengurangkan `TreatyIn.EGNPI` terhadap `TreatyIn.OLDDATA.EGNPI`,
> dan **keduanya masih berisi nilai kontrak asal**, sebab jalur addendum premi tidak menyentuh
> `TreatyIn.EGNPI`.
>
> Angka yang justru menjadi alasan jenis addendum itu ada — premi yang disesuaikan dari estimasi ke
> aktual — **tidak pernah muncul sebagai selisih**. Yang muncul hanya selisih share dan brokerage,
> yang memang dibaca dari `ActualValue` (§5.4 di atas).
>
> Dicatat sebagai **TDA-16**. Akibatnya bagi rancangan dinyatakan di **C3**, dan ia menjadi syarat
> yang wajib dijawab di sana: di bawah GRL-12, sebuah versi `PENYESUAIAN_PREMI` yang **hanya**
> mengubah premi aktual akan terbaca **non-material** — bertentangan dengan perilaku sistem lama
> yang memaksanya material.

Selebihnya bersyarat hanya pada hal teknis: penjumlahan EGNPI per mata uang
(`.Currency == local.Currency`, `local.Appendflag == "0"`, dan keluar-iterasi bila total nol),
pembuangan potongan bernilai nol (`.Deduction < 1`), dan tujuh langkah ber-blok `//` yang **mati**
di `TreatyEDMDifferenceDeduction` (langkah 4 dan 5, berlabel *"Facultative share not enable yet in
this edm"*) serta `TreatyEDMProRateCalculation`.

Dua langkah di `TreatyEDMDifferencePremium` membawa teks syarat `param.type == "egnpi"`, tetapi
`pyStepsPreCondition` keduanya bernilai **`false`** - syaratnya dinonaktifkan dan langkahnya berjalan
selalu. Lihat catatan di §6.2 tentang dua bentuk "mati" yang berbeda.

### 5.5 Sebagian selisih ditulis ke pohon yang salah

`TreatyEDMDifferenceShare` langkah 4 menghitung selisih tetapi menuliskannya ke
`TreatyIn.ActualValue.LimitShareSummaryList(...)`, bukan ke `ValueDifference`:

```
TreatyIn.ActualValue.LimitShareSummaryList(<CURRENT>).Limit
      = .Limit - TreatyIn.OLDDATA.LimitShareSummaryList(<CURRENT>).Limit
```

Ini penyakit yang sama yang sudah ditemukan di modul induk pada `TreatyInDifferenceFacShare`
(tercatat sebagai tambalan yang **tidak dibawa** ke sistem baru). Ia muncul lagi di jalur EDM, pada
aturan yang berbeda. Digabung menjadi satu butir: **"selisih ditulis ke `ActualValue`"** adalah pola
berulang, bukan kejadian tunggal.

Dan hasilnya **selalu dimusnahkan**, karena dua fakta bertemu:

* langkah 4 itu **tidak bersyarat**, jadi ia berjalan setiap kali mesin selisih dipanggil;
* mesin selisih dipanggil dari jalur pengajuan hanya ketika `EDMState` 1 atau 2
  (`TreatyInSubmitEDM` langkah 7, precondition `TreatyIn.EDMState == "1" || TreatyIn.EDMState == "2"`);
* dan untuk `EDMState` 1 dan 2 itulah `SaveTreatyIn_EDM_Act` **menimpa seluruh `ActualValue`**
  (§2.2), tepat sesudahnya (langkah 8 pengajuan).

Arah dampaknya: angka ringkasan share yang sempat terlihat di layar **tidak pernah tersimpan**.
Yang perlu diperiksa dengan data - bukan dengan kode - adalah apakah ada pemanggil lain mesin
selisih yang lolos dari urutan ini; `TreatyEDMCalculateDifference` juga dirujuk lima seksi layar
(`TreatyInActionButtons`, `TreatyInNONProportional`, `TreatyInTabsNonProportional`,
`TreatyInTabsNonProportionalValueDifference`, `InputTreatyInOffer`).

### 5.6 Daftar 33 akar `ValueDifference`, beserta penulisnya

Dihitung ulang atas ekspor Adjustment (33 akar, 165 jalur daun sesudah subscript dinormalkan).
Kolom penulis diambil dari aturan di `Activity/` dan `DataTransform/` yang menyebut akar itu;
`TreatyIn*` adalah aturan milik Treaty In yang dipakai ulang, `TreatyEDM*` khas jalur addendum.

| # | Akar | Penulis |
|---|---|---|
| 1 | `BrokeragePercent` | `TreatyEDMDifferenceShare`, `TreatyInDifferenceShare` |
| 2 | `EGNPI` | `TreatyEDMDifferencePremium`, `TreatyInDifferencePremium` |
| 3 | `FacultativeShare` | `TreatyInDifferenceFacShare` |
| 4 | `FacultativeShareBrokerage` | `TreatyInDifferenceFacShare` |
| 5 | `FacultativeShareList` | `TreatyEDMDifferenceDeduction`, `TreatyInDifferenceDeduction`, `TreatyInDifferenceFacShare`, `TreatyInDifferenceFacShareTotal` |
| 6 | `Installment` | `TreatyEDMProRateCalculation`, `TreatyInSetValueDifferenceInstallment` |
| 7 | `LimitFacShareSummaryList` | `TreatyInDifferenceSummaryFacShare` |
| 8 | `LimitShareSummaryList` | `TreatyEDMProRateCalculation`, `TreatyInDifferenceShareSumary` |
| 9 | `LimitSummaryList` | `TreatyEDMDifferenceLimits`, `TreatyInDifferenceLimits`, `TreatyInDifferenceLimitsSumary` |
| 10 | `Limits` | `TreatyEDMDifferenceLimits`, `TreatyInDifferenceLimits`, `TreatyInDifferenceLimitsSetTotal`, `TreatyInDifferenceLimitsSumary` |
| 11 | `RNMShare` | `TreatyEDMDifferenceShare`, `TreatyInDifferenceShare` |
| 12 | `Share` | `TreatyEDMDifferenceShare`, `TreatyEDMDifferenceDeduction`, `GetNilaiTotal`, `GetNilaiTotalAdjustment`, `FetchQSfromMasterXOL` |
| 13 | `TotalEgnpiAmount` | `TreatyEDMDifferencePremium`, `TreatyInDifferencePremium` |
| 14 | `TotalEgnpiAmountNP` | idem |
| 15 | `TotalEgnpiProportion` | idem |
| 16 | `TotalFacShareDeductionNP` | `TreatyInDifferenceFacShareTotal` |
| 17 | `TotalFacShareGrossNP` | idem |
| 18 | `TotalFacShareNetNP` | idem |
| 19 | `TotalFacShareRnmNP` | idem |
| 20 | `TotalInstallmentNP` | `TreatyEDMProRateCalculation`, `TreatyInSetValueDifferenceInstallment` |
| 21 | `TotalLimitDeductblNP` | `TreatyEDMDifferenceLimits`, `TreatyInDifferenceLimits`, `TreatyInDifferenceLimitsSetTotal` |
| 22 | `TotalLimitIOONP` | idem |
| 23 | `TotalLimitMDPNP` | idem |
| 24 | `TotalLimitPremiEarnNP` | idem |
| 25 | `TotalLimitsROL` | idem |
| 26 | `TotalShareDeductionNP` | `TreatyEDMProRateCalculation`, `TreatyInDifferenceShareTotal` |
| 27 | `TotalShareGrossNP` | idem |
| 28 | `TotalShareNetNP` | idem, ditambah `TreatyInSetValueDifferenceInstallment` |
| 29 | `TotalShareRnmNP` | `TreatyEDMProRateCalculation`, `TreatyInDifferenceShareTotal` |
| 30 | `TotalSpreadedNetPremi` | `GetNilaiTotalAdjustment`, `TreatyEDMProRateCalculation`, `TreatyInDifferenceShareTotal` |
| 31 | `TotalSpreadedNetPremiRI` | idem |
| 32 | `TotalSpreadedRnmProp` | idem |
| 33 | `TotalSpreadedRnmRIProp` | idem |

**Dua hal yang terbaca dari daftarnya, dan tidak dari hitungannya.**

1. **Dua puluh satu dari 33 akar adalah agregat `Total*`** — turunan dari keduabelas akar lainnya.
   Bila selisih disimpan, hanya dua belas akar pertama yang merupakan fakta; sisanya dapat
   dihitung ulang darinya.
2. **Sebagian besar akar ditulis oleh aturan milik Treaty In, bukan oleh aturan khas addendum.**
   Hanya tujuh akar yang punya penulis ber-awalan `TreatyEDM*`; sisanya ditulis `TreatyIn*` yang
   dipakai bersama. Ini menguatkan §1.3: permukaan khas Adjustment jauh lebih tipis daripada
   tampaknya.

---

## 6. Penyimpanan

### 6.1 Tabel yang disentuh jalur addendum

| Aturan | Sasaran |
|---|---|
| `SaveTreatyInEDM` -> `POOLDATA.PEGA_M_TREATY_IN_EDM` | `M_TREATY_IN_EDM` (JSON) **dan** `TREATY_IN_EDM` (24 kolom datar) |
| `SaveTreatyInDetailEdm` -> `POOLDATA.PEGA_M_TREATY_IN_DETAIL_EDM` | `M_TREATY_IN_DETAIL_EDM` dan `TREATYINDETAILEDM` |
| `RemoveTreatyInEDM` | `DELETE FROM M_TREATY_IN_EDM WHERE ID = …` |
| `RemoveTreatyInEDM2` | `DELETE FROM TREATY_IN_EDM WHERE ID = …` |
| `RemoveTreatyInDetailEdm` | `DELETE FROM TREATYINDETAILEDM` **dan** `DELETE FROM M_TREATY_IN_DETAIL_EDM WHERE JSONDATA.TREATYID = …` |
| `BrowseTreatyInEDM` | `select JSONDATA from M_TREATY_IN_EDM where ID = …` |
| `EDMCheckExistingData` | `select count(ID) from treaty_in_edm where id = …` |
| `CopyAllAttachment2_Sql` | `select … from M_ATTACHMENTTREATY_2 where treatyid = {TreatyIn.OLDID}` |

Pola pasangan JSON + tabel datar sama persis dengan modul induk. Yang perlu dicatat untuk migrasi:
**tabel datar `TREATY_IN_EDM` tidak memuat `JSONDATA`**, sehingga siapa pun yang menyusun model dari
tabel relasional akan kehilangan seluruh isi kontraknya — pengulangan persis masalah yang sudah
ditemukan di Treaty In.

### 6.2 Sisip-atau-perbarui ditentukan sebuah hitungan, dan tabrakan berakhir sebagai perbaruan

`SaveTreatyIn_EDM_Act` langkah 5–8:

```
5  RDB-List EDMCheckExistingData        -> ExistingDataResult          (blok "jmp")
6  InputParam.STSCARI1 <- ExistingDataResult.pxResults(1).HASIL1
7  InputParam.DATAPEGA <- @ASM.GetPageJSONString()   [when PoductName.TYPE=="" || GRUP=="" -> LEWATI]
8  RDB-List SaveTreatyInEDM
```

Dan di dalam `POOLDATA.PEGA_M_TREATY_IN_EDM`:

```sql
IF STSINPUT = '0' THEN  INSERT INTO M_TREATY_IN_EDM …;  INSERT INTO TREATY_IN_EDM …;
ELSE                    UPDATE M_TREATY_IN_EDM SET JSONDATA = DataPega WHERE ID = IDPega;
                        UPDATE TREATY_IN_EDM SET … WHERE ID = IDPega;
END IF;
…
ErrMsg := 'Data Sudah Disimpan Dengan ID : ' || id_TreatyIn_out;  StsSave := 1;
```

Maka pengenal yang bertabrakan (§4.4) **tidak ditolak** — ia masuk ke cabang `UPDATE` dan
**menimpa addendum yang sudah ada**, lalu melaporkan berhasil. Arah dampaknya: kehilangan data
tanpa galat, pada addendum yang **bukan** yang sedang dikerjakan pengguna.

Satu hal yang **hampir** saya laporkan sebagai cacat dan ternyata bukan: langkah 7 dan 9 memuat
teks precondition `PoductName.TYPE == "" || PoductName.GRUP == ""`, dan halaman `PoductName` tidak
pernah diisi di mana pun dalam ekspor — **diuji ulang 24 September 2026 dengan sapuan terkalibrasi
dan tetap berdiri**: kata `PoductName` muncul **delapan kali di badan kedua ekspor, semuanya sebagai
teks precondition** (`pyStepsPreCondParamsWhen`) di `SaveTreatyIn_Act` dan `SaveTreatyIn_EDM_Act`,
dan **nol kali** sebagai sasaran tulis pada kelima jenis penulis. **Titik butanya kemudian ditutup
(25 September 2026):** ketiga puluh dua `pyStepsJavaSource` disapu sebagai teks - kata `PoductName`
**tidak muncul satu kali pun** - dan satu-satunya penulisan generik di sana, `adoptJSONObject`,
bekerja pada halaman `TreatyIn` sehingga **tidak dapat menciptakan halaman `PoductName`** yang
merupakan halaman tingkat atas, bukan properti `TreatyIn` — sehingga kondisinya selalu benar dan kedua langkah itu akan
selalu dilewati. Tetapi `pyStepsPreCondition` pada keduanya bernilai **`false`**: **preconditionnya
dinonaktifkan**, dan teksnya hanya sisa. Kedua langkah berjalan tanpa syarat.

> Ini pelajaran membaca, bukan temuan sistem: sebuah langkah punya **dua** penanda mati yang berbeda
> — blok `//` mematikan langkahnya, `pyStepsPreCondition = false` mematikan **syaratnya** sementara
> langkahnya tetap berjalan. Membaca teks syarat tanpa memeriksa penandanya membalik kesimpulan.

### 6.3 Dua langkah yang mengganti baris kontrak — dimatikan

`SaveTreatyIn_EDM_Act` langkah 12 dan 13 ber-blok `//`:

| # | Langkah | Precondition tertulis |
|---|---|---|
| 12 | `RDB-List InsertToMTreatyIn` — *"when Resolve Complete insert into m treaty in"* | `TreatyIn.StatusAkseptasi == "Resolve Complete"` |
| 13 | `call SaveTreatyInDetail_Act` — *"save treatyindetail"* | idem |
| 14 | `call SaveTreatyInDetailEdm_Act` — *"save treatyindetailedm"* | idem, **HIDUP** |

Ini menegaskan temuan G1b sesi sebelumnya dengan kondisi yang kini terbaca: kedua langkah itu
**punya syarat keadaan `Resolve Complete`**, artinya ditulis untuk maksud "addendum yang disetujui
menggantikan kontrak", lalu dimatikan. Yang tersisa hidup hanya penulisan ke tabel **detail
addendum**. Di sistem lama, **addendum yang disetujui tidak menggantikan kontraknya**.

### 6.4 `Force Resolve Complete(dev)` di layar addendum menyimpan ke tabel KONTRAK - dan tidak menulis apa pun

Tombol itu memanggil `SaveTreatyIn_Act`, yaitu jalur simpan **kontrak**. Ketika halaman clipboard
adalah addendum, rantainya terbaca utuh:

| Langkah | Yang terjadi |
|---|---|
| `SaveTreatyIn_Act` 3 | `InputParam.IDPEGA = TreatyIn.ID` - untuk addendum bernilai `‹kontrak›/Rnn` |
| `SaveTreatyIn_Act` 4 | `RDB-List SaveTreatyIn` -> `POOLDATA.PEGA_TREATY_IN(...)` |
| `PEGA_TREATY_IN` | `IF IDPega = 'UnknownId' THEN` sisip `M_TREATY_IN` + `TREATY_IN`; `ELSE` **`UPDATE … WHERE ID = IDPega`** pada kedua tabel itu |
| `SaveTreatyIn_Act` 10 | `TreatyIn.ID = Param.IDPEGAOUT` |

Pengenal addendum bukan `'UnknownId'`, jadi ia **selalu masuk cabang `ELSE`**. Dan
`UPDATE ... WHERE ID = 'XXXXXXX/R01'` pada `M_TREATY_IN` dan `TREATY_IN` **tidak menemukan baris**,
sebab addendum tinggal di `M_TREATY_IN_EDM` / `TREATY_IN_EDM`.

**Di Oracle, `UPDATE` yang tidak mengenai baris bukan galat.** Tidak ada exception, `COMMIT`
tetap dijalankan, lalu:

```
StsSave   := 1;
ErrMsg    := 'Data Sudah Disimpan Dengan ID : ' || id_TreatyIn_out;
IDPegaOut := IDPega;
```

Jadi hasilnya dua-duanya:

* **Baris kontrak asalnya tidak tersentuh.** Penjaganya `WHERE ID = IDPega`, dan pengenal addendum
  tidak pernah sama dengan pengenal kontrak (tujuh karakter). Kekhawatiran "menimpa kontraknya"
  **tidak berdasar**.
* **Addendumnya juga tidak tersimpan.** Perubahan `StatusAkseptasi = "Resolve Complete"`,
  `Position = ""`, `PositionUsername = ""` yang baru saja ditulis `TreatyInForceResolveComplete`
  hanya ada di clipboard, lalu hilang begitu halaman ditutup - sementara layar melaporkan
  *"Data Sudah Disimpan"*.

Ini **kegagalan senyap** dalam bentuknya yang paling murni (`METODE` §2.2): jalur yang salah,
laporan yang benar. Sisinya: aturannya `SaveTreatyIn_Act` ada di **irisan**, tetapi
**keterjangkauannya dari layar addendum 56-KHAS** - `Section/InputTreatyInAdjustment` hanya ada di
ekspor Adjustment.

**Tidak ada keadaan kedua yang dapat bercabang.** Diperiksa 25 September 2026: sel tombol itu
memuat **satu** aksi saja - `pyAction = refresh`, `pyTarget = thisSection`, dengan
`pyPreDataTransform = TreatyInForceResolveComplete` dan `pyActivity = SaveTreatyIn_Act`. Tidak ada
flow action, tidak ada penyelesaian assignment, tidak ada `Obj-Save`.

Dan itu bukan kebetulan tombol ini: **modul ini tidak memakai objek kerja Pega sama sekali.** Di
seluruh ekspor Adjustment tidak ada satu pun langkah `Obj-Save`, `Work-.*`, maupun penyelesaian
assignment; `Assign-Worklist` muncul hanya sebagai kelas *applies-to* sebuah RDB-List pinjaman
(`GetCountClaim`), dan `pyWorkPage` hanya sebagai nama halaman di aktivitas riwayat klaim dan
Report Definition. Kelas modulnya `Data-Portal`, kelas **data**, bukan kelas kerja; seluruh keadaan
tinggal di `JSONDATA` POOLDATA.

**Maka pengukuran dari tabel objek kerja Pega tidak mungkin** - bukan karena sulit, melainkan
karena objeknya tidak ada. Kegagalannya benar-benar hanya "tidak tersimpan": tidak ada keadaan
perkara yang dapat bercabang dari keadaan baris, sebab tidak ada perkara.

**Yang tidak terbaca:** apakah pernah ada yang menekannya di layar addendum. Tidak ada jejak yang
ditinggalkan - justru karena tidak ada yang tersimpan. Ia karena itu **tidak dapat diukur UA mana
pun**, dan itu harus dikatakan apa adanya.

---

## 7. Persetujuan

### 7.1 Rantainya sama dengan kontrak — hanya pembungkusnya yang berbeda

Tiga aturan khas EDM (`TreatyInSubmitEDM`, `TreatyInAkseptasiEDM_Act`,
`TreatyInDeclineConfirmation_postactEDM`) semuanya memanggil **`Akseptasi_DT` yang sama** dengan
jalur kontrak. Mesin keadaannya satu, bukan dua.

`Akseptasi_DT` bercabang di puncak pada `RevisionState`, bukan pada `EDMState`:

| Cabang | Syarat | Rantai |
|---|---|---|
| 1 | `RevisionState != 1 && StatusAkseptasi != "Resolve Complete"` | Admin -> SecHead -> DeptHead -> Director -> `Resolve Complete` |
| 2 | `RevisionState == 1 && StatusAkseptasi != "Resolve Complete"` | Admin -> SecHead -> `Resolve Complete` |

Cabang untuk `ReasTreatyInGroupLeader` ada di antara keduanya dan **seluruhnya mati** (blok `//`).

Nama orang disemat langsung di dalam aturan (`"IRVANDY"`, `"AGUNGPUTRAANDALAS"`,
`"YOHANESKRISTIAWAN"`, `"NANDINA"`), dan perannya dibaca dari `OperatorID.pyTelephone`
(`"SPVTREATY1"`, `"TREATY1"`, …) — kolom nomor telepon dipakai sebagai kolom peran.

### 7.2 KOREKSI 24 September 2026 — panjang rantai untuk addendum TIDAK terbaca

**Versi pertama bagian ini salah, dan dicabut.** Ia menyatakan bahwa addendum yang lahir dari
pemilihan sebuah kontrak menempuh rantai dua tingkat, karena `SetTreatyIn_Act` menetapkan
`RevisionState = 1`. Dasar itu runtuh bersama §4.5: **jalur Adjustment tidak pernah mengirim
`param.revisionstate`**, sehingga langkah yang menuliskannya dilewati.

**Yang benar, dan yang tidak terbaca.**

* `Akseptasi_DT` memang bercabang pada `TreatyIn.RevisionState`, bukan pada `EDMState` — itu tetap
  berdiri, dan tetap berarti **panjang rantai tidak ditentukan jenis addendum**.
* Tetapi **nilai `RevisionState` pada sebuah addendum tidak ditulis oleh jalur addendum mana pun**.
  Ia **diwarisi** dari `JSONDATA` baris yang dimuat, dan tidak ada aturan di ekspor ini yang
  mengosongkannya saat addendum dibuat.
* Maka rantai dua tingkat **dapat** terjadi pada sebuah addendum — bila baris asalnya kebetulan
  membawa `RevisionState = 1` — tetapi apakah itu terjadi, dan seberapa sering, **tidak dapat
  dijawab ekspor**. Ia pertanyaan data: **UA-9**.

Yang tetap berdiri tanpa syarat: cabang 2 `Akseptasi_DT` **mengosongkan `RevisionState`** begitu
`Resolve Complete` tercapai (aksi `2.2.1.4`), sehingga jejak "mengapa rantainya pendek" hilang
bersama nilainya. Itu tetap TDA-08.

Pertanyaan bisnisnya tidak berubah bentuknya, dan tetap **daftar untuk dibantah**: *"revisi internal
berhenti di kepala seksi, penyesuaian premi naik sampai direktur — benar begitu?"*

### 7.3 Penolakan menghapus addendum, tanpa jejak

`Activity/TreatyInDeclineConfirmation_postactEDM.xml`:

| # | Langkah | Keadaan |
|---|---|---|
| 1 | set `Param.SkipVisibility/Info/Comment` | **MATI** |
| 2 | `Call AddCommentList_Act` | **MATI** |
| 3 | `CommentList(<LAST>).IsApproved <- "Decline"`, `StatusAkseptasi <- "Decline"` | **MATI** |
| 4 | `Call SaveTreatyIn_Act` | **MATI** |
| 5 | `InputData.CARI1 <- TreatyIn.ID` | hidup |
| 6 | `RDB-List RemoveTreatyInEDM` -> `DELETE FROM M_TREATY_IN_EDM` | hidup |
| 7 | `RDB-List RemoveTreatyInEDM2` -> `DELETE FROM TREATY_IN_EDM` | hidup |
| 8 | `call TreatyInInputVis` | hidup |

Bandingkan dengan jalur kontrak, `TreatyInDeclineConfirmation_postact`, yang **empat langkah
itu justru hidup** dan tidak menghapus apa pun.

> **Menolak sebuah addendum memusnahkan barisnya dari kedua tabel.** Tidak ada `StatusAkseptasi =
> "Decline"` yang tersimpan, tidak ada komentar, tidak ada baris tinggal. Langkah-langkah yang
> dulu mencatatnya masih ada di aturannya, dimatikan.

Tiga hal yang ikut lenyap bersamanya, dan tidak ada satu pun yang membersihkannya:
`TREATYINDETAILEDM` / `M_TREATY_IN_DETAIL_EDM` (dihapus hanya oleh `RemoveTreatyInDetailEdm`, yang
**tidak dipanggil** di jalur ini), lampiran yang sudah disalin `TreatyRevisionCopyAttachment`, dan
`RevisionState = 1` yang terlanjur tertulis di baris kontraknya (§4.5) — kontraknya tetap tertandai
sedang direvisi meski revisinya sudah tidak ada.

---

### 7.4 KOREKSI 24 September 2026 - daur hidup `RevisionState`, terbaca lengkap

§7.2 menyatakan bahwa nilai `RevisionState` sebuah addendum **diwarisi** dan bahwa bagaimana ia
sampai di situ *"tidak dapat dijawab ekspor"*. Bagian pertama tetap berdiri. Bagian kedua
**dicabut**: dengan menyapu seluruh penulis `RevisionState` di kedua ekspor - bukan mencari
namanya, melainkan mendaftar setiap aksi tulis - jalannya terbaca utuh.

Sapuan sebelumnya luput karena memakai nama tag yang salah. Sasaran tulis sebuah **Data Transform**
bernama `pyPropertiesName` / `pyPropertiesValue`, bukan `pyName` atau `pyTarget`; sapuan dengan tag
yang salah mengembalikan nol, dan nol itu sempat saya baca sebagai jawaban (`MA-05`).

**Setiap penulis `TreatyIn.RevisionState` di kedua ekspor - daftar tertutup:**

| Aturan | Langkah | Nilai | Yang ikut ditulis di langkah yang sama |
|---|---|---|---|
| `SetTreatyIn_Act` | 7 | **`1`** | `StatusAkseptasi = ""`, `PositionUsername = OperatorID.pxInsName` |
| `Akseptasi_DT` | 2.2.1.4 | `""` | SecHead **Accept** -> `StatusAkseptasi = "Resolve Complete"`, `ViewState = ""` |
| `Akseptasi_DT` | 2.2.2.4 | **`1`** | SecHead **Reject** -> `StatusAkseptasi = "Reject"`, `Position = "ReasTreatyInAdmin"` |
| `Akseptasi_DT` | 2.2.3.4 | `""` | SecHead **Decline** -> `StatusAkseptasi = "Decline"` |
| `TreatyInCopy` | 1 | `""` | **ditambahkan 24 September 2026 (audit `TA-04`)** - aktivitas "salin kontrak", dirujuk dua kali di badan `Section/InputTreatyInOffer`. Ia **mengosongkan** penanda pada salinan, dan **tidak ada di jalur addendum** |

Tidak ada penulis lain — disapu atas kelima jenis penulis dengan perkakas terkalibrasi. Maka:

1. **`RevisionState = 1` hanya lahir di satu tempat** - langkah 7 `SetTreatyIn_Act`, yang
   preconditionnya `param.revisionstate==1`, dikirim hanya oleh kontrol 7 dan 8
   `Section/InputTreatyInOffer` (§4.5).
2. **Langkah yang sama mengosongkan `StatusAkseptasi`.** Jadi sebuah kontrak yang sedang
   direvisi-di-tempat **tidak lagi berstatus `Resolve Complete`** sejak tombolnya ditekan.
3. **Penolakan TIDAK mengosongkannya.** Cabang `2.2.2` justru **menulis ulang `"1"`** - kontrak
   kembali ke inputor dengan tanda "sedang direvisi" tetap menempel. Yang mengosongkannya hanya
   **Accept** (lewat `Resolve Complete`) dan **Decline**. Ini menjawab pertanyaan pemeriksaan
   secara langsung, dan jawabannya **tidak** seperti yang diduga.
4. **Tidak ada cabang yang berjalan** bila `RevisionState == 1` **dan** `StatusAkseptasi ==
   "Resolve Complete"` sekaligus: cabang 1 menuntut `RevisionState != 1`, cabang 2 menuntut
   `StatusAkseptasi != "Resolve Complete"`. Keadaan itu **membeku** - `Akseptasi_DT` tidak
   melakukan apa pun atasnya.

**Pewarisan ke addendum, terbukti ujung ke ujung.** `TreatyInEDMSetValue` - aktivitas di balik
tombol *Choose* pada picker - langkah 3 menyalin halaman kontrak dan menetapkan `OLDID`,
`CommentList = ""`, `EDMState`, `EDMMaterialType`, `RevisionDate`. **`RevisionState` tidak
disentuh.** Langkah 5 menerapkan `TreatyInSetEdit`, yang juga tidak mengosongkannya - dan yang
langkah 2-nya justru **membacanya**:

```
TreatyInSetEdit
  1   ViewState        = 0
  2   WHEN TreatyIn.RevisionState == 1
  2.1     ViewState    = 1
  3   Position         = "ReasTreatyInAdmin"
  6   StatusAkseptasi  = ""
```

Maka sebuah addendum yang mewarisi `RevisionState = 1` bukan hanya menempuh **rantai dua tingkat**
(`Akseptasi_DT` cabang 2), ia juga **lahir dalam keadaan terkunci** (`ViewState = 1`). Dua akibat,
satu sebab.

**Dan sebabnya dapat dimasuki.** Dua jalan, keduanya hidup:

* **Lewat daftar pilihan yang tidak menyaring apa pun** (§7.6, TDA-11): kontrak yang sedang
  direvisi-di-tempat tetap muncul di picker, karena daftarnya tidak memuat saringan keadaan sama
  sekali.
* **Lewat tombol `Force Resolve Complete(dev)`** (§7.5): ia menulis `StatusAkseptasi =
  "Resolve Complete"` **tanpa menyentuh `RevisionState`**, lalu menyimpan. Hasilnya persis keadaan
  beku butir 4 di atas - kontrak yang tampak selesai tetapi tetap bertanda sedang direvisi, dan
  setiap addendum yang lahir darinya mewarisi rantai pendek.

Berapa banyak baris yang benar-benar membawanya tetap pertanyaan data: **UA-9**, kini dipertajam.

### 7.5 Tiga tombol "(dev)" di `Section/TreatyInActionButtons`, dan dua di antaranya tidak sesuai labelnya

> **Atribusi.** Keempat "pintu samping" sudah didaftar **ADR-0055 §4** pada 23 September 2026,
> termasuk catatan bahwa `TreatyInForceEdit` terlihat oleh setiap pengguna layar penawaran. Bagian
> ini **menambah dan mengoreksi** daftar itu; ia bukan penemuan ulang (`MA-06`).

Seksi `TreatyInActionButtons` disertakan oleh `Section/InputTreatyInOffer` **dan**
`Section/InputTreatyInAdjustment` - badan empat kemunculan pada masing-masing. Jadi tombol-tombol
ini ada di layar addendum juga, bukan hanya di layar kontrak. ADR-0055 menyebut layar penawaran
saja.

> **KOREKSI 24 September 2026, dan ia membalik klaim saya sendiri.** Rumusan pertama bagian ini
> berbunyi *"`Force Edit (dev)` terlihat oleh semua orang"*, atas dasar `pyVisible = ALWAYS` pada
> **selnya**. Rantai syaratnya ditelusuri sampai akar seksi, dan ada satu lapis di atasnya:
>
> ```
> Section/TreatyInActionButtons.xml, layout pySectionId = S3
>   pyIsVisibilityOption   = CONDITION
>   pyContainerVisibleWhen = OperatorID.pyOrgDivision = 'IT'
> ```
>
> **Ketiga tombol "(dev)" berada di dalam layout itu.** Syarat sel dan syarat wadah berlaku
> **bersamaan** - keduanya dikalikan, bukan yang satu menggantikan yang lain. Arti `pyVisible` terbaca dari sebarannya di ekspor Adjustment:
>
> | Nilai | Jumlah | Yang menyimpan `pyCondition` |
> |---|---|---|
> | `ALWAYS` | 6663 | **37** - sisa, tidak dievaluasi |
> | `OTHER` | 553 | **553 - seluruhnya** |
> | `NOTBLANK` | 5 | 0 |
>
> Pemadanan sempurna pada `OTHER` (553 dari 553) menunjukkan `pyVisible` adalah **pemilih** aturan
> tampil: `OTHER` berarti "pakai `pyCondition`", `ALWAYS` berarti "abaikan apa pun yang tersimpan
> di sebelahnya". Maka hasil akhirnya per tombol:
>
> | Tombol | Syarat wadah | Syarat sel | **Terlihat oleh** |
> |---|---|---|---|
> | `Force Resolve Complete(dev)` | divisi `IT` | `OTHER` - dua nama | **operator divisi IT yang bernama ALDO SAPUTRA atau Daniel Suhana** |
> | `ReturnToInputor(dev)` | divisi `IT` | `OTHER` - dua nama | idem |
> | `Force Edit (dev)` | divisi `IT` | `ALWAYS` - nama diabaikan | **setiap operator divisi IT** |
>
> Ini mengoreksi **ADR-0055 §4**, yang mencatat `TreatyInForceEdit` *"terlihat oleh setiap pengguna
> layar penawaran"*.
>
> **Penjaganya nama lengkap yang ditanam di dalam aturan** - `'ALDO SAPUTRA'`, `'Daniel Suhana'` -
> bukan peran, bukan privilese. Itu keluarga **TDA-09**: peran dan orang disemat sebagai teks di
> aturan, sehingga berganti nama tampilan mengubah arti penjaganya tanpa ada yang menyentuh kode.
>
> `pyAssociatedPrivileges = false` pada **setiap** lapis: tidak ada privilese sama sekali; penjaganya
> hanya divisi dan nama.
>
> Pada `Section/InputTreatyInAdjustment` ada satu lapis lagi: tiga dari empat penyertaan
> `TreatyInActionButtons` berada di layout `S14` ber-`pyContainerVisibleWhen = OutputParam.DATASHOW=1`.

| Tombol | `pyVisible` sel | Kondisi sel | Yang dijalankan | Menyimpan? |
|---|---|---|---|---|
| `Force Edit (dev)` | **`ALWAYS`** | `OperatorID.pyUserName = 'ALDO SAPUTRA' \|\| … 'Daniel Suhana'` | `TreatyInForceEdit` -> `ViewState = 0` | tidak |
| `Force Resolve Complete(dev)` | `OTHER` | idem | `TreatyInForceResolveComplete` -> `StatusAkseptasi = "Resolve Complete"`, `Position = ""`, `PositionUsername = ""` | **ya** - `pyActivity = SaveTreatyIn_Act` |
| `ReturnToInputor(dev)` | `OTHER` | idem | `TreatyInReturntoInputor` -> `StatusAkseptasi = ""`, `Position = "ReasTreatyInAdmin"` | tidak |

Dua hal yang harus dibaca terpisah:

1. **`Force Edit (dev)` terlihat oleh setiap operator divisi `IT`** - bukan dua orang, dan bukan
   semua orang. `pyVisible = ALWAYS` membuang kondisi namanya; gerbang divisi di atasnya tetap
   berlaku. Akibatnya setiap orang di divisi IT dapat membuka kunci formulir yang seharusnya
   terkunci (`ViewState = 0`). Tombol itu sendiri tidak menyimpan; yang menyimpan adalah tombol
   simpan berikutnya, yang kini menjadi dapat ditekan.
2. **`Force Resolve Complete(dev)` menyimpan, dan meninggalkan `RevisionState`.** Ia satu-satunya
   aturan yang dapat menghasilkan pasangan `Resolve Complete` + `RevisionState = 1`, yaitu keadaan
   beku §7.4 butir 4. Kondisinya memang dievaluasi (`pyVisible = OTHER`), tetapi penjaganya adalah
   **pembandingan nama orang sebagai teks**, bukan peran maupun privilese.

Ketiganya berlabel "(dev)" dan tidak satu pun dijaga oleh privilese. Penjaganya **nama orang
sebagai teks**; bila nama tampilan orangnya berubah, penjaganya berubah arti tanpa ada yang
menyentuh kode.

**Pintu keempat ADR-0055 terkunci dari dalam.** ADR-0055 §4 mencatat `TreatyInSetToDirector`
"melompati dua tingkat", terjangkau lewat `TreatyInTestAgent` di layar penawaran. Jalan menujunya
memang hidup - `TreatyInTestAgent` dirujuk **dua kali di badan** `Section/InputTreatyInOffer`, dan
satu-satunya langkahnya memanggil `TreatyInSetToDirector`. Aktivitas itu sendiri, di kedua ekspor:

| Langkah | `pyStepsBlockName` | Isi |
|---|---|---|
| 1, 2, 3, 4 - **seluruh langkah tingkat atas** | **`//`** | ambil waktu Oracle, bandingkan, dan blok bersyarat |
| 4.1 - 4.7 | bersarang di dalam langkah 4 yang mati | 4.5 menulis `Position = ReasTreatyInDirector`, `PositionUsername = "NANDINA"` |

**Tidak ada satu langkah hidup pun.** Dan seandainya hidup, langkah 4.5 menuntut
`TreatyIn.Position == "ReasTreatyInGroupLeader"` - keadaan yang tidak dapat dimasuki. Dua lapis
mati sekaligus, dan yang kedua tersembunyi di balik yang pertama.

### 7.6 Daftar pilihan revisi: tidak ada saringan keadaan sama sekali - TDA-11 dikuatkan

Grid yang dilihat pengguna di `Section/PickerTreatyInMasterRevisi` bersumber pada
`pyPageListProperty = TempMasterList.pxResults`, dan `TempMasterList` diisi sebuah layout
ber-`pyLoadDeferred = true` yang memanggil `pyDeferLoadRetrievalActivity =
TreatyLoadMasterJoinEdm`. SQL-nya utuh:

```sql
select a.ID as CARI1, ... from pooldata.treaty_in a
UNION
select b.ID as CARI1, ... from pooldata.treaty_in_edm b
```

**Tanpa `WHERE` sama sekali** - tanpa saringan `StatusAkseptasi`, tanpa `RevisionState`, tanpa
pembeda kontrak-versus-addendum. Varian `TreatyLoadMasterJoinEdmXOL` menambah satu-satunya saringan
yang ada di keluarga ini, `PROPORTIONTYPE = 'NonProportional'`.

**Ada konfigurasi kedua yang bertentangan, dan tag pemilihnya menyebut siapa yang menang.**
Grid yang sama membawa `pyGridProps` berisi `pyRDName = BrowseTREATY_IN` dengan `pyRDParams` yang
memuat `StatusAkseptasi = "Resolve Complete"`. Yang **memilih** sumber adalah `pySourceType`, dan
nilainya tertulis dua kali - di baris layout dan di dalam `pyGridProps` - dengan nilai yang sama:

| Baris | Tag | Nilai |
|---|---|---|
| `Section/PickerTreatyInMasterRevisi.xml:5044` | `pyPageListProperty` | `TempMasterList.pxResults` |
| `Section/PickerTreatyInMasterRevisi.xml:5052` | **`pySourceType`** | **`Property`** |
| `Section/PickerTreatyInMasterRevisi.xml:7923` | `pyRDName` | `BrowseTREATY_IN` |
| `Section/PickerTreatyInMasterRevisi.xml:7991` | **`pySourceType`** | **`Property`** |

`pySourceType = Property` berarti sumbernya **properti daftar halaman**, bukan Report Definition.
Dua tanda pendukung di baris yang sama: `pyRDResultsPageExists = false` dan
`pyVirtualRDExists = false`. Jadi ini **terbaca**, bukan disimpulkan dari akibatnya; `pyRDName`
beserta parameternya adalah **sisa pengaturan lama** (`METODE` §2.0).

Konsekuensinya untuk §7.4: kontrak yang sedang direvisi-di-tempat - `StatusAkseptasi = ""`,
`RevisionState = 1` - **tetap muncul** di daftar, dan addendum yang dibuat darinya mewarisi rantai
pendek sekaligus lahir terkunci. Tidak perlu tombol "(dev)" untuk itu.

---

### 7.7 Apa yang sebenarnya dibuka tombol `Revision` — kuncinya selektif, dan ia bocor

`Revision` mengirim **dua** parameter, bukan satu: `viewstate=1` **dan** `revisionstate=1`. Maka
`SetTreatyIn_Act` menjalankan **dua** langkah:

| Langkah | Precondition | Yang ditulis |
|---|---|---|
| 6 | `param.viewstate==1` | `TreatyIn.ViewState = 1`, `TreatyIn.IsEditData = 1` |
| 7 | `param.revisionstate==1` | `TreatyIn.RevisionState = 1`, `StatusAkseptasi = ""`, `PositionUsername`, dan seterusnya |

Langkah 6 itu yang selama ini terlewat dari pembacaan saya: **`Revision` tidak hanya memendekkan
rantai, ia juga mengunci sesuatu.**

**Seberapa luas kuncinya — dihitung, bukan dikira.** Dari seluruh sel FIELD yang dapat disunting di
seluruh `Section/` ekspor Adjustment (badan aturan, `.pyTemplate*` dikecualikan):

| | Jumlah sel |
|---|---|
| **dijaga** `ViewState` (lewat `pyDisabledWhen` atau `pyReadOnlyCondition`) | **96** |
| **tanpa penjaga** `ViewState` | **882** |

Pola kondisinya dua macam: `TreatyIn.ViewState = 1` (dan variasi tanda kutipnya) sebanyak 183
kemunculan, dan `TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2` sebanyak 74 - satu penjaga
yang memikul **dua** maksud sekaligus.

**Yang dikunci terpusat pada uang di muka kontrak:**

| Berkas | Sel terkunci | Isinya |
|---|---|---|
| `DetailLimits` | 20 | mata uang dan nilai batas per bahaya - gempa, banjir, `RSMD` |
| `Layers` | 20 | `Deductible`, `AgregateLimit`, `AdjRate`, mata uang lapisan |
| `TreatyInTabsNonProportional` | 14 | `BrokeragePercent`, mata uang bahaya |
| `TreatyInTabsProportional` | 13 | `FacShare`, `FacShareBrokerage`, `BrokeragePercentP`, nama reasuradur dan broker |
| sebelas berkas lain | 1-3 masing-masing | `MaxRetention`, `EGNPI`, `Share`, `ShareRetro`, `RetroList` |

**Dan yang tidak dikunci — di situlah bocornya.** Ketiga field lapisan beku yang dapat disunting
(`Ceding`, `Commencement`, `Termination`) **tidak punya penjaga `ViewState` sama sekali**, di kedua
layar. Begitu pula 882 sel lain, termasuk banyak baris rincian di bawah seksi uang yang justru
dikunci di tingkat atasnya.

**Dicocokkan dengan daftar 33 akar `ValueDifference` (§5.6), dan hasilnya membatalkan bacaan
"kunci uang" itu sendiri.** Dua puluh satu dari 33 akar adalah agregat `Total*` yang tidak disunting
langsung; yang tersisa **dua belas akar** adalah besaran uang dan persentase yang benar-benar dapat
disesuaikan. Dihitung per seksi yang memasoknya:

| Seksi | Dijaga `ViewState` | Tanpa penjaga | Akar §5.6 yang disentuhnya |
|---|---|---|---|
| `Installments` | **0** | 2 — `.AmountTotal`, `.PctTotal` | #6 `Installment` |
| `DetailShareRetro` | **0** | 2 — `.Currency`, `.Value` | #5 `FacultativeShareList` |
| `ShareRetro` | 1 | **27** — termasuk `.Value`, `.Pct`, `.Currency` | #12 `Share` |
| `Share` | 2 | **33** — termasuk `.Value`, `.Pct`, `.Currency` | #11 `RNMShare`, #12 `Share` |
| `LayersEDM` | 3 | **18** — termasuk `.Deductible`, `.AgregateLimit`, `.AdjRate` | #10 `Limits` |
| `DetailShare` | 2 | 6 — termasuk `.Value`, `.Currency` | #5, #7 |
| `DetailEGNPI` | 2 | 2 — `.AmountIDR`, `.Proportion` | #2 `EGNPI` |
| `DetailLimits` | 20 | 5 — termasuk **`.Value`**, **`.Currency`** | #9, #10 |
| `Layers` | 20 | 3 — termasuk `.Value` | #10 `Limits` |
| `TreatyInTabsNonProportional` | 14 | **92** — termasuk `.Deductible`, `.AggregateLimit`, `.Amount` | #1, #10 |
| `TreatyInTabsProportional` | 13 | 8 | #3, #4 |
| `MaxRetention` | 2 | 0 | — |

**Tidak satu pun dari dua belas akar itu terkunci seluruhnya.** Dan pasangan yang paling telak:
`Layers.xml` menjaga `.Deductible`, `.AgregateLimit`, dan `.AdjRate` — sedangkan **`LayersEDM.xml`,
seksi kembarannya, tidak menjaga satu pun dari ketiganya.**

**Maka bacaan "kunci uang" itu gugur.** Yang tersisa bukan kebijakan melainkan **ketidakkonsistenan**:
field uang yang sama dikunci di satu seksi dan terbuka di seksi saudaranya. Sesudah `Revision`,
seorang pengisi masih dapat mengubah deductible, batas agregat, nilai dan persentase bagian, serta
total angsuran — dan perubahan itu menempuh rantai dua tingkat.

> **KOREKSI 27 September 2026 — penanda nonaktif ditafsirkan, dan daftar bocor menyusut.**
>
> Hitungan "96 dijaga / 882 tidak" **mengabaikan penanda nonaktif tingkat mode**. Ditafsirkan dulu,
> dengan kalibrasi (`TA-04`): `Installments..AmountTotal` - kolom total hasil hitung, diketahui
> hanya-baca - dan `InputTreatyInAdjustment..Ceding` - diketahui dapat disunting.
>
> | `pyDisabledNew` | `pyDisabled` | `pyDisabledWhen` | Jumlah mode | Arti |
> |---|---|---|---|---|
> | `false` | `-` atau `false` | - | 1.302 | **aktif** - `.Ceding` ada di sini |
> | **`always`** | `true` | - | **146** | **nonaktif selalu** - `.AmountTotal` ada di sini |
> | **`true`** | `true` | **ada** | **249** | **nonaktif bersyarat** - `pyDisabledWhen` yang menentukan |
> | `-` | `-` | - | 11.854 | mode tanpa setelan |
>
> **Ini terbaca, bukan disimpulkan:** 146 + 249 = **395**, tepat sama dengan jumlah `pyDisabled =
> true` di badan `Section/`; dan tidak ada satu pun mode ber-`pyDisabledNew = true` yang tanpa
> `pyDisabledWhen`, maupun ber-`always` yang punya syarat. Kedua kelompok tidak bertumpang tindih.
>
> **Sensus sel FIELD diulang dengan tafsir itu:**
>
> | | Jumlah |
> |---|---|
> | **aktif tanpa syarat** | **762** |
> | **selalu nonaktif** | **116** |
> | nonaktif bersyarat - `ViewState` | **83** |
> | nonaktif bersyarat - syarat lain | 17 |
>
> **Dan daftar bocor `DB-10` menyusut. Dugaan pemilik proses tentang kolom total benar:**
>
> | Yang saya sebut bocor | Keadaan sebenarnya |
> |---|---|
> | `Installments` `.AmountTotal`, `.PctTotal` | **selalu nonaktif** - kolom total hasil hitung. **Dicabut dari daftar bocor** |
> | `LayersEDM` `.Deductible`, `.AgregateLimit`, `.Currency`, dan sebelas lainnya | **selalu nonaktif** - empat belas sel. **Dicabut** |
> | `DetailEGNPI` `.AmountIDR`, `.Proportion` | **selalu nonaktif**. **Dicabut** |
>
> **Dan ini membalik kalimat paling tajam saya:** *"`Layers` menjaga `.Deductible`; kembarannya
> `LayersEDM` tidak"* - **terbalik**. `LayersEDM` justru **lebih ketat**: empat belas selnya
> nonaktif permanen, sedangkan `Layers` hanya menjaganya bersyarat. Kalimat itu dicabut.
>
> **Yang tetap bocor sesudah koreksi - dan ia tetap memuat angka uang:**
>
> | Seksi | Sel aktif tanpa syarat |
> |---|---|
> | `TreatyInTabsNonProportional` | **36** - termasuk `.Deductible`, `.AggregateLimit`, `.Amount`, `.Currency` |
> | `TreatyRetroList` | **26** - termasuk `.Deductible`, `.Limit`, `.NetPremi` |
> | `TreatyInTabsProportional` | 9 |
> | `Share`, `ShareRetro` | 5 masing-masing - `.Value`, `.Pct`, `.Currency` |
> | `DetailLimits` | 4 - **`.Value`**, **`.Currency`**, `.Achievement`, `.CurrencyAchievement` |
> | `DetailShare`, `DetailShareRetro`, `MaxRetention` | 2 masing-masing - `.Value`/`.Amount`, `.Currency` |
> | `Layers` | 2 - `.Value`, `.TreatyGroup` |
> | `LayersEDM` | 1 - `.AdjRate` |
>
> **Kesimpulan `DB-10` tidak berubah, tetapi daftarnya berubah:** deductible dan batas agregat
> memang masih dapat diubah - **lewat `TreatyInTabsNonProportional` dan `TreatyRetroList`**, bukan
> lewat `LayersEDM`; nilai dan persentase bagian masih dapat diubah; **total angsuran tidak.**
> **GRL-08 butir iii** ikut memakai daftar yang sudah dikoreksi ini.

**Dipisah menurut pohon layar, dan pemisahan itu TIDAK mengubah kesimpulannya.** `Revision` bekerja
di layar **penawaran**, sedangkan penguncian warisan GRL-08 butir iii terjadi di layar **addendum**.
Maka daftar seksi kedua pohon disusun dari graf penyertaan seksi:

| | Seksi yang dapat dicapai |
|---|---|
| **pohon penawaran** (`InputTreatyInOffer`) | `DetailLimits`, `Layers`, **`LayersEDM`**, `Share`, `ShareRetro`, `DetailShare`, `Installments`, `DetailEGNPI`, `MaxRetention`, `TreatyInTabsNonProportional`, `TreatyInTabsProportional`, `TreatyInNONProportional`, `TreatyRetroList` |
| **pohon addendum** (`InputTreatyInAdjustment`) | **seluruh seksi di atas**, ditambah empat belas seksi `*OldData` |

**`LayersEDM` ada di KEDUA pohon** - dugaan bahwa ia hanya dirender di layar addendum **dibantah**.
Dan tidak ada satu pun seksi yang hanya ada di pohon penawaran: pohon addendum adalah **superset**
pohon penawaran ditambah kolom nilai lama.

Akibatnya dua-duanya:

* **Bunyi `DB-10` tidak berubah.** Seluruh seksi uang yang tidak terjaga - `Installments`, `Share`,
  `ShareRetro`, `LayersEDM`, dan sisanya - berada di pohon penawaran, yaitu tempat `Revision`
  bekerja.
* **GRL-08 butir iii ikut mengecil.** Addendum yang "lahir terkunci" karena mewarisi
  `RevisionState = 1` ternyata **hanya terkunci sebagian**, dengan daftar bocor yang sama.
  Kuncinya bukan gembok; ia penutup di sebagian kontrol.

> **Batas perkakasnya, dilaporkan sebagaimana adanya.** Graf di atas disusun dari penyertaan seksi
> (`pySectionName` dan sejenisnya). Ia **tidak** menangkap seksi yang dibuka sebagai **local
> action** - `ShowSummary` termasuk yang begitu, dan karena itu muncul "tidak terjangkau" pada
> kedua kolom padahal ia dibuka dari layar penawaran (§4.6). Satu seksi lain, `DetailShareRetro`,
> juga tidak tertangkap. Keduanya **bukan lubang pada kesimpulan di atas** - keduanya di luar
> daftar seksi uang yang dijaga - tetapi disebut supaya tidak dibaca sebagai bukti ketidakhadiran.

**Bacaan yang benar, akhirnya:** `ViewState = 1` **bukan** penanda "revisi tanpa uang". Ia penjaga
yang dipasang **sebagian, pada sebagian kontrol, tanpa pola yang dapat dipertanggungjawabkan** —
bentuk yang sama dengan penjaga-penjaga lain di korpus ini.

**Siapa yang melepas kuncinya.** Penulis `ViewState` disapu tertutup atas kelima jenis penulis:

| Aturan | Nilai | Catatan |
|---|---|---|
| `SetTreatyIn_Act` 6, `SetTreatyInEDM_Act` | `1` | dipasang saat `Revision` dan saat addendum dimuat |
| `TreatyInSetEdit` 1 | `0` | melepas - **tetapi** langkah 2 memasangnya kembali bila `RevisionState == 1` |
| `TreatyInForceEdit` | `0` | tombol "(dev)", operator divisi `IT`; **tidak** memeriksa `RevisionState` |
| `TreatyInCopy` 1 | `0` | jalur salin kontrak |
| `Akseptasi_DT` 2.2.1.5 | `""` | dilepas saat revisi selesai disetujui |

Jadi selama revisi berjalan kuncinya bertahan, dan ia dilepas sendiri ketika kepala seksi menyetujui.
Satu-satunya yang melepasnya **di tengah jalan tanpa memeriksa apa pun** adalah `Force Edit (dev)`.

**Akibatnya untuk pertanyaan bisnis, sesudah pencocokan di atas.** Dugaan bahwa rantai dua tingkat
dibenarkan oleh "revisi tanpa uang" **tidak dapat dipertahankan**: angka uang tetap dapat diubah
sesudah `Revision`, di banyak tempat. `DB-10` karena itu ditulis ulang **untuk kedua kalinya**, dan
kali ini tidak mengandaikan kebijakan apa pun - ia menanyakan apakah kebijakannya ada, dan
menyebutkan apa yang sebenarnya dapat berubah.

---

## 8. Layar `OldData`

Dua puluh berkas, seluruhnya cermin baca `TreatyIn.OLDDATA`. Tidak ada satu pun yang menghitung.

| Section | Flow Action | Harness | Kelas |
|---|---|---|---|
| `CoBListOldData` | ya | | `Data-TreatyInLimitLayer` |
| `DetailEGNPIOldData` | ya | | `Data-TreatyInEGNPI` |
| `DetailLimitsOldData` | ya | | `Data-TreatyInLimitsDetail` |
| `DetailShareOldData` | ya | | `Data-TreatyInLimitsDetail` |
| `LayersOldData` | ya | | `Data-TreatyInLimits` |
| `LimitProportionalOldData` | ya | | `Data-TreatyInLimits` |
| `MaxRetentionOldData` | ya | | `Data-TreatyInRetention` |
| `TotalLimitsOldData` | ya | | `Data-TreatyInLimits` |
| `ShareOldData` | ya | | (ada di kedua ekspor) |
| `TreatyInFacultativeShareCalculationOldData` | | ya | `Data-Portal` |
| `TreatyInNONProportionalOldData` | | | `Data-Portal` |
| `TreatyInTabsNonProportionalOldData` | | | `Data-Portal` |
| `TreatyInTabsNonProportionalOldDataShare` | | | `Data-Portal` |
| `TreatyInTabsProportionalOldData` | | | `Data-Portal` |
| `Installments_ReadOnly` | ya | | `Data-TreatyInInstallment` |

Cara mereka dibuat terbaca dari memo aturannya sendiri: *"save as"*, *"Initial Creation"*,
*"added section with all disabled inputs"*, *"readonly untuk olddata"*, *"ganti jadi old data"*.
Yaitu **salinan seksi penyuntingan dengan inputnya dimatikan**, bukan komponen tampilan tersendiri.

Akibatnya untuk migrasi React: **20 berkas ini tidak berpasangan satu-lawan-satu dengan komponen
baru.** Ia satu komponen tampilan yang sama dengan mode `readOnly`, diberi sumber data `OLDDATA`.
Menerjemahkannya berkas-per-berkas akan menggandakan dua puluh komponen tanpa satu pun perilaku
baru.

---

## 9. Utang teknis dan cacat

Disusun menurut aturan pelaporan proyek ini: **arah dampak disebut**, temuan sejenis digabung
jadi satu butir, dan **"masih bisa terjadi"** (terbaca dari aturan) dipisah dari **"pernah
terjadi"** (hanya dapat dijawab data). Tidak satu pun butir di bawah menyatakan dirinya terverifikasi
atas data produksi.

### 9.1 Yang terbaca dari aturan — "masih bisa terjadi"

| # | Butir | Bukti | Arah dampak |
|---|---|---|---|
| **TDA-01** | Nomor revisi dihitung dari baris yang **dipilih**, bukan dari revisi terakhir; penjaga duplikatnya mati; tabrakan berakhir sebagai `UPDATE` yang melaporkan berhasil | `TreatyInRevisi_post` langkah 5; `TreatyInEDMSetValue` langkah 7 blok `//`; `PEGA_M_TREATY_IN_EDM` cabang `ELSE` | **kehilangan data diam-diam** — addendum lain tertimpa, pengguna melihat pesan berhasil |
| **TDA-02** | Penolakan addendum **menghapus barisnya** dari `M_TREATY_IN_EDM` dan `TREATY_IN_EDM`; langkah pencatat penolakan mati | `TreatyInDeclineConfirmation_postactEDM` langkah 1–4 mati, 6–7 hidup | **kehilangan jejak** — tidak ada riwayat penolakan sama sekali |
| **TDA-03** | Penghapusan itu **tidak membersihkan** detail addendum, lampiran salinan, maupun `RevisionState = 1` di baris kontraknya | `RemoveTreatyInDetailEdm` tidak dipanggil di jalur itu; `TreatyRevisionCopyAttachment` tidak punya pasangan hapus | **sampah tertinggal**; kontrak tetap tertandai sedang direvisi |
| **TDA-04** | Selisih dipadankan menurut **posisi baris**, bukan kunci bisnis; panjang daftar tidak pernah dibandingkan | seluruh `TreatyEDMDifference*` | **angka salah tanpa galat** bila baris disisipkan/dihapus di tengah |
| **TDA-05** | Selisih dikurangkan **tanpa memeriksa mata uang** kedua sisi | seluruh `TreatyEDMDifference*` | **angka tak bermakna** bila mata uang berubah antar versi |
| **TDA-06** | Sebagian selisih ditulis ke `ActualValue`, yang lalu ditimpa habis saat menyimpan | `TreatyEDMDifferenceShare` langkah 4 + `SaveTreatyIn_EDM_Act` langkah 2–4 | **hasil hilang** — angka yang terlihat di layar tidak tersimpan |
| **TDA-07** | ~~Membuka kontrak dari picker Adjustment menulis kembali ke baris kontrak.~~ **DICABUT 24 Sep 2026** (§4.5). Yang berdiri menggantikannya: **dua kontrol di layar `InputTreatyInOffer` mengosongkan `StatusAkseptasi` sebuah kontrak lalu menyimpannya**, tanpa versi baru dan tanpa jejak selain satu komentar | `SetTreatyIn_Act` langkah 7–11 bersyarat `param.revisionstate == 1`; hanya kontrol 7 dan 8 `InputTreatyInOffer` yang mengirimnya | **temuan Treaty In, bukan Adjustment**: kontrak yang sudah disetujui dapat kehilangan status akseptasinya |
| **TDA-08** | `Akseptasi_DT` bercabang pada `RevisionState`, bukan pada jenis addendum: `RevisionState = 1` memotong rantai empat tingkat menjadi dua, lalu **nilainya dikosongkan** saat `Resolve Complete`. Untuk addendum, nilai itu **tidak pernah ditulis jalur addendum** — ia diwarisi dari baris yang dimuat (§4.5, §7.2) | `Akseptasi_DT` cabang 2 dan aksi `2.2.1.4`; `SetTreatyIn_Act` langkah 7 hanya menyala lewat kontrol 7-8 `InputTreatyInOffer` | **panjang rantai tidak ditentukan jenis perubahan**, dan **jejaknya hilang** — sesudah selesai, tidak ada yang menyatakan mengapa rantainya pendek. Seberapa sering addendum mewarisi `1` **tidak terbaca**: UA-9 |
| **TDA-09** | Peran dibaca dari `OperatorID.pyWorkBasketList(2)` (indeks tetap) dan `OperatorID.pyTelephone`; nama orang disemat di aturan; ada jalur lolos `pyPosition = "IT Developer"` | `IsTreatyUser`, `Akseptasi_DT` | **wewenang tidak dapat diaudit**; berpindah jabatan menuntut ubah aturan |
| **TDA-10** | Sifat material/non-material hanya ditegakkan sebagai `pyDisabledWhen` di layar; **tidak ada pemeriksaan di sisi penyimpanan** | ratusan `pyDisabledWhen`, nol pemeriksaan di `SaveTreatyIn_EDM_Act` | **aturan dapat dilanggar** lewat jalur mana pun yang bukan layar itu |
| **TDA-11** | Daftar pilihan menyatukan kontrak dan addendum tanpa kolom pembeda, tanpa saringan keadaan; varian non-XOL tidak menyaring `ProportionType` | `TreatyLoadMasterJoinEdm` vs `…XOL` | **salah pilih** tidak tertangkap apa pun |
| **TDA-12** | Jenis baris ditentukan **panjang teks pengenal** (`@length(Param.ID) == 7`), dan nomor revisi diurai dengan **offset yang meleset satu** — `@substring(ID,10,12)` membaca digit satuan saja, padahal digit puluhan ada di indeks 9 | `TreatyInRevisi_post` **baris 848**; `TreatyInEDMSetValue` langkah 1–2 | **patah pada revisi kesepuluh**, dengan dua akhir yang keduanya rusak: **diam** (menimpa `/R02`) atau **berisik** (revisi-dari-revisi selalu galat). Mana yang terjadi **tidak terbaca** — UA-1. Pengenal kontrak yang panjangnya bukan tujuh merusak seluruh percabangan |
| **TDA-13** | Jenis addendum **dan** sifat materialnya dapat diubah pengguna lewat dua radio group di `PickerTreatyInMasterRevisi` — tanpa aturan, tanpa jejak, tanpa kondisi tampil, tanpa kondisi mematikan | `pyReadOnly = false`, `pyDisabledNew = false`, nol `pyVisible`/`pyCondition`/`pyVisibleWhen` sampai akar seksi (§3.2) | **klasifikasi addendum ditentukan pilihan bebas di layar.** Nilai apa saja yang ditawarkan radio **tidak terbaca** — Field Value tidak ter-ekspor; `EDMState = "2"` karena itu **masih bisa terjadi**, belum terbukti |
| **TDA-14** | Nilai disimpan sebagai teks di sebagian tempat — terbukti dari `@toDecimal()` yang dipasang di kedua sisi pengurangan `Deductible2` | `TreatyEDMDifferenceLimits` langkah 1.1 | presisi tidak dapat dipulihkan (penyakit yang sama dengan modul induk) |
| **TDA-15** | Penggandaan data: satu baris addendum memuat pohon kontrak **dua sampai tiga kali** (`TreatyIn`, `OLDDATA`, `ActualValue`) ditambah `ValueDifference` | `SaveTreatyIn_EDM_Act` langkah 2–7 | ukuran `JSONDATA` berlipat; biaya baca/tulis dan penyimpanan |
| **TDA-16** | **Selisih EGNPI pada addendum premi selalu nol.** Pengguna menyunting `TreatyIn.ActualValue.EGNPI`, tetapi `TreatyEDMDifferencePremium` beriterasi atas `TreatyIn.EGNPI` dan mengurangkannya terhadap `TreatyIn.OLDDATA.EGNPI` — keduanya masih nilai kontrak asal | `TreatyEDMDifferencePremium` langkah 1 (`pyStepsObjectName = TreatyIn.EGNPI`) versus `TreatyEDMDifferenceShare` langkah 1 (`= TreatyIn.ActualValue`) | **angka pokok jenis addendum itu tidak pernah terlihat sebagai selisih**; di bawah GRL-12 ia akan terbaca non-material |

### 9.2 Yang hanya dapat dijawab data — "pernah terjadi"

Daftar ini **tidak memblokir rancangan**. Ia memblokir daftar eskalasi ke manajemen dan perkiraan
pekerjaan perbaikan data.

| # | Pertanyaan | Kueri yang menjawabnya |
|---|---|---|
| **UA-1** | Bentuk pengenal yang benar-benar terjadi, **tiga hitungan**: (a) berapa kontrak yang punya lebih dari sembilan addendum, dan bentuk pengenal apa yang dipakai addendum kesepuluh dan seterusnya; (b) berapa pengenal yang **berulang** di dalam satu kontrak; (c) **(diganti 26 September 2026)** berapa addendum yang `OLDID`-nya menunjuk **baris addendum**, dan - yang sebenarnya menandai percabangan - berapa addendum yang `OLDID`-nya **bukan versi terakhir yang ada pada saat ia dibuat**. Butir (c) yang lama, *"pengenal yang memuat lebih dari satu `/R`"*, **dicabut**: pengenal bersusun tidak pernah terbentuk (lihat koreksi §4.4), dan percabangan **tidak terlihat dari bentuk pengenal** - hanya dari `OLDID` | `TREATY_IN_EDM`; (c) membandingkan `OLDID` dengan `EDMDATE` baris-baris sekontrak |
| **UA-2** | Sebaran `EDMSTATE` × `EDMMATERIALTYPE` di produksi — `group by` keduanya | `TREATY_IN_EDM`. **Kepala kueri ini menyebut apa yang dianggap sah, supaya sisanya terbaca sebagai anomali** (`EXP-1`, 26 September 2026): kombinasi yang **dapat dibuat lewat layar hanya lima** — **(1,1), (1,2), (2,1), (2,2)** dari radio picker, ditambah **(3,1)** dari tombol Penyesuaian. Bila kueri menemukan **(3,2)**, atau nilai di luar {1,2,3}, atau kombinasi kosong, itu **anomali data** — bukan jenis yang terlewat dibaca. Kadar ini berlaku untuk versi ruleset 01-01-56 |
| **UA-3** | Berapa addendum non-material (`EDMMATERIALTYPE = 2`) yang nilai uangnya tetap berubah? | bandingkan `JSONDATA` dengan `JSONDATA.OLDDATA` pada baris ber-`EDMMATERIALTYPE = 2` |
| **UA-4** | Berapa kontrak yang `RevisionState`-nya tertinggal `1` tanpa addendum hidup? | `TREATY_IN` × `TREATY_IN_EDM` |
| **UA-5** | Berapa addendum yang daftar barisnya berubah panjang terhadap `OLDDATA`-nya? | `JSON_TABLE` atas kedua sub-pohon, bandingkan cacah baris per daftar |
| **UA-6** | Berapa addendum yang mata uangnya berbeda antara baru dan lama pada besaran yang sama? | idem |
| **UA-7** | Berapa baris `TREATYINDETAILEDM` / `M_TREATY_IN_DETAIL_EDM` yang induknya sudah tidak ada? | anti-join ke `TREATY_IN_EDM` |
| **UA-9** | **(dipertajam 24 September 2026)** Tiga hitungan sekaligus, karena §7.4 kini menyebut sebabnya: (a) berapa **addendum** yang membawa `RevisionState = 1` warisan - sehingga menempuh rantai dua tingkat **dan** lahir dengan `ViewState = 1`; (b) berapa **kontrak** yang membawa `RevisionState = 1` bersama `StatusAkseptasi = "Resolve Complete"` - keadaan beku yang hanya dapat lahir dari `Force Resolve Complete(dev)`; (c) berapa **kontrak** yang membawa `RevisionState = 1` bersama `StatusAkseptasi = "Reject"` - revisi yang ditolak dan dikembalikan, yang menurut aturan **tetap** bertanda sedang direvisi | `TREATY_IN` dan `TREATY_IN_EDM` lewat `JSON_TABLE` atas `JSONDATA`; (b) mengukur pemakaian tombol paksa, (c) memisahkan yang wajar dari yang tersangkut |
| **UA-10** | Berapa addendum yang **kelima** field lapisan beku-nya berbeda dari kontrak induknya — cedant (**nama dan ID diukur terpisah**), source of business, `ProportionType`, tanggal mulai, tanggal berakhir? Dan berapa yang **nama cedant-nya berubah tetapi `CEDINGID`-nya tidak**? | bandingkan `TREATY_IN_EDM` dengan `TREATY_IN` pada keenam kolom; **kelima field diukur**, bukan hanya tiga yang terbukti terbuka di layar — tidak ditemukan terbuka bukan bukti tidak pernah berubah (`METODE` §3.1). Baris yang namanya berubah tanpa ID-nya berubah adalah tanda penyuntingan lewat kontrol yang tidak menulis ID (§4.6) |
| **UA-11** | Berapa kontrak yang status akseptasinya kosong padahal pernah disetujui, atau yang membawa komentar "Create Revision"? | `TREATY_IN` lewat `JSON_TABLE`; mengukur akibat kontrol 7-8 `InputTreatyInOffer` (§4.5) |
| **UA-12** | **(diperluas 24 September 2026)** Mengukur **sejarah**, bukan hanya keadaan terakhir. `Position` hanya menyimpan tahap yang sedang berjalan: baris yang dulu melewati peran kelompok lalu naik ke tingkat berikutnya **tidak lagi** menyimpan nilai itu. Tiga hitungan: (a) baris yang `Position`-nya bernilai `ReasTreatyInGroupLeader` sekarang; (b) baris yang `PositionUsername`-nya bernilai **`"BERNARD"`**; (c) baris yang `CommentList(*).OperatorName`-nya menunjuk orang yang sama. **Bentuk nama berbeda antar penulis dan kuerinya harus mencocokkan bentuk yang sama:** `PositionUsername` diisi dua bentuk - teks yang ditanam di aturan (`"BERNARD"`, `"IRVANDY"`, `"NANDINA"`, `"AGUNGPUTRAANDALAS"`, `"YOHANESKRISTIAWAN"`: tanpa spasi, huruf besar, berciri **ID operator**) dan `OperatorID.pxInsName` (`SetTreatyIn_Act` langkah 7, `TreatyInSetEdit` langkah 5) yang juga ID operator; sedangkan `CommentList(*).OperatorName` diisi `OperatorID.pxInsName` pada `SetTreatyIn_Act` langkah 9 tetapi **dibaca balik** ke `PositionUsername` oleh `Akseptasi_DT`. Kueri (c) karena itu mencocokkan **ID operator**, bukan nama lengkap | `TREATY_IN` dan `TREATY_IN_EDM` lewat `JSON_TABLE`. **Batas yang harus disebut:** `CommentList` hanya punya empat ruas - `OperatorName`, `IsApproved`, `Suggest`, `Date` - dan **tidak merekam tingkat penyetuju**. Maka jejak persetujuan tingkat kelompok hanya dapat dikenali lewat **nama orangnya**, dan ekspor menyebut tepat satu nama. **Kadar pembuktiannya terbatas, dan itu ditulis di kepalanya sendiri - sama seperti `UA-3`:** hitungan (b) dan (c) dapat **mematahkan** dugaan *"jalur kelompok tidak pernah dipakai"*, tetapi **tidak dapat mengesahkan** bahwa jalur itu pernah dipakai - orang yang sama dapat menyetujui di tingkat lain, dan `CommentList` hanya punya empat ruas (`OperatorName`, `IsApproved`, `Suggest`, `Date`) sehingga **tidak merekam tingkat penyetuju**. Hanya hitungan (a) yang berbicara langsung tentang keadaan. Bila (a) tidak nol, itu keadaan warisan tak terpetakan - ADR-0054, cabang I |
| **UA-13** | Berapa kali revisi-di-tempat dijalankan per bulan, dan pada berapa kontrak berbeda? | `TREATY_IN` lewat `JSON_TABLE`: hitung baris yang `CommentList`-nya memuat `Suggest = "Create Revision"` - satu-satunya jejak yang ditinggalkan `SetTreatyIn_Act` langkah 9 - dikelompokkan menurut `Date`. **Tidak memblokir apa pun**: ia memperkirakan tambahan beban penyetuju sesudah GRL-07 butir 4 menghapus pemendekan rantai, dan menentukan siapa yang harus diberi tahu sebelum cut-over (cabang D) |
| **UA-14** | Berapa addendum yang **isinya tertimpa** akibat tabrakan pengenal (TDA-01), dan mana saja? | `M_TREATY_IN_EDM` lewat `JSON_TABLE`: (a) baris yang `RevisionDate`-nya **lebih baru** daripada baris bernomor lebih tinggi pada kontrak yang sama - `TreatyInEDMSetValue` langkah 3 menulis `RevisionDate = @CurrentDateTime()` pada setiap pembuatan, jadi urutan waktu yang terbalik adalah tanda; (b) baris ber-`StatusAkseptasi` bukan draf yang `CommentList`-nya **kosong** - langkah yang sama mengosongkannya saat pembuatan, sehingga addendum yang pernah disetujui seharusnya membawa jejak. **`OLDID` tidak dapat dipakai** - ia hanya menunjuk induk dan tidak berubah saat tertimpa. **Batasnya:** isi yang hilang itu sendiri **tidak terukur** - ia tertimpa tanpa salinan. Yang terukur hanyalah **berapa baris yang tertimpa**, bukan apa yang hilang |
| **UA-15** | Adakah operator yang `pyUserName`-nya persis **`ALDO SAPUTRA1`**, dan adakah yang `pyUserIdentifier`-nya **`aldo_saputra`**? | tabel operator Pega (`Data-Admin-Operator-ID`). Yang pertama menentukan apakah layar `ShowSummary` - satu-satunya penyunting `CedingID` - **terjangkau seseorang**; yang kedua untuk tombol `TestCopyDifference`. **Ini pertanyaan data, bukan ekspor** (`METODE` §3.1, §4.6): nama berakhiran angka lazim dipakai akun uji, jadi ketiadaannya tidak dapat disimpulkan dari bentuk namanya |
| **UA-16** | Berapa addendum yang **disimpan sesudah disetujui**, yaitu disunting di tempat lewat `Force Edit (dev)` + `Save EDM(dev)`? | **Tanda (b) adalah tanda utama; tanda (a) hanya perkiraan kapan.** Sebabnya: satu `CommentList` dapat memuat entri berzona **Pega (GMT)** dan entri berzona **Oracle** sekaligus, dan langkah jam Oracle di `AddCommentList_Act` yang kini **mati** menunjukkan zonanya **pernah berganti** - kapan pergantiannya **tidak terbaca**. Maka setiap baris yang dibuat sebelum pergantian itu punya pembanding berzona berbeda dari yang sesudahnya, dan tidak ada cara memisahkannya dari bahan yang ada. Dua tanda, dengan kadar yang berbeda: **(a)** `TREATY_IN_EDM.EDMDATE` diperbarui `SYSDATE` pada **setiap** penyimpanan (`PEGA_M_TREATY_IN_EDM` cabang `ELSE`). **Zona waktunya harus disamakan lebih dulu, dan kedua sisinya BERBEDA.** `EDMDATE` diisi `SYSDATE` - **jam server Oracle**. `CommentList.Date` pada komentar **persetujuan** diisi `AddCommentList_Act` langkah 2 dan 3 dengan **`@CurrentDateTime()`** - **Pega, GMT**. Sedangkan `CommentList.Date` pada komentar **"Create Revision"** diisi `SetTreatyIn_Act` langkah 9 dengan `@toDateTime(time.pxResults(1).CARI1)` dari `GetCurrentDate` = `select sysdate from dual` - **jam Oracle lagi**. Jadi satu `CommentList` dapat memuat dua entri berzona berbeda, dan jejak peralihannya masih terlihat: `AddCommentList_Act` langkah 1 adalah `RDB-List` ber-blok `//`, panggilan jam Oracle yang **dimatikan** saat penulisnya berpindah ke `@CurrentDateTime()`. **Kuerinya karena itu harus menormalkan dulu** - misalnya `FROM_TZ(CAST(EDMDATE AS TIMESTAMP), <zona server>) AT TIME ZONE 'UTC'` dibandingkan dengan `CommentList.Date` yang dibaca sebagai GMT. **Zona server Oracle tidak terbaca dari bahan yang ada** dan ditandai **perlu dicek** ke DBA; tanpa penyamaan itu toleransi lima menit tidak berarti, sebab selisih beberapa jam akan menandai hampir setiap baris. **Dan pembandingnya tetap harus bertoleransi, sebabnya terbaca:** `TreatyInSubmitEDM` langkah 6 menambah komentar, langkah 7 menghitung selisih, langkah 8 memanggil `SaveTreatyIn_EDM_Act` - jadi **persetujuan itu sendiri menyimpan barisnya**, dan `EDMDATE` selalu **sedikit lebih baru** daripada tanggal komentar persetujuannya, dalam transaksi yang sama. Maka ujinya: `EDMDATE > MAX(CommentList(*).Date) + 5 menit`. Toleransi lima menit jauh melampaui jarak satu transaksi dan jauh di bawah jarak penyuntingan manusia; **(b)** `SaveTreatyIn_EDM_Act` **tidak memanggil mesin selisih**; `TreatyEDMCalculateDifference` dipanggil `TreatyInSubmitEDM` **langkah 7** dan kontrol tab. Karena setiap pengajuan menghitung ulang selisihnya, baris yang disunting lalu disimpan lewat `Save EDM(dev)` akan punya `ValueDifference` yang **tidak sama** dengan nilai baru dikurangi `OLDDATA`. **Tanda ini tidak perlu toleransi waktu**, dan karena itu lebih kuat daripada (a); (a) dipakai untuk memperkirakan **kapan**, (b) untuk memastikan **apakah** |
| **UA-17** | Berapa kontrak yang punya **revisi-di-tempat sesudah addendum terakhirnya disetujui**, dan berapa yang urutannya terbalik? | `TREATY_IN` dan `TREATY_IN_EDM` lewat `JSON_TABLE`: bandingkan tanggal komentar **"Create Revision"** pada baris kontrak dengan `RevisionDate` addendum terakhirnya. **Ia menjawab pertanyaan yang tidak dapat dijawab bentuk data:** di sistem lama baris `TREATY_IN` tidak pernah digantikan addendum, tetapi **dapat disunting di tempat** lewat `Revision` - termasuk sesudah addendum terakhir dibuat. Urutan antara keduanya hanya dapat direkonstruksi dari cap waktu. Menentukan arti **operasional** "versi berlaku" untuk kontrak warisan (GRL-11, titipan I2, `DB-12`). **Batas:** kedua cap waktu itu berasal dari jam yang sama (Oracle `sysdate`), sehingga perbandingannya sah - tidak seperti `UA-16` |
| **UA-18** | Untuk **setiap entitas anak**, berapa daftar di `JSONDATA` yang memuat **dua baris dengan `KUNCI_PADANAN` yang sama**? | `M_TREATY_IN` dan `M_TREATY_IN_EDM` lewat `JSON_TABLE`, satu kueri per daftar. **Kunci padanan hanya bekerja bila ia unik di dalam daftarnya**, dan keunikan itu **pertanyaan data, bukan rancangan**: `SEAM-ADJUSTMENT.md` §3 menetapkan bentuknya, tidak menjamin data warisan mematuhinya. Bentuk majemuk ikut diuji - untuk `DETAIL_PROPORSIONAL` berarti (nomor layer + bagian layer) + kelompok treaty, untuk `POTONGAN` ditambah penunjuk induk mana. Bila ada yang tidak unik, pemadanan selisih **akan memasangkan baris yang salah** pada baris warisan, dan itu milik cabang I |
| **UA-8** | Berapa panjang rata-rata dan maksimum `JSONDATA` di `M_TREATY_IN_EDM`, dan berapa bagiannya yang `OLDDATA` + `ActualValue`? | pengukuran untuk perkiraan ukuran basis data baru |

**UA-3 punya sifat pembuktian yang terbatas dan itu harus disebut di kepalanya sendiri:** ia dapat
**mematahkan** dugaan "non-material tidak pernah mengubah nilai", tetapi tidak dapat
**mengesahkannya** — penyuntingan yang ditimpa tidak meninggalkan jejak.

### 9.3 Butir parkiran sesi sebelumnya, ditutup

`TEMUAN-ADJUSTMENT-DITUNDA.md` memarkir empat butir khas Adjustment. Ketiganya kini terjawab:

| Butir parkir | Nasib |
|---|---|
| **A-1** nomor revisi dari potongan teks | berdiri; menjadi bagian **TDA-12** |
| **A-2** `ROWNUM = 1` mendahului `ORDER BY` | **bobotnya turun** — `HASIL1` hanya diuji kosong/tidak, nilainya tidak pernah dipakai (§4.4). Cacat yang sebenarnya ada di tempat lain: **TDA-01** |
| **B-1** picker meng-UNION tanpa pembeda | berdiri; menjadi **TDA-11**, ditambah temuan saringan `ProportionType` yang timpang |
| **D-1** `PositionUsername` di kelas addendum | berdiri; sama seperti di modul induk, dan keputusan modul induk (dibuang) berlaku wajar di sini |
| **G4** pengenal dibentuk dari teks, offset tetap | berdiri; menjadi **TDA-12** |

Dan pertanyaan G4 yang sengaja ditunda — *"untuk revisi kedua dan seterusnya, apakah `OLDID`
menunjuk kontrak atau revisi sebelumnya?"* — kini **terjawab dari aturan**:
`TreatyInRevisi_post` langkah 5 menetapkan `TreatyIn.OLDID <- TreatyIn.ID` **sebelum** `ID` diubah,
sehingga `OLDID` menunjuk **baris yang dipilih pengguna**. Bila yang dipilih revisi sebelumnya,
`OLDID` menunjuk revisi itu; bila yang dipilih kontrak, ia menunjuk kontrak. **Rantainya karena itu
tidak seragam**, dan bentuknya ditentukan pilihan pengguna, bukan aturan.

---

## 10. Akibat untuk migrasi React / Golang / Oracle

Bagian ini **bukan rancangan** dan tidak menetapkan apa pun. Ia mencatat hal-hal yang, bila
terlewat, akan menghasilkan rancangan yang salah.

### 10.1 Yang dikoreksi dari masukan rancangan yang sudah ada

`SEAM-ADJUSTMENT.md` dibangun di atas anggapan bahwa **sistem lama tidak menyimpan nilai lama**,
sehingga nilai lama harus di-SELECT dari versi terdahulu. §2.1 menunjukkan yang sebaliknya: sistem
lama **menyimpan potret penuh** di `OLDDATA`.

Itu **tidak membatalkan** arah seam-nya — memilih rujukan alih-alih salinan tetap keputusan
rancangan yang sah, dan ADR-0036 tetap berdiri. Yang berubah adalah **sifat pekerjaan migrasinya**:

* migrasi **tidak perlu merekonstruksi** nilai lama dari tabel yang bergerak — ia sudah ada, beku,
  di dalam `JSONDATA` tiap addendum;
* karena itu nilai lama yang dimigrasikan **dapat dicocokkan** dengan nilai versi terdahulu yang
  hasil migrasinya, dan **ketidakcocokan menjadi alat deteksi** mutu migrasi;
* dan `OLDDATA` **tidak dibawa** sebagai bentuk tersimpan, karena di model versi ia rangkap.

### 10.2 Yang harus ada di model baru dan tidak ada padanannya di sistem lama

| Hal | Karena |
|---|---|
| Kunci padanan baris lintas versi | TDA-04: sistem lama memadankan menurut posisi |
| Pengenal versi buatan sistem | TDA-12: pengenal lama dibentuk dari teks dan diurai dengan offset tetap |
| Penolakan sebagai **keadaan**, bukan penghapusan | TDA-02: menolak berarti menghapus baris |
| Pemeriksaan material/non-material di sisi penyimpanan | TDA-10: sekarang hanya kondisi tampilan |
| Mata uang sebagai bagian tak terpisahkan dari nilai | TDA-05 |

### 10.3 Yang menyusut, bukan berkembang

* **20 berkas `OldData` bukan 20 komponen React.** Ia satu komponen dengan mode baca dan sumber
  data berbeda (§8).
* **Tidak ada entitas "Adjustment" yang perlu dibuat** — kesimpulan G1 sesi sebelumnya bertahan,
  dan sekarang bertambah bukti: 19 properti kelas `Int-treaty_in_edm` seluruhnya lapisan kepala
  kontrak ditambah tiga penanda (`OLDID`, `EDMState`, `EDMMaterialType`).
* **Mesin selisihnya satu bentuk, bukan lima.** Kelima aktivitas `TreatyEDMDifference*` menjalankan
  rumus yang sama atas pohon yang berbeda; satu fungsi generik yang menelusuri dua pohon dan
  memadankan barisnya menggantikan seluruhnya.

### 10.4 Angka untuk perkiraan

| Hal | Angka |
|---|---|
| Berkas khas Adjustment | 56 dari 379 |
| Aturan yang benar-benar logika (bukan layar, bukan platform) | 12 Activity + 4 Data Transform |
| Akar `ValueDifference` yang dicerminkan | **33** (21 di antaranya `Total*` agregat), **165 jalur daun** — daftarnya §5.6 |
| Kolom kepala tabel datar addendum | 24 |
| Parameter prosedur simpan | 25 masuk, 3 keluar |

> Catatan atas angka 33: sesi sebelumnya menghitung **31 akar / 161 jalur** dari ekspor Treaty In.
> Dihitung ulang atas ekspor Adjustment hasilnya **33 akar / 165 jalur daun**; selisih akarnya
> `LimitShareSummaryList` dan `LimitFacShareSummaryList`. Kedua angka dihitung dengan skrip,
> bukan disalin — dan **daftarnya** ada di §5.6, karena hitungan yang tidak pernah dituangkan
> sebagai daftar belum ditinjau siapa pun.
