# ADR-0040 — Identitas kontrak, addendum, dan kunci alami yang memperingatkan

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Sistem lama merangkai kontrak dengan satu kolom `OLDID` yang dipakai untuk dua hubungan sekaligus:
salinan sebuah kontrak, dan addendum atas sebuah kontrak. Addendum tinggal di tabel terpisah
(`M_TREATY_IN_EDM`, `TREATY_IN_EDM`) dan muncul sebagai baris setara di daftar.

Tidak ada deteksi duplikat yang berjalan: seluruh logikanya di `CheckDuplicateOffer` dimatikan.

## Keputusan

1. **`OLDID` pecah menjadi dua hubungan yang dinamai berbeda.** "Salinan dari" dan "addendum atas"
   adalah dua hal berbeda dan tidak boleh berbagi kolom. "Addendum atas" berhenti menjadi kolom
   dan menjadi hubungan versi.
2. **Kunci alami** sebuah kontrak adalah cedant + source of business + periode + `ProportionType`.
   Ia **memperingatkan, tidak melarang** — ada keadaan sah di mana dua kontrak berbagi kunci, dan
   sistem tidak boleh memutuskan itu untuk penggunanya.
3. **Lapisan beku addendum.** Untuk setiap addendum, lima hal tidak boleh berubah: cedant, source
   of business, `ProportionType`, tanggal mulai, dan tanggal berakhir. Perubahan atas salah satunya
   berarti kontrak lain, bukan addendum.
4. **Perpanjangan di tengah periode bukan addendum.** Ia operasi tersendiri dengan namanya sendiri.
5. **Materialitas addendum diturunkan**, tidak diketik: bila perubahan itu menghasilkan baris
   selisih, ia material. Kecuali untuk baris warisan — lihat ADR-0042.
6. **Bagian NuRe boleh berubah** lewat addendum, dan perubahannya berlaku ke depan.

## Konsekuensi

- Deteksi duplikat adalah **kemampuan baru**, bukan pelestarian. Jangan dihitung gratis.
- Peringatan kunci alami harus bisa diabaikan oleh orang, dan pengabaiannya tercatat.
- Riwayat versi kontrak menjadi struktur pertama-kelas, bukan tabel terpisah yang dirangkai kolom.

## Penandaan 23 September 2026 — keputusan (1) adalah PERUBAHAN, bukan pelestarian

Dibuktikan saat merancang sambungan ke modul Adjustment, dan ditandai di sini supaya tidak
ditemukan ulang sebagai "spesifikasi yang keliru".

**Di sistem lama, addendum yang disetujui tidak menggantikan apa pun.** Aktivitas
`SaveTreatyIn_EDM_Act` memuat dua langkah yang akan mempromosikannya — satu menulis ke
`M_TREATY_IN`, satu memanggil penyimpan `TREATYINDETAIL` — dan **keduanya ber-blok `//`, mati**.
Yang hidup hanya penyimpan `TREATYINDETAILEDM`. Addendum karena itu tersimpan di tabelnya sendiri,
dan baris kontrak tetap seperti semula.

**Di sistem baru, versi yang disetujui menggantikan pendahulunya**, dan hilir membaca versi yang
berlaku. Tidak ada baris "kontrak saat ini" yang terpisah untuk ditulis, sehingga tidak ada langkah
penyalin yang dapat dimatikan.

Keputusan ini **tidak bersandar pada bukti sistem lama** — bukti sistem lama justru melawannya.
Ia keputusan rancangan, diambil sadar. Seberapa banyak data yang terdampak diukur **Uji Z**.
