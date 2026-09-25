# Prompt kerja — to-ticket (Treaty In) dan struktur data · ERD · tabel datar (Treaty In + Adjustment)

**Ditetapkan:** 24 September 2026 · **Pemilik:** pemilik proses migrasi NuRe
**Pasangannya:** `PENGETAHUAN-MIGRASI-TREATY-IN.md` berisi **APA** · `_migration-docs/METODE-GRILLING.md`
berisi **CARA grilling** · berkas ini berisi **CARA dua alur berikutnya**.

> ### Ini bukan spesifikasi, bukan tiket, dan bukan temuan.
> Ia **aturan menjalankan dua alur**. Tidak ada satu pun klaim tentang sistem lama yang lahir di
> sini. Bila berkas ini tampak menyatakan fakta tentang Pega, itu kekeliruan penulisan — yang
> mengikat adalah berkas yang dirujuknya.

---

## 0. Lingkup — ditetapkan pemilik proses, 24 September 2026

Dua alur, dan **lingkupnya berbeda**. Ini keputusan pemilik proses, bukan turunan dari berkas mana
pun, dan karena itu ia ditulis paling atas:

| Alur | Lingkup modul | Keluaran |
|---|---|---|
| **A — to-ticket** | **Treaty In SAJA** | tiket di `5-tiket/issues/` + papan `5-tiket/README.md` |
| **B — struktur data, ERD, tabel datar** | **Treaty In DAN Treaty In Adjustment** | pembaruan `4-erd-dan-tabel-datar/` |

> **Adjustment tidak memperoleh tiket.** Ia masuk alur B sepenuhnya dan tidak masuk alur A sama
> sekali. Siapa pun yang menulis tiket Adjustment dari berkas ini melanggar §0.

Dua path yang boleh dibaca:

```
D:\XML_NURE\_migration-docs\treaty-in\
D:\XML_NURE\_migration-docs\treaty-in-adjustment\
```

Alur A membaca path pertama. Alur B membaca **keduanya**.

### 0.1 Kenapa Adjustment boleh masuk alur B padahal grilling-nya belum selesai

`treaty-in-adjustment/CABANG-K-PEMETAAN-TO-SPEC.md` §4 menahan **to-spec** Adjustment sampai tujuh
butir GRILL terkunci dan paket `REV` ditanggapi. Penahan itu **tetap berlaku**, dan alur B tidak
melanggarnya, karena alur B hanya memindahkan apa yang **sudah terkunci** — butir bergolongan
KONFIRMASI dan keputusan GRL yang sudah berputus.

Pemisahnya satu kalimat, dan ia yang menentukan setiap kasus ragu di alur B:

> **Yang sudah terkunci MASUK sebagai bentuk. Yang masih GRILL MASUK sebagai LUBANG BERTANDA —
> dengan pemilik dan saat penagihannya — dan tidak pernah sebagai atribut, entitas, atau relasi.**

Empat butir yang **masih GRILL** pada tanggal berkas ini, dan karena itu **tidak boleh dimodelkan**:

| Butir | Isinya | Yang berubah bila jawabannya datang |
|---|---|---|
| **C1** | definisi materialitas addendum | materialitas menjadi atribut tersimpan, turunan jenis, atau hasil pemeriksaan |
| **C2** | himpunan nilai `JENIS_ADDENDUM` | nilai enum pada §10.2 |
| **C3** | arti bisnis `ActualValue` | **jumlah entitas** — dibawa, dibuang, atau menjadi turunan |
| **E3a** | pro rata (`EDMEffective`) | **"berlaku sejak" menjadi atribut versi** — perubahan bentuk |

**C3 dan E3a mengubah bentuk.** Keduanya wajib muncul di keluaran alur B sebagai lubang bernomor,
bukan sebagai kotak yang digambar dengan garis putus lalu dilupakan.

---

## 1. Urutan wewenang bila dua masukan bertentangan

Diambil dari `5-tiket/KEPUTUSAN-PEMBAGIAN-TIKET.md` §0.2, dengan keadaannya pada hari ini:

| # | Sumber | Keadaan |
|---|---|---|
| 1 | `CONTEXT.md` — aturan kerja §2.0…§2.9b, penamaan §16 | berlaku |
| 2 | ADR di `docs/adr/` — 0034…0055 | berlaku |
| ~~3~~ | ~~folder `ddl/`~~ | **DICORET** — foldernya belum ada (L-1). Anak tangga yang tidak ada akan dipakai orang seolah ia ada |
| 4 | `SPEC-MODEL-DATA.md` §10 | berlaku — **28 entitas**, §10.23c |
| 5 | `SPEC-INVARIAN.md` | berlaku — 63 + 2 |
| 6 | `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` | **daftar entitas yang mengikat** — bukan §2.3 (L-6) |
| 7 | `4-erd-dan-tabel-datar/ERD.md` **teks §2** | mengikat untuk relasi; §6 gambar tidak |

Dan satu yang mendahului seluruhnya karena ia soal bukti, bukan soal isi:

> **Masa lalu tidak diputuskan, ia dibuktikan.** Setiap klaim tentang apa yang sistem lama lakukan
> menuntut `EVIDENCED`. `DECIDED` sah untuk putusan bisnis dan rancangan kita, **tidak pernah**
> untuk pernyataan tentang masa lalu.

---

## 2. Masukan wajib — diperiksa ada sebelum langkah pertama

Bila salah satu tidak ada, **itu bukan alasan menunda**: catat ketiadaannya sebagai lubang bernomor
dan kerjakan sisanya (§6, aturan berhenti tidak memuat "masukan kurang").

### Alur A

| Berkas | Perannya |
|---|---|
| `5-tiket/DAFTAR-PEKERJAAN.md` | **sumber tiket** — kemampuan `P-01`…`P-59` |
| `5-tiket/BENTUK-TIKET.md` §6 | **tiga belas medan**, dan §7 penjaganya |
| `5-tiket/KEPUTUSAN-PEMBAGIAN-TIKET.md` | gerbang G1 (besar tiket), G2 (isi penyerahan pertama), G3 (tiket tertahan) |
| `5-tiket/LUBANG-SPESIFIKASI.md` | L-1…L-10 beserta pemilik dan saat penagihannya |
| `SPEC-MODEL-DATA.md` §10 | nama atribut, tipe, keterisian — **yang membuat kriteria selesai dapat gagal** |
| `SPEC-INVARIAN.md` | `INV-01`…`INV-63` + 2 tidak-dapat-dilanggar |
| `docs/adr/0034…0055` | dasar `DECIDED` |
| `PENGETAHUAN.md`, `TETAPAN-DI-KODE.md`, `4-erd-dan-tabel-datar/CABANG-MATI-DI-JALUR-PERSETUJUAN.md` | dasar `EVIDENCED` |
| `komite-claim-non-prop/3-to-tickets/issues/` | **contoh BENTUK**, bukan sumber isi |

### Alur B

Seluruh masukan alur A, **ditambah** dari sisi Adjustment:

| Berkas | Perannya |
|---|---|
| `treaty-in-adjustment/CABANG-K-PEMETAAN-TO-SPEC.md` | ke mana tiap keputusan Adjustment mendarat; §4 prasyarat |
| `treaty-in-adjustment/KEPUTUSAN-GRILLING-ADJUSTMENT.md` | putusan `GRL-01`…`GRL-12` |
| `treaty-in-adjustment/PEMILAHAN-SISA-GRILLING.md` | golongan tiap butir — **yang membedakan terkunci dari GRILL** |
| `treaty-in-adjustment/GRILL-A/06-PUTUSAN.md`, `GRILL-B/06-PUTUSAN.md` | bahan to-spec `B-1`…`B-4`, `E-1`, `R-F1` |
| `treaty-in-adjustment/USULAN-REVISI-ADR.md` | `REV-1`…`REV-3`, `IND-1`…`IND-3` — **belum ditanggapi pemilik ADR** |
| `treaty-in-adjustment/PETA-SUMBER-INDUK.md` | berkas induk mana yang disentuh keputusan Adjustment |
| `treaty-in/SEAM-ADJUSTMENT.md` | `KUNCI_PADANAN`, Bentuk A dan Bentuk B |
| `4-erd-dan-tabel-datar/TEMUAN-ADJUSTMENT-DITUNDA.md` | delapan temuan yang **tidak diadili** |

---

## 3. Alur A — urutan langkah, dan keluaran ditulis saat LANGKAHNYA selesai

> **Keluaran ditulis saat langkahnya selesai, bukan saat sesinya selesai.** Sesi dapat berakhir
> kapan saja; apa yang hanya ada di konteks kerja **tidak ada**.

