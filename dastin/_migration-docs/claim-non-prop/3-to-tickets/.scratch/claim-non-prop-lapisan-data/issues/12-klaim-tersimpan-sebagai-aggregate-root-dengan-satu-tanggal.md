---
status: selesai
menunggu-luar: [REQ-018]
---

# 12: Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


> **SELESAI 19 September 2026.** `ddl-usulan/01_KLAIM.sql` ditulis ulang lengkap.
>
> **Yang berubah dari versi sebelumnya — dua kolom baru, keduanya konsekuensi keputusan yang turun hari ini:**
>
> | Kolom | Dari | Kenapa di sini |
> |---|---|---|
> | `LINI_USAHA` | SPEC §21.3 (menutup G1 dan A15) | Penentu syariah lama ada dua — node yang mengeksekusi dan nama akun notifikasi — dan **tidak satu pun ikut pindah**. Apa pun jawaban G1, sistem baru perlu penanda yang berdiri sendiri. Disediakan sekarang meski belum ada yang mengisinya, sebagaimana ADR-0025 menyediakan kolom lingkup sebelum ada yang memerlukannya |
> | `PENUTUPAN_LAMA` | SPEC §21.6 (menutup B9 dan A7) | Jalur `CloseClaimMD` **tidak dimigrasi**, tetapi klaim yang pernah lewat sana dimigrasi sebagai tertutup **dan ditandai**. Tanpa kolom ini, keputusan "ditandai" tidak punya tempat, dan asal-usul baris hilang saat jalurnya dibuang |
>
> Keduanya **tanpa CHECK domain**, dan itu disengaja: nilainya belum terbaca, dan menutup domain yang belum diketahui adalah cacat yang sama dengan yang membuat `STATUS_KLAIM` dibiarkan terbuka (D1, D35, SPEC §21.1).
>
> **Nama constraint final**: `PK_KLAIM`, `UQ_KLAIM_1`, `CK_KLAIM_1`, `CK_KLAIM_2`, `CK_KLAIM_3`, `IX_KLAIM_1` — bentuk `<peran>_<tabel>[_n]`, tidak mengeja kolom. Terpanjang **10 byte**.
>
> **Empat aturan dinyatakan jatuh ke aplikasi**, bukan diam-diam tidak ditegakkan — ketetapan `TANGGAL_KEJADIAN`, domain `STATUS_KLAIM` dan `LINI_USAHA`, keterkaitan `PENUTUPAN_LAMA`, dan ketetapan urutan penulisan. Ditulis di kaki berkas DDL-nya.
>
> **REQ-018 tetap menunggui, tidak menahan.** `UQ_KLAIM_1` dipasang sekarang; bila data lama melanggarnya, yang berubah **penanganan migrasi**, bukan kuncinya.

*Asal: `T-03` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Klaim dapat disisipkan; nomor klaim ganda ditolak; masa berlaku treaty terbalik ditolak.

`KLAIM` beserta `PK_KLAIM`, `UQ_KLAIM_1`, `CK_KLAIM_1`, `CK_KLAIM_2`, `CK_KLAIM_3`, `IX_KLAIM_1` — `ddl-usulan/01_KLAIM.sql`.

**Tidak termasuk:** `CHECK` domain untuk `STATUS_KLAIM` — sengaja tidak ada, lihat T-12. Kolom jejak pelaku dipasang di sini tetapi tanpa foreign key; tabel pengguna tidak dirancang (ADR-0006).

**Blocked by:**

- **REQ-018** — **menunggui, tidak menahan** (dari luar papan): kuncinya sudah ditetapkan ADR-0024; yang berubah bila datanya melanggar adalah **penanganan migrasi**, bukan kunci
- ~~`01`~~ *(selesai)* — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis


**Dasar:** DECIDED(ADR-0001, ADR-0021, ADR-0022). EVIDENCED: `BLUEPRINT.md` §2.2 — sapuan 114 activity, **tidak satu pun `Property-Set` menulis `.DateOfLoss`**.

- [x] Klaim dengan nomor yang sudah ada **ditolak** — `UQ_KLAIM_1`.
- [x] Klaim ber-`TREATY_MULAI` sesudah `TREATY_AKHIR` **ditolak** — `CK_KLAIM_1`.
- [x] Klaim dengan `KEADAAN_BARIS` di luar tiga nilai domain **ditolak** — `CK_KLAIM_3`.
- [x] Klaim tanpa `DIBUAT_OLEH` **ditolak** (`NOT NULL`); `DIBUAT_ATAS_NAMA` kosong **diterima** — kosong berarti tidak ada perwakilan, bukan tidak diketahui.

**Ketidakpastian:** **REQ-018** (duplikat `PYID`) belum terjawab. `UQ_KLAIM_1` dipasang sekarang. Bila data lama melanggarnya, yang berubah adalah **penanganan migrasi**, bukan kuncinya — ADR-0024 sudah menetapkan pemilahan tiga kelompoknya.
