# Cocok-silang cacah atribut §10 — dan apa yang sapuannya temukan

**Tanggal:** 24 September 2026 · **Langkah 1 putaran penutup to-spec**
**Dibangkitkan:** `alat/periksa-cacah-atribut.py` dan `alat/cocok-skalar-tanpa-rumah.py`.
Naskahnya ditulis tangan; **angkanya dari perkakas**, dan dapat dijalankan ulang.

> ### KENAPA DISAPU, BUKAN DITAMBAL SATU PER SATU
>
> `F-2` (§10.2 mengaku 47/48 memuat 44) dan `F-3` (§10.21 mengaku 6 memuat 5) ditemukan satu per
> satu. **Tiga instans di satu berkas bukan kebetulan** — jadi polanya disapu secara mekanis, bukan
> dicari dengan mata.
>
> Hasilnya: **dua instans lagi** yang `F-2` dan `F-3` tidak lihat, dan **keduanya ternyata bukan
> lubang**. Keduanya tetap menghasilkan temuan, dan justru itu gunanya menyapu.

---

## 1. Tabel cocok-silang — SELURUH §10.x, termasuk yang cocok

Dua cacah dicetak terpisah, dan itu keputusan: **BARIS** menjawab *"berapa baris ditulis"*, **NAMA**
menjawab *"berapa kolom akan berdiri di DDL"*. Keduanya berbeda ketika satu baris memuat lebih dari
satu nama.

| § | Entitas | Diumumkan | Baris | Nama | Hasil |
|---|---|---:|---:|---:|---|
| 10.1 | `KONTRAK` | 8 | 8 | 8 | cocok |
| **10.2** | **`VERSI_KONTRAK`** | ~~47~~ **44** | **44** | **44** | **COCOK sejak 24 Sep 2026** — `F-2` ditutup; angka judulnya yang diturunkan, sesudah pencariannya diselesaikan (§4) |
| **10.3** | `LAYER` | 19 | 25 | 25 | **TIDAK COCOK — lihat §2.1, bukan lubang** |
| 10.4 | `DETAIL_PROPORSIONAL` | 12 | 12 | 12 | cocok |
| 10.5 | `BAGIAN` | 8 | 8 | 8 | cocok |
| 10.6 | `PENYEBARAN` | 5 | 5 | 5 | cocok |
| 10.7 | `RINCIAN_PENYEBARAN` | 4 | 4 | 4 | cocok |
| 10.8 | `NILAI_PENYEBARAN` | — | 3 | 3 | judul tanpa cacah |
| **10.9** | **`POTONGAN`** | 5 | **0** | **0** | **tabelnya di §14.2 — lihat §2.2, bukan lubang** |
| 10.10 | `MATA_UANG_KONTRAK` | 6 | 6 | 6 | cocok |
| 10.11 | `RETENSI_CEDANT` | 6 | 6 | 6 | cocok |
| 10.12 | `EGNPI` | 8 | 8 | 8 | cocok |
| 10.13 | `PORTOFOLIO` | 5 | 5 | 5 | cocok |
| 10.14 | `PERIODE_PELAPORAN` | 7 | 7 | 7 | cocok |
| 10.15 | `PERIODE_AKUMULASI` | 6 | 6 | 6 | cocok |
| 10.16 | `TERMIN` | 8 | 8 | 8 | cocok |
| 10.17 | `SKALA_KOASURANSI` | 4 | 4 | 4 | cocok |
| 10.18 | `BATAS_PER_BAHAYA` | 5 | 5 | 5 | cocok |
| 10.19 | `DOKUMEN_KONTRAK` | 5 | 5 | 5 | cocok |
| 10.20 | `CATATAN_PERSETUJUAN` | 6 | 6 | 6 | cocok |
| **10.21** | **`JEJAK_PERUBAHAN`** | ~~6~~ **8** | **8** | **8** | **COCOK sejak 24 Sep 2026** — barisnya dipecah dua, dan `PERAN_PELAKU` ditambahkan atas ADR-0045 |
| 10.21a | `PERISTIWA_KONTRAK` | 5 | 5 | 5 | cocok |
| 10.22 | enam tabel acuan | — | 5 | 5 | judul tanpa cacah; bentuk seragam × 6 |

**17 cocok · 3 tidak cocok · 2 tanpa cacah di judul.** — *keadaan 24 September 2026 pagi.*
**Sesudah putaran penutup: 19 cocok · 1 tidak cocok (§10.3, dan sebabnya bukan lubang — §2.1) · 2
tanpa cacah di judul.**

