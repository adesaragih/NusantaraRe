> Modul  : Treaty In Adjustment · ronde C
> Dibuat : 2026-09-24
> Peran  : ditulis dengan peran penilai. Pertanyaan diajukan penggrill, dijawab dengan bukti,
>          diputuskan pemilik proses pada 24 September 2026 ("ikuti rekomendasi").
> Masukan: `00-LINGKUP.md`, `01-TEMUAN.md` (NC-01…NC-07), `PENGETAHUAN.md`, ADR induk 0036, 0037,
>          0040, 0042, 0043, 0049, 0053, 0054, 0055, dan GRL-01…GRL-13
> Status : **TERKUNCI** — empat putusan, `GRL-14` … `GRL-17`
> Sifat  : grilling. Tidak ada spesifikasi, DDL, struct Go, komponen React, endpoint, atau tiket.

# Putusan ronde C

| # | Judul | Label | Pembatalnya |
|---|---|---|---|
| **GRL-14** | `ActualValue` dipecah menurut artinya; tidak ada pohon ketiga di sistem baru | PERUBAHAN | `DB-15`, `DB-19` |
| **GRL-15** | Pro rata: atributnya dibawa sekarang, mesinnya **sengaja tidak dibangun** | atribut PERUBAHAN · mesin **keputusan untuk tidak membangun** | `DB-5` |
| **GRL-16** | Share fakultatif tidak memerlukan kemampuan tersendiri; yang hilang adalah **kelas** cacatnya | PERUBAHAN | `DB-18` |
| **GRL-17** | `NOMOR_URUT_VERSI` warisan diberikan ulang menurut kronologi; `/Rnn` dilestarikan sebagai pengenal | PERUBAHAN | `UA-21` mengukur luasnya, tidak membatalkan |

---

## GRL-14 — `ActualValue` dipecah menurut artinya

### Keputusan

**Tidak ada `ActualValue` di sistem baru**, dalam bentuk apa pun — bukan entitas, bukan pohon, bukan
kolom. Kedua artinya mendarat di tempat yang berbeda:

| `EDMState` lama | Arti yang terbaca (NC-05) | Ke mana ia pergi |
|---|---|---|
| 1, 2 | **potret sesudah-simpan** — keluaran | **dibuang.** Fungsinya sudah dipegang versi itu sendiri: versi disimpan penuh dan beku sesudah disetujui (ADR-0036, GRL-02). Potret dari sesuatu yang sudah beku adalah salinan, bukan fakta |
| 3 | **masukan pengguna** — nilai yang sedang diajukan | **menjadi nilai versi itu.** Premi aktual yang diketik pengisi adalah nilai `VERSI_KONTRAK`, sama seperti besaran lain |

### Kenapa begini

Sistem lama memakai satu nama untuk dua benda yang berlawanan arah — satu keluaran, satu masukan —
dan pembedanya sebuah kode jenis. Membawa nama itu ke sistem baru berarti menanam cacat yang sama
dengan tangan sendiri.

Dan pemecahannya **menghapus TDA-16 tanpa mekanisme tambahan**. Di sistem lama,
`TreatyEDMDifferencePremium` beriterasi `TreatyIn.EGNPI` sementara pengisi menyunting
`TreatyIn.ActualValue.EGNPI`, sehingga selisih EGNPI pada addendum premi **selalu nol** — angka yang
menjadi alasan jenis addendum itu ada tidak pernah muncul sebagai selisih. Begitu premi aktual
menjadi nilai versi, mengubahnya menghasilkan baris `NILAI_SELISIH` dengan sendirinya. Tidak ada
pengecualian yang perlu ditulis tangan, dan pengecualian bertangan adalah persis bentuk yang GRL-12
dibuat untuk menghapus.

### Label

**PERUBAHAN.** Berubah dari: *"pada addendum premi, nilai aktual hidup di pohon terpisah dan tidak
pernah masuk selisih"*.

### Batas kasus yang BERUBAH, dan ia dinyatakan di sini karena ia tidak lestari

GRL-13 mencatat *"penyesuaian premi selalu material"*. Sesudah GRL-12 dan GRL-14, kalimat itu
**menyusut dengan sendirinya** menjadi:

> **Penyesuaian premi yang mengubah premi aktual selalu material.**

