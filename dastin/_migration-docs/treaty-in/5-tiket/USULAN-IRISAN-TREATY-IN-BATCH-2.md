# Usulan irisan — batch 2 Treaty In, daur hidup dan warisan

**Tanggal:** 25 September 2026 · **Langkah 2 to-ticket** · **Status: MENUNGGU PERSETUJUAN**
**Sifat:** usulan. **Nol berkas tiket ditulis.**

> ### PENAHANNYA UKURAN SESI, BUKAN `REV-3`
>
> Batch ini pernah dicatat *"menunggu tanggapan `REV-3`"* di tiga berkas. **Itu keliru** —
> `REV-3` menyatakan keputusan ADR-0055 **tetap seluruhnya**, dan `KTV-4` sudah mengadilinya:
> ***"nol dari enam menentukan letak kolom."*** Ralatnya bertanggal 25 September di
> `issues/README.md`, `ASUMSI-CLEAR.md`, dan `SERAH-TERIMA-TO-TICKET-BATCH-1.md` `K-1`.
>
> **Batch ini dapat dikerjakan sekarang.** Yang memisahkannya dari batch 1 hanya bahwa 51 kemampuan
> tidak muat satu jendela konteks.

---

## 1. Adjudikasi terhadap 44 tiket yang sudah ada

**Sumbernya sapuan mekanis**, bukan hitungan tangan: medan `*Asal:*` ke-44 berkas tiket dicocokkan
terhadap daftar `P-nn` di `DAFTAR-PEKERJAAN.md`.

| | |
|---|---:|
| kemampuan aktif | **62** |
| disebut tiket mana pun | **40** |
| **belum disebut satu tiket pun** | **22** |
| ditambah pecahan yang separuhnya sudah ditiketkan | **+2** — `P-01`, `P-50` |
| **lingkup batch 2** | **24 butir** |

### Hitungan tiga golongan

| Golongan | Jumlah |
|---|---:|
| **SUDAH DITIKETKAN** | **0** |
| **BERSINGGUNGAN** | **4** |
| **BARU** | **20** |

> **Nol yang sudah ditiketkan, dan itu bukan tanda adjudikasi tidak dijalankan.** Batch 1 memang
> dipilih sebagai **komplemen** batch 2 — kemampuan daur hidup disisihkan dengan sengaja. Yang
> membuktikan adjudikasinya berjalan adalah **empat yang bersinggungan**, dan ketiganya di antaranya
> tidak akan terlihat tanpa memeriksa isi tiketnya.

### 1.1 BERSINGGUNGAN — empat, dan bagian yang tersisa disebut

| Kemampuan | Tiket yang menyentuhnya | Apa yang **sudah** ada | Apa yang **tersisa** |
|---|---|---|---|
| **`P-01`** *(pecahan)* | `14` | **kolom** `KEADAAN_SIKLUS_HIDUP` berdiri — tabelnya tidak dapat berdiri tanpanya | **himpunan nilai sahnya** dan **mesin perpindahannya** |
| **`P-50`** *(pecahan)* | `44` | isi **enam tabel acuan** dipindahkan | pemindahan **kepala kontrak** dan seluruh **daftar anaknya** — dan kepala **menulis `KEADAAN_SIKLUS_HIDUP`** |
| **`P-38`** | `14` | `INV-24` **disebut** di tiket `14` | penegakannya **atas seluruh entitas anak** — trigger `BEFORE INSERT OR UPDATE OR DELETE`. Menyebut bukan menegakkan |
| **`P-27`** | `20` | kurs **tersimpan** pada `MATA_UANG_KONTRAK` beserta periodenya | **pembekuannya**: angka rupiah pada versi `DISETUJUI` tidak dibaca ulang saat ditampilkan |

