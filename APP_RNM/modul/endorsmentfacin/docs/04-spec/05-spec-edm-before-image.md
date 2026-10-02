# Spesifikasi — Modul Before-Image Endorsement (EDM)

> **Siklus:** Endorsement (EDM). **Modul pertama** yang dispec dari lingkup EDM (708 berkas, K-043).
>
> **Sumber bahan:** `07-edm\08-bahan-spec-before-image.md` — seluruh angka di sana dihitung dari
> korpus, bukan diperkirakan. Discovery pendukung: `07-edm\02-e1` … `07-edm\07-e6`.
>
> **Keputusan yang mengikat:** K-010/K-012 (uang = `Money`) · K-018 (skala rasio) · K-027 (koma
> desimal) · K-043 (lingkup 708) · K-044 (Life ikut; pilih-tertanggung dibuang) · K-045 · K-046
> (kode usang diport apa adanya, kecuali A.5).
>
> ✅ **Seam disetujui work owner 19 September 2026.** Satu seam baru — **Seam 4**
> `services/endorsement.PrepareBeforeImage` — mencakup ketiga lapis. Selisih **tidak** mendapat seam
> baru; ia diuji lewat **Seam 3** (`services/premium.Calculate`) yang sudah ada. Tiga seam NB yang
> disetujui 16 September 2026 **tidak berubah**. Total tetap **4 seam**.

---

## Problem Statement

Endorsement bukan penerbitan polis baru. Ia **mengubah polis yang sudah berjalan**, dan yang
dipertanggungjawabkan ke reasuradur bukan nilai penuh polis melainkan **selisih** antara keadaan
sesudah dan keadaan sebelum.

Selisih tidak dapat dihitung tanpa "keadaan sebelum". Di sistem lama, "keadaan sebelum" bukan satu
hal — ia **tiga mekanisme berbeda** yang disimpan di tempat berbeda, diisi oleh rule berbeda, pada
waktu berbeda, dan dikonsumsi oleh bagian berbeda. Ketiganya sering tertukar ketika dibaca sekilas,
karena semuanya "nilai lama".

Bila ketiganya disatukan atau salah satunya dilewatkan, angka selisih berubah — dan perubahan itu
**tidak akan terdeteksi sebagai galat**. Ia hanya muncul sebagai selisih yang tidak dapat dijelaskan
saat rekonsiliasi paralel run, setelah sistem baru sudah dipakai.

Underwriter dan bagian produksi tidak punya cara memverifikasi sendiri apakah "nilai lama" yang
dipakai sistem baru sama dengan yang dipakai Pega. Mereka hanya melihat angka akhir.

---

## Solution

Bangun modul **before-image** sebagai fondasi siklus endorsement: satu modul yang menyiapkan
"keadaan sebelum" dengan **tiga lapis yang tetap terpisah**, persis seperti sistem lama, lalu
menyediakan satu kontrak selisih di atasnya.

Ketiga lapis dipertahankan sebagai konsep terpisah di dalam model domain — **bukan** disederhanakan
menjadi satu "snapshot":

- **Lapis A — dokumen polis lama, utuh.** Versi terakhir polis dimuat dari basis data ke dalam
  agregat penawaran, lalu **51 penyalinan** membuat data kerja endorsement berangkat dari keadaan
  polis lama, bukan dari kosong. Inilah sumber angka yang dipakai perhitungan selisih produksi.
- **Lapis B — nilai lama per baris.** Untuk setiap objek, coverage, dan cedant, nilai lama disimpan
  sebagai properti bersaudara di dalam daftar yang sedang diedit, sehingga layar dapat menampilkan
  "dari → menjadi" berdampingan. Diisi ulang **setiap kali layar endorsement dibuka**.
- **Lapis C — penanda baris warisan.** Setiap simpul daftar yang berasal dari polis lama ditandai,
  sehingga baris lama dapat dibedakan dari baris yang ditambahkan saat endorsement.

Di atas ketiganya berdiri **satu kontrak selisih** yang tunggal dan eksplisit: nilai baru diprorata
menurut sisa periode, lalu dikurangi nilai lama yang jenis treaty-nya sama.

Hasilnya: underwriter melihat "dari → menjadi" di layar, bagian produksi memperoleh selisih yang
dapat direkonsiliasi baris demi baris terhadap Pega, dan tim migrasi punya satu titik uji tunggal
untuk membuktikan keduanya.

---

## User Stories

### Lapis A — dokumen polis lama

