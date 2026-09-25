---
status: menunggu-instance
---

# 32: Uji: kunci alami akseptasi

> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-24` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa baris ganda tidak dapat lahir.

**Blocked by:**

- ~~`18`~~ *(selesai)* — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran


**Dasar:** DECIDED(ADR-0024). Bentuk uji mengikuti `pengetahuan/PREFLIGHT.sql`: SQL berblok, tiap blok berdiri sendiri, dapat dijalankan dari TOAD. Itu satu-satunya prior art — `MEMORI_PEMAHAMAN.MD` §10.8 mencatat ketiadaan kerangka uji sebagai "aspek yang tidak ada".

- [ ] `INSERT` kedua dengan (klaim, layer, mata uang) sama **ditolak**; `INSERT` dengan `LayerPart` berbeda **diterima**, dan itu perilaku yang dimaksud — bukan celah.

**Ketidakpastian:** **REQ-033**.
