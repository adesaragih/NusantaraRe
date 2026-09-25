> Modul  : Komite Claim Non Prop · Tahap 2 · 2026-09-21
> Peran  : penulis spesifikasi
> Ronde  : lapis persyaratan, 2026-09-21 — penambahan di atas isi naratif 2026-09-21.
>          Nol bagian naratif ditulis ulang; yang disunting hanya *Testing Decisions*,
>          tabel pagar pada *Out of Scope*, dan cerita 42 — ketiganya atas perintah yang
>          tercatat pada `AUDIT-SPEC-01.md` bagian 7
> Ronde  : penutupan, 2026-09-21 — model data per objek, tabel medan layar, keadaan dua
>          baris ketertelusuran, dan tiga deviasi turunan `K6-4`. Penambahan; nol
>          persyaratan baru
> Ronde  : objek ketujuh, 2026-09-21 — `PERISTIWA_ROSTER` ditambahkan sebagai tempat
>          mendarat `S-052`, yang sebelumnya tidak punya tabel mana pun. Nol persyaratan
>          baru; `S-052` sudah ada dan bunyinya tidak diubah
> Masukan: SPEC-KOMITE-01.md (isi naratif, 2026-09-21) · AUDIT-SPEC-01.md (2026-09-21) ·
>          KETETAPAN.md · CONTEXT.md · INVENTARIS-BUKTI.md · REGISTER-PAGAR.md ·
>          REGISTER-DEVIASI.md · GRILL-05/04-PARITAS.md · PENGETAHUAN.md §2.3 dan §13 ·
>          docs/adr/0030–0032 · SPEC-MODEL-DATA.md §16 (sisi Klaim) ·
>          ddl-usulan/01_KLAIM.sql dan 06_ADJUSTMENT.sql · PENGETAHUAN.md §7.1, §7.2, §10.1
> Status : TERBUKA
> Sifat  : HIDUP

# SPEC-KOMITE-01

Kosakata berkas ini diambil dari `CONTEXT.md`. Ketetapan dirujuk dengan nomor dari
`KETETAPAN.md`; bunyinya tidak ditulis ulang di sini. Berkas ini **tidak memutuskan apa pun
tentang uang** — nilai, kurs, dan pembayaran milik modul lain dan aliran yang berpagar.

---

## Cakupan dan cacah

Apa yang ditulis penuh, apa yang ditulis sebagian, dan apa yang sengaja tidak ditulis.

| Aliran | Isinya | Keadaan di spesifikasi ini |
|---|---|---|
| A-1a | Mekanika memutus, urutan giliran, dan wewenang | **DITULIS PENUH** |
| A-1b | Pembentukan sirkulasi, pemilihan jenjang, roster, tabel seleksi | **DITULIS PENUH** |
| A-2 | Layar persetujuan | **DITULIS SEBAGIAN** — aturannya ditulis; tabel medan per kendali belum |
| A-3 | Penomoran akseptasi | **DITULIS, AMBANGNYA BERPAGAR** — mekanika penuh; nilai ambang periode buku di balik PAGAR-01 |
| A-4 | Kiriman pembayaran ke Kasir | **BEKU** — `PG-04`; PAGAR-05, PAGAR-06 |
| A-5 | Akseptasi, rincian layer, pembalikan | **BEKU** — `PG-05`; PAGAR-02, PAGAR-03, PAGAR-04 |
| A-6 | Surat dan dokumen | **ISINYA BEKU** — `PG-06`; PAGAR-07, PAGAR-08. Kapan terbit dan siapa penerimanya ditulis |

| Yang dicacah | Cacah |
|---|---:|
| Persyaratan `S-xxx` | **70** |
| — aliran A-1a | 30 |
| — aliran A-1b | 27 |
| — aliran A-2 | 6 |
| — aliran A-3 | 7 |
| Cerita yang berlaku | 43 dari 44 nomor; cerita 42 dicabut |
| Pagar | **8**, atas empat aliran |
| Baris ketertelusuran | **63** |
| Deviasi bernomor yang dirujuk | **25** — 24 berlaku, 1 gugur karena dilebur |
| Objek skema baru | **7**, berisi **64 kolom** |
| Perubahan pada tabel milik sisi Klaim | 1 kolom baru, 1 constraint baru |
| Medan layar | **7** dapat disunting, **11 kelompok** baca-saja, **5** dibuang |
| Invarian | 11, tiga di antaranya pemeriksaan skema |

Berkas ini **tidak memutuskan apa pun tentang uang**, dan tidak memuat satu pun persyaratan
untuk A-4, A-5, maupun isi A-6.

---

## Seam pengujian — mohon diperiksa sebelum dibangun

**Satu seam.** Batas aplikasi modul Komite: dua perintah dan satu kueri, dipanggil langsung
dalam proses, terhadap skema `KLAIMNP` yang nyata.

```
BentukSirkulasi(idKlaim, idUsulan, jenisSirkulasi, maksud, pelaku, kunciIdempotensi)
CatatKeputusan(idSirkulasi, hasil, komentar, penandaUsulan?, pelaku, kunciIdempotensi)
BacaSirkulasi(idSirkulasi | idKlaim)
```

Alasan memilih titik ini dan bukan yang lain:

- **Lebih tinggi dari HTTP.** Kontrak HTTP adalah adapter tipis di atas ketiganya. Menguji
  lewat HTTP menambah serialisasi dan routing ke dalam setiap uji tanpa menambah satu pun
  perilaku yang diuji.
- **Lebih tinggi dari repositori.** `ADR-0032` menetapkan sirkulasi, jenjang, dan keputusan
  sebagai satu agregat yang berubah bersama; `ADR-0031` menetapkan pembentukan dan akibatnya
  sebagai satu transaksi. Seam di bawah batas transaksi tidak dapat menguji keduanya.
- **Basis data nyata, bukan ganda.** Tiga keputusan ditegakkan oleh basis data, bukan oleh
  kode: idempotensi nomor akseptasi lewat constraint (`D-4`), keunikan derajat per sirkulasi,
  dan keunikan nomor urut sirkulasi per usulan (`F-11`). Ganda dalam-memori akan meloloskan
  ketiganya.
- **Layar tidak menjadi seam.** `F-1`, `F-6`, `F-9` menunjukkan seluruh panel nilai baca-saja;
  satu-satunya masukan layar adalah keputusan, komentar, dan penanda usulan. Menguji layar
  berarti menguji tiga medan, dan ketiganya sudah diuji di seam ini.

**Bila ini tidak sesuai harapan Anda, katakan sekarang** — seluruh bagian "Testing Decisions"
bergantung padanya.

### Seam — diratifikasi 2026-09-21

Pertanyaan di atas sudah dijawab, dan jawabannya mengubah tiga hal. Bagian di atas dibiarkan
apa adanya agar terbaca apa yang ditanyakan; yang berlaku adalah yang di bawah ini.

**Lima operasi, bukan tiga:**

```
BentukSirkulasi   CatatKeputusan   BacaSirkulasi   KelolaRoster   KelolaTabelSeleksi
```

Dua yang terakhir masuk karena menopang cerita 35, 36, 37, 38, 39, dan 40, dan karena `K5-5`
menjadikan tabel seleksi **data** — data yang dikelola adalah perilaku, bukan pengaturan
mati. Keduanya sudah muncul pada *Kontrak API*; yang kurang hanya pengakuannya sebagai seam.

**"Adapter tipis" berlaku untuk HTTP, tidak untuk layar.** Spesifikasi ini menuntut tiap
kunci ditegakkan dua kali, dan penegakan pertama hidup di luar seam perintah. Layar karena
itu mendapat seam sendiri yang **sempit**: ujinya memeriksa bahwa penegakan pertama **ada**,
bukan bahwa aturannya benar. Kebenaran aturannya diuji di seam perintah, satu kali.

**`I-2`, `I-8`, dan `I-10` adalah pemeriksaan skema, bukan pemeriksaan perilaku.** Ketiganya
diuji sebagai uji DDL terhadap katalog skema — lihat `S-064`, `S-067`, `S-068`. Bentuk ini
dipilih karena constraint yang membuat keadaan terlarang **tidak dapat ditulis** lebih kuat
daripada uji yang memanggil operasi, dan sejalan `D-4` yang memilih constraint di atas
baca-lalu-tulis. Menandainya penting agar bentuk ujinya tidak dikira uji perilaku yang
tertinggal.

---

## Problem Statement

Komite klaim non-proporsional hari ini bekerja di dalam Pega, dan cara kerjanya menyimpan
enam masalah yang dirasakan langsung oleh orangnya:

**Pemutus tidak tahu apakah gilirannya benar-benar giliran dia.** Pemeriksaan wewenang di
sistem lama hanya memasang pesan lalu melanjutkan; siapa pun yang dapat membuka layar dapat
mengirim keputusan, dan keputusan itu tersimpan. Pengenal pengguna dicocokkan sebagai
potongan teks, sehingga seseorang yang namanya kebetulan terkandung di dalam nama pemegang
jenjang ikut lolos.

**Orang menemukan namanya pada keputusan yang tidak pernah ia ambil.** Ketika satu jenjang
menolak, sistem lama menuliskan penolakan atas nama seluruh jenjang di bawahnya — lengkap
dengan komentar dan jam milik si penolak. Riwayat menjadi bukti yang salah tentang orang.

**Sirkulasi kadang lahir tanpa seorang pun di dalamnya.** Bila bagian treaty atau nilai
usulan jatuh di kombinasi yang tidak tercakup aturan, penyaring roster dipanggil tanpa
ambang dan mengembalikan nol baris. Berkasnya tetap terbentuk, lalu diam.

**Pembentukan yang gagal tidak terlihat gagal.** Ketika daftar penyebaran kosong, prosesnya
keluar diam-diam sambil menyimpan penanda galat yang tidak dibaca siapa pun. Pengguna
menekan tombol dan tidak terjadi apa-apa. Pada jalur galat yang lain, pembuatan dilewati
tetapi langkah sesudahnya tetap berjalan, sehingga klaim menerima nomor sirkulasi milik
berkas lain, menyimpannya, dan mengirim surat atasnya.

**Klaim terlihat sudah ditutup padahal komite belum memutus apa pun.** Penanda "ditutup" dan
"ditolak" ditulis ke klaim pada saat pengajuan, di langkah yang sama yang membuat berkas
sirkulasi — dan tetap tertinggal di sana ketika sirkulasinya gagal lahir.

**Siapa yang harus memutus ditentukan oleh angka di dalam kode.** Kelas kewenangan dihitung
dari empat tetapan yang ditulis di dalam langkah, satu cabangnya tidak pernah dapat benar,
dan satu substitusi orang ke orang dipaku di dalam aturan. Ketika seseorang pindah jabatan,
tidak ada tempat untuk mengubahnya selain kode.

---

## Solution

Modul Komite dibangun ulang sebagai modul di dalam konteks Klaim (`ADR-0032`), dengan
**sirkulasi** sebagai agregat sendiri: satu berkas yang mengedarkan satu **usulan** kepada
sejumlah **jenjang** untuk diputuskan.

Dari sudut pandang orang yang memakainya:

- **Giliran itu pasti.** Satu-satunya yang boleh memutus adalah **pemegang** **jenjang
  aktif**, yaitu jenjang berderajat terendah yang belum memutuskan. Permintaan dari orang
  lain **ditolak**, bukan ditandai, dan penolakannya terjadi di server (`K5-1`, keputusan
  beku no. 1).
- **Riwayat hanya memuat tindakan yang benar-benar terjadi.** Jenjang yang tidak sempat
  mendapat giliran berkeadaan `TIDAK_SAMPAI` — bukan "menolak" (keputusan beku no. 2).
  Peran dan jabatan disalin ke keputusan pada saat diambil, sehingga mutasi jabatan tidak
  mengubah masa lalu (keputusan beku no. 4).
- **Siapa yang berwenang adalah data.** **Tabel seleksi** menentukan **kelas kewenangan**
  dari tiga masukan: nilai usulan, **bagian treaty**, dan bersyarat (`K5-5`, `K6-3`). Setiap
  kombinasi punya satu **aturan seleksi**; kombinasi yang tidak tercakup adalah **galat
  konfigurasi** yang menolak pembentukan, bukan sirkulasi kosong. **Roster** dimiliki sistem
  baru dan diubah tanpa menyentuh kode (`D-5`), termasuk **delegasi tetap** (`H-3`).
- **Pembentukan berhasil seluruhnya atau tidak sama sekali.** Sirkulasi, penulisan
  pengenalnya ke klaim, penyimpanan klaim, dan pemberitahuan adalah satu transaksi
  (`K5-2`, `ADR-0031`). Kegagalan terlihat oleh pengguna dan tercatat (`K5-3`). Sirkulasi
  tanpa jenjang ditolak saat dibuat (`H-6`, keputusan beku no. 8).
- **Maksud dan akibat dipisahkan.** "Diajukan untuk ditutup" tercatat sejak pengajuan;
  "ditutup" hanya lahir dari keputusan komite lewat **fungsi akibat** (`K5-7`, `D-3`, `J-3`).
  Klaim yang sirkulasinya gagal lahir membawa maksudnya dan tidak membawa akibatnya.
- **Satu tempat menghitung akibat.** **Hasil** selalu berarti putusan atas usulan, tidak
  pernah nasib klaim. Nasib klaim dihitung dari (jenis sirkulasi, bersyarat, hasil) di satu
  tempat, dan model baca menyajikan tiga nilai tampil dari fungsi yang sama (`D-3`, `J-3`,
  keputusan beku no. 9).

---

## User Stories

**Membentuk sirkulasi**

1. Sebagai analis klaim, saya ingin mengirim satu usulan pembayaran ke komite, agar ia
   diputuskan oleh jenjang yang berwenang atas nilai itu.
2. Sebagai analis klaim, saya ingin mengirim usulan menutup klaim ke komite, agar penutupan
   punya jejak keputusan, bukan hanya jejak perintah.
3. Sebagai analis klaim, saya ingin mengirim usulan menolak klaim ke komite, dengan alasan
   yang sama.
4. Sebagai analis klaim, saya ingin jenis sirkulasi ditentukan oleh tindakan yang saya pilih,
   bukan ditebak dari isi sebuah kolom analisis, agar dua maksud berbeda tidak tertukar.
5. Sebagai analis klaim, saya ingin sistem memilih jenjang mana yang harus memutus, agar saya
   tidak perlu menghafal ambang kewenangan.
6. Sebagai analis klaim, saya ingin pembentukan yang gagal memberi tahu saya alasannya, agar
   saya tidak menunggu berkas yang tidak pernah lahir.
7. Sebagai analis klaim, saya ingin pembentukan yang gagal tidak meninggalkan apa pun pada
   klaim, agar klaim tidak terlihat sedang disirkulasikan padahal tidak.
8. Sebagai analis klaim, saya ingin menekan tombol dua kali tidak melahirkan dua sirkulasi
   atas usulan yang sama.
9. Sebagai analis klaim, saya ingin mengirim ulang usulan yang pernah ditolak komite, agar
   perbaikan dapat diajukan kembali.
10. Sebagai analis klaim, saya ingin setiap pengajuan ulang punya nomor urut tersendiri pada
    usulan itu, agar dua putaran dapat dibedakan tanpa mengandalkan jam.
11. Sebagai analis klaim, saya ingin diberi tahu ketika kombinasi nilai dan bagian treaty
    tidak tercakup aturan seleksi, agar yang salah diperbaiki di data acuan, bukan dibiarkan
    menjadi sirkulasi kosong.

**Memutus**

12. Sebagai pemegang jenjang, saya ingin melihat daftar sirkulasi yang menunggu keputusan
    saya, agar saya tahu apa yang menjadi giliran saya hari ini.
13. Sebagai pemegang jenjang, saya ingin menyetujui usulan dengan komentar, agar alasan saya
    tersimpan bersama keputusannya.
14. Sebagai pemegang jenjang, saya ingin menolak usulan dengan komentar, dengan cara yang sama.
15. Sebagai pemegang jenjang, saya ingin menyetujui usulan secara bersyarat beserta catatan
    syaratnya, agar persetujuan dapat diberikan tanpa menerbitkan kewajiban yang belum layak
    terbit.
16. Sebagai pemegang jenjang pertama, saya ingin menyunting penanda usulan sambil memutus,
    agar koreksi kecil tidak memerlukan putaran baru.
17. Sebagai pemegang jenjang kedua dan seterusnya, saya ingin yakin bahwa yang saya putuskan
    sama persis dengan yang dilihat jenjang sebelum saya.
18. Sebagai pemegang jenjang, saya ingin permintaan saya ditolak bila giliran bukan milik
    saya, agar saya tidak mengira keputusan saya tercatat padahal tidak.
19. Sebagai pemegang jenjang, saya ingin keputusan saya gagal seluruhnya ketika ada gangguan,
    bukan tersimpan separuh.
20. Sebagai pemegang jenjang, saya ingin membedakan "sistem sedang bermasalah" dari "Anda
    tidak berwenang", agar saya tahu harus menunggu atau menghubungi orang.
21. Sebagai pemegang jenjang, saya ingin masukan saya tidak hilang ketika penyimpanan gagal.
22. Sebagai pemegang jenjang, saya ingin percobaan ulang setelah tanggapan hilang tidak
    melahirkan keputusan kedua.
23. Sebagai pemegang jenjang, saya ingin keputusan saya tidak dapat disunting sesudah
    tercatat, agar tidak ada yang mengubah apa yang saya putuskan.
24. Sebagai pemegang jenjang terakhir, saya ingin persetujuan saya menutup sirkulasi dan
    menerbitkan **nomor akseptasi** bila usulannya tidak bersyarat.
25. Sebagai pemegang jenjang terakhir, saya ingin persetujuan bersyarat **tidak** menerbitkan
    nomor akseptasi, agar kewajiban tidak terbit sebelum syaratnya jelas.

**Riwayat dan keterbacaan**

26. Sebagai anggota komite, saya ingin jenjang yang tidak sempat memutus tercatat sebagai
    `TIDAK_SAMPAI`, agar tidak ada penolakan yang dituliskan atas nama saya.
27. Sebagai anggota komite, saya ingin jabatan yang tercatat pada keputusan saya adalah
    jabatan saya saat itu, bukan jabatan saya hari ini.
28. Sebagai anggota komite, saya ingin satu waktu keputusan per jenjang, agar tidak ada dua
    kolom jam yang perlahan berbeda.
29. Sebagai pembaca klaim, saya ingin melihat seluruh riwayat sirkulasi sebuah klaim, agar
    saya paham mengapa klaim itu berada di keadaan sekarang.
30. Sebagai pembaca klaim, saya ingin hak baca sirkulasi mengikuti hak baca klaim, agar saya
    tidak perlu diberi peran tambahan untuk melihat berkas yang sudah boleh saya buka.
31. Sebagai pembaca klaim, saya ingin melihat apakah sebuah usulan disetujui, disetujui
    bersyarat, atau ditolak, tanpa harus menafsirkan sendiri arti sebuah kode.
32. Sebagai pembaca klaim, saya ingin "komite menyetujui usulan penutupan" tidak terbaca
    sebagai "klaim disetujui".
33. Sebagai pembaca klaim, saya ingin melihat sudah berapa lama sebuah sirkulasi berjalan dan
    berapa lama jenjang aktif menunggu, tanpa ada yang perlu menuliskannya.
34. Sebagai pembaca klaim, saya ingin melihat siapa yang mengajukan sebuah usulan dan apa
    maksudnya, terpisah dari apa akibatnya.

**Roster dan wewenang**

35. Sebagai pemilik proses, saya ingin mengubah siapa menduduki sebuah jenjang tanpa
    menyentuh kode, agar mutasi jabatan tidak menjadi permintaan perubahan.
36. Sebagai pemilik proses, saya ingin mencatat delegasi tetap pada roster, agar penggantian
    sementara tidak dikerjakan dengan meminjam akun.
37. Sebagai pemilik proses, saya ingin mengubah aturan seleksi sebagai data, agar perubahan
    kebijakan kewenangan tidak menunggu rilis.
38. Sebagai pemilik proses, saya ingin sistem menolak konfigurasi seleksi yang tidak menutupi
    seluruh kombinasi, agar lubang tidak ditemukan oleh klaim pertama yang jatuh ke dalamnya.
39. Sebagai pemilik proses, saya ingin tidak ada satu pun aturan yang bercabang berdasarkan
    nama orang, agar perilaku sistem dapat dijelaskan tanpa daftar pengecualian.
