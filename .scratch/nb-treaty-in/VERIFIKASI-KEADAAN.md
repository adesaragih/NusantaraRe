# Verifikasi `KEADAAN-NB-TREATY-IN.md`
## Sepuluh pernyataan diuji — sembilan cocok, satu hanya separuh dapat diperiksa dari korpus

> ⛔⛔ **RALAT MENYELURUH ATAS BERKAS INI — P18 DITARIK.** `[penyimpangan sadar]` 2026-09-22
>
> Berkas ini memuat **4 pernyataan** yang bersandar pada premis ⛔ *"isi langkah penetapan nilai
> tidak ikut terekspor"* — di baris **221 · 223 · 224 · 288**. **Premis itu keliru.**
>
> Isi langkah **ada di dalam ekspor**, di tag `PropertiesName`/`PropertiesValue` — **tanpa awalan
> `py`**. Tim migrasi memeriksa varian **dengan** awalan, menemukan nol, dan mempercayainya.
> ⭐ **Sebabnya dua huruf.** **938 dari 939** langkah `Property-Set` membawa isinya sendiri;
> **2.481** pasangan nama=nilai terisi; **nol** tanda terpotong. Lihat
> `VERIFIKASI-P18.md` dan `KOREKSI-P18-DIJALANKAN.md`.
>
> ⚠️ **Kalimat-kalimat itu sengaja TIDAK disunting satu per satu**, sebab sebagian berada di dalam
> **jawaban yang sudah diberikan pemiliknya** — mengubah jawaban orang lain bukan wewenang tim
> migrasi. Ralat ini berlaku atas seluruhnya.

---


