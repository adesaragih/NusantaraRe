> Modul  : Treaty In Adjustment · paket usulan revisi ADR induk
> Dibuat : 2026-09-25 · sesi grilling, sesudah ronde A ditutup
> Status : **DRAF — ADR induk TIDAK disunting**
> Sifat  : TAMBAH-SAJA sampai diserahkan

# Usulan revisi ADR induk — paket `REV`

## Kenapa berkas ini ada, dan apa yang TIDAK dilakukannya

Grilling modul Adjustment membaca ekspor Pega yang sama yang mendasari ADR-ADR induk, dan pada
beberapa titik menemukan bahwa **pernyataan fakta di dalam sebuah ADR keliru** — bukan
keputusannya. `METODE` §3.7 melarang ADR memutuskan fakta tentang sistem lama; berkas ini
memperbaiki faktanya dan **tidak menyentuh satu pun keputusan**.

**Aturan main yang mengikat berkas ini:**

1. **ADR induk tidak disunting.** Tidak sekarang, tidak diam-diam, tidak sebagian.
2. **Draf ditulis begitu temuannya lahir**, bukan dikumpulkan di akhir (`METODE` §6.1) — karena itu
   berkas ini bersifat tambah-saja dan bertumbuh sepanjang grilling.
3. **Diserahkan sekaligus** kepada pemilik ADR **di akhir grilling**.
4. **Pengecualian tunggal:** bila `PP-1` menunjukkan sesi to-spec induk sedang berjalan, paket ini
   diserahkan **lebih cepat** — sebab menyunting spesifikasi di atas fakta yang keliru lebih mahal
   daripada mengganggu jadwal.
5. **Prasyarat cabang K:** *"paket `REV` sudah diserahkan dan ditanggapi"* masuk sebagai prasyarat
   masuk ke to-spec.

| # | ADR | Sifat usulan | Menyentuh keputusan? |
|---|---|---|---|
| `REV-1` | 0036 — beku saat disetujui | penulisan ulang **alasan** | tidak |
| `REV-2` | 0052 — satu rantai persetujuan | koreksi **Konteks** | tidak |
| `REV-3` | 0055 — daftar keadaan dan perpindahan | koreksi **empat pernyataan fakta** di §4 | tidak |
| `REV-4` | 0049 — jenis addendum satu sumbu | koreksi **sumber turunan materialitas**: dari jenis menjadi dari baris selisih tersimpan (GRL-12) | tidak |
| `REV-5` | 0049 — himpunan jenisnya | **menunggu C2.** Sesudah GRL-12, nilai "administratif" runtuh menjadi non-material turunan, sehingga himpunan tiga jenis ADR-0049 tidak lagi utuh apa pun jawaban C2 | ya — himpunan nilai enum |

---

## `REV-1` — ADR-0036: alasannya benar, tetapi memerikan sistem lama secara keliru

**Yang diusulkan berubah:** bagian **"Kenapa ini bukan sekadar kebersihan model"**. Keputusannya —
*"Angka yang menjadi dasar suatu persetujuan dibekukan pada saat persetujuan itu"* — dan kelima
penerapan yang mengikat **tidak disentuh**.

**Bunyi sekarang, kutipan persis:**

> ## Kenapa ini bukan sekadar kebersihan model
>
> Pemilik proses menetapkan bahwa Treaty In Adjustment akan **mengambil data lama lewat SELECT**,
> bukan menyalinnya. Itu hanya aman bila data lama tidak bergerak. Kalau nilai kontrak masih bisa
> berubah sesudah akseptasi, selisih yang dihitung hari ini dan bulan depan atas kontrak yang sama
> akan berbeda tanpa ada yang mengubah apa pun.
>
> Keputusan ini karena itu adalah **prasyarat teknis** bagi cara kerja Adjustment yang sudah
> ditetapkan, bukan pilihan rancangan yang bisa ditunda.

