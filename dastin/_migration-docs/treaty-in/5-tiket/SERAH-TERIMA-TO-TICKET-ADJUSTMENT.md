# Serah terima — to-ticket kemampuan Adjustment

**Tanggal:** 24 September 2026 · **Fase:** to-ticket, ronde kemampuan Adjustment · **Keadaan:** selesai
**Papan:** [`issues/`](issues/) — **papan bersama**, bukan papan Adjustment.
`GRL-01`: satu model, satu spesifikasi, **satu papan**.

> **Modul Adjustment tidak punya papan sendiri, sebagaimana ia tidak punya spesifikasi sendiri.**
> Tiket `01`…`13` lahir dari kemampuan Adjustment, tetapi tinggal di `treaty-in/5-tiket/` bersama
> tiket Treaty In. Siapa pun yang mencari papan Adjustment di folder Adjustment **tidak akan
> menemukannya, dan itu disengaja**.

---

## 1. Gerbang selesai — dijawab satu per satu

| Gerbang | Keadaan |
|---|---|
| `ASUMSI-CLEAR.md` berdiri: asumsi, tiket yang bersandar, apa yang ditinjau bila patah | **YA** — empat asumsi, masing-masing berpenanda sapuan. Kolom *"tiket yang bersandar"* **belum diisi**; lihat §5 |
| `DAFTAR-PEKERJAAN.md` diperpanjang dari `P-60`; nomor lama tidak disentuh | **YA** — `P-60`…`P-66`, tujuh kemampuan baru |
| Pengisian `NOMOR_DOKUMEN` warisan dari arsip kertas punya tiketnya sendiri | **YA** — tiket `09`, dan pelakunya **orang**, bukan sistem |
| Daftar irisan diajukan dan disetujui sebelum berkas tiket ditulis | **YA** — [`USULAN-IRISAN-ADJUSTMENT.md`](USULAN-IRISAN-ADJUSTMENT.md), disetujui dengan dua koreksi |
| Tiap tiket memakai sembilan medan, tanpa kecuali | **YA** — **13 dari 13 memuat 9/9** |
| Papan berdiri, memisahkan menahan dari menunggui | **YA** — [`issues/README.md`](issues/README.md) |
| Penahan luar papan tertulis dengan nama aslinya | **YA** — `REV-3`, `DB-20`, `DB-16a`, `DB-16b`; tidak diterjemahkan menjadi nomor tiket |

---

## 2. Yang dihasilkan

| Berkas | Isi |
|---|---|
| `issues/01-…` … `issues/13-…` | **13 berkas tiket**, satu per tiket, sembilan medan |
| `issues/README.md` | papan — pencacah, lingkup, ukuran selesai, daftar tertahan |
| `USULAN-IRISAN-ADJUSTMENT.md` | daftar irisan yang disetujui, beserta alasan tiap pemecahan |
| `ADJUDIKASI-KEMAMPUAN-ADJUSTMENT.md` | 12 calon → **2 sudah ada · 3 perlu diubah · 7 baru** |
| `ASUMSI-CLEAR.md` | empat asumsi, penanda sapuannya, dan akibat bila patah |

**Sepuluh kemampuan menjadi tiga belas tiket.** `P-65` dipecah tiga — perluas · pindahkan ·
kerutkan — karena ia **perubahan lebar**, bukan irisan tegak. Dan irisan materialitas dipecah dua
sesudah koreksi: penguncian ruas tidak memerlukan `NILAI_SELISIH`, penegakan `INV-69` memerlukannya.

### Adjudikasi yang menahan tiket kembar

Sebelum satu baris `P-60` ditulis, kedua belas calon diadili terhadap `P-01`…`P-59`. Sebabnya:
`DAFTAR-PEKERJAAN.md` induk **sudah menyebut** addendum 14×, versi 43×, selisih 4×, materialitas 2×.

Golongan **SUDAH ADA** — dua — disimpan, tidak dihapus. Ia bukti bahwa kemampuan itu **diperiksa dan
sengaja tidak ditambahkan**, bukan terlewat.

---

## 3. Yang MENAHAN