40. Sebagai auditor, saya ingin setiap pergantian pemegang jenjang tercatat sebagai peristiwa
    dengan pelaku, alasan, dan waktunya.
41. Sebagai auditor, saya ingin setiap penyuntingan penanda usulan tercatat dengan nilai
    sebelum dan sesudah, agar perubahan saat memutus tidak tersembunyi.

**Migrasi dan pembandingan**

42. **Dicabut 2026-09-21.** Cerita ini menjanjikan yang tidak dispesifikasikan: migrasi
    data berada di *Out of Scope*, dan pagar menang atas cerita. Nomornya tidak dipakai
    ulang; sebab pencabutannya ada di *Further Notes*.
43. Sebagai tim migrasi, saya ingin tali antara klaim dan sirkulasi dibaca dari sisi
    sirkulasi, agar pengenal yang salah tulis pada klaim lama tidak ikut terbawa.
44. Sebagai penguji paritas, saya ingin setiap perbedaan perilaku yang disengaja punya uji
    yang memang dirancang gagal, agar perbedaan itu terbukti dimaksudkan.

---

## Implementation Decisions

### Batas modul dan transaksi

Komite adalah **modul di dalam konteks Klaim**, bukan konteks berjaringan (`D-2`,
`ADR-0032`). Sirkulasi, jenjang, keputusan, dan peristiwanya membentuk satu agregat yang
berubah bersama atau tidak sama sekali. Perubahan yang menyeberang ke klaim terjadi di dalam
batas transaksi yang sama — tidak ada pesan antar-konteks, tidak ada konsistensi akhir,
tidak ada tautan basis data (`ADR-0016`).

Efek ke luar yang tidak dapat dibatalkan berada **di luar** batas transaksi, dijalankan
sesudahnya, dengan niat dicatat sebelum panggilan dan kunci idempotensi per efek (keputusan
beku no. 6). Jenis dan penanganannya per aliran ditetapkan setelah pagarnya dibuka
(lihat **Out of Scope**).

### Mesin keadaan

Diambil dari purwarupa mesin keadaan; bentuk ini mengikat lebih tepat daripada prosa.

```
stateDiagram-v2
    [*] --> BERJALAN : BentukSirkulasi berhasil
    BERJALAN --> BERJALAN : jenjang aktif menyetujui, masih ada jenjang belum memutus
    BERJALAN --> DISETUJUI : jenjang aktif menyetujui, tidak ada jenjang tersisa
    BERJALAN --> DITOLAK : jenjang aktif menolak
    DISETUJUI --> [*]
    DITOLAK --> [*]
```

Keadaan jenjang: `BELUM_MEMUTUS` · `DISETUJUI` · `DITOLAK` · `TIDAK_SAMPAI`.

Tidak ada keadaan awal "dibentuk sebagian": pembentukan adalah satu transaksi, sehingga
sirkulasi lahir sudah lengkap dengan jenjangnya. Tidak ada keadaan kedaluwarsa dan tidak ada
pekerja latar (`E-5`) — umur adalah atribut turunan dari cap waktu yang sudah tersimpan.

**Invarian yang ditegakkan:**

| # | Invarian | Cara membuktikannya dilanggar |
|---|---|---|
| I-1 | Setiap sirkulasi punya sekurang-kurangnya satu jenjang | Bentuk sirkulasi pada kombinasi yang aturan seleksinya menghasilkan roster kosong; harapkan penolakan, bukan berkas |
| I-2 | Derajat unik di dalam satu sirkulasi | Sisipkan dua jenjang berderajat sama; harapkan constraint menolak |
| I-3 | Jenjang aktif adalah jenjang berderajat terendah berkeadaan `BELUM_MEMUTUS` | Setujui sebagai pemegang jenjang berderajat lebih tinggi; harapkan penolakan wewenang |
| I-4 | Paling banyak satu jenjang berkeadaan `DITOLAK` per sirkulasi | Coba catat keputusan kedua pada sirkulasi yang sudah `DITOLAK`; harapkan penolakan keadaan |
| I-5 | Sirkulasi `DITOLAK` tidak menyisakan jenjang `BELUM_MEMUTUS` | Tolak di jenjang tengah; periksa seluruh jenjang di bawahnya berkeadaan `TIDAK_SAMPAI` |
| I-6 | Sirkulasi `DISETUJUI` berarti seluruh jenjangnya `DISETUJUI` | Cari sirkulasi `DISETUJUI` yang punya jenjang bukan `DISETUJUI`; harapkan nol baris |
| I-7 | Sirkulasi `BERJALAN` punya tepat satu jenjang aktif | Hitung jenjang aktif per sirkulasi berjalan; harapkan selalu satu |
| I-8 | Jumlah jenjang tidak pernah disimpan sebagai kolom | Cari kolom yang menyimpannya; harapkan tidak ada. Hitungan adalah turunan (`K5-1`, `F-13`) |
| I-9 | Keputusan tidak pernah berubah sesudah tercatat | Kirim keputusan kedua atas jenjang yang sudah memutus; harapkan penolakan |
| I-10 | Satu waktu keputusan per jenjang | Cari kolom waktu kedua; harapkan tidak ada (keputusan beku no. 3) |
| I-11 | Sirkulasi yang gagal lahir tidak meninggalkan akibat pada klaim | Paksa kegagalan pembentukan; periksa klaim tidak berubah selain maksud (`K5-7`) |

### Jenjang aktif — satu penentu

`ADR-0030` menetapkan satu penentu: jenjang berderajat terendah yang belum memutuskan.
Hitungan jenjang adalah **turunan** dari daftar jenjang, tidak pernah menjadi syarat
penugasan maupun syarat penutupan, dan tidak pernah disimpan sebagai penyimpan kedua.
Penutupan sirkulasi terjadi dari **keadaan keputusan**, bukan dari perbandingan dua bilangan.

### Tabel seleksi dan roster

**Tabel seleksi** menerima tiga masukan — nilai usulan, bagian treaty, bersyarat — dan tiap
**aturan seleksi** menghasilkan satu **kelas kewenangan** (`K5-5`, `K6-3`). Dua kelas
dipertahankan; ambang berbasis nilai tidak dipulihkan sekarang (`K6-1`). Kombinasi yang
tidak tercakup adalah **galat konfigurasi** yang menolak pembentukan.

Nilai penentu kelas adalah nilai **usulan yang sedang diajukan**, bukan nilai baris terakhir
sebuah daftar (`F-17`). **Derajat** diambil dari roster saat pembentukan (`F-14`).
Pembentukan roster sebuah sirkulasi **idempoten**: membentuk ulang menyusun ulang, tidak
menambahkan (`K5-4`).

Isi awal tabel seleksi disemai agar mereproduksi perilaku dua kelas yang berlaku sekarang.
**Angkanya tidak ditulis di spesifikasi ini** dan tidak diambil dari deskripsi langkah
(`K6-1`); ia disemai dari nilai yang berpenulis di sistem lama dan hidup sebagai data.

Jalur `TUTUP_KLAIM` dan `TOLAK_KLAIM` tetap satu jenjang, pemegangnya diselesaikan dari
roster **berdasarkan peran**, dan jenis sirkulasi dinyatakan oleh pemanggil (`K6-4`). Bila
tidak ada baris roster aktif yang memegang peran itu, sirkulasi ditolak saat dibuat.

Roster dimiliki sistem baru, disemai sekali, tanpa pembacaan runtime ke tabel lama (`D-5`).
**Delegasi tetap** adalah atribut roster (`H-3`). **Pengalihan sesaat** tidak diekspos:
modelnya menampung, permukaannya tidak ada (`H-4`, `E-3`).

### Wewenang

`ADR-0006` berlaku penuh: tidak ada hak yang diwarisi dari sistem lama, karena tidak ada
yang tertulis di sana. Gerbang wewenang **menolak**, ditegakkan di server, dan mencocokkan
identitas dengan **kesetaraan penuh** (keputusan beku no. 1). Jenjang adalah **kedudukan**
yang diisi seseorang, bukan seseorang; pergantian pemegang adalah peristiwa tersendiri
dengan pelaku, alasan, dan waktu (`E-3`).

Hak baca sirkulasi **diturunkan** dari hak baca klaim dan bukan peran tersendiri (`K6-2`,
`J-4`).

Tidak ada cabang berdasarkan identitas orang di jalur mana pun (keputusan beku no. 5).

### Pencatatan keputusan

Keputusan adalah **catatan tetap per jenjang**, bukan satu medan yang dipakai bergantian.
Peran dan jabatan disalin ke keputusan pada saat diambil (keputusan beku no. 4). Satu kolom
waktu per jenjang (keputusan beku no. 3).

**Versi usulan** naik hanya ketika penanda usulan berubah. Hanya jenjang berderajat terendah
yang boleh menyuntingnya, dan penyuntingan itu terekam sebagai dua peristiwa dalam satu
transaksi: usulan diubah (membawa nilai sebelum dan sesudah) dan jenjang memutus (`E-1`).
Invarian "jenjang berikutnya melihat keadaan identik" dijamin bentuk **dan** ditegakkan
sebagai penjaga murah di batas tulis: keputusan ditolak bila versi usulan berubah oleh
jenjang selain yang pertama (`J-1`).

Penanda dua-keadaan disimpan sebagai satu boolean dengan satu ejaan di seluruh lapisan;
kosong dan nol tidak pernah berbagi satu kolom (`K5-6`).

### Maksud dan akibat

**Maksud** dicatat pada klaim saat pengajuan dan tidak pernah menjadi **akibat** (`K5-7`).
Keduanya menempati kolom berbeda. Sirkulasi yang gagal lahir atau ditolak meninggalkan
maksudnya tercatat dan akibatnya tidak.

**Fungsi akibat** ditulis sekali di satu tempat dan menerima tiga masukan (`D-3`, `J-3`):

| Jenis sirkulasi | Bersyarat | Hasil | Akibat pada klaim | Nilai tampil |
|---|---|---|---|---|
| `ADJUSTMENT` | tidak | `DISETUJUI` | usulan diakseptasi; nomor akseptasi terbit | Disetujui |
| `ADJUSTMENT` | ya | `DISETUJUI` | usulan diakseptasi bersyarat; nomor **tidak** terbit | Disetujui bersyarat |
| `ADJUSTMENT` | tidak/ya | `DITOLAK` | klaim ditolak | Ditolak |
| `TUTUP_KLAIM` | tidak berlaku | `DISETUJUI` | klaim ditutup | Disetujui |
| `TUTUP_KLAIM` | tidak berlaku | `DITOLAK` | klaim tetap berjalan | Ditolak |
| `TOLAK_KLAIM` | tidak berlaku | `DISETUJUI` | klaim ditolak | Disetujui |
| `TOLAK_KLAIM` | tidak berlaku | `DITOLAK` | klaim tetap berjalan | Ditolak |

Model baca menyajikan ketiga nilai tampil dari fungsi yang sama; laporan tidak pernah
membaca hasil telanjang (`J-3`, keputusan beku no. 9).

### Tulis-balik ke klaim

Kontrak tulis sistem lama dipetakan satu per satu. Seluruhnya terjadi di dalam transaksi
keputusan.

| Sasaran lama | Keputusan |
|---|---|
| Daftar jenjang per usulan | Menjadi jenjang milik sirkulasi; tidak disalin ke usulan |
| Daftar jenjang per klaim | **Turunan tampilan**, tidak disimpan (`CONTEXT.md` bagian 6) |
| Keadaan keputusan pada usulan | Mendarat pada kolom keputusan yang sudah ada di tabel usulan sisi Klaim |
| Nomor dan tanggal akseptasi pada usulan | Mendarat pada kolom yang sudah ada; terbit sekali |
| Penanda dan catatan bersyarat pada usulan | Atribut usulan (`J-3`), kolom yang sudah ada |
| Penanda "sedang disirkulasikan" | **Tidak dimigrasikan sebagai kolom** — turunan dari jenjang, dihitung saat dibaca (`D-2`) |
| Penanda maksud pada klaim | **Kolom baru** pada tabel klaim, terpisah dari keadaan klaim (`K5-7`) |
| Keadaan klaim | Kolom yang sudah ada; ditulis oleh fungsi akibat, bukan oleh langkah keputusan |
| Kronologi | Tabel kronologi yang sudah ada |
| Penanda final XOL dan penanda tolak lawas | **Tidak dibawa.** Keduanya sasaran tulis tanpa kolom pendaratan; yang pertama milik aliran yang berpagar, yang kedua dibaca sesudah halamannya dikosongkan sehingga tidak pernah punya nilai yang berarti |

Keputusan pada jalur bersyarat mendarat pada **jenjang yang memutus**, bukan pada baris
terakhir sebuah daftar, dan bukan pada halaman yang berbeda (`D-1`, `K-07`).

### Skema

Tujuh objek baru di skema `KLAIMNP`. Penamaan mengikuti `SPEC-MODEL-DATA` bagian 16: kata
utuh bahasa Indonesia, `UPPER_SNAKE_CASE`, di bawah 30 byte, constraint berbentuk
`<peran>_<tabel>[_n]` tanpa mengeja kolom. Nilai uang `NUMBER(38,20)` (`ADR-0003`).
Pengenal dari sequence `SQ_<tabel>` (bagian 16.2a). Zona waktu `Asia/Jakarta`.

| Objek | Isi | Asal-usul |
|---|---|---|
| `SIRKULASI` | satu baris per usulan yang diedarkan: pengenal klaim dan usulan, jenis sirkulasi, nomor urut per usulan, kelas kewenangan, keadaan, hasil, cap waktu | agregat `ADR-0032`; nomor urut dari `F-11` |
| `JENJANG_SIRKULASI` | satu baris per jenjang: derajat, pemegang, peran dan jabatan salinan, keadaan jenjang, satu waktu keputusan, komentar | keputusan beku no. 2, 3, 4; `K5-1` |
| `PERISTIWA_SIRKULASI` | satu baris per peristiwa: jenis, pelaku, alasan, waktu, versi usulan yang berlaku | `E-1`, `E-3` |
| `VERSI_USULAN` | satu baris per versi penanda usulan, dengan nilai sebelum dan sesudah | `E-1`, `J-1` |
| `ROSTER_JENJANG` | siapa menduduki jenjang apa: derajat, peran, kelas kewenangan, delegasi tetap, keaktifan | `D-5`, `H-3` |
| `ATURAN_SELEKSI` | satu baris per kombinasi masukan → kelas kewenangan | `K5-5`, `K6-1`, `K6-3` |
| `PERISTIWA_ROSTER` | satu baris per pergantian pemegang sebuah baris roster: jenis, pemegang sebelum dan sesudah, pelaku, alasan, waktu | `E-3`, `H-3`; menggantikan tambalan `F-19` |

Perubahan pada tabel milik sisi Klaim: satu kolom maksud pengajuan pada tabel klaim
(`K5-7`). Kolom keputusan, nomor akseptasi, tanggal akseptasi, penanda bersyarat, dan
catatan bersyarat pada tabel usulan **sudah ada** dan tidak diubah bentuknya; kolom
keputusan di sana sengaja dibiarkan tanpa domain tertutup oleh sisi Klaim karena modul ini
yang menulisnya.

Satu catatan penamaan yang tidak saya ubah: kolom keputusan pada tabel usulan sisi Klaim
memakai kata yang glosarium modul ini larang sebagai nama. Kolom itu milik modul lain dan
sudah ada; spesifikasi ini tidak menamainya ulang, dan tidak membuat satu pun nama baru yang
memakai kata itu.

### Kontrak API

Adapter tipis di atas seam. Hak baca diturunkan dari hak baca klaim (`K6-2`).

| Metode | Jalur | Masukan | Keluaran | Galat | Idempotensi |
|---|---|---|---|---|---|
| POST | `/klaim/{idKlaim}/sirkulasi` | id usulan, jenis sirkulasi, maksud, kunci idempotensi | sirkulasi beserta jenjangnya | `403` bukan pengaju berhak · `409` sudah ada sirkulasi berjalan atas usulan itu · `422` galat konfigurasi seleksi · `422` tidak ada pemegang peran · `503` gangguan | kunci idempotensi + keunikan nomor urut per usulan |
| GET | `/sirkulasi/{id}` | — | sirkulasi, jenjang, keputusan, umur turunan, nilai tampil | `403` · `404` | — |
| GET | `/klaim/{idKlaim}/sirkulasi` | — | daftar sirkulasi klaim itu | `403` · `404` | — |
| GET | `/sirkulasi/menunggu-saya` | — | sirkulasi yang jenjang aktifnya dipegang pemanggil | `403` | — |
| POST | `/sirkulasi/{id}/keputusan` | hasil, komentar, penanda bersyarat, penanda usulan (opsional, hanya jenjang pertama), kunci idempotensi | sirkulasi sesudah keputusan, nomor akseptasi bila terbit | `403` bukan pemegang jenjang aktif · `409` jenjang sudah memutus atau sirkulasi selesai · `409` versi usulan berubah oleh jenjang selain pertama · `422` penanda usulan disunting oleh jenjang bukan pertama · `503` gangguan | kunci idempotensi; percobaan ulang mengembalikan hasil yang sama, bukan nomor kedua (`E-2`) |
| GET · PUT | `/roster-jenjang` | baris roster | roster berlaku | `403` bukan peran administratif | penuh |
| GET · PUT | `/aturan-seleksi` | baris aturan seleksi | tabel seleksi berlaku | `403` · `422` cakupan kombinasi tidak lengkap | penuh |

Galat gangguan **dapat dibedakan** pengguna dari penolakan wewenang (`E-2`), dan masukan
pemutus tidak hilang saat gagal.

**Port hilir — dinamai, kosong, berpagar.** Tidak satu pun jalur berakhir di efek hilir
tanpa port yang dinamai:

| Port | Dipanggil sesudah | Keadaan |
|---|---|---|
| `PORT_PENCATATAN_AKSEPTASI` | nomor akseptasi terbit | kosong — **PAGAR-02** |
| `PORT_PEMBALIKAN_AKSEPTASI` | akibat membatalkan akseptasi yang sudah terbit | kosong — **PAGAR-03** |
| `PORT_PENGIRIMAN_KASIR` | akseptasi tercatat | kosong — **PAGAR-05** |
| `PORT_PENERBITAN_SURAT` | giliran berpindah atau sirkulasi selesai | kosong — **PAGAR-07** |
| `PORT_PENGUNGGAHAN_DOKUMEN` | dokumen terbit | kosong — **PAGAR-08** |

### Layar

Satu-satunya masukan komite adalah **keputusan, komentar, dan penanda usulan**. Seluruh
panel nilai adalah proyeksi baca-saja: komite tidak dapat mengubah satu pun angka uang
(`F-9`), dan tidak pernah menyunting rekening penerima (`F-1`, `F-6`).

Penanda XOL usulan dikunci ke jenjang pertama; label mata uang baca-saja seluruhnya (`J-2`).
Medan induk yang baca-saja **tidak memblokir** keputusan (`E-4`, `F-7`); komentar wajib.
Blok yang di sistem lama bersyarat `1=2`, `1=3`, atau `NEVER` dibuang, dan pembuangannya
dicatat. Tiap kunci ditegakkan **dua kali** — di klien untuk kejelasan, di server untuk
kebenaran; kunci yang hanya ada di klien bukan kunci.

### Penomoran akseptasi

Syarat terbit: `DISETUJUI` pada jenjang terakhir **dan** usulan tidak bersyarat (`J-3`).
Nomor diterbitkan **di dalam transaksi keputusan**, memakai sumber urutan Oracle yang ada,
dengan pembungkus pencatat di `KLAIMNP` (`D-4`). Idempotensi ditegakkan **constraint**,
bukan baca-lalu-tulis. Nomor terbit sekali dan tidak pernah disunting sesudahnya (keputusan
beku no. 7). Gagal menerbitkan berarti **seluruh** transaksi keputusan batal (`E-2`).

**Periode buku** ditentukan di satu lapis, Oracle; sisi aplikasi tidak menggeser ulang dan
tidak menambal (`D-4`). Ambangnya adalah **parameter bernama**, bukan angka di dalam kode —
nilainya berpagar (**PAGAR-01**). Panjang nomor dijaga hilir; fase 1 memakai prosedur lama
apa adanya, dan dasar faktualnya bertanda `TAFSIR (menunggu C-01)`.

---

## Persyaratan

Lapisan ini **diturunkan dari bagian-bagian di atas**, bukan dari pembacaan XML baru dan
bukan dari ketetapan yang ditafsir ulang. Aturan yang tidak ada di *Implementation
Decisions* atau *User Stories* tidak lahir di sini; ia menjadi pagar.

