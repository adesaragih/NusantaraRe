---
status: menunggu-instance
---

# 38: Uji: bentuk view kompatibilitas

> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-30` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa hilir menerima bentuk yang sama dengan hari ini, dan menerima selisih pada muatan kedua.

**Blocked by:**

- ~~`29`~~ *(selesai)* — Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam
- ~~`05`~~ *(selesai)* — Sapuan S1: kolom yang sesungguhnya diterima Arasapas dan kasir


**Dasar:** EVIDENCED: `SPEC-MODEL-DATA.md` bagian 18. **Ketidakpastian** Pola selisih `InputParamOs.Value - OutOSAcc...` — lihat T-16.

- [ ] Kolom dan tipe yang dihasilkan sama dengan yang diterima hilir hari ini, termasuk `CASEID`. Dibandingkan terhadap 17 parameter `InputParamOs.*` hasil sapuan S1 — **dibaca, bukan dirancang**.
- [ ] **Muatan kedua atas klaim, jenis reasuransi, dan mata uang yang sama menghasilkan selisih, bukan nilai penuh.**
- [ ] Muatan yang ditandai ditolak **tidak** ikut mengurangi muatan berikutnya.

**Ketidakpastian:** Pola selisih `InputParamOs.Value - OutOSAcc...` — lihat T-16.
