# Struktur data — Treaty In dan Treaty In Adjustment

**Tanggal:** 23 September 2026
**Skema:** `TREATY_MASUK` · akun aplikasi `TREATY_MASUK_APP`
**Mendahului:** `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` — G1, G1b, G1c, G2, G3, G4 seluruhnya tertutup

> **Isi berkas ini: nama, arti satu kalimat, kunci utama, kunci alami, induk.**
> Atribut selengkapnya ada di `SPEC-MODEL-DATA.md` §10 dan **tidak** diulang di sini.

---

## 0. Sekat antarmodul, dan cara membacanya di seluruh berkas ini

| Tanda | Artinya |
|---|---|
| **TI** | milik Treaty In |
| **ADJ** | milik Treaty In Adjustment |
| **B** | dipakai **bersama** — mengikat kedua sesi, tidak boleh diubah sepihak |
| **[luar]** | di luar skema `TREATY_MASUK` |
| **[G2]** / **[G3]** | gelombang berikutnya — dimodelkan, tidak dibangun sekarang |

**Sekatnya tipis, dan itu hasil G1.** Adjustment tidak membawa entitas kontrak sendiri: sebuah
penyesuaian **adalah** sebuah `VERSI_KONTRAK`. Yang benar-benar milik Adjustment hanya catatan
selisihnya.

---

## 1. Entitas Treaty In

### 1.1 Tulang punggung

| Entitas | Bita | Arti | Kunci utama | Kunci alami | Induk |
|---|---|---|---|---|---|
| `KONTRAK` | 7 | satu kontrak treaty masuk, pada lapisan yang tidak pernah berubah sepanjang hidupnya | `ID_KONTRAK` | cedant + asal bisnis + periode + sifat proporsi — **memperingatkan, tidak melarang** | — |
| `VERSI_KONTRAK` | 13 | keadaan lengkap sebuah kontrak pada satu titik, dengan persetujuannya sendiri | `ID_VERSI_KONTRAK` | `NOMOR_URUT_VERSI` di dalam kontraknya | `KONTRAK` |

**`VERSI_KONTRAK` memikul dua peran sekaligus**, dan itu keputusan G1: ia versi biasa **dan**
penyesuaian. Versi yang lahir dari penyesuaian dikenali dari terisinya `ID_VERSI_KONTRAK_DASAR`.

> **`ID_VERSI_KONTRAK_DASAR` adalah satu-satunya penunjuk yang disimpan, dan ia tidak punya
> pasangan.** Sisi "baru" dari sebuah selisih **adalah versi induk barisnya sendiri** — tidak perlu
> kolom untuk itu, dan pembaca berikutnya tidak perlu mencarinya.
>
> Ia disimpan **eksplisit**, bukan disimpulkan dari `NOMOR_URUT_VERSI - 1`, karena *picker* sistem
> lama meng-UNION kontrak dengan addendum sehingga dasar sebuah penyesuaian **belum tentu** versi
> tepat sebelumnya (G2, G4).

### 1.2 Anak langsung versi — dipakai kedua cabang

| Entitas | Bita | Arti | Kunci alami di dalam versinya |
|---|---|---|---|
| `MATA_UANG_KONTRAK` | 17 | satu mata uang yang berlaku pada kontrak, beserta kurs dan periode berlakunya | kode mata uang |
| `RETENSI_CEDANT` | 14 | retensi maksimum cedant per kelompok treaty | kelompok treaty + mata uang |
| `EGNPI` | 5 | perkiraan *gross net premium income* per kelompok treaty | kelompok treaty + mata uang |
| `PORTOFOLIO` | 10 | portofolio masuk atau keluar yang menyertai kontrak | jenis + jenis portofolio |
| `PERIODE_PELAPORAN` | 17 | satu periode pelaporan beserta tanggal jatuh temponya | periode |
| `PERIODE_AKUMULASI` | 17 | satu periode akumulasi | periode |
| `TERMIN` | 6 | satu termin pembayaran premi | nomor termin |
| `SKALA_KOASURANSI` | 16 | satu baris skala ko-asuransi: batas persen dan bagiannya | persen limit |
| `BATAS_PER_BAHAYA` | 16 | batas tanggungan untuk satu bahaya bernama | bahaya |
| `CATATAN_PERSETUJUAN` | 18 | satu keputusan persetujuan: siapa, kapan, dan apa alasannya | (tidak ada — urutan waktu) |
| `DOKUMEN_KONTRAK` | 15 | rujukan ke satu dokumen lampiran — **dirujuk, tidak dimiliki** | pengenal dokumen |
| `JEJAK_PERUBAHAN` | 15 | satu perubahan atas fakta kontrak: siapa, kapan, apa | (tidak ada — urutan waktu) |