Sisanya — sebuah `PENYESUAIAN_PREMI` yang **tidak mengubah apa pun** — terbaca **non-material**,
karena ia tidak punya baris `NILAI_SELISIH`. Sistem lama menyebutnya material, sebab materialitasnya
melekat pada **jenis**, bukan pada akibat.

Itu **PERUBAHAN, bukan cacat**, dan arahnya benar: versi yang tidak mengubah apa pun memang tidak
mengubah apa pun. Tetapi ia tidak boleh muncul sebagai kejutan di hari peralihan, jadi ia:

- ditanyakan ke bisnis sebagai **`DB-19`** — apakah penyesuaian premi yang tidak mengubah angka
  memang pernah dibuat, dan apakah ia dianggap material;
- dan bila `DB-19` dibantah, yang berubah **bukan** GRL-14 melainkan GRL-12, ke arah yang sudah
  disediakan pilihan (c)-nya.

### Titipan dan akibat

| Ke mana | Isi |
|---|---|
| **E3c** (TDA-06) | **larut.** "Ringkasan share yang ditulis ke `ActualValue` lalu ditimpa" berhenti menjadi pilihan rancangan dan menjadi pemerian cacat lama, karena tidak ada `ActualValue` yang dapat ditimpa. Nasib TDA-06: **diperbaiki** |
| **TDA-16** | **diperbaiki** — oleh bentuk, bukan oleh tambalan |
| bahan to-spec | premi aktual adalah besaran `VERSI_KONTRAK` biasa; ia masuk `BESARAN_DAPAT_DISESUAIKAN` seperti besaran lain |
| sisi | **IRISAN** — `SaveTreatyIn_EDM_Act` dan mesin selisih ada di ekspor induk juga (NC-07); usulan untuk langkah 8–10 induk |

---

### Penguatan 24 September 2026 — buktinya kini TERKALIBRASI, dan satu berkas induk dikoreksi

GRL-14 dikunci dengan bukti bahwa `ActualValue` punya dua arti. **Sesi to-spec Treaty In 24
September menemukan sebab mekanisnya**, dan ia menguatkan putusan ini tanpa mengubahnya.

**Arah prasyaratnya terbaca.** `SaveTreatyIn_EDM_Act` langkah 2 berkondisi `TreatyIn.EDMState=="3"`
dengan `pyStepsPreCondParamsWhenTrue = 1` dan parameter `jmp`; langkah 5 ber-`pyStepsBlockName = jmp`.

Arti kode `1` **dikalibrasi, bukan diandaikan** — disapu atas **kedua ekspor**, seluruh
`pyStepsPreCondParamsWhenTrue`/`WhenFalse`:

| Kode | Jumlah | Membawa nama blok? |
|---:|---:|---|
| **1** | **66** | **66 dari 66 — selalu** |
| 2 | 18.412 | tidak pernah — dan ia muncul pada langkah yang **tidak berkondisi sama sekali** |
| 3 · 4 · 5 | 4.890 · 412 · 16 | tidak pernah |

**Maka kode `1` = lompat ke blok**, dan untuk `EDMState == "3"` **langkah 2, 3, dan 4 terlewat
seluruhnya**.

**Akibatnya pada pembacaan dua arti:**

| Jenis | Langkah 2-3-4 | `ActualValue` pada penyimpanan |
|---|---|---|
| **1, 2** | berjalan | **ditimpa salinan utuh pohon utama setiap simpan** — karena itu ia **keluaran**, dan potret dari sesuatu yang sudah beku adalah salinan |
| **3** | **dilompati** | **tidak disentuh** — isinya tetap apa yang diketik pengisi; karena itu ia **masukan**, dan **tidak punya rumah lain** |

Pembedaan keluaran-versus-masukan di tabel putusan di atas kini **terbaca dari mesinnya**, bukan
hanya dari arti bisnisnya.

> **Satu berkas induk berbunyi terbalik selama ini, dan sudah dikoreksi.**
> `../treaty-in/4-erd-dan-tabel-datar/7-2-ACTUALVALUE.md` §1.2–§1.4 menyatakan jenis 3 yang dipotret
> dan jenis 1–2 yang tidak. Koreksinya diterapkan 24 September 2026 atas persetujuan pemilik proses,
> beserta pemindahan alamat cacat §1.3 dan sarang §1.4. Uraian lama dipertahankan sebagai alasannya.
>
> **Aturan bacanya diwariskan:** `../treaty-in/CONTEXT.md` §2.1 kini menyatakan bahwa prasyarat yang
> **hidup** pun harus dibaca bersama **kode tindakannya** — kondisi tanpa kode tindakan hanya separuh
> kalimat, dan separuh yang tersisa dapat membalik artinya.

