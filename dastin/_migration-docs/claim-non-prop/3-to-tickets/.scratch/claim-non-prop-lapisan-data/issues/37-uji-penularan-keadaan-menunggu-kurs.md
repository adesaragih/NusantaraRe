---
status: menunggu-instance
---

# 37: Uji: penularan keadaan menunggu kurs

> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-28` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa ketiadaan kurs merambat ke seluruh turunannya, bukan berhenti di baris asalnya.

**Blocked by:**

- ~~`17`~~ *(selesai)* — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang
- ~~`18`~~ *(selesai)* — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran
- ~~`26`~~ *(selesai)* — Premi pemulihan dapat dihitung ulang dari barisnya sendiri


**Dasar:** DECIDED(ADR-0014). **Ketidakpastian** Penularannya aturan prosedur; basis data menyimpan keadaannya, tidak menegakkan penularannya. Uji ini karena itu menguji **hasil**, bukan mekanismenya.

- [ ] Baris turunan dari nilai yang menunggu kurs **ikut** bertanda menunggu kurs, di seluruh tabel turunannya — alokasi, Retensi Cedant, akseptasi, Adjustment, premi pemulihan.

**Ketidakpastian:** Penularannya aturan prosedur; basis data menyimpan keadaannya, tidak menegakkan penularannya. Uji ini karena itu menguji **hasil**, bukan mekanismenya.