### 1.3 Cabang — satu-satunya tempat kedua sisi berpisah

> ### DITAMBAHKAN 24 September 2026 — dua entitas yang lahir sesudah berkas ini ditulis
>
> Berkas ini bertanggal **23 September** dan **mengikat** soal daftar entitas
> (`5-tiket/KEPUTUSAN-PEMBAGIAN-TIKET.md` §0.2 butir 4). Dua entitas diterima **24 September** dan
> tidak pernah menyusul masuk — sehingga wewenang tertinggi menyebut **27** sementara
> `SPEC-MODEL-DATA.md` §10.23c memutuskan **28**, dan `2-to-spec/ddl-usulan/` berdiri dengan **29**
> tabel.
>
> | Entitas | Bita | Arti | Kunci alami | Induk |
> |---|---|---|---|---|
> | **`PEMULIHAN_LIMIT`** | 15 | satu ketentuan pemulihan limit pada sebuah layer: berapa porsi limit yang dipulihkan dan dengan tarif premi berapa | **belum bernomor** — `SPEC-INVARIAN.md` §7.3 | `LAYER` |
> | **`PERISTIWA_KONTRAK`** | 17 | satu peristiwa yang dialami sebuah versi kontrak, dicatat apa adanya | **tidak ada, dan itu keputusan** — peristiwa yang sama dapat terjadi dua kali pada versi yang sama (§10.21a) | `VERSI_KONTRAK` |
>
> `PEMULIHAN_LIMIT` lahir dari §10.23c butir 1; kemampuan *"syarat berbeda tiap pemulihan"*
> bergolongan **BARU** (P-59), sebab sistem lama menulis kedua persentasenya **tetap `"100"` di
> dalam kode**. `PERISTIWA_KONTRAK` lahir dari §10.21a.
>
> **Kunci alami keduanya tidak saya karang** — keadaannya diambil apa adanya dari
> `SPEC-INVARIAN.md` §7.3 dan §7.4.

> ### DITAMBAHKAN 24 September 2026 — `DOKUMEN_ADDENDUM` · **diff `D-2` / `D-3`**
>
> Entitas **BARU**, lahir dari **`GRL-19`** (ronde D grilling Adjustment). `GRL-01` **dikoreksi,
> bukan dibongkar**: pernyataan *"addendum adalah `VERSI_KONTRAK`"* tetap benar — dokumen bukan
> addendum, ia **pembungkus** beberapa addendum.
>
> | Entitas | Bita | Arti | Kunci alami | Induk |
> |---|---:|---|---|---|
> | **`DOKUMEN_ADDENDUM`** | 17 | satu dokumen addendum yang disepakati cedant, memayungi satu atau beberapa versi kontrak | **nomor dokumen**, unik **global** (`INV-71`) | **— berdiri sendiri** |
>
> **Relasi:** dokumen **(1) → (N)** versi, **lintas kontrak**. `VERSI_KONTRAK` memperoleh
> `ID_DOKUMEN_ADDENDUM` yang **boleh kosong**. **Persetujuan tetap melekat pada versi**, per kontrak
> satu per satu — jawaban bisnis *"1 kontrak di aksep 1 per 1"*.
>
> **Label: BARU.** Sapuan nomor dokumen atas 28 nama calon, lima bentuk penulis, 708 berkas, kedua
> ekspor — **NOL**, dan kalibrasinya lulus (`EDMState`, `EDMMaterialType`, `OLDID` ketiganya
> ditemukan). **Sistem lama tidak pernah merekam nomor dokumen addendum.**
>
> **Arah dampak bila salah:** bila dokumen ternyata tidak perlu menjadi entitas — nomornya cukup
> sebagai atribut teks pada versi — ongkosnya **satu tabel dibuang dan kolomnya dipindahkan**, tanpa
> kehilangan data. Sebaliknya, satu dokumen yang menyentuh lima kontrak akan tersimpan **lima kali
> sebagai teks**, dan tidak ada yang tahu kelimanya satu benda.
>
> **Syarat pembalikan:** `DB-16a` dan `DB-16b` menyempitkan bentuknya; `UA-19` menguji keunikan
> nomornya. **Tidak ditandai terverifikasi.**

