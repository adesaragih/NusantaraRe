> Modul  : Treaty In Adjustment · Ronde A · 2026-09-24
> Peran  : penilai
> Masukan: seluruh berkas ronde A · `KEPUTUSAN-GRILLING-ADJUSTMENT.md` · `METODE-GRILLING.md` Bagian IX
> Status : TERBUKA — ronde A belum ditutup
> Sifat  : TAMBAH-SAJA

# 07 · AUDIT RONDE A

## 1. Buku besar

### 1.1 Temuan `NA-01 … NA-06`

| # | Status akhir | Berlabuh di |
|---|---|---|
| NA-01 | berdiri, kadar dibatasi | TDA-13, `EXP-1`, cabang C |
| NA-02 | berdiri, dua akhir yang mungkin | TDA-12, `UA-1`, cabang B |
| NA-03 | berdiri | §1.2, sidang `METODE` §8.1 |
| NA-04 | berdiri | kolom sisi Lacak TDA, GRL-01 butir 5 |
| NA-05 | berdiri | §5.6, cabang E |
| NA-06 | berdiri — **mencabut TDA-07**, lolos penyisiran tertutup | §4.5, §7.2, TDA-08, `UA-9`, `UA-11` |
| NA-07 | tiga dari lima field lapisan beku dapat disunting di layar addendum, tanpa syarat apa pun | GRL-05, `UA-10` |
| NA-08 | dua dari tiga kontrol `Ceding` tidak menulis `CedingID`; ketidaksinkronan tersimpan ke tabel datar | §4.6, GRL-05, `UA-10`, titipan I/F, `DB-8`, `DB-9` |
| NA-09 | dua hitungan ikut menghitung salinan terbungkus | `MA-04`, §4.3, §4.5, `00-LINGKUP` §5a |
| NA-10 | daftar tertutup penulis `TreatyIn.Position` — **penguatan mandiri** atas ADR-0055 §4.1, bukan temuan baru; yang baru jejak `"BERNARD"` | GRL-07, ADR-0052, ADR-0055, `UA-12` |
| NA-11 | pintu samping ADR-0055 §4 juga ada di layar addendum; `Force Resolve Complete(dev)` meninggalkan `RevisionState`; `TreatyInSetToDirector` seluruhnya mati | §7.5, ADR-0055 §4, `UA-9`, titipan D |
| NA-12 | daur hidup `RevisionState`; penolakan **tidak** mengosongkannya; warisannya juga mengunci addendum | §7.4, TDA-08, `UA-9`, GRL-07 |
| NA-13 | TDA-11 dikuatkan; sumber grid picker diperbaiki, konfigurasi RD saingan tidak dapat hidup | §7.6, TDA-11, titipan B |
| NA-14 | audit `TA-04`: lima jenis penulis dikalibrasi; tiga klaim berdiri, tiga angka keliru, satu daftar bocor | §4.3, §4.6, §6.2, §7.4, `tools/tulis.py` |
| NA-15 | tombol paksa di layar addendum menyimpan ke tabel kontrak dan **tidak menulis apa pun**, sambil melapor berhasil | §6.4, `REV-3` butir iv |
| NA-16 | keterlihatan = wadah × sel: dua tombol untuk dua orang bernama di divisi IT, `Force Edit` untuk setiap operator divisi IT | §7.5, `MA-07`, `REV-3` butir iii |
| NA-17a | `CedingID` teks bebas tanpa validasi di `ShowSummary` — **layarnya dijaga satu nama operator**, jadi kadarnya jauh lebih rendah daripada laporan pertama | §4.6, `IND-1`, `UA-10` |
| NA-17b | tiga field lapisan beku dapat disunting **tanpa syarat keadaan** di kedua layar, termasuk sesudah kontrak disetujui dibuka kembali | GRL-05, `UA-10` |
| NA-19 | `Save EDM(dev)` menyimpan addendum tanpa syarat status; bersama `Force Edit (dev)` memberi jalan menyunting addendum yang sudah disetujui | GRL-02 sumber PERUBAHAN ketiga, `REV-1` |
| NA-20 | penggeser `Termination` otomatis ada di seksi layar addendum tetapi pemicunya hanya-baca — **tidak terjangkau** | GRL-05 |
| NA-21 | `Revision` juga mengunci sebagian kontrol lewat `ViewState = 1` | §7.7, GRL-07 |
| NA-22 | kunci itu **tidak konsisten**; tidak satu pun dari dua belas akar uang terkunci seluruhnya | §7.7, `DB-10` ditulis ulang kedua kalinya |
| NA-18 | keempat kontrol `SetTreatyIn_Act` dibatasi peran inputor; hanya `Revision` menyentuh kontrak yang sudah disetujui | §4.5, diff eskalasi butir 1 |

