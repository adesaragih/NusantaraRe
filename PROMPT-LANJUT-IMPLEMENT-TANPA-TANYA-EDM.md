# Prompt — lanjut implement tanpa banyak tanya · EDM (sesi `nusantarare-55`)

> Perintah work owner, 01-10-2026. Berlaku sampai work owner mencabutnya.

## Aturan kerja

1. **Langsung implement.** Jangan bertanya untuk hal yang bisa diputuskan dengan rekomendasi terbaik.
2. **Ambiguitas → putuskan sendiri, catat, lanjut.** Pilih opsi yang akan Anda rekomendasikan, catat di
   `docs/LAPORAN-IMPLEMENTASI-BEFORE-IMAGE.md` (atau register keputusan modul) sebagai *"Axx — keputusan
   agent, menunggu konfirmasi"* beserta bukti `path + rule`, lalu lanjut.
3. **`AskUserQuestion` hanya boleh untuk:** (a) operasi destruktif / menghapus berkas; (b) menjalankan
   migrasi ke database Oracle nyata; (c) mengubah modul di luar jatah Anda; (d) authn/authz. Selain itu
   **dilarang bertanya**. Menulis berkas migrasi di rentang modul (260-299) boleh; menjalankannya ke
   Oracle tidak.
4. **Korpus diam → jangan tebak, jangan berhenti.** Angka uang, arti kode, isi SP, `DecisionTable` tanpa
   baris: tulis `belum terverifikasi`, beri galat/kosong eksplisit di kode, catat sebagai pertanyaan
   terbuka, lanjut ke tiket berikutnya.
5. **Tiket kurang/tidak jelas → baca XML** di `D:\migrasi\RNM\Endorsment Fac In\`. Bila perlu tiket baru,
   tulis ringkas dulu di `docs/issues/` bertanda *"disusun agent dari XML atas perintah work owner — bukan
   hasil /to-tickets"*, lalu implement. Ini perintah eksplisit work owner, bukan meniru skill diam-diam.
   Bila korpus berbeda dari tiket: ikuti korpus, catat bedanya (seperti yang sudah Anda lakukan).
6. **CLAUDE.md §4, §4a, §7 tetap berlaku penuh:** bukti `path + rule`, label `[terverifikasi]` /
   `[dugaan]` / `[pertanyaan terbuka]`, sensus dua cara, uang tanpa `float`.
7. **Git: jangan commit, push, atau PR.** Kerja lokal saja. Commit `9ca6f06` dibiarkan apa adanya.
8. **Uji:** `go test` per paket selama kerja, full suite di akhir. Code review dua sumbu per kelompok
   tiket; perbaiki temuannya tanpa bertanya.
9. **Satu laporan di akhir**, bukan per tiket. Semua pertanyaan dikumpulkan di bagian penutup
   *"Menunggu work owner"*, tidak diajukan di tengah kerja.
10. **Perbarui baris `Status:`** setiap tiket — termasuk yang `blocked` padahal penghalangnya sudah
    selesai.

## Jatah: hanya `modul/endorsmentfacin`

⛔ Jangan menyentuh `modul/nbfacin`, `modul/rnwfacin`, atau kontrak mesin NB di `inti/backend/kontrak` —
itu jatah sesi `nusantarare-c3`. ⛔ Jangan membuat salinan/tiruan mesin NB di dalam EDM, dan jangan
mengimpor `modul/nbfacin` langsung.

## Penghalang yang sudah gugur

Tiket NB 01, 03, 04, 05, 08 sudah dikerjakan; 06 (layering) ditutup — tidak dipakai di sistem baru.
NB-15/16 sudah dibangun (menunggu tinjauan work owner).

## Urutan

1. **Tanpa menunggu kontrak NB:** tuntaskan bagian `sebagian` E06, E08, E09, E10 sejauh korpus
   mengizinkan; lalu E18 (jalur Life) → E20 (skoring medis).
2. **Begitu kontrak NB ada** — sesi `nusantarare-c3` akan mengirim pesan, atau periksa sendiri
   `inti/backend/kontrak/`: E01 → E02 → E16; E03 → E04 → E05; E12 → E13 → E14 → E15; E19; E22.
3. Bila semua tiket yang tidak butuh kontrak habis dan kontrak belum ada: tulis laporan, berhenti.
   Jangan membangun kontraknya sendiri.

## Tetap tertahan — lewati, sebut di laporan

E17 (seam repository belum ada), E21 (tabel flat P-10), pembulatan `@Math.divide` pada porsi periode
(tetap galat eksplisit), bentuk seam 4 dua fungsi (pertahankan, catat sebagai keputusan agent).