| Entitas | Bita | Arti | Kunci alami | Induk |
|---|---|---|---|---|
| `LAYER` | 5 | satu layer non-proporsional, atau satu kelompok limit proporsional | nomor layer + bagian layer | `VERSI_KONTRAK` |
| `DETAIL_PROPORSIONAL` | 19 | ketentuan proporsional untuk satu kelompok treaty di dalam sebuah layer | kelompok treaty | `LAYER` |
| `BAGIAN` | 6 | bagian NuRe atas sebuah layer non-proporsional | (satu per layer) | `LAYER` |

**`BAGIAN` hanya ada di cabang non-proporsional.** Bagian NuRe pada cabang proporsional adalah
**atribut** `DETAIL_PROPORSIONAL`, bukan baris — kedua sisi tidak sebanding pada sumbu
masukan-versus-turunan.

### 1.3a Tabel anak paket uang per mata uang — **`F-18`, didaftarkan 24 September 2026**

| Entitas | Bita | Arti | Kunci alami | Induk |
|---|---|---|---|---|
| `NILAI_MDP` | 9 | nilai *minimum deposit premium* sebuah layer pada satu mata uang | mata uang | `LAYER` |
| `NILAI_MDP_MINIMUM` | 17 | nilai minimum dari MDP itu, pada satu mata uang | mata uang | `LAYER` |
| `NILAI_PREMI_BRUTO` | 17 | premi bruto sebuah bagian non-proporsional pada satu mata uang | mata uang | `BAGIAN` |
| `NILAI_PREMI_BRUTO_MINIMUM` | 25 | minimum dari premi bruto itu, pada satu mata uang | mata uang | `BAGIAN` |
| `NILAI_CADANGAN_PREMI` | 20 | cadangan premi sebuah detail proporsional pada satu mata uang | mata uang | `DETAIL_PROPORSIONAL` |

Kelimanya lahir dari **`P-8` golongan A**: paket uang yang di sistem lama berbentuk **daftar per
mata uang** (`MDPList[]`, `MDPMinList[]`, `GrossPremiumList[]`, `GrossPremiumMinList[]`,
`ReservePremiumList[]`) menjadi **tabel anak**, bukan kolom. Bentuknya sama persis dengan
`NILAI_PENYEBARAN` di §1.4, yang sudah berdiri lebih dulu.

> ### `F-18` — kenapa kelimanya baru terdaftar sekarang
>
> **Mereka sudah berdiri di `KAMUS-KOLOM.md` dan di `ddl-usulan/` sejak `P-8` diputuskan, dan tidak
> pernah masuk ke sini.** Sebabnya mekanis: keduanya **dibangkitkan** dari berkas definisi, yang
> membuat tabel anak itu **di dalam perkakas**; berkas ini **ditulis tangan**. Yang dibangkitkan
> tumbuh, yang ditulis tangan tidak.
>
> **Ini instans ketiga dari bentuk yang sama.** `F-4` — `PEMULIHAN_LIMIT` dan `PERISTIWA_KONTRAK`
> hilang dari sini. `L-6` — sebab yang sama melahirkan entitas yang berdiri dua kali. Dan sekarang
> lima sekaligus.
>
> | | |
> |---|---|
> | **Kenapa ini berbahaya, bukan sekadar tidak rapi** | butir 4 urutan wewenang menetapkan **berkas ini yang MENGIKAT** soal daftar entitas. Siapa pun yang menghitung ruang lingkup membukanya dan menemukan **lima entitas lebih sedikit** daripada yang ada di DDL |
> | **Arah dampak bila salah** | kelimanya sudah diputuskan `P-8`, jadi mendaftarkannya **tidak menambah keputusan apa pun** — ia menyalin keputusan yang sudah ada ke tempat yang mengikat |
> | **Yang mencegahnya terulang** | pemeriksaan **dua arah** antara berkas ini dan `KAMUS-KOLOM.md`, dijalankan setiap kali `ddl-usulan/` dibangkitkan ulang. Putaran ini menjalankannya, dan **itulah yang menemukan kelimanya** |

