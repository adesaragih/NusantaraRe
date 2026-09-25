---
status: aktif
---

# 46: Versi berkeadaan terminal beku, dan pembekuannya menjangkau seluruh entitas anaknya

*Asal: `DAFTAR-PEKERJAAN.md` `P-38` · `INV-24` · `ADR-0036`.*

**What to build:** Versi berkeadaan `DISETUJUI`, `DITOLAK`, atau `DIBATALKAN` **tidak dapat diubah nilainya
oleh siapa pun** — dan larangan itu menjangkau **seluruh entitas anaknya**, bukan hanya barisnya
sendiri.

Trigger `BEFORE INSERT OR UPDATE OR DELETE` pada tiap tabel anak.

**Persyaratan:** `INV-24` · `ADR-0036` (angka dasar persetujuan dibekukan saat disetujui)

**Tidak termasuk:** **Daftar perpindahan** — tiket `45`. Tiket ini menjaga **nilai**, bukan perpindahan.
**Pembekuan jenis dan materialitas sejak `AJUKAN`** — tiket `07`, dan ia **lebih awal** daripada
pembekuan ini.

**Jalur gagal:** `UPDATE` pada versi `DISETUJUI` -> ditolak · **`INSERT` baris layer baru** pada versi
`DISETUJUI` -> **ditolak** · `DELETE` baris penyebaran pada versi `DITOLAK` -> ditolak.

**Uji:** **Negatif:** `UPDATE`, `INSERT`, dan `DELETE` diuji **pada tiap tabel anak**, bukan hanya
pada `VERSI_KONTRAK`.
**Positif:** versi `DRAFT` dan versi yang sedang menunggu persetujuan **tetap dapat diubah** —
pembekuan tidak boleh bocor ke keadaan non-terminal.

> **`UPDATE`/`DELETE` saja masih mengizinkan versi disetujui BERTAMBAH baris.** Itu sebab
> triggernya harus `BEFORE INSERT OR UPDATE OR DELETE`, dan sebab uji negatifnya wajib menyentuh
> `INSERT`.

**Menggantikan:** `TDA-07` penggantinya — *dua kontrol layar mengosongkan status akseptasi kontrak yang
sudah disetujui.* Dan `Force Edit (dev)` ditambah `Save EDM(dev)`: operator divisi IT dapat membuka
kunci addendum yang **sudah disetujui**, menyuntingnya, lalu menyimpannya di tempat — tanpa mesin
selisih dihitung ulang.

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(TreatyInForceEdit@ekspor-2026-09 - pyVisible=ALWAYS di atas gerbang divisi IT)
        EVIDENCED(SaveTreatyIn_EDM_Act@ekspor-2026-09 - tidak menghitung ulang selisih)
        DECIDED(INV-24, ADR-0036)
```

- [ ] trigger `BEFORE INSERT OR UPDATE OR DELETE` terpasang pada **setiap** tabel anak versi
- [ ] uji negatif mencakup `INSERT`, bukan hanya `UPDATE`
- [ ] uji positif lulus: versi non-terminal tetap dapat diubah
- [ ] daftar tabel anak yang dijangkau **tertulis**, dan dicocokkan terhadap `KAMUS-KOLOM.md` — bukan disusun dari ingatan
