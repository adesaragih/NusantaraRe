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

**Status:** ready-for-human — sebagian 01-10-2026: premi renewal 4/4 cocok lewat kerangka NB-16; nilai dasar akseptasi belum tercakup · ⏸ pengerjaan rnwfacin dihentikan work owner 01-10-2026 (`PROMPT-LANJUT-IMPLEMENT-NB-SAJA.md`)

- [x] Kasus renewal dijalankan lewat **kerangka NB-16**; tidak ada pembanding baru yang dibuat
- [x] Perbandingan **eksak, tanpa toleransi**; tidak ada parameter epsilon
- [x] Masukannya **hanya fixture ter-de-identifikasi** (NB-15); tidak ada jalur yang membuka berkas kasus mentah
- [x] Selisih dilaporkan dengan menyebut **kasus, lini bisnis, dan langkah** tempat angkanya mulai berbeda
- [ ] Nilai dasar akseptasi renewal terbukti **nilai pertanggungan penuh**, bukan selisih — sesuai keputusan yang sudah diambil
- [ ] Bila ada selisih, **keputusan membuang `pySaveSQL` `PROSESCOPY` (R06) diperiksa lebih dulu** sebagai tersangka
- [x] Hasilnya deterministik dan dapat dijalankan ulang

## Comments

### 2026-10-01 — sebagian (agent)

Kerangka NB-16 (`nbfacin/backend/services/rekonsiliasi`, `TestBerkasKasusP5`) menjalankan fixture `rnw-fire-1`: **4/4
coverage FIRE cocok sebagai teks persis**, tanpa pembanding baru dan tanpa toleransi. Belum: **nilai dasar akseptasi**
(kriteria 5) — nilai tangga lama (`CARID2`, `LetterNo`) tidak terekam di fixture, dilaporkan "belum tercakup". Kriteria
6 (tersangka `pySaveSQL`) tidak berlaku selama tidak ada selisih; R06 sendiri diusulkan wontfix (A41). Blocker formal R04
(layar) tetap.
