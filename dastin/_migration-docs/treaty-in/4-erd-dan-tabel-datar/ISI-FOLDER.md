# Struktur — ERD dan tabel datar · Treaty Masuk

Bentuk data — **dirancang** dan **dibaca apa adanya dari sistem lama**. Sepuluh berkas di antaranya
**TURUNAN**: dibangkitkan oleh alat di [`alat/`](../alat/) dan **tidak pernah disunting tangan**.
Bila salah satunya berbeda dari sumbernya, **alatnya yang salah**.

~~**16 berkas**~~ **— diperbarui 25 September 2026.** Cacahnya **tidak lagi ditulis di sini**, dan
itu keputusan: ia sudah salah dua kali. Yang berlaku: **isi foldernya sendiri**.

> ### DUA KELUARGA BERKAS DI FOLDER INI, DAN JANGAN TERTUKAR — 25 September 2026
>
> | | Sumbernya | Isinya | Berkasnya |
> |---|---|---|---|
> | **Keluarga A — dari ERD v2** | `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`, **persis, tanpa perubahan** (perintah pemilik proses) | **37 tabel `T_…` SISTEM LAMA**, empat lembar Prop/Non Prop/EDM Prop/EDM Non Prop | `ERD-TREATY-IN-DAN-EDM.html` · `STRUKTUR-DATA-ERD-V2.md` · `TABEL-DATAR-ERD-V2.md` |
> | **Keluarga B — dari to-spec** | `2-to-spec/ddl-usulan/` + `STRUKTUR-DATA.md` + to-spec Adjustment | **37 entitas MODEL BARU** (`KONTRAK`, `VERSI_KONTRAK`, …) | `ERD-SKEMA-BARU.html` · `STRUKTUR-DATA-SKEMA-BARU.xlsx` · `COCOK-ENAM-SUMBER.md` |
>
> **Keduanya kebetulan sama-sama 37, dan itu kebetulan belaka** — yang satu tabel datar sistem
> lama, yang satu entitas model baru. Tiap berkas menyebut keluarganya di kepalanya sendiri.
> Keluarga A berbanner **POTRET SISTEM LAMA**; keluarga B berbanner **RANCANGAN**.

> ### YANG MUTAKHIR, DIBACA LEBIH DULU — 25 September 2026
>
> | Kalau Anda mencari… | Buka |
> |---|---|
> | **struktur data menurut ERD v2** *(keluarga A)* | **`STRUKTUR-DATA-ERD-V2.md`** — 37 tabel `T_…`, kunci, jalur Pega, padanan model baru |
> | **tabel datar menurut ERD v2** *(keluarga A)* | **`TABEL-DATAR-ERD-V2.md`** |
> | **ERD HTML menurut ERD v2, BERGARIS PENGHUBUNG** *(keluarga A)* | **`ERD-TREATY-IN-DAN-EDM.html`** — 130 kotak di empat lembar |
> | **struktur data, bentuk ERD berkotak** | **`STRUKTUR-DATA-SKEMA-BARU.xlsx`** — **37 entitas**, enam lembar, **bentuknya meniru `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`** tetapi isinya model baru |
> | **gambar skema baru** | **`ERD-SKEMA-BARU.html`** — **37 tabel**, enam bagian, dibangkitkan dari **enam sumber**. Inilah yang mutakhir |
> | **selisih antar sumber struktur** | **`COCOK-ENAM-SUMBER.md`** — hasil gerbang, tiga temuan `S-1`…`S-3` |
> | gambar skema baru, bentuk Excel | `ERD-SKEMA-BARU.xlsx` (tujuh lembar) — **35 tabel**, dari `ddl-usulan/` saja; **tanpa** kedua entitas Adjustment |
> | **nama dan tipe kolom** | `2-to-spec/KAMUS-KOLOM.md` — **mengikat**, urutan wewenang butir 5 |
> | **daftar entitas** | `STRUKTUR-DATA.md` — **mengikat**, butir 4 |
> | **relasi dan perilaku hapus** | `ERD.md` §2 — **mengikat**, ditulis tangan |
> | **selisihnya terhadap DDL** | `ERD.md` §2z — **dibangkitkan** |
>
> **Berkas POTRET SISTEM LAMA dan BASI di folder ini seluruhnya BERBANNER**, dan dua di antaranya
> sudah dipindahkan ke [`_arsip/`](_arsip/) pada 25 September 2026 karena **tergantikan utuh**.
> Yang tinggal di sini tetap dirujuk 5–7 berkas, dan banner adalah perlakuan yang lebih murah
> daripada memindahkannya.

