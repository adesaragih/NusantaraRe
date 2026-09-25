---
status: selesai
---

# 14: Nilai yang tidak dapat diurai tercatat, tidak dibulatkan, tidak dibuang


> **SELESAI 19 September 2026, dikerjakan bersama tiket `13`.**
>
> **Yang bertambah: kolom `ID_PENDARATAN`, dan itulah seam-nya.** Tanpa kolom itu, "tercatat **beserta asalnya**" (`SPEC-MODEL-DATA.md` §21.8) hanya berlaku sampai tingkat nama sumber — dan nama sumber tidak menunjuk baris mana.
>
> `FK_MIGRASI_NILAI_DITOLAK_1` menegakkan bahwa muatan yang ditunjuk memang ada; `IX_MIGRASI_NILAI_DITOLAK_1` menopangnya sesuai aturan tiket `04`.
>
> **Kolomnya boleh kosong, dan kosong punya arti**: sebagian nilai ditolak berasal langsung dari kolom tabel lama, bukan dari muatan. Kosong berarti **bukan dari muatan** — bukan **tidak diketahui** (ADR-0019). `SUMBER_LAMA` dan `PENGENAL_LAMA` tetap wajib, jadi asalnya selalu terbaca.
>
> **Arah relasinya dari sini ke sana**, karena satu muatan dapat melahirkan banyak nilai ditolak.
>
> `NILAI_MENTAH` tetap `VARCHAR2`, bukan `NUMBER`: mengubahnya jadi angka adalah persis hal yang gagal dilakukan.

*Asal: `T-22` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada angka yang hilang tanpa jejak.

`MIGRASI_NILAI_DITOLAK` beserta `PK_MIGRASI_NILAI_DITOLAK`, `FK_MIGRASI_NILAI_DITOLAK_1`, `IX_MIGRASI_NILAI_DITOLAK_1` — `ddl-usulan/20_MIGRASI_NILAI_DITOLAK.sql`.

**Blocked by:**

- ~~`01`~~ *(selesai)* — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis


**Dasar:** DECIDED(ADR-0014). DECIDED-TEKNIS(AK-4): migrasi berhenti bila ada **satu** baris yang tidak dapat diurai **dan** tidak dapat diselesaikan. Tidak ada baris yang dibuang karena "cuma sedikit".

- [x] Nilai mentah tersimpan **apa adanya**, tanpa pembulatan dan tanpa konversi.
- [x] Setiap baris membawa sebab penolakan dan pengenal barisnya di sistem lama.
- [x] Nilai yang berasal dari sebuah muatan menunjuk **baris muatannya**, bukan hanya nama sumbernya — `ID_PENDARATAN`.

**Ketidakpastian:** Tidak ada REQ. Volume belum terukur — **REQ-031**.