**Yang tidak masuk pemeriksaan, dan atas dasar apa** — dicetak perkakasnya sendiri:

| Sebab | Jumlah | Bagian |
|---|---:|---|
| tanpa tabel atribut | 3 | §10.0, §10.9, §10.23 |
| judul tidak menyebut "N atribut" | 2 | §10.8, §10.22 |

> **Batas perkakas, dan ia tertulis di dalam keluarannya sendiri:** ia **tidak dapat menyatakan
> sebuah cacah BENAR.** Ia hanya menyatakan dua pernyataan di dalam berkas yang sama **cocok** atau
> **tidak cocok**. Bila judul dan tabel sama-sama salah, ia **diam** — dan diamnya bukan pernyataan
> bahwa keduanya benar.

---

## 2. Dua yang ternyata BUKAN lubang — dan temuan yang tetap lahir darinya

### 2.1 §10.3 — dua entitas di dalam satu bagian

19 + 6 = 25. Tabel pertama `LAYER` (**19 baris**, cocok dengan judulnya); tabel kedua
**`PEMULIHAN_LIMIT`** (**6 baris**) — entitas tersendiri yang lahir dari §10.23c keputusan 1, punya
berkas DDL sendiri (`31_PEMULIHAN_LIMIT.sql`), dan **tinggal menumpang di dalam §10.3**.

**Temuan yang tetap lahir — `F-11`:**

> **`PEMULIHAN_LIMIT` tidak punya bagian §10.x-nya sendiri.**
>
> Akibatnya tiga, dan ketiganya praktis: siapa pun yang **mencacah entitas dari judul §10**
> kehilangan satu; siapa pun yang membaca §10.3 mengira `LAYER` punya 25 atribut; dan pengurai yang
> membangkitkan kamus kolom **harus tahu memecahnya**, padahal tidak ada yang menyuruhnya.

Ia entitas yang **paling terakhir ditambahkan** dan yang **paling sering tertinggal** — ia juga yang
hilang dari `STRUKTUR-DATA.md` (`F-4`).

**Usulan:** `PEMULIHAN_LIMIT` diberi bagian **§10.3a**, sejajar dengan cara §10.2a dan §10.21a
dibuat. Ini **usulan suntingan berkas induk** dan menunggu persetujuan.

### 2.2 §10.9 `POTONGAN` — nol tabel, dan itu BENAR

Bagiannya berbunyi: *"Atributnya sudah ditetapkan §14.2 dan tidak diulang di sini — mengulangnya
membuat dua tempat, satu akan basi."*

§14.2 memuat **tepat 5 baris** — `ID_POTONGAN`, `ID_INDUK_POTONGAN`, `ID_JENIS_POTONGAN`,
`DASAR_PERHITUNGAN`, `PERSEN_POTONGAN` — **cocok dengan judul §10.9.**

Keputusan "satu fakta satu tempat" itu **benar dan dipertahankan**. Yang kurang adalah bahwa
**perkakas tidak dapat mengikuti penunjuknya**, dan itu batas perkakas, bukan cacat berkas.

> **Nilai baris ini justru pada penolakannya.** Kalau perkakasnya diam-diam membuang §10.9 tanpa
> melaporkannya, tidak akan ada yang memeriksanya — dan "§10.9 tanpa tabel" akan tetap tidak
> diketahui apakah ia keputusan atau kelalaian.

### 2.3 §10.21 — tambalannya benar, tetapi memunculkan ambiguitas baru

`F-3` sudah ditambal: baris keenam **ada**, berbunyi `NILAI_SEBELUM` / `NILAI_SESUDAH`.

**Satu baris, dua nama.** Judul berbunyi 6; DDL akan berdiri dengan **7 kolom**.

Bukan kesalahan — baris itu sengaja menggabungkan pasangan yang bentuknya sama. Tetapi **cacah yang
dipakai orang berikutnya bergantung pada cara ia menghitung**, dan tidak ada yang menyuruhnya
memilih yang mana.

**Usulan:** judul §10.21 berbunyi **"6 baris · 7 kolom"**, atau barisnya dipecah dua. Menunggu
persetujuan.

---

## 3. §10.2 — satu-satunya lubang yang nyata, dan pencarian keempatnya

### 3.1 Dilokalisasi per kelompok

| Kelompok | Baris |
|---|---:|
| Identitas dan keadaan | 7 |
| Penamaan dan lingkup | 6 |
| Bagian NuRe dan biaya | 6 |
| Ketentuan tertulis | 6 |
| Cara pembukuan dan pelaporan | 11 |
| Termin dan prorata | 2 |
| Kapasitas dan batas | 3 |
| Retro dan status master | 2 |
| **§10.2a** `SIFAT_MATERIAL_ADDENDUM` | 1 |
| **Jumlah** | **44** |

**Tidak ada satu kelompok yang "kurang".** Kekurangannya tersebar atau berada di luar seluruh
kelompok — yang berarti ia **tidak dapat ditemukan dengan membaca §10.2 saja**.

### 3.2 Sumber yang dipakai, dan kenapa ia harus MANDIRI dari §3.2

Blok koreksi §10.2 menyatakan kelima belas jalur §3.2 yang menganggur **sudah diadili satu per satu
dan tidak satu pun menjelaskan keempatnya**. Maka §3.2 sudah habis.

Dan `F-8` membuktikan semesta §3.2 memang **bolong**: `TreatyIn.BrokeragePct` **tidak ada di
dalamnya sama sekali**.

Sumber mandiri yang dipakai: **`datar-treatyin-lama.csv`** — dibangkitkan `buat-pohon-treatyin.py`
dari sapuan 329 berkas ekspor, **bukan** dari §3.2 dan **bukan** dari §10.

**Yang dibandingkan:** 99 skalar kedalaman-1 `TreatyIn` terhadap 131 nama *Asal* yang terkumpul dari
seluruh tabel §10.

**Hasil: 42 skalar tidak muncul sebagai *Asal* mana pun.** Penolakan perkakasnya: 43 simpul bukan
skalar (Page / Page List), dan empat pohon cermin dikeluarkan bernama.

### 3.3 Empat puluh dua calon, digolongkan — dan sebagian besar sudah punya sebab

| Golongan | Jumlah | Contoh | Sebab yang sudah tertulis |
|---|---:|---|---|
| agregat `Total*` | 14 | `TotalShareNet`, `TotalLimitsROL` | turunan, tidak disimpan (ADR-0037) |
| denominasi `Currency*` | 6 | `CurrencyEarthquake`, `CurrencyRSMD` | **melebur** ke paket uangnya (§10.0a) |
| pro rata | 3 | `ProRatePercent`, `ProRateDays` | turunan; mesinnya sengaja tidak dibangun (`GRL-15`) |
| penanda mekanisme lama | 6 | `ViewState`, `IsEditData`, `RevisionState`, `PositionUsername`, `ChooseStatusAkseptasi`, `AddendumPremi` | **tidak dibawa** (`GRL-08`, `GRL-18`) |
| sudah diadili dan dibuang | 5 | `CedingStatusActive`, `SourceStatusActive`, `ContractRefNo`, `LeadingReinsName`, `BrokeragePercentP` | blok koreksi §10.2 |
| nol penulis, kolom selalu kosong | 2 | `BrokeragePct`, `NusareSharePct` | **`F-8`** |
| perkakas Pega | 1 | `pyErrMsg` | bukan fakta bisnis |
| **BELUM punya sebab tertulis** | **5** | **`TreatyYear` · `RevisionDate` · `CoInScale` · `Ceding` · `LeadingReinsSource`** | — |

### 3.4 Lima yang belum punya sebab — **BELUM diputuskan, dan tidak dikarang**

> **Ini calon, bukan putusan.** Perkakasnya menyatakan namanya **tidak muncul** di kolom *Asal* mana
> pun — **bukan** bahwa ruasnya tidak punya rumah. Ia mungkin berumah dengan nama *Asal* yang
> ditulis berbeda.

Yang paling kuat di antaranya, dan sebabnya disebut supaya yang meneruskan tahu ke mana melihat:

| Calon | Kenapa ia kuat |
|---|---|
| **`TreatyYear`** | ia **salah satu dari 25 parameter** yang diserahkan ke prosedur `PEGA_M_TREATY_IN_EDM`, jadi ia **ada di tabel relasional lama** — bukan hanya di clipboard |
| **`RevisionDate`** | ditulis `TreatyInRevisi_post`, jalur khas addendum |
| **`CoInScale`** | skalar kedalaman-1 yang **bernama sama** dengan daftar `CoInScale[]`; kemungkinan besar bukan atribut versi melainkan penampung layar — **perlu dibaca penulisnya** |
| **`Ceding`, `LeadingReinsSource`** | keduanya **nama**, sementara §10.2 kemungkinan hanya membawa **ID**-nya. Bila benar, keduanya sengaja tidak disimpan — tetapi **itu belum tertulis di mana pun** |