### 1.4 Potongan dan penyebaran — satu konsep, dua pelekatan

| Entitas | Bita | Arti | Kunci alami | Induk |
|---|---|---|---|---|
| `POTONGAN` | 8 | satu potongan yang disepakati atas premi | jenis potongan | `BAGIAN` **atau** `DETAIL_PROPORSIONAL` |
| `PENYEBARAN` | 10 | satu baris penyebaran bagian NuRe ke susunan retro internal | jenis reasuransi | `BAGIAN` **atau** `DETAIL_PROPORSIONAL` |
| `RINCIAN_PENYEBARAN` | 18 | rincian penyebaran per **jenis reasuransi** satu tingkat di bawah induknya | jenis reasuransi | `PENYEBARAN` |
| `NILAI_PENYEBARAN` | 16 | nilai sebuah rincian penyebaran pada satu mata uang | mata uang | `RINCIAN_PENYEBARAN` |

Bentuk fisik "dua pelekatan" ditetapkan sesi DDL, **dengan syarat mengikat: aturannya tidak boleh
tertulis dua kali** (INV-63).

> ### Dikoreksi 24 September 2026 — `RINCIAN_PENYEBARAN` bukan "per pihak"
>
> Kedua baris itu semula berbunyi **per pihak**, dengan kunci alami **pihak**. Keduanya tidak
> terbaca dari sumber, dan sumbernya baru terbaca setelah L-8 ditutup: `BreakDownSprdList` tidak
> pernah ada di pohon 985 simpul.
>
> `Activity/SetSpreadName.xml` mengisi `BreakDownSprdList(…).ReinsID` dari
> `pxResults(…).ReinsTypeID`. **Nama ruasnya menyebut pihak; yang ditulis ke dalamnya jenis.**
> Seluruh 866 pasangan kelas–properti disapu: identitas pihak **tidak ada sama sekali** di keluarga
> penyebaran — ia hanya ada di kelas retro (GEL-2) dan di master luar skema.
>
> **Akibat di berkas lain:** INV-47 dan INV-50 menyebut "sumbu pihak" dan sudah diperbaiki —
> `SPEC-INVARIAN.md` §4.4a. Uraian lengkap `TITIK-BUTA-POHON.md` §6.
>
> **Kedudukan entitasnya tidak berubah.** Ia berdiri sebagai **fakta terbukukan yang membawa
> penunjuk asalnya** — ADR-0036 menyebut penyebaran dengan nama, dan INV-58 menyediakan
> pengecualiannya. Yang ditentukan uji produksi hanyalah **golongan kemampuannya**, bukan bentuknya.

### 1.5 Di luar batas gelombang ini

| Entitas | Bita | Arti | Kedudukan |
|---|---|---|---|
| `RETRO_KELUAR` | 12 | bagian yang diserahkan NuRe ke retrosesioner — **arah keluar** | **[G2]** |
| `PENCAPAIAN` | 10 | realisasi premi terhadap perkiraan, per periode | **[G3]** |

`RETRO_KELUAR` **bukan** fakultatif masuk. Pembedanya **arah**, dan karena itu ADR-0020 tidak
berlaku di sini.

---

## 2. Entitas Treaty In Adjustment

**Dua. Tidak ada yang ketiga**, dan itu hasil G1c: kesembilan belas properti kelas addendum sistem
lama seluruhnya sudah punya rumah di `KONTRAK` dan `VERSI_KONTRAK`.

| Entitas | Bita | Sekat | Arti | Kunci utama | Kunci alami | Induk |
|---|---|---|---|---|---|---|
| `NILAI_SELISIH` | 13 | **ADJ** | satu besaran yang berubah pada sebuah versi terhadap versi dasarnya | `ID_NILAI_SELISIH` | besaran + kunci padanan baris | `VERSI_KONTRAK` |
| `BESARAN_DAPAT_DISESUAIKAN` | 25 | **B** | tabel acuan: besaran apa yang boleh disesuaikan, dan satuannya | `ID_BESARAN` | nama besaran | — |

### Bentuknya, dan kenapa sempit