Tiap persyaratan lolos enam uji sebelum ditulis: dapat dikerjakan sendiri · punya jalur
gagal · punya uji · tanda buktinya sah · merujuk ketetapan dengan nomor · tidak menyeberang
pagar. Yang gagal uji pertama **tidak ditulis**; daftarnya ada di akhir bagian ini.

Kolom `Uji` menyebut skenario paritas ronde 5 (`GRILL-05/04-PARITAS.md`) bila ada, dan
`BARU` bila skenarionya belum pernah ditulis. Kolom `Beda AS-IS` merujuk nomor deviasi:
**1 – 13** hidup di `KETETAPAN.md` bagian 9, **14 – 22** di `REGISTER-DEVIASI.md`.

---

### A · Pembentukan sirkulasi — seam `BentukSirkulasi`, aliran A-1b

```
S-001 · Terima jenis sirkulasi dari pemanggil
  Aliran     : A-1b
  Aturan     : BentukSirkulasi menerima jenis sirkulasi sebagai masukan wajib; jenis tidak
               pernah diturunkan dari isi kolom analisis mana pun.
  Bila gagal : jenis kosong atau di luar ketiga nilai yang dikenal → tolak `422` dengan
               galat masukan yang menyebut medan jenis. Tidak ada sirkulasi yang lahir,
               dan klaim tidak berubah — termasuk maksudnya.
  Tanda      : DIPUTUSKAN K6-4
  Beda AS-IS : BERUBAH — sistem lama membedakan jenis dari kolom analisis (F-20). Tidak
               punya nomor deviasi; K6-4 lahir sesudah register deviasi terakhir ditutup.
               **Diberi nomor 2026-09-21: deviasi 23**
  Uji        : P5-04, P5-05, PS-01
  Cerita     : 4
  Seam       : BentukSirkulasi
```

```
S-002 · Tentukan kelas kewenangan dari tabel seleksi
  Aliran     : A-1b
  Aturan     : Kelas kewenangan sebuah sirkulasi ditentukan oleh satu pencarian ke tabel
               seleksi dengan tiga masukan — nilai usulan, bagian treaty, dan bersyarat.
               Tidak ada cabang kelas di dalam kode.
  Bila gagal : salah satu masukan tidak tersedia → tolak `422`; galat menyebut masukan mana
               yang tidak terbaca. Tidak ada sirkulasi yang lahir.
  Tanda      : DIPUTUSKAN K5-5, K6-3
  Beda AS-IS : BERUBAH — deviasi 18; sistem lama menghitung kelas dari tetapan di dalam
               langkah (F-16)
  Uji        : P5-09, P5-10, P5-11
  Cerita     : 1, 5
  Seam       : BentukSirkulasi, KelolaTabelSeleksi
```

```
S-003 · Tolak kombinasi yang tidak tercakup aturan seleksi
  Aliran     : A-1b
  Aturan     : Kombinasi masukan yang tidak cocok dengan satu pun aturan seleksi adalah
               galat konfigurasi, dan galat konfigurasi menolak pembentukan.
  Bila gagal : tolak `422` dengan galat yang menyebut ketiga nilai masukan yang tidak
               tercakup, agar yang salah diperbaiki di data acuan. Tidak ada sirkulasi
               yang lahir, tidak ada roster yang tersusun, klaim tidak berubah.
  Tanda      : DIPUTUSKAN K5-5
  Beda AS-IS : BERUBAH — deviasi 12, deviasi 18; sistem lama memanggil penyaring roster
               dengan ambang tak terisi dan melahirkan sirkulasi tanpa jenjang
  Uji        : P5-09, P5-11
  Cerita     : 11
  Seam       : BentukSirkulasi
```

```
S-004 · Pakai nilai usulan yang sedang diajukan sebagai penentu kelas
  Aliran     : A-1b
  Aturan     : Masukan nilai bagi tabel seleksi adalah nilai usulan yang sedang diajukan,
               bukan nilai baris terakhir sebuah daftar dan bukan jumlah seluruh usulan.
  Bila gagal : usulan yang diajukan tidak dapat ditunjuk dari masukan → tolak `422`;
               pembentukan tidak berjalan.
  Tanda      : EVIDENCED CreateChildKomiteCNP_Act·11 (F-17)
  Beda AS-IS : BERUBAH — deviasi 8. Ujinya wajib memakai klaim dengan sekurang-kurangnya
               dua usulan berbeda kelas; pada klaim berusulan tunggal deviasi ini lolos
               tanpa terlihat
  Uji        : P5-13
  Cerita     : 5
  Seam       : BentukSirkulasi
```

```
S-005 · Ambil derajat dari roster pada saat pembentukan
  Aliran     : A-1b
  Aturan     : Derajat tiap jenjang disalin dari roster pada saat sirkulasi dibentuk, urut
               menaik, dan tidak dihitung ulang sesudahnya. Jumlah jenjang adalah akibat
               dari isi roster, bukan masukan.
  Bila gagal : dua baris roster yang terpilih membawa derajat yang sama → tolak `422`
               sebagai galat konfigurasi roster; tidak ada sirkulasi yang lahir.
  Tanda      : EVIDENCED FilterEmailKomiteWithLimit (F-14); CreateChildKomiteCNP_Act·9,
               26.13–26.15 (F-13)
  Beda AS-IS : SAMA — urutan derajat menaik dipertahankan
  Uji        : P5-01, P5-03
  Cerita     : 5
  Seam       : BentukSirkulasi
```

```
S-006 · Susun ulang roster sirkulasi, jangan menambahkan
  Aliran     : A-1b
  Aturan     : Membentuk roster sebuah sirkulasi menyusun ulang isinya dari awal;
               pembentukan kedua atas perkara yang sama menghasilkan roster yang sama
               panjangnya, bukan roster yang bertambah.
  Bila gagal : penyusunan ulang tidak dapat diselesaikan → seluruh pembentukan batal;
               tidak ada roster separuh yang tertinggal.
  Tanda      : DIPUTUSKAN K5-4
  Beda AS-IS : BERUBAH — deviasi 17; pada jalur bersyarat sistem lama melewati pembersihan
               lalu menambahkan
  Uji        : P5-06
  Cerita     : 8
  Seam       : BentukSirkulasi
```

```
S-007 · Tolak sirkulasi yang tidak punya jenjang
  Aliran     : A-1b
  Aturan     : Sirkulasi tanpa sekurang-kurangnya satu jenjang ditolak pada saat dibuat.
               Tidak ada jalur yang menutup sirkulasi tanpa keputusan.
  Bila gagal : roster terpilih kosong → tolak `422` dengan galat yang menyebut kelas
               kewenangan yang tidak punya pemegang aktif. Tidak ada berkas yang lahir.
  Tanda      : DIPUTUSKAN H-6, keputusan beku no. 8
  Beda AS-IS : BERUBAH — deviasi 6; sistem lama tetap membuat berkasnya lalu diam (F-22)
  Uji        : P5-09, P5-11
  Cerita     : 11
  Seam       : BentukSirkulasi
```

```
S-008 · Bentuk jalur tutup dan tolak klaim dari peran, bukan dari nama
  Aliran     : A-1b
  Aturan     : Sirkulasi menutup dan menolak klaim tetap satu jenjang, dan pemegangnya
               diselesaikan dari roster berdasarkan peran.
  Bila gagal : tidak ada baris roster aktif yang memegang peran itu → tolak `422` dengan
               galat yang menyebut peran yang kosong. Tidak ada sirkulasi yang lahir, dan
               maksud tidak ditulis ke klaim.
  Tanda      : DIPUTUSKAN K6-4
  Beda AS-IS : BERUBAH — nama orang, surel, inisial, dan dua sebutan jabatan untuk satu
               kedudukan ditulis di dalam rule (F-20). Sebagian tertutup deviasi 20; bagian
               yang khusus jalur ini tidak punya nomor deviasi, karena K6-4 lahir sesudah
               register deviasi terakhir ditutup. **Diberi nomor 2026-09-21: deviasi 24
               untuk pemegang dari peran, deviasi 25 untuk penolakan saat peran kosong**
  Uji        : P5-04, P5-05, PS-02, PS-03
  Cerita     : 2, 3
  Seam       : BentukSirkulasi, KelolaRoster
```

```
S-009 · Bentuk sirkulasi dan akibatnya dalam satu transaksi
  Aliran     : A-1b
  Aturan     : Pembentukan sirkulasi, penulisan pengenalnya ke klaim, penyimpanan klaim,
               dan pencatatan niat pemberitahuan berhasil bersama atau gagal bersama.
               Tidak ada langkah sesudah pembentukan yang berjalan ketika pembentukan
               tidak jadi.
  Bila gagal : gangguan di tengah → seluruh transaksi batal dan pemanggil menerima `503`
               yang dapat dibedakan dari penolakan. Klaim tidak membawa pengenal
               sirkulasi, dan tidak ada niat pemberitahuan yang tercatat.
  Tanda      : DIPUTUSKAN K5-2, ADR-0031
  Beda AS-IS : BERUBAH — deviasi 15; penjaga sistem lama hanya melewati satu langkah,
               langkah sesudahnya menulis pengenal milik sirkulasi lain lalu menyimpan
  Uji        : P5-07
  Cerita     : 7
  Seam       : BentukSirkulasi
```

```
S-010 · Buat kegagalan pembentukan terlihat dan tercatat
  Aliran     : A-1b
  Aturan     : Setiap jalur pembentukan yang berakhir tanpa sirkulasi berakhir dengan
               kegagalan yang terlihat oleh pengaju dan tercatat. Penanda yang tidak dibaca
               siapa pun bukan penanganan galat.
  Bila gagal : jalur keluar yang tidak punya galat bernama adalah cacat persyaratan ini
               sendiri; uji memeriksa bahwa tidak ada jalur keluar senyap — pemanggil
               selalu menerima `422` atau `503`, tidak pernah `200` tanpa sirkulasi.
  Tanda      : DIPUTUSKAN K5-3
  Beda AS-IS : BERUBAH — deviasi 16; sistem lama keluar diam-diam dan menyimpan penanda
               yang tidak dibaca satu rule pun
  Uji        : P5-08
  Cerita     : 6
  Seam       : BentukSirkulasi
```

```
S-011 · Jangan tinggalkan akibat pada klaim ketika pembentukan gagal
  Aliran     : A-1b
  Aturan     : Klaim yang pembentukan sirkulasinya gagal tidak berubah selain maksudnya —
               tidak menerima pengenal sirkulasi, tidak berpindah keadaan, dan tidak
               menerima akibat apa pun.
  Bila gagal : ditemukan klaim yang membawa pengenal sirkulasi tanpa sirkulasi yang
               bersesuaian → itu pelanggaran invarian I-11 dan uji gagal.
  Tanda      : DIPUTUSKAN K5-2, K5-7
  Beda AS-IS : BERUBAH — deviasi 22
  Uji        : P5-07, P5-08
  Cerita     : 7
  Seam       : BentukSirkulasi
```

```
S-012 · Catat maksud pada kolom tersendiri saat pengajuan
  Aliran     : A-1b
  Aturan     : Maksud pengajuan dicatat pada klaim pada saat pengajuan, pada kolom yang
               terpisah dari keadaan klaim, dan tidak pernah menjadi akibat.
  Bila gagal : maksud tidak dapat ditulis → seluruh pembentukan batal; klaim tidak berubah
               sama sekali.
  Tanda      : DIPUTUSKAN K5-7
  Beda AS-IS : BERUBAH — deviasi 22; sistem lama menulis penanda maksud dan penanda akibat
               ke tempat yang sama, pada langkah yang membuat berkas anak
  Uji        : P5-07, P5-08
  Cerita     : 7, 34
  Seam       : BentukSirkulasi
```

```
S-013 · Jadikan pembentukan idempoten terhadap kuncinya
  Aliran     : A-1b
  Aturan     : BentukSirkulasi membawa kunci idempotensi; permintaan berulang dengan kunci
               yang sama mengembalikan sirkulasi yang sama, bukan sirkulasi kedua.
  Bila gagal : kunci sama dengan muatan berbeda → tolak `409`; sirkulasi yang sudah ada
               tidak berubah.
  Tanda      : DIPUTUSKAN E-2 (bentuk kunci), K5-2
  Beda AS-IS : BERUBAH — tidak punya nomor deviasi: sistem lama tidak punya kunci
               idempotensi sama sekali, sehingga tidak ada perilaku yang dapat dibandingkan
  Uji        : BARU
  Cerita     : 8
  Seam       : BentukSirkulasi
```

```
S-014 · Beri tiap sirkulasi nomor urut pada usulannya
  Aliran     : A-1b
  Aturan     : Satu usulan boleh disirkulasikan lebih dari sekali. Tiap sirkulasi atas
               usulan yang sama membawa nomor urut tersendiri pada usulan itu, dan dua
               putaran dibedakan oleh nomor urut, bukan oleh cap waktu.
  Bila gagal : sirkulasi baru diajukan sementara satu sirkulasi atas usulan itu masih
               berjalan → tolak `409`; sirkulasi yang berjalan tidak berubah.
  Tanda      : EVIDENCED KomitePostAdjustment·16; AdjustmentDetailNP (F-11)
  Beda AS-IS : SAMA pada kemampuannya; BERUBAH pada pembedanya — nomor urut adalah ikatan
               tambahan F-11, tanpa nomor deviasi karena ia menambah, bukan mengganti
  Uji        : BARU
  Cerita     : 9, 10
  Seam       : BentukSirkulasi
```

```
S-015 · Tolak pada validasi, jangan menandai
  Aliran     : A-1b
  Aturan     : Validasi pembentukan adalah bagian dari kasus-guna yang menolak. Tidak ada
               jalur yang memasang pesan lalu melanjutkan.
  Bila gagal : validasi tidak terpenuhi → tolak `422` dan sebutkan medan yang gagal;
               tidak ada berkas, tidak ada roster, tidak ada tulisan ke klaim.
  Tanda      : DIPUTUSKAN H-5, K5-2
  Beda AS-IS : BERUBAH — deviasi 13; gerbang sistem lama menandai galat dan bergantung
               pada penjaga di langkah lain (F-21)
  Uji        : BARU
  Cerita     : 6
  Seam       : BentukSirkulasi
```

```
S-016 · Periksa hak pengaju sebelum membentuk
  Aliran     : A-1b
  Aturan     : Pembentukan sirkulasi ditolak bila pemanggil tidak berhak mengajukan usulan
               atas klaim itu. Pemeriksaan terjadi di server.
  Bila gagal : tolak `403`, dapat dibedakan dari `503`; tidak ada sirkulasi, dan klaim
               tidak berubah.
  Tanda      : DIPUTUSKAN ADR-0006; EVIDENCED nol pyPrivilegeName berisi pada 59 rule (F-24)
  Beda AS-IS : BERUBAH — deviasi 21
  Uji        : P5-16
  Cerita     : 1
  Seam       : BentukSirkulasi
```

---

### B · Memutus — seam `CatatKeputusan`, aliran A-1a

```
S-017 · Terima keputusan hanya dari pemegang jenjang aktif
  Aliran     : A-1a
  Aturan     : Satu-satunya yang boleh memutus adalah pemegang jenjang aktif. Pencocokan
               identitas memakai kesetaraan penuh, ditegakkan di server, dan permintaan
               yang gagal ditolak.
  Bila gagal : tolak `403` dengan pesan yang menyatakan giliran bukan miliknya, berbeda
               bentuk dari `503`. Tidak ada keputusan yang tersimpan, keadaan jenjang dan
               keadaan sirkulasi tidak berubah.
  Tanda      : DIPUTUSKAN keputusan beku no. 1, K5-1
  Beda AS-IS : BERUBAH — deviasi 1, deviasi 21; sistem lama mencocokkan sebagai potongan
               teks dan hanya memasang pesan
  Uji        : BARU; P5-16 untuk sisi haknya
  Cerita     : 18
  Seam       : CatatKeputusan
```

```
S-018 · Hitung jenjang aktif dari satu penentu
  Aliran     : A-1a
  Aturan     : Jenjang aktif adalah jenjang berderajat terendah yang belum memutuskan.
               Jumlah jenjang adalah turunan, tidak pernah menjadi syarat penugasan maupun
               syarat penutupan, dan tidak pernah disimpan sebagai penyimpan kedua.
  Bila gagal : sirkulasi berjalan yang tidak punya tepat satu jenjang aktif adalah
               pelanggaran invarian I-7; uji menghitung per sirkulasi dan harapkan satu.
  Tanda      : DIPUTUSKAN K5-1, ADR-0030
  Beda AS-IS : BERUBAH — deviasi 14; sistem lama memakai dua penentu yang tidak saling
               memeriksa
  Uji        : P5-01, P5-02, P5-02b
  Cerita     : 12, 18
  Seam       : CatatKeputusan, BacaSirkulasi
```

```
S-019 · Jadikan keputusan catatan tetap
  Aliran     : A-1a
  Aturan     : Keputusan sebuah jenjang — hasil, komentar, dan waktunya — tercatat sebagai
               catatan tetap milik jenjang itu, dan tidak dapat disunting sesudah tercatat.
  Bila gagal : keputusan kedua atas jenjang yang sudah memutus → tolak `409`; catatan yang
               ada tidak berubah.
  Tanda      : DIPUTUSKAN keputusan beku no. 2, E-1
  Beda AS-IS : SAMA pada kemampuannya; ketetapan yang mengikat adalah invarian I-9
  Uji        : BARU
  Cerita     : 13, 14, 23
  Seam       : CatatKeputusan
```

```
S-020 · Simpan satu waktu keputusan per jenjang
  Aliran     : A-1a
  Aturan     : Waktu sebuah keputusan disimpan sekali, pada satu kolom milik jenjang itu.
  Bila gagal : ditemukan kolom waktu kedua pada katalog skema → uji DDL gagal.
  Tanda      : DIPUTUSKAN keputusan beku no. 3
  Beda AS-IS : BERUBAH — deviasi 5; sistem lama menulis satu fakta ke dua kolom tanggal
  Uji        : BARU — uji DDL, lihat S-068
  Cerita     : 28
  Seam       : CatatKeputusan, DDL
```

```
S-021 · Salin peran dan jabatan pada saat keputusan diambil
  Aliran     : A-1a
  Aturan     : Peran dan jabatan pemegang disalin ke keputusan pada saat keputusan diambil,
               dari roster, sehingga mutasi jabatan tidak mengubah riwayat.
  Bila gagal : baris roster pemegang tidak terbaca pada saat memutus → tolak `503`;
               keputusan tidak tersimpan.
  Tanda      : DIPUTUSKAN keputusan beku no. 4
  Beda AS-IS : BERUBAH — deviasi 2, deviasi 9, deviasi 20; sistem lama menimpa jabatan
               dengan teks di dalam rule dan dengan pengenal baris roster (F-18, F-19)
  Uji        : P5-14, P5-15
  Cerita     : 27
  Seam       : CatatKeputusan
```

```
S-022 · Tandai jenjang yang tidak sempat memutus sebagai TIDAK_SAMPAI
  Aliran     : A-1a
  Aturan     : Penolakan oleh jenjang aktif menutup sirkulasi. Seluruh jenjang berderajat
               lebih tinggi yang belum memutus berkeadaan `TIDAK_SAMPAI` — bukan menolak,
               dan tanpa komentar maupun waktu milik orang lain.
  Bila gagal : sirkulasi berkeadaan ditolak yang masih menyisakan jenjang belum memutus
               adalah pelanggaran invarian I-5; uji memeriksa seluruh jenjang di bawahnya.
  Tanda      : DIPUTUSKAN keputusan beku no. 2
  Beda AS-IS : BERUBAH — deviasi 3; sistem lama menuliskan penolakan atas nama jenjang
               yang tidak pernah bertindak, lengkap dengan komentar dan jam milik penolak
  Uji        : BARU
  Cerita     : 26
  Seam       : CatatKeputusan
```

```
S-023 · Tutup sirkulasi ketika jenjang terakhir menyetujui
  Aliran     : A-1a
  Aturan     : Persetujuan oleh jenjang aktif ketika tidak ada jenjang tersisa menutup
               sirkulasi sebagai disetujui. Penutupan lahir dari keadaan keputusan, bukan
               dari perbandingan dua bilangan.
  Bila gagal : sirkulasi disetujui yang memuat jenjang bukan disetujui adalah pelanggaran
               invarian I-6; uji mencari baris semacam itu dan harapkan nol.
  Tanda      : DIPUTUSKAN K5-1, ADR-0030
  Beda AS-IS : BERUBAH — deviasi 14; penutupan sistem lama bergantung pada perbandingan
               hitungan yang dipalsukan saat penolakan
  Uji        : P5-02
  Cerita     : 24
  Seam       : CatatKeputusan
```

