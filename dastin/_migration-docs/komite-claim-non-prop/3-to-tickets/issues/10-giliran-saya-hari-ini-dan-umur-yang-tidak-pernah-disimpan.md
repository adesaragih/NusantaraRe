---
status: tertahan
---

# 10: Giliran saya hari ini, dan umur yang tidak pernah disimpan

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemegang jenjang membuka daftar sirkulasi yang menunggu keputusannya, dan
melihat sudah berapa lama tiap sirkulasi berjalan serta berapa lama jenjang aktif menunggu —
tanpa ada yang perlu menuliskannya, dan tanpa satu pun kolom umur tersimpan.

Jalur baca daftar menunggu pemanggil; umur sebagai atribut turunan.

**Persyaratan:** `S-045`, `S-048`.

**Tidak termasuk:** masa tenggat dan eskalasi. Tidak ada keadaan kedaluwarsa dan tidak ada
pekerja latar (`E-5`); umur dikumpulkan sebagai data agar aturan tenggat kelak diputuskan
dengan bukti, bukan dengan angka karangan.

**Jalur gagal:** pemanggil bukan pemegang jenjang mana pun → **daftar kosong, bukan galat** ·
pemanggil tidak berhak membaca klaimnya → `403`, nol sebagian isi bocor lewat pesan galat.

**Uji:** P5-03 — empat nama keranjang sistem lama tidak terpakai; daftar dihitung dari jenjang
aktif · BARU — uji DDL: nol kolom umur tersimpan.

**Menggantikan:** keranjang kerja Pega — daftar dihitung dari jenjang aktif, bukan dari
keranjang yang ditulis terpisah · `KomiteTreaty_Flow` — nol SLA, timer, eskalasi, reassign,
dan delegasi di seluruh modul (`F-4`); **tidak ada yang menggantikannya**, umur menjadi
turunan.

**Blocked by:**

- `04` — Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak


**Dasar:** DECIDED(`E-5`, `K5-1`, `K6-2`, `J-4`). EVIDENCED: `F-4`.

- [ ] Daftar dihitung dari **jenjang aktif**, bukan dari keranjang tersimpan.
- [ ] Umur sirkulasi dan umur tunggu jenjang aktif dihitung **saat dibaca**.
- [ ] Nol kolom umur tersimpan — dibuktikan pencarian atas katalog skema.
- [ ] Nol keadaan kedaluwarsa dan nol pekerja latar.

**Ketidakpastian:** Tidak ada.
