# Usulan irisan — kemampuan Adjustment

**Tanggal:** 24 September 2026 · **Langkah 2 to-ticket**
**Status: DISETUJUI 24 September 2026, dengan dua koreksi** — tiket ditulis di
[`issues/`](issues/).

> ## DUA KORESI PEMILIK PROSES, DAN AKIBATNYA PADA DAFTAR DI BAWAH
>
> **Daftar di bawah dipertahankan apa adanya sebagai bentuk yang disetujui.** Yang berubah dicatat
> di sini, bukan disunting diam-diam ke dalam tabelnya — supaya yang sudah membacanya tahu apa yang
> bergeser.
>
> ### Koreksi 1 — irisan 02 dipecah. **Dua belas menjadi tiga belas.**
>
> `INV-69` menegakkan **atas baris `NILAI_SELISIH`**, dan `NILAI_SELISIH` baru lahir di irisan
> selisih. Irisan 02 semula memuat **dua** hal dengan **satu** sisi penghalang.
>
> **Pilihan saya: PECAH**, bukan memblokir 02 dengan 08. Sebabnya dapat diperiksa — **penguncian
> ruas tidak memerlukan `NILAI_SELISIH` sama sekali**: ia memeriksa **kolom mana yang berubah** pada
> penyimpanan, bukan baris selisih yang dihasilkannya. Memblokirnya di belakang tabel selisih
> menunda sesuatu yang tidak bergantung padanya, dan menggabungkannya membuat satu tiket dinyatakan
> selesai sementara separuh penegakannya belum ada — persis `TDA-10`.
>
> | Semula | Menjadi |
> |---|---|
> | 02 *(penguncian + `INV-69`)* | **02** penguncian ruas, ditegakkan sisi simpan · **11** `INV-69`/`INV-70` + pemantau kebasian |
>
> ### Koreksi 2 — irisan reinstatement ditulis ulang, "terbaca" dibuang
>
> *"Terbaca Material"* adalah bahasa `GRL-12`, dan **`GRL-12` BATAL**. Di bawah `GRL-20` materialitas
> **dinyatakan pengisi**; tidak ada yang terbaca. Perilakunya kini berbunyi: versi yang mengubah
> persentase pemulihan limit **dan dinyatakan `TIDAK_MATERIAL` ditolak saat simpan** oleh `INV-69` —
> dan di sistem lama hal yang sama **lolos tanpa sepatah galat**. Itu yang masuk pemberitahuan
> pra-peralihan.
>
> ### Penomoran sesudah kedua koreksi
>
> `01` P-42 · `02` penguncian ruas · `03` P-45 · `04` dokumen · `05` perluas · `06` selisih ·
> `07` beku+jejak · `08` tanggal dokumen · `09` arsip kertas · `10` pindahkan · `11` `INV-69`+pemantau ·
> `12` kerutkan · `13` reinstatement.
>
> ### Hitungan yang berlaku — dan ia berbeda dari yang pemilik proses tulis
>
> Perintahnya menetapkan **"dapat dimulai: 3"**. Angka itu benar **untuk pilihan yang tidak saya
> ambil** — bila 02 diblokir 08. Karena 02 **dipecah** dan penguncian ruas tidak berpenghalang,
> angkanya **4**: `01`, `02`, `03`, `05`.
>
> Dilaporkan begini, bukan diselaraskan diam-diam, sebab selisihnya **akibat dari pilihan yang
> didelegasikan kepada saya** — dan bila pilihan itu ditolak, angkanya kembali 3.
>
> | | |
> |---|---:|
> | tiket | **13** |
> | aktif · tertahan | 10 · 3 |
> | **dapat dimulai** | **4** |

> **Usulan yang menunggu persetujuan ADALAH keluaran.** Ia ditulis ke disk karena usulan yang hanya
> hidup di balasan **hilang begitu sesinya tutup** — dan yang hilang akan disusun ulang dari awal
> oleh orang yang tidak tahu ia pernah ada.

## Lingkup

**Sepuluh kemampuan**: tujuh **BARU** (`P-60`…`P-66`) dan tiga **PERLU DIUBAH** (`P-04`, `P-42`,
`P-45`). Dari **66 kemampuan** di `DAFTAR-PEKERJAAN.md`; `P-01`…`P-59` milik Treaty In **belum
ditiketkan**.

**Dua belas irisan dari sepuluh kemampuan** — `P-65` berubah-lebar dan dipecah tiga.

## Ukuran selesai di papan ini

Bukan "dapat diperagakan di layar". `L-4` menyatakan tidak ada spesifikasi layar di mana pun, dan
lubang itu milik **batch lapisan aplikasi**. Irisan di sini tegak melalui **lapisan yang ada**:

```
skema  →  constraint dan invarian  →  jalur migrasi  →  uji negatif DAN positif
```

---

## Daftar irisan

| # | Judul | Dari | Blocked by | Yang diantarkan | Penahan luar papan | Menahan / menunggui |
|---:|---|---|---|---|---|---|
| **01** | Versi baru lahir dari versi berlaku terakhir, dan rujukan dasarnya disimpan eksplisit | `P-42` **diubah** | None | Membuat versi penyesuaian menghasilkan baris ber-`ID_VERSI_KONTRAK_DASAR` terisi; versi pertama menolak terisi; menunjuk versi `DITOLAK` atau `DIBATALKAN` **ditolak** | `REV-3` | **menunggui** — yang berubah **keadaan mana yang dihitung berlaku**, bukan bahwa ada penunjuk dasar |
| **02** | Materialitas mengunci ruas yang boleh disunting, dan penolakannya ditegakkan di sisi simpan | `P-60` BARU | None | Versi bermaterialitas *Non Material* **menolak simpan** yang mengubah ruas uang; *Material* menerimanya. **PEMBUAT PERTAMA** `SIFAT_MATERIAL_ADDENDUM` | `DB-20` | **menunggui** — `DB-20` menentukan **titik bekunya**, dan titik beku milik irisan 03. Penguncian dan penegakannya tidak berubah |
| **03** | Jenis dan materialitas beku sejak diajukan, dan tiap perubahan meninggalkan jejak | `P-04` **diubah** | 02 | Mengubah jenis atau materialitas pada versi `DRAFT` **berhasil dan berjejak**; pada versi yang sudah diajukan **ditolak**, dan pesannya menyebut pengembalian ke `DRAFT` | `DB-20` | **MENAHAN** — titik beku adalah **pokok** irisan ini. Bila `DB-20` menjawab "beku sejak lahir", kriteria selesainya berbalik |
| **04** | Tanggal berlaku pada versi ditolak bila jatuh di luar periode kontrak | `P-45` **diubah** | None | Tanggal berlaku versi di luar periode kontrak ditolak; tanggal **pada dokumen** dinyatakan **bukan bagian irisan ini** | — | — |
| **05** | Dokumen addendum berdiri bernomor sendiri, dan satu dokumen menyentuh beberapa kontrak | `P-61` BARU | None | `DOKUMEN_ADDENDUM` ada; nomor ganda **ditolak** (`INV-71`); satu dokumen menunjuk versi di beberapa kontrak sekaligus, dan **persetujuan tetap per kontrak**. **PEMBUAT PERTAMA** | `DB-16a` | **MENAHAN** — bila dokumen ternyata dikirim ke luar, jenis bernilai **tiga** dan penanda "dikirim ke luar" kembali. Itu kolom, bukan penyetelan |
| **06** | Dokumen addendum membawa tanggal berlakunya sendiri, boleh kosong | `P-62` BARU | 05 | Dokumen dapat menyimpan tanggal berlakunya; kosong diterima; tidak dicampur dengan tanggal berlaku versi | `DB-16b` | **MENAHAN** — tanggal itu **adalah** pokok irisan ini, dan `KTV-2` menuntut kolomnya **dicabut sebelum data masuk** bila dibantah |
| **07** | Nomor dokumen baris warisan diisi dari arsip kertas | `P-66` BARU | 05 | **PK** dapat mengisi nomor dokumen pada baris warisan; yang belum terisi terlihat sebagai **daftar**, bukan sebagai diam. Migrasi **tidak** mengisinya — ia tidak punya sumber | `DB-16a` | **menunggui** — bentuk kolomnya ikut irisan 05; pekerjaan pengisiannya tidak berubah apa pun jawabannya |
| **08** | Baris selisih berkunci bisnis, dan mata uang ada di dalam kuncinya | `P-63` BARU | 01 | Menyimpan versi penyesuaian menghasilkan baris `NILAI_SELISIH` per besaran yang berubah, dipadankan **kunci bisnis** bukan posisi baris; baris disisipkan di tengah daftar **tidak** menggeser padanannya. **PEMBUAT PERTAMA** `NILAI_SELISIH` | — | — |
| **09** | Perubahan persentase reinstatement menghasilkan baris selisih | `P-64` BARU | 08 | Mengubah persentase pemulihan limit memunculkan baris selisih; versi yang **hanya** mengubahnya terbaca **Material** | — | — |
| **10** | *Perluas* — kolom nomor urut versi baru berdiri berdampingan dengan pengenal warisan | `P-65` BARU | None | Kolom nomor urut baru ada dan kosong; pengenal `/Rnn` tetap; **tidak ada** pembaca yang berubah. Hijau | — | — |
| **11** | *Pindahkan* — nomor urut warisan diisi menurut kronologi, per batch kontrak | `P-65` BARU | 10 | Batch kontrak warisan memperoleh nomor urut dari tanggal komentar pembuatan; baris yang **tidak dapat diurutkan** menjadi **pengecualian migrasi bernomor**, dilaporkan per kontrak. Hijau batch demi batch | — | — |
| **12** | *Kerutkan* — tidak ada lagi pembaca yang menurunkan urutan dari pengenal | `P-65` BARU | 11 | Turunan "versi berlaku" dibaca **hanya** dari nomor urut; membacanya dari pengenal `/Rnn` tidak ada lagi di mana pun | — | — |