**Siapa menutup:** sesi to-spec, dengan membaca penulis kelimanya lewat `alat/sapu-penulis-properti.py`
dan memeriksa apakah masing-masing muncul di DDL `TREATY_IN` lama.

**Yang menagih:** `KAMUS-KOLOM.md` §0 — tabel selisihnya sudah berdiri di muka, dan ia akan terus
berbunyi **"empat diumumkan, tidak pernah didaftar"** sampai kelima calon ini diadili.

> **Tidak satu pun dari kelima ditambahkan ke §10.2 dalam putaran ini.** Mengarang atribut membuat
> kamus kolom berbohong dengan cara yang paling sulit ditemukan — pembacanya tidak punya cara
> membedakan mana yang berasal dari §10 dan mana yang ditambahkan.

---

## 4. Kelima calon §3.4 — **DIADILI 24 September 2026**

**Sumber:** `alat/datar-penulis-properti.csv` — keluaran `alat/sapu-penulis-properti.py`, kelima
bentuk penulis, kedua ekspor, 708 berkas. Tidak ada berkas XML yang dibaca langsung.

> **Yang dicari bukan "apakah namanya ada", melainkan SIAPA MENULISNYA dan DARI APA.** Sebuah nama
> yang tidak muncul di kolom *Asal* mana pun dapat berarti tiga hal yang sangat berbeda: ia atribut
> yang hilang, ia turunan, atau ia bukan atribut sama sekali. Hanya penulisnya yang memisahkan
> ketiganya.

| Calon | Putusan | Bukti, dikutip apa adanya |
|---|---|---|
| **`TreatyYear`** | **TURUNAN — tidak disimpan** (ADR-0037). **Bacaan pokok, dengan sisa yang disebut** — lihat §4.3 | `DataTransform/TreatyInSetTreatyYear.xml`: `TreatyIn.TreatyYear = @substring(TreatyIn.Commencement,0,4)` — **empat aksara pertama tanggal mulai**, yang sudah tersimpan sebagai `KONTRAK.TANGGAL_MULAI` |
| **`RevisionDate`** | **DIBUANG — jejak mesin** | satu-satunya penulisnya `Activity/TreatyInEDMSetValue.xml`: `TreatyIn.RevisionDate = @CurrentDateTime()`. Sekeluarga dengan `EDMDate` (ADR-0006) — tanggal **sentuh**, bukan tanggal fakta. Pohon mencatatnya `hanya-adjustment`, 0 rujukan, 0 berkas |
| **`CoInScale`** | **BUKAN ATRIBUT — penampung layar** | **nol penulis** atas kelima bentuk. Pohon: `Skalar`, kelas Pega **`(tidak dideklarasikan)`**, muncul hanya di `Harness/InputTreatyInOffer.xml`, `Section/InputTreatyInOffer.xml`, `Section/TreatyInNONProportional.xml`. Skala ko-asuransi yang sebenarnya adalah daftar `CoInScale[]` → §10.17 `SKALA_KOASURANSI` |
| **`Ceding`** | **SUDAH punya sebab tertulis — di §10.1, bukan §10.2** | §10.1: *"**`Ceding` dan `LeadingReinsSource` (nama) tidak disimpan.** Keduanya salinan dari master"*. Yang disimpan `ID_CEDANT` ← `CedingID`. Penulisnya menguatkan: `SaveData.CEDING = TreatyIn.Ceding` **berdampingan** dengan `SaveData.CEDINGID = TreatyIn.CedingID` |
| **`LeadingReinsSource`** | idem | §10.1 `ID_ASAL_BISNIS` ← `LeadingReinsSourceID`, berketerangan *"source of business"*. Penulisnya: `SaveData.SOB` / `SaveData.SOBID` — **nama kolom datarnya sendiri berbunyi SOB**, bukan "leading reins" |

**Lima diadili — NOL menjadi atribut `VERSI_KONTRAK`.**

### 4.1 Satu koreksi terhadap §3.4, dinyatakan bukan diselaraskan diam-diam

§3.4 menulis tentang `Ceding` dan `LeadingReinsSource`: *"Bila benar, keduanya sengaja tidak
disimpan — tetapi **itu belum tertulis di mana pun**."*

