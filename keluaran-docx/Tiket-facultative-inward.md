# Tiket - Fakultatif Masuk dan Keluar

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| NB Fac In | **16** | 16 | 0 | 0 | 0 | 0 |
| RNW Fac In | **8** | 8 | 0 | 0 | 0 | 0 |
| Endorsment Fac In | **22** | 0 | 22 | 0 | 0 | 0 |
| Fac Out | **14** | 0 | 14 | 0 | 0 | 0 |
| **Jumlah** | **60** | **24** | **36** | **0** | **0** | **0** |

### Tiket tertahan dan gerbangnya - 36 tiket

| Modul | # | Judul | Tertahan oleh |
| --- | ---: | --- | --- |
| Endorsment Fac In | E01 | Registry predikat EDM — sumber data per-rule | `..\08-registry-rules-eval.md` |
| Endorsment Fac In | E02 | Tiga predikat yang perilakunya berbeda di endorsement | E01 |
| Endorsment Fac In | E03 | Tracer — alur masuk endorsement sampai kasus berfase Policy | `..\08-registry-rules-eval.md` · E01 |
| Endorsment Fac In | E04 | Enam gerbang penolakan dan empat klep pembatalnya | E03 |
| Endorsment Fac In | E05 | Validasi tanggal endorsement | E03 |
| Endorsment Fac In | E06 | Tracer — before-image tiga lapis untuk satu lini | `..\01-tracer-money-ratio-premi-pa.md` · E01 |
| Endorsment Fac In | E07 | Lapis A — sembilan penyalinan yang menyimpang dari pola | E06 |
| Endorsment Fac In | E08 | Lapis B — nilai lama per baris untuk enam lini sisanya | E06 |
| Endorsment Fac In | E09 | Dua guard yang dipertahankan dan satu perbaikan sadar | E08 |
| Endorsment Fac In | E10 | Lapis C — penanda baris warisan, tujuh varian | E06 |
| Endorsment Fac In | E11 | Tiga gerbang keluar lapis B | E08 |
| Endorsment Fac In | E12 | Prefactor — Seam 3 menjadi satu pintu dengan empat bentuk | `..\01-tracer-money-ratio-premi-pa.md` · `..\03-resolver-cob-skala.md` · |
| Endorsment Fac In | E13 | Selisih pembayaran tingkat mata uang | E12 · E06 |
| Endorsment Fac In | E14 | Tabel rumus premi-menjadi per lini, jenis endorsement, dan jalur | E13 |
| Endorsment Fac In | E15 | Selisih per baris spreading dan pemetaan kolomnya | E12 · E06 |
| Endorsment Fac In | E16 | Berkas pendukung layar endorsement | E01 · E03 |
| Endorsment Fac In | E17 | Lapisan query dan lookup endorsement ⛔ BLOCKED | ⛔ **seam repository** (belum ada) |
| Endorsment Fac In | E18 | Jalur Life sebagai alur tersendiri | E06 · E08 · E10 |
| Endorsment Fac In | E19 | Perhitungan premi jiwa | E12 · E18 |
| Endorsment Fac In | E20 | Skoring medis — Seam 6 | E18 |
| Endorsment Fac In | E21 | Jalur produksi endorsement ⛔ BLOCKED | ⛔ **tabel flat (P-10)** · E13 · E15 |
| Endorsment Fac In | E22 | Rekonsiliasi eksak siklus endorsement | `..\15-de-identifikasi-berkas-kasus.md` · `..\16-rekonsiliasi-eksak-tahap-1.md`  |
| Fac Out | F01 | Predikat Fac Out masuk registry | `..\08-registry-rules-eval.md` · `..\09-predikat-sikap-khusus.md` |
| Fac Out | F02 | Tracer — derivasi Fac Retro saat UW Accept (lini Fire) | F01 · `..\01-tracer-money-ratio-premi-pa.md` · `..\11-tangga-akseptasi-bentuk-a. |
| Fac Out | F03 | Daftar pengecualian okupasi — 26 alternatif | F02 |
| Fac Out | F04 | Derivasi Fac Retro untuk tiga lini bisnis sisanya | F02 |
| Fac Out | F05 | Empat struktur penampung Fac Retro — dipisahkan, bukan disatukan | F02 |
| Fac Out | F06 | Perhitungan premi retrosesi | F05 · `..\01-tracer-money-ratio-premi-pa.md` · `..\03-resolver-cob-skala.md` |
| Fac Out | F07 | Validasi Fac Out sebelum penawaran dikirim | F05 |
| Fac Out | F08 | Tangga persetujuan retrosesi — empat tingkat | F05 |
| Fac Out | F09 | Dua generator nomor retrosesi — tidak disatukan | F08 |
| Fac Out | F10 | Produksi Fac Out ke tabel `FACOUTPRODUCTION` | F06 · F08 · F09 |
| Fac Out | F11 | Jalur Fac Out pada siklus endorsement | F10 · `..\edm\E21-jalur-produksi-edm.md` |
| Fac Out | F12 | Layar dan daftar Fac Out | F01 · F05 |
| Fac Out | F13 | Cetak R/I Slip dan kirim surat retrosesi ⛔ BLOCKED EKSTERNAL | F08 · ⛔ **eksternal — isi `M_LINK_SERVICE`** |
| Fac Out | F14 | Rekonsiliasi eksak Fac Out | F10 · `..\15-de-identifikasi-berkas-kasus.md` · `..\16-rekonsiliasi-eksak-tahap- |

---

# NB Fac In

Jumlah tiket: **16**

## NB Fac In - 01 - Tracer — uang, rasio, dan satu jalur premi PA yang hidup ujung ke ujung

**What to build:** Satu kasus PA sederhana masuk, satu angka premi keluar — lewat tipe uang dan rasio
yang sesungguhnya, bukan angka telanjang. Ini irisan **tertipis yang lengkap**: begitu ia hijau,
seluruh rantai (parsing masukan → uang bermata-uang → rasio berskala → rumus → pembulatan → hasil)
terbukti hidup, dan tiket berikutnya tinggal menambah bentuk rumus.

Uang tidak pernah `float`. Rasio membawa satuannya sendiri. Keduanya **tipe berbeda yang tidak dapat
dijumlahkan** — satu-satunya jembatan adalah perkalian eksplisit `uang × rasio → uang`. Nilai masuk
sebagai teks berkoma desimal dan dikonversi **di batas input**, bukan tersebar di dalam perhitungan.

Bentuk rumus yang dipakai tracer ini adalah yang paling sederhana dari PA — pembagi 1.000 langsung,
tanpa pro-rata: `CalculatePremiPA_FacIn` **L1003**. Bentuk PA lainnya menyusul di tiket 04.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Nilai uang memakai tipe desimal dengan **mata uang wajib menyertai**; tidak ada `float` di jalur uang mana pun
- [ ] Rasio memakai tipe terpisah yang **membawa skalanya sendiri** (per mille / persen)
- [ ] Menjumlahkan uang dengan rasio **gagal saat kompilasi** — dibuktikan berkas uji yang wajib tidak terkompilasi, bukan unit test runtime
- [ ] Satu-satunya jembatan uang↔rasio adalah operasi perkalian eksplisit
- [ ] Parser masukan menerima **koma sebagai pemisah desimal** dan memangkas spasi di ujung; konvensi ditetapkan eksplisit, **tidak diserahkan ke locale** (K-027)
- [ ] Resolver lini bisnis → skala berisi **PA saja** pada tiket ini
- [ ] Pintu masuk perhitungan premi menerima satu kasus PA dan mengembalikan uang
- [ ] **Rekonsiliasi eksak**: fixture PA bentuk `L1003` cocok dengan sistem lama **sampai digit terakhir**, tanpa toleransi (ADR-F-0001)
- [ ] Pembulatan menuliskan presisinya di tempatnya, disertai komentar yang menyebut rule Pega asal dan nomor langkahnya

## NB Fac In - 02 - Mata uang yang tidak diketahui sebagai keadaan eksplisit

**What to build:** Sistem baru dapat membaca dan menampilkan nilai uang yang **tidak membawa mata
uang**, sebagaimana sistem lama melakukannya hari ini — tanpa menebak, dan tanpa berhenti.

Dasarnya terukur: korpus memang kehilangan informasi ini (112 layar menampilkan uang tanpa field mata
uang mana pun), dan pengukuran produksi D1 menemukan **4 dari 691.925 baris** tiba tanpa mata uang.
Langka, tetapi nyata — dan justru kelangkaannya yang berbahaya: cacat yang muncul 4 kali dari 691.925
**tidak akan tertangkap pengujian sampel**.

Karena itu keadaan tidak-diketahui dibuat **eksplisit dan terlihat**, lalu gagal keras **hanya di
titik yang benar-benar memerlukan jawabannya**. Default diam-diam ke satu mata uang ditolak: itu
menebak, dan tebakannya akan salah persis pada kasus paling mahal.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Keadaan mata uang tidak diketahui adalah **nilai sah yang eksplisit**, bukan kosong dan bukan kegagalan
- [ ] Baca, tampilkan, dan simpan di memori: **diizinkan**
- [ ] Aritmetika antar nilai yang sama-sama tidak diketahui mata uangnya: **diizinkan**
- [ ] Aritmetika **lintas mata uang** yang melibatkannya: **gagal keras**
- [ ] **Tidak ada** default diam-diam ke mata uang mana pun di kode
- [ ] Implementasi menghasilkan **hitungan** berapa banyak nilai tiba tanpa mata uang — sebagai ukuran, bukan dugaan
- [ ] Kewajiban **"gagal keras saat menulis ke Oracle"** dicatat sebagai TODO yang terikat pada tiket jalur tulis produksi; **tidak** diimplementasikan sekarang karena seam-nya belum ada, dan **tidak** dihapus dari daftar kewajiban

## NB Fac In - 03 - Resolver lini bisnis → skala rasio, dan aturan pembagi komposit

**What to build:** Setiap rasio yang masuk perhitungan mendapat skalanya dari **satu resolver
terpusat** yang memetakan lini bisnis (COB) → per mille atau persen. Tidak ada tempat kedua yang bisa
menyimpang.

Properti rate yang sama memakai **dua satuan berbeda menurut lini bisnis**. Menyeragamkannya
menggeser premi **satu ordo besaran 10** pada separuh portofolio — jadi bahayanya bukan salah memilih
satuan, melainkan menyamakannya.

**Aturan penguraian pembagi komposit — inti tiket ini.** Pembagi gabungan pada satu operasi pembagian
adalah **hasil kali** sumbangan tiap faktor bersatuan, bukan satu satuan tunggal. Satuan rate dibaca
dari **faktor yang menempel padanya saja**, tidak pernah dari pembagi total:

```
@Math.divide((.TSI * .Rate * ProRatePercent), 100000, 4)
        100000  =  1000 (rate ber-‰)  ×  100 (ProRatePercent ber-%)
```

Pembagi `100000` karena itu **menegaskan** rate = per mille — bukan membantahnya. Siapa pun yang
membacanya sebagai "satuan lebih kecil dari per mille" akan salah satu ordo besaran.

Resolver mengikuti **rumus**, bukan label layar: tiga label layar terbukti bertentangan dengan
rumusnya, dan rumus yang menang di ketiganya.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Peta skala lengkap: **PA · Layering · FIRE = per mille**; **MBU · ANEKA · BONDING · GOLF · MARINE CARGO = persen**
- [ ] Aturan pembagi komposit diterapkan: sumbangan per faktor, bukan pembagi total
- [ ] Kasus uji membuktikan pembagi `100000` menghasilkan **per mille**, bukan satuan lain
- [ ] **Lini bisnis yang belum ada di peta → gagal keras**, bukan skala default
- [ ] Skala **tidak** dibaca dari tabel konfigurasi yang dapat diubah tanpa deployment — ia fakta struktural korpus, bukan parameter bisnis
- [ ] Resolver adalah **satu-satunya** pengisi skala; tidak ada pemanggil yang menyetel skala sendiri

## NB Fac In - 04 - Rumus premi PA — keempat bentuk terverifikasi

**What to build:** Premi PA dihitung benar untuk **keempat bentuk rumus** yang benar-benar ada di
sistem lama, bukan satu bentuk yang mewakili. Keempatnya hidup berdampingan di satu rule dan sepakat
bahwa rate PA berskala **per mille** — dua memberi pembagi 1.000 langsung, dua memberi 100.000 yang
terurai menjadi 1.000 × 100.

Asal, `CalculatePremiPA_FacIn`:

```
L713   (@Math.divide((.TSI * .Rate),100000,4) * ProRatePercent) - .Discount
L858   (@Math.divide((.TSI*.Rate),1000,4) * @Math.divide((.PctShortPeriod),100,4)) - .Discount
L1003  (@Math.divide((.TSI*.Rate),1000,4)*1) - .Discount
L1146  @Math.divide((.TSI*ProRatePercent*.Rate),100000,20) - .Discount
```

⚠️ **Presisi berbeda di dalam rule yang sama** — `L1146` membulatkan **20 desimal**, tiga lainnya
**4 desimal**. Itu bukan kekeliruan yang perlu dirapikan; itu perilaku terekam. Presisi ditulis
**literal di tempatnya**, bukan disentralkan ke registry: satu rule tidak punya satu presisi.

**Urutan operasi dipertahankan persis** — `round(a,4) * b`, bukan `round(a*b,4)`. Ini yang menentukan
digit terakhir, dan digit terakhir adalah ukuran keberhasilan.

**Blocked by:** 03

**Status:** ready-for-agent

- [ ] Keempat bentuk diimplementasikan sebagai cabang internal di balik satu pintu masuk perhitungan — **bukan** empat seam
- [ ] Presisi tiap pembulatan ditulis **literal di tempatnya**, disertai komentar yang menyebut rule asal dan nomor langkahnya
- [ ] Presisi **20 desimal** pada bentuk `L1146` direproduksi apa adanya, tidak diseragamkan ke 4
- [ ] Urutan operasi tiap bentuk identik dengan sumbernya
- [ ] **Rekonsiliasi eksak** per bentuk: keempat fixture cocok sampai digit terakhir, tanpa toleransi
- [ ] Pengurangan diskon terjadi di posisi yang sama seperti sumbernya

## NB Fac In - 05 - Rumus premi MBU — bentuk bersarang

**What to build:** Premi MBU dihitung benar, dan bersamanya **aturan pembagi komposit terbukti dari
struktur ekspresinya sendiri** — bukan dari tafsir kita.

MBU memakai bentuk **bersarang**, dan di situlah nilainya. Pembagi **dalam** menempel pada rate;
pembagi **luar** menempel pada pro-rata. Keterikatan faktor→pembagi karena itu **tertulis eksplisit**,
tidak perlu disimpulkan:

```
@Math.divide( @Math.divide((.TSI*(.Rate+.Loading)),100,4) * @if(.ProRatePercent=="",100,.ProRatePercent), 100, 4)
                                          ^^^                                                            ^^^
                              pembagi DALAM → .Rate (%)                        pembagi LUAR → ProRatePercent
```

Asal: `FillPremiMBU_FacIn` **L1144** dan **L1626**.

Dua hal yang mudah terlewat: `.Loading` **dijumlahkan ke rate sebelum dikalikan** (bukan sesudah),
dan pro-rata kosong diperlakukan sebagai **100**, bukan nol — membalik keduanya mengubah angka tanpa
gejala.

**Blocked by:** 03

**Status:** ready-for-agent

- [ ] Kedua bentuk (`L1144`, `L1626`) diimplementasikan sebagai cabang internal di balik pintu masuk yang sama
- [ ] Struktur **bersarang** dipertahankan — bukan diratakan menjadi satu pembagian dengan pembagi 10.000
- [ ] `.Loading` dijumlahkan ke rate **sebelum** perkalian, sesuai sumbernya
- [ ] Pro-rata kosong → **100**, direproduksi apa adanya
- [ ] Skala MBU **persen** diambil dari resolver (tiket 03), bukan ditulis ulang di sini
- [ ] **Rekonsiliasi eksak** per bentuk, tanpa toleransi
- [ ] Kasus uji yang membuktikan pembagi dalam terikat ke rate dan pembagi luar ke pro-rata — bila keduanya tertukar, uji harus **gagal**

## NB Fac In - 06 - Rumus premi Layering

**What to build:** Premi untuk struktur berlapis dihitung benar, dengan rate berskala **per mille**.

Asal: `GenerateLayerList_ACT` **L909**:

```
@Math.divide((.TSI * .Rate * pyWorkPage.OfferFacIn.ProRatePercent), 100000, 4)
```

Bentuknya identik dengan salah satu varian PA, dan terurai sama: `100000 = 1.000 × 100`.

⚠️ **Ini satu-satunya ekspresi ber-rate di seluruh berkas bernuansa Layer di korpus NB** — tidak ada
varian yang bertentangan. Kepastian itu penting justru karena label layar daftar Layer menulis
`Rate (%)`, yang **bertentangan dengan rumusnya**. Rumus yang menang; labelnya yang diperbaiki, dan
perbaikan label **tidak mengubah satu angka pun**.

**Blocked by:** 03

**Status:** ready-for-agent

- [ ] Rumus diimplementasikan sebagai cabang internal di balik pintu masuk perhitungan yang sama
- [ ] Skala **per mille** diambil dari resolver (tiket 03)
- [ ] Pembagi `100000` diuraikan sebagai `1.000 × 100`, bukan diperlakukan sebagai satu satuan
- [ ] **Rekonsiliasi eksak**, tanpa toleransi
- [ ] Label tampilan daftar Layer diperbaiki dari `(%)` menjadi `(‰)` — **perbaikan label saja**, dicatat bahwa nol angka berubah sehingga rekonsiliasi tidak terganggu

## NB Fac In - 07 - Pembulatan di dalam loop akumulasi — kasus uji wajib

**What to build:** Bukti bahwa perhitungan premi tetap cocok **untuk daftar panjang**, bukan hanya
untuk satu nilai.

Ini tiket pengujian, dan ia berdiri sendiri karena alasan yang spesifik: di sistem lama, sebagian
pembulatan berada **di dalam loop akumulasi**, sehingga galatnya **menumpuk per iterasi**. Port yang
benar untuk satu nilai masih bisa meleset untuk daftar panjang — dan itu bentuk kesalahan yang
**paling sulit terlihat**, karena setiap nilai tunggalnya lulus.

Lingkupnya **lintas lini bisnis**: galat akumulasi tidak mengenal batas COB, jadi menguji satu COB
saja tidak membuktikan apa pun tentang yang lain.

Dua sifat sistem lama yang memperkuat kebutuhan ini: validasi sisa spreading menuntut kesamaan
**persis** (jumlah bagian = 100, jumlah premi tersebar = premi) tanpa distribusi galat pembulatan,
dan pembulatan dilakukan lebih dulu sebelum dibandingkan.

**Blocked by:** 04, 05, 06

**Status:** ready-for-agent

- [ ] Kasus uji akumulasi untuk **tiap lini bisnis** yang punya rumus di tiket 04–06
- [ ] Setiap kasus memakai daftar yang cukup panjang untuk memunculkan penumpukan galat — bukan dua atau tiga baris
- [ ] Pembulatan terjadi **di dalam** loop, di posisi yang sama seperti sumbernya — bukan sekali di akhir
- [ ] **Rekonsiliasi eksak** pada nilai agregat, bukan hanya pada tiap elemen
- [ ] Kasus uji yang membuktikan pemindahan pembulatan ke luar loop **membuat uji gagal** — kalau tidak, ujinya tidak menguji apa pun
- [ ] Validasi sisa direproduksi tanpa toleransi: kesamaan **persis** setelah pembulatan, tanpa distribusi galat

## NB Fac In - 08 - Registry predikat — satu seam untuk seluruh rule `When`

**What to build:** Seluruh gerbang keputusan sistem lama hidup di **satu registry** yang dapat
ditanyai dengan nama, dan setiap predikat menyebut rule Pega asalnya. Tidak ada salinan logika
gerbang yang tersebar dan bisa menyimpang.

⚠️ **Kondisi rule `When` ada di DUA tag, dan keduanya wajib dibaca.** Membaca hanya yang pertama
membuat rule tampak "tanpa kondisi" padahal kondisinya ada — **170 berkas (28,3 %)** menyembunyikan
kondisinya di tag kedua: 13 tagnya kosong, 157 berisi teks placeholder. Seluruh **601** berkas punya
tag kedua terisi, dan **tidak ada satu pun** rule yang kondisinya benar-benar tidak terbaca.

Klaim lama "lima rule `When` kondisinya kosong" berasal dari membaca satu tag saja, dan **keliru**.

Predikat individual **bukan** seam. Satu pintu masuk, diuji table-driven.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Satu pintu masuk evaluasi bernama, melayani seluruh registry — bukan satu fungsi publik per predikat
- [ ] Kondisi dibaca dari **kedua** tag; rule yang kondisinya hanya ada di tag kedua tetap terbaca benar
- [ ] Teks placeholder pada tag pertama **tidak** diperlakukan sebagai kondisi
- [ ] **Nama predikat yang tidak dikenal → gagal keras**, bukan mengembalikan salah — salah ketik tidak boleh berubah menjadi gerbang yang selalu tertutup
- [ ] Setiap predikat menyebut rule Pega asalnya dalam komentar
- [ ] Uji table-driven mencakup contoh dari kedua kelompok: yang kondisinya terbaca di tag pertama, dan yang tersembunyi di tag kedua

## NB Fac In - 09 - Empat predikat yang sikapnya sudah diputuskan, dan tidak boleh disamakan

**What to build:** Empat predikat yang masing-masing sudah punya keputusan tersendiri, dijalankan
persis sebagaimana diputuskan — bukan diseragamkan menjadi "implementasikan saja kondisinya".

**`IsPKSASM` → gagal keras.** Kondisinya terbaca, tetapi **labelnya belum ter-resolve**, tag kondisi
pertamanya masih berisi placeholder, dan ia bertanda sementara — di ketiga folder korpus. Ia gagal
keras bukan karena kondisinya hilang, melainkan karena **belum terbukti sebagai yang dieksekusi**.

⛔ Daftar nilai sah untuk properti yang diujinya sudah diketahui dari basis data. **Itu tidak
mencabut gagal-kerasnya.** Mengetahui nilai apa yang sah dan membuktikan kondisi itu dieksekusi
adalah dua pertanyaan berbeda; hanya yang kedua yang jadi penyebab, dan yang kedua belum terjawab.

**`IsOfferFacIn` → ekspresi tersimpan yang berlaku.** Rule ini memuat **dua kondisi berbeda**: teks
tampilannya menyebut daftar kode bisnis, ekspresi tersimpannya menguji penanda fakultatif. Yang
berlaku adalah **ekspresi tersimpan**. Teks tampilan dicatat sebagai kandidat perbaikan, **tidak**
diimplementasikan — dan menjadi tersangka pertama bila paralel run memperlihatkan selisih pada
gerbang masuk siklus.

**`IsSpreadingDepan` → implementasikan, catat asal salinannya.** Rule ini hilang dari folder NB tetapi
terbaca penuh di folder siklus lain. Ia **bukan** rule yang kondisinya tidak diketahui, jadi tidak
perlu gagal keras. ⚠️ Salinan yang terbaca berversi lebih lama, dan 95 rule di korpus terbukti berbeda
versi antar folder — asal salinan **wajib dicatat** di komentar.

**`IsFacout` → diport apa adanya, nama dipertahankan.** Fiturnya usang secara bisnis, tetapi kodenya
masih aktif dan masih dirujuk sepuluh berkas. ⚠️ Namanya menyiratkan *fac out*; **isinya menguji
hasil keputusan Banding**. Nama dipertahankan demi ketertelusuran ke rule asal — **tetapi namanya
bukan dokumentasi artinya**, dan komentar wajib menyebutkan itu.

**Blocked by:** 08

**Status:** ready-for-agent

- [ ] `IsPKSASM` **gagal keras** bila dievaluasi; pesannya menyebut alasannya (kondisi belum terbukti dieksekusi), bukan "nilai tidak dikenal"
- [ ] `IsOfferFacIn` memakai ekspresi tersimpan; teks tampilan dicatat sebagai kandidat perbaikan di komentar, tidak dieksekusi
- [ ] `IsSpreadingDepan` diimplementasikan sesuai kondisi terbacanya, dengan komentar menyebut **folder asal salinan dan versinya**
- [ ] `IsFacout` diport apa adanya dengan nama aslinya; komentar menyatakan isinya menguji **Banding**, bukan fac out
- [ ] Keempatnya punya kasus uji sendiri — termasuk uji bahwa `IsPKSASM` **memang** gagal keras

## NB Fac In - 10 - Jembatan predikat lini bisnis → resolver skala

**What to build:** Lini bisnis sebuah kasus **ditentukan oleh predikat yang sesungguhnya dipakai
sistem lama**, lalu mengalir ke resolver skala — bukan diisi pemanggil sebagai parameter yang
diasumsikan benar.

Ini titik integrasi antara registry predikat dan perhitungan premi: dua bagian yang sudah bekerja
sendiri-sendiri, disambungkan di satu tempat yang dapat diuji. Predikat yang terlibat mengenali
FIRE, PA, MBU, ANEKA, MARINE CARGO, GOLF, dan BONDING.

⚠️ **Kehati-hatian yang diwarisi dari discovery.** Lini bisnis ditentukan satu properti, **tetapi dua
predikat membacanya lewat jalur berbeda di dalam agregat** — sebagian memakai dua sampai tiga jalur.
Apakah ketiga jalur selalu sinkron **belum terjawab**. Sampai terjawab, jalur-jalur itu
**dipertahankan apa adanya** dan perbedaannya dicatat saat runtime; properti itu **tidak boleh
dinormalisasi** menjadi satu field.

⚠️ Sebagian nilai lini bisnis muncul **dengan spasi di depan atau belakang**. Dipangkas di batas
input, dan dicatat sebagai kandidat perbaikan.

**Blocked by:** 08, 03

**Status:** ready-for-agent

- [ ] Lini bisnis diturunkan lewat registry predikat, bukan diterima mentah dari pemanggil
- [ ] Hasilnya mengalir ke resolver skala; **tidak ada** jalur yang melewati resolver
- [ ] Jalur baca ganda pada properti penentu **dipertahankan**, tidak dinormalisasi
- [ ] Ketidaksinkronan antar jalur **tercatat saat runtime** alih-alih dipilih diam-diam salah satunya
- [ ] Spasi di ujung nilai dipangkas di batas input, dicatat sebagai kandidat perbaikan
- [ ] Kasus yang lini bisnisnya tidak dikenali **gagal keras** lewat resolver (tiket 03), bukan memakai default

