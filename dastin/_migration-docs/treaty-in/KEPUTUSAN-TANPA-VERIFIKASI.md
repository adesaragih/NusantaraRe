# Keputusan tanpa verifikasi — Treaty In

**Menggantikan** `BELUM-BISA-DI-SPEC.md`, yang dihapus.
**Untuk:** sesi penulisan spesifikasi, yang berjalan **penuh** tanpa menunggu uji maupun wawancara.
**Tanggal:** 23 September 2026

Seluruh butir yang sebelumnya tertahan **sudah diputuskan**. Setiap keputusan di sini menyebut
**dasarnya** dan **syarat pembalikannya**. **Tidak ada satu pun yang terverifikasi.**

---

## Halaman muka — seberapa jauh ini bisa dipercaya

| | Jumlah |
|---|---|
| Keputusan mengikat, kumulatif seluruh sesi | **55** |
| — bersandar pada **bukti artefak** yang dibaca langsung dari ekspor atau DDL | **21** |
| — bersandar pada **penalaran tanpa verifikasi** | **32** |
| — bersandar pada **anggapan bahwa sebuah uji kembali bersih** (Uji K dan Q) | **2** |
| Uji data yang disiapkan | 22 (Uji A–V) |
| — memblokir model | **0** — kelimanya sudah diputuskan |
| — kini berstatus **pengukur kerugian historis**, pindah ke berkas eskalasi | **15** |
| — dianggap bersih | 2 |
| Pertanyaan wawancara | 15 — **seluruhnya diputuskan**, 4 dicoret karena tidak relevan lagi |

### Berapa yang akan berubah bila uji kelak dijalankan dan hasilnya berbeda

Dari 55 keputusan, **6 akan mengubah BENTUK MODEL**. Sisanya berubah sebagai **data, baris tabel,
atau penyajian ulang angka** — bukan pembongkaran rancangan.

| # | Keputusan | Bila terbukti sebaliknya |
|---|---|---|
| 1 | Jenis addendum **satu sumbu**; materialitas diturunkan | kolom materialitas kembali sebagai sumbu kedua yang harus dijaga konsisten selamanya |
| 2 | `CurrencyRelation` disimpan tetapi **tidak menggerakkan angka** | ia naik menjadi masukan perhitungan; rumus limit bercabang |
| 3 | Adapter acuan kapasitas **gagal keras** pada ambiguitas periode | adapter dilonggarkan; aturan pemilihan susunan masuk model |
| 4 | **"Rangkaian treaty" dicoret** | benda baru muncul beserta hubungannya antar-tahun |
| 5 | Porsi limit yang dipulihkan adalah **masukan** | ia menjadi turunan, dan menuntut umpan dari modul Klaim — batas konteks baru |
| 6 | EPI dicatat pada **tingkat 100% treaty** | bukan perubahan model, melainkan **penyajian ulang angka historis** — tetapi dampaknya berkali lipat, jadi tetap dihitung di sini |

Empat dari enam **menambah** sesuatu yang belum ada; hanya nomor 1 dan 3 yang menuntut pembongkaran
bentuk yang sudah dibangun.

---

## §0 Pemisahan yang meruntuhkan sebagian besar daftar

Daftar sebelumnya mencampur dua hal yang berbeda pembacanya:

| | Uji **pemblokir model** | Uji **pengukur kerugian** |
|---|---|---|
| Tanpa jawabannya | bentuk datanya tidak bisa ditentukan | bentuknya sudah pasti |
| Yang belum diketahui | bagaimana bendanya berbentuk | berapa banyak kontrak yang terlanjur salah |
| Yang terhambat | spesifikasi | daftar eskalasi dan perkiraan pekerjaan perbaikan |
| Pembacanya | perancang | manajemen |

**Golongan kedua tidak pernah memblokir spec.** Setelah dipisah, dari 22 uji hanya **lima
kelompok** yang benar-benar memblokir model — dan kelimanya diputuskan di §1. Lima belas sisanya
pindah ke `DAFTAR-ESKALASI-MANAJEMEN.md` sebagai pengukur kerugian.

---

## §1 Lima yang benar-benar memblokir — diputuskan

Kelimanya diangkat menjadi ADR tersendiri: **ADR-0049** sampai **ADR-0053**. Ringkasannya di sini.

### 1.1 Jenis addendum dan materialitas → **satu sumbu** (ADR-0049)

Jenis addendum punya tiga nilai: perubahan estimasi, penyesuaian ke nilai aktual, administratif.
**Materialitas diturunkan dari jenis** — administratif non-material, dua lainnya material.
Kolom materialitas **tidak ada** di model baru.

**Dasar.** Kombinasi (administratif, material) tidak punya arti bisnis apa pun, dan tidak ada yang
mencegahnya terbentuk. Dua kolom yang bisa saling bertentangan, dengan satu kombinasi tak bermakna,
adalah **satu gagasan yang tumbuh dua kali**, bukan dua sumbu.

**Data warisan.** Nilai materialitas lama disimpan sebagai nilai warisan dengan asal-usulnya, tidak
dipakai menurunkan apa pun, dan tidak menjadi kolom hidup.

**Pembalikan.** Bisnis menyebut satu kasus nyata di mana addendum administratif harus dinyatakan
material, atau menyebut pembeda material yang bukan tentang akibat premi.

### 1.2 Tabel acuan kapasitas → sumber luar yang tidak dipercaya (ADR-0050)

Treaty In **tidak memilikinya**. Dibaca lewat adapter, bertipe tegas di perbatasan, **gagal keras**
bila tidak terurai. Bila lebih dari satu susunan cocok untuk satu tanggal, **adapter gagal** —
tidak memilih yang pertama, tidak memilih yang terbaru. Ambiguitas dilaporkan, bukan diselesaikan
diam-diam.

**Dasar.** Tidak ada satu pun aturan di ekspor yang menulis ke tabel itu; bentuk kolomnya bentuk
kontrak, bukan lookup; dan ketiadaan primary key maupun indeks berarti tumpang tindih memang
mungkin terjadi.

**Pembalikan.** Pemelihara ditemukan dan menyanggupi constraint — adapter boleh dilonggarkan.

### 1.3 Pembaca hilir → anggap ada setidaknya satu yang tidak dikenal (ADR-0051)

Bentuk terbitan hilir dipertahankan sebagai **turunan** dari bentuk kanonik, dan view itu
**memaparkan ID lama dalam format lama** — kode situs ditambah enam digit — supaya pemotong tujuh
karakter tetap bekerja. Setiap kontrak yang dimigrasi **menyimpan ID lamanya secara permanen**.

**Dasar — ongkos yang tidak setangkup.** Salah karena menganggap ada konsumen: memelihara view yang
tidak dibaca siapa pun, **murah**. Salah karena menganggap tidak ada: integrasi putus dan baru
ketahuan di produksi, **mahal**. Bila ongkosnya tidak setangkup, ambil sisi yang murah.

**Pembalikan.** Sensus konsumen menunjukkan tidak ada pembaca sama sekali — view dipensiunkan.
Pemensiunan mudah, pemulihan tidak.

### 1.4 Jalur persetujuan → satu rantai bawaan empat tingkat (ADR-0052)

