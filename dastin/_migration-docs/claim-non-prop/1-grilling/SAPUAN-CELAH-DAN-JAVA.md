# Penutupan celah cakupan, Java tertanam, dan tiket `05` dan `06`

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, **12 jenis rule**.
> Folder `Komite Claim Non Prop` **tidak dibuka** — dinyatakan tertutup untuk batch ini pada 19 September 2026, bukan ditunda. Setiap pernyataan nihil di dokumen ini berbatas pada folder `Claim Non Prop`, dan batas itu ditulis, bukan disimpulkan pembaca.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut.

Dijalankan 19 September 2026, lima langkah berurutan sesuai arahan Q5.

**Gerbang ADR-0028 terangkat** — ia `accepted` sejak hari ini. **Gerbang penamaan berdiri**: tidak ada nama constraint, nama index, maupun singkatan di dokumen ini.

---

# 1. Tag yang dilipat ke indeks

Lima belas tag baca dan dua pola tulis, seluruhnya dari daftar celah pemeriksaan cakupan.

| Kelompok | Tag |
|---|---|
| kondisi `Rule-Obj-When` | `pyLogic`, `pyConditionFieldName`, `pyConditionOperation`, `pyConditionLabel`, `pyDesignatedProperty`, `pyParametersParamValue` |
| penjaga transisi langkah | `pyStepsTransParamsWhen` |
| Page-Copy | `CopyFrom`, `CopyInto` |
| sasaran dan penyaring | `Property`, `pyPropertyTarget`, `pyFilterName`, `Field`, `Message` |

Dua pola tulis baru: `CopyFrom` → `CopyInto` (sasarannya yang **kedua**, bukan yang pertama — arah ini diperiksa terhadap XML sebelum dipakai) dan `pyPropertyTarget` → `pyPropertyValue`.

Lima tag menyebut properti **tanpa titik di depan**; untuk kelimanya nama telanjang ikut dihitung sebagai pembacaan: `pyDesignatedProperty`, `pyConditionFieldName`, `pyFilterName`, `Property`, `Field`.

| | Sebelum | Sesudah |
|---|---|---|
| nama ditulis | 721 | **733** |
| nama dibaca | 845 | **921** |
| berkas terbaca | 279 | 279 |

---

# 2. S11, S12, S16 dijalankan ulang — selisihnya

Dijalankan sekali di atas indeks utuh. **Yang dilaporkan selisihnya, bukan hasil barunya.**

## S11 — penanda yatim

| Golongan | Lama | Baru | Selisih |
|---|---|---|---|
| 1 — ditulis **dan** dibaca | 40 | **41** | `AktifButton` masuk |
| 2 — ditulis, tidak dibaca | 14 | **13** | `AktifButton` keluar |
| 3 — dibaca, tidak ditulis | 48 | **49** | `VirusCheckStatus` masuk |

**Satu klaim yatim gugur, dan dicabut di tempatnya.** `AktifButton` ternyata **dibaca** — `Activity\ProteksiSendKomiteCNP_Act.xml` baris 4491 dan 4969, lewat tag `<Field>.AktifButton</Field>`, yaitu medan Report Definition. Tag itu tidak pernah dibaca indeks lama. Ia sekarang golongan 1.

**Satu masuk, dengan keterangan.** `VirusCheckStatus` muncul di `AttachRISlipToWork` baris 874 sebagai `<pyStepsTransParamsWhen>@startsWith(param.VirusCheckStatus,"Virus")</pyStepsTransParamsWhen>` — persis kelas tag yang hilang. Tetapi ia **parameter**, bukan properti objek kerja, jadi ia bukan penanda yatim dalam arti yang sama. Dicatat supaya tidak salah dibaca sebagai cacat.

**Yang TIDAK bergerak, dan inilah yang penting:**

| Penanda | Keadaan sesudah jaring diperlebar |
|---|---|
| `IsEditClaim` | ditulis **3**, dibaca **0** |
| `CNPStatusCase` | ditulis **3**, dibaca **0** |
| `IsAnyAcceptation` | ditulis **0**, dibaca **12** |

Ketiganya bertahan. Klaim EVIDENCED-NIHIL atasnya **menguat**, karena kini diuji terhadap 45 tag alih-alih 30.

## S12 — ketidakkonsistenan tipe

**22 properti, tidak berubah sama sekali.** Nol masuk, nol keluar. Tag yang dilipat tidak memuat perbandingan literal baru.

