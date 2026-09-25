---
status: aktif
---

# 51: Direktur menyetujui versi yang menunggunya, dan versinya menjadi disetujui bila seluruh K2 terpenuhi

*Asal: `DAFTAR-PEKERJAAN.md` `P-33` · `ADR-0055` (`SETUJUI` tingkat 3) · `K2-1`…`K2-6`.*

**What to build:** **DR** menyetujui versi di antriannya; versinya menjadi **`DISETUJUI`** — tetapi **hanya
bila seluruh `K2` terpenuhi**. Perpindahan ini membuat versinya **terminal**, dan karena itu ia
memicu pembekuan tiket `46`.

Pecahan **ketiga dari tiga**.

**Persyaratan:** `ADR-0055` (`SETUJUI` tingkat 3) · `K2-1`…`K2-6` · `INV-24` (versinya menjadi terminal)

**Tidak termasuk:** **Pembekuan nilainya** — tiket `46` yang memasangnya; tiket ini hanya memicu keadaannya.

**Jalur gagal:** DR menyetujui sementara satu syarat `K2` tidak terpenuhi -> **ditolak**, pesannya menyebut
syarat mana · Sesudah `DISETUJUI`, mengubah nilai versinya -> ditolak oleh `46`.

**Uji:** **Negatif:** keenam syarat `K2` diuji satu per satu.
**Positif:** DR menyetujui versi yang memenuhi seluruh `K2` -> `DISETUJUI`, **dan pembekuan `46`
langsung berlaku** — diuji dengan mencoba mengubah satu nilai sesudahnya.

**Menggantikan:** **`TreatyInForceResolveComplete`** — *persetujuan dapat terjadi tanpa penyetuju.* Tombol
itu menyetel selesai-disetujui **tanpa satu pun penyetuju**, dan `TreatyInForceEdit` yang membukanya
ber-`pyVisible = ALWAYS`. Dan `TreatyInSetToDirector` **seluruh langkahnya mati** — tingkat direktur
di sistem lama tidak pernah benar-benar dijalankan lewat jalur itu.

**Blocked by:** `50`

**Dasar:**
```
EVIDENCED(TreatyInForceResolveComplete@ekspor-2026-09 - menyetel selesai tanpa penyetuju)
        EVIDENCED(TreatyInSetToDirector@ekspor-2026-09 - langkah 1-4 blok //, 4.1-4.7 ikut mati)
        DECIDED(ADR-0055, INV-24)
```

- [ ] DR menyetujui -> `DISETUJUI`, hanya bila seluruh `K2` terpenuhi
- [ ] keenam `K2` diuji negatif satu per satu
- [ ] uji positif membuktikan **pembekuan `46` langsung berlaku** sesudahnya
- [ ] **tidak ada jalan** menyetel `DISETUJUI` selain lewat perpindahan ini