1. Sebagai **underwriter endorsement**, saya ingin formulir endorsement terbuka sudah berisi data
   polis yang berjalan, sehingga saya hanya mengubah yang perlu diubah dan tidak mengetik ulang
   seluruh polis.
2. Sebagai **underwriter endorsement**, saya ingin identitas bisnis polis lama (kode dan nama
   bisnis, ceding, sumber bisnis, marketing, tertanggung, nomor slip penawaran, tipe polis) terbawa
   apa adanya, sehingga endorsement tetap melekat pada polis yang benar.
3. Sebagai **underwriter endorsement**, saya ingin periode polis lama (tanggal mulai, tanggal
   berakhir, tanggal penawaran, tanggal produksi) terbawa, sehingga prorata endorsement dihitung
   terhadap periode yang benar.
4. Sebagai **underwriter endorsement**, saya ingin struktur share dan kapasitas treaty polis lama
   (mata uang, daftar mata uang, skala inward, persentase share, kapasitas treaty, tahun berjalan,
   tipe prorata, daftar cedant) terbawa, sehingga perubahan saya dihitung di atas struktur yang sama.
5. Sebagai **underwriter endorsement**, saya ingin ketentuan pembayaran polis lama (angsuran, komisi
   reasuransi, persentase brokerage) terbawa, sehingga selisih komisi dihitung terhadap dasar yang
   benar.
6. Sebagai **bagian produksi**, saya ingin nilai lama yang dipakai menghitung selisih diambil dari
   dokumen polis lama yang utuh, bukan dari nilai per baris, sehingga angka yang masuk tabel produksi
   dapat direkonsiliasi terhadap Pega.
7. Sebagai **tim migrasi**, saya ingin versi polis yang dimuat adalah versi terakhir menurut aturan
   yang sama dengan sistem lama, sehingga before-image tidak pernah mengambil versi yang salah.
8. Sebagai **tim migrasi**, saya ingin penyalinan ke halaman selain agregat penawaran dilakukan
   eksplisit, sehingga halaman-halaman itu tidak kosong di sistem baru.

### Lapis B — nilai lama per baris

9. Sebagai **underwriter endorsement**, saya ingin melihat nilai lama dan nilai baru berdampingan
   pada tiap baris objek, sehingga saya tahu apa yang sedang saya ubah.
10. Sebagai **underwriter endorsement**, saya ingin nilai lama diperbarui setiap kali saya membuka
    kembali layar endorsement, sehingga tampilan tidak basi terhadap keadaan tersimpan.
11. Sebagai **underwriter kebakaran**, saya ingin nilai lama tersedia untuk TSI objek, premi bruto,
    premi Nusantara Re, TSI total, serta TSI, premi, dan rate pada ringkasan, sehingga perubahan
    pada tingkat objek maupun ringkasan sama-sama terlihat.
12. Sebagai **underwriter kendaraan bermotor**, saya ingin nilai lama tersedia untuk TSI, premi,
    premi Nusantara Re, premi rupiah, dan premi bruto setelah diskon armada, sehingga struktur
    diskon armada ikut terbandingkan.
13. Sebagai **underwriter marine cargo**, saya ingin nilai lama tersedia untuk TSI, premi, dan premi
    Nusantara Re per coverage, sehingga perubahan pengangkutan terlihat per coverage.
14. Sebagai **underwriter aneka, golf, dan personal accident**, saya ingin nilai lama tersedia untuk
    TSI dan premi pada tingkat objek maupun ringkasan, sehingga perubahan terlihat di kedua tingkat.
15. Sebagai **underwriter travel**, saya ingin nilai lama tersedia untuk TSI dan premi per coverage
    peserta, sehingga perubahan peserta terlihat.
16. Sebagai **underwriter**, saya ingin nilai lama pada tingkat cedant tersedia untuk TSI dan premi
    di semua lini bisnis, sehingga pembagian ke cedant ikut terbandingkan.
17. Sebagai **underwriter**, saya ingin nilai lama yang tidak ada di polis lama ditampilkan sebagai
    nol, bukan kosong, sehingga aritmetika selisih tidak menghasilkan hasil yang tidak terdefinisi.

### Lapis C — penanda baris warisan

18. Sebagai **underwriter endorsement**, saya ingin baris yang berasal dari polis lama dapat
    dibedakan dari baris yang saya tambahkan saat endorsement, sehingga saya tahu mana yang warisan.
