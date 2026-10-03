# 06: Perhitungan selisih di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** **02** · **03**
**Bergantung pada tiket NB:** **16** *(kerangka penyimpanan)*
**Menutup:** AC **16–22** *(7 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-4 · ID-28 · ID-29 · ID-30 · ID-31 · ID-32

## Hasil & nilai pengguna

Selisih antar generasi dihitung **di lapisan layanan**, dengan **satu rumus** untuk semua kasus —
endorsemen pertama maupun berlapis, proporsional maupun non-proporsional:

```
selisih.X = baris_ini.X − baris(generasi sebelumnya).X
```

⛔⛔ **Rumus ini TIDAK pernah ditulis di dalam basis data.** View yang menghitung selisih
**dibatalkan** — rumus di dua tempat cepat atau lambat bercabang.
⛔ **Tiket ini tidak boleh menghidupkannya kembali.**

## Yang dibangun

Perhitungan selisih dengan empat aturan turunan:

| Golongan medan | Perlakuan |
| --- | --- |
| **uang** | dikurangi |
| **persentase** | ⭐ **disalin**, tidak dikurangi |
| **kunci** | disalin apa adanya |
| **arah utang-piutang** | diturunkan dari **tanda** selisih |

⭐ Pada selisih lapisan, arah utang-piutang diturunkan dari **jumlah sepanjang daftar lapisan**,
bukan dari satu lapisan saja.

⭐ Lapisan penyimpanan hanya mengambil **dua baris** — baris ini dan baris yang ditunjuknya.
⛔ Nol agregasi di basis data.

⚠️ `[penyimpangan sadar]` **Sistem lama punya DUA rumus di sisi proporsional**, dan yang kedua
mengurangi terhadap **selisih** generasi lampau, bukan terhadap **nilainya**. Dengan angka contoh:
nilai baru 180, nilai lama 150, selisih lama 50 ⇒ varian kedua menghasilkan **130**, padahal yang
benar **30**. ⛔ **Varian kedua tidak ditiru** — `[keputusan work owner]`.

## Batas — yang TIDAK termasuk

⛔ Penulisan hasilnya ke tabel proyeksi — tiket **07**.
⛔ Penandaan baris yang lahir dari rumus lama — tiket **09**.

## Cara mengujinya

Operannya lewat seam `repository`; rumusnya lewat lapisan layanan.

⭐ **Uji regresi wajib:** sisi non-proporsional harus menghasilkan angka yang **sama** sebelum dan
sesudah penyeragaman — `[terverifikasi]` ia memang sudah seragam, dan penyeragaman tidak boleh
mengubahnya. Uji yang menemukan perubahan berarti rumusnya salah diterjemahkan.

## Acceptance criteria

- [ ] **AC 16** — satu rumus untuk **seluruh** generasi
- [ ] **AC 17** — medan **uang** berisi hasil pengurangan
- [ ] **AC 18** — medan **persentase** disalin, tidak dikurangi
- [ ] **AC 19** — medan **kunci** disalin apa adanya
- [ ] **AC 20** — arah utang-piutang diturunkan dari **tanda** selisih
- [ ] **AC 21** — pada selisih lapisan, arah itu diturunkan dari **jumlah** sepanjang daftar
- [ ] **AC 22** — sisi non-proporsional menghasilkan angka **sama** sebelum dan sesudah

## ⚠️ Satu hal yang tidak terbaca dari sistem lama

`[terbuka]` Kedua blok rumus di sistem lama **sama-sama bergerbang aktif**, dan pemilihnya hanya
tertulis di **keterangan langkah**, bukan di gerbang yang dijalankan. ⇒ Mana yang benar-benar
berjalan **tidak dapat dinyatakan dari ekspor**.

⭐ **Tidak menahan tiket ini** — varian kedua tidak ditiru apa pun jawabannya.