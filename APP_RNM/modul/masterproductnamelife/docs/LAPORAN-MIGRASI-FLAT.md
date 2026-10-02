# Laporan uji kering pindah flat — DEV 02-10-2026

> Untuk: work owner dan DBA. Brief `PROMPT-PINDAH-FLAT-MASTER-PRODUCT-NAME-LIFE.md` T7. Alat
> `backend/alat/pindahflat` mode **`-uji`** (bawaan) dijalankan asisten terhadap DEV pada **02-10-2026 18:50**, sesudah
> perbaikan `/code-review`: **hanya SELECT** atas `M_PRODUCT_LIFE` dan `M_PRODUCTINWARD_LIFE`, nol tulisan; tabel flat
> belum ada di DEV (migrasi 140–147 belum dijalankan). Laporan alat **agregat**: ID produk dan nama kolom/kunci, tidak
> satu nilai data. Oracle DEV: **12.2.0.1.0** Enterprise Edition (`PRODUCT_COMPONENT_VERSION`, `V$VERSION`) — pengenal
> ≤ 30 byte aman, `FETCH FIRST` tersedia.
>
> ⚠️ **Putaran pertama (18:25, sebelum `/code-review`) melaporkan `gagal: 0`.** Itu keliru-aman: aturan lama mencacah
> SEMUA kunci baris yang tidak dikelola sebagai "kunci internal Pega" (dibuang D2) dan menerima setiap `MATURE` yang tidak
> berbentuk `dd/MM/yyyy` sebagai K3. Aturan sekarang hanya membuang kunci `px*`/`py*`/`pz*` dan menggagalkan kunci lain
> yang berisi; K3 hanya untuk nilai yang memang tidak terbaca. Hasilnya di bawah.

## Keluaran alat (apa adanya, putaran 18:50)

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
  K3 nilai tidak sah di-NULL-kan: 1
    100017 M_PRODUCTNAME_LIFE_UWLIMIT[25].MAXINSURED: bukan angka
  K4 objek OutwardList kosong dibuang: 189
  normalisasi teks (menunggu keputusan work owner): 2
    M_PRODUCTNAME_LIFE.RICOMM (koma desimal)                   1
    M_PRODUCTNAME_LIFE.RICOMM (nol depan)                      1
    contoh: 100119 M_PRODUCTNAME_LIFE.RICOMM
    contoh: 100174 M_PRODUCTNAME_LIFE.RICOMM
  medan tanpa kolom yang terisi: 752
    baris OutwardList kunci OUTWARDNAME                        189
    baris OutwardList kunci OUTWARDNAMEID                      189
    baris OutwardList kunci OUTWARDRATE                        187
    baris OutwardList kunci OUTWARDRATEID                      187
  dibuang D2 (bukan data): 12561
    Comment (masukan popup; isinya sudah baris CommentList)    66
    IsView (keadaan layar)                                     181
    baris CommentList kunci pxObjClass                         191
    baris DocumentClaim kunci pxCreateDateTime                 1129
    baris DocumentClaim kunci pxCreateOpName                   1129
    baris DocumentClaim kunci pxCreateOperator                 1129
    baris DocumentClaim kunci pxCreateSystemID                 1129
    baris DocumentClaim kunci pxObjClass                       1129
    baris FinancialUnderwritingList kunci pxCreateDateTime     42
    baris FinancialUnderwritingList kunci pxCreateOpName       42
    baris FinancialUnderwritingList kunci pxCreateOperator     42
    baris FinancialUnderwritingList kunci pxCreateSystemID     42
    baris FinancialUnderwritingList kunci pxObjClass           52
    baris LienClause kunci pxCreateDateTime                    5
    baris LienClause kunci pxCreateOpName                      5
    baris LienClause kunci pxCreateOperator                    5
    baris LienClause kunci pxCreateSystemID                    5
    baris LienClause kunci pxObjClass                          5
    baris OutwardList kunci pxObjClass                         191
    baris PlanList kunci pxCreateDateTime                      287
    baris PlanList kunci pxCreateOpName                        287
    baris PlanList kunci pxCreateOperator                      287
    baris PlanList kunci pxCreateSystemID                      287
    baris PlanList kunci pxListSubscript                       286
    baris PlanList kunci pxObjClass                            287
    baris UnderwritingLimitList kunci pxCreateDateTime         440
    baris UnderwritingLimitList kunci pxCreateOpName           440
    baris UnderwritingLimitList kunci pxCreateOperator         440
    baris UnderwritingLimitList kunci pxCreateSystemID         440
    baris UnderwritingLimitList kunci pxObjClass               2169
    kunci halaman inward pxObjClass                            196
    kunci halaman umum pxObjClass                              196
  gagal: 753 (752 tanpa kolom flat)
    100175 M_PRODUCTNAME_LIFE.MATURE: bukan tanggal (bentuk lain yang dapat diselamatkan - bukan K3, menunggu keputusan)
    100004: baris OutwardList kunci OUTWARDRATEID terisi tetapi tidak punya kolom flat
    100004: baris OutwardList kunci OUTWARDNAME terisi tetapi tidak punya kolom flat
    100004: baris OutwardList kunci OUTWARDNAMEID terisi tetapi tidak punya kolom flat
    100004: baris OutwardList kunci OUTWARDRATE terisi tetapi tidak punya kolom flat
    100006: baris OutwardList kunci OUTWARDNAME terisi tetapi tidak punya kolom flat
    100006: baris OutwardList kunci OUTWARDNAMEID terisi tetapi tidak punya kolom flat
    100006: baris OutwardList kunci OUTWARDRATE terisi tetapi tidak punya kolom flat
    100006: baris OutwardList kunci OUTWARDRATEID terisi tetapi tidak punya kolom flat
    ... dan 744 lainnya tanpa kolom flat (cacah per kunci di atas)
  rekonsiliasi: TIDAK LOLOS; ditulis: false