### Rancangan sistem baru

| Berkas | Isinya | Sifat |
|---|---|---|
| [`KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md`](KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md) | gerbang G1, G1b, G1c, G2, G3, G4 beserta buktinya — **ditulis sebelum satu kotak pun digambar** | ditulis tangan |
| [`STRUKTUR-DATA.md`](STRUKTUR-DATA.md) | daftar entitas kedua modul: arti, kunci utama, kunci alami, induk | ditulis tangan |
| [`ERD.md`](ERD.md) | relasi **mengikat** — kardinalitas, keterisian, perilaku hapus, dan sekat antarmodul | **§1–§5 ditulis tangan · §2z dan §6 DIBANGKITKAN** |
| **`STRUKTUR-DATA-SKEMA-BARU.xlsx`** | **struktur data model baru dalam bentuk ERD berkotak** — 37 entitas, **enam lembar**: *Proporsional* (22) · *Non-proporsional* (26) · *Penyesuaian* (3) · *Acuan* (7) · *Daftar Relasi* (41) · *Catatan & Batas*. Kotak tiga baris: nama + kardinalitas + cacah kolom · PK/FK + rujukan acuan · **asal di sistem lama + nama tabel datar lama + kunci alami**. Dua kotak **merah** = BELUM BER-DDL; `ON DELETE` yang tidak dinyatakan **berlatar kuning** | **turunan** dari [`alat/bangkitkan-struktur-data-xlsx.py`](../alat/bangkitkan-struktur-data-xlsx.py) |
| **`ERD-SKEMA-BARU.html`** | **gambar skema baru — 37 tabel · 41 relasi · 269 kolom · enam lembar.** Kartu entitas lengkap dengan kolom, tipe, keterisian, kunci alami, **asal di sistem lama**, dan **nama tabel datar lama**; daftar relasi; **30 rujukan yang SEHARUSNYA ADA**; **17 padanan peta nama yang belum terpadankan**; cakupan nasib; entitas luar; penolakan; dan bagian Batas. Dua kartu bertanda **BELUM BER-DDL** | **turunan** dari [`alat/bangkitkan-erd-html.py`](../alat/bangkitkan-erd-html.py), rantai **enam sumber** |
| **[`COCOK-ENAM-SUMBER.md`](COCOK-ENAM-SUMBER.md)** | **hasil gerbang enam sumber** — enam pemeriksaan `A`…`F`, tiga temuan `S-1`…`S-3`, dan `L-8` diukur ulang secara mandiri (**358**, terhadap **340** yang diklaim) | **turunan** dari [`alat/cocok-enam-sumber-struktur.py`](../alat/cocok-enam-sumber-struktur.py) |
| `erd-skema-terpadu.json` | bahan mesin-baca gambar 37 tabel, beserta kutipan sumber tiap kolom Adjustment | **turunan** |
| ~~`ERD-SKEMA-BARU.xlsx`~~ | gambar skema baru — tujuh lembar, **35 tabel** | **KURANG DUA** — dibangkitkan hanya dari `ddl-usulan/`, sehingga **tanpa** `NILAI_SELISIH` dan `BESARAN_DAPAT_DISESUAIKAN`. Lihat `S-1` |
| `erd-skema-baru.json` | bahan mesin-baca xlsx di atas — **35 tabel** | **turunan**, masukan bagi `erd-skema-terpadu.json` |
| [`TEMUAN-F19-F21-RUJUKAN-YANG-HILANG.md`](TEMUAN-F19-F21-RUJUKAN-YANG-HILANG.md) | tiga temuan ronde ERD: kolom sambungan Adjustment yang tidak ada, 22 kunci asing yang hilang, dan 38 `ON DELETE` yang tidak sampai ke DDL | ditulis tangan |
| [`TABEL-DATAR.md`](TABEL-DATAR.md) | rancangan tabel datar: satu baris mewakili apa, kuncinya, asal tiap kolom, cara pengisian | ditulis tangan |
| [`TEMUAN-ADJUSTMENT-DITUNDA.md`](TEMUAN-ADJUSTMENT-DITUNDA.md) | temuan perilaku Adjustment yang **tidak diadili**, dipilah menurut garis 323/56 | ditulis tangan |
| ~~`Diagram-Skema-Tabel-TreatyMasuk.xlsx`~~ | peta kotak — 31 tabel, 42 relasi | **BASI, BERBANNER** — nol `DOKUMEN_ADDENDUM`, nol tabel `NILAI_*` |
| ~~`ERD-TREATY-MASUK.html`~~ | peta yang sama sebagai SVG | **BASI, BERBANNER** — idem |