**Sisa yang tidak dihitung sebagai cacat:** 20 dari 66 pemakaian kode `1` menunjuk nama blok yang
tidak ditemukan di aktivitas yang sama. Sebabnya belum diperiksa.

---

## GRL-15 — pro rata: atributnya dibawa, mesinnya sengaja tidak dibangun

### Keputusan

Dua bagian, dan keduanya perlu dibaca bersama:

1. **`VERSI_KONTRAK` membawa atribut "berlaku sejak"** sejak hari pertama, dengan nilai bawaan =
   tanggal mulai kontrak. Nama finalnya bahan to-spec.
2. **Perhitungan prorata TIDAK dibangun**, dan itu **keputusan**, bukan pekerjaan yang tertunda.

### Kenapa begini

Tiga sapuan ronde ini menunjukkan pro rata di sistem lama **hidup, dipanggil, dan tidak pernah dapat
menghasilkan apa pun selain 100 %**:

```
EDMEffective  : satu penulis nilai — TreatyInSetEditPre 1.5 → = Commencement
                kontrol layarnya pxDateTime HANYA-BACA di kedua ekspor          (NC-01)
IsProRate     : dua penulisan, keduanya pxCheckbox layar                        (NC-02)
ProRatePercent: @divide(ProRateDays, ProRateTotalDays) × 100
                ProRateDays      = selisih(EDMEffective, Termination)
                ProRateTotalDays = selisih(Commencement, Termination)
                dan EDMEffective == Commencement  →  1 × 100 = 100              (NC-03)
```

Maka membangun mesinnya sekarang bukan pelestarian: ia **menciptakan kemampuan yang belum pernah ada
yang meminta**. Sementara itu, atributnya murah dibawa dan mahal dipasang belakangan — ia mengubah
bentuk setiap baris selisih yang sudah beku.

### Bentuk "sengaja tidak dibangun", ditulis sebagai keputusan

Supaya orang berikutnya **tidak memperbaikinya dengan cara yang merusaknya**:

| | |
|---|---|
| **Apa yang sengaja tidak dilakukan** | rumus prorata: membagi besaran versi menurut porsi periode yang tersisa sejak "berlaku sejak" |
| **Kenapa** | sistem lama tidak pernah menjalankannya sungguhan (NC-01+NC-03), dan `DB-5` — apakah addendum pernah berlaku di tengah periode — belum dijawab bisnis |
| **Akibatnya pada angka** | setiap versi berlaku penuh sejak tanggal mulai kontrak. Tidak ada angka yang berbeda dari sistem lama, karena sistem lama pun selalu 100 % |
| **Di mana selisihnya terlihat** | `UA-19` — sebaran waktu baris ber-`ProRatePercent ≠ 100` di produksi |
| **Kapan ia ditagih** | saat `DB-5` dijawab. Bila dibantah, yang dibangun adalah rumusnya saja; atributnya sudah ada |

### Label

- **Atribut "berlaku sejak" → PERUBAHAN.** ADR-0037 menetapkan *"faktor prorata addendum = turunan
  dari tanggal, bukan centang"*; sistem lama adalah **centang** (NC-02). ADR-nya benar, labelnya yang
  belum ada — diajukan sebagai **`REV-6`**.
- **Mesin prorata → tidak dibangun**, dan itu bukan salah satu dari tiga label. Ia keputusan
  bertanggal dengan saat penagihan.

### Sisi

**NC-02 dan NC-03 berada di IRISAN.** `TreatyCalculateProratePct` dan `Section/TreatyInNONProportional`
ada di ekspor Treaty In juga, dan **tidak** ada di katalog 56. Maka pro rata yang tidak pernah dapat
menghasilkan selain 100 % berjalan di **modul induk hari ini** — usulan untuk langkah 8–10 induk.

---

## GRL-16 — share fakultatif: yang hilang adalah KELAS cacatnya

### Keputusan

**Tidak ada kemampuan yang perlu dibangun untuk share fakultatif.** `NILAI_SELISIH` diturunkan per
`KUNCI_PADANAN` atas baris anak yang tersimpan (`SEAM-ADJUSTMENT.md` §3, E1 KONFIRMASI), sehingga
share fakultatif ikut terbawa seperti besaran lain — tanpa satu baris pun yang menyebut namanya.