```
S-024 · Izinkan hanya jenjang berderajat terendah menyunting penanda usulan
  Aliran     : A-1a
  Aturan     : Penanda usulan hanya dapat disunting oleh pemegang jenjang berderajat
               terendah, dan hanya pada saat ia memutus. Penyuntingan terekam sebagai dua
               peristiwa dalam satu transaksi: usulan diubah — membawa nilai sebelum dan
               sesudah — dan jenjang memutus.
  Bila gagal : penyuntingan datang dari jenjang bukan yang pertama → tolak `422`; penanda
               usulan, versi usulan, dan keputusan sama-sama tidak tersimpan.
  Tanda      : DIPUTUSKAN E-1
  Beda AS-IS : SAMA — paritas F-2; yang berubah adalah tempat penegakannya, lihat S-025
  Uji        : BARU
  Cerita     : 16, 41
  Seam       : CatatKeputusan
```

```
S-025 · Tegakkan versi usulan sebagai penjaga di batas tulis
  Aliran     : A-1a
  Aturan     : Versi usulan naik hanya ketika penanda usulan berubah. Tiap keputusan
               mencatat versi yang diputusnya, dan keputusan ditolak bila versi usulan
               berubah oleh jenjang selain yang pertama.
  Bila gagal : tolak `409` dengan pesan yang menyebut versi yang diharapkan dan versi yang
               ditemukan; keputusan tidak tersimpan dan versi usulan tidak naik.
  Tanda      : DIPUTUSKAN J-1
  Beda AS-IS : BERUBAH — deviasi 7; sistem lama menjamin ini lewat bentuk layar saja, dan
               bentuk yang dijamin 296 kendali layar berubah tanpa ada yang memberi tahu
  Uji        : BARU
  Cerita     : 17
  Seam       : CatatKeputusan
```

```
S-026 · Simpan penanda dua-keadaan sebagai satu boolean
  Aliran     : A-1a
  Aturan     : Penanda yang hanya punya dua keadaan disimpan sebagai satu boolean dengan
               satu ejaan pada seluruh lapisan. Kosong dan nol tidak pernah berbagi satu
               kolom, dan tidak ada kolom yang artinya bergantung pada kosong atau berisi.
  Bila gagal : nilai di luar kedua keadaan itu ditolak `422` di batas tulis; tidak ada
               yang tersimpan.
  Tanda      : DIPUTUSKAN K5-6
  Beda AS-IS : BERUBAH — deviasi 19
  Uji        : P5-05
  Cerita     : 31
  Seam       : CatatKeputusan, DDL
```

```
S-027 · Jadikan pencatatan keputusan idempoten terhadap kuncinya
  Aliran     : A-1a
  Aturan     : CatatKeputusan membawa kunci idempotensi. Percobaan ulang dengan kunci yang
               sama mengembalikan hasil yang sama — termasuk nomor akseptasi yang sama —
               bukan keputusan kedua dan bukan nomor kedua.
  Bila gagal : kunci sama dengan muatan berbeda → tolak `409`; keadaan tidak berubah.
  Tanda      : DIPUTUSKAN E-2
  Beda AS-IS : BERUBAH — tidak punya nomor deviasi: sistem lama tidak punya kunci
               idempotensi, sehingga tidak ada perilaku yang dapat dibandingkan
  Uji        : BARU
  Cerita     : 22
  Seam       : CatatKeputusan
```

```
S-028 · Simpan keputusan seluruhnya atau tidak sama sekali
  Aliran     : A-1a
  Aturan     : Keputusan, akibatnya pada klaim, penerbitan nomor bila syaratnya terpenuhi,
               dan pencatatan niat efek hilir berada dalam satu transaksi.
  Bila gagal : gangguan di tengah → `503`; tidak ada keputusan tersimpan, tidak ada nomor
               terbit, tidak ada akibat pada klaim, tidak ada niat tercatat.
  Tanda      : DIPUTUSKAN E-2, ADR-0031
  Beda AS-IS : BERUBAH — deviasi 15 pada bentuknya
  Uji        : BARU
  Cerita     : 19
  Seam       : CatatKeputusan
```

```
S-029 · Bedakan gangguan dari penolakan wewenang
  Aliran     : A-1a
  Aturan     : Galat gangguan sistem dan penolakan wewenang adalah dua bentuk galat yang
               berbeda di permukaan, sehingga pemutus tahu harus menunggu atau menghubungi
               orang.
  Bila gagal : satu bentuk galat dipakai untuk keduanya → uji gagal; uji memeriksa kode
               dan bentuk pesan, bukan kalimatnya.
  Tanda      : DIPUTUSKAN E-2
  Beda AS-IS : BERUBAH — tidak punya nomor deviasi: sistem lama tidak menolak sama sekali,
               sehingga tidak ada pembedaan yang dapat dibandingkan
  Uji        : BARU
  Cerita     : 20
  Seam       : CatatKeputusan
```

```
S-030 · Pertahankan masukan pemutus ketika penyimpanan gagal
  Aliran     : A-1a, A-2
  Aturan     : Kegagalan menyimpan keputusan tidak menghapus komentar, penanda bersyarat,
               maupun penanda usulan yang sudah diisi pemutus.
  Bila gagal : masukan hilang sesudah galat → uji gagal.
  Tanda      : DIPUTUSKAN E-2
  Beda AS-IS : BERUBAH — tidak punya nomor deviasi
  Uji        : BARU
  Cerita     : 21
  Seam       : CatatKeputusan, Layar
```

```
S-031 · Hitung akibat di satu tempat lewat fungsi akibat
  Aliran     : A-1a
  Aturan     : Akibat pada klaim dihitung dari tiga masukan — jenis sirkulasi, bersyarat,
               dan hasil — oleh satu fungsi yang ditulis sekali. Hasil tidak pernah berarti
               nasib klaim.
  Bila gagal : kombinasi yang tidak tercakup ketujuh baris tabel fungsi akibat → tolak
               `422` sebagai galat konfigurasi; keputusan tidak tersimpan.
  Tanda      : DIPUTUSKAN D-3, J-3
  Beda AS-IS : SAMA pada keluarannya; BERUBAH pada tempatnya — sistem lama menyebarkannya
               di dua jalur. Tanpa nomor deviasi karena perilakunya identik
  Uji        : BARU — ketujuh baris
  Cerita     : 13, 14, 32
  Seam       : CatatKeputusan, BacaSirkulasi
```

```
S-032 · Simpan keadaan sirkulasi terpisah dari akibatnya pada klaim
  Aliran     : A-1a
  Aturan     : Keadaan sirkulasi dan akibat pada klaim menempati kolom yang berbeda,
               sehingga "komite menyetujui usulan penutupan" tidak pernah terbaca sebagai
               "klaim disetujui".
  Bila gagal : satu kolom dipakai untuk keduanya → uji DDL gagal.
  Tanda      : DIPUTUSKAN keputusan beku no. 9
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; pemisahan ini bentuk simpan, bukan perbedaan
               keluaran yang dapat diuji paritasnya
  Uji        : BARU — uji DDL
  Cerita     : 32
  Seam       : DDL, BacaSirkulasi
```

```
S-033 · Tulis balik ke klaim di dalam transaksi keputusan
  Aliran     : A-1a
  Aturan     : Seluruh tulis-balik ke sisi Klaim — keadaan keputusan pada usulan, nomor dan
               tanggal akseptasi, penanda dan catatan bersyarat, keadaan klaim, dan
               kronologi — terjadi di dalam transaksi keputusan, tanpa pesan antar-konteks
               dan tanpa tautan basis data.
  Bila gagal : salah satu tulisan gagal → seluruh transaksi batal; `503`.
  Tanda      : DIPUTUSKAN D-2, ADR-0016, ADR-0032
  Beda AS-IS : SAMA — di sistem lama pun satu commit
  Uji        : BARU
  Cerita     : 29
  Seam       : CatatKeputusan
```

```
S-034 · Jangan simpan penanda "sedang disirkulasikan"
  Aliran     : A-1a
  Aturan     : Keadaan "sedang disirkulasikan" dihitung dari jenjang saat dibaca; ia tidak
               disimpan sebagai kolom.
  Bila gagal : ditemukan kolom yang menyimpannya → uji DDL gagal.
  Tanda      : DIPUTUSKAN D-2
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; kolom lama tidak dimigrasikan, sehingga
               tidak ada keluaran yang dapat dibandingkan
  Uji        : BARU — uji DDL
  Cerita     : 33
  Seam       : DDL
```

```
S-035 · Daratkan keputusan bersyarat pada jenjang yang memutus
  Aliran     : A-1a
  Aturan     : Pada jalur bersyarat, keputusan tercatat pada jenjang yang memutus — bukan
               pada baris terakhir sebuah daftar, dan bukan pada halaman yang berbeda.
  Bila gagal : keputusan mendarat pada jenjang lain → pelanggaran invarian I-3 dan uji
               gagal.
  Tanda      : DIPUTUSKAN D-1 (K-07); DIKUKUHKAN GRILL-06/01-PEMBACAAN P6-4
  Beda AS-IS : BERUBAH — deviasi 4
  Uji        : BARU
  Cerita     : 15
  Seam       : CatatKeputusan
```

```
S-036 · Jangan bercabang berdasarkan identitas orang
  Aliran     : A-1a, A-1b
  Aturan     : Tidak ada aturan di jalur mana pun yang bercabang berdasarkan identitas
               seseorang. Nilai lingkungan, kode akuntansi, alamat layanan, dan daftar
               penerima surat adalah konfigurasi.
  Bila gagal : ditemukan cabang semacam itu → uji gagal; ujinya adalah pencarian atas
               kode, bukan pemanggilan operasi.
  Tanda      : DIPUTUSKAN keputusan beku no. 5, H-7
  Beda AS-IS : BERUBAH — deviasi 10, deviasi 20
  Uji        : P5-15
  Cerita     : 39
  Seam       : CatatKeputusan, KelolaRoster
```

```
S-037 · Catat niat efek hilir sebelum memanggilnya
  Aliran     : A-1a
  Aturan     : Efek ke luar yang tidak dapat dibatalkan berada di luar batas transaksi dan
               dijalankan sesudahnya. Niat memanggilnya dicatat sebelum panggilan, dengan
               kunci idempotensi per efek. Tiap efek melewati port yang dinamai.
  Bila gagal : port dipanggil tanpa niat tercatat lebih dulu → uji gagal. Kegagalan port
               tidak membatalkan keputusan yang sudah tercatat; ia meninggalkan niat yang
               belum terpenuhi dan terbaca.
  Tanda      : DIPUTUSKAN keputusan beku no. 6
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; isi tiap efek berpagar, sehingga
               pembandingnya tidak dipegang
  Uji        : BARU — port diuji sebagai port kosong
  Cerita     : —  (lihat *Further Notes*, keputusan teknis no. 2)
  Seam       : CatatKeputusan
```

---

### C · Penomoran akseptasi — seam `CatatKeputusan` dan `DDL`, aliran A-3

```
S-038 · Terbitkan nomor hanya pada persetujuan terakhir yang tidak bersyarat
  Aliran     : A-3
  Aturan     : Nomor akseptasi terbit bila dan hanya bila jenjang terakhir menyetujui dan
               usulannya tidak bersyarat.
  Bila gagal : usulan bersyarat pada persetujuan terakhir → sirkulasi tetap tertutup
               sebagai disetujui, nomor tidak terbit, dan nilai tampil berbunyi disetujui
               bersyarat.
  Tanda      : DIPUTUSKAN J-3
  Beda AS-IS : SAMA — paritas F-5
  Uji        : BARU
  Cerita     : 24, 25
  Seam       : CatatKeputusan
```

```
S-039 · Terbitkan nomor di dalam transaksi keputusan
  Aliran     : A-3
  Aturan     : Nomor diterbitkan di dalam transaksi keputusan, memakai sumber urutan Oracle
               yang ada, lewat satu pembungkus pencatat di skema modul ini. Nomor bukan
               efek luar; ia fakta internal yang menjadi syarat efek luar.
  Bila gagal : sumber urutan tidak menjawab → `503`, seluruh transaksi batal, tidak ada
               keputusan tersimpan.
  Tanda      : DIPUTUSKAN D-4 — TAFSIR (menunggu C-01) pada dasar faktualnya
  Beda AS-IS : SAMA — prosedur lama dipakai apa adanya pada fase 1
  Uji        : BARU
  Cerita     : 24
  Seam       : CatatKeputusan
```

```
S-040 · Tegakkan keunikan nomor dengan constraint
  Aliran     : A-3
  Aturan     : Idempotensi penerbitan nomor ditegakkan oleh constraint basis data, bukan
               oleh baca-lalu-tulis.
  Bila gagal : penerbitan kedua atas usulan yang sama ditolak oleh constraint; transaksi
               batal dan pemanggil menerima hasil penerbitan pertama.
  Tanda      : DIPUTUSKAN D-4 — TAFSIR (menunggu C-01)
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; bentuk penegakan, bukan keluaran
  Uji        : BARU — uji DDL, lihat S-066
  Cerita     : 24
  Seam       : DDL
```

```
S-041 · Terbitkan nomor sekali dan jangan pernah menyuntingnya
  Aliran     : A-3
  Aturan     : Nomor akseptasi diterbitkan sekali, dari nilai yang seluruhnya berpenulis,
               dan tidak pernah disunting sesudah terbit.
  Bila gagal : permintaan menyunting nomor → tolak `409`; nilai yang ada tidak berubah.
  Tanda      : DIPUTUSKAN keputusan beku no. 7
  Beda AS-IS : SAMA — paritas mutlak atas angka yang sudah terbit (D-1)
  Uji        : BARU
  Cerita     : 24
  Seam       : CatatKeputusan
```

```
S-042 · Batalkan seluruh keputusan ketika penerbitan gagal
  Aliran     : A-3
  Aturan     : Gagal menerbitkan nomor membatalkan seluruh transaksi keputusan. Tidak ada
               keputusan tersimpan tanpa nomor ketika syarat terbit terpenuhi.
  Bila gagal : ditemukan keputusan final tidak bersyarat tanpa nomor → uji gagal.
  Tanda      : DIPUTUSKAN E-2
  Beda AS-IS : BERUBAH — tanpa nomor deviasi
  Uji        : BARU
  Cerita     : 19, 24
  Seam       : CatatKeputusan
```

```
S-043 · Tentukan periode buku di satu lapis
  Aliran     : A-3
  Aturan     : Periode buku ditentukan di satu lapis, yaitu Oracle. Sisi aplikasi tidak
               menggeser ulang dan tidak menambal. Ambang pergeseran periode adalah
               parameter bernama, bukan angka di dalam kode.
  Bila gagal : parameter tidak terisi → penerbitan tidak dijalankan dan transaksi batal
               dengan `503`; tidak ada nomor yang terbit dengan periode terkaan.
  Tanda      : DIPUTUSKAN D-4 — TAFSIR (menunggu C-01)
  Beda AS-IS : SAMA pada mekanikanya. **Nilai** ambangnya tidak ditulis — PAGAR-01
  Uji        : BARU — mekanika saja; uji nilai menunggu PAGAR-01 dibuka
  Cerita     : 24
  Seam       : CatatKeputusan
```

---

### D · Membaca — seam `BacaSirkulasi`

```
S-044 · Turunkan hak baca sirkulasi dari hak baca klaim
  Aliran     : A-1a
  Aturan     : Siapa pun yang berhak membaca sebuah klaim berhak membaca sirkulasinya. Hak
               baca sirkulasi bukan peran tersendiri.
  Bila gagal : pemanggil tidak berhak membaca klaimnya → `403`; tidak ada sebagian isi
               yang bocor lewat pesan galat.
  Tanda      : DIPUTUSKAN K6-2, J-4
  Beda AS-IS : SAMA pada cakupan pembacanya; BERUBAH pada adanya pemeriksaan — deviasi 21
  Uji        : P5-16
  Cerita     : 30
  Seam       : BacaSirkulasi
```

```
S-045 · Hitung umur saat dibaca
  Aliran     : A-1a
  Aturan     : Umur sirkulasi dan umur tunggu jenjang aktif adalah atribut turunan dari cap
               waktu yang sudah tersimpan, dihitung saat dibaca. Tidak ada keadaan
               kedaluwarsa dan tidak ada pekerja latar.
  Bila gagal : ditemukan kolom umur yang tersimpan → uji DDL gagal.
  Tanda      : DIPUTUSKAN E-5
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; sistem lama tidak punya umur sama sekali
               (F-4), sehingga tidak ada pembanding
  Uji        : BARU
  Cerita     : 33
  Seam       : BacaSirkulasi
```

```
S-046 · Sajikan nilai tampil dari fungsi akibat
  Aliran     : A-1a
  Aturan     : Model baca menyajikan nilai tampil — disetujui, disetujui bersyarat,
               ditolak — dari fungsi akibat yang sama. Laporan tidak pernah membaca hasil
               telanjang.
  Bila gagal : ditemukan pembaca yang menafsirkan hasil sendiri → uji gagal.
  Tanda      : DIPUTUSKAN J-3, keputusan beku no. 9
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; sistem lama tidak punya model baca tersendiri
  Uji        : BARU — ketujuh baris fungsi akibat diperiksa dari model baca
  Cerita     : 31, 32
  Seam       : BacaSirkulasi
```

```
S-047 · Sajikan riwayat sirkulasi sebuah klaim
  Aliran     : A-1a
  Aturan     : Pembaca klaim dapat melihat seluruh sirkulasi klaim itu beserta jenjang,
               keputusan, dan peristiwanya, urut waktu.
  Bila gagal : klaim tidak ada atau tidak boleh dibaca → `404` atau `403`.
  Tanda      : DIPUTUSKAN ADR-0032
  Beda AS-IS : SAMA — komentar komite sudah menjadi bagian berkas klaim (J-4)
  Uji        : BARU
  Cerita     : 29
  Seam       : BacaSirkulasi
```

```
S-048 · Sajikan daftar sirkulasi yang menunggu pemanggil
  Aliran     : A-1a
  Aturan     : Pemegang jenjang dapat membaca daftar sirkulasi yang jenjang aktifnya
               dipegangnya. Daftar itu dihitung dari jenjang aktif, bukan dari keranjang
               yang ditulis terpisah.
  Bila gagal : pemanggil bukan pemegang jenjang mana pun → daftar kosong, bukan galat.
  Tanda      : DIPUTUSKAN K5-1; EVIDENCED empat nama keranjang tidak terpakai (P5-03)
  Beda AS-IS : SAMA
  Uji        : P5-03
  Cerita     : 12
  Seam       : BacaSirkulasi
```

```
S-049 · Sajikan pengaju dan maksud terpisah dari akibat
  Aliran     : A-1a
  Aturan     : Pembaca melihat siapa mengajukan sebuah usulan dan apa maksudnya, sebagai
               dua medan yang terpisah dari akibat pada klaim.
  Bila gagal : maksud dan akibat tersaji sebagai satu medan → uji gagal.
  Tanda      : DIPUTUSKAN K5-7
  Beda AS-IS : BERUBAH — deviasi 22
  Uji        : P5-07
  Cerita     : 34
  Seam       : BacaSirkulasi
```

---

### E · Roster — seam `KelolaRoster`

```
S-050 · Miliki roster di sistem baru
  Aliran     : A-1b
  Aturan     : Roster dimiliki sistem baru dan disemai sekali. Tidak ada pembacaan runtime
               ke tabel sistem lama, dan tidak ada tautan basis data ke sistem lain.
  Bila gagal : roster tidak terbaca saat pembentukan → `503`; tidak ada sirkulasi yang
               lahir.
  Tanda      : DIPUTUSKAN D-5, ADR-0016
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; yang berubah adalah pemilik rumahnya, bukan
               keluarannya
  Uji        : BARU
  Cerita     : 35
  Seam       : KelolaRoster
```

```
S-051 · Simpan delegasi tetap sebagai atribut roster
  Aliran     : A-1b
  Aturan     : Seorang yang memegang jenjang milik orang lain sampai dicabut dicatat
               sebagai atribut pada baris roster, bukan sebagai penukaran nama di dalam
               aturan.
  Bila gagal : delegasi menunjuk orang yang tidak aktif di roster → tolak `422` saat
               disimpan; baris roster tidak berubah.
  Tanda      : DIPUTUSKAN H-3
  Beda AS-IS : BERUBAH — deviasi 20; sistem lama menambal satu baris kode yang menukar
               satu nama dengan nama lain (F-19)
  Uji        : P5-15
  Cerita     : 36
  Seam       : KelolaRoster
```