Modul: **hanya `D:\XML\RNM_BRD\NB Treaty In`**. Korpus **READ-ONLY**; menulis hanya ke
`OUTPUT_HASIL_RNM\`. ⛔ `D:\XML\nusantara-re\` **tidak disentuh**.
⛔ `KEADAAN-NB-TREATY-IN.md` **tidak disunting**. ⛔ `grilling-ronde-1..4.md` **tidak disunting**.

⭐ **Ronde ini menemukan kesalahan — dan seluruhnya milik ronde 1–4, bukan milik berkas keadaan.**

---

## Bab 1 — Hasil sepuluh pernyataan

| # | Pernyataan | Ditulis | Terukur | Cocok? | Cara |
| ---: | --- | ---: | ---: | :---: | --- |
| 1 | berkas dalam modul | **278** | **278** | ✅ | `glob('**/*.xml')` rekursif, Python |
| 2 | ukuran modul | **35.492.317** | **35.492.317** | ✅ | jumlah `os.path.getsize` atas berkas yang sama |
| 3 | `Property-Set` di 92 `Activity` | **938** | **938** | ✅ | ⭐ **dua cara** — teks tag utuh · pengurai `ElementTree` |
| 4 | `Page-Clear-Messages` | **16** | **16** | ✅ | ⭐ dua cara, sama seperti #3 |
| 5 | baris `BusinessType_DeT` | **36** | **36** | ✅ | `pyPropertyValues` dan `pyResults`, masing-masing **36 `rowdata`** |
| 6 | `Activity` terjangkau / yatim | **64 / 28** | **64 / 28** | ✅ | ⭐ **dua cara** — urai langkah `Call` · cari nama sebagai teks |
| 7 | `Property-Set` terjangkau / yatim | **377 / 561** | **377 / 561** | ✅ | ⭐ dua cara, sama seperti #6 |
| 8 | medan wajib `pyRequired`=true | **27** medan di **6** layar | **27** di **6** | ✅ | ⭐ **dua cara** — per layar · properti unik lintas layar |
| 9 | medan terkunci permanen | **38** — 36 Dept Head, 2 biasa | **38** — **36 · 2** | ✅ | `pyReadOnlyCondition` selalu-benar, pengurai |
| 10 | kolom view `POOLDATA.TREATYINDETAILJOINEDM` | **39**; **33/33** medan | ⚠️ **33/33 ✅** · ⛔ **39 tidak dapat diperiksa** | ⚠️ **separuh** | `pyFieldName` unik di `BrowseTreatyJoinEDM` |

### 1.1 ⭐ Kenapa angka ronde 1–4 berbeda — dan sebabnya tepat satu berkas

⭐ `[terverifikasi]` Berkas keadaan mencatat bahwa `Protection_Act` diperbaiki pada 22 September.
⭐ **Akibatnya dapat dihitung dan cocok sampai satuan byte:**

| | Nilai |
| --- | ---: |
| ukuran modul menurut ronde 2 dan 4 | 35.678.284 |
| ukuran modul sekarang | **35.492.317** |
| **selisih** | ⭐ **185.967** |
| `Protection_Act` versi Fac In *(salah)* − versi Treaty *(benar)* | ⭐ **343.271 − 157.304 = 185.967** |

✅ **Cocok tepat.** ⭐ `Protection_Act` kini **157.304 B** berkelas
**`ASM-FW-GISFW-Data-PolicyTreatyIn`** — versi Treaty yang benar.

⭐ **md5 korpus berubah**, dan itu wajar: `62a3735ebb2eabb788c8cd8abdda78ca` *(ronde 1–4)* →
**`96271809212b68bde10ae8612f414baf`** *(sekarang)*.
⇒ ⭐ **Seluruh selisih angka ronde 1–4 terhadap berkas keadaan terjelaskan oleh satu pertukaran
berkas**: `Property-Set` 942→**938** *(−4)*, `Page-Clear-Messages` 17→**16** *(−1)*.
⛔ **Bukan kekeliruan pengukuran** — melainkan korpus yang berubah di antara dua pengukuran.

### 1.2 ⛔⛔ Pernyataan 5 — ronde 4 SALAH, berkas keadaan BENAR

> ⛔ **Ronde 4 §0.3 menulis:** *"nilai dan hasil tiap baris tidak terekspor"*, dan
> *"`pyCondition` · `pyOrConditions` · `pyResults` ⛔ 0 berisi"*.

⛔ **Itu keliru, dan sebabnya wajib dicatat sebagai cara kerja, bukan sebagai angka.**
⭐ Ketiga tag itu **bukan tag berteks** — ia **wadah `PageList` berisi `rowdata`**. Saya mengujinya
dengan memeriksa **teks langsung di dalam tag**, sehingga hasilnya nol.

`[terverifikasi]` Yang sebenarnya ada di `DecisionTable\BusinessType_DeT.xml`:

| Wadah | `rowdata` |
| --- | ---: |
| ⭐ `pyPropertyValues` | ⭐ **36** |
| ⭐ `pyResults` | ⭐ **36** |
| `pyCondition` | **36** *(dua wadah)* |
| `pyOrCondition` | 39 · 21 · 20 · 10 · 7 · 4 · 3 · 2 · 1 |

⭐ **Dan hasilnya terbaca penuh** — 36 nilai keluaran, contoh sepuluh pertama:
`"Medicare"` · `"PA"` · `"Bonding"` · `"BondingKBG"` · `"HE"` · `"MarineHull"` · `"Glass"` ·
`"Liability"` · `"AllRisk"` · `"AviationHull"`.

⭐⭐ **Dan berkas keadaan bahkan menjelaskan kekeliruan saya lebih tepat daripada saya sendiri:**
§5-nya menulis *"`pyRowNum` berbasis nol menandai baris ber-daftar-OR, bukan nomor baris"*.
✅ **Terkonfirmasi:** ke-9 nilai `pyRowNum` **−1, 2, 3, 5, 7, 22, 29, 30, 34** berjumlah **9**, dan
wadah `pyOrCondition` juga berjumlah **9** *(di luar satu wadah ber-1 rowdata)* — ⭐ keduanya
**menandai baris yang punya daftar-OR**, bukan mencacah baris.

⚠️ **Ini kegagalan instrumen kesepuluh sepanjang lima ronde, dan yang KEDUA yang tidak saya tangkap
sendiri.** ⛔ Lebih buruk: **brief ronde 4 sudah menyebut ketiga tag itu**, dan **jebakan #2 brief
ini sudah memperingatkan justru tentang `pyRowNum` berbasis nol.** ⭐ Saya membaca peringatannya,
lalu tetap salah — karena menguji **keberadaan teks**, bukan **keberadaan wadah**.

### 1.3 ⭐ Pernyataan 6 dan 7 — lebih baik daripada yang brief khawatirkan

⚠️ Brief memperingatkan bahwa 64/28 mungkin **batas atas**, sebab penyusunnya *"mencari nama
aturan sebagai teks"*. ⭐ **Saya uji dengan cara yang lebih tepat, dan hasilnya justru sama persis:**

| Cara | `Activity` terjangkau / yatim | `Property-Set` terjangkau / yatim |
| --- | ---: | ---: |
| ⭐ **A — urai langkah `Call`**, telusuri dari titik masuk | ⭐ **64 / 28** | ⭐ **377 / 561** |
| B — cari nama aturan sebagai **teks** di berkas lain | **83 / 9** | **878 / 60** |
| **ditulis berkas keadaan** | **64 / 28** | **377 / 561** |

✅ **Berkas keadaan cocok dengan cara A**, bukan cara B.
⭐ ⇒ **Angkanya BUKAN batas atas.** ⛔ Justru cara teks yang melonggar — ia menyatakan **19 aturan
tambahan** terjangkau yang **tidak pernah benar-benar dipanggil**, dan menggeser `Property-Set`
terjangkau dari 377 menjadi 878.

⭐ **Cara A, apa adanya:** titik masuk = **36** `Activity` yang namanya muncul di `Flow`, `Section`,
`Harness`, `FlowAction`, atau `DataTransform`; lalu diikuti lewat langkah ber-pola
`call [<Kelas>.]<Nama>` sampai tidak ada yang baru.

### 1.4 ⭐ Pernyataan 8 — 27 benar, dan dua cara menjelaskan kenapa 160 juga benar

`[terverifikasi]` `pyRequired` = `true` muncul **160 kali** di **6 layar**. ⭐ Itu **kemunculan**,
bukan medan.

| Cara | Hasil |
| --- | ---: |
| kemunculan literal `<pyRequired>true</pyRequired>` | **160** di 6 layar |
| ⭐ **A — properti unik PER LAYAR, dijumlahkan** | **80** |
| ⭐⭐ **B — properti unik LINTAS layar** | ⭐ **27** |

✅ **Berkas keadaan memakai cara B, dan 27 tepat.** ⭐ Selisih 80 lawan 27 wajar: medan yang sama
muncul di beberapa layar — misalnya layar `Detail…` dan `General…` yang berpasangan.

⭐ **Ke-27 medan wajib, apa adanya:** `.Claim` · `.ClaimPaymentType` · `.ClaimType` ·
`.DateofSurvey` · `.Deduction1` · `.Deduction2` · `.EndDate` · `.ExcessLoss` · `.IsApproved` ·
`.OutstandingClaim` · `.OveriddingCommOgp` · `.OveriddingCommOnp` · `.PremiOgp` · `.PremiOnp` ·
`.ProductionDate` · `.Quartal` · `.QuotationData.IsSurveyReport` · `.QuotationData.MOID` ·
`.ResultOnp1` · `.RiCommOgp` · `.RiCommOnp` · `.SalvageValue` · `.StartDate` · `.StatementDate` ·
`.Suggest` · `.TypeTax` · `.YearOfQuartal`.

⚠️ **Sebaran per layar:** `DetailPolicyTreatyIn` 22 · `GeneralPolicyTreatyIn` 22 ·
`DetailDeptHeadTreatyIn_UW` 16 · `GeneralDeptHeadTreatyIn_UW` 16 · `ListSuggest` 3 ·
`InputHistoricalSurveyReportDtl` 1.

> ⛔ **RALAT ronde 4.** Ronde 4 Bab A.3 menulis *"medan wajib pada sel `FIELD` = 0 … saya tidak
> menyimpulkan bahwa modul ini tidak punya medan wajib"*.
> ⭐ **Kehati-hatiannya benar, pengukurannya tidak.** ⛔ Saya menambatkan `pyRequired` pada
> `pyUserData` anak sel `FIELD`; ⭐ **tempat sebenarnya lebih tinggi** — ia ditemukan dengan
> menaiki simpul induk sampai bertemu `pyValue` berawalan titik.

### 1.5 ⚠️ Pernyataan 10 — separuh dapat diperiksa, dan separuhnya lagi bukan urusan korpus

| Bagian | Hasil |
| --- | --- |
| ⭐ **"33 dari 33 medan laporan tersedia"** | ✅ `ReportDefinition\BrowseTreatyJoinEDM.xml`, kelas `ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM`, memuat **`pyFieldName` berisi 66, unik 33** |
| ⛔ **"view `POOLDATA.TREATYINDETAILJOINEDM` berkolom 39"** | ⛔ **TIDAK dapat diperiksa dari korpus** |

⚠️ **Lingkup penelusuran, disebut sesuai jebakan #4:** disisir **seluruh 278 berkas** modul untuk
teks `TREATYINDETAILJOINEDM` ⇒ **11 berkas** menyebutnya, ⛔ **nol** di antaranya memuat definisi
view atau daftar kolomnya. ⭐ Ke-41 naskah `pyBrowseSQL` juga disisir ⇒ ⛔ **nol** menyentuh view
itu dengan nama tersebut.
⇒ ⭐ **Angka 39 adalah `[data DBA]`**, dan **tidak dapat dibantah maupun dibenarkan dari sini.**

---

## Bab 2 — Pertentangan antar sumber

⭐ **Urutan kewenangan dipakai apa adanya:** ① lembar jawaban *(mengikat)* → ② berkas keadaan →
③ grilling ronde 1–4 *(tersegel)*.

### 2.1 ⭐ Antara sumber ① dan ② — NOL pertentangan

`[terverifikasi]` Ke-**12 ketetapan** di `KEADAAN` §4 merujuk **13 nomor pertanyaan** —
P2, P3, P4, P6, P23, P24, P25, P29, P32, P33, P45, P46, P47.
✅ **Ketiga belasnya SUDAH TERJAWAB** di lembar pemiliknya.
⛔ **Tidak ada satu pun ketetapan yang bersandar pada pertanyaan yang masih kosong.**

⭐ Itu pemeriksaan yang paling menentukan untuk Bab 5, dan hasilnya bersih.

### 2.2 ⭐ Antara sumber ② dan ③ — daftar koreksi berkas keadaan diperiksa kelengkapannya

`KEADAAN` §5 mendaftar **8 koreksi** atas ronde 1–4. ⭐ **Diuji satu per satu:**

| Koreksi di §5 | Diuji ulang | |
| --- | --- | :---: |
| *"942 Property-Set"* → **938** | **938**, dua cara | ✅ |
| *"17 Page-Clear-Messages"* → **16** | **16**, dua cara | ✅ |
| *"lapisan layar nol dibuka"* → ronde 2 sudah menyisir 31 | ✅ ronde 4 §0.2 juga sudah meralatnya | ✅ |
| *"naskah SQL tidak terekspor"* → ada di `<pyBrowseSQL>` | ⭐ **41 di modul ini** | ✅ |
| *"baris tabel keputusan tidak terkirim"* → ada di `pyCondition`/`pyOrConditions`/`pyResults` | ⭐ **36 `rowdata`** di `pyPropertyValues` dan `pyResults` | ✅ |
| *"BusinessType_DeT tinggal 9 dari 34"* → **36 utuh** | ✅ | ✅ |
| *"160 setelan wajib tidak menempel medan"* → **27 medan di 6 layar** | ✅ | ✅ |
| *"seluruh 38 terkunci di layar Dept Head"* → **36 di sana, 2 biasa** | ✅ **36 · 2** | ✅ |

⭐ **Dua koreksi yang TERLEWAT dari §5, dan saya tambahkan:**

| # | Yang terlewat | Bunyi ronde 4 | Yang benar |
| --- | --- | --- | --- |
| ⭐ **a** | **ukuran modul** | ronde 2 dan 4 memakai **35.678.284 B** dan **46,6 %** untuk `Section`+`Harness` | ⭐ ukuran modul kini **35.492.317**; ⚠️ maka **46,6 % menjadi 46,9 %** — `Section`+`Harness` tetap 16.635.834 B tetapi penyebutnya mengecil |
| ⭐ **b** | **md5 korpus** | ronde 1–4 mengunci invarian pada `62a3735ebb2eabb788c8cd8abdda78ca` | ⭐ kini **`96271809212b68bde10ae8612f414baf`**; ⛔ setiap pemeriksaan ulang ronde 1–4 akan **gagal** bila memakai md5 lama |

⚠️ **Butir b penting bukan karena angkanya**, melainkan karena ⛔ **empat berkas grilling menyatakan
"korpus tidak berubah" dengan md5 yang kini usang.** ⭐ Pernyataan itu **benar pada saat ditulis**
dan **tidak boleh dibaca sebagai berlaku sekarang**.

### 2.3 ⚠️ Satu nomor yang berubah isi antar berkas — dan ia BUKAN pertentangan

⚠️ Pembanding otomatis saya menandai **P42** sebagai berbeda antara `PERTANYAAN-RONDE-4.md` dan
lembar pemiliknya. ⛔ **Diperiksa: itu galat pembanding saya**, bukan perbedaan isi — pemisah blok
saya ikut menyerap judul bab berikutnya *(`# Untuk Finance`)*.
✅ **Isi jawabannya identik.**

