# ADR-0049 — Jenis addendum adalah satu sumbu; materialitas diturunkan

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Sistem lama menyimpan dua penanda pada setiap addendum: sebuah penanda keadaan dengan tiga nilai,
dan sebuah penanda materialitas. Keduanya di-AND di dalam kondisi aturan, yang membuat *kode*
memperlakukannya sebagai dua besaran independen.

Perdebatan selama penggalian: apakah dua kolom itu dua gagasan bisnis, atau satu gagasan yang
tumbuh dua kali. Uji data yang dirancang untuk memutuskannya (Uji A, B, C) **tidak dijalankan**.

## Keputusan

**Satu sumbu.** Jenis addendum punya tiga nilai:

| Jenis | Materialitas turunannya |
|---|---|
| perubahan estimasi | material |
| penyesuaian ke nilai aktual | material |
| administratif | non-material |

**Kolom materialitas tidak ada** di model baru. Ia diturunkan dari jenis.

## Dasar — penalaran, bukan verifikasi

Kombinasi *(administratif, material)* tidak punya arti bisnis apa pun, dan tidak ada apa pun di
sistem lama yang mencegahnya terbentuk. Dua kolom yang dapat saling bertentangan, dengan sedikitnya
satu kombinasi yang tak bermakna, adalah tanda **satu gagasan yang tumbuh dua kali** — bukan dua
sumbu yang saling melengkapi.

Bahwa kode meng-AND keduanya membuktikan kode memperlakukannya independen. Itu **bukan** bukti
bisnis punya dua gagasan.

## Data warisan

Nilai materialitas lama **disimpan sebagai nilai warisan** dengan asal-usulnya (ADR-0042). Ia tidak
dipakai menurunkan apa pun dan tidak menjadi kolom hidup. Aturan penurunan yang baru **tidak
dijalankan** atas baris warisan — menjalankannya akan mengklasifikasi ulang sejarah.

## Syarat pembalikan

Bisnis menyebut **satu kasus nyata** di mana addendum administratif harus dinyatakan material, atau
menyebut pembeda material yang bukan tentang akibat premi.

**Bila terbalik:** kolom materialitas kembali sebagai sumbu kedua yang harus dijaga konsisten
selamanya. Ini salah satu dari enam keputusan yang **mengubah bentuk model**.