| # | Langkah | Keluaran, ditulis begitu langkah selesai |
|---|---|---|
| **A1** | Cocokkan `DAFTAR-PEKERJAAN.md` dengan `SPEC-MODEL-DATA.md` §10 **dua arah**: kemampuan yang menyentuh entitas tanpa §10, dan entitas ber-§10 yang tidak disentuh kemampuan mana pun | `5-tiket/PRA-TIKET-PENCOCOKAN.md` — dua daftar, dan **hasil bersih ikut dicatat** |
| **A2** | Tetapkan pemecahan G3: kemampuan mana menjadi **lebih dari satu** tiket, dan garis pemotongnya | bagian di berkas yang sama, §2 |
| **A3** | Tulis tiket, **satu berkas per tiket**, nomor berurut mulai `01`, nama berkas `NN-judul-dipotong-sekitar-60-karakter.md` | `5-tiket/issues/NN-….md` — **ditulis satu per satu, bukan dikumpulkan** |
| **A4** | Bangkitkan papan dari berkas tiket — dibaca dari medan `status:`, **tidak pernah menghapus tiket** | `5-tiket/README.md` |
| **A5** | Peta liputan: tiap `P-xx` menunjuk tiket mana, dan tiap tiket menunjuk `P-xx` mana — **dua arah** | `5-tiket/PETA-LIPUTAN.md` |
| **A6** | Daftar tiket yang `DASAR:`-nya seluruhnya `DECIDED` padahal `golongan: pelestarian` | bagian di `PETA-LIPUTAN.md` §3 — lihat aturan putusan P-4 |

**A1 mendahului A3 dan tidak boleh dilompati.** Sebabnya sudah terbukti sekali di proyek ini: lima
entitas berkunci alami tanpa invarian bernomor lolos **tujuh sesi**, bukan karena ada yang lalai
melainkan karena tidak ada satu pun langkah yang mencocokkan dua berkas **ke dua arah**.

### 3.1 Bentuk tiket — tiga belas medan, tidak ditawar

Bentuknya ada di `BENTUK-TIKET.md` §6 dan **tidak diulang di sini** — satu fakta, satu tempat.
Yang diulang hanya tiga penjaganya, karena ketiganya yang paling mudah luntur saat menulis tiket
kelima puluh:

1. **`KENAPA BEGINI` harus MUSTAHIL ditulis dari medan lain.** Bila isinya dapat disusun dari
   `SATU KALIMAT` + `INVARIAN` + `ASALNYA DARI MANA`, ia **kosong**. Yang sah: menyebut apa yang
   terjadi di sistem lama, atau akibat yang akan **mengejutkan** pembaca yang tidak ikut sesi mana
   pun.
2. **`MENGGANTIKAN:` punya kewajiban keterisian menurut golongan** — `perubahan` dan `pelestarian`
   wajib menyebut aturan lama beserta kode cacatnya; `baru` wajib berbunyi **"tidak ada"**.
   Kekosongan yang **dinyatakan**, bukan kekosongan.
3. **`EVIDENCED(...)` menyebut EKSPOR MANA**, bukan hanya berkasnya — `EVIDENCED(SetSpreadName@ekspor-2026-09)`.
   Sebabnya L-9: asal lingkungan ekspor belum diketahui, dan bila jawabannya kelak "QA", satu sapuan
   atas nama ekspor mengeluarkan daftar tiket yang terdampak.

Dan kriteria selesai ditulis sebagai **pemeriksaan yang dapat gagal**, bukan susunan layar — L-4
menyatakan tidak ada spesifikasi layar di mana pun.

### 3.2 Yang TIDAK ikut ditiru dari folder contoh

Yang ditiru **bentuknya**; yang menyebut **dunia nyata** adalah isi, dan isi tidak ikut. Ujinya saat
ragu: **apakah ia menyebut sesuatu di luar berkas itu?** Nama medan, susunan bagian, konvensi nama
berkas — bentuk. Kode status HTTP, nama tabel, nama peran, angka ambang — dunia nyata, dan dunia
nyata modul contoh bukan dunia nyata modul ini.

---

## 4. Alur B — urutan langkah

