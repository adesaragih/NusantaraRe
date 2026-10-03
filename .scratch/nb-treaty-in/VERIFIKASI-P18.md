# VERIFIKASI-P18 — apakah P18 berdiri di atas premis yang keliru?

> **Untuk:** pemilik pekerjaan migrasi Treaty Inward, dan siapa pun yang akan memutuskan apakah
> surat P18 jadi dikirim.
>
> **Lingkup:** `NB Treaty In` **dan** `EDM Treaty In`. Korpus dibaca saja; tidak satu berkas pun
> di dalamnya disentuh.
>
> ⛔ **Ronde ini tidak menutup P18.** Penutupan milik `[pemilik export Pega]` dan work owner.
> Yang diberikan di sini adalah **vonis teknis** atas premisnya, berikut daftar apa yang harus
> berubah bila vonis itu diterima.

---

## Bab 1 — Vonis dalam satu halaman

⭐⭐ **Premis P18 KELIRU. Isi 268 langkah penetapan nilai ADA di dalam ekspor yang sudah kami
terima — lengkap, dan tidak terpotong.**

P18 berbunyi, apa adanya, di `PERTANYAAN-untuk-pemilik-export-Pega.md`:

> *"Kami dapat melihat **bahwa** langkah itu ada, dan **urutannya**, ⛔ **tetapi tidak satu pun
> rumus atau nilai yang ditetapkannya ikut terkirim.** Yang ada hanya salinan sementara yang
> tertinggal dari layar penyunting, dan salinan itu **terpotong di tengah** sehingga tidak dapat
> dipakai."*

**RALAT.** Kalimat itu salah pada kedua bagiannya sekaligus:

| Yang didalilkan P18 | Yang terukur | Penanda |
| --- | --- | --- |
| "tidak satu pun rumus ikut terkirim" | ⭐ **2.481 pasangan nama=nilai terisi** di `NB Treaty In`, **1.309** di `EDM Treaty In` | `[terverifikasi]` |
| "hanya salinan sementara dari layar penyunting" | ⭐ salinan itu **nyata ada** (`pyExpressionGadget`) — tetapi ia **bukan satu-satunya**; isi sebenarnya ada di tempat lain dan utuh | `[terverifikasi]` |
| "terpotong di tengah" | ⭐ **nol** tanda terpotong pada tujuh uji keutuhan | `[terverifikasi]` |
| "94 dari 100 langkah tidak dapat ditiru" | ⭐ **938 dari 939** langkah `Property-Set` membawa isinya sendiri — **99,9 %** | `[terverifikasi]` |
| "layar Kepala Departemen tidak akan dapat disimpan sama sekali" | ⭐ 28 dari 34 medan wajib+terkunci punya langkah pengisi; 3 medan sisanya **dapat diketik** di layar lain | `[terverifikasi]` |

### ⛔ Sebab kekeliruannya — dua huruf

Isi langkah `Activity` tersimpan di tag **`PropertiesName`** dan **`PropertiesValue`**, **tanpa
awalan `py`**, di dalam `pyParamArray` milik langkah itu.

Ronde 2, 3, dan 4 memeriksa tag **`pyPropRef`**, **`pyPropertiesName`**, **`pyPropertiesValue`** —
**dengan** awalan `py`. Ketiganya memang kosong. Kesimpulan "isinya tidak dikirim" ditarik dari
kekosongan itu.

⭐⭐ **Kedua nama tag itu hidup berdampingan di berkas yang sama.** Bukti paling telanjang ada di
`SumTSIPremiSpreadedRNM_Act`, tempat `<pyPropertiesName/>` yang **kosong** terselip persis di
antara pasangan tanpa awalan yang **terisi**:

```
<PropertiesName>Local.TsiTopRisk</PropertiesName>
<pyPropertiesName/>                        <- yang saya periksa di ronde 2-4. Selalu kosong.
<PropertiesValue>.CARI3</PropertiesValue>  <- yang sebenarnya membawa isi.
```

