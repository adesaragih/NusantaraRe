> Modul  : Komite Claim Non Prop · Ronde 05 · 2026-09-20
> Peran  : interogator
> Masukan: `01-TEMUAN.md` · `02-SIDANG.md` · `INVENTARIS-BUKTI.md` (2026-09-20) · `REGISTER-PAGAR.md`
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

# 03 · LUBANG

Dua daftar, dipisah tegas: **pertanyaan** yang menunggu keputusan orang, dan **pengambilan
bukti** yang menunggu berkas atau baris data. Sebuah perkara tidak boleh berada di keduanya.

---

## 1. Yang sudah ditutup ronde ini dengan membaca — bukan dengan bertanya

Aturan biaya langkah 1 berbunyi: berkas yang sudah dipegang **dibaca dan dilaporkan sebagai
fakta**. Empat perkara yang menginap sejak `PENGETAHUAN.md` §14 ditutup di sini tanpa satu
pun pertanyaan.

| Perkara | Ditutup oleh | Hasil |
|---|---|---|
| Q-6 — pintu masuk lain `KomitePostAdjustmentCWP` | pembacaan 338 berkas | N-13 · nol pemanggil lain di dalam ekspor |
| Q-7 — pemicu `komiteAccept_ticket` | `KomiteTreaty_Flow.xml` | N-14 · ticket menempel pada gerbang `KomiteLoop`; penaiknya di luar ekspor |
| Q-3 — arti empat nilai `FlagProrate` | `ReinstatementPremiumDetails.xml` | N-15 · empat nilai terpakai; artinya menunggu Field Value |
| Q-2 — nilai `TransferType` dan apakah masih dipakai | `CreateChildKomiteCNP_Act.xml`, `CreateChildKomiteCloseNP_Act.xml`, `KomiteRouter.xml` | N-04 · masih dipakai, satu penulis, satu jalur tanpa penulis; `'2'` menghidupkan jalur banyak-jenjang |

Ketiganya semula terdaftar sebagai "perlu kueri" atau "perlu ekspor". Tidak satu pun
memerlukannya.

---

## 2. Pertanyaan `Q5-x`

**Satu** pertanyaan, dan ia membawa rekomendasi. Sebuah pertanyaan kedua sempat disiapkan —
urutan giliran jenjang — lalu dibatalkan karena jawabannya ternyata terbaca di
`KomiteRouter.xml` langkah 6.1: giliran jatuh pada baris **pertama** yang belum memutuskan,
urut `.DEGREE ASC`. Pertanyaan yang jawabannya ada di berkas yang dipegang bukan pertanyaan.

```
Q5-1 · Klaim bersyarat dan ambang
  Perkara   : N-10
  Pertanyaan: Klaim bersyarat selalu memakai roster terluas (LIMIT_BOTTOM = 0), menimpa
              kelas 25 juta. Apakah itu kebijakan yang dibawa, atau akibat urutan langkah
              yang tidak disengaja?
  Pilihan   : (a) Kebijakan — "bersyarat" menjadi dimensi kedua pada tabel seleksi H-1.
              (b) Bukan kebijakan — bersyarat mengikuti kelas nilainya seperti klaim lain.
  Rekomendasi: (a). Menimpa itu berdiri di sub-langkah tersendiri yang berpra-syarat
              eksplisit `Local.Subjectivity==true`, bukan di sela-sela langkah lain —
              bentuknya bentuk aturan, bukan bentuk kecelakaan.
  Pemegang  : pemilik proses
  Memblokir : tidak memblokir — H-1 sudah menetapkan seleksi sebagai data; ini menambah
              satu dimensi pada tabel yang sama
```

Pertanyaan ini digabungkan dengan perkara pemilik proses yang sudah terbuka di
`07-AUDIT.md` bagian 4, agar pemilik proses menerima satu daftar, bukan tiga.

---

## 3. Pengambilan bukti