### Kenapa begini

Di sistem lama ada **dua ketiadaan yang berdiri sendiri**, dan yang pertama tidak menutup yang kedua
(NC-04):

| # | Ketiadaan | Sifatnya |
|---|---|---|
| 1 | `TreatyEDMDifferenceDeduction` langkah 4–5 **MATI**, label aslinya *"Facultative share not enable yet in this edm"* | dimatikan **dengan sengaja** |
| 2 | `TreatyInDifferenceFacShare` **tidak pernah dipanggil** dari rantai `TreatyEDMCalculateDifference` | **tidak pernah disambungkan** |

Akar `FacultativeShare` dan `FacultativeShareBrokerage` berpenulis tunggal aturan di butir 2, maka
keduanya tidak pernah terisi lewat jalur pengajuan addendum. Menyalakan blok mati **tidak** menutup
butir 2 — dan itulah sebab jawabannya bukan "nyalakan".

Yang memutuskan bukan nilai share fakultatifnya melainkan **bentuk mesinnya**. Mesin lama mendaftar
besaran satu per satu, sehingga setiap besaran yang lupa didaftarkan menjadi ketiadaan tersendiri.
Mesin yang menurunkan selisih dari baris tersimpan **tidak punya daftar untuk dilupakan**.

### Yang WAJIB menyertainya — dan tanpa ini keputusannya hanya klaim

*"Terbawa sendiri"* adalah pernyataan rancangan. Pernyataan rancangan yang tidak pernah diperiksa
berperilaku persis seperti kemampuan yang tidak pernah disambungkan — yaitu **seperti butir 2 di
atas**.

Maka satu **uji negatif** ikut ditulis sebagai bahan to-spec:

> Mengubah share fakultatif pada sebuah versi **harus** menghasilkan sedikitnya satu baris
> `NILAI_SELISIH`. Uji yang tidak pernah gagal atas besaran ini berarti besaran itu tidak ada di
> dalam himpunan yang dipadankan.

### Label

**PERUBAHAN.** Berubah dari: *"selisih share fakultatif tidak pernah dihitung pada jalur pengajuan
addendum"* — bukan dari "dihitung dengan cara lain".

### Pembatal

**`DB-18`** — bila bisnis menyatakan share fakultatif memang **sengaja** dikecualikan dari selisih
addendum, yang dibutuhkan satu penyaring, bukan pembongkaran. Keputusan ini tidak berubah bentuk;
hanya himpunan besarannya yang menyusut.

---

## GRL-17 — `NOMOR_URUT_VERSI` warisan diberikan ulang menurut kronologi

### Payungnya lebih dulu: I2 turun menjadi KONFIRMASI

Ketiga ADR dibaca (prasyarat `MA-10`, NC-06). Pertanyaan umum I2 — *"aturan baru mana yang boleh
dijalankan atas baris warisan, dan asal-usulnya dicatat bagaimana"* — **sudah dijawab**:

| Sumber | Yang menjawab |
|---|---|
| **ADR-0042** butir 2 | invarian ditegakkan pada **penulisan**, bukan pada baris; baris warisan ditulis sekali dan tidak dapat disunting |
| **ADR-0042** butir 3 | **sentuh-perbaiki** — begitu kontrak warisan disentuh, ia memenuhi invarian saat itu juga |
| **ADR-0042** §materialitas | kolom membawa **asal-usulnya**: diturunkan aturan, atau diimpor sebagai warisan |
| **ADR-0043** | migrasi **memindahkan**; menghitung ulang adalah **peristiwa bisnis tersendiri** |
| **ADR-0054** | nilai warisan tanpa padanan → `WARISAN_TAK_TERPETAKAN`, nilai asli tersimpan |

### Keputusan atas residunya

Ketiga ADR menjawab *apa yang boleh dilakukan atas **nilai** warisan*. Tidak satu pun menjawab
*berapa `NOMOR_URUT_VERSI` sebuah baris warisan* — dan itu **bukan nilai warisan**, melainkan
**mekanisme sistem baru** yang harus diberikan kepada baris lama.