| Tag di folder `Activity` | NB: tag | NB: **berisi** | EDM: tag | EDM: **berisi** |
| --- | ---: | ---: | ---: | ---: |
| `PropertiesName` *(tanpa awalan)* | 2.507 | ⭐ **2.485** | 1.329 | ⭐ **1.310** |
| `PropertiesValue` *(tanpa awalan)* | 2.507 | ⭐ **2.482** | 1.329 | ⭐ **1.311** |
| `pyPropertiesName` | 1 | ⛔ **0** | 0 | ⛔ **0** |
| `pyPropertiesValue` | 0 | ⛔ **0** | 0 | ⛔ **0** |
| `pyPropRef` | 616 | ⛔ **0** | 437 | ⛔ **0** |

⚠️ **Jebakan baru untuk `CLAUDE.md` §4a — jebakan #7: _nama tag berbeda menurut JENIS aturan._**
Di `Activity` yang berisi adalah bentuk **tanpa** awalan. Di `DataTransform` justru **kebalikannya**
— yang berisi adalah `pyPropertiesName` / `pyPropertiesValue` **dengan** awalan, dan bentuk tanpa
awalan nol. Memeriksa satu jenis aturan lalu menyimpulkan untuk jenis lain melahirkan nol palsu.

---

## Bab 2 — Sensus: dua cara yang benar-benar berbeda

**Jendela:** folder `Activity` pada kedua modul. Blok `pyExpressionGadget` **dibuang lebih dulu di
kedua cara** — ia peninggalan layar penyunting, dan itu aturan baca yang sudah mengikat sejak
ronde 3.

- **CARA A — `ElementTree`, sarang dihormati.** Pasangan dihitung hanya bila `PropertiesName` dan
  `PropertiesValue` berada dalam **satu blok `rowdata` yang sama** di bawah `pyParamArray` milik
  langkah. Bukan menurut urutan kemunculan.
- **CARA B — pola teks, tanpa `ElementTree` sama sekali.** Mencacah tiap tag secara terpisah, lalu
  pasangan yang bersebelahan di teks.

### 2.1 Hasil

| | NB Treaty In | EDM Treaty In |
| --- | ---: | ---: |
| berkas `Activity` | **92** | **66** |
| langkah seluruh metode | **1.342** | **618** |
| langkah `Property-Set` | **939** | **380** |
| baris `pyParamArray` | 2.895 | 1.556 |
| ⭐ **CARA A — pasangan nama+nilai terisi** | ⭐ **2.481** | ⭐ **1.309** |
| CARA B — `PropertiesName` terisi | 2.485 | 1.310 |
| CARA B — `PropertiesValue` terisi | 2.482 | 1.311 |
| CARA B — bersebelahan & dua-duanya terisi | 2.480 | 1.309 |

### 2.2 ⛔ Dua cara sempat berselisih. Selisihnya dikejar sampai nol.

⚠️ **Percobaan pertama saya gagal, dan kegagalannya dicatat utuh** — bukan dihapus.

Percobaan pertama membaca **1.375** langkah `Property-Set` di NB *(bukan 939)* dan **3.458**
pasangan *(bukan 2.481)*. Sebabnya: saya menelusuri **setiap** elemen `<pySteps>` dengan `iter()`.
Ternyata **1.736 dari 1.828** elemen `<pySteps>` di NB adalah **sarang** — setiap langkah membawa
wadah sub-langkahnya sendiri — sehingga isi yang sama ditelusuri berkali-kali. Yang benar: hanya
**92 `pySteps` puncak untuk 92 berkas**. Sesudah diperbaiki, angka langkah jatuh ke 939.

Sesudah perbaikan itu, sisa selisih NB adalah **2.481 − 2.480 = 1** dan **2.485 − 2.481 = 4**.
Semuanya dijelaskan tuntas, bukan ditoleransi:

| Selisih | Jumlah | Sebab | Penanda |
| --- | ---: | --- | --- |
| pasangan terisi tapi **tidak bersebelahan** di teks | **1** | `<pyPropertiesName/>` kosong menyisip di antaranya *(`SumTSIPremiSpreadedRNM_Act`)* | `[terverifikasi]` |
| nama terisi, **nilai kosong** | **4** | baris `pyParamArray` setengah terisi | `[terverifikasi]` |
| nilai terisi, **nama kosong** | **1** | idem | `[terverifikasi]` |
| pasangan di langkah yang **nama metodenya kosong** di ekspor | **11** | 4 aturan; isinya tetap terbaca penuh | `[terverifikasi]` |

