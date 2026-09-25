---
status: aktif
---

# 47: Paling banyak satu versi tak-terminal per kontrak

*Asal: `DAFTAR-PEKERJAAN.md` `P-39` · `INV-25`.*

**What to build:** Sebuah kontrak hanya dapat punya **satu** versi yang belum selesai. Membuat versi kedua
yang belum selesai **ditolak**, dan pesannya menyebut versi mana yang masih terbuka.

**Persyaratan:** `INV-25` · `ADR-0055`

**Tidak termasuk:** **Pembatalan draf** yang membebaskan slotnya — tiket `56`.
**Penomoran versi** — sudah dipegang `INV-04` di tiket `14`.

**Jalur gagal:** Membuat versi penyesuaian sementara versi sebelumnya masih `DRAFT` atau menunggu persetujuan -> **ditolak**, pesannya menyebut nomor versi yang masih terbuka.

**Uji:** **Negatif:** buat versi kedua saat yang pertama `DRAFT`; saat menunggu SH; saat menunggu DR.
**Positif:** sesudah versi pertama menjadi `DISETUJUI`, `DITOLAK`, **atau `DIBATALKAN`**, versi
berikutnya **diterima** — ketiga keadaan terminal membebaskan slotnya, dan menguji hanya
`DISETUJUI` akan melewatkan dua.

**Menggantikan:** `TDA-01` — *penjaga duplikat `TreatyInEdmCheckDuplicate` **mati**; tabrakan masuk cabang
`UPDATE`, menimpa addendum lain, dan melapor berhasil.*

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(TreatyInEdmCheckDuplicate@ekspor-2026-09 - MATI)
        DECIDED(INV-25, ADR-0055)
```

- [ ] constraint menolak versi tak-terminal kedua, pesannya menyebut versi yang terbuka
- [ ] uji positif mencakup **ketiga** keadaan terminal, bukan hanya `DISETUJUI`