Tabel peruteannya **kosong dari pengecualian**. Jalur alternatif Group Leader **tidak dibawa**.
Jalur revisi **tidak dipendekkan**: jalur untuk suatu perubahan tidak boleh lebih pendek daripada
jalur yang diperlukan untuk nilai hasilnya.

**Dasar.** Tidak ada bukti jalur alternatif itu pernah disahkan, dan pemendekan jalur revisi dipilih
lewat tombol tanpa aturan apa pun. Membuang keduanya **dapat dibatalkan dengan menambah satu baris
tabel**; mempertahankannya tidak dapat dibatalkan tanpa mengaudit ulang kontrak yang sudah
melewatinya.

**Pembalikan.** Satu baris di tabel perutean.

### 1.5 Mata uang dan `CurrencyRelation` (ADR-0053)

| | Keputusan |
|---|---|
| **a. Bentuk** | Mata uang per layer **tidak dibatasi** — daftar, bukan dua slot. Batas dua, bila kelak dikonfirmasi, menjadi **aturan validasi**, bukan bentuk tabel |
| **b. `CurrencyRelation`** | **Bertahan** sebagai atribut bernilai bernama, **disimpan dan ditampilkan**, tetapi **tidak dipakai dalam perhitungan apa pun** sampai bisnis menamai artinya |
| **c. Penyebut Rate on Line** | Dihitung **sekali**: tiap mata uang dikonversi sekali lalu dijumlahkan. Tidak ada percabangan atas `CurrencyRelation`, tidak ada penjumlahan per mata uang EGNPI |

**Dasar bagian b, dan ini yang terpenting:** atribut yang artinya tidak diketahui **tidak boleh
menggerakkan angka**. Menyimpannya aman; membiarkannya bercabang di rumus adalah mewarisi perilaku
yang tidak bisa dijelaskan siapa pun.

**Dasar bagian a.** Asimetri di cabang lama lebih mungkin jejak asumsi slot IDR daripada semantik
AND/OR, dan bentuk daftar menampung kedua kemungkinan tanpa kehilangan apa pun.

**Pembalikan.** Begitu bisnis menamai AND dan OR, `CurrencyRelation` boleh naik menjadi masukan
perhitungan.

---

## §2 Sembilan uji yang ternyata tidak memblokir model

Seluruhnya pindah ke `DAFTAR-ESKALASI-MANAJEMEN.md` sebagai **pengukur kerugian**. Modelnya sudah
pasti tanpa mereka.

| Uji | Kenapa modelnya sudah pasti |
|---|---|
| **L** — deposit versus minimum premium | Lantai adalah sifat besaran; selisih selalu bertanda; tidak ada penjepitan di mana pun. Deposit dan minimum adalah **dua properti terpisah**, jadi keduanya **harus** bisa berbeda — terlepas dari apakah pernah berbeda |
| **M** — penyebaran manual pilihan atau kejatuhan | Menyemai adalah tindakan eksplisit, suntingan divalidasi, menyimpan tidak pernah menarik ulang. Alasan orang memilih manual mengubah **pelatihan dan harapan kualitas data**, bukan bentuk data |
| **N** — penyimpangan penyebaran otomatis | murni pengukur kerugian |
| **O** — reinstatement bukan 100% | dua besaran mandiri, dihubungkan perkalian bukan persamaan |
| **P** — Rate on Line terbagi jumlah mata uang | perbaikannya sudah ditetapkan di §1.5c; uji hanya mengukur berapa yang terlanjur |
| **R** — kurs beku setengah ada | **sudah terjawab dari ekspor**: tidak ada satu pun aturan berjalan yang menulis `Conversion`; hanya konversi satu kali dari tabel lama. Maka **dianggap belum ada dan dibangun baru**. Baris lama yang kebetulan terisi diperlakukan sebagai **catatan kurs warisan** |
| **S** — apa yang lolos saat penjaga mati | aturan sentuh-perbaiki sudah menanganinya; uji memperkirakan berapa yang akan menuntut perbaikan |
| **T** — kapasitas proporsional pernah jalan atau tidak | model menghitung kapasitas sebagai turunan; migrasi memindahkan nilai tersimpan apa adanya; **run paritas** yang akan memunculkan ketidakcocokan |
| **U, V** — potongan basi, bagian historis | modelnya sudah pasti; keduanya pengukur kerugian |

---

## §3 Wawancara — seluruhnya diputuskan

### Diputuskan

| Pertanyaan | Keputusan | Dasar | Pembalikan |
|---|---|---|---|
| EPI pada tingkat mana | **100% treaty** | rumus pencapaian membagi premi dengan bagian NuRe sebelum membandingkannya dengan EPI; pembagian itu hanya masuk akal bila EPI di tingkat 100%. Bila di tingkat bagian, angka pencapaian berlebih berkali lipat — selisih sebesar itu terhadap target tahunan tidak mungkin luput bertahun-tahun | bila underwriting menyatakan sebaliknya, ini **bukan perubahan model** melainkan **penyajian ulang angka historis**, dan langsung naik ke daftar eskalasi |
| Penyimpangan dari susunan baku | perubahan yang menuntut **persetujuan yang sama dengan kontraknya**, pelakunya tercatat. Tidak ada model izin tersendiri | ia mengubah pembagian kapasitas NuRe, jadi bobotnya setara dengan syarat kontrak lain | bisnis menetapkan jenjang izin tersendiri |
| ROL dibaca dari mana | **dibangun sebagai besaran kelas satu**, apa pun jawabannya | seseorang mengganti rumusnya **di sistem produksi** pada April 2025. Orang tidak menyunting rumus di produksi untuk angka pajangan | — |
| Pro rata waktu pada reinstatement | model membawa **penanda prorata waktu per reinstatement, bawaannya mati** | murah dibawa sekarang, mahal dipasang belakangan. Bila tidak pernah dipakai, ia mati selamanya dan tidak merugikan siapa pun | — |
| Boleh diajukan dengan EPI / installment / jadwal kosong | **boleh** | ketiganya kelengkapan transisi, bukan invarian, dan transisi tempatnya bukan di pengajuan. Sistem lama mengizinkannya tanpa bukti kerugian. Daftar kelengkapan berbentuk data, jadi memperketatnya kelak adalah **mengisi baris**, bukan mengubah rancangan | mengisi baris di tabel kelengkapan |
| Riwayat lintas tahun pernah diminta | **"rangkaian treaty" dicoret.** Amandemen (b) pada keputusan identitas kontrak dibatalkan; tersisa hubungan "disalin dari" | tidak ada bukti permintaan; menambah benda yang tidak diminta siapa pun adalah menambah lingkup. Rantai "disalin dari" memberi titik awal bila kelak dibutuhkan | permintaan riwayat lintas tahun muncul — benda baru, bukan perubahan yang ada |
| Keuangan — pengembalian premi bisa dibukukan? | Treaty In menghasilkan **selisih bertanda sebagai fakta**. Bentuk dokumen pembukuannya urusan modul keuangan | bila keuangan tidak bisa menerima nilai negatif, itu **batasan integrasi** yang ditemukan kemudian, bukan perubahan model | — |
| Keuangan — kurs mana yang dipakai buku besar? | **tidak memblokir apa pun.** Modul memakai tabel kursnya sendiri, dan **sumbernya tercatat pada setiap nilai** | karena setiap nilai uang membawa kurs **dan sumbernya** (ADR-0039), mengganti sumber kelak adalah perubahan **data**, bukan model. Ketidaksepakatan dengan buku besar, bila ada, dapat ditemukan tanpa membongkar apa pun | mengganti sumber kurs pada data |
| Dokumen batas wewenang | **tabel perutean dengan satu rantai bawaan, tanpa ambang.** Ketika dokumennya muncul, ambang menjadi baris tabel | model tidak tersentuh oleh isi dokumen itu | satu baris tabel |
| Masa simpan arsip JSON | **disimpan tanpa batas** sampai keputusan diambil, dan masa simpan itu **parameter**, bukan konstanta di kode | menyimpan terlalu lama dapat diperbaiki; menghapus terlalu awal tidak | mengisi parameter |
| Pemilik tabel acuan dan layanan dokumen | Treaty In **tidak memiliki keduanya** dan berbicara lewat adapter | siapa yang akhirnya memiliki **tidak mengubah bentuk adapter** | — |
| Dua konteks baru di peta jalan | tetap pertanyaan **tingkat program**, tetapi **tidak memblokir spec**. Pindah ke berkas eskalasi | Treaty In dirancang terhadap kontrak baca dan adapter, apa pun yang ada di seberangnya | — |