⭐ **Rekonsiliasi tepat:** 2.470 *(di bawah langkah bermetode)* + 11 *(di bawah langkah tanpa nama
metode)* = **2.481**. Untuk EDM kedua cara langsung sepakat di **1.309** — selisih **nol**.

### 2.3 Angka yang diklaim brief vs angka yang terukur

⚠️ Brief menyuruh **tidak mempercayai angkanya dan mengukur sendiri**. Sudah. Hasilnya berbeda,
dan **ke arah yang lebih menguntungkan**:

| Pernyataan | klaim NB | **ukur NB** | klaim EDM | **ukur EDM** |
| --- | ---: | ---: | ---: | ---: |
| `Activity` | 92 | ✅ **92** | 66 | ✅ **66** |
| langkah `Property-Set` | 938 | **939** *(938 berisi)* | 378 | **380** *(379 berisi)* |
| pasangan terisi | 2.080 | ⭐ **2.481** | 1.062 | ⭐ **1.309** |
| aturan punya langkah tapi **nol** pasangan | 0 | ✅ **0** | — | ✅ **0** |
| nilai tampak terpotong | 0 | ✅ **0** | — | ✅ **0** |

⚠️ Angka mapan **938** ternyata adalah jumlah langkah `Property-Set` yang **berisi**, bukan jumlah
seluruhnya. Seluruhnya **939**; satu di antaranya kosong *(Bab 4)*.

---

## Bab 3 — Uji keutuhan: adakah tanda nilai terpotong?

Tujuh uji, dijalankan atas **2.481** nilai NB dan **1.309** nilai EDM.

| Uji | NB | EDM | Bacaan |
| --- | ---: | ---: | --- |
| kurung `( )` tidak seimbang | **0** | **0** | ✅ |
| kurung `[ ]` tidak seimbang | **0** | **0** | ✅ |
| nilai berakhir dengan operator `+ - * / ,` | **0** | **0** | ✅ |
| petik ganda ganjil *(tak tertutup)* | **0** | **0** | ✅ |
| memuat elipsis, atau berakhir tiga titik | **0** | **0** | ✅ |
| tepat di batas 255 / 256 / 512 / 1024 aksara | **0** | **0** | ✅ |
| petik tunggal ganjil | 56 | 4 | ⚠️ **bukan pemotongan** |

⚠️ **Lima puluh enam itu alarm palsu, dan saya periksa satu per satu.** Semuanya apostrof Inggris
di dalam literal teks — `"Reff Number Can't Empty "`, `"... isn't match with premium RNM! "`,
`"Error: Nopolis is incorrect / doesn't exists"`. Tidak satu pun tanda terpotong.

### Sebaran panjang nilai

| | NB | EDM |
| --- | ---: | ---: |
| terpendek | 1 | 1 |
| median | **16** | **16** |
| rerata | 25,1 | 34,3 |
| persentil 90 / 95 / 99 | 52 / 82 / 143 | 89 / 115 / 287 |
| **terpanjang** | **393** | **367** |
| nilai ≥ 255 aksara | 3 | 18 |
| nilai ≥ 1024 aksara | **0** | **0** |

⭐ **Sebaran ini sendiri adalah buktinya.** Ekspor yang memotong akan menumpuk nilai persis di satu
angka — 255, 256, atau 1024. Di sini tidak ada tumpukan seperti itu sama sekali; puncaknya di 393
dan 367, angka yang tidak istimewa bagi mesin mana pun.

Nilai **terpanjang** di NB — 393 aksara, empat lapis `@if`, kurung tertutup rapi sampai huruf
terakhir:

```
@if(Local.DateToUW<=Local.Begindate,
    @if(.OfferFacIn.QuotationData.IsTBA=="true",
        @DateTime.addCalendar(@FormatDateTime(Local.LastAksep,"yyyyMMdd","Asia/Jakarta","in_ID"),0,3,0,0,0,0,0),
        @DateTime.addCalendar(@FormatDateTime(Local.Begindate,"yyyyMMdd","Asia/Jakarta","in_ID"),0,1,0,0,0,0,0)),
    @DateTime.addCalendar(@FormatDateTime(Local.LastAksep,"yyyyMMdd","Asia/Jakarta","in_ID"),0,1,0,0,0,0,0))
```