## NB Fac In - 11 - Tangga akseptasi bentuk standar — satu keputusan, satu transisi

**What to build:** Seorang underwriter mengambil satu keputusan, dan kasus berpindah **satu langkah**
pada tangga akseptasi — persis seperti hari ini.

⛔ **Ini mesin keadaan satu langkah, bukan loop.** Jangan merancangnya sebagai proses yang menghitung
seluruh rantai approver sekaligus: bentuk itu **tidak ada di sistem lama** dan berperilaku berbeda
ketika rantai terputus di tengah.

**Tiga field state yang tidak boleh disatukan menjadi satu "status":**

| Field | Isinya |
| --- | --- |
| **antrean** | peran yang **sedang** memegang kasus (ruang nama tersendiri) |
| **kode jabatan tujuan** (`next_approver_position`) | jabatan **berikutnya** dalam tangga — ruang nama berbeda |
| **hasil keputusan** | keputusan underwriting terakhir |

⚠️ Kode jabatan tujuan disimpan di properti warisan bernama `LetterNo`, yang **tidak berisi nomor
surat**. Antrean dan kode jabatan diberi **tipe berbeda** supaya kompilator menangkap pertukarannya —
salah satu tempat sistem tipe menangkap cacat warisan.

⚠️ **Tangga selesai adalah penyelesaian normal, bukan galat.** Ketika tidak ada jabatan tujuan yang
cocok, wewenang sudah cukup dan kasus berhenti di situ.

Tabel limit **disuntikkan sebagai fixture**, dan fixture itu sekaligus **kontrak bentuk data** ke DBA.
Yang hilang dari tangga ini dulu adalah datanya, bukan logikanya — kini datanya sudah ada, dengan
ejaan jabatan **tanpa spasi** yang terbukti cocok dengan token routing di rule.

⛔ **Kolom nama dan login WAJIB dibuang** saat fixture dibangun dari berkas data — keduanya memuat
nama orang. Dibuang, bukan disamarkan.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Satu keputusan manusia menghasilkan **tepat satu** transisi; tidak ada kasus uji yang mengharapkan rantai
- [ ] Ketiga field state terpisah dan bergerak independen
- [ ] Antrean dan kode jabatan bertipe **berbeda**; menukarnya **gagal saat kompilasi**
- [ ] Tidak ada jabatan tujuan yang cocok → **tangga selesai**, ditandai penyelesaian normal, bukan galat
- [ ] Urutan eskalasi mengikuti ambang limit menaik; baris pertama yang cocok = approver berikutnya
- [ ] Fixture memuat kolom yang **dibaca query** saja — **tanpa kolom nama dan login**
- [ ] Ejaan nilai jabatan **tanpa spasi**, sesuai data nyata
- [ ] Pemeriksaan kebocoran nama dijalankan atas fixture sebelum di-commit, memakai pencocokan **batas kata** — bukan substring
- [ ] ⚠️ Lima belas tautologi pembanding limit dan empat nomor polis literal yang ada di rule alur **direproduksi apa adanya** dan ditandai kandidat perbaikan — bukan dirapikan

## NB Fac In - 12 - Akseptasi lini financial — filter ambang limit, tanpa tangga berjenjang

**What to build:** Kasus lini financial menemukan jabatan berwenangnya lewat **penyaringan ambang
limit**, bukan lewat tangga bertingkat. Lini ini **tidak bereskalasi**.

⛔ **Bentuk tabelnya berbeda, dan menyamakannya membuat tangga financial macet total.** Tabel limit
financial **tidak punya** kolom jabatan atasan, tidak punya kolom antrean, dan kolom limitnya bukan
ambang tunggal melainkan **empat ambang per jenis pertanggungan**: bond, credit CL, credit NCL, trade.

⛔ **Ejaan jabatannya PAKAI SPASI**, berbeda dari bentuk standar yang tanpa spasi. Bila kode
mencocokkan dengan token tanpa spasi, **tidak ada approver financial yang pernah ditemukan**.
**Normalisasi ejaan dilarang** — menghapus spasi agar seragam adalah perubahan perilaku, bukan
migrasi.

**Yang tidak ada** adalah eskalasi berjenjang dan antrean per jabatan. **Yang tetap ada dan wajib
diport** adalah **filter ambangnya**: tiga rule SQL menyaring dengan `WHERE <kolom_limit> <= nilai`,
masing-masing untuk satu jenis pertanggungan, tanpa penggabungan tabel dan tanpa pengurutan.
Membuang `WHERE` menghapus satu-satunya kontrol wewenang yang dimiliki lini ini.

Keluarannya karena itu **daftar jabatan yang limitnya menampung nilai** — bukan satu jabatan tujuan
berikutnya.

⚠️ Kedua mekanisme **tidak disatukan di balik satu abstraksi**. Memaksa bentuk ini ke dalam bentuk
tangga berarti mengarang langkah naik yang tidak ada.

**Blocked by:** 11

**Status:** ready-for-agent

- [ ] Satu query per jenis pertanggungan (bond, credit CL, credit NCL), memakai kolom ambangnya masing-masing
- [ ] **`WHERE` dipertahankan apa adanya**; tidak ada penyaringan tambahan di luar ambang limit
- [ ] Pencocokan jabatan memakai ejaan **berspasi**; tidak ada normalisasi
- [ ] Keluaran berupa **daftar jabatan**, bukan jabatan tujuan berikutnya
- [ ] Tidak ada langkah naik dan tidak ada token antrean untuk lini ini
- [ ] Kedua bentuk tetap dua jalur terpisah di kode; kasus uji membuktikan bentuk standar **tidak** dipakai untuk lini financial dan sebaliknya
- [ ] Fixture bentuk ini juga **tanpa kolom nama**

## NB Fac In - 13 - Domain hasil keputusan, dan satu flag fase untuk jalur yang belum terverifikasi

**What to build:** Hasil keputusan underwriting punya **enam nilai yang sah dan artinya pasti**, dan
sistem berperilaku berbeda-tetapi-terkendali pada jalur yang belum terverifikasi.

**Enam nilai, artinya terverifikasi dari dekoder literal di korpus:** Accept · Reject · Ask ·
**Banding** · Decline · Revise. Dua nilai lain milik ranah ceding, hidup di properti terpisah, dan
tidak pernah menulis ke kolom ini; satu nilai lagi nol jejak di seluruh korpus.

**`Reject` dan `Decline` tidak disatukan.** Decline adalah penolakan **final tanpa hak banding**;
Reject masih **dapat dibanding**. Efek samping Reject yang terverifikasi — mematikan konfirmasi
binding dan penerimaan slip, lalu menjadi syarat munculnya jalur banding lewat hitungan riwayat
berstatus reject — adalah **perilaku yang dikehendaki**, direproduksi apa adanya.

⚠️ **Validasi domain ini adalah perilaku BARU.** Sistem lama tidak punya validasi domain pada kolom
ini — tidak punya validasi apa pun selain wajib-isi. Karena itu ia **tidak boleh** berbentuk penolakan
diam-diam.

**Satu flag fase, bukan dua basis kode:**

| Situasi | Paralel run | Produksi |
| --- | --- | --- |
| Hasil keputusan **dikenali**, baris keputusannya belum terverifikasi | gagal keras | `decline` + catatan |
| Kondisi **tidak dikenali sama sekali** | gagal keras | gagal keras |
| Nilai **di luar domain** enam nilai | gagal keras + catatan | gagal keras di jalur tulis + catatan |

`decline` dipilih di produksi karena ia **hasil default yang terekam** — memilihnya adalah reproduksi,
bukan tebakan. Transisi sebuah jalur dari gagal-keras ke `decline` **hanya** boleh setelah jalur itu
diverifikasi terhadap ekspor produksi, bukan karena ia sering muncul dan mengganggu.

**Blocked by:** 11

**Status:** ready-for-agent

- [ ] Enum enam nilai, masing-masing menyebut rule Pega asal dekodernya dalam komentar
- [ ] Reject dan Decline **tetap terpisah**, dengan efek samping Reject direproduksi
- [ ] **Satu** flag fase eksplisit mengendalikan perbedaan perilaku — bukan dua basis kode
- [ ] Jalur "dikenali tetapi belum terverifikasi" berperilaku sesuai fase
- [ ] Jalur "tidak dikenali sama sekali" gagal keras di **kedua** fase
- [ ] Nilai di luar domain **tidak pernah ditolak diam-diam**; selalu disertai catatan
- [ ] Dicatat sebagai kandidat perbaikan bahwa validasi domain ini **tidak ada di sistem lama**
- [ ] Bila paralel run menemukan nilai lain benar-benar ada di data produksi, itu **membatalkan premis** dan dibuka sebagai keputusan baru — **bukan** alasan melonggarkan validasi diam-diam

## NB Fac In - 14 - Special Acceptance dan tangga akseptasi putaran kedua

**What to build:** Dua cabang inti tangga akseptasi yang sempat **ditangguhkan** kini berjalan penuh —
jalur **Special Acceptance** dan **tangga putaran kedua**.

Keduanya pernah ditangguhkan karena rule yang dipanggilnya tidak ada di ekspor yang kami terima.
Rule-nya kini tiba, rujukannya sudah diverifikasi ulang ke korpus, dan aturannya berlaku: **rule ada →
cabang dipulihkan**.

⚠️ **Dua rule Special Acceptance BUKAN Activity.** Keduanya dieksekusi sebagai **`RequestType` pada
langkah RDB-List**, berkelas integrasi, berdampingan dengan penanda akses dan halaman browse. Artinya
yang perlu ditulis adalah **query dan pemanggilan bacanya** — bukan sebuah fungsi layanan. Menyebutnya
"activity" akan menyesatkan porting.

Cabang putaran kedua berbeda: ia **benar-benar** panggilan activity, dipanggil dari post-activity
validasi tanggal.

⚠️ **Jalur ini hanya membaca.** Tidak ada bagian tiket ini yang menulis ke basis data; jalur tulis
produksi masih menunggu dan berada di luar lingkup.

**Blocked by:** 11, 13

**Status:** ready-for-agent

- [ ] Jalur Special Acceptance berjalan lewat pemanggilan **baca** bergaya `RequestType`, bukan dipaksa menjadi panggilan fungsi layanan
- [ ] Kedua rule Special Acceptance diperlakukan sebagai **rule integrasi/SQL**, dengan komentar menyebut asal dan kelasnya
- [ ] Tangga putaran kedua berjalan sebagai panggilan activity dari post-activity validasi tanggal
- [ ] Keduanya tetap mematuhi aturan **satu keputusan = satu transisi** (tiket 11)
- [ ] Tidak ada operasi tulis di jalur ini
- [ ] Fixture menutupi keduanya; kasus uji membuktikan cabangnya **benar-benar tercapai**, bukan sekadar ada
- [ ] Status "ditangguhkan" dicabut di dokumentasi rancangan, dengan menyebut bukti rujukan yang memulihkannya

## NB Fac In - 15 - De-identifikasi berkas kasus menjadi fixture rekonsiliasi

**What to build:** Lima berkas kasus nyata berubah menjadi fixture yang **aman disimpan di
repositori** — identitas hilang, **setiap angka utuh**.

Kelimanya berformat **JSON** dan memuat data pelanggan: nama tertanggung, alamat, serta **nama
operator** yang ikut terbawa metadata di kelimanya. Aturan proyek melarang nama orang dan data
pelanggan masuk artefak mana pun, termasuk fixture dan test — jadi berkas mentahnya **tidak boleh**
dipakai langsung oleh tiket rekonsiliasi.

Nilainya tetap besar: angka-angka di dalamnya adalah **masukan rekonsiliasi yang sesungguhnya**,
mencakup beberapa lini bisnis dan ketiga siklus.

#### ⛔ Penyaringan berbasis pola DILARANG

Ini kriteria terpenting tiket ini, dan alasannya konkret. Menyaring dengan pencocokan kata kunci akan
**membuang angka yang justru diuji**: `TotalSumInsured` tertangkap filter *"insured"*, sementara
`CedingRetention`, `ShareCeding`, dan `ShareOfCeding` tertangkap filter *"ceding"* — keempatnya
**angka yang wajib dipertahankan**.

Ini jebakan yang sama bentuknya dengan kekeliruan pencocokan awalan yang sudah dua kali terjadi di
proyek ini. Bedanya, kali ini taruhannya angka rekonsiliasi.

**Gunakan daftar kunci eksplisit, ditinjau satu per satu.**

**BUANG (12 kunci, disetujui):**
`InsuredName` · `InsuredID` · `ASMAddress` · `SelectedLocationAddress` · `RoadName` · `ASMZipCode` ·
`RiskZipCode` · `Email` · `pxCreateOpName` · `MarketingName` · `CedingCoName` · `CoinsName`

**PERTAHANKAN:** seluruh field bernilai angka — termasuk keempat yang tertangkap filter palsu di atas.

**Bila ditemukan field identitas lain saat implementasi:** tambahkan ke daftar lewat **tinjauan
manual**. **Jangan beralih ke filter pola.**

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Kelima berkas menghasilkan fixture turunan; berkas mentah **tidak** ikut masuk repositori
- [ ] Penghapusan memakai **daftar kunci eksplisit**; tidak ada pencocokan pola, substring, atau regex kata kunci di jalur penyaringan
- [ ] Kedua belas kunci pada daftar BUANG hilang seluruhnya dari fixture
- [ ] **Nol nama orang** dan **nol alamat** di fixture, diperiksa dengan pencocokan **batas kata**
- [ ] **Setiap angka identik dengan sumbernya** — dibuktikan pembandingan otomatis antara fixture dan berkas asal, bukan pemeriksaan mata
- [ ] Struktur dokumen dipertahankan; hanya nilai kunci identitas yang hilang, bukan bentuknya
- [ ] Daftar kunci yang dibuang didokumentasikan bersama fixture, sehingga penambahan berikutnya dapat ditinjau

## NB Fac In - 16 - Kerangka rekonsiliasi eksak — tahap 1

**What to build:** Sebuah pembanding yang menjalankan kasus terekam lewat perhitungan sistem baru dan
menyatakan, tanpa ruang tafsir, apakah hasilnya **identik sampai digit terakhir** dengan sistem lama.

Ini ukuran keberhasilan migrasi, bukan pelengkap. Karena seluruh jalur tulis produksi melewati
prosedur basis data yang isinya baru sebagian kami miliki, **membaca kode saja tidak dapat
membuktikan port-nya benar** — hanya membandingkan keluaran dua sistem atas masukan yang sama yang
bisa.

**Tahap 1 = perhitungan murni atas masukan terekam.** Tidak menyentuh alur, tidak menyentuh jalur
tulis, tidak menunggu ekspor produksi tunggal. Tahap berikutnya (per-modul dengan fixture tangga, lalu
end-to-end) menyusul setelah prasyaratnya tiba.

⛔ **Toleransi ditolak.** Bila urutan operasi dan presisi per-langkah direproduksi apa adanya, nilai
desimal bersifat deterministik dan hasilnya **harus** identik. Selisih sekecil apa pun berarti ada
salah-port — dan toleransi hanya menyembunyikannya, justru pada fase yang dirancang untuk
menemukannya.

⛔ **Kerangka ini tidak pernah menyentuh berkas kasus mentah.** Masukannya **hanya** fixture hasil
tiket 15. Menunjuk langsung ke berkas asal akan menarik data pelanggan ke dalam jalur pengujian.

**Blocked by:** 15, 04, 05, 06, 11

**Status:** ready-for-agent

- [ ] Pembanding membaca **fixture ter-de-identifikasi saja**; tidak ada jalur yang membuka berkas kasus mentah
- [ ] Perbandingan **eksak, tanpa toleransi**; tidak ada parameter epsilon di mana pun
- [ ] Selisih dilaporkan dengan menyebut **kasus, lini bisnis, dan langkah** tempat angkanya mulai berbeda — bukan hanya "tidak cocok"
- [ ] Mencakup lini bisnis yang rumusnya sudah ada (tiket 04–06) dan nilai dasar akseptasi (tiket 11)
- [ ] Kasus yang lini bisnisnya belum punya rumus dilaporkan sebagai **belum tercakup**, bukan lulus
- [ ] Hasilnya dapat dijalankan ulang dan deterministik
- [ ] Dicatat bahwa ini **tahap 1**, dengan prasyarat tahap berikutnya disebut eksplisit

# RNW Fac In

Jumlah tiket: **8**

## RNW Fac In - R01 - Tracer — alur masuk renewal sampai kasus diterima tangga akseptasi

**What to build:** Seorang underwriter memasukkan **No. Polis + Renewal Date + Note**, menekan **OK**,
dan sebuah kasus renewal tercipta — dengan data polis lama sudah terisi sebagai nilai awal — lalu
**diterima tangga akseptasi tanpa penyesuaian apa pun**.

Ini irisan **tertipis yang lengkap** untuk renewal. Begitu ia hijau, seluruh premisnya terbukti: bahwa
renewal hanyalah pintu masuk lain menuju mesin yang sudah ada. Tiket RNW berikutnya tinggal menambah
layar.

⚠️ **Gerbang masuknya harus ditetapkan, bukan diport.** `[terverifikasi]` predikat gerbang masuk
siklus New Business **dirujuk nol kali** di korpus renewal, dan flow renewal **tidak dirujuk berkas
mana pun** — penentunya ada di konfigurasi work type/portal yang **tidak ikut terekspor**. Karena
tidak ada perilaku terekam untuk direproduksi, sistem baru menetapkannya eksplisit sesuai alur di atas.

`[terverifikasi]` Asal, `FlowAction\Renewal_FlowAct`: section `InputRenewal`, pre-activity
`InputOfferFacInEngineer_preACT`, post-activity `SetValidateDate_PostAct`, transform
`AddToListSuggestOfferFacIn_DT`. Flow `InputRenewalFacultativeIn` memuat 75 shape dan merujuk
**19 rule `When`** — seluruhnya predikat routing yang diwarisi dari NB, tidak ada predikat gerbang-masuk
di antaranya.

Berkas gambar `halaman depan renewal.JPG` **ada** sebagai referensi layar; isinya tidak dibaca sebagai
fakta.

**Blocked by:** **NB-08** (registry predikat `rules.Eval`) · **NB-11** (`acceptance.Next`)

**Status:** ready-for-agent

- [ ] Layar masuk menerima **No. Polis**, **Renewal Date**, dan **Note**; tombol **OK** memicu pembuatan kasus
- [ ] Kasus baru tercipta dengan pembeda siklus **`StatusBusiness = 2`**
- [ ] **Data polis lama tersalin sebagai nilai awal** pada kasus baru — underwriter cukup menyunting yang berubah
- [ ] Kasus hasil **diterima `acceptance.Next` tanpa penyesuaian** — tidak ada cabang khusus renewal di tangga akseptasi
- [ ] Ke-19 predikat routing dievaluasi lewat registry NB (`rules.Eval`), **bukan** salinan predikat baru
- [ ] Gerbang masuk **ditetapkan eksplisit** dan didokumentasikan sebagai keputusan rancangan — bukan diklaim sebagai perilaku terekam
- [ ] Setiap transisi menyebut rule Pega asalnya dalam komentar
- [ ] Tidak ada modul perhitungan baru yang dibuat oleh tiket ini

## RNW Fac In - R02 - Layar periode renewal

**What to build:** Underwriter melihat dan menyunting **periode polis baru** berdampingan dengan
rujukan ke **polis lama** — nomor polis lama, tanggal renewal, serta tanggal berlaku dan berakhir —
dalam satu layar.

`[terverifikasi]` `Section\PeriodeRenewal` mengikat **22 properti**, di antaranya
`.QuotationData.OldPolicyNo` · `.QuotationData.RNWDate` · `.QuotationData.StatusBusiness` ·
`.PolicyData.StartDateTime` · `.PolicyData.EndDateTime` · `.PolicyData.OfferingDate` ·
`.IsSpecialAcceptance`. Varian `_IsUW` mengikat 18 properti.

⚠️ **Dua komponen terpisah, bukan satu komponen dua mode.** `PeriodeRenewal` dan
`PeriodeRenewal_IsUW` diport masing-masing sebagai komponen sendiri.

Ini **keputusan sadar work owner untuk siklus renewal**, dan ia **sengaja berbeda** dari rekomendasi
spec model-data NB §6 (satu komponen dengan prop `mode`). Alasannya reproduksi perilaku terekam
(`CLAUDE.md` §1): sistem lama memang punya dua section. **Bukan inkonsistensi** — perbedaannya
disengaja dan tercatat di sini.

**Blocked by:** R01

**Status:** ready-for-agent

- [ ] `PeriodeRenewal` dan `PeriodeRenewal_IsUW` diport sebagai **dua komponen terpisah**
- [ ] Nomor polis lama, tanggal renewal, dan pembeda siklus terlihat dan terikat benar
- [ ] Tanggal berlaku, berakhir, dan penawaran dapat disunting sesuai perilaku lama
- [ ] Nilai awal berasal dari salinan polis lama (R01), bukan diisi ulang manual
- [ ] Komentar menyebut section Pega asal tiap komponen
- [ ] Keputusan "dua komponen terpisah" dicatat di kode sebagai pilihan sadar, dengan rujukan ke alasannya

## RNW Fac In - R03 - Layar input renewal — wadah, bukan duplikat

**What to build:** Layar input renewal menampilkan penawaran **dengan menyisipkan komponen yang sudah
ada dari New Business** — bukan dengan menyalinnya. Underwriter melihat layar yang sama seperti yang
ia kenal, dan tim tidak memelihara dua salinan yang bisa menyimpang.

Ini kriteria yang paling mudah dilanggar tanpa sengaja, dan paling mahal bila dilanggar: menyalin
komponen NB terasa lebih cepat, lalu perbaikan pada satu sisi diam-diam tidak sampai ke sisi lain.

`[terverifikasi]` Bukti bahwa layar renewal memang wadah: `Section\InputRenewal` hanya punya
**2 properti sendiri** (`.IsShowDetail` + template), dan menyisipkan **7 section**:
`AllSummarySection` · `PeriodeRenewal` · `InputRenewalDtl` · `InputDtlObject_FacIn` ·
`InputInwardFacultativeSuggest` · `EmailSection` · `EmailSectionCeding`.
Varian `_IsUW` menyisipkan **6 section**, juga dengan 2 properti sendiri.

📌 Empat section yang disisipkan **tidak ada di folder renewal** — `InputInwardFacultativeSuggest` ·
`InwardFacIn` · `OfferFacIn_NusaRe` · `OfferFacIn_NusaRe_IsUW`. Itu **bukan kekurangan**: keempatnya
ada di New Business dan **diwarisi apa adanya** (K-040), persis seperti 1.907 berkas identik lainnya.
Bangun sekali, pakai di kedua siklus.

⚠️ **Dua komponen terpisah** — `InputRenewal` dan `InputRenewal_IsUW` masing-masing komponen sendiri.
Keputusan sadar work owner, sengaja berbeda dari rekomendasi spec model-data NB §6; alasannya
reproduksi perilaku terekam (`CLAUDE.md` §1).

**Blocked by:** R02

**Status:** ready-for-agent

- [ ] `InputRenewal` dan `InputRenewal_IsUW` diport sebagai **dua komponen terpisah**
- [ ] Keduanya **menyisipkan** komponen NB — **tidak ada komponen NB yang disalin** ke dalam kode renewal
- [ ] **Kriteria pembukti:** perubahan pada sebuah komponen NB **terlihat di layar renewal** tanpa perubahan kode renewal. Bila tidak terlihat, komponennya tersalin — dan tiket ini gagal
- [ ] Keempat section yang diwarisi dari NB dipakai dari sumber NB, bukan dibuat ulang
- [ ] Layar periode (R02) tersisip, bukan diduplikasi
- [ ] Komentar menyebut section Pega asal

## RNW Fac In - R04 - Layar detail renewal

**What to build:** Layar detail renewal menampilkan seluruh rincian penawaran — objek, coverage,
spreading, komisi, pembayaran — dengan **menyisipkan komponen New Business**, sama seperti R03, dan
menambahkan hanya field yang memang khas renewal.

Ini layar terbesar dalam lingkup renewal.

`[terverifikasi]` `Section\InputRenewalDtl` menyisipkan **14 section** dan mengikat **37 properti
sendiri**; varian `_IsUW` menyisipkan **17 section** dan mengikat 36 properti. Yang disisipkan antara
lain `CoverageList` · `CoverageSpreadingList` · `CoverageCommisionList` · `PaymentCurrencyList` ·
`OfferFacIn_NusaRe` · `SummaryCoverage_Section` · `SummarySpreading_Section` · `CallDualScoringRisk` ·
`InputDtlObject_FacIn`.

⚠️ **Tiket ini sengaja dibiarkan utuh, tidak dipecah.** Pengerja **boleh** memecah varian `_IsUW`
menjadi sub-tugas bila satu context window tidak cukup — tetapi pemecahan itu keputusan pelaksanaan,
bukan pemecahan tiket.

⚠️ **Dua komponen terpisah** — `InputRenewalDtl` dan `InputRenewalDtl_IsUW` masing-masing komponen
sendiri. Keputusan sadar work owner, sengaja berbeda dari spec model-data NB §6; alasannya reproduksi
perilaku terekam (`CLAUDE.md` §1).

##### ⏸ `InputDtlObject` — usang, tetapi diport apa adanya

`[terverifikasi]` `Section\InputRenewalDtl` menyisipkan `InputDtlObject` **empat kali** (L2764 · L2944 ·
L3207 · L3390), sementara section itu **tidak ada di NB maupun RNW**. Penggantinya
`InputDtlObject_FacIn` **ada di keduanya** dan muncul **26 kali** di berkas yang sama.

Keputusan work owner (K-041): **`InputDtlObject` usang, digantikan `InputDtlObject_FacIn`**; keempat
sisipan adalah sisa kelewat.

⛔ **Jangan hapus keempat sisipan itu diam-diam** (`CLAUDE.md` §1). Karena section tujuannya tidak ada,
sisipan itu tidak menghasilkan apa pun dan menjadi tidak berdampak dengan sendirinya.

⚠️ **`InputDtlObject` ≠ `InputDtlObject_FacIn`** — dua section berbeda, dibedakan hanya oleh sufiks.
Jangan disatukan, jangan diperlakukan sebagai salah ketik.

**Blocked by:** R03

**Status:** ready-for-agent