## S16 — properti Pega

46 → **51 nama**. Lima masuk, nol keluar:

| Properti | Lewat tag | Kelasnya |
|---|---|---|
| `pxAttachedBy` | `<Property>` | **kolom pelaku** — kelas D32 |
| `pxAttachKey` | `<Property>` | pegangan lampiran |
| `pxAttachName` | `<Property>` | nama lampiran |
| `pxRefObjectKey` | `<Property>` | pegangan objek rujukan |
| `pyAttachStream` | `<Property>`, `PropertiesValue` | isi lampiran |

**`pxAttachedBy` memperbesar D32.** Daftar kolom pelaku yang tersimpan di properti `px` kini bertambah satu: siapa yang melampirkan dokumen. Membuang `px` mentah bukan hanya menghilangkan pencatat baris tertanam, tetapi juga **pencatat lampiran**.

---

# 3. 480 baris Java dibaca sebagai kode

**28 blok, 18 berkas, 480 baris.** Urutan pencarian sesuai arahan.

## 3.1 `double` dan `float` — **nol, di mana pun**

| Dicari | Hit |
|---|---|
| `double`, `float`, `Double`, `Float` | **0** |
| `parseDouble`, `parseFloat` | **0** |
| `BigDecimal` | 0 |
| `String` | 92 |
| `int` | 8 |

**Tidak ada tipe biner pecahan di jalur uang.** Ini hasil bersih dan dicatat tanpa diminta.

**Tetapi ketiadaan `double` bukan berarti presisi terjaga.** Uang melintasi Java sebagai **`String`**, 92 kali, lewat `.toString()` atas properti Pega. Tidak ada `BigDecimal` sama sekali. Java di sini tidak menghitung uang — ia **memindahkannya sebagai teks**, dan itu masalah yang berbeda, bukan masalah yang tidak ada. Lihat 3.3.

## 3.2 Tujuh blok adalah satu rutin yang disalin-tempel — dan ia **menghapus baris**

Tujuh dari 28 blok adalah rutin yang sama: membaca sebuah PageList dari belakang, merangkai beberapa properti jadi satu kunci berpemisah `#`, memasukkannya ke `HashSet`, dan **membuang baris yang kuncinya sudah ada**.

```java
boolean isAdded = values.add(aValue);
if(!isAdded){
  KurssList.remove(i);
}
```

Komentar `// read all 10 properties from aPage` dipertahankan di ketujuhnya meski yang dibaca 1 sampai 4 properti — bukti salin-tempel yang tidak diedit.

**Kunci alami yang ditegakkan di Java, bukan di data:**

| Rule | Halaman | Kunci dedup |
|---|---|---|
| `AddAkseptasiCNP_Act` | `Kurs` | `(CurrencyID, Currency)` |
| `AdjClaimAmount_Act` | `tempCurrency` | `(Currency)` |
| `AdjClaimCNP_Act` | `tempCurrency` | `(Currency)` |
| `CountSpreading_Act` | `Estimate` | `(CurrencyID, Currency, TypeLossID, TypeLoss)` |
| `CountTotalInsterest_Act` | `Cuan` | `(CurrencyID, Currency)` |
| `CountTotalInsterest_Act` | **`.ClaimData.TotalInterestInsured`** | `(CurrencyID, Currency)` |
| `GetHistoryMasterID_NP` | `TempTotal` | `(XOL, Currency)` |

**Dua hal yang perlu dinyatakan terpisah.**

Pertama, ini **daftar kunci alami berbukti** untuk enam kumpulan yang akan menjadi tabel. Sebelumnya kunci-kunci itu disimpulkan dari pemakaian; sekarang terbaca sebagai kode.

Kedua, dan lebih berat: **satu di antaranya membuang baris dari PageList milik objek kerja**, bukan dari halaman sementara. `.ClaimData.TotalInterestInsured` adalah data klaim, dan barisnya dihapus diam-diam saat rule berjalan. Tidak ada pencatatan, tidak ada penanda, tidak ada pesan. Baris yang hilang tidak meninggalkan jejak.

## 3.3 Muatan kasir dirakit dengan tangan — dan angkanya tidak berkutip

`HitServiceToKasir_Act`, dua blok identik 54 baris, membangun JSON dengan penyambungan teks:

```java
String Nett = DataJSON.getProperty(".CARI9").toString();
...
+"\"Nett\":"+Nett+","
+"\"Deductible\":"+Deductible+","
+"\"KaliDeduct\":"+KaliDeduct+","
```

