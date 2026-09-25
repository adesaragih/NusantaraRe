> Modul  : Treaty In Adjustment · pemilahan sisa grilling
> Dibuat : 2026-09-26, sesudah GRL-10 · **dilengkapi 26 September 2026**
> Status : **DITERAPKAN** — seluruh 43 baris anggaran asli digolongkan; anggaran baru berlaku
> Dasar  : `METODE` §7.1, §7.2, §6.4, §3.9, §3.2, §6.6

# Pemilahan sisa grilling

## 0. Dari mana daftarnya diambil

Tabel anggaran 43 pertanyaan **dibaca ulang dari catatan sesi**, bukan dari ingatan — pelajaran
`MA-09`. Sumbernya jawaban saya sendiri pada giliran A5, bagian *"Anggaran pertanyaan — lengkap"*,
yang memuat sebelas baris cabang beserta rincian butirnya. **Seluruh 43 butir muncul di bawah**,
masing-masing dengan **tepat satu** golongan.

Kriteria: grilling selesai bila **tidak ada lagi yang menghalangi penyusunan spesifikasi** — bukan
bila semua pertanyaan terjawab (`METODE` §7.1, §7.2).

**Koreksi definisi golongan (`PG-06`):** keputusan **wewenang** tidak masuk GRILL. Ia milik orang
berwenang (`METODE` §6.6); golongannya **DITUNDA**, pemiliknya manajemen, tempatnya daftar eskalasi.

## 1. Seluruh 43 butir, digolongkan

### Cabang A — 4 butir · **selesai, ronde A**

| Butir | Golongan | Di mana |
|---|---|---|
| ADR-0040 | selesai | GRL-05 |
| ADR-0049 | selesai | GRL-06 |
| ADR-0052 | selesai | GRL-07 |
| ADR-0055 | selesai | GRL-08 |

### Cabang B — 5 butir · **selesai, ronde B**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| identitas versi | **KONFIRMASI** | ADR-0040 butir 1 — "addendum atas" menjadi hubungan versi |
| nomor tampilan `/Rnn` | selesai | **GRL-09** |
| bentuk rantai | selesai | **GRL-10**; bentuk simpannya KONFIRMASI lewat `STRUKTUR-DATA.md` §1.1 |
| addendum ganda | **KONFIRMASI** | `SPEC-INVARIAN.md` INV-25 — paling banyak satu versi tak-terminal per kontrak |
| penjaga tabrakan | **REKOMENDASI TO-SPEC** | sudah ditulis sebagai bahan to-spec **B-1** dan **B-2** di `GRILL-B/06-PUTUSAN.md` |

### Cabang C — 4 butir · **3 GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| **C1 — definisi materialitas** | **GRILL** | menentukan apakah materialitas atribut tersimpan, turunan jenis, atau hasil pemeriksaan. Gerbang bagi E |
| **C2 — daftar jenis**, mencakup **penyesuaian premi** | **GRILL** | menentukan himpunan nilai `JENIS_ADDENDUM`. **Dua baris anggaran asli digabung** menjadi satu pertanyaan: "penyesuaian premi" adalah pertanyaan tentang apakah jenis ketiga berdiri sendiri, yaitu isi pertanyaan daftar jenis. Bila pemilik proses menghendaki keduanya terpisah, anggaran GRILL menjadi 8 |
| **C3 — arti bisnis `ActualValue`** | **GRILL** | pohon ketiga di `JSONDATA`: dibawa sebagai atribut, dibuang, atau menjadi turunan — menentukan jumlah entitas |