```
S-052 · Catat pergantian pemegang sebagai peristiwa
  Aliran     : A-1b
  Aturan     : Setiap pergantian pemegang sebuah jenjang tercatat sebagai peristiwa dengan
               pelaku, alasan, dan waktunya. Jenjang adalah kedudukan yang diisi seseorang,
               bukan seseorang.
  Bila gagal : peristiwa tidak dapat ditulis → penyimpanan roster batal seluruhnya.
  Tanda      : DIPUTUSKAN E-3
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; sistem lama tidak mencatat pergantian
  Uji        : BARU
  Cerita     : 40
  Seam       : KelolaRoster
  Catatan    : mendarat di `PERISTIWA_ROSTER` — objek ketujuh, ditambahkan 2026-09-21 karena
               `PERISTIWA_SIRKULASI` tidak dapat menampungnya. Bunyi `Aturan`, `Bila gagal`,
               `Tanda`, dan `Beda AS-IS` tidak berubah; yang berubah hanya bahwa persyaratan
               ini kini punya tempat mendarat, sehingga ia **lolos `T-1`**
```

```
S-053 · Batasi pengelolaan roster pada peran administratif
  Aliran     : A-1b
  Aturan     : Membaca dan mengubah roster menolak pemanggil yang tidak memegang peran
               administratif.
  Bila gagal : tolak `403`; roster tidak berubah.
  Tanda      : DIPUTUSKAN ADR-0006; bentuknya dari *Kontrak API*
  Beda AS-IS : BERUBAH — deviasi 21
  Uji        : P5-16
  Cerita     : 35, 37
  Seam       : KelolaRoster, KelolaTabelSeleksi
  Catatan    : **nama** peran administratifnya tidak ditetapkan di sini — lihat daftar
               gagal T-1 di akhir bagian ini
```

```
S-054 · Jangan bawa pengecualian per orang ke dalam roster
  Aliran     : A-1b
  Aturan     : Roster tidak memuat baris yang berlaku bagi satu orang tertentu di luar
               aturan umum. Jalur khusus satu operator pembuat berkas tidak dimigrasikan.
  Bila gagal : baris semacam itu tidak dapat dibentuk; penyimpanan ditolak `422`.
  Tanda      : DIPUTUSKAN H-7, keputusan beku no. 5
  Beda AS-IS : BERUBAH — deviasi 10
  Uji        : BARU
  Cerita     : 39
  Seam       : KelolaRoster
```

---

### F · Tabel seleksi — seam `KelolaTabelSeleksi`

```
S-055 · Simpan aturan seleksi sebagai data
  Aliran     : A-1b
  Aturan     : Satu aturan seleksi adalah satu baris: satu kombinasi masukan, satu kelas
               kewenangan. Aturan diubah sebagai data, tanpa menunggu rilis.
  Bila gagal : baris yang menduplikasi kombinasi yang sudah ada → tolak `422`; tabel tidak
               berubah.
  Tanda      : DIPUTUSKAN K5-5, H-1
  Beda AS-IS : BERUBAH — deviasi 18; sistem lama menyebarkannya sebagai tetapan di dalam
               langkah
  Uji        : P5-09, P5-10, P5-11
  Cerita     : 37
  Seam       : KelolaTabelSeleksi
```

```
S-056 · Tolak konfigurasi seleksi yang cakupannya tidak lengkap
  Aliran     : A-1b
  Aturan     : Penyimpanan tabel seleksi ditolak bila hasilnya meninggalkan satu pun
               kombinasi masukan tanpa aturan.
  Bila gagal : tolak `422` dengan galat yang menyebut kombinasi yang tidak tercakup; tabel
               yang berlaku tidak berubah.
  Tanda      : DIPUTUSKAN K5-5
  Beda AS-IS : BERUBAH — deviasi 18; lubangnya ditemukan oleh klaim pertama yang jatuh ke
               dalamnya, bukan oleh penyimpanan konfigurasi
  Uji        : BARU
  Cerita     : 38
  Seam       : KelolaTabelSeleksi
```

```
S-057 · Semai isi awal dari nilai yang berpenulis
  Aliran     : A-1b
  Aturan     : Isi awal tabel seleksi disemai agar mereproduksi perilaku dua kelas yang
               berlaku sekarang, dari nilai yang berpenulis di sistem lama. Angka yang
               hanya ada pada deskripsi langkah tidak dipakai, dan tidak ada angka ambang
               di dalam kode.
  Bila gagal : nilai semai tidak dapat ditunjuk ke langkah yang menuliskannya → semai
               tidak dijalankan, dan itu dilaporkan sebagai galat konfigurasi, bukan
               ditebak.
  Tanda      : DIPUTUSKAN K6-1, K6-3
  Beda AS-IS : SAMA — isi awal mereproduksi perilaku sekarang; yang berubah hanya tempat
               tinggalnya
  Uji        : P5-10, P5-12
  Cerita     : 37
  Seam       : KelolaTabelSeleksi
```

---

### G · Layar — seam `Layar`, aliran A-2

Seam layar **sempit dengan sengaja**. Uji di sini memeriksa bahwa penegakan pertama **ada**;
kebenaran aturannya diuji di seam perintah. Kunci yang hanya ada di klien bukan kunci.

```
S-058 · Batasi masukan layar pada tiga medan
  Aliran     : A-2
  Aturan     : Satu-satunya masukan komite di layar adalah keputusan, komentar, dan penanda
               usulan.
  Bila gagal : medan lain terkirim → diabaikan di server, dan permintaan ditolak `422` bila
               medan itu berbeda dari nilai tersimpan.
  Tanda      : EVIDENCED ShowTransfer (F-1, F-6, F-9)
  Beda AS-IS : SAMA
  Uji        : BARU
  Cerita     : 13, 16
  Seam       : Layar, CatatKeputusan
```

```
S-059 · Sajikan seluruh panel nilai sebagai baca-saja
  Aliran     : A-2
  Aturan     : Seluruh panel nilai adalah proyeksi baca-saja. Komite tidak dapat mengubah
               satu pun angka uang dan tidak pernah menyunting rekening penerima.
  Bila gagal : kendali nilai yang dapat disunting ditemukan di layar → uji gagal.
  Tanda      : EVIDENCED ShowTransfer — 16 properti nilai baca-saja (F-9); 16 kemunculan
               medan rekening baca-saja (F-6)
  Beda AS-IS : SAMA
  Uji        : BARU
  Cerita     : 13
  Seam       : Layar
```

```
S-060 · Kunci penanda XOL ke jenjang pertama dan label mata uang seluruhnya
  Aliran     : A-2
  Aturan     : Penanda XOL usulan hanya dapat disunting oleh jenjang berderajat terendah.
               Label mata uang baca-saja bagi seluruh jenjang.
  Bila gagal : penyuntingan dari jenjang lain ditolak `422` di server, meski klien
               mengizinkannya.
  Tanda      : DIPUTUSKAN J-2; EVIDENCED ShowTransfer (F-8)
  Beda AS-IS : BERUBAH — deviasi 7
  Uji        : BARU
  Cerita     : 16, 17
  Seam       : Layar, CatatKeputusan
```

```
S-061 · Jangan biarkan medan induk baca-saja memblokir keputusan
  Aliran     : A-2
  Aturan     : Medan induk yang baca-saja tidak memblokir keputusan. Komentar wajib.
  Bila gagal : keputusan tanpa komentar → tolak `422`; keputusan tidak tersimpan.
  Tanda      : DIPUTUSKAN E-4; EVIDENCED ShowTransfer — lima medan induk baca-saja dan
               wajib sekaligus (F-7)
  Beda AS-IS : SAMA — aturan wajibnya memang tidak pernah menyala di sistem lama
  Uji        : BARU
  Cerita     : 13, 14
  Seam       : Layar, CatatKeputusan
```

```
S-062 · Buang blok yang syaratnya tidak pernah menyala, dan catat pembuangannya
  Aliran     : A-2
  Aturan     : Blok layar yang di sistem lama bersyarat pada perbandingan yang tidak pernah
               dapat benar, atau bersyarat "tidak pernah", dibuang. Tiap pembuangan
               tercatat beserta syarat lamanya.
  Bila gagal : blok terbawa tanpa catatan → uji gagal.
  Tanda      : EVIDENCED ShowTransfer
  Beda AS-IS : SAMA pada apa yang terlihat — blok itu tidak pernah tampil
  Uji        : BARU
  Cerita     : —  (lihat *Further Notes*, keputusan teknis no. 1)
  Seam       : Layar
```

```
S-063 · Tegakkan tiap kunci dua kali
  Aliran     : A-2
  Aturan     : Tiap kunci ditegakkan di klien untuk kejelasan dan di server untuk
               kebenaran.
  Bila gagal : kunci yang hanya ada di klien → uji gagal. Uji layar memeriksa bahwa
               penegakan klien **ada**; kebenaran aturannya diuji di seam perintah.
  Tanda      : DIPUTUSKAN J-1 (alasannya), ADR-0006
  Beda AS-IS : BERUBAH — deviasi 7, deviasi 21; sistem lama menegakkan hanya di klien
  Uji        : BARU
  Cerita     : 18, 23
  Seam       : Layar
```

---

### H · Skema — seam `DDL`

Enam persyaratan berikut diuji sebagai **pemeriksaan skema**, bukan pemeriksaan perilaku.
Constraint yang membuat keadaan terlarang tidak dapat ditulis lebih kuat daripada uji yang
memanggil operasi, dan sejalan `D-4` yang memilih constraint di atas baca-lalu-tulis.

```
S-064 · Jadikan derajat unik di dalam satu sirkulasi
  Aliran     : A-1b
  Aturan     : Constraint menolak dua jenjang berderajat sama pada satu sirkulasi.
  Bila gagal : penyisipan kedua ditolak basis data; transaksi pembentukan batal.
  Tanda      : DIPUTUSKAN K5-1 — invarian I-2
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; sistem lama tidak punya penegak bentuk
  Uji        : BARU — uji DDL
  Cerita     : 5
  Seam       : DDL
```

```
S-065 · Jadikan nomor urut sirkulasi unik per usulan
  Aliran     : A-1b
  Aturan     : Constraint menolak dua sirkulasi bernomor urut sama atas satu usulan.
  Bila gagal : penyisipan kedua ditolak basis data; pembentukan batal dengan `409`.
  Tanda      : EVIDENCED KomitePostAdjustment·16 (F-11)
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; ikatan tambahan F-11
  Uji        : BARU — uji DDL
  Cerita     : 10
  Seam       : DDL
```

```
S-066 · Jadikan nomor akseptasi unik
  Aliran     : A-3
  Aturan     : Constraint menolak nomor akseptasi yang sama terbit dua kali.
  Bila gagal : penerbitan kedua ditolak basis data; transaksi keputusan batal.
  Tanda      : DIPUTUSKAN D-4, keputusan beku no. 7 — TAFSIR (menunggu C-01)
  Beda AS-IS : SAMA pada angkanya; BERUBAH pada penegaknya
  Uji        : BARU — uji DDL
  Cerita     : 24
  Seam       : DDL
```

```
S-067 · Jangan simpan jumlah jenjang sebagai kolom
  Aliran     : A-1a
  Aturan     : Tidak ada kolom yang menyimpan jumlah jenjang sebuah sirkulasi. Hitungan
               adalah turunan.
  Bila gagal : kolom semacam itu ditemukan di katalog skema → uji DDL gagal.
  Tanda      : DIPUTUSKAN K5-1 — invarian I-8; EVIDENCED F-13
  Beda AS-IS : BERUBAH — deviasi 14
  Uji        : P5-02b — sisi perilakunya; uji DDL untuk bentuknya
  Cerita     : 28
  Seam       : DDL
```

```
S-068 · Simpan satu kolom waktu keputusan per jenjang
  Aliran     : A-1a
  Aturan     : Tabel jenjang memuat tepat satu kolom waktu keputusan.
  Bila gagal : kolom waktu kedua ditemukan di katalog skema → uji DDL gagal.
  Tanda      : DIPUTUSKAN keputusan beku no. 3 — invarian I-10
  Beda AS-IS : BERUBAH — deviasi 5
  Uji        : BARU — uji DDL
  Cerita     : 28
  Seam       : DDL
```

```
S-069 · Namai objek menurut aturan penamaan yang berlaku
  Aliran     : A-1a, A-1b, A-3
  Aturan     : Nama tabel, kolom, constraint, dan sequence mengikuti `SPEC-MODEL-DATA`
               bagian 16. Nilai uang memakai presisi desimal yang ditetapkan, tidak pernah
               bilangan mengambang. Cap waktu memakai zona `Asia/Jakarta`.
  Bila gagal : pengenal melewati batas panjang atau memakai kata yang dilarang glosarium →
               alat pemeriksa penamaan keluar dengan kode gagal.
  Tanda      : DIPUTUSKAN aturan kerja F, aturan kerja G, ADR-0003
  Beda AS-IS : BERUBAH — tanpa nomor deviasi; penamaan bukan perilaku
  Uji        : BARU — uji DDL, prior art di sisi Klaim
  Cerita     : —  (lihat *Further Notes*, keputusan teknis no. 3)
  Seam       : DDL
```

```
S-070 · Simpan tali klaim–sirkulasi di sisi sirkulasi
  Aliran     : A-1b
  Aturan     : Hubungan antara klaim dan sirkulasi dibaca dari sisi sirkulasi. Pengenal
               sirkulasi yang tersimpan pada klaim lama tidak dipercaya sebagai tali.
  Bila gagal : pembacaan yang bersandar pada pengenal di sisi klaim → uji gagal.
  Tanda      : DIPUTUSKAN ADR-0031; EVIDENCED CreateChildKomiteCNP_Act·29 (N-07)
  Beda AS-IS : BERUBAH — deviasi 15
  Uji        : P5-07
  Cerita     : 43
  Seam       : DDL, BacaSirkulasi
```

---

### Yang tidak ditulis sebagai persyaratan karena gagal T-1

Ketiganya menunggu keputusan yang belum diambil. Menuliskannya sekarang berarti mengarang
jawabannya, dan tiket yang lahir darinya akan berhenti di tengah.

| Yang tidak ditulis | Apa yang ditunggunya |
|---|---|
| **Nama peran administratif** yang berhak mengelola roster dan tabel seleksi | Keputusan pemilik proses. `H-4` menyebut default "peran administratif tersendiri di luar keanggotaan komite" sebagai default, bukan sebagai ketetapan. `S-053` menulis gerbangnya; yang tidak ditulis adalah nama perannya |
| **Bentuk migrasi sirkulasi berjalan** | Satu cacah yang belum diambil: berapa klaim di produksi membawa pengenal sirkulasi yang bukan miliknya (`INVENTARIS-BUKTI.md` §2.5 baris 6). Yang sudah diputuskan ditulis sebagai `S-070`; bentuk migrasinya tidak |
| **Keadaan `KonversiKlaim_Act` dan `InsertJsonClaimTreatyNonProp_act`** | Keputusan pemilik proses, **bukan bukti**. Berkas keempat rule-nya ada di repo (`INVENTARIS-BUKTI.md` §2.4, dimutakhirkan 2026-09-21), sehingga ini bukan lubang bukti dan tidak boleh dipagari. Yang belum ada adalah port hilirnya. Lihat *Ketertelusuran*, catatan di bawah tabel |


---

## Model data

Detail atas sepuluh persyaratan ber-seam `DDL` (`S-064` … `S-070`) dan objek yang disebut
*Implementation Decisions → Skema*. **Bukan persyaratan baru**: tiap kolom menunjuk `S-xxx`
yang menuntutnya, atau menyebut ketetapan yang melahirkannya.

Penamaan mengikuti `SPEC-MODEL-DATA` bagian 16 — kata utuh bahasa Indonesia,
`UPPER_SNAKE_CASE`, **di bawah 30 byte**, constraint berbentuk `<peran>_<tabel>[_n]` yang
**tidak mengeja kolom**, awalan tetap `PK_` `UQ_` `FK_` `CK_` `IX_` `SQ_`. Pengenal dari
sequence `SQ_<tabel>` (bagian 16.2a). Cap waktu bertimezone, zona `Asia/Jakarta` (aturan
kerja `G`). Akar kata diambil dari `CONTEXT.md`, tidak dibentuk ulang.

**Nol objek untuk A-4, A-5, dan A-6.** Tidak ada tabel akseptasi, tabel pembalik, log kiriman
kasir, maupun tabel dokumen di sini — keempatnya berpagar, dan yang menyentuhnya hanyalah
port kosong `S-037`.

**Satu-satunya kolom uang ada di `ATURAN_SELEKSI`**, dan ia batas, bukan nilai. Lima objek
lainnya **tidak memuat satu pun kolom uang**: nilai usulan tinggal di sisi Klaim, komite
tidak dapat mengubah satu pun angka uang (`F-9`), dan berkas ini tidak memutuskan apa pun
tentang uang.

---

### `SIRKULASI`

Satu baris per usulan yang diedarkan. Pengenal dari `SQ_SIRKULASI`.

| Kolom | Tipe | Wajib | Kunci / constraint | S-xxx | Asal-usul di sistem lama |
|---|---|:-:|---|---|---|
| `ID_SIRKULASI` | `NUMBER(19)` | ya | `PK_SIRKULASI` | `S-069` | tidak ada padanan — pengenal berkas anak Pega; lahir dari `ADR-0032` |
| `ID_KLAIM` | `NUMBER(19)` | ya | `FK_SIRKULASI_1` → `KLAIM` | `S-070` | `pxCoverInsKey`, satu-satunya tali ke Klaim (`PENGETAHUAN.md` §7.1) |
| `ID_USULAN` | `NUMBER(19)` | ya | `FK_SIRKULASI_2` → `ADJUSTMENT` | `S-070` | `ClaimData.AdjustmentList(n)` |
| `JENIS_SIRKULASI` | `VARCHAR2(16 CHAR)` | ya | `CK_SIRKULASI_1` — domain tertutup tiga nilai | `S-001` | `TransferType` dan `TypeComentAnalysis`, dua sumber yang dilebur (`F-20`) |
| `NOMOR_URUT` | `NUMBER(5)` | ya | `UQ_SIRKULASI_1` bersama `ID_USULAN` | `S-014`, `S-065` | tidak ada padanan — ikatan tambahan `F-11`; sistem lama hanya punya cap waktu |
| `KELAS_KEWENANGAN` | `VARCHAR2(32 CHAR)` | ya | — | `S-002` | `Flagkomite` diterjemahkan menjadi `Param.LIMIT_BOTTOM` |
| `KEADAAN` | `VARCHAR2(16 CHAR)` | ya | `CK_SIRKULASI_2` — `BERJALAN`, `DISETUJUI`, `DITOLAK` | `S-023` | keadaan berkas Pega; mesin keadaan `KomiteTreaty_Flow` |
| `HASIL` | `VARCHAR2(16 CHAR)` | tidak | `CK_SIRKULASI_3` — `DISETUJUI`, `DITOLAK` | `S-031` | `AcceptStatus`, `"1"` terima dan `"2"` tolak |
| `PENANDA_BERSYARAT` | `NUMBER(1)` | ya | `CK_SIRKULASI_4` — `IN (0,1)` | `S-026`, `S-038` | `IsSubjectivity` |
| `CATATAN_BERSYARAT` | `VARCHAR2(2000 CHAR)` | tidak | — | `S-038` | `SubjectivityNote` |
| `DIBENTUK_OLEH` | `VARCHAR2(64 CHAR)` | ya | — | `S-016`, `S-049` | tidak ada padanan — pengaju tidak pernah dicatat pada berkas sirkulasi |
| `DIBENTUK_PADA` | `TIMESTAMP WITH TIME ZONE` | ya | — | `S-045` | `pxCreateDateTime` |
| `DISELESAIKAN_PADA` | `TIMESTAMP WITH TIME ZONE` | tidak | — | `S-045` | tidak ada padanan — lahir dari `E-5` |
| `KUNCI_IDEMPOTENSI` | `VARCHAR2(64 CHAR)` | ya | `UQ_SIRKULASI_2` | `S-013` | tidak ada padanan — lahir dari `E-2` |

`HASIL` kosong selama `KEADAAN = BERJALAN`, dan kosong **bukan** nol (`ADR-0019`). Domain
`KELAS_KEWENANGAN` sengaja **tidak** ditutup `CHECK`: kelas adalah data yang dikelola
(`K5-5`), dan menutupnya di DDL memindahkan kembali kebijakan ke tempat yang butuh rilis.
Yang ditegakkan adalah keberadaan pemegangnya, lewat `S-007` pada saat pembentukan.

