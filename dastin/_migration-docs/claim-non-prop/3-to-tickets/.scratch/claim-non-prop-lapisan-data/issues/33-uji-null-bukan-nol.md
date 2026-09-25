---
status: menunggu-instance
---

# 33: Uji: NULL bukan nol

> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-27` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa nilai yang belum dihitung tidak pernah terbaca sebagai nol.

**Blocked by:**

- ~~`16`~~ *(selesai)* — Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs


**Dasar:** DECIDED(ADR-0019). **Ketidakpastian** Tidak ada.

- [ ] Baris menunggu kurs ber-IDR `NULL`, `KEADAAN_BARIS` bukan `LENGKAP`, dan agregasi atasnya **tidak** memperlakukannya sebagai nol. Baris bernilai nol yang sudah dihitung dapat dibedakan darinya lewat satu kueri.

**Ketidakpastian:** Tidak ada.