- [ ] `InputRenewalDtl` dan `InputRenewalDtl_IsUW` diport sebagai **dua komponen terpisah**
- [ ] Keduanya **menyisipkan** komponen NB; tidak ada komponen NB yang disalin
- [ ] Ke-37 dan 36 properti khas layar ini terikat benar
- [ ] Keempat sisipan `InputDtlObject` **diport apa adanya**, dengan komentar menyatakan statusnya usang dan section penggantinya
- [ ] Tidak ada penyatuan `InputDtlObject` dengan `InputDtlObject_FacIn`
- [ ] Tidak ada perhitungan premi, spreading, atau komisi yang ditulis ulang di tiket ini — seluruhnya diwarisi NB

## RNW Fac In - R05 - Daftar kandidat renewal dan portalnya

**What to build:** Underwriter membuka daftar polis yang mendekati jatuh tempo dan memilih mana yang
akan diperpanjang — tanpa perlu tahu nomor polisnya lebih dulu. Daftar yang sama muncul di portal.

`[terverifikasi]` `ReportDefinition\RenewalList_RD` — **11 kolom**:

```
.pyID · .pzInsKey · .pxCreateDateTime · .pxCreateOperator
.Quotation.OldPolicyNo · .Quotation.InsuredName · .Quotation.MarketingName
.OfferFacIn.PolicyData.StartDateTime · .OfferFacIn.PolicyData.EndDateTime
.NBStatus · .NBStatusNew
```

Berkunci **`OldPolicyNo`** dan tanggal berakhirnya polis — itulah yang membuatnya daftar *kandidat*,
bukan sekadar daftar kasus.

`[terverifikasi]` `Section\SFAPortal_Renewal` mengikat 9 properti dengan bentuk yang sama.

##### ⚠️ `Work-Renewal` adalah irisan pelaporan, bukan kelas kerja

`[terverifikasi]` `RenewalList_RD` berkelas `ASM-FW-GISFW-Work-Renewal`, **tetapi** `pyWorkClass` pada
flow renewal adalah `ASM-FW-GISFW-Work` — **sama dengan NB** — dan `Work-Renewal` tidak dipakai rule
lain mana pun.

⛔ **Jangan membuat kelas kerja atau antrean tersendiri untuk renewal.** Kasus renewal hidup di kelas
kerja yang sama dengan NB; `Work-Renewal` hanya cara melaporkannya.

**Blocked by:** R01

**Status:** ready-for-agent

- [ ] Daftar kandidat menampilkan kesebelas kolom, berkunci nomor polis lama dan tanggal berakhir
- [ ] Portal menampilkan daftar yang sama
- [ ] Memilih satu baris mengarah ke alur masuk renewal (R01) dengan nomor polis terisi
- [ ] **Tidak ada kelas kerja atau antrean baru** untuk renewal — kasus tetap di kelas kerja NB
- [ ] Penyaringan status dan kelompok tim direproduksi sesuai perilaku lama
- [ ] ⚠️ Kolom nama tertanggung dan nama marketing **ditampilkan sesuai perilaku lama**, tetapi **nilainya tidak pernah masuk fixture, test, atau artefak mana pun**

## RNW Fac In - R06 - Kelompok bisnis pada kasus renewal

**What to build:** Sistem menurunkan **kelompok bisnis** sebuah kasus dari kode bisnis penawarannya,
agar klasifikasi kasus renewal konsisten dengan sistem lama.

`[terverifikasi]` `Activity\GetBusinessGroup_Act`, 4 langkah kedalaman penuh:

| Langkah | Metode | Isi |
| ---: | --- | --- |
| 1 | `Property-Set` | `InputData.CARI2` ← `…QuotationData.BusinessCode` |
| 2 | `RDB-List` | → `CariBusinessGID` |
| 3 | `Property-Set` | `ParamBis.CARI10` ← `OutBis.pxResults(1).CARI3` |
| 4 | `Page-Remove` | bersihkan halaman |

`[terverifikasi]` `RDBList\CariBusinessGID` — `<pyBrowseSQL>`:

```sql
select ID, OLDID, NOTE as "Note", GROUPPANEL as "GroupPanel", BusinessGroupID as CARI3
from business where ID = {InputData.CARI2}
```

##### ⚠️ Ini kelompok bisnis, BUKAN lini bisnis (COB)

Jangan tertukar. **Kelompok bisnis** diturunkan di sini dari tabel `business`. **Lini bisnis (COB)** —
yang menggerakkan skala rasio per K-018 — tetap ditentukan predikat `IsFire` / `IsPA` / `IsMBU` / …
yang **diwarisi dari NB**. Menyatukan keduanya akan merusak pemilihan skala.

`[terverifikasi]` DecisionTable `BusinessType_DeT` **dirujuk nol kali** di korpus renewal.

##### ⛔ Blok `pySaveSQL` DIBUANG — perubahan perilaku yang disengaja

Rule yang sama memuat `<pySaveSQL>` berupa blok PL/SQL yang memanggil `POOLDATA.PROSESCOPY`.
**Blok itu tidak diport** (K-037).

Alasannya, seluruhnya `[terverifikasi]`: prosedur `PROSESCOPY` **tidak ada di basis data** (dikonfirmasi
DBA) · seluruh argumennya **literal keras** · keluarannya ke `dbms_output` yang tidak dibaca aplikasi ·
**tidak berhubungan** dengan `pyBrowseSQL` di rule yang sama.

⚠️ Ini **bukan** porting apa adanya, dan **bukan** perbaikan diam-diam — ia keputusan sadar yang sudah
dicatat bernomor. `<pyBrowseSQL>` **tetap diport apa adanya**.

**Blocked by:** R01

**Status:** ready-for-agent

- [ ] Kelompok bisnis diturunkan: kode bisnis → tabel `business` → `BusinessGroupID`
- [ ] `<pyBrowseSQL>` diport apa adanya, termasuk alias kolomnya
- [ ] **Blok `pySaveSQL` yang memanggil `PROSESCOPY` tidak diport**, dan ketiadaannya dicatat di komentar sebagai keputusan bernomor — bukan dihapus tanpa jejak
- [ ] Kelompok bisnis **tidak** dipakai memilih skala rasio; pemilihan skala tetap lewat resolver COB NB
- [ ] Komentar menyebut activity dan rule RDB asalnya
- [ ] ⚠️ Bila paralel run memperlihatkan selisih yang menunjuk jalur ini, keputusan membuang `pySaveSQL` adalah **tersangka pertama** — catat itu di komentar

## RNW Fac In - R07 - Konversi produksi kasus renewal

**What to build:** Kasus renewal yang disetujui dikonversi ke produksi **lewat jalur yang sama persis
dengan New Business** — bukan jalur kedua yang bisa menyimpang.

`[terverifikasi]` Dasarnya kuat: `Activity\SaveJsonPolicyFacIn_Act` — activity simpan produksi —
**identik byte-per-byte** antara NB dan RNW (371.819 byte, SHA-256 sama). Renewal **tidak memerlukan
penghasil nomor polis tersendiri**.

`[terverifikasi]` `Activity\serviceInsertArasapasRNW_act` — kelas `ASM-FW-GISFW-Work`, **10 langkah**:

```
1 Call serviceInsertArasapas_act    6 Connect-REST
2 Property-Set                      7 Property-Set
3 RDB-List                          8 Call InsertLogServiceProd
4 Property-Set                      9 Page-Remove
5 Call …Int-M_LINK_SERVICE.GetLinkService   10 RDB-List
```

Polanya: **endpoint diambil dari tabel layanan → panggilan REST → catat log layanan**. Ini menguatkan
`CLAUDE.md` §4.4 langsung dari korpus — endpoint memang datang dari tabel, bukan literal.

#### ⛔ Blocker eksternal: jalur produksi NB belum ditiketkan

Tiket ini **ditulis penuh sekarang, tetapi belum dapat dieksekusi.** Jalur simpan produksi berada di
**Out of Scope butir 3 spec NB**, menunggu **struktur tabel flat** (keputusan work owner) dan
**`ALL_SOURCE`** untuk prosedur yang dilewati seluruh tulisan produksi.

⚠️ Celah **K-004** juga tetap terbuka: rule konversi asli kelas `Work` (`serviceInsertArasapas_act`
tanpa sufiks) **tidak ada di korpus mana pun**. Yang tersedia hanya dua saudara yang memperlihatkan
polanya. Arah rancangan mengikuti pola NB (K-035); rule aslinya belum pernah terbaca.

**Blocked by:** R01 · ⛔ **EKSTERNAL — jalur produksi NB** (Out of Scope butir 3 spec NB: tabel flat +
`ALL_SOURCE`; belum ada tiketnya)

**Status:** ready-for-agent

- [ ] Konversi renewal memakai **jalur simpan yang sama** dengan NB — satu implementasi, bukan dua
- [ ] Endpoint diambil dari tabel layanan saat runtime; **tidak ada endpoint literal di kode**
- [ ] Panggilan servis dicatat ke log layanan sesuai perilaku lama
- [ ] Urutan kesepuluh langkah dipertahankan
- [ ] Komentar mencatat bahwa **rule kelas `Work` aslinya tidak ada di korpus** (K-004), dan bahwa rancangan ini mengikuti pola dari dua saudara yang terbaca
- [ ] ⛔ Tiket **tidak dinyatakan selesai** sebelum jalur produksi NB tersedia — blocker eksternalnya dicatat, bukan disiasati

## RNW Fac In - R08 - Rekonsiliasi kasus renewal

**What to build:** Kasus renewal yang dijalankan sistem baru menghasilkan angka yang **identik sampai
digit terakhir** dengan sistem lama — dibuktikan memakai kerangka rekonsiliasi yang sudah ada, **tanpa
pembanding baru**.

Ini yang menutup lingkaran: R01–R06 membangun pintu masuk dan layar; tiket ini membuktikan bahwa di
balik pintu itu, angkanya memang tidak berubah.

#### Mengapa tidak ada pembanding baru

`[terverifikasi]` Renewal **tidak punya mesin hitung ulang**. Hanya empat activity yang menggerbangi
pembeda siklus — masa berlaku tanggal, validasi tanggal, proteksi spreading, input pembayaran — dan
**tidak satu pun perhitungan premi**. Keempatnya berkas bersama yang identik dengan NB.

Membangun pembanding kedua karena itu bukan sekadar mubazir — ia **berbahaya**: dua pembanding dapat
menyimpang, lalu keduanya tampak benar.

📌 Masukannya **sudah tersedia**: salah satu dari lima berkas kasus yang diurus tiket NB-15 adalah
kasus renewal nyata. Tiket ini memakainya lewat kerangka NB-16.

⛔ **Toleransi ditolak** (ADR-F-0001). Bila urutan operasi dan presisi per-langkah direproduksi apa
adanya, nilai desimal bersifat deterministik dan hasilnya **harus** identik.

⛔ **Tidak pernah menyentuh berkas kasus mentah** — hanya fixture ter-de-identifikasi hasil NB-15.

**Blocked by:** **NB-16** (kerangka rekonsiliasi eksak tahap 1) · R01 · R04

**Status:** ready-for-agent

- [ ] Kasus renewal dijalankan lewat **kerangka NB-16**; tidak ada pembanding baru yang dibuat
- [ ] Perbandingan **eksak, tanpa toleransi**; tidak ada parameter epsilon
- [ ] Masukannya **hanya fixture ter-de-identifikasi** (NB-15); tidak ada jalur yang membuka berkas kasus mentah
- [ ] Selisih dilaporkan dengan menyebut **kasus, lini bisnis, dan langkah** tempat angkanya mulai berbeda
- [ ] Nilai dasar akseptasi renewal terbukti **nilai pertanggungan penuh**, bukan selisih — sesuai keputusan yang sudah diambil
- [ ] Bila ada selisih, **keputusan membuang `pySaveSQL` `PROSESCOPY` (R06) diperiksa lebih dulu** sebagai tersangka
- [ ] Hasilnya deterministik dan dapat dijalankan ulang

# Endorsment Fac In

Jumlah tiket: **22**

## Endorsment Fac In - E01 - Registry predikat EDM — sumber data per-rule

**What to build:** Registry predikat yang sudah ada diperluas sehingga **tiap predikat mendeklarasikan
sumber datanya sendiri**. Sebagian besar predikat endorsement membaca properti kasus, persis seperti
New Business; **enam predikat lini bisnis** membaca kolom hasil query atas **polis lama**.

`[terverifikasi]` Dari **202** rule `When` endorsement, **196 membaca properti kasus** dan **hanya 6
membaca hasil query**. Registry tetap **satu** — sumber data adalah atribut **per-rule**, bukan
per-siklus.

`[terverifikasi]` Kolom yang dibaca keenamnya berisi **jenis bisnis polis lama**, diisi satu query
yang berjalan di alur masuk sebelum kasus lahir. Karena itu **urutan eksekusi mengikat**: query
dijalankan lebih dulu, baru predikat dapat dievaluasi.

**Asal (Pega).** `When\IsAneka` · `IsFire` · `IsPA` · `IsMBU` · `IsMarineCargo` · `IsGolfInsurance` ·
`When\IsEDM` · `When\IsNotEDM` · pengisi `RDBList\GetBusinessType_Sql` lewat
`Activity\SetErrorBatalEndorsement_Act` langkah 7.

**Keputusan.** K-050 · K-018 · K-029 · `CLAUDE.md` §4.5, §4.6

**Blocked by:** `..\08-registry-rules-eval.md`
⛔ Tiket itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] Registry **satu**; sumber data dideklarasikan **per-rule**
- [ ] **37 cabang** lini bisnis terimplementasi (5 · 25 · 1 · 4 · 1 · 1) — **bukan 203**
- [ ] Kolom sumber = jenis bisnis **polis lama**, versi terakhir (`COUNT−1`; **aman**, tidak ada penghapusan baris — K-050)
- [ ] Urutan mengikat: query pengisi berjalan **sebelum** predikat lini bisnis dievaluasi; bila belum, keenamnya bernilai salah dan percabangan runtuh **diam-diam**
- [ ] **Kedua** tag kondisi dibaca; yang mengikat **ekspresi tersimpan**, bukan label tampilan
- [ ] Nama rule tak dikenal → `panic`
- [ ] Predikat lini bisnis menyalakan resolver skala (K-018) — salah lini berarti salah satuan rate
- [ ] Tiap predikat menyebut rule Pega asalnya dalam komentar (§4.6)
- [ ] **K-046** `K046_IsMBU_CabangGanda_HasilTidakBerubah` — dua nilai muncul dua kali; redundan, hasil tidak berubah
- [ ] **K-046** `K046_IsAneka_LabelKodeBisnis_Diabaikan_IkutiValue1` — label tampilan tidak berkaitan dengan ekspresi tersimpan
- [ ] **K-046** `K046_COB_EDM_TanpaDelegasiSubRule` — endorsement membandingkan literal, tidak memanggil sub-rule

## Endorsment Fac In - E02 - Tiga predikat yang perilakunya berbeda di endorsement

**What to build:** Tiga predikat bernama sama dengan New Business tetapi **berperilaku berbeda** di
endorsement, ditambah satu penegasan yang mudah salah diimplementasikan.

`[terverifikasi]` Di endorsement, operator ber-workbasket **marketing dihitung sebagai underwriter** —
menentukan **siapa yang boleh menyetujui**. Predikat klaim menguji **keberadaan halaman data**, bukan
prefiks ID kasus. Predikat travel menerima **dua sumber**, bukan satu.

⛔ **Predikat "bukan endorsement" BUKAN negasi predikat "endorsement".** Keduanya membaca **jalur
properti berbeda** — satu agregat tersimpan, satu halaman aktif. Keduanya dapat bernilai benar
bersamaan, atau salah bersamaan.

**Asal (Pega).** `When\IsUW` · `When\IsClaim` · `When\IsTravel` · `When\IsEDM` · `When\IsNotEDM`

**Keputusan.** K-046 · K-050 · `CLAUDE.md` §4.5

**Blocked by:** E01

**Status:** blocked

- [ ] Predikat underwriter endorsement mengakui **satu workbasket lebih banyak** daripada New Business
- [ ] Predikat klaim endorsement menguji keberadaan halaman data; New Business menguji prefiks ID
- [ ] Predikat travel endorsement menerima dua sumber; New Business satu
- [ ] Predikat "bukan endorsement" punya **implementasi sendiri** — membaca halaman aktif
- [ ] Kasus uji: kedua halaman sengaja **tidak sinkron** → membuktikan keduanya dapat bernilai sama
- [ ] **K-046** `K046_IsUW_MarketingSebagaiUnderwriter_EDM`
- [ ] **K-046** `K046_IsClaim_UjiPrefiks_vs_UjiHalaman`
- [ ] **K-046** `K046_IsTravel_LabelTidakDieksekusi` — label tampilan menyebut kode bisnis yang tidak dieksekusi

⚠️ Siapa yang boleh menyetujui di endorsement masih **`[pertanyaan terbuka]`** milik Underwriting.
Sampai dijawab, perilakunya **diport apa adanya**.

## Endorsment Fac In - E03 - Tracer — alur masuk endorsement sampai kasus berfase Policy

**What to build:** Seorang petugas memilih **nomor polis + jenis endorsement + tanggal**, dan sebuah
kasus endorsement tercipta — atau **ditolak dengan alasan spesifik**.

Ini irisan **tertipis yang lengkap** untuk alur masuk. Begitu hijau, premisnya terbukti: endorsement
lahir **langsung di fase polis**, tanpa fase penawaran maupun binding.

`[terverifikasi]` Kasus lahir dari kelas portal menuju kelas kerja berprefiks `EDM-`, dengan **tautan
tiga arah** antara kasus portal, kasus endorsement, dan handle assignment. Tidak ada satu pun
konektor di flow yang menetapkan fase penawaran atau binding.

⛔ **Penolakan adalah keluaran bisnis yang sah**, bukan kegagalan teknis — keduanya harus dapat
dibedakan.

**Asal (Pega).** `Activity\SetValueToEDMWork` langkah 7–13 · `DataTransform\DataToEDM` ·
`Flow\InputAddendumFacultativeIn` konektor `Start2 → Assignment7` · `Activity\SetEdmType`

**Keputusan.** K-049 (Seam 5) · K-029 · K-044

**Blocked by:** `..\08-registry-rules-eval.md` · E01
⛔ Tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] **Seam 5 `services/endorsement.OpenCase` hidup**
- [ ] Kasus lahir berkelas kerja endorsement, berprefiks **`EDM-`**
- [ ] **Tautan tiga arah** terisi lengkap — kasus portal ↔ kasus endorsement ↔ handle assignment
- [ ] Nomor polis lama dan alasan endorsement terisi
- [ ] Kasus **berfase Policy**; kasus uji negatif: **tidak pernah** melewati fase penawaran atau binding
- [ ] Ditempatkan di workbasket marketing dengan tiket admin polis
- [ ] Jenis endorsement dari domain terkunci (K-029); kode usang **tetap diport**, cabangnya tidak dihapus
- [ ] Penolakan **terbedakan** dari kegagalan sistem, dengan alasan spesifik per gerbang
- [ ] Tiap transisi menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E04 - Enam gerbang penolakan dan empat klep pembatalnya

**What to build:** Enam pemeriksaan yang dapat **menolak** pembuatan kasus endorsement, empat di
antaranya dapat **dibatalkan** oleh tabel pengecualian.

`[terverifikasi]` Polanya seragam per gerbang: ambil data pemeriksa → siapkan penanda klep → panggil
laporan pembatal → pasang galat **kecuali** klep terisi.

⛔ **Hanya empat dari enam gerbang punya klep.** Gerbang "sudah ada endorsement lain yang belum
selesai" dan gerbang "tidak ter-spreading fac out" **tidak dapat dibatalkan**.

⛔ **Klep kosong berarti penolakan berlaku penuh.** Itu perilaku sistem lama dan arahnya aman —
diport apa adanya.

**Asal (Pega).** `Activity\SetErrorBatalEndorsement_Act` langkah 5 (24 sub-langkah) ·
`RDBList\GetEDMStatus_SQL` · `GetDataClaim_SQL` · `SearcStatusBayarArasaps_SQL` ·
`GetListRNWbyNopolis_SQL` · `GetFacoutList_SQL` · ReportDefinition `GetListEdm` ·
`BrowseOpenProteksiEdm_RD`

**Keputusan.** K-049 · K-046 · K-006

**Blocked by:** E03

**Status:** blocked

- [ ] Keenam gerbang menolak pada kondisi yang benar **dan hanya** pada kondisi itu
- [ ] Gerbang pembayaran **hanya** berlaku untuk jenis endorsement pembatalan
- [ ] **Empat klep** membatalkan penolakan yang benar, masing-masing dengan penanda jenisnya sendiri
- [ ] **Gerbang kedua dan keenam tetap menolak** meski tabel pengecualian terisi — keduanya tanpa klep
- [ ] Klep kosong atau tabel tak terjangkau → **penolakan berlaku penuh**
- [ ] Bypass "endorsement RI slip" melewati gerbang yang benar — diport apa adanya
- [ ] Rule pengecekan pembatal **diimplementasikan**; **isi tabelnya di luar lingkup** dan bukan blocker
- [ ] Rujukan laporan bersifat **runtime** — tidak dapat divalidasi saat kompilasi; ketiadaannya **bukan `panic`** karena tabelnya punya pemilik (K-006)
- [ ] Tiap gerbang menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E05 - Validasi tanggal endorsement

**What to build:** Tanggal endorsement wajib berada **di dalam periode polis** yang di-endors.
Sejumlah nomor polis tertentu dikecualikan dari validasi ini.

`[terverifikasi]` Periode dibaca dari tabel produksi; tanggal di luar rentangnya memunculkan galat.
Pengecualiannya berupa **daftar nomor polis yang tertanam di kode lama**.

⛔ **Nomor polis produksi tidak boleh menjadi literal di kode target.** Ia menjadi **data konfigurasi
yang dapat diaudit** — itu keputusan rancangan, **bukan** perubahan perilaku. Daftar
pengecualiannya tetap sama persis.

**Asal (Pega).** `Activity\CheckEDMPolisDate` (6 langkah) · `RDBList\GetStartDate`

**Keputusan.** K-046 · `CLAUDE.md` §3.5, §4.6

**Blocked by:** E03

**Status:** blocked

- [ ] Tanggal di luar periode polis → galat dengan pesan yang dapat dibaca petugas
- [ ] Daftar pengecualian menjadi **konfigurasi**, bukan literal di kode
- [ ] **Perilaku tidak berubah** — kasus yang dikecualikan sistem lama tetap dikecualikan
- [ ] ⛔ Nomor polis **tidak pernah** muncul di kode, tiket, test, fixture, maupun log
- [ ] Test memakai **nomor sintetis**, bukan nomor produksi
- [ ] **K-046** `K046_ValidasiTanggal_PerbandinganSubstring8Karakter` — sistem lama membandingkan **potongan string tanggal**, bukan objek tanggal; diport apa adanya
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E06 - Tracer — before-image tiga lapis untuk satu lini

**What to build:** Ketika kasus endorsement dibuka, "keadaan sebelum" tersedia dalam **tiga lapis
terpisah**, dan satu berkas kasus kebakaran nyata **rekonsiliasi eksak** terhadap Pega.

Ini irisan **tertipis yang lengkap** untuk before-image. Begitu hijau, seluruh premisnya terbukti:
tiga mekanisme berbeda, siklus hidup berbeda, konsumen berbeda.

`[terverifikasi]` **Lapis A** memuat dokumen polis versi terakhir lalu menjalankan penyalinan
sehingga data kerja **identik** dengan polis lama. **Lapis B** mengisi nilai lama per baris.
**Lapis C** menandai baris warisan.

⛔ **Delta dihitung terhadap polis lama, bukan terhadap kosong.** Implementasi **tidak boleh** mulai
dari struktur kosong lalu menambahkan perubahan.

⚠️ Istilah "before-image" adalah **nama konsep**; tidak ada properti bernama itu di korpus.

**Asal (Pega).** `Activity\SetValueToEDMWork` langkah 14 (13 sub-langkah) dan langkah 15 ·
`RDBList\GetEDMOldData_SQL` · `Activity\SetOldData` blok 4.1–4.2 ·
`Activity\SetOLDValueToEDMWork_FIRE` · `Activity\InputAddendumFacIn_PreAct` langkah 28

**Keputusan.** K-047 (Seam 4) · K-010/K-012 · K-018 · K-027 · K-048

**Blocked by:** `..\01-tracer-money-ratio-premi-pa.md` · E01
⛔ Tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] **Seam 4 `services/endorsement.PrepareBeforeImage` hidup**
- [ ] **Tiga lapis tetap terpisah** di model domain — menyatukannya mengubah angka
- [ ] Lapis A memuat **versi terakhir** polis dan menjalankan penyalinan berpola
- [ ] Data kerja setelah penyalinan **identik** dengan polis lama
- [ ] Lapis B mengisi properti nilai lama lini kebakaran; nilai kosong → **nol**, bukan nilai tak-ada
- [ ] Lapis C menandai baris warisan sesuai kedalaman daftar
- [ ] Porsi periode: presisi desimal tinggi + **guard pembagian nol**; bertipe **`Ratio`**, bukan `Money`
- [ ] ⚠️ Porsi periode **bersifat sementara** — ditimpa modul selisih; **urutan eksekusi before-image → selisih wajib dijaga** (K-048)
- [ ] Nilai uang bertipe `Money`; `Money + Ratio` **gagal saat kompilasi**
- [ ] **Satu fixture kebakaran rekonsiliasi eksak**
- [ ] Tiap penyalinan menyebut rule Pega asalnya (§4.6)

## Endorsment Fac In - E07 - Lapis A — sembilan penyalinan yang menyimpang dari pola

**What to build:** Penyalinan dokumen polis lama tidak seluruhnya berbentuk "salin field bernama
sama". Sembilan di antaranya berbeda, dan **empat halaman akan kosong** bila diabaikan.

`[terverifikasi]` Dari **54 penugasan**: **51** bersumber dokumen polis lama, **3** tidak. Dari yang
51, **45** berpola seragam dan **6 menyimpang** — menulis ke halaman **selain** agregat penawaran.

> 🔁 **Arsip lintas-siklus keliru di titik ini.** Ia menyatakan seluruh 54 berbentuk seragam. Itu
> **tidak berlaku untuk 9 dari 54**. Implementasi yang mengikuti arsip meninggalkan empat halaman
> kosong.

**Asal (Pega).** `Activity\SetValueToEDMWork` langkah 14.3

**Keputusan.** K-047 · `CLAUDE.md` §4.6

**Blocked by:** E06

**Status:** blocked

