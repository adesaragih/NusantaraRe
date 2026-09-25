---
status: tertahan
---

# 06: Versi usulan naik saat penanda berubah, dan penjaga menolak versi yang berselisih

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemegang jenjang berderajat terendah menyunting penanda usulan sambil
memutus, dan suntingan itu terekam sebagai **dua peristiwa dalam satu transaksi** — usulan
diubah, membawa nilai sebelum dan sesudah, dan jenjang memutus. Jenjang berikutnya yakin
bahwa yang ia putuskan sama persis dengan yang dilihat jenjang sebelumnya.

`VERSI_USULAN` dan `PERISTIWA_SIRKULASI` beserta constraint-nya; penjaga versi di batas tulis.

**Persyaratan:** `S-024`, `S-025`.

**Tidak termasuk:** `S-052` — peristiwa pergantian pemegang pada roster. Ia menunggu objek
`PERISTIWA_ROSTER`; `PERISTIWA_SIRKULASI` tidak dapat menampungnya.

**Jalur gagal:** penyuntingan datang dari jenjang bukan yang pertama → `422`, penanda usulan
**dan** versi **dan** keputusan sama-sama tidak tersimpan · versi usulan berubah oleh jenjang
selain yang pertama → `409` menyebut versi yang diharapkan dan versi yang ditemukan.

**Uji:** BARU — penyuntingan oleh jenjang bukan pertama ditolak · BARU — dua peristiwa lahir
dalam satu transaksi, dan keduanya hilang bersama ketika transaksi batal.

**Menggantikan:** suntingan yang menimpa nilai sebelumnya **tanpa jejak** — tidak ada padanan
di sistem lama; bentuk yang dijamin 296 kendali layar saja, kini ditegakkan juga sebagai
penjaga murah di batas tulis (`J-1`).

**Blocked by:**

- `04` — Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak


**Dasar:** DECIDED(`E-1`, `J-1`, `J-2`, `E-3`). EVIDENCED: `F-2`, `F-8`;
`3-to-tickets/FAKTA-LAYAR-01.md` — keempat medan berkunci jenjang pertama bernama.

- [ ] Versi usulan naik **hanya** ketika penanda usulan berubah.
- [ ] Tiap keputusan mencatat versi yang diputusnya.
- [ ] Penyuntingan terekam dua peristiwa dalam **satu** transaksi; keduanya batal bersama.
- [ ] Jenjang pertama yang **menolak** tetap meninggalkan suntingannya tercatat sebagai
      peristiwa, tanpa akibat hilir.
- [ ] Nilai sebelum dan sesudah menyimpan **penanda**, bukan uang — nol kolom uang di kedua
      objek.

**Ketidakpastian:** Tidak ada.