| # | Langkah | Keluaran |
|---|---|---|
| **B1** | **Inventarisasi apa yang basi.** Bandingkan tanggal dan isi tiap berkas `4-erd-dan-tabel-datar/` terhadap keputusan 24 September — §10 28 entitas, §10.23c, `SUMBU_REKONSILIASI`, koreksi `RINCIAN_PENYEBARAN`, `KUNCI_PADANAN` `POTONGAN` | `4-erd-dan-tabel-datar/AUDIT-KEBASIAN.md` — satu baris per berkas, **termasuk yang ternyata segar** |
| **B2** | Perbarui `STRUKTUR-DATA.md`: entitas kedua modul, dengan **tiga hal per entitas** — kunci alami + lingkupnya, `SUMBU_REKONSILIASI`, dan `KUNCI_PADANAN` | `STRUKTUR-DATA.md` |
| **B3** | Perbarui `ERD.md` **teks §2** — kardinalitas, keterisian, perilaku hapus, penanda ⟦SEKAT⟧ — lalu §6 gambarnya | `ERD.md` |
| **B4** | Perbarui `TABEL-DATAR.md` — butiran tiap tabel datar, kunci, asal tiap kolom, salin-atau-hitung, kesegaran | `TABEL-DATAR.md` |
| **B5** | Jalankan ulang perkakas turunan; **suntingan dilakukan di dalam skripnya, bukan di hasilnya** | `Diagram-Skema-Tabel-TreatyMasuk.xlsx`, `ERD-TREATY-MASUK.html`, `PETA-NAMA-TABEL-TREATYIN.md`, `peta-nama-tabel-treatyin.tsv`, `ERD-STRUKTUR-TREATYIN.html` |
| **B6** | Perbarui `ISI-FOLDER.md` — cacah berkas, mana turunan mana tulisan tangan, dan **daftar penyimpangan dari folder contoh ke dua arah** | `ISI-FOLDER.md` |
| **B7** | Kumpulkan lubang Adjustment yang **tidak dimodelkan** — C1, C2, C3, E3a dan sisanya — masing-masing dengan pemilik dan saat penagihan | `4-erd-dan-tabel-datar/LUBANG-ADJUSTMENT-DI-STRUKTUR.md` |

### 4.1 Tiga kewajiban per entitas di B2 — dan pembacanya berbeda-beda

Ini yang membuat B2 bukan penyalinan:

| Kewajiban | Bentuknya | **Siapa yang membacanya** |
|---|---|---|
| **kunci alami + lingkup** | Bentuk A (unik di dalam versi) atau Bentuk B (unik di dalam induk langsung) | yang membangun **pemadanan antar versi** — seam Adjustment |
| **`SUMBU_REKONSILIASI`** | **nama kolom** pengelompokan saat baris turunannya dijumlahkan | yang menulis **`GROUP BY`** di sesi DDL |
| **`KUNCI_PADANAN`** | identitas bisnis yang stabil lintas versi; entitas berinduk dua menyebut **induk mana** | yang menghitung **nilai selisih per baris** |

> **Entitas yang TIDAK punya salah satunya menyatakan itu juga, beserta alasannya.** Entitas tanpa
> kunci alami yang tidak menyatakan dirinya begitu **terlihat persis sama** dengan entitas yang lupa
> diberi kunci — dan bedanya baru ketahuan jauh di hilir, saat pemadanannya tidak bekerja.

Dan lingkupnya **dibaca dari penulisnya, bukan dari bentuk yang terlihat**: sebelum menulis lingkup
sebuah keunikan, buka aktivitas yang **menambah baris** ke daftar itu dan lihat apa yang dilewati
gelungnya. Tiga invarian sudah pernah salah lingkup dalam satu hari karena aturan ini dilewati —
sebuah daftar bernama "termin" ternyata baris **per mata uang**.

### 4.2 Tabel datar — dua lapisan nama, dan keduanya hidup berdampingan

| Lapisan | Contoh | Aturan yang berlaku |
|---|---|---|
| **`NAMA_T`** — tabel datar, pemetaan pohon clipboard Pega | `T_TREATY_IN_LIMITS` | awalan `T_` **dipakai**; ia yang menandai bahwa ia bukan entitas rancangan |
| **`PADANAN_DDL`** — objek Oracle | `LAYER` | §16 penuh: Indonesia, `UPPER_SNAKE_CASE`, kata utuh, **≤ 30 bita dihitung** |

Menghapus `T_` tidak menegakkan §16; ia menghapus satu lapisan pemetaan dan membuat tabel datar tak
lagi dapat dibedakan dari entitas rancangan. `ERD.md`, `STRUKTUR-DATA.md`, dan kedua berkas gambar
rancangan tetap memakai nama §16.

