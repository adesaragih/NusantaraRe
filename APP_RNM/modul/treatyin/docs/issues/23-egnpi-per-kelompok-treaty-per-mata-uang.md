---
status: aktif
golongan: pelestarian
---

# 23: EGNPI dicatat per kelompok treaty per mata uang

*Asal: `DAFTAR-PEKERJAAN.md` `P-08` · `SPEC-MODEL-DATA.md` §10.12 · `SPEC-INVARIAN.md` `INV-09`.*

**What to build:** **PK** mencatat **perkiraan premi bruto bersih tahunan** — EGNPI — per kelompok
treaty per mata uang. Baris kembar pada sumbu itu ditolak.

Artefak: entitas `EGNPI`, kunci asingnya, kolom `KODE_MATA_UANG`, dan `INV-09`.

**PEMBUAT PERTAMA** untuk `EGNPI`.

**Kenapa begini:** EGNPI adalah **penyebut hampir setiap perhitungan premi** di modul ini —
`DetailCalculation` menggelung `.EgnpiTotalList` dan mengalikannya dengan tarif penyesuaian untuk
menghasilkan premi diperoleh. Besaran yang menjadi penyebut orang lain harus punya **kunci yang tidak
dapat bertabrakan**, dan kuncinya menyebut mata uang dengan alasan yang sama seperti `RETENSI_CEDANT`.

**Persyaratan:** `INV-09` (kelompok treaty + mata uang unik di dalam satu versi) · `INV-36` · `INV-44` ·
`INV-39` dan `INV-40` **tidak dipasang di sini** — lihat `Tidak termasuk`

**Tidak termasuk:** **Tingkat pencatatan `NILAI_EGNPI`.** `INV-39` dan `INV-40` menuntut setiap paket
uang menyatakan apakah ia `TREATY_100_PERSEN` atau `BAGIAN_NURE`; untuk EGNPI **tingkatnya belum
ditentukan**, dan yang menjawabnya **`T-4`** dari teknik treaty. Constraint yang dipasang sekarang
**menuntut nilai yang migrasi tidak tahu cara mengisinya** — itu sebab ketiadaannya berdiri sebagai
**pernyataan keputusan** di `ddl-usulan/`, bukan sebagai pekerjaan tertunda.
**Premi diperoleh yang diturunkan dari EGNPI** — `F-15`, menunggu pemilik proses.

**Jalur gagal:** Dua baris berkelompok treaty **dan** mata uang sama -> **ditolak** `INV-09` · Nilai
EGNPI terisi tanpa mata uang -> ditolak `INV-36`.

**Uji:** **Negatif:** baris kembar pada sumbu penuh; nilai tanpa mata uang.
**Positif:** satu versi dengan dua baris berkelompok treaty sama tetapi **mata uang berbeda** ->
diterima.
**Positif kedua — dan ia khas entitas ini:** sebuah kontrak **tanpa satu pun baris EGNPI** tetap dapat
disimpan dan disetujui. EGNPI adalah perkiraan, bukan syarat; constraint yang menuntutnya ada akan
menolak kontrak yang sah pada hari pertama.

**Menggantikan:** `P-08` melestarikan daftar `EGNPI` sistem lama. Dua kolom yang **tidak dibawa**:
`CurrencyEgnpiAmount` dan `TotalEgnpiAmount` — keduanya agregat turunan yang tidak disimpan
(`SPEC-MODEL-DATA.md` §10.0c, `ADR-0037`).

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(TreatyIn.EGNPI@ekspor-2026-09; DetailCalculation@ekspor-2026-09 - gelung .EgnpiTotalList)
        DECIDED(INV-09, ADR-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-4)
```

- [ ] `EGNPI` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md` dengan `KODE_MATA_UANG` pada barisnya
- [ ] `INV-09` terpasang sebagai `UNIQUE` bertiga kolom
- [ ] kedua uji positif lulus — mata uang berbeda diterima, **dan** kontrak tanpa EGNPI diterima
- [ ] ketiadaan `INV-39`/`INV-40` tertulis sebagai **pernyataan keputusan** yang menyebut `T-4` sebagai penagihnya
- [ ] `TotalEgnpiAmount` dan `CurrencyEgnpiAmount` **tidak ada** sebagai kolom — diperiksa dengan sapuan
- [ ] `KTV-A` dan `T-4` tercatat di `ASUMSI-CLEAR.md`
