---
status: accepted
tanggal: 2026-09-23
sumber: spec penyimpanan NB Treaty In dan spec Claim Prop, keputusan work owner
---

# Nama skema ditulis eksplisit pada setiap query

`[keputusan work owner]` Setiap query menyebut **nama skema secara eksplisit**. Tidak ada query
yang bergantung pada skema bawaan sesi.

## Kenapa

`[terverifikasi]` Basis data yang dipakai memuat **lebih dari satu skema**, dan sebagian besar
query di sistem lama menulis nama tabel **tanpa awalan skema**. Query semacam itu benar hanya
selama sesi kebetulan menunjuk skema yang tepat.

Ketergantungan itu tidak terlihat di teks query. Ia baru muncul sebagai kesalahan ketika
penyetelan sambungan berubah - biasanya saat pindah lingkungan, dan biasanya jauh dari orang yang
menulisnya.

## Akibat

1. Nol nama tabel tanpa awalan skema di dalam teks query.
2. Nama skema adalah **konfigurasi**, bukan tulisan tetap yang tersebar - sehingga satu tempat
   yang berubah bila skema dipindah.
3. Test yang menemukan query tanpa awalan skema **gagal**.