### `JENJANG_SIRKULASI`

Satu baris per jenjang. Pengenal dari `SQ_JENJANG_SIRKULASI`.

| Kolom | Tipe | Wajib | Kunci / constraint | S-xxx | Asal-usul di sistem lama |
|---|---|:-:|---|---|---|
| `ID_JENJANG` | `NUMBER(19)` | ya | `PK_JENJANG_SIRKULASI` | `S-069` | tidak ada padanan — baris `KomiteList` tidak berpengenal sendiri |
| `ID_SIRKULASI` | `NUMBER(19)` | ya | `FK_JENJANG_SIRKULASI_1` | `S-005` | keanggotaan `pyWorkPage.KomiteList` |
| `DERAJAT` | `NUMBER(3)` | ya | **`UQ_JENJANG_SIRKULASI_1`** bersama `ID_SIRKULASI` — **mewujudkan `I-2`** | `S-005`, `S-064` | `.DEGREE`, satu-satunya kolom berurutan pada penyaring roster (`F-14`) |
| `PEMEGANG` | `VARCHAR2(64 CHAR)` | ya | — | `S-017` | `KomiteID`, berisi nama pengguna |
| `PERAN_SALINAN` | `VARCHAR2(64 CHAR)` | ya | — | `S-021` | `KomitePost` |
| `JABATAN_SALINAN` | `VARCHAR2(128 CHAR)` | ya | — | `S-021` | `IDKomite`, yang di sistem lama justru ditimpa pengenal baris roster (`F-18`) |
| `KEADAAN_JENJANG` | `VARCHAR2(16 CHAR)` | ya | `CK_JENJANG_SIRKULASI_1` — `BELUM_MEMUTUS`, `DISETUJUI`, `DITOLAK`, `TIDAK_SAMPAI` | `S-022`, `S-026` | `KomiteAproval` bernilai `0`/`1`/`2`; keadaan keempat tidak ada padanan, lahir dari keputusan beku no. 2 |
| `KOMENTAR` | `VARCHAR2(2000 CHAR)` | tidak | — | `S-019` | `KomiteComment` |
| `DIPUTUS_PADA` | `TIMESTAMP WITH TIME ZONE` | tidak | **satu kolom — mewujudkan `I-10`** | `S-020`, `S-068` | `DateApproval` **dan** `DateApprove`, dua kolom untuk satu fakta, dilebur jadi satu |
| `ID_VERSI_USULAN` | `NUMBER(19)` | tidak | `FK_JENJANG_SIRKULASI_2` → `VERSI_USULAN` | `S-025` | tidak ada padanan — lahir dari `E-1` dan `J-1` |
| `KUNCI_IDEMPOTENSI` | `VARCHAR2(64 CHAR)` | tidak | `UQ_JENJANG_SIRKULASI_2` | `S-027` | tidak ada padanan — lahir dari `E-2` |

**`I-8` diwujudkan sebagai ketiadaan.** Tidak ada kolom yang menyimpan jumlah jenjang, di
tabel ini maupun di `SIRKULASI`. Ujinya adalah pencarian atas katalog skema yang harus pulang
kosong (`S-067`). `KOMITECOUNT` dan `FLAGONGOINGCOMMITTE` **tidak dimigrasikan sebagai
kolom** — keduanya turunan, dihitung saat dibaca (`D-2`).

### `PERISTIWA_SIRKULASI`

Satu baris per peristiwa. Pengenal dari `SQ_PERISTIWA_SIRKULASI`.

| Kolom | Tipe | Wajib | Kunci / constraint | S-xxx | Asal-usul di sistem lama |
|---|---|:-:|---|---|---|
| `ID_PERISTIWA` | `NUMBER(19)` | ya | `PK_PERISTIWA_SIRKULASI` | `S-069` | tidak ada padanan — lahir dari `E-1`, `E-3` |
| `ID_SIRKULASI` | `NUMBER(19)` | ya | `FK_PERISTIWA_SIRKULASI_1` | `S-052` | tidak ada padanan |
| `JENIS_PERISTIWA` | `VARCHAR2(32 CHAR)` | ya | `CK_PERISTIWA_SIRKULASI_1` — domain tertutup | `S-024`, `S-052` | tidak ada padanan — sistem lama tidak mencatat peristiwa apa pun |
| `PELAKU` | `VARCHAR2(64 CHAR)` | ya | — | `S-052` | tidak ada padanan |
| `ALASAN` | `VARCHAR2(2000 CHAR)` | tidak | — | `S-052` | tidak ada padanan |
| `TERJADI_PADA` | `TIMESTAMP WITH TIME ZONE` | ya | — | `S-045` | tidak ada padanan |
| `ID_VERSI_USULAN` | `NUMBER(19)` | tidak | `FK_PERISTIWA_SIRKULASI_2` | `S-024` | tidak ada padanan |

### `VERSI_USULAN`

Satu baris per versi penanda usulan. Pengenal dari `SQ_VERSI_USULAN`.

| Kolom | Tipe | Wajib | Kunci / constraint | S-xxx | Asal-usul di sistem lama |
|---|---|:-:|---|---|---|
| `ID_VERSI` | `NUMBER(19)` | ya | `PK_VERSI_USULAN` | `S-069` | tidak ada padanan — suntingan sistem lama menimpa nilai sebelumnya tanpa jejak |
| `ID_SIRKULASI` | `NUMBER(19)` | ya | `FK_VERSI_USULAN_1` | `S-025` | tidak ada padanan |
| `NOMOR_VERSI` | `NUMBER(5)` | ya | `UQ_VERSI_USULAN_1` bersama `ID_SIRKULASI` | `S-025` | tidak ada padanan |
| `PENANDA` | `VARCHAR2(64 CHAR)` | ya | `CK_VERSI_USULAN_1` — domain tertutup penanda usulan | `S-024` | nama properti penanda pada berkas anak |
| `NILAI_SEBELUM` | `VARCHAR2(256 CHAR)` | tidak | — | `S-024` | tidak ada padanan |
| `NILAI_SESUDAH` | `VARCHAR2(256 CHAR)` | tidak | — | `S-024` | tidak ada padanan |
| `DIUBAH_OLEH` | `VARCHAR2(64 CHAR)` | ya | — | `S-024` | tidak ada padanan |
| `DIUBAH_PADA` | `TIMESTAMP WITH TIME ZONE` | ya | — | `S-024` | tidak ada padanan |

`NILAI_SEBELUM` dan `NILAI_SESUDAH` menyimpan **penanda**, bukan uang — komite tidak dapat
mengubah satu pun angka uang (`F-9`), sehingga tidak ada nilai uang yang pernah berversi.

### `ROSTER_JENJANG`

Daftar acuan: siapa menduduki jenjang apa. Pengenal dari `SQ_ROSTER_JENJANG`.

| Kolom | Tipe | Wajib | Kunci / constraint | S-xxx | Asal-usul di sistem lama |
|---|---|:-:|---|---|---|
| `ID_ROSTER` | `NUMBER(19)` | ya | `PK_ROSTER_JENJANG` | `S-069` | baris `EMAILKOMITE` |
| `PEMEGANG` | `VARCHAR2(64 CHAR)` | ya | — | `S-050` | `EMAILKOMITE`, kolom pengguna |
| `DERAJAT` | `NUMBER(3)` | ya | `UQ_ROSTER_JENJANG_1` bersama `KELAS_KEWENANGAN`, atas baris aktif saja | `S-005` | `.DEGREE` |
| `PERAN` | `VARCHAR2(64 CHAR)` | ya | — | `S-008` | `KomitePost`; di jalur tutup/tolak ditulis sebagai nama orang di dalam aturan (`F-20`) |
| `JABATAN` | `VARCHAR2(128 CHAR)` | ya | — | `S-021` | `.JABATAN`, yang ditimpa daftar berisi `.ID` (`F-18`) |
| `KELAS_KEWENANGAN` | `VARCHAR2(32 CHAR)` | ya | — | `S-002` | `.LIMIT_BOTTOM` hasil terjemahan `Flagkomite` |
| `DELEGASI_KEPADA` | `VARCHAR2(64 CHAR)` | tidak | `CK_ROSTER_JENJANG_1` — tidak boleh sama dengan `PEMEGANG` | `S-051` | tidak ada padanan — lahir dari `H-3`; sistem lama menambalnya satu baris kode (`F-19`) |
| `AKTIF` | `NUMBER(1)` | ya | `CK_ROSTER_JENJANG_2` — `IN (0,1)` | `S-026`, `S-054` | tidak ada padanan |

`UQ_ROSTER_JENJANG_1` berlaku atas baris aktif saja dan diwujudkan sebagai index unik
berbasis fungsi; baris yang dinonaktifkan tetap tersimpan, karena mencabut seseorang dari
roster tidak boleh menghapus jejaknya. `DELEGASI_KEPADA` kosong berarti tidak ada delegasi,
dan kosong di sini **bukan** nol (`ADR-0019`, `K5-6`).

### `ATURAN_SELEKSI`

Satu baris per kombinasi masukan. Pengenal dari `SQ_ATURAN_SELEKSI`.

| Kolom | Tipe | Wajib | Kunci / constraint | S-xxx | Asal-usul di sistem lama |
|---|---|:-:|---|---|---|
| `ID_ATURAN` | `NUMBER(19)` | ya | `PK_ATURAN_SELEKSI` | `S-069` | tidak ada padanan — lahir dari `K5-5` |
| `BATAS_NILAI_BAWAH` | `NUMBER(38,20)` | tidak | — | `S-002`, `S-057` | empat tetapan di dalam langkah (`F-16`); nilainya tidak disalin ke sini |
| `BATAS_NILAI_ATAS` | `NUMBER(38,20)` | tidak | — | `S-002`, `S-057` | idem |
| `BATAS_BAGIAN_BAWAH` | `NUMBER(9,6)` | tidak | — | `S-002` | pembanding bagian treaty (`F-16`) |
| `BATAS_BAGIAN_ATAS` | `NUMBER(9,6)` | tidak | — | `S-002` | idem |
| `PENANDA_BERSYARAT` | `NUMBER(1)` | tidak | `CK_ATURAN_SELEKSI_1` — `IN (0,1)` | `S-002`, `S-026` | sub-langkah bersyarat yang deskripsinya sama persis dengan sub-langkah biasa (`K6-3`) |
| `KELAS_KEWENANGAN` | `VARCHAR2(32 CHAR)` | ya | — | `S-002` | `Flagkomite` |
| `AKTIF` | `NUMBER(1)` | ya | `CK_ATURAN_SELEKSI_2` — `IN (0,1)` | `S-026` | tidak ada padanan |
| — | — | — | `UQ_ATURAN_SELEKSI_1` atas keempat batas dan penanda, baris aktif saja | `S-055` | tidak ada padanan |

Keempat batas dan `PENANDA_BERSYARAT` boleh kosong, dan **kosong berarti "tidak membatasi"**,
bukan nol (`ADR-0019`). Itulah bentuk yang membuat isi awal dapat mereproduksi perilaku dua
kelas tanpa memulihkan ambang berbasis nilai (`K6-1`): baris semai meninggalkan batas
nilainya kosong. Keempat kolom batas nilai memakai presisi uang `NUMBER(38,20)` (`ADR-0003`)
agar ambang yang kelak diputuskan pemilik proses tidak perlu mengubah bentuk tabel.

`CK_ATURAN_SELEKSI_1` menjaga `K5-6` di kolom yang memang punya tiga keadaan berarti — `0`,
`1`, dan kosong — sehingga kosong dan nol tidak pernah berbagi arti.

### `PERISTIWA_ROSTER`

Satu baris per pergantian pemegang sebuah baris roster. Pengenal dari `SQ_PERISTIWA_ROSTER`.

Objek ini ada karena `S-052` tidak punya tempat mendarat. `PERISTIWA_SIRKULASI` tidak dapat
menampungnya: pengenal sirkulasinya **wajib**, dan pergantian roster tidak punya sirkulasi.
Menjadikan pengenal itu boleh kosong akan melahirkan satu kolom yang artinya bergantung pada
kosong-atau-berisi — persis yang `K5-6` larang — dan menggabungkan dua jenis peristiwa yang
tidak pernah dibaca bersama, hidup di seam berbeda, dan berumur berbeda.

| Kolom | Tipe | Wajib | Kunci / constraint | S-xxx | Asal-usul di sistem lama |
|---|---|:-:|---|---|---|
| `ID_PERISTIWA` | `NUMBER(19)` | ya | `PK_PERISTIWA_ROSTER` | `S-069` | tidak ada padanan — sistem lama tidak mencatat pergantian pemegang sama sekali |
| `ID_ROSTER` | `NUMBER(19)` | ya | `FK_PERISTIWA_ROSTER_1` → `ROSTER_JENJANG` | `S-052` | tidak ada padanan — barisnya ditimpa di tempat, tanpa jejak |
| `JENIS_PERISTIWA` | `VARCHAR2(32 CHAR)` | ya | `CK_PERISTIWA_ROSTER_1` — domain tertutup | `S-052` | **delegasi tetap** ditambal satu baris kode yang menukar satu nama dengan nama lain saat pembentukan (`F-19`); jenis lainnya tidak ada padanan |
| `PEMEGANG_SEBELUM` | `VARCHAR2(64 CHAR)` | tidak | — | `S-052` | `KomiteID` sebelum ditimpa — nilai yang di sistem lama hilang begitu ditimpa |
| `PEMEGANG_SESUDAH` | `VARCHAR2(64 CHAR)` | tidak | `CK_PERISTIWA_ROSTER_2` — sekurang-kurangnya satu dari kedua pemegang terisi | `S-052` | `KomiteID` sesudah ditimpa (`F-19`) |
| `PELAKU` | `VARCHAR2(64 CHAR)` | ya | — | `S-052` | tidak ada padanan — `E-3` menuntutnya, sistem lama tidak menyimpannya |
| `ALASAN` | `VARCHAR2(2000 CHAR)` | tidak | — | `S-052` | tidak ada padanan — `E-3` |
| `TERJADI_PADA` | `TIMESTAMP WITH TIME ZONE` | ya | — | `S-052` | tidak ada padanan — `E-3`; zona `Asia/Jakarta` (aturan kerja `G`) |

Kedua kolom pemegang boleh kosong, dan **kosong di sini berarti tidak ada pemegang** —
kedudukan yang belum terisi pada `PEMEGANG_SEBELUM`, kedudukan yang dikosongkan pada
`PEMEGANG_SESUDAH`. Keduanya teks, sehingga tidak ada kolom di objek ini yang kosong dan
nolnya dapat tertukar (`K5-6`, `ADR-0019`). `CK_PERISTIWA_ROSTER_2` menolak baris yang kedua
pemegangnya kosong sekaligus, karena baris semacam itu tidak mencatat pergantian apa pun.

**Objek ini tidak menampung pengalihan sesaat, dan itu disengaja.** `H-4` menetapkan
pengalihan sesaat tetap tidak diekspos: modelnya menampung lewat `E-3`, permukaannya tidak
ada. Kolomnya **tidak kurang** — yang tidak ada adalah jalur yang menulisnya. Dicatat di sini
agar tidak ada yang menambahkannya sebagai kolom yang dikira terlupakan.

**Peristiwa roster tidak ikut dibaca saat sirkulasi dibaca.** Ia hidup di seam `KelolaRoster`,
bukan `BacaSirkulasi`, dan tidak muncul pada model baca sirkulasi mana pun.

### Perubahan pada tabel milik sisi Klaim

Dua, dan keduanya perlu persetujuan pemilik tabelnya.

| Perubahan | Tabel | Bentuk | S-xxx | Sebab |
|---|---|---|---|---|
| **Kolom baru** `MAKSUD_PENGAJUAN` | `KLAIM` | `VARCHAR2(16 CHAR)`, boleh kosong, `CK_KLAIM_<n>` domain tertutup | `S-012` | `IsCloseFile` dan `IsReject` ditulis ke klaim saat pengajuan (`K5-7`). **Tiga keadaan berarti**, bukan dua — tidak diajukan, diajukan untuk ditutup, diajukan untuk ditolak — sehingga `K5-6` justru **melarang** membuatnya boolean |
| **Constraint baru** `UQ_ADJUSTMENT_<n>` | `ADJUSTMENT` | keunikan atas kolom nomor akseptasi yang **sudah ada** | `S-040`, `S-066` | `D-4` menegakkan idempotensi lewat constraint, bukan baca-lalu-tulis. Bentuk kolomnya tidak berubah; yang ditambahkan hanya penegaknya |

Kolom pendaratan yang **sudah ada dan tidak diubah**: keputusan komite, nomor akseptasi,
tanggal akseptasi, penanda bersyarat, catatan bersyarat, dan usul tutup pada tabel usulan.
Domain kolom keputusan di sana sengaja dibiarkan terbuka oleh sisi Klaim karena modul ini
yang menulisnya.

**Yang tidak dimigrasikan sebagai kolom, dan disebut begitu di sini agar tidak dicari:**
`KOMITECOUNT` dan `FLAGONGOINGCOMMITTE` — turunan, dihitung saat dibaca (`D-2`) · penanda
"sedang disirkulasikan" — turunan dari jenjang (`S-034`) · `.KomiteNo` sebagai tali ke
sirkulasi — tali yang sah adalah relasi dari sisi sirkulasi (`ADR-0031`, `S-070`) · penanda
final XOL dan penanda tolak lawas — keduanya sasaran tulis tanpa kolom pendaratan.

---

## Medan layar

Detail atas enam persyaratan ber-seam `Layar` (`S-058` … `S-063`). **Bukan persyaratan
baru**; tiap medan menunjuk `S-xxx` yang menuntutnya.

**Cakupan, dinyatakan lebih dulu.** Satu-satunya masukan komite adalah **keputusan,
komentar, dan penanda usulan**. Seluruh panel nilai adalah proyeksi baca-saja: komite tidak
dapat mengubah satu pun angka uang (`F-9`), dan tidak pernah menyunting rekening penerima
(`F-1`, `F-6`). Enam belas kemunculan medan rekening seluruhnya baca-saja, dan enam belas
properti nilai juga.

**Tiap kunci ditegakkan dua kali.** Kolom *validasi klien* adalah penegakan pertama — untuk
kejelasan. Kolom *validasi server* menunjuk `S-xxx` di seam perintah — untuk kebenaran.
**Kunci yang hanya ada di klien bukan kunci** (`S-063`).

### Medan yang dapat disunting

| Medan | Sumber lama | Kendali | Wajib | Baca-saja bila | Tampil bila | Validasi klien | Validasi server | S-xxx |
|---|---|---|:-:|---|---|---|---|---|
| Keputusan | `.AcceptStatus` | pilihan dua nilai | ya | pemanggil bukan pemegang jenjang aktif | selalu | dua nilai saja; tombol lanjut mati sampai terisi | `S-017` menolak bukan pemegang; `S-031` menolak nilai di luar domain | `S-058` |
| Komentar | `.Comment` | teks panjang | ya | idem | selalu | tidak boleh kosong | `S-061` menolak keputusan tanpa komentar | `S-058`, `S-061` |
| Penanda bersyarat | `.IsSubjectivity` | dua keadaan | tidak | derajat pemegang bukan yang terendah | jenis sirkulasi adalah usulan pembayaran | satu boolean, satu ejaan | `S-024` menolak penyuntingan dari jenjang bukan pertama; `S-026` menolak nilai di luar dua keadaan | `S-058` |
| Catatan syarat | `.SubjectivityNote` | teks panjang | ya bila penanda bersyarat menyala | idem | penanda bersyarat menyala | wajib mengikuti penandanya | `S-024`; `S-038` memakainya sebagai syarat tidak terbitnya nomor | `S-058` |
| Penanda usul tutup | `.IsProposeClose` | dua keadaan | tidak | idem | jenis sirkulasi adalah usulan pembayaran | satu boolean | `S-024`, `S-026` | `S-058` |
| Penanda usul cadang | `.IsPropReserved` | dua keadaan | tidak | idem | idem | satu boolean | `S-024`, `S-026` | `S-058` |
| Penanda XOL | `.CNPFlagXOL` | dua keadaan | tidak | idem — **dikunci `J-2`** | idem | satu boolean | `S-060` menolak penyuntingan dari jenjang lain; masuk versi usulan lewat `S-025` | `S-060` |