Dan satu aturan yang tidak boleh luntur: **tabel datar tidak pernah ditulis aplikasi, ia diturunkan.**
Akun `TREATY_MASUK_APP` diberi `SELECT` saja. Sistem lama melanggar aturan ini dan akibatnya sudah
terlihat — `TREATYINDETAIL` diisi dari dalam jalur penyimpanan, sehingga ketika satu langkah
dimatikan di jalur addendum, tabelnya diam-diam berhenti menerima addendum **tanpa satu pun galat**.

---

## 5. Aturan putusan — ditetapkan DI MUKA, satu baris per cabang hasil

Supaya pekerjaan tidak berhenti untuk bertanya pada cabang yang sudah terjawab, dan berhenti tepat
pada cabang yang memang menuntut orang.

| # | Bila ditemukan… | Putusan |
|---|---|---|
| **P-1** | kemampuan menyentuh entitas yang **tingkat pencatatannya** `BELUM DITENTUKAN` (§10.23a) | tiket **tetap ditulis**; yang belum diketahui menjadi medan `PENGHALANG` dengan pemilik dan saat penagihan. **Bukan penghenti** |
| **P-2** | satu kemampuan jelas menjadi dua tiket atau lebih menurut G3 | **pecah**, nomor berlanjut, dan `PETA-LIPUTAN.md` mencatat `P-xx` yang sama menunjuk beberapa tiket |
| **P-3** | dua kemampuan ternyata menggantikan **aturan lama yang sama** | **satu tiket**, dan keduanya disebut di `ASALNYA DARI MANA`. Dua tiket yang diam-diam menggantikan aturan yang sama adalah cacat yang medan `MENGGANTIKAN:` memang dipasang untuk menangkapnya |
| **P-4** | tiket `golongan: pelestarian` yang seluruh `DASAR:`-nya `DECIDED` | **baca ekspornya dulu.** Bila buktinya ada, tulis `EVIDENCED`. Bila **tidak ada**, golongannya berubah menjadi calon `perubahan` atau `baru` dan itu **temuan bernomor**, bukan koreksi diam-diam |
| **P-5** | angka di sebuah berkas berbeda dari angka di berkas lain | **keduanya dikutip**, yang lebih baru menang, dan selisihnya masuk `AUDIT-KEBASIAN.md`. Tidak ada angka yang diperbaiki tanpa barisnya sendiri |
| **P-6** | butir Adjustment yang masih **GRILL** menyentuh bentuk | **jangan dimodelkan.** Ia masuk `LUBANG-ADJUSTMENT-DI-STRUKTUR.md` dengan pemilik dan saat penagihan |
| **P-7** | butir Adjustment bergolongan **KONFIRMASI** atau **putusan GRL** | **dimodelkan**, dan asal putusannya disebut |
| **P-8** | entitas yang atributnya **jauh lebih sedikit daripada perannya menuntut** | itu tanda **batas metode**, bukan entitas sederhana. Tambal dengan **sumber yang mandiri**, bukan penyisiran yang sama dijalankan lebih teliti |
| **P-9** | sebuah daftar di sistem lama tampak layak menjadi entitas baru | **daftar bukan entitas sampai terbukti punya identitas sendiri.** Tiga usul semacam itu sudah gugur di §10.23c — dua karena daftarnya **perancah clipboard**, satu karena buktinya **tidak memisahkan apa pun** |
| **P-10** | sebuah invarian dirumuskan tanpa menghasilkan **nama kolom** | tolak rumusannya. Kriteria yang tidak bisa gagal bukan kriteria, dan rumusan yang tidak dapat dikompilasi memindahkan keputusannya ke orang yang tidak tahu bahwa ia sedang memutuskan sesuatu |
| **P-11** | satu properti lama ternyata punya **lebih dari dua** arti | **BERHENTI dan laporkan** (§6 butir 1) |
| **P-12** | penahan menunggu pihak luar — instans Oracle, izin kueri, ekspor kedua, dokumen kebijakan | langkahnya **tetap dikerjakan**; bagian yang menunggu **dipecah keluar** menjadi butir tertahan bernomor, **tanpa taksiran**, muncul di papan |
| **P-13** | satu nilai lama dapat dibaca sebagai **satu benda** atau **dua benda**, dan ekspor tidak memisahkannya | **ambil sisi yang ongkos salahnya murah** — lihat §5.1 |

### 5.1 Satu benda atau dua — aturannya ongkos, bukan kemungkinan

Bentuk ini sudah muncul dan akan muncul lagi: sebuah nilai warisan yang mungkin sama dengan atribut
yang sudah dimodelkan, mungkin pula benda lain. Jangan mencari yang paling mungkin benar; bandingkan
**ongkos salahnya ke dua arah**.

