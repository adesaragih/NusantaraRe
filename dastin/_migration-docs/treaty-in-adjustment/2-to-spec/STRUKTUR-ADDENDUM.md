# Struktur addendum — DELTA atas spesifikasi induk

**Tanggal:** 24 September 2026 · **Fase:** to-spec

> # BUKAN SUMBER KEBENARAN
>
> **Struktur kanonik ada di `../../treaty-in/2-to-spec/`.** Berkas ini **hanya memuat yang khas
> addendum**.
>
> `GRL-01` mengikat: **modul ini tidak punya spesifikasi sendiri.** Addendum **adalah**
> `VERSI_KONTRAK`. Dua salinan penuh akan berbeda dalam sebulan, dan pembacanya tidak akan tahu mana
> yang menang.
>
> Hasil akhirnya ditulis ke **lima artefak induk**; folder ini menerima **catatan kerja dan usulan
> diff**. Artefak induk disunting **hanya sesudah diff-nya disetujui** — usulannya di
> `USULAN-DIFF-KE-INDUK.md`.

---

## 1. Keempat pohon menjadi apa — termasuk yang sengaja hilang

**Setiap yang hilang punya baris keputusan.** Tidak satu pun lenyap tanpa tercatat.

### 1.1 `OLDDATA` — tidak ada tabel

| | |
|---|---|
| **Nasib** | **DIBUANG** — 142 simpul, nol kolom |
| **Penggantinya** | sisi lama **di-`SELECT` ke versi dasar** lewat `ID_VERSI_KONTRAK_DASAR` |
| **Dasar** | ADR-0048 butir 2 · `GRL-03` · `GRL-10` |
| **`SISI`** | **56-KHAS** — satu penulis, hanya ada di ekspor Adjustment |
| **Label** | **PERUBAHAN.** Berubah dari: *"nilai lama disalin ke pohon kembar di dalam `JSONDATA` yang sama"* |
| **Arah dampak bila salah** | bila versi dasar ternyata **dapat berubah** sesudah versi anaknya lahir, `SELECT` memberi angka yang berbeda dari yang dulu dipakai. **`INV-24` menutupnya**: versi terminal beku. Bila `INV-24` dicabut, keputusan ini dibuka kembali |

### 1.2 `ActualValue` — dibuang seluruhnya

| | |
|---|---|
| **Nasib** | **DIBUANG** — **108 jalur cermin tidak melahirkan satu kolom pun** |
| **Penggantinya** | premi aktual menjadi **nilai versinya sendiri**; selisihnya lahir terhadap versi dasar lewat mesin yang sama dengan besaran lain |
| **Dasar** | **`GRL-14`** |
| **`SISI`** | KEDUANYA — 19 penulis, 18 berkas bersama + 1 khas |
| **Label** | **PERUBAHAN.** Berubah dari: *"pada addendum premi, nilai aktual hidup di pohon terpisah dan tidak pernah masuk selisih"* |

> ### `TDA-16` hilang bersamanya — dan itu dinyatakan, bukan disimpulkan
>
> Di sistem lama `TreatyEDMDifferencePremium` beriterasi atas `TreatyIn.EGNPI` sementara pengisi
> menyunting `TreatyIn.ActualValue.EGNPI`, sehingga **selisih EGNPI pada addendum premi selalu
> nol** — angka yang menjadi alasan jenis addendum itu ada **tidak pernah muncul sebagai selisih**.
>
> Begitu premi aktual menjadi nilai versi, mengubahnya **menghasilkan baris `NILAI_SELISIH` dengan
> sendirinya**. **Tidak ada pengecualian yang perlu ditulis tangan** — dan pengecualian bertangan
> persis bentuk yang `GRL-20` dibuat untuk menghapus.

### 1.3 `ValueDifference` — menjadi tabel selisih tersendiri

| | |
|---|---|
| **Nasib** | **DISIMPAN** untuk besaran pokok; **TURUNAN** untuk agregat |
| **Penggantinya** | `NILAI_SELISIH`, berkunci **`KUNCI_PADANAN`** |
| **Dasar** | ADR-0048 butir 3 · `METODE` §8.3 |
| **`SISI`** | **IRISAN** — **16 dari 16 penulisnya di berkas yang ada di kedua ekspor** |

> **Mesin selisih milik Treaty In, bukan Adjustment.** Maka `TDA-04`, `TDA-05`, dan `TDA-16`
> **temuan induk** — diadili di sini untuk jalur addendum saja, dan dicatat sebagai **usulan untuk
> induk**.

### 1.4 `ValueBeforeProrate` — tidak ada

