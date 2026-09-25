# R08: Rekonsiliasi kasus renewal

**What to build:** Kasus renewal yang dijalankan sistem baru menghasilkan angka yang **identik sampai
digit terakhir** dengan sistem lama — dibuktikan memakai kerangka rekonsiliasi yang sudah ada, **tanpa
pembanding baru**.

Ini yang menutup lingkaran: R01–R06 membangun pintu masuk dan layar; tiket ini membuktikan bahwa di
balik pintu itu, angkanya memang tidak berubah.

## Mengapa tidak ada pembanding baru

`[terverifikasi]` Renewal **tidak punya mesin hitung ulang**. Hanya empat activity yang menggerbangi
pembeda siklus — masa berlaku tanggal, validasi tanggal, proteksi spreading, input pembayaran — dan
**tidak satu pun perhitungan premi**. Keempatnya berkas bersama yang identik dengan NB.

Membangun pembanding kedua karena itu bukan sekadar mubazir — ia **berbahaya**: dua pembanding dapat
menyimpang, lalu keduanya tampak benar.

📌 Masukannya **sudah tersedia**: salah satu dari lima berkas kasus yang diurus tiket NB-15 adalah
kasus renewal nyata. Tiket ini memakainya lewat kerangka NB-16.

⛔ **Toleransi ditolak** (ADR-0001). Bila urutan operasi dan presisi per-langkah direproduksi apa
adanya, nilai desimal bersifat deterministik dan hasilnya **harus** identik.

⛔ **Tidak pernah menyentuh berkas kasus mentah** — hanya fixture ter-de-identifikasi hasil NB-15.

**Blocked by:** **NB-16** (kerangka rekonsiliasi eksak tahap 1) · R01 · R04

**Status:** ready-for-agent

- [ ] Kasus renewal dijalankan lewat **kerangka NB-16**; tidak ada pembanding baru yang dibuat
- [ ] Perbandingan **eksak, tanpa toleransi**; tidak ada parameter epsilon
- [ ] Masukannya **hanya fixture ter-de-identifikasi** (NB-15); tidak ada jalur yang membuka berkas kasus mentah
- [ ] Selisih dilaporkan dengan menyebut **kasus, lini bisnis, dan langkah** tempat angkanya mulai berbeda
- [ ] Nilai dasar akseptasi renewal terbukti **nilai pertanggungan penuh**, bukan selisih — sesuai keputusan yang sudah diambil
- [ ] Bila ada selisih, **keputusan membuang `pySaveSQL` `PROSESCOPY` (R06) diperiksa lebih dulu** sebagai tersangka
- [ ] Hasilnya deterministik dan dapat dijalankan ulang