| Pilihan | Bila ternyata salah | Dapat dibatalkan? |
|---|---|---|
| **dua atribut** | satu kolom yang ternyata selalu kosong atau menggandakan tetangganya | **ya** — kolomnya dicabut |
| **satu atribut** | dua pengenal warisan yang berbeda **larut menjadi satu**, dan yang kalah **musnah** | **tidak** — tidak ada tempat lain yang memuatnya |

Karena ongkosnya tidak setangkup, **ambil dua atribut**, catat dasarnya `DECIDED`, dan tulis syarat
pembalikannya. Dan pada kasus seperti ini wajib ditulis pula **aturan pembedanya per baris** — apa
yang dibaca migrasi untuk memutuskan baris ini termasuk yang mana — karena tanpa itu, dua kolom
hanya memindahkan kebingungan satu tingkat ke hilir.

---

## 6. Aturan berhenti — DAFTAR TERTUTUP

Tiga keadaan, dan hanya tiga:

1. cabang hasil yang **ditandai "berhenti"** oleh aturan putusan §5 — hari ini hanya **P-11**;
2. sesuatu yang ditemukan **MEMBATALKAN** sebuah keputusan yang sudah ditetapkan — bukan menambah,
   **membatalkan**;
3. sebuah blok tidak dapat dikerjakan tanpa memutuskan hal yang **wewenangnya ada pada pemilik
   proses**, **dan** pemecahan pekerjaan tidak menolong karena yang belum diputuskan duduk **tepat
   di tengahnya** — bentuknya seperti `P-41`, yang sudah diuji dan memang tidak dapat dipecah.

> **Selain ketiganya: kerjakan, catat, lanjut.**

Tanpa kalimat penutup itu, setiap keraguan berubah menjadi alasan berhenti yang tampak sah.

**Uji butir 3 sebelum memakainya**, karena ia yang paling sering dipakai keliru: bila sisa blok masih
dapat dikerjakan tanpa jawaban itu, maka yang belum diputuskan **tidak** duduk di tengah — pertanyaannya
diajukan, pekerjaan yang tidak bergantung padanya **tetap berjalan**, dan yang bergantung dipecah
keluar sebagai butir tertahan. Mengajukan pertanyaan benar; menghentikan seluruh blok karenanya
tidak.

**Laporan di setiap batas blok: ringkas, bukan lengkap** — lalu lanjut tanpa menunggu jawaban.
Laporan yang lengkap di tengah pekerjaan adalah bentuk lain dari berhenti.

---

## 7. Penjaga per keluaran

### 7.1 Setiap perkakas melaporkan apa yang DITOLAKNYA dan berapa banyak

> Penolakan yang diam adalah bentuk perkakas dari *"diam bukan bukti"*: yang dibuang tanpa hitungan
> membuat yang tersisa **terbaca sebagai keseluruhan**.

Ini sudah menelan korban di proyek ini: satu penjaga yang **benar** membuang **414 properti** tanpa
jejak, dan lubangnya lolos tujuh sesi. Satu baris `"414 properti ditolak"` akan menemukannya di hari
pertama.

Maka tiap penyisir, pengklasifikasi, atau pembangkit mencetak — di samping hasilnya — berapa banyak
yang **tidak masuk hasil dan atas dasar apa**, satu baris per sebab. Dan dua hal dipisah saat
melaporkannya: **tidak ada sama sekali** versus **ada, tetapi tidak di tempat yang diperiksa**.

### 7.2 Setiap perkakas menyatakan BATAS KEMAMPUANNYA di dalam dirinya sendiri

> Perkakas ini **tidak dapat menyatakan sebuah nama benar**. Ia hanya dapat menyatakan sebuah nama
> **SALAH**, atau **BELUM TERPERIKSA**.

Tanpa kalimat itu, "nol temuan" akan dibaca sebagai "semuanya benar" — kesimpulan yang perkakasnya
tidak pernah sanggup berikan.

### 7.3 Setiap sapuan menyebut NAMA TAG yang dipakainya

Diperiksa lebih dulu terhadap `alat/datar-nama-elemen.csv`. Sebabnya sudah terbukti berulang: satu
sapuan mencari `<PropertiesName>` saja dan **73 berkas** memakai `<pyPropertiesName>`; satu alat
mengenal `pyStepsBlockName` saja dan Section/Harness/Data Transform memakai `pyDisabled`; dan bentuk
penulis kelima — `<pyStepsCallParams><CopyInto>` — hanya ketahuan lewat **nol yang mustahil**.

