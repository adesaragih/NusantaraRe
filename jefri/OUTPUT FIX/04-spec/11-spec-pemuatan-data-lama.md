# Spec 11 — Program pemuatan data lama ke tabel flat

> **25 September 2026.** Disusun lewat `/to-spec` atas permintaan work owner.
>
> **Label triase:** `ready-for-agent`.
>
> ⛔ **Penyimpangan tempat terbit, disengaja.** `docs/agents/issue-tracker.md` menempatkan tiket dan
> spec di `.scratch\<slug>\`. Work owner **melarang menulis ke `.scratch\`**; seluruh output hanya ke
> `RNM\OUTPUT\`. Spec ini karena itu terbit di `04-spec\`, mengikuti sepuluh spec sebelumnya.
>
> ⛔ **Penyimpangan kedua, juga disengaja.** Templat `/to-spec` melarang mencantumkan path berkas.
> Larangan itu ditujukan pada **berkas kode sasaran**, yang cepat basi — dan dipatuhi penuh di sini.
> Rujukan ke **korpus Pega dan DDL** tetap dicantumkan karena `CLAUDE.md` §3 dan §4.6 **mewajibkannya**:
> tanpa itu rekonsiliasi paralel run tidak dapat ditelusuri. Korpus bersifat read-only dan tidak basi.
>
> **Dasar keputusan:** K-063 · K-066 · K-067 · K-068 · K-069 · **K-070** · **K-071** · **K-072** ·
> **K-073** · ADR-0001 · ADR-0004 · ADR-0005 · ADR-0006 · ADR-0007 · V-16 · V-19 · V-23 · V-25 ·
> V-27 · V-28 · V-40 · V-41 · V-42 · V-43 · V-45.
> **Bahan:** `08-flat\BAHAN-SPEC-PEMUATAN.md` (revisi 25 September).

---

## Problem Statement

Seluruh data bisnis sebuah **penawaran** (`OfferFacIn`) hari ini tersimpan sebagai **satu dokumen
JSON utuh** di satu kolom CLOB. Bentuk itu ikut mati bersama Pega: sistem baru membaca dan menulis
**78 tabel flat** (1.329 kolom, 220+ relasi), bukan satu dokumen.

Akibatnya, tanpa program pemuatan:

- Sistem baru berdiri **kosong**. Tidak ada satu pun polis lama yang dapat dibuka, di-endors, atau
  di-renew di dalamnya.
- **Rekonsiliasi paralel run mustahil.** ADR-0001 menuntut nol selisih sampai digit terakhir antara
  sistem lama dan sistem baru; tanpa data yang sama di kedua sisi, tidak ada yang dapat dibandingkan.
- Pengujian hanya dapat memakai data karangan, yang `[terverifikasi]` tidak mewakili keadaan nyata —
  korpus memuat TSI puluhan miliar, nilai berkoma desimal, dan baris tanpa mata uang.

Pemindahannya bukan penyalinan kolom-ke-kolom. Dokumen sumber adalah **pohon bersarang sampai
delapan tingkat** dengan wadah berulang; sasarannya tabel relasional. Satu penawaran rata-rata
menjadi **±197 baris**, dan yang terbesar **13.249 baris**.

---

## Solution

Sebuah program pemuatan yang mengubah satu dokumen penawaran menjadi himpunan baris untuk 78 tabel
flat, dijalankan dalam **dua fase** sesuai keputusan work owner **K-070** dan **K-073**.

| | **Fase 1** | **Fase 2** |
| --- | --- | --- |
| Populasi | **versi polis terakhir** tiap polis | ⛔ **seluruh versi polis**, dalam satu jalan |
| Tabel flat | dimuat dari keadaan kosong | ⛔ **dikosongkan lebih dulu**, lalu dimuat penuh |
| Kapan | bersamaan dengan pembangunan sistem | **sambil menguji**, sebelum sistem dipakai sungguhan |
| Identitas baris | dibangkitkan, ⛔ **bersifat sementara** | dibangkitkan **sekali**, lalu **beku** |
| Rekonsiliasi ADR-0001 | tidak ditegakkan sebagai putusan akhir | ⛔ **di sinilah** nol selisih ditegakkan |

Bentuk dua fase ini dipilih karena membuat keputusannya **dapat dibatalkan**: selama fase 2 belum
lewat, tabel masih boleh dikosongkan, sehingga identitas baris tidak pernah perlu **dijodohkan ke
belakang**. Pilihan sebaliknya — menjodohkan identitas baris lama ke identitas yang sudah beku —
adalah pintu satu arah yang kesalahannya **tidak bersuara**.

Dari sisi orang yang memakainya, hasilnya: polis lama terbuka di sistem baru, endorsement berikutnya
melanjutkan penomoran yang benar, dan angka di kedua sistem dapat diadu satu lawan satu.

---

## User Stories

### Menjalankan pemuatan

1. Sebagai **insinyur pemuatan**, saya ingin menjalankan fase 1 atas seluruh polis dan mendapat
   laporan berisi jumlah polis diproses, jumlah baris per tabel, dan jumlah kegagalan, supaya saya
   tahu pemuatan berhasil tanpa harus memeriksa tabel satu per satu.
2. Sebagai **insinyur pemuatan**, saya ingin pemuatan satu penawaran berlangsung dalam **satu
   transaksi**, supaya sebuah penawaran tidak pernah tersimpan setengah jadi.
3. Sebagai **insinyur pemuatan**, saya ingin kegagalan pada satu penawaran **tidak menghentikan**
   penawaran lain, supaya satu dokumen rusak tidak membatalkan pemuatan sepuluh ribu dokumen lain.
4. Sebagai **insinyur pemuatan**, saya ingin daftar penawaran yang gagal beserta **alasannya**,
   supaya saya dapat memperbaikinya dan memuat ulang hanya yang gagal.
5. Sebagai **insinyur pemuatan**, saya ingin menjalankan pemuatan atas **satu penawaran tertentu**,
   supaya saya dapat menelusuri satu kasus tanpa memproses seluruh populasi.
6. Sebagai **insinyur pemuatan**, saya ingin fase 2 **mengosongkan tabel flat lebih dulu** lalu
   memuat seluruh versi polis dalam satu jalan, supaya identitas baris lahir dari pandangan atas
   riwayat yang lengkap.
7. Sebagai **insinyur pemuatan**, saya ingin pemuatan bersifat **dapat diulang** — menjalankannya dua
   kali atas sumber yang sama menghasilkan isi tabel yang sama — supaya percobaan berulang saat
   pengujian tidak menumpuk baris ganda.

### Memilih versi polis

8. Sebagai **work owner**, saya ingin fase 1 memuat **versi polis terakhir** tiap polis, supaya
   sistem baru dapat diuji lebih cepat tanpa menunggu seluruh riwayat dipindahkan.
9. Sebagai **insinyur pemuatan**, saya ingin versi terakhir ditentukan oleh **tanggal masuk
   terbaru**, bukan oleh nomor versi, supaya pengurutan tidak terjebak perbandingan teks.
10. Sebagai **tim rekonsiliasi**, saya ingin nomor versi tetap **diperiksa silang** terhadap tanggal
    masuk, dan pemuatan **berhenti keras** bila keduanya menunjuk baris berbeda, supaya anggapan
    "keduanya sama-sama tertinggi" tidak patah diam-diam.
11. Sebagai **insinyur pemuatan**, saya ingin polis dikelompokkan menurut **nomor polis**, supaya
    "versi terakhir" punya cakupan yang jelas dan bukan sekadar nomor urut yang mengambang.
12. Sebagai **work owner**, saya ingin **nomor versi disimpan apa adanya** — polis yang sudah sampai
    versi ketujuh tercatat sebagai ketujuh, bukan dinomori ulang jadi pertama — supaya riwayat tidak
    ditulis ulang dan fase 2 tidak perlu membongkarnya.

### Meratakan pohon menjadi baris

13. Sebagai **insinyur pemuatan**, saya ingin tiap wadah berulang di dokumen sumber menjadi
    **baris-baris** pada tabel flat pasangannya, supaya struktur bersarang berubah menjadi relasional
    tanpa kehilangan isi.
14. Sebagai **insinyur pemuatan**, saya ingin pohon ditelusuri sampai **delapan tingkat**, supaya
    dokumen terdalam tetap terbaca utuh.
15. Sebagai **insinyur pemuatan**, saya ingin tiap baris anak membawa penunjuk ke **baris induknya**,
    supaya hubungan asuh di dokumen sumber tetap terbaca di tabel.
16. Sebagai **insinyur pemuatan**, saya ingin dua belas tabel yang induknya bisa lebih dari satu juga
    membawa **nama tabel induknya**, supaya baris anak tidak menjadi yatim yang tidak dapat
    ditelusuri.
17. Sebagai **tim rekonsiliasi**, saya ingin pemuatan **berhenti keras** bila sebuah baris berinduk
    ganda tidak dapat menetapkan tabel induknya, supaya kesalahan yang tidak dijaga basis data tidak
    lolos diam-diam.
18. Sebagai **insinyur pemuatan**, saya ingin urutan baris di dalam satu wadah **dipertahankan** dan
    dicatat sebagai nomor urut, supaya susunan yang terekam di sistem lama tetap terbaca.
19. Sebagai **insinyur pemuatan**, saya ingin tabel yang berperan sebagai **wadah murni** tetap
    memperoleh barisnya meski tidak punya satu pun medan terisi, supaya anak-anaknya punya tempat
    menggantung.
20. Sebagai **insinyur pemuatan**, saya ingin tabel yang dapat dicapai lewat lebih dari satu jalur
    mencatat **jalur sumbernya**, supaya asal sebuah baris tetap dapat ditelusuri.

### Nilai, uang, dan mata uang

21. Sebagai **tim rekonsiliasi**, saya ingin nilai masuk **apa adanya tanpa dibulatkan**, supaya
    presisi dan urutan operasi sistem lama tidak berubah.
22. Sebagai **insinyur pemuatan**, saya ingin titik diperlakukan sebagai **pemisah desimal** dan
    ketiadaan pemisah ribuan diperlakukan sebagai kaidah, supaya angka tidak salah baca.
23. Sebagai **insinyur pemuatan**, saya ingin nilai berkoma desimal ditangani sebagai
    **pengecualian yang dikenali**, bukan sebagai bentuk baku, supaya kasus langka tidak diam-diam
    menjadi nol.
24. Sebagai **insinyur pemuatan**, saya ingin parser dokumen penawaran **terpisah** dari parser tabel
    lookup, supaya sumber yang memakai koma desimal tidak mencemari sumber yang memakai titik.
25. Sebagai **work owner**, saya ingin tiap baris pembawa uang juga membawa **kode mata uangnya**,
    supaya nilai moneter tidak pernah berdiri tanpa satuan.
26. Sebagai **work owner**, saya ingin baris lama yang datang **tanpa mata uang** tersimpan sebagai
    keadaan tidak-diketahui yang eksplisit, bukan diisi mata uang tebakan, supaya kekurangannya
    terhitung dan bersuara.
27. Sebagai **tim rekonsiliasi**, saya ingin jumlah baris yang masuk dengan mata uang tidak diketahui
    **dilaporkan**, supaya besarnya masalah terukur, bukan terasa.
28. Sebagai **insinyur pemuatan**, saya ingin aturan berhenti-keras saat menulis mata uang tidak
    diketahui **tidak berlaku** bagi baris lama yang dimigrasikan, supaya pemuatan tidak berhenti
    oleh keadaan yang justru ingin direkam.
29. Sebagai **insinyur pemuatan**, saya ingin satu baris cukup punya **satu kolom mata uang**, supaya
    baris tidak perlu dipecah.

### Kelompok lini bisnis

30. Sebagai **insinyur pemuatan**, saya ingin kelompok lini bisnis ditentukan oleh **aturan berjenjang
    yang berhenti di kecocokan pertama**, supaya satu penawaran tidak masuk dua kelompok.
31. Sebagai **insinyur pemuatan**, saya ingin empat langkah pertama aturan itu membaca **bentuk
    dokumen**, supaya penentuannya tidak bergantung nilai yang bisa kosong.
32. Sebagai **insinyur pemuatan**, saya ingin langkah terakhir — satu-satunya yang membaca nilai —
    punya perlakuan tegas bila nilainya **muncul lebih dari sekali** dalam satu dokumen, supaya
    hasilnya tidak bergantung pada urutan pembacaan.
33. Sebagai **tim rekonsiliasi**, saya ingin pemuatan **berhenti keras** bila nilai penentu itu
    berisi kode yang tidak dikenal, supaya kelompok lini bisnis tidak ditebak.

### Apa yang tidak ikut dimuat

34. Sebagai **work owner**, saya ingin cabang salinan-kerja **tidak ikut dimuat**, supaya potret
    sementara tidak tersimpan seolah-olah nilai sebenarnya.
35. Sebagai **work owner**, saya ingin medan yang hanya melayani layar **tidak ikut dimuat**, supaya
    tabel tidak menampung parameter tampilan.
36. Sebagai **work owner**, saya ingin metadata ekspor Pega **tidak ikut dimuat**, supaya jejak alat
    ekspor tidak menjadi data bisnis.
37. Sebagai **work owner**, saya ingin ketiga tabel penjumlahan **tidak dibuat**, karena nilainya
    informasi saja dan tidak dipakai untuk pembayaran.
38. Sebagai **tim rekonsiliasi**, saya ingin lima nilai total pada tabel akar **disimpan sebagaimana
    terekam dan tidak dihitung ulang**, supaya tidak muncul selisih yang tidak berasal dari data.

### Tautan antar versi

39. Sebagai **insinyur pemuatan**, saya ingin penunjuk ke polis yang di-renew **hanya terisi pada
    renewal**, dan kosong pada new business maupun endorsement, supaya kolom itu tidak dipaksa
    memikul arti yang bukan miliknya.
40. Sebagai **insinyur pemuatan**, saya ingin generasi endorsement tertaut lewat **nomor polis yang
    sama dengan nomor versi yang menaik**, supaya rantainya terbaca tanpa kolom tambahan.
41. Sebagai **work owner**, saya ingin penunjuk renewal **kosong di seluruh baris fase 1**, karena
    generasi sebelumnya memang belum ada untuk ditunjuk.

### Identitas baris yang bersifat sementara

42. Sebagai **insinyur pemuatan**, saya ingin tiap baris wadah berulang punya identitas sendiri,
    supaya baris dapat dirujuk tanpa bergantung pada posisinya.
43. Sebagai **work owner**, saya ingin identitas baris hasil fase 1 **tidak pernah keluar** dari tabel
    flat — tidak diekspor, tidak dikirim ke tim lain, tidak ditampilkan sebagai identitas, tidak
    dijadikan kunci integrasi — supaya pengosongan tabel di fase 2 tetap menjadi tindakan internal.
44. Sebagai **work owner**, saya ingin ada **satu titik beku yang diumumkan** setelah fase 2 lolos
    rekonsiliasi, sesudahnya tidak ada pemuatan ulang lagi, supaya kelonggaran masa uji tidak
    menjalar ke masa pakai.

### Bentuk masukan

45. Sebagai **insinyur pemuatan**, saya ingin program membaca bentuk **JSON** yang dipakai produksi,
    bukan hanya bentuk ekspor XML yang dipakai contoh, supaya yang diuji sama dengan yang dijalankan.
46. Sebagai **insinyur pemuatan**, saya ingin satu perbedaan bentuk yang sudah diketahui — daftar
    bertakik kode mata uang pada XML versus larik pada JSON — **ditangani secara tegas**, supaya jalur
    sumber yang diturunkan dari XML tetap sah atas masukan JSON.

### Keamanan dan kepatuhan

47. Sebagai **work owner**, saya ingin nama orang, alamat surel, tanggal lahir, plat nomor, dan alamat
    tertanggung **tidak pernah muncul di log, laporan, pesan galat, maupun berkas uji**, meski
    kolomnya tetap dimuat ke tabel.
48. Sebagai **work owner**, saya ingin laporan pemuatan hanya memuat **hitungan, nama tabel, dan nama
    kolom**, supaya laporan dapat dibagikan tanpa menyaring ulang.

### Rekonsiliasi

49. Sebagai **tim rekonsiliasi**, saya ingin jumlah baris per tabel dapat diadu dengan jumlah simpul
    di dokumen sumber, supaya kehilangan baris ketahuan sebagai angka, bukan sebagai firasat.
50. Sebagai **tim rekonsiliasi**, saya ingin rekonsiliasi nol-selisih dijalankan **setelah fase 2**
    atas populasi lengkap, supaya yang diuji adalah keadaan sasaran, bukan keadaan sementara.

---

## Implementation Decisions

### Modul yang dibangun

| Modul | Tanggung jawab |
| --- | --- |
| **`loader`** | Mengubah **satu dokumen penawaran** menjadi himpunan baris. ⛔ **Murni** — tidak menyentuh basis data, tidak menyentuh jam, tidak menyentuh berkas |
| **`repository`** (diperluas) | Menulis himpunan baris ke Oracle dalam satu transaksi; menyediakan pemilihan versi polis terakhir |
| **`models`** (diperluas) | Bentuk baris per tabel, memakai tipe uang dan rasio yang sudah ada |

Arah ketergantungan mengikuti aturan yang berlaku: **handlers → services → repository**. Program
pemuatan berdiri sebagai perkakas di atas `repository`, dan **tidak** dipanggil dari `handlers`.

### Seam pengujian — **satu seam baru**, disetujui work owner 25 September 2026

| Seam | Bentuk | Status |
| --- | --- | --- |
| **`loader.Flatten`** | satu dokumen penawaran → himpunan baris untuk 78 tabel | ⭐ **baru** |
| **seam `repository`** | himpunan baris → Oracle | **sudah dicadangkan** spec 03 §5 butir 5, dulu tertunda karena bentuk tabel flat belum ada |

⛔ **Tidak ada seam ketiga.** Tiga seam runtime yang sudah ada (`rules.Eval`, `acceptance.Next`,
`premium.Calculate`) **tidak disentuh** — seluruhnya soal perhitungan, bukan pemindahan data.

**Mengapa `Flatten` murni:** kebenaran pemetaan pohon ke 78 tabel adalah bagian yang paling mudah
salah dan paling sunyi kesalahannya, terutama pada dua belas tabel berinduk ganda yang **tidak dijaga
constraint apa pun**. Memurnikannya membuat seluruh pemetaan itu dapat diuji **tanpa basis data dan
tanpa tiruan**, dengan kasus uji berupa dokumen masuk dan baris keluar.

### Bentuk keluaran seam `Flatten`

Bentuk berikut **mengunci keputusan** dan karena itu dicantumkan; ia bukan contoh kode yang berjalan.

```
FlattenResult
  Rows        peta: nama tabel -> deretan baris, urutan dipertahankan
  Diagnostics hitungan: baris per tabel, mata uang tidak diketahui, medan dibuang
```

Tiap baris membawa kolom sistem yang wajib: identitas surrogate, kunci kerja penawaran, penunjuk
induk, **nama tabel induk** (wajib bagi dua belas tabel berinduk ganda), nomor urut baris, jalur
sumber (hanya bagi tabel berjalur ganda), kelompok lini bisnis, dan identitas baris.

### Keputusan yang mengikat

1. **Dua fase, dengan pengosongan di fase 2** (K-070, K-073). Fase 2 bukan penambahan — ia
   **mengosongkan lalu memuat penuh**. Program pemuatan karena itu wajib mendukung pengosongan
   sebagai operasi yang disengaja dan tercatat, bukan sebagai efek samping.
2. **Identitas baris fase 1 bersifat sementara** (K-073). Kontrak ini ditegakkan dengan tidak
   memaparkannya lewat antarmuka apa pun di luar `repository`.
3. **Pemilihan versi terakhir** (K-071, K-072): kelompokkan menurut **nomor polis**, urutkan menurut
   **tanggal masuk menurun**, periksa silang dengan **nomor versi tertinggi**. Ketidaksepakatan
   keduanya → **berhenti keras**. Asal aturan urut: `GetProdKeOldData_SQL` memakai
   `ORDER BY TGL_INPUT DESC`.
   ⛔ **Nomor versi tidak boleh menjadi penentu urutan** — `[terverifikasi]` `DDL\JSON_POLIS.txt`
   mendefinisikannya `VARCHAR2(5)`, sehingga pengurutannya adalah perbandingan teks dan `'9'` melebihi
   `'10'`. Pola yang sama dengan `FACINOFFER.RATE VARCHAR2(100)`.
4. **Nomor versi berubah tipe saat dimuat** (K-071 J-3): sumber teks, sasaran numerik. Nilai yang
   bukan angka → **berhenti keras**, bukan kosong dan bukan nol.
5. **Penunjuk renewal hanya untuk renewal** (K-071 J-5, mempersempit K-068). New business dan
   endorsement membiarkannya kosong; pada fase 1 seluruh baris kosong. Kendala unik atas kolom itu
   tetap sah karena Oracle mengizinkan banyak nilai kosong pada indeks unik.
6. **Generasi endorsement tertaut lewat nomor polis + nomor versi** (K-072 J-11).
   ⚠️ **Aturan ini `[dugaan]` dan tidak dapat diuji dari korpus.** `[terverifikasi]` 115 contoh memuat
   **113 nomor polis unik**; dua yang berulang adalah dua pasang duplikat byte-identik yang sudah
   dikenal (V-45), bukan generasi berbeda. ⛔ **Nol pasangan generasi sejati** tersedia untuk diuji.
7. **Mata uang**: satu kolom per baris sudah memadai — `[terverifikasi]` dari 13.666 baris, **432**
   membawa kode mata uang dan **seluruhnya tepat satu**, nol yang lebih dari satu. **24 kolom** mata
   uang berkeadaan wajib-terisi dengan nilai bawaan *tidak diketahui* (K-069), sehingga baris lama
   tanpa mata uang **terekam dan terhitung**, bukan ditolak. Aturan berhenti-keras ADR-0006 berlaku
   bagi tulisan **baru aplikasi**, bukan bagi baris lama yang dimigrasikan.
8. **Uang dan rasio memakai tipe yang sudah ada** (ADR-0004): nilai moneter selalu berpasangan dengan
   mata uangnya; rasio selalu berpasangan dengan skalanya. ⛔ Tidak ada bilangan pecahan biner di
   jalur uang, baik di Go maupun di Oracle.
9. **Agregat lintas mata uang memakai tipe ketiga** (ADR-0007) — bukan tipe uang biasa. Berlaku bagi
   lima nilai total pada tabel akar, yang `[terverifikasi]` **disimpan sebagaimana terekam dan tidak
   dihitung ulang** (K-067).
10. **Presisi tidak diubah** (ADR-0005, ADR-0001): nilai masuk apa adanya, program pemuatan **tidak
    membulatkan**. Format angka `[terverifikasi]` atas 115 contoh: **185.389** bulat · **58.124**
    titik desimal · **71** koma desimal · **nol** pemisah ribuan dalam bentuk apa pun.
11. **Parser dokumen penawaran terpisah dari parser tabel lookup.** `[terverifikasi]`
    `FACINOFFER.RATE` dan `TABLEOFLIMIT.PCTLIMIT` memakai **koma** desimal, sedangkan dokumen
    penawaran memakai **titik**. Menyatukan keduanya akan salah membaca salah satunya.
12. **Kelompok lini bisnis** mengikuti aturan berjenjang V-16, berhenti di kecocokan pertama; empat
    langkah pertama membaca bentuk, langkah kelima membaca nilai. `[terverifikasi]` nilai penentu itu
    muncul **128 kali pada 115 berkas** dengan **18 nilai berbeda** — jadi sebagian dokumen memuatnya
    **lebih dari sekali**, dan aturan "mana yang mengikat" wajib ditetapkan, bukan diserahkan pada
    urutan pembacaan. Kode di luar delapan belas nilai itu → **berhenti keras**.
13. **Yang dikecualikan**, seluruhnya keputusan yang sudah ada: cabang salinan-kerja (V-19,
    `[terverifikasi]` **326 kemunculan** dibuang pada satu sapuan sembilan medan) · cabang retro tanpa
    daftar (V-25) · parameter layar (V-34) · tiga tabel penjumlahan (V-28, K-067) · medan bersufiks
    lama (V-28b, V-45) · metadata ekspor 23 tag (K-042/K-043).
    ⛔ **Cabang salinan-kerja bukan pengganti versi sebelumnya** — V-43a `[terverifikasi]` menyatakan
    ia potret **sesudah** baris baru dibentuk.
14. **Masukan produksi adalah JSON** (K-066). `[terverifikasi]` satu kasus yang tersedia dalam kedua
    bentuk cocok **621 dari 621** jalur daun setelah ruas kode mata uang dinormalkan; satu-satunya
    jalur khas XML adalah metadata Pega. ⛔ Perbedaan **bentuk** tetap wajib ditangani: pada XML kode
    mata uang menjadi **nama elemen**, pada JSON ia unsur **larik** dengan kodenya di dalam medan
    nama. Pembawa kode mata uang tetap medan yang sama di kedua bentuk, sehingga kolom mata uang
    tetap sah.
15. **Urutan muat tabel** diturunkan dari daftar relasi: tabel kerja → tabel akar (berbagi kunci
    utama, satu lawan satu) → anak langsung → menurun mengikuti pohon sampai delapan tingkat. Dua
    belas tabel berinduk ganda mengisi penunjuk induk **dan** nama tabel induk dalam transaksi yang
    sama.
16. **Tiga tabel adalah wadah murni** — `[terverifikasi]` nol medan skalar terisi. Barisnya tetap
    dibuat (V-27) berisi kolom sistem saja, karena anak-anaknya menggantung padanya.
17. **Jalur tulis produksi tetap wajib ada di sistem baru**, dengan alasan yang baru diketahui:
    `[keterangan work owner]` keluarannya **dipakai tim lain**. Ini menguatkan lingkup yang sudah ada,
    bukan menambahnya.
18. **Pasangan nilai-sesudah/selisih diport apa adanya.** `[terverifikasi]` dari 80 berkas DDL,
    **5** memuat pasangan itu: **9** basis bernilai-sesudah, **15** basis berselisih, **9**
    berpasangan lengkap. ⛔ Enam basis berselisih **tanpa** pasangan nilai-sesudah — `[keterangan work
    owner]` keenamnya menandai **penambahan lokasi**, yaitu keadaan yang memang **tidak punya nilai
    sebelumnya**. Karena itu keenamnya **bukan pasangan yang hilang**; ⛔ **jangan dibuatkan kolom
    nilai-sesudah baru** — itu perbaikan diam-diam yang dilarang.

### Yang wajib berhenti keras

Berhenti keras dipilih di tiap titik berikut karena aplikasi yang berjalan dan salah diam-diam jauh
lebih berbahaya daripada aplikasi yang berhenti:

| Keadaan | Alasan |
| --- | --- |
| Tanggal masuk terbaru dan nomor versi tertinggi **menunjuk baris berbeda** | anggapan pemilihan versi patah |
| Nomor versi **bukan angka** | konversi tipe tidak dapat dipercaya |
| Baris berinduk ganda **tidak dapat menetapkan tabel induknya** | tidak ada constraint yang akan mengeluh |
| Kode penentu lini bisnis **di luar delapan belas nilai yang terbaca** | arti kode tidak boleh ditebak |
| Kedalaman pohon **melebihi delapan tingkat** | melampaui yang pernah terukur |
| Kode enumerasi yang **artinya belum dijawab** tercapai | `BusinessCode` dan `BusinessOldId` (98 dan 87 kode) |

---

## Testing Decisions

### Apa yang membuat sebuah test baik di sini

Test menguji **perilaku yang terlihat dari luar seam**, bukan cara kerja di dalamnya. Untuk
`loader.Flatten` itu berarti: masukan satu dokumen penawaran, keluaran himpunan baris. Test **tidak
boleh** memanggil fungsi bagian dalam, tidak boleh memeriksa urutan pemanggilan, dan tidak boleh
menyentuh basis data.

**Prior art:** ketiga seam yang sudah ada — registry predikat, satu transisi tangga akseptasi, dan
satu pintu masuk perhitungan premi. Ketiganya table-driven, bentuknya harus diikuti, dan spec 03
menegaskan ketiganya **menjadi** prior art bagi modul berikutnya.

### Yang diuji di seam `loader.Flatten`

Table-driven, tiap kasus berupa dokumen masuk dan baris keluar yang diharapkan. Kasus wajib:

1. **Satu penawaran sederhana** — memastikan bentuk dasar baris dan kolom sistem.
2. **Wadah berulang bersarang sampai delapan tingkat** — memastikan penelusuran tidak berhenti dini.
3. **Dua belas tabel berinduk ganda**, tiap tabel setidaknya sekali, termasuk tabel yang penunjuk
   induknya dapat menunjuk **lima** tabel berbeda. Memeriksa penunjuk induk **dan** nama tabel
   induknya.
4. **Tabel wadah murni** — memastikan barisnya tetap dibuat meski tanpa medan terisi.
5. **Tiap kelompok lini bisnis**, satu kasus per kelompok, termasuk kasus yang dokumennya memuat
   nilai penentu **lebih dari sekali**.
6. **Baris tanpa mata uang** — memastikan nilai bawaan *tidak diketahui* terpasang dan terhitung di
   laporan, bukan menggagalkan pemuatan.
7. **Nilai berkoma desimal** — memastikan pengecualian itu terbaca benar, dan nilai bertitik desimal
   tidak ikut berubah.
8. **Cabang yang dikecualikan** — memastikan cabang salinan-kerja, parameter layar, dan metadata
   ekspor **tidak** menghasilkan baris.
9. **Tiga tabel penjumlahan** — memastikan tidak ada baris yang dibuat untuk ketiganya.
10. **Tiap keadaan berhenti keras** pada tabel di atas, satu kasus masing-masing, memastikan ia
    **memang** berhenti dan tidak mengembalikan nilai bawaan.
11. **Kedua bentuk masukan** atas kasus yang sama — memastikan JSON dan XML menghasilkan himpunan
    baris yang **sama**, termasuk pada daftar bertakik kode mata uang.

### Yang diuji di seam `repository`

1. **Satu penawaran ditulis dalam satu transaksi**; kegagalan di tengah tidak meninggalkan baris.
2. **Pemilihan versi terakhir** — kelompok dengan beberapa versi, memastikan yang terpilih adalah
   tanggal masuk terbaru, dan memastikan **berhenti keras** saat nomor versi tidak sepakat.
3. **Konversi tipe nomor versi**, termasuk nilai bukan angka yang wajib berhenti keras.
4. **Dapat diulang** — menjalankan pemuatan dua kali menghasilkan isi tabel yang sama.
5. **Pengosongan fase 2** — memastikan ia mengosongkan seluruh 78 tabel dan bukan sebagian.

### Yang tidak dapat diuji sekarang, dan alasannya dicatat

⛔ **Tautan antar generasi endorsement tidak dapat diuji dari korpus.** `[terverifikasi]` tidak ada
satu pun pasangan generasi sejati di antara 115 contoh. Test untuknya wajib memakai dokumen yang
**disusun sendiri**, dan statusnya `[dugaan]` sampai ada contoh nyata berpasangan. Ini dicatat sebagai
kekurangan yang diketahui, bukan sebagai cakupan yang lengkap.

---

## Out of Scope

- **Jalur tulis produksi** ke tabel produksi fac in dan fac out. Keduanya sudah flat dan sudah ada;
  jalurnya punya tiketnya sendiri.
- **Cakupan Treaty In.** Korpusnya tidak dibaca (K-005). Dua tabel akar `[terverifikasi]` dipakai
  bersama Treaty In (K-064), tetapi kolomnya dibangun **hanya dari sisi Facultative Inward**; irisan
  kolom dan tabrakan ruang kunci baru akan muncul saat integrasi — risiko yang diterima sadar (K-065).
- **Kolom penanda lini bisnis facultative-versus-treaty.** ⛔ Work owner menetapkan itu **di luar
  proyek ini** (K-069); kolomnya sudah dicabut dari rancangan.
- **Penjadwalan fase 2.** Spec ini menetapkan **bentuknya**; kapan ia dijalankan adalah keputusan
  pelaksanaan. Syarat yang mengikat tercatat di K-072.
- **Perubahan skema Oracle.** Skema tabel produksi dan lookup tidak berubah; tabel flat sudah
  ber-DDL draf dan tidak diputuskan ulang di sini.
- **Antarmuka pengguna.** Program pemuatan adalah perkakas, bukan layar.
- **Perbaikan perilaku sistem lama.** Segala kejanggalan yang ditemui **diport apa adanya**;
  perbaikan dipisahkan dari migrasi dan menjadi keputusan bisnis.

---

## Further Notes

### Dua butir yang masih terbuka, keduanya tidak memblokir

| # | Butir | Dampak |
| :-: | --- | --- |
| **J-16** | Penyaring "endorsement sudah jadi" untuk fase 1. ⛔ Status konversi **bukan** jawabannya — `[keterangan work owner]` kolom itu menandai penyerahan ke tim lain, bukan selesainya endorsement | tidak menentukan hasil akhir karena fase 2 memuat ulang; **tetapi** menguji atas generasi yang salah menghasilkan pengujian yang menyesatkan |
| **J-17** | Arti dua kolom bernama "master" pada tabel akar — `[terverifikasi]` terisi **1 dari 115** dan berbeda dari nomor polis | tidak memblokir; artinya `belum terverifikasi` |

### Risiko yang diterima sadar

⚠️ **Selama fase 1, sistem diuji atas data yang tidak lengkap.** Perilaku yang bergantung pada
riwayat — jalur banding dan penomoran endorsement berikutnya — akan berlaku berbeda dibanding saat
sistem dipakai sungguhan. ⛔ **Lolos uji di fase 1 bukan jaminan.** Peredamnya satu: jalankan fase 2
**sedini mungkin**, bukan selambat mungkin.

⚠️ **Syarat yang paling mudah dilanggar tanpa sadar** adalah bahwa sumber lama wajib **tetap dapat
dibaca** sampai fase 2 selesai. Mematikan sistem lama terasa seperti akhir pekerjaan, padahal fase 2
justru membutuhkannya hidup. ⛔ Bila dilanggar, generasi lama tidak dapat dimuat lewat jalan mana pun.

### Angka yang belum terekonsiliasi

⚠️ Dua cara menghitung simpul calon baris memberi **22.139** dan **20.937**, berselisih **1.202**,
padahal hanya ada satu ejaan. Angka yang dipakai untuk menentukan ukuran pekerjaan —
**20.937 + 1.685 = 22.622** — adalah yang konsisten secara internal. ⛔ **Memadai untuk menentukan
ukuran; jangan dipakai sebagai angka rekonsiliasi.**

### Utang dokumen yang menyentuh spec ini

⚠️ **Bunyi V-38 dan V-40 masih menyebut pencarian nomor versi tertinggi** sebagai cara menaut versi.
Keduanya digantikan: penunjuk renewal untuk renewal (K-068, dipersempit K-071), nomor polis + nomor
versi untuk endorsement (K-072). Dua keputusan sah yang akan terus terbaca bertentangan sampai
bunyinya diperbaiki.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
