# Usulan untuk tim inti — presisi kolom angka tabel flat Fac In

> **Usulan, bukan keputusan.** Presisi tabel flat adalah keputusan tim inti (butir 66; keputusan work owner
> 02-10-2026 butir 68.2, diteruskan sesi `nusantarare-0f`). Ditulis agent modul `nbfacin`, 2 Oktober 2026. Tiket 23
> (berkas migrasi 180–219) tetap **ditahan** sampai ada keputusan.

## Masalahnya

`[terverifikasi]` DDL draf (`D:\migrasi\RNM\OUTPUT\08-flat\DDL-tabel-flat-draf.sql`) memakai **415** kolom `NUMBER` tanpa
presisi, 1 `NUMBER(3)`, 50 `NUMBER(5)`. Penjaga inti `TestNolNumberTanpaPresisi` hanya mengizinkan `NUMBER(38,8)`,
`(5)`, `(10)`, `(19)`. Mengikuti penjaga berarti uang dan rate menjadi `NUMBER(38,8)`. Data lama tidak muat di sana.

## Ukurannya

Diukur 02-10-2026 atas **keluaran `loader.Flatten`** (jadi cabang `OldData`/`Total*` dkk. sudah dibuang; delapan kolom
teks butir 68.1 tidak dihitung). Dua jendela: **5 fixture** de-identifikasi di repositori, dan **115 contoh XML**
`D:\migrasi\RNM\DDL\CONTOH` yang diubah ke JSON di luar repositori lalu dihapus.

**A. Lebih dari 8 desimal BERMAKNA** — yang akan **dibulatkan** `NUMBER(38,8)`. Nol di belakang tidak dihitung
(`0.00000000000000000000` bukan temuan). ⚠️ Putaran pertama pengukuran ikut menghitung nol di belakang dan memberi angka
lebih besar; angka di bawah adalah putaran kedua (koefisien dinormalkan dulu).

| Jendela | Nilai | Kolom | Terbanyak |
| --- | ---: | ---: | --- |
| 5 fixture | **458** | **20** | `T_SPREADINGLIST.PREMIUM_SPREADED` 95, `SHARE_PERCENTAGE` 90, `TSI_SPREADED` 90, `T_COVERAGELIST.PREMIUM` / `PREMI_NUSANTARA_RE` 50 |
| 115 contoh korpus | **10.876** | **82** | `T_SPREADINGLIST.PREMIUM_SPREADED` 2.249, `T_COVERAGELIST.PREMI_NUSANTARA_RE` 2.241, `T_COVERAGELIST.PREMIUM` 2.224, `T_COVERAGELIST.RATE` 1.813, `T_PROPERTYITEMLIST.TOTAL_PREMIUM_NUSANTARA_RE` 440 |

**B. Lebih dari 38 digit BERMAKNA** — melampaui presisi `NUMBER` Oracle (1–38 digit bermakna).

| Jendela | Kolom | Nilai |
| --- | --- | ---: |
| 115 contoh korpus | `T_CURRENCYLIST.PAY_NET_PREMIUM` | **3** |
| 5 fixture | `T_CURRENCYLIST.PAY_NET_PREMIUM` (`edm-fire-1`) | **1** |

> ⚠️ **Ralat butir 69** (verifikasi independen sesi `nusantarare-0f`, diukur ulang agent 02-10-2026). Versi pertama
> tabel ini menghitung **koefisien mentah, termasuk nol di belakang koma**, dan memuat **217** nilai korpus di sembilan
> kolom (`T_COVERAGELIST.NET_RATE` 185, tiga kolom `T_CURRENCYLIST` 8 masing-masing, lima kolom `T_FR_*` 8 bersama) dan
> 3 nilai fixture. Dengan nol ujung dibuang hanya **3 + 1** di atas yang sungguh > 38 digit bermakna. Kalimat lama
> "`NUMBER` polos pun membulatkannya" karena itu **keliru untuk 214 dari 217** nilai itu: `[dugaan]` — tidak diuji ke
> Oracle — `NUMBER` menyimpan menurut digit **bermakna**, sehingga nol di ujung tidak memakan presisi. Hitungan
> `Diagnostik.LebihDari38Digit` di Flatten diralat dengan cara yang sama (digit bermakna).

Desimal terpanjang yang **tertulis** (termasuk nol ujung): **60** (`PAY_NET_PREMIUM`, `T_FR_PAYMENT.NET_PREMIUM`).

## Usulan (keputusan work owner 68.2 meminta usulan ini ditulis)

1. Kolom uang/rate yang nilainya **dapat** melebihi 38 digit atau 8 desimal bermakna disimpan sebagai **teks mentah**
   apa adanya dari dokumen — tidak dibulatkan, tidak dinormalkan.
2. Bila sebuah kolom perlu dihitung atau diindeks sebagai angka di Oracle, tambahkan **kolom `NUMBER` turunan** di
   sampingnya, dengan presisi yang tim inti tetapkan dan aturan pembulatan yang tertulis. Teks mentah tetap sumber
   kebenaran.
3. Di dalam aplikasi nilai tetap **desimal eksak** (`apd`), diurai dari teks oleh satu fungsi bersama.

⚠️ **Preseden dan batasnya.** ADR-0034 *"uang melintasi batas procedure sebagai teks"* adalah preseden untuk **teks di
batas procedure** — hanya batas pemanggilan stored procedure. Tetapi tabel di ADR itu sendiri menulis bentuk **di kolom basis data: desimal berskala tetap**. Usulan
butir 1 karena itu **melampaui** ADR-0034 dan butuh ADR baru atau amandemennya — bukan penerapan ADR yang sudah ada.

**Akibat bagi `loader.Flatten` bila usulan diterima:** saat ini Flatten mengurai kolom `NUMBER` menjadi desimal eksak
dan **menormalkan koma desimal** (BAHAN §7) — teks mentahnya tidak disimpan. Kolom yang pindah ke teks mentah harus
dimasukkan ke daftar teks (seperti delapan kolom butir 68.1), sehingga teksnya tersimpan tanpa normalisasi.

## Yang perlu diputuskan tim inti

- Apakah usulan di atas diterima, atau kolom diberi presisi lain (misalnya `NUMBER` polos, dengan konsekuensi tabel B).
- Daftar kolom yang masuk butir 1: semua kolom uang/rate, atau hanya yang terukur melampaui batas (82 kolom di korpus,
  dari 115 contoh — **bukan** bukti bahwa kolom lain aman di produksi).
- Apakah `TestNolNumberTanpaPresisi` perlu pengecualian untuk tabel flat.

*Tanpa nama orang, tanpa nilai kasus — hanya nama kolom dan hitungan.*