- [ ] **54 penugasan** terimplementasi, tidak diringkas
- [ ] **Enam cermin lintas-halaman** ditulis eksplisit: jenis bisnis ke halaman kutipan · tiga field pembayaran ke halaman polis · penanda B2B ke akar objek kerja · nama marketing ke halaman parameter kredit
- [ ] **Kriteria kunci:** keempat halaman itu **tidak boleh kosong** setelah lapis A berjalan
- [ ] Tiga field pembayaran ditulis **dua kali** — ke agregat penawaran **dan** ke halaman polis
- [ ] Tiga penugasan yang **bukan** dari dokumen lama tetap dijalankan
- [ ] Empat kelompok field terisi lengkap: identitas bisnis · periode · struktur share dan kapasitas · pembayaran
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E08 - Lapis B — nilai lama per baris untuk enam lini sisanya

**What to build:** Nilai lama per baris untuk golf, aneka, personal accident, marine cargo,
kendaraan bermotor, dan travel — plus tingkat cedant yang sama di ketujuh lini.

`[terverifikasi]` Pola nilainya **tunggal dan tanpa pengecualian**: bila properti sumber tidak
kosong pakai nilainya, bila kosong pakai **nol**.

⛔ **Kosong dipetakan ke nol, bukan ke nilai tak-ada.** Aritmetika selisih bergantung pada ini.

📌 Premi Nusantara Re lama **hanya ada di dua lini** — marine cargo dan kendaraan bermotor.

**Asal (Pega).** `Activity\SetOldData` blok 4.3–4.14

**Keputusan.** K-047 · K-010/K-012 · K-018 · `CLAUDE.md` §4.6

**Blocked by:** E06

**Status:** blocked

- [ ] Keenam lini mengisi properti nilai lama **masing-masing sesuai petanya**, dan hanya itu
- [ ] Tingkat cedant identik di ketujuh lini
- [ ] **Seluruh** penugasan nilai lama memakai pola kosong→nol; **nol pengecualian**
- [ ] Premi Nusantara Re lama muncul **hanya** di dua lini, bukan tujuh
- [ ] Nilai uang bertipe `Money` dengan mata uangnya; rate lama bertipe `Ratio` dengan **skala per lini** (K-018)
- [ ] Lapis B diisi ulang **setiap kali layar endorsement dibuka** — bukan sekali saat kasus lahir
- [ ] ⚠️ Lapis B menangani **travel tetapi bukan jiwa**; himpunannya berbeda dari lapis C
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E09 - Dua guard yang dipertahankan dan satu perbaikan sadar

**What to build:** Dua guard di sistem lama menguji **properti tujuan** alih-alih properti sumber.
Keduanya **diport apa adanya**. Satu hal diperbaiki — dan hanya satu.

`[terverifikasi]` Guard rate lama berperilaku demikian pada **seluruh** kemunculannya; guard premi
lama hanya pada **satu** dari delapan belas kemunculan.

⛔ **Keduanya perilaku yang BENAR menurut work owner (K-046).** Meluruskannya mengubah angka selisih
dan memutus rekonsiliasi paralel run.

⚠️ **A.5 — satu-satunya perbaikan sadar di seluruh modul before-image.** Nilai cadangan premi
Nusantara Re lama memakai **string** di sistem lama; di sistem baru ia **angka nol**, karena nilai
uang bertipe `Money` dan `Money` tidak boleh menampung string.

**Asal (Pega).** `Activity\SetOldData` — penugasan premi Nusantara Re lama (dua tempat) dan guard
rate/premi lama

**Keputusan.** **K-046** · **A.5** · K-010/K-012 · ADR-F-0001

**Blocked by:** E08

**Status:** blocked

- [ ] **K-046** `K046_RateOld_GuardSelfReferential_6dari6` — guard menguji properti tujuan; **perilaku benar**, bukan cacat
- [ ] **K-046** `K046_PremiumOld_SelfReferential_HanyaFire` — hanya satu dari delapan belas kemunculan
- [ ] Komentar kode menyebut rule Pega asal **dan** K-046, agar tidak "diluruskan" pembaca berikutnya
- [ ] **A.5:** nilai cadangan premi Nusantara Re lama adalah **angka nol**, bukan string — ditandai **perbaikan, bukan port**
- [ ] Pembanding rekonsiliasi **menormalkan** string dan angka sebelum membandingkan
- [ ] Selisih tipe yang muncul **dijelaskan oleh A.5** dalam prosedur rekonsiliasi — bukan ditandai cacat
- [ ] ⚠️ Cache ekspresi di korpus memuat bentuk berbeda dari yang dieksekusi; **yang mengikat ekspresi tersimpan**

## Endorsment Fac In - E10 - Lapis C — penanda baris warisan, tujuh varian

**What to build:** Setiap simpul daftar yang berasal dari polis lama ditandai, sehingga baris warisan
dapat dibedakan dari baris yang ditambahkan saat endorsement.

`[terverifikasi]` Tujuh varian, satu per lini bisnis. **Cacah penanda mengikuti kedalaman model data
lini itu**, bukan kelengkapan implementasi — lini dengan empat tingkat daftar bersarang menandai
lebih banyak simpul daripada lini dengan satu tingkat.

⚠️ Nilai penanda tersimpan **berkutip**. Pola pencarian tanpa kutip memberi **nol hasil** dan pernah
menyesatkan.

**Asal (Pega).** `Activity\SetOLDValueToEDMWork_{FIRE,MC,Aneka,MBU,LIFE,PA,GOLF}` ·
dipanggil `Activity\SetValueToEDMWork` langkah 14.7–14.13

**Keputusan.** K-047 · K-046 · `CLAUDE.md` §4.6

**Blocked by:** E06

**Status:** blocked

- [ ] Ketujuh varian menandai simpul sesuai **kedalaman daftar** lini masing-masing
- [ ] Nilai penanda sesuai yang tersimpan di korpus, **berkutip**
- [ ] **Hanya varian kebakaran** yang menyetel penanda prorata tingkat agregat; enam lainnya tidak
- [ ] Penanda prorata itu **tidak** termasuk field yang disalin lapis A
- [ ] **K-046** `K046_Life_DuaPenandaOldData` — varian jiwa menyetel **dua** penanda berbeda; enam lainnya satu
- [ ] ⚠️ Lapis C menangani **jiwa tetapi bukan travel** — himpunannya berbeda dari lapis B (E08)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E11 - Tiga gerbang keluar lapis B

**What to build:** Pengisian nilai lama per baris **berhenti sebelum berjalan** pada tiga keadaan.
Ketiganya menghasilkan lapis B kosong — dan itu **perilaku yang benar**, bukan galat.

`[terverifikasi]` Gerbang pertama mengeluarkan kasus **jiwa**; gerbang kedua mengeluarkan kasus yang
**bukan endorsement**; gerbang ketiga mengeluarkan polis dengan **lokasi melebihi ambang**.

⛔ **Polis besar tidak mendapat nilai lama per baris sama sekali** — layar dan perhitungan per baris
menampilkan nol. Diport apa adanya; alasan ambangnya tidak diketahui dari korpus.

📌 Gerbang kedua memakai predikat yang membaca **agregat tersimpan**. Predikat kebalikannya membaca
**halaman aktif** — lihat E02; jangan diimplementasikan sebagai negasi.

**Asal (Pega).** `Activity\SetOldData` langkah 1–3

**Keputusan.** K-044 · K-047 · K-046

**Blocked by:** E08

**Status:** blocked

- [ ] Kasus **jiwa** keluar sebelum blok pengisian — jiwa memakai jalurnya sendiri (E18)
- [ ] Kasus **bukan endorsement** keluar
- [ ] Polis dengan **lokasi melebihi ambang** keluar
- [ ] Ketiganya menghasilkan lapis B **kosong tanpa galat** — bukan kegagalan
- [ ] Ambang lokasi diambil dari korpus apa adanya; alasannya tetap **`[pertanyaan terbuka]`**
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E12 - Prefactor — Seam 3 menjadi satu pintu dengan empat bentuk

**What to build:** Pintu perhitungan premi yang sudah ada diperluas agar dapat memilih di antara
**empat bentuk perhitungan**, tanpa menambah seam. *"Make the change easy, then make the easy
change."*

`[terverifikasi]` Endorsement **memakai ulang** mesin premi New Business — puluhan activity
perhitungan ada di kedua korpus. Yang ditambahkan endorsement adalah **bentuk**, bukan mesin baru.

⛔ **SATU pintu, empat bentuk — bukan seam per bentuk.**

| Bentuk | Isi |
| --- | --- |
| **1 — dasar** | `TSI × Rate × faktor ÷ pembagi-komposit`; pembagi diturunkan **resolver K-018** dari faktor yang ikut |
| **2 — dua-bagian endorsement** | porsi baru untuk sisa periode **+** porsi lama untuk periode berjalan |
| **3 — tabel tarif** | premi dari lookup **+** tambahan per minggu; masukan dari **repository** |
| **4 — dekomposisi delta coverage** | `(Δrate × TSI_lama) + (ΔTSI × rate)`, dengan cabang periode pendek |

⛔ **Dua keluarga rasio prorata tetap terpisah** — yang untuk premi berbeda penyebut, satuan, dan
skala dari yang untuk delta spreading. Namanya nyaris sama; menyatukannya menggeser angka.

⛔ **Porsi periode dipakai sebagai pecahan, langsung** — tidak dibagi seratus lagi.

**Asal (Pega).** `Activity\CalculatePremiFire` · `HitungPremiOnChange` · `calculatePremiPA` ·
`CalculatePremiumTravel` · `SetLocalNonMbuProrate` · `CalculatePremiPA_FacIn`

**Keputusan.** **K-051** · K-018 · K-010/K-012 · K-027 · ADR-F-0004/0005

**Blocked by:** `..\01-tracer-money-ratio-premi-pa.md` · `..\03-resolver-cob-skala.md` ·
`..\04-rumus-premi-pa.md` · `..\05-rumus-premi-mbu.md` · `..\06-rumus-premi-layering.md`
⛔ Kelimanya **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] Satu pintu memilih di antara **empat bentuk** berdasarkan lini bisnis dan konteks endorsement
- [ ] Bentuk 1 tetap **satu implementasi**, dipakai New Business, Renewal, dan Endorsement
- [ ] Pembagi komposit diturunkan resolver (K-018) — **bukan** angka hafalan per lini
- [ ] Bentuk 2 adalah **komposisi** dua panggilan bentuk 1, bukan cabang di dalam rumus
- [ ] Bentuk 3 menerima hasil lookup sebagai **masukan**, tidak menghitungnya sendiri
- [ ] **Dua keluarga rasio prorata terpisah** dan tidak dapat tertukar
- [ ] Porsi periode **pecahan, dipakai langsung**
- [ ] Pembulatan literal per langkah (K-011); nilai berkoma desimal (K-027)
- [ ] `Money + Ratio` **gagal saat kompilasi**; jembatan tunggal `Money × Ratio → Money`

**Catatan sebelum implementasi** — bukan blocker spec:
rumus PA & travel penuh (tarif dari repository, sumbernya nihil di korpus) ·
perhitungan PA bersama sudah terbukti **identik** New Business↔Endorsement ·
**keputusan desain:** empat bentuk sebagai satu fungsi bercabang **atau** empat implementasi di balik
satu pintu · **lima dari enam** salinan mesin premi berbeda isinya dan perlu dibedah sebelum dipakai
ulang.

## Endorsment Fac In - E13 - Selisih pembayaran tingkat mata uang

**What to build:** Untuk setiap mata uang pada kasus endorsement, hitung **selisih pembayaran** —
nilai sesudah dikurangi nilai sebelum.

`[terverifikasi]` Rumus kanoniknya tunggal dan berulang di belasan penugasan. Nilai "sebelum"
diambil dari **lapis A** (dokumen polis lama), bukan dari nilai lama per baris.

⛔ **Netto menambahkan PPh dan PPN, tidak mengurangkannya.** Terlihat berlawanan intuisi pajak;
diport apa adanya.

⚠️ **Urutan mengikat:** transformasi perhitungan premi **menimpa** porsi periode yang disetel modul
before-image sebelum jalur produksi membacanya (K-048).

**Asal (Pega).** `Activity\CountEndorsementData` · `CountDataEDMElse` · `CountPaymentEdm_Act` ·
`CountPaymentEdmTSIObj_Act` · `DataTransform\CountPremiEDM_DT`

**Keputusan.** K-048 · K-046 · K-010/K-012

**Blocked by:** E12 · E06

**Status:** blocked

- [ ] Selisih pembayaran = nilai sesudah − nilai sebelum, per mata uang
- [ ] Netto = premi − komisi − brokerage − potongan **+ PPh + PPN**
- [ ] Nilai sebelum bersumber **lapis A**, bukan lapis B
- [ ] Lima komponen "lama" tingkat pembayaran ikut terisi: komisi · PPN · PPh · brokerage · potongan
- [ ] Mata uang yang **tidak ada** di polis lama → seluruh nilai dianggap **baru**, nilai sebelum dipaksa nol
- [ ] Seluruh nilai bertipe `Money` dengan mata uangnya
- [ ] ⚠️ Urutan eksekusi terhadap modul before-image **dijaga** (K-048)
- [ ] **K-046** `K046_NettoPremi_PPhPPNDitambahkan`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E14 - Tabel rumus premi-menjadi per lini, jenis endorsement, dan jalur

**What to build:** Properti "premi menjadi" **tidak dihitung dengan satu rumus**. Ia punya belasan
bentuk berbeda yang dipilih menurut **lini bisnis × jenis endorsement × jalur**.

`[terverifikasi]` **Dua puluh dua penugasan** di enam berkas menghasilkan **sebelas rumus berbeda
secara semantik**.

⛔ **JANGAN diseragamkan.** Modelkan sebagai **tabel keputusan**, bukan satu fungsi dengan cabang
tersembunyi. Menyeragamkannya mengubah angka pada lini tertentu.

📌 Pembagian "pakai nilai pembayaran lama vs premi lama" **bukan sifat lini** — ia berlaku hanya pada
satu cabang tertentu, dan di sana hanya **dua lini** yang memakai bentuk berbeda.

**Asal (Pega).** `Activity\CountEndorsementData` · `CountDataEDMElse` · `CountPaymentEdm_Act` ·
`CountPaymentEdmTSIObj_Act` · `ReCountPremiLifeEDM` · `CopyAllObj_ACT` ·
`DataTransform\CountPremiEDM_DT`

**Keputusan.** **K-046** · K-029 · K-048

**Blocked by:** E13

**Status:** blocked

- [ ] Terimplementasi sebagai **tabel** (lini × jenis endorsement × jalur), bukan satu rumus
- [ ] Kedua puluh dua penugasan terwakili; **sebelas bentuk semantik** dapat dibedakan
- [ ] Cabang "batal": **dua lini** memakai nilai pembayaran lama, empat lainnya premi lama
- [ ] Cabang "batal sejak semula": ketujuh lini memakai bentuk yang **sama**
- [ ] **K-046** `K046_EDMPremiMenjadi_TidakSeragam_PerLiniDanEdmType`
- [ ] **K-046** `K046_EdmType1_HanyaLife_LiniLainTidakDihitungUlang` — hanya satu lini yang cabang perhitungannya berjalan untuk jenis endorsement pertama
- [ ] **K-046** `K046_LabelMBU_GerbangIsMarineCargo_IkutiKondisi` — label langkah menyebut satu lini, gerbangnya lini lain; **ikuti gerbang**
- [ ] **K-046** `K046_LabelBatalSejakSemula_Gerbang2_BukanSatu` — label dan gerbang bertentangan di lima lini
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E15 - Selisih per baris spreading dan pemetaan kolomnya

**What to build:** Untuk setiap baris spreading, hitung selisih antara nilai sesudah dan nilai
sebelum, lalu petakan ke pasangan kolom produksi.

⛔ Tiket ini **menghasilkan nilai**; **menuliskannya ke tabel produksi milik E21** yang masih
ter-block.

`[terverifikasi]` Rumus kanoniknya: nilai baru **diprorata lebih dulu**, baru dikurangi nilai lama
yang **jenis treaty-nya sama**.

⛔ **Nilai lama diambil dari lapis A**, bukan dari nilai lama per baris. Ini titik yang paling mudah
salah diimplementasikan.

⛔ **Pasangan dicocokkan menurut jenis treaty, bukan indeks posisi.** Penambahan atau penghapusan
baris tidak boleh menggeser pasangan.

**Asal (Pega).** `Activity\SaveFacinProdEDMFire_Act` (puluhan penugasan pengurangan) dan lima cabang
lini lainnya · `RDBList\InsertTreatyProduction_Sql` untuk pemetaan kolom

**Keputusan.** K-048 · K-027 · K-046 · K-010/K-012

**Blocked by:** E12 · E06

**Status:** blocked

- [ ] Selisih = (nilai baru × porsi sisa periode) − nilai lama berjenis treaty sama
- [ ] Nilai lama dari **lapis A** — kasus uji dengan lapis A dan lapis B sengaja **berbeda**
- [ ] Pencocokan menurut **jenis treaty**; kasus uji dengan baris ditambah dan dihapus di tengah
- [ ] Tujuh varian menyimpang terimplementasi: penyesuaian spreading · penyesuaian mata uang (dua arah) · penyesuaian rate · coverage hilang → negatif · jenis endorsement pembatalan → balik tanda
- [ ] Sembilan pasangan kolom terpetakan
- [ ] Nilai berkoma desimal dibaca sesuai **K-027**
- [ ] **K-046** `K046_PasanganMenjadiSelisih_TidakSetangkup` — kolom "sesudah" lebih sedikit daripada kolom "selisih"; tiga pasangan memakai kolom tanpa sufiks
- [ ] **K-046** `K046_EjaanKolomPCT_BROKERGARE_FEE_Dipertahankan` — kolom persen brokerage salah eja di skema; skema tidak berubah (§4.3)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E16 - Berkas pendukung layar endorsement

**What to build:** Ratusan berkas pendukung khas endorsement — layar, aksi layar, transformasi data,
dan aktivitas pembantu yang tidak punya padanan di New Business.

`[terverifikasi]` **350 berkas diport.** Empat berkas fitur pilih-tertanggung **tidak diport**
(K-044) meskipun ada di korpus — itu keputusan sadar, bukan kelalaian.

⛔ **Tidak satu pun kelompok menjadi modul domain baru.** Yang menyentuh perhitungan masuk ke modul
premi lewat Seam 3; sisanya lapisan tampilan dan pendukung.

⚠️ **Klasifikasi wajib dari tipe rule sebenarnya, bukan nama folder.** Seluruh berkas di folder
bernama "RDBList" ternyata bertipe rule SQL — itu aturan di korpus ini, bukan pengecualian.

**Asal (Pega).** 94 Activity · 81 Section · 54 FlowAction · 29 DataTransform EDM-only

**Keputusan.** K-051 · K-044 · K-046 · K-006 · `CLAUDE.md` §3.3, §4.6

**Blocked by:** E01 · E03

**Status:** blocked

- [ ] **350 berkas diport**; keempat berkas pilih-tertanggung **tidak** — tercatat sebagai keputusan
- [ ] Pasangan layar berakhiran penanda underwriter tetap **dua komponen terpisah** — ikut sistem lama
- [ ] Berkas yang bernama sama dengan New Business **tetapi bertipe rule berbeda** tidak memakai ulang implementasi New Business
- [ ] Klasifikasi dari tipe rule sebenarnya, **bukan** nama folder
- [ ] **Lima berkas tanpa rujukan** diport apa adanya bila terbaca — ⛔ **tidak divonis usang** (K-006); statusnya tetap **`[pertanyaan terbuka]`**
- [ ] Tiap berkas menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E17 - Lapisan query dan lookup endorsement ⛔ BLOCKED

**What to build:** Query, halaman data, dan definisi laporan khas endorsement — lapisan yang membaca
Oracle untuk kebutuhan layar dan perhitungan.

⛔ **BLOCKED.** Seam repository **belum ada**, dan bentuknya belum dapat dirancang tanpa skema yang
sudah pasti. Ini **kewajiban tertunda sejak spec New Business**, bukan kelalaian putaran ini.

`[terverifikasi]` Termasuk di dalamnya **tarif travel**, yang rule query-nya **nihil di seluruh
korpus**. Struktur tabelnya tidak dapat diketahui dari korpus dan **tidak ditebak**.

**Asal (Pega).** 22 rule SQL (folder `RDBList\`) · 14 halaman data · 28 definisi laporan EDM-only

**Keputusan.** K-051 · K-027 · `CLAUDE.md` §4.3, §4.5

**Blocked by:** ⛔ **seam repository** (belum ada)

**Status:** blocked

- [ ] ⛔ Menunggu seam repository dirancang
- [ ] Tarif travel menjadi **masukan dari repository** — bukan dihitung di modul premi (lihat E12 bentuk 3)
- [ ] Rule query tarif travel **nihil di korpus** → **`[pertanyaan terbuka]`**, jangan ditebak
- [ ] Tabel pita tarif personal accident juga **belum diketahui** strukturnya
- [ ] Seluruh nilai berkoma desimal dibaca sesuai **K-027**
- [ ] Endpoint dan host dari **konfigurasi**, tidak pernah literal (`CLAUDE.md` §4.4)
- [ ] Skema Oracle **tidak berubah** (§4.3) — termasuk nama kolom yang salah eja
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E18 - Jalur Life sebagai alur tersendiri

**What to build:** Endorsement jiwa berjalan lewat **flow `InputEDMLife` tersendiri** — bukan cabang
kondisi di dalam alur lini umum.

`[terverifikasi]` Empat titik cabang memisahkannya: pengisian nilai lama per baris **keluar** untuk
jiwa · penomoran endorsement memakai pasangan langkah **komplementer** · jalur produksi jiwa berjalan
**sejajar**, bukan di dalam percabangan tujuh lini · penanda baris warisan punya varian sendiri.

⛔ **Jiwa melewati tangga akseptasi** (K-044). Korpus tidak memuat rule tangga akseptasi untuk jiwa.
⚠️ **`[pertanyaan terbuka]`** apakah itu memang tanpa persetujuan — milik work owner dan
Underwriting. Sampai dijawab, **diport apa adanya**.

⚠️ Tiket ini **tidak bersandar** pada klaim arsip bahwa nilai dasar akseptasi adalah selisih TSI —
klaim itu tetap **`[belum diuji]`**.

**Asal (Pega).** `When\IsLife` (enam belas cabang atas kode bisnis lama) ·
`Activity\SetOldData` langkah 1 · `Activity\SaveEDMToJsonPolicy_Act` langkah 12–13, 15, 28–29 ·
`Activity\SetOLDValueToEDMWork_LIFE` · `RDBList\GetKodeProdLife_SQL`

**Keputusan.** **K-044** · K-046 · K-006

**Blocked by:** E06 · E08 · E10

**Status:** blocked

- [ ] Flow **`InputEDMLife` tersendiri** — dimodelkan sebagai alur, bukan cabang `if`
- [ ] Predikat jiwa membaca **kode bisnis lama**, bukan jenis bisnis dari hasil query (berbeda dari predikat lini lain, E01)
- [ ] Jiwa **melewati lapis B** — gerbang keluar E11
- [ ] Penomoran endorsement jiwa lewat query kode produk jiwa; rule generator bernama-jiwa **nihil di korpus** dan itu **bukan `panic`** — penomorannya tertanam di pemanggil
- [ ] Jalur produksi jiwa berjalan **sejajar**, bukan di dalam percabangan tujuh lini
- [ ] ⛔ Jiwa **tanpa tangga akseptasi** — Seam 2 tidak dipanggil di jalur ini
- [ ] Status "tanpa persetujuan" tercatat sebagai **`[pertanyaan terbuka]`**, bukan diputuskan di kode
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E19 - Perhitungan premi jiwa

**What to build:** Premi endorsement jiwa dihitung lewat **Seam 3 bentuk 1 yang terparameterisasi** —
bukan bentuk kelima.

`[terverifikasi]` Bentuknya tetap keluarga rumus dasar, tetapi dengan **tiga perbedaan
parameter**: basisnya nilai pertanggungan **setelah dikurangi bagian ceding**, rate diambil dari
**properti rate rata-rata jiwa**, dan ada **satu faktor pembebanan tambahan**.

⛔ **Properti bernama "rate" di rumus ini berperan sebagai faktor pembebanan**, bukan rate premi.
Membacanya sebagai rate premi akan salah total.

✅ Pembagi yang lebih kecil daripada lini lain adalah **konsekuensi K-018**, bukan rumus lain —
hanya satu faktor yang berkontribusi ke satuan.

**Asal (Pega).** `Activity\ReCountPremiLifeEDM` — penetapan nilai pertanggungan, nilai pertanggungan
liability, premi, dan nilai ter-spreading

**Keputusan.** K-051 · K-018 · K-010/K-012 · K-046

**Blocked by:** E12 · E18

**Status:** blocked

- [ ] Premi jiwa dihitung lewat **Seam 3**, sebagai **bentuk 1 terparameterisasi** — bukan bentuk baru
- [ ] Basis = nilai pertanggungan **liability**, bukan nilai pertanggungan penuh
- [ ] Rate dari properti rate rata-rata jiwa, **bukan** properti rate umum
- [ ] Pembagi diturunkan **resolver K-018** dari faktor yang ikut
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`
- [ ] **K-046** `K046_Life_RateSebagaiPembebanan` — properti "rate" adalah faktor pembebanan
- [ ] **K-046** `K046_Life_TreatyType1000036_SpreadingNol` — satu kode jenis treaty memaksa nilai ter-spreading menjadi nol; arti kodenya **`[pertanyaan terbuka]`**, diport apa adanya
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E20 - Skoring medis — Seam 6

**What to build:** Hasil pemeriksaan laboratorium dan usia tertanggung diubah menjadi **keputusan
underwriting**: diterima standar, ditunda, atau ditolak.

⛔ **Ini BUKAN perhitungan premi.** `[terverifikasi]` Berkas sumbernya **tidak menyentuh** rate,
nilai pertanggungan, premi, maupun pembagian — nol kemunculan keempatnya. Keluarannya **keputusan**,
bukan uang. Karena itu ia tidak masuk seam perhitungan mana pun.

⚠️ **Ambang klinisnya adalah aturan medis dan bisnis**, bukan konstanta teknis. **Diport apa
adanya** — jangan dibulatkan, jangan diseragamkan, jangan ditebak artinya.

⛔ **Domain sensitif.** Nilai ambang, hasil laboratorium, dan data medis **tidak pernah** disalin ke
tiket, test, fixture, log, maupun dokumen. Test memakai **nilai sintetis**.

⚠️ Skoring **risiko** yang bersama New Business adalah sistem **berbeda** dan **bukan** lingkup tiket
ini.

**Asal (Pega).** `Activity\CalculateScorLife_Act` · `SetParamLab_Act` · `CalculatePhysicalExam` ·
`SaveMedical` · `Harness\Medical_Harnes` · `Section\Medical_Sec`

**Keputusan.** **K-052** (Seam 6) · K-046 · `CLAUDE.md` §3.4, §3.5, §4.6

**Blocked by:** E18

**Status:** blocked

- [ ] **Seam 6 `services/underwriting.ScoreMedical` hidup**
- [ ] Keluaran berupa **keputusan** — diterima standar, ditunda, atau ditolak — **bukan angka uang**
- [ ] Skor yang bergantung **usia** bercabang pada rentang yang benar
- [ ] ⛔ Ambang klinis **diport apa adanya**; nilainya dibaca dari korpus, tidak ditebak
- [ ] ⛔ **Nilai ambang, hasil lab, dan data medis TIDAK disalin** ke artefak mana pun
- [ ] Test memakai **nilai sintetis** yang menguji batas, bukan data pasien
- [ ] Keluaran "ditolak" dapat muncul — perilakunya teruji sebagai **keputusan**, bukan galat
- [ ] **K-046** `K046_SkorMedis_AmbangKlinisDipertahankan`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

⚠️ **`[pertanyaan terbuka]`** — apa yang terjadi pada kasus berskor "ditolak" bila jiwa **melewati
tangga akseptasi** (E18). Milik work owner dan Underwriting.

## Endorsment Fac In - E21 - Jalur produksi endorsement ⛔ BLOCKED

**What to build:** Menuliskan selisih yang dihitung E13–E15 ke tabel produksi, beserta penomoran
versi polis.

⛔ **BLOCKED — menunggu tabel flat (P-10).** Prosedur tersimpan yang menulis dokumen polis di sistem
lama **diganti insert-ke-tabel-flat**, dan bentuknya dirancang **bersama New Business dan Renewal**,
bukan sendiri untuk endorsement.

⚠️ **Tiket ini tidak dihapus.** Lingkupnya sudah diketahui; menghilangkannya akan menyembunyikan
pekerjaan yang pasti ada.

`[terverifikasi]` Rantainya: penulisan dokumen polis → penyimpanan produksi → percabangan tujuh lini.
Ada **dua gerbang idempotensi** di lapis berbeda, sehingga rantai aman diulang — penting untuk
rekonsiliasi paralel run.

**Asal (Pega).** `Activity\SaveEDMToJsonPolicy_Act` · `SaveTreatyProduction_Act` ·
`SaveFacinProdAllEDM_Act` (tujuh cabang lini) · `InsertFacoutProductionEDM` · `GetEdmProdKe_Act`

**Keputusan.** K-046 · K-027 · K-050 · K-006 · `CLAUDE.md` §4.3

**Blocked by:** ⛔ **tabel flat (P-10)** · E13 · E15

**Status:** blocked

- [ ] ⛔ Menunggu rancangan tabel flat bersama New Business dan Renewal
- [ ] **Dua gerbang idempotensi** di lapis berbeda; rantai aman diulang
- [ ] Versi polis nol-berbasis, versi baru = terakhir + 1 — **cara hitung versi terakhir aman** karena tidak ada penghapusan baris (K-050)
- [ ] Skema Oracle **tidak berubah** (§4.3)
- [ ] Nilai berkoma desimal (K-027)
- [ ] **K-046** `K046_GuardFacOut_HanyaType7` — gerbang idempotensi fac out hanya dipakai keluar pada satu jenis penyesuaian
- [ ] **K-046** `K046_JsonPolisMonitoring_HanyaPrefiksEDMT` — tabel pemantauan hanya ditulis untuk satu prefiks kasus
- [ ] **K-046** `K046_DuaSemantikPRODKE` — dua cara menentukan versi terakhir hidup berdampingan
- [ ] ⚠️ Satu cabang lini bisnis rule-nya **tidak ada di korpus** tetapi ada di folder pelengkap work owner — kandidat amandemen **K-006**, bukan `panic`
- [ ] Endpoint dari konfigurasi, tidak pernah literal (§4.4)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Endorsment Fac In - E22 - Rekonsiliasi eksak siklus endorsement

**What to build:** Pembanding **nol-selisih** untuk kasus endorsement, memakai kerangka yang sudah
dibangun untuk New Business — **tanpa pembanding baru**.

⛔ **Seluruh kejanggalan yang diport apa adanya harus muncul sebagai selisih NOL**, bukan sebagai
selisih yang dimaafkan. Bila sistem baru "memperbaiki" salah satunya diam-diam, pembanding inilah
yang menangkapnya.

⚠️ **Satu pengecualian yang sah:** perbaikan A.5 mengubah tipe nilai cadangan premi Nusantara Re dari
string menjadi angka. Pembandingan mentah akan menunjukkan **beda tipe meski nilainya sama**.
Pembanding **menormalkan** sebelum membandingkan, dan selisih semacam ini **dijelaskan oleh A.5** —
bukan ditandai cacat.

⛔ **Tidak pernah menyentuh berkas mentah ber-PII.** Fixture yang dipakai adalah yang sudah
ter-de-identifikasi.

⚠️ Tiket ini **tidak bersandar** pada klaim arsip bahwa nilai dasar akseptasi adalah selisih TSI —
klaim itu tetap **`[belum diuji]`**.

**Asal (Pega).** — (pembanding, bukan port rule)

**Keputusan.** K-046 · A.5 · ADR-F-0001 · K-025

**Blocked by:** `..\15-de-identifikasi-berkas-kasus.md` · `..\16-rekonsiliasi-eksak-tahap-1.md` ·
E09 · E14 · E15 · E19
⛔ Kedua tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] Memakai **kerangka rekonsiliasi New Business**, tanpa pembanding baru
- [ ] ⛔ **Tidak pernah** membaca berkas mentah ber-PII; hanya fixture ter-de-identifikasi
- [ ] Seluruh kejanggalan K-046 menghasilkan **selisih nol**
- [ ] Pembanding **menormalkan** string dan angka; selisih tipe dari **A.5** dijelaskan, bukan ditandai cacat
- [ ] Mencakup ketiga lapis before-image, selisih pembayaran, selisih baris spreading, dan premi jiwa
- [ ] Setiap selisih yang tersisa **dapat ditelusuri** ke rule Pega asalnya (§4.6)
- [ ] Nilai medis, nomor polis, dan nama orang **tidak pernah** muncul di laporan pembanding

