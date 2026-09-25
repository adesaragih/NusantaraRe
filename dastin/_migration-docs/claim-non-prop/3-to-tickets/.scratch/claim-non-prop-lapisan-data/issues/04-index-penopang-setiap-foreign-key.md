---
status: selesai
---

# 04: Index penopang setiap foreign key

> **SELESAI 2026-09-18.** IX_ADJUSTMENT_1 terpasang; pembangkit kini memberi index penopang untuk setiap FK.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.
>
> **SELESAI PENUH 19 September 2026.** Pilahan *"aturannya selesai, namanya draf"* dicabut.
>
> Nama index tidak lagi menunggu apa pun: daftar singkatan dibatalkan dan REQ-032 turun jadi verifikasi. Index kini dinamai `IX_<tabel>_<n>` sesuai aturan 2 di tiket `03`; `IX_ADJ_REKENING` menjadi **`IX_ADJUSTMENT_1`**, penopang `FK_ADJUSTMENT_2`.
>
> Yang **tidak** berubah: kewajibannya. Setiap FK tetap harus punya index atau unique yang kolom pertamanya sama dengan kolom pertama FK itu, dan pembangkit tetap menegakkannya.

*Asal: `T-15` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada FK yang kolomnya bukan awalan sebuah index atau unique.

`IX_ADJUSTMENT_1` untuk `FK_ADJUSTMENT_2`; sapuan ulang seluruh FK oleh pembangkit.

**Blocked by:**

- None (can start immediately)


**Dasar:** DERIVED: di Oracle, FK tanpa index membuat `DELETE` dan `UPDATE` pada tabel induk mengambil kunci di tingkat tabel.

- [x] Setiap FK punya index atau unique yang kolom pertamanya sama dengan kolom pertama FK itu.
- [x] Empat belas FK ke `KLAIM` tetap tertopang `UNIQUE` berawalan `ID_KLAIM` — **tidak boleh hilang** saat kunci disentuh lagi.

**Ketidakpastian:** Tidak ada.
