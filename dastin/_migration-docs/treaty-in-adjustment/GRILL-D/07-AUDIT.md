> Modul  : Treaty In Adjustment · Ronde D · 2026-09-24

# 07 · AUDIT RONDE D

## 1. Buku besar ronde

| | |
|---|---:|
| putusan **dibatalkan** | 2 — `GRL-04`, `GRL-12` |
| putusan **baru** | 2 — `GRL-19`, `GRL-20` |
| putusan diperiksa dan dinyatakan **tidak tersentuh** | 6 — `GRL-07`, `GRL-08`, `GRL-10`, `GRL-11`, `GRL-14`, dan `GRL-01` (dikoreksi, bukan dibongkar) |
| putusan **tersentuh sebagian** | 2 — `GRL-13` (lingkup dipertegas), `GRL-18` (satu bagian jatuh) |
| temuan ekspor | 4 — `TD-01` … `TD-04` |
| lubang | 6 — `LD-1` … `LD-6` |
| invarian diusulkan | 2 — `INV-69`, `INV-70`, masing-masing dengan uji negatif **dan** positif |
| butir daftar-untuk-dibantah baru | 3 — `DB-20`, `DB-16a`, `DB-16b` |
| uji data | 2 — `UA-19` baru; `UA-3` **ditulis ulang kepalanya** |

## 2. Kesalahan yang dicatat

### MA-12 — butir daftar-untuk-dibantah memuat EMPAT klaim dalam satu kalimat

**`DB-15`** disusun berbunyi *"Materialitas tidak dipakai di luar layar pengisian"* dengan empat
klaim di dalamnya: jalur persetujuan, dokumen ke cedant, akuntansi, dan laporan.

**Akibatnya terukur:** jawaban **"TIDAK"** tidak memberi tahu **klaim mana yang jatuh**, dan
pendalaman C terpaksa dikirim susulan. Satu putaran ke pemilik proses terbuang — dan satu putaran
ke orang diukur dalam hari.

> **Aturan yang berlaku seterusnya: SATU FAKTA SATU BUTIR.**
>
> Butir yang memuat lebih dari satu klaim **tidak dapat dijawab dengan "tidak"** — jawabannya tidak
> menunjuk apa pun. Bila sebuah pernyataan memuat kata **"dan"**, **"maupun"**, atau daftar, ia
> dipecah sebelum dikirim.

**Sudah diterapkan di ronde ini:** `DB-16` yang lama — juga empat klaim — dipecah menjadi `DB-16a`
dan `DB-16b` (`05-GERBANG` §5), dan `DB-20` ditulis sebagai satu klaim tunggal sejak awal.

### MA-13 — angka yang tidak dapat direproduksi sempat dipakai sebagai penguat

`GERBANG-0-TO-SPEC.md` memakai **"220 kondisi `pyDisabledWhen` di 18 seksi"** sebagai penguat
jawaban C. Sapuan ronde ini tidak dapat mereproduksinya: **758 kemunculan, 18 kondisi berbeda,
26 seksi**, dengan tiap angka per seksi **tepat dua kali** angka gerbang 0.

**Arah klaimnya tidak berubah** — penguncian memang nyata, dan 660 dari 758 menyebut
`EDMMaterialType`. Yang salah **bukan kesimpulannya, melainkan angkanya**, dan angka yang tidak
dapat direproduksi tidak boleh berdiri sebagai bukti.

> **Aturan:** setiap angka yang dipakai sebagai penguat menyebut **perkakas yang menghasilkannya**
> dan **ekspor mana yang disapu**. Gerbang 0 menyebut keduanya tidak.

## 3. Usulan aturan

### TA-10 — jawaban bisnis yang datang lewat jalur lain WAJIB masuk ronde

Tiga butir dijawab 24 September lewat jalur di luar sesi grilling. Jawabannya **tidak pernah masuk
ke sesi mana pun**, sementara indeks tetap mencatat pembatalnya **terbuka** dan keputusan yang
dibatalkannya tetap **terkunci** — selama sehari penuh, dan to-spec sempat dimulai di atasnya.

> **Jawaban atas butir daftar-untuk-dibantah adalah MASUKAN GRILLING, dari jalur mana pun ia
> datang.** Ia tidak menjadi keputusan dengan sendirinya; ia **menyalakan pembatal**, dan pembatal
> yang menyala menuntut ronde — sependek apa pun.
>
> Uji termurahnya satu perintah: **daftar butir `DB` yang sudah punya jawaban, dan bandingkan dengan
> daftar pembatal yang masih bertanda terbuka di indeks.** Selisihnya adalah pekerjaan yang belum
> dilakukan.

### TA-11 — pembatal berbentuk "A atau B" diperiksa PER CABANG

Pembatal `GRL-04` berbunyi *"lebih dari satu kontrak **atau** satu revisi menggabungkan lebih dari
satu dokumen"*. Cukup satu cabang untuk menyalakannya — tetapi **memeriksa keduanya memberi tahu
seberapa jauh putusannya jatuh**.

Di sini keduanya menyala, dan itu yang membuat jelas bahwa yang lahir adalah **benda di atas versi**,
bukan benda di bawahnya — sehingga `PENYESUAIAN` tetap tidak kembali.

## 4. Yang ronde ini TIDAK lakukan, dan itu disengaja

| | |
|---|---|
| tidak menulis spesifikasi, DDL, struct, komponen, atau tiket | fase grilling |
| tidak memutuskan **titik beku materialitas** | tidak terbaca dari ekspor — `DB-20` |
| tidak memutuskan **perlakuan atas baris warisan yang melanggar** | angkanya belum ada — `UA-3`, `LD-4` |
| tidak memperbarui `CABANG-K-PEMETAAN-TO-SPEC.md` | pekerjaan sesi to-spec; dan ia menunggu `DB-20`, `DB-16a`, `DB-16b` |
| tidak menyentuh `SPEC-MODEL-DATA.md` maupun `SPEC-INVARIAN.md` induk | `INV-69` dan `INV-70` **usulan**, bukan suntingan |
