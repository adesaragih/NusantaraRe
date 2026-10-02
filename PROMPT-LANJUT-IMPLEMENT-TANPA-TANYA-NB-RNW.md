# Prompt — lanjut implement tanpa banyak tanya · NB + RNW (sesi `nusantarare-c3`)

> Perintah work owner, 01-10-2026. Berlaku sampai work owner mencabutnya.

## Aturan kerja

1. **Langsung implement.** Jangan bertanya untuk hal yang bisa diputuskan dengan rekomendasi terbaik.
2. **Ambiguitas → putuskan sendiri, catat, lanjut.** Pilih opsi yang akan Anda rekomendasikan, catat di
   register keputusan modul (`docs/KEPUTUSAN-30-09-2026.md`) sebagai *"Axx — keputusan agent, menunggu
   konfirmasi"* beserta bukti `path + rule`, lalu lanjut.
3. **`AskUserQuestion` hanya boleh untuk:** (a) operasi destruktif / menghapus berkas; (b) menjalankan
   migrasi ke database Oracle nyata; (c) mengubah modul di luar jatah Anda; (d) authn/authz. Selain itu
   **dilarang bertanya**. Menulis berkas migrasi di rentang modul boleh; menjalankannya ke Oracle tidak.
4. **Korpus diam → jangan tebak, jangan berhenti.** Angka uang, arti kode, isi SP, `DecisionTable` tanpa
   baris: tulis `belum terverifikasi`, beri galat/kosong eksplisit di kode, catat sebagai pertanyaan
   terbuka, lanjut ke tiket berikutnya.
5. **Tiket belum ada → baca XML.** Baca korpus di `D:\migrasi\RNM\<modul>\`, tulis tiket ringkas lebih
   dulu di `docs/issues/` bertanda *"disusun agent dari XML atas perintah work owner — bukan hasil
   /to-tickets"*, lalu implement. Ini perintah eksplisit work owner, bukan meniru skill diam-diam.
6. **CLAUDE.md §4, §4a, §7 tetap berlaku penuh:** bukti `path + rule`, label `[terverifikasi]` /
   `[dugaan]` / `[pertanyaan terbuka]`, sensus dua cara, uang tanpa `float`.
7. **Git: jangan commit, push, atau PR.** Kerja lokal saja.
8. **Uji:** `go test` per paket selama kerja, full suite di akhir. Code review dua sumbu per kelompok
   tiket; perbaiki temuannya tanpa bertanya.
9. **Satu laporan di akhir**, bukan per tiket. Semua pertanyaan dikumpulkan di bagian penutup
   *"Menunggu work owner"*, tidak diajukan di tengah kerja.
10. **Perbarui baris `Status:`** setiap tiket yang dikerjakan atau ditutup.

## Jatah: `modul/nbfacin`, `modul/rnwfacin`, kontrak mesin NB di `inti/backend/kontrak`

⛔ Jangan menyentuh `modul/endorsmentfacin` — itu jatah sesi `nusantarare-55`.

## Urutan

1. **Rapikan status NB.** Tiket 01, 03, 04, 05, 07, 08, 09, 10, 13 sudah dikerjakan tetapi masih
   `ready-for-agent` — perbarui (04: sebagian, PA metode 2 menunggu DBA). Indeks menulis "16 tiket",
   berkasnya 17 — betulkan.
2. **Kontrak mesin NB di `inti/backend/kontrak`** (`rules.Eval`, `premium.Calculate`,
   `acceptance.Next`, resolver COB), disambung lewat `Pendaftaran()` + perakit, lolos penjaga
   impor lintas modul. Selesaikan R01 sebatas layanan.
3. **Begitu kontrak lulus uji:** kirim pesan ke sesi `nusantarare-55` (SendMessage) berisi berkas
   kontrak, nama antarmuka, dan cara memakainya. Sesi itu menunggu kontrak ini untuk tiket EDM.
4. **Rumus premi lini lain dari XML** — FIRE, ANEKA, BONDING, GOLF, MARINE CARGO. Mulai dari
   `NB FacIn\Activity\`: `CountNetPremiFOFire.xml`, `fillPremiAneka.xml`, `fillPremiBond.xml`,
   `FillPremiGolf.xml`, `PremiPaymentMarine.xml`, `CountPremiAndTSIRNM*_ACT.xml`,
   `SumTSIPremiSpreadedRNM_*_Act.xml`. Tulis tiketnya dulu (aturan 5).
   **Patokan lulus: 5 kasus nyata P-5 di kerangka rekonsiliasi** berubah dari "belum tercakup" menjadi
   "cocok" sampai digit terakhir. Yang tetap selisih dilaporkan, bukan disesuaikan.
5. **RNW R02–R08**, sebatas layanan.

## Tetap tertahan — lewati, sebut di laporan

Tiket 04 (data DBA), 17 (pemuat repository), kasus MBU #5 (pemilik Pega), A31 (konfirmasi work owner).
