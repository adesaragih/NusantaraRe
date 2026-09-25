---
status: selesai
menunggu-luar: [REQ-018]
---

# 15: Jembatan ke sistem lama berdiri sebagai tabel terpisah


> **SELESAI 19 September 2026.** `ddl-usulan/19_MIGRASI_KORELASI.sql` ditulis ulang lengkap.
>
> **Dua hal bertambah, dan keduanya lahir dari membaca ulang janji tiket ini.**
>
> **1. `CK_MIGRASI_KORELASI_1` — sekurangnya satu pengenal lama wajib ada.** Kelima kolom boleh kosong sendiri-sendiri (itu memang kriteria penerimaan kedua), tetapi **kosong semua** menghasilkan baris jembatan yang tidak menjembatani apa pun. Ini satu-satunya keutuhan yang dapat ditegakkan basis data di tabel ini — karena itu dipasang.
>
> **2. `IX_MIGRASI_KORELASI_2` dan `_3` — arah lama → baru.** Versi sebelumnya hanya punya index `(TABEL_TUJUAN, ID_TUJUAN)`, yaitu arah **baru → lama**. Judul tiket ini menjanjikan arah yang **berlawanan**: *"dapat ditunjuk balik lewat pengenal lamanya"*. Tanpa index atas `PZINSKEY_LAMA` dan `PYID_LAMA`, janji itu berjalan lewat pemindaian penuh. Ditambahkan.
>
> **Tiga hal dinyatakan jatuh ke aplikasi**, bukan diam-diam tidak ditegakkan: `UNIQUE` kombinasi pengenal lama (menunggu REQ-018), keutuhan rujukan `TABEL_TUJUAN`/`ID_TUJUAN` (polimorfik — Oracle tidak punya FK polimorfik, dan memalsukannya dengan 22 kolom nullable lebih buruk daripada tidak punya), dan umur baris.
>
> **Istilah `_Avoid_` di nama kolom dijelaskan di badannya**, bukan dibiarkan tanpa alasan: `CASEID` dan `pyID` dipakai sebagai **nama lama bersufiks `_LAMA`**, di tabel yang seluruh tugasnya memang menyimpan nama lama — pengecualian yang sama dengan lapisan view kompatibilitas (ADR-0021). Tidak satu pun muncul di tabel kanonik.
>
> **REQ-018 tetap menunggui, tidak menahan.** Yang berubah bila jawabannya lain: **constraint, bukan kolom** — dan itu sebabnya `UNIQUE`-nya sengaja belum dipasang. Menebaknya sekarang menggugurkan migrasi di tengah jalan, tepat ketika ia paling mahal dibatalkan.

*Asal: `T-23` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Setiap baris kanonik dapat ditunjuk balik lewat pengenal lamanya, tanpa mengurai teks.

`MIGRASI_KORELASI` dengan **lima kolom terpisah**, `PK_MIGRASI_KORELASI`, `CK_MIGRASI_KORELASI_1`, dan tiga index — `IX_MIGRASI_KORELASI_1` (baru → lama), `_2` dan `_3` (lama → baru). `ddl-usulan/19_MIGRASI_KORELASI.sql`.

**Tidak termasuk:** `UNIQUE` atas kombinasi pengenal lama — **ditahan, dengan alasan**: kolom mana yang ada sudah EVIDENCED; **kombinasi mana yang unik** menunggu REQ-018.

**Blocked by:**

- **REQ-018** — **menunggui, tidak menahan** (dari luar papan): `UNIQUE` atas kombinasi pengenal lama sudah **di luar lingkup tiket ini**; kolomnya EVIDENCED, kombinasinya menyusul
- ~~`01`~~ *(selesai)* — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis


**Dasar:** DECIDED(ADR-0005, ADR-0023). EVIDENCED: kelima pengenal lama — `PZINSKEY` dan `PYID` dari DDL tabel work, `CASEID` dari `OS_AKSEPTASI_KLAIM` dan §13.2, `IndexObject` dari §8.4.

- [x] Kelima pengenal terbaca sebagai kolom tersendiri dan dapat di-`JOIN` tanpa penguraian teks.
- [x] Baris yang bukan berasal dari work object Pega **diterima** dengan `PZINSKEY_LAMA` kosong — kosong berarti bukan dari sana, bukan tidak diketahui.
- [x] Baris tanpa **satu pun** pengenal lama **ditolak** — `CK_MIGRASI_KORELASI_1`.
- [x] Pencarian **dari pengenal lama** tertopang index, bukan pemindaian penuh — `IX_MIGRASI_KORELASI_2`, `_3`.

**Ketidakpastian:** **REQ-018**, BLOCKER, OPEN. Yang berubah bila jawabannya lain: **constraint**, bukan kolom.
