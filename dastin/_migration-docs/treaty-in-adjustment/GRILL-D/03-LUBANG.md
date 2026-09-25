> Modul  : Treaty In Adjustment · Ronde D · 2026-09-24
> Sifat  : TAMBAH-SAJA. Lubang dilaporkan, tidak ditambal.

# 03 · LUBANG RONDE D

Setiap lubang membawa **siapa yang menutupnya** dan **apa yang membuatnya ditagih**.

## LD-1 — Kolom dokumen akan KOSONG untuk seluruh baris warisan

| | |
|---|---|
| **Apa** | `GRL-19` melahirkan `DOKUMEN_ADDENDUM`, dan `TD-01` membuktikan sistem lama **tidak pernah merekam nomornya** |
| **Akibat** | setiap versi warisan berujung pada rujukan dokumen yang **kosong**, dan tidak ada cara mengisinya kecuali orang mengetik ulang dari arsip kertas |
| **Kenapa ini bukan cacat rancangan** | entitasnya **BARU**; tidak ada yang hilang. Yang perlu diketahui adalah bahwa **pengisiannya pekerjaan orang, bukan pekerjaan migrasi** |
| **Siapa menutup** | **pemilik proses** — ia keputusan tentang apakah pengisian ulang itu sepadan |
| **Yang menagih** | butir a `GRL-19` menjadikan kolomnya **boleh kosong**, sehingga migrasi tetap berjalan. Penagihnya laporan pertama yang mencoba mengelompokkan versi menurut dokumennya dan mendapati seluruh warisan berkelompok "tanpa dokumen" |

## LD-2 — Keunikan nomor dokumen dipilih GLOBAL tanpa data

| | |
|---|---|
| **Apa** | `GRL-19` butir c memilih keunikan **global**, dan arah salahnya **menolak data yang sah** |
| **Kenapa tetap dipilih** | dokumen **berdiri sendiri**; keunikan per cedant **tidak dapat dikompilasi** pada entitas yang tidak berinduk cedant |
| **Siapa menutup** | **pemilik data** — `UA-19` |
| **Yang menagih** | `UA-19` berada di daftar uji yang sama dengan uji migrasi lain; bila ia menemukan pengulangan, keunikannya **dilingkupi**, bukan dicabut |

## LD-3 — Titik beku materialitas tidak terbaca dari ekspor

| | |
|---|---|
| **Apa** | keterangan *"tidak boleh diubah"* dapat berarti **field**-nya atau **pilihannya sendiri**. `TD-02` membuktikan yang pertama ada; ia **tidak** membuktikan yang kedua tidak ada |
| **Akibat** | bila materialitas beku sejak diajukan, ia sejalan dengan jenis (`GRL-18` bagian 2). Bila tidak, **kedua sumbu punya titik beku berbeda**, dan itu keputusan tersendiri |
| **Siapa menutup** | **bisnis** — `DB-20`, ditulis sebagai **satu fakta satu butir** |
| **Yang menagih** | `GRL-20` butir b menyatakan terang bahwa ia **tidak dipilih sendiri**; siapa pun yang menulis invariannya akan menemukan titik bekunya belum ada |

## LD-4 — Baris warisan mungkin melanggar invarian baru, dan angkanya belum ada

| | |
|---|---|
| **Apa** | `TDA-10` — di sistem lama materialitas **tidak pernah ditegakkan di sisi simpan**. Penguncian layar dapat dilewati |
| **Akibat** | `INV-69` dan `INV-70` **tidak dapat langsung ditegakkan atas data lama**. Perlakuan atas yang melanggar — ditolak, ditandai, atau dinaikkan menjadi Material — **tidak dapat diputuskan sebelum angkanya ada** |
| **Siapa menutup** | **pemilik data** lewat `UA-3` yang sudah ditulis ulang kepalanya, lalu **pemilik proses** untuk perlakuannya |
| **Yang menagih** | uji positif `N-7c+` menuntut versi warisan yang melanggar **tetap dapat dimuat**; menulis uji itu memaksa pertanyaannya muncul |

