# Fixture tabel limit akseptasi

Ekspor `POOLDATA.M_LIMIT_*` dari work owner, 1 Oktober 2026 (`docs/KEPUTUSAN-30-09-2026.md` butir 42),
disalin **byte demi byte** dari `D:\migrasi\RNM\DDL\*.csv`. `NAMA` dan `LOGIN` sengaja tidak diekspor —
dijaga `TestFixtureLimitTanpaNamaLogin`. Kolom lain = yang dibaca rule SQL Pega (`NB FacIn\RDBList\`),
**ditambah** `WORKBASKET` dan `JABATAN_ATASAN` di kelima tabel bentuk A: keduanya hanya dibaca
`GetAksepBanding_SQL` (atas `M_LIMIT_PROPERTYY`) dan **tidak** dipakai `acceptance.Next`.

| Berkas | Dipakai | Baris |
| --- | --- | ---: |
| `M_LIMIT_PROPERTYY.csv` | NB-11 bentuk A | 37 |
| `M_LIMIT_PROPERTY_NON_PREFERREDD.csv` | NB-11 bentuk A | 37 |
| `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL.csv` | NB-11 bentuk A | 37 |
| `M_LIMIT_ENGINEERINGG.csv` | NB-11 bentuk A | 37 |
| `M_LIMIT_NONPROPANDENGG.csv` | NB-11 bentuk A | 37 |
| `M_LIMIT_FINANCIALINS.csv` | NB-12 bentuk B (Bond/Kredit) | 4 |

**Format — wajib diketahui pemuat:**

- UTF-8 tanpa BOM, akhir baris CRLF, pemisah **`;`**, satu baris header.
- Angka memakai **titik sebagai pemisah ribuan** (`1.000.000`). `[terverifikasi]` setiap nilai bertitik
  punya ≥ 2 titik berpola ribuan (541 nilai), 581 bilangan bulat polos, nol bentuk lain — jadi tidak ada
  desimal. Hapus titiknya sebelum diurai; **jangan** membacanya sebagai desimal.
- `JABATAN` bentuk A **tanpa spasi** (`DIREKTURTEKNIK`); `M_LIMIT_FINANCIALINS` **berspasi**
  (`DIREKTUR TEKNIK`). Dibandingkan persis, tanpa *trim*.
