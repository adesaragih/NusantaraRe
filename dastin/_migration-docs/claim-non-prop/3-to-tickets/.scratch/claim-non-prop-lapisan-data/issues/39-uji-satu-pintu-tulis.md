---
status: menunggu-instance
---

# 39: Uji: satu pintu tulis

> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-29` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa satu pintu tulis berlaku sebagai hak akses, bukan sebagai kesepakatan.

**Blocked by:**

- `35` — Satu pintu tulis ditegakkan hak akses, bukan kesepakatan


**Dasar:** DECIDED(ADR-0017). **Ketidakpastian** **REQ-021**.

- [ ] Akun selain akun aplikasi **gagal** `INSERT` ke tabel kanonik dan **berhasil** `SELECT` lewat view.

**Ketidakpastian:** **REQ-021**.