### Bentuk sistem lama, apa adanya

Ditambahkan **23 September 2026**. Seluruhnya turunan dari
[`alat/buat-pohon-treatyin.py`](../alat/buat-pohon-treatyin.py) — sapuan **329 berkas** ekspor
Treaty In, ditambah 379 berkas Adjustment **hanya** untuk menandai simpul mana yang juga ada di sana.

| Berkas | Isinya | Sifat |
|---|---|---|
| [`struktur-treatyin-lama.md`](struktur-treatyin-lama.md) | pohon `TreatyIn` sampai kedalaman 4 + kamus field per simpul + batas bukti | **turunan** (naskahnya ditulis tangan, angkanya dari alat) |
| `pohon-treatyin.txt` | pohon §4 siap tempel — 187 baris | **turunan** |
| `datar-treatyin-lama.csv` | 985 simpul, satu baris per simpul | **turunan** |
| `datar-treatyin-kelas.csv` | 17 kelas Pega dan simpul yang memakainya | **turunan** |
| `datar-treatyin-muatan-keluar.csv` | 175 argumen ke empat prosedur Oracle | **turunan** |

### Tabel datar `T_…`

| Berkas | Isinya | Sifat |
|---|---|---|
| [`PETA-NAMA-TABEL-TREATYIN.md`](PETA-NAMA-TABEL-TREATYIN.md) | 44 tabel: kepala `TREATY_IN` + 43 anak `T_…`, dengan padanan DDL-nya + 4 keputusan bentuk | **turunan** dari [`alat/buat-peta-nama-tabel-treatyin.py`](../alat/buat-peta-nama-tabel-treatyin.py) |
| `peta-nama-tabel-treatyin.tsv` | bentuk mesin-baca yang sama, berkolom `JALUR_PEGA` | **turunan** |
| `ERD-STRUKTUR-TREATYIN.html` | ERD 44 tabel — kepala `TREATY_IN` yang sudah ada + anak `T_…`, **indentasinya = kedalaman simpul di pohon Pega**; isi tiap kotak dibaca dari hasil sapuan | **turunan** dari [`alat/buat-erd-struktur-treatyin.py`](../alat/buat-erd-struktur-treatyin.py) |

#### Dua ERD, dan bedanya bukan gaya

| | [`ERD-TREATY-MASUK.html`](ERD-TREATY-MASUK.html) | [`ERD-STRUKTUR-TREATYIN.html`](ERD-STRUKTUR-TREATYIN.html) |
|---|---|---|
| Menggambar | **entitas rancangan** | **tabel datar** turunan pohon Pega |
| Nama | §16 — `KONTRAK`, `VERSI_KONTRAK` | `T_…` + padanan §16-nya |
| Kotak | 31 | 44 |
| Relasi | 42 | 39 |
| Isi kotak | ditulis tangan dari `ERD.md` | **dibaca dari sapuan XML** — nama field dan cacah rujukan tidak diketik ulang |
| Mengikat? | tidak — `ERD.md` yang mengikat | tidak — ia laporan bentuk lama |

Keduanya **tidak saling menggantikan**. Yang pertama menjawab *"apa yang akan dibangun"*, yang
kedua *"apa bentuknya sekarang"*. Jumlah relasi lintas-sekatnya pun berbeda (2 lawan 1), dan itu
dijelaskan di dalam berkasnya: relasi kedua berpangkal pada `BESARAN_DAPAT_DISESUAIKAN`, tabel
acuan rancangan yang tidak pernah ada di sistem lama sehingga tidak punya kotak di peta struktur.

> **`ERD.md` yang mengikat untuk RANCANGAN.** Berkas gambar adalah alat bantu baca; bila berbeda,
> teksnya yang menang dan gambarnya yang salah.
>
> **`struktur-treatyin-lama.md` tidak mengikat apa pun — ia LAPORAN.** Ia memerikan bentuk sistem
> lama apa adanya, termasuk salah eja dan nama yang tidak cocok dengan isinya. Jangan pakai ia
> sebagai rancangan, dan jangan perbaiki namanya di sana.

---

## Bentuknya meniru contoh claim-non-prop — dan di tiga titik ia menyimpang

Bentuk berkas gambar mengikuti
`_migration-docs/claim-non-prop/4-erd-dan-tabel-datar/Diagram-Skema-Tabel-Gabungan.xlsx` dan
`ERD-GABUNGAN.html`: kanvas kolom sempit, kotak tabel bertumpuk, baris PK berlatar krem, baris FK
berlatar biru pucat, pita lajur berwarna, dan sheet daftar di belakangnya.

**Isinya tidak disalin.** Contoh itu bentuk, bukan sumber; klaim domain apa pun dari sana tidak
berlaku di sini.

Semula dicatat empat penyimpangan. **Satu dicabut**, jadi tinggal tiga — seluruhnya karena §16 atau `ERD.md` menang atas bentuk contohnya:

| # | Di contoh | Di sini | Kenapa |
|---|---|---|---|
| 1 | ~~nama tabel berawalan `T_`~~ | ~~tanpa awalan~~ | **DICABUT 23 September 2026 — lihat di bawah** |
| 2 | nama kolom bercampur Inggris — `CLAIM_ID`, `COVER_KEY` | seluruhnya Indonesia — `ID_KONTRAK`, `ID_VERSI_KONTRAK_DASAR` | §16 |
| 3 | daftar relasi **tanpa** kolom perilaku hapus | ada kolom **Perilaku hapus** | `ERD.md` §1 menuntut **empat** hal per relasi; tanpa kolom itu satu di antaranya hilang |
| 4 | sheet `Tali Penghubung` — pertemuan dua **berkas** | sheet `Sambungan Antarmodul` — pertemuan dua **modul** | perannya sama, subjeknya berbeda: di sini yang bertemu Treaty In dan Adjustment, bukan dua workbook |

### Penyimpangan 1 dicabut — awalan `T_` dipakai

Penyimpangan 1 berdiri di atas premis yang **salah**: bahwa `T_WORK_CLAIM` dan `KLAIM` adalah dua
pilihan yang bersaing, sehingga §16 harus memenangkan salah satunya. Keduanya sebenarnya hidup
berdampingan di satu baris yang sama — `peta-nama-tabel.tsv` modul claim-non-prop punya kolom
`NAMA_T` **dan** kolom `PADANAN_DDL`:

| `NAMA_T` | `PADANAN_DDL` |
|---|---|
| `T_CLAIM_SPREADING_RISK` | `ALOKASI_LAYER` |
| `T_GENERAL_CLAIM` | `KLAIM` |

`NAMA_T` adalah nama **tabel datar** — bentuk pipih yang memetakan pohon clipboard Pega satu lawan
satu. `PADANAN_DDL` adalah nama **objek Oracle** yang tunduk penuh pada §16. Menghapus `T_` tidak
menegakkan §16; ia menghapus satu lapisan pemetaan dan membuat tabel datar tak lagi dapat dibedakan
dari entitas rancangan.

Karena itu [`PETA-NAMA-TABEL-TREATYIN.md`](PETA-NAMA-TABEL-TREATYIN.md) memakai **kedua** kolom, dan
alatnya menolak berjalan bila ada pengenal melebihi 30 bita atau ada `NAMA_T` ganda — **diperiksa,
bukan diandaikan**.

`ERD.md`, `STRUKTUR-DATA.md`, dan kedua berkas gambar tetap memakai nama §16, karena yang mereka
gambarkan adalah **entitas rancangan**, bukan tabel datar.

---

## Yang digambar sebagai lubang, bukan sebagai kotak

Satu kotak bergaris putus merah di kaki peta: **batas penerbitan ke luar, belum diketahui isinya.**

Ia bukan entitas dan tidak akan menjadi tabel. Ia ditandai karena `TREATYINOFFER` — yang selama ini
dianggap jalur penerbitan Treaty In ke sistem lain — tidak punya satu pun penulis yang terjangkau di
sistem berjalan.

Lubang yang tergambar lebih baik daripada lubang yang tidak terlihat.

---

## Menjalankan ulang alatnya

```
cd D:\XML_NURE\_migration-docs\treaty-in
PYTHONIOENCODING=utf-8 python alat/buat-skema-treaty-masuk.py
```

Menuntut `openpyxl`. Empat alat lain tidak menuntut apa pun di luar pustaka baku:

```
PYTHONIOENCODING=utf-8 python alat/buat-pohon-treatyin.py
PYTHONIOENCODING=utf-8 python alat/buat-peta-nama-tabel-treatyin.py
PYTHONIOENCODING=utf-8 python alat/buat-erd-struktur-treatyin.py
PYTHONIOENCODING=utf-8 python alat/langkah-hidup.py "Treaty In/Activity/NamaActivity.xml"
```

Seluruh isi tiap keluaran didefinisikan di satu tempat di dalam skripnya — sunting
**di sana**, jangan di hasilnya.

### Rantai gambar skema baru — **lima langkah, berurutan**

```
cd D:\XML_NURE\_migration-docs\treaty-in
set PYTHONIOENCODING=utf-8
python alat\bangkitkan-erd-skema-baru.py      # ddl-usulan/  -> erd-skema-baru.json     (35 tabel)
python alat\cocok-enam-sumber-struktur.py     # enam sumber  -> himpunan-struktur.json  (gerbang)
python alat\lengkapi-erd-adjustment.py        # + Adjustment -> erd-skema-terpadu.json  (37 tabel)
python alat\bangkitkan-erd-html.py            #              -> ERD-SKEMA-BARU.html
python alat\bangkitkan-struktur-data-xlsx.py  #              -> STRUKTUR-DATA-SKEMA-BARU.xlsx
```

**Dua keluaran terakhir memakai bahan yang sama**, `erd-skema-terpadu.json`, dan karena itu tidak
dapat berbeda isinya. Yang berbeda hanya bentuknya: satu HTML berkartu, satu Excel berkotak.

### Rantai keluarga A — **dua langkah**, dari ERD v2

```
cd D:\XML_NURE\_migration-docs\treaty-in
set PYTHONIOENCODING=utf-8
python alat\urai-erd-v2.py              # Diagram-…-v2.xlsx -> alat/erd-v2.json  (salinan apa adanya)
python alat\bangkitkan-dari-erd-v2.py   # erd-v2.json -> HTML + STRUKTUR-DATA + TABEL-DATAR
```

**Ketiga keluarannya lahir dari satu perintah dan satu berkas**, sehingga tidak dapat saling
berbeda. Satu-satunya hal yang ditambahkan terhadap sumbernya adalah **garis penghubung** di HTML —
dan itu memenuhi janji legenda sumbernya sendiri, yang berbunyi *"garis tegak turun dari kotak
induk, lalu siku mendatar menyentuh kotak anak"* padahal di xlsx-nya garis itu tidak ada.

**Urutannya mengikat.** Langkah ketiga membaca keluaran langkah kedua; menjalankan langkah keempat
tanpa langkah ketiga menghasilkan gambar **35 tabel** — dan perkakasnya **mengatakannya**, alih-alih
diam.


---

## Berbanner 25 September 2026 — **seluruhnya, dan dua di antaranya diarsipkan**

Berkas yang **mudah dikira rancangan** sekarang menyatakan dirinya bukan. Banner disisipkan di
**tiap lembar** untuk xlsx, dan tepat sesudah `<body>` untuk HTML.

| Berkas | Banner | Keadaan |
|---|---|---|
| `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx` | **POTRET SISTEM LAMA** | berbanner, 6 lembar — **tetap di sini** |
| `Diagram-Skema-Tabel-TreatyMasuk.xlsx` | **BASI** | berbanner, 4 lembar |
| `ERD-TREATY-MASUK.html` | **BASI** | berbanner — **digantikan `ERD-SKEMA-BARU.html`** |
| `ERD-STRUKTUR-TREATYIN.html` | **POTRET SISTEM LAMA** | berbanner — gambar **sistem lama**, bukan digantikan |
| `_arsip/ERD-TREATY-IN-DAN-EDM.xlsx` | POTRET SISTEM LAMA | berbanner, 7 lembar — **DIARSIPKAN** |
| `_arsip/Diagram-Skema-Tabel-TreatyIn-dan-EDM.xlsx` | POTRET SISTEM LAMA | berbanner, 6 lembar — **DIARSIPKAN** |

### Dua yang diarsipkan — 25 September 2026

Keduanya **tergantikan utuh** oleh `…-v2.xlsx`, dan dipindahkan ke [`_arsip/`](_arsip/) sesudah
Excel ditutup. **Disimpan, bukan dihapus**: berkas yang pernah ada adalah bukti bahwa bentuknya
pernah dicoba dan kenapa ia ditinggalkan.

