---
status: aktif
golongan: pelestarian
---

# 34: Ketentuan proporsional dicatat per kelompok treaty di dalam sebuah layer

*Asal: `DAFTAR-PEKERJAAN.md` `P-21` · `SPEC-MODEL-DATA.md` §10.4, §12.3 · `SPEC-INVARIAN.md` `INV-06`, `INV-32`.*

**What to build:** **PK** mencatat ketentuan cabang proporsional — bagian NuRe, jenis treaty, dan
cadangan premi per mata uang — **per kelompok treaty di dalam sebuah layer**. Baris kembar pada sumbu
itu ditolak, dan cabang ini **hanya ada di bawah kontrak proporsional**.

Artefak: entitas `DETAIL_PROPORSIONAL` dan tabel anak `NILAI_CADANGAN_PREMI`; `INV-06`; `INV-32`.

**PEMBUAT PERTAMA** untuk `DETAIL_PROPORSIONAL`, `NILAI_CADANGAN_PREMI`.

**Kenapa begini:** **Bagian NuRe pada cabang proporsional tidak punya baris sendiri** — ia atribut
`RNMShare` **pada** baris ini, terisi bila `RNMShareAcrossTheBoard` mati. §12.3 menetapkannya sesudah
uji semula dijalankan pada sumbu yang salah. Dan lingkupnya **di dalam layer**, bukan di dalam versi:
`INV-06` berbunyi *"`ID_KELOMPOK_TREATY` unik di dalam satu `LAYER`"* — menulisnya *"di dalam versi"*
akan menolak dua layer yang sama-sama memuat kelompok treaty yang sama, dan itu susunan yang lazim.

**Persyaratan:** `INV-06` (kelompok treaty unik di dalam satu `LAYER`) · **`INV-32`**
(`DETAIL_PROPORSIONAL` hanya di bawah kontrak ber-`SIFAT_PROPORSI` `PROPORSIONAL`), bergolongan
**INDEKS UNIK** · `INV-30` (`JENIS_TREATY` dua nilai) · `INV-36` · `INV-41`

**Tidak termasuk:** **Aturan tepat-satu antara persen quota share dan jumlah lines surplus** — irisan
`35`. Kolomnya berdiri di sini; **aturannya** irisan berikutnya.
**Baris surplus yang menuntut adanya baris quota share** — irisan `36`.
**`KAPASITAS_SURPLUS` sebagai kolom tersimpan** — **tidak ada**: ia **kelipatan** retensi × jumlah
lines, dan `INV-48` menamainya pengecualian bernama terhadap `INV-47`. Turunan tidak disimpan
(`ADR-0037`).
**Potongan dan penyebaran** — irisan `37` dan `38`.

**Jalur gagal:** Dua baris berkelompok treaty sama pada satu layer -> **ditolak** `INV-06` ·
`DETAIL_PROPORSIONAL` di bawah kontrak non-proporsional -> **ditolak** `INV-32` · Cadangan premi
terisi tanpa mata uang -> ditolak `INV-36`.

**Uji:** **Negatif:** baris kembar di dalam satu layer; baris di bawah kontrak non-proporsional;
cadangan premi tanpa mata uang.
**Positif — dan ia yang menangkap lingkup yang terlalu lebar:** **dua layer pada versi yang sama,
masing-masing memuat kelompok treaty yang sama** -> **diterima**. `UNIQUE` yang ditulis di dalam
versi alih-alih di dalam layer menolaknya, dan lulus uji negatif di atas. Ini bentuk yang sama dengan
ketiga kekeliruan lingkup yang sudah tercatat di modul ini.
**Positif kedua:** `BAGIAN` pada kontrak **non-proporsional** tetap diterima — `INV-32` dan `INV-33`
tidak saling menelan.

**Menggantikan:** `P-21` melestarikan `Data-TreatyInLimitsDetail`. Yang bergeser: **cadangan premi
menjadi tabel anak per mata uang** (`P-8` golongan A), dan bagian NuRe proporsional **dinyatakan
sebagai atribut**, bukan entitas — koreksi §12.3 terhadap §8.

**Blocked by:** 31 · 15

**Dasar:**
```
EVIDENCED(Data-TreatyInLimitsDetail@ekspor-2026-09 - RNMShare per baris bila RNMShareAcrossTheBoard mati)
        DECIDED(INV-06, INV-32, INV-48, ADR-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `DETAIL_PROPORSIONAL` dan `NILAI_CADANGAN_PREMI` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-06` terpasang **di dalam `LAYER`**, bukan di dalam versi — diperiksa dengan uji positif, bukan dibaca dari rumusannya
- [ ] `INV-32` terpasang sebagai indeks unik, dan tidak menelan `INV-33`
- [ ] kedua uji positif lulus
- [ ] sapuan membuktikan **tidak ada** kolom `KAPASITAS_SURPLUS` tersimpan; `INV-48` tertulis sebagai pengecualian bernama
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