### Cabang D — 5 butir · **tidak ada GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| daftar keadaan | **KONFIRMASI** | ADR-0055 + GRL-08 |
| rantai per jenis | **DITUNDA** | menunggu `DB-10`; bila dibantah ia menjadi usulan revisi ADR-0052 butir 4, bukan pertanyaan grilling |
| penolakan | **KONFIRMASI** | ADR-0055 §4.2, §4.3, dan perubahan 24 Sep (syarat a — baris tidak dihapus) |
| penarikan sesudah diajukan | **DITUNDA** | ADR-0055 perubahan 24 Sep menyatakannya **sengaja tidak diputuskan**: *"ia pertanyaan tersendiri dengan jawabannya sendiri"*. Pemiliknya bisnis; butir daftar untuk dibantah **`DB-13`** |
| jejak | **KONFIRMASI** | ADR-0045 — jejak perubahan sebagai fakta mesin |
| *(tambahan ronde A)* penyalaan "versi menggantikan" | **REKOMENDASI TO-SPEC** | **R-D1** di §5 |
| *(tambahan ronde A)* tahap peringatan lapisan beku | **KONFIRMASI + DITUNDA** | lihat §5, `R-D2` **dicabut** |
| *(tambahan ronde A)* pemberitahuan penyetuju dan pembuat | **DITUNDA** | menunggu `UA-13`, `UA-9(a)` |
| *(tambahan ronde A)* wewenang tombol "(dev)" | **DITUNDA** | `PG-06` — wewenang milik manajemen; sudah di daftar eskalasi butir 1 |

### Cabang E — 5 butir · **2 GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| ~~**E1 — kunci padanan per daftar**~~ | **KONFIRMASI** *(turun 26 September 2026)* | `SEAM-ADJUSTMENT.md` §3 sudah menetapkan **`KUNCI_PADANAN`** — identitas bisnis yang stabil lintas versi — dan tambahan 24 September memisahkan **Bentuk A** (sembilan entitas, kunci alaminya cukup) dari **Bentuk B** (`DETAIL_PROPORSIONAL`, `POTONGAN`, `PENYEBARAN` — kunci padanan **induknya** lebih dulu), beserta ruas tambahan bagi `POTONGAN` yang berinduk dua. Itu persis isi E1. **Sisa yang belum tertulis** menjadi bahan to-spec **E-1**, bukan pertanyaan |
| mata uang (TDA-05) | **KONFIRMASI** | ADR-0053 — bentuk daftar, tiap mata uang dikonversi sekali, `CurrencyRelation` tidak menggerakkan angka. **Penegakannya menumpang E1**: mata uang masuk kunci padanan |
| presisi (TDA-14) | **REKOMENDASI TO-SPEC** | tipe angka adalah pengisian bentuk; `SPEC-MODEL-DATA` sudah menyediakan kolom tipe |
| **E3 — kemampuan yang mati atau hilang: dibangun atau dibuang** | **GRILL** | satu pertanyaan, **tiga baris keputusan terpisah** — masing-masing dapat dijawab berbeda: |
| E3a — **pro rata** (`EDMEffective` -> `TreatyCalculateProratePct`) | — | bila dibangun, **"berlaku sejak" menjadi atribut versi** — perubahan bentuk. Bila dibuang, addendum selalu berlaku sejak tanggal mulai kontrak (`DB-5`) |
| E3b — **share fakultatif** | — | `TreatyEDMDifferenceDeduction` langkah 4 dan 5 ber-blok `//`, berlabel *"Facultative share not enable yet in adjustment"* — kemampuan yang **belum pernah menyala**. Menyalakannya `METODE` §3.8: kemampuan **BARU**, bukan pelestarian |
| E3c — **ringkasan yang tidak pernah tersimpan** (TDA-06) | — | sebagian selisih ditulis ke `ActualValue` lalu ditimpa. Membangunnya berarti menyimpan angka yang di sistem lama hilang — **PERUBAHAN**, dan sumbernya sudah tercatat di `REV-1` |
| aturan induk yang dipakai ulang | **REKOMENDASI TO-SPEC** | mengisi bentuk yang sudah ada; tidak ada entitas baru |
| bentuk tabel selisih | **KONFIRMASI** | `METODE` §8.3 — bentuk sempit, menggantung pada `VERSI_KONTRAK` |
| **mekanisme pembekuan** — bekukan hasil beserta pemadanannya, atau versikan rumus | **KONFIRMASI** *(ditambahkan 27 September 2026)* | **Butir ini sebelumnya titipan tanpa golongan** — ia tidak ada di tabel anggaran 43, dan pemilik proses benar menuntutnya diberi tempat. Jawabannya: `STRUKTUR-DATA.md` §1.1 memberi `NILAI_SELISIH` kunci alami **"besaran + kunci padanan baris"**, berinduk `VERSI_KONTRAK`. Jadi **kunci padanan ikut tersimpan pada baris selisih itu sendiri**, bukan hanya hidup di bentuk baca `NILAI_VERSI_KONTRAK`. Maka mekanismenya **"bekukan hasil beserta pemadanannya"** — sejalan GRL-02 (angka dasar persetujuan beku) dan **GRL-12 syarat 1** (materialitas dibaca dari baris tersimpan). Perubahan cara penurunan bentuk baca kelak **tidak** dapat mengubah pemadanan yang sudah disetujui |