### Dicoret — tidak relevan lagi

| Pertanyaan | Sebab |
|---|---|
| Alur kerja setelah memilih susunan penyebaran | diselesaikan oleh keputusan menyemai-eksplisit |
| Bagaimana tahu laporan terlambat | diputuskan: pembacaan turunan, dihitung saat dibaca, tidak disimpan |
| Dasar perhitungan tiap nama potongan | **sudah terjawab dari artefak.** Keempat nama yang benar-benar ditulis ke daftar potongan — `"Brokerage fee"`, `"Facultative Brokerage fee"`, `"Overiding Commision"`, `"Comm to NuRe"` — seluruhnya **persentase dari premi bruto**, dan `CalculateDeduction` mengonfirmasinya (`.Deduction = DeductionPct/100 × .GrossPremium`). Profit commission tidak pernah masuk daftar itu. **Yang tetap dibawa: dimensi dasar perhitungan pada setiap baris**, supaya potongan berbasis laba atau berjumlah disepakati muat kelak tanpa mengubah bentuk |
| Kosong dibandingkan `false` di Pega | tidak ada bendera semacam itu di model baru |

---

## §3b Satu butir yang terlewat dari penyisiran, dan saya putuskan dengan pola yang sama

Dari lima belas open question, empat belas tersentuh §3. **Satu tidak**, dan ia justru salah satu
yang menyentuh batas konteks:

> **Porsi limit yang dipulihkan sebuah reinstatement — syarat yang diketik, atau turunan dari
> seberapa besar limit tergerus klaim?**

Di banyak wording pasar, limit dipulihkan sebesar yang tergerus, dan yang dinegosiasikan hanya
tarif preminya. Bila itu yang berlaku di NuRe, porsi pemulihan pindah dari masukan ke turunan —
dan ia menuntut umpan dari modul Klaim.

**Keputusan: porsi limit yang dipulihkan adalah MASUKAN — syarat kontrak yang diketik.**

**Dasar, tiga lapis:**

1. **Artefak.** Di sistem lama, kedua persentase reinstatement adalah field yang diketik orang.
   Tidak ada satu pun aturan di ekspor yang menghubungkan reinstatement dengan data klaim — tidak
   ada pembacaan, tidak ada parameter, tidak ada rujukan.
2. **Batas konteks.** Menjadikannya turunan menuntut umpan lintas konteks dari Klaim yang hari ini
   **tidak ada sama sekali**. Membangunnya berarti membuka batas konteks baru di gelombang 1 atas
   dasar dugaan.
3. **Ongkos yang tidak setangkup.** Bila kelak terbukti turunan, penambahannya **bersifat aditif**:
   nilai yang diketik menjadi nilai bawaan, dan penurunan dari gerusan klaim menimpanya. Sebaliknya,
   bila dibangun sebagai turunan dan ternyata masukan, tidak ada tempat menyimpan angka yang
   dinegosiasikan.

**Pembalikan.** Underwriting memastikan wording NuRe menyatakan limit dipulihkan sebesar yang
tergerus. Saat itu porsi pemulihan menjadi turunan, dan **batas konteks dengan Klaim dibuka** —
ini satu dari enam keputusan yang mengubah bentuk model.

---

## §4 Pemeriksaan di Pega — tidak ada yang memblokir

| Butir | Keputusan |
|---|---|
| Persentase bagian dan retro terbuka di jalur revisi | model **mengunci keduanya** terlepas dari hasil pemeriksaan. Pemeriksaan hanya mengukur celah historis → pindah ke eskalasi |
| Delapan ejaan bendera tampilan | tidak ada yang diangkut; bendera itu sendiri dihapus (ADR-0046) |
| Kosong dibandingkan `false` | tidak ada bendera semacam itu di model baru — dicoret |
| Tombol paksa di layar penawaran — **siapa yang bisa melihatnya** | **dikerjakan dan selesai 23 Sep 2026.** Hasilnya di bawah |

### Tombol paksa: dua pertanyaan, dipisah karena pemilik dan waktu tunggunya berbeda

Butir ini semula satu. Ia sebenarnya dua, dan menggabungkannya membuat yang satu menunggu yang lain
tanpa perlu.

| Pertanyaan | Dijawab oleh | Status |
|---|---|---|
| **Apakah ia terjangkau hari ini?** | pemeriksaan layar | **selesai** — hasilnya di tabel bawah |
| **Apakah ia pernah dipakai?** | data | **Uji H-4 dan H-5**, berjalan bersama uji lain → butir 1 daftar eskalasi |

Hasil pemeriksaan layar, dibaca dari **kondisi**-nya, bukan dari labelnya — keenamnya berlabel
`(dev)` dan label itu tidak menentukan apa pun:

| Tombol | Yang dilakukannya | `pyVisible` | Terlihat oleh |
|---|---|---|---|
| `TreatyInForceResolveComplete` | status → `Resolve Complete`, posisi dikosongkan, lalu **disimpan** | `OTHER` | dua nama: `ALDO SAPUTRA`, `Daniel Suhana` |
| `TreatyInReturntoInputor` | status dikosongkan, posisi → admin | `OTHER` | dua nama yang sama |
| `TreatyInForceEdit` | `ViewState = 0` — **membuka kunci** | **`ALWAYS`** | **semua pengguna** |
| `TreatyInSetToDirector` | posisi → direktur, **melompati dua tingkat** | dipanggil `TreatyInTestAgent` | lewat layar penawaran |
| `TreatyInSetValue` | status → `"test"` | — | enam layar |

**Tiga hal yang keluar dari pemeriksaan ini, dan dua di antaranya membalik dugaan awal:**

1. **`pxLimitedAccess = Dev` bukan pembeda.** Ia melekat pada **seluruh** aturan di ruleset `GISFW`,
   termasuk mesin persetujuan `Akseptasi_DT` itu sendiri. Memakainya sebagai tanda "ini alat
   pengembang" akan menghasilkan kesimpulan yang salah.

