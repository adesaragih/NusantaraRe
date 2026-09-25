# Tabel datar — Treaty In

**Tanggal:** 23 September 2026 · **Skema:** `TREATY_MASUK`

> **Tabel datar adalah bentuk PIPIH untuk dibaca dan dilaporkan, BUKAN sumber kebenaran.**

---

## 0. Aturan yang mengikat seluruh berkas ini

> ### TABEL DATAR TIDAK PERNAH DITULIS APLIKASI. IA DITURUNKAN.
>
> Satu fakta satu penulis (ADR-0041), dan penulisnya bukan tabel ini. Tidak ada jalur tulis ke
> tabel datar dari kode aplikasi, sekarang maupun nanti. Akun `TREATY_MASUK_APP` diberi hak
> **BACA saja** atasnya.

Sistem lama melanggar aturan ini, dan akibatnya sudah terlihat: `TREATYINDETAIL` diisi oleh langkah
di dalam jalur penyimpanan, sehingga ketika satu langkah itu dimatikan pada jalur addendum, tabelnya
diam-diam berhenti menerima addendum — tanpa satu pun galat.

---

## 1. `DATAR_VERSI_KONTRAK` — satu baris mewakili apa

> **SATU BARIS = SATU VERSI KONTRAK.**

Bukan satu kontrak, dan bukan satu versi per layer. Alasannya: versi adalah **satuan terkecil yang
punya keadaan persetujuan sendiri dan nilainya beku setelah disetujui** — ia satu-satunya butiran
yang tidak bergeser di bawah pembacanya.

Satu baris per kontrak akan memaksa memilih versi mana yang diwakilinya, dan pilihan itu berubah
tiap kali ada addendum. Satu baris per layer akan menggandakan seluruh nilai kepala.

### Kuncinya

| Kunci | Kenapa unik |
|---|---|
| `ID_VERSI_KONTRAK` | pengenal buatan sistem, unik seumur hidup sistem, tidak pernah dipakai ulang (INV-02, INV-03) |

Kunci alternatif yang **juga** unik dan lebih terbaca manusia: `ID_KONTRAK` + `NOMOR_URUT_VERSI`
(INV-04). Ia dicantumkan sebagai kolom, tetapi **bukan** kunci tabel — kunci tunggal lebih murah
dijoin dan tidak dapat terurai.

### Kolomnya

| Kolom | Dari entitas | Salin atau hitung |
|---|---|---|
| `ID_VERSI_KONTRAK` | `VERSI_KONTRAK` | salin |
| `ID_KONTRAK` | `KONTRAK` | salin |
| `NOMOR_URUT_VERSI` | `VERSI_KONTRAK` | salin |
| `ID_VERSI_KONTRAK_DASAR` | `VERSI_KONTRAK` | salin |
| `ADALAH_PENYESUAIAN` | `VERSI_KONTRAK` | **hitung** — benar bila `ID_VERSI_KONTRAK_DASAR` terisi |
| `NOMOR_KONTRAK_WARISAN` | `KONTRAK` | salin |
| `ID_CEDANT` | `KONTRAK` | salin |
| `NAMA_CEDANT` | `CEDANT` [luar] | **hitung** — dibaca dari master saat penurunan, tidak disimpan di model |
| `ID_ASAL_BISNIS` | `KONTRAK` | salin |
| `SIFAT_PROPORSI` | `KONTRAK` | salin |
| `TANGGAL_MULAI` · `TANGGAL_BERAKHIR` | `KONTRAK` | salin |
| `NAMA_KONTRAK` | `VERSI_KONTRAK` | salin |
| `KEADAAN_SIKLUS_HIDUP` | `VERSI_KONTRAK` | salin |
| `TANGGAL_DISETUJUI` | `CATATAN_PERSETUJUAN` | **hitung** — waktu keputusan terakhir yang menyetujui |
| `PERSEN_BAGIAN_NURE` | `VERSI_KONTRAK` | salin |
| `KODE_MATA_UANG_KONTRAK` | `VERSI_KONTRAK` | salin |
| `JUMLAH_LAYER` | `LAYER` | **hitung** — cacah |
| `JUMLAH_MATA_UANG` | `MATA_UANG_KONTRAK` | **hitung** — cacah |
| `TOTAL_LIMIT_IDR` | `LAYER` | **hitung** — jumlah nilai IDR, tingkat 100% treaty |
| `TOTAL_PREMI_BRUTO_IDR` | `BAGIAN` / `DETAIL_PROPORSIONAL` | **hitung** |
| `JUMLAH_BESARAN_BERUBAH` | `NILAI_SELISIH` | **hitung** — cacah baris selisih pada versi ini ⟦**dari sisi ADJUSTMENT**⟧ |
| `TOTAL_SELISIH_IDR` | `NILAI_SELISIH` | **hitung** — jumlah selisih bersatuan `UANG`, dalam IDR ⟦**dari sisi ADJUSTMENT**⟧ |

**Dua kolom terakhir datang dari sisi Adjustment**, dan itu ditandai karena ia satu-satunya tempat
tabel datar Treaty In menyeberangi sekat. Keduanya **kosong** pada versi yang bukan penyesuaian —
bukan nol, karena nol berarti "tidak ada yang berubah" dan kosong berarti "tidak ada penyesuaian"
(ADR-0035).

