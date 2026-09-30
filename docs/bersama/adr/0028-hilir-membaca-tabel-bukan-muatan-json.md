---
status: accepted
tanggal: 2026-09-23
sumber: spec Claim Life dan spec penyimpanan Treaty In, keputusan work owner
---

# Sistem hilir membaca tabel, bukan muatan JSON

`[keputusan work owner]` Sistem hilir - produksi, pelaporan, layanan luar - membaca **langsung
dari tabel** sistem baru.

Seluruh penyusun muatan JSON keluar yang ada di sistem lama **tidak dimigrasikan**.

## Kenapa

1. Muatan JSON di sistem lama adalah **akibat** dari dokumen JSON sebagai bentuk simpan. Begitu
   data tersimpan relasional, muatan itu kehilangan alasan keberadaannya.
2. Menyusun ulang JSON dari tabel berarti **menulis bentuk yang sama dua kali** - sekali sebagai
   tabel, sekali sebagai muatan - dan keduanya akan bercabang.
3. Pembaca hilir yang membaca tabel mendapat kolom bertipe benar. Pembaca muatan JSON mendapat
   teks, dan harus mengurai ulang.

## Yang tetap dipertahankan

Kontrak efek keluar - kapan dikirim, ke mana, apa yang menandakan berhasil - **tidak berubah**.
Yang berubah hanya **bentuk** bacaannya.

Tabel hilir yang sudah datar dan sudah dipakai sistem lain **tidak dibuat ulang**. Ia diisi apa
adanya.

## Akibat

1. Rule penyusun muatan JSON keluar tidak masuk lingkup migrasi. Ia masuk daftar yang **sengaja
   tidak dibangun**.
2. Tabel yang dibaca hilir membawa **kunci penyaring** supaya pembaca tidak perlu menggabungkan
   balik.
3. Perubahan bentuk tabel yang dibaca hilir adalah perubahan kontrak - diperlakukan sebagai
   perubahan luar, bukan perubahan dalam.
