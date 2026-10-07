# KEADAAN — EDM Treaty In

## Migrasi Treaty Inward — Realisasi & Endorsement · modul *EDM Treaty In*

> **Apa berkas ini.** Keadaan terukur modul `EDM Treaty In` per 22 September 2026, sebagai sumber
> tunggal bagi ronde `to-spec`. Setiap angka dihitung **dua cara berbeda**, dan perintahnya ada di
> Bab 8.
>
> ⛔ **Bukan** berkas ronde. `grilling-ronde-1.md` **tersegel dan tidak disentuh**; koreksinya hidup
> di Bab 6 berkas ini.
>
> ⛔ **Berkas ini tidak menutup satu pun butir `[terbuka]`.**

---

## 1 · Urutan kewenangan

Bila dua sumber berselisih, yang di atas menang.

| # | Sumber | Keterangan |
| ---: | --- | --- |
| **1** | `[keputusan work owner]` | keputusan yang sudah diambil — **tidak boleh diubah berkas ini** |
| **2** | `[data DBA]` | naskah stored procedure, view, dan contoh `DATA_JSON` sungguhan |
| **3** | `[terverifikasi]` dari korpus | terbaca dari berkas ekspor, dengan perintah ujinya |
| **4** | `modul/nbtreatyin/docs/KEADAAN-NB-TREATY-IN.md` *(jalur — koreksi 06-10)* | ketetapan modul saudara — **berlaku di sini kecuali terbukti sebaliknya** |
| **5** | `grilling-ronde-1.md` EDM | ⚠️ **empat pernyataannya sudah terbukti keliru** — lihat Bab 6 |
| ⛔ | `[dugaan]` | ⛔ **tidak dipakai sebagai dasar keputusan mana pun** |

---

## 2 · Angka yang berlaku — seluruhnya diuji ulang dua cara

### 2.1 Berkas dan ukuran

| Ukuran | Nilai | Cara A | Cara B | |
| --- | ---: | --- | --- | :---: |
| berkas `.xml` | **163** | `glob **/*.xml` | `os.walk` seluruh isi → **164 berkas**, 12 folder | ✅ |
| ukuran modul | **17.639.259 B** | jumlah `getsize` ke-163 `.xml` | md5 gabungan `24afd5726e083c9df6a692da94728139` | ✅ |
| tipe rule | **12** | nama folder induk | cacah berkas per folder | ✅ |

⚠️ **Berkas ke-164 bukan XML:** `Struktur_InputAddendumTreatyIn.xlsx` *(118.333 B)* — lembar kerja,
bukan aturan Pega. **Tidak masuk hitungan aturan mana pun.** `[terverifikasi]`

> ⛔ **KOREKSI 06-10-2026 — korpus berubah.** Byte kini **17.640.859** (+1.600); md5 gabungan lama tidak
> tereproduksi. Dua berkas berubah sesudah 22-09: `Activity/CountSpreading_Act.xml` (22-09 17:50) dan
> `Activity/InsetTreatyInProdAddendum_Act.xml` (23-09). Cacah 163 / 12 tipe / 380 / 1.309 **tetap** — Bab 11.

| Tipe rule | Jumlah | | Tipe rule | Jumlah |
| --- | ---: | --- | --- | ---: |
| `Activity` | **66** | | `FlowAction` | 6 |
| `RDBList` | **36** | | `Harness` | 3 |
| `Section` | 18 | | ⭐ `ConnectREST` | **1** |
| `When` | 11 | | `DecisionTable` | 1 |
| `DataTransform` | 10 | | `Flow` | **1** |
| `ReportDefinition` | 9 | | ⭐ `SystemSettings` | **1** |

### 2.2 Langkah dan isinya

| Ukuran | **Yang berlaku** | Cara A | Cara B | |
| --- | ---: | ---: | ---: | :---: |
| langkah seluruh metode | **618** | 618 | — | |
| langkah `Property-Set` | ⭐ **380** | **380** | **380** | ✅ |
| `Property-Set` **membawa isinya** | ⭐ **379** *(99,7 %)* | 379 | — | |
| `Property-Set` kosong | **1** | `FillPaymentInstallment` langkah `4.3.2` | — | |
| pasangan nama=nilai terisi | ⭐ **1.309** | **1.309** | **1.309** | ✅ |
| `PropertiesName` terisi | 1.310 | — | pola teks | |
| `PropertiesValue` terisi | 1.311 | — | pola teks | |

**Cara A** — `ElementTree`, `pySteps` **puncak saja**, rekursi ke sub-langkah, pasangan hanya
dihitung bila `PropertiesName` dan `PropertiesValue` berada dalam **satu blok `rowdata` yang sama**.
**Cara B** — pola teks murni, tanpa `ElementTree`.

**Metode langkah teratas:** `PROPERTY-SET` 380 · `RDB-LIST` 60 · `PAGE-NEW` 21 ·
`PROPERTY-SET-MESSAGES` 18 · `PAGE-REMOVE` 15 · `PROPERTY-REMOVE` 12 · `PAGE-SET-MESSAGES` 10 ·
`CALL COUNTNETPREMI_ACT` 9 · `PAGE-CLEAR-MESSAGES` 8 · `OBJ-SAVE` 6.

⭐ **Seluruh 1.309 pasangan berada di langkah `Property-Set`.** Nol di metode lain.

**Kedalaman sarang langkah `Property-Set`:** tingkat 1 = 165 · 2 = 83 · 3 = 65 · 4 = 29 · 5 = 20 ·
6 = 18. ⚠️ **Enam tingkat** — pembaca yang hanya menelusuri satu tingkat akan kehilangan 215 langkah.

### 2.3 ⛔ Selisih 378/380 dan 1.062/1.309 — **dijelaskan, bukan dipilih**

Brief memuat angka **378** dan **1.062** dari sesi terdahulu, dengan peringatan jangan dipercaya.
Saya mencoba **enam rekonstruksi** untuk menemukan pola yang menghasilkannya:

| Rekonstruksi | Hasil |
| --- | ---: |
| blok `rowdata` non-greedy, satu pasangan per blok *(jebakan #2)* | 1.309 |
| idem, blok penyunting **tidak** dibuang | 1.309 |
| hanya nama berawalan `.` / `Local.` / `Param.` / `pyWorkPage` | 925 |
| hanya nilai yang bukan angka harfiah | 1.175 |
| bersebelahan tanpa spasi/baris-baru | 0 |
| satu pasangan per langkah | 379 |
| — untuk langkah: teks `<pyStepsActivityName>` · gadget tidak dibuang · blok non-greedy | **380 · 380 · 380** |

⛔⛔ **Tidak satu pun menghasilkan 378 atau 1.062.** Bahkan pola yang sengaja saya rusak tetap
memberi 380 dan 1.309.

⭐ **Kesimpulan: 378 dan 1.062 tidak dapat direproduksi dan ditarik dari peredaran.**
⚠️ **Sebab pastinya tidak dapat saya nyatakan** — sesi yang menghasilkannya tidak menuliskan
perintahnya. ⛔ **Tidak dikarang.** Yang berlaku: **380** dan **1.309**, sepakat dua cara.

### 2.4 Penjangkauan — ⭐ **nol yatim**

| Ukuran | **EDM** | NB *(pembanding)* |
| --- | ---: | ---: |
| `Activity` | **66** | 92 |
| titik masuk *(`Flow` · `Section` · `Harness` · `FlowAction` · `DataTransform`)* | 34 | 36 |
| ⭐ **terjangkau / yatim — CARA A** *(urai `Call`)* | ⭐ **66 / 0** | 64 / 28 |
| terjangkau / yatim — CARA B *(nama disebut >1 kali)* | **66 / 0** | 92 / 0 |
| langkah `Property-Set` terjangkau / yatim | ⭐ **380 / 0** | 377 / 561 |

⚠️ **Untuk NB, kedua cara berselisih** *(64/28 lawan 92/0)*, sebab nama aturan muncul juga di
metadata ekspor. **Cara A yang dipakai untuk NB.** ⭐ **Untuk EDM keduanya sepakat**, sehingga
`66 / 0` berdiri kuat.

---

## 3 · Lingkup — ⭐⭐ **EDM tidak akan menyusut seperti NB**

⭐⭐ **Nyatakan terang-terangan supaya tidak ada yang mengharapkan penyusutan serupa.**

| | NB Treaty In | **EDM Treaty In** |
| --- | ---: | ---: |
| `Activity` seluruhnya | 92 | **66** |
| ⛔ yatim, dibuang | **28** *(−30 %)* | ⭐ **0** |
| langkah `Property-Set` seluruhnya | 938 | **380** |
| ⛔ langkah yatim, dibuang | **561** *(−60 %)* | ⭐ **0** |
| langkah yang diminta/dikerjakan | 268 *(−72 %)* | ⭐ **380** *(−0 %)* |

⭐ **Folder NB adalah *dependency closure* yang memuat aturan milik modul lain; folder EDM tidak.**
Setiap aturan EDM terjangkau dari titik masuk nyata. ⇒ **Lingkup EDM adalah lingkup penuh.**

### Tiga hal yang tidak dimiliki NB

| Hal | Bukti |
| --- | --- |
| ⭐ tipe rule **`ConnectREST`** | `convertJsonNusareToProduction` — 14.283 B, kelas `ASM-FW-GISFW-Work` |
| ⭐ tipe rule **`SystemSettings`** | `LinkService` — 6.677 B, `LINKSERVICE!LINKSERVICE` |
| ⭐ **kelas kerja sendiri** | `Flow\InputAddendumTreatyIn` berkelas `ASM-FW-GISFW-Work-EndorsementTreaty` |

⭐ **Endorsemen adalah jenis kasus tersendiri**, bukan tahap di dalam kasus polis baru.
`[terverifikasi]` Alurnya *(202.269 B)* memuat **8 kotak keputusan**, **4 penugasan**,
**2 utilitas**, ~~**22 sambungan**~~ **24 baris sambungan** *(koreksi 06-10: 23 `TransitionN` + 1 sisa penyuntingan)*, dan menyebut **tiga antrean**:
`ReasTreatyInAdmin` · `ReasTreatyInSecHead` · `ReasTreatyInDeptHead`.

⇒ ⭐ **Sejalan dengan keputusan tangga tiga jenjang pada P13** — ~~tidak ada penyimpangan~~.
⛔ *Koreksi 06-10 (XML):* **ada** — `Decision9` mengirim berkas langsung ke `Utility1` bila `LetterNo` ≠
`"TREATYINDEPTHEAD"` (≤ 200 juta, `CekLimitTreatyAcc_Act` langkah 2.4), sehingga Sec Head dapat menyelesaikan
sendiri. Tiga jenjang tetap = **keputusan** (penyimpangan sadar) — Bab 11.

---

## 4 · Sumber data

| Ukuran | Nilai |
| --- | ---: |
| `RDBList` | **36** |
| ⭐ punya naskah SQL `pyBrowseSQL` | ⭐ **36 dari 36** *(100 %)* |
| pernyataan `FROM` terbaca | 31 |
| stored procedure dipanggil | **6** |

**Tabel dan view teratas:** `POOLDATA.JSON_POLIS` 4 · `POOLDATA.M_TREATY_IN_EDM` 3 ·
`DUAL` 5 · `POOLDATA.M_TREATY_IN` · `POOLDATA.M_TREATY_OUT` · `POOLDATA.TREATYINPRODUCTION` ·
`POOLDATA.TREATYGROUP`.

### ⛔⛔ Dua stored procedure yang naskahnya BELUM ada

| Procedure | Naskah |
| --- | --- |
| `POOLDATA.PEGA_TREATY_IN` | ✅ diterima — **P1** |
| `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | ✅ diterima — **P1** |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | ✅ diterima — **P1** |
| `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` | ✅ diterima — **P1** |
| ⛔ **`POOLDATA.INSERTJSONPOLISMONITORING`** | ⛔ **belum pernah diminta** |
| ⛔ **`POOLDATA.INSERTUPDATEACHIEVMENT`** | ⛔ **belum pernah diminta** |

⭐ `[terbuka]` **Dua procedure ini dipanggil EDM tetapi tidak termasuk keempat yang diminta P1.**
`[terverifikasi]` ⛔ **Tidak ditutup, tidak diangkat menjadi pertanyaan bernomor di sini** —
dicatat supaya tidak hilang sebelum `to-spec` dijalankan. Pemiliknya `[data DBA]`.

---

## 5 · Keputusan yang mengikat

### 5.1 Dua belas ketetapan NB berlaku juga di sini

Seluruh ketetapan `KEADAAN-NB-TREATY-IN.md` Bab 4 berlaku di EDM **kecuali terbukti sebaliknya** —
`IsApproved` sebagai teks · penolakan bertingkat · uang berpresisi penuh · empat medan persentase ·
tanggal dibaca apa adanya · skema `POOLDATA.` eksplisit · satu transaksi · `OPERATORID` dari
identitas login · pencarian nomor urut antrean tidak dimigrasi · syarat yang dijalankan menang atas
keterangannya · medan terkunci dan wajib · `TypeTax` ditiru apa adanya.

### 5.2 ⚠️ Satu ketetapan yang EDM **langgar** — ini temuan

⛔ **Ketetapan 6** berbunyi *"setiap query menulis skema `POOLDATA.` eksplisit"*.
`[terverifikasi]` EDM **tidak konsisten**: dari 31 pernyataan `FROM`, hanya **16** memakai
`POOLDATA.`; sisanya merujuk tabel **tanpa skema**:

| Tabel tanpa skema | Kali | Contoh berkas |
| --- | ---: | --- |
| `JSON_POLIS` | 3 | `GetPolicyNoByCaseId` |
| `CURRENCY` | 2 | `GetCurrency` · `GetDataCurrencyByName_SQL` |
| `REINSURANCETYPE` | 2 | — |
| `TREATYINPRODUCTION` | 1 | `GetDataTreatyInProd_SQL` |
| `BUSINESS` | 1 | `GetOldIDBusiness_SQL` |
| ⚠️ `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | 1 | `GetCountClaim` — **skema lain**, bukan `POOLDATA` |

⭐ **Ketetapan 6 tetap berlaku untuk sistem baru** — ia ketetapan penulisan, bukan pemerian sistem
lama. ⚠️ Yang perlu diketahui: **`DATAPEGA` adalah skema kedua** yang disentuh modul ini, dan
`JSON_POLIS` dirujuk dengan **dua ejaan** *(dengan dan tanpa skema)* — persoalan yang sama dengan
**P3** di NB. ⛔ `[terbuka]`, **tidak ditutup di sini**.

### 5.3 Keputusan EDM sendiri — P50 sampai P60, **seluruhnya terjawab**

| Butir | Ketetapan |
| --- | --- |
| **P50** | EDM tidak punya FacOut; pengiriman kedua **tidak pernah berjalan** *(⚠️ koreksi 06-10: premis tidak terbukti — `IsFacRetro` diset di jalur EDM; Bab 11)* |
| **P51** | data produksi **dihapus** bila konversi gagal — ⭐ **bersyarat**: hanya bila `STS_KONVERSI` belum `1` |
| **P52** | langkah 10 dihidupkan, langkah 18 tetap mati |
| **P53** | pengiriman ke produksi tanpa autentikasi — **benar**, jalur internal |
| **P54** | `TANGGAL_CLOSING` dibaca `WHERE ROWNUM = 1` — **satu baris berlaku global** |
| **P55** | penomoran satu jalur; adendum dan premi tambahan diikutkan |
| **P56** | pembatalan **jenis endorsemen**, dipilih di awal |
| **P57** | endorsemen **berlapis**; selisih dihitung terhadap keadaan **tepat sebelumnya** |
| **P58** | data lama **disalin dan beku** saat berkas dibuat |
| **P59** | lingkungan uji dan produksi memakai basis data berbeda |
| **P60** | ⭐ ditutup **dari korpus** — lihat Bab 6.5 |

---

## 6 · Koreksi atas ronde 1

⛔ `grilling-ronde-1.md` **tersegel dan tidak disunting.** Bunyi lamanya dikutip di sini, tidak
dihapus.

### 6.1 ⛔⛔ *"Isi langkah penetapan nilai tidak terekspor"* — **KELIRU**

> Bunyi ronde 1: ⛔ *"nilai `Property-Set` tidak terekspor"*

⭐⭐ **Isinya ADA**, di tag **`PropertiesName`** dan **`PropertiesValue`** — **tanpa awalan `py`** —
di dalam `pyParamArray` milik tiap langkah. `[terverifikasi]`

| | EDM |
| --- | ---: |
| `Property-Set` | **380** |
| ⭐ membawa isinya sendiri | ⭐ **379** *(99,7 %)* |
| pasangan terisi | **1.309** |
| `pyPropRef` terisi | ⛔ **0** *(437 tag, semuanya kosong)* |
| `pyPropertiesName` / `pyPropertiesValue` terisi | ⛔ **0 / 0** |

⛔ **Sebabnya dua huruf `py`.** Ronde 1 memeriksa varian **dengan** awalan, menemukan nol, dan
mempercayainya. Rinciannya di ~~`..\nb-treaty-in\VERIFIKASI-P18.md`~~ `modul/nbtreatyin/docs/VERIFIKASI-P18.md` *(jalur — koreksi 06-10)*.

### 6.2 ⛔ *"Pemeriksaan keberhasilan hanya memeriksa FacIn"* — **KELIRU**

> Bunyi ronde 1: ⛔ *"pemeriksaan keberhasilan hanya memeriksa FacIn"*

`When\IsSuccessHitService.xml` memeriksa **keduanya**, dalam **tiga cabang**:
`FacIn = 1` **dan** *(`FacOut = ""` **atau** `FacOut = 1`)*.

⛔ **Sebabnya: ronde 1 membaca `pyConditionString`** — keterangan untuk manusia — bukan bentuk
bertanda kurung yang dijalankan. ⭐ **Itu persis jebakan P23**, dan ketetapan 10 sudah menetapkan
**yang dijalankan benar**.

### 6.3 ⚠️ *"Beda hanya metadata ekspor"* untuk 69 aturan — **sebagian benar; angka brief tidak tereproduksi**

`[terverifikasi]` Dari **73** aturan senama di kedua modul, diukur **dua cara**:

| Cara | Hasil |
| --- | --- |
| **A** — tanda tangan **perilaku** *(metode langkah + pasangan nama=nilai + gerbang + naskah SQL + syarat bertanda kurung; nol metadata)* | ⭐ **69 identik · 4 BEDA** |
| **B** — seluruh tag keterangan *(`pyStepsDescription`, `pyRuleDescription`, `pyNote`, `pyUsage`, `pyComments`, `pyLongDescription`, `pyStepsBlockName`)* | **3 berkas beda** |

**Empat yang beda perilaku:** `CountSpreading_Act` · `SetCategoryAttach` ·
`serviceInsertArasapas_act` · `IsUW`.

⚠️ **Koreksi atas brief.** Brief menyatakan *"enam berkas punya catatan pengembang berbeda, empat di
antaranya di dalam kelompok 69 itu"*. `[terverifikasi]` Terukur: **3 berkas**, dan **ke-3-nya berada
di dalam kelompok 4 yang beda perilaku** — ⭐ **nol** di dalam kelompok 69.

⇒ ⭐ **Kelompok 69 itu benar-benar identik**, baik perilaku maupun catatan. Pernyataan ronde 1
*"beda hanya metadata ekspor"* **bertahan untuk ke-69-nya**; yang keliru hanyalah bila ia diperluas
ke seluruh 73.

### 6.4 ⛔ *"`BusinessType_DeT` tinggal 9 baris dari 34"* — **KELIRU**

**36 baris utuh.** `pyRowNum` pada `pyOrConditions` berbasis **nol** dan menandai baris yang punya
daftar OR — **bukan** nomor baris.

⭐⭐ **Terverifikasi ulang dari data produksi:** contoh `DATA_JSON` ber-`GroupPanel = 006` dan
`QuotationData.BusinessOldId = 01` disimulasikan pada tabel yang direkonstruksi — baris 29 tidak
cocok, baris 30 tidak cocok, **baris 31 cocok** lewat daftar OR-nya, menghasilkan `"FireStyle2"` —
dan dokumennya memang berisi `"FireStyle2"`. `[data DBA]` `[terverifikasi]`

### 6.5 ✅ P60 — **jawabannya sudah tercatat, dan angkanya saya periksa ulang**

| Yang harus benar | Terukur | |
| --- | --- | :---: |
| 8 langkah EDM lawan 7 NB | **8 lawan 7** *(langkah bermetode)* | ✅ |
| langkah tambahan `.SpreadingRiskList(1).SplitRNMSharePct = 100` | **sama persis**, langkah 4, berketerangan *"JIKA SPREADINGLIST 1 DAN SplitRNMSharePct NULL \|\| 0"* | ✅ |
| presisi 20 lawan 10 | EDM `@divide(.SplitRNMSharePct, Primary.RNMShare, **20**)`; NB `100/@LengthOfPageList(...)` dengan `@divide(...,100,**10**)` | ✅ |

⭐ **Sebab EDM memerlukannya dan NB tidak, terbaca dari korpus:**
`TreatyNonPropSetSpreading` — yang mengisi `SplitRNMSharePct` — **ada di NB dan TIDAK ADA di EDM**.
Di polis baru daftar penyebaran dibangun dari nol dan medannya diisi di muka; di endorsemen daftar
itu datang dari polis yang diendorse dan medannya dapat tiba kosong — sementara rumus EDM
**membaginya**.

⛔ `[terbuka]` **Maksud dagangnya** — kenapa dasar `SplitRNMSharePct / RNMShare` dipilih untuk
endorsemen — **tidak terbaca dari korpus**, milik `[Product+Underwriting]`, **tidak dikarang**.

---

## 7 · Tiga temuan dari data produksi yang menyentuh EDM

Tiga contoh `DATA_JSON` sungguhan diterima work owner 22 September 2026.
⛔ **Tidak satu pun disimpan di berkas mana pun** — ketiganya memuat nama orang pada medan pemasar,
operator, dan tertanggung. Yang dicatat hanya **bentuk**, **nama medan**, dan **jumlah**.

### 7.1 ⛔⛔ Bentuk dokumen mengikuti jenis treaty

| | Contoh 1 | Contoh 2 | Contoh 3 |
| --- | ---: | ---: | ---: |
| jenis | proporsional | proporsional *(SOA)* | ⛔ **non-proporsional XOL** |
| medan skalar tingkat atas | 37 | 64 | 47 |
| ⭐ ada di **ketiganya** | | **23** | |
| ⭐ gabungan ketiganya | | **74** | |

**Penentunya:** `QuotationData.ProportionalType` bernilai `"Proportional"` lawan
`"NonProportional"`; `IsNewPolicyNonProp` `"0"` lawan `"1"`. `[terverifikasi]` Korpus menyebut
`ProportionalType` **170 kali**, `IsNewPolicyNonProp` **87 kali**.

⛔⛔ **Yang paling berbahaya: `ListInstallment` bersarang pada satu bentuk, datar pada bentuk lain.**
`[terverifikasi]` **Kedua bentuk nyata di dalam aturan**, disapu NB + EDM:

| Bentuk | Rujukan | Berkas |
| --- | ---: | ---: |
| ⛔ `ListInstallment(n).InstallmentList` — **bersarang** | **58** | 6 |
| ⛔ `ListInstallment(n).<medan skalar>` — **datar** | **101** | 11 |

⇒ ⛔ **Pengurai yang menganggap bentuknya seragam akan patah.**

### 7.2 ⭐⭐ `OldData` di dokumen polis baru — **dan korpus menjelaskannya**

Contoh 3 memuat **`OldData`** di dalam dokumen **polis baru**, berisi `TreatyXOLList`.

⭐⭐ **Sebabnya terbaca dari korpus, bukan misteri.** `[terverifikasi]`
`Activity\TreatyRealizationCheckXOLList` **di NB Treaty In** menetapkan:

```
pyWorkPage.PolicyTreatyIn.OldData.TreatyXOLList = pyWorkPage.PolicyTreatyIn.TreatyXOLList
```

⇒ **`OldData` ditulis saat realisasi polis baru**, jauh sebelum endorsemen mana pun.
`OldData.TreatyXOLList` dirujuk **130 kali di 9 berkas**.

### ⭐ Dampaknya ke P57 dan P58 — **diperiksa: keduanya tetap berdiri**

| Butir | Masih berdiri? | Sebab |
| --- | :---: | --- |
| **P57** *(endorsemen berlapis; pemilih layar `.OldData.EDMNo = ''` lawan `!= ''`)* | ✅ **ya** | Pemilihnya menguji **isi `EDMNo`**, bukan **keberadaan `OldData`**. ⭐ Temuan ini justru **menjelaskan kenapa** ujinya ditulis begitu — `OldData` ada bahkan pada polis baru, sehingga uji keberadaan akan selalu benar dan tidak berguna |
| **P58** *(data lama disalin dan beku)* | ✅ **ya** | `CreateEDMT` langkah 10 `Page-Copy` *"from TempPage.PolicyTreatyIn to MergePage.PolicyTreatyIn.OldData"*, lalu langkah 13 membuang halaman sumbernya. Penyalinan dan pembekuannya tidak berubah |

⛔⛔ **Tetapi satu butir `[terbuka]` BARU lahir dari keduanya digabung.** `[terverifikasi]`
Langkah 10 menyalin **seluruh halaman polis** ke dalam `OldData`. Karena halaman itu **sudah**
memuat `OldData` miliknya sendiri, hasilnya adalah **`OldData` di dalam `OldData`**.

⚠️ **Dan korpus tidak pernah membacanya:** rujukan berbentuk `OldData…OldData` = ⛔ **0**.

⇒ Susunan bersarang itu **dibuat tetapi tidak pernah dipakai** — menumpuk pada tiap endorsemen
berlapis. ⛔ **Perlu dipastikan apakah sengaja.** Pemiliknya `[work owner]`.
⛔ **Tidak ditutup, tidak diangkat menjadi pertanyaan bernomor di sini.**

### 7.3 ⛔⛔ Galat uang tersimpan — **berlaku juga untuk EDM**

```
premium angsuran  148157378.220000069   x 4  =  592629512.880000276
NetPremium        592629512.880000276        <- sama persis
selisih dari nilai bulat dua desimal         =  2,76 x 10^-7
```

`[terverifikasi]` Galatnya **lahir di rantai perhitungan, bukan di penyimpanan** — ia **57 kali
lebih besar** daripada galat `float64` untuk besaran ini dan **arahnya berlawanan**, dan
`592629512.88 ÷ 4` dalam desimal **tepat tanpa sisa**.

⭐⭐ **Butir ini berlaku juga untuk EDM, dan di sini akibatnya lebih tajam** — EDM **menghitung
selisih**. Selisih dua angka yang masing-masing membawa galat akan membawa **jumlah** galatnya, dan
endorsemen berlapis *(P57)* menumpuknya lagi.

⛔ **Tercatat sebagai `[terbuka]` pada P29 NB dengan tiga pilihan** *(ikuti apa adanya · bulatkan
saat migrasi · bulatkan saat dihitung)*. **Belum dipilih work owner.**
⛔ **Tidak dipilih di sini.** Pemutusnya `[work owner]`; presisi wajar per mata uang `[Finance]`.

⭐ **ADR-0003 tetap ditegakkan apa pun pilihannya** — sistem baru tidak memakai `float`, sehingga
**tidak menambah galat baru**. ⛔ Tidak ada ADR baru yang diperlukan.

---

## 8 · Cara menguji ulang berkas ini

Seluruh angka dapat diperiksa ulang dari korpus. ⚠️ **Pakai Python, bukan loop shell** — jalur
korpus memuat spasi. ⚠️ **Buang blok `pyExpressionGadget` lebih dulu** — ia peninggalan layar
penyunting.

```python
import os, re, glob, html
import xml.etree.ElementTree as ET
EDM = r"D:\XML\RNM_BRD\EDM Treaty In"
G = re.compile(r"<pyExpressionGadget>.*?</pyExpressionGadget>", re.S)
def load(p): return G.sub("", open(p, encoding="utf-8", errors="replace").read())
def txt(e, t):
    x = e.find(t)
    return html.unescape((x.text or "").strip()) if x is not None and x.text else ""

def top(r):                       # hanya pySteps PUNCAK — jebakan sarang
    par = {c: pa for pa in r.iter() for c in pa}
    out = []
    for e in r.iter("pySteps"):
        a = par.get(e); dalam = False
        while a is not None:
            if a.tag == "pySteps": dalam = True; break
            a = par.get(a)
        if not dalam: out.append(e)
    return out

def steps(ps, path=""):           # rekursi ke sub-langkah
    o = []
    for i, rd in enumerate(ps.findall("rowdata"), 1):
        no = f"{path}.{i}" if path else str(i)
        o.append((no, rd))
        sub = rd.find("pySteps")
        if sub is not None: o += steps(sub, no)
    return o

ps = pair = 0
for p in glob.glob(os.path.join(EDM, "Activity", "*.xml")):
    r = ET.fromstring(load(p))
    for t in top(r):
        for no, st in steps(t):
            m = (txt(st, "pyStepsActivityName") or txt(st, "pyStepsActivityNameUC")).upper()
            if m != "PROPERTY-SET": continue
            ps += 1
            pa = st.find("pyParamArray")
            if pa is None: continue
            pair += sum(1 for x in pa.findall("rowdata")
                        if txt(x, "PropertiesName") and txt(x, "PropertiesValue"))
print(ps, pair)          # harus 380 1309
```

**Cara kedua — pola teks murni, tanpa `ElementTree`:**

```python
ADJ = re.compile(r"<PropertiesName>(.*?)</PropertiesName>\s*"
                 r"<PropertiesValue>(.*?)</PropertiesValue>", re.S)
n = sum(len([1 for a, b in ADJ.findall(load(p)) if a.strip() and b.strip()])
        for p in glob.glob(os.path.join(EDM, "Activity", "*.xml")))
print(n)                 # harus 1309
```

### Invarian korpus

| Ukuran | Nilai |
| --- | --- |
| berkas `.xml` | **163** |
| byte | **17.639.259** |
| md5 gabungan | `24afd5726e083c9df6a692da94728139` |
| md5 folder `Activity` *(66 berkas)* | `22cf762f4ea4c513ef0e109dfedf94fb` |

> ⛔ *Koreksi 06-10:* invarian di atas berlaku **22-09**; md5 berubah — diukur ulang di **Bab 11**.

⚠️ **Bila md5 berubah, seluruh angka berkas ini wajib diukur ulang.** Korpus proyek ini **pernah
berubah di tengah jalan** — berkas `Protection_Act` pernah tertukar versi.

### ⛔ Enam jebakan yang sudah menjerat proyek ini

1. ⛔⛔ **Menebak nama tag.** `pyPropRef`, `pyPropertiesValue`, `pyStepsPage`, `pyCriteriaValue`,
   `pyResult`, `pyShapeName` — semuanya pernah dicari dan tidak ada.
   ⭐ **`Activity` memakai `PropertiesName`/`PropertiesValue` TANPA awalan `py`;
   `Flow` dan `DataTransform` memakai `pyPropertiesName`/`pyPropertiesValue` DENGAN awalan.**
   **Periksa daftar tag yang benar-benar ada lebih dulu.**
2. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>`.
3. `pyRowNum` pada `pyOrConditions` berbasis **nol**; baris tabel berbasis satu.
4. Memasangkan dua daftar menurut **urutan**, bukan menurut blok `rowdata` yang **sama**.
5. Membaca `pyConditionString` — keterangan manusia — alih-alih bentuk bertanda kurung yang
   dijalankan.
6. Menyatakan sesuatu nihil **tanpa menyebut lingkup penelusuran**.

⭐ **Jebakan ketujuh, ditambahkan hari ini:** `<pagedata>` membungkus seluruh berkas, sehingga pola
`<(\w+)>(.*?)</\1>` **menelan isi berkas** dan melaporkan nol untuk tag yang jelas berisi.

---

## 9 · Butir `[terbuka]` yang tercatat berkas ini

⛔⛔ **Berkas ini tidak menutup satu pun butir.** Tiga di bawah **dicatat**, bukan diangkat menjadi
pertanyaan bernomor — penomoran milik ronde pertanyaan.

| # | Butir | Pemilik |
| ---: | --- | --- |
| 1 | ⛔ **Dua stored procedure yang naskahnya belum pernah diminta** — `INSERTJSONPOLISMONITORING` dan `INSERTUPDATEACHIEVMENT` | `[data DBA]` |
| 2 | ⚠️ **`OldData` di dalam `OldData`** — dibuat `CreateEDMT` langkah 10, **tidak pernah dibaca korpus** *(0 rujukan)*, menumpuk pada endorsemen berlapis | `[work owner]` |
| 3 | ⚠️ **Skema tidak konsisten di SQL EDM** — 16 dari 31 `FROM` memakai `POOLDATA.`; `JSON_POLIS` muncul **dua ejaan**; `DATAPEGA` adalah **skema kedua** | `[data DBA]` |

### ⭐ Pemutakhiran 22 September 2026 — butir 1 **TERTUTUP**

Work owner menyerahkan naskah `INSERTJSONPOLISMONITORING` dan `InsertUpdateAchievment` pada hari
yang sama. Tercatat di jawaban **P1**. ⛔ **Butir 1 tidak lagi menahan `to-spec`.**

Butir 2 dan 3 **tetap terbuka**, dan dua butir **baru** lahir dari naskah itu:

| Butir baru | Pemilik |
| --- | --- |
| ⛔⛔ Penjaga anti-duplikat `ACHIEVEMENT` membandingkan **seluruh 17 kolom termasuk enam kolom uang** — dipatahkan galat `2,76 x 10^-7` yang terbukti ada di data produksi, sehingga baris yang seharusnya sama tersisip dua kali | `[work owner]` |
| ⚠️ Pesan galat `INSERTJSONPOLISMONITORING` memuat **markah HTML** dan membaca penanda lingkungan `p_ProductionLevel` — menyentuh **ADR-0005** | `[work owner]` |

⭐ Butir pertama **terikat langsung** pada pilihan (a)/(b)/(c) galat uang di **P29**: selama presisi
belum diputuskan, penjaga duplikat tidak dapat dirancang.

**Yang sudah terbuka di tempat lain dan berlaku juga di sini:** galat angka uang tersimpan *(P29)* ·
maksud dagang dasar `SplitRNMSharePct / RNMShare` *(P60)* · selisih mana yang dipakai laporan
*(P57)* · batas berapa kali satu polis boleh diendorse *(P57)* · tiga pemeriksaan pencegah yang
dimatikan *(P58)*.

---

## 10 · TELEMETRI EKSEKUSI

### Cara tiap angka diperoleh

| Angka | Cara A | Cara B | Sepakat? |
| --- | --- | --- | --- |
| berkas & byte | `glob **/*.xml` | `os.walk` + md5 gabungan | ✅ 163 |
| langkah `Property-Set` | `ElementTree`, `pySteps` puncak, rekursi | pola teks `<pyStepsActivityName>` | ✅ 380 |
| pasangan nama=nilai | pasangan dalam satu `rowdata` | pola teks bersebelahan | ✅ 1.309 |
| penjangkauan `Activity` | urai `Call` dari titik masuk | nama disebut >1 kali | ✅ 66 / 0 |
| beda dengan NB | tanda tangan **perilaku**, nol metadata | seluruh tag keterangan | ⚠️ 4 lawan 3 — **dijelaskan** |
| galat uang | `Decimal` presisi 40 | pembanding `float64` | ✅ |

### ⛔ Angka yang TIDAK dapat direproduksi

| Angka brief | Terukur | Rekonstruksi dicoba | Hasil |
| --- | ---: | ---: | --- |
| langkah `Property-Set` **378** | **380** | 3 | ⛔ **tidak satu pun menghasilkan 378** |
| pasangan **1.062** | **1.309** | 6 | ⛔ **tidak satu pun menghasilkan 1.062** |
| *"enam berkas beda catatan, empat di kelompok 69"* | **3 berkas, nol di kelompok 69** | 2 | ⛔ **tidak tereproduksi** |

⭐ **Ketiganya ditarik dari peredaran.** ⚠️ **Sebab pastinya tidak dapat saya nyatakan** — sesi yang
menghasilkannya tidak menuliskan perintahnya. ⛔ **Tidak dikarang.**

### Batas ronde ini

| ⛔ Larangan | Dipatuhi? | Bukti |
| --- | :---: | --- |
| `grilling-ronde-1.md` tidak disentuh | ✅ | nol tulis; koreksinya hidup di Bab 6 |
| korpus `D:\XML\RNM_BRD\` read-only | ✅ | nol berkas korpus disunting |
| jangan menutup butir `[terbuka]` | ✅ | **0** ditutup; **3** dicatat di Bab 9 |
| jangan membuat ADR baru | ✅ | **0** — dan tidak ada yang perlu |
| nol kode, nol DDL | ✅ | kode di Bab 8 adalah **perintah uji**, bukan kode sistem |
| nol nama orang tersalin | ✅ | ⛔ **0**. Nilai `Flow` berbentuk *"… IS IN ‹nama›'S INBOX"* **tidak disalin** — hanya nama medan dan jumlahnya |
| nol contoh JSON tersimpan | ✅ | hanya **bentuk**, nama medan, jumlah, dan **lima angka medan uang** |

### Ongkos

| Ukuran | Nilai |
| --- | ---: |
| panggilan model | **40** |
| token keluar | **128.920** |
| token cache ditulis | **214.219** |
| token cache dibaca | **20.390.775** |
| berkas korpus dibaca | **163** |
| berkas keluaran | ****1** — `KEADAAN-EDM-TREATY-IN.md`; ditambah `PERTANYAAN-RONDE-1.md` disunting** |
| berkas korpus disunting | **0** |
| butir `[terbuka]` ditutup | **0** |

⚠️ **Telemetri TIDAK diukur dari luar.** `claude --print --output-format json` akan menjalankan sesi
**baru dan terpisah**, yang ongkosnya bukan ongkos ronde ini. ⭐ Angka di atas **hasil pengurangan
terhadap baseline** yang diambil dari catatan sesi pada awal ronde. ⚠️ Angka token sejati tidak
terlihat dari dalam sesi.

---

*Disusun 22 September 2026. Ronde berikutnya: `to-spec` EDM Treaty In, yang membaca berkas ini
sebagai sumber keadaan terukur.*

---

## 11 · Koreksi 06-10-2026

⛔ `grilling-ronde-1.md` **tetap tersegel**; koreksinya ditulis di sini (11.3). Log seluruh koreksi dokumen:
`KOREKSI-DOKUMEN-2026-10-06.md`; peta per AC: `REKONSILIASI-AC.md`. ⛔ Nol butir `[terbuka]` milik WO/DBA ditutup.

### 11.1 Ukur ulang korpus — md5 **berubah**, sensus **tetap**

| Ukuran | 22-09 (Bab 2, 8) | **06-10** | Cara |
| --- | ---: | ---: | --- |
| berkas `.xml` | 163 | **163** | `glob **/*.xml` = `os.walk` (164 berkas, 1 `.xlsx`) |
| tipe rule | 12 | **12** | nama folder induk |
| byte `.xml` | 17.639.259 | ⛔ **17.640.859** (+1.600) | jumlah `getsize` |
| md5 gabungan | `24afd572…94728139` | ⛔ **tidak tereproduksi** | lihat di bawah |
| byte `Activity` (66) | — | 7.383.092 | jumlah `getsize` |
| langkah bermetode | 618 | **618** | ElementTree, `pySteps` puncak + rekursi |
| langkah `Property-Set` | 380 | **380 · 380** | **A** ElementTree (resep Bab 8) · **B** pola teks `<pyStepsActivityName>Property-Set</…>` |
| pasangan nama=nilai | 1.309 | **1.309 · 1.309** | **A** satu `rowdata` · **B** pola teks bersebelahan (resep Bab 8) |
| `Property-Set` berisi / kosong | 379 / 1 | **379 / 1** (`FillPaymentInstallment` 4.3.2) | A |
| `PropertiesName` / `PropertiesValue` terisi | 1.310 / 1.311 | **1.310 / 1.311** | pola teks |
| kedalaman 1…6 | 165·83·65·29·20·18 | **165·83·65·29·20·18** | A |
| metode teratas | Bab 2.2 | **sama** (`PROPERTY-SET` 380 · `RDB-LIST` 60 · `PAGE-NEW` 21 …) | A |

**Sebab md5 berubah** — dua berkas bercap waktu sesudah ukuran 22-09 (161 berkas lain 08-09-2026):
`Activity/CountSpreading_Act.xml` (22-09-2026 17:50, 91.357 B) dan `Activity/InsetTreatyInProdAddendum_Act.xml`
(23-09-2026 14:31, 309.257 B; `pxUpdateDateTime` ekspor 23-09). Perintah: `find "EDM Treaty In" -type f -newermt 2026-09-21`.

⚠️ **Resep md5 22-09 tidak tertulis** di Bab 8. Lima varian diukur 06-10 — tidak satu pun mungkin cocok karena
isi berubah: concat byte urut jalur penuh `57926cc1b3f4524cb0cf73e1995ca211` · urut relpath huruf kecil
`f76847aeb328e328b5e0b5231dab7bcf` · urut nama berkas `5cd063b3286beb23864f078794f0ff0c` · md5-dari-md5
`f57c2a3284437023325abbd07b45a9b6` · relpath+byte `1375a54a19f111a8cc37524b5c5179c1`; folder `Activity` concat
`2d2f35d333c326505de75e8603717436`. ⭐ **Invarian baru yang diusulkan** (resep tertulis): md5 concat byte seluruh
`.xml` urut jalur penuh = **`57926cc1b3f4524cb0cf73e1995ca211`**, 163 berkas, 17.640.859 B.

```python
import os, glob, hashlib
EDM = r"D:\XML\RNM_BRD\EDM Treaty In"
h = hashlib.md5()
for p in sorted(glob.glob(os.path.join(EDM, "**", "*.xml"), recursive=True)):
    h.update(open(p, "rb").read())
print(h.hexdigest())          # 06-10-2026: 57926cc1b3f4524cb0cf73e1995ca211
```

⭐ Karena 380 / 1.309 sepakat dua cara sesudah perubahan, **angka Bab 2.2 tetap berlaku**. Angka Bab 2.4
(penjangkauan) dan 6.3 (73 senama) **tidak diukur ulang** ronde ini.

### 11.2 Koreksi atas berkas ini

| Bab | Bunyi lama (dikutip) | Bunyi baru | Bukti |
| --- | --- | --- | --- |
| 3 | *"22 sambungan"* | **24** baris `pyConnectors` (23 `TransitionN` + 1 sisa penyuntingan) | `Flow/InputAddendumTreatyIn.xml`; sejalan `grilling-ronde-1.md` baris 334 (24) |
| 3 | *"Sejalan dengan keputusan tangga tiga jenjang pada P13 — tidak ada penyimpangan"* | XML punya cabang **Sec Head menyelesaikan sendiri**: `Decision9` → `Utility1` (`Transition24`) bila `LetterNo` ≠ `"TREATYINDEPTHEAD"`; `LetterNo` diisi hanya bila \|`TotalPremium`\| × kurs > 200.000.000. Tiga jenjang = **keputusan** (brief F0-F1 06-10, prompt WO; NB K2) ⇒ **penyimpangan sadar** | `When/ToTREATYDEPTHEAD.xml`; `Activity/CekLimitTreatyAcc_Act.xml` langkah 2.4 (2.2 ber-`//`); NB `docs/INVENTARIS-XML.md` baris 70 |
| 5.3 | P50 *"pengiriman kedua tidak pernah berjalan"* | premis tidak terbukti: `IsFacRetro` diset `InputPolicyTreatyEDMDetail_NP`/`_AdjPremi` langkah 3 bila master `FacultativeShare > 0` ⇒ `serviceInsertArasapas_act` langkah 7 dapat jalan. Keputusan P50 tidak diubah — **butir WO** | berkas-berkas itu; `When/IsFacRetro.xml` |
| 1, 6.1 | `KEADAAN-NB-TREATY-IN.md`, `..\nb-treaty-in\VERIFIKASI-P18.md` | `modul/nbtreatyin/docs/…` | folder NB |
| 2.1, 8 | byte 17.639.259, md5 `24afd572…` | 11.1 di atas | 11.1 |

### 11.3 Koreksi atas `grilling-ronde-1.md` (tersegel — dikutip, tidak disunting)

| Baris | Bunyi ronde 1 (dikutip) | Yang benar | Bukti |
| ---: | --- | --- | --- |
| 342 | *"Keputusan work owner untuk NB (tangga tiga jenjang) sudah menjadi kenyataan di EDM."* | ⛔ **Tidak di XML**: Sec Head dapat menyelesaikan sendiri bila premi ≤ 200 juta (11.2). Tiga jenjang menjadi kenyataan **di sistem baru** sebagai keputusan, bukan sebagai paritas | 11.2 baris kedua |
| 312 | *"Nol rule yang menguji status FACOUT terbaca."* | sudah dikoreksi Bab 6.2 (`IsSuccessHitService` menguji keduanya); **tambahan 06-10**: gerbang pengiriman FACOUT (`IsFacRetro`) **diset** di jalur EDM | 11.2 baris ketiga |
| 334 | *"Sambungan … 24"* | ✅ **benar** — justru Bab 3 berkas ini (22) yang keliru | `Flow/InputAddendumTreatyIn.xml` |

### 11.4 Butir yang lahir dari koreksi ini — dicatat, **tidak ditutup**

| # | Butir | Pemilik |
| ---: | --- | --- |
| 1 | Cakupan *"satu rumus untuk semua"* (ID-28) atas NonProp — XML memakai batas bawah 0, prorata, pajak ulang, mentah untuk jenis 4/2; AC 22 menuntut angka NonProp tak berubah | `[work owner]` |
| 2 | P50 dengan `IsFacRetro` = 1 (master `FacultativeShare > 0`) | `[work owner]` |
| 3 | Pembatalan: tiru 16 medan `SetEDMTCancel` atau nolkan seluruh uang (AC 14) | `[work owner]` |
| 4 | Sumber rincian angsuran NonProp: salinan master (`SetInstallmentValue`) atau bangun ulang dari selisih XOL (`FillPaymentInstallmentEDMT`) — AC 38 | `[work owner]` · asisten utama |
| 5 | Daftar kolom khas EDM (AC 32) disusun ulang dari katalog NB | asisten utama |
| 6 | `T_POLIS_SURVEY` dan `HISTORYAKSEPTASIPRODUCTION` pada generasi EDM — cacah AC 55 | `[work owner]` |
| 7 | Penghapusan produksi bersyarat `IsPEGAPROD`; langkah 10 berprasyarat `IsTreatyIn` (AC 38, AC 42) | `[work owner]` |
| 8 | P29 (a) untuk EDM — penunjuk ke NB ID-15 / F20 (spec.md 8.2) | `[work owner]` · `[Finance]` |

### 11.6 Jawaban work owner 07-10-2026 atas butir 11.4

Work owner 07-10-2026 (kedua): *"ikuti rekomendasi semua"* - butir 2, 4, 7, 8 ditutup seperti tabel di bawah. Bila
konversi Arasapas kelak disambung, butir 2 (P50 `IsFacRetro`: kirim FacOut bila master `FacultativeShare > 0`) dan 7
(`IsPEGAPROD` / `IsTreatyIn`) dibangun **mengikuti XML**.

| # 11.4 | Keadaan sesudah 07-10-2026 | Dasar |
| ---: | --- | --- |
| 1 | **ditutup** - Prop satu rumus, NonProp `CalculateDifferenceEDM_act` apa adanya (migrasi 363 mengikuti XML) | keputusan WO ID-28/30 06-10-2026 |
| 2 | **gugur** - konversi Arasapas tetap tidak disambung, langkah 7 FacOut tidak pernah jalan | WO 07-10-2026 butir 4 |
| 3 | **ditutup** - tiru 16 medan `SetEDMTCancel` | WO 07-10-2026 butir 6 *"ikuti XML dulu"* |
| 4 | **ditutup** - mengikuti XML (`FillPaymentInstallmentEDMT` membangun ulang) | XML; dijelaskan ke WO 07-10-2026 butir 13 |
| 5 | **terbuka** - daftar kolom khas AC 32 belum disusun ulang | asisten utama |
| 6 | **ditutup** - generasi EDM menulis `T_POLIS_SURVEY` dan `HISTORYAKSEPTASIPRODUCTION` (pola NB); cacah AC 55 = 🟡 | WO 07-10-2026 butir 9 *"YA"* |
| 7 | **gugur** - penghapusan produksi dan pesan gagal konversi tidak dibangun (konversi tidak disambung) | WO 07-10-2026 butir 4 |
| 8 | **ditutup** - nilai lama disalin apa adanya (pemuat SUMBER 'PEGA', AC 39), sama dengan NB ID-15 | NB ID-15; dijelaskan ke WO 07-10-2026 butir 13 |

### 11.5 TELEMETRI EKSEKUSI ronde koreksi

| Ukuran | Keadaan |
| --- | --- |
| berkas korpus dibaca | 163 `.xml` (sensus) + 17 activity / 1 flow / 1 When / 1 DT dibaca rinci per langkah |
| berkas korpus disunting | **0** |
| `grilling-ronde-1.md` disunting | **0** |
| token · durasi · biaya | ⛔ **tidak diukur** — angka token sejati tidak terlihat dari dalam sesi; tidak ditaksir |