> **`P-38` adalah contoh paling bersih kenapa golongan ini perlu ada.** `§0.2` batch 1 menyebutnya
> *"sudah tersentuh papan"*, dan sapuan medan `*Asal:*` menyebutnya **belum ditiketkan sama sekali**.
> **Keduanya benar**: tiket `14` menyebut `INV-24` sebagai invarian yang harus dipenuhi, dan tidak
> membangun penegakannya. **Tersentuh bukan tercakup.**

### 1.2 Dua kemampuan yang tidak masuk batch mana pun — `K-3`

`P-20` dan `P-49` **tidak muncul di batch 1 maupun batch 2** sampai sapuan ini dijalankan. Keduanya
`TERTAHAN` jawaban luar, dan keduanya **bukan** kemampuan daur hidup — sehingga tidak tersaring ke
batch 2 dan tidak tertinggal di batch 1.

**Kemampuan yang tidak masuk batch mana pun tidak akan pernah ditiketkan**, dan tidak ada pencacah
yang memperlihatkannya. Keduanya dimasukkan di sini.

---

## 2. Pertanyaan yang harus dijawab terang: daur hidup satu irisan atau dua puluh?

**Jawaban: satu irisan untuk MESINNYA, dan irisan tersendiri untuk tiap AKIBATNYA.**
Batch ini **18 irisan**, bukan 24 dan bukan 8.

### Kenapa mesinnya tidak boleh dipecah

Daftar keadaan dan daftar perpindahan **saling mengunci**: sebuah daftar putih berisi tiga belas
perpindahan adalah **satu artefak**, dan membangunnya bertahap berarti daftarnya **salah di setiap
langkah antara**.

> Daftar putih yang baru memuat lima perpindahan **menolak delapan perpindahan yang sah**. Ia akan
> **lulus setiap uji negatif** yang pernah ditulis untuknya — sebab uji negatif hanya memeriksa
> bahwa yang terlarang ditolak — dan **gagal hanya pada uji positif**.
>
> Itu persis bentuk kegagalan yang kepala papan ini peringatkan, dan ia akan lolos berminggu-minggu.

### Kenapa akibatnya HARUS dipecah

Tiga hal yang tampak bagian dari mesin ternyata **artefak tersendiri yang dapat gagal
sendiri-sendiri**:

| Bukan mesin | Ia sebenarnya |
|---|---|
| `P-38` terminal beku | **trigger pada entitas anak** — `UPDATE`/`DELETE` saja masih mengizinkan versi disetujui **bertambah** baris |
| `P-39` satu versi tak-terminal | **constraint keunikan bersyarat**, bukan baris di daftar perpindahan |
| `P-31`…`P-33` tiga tingkat | **siapa boleh memicu perpindahan mana** — wewenang, bukan daftar. ADR-0044: keadaan menjawab *"mungkinkah sekarang"*, peran menjawab *"bolehkah orang ini"*, dan **setiap larangan harus punya tepat satu sebab** |

**Memaksa ketiganya masuk irisan mesin** menghasilkan satu tiket yang tidak dapat dinyatakan selesai
tanpa menyelesaikan seluruh batch — dan itu melanggar `U-4`.

---

## 3. Daftar irisan — **18, nomor mulai 45**