19. Sebagai **bagian produksi**, saya ingin penanda warisan tersedia di setiap tingkat daftar sesuai
    kedalaman model data lini bisnisnya, sehingga penandaan tidak hilang pada daftar bersarang.

### Selisih

20. Sebagai **bagian produksi**, saya ingin selisih dihitung dengan satu rumus yang sama di seluruh
    lini bisnis, sehingga tidak ada lini yang diam-diam memakai aturan sendiri.
21. Sebagai **bagian produksi**, saya ingin nilai baru dikalikan porsi sisa periode sebelum
    dikurangi nilai lama, sehingga endorsement di tengah periode hanya membebankan porsi yang
    tersisa.
22. Sebagai **bagian produksi**, saya ingin pasangan nilai lama dan nilai baru dicocokkan menurut
    jenis treaty, bukan menurut urutan baris, sehingga penambahan atau penghapusan baris tidak
    menggeser pasangan.
23. Sebagai **aktuaris**, saya ingin porsi periode dihitung dengan presisi desimal tinggi dan
    terlindung dari pembagian nol, sehingga pembulatan tidak menggerus angka premi.
24. Sebagai **aktuaris**, saya ingin varian perhitungan yang menyimpang dari rumus umum
    (penyesuaian spreading, perubahan mata uang, penyesuaian rate, coverage yang dihapus,
    pembatalan) tetap berperilaku seperti sistem lama, sehingga rekonsiliasi paralel run bersih.

### Batas dan gerbang

25. Sebagai **tim migrasi**, saya ingin persiapan nilai lama per baris hanya berjalan untuk kasus
    endorsement, sehingga siklus lain tidak ikut terpengaruh.
26. Sebagai **tim migrasi**, saya ingin endorsement jiwa melewati mekanisme nilai lama per baris dan
    memakai jalurnya sendiri, sehingga perilakunya sama dengan sistem lama.
27. Sebagai **tim migrasi**, saya ingin polis dengan lokasi sangat banyak tetap berperilaku seperti
    sistem lama — nilai lama per baris tidak diisi — sehingga tidak ada perbedaan yang tak dapat
    dijelaskan.
28. Sebagai **tim migrasi**, saya ingin predikat "bukan endorsement" tidak diimplementasikan sebagai
    kebalikan logis dari predikat "endorsement", sehingga perbedaan sumber data keduanya
    dipertahankan.

### Tipe dan ketelitian

29. Sebagai **aktuaris**, saya ingin setiap nilai lama yang berupa uang membawa mata uangnya sendiri,
    sehingga selisih antar mata uang tidak pernah dijumlahkan diam-diam.
30. Sebagai **aktuaris**, saya ingin rate lama membawa skalanya sendiri sesuai lini bisnis, sehingga
    per-mil tidak pernah tertukar dengan persen.
31. Sebagai **pengembang**, saya ingin penjumlahan antara nilai uang dan rasio **gagal saat
    kompilasi**, sehingga kekeliruan satuan tidak pernah sampai ke produksi.
32. Sebagai **aktuaris**, saya ingin nilai desimal berkoma dari sistem lama terbaca sesuai kontrak
    yang sudah disepakati, sehingga tidak ada nilai yang salah terbaca faktor seribu.

### Rekonsiliasi dan ketertelusuran

33. Sebagai **tim rekonsiliasi**, saya ingin setiap penyalinan, guard, dan rumus menyebut rule Pega
    asalnya, sehingga selisih dapat ditelusuri ke sumbernya.
34. Sebagai **tim rekonsiliasi**, saya ingin kejanggalan sistem lama yang sengaja dipertahankan
    tercatat sebagai keputusan, sehingga tidak ada yang "memperbaikinya" di kemudian hari.
35. Sebagai **tim rekonsiliasi**, saya ingin satu-satunya perbaikan sadar di modul ini terdokumentasi
    beserta dampaknya pada pembandingan, sehingga selisih tipe yang muncul dapat dijelaskan dan
    bukan dikira bug.

---

## Implementation Decisions

### 1. Tiga lapis tetap tiga, bukan satu

Model domain mempertahankan ketiganya sebagai konsep terpisah. Menyatukannya menjadi satu "snapshot"
akan menghapus perbedaan konsumen yang sudah terverifikasi: **selisih produksi membaca lapis A**,
sementara **layar dan perhitungan per baris membaca lapis B**. Kedua angka itu tidak selalu sama.

