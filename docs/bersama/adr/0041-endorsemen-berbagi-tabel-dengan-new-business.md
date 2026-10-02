---
status: accepted
tanggal: 2026-09-23
sumber: spec Endorsement Life, keputusan work owner
---

# Endorsemen Life berbagi tabel dengan new business dan dibedakan kolom versi

`[keputusan work owner]` Endorsemen **tidak punya tabel sendiri**. Ia berbagi tabel dengan new
business, dan dibedakan oleh **kolom pembeda versi**.

Pencocokan baris antar versi lewat **penunjuk ke baris induk**, bukan lewat urutan indeks.

## Kenapa bukan tabel terpisah

1. Bentuk datanya sama. Tabel terpisah berarti bentuk yang sama dipelihara di dua tempat.
2. Pembaca hilir membaca satu tempat, bukan menggabungkan dua tabel yang bentuknya kebetulan sama.
3. Rantai versi menjadi hubungan di dalam satu tabel, bukan jembatan antar tabel.

## Kenapa penunjuk induk, bukan indeks

`[penyimpangan sadar]` Urutan baris **dapat berubah** antar versi. Pencocokan lewat indeks putus
persis ketika urutannya bergeser - dan pergeseran itu tidak menimbulkan galat apa pun, hanya angka
yang salah pasang.

## Hubungannya dengan modul lain

ADR-0018 menetapkan bentuk rantai generasi lewat penunjuk, dan ADR-0019 menetapkan pemasangan
baris anak lewat nomor urut. Catatan ini adalah bentuk yang **berbeda** untuk persoalan yang sama,
pada modul yang berbeda.

Perbedaan itu **disengaja dan dicatat**: modul ini berbagi tabel dengan new business, sedangkan
modul treaty punya tabel sendiri untuk tiap generasi. Menyeragamkan keduanya sesudah keduanya
berjalan akan lebih mahal daripada mencatat perbedaannya di sini.

## Data lama tidak disimpan sebagai salinan

`[keputusan work owner]` Halaman kerja yang menampilkan keadaan sebelum dan sesudah **tidak
di-persist**. Data lama dibaca dari **versi sebelumnya** lewat penunjuk induk.

## Akibat

1. Nol tabel salinan data lama.
2. Setiap pembacaan wajib menyaring pada kolom pembeda versi.
3. Test wajib: urutan baris bergeser antar versi, dan pasangannya **tetap benar**.
