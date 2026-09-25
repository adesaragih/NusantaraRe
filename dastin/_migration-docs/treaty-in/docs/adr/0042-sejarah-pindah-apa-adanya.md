# ADR-0042 — Sejarah pindah apa adanya: baris warisan dan aturan sentuh-perbaiki

**Status:** diterima, 23 September 2026
**Berlaku untuk:** migrasi data Treaty In

## Konteks

Model baru menetapkan invarian yang tidak dipenuhi sebagian data lama — antara lain "potongan tidak
boleh melebihi bruto yang menjadi dasarnya". Pilihannya: bersihkan sejarah supaya memenuhi,
tinggalkan sejarah di sistem lama, atau pindahkan apa adanya.

Membersihkan sejarah agar memenuhi invarian berarti **mengubah angka yang pernah disetujui dan
dibukukan** — persis penyakit yang seluruh sesi ini dipakai untuk mendiagnosis. Menyembuhkannya
dengan melakukannya sendiri, dalam skala besar, sekali jalan, atas seluruh sejarah, tidak bisa
dibenarkan.

Meninggalkan sejarah di sistem lama juga gugur: Treaty In Adjustment ditetapkan mengambil data lama
lewat SELECT, dan ADR-0016 melarang database link. Kontrak lama yang tidak ikut pindah tidak akan
bisa dijangkau sama sekali.

## Keputusan

Seluruh kontrak dan addendum pindah **lengkap dan apa adanya, termasuk cacatnya**, dengan tiga
aturan:

1. **Penanda warisan.** Baris hasil migrasi ditandai sebagai warisan, dan penanda itu menyatakan:
   *diterima apa adanya, tidak divalidasi terhadap invarian yang berlaku sekarang*. Penandanya
   terbaca oleh siapa pun yang membaca datanya, bukan tersembunyi di tabel teknis.
2. **Invarian ditegakkan pada penulisan, bukan pada baris.** Baris warisan ditulis sekali oleh
   migrasi dan tidak bisa disunting, jadi ia tidak pernah melanggar apa pun.
3. **Aturan sentuh-perbaiki.** Begitu sebuah kontrak warisan **disentuh** — addendum, revisi,
   adjustment, apa pun yang menulis — ia harus memenuhi invarian saat itu juga, dan penanda
   warisannya dicabut. Yang tidak disentuh tetap warisan selamanya.

## Materialitas addendum historis

Aturan penurunan materialitas yang baru **tidak dijalankan** atas baris warisan. Menjalankannya
akan **mengklasifikasi ulang sejarah**, dan itu bentuk lain dari mengubah masa lalu.

- `EDMMATERIALTYPE` diambil **apa adanya** sebagai nilai warisan yang tidak diverifikasi.
- Kolom itu membawa **asal-usulnya**: diturunkan oleh aturan, atau diimpor sebagai warisan.

Satu kolom memang akan menyimpan sesuatu yang aturannya sendiri tidak menghasilkan, dan itu benar
selama kolom itu jujur menyatakan dari mana nilainya datang. Yang salah bukan menyimpan nilai
warisan; yang salah adalah menyimpannya seolah ia hasil aturan.

## Konsekuensi

- **Tidak ada anggaran pembersihan di muka.** Biayanya tersebar dan hanya dibayar untuk kontrak
  yang memang masih dipakai.
- Yang diperbaiki adalah yang **diperiksa orang**, bukan yang ditebak skrip.
- Ada kemungkinan **tiga generasi data**, bukan dua: sebagian data sekarang sendiri hasil migrasi
  sebelumnya dari tabel yang lebih lama. Penanda warisan perlu menyatakan generasi mana.