| Lapis | Isi | Siklus hidup | Konsumen |
| --- | --- | --- | --- |
| A | dokumen polis lama utuh | sekali, saat kasus endorsement lahir | perhitungan selisih produksi |
| B | nilai lama per baris | **setiap kali layar dibuka** | layar + perhitungan per baris |
| C | penanda baris warisan | sekali, saat kasus lahir | pembeda baris lama vs baru |

Asal: `SetValueToEDMWork` langkah 14 (A dan C), `SetOldData` (B), tujuh varian
`SetOLDValueToEDMWork_<LOB>` (C).

⚠️ Istilah "before-image" adalah **nama konsep di dokumen ini**. Tidak ada properti bernama itu di
korpus — jangan menamai field demikian.

### 2. Modul yang dibangun

| Modul | Isi | Seam |
| --- | --- | --- |
| `services/endorsement` | **baru** — penyiapan before-image tiga lapis + kontrak selisih | **Seam 4** (§Testing) |
| `services/premium` | **diperluas** — menerima masukan endorsement (nilai lama + rasio prorata) | Seam 3 yang sudah ada, tidak berubah bentuknya |
| `pkg/money`, `pkg/ratio` | **dipakai ulang tanpa perubahan** | tidak punya seam sendiri (keputusan NB) |
| `internal/rules` | **dipakai ulang** — predikat lini bisnis dan gerbang siklus | Seam 1 yang sudah ada |

`services/acceptance` **tidak disentuh** di spec ini.

### 3. Kontrak lapis A — penyalinan yang tidak boleh diringkas

Pemuatan dokumen polis lama diikuti **54 penugasan** yang harus diport eksplisit:

- **51** bersumber dari dokumen polis lama;
- **3** tidak bersumber dari dokumen lama (penyelarasan internal antar-halaman).

Dari 51 itu, **45** mengikuti pola "salin field yang namanya sama", dan **6 menyimpang** — menulis ke
halaman selain agregat penawaran: tipe bisnis ke halaman kutipan, tiga field pembayaran dicerminkan
ke halaman polis, penanda B2B ke akar objek kerja, dan nama marketing ke halaman parameter kredit.

⛔ **Keenam penyimpangan wajib disebut eksplisit.** Arsip lintas-siklus menyatakan seluruh 54
berbentuk seragam; itu **tidak berlaku untuk 9 dari 54**. Implementasi yang mengikuti arsip akan
meninggalkan halaman-halaman itu kosong.

Empat kelompok field dari 45 yang berpola seragam:

| Kelompok | Cacah |
| --- | ---: |
| Identitas bisnis | **20** |
| Periode & tanggal | **4** |
| Struktur share & kapasitas | **18** |
| Pembayaran | **3** |

**Akibat yang mengikat:** setelah penyalinan, data kerja endorsement **identik** dengan polis lama.
Implementasi **tidak boleh** memulai dari struktur kosong lalu menambahkan perubahan; ia harus
menyalin lapis A lebih dulu, baru menerima perubahan pengguna. Delta dihitung terhadap polis lama,
bukan terhadap kosong.

### 4. Kontrak lapis B — 52 penugasan, tujuh lini, pola nilai tunggal

Nilai lama per baris diisi untuk **tujuh lini**: kebakaran, golf, aneka, personal accident, marine
cargo, kendaraan bermotor, travel — masing-masing berpasangan (tingkat objek + tingkat cedant).
Tingkat cedant identik di ketujuh lini: TSI lama dan premi lama.

**Pola nilai tunggal, 52 dari 52, tanpa pengecualian:**

```
<properti>Old = jika <properti sumber> tidak kosong → <properti sumber>, selain itu → 0
```

⛔ **Kosong dipetakan ke nol, bukan ke nilai tak-ada.** Di Go, properti `*Old` yang tidak terisi
bernilai nol, bukan pointer nil. Aritmetika selisih bergantung pada ini.

Premi Nusantara Re lama hanya ada pada **dua lini** — marine cargo dan kendaraan bermotor — bukan
pada lima lini lainnya.

### 5. Dua guard yang sengaja dipertahankan (K-046)

Dua guard di sistem lama menguji **properti tujuan**, bukan properti sumber:

- **Rate lama** — pada **6 dari 6** kemunculan;
- **Premi lama** — hanya pada lini kebakaran; **17 kemunculan lain** memakai bentuk yang menguji
  properti sumber.

