# `F-19` · `F-20` · `F-21` — rujukan yang hilang di jalan menuju DDL

**Tanggal:** 25 September 2026 · **Ronde:** perapian artefak ERD
**Sifat:** temuan. **Dilaporkan, TIDAK ditambal** — sebabnya di §4.
**Ditemukan oleh:** `alat/cocok-relasi-erd-versus-ddl.py`, yang dapat dijalankan ulang.

> ## KENAPA KETIGANYA BARU KETAHUAN SEKARANG
>
> **Gerbang 0 ronde ini LULUS** — `KAMUS-KOLOM.md`, `ddl-usulan/`, dan `STRUKTUR-DATA.md` menyebut
> **35 entitas yang sama persis**, nol selisih. Dan batas gerbang itu tercetak di dalam keluarannya
> sendiri:
>
> > *"Ia mencocokkan **NAMA**, bukan **ISI**. Dua sumber yang menyebut nama yang sama dengan kolom
> > yang berbeda **LULUS** perkakas ini."*
>
> Ketiga temuan di bawah **seluruhnya** hidup di dalam celah itu. Tidak satu pun entitas hilang;
> yang hilang **rujukan di antara entitas yang semuanya hadir** — dan itu tidak dapat dilihat oleh
> pemeriksaan yang mencocokkan nama entitas.
>
> **Pelajaran yang dibawa, dan ia berlaku untuk setiap pemeriksaan di modul ini:** sebuah gerbang
> yang lulus **hanya menjamin apa yang diperiksanya**. Yang membuat gerbang berguna bukan angka
> "lulus", melainkan **kalimat batas yang ia cetak sendiri** — dan kalimat itulah yang menunjuk
> pemeriksaan berikutnya yang harus dibuat.

---

## `F-19` — `ID_VERSI_KONTRAK_DASAR` **tidak ada di mana pun yang mengikat**

### Apa

Kolom yang membawa **seluruh sambungan modul Adjustment** tidak ada di satu pun artefak yang
mengikat:

| Sumber | Memuatnya? |
|---|---|
| `SPEC-MODEL-DATA.md` §10.2 | **TIDAK** — nol kemunculan di seluruh berkas |
| `2-to-spec/KAMUS-KOLOM.md` | **TIDAK** — dan berkas ini **mengikat** soal nama kolom (urutan wewenang butir 5) |
| `2-to-spec/ddl-usulan/09_VERSI_KONTRAK.sql` | **TIDAK** — `VERSI_KONTRAK` hanya punya **dua** kunci asing: `ID_KONTRAK` dan `ID_DOKUMEN_ADDENDUM` |

Sementara itu ia **disebut di sebelas berkas**, dan tiga di antaranya menggantungkan seluruh isinya
padanya:

| Berkas | Bunyinya |
|---|---|
| `ERD.md` §2.2 | *"`VERSI_KONTRAK 1──o< VERSI_KONTRAK [hapus: tolak]`"* — dan judul bagiannya: **"dan ia inti sambungan Adjustment"** |
| `5-tiket/issues/01-…` | *"Artefak: **kolom `ID_VERSI_KONTRAK_DASAR`** pada `VERSI_KONTRAK`, **kunci asingnya**, dan constraint yang menjaga keterisiannya"* |
| `5-tiket/issues/40-…` | membacanya; `Blocked by: 14 · 01` |

`PETA-TELUSUR-JSON.md` memetakan `OLDID` ke sana, dan `GRL-10`/`GRL-17` memutuskannya disimpan
**tepat di satu tempat**.

### Kenapa ia berbahaya, bukan sekadar tidak rapi

**Tiket `01` adalah tiket yang paling dini dapat dimulai sesudah `14`**, dan kriteria selesainya
berbunyi *"`ID_VERSI_KONTRAK_DASAR` berdiri di `VERSI_KONTRAK` dengan kunci asingnya"*. Yang
mengerjakannya akan membuka `KAMUS-KOLOM.md` — berkas yang **mengikat** soal nama kolom — dan
**tidak menemukannya**.

Dua hal dapat terjadi, dan keduanya buruk:

| Kemungkinan | Akibatnya |
|---|---|
| ia **menamainya sendiri** | nama kolom yang tidak pernah ditinjau siapa pun masuk ke tabel yang seluruh modul Adjustment gantungkan padanya |
| ia **berhenti dan bertanya** | tiket paling dini di papan tertahan oleh sesuatu yang seharusnya sudah selesai di to-spec |

### Siapa menutup, dan apa yang dituntutnya