Ketiganya **tanpa kutip** — dikirim sebagai angka JSON. Delapan belas medan lain berkutip sebagai teks.

**Akibat yang terbaca langsung:**

- Bila `CARI9` kosong, hasilnya `"Nett":,` — **JSON rusak**, dan kegagalannya terjadi di sisi penerima, bukan di sini.
- Tidak ada pemeriksaan bahwa isinya memang angka. `.toString()` atas properti Pega mengembalikan apa pun yang ada.
- Tidak ada kendali skala. Berapa desimal yang terkirim ditentukan cara Pega menuliskan propertinya, bukan oleh keputusan siapa pun.
- Tidak ada pelolosan karakter untuk 18 medan teks. Nama perusahaan bertanda kutip merusak muatan.

## 3.4 JSON asing diadopsi utuh ke objek kerja

Tiga rule memanggil `adoptJSONObject` atas teks dari kolom `HASIL1`:

| Rule | Halaman tujuan |
|---|---|
| `GeneratePlaCNP_Act` | `TempTreatyOut` |
| `GetDetailPolis_act` | `TempWorkPage.OfferTreatyIn` |
| `SetValueClaimTNP_Act` | `pyWorkPage.TreatyInMaster` |

Bentuk JSON-nya **tidak diperiksa**. Apa pun yang ada di `HASIL1` menjadi properti. Ini menyentuh **B4**: struktur JSON bukan hanya tidak terdokumentasi, ia **tidak dibatasi**.

`SetValueClaimTNP_Act` menulis ke `pyWorkPage.TreatyInMaster` — halaman objek kerja, bukan halaman sementara.

## 3.5 Panggilan keluar tanpa batas waktu

`GetBase64Attachment` membuka `java.net.URL` dari nilai parameter dan membaca seluruh isinya ke memori:

```java
String fileUrl = tools.getParamValue("FileURL");
java.net.URL url = new java.net.URL(fileUrl);
java.io.InputStream is = url.openStream();
```

Tanpa batas waktu, tanpa batas ukuran, tanpa daftar alamat yang diizinkan. Dicatat sebagai keadaan sistem lama; **bukan lingkup lapisan data**, dan tidak ditindaklanjuti di sini.

## 3.6 Java **tidak** membaca satu pun penanda yang dinyatakan tidak pernah dibaca

| Penanda | Kemunculan di 480 baris Java |
|---|---|
| `IsEditClaim` | **0** |
| `IsAnyAcceptation` | **0** |
| `CNPStatusCase` | **0** |
| `IsTreatyIn`, `FlagProrate`, `ReporterStatus`, `StsKatastrofe`, `AcceptanceStatus` | **0** |

Inilah yang diramalkan: **F2, H2, dan I2 menguat — bukan tertutup.** Java adalah tempat terakhir di folder terbuka yang dapat memuat pembaca tersembunyi, dan ia bersih. Dicatat di kolom bukti masing-masing.

---

# 4. Jenis rule yang terekspor tanpa isi

Diukur sebagai persentase karakter di luar tag metadata Pega, per jenis rule. Hasilnya **seragam, 5,5% sampai 11,8%** — tidak ada jenis yang menonjol kosong. Dugaan bahwa `Rule-Obj-Flow` atau `Rule-Declare-DecisionTable` terekspor lebih kosong dari yang lain **tidak terbukti**; keduanya justru di ujung atas.

## Tetapi satu berkas memang kehilangan isinya — dan klaim saya sebelumnya salah

> **KOREKSI.** Saya sebelumnya menulis `GetMimeType.xml` memuat **"nol sel keputusan"**. **Itu salah.** Pemeriksaan strukturnya menemukan **42 slot baris**, bernomor 1 sampai 42, hadir dalam tiga blok sejajar.

Yang benar: **42 baris ada, dan nilainya tidak ada.** Blok yang seharusnya memuat nilai tiap baris seluruhnya **menutup sendiri** — `<rowdata REPEATINGINDEX="1"/>` sampai `"42"/>`, empat puluh dua tag kosong berturut-turut. Yang tersisa di ekspor hanya kerangkanya, lebar kolom, dan satu nama parameter: `ext`.

Tabel ini memetakan **ekstensi berkas ke tipe MIME** pada class `ASM-FW-GISFW-Int-T_STORAGE_IMAGE`. Empat puluh dua pemetaan, dan tidak satu pun terbaca.

