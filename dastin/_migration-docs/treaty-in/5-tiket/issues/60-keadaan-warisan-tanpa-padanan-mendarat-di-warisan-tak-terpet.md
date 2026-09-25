---
status: aktif
---

# 60: Keadaan warisan yang tidak punya padanan mendarat di warisan-tak-terpetakan, dengan nilai aslinya tersimpan

*Asal: `DAFTAR-PEKERJAAN.md` `P-51` · `ADR-0054` · `INV-21`, `INV-26`.*

**What to build:** Nilai keadaan lama yang tidak punya padanan — seperti `"test"` — mendarat di
**`WARISAN_TAK_TERPETAKAN`**, dan **nilai aslinya tersimpan** di `KEADAAN_WARISAN_ASLI`.

**Persyaratan:** `ADR-0054` · `INV-21` · `INV-26`

**Tidak termasuk:** **Perbaikannya** — tiket `61`. Tiket ini hanya menempatkan; memperbaiki adalah perbuatan orang.

**Jalur gagal:** Nilai liar dipetakan diam-diam ke keadaan sah -> **cacat**; itu membersihkan sejarah ·
`KEADAAN_WARISAN_ASLI` terisi pada baris yang **bukan** `WARISAN_TAK_TERPETAKAN` -> ditolak.

**Uji:** **Negatif:** nilai `"test"` dipetakan ke `DRAFT` -> harus tidak terjadi.
**Positif:** baris ber-`WARISAN_TAK_TERPETAKAN` **tidak dapat** memasuki alur kerja biasa — tidak
ada perpindahan keluar selain `PERBAIKAN_WARISAN`.

**Menggantikan:** **`TreatyInSetValue` menyetel `StatusAkseptasi = "test"`**, dan aturan itu terpasang di
**enam layar**. Nilai itu tidak punya padanan di antara keadaan sah mana pun.

**Blocked by:** `59`

**Dasar:**
```
EVIDENCED(TreatyInSetValue@ekspor-2026-09 - StatusAkseptasi="test" di enam layar)
        DECIDED(ADR-0054, INV-21, INV-26)
```

- [ ] nilai tanpa padanan mendarat di `WARISAN_TAK_TERPETAKAN`
- [ ] `KEADAAN_WARISAN_ASLI` terisi **hanya** di baris itu
- [ ] nilai aslinya **tidak pernah dibaca** perhitungan mana pun
- [ ] uji positif: tidak ada perpindahan keluar selain `PERBAIKAN_WARISAN`