⛔ **Diport apa adanya.** Work owner menetapkan ini **perilaku yang benar**, bukan cacat.
Meluruskannya akan mengubah angka selisih rate dan memutus rekonsiliasi. Keduanya wajib disertai
komentar yang menyebut rule Pega asalnya dan menunjuk keputusan K-046, agar tidak "diperbaiki" oleh
pembaca berikutnya.

### 6. ⚠️ A.5 — satu-satunya perbaikan sadar di modul ini

Premi Nusantara Re lama memakai nilai cadangan berupa **string** `"0"` pada dua penugasannya,
sementara **50 penugasan lain** memakai **angka** `0`.

⛔ **Di sistem baru: angka `0`.** Premi Nusantara Re adalah nilai uang (K-010/K-012), dan tipe uang
tidak boleh menampung string.

**Ini perbaikan, bukan port.** Satu-satunya di modul ini, dan wajib ditandai demikian di kode.

**Konsekuensi rekonsiliasi — wajib masuk prosedur pembandingan:** pada kasus premi Nusantara Re
kosong di marine cargo atau kendaraan bermotor, sistem lama menyimpan string dan sistem baru
menyimpan angka. Pembandingan mentah akan menunjukkan **beda tipe meski nilainya sama**. Pembanding
rekonsiliasi harus menormalkan sebelum membandingkan, dan selisih semacam ini **dijelaskan oleh
A.5** — bukan ditandai sebagai cacat.

### 7. Tipe data

Penerapan keputusan yang sudah ada, bukan keputusan baru:

| Kelompok | Tipe |
| --- | --- |
| Seluruh nilai lama yang berupa uang — TSI, TSI objek, premi, premi bruto, premi Nusantara Re, premi rupiah, premi bruto setelah diskon armada | **`Money{Amount decimal, Currency}`** (K-010/K-012) |
| Rate lama | **`Ratio{Value decimal, Scale}`**, skala per lini mengikuti K-018 |
| Penanda baris warisan | string, nilai literal tetap |
| Porsi periode sebelum dan sesudah tanggal endorsement | **`Ratio`**, presisi desimal tinggi |

⛔ `Money` dan `Ratio` **tidak dapat dijumlahkan**; satu-satunya jembatan **`Money × Ratio → Money`**.
Pelanggaran wajib **gagal saat kompilasi** (kontrak NB, dipakai ulang di sini).

Nilai desimal berkoma dari sistem lama dibaca mengikuti kontrak K-027.

### 8. Kontrak selisih

Rumus kanonik:

> **`SELISIH = (nilai_baru × porsi_sisa_periode) − nilai_lama_dengan_jenis_treaty_sama`**

Tiga keputusan yang melekat padanya:

1. **Nilai lama diambil dari lapis A**, bukan dari nilai lama per baris (lapis B). Ini titik yang
   paling mudah salah diimplementasikan.
2. **Pasangan dicocokkan menurut jenis treaty**, bukan menurut indeks posisi. Penambahan atau
   penghapusan baris tidak boleh menggeser pasangan.
3. **Porsi periode dipakai sebagai pengali nilai baru**, bukan penjumlah. Dua rasio berbeda peran:
   porsi **sesudah** tanggal endorsement adalah pengali umum; porsi **sebelum** hanya dipakai pada
   satu jenis penyesuaian (penyesuaian rate).

Perhitungan porsi periode wajib menyertakan **guard pembagian nol** (periode nol diperlakukan sebagai
satu) dan presisi desimal tinggi.

#### ⚠️ 8.1 Porsi periode yang disetel modul ini bersifat SEMENTARA

> Ditambahkan atas **K-048 §8.1**, setelah bahan spec modul selisih menemukannya.

`[terverifikasi]` Properti porsi-sisa-periode (`OfferFacIn.ProrateEDMEnd`) yang disetel modul ini
**bukan nilai final**. Ia **ditulis oleh lima rule (12 penugasan)** dan **dibaca oleh 13 rule** di
seluruh korpus EDM. Modul selisih **menghitung ulang dan menimpanya** sebelum jalur produksi
membacanya:

| Penulis | Perlakuan |
| --- | --- |
| **Modul ini** — `SetValueToEDMWork` langkah 15 | dihitung dari selisih **DateTime** periode polis lama |
| `CountPremiEDM_DT` langkah 9 | **dihitung ulang** dari selisih **hari** (`@DateTimeDifference(…,"D")`) lalu **ditimpa** |
| `CountPremiEDM_DT` langkah 11.3 | dipaksa **`1`** bila `EdmType==1` |
| `CountPremiEDM_DT` langkah 12.1 | ditulisi **rasio porsi AWAL** bila `EdmType==2` |

