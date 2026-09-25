---
status: menunggu-instance
---

# 34: Uji: paritas

> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-31` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa setiap baris hasil migrasi dapat ditunjuk balik ke barisnya di sistem lama, tepat satu.

**Blocked by:**

- ~~`23`~~ *(selesai)* — Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama
- ~~`15`~~ *(selesai)* — Jembatan ke sistem lama berdiri sebagai tabel terpisah


**Dasar:** DECIDED(ADR-0005), lewat ADR-0004. **Ketidakpastian** **REQ-018**; **REQ-015** mengukur status buka/tutup 8 klaim bertambalan.

- [ ] Untuk satu klaim contoh, tiap baris baru punya **tepat satu** pasangan di `MIGRASI_KORELASI`. Kasus uji utama **`CLMNP-975`** — tambalan terbaru, 2026-07-16, EVIDENCED `BLUEPRINT.md` §7.2.

**Ketidakpastian:** **REQ-018**; **REQ-015** mengukur status buka/tutup 8 klaim bertambalan.
