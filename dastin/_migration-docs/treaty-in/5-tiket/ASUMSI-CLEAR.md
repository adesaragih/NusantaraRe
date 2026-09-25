# Asumsi yang dianggap clear — dan apa yang ditinjau bila ia patah

**Tanggal:** 24 September 2026 · **Dibuat saat:** to-ticket kemampuan Adjustment dimulai
**Diperbarui:** 24 September 2026 sore — kolom *"tiket yang bersandar"* **diisi dari sapuan**, bukan dari ingatan, sesudah tiket `14`…`44` ditulis.
**Sifat:** catatan gerbang. **Bukan** pernyataan bahwa penahannya tertutup.

> ## KENAPA BERKAS INI ADA
>
> Pemilik proses memerintahkan to-ticket berjalan **dengan sejumlah penahan dianggap clear**.
> Itu keputusan yang sah dan tidak dipersoalkan di sini.
>
> Yang dipersoalkan hanya satu hal: **asumsi yang tidak tercatat tidak dapat disapu balik.** Bila
> salah satunya patah tiga bulan lagi, pertanyaannya menjadi *"tiket mana yang bersandar padanya"* —
> dan tanpa berkas ini, jawabannya menuntut membaca ulang setiap tiket.
>
> **Tanpa berkas ini, "anggap clear" berubah menjadi "lupa".**

## Cara memakainya

Setiap tiket yang bersandar pada salah satu asumsi menuliskannya di medan `Dasar:` dengan bentuk
yang **dapat disapu satu perintah**:

```
Dasar:  EVIDENCED(TreatyEDMDifferenceLimits@ekspor-2026-09)
        DECIDED(GRL-20)
        DIASUMSIKAN-CLEAR(REV-3)        <- inilah yang dapat disapu
```

Ketika sebuah asumsi patah, satu sapuan atas `DIASUMSIKAN-CLEAR(<kode>)` mengeluarkan daftar tiket
yang harus ditinjau. **Tidak ada yang perlu diingat.**

> ### SAPUANNYA KINI MENUNTUT DUA FOLDER — 25 September 2026
>
> Papan dipisah atas perintah pemilik proses: tiket `01`…`13` pindah ke
> `treaty-in-adjustment/5-tiket/issues/`, tiket `14`…`64` tetap di `treaty-in/5-tiket/issues/`.
>
> ```
> grep -r "DIASUMSIKAN-CLEAR(" treaty-in/5-tiket/issues/ treaty-in-adjustment/5-tiket/issues/
> ```
>
> **Berkas ini TIDAK ikut dipisah**, dan itu keputusan: enam dari tiga belas kodenya dipakai **kedua
> papan** — `KTV-A` oleh 19 tiket papan induk, `DB-20` oleh empat tiket papan Adjustment. Memecahnya
> akan menghasilkan dua daftar yang masing-masing **terlihat lengkap**.
>
> **Ongkosnya:** sapuan yang hanya menyisir satu folder mengembalikan daftar yang **tidak lengkap
> tanpa mengatakannya**. Kolom *"tiket yang bersandar"* di bawah kini menandai papannya.

---

## 1. Empat asumsi yang pemilik proses sebut

| # | Kode sapuan | Diasumsikan clear | Kenyataannya di berkas | Tiket yang bersandar |
|---|---|---|---|---|
| 1 | `D-7` | kalimat pembaca arsip sudah mendarat di induk | **masih diparkir** — dasar parkirnya sudah lewat, sebab `F-1` ternyata sudah diterapkan; tinggal satu kalimat izin | **`42`** |
| ~~2~~ | ~~`D-1`~~ | ~~catatan `SIFAT_MATERIAL_ADDENDUM` sudah diganti di §10.2a induk~~ | **DICORET 24 September 2026 oleh pemilik proses.** `F-2` ditutup di blok kepala §10.2 — pencarian keempat atribut diselesaikan atas sumber mandiri, nol calon tersisa, angka judul diturunkan ke **45** yang dihitung (43 §10.2 + 1 §10.2a + 1 §10.19a). Parkir `D-1` dicabut, dan **diffnya sudah diterapkan** di §10.2a | **tidak ada** — asumsi ini tidak lagi dipakai di `Dasar:` tiket mana pun |
| 3 | `REV-3` | paket `REV-1`…`REV-6` sudah ditanggapi pemilik ADR | **belum diserahkan sama sekali** | **`01`, `07`** — **dan hanya keduanya.** ~~*dan 20 butir batch 2 menunggunya*~~ **RALAT 25 Sep 2026**: batch 2 **tidak** menunggu `REV-3`; penahannya **ukuran sesi**. `KTV-4`: *"nol dari enam menentukan letak kolom"* |
| 4 | `DB-20` | titik beku materialitas | **belum dikirim** | **`02`, `07`, `11`, `13`** |
| 4a | `DB-16a` | dokumen addendum dikirim ke luar atau tidak | **belum dikirim** | **`04`, `09`** |
| 4b | `DB-16b` | dokumen punya tanggal berlaku sendiri atau tidak | **belum dikirim** | **`08`** — **bertenggat**, §3 |

