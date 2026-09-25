---
status: aktif
---

# 62: Kontrak dicari menurut cedant, asal bisnis, periode, sifat proporsi, dan keadaan siklus hidupnya

*Asal: `DAFTAR-PEKERJAAN.md` `P-56` · `SPEC-MODEL-DATA.md` §10.1 dan §10.2 · `ADR-0055`.*

**What to build:** **PK** mencari kontrak yang ada menurut lima sumbu — termasuk **keadaan siklus hidupnya** —
dan melihat hasilnya sebagai daftar.

**Persyaratan:** `SPEC-MODEL-DATA.md` §10.1, §10.2 · `ADR-0055` (himpunan keadaan yang dapat dicari)

**Tidak termasuk:** **Bentuk layarnya** — `L-4`: tidak ada spesifikasi layar di mana pun, dan ia milik batch
lapisan aplikasi. Yang dibangun di sini **perilaku pencariannya**, bukan tampilannya.

**Jalur gagal:** Mencari menurut keadaan yang tidak ada di daftar -> mengembalikan kosong, bukan galat · Kontrak tanpa versi berlaku -> tetap ditemukan, dengan keadaannya apa adanya.

**Uji:** **Negatif:** keadaan di luar delapan.
**Positif:** kontrak ber-`WARISAN_TAK_TERPETAKAN` **ditemukan** — bila ia tidak muncul di pencarian,
tidak akan ada yang tahu ia perlu diperbaiki.

**Menggantikan:** `TDA-11` — *daftar pilihan menyatukan kontrak dan addendum tanpa pembeda jenis dan
**tanpa saringan keadaan sama sekali***: `TreatyLoadMasterJoinEdm` meng-UNION `treaty_in` dengan
`treaty_in_edm` **tanpa `WHERE`**.

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(TreatyLoadMasterJoinEdm@ekspor-2026-09 - UNION tanpa WHERE, tanpa pembeda jenis)
        DECIDED(ADR-0055)
```

- [ ] kelima sumbu pencarian bekerja, termasuk keadaan
- [ ] uji positif: baris `WARISAN_TAK_TERPETAKAN` ikut ditemukan
- [ ] hasil pencarian **tidak** mencampur kontrak dengan versi tanpa pembeda
