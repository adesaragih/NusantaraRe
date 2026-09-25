# ADR-0034 — Arsip JSON sistem lama tidak punya jalur baca

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Sistem lama menyimpan seluruh halaman clipboard sebagai satu dokumen JSON di `M_TREATY_IN.JSONDATA`
(dan `M_TREATY_IN_EDM.JSONDATA` untuk addendum). Dokumen itu adalah bentuk kanonik sistem lama.

Saat migrasi, dokumen-dokumen itu akan tetap ada. Godaannya jelas: menyimpannya "untuk jaga-jaga"
dan membiarkan aplikasi baru membacanya ketika ada yang tidak ketemu di model baru.

## Keputusan

Arsip JSON sistem lama **tidak punya jalur baca apa pun** di aplikasi hasil migrasi. Bukan
"dihindari", bukan "hanya untuk kasus khusus" — **tidak ada kode yang membacanya**.

Arsip itu punya **masa simpan** yang ditetapkan, dan setelah masa itu lewat ia dihapus.

## Alasan

Jalur baca cadangan ke bentuk lama menciptakan bentuk kanonik kedua, yang dilarang ADR-0023. Lebih
buruk lagi, ia menyembunyikan kegagalan migrasi: data yang tidak ikut pindah tidak akan pernah
ketahuan selama aplikasi masih bisa mengambilnya dari arsip.

Tanpa jalur baca, setiap data yang dibutuhkan **harus** ada di model baru, dan ketidakhadirannya
muncul sebagai kegagalan yang terlihat.

## Konsekuensi

- Inventaris struktur data harus lengkap sebelum migrasi, karena tidak ada jaring pengaman.
- Arsip disimpan di luar jangkauan aplikasi, dengan akses bernama untuk keperluan forensik saja.
- Masa simpan harus ditetapkan oleh pemilik proses dan kepatuhan, bukan oleh tim teknis.