| | |
|---|---|
| **Nasib** | **DIBUANG** — 15 simpul, nol kolom |
| **Dasar** | **`GRL-15`** |
| **Yang dibawa** | atribut **"berlaku sejak"** (`TANGGAL_BERLAKU_ADDENDUM`) |
| **Yang sengaja TIDAK dibangun** | **mesin pro rata.** Lima langkah pro rata **mati** di sistem lama |
| **Label** | **PERUBAHAN**, dan penghilangan ini **disengaja** |
| **Arah dampak bila salah** | bila pro rata dibutuhkan, ia **kemampuan BARU** yang dibangun di atas atribut yang sudah ada — bukan penggalian ulang. Atributnya sengaja dibawa justru untuk itu |

---

## 2. Tabel selisih — dan kenapa rumusnya berubah

Rumus lama dibaca dari ekspor: **58 penugasan pengurangan, 58 dari 58** terhadap `TreatyIn.OLDDATA.*`.
**`EVIDENCED(TreatyEDMDifference*@ekspor-2026-09)`**

| Sistem lama | Sistem baru | Dasar |
|---|---|---|
| sisi lama = salinan `OLDDATA` di `JSONDATA` yang sama | **`SELECT` ke versi dasar** | ADR-0048 butir 2 |
| dipadankan menurut **posisi** — `<CURRENT>` 49×, `local.subscript` 16× | **`KUNCI_PADANAN`** | `SEAM-ADJUSTMENT` §3 · `E1` |
| **mata uang disalin, tidak dibandingkan** | mata uang **masuk ke dalam kunci** | `E-1a` · ADR-0053 |
| **8 hasil pengurangan ditulis ke `ActualValue.LimitShareSummaryList`** | **tidak dibawa** | keputusan pemilik proses 22 Sep 2026 |

### 2.1 Akibat yang harus ditulis terang supaya tidak dikira cacat

> Karena **mata uang masuk ke dalam kunci padanan**, mengubah mata uang sebuah baris akan tampil
> sebagai **satu baris dihapus dan satu baris ditambah** — bukan sebagai selisih nilai.
>
> **Itu perilaku yang benar dan disengaja**: selisih antara jumlah berdenominasi berbeda memang
> tidak bermakna (ADR-0053). Tetapi ia akan terlihat aneh bagi pembaca yang mengharapkan satu baris
> berubah, dan karena itu ditulis di spesifikasinya, bukan ditemukan sendiri.

### 2.2 `ReinstatementPct` — TEMUAN, dan sapuannya dikalibrasi

**Sapuan:** seluruh korpus, 708 berkas, kedua ekspor, dua bentuk penulis.
**Kalibrasi:** `Limit` — besaran yang **pasti** dikurangi — ditemukan **6 pengurangan nyata**
terhadap `OLDDATA`. **LULUS**, jadi nolnya bermakna.

| Besaran | Penugasan | Berupa pengurangan terhadap `OLDDATA` |
|---|---:|---:|
| `Limit` *(kalibrasi)* | 78 | **6** |
| **`ReinstatementPct`** | **10** | **0** |
| **`ReinstatementValue`** | **4** | **0** |

Di `TreatyEDMDifferenceLimits`, bentuknya **penyalinan**, bukan pengurangan:

```
TreatyIn.ValueDifference.Limits(<CURRENT>).ReinstatementPct  =  .ReinstatementPct
```

> ### Maka: perubahan persentase reinstatement TIDAK PERNAH menghasilkan baris selisih
>
> **`EVIDENCED(TreatyEDMDifferenceLimits@ekspor-2026-09)`** · **`SISI` — IRISAN**, sebab
> `TreatyEDMDifferenceLimits` ada di kedua ekspor.
>
> Akibatnya menyentuh `GRL-20`: di bawah `INV-69`, versi **Non Material** dilarang punya baris
> selisih bertipe uang atau porsi. Persentase reinstatement **adalah porsi** — dan di sistem lama ia
> **tidak pernah** menghasilkan baris selisih. Maka versi yang **hanya** mengubah reinstatement akan
> terbaca **Non Material** di sistem lama, dan **Material** di sistem baru.
>
> **Itu perubahan perilaku, dan ia disengaja** — sistem baru menangkap perubahan yang sistem lama
> lewatkan. Dicatat sebagai **`TDA-18` calon**, `SISI` **IRISAN**, dan **diusulkan ke induk**; ia
> **belum melewati ronde TDA mana pun**, sama seperti `TDA-17`.

### 2.3 `UA-18` — keunikan `KUNCI_PADANAN` pada data warisan

`SEAM-ADJUSTMENT.md` §3 **menetapkan bentuk** kunci padanan. Ia **tidak menjamin** data warisan
mematuhinya.

