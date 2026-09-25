# Titik buta pohon 985 simpul — ditemukan 24 September 2026

**Ditemukan saat menulis `SPEC-MODEL-DATA.md` §10 untuk entitas `RINCIAN_PENYEBARAN`**, yaitu satu
dari tiga entitas yang didahulukan. Ia tidak dicari; ia menabrak pekerjaan.

> **Yang salah bukan penjaganya. Yang salah adalah apa yang dilakukan pada yang ditolak penjaga.**

---

## 1. Bagaimana ia terjadi

`alat/buat-pohon-treatyin.py` memuat penjaga ini, dan penjaga itu **benar**:

```python
# `Primary` hanya berarti TreatyIn bila ATURANNYA SENDIRI applies-to kelas itu.
m = re.search(r'<pyClassName>(.*?)</pyClassName>', s, re.S)
kelas_rule = m.group(1).strip().upper() if m else ''
primary_ok = (kelas_rule == KELAS_AKAR)
```

Tanpa penjaga itu, `Primary.SpreadingListXOL` dari sebuah aturan ber-*applies-to*
`Data-TreatyInShare` akan masuk sebagai **`TreatyIn.SpreadingListXOL` di tingkat akar** — simpul
palsu, karena tempatnya sebenarnya di bawah `Share[]`.

**Akibat yang tidak pernah diperiksa:** yang ditolak penjaga **tidak dipindahkan ke bawah kelas
anaknya**. Ia dibuang. Aturan yang applies-to kelas anak karena itu **tidak menyumbang satu simpul
pun**, di mana pun.

### Contoh yang membongkarnya

| | |
|---|---|
| Berkas | `Treaty In/Activity/SetSpreadingXOL.xml` |
| `pyClassName` pertama | `ASM-FW-GISFW-Data-TreatyInShare` — **bukan** kelas akar |
| Jalur di dalamnya | `Primary.SpreadingListXOL(idx).BreakDownSprdListXOL` |
| Di pohon 985 simpul | **nol** |
| Di ekspor | `BreakDownSprdList` **28 kali**, `BreakDownSprdListXOL` **26 kali**, di 6 berkas Treaty In |

`BreakDownSprdList` adalah **sumber entitas `RINCIAN_PENYEBARAN`** — `SPEC-MODEL-DATA.md` §3.5
menyebutnya dengan nama. Entitas itu selama ini berdiri di atas jalur yang tidak pernah terbaca
perkakas mana pun.

### Satu baris yang akan menangkapnya di hari pertama

Penjaga itu kini **melapor**, dan perkakas dijalankan ulang tanpa perubahan hasil (985 / 927 / 17
identik):

```
DITOLAK penjaga primary_ok: 584 rujukan jalur
  dari 16 berkas, 6 kelas aturan bukan-akar
     ASM-FW-GISFW-DATA-TREATYINSHARE              256
     ASM-FW-GISFW-DATA-TREATYINRETROSHARE         160
     ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL       146
     ASM-FW-GISFW-DATA-TREATYINLIMITS              12
     ASM-FW-GISFW-DATA-TREATYINSHAREREINS           8
     ASM-FW-GISFW-DATA-TREATYINEGNPI                2
```

Tiga angka yang tampak berbeda dan ketiganya benar: **584 rujukan jalur** ditolak (satu properti
dapat dirujuk berkali-kali), yang menghasilkan **414 pasangan (kelas, properti)** dan **377 nama
properti berbeda** yang tidak ada di pohon.

> **Aturan yang ditetapkan pemilik proses, dan berlaku untuk perkakas mana pun:**
> **perkakas yang menolak sesuatu harus melaporkan apa yang ditolaknya dan berapa banyak.**
> Penolakan yang diam adalah bentuk perkakas dari *"diam bukan bukti"* — yang dibuang tanpa hitungan
> membuat yang tersisa terbaca sebagai keseluruhan.

---

## 2. Seberapa besar — diukur, bukan diduga

Perkakas baru `alat/sapu-properti-per-kelas.py` memakai **sumber yang mandiri dari penyisiran
jalur**: entri indeks rujukan aturan, yang memasangkan `pxRuleClassName` (kelas pemilik properti)
dengan `pyRuleName` (nama propertinya). Ia tidak pernah melewati awalan `TreatyIn.` sama sekali,
sehingga ia tidak dapat mewarisi cacat yang sama.

| | |
|---|---:|
| Berkas disapu (Treaty In + Adjustment) | 708 |
| Kelas yang punya properti | 56 |
| Pasangan (kelas, properti) | 866 |
| **Properti yang TIDAK ADA di pohon sama sekali** | **414** |

Angka 414 mencakup kelas yang memang di luar modul (pencarian, penawaran, master). **Dibatasi pada
14 kelas yang menyandang 27 entitas penyerahan pertama, angkanya 74**:

| Kelas Pega | Entitas | Hilang |
|---|---|---:|
| `Data-TreatyInLimitsDetail` | `DETAIL_PROPORSIONAL` | **38** |
| `Data-TreatyInLimits` | `LAYER` | **16** |
| `Data-TreatyInShare` | `BAGIAN` | 6 |
| `Data-TreatyInCurrencyList` | `MATA_UANG_KONTRAK` | 3 |
| `Data-TreatyInAccumulation` | `PERIODE_AKUMULASI` | 2 |
| `Data-TreatyInCoInScale` | `SKALA_KOASURANSI` | 2 |
| `Data-TreatyInInstallment` | `TERMIN` | 2 |
| `Data-TreatyInLimitsSpreading` | `PENYEBARAN` / `RINCIAN_PENYEBARAN` / `NILAI_PENYEBARAN` | 2 |
| `Data-TreatyInAccountReport` | `PERIODE_PELAPORAN` | 1 |
| `Data-TreatyInDeduction` | `POTONGAN` | 1 |
| `Data-SuggestList` | `CATATAN_PERSETUJUAN` | 1 |
| `Data-TreatyInPortfolio` · `Data-TreatyInRetention` · `Data-TreatyInEGNPI` | — | 0 |
| | **Total** | **74** |

### Dari 74, hanya 41 menuntut keputusan

| Golongan | Jumlah | Perlakuan |
|---|---:|---|
| **PERLU DIPUTUSKAN** | **41** | diadili di §10 entitasnya masing-masing |
| agregat dan turunan (`Total…`, `Sum…`) | 19 | tidak disimpan — ADR-0037, §4.2 |
| pencapaian — GEL-3 | 7 | di luar gelombang ini |
| sudah ditambal §12.4 dari inventaris kelas | 7 | `SubDays`, `SubDueDate`, `PctLimit`, `CoInShare`, `Amount`, `ID`, `CurrencyID` |

**Tujuh yang sudah ditambal §12.4 adalah buktinya bahwa penambalan itu bekerja** — dan sekaligus
buktinya bahwa ia dijalankan **hanya untuk lima entitas yang kebetulan ditanyakan**, bukan untuk
seluruhnya. Perkakas ini menjalankannya untuk seluruhnya.

---

## 3. Apa yang ikut tersangkut

| Artefak | Akibatnya |
|---|---|
| `struktur-treatyin-lama.md` | angka **985 simpul** adalah jumlah yang **terjangkau dari akar**, bukan jumlah model. Batas itu sudah disebut §8 sebagai batas bukti, tetapi **tidak pernah diukur** |
| `peta-nama-tabel-treatyin.tsv` · `ERD-STRUKTUR-TREATYIN.html` | **tidak ada tabel datar untuk `BreakDownSprdList`/`…XOL`**. 44 tabel adalah 44 tabel *yang terlihat* |
| `SPEC-MODEL-DATA.md` §3.3, §3.4, §3.6 | daftar jalur masukan `LAYER`, `DETAIL_PROPORSIONAL`, dan `BAGIAN` tidak lengkap |
| **`SPEC-MODEL-DATA.md` §10.3 dan §10.4** | **dua entitas yang sudah dinyatakan selesai**. `LAYER` 13 atribut dan `DETAIL_PROPORSIONAL` 9 atribut disusun di atas daftar jalur yang berlubang |

> **§10.3 dan §10.4 dibuka kembali.** Itu bukan pilihan: menulis §10 untuk 23 entitas dengan metode
> yang lebih baik daripada metode yang dipakai keempat entitas pertama menghasilkan berkas yang
> **dua bagiannya tidak sebanding**, dan tidak ada yang akan tahu bagian mana.

---

## 4. Yang TIDAK berubah

- **Penjaga `primary_ok` tetap.** Mencabutnya mengembalikan simpul palsu di tingkat akar, dan itu
  kesalahan yang lebih buruk: simpul palsu **terbaca sebagai fakta**, sedangkan simpul hilang
  setidaknya tidak mengarang apa pun.
- **Keempat salinan rekursif** (`ValueDifference`, `ActualValue`, `OLDDATA`, `ValueBeforeProrate`)
  dan temuan bahwa `OLDDATA` tidak dideklarasikan **tidak tersentuh** — keduanya dibaca dari
  jalur akar, yang memang terjangkau.
- **Angka 27 entitas tidak berubah.** Yang berubah jumlah **atribut** per entitas, bukan jumlah
  entitasnya. Tidak satu pun dari 41 properti itu berbentuk entitas baru kecuali
  `BreakDownSprdList`, dan entitas itu **sudah ada** namanya: `RINCIAN_PENYEBARAN`.

---

## 5. Cara menjalankan ulang

```
python alat/sapu-properti-per-kelas.py
```

Keluarannya dua:

| Berkas | Isi |
|---|---|
| `datar-properti-per-kelas.csv` | 56 kelas, seluruh propertinya |
| `datar-titik-buta-pohon.csv` | yang tidak tergantung di kelasnya pada pohon, dengan kolom `ADA_DI_POHON_DI_TEMPAT_LAIN` |

Kolom terakhir itu memisahkan dua hal yang sangat berbeda dan mudah tertukar: **`ya`** berarti
propertinya ada di pohon, hanya tidak tergantung pada kelas itu — biasanya karena nama yang sama
dipakai dua kelas. **`TIDAK`** berarti ia benar-benar tidak pernah terlihat. Hanya yang `TIDAK`
yang dihitung sebagai titik buta.

---

## 6. Yang langsung terbuka di baliknya — dan ia menyentuh ketiga entitas yang didahulukan

Begitu `BreakDownSprdList` terbaca, isinya menjawab pertanyaan yang tidak sempat ditanyakan.

### 6.1 Kedua cabang meletakkannya di tempat yang sama

| Cabang | Jalur | Ruasnya |
|---|---|---|
| proporsional | `Limits.Detail[].SpreadingList[].BreakDownSprdList[]` | `ReinsID`, `ReinsName`, `SharePct`, `Currency`, `Amount` |
| non-proporsional | `Share[].SpreadingListXOL[].BreakDownSprdListXOL[]` | `ReinsTypeID`, `ReinsTypeName`, `SharePct`, `Currency`, `Amount` |

Sekilas kedua cabang memecah pada **sumbu yang berbeda**: sisi proporsional per **pihak**
(`ReinsID`), sisi non-proporsional per **jenis** (`ReinsTypeID`). Bila benar, `RINCIAN_PENYEBARAN`
tidak dapat menjadi satu entitas — persis bentuk yang §12.3 tolak untuk `BAGIAN`.

### 6.2 Tetapi labelnya berbohong, dan kondisinya yang menang

`Activity/SetSpreadName.xml` — *applies-to* `Data-TreatyInLimitsDetail`, sisi **proporsional** —
mengisi kedua ruas itu begini:

```
Primary.SpreadingList(idx).BreakDownSprdList(<APPEND>).ReinsName
    = pyReportContentPage.pxResults(idx).ReinsTypeName
Primary.SpreadingList(idx).BreakDownSprdList(<LAST>).ReinsID
    = pyReportContentPage.pxResults(idx).ReinsTypeID
```

**`ReinsID` dan `ReinsName` diisi dari `ReinsTypeID` dan `ReinsTypeName`.** Namanya menyebut
*pihak*; yang ditulis ke dalamnya *jenis*. Tidak ada asimetri sumbu — yang ada asimetri **nama**.

> Ini keluarga kesalahan yang sudah tercatat di `PENGETAHUAN.md`: **label tidak dipercaya, kondisi
> yang dibaca.** Kali ini labelnya bukan deskripsi langkah melainkan **nama properti**, dan ia
> menyesatkan dengan cara yang sama.

**Akibatnya pada `STRUKTUR-DATA.md` §1.4:** `RINCIAN_PENYEBARAN` tertulis *"rincian penyebaran per
pihak"* dengan kunci alami **pihak**. Keduanya **tidak terbaca dari sumber**. Yang terbaca: rincian
penyebaran per **jenis reasuransi**, satu tingkat di bawah `PENYEBARAN` yang sumbunya juga jenis
reasuransi.

### 6.3 Dan ia mungkin bukan masukan sama sekali

Ketiga ruas sisanya diisi di aktivitas yang sama:

```
BreakDownSprdList(<LAST>).SharePct = pyReportContentPage.pxResults(idx).Pct
BreakDownSprdList(<LAST>).Currency = .Currency
BreakDownSprdList(<LAST>).Amount
    = @divide(.Value * Local.SplitRNMShare * pyReportContentPage.pxResults(idx).Pct, 100, 20)
```

`pyReportContentPage.pxResults` adalah hasil **laporan atas tabel master**. Maka:

| Ruas | Sifatnya |
|---|---|
| `ReinsID` / `ReinsName` | **salinan** pengenal jenis dari master |
| `SharePct` | **salinan** persentase dari master |
| `Amount` | **turunan** — `nilai × bagian NuRe × persen master` |
| `Currency` | ikut nilainya |

**Tidak satu pun ruasnya diketik orang.** Pada sumbu **masukan versus turunan** — sumbu yang §12.3
pakai untuk memutuskan `BAGIAN` — `RINCIAN_PENYEBARAN` jatuh seluruhnya di sisi turunan.

Dan itu menjelaskan satu hal yang sudah tercatat sebagai pola di `PENGETAHUAN.md`: *jalur
perhitungan otomatis tidak memerlukan validasi — **selama** persentase di tabel master berjumlah
100.* **INV-50** ("jumlah persen penyebaran per sumbu pihak sama dengan 100") berdiri persis di atas
andaian itu.