### Bentuk nilai — apakah ini benar-benar rumus?

| Bentuk | NB | EDM |
| --- | ---: | ---: |
| rujukan properti *(`.Currency`, `Local.X`, `Param.Y`)* | 1.080 | 612 |
| ⭐ ungkapan / fungsi *(`@if`, `@divide`, tambah-kurang-kali-bagi)* | ⭐ **557** | ⭐ **272** |
| angka harfiah | 488 | 134 |
| literal teks | 237 | 121 |
| lain | 119 | 170 |

⭐ **557 ungkapan berhitung di NB dan 272 di EDM.** Inilah barang yang P18 nyatakan hilang.

---

## Bab 4 — Cukupan per langkah: apakah tiap langkah membawa isinya sendiri?

⚠️ Pemeriksaan ini dijalankan **di dalam blok `rowdata` milik tiap langkah**, lewat
`pyParamArray`-nya sendiri — **bukan** menurut urutan kemunculan di berkas, dan **bukan** dengan
memasangkan dua daftar sejajar. *(Memasangkan menurut urutan adalah jebakan §4a #6, dan ia sudah
sekali melahirkan hasil omong kosong di ronde 3.)*

| | NB Treaty In | EDM Treaty In |
| --- | ---: | ---: |
| langkah `Property-Set` | 939 | 380 |
| ⭐ **punya ≥1 pasangan terisi di bloknya sendiri** | ⭐ **938** *(99,9 %)* | ⭐ **379** *(99,7 %)* |
| ⛔ nol pasangan terisi | **1** | **1** |
| aturan yang punya langkah tapi nol pasangan | **0** | **0** |

### Satu langkah kosong itu — apa isinya?

`FillPaymentInstallment`, langkah **4.3.2**. Satu baris `pyParamArray`, nama kosong, nilai kosong,
**tanpa keterangan, tanpa gerbang, tanpa blok**. ⭐ Ini langkah yang memang kosong di sistem lama —
tempat yang disiapkan lalu tidak diisi — bukan isi yang hilang saat ekspor. `[terverifikasi]`

### ⭐ Lingkup yang benar-benar diminta P18: 51 aturan / 268 langkah

| Ukuran | Nilai |
| --- | ---: |
| aturan pada lampiran, ditemukan seluruhnya | **51 / 51** |
| langkah `Property-Set` terukur | **269** *(lampiran menulis 268)* |
| baris `pyParamArray` | 708 |
| ⭐ **pasangan nama=nilai terisi** | ⭐ **707** |
| aturan yang jumlah langkahnya beda dari lampiran | **1** — `SetValidateInstallment_Act` *(lampiran 2, terukur 3)* |

⭐ **707 pasangan untuk 269 langkah. Nol yang hilang.**

### Akibat kedua yang didalilkan P18 — sembilan medan terkunci

P18 menambahkan: *"sembilan medan pada layar Kepala Departemen wajib diisi tetapi terkunci ...
⛔ layar Kepala Departemen tidak akan dapat disimpan sama sekali."*

Terukur: **34** medan wajib **dan** terkunci di seluruh `Section` NB.

| | Jumlah | Keterangan |
| --- | ---: | --- |
| ⭐ ada langkah yang mengisinya | **28** | mis. `.PremiOnp` ← `CountOGPONP_Act`; `.RiCommOgp` ← `CalculatePremi_Act`; `.StatementDate` ← `InputPolicyTreatyInPre_Act` |
| ⚠️ tidak ada langkah pengisinya | **6** *(3 medan × 2 layar)* | `.ExcessLoss`, `.OutstandingClaim`, `.SalvageValue` |

