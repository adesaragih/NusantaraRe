# ADR-0037 — Masukan versus turunan; tidak ada "mode manual"

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In, dan diusulkan untuk seluruh modul

## Konteks

Pola yang sama muncul tiga kali:

- **Penyebaran** punya jalur otomatis dan jalur manual, dipilih oleh keterisian sebuah field yang
  bisa mengosongkan dirinya sendiri.
- **Reinstatement** punya empat field yang saling memicu dalam satu cincin: persen menggerakkan
  jumlah, jumlah menggerakkan persen yang lain.
- **Kapasitas proporsional** punya jalur "man" untuk Surplus yang diketik dan jalur turunan,
  dengan sebuah sakelar yang di dua pertiga titik pemanggilan diberi makan nomor layer.

Sistem lama tidak pernah memutuskan mana masukan dan mana turunan. Setiap kali sebuah turunan
ternyata perlu disepakati, jawabannya adalah **membuka suntingan** — dan lahirlah cincin, jalur
manual, dan penimpaan. Semuanya gejala dari satu keputusan yang tidak pernah diambil.

## Keputusan

1. Besaran yang **dinegosiasikan** adalah **masukan** dan diketik orang.
2. Besaran yang **dihitung** adalah **turunan**, tidak bisa disunting, dan tidak pernah menulis
   balik ke masukannya.
3. **"Mode manual" tidak ada.** Bila sebuah besaran boleh disepakati, ia **dinaikkan menjadi
   masukan** — bukan dibiarkan sebagai turunan yang boleh disunting.

## Penerapan yang sudah ditetapkan

| Besaran | Sifat |
|---|---|
| Quota Share: persentase QS | masukan (satu bilangan; retensi dan sesi turunannya) |
| Surplus: retensi sebagai jumlah uang, dan jumlah lines | **dua masukan**; kapasitas turunan |
| Reinstatement: porsi limit yang dipulihkan, tarif premi reinstatement | dua masukan, dinegosiasikan terpisah; hubungannya perkalian |
| Jumlah limit dipulihkan, premi tambahan | turunan |
| Potongan berbasis premi bruto | turunan |
| Potongan berjumlah disepakati | masukan, dan **tidak diprorata** |
| Penyebaran per baris | masukan bila diisi tangan; turunan bila disemai dari acuan baku |
| Faktor prorata addendum | turunan dari tanggal, bukan centang |
| "Boleh diubah" | turunan dari keadaan siklus hidup |

## Konsekuensi

- **Prorata mengenai dasarnya, bukan hasilnya.** Ada satu besaran yang menyusut — paparan waktu
  atas premi. Prorata itu sekali, lalu seluruh turunan dihitung ulang dari dasar yang baru. Dua
  puluh blok prorata terpisah menjadi satu perkalian, dan pertanyaan "apa yang ikut menyusut"
  lenyap bersama bentuk yang melahirkannya.
- Cincin `RetentionPct`/`CessionPct` tidak diangkut. Satu-satunya akibat nyata dari cincin itu
  adalah membuka pintu suntingan.
- `CessionPct` **pecah**: di Quota Share ia persentase sesi yang bermakna; di Surplus yang ada
  adalah jumlah lines. Keduanya tidak berbagi kolom maupun nama. Pemecahan ini berlaku **tanpa
  syarat**, bukan bergantung pada jawaban soal manual.