**Premisnya TIDAK dicabut, dan draf kedua `REV-1` keliru menyatakannya begitu.** **GRL-03 butir 2**
mempertahankan ADR-0048 butir 2 — *"data lama tidak disimpan ulang, diambil lewat SELECT"* — sebagai
**aturan yang mengikat pemilik proses** (`METODE` §8.5), dan melabelinya **PERUBAHAN** justru karena
sistem lama melakukan sebaliknya. Dengan versi yang beku, "mengambil nilai lama lewat SELECT"
berarti **membaca versi dasar yang tidak dapat berubah**. Rantai alasannya utuh.

**Yang benar-benar perlu diperbaiki hanya dua hal, dan keduanya soal pemerian, bukan penalaran:**

1. **Kalimat itu terbaca sebagai pemerian sistem lama, padahal ia arahan rancangan.** Pembaca
   berikutnya akan menyimpulkan bahwa modul Adjustment yang sekarang membaca-ulang lewat `SELECT`.
   **Ia tidak.** `TreatyInSetAddendumToHistory` menyalin seluruh halaman kontrak ke
   `TreatyIn.OLDDATA` (`Page-Copy` -> `Page-Copy` -> `Page-Remove`), dan salinan itu tersimpan di
   dalam `JSONDATA` baris addendum itu sendiri — **potret, bukan rujukan** (§2.1). **G2 sesi
   sebelumnya salah persis begitu.**
2. **Arah sebab-akibatnya perlu dibalik supaya tidak melingkar.** Bukan "SELECT menuntut
   pembekuan", melainkan **"pembekuanlah yang membuat SELECT aman"** — dan karena itu keputusan
   ADR-0036 berdiri lebih dulu, bukan sebagai pelayan sebuah pilihan bentuk.

**Usulan bunyinya:**

> ## Kenapa ini bukan sekadar kebersihan model
>
> Pemilik proses menetapkan bahwa Treaty In Adjustment **tidak menyimpan ulang data lama**,
> melainkan **mengambilnya lewat SELECT** ke versi dasarnya (ADR-0048 butir 2). Arahan itu hanya
> aman bila versi dasar tidak dapat bergerak — dan **pembekuan inilah yang membuatnya aman**, bukan
> sebaliknya. Kalau nilai sebuah versi masih bisa berubah sesudah akseptasi, selisih yang dihitung
> hari ini dan bulan depan atas versi yang sama akan berbeda tanpa ada yang mengubah apa pun.
>
> **Perlu dicatat supaya tidak disalahbaca: arahan "SELECT, bukan menyalin" adalah PERUBAHAN
> terhadap sistem lama, bukan pemerian sistem lama.** Sistem lama justru **menyalin**: addendum
> membawa potret penuh kontrak di dalam barisnya sendiri (`TreatyIn.OLDDATA`). Maka niat
> membekukan angka dasar persetujuan adalah **PELESTARIAN** — sistem lama sudah membekukannya —
> sedangkan **bentuk penyimpanannya** yang berubah.

**Yang tetap PERUBAHAN, dan sumbernya disebut satu per satu:**

| Bagian | Sumber perubahannya |
|---|---|
| bentuk simpan: rujukan ke versi dasar, bukan potret | **ADR-0048 butir 2**, dipertahankan **GRL-03 butir 2** |
| baris yang sudah disetujui tidak boleh ditimpa | **TDA-01** — tabrakan pengenal masuk cabang `UPDATE` dan menimpa baris lain sambil melapor berhasil |
| setiap angka ringkasan tersimpan | **TDA-06** — sebagian ditulis ke `ActualValue` lalu ditimpa pada penyimpanan berikutnya |
| baris addendum yang sudah disetujui tidak dapat disunting di tempat | **NA-19** — `Force Edit (dev)` membuka kuncinya dan `Save EDM(dev)` menyimpannya tanpa syarat status; keduanya terlihat operator divisi `IT`. Diukur `UA-16` |

**Satu hal yang sengaja TIDAK dicampur.** Paruh kedua keputusan ADR-0036 — *"angka yang
menggambarkan posisi terkini dihitung saat dibaca"* — mengatur **angka posisi**, bukan angka dasar
persetujuan. Ia tidak tersentuh usulan ini, dan tidak boleh dipakai sebagai argumen tentang
pembekuan.