---

## Bab 3 — Kelengkapan dan sinkronisasi lembar jawaban

### 3.1 ⭐ Cacahan — DUA CARA

| Lembar | ⭐ A: judul `##` | ⭐ B: slot `**Jawaban:**` | Terjawab | Kosong |
| --- | ---: | ---: | ---: | ---: |
| DBA | **6** | 6 | **5** | ⛔ **1** |
| Finance | 3 | 3 | **3** | 0 |
| IAM | 5 | 5 | **5** | 0 |
| Product & Underwriting | **17** | 17 | **17** | 0 |
| pemilik export Pega | **5** | ⚠️ **4** | **4** | ⛔ **1** |
| pengembang Pega lama | **15** | ⚠️ **14** | **15** | 0 |
| ⭐ **TOTAL** | ⭐ **51** | **49** | ⭐ **49** | ⭐ **2** |

⚠️ **Dua lembar memberi A ≠ B, dan itu BENAR, bukan cacat:** keduanya memuat satu pertanyaan
berjudul **DITARIK** — **P41** di lembar pemilik export, **P48** di lembar pengembang Pega lama.
⭐ Pertanyaan yang ditarik **tidak berbadan lima butir dan tidak berkolom jawaban**, jadi ia
terhitung di cara A tetapi tidak di cara B.
⇒ ⭐ **51 judul − 2 ditarik = 49 yang dapat dijawab**, dan ⭐ **49 terjawab**.

### 3.2 ⛔ Dua yang masih kosong — dan keduanya persis yang berkas keadaan sebut

| # | Lembar | Judulnya |
| --- | --- | --- |
| ⛔ **P1** | DBA | isi tiga program penyimpan data di basis data |
| ⛔ **P18** | pemilik export Pega | isi **268** langkah penghitungan |

✅ **Cocok tepat dengan `KEADAAN` §6**, yang menyebut hanya P1 dan P18 sebagai penahan.
⭐ Perhatikan judul P18 kini berbunyi **268**, bukan 942 — ⭐ sejalan dengan `KEADAAN` §5
*("938 di folder; **268** yang diminta sesudah penyaringan")*.

### 3.3 ⛔ Sinkronisasi salinan — satu salinan sudah USANG

| Berkas | Sifat | Keadaan |
| --- | --- | --- |
| enam lembar `PERTANYAAN-untuk-*.md` | ⭐ **sumber** | ✅ mutakhir |
| `PERTANYAAN-RONDE-4.md` | salinan | ✅ **selaras** — 8 pertanyaan P42–P49, isinya sama dengan lembar pemilik |
| ⛔ **`PERTANYAAN-YANG-MASIH-KOSONG.md`** | salinan | ⛔⛔ **USANG** |

⛔ **Berkas itu dibuat sebagai cuplikan pada pukul 14.45** dan memuat **19 pertanyaan** yang saat
itu kosong. ⭐ Sejak itu **17 di antaranya sudah terjawab** di lembar pemiliknya.
⇒ ⛔ **Siapa pun yang memakainya sekarang akan mengerjakan ulang 17 pertanyaan yang sudah selesai.**

⭐ **Saran, dan ia keputusan Anda bukan saya:** berkas itu **dibuat ulang** dari lembar pemilik,
atau **dihapus**. ⛔ Saya tidak melakukan keduanya — brief ronde ini hanya mengizinkan satu berkas
keluaran.

---

## Bab 4 — Kandidat berkas korpus yang mencurigakan

⭐ **NOL kandidat ditemukan**, dengan cara yang **berbeda** dari cara penyusun.

⭐ **Cara yang dijalankan** — bukan mencari kelas berbagi, melainkan **membandingkan isi**:

1. md5 dihitung untuk **ke-278 berkas** NB Treaty In;
2. untuk tiap berkas, dicari berkas **bernama sama pada jalur relatif sama** di **20 modul lain**
   *(⛔ `nusantara-re` dikecualikan)*;
3. berkas NB yang kelasnya **BUKAN** Treaty disaring — **205 berkas**;
4. dari 205 itu, ditandai yang **modul Treaty lain punya varian berkelas Treaty dengan md5
   berbeda** — ⭐ itulah pola yang menjerat `Protection_Act`.

⇒ ⛔ **Hasil: 0 kandidat.**

⚠️ **Lingkup dinyatakan sesuai jebakan #4.** ⛔ Cara ini **tidak** menangkap: berkas salah yang
**tidak punya kembaran bernama sama** di modul lain; berkas yang **salah versi tetapi berkelas
benar**; dan berkas yang salah di modul **selain** NB Treaty In.
⭐ Jadi *"nol"* berarti **nol menurut cara ini**, bukan nol mutlak.

⭐ **Dan satu bukti tambahan bahwa perbaikan `Protection_Act` memang beres:** berkasnya kini
**157.304 B**, berkelas **`ASM-FW-GISFW-Data-PolicyTreatyIn`**, dan **selisih ukuran modul cocok
sampai satuan byte** *(Bab 1.1)*.

---

## Bab 5 — Putusan

> ⭐⭐ **`KEADAAN-NB-TREATY-IN.md` LAYAK dipakai sebagai dasar penulisan spec.**

⭐ **Sembilan dari sepuluh pernyataan cocok tepat**, empat di antaranya diuji **dua cara yang
benar-benar berbeda**; pernyataan kesepuluh cocok pada bagian yang **dapat** diperiksa dari korpus,
dan bagian sisanya **bukan urusan korpus**. ⛔ **Nol pertentangan** antara lembar jawaban dan
berkas keadaan, dan **ke-13 nomor pertanyaan yang menopang 12 ketetapannya sudah terjawab**.

⚠️ **Tiga hal yang sebaiknya dibereskan lebih dulu — ⛔ tidak satu pun menahan penulisan spec:**

| # | Yang perlu dibereskan | Sifat |
| --- | --- | --- |
| ⛔ **1** | **`PERTANYAAN-YANG-MASIH-KOSONG.md` usang** — 17 dari 19 isinya sudah terjawab | ⚠️ **berisiko**: menyesatkan penjawab, bukan menyesatkan spec |
| ⭐ **2** | **§5 berkas keadaan tambahkan dua koreksi** — ukuran modul *(46,6 % → 46,9 %)* dan **md5 korpus baru** | ⭐ kerapian jejak |
| ⚠️ **3** | **Angka 39 kolom view** ditandai `[data DBA]`, bukan terukur | ⭐ kejujuran label |

