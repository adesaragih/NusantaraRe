---
status: aktif
---

# 48: Pengajuan menolak dua syarat dan memperingatkan enam, dan peringatannya tercatat

*Asal: `DAFTAR-PEKERJAAN.md` `P-29` dan `P-30` · `INV-27` · `K1-1`…`K1-8` · `CO-6`…`CO-8`.*

**What to build:** **PK** mengajukan versi `DRAFT` untuk persetujuan. `K1-1` dan `K1-2` **menolak**
pengajuan, dengan pesan yang menyebut apa yang kurang. `K1-3` sampai `K1-8` menghasilkan
**peringatan yang tercatat** — belum penolakan.

**Persyaratan:** `INV-27` · `K1-1`…`K1-8` · `CO-6`…`CO-8` · `ADR-0055` (`AJUKAN`)

**Tidak termasuk:** **Menaikkan `K1-3`…`K1-8` menjadi penolakan** — itu kemampuan **BARU** tersendiri, dan
`CARA MENYALAKANNYA` miliknya belum ditulis. Tiket ini hanya memasang mode peringatan.

**Jalur gagal:** Pengajuan tanpa syarat `K1-1` -> **ditolak**, pesannya menyebut apa yang kurang ·
Pengajuan yang melanggar `K1-5` -> **diterima**, dan peringatannya **tercatat dan dapat dibaca**.

**Uji:** **Negatif:** kedelapan syarat diuji satu per satu; dua menolak, enam meloloskan.
**Positif:** pengajuan yang memenuhi seluruh delapan **diterima tanpa satu peringatan pun** —
peringatan palsu sama merusaknya dengan penolakan palsu.

> **Mode peringatan tanpa pembaca adalah fitur yang dimatikan, ditambah biaya log.** Kriteria
> selesai menuntut peringatannya **dapat dibaca seseorang**, bukan hanya tertulis.

**Menggantikan:** **`TreatyInSubmitEDM`** — pada jalur addendum **kedelapan syarat `K1` MATI**: `CheckID`,
`CheckError`, dan keluar-bergalat seluruhnya dinonaktifkan. Addendum diajukan **tanpa satu pun
pemeriksaan**.

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(TreatyInSubmitEDM@ekspor-2026-09 - CheckID/CheckError/exit MATI)
        DECIDED(INV-27, ADR-0055)
```

- [ ] `K1-1` dan `K1-2` menolak, pesannya menyebut apa yang kurang
- [ ] `K1-3`…`K1-8` menghasilkan peringatan **tersimpan**
- [ ] peringatan **dapat dibaca seseorang** — pembacanya disebut namanya, bukan "tersedia di log"
- [ ] uji positif lulus: pengajuan yang memenuhi delapan syarat **nol peringatan**