```

## Bacaan

| Hal | Hasil | Sesuai brief bab 0? |
| --- | --- | --- |
| Sumber | 196 produk, 196 baris inward; 0 produk tanpa inward; 0 inward yatim | ya |
| Inward ber-ID lain | 2 — cacahnya cocok dengan pasangan bersilang 100079 ↔ 100081 brief bab 0 (alat mencetak cacah, bukan ID-nya); di tabel flat sisi inward ada di baris produknya sendiri | ya |
| Baris anak | LIEN 5, DOCCLAIM 1129, PLAN 287, FINUW 52, UWLIMIT 2169, OUTWARD 2, COMMENT 191 | ya — sama dengan cacah larik bab 0 |
| K3 | **1** nilai di-NULL-kan: `UnderwritingLimitList[*].MaxInsured` produk 100017 baris 25 (tidak terbaca sebagai angka, juga tanpa pemisah) | ⚠️ brief: 1 + 1 |
| `MATURE` produk 100175 | **bukan K3**: tidak berbentuk `dd/MM/yyyy`, tetapi terbaca sebagai tanggal dalam bentuk lain — dapat diselamatkan | ⚠️ premis K3 keliru → **OQ-FLAT-09** |
| K4 | 189 objek `OutwardList` "kosong" (keenam kunci OR kosong) — tetapi **berisi** `OUTWARDNAME` 189, `OUTWARDNAMEID` 189, `OUTWARDRATE` 187, `OUTWARDRATEID` 187 (agregat terpisah, SELECT `JSON_TABLE`: 4 nama berbeda ≤ 11 karakter, 4 ID angka ≤ 7, 3 rate berbeda ≤ 21 karakter, 3 ID rate angka ≤ 7; **nol** di kedua objek OR ber-`REINSTYPEID`) | ⚠️ premis K4 ("189 objek kosong") keliru → **OQ-FLAT-08** |
| Dibuang D2 | kunci `px*` halaman dan baris (`pxObjClass`, jejak `pxCreate*`, `pxListSubscript`), `IsView` 181, `Comment` 66 | — bukan data |
| Normalisasi | **2** — `RICOMM` produk 100119 (koma desimal) dan 100174 (nol depan): nilainya sama secara angka | **OQ-FLAT-07** |
| Gagal | **753** = 752 kunci `OUTWARD*` tanpa kolom flat (OQ-FLAT-08) + 1 `MATURE` (OQ-FLAT-09) | — |

**Kesimpulan:** `-jalankan` DITAHAN alat — sebagaimana mestinya — sampai tiga keputusan work owner:

| OQ | Keputusan yang diminta | Rekomendasi |
| --- | --- | --- |
| OQ-FLAT-07 | 2 normalisasi `RICOMM` (koma desimal, nol depan) diterima dalam bentuk kanonik? | terima kedua jenis (`-terima-normalisasi="koma desimal,nol depan"`): nilai sama, `RICOMM` persen ≤ 5 karakter |
| OQ-FLAT-08 | 189 objek `OutwardList` ber-`OUTWARDNAME`/`OUTWARDRATE` (bukan baris OR): dipindah atau dibuang? | **pindahkan**: tambah empat kolom `OUTWARDNAMEID`, `OUTWARDNAME`, `OUTWARDRATEID`, `OUTWARDRATE` ke `M_PRODUCTNAME_LIFE_OUTWARD` (migrasi 146, belum pernah dijalankan di DEV) dan K4 hanya membuang objek yang SEMUA kuncinya kosong |
| OQ-FLAT-09 | `MATURE` produk 100175 berbentuk tanggal lain: dikonversi ke tanggal atau di-NULL-kan (K3)? | **konversi** (nilai terselamatkan); bila work owner tetap memilih NULL, K3 tetap dicatat |