| | |
|---|---|
| **Siapa** | **pemilik proses**, bersama sesi yang membuka §10 |
| **Apa yang dituntut** | satu baris atribut di §10.2 — nama, asal (`OLDID` pecahan 2), tipe `R`, boleh kosong **ya** (kosong pada versi pertama), beserta catatannya. Sesudah itu kamus dan DDL **dibangkitkan ulang** dan kolomnya muncul sendiri |
| **Yang menagih** | berkas ini, `ERD.md` §2z.2, dan kriteria selesai tiket `01` yang menyebut nama kolom yang tidak ada |
| **Kenapa tidak ditambal di sini** | menambah atribut ke §10 adalah **keputusan model**, bukan perapian gambar. Dan §10 sendiri melarangnya: *"mengarang atribut membuat kamus kolom berbohong dengan cara yang paling sulit ditemukan"* |

---

## `F-20` — kunci asing diturunkan dari **pola nama**, dan 22 rujukan karena itu tidak terjaga

### Apa

Pembangkit DDL menurunkan `FOREIGN KEY` dari pola **`ID_<ENTITAS>`**. Setiap kolom rujukan yang
dinamai lain **tidak memperoleh kunci asing sama sekali** — dan ketiadaannya **tidak berbunyi di
mana pun**: tabelnya tetap berdiri, DDL-nya tetap sah, dan tidak ada pemeriksaan yang gagal.

**Sensusnya: 30 kolom rujukan tanpa kunci asing, 22 di antaranya seharusnya ada.**

| Sasaran | Berapa | Catatan |
|---|---:|---|
| **`MATA_UANG.KODE`** | **18** | `KODE_MATA_UANG` pada sepuluh tabel, `KODE_MATA_UANG_KONTRAK`, dan lima kolom mata uang `KTV-C`. **`INV-44` menuntut kode mata uang merujuk tabel acuan mata uang** |
| `KELAS_BISNIS` | 1 | `VERSI_KONTRAK.KELAS_BISNIS_KONTRAK` |
| `KONTRAK` (rujukan ke tabelnya sendiri) | 1 | `ID_KONTRAK_DISALIN_DARI` — `ERD.md` §2.1 memutuskannya `[hapus: putus]` |
| `JENIS_REASURANSI` (bersusun + rujukan sendiri) | 2 | `JENIS_REASURANSI.ID_INDUK`, `PENYEBARAN.ID_JENIS_REASURANSI_INDUK` |
| **di luar skema — benar tidak ada** | 8 | `ID_CEDANT`, `ID_ASAL_BISNIS`, `ID_REASURADUR_PEMIMPIN`, `ID_KETUA_TREATY`, `ID_DOKUMEN`, `ID_DAFTAR_RETRO`, `ID_SUSUNAN_RETRO` ×2 |

Daftar lengkapnya, per kolom, di `ERD.md` **§2z.3** dan di lembar **Daftar Relasi** pada
`ERD-SKEMA-BARU.xlsx` — keduanya **dibangkitkan**, jadi ia tidak dapat basi.

### Kenapa delapan belas yang menunjuk mata uang paling berat

`ADR-0053` menjadikan mata uang sebuah **daftar**, dan `INV-44` menjaganya. Tanpa kunci asing, sebuah
baris dapat menyimpan kode mata uang yang **tidak ada di tabel acuan** — dan yang terjadi bukan galat
melainkan **jumlah uang yang denominasinya tidak dapat ditafsirkan**, persis keadaan yang
`ADR-0039` dan tiga invarian kunci alami ada untuk mencegahnya.

**Dan ia tidak akan ketahuan dari uji negatif mana pun** yang ditulis untuk tiket `20`, sebab uji itu
memeriksa keunikan mata uang di dalam versi — bukan keberadaannya di tabel acuan.

### Siapa menutup

| | |
|---|---|
| **Siapa** | sesi yang memegang `alat/buat-kamus-dan-ddl.py` dan `alat/tulis-kamus-dan-ddl.py` |
| **Apa yang dituntut** | aturan penurunan kunci asing berhenti bersandar pada **pola nama**, dan memakai **peta rujukan yang dinyatakan** — kolom → tabel sasaran — yang dapat diperiksa orang |
| **Yang menagih** | `ERD.md` §2z.3, dan lembar Daftar Relasi yang mewarnai kolom **SEHARUSNYA ADA** |
| **Kenapa tidak ditambal di sini** | menambah 22 kunci asing **mengubah `ddl-usulan/`**, berkas yang **64 tiket** bersandar padanya dan yang kriteria selesai tiket `14` dan `15` ditulis terhadapnya. Itu perubahan artefak mengikat, bukan perapian gambar |

---

## `F-21` — perilaku hapus **diputuskan**, dan **nol** sampai ke DDL

### Apa

`ERD.md` §1 menyatakannya sebagai aturan: *"**Perilaku hapus ditetapkan sadar, tidak dibiarkan
bawaan.** Tiga nilai dipakai."* Dan §2 memenuhinya — **seluruh 36 relasi dalam-skema** membawa
`[hapus: …]`, sebagian dengan alasannya tertulis:

> *"`tolak` pada yang pertama disengaja: kontrak yang punya versi **tidak boleh** hilang, karena
> versinya memuat angka yang pernah dibukukan."*
>
> *"`CATATAN_PERSETUJUAN` dan `JEJAK_PERUBAHAN` memuat **siapa melakukan apa dan kapan**.
> Menghapusnya bersama induknya akan menghapus jejak, dan **jejak yang dapat dihapus bersama
> bendanya bukan jejak**."*

**`2-to-spec/ddl-usulan/` memuat NOL `ON DELETE`** — nol dari 38 kunci asing.

| | |
|---|---:|
| relasi dalam-skema yang `ERD.md` §2 **putuskan** perilaku hapusnya | **36 dari 36** |
| kunci asing yang **membawa** `ON DELETE` di DDL | **0 dari 38** |

### Kenapa ini bentuk kegagalan yang khas, dan layak dicatat begitu

Ia **bukan** keputusan yang belum diambil. Ia keputusan yang **sudah diambil, lengkap dengan
alasannya, lalu hilang di jalan** menuju berkas yang akan dibangun — karena pembangkit DDL tidak
pernah membaca `ERD.md`.

> **`INV-18` menuntut perilaku hapus ditetapkan sadar.** Sebuah kunci asing tanpa `ON DELETE`
> memakai bawaan Oracle — `NO ACTION` — yang **kebetulan** sama artinya dengan `tolak`. Maka untuk
> relasi ber-`[hapus: tolak]` hasilnya **kebetulan benar**, dan untuk **yang ber-`ikut hapus`**
> hasilnya **diam-diam salah**: penghapusan induk akan **ditolak** alih-alih menghapus anaknya.

Dan itu **tidak akan terlihat** sampai seseorang mencoba menghapus sebuah versi kontrak.

### Siapa menutup

| | |
|---|---|
| **Siapa** | sesi yang memegang pembangkit DDL, bersama gerbang sesi DDL |
| **Apa yang dituntut** | pembangkit membaca perilaku hapus dari sumber yang **satu** — dan sumbernya sudah ada: `ERD.md` §2. Sesudah itu `ERD.md` §2 berhenti menjadi satu-satunya tempat keputusan itu hidup |
| **Yang menagih** | `ERD.md` §2z.1, yang mencetak kolom `ON DELETE` kosong untuk **ketiga puluh delapan** barisnya |

---

## 4. Kenapa ketiganya dilaporkan dan TIDAK ditambal

Tiga sebab, dan ketiganya mengikat:

1. **`F-19` menuntut menambah atribut ke §10** — keputusan model, wewenang pemilik proses. §10
   sendiri melarang mengarang atribut.
2. **`F-20` dan `F-21` mengubah `ddl-usulan/`**, artefak yang **64 tiket** bersandar padanya. Kriteria
   selesai tiket `14` dan `15` ditulis terhadap bentuknya yang sekarang.
3. **Ronde ini bernama perapian artefak ERD**, dan perintahnya menyebut tegas apa yang boleh ditulis.
   Menambah 22 kunci asing dan 38 `ON DELETE` bukan perapian gambar.

**Yang ronde ini lakukan** adalah membuat ketiganya **berbunyi** — di `ERD.md` §2z yang
**dibangkitkan**, di lembar Daftar Relasi yang **mewarnai** kolom SEHARUSNYA ADA, dan di berkas ini.

> **Dan itu bukan hal kecil.** Sebelum ronde ini, ketiganya **tidak berbunyi di mana pun** — tidak
> ada pemeriksaan yang gagal, tidak ada angka yang menyimpang, dan Gerbang 0 melaporkan **lulus**.

---

## 5. Pemeriksaan KETIGA yang belum ada, dan ia ditulis di sini supaya tidak hilang

`alat/cocok-relasi-erd-versus-ddl.py` mencetak batasnya sendiri:

> *"Ia **TIDAK** memeriksa kardinalitas yang `ERD.md` tulis (lambang `1` / `o` / `<`) terhadap
> `UNIQUE` di DDL. Itu pemeriksaan **ketiga**, dan ia **belum ada**."*

Contoh yang membuatnya bukan kekhawatiran kosong: `ERD.md` §2.4 menyatakan
`LAYER 1──o1 BAGIAN` — **paling banyak satu** bagian per layer. DDL memang memuat
`UNIQUE (ID_LAYER)` pada `BAGIAN`, sehingga keduanya **kebetulan sepakat**. Tidak ada yang
memeriksanya, dan tidak ada yang akan tahu bila salah satunya berubah.

| | |
|---|---|
| **Siapa membuat** | sesi berikutnya yang menyentuh `ERD.md` atau `ddl-usulan/` |
| **Yang menagih** | kalimat batas di dalam keluaran perkakasnya sendiri, yang tercetak setiap kali ia dijalankan |
