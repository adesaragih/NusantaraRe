---
status: accepted
tanggal: 2026-09-23
sumber: spec Komite Claim Fac In ID-3, keputusan work owner
---

# Tangga persetujuan berakhir karena keadaan jenjang, bukan karena pencacah

`[keputusan work owner]` Tangga persetujuan berakhir bila **seluruh jenjang menyetujui**. Satu
menolak, kasus **tutup seketika**.

Syarat berakhir itu **dinyatakan eksplisit** dan dapat diuji.

## Apa yang diganti

`[terverifikasi]` Di sistem lama penugasan berhenti karena penentu penerima **tidak menghasilkan
siapa pun** ketika setiap jenjang sudah memutuskan. Itu **perilaku diam**: benar dalam praktik,
tetapi tidak dapat diuji dan tidak dapat dijelaskan.

`[penyimpangan sadar]` Sistem baru menyatakan syaratnya, bukan menyimpulkannya dari ketiadaan
penerima.

## Pencacah jenjang bukan kebenaran

Pencacah jumlah jenjang **boleh ada sebagai tampilan**, tetapi **bukan** dasar keputusan.

`[terverifikasi]` Alasannya nyata: daftar jenjang **dapat berubah di awal** - pengaju dikeluarkan
dari daftar - dan pencacah tidak ikut tahu. Tangga yang berhenti berdasarkan pencacah akan
berhenti di tempat yang salah persis pada kasus yang paling perlu benar.

## Akibat

1. Keputusan berakhir diambil dari **keadaan setiap jenjang**, bukan dari cacah.
2. Test wajib mencakup: daftar jenjang yang menyusut di awal, dan tangga yang tetap berakhir pada
   jenjang yang benar.
3. Penolakan menutup kasus seketika - bukan melanjutkan ke jenjang berikutnya.