### Cabang F — 4 butir · **tidak ada GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| bentuk simpan — relasional atau JSON | **KONFIRMASI** | `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` adalah model **relasional** dengan entitas dan atribut bernama; ADR-0034 menutup jalur baca arsip JSON. **Butir ini sempat hilang dari daftar saya dan digantikan diam-diam oleh "bentuk simpan cedant" — itu keliru, dan keduanya kini berdiri sendiri** |
| *(tambahan ronde A)* bentuk simpan cedant | **KONFIRMASI** | `STRUKTUR-DATA.md` mendaftar **`CEDANT [luar]`** sebagai entitas luar yang dirujuk `KONTRAK`. "Rujukan saja" sudah diputuskan; **F1 turun dari GRILL** |
| penggandaan (TDA-15) | **KONFIRMASI** | GRL-01 + `STRUKTUR-DATA.md` §1.1 — `OLDDATA` tidak dibawa sebagai pohon kedua; sisi lama dibaca lewat `ID_VERSI_KONTRAK_DASAR` |
| prosedur `POOLDATA` | **REKOMENDASI TO-SPEC** | jalur tulis adalah pengisian bentuk |
| penguncian | **REKOMENDASI TO-SPEC** | tidak mengubah entitas; **R-F1** di §5 |

### Cabang G — 3 butir · **tidak ada GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| sumber peran (TDA-09) | **KONFIRMASI** | ADR-0044 — entitas peran dan penugasan bertanggal masuk gelombang 1; *"Nama orang tidak pernah muncul di dalam aturan"* |
| jalur lolos `IT Developer` | **KONFIRMASI** | idem — keadaan dan peran terpisah tegas, setiap larangan punya tepat satu sebab |
| pemisahan pembuat dan penyetuju | **DITUNDA** | kebijakan bisnis, bukan bentuk; butir daftar untuk dibantah **`DB-14`** |

### Cabang H — 4 butir · **tidak ada GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| tiga kolom | **REKOMENDASI TO-SPEC** | susunan layar; tidak mengubah entitas |
| komponen `OldData` | **REKOMENDASI TO-SPEC** | idem |
| picker (TDA-11) | **KONFIRMASI** | INV-25 membuat keadaan yang dipilih picker lama mustahil di model baru |
| sumber aturan editabilitas | **KONFIRMASI** | ADR-0038 — aturan yang bisa berubah disimpan sebagai data. **Isinya** datang dari C1, bukan bentuknya |