### 1.2 Temuan turunan `TA-01 … TA-02`

| # | Isi | Tindakan |
|---|---|---|
| TA-01 | dua bentuk "mati" yang berbeda | usulan tambahan untuk `METODE` §2.0 |
| TA-02 | dua ekspor dapat memuat versi aturan berbeda | usulan penajaman `METODE` §8.1 |
| TA-03 | salinan terbungkus (`pyIncludedRuleXML`, `pyRuleVersionsList`) tidak dihitung dan tidak jadi bukti | usulan tambahan untuk `METODE` §2.0a |
| TA-04 | bentuk "mati" yang **ketiga** (`pyDisabled = true` pada Data Transform), dan aturan sapuan-nol | usulan perluasan `TA-01` dan tambahan untuk `METODE` §3.1 |
| TA-05 | sebelum menulis temuan, baca setiap ADR yang masih terbuka di daftar rekonsiliasi cabang yang sedang dikerjakan | usulan tambahan untuk `METODE` §3.6 — lahir dari `MA-06` |
| TA-06 | klaim keterlihatan menyebut **hasil perkalian** seluruh tingkat — sel dan setiap wadah — beserta nilai `pyVisible`, `pyContainerVisibleWhen`, dan `pyAssociatedPrivileges` tiap tingkat | usulan tambahan untuk `METODE` §2.0 — lahir dari `MA-07` dan ralatnya |
| TA-07 | parameter dianggap terkirim hanya bila **nilainya** dibaca; pencarian berbasis teks di dalam subpohon dinyatakan sebagai pencarian teks, bukan pembacaan struktur. Ditambah: pembalikan klaim tidak boleh hanya muncul sebagai perubahan artefak | usulan untuk `METODE` §2.0a dan §5.4 — lahir dari `MA-08` |

### 1.3 Pertanyaan ronde ini

| # | Keadaan |
|---|---|
| QA-1 s.d. QA-8 | **tertutup** — GRL-01 s.d. GRL-08. **Cabang A ditutup.** |

### 1.4 Perkara lama yang berubah status ronde ini

| Perkara | Dari | Menjadi |
|---|---|---|
| G2 | kesimpulan sesi ERD | **SALAH KAPRAH** |
| G4 "tabrakan mustahil" | fakta | **SALAH KAPRAH sebagai fakta**, bertahan sebagai rancangan |
| `METODE` §8.1 | fakta | **SALAH KAPRAH dua kali** |
| `METODE` §8.4 (OLDID, R01..R99) | fakta | **SALAH KAPRAH** |
| `METODE` §8.4 ROWNUM · `TEMUAN` A-2 | cacat | **DIPERKECIL** |
| `METODE` §8.6 Q2 · `TEMUAN` G4 | terbuka | **DITUTUP** |
| TDA-07 | temuan | **DICABUT** sesudah penyisiran tertutup; diganti temuan Treaty In yang lebih sempit |
| Butir eskalasi induk "selisih tidak dapat dihitung ulang" | berlaku | **DIRALAT** |

## 2. Perkara yang tertutup dua kali atau bertabrakan

| Perkara | Keadaan |
|---|---|
| "Selisih disimpan atau dihitung?" | ditutup **dua kali** dengan jawaban berbeda: GRL-01 butir 4 dan GRL-02 butir 3 menyerahkannya ke E, lalu GRL-03 menunjukkan ADR-0048 butir 3 sudah mengikat. **Diselesaikan** lewat ralat pada kedua butir; sumber kebenarannya sekarang GRL-03. |
| Definisi materialitas | **tiga** definisi hidup berdampingan: diketik pengguna (sistem lama, NA-01), diturunkan dari jenis (ADR-0049), diturunkan dari hasil selisih (ADR-0040 butir 5). **Belum diselesaikan** — milik C1, dan C1 gerbang bagi E. |
| Label butir 4 GRL-03 | PELESTARIAN **sementara**; berubah menjadi PERUBAHAN bila cabang B memilih rantai linear. Ditandai, bukan disembunyikan. |

**Tidak ada perkara yang tertutup dua kali dan dibiarkan.**