⛔ **Dan yang tetap menahan:** **P1** saja. ⭐ Spec untuk **alur, tahap, wewenang, kronologi,
masa berlaku, layar, dan — sejak 2026-09-22 — perhitungan uang** dapat ditulis sekarang;
⛔ spec **penyimpanan** tidak.

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat ini semula berbunyi:
> > *"⛔ **Dan yang tetap menahan, tidak berubah:** **P1** dan **P18**. ⭐ Spec untuk **alur, tahap,
> > wewenang, kronologi, masa berlaku, dan layar** dapat ditulis sekarang; ⛔ spec **penyimpanan**
> > dan **perhitungan uang** tidak."*
>
> **P18 ditarik.** Perhitungan uang tidak lagi tertahan olehnya.

---

## Bab 6 — TELEMETRI EKSEKUSI

⛔⛔ **Pengukuran dari luar TIDAK dilakukan.** ⭐ Brief meminta `claude --print --output-format json`
dijalankan dari luar sesi; ronde ini berjalan **di dalam sesi interaktif**, sehingga
**`total_cost_usd` dan `duration_ms` tidak tersedia**. ⛔ **Tidak saya taksir.**

| Yang dicatat | Nilai | Sumber |
| --- | ---: | --- |
| ⭐ **token keluaran** | ⭐ **92.996** | transkrip `.jsonl`, SELISIH terhadap baseline |
| ⭐ **token cache-read** | ⭐ **20.800.918** | sama |
| token cache-write | **2.717.760** | sama |
| token masuk | **70** | sama |
| ⭐ **jumlah panggilan alat** | ⭐ **35** | sama |
| ⛔ durasi | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar |
| ⛔ biaya | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar |
| ⭐ **byte dibaca** | ⭐ **35.492.317** *(278 berkas modul, diurai penuh)* ditambah **20 modul lain** untuk Bab 4 | perintah audit di tiap bab |

⚠️ **Tiga keterbatasan pengukuran dari dalam:** *(a)* pesan terakhir belum tertulis ke transkrip
saat pengukuran kedua — angkanya **kurang satu pesan**; *(b)* jalur transkrip **diperiksa** dan
cocok dengan sesi ini; *(c)* yang dilaporkan **SELISIH**, bukan total sesi.

---

## Lampiran — bukti isolasi

| Yang dibuktikan | Keadaan |
| --- | --- |
| ⛔ `D:\XML\nusantara-re\` | ✅ **NOL disentuh** — disaring keluar dari daftar modul |
| korpus `D:\XML\RNM_BRD\` | ✅ **READ-ONLY** — hanya dibaca |
| ⛔ `KEADAAN-NB-TREATY-IN.md` disunting | ✅ **NOL** |
| ⛔ `grilling-ronde-1..4.md` disunting | ✅ **NOL** |
| ⛔ lembar jawaban disunting | ✅ **NOL** — hanya dibaca dan dicacah |
| berkas dibuat | ✅ **tepat SATU** — berkas ini |
| pertanyaan terbuka ditutup | ✅ **NOL** |
| nilai nama orang disalin | ✅ **NOL** |
| rahasia · data nasabah · nomor polis apa adanya | ✅ **NOL** |

⭐ **Penanda dipakai:** `[terverifikasi]` untuk yang berperintah audit · `[data DBA]` untuk angka
39 kolom view. ⛔ Nol `[dugaan]` dipakai sebagai dasar putusan.

⭐ **Bandingkan dengan lima ronde sebelumnya** — seluruhnya SELISIH, alat yang sama:

| | R1 | R2 | R3 | R4 | ⭐ **Verifikasi** |
| --- | ---: | ---: | ---: | ---: | ---: |
| keluaran | 192.021 | 225.898 | 178.857 | 154.518 | ⭐ **92.996** |
| cache-read | 13.532.360 | 22.688.258 | 16.725.320 | 26.189.892 | **20.800.918** |
| panggilan | 69 | 72 | 39 | 47 | ⭐ **35** |

⭐ **Keluaran dan panggilan terendah dari seluruh ronde** — ⭐ sebab ronde ini **menguji sepuluh
pernyataan yang sudah dirumuskan**, bukan menjelajah korpus dari nol. ⚠️ Cache-read tetap besar
karena ke-278 berkas modul diurai penuh, ditambah 20 modul lain untuk Bab 4.