### Cabang I — 4 butir · **2 GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| **I1 — versi berlaku untuk kontrak warisan** | **GRILL** | menentukan apakah `KONTRAK` butuh penunjuk versi berlaku atau ia turunan. **Diperiksa: tidak ada artefak induk yang memutuskannya** — frasa "versi berlaku" tidak muncul di `STRUKTUR-DATA.md`, `ERD.md`, maupun ADR mana pun |
| **I2 — aturan baru mana yang boleh dijalankan atas baris warisan** | **GRILL** | menentukan ada-tidaknya kolom asal-usul dan berapa banyak. GRL-06 hanya menutup materialitas. Menampung titipan `B-3`/`B-4` warisan dan `ID_VERSI_KONTRAK_DASAR` untuk baris warisan |
| strategi peralihan | **REKOMENDASI TO-SPEC** | urutan kerja, bukan bentuk |
| ambang ketidakcocokan | **REKOMENDASI TO-SPEC** | angka ambang; **R-I1** di §5 |
| penanganan temuan data | **DITUNDA** | menunggu `UA-1`, `UA-3`, `UA-10` |
| riwayat yang terhapus | **DITUNDA** | menunggu `UA-7`; mekanismenya sudah hilang (ADR-0055 syarat a) |

### Cabang J — 3 butir · **tidak ada GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| lampiran disalin atau dirujuk | **KONFIRMASI** | `STRUKTUR-DATA.md` §1.2 — `DOKUMEN_KONTRAK`: *"rujukan ke satu dokumen lampiran — **dirujuk, tidak dimiliki**"* |
| jalur kategori | **REKOMENDASI TO-SPEC** | pengisian bentuk |
| nasib lampiran saat ditolak | **KONFIRMASI** | ADR-0055 perubahan 24 Sep syarat (a) — baris tidak dihapus, sehingga lampiran yang merujuknya tetap sah |

### Cabang K — 2 butir · **tidak ada GRILL**

| Butir | Golongan | Alasan / sumber |
|---|---|---|
| prasyarat to-spec | **DITUNDA** | termasuk *"paket `REV` diserahkan dan ditanggapi"* |
| keputusan induk yang harus final | **DITUNDA** | pemiliknya pemilik ADR |

## 2. Anggaran baru

**Enam butir GRILL** *(turun dari tujuh, 26 September 2026 — E1 sudah dijawab induk)*, menggantikan batas atas 40.

| Cabang | Butir GRILL |
|---|---|
| C | C1 ✅ GRL-12, C2, C3 |
| E | E3 (tiga baris keputusan: E3a, E3b, E3c) |
| I | I1 ✅ GRL-11, I2 |

Titik periksa sesudah butir kelima. Penutupan sesudah keenam.

**Terpakai dua** — `I1` (GRL-11) dan `C1` (GRL-12). **Sisa empat:** `C2`, `C3`, `E3`, `I2`.

### Bahan to-spec yang lahir dari penurunan E1

| # | Bunyi | Sebabnya |
|---|---|---|
| **E-1** | `KUNCI_PADANAN` sebuah entitas yang membawa mata uang **menyertakan mata uang itu** | **TDA-05**: mesin selisih lama mengurangkan tanpa memeriksa mata uang. `SEAM-ADJUSTMENT.md` §3 menyebut nomor layer dan pengenal pihak, **tidak menyebut mata uang** — dan tanpa itu dua baris bermata uang berbeda dapat dipadankan |
| **E-1a** | **Akibat E-1 yang harus ditulis terang supaya tidak dikira cacat:** karena mata uang masuk ke dalam kunci, **mengubah mata uang sebuah baris akan tampil sebagai satu baris dihapus dan satu baris ditambah**, bukan sebagai selisih nilai. Itu perilaku yang benar dan disengaja — selisih antara jumlah berdenominasi berbeda memang tidak bermakna (ADR-0053) — tetapi ia akan terlihat aneh bagi pembaca yang mengharapkan satu baris berubah | konsekuensi langsung E-1 |

## 3. Perpindahan terhadap usulan sebelumnya