2. **Yang menyetujui tanpa approver justru TERKUNCI; yang membuka kunci justru TERBUKA.** Dugaan
   awal menempatkan `ForceResolveComplete` sebagai temuan terbesar. Kondisinya membantah itu: ia
   dibatasi dua nama. Sebaliknya `ForceEdit` ber-`pyVisible = ALWAYS`, yang membuat kondisi di
   sebelahnya **tidak dibaca sama sekali** — ia terlihat semua orang. `ViewState = 1` mematikan
   sekitar 178 kendali dan menutup 142 kondisi tampilan; `ForceEdit` membuka semuanya sekaligus.

3. **Pintunya terbuka, jalan keluarnya tidak terlihat.** Setelah kunci terbuka pada kontrak
   berstatus `Resolve Complete`, seluruh kendali penyimpanan yang ditemukan pada `Section` yang
   diperiksa **tidak ikut muncul** bagi pengguna biasa — syaratnya menuntut status bukan
   `Resolve Complete`, atau `ViewState == '1'`, atau divisi IT. Batas pernyataan ini disebutkan apa
   adanya: yang diperiksa adalah **definisi tombol**, bukan sistem yang berjalan.

Ketiganya tidak mengubah satu pun keputusan model — keenam tombol tidak dibawa (ADR-0055 §4) — dan
mengubah **bobot** butir eskalasinya.

---

### §4b Sapuan `NEVER &&` — kemampuan yang pernah ada lalu dimatikan

Bentuk `NEVER && ‹syarat sungguhan›` bukan kebetulan: ia yang terjadi ketika seseorang **mematikan**
sesuatu dengan menambahkan awalan, bukan menghapusnya. Syarat di belakang `NEVER` karena itu
**merekam penjaga yang dulu berlaku** — dan penjaga itu tidak tertulis di dokumen mana pun.

Sapuan kedua ekspor menemukan **delapan**, dan masing-masing dua hal sekaligus: kemampuan yang
pernah ada, dan jejak wewenang yang tidak terekam di tempat lain.

| Syarat yang dimatikan | Di mana | Penjaga yang terekam di belakangnya |
|---|---|---|
| `NEVER && TreatyMasterInEDM` | `InputTreatyInOffer`, `TreatyInNONProportional`, `TreatyInTabsProportional` | berlaku ketika kontrak dibuka dalam modus master addendum |
| `NEVER && OperatorID.pyUserName = 'ALDO SAPUTRA'` | `InputTreatyInOffer` (tombol *populate treatyindetail*) | satu orang bernama |
| `NEVER && OperatorID.pyUserName = 'ALDO SAPUTRA'` | `TreatyInActionButtons` (*RemoveLastComment*) | satu orang bernama |
| `FALSE && …pyWorkBasketName = 'ReasTreatyInAdmin' && Position = '' && StatusAkseptasi = 'Resolve Complete'` | `TreatyInActionButtons` (*Submit Revision*) | **admin, pada kontrak yang sudah disetujui** — yaitu tombol pengajuan revisi |
| `FALSE && TreatyIn.ViewState != '1'` | `TreatyInTabsNPValueDifferenceProRate`, `…_NoProRate` | berlaku selagi kontrak dapat disunting — layar selisih prorata |

**Tidak satu pun diadili di sini.** Daftar ini dua hal: kemampuan yang mungkin masih dibutuhkan
orang dan karena itu perlu ditanyakan sebelum dianggap tidak ada, dan jejak siapa yang dulu boleh
memakainya.

Yang paling layak ditanyakan: **tombol pengajuan revisi** (`Submit Revision`). Penjaganya masuk
akal dan lengkap — admin, kontrak sudah disetujui — dan ia dimatikan. Bila revisi memang dilakukan
orang hari ini, ia dilakukan lewat jalur lain.

Dicatat terpisah: `1=2`, `FALSE`, dan `NEVER` **telanjang** — tanpa syarat di belakangnya — muncul
ratusan kali dan **tidak** merekam apa pun. Keduanya tidak dicampur.

---

### §4c Pemeriksaan ulang klaim yang dipakai MEMBUANG

Klaim berbentuk "tidak pernah ditulis" aman ketika dipakai **menunda** dan berbahaya ketika dipakai
**membuang**: bila ia salah, sebuah fakta nyata hilang dari model diam-diam. Karena satu klaim
sejenis sudah terbukti salah (`CurrencyList`, lihat `INVENTARIS-STRUKTUR-DATA.md`), yang dipakai
membuang diperiksa ulang dengan penyisiran dua tingkat.

| Atribut | Alasan buangnya | Hasil pemeriksaan ulang |
|---|---|---|
| `ContractRefNo` | tidak pernah ditulis + tidak ada di DDL + `pyReadOnly` | **nol penugasan**, dalam bentuk apa pun, awalan apa pun — bertahan |
| `BrokeragePct` | nama parameter panggilan Oracle, bukan properti | **nol penugasan** — bertahan |
| `CedingStatusActive` | salinan status master | **nol penugasan** — bertahan |
| `SourceStatusActive` | salinan status master | **nol penugasan** — bertahan |
| `FacShare`, `FacShareBrokerage` | salinan dari `FacultativeShare` | penugasan ada dan **HIDUP** (`TreatyInXOLAddSpreading` langkah 6) — salinannya nyata, buangnya bertahan, dan temuan *salinan basi* bertahan |

Pencariannya diperluas melampaui awalan `TreatyIn.`: bentuk relatif dan pemetaan kolom RDB ikut
disapu. Kelimanya bertahan.

Yang dipakai **menunda** — `AccountingMode` dan pasangannya — sengaja **tidak** diperiksa ulang:
bila klaimnya salah, akibatnya hanya satu pertanyaan yang tidak perlu.

---

## §5 Yang tetap tidak dikerjakan, dan itu benar

| Hal | Sebab |
|---|---|
| Tabel selisih Adjustment, layar, dan alurnya | **embargo tetap berlaku**; hanya *seam*-nya (ADR-0048) |
| ERD | sesi berikutnya; bahannya di `INVENTARIS-STRUKTUR-DATA.md` |
| Perilaku layar dinamis | direkonstruksi dari definisi Section saat menulis spec |

### Dua yang dipecah, dan bagian pertamanya dikerjakan sekarang

**Daftar harapan selisih untuk uji paritas.** Dipecah dua:

| Bagian | Kapan | Isinya |
|---|---|---|
| **Bentuk** | **sekarang, sebelum run paritas** | per kelas cacat: arah selisihnya, rumus penyimpangannya, dan tanda pengenalnya |
| **Populasi** | menunggu run paritas | kontrak mana saja yang terdampak, dan berapa rupiah |

> **Penjaga yang ditulis setelah melihat yang dijaganya tidak menjaga apa pun.** Bentuk harus
> ditulis lebih dulu, atau uji terima berubah jadi latihan menjelaskan.

**Kumpulan data uji.** Kumpulan kurasi tetap keluaran uji — tetapi **run paritas sendiri tidak
menunggu**. Ia berjalan atas seluruh data yang dimigrasi, dan justru run itulah yang **menamai**
kontrak mana yang masuk kumpulan kurasi. Urutannya: migrasi → run paritas → kumpulan kurasi, bukan
sebaliknya.