**Dua sasaran sebelumnya dicabut, dan keduanya kesalahan saya.** Draf pertama menyasar rumusan
*"angka selisih dapat bergeser"* — itu **label GRL-02 buatan saya**, bukan teks ADR-0036. Draf kedua
menyatakan premis SELECT **sudah dicabut** — itu **bertentangan dengan GRL-03 butir 2**, yang
mempertahankannya.

**Menunggu:** `DB-1`. Bila penyetuju menyatakan selisih hanya bahan pertimbangan dan boleh dihitung
ulang, bagian PELESTARIAN gugur dan alasannya menyusut menjadi sumber-sumber PERUBAHAN saja.

---

## `REV-2` — ADR-0052: Konteksnya menyebut dua penyimpangan; yang berjalan satu

**Yang diusulkan berubah:** kalimat pertama bagian **Konteks**.

**Bunyi sekarang:**

> Sistem lama punya dua penyimpangan dari rantai empat tingkat: sebuah jalur alternatif lewat
> penerima tugas kelompok, dan sebuah jalur revisi yang memendekkan empat tingkat menjadi dua.

**Usulan bunyinya:**

> Sistem lama punya satu penyimpangan yang **berjalan** dari rantai empat tingkat — jalur revisi
> yang memendekkan empat tingkat menjadi dua. Ada pula sisa sebuah jalur alternatif lewat penerima
> tugas kelompok, tetapi **dalam versi aturan yang ter-ekspor** keadaan itu **tidak dapat
> dimasuki**: tidak ada aturan di kedua ekspor yang menulis
> `Position = "ReasTreatyInGroupLeader"`, dan cabang yang memindahkannya ke tingkat berikutnya
> dinonaktifkan (`Akseptasi_DT` 1.4, `pyDisabled = true`). Membuangnya karena itu **pelestarian
> terhadap perilaku hari ini**, bukan perubahan. Lihat ADR-0055 §4 butir 1, yang sudah menetapkan
> hal ini.
>
> Pernyataan ini **tidak** berlaku untuk versi aturan sebelumnya: cabang 1.4 dibuat 2019 dan
> dinonaktifkan kemudian, dan `TreatyInSetValue` langkah 2.6 masih menanam nama orang untuk peran
> itu. Bila data menunjukkan baris yang pernah melewatinya, yang benar bukan "tidak pernah ada"
> melainkan "tidak lagi dapat dimasuki".

**Sebabnya, dan kenapa ia ringan.** ADR-0055 §4 butir 1 **sudah** membuktikannya pada 23 September
2026 dan sudah menyatakan bahwa dasar ADR-0052 kini lebih kuat. Yang belum dikerjakan hanyalah
ADR-0052 sendiri masih berbunyi "dua penyimpangan". `REV-2` karena itu berbentuk **rujukan silang**,
bukan penyelidikan baru.

**Kenapa ini bukan kosmetik.** Pembaca berikutnya yang membaca "dua penyimpangan" akan mengira ada
jalur kelompok yang pernah dipakai, mencari datanya, lalu merancang pemetaan untuk sesuatu yang
tidak pernah berjalan (`METODE` §3.8, arah kebalikannya).

**Satu penajaman yang ikut diusulkan.** Konteks berbunyi *"Pemilihan jalur revisi dilakukan lewat
tombol"*. Itu benar untuk **kontrak**; untuk **addendum** nilainya **diwarisi** dari baris yang
dimuat, tidak dipilih siapa pun (§7.4, NA-12). Usulan tambahannya satu kalimat:

> Pada addendum, jalur pendek itu tidak dipilih siapa pun: nilainya diwarisi dari baris yang dimuat.

**Menunggu:** `DB-10` dan `UA-12`, dan keduanya tidak setara.

* **`DB-10`** mempersempit butir 4 — ia menanyakan kehendak bisnis di balik rantai pendek, bukan
  fakta tentang jalur kelompok.