# Fac Out

Jumlah tiket: **14**

## Fac Out - F01 - Predikat Fac Out masuk registry

**What to build:** Sistem dapat menjawab dua pertanyaan tentang sebuah kasus — *apakah kasus ini punya
retrosesi keluar?* dan *apakah baris spreading ini membawa penanda Fac Out?* — lewat **registry
predikat yang sudah ada**, tanpa registry kedua dan tanpa seam baru.

⛔ **`IsFacout` TIDAK ditangani di sini.** `[terverifikasi]` Ia **tidak terpanggil sama sekali** di
jalur Fac Out: `NB FacIn\Flow\InputInwardFacultativeOffer.xml` memuatnya **0 kali**, peka huruf maupun
tidak. *(Berkas flow itu `[terverifikasi]` hanya ada di folder **NB**; tidak ada di RNW maupun
Endorsment.)*
Isinya menguji **banding** (`ProposalAcceptStatus = 4`, K-029), bukan fac out. Porting-nya sudah
menjadi lingkup tiket New Business `..\09-predikat-sikap-khusus.md` di bawah K-019 — **jangan
diulang**.

`[terverifikasi]` Dua predikat yang masuk lingkup, beserta kemunculannya di flow: `IsFacRetro` **17×**,
`IsInputFacRetro` **9×**.

##### ⛔ `IsFacRetro` adalah **DUA rule di kelas berbeda** — catatan asal-usul `IsUW` hanya berlaku untuk satu

> ⛔ **Dipersempit 21 September 2026.** Rumusan sebelumnya menuliskan asal-usul salinan `IsUW` seolah
> berlaku **umum** untuk `IsFacRetro`. Itu **hanya benar untuk satu dari dua rule**. Rumusan lama
> tidak dihapus — ia **dibatasi lingkupnya** (`PANDUAN-KERJA` §7).

`[terverifikasi]` Identitas rule ditentukan **basis `pzInsKey` + `pyClassName`**, bukan nama berkas
(`_ARSIP-lintas-siklus\_BACA-INI.md` Koreksi R1). `When\IsFacRetro` berwujud **dua rule**:

| | Varian **OfferFacIn** | Varian **Work** |
| --- | --- | --- |
| `pyClassName` | `ASM-FW-GISFW-Data-OfferFacIn` | `ASM-FW-GISFW-Work` |
| Terekspor di | **NB** dan **RNW** | **Endorsment** |
| `<pyLogic>` `A` menguji | `.IsFacRetro = 1` | `pyWorkPage.OfferFacIn.IsFacRetro = 1` |
| `<pyNestedConditions>` | ⛔ `rowdata(1)` **masih memuat sisa** pemeriksaan workbasket (`pxRequestor…pyWorkBasketName = ReasFacInGroupLeader`) | ✅ **BERSIH** — tidak ada sisa |
| `<pzOriginalInstanceKey>` | menunjuk rule **`ISUW`** | — |

⚠️ **Asal-usul salinan `IsUW` HANYA berlaku pada varian `ASM-FW-GISFW-Data-OfferFacIn`.** Di varian
itu, yang mengikat adalah **ekspresi tersimpan**, bukan grid kondisinya: yang dieksekusi
`<pyLogic>` = `A` → `<pyConditionValue1>` = `compareTwoValues(.IsFacRetro,"=",1)`; grid kondisi
warisan **diabaikan**. Sikapnya **sama persis dengan `IsOfferFacIn`** (K-002, tiket NB
`..\09-predikat-sikap-khusus.md`).

⛔ **Varian `ASM-FW-GISFW-Work` TIDAK memuat sisa itu** dan **tidak perlu perlakuan khusus** —
kondisinya bersih. Memperlakukan kedua varian sama akan memasang penanganan warisan pada rule yang
tidak memerlukannya.

`[dugaan]` **Varian `Work` kemungkinan yang dipanggil flow.** Indeks rujukan
`NB FacIn\Flow\InputInwardFacultativeOffer.xml` menyebut kelas **`ASM-FW-GISFW-Work`**. ⚠️ Indeks merekam resolusi
**saat berkas terakhir disimpan** — ia **bukan bukti** rule mana yang resolve saat runtime, sehingga
kata "kemungkinan" **tetap `[dugaan]`**. **Keduanya diport** (pola rule kembar, `..\..\10-audit\10-rule-kembar-pzinskey.md`).

✅ **K-057 TIDAK terpengaruh.** Kedua varian menguji **flag yang sama** (`.IsFacRetro = 1`) pada
**halaman yang sama** (`OfferFacIn`) — hanya jalur penulisan propertinya berbeda (relatif vs lewat
`pyWorkPage`). Pemicu Fac Out tetap seperti ditetapkan K-057.

⛔ **Tiga bentuk penanda spreading TIDAK DIGABUNG** (K-054):

| Bentuk | Ekspresi | Di mana |
| --- | --- | --- |
| **ketat** | `.TreatyType=="10015"` | `SetDataFacOutFire_Act` · `CopyFacRetroFire_ACT` |
| **contains tunggal** | `@contains(.TreatyType,"10015")` | di dalam `CopyToAllSpreading_ACT` · `CopyToAllLocSpreading_ACT` |
| **umum tiga-arah** | `@contains(.TreatyType,"10015")\|\|@contains(.TreatyName,"SPL")\|\|@contains(.TreatyType,"10007")` | `CopyToAllSpreading_ACT` · `CopyToAllLocSpreading_ACT` |

⚠️ **`"SPL"` diuji pada `TreatyName`, bukan `TreatyType`.** Menyamakan ketiga bidang uji akan menarik
baris spreading yang tidak seharusnya ikut.

##### ⛔ Empat properti bernama mirip — registry harus membedakannya

`[terverifikasi]` Registry **tidak boleh** memperlakukan keempatnya sebagai satu:

| Properti | Gaya nilai | Perannya |
| --- | --- | --- |
| `OfferFacIn.IsFacRetro` | `0`/`1` telanjang | **pemicu Fac Out** — dibaca `When\IsFacRetro`, menjadi `<pyTaskWhen>` di flow |
| `OfferFacIn.IsInputFacRetro` | `0`/`1` | gerbang **re-input** saat UW — dibaca `When\IsInputFacRetro` |
| `pyWorkPage.IsInFacRetro` | dibanding `!= 1` | **kondisi tampilan** Section; **bukan** rule `When` |
| `.IsFacRetroOffer` | **`"0"` berkutip** | penanda penawaran retro; disetel `OfferFacOut_PreAct` |

⚠️ Hanya **dua** di antaranya yang berupa rule `When` dan masuk registry. `pyWorkPage.IsInFacRetro`
dan `.IsFacRetroOffer` adalah **properti biasa** — memasukkannya ke registry adalah kekeliruan
kategori. Rincian di `..\..\10-audit\03-struktur-facretro-dan-populasi.md` dan tiket
`F05-empat-struktur-staging.md`.

**Asal (Pega).** `When\IsFacRetro` · `When\IsInputFacRetro` · penanda di
`Activity\SetDataFacOutFire_Act` · `CopyFacRetroFire_ACT` · `CopyToAllSpreading_ACT` ·
`CopyToAllLocSpreading_ACT`

**Keputusan.** **K-054** · K-050 (satu registry, sumber data per-rule) · K-002 · K-019 ·
**K-057** (pemicu = flag `.IsFacRetro`; rumusan `IsFacout` **dicabut**) · `CLAUDE.md` §4.5, §4.6

**Blocked by:** `..\08-registry-rules-eval.md` · `..\09-predikat-sikap-khusus.md`
⛔ Keduanya **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] `IsFacRetro` dan `IsInputFacRetro` hidup lewat **Seam 1 `rules.Eval`** — tanpa registry kedua
- [ ] `IsFacRetro` memakai **ekspresi tersimpan** `.IsFacRetro = 1`; grid kondisi warisan **diabaikan**
- [ ] ⛔ **DUA rule `IsFacRetro` diport** — kelas `Data-OfferFacIn` **dan** `GISFW-Work`; komentar §4.6 menyebut **kelasnya**, nama saja tidak cukup
- [ ] Komentar menyebut asal-usul salinan `IsUW` **hanya pada varian `Data-OfferFacIn`** — varian `GISFW-Work` kondisinya bersih dan **tidak** diberi penanganan warisan
- [ ] **Tiga bentuk penanda spreading terpisah** dan tidak dapat tertukar
- [ ] `"SPL"` diuji pada **`TreatyName`**; kedua kode lain pada `TreatyType`
- [ ] ⛔ `IsFacout` **tidak diimplementasikan di sini** — komentar menunjuk `..\09-predikat-sikap-khusus.md`
- [ ] `CoverageBasis == 5` dikenali sebagai **Layering Basis** `[terverifikasi]` `DDL\CoverageBasis.xml`
- [ ] Arti `10015` = FACOUT dan `10007` = ORS ditandai **keterangan work owner (K-054)**, bukan terverifikasi korpus
- [ ] `[pertanyaan terbuka]` arti `TreatyName == "SPL"` **tetap terbuka**; nilainya diport apa adanya
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F02 - Tracer — derivasi Fac Retro saat UW Accept (lini Fire)

**What to build:** Saat Underwriter menekan **Accept**, lokasi yang punya spreading Fac Out otomatis
tersalin ke daftar retrosesi, kasus ditandai sebagai punya Fac Out, dan alur membelok ke menu Fac Out
**sebelum** produksi. Tiket ini menembus seluruh lapisan untuk **satu lini bisnis — Fire** — sebagai
tracer.

⛔ **PEMICUNYA BUKAN `IsFacout`, DAN BUKAN `ProposalAcceptStatus = 4`.** Rumusan lama salah pada dua
lapis dan sudah dibatalkan (`09-facout\01` §0.1):

1. `[terverifikasi]` `IsFacout` muncul **0 kali** di `NB FacIn\Flow\InputInwardFacultativeOffer.xml`
   (`[terverifikasi]` berkas itu hanya ada di folder **NB**). Yang ada
   di `SetValidateDateUW_PostAct` hanyalah **variabel lokal** bernama `isFacOut` — tag
   `<pyLocalParameters>` → `<pyParametersParamName>`, bertipe `String`. Nama variabel lokal yang mirip
   nama rule; **bukan** pemanggilan rule.
2. `[terverifikasi]` `ProposalAcceptStatus = 4` berarti **Banding**, bukan Accept (**K-029**,
   `GLOSARIUM.md`).

`[terverifikasi]` **Rantai yang sebenarnya**, urutan dibaca dari `<pyStepPageReference>` — bukan dari
nomor baris:

| # | Berkas | Tag pembawa | Isi |
| ---: | --- | --- | --- |
| 1 | `Activity\SetValidateDateUW_PostAct` | `pyStepPageReference` = `RH_1.pySteps(3)` | `Call SetDataFacOut_Act`, **TANPA precondition** |
| 2 | `Activity\SetDataFacOut_Act` | `pyStepsActivityName` · `PropertiesName` | `Property-Remove` `pyWorkPage.OfferFacIn.FacRetroList` · `Property-Set` `.OfferFacIn.IsFacRetro = 0` (**reset**) · langkah 3–6 panggil per COB |
| 3 | `Activity\SetDataFacOutFire_Act` | `pyStepsPreCondParamsWhen` · `PropertiesName` | `.OfferFacIn.IsFacRetro = 1`, bergerbang `.TreatyType=="10015"` / `Local.isFacOut==1` / `Local.isFacOutLoc==1` |
| 4 | `Flow\InputInwardFacultativeOffer` | `pyTaskStatusOrWhen`=`WHEN` · `pyTaskWhen`=`IsFacRetro` · `pyLikelihood`=100 | transisi **`Transition92`** → shape `pyMOName` = `[Fac Out]` |
| 5 | shape `[Fac Out]` | **`pyImplementation`** = `OfferFacRetro` | ⚠️ **bukan** `pySubFlowName` |

📌 **Penyaringan terjadi di precondition per COB, bukan di titik masuk.** Titik masuk tidak
bergerbang sama sekali — `SetDataFacOut_Act` selalu dipanggil, lalu mereset flag ke `0`, dan hanya
sub-aktivitas per lini yang menyalakannya ke `1`.

`[terverifikasi]` Penanda penaut ke objek Fac In induk **lahir di sini**, bukan saat produksi:
`SetDataFacOutFire_Act` menyetel `FacRetro.LocationList(<LAST>).Property.ObjectNo` dan
`.Property.ObjectNoFacIn`.

**Asal (Pega).** `Activity\SetValidateDateUW_PostAct` `RH_1.pySteps(3)` · `Activity\SetDataFacOut_Act` ·
`Activity\SetDataFacOutFire_Act` · `Flow\InputInwardFacultativeOffer` `Transition92`

**Keputusan.** **K-057** (pemicu = flag `.IsFacRetro`; rumusan lama **dicabut**) · K-053 · K-054 · K-019 · K-029 ·
K-010/K-012 · `CLAUDE.md` §4.6

**Blocked by:** F01 · `..\01-tracer-money-ratio-premi-pa.md` · `..\11-tangga-akseptasi-bentuk-a.md`
⛔ Kedua tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] Derivasi dipanggil dari jalur UW Accept **tanpa gerbang di titik masuk**
- [ ] Langkah reset menghapus daftar retrosesi **dan** menihilkan flag sebelum pengisian ulang
- [ ] Flag `.IsFacRetro` menyala **per lini bisnis**, hanya bila penanda spreading ditemukan
- [ ] Flag itu menjadi **gerbang belok** ke menu Fac Out — perilaku `Transition92`
- [ ] ⛔ `IsFacout` **tidak dipakai** di jalur ini; komentar menyatakan sebabnya
- [ ] ⛔ `ProposalAcceptStatus` **tidak menggerbangi** jalur ini
- [ ] Penaut `ObjectNo` + `ObjectNoFacIn` tercipta **saat derivasi**, bukan saat produksi
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`
- [ ] Demo: satu kasus Fire dengan spreading ber-penanda → daftar retrosesi terisi, flag menyala, alur membelok
- [ ] Demo tandingan: kasus Fire **tanpa** spreading ber-penanda → flag tetap `0`, alur **tidak** membelok
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F03 - Daftar pengecualian okupasi — 26 alternatif

**What to build:** Sebagian kode okupasi menandai risiko sebagai **pengecualian**, dan penandaan itu
ikut berjalan di jalur derivasi Fac Out lini Fire.

`[terverifikasi]` Gerbangnya **satu ekspresi dengan 26 alternatif `||`**: **20** memakai
`@startsWith` dan **6** memakai kesetaraan ketat. Keluarannya menyetel
`pyWorkPage.OfferFacIn.IsOccupException`.

```powershell
### menghasilkan 26 / 20 / 6
$c=[IO.File]::ReadAllText("D:\migrasi\RNM\NB FacIn\Activity\SetDataFacOutFire_Act.xml")
$occ=[System.Net.WebUtility]::HtmlDecode(
  [regex]::Match($c,'<pyStepsPreCondParamsWhen>([^<]*OccupationId[^<]*)</pyStepsPreCondParamsWhen>').Groups[1].Value)
($occ -split '\|\|').Count
([regex]::Matches($occ,'@startsWith')).Count
([regex]::Matches($occ,'\.OccupationId==')).Count
```

⛔ **Gaya penulisannya tidak konsisten, dan itu DIPORT APA ADANYA.** `[terverifikasi]` Sebagian kode
diuji sebagai **awalan** (`251`, `252`, `257`, `258`, `259`), sebagian lagi sebagai **kesetaraan
penuh** (`256`) — padahal berada di rentang yang sama. Menyeragamkannya **mengubah himpunan kode yang
tertangkap**, dan itu perubahan perilaku (`CLAUDE.md` §1).

⚠️ **Awalan yang lebih pendek menelan yang lebih panjang.** `[terverifikasi]` Daftar memuat awalan
`26` dan `28` berdampingan dengan kesetaraan penuh pada kode empat-lima digit. Urutan dan bentuk
ujinya **dipertahankan persis**, tidak diringkas menjadi rentang.

**Asal (Pega).** `Activity\SetDataFacOutFire_Act` — `pyStepsPreCondParamsWhen` atas `.OccupationId`;
target `pyWorkPage.OfferFacIn.IsOccupException`

**Keputusan.** K-046 (pola kode usang diport apa adanya) · K-054 · `CLAUDE.md` §1, §4.6

**Blocked by:** F02

**Status:** blocked

- [ ] Daftar **26 alternatif** lengkap, tidak satu pun dihilangkan
- [ ] **20 uji awalan** dan **6 uji kesetaraan penuh** tetap sebagai bentuk aslinya
- [ ] ⛔ Daftar **tidak diringkas** menjadi rentang, regex, atau tabel lookup
- [ ] Kode okupasi disimpan sebagai **string**, bukan angka — uji awalan menuntutnya
- [ ] Hasilnya menyetel penanda pengecualian okupasi, bukan menggerakkan alur
- [ ] Test menguji **batas**: kode yang persis cocok, yang cocok sebagai awalan, dan yang di luar daftar
- [ ] **K-046** `K046_OkupasiFacOut_GayaUjiTidakSeragam` — awalan dan kesetaraan bercampur, sengaja dipertahankan
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F04 - Derivasi Fac Retro untuk tiga lini bisnis sisanya

**What to build:** Derivasi retrosesi yang sudah hidup untuk Fire (F02) diperluas ke **Aneka/Golf**,
**Marine Cargo/MBU**, dan **PA/Travel**, sehingga seluruh lini bisnis dapat menurunkan Fac Out.

`[terverifikasi]` `Activity\SetDataFacOut_Act` memanggil **empat** sub-aktivitas per lini —
`SetDataFacOutFire_Act` (F02) plus `SetDataFacOutAnekaGolf_Act`, `SetDataFacOutCargoMBU_Act`,
`SetDataFacOutPATravel_Act` — masing-masing bergerbang predikat lini bisnisnya sendiri.

⚠️ **Tiap lini menulis ke koleksi anak yang berbeda.** `[terverifikasi]` `FacRetroList` punya **enam**
koleksi anak, bukan satu: `LocationList` · `CargoList` · `PersonList` · `VehicleList` ·
`Property.RiskLocation.AnekaList` · `Property.RiskLocation.OccupationList(n).AnekaList`. Menyamakan
keenamnya ke satu daftar akan kehilangan data.

`[terverifikasi]` **Penjumlah PA berdiri sendiri.** `Activity\SumFacOutPA_Act` bergerbang
`.TreatyType=="10015"` dan `IsPA`, menyetel `.PremiumSpreaded` dan `.TSISpreaded`. Ia **bukan**
bagian dari `SetDataFacOutPATravel_Act` dan tidak boleh dilebur ke dalamnya.

⚠️ **Dua dari empat sub-aktivitas berbeda isi di Endorsement.** `[terverifikasi]` kontrak 23 tag:
`SetDataFacOutFire_Act` dan `SetDataFacOutAnekaGolf_Act` **berbeda** NB↔EDM; `SetDataFacOutCargoMBU_Act`
dan `SetDataFacOutPATravel_Act` **identik**. Perbedaannya ditangani **F11**, bukan di sini.

**Asal (Pega).** `Activity\SetDataFacOutAnekaGolf_Act` · `SetDataFacOutCargoMBU_Act` ·
`SetDataFacOutPATravel_Act` · `SumFacOutPA_Act` · induk `SetDataFacOut_Act`

**Keputusan.** K-053 · K-054 · K-018 (skala rasio per lini) · K-010/K-012 · `CLAUDE.md` §4.6

**Blocked by:** F02

**Status:** blocked

- [ ] Keempat lini menurunkan Fac Out lewat **satu jalur derivasi yang sama**, bercabang di gerbang lini
- [ ] **Enam koleksi anak** `FacRetroList` terisi sesuai lininya, tidak dileburkan
- [ ] Penjumlah PA tetap **langkah tersendiri**, bergerbang penanda ketat `"10015"` **dan** lini PA
- [ ] Skala rasio mengikuti **resolver K-018**, bukan angka hafalan per lini
- [ ] Lini yang tidak ada di peta resolver → **`panic`** (`CLAUDE.md` §4.5)
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`
- [ ] Perbedaan Endorsement pada dua sub-aktivitas **tidak ditangani di sini** — komentar menunjuk F11
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F05 - Empat struktur penampung Fac Retro — dipisahkan, bukan disatukan