---

## §6 `OLDID` arti ketiga — DUA atribut, bukan satu

**Diputuskan 24 September 2026, oleh pemilik proses, atas pertanyaan blok B butir 7.1.**
**Bacaan (ii): benda berbeda.** `NOMOR_PENAWARAN_WARISAN` berdiri sendiri, di samping
`NOMOR_KONTRAK_WARISAN` yang tetap berasal dari `ID`.

### Dasar — dan dasarnya ONGKOS, bukan kemungkinan

Ekspor memang tidak memisahkan kedua bacaan, dan itu diterima apa adanya. Yang memutuskan bukan
bacaan mana yang lebih mungkin, melainkan **ongkos salahnya ke dua arah**, yang tidak setangkup:

| Pilihan | Bila ternyata salah | Dapat dibatalkan? |
|---|---|---|
| **dua atribut** (dipilih) | satu kolom yang ternyata selalu kosong atau menggandakan tetangganya | **ya** — kolomnya dicabut, tanpa kehilangan apa pun |
| satu atribut | dua pengenal warisan larut jadi satu; nomor penawaran pra-Pega **musnah** | **tidak** — penulisnya konversi sekali jalan Juli 2019, dan pembacanya (`TREATYINOFFER`) sudah mati serta isinya bercampur |

> Nilai yang penulisnya sudah tidak berjalan dan pembacanya sudah mati **tidak punya tempat kedua**.
> Melarutkannya adalah satu-satunya kesalahan di sini yang tidak dapat dibatalkan.

Dan ia sejalan dengan dua hal yang sudah ditetapkan sebelumnya, sehingga bukan keputusan baru
melainkan penerapan: `PENGETAHUAN-MIGRASI-TREATY-IN.md` §5.2 sudah menyatakan **`OLDID` menjadi tiga
atribut**, dan larangan *"nama yang berarti dua hal di sistem lama tidak dibawa"* melarang dua arti
berbagi satu kolom.

### Yang WAJIB menyertainya — aturan pembeda per baris

Dua kolom tanpa aturan pembeda hanya memindahkan kebingungan satu tingkat ke hilir. Maka migrasi
membaca, per baris, dalam urutan ini:

| Urutan | Bacaan | Penandanya di data lama |
|---|---|---|
| 1 | arti 2 — kontrak asal salinan | `ID` bernilai penanda salinan (`"UnknownId"`) |
| 2 | arti 3 — nomor penawaran warisan | `Information` memuat penggal `" No Offer: "` **dan** nilainya sama dengan `OLDID` |
| 3 | arti 1 — versi pendahulu | selebihnya |

Penanda kedua **terbaca dari ekspor** — `TreatyInMappingDataconvert` menulis keduanya dari nilai yang
sama — sehingga bagian ini **`EVIDENCED`, bukan `DECIDED`**. Yang `DECIDED` hanya jumlah atributnya.

Baris yang **tidak cocok dengan satu pun** dari ketiganya tidak ditebak: ia mendarat di keluaran
migrasi sebagai baris tak terpetakan, sejalan ADR-0054.

### Syarat pembalikan

| Bila… | Maka |
|---|---|
| **Uji AL** menunjukkan populasi bernomor penawaran **persis sama** dengan populasi `ID`, dan nilainya identik baris per baris | `NOMOR_PENAWARAN_WARISAN` **dicabut**; `NOMOR_KONTRAK_WARISAN` tetap berasal dari `ID`. Satu kolom dihapus, tidak ada nilai yang hilang |
| Uji AL menunjukkan kedua populasi **beririsan sebagian** | keduanya bertahan, dan irisannya menjadi temuan tersendiri — bukan alasan melarutkan |
| Uji AL **tidak dapat dijalankan** (izin kueri tidak turun) | keputusan ini **tetap berlaku**, karena sisi murahnya memang dipilih untuk keadaan tanpa jawaban |

> Ketiga cabang sudah punya putusannya, jadi hasil Uji AL **tidak menahan** satu pun pekerjaan.
> Ia mengubah jumlah kolom, bukan bentuk model.

### Yang TIDAK berubah karena keputusan ini

- **Uji AC selamat** — nomor revisi dibaca dari pola `ID` (`@substring(ID,10,12)`), bukan dari `OLDID`.
- **Eskalasi butir 3 bertambah kalimat kedua**, dan itu sudah benar: `TREATYINOFFER` bukan hanya
  berhenti diisi — **isinya bercampur tiga arti**, karena `SaveData1.NOOFFER = TreatyIn.OLDID`
  ditulis tanpa syarat. Sebuah tabel yang bercampur arti bukan dasar rekonsiliasi, dan itu keterangan
  yang dibutuhkan pembaca eskalasi.

---

# §7 Empat keputusan bentuk fisik — `KTV-A`, `KTV-B`, `KTV-C`, `KTV-D`

**Ditambahkan 24 September 2026, putaran penutup to-spec.**
**Terverifikasi: NOL.** Ketiganya diputuskan **sesudah** `TDA-17`, `F-2`, dan `F-3` ditutup, sebab
presisi ditetapkan **per kolom** dan daftar kolomnya belum lengkap sebelum itu.

> **Alat yang dipakai ketiganya sama**, dan ia yang membuat ketiganya dapat diputuskan tanpa uji:
> **ambil sisi yang ongkos salahnya murah, lalu tulis kapan ia masih dapat dibalik.** Ketiganya
> dapat dibalik **sebelum data dimuat**, dan tidak satu pun sesudahnya. Kalimat itu bagian dari
> keputusannya, bukan hiasan.

---

## KTV-A — Presisi longgar per kelompok tipe, **beserta alasan tiap angkanya**

**Label: BARU.** Tidak ada padanannya di sistem lama — sistem lama sebagian besar memakai `NUMBER`
tanpa presisi sama sekali.

### Dasar umum

Terlalu lebar di Oracle **murah**: `NUMBER` disimpan panjang-berubah, jadi presisi yang
dideklarasikan **tidak memakan tempat**; `VARCHAR2` menyimpan sepanjang isinya. Terlalu sempit
**memotong data**, dan potongannya **baru ketahuan sesudah data masuk**. **Ongkosnya tidak
setangkup**, dan pada ketidaksetangkupan seperti itu sisi murahnya yang diambil.

### Tabel presisi — setiap angka membawa alasannya

