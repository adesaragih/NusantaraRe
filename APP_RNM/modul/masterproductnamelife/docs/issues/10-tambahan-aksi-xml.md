# 10 (tambahan): aksi XML yang belum dicakup tiket 01–09

**Status:** selesai (01-10-2026) — paket 9 `14d3d48`, layar paket 10 (`1ada8d7`)

**Sumber:** brief bab 5 baris *"tambahan"* — `CopyProduct`, `ViewRate`, `View Reinstype`, `Generate`, `Refresh`,
`EditProductName_Confirm`, `SaveProductName_Confirm`: *"bila XML menunjukkan perilakunya dan tiket belum mencakupnya,
**dibangun** dan dicatat sebagai tambahan tiket"*. Ditambah checkbox `On Retention` (RALAT R9/P3, dijadwalkan paket 9).

Nomor baris `bNNN` = `sed -e 's/></>\n</g'`; blok langkah `·` = `pyStepsBlockName` kosong (langkah jalan), `//` = ter-remark.

## Hasil per aksi

| Aksi (bNNN · teks VERBATIM) | Bukti XML | Dibangun | Keadaan |
| --- | --- | --- | --- |
| **`Copy`** b59854 | DataTransform `CopyProduct` b59965: 1 b144 `ProductName.ID := ""`, 2 b173 `ProductNameInward.ID := ""`, 3 b203 `OutputParam.STSSAVE := 99`, 4 b233/b236 `OutputParam.ERRMSG := "Data sudah dicopy, silakan melakukan perubahan dan tekan SAVE untuk menyimpan"` (wadah b1269 `STSSAVE==99`) | `POST /produk` dengan `salinanDari`: produk baru (ID dari sequence, inward ID = ID produk R14); medan milik server (medan mati, `CommentList`, `OutwardList`, medan inward mati) **diwarisi** dari produk asal seperti halaman Pega yang tersalin utuh; `CREATEOP` = pelaku (OQ-MPNL-13). Pesan VERBATIM dan pengosongan ID = layar (paket 10) | ✅ backend |
| **`Generate`** b60122 | `GenerateUpload_Act` b60228: 1 b226 `·` Page-Remove; 2 b332 `·` `CARI1..CARI40`; 3 b1285 `·` `pxConvertResultsToCSV` `FileName=SeeDetail` b1335, `CSVPropHeaders` b1340, `CSVProperties` b1343 | `POST /produk/generate` (badan = isi form saat itu, tanpa simpan dan tanpa validasi) → `SeeDetail.csv`, 39 judul VERBATIM (uji membaca b1340 korpus), nilai menurut nama judul, `@if` `Basic/Rider`, `QS/SUPRLUS`, `KREDIT LIFE/NON KREDIT LIFE/PA/HEALTH`, `Annual/Semi Annual/Quarterly/Monthly/Single` | ✅ backend — OQ-MPNL-06 |
| checkbox **`On Retention`** b47312 (`pyCheckboxCaption` b47312) | onChange `GetReinsTypeOR_Life` b47488: 1 b236 `·` PRE=false hapus `OutwardList`; 2 b383 `·` `CARI1/2 ← BEGIN/MATURE`; 3 b534 `·` PRE=false `BrowseReinstypeOR_SQL`; 4 b731 `·` PRE=false ulang; 4.1 b770 `·` tambah baris | `hitungOutward` di badan simpan → `OutwardList` diganti hasil `BrowseReinstypeOR_SQL` b84 (`REINSTYPEID = '10200'`, periode meliputi `BEGIN`–`MATURE`), **dicentang maupun dilepas** (PRE=false); `OVR_COMM = ""`; `TREATYCONTRACTID = ""` (OQ-MPNL-15) | ✅ backend — OQ-MPNL-09 |
| **`View Rate`** b34113 | `SetParamRate` b34310 + `localAction ViewRate` b34354 — sumber rate `RATE_LIFE` (view atas JSON) | `GET /rate` = 503 berkalimat | ⏸️ OQ-MPNL-03 (paket 6) |
| **`View Reinstype`** b50764 / b55382 | grid `OUTWARD` di wadah `1==2` (b47876 / b52491); `SetParamOutward` tidak ada di korpus | — | ➖ **mati** — tidak dibangun |
| **`Refresh`** b65270 | `LoadAttachmentProdName` b65401 | `GET /produk/{id}/lampiran` | ✅ paket 8 |
| **`Edit`** b59489 → `EditProductName_Confirm` | FlowAction b19 submit **`Edit`**, b18 **`Cancel`**, section b496 **`Do you want to Edit the data?`**, transform `SetViewEdit` b71 → 1 b151 `IsView := "false"` | dialog konfirmasi di layar; simpan sesudah `Edit` menulis `IsView = "false"` (paket 3) | ⏳ layar paket 10 |
| **`Save`** b59041 → `SaveProductName_Confirm` | FlowAction b34 submit **`Save`**, b33 **`Cancel`**, section b519 **`Do you want to save the data?`**, b1025 area teks **`Comment`** | dialog konfirmasi + `Comment` → `POST`/`PUT` (`comment` dipakai `AddCommentList_Act`) | ⏳ layar paket 10 |

## Acceptance criteria

- [x] Salinan = produk baru: ID dari sequence, kedua ID tidak pernah dari produk asal; produk asal tidak tersentuh.
- [x] Salinan mewarisi medan milik server produk asal; `CREATEOP` = pelaku; `salinanDari` pada `PUT` → 400; produk asal tidak ada → 404 berkalimat.
- [x] `Generate` tidak menyimpan apa pun, menuntut identitas, judul = b1340 korpus, CRLF.
- [x] `On Retention` diubah → `OutwardList` dari kontrak OR yang periodenya meliputi `BEGIN`–`MATURE`; tidak diubah → `OutwardList` tersimpan dipertahankan; tanggal kosong → kosong; kontrak tak terbaca → simpan gagal terang (503), bukan daftar kosong diam-diam.
- [x] Penanda permintaan (`salinanDari`, `hitungOutward`) tidak pernah tersimpan di `JSONDATA`.

## Seam & perintah verifikasi

```
go test ./modul/masterproductnamelife/...
go test -tags db -p 1 ./modul/masterproductnamelife/...   (MELEWATI tanpa ORACLE_DSN)
```

## Status 01-10-2026 — K1 keputusan work owner 01-10-2026 (OQ-MPNL-03)

Kalimat lama *"`GET /rate` = 503 berkalimat"* tidak berlaku lagi: `Choose R/I Rate` membaca view `RATE_LIFE_SUMMARY` dan `View Rate`
membaca view `RATE_LIFE`, baca saja, kolom RD saja (`GET /master/ri-rate`, `GET /rate` = 200). Baris `PLAN LIST` baru dapat
diberi R/I Rate; pilihan baru wajib ada di view dan namanya diambil dari master (`SetRIRate` b2448). `View Rate` disaring
`RIRATEID` baris plan — penyimpangan sadar dari `ViewRate.xml` b1024 (`ParamID.OUTWARDRATEID` tidak pernah diisi; PARITAS bab
keputusan OQ 01-10-2026).