* **`UA-12(a)`** dapat **membantah sebagian kalimat usulan di atas.** Bila ada baris yang
  `Position`-nya bernilai `ReasTreatyInGroupLeader` hari ini, maka keadaan itu **pernah dimasuki**,
  dan frasa "dalam versi aturan yang ter-ekspor" menjadi satu-satunya yang menyelamatkan kalimat
  itu dari salah. Dalam hal itu Konteks harus menyebut keadaan itu sebagai **keadaan warisan yang
  ada pada data** (ADR-0054), bukan hanya sebagai kode yang tidak terpakai.

---

## `REV-3` — ADR-0055: empat pernyataan fakta di §4 yang tidak cocok dengan ekspor

**Yang diusulkan berubah:** tabel **§4 "Tidak ada pintu yang membatalkan keadaan terminal"**.
Keputusannya — *"di model baru tidak ada satu pun jalan menyetel keadaan selain melalui perpindahan
di daftar"* — **tetap, seluruhnya**.

**Bunyi tabel sekarang:**

| Aturan | Yang dilakukannya | Terlihat oleh |
|---|---|---|
| `TreatyInForceEdit` | membuka kunci kontrak yang sudah disetujui | **setiap pengguna layar penawaran** |
| `TreatyInForceResolveComplete` | menyetel selesai disetujui tanpa approver | dua nama pengembang |
| `TreatyInReturntoInputor` | mengembalikan ke admin, status dikosongkan | dua nama pengembang |
| `TreatyInSetToDirector` | melompati dua tingkat | lewat `TreatyInTestAgent` di layar penawaran |

**Usulan tabel penggantinya:**

| Aturan | Yang dilakukannya | Terlihat oleh | Catatan |
|---|---|---|---|
| `TreatyInForceEdit` | membuka kunci (`ViewState = 0`) | **setiap operator berdivisi `IT`** | selnya `pyVisible = ALWAYS`, sehingga kondisi dua nama di sebelahnya **tidak dievaluasi**; gerbang divisi pada layout `S3` tetap berlaku |
| `TreatyInForceResolveComplete` | menyetel selesai disetujui tanpa approver, **dan menyimpan** (`pyActivity = SaveTreatyIn_Act`) | **operator divisi `IT` yang bernama `ALDO SAPUTRA` atau `Daniel Suhana`** | keterlihatan = wadah × sel |
| `TreatyInReturntoInputor` | mengembalikan ke admin, status dikosongkan | idem | tidak menyimpan |
| `TreatyInSetToDirector` | — **tidak melakukan apa pun** | — | **seluruh langkah tingkat atasnya ber-blok `//`**; langkah 4.1–4.7 bersarang di dalam langkah 4 yang mati. Dan seandainya hidup, 4.5 menuntut `Position == "ReasTreatyInGroupLeader"`, keadaan yang tidak dapat dimasuki |

**Empat butir usulan, terpisah supaya dapat diterima satu per satu:**

| # | Usulan | Dasar |
|---|---|---|
| **i** | `TreatyInSetToDirector` dipindahkan dari daftar pintu yang berjalan ke catatan "pintu mati" | seluruh langkah tingkat atas ber-blok `//` |
| **ii** | kalimat "layar penawaran" diperluas: `Section/TreatyInActionButtons` juga disertakan `Section/InputTreatyInAdjustment`, jadi **pintu samping ada di layar addendum juga** | empat kemunculan di badan masing-masing |
| **iii** | keterlihatan ditulis sebagai **hasil perkalian wadah dan sel**, per tombol seperti tabel di atas | sebaran `pyVisible`: 553 dari 553 `OTHER` menyimpan kondisi; hanya 37 dari 6663 `ALWAYS` |
| **iv** | ditambahkan satu baris: di layar addendum, `Force Resolve Complete(dev)` memanggil prosedur **kontrak** `PEGA_TREATY_IN`, yang masuk cabang `UPDATE … WHERE ID = IDPega`; pengenal addendum tidak ada di tabel kontrak, sehingga **tidak ada baris yang berubah dan tidak ada galat** — layar melapor *"Data Sudah Disimpan"*. Baris kontrak asalnya **tidak** tersentuh | §6.4, `PEGA_TREATY_IN.txt` |