⛔ **Keputusan K-048: diport apa adanya** — satu properti yang ditimpa berurutan, bukan dipecah
menjadi beberapa rasio bernama berbeda.

⛔ **Konsekuensi yang mengikat implementasi:** **urutan eksekusi modul before-image → modul selisih
WAJIB dijaga.** Membalik urutannya, menjalankan keduanya paralel, atau meng-cache nilai porsi periode
dari modul ini akan menghasilkan angka yang berbeda.

⚠️ Spec modul ini **tetap** mendefinisikan perhitungan porsi periode sebagaimana di §8 — itu yang
dilakukan sistem lama di titik ini. Yang ditambahkan catatan ini hanyalah: **nilainya tidak boleh
diperlakukan sebagai final oleh konsumen mana pun.**

### 9. Lima varian yang menyimpang dari rumus kanonik — diport apa adanya

| Kondisi | Perlakuan |
| --- | --- |
| Penyesuaian spreading | selisih **sama dengan** nilai sesudah — seluruh nilai dianggap baru |
| Penyesuaian mata uang dengan mata uang berubah | idem |
| Penyesuaian rate | porsi sebelum tanggal endorsement ikut ditambahkan sebelum pengurangan |
| Coverage hilang dari data baru | selisih = nilai lama dikali negatif satu |
| Pembatalan (tiga kode jenis endorsement) | selisih persentase komisi, brokerage, dan potongan **dibalik tandanya** |

### 10. Tiga gerbang keluar dan konsekuensinya

Persiapan nilai lama per baris berhenti sebelum berjalan bila salah satu terpenuhi:

| Gerbang | Konsekuensi desain |
| --- | --- |
| **Kasus jiwa** | Endorsement jiwa **melewati lapis B sepenuhnya** dan memakai jalurnya sendiri. Per **K-044 Life IKUT lingkup** — jadi ini **alur tersendiri yang dimodelkan eksplisit**, bukan cabang kondisi di dalam lapis B |
| **Bukan endorsement** | Predikat siklus membaca agregat tersimpan. ⛔ Predikat "bukan endorsement" membaca **halaman aktif** — sumber berbeda, sehingga **tidak boleh** diimplementasikan sebagai negasi logis. Keduanya bisa benar atau salah bersamaan |
| **Lokasi melebihi ambang** | Polis besar **tidak mendapat nilai lama per baris sama sekali**; layar dan perhitungan per baris menampilkan nol. Diport apa adanya; alasan ambangnya tidak diketahui dari korpus |

### 11. Cakupan lini berbeda antar lapis

Lapis B menangani **travel tetapi bukan jiwa**; lapis C menangani **jiwa tetapi bukan travel**.
Ini bukan kekeliruan pembacaan — himpunannya memang berbeda, dan implementasi harus mengikutinya.

Dua kekhususan lapis C yang wajib diport:

- Hanya varian **kebakaran** yang menyetel penanda prorata tingkat agregat. Enam varian lain tidak,
  dan penanda itu **tidak** termasuk field yang disalin lapis A.
- Varian **jiwa** menyetel **dua** penanda berbeda; enam varian lain hanya satu.

### 12. Ketertelusuran

Setiap penyalinan, guard, rumus, dan gerbang menyebut **rule Pega asalnya** dalam komentar
(`CLAUDE.md` §4.6). Tanpa itu rekonsiliasi paralel run tidak mungkin.

---

## Testing Decisions

### Apa yang membuat test baik di proyek ini

**Hanya perilaku eksternal lewat seam yang disetujui.** Test tidak menyentuh fungsi internal, tidak
memeriksa urutan pemanggilan, dan tidak menegaskan bentuk struct antara. Yang diuji: masukan →
keluaran pada satu titik tertinggi.

Alasannya khas migrasi: test di sini berfungsi sebagai **alat rekonsiliasi**, bukan hanya jaring
pengaman regresi. Test yang terikat implementasi akan menghalangi penyesuaian ketika ekspor produksi
tunggal akhirnya tiba.

### Seam

Tiga seam NB yang disetujui 16 September 2026 **tidak berubah**:

1. `internal/rules.Eval` — registry predikat
2. `services/acceptance.Next` — satu transisi tangga
3. `services/premium.Calculate` — satu pintu masuk perhitungan

**Diusulkan satu seam baru — Seam 4:**

```
services/endorsement.PrepareBeforeImage(polisLama, tanggalEndorsement) → AgregatPenawaran, error
```

