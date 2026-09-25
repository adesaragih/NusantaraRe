---
status: selesai
---

# 22: Koreksi bernilai tercatat menggantikan tambalan di dalam kode


> **SELESAI 19 September 2026.**
>
> **Butir 4 tidak terpenuhi oleh index yang ada.** Ia menjanjikan *"riwayat koreksi atas **satu kolom tertentu** terbaca lewat satu kueri"*, sementara `IX_KOREKSI_NILAI_1` hanya mencakup `(TABEL_SASARAN, ID_SASARAN)` — tanpa `KOLOM_SASARAN`. Janji itu berjalan lewat pemindaian seluruh riwayat baris. Index diperluas menjadi `(TABEL_SASARAN, ID_SASARAN, KOLOM_SASARAN, DIBUAT_PADA)`; `DIBUAT_PADA` ikut karena riwayat dibaca berurut waktu.
>
> **Satu constraint bertambah**: `CK_KOREKSI_NILAI_1` menolak baris yang `NILAI_SEBELUM` **dan** `NILAI_SESUDAH`-nya kosong. Keduanya boleh kosong sendiri-sendiri — `NILAI_SEBELUM` kosong berarti kolomnya belum pernah terisi, `NILAI_SESUDAH` kosong berarti koreksi ini **membatalkan** nilai (butir 2) — tetapi kosong keduanya **tidak menyatakan apa pun**. Sebuah koreksi yang tidak mengubah apa pun bukan koreksi.
>
> Kriteria lain terpenuhi bentuk yang ada: `ALASAN NOT NULL`, `DIBUAT_OLEH NOT NULL`, dan `UQ_KLAIM_PENJAGA_TANGGAL_1` atas `ID_KLAIM` — satu klaim muncul sekali.

*Asal: `T-11` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Setiap penyimpangan nilai punya nilai sebelum, nilai sesudah, alasan, pelaku, dan waktu.

`KOREKSI_NILAI`, `KLAIM_PENJAGA_TANGGAL`, `IX_KOREKSI_NILAI_1`.

**Tidak termasuk:** Pengisian daftar klaim terdampak — menunggu migrasi.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-0004, ADR-0018). EVIDENCED: `BLUEPRINT.md` §7 — **29 langkah, 18 ekspresi unik, 8 rule** menambal per-case; tambalan terbaru `CLMNP-975` bertanggal **2026-07-16**.

- [x] Koreksi tanpa alasan **ditolak**; tanpa pelaku **ditolak** — keduanya `NOT NULL`.
- [x] Koreksi dengan `NILAI_SESUDAH` kosong **diterima** — koreksi yang membatalkan nilai.
- [x] Koreksi dengan **kedua** nilai kosong **ditolak** — `CK_KOREKSI_NILAI_1`, baru.
- [x] Satu klaim hanya dapat muncul sekali di daftar klaim penjaga tanggal — `UQ_KLAIM_PENJAGA_TANGGAL_1`.
- [x] Riwayat koreksi atas satu kolom tertentu terbaca lewat satu kueri — `IX_KOREKSI_NILAI_1` **diperluas** dengan `KOLOM_SASARAN` dan `DIBUAT_PADA`.

**Ketidakpastian:** Tidak ada REQ. Bila `KOREKSI_NILAI` dan kolom suntingan di tingkat baris berbeda, **`KOREKSI_NILAI` yang berlaku** — `SPEC-MODEL-DATA.md` bagian 20.