**Itu keliru. Ia tertulis**, di §10.1, dalam satu kalimat yang menyebut keduanya bersama. Yang
terjadi bukan sebab yang hilang melainkan **sebab yang dicari di bagian yang salah**: keduanya
atribut `KONTRAK`, bukan `VERSI_KONTRAK`, dan sapuannya membandingkan skalar terhadap *seluruh*
§10 tanpa memisahkan entitasnya.

> **Pelajaran yang dibawa ke sapuan berikutnya:** perkakas yang mencocokkan nama terhadap semesta
> gabungan akan melaporkan "tidak punya rumah" untuk sesuatu yang berumah **di tetangga**. Yang
> dicetaknya bukan lubang, melainkan **daftar tempat mencari**.

### 4.2 `LeadingReinsSource` — nama yang berarti lain daripada bunyinya

| | |
|---|---|
| **Bunyinya** | *sumber reasuradur pemimpin* |
| **Isinya** | **source of business** — `TreatyInMappingDataconvert.xml`: `TreatyIn.LeadingReinsSource = .SobName`, dan kolom datarnya `SOB` / `SOBID` |
| **Kenapa ini dicatat** | §10.2 punya atribut **lain** bernama `ID_REASURADUR_PEMIMPIN` ← `LeadingReinsID`. Dua nama yang hampir sama untuk dua benda yang tidak berhubungan. `ID_ASAL_BISNIS` di §10.1 sudah memakai nama yang menyebut isinya — **keputusan itu benar dan dikuatkan di sini** |

> Aturan yang berlaku: **pengenal sistem lama dikutip apa adanya dan tidak dirapikan.** Maka
> `LeadingReinsSource` tetap ditulis begitu di kolom *Asal*; yang tidak dibawa adalah **artinya yang
> menyesatkan**, bukan ejaannya.

### 4.3 `TreatyYear` — satu sisa yang TIDAK ditutup, dan kenapa `T-7` tetap dikirim

Penulis `TreatyIn.TreatyYear` ada **dua**, bukan satu:

| Penulis | Yang ditulisnya |
|---|---|
| `DataTransform/TreatyInSetTreatyYear.xml` | `@substring(TreatyIn.Commencement,0,4)` — **diturunkan** |
| `Activity/TreatyInMappingDataconvert.xml` | `TreatyOffer.TreatyYear` — **disalin dari borang penawaran** |

**Yang tidak terbaca dari ekspor: urutannya.** Bila penawaran memuat tahun yang **berbeda** dari
tahun tanggal mulai, ekspor tidak memberi tahu apakah nilai itu bertahan atau ditimpa oleh
*data transform* di atas.

> **Ini bentuk "uji yang dijawab sama oleh kedua kemungkinan".** Bila seluruh kontrak yang ada
> kebetulan bertahun-treaty sama dengan tahun mulainya, maka **"selalu diturunkan"** dan
> **"kadang disepakati, tetapi selama ini kebetulan sama"** menghasilkan data yang **persis sama**.
> Membaca data saja tidak memisahkannya.

Maka **`T-7` tetap dikirim** — tetapi **dipersempit**, dan bentuknya diubah: ia tidak lagi meminta
pernyataan, ia meminta **satu contoh alur kerja**. Pertanyaan bentuk pertama menghasilkan ingatan
yang membenarkan diri; bentuk kedua menghasilkan perilaku yang dapat diperiksa.

**Akibat bila bacaan pokoknya salah:** `TAHUN_TREATY` menjadi atribut `KONTRAK`, satu kolom
tambahan. **Murah** — dan itu sebab menutup `F-2` tidak perlu menunggu jawaban `T-7`.


---

## 5. Yang belum dikerjakan di langkah 1


| # | Sisa | Keadaan |
|---|---|---|
| 1 | kelima calon §3.4 diadili | **SELESAI 24 Sep 2026** — §4. Nol menjadi atribut; `F-2` ditutup di §10.2 |
| 2 | `DEDUCTIBLE2`, `PREMIUM_EARNED`, `ROL_PCT` — tiga kolom tanpa rumah dari modul Adjustment | **satu mendarat, dua dipertikaikan** — `2-to-spec/F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` |
| 3 | `KAMUS-KOLOM.md` dan `ddl-usulan/` dibangkitkan ulang | **SELESAI** — selisihnya dicetak di `SISA-DAN-SERAH-TERIMA.md` §12 |

Butir 3 sengaja **tidak** dijalankan: membangkitkan ulang tanpa perubahan di §10 hanya menghasilkan
keluaran yang sama, dan selisih nol akan terbaca sebagai "sudah lengkap".
