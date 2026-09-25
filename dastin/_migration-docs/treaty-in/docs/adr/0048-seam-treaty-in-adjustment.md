# ADR-0048 — Sambungan (seam) ke Treaty In Adjustment

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In
**Batas:** modul Treaty In Adjustment sendiri **berada di bawah embargo** dan tidak dibedah. ADR
ini hanya menetapkan apa yang harus **disediakan** Treaty In.

## Arahan pemilik proses yang mengikat

1. Saat sebuah Treaty In Adjustment dibuat, ID lama yang disimpan padanya adalah **ID Treaty In
   yang disesuaikan**. Sambungannya eksplisit, bukan disimpulkan.
2. **Data lama tidak boleh disimpan ulang.** Ia harus diambil lewat SELECT dari data lama. Tidak
   ada penyalinan nilai kontrak ke dalam catatan Adjustment.
3. **Selisih = nilai sekarang − nilai lama**, dan ia disimpan di **tabel tersendiri**.

Aturan (2) adalah ADR-0023 yang diterapkan oleh pemilik prosesnya sendiri: satu bentuk kanonik,
hilir diberi rujukan bukan salinan.

## Yang harus disediakan Treaty In, dan hanya ini

1. **Identitas versi kontrak** yang stabil, tidak berubah, dan bisa ditunjuk dari luar konteks.
2. **Jaminan** bahwa versi yang sudah disetujui tidak pernah berubah nilainya (ADR-0036, dibuktikan
   oleh ADR-0045).
3. **Bentuk baca** yang bisa dipakai konteks lain untuk mengambil nilai kontrak pada versi
   tertentu — lengkap dengan tingkat, mata uang, kurs, dan bagian yang menyertainya (ADR-0039).

Tabel selisih, layarnya, dan alur kerjanya **bukan urusan sesi ini dan tidak dirancang**.

## Sifat yang diinginkan, dicatat sekarang

Selisih adalah **turunan yang dibukukan**, jadi ia fakta tercatat — bentuk yang sama dengan potret
penyebaran. Karena itu baris selisih membawa **penunjuk ke versi kontrak mana** ia dihitung dan
**kapan**. Tanpa itu ia angka tanpa asal.

Dan karena nilai lama bisa diambil lewat SELECT sementara nilai sekarang diketahui, **selisih yang
tersimpan dapat direkonsiliasi ulang kapan saja**. Ketidakcocokan antara selisih tersimpan dan
selisih terhitung bukan masalah — **ia alat deteksi**. Rancang supaya rekonsiliasi itu mungkin.

## Kenapa pembekuan menjadi prasyarat, bukan kebersihan

Aturan (2) hanya aman bila data lama **tidak bergerak**. Selama nilai kontrak masih bisa berubah
sesudah akseptasi — lewat penyusunan ulang penyebaran, suntingan baris acuan, kurs tahun berjalan,
atau bagian NuRe yang berubah — maka selisih yang dihitung hari ini dan yang dihitung bulan depan
atas kontrak yang sama akan berbeda tanpa ada yang mengubah apa pun, dan selisih yang dibukukan
tidak bisa dipertanggungjawabkan.

ADR-0036 karena itu bukan pilihan rancangan yang bisa ditunda. Ia **prasyarat teknis** bagi cara
kerja Adjustment yang sudah ditetapkan.