**Satu hal yang TIDAK diusulkan berubah, dan sebabnya disebut.** §4 boleh tetap menyebut keempat
aturan itu sebagai "pintu samping", termasuk yang mati. Menghapus yang mati dari daftar akan
membuat pembaca berikutnya menemukannya lagi dan mengira ia terlewat. Yang diperlukan hanya
**menandainya mati**, bukan membuangnya (`METODE` §3.2).

---

## Bukan revisi ADR — temuan untuk modul induk (`IND`)

Bernomor supaya dapat dilacak (`METODE` §3.9). Ketiganya menyentuh **Treaty In**, bukan jalur
addendum, dan tidak satu pun sudah tercatat di dokumen induk mana pun.

| # | Isi | Calon pemilik | Kadar |
|---|---|---|---|
| **IND-1** | `Section/ShowSummary` memuat `pxTextInput` **dapat disunting** untuk kelima field lapisan beku, termasuk **`CedingID`** dan **`ProportionType`** yang tidak punya penyunting di layar mana pun lagi. **Tetapi layarnya dijaga satu nama operator** — `pyCondition = OperatorID.pyUserName = 'ALDO SAPUTRA1'`, berakhiran `1`, sekeluarga dengan pola `1=2`. Jadi ia **layar pengembang**, bukan jalur bisnis | ADR-0040, atau catatan di `PENGETAHUAN.md` induk | rendah — tetapi ia satu-satunya penyunting `CedingID` di korpus |
| **IND-2** | Keempat kontrol `SetTreatyIn_Act` **dibatasi peran inputor** lewat `pyUserData/pyCondition`, bukan lewat privilese; labelnya terbaca `Revision`, `Copy`, `View`, `Edit`; dan **hanya `Revision`** mengirim `revisionstate=1` | daftar eskalasi induk butir 1 — diff menunggu persetujuan | tinggi — ia mengoreksi teks yang sudah beredar |
| **IND-5** | **`pyDisabled` di `Section` menandai MODE, bukan kendali.** `SAPUAN-DAN-NAMA-TAGNYA.md` §4 menyebutnya *"hampir pasti kendali layar yang dimatikan"*. Di ekspor Adjustment ia duduk di `pyModes/rowdata` ber-`pxObjClass = Embed-Control-Mode`, dan **100% sel yang dijaga kondisi `ViewState` memilikinya** (96 dari 96) — sehingga ia tidak dapat berarti kendali mati. Angka 1.082 + 198 di Section dan Harness karena itu **bukan** jumlah kendali mati | `SAPUAN-DAN-NAMA-TAGNYA.md` §4 | tinggi — angka itu dipakai sebagai batas yang harus disebut |
| **IND-4** | **`BENTUK-PENULIS-PROPERTI.md` mendaftar empat bentuk penulis, seluruhnya di sisi aturan. Bentuk kelima ada di sisi antarmuka:** pengikatan sel — `pyValue` pada sel ber-`pyReadOnly != true` — yaitu penulisan yang terjadi saat pengguna mengetik. Bentuk itu yang memunculkan `NA-01`, `NA-17b`, dan analisis `ViewState` `NA-21`/`NA-22`. Perkakasnya `treaty-in-adjustment/tools/tulis.py`, jenis `BIND` | `BENTUK-PENULIS-PROPERTI.md` §2 | tinggi — tanpa bentuk kelima, "siapa yang dapat mengubah field ini" tidak terjawab |
| **IND-3** | **`Save EDM(dev)`** menyimpan addendum **tanpa syarat status apa pun**; bersama `Force Edit (dev)` ia memberi operator divisi `IT` jalan menyunting-lalu-menyimpan addendum yang **sudah disetujui**, tanpa versi baru dan tanpa persetujuan ulang | daftar eskalasi induk butir 1, dan `REV-1` sumber PERUBAHAN ketiga | tinggi |
