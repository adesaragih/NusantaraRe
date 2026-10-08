# 10: Pemuat migrasi endorsemen

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Seluruh dokumen polis dipindahkan, setiap generasinya.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `modul/nbtreatyin/docs/KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 5.

> ⛔ **KOREKSI 06-10-2026** (log `../KOREKSI-DOKUMEN-2026-10-06.md`).
> 1. Bab *Kenapa tiket ini `blocked`* di bawah dicoret — penahannya gugur 23-09.
> 2. Keadaan NB nyata: pemuat dokumen lama NB menolak generasi `PRODKE > 0` dengan
>    `ErrGenerasiEndorsemen` *"milik pemuat EDM tiket 10"* (`modul/nbtreatyin/backend/models/dokumenlama.go`
>    baris 113–115) ⇒ seluruh generasi endorsemen lama memang menunggu tiket ini. Sumbernya dokumen lama
>    `JSON_POLIS.DATA_JSON` berbaris `PRODKE ≥ 1`; sistem baru **tidak menulis** `DATA_JSON` (NB
>    `models/produksi.go` baris 63) — hasil muatnya baris tabel, tanpa JSON.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** **09** · ⛔ `[work owner]` **lingkup pemindahan dokumen lama**~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **22** *(pemuat dokumen lama)*
**Menutup:** AC **44** *(1 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-3

## Hasil & nilai pengguna

Dokumen endorsemen lama dimuat **lewat antarmuka yang sama** dengan jalur biasa, sehingga data lama
melewati pemeriksaan yang sama — termasuk aturan keutuhan nomor urut dan larangan percabangan.

⛔ **Nol pintu belakang.** Pemuat yang menulis langsung ke tabel akan melewati seluruh penjaga yang
dibangun tiket **02**, **03**, dan **04**.

## Yang dibangun

Pemuat massal endorsemen yang memakai **antarmuka penyimpanan yang sama**, memulihkan rantai
generasi dari dokumen lama secara berurutan, dan menyalakan kedua penanda migrasi.

⚠️ **Urutan pemuatan penting:** generasi harus dimuat menurut nomornya, sebab penunjuk generasi
menuntut generasi sebelumnya sudah ada.

## Batas — yang TIDAK termasuk

⛔ Kedua penanda migrasi itu sendiri — tiket **09**.
⛔ Pemuat dokumen polis baru — tiket NB **22**.
⛔ Keputusan lingkup pemindahan — `[work owner]`.

## Cara mengujinya

Lewat seam `repository` yang sama. ⭐ **Uji utama:** rantai tiga generasi dimuat dari dokumen lama,
lalu diperiksa bahwa penjaga percabangan dan penjaga keutuhan **benar-benar berjalan** atasnya —
bukan dilewati karena ini jalur migrasi.

## Acceptance criteria

- [ ] **AC 44** — pemuat migrasi menulis lewat antarmuka penyimpanan yang **sama**

## ⛔ ~~Kenapa tiket ini `blocked`~~ — ✅ penahan gugur 23-09 *(dicoret koreksi 06-10)*

~~`[work owner]` **Lingkup pemindahan belum diputuskan** — sama dengan penahan tiket **09**. Memuat
seluruh riwayat dan memuat dua tahun terakhir adalah pekerjaan yang berbeda besarnya.~~
⭐ `[keputusan work owner]` 23-09: **seluruh** polis, **setiap** generasi.

⚠️ **Tiket ini sengaja kecil** — satu AC — sebab isinya hampir seluruhnya dipakai bersama tiket NB
**22**. Yang khas endorsemen hanya **urutan pemuatan generasi**.