**Ini celah bahan, bukan celah metode.** Tidak ada pola sapuan yang dapat memulihkannya — isinya memang tidak dikirim. Diajukan sebagai **REQ-035**.

---

# 5. Tiket `05` dan `06`

## 5.1 `05` — sisi kasir, terpetakan lengkap

Inilah yang membuka tiket ini kembali: sisi Arasapas sudah terpetakan 17 parameter, sisi kasir belum sama sekali. Java menutupnya — nama bisnis tiap slot ada di dalamnya.

**Muatan kasir: 21 medan, dari 25 slot.**

| Slot | Nama di muatan | Bentuk JSON | Diisi dari |
|---|---|---|---|
| `CARI1` | `NoTrans` | teks | `Primary.AcceptedNo` |
| `CARI2` | `NoKlaim` | teks | `pyWorkPage.ClaimData.NoClaim` |
| `CARI3` | `LbuId` | teks | `OfferFacIn.QuotationData.BusinessOldId` |
| `CARI4` | `NoPolis` | teks | `ClaimData.PolicyData.PolicyNo` |
| `CARI5` | `AcceptType` | teks | `Primary.PaymentType` |
| `CARI6` | `Kepada` | teks | `Primary.PayableTo` |
| `CARI7` | `AccountNo` | teks | `.NoAccount` dengan non-angka dibuang |
| `CARI8` | `TglAksep` | teks | `Primary.AcceptedDate` dirakit ulang |
| **`CARI9`** | **`Nett`** | **ANGKA** | **`.TotalClaim - .PremiumSpreaded` + satu besaran menurut `PaymentType`** |
| `CARI10` | `Deductible` | **ANGKA** | **`0` tertanam** |
| `CARI11` | `KaliDeduct` | **ANGKA** | **`0` tertanam** |
| `CARI12` | `StsSyariah` | teks | **`0` tertanam** |
| `CARI13` | `CompanyName` | teks | **`"NUSARE"` tertanam** |
| `CARI14` | `LjtdId` | teks | **`"D0031"` tertanam** |
| `CARI15` | `LdcId` | teks | **`"100081"` tertanam** |
| `CARI16` | `StsAp` | teks | **`"0"` tertanam** |
| `CARI17` | `LkuId` | teks | `.CurrencyID` |
| `CARI18` | `LbgID` | teks | `.IDOfBank` |
| `CARI22` | `TglBolehBayar` | teks | rakitan `CARI19`-`CARI20`-`CARI21` |
| `CARI23` | `Email` | teks | `Email.pxResults(1).HASIL1` |
| `CARI24` | `UserInput` | teks | `.pxCreateOperator`, jatuh ke `pxRequestor.pyUserIdentifier` bila kosong |

**Tiga slot tidak ikut muatan**: `CARI19`, `CARI20`, `CARI21` — hari, bulan, tahun, dipakai hanya untuk merakit `CARI22`. Ditambah `CARIDATETIME` sebagai penampung sementara.

### Apa yang dibuka pemetaan ini

**Enam nilai tertanam di kode.** `Deductible`, `KaliDeduct`, `StsSyariah`, `CompanyName`, `LjtdId`, `LdcId`, `StsAp` tidak pernah berasal dari data. Dua di antaranya berat:

- **`StsSyariah` selalu `0` dan `CompanyName` selalu `"NUSARE"`.** Muatan ke kasir **tidak pernah menyatakan syariah**, apa pun keadaan klaimnya. Ini menyentuh **G1** dari arah baru: pembedaan syariah ada di `pyNotifyAccountName` untuk surel, tetapi **tidak ada sama sekali** di jalur pembayaran.
- `LjtdId` `"D0031"` dan `LdcId` `"100081"` adalah pengenal sistem hilir yang dipatri. Bila keduanya berubah di sisi kasir, tidak ada yang memberi tahu sistem ini.

**Satu hal yang saya tidak nyatakan sebagai temuan, karena tidak dapat dipastikan dari ekspor**: `CARI20` berbunyi `@toDecimal(@month(TempKasir.CARIDATETIME)) + 1`. Penambahan satu itu benar bila `@month()` berbasis nol, dan menghasilkan bulan ke-13 bila berbasis satu. Perilaku fungsi bawaan Pega **tidak ada di ekspor**. Diajukan sebagai **REQ-036**, bukan disimpulkan.