| Perpindahan | Arah | Alasan |
|---|---|---|
| **F1 bentuk simpan cedant** | GRILL -> **KONFIRMASI** | `CEDANT [luar]` sudah entitas luar yang dirujuk (`STRUKTUR-DATA.md` §3) |
| **F bentuk simpan relasional/JSON** | hilang -> **KONFIRMASI** | butir yang sempat lenyap dari daftar saya; sudah dijawab model induk |
| **J lampiran, G peran** | tidak pernah didaftar -> **KONFIRMASI** | ADR-0044, `STRUKTUR-DATA.md` §1.2 |
| **I1** | GRILL -> **GRILL** | diperiksa ke artefak induk; **tidak ada yang menjawabnya** |
| **C penyesuaian premi** | hilang -> digabung ke **C2** | pertanyaannya sama: apakah jenis ketiga berdiri sendiri |
| **E yang mati/hilang** | hilang -> digabung ke **E3** | satu pertanyaan bentuk: dibangun atau dibuang |

Semua perpindahan berada di dalam yang diizinkan pemilik proses: F1, J, dan G turun ke KONFIRMASI;
F-bentuk-simpan diperiksa ke induk dan ternyata sudah dijawab. **Tidak ada butir yang naik ke
GRILL.**

## 4. Urutan menurut ketergantungan

```
C1 ──gerbang──> C2        C1 ──gerbang──> E1 ──> E3
I1 ──>  I2      (mandiri terhadap C dan E)
C3 (mandiri)
```

**`C1`, `C2`, `E1`, dan `E3` menunggu `EXP-1`.** Sampai ekspor itu datang, urutannya dimulai dari
butir mandiri: **`I1` -> `I2` -> `C3`**. Bila `EXP-1` datang di tengah, C1 disisipkan pada gilirannya
berikutnya. **Tidak ada penebakan himpunan nilai.**

## 5. Butir REKOMENDASI TO-SPEC yang tidak punya ronde

Cabang D, F, H, I, dan J tidak punya ronde, sehingga rekomendasinya ditulis di sini — dengan bentuk
putusan: label, satu sampai tiga kalimat **kenapa begini** (`METODE` §6.7), dan catatan bahwa
keputusan akhirnya di to-spec. Butir yang jatuh di cabang beronde ditulis di berkas putusan rondenya
(B-1, B-2 di `GRILL-B/06-PUTUSAN.md`).

### R-D1 · Penyalaan "versi menggantikan" dilakukan sesudah I1 terkunci

* **Label** — **BARU**. Sistem lama tidak pernah mengganti kontrak dengan addendum yang disetujui
  (§6.3).
* **Rekomendasi** — dinyalakan **sesudah** migrasi baris warisan selesai dan `I1` memutuskan versi
  berlaku, bukan bersamaan dengannya.
* **Kenapa begini** — menyalakannya bersamaan berarti dua sumber perubahan bekerja pada baris yang
  sama di hari yang sama, dan bila hasilnya salah tidak ada cara memisahkan sebabnya.
* **Dan ia mengubah angka yang dibaca hilir** (`METODE` §4.5). Butir daftar untuk dibantah:

  > **`DB-12`.** *"Angka kontrak yang keluar ke akuntansi dan ke pihak lawan selama ini **tidak**
  > memperhitungkan addendum yang sudah disetujui — ia memakai angka kontrak aslinya. Mulai sistem
  > baru, angka itu akan memperhitungkannya. Siapa yang perlu tahu sebelum hari peralihan, dan
  > adakah laporan yang harus dihitung ulang?"*

  Dikaitkan dengan **I1**: yang menentukan angka mana yang keluar adalah definisi versi berlaku.
* **Keputusan akhir di to-spec.**

### ~~R-D2~~ · Tahap peringatan lapisan beku — **DICABUT**

Saya menulis *"ADR-0040 sudah memilih memperingatkan"*. **Salah, dan kutipannya membuktikannya.**
Frasa *"Ia **memperingatkan, tidak melarang**"* ada pada **butir 2**, tentang **kunci alami**.
Lapisan beku ada di **butir 3**, dan bunyinya larangan:

> **Lapisan beku addendum.** Untuk setiap addendum, lima hal tidak boleh berubah: cedant, source of
> business, `ProportionType`, tanggal mulai, dan tanggal berakhir. Perubahan atas salah satunya
> berarti kontrak lain, bukan addendum.