`NILAI_SELISIH` berbentuk **sempit** — satu baris per besaran yang berubah, bukan satu kolom per
besaran. Keberatan pokok terhadap bentuk sempit, yaitu tipe tidak dapat ditegakkan karena satu
kolom nilai harus memuat uang, tanggal, dan persentase sekaligus, **tidak berlaku di sini**:
seluruh besaran yang dapat disesuaikan adalah **uang atau persentase**, nol tanggal, nol teks.

Tipenya dijaga tiga hal: `ID_BESARAN` merujuk tabel acuan bukan teks bebas; tabel acuan itu membawa
`SATUAN_BESARAN` (`UANG` / `PERSENTASE`) yang mengikat bentuk barisnya lewat constraint; dan kolom
nilainya bilangan eksak, tidak pernah teks.

### Isi awal tabel acuan

`BESARAN_DAPAT_DISESUAIKAN` **diisi 31 sebagai awal, bukan sebagai batas.** Ketiadaan sebuah
besaran di sana tidak diperlakukan sebagai larangan sampai dikonfirmasi bisnis — ke-31 itu
mencerminkan apa yang pernah dibangun orang, bukan apa yang boleh disesuaikan menurut bisnis.

---

## 3. Tabel acuan bersama

| Entitas | Bita | Sekat | Arti |
|---|---|---|---|
| `MATA_UANG` | 9 | **B** | kode dan nama mata uang |
| `JENIS_POTONGAN` | 14 | **B** | jenis potongan yang dikenal, beserta dasar perhitungannya |
| `JENIS_REASURANSI` | 16 | **B** | jenis reasuransi pada penyebaran |
| `BAHAYA` | 6 | **B** | bahaya bernama yang dapat punya batas tersendiri |
| `KELOMPOK_TREATY` | 15 | **B** | kelompok treaty |
| `KELAS_BISNIS` | 12 | **B** | kelas bisnis |

Keenamnya **wajib tabel, bukan `CHECK`** — himpunannya bertambah tanpa mengubah arti apa pun
(ADR-0038, INV-62).

---

## 4. Rujukan ke luar skema

| Entitas | Bita | Dirujuk oleh |
|---|---|---|
| `CEDANT` [luar] | 6 | `KONTRAK` |
| `ASAL_BISNIS` [luar] | 11 | `KONTRAK` |
| `REASURADUR` [luar] | 10 | `VERSI_KONTRAK` |
| `DOKUMEN` [luar] | 7 | `DOKUMEN_KONTRAK` |

**Namanya tidak disalin, hanya pengenalnya** (ADR-0023, ADR-0041). Apakah rujukannya kunci asing
lintas skema ke `POOLDATA` atau nilai yang divalidasi aplikasi adalah keputusan sesi DDL.

---

## 5. Batas penerbitan yang BELUM DIKETAHUI ISINYA

Ini bukan entitas. Ia **lubang yang sengaja digambar**, karena lubang yang tergambar lebih baik
daripada lubang yang tidak terlihat.

`TREATYINOFFER` — tabel yang selama ini dianggap sebagai penerbitan Treaty In ke sistem lain —
**tidak punya satu pun jalur penulis yang terjangkau** di sistem berjalan. Penulisnya dua,
pemanggilnya satu, dan kedua pemanggil pemanggilnya tertutup: satu ber-blok mati, satu di balik
tombol bersyarat `NEVER`.

Maka **sambungan keluar Treaty In yang sebenarnya belum diketahui.** Dua keadaan mungkin, dan
keduanya penting:

| | Artinya |
|---|---|
| hilir juga mati | tabel itu bangkai, dan ia **tidak ikut dimigrasikan** |
| ada jalur penerbitan lain | selama tujuh sesi kita salah mengira di mana batas modul ini berada, dan **sistem baru harus menggantikan jalur yang sebenarnya dipakai** |

Ditutup dua hal: **Uji AB** mengukur kapan tabel itu membeku dan berapa kontrak yang tidak pernah
terbit; dan satu pertanyaan bisnis yang bentuknya sengaja tidak menyebut nama tabel —
**angka kontrak yang keluar dari modul ini, ke akuntansi, ke retro, ke pihak lawan, datangnya dari
mana dan siapa yang mengerjakannya.**

`TREATYINOFFER` **tidak akan menjadi entitas** di model baru. Yang digambar bukan tabelnya,
melainkan **batasnya** — lihat `ERD.md` §4.

---

## 6. Hitungan