| Berkas | Kenapa ditinggalkan |
|---|---|
| `ERD-TREATY-IN-DAN-EDM.xlsx` | bentuk tabel biasa, **bukan** bentuk ERD `contooh.xlsx` |
| `Diagram-Skema-Tabel-TreatyIn-dan-EDM.xlsx` | versi pertama bentuk ERD — garis relasi belum tergambar, akar bertanda `1:1`, penanda "bersama" terpotong |

**Perujuknya diperiksa lebih dulu, dan hanya ada dua** — berkas ini dan
`alat/banner-potret-sistem-lama.py`. Keduanya sudah diperbarui ke jalur `_arsip/`.

> **Rujukan markdown yang putus tidak menimbulkan galat**, dan tidak ada yang menjalankannya untuk
> mengetahui. Karena itu perujuk dihitung **sebelum** berkasnya dipindahkan, bukan sesudah.

### Yang TIDAK diarsipkan, dan sebabnya disebut

Ketiga berkas **basi** tetap di tempatnya: masing-masing dirujuk **5 sampai 7 berkas**, dan banner
adalah perlakuan yang **lebih murah dan lebih aman** daripada memindahkan berkas yang dirujuk tujuh
tempat. **Basi bukan berarti tidak dipakai** — berkas basi yang berbanner tetap menjawab *"seperti
apa bentuknya sebelum ini"*.

Menjalankan ulang perkakasnya aman — ia **melewati lembar yang sudah berbanner**:

```
PYTHONIOENCODING=utf-8 python alat/banner-potret-sistem-lama.py
```

## Perkakas ronde ini — seluruhnya dapat dijalankan ulang

| Perkakas | Apa yang dilakukannya |
|---|---|
| **[`alat/cocok-enam-sumber-struktur.py`](../alat/cocok-enam-sumber-struktur.py)** | **Gerbang 0 diperlebar** — enam sumber struktur diadu, enam pemeriksaan. **Keluar dengan kode 1** selama gerbangnya gagal, sehingga dapat dipasang sebagai penjaga |
| **[`alat/lengkapi-erd-adjustment.py`](../alat/lengkapi-erd-adjustment.py)** | menambahkan kedua entitas Adjustment beserta **kutipan sumber tiap kolomnya**, dan menempelkan **asal Pega** serta **nama tabel datar lama** ke tiap entitas → `erd-skema-terpadu.json` |
| **[`alat/urai-erd-v2.py`](../alat/urai-erd-v2.py)** | **keluarga A** — mengurai `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx` **apa adanya** menjadi `alat/erd-v2.json`: 130 kotak, 41 relasi, 13 catatan, **banner ikut tersalin**. Ia menyalin, ia tidak menilai |
| **[`alat/bangkitkan-dari-erd-v2.py`](../alat/bangkitkan-dari-erd-v2.py)** | **keluarga A** — dari **satu** sumber `erd-v2.json` membangkitkan **ketiganya sekaligus**: `ERD-TREATY-IN-DAN-EDM.html` (bergaris penghubung), `STRUKTUR-DATA-ERD-V2.md`, `TABEL-DATAR-ERD-V2.md`. **Ketiganya tidak dapat berbeda isinya** |
| **[`alat/bangkitkan-struktur-data-xlsx.py`](../alat/bangkitkan-struktur-data-xlsx.py)** | membangkitkan `STRUKTUR-DATA-SKEMA-BARU.xlsx` — **bentuk ERD berkotak yang ditiru dari `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`**: kotak tiga baris, indentasi 4 kolom per kedalaman, lebar kolom 2.3, kotak dilebur 46 kolom, garis kisi mati, penanda `BERSAMA →`. **Yang ditiru bentuknya; isinya model baru, bukan potret** |
| **[`alat/hitung-perujuk-erd.py`](../alat/hitung-perujuk-erd.py)** | **dijalankan SEBELUM memindahkan berkas apa pun** — menghitung siapa merujuk siapa di seluruh korpus, memisahkan perujuk **katalog** dari perujuk **kerja**, dan mengukur sendiri dua cacatnya: nama yang bertabrakan dan nama yang cocok sebagai substring |
| ~~[`alat/cocok-tiga-sumber-entitas.py`](../alat/_arsip/cocok-tiga-sumber-entitas.py)~~ | Gerbang 0 versi lama. **DIARSIPKAN 25 Sep 2026 ke [`alat/_arsip/`](../alat/_arsip/)** — ia mencetak `lulus = True` karena **mengecualikan** kedua entitas Adjustment ke daftar `LUAR` atas dasar embargo §11.3 yang sudah lewat |
| [`alat/bangkitkan-erd-skema-baru.py`](../alat/bangkitkan-erd-skema-baru.py) | membangkitkan `ERD-SKEMA-BARU.xlsx` dan `erd-skema-baru.json` — **35 tabel, dari DDL saja** |
| [`alat/bangkitkan-erd-html.py`](../alat/bangkitkan-erd-html.py) | membangkitkan `ERD-SKEMA-BARU.html` dari `erd-skema-terpadu.json`, dan **mundur ke** `erd-skema-baru.json` sambil melaporkannya bila yang terpadu tidak ada. Nol angka ditulis tangan |
| [`alat/cocok-relasi-erd-versus-ddl.py`](../alat/cocok-relasi-erd-versus-ddl.py) | mencocokkan **relasi** §2 terhadap `FOREIGN KEY` — perkakas yang menemukan `F-19`…`F-21` |
| [`alat/tulis-blok-erd-md.py`](../alat/tulis-blok-erd-md.py) | menulis §2z dan gambar §6 ke dalam `ERD.md` |
| [`alat/banner-potret-sistem-lama.py`](../alat/banner-potret-sistem-lama.py) | banner pada xlsx potret dan basi |
| [`alat/banner-html-basi.py`](../alat/banner-html-basi.py) | banner pada HTML |

