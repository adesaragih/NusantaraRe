---
status: selesai
---

# 26: Premi pemulihan dapat dihitung ulang dari barisnya sendiri


> **SELESAI 19 September 2026 — tanpa perubahan struktur.** Ketiga kriteria sudah terpenuhi bentuk yang ada; yang bertambah **catatan di badan DDL-nya**, dan itu perlu.
>
> Sistem lama punya **dua** rumus premi pemulihan yang bekerja atas nilai **berbeda** (FINDING-007). Yang diwarisi `CountReinstatement_Act`; yang **tidak** diwarisi `AdjClaimCNP_Act` baris 3109 (AK-2b). Itu tidak boleh tinggal di dokumen saja: kolom tabel ini adalah masukan **rumus pertama**, dan bila kelak ada yang mengisinya dari rumus kedua, **angkanya akan tampak sah dan hasilnya berbeda tanpa satu pun tanda**. Sekarang tertulis di kepala berkasnya.
>
> `LIMIT_LAYER` adalah **penyebut**. Nol di penyebut bukan nilai yang aneh — ia perhitungan yang **tidak dapat dijalankan**. `CK_PREMI_PEMULIHAN_1` menolaknya di baris, bukan menyerahkannya ke kode yang harus ingat memeriksanya.

*Asal: `T-07` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Setiap masukan rumus tersimpan bersama hasilnya; limit layer nol ditolak.

`PREMI_PEMULIHAN` beserta seluruh masukan rumus; `CK_PREMI_PEMULIHAN_1`; `UQ_PREMI_PEMULIHAN_1`.

**Tidak termasuk:** Perhitungannya sendiri — ini lapisan data.

**Blocked by:**

- ~~`17`~~ *(selesai)* — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang


**Dasar:** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.3 — `CNPLimit` adalah **penyebut**. DECIDED-TEKNIS(AK-2b): acuannya `CountReinstatement_Act`; `AdjClaimCNP_Act` baris 3109 **tidak diwarisi** (FINDING-007).

- [x] Baris ber-`LIMIT_LAYER` nol **ditolak** — `CK_PREMI_PEMULIHAN_1`.
- [x] Baris kedua untuk (klaim, layer, mata uang) sama **ditolak** — `UQ_PREMI_PEMULIHAN_1`, yang menegakkan **keempat** field layer.
- [x] Seluruh masukan rumus dapat dibaca dari satu baris, tanpa menyentuh tabel master treaty — tujuh kolom, empat di antaranya `NOT NULL` karena rumus tanpa salah satunya tidak dapat dijalankan sama sekali.

**Ketidakpastian:** AK-2b menempatkan selisih terhadap sistem lama sebagai **pengecualian bernama** pada shadow-run, bukan kegagalan cutover. Besaran selisihnya belum terukur.