**What to build:** Data retrosesi hidup di **empat** struktur berbeda sepanjang satu kasus. Tiket ini
memodelkan keempatnya beserta perannya, sehingga tidak ada yang tertukar.

⛔ **EMPAT, bukan dua.** Catatan lama yang menyebut "dua struktur" (`09-facout\01` §2.1) **tidak
lengkap**:

| # | Struktur | Bentuk | Diisi oleh |
| ---: | --- | --- | --- |
| 1 | `OfferFacIn.FacRetro` | **tunggal** (staging) | `SetDataFacOutFire_Act` — `FacRetro.LocationList(<APPEND>)` |
| 2 | `OfferFacIn.FacRetroList` | **list berindeks** | `CopyFacRetroFire_ACT` — `FacRetroList(counterFacRetro)…` |
| 3 | `OfferFacIn.FacRetroDetails` | tunggal, 4 properti | `OfferFacOut_PreAct` · `GetRISlipDataFromDB_Act` · `SetPropertyToOfferFacIn_Act` |
| 4 | `pyWorkPage.Policy.FacOfferList(n)` | list berindeks | sumber salinan bagi struktur 3 |

`[terverifikasi]` **`FacRetroDetails` punya tepat empat properti**, seluruhnya dari tag
`<PropertiesName>` / `<PropertiesValue>` / `<pyValue>` / `<pyStepsPreCondParamsWhen>`:

| Properti | Ditulis | Dibaca |
| --- | --- | --- |
| `.BackUpStatus` | `GetRISlipDataFromDB_Act` | `CheckProtectFacout_Act` (`==3` **telanjang**) · `OfferFacOut_PostAct` (`=="3"` **berkutip**) |
| `.FacNo` | `OfferFacOut_PreAct` | — |
| `.AdditionalInfo` | `OfferFacOut_PreAct` | — |
| `.DocumentPosition` | `OfferFacOut_PreAct` · `SetPropertyToOfferFacIn_Act` | `Harness\ViewLetter` · `Section\FacOutPrintRISlipSectionInside` · `Section\FacultativeLetter` |

⚠️ **Nilai yang sama dibandingkan dua gaya** — `==3` telanjang dan `=="3"` berkutip. Angka
diperlakukan sebagai string di satu tempat dan sebagai angka di tempat lain. **Kandidat K-046, diport
apa adanya.**

##### Anak langsung `FacRetroList` — **19**, bukan enam

> ⛔ **KOREKSI 21 September 2026.** Tiket ini semula menulis ~~"`FacRetroList` punya **enam** koleksi
> anak: `LocationList` · `CargoList` · `PersonList` · `VehicleList` ·
> `Property.RiskLocation.AnekaList` · `Property.RiskLocation.OccupationList(n).AnekaList`"~~.
> **Salah pada dua hal:** dua yang terakhir **bukan anak langsung** (letaknya bersarang di bawah
> `LocationList(n).Property.RiskLocation`), dan daftarnya **melewatkan `CurrencyList`**.

`[terverifikasi]` `FacRetroList(n)` punya **19 anak langsung** — **5 koleksi** dan **14 non-koleksi**:

| Koleksi (`*List`) — 5 | Non-koleksi — 14 |
| --- | --- |
| `CargoList` · **`CurrencyList`** · `LocationList` · `PersonList` · `VehicleList` | `EndPeriod` · `OurRef` · `PctPremiAllObj` · `PctPremiAllObjUSD` · `PCTPremiIDR` · `PctShareAllObj` · `PrintRISlip` · `ReinsurerID` · `ReinsurerName` · `RiCommAllObj` · `StartPeriod` · `TFAllObj` · `UjrahAllObj` · `pyExpanded` |

⚠️ `pyExpanded` adalah properti tampilan Pega, bukan data bisnis — dicatat agar sensusnya dapat
direproduksi.

📌 Lini Aneka memang dijangkau lewat `LocationList(n).Property.RiskLocation.AnekaList` dan
`…OccupationList(o).AnekaList`, tetapi itu **jalur bersarang**, bukan anak `FacRetroList`.

**Sumber angka:** `..\..\10-audit\03-struktur-facretro-dan-populasi.md` §2.1.

`[terverifikasi]` **`FacOutTSI` melekat di coverage, bukan di `FacOutObjectList`.** Seluruh jalur
berbentuk `.CoverageList(1).FacOutTSI` atau saudara PA-nya `.ASMCoverage(1).FacOutTSI`; **nol** jalur
berbentuk `FacOutObjectList(…).FacOutTSI`. Anggota `FacOutObjectList(1)` yang sebenarnya: `Rate`,
`ShareOffered`, `PercentOffered`, `ObjectPremi`, `RiComm`, `commision`, `TF`, `Ujrah`.

⛔ **`RiComm` di `FacOutObjectList` BUKAN `.RIComm` di coverage.** Beda kapitalisasi, beda properti,
beda pemilik. Begitu pula `commision` (ejaan pendek) terhadap `COMMISION` di tabel produksi.

⚠️ **Gaya nilai flag berbeda antar properti.** `[terverifikasi]` `.IsFacRetroOffer` disetel **`"0"`
berkutip** di `OfferFacOut_PreAct`, sedangkan `.IsFacRetro` memakai `0`/`1` **telanjang**. Properti
berbeda, gaya berbeda — **diport apa adanya**.

⚠️ `[terverifikasi]` Gerbang `CopyFacRetroFire_ACT` menguji **`local.Facout=="0"`** — angka
dibandingkan **sebagai string**. Kandidat K-046, diport apa adanya.

📌 `[terverifikasi]` **Nol kemunculan `FacRetroDetails` di `DDL\`** → ia struktur clipboard/JSON,
**bukan** tabel Oracle. Tidak ada DDL yang perlu dicari untuknya.

##### ⛔ Empat properti bernama mirip — **TIDAK BOLEH disamakan**

`[terverifikasi]` Empat properti berbeda dengan nama yang nyaris sama hidup berdampingan. Menyamakan
dua di antaranya akan menggeser gerbang alur:

| Properti | Nilai & gaya | Ditulis / dibaca di | Perannya |
| --- | --- | --- | --- |
| `OfferFacIn.IsFacRetro` | `0` / `1` **telanjang** | disetel `0` di `SetDataFacOut_Act` (`RH_1.pySteps(3)`, `<PropertiesName>`); disetel `1` di **keempat** `SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act` | **pemicu Fac Out** — menjadi `<pyTaskWhen>` transisi di flow masuk siklus ybs. |
| `OfferFacIn.IsInputFacRetro` | `0` / `1` | `<pyStepsPreCondParamsWhen>` `Param.Status=="UW" && …IsInputFacRetro==1` | gerbang **re-input** saat UW; **9** kemunculan di `NB FacIn\Flow\InputInwardFacultativeOffer` |
| `pyWorkPage.IsInFacRetro` | dibanding `!= 1` | `<pyContainerVisibleWhen>` di `Section\InputCoverageFire` | **kondisi tampilan** layar, bukan alur |
| `.IsFacRetroOffer` | **`"0"` berkutip** | disetel di `Activity\OfferFacOut_PreAct` (`<PropertiesName>`) | penanda penawaran retro |

⚠️ **Gayanya berbeda dan itu disengaja dipertahankan:** `.IsFacRetro` memakai angka telanjang,
`.IsFacRetroOffer` memakai string berkutip. **Diport apa adanya** (K-046).

⚠️ `[terverifikasi]` Reset memakai jalur **relatif** `.OfferFacIn.IsFacRetro`; penyalaan memakai jalur
**absolut** `pyWorkPage.OfferFacIn.IsFacRetro`. Dua bentuk penulisan untuk properti yang sama.

**Sumber angka:** `..\..\10-audit\03-struktur-facretro-dan-populasi.md` ·
`..\..\10-audit\06-tinjau-ulang-suntingan-facout.md` §2.5.

**Asal (Pega).** `Activity\SetDataFacOutFire_Act` · `CopyFacRetroFire_ACT` · `CopyFacRetro_ACT` ·
keluarga `CopyAllObjFacOut{FireAneka,GolfCargo,PAMBU}_ACT` · `OfferFacOut_PreAct` ·
`GetRISlipDataFromDB_Act` · `SetPropertyToOfferFacIn_Act`

**Keputusan.** K-053 · **K-046** · K-010/K-012 · `CLAUDE.md` §1, §4.6

**Blocked by:** F02

**Status:** blocked

- [ ] **Empat struktur** dimodelkan terpisah; tidak ada yang dilebur
- [ ] `FacRetroDetails` punya **tepat empat properti**, tidak lebih
- [ ] **Enam koleksi anak** `FacRetroList` terwakili
- [ ] `FacOutTSI` melekat di **coverage**, bukan di dalam `FacOutObjectList`
- [ ] `RiComm` (obyek fac out) dan `.RIComm` (coverage) adalah **dua properti berbeda** yang tidak dapat tertukar
- [ ] **K-046** `K046_BackUpStatus_DuaGayaPembanding` — `==3` dan `=="3"` berdampingan
- [ ] **K-046** `K046_LocalFacout_AngkaSebagaiString` — `local.Facout=="0"`
- [ ] **K-046** `K046_IsFacRetroOffer_NolBerkutip` — `"0"` berkutip vs `0` telanjang
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F06 - Perhitungan premi retrosesi

**What to build:** Premi yang diteruskan ke reasuradur retro dihitung per coverage, lengkap dengan
komisi RI dan potongan diskon, sehingga nilainya siap ditulis ke produksi Fac Out.

#### ⛔ DUA rule, bukan satu — dipilih menurut COB

**K-060 (amandemen):** `CountRateRetroCov` adalah **dua rule berbeda di kelas berbeda**, dipilih
menurut **lini bisnis**. **Keduanya berlaku dan keduanya harus diport** — bukan dipilih salah satu.

⛔ **Identitas rule = `pzInsKey` + `pyClassName`, BUKAN nama berkas.**

`[terverifikasi]`

| | Varian **FIRE** | Varian **ANEKA** |
| --- | --- | --- |
| `pyClassName` | `ASM-FW-GISFW-Data-PropertyItem` | `ASM-FW-GISFW-Data-Aneka` |
| Basis `pzInsKey` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-PROPERTYITEM COUNTRATERETROCOV` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-ANEKA COUNTRATERETROCOV` |
| **Berkas berlaku** | **`DDL\CountRateRetroCov.xml`** | **`DDL\CountRateRetroCov(ANEKA).xml`** |
| Salinan korpus | NB · RNW — ⛔ **BASI** | Endorsment — ✅ **mutakhir** (nol beda) |
| Langkah | **14** (10 puncak + 4 bersarang) | **12** (9 puncak + 3 bersarang) |

`[terverifikasi]` `pzOriginalInstanceKey` keduanya **sama** — varian Aneka lahir dari menyalin varian
PropertyItem lalu menyimpang.

##### Perbedaan yang mengubah angka

| Aspek | **FIRE** | **ANEKA** |
| --- | --- | --- |
| **Pembagi premi** | **100.000** | ⛔ **10.000** |
| **Uji mata uang** | `.Currency=="IDR"` · `=="USD"` | ⛔ `.Currency.Name=="IDR"` · `=="USD"` |
| `Local.RIComIN` diakumulasi | ✅ ada | ⛔ **tidak ada** |
| `Local.TotalPremiCov` | ada | ⛔ tidak ada |
| Rumus koreksi langkah 10 | ✅ ada | ⛔ **tidak ada** |
| Gerbang `pySteps(2)`/`(3)`/`(4)` | `IsEDM` / `IsEdmExtendPeriod` / `IsEDM` | **identik** |

⛔ **Pembagi berbeda satu orde.** Memakai satu pembagi untuk kedua jalur menghasilkan premi retro
**10× salah** pada salah satunya.

📌 **Tabel langkah lengkap kedua varian ada di
`..\..\10-audit\01-pohon-langkah-countrateretrocov.md` §5 (FIRE) dan §7A (ANEKA); banding
berdampingan di §7B.**

#### Rumus inti — varian FIRE

```
.PremiumRetro     = @Math.divide((ShareOffered * .Rate * Local.prorate), 100000, 20)
.PremiumRetro     = .PremiumRetro - @Math.divide((.PremiumRetro * .DiscountPercentage), 100, 20)
.RICommPercentage = @Math.divide(Local.RIComIN, Local.LengCov, 20)
.RIComm           = @Math.divide((.PremiumRetro * .RICommPercentage), 100, 20)
.Rate             = @Math.divide(Local.RateOut, Local.LengCov, 20)
```

`[terverifikasi]` **Kedua rata-rata itu sejajar.** `Local.RateIN` **dan** `Local.RIComIN` sama-sama
**diakumulasi per coverage** di `RH_2.pySteps(8)`, lalu sama-sama dibagi `Local.LengCov`:

```
RH_2.pySteps(8)  (Property-Set, iterasi .CoverageList, tanpa gerbang)
  [1] Local.RateIN  := Local.RateIN + .Rate
  [2] Local.LengCov := .pxListSubscript
  [3] Local.RIComIN := Local.RIComIN + .RIComm
```

⛔ **URUTAN OPERASI DIKUNCI: diskon dikurangkan LEBIH DULU, `.RIComm` dihitung SESUDAHNYA.**
`[terverifikasi]` Keduanya berada di langkah **berbeda** — diskon di `RH_2.pySteps(9).pySteps(2)`,
komisi di `RH_2.pySteps(9).pySteps(3)` — dan **keduanya tanpa gerbang**, sehingga tidak ada percabangan
yang dapat membalik urutannya. Menghitung `.RIComm` dari premi **sebelum** diskon menghasilkan angka
berbeda.

⚠️ **Gerbang langkah 3 berbeda antara kedua salinan.** `[terverifikasi]` Salinan **DDL yang berlaku**
menggerbangi `Local.prorate := 100` **hanya** dengan `IsEdmExtendPeriod`; salinan korpus yang basi
menuntut `IsEDM` **dan** `IsEdmExtendPeriod`. **Ikuti DDL.** Perincian di audit 01 §4.2.

✅ Pembagi **100.000** = `1.000 × 100` adalah konsekuensi **aturan pembagi komposit K-018**, bukan
rumus lain.

#### `Local.prorate` — enam penugasan, tetapi **nilai akhirnya sudah pasti**

> ⛔ **Keputusan work owner 21 September 2026.** Pada `CountRateRetroCov`, `Local.prorate` **bernilai
> `1`** keluar dari langkah 1. **Dua penugasan sebelumnya di langkah yang sama TIDAK berlaku.**
>
> Rumusan lama di tiket ini menyajikan keenam penugasan sebagai sama-sama mungkin dan menggantungkan
> jawabannya pada `REPEATINGINDEX`. **Digantikan bagian ini**; tidak dihapus dari riwayat
> (`PANDUAN-KERJA` §7).

`[terverifikasi]` Keempat langkah prorata **identik di kedua varian** — gerbang, T/F dan penugasannya
sama persis di `DDL\CountRateRetroCov.xml` dan `DDL\CountRateRetroCov(ANEKA).xml`, diukur dengan
parser yang membedakan **anak langsung** dari keturunan:

| Langkah | Gerbang | T | F | Penugasan |
| --- | --- | ---: | ---: | --- |
| `pySteps(1)` | — | 2 | 2 | `[1]` `EndPeriod − StartPeriod` ⛔ **mati**<br>`[2]` `@divide(prorate, @if(EDMDay=="",365,@toDecimal(EDMDay)), 20) * 100` ⛔ **mati**<br>`[3]` **`1`** ✅ berlaku |
| `pySteps(2)` | `IsEDM` | 2 | 3 | `(ProrateEDMEnd + ProrateStartEDM) * 100` |
| `pySteps(3)` | **`IsEdmExtendPeriod` saja** — ⛔ salinan korpus yang basi menuntut `IsEDM` **dan** `IsEdmExtendPeriod` | 2 | 3 | `100` |
| `pySteps(4)` | `IsEDM` | **3** | **2** | `100` — ⛔ **terbalik**: jalan saat `IsEDM` **SALAH** |

##### Nilai akhir `Local.prorate` per jalur

`[terverifikasi]` dari tabel di atas, dengan penambatan kode transisi (`2` = jalankan lalu lanjut,
`3` = lewati langkah — **`[dugaan kuat]`**, `..\..\10-audit\02-peta-kode-transisi.md` §2.2):

| Jalur | Langkah yang berjalan | **Nilai akhir `Local.prorate`** |
| --- | --- | --- |
| **NB / RNW** (non-EDM) | 1 → 4 | **`100`** |
| **EDM biasa** | 1 → 2 | **`(ProrateEDMEnd + ProrateStartEDM) × 100`** |
| **EDM extend period** | 1 → 2 → 3 | **`100`** |

⛔ **Pada jalur non-EDM, `Local.prorate` BUKAN prorata.** Nilainya selalu `100` dan berfungsi sebagai
**konstanta skala** yang menyatu dengan pembagi — **bukan** pembagian waktu. Kosakatanya dikunci di
`..\..\steering\GLOSARIUM.md` bagian *Kosakata yang dihindari*.

📌 `[dugaan]` `<pyStepsDescription>` langkah 4 berbunyi **`IsNB`** di **kedua** varian — sejalan dengan
pembacaan "jalan saat bukan EDM". **Deskripsi adalah label, bukan bukti** (`PANDUAN-KERJA` §3);
dicatat sebagai penguat saja, tidak dijadikan dasar.

##### ⛔ Dua rumus mati di langkah 1 — **tetap diport**

`[terverifikasi]` Ketiga penugasan berada di `pyParamArray` yang sama, ber-atribut `REPEATINGINDEX`
**1**, **2**, **3**; yang ber-indeks **3** adalah `Local.prorate := 1`.

Kedua rumus pertama — selisih periode, lalu pembagian terhadap `EDMDay`/365 — **tidak berpengaruh
pada jalur mana pun**. Keduanya **diport apa adanya** (`CLAUDE.md` §1), ditandai **KODE MATI**, dan
didaftarkan sebagai **kandidat perbaikan (K-046)** bersama keanehan `Local.RIComIN` varian Aneka di
bawah — `K046_Prorate_DuaRumusMati_Langkah1`.

⚠️ **Yang dijawab work owner adalah NILAI `Local.prorate`, bukan aturan umum** bahwa urutan
`REPEATINGINDEX` = urutan eksekusi. Aturan umum itu **tetap `[pertanyaan terbuka]`** di tempat lain
yang memakainya, dan **tidak ikut ditutup** oleh keputusan ini.

##### Pemeriksaan silang K-018 — kedua pembagi benar

`[terverifikasi]` Pada jalur non-EDM, dengan `Local.prorate = 100`:

| Varian | Rumus efektif | Skala `.Rate` | Sesuai K-018? |
| --- | --- | :-: | :-: |
| **FIRE** (`Data-PropertyItem`) | `Share × Rate × 100 / 100000` = `Share × Rate / 1000` | **‰** | ✅ |
| **ANEKA** (`Data-Aneka`) | `Share × Rate × 100 / 10000` = `Share × Rate / 100` | **%** | ✅ |

⛔ Keduanya konsisten dengan **skala rasio per lini bisnis** yang dikunci **K-018** — FIRE per mille,
ANEKA persen. Pembagi yang berbeda **bukan kekeliruan ekspor**; ia konsekuensi langsung **aturan
pembagi komposit**. Menyeragamkannya menghasilkan premi **10× salah** pada salah satu jalur.

#### Penskalaan mata uang sebelum dipakai

`[terverifikasi]` `Local.ShareOffered` diskala ulang **sebelum** masuk rumus premi:

| Gerbang | Penskalaan |
| --- | --- |
| `.Currency=="IDR" && …PctPremiAllObj != ""` | `ShareOffered * PctPremiAllObj / 100` |
| `.Currency=="USD" && …PctPremiAllObjUSD != ""` | `ShareOffered * PctPremiAllObjUSD / 100` |

⛔ Hanya **dua** mata uang yang punya cabang. Mata uang lain melewati keduanya tanpa penskalaan —
**diport apa adanya**, bukan digeneralisasi.

#### Rumus kedua adalah KOREKSI, bukan alternatif

`[terverifikasi]` Langkah `RH_1.pySteps(10)` bergerbang kesetaraan **4 desimal**:
`@Math.divide(Local.TotalPremiCov,1,4) == @Math.divide(.CoverageList(1).FacOutObjectList(1).ObjectPremi,1,4)`
dengan T=3 / F=2 → langkah bersarang `10.1` berjalan **ketika keduanya TIDAK sama**:

```
.PremiumRetro = @Math.divide((@toDecimal(Primary.CoverageList(1).FacOutObjectList(1).ShareOffered)
                              * .Rate * Local.prorate), 100000, 20)
```

⛔ Rumus koreksi ini memakai `ShareOffered` **langsung dari halaman Primary**, sehingga **MELEWATI
penskalaan mata uang** di atas. Itu perbedaan nyata, bukan penyederhanaan.

⚠️ Ada **dua gerbang kesetaraan pada presisi berbeda**: rate pada **10 desimal**
(`RH_1.pySteps(9).pySteps(1)`), premi pada **4 desimal** (`RH_1.pySteps(10)`). Presisinya **literal di
tempatnya** (ADR-F-0005), tidak diseragamkan.

#### ⛔ Keanehan `Local.RIComIN` — **hidup di varian ANEKA saja**

> Sesi sebelumnya **mencabut keanehan ini seluruhnya**. ⛔ **Itu terlalu jauh, dan dibatasi ulang.**
> Pencabutan **benar untuk varian FIRE**, tetapi **salah untuk varian ANEKA**.

`[terverifikasi]` Pengukuran per varian:

| Salinan | Kelas | `<PropertiesName>` (menyetel) | `<PropertiesValue>` (membaca) | Keanehan |
| --- | --- | ---: | ---: | :-: |
| `DDL\CountRateRetroCov.xml` (FIRE, berlaku) | `Data-PropertyItem` | **1** | 2 | ✅ tidak ada |
| `NB FacIn\Activity\CountRateRetroCov.xml` (FIRE, basi) | `Data-PropertyItem` | 0 | 1 | *(artefak salinan basi)* |
| **`DDL\CountRateRetroCov(ANEKA).xml` (ANEKA, berlaku)** | `Data-Aneka` | **0** | **1** | ⛔ **ADA** |

⛔ **Di jalur ANEKA, `.RICommPercentage = @Math.divide(Local.RIComIN, Local.LengCov, 20)` membagi dari
variabel yang TIDAK PERNAH DIISI.** Langkah 8 varian Aneka hanya mengakumulasi `Local.RateIN` dan
`Local.LengCov`; **tidak ada** `Local.RIComIN := Local.RIComIN + .RIComm`.

**Diport apa adanya; kandidat perbaikan (K-046) — terbatas varian ANEKA.**

```powershell
### FIRE(DDL) setel=1 baca=2 ; ANEKA(DDL) setel=0 baca=1
foreach($p in @('D:\migrasi\RNM\DDL\CountRateRetroCov.xml',
                'D:\migrasi\RNM\DDL\CountRateRetroCov(ANEKA).xml')){
  $t=[IO.File]::ReadAllText($p); $s=0;$b=0
  foreach($m in [regex]::Matches($t,'<PropertiesName>([^<]*)</PropertiesName>')){ if(([regex]'RIComIN').IsMatch($m.Groups[1].Value)){$s++} }
  foreach($m in [regex]::Matches($t,'<PropertiesValue>([^<]*)</PropertiesValue>')){ if(([regex]'RIComIN').IsMatch($m.Groups[1].Value)){$b++} }
  "$([IO.Path]::GetFileName($p)) : setel=$s baca=$b" }
```

⚠️ `[pertanyaan terbuka]` Bahkan di varian FIRE, `Local.RIComIN` **tidak pernah di-nol-kan** —
`pySteps(5)` menginisialisasi `Local.RateIN := 0` dan `Local.TotalPremiCov := 0`, tetapi tidak ada
`Local.RIComIN := 0`. Apakah variabel lokal Pega otomatis nol adalah **`[di luar korpus]`**.

#### ⚠️ Jalur properti mata uang berbeda antar varian

`[terverifikasi]` FIRE menguji **`.Currency`**; ANEKA menguji **`.Currency.Name`**.

⛔ **Satu implementasi yang membaca mata uang lewat satu jalur akan SALAH pada salah satu varian** —
gerbangnya tidak pernah benar, sehingga **penskalaan mata uang terlewat diam-diam** dan premi keluar
tanpa diskala. Kegagalannya **senyap**, bukan galat.

#### `[pertanyaan terbuka]` Kemungkinan varian KETIGA — kelas `Data-Cargo`

`[terverifikasi]` Indeks `Embed-Reference-Rule` pada **7 berkas** mencatat `CountRateRetroCov` pada
**tiga** kelas: `ASM-FW-GISFW-Data-PropertyItem`, `ASM-FW-GISFW-Data-Aneka`, dan
**`ASM-FW-GISFW-Data-Cargo`**. Tiap pemanggil produksi memuat **tiga** langkah `Call
CountRateRetroCov`, seluruhnya bergerbang `.CoverageList(1).PremiumRetro==""`.

⛔ **Varian berkelas `Data-Cargo` TIDAK ADA** di ketiga folder korpus maupun di `DDL\`.

⚠️ Work owner menyatakan hanya ada **dua** rule di Pega. Indeks rujukan merekam resolusi **saat berkas
terakhir disimpan**, jadi ia **bukan bukti** rule itu masih ada sekarang. **Jangan disimpulkan salah
satu.** Rincian: `..\..\10-audit\09-pemanggil-countrateretrocov.md`.

**Asal (Pega).** **`DDL\CountRateRetroCov.xml`** (varian FIRE) · **`DDL\CountRateRetroCov(ANEKA).xml`**
(varian ANEKA) — keduanya berlaku (K-060). Pembanding basi:
`NB FacIn\Activity\CountRateRetroCov.xml`.

**Keputusan.** **K-060** · **K-058** · **K-048** §8.1 · K-018 · K-010/K-012 ·
K-027 · **K-046** · ADR-F-0004/0005 · `CLAUDE.md` §1, §4.1, §4.6

**Blocked by:** F05 · `..\01-tracer-money-ratio-premi-pa.md` · `..\03-resolver-cob-skala.md`
⛔ Kedua tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] ⛔ **Diskon sebelum komisi RI** — urutan terkunci dan teruji secara terpisah
- [ ] Pembagi **100.000** diturunkan **resolver K-018**, bukan angka hafalan
- [ ] Enam penugasan `prorate` **ditimpa berurutan**, tidak dipecah jadi rasio bernama berbeda (K-048 §8.1)
- [ ] ⛔ `Local.prorate` keluar langkah 1 bernilai **`1`** (keputusan work owner 21 Sept) — dua rumus pertama **KODE MATI**, **tetap diport**, `K046_Prorate_DuaRumusMati_Langkah1`
- [ ] Nilai akhir prorata per jalur **diuji terpisah**: non-EDM `100` · EDM `(ProrateEDMEnd+ProrateStartEDM)×100` · EDM extend `100`
- [ ] Pada jalur non-EDM, `100` diperlakukan sebagai **konstanta skala**, bukan prorata — tidak dinamai `prorate` di kode baru
- [ ] Gerbang langkah 4 **terbalik** (jalan saat bukan EDM) — arah dibaca dari kode transisi, bukan dari label
- [ ] Penskalaan mata uang hanya untuk **dua** mata uang; sisanya lewat tanpa penskalaan
- [ ] Rumus koreksi berjalan **saat total TIDAK sama**, dan **melewati** penskalaan mata uang
- [ ] Dua presisi pembanding (**10** dan **4** desimal) **literal di tempatnya**, tidak diseragamkan (ADR-F-0005)
- [ ] ⛔ **DUA varian diport, bukan satu** — dipilih menurut **COB** (K-060 amandemen)
- [ ] Pembagi **100.000** (FIRE) dan **10.000** (ANEKA) **tidak diseragamkan**
- [ ] Mata uang dibaca lewat **jalur yang benar per varian** — `.Currency` (FIRE) vs `.Currency.Name` (ANEKA)
- [ ] Varian ANEKA **tanpa** rumus koreksi dan **tanpa** `Local.TotalPremiCov`
- [ ] `Local.RIComIN` diakumulasi **hanya** di varian FIRE; **K-046** `K046_RIComIN_TakDiisi_HanyaAneka`
- [ ] Gerbang langkah 3 mengikuti **salinan DDL**: `IsEdmExtendPeriod` **saja**
- [ ] **K-046** `K046_Prorate_EnamPenugasanSalingMenimpa`
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`; `Money + Ratio` **gagal saat kompilasi**
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6), **beserta kelasnya** — nama saja tidak cukup