| # | Penahan | Siapa mencabutnya | Tiket yang kena |
|---|---|---|---|
| **1** | **`14` belum mendarat** — tidak satu pun tiket `01`…`13` membuat `KONTRAK` dan `VERSI_KONTRAK` | ronde Treaty In batch 1 | **`01`…`05`**, dan lewat mereka seluruh papan |
| **2** | **`D-7` diparkir** — dasarnya sudah lewat, `F-1` sudah diterapkan | **pemilik proses**, satu kalimat | kalimat pembaca arsip |
| **3** | **`DB-20`** — titik beku materialitas | bisnis | `02`, `07`, `11`, `13` |
| **4** | **`DB-16a`** — apakah revisi internal disertai dokumen | bisnis | `04`, `09` |
| **5** | **`DB-16b`** — apakah dokumen punya tanggal berlaku sendiri | bisnis | `08` |

### Paket `REV` BUKAN penahan tiket — ralat 25 September 2026

Draf pertama berkas ini mencantumkan *"tanggapan paket `REV-1`…`REV-6`"* sebagai penahan, dengan
alasan *"`REV-3` dapat mengubah kolom keadaan"*. **Itu keliru, dan ralatnya dicatat di sini alih-alih
disunting diam-diam.**

`REV-3` menyatakan di barisnya sendiri:

> **Keputusannya — *"di model baru tidak ada satu pun jalan menyetel keadaan selain melalui
> perpindahan di daftar"* — tetap, seluruhnya.**

Yang diusulkan berubah hanya **tabel §4 ADR-0055**, dan keempat butirnya seluruhnya **koreksi atas
pemerian sistem lama**: siapa yang melihat `Force Edit`, bahwa `Force Resolve` ikut menyimpan, bahwa
`SetToDirector` mati, dan bahwa pintu samping ada juga di layar addendum. **Tidak satu pun mengubah
daftar keadaan, kolom keadaan, maupun perpindahannya.**

Dan proyek ini **sudah mengadilinya lebih dulu** — `KTV-4` memeriksa keenam `REV` satu per satu:
*"nol dari enam menentukan letak kolom"*. Ralat ini hanya menyesuaikan berkas ini dengan adjudikasi
yang sudah ada.

**Akibatnya:** paket `REV` tetap harus diserahkan — ia memperbaiki **pemerian sistem lama** di ADR
induk, dan itu **utang dokumentasi**, bukan utang tiket. Tetapi ia **tidak menahan satu tiket pun**,
dan **tidak menunda batch 2**.

> **Butir 6 bertenggat.** `KTV-2` mensyaratkan kolom tanggal berlaku **dicabut sebelum data dimuat**.
> Sesudah migrasi berjalan, mencabut kolom berongkos.

### Butir 1 lahir dari ronde berikutnya, dan ia mengoreksi papan ini

Ronde Treaty In batch 1 menemukan: **tidak ada satu pun dari tiga belas tiket yang membuat
`VERSI_KONTRAK`.** Tiket `01` menambahkan kolom *"pada `VERSI_KONTRAK`"*, `02` dan `05` idem —
**ketiganya mengandaikan tabelnya sudah ada**, dan **PEMBUAT PERTAMA**-nya milik `P-01`, kemampuan
Treaty In yang ronde ini memang di luar lingkup.

| | |
|---|---|
| **Akibatnya** | `issues/README.md` sempat mencetak *"dapat dimulai: 4"*. **Keempatnya tidak dapat dimulai** |
| **Kenapa tidak terlihat** | angka itu dihitung atas **penghalang di dalam papan**, dan penghalangnya ada **di luar** |
| **Perbaikannya** | tiket `01`…`05` memperoleh `Blocked by: 14` |

**Pelajarannya berlaku untuk setiap batch berikutnya:** *"dapat dimulai"* yang dihitung hanya atas
tepi di dalam papan **akan salah** selama ada entitas yang PEMBUAT PERTAMA-nya di luar lingkup batch.

---

## 4. Yang hanya MENGUKUR KERUSAKAN — tidak menahan satu tiket pun