Maka:

* **Aturannya KONFIRMASI** — ADR-0040 butir 3 sudah memutuskan: perubahan atas kelimanya bukan
  addendum. Tidak ada rekomendasi yang perlu ditulis.
* **Yang tersisa bukan "kapan mulai memblokir", melainkan nasib baris warisan yang sudah
  melanggar** — dan itu **DITUNDA**, menunggu `UA-10`, milik cabang **I**.
* Bila kelak seseorang mengusulkan tahap peringatan, itu **penyimpangan dari ADR-0040 butir 3** dan
  harus diajukan sebagai usulan revisi ADR, bukan sebagai rekomendasi to-spec.

### R-F1 · Penguncian baris: optimistis, dengan penanda versi baris

* **Label** — **BARU**. Sistem lama tidak punya penguncian sama sekali: `PEGA_TREATY_IN` dan
  `PEGA_M_TREATY_IN_EDM` menjalankan `UPDATE … WHERE ID = …` tanpa memeriksa apakah baris berubah
  sejak dibaca.
* **Rekomendasi** — penguncian **optimistis**: setiap baris membawa penanda versi baris, dan
  penyimpanan yang penandanya tidak cocok **ditolak dengan pesan yang menyebut siapa yang mengubah
  lebih dulu**.
* **Kenapa begini** — penguncian pesimistis menahan baris selama formulir terbuka, dan formulir
  modul ini besar serta lama diisi. Optimistis menyelamatkan kasus yang justru terjadi di sistem
  lama: dua orang menyimpan kontrak yang sama dan yang terakhir menang tanpa ada yang tahu.
* **Keputusan akhir di to-spec.**

### R-I1 · Ambang ketidakcocokan disebut sebagai angka, bukan "sedikit"

* **Label** — **BARU**.
* **Rekomendasi** — migrasi berhenti dan dilaporkan bila baris yang tidak dapat dipetakan melebihi
  **satu persen** dari jumlah baris, atau bila **satu saja** kontrak yang berstatus disetujui tidak
  dapat dipetakan.
* **Kenapa begini** — ambang tanpa angka akan diputuskan diam-diam oleh orang yang menjalankan
  migrasi pada tengah malam. Dua ambang dipakai karena dua jenis kerusakan berbeda: banyak baris
  remeh yang gagal adalah masalah kualitas data; satu kontrak berlaku yang gagal adalah masalah
  bisnis.
* **Keputusan akhir di to-spec.**

## 6. Status cabang sesudah pemilahan

| Cabang | Status |
|---|---|
| **A** | **ditutup dengan ronde** — GRL-05 s.d. GRL-08 |
| **B** | **ditutup dengan ronde** — GRL-09, GRL-10 |
| **C** | **terbuka** — C1, C2, C3 |
| **D** | **ditutup tanpa ronde** — 2 KONFIRMASI, 1 REKOMENDASI (`R-D1`), 4 DITUNDA |
| **E** | **terbuka** — E1, E3 |
| **F** | **ditutup tanpa ronde** — 3 KONFIRMASI, 2 REKOMENDASI (`R-F1`, prosedur) |
| **G** | **ditutup tanpa ronde** — 2 KONFIRMASI, 1 DITUNDA (`DB-14`) |
| **H** | **ditutup tanpa ronde** — 2 KONFIRMASI, 2 REKOMENDASI |
| **I** | **terbuka** — I1, I2 |
| **J** | **ditutup tanpa ronde** — 2 KONFIRMASI, 1 REKOMENDASI |
| **K** | **ditutup tanpa ronde** — 2 DITUNDA |

## 7. Penutupan grilling

Grilling ditutup bila **ketujuh butir GRILL terkunci** dan setiap butir di tiga golongan lain sudah
punya tempat — yang sudah dipenuhi berkas ini. Sesudah itu disusun serah terima dengan pemisahan
terang antara **yang masih menghalangi** dan **yang hanya mengukur kerusakan**.