| Kel. | Isi | Ditetapkan | Alasan angkanya |
|---|---|---|---|
| `U1` | nilai uang | `NUMBER(38,20)` | **38 adalah presisi maksimum Oracle** — tidak ada yang lebih longgar untuk dipilih. Skala **20** bukan angka yang diambil dari udara: sistem lama membulatkan rantai perhitungannya pada **12** desimal (`@divide(TreatyIn.ProRatePercent,100,12)`) dan **8** desimal (`@divide(local.TotalPremi, Local.TotalLimit, 8)`); hasil kali dua pembulatan itu menuntut **20**. Tabel datar sistem lama sendiri hanya memakai `NUMBER(20,4)` — **kita lebih longgar daripada sumbernya, bukan lebih sempit** |
| `P1` | persentase dan porsi | `NUMBER(38,20)` | idem. Porsi ikut terkalikan di rantai yang sama |
| `P2` | tarif brokerase dan pajak | ~~`NUMBER(11,8)`~~ **`NUMBER(38,20)` — DISERAGAMKAN** | dua alasan, dan yang pertama berupa data nyata: **`9989998` tidak muat** di `NUMBER(11,8)`, yang hanya menampung tiga angka di depan koma — dan angka itu **ditulis sistem lama ke `ROLPct`** setiap kali limitnya nol (`TETAPAN-DI-KODE.md` §2b). Kedua: dua kelompok yang sama-sama berarti "persentase" dengan lebar berbeda **mengundang orang berikutnya memilih yang salah**, dan ia tidak akan tahu apa yang dirusaknya |
| `K1` | kurs | `NUMBER(38,20)` | idem `U1`. Kurs adalah pengali, jadi ia **menyumbang** desimal ke rantai, bukan menerimanya |
| `L1` | limit layer dan premi deposit | `NUMBER(38,20)` | idem `U1` |
| `C1` | cacah dan nomor urut | `NUMBER(9)` | sembilan angka ≈ satu miliar. Tidak ada cacah di modul ini yang mendekatinya |
| **`C1` pada kolom bernama `ID_*`** | pengenal buatan | **`NUMBER(19)`** | **ini koreksi cacat, bukan pelonggaran.** §10 memakai `C1` untuk kunci utama **dan** untuk cacah, sementara `R` (kunci asing) bertipe `NUMBER(19)`. Akibatnya kunci asing `NUMBER(19)` menunjuk kunci utama `NUMBER(9)` — **lebarnya berbeda di dua ujung satu relasi**. Aturan yang berlaku sekarang: **kolom yang namanya diawali `ID_` memakai `NUMBER(19)`, apa pun kelompok tipenya di §10** |
| `R` | kunci asing | `NUMBER(19)` | tidak berubah; ia yang menjadi acuan baris di atasnya |
| `T` | teks | ~~`VARCHAR2(255 CHAR)`~~ **`VARCHAR2(1000 CHAR)`** | **255 lebih sempit daripada sistem lama**, dan itu terbaca dari DDL-nya sendiri: `TREATYINDETAIL` dan `PROPORTIONALARRG` memakai **`VARCHAR2(1000)`** pada 52 kolom. Memigrasikan kolom 1000 ke kolom 255 **memotong data pada hari peralihan** |
| **`T` pada tujuh kolom teks bebas** | narasi manusia | **`VARCHAR2(4000 CHAR)`** | `KETERANGAN`, `CATATAN`, `PENGECUALIAN`, `KETENTUAN_KHUSUS`, `CATATAN_BORDEREAUX` pada `VERSI_KONTRAK`; `NILAI_SEBELUM` dan `NILAI_SESUDAH` pada `JEJAK_PERUBAHAN`. Sistem lama memakai **`VARCHAR2(4000)`** untuk teks bebasnya (`ACHIEVEMENT`, `T_STORAGE_IMAGE`). Kedua kolom jejak menampung **nilai ruas apa pun sebagai teks**, jadi lebarnya harus sebesar ruas terlebar yang mungkin berubah |
| `D` | tanggal | `DATE` | sistem lama tidak menyimpan zona waktu di mana pun, jadi `TIMESTAMP WITH TIME ZONE` akan **mengarang informasi**. `DATE` Oracle membawa jam-menit-detik, yang cukup untuk `WAKTU_PERUBAHAN` |
| `E` | himpunan tertutup | `VARCHAR2(40 CHAR)` | nilai enumerasi **terpanjang di seluruh korpus** adalah `TBD_PLA_ARTI_BELUM_DIKETAHUI` — **28 bita**. 40 memberi ruang untuk satu nilai yang lebih panjang tanpa mengubah skema |

### Satu batas yang dinyatakan, bukan disembunyikan

`VARCHAR2(4000 CHAR)` **sah dideklarasikan**, tetapi pada basis data ber-`AL32UTF8` batas
**bita**-nya tetap 4000. Teks 4.000 aksara yang memuat aksara berbita ganda **tetap dapat ditolak
saat disisipkan**, kecuali `MAX_STRING_SIZE = EXTENDED`.

**Tidak ditambal di sini**, sebab tidak ada instans yang dapat ditanyai (`L-3`). Yang mengukurnya
**`Uji AP`** — panjang teks terpanjang yang benar-benar ada di ketujuh kolom itu.

### Arah dampak bila salah

Salah ke arah longgar: kolom lebih lebar daripada yang pernah dipakai. **Tidak ada data yang rusak
dan tidak ada ruang yang terbuang** — Oracle menyimpan sepanjang isinya. Salah ke arah sempit:
**data terpotong pada hari peralihan**, dan yang terpotong tidak dapat dikembalikan.

### Syarat pembalikan

**Sesi DDL berikutnya boleh mempersempit, dan hanya SEBELUM data dimuat.** Sesudah itu tidak:
mempersempit kolom yang sudah berisi menuntut memeriksa setiap baris, dan baris yang tidak muat
tidak punya tempat untuk pergi. **`Uji AP` dan `Uji AQ`** memberi angkanya.

---

## KTV-B — Dua induk polimorfik menjadi **dua kolom bernama + `CHECK` tepat satu terisi**

**Label: BARU.** Bentuk fisiknya memang belum pernah ada.

### Keputusan

`POTONGAN` dan `PENYEBARAN` masing-masing berinduk **dua** — `BAGIAN` (cabang non-proporsional) dan
`DETAIL_PROPORSIONAL` (cabang proporsional), §12.3 dan §14.1.

| Yang dihapus | Yang menggantikannya |
|---|---|
| `POTONGAN.ID_INDUK_POTONGAN` | `ID_BAGIAN` **nullable** + `ID_DETAIL_PROPORSIONAL` **nullable**, masing-masing dengan **kunci asing sungguhan**, ditambah `CHECK` yang menuntut **tepat satu** terisi |
| `PENYEBARAN.ID_INDUK_PENYEBARAN` | idem |

**Bukan tabel jembatan, bukan kolom diskriminator.**

### Dasar

1. **`INV-17` melarang rujukan yang sasarannya bergantung nilai kolom lain.** Diskriminator — satu
   kolom pengenal ditambah satu kolom "induknya yang mana" — **melanggarnya secara harfiah**. Dua
   kolom bernama **tidak melanggarnya**: masing-masing menunjuk satu tabel dan hanya satu tabel.
2. **Basis data ikut menjaga.** Pada bentuk lama, `ddl-usulan/` memuat pernyataan keputusan
   *"basis data TIDAK menolak baris yatim pada kolom ini"*. Dengan dua kolom bernama, **dua kunci
   asing berdiri**, dan baris yatim menjadi mustahil.
3. **Arahnya dapat dibalik.** Dua kolom dapat dilebur menjadi jembatan kelak; **jembatan sulit
   dibongkar** — ia membawa barisnya sendiri, dan barisnya sudah punya pembaca.
4. **`INV-63` selamat.** Syarat mengikat §14.1 berbunyi *"aturannya tidak boleh tertulis dua kali"*.
   Yang tidak boleh berganda adalah **aturan potongan**, bukan kolom kunci asingnya. Satu tabel
   `POTONGAN`, satu rumus, dua pelekatan — persis seperti yang `CalculateDeduction` lakukan.

