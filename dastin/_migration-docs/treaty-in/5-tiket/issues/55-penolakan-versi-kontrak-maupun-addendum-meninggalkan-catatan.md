---
status: aktif
---

# 55: Penolakan versi — kontrak maupun addendum — meninggalkan catatan, dan barisnya tidak dihapus

*Asal: `DAFTAR-PEKERJAAN.md` `P-35` dan `P-58` · `ADR-0055` (`TOLAK`).*

**What to build:** **SH / DH / DR** menolak versi. Penolakan tanpa alasan **ditolak**. Versinya menjadi
`DITOLAK` — **terminal** — dan **barisnya tetap tersimpan**, nomornya **tidak dipakai ulang**.

**Persyaratan:** `ADR-0055` (`TOLAK`, `DITOLAK` terminal) · `INV-04` (nomor tidak dipakai ulang) · `INV-24`

**Tidak termasuk:** **Pengembalian ke `DRAFT`** — tiket `53`. Menolak dan mengembalikan berakibat berbeda: yang satu terminal, yang lain tidak.

**Jalur gagal:** Penolakan tanpa alasan -> ditolak · Sesudah `DITOLAK`, mengubah versinya -> ditolak oleh
`46` · Nomor versi yang ditolak dipakai ulang -> ditolak oleh `INV-04`.

**Uji:** **Negatif:** alasan kosong; mencoba memakai ulang nomor versi yang ditolak.
**Positif:** versi yang ditolak **masih dapat dibaca** — beserta alasannya, pelakunya, dan
waktunya. Itu yang membuktikan barisnya tidak dihapus.

**Menggantikan:** `TDA-02` — ***menolak addendum MENGHAPUS barisnya.*** Dan kode pencatatannya **ada dan
dimatikan**: `TreatyInDeclineConfirmation_postactEDM` langkah 1–4 mati, dua langkah `RDB remove`
hidup. Akibatnya nomor revisi mungkin dipakai ulang, dan **tidak ada yang pernah dapat menjawab
berapa kali sebuah kontrak gagal diubah**.

**Blocked by:** `49` · `54`

**Dasar:**
```
EVIDENCED(TreatyInDeclineConfirmation_postactEDM@ekspor-2026-09 - RDB remove hidup, pencatatan MATI)
        DECIDED(ADR-0055, INV-04, INV-24)
```

- [ ] penolakan tanpa alasan ditolak
- [ ] versi `DITOLAK` **tetap tersimpan** dan dapat dibaca beserta alasannya
- [ ] nomornya **tidak dapat** dipakai ulang
- [ ] berlaku sama untuk versi kontrak **dan** versi addendum — diuji keduanya