**`[pertanyaan terbuka]` yang tersisa:**
1. ✅ **DIJAWAB WORK OWNER 21 September 2026 — ditutup untuk tiket ini.** ~~Apakah urutan
   `REPEATINGINDEX` dalam satu `Property-Set` = urutan eksekusi; bila ya, dua penugasan `prorate`
   pertama di langkah 1 mati.~~ Work owner menetapkan `Local.prorate` **bernilai `1`** keluar dari
   langkah 1, sehingga kedua rumus itu **memang mati** — lihat bagian *Dua rumus mati di langkah 1*.
   ⚠️ **Yang ditutup adalah NILAI `Local.prorate` pada rule ini**, **bukan** aturan umum bahwa
   `REPEATINGINDEX` = urutan eksekusi. Aturan umum itu tetap **`[di luar korpus]`** dan tetap terbuka
   di tiket/dokumen lain yang menyandarinya.
2. Mengapa gerbang `IsEDM` hilang di langkah 3 salinan DDL varian FIRE — disengaja atau kekeliruan
   ekspor. Korpus tidak dapat menjawab.
3. Apakah `Local.RIComIN` perlu di-nol-kan eksplisit di awal.
4. **Apakah ada varian ketiga berkelas `Data-Cargo`** — lihat bagian di atas.
5. **Apa yang memilih varian saat runtime.** Korpus tidak memuat gerbang pemilih; pemilihan terjadi
   lewat **kelas halaman** tempat rule dipanggil. Pemetaan COB → kelas belum terbaca dari korpus.

## Fac Out - F07 - Validasi Fac Out sebelum penawaran dikirim

**What to build:** Sebelum penawaran retrosesi berjalan, sistem memeriksa bahwa daftar objek Fac Out
konsisten dengan objek Fac In induknya dan bahwa reasuradur tujuannya sudah terisi — dan menolak
dengan pesan bila tidak.

`[terverifikasi]` Dua pemeriksa terpisah, keduanya `Rule-Obj-Activity`:

| Pemeriksa | Gerbang yang terbaca | Akibat |
| --- | --- | --- |
| `CheckDataFacOut_Act` | `Local.ListObjectFacout != Local.ListObjectFacIn` · `Local.FlagError==1` · `pyWorkPage.OfferFacIn.FlagSaveFO=="1"` | `Page-Set-Messages` — pesan galat |
| `CheckProtectFacout_Act` | `.ReinsurerName==""` · `Local.reas<1` · `ProtectReasFacOut.CARI2=="1"` · `…FacRetroDetails.BackUpStatus==3` | `Property-Set-Messages` + `Obj-Save`; menyetel `QuotationData.FacOutStatus` |

⚠️ **Perbandingan daftar objek dilakukan sebagai perbandingan dua nilai terangkai**, bukan sebagai
himpunan. `[terverifikasi]` Keduanya variabel lokal yang dibandingkan langsung dengan `!=`. Bentuk
perakitannya **diport apa adanya** — mengubahnya menjadi perbandingan himpunan akan mengubah kapan
galat muncul (urutan dan duplikat mulai berpengaruh).

`[terverifikasi]` **Ambang jumlah lokasi.** `Activity\OfferFacOut_PreAct` memuat gerbang
`@Utilities.SizeOfPropertyList(.OfferFacIn.LocationList)>50`. Ambang **50** diport apa adanya; artinya
**`[pertanyaan terbuka]`** — korpus tidak menjelaskan mengapa 50.

⚠️ `.FlagSaveFO=="1"` — angka dibandingkan **sebagai string**, sejalan dengan pola yang sudah dicatat
di F05. Kandidat K-046.

**Asal (Pega).** `Activity\CheckDataFacOut_Act` · `CheckProtectFacout_Act` · `SetErrorMessageFacOut` ·
gerbang ambang lokasi di `OfferFacOut_PreAct`

**Keputusan.** **K-046** · K-053 · `CLAUDE.md` §1, §4.6

**Blocked by:** F05

**Status:** blocked

- [ ] Ketidaksesuaian daftar objek Fac Out vs Fac In **menghasilkan pesan galat**, bukan penolakan diam
- [ ] Reasuradur kosong menghalangi lanjut, dengan pesan yang dapat dibaca pengguna
- [ ] Ambang **50 lokasi** diport apa adanya; komentar menandainya `[pertanyaan terbuka]`
- [ ] Perbandingan daftar objek **mempertahankan bentuk aslinya**, tidak diubah menjadi perbandingan himpunan
- [ ] **K-046** `K046_FlagSaveFO_AngkaSebagaiString`
- [ ] Pesan galat **tidak memuat** nama orang, nomor polis, atau data pelanggan (K-025)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F08 - Tangga persetujuan retrosesi — empat tingkat

**What to build:** Penawaran retrosesi berjalan melalui tangga peran tersendiri sampai disetujui atau
ditolak, lengkap dengan percabangan grup dan jejak persetujuan.

⛔ **Ini BUKAN Seam 2.** `[terverifikasi]` `Flow\OfferFacRetro` memuat **nol shape Approval** di
keempat salinannya — ia **tidak** memakai mesin akseptasi limit-berjenjang Fac In. Tangga retro adalah
mesin tersendiri (**K-055**).

#### Tangga empat tingkat berlaku **MENYELURUH untuk ketiga siklus**

✅ **K-058** menetapkan: rule yang sudah **diekspor ulang ke `DDL\` MENANG** atas salinan korpus.
`DDL\OfferFacRetro.xml` memuat tangga **empat tingkat**, maka **empat tingkat berlaku untuk NB, RNW
dan EDM** — tanpa pembatas lingkup.

`[terverifikasi]` Empat salinan, **keempatnya ber-`pyRuleSetVersion` IDENTIK `01-01-95`**:

| Salinan | Ukuran | Assignment | Gateway | Workbasket | Status |
| --- | ---: | ---: | ---: | --- | --- |
| `DDL\OfferFacRetro.xml` (ekspor ulang, commit 2026-08-31) | 171.308 B | 5 | 6 | Admin/Head/GroupLeader/TechnicalDirector | ✅ **BERLAKU** |
| `Endorsment Fac In\Flow\OfferFacRetro.xml` | 170.966 B | 5 | 6 | idem — **4 tingkat** | sejalan |
| `NB FacIn\Flow\OfferFacRetro.xml` | 134.214 B | 3 | 4 | Admin + Head — **2 tingkat** | ⛔ **BASI** |
| `RNW Fac In\Flow\OfferFacRetro.xml` | 134.214 B | 3 | 4 | **identik byte** dengan NB | ⛔ **BASI** |

⚠️ **Salinan NB/RNW yang 2 tingkat dinyatakan BASI — dicatat, bukan dihapus** (`PANDUAN-KERJA` §7).
Ia tetap berguna sebagai pembanding untuk mengukur apa yang berubah.

⛔ **Nomor versi tidak menolong di sini** — keempatnya `01-01-95`, isinya berbeda. Ini salah satu
dasar pertanyaan terbuka K-058 tentang metode deteksi drift.

#### ✅ Layar keempat peran SAMA

**K-058:** perbedaan antar salinan **hanya pada ALURNYA** (2 tingkat vs 4 tingkat). **Tampilan menu
dan isinya SAMA untuk keempat peran** — `ReasFacOutGroupLeader` dan `ReasFacOutTechnicalDirector`
memakai **layar yang sama** dengan `ReasFacOutAdmin` dan `ReasFacOutHead`.

📌 Konsekuensi implementasi: **satu set layar**, empat tahap alur. Jangan membangun layar terpisah
per peran.

#### Isi tangga

`[terverifikasi]` Dari `<pyMOName>` pada `DDL\OfferFacRetro.xml`: empat peran
`ReasFacOutAdmin` · `ReasFacOutHead` · `ReasFacOutGroupLeader` · `ReasFacOutTechnicalDirector`,
masing-masing **satu kali**; `Is it group?` / `IsGroup`; `Accept?` **4×**; `reject` **4×**;
`confirm` **4×**; `IsRISlip` **2×**; `OfferFacOut` **4×**; routing `pyRouteTo` = `Custom`.

`[terverifikasi]` Siklus hidup layarnya: `Activity\OfferFacOut_PreAct` bergerbang
`.pyWorkBasketName=="ReasFacOutAdmin" || .pyWorkBasketName=="ReasFacOutHead"`, memanggil
`GetOldDataRetro_ACT`, `SumCurrencyListAllRetro_Act`, `GetHistoryAkseptasiPega_Act`,
`ProtekReinsurerList_Act`, `CheckDataFacOut_Act`. `OfferFacOut_PostAct` menulis jejak
`ViewSuggest(<APPEND>).Approval` / `.CommentSuggest` / `.DateSuggest` dan membaca
`EmailTypeRetro` / `EmailTypeRetroSlip` bernilai `1` / `2` / `7`.

⚠️ **Satu keputusan manusia = satu transisi** (`GLOSARIUM.md`) — bukan satu proses yang menghitung
seluruh rantai sekaligus.

**Asal (Pega).** `DDL\OfferFacRetro.xml` · `Flow\OfferFacRetro` (tiga salinan korpus) ·
`Activity\OfferFacOut_PreAct` · `OfferFacOut_PostAct`

**Keputusan.** **K-055** · **K-058** (salinan `DDL\` menang; layar keempat peran sama) · K-053 ·
`CLAUDE.md` §4.6

**Blocked by:** F05

**Status:** blocked

- [ ] Empat tingkat berurutan, **satu keputusan = satu transisi**
- [ ] Percabangan grup hidup sebagai gerbang tersendiri
- [ ] Tiap tingkat punya jalur **terima** dan **tolak**
- [ ] ⛔ **Tidak memakai mesin akseptasi Fac In (Seam 2)** — nol shape Approval
- [ ] Jejak persetujuan tercatat per langkah (siapa, kapan, komentar) — **identitas disimpan sebagai rujukan, nilainya tidak disalin ke artefak** (K-025)
- [ ] Enumerasi tipe email `1`/`2`/`7` **diport apa adanya**; artinya `[pertanyaan terbuka]` → `panic` bila nilai lain muncul
- [ ] Empat tingkat berlaku **menyeluruh** untuk NB, RNW dan EDM (K-058) — tanpa cabang per siklus
- [ ] **Satu set layar** dipakai keempat peran; GroupLeader dan TechnicalDirector memakai layar Admin/Head
- [ ] Komentar menyebut salinan `DDL\OfferFacRetro.xml` sebagai sumber, dan menandai salinan NB/RNW 2 tingkat sebagai **basi**
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F09 - Dua generator nomor retrosesi — tidak disatukan

**What to build:** Penawaran retrosesi mendapat nomor slip yang unik dan berurut, memakai **dua skema
penomoran berbeda** yang hidup berdampingan di sistem lama.

⛔ **DUA generator, bukan satu.** `[terverifikasi]` Skemanya berbeda, urutannya berbeda, dan
**sequence-nya berbeda**:

| | Generator 1 | Generator 2 |
| --- | --- | --- |
| Berkas | `DDL\GENERATE_FACRETRO_NO.txt` (fungsi Oracle) | `RDBList\GenerateOurRefFacOut_SQL` (`Rule-Connect-SQL`, tag `pyBrowseSQL`) |
| Bentuk | `RNM-` + FacCode + BusinessCode + `.` + bulan + `.` + tahun + `.` + `LPAD(seq,5,'0')` | `'Y'` + `{InputData.CARI21}` + `{…BusinessOldId}` + `'RNM'` + `to_char(sysdate,'yy')` + seq |
| Sequence | `FACRETRO_SEQ` | **`T_FAC_OFFER_SEQ`** |
| Padding | `LPAD(…,5,'0')` | **tanpa padding** |
| Posisi `'Y'` | di tengah (sebagai FacCode) | **di depan** |

`[terverifikasi]` Pemetaan FacCode: `'Y'` bila `FacType = 'FAKULTATIF RETROSESI'`; `'F'` bila
`'FAKULTATIF INWARD'`. Jadi retrosesi berawalan **`RNM-Y…`**, berbeda dari Fac In `RNM-F…`.

⛔ **Tidak ada cabang `ELSE`.** `[terverifikasi]` Bila `FacType` bukan salah satu dari kedua nilai itu,
`FacCode` tetap **NULL**. Perilaku itu **diport apa adanya** — jangan ditambahi cabang default.

⚠️ **`FacRetroSlipNumber` dideklarasikan `VARCHAR2(21)`.** `[terverifikasi]` Panjang keluarannya
bergantung panjang `BusinessCode`; bila melebihi 21 karakter, fungsi Oracle gagal saat runtime.
**Diport apa adanya**, dengan batas panjangnya diuji. `[pertanyaan terbuka]` apakah `BusinessCode`
dijamin pendek — `GLOSARIUM.md` mencatat 98 kode dengan arti **belum terverifikasi**.

📌 `[pertanyaan terbuka]` Korpus **tidak menjelaskan kapan** masing-masing generator dipakai.
Keduanya diport; pemilihannya mengikuti pemanggil, bukan ditebak.

**Asal (Pega).** `DDL\GENERATE_FACRETRO_NO.txt` · `RDBList\GenerateOurRefFacOut_SQL`

**Keputusan.** K-055 · K-056 · `CLAUDE.md` §1, §4.3, §4.6

**Blocked by:** F08

**Status:** blocked

- [ ] **Kedua generator** hidup terpisah; tidak diseragamkan menjadi satu
- [ ] Bentuk nomor retro `RNM-Y…` berbeda dari Fac In `RNM-F…`
- [ ] ⛔ **Tanpa cabang `ELSE`** — `FacType` tak dikenal menghasilkan kode kosong, **bukan** nilai default
- [ ] Padding 5 digit hanya pada generator 1; generator 2 **tanpa padding**
- [ ] Nomor urut diambil dari sequence basis data, **bukan** dihitung aplikasi
- [ ] Batas panjang `VARCHAR2(21)` diuji sebagai **batas nyata**, bukan diabaikan
- [ ] **K-046** `K046_GenerateFacRetroNo_TanpaElse`
- [ ] Nomor polis produksi **tidak pernah** muncul di test; memakai nilai sintetis (K-025)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F10 - Produksi Fac Out ke tabel `FACOUTPRODUCTION`

**What to build:** Penawaran retrosesi yang sudah disetujui tersimpan sebagai baris produksi Fac Out,
tertaut ke objek Fac In induknya, dan **aman diulang** tanpa menghasilkan baris ganda.

✅ **Tabel produksinya sudah ada dan sudah flat** — 65 kolom, multi-lini, dengan pasangan
`_MENJADI`/`_SELISIH`. **Tidak perlu tabel flat baru** (K-056). Ini membedakan F10 dari `edm\E21`, yang
masih menunggu rancangan tabel flat.

⛔ **Insert = Connect-SQL langsung, BUKAN stored procedure.** `[terverifikasi]`
`DDL\InsertTreatyProd_Sql.xml`, `<pxObjClass>` = `Rule-Connect-SQL`, tag `<pyBrowseSQL>`:
`BEGIN INSERT INTO FACOUTPRODUCTION (…65 kolom…) VALUES (…); COMMIT; END;`. Catatan lama "via SP
`INSERTFACOUTPRODUCTION`" sudah **dibatalkan** — itu nama request/connector, bukan SP basis data.

#### Gerbang idempotensi

`[terverifikasi]` `Activity\InsertFacoutProduction` memanggil `RDBList\CekFacoutProd_Sql`:

```sql
select DISTINCT IDPEGA from FACoutPRODUCTION
 where RISLIPNO = {DataIN.CARI7} and IDPEGA = {pyWorkPage.pzInsKey}