> **Daftar nama elemen menjawab "nama mana lagi untuk gagasan ini"; ia tidak menjawab "gagasan ini
> berbentuk apa lagi".** Yang menjawab yang kedua adalah **properti yang punya layarnya sendiri
> tetapi nol penulis**.

### 7.4 Penanda atas sesuatu yang sengaja tidak dikerjakan ditulis sebagai KEPUTUSAN

Bukan sebagai pekerjaan tertunda. Ia menyebut: apa yang sengaja tidak dilakukan, **kenapa**, apa
akibatnya pada angka yang dibaca, di mana selisihnya dapat dilihat, dan **kapan ia ditagih**.

Penanda berbunyi *"BELUM TERPASANG"* akan dipasang orang berikutnya — **tanpa tahu kenapa ia tidak
dipasang, dan tanpa dapat mengujinya juga.**

### 7.5 Papan dibangkitkan dari berkas, dan tidak pernah menghapus tiket

Status adalah **field di dalam berkas**, bukan lokasi foldernya. Penahan dari luar papan memakai
**nama aslinya** — `L-3`, `CO-3` — tidak diterjemahkan jadi nomor tiket. Penahan yang selesai
**dicoret, tidak dihapus**. **Nomor yang lompat bukan kekeliruan**; menyusun ulang penomoran memutus
setiap rujukan yang sudah ada. Dan **"menunggui, tidak menahan"** dipisahkan dari penahan sungguhan,
supaya papan tidak berteriak serigala.

---

## 8. Larangan

| | Dilarang | Sebab |
|---|---|---|
| 1 | **tiket untuk Adjustment** | §0 — keputusan pemilik proses |
| 2 | **DDL, struct Golang, komponen React, rancangan endpoint** | sesi DDL belum berjalan (L-1); urutan wewenang butir 3 dicoret |
| 3 | **spesifikasi layar sebagai kriteria selesai** | L-4 — tidak ada spesifikasi layar di mana pun |
| 4 | **memodelkan butir Adjustment yang masih GRILL** | §0.1 — C1, C2, C3, E3a |
| 5 | **memperbaiki cacat diam-diam** | cacat ditemukan, dicatat, lalu **diputuskan** |
| 6 | **menambal lubang dari ekspor** | menambalnya menjadi keputusan tak bertanda yang tidak pernah ditinjau siapa pun |
| 7 | **menyunting berkas turunan** | sunting **di dalam skripnya**. Bila hasil berbeda dari sumbernya, **alatnya yang salah** |
| 8 | **menyalin klaim domain dari folder contoh** | yang ditiru bentuknya; isi tidak ikut |
| 9 | **menyebut nama orang sebagai peran** | sistem lama sudah melakukannya empat kali dan akibatnya ada di §6.1 butir 9 `PENGETAHUAN-MIGRASI-TREATY-IN.md` |
| 10 | **membersihkan `ActualValue` bersarang sebelum dibaca sekali** | tiap lapis sarang adalah potret kontrak pada satu addendum sebelumnya — satu-satunya sisa dari nilai disetujui yang tidak dapat direproduksi |
| 11 | **membawa satu nama untuk dua arti** | menanam cacat sistem lama dengan tangan sendiri. Tiap arti diberi nama yang **menyebut dirinya** |

---

## 9. Ketidakcocokan yang SUDAH terbaca — dicocokkan, bukan ditambal

Ditemukan saat prompt ini disusun, 24 September 2026, dengan membandingkan berkas terhadap berkas.
Seluruhnya **dikutip, bukan disimpulkan**. Masing-masing wajib muncul kembali di `AUDIT-KEBASIAN.md`
(B1) atau `PRA-TIKET-PENCOCOKAN.md` (A1) dengan putusannya, **termasuk bila putusannya "tidak ada
yang berubah"**.