Satu pintu masuk yang mencakup **ketiga lapis sekaligus**. Lapis A, B, dan C **bukan** seam
tersendiri — memberi masing-masing seam akan mengunci pembagian internal yang belum tentu bertahan.

**Selisih tidak mendapat seam baru.** Ia diuji lewat **Seam 3 yang sudah ada**
(`services/premium.Calculate`), diperluas menerima masukan endorsement. Ini menjaga jumlah seam tetap
minimum dan menempatkan titik rekonsiliasi premi di satu tempat, konsisten dengan keputusan NB.

✅ **Seam 4 disetujui work owner 19 September 2026, tanpa perubahan.** Tidak ada seam lain yang boleh
ditambahkan untuk modul ini tanpa keputusan baru.

### Yang diuji di Seam 4

**Table-driven per lini bisnis** (tujuh lini untuk lapis B, tujuh untuk lapis C, himpunannya
berbeda). Kasus wajib:

- Lapis A memuat versi polis yang benar, dan **54 penugasan** terjadi — termasuk **keenam
  penyimpangan** ke halaman lain. Kasus negatif: halaman-halaman itu **tidak boleh** kosong.
- Data kerja setelah penyalinan **identik** dengan polis lama, bukan kosong.
- Lapis B mengisi properti nilai lama yang benar **per lini**, dan hanya per lini itu. Premi
  Nusantara Re lama muncul **hanya** pada dua lini.
- Nilai kosong menjadi **nol**, bukan nilai tak-ada — untuk **seluruh** properti nilai lama.
- **Kedua guard yang menguji properti tujuan berperilaku seperti sistem lama** — ini kasus uji
  eksplisit, bukan efek samping. Test-nya menegaskan hasil yang "salah menurut intuisi" dan menyebut
  K-046 dalam namanya, agar tidak ada yang mengubahnya tanpa membaca keputusan.
- **A.5**: premi Nusantara Re kosong menghasilkan **angka** nol. Disertai satu test rekonsiliasi yang
  menegaskan pembanding menormalkan string dan angka sebelum membandingkan.
- Ketiga gerbang keluar: kasus jiwa, kasus bukan-endorsement, kasus lokasi melebihi ambang —
  masing-masing menghasilkan lapis B **kosong**, dan itu perilaku yang benar.
- Predikat "bukan endorsement" **tidak** setara dengan negasi predikat "endorsement": kasus uji
  dengan kedua halaman tidak sinkron, yang membuktikan keduanya bisa bernilai sama.
- Lapis C menandai simpul sesuai **kedalaman daftar** lini bisnisnya, dan dua kekhususan lini
  kebakaran serta lini jiwa.

### Yang diuji di Seam 3 untuk endorsement

- Rumus kanonik: nilai baru diprorata lalu dikurangi nilai lama **berjenis treaty sama**.
- **Nilai lama berasal dari lapis A**, bukan lapis B — kasus uji dengan keduanya sengaja dibuat
  berbeda, memastikan implementasi mengambil yang benar.
- Pencocokan menurut jenis treaty, bukan indeks: kasus dengan baris ditambah dan dihapus di tengah.
- Kelima varian menyimpang, satu kasus masing-masing.
- Guard pembagian nol pada periode.
- Presisi desimal porsi periode tidak tergerus pembulatan.

### Pemeriksaan waktu-kompilasi

Kontrak NB dipakai ulang tanpa perubahan: penjumlahan `Money` dengan `Ratio` **gagal saat
kompilasi**. Modul ini menambah nilai lama sebagai `Money` dan porsi periode sebagai `Ratio`, jadi
pemeriksaan itu kini melindungi lebih banyak titik.

### Prior art

Spec NB (`03-spec-modul-terverifikasi.md`) dan spec RNW (`04-spec-rnw.md`). Bentuk test di ketiga
seam NB **menjadi** prior art — test modul ini harus mengikuti bentuknya, bukan menciptakan gaya
baru.

---

## Out of Scope

1. **Jalur produksi endorsement** — penulisan ke tabel produksi fac in dan fac out, serta penomoran
   versi polis. Ditunda menunggu tabel flat dan sumber prosedur tersimpan dari DBA. Modul ini
   **menghasilkan** selisih; menuliskannya adalah modul lain.
2. **Tangga akseptasi endorsement.** Termasuk klaim arsip bahwa nilai dasar akseptasi endorsement
   adalah selisih TSI — klaim itu **belum diuji ke korpus** dan **tidak boleh** dijadikan dasar di
   spec ini.
