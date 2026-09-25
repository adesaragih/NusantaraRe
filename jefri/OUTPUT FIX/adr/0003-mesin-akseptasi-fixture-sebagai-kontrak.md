---
status: accepted
---

# Mesin akseptasi ditulis lebih dulu; fixture tabel limit menjadi kontrak data ke DBA

`services/acceptance` ditulis **sekarang** dengan fixture tabel limit, tidak menunggu isi
`M_LIMIT_*` tiba dari DBA. Yang hilang dari tangga akseptasi adalah **datanya**, bukan **logikanya**:
query SQL-nya terbaca penuh, urutan eskalasinya terverifikasi (`LIMIT_BOTTOM` menaik, baris pertama =
approver berikutnya), cara berhentinya terverifikasi (tidak ada `To*` yang cocok → cabang `Else` →
tangga selesai, bukan galat), dan ketiga folder korpus berperilaku identik. Fixture yang kita tulis
sekaligus **menjadi kontrak bentuk data** yang dikirim ke DBA.

## Consequences

Ini membalik arah ketergantungan menjadi menguntungkan: alih-alih menunggu data lalu menemukan
bentuknya tidak seperti dugaan, kita menyatakan bentuk yang kita harapkan lebih dulu dan
ketidakcocokan muncul **saat permintaan dikirim**, bukan di akhir.

~~Fixture wajib memuat kolom yang dibaca query: `JABATAN`, `JABATAN_ATASAN`, `LIMIT_BOTTOM`,
`LIMIT_BOTTOM2`, `MAX_LIMIT_IDR`, `MAX_LIMIT_USD`, `BATAS_WAKTU`, `LOGIN` — untuk ketujuh tabel limit
per lini bisnis.~~ → **dikoreksi amandemen di bawah: ada DUA bentuk, bukan satu, dan enam tabel, bukan tujuh.**

~~⚠️ Satu hal yang fixture **tidak** dapat tebak dan harus datang dari DBA: **ejaan pasti nilai kolom
`JABATAN`**.~~ → **sudah datang; lihat amandemen.** Token routing di rule ditulis tanpa spasi
(`DIREKTURTEKNIK`); bila ejaan Oracle berbeda, tidak ada approver yang pernah ditemukan dan tangga
macet total.

Keputusan ini **mencabut** sikap sebelumnya yang menunda `services/acceptance` sampai tabel limit
tiba (lihat K-009).

---

## Amandemen — 17 September 2026, setelah DDL + isi tabel limit diterima

**Keputusan intinya tidak berubah.** Menulis fixture lebih dulu terbukti benar: ketidakcocokan bentuk
memang muncul saat data datang, dan itu tepat tujuannya. Yang berubah adalah **asumsi bentuknya**.

Sumber: `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` (DDL) dan `M_LIMIT_*.xls` (isi), diverifikasi
17 September 2026.

### ⛔ Asumsi "satu bentuk untuk ketujuh tabel" **SALAH**

`[terverifikasi]` Ada **enam** tabel limit di basis data, dalam **dua bentuk yang tidak kompatibel**.

#### Bentuk A — standar, 14 kolom, 5 tabel

`M_LIMIT_PROPERTYY` · `M_LIMIT_ENGINEERINGG` · `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` ·
`M_LIMIT_PROPERTY_NON_PREFERREDD` · `M_LIMIT_NONPROPANDENGG`

```
ID · JABATAN · MAX_LIMIT_IDR · LIMIT_BOTTOM · MAX_LIMIT_USD · BATAS_WAKTU · NAMA
TGL_UPDATE · EFFECTIVE_DATE · TEAM_GROUP · LOGIN · WORKBASKET · JABATAN_ATASAN · LIMIT_BOTTOM2
```

`[terverifikasi]` **Ejaan `JABATAN` TANPA SPASI** — cocok dengan token routing di rule:
`DIREKTURTEKNIK` · `DIREKTURMARKETING` · `KADIVTEKNIK` · `KADIVFACULTATIVE` · `DEPHEADUNDERWRITER` ·
`MANAGERTEKNIK` · `SENIORUW` · `UNDERWRITER` · `LEADER` · `JUW_B`.

`WORKBASKET` **terisi** dan cocok dengan ruang nama antrean (`ReasFacInUnderwriting`,
`ReasFacInSeniorUnderwriting`, …), sehingga kedua ruang nama — antrean dan kode jabatan — tersedia
langsung dari tabel ini.

#### Bentuk B — financial, 7 kolom, 1 tabel

`M_LIMIT_FINANCIALINS`

```
ID · JABATAN · LIMITBOND_BOTTOM · LIMITCREDITCL_BOTTOM · LIMITCREDITNCL_BOTTOM · LIMITTRADE_BOTTOM · NAMA
```

`[terverifikasi]` Bentuk ini **tidak punya** `WORKBASKET`, `JABATAN_ATASAN`, `TEAM_GROUP`,
`EFFECTIVE_DATE`, maupun `LOGIN`. Kolom limitnya **sama sekali berbeda**: empat ambang per jenis
pertanggungan (bond · credit CL · credit NCL · trade), bukan `MAX_LIMIT_IDR`/`LIMIT_BOTTOM`.

⛔ `[terverifikasi]` **Ejaan `JABATAN` PAKAI SPASI**: `DIREKTUR TEKNIK` · `DIREKTUR MARKETING` ·
`KADIV KEUANGAN` · `SENIOR UNDERWRITER`.

### Konsekuensi yang mengikat implementasi

1. **Mesin akseptasi financial mencocokkan jabatan berspasi, bukan token tanpa spasi.** Bila kedua
   bentuk disamakan — satu fungsi pencocokan untuk semua lini — **tidak ada approver financial yang
   pernah ditemukan**, dan seluruh tangga lini financial macet. Ini persis kegagalan yang ADR ini
   peringatkan; ia nyata, dan sekarang terbukti.
2. **Normalisasi ejaan DILARANG.** Menghapus spasi agar "seragam" adalah perubahan perilaku, bukan
   migrasi (`CLAUDE.md` §1). Kedua ejaan direproduksi apa adanya, masing-masing pada bentuknya.
3. **Bentuk B tidak punya rantai atasan.** Tanpa `JABATAN_ATASAN`, tangga financial **tidak dapat
   menaik lewat mekanisme yang sama**. Bagaimana eskalasinya bekerja — atau apakah memang tidak ada
   eskalasi — **belum terjawab dari data ini** dan ditandai `[pertanyaan terbuka]`.
4. **Bentuk B tidak punya antrean.** Tanpa `WORKBASKET`, tidak ada token antrean yang dapat dibaca
   untuk lini financial.
5. **Fixture tetap dipakai untuk bentuk B** sampai butir 3 dan 4 terjawab; bentuk A boleh berpindah ke
   data nyata.

### `M_LIMIT_LIFE` bukan kekurangan

`[terverifikasi]` Dikonfirmasi **tidak ada** di basis data. Target efektif **enam** tabel, bukan tujuh.

### Kolom yang tidak boleh masuk fixture

⛔ Kolom **`NAMA`** dan **`LOGIN`** memuat **nama orang**. Keduanya **wajib dibuang** saat fixture
dibangun dari berkas isi tabel — `CLAUDE.md` §3 butir 5 berlaku penuh pada fixture dan test. Lihat
**K-025**.