### 6.4 Kedudukannya DIPUTUSKAN — dan bukan oleh saya, melainkan oleh ADR-0036

Saya semula menahan kedudukan `RINCIAN_PENYEBARAN` pada uji produksi yang belum terjangkau. **Itu
keliru, dan ADR yang sudah diterima menjawabnya.**

ADR-0036 tidak sekadar menetapkan asas "beku saat disetujui, hitung saat dibaca". Ia menyebut
penyebaran **dengan nama**, sebagai penerapan yang mengikat:

> Penyebaran adalah **turunan** selama kontrak belum disetujui, dan menjadi **fakta tercatat**
> begitu disetujui. Sesudah akseptasi, perubahan pada susunan baku **tidak menjalar**.
>
> Hasil yang dibekukan membawa **penunjuk** ke masukan yang dipakai … serta **asal-usulnya**
> (disemai dari acuan baku, diisi tangan, atau disemai lalu disunting).

Dan `SPEC-INVARIAN.md` §3 daftar K2 memeriksa **INV-50** sebagai syarat persetujuan akhir. Maka
persentase penyebaran **adalah** angka yang menjadi dasar persetujuan, ia wajib dibekukan pada
versinya, dan sesuatu yang dibekukan menuntut tempat untuk disimpan.

> **`RINCIAN_PENYEBARAN` BERDIRI**, apa pun hasil uji produksinya — sebagai **fakta terbukukan yang
> membawa penunjuk asalnya**, yaitu persis pengecualian yang **INV-58** sediakan. `NILAI_PENYEBARAN`
> berdiri bersamanya sebagai anaknya.

**Kunci alaminya JENIS REASURANSI, bukan pihak** — §6.2. `STRUKTUR-DATA.md` §1.4 diperbaiki.

### 6.5 Yang uji produksi tentukan hanyalah GOLONGANNYA

| Hasil uji | Artinya | Golongan kemampuannya |
|---|---|---|
| persentase **tidak pernah** menyimpang dari master | di sistem lama ia murni turunan | membekukannya adalah **PERUBAHAN** |
| persentase **pernah** menyimpang | ada yang menyuntingnya; ia sudah masukan sejak dulu | membekukannya adalah **PELESTARIAN** |

Golongan menentukan **cara menguji dan besar taksiran**. Ia **tidak** menentukan bentuk model.

Maka §10 untuk `RINCIAN_PENYEBARAN` dikerjakan **sekarang**, dengan golongan kemampuannya ditandai
**TERTAHAN** pada uji itu — bentuk G3 *"dipecah"*, diterapkan pada entitas alih-alih pada tiket.

Ujinya **dapat mematahkan, tidak dapat mengesahkan**: suntingan yang ditimpa perhitungan ulang tidak
menyisakan jejak. Sifat itu wajib tertulis di kepala ujinya sendiri.

### 6.6 Satu akibat lagi, di berkas yang berbeda — dan ia sudah diperbaiki

`ReinsID` berisi jenis, bukan pihak. Akibatnya bukan hanya pada `STRUKTUR-DATA.md` §1.4:

**INV-47 dan INV-50 menyebut "sumbu pihak"** — sumbu yang tidak ada di keluarga penyebaran. Sesi DDL
menulis `GROUP BY` dari kalimat itu, dan pengelompokan yang keliru **lulus pada data yang seharusnya
gagal**. Keduanya diperbaiki 24 September 2026; pemeriksaan dan hasilnya di `SPEC-INVARIAN.md`
§4.4a. Ke-866 pasangan kelas–properti disapu: identitas pihak hanya ada di kelas retro (GEL-2), di
master luar skema, dan pada satu ruas `ReisuredParticipant` di `Data-TreatyInLimitsDetail` yang
belum pernah terlihat.

---

## 7. Satu baris yang merangkum seluruh L-8

`PLA` berstatus *"kepanjangannya tidak ada di mana pun"* selama tujuh sesi, dan diberi nama
sementara yang sengaja tidak dapat dikirim ke produksi. Ia terpecahkan pada hari L-8 ditutup, dari
kelas `Data-TreatyInBusinessLimit` — salah satu dari 39 kelas yang tidak pernah muncul di pohon.

> **Ia tidak pernah tersembunyi. Ia hanya ada di kelas yang pohonnya tidak pernah lihat.**

Itu berlaku untuk seluruh 414. Tidak satu pun disembunyikan siapa pun, tidak satu pun hilang dari
ekspor, dan tidak satu pun menuntut sumber baru untuk ditemukan. Yang kurang hanyalah **satu jalan
kedua menuju hal yang sama** — dan satu baris laporan dari perkakas yang menolaknya.