| # | Judul | Dari | Blocked by | Yang diantarkan | Penahan luar papan | Menahan / menunggui |
|---:|---|---|---|---|---|---|
| **45** | Daftar keadaan dan perpindahan sah berdiri; perpindahan di luar daftar ditolak | `P-01` pecahan · `P-37` | `14` | Kolom keadaan hanya menerima delapan nilai sah; **tiga belas** perpindahan diterima dan **setiap yang lain ditolak**, termasuk setiap perpindahan keluar dari `DISETUJUI` dan `DITOLAK` | — | — |
| **46** | Versi terminal beku, dan pembekuannya menjangkau seluruh entitas anaknya | `P-38` | `45` | Mengubah nilai versi `DISETUJUI` ditolak; **menambah baris anak** pada versi terminal juga ditolak | — | — |
| **47** | Paling banyak satu versi tak-terminal per kontrak | `P-39` | `45` | Membuat versi kedua yang belum selesai ditolak, pesannya menyebut versi yang masih terbuka | — | — |
| **48** | Pengajuan menolak dua syarat dan memperingatkan enam | `P-29` · `P-30` | `45` | `K1-1` dan `K1-2` **menolak**; `K1-3`…`K1-8` menghasilkan **peringatan tercatat**, belum penolakan | — | — |
| **49** | Persetujuan tiga tingkat, tiap tingkat memindahkan ke antrian berikutnya | `P-31` · `P-32` · `P-33` | `45` · `48` | SH → antrian DH → antrian DR → `DISETUJUI`, dan tingkat terakhir hanya bila seluruh `K2` terpenuhi | — | — |
| **50** | Pengaju tidak menyetujui, tiap tingkat orang berbeda, tanpa jalan pintas | `P-36` · `P-40` | `49` | Penyetuju yang sama dengan pengaju **ditolak**; sama dengan tingkat sebelumnya **ditolak**; melompati tingkat **ditolak** | — | — |
| **51** | Pengembalian ke `DRAFT` dengan alasan wajib | `P-34` | `49` | SH/DH/DR mengembalikan versi ke `DRAFT`; pengembalian tanpa alasan **ditolak** | — | — |
| **52** | Penolakan versi — kontrak maupun addendum — meninggalkan catatan | `P-35` · `P-58` | `49` | Penolakan tanpa alasan ditolak; **barisnya tidak dihapus**, dan `DITOLAK` terminal | — | — |
| **53** | Pembatalan draf oleh pembuatnya sendiri | `P-57` | `45` | `DRAFT → DIBATALKAN` oleh pembuatnya; barisnya tetap tersimpan, **nomornya tidak dipakai ulang** | — | — |
| **54** | Tiap perpindahan meninggalkan satu catatan persetujuan, tanpa dapat dilewati | `P-46` | `45` | Setiap perpindahan menulis pelaku dan waktunya; **tidak ada jalur yang melewatinya** | — | — |
| **55** | Angka rupiah pada versi disetujui tidak bergeser, hari ini maupun tahun depan | `P-27` | `46` · `20` | Kurs **dibekukan pada versinya**, tidak dibaca ulang saat ditampilkan | — | — |
| **56** | Wewenang persetujuan dibatasi nilai kontrak | `P-41` | `49` | Penyetuju yang wewenangnya di bawah nilai kontrak **ditolak**, pesannya menyebut batasnya | **dokumen batas wewenang** | **MENAHAN** — angka ambangnya **adalah** pokok tiketnya, dan ia **tidak dapat dipecah** |
| **57** | Kepala kontrak warisan dipindahkan beserta keadaannya | `P-50` pecahan | `45` · `44` | Kepala kontrak dan daftar anaknya pindah **apa adanya**; keadaan warisan ditulis ke kolom yang sudah bernilai sah | `KTV-A` | **menunggui** — **lebar kolom**, dan ia **bertenggat**: irisan ini memuat data |
| **58** | Keadaan warisan tanpa padanan mendarat di `WARISAN_TAK_TERPETAKAN` | `P-51` | `57` | Nilai liar seperti `"test"` mendarat di keadaan bernama, **nilai aslinya tersimpan** dan tidak pernah dibaca perhitungan | — | — |
| **59** | Baris warisan diperbaiki ke keadaan sah yang dipilih eksplisit | `P-52` | `58` | Orang memilih keadaan sahnya; pilihannya **tercatat di jejak** beserta siapa dan kapan | — | — |
| **60** | Kontrak dicari, termasuk menurut keadaan siklus hidupnya | `P-56` | `45` | Pencarian menurut cedant, asal bisnis, periode, sifat proporsi, **dan keadaan** | — | — |
| **61** | Cara pembukuan XOL dicatat, dan hanya pada kontrak non-proporsional | `P-20` | `14` | `INV-34` menolak cara pembukuan XOL pada kontrak proporsional | **`Uji X-2`** | **MENAHAN** — uji itu menentukan apakah kolomnya **dipakai sama sekali** |
| **62** | Tabel acuan pembagian kapasitas punya pemilik yang bertanggung jawab | `P-49` | `15` | Tabel acuan kapasitas membawa pemiliknya, dan perubahannya berjejak | **penetapan pemilik** | **MENAHAN** — tanpa pemilik yang ditetapkan, tidak ada yang dapat dinyatakan selesai |