Ketujuhnya adalah **penanda usulan** kecuali dua yang pertama, yang adalah keputusan dan
komentar. Tidak ada medan kedelapan: `.CurrencyID`, yang `F-8` temukan dapat disunting pada
satu render, **dijadikan baca-saja seluruhnya** oleh `J-2` — mengganti label mata uang
sementara nilainya beku bukan kemampuan.

### Medan dan blok yang baca-saja

| Kelompok | Sumber lama | Cacah | Baca-saja bagi | Tampil bila | S-xxx |
|---|---|---:|---|---|---|
| Rekening penerima | `PayableTo`, `NameOfBank`, `BranchOfBank`, `NoAccount`, `SwiftCode`, `Currency` | 16 kemunculan | **semua jenjang** (`F-1`, `F-6`) | selalu; baris Swift hanya bila terisi | `S-059` |
| Panel nilai | properti nilai pada berkas anak | 16 properti | **semua jenjang** (`F-9`) | selalu | `S-059` |
| Label mata uang | `.CurrencyID` | 1 render | **semua jenjang** (`J-2`) | selalu | `S-059`, `S-060` |
| Daftar jenjang | `pyWorkPage.KomiteList` | 1 daftar | semua | selalu | `S-047` |
| Alokasi per layer dan treaty | `.Adjustment.SpreadingRisk` | 1 daftar | semua | selalu | `S-059` |
| Penyebaran keluar | `.Adjustment.SpreadingAdjustment` | 1 daftar | semua | selalu | `S-059` |
| Penyebaran porsi tetap | `.Adjustment.SpreadingQuotaShare` | 1 daftar | semua | selalu | `S-059` |
| Pembagian kerugian | `.Adjustment.CNPSpreadLoss` | 1 daftar | semua | selalu | `S-059` |
| Nilai klaim per mata uang | `.Adjustment.ListClaimAcceptation` | 1 daftar | semua | selalu | `S-059` |
| Alokasi yang sudah pernah dibayar | `.Adjustment.AlokasiXOLPaid` | 1 daftar | semua | daftar itu tidak kosong | `S-059` |
| Lima medan induk | medan induk `pyReadOnly` **dan** `pyRequired` | 5 | semua | selalu | `S-061` |

**Lima medan induk tidak memblokir apa pun.** Aturan wajibnya tidak pernah menyala di sistem
lama karena medannya baca-saja; medan wajib yang pelakunya tidak dapat mengisi adalah
jebakan, bukan validasi (`E-4`, `F-7`). Ini **paritas**, bukan deviasi.

### Blok yang dibuang, dan sebabnya

Dicatat, bukan dihilangkan diam-diam (`S-062`).

| Blok | Syarat lamanya | Sebab dibuang |
|---|---|---|
| Blok bersyarat `1=2` | perbandingan yang tidak pernah dapat benar | tidak pernah tampil; membawanya membawa kembali pertanyaan tentang kapan ia menyala |
| Blok bersyarat `1=3` | idem | idem |
| Blok bersyarat `NEVER` | dinyatakan tidak pernah | idem |
| Medan bersufiks dua | pasangan tampilan atas data hulu | nol pembaca di Activity modul ini (`F-3`) |
| Kunci dari penanda luar | `pyWorkCover.IsOutstanding`, `.CNPFlagOuts` | sasarannya tidak dapat ditentukan dari berkas — kedekatan posisi bukan bukti (`F-12`), dan `J-2` melarutkan `F-12` |

### Yang belum tertulis di tabel ini

Dua medan dari enam yang `F-2` cacah sebagai dapat-disunting **tidak dinamai** di satu pun
berkas yang dipegang: `PENGETAHUAN.md` §10.1 mendaftarkan medan rekening sebagai dapat
disunting, dan `F-1` serta `F-6` membatalkan pendaftaran itu. Nama kedua medan sisanya
terbaca dari `Section/ShowTransfer.xml`, **yang ada di repo** — karena itu ini **bukan
pagar** dan bukan lubang inventaris; ia pembacaan yang belum dilakukan, dan pagar yang
menunjuk sesuatu di repo gugur dengan sendirinya. Tujuh medan di tabel di atas seluruhnya
berpenulis; yang belum diperiksa adalah apakah daftarnya sudah lengkap.

---

## Testing Decisions

### Apa yang membuat sebuah uji baik di sini

Uji memeriksa **perilaku yang terlihat dari luar seam**: apa yang dikembalikan perintah, apa
yang berubah di basis data, dan galat apa yang keluar. Uji **tidak** memeriksa nama fungsi
internal, urutan pemanggilan, atau bentuk struktur dalam memori — seluruh temuan yang
melahirkan modul ini justru lahir karena sistem lama menegakkan aturan lewat bentuk internal
yang tidak ada yang dapat memeriksanya dari luar.

Tiga sifat tambahan yang khusus untuk modul ini:

- **Jalur gagal adalah spesifikasi, bukan catatan kaki.** Tiap perintah punya uji untuk
  keadaan gagalnya, dan uji itu memeriksa dua hal: galat yang benar keluar, **dan** keadaan
  yang seharusnya tidak berubah memang tidak berubah.
- **Uji deviasi dirancang gagal terhadap sistem lama.** Uji semacam itu ada bukan untuk
  dilewati; ia ada untuk membuktikan perbedaannya memang dimaksudkan. Uji deviasi yang
  **lulus** paritas adalah tanda deviasinya tidak terpasang.
- **Keadaan masukan disebut, bukan hanya perbedaannya.** Sebuah uji yang tidak menyebut
  kombinasi masukan yang membuatnya gagal tidak dapat dijalankan orang lain.

### Yang diuji

| Yang diuji | Lewat seam | Persyaratan | Catatan |
|---|---|---|---|
| Pembentukan sirkulasi, seluruh jenisnya | `BentukSirkulasi` | S-001 … S-008 | termasuk galat konfigurasi seleksi dan ketiadaan pemegang peran |
| Kegagalan pembentukan | `BentukSirkulasi` | S-009 … S-016 | memeriksa klaim tidak berubah selain maksud |
| Urutan giliran dan gerbang wewenang | `CatatKeputusan` | S-017, S-018 | pemegang bukan jenjang aktif ditolak |
| Seluruh jalur mesin keadaan dan **delapan** invarian perilaku | kelima operasi | S-018 … S-023 | tabel invarian di atas menyebut cara membuktikan tiap pelanggaran. `I-2`, `I-8`, `I-10` **tidak** di sini — lihat baris uji DDL |
| Fungsi akibat, ketujuh barisnya | `CatatKeputusan` + `BacaSirkulasi` | S-031, S-046 | tiga nilai tampil diperiksa dari model baca, bukan dari kode |
| Penerbitan nomor akseptasi dan kegagalannya | `CatatKeputusan` | S-038 … S-043 | idempotensi diuji dengan mengulang perintah ber-kunci sama |
| Versi usulan dan penjaga di batas tulis | `CatatKeputusan` | S-024, S-025 | penyuntingan oleh jenjang bukan pertama ditolak |
| Idempotensi kedua perintah | `BentukSirkulasi` + `CatatKeputusan` | S-013, S-027 | percobaan ulang mengembalikan hasil yang sama |
| Umur sebagai turunan | `BacaSirkulasi` | S-045 | memeriksa tidak ada kolom umur yang tersimpan |
| Hak baca diturunkan | `BacaSirkulasi` | S-044 | pembaca yang boleh membuka klaim boleh membaca sirkulasinya |
| Pengelolaan roster dan delegasi tetap | `KelolaRoster` | S-050 … S-054 | pergantian pemegang diperiksa sebagai peristiwa, bukan sebagai penimpaan |
| Pengelolaan tabel seleksi dan cakupannya | `KelolaTabelSeleksi` | S-055 … S-057 | penyimpanan yang meninggalkan kombinasi tak tercakup ditolak |
| Penegakan kunci sisi klien | `Layar` | S-058 … S-063 | **seam sempit**: memeriksa penegakan pertama **ada**, bukan bahwa aturannya benar |
| Bentuk skema — `I-2`, `I-8`, `I-10`, keunikan, penamaan | `DDL` | S-064 … S-070 | pemeriksaan katalog skema dan constraint, bukan pemanggilan operasi |

Port hilir diuji sebagai **port kosong**: uji memeriksa bahwa port dipanggil dengan niat
yang tercatat lebih dulu dan kunci idempotensi, bukan memeriksa apa yang dikirimnya — isinya
berpagar (`S-037`).

**Dua bentuk uji yang bukan uji perilaku, dan sebabnya.** Uji DDL memeriksa katalog skema
dan constraint; ia dipilih karena keadaan terlarang yang **tidak dapat ditulis** lebih kuat
daripada keadaan terlarang yang ditolak kode, dan sejalan `D-4` yang memilih constraint di
atas baca-lalu-tulis. Uji layar memeriksa bahwa penegakan sisi klien **ada**; kebenaran
aturannya diuji sekali di seam perintah, tidak dua kali. Keduanya dinyatakan agar tidak
dikira uji perilaku yang tertinggal.

### Prior art

Register deviasi modul ini sudah memuat dua puluh dua baris, masing-masing dengan uji
paritas yang dirancang gagal dan keadaan masukan yang membuatnya gagal. Uji di atas mengacu
padanya dengan nomor, bukan menyalinnya. Pemetaan tiap nomor deviasi ke persyaratan yang
mewujudkannya ada di bagian **Deviasi**.

Dua puluh dari dua puluh dua uji itu berstatus **belum diratifikasi** — disusun juru catat
dari bunyi deviasinya, bukan disalin dari sumber keputusannya. Status itu tidak berubah oleh
spesifikasi ini, dan uji yang belum diratifikasi tetap dijalankan sambil tetap bertanda
demikian.

Skenario paritas ronde 5 (`P5-01` … `P5-16`) menyediakan bentuk yang sudah dipakai: skenario,
masukan, keluaran yang diamati, tempat memeriksanya, dan penanda paritas atau deviasi.

Sisi Claim menyediakan prior art untuk uji tingkat skema: berkas DDL-nya memasang constraint
sebagai penegak aturan, bukan sebagai kerapian, dan alat pemeriksa penamaannya keluar dengan
kode gagal ketika sebuah pengenal melewati batas. Modul ini memakai pola yang sama.

---

## Out of Scope

**Empat aliran berpagar** — A-3, A-4, A-5, A-6 — lewat delapan pagar. A-3 **terbatas**,
bukan beku: mekanika penomoran ditulis penuh dan hanya nilai ambang periode buku yang
ditahan. A-4, A-5, dan isi A-6 beku seluruhnya.

Pagar bukan penundaan bersyarat: tidak ada pendekatan sementara, dan "sementara pakai
pendekatan X" adalah pelanggaran karena pendekatan sementara mengeras jadi keputusan.

Kolom **Titik sentuh** menyebut persyaratan mana yang berhenti di pagar itu.

| # | Yang tidak ditulis | Aliran | Register | Baris inventaris — bagian dan **nama objek** | Titik sentuh | Sementara ini |
|---|---|---|---|---|---|---|
| PAGAR-01 | Nilai ambang periode buku | A-3 | `PG-03` | §2.5 baris 4 — isi badan **`PROC_GENERATE_SEQUENCE_NUMBER`** pada basis data berjalan belum pernah dibaca; yang dipegang hasil reverse-engineer | `S-043` | Parameter bernama tanpa nilai; mekanika penomoran ditulis penuh |
| PAGAR-02 | Isi baris akseptasi dan bentuk simpannya | A-5 | `PG-05` | §2.1 — **`PEGA_JSON_OS_AKSEP_KLAIM`**, **`PEGA_JSON_OS_AKSEP_SUBJECTIVITY`**, **`HISTORYAKSEPTASIPEGA`**: DDL ketiganya tidak dipegang | `S-037` — `PORT_PENCATATAN_AKSEPTASI` | `PORT_PENCATATAN_AKSEPTASI` kosong |
| PAGAR-03 | Cakupan dan bentuk baris pembalikan | A-5 | `PG-05` | §2.1 — **`PEGA_JSON_OS_AKSEP_KLAIM`** tidak dipegang; §2.5 baris 2 — nol baris **`OS_AKSEPTASI_KLAIM.DATA_JSON`** pernah diambil | `S-037` — `PORT_PEMBALIKAN_AKSEPTASI` | `PORT_PEMBALIKAN_AKSEPTASI` kosong |
| PAGAR-04 | Rincian layer untuk usulan bersyarat | A-5 | `PG-05` | §2.1 — **`XOL2_AKSEP_KLAIM`** tidak dipegang; §2.3 — **`SetProtectionEstimation`** dipanggil, berkasnya tidak ada | `S-037` | Tidak ada persyaratan; lubang dinyatakan |
| PAGAR-05 | Muatan kiriman pembayaran dan jumlah kiriman per akseptasi | A-4 | `PG-04` | §2.5 baris 1 — nol baris **`POOLDATA.DIRECTTOKASIR_LOG`**; DDL-nya dipegang, isinya tidak | `S-037` — `PORT_PENGIRIMAN_KASIR` | `PORT_PENGIRIMAN_KASIR` kosong |
| PAGAR-06 | Penjaga kiriman ganda | A-4 | `PG-04` | §2.5 baris 1 — nol baris **`POOLDATA.DIRECTTOKASIR_LOG`**; berapa kiriman terjadi per nomor akseptasi dijawab `SELECT` yang sama | `S-037` | Niat dicatat sebelum panggilan dengan kunci idempotensi per efek; penjaganya tidak dirancang |
| PAGAR-07 | Isi surat dan dokumen | A-6 | `PG-06` | §2.3 — **`PostEmailKomiteCNP`** dipanggil `KomitePostAdjustment` 19.1, berkasnya tidak ada | `S-037` — `PORT_PENERBITAN_SURAT` | `PORT_PENERBITAN_SURAT` kosong. **Kapan** surat terbit dan **siapa** penerimanya ditulis; isinya tidak |
| PAGAR-08 | Lampiran dan kanal pengiriman | A-6 | `PG-06` | §2.3 — **`SendEmailWithAttachments`** dipanggil lima kali, berkasnya tidak ada | `S-037` — `PORT_PENGUNGGAHAN_DOKUMEN` | `PORT_PENGUNGGAHAN_DOKUMEN` kosong |

Selain kedelapan pagar itu, di luar cakupan:

- **Pengalihan sesaat.** Modelnya menampung, permukaannya tidak ada (`H-4`). Peran yang
  berwenang mengalihkan adalah keputusan pemilik proses yang belum diambil.
- **Jalur khusus satu operator pembuat case.** Tidak dimigrasikan (`H-7`).
- **Masa tenggat dan eskalasi.** Tidak ada keadaan kedaluwarsa dan tidak ada pekerja latar
  (`E-5`). Umur dikumpulkan sebagai data agar aturan tenggat kelak diputuskan dengan bukti.
- **Perhitungan nilai, kurs, rumus pemulihan, dan alokasi layer.** Milik modul Klaim.
- **Migrasi data.** Bentuknya bergantung pada satu cacah yang belum diambil: berapa klaim di
  produksi membawa pengenal sirkulasi yang bukan miliknya (§2.5 baris 6). Yang sudah
  diputuskan: tali klaim–sirkulasi dibaca dari sisi sirkulasi, dan pengenal yang tersimpan di
  klaim lama tidak dipercaya.

---

## Ketertelusuran

Rantai **rule lama → persyaratan → uji**. Tanpa peta ini shadow-run kehilangan dasar
pembandingnya, dan tiket tidak dapat menyebut apa yang digantikannya.

Tiap rule lama punya baris. Rule yang jatuh ke aliran berpagar diberi nomor pagar, **bukan
dikosongkan** — baris kosong terbaca sebagai terlupakan.

Keadaan: **DIGANTIKAN** · **DIBUANG** dengan sebabnya · **DIPAGARI** dengan nomor pagarnya.

### Inti — sebelas rule yang menjalankan modul

| Rule lama (berkas·langkah) | S-xxx | P-xx | Keadaan |
|---|---|---|---|
| `KomiteTreaty_Flow` — mesin keadaan sirkulasi | S-018, S-022, S-023 | P5-01, P5-02 | **DIGANTIKAN** — mesin keadaan pada *Implementation Decisions* |
| `KomiteTreaty_Flow` — `pyRouteTo=Custom`, nol SLA dan nol eskalasi (F-4) | S-045 | — | **DIBUANG** — tidak ada yang menggantikannya; umur menjadi turunan (`E-5`) |
| `IsKomiteLoop` — gerbang `KomiteLoop` | S-018 | P5-02b | **DIBUANG** — `K5-1`; hitungan jenjang berhenti menjadi penentu |
| `IsKomiteLoop` — ticket `komiteAccept_ticket` | — | — | **DIBUANG** — penaiknya tidak ada di antara 338 berkas; pemiliknya di luar ekspor (`KETETAPAN.md` §10.1) |
| `KomiteRouter`·1–4 — penugasan lewat `.KomiteCount` | S-018 | P5-02b | **DIBUANG** — `K5-1`; penentu kedua dihapus |
| `KomiteRouter`·6, 6.1 — baris pertama yang belum memutuskan | S-018 | P5-01, P5-02 | **DIGANTIKAN** — paritas; urutan derajat menaik dipertahankan |
| `KomiteRouter` — empat cabang mati (K-02) | — | — | **DIBUANG** — cabang tidak pernah menyala |
| `CreateChildKomiteCNP_Act`·9, 26.13–26.15 — jumlah jenjang dari cacah baris laporan | S-005, S-067 | P5-01, P5-03 | **DIGANTIKAN** — jumlah menjadi akibat isi roster (F-13) |
| `CreateChildKomiteCNP_Act`·10 — bagian treaty ke pembanding | S-002 | P5-09 | **DIGANTIKAN** — bagian treaty menjadi masukan tabel seleksi |
| `CreateChildKomiteCNP_Act`·11 — nilai baris terakhir sebagai penentu | S-004 | P5-13 | **DIGANTIKAN** — nilai usulan yang diajukan (F-17) |
| `CreateChildKomiteCNP_Act`·12, 13 — dua tetapan kelas di dalam langkah | S-002, S-055, S-057 | P5-10 | **DIGANTIKAN** — aturan seleksi sebagai data |
| `CreateChildKomiteCNP_Act`·14 — cabang yang syaratnya tidak pernah dapat benar | S-003 | P5-09 | **DIBUANG** — cabang mati (F-16); perintahnya diserap `K5-5` lewat deviasi 18 |
| `CreateChildKomiteCNP_Act`·14 — jalur yang meninggalkan ambang tak terisi | S-003, S-007 | P5-11 | **DIGANTIKAN** — galat konfigurasi, bukan roster kosong |
| `CreateChildKomiteCNP_Act`·26.3–26.6 — pembanding ambang yang dimatikan | S-057 | P5-12 | **DIGANTIKAN** — bersyarat menjadi masukan ketiga (`K6-3`) |
| `CreateChildKomiteCNP_Act`·26.8, 26.13 — `.JABATAN` ditimpa daftar berisi `.ID` | S-021 | P5-14 | **DIGANTIKAN** — jabatan disalin dari roster (F-18) |
| `CreateChildKomiteCNP_Act`·26.8.3 — satu substitusi orang ke orang | S-051, S-052, S-036 | P5-15 | **DIGANTIKAN** — cabangnya dibuang oleh keputusan beku no. 5, tetapi kebutuhan yang ditambalnya digantikan: delegasi tetap sebagai atribut roster (`S-051`, `H-3`), dan pergantiannya tercatat sebagai peristiwa pada `PERISTIWA_ROSTER` (`S-052`, `E-3`) |
| `CreateChildKomiteCNP_Act`·26.9, 26.10, 26.15 — perlakuan khusus satu operator | S-054 | — | **DIBUANG** — `H-7` |
| `CreateChildKomiteCNP_Act`·28 — pembatalan hanya pada alokasi layer kosong | S-007, S-010 | P5-08 | **DIGANTIKAN** — setiap jalur keluar berakhir dengan kegagalan yang terlihat (F-22) |
| `CreateChildKomiteCNP_Act`·29 — penjaga `@hasMessages` melewati satu langkah | S-009, S-015, S-070 | P5-07 | **DIGANTIKAN** — satu transaksi (`K5-2`) |
| `CreateChildKomiteCloseNP_Act`·6 — penanda maksud ditulis ke induk saat pembuatan | S-011, S-012 | P5-07, P5-08 | **DIGANTIKAN** — maksud pada kolom tersendiri (`K5-7`) |
| `CreateChildKomiteCloseNP_Act`·8, 9 — satu jenjang, pemegang dari nama di dalam rule | S-008 | P5-04 | **DIGANTIKAN** — pemegang diselesaikan dari roster berdasarkan peran (`K6-4`) |
| `CreateChildKomiteCloseNP_Act`·12 — jenis dibedakan dari kolom analisis | S-001 | P5-04, P5-05 | **DIGANTIKAN** — jenis dinyatakan pemanggil (`K6-4`, F-20) |
| `FilterEmailKomiteWithLimit` — penyaring roster, urut derajat menaik | S-005, S-055 | P5-10 | **DIGANTIKAN** — pencarian ke tabel seleksi lalu ke roster |
| `FilterEmailKomiteWithLimit` — `.LIMIT_TOP` diambil tetapi tidak pernah menyaring | — | — | **DIBUANG** — tidak pernah berpengaruh (F-14) |
| `ShowTransfer` — panel nilai dan rekening, seluruhnya baca-saja | S-059 | — | **DIGANTIKAN** — layar Persetujuan Komite (F-1, F-6, F-9) |
| `ShowTransfer` — enam medan dapat-sunting, empat berkunci jenjang pertama | S-058, S-060 | — | **DIGANTIKAN** — tiga medan masukan; kunci ditegakkan dua kali (F-2, F-8) |
| `ShowTransfer` — lima medan induk baca-saja dan wajib sekaligus | S-061 | — | **DIGANTIKAN** — tidak memblokir keputusan (`E-4`, F-7) |
| `ShowTransfer` — blok bersyarat yang tidak pernah menyala | S-062 | — | **DIBUANG** — syaratnya tidak pernah dapat benar |
| `ShowTransfer` — medan bersufiks dua, tampilan atas data hulu | — | — | **DIBUANG** — nol pembaca di Activity modul ini (F-3) |
| `KomitePostAdjustment`·14.12 — penanda usulan tutup mendarat di induk pada persetujuan terakhir non-bersyarat | S-031, S-038 | — | **DIGANTIKAN** — fungsi akibat (F-5) |
| `KomitePostAdjustment`·16 — nomor urut sirkulasi atas satu usulan | S-014, S-065 | — | **DIGANTIKAN** — constraint keunikan (F-11) |
| `KomitePostAdjustment`·20 — penolakan satu usulan menolak seluruh klaim | S-031 | — | **DIGANTIKAN** — baris ketiga tabel fungsi akibat; ditandai `D-1` sebagai paritas-dengan-pertanyaan-bisnis, bukan cacat |
| `KomitePostAdjustment` — jalur A akibat pada klaim, tersebar | S-031, S-033 | — | **DIGANTIKAN** — fungsi akibat di satu tempat (`D-3`) |
| `KomitePostAdjustment`·19.1 → `PostEmailKomiteCNP` | S-037 | — | **DIPAGARI PAGAR-07** — `PORT_PENERBITAN_SURAT` kosong |
| `KomitePostAdjustment`·14.20 → `GetBase64Attachment` | — | — | **DIBUANG** — rule bawaan, bukan kandidat penulisan ulang |
| `KomitePostAdjustmentCWP` — jalur B akibat pada klaim | S-031, S-032 | P5-04, P5-05 | **DIGANTIKAN** — fungsi akibat; keadaan sirkulasi terpisah dari akibat (keputusan beku no. 9) |
| `KomitePostAdjustmentCWP` → `insertClaimReject_NP`, `insertClaimFinalOrClosed_NP` | S-037 | — | **DIPAGARI PAGAR-07** — efek ke luar, isinya berpagar |
| `InsertXOLKlaimCNP` → `SaveXOLClaim_SQL` | S-037 | — | **DIPAGARI PAGAR-04** — sasaran rincian layer `XOL2_AKSEP_KLAIM` tidak dipegang |
| `HitServiceToKasirKMT_Act`·10.3, 10.7 — kiriman ke Kasir di dalam loop | S-037 | — | **DIPAGARI PAGAR-05** — `PORT_PENGIRIMAN_KASIR` kosong |
| `HitServiceToKasirKMT_Act`·10.9 — penulisan log kiriman tanpa pra-syarat | S-037 | — | **DIPAGARI PAGAR-06** — jumlah kiriman per akseptasi berpagar |

