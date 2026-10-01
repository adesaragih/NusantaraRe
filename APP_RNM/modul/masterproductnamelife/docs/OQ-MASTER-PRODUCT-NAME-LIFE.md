# Pertanyaan terbuka — Master Product Name Life

> Register OQ modul ini (awalan `OQ-MPNL-`), dibuka 01-10-2026 sesi implementasi. OQ lintas proyek tetap di
> `discovery/open-questions.md`. ⛔ Hanya work owner (atau pemilik yang disebut) yang menutup OQ; asisten mencatat **bawaan**
> yang dibangun sampai jawaban datang, dan bawaan itu dapat dibalik tanpa migrasi skema.
> Bukti setiap butir: `PARITAS-LAYAR-DAN-AKSI.md` dan `RALAT-DEV-30-09-2026.md`.

| OQ | Pertanyaan | Bawaan sampai dijawab | Pemilik | Status |
| --- | --- | --- | --- | --- |
| OQ-MPNL-01 | Produk tetap JSON di dua tabel lama seperti Pega, atau dipindah ke tabel relasional baru (tiket 01)? Tabel baru membuat tiga view, dua prosedur, dan Claim Life tidak melihat produk baru, kecuali kedua bentuk ditulis bersamaan. | **JSON seperti Pega** (P1): `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE`, `JSONDATA` berkunci Pega, kolom datar `RIRISKID`/`RIRISK`/`PRODUCTNAME`/`BEGIN_DATE`; nol tabel baru, nol DDL | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — produk tetap **JSON seperti Pega**; tiket 01 tetap ditangguhkan |
| OQ-MPNL-02 | `ID` baris inward: Pega mengambilnya dari `M_PRODUCT_INWARD_LIFE_SEQ` (prosedur inward, `SaveProductName_Act` 16 b3052) dan menautkan lewat `PRODUCTID`, sedangkan Claim Life mencari inward `WHERE ID = <ID produk>` dan fakta bisnis menyebut kedua sisi berbagi identitas. DEV 196 vs 197 baris — keduanya pernah tidak sejalan. | `ID` inward = `ID` produk = `PRODUCTID`; `M_PRODUCT_INWARD_LIFE_SEQ` tidak dipakai; produk baru yang ID-nya sudah ada di salah satu tabel **ditolak** (bukan digandakan). Produk lama dibaca inward-nya lewat `PRODUCTID` (seperti `BrowseProductInward` b840) | work owner / DBA | ditutup 01-10-2026 (data DEV, lanjutan 1 L2) |
| OQ-MPNL-03 | Pemilih **R/I Rate** (`BrowseRateLifeSummary`, kelas `RATE_LIFE_SUMMARY`) dan **View Rate** (`BrowseRateLife_RD`, kelas `M_RATE_LIFE`) membaca view atas JSON rate — penjaga Claim Life `TestMasterViewTidakDisentuh` menuntut persetujuan (preseden OQ-MCRL-13). Boleh dibaca modul ini? | rute dan layar dibangun, menjawab **503 berkalimat**; `RIRATE`/`RIRATEID` baris plan yang sudah tersimpan dipertahankan. ⚠️ Akibat: baris plan **baru** tidak dapat diberi R/I Rate, dan `ProteksiPlanListLife` (*"RI/RATE tidak boleh kosong"*) menolak simpannya | work owner | ✅ **ditutup 01-10-2026 — K1 keputusan work owner 01-10-2026** (*"izinkan membaca view rate, baca-saja"*), dibangun `628c148` *(riwayat: terbuka — ✅ **ditutup 01-10-2026 — K1 keputusan work owner 01-10-2026** (*"izinkan membaca view rate, baca-saja"*): `Choose R/I Rate` membaca view `RATE_LIFE_SUMMARY`, `View Rate` membaca view `RATE_LIFE`, kolom RD saja; baris plan baru dapat diberi R/I Rate)* |
| OQ-MPNL-04 | Objek fisik master dibaca menurut konvensi nama kelas `ASM-FW-GISFW-Int-<X>`: `AGENT` (Ceding, SOB), `CLIENT` (Policy Holder), `CURRENCY`, `CAUSEOFLOSS_LIFE`, `RIRISK_LIFE_SUMMARY`, `PRODUCT_TYPE_LIFE` (plan). Tidak satu pun terbukti di katalog yang tercatat (inventori `[dugaan]`). | dibaca dengan nama itu; objek yang tidak ada = **503** yang menyebut objeknya (bukan daftar kosong). Nilai master diverifikasi hanya bila berubah dari yang tersimpan | DBA | ditutup 01-10-2026 (katalog DEV, L3) |
| OQ-MPNL-05 | Daftar pilihan bersumber `associated` (Property tidak ikut ekspor): `Product Name` b3620, `Premium Payment Method` b25611, `Birthday` b27960, `Document List` b15286. | `Product Name`, `Birthday`, `Document List`: isian teks bebas (nilai lama tampil apa adanya). `Premium Payment Method`: kode `1` Annual, `2` Semi Annual, `3` Quarterly, `4` Monthly (dari `GenerateUpload_Act` b332 `CARI37`); kode untuk *Single* tidak diketahui — nilai lain ditampilkan apa adanya | pemilik ekspor Pega | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — ⏳ **konfirmasi menyusul: pemilik ekspor Pega (daftar pilihan `associated`)** |
| OQ-MPNL-06 | `Generate` (`GenerateUpload_Act` 3 b1285): 39 judul (b1340) untuk 40 properti (b1343), dan kolom ke-12 `BENEFIT` memetakan `CARI2` — `pxConvertResultsToCSV` mengisi kolom menurut posisi, sehingga `BENEFIT`…`TREATYNUMBER` bergeser. | `SeeDetail.csv` dengan **39 judul VERBATIM**, tiap kolom berisi nilai yang **namanya** disebut judul (`LIENCLAUSE` = `CARI36`, tanpa judul, tidak ikut); ekspresi `@if` (`Basic/Rider`, `QS/SUPRLUS`, `KREDIT LIFE/…/HEALTH`, `Annual/…/Single`) ditiru | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MPNL-07 | `Download All` b67657 menjalankan jalur treaty-in (`M_ATTACHMENTTREATY_2` milik `TreatyIn.ID`) — salah ekspor (OQ-056). | mengunduh seluruh lampiran **produk ini** sebagai satu arsip zip (AC tiket 08) | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MPNL-08 | Kolom datar `PRODUCTNAME` dan `BEGIN_DATE` `M_PRODUCT_LIFE`: korpus nol penulis (`SaveProductNameLIfeFlat` hanya `RIRISKID`/`RIRISK`), prosedur hanya `JSONDATA` (`dba-procedures-and-ddl.md` §1). | ditulis menurut P1: `PRODUCTNAME` = JSON `PRODUCTNAME`, `BEGIN_DATE` = `BEGIN` inward (`dd/MM/yyyy` → `DATE`), NULL bila kosong | work owner / DBA | ditutup 01-10-2026 (katalog DEV, L1/L4) |
| OQ-MPNL-09 | `OutwardList`: Pega mengisinya saat checkbox `On Retention` diubah (nilai `BEGIN`/`MATURE` saat itu). `OVR_COMM` yang dibaca view tidak punya penulis di korpus (`SetParamOutward` tidak diekspor). | dihitung **saat simpan** bila checkbox diubah dalam sesi sunting itu (memakai `BEGIN`/`MATURE` yang disimpan), selain itu dipertahankan; `OVR_COMM` `""` | work owner | ditutup 01-10-2026 (data DEV, L5) |
| OQ-MPNL-10 | Unggah lampiran: `InsertGoogleStorage_Act` mengirim berkas ke layanan luar (`ServiceGoogle` + `LinkService`) dan mencatat `URLPUBLIC`. Alamatnya di `M_LINK_SERVICE` (OQ-047). | **stub outbox**: rekam `M_ATTACHMENTPRODUCTNAME` + antre efek di satu transaksi; pelaksana stub menyimpan berkas di folder lokal `UNGGAHAN_DIR` dan mencatat `T_STORAGE_IMAGE` (`URLPUBLIC` kosong, `APPFOLDER = Contract`, `STORAGE = standard`). Gagal → status `gagal`, dapat diulang | work owner / tim inti | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MPNL-11 | `View Office Online` membungkus URL bertanda tangan ke penampil kantor di luar (`DownloadAttProdName_Act` 7 b1080). | tautan tampil untuk jenis `xls/xlsx/doc/docx/ppt/pptx` (b69291); rute menjawab **503 berkalimat**: penampil luar tidak dipanggil, alamatnya tidak ditulis | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MPNL-12 | `Medical` baris batas underwriting: tiket 06 menuntut hanya `FCL`/`NM`/`MEDIS`, XML teks bebas (b44814). | **teks bebas** (R16) | Product + UW | ditutup 01-10-2026 (data DEV, L6) |
| OQ-MPNL-13 | `Copy` (`CopyProduct`) mengosongkan kedua `ID` tetapi tidak `CREATEOP`; `SaveProductName_Act` 6 b1370 hanya mengisi `CREATEOP` bila kosong — salinan mewarisi pembuat produk asal. | produk **baru** (termasuk salinan) ber-`CREATEOP` = pelaku (ADR-0007); selebihnya salinan ikut XML — medan mati, `CommentList`, `OutwardList` produk asal diwarisi (paket 9) | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MPNL-14 | `AddCommentList_Act` dipanggil di **setiap** simpan (`SaveProductName_Act` 7 b1515, tanpa prakondisi) — juga bila `Comment` kosong. | ikut XML: setiap simpan menambah satu baris `CommentList` (`Date`, `OperatorName` = akun pelaku, `Suggest` = komentar, boleh kosong) | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MPNL-15 | `GetReinsTypeOR_Life` 4.1 b770 menyalin `.TREATYCONTRACTID` baris hasil `BrowseReinstypeOR_SQL` b84 (`SELECT tc.*, ty.*`), padahal `TREATYCONTRACT_LIFE` dan `TREATYYEAR_LIFE` tidak punya kolom bernama `TREATYCONTRACTID` (DDL `[data DBA]`) — di Pega nilainya selalu kosong. | ikut XML: `OutwardList[*].TREATYCONTRACTID = ""`; alternatifnya `tc.ID` | work owner | ditutup 01-10-2026 (katalog DEV, L7) |

