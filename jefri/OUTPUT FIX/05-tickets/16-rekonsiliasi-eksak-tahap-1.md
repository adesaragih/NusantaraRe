# 16: Kerangka rekonsiliasi eksak — tahap 1

**What to build:** Sebuah pembanding yang menjalankan kasus terekam lewat perhitungan sistem baru dan
menyatakan, tanpa ruang tafsir, apakah hasilnya **identik sampai digit terakhir** dengan sistem lama.

Ini ukuran keberhasilan migrasi, bukan pelengkap. Karena seluruh jalur tulis produksi melewati
prosedur basis data yang isinya baru sebagian kami miliki, **membaca kode saja tidak dapat
membuktikan port-nya benar** — hanya membandingkan keluaran dua sistem atas masukan yang sama yang
bisa.

**Tahap 1 = perhitungan murni atas masukan terekam.** Tidak menyentuh alur, tidak menyentuh jalur
tulis, tidak menunggu ekspor produksi tunggal. Tahap berikutnya (per-modul dengan fixture tangga, lalu
end-to-end) menyusul setelah prasyaratnya tiba.

⛔ **Toleransi ditolak.** Bila urutan operasi dan presisi per-langkah direproduksi apa adanya, nilai
desimal bersifat deterministik dan hasilnya **harus** identik. Selisih sekecil apa pun berarti ada
salah-port — dan toleransi hanya menyembunyikannya, justru pada fase yang dirancang untuk
menemukannya.

⛔ **Kerangka ini tidak pernah menyentuh berkas kasus mentah.** Masukannya **hanya** fixture hasil
tiket 15. Menunjuk langsung ke berkas asal akan menarik data pelanggan ke dalam jalur pengujian.

**Blocked by:** 15, 04, 05, 06, 11

**Status:** ready-for-agent

- [ ] Pembanding membaca **fixture ter-de-identifikasi saja**; tidak ada jalur yang membuka berkas kasus mentah
- [ ] Perbandingan **eksak, tanpa toleransi**; tidak ada parameter epsilon di mana pun
- [ ] Selisih dilaporkan dengan menyebut **kasus, lini bisnis, dan langkah** tempat angkanya mulai berbeda — bukan hanya "tidak cocok"
- [ ] Mencakup lini bisnis yang rumusnya sudah ada (tiket 04–06) dan nilai dasar akseptasi (tiket 11)
- [ ] Kasus yang lini bisnisnya belum punya rumus dilaporkan sebagai **belum tercakup**, bukan lulus
- [ ] Hasilnya dapat dijalankan ulang dan deterministik
- [ ] Dicatat bahwa ini **tahap 1**, dengan prasyarat tahap berikutnya disebut eksplisit
