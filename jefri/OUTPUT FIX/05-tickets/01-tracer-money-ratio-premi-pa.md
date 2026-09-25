# 01: Tracer — uang, rasio, dan satu jalur premi PA yang hidup ujung ke ujung

**What to build:** Satu kasus PA sederhana masuk, satu angka premi keluar — lewat tipe uang dan rasio
yang sesungguhnya, bukan angka telanjang. Ini irisan **tertipis yang lengkap**: begitu ia hijau,
seluruh rantai (parsing masukan → uang bermata-uang → rasio berskala → rumus → pembulatan → hasil)
terbukti hidup, dan tiket berikutnya tinggal menambah bentuk rumus.

Uang tidak pernah `float`. Rasio membawa satuannya sendiri. Keduanya **tipe berbeda yang tidak dapat
dijumlahkan** — satu-satunya jembatan adalah perkalian eksplisit `uang × rasio → uang`. Nilai masuk
sebagai teks berkoma desimal dan dikonversi **di batas input**, bukan tersebar di dalam perhitungan.

Bentuk rumus yang dipakai tracer ini adalah yang paling sederhana dari PA — pembagi 1.000 langsung,
tanpa pro-rata: `CalculatePremiPA_FacIn` **L1003**. Bentuk PA lainnya menyusul di tiket 04.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Nilai uang memakai tipe desimal dengan **mata uang wajib menyertai**; tidak ada `float` di jalur uang mana pun
- [ ] Rasio memakai tipe terpisah yang **membawa skalanya sendiri** (per mille / persen)
- [ ] Menjumlahkan uang dengan rasio **gagal saat kompilasi** — dibuktikan berkas uji yang wajib tidak terkompilasi, bukan unit test runtime
- [ ] Satu-satunya jembatan uang↔rasio adalah operasi perkalian eksplisit
- [ ] Parser masukan menerima **koma sebagai pemisah desimal** dan memangkas spasi di ujung; konvensi ditetapkan eksplisit, **tidak diserahkan ke locale** (K-027)
- [ ] Resolver lini bisnis → skala berisi **PA saja** pada tiket ini
- [ ] Pintu masuk perhitungan premi menerima satu kasus PA dan mengembalikan uang
- [ ] **Rekonsiliasi eksak**: fixture PA bentuk `L1003` cocok dengan sistem lama **sampai digit terakhir**, tanpa toleransi (ADR-0001)
- [ ] Pembulatan menuliskan presisinya di tempatnya, disertai komentar yang menyebut rule Pega asal dan nomor langkahnya