## Pemeriksaan bertanggal 01-10-2026 — paket 11

Kelima belas butir **terbuka**; bawaan masing-masing sudah dibangun (paket 1–10) dan dapat dibalik tanpa migrasi skema.
Yang paling membatasi pemakaian hari ini:

| OQ | Akibat bila tidak dijawab |
| --- | --- |
| OQ-MPNL-03 | baris `PLAN LIST` **baru** tidak dapat diberi R/I Rate, sehingga `ProteksiPlanListLife` (`RI/RATE tidak boleh kosong`) menolak simpannya — produk baru hanya dapat disimpan tanpa baris plan *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: akibat ini gugur — baris plan baru dapat diberi R/I Rate dari view)* |
| OQ-MPNL-10 / OQ-047 | berkas lampiran tersimpan di folder stub `UNGGAHAN_DIR`, bukan di penyimpanan nyata; `URLPUBLIC` kosong |
| OQ-MPNL-11 | `View Office Online` menjawab 503 berkalimat |
| OQ-MPNL-05 | `Product Name`, `Document List`, `Birthday` isian teks (daftar `associated` tidak ikut ekspor); `Premium Payment Method` hanya kode 1–4 |
| OQ-MPNL-15 | `OutwardList[*].TREATYCONTRACTID` kosong seperti Pega |

