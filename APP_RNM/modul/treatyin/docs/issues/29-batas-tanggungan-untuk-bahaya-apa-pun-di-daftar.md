---
status: aktif
golongan: perubahan
---

# 29: Batas tanggungan dicatat untuk bahaya apa pun yang ada di daftar bahaya, termasuk yang ditambahkan kemudian

*Asal: `DAFTAR-PEKERJAAN.md` `P-14` · `SPEC-MODEL-DATA.md` §10.0b, §10.18 · `ADR-0038` · `SPEC-INVARIAN.md` `INV-14`.*

**What to build:** **PK** mencatat batas tanggungan **per bahaya**, untuk **bahaya apa pun yang ada di
tabel acuan** — termasuk bahaya yang ditambahkan sesudah sistem berjalan, **tanpa perubahan skema**.

Artefak: entitas `BATAS_PER_BAHAYA`, kunci asingnya ke `VERSI_KONTRAK` dan `BAHAYA`, kolom
`KODE_MATA_UANG`, dan `INV-14`.

**PEMBUAT PERTAMA** untuk `BATAS_PER_BAHAYA`.

**Kenapa begini:** Sistem lama memakai **delapan kolom untuk empat bahaya** yang **namanya menjadi
nama kolom** — `Earthquake`, `FloodJab`, `FloodNation`, `RSMDLimit`, masing-masing berpasangan dengan
kolom mata uangnya. Bahaya kelima menuntut kolom kesembilan dan kesepuluh, dan bahaya baru adalah hal
yang pasti terjadi di reasuransi. **Bahayanya adalah nilai, bukan nama kolom** — dan itu `ADR-0038`.

**Persyaratan:** `INV-14` (bahaya unik di dalam satu versi) · `INV-36` · `INV-44` · `ADR-0038` ·
`INV-46` (tidak ada pasangan kolom kembar untuk dua mata uang)

**Tidak termasuk:** **Tingkat pencatatan `NILAI_BATAS`.** `INV-39` dan `INV-40` **tidak dipasang**;
yang menjawabnya **`T-2`** dari teknik treaty. Ketiadaannya berdiri sebagai **pernyataan keputusan**
di `ddl-usulan/`, sebab constraint yang dipasang sekarang menuntut nilai yang migrasi tidak tahu cara
mengisinya.
**Glosarium `RSMD`** — istilah pasar yang **dipertahankan** dan **wajib punya entri glosarium**
(`CONTEXT.md` §3.3b). Kepanjangannya **belum tercatat di mana pun**; korpusnya sudah disapu habis dan
hasilnya nol. Itu butir wawancara, bukan pekerjaan irisan ini.

**Jalur gagal:** Dua baris berbahaya sama pada satu versi -> **ditolak** `INV-14` · Nilai batas terisi
tanpa mata uang -> ditolak `INV-36` · Bahaya yang tidak ada di tabel acuan -> ditolak · Sebuah kolom
bernama bahaya di `VERSI_KONTRAK` -> **tidak ada**, diperiksa dengan sapuan.

**Uji:** **Negatif:** baris berbahaya kembar; nilai tanpa mata uang; bahaya tak terdaftar.
**Positif — dan ia pokok irisan ini:** tambahkan **bahaya kesembilan** ke tabel acuan lewat irisan
`15`, lalu catat batas tanggungan untuknya pada sebuah kontrak. **Berhasil tanpa satu pun perubahan
skema.** Inilah satu-satunya uji yang memisahkan tabel acuan yang sungguhan dari daftar nilai yang
diam-diam masih tertanam di tempat lain.

**Menggantikan:** delapan kolom bernama bahaya di kepala `TreatyIn`, beserta keempat kolom
`Currency*`-nya yang **melebur ke paket uangnya** (`COCOK-SILANG-CACAH-ATRIBUT.md` §3.3).
Golongannya **PERUBAHAN** — bentuknya sengaja diubah, jadi data lama menunjukkan bentuk yang dibuang.

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(SPEC-MODEL-DATA.md §10.0b - Earthquake/FloodJab/FloodNation/RSMDLimit + empat kolom Currency*)
        DECIDED(ADR-0038, INV-14, INV-46, KTV-A, KTV-C)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-2)
```

- [ ] `BATAS_PER_BAHAYA` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`, dengan `KODE_MATA_UANG` **boleh kosong** — ia bukan bagian kunci alami (`KTV-C` koreksi 1)
- [ ] `INV-14` terpasang atas bahaya di dalam satu versi
- [ ] uji positif lulus: **bahaya kesembilan** dapat dipakai tanpa perubahan skema
- [ ] sapuan membuktikan **nol** kolom bernama bahaya di `VERSI_KONTRAK`
- [ ] ketiadaan `INV-39`/`INV-40` tertulis sebagai **pernyataan keputusan** yang menyebut `T-2`
- [ ] `RSMD` **terdaftar sebagai butir glosarium yang belum terisi**, bukan didiamkan
- [ ] `KTV-A` dan `T-2` tercatat di `ASUMSI-CLEAR.md`
