# 07: Pembulatan di dalam loop akumulasi — kasus uji wajib

**What to build:** Bukti bahwa perhitungan premi tetap cocok **untuk daftar panjang**, bukan hanya
untuk satu nilai.

Ini tiket pengujian, dan ia berdiri sendiri karena alasan yang spesifik: di sistem lama, sebagian
pembulatan berada **di dalam loop akumulasi**, sehingga galatnya **menumpuk per iterasi**. Port yang
benar untuk satu nilai masih bisa meleset untuk daftar panjang — dan itu bentuk kesalahan yang
**paling sulit terlihat**, karena setiap nilai tunggalnya lulus.

Lingkupnya **lintas lini bisnis**: galat akumulasi tidak mengenal batas COB, jadi menguji satu COB
saja tidak membuktikan apa pun tentang yang lain.

Dua sifat sistem lama yang memperkuat kebutuhan ini: validasi sisa spreading menuntut kesamaan
**persis** (jumlah bagian = 100, jumlah premi tersebar = premi) tanpa distribusi galat pembulatan,
dan pembulatan dilakukan lebih dulu sebelum dibandingkan.

**Blocked by:** 04, 05, 06

**Status:** ready-for-agent

- [ ] Kasus uji akumulasi untuk **tiap lini bisnis** yang punya rumus di tiket 04–06
- [ ] Setiap kasus memakai daftar yang cukup panjang untuk memunculkan penumpukan galat — bukan dua atau tiga baris
- [ ] Pembulatan terjadi **di dalam** loop, di posisi yang sama seperti sumbernya — bukan sekali di akhir
- [ ] **Rekonsiliasi eksak** pada nilai agregat, bukan hanya pada tiap elemen
- [ ] Kasus uji yang membuktikan pemindahan pembulatan ke luar loop **membuat uji gagal** — kalau tidak, ujinya tidak menguji apa pun
- [ ] Validasi sisa direproduksi tanpa toleransi: kesamaan **persis** setelah pembulatan, tanpa distribusi galat