| Uji | Yang diukurnya |
|---|---|
| `UA-3` | berapa baris warisan **melanggar** `INV-69`/`INV-70` |
| `UA-18` | berapa daftar memuat dua baris berkunci padanan sama |
| `UA-19` | berapa nomor dokumen berulang |
| `Uji AN` | baris ganda `TREATYINDETAIL` — akibat prosedur tanpa cabang `ELSE` |
| — | ongkos pengisian ulang nomor dokumen dari arsip kertas |

Keempat kueri **belum ditulis**, dan sebelum ditulis harus **dirukunkan** dengan register uji induk
`Uji A`…`Uji AN`. Kueri kembar ke DBA adalah putaran yang terbuang, dan satu putaran diukur dalam
hari.

---

## 5. Kesalahan yang dicatat — tiga, dan dua milik penilai

| # | Kesalahan | Akibatnya |
|---|---|---|
| **1** | **Penanda `DIASUMSIKAN-CLEAR(DB-20, REV-3)`** di tiket `07` — dua kode dalam satu penanda | Sapuan `DIASUMSIKAN-CLEAR(REV-3)` **tidak menemukannya**. Seluruh alasan penanda ini ada adalah sapuan satu perintah. **Belum diperbaiki** |
| **2** | *(penilai)* **`KTV-2`/`KTV-3` dipasangkan terbalik** dengan `DB-16a`/`DB-16b` di perintah ronde sebelumnya | Tertangkap sesi to-spec dan dibetulkan sebelum mendarat. Ongkosnya nol, tetapi ia **tidak akan tertangkap** bila sesi itu menerima perintah apa adanya |
| **3** | *(penilai)* **Sisi penghalang ditinjau hanya di dalam batch** | Tepi ke entitas **di luar** batch tidak diperiksa — dan itulah yang melahirkan butir 1 §3. Ditemukan ronde berikutnya, bukan ronde ini |

### Yang belum dikerjakan dan harus dikerjakan

**Kolom *"tiket yang bersandar"* di `ASUMSI-CLEAR.md` masih berbunyi *"diisi sesudah tiket
dinomori"*.** Tiketnya sudah bernomor. Sapuannya sekarang menghasilkan:

| Asumsi | Tiket |
|---|---|
| `DB-20` | `02`, `07`, `11`, `13` |
| `DB-16a` | `04`, `09` |
| `DB-16b` | `08` |
| `REV-3` | `01`, `07` |
| ~~`F-2` / `D-1`~~ | **dicoret** — `F-2` ditutup sore 24 Sep, `D-1` mendarat |

---

## 6. Satu temuan yang naik ke modul induk

**`INV-47`, `INV-50`, dan `INV-51` ditegakkan lewat *materialized view*, dan tidak satu pun punya
pemantau kebasian.**

> MV yang gagal me-*refresh* **berhenti menegakkan tanpa satu galat pun.** Invarian yang berhenti
> menegakkan diam-diam lebih berbahaya daripada invarian yang tidak pernah dipasang — sebab yang
> pertama **tercatat sebagai terpasang**.

Tiket `11` memasang pemantau untuk `INV-69`/`INV-70` dan **menetapkan polanya**. Ketiga invarian
induk itu **belum punya tiket** — ia lubang induk, dan pemiliknya sesi to-spec Treaty In.

---

## 7. Yang ditagih sebelum implementasi

1. **Tiket `14` mendarat** — tanpa `KONTRAK` dan `VERSI_KONTRAK`, tidak satu pun tiket di papan ini
   dapat dimulai.
2. **Tiga `DB` dikirim.** `DB-16b` **bertenggat** — sebelum data dimuat. Ini yang paling mendesak;
   ketiganya menahan **lima** tiket.
3. **Paket `REV` diserahkan** — utang dokumentasi ADR induk, **bukan penahan tiket**. Lihat ralat
   di §3.
4. **Parkir `D-7` dicabut** — satu kalimat.
5. **Penanda dua-kode di tiket `07` dipecah**, dan kolom *"tiket yang bersandar"* diisi.

**Implementasi tidak dimulai.** Ia hanya dimulai atas perintah pemilik proses.
