---
status: aktif
golongan: pelestarian
---

# 37: Potongan atas premi dicatat, dengan aturan yang sama di kedua tempat ia melekat

*Asal: `DAFTAR-PEKERJAAN.md` `P-19` · `SPEC-MODEL-DATA.md` §14.1, §14.2 · `SPEC-INVARIAN.md` `INV-15`, `INV-17`, `INV-63`.*

**What to build:** **PK** mencatat potongan atas premi — brokerase, komisi, dan sejenisnya — pada
**bagian non-proporsional maupun ketentuan proporsional**, dan **rumus perhitungannya tertulis
sekali** serta berlaku di kedua pelekatan.

Artefak: entitas `POTONGAN` dengan **dua kolom induk bernama**, `CHECK` tepat satu terisi, dan
`INV-15` sebagai **dua** `UNIQUE`.

**PEMBUAT PERTAMA** untuk `POTONGAN`.

**Kenapa begini:** `POTONGAN` **satu konsep dengan dua pelekatan**, dan buktinya satu aktivitas
melayani keduanya dengan satu rumus: `CalculateDeduction` menghitung
`Deduction = @divide(DeductionPct, 100, 4) * GrossPremiumList(1).Value`, dan pemanggilnya mencakup
kedua sisi — `Share.xml`, `ShareRetro.xml`, `TreatyInXOLAddSpreadingDetail` di non-proporsional,
`DetailLimits.xml` di proporsional. **Rumusnya tidak bercabang.** Menyalinnya menjadi dua akan
membuat keduanya berpisah dalam dua tahun.

Bentuk fisiknya **`KTV-B`**: dua kolom bernama, bukan diskriminator. `INV-17` melarang rujukan yang
sasarannya bergantung nilai kolom lain — diskriminator melanggarnya secara harfiah, dua kolom bernama
tidak. Dan keuntungannya bukan sampingan: pada bentuk lama, **basis data tidak menolak baris yatim**;
dengan dua kolom bernama, **dua kunci asing berdiri** dan baris yatim menjadi mustahil.

**Persyaratan:** **`INV-15`** (jenis potongan unik di dalam satu induk) — terpasang sebagai **dua**
`UNIQUE`, satu per pelekatan · **`INV-17`** · **`INV-63`** (aturan potongan tertulis **sekali**,
berlaku pada kedua pelekatannya), bergolongan **APLIKASI, tinjauan skema** · `INV-52` (premi bruto
dikurangi seluruh potongan sama dengan premi bersih) · `INV-41`

**Tidak termasuk:** **Nama potongan sebagai teks.** Sistem lama menuliskannya tetap —
`DeductionList(1).Comment = "Comm to NuRe"`, `"Brokerage fee"`, `"Facultative Brokerage fee"` — dan di
model baru ia **rujukan ke tabel acuan `JENIS_POTONGAN`** (§14.2). Tetapan itu **hilang dengan
sendirinya**; `INV-59` melarang menyalin namanya.
**Penyebaran** — irisan `38`, bentuk induk gandanya sama tetapi entitasnya lain.
**Kolom `ID_INDUK_POTONGAN`** — **tidak ada lagi**. Ia diganti `ID_BAGIAN` + `ID_DETAIL_PROPORSIONAL`
oleh `KTV-B`; siapa pun yang memulihkannya membalikkan keputusan itu tanpa membukanya.

**Jalur gagal:** Kedua kolom induk terisi -> **ditolak** `CHECK` · Kedua kolom induk kosong ->
ditolak · Dua potongan berjenis sama pada induk yang sama -> **ditolak** `INV-15` · Baris potongan
menunjuk induk yang tidak ada -> ditolak kunci asing.

**Uji:** **Negatif:** kedua induk terisi; kedua induk kosong; jenis potongan kembar pada satu
`BAGIAN`; jenis potongan kembar pada satu `DETAIL_PROPORSIONAL`; induk yatim.
**Positif — dan ia yang menangkap `UNIQUE` yang mencampur dua ruang pengenal:** sebuah `BAGIAN`
bernomor 7 dan sebuah `DETAIL_PROPORSIONAL` bernomor 7, **masing-masing dengan potongan berjenis
sama** -> **keduanya diterima**. Pada bentuk lama — satu `UNIQUE` atas `ID_INDUK_POTONGAN` — keduanya
bertabrakan meskipun keduanya sah, dan itu **menolak data yang sah**.
**Positif kedua:** rumus potongan yang sama menghasilkan angka yang sama di kedua pelekatan,
diperiksa dengan **satu** himpunan data uji yang dijalankan dua kali — bukan dua himpunan.

**Menggantikan:** `P-19` melestarikan `DeductionList` dan rumus `CalculateDeduction`. Yang bergeser:
namanya menjadi **rujukan tabel acuan**, dan induk gandanya memperoleh **kunci asing sungguhan** yang
sebelumnya tidak ada sama sekali.

**Blocked by:** 33 · 34 · 15

**Dasar:**
```
EVIDENCED(CalculateDeduction@ekspor-2026-09 - satu rumus, pemanggil di kedua cabang: Share.xml, ShareRetro.xml, TreatyInXOLAddSpreadingDetail, DetailLimits.xml)
        EVIDENCED(DeductionList Comment="Comm to NuRe"/"Brokerage fee" - nama sebagai tetapan)
        DECIDED(KTV-B, INV-15, INV-17, INV-63, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `POTONGAN` berdiri dengan `ID_BAGIAN` dan `ID_DETAIL_PROPORSIONAL` **keduanya nullable**, masing-masing berkunci asing
- [ ] `CHECK` tepat-satu-terisi terpasang, dan **kedua arah** pelanggarannya diuji
- [ ] `INV-15` terpasang sebagai **dua** `UNIQUE`; uji positif membuktikan keduanya tidak saling mencampur
- [ ] `INV-63` diperiksa lewat **tinjauan skema tertulis**: rumus potongan ditemukan **tepat satu kali** di seluruh basis kode
- [ ] uji positif kedua lulus: satu himpunan data uji, dijalankan di kedua pelekatan, hasilnya sama
- [ ] sapuan membuktikan **tidak ada** kolom `ID_INDUK_POTONGAN` yang tertinggal
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