### Akibat pada kunci alami, dan ia justru membaik

`INV-15` dan `INV-16` berbunyi *"unik di dalam induknya"*. Sebelumnya satu `UNIQUE`; sekarang
**dua**, satu per pelekatan:

```
UNIQUE (ID_BAGIAN,              ID_JENIS_POTONGAN)
UNIQUE (ID_DETAIL_PROPORSIONAL, ID_JENIS_POTONGAN)
```

Oracle **melewatkan** baris yang salah satu kolom kuncinya `NULL`, sehingga tiap constraint hanya
menjaga pelekatannya sendiri — yang memang artinya. Pada bentuk lama, satu `UNIQUE` atas
`ID_INDUK_POTONGAN` **mencampur dua ruang pengenal**: baris `BAGIAN` bernomor 7 dan baris
`DETAIL_PROPORSIONAL` bernomor 7 akan bertabrakan meskipun keduanya sah.

> **Ini bukan keuntungan sampingan, ia pembetulan.** Bentuk lama **menolak data yang sah** —
> kesalahan yang sama bentuknya dengan ketiga kekeliruan lingkup yang sudah tercatat di modul ini.

### Arah dampak bila salah

Salah: dua kolom yang salah satunya selalu kosong — **satu kolom kosong per baris**. Salah ke arah
jembatan yang ternyata tidak perlu: sebuah tabel yang harus dibongkar sesudah punya pembaca.

### Syarat pembalikan

**Bila induknya bertambah lebih dari dua**, bentuknya ditinjau: tiga kolom masih dapat dibaca, lima
tidak. Titik tinjaunya **induk ketiga**, bukan jumlah baris.

---

## KTV-C — Paket uang tanpa denominasi memperoleh kolom mata uang **bernama, boleh kosong**

**Label: BARU.**

### Koreksi terhadap rumusan yang diterima — dinyatakan, bukan diselaraskan diam-diam

Perintah yang diterima menyebut **enam** paket golongan C, dan meminta masing-masing diberi kolom
bernama **`KODE_MATA_UANG`**, *"sama dengan golongan B pada `P-8`"*. **Dua hal diubah, dan keduanya
punya sebabnya.**

**Koreksi 1 — `NILAI_BATAS` bukan golongan C; ia golongan B.** §10.18 menulis asalnya sendiri:

> `NILAI_BATAS` | `Earthquake` / `FloodJab` / `FloodNation` / `RSMDLimit` **+ mata uangnya**

Mata uangnya **ada** di sistem lama — `CurrencyEarthquake`, `CurrencyFloodJab`, `CurrencyFloodNat`,
`CurrencyRSMD` — dan keempatnya tercatat di §3.3 `COCOK-SILANG-CACAH-ATRIBUT.md` sebagai *"melebur
ke paket uangnya"*. Satu baris `BATAS_PER_BAHAYA` memuat **satu** paket uang, jadi bentuknya persis
golongan B. **Ia dipindahkan, dan golongan C tinggal LIMA.**

Satu beda dari kelima kolom B yang lain, dan sebabnya ditulis: `KODE_MATA_UANG` di sini **boleh
kosong**, sebab ia **bukan bagian kunci alami** mana pun — kunci `BATAS_PER_BAHAYA` adalah
`ID_BAHAYA` (INV-14). Memaksanya `NOT NULL` berarti migrasi gagal pada baris warisan yang mata
uangnya kosong.

**Koreksi 2 — namanya menyebut paketnya, bukan `KODE_MATA_UANG` polos.** Golongan B dapat memakai
nama polos karena **satu baris memuat tepat satu paket uang**. Kelima kolom C tidak begitu:

| Entitas | Paket golongan C di baris yang sama | Kenapa nama polos tidak dapat dipakai |
|---|---|---|
| `VERSI_KONTRAK` | **tiga** — `BATAS_MAKSIMUM_KELOMPOK`, `BATAS_MAKSIMUM_NON_KELOMPOK`, `BATAS_PILIHAN` | satu kolom mata uang untuk tiga jumlah berarti **memutuskan bahwa ketiganya selalu semata uang** — fakta yang tidak ada di mana pun. Dan baris ini **sudah** punya `KODE_MATA_UANG_KONTRAK` |
| `LAYER` | **dua** — `DEDUCTIBLE`, `LIMIT_AGREGAT` | baris ini **sudah** punya `KODE_MATA_UANG` milik `LIMIT` (P-8 golongan B) |

**Nama yang berarti dua hal di sistem lama tidak dibawa** — dan nama yang akan berarti dua hal di
sistem baru tidak ditanam. Maka:

| Entitas | Paket uang | Kolom mata uang | Bita |
|---|---|---|---:|
| `VERSI_KONTRAK` | `BATAS_MAKSIMUM_KELOMPOK` | `MATA_UANG_BATAS_KELOMPOK` | 24 |
| `VERSI_KONTRAK` | `BATAS_MAKSIMUM_NON_KELOMPOK` | `MATA_UANG_BATAS_NON_KELOMPOK` | 28 |
| `VERSI_KONTRAK` | `BATAS_PILIHAN` | `MATA_UANG_BATAS_PILIHAN` | 23 |
| `LAYER` | `DEDUCTIBLE` | `MATA_UANG_DEDUCTIBLE` | 20 |
| `LAYER` | `LIMIT_AGREGAT` | `MATA_UANG_LIMIT_AGREGAT` | 23 |

Kelimanya **boleh kosong**, dan tidak satu pun melewati batas 30 bita — diperiksa perkakasnya, yang
**berhenti** bila ada yang melewatinya.

> **Satu ketaksamaan yang tersisa, dan ia dilaporkan bukan dirapikan.** `LAYER.KODE_MATA_UANG`
> (golongan B, milik `LIMIT`) memakai nama polos sementara kedua tetangganya bernama. Itu keputusan
> `P-8` yang **sudah mendarat**, dan mengubahnya bukan wewenang putaran ini. Bila pemilik proses
> menghendaki keseragaman, sesi DDL menamainya `MATA_UANG_LIMIT` — **sebelum data dimuat**, satu
> penggantian nama tanpa kehilangan apa pun.

### Dasar

Kolom *nullable* yang tidak terpakai **murah** — ia kolom kosong. Menambahkannya **sesudah migrasi
berjalan** mahal: setiap baris warisan harus disentuh ulang, dan **tidak ada sumber untuk
mengisinya**, sebab sistem lama memang tidak merekam mata uangnya. Polanya sendiri **sudah
diputuskan `P-8`**; yang ditambahkan di sini hanya penerapannya pada golongan yang waktu itu
ditahan.

### Syarat pembalikan

**`T-6` menyempitkan, tidak lagi memblokir.** Bila teknik treaty menjawab bahwa kelima besaran itu
memang selalu dalam mata uang kontrak, kolomnya **dicabut sebelum data dimuat** dan bacaannya
diserahkan ke `KODE_MATA_UANG_KONTRAK`. Mencabut kolom kosong tidak berongkos.

### Arah dampak bila salah

