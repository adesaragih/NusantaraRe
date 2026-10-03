# 25: Perhitungan selisih di lapisan layanan

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\06`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~ready-for-agent~~ — **DIGANTIKAN**
**Blocked by:** **24**
**Menutup:** EDM AC **16–22** *(7 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-4 · ID-28..ID-32

## Hasil & nilai pengguna

Selisih antara generasi dihitung **di lapisan layanan**, dengan **satu rumus** untuk semua kasus —
endorsemen pertama maupun berlapis, proporsional maupun non-proporsional.

```
selisih.X = baris_ini.X − baris(generasi sebelumnya).X
```

⚠️ `[penyimpangan sadar]` **Sistem lama punya DUA rumus di sisi proporsional**, dan yang kedua
mengurangi terhadap **selisih** generasi lampau, bukan terhadap **nilainya**. Dengan angka contoh:
nilai baru 180, nilai lama 150, selisih lama 50 ⇒ varian kedua menghasilkan **130**, padahal yang
benar **30**. ⛔ **Varian kedua tidak ditiru** — `[keputusan work owner]`.

## Yang dibangun

Perhitungan selisih di lapisan layanan, dengan empat aturan turunan:

| Golongan medan | Perlakuan |
| --- | --- |
| **uang** | dikurangi |
| **persentase** | ⭐ **disalin**, tidak dikurangi |
| **kunci** | disalin |
| **arah utang-piutang** | diturunkan dari **tanda** selisih |

⭐ Pada selisih lapisan, arah utang-piutang diturunkan dari **jumlah sepanjang daftar lapisan**,
bukan dari satu lapisan.

⭐ Lapisan penyimpanan hanya mengambil **dua baris** — baris ini dan baris yang ditunjuknya.
⛔ Tidak ada agregasi di basis data.

## Batas — yang TIDAK termasuk

⛔ Penulisan hasilnya ke tabel proyeksi — tiket **26**.
⛔ Penanda migrasi untuk baris hasil rumus lama — tiket **28**.

## Cara mengujinya

Lewat seam `repository` untuk operannya, dan lewat lapisan layanan untuk rumusnya.
⭐ **Uji regresi wajib:** sisi non-proporsional harus menghasilkan angka yang **sama** sebelum dan
sesudah penyeragaman — ia memang sudah seragam, dan penyeragaman tidak boleh mengubahnya.

## Acceptance criteria

- [ ] **AC 16** — satu rumus untuk **seluruh** generasi
- [ ] **AC 17** — medan **uang** berisi hasil pengurangan
- [ ] **AC 18** — medan **persentase** disalin, tidak dikurangi
- [ ] **AC 19** — medan **kunci** disalin apa adanya
- [ ] **AC 20** — arah utang-piutang diturunkan dari **tanda** selisih
- [ ] **AC 21** — pada selisih lapisan, arah itu diturunkan dari **jumlah** sepanjang daftar
- [ ] **AC 22** — sisi non-proporsional menghasilkan angka yang **sama** sebelum dan sesudah

## ⚠️ Satu hal yang tidak terbaca dari sistem lama

`[terbuka]` Kedua blok rumus di sistem lama **sama-sama bergerbang aktif**, dan pemilihnya hanya
tertulis di **keterangan langkah**, bukan di gerbang yang dijalankan. ⇒ Mana yang benar-benar
berjalan **tidak dapat dinyatakan dari ekspor**.

⭐ **Tidak menahan tiket ini** — varian kedua tidak ditiru apa pun jawabannya.