⭐ **Ketiga medan sisa itu pun bukan lubang.** Pada `DetailPolicyTreatyIn` dan
`GeneralPolicyTreatyIn` ketiganya **`bacasaja=false`** — **diketik langsung oleh underwriter** pada
tahap lebih awal, lalu hanya **ditampilkan terkunci** di layar Kepala Departemen. Di
`EDM Treaty In` ketiganya bahkan ditulis oleh langkah *(`CopyGeneralDataEDM_act`, `SetEDMTCancel`,
`EDMTCalculateTreatyDifference`)*. `[terverifikasi]`

⛔ **Jadi dalil "layar Dept Head tidak dapat disimpan sama sekali" juga runtuh.**

---

## Bab 5 — Lima aturan teratas, seluruh isinya, apa adanya

Diambil dari **puncak** `LAMPIRAN-P18-ATURAN-DIMINTA.md` — lima aturan terberat, **37 %** dari seluruh
langkah yang diminta — dan **306** pasangan terisi di antara mereka saja.

| # | Aturan | Langkah *(lampiran)* | ⭐ Pasangan terisi terukur |
| ---: | --- | ---: | ---: |
| 1 | `InsertToTreatyOutXOLList` | 30 | **107** |
| 2 | `InsertToTreatyXOLListRetroShare` | 23 | **72** |
| 3 | `GeneratePolicyNoTreaty_Act` | 17 | **27** |
| 4 | `InsertToTreatyXOLList` | 16 | **70** |
| 5 | `FillPaymentInstallment` | 15 | **30** |

### 5.1 Contoh — inti perhitungan pajak, terbaca utuh

Dari `InsertToTreatyXOLList` langkah **3.2.4**, *"Insert Every Value to the property"*:

```
local.currency               = InputXOL.CARI31
local.grosspremi             = local.grosspremi + @toDecimal(InputXOL.CARI32)
local.netpremi               = local.netpremi   + @toDecimal(InputXOL.CARI13)
local.duetovalue             = local.duetovalue + @toDecimal(InputXOL.CARI47)
local.deduction              = local.deduction  + @toDecimal(InputXOL.CARI44)
Local.BrokerageFeeSebenarnya = @if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",
                                   @divide(local.deduction,@divide(102.2,100,8),8),
                                   local.deduction)
Local.PPHValue               = Local.BrokerageFeeSebenarnya * @divide(2,  100,8)
Local.PPNValue               = Local.BrokerageFeeSebenarnya * @divide(2.2,100,8)
```

⭐ **Di sini terbaca sekaligus:** pembagi **102,2** untuk `TypeTax = "Inclusive"`, PPH **2 %**,
PPN **2,2 %**, dan **presisi 8 angka** pada tiap `@divide`. Itu tepat angka-angka yang selama ini
dikarantina di `spec.md` §8.2 sebagai "belum diketahui".

Dan langkah **3.6** yang menyusulnya, di balik gerbang `FlagPPH=="true"`:

```
...TreatyXOLList(<CURRENT>).NetPremiAfterPPH = local.netpremi + Local.PPHValue
...TreatyXOLList(<CURRENT>).NetPremiAfterPPN = local.netpremi + Local.PPNValue
...TreatyXOLList(<CURRENT>).NetPremiAfterTax = local.netpremi + Local.PPHValue + Local.PPNValue
```

### 5.2 Contoh — pembentukan nomor polis, terbaca sampai potongan terakhir

Dari `GeneratePolicyNoTreaty_Act`, langkah 7-9, 18-19, 26, 28:

```
InputData.CARI20  = "QR"            <- when premium
InputData.CARI20  = "QP"            <- when pay
InputData.CARI20  = "TP"            <- when ClaimType == "XOL Retro"
InputData.CARI23  = @addCalendar(TempDate.pxResults(1).CARI1,'0','1','0','0','0','0','0')
InputData.CARI24  = @substring(InputData.CARI23,4,6) + "." + @substring(InputData.CARI23,0,4)
ParamSeq.CARI4    = ParamSeq.HASIL3 + InputData.CARI20
...PolicyNo       = ParamSeq.CARI4 + ".T" + ...OJKBusinessID + "." + ParamSeq.HASIL1 + "." + ParamSeq.HASIL2
```

