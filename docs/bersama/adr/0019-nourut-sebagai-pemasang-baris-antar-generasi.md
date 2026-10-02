---
status: accepted
tanggal: 2026-09-23
sumber: spec penyimpanan NB dan EDM Treaty In, keputusan work owner
---

# Baris anak dipasangkan antar generasi lewat nomor urut, bukan kunci dagang

Setiap tabel anak membawa **nomor urut baris** di dalam induknya. Nomor urut itulah yang
memasangkan baris generasi lama dengan baris generasi baru saat selisih dihitung.

**Bukan kunci dagang.** Kunci dagang - nama treaty, mata uang, nomor angsuran, lapisan - **boleh
berubah** di dalam endorsemen. Memakainya sebagai pemasang berarti pasangan putus persis ketika
endorsemen mengubahnya, yaitu ketika selisih paling dibutuhkan.

## Dua perilaku yang berbeda, sengaja

| | Polis baru | Endorsemen |
| --- | --- | --- |
| Menghapus baris | **boleh** | **tidak bisa** |
| Nomor urut sesudah penghapusan | **dinomori ulang rapat** - 3 baris, yang ke-2 dihapus, sisanya jadi 1 dan 2 | tidak berlaku |
| Baris baru | di belakang | di belakang, `maksimum + 1` |

Perbedaan ini bukan ketidakkonsistenan. Pada polis baru belum ada generasi sebelumnya, sehingga
penomoran ulang tidak memutus pasangan apa pun. Pada endorsemen, penomoran ulang **akan**
memutusnya.

## Aturan keutuhan

Generasi `n+1` wajib memuat **setiap nomor urut** yang ada pada generasi `n`. Generasi yang
kehilangan satu nomor urut **ditolak** - karena baris yang hilang berarti selisih yang tidak dapat
dihitung.

## Akibat

1. Nomor urut ditetapkan lapisan `repository` saat menulis, bukan dikirim dari layar.
2. Lapisan layanan menolak generasi yang kehilangan nomor urut, dan menolak penghapusan baris pada
   endorsemen.
3. Test wajib mencakup: tiga baris, yang tengah dihapus pada polis baru - sisanya bernomor 1 dan 2;
   dan generasi endorsemen yang kehilangan satu nomor urut - **ditolak**.