### Hitungan

| | |
|---|---:|
| irisan | **18** |
| dari kemampuan | **24** |
| **MENAHAN** | **3** — `56`, `61`, `62` |
| menunggui | 1 — `57` |
| **Yang sebenarnya dapat dimulai** | **NOL** |

### Kenapa NOL, dan kenapa angka itu benar kali ini

Pelajaran `K-2` diterapkan: untuk tiap irisan yang mengaku dapat dimulai, **sebutkan entitas yang
disentuhnya dan tunjukkan tiket mana yang membuatnya.**

| Irisan tanpa penghalang di dalam batch ini | Entitas yang disentuhnya | Tiket yang membuatnya |
|---|---|---|
| `45` | `VERSI_KONTRAK` | **`14`** — belum selesai |
| `61` | `VERSI_KONTRAK` | **`14`** — belum selesai |
| `62` | tabel acuan kapasitas | **`15`** — belum selesai |

**Tidak satu pun berdiri sendiri**, dan itu bukan kelemahan batch ini: ia akibat langsung dari
`14` dan `15` sebagai **PEMBUAT PERTAMA** seluruh papan.

> **Begitu `14` dan `15` mendarat**, yang dapat dimulai di batch ini menjadi **tiga** — `45`, `61`,
> `62` — dan begitu `45` mendarat, melompat ke **delapan**.

---

## 4. Yang TIDAK menjadi irisan

| Hal | Kenapa |
|---|---|
| perubahan lebar | **tidak ada di batch ini.** Penomoran ulang `NOMOR_URUT_VERSI` sudah dipecah perluas–pindahkan–kerutkan di batch Adjustment (`05`, `10`, `12`) |
| entitas peran dan penugasan bertanggal | `ADR-0044` menaruhnya di gelombang 1, dan **§10 tidak memuatnya sama sekali** — itu `F-16`, lubang spesifikasi, bukan irisan yang dapat ditulis |
| paket `REV-1`…`REV-6` | **utang dokumentasi ADR induk**, bukan utang tiket. Ia tidak menahan satu irisan pun |

---

## 5. Tiga pertanyaan

1. **Kekasarannya pas?** Yang paling mungkin Anda tolak: `49` menggabungkan **tiga** tingkat
   persetujuan menjadi satu irisan. Alasan saya: ketiganya satu rantai dengan satu bentuk, dan
   memisahkannya menghasilkan tiga tiket yang kriteria selesainya hampir identik. Alasan menolaknya
   sama kuatnya: `U-2` melarang satu tiket melintasi lebih dari **satu** perpindahan, dan `49`
   melintasi **tiga**. **Bila `U-2` ditegakkan harfiah, `49` pecah tiga dan batch ini menjadi 20.**
2. **Sisi penghalangnya benar?** Yang paling saya ragukan `55` — ia bergantung pada `46` *(terminal
   beku)* **dan** `20` *(mata uang berkurs)*, dan `20` ada di batch 1. Tepi ke luar batch itu persis
   yang keliru di `K-2`.
3. **Ada yang perlu digabung atau dipecah?** Kandidat gabung: `58`+`59` — keduanya tentang baris
   warisan tak terpetakan, dan yang kedua tidak berarti tanpa yang pertama. Kandidat pecah: `45`,
   bila Anda menghendaki daftar keadaan dan daftar perpindahan berdiri terpisah.