### Hitungan

| | |
|---|---:|
| irisan | **12** |
| **MENAHAN** | **3** — 03, 05, 06 |
| menunggui, boleh jalan | 3 — 01, 02, 07 |
| **Yang sebenarnya dapat dimulai** | **4** — 01, 02, 04, 10 |

> **Empat dari dua belas, dan itu angka yang menggambarkan keadaan** — bukan "aktif 12". Tiga
> tertahan pertanyaan bisnis yang **belum dikirim**; lima menunggu tetangganya di papan.

---

## Tiga irisan dari PERLU DIUBAH — bunyi lamanya, supaya yang bergeser terlihat

`Menggantikan:` pada ketiganya menyebut **bunyi lama kemampuannya**, bukan hanya kode cacat. Yang
sudah membaca versi lama perlu tahu apa yang bergeser, dan kode cacat tidak memberitahunya.

### Irisan 01 — `P-42`

| | |
|---|---|
| **Bunyi lama** | *"**PK** dapat membuat versi baru **dari versi yang sudah `DISETUJUI`**, dan versi baru itu mulai dari `DRAFT` serta melewati keempat tingkat"* |
| **Yang bergeser** | dasarnya menjadi **versi berlaku terakhir**, dan rujukannya **disimpan eksplisit** |
| **Kenapa** | sebuah kontrak dapat punya **beberapa** versi `DISETUJUI` dalam sejarahnya. Bunyi lama mengizinkan memilih yang bukan terakhir — dan selisih yang dihitung terhadap dasar yang salah **tidak menghasilkan galat** |

### Irisan 03 — `P-04`

| | |
|---|---|
| **Bunyi lama** | *"**PK** dapat mengisi dan mengubah seluruh kepala versi **selama versinya `DRAFT`**"* |
| **Yang bergeser** | pembekuan sesudah `DRAFT` **dinyatakan**, bukan disiratkan; dan perubahan jenis maupun materialitas **meninggalkan jejak** |
| **Kenapa** | bunyi lama diam tentang apa yang terjadi sesudah `DRAFT`, dan diam sama sekali tentang jejak. Sistem lama mengubah jenis **kapan saja tanpa jejak** — pertanyaan *"jenis apa versi ini ketika kepala seksi menyetujuinya"* **tidak dapat dijawab** untuk satu pun baris warisan |

### Irisan 04 — `P-45`

| | |
|---|---|
| **Bunyi lama** | *"**Tanggal berlaku addendum** yang jatuh di luar periode kontraknya **ditolak**"* |
| **Yang bergeser** | pembawanya **disebut** — *"tanggal berlaku **pada versi**"* |
| **Kenapa** | `KTV-2` memberi dokumen tanggal berlakunya sendiri. Tanpa menyebut pembawanya, irisan 04 dan irisan 06 **sama-sama mengklaim ruas bernama sama**, dan yang menang adalah yang ditulis lebih dulu |

---

## Kewajiban yang melekat pada irisan tertentu

### Irisan 02 — `INV-69` dan `INV-70` berstatus KLAIM, bukan penegakan

Keduanya ditegakkan lewat *materialized view* ber-`REFRESH ON COMMIT`, dan **belum pernah dijalankan
di Oracle mana pun** (`L-3`).

> **MV yang gagal me-refresh berhenti menegakkan tanpa satu galat pun.**
>
> Irisan yang hanya memasang constraint-nya akan **dinyatakan selesai padahal tidak menegakkan apa
> pun** — dan tidak ada yang akan tahu sampai data yang melanggar sudah masuk.

Maka kriteria selesai irisan 02 memuat **tiga**, bukan satu:

| # | Kriteria |
|---|---|
| 1 | constraint terpasang |
| 2 | **uji negatif dijalankan** dan gagal sebagaimana seharusnya — bukan diargumentasikan |
| 3 | **pemantau kebasian** berdiri: `REFRESH_MODE`, `STALENESS`, terjadwal, dan **berbunyi kepada seseorang** |

Butir 3 menyentuh lubang yang lebih luas daripada irisan ini: **induk sudah memakai MV untuk
`INV-47`, `INV-50`, dan `INV-51`, dan tidak satu pun `P-nn` menyebut pemantau kebasiannya.** Dicatat
di `ADJUDIKASI-KEMAMPUAN-ADJUSTMENT.md` §3.1 sebagai **temuan**, menunggu satu kalimat pemilik
proses.

### Irisan 08 — uji POSITIF untuk share fakultatif

`GRL-16` memutuskan tidak ada kemampuan yang dibangun untuk share fakultatif: ia terbawa sendiri
oleh bentuk `NILAI_SELISIH`. **Pernyataan rancangan yang tidak pernah diperiksa berperilaku persis
seperti kemampuan yang tidak pernah disambungkan** — yaitu seperti cacat yang `GRL-16` sendiri
temukan.

Maka irisan 08 membawa uji positif: **mengubah share fakultatif harus menghasilkan sedikitnya satu
baris selisih.** Uji yang tidak pernah gagal atas besaran itu berarti besaran itu **tidak ada di
dalam himpunan yang dipadankan**.

### Irisan 09 — PERUBAHAN perilaku yang dilihat penyetuju di hari pertama

Versi yang **hanya** mengubah persentase reinstatement terbaca **Non Material** di sistem lama dan
**Material** di sistem baru. Itu disengaja — sistem baru menangkap yang lama lewatkan.

Kriteria selesainya memuat **pemberitahuan pra-peralihan**, bukan hanya perilakunya. Perubahan
klasifikasi yang muncul tanpa diberitahukan akan dibaca sebagai kerusakan.

### Irisan 11 — pengecualian migrasi adalah KELUARAN, bukan kegagalan

Baris warisan yang kronologinya tidak dapat ditentukan **tidak dipaksa masuk urutan**. Ia menjadi
pengecualian bernomor, dilaporkan **per kontrak**, dan kontraknya tidak memperoleh turunan "versi
berlaku" sampai seseorang memutuskannya.

Dan satu hal yang wajib tertulis di irisan ini atau uji paritas akan menandai keberhasilan sebagai
kegagalan: **`NOMOR_URUT_VERSI` dikecualikan dari kriteria identik ADR-0043.** Ia mekanisme yang
**diberikan**, bukan angka yang **dipindahkan**; yang tunduk pada kriteria identik adalah pengenal
`/Rnn`.

---

## Yang TIDAK menjadi irisan, dan sebabnya

| Hal | Kenapa bukan irisan |
|---|---|
| logika dua prosedur yang pindah ke Golang | `U-1` — tidak ada pelaku di luar tiket. Ia **ikut** irisan pertama yang memerlukannya, dan irisan itu menyatakannya di **YANG TEGAS BUKAN BAGIAN TIKET INI** milik tetangganya |
| kelas cacat share fakultatif (`GRL-16`) | tidak ada yang dibangun. Ia menjadi **uji positif** di dalam irisan 08 |
| mesin pro rata (`GRL-15`) | **sengaja tidak dibangun**. Gagal `U-3` — tidak ada kalimat *"pelaku X dapat Y"* untuk sesuatu yang tidak dibangun. Ia masuk **YANG TEGAS BUKAN BAGIAN TIKET INI** pada irisan 04, ditulis sebagai **pernyataan keputusan**, bukan `TODO` |
| bentuk paket uang | **sifat** (`U-6`) — dibawa tiap irisan sebagai invarian; bentuk fisiknya dibangun **PEMBUAT PERTAMA** |

---

## Satu asumsi yang mungkin sudah dapat dicoret

`ASUMSI-CLEAR.md` asumsi **2** berbunyi: *`D-1` menunggu `F-2`*.

`2-to-spec/COCOK-SILANG-CACAH-ATRIBUT.md` kini menyatakan **`F-2` DITUTUP** — §10.2 cocok di 44
sesudah angka judulnya diturunkan dan pencariannya diselesaikan.

**Bila `D-1` memang hanya menunggu `F-2`, asumsi 2 dicoret** — dan ia dicoret, **tidak dihapus**,
dengan tanggal dan siapa yang menutupnya. Tidak saya lakukan sendiri: `D-1` milik sesi to-spec
induk, dan mencoretnya dari sini berarti menyatakan sesuatu selesai atas nama orang lain.