## Pemeriksaan bertanggal 01-10-2026 — lanjutan 1 (data DEV)

> Sumber: brief `PROMPT-LANJUTAN-MASTER-PRODUCT-NAME-LIFE-1.md` §1–§2 (katalog DEV `ALL_TAB_COLUMNS`/`ALL_OBJECTS` dan agregat `JSONDATA`, dibaca asisten, baca-saja). Penutupan atas permintaan work owner di brief itu (§2); baris register di atas memuat statusnya, pertanyaan lamanya tetap.

| OQ | Bukti DEV | Keputusan |
| --- | --- | --- |
| OQ-MPNL-02 | 196/196 produk punya baris inward ber-`ID` sama; 197/199 baris inward ber-`PRODUCTID` = `ID` | inward ber-`ID` = ID produk, `PRODUCTID` = `ID` (sudah begitu, `mpnl_identitas.go`); uji `TestInwardBerIDSamaDenganProduk` |
| OQ-MPNL-04 | tabel `AGENT`, `CLIENT`; view `CURRENCY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE`, `RIRISK_LIFE_SUMMARY` | objek bernama persis itu dibaca; uji `TestObjekMasterAdaDiKatalogDEV`. ⚠️ Yang terbukti hanya NAMA objek — kolom yang dibaca pemilih (`CLIENTNAME`, `STATUSACTIVE`, `BU_NOTE`, `USEDBY`, …) dan `BrowseReinstypeOR_SQL` belum dicocokkan katalog DEV (catatan code-review) |
| OQ-MPNL-08 | `M_PRODUCT_LIFE` hanya `ID`, `JSONDATA`, `RIRISKID`, `RIRISK` | kolom datar hanya dua (`6fd539c`) |
| OQ-MPNL-09 | 153 produk ber-`OutwardList[0]`, 0 ber-`OVR_COMM` terisi | `OVR_COMM` kosong seperti Pega; uji `TestOutwardOvrCommDanTreatyContractIDKosong` |
| OQ-MPNL-12 | > 60 nilai berbeda `Medical` | teks bebas; tiket 06 diralat |
| OQ-MPNL-15 | `TREATYCONTRACTID` tidak ada di `TREATYCONTRACT_LIFE` maupun `TREATYYEAR_LIFE` | tetap kosong seperti Pega |
| OQ-MPNL-03 | view `RATE_LIFE_SUMMARY` dan `RATE_LIFE` ada (`M_RATE_LIFE_SUMMARY` 339 baris, `M_RATE_LIFE` 96.038) | **tetap terbuka** — register ini tidak memuat izin bertanggal work owner; `R/I Rate`/`View Rate` tetap 503 *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: izin tercatat; ditutup — DEV baca-saja: `ri-rate` 200, 346 baris; `rate` 200, 59 baris untuk satu RIRATEID)* |
| OQ-MPNL-05 | (temuan uji manual) DEV memuat `PAYMENT = 5` (mis. produk `100023`) | kemungkinan kode `Single` — tetap terbuka; nilai di luar daftar ditampilkan, tidak dibuang |

