---
status: accepted
tanggal: 2026-09-23
sumber: spec Claim Life dan spec penyimpanan NB Treaty In, keputusan work owner
---

# Tabel kerja adalah akar lintas-lini, bukan tabel di samping

Satu **tabel kerja** menjadi akar bagi lebih dari satu lini pekerjaan. Ia tidak berdiri di samping
pohon sebagai pelengkap - ia **akarnya**.

Baris dari lini yang berbeda hidup di tabel yang sama, dibedakan oleh **kolom pembeda lini**.

## Contoh yang sudah berjalan

| Tabel kerja | Menaungi |
| --- | --- |
| tabel kerja klaim | klaim **dan** kasus komite, dibedakan kolom pembeda |
| tabel kerja polis | polis treaty **dan** daftar premi |

`[terverifikasi]` Keduanya sudah dipakai lebih dari satu modul sebelum keputusan ini ditulis.
Catatan ini mengesahkan yang sudah terjadi, bukan mengusulkan yang baru.

## Kenapa

1. Objek kerja punya daur hidup, pemilik, dan jejak audit yang **bentuknya sama** lintas lini.
   Menduplikasinya per lini berarti menduplikasi daur hidup itu.
2. Penyerahan pekerjaan antar lini - misalnya klaim diserahkan ke komite - menjadi hubungan induk
   ke anak di dalam **satu** tabel, bukan jembatan antar tabel.
3. Laporan lintas lini tidak perlu menggabungkan tabel yang bentuknya berbeda-beda.

## Akibat

1. **Setiap modul baru memeriksa dulu** apakah tabel kerja yang ada sudah menaunginya, sebelum
   mengusulkan tabel kerja sendiri.
2. Kolom pembeda lini **wajib** ada dan **wajib** terisi. Baris tanpa pembeda tidak sah.
3. Penghapusan baris anak di dalam tabel kerja ditegakkan di lapisan layanan, bukan oleh kaskade
   basis data - sebab satu tabel menaungi lebih dari satu arti.
4. Test yang menemukan baris tanpa kolom pembeda **gagal**.