### Apa yang ditinjau ulang bila masing-masing patah

| Asumsi | Yang ditinjau, dan seberapa dalam |
|---|---|
| **`D-7`** | kalimat pembaca arsip lama. **Paling dangkal** — ia satu paragraf di artefak induk, tidak menyentuh kolom maupun invarian mana pun |
| **`D-1`** | catatan pada `SIFAT_MATERIAL_ADDENDUM`. Tidak mengubah kolomnya; mengubah **keterangannya**. Dangkal, **kecuali** bila `F-2` ternyata memulihkan atribut yang bertabrakan namanya |
| **`REV-3`** | **paling dalam.** `REV-3` mengusulkan koreksi ADR-0055 §4, dan ADR-0055 adalah sumber **daftar keadaan**. Bila pemilik ADR menolaknya atau mengubahnya, **kolom keadaan dapat berubah** — dan setiap tiket daur hidup, persetujuan, penolakan, dan pembatalan ikut ditinjau |
| **`DB-20`** | titik beku materialitas. Bila bisnis menyatakan materialitas beku **sejak lahir**, bukan sejak diajukan, maka `KTV-1` berubah dan tiket yang mengizinkan pembetulan selama `DRAFT` **menolak sesuatu yang seharusnya diizinkan, atau sebaliknya** |
| **`DB-16a`** | apakah dokumen addendum dikirim ke luar. Bila ya, `KTV-3` berubah: jenis bernilai **tiga**, bukan dua, dan penanda "dikirim ke luar" kembali |
| **`DB-16b`** | tanggal berlaku dokumen. Bila dibantah, `KTV-2` menuntut kolomnya **dicabut sebelum data masuk** — dan itu **bertenggat**, lihat §3 |

---

## 1a. Sembilan asumsi lagi yang lahir dari batch 1 Treaty In — 24 September 2026

Keempat asumsi §1 lahir ketika papan hanya memuat kemampuan Adjustment. Tiket `14`…`44` menambah
sembilan kode lagi, dan **seluruhnya diambil dari sapuan `DIASUMSIKAN-CLEAR(` atas folder
`issues/`** — bukan dari ingatan.

| # | Kode sapuan | Diasumsikan clear | Kenyataannya | Tiket yang bersandar |
|---|---|---|---|---|
| 7 | **`KTV-A`** | presisi kolom sudah ditetapkan gerbang sesi DDL | **diputuskan tanpa verifikasi** — `KEPUTUSAN-TANPA-VERIFIKASI.md` §7. **BERTENGGAT**, lihat §3a | **19 tiket** — `14`, `15`, `20`, `22`…`31`, `33`, `34`, `37`, `38`, `39`, `44` |
| 8 | **`T-6`** | kelima paket uang golongan C memang perlu mata uang sendiri | **belum dikirim** ke teknik treaty | `14`, `31` |
| 9 | **`T-2`** | tingkat pencatatan `NILAI_BATAS` diketahui | idem | `29` |
| 10 | **`T-3`** | tingkat pencatatan `LIMIT_AGREGAT` diketahui | idem | `31` |
| 11 | **`T-4`** | tingkat pencatatan `NILAI_EGNPI` diketahui | idem | `23` |
| 12 | **`F-15`** | nasib `PERSEN_ROL` dan premi diperoleh sudah diputuskan | **menunggu satu kalimat pemilik proses** — `2-to-spec/F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` §5 | `31` |
| 13 | **`F-16`** | entitas peran dan penugasan bertanggal ada di §10 | **tidak ada satu pun**, padahal `ADR-0044` menaruh keduanya di gelombang 1 | `39` |
| 14 | **`L-3`** | ada instans Oracle yang terjangkau untuk menjalankan uji | **tidak ada** — lubang **lingkungan**, bukan lubang spesifikasi | `43`, dan **`38`** lewat medan `PENGHALANG`-nya |
| 15 | **`Uji AD`** | golongan `P-18` sudah pasti `L` atau `B` | **belum dijalankan**; izin kueri baca-saja belum turun | `38` lewat medan `PENGHALANG`-nya |

