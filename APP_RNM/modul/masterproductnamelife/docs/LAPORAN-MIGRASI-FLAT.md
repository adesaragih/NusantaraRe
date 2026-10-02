# Laporan uji kering pindah flat — DEV 02-10-2026

> Untuk: work owner dan DBA. Brief `PROMPT-PINDAH-FLAT-MASTER-PRODUCT-NAME-LIFE.md` T7. Alat
> `backend/alat/pindahflat` mode **`-uji`** (bawaan) dijalankan asisten terhadap DEV pada **02-10-2026 18:25**, kode
> commit `353f9e80`: **hanya SELECT** atas `M_PRODUCT_LIFE` dan `M_PRODUCTINWARD_LIFE`, nol tulisan, tabel flat belum
> ada di DEV (migrasi 140–147 belum dijalankan). Laporan alat **agregat**: ID produk dan nama kolom, tidak satu nilai data.
> Oracle DEV: **12.2.0.1.0** Enterprise Edition (`PRODUCT_COMPONENT_VERSION`, `V$VERSION`) — pengenal ≤ 30 byte aman,
> `FETCH FIRST` tersedia.

## Keluaran alat (apa adanya)

```
pindahflat -uji
  sumber: M_PRODUCT_LIFE 196 baris, M_PRODUCTINWARD_LIFE 196 baris
  produk tanpa inward: 0; inward ber-ID lain (dicari lewat PRODUCTID): 2; inward yatim: 0
  baris yang (akan) ditulis:
    M_PRODUCTNAME_LIFE             196
    M_PRODUCTNAME_LIFE_LIEN        5
    M_PRODUCTNAME_LIFE_DOCCLAIM    1129
    M_PRODUCTNAME_LIFE_PLAN        287
    M_PRODUCTNAME_LIFE_FINUW       52
    M_PRODUCTNAME_LIFE_UWLIMIT     2169
    M_PRODUCTNAME_LIFE_OUTWARD     2
    M_PRODUCTNAME_LIFE_COMMENT     191
  K3 nilai tidak sah di-NULL-kan: 2
    100017 M_PRODUCTNAME_LIFE_UWLIMIT[25].MAXINSURED: bukan angka
    100175 M_PRODUCTNAME_LIFE.MATURE: bukan tanggal
  K4 objek OutwardList kosong dibuang: 189
  normalisasi teks (menunggu keputusan work owner): 2
    M_PRODUCTNAME_LIFE.RICOMM (koma desimal)                   1
    M_PRODUCTNAME_LIFE.RICOMM (nol depan)                      1
    contoh: 100119 M_PRODUCTNAME_LIFE.RICOMM
    contoh: 100174 M_PRODUCTNAME_LIFE.RICOMM
  medan tanpa kolom yang terisi: 0
  dibuang D2 (bukan data): 4663
    Comment (masukan popup; isinya sudah baris CommentList)    66
    IsView (keadaan layar)                                     181
    baris CommentList berkunci internal Pega                   191
    baris DocumentClaim berkunci internal Pega                 1129
    baris FinancialUnderwritingList berkunci internal Pega     52
    baris LienClause berkunci internal Pega                    5
    baris OutwardList berkunci internal Pega                   191
    baris PlanList berkunci internal Pega                      287
    baris UnderwritingLimitList berkunci internal Pega         2169
    kunci halaman inward pxObjClass                            196
    kunci halaman umum pxObjClass                              196
  gagal: 0
  rekonsiliasi: TIDAK LOLOS; ditulis: false
```

## Bacaan

| Hal | Hasil | Sesuai brief bab 0? |
| --- | --- | --- |
| Sumber | 196 produk, 196 baris inward; 0 produk tanpa inward; 0 inward yatim | ya (196/196, 3 yatim sudah dihapus work owner) |
| Inward ber-ID lain | 2 — cacahnya cocok dengan pasangan bersilang 100079 ↔ 100081 brief bab 0 (alat mencetak cacah, bukan ID-nya); dipasangkan lewat `PRODUCTID`; di tabel flat sisi inward ada di baris produknya sendiri | ya |
| Baris anak | LIEN 5, DOCCLAIM 1129, PLAN 287, FINUW 52, UWLIMIT 2169, OUTWARD 2, COMMENT 191 | ya — sama dengan cacah larik bab 0 (OUTWARD sesudah K4) |
| K3 | 2 nilai di-NULL-kan: 1 `MATURE` bukan tanggal (produk 100175), 1 `UnderwritingLimitList[*].MaxInsured` bukan angka (produk 100017, baris 25) | ya (1 + 1) |
| K4 | 189 objek `OutwardList` kosong dibuang (dari 191; 2 berisi dipindah) | ya |
| Medan tanpa kolom yang terisi | 0 — tidak ada nilai yang hilang karena tidak berkolom | ya (bab 0: nol terisi) |
| Dibuang D2 | 4663 kunci internal Pega (`pxObjClass` halaman dan baris), `IsView` (keadaan layar) 181, `Comment` masukan popup 66 (isinya sudah baris `CommentList`) | — bukan data |
| Gagal | **0** | — |
| Normalisasi | **2** — `RICOMM` produk 100119 (koma desimal) dan 100174 (nol depan): nilainya sama secara angka, teksnya berubah ke bentuk kanonik | **menahan `-jalankan`** — keputusan **OQ-FLAT-07** |

**Kesimpulan:** rekonsiliasi teks demi teks lolos untuk seluruh 196 produk kecuali dua normalisasi `RICOMM` yang menunggu
keputusan work owner. Sesudah OQ-FLAT-07 diputuskan "terima", `-jalankan -terima-normalisasi` dapat dijalankan menurut
`PANDUAN-PINDAH-FLAT.md`; bila "jangan", kedua nilai diperbaiki di sumber lebih dulu (di luar alat ini).