## LD-5 — `TDA-10` kembali terbuka, dan sebabnya terbalik dari dugaan

| | |
|---|---|
| **Apa** | `TDA-10` sempat ditutup oleh `GRL-12` dengan alasan *"materialitas menjadi turunan, tidak ada pernyataan yang perlu ditegakkan"*. `GRL-12` batal, maka penutupnya jatuh |
| **Yang berubah, dan ini bukan sekadar membuka kembali** | jawaban C menyatakan penguncian layar itu **kehendak bisnis**. Maka `TDA-10` bukan lagi *"aturan yang seharusnya tidak ada ditegakkan di tempat yang salah"*, melainkan **"aturan yang benar ditegakkan di satu-satunya tempat yang dapat dilewati"** |
| **Ditutup oleh** | `GRL-20` butir 3 — aturannya naik menjadi invarian yang ditegakkan saat simpan |
| **Yang tersisa** | bukan keputusannya, melainkan **berapa data lama yang terlanjur melanggar** — `LD-4` |

## LD-6 — Angka gerbang 0 tidak dapat direproduksi

| | |
|---|---|
| **Apa** | `GERBANG-0-TO-SPEC.md` menyebut 220 kondisi di 18 seksi; sapuan ronde ini memberi **758 kemunculan, 18 kondisi berbeda, 26 seksi**, dan tiap angka per seksi **tepat dua kali** angka gerbang 0 |
| **Dugaan yang paling sesuai** | gerbang 0 menyapu **satu ekspor**, bukan keduanya; dan **"18 seksi"** kemungkinan sebenarnya **18 kondisi berbeda** |
| **Kenapa dilaporkan alih-alih dikoreksi diam-diam** | arah klaimnya tidak berubah, tetapi angka yang tidak dapat direproduksi **tidak boleh dipakai sebagai bukti** |
| **Siapa menutup** | sesi mana pun yang menyentuh `GERBANG-0-TO-SPEC.md` berikutnya |
| **Yang menagih** | `TD-03` memuat kedua angka berdampingan |

---

## Pembersihan angka 220/18 — apa yang disentuh, dan apa yang TIDAK

**Sapuan:** seluruh `.md` kedua modul. **Ditemukan: 8 kemunculan.**

### Disunting — 4 kemunculan di 3 berkas

| Berkas | Kemunculan | Diganti dengan |
|---|---:|---|
| `GERBANG-0-TO-SPEC.md` | 2 | rujukan ke `TD-02`, dengan blok pencabutan bertanggal |
| `PENGETAHUAN-PENGGRILL-ADJUSTMENT.md` | 1 | rujukan, dan ditambah bahwa materialitas **sakelar dua arah** |
| `PENGETAHUAN.md` | 1 | blok koreksi bertanggal; angkanya **dicabut, bukan diganti angka baru** |

### TIDAK disunting — 4 kemunculan, masing-masing dengan sebabnya

| Berkas | Baris | Kenapa tidak disentuh |
|---|---:|---|
| `GRILL-A/01-TEMUAN.md` | 165 | **berkas ronde yang sudah DITUTUP, dan sifatnya TAMBAH-SAJA.** Menyunting temuan ronde tertutup menghapus jejak apa yang diketahui saat putusannya diambil. Koreksinya hidup di `GRILL-D`, dan indeks menunjuk keduanya |
| `GRILL-D/01-TEMUAN.md` | 118 | ia **melaporkan** selisihnya — angka lama harus tetap terbaca di sebelah angka baru |
| `GRILL-D/03-LUBANG.md` | 58 | idem — `LD-6` |
| `GRILL-D/07-AUDIT.md` | 41 | idem — `MA-13` |

> **Tidak ada angka baru yang disebarkan.** Keempat suntingan mengganti angka dengan **rujukan**,
> sebab angka yang belum stabil tidak boleh beredar sebagai bukti — dan itu justru kesalahan yang
> `MA-13` catat.
