---
status: selesai
---

# 28: Pasangan IDR–kurs

> **SELESAI 2026-09-18.** RANCANGAN selesai — 23 CHECK pasangan sudah ada di DDL, dipasang pembangkit. Pemasangannya ikut tiket tabelnya.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-13` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada nilai IDR yang dapat lahir tanpa kurs yang menghasilkannya, di tabel mana pun.

**23 `CHECK`** pasangan, dipasang **oleh pembangkit** atas setiap kolom `*_IDR` di tabel yang punya `KURS`.

**Tidak termasuk:** Penularan keadaan ke baris turunan — itu aturan prosedur, diuji di T-28.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`22`~~ *(selesai)* — Koreksi bernilai tercatat menggantikan tambalan di dalam kode


**Dasar:** DECIDED(ADR-0007, ADR-0014). EVIDENCED: FINDING-006 — `RETURN 1` pada kurs yang tidak ditemukan.

- [ ] Di **setiap** tabel bernilai uang, `*_IDR` terisi tanpa `KURS` **ditolak**.
- [ ] Menambahkan tabel bernilai uang baru tanpa `CHECK`-nya **tidak mungkin** — pembangkit memasangnya sendiri.

**Ketidakpastian:** **T-35** dapat menambah kolom `*_IDR` di `ALOKASI_LAYER`; bila itu terjadi, `CHECK`-nya ikut terpasang tanpa perubahan tangan.