### Sekitarnya — rule yang dipanggil dari jalur di atas

| Rule lama (berkas·langkah) | S-xxx | P-xx | Keadaan |
|---|---|---|---|
| `ViewTransferDtl` — Flow Action pembuka layar | S-058 | — | **DIGANTIKAN** — layar Persetujuan Komite |
| `SetKomiteList_Act` — pemetaan nama orang ke jabatan | S-021, S-036 | P5-14, P5-15 | **DIGANTIKAN** — roster sebagai data (K-03, F-18) |
| `ShowDetailXOL` + `ReinstatementPremiumDetails` | S-059 | — | **DIGANTIKAN** — proyeksi baca-saja; rumusnya milik modul Klaim |
| `ResetSubjectivityNote` | S-026, S-030 | P5-05 | **DIGANTIKAN** — penanda dua-keadaan satu boolean (`K5-6`) |
| `ProteksiSendKomiteCNP_Act` — gerbang yang menandai, tidak menghentikan | S-015 | — | **DIGANTIKAN** — validasi menolak (`H-5`, F-21) |
| `ProteksiSendKomiteCNP_Act`·21 → `SetProtectionEstimation` | — | — | **DIPAGARI PAGAR-04** — berkasnya tidak ada (`INVENTARIS-BUKTI.md` §2.3) |
| `GetSequenceNumber_SQL` → `PROC_GENERATE_SEQUENCE_NUMBER` | S-039, S-040, S-041, S-066 | — | **DIGANTIKAN** — pembungkus pencatat; prosedur lama dipakai apa adanya (`D-4`) |
| `PROC_GENERATE_SEQUENCE_NUMBER` — penentuan periode buku dari `TANGGAL_CLOSING` | S-043 | — | **DIPAGARI PAGAR-01** — nilai ambangnya tidak ditulis |
| `GetKodeProdNonLife_SQL` — kode produksi | S-043 | — | **DIPAGARI PAGAR-01** — bersandar pada periode buku yang sama |
| `InsertHistoryAkseptasiPega_Sql` → `HISTORYAKSEPTASIPEGA` | S-037 | — | **DIPAGARI PAGAR-02** — DDL objeknya tidak dipegang (§2.1) |
| `InsertOSKlaimCNP` → `SaveOSClaim_SQL` | S-037 | — | **DIPAGARI PAGAR-02** — `PORT_PENCATATAN_AKSEPTASI` kosong |
| `InsertOSSubjectivityCNP` → `SaveOSSubjectivity_SQL` | S-037 | — | **DIPAGARI PAGAR-02** — jalur bersyarat, objeknya tidak dipegang |
| `SaveRejectOSKomiteCNP` → `GetDataCNPOS`, `SaveDataToOsAkseptasiNP` | S-037 | — | **DIPAGARI PAGAR-03** — `PORT_PEMBALIKAN_AKSEPTASI` kosong |
| `GenerateAccCNP_act` → `InsertDocument_Act` → `InsertGoogleStorage_Act` | S-037 | — | **DIPAGARI PAGAR-08** — `PORT_PENGUNGGAHAN_DOKUMEN` kosong |
| `SendEmailKlaim_KMT`, `SendEmailKlaimRejectClose_KMT` → `GetEmailUser_RD` | S-037 | — | **DIPAGARI PAGAR-07** — kapan terbit dan siapa penerimanya ditulis; isinya tidak |
| `SendAcceptationToKasir`, `SendErrorDirectKasir`, `InsertLOGDirectKasir_SQL` | S-037 | — | **DIPAGARI PAGAR-05, PAGAR-06** |
| `getStatusKonversi_Act` → `getStatusKonversi_SQL` | S-037 | — | **DIPAGARI PAGAR-05** — dipanggil dari dalam jalur kiriman Kasir |
| `InsertChronology_DT` | S-033 | — | **DIGANTIKAN** — tabel kronologi yang sudah ada, tanpa pengecualian per peran |
| `IsPEGAPROD`, `IsPEGASyariah`, `GetLinkService`, `LinkService` | S-036 | — | **DIGANTIKAN** — konfigurasi, bukan cabang di dalam aturan (keputusan beku no. 5) |
| `BrowseCurrency_RD`, `BrowseBankGroup`, `GetEmailUser_RD` | S-059 | — | **DIGANTIKAN** — endpoint referensi baca-saja |
| `GetMimeType` | — | — | **DIPAGARI PAGAR-08** — milik layanan dokumen |

### Dua rule yang keadaannya menunggu keputusan

Diperiksa 2026-09-21, dan hasilnya membalik dugaan sebelumnya: **berkas keempat rule ini ada
di repo**. Keduanya karena itu **bukan** lubang bukti dan **tidak boleh dipagari** — pagar
yang menunjuk sesuatu di repo gugur dengan sendirinya (aturan kerja `B`). Lihat
`INVENTARIS-BUKTI.md` §2.4 dan §7.1.

| Rule lama (berkas·langkah) | S-xxx | P-xx | Keadaan |
|---|---|---|---|
| `KonversiKlaim_Act` → `KonversiKlaimNonLife` — dipanggil dari `KomitePostAdjustment`, nomor langkah tidak terbaca | — | — | **BELUM DIPUTUSKAN** — efek ke luar yang tidak terpeta ke satu pun dari kelima port hilir. Menunggu **keputusan**, bukan bukti |
| `InsertJsonClaimTreatyNonProp_act` → `InsertClaimPNC` — dipanggil dua kali dari jalur keputusan, nomor langkah tidak terbaca | — | — | **BELUM DIPUTUSKAN** — nasibnya bergantung pada apakah JSON klaim adalah artefak Pega (`PENGETAHUAN.md` §13, "perlu keputusan"). Menunggu **keputusan**, bukan bukti |

Keadaan keempat ini dipakai dengan sengaja, dan bukan keadaan kosong: ketiga keadaan yang
lain menyatakan apa yang **sudah** diputuskan, sementara kedua baris ini menyatakan bahwa
keputusannya milik pemilik proses dan belum diambil. Keduanya tercatat pada daftar gagal
`T-1` di akhir bagian *Persyaratan*, dan tetap menjadi hal yang harus tertutup sebelum
spesifikasi ini dipakai sebagai rujukan tetap.

Nomor langkah pemanggilnya **tidak disebut** karena tidak terbaca: pohon panggilan
`PENGETAHUAN.md` §2.3 bernomor hierarkis penelusuran dan, menurut `INVENTARIS-BUKTI.md` §1,
"tidak menyatakan langkah mana memanggil rule mana". Menyebut nomor langkah dari nomor pohon
akan menaikkan indeks buatan manusia menjadi bukti langkah.

---

## Deviasi

Dua puluh lima deviasi bernomor ditarik ke persyaratan yang mewujudkannya. **Bunyinya tidak
disalin** — dua salinan akan perlahan berbeda. Baris **1 – 13** hidup di `KETETAPAN.md`
bagian 9; baris **14 – 22** di `REGISTER-DEVIASI.md` bagian 2 dan 5; baris **23 – 25** di
`REGISTER-DEVIASI.md` bagian 6, diberi nomor 2026-09-21.

Uji deviasi **dirancang gagal** terhadap sistem lama. Uji semacam itu yang **lulus** paritas
adalah tanda deviasinya tidak terpasang.

| # deviasi | S-xxx yang mewujudkannya | P-xx | Sifat uji | Ratifikasi |
|---:|---|---|---|---|
| 1 | S-017 | BARU | DEVIASI DIHARAPKAN | belum diratifikasi |
| 2 | S-021 | P5-14 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 3 | S-022 | BARU | DEVIASI DIHARAPKAN | belum diratifikasi |
| 4 | S-035 | BARU | DEVIASI DIHARAPKAN | belum diratifikasi |
| 5 | S-020, S-068 | BARU — uji DDL | DEVIASI DIHARAPKAN | belum diratifikasi |
| 6 | S-007 | P5-09, P5-11 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 7 | S-025, S-060, S-063 | BARU | DEVIASI DIHARAPKAN | belum diratifikasi |
| 8 | S-004 | P5-13 | DEVIASI DIHARAPKAN | **diratifikasi** — satu-satunya uji yang pernah dirancang eksplisit |
| 9 | S-021 | P5-14 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 10 | S-054, S-036 | BARU | DEVIASI DIHARAPKAN | belum diratifikasi |
| 11 | — | — | **GUGUR** — dilebur ke 18 | tidak berlaku |
| 12 | S-003 | P5-11 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 13 | S-015 | BARU | DEVIASI DIHARAPKAN | belum diratifikasi |
| 14 | S-018, S-023, S-067 | P5-02b | DEVIASI DIHARAPKAN | belum diratifikasi |
| 15 | S-009, S-028, S-070 | P5-07 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 16 | S-010 | P5-08 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 17 | S-006 | P5-06 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 18 | S-002, S-003, S-055, S-056 | P5-09, P5-11 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 19 | S-026 | P5-05 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 20 | S-021, S-051 | P5-15 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 21 | S-016, S-044, S-053 | P5-16 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 22 | S-011, S-012, S-049 | P5-07, P5-08 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 23 | S-001 | PS-01 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 24 | S-008 | PS-02 | DEVIASI DIHARAPKAN | belum diratifikasi |
| 25 | S-008 | PS-03 | DEVIASI DIHARAPKAN | belum diratifikasi |

**Nol deviasi yatim.** Setiap baris yang masih berlaku punya sekurang-kurangnya satu
persyaratan yang mewujudkannya; baris 11 gugur karena dilebur, bukan karena tidak ada yang
mengerjakannya.

**Dua puluh tiga dari dua puluh lima tetap bertanda belum diratifikasi.** Hanya baris 8 yang
pernah dirancang eksplisit; baris 11 gugur. Statusnya tidak berubah oleh bagian ini, dan uji
yang belum diratifikasi tetap dijalankan sambil tetap bertanda demikian.

**Tiga perubahan perilaku yang dulu tanpa nomor kini bernomor.** Ketiganya lahir dari `K6-4`,
yang ditetapkan sesudah pemutakhiran terakhir register: jenis sirkulasi dinyatakan pemanggil
(`S-001` → **deviasi 23**), pemegang jalur tutup dan tolak diselesaikan dari roster
berdasarkan peran (`S-008` → **deviasi 24**), dan sirkulasi ditolak bila tidak ada pemegang
peran itu (`S-008` → **deviasi 25**). Ketiganya kini punya baris di `REGISTER-DEVIASI.md`
bagian 6, masing-masing dengan uji yang dirancang gagal dan keadaan masukan yang membuatnya
gagal — ditulis **di register**, bukan di sini.

Deret `PS-xx` dibuka bersamanya: `P` paritas, `S` tahap spesifikasi. Ia **bukan** ronde tujuh;
grilling tetap tertutup. Deret tersendiri dipakai agar tidak mengulang tabrakan
`G-01 … G-20` melawan `G-1 … G-4`.


---

## Further Notes

**Yang paling mudah salah dibaca di modul ini.** Pada usulan menutup atau menolak klaim,
"komite menyetujui" berarti **klaim ditutup atau ditolak** — karena yang disetujui adalah
usulan penutupan. Salah membaca ini membalik seluruh logika. Itu sebabnya hasil dan akibat
dipisah tegas, fungsi akibat hidup di satu tempat, dan laporan tidak pernah membaca hasil
telanjang.

**Sebelas dari perbedaan perilaku di sini adalah perbaikan, bukan pilihan.** `D-1`
menetapkan pembagiannya: angka yang sudah terbit dipertahankan mutlak; kebenaran catatan dan
wewenang diperbaiki dan paritas ditolak. Di mana grilling menghasilkan cacat berbukti
langkah, perbaikan menjadi default dan paritas yang harus dibela.

**Satu ketetapan lama tidak mengikat.** `H-2` berstatus `BUNYI HILANG`; jalur menutup dan
menolak klaim diatur `K6-4`, yang menggantikannya. Ketetapan yang bunyinya hilang tidak
mengikat siapa pun.

**Satu dasar faktual masih menunggu.** `C-01`: berkas DDL yang dipegang adalah hasil
reverse-engineer, bukan ekspor langsung. Lima penutupan yang bersandar padanya — termasuk
mekanika periode buku — bertanda `TAFSIR (menunggu C-01)` sampai isi badan prosedur penerbit
dibandingkan terhadap basis data berjalan.

**Cerita 42 dicabut, bukan didamaikan.** Cerita itu meminta sirkulasi lama terbaca di sistem
baru dengan riwayat keputusannya. Migrasi data berada di *Out of Scope*, dan bentuknya
bergantung pada satu cacah yang belum diambil (`INVENTARIS-BUKTI.md` §2.5 baris 6). Ketika
sebuah cerita bertentangan dengan pagar, **pagar menang** — mendamaikan keduanya berarti
menuliskan persyaratan atas bukti yang tidak dipegang. Nomor 42 tidak dipakai ulang, dan
barisnya tetap ada agar pencabutannya terbaca. Yang **tetap** berlaku dari wilayah itu
adalah `S-070`: tali klaim–sirkulasi dibaca dari sisi sirkulasi.

**Empat keputusan yang tidak berasal dari permintaan siapa pun.** Audit mencatat empat
keputusan tanpa cerita. Satu di antaranya ternyata punya cerita; tiga sisanya dinyatakan di
sini sebagai keputusan teknis, agar tidak menggantung:

1. **Blok layar yang syaratnya tidak pernah menyala dibuang** (`S-062`) — keputusan teknis.
   Tidak ada pemakai yang meminta, karena tidak ada pemakai yang pernah melihatnya. Yang
   diminta hanyalah pembuangannya tercatat.
2. **Niat efek hilir dicatat sebelum panggilan, lewat port yang dinamai** (`S-037`) —
   keputusan teknis. Isinya berpagar, sehingga tidak ada cerita yang dapat menyebut apa yang
   diuntungkan tanpa menyeberang pagar.
3. **Penamaan, presisi uang, dan zona waktu** (`S-069`) — keputusan teknis, berasal dari
   ketetapan luar yang mengikat (aturan kerja `F` dan `G`).
4. **Tiap kunci ditegakkan dua kali** (`S-063`) — **punya cerita**: 18 dan 23. Penegakan
   kedua adalah satu-satunya sebab cerita 18 dapat dipenuhi ketika permintaan datang bukan
   lewat layar.

**Yang masih terbuka sesudah lapisan ini ditambahkan.** Empat hal, dicatat agar tidak hilang:

- ~~**Kolom, tipe, kunci, dan constraint per tabel**~~ — **ditutup 2026-09-21**, bagian
  *Model data*: tujuh objek, 64 kolom, masing-masing menunjuk `S-xxx`.
- ~~**Tabel medan layar**~~ — **ditutup 2026-09-21**, bagian *Medan layar*: tujuh medan
  dapat-sunting, sebelas kelompok baca-saja, lima blok dibuang dengan sebabnya.
- **Nama peran administratif** yang mengelola roster dan tabel seleksi: keputusan pemilik
  proses. Gerbangnya ditulis (`S-053`); nama perannya tidak. **Masih terbuka.**
- **Dua rule yang keadaannya menunggu keputusan** — jalur konversi dan penulis JSON klaim.
  Diperiksa 2026-09-21: **berkasnya ada di repo**, sehingga keduanya bukan lubang bukti dan
  tidak boleh dipagari. Yang belum ada adalah port hilirnya, dan itu keputusan pemilik
  proses. **Masih terbuka.**
- **Dua nama medan layar** dari enam yang `F-2` cacah belum dinamai di berkas dokumentasi
  mana pun. Namanya terbaca dari `Section/ShowTransfer.xml`, yang ada di repo — pembacaan
  yang belum dilakukan, bukan bukti yang belum dipegang, dan karena itu bukan pagar.
  **Masih terbuka.**

**Tiga perubahan perilaku kini bernomor.** Ketiganya lahir dari `K6-4` dan kini punya baris
di `REGISTER-DEVIASI.md` bagian 6 sebagai deviasi 23, 24, dan 25, masing-masing dengan uji
`PS-xx` yang dirancang gagal. Barisnya ditulis di register, bukan di sini — register bukan
milik spesifikasi, dan dua salinan akan perlahan berbeda.

**Pintu satu-satunya untuk membuka kembali penyelidikan.** Grilling ditutup. Ronde baru
dibuka hanya bila salah satu kueri §2.5 pulang dan hasilnya **bertentangan** dengan ketetapan
yang sudah ada. Lubang yang muncul saat membangun ditangani sebagai **pagar**, bukan sebagai
ronde.