```

Gerbangnya `@SizeOfPropertyList(FacoutList.pxResults)>0` dengan
`pyStepsPreCondParamsWhenTrue` = **6** / `pyStepsPreCondParamsWhenFalse` = **2**.

`[terverifikasi]` Gerbang itu berada di alamat **`RH_1.pySteps(7)`**, dan
`<pyStepsPreCondParamsWhen>`, `…WhenTrue`, `…WhenFalse` berada di **`rowdata` yang sama** — keanggotaan
struktural, bukan kedekatan teks.

`[dugaan kuat]` Penambatan kode transisi dari dua jangkar independen (`SaveEDMToJsonPolicy_Act`
`RH_1.pySteps(10)` T=6, dan `SaveTreatyProduction_Act` `RH_1.pySteps(4)` T=2/F=6): **`2` = jalankan
lalu lanjut · `3` = lewati langkah · `6` = keluar activity**. Maka: **baris sudah ada → keluar, tidak
menyisipkan; belum ada → lanjut dan sisipkan.**

⚠️ Ditandai **`[dugaan kuat]`, bukan `[terverifikasi]`**. Yang sudah terbukti adalah **keanggotaan
struktural** gerbang dan T/F pada langkah yang sama (ditunjukkan lewat `<pyStepPageReference>` untuk
kelima jangkar) serta pola distribusinya. Yang **belum** terbukti adalah **arti** kodenya — korpus
tidak memuat berkas yang memetakan kode→perilaku secara eksplisit. Kenaikan label memerlukan
konfirmasi UI Pega work owner.

Kode **1**, **4** dan **5** **belum tertambat** — `[pertanyaan terbuka]`, butuh UI Pega work owner.
Kode **5** yang paling sering di antaranya (908 kemunculan pada satu pasangan).

📌 **Sensus lengkap 31 pasangan (T,F), bukti keanggotaan kelima jangkar, dan perintah auditnya ada di
`..\..\10-audit\02-peta-kode-transisi.md`.**

#### Pemetaan 65 kolom — disalin dari `09-facout\01` §4.3.1

`[terverifikasi]` **65 kolom seimbang dengan 65 nilai**, parameter dari **lima halaman**:
`DataIN` 35 · `DataIN1` 14 · `DataINCargo` 9 · `DataINMBU` 5 · `DataINPA` 2.

| # | Kolom | Tipe Oracle | Parameter |
| ---: | --- | --- | --- |
| 1 | `IDPEGA` | `VARCHAR2(150)` | `{DataIN.CARI1}` |
| 2 | `POLICYNO` | `VARCHAR2(50)` | `{DataIN.CARI2}` |
| 3 | `GROUPPANEL` | `VARCHAR2(10)` | `{DataIN.CARI3}` |
| 4 | `REINSURER_ID` | `VARCHAR2(15)` | `{DataIN.CARI4}` |
| 5 | `REINSURER_NAME` | `VARCHAR2(150)` | `{DataIN.CARI5}` |
| 6 | `TGL_PRINT` | `DATE` | `To_date({DataIN.CARI6}, 'DD/MM/YYYY HH24:MI:SS')` |
| 7 | `RISTARTPERIOD` | `DATE` | `To_date({DataIN.CARI8}, …)` |
| 8 | `RIENDPERIOD` | `DATE` | `To_date({DataIN.CARI9}, …)` |
| 9 | `START_DATE` | `DATE` | `To_date({DataIN.CARI11}, …)` |
| 10 | `END_DATE` | `DATE` | `To_date({DataIN.CARI12}, …)` |
| 11 | `RISLIPNO` | `VARCHAR2(600)` | `{DataIN.CARI7}` |
| 12 | `NOENDORS` | `VARCHAR2(100)` | `{DataIN.CARI13}` |
| 13 | `PACKINGID` | `VARCHAR2(200)` | `{DataINCargo.CARI1}` |
| 14 | `PACKINGNOTE` | `VARCHAR2(900)` | `{DataINCargo.CARI2}` |
| 15 | `GOODNOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI3}` |
| 16 | `TRADINGNOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI4}` |
| 17 | `SHIPID` | `VARCHAR2(100)` | `{DataINCargo.CARI6}` |
| 18 | `FROMRUTE` | `VARCHAR2(500)` | `{DataINCargo.CARI7}` |
| 19 | `TORUTE` | `VARCHAR2(500)` | `{DataINCargo.CARI8}` |
| 20 | `SAILDATE` | `DATE` | `To_date({DataINCargo.CARI9}, …)` |
| 21 | `CONVEYANCENOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI5}` |
| 22 | `OBJECTNO` | `VARCHAR2(50)` | `{DataIN1.CARI1}` |
| 23 | `ZIPCODE` | `VARCHAR2(50)` | `{DataIN1.CARI2}` |
| 24 | `PROVINCE` | `VARCHAR2(100)` | `{DataIN1.CARI3}` |
| 25 | `CITY` | `VARCHAR2(100)` | `{DataIN1.CARI4}` |
| 26 | `DISTRICT` | `VARCHAR2(100)` | `{DataIN1.CARI5}` |
| 27 | `RW` | `VARCHAR2(1000)` | `{DataIN1.CARI6}` |
| 28 | `ADDRESS` | `VARCHAR2(4000)` | `{DataIN1.CARI7}` |
| 29 | `BUILDINGNO` | `VARCHAR2(30)` | `{DataIN1.CARI8}` |
| 30 | `ROADNAME` | `VARCHAR2(4000)` | `{DataIN1.CARI9}` |
| 31 | `OBJECTNAME` | `VARCHAR2(4000 CHAR)` | `{DataIN1.CARI10}` |
| 32 | `OCCUPATION` | `VARCHAR2(4000)` | `{DataIN1.CARI11}` |
| 33 | `OBJECTITEM` | `VARCHAR2(4000)` | `{DataIN1.CARI12}` |
| 34 | `OBJECTITEMID` | `VARCHAR2(100)` | `{DataIN1.CARI13}` |
| 35 | `BRANDNAME` | `VARCHAR2(500 CHAR)` | `{DataINMBU.CARI2}` |
| 36 | `LICENSEPLATE` | `VARCHAR2(50)` | `{DataINMBU.CARI3}` |
| 37 | `TYPENAME` | `VARCHAR2(50)` | `{DataINMBU.CARI4}` |
| 38 | `MODELNAME` | `VARCHAR2(500)` | `{DataINMBU.CARI5}` |
| 39 | `ENGINENUMBER` | `VARCHAR2(100)` | `{DataINMBU.CARI6}` |
| 40 | `CLASSPA` | `VARCHAR2(100)` | `{DataINPA.CARI1}` |
| 41 | `DOB` | `VARCHAR2(100)` | `{DataINPA.CARI2}` |
| 42 | `CURRENCY` | `VARCHAR2(10)` | `{DataIN.CARI23}` |
| 43 | `CURRENCYID` | `VARCHAR2(10)` | `{DataIN.CARI29}` |
| 44 | `TSIRNM` | `NUMBER(20,4)` | `{DataIN.CARI30}` |
| 45 | `TSISPREADED` | `NUMBER(20,4)` | `{DataIN.CARI31}` |
| 46 | `PCTOFFERED` | `NUMBER(20,4)` | `{DataIN.CARI32}` |
| 47 | `SHAREOFFERED` | `NUMBER(20,4)` | `{DataIN.CARI24}` |
| 48 | `OBJECTPREMI` | `NUMBER(20,4)` | `{DataIN.CARI25}` |
| 49 | `RICOMM` | `NUMBER(20,4)` | `{DataIN.CARI26}` |
| 50 | `COMMISION` | `NUMBER(20,4)` | `{DataIN.CARI27}` |
| 51 | `RATE` | `NUMBER(25,20)` | `{DataIN.CARI28}` |
| 52 | `SHAREOFFERED_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI33}` |
| 53 | `OBJECTPREMI_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI34}` |
| 54 | `COMMISION_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI35}` |
| 55 | `TYPEFACULTATIVE` | `VARCHAR2(50)` | `{DataIN.CARI39}` |
| 56 | `PRORATE` | `NUMBER(25,20)` | `{DataIN.CARI40}` |
| 57 | `RATE_COVERAGE` | `NUMBER(20,4)` | `{DataIN.CARI36}` |
| 58 | `PREMI_COVERAGE_MENJADI` | `NUMBER(20,8)` | `{DataIN.CARI37}` |
| 59 | `PREMI_COVERAGE_SELISIH` | `NUMBER(20,8)` | `{DataIN.CARI38}` |
| 60 | `COVERAGE_NAME` | `VARCHAR2(1000)` | `{DataIN.CARI42}` |
| 61 | `COVERAGE_ID` | `VARCHAR2(100)` | `{DataIN.CARI41}` |
| 62 | `COMMISION_COVERAGE_PCT` | `NUMBER(20,8)` | `{DataIN.CARI43}` |
| 63 | `COMMISION_COVERAGE_MENJADI` | `NUMBER(20,8)` | `{DataIN.CARI44}` |
| 64 | `COMMISION_COVERAGE_SELISIH` | `NUMBER(20,8)` | `{DataIN.CARI45}` |
| 65 | `OBJECTNO_FACIN` | `VARCHAR2(10)` | `{DataIN1.CARI14}` |

⛔ **Pemetaannya TIDAK berurutan** — lihat kolom 6–11, 17–21, dan 42–57. **Diport apa adanya; jangan
"dirapikan".** Satu geseran kolom = korupsi data produksi yang senyap.

`[terverifikasi]` **Sebelas parameter tidak terpakai:** `DataIN.CARI10` dan `CARI14`…`CARI22`
(sepuluh), plus `DataINMBU.CARI1`.

⛔ **Interpolasi string diganti parameter terikat.** `{Halaman.CARIn}` di sistem lama adalah
**substitusi teks** ke dalam SQL. Di sistem baru ia **wajib** menjadi parameter terikat. Itu perubahan
mekanisme penyampaian nilai, **bukan** perubahan perilaku — nilai yang tertulis ke kolom tetap sama,
sehingga tidak melanggar `CLAUDE.md` §1.

⚠️ **Kolom ber-PII:** `DOB`, `LICENSEPLATE`, `ENGINENUMBER`, `ADDRESS`, `OBJECTNAME`,
`REINSURER_NAME`. Nilainya **tidak pernah** disalin ke test, fixture, log, atau dokumen (K-025).

**Asal (Pega).** `DDL\InsertTreatyProd_Sql.xml` · `DDL\FACOUTPRODUCTION.txt` ·
`Activity\InsertFacoutProd` (pengirim) · `InsertFacoutProduction` · `RDBList\CekFacoutProd_Sql` ·
`Activity\UpdateStsKonversiFacOut_Act`

**Keputusan.** **K-056** · K-053 · K-027 · K-010/K-012 · `CLAUDE.md` §4.1, §4.3, §4.4, §4.6

**Blocked by:** F06 · F08 · F09

**Status:** blocked

- [ ] Baris produksi tertulis dengan **65 kolom** sesuai tabel di atas, **urutan pemetaan apa adanya**
- [ ] ⛔ Pemetaan diambil dari **tabel**, bukan dari urutan kolom — tidak ada penyesuaian "supaya rapi"
- [ ] Sebelas parameter tak terpakai **tetap tak terpakai**
- [ ] **Parameter terikat**, bukan interpolasi string; alasannya dicatat dalam komentar
- [ ] Gerbang idempotensi: baris sudah ada → **tidak menyisipkan**; rantai **aman diulang**
- [ ] Penaut `OBJECTNO_FACIN` terisi dari nilai yang lahir di **F02**, bukan dihitung ulang
- [ ] Uang bertipe `Money` dengan mata uang; ditulis ke `NUMBER` Oracle — **tidak pernah `float`**
- [ ] Skema Oracle **tidak berubah** (§4.3)
- [ ] ⛔ Nilai `DOB`/`LICENSEPLATE`/`ENGINENUMBER`/`ADDRESS` **tidak pernah** masuk test atau log
- [ ] `[pertanyaan terbuka]` kode transisi **1**, **4**, **5** belum tertambat — bila tercapai, **`panic`**
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F11 - Jalur Fac Out pada siklus endorsement

**What to build:** Endorsement yang mengubah porsi retrosesi menghasilkan **selisih** Fac Out yang
tertulis berdampingan dengan nilai sesudahnya, memakai jalur produksi Fac Out yang sudah hidup (F10).

#### Populasi kerja — **55, bukan 39**

⛔ **Populasi berbasis nama TIDAK cukup.** `[terverifikasi]`:

| Ukuran | NB | RNW | EDM |
| --- | ---: | ---: | ---: |
| Activity **bernama** `FacOut`/`FacRetro` (Ordinal, peka huruf) | 41 | 40 | 39 |
| Hadir di ketiga folder | — | **39** | — |
| Activity membawa penanda `10015` **tanpa** nama Fac Out | 16 | 15 | 16 |
| **Permukaan fungsional** | ≈57 | **≈55** | **≈55** |

Yang tanpa nama Fac Out termasuk `CopyToAllSpreading_ACT`, `CopyToAllLocSpreading_ACT`,
`cekSpreadingFactIn`, `CountPremiumNet`, `SaveJsonPolicyFacIn_Act`, `SetReinsurerEndorsement_Act`.

✅ **K-059:** keduanya **diperlakukan sama**, karena **`10015` adalah kode Fac Out**. Pembawa kode itu
**bukan** kelompok kelas dua yang "perlu ditinjau" — ia **bagian penuh** dari populasi kerja.

⛔ Angka lama **18/18 dan 10/18 DIBATALKAN** — populasinya tidak pernah terdefinisi.

📌 **Sumber angka:** `..\..\10-audit\03-struktur-facretro-dan-populasi.md` §3 (populasi + banding 23
tag, lengkap dengan uji silang wajib) dan `..\..\10-audit\07-hitung-berkas-dan-drift.md` §1 (jumlah
berkas per subfolder) serta §3 (drift). Keduanya memuat perintah audit yang dapat dijalankan ulang.

#### Sembilan activity yang berbeda di Endorsement

`[terverifikasi]` Dari **39** yang hadir di ketiga folder, kontrak 23 tag (K-042/K-043, termasuk
amandemen "isi blok `pzIndexes` diabaikan"): **NB vs RNW 39/39 identik**; **NB vs EDM 30/39 identik,
9 berbeda**:

`CopyAllObjFacOutFireAneka_ACT` · `CopyFacRetroAnekaGolf_ACT` · `InsertFacoutProd` ·
`InsertFacoutProduction` · `InsertFacoutProductionEDM` · `InsertFacoutProductionEDMCurr` ·
`OfferFacOut_PreAct` · `SetDataFacOutAnekaGolf_Act` · `SetDataFacOutFire_Act`

📌 **NB dan RNW berbagi satu implementasi** — nol perbedaan. Hanya Endorsement yang butuh cabang.

#### Isi khusus endorsement

`[terverifikasi]` `Activity\InsertFacoutProd` adalah pengirim: `Call InsertFacoutProduction` dan
`Call InsertFacoutProductionEDM` (bergerbang `IsEDM`), lalu `Obj-Refresh-And-Lock`, lalu
`Call UpdatePolisAddendum_Act`. Gerbang prefiks `@contains(pyWorkPage.pyWorkIDPrefix,"EDM-")` ada di
`InsertFacoutProduction`.

`[terverifikasi]` Kolom selisih di `FACOUTPRODUCTION`: `SHAREOFFERED_SELISIH`, `OBJECTPREMI_SELISIH`,
`COMMISION_SELISIH`, `PREMI_COVERAGE_MENJADI`/`_SELISIH`, `COMMISION_COVERAGE_MENJADI`/`_SELISIH`,
`RATE_COVERAGE`.

⚠️ `[terverifikasi]` **Ketidaksetangkupan pasangan `_MENJADI` vs `_SELISIH`** — tidak setiap kolom
selisih punya pasangan "menjadi". Sudah tercatat sebagai butir K-046 (nomor 7 pada daftar 12 kode
usang). **Diport apa adanya.**

⚠️ **K-046** `K046_GuardFacOut_HanyaType7` — gerbang idempotensi fac out hanya dipakai keluar pada
`Type=="7"`. Untuk `Type` lain, penulisan ganda **tidak dicegah**. `[terverifikasi]`
`InsertFacoutProductionEDM`. Diport apa adanya; risikonya tercatat sebagai E-Q32.

**Asal (Pega).** `Activity\InsertFacoutProductionEDM` · `InsertFacoutProductionEDMCurr` ·
`CopyFacRetroEDMLoc_ACT` · `ProtectionEDMFacout_Act` · `DataTransform\CountPremiEDMFacOut_DT`
(`Rule-Obj-Model`) · sembilan activity berbeda di atas

**Keputusan.** **K-059** (populasi ≈55 — nama Fac Out **dan** pembawa `10015` diperlakukan sama,
karena `10015` **adalah** kode Fac Out) · K-056 · **K-046** ·
K-042/K-043 (kontrak 23 tag) · `CLAUDE.md` §4.3, §4.6

**Blocked by:** F10 · `..\edm\E21-jalur-produksi-edm.md`
⛔ `E21` sendiri **blocked** menunggu tabel flat (P-10) — lihat indeks EDM.

**Status:** blocked

- [ ] Populasi kerja **≈55 per folder**, bukan 39 — pembawa `10015` **diperlakukan sama** dengan yang bernama Fac Out (K-059)
- [ ] **Sembilan** activity yang berbeda di Endorsement punya cabang tersendiri; **tiga puluh** sisanya dipakai ulang
- [ ] ⛔ **NB dan RNW berbagi satu implementasi** — tidak ada cabang siklus di antara keduanya
- [ ] Kolom selisih terisi berpasangan dengan nilai sesudahnya, **sesuai ketaksetangkupan aslinya**
- [ ] **K-046** `K046_GuardFacOut_HanyaType7` — gerbang idempotensi hanya aktif pada satu jenis penyesuaian
- [ ] **K-046** `K046_MenjadiSelisih_TidakSetangkup`
- [ ] Nilai berkoma desimal (K-027); uang bertipe `Money`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F12 - Layar dan daftar Fac Out

**What to build:** Pengguna dapat membuka menu Fac Out, melihat daftar retrosesi sebuah polis, dan
menelusuri objek serta coverage yang diretrokan — sebagai menu **terpisah** dari Fac In namun
**tertaut** ke kasus induknya (K-053).

`[terverifikasi]` Inventaris berkas Fac Out/FacRetro per tipe rule, ketiga folder digabung, dihitung
**per subfolder** (listing rekursif terpotong di batas entri):

| Tipe | Jumlah |
| --- | ---: |
| Section | 131 |
| Activity | 120 |
| FlowAction | 57 |
| DataTransform | 12 |
| RDBList | 10 |
| When | 9 |
| Harness | 9 |
| DataPage | 6 |
| Flow | 5 |
| **TOTAL** | **359** |

Nama unik: **46 Section · 19 FlowAction · 3 Harness · 4 RDBList · 2 Flow · 2 DataPage**.

`[terverifikasi]` **Daftar Fac Out** dibaca dari `RDBList\GetFacoutList_SQL` (`Rule-Connect-SQL`, tag
`pyBrowseSQL`):

```sql
select IDPEGA as HASIL1, policyno HASIL2 from facoutproduction where policyno = {InputData.CARI17}
```

Jadi daftarnya bersumber dari **tabel produksi**, berkunci nomor polis — bukan dari staging.

#### ⛔ `GetOPFacOut_Sql` — TIDAK diport, **DIHAPUS** (K-062)

✅ **K-062 sudah diputuskan.** Klep ini **tidak diport dan dihapus, tidak dipakai lagi**. **Kolom
pelaksana TETAP ADA di layar Fac Out, tetapi ISINYA KOSONG.**

⛔ **Ini bukan blocker eksternal lagi.** Tidak ada yang perlu ditunggu.

##### ⚠️ Pengecualian yang disengaja terhadap K-006

> **K-062 adalah pengecualian PERTAMA yang disengaja** terhadap pola "ditangguhkan, bukan dihapus"
> (**K-006**, `PANDUAN-KERJA` §5). Pola baku menuntut `panic` bila cabang tercapai; **di sini tidak**.
>
> **Sesi berikutnya JANGAN membacanya sebagai inkonsistensi lalu membalikkannya.** Membalikkannya
> memerlukan keputusan baru.
>
> Alasannya dapat dipertanggungjawabkan: yang hilang adalah **tampilan identitas pelaksana**, bukan
> angka uang atau gerbang alur — mengosongkan kolom **tidak** menghasilkan selisih rekonsiliasi
> (ADR-F-0001).

##### Sisi pemanggil wajib menghasilkan KOSONG, bukan galat

`[terverifikasi]` Pemindaian **peka huruf** atas **seluruh 6.071 berkas** (14 subfolder × 3 folder):
**10 berkas** memuat `GetOPFacOut`, **3** di antaranya berkas definisinya sendiri →
**7 perujuk sejati**.

| Perujuk | NB | RNW | EDM | Tag pembawa |
| --- | :-: | :-: | :-: | --- |
| `Activity\OfferFacOut_PreAct` | ✓ | ✓ | ✓ | `<RequestType>` = `GetOPFacOut_Sql` · `<pyRuleName>` = `ASM-FW-GISFW-Work ASM GetOPFacOut_Sql` |
| `Activity\PrintRISlipPre_act` | ✓ | ✓ | ✓ | idem |
| `Activity\SumCurrencyListAllRetro_Act` | — | — | ✓ | idem |
| `RDBList\GetOPFacOut_Sql` *(definisi)* | ✓ | ✓ | ✓ | `<pyRequestType>` · `<pyLabel>` · `<pxTabLabel>` · `<pyRuleName>` |

⛔ **Ketiga pemanggil itu wajib menghasilkan nilai kosong, BUKAN galat.** `PrintRISlipPre_act`
menyentuh jalur cetak RI Slip (**F13**); `SumCurrencyListAllRetro_Act` hanya ada di Endorsement.

`[terverifikasi]` Yang dibacanya: identitas pelaksana dari riwayat penugasan kasus, disaring pada
langkah bernama `OfferFacOut`/`OfferRetro`, diurutkan menurut waktu commit, dari
**`datapega.pc_history_asm_fw_gisfw_work`** — tabel **internal Pega**. `CLAUDE.md` §4.3:
`DATAPEGA.PC_*` **hilang bersama Pega**. ⚠️ Nilainya **identitas orang** — mekanismenya saja yang
dicatat (K-025).

⚠️ **Nama folder bukan tipe rule.** `[terverifikasi]` `CekFacoutProd_Sql`, `GetOPFacOut_Sql` dan
`GenerateOurRefFacOut_SQL` berada di folder `RDBList\` tetapi ber-`<pxObjClass>` = `Rule-Connect-SQL`
— seluruh **605** berkas di `RDBList\` berkelas itu. Begitu pula seluruh **443** berkas di
`DataTransform\` berkelas `Rule-Obj-Model`. Tipe dibaca dari tag, tidak pernah dari nama folder.

⚠️ `[terverifikasi]` Section Fac Out hadir **berpasangan** `X` dan `X_IsUW` (mis.
`ShowCoverageFacOut` / `ShowCoverageFacOut_IsUW`). Sejalan keputusan work owner pada siklus Renewal,
pasangan `_IsUW` adalah **dua komponen terpisah**, bukan satu komponen dua mode.

**Asal (Pega).** `Harness\ViewFacretro` · `ViewCoverageFacOut` · `ViewCoverageFacOut_isUW` ·
46 `Section\*FacOut*` · 19 `FlowAction\*FacOut*` · `RDBList\GetFacoutList_SQL` ·
`RDBList\GetOPFacOut_Sql` · `DataPage\D_CoverageFacOut` · `D_FacOutFromVehicle`

**Keputusan.** **K-062** (klep dihapus; kolom pelaksana kosong — ⛔ **pengecualian pertama terhadap
K-006**) · K-053 · `CLAUDE.md` §4.3, §4.6

**Blocked by:** F01 · F05

**Status:** blocked

- [ ] Menu Fac Out **terpisah** dari menu Fac In di navigasi, tetapi tertaut lewat penaut objek induk
- [ ] Daftar Fac Out dibaca dari **tabel produksi** berkunci nomor polis
- [ ] ⛔ `GetOPFacOut_Sql` **dihapus, tidak dipakai** (K-062) — **bukan** `panic`, **bukan** ditangguhkan
- [ ] Kolom pelaksana **tetap ada** di layar Fac Out dengan **isi kosong**
- [ ] Ketiga pemanggil (`OfferFacOut_PreAct`, `PrintRISlipPre_act`, `SumCurrencyListAllRetro_Act`) menghasilkan **nilai kosong, bukan galat**
- [ ] Komentar menyatakan sebabnya **dan** menandai ini **pengecualian sadar terhadap K-006**, agar tidak dibalikkan
- [ ] Pasangan `X` / `X_IsUW` tetap **dua komponen terpisah**
- [ ] Tipe rule tiap berkas dibaca dari `<pxObjClass>`, **tidak pernah** dari nama folder
- [ ] Layar **tidak menampilkan** nilai `DOB`, plat nomor, nomor mesin, atau alamat di test/fixture (K-025)
- [ ] Endpoint dari konfigurasi, tidak pernah literal (§4.4)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

## Fac Out - F13 - Cetak R/I Slip dan kirim surat retrosesi ⛔ BLOCKED EKSTERNAL

**What to build:** Setelah tangga retro selesai, sistem mencetak R/I Slip retrosesi dan mengirimkannya
ke reasuradur tujuan.

⛔ **BLOCKED — menunggu isi `M_LINK_SERVICE`.** `CLAUDE.md` §4.4: daftar endpoint sesungguhnya ada di
tabel Oracle `M_LINK_SERVICE` (kunci `KATEGORI_1` + `KATEGORI_2`), yang **isinya tidak ada di korpus**.
Tanpa itu, tujuan kirim tidak dapat ditentukan, dan endpoint **tidak boleh** ditulis literal.

⚠️ **Tiket ini tidak dihapus.** Lingkupnya sudah diketahui; menghilangkannya akan menyembunyikan
pekerjaan yang pasti ada. Ini **pola yang sama dengan `..\rnw\R07-konversi-produksi-renewal.md`** —
tiket ditulis penuh, hanya eksekusinya menunggu.

`[terverifikasi]` Yang ada di korpus:

| Bagian | Bukti |
| --- | --- |
| Shape alur | `DDL\OfferFacRetro.xml` `<pyMOName>`: `PRINT R/I SLIP`, `Send Email Print RI Slip`, `IsRISlip` (2×), `PrintRISlip_FlowAction` |
| Activity PDF | `GeneratePDFFacOutMemoPlacing` · `GeneratePDFFacOutOfferStatus` · `GeneratePDFFacOutViewOffer` · `DeletePDFFacOutMemoPlacing` |
| Section cetak | `FacOutPrintRISlipSectionInside` · `PrintRISlip` · `Harness\PrintRISlips` |
| Data slip | `.OfferFacIn.FacRetroDetails.DocumentPosition` (F05) · nomor slip dari **F09** |
| Tipe email | `OfferFacOut_PostAct` membaca `EmailTypeRetro` / `EmailTypeRetroSlip` bernilai `1` / `2` / `7` |

⚠️ `[pertanyaan terbuka]` Arti `1` / `2` / `7` pada tipe email **tidak dijelaskan korpus**. Nilainya
diport apa adanya; nilai di luar ketiganya → **`panic`** (`CLAUDE.md` §4.5).

⛔ **Empat activity PDF ada, tetapi isi templatnya tidak.** `[dugaan]` Templat dokumen kemungkinan
berada di luar ekspor rule. Selama belum terbukti ada, bentuk keluaran PDF **tidak boleh ditebak**.

⛔ **Alamat email tujuan tidak pernah disalin** ke tiket, test, fixture, log, maupun dokumen (K-025).
Yang dicatat hanya **jumlah dan mekanismenya**.

**Asal (Pega).** `DDL\OfferFacRetro.xml` (shape cetak & kirim) · `Activity\GeneratePDFFacOutMemoPlacing` ·
`GeneratePDFFacOutOfferStatus` · `GeneratePDFFacOutViewOffer` · `DeletePDFFacOutMemoPlacing` ·
`Section\FacOutPrintRISlipSectionInside` · `Activity\OfferFacOut_PostAct`

**Keputusan.** **K-061** (ditulis penuh, ditandai menunggu, pola R07 — **tidak dihapus**) · K-055 ·
K-062 (pemanggil `PrintRISlipPre_act` menghasilkan kosong) · K-025 · `CLAUDE.md` §4.4, §4.5, §4.6

**Blocked by:** F08 · ⛔ **eksternal — isi `M_LINK_SERVICE`**

**Status:** blocked

- [ ] ⛔ Menunggu isi `M_LINK_SERVICE`; **tidak ada endpoint literal** di kode
- [ ] R/I Slip memuat nomor slip dari **F09** dan posisi dokumen dari **F05**
- [ ] Enumerasi tipe email `1`/`2`/`7` diport apa adanya; nilai lain → **`panic`**
- [ ] ⛔ Bentuk keluaran PDF **tidak ditebak** selama templatnya belum terbukti ada di korpus
- [ ] ⛔ **Alamat email dan nama orang tidak pernah** muncul di artefak mana pun (K-025)
- [ ] Endpoint dan kredensial dari konfigurasi (§4.4)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

✅ **K-061 sudah diputuskan: pola R07.** Tiket ditulis penuh, **ditandai menunggu**, **tidak dihapus**.

⛔ **Blocker eksternalnya, dinyatakan tegas:** isi tabel Oracle **`M_LINK_SERVICE`** (kunci
`KATEGORI_1` + `KATEGORI_2`) **tidak ada di korpus**. Tanpa itu tujuan kirim tidak dapat ditentukan,
dan `CLAUDE.md` §4.4 melarang endpoint literal. **Pemiliknya DBA.** Sampai itu tersedia, F13
**tidak dapat dieksekusi** — hanya ditulis.

⚠️ **K-062 menyentuh tiket ini:** `Activity\PrintRISlipPre_act` adalah salah satu dari tiga pemanggil
`GetOPFacOut_Sql`. Klep itu **dihapus**, sehingga pemanggil di jalur cetak ini **wajib menghasilkan
nilai kosong, bukan galat** — lihat **F12**.

## Fac Out - F14 - Rekonsiliasi eksak Fac Out

**What to build:** Pembanding **nol-selisih** untuk baris produksi Fac Out, memakai **kerangka
rekonsiliasi New Business** — tanpa pembanding baru.

`[terverifikasi]` Sasaran pembandingnya sudah ada dan sudah flat: `DDL\FACOUTPRODUCTION.txt` — **65
kolom**, dengan pasangan `_MENJADI`/`_SELISIH`. Pembanding membaca baris yang ditulis **F10**, bukan
struktur staging.

`[terverifikasi]` Uang di tabel itu bertipe `NUMBER` dengan **tiga skala berbeda**: `NUMBER(20,4)`
(nilai pokok), `NUMBER(20,8)` (premi & komisi coverage), `NUMBER(25,20)` (`RATE`, `PRORATE`).
Pembanding **menghormati ketiganya apa adanya** dan tidak membulatkan ke satu skala bersama.

⛔ **Seluruh kejanggalan yang diport apa adanya harus muncul sebagai selisih NOL**, bukan sebagai
selisih yang dimaafkan. Bila sistem baru "memperbaiki" salah satunya diam-diam, pembanding inilah yang
menangkapnya.

Kejanggalan Fac Out yang **wajib** menghasilkan selisih nol:

| Penanda | Isi | Tiket |
| --- | --- | --- |
| ~~`K046_RIComIN_TidakPernahDisetel`~~ | ⛔ **DICABUT (K-060)** — keanehan tidak ada pada salinan `DDL\` yang berlaku; ia akibat salinan korpus basi. **Bukan** kandidat perbaikan, **tidak** diuji sebagai keanehan | — |
| `K046_Prorate_EnamPenugasanSalingMenimpa` | satu properti ditimpa berurutan | F06 |
| `K046_BackUpStatus_DuaGayaPembanding` | `==3` telanjang vs `=="3"` berkutip | F05 |
| `K046_LocalFacout_AngkaSebagaiString` | `local.Facout=="0"` | F05 |
| `K046_IsFacRetroOffer_NolBerkutip` | `"0"` berkutip vs `0` telanjang | F05 |
| `K046_OkupasiFacOut_GayaUjiTidakSeragam` | awalan dan kesetaraan bercampur | F03 |
| `K046_GenerateFacRetroNo_TanpaElse` | `FacType` tak dikenal → kode kosong | F09 |
| `K046_GuardFacOut_HanyaType7` | gerbang idempotensi hanya aktif satu jenis penyesuaian | F11 |
| `K046_MenjadiSelisih_TidakSetangkup` | pasangan kolom tidak setangkup | F11 |

⛔ **Tidak pernah menyentuh berkas mentah ber-PII.** Fixture yang dipakai **sudah
ter-de-identifikasi**. Kolom `FACOUTPRODUCTION` yang ber-PII — `DOB`, `LICENSEPLATE`, `ENGINENUMBER`,
`ADDRESS`, `OBJECTNAME`, `REINSURER_NAME` — **tidak pernah** muncul di laporan pembanding, hanya
dihitung ada/tidaknya.

⚠️ **Uji dengan fixture, bukan dengan mengandalkan data produksi.** K-019 mencatat fitur fac out
**usang secara bisnis**, sehingga volume kasus yang melewati jalur ini **mungkin nol di data mutakhir**.
Bila paralel run tidak pernah menyentuh jalur Fac Out, itu **bukan bukti port-nya benar** — hanya
bukti jalurnya tidak terpakai. `[pertanyaan terbuka]` berapa banyak kasus Fac Out yang benar-benar ada
di data mutakhir — korpus tidak memuat data produksi, jadi angkanya hanya dapat datang dari DBA.

⚠️ **Dua presisi pembanding di F06 tidak diseragamkan** — rate pada 10 desimal, premi pada 4. Pembanding
menghormati keduanya apa adanya (ADR-F-0005).

⚠️ Selisih tipe dari **A.5** (string `"0"` → angka `0`) adalah **perbaikan sadar** yang sudah
dijelaskan K-046; pembanding **menormalkan** sebelum membandingkan dan **tidak** menandainya cacat.

**Asal (Pega).** — (pembanding, bukan port rule)

**Keputusan.** K-046 · K-019 · K-007 · K-025 · ADR-F-0001 · ADR-F-0005

**Blocked by:** F10 · `..\15-de-identifikasi-berkas-kasus.md` · `..\16-rekonsiliasi-eksak-tahap-1.md`
⛔ Kedua tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] Memakai **kerangka rekonsiliasi New Business**, tanpa pembanding baru
- [ ] ⛔ **Tidak pernah** membaca berkas mentah ber-PII; hanya fixture ter-de-identifikasi
- [ ] **Delapan** penanda K-046 yang masih berlaku menghasilkan **selisih nol** (semula sembilan; `K046_RIComIN_TidakPernahDisetel` dicabut oleh K-060)
- [ ] ⛔ Jalur Fac Out diuji dengan **fixture**, tidak mengandalkan ada-tidaknya kasus di data produksi
- [ ] Dua presisi pembanding (10 dan 4 desimal) dihormati apa adanya
- [ ] Setiap selisih yang tersisa **dapat ditelusuri** ke rule Pega asalnya (§4.6)
- [ ] Nilai `DOB`, plat nomor, nomor mesin, alamat, nama orang, dan nomor polis produksi **tidak pernah** muncul di laporan pembanding