Masih terbuka: OQ-MPNL-01, 03, 05, 06, 07, 10, 11, 13, 14. *(ralat 01-10-2026: seluruhnya ditutup — lihat bab keputusan work owner 01-10-2026 di bawah)*

## Keputusan work owner 01-10-2026 (`PROMPT-LANJUTAN-TIGA-MODUL-LIFE-KEPUTUSAN-OQ.md` §2)

*"Ikuti rekomendasi"* — 8 OQ ditutup dengan bawaan yang sudah dibangun; OQ-MPNL-03 diputuskan berbeda (K1, R/I Rate dan View Rate dari
view rate, `628c148`). Status lama dikutip:

| OQ | Status sebelum 01-10-2026 |
| --- | --- |
| OQ-MPNL-01 | terbuka |
| OQ-MPNL-05 | terbuka |
| OQ-MPNL-06 | terbuka |
| OQ-MPNL-07 | terbuka |
| OQ-MPNL-10 | terbuka |
| OQ-MPNL-11 | terbuka |
| OQ-MPNL-13 | terbuka |
| OQ-MPNL-14 | terbuka |

**Konfirmasi menyusul**: OQ-MPNL-05 — pemilik ekspor Pega (daftar pilihan `associated` `Product Name`, `Premium Payment Method`,
`Birthday`, `Document List`; kode `PAYMENT = 5` di DEV).

Masih terbuka: **nol** OQ modul ini.