| | Jumlah |
|---|---|
| Entitas Treaty In, dalam batas gelombang ini | **29** *(diperbarui 24 Sep 2026: +`PEMULIHAN_LIMIT`, +`PERISTIWA_KONTRAK`, +`DOKUMEN_ADDENDUM`, +5 tabel anak paket uang §1.3a)* |
| Entitas Treaty In Adjustment | **1** (`NILAI_SELISIH`) |
| Tabel acuan bersama | **7** (6 + `BESARAN_DAPAT_DISESUAIKAN`) |
| Rujukan ke luar skema | **4** |
| Di luar gelombang ini | **2** (`RETRO_KELUAR`, `PENCAPAIAN`) |
| Batas yang belum diketahui | **1** |

**Seluruh nama diperiksa terhadap batas 30 bita.** Terpanjang: `BESARAN_DAPAT_DISESUAIKAN` dan
`NILAI_PREMI_BRUTO_MINIMUM`, **25** — keduanya.

> **Cacah 29 dicocokkan dua arah terhadap `2-to-spec/KAMUS-KOLOM.md`** pada 24 September 2026:
> 29 entitas Treaty In + 6 tabel acuan = **35**, sama dengan cacah kamus. Selisihnya waktu itu
> **lima** — seluruhnya tabel anak paket uang — dan itu `F-18`.

### Pencocokan itu MELEWATKAN dua entitas — `S-1`, 25 September 2026

**Kalimat di atas mencocokkan 29 + 6.** Ia tidak mencocokkan **`NILAI_SELISIH`** dan
**`BESARAN_DAPAT_DISESUAIKAN`** — kedua entitas yang berkas ini sendiri daftarkan di §2 dan §3.
Diadu utuh, cacah yang benar di dalam gelombang ini adalah **37**:

| | Jumlah | Di `KAMUS-KOLOM.md` | Di `ddl-usulan/` |
|---|---:|---:|---:|
| Entitas Treaty In | 29 | 29 | 29 |
| Entitas Treaty In Adjustment — `NILAI_SELISIH` | 1 | **0** | **0** |
| Tabel acuan bersama | 7 | **6** | **6** |
| **Di dalam gelombang ini** | **37** | **35** | **35** |

**Sebabnya terbaca, dan ia penagih yang sudah dijawab:** `SPEC-MODEL-DATA.md` §11.3 menggolongkan
keduanya *"sengaja TIDAK dikerjakan"* dengan sebab **embargo**. Embargo itu sudah lewat — modul
Adjustment sudah digrilling, sudah to-spec (enam diff **DITERAPKAN**), dan sudah to-ticket; tiket
**`06`** bertanda **PEMBUAT PERTAMA `NILAI_SELISIH`**. Sebabnya hilang, penagihnya tidak dicabut.

**Akibat yang terukur:** siapa pun yang membangun dari `ddl-usulan/` membangun 35 tabel, lalu tiket
`06` **tidak punya tabel untuk dibuat** — dan seluruh papan Adjustment bersandar padanya.

**Yang sudah dikerjakan 25 September 2026, dan yang BELUM:**

| | |
|---|---|
| **Sudah** | keduanya **digambar** di [`ERD-SKEMA-BARU.html`](ERD-SKEMA-BARU.html) beserta relasinya, dengan **5 + 7 kolom** yang tiap-tiapnya **membawa kutipan sumbernya**, dan bertanda **BELUM BER-DDL** |
| **Belum** | **nol berkas** di `2-to-spec/ddl-usulan/`, **nol pasal** di `2-to-spec/KAMUS-KOLOM.md`. **Menggambar sebuah tabel bukan membangunnya** |
| **Yang menahan** | keduanya belum melewati **§10** `SPEC-MODEL-DATA.md` — penamaan, tipe, dan keterisian tingkat atribut. Tiga hal yang sumbernya **diam** didaftar di [`COCOK-ENAM-SUMBER.md`](COCOK-ENAM-SUMBER.md) §3 |

Uraian lengkap, beserta dua temuan lain yang lahir dari pencocokan yang sama, ada di
[`COCOK-ENAM-SUMBER.md`](COCOK-ENAM-SUMBER.md). Pemeriksanya
[`alat/cocok-enam-sumber-struktur.py`](../alat/cocok-enam-sumber-struktur.py) **keluar dengan kode
1** selama selisih ini belum tertutup.
