# Adjudikasi kemampuan Adjustment terhadap `P-01` … `P-59`

**Tanggal:** 24 September 2026 · **Dijalankan sebelum** satu baris `P-60` ditulis
**Dasar:** `GRL-01` — satu model, satu spesifikasi, satu papan.

> ## KENAPA BERKAS INI ADA
>
> `DAFTAR-PEKERJAAN.md` induk sudah menyebut **addendum**, **versi**, **selisih**, **materialitas**,
> dan **dokumen**. Menambahkan `P-60` tanpa memeriksa menghasilkan **kemampuan kembar** — dan
> kemampuan kembar melahirkan dua tiket yang **kriteria selesainya dapat bertentangan**, sementara
> keduanya tampak sah di papan.
>
> Golongan **SUDAH ADA** adalah yang paling mudah dihilangkan dan paling berguna disimpan: ia bukti
> bahwa kemampuan itu **diperiksa dan sengaja tidak ditambahkan**, bukan terlewat.

## Hitungan

| Golongan | Jumlah |
|---|---:|
| calon diadili | **12** |
| **SUDAH ADA** | **2** |
| **PERLU DIUBAH** | **3** |
| **BARU** | **7** |

> **Tujuh dari dua belas bergolongan BARU, dan itu tinggi.** Perintahnya berbunyi: *"bila BARU
> mendekati jumlah calon, kemungkinan besar adjudikasinya tidak dijalankan — periksa ulang."*
>
> **Diperiksa ulang, dan ketujuhnya bertahan.** Sebabnya satu dan dapat diperiksa: enam dari tujuh
> lahir dari keputusan yang **belum ada** ketika `DAFTAR-PEKERJAAN.md` ditulis 24 September pagi —
> `GRL-14` sampai `GRL-20`, `KTV-1` sampai `KTV-4`, `INV-69` sampai `INV-71`, dan `TDA-18`. Yang
> ketujuh, pemantau kebasian, adalah **lubang induk** yang baru terlihat dari sini.
>
> Dan dua dugaan saya sendiri **gugur** dalam pemeriksaan ini — lihat C1 dan C6. Itu tanda
> adjudikasinya berjalan, bukan distempel.

---

## 1. SUDAH ADA — diperiksa, sengaja tidak ditambahkan

### C2 — jenis addendum dinyatakan dan dibetulkan selama `DRAFT` → **`P-04`**

`P-04` berbunyi: *"**PK** dapat mengisi dan mengubah **seluruh kepala versi** selama versinya
`DRAFT`."*

`JENIS_ADDENDUM` adalah atribut kepala `VERSI_KONTRAK` (§10.2 induk). **"Seluruh kepala versi"
memuatnya.** Menambah baris baru untuk "menyatakan jenis" akan menghasilkan dua tiket yang keduanya
mengklaim penyuntingan kepala versi selama `DRAFT` — dan kriteria selesai keduanya akan menyentuh
ruas yang sama.

**Yang tidak dicakup `P-04`** — pembekuan pada `AJUKAN` dan jejaknya — bukan baris baru melainkan
**perubahan pada `P-04`**. Lihat C3.

### C9 — atribut *"berlaku sejak"* pada versi → **`P-45`**

`P-45` berbunyi: *"Tanggal berlaku addendum yang jatuh di luar periode kontraknya **ditolak**."*

Sebuah kemampuan yang **memvalidasi** tanggal berlaku addendum mengandaikan tanggal itu **ada dan
dapat diisi**. `GRL-15` tidak menambahkan ruas baru; ia menegaskan bawaannya = tanggal mulai
kontrak.

> **Dan ketiadaan mesin pro rata BUKAN kemampuan.** Ia gagal `U-3` — tidak ada kalimat *"pelaku X
> dapat Y"* yang dapat ditulis untuk sesuatu yang sengaja tidak dibangun — dan gagal `U-1`, sebab
> tidak ada pelaku di luar tiket yang menikmatinya.
>
> Perlakuannya: ia masuk medan **YANG TEGAS BUKAN BAGIAN TIKET INI** pada tiket `P-45`, ditulis
> sebagai **pernyataan keputusan** — apa yang sengaja tidak dilakukan, kenapa, akibatnya pada angka,
> dan kapan ia ditagih (`DB-5`). **Bukan sebagai `TODO`.**

---

## 2. PERLU DIUBAH — nomornya tetap, bunyinya bergeser

### C1 — `P-42`: dasar sebuah versi baru

**Dugaan saya semula: BARU. Gugur.**

