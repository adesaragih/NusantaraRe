---
status: accepted
tanggal: 2026-09-23
sumber: spec Claim Fac In Bab 15 dan spec Claim Prop - keputusan work owner atas rekomendasi asisten
---

# Aturan peran ditetapkan sekali dan berlaku lintas modul

`[keputusan work owner]` Aturan peran sistem baru **ditetapkan sekali** dan berlaku lintas modul.
Ia **tidak** dirancang ulang per modul.

Penegakannya di **lapisan layanan**, bukan di layar.

## Kenapa ini bukan soal kerapian

`[terverifikasi]` Sensus korpus menemukan lubang yang sama pada **dua modul beruntun**: seluruh
medan hak akses ada di layar tetapi **kosong seluruhnya**. Ratusan kemunculan nama medan hak
akses, **nol** yang terisi. Satu-satunya yang terisi hanya menyebut **nama kelas**, bukan nama hak.

Artinya gerbang wewenang di sistem lama **hidup di layar, dan sebagian besar tidak hidup sama
sekali**. Merancang ulang per modul berarti menyalin lubang itu berkali-kali, dan menemukannya
kembali berkali-kali.

## Bentuk penegakan

1. Wewenang diperiksa **di lapisan layanan**, pada setiap jalur yang mengubah keadaan. Layar boleh
   menyembunyikan tombol, tetapi penyembunyian **bukan** penegakan.
2. Sumber peran adalah **satu tabel**, bukan nama orang yang ditanam di kode.
3. `[penyimpangan sadar]` Nama orang yang ditanam di kode sistem lama **dibuang**. Perilaku tidak
   berubah - orang yang sama tetap mendapat tingkat yang sama - tetapi pergantian pemegang jabatan
   cukup lewat baris tabel.

## Akibat

1. Nol nama orang di dalam kode maupun di dalam teks SQL.
2. Test wajib mencakup jalur yang **ditolak** karena wewenang, bukan hanya jalur yang berhasil.
3. Modul baru memakai aturan peran yang sudah ada. Bila ia menuntut peran baru, peran itu
   ditambahkan ke aturan bersama - bukan dibuat sendiri di dalam modul.