### Yang sengaja TIDAK ada di tabel ini

| Tidak ada | Sebab |
|---|---|
| Kolom per layer, per mata uang, per potongan | butirannya versi; merinci akan menggandakan baris dan mengingkari §1 |
| Nilai uang tanpa mata uang | setiap nilai agregat dinyatakan **dalam IDR**, dan namanya menyebutnya |
| Nama orang pemegang kontrak | turunan yang usang (CONTEXT.md §2.9) |
| Bendera "terlambat" | keterlambatan dihitung saat dibaca, tidak disimpan |

---

## 2. `DATAR_LAYER` — butiran kedua

> **SATU BARIS = SATU LAYER PADA SATU VERSI.**

**Kunci:** `ID_LAYER`. Unik karena layer milik tepat satu versi, dan nomor layer unik di dalam
versinya (INV-05).

| Kolom | Dari entitas | Salin atau hitung |
|---|---|---|
| `ID_LAYER` · `ID_VERSI_KONTRAK` | `LAYER` | salin |
| `ID_KONTRAK` · `NOMOR_URUT_VERSI` | `KONTRAK`, `VERSI_KONTRAK` | salin — agar dapat dibaca tanpa join |
| `NOMOR_LAYER` · `BAGIAN_LAYER` · `JENIS_LAYER` | `LAYER` | salin |
| `LIMIT_IDR` · `DEDUCTIBLE_IDR` | `LAYER` | **hitung** — nilai × kurs yang tercatat pada versi |
| `KODE_MATA_UANG` | `LAYER` | salin |
| `PERSEN_MINIMUM_DEPOSIT` | `LAYER` | salin |
| `JUMLAH_DETAIL_PROPORSIONAL` | `DETAIL_PROPORSIONAL` | **hitung** — cacah; kosong pada cabang non-proporsional |
| `ADA_BAGIAN` | `BAGIAN` | **hitung**; kosong pada cabang proporsional |

---

## 3. Cara ia diisi, dan kapan ia dianggap benar

### Bentuknya: **view biasa**, bukan materialized view

| Pilihan | Putusan | Alasan |
|---|---|---|
| **View** | **dipilih** | selalu segar menurut definisi; tidak dapat basi; tidak menuntut pemantau |
| Materialized view | tidak | menuntut uji negatif dan pemantau kebasian, dan ongkos itu hanya berbayar bila pembacaannya berat |
| Proses terjadwal | tidak | ia jalur tulis, dan jalur tulis ke tabel datar dilarang §0 |

**Kesegarannya: SEGAR SAAT DIBACA.** Sebuah view tidak punya keadaan sendiri; apa yang dibaca
adalah keadaan entitas pada saat kueri berjalan.

> Tabel datar yang tidak menyatakan kesegarannya akan dibaca orang sebagai angka terkini, dan suatu
> hari ia bukan. Karena itu kesegarannya dinyatakan di sini, dan dipilih bentuk yang **tidak dapat**
> tertinggal.

### Bila kelak ia harus menjadi materialized view

Hanya bila profil beban menunjukkan pembacaannya benar-benar berat. Dan bila itu terjadi, **seluruh
aturan langkah 8 berlaku tanpa kecuali**:

| Wajib | Rujukan |
|---|---|
| **uji negatif** yang dijalankan, bukan diargumentasikan | `UJI-NEGATIF-INVARIAN.md` §2 |
| **pemantau kebasian** — `REFRESH_MODE`, `STALENESS`, terjadwal | `UJI-NEGATIF-INVARIAN.md` §3 |
| **pernyataan apa yang hilang** bila ia dimatikan | `UJI-NEGATIF-INVARIAN.md` §4 |

Perpindahan dari view ke materialized view adalah **keputusan yang ditulis**, bukan penyetelan
operasional.

---

## 4. Hak akses

| Akun | Hak atas tabel datar |
|---|---|
| pemilik objek `TREATY_MASUK` | membuat dan mengubah definisinya |
| `TREATY_MASUK_APP` | **`SELECT` saja** |
| akun Pega, selama masa berdampingan | **`SELECT` saja**, dan hanya pada bentuk baca kompatibilitas |

Tidak ada akun yang diberi `INSERT`, `UPDATE`, atau `DELETE` atas tabel datar. Larangan §0
ditegakkan **izin basis data, bukan kebijakan** (ADR-0041).

---

## 5. Apa yang tabel datar ini BUKAN

Dinyatakan supaya tidak dipakai sebagai sesuatu yang bukan dirinya:

| Bukan | Sebab |
|---|---|
| **bukan sumber migrasi** | ia turunan; sumbernya entitas di belakangnya |
| **bukan dasar rekonsiliasi terhadap sistem lama** | `TREATYINDETAIL` sistem lama tidak menerima addendum, dan `TREATYINOFFER` tidak punya penulis yang terjangkau — mencocokkan terhadap keduanya mengukur kerusakan, bukan kebenaran |
| **bukan jalur penerbitan ke luar** | jalur itu **belum diketahui** — `ERD.md` §4 |

Yang terakhir penting: tabel datar ini **tidak** menutup lubang penerbitan. Ia bentuk baca internal.
Apa yang menggantikan jalur penerbitan lama ditetapkan setelah pertanyaan bisnis di `ERD.md` §4
terjawab.