| | |
|---|---|
| **Bunyi lama** | *"**PK** dapat membuat versi baru dari versi yang sudah `DISETUJUI`, dan versi baru itu mulai dari `DRAFT` serta melewati keempat tingkat"* |
| **Kenapa berubah** | `GRL-10` menetapkan dasar = **versi berlaku terakhir**, bukan "versi yang sudah `DISETUJUI`" mana pun. Bedanya nyata: sebuah kontrak dapat punya beberapa versi `DISETUJUI` dalam sejarahnya, dan memilih yang bukan terakhir **menghasilkan selisih terhadap dasar yang salah** |
| **Bunyi baru** | *"…dari **versi berlaku terakhir**, dan rujukan dasarnya **disimpan eksplisit**; versi pertama tidak punya dasar, versi penyesuaian wajib punya"* |
| **Berubah dari** | *"dari versi yang sudah `DISETUJUI`"* |
| **Asal** | `GRL-10`, bahan to-spec `B-3`, `B-4` |

### C3 — `P-04`: pembekuan pada `AJUKAN`, dan jejaknya

| | |
|---|---|
| **Bunyi lama** | *"…dapat mengisi dan mengubah seluruh kepala versi **selama versinya `DRAFT`**"* |
| **Kenapa berubah** | Bunyi lamanya **menyiratkan** pembekuan sesudah `DRAFT` tanpa menyatakannya, dan tidak menyebut jejak sama sekali. `GRL-18` dan `KTV-1` menetapkan **jenis dan materialitas beku sejak `AJUKAN`** — lebih awal daripada `P-38` yang membekukan pada `DISETUJUI` — dan tiap perubahan meninggalkan baris jejak (ADR-0045) |
| **Bunyi baru** | *"…dan perubahan atas **jenis** maupun **materialitas** meninggalkan baris jejak; sesudah diajukan keduanya **beku**, dan mengubahnya menuntut versinya dikembalikan ke `DRAFT` lebih dulu"* |
| **Berubah dari** | *"dapat diubah selama `DRAFT`"*, yang diam tentang apa yang terjadi sesudahnya |

> **Kenapa ini bukan baris baru:** membuat `P-60` untuk "pembekuan" menghasilkan tiket yang
> kriteria selesainya **menolak** hal yang kriteria selesai `P-04` **mengizinkan**, dan tidak ada di
> papan yang menunjukkan keduanya bicara tentang ruas yang sama.

### C6 — `P-45`: tanggal berlaku ada di DUA pembawa, dan itu harus disebut

**Dugaan saya semula: BARU. Gugur — tetapi ia memunculkan tabrakan yang lebih penting.**

`P-45` sudah bicara tentang *"tanggal berlaku addendum"*. `KTV-2` menetapkan **`DOKUMEN_ADDENDUM`
membawa tanggal berlaku sendiri, boleh kosong**.

> **Maka "tanggal berlaku" kini punya DUA pembawa** — versi dan dokumen — dan `P-45` tidak menyebut
> yang mana. Dibiarkan begitu, tiket `P-45` dan tiket dokumen akan **sama-sama mengklaim ruas
> bernama sama**, dan yang menang adalah yang ditulis lebih dulu.

| | |
|---|---|
| **Bunyi baru** | *"Tanggal berlaku **pada versi** yang jatuh di luar periode kontraknya ditolak — **dan tanggal berlaku pada dokumen addendum adalah ruas yang berbeda**, §`KTV-2`"* |
| **Berubah dari** | *"tanggal berlaku addendum"* tanpa pembawa |

---

## 3. BARU — `P-60` … `P-66`

| Kode | Kemampuan | Gol | Asal | Kenapa ia BUKAN kembar |
|---|---|---|---|---|
| **`P-60`** | **PK** dapat menyatakan **materialitas** sebuah versi, dan materialitas itu **mengunci ruas mana yang boleh disunting** — dua arah, dan penolakannya ditegakkan **di sisi simpan**, bukan hanya di layar | **B** | `GRL-20`, `INV-69`, `INV-70` | tidak ada `P-nn` yang menyebut materialitas sebagai masukan. `GRL-12` yang menjadikannya turunan **BATAL** |
| **`P-61`** | **PK** dapat mencatat **dokumen addendum** bernomor sendiri, dan **satu dokumen dapat menyentuh beberapa kontrak** — persetujuannya tetap per kontrak | **B** | `GRL-19`, `INV-71` | **bukan `P-15`.** `P-15` melampirkan berkas pada satu kontrak, dirujuk tidak disalin. Ini **entitas dengan pengenal bisnisnya sendiri** yang berdiri **di atas** versi dan melintasi kontrak |
| **`P-62`** | Dokumen addendum dapat membawa **tanggal berlakunya sendiri**, boleh kosong | **B** | `KTV-2` | pembawa yang berbeda dari `P-45` — lihat C6 |
| **`P-63`** | **PK** melihat **baris selisih per besaran yang berubah**, dipadankan lewat **kunci bisnis termasuk mata uang**, bukan menurut posisi baris | U | ADR-0048 butir 3, `E-1a`, `GRL-14` | **bukan `P-44`.** `P-44` menampilkan **nilai versi sebelumnya berdampingan**; ini **baris selisih tersimpan dan beku**, yang materialitas dibaca darinya |
| **`P-64`** | Perubahan **persentase reinstatement menghasilkan** baris selisih | U | `TDA-18` | **PERUBAHAN perilaku yang disengaja.** Sistem lama menyalin, tidak mengurangi — 10 penugasan, **nol pengurangan**, kalibrasi lulus |
| **`P-65`** | **PM** dapat memindahkan addendum warisan dengan **nomor urut diberikan ulang menurut kronologi**, sementara pengenal `/Rnn` **dilestarikan apa adanya** | U | `GRL-17` | **pecahan payung `P-50`**, sejajar `P-51` dan `P-52` yang sudah berdiri sebagai baris sendiri |
| **`P-66`** | **PK** dapat mengisi **nomor dokumen** baris warisan **dari arsip kertas** | **B** | `GRL-19` §3.1 | **pekerjaan orang, bukan sistem.** Migrasi tidak punya sumber — sapuan ekspor mengembalikan **nol** properti penyimpan nomor dokumen |