### Apa yang ditinjau ulang bila masing-masing patah

| Asumsi | Yang ditinjau, dan seberapa dalam |
|---|---|
| **`KTV-A`** | **paling luas, tetapi paling dangkal per tiket**: sembilan belas tiket, dan yang berubah **lebar kolom**, bukan bentuk tabel. Yang membuatnya berbahaya bukan kedalamannya melainkan **tenggatnya** |
| **`T-6`** | lima kolom mata uang **dicabut**, bukan diubah. `14` kehilangan tiga, `31` kehilangan dua. Tidak ada data yang hilang — kolomnya memang kosong |
| **`T-2`, `T-3`, `T-4`** | satu `CHECK` per paket uang **ditambahkan** (`INV-39`, `INV-40`). Menambah constraint pada kolom yang sudah ada **berongkos**, sebab ia menuntut memeriksa baris yang sudah masuk |
| **`F-15`** | bila **turunan**: `PERSEN_ROL` dicabut dari `LAYER`, dan tabel anak premi diperoleh **tidak pernah dibuat** — tiket `31` berkurang satu kolom. Bila **potret beku**: satu entitas **bertambah**, dan `STRUKTUR-DATA.md` ikut disunting |
| **`F-16`** | **bukan bentuk `PERAN_PELAKU`** — itu sudah diputuskan `KTV-D` sebagai potret, dan potret tidak merujuk apa pun. Yang ditinjau **dari mana nilainya diambil** saat jejak ditulis |
| **`L-3`** | tiap uji negatif yang hari ini **diargumentasikan** harus **dijalankan**. Yang paling menggantung: uji negatif *materialized view* untuk `INV-47`, `INV-50`, `INV-51` di tiket `38`, dan `INV-69`/`INV-70` di tiket `11` |
| **`Uji AD`** | golongan `P-18`. Bila **B**, tiket `38` **menuntut medan `CARA MENYALAKANNYA`** dengan keempat butirnya. Bentuk tabelnya tidak berubah |

---

## 2. Dua penahan lagi yang ada di berkas dan TIDAK disebut sebagai asumsi

`2-to-spec/SERAH-TERIMA.md` §3 mendaftar **enam** penahan, bukan empat. Kedua sisanya dicatat di sini
supaya daftar ini tidak lebih pendek daripada kenyataannya:

| # | Penahan | Keadaan 24 September 2026 |
|---|---|---|
| 5 | angka **presisi** dan bentuk **induk polimorfik** | **terbuka** — gerbang sesi DDL induk. Ia menahan bentuk fisik, **bukan nama kolom**: urutan wewenang butir 5 sudah memberi `KAMUS-KOLOM.md` sebagai yang mengikat soal nama dan tipe |
| 6 | **tiga kolom `TDA-17` tanpa rumah** | **SATU DITUTUP, DUA DIPERTIKAIKAN — diperbarui 24 September 2026 sore.** `DEDUCTIBLE2` → `LAYER.DEDUCTIBLE_KEDUA` **berdiri, dan sapuan penulisnya bersih**. `ROL_PCT` → `LAYER.PERSEN_ROL` **berdiri tetapi ditandai `F-15`** — `DetailCalculationROL` menghitungnya `premi ÷ limit × 100`. `PREMIUM_EARNED` **tidak dibuat** — `DetailCalculation` langkah 25 menghitungnya `EGNPI × tarif penyesuaian`, dan dua artefak induk sudah mengadilinya **DITURUNKAN** ber-`EVIDENCED`. **Keduanya menunggu satu kalimat pemilik proses**, `2-to-spec/F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` §5 |