3. **Modul premi endorsement** — sebelas bentuk perhitungan premi-menjadi di tujuh berkas. Modul
   tersendiri, meskipun memakai Seam 3 yang sama.
4. **Fitur pilih-tertanggung** — **K-044: DIBUANG**. Keempat berkasnya ada di korpus tetapi tidak
   diport.
5. **Jalur perbaikan pemuatan dokumen lama** ketika kode bisnis kosong. Ada di korpus, belum dibedah
   langkah demi langkah.
6. **Lapisan tampilan.** Spec ini menyediakan data "dari → menjadi"; merancang layarnya terpisah.
   Hanya dua section di korpus yang menampilkan nilai lama, padahal lapis B mengisi 52 properti —
   apakah sisanya memang hanya untuk perhitungan masih pertanyaan terbuka.
7. **Repository Oracle.** Modul ini bicara tentang agregat di memori. Kontrak `panic` saat menulis
   mata uang tak dikenal tetap menjadi kewajiban yang belum punya tempat, sama seperti di spec NB.

---

## Further Notes

### Mengapa modul ini didahulukan

Dari seluruh lingkup EDM, before-image adalah satu-satunya bagian besar yang **tidak menunggu siapa
pun**. Mekanismenya terbaca penuh dari korpus; yang tertahan di DBA adalah penulisan ke tabel
produksi, bukan perhitungan nilainya. Mendahulukannya berarti pekerjaan berjalan sementara paket DBA
dan IT masih diproses.

Ia juga **fondasi**: modul premi endorsement, tangga akseptasi endorsement, dan jalur produksi
semuanya membaca keluarannya. Salah di sini menjalar ke semuanya.

### Angka di spec ini dihitung, bukan diperkirakan

Setiap cacah — 54, 51, 45, 6, 3, 52, 2 dari 52 — berasal dari penghitungan langsung atas korpus,
bukan dari ringkasan dokumen sebelumnya. Bahan beserta perintah auditnya ada di
`07-edm\08-bahan-spec-before-image.md`.

Ini bukan kehati-hatian berlebihan. Sepanjang discovery EDM, **empat angka pernah salah** karena
diwarisi dari ringkasan alih-alih dihitung ulang, dan dua di antaranya sempat menjadi dasar
keputusan sebelum dikoreksi.

### Arsip lintas-siklus keliru pada satu titik di modul ini

Arsip menyatakan seluruh 54 penyalinan lapis A berbentuk seragam "salin field bernama sama".
Penghitungan ulang menunjukkan **9 dari 54 tidak demikian**. Implementasi yang mengikuti arsip akan
meninggalkan empat halaman kosong. Koreksinya sudah masuk §Implementation Decisions butir 3.

### Perbaikan dipisahkan dari migrasi — dan di sini pengecualiannya tunggal

Modul ini memuat kejanggalan yang jelas terlihat: dua guard yang menguji properti tujuan alih-alih
properti sumber. **Keduanya diport apa adanya** atas keputusan work owner (K-046).

Hanya **satu** perbaikan sadar dilakukan — A.5, string nol menjadi angka nol — dan itu pun karena
tipe uang tidak dapat menampung string. Perbaikan itu ditandai di kode, di spec, dan di prosedur
rekonsiliasi.

Aturannya tetap: migrasi yang diam-diam "memperbaiki" menghasilkan selisih angka yang tidak dapat
dijelaskan saat paralel run.

### Tentang penamaan di dokumen ini

Nama modul (`services/endorsement`, `services/premium`) adalah **struktur target**, belum ada. Nama
rule Pega disebut di bahan dan wajib masuk komentar kode demi ketertelusuran, tetapi **tidak** menjadi
nama tipe atau fungsi di sistem baru. Nama rule bukan bukti perilaku — sepanjang proyek ini tujuh
keluarga jebakan penamaan sudah terbukti menghasilkan kesimpulan salah.

### Yang membuat spec ini dapat dikerjakan sekarang

Seluruh perilaku yang dispec sudah `[terverifikasi]` ke korpus dengan path, langkah, dan isi tag.
Empat pertanyaan terbuka yang tersisa menyangkut **alasan** — mengapa ambang lokasi sekian, mengapa
hanya satu lini menyetel penanda prorata — bukan **perilaku**. Perilakunya terbaca dan dapat diport
tanpa menunggu jawabannya.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