| # | Ketidakcocokan | Di mana |
|---|---|---|
| **S-1** | **Cacah kemampuan berbeda di dalam satu berkas.** §2 berjudul *"60 didaftar, 57 aktif"*; §4 mencatat *"didaftar 58 · aktif 55"*. Selisihnya **tepat dua** — `P-58` dan `P-59`, yang ditambahkan 24 September | `DAFTAR-PEKERJAAN.md` §2 vs §4 |
| **S-2** | **Jumlah entitas.** `DAFTAR-PEKERJAAN.md` §5.3 dan `PENGETAHUAN-MIGRASI-TREATY-IN.md` §7.1 menyebut **27**; `SPEC-MODEL-DATA.md` §10.23c berputus **28** — `PEMULIHAN_LIMIT` diterima. §7.1 **sudah memuat** `PEMULIHAN_LIMIT` di daftar anaknya tetapi angkanya belum ikut | §10.23c: *"27 menjadi 28, bukan 30"* |
| **S-3** | **`ERD.md` dan `TABEL-DATAR.md` bertanggal 23 September**, sedangkan `SPEC-MODEL-DATA.md` §11.1 menyatakan sendiri ERD *"dijalankan ulang sesudah keputusan §10.23c"* — dan §10.23c berputus **24** September. Keduanya basi **menurut pernyataan berkasnya sendiri**, bukan menurut dugaan | §11.1 baris langkah 10 |
| **S-4** | **Judul tabel berselisih dengan isinya.** §10.23b berjudul *"hanya empat entitas punya baris turunan"*; tabelnya memuat **enam** | `SPEC-MODEL-DATA.md` §10.23b |
| **S-5** | **`ISI-FOLDER.md` menyebut "16 berkas"**; foldernya kini memuat lebih dari dua puluh, sebagian ditambahkan 24 September | `ISI-FOLDER.md` baris kedua |
| **S-6** | **Tanggal berkas sisi Adjustment berada di depan hari ini** — beberapa bertanggal **26 dan 27 September 2026**, sedangkan hari ini **24 September 2026**. Akibat praktisnya satu: *"mana yang lebih baru"* **tidak dapat dibaca dari tanggal di dalam berkas** untuk folder itu | `CABANG-K-PEMETAAN-TO-SPEC.md`, `PEMILAHAN-SISA-GRILLING.md` |

**S-1 sampai S-5 berbentuk sama**, dan bentuk itu layak disebut karena ia akan terulang: sebuah
angka diperbarui di satu tempat dan tidak di tempat lain. Itu persis kegagalan yang aturan **satu
fakta, satu tempat** dipasang untuk mencegah — dan pada S-1 dan S-4 ia terjadi **di dalam satu
berkas**.

> Yang dikerjakan bukan menyeragamkan angkanya. Yang dikerjakan adalah **menetapkan mana yang
> mengikat, menulis sebabnya, dan mencabut yang tidak** — supaya pembaca berikutnya tidak menemukan
> tiga angka dan memilih salah satu.

---

## 10. Pelaporan

Di setiap batas blok — sesudah A1, A3 (tiap sepuluh tiket), A5, B1, B3, B5, B7 — satu laporan
**ringkas**: apa yang selesai, berapa keluaran yang ditulis, apa yang ditemukan, dan apa yang
**ditolak beserta jumlahnya**. Lalu lanjut.

Di akhir: `5-tiket/RINGKASAN-TO-TICKET.md` dan bagian penutup `AUDIT-KEBASIAN.md`, keduanya memuat
hal yang sama bentuknya —

| Yang dilaporkan | Bukan |
|---|---|
| **daftar**, dengan alasan per baris | persentase |
| **hitungan yang diverifikasi ulang** dari sumbernya | hitungan yang diingat |
| **apa yang tidak dikerjakan dan kenapa** | diam |

> Daftar dengan alasan per baris dapat diperiksa orang lain. Persentase tidak.
>
> Hitungan yang tidak pernah dituangkan sebagai daftar **belum ditinjau siapa pun**.

---

## 11. Kalimat yang dipakai saat ragu

> Ini kebutuhan bisnis, atau akibat dari cara sistem lama dibangun?

> Struktur yang terlihat bukan struktur yang berlaku.

> Rujukan bukan panggilan. Panggilan bukan langkah hidup.

> Tidak ada di ekspor bukan tidak ada di sistem.

> Masa lalu tidak diputuskan, ia dibuktikan.

> Uji yang dijawab sama oleh kedua kemungkinan bukan uji.

> Kriteria yang tidak bisa gagal bukan kriteria.

> Lingkup ditulis dari PENULISNYA, bukan dari bentuk yang terlihat.

> Daftar bukan entitas sampai terbukti punya identitas sendiri.

> Nama yang berarti dua hal di sistem lama tidak dibawa.

> Keluaran ditulis saat langkahnya selesai, bukan saat sesinya selesai.