**Butir 5 tidak diberi kode sapuan** karena tidak ada tiket yang boleh bersandar padanya: butir 3
urutan wewenang **tetap dicoret**, sehingga tidak ada tiket yang sah mengaku bersandar pada bentuk
tabel dari berkas DDL.

**Butir 6 diberi kode `F-15`**, bukan lagi `TDA-17`, untuk kedua kolom yang dipertikaikan — satu
kode per penanda, dan `TDA-17` sudah menjadi putusan yang sebagian dibantah. Irisan **31** bersandar
padanya.

---

## 3. Satu asumsi yang BERTENGGAT, dan tenggatnya bukan tanggal

`KTV-2` mensyaratkan kolom tanggal berlaku dokumen **dicabut sebelum data masuk** bila `DB-16b`
dibantah.

> Sesudah migrasi berjalan, mencabut kolom **berongkos** — ia menuntut memeriksa apakah ada yang
> sudah mengisinya, dan apa yang terjadi pada baris yang mengisinya.

Maka tenggat asumsi ini adalah **momen pemuatan data pertama**, bukan sebuah tanggal. Tiket yang
membangun `DOKUMEN_ADDENDUM` wajib menyebutnya, dan tiket migrasi wajib **memeriksa berkas ini**
sebelum berjalan.

---

## 3a. Asumsi bertenggat yang KEDUA — `KTV-A`, dan tenggatnya sama dengan `KTV-2`

`KTV-A` memutuskan presisi kolom **tanpa verifikasi**, dengan syarat pembalikan yang berbunyi: *"sesi
DDL berikutnya boleh mempersempit, dan **hanya sebelum data dimuat**."*

> Sesudah data masuk, mempersempit kolom menuntut **memeriksa setiap baris**, dan baris yang tidak
> muat **tidak punya tempat pergi**.

**Yang memuat data pertama adalah tiket `44`** — pemindahan isi enam tabel acuan. Maka tenggat
`KTV-A` adalah **momen tiket `44` dijalankan**, dan tiket `44` **bertanda tertahan justru karena
itu**.

**Dua asumsi kini bertenggat pada momen yang sama**, dan keduanya wajib diperiksa sebelum tiket
migrasi mana pun berjalan:

| Asumsi | Tenggatnya | Tiket yang menutup tenggatnya |
|---|---|---|
| `KTV-2` | momen pemuatan data pertama | tiket `08` menyatakannya; tiket migrasi memeriksanya |
| **`KTV-A`** | idem | **tiket `44`** |

### Dan satu asumsi yang bertenggat LEBIH AWAL — `KTV-D`

`KTV-D` menetapkan `PERAN_PELAKU` sebagai **potret**. Tenggatnya **bukan** pemuatan data pertama
melainkan **baris jejak pertama**, yang datang lebih awal: begitu tiket `39` menyala, jejaknya mulai
terisi.

> Sesudah itu, membalikkannya menuntut **menulis ulang sejarah** — hal yang jejak ini ada untuk
> mencegahnya. `KTV-D` **tidak diberi kode sapuan** sebab tidak ada yang menunggunya; ia diputuskan,
> bukan diasumsikan clear.

---

## 4. Cara memeriksa berkas ini tetap benar

**Dua arah, pada tiap titik periksa:**

| Arah | Pertanyaannya |
|---|---|
| **1** | adakah `DIASUMSIKAN-CLEAR(<kode>)` di tiket yang **kodenya tidak ada** di berkas ini? |
| **2** | adakah asumsi di berkas ini yang **kolom tiketnya kosong** padahal tiketnya sudah ditulis? |

Arah kedua yang paling mudah terlewat — dan ia yang membuat berkas ini menjadi hiasan bila
dibiarkan.

**Dan ketika sebuah asumsi ditutup:** barisnya **dicoret, tidak dihapus**, dengan tanggal dan siapa
yang menutupnya. Rantainya harus tetap terbaca oleh orang yang membaca tiket lamanya.