**Kesimpulan tiket `05`**: sisi Arasapas 17 parameter bernama, sisi kasir 21 medan bernama dari 25 slot. Keduanya terbaca. **Tiket ditutup.**

## 5.2 `06` — keadaan kasus, dengan batas tertulis

Dijalankan di atas indeks utuh 45 tag.

| | Hasil |
|---|---|
| Nilai yang benar-benar ditulis | **dua**: `"COMITEE ACCEPTANCE (DEPT. HEAD)"`, `"INPUT ACCEPTATION CLAIM"` |
| Ditulis di | `CreateChildKomiteCNP_Act` 13793, `CreateChildKomiteCloseNP_Act` 2802, `InputOutStandingCTNP_PostAct` 400 |
| Dibaca | **nol** |
| Dibandingkan | **nol** |
| `"CLAIM ACCEPTED"` / `"CLAIM REJECTED"` | **nol kemunculan, di mana pun** |

### Jawabannya, beserta batasnya

> **Tidak ada pembaca `CNPStatusCase` di folder `Claim Non Prop`**, diuji terhadap 45 tag pada 279 berkas dan 12 jenis rule, ditambah 480 baris Java tertanam.

Batas itu bagian dari jawaban, bukan tambahan. Folder `Komite Claim Non Prop` tertutup untuk batch ini, dan ketiga penulisnya ada di rule yang berurusan dengan Komite — jadi pembacanya, bila ada, kemungkinan besar justru di sana.

**Yang dapat dinyatakan tanpa syarat**: `CNPStatusCase` **bukan enum yang perlu dilengkapi**. Ia kolom teks yang diisi tiga tempat dan tidak mengatur apa pun di dalam modul ini. Model memperlakukannya sebagaimana kebijakan **E17**: dimigrasi apa adanya, ditandai tak berpemilik, tanpa domain tertutup.

**Ejaan tersimpan `COMITEE`, satu T.** Ejaan `COMMITTEE` dua T muncul 10 kali sebagai teks tampilan di 7 berkas. Yang masuk basis data adalah yang salah eja, dan migrasi memindahkannya apa adanya.

**Tiket ditutup**, dengan batas di atas tertulis di badan jawabannya.

---

# 6. Kode D baru

| Kode | Temuan | Menyentuh |
|---|---|---|
| **D37** | `AktifButton` **bukan** penanda yatim — dibaca lewat `<Field>` di `ProteksiSendKomiteCNP_Act` 4491, 4969. Klaim lama dicabut | S11, E17 |
| **D38** | `IsEditClaim`, `CNPStatusCase`, `IsAnyAcceptation` **bertahan** setelah jaring diperlebar dari 30 ke 45 tag dan setelah 480 baris Java dibaca | F2, H2, I2, `17`, `31` |
| **D39** | **Nol `double`, nol `float`, nol `BigDecimal`** di 480 baris Java. Uang melintas sebagai `String`, 92 kali | ADR-0003 |
| **D40** | Tujuh blok Java adalah satu rutin dedup yang disalin-tempel, menegakkan **kunci alami di kode**. Salah satunya membuang baris dari `.ClaimData.TotalInterestInsured` — **data klaim, tanpa jejak** | `16`, `17`, `18`, `19` |
| **D41** | Muatan kasir dirakit dengan penyambungan teks; `Nett`, `Deductible`, `KaliDeduct` **tanpa kutip**. Nilai kosong menghasilkan JSON rusak | `25`, `29`, `05` |
| **D42** | Tujuh nilai muatan kasir **tertanam di kode**, termasuk `StsSyariah = 0` dan `CompanyName = "NUSARE"` — jalur pembayaran **tidak pernah menyatakan syariah** | G1, `25`, `29` |
| **D43** | `adoptJSONObject` di tiga rule mengadopsi JSON **tanpa pemeriksaan bentuk**; satu di antaranya ke halaman objek kerja | B4 |
| **D44** | `pxAttachedBy` adalah **kolom pelaku lampiran** — memperbesar D32 | `20` |
| **D45** | **KOREKSI.** `GetMimeType` bukan "nol sel"; ia **42 slot baris yang nilainya tidak ikut terekspor**. Celah bahan, bukan celah metode | REQ-035 |
| **D46** | Sisi kasir terpetakan lengkap: **21 medan bernama dari 25 slot** `CARI`; tiga slot hanya perakit tanggal | `05`, `25`, `29` |