## 2a. Kesalahan penggrill, dicatat atas permintaannya sendiri

`METODE` §1.3: *"kesalahan yang dinyatakan terang lebih berharga daripada kebenaran yang dipakai
diam-diam"* — dan §1.3 menyebut penggrill punya kegagalan khasnya sendiri.

| # | Kesalahan | Akibatnya | Koreksi |
|---|---|---|---|
| **PG-01** | `METODE-GRILLING.md` dinyatakan sudah ada di direktori kerja, padahal ada di `C:\Users\Administrator\Downloads\` | satu putaran berjalan tanpa metode; seluruh rujukan §-nya sempat ditandai "atas pernyataan pemilik proses" | berkasnya disalin ke `_migration-docs/METODE-GRILLING.md`; seluruh rujukan diperiksa ulang dan cocok, kecuali §8.1 yang memang salah |
| **PG-02** | `NILAI_SELISIH` diminta ditandai "sementara" di GRL-01 butir 4, dan "disimpan atau dihitung" diserahkan ke cabang E di GRL-02 butir 3 | satu perkara tertutup dua kali dengan jawaban berbeda | diralat di kedua butir; sumber kebenarannya GRL-03 |
| **PG-03** | Satu berkas per GRL diperintahkan (GRL-01 -> GRILL-01) sebelum bentuk contoh `GRILL-05` diperiksa | susunan yang diperintahkan akan memalsukan bentuk ronde | pembedah berhenti dan melaporkan ketidakcocokannya; susunan folder-per-ronde dipakai |
| **PG-04** | Nomor `UA-9` diberikan untuk kontrak yang kehilangan status akseptasi tanpa memeriksa nomor yang sudah terpakai | tiga pengukuran berebut satu nomor | dirapikan menjadi `UA-9`, `UA-10`, `UA-11` |
| **PG-06** | Definisi golongan **GRILL** memasukkan "keputusan wewenang yang belum dijawab ADR" | bertentangan dengan `METODE` §6.6 — wewenang milik orang berwenang, dan grilling tidak dapat memutuskannya | keputusan wewenang masuk **DITUNDA**, pemilik manajemen, dicatat di daftar eskalasi — seperti tombol "(dev)" |
| **PG-05** | Perintah kerja menyebut "ketiga aturan Connect-REST"; `PENGETAHUAN.md` §1.1 mencatat **satu** per ekspor | pembedah sempat mencari dua aturan yang tidak ada | disapu apa adanya: satu aturan (`ServiceGoogle`), enam **langkah** yang memanggilnya |

Keduanya ditemukan **bukan** oleh penggrill maupun pembedah sendiri, melainkan oleh **tabrakan
antar putusan** yang muncul saat GRL-03 disusun. Itu argumen untuk mempertahankan bagian
"perkara yang tertutup dua kali" di setiap ronde.

## 3. Pertanyaan yang hanya manusia dapat jawab — daftar untuk dibantah, SIAP KIRIM

Ditulis dalam bahasa bisnis, tanpa nama tabel, kolom, atau aturan (`METODE` §4.7). Setiap butir
adalah **pernyataan untuk dibenarkan atau dibantah**, bukan pertanyaan terbuka. Tidak ada berkas
terpisah: daftarnya tinggal di sini, mengikuti pola contoh yang menaruhnya di `07-AUDIT`.

### 3.1 Untuk bagian teknik treaty

| # | Pernyataan | Bila dibantah, yang berubah |
|---|---|---|
| **DB-3** | "Satu dokumen addendum selalu mengubah tepat satu kontrak treaty — tidak pernah satu dokumen untuk beberapa kontrak atau beberapa tahun treaty sekaligus." | **GRL-04 batal.** Penyesuaian menjadi entitas tersendiri, dan jumlah tabel bertambah |
| **DB-4** | "Dua addendum yang datang terpisah tidak pernah digabung menjadi satu revisi." | **GRL-04 batal**, sebab yang sama |
| **DB-5** | "Setiap addendum berlaku sejak tanggal mulai kontrak, tidak pernah mulai berlaku di tengah periode." | **GRL-04 tetap berdiri**, tetapi cabang C dan E berubah: "berlaku sejak" menjadi atribut versi, dan perhitungan pro rata menjadi kemampuan yang harus diputuskan — dibangun atau dibuang |
| **DB-6** | "Perpanjangan atau pemendekan periode treaty tidak pernah dilakukan lewat addendum — selalu lewat kontrak baru." | **GRL-05 butir 3 ditinjau**; dicatat sebagai usulan revisi ADR-0040, bukan penyimpangan diam-diam |
| **DB-7** | "Pergantian cedant, misalnya karena merger atau pengalihan portofolio, selalu dicatat sebagai kontrak baru." | idem |

### 3.2 Untuk penyetuju addendum — kepala seksi, kepala departemen, direktur

| # | Pernyataan | Bila dibantah, yang berubah |
|---|---|---|
| **DB-1** | "Penyetuju memutuskan addendum berdasarkan selisih yang ditampilkan kepadanya; angka selisih itu harus tetap sama bila addendum dibuka lagi kapan pun." | **GRL-02 ditinjau.** Bila selisih hanya bahan pertimbangan dan boleh dihitung ulang, pembekuan angka persetujuan tidak diperlukan |

### 3.2a Untuk bagian akuntansi, pelaporan, dan hubungan dengan pihak lawan

| # | Pernyataan | Bila dibantah, yang berubah |
|---|---|---|
| **DB-8** | "Laporan, pembukuan, dan surat ke pihak lawan mengacu pada **nama** cedant — bukan pada nomor pengenalnya." | menentukan mana yang menjadi kebenaran saat keduanya tidak sinkron, dan karena itu menentukan arah perbaikan data warisan (cabang I) |
| **DB-9** | "Bila nama cedant pada sebuah addendum berbeda dari nama pada kontrak induknya, itu **selalu** kekeliruan pengetikan — bukan pergantian cedant yang sah." | bila dibantah, sebagian ketidaksinkronan adalah fakta bisnis yang harus dipertahankan, bukan data rusak yang harus diperbaiki |

### 3.2a Untuk penyetuju revisi — kepala seksi, kepala departemen, direktur

| # | Pernyataan | Bila dibantah, yang berubah |
|---|---|---|
| **DB-10** | **(ditulis ulang kedua kalinya, 26 September 2026, sesudah NA-22)** "Ada satu jenis revisi atas kontrak yang sudah disetujui yang memang **dirancang cukup disetujui kepala seksi**, tanpa naik ke kepala departemen dan direktur. Bila ada: apa batasnya?" **Yang perlu disebut kepada penjawabnya, karena layarnya tidak membatasinya:** sesudah tombol itu ditekan, yang masih dapat diubah antara lain **deductible, batas agregat, nilai dan persentase bagian, dan total angsuran** — bukan hanya keterangan. | Bila **ada batasnya**: batas itu ditulis sebagai aturan, dan layar sekarang tidak menegakkannya — cacat. Bila **tidak ada**: rantai pendek adalah akibat cara sistem dibangun, GRL-07 butir 4 berdiri utuh, dan yang tersisa hanya memberitahukan tambahan beban penyetuju (`UA-13`) |
| ~~**DB-10b**~~ | **DICABUT 26 September 2026.** Ia lahir dari pemisahan "niat versus kebocoran"; sesudah NA-22 tidak ada niat yang perlu dipisahkan | — |

**Butir ini bukan pertanyaan rancangan.** ADR-0052 butir 4 sudah memutuskan apa yang dibangun; yang
belum ada adalah **jawaban bisnis** yang mendasarinya. `METODE` §1.2 menuntut pemisahan itu: apakah
rantai pendek sebuah **kebutuhan bisnis** atau sekadar **akibat cara sistem lama dibangun**.

### 3.3 Untuk pemilik proses

| # | Pernyataan | Bila dibantah, yang berubah |
|---|---|---|
| **DB-2** | **(ditulis ulang 26 September 2026)** "Revisi selalu dibuat dari versi terakhir; tidak pernah ada yang **sengaja** membuat revisi dari versi yang lebih lama." | Bunyi lama — *"rujukan lama menunjuk versi yang disesuaikan"* — **dicabut**: addendum atas addendum adalah hal biasa dan sudah menjadi bentuk model baru (GRL-10), jadi ia bukan lagi pertanyaan. Yang tersisa adalah **percabangan**: bila dibantah, `B-4` perlu satu keadaan sah tambahan, dan baris warisan yang bercabang (`UA-1(c)`) bukan kerusakan melainkan fakta |

**Dua belas butir, dan hanya dua yang dapat membatalkan sebuah putusan** — `DB-3` dan `DB-4`, yang
membatalkan GRL-04. Selebihnya mempersempit atau memicu usulan revisi ADR, bukan membatalkan.
`DB-10` berdiri sendiri: ia satu-satunya yang menanyakan **kehendak bisnis** di balik sebuah
keputusan ADR yang sudah diambil.

## 3a. Koreksi atas butir 1 daftar eskalasi induk — DISIMPAN 24 September 2026

Butir 1 `DAFTAR-ESKALASI-MANAJEMEN.md` sudah memuat tombol yang **membuka kembali kunci** kontrak
yang sudah disetujui, dan menyebut peredamnya: *"setelah kunci terbuka, tombol simpan untuk kontrak
yang sudah disetujui **tidak ikut muncul** bagi pengguna biasa"*.

Temuan NA-06 menunjukkan peredam itu **tidak berlaku untuk dua tombol lain di layar yang sama**:
kontrol 7 dan 8 mengosongkan status akseptasi **dan menyimpan** dalam satu tindakan, sehingga tidak
ada tombol simpan terpisah yang perlu muncul. Keduanya tanpa pembatasan peran dan tanpa syarat
tampil.

**Tidak ada butir eskalasi baru.** Yang ditambahkan hanya tiga paragraf pada butir 1, disetujui
pemilik proses dan disimpan 24 September 2026, memuat: batas pemeriksaannya (terbaca dari aturan,
**belum dijalankan**), cara mengenali kontrak yang terdampak (komentar "Create Revision"), rujukan
ke `UA-11`, dan jumlah tempat yang benar — **satu**, karena delapan kemunculan pada berkas
pembungkus layar seluruhnya berada di dalam `pyIncludedRuleXML`, yaitu salinan susunan yang sama.

## 4. Sisa pengambilan bukti

| Jenis | Butir |
|---|---|
| Ekspor tambahan | `EXP-1` Field Value `EDMState` / `EDMMaterialType` |
| Kueri data | `UA-1` … `UA-11` — **tidak satu pun memblokir rancangan** (`METODE` §7.2) |
| Pertanyaan pemilik proses | `PP-1` pemilik pengerjaan to-spec induk |
| Usulan revisi ADR | `REV-1` alasan ADR-0036 |
| Usulan revisi metode | `TA-01` §2.0, `TA-02` §8.1 |

## 5. Yang TIDAK dikerjakan ronde ini, disebut supaya tidak dikira sudah

* Rekonsiliasi ADR-0040, 0049, 0052, 0055.
* Pemetaan atribut `Int-treaty_in_edm` satu per satu; peta telusur JSON modul ini; DDL.
* Verifikasi apa pun atas data produksi.
* Cabang B sampai K.

---

## 4. Titik periksa kedua — sesudah A6, A7, A8

Ditulis 24 September 2026, saat cabang A ditutup. Bentuknya sama dengan titik periksa pertama:
apa yang dikunci, apa yang salah, apa yang berubah pada cara kerja, dan apa yang diminta dari
pemilik proses.

### 4.1 Yang dikunci pada tiga pertanyaan ini

| # | Putusan | Label | Sisa yang terbuka |
|---|---|---|---|
| **GRL-06** | klasifikasi materialitas warisan tidak dihitung ulang | PELESTARIAN | batas larangannya sudah dinyatakan; titipan ke I dan C |
| **GRL-07** | ADR-0052: satu rantai; penyimpangan yang **berjalan** hanya satu | butir 1-3 PELESTARIAN, butir 4 PERUBAHAN | `DB-10`, `REV-2`, `UA-13` |
| **GRL-08** | ADR-0055: daftar keadaannya dipakai apa adanya untuk versi addendum | label **per perpindahan** | `REV-3` (empat butir) |

### 4.2 Kesalahan pembedah pada rentang ini — empat, dan tiga di antaranya satu keluarga

| # | Isi | Bentuknya | Yang berubah pada cara kerja |
|---|---|---|---|
| **MA-04** | salinan terbungkus dihitung sebagai isi aturan | menghitung yang bukan bukti | aturan bukti `00-LINGKUP` §5a; setiap hitungan memisahkan badan dan salinan |
| **MA-05** | sapuan penulis memakai nama tag yang salah; **nol dibaca sebagai jawaban** | perkakas tidak menengok ke tempat yang benar | `TA-04`: sapuan bernilai nol dikalibrasi dulu atas kasus positif. Diwujudkan `tools/tulis.py` |
| **MA-06** | temuan disebut baru tanpa membaca ADR yang masih di daftar rekonsiliasi sendiri | tidak membaca sumber yang sudah menjawabnya | `TA-05`: baca dulu setiap ADR terbuka di cabang yang sedang dikerjakan |
| **MA-07** | syarat tampil dibaca dari sel tanpa wadah — **lalu diralat keliru ke arah sebaliknya**, wadah tanpa sel | satu tingkat dibaca sebagai seluruh rantai, dua kali berturut-turut | `TA-06`: klaim keterlihatan menyebut **hasil perkalian** seluruh tingkat, beserta nilai `pyVisible` tiap tingkat |
| **MA-08** | klaim "kontrol revisi tidak berpembatas" dicabut **hanya di dalam sebuah diff**; dan pembalikan itu sendiri memuat kesalahan baru — nama parameter dibaca sebagai nilainya | daftar tag yang kurang, lalu keberadaan dibaca sebagai nilai | `TA-07` + penajaman `METODE` §5.4: pembalikan klaim dinyatakan terang di badan laporan, dengan klaim lamanya dikutip |

**Tiga dari empat adalah satu keluarga:** semuanya klaim **negatif** — "nol", "satu-satunya",
"terlihat semua orang" — yang dibuat dari pemeriksaan yang lebih sempit daripada klaimnya. Itu pola,
bukan kebetulan, dan `TA-04` sampai `TA-06` ketiganya menyerang pola yang sama dari sisi berbeda.

### 4.3 Usulan untuk `METODE`, dikumpulkan dan belum diputuskan

`TA-01` dua bentuk mati · `TA-02` versi aturan berbeda antar ekspor · `TA-03` salinan terbungkus ·
`TA-04` kalibrasi sebelum percaya nol · `TA-05` baca ADR terbuka lebih dulu · `TA-06` rantai
keterlihatan · `TA-07` nilai bukan nama, dan pembalikan dinyatakan terang. **Diputuskan di akhir
grilling**, bukan sekarang.

`TA-03` dan `TA-04` **disimpan bersama** atas permintaan pemilik proses; `METODE` tidak disunting
sekarang.

### 4.4 Hasil audit `TA-04` — ringkasannya

Tujuh klaim diperiksa ulang: **tiga berdiri**, **empat keliru angkanya**, dan **satu daftar
"tertutup" bocor satu penulis**. Tidak satu pun kesimpulan pokok yang runtuh — yang runtuh adalah
**kalimat-kalimat yang lebih kuat daripada buktinya**. Rinciannya `01-TEMUAN` NA-14.

**Titik butanya ditutup 25 September 2026**: langkah `Java` (32) disapu sebagai teks, pemetaan
respons Connect-REST diperiksa, dan 191 langkah SQL/REST dicirikan. Kadar setiap klaim yang berdiri
kini **delapan jenis penulis**, bukan lima. Rumusan "hanya memuat, tidak mencipta" **tidak dapat
dipakai** — empat aturan SQL mencipta nilai baru — dan diganti rumusan yang lebih tepat: langkah
SQL tidak pernah menulis langsung ke `TreatyIn.*`; perpindahannya selalu lewat Property-Set atau
Java.

Perkakasnya kini berkas tetap, [`tools/tulis.py`](../tools/tulis.py), dengan mode kalibrasi dan
laporan titik buta (`Java` 32 langkah, SQL/REST 191 langkah) yang dicetak setiap kali.

### 4.5 Anggaran

| Saat | Nilai | Sebab |
|---|---|---|
| ditetapkan | ±40 | `00-LINGKUP` §6 |
| sesudah titik periksa pertama | 43 -> 42 | A6 ditutup dengan rujukan |
| **sekarang** | **42 -> 40** | dua butir cabang D (**TDA-02**, **TDA-03**) sudah dijawab ADR-0055 |

**Terpakai A1 s.d. A8 = 8. Sisa 32.** Cabang A selesai dengan delapan pertanyaan; sepuluh cabang
tersisa (B s.d. K) untuk 32 pertanyaan, yaitu rata-rata tiga per cabang. Itu ketat tetapi masuk
akal, sebab sebagian besar cabang mewarisi putusan ronde A alih-alih membukanya dari nol — B, D,
dan I sudah punya sumber induk yang menjawab sebagian pertanyaannya.

**Yang diminta dari pemilik proses pada titik periksa ini:** apakah angka 40 itu tetap, dan apakah
`REV-1`, `REV-2`, `REV-3` dikerjakan sekarang atau dikumpulkan sampai grilling selesai.