> **`UA-18`.** Untuk setiap entitas anak, berapa daftar di `JSONDATA` yang memuat **dua baris dengan
> kunci padanan yang sama**? Bentuk majemuk ikut diuji — untuk `DETAIL_PROPORSIONAL` berarti (nomor
> layer + bagian layer) + kelompok treaty; untuk `POTONGAN` ditambah penunjuk induk mana.

**Bila ada yang tidak unik, pemadanan selisih akan memasangkan baris yang salah** pada baris
warisan — dan hasilnya angka selisih yang salah, bukan kegagalan yang terlihat.

---

## 3. `DOKUMEN_ADDENDUM` — entitas BARU di atas versi

`GRL-19`. **Bukan `PENYESUAIAN` yang kembali:** `PENYESUAIAN` ada **di bawah** versi; dokumen **di
atasnya**. Penelusuran 19 properti (`METODE` §8.3) yang menolak `PENYESUAIAN` **tetap berdiri**.

| Yang ditetapkan | Isi |
|---|---|
| **Label** | **BARU** — sapuan nomor dokumen **nol**, terkalibrasi: 28 nama, lima bentuk penulis, 708 berkas, kedua ekspor |
| **Relasi** | **dokumen (1) → (N) versi**, lintas kontrak |
| **Persetujuan** | **tetap melekat pada versi**, per kontrak satu per satu |
| **Rujukan dari `VERSI_KONTRAK`** | **boleh kosong** (`GRL-19` butir a); `DB-16a` menyempitkannya |
| **Tanggal berlaku dokumen** | kolom **nullable** (`KTV-2`); **dicabut bila `DB-16b` dibantah** |
| **Penanda "dikirim ke luar"** | **tidak ada** (`KTV-3`); satu nilai enum ditambahkan bila `DB-16a` dibantah |
| **Keunikan nomor** | **global**, diuji `UA-19` — dan arah salahnya **menolak data yang sah** |
| **Nama** | **`DOKUMEN_ADDENDUM`**, bukan `DOKUMEN_KONTRAK` yang sudah dipakai induk untuk lampiran |

### 3.1 Akibat migrasi — masuk RENCANA PERALIHAN, bukan spesifikasi

> **Seluruh baris warisan akan berkolom dokumen KOSONG.** Migrasi **tidak punya sumber** untuk
> mengisinya: `TD-01` membuktikan sistem lama tidak pernah merekam nomornya. Pengisiannya
> **pekerjaan orang dari arsip kertas**.

> ### Arti kalimat itu — DIBACA BEGINI, bukan begitu (syarat 2, 24 Sep 2026)
>
> | Bacaan | Akibat pada `NOMOR_DOKUMEN NOT NULL` |
> |---|---|
> | **✅ BENAR** — tabel `DOKUMEN_ADDENDUM` **tidak punya satu pun baris warisan**; yang kosong adalah **`VERSI_KONTRAK.ID_DOKUMEN_ADDENDUM`**, bernilai `NULL` | **aman** |
> | ❌ salah — barisnya **ada**, nomornya kosong | **pecah saat migrasi** |
>
> **Migrasi tidak menyisipkan satu pun baris `DOKUMEN_ADDENDUM`.**
>
> Tanpa pernyataan ini, orang berikutnya akan **mencabut `NOT NULL`** dari `NOMOR_DOKUMEN` — dan
> **kunci alaminya ikut hilang**, sehingga `INV-71` tidak dapat ditegakkan dan dokumen kehilangan
> satu-satunya pengenal yang dipakai manusia. Kerusakannya tidak terlihat sampai ada dua dokumen
> bernomor sama. Ditulis juga di `SPEC-MODEL-DATA.md` §10.19b.

**Ia harus muncul di rencana peralihan, bukan hanya di sini.** Kolom baru yang kosong untuk seluruh
warisan adalah hal yang **direncanakan**, bukan efek samping.

### 3.2 `GRL-01` dikoreksi, bukan dibongkar

Jumlah entitas **bertambah satu**. Pernyataan *"addendum adalah `VERSI_KONTRAK`"* **tetap benar** —
dokumen bukan addendum; ia **pembungkus** beberapa addendum.

**Selaraskan `STRUKTUR-DATA.md` induk lebih dulu** — ia yang **mengikat** soal daftar entitas
(urutan wewenang butir 4). Usulan diff-nya di `USULAN-DIFF-KE-INDUK.md`.

---

## 4. Materialitas adalah MASUKAN, dan ia sakelar DUA ARAH

`GRL-20` menggantikan `GRL-12`. **Alasannya:** pada saat materialitas dipakai, **akibatnya belum
ada** — sehingga ia tidak dapat diturunkan dari akibat.

> ### JANGAN pakai rumusan lama
>
> *"Non Material berarti uang tidak berubah"* hanya menyebut **satu arah**. Kedua pilihan mengunci
> himpunan yang **saling lepas**:

| Pilihan | Yang **tidak boleh** berubah |
|---|---|
| **Material** | teks kesepakatan — `Exclusions`, `SpecialConditions`, `TreatyContractName`, `ContractRefNo` |
| **Non Material** | besaran uang dan porsi — **428 sel bernama** |

**`EVIDENCED(Section@ekspor-2026-09)`** — hitungannya `GRILL-D/01-TEMUAN.md` `TD-02`.
**Angka "220 kondisi di 18 seksi" DILARANG dipakai** (`MA-13`): ia tidak dapat direproduksi.

| Yang mendarat | Bentuknya |
|---|---|
| `SIFAT_MATERIAL_ADDENDUM` | atribut **masukan** pada §10.2 — **kembali**; bagian `GRL-18` yang menghapusnya **gugur** |
| Titik beku | beku sejak `AJUKAN`, berjejak selama `DRAFT` (`KTV-1`); `DB-20` menyempitkan |
| **`INV-69`** | versi **Non Material** tidak boleh punya baris `NILAI_SELISIH` bertipe uang atau porsi |
| **`INV-70`** | versi **Material** tidak mengubah teks kesepakatan dan identitas kontrak terhadap versi dasarnya |

### 4.1 Uji — negatif DAN positif

| # | Yang disisipkan | Harus |
|---|---|---|
| `N-7a` | versi Non Material + satu baris selisih premi | **GAGAL** |
| `N-7a+` | versi Non Material + satu baris selisih **tanggal pelaporan** | **BERHASIL** |
| `N-7b` | versi Material + perubahan `PENGECUALIAN` | **GAGAL** |
| `N-7b+` | versi Material + perubahan limit, pengecualian tidak disentuh | **BERHASIL** |
| **`N-7c+`** | **versi warisan yang melanggar keduanya** | **BERHASIL DIMUAT** |

> **`N-7c+` bukan kelengkapan formal.** Tanpa ia, invariannya **menolak data lama pada hari
> peralihan**. `UA-3` yang mengukur berapa banyak — kepalanya sudah ditulis ulang.

---

## 5. `TDA-17` — empat punya rumah, tiga adalah lubang INDUK

| Kolom datar lama | Rumah di skema baru |
|---|---|
| `DEDUCTIBLE` | `LAYER.DEDUCTIBLE` |
| `MDP` | `LAYER.MDP` |
| `MDP_PCT` | `LAYER.PERSEN_MINIMUM_DEPOSIT` |
| `ADJ_RATE` | `LAYER.PERSEN_PENYESUAIAN` |
| **`DEDUCTIBLE2`** | **TIDAK ADA** |
| **`PREMIUM_EARNED`** | **TIDAK ADA** |
| **`ROL_PCT`** | **TIDAK ADA** |

> ### PENANDA YANG TIDAK DAPAT DISALAHBACA
>
> **Ketiganya BUKAN sengaja dibuang.** Tidak ada keputusan mana pun yang membuangnya.
>
> Mereka **tidak punya rumah** karena `L-8`: ketiganya muncul di peta telusur **hanya di dalam pohon
> cermin**, tidak pernah di pohon utama, dan dua ada di daftar titik buta bertanda `TIDAK`. **§10
> tidak melewatkannya — §10 tidak pernah diperlihatkan kepadanya.**
>
> **Pemiliknya sesi to-spec INDUK.** Tidak dapat diputuskan di modul Adjustment.
>
> Tanpa penanda ini, pembaca berikutnya akan menyimpulkan **pembuangan yang tidak pernah diputuskan
> siapa pun**.

---

## 6. Kalimat untuk pembaca arsip lama

`JSONDATA` menjadi arsip tanpa jalur baca (ADR-0034). **Kalimat berikut disimpan bersama arsipnya:**

> **`ActualValue` di dalam `JSONDATA` BERLAPIS** — `ActualValue(ActualValue(…))` — dan kedalamannya
> bertambah satu tiap kali addendum dibuat. **Pelapisan itu tumbuh pada jenis 1 (Internal) dan 2
> (External), BUKAN pada addendum premi:** prasyarat `EDMState=="3"` pada langkah 2
> `SaveTreatyIn_EDM_Act` **melompat** ke blok `jmp`, sehingga jenis 3 **melewati** langkah
> pemotretan — beserta dua langkah sesudahnya.
>
> Siapa pun yang membacanya tanpa mengetahui ini akan **membaca potret yang salah**.

**Dan arahnya penting:** berkas induk `4-erd-dan-tabel-datar/7-2-ACTUALVALUE.md` §1.2 pernah memuat
arah yang **terbalik**; koreksinya sudah diterapkan, tetapi arsipnya akan dibaca orang yang tidak
membuka berkas itu.