**Keenamnya melaporkan apa yang ditolaknya beserta jumlahnya, dan mencetak batas kemampuannya
sendiri di dalam keluarannya.**


---

## Dipisah per modul — 25 September 2026

`STRUKTUR-DATA-ERD-V2.md` memuat **keempat lembar** potret sistem lama sekaligus, dua di antaranya
milik jalur addendum. Atas perintah pemilik proses ia **dipecah dua**, masing-masing disimpan **di
folder modulnya sendiri**:

| Modul | Berkas | Lembar | Isinya |
|---|---|---|---|
| **Treaty In** | [`../STRUKTUR-DATA-TREATY-IN.md`](../STRUKTUR-DATA-TREATY-IN.md) | *Treaty In Prop* · *Treaty In Non Prop* | 36 tabel unik · 70 kotak · 35 relasi |
| **Treaty In Adjustment** | [`../../treaty-in-adjustment/STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md`](../../treaty-in-adjustment/STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md) | *Treaty In EDM Prop* · *Treaty In EDM Non Prop* | 35 tabel unik · 60 kotak · 34 relasi |

**Keduanya DIBANGKITKAN** [`alat/pisah-struktur-erd-v2-per-modul.py`](../alat/pisah-struktur-erd-v2-per-modul.py)
dari `alat/erd-v2.json` — **sumber yang sama** dengan berkas gabungannya, bukan dipotong dari
keluarannya. Memotong keluaran akan melahirkan dua berkas yang tidak dapat dibangkitkan ulang.

> ### YANG DIPISAH GAMBARNYA, BUKAN MODELNYA
>
> **34 dari 37 tabel dipakai kedua modul**, dan karena itu muncul di **kedua** pecahan. Itu bukan
> penggandaan melainkan kenyataan sistem lama: satu tabel yang sama dipakai jalur kontrak **dan**
> jalur addendum. Tiap baris membawa kolom **Bersama modul lain**.
>
> **Yang khas per modul hanya tiga:** `T_TREATY_FAC_LIMITS` dan `T_TREATY_INSTALLMENT_ITEM`
> (Treaty In), `T_TREATY_VALUE_BEFORE_PRORATE` (Adjustment).

**Cacahnya menutup dua arah**, dan itu diperiksa bukan diandaikan: 70 + 60 = **130 kotak**, sama
dengan sumbernya; gabungan tabel uniknya **37**, sama dengan sumbernya; **39 relasi sejati** mendarat
seluruhnya. Dua baris di *Daftar Relasi* sumbernya **bukan relasi** melainkan catatan kaki, dan
keduanya ditolak dengan sebabnya. **Tiga relasi ber-*(tidak terklasifikasi)*** tidak muncul di lembar
mana pun dan **dibawa ke kedua berkas** — membuangnya dari keduanya akan menghilangkannya, dan
memilih salah satu berarti memutuskan sesuatu yang sumbernya tidak nyatakan.