Urutan biaya dipatuhi: tidak ada satu pun usulan Tracer di daftar ini. Dua usulan yang
sempat berbentuk Tracer diganti — penggantinya disebut pada kolom terakhir.

| # | Bukti | Tingkat biaya | Pemegang | Menutup | Pengganti Tracer |
|---|---|---|---|---|---|
| B5-1 | Baris `INVENTARIS-BUKTI.md` §2.5 untuk **data produksi yang belum diambil** | 1 — tulisan sendiri | kamu sendiri | mengesahkan seluruh pagar A-4 dan A-5 | — |
| B5-2 | Ekspor `Rule-Obj-FieldValue` untuk `FlagProrate` dan `.Type` | 1 — ekspor, bukan jejak | admin Pega | N-15, sisa N-04 | menggantikan "jalankan Tracer untuk melihat nilai yang lewat" |
| B5-3 | Pencarian nama `komiteAccept_ticket` pada ruleset di luar dua folder | 1 — pencarian teks pada ekspor tambahan | admin Pega | sisa N-14 | menggantikan "jalankan Tracer dan tunggu ticket naik" |
| B5-4 | Ekspor `PostEmailKomiteCNP` | 1 | admin Pega | G-13 bagian P1 | — |
| B5-5 | Ekspor ulang dua activity inti dengan parameter metode | 1 | admin Pega | sisa G-01, Q-9 | — |
| B5-6 | Ekspor rule `RDB-List` pengisi halaman `TglProd` | 1 | admin Pega | sisa Q-4 | — |
| B5-7 | Perbandingan isi badan `PROC_GENERATE_SEQUENCE_NUMBER` terhadap basis data berjalan | 3 — source objek basis data | DBA | C-01, dan lima penutupan yang bersandar padanya | — |
| B5-8 | Satu `SELECT` distribusi nilai `.Type` pada adjustment yang pernah masuk komite | 2 — satu `SELECT` baca-saja | DBA | pelengkap N-04: apakah nilai selain `'2'` dan `'3'` pernah terpakai | menggantikan "Tracer pada `CreateChildKomiteCNP_Act`" |

**B5-1 dikerjakan ronde ini**, di `INVENTARIS-BUKTI.md` §5. Sisanya tidak memblokir
penulisan spesifikasi mekanika keputusan, layar, pembentukan sirkulasi, maupun penomoran.

---

## 4. Lubang yang bertahan di balik pagar

| # | Lubang | Pagar | Baris inventaris | Pembuka |
|---|---|---|---|---|
| L5-1 | `SendEmailKlaimRejectClose_KMT.xml` membaca `childPageKomite.TransferType` sedangkan tak satu pun jalur tutup/tolak menulisnya — akibatnya terhadap isi surat tidak dinilai | PG-06 (A-6) | `INVENTARIS-BUKTI.md` §2.3 — `PostEmailKomiteCNP` tidak ada, isi surat tidak dapat dibaca utuh | dibuka bersama isi A-6 |
| L5-2 | Apakah muatan yang dikirim ke Kasir membawa penanda jalur yang sama | PG-04 (A-4) | `INVENTARIS-BUKTI.md` §2.5 baris 1 — nol baris `POOLDATA.DIRECTTOKASIR_LOG` | satu `SELECT` atas log kiriman |
| L5-3 | Apakah baris akseptasi menyimpan urutan jenjang, sehingga urutan lama dapat dibuktikan dari data | PG-05 (A-5) | `INVENTARIS-BUKTI.md` §2.1 — `PEGA_JSON_OS_AKSEP_KLAIM` tidak ada | source ketiga prosedur akseptasi |

Ketiganya dicatat **tanpa perlakuan dan tanpa alasan yang menembus pagar**. L5-3 khususnya
menggoda: ia dapat membuktikan Q5-1 dari data, tetapi menariknya ke dalam argumen berarti
memakai A-5 sebagai alasan, dan alasan pun tunduk pada pagar.