⭐ Termasuk **aturan tanggal 25**: bila sudah lewat tanggal 25 berjalan, bulan digeser satu
*(`@addCalendar(...,'0','1',...)`)*. Dan langkah **5.1** yang menetapkan `ProductionDate` ke waktu
kini terbaca **`pyStepsPreCondition = "false"` — gerbangnya dimatikan**, sehingga ia tidak pernah
jalan.

### 5.3 ⭐ Cukupkah ini untuk menulis ulang perhitungannya?

**Cukup.** Lima aturan itu memberi, untuk tiap langkah: **sisi kiri** *(properti yang ditulis,
lengkap dengan indeks `<CURRENT>` / `<APPEND>` / `<LAST>`)*, **sisi kanan** *(ungkapan penuh)*,
**urutan**, **kedalaman sarang** *(nomor seperti `4.7.2.3.2.1`)*, **halaman yang di-loop**,
**gerbang beserta arahnya**, dan **keterangan penulis aslinya**.

⚠️ **Tiga hal yang tetap perlu dari luar** — dan **tidak satu pun milik P18**:

| Yang perlu | Dari siapa | Sudah ada? |
| --- | --- | --- |
| arti medan `CARIn` / `HASILn` pada RDB-List | naskah SQL di folder `RDBList` | ✅ ada di ekspor |
| pembulatan & tipe uang di basis data | `[data DBA]` — **P1** | ⛔ masih terbuka, **terpisah dari P18** |
| apakah 2 % / 2,2 % boleh berubah menurut waktu | `[Product+Underwriting]` | ⚠️ pertanyaan baru, bukan P18 |

### 5.4 `DataTransform` dan `FlowAction` — pola yang sama atau tidak?

| Jenis aturan | NB | EDM | Pola |
| --- | ---: | ---: | --- |
| `Activity` | 92 berkas | 66 berkas | ⭐ `PropertiesName` / `PropertiesValue` **tanpa awalan** |
| `DataTransform` | 12 berkas | 10 berkas | ⚠️ **kebalikannya** — `pyPropertiesName` / `pyPropertiesValue` **dengan awalan**; bentuk tanpa awalan **nol** |
| `FlowAction` | 9 berkas | 6 berkas | ⛔ **tidak memakai pasangan sama sekali** — ia merujuk `pyStreamName`, bukan menetapkan nilai |
| `When` | 75 berkas | 11 berkas | `pyParametersParamName` / `pyParametersParamValue` — nama ketiga lagi |

`DataTransform` **terbaca**: **53** pasangan nama+nilai di NB *(dari 80 `pyPropertiesName` terisi)*
dan **49** di EDM. Contohnya `.PolicyTreatyIn.OperatorName = OperatorID.pyUserIdentifier`.

⚠️ **Tetapi `DataTransform` belum disensus setara** — belum diuji keutuhannya, belum dicocokkan
dua cara. Itu `[terbuka]`, dan sengaja tidak saya paksakan masuk vonis ronde ini.

---

## Bab 6 — ⭐⭐ Vonis: P18 masih diperlukan, atau tidak?

# ⭐⭐ TIDAK DIPERLUKAN

Bukan "diperlukan sebagian". **Tidak diperlukan.** Alasannya tunggal dan terukur: **seluruh** yang
diminta P18 — isi 268 langkah penetapan nilai pada 51 aturan — **sudah ada di dalam ekspor yang kami
pegang**, sejumlah **707 pasangan nama=nilai terisi pada 269 langkah**, tanpa satu pun tanda
terpotong. Tidak ada yang perlu dikirim ulang.

⭐ **Ronde ini berhasil dengan cara yang paling mahal bagi saya: ia membuktikan sesi asisten
sebelumnya — yaitu saya — keliru selama tiga ronde.** Itu tetap hasil yang benar, dan jauh lebih
murah diketahui sekarang daripada sesudah suratnya telanjur dikirim.

### 6.1 Apa lagi yang berubah bila vonis ini diterima

⛔ **Tidak satu pun berkas di bawah ini saya sunting** — ronde ini melarangnya. Ini daftar kerja
untuk ronde berikutnya, sesudah work owner memutuskan.