Salah ke arah longgar: lima kolom kosong selamanya. Salah ke arah sebaliknya: sebuah batas kapasitas
dalam mata uang yang berbeda dari mata uang kontrak **tidak dapat dinyatakan**, dan yang terjadi
adalah orang menuliskan angkanya dalam mata uang kontrak setelah mengalikannya sendiri — **jumlah
uang yang sudah dikonversi diam-diam, tanpa kurs yang tercatat.**

---

---

## KTV-D — `PERAN_PELAKU` pada jejak perubahan adalah **POTRET**, bukan rujukan

**Ditambahkan 24 September 2026, putaran penulisan tiket batch 1. Terverifikasi: NOL.**
**Label: BARU.** Entitas `JEJAK_PERUBAHAN` seluruhnya baru; sistem lama tidak punya jejak perubahan.

### Kenapa ia diputuskan sekarang, bukan ditunggu

Irisan `39` semula ditandai **MENAHAN** atas `F-16` — ketiadaan entitas peran dan penugasan bertanggal
di §10, padahal `ADR-0044` menempatkan keduanya di gelombang 1. Pemilik proses mengoreksi
penandaannya, dan koreksinya tajam:

> **Yang menahan irisan `39` bukan ketiadaan entitas peran.** Yang menahannya satu keputusan yang
> belum diambil: **`PERAN_PELAKU` itu potret, atau rujukan?**

Dan keputusan itu **tidak dapat ditunda sampai entitas perannya lahir**, sebab jejak perubahan
bersifat **tambah-saja** dan `ADR-0045` menyebutnya **fakta mesin**:

> **Baris pertama yang ditulis dengan bentuk yang salah tidak dapat diperbaiki tanpa menulis ulang
> sejarah** — hal yang jejak ini ada untuk mencegahnya.

Ini bentuk ketidaksetangkupan yang sama dengan `KTV-A`, hanya lebih tajam: di sana yang hilang
**data yang terpotong**; di sini yang hilang **kemampuan membuktikan apa pun**, sebab jejak yang
pernah ditulis ulang berhenti menjadi bukti.

### Keputusan

> **`PERAN_PELAKU` adalah POTRET.** Ia bertipe teks, **final**, dan **tidak dinormalisasi ulang**
> ketika entitas peran kelak lahir.

Ia sekeluarga dengan `PELAKU` pada entitas yang sama dan dengan `NAMA_PEMUTUS` pada
`CATATAN_PERSETUJUAN` — ketiganya dibekukan sebagai teks **dengan sengaja**, dan §12.5 sudah
menyatakannya untuk yang ketiga.

### Dasar

`ADR-0045` menuntut jejak menyatakan peran yang berlaku **saat itu**, *"bukan peran orangnya hari
ini"*. Rujukan ke tabel peran **tidak dapat memenuhinya**:

| Bentuk | Yang terjadi ketika peran seseorang diganti |
|---|---|
| **rujukan** ke entitas peran | jejaknya **ikut berubah** — baris tahun lalu tiba-tiba berbunyi bahwa orang itu sudah direktur waktu itu |
| **potret** teks | jejaknya **tetap**, dan ia menyatakan apa yang benar pada saat perubahan terjadi |

Bentuk pertama **menghapus persis fakta yang jejak ini ada untuk merekamnya**. Dan ia melakukannya
**tanpa satu galat pun** — tidak ada yang tahu bahwa sejarah sudah berubah.

### Arah dampak bila salah

Salah ke arah potret: sebuah kolom teks yang nilainya dapat dihitung ulang dari entitas peran, dan
**ejaannya dapat tidak seragam** antar baris — misalnya `"Direktur"` berdampingan dengan
`"DIREKTUR"`. Itu **menyulitkan pengelompokan**, dan itu seluruh ongkosnya.

Salah ke arah rujukan: **jejak berhenti membuktikan apa pun**, dan kerusakannya **tidak dapat
dipulihkan** — nilai lama yang benar tidak tersimpan di mana pun untuk dikembalikan.

### Syarat pembalikan

> Bila bisnis menuntut jejak menampilkan **peran terkini** alih-alih peran saat itu, bentuknya
> ditinjau — dan itu **PERUBAHAN ARTI, bukan perubahan bentuk.**

Kalimat terakhir yang mengikat. Peninjauannya **tidak boleh** dikerjakan sebagai penyeragaman skema:
ia menuntut pernyataan tertulis bahwa jejak modul ini **berhenti menjadi bukti tentang masa lalu** dan
menjadi tampilan tentang keadaan sekarang. Siapa pun yang mengubahnya tanpa pernyataan itu
membalikkan `ADR-0045` tanpa membukanya.

### Hubungannya dengan `F-16`, dan kenapa `F-16` tetap terbuka

`F-16` — §10 tidak memuat entitas peran maupun penugasan bertanggal, padahal `ADR-0044` menempatkan
keduanya di gelombang 1 — **tetap lubang terbuka dan tetap berpemilik**.

Yang berubah hanya **kedudukannya**: ia **tidak lagi menahan** irisan `39`. Bentuk `PERAN_PELAKU`
sudah diputuskan, dan bentuk itu **tidak berubah** apa pun jadinya entitas peran nanti — sebab
potret memang tidak merujuk apa-apa.

> **Yang masih ditagih `F-16`:** dari mana nilai `PERAN_PELAKU` **diambil** pada saat jejak ditulis.
> Tanpa entitas peran, ia diambil dari mekanisme wewenang apa pun yang berlaku saat itu, dan
> **mekanisme itu belum ada**. Irisan `39` menyatakannya sebagai asumsi, bukan menyelesaikannya.

| Butir | Kedudukan | Yang menyempitkannya | Kapan masih dapat dibalik |
|---|---|---|---|
| `KTV-D` | diputuskan | pernyataan bisnis tentang **arti** jejak | **sebelum baris jejak pertama ditulis** — sesudah itu, membalikkannya menuntut menulis ulang sejarah |

**Terverifikasi: tidak.**

---

## Ringkasan kedudukan — §7

| Butir | Kedudukan | Yang menyempitkannya | **Kapan masih dapat dibalik** |
|---|---|---|---|
| `KTV-A` | diputuskan | gerbang sesi DDL; `Uji AP`, `Uji AQ` | **sebelum data dimuat** |
| `KTV-B` | diputuskan | induk **ketiga**, bila kelak ada | sebelum data dimuat; sesudahnya masih mungkin tetapi menuntut pemindahan baris |
| `KTV-C` | diputuskan | `T-6` | **sebelum data dimuat** |
| `KTV-D` | diputuskan | pernyataan bisnis tentang **arti** jejak | **sebelum baris jejak pertama ditulis** |

**Nol di antaranya terverifikasi.**

> **`KTV-D` ditambahkan belakangan**, pada putaran penulisan tiket, dan tenggatnya **berbeda jenis**
> dari ketiga yang lain. `KTV-A`, `KTV-B`, dan `KTV-C` bertenggat pada **pemuatan data pertama** —
> yaitu tiket `44`. `KTV-D` bertenggat pada **baris jejak pertama**, yang datang lebih awal: begitu
> irisan `39` menyala, jejaknya mulai terisi, dan sesudah itu membalikkannya menuntut menulis ulang
> sejarah.
