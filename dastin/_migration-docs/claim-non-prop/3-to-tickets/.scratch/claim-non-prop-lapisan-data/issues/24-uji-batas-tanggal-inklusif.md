---
status: menunggu-instance
---

# 24: Uji: batas tanggal inklusif

> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-26` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa klaim yang jatuh tepat di hari terakhir masa berlaku treaty tidak tertolak.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-0022). EVIDENCED: FINDING-005 — di sistem lama perbandingannya dilakukan antara dua format teks yang berbeda.

- [ ] Klaim ber-Tanggal Kejadian **tepat di hari terakhir** masa berlaku treaty **diterima**; sehari sesudahnya **ditolak**.

**Ketidakpastian:** **REQ-020** mengukur berapa klaim lama yang terdampak; tidak menahan uji ini.