| Berkas | Yang harus berubah |
| --- | --- |
| ⛔⛔ `SURAT-P18-KE-PEMILIK-EXPORT-PEGA.md` | **dibatalkan, jangan dikirim.** Suratnya meminta barang yang sudah kami punya |
| ⛔⛔ `LAMPIRAN-P18-ATURAN-DIMINTA.md` | **dibatalkan** bersama suratnya |
| ⛔ `spec.md` **§8.2** *(Karantina rantai uang)* | karantina **dibuka**. "Empat akibat bila P18 tidak dijawab" gugur seluruhnya |
| ⛔ `spec.md` baris 51, 113, 805-807, 881, 949 | P18 tidak lagi menahan; butir `[terbuka]` yang MENAHAN turun dari **2** ke **1** *(tinggal P1)* |
| ⛔ `spec.md` **AC 79** | *"rantai perhitungan tidak dibangun sebelum P18 dijawab"* — **dicabut**; pekerjaannya boleh mulai |
| ⛔ `issues/13-rantai-perhitungan-uang-dikarantina.md` | `Status: blocked` → **siap dikerjakan** atas sisi P18. ⚠️ **P8 tetap menahan sebagian** |
| ⛔ `issues/12-layar-jenjang-ketiga.md` | `Status: blocked` → **siap dikerjakan**. Dalil "layar tidak dapat disimpan" gugur *(Bab 4)* |
| ⛔ `issues/00-PETA-AC.md` baris 83 | "Tertahan P18: tiket 12 dan 13" → **nihil** |
| ⛔ `issues/07`, `issues/15` | rujukan "yang menunggu P18" disesuaikan |
| ⛔ `PERTANYAAN-untuk-pemilik-export-Pega.md` | ⚠️ P18 **ditarik oleh tim migrasi, bukan dijawab**. Lembar itu tinggal 3 pertanyaan |
| ⛔ `KEADAAN-NB-TREATY-IN.md` | pernyataan tentang isi langkah yang hilang perlu diralat |
| ⛔ `grilling-ronde-2.md`, `-3.md`, `-4.md` | ⚠️ **tersegel — jangan disunting.** Ralatnya hidup di berkas ini |
| ⛔ `CLAUDE.md` §4a | tambahkan **jebakan #7**: *nama tag berbeda menurut jenis aturan; periksa daftar tag yang benar-benar ada sebelum menyimpulkan nol* |

### 6.2 Neraca pertanyaan sesudah vonis ini

| | Sebelum | Sesudah |
| --- | ---: | ---: |
| nomor terpakai P1-P60 | 60 | 60 |
| ditarik | 2 *(P41, P48)* | ⭐ **3** *(+P18)* |
| dapat dijawab | 58 | **57** |
| terjawab | 45 | 45 |
| ⭐ masih terbuka | 13 | ⭐ **12** |
| ⛔ **yang MENAHAN tiket** | **2** *(P1, P18)* | ⭐ **1** *(P1 saja)* |

⭐⭐ **Akibat terbesarnya:** penghalang terberat proyek ini hilang, dan `[data DBA]` **P1** naik
menjadi **satu-satunya** pertanyaan yang menahan pekerjaan yang sudah tertulis di 16 tiket.

### 6.3 ⚠️ Satu hal yang TIDAK boleh disimpulkan dari ronde ini

⛔ Vonis ini berkata **isi langkah ada dan terbaca**. Ia **tidak** berkata perhitungan lama itu
**benar**, **tidak** berkata ia boleh disalin apa adanya ke Go, dan **tidak** menggantikan
`[keputusan work owner]` mana pun. Pertanyaan seperti *"apakah pembagi 102,2 masih berlaku untuk
tahun berjalan"* adalah pertanyaan **baru** untuk `[Product+Underwriting]`, bukan P18 yang
dihidupkan kembali.

---

## TELEMETRI EKSEKUSI

### Invarian korpus — terukur pada ronde ini

| Modul | Berkas | Byte | md5 gabungan |
| --- | ---: | ---: | --- |
| `NB Treaty In` | 278 | 35.492.317 | `96271809212b68bde10ae8612f414baf` |
| ↳ folder `Activity` | 92 | — | `9331367d84437a97c461c30d2b0b1998` |
| `EDM Treaty In` | 163 | **17.639.259** | `24afd5726e083c9df6a692da94728139` |
| ↳ folder `Activity` | 66 | — | `22cf762f4ea4c513ef0e109dfedf94fb` |

⚠️ **Dicatat tanpa ditafsirkan:** ukuran `EDM Treaty In` terbaca **17.639.259 B**, sedangkan catatan
EDM ronde 1 menulis **17.639.246 B** — selisih **13 B**. ⛔ Saya **tidak** menyimpulkan sebabnya.
Menyimpulkan riwayat dari struktur adalah kekeliruan yang sudah dua kali saya lakukan *(P48, P41)*
dan dua kali diralat oleh work owner. Butir ini `[terbuka]`.

### Cara tiap angka diperoleh

| Angka | Cara A | Cara B | Sepakat? |
| --- | --- | --- | --- |
| pasangan terisi | `ElementTree`, sarang dihormati, pasangan dalam satu `rowdata` | pola teks murni, tanpa `ElementTree` | ✅ sesudah 5 selisih dijelaskan |
| langkah `Property-Set` | `pySteps` **puncak saja**, rekursi ke sub-langkah | cacah metode per langkah | ✅ 939 / 380 |
| keutuhan nilai | tujuh uji struktural | sebaran panjang + batas mencurigakan | ✅ nol tanda potong |
| cukupan per langkah | `pyParamArray` milik langkah sendiri | cacah aturan bernilai nol | ✅ 0 aturan kosong |

### Kegagalan alat pada ronde ini — dicatat, bukan disembunyikan

| # | Kegagalan | Ditangkap oleh | Akibat bila lolos |
| ---: | --- | --- | --- |
| 1 | Regex tag-berpasangan tertelan `<pagedata>` sepanjang 65.237 aksara → **melaporkan "berisi 0" untuk tag yang jelas berisi** | uji saya sendiri, satu langkah kemudian | ⛔⛔ **akan mengukuhkan P18 yang keliru** |
| 2 | `iter("pySteps")` menelusuri sarang berulang → 1.375 langkah, 3.458 pasangan | ⭐ perselisihan antara cara A dan cara B | angka langkah 46 % terlalu besar |
| 3 | "56 petik tunggal ganjil" sempat terbaca sebagai tanda terpotong | pembacaan nilainya satu per satu | ⛔ vonis "sebagian terpotong" yang palsu |

⭐ **Kegagalan #1 adalah pelajaran utama ronde ini**: ia **bentuk yang sama persis** dengan
kekeliruan yang melahirkan P18 — alat yang salah menghasilkan nol, dan nol itu dipercaya. Bedanya,
kali ini ia tertangkap karena brief memaksa **memeriksa daftar tag yang benar-benar ada lebih
dulu**, bukan menebaknya.

### Ongkos ronde ini

| Ukuran | Nilai |
| --- | ---: |
| panggilan model | **63** |
| token keluar | **188.485** |
| token cache ditulis | **461.941** |
| token cache dibaca | **6.605.809** |
| berkas korpus dibaca | 441 *(278 NB + 163 EDM)* |
| berkas keluaran | **1** — berkas ini |
| berkas korpus disunting | **0** |
| pertanyaan terbuka ditutup sendiri | **0** |

⚠️ Angka token diukur dari catatan sesi sebagai **selisih** terhadap titik awal ronde. Jendela hitung
baris: **seluruh berkas dikurangi bab TELEMETRI EKSEKUSI** = **400 baris**. ⚠️ Angka token
sejati tidak terlihat dari dalam sesi; ini pendekatan terbaik yang tersedia.

### Batas ronde ini

- ⛔ Tidak menutup P18. Penarikannya **keputusan work owner**.
- ⛔ Tidak menyunting `spec.md`, tiket, `KEADAAN-...`, lembar pertanyaan, atau berkas grilling.
- ⛔ Nol kode, nol DDL, nol usulan kolom. Cuplikan di Bab 5 adalah **kutipan korpus apa adanya**.
- ⛔ Nol nama orang, nol nomor polis, nol nilai kredensial, nol data produksi.
- ⚠️ Vonis ini berlaku untuk `Activity`. `DataTransform` memakai nama tag berbeda *(Bab 5.4)* dan
  **belum** disensus setara. `[terbuka]`