### 3.1 Satu BARU yang ternyata lubang INDUK, bukan lubang Adjustment

`P-60` menuntut **pemantau kebasian** karena `INV-69` dan `INV-70` ditegakkan lewat *materialized
view* yang **belum pernah dijalankan** (`L-3`).

Saat memeriksa apakah kemampuan itu sudah ada, ditemukan hal yang lebih luas:

> **Induk sudah memakai *materialized view* untuk `INV-47`, `INV-50`, dan `INV-51` — dan tidak ada
> satu pun `P-nn` yang menyebut pemantau kebasiannya.**
>
> *MV* yang gagal me-refresh **berhenti menegakkan tanpa satu galat pun**. Aturan induk sendiri
> menuntut uji negatif **dan** pemantau kebasian; yang kedua tidak pernah jadi kemampuan.

**Tidak saya tambahkan sebagai `P-67`**, karena ia **bukan kemampuan Adjustment** dan menambahkannya
dari sini berarti memutuskan sesuatu untuk daftar induk di luar lingkup yang diperintahkan.

| | |
|---|---|
| **Golongan** | **temuan**, bukan kemampuan |
| **Siapa menutup** | pemilik proses — satu kalimat menentukan apakah ia baris `P-nn` tersendiri atau sifat yang dibawa tiap tiket ber-*MV* |
| **Yang menagih** | berkas ini, dan kriteria selesai `P-60` yang akan menyebut pemantau kebasian tanpa punya nomor kemampuan untuk dirujuk |

---

## 4. Yang diperiksa dan TIDAK menghasilkan perubahan

Dicatat supaya tidak diperiksa ulang tanpa sebab:

| Dugaan | Hasil |
|---|---|
| daur hidup versi dan perpindahan keadaan | **`P-37`** sudah menolak perpindahan di luar tiga belas yang sah, termasuk keluar dari `DISETUJUI` dan `DITOLAK`. `GRL-08` tidak menambah keadaan baru di luar yang sudah ADR-0055 muat |
| rantai persetujuan empat tingkat | **`P-31`…`P-33`, `P-40`** — ADR-0052 milik induk, tidak berubah |
| pengenal versi, rantai tanpa percabangan | **`P-42`** sesudah diubah (C1); tidak butuh baris kedua |
| versi berlaku sebagai turunan | `GRL-11` menetapkan **tidak ada kolom penunjuk**. Sebuah turunan tanpa kolom **tidak punya kemampuan tersendiri** — ia dibaca, dan yang membacanya `P-56` (pencarian) dan `P-54` (pembaca hilir) |
| penolakan versi addendum | **`P-58`** sudah ada, ditambahkan 24 Sep |
| pembatalan draf | **`P-57`** sudah ada |
| kemampuan yang menyebut materialitas dan perlu diubah karena `GRL-12` batal | **nol** — tidak satu pun dari 59 baris menyebut materialitas di kalimat kemampuannya. Diperiksa, bersih |
| logika dua prosedur yang pindah ke Golang | **bukan kemampuan** (`U-1`) — tidak ada pelaku di luar tiket. Ia **ikut** kemampuan pertama yang memerlukannya, dan tiket itu menyatakannya di **YANG TEGAS BUKAN BAGIAN TIKET INI** milik tetangganya |
| kelas cacat share fakultatif `GRL-16` | **bukan kemampuan** — tidak ada yang dibangun. Ia menjadi **kewajiban uji positif** di dalam `P-63`: mengubah share fakultatif **harus** menghasilkan baris selisih |