| Yang disimpan | Dari mana |
|---|---|
| **pengenal tampilan `/Rnn`** | dilestarikan **apa adanya** — sudah diputuskan GRL-09, dan nomornya sudah beredar di luar sistem |
| **`NOMOR_URUT_VERSI`** | **diberikan ulang menurut kronologi**: sumber utama tanggal komentar pembuatan (*"Had Created …"*, ditulis `TreatyInSetEditPre`), cadangan `EDMDATE` |

### Kenapa begini

GRL-11 menurunkan versi berlaku dari `NOMOR_URUT_VERSI` **tertinggi**, bersyarat *"urutan nomor =
urutan persetujuan"*. Untuk baris warisan syarat itu **tidak dapat dipenuhi dari nomor lamanya**:
`PENGETAHUAN.md` §4.4 membuktikan penomoran `/Rnn` meleset dan dapat dipakai ulang — memilih
`…/R01` ketika R02 sudah ada menghasilkan **R02 lagi**, dan tabrakannya berakhir sebagai `UPDATE`
yang menimpa (TDA-01).

> **Nomor yang dapat dipakai ulang bukan urutan.** Menurunkan keadaan hukum sebuah kontrak darinya
> berarti menurunkannya dari penghitung yang sudah terbukti salah.

Dan `EDMDATE` hanya **cadangan**, bukan sumber utama, karena `EDMDATE = SYSDATE` pada **setiap**
pembaruan — ia waktu sentuh terakhir, bukan waktu lahir.

### Kelas yang harus DINYATAKAN, bukan ditebak

Baris warisan yang kronologinya **tidak dapat ditentukan dari sumber mana pun** tidak dipaksa masuk
urutan. Ia menjadi **pengecualian migrasi bernomor**, dilaporkan **per kontrak**, dan kontrak yang
memuatnya tidak memperoleh "versi berlaku" turunan sampai seseorang memutuskannya.

Bentuknya meniru ADR-0054 — sebuah kelas bernama untuk yang tidak berpadanan — meski ADR itu tentang
**keadaan** dan ini tentang **urutan**. Kemiripan bentuknya disebut supaya tidak dikira instans dari
ADR-0054; ia bukan.

### Satu tabrakan dengan ADR-0043 yang WAJIB tertulis

ADR-0043 menetapkan ukuran keberhasilan migrasi: **"setiap angka identik. Bukan mirip, bukan dalam
toleransi."** Dibaca polos, memberikan nomor urut baru **melanggarnya**.

Ia tidak melanggar, dan sebabnya harus tertulis atau uji paritas akan menandainya sebagai kegagalan:

> `NOMOR_URUT_VERSI` **tidak punya padanan di sistem lama**. Sistem lama tidak menyimpan nomor urut;
> ia hanya menyimpan pengenal. Maka nomor urut bukan angka yang **dipindahkan** melainkan mekanisme
> yang **diberikan**, dan kriteria identik berlaku atas angka yang dipindahkan.
>
> Yang tunduk pada kriteria identik adalah **pengenal `/Rnn`** — dan itu memang dilestarikan apa
> adanya.

Kalimat ini masuk daftar harapan uji paritas (ADR-0043 §"syarat yang tidak boleh dilewati"), supaya
pembacanya tidak merasionalisasi selisih yang muncul.

### Label

**PERUBAHAN.** Berubah dari: *"nomor urut sebuah versi adalah angka yang tertulis pada pengenalnya"*.

### Yang diukur, dan apa yang tidak berubah karenanya

**`UA-21`** mengukur berapa kontrak memuat nomor `/Rnn` ganda atau melompat, dan berapa addendum
yang tertimpa. Hasilnya **tidak mengubah keputusan ini** — bila nol, aturannya tetap ditulis, karena
ia menutup kelas persoalan dan bukan satu nilai. Yang berubah hanya perkiraan berapa banyak
pengecualian migrasi yang menunggu.

---

## Yang ronde ini TIDAK putuskan

| Hal | Sebab |
|---|---|
| nama final atribut "berlaku sejak" dan penamaan besaran premi aktual | bahan to-spec, bukan grilling |
| bentuk fisik `NILAI_SELISIH` | sesi DDL; keputusan logisnya sudah ada (E1 KONFIRMASI) |
| nasib tujuh TDA yang cabangnya sudah ditutup | **ronde D** — lihat `03-LUBANG.md` §4 |
| wewenang tombol "(dev)" | manajemen; sudah di daftar eskalasi induk |
