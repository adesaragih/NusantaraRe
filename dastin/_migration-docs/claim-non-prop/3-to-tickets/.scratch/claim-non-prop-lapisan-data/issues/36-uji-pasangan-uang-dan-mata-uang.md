---
status: menunggu-instance
---

# 36: Uji: pasangan uang dan mata uang

> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-25` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa tidak ada nilai uang yang dapat hidup tanpa satuannya, di tabel mana pun.

**Blocked by:**

- ~~`16`~~ *(selesai)* — Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs
- ~~`28`~~ *(selesai)* — Pasangan IDR–kurs


**Dasar:** DECIDED(ADR-0007). **Ketidakpastian** T-35 dapat menambah kolom yang ikut diuji.

- [ ] Nilai uang tanpa mata uang **ditolak**; nilai IDR tanpa kurs **ditolak**; kurs tanpa IDR **diterima**. Diuji di **setiap** tabel bernilai uang, bukan satu.

**Ketidakpastian:** T-35 dapat menambah kolom yang ikut diuji.
