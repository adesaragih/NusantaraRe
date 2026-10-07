# Pertanyaan terbuka — Master Product Name Life

> Register OQ modul ini (awalan `OQ-MPNL-`), dibuka 01-10-2026 sesi implementasi. OQ lintas proyek tetap di
> `discovery/open-questions.md`. ⛔ Hanya work owner (atau pemilik yang disebut) yang menutup OQ; asisten mencatat **bawaan**
> yang dibangun sampai jawaban datang, dan bawaan itu dapat dibalik tanpa migrasi skema.
> Bukti setiap butir: `PARITAS-LAYAR-DAN-AKSI.md` dan `RALAT-DEV-30-09-2026.md`.

| OQ | Pertanyaan | Bawaan sampai dijawab | Pemilik | Status |
| --- | --- | --- | --- | --- |
| OQ-MPNL-01 | Produk tetap JSON di dua tabel lama seperti Pega, atau dipindah ke tabel relasional baru (tiket 01)? Tabel baru membuat tiga view, dua prosedur, dan Claim Life tidak melihat produk baru, kecuali kedua bentuk ditulis bersamaan. | **JSON seperti Pega** (P1): `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE`, `JSONDATA` berkunci Pega, kolom datar `RIRISKID`/`RIRISK` *(ralat audit 02-10-2026: `PRODUCTNAME`/`BEGIN_DATE` dicabut lanjutan 1 L1 - tidak ada di katalog DEV)*; nol tabel baru, nol DDL | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — produk tetap **JSON seperti Pega**; tiket 01 tetap ditangguhkan |
| OQ-MPNL-02 | `ID` baris inward: Pega mengambilnya dari `M_PRODUCT_INWARD_LIFE_SEQ` (prosedur inward, `SaveProductName_Act` 16 b3052) dan menautkan lewat `PRODUCTID`, sedangkan Claim Life mencari inward `WHERE ID = <ID produk>` dan fakta bisnis menyebut kedua sisi berbagi identitas. DEV 196 vs 197 baris — keduanya pernah tidak sejalan. | `ID` inward = `ID` produk = `PRODUCTID`; `M_PRODUCT_INWARD_LIFE_SEQ` tidak dipakai; produk baru yang ID-nya sudah ada di salah satu tabel **ditolak** (bukan digandakan). Produk lama dibaca inward-nya lewat `PRODUCTID` (seperti `BrowseProductInward` b840) | work owner / DBA | ditutup 01-10-2026 (data DEV, lanjutan 1 L2) |
| OQ-MPNL-03 | Pemilih **R/I Rate** (`BrowseRateLifeSummary`, kelas `RATE_LIFE_SUMMARY`) dan **View Rate** (`BrowseRateLife_RD`, kelas `M_RATE_LIFE`) membaca view atas JSON rate — penjaga Claim Life `TestMasterViewTidakDisentuh` menuntut persetujuan (preseden OQ-MCRL-13). Boleh dibaca modul ini? | rute dan layar dibangun, menjawab **503 berkalimat**; `RIRATE`/`RIRATEID` baris plan yang sudah tersimpan dipertahankan. ⚠️ Akibat: baris plan **baru** tidak dapat diberi R/I Rate, dan `ProteksiPlanListLife` (*"RI/RATE tidak boleh kosong"*) menolak simpannya | work owner | ✅ **ditutup 01-10-2026 — K1 keputusan work owner 01-10-2026** (*"izinkan membaca view rate, baca-saja"*), dibangun `628c148` *(riwayat: terbuka — ✅ **ditutup 01-10-2026 — K1 keputusan work owner 01-10-2026** (*"izinkan membaca view rate, baca-saja"*): `Choose R/I Rate` membaca view `RATE_LIFE_SUMMARY`, `View Rate` membaca view `RATE_LIFE`, kolom RD saja; baris plan baru dapat diberi R/I Rate)* *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)* |
| OQ-MPNL-04 | Objek fisik master dibaca menurut konvensi nama kelas `ASM-FW-GISFW-Int-<X>`: `AGENT` (Ceding, SOB), `CLIENT` (Policy Holder), `CURRENCY`, `CAUSEOFLOSS_LIFE`, `RIRISK_LIFE_SUMMARY`, `PRODUCT_TYPE_LIFE` (plan). Tidak satu pun terbukti di katalog yang tercatat (inventori `[dugaan]`). | dibaca dengan nama itu; objek yang tidak ada = **503** yang menyebut objeknya (bukan daftar kosong). Nilai master diverifikasi hanya bila berubah dari yang tersimpan | DBA | ditutup 01-10-2026 (katalog DEV, L3) |
| OQ-MPNL-05 | Daftar pilihan bersumber `associated` (Property tidak ikut ekspor): `Product Name` b3620, `Premium Payment Method` b25611, `Birthday` b27960, `Document List` b15286. | `Product Name`, `Birthday`, `Document List`: isian teks bebas (nilai lama tampil apa adanya). `Premium Payment Method`: kode `1` Annual, `2` Semi Annual, `3` Quarterly, `4` Monthly (dari `GenerateUpload_Act` b332 `CARI37`); kode untuk *Single* tidak diketahui — nilai lain ditampilkan apa adanya | pemilik ekspor Pega | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — ⏳ **konfirmasi menyusul: pemilik ekspor Pega (daftar pilihan `associated`)** |
| OQ-MPNL-06 | `Generate` (`GenerateUpload_Act` 3 b1285): 39 judul (b1340) untuk 40 properti (b1343), dan kolom ke-12 `BENEFIT` memetakan `CARI2` — `pxConvertResultsToCSV` mengisi kolom menurut posisi, sehingga `BENEFIT`…`TREATYNUMBER` bergeser. | `SeeDetail.csv` dengan **39 judul VERBATIM**, tiap kolom berisi nilai yang **namanya** disebut judul (`LIENCLAUSE` = `CARI36`, tanpa judul, tidak ikut); ekspresi `@if` (`Basic/Rider`, `QS/SUPRLUS`, `KREDIT LIFE/…/HEALTH`, `Annual/…/Single`) ditiru | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MPNL-07 | `Download All` b67657 menjalankan jalur treaty-in (`M_ATTACHMENTTREATY_2` milik `TreatyIn.ID`) — salah ekspor (OQ-056). | mengunduh seluruh lampiran **produk ini** sebagai satu arsip zip (AC tiket 08) | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MPNL-08 | Kolom datar `PRODUCTNAME` dan `BEGIN_DATE` `M_PRODUCT_LIFE`: korpus nol penulis (`SaveProductNameLIfeFlat` hanya `RIRISKID`/`RIRISK`), prosedur hanya `JSONDATA` (`dba-procedures-and-ddl.md` §1). | **dicabut (lanjutan 1 L1; ralat audit 02-10-2026): kedua kolom tidak ada di `M_PRODUCT_LIFE` DEV, tidak ditulis.** Bawaan lama: ditulis menurut P1: `PRODUCTNAME` = JSON `PRODUCTNAME`, `BEGIN_DATE` = `BEGIN` inward (`dd/MM/yyyy` → `DATE`), NULL bila kosong | work owner / DBA | ditutup 01-10-2026 (katalog DEV, L1/L4) |
| OQ-MPNL-09 | `OutwardList`: Pega mengisinya saat checkbox `On Retention` diubah (nilai `BEGIN`/`MATURE` saat itu). `OVR_COMM` yang dibaca view tidak punya penulis di korpus (`SetParamOutward` tidak diekspor). | dihitung **saat simpan** bila checkbox diubah dalam sesi sunting itu (memakai `BEGIN`/`MATURE` yang disimpan), selain itu dipertahankan; `OVR_COMM` `""` | work owner | ditutup 01-10-2026 (data DEV, L5) |
| OQ-MPNL-10 | Unggah lampiran: `InsertGoogleStorage_Act` mengirim berkas ke layanan luar (`ServiceGoogle` + `LinkService`) dan mencatat `URLPUBLIC`. Alamatnya di `M_LINK_SERVICE` (OQ-047). | **stub outbox**: rekam `M_ATTACHMENTPRODUCTNAME` + antre efek di satu transaksi; pelaksana stub menyimpan berkas di folder lokal `UNGGAHAN_DIR` dan mencatat `T_STORAGE_IMAGE` (`URLPUBLIC` kosong, `APPFOLDER = Contract`, `STORAGE = standard`). Gagal → status `gagal`, dapat diulang | work owner / tim inti | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** *(Ralat 03-10-2026: dibalik — keputusan work owner "untuk document masih belum berfungsi, ikuti dari XML nya aja"; penyimpanan nyata seperti XML di balik `PELAKSANA_STORAGE=nyata`, bab bertanggal 03-10-2026 di bawah.)* |
| OQ-MPNL-11 | `View Office Online` membungkus URL bertanda tangan ke penampil kantor di luar (`DownloadAttProdName_Act` 7 b1080). | tautan tampil untuk jenis `xls/xlsx/doc/docx/ppt/pptx` (b69291); rute menjawab **503 berkalimat**: penampil luar tidak dipanggil, alamatnya tidak ditulis | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** *(Ralat 03-10-2026: dibalik — keputusan work owner "izinkan ditulis di kode"; penampil dibuka seperti XML b1103, bab 03-10-2026 sore di bawah.)* |
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
| OQ-MPNL-03 | view `RATE_LIFE_SUMMARY` dan `RATE_LIFE` ada (`M_RATE_LIFE_SUMMARY` 339 baris, `M_RATE_LIFE` 96.038) | **tetap terbuka** — register ini tidak memuat izin bertanggal work owner; `R/I Rate`/`View Rate` tetap 503 *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: izin tercatat; ditutup — DEV baca-saja: `ri-rate` 200, 346 baris; `rate` 200, 59 baris untuk satu RIRATEID)* *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)* |
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

## Keputusan work owner 01-10-2026 atas laporan K1 — penyimpangan sadar

Work owner 01-10-2026: *"ikuti rekomendasi semua"* atas laporan lanjutan keputusan OQ (commit `d0c7b0a`, `628c148`, `a64828f`, `9af0be9`).

| OQ | Hal | Keputusan | Status |
| --- | --- | --- | --- |
| OQ-MPNL-16 | `View Rate` menyaring dengan `RIRATEID` baris plan. Di XML grid menyaring `ParamID.OUTWARDRATEID` (`ViewRate` b1024) yang tidak pernah diisi rule mana pun, sehingga di Pega dialog itu selalu kosong. | **dipertahankan** — meniru dialog yang selalu kosong tidak berguna | ditutup 01-10-2026 |
| OQ-MPNL-17 | Server menolak `RIRATEID` pilihan baru yang tidak ada di view `RATE_LIFE_SUMMARY`. Pega tidak memeriksanya. | **dipertahankan** — pemilih hanya menawarkan rate yang ada; pemeriksaan mencegah data rusak | ditutup 01-10-2026 *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)* |

## Keputusan work owner 02-10-2026 — pindah ke tabel flat (`PROMPT-PINDAH-FLAT-MASTER-PRODUCT-NAME-LIFE.md`)

**OQ-MPNL-01 dibuka ulang dan ditutup 02-10-2026: FLAT** (K5) — menggantikan penutupan 01-10-2026 (*"produk tetap **JSON seperti
Pega**; tiket 01 tetap ditangguhkan"*, dikutip). Keputusan grilling D2 dan Q1b berlaku lagi; rinciannya tiket 01 bab bertanggal 02-10-2026.

| OQ | Pertanyaan | Bawaan / keputusan | Pemilik | Status |
| --- | --- | --- | --- | --- |
| OQ-FLAT-01 | FK `M_ATTACHMENTPRODUCTNAME.TREATYID` → `M_PRODUCTNAME_LIFE.ID` ditambahkan? Tabel lampiran bukan milik migrasi ini | tidak | work owner | ✅ ditutup 02-10-2026 — keputusan work owner: *"ikuti rekomendasi"* |
| OQ-FLAT-02 | Ekspor flat → JSON untuk jalur mundur **sesudah** ada tulisan baru dibangun? | tidak (jalur mundur hanya sebelum tulisan baru: `-migrate-down` + versi aplikasi lama) | work owner | ✅ ditutup 02-10-2026 — keputusan work owner: *"ikuti rekomendasi"* |
| OQ-FLAT-03 | Siapa dan kapan menghentikan penulisan Pega ke layar Product Name Life (syarat peralihan)? *(Diperjelas 02-10-2026: yang dimaksud layar Pega LAMA — menu MASTER → **Master Product Name Life** (section `InboxProductName`). Tombol `Save`-nya (`SaveProductName_Act` langkah 9, 10, 15 — ketiganya hidup, `pyStepsBlockName` kosong) menulis `M_PRODUCT_LIFE.JSONDATA` + `RIRISKID`/`RIRISK` dan `M_PRODUCTINWARD_LIFE.JSONDATA`. Satu-satunya penulis lain di korpus, `SaveInwardProductName_Act`, hanya dicapai tombol `Inward` b75368 ber-vis `1=2` (mati). Sesudah peralihan aplikasi baru membaca tabel flat saja: simpan di layar Pega itu tidak terlihat di aplikasi baru, dan produk yang sudah dipindah lalu diubah di Pega membuat alat pindah menolak putaran ulang.)* | — | work owner | terbuka |
| OQ-FLAT-04 | Claim Life `ambangproduk.go` membaca view `PRODUCTINWARD_LIFE` atas `M_PRODUCTINWARD_LIFE.JSONDATA`; sesudah peralihan tabel JSON tidak diperbarui lagi dan view **tidak** dibangun ulang (K7). Claim Life dialihkan membaca `M_PRODUCTNAME_LIFE` (brief Claim Life), atau view dibangun ulang oleh DBA? | belum ada — di luar folder modul ini | work owner / pemilik Claim Life | ⏸️ ditunda 02-10-2026 — keputusan work owner: *"OQ-FLAT-04 :  BIARKAN SAJA, NANTI PAS DEVELOP BAGIAN ITU AKAN DIPEERBAIKI!"* (diperbaiki saat bagian Claim Life itu dikembangkan; modul ini tidak menyentuh Claim Life) |
| OQ-FLAT-07 | Uji kering DEV 02-10-2026 (`LAPORAN-MIGRASI-FLAT.md`): 2 nilai `RICOMM` berubah teks bila dipindah ke NUMBER(38,8) — 1 berkoma desimal, 1 bernol depan; nilai angkanya sama. Diterima dalam bentuk kanonik? | **terima** (rekomendasi): services sendiri sudah menyimpan koma sebagai titik saat simpan; nilai uang tidak berubah. Jalankan `-jalankan -terima-normalisasi="koma desimal,nol depan"` (per JENIS: menerima satu jenis tidak menerima jenis lain) | work owner | ✅ ditutup 02-10-2026 — keputusan work owner: *"ikuti rekomendasi"* |
| OQ-FLAT-08 | Uji kering DEV 02-10-2026 18:50: 189 objek `OutwardList` yang dianggap K4 "kosong" (keenam kunci OR kosong) berisi `OUTWARDNAME`/`OUTWARDNAMEID` (189) dan `OUTWARDRATE`/`OUTWARDRATEID` (187) - 4 nama dan 3 rate berbeda; nol di kedua objek OR. Premis K4 keliru. Dipindah atau dibuang? | **pindahkan**: empat kolom baru di `M_PRODUCTNAME_LIFE_OUTWARD` (migrasi 146, belum dijalankan di DEV); K4 hanya objek yang SEMUA kuncinya kosong. Sampai diputuskan, alat menolak `-jalankan` (752 nilai tanpa kolom) | work owner | ✅ ditutup 02-10-2026 — keputusan work owner: *"ikuti rekomendasi"* |
| OQ-FLAT-09 | `MATURE` produk 100175 (K3 "bukan tanggal") ternyata tanggal dalam bentuk lain (bukan `dd/MM/yyyy`) - dapat diselamatkan. Dikonversi atau di-NULL-kan? | **konversi**; alat menolak `-jalankan` sampai diputuskan | work owner | ✅ ditutup 02-10-2026 — keputusan work owner: *"ikuti rekomendasi"* |
| OQ-FLAT-05 | Tipe angka: rancangan `NUMBER` tanpa presisi dan `NUMBER(1/3/4)` ditolak penjaga inti `TestNolNumberTanpaPresisi` | **patuhi penjaga**: `NUMBER(38,8)` / `NUMBER(5)` (K6) | work owner | ✅ ditutup 02-10-2026 |
| OQ-FLAT-06 | View `CREATE OR REPLACE VIEW` di migrasi ditolak penjaga inti `TestSeluruhCreateDapatDibacaNamanya` | *"tidak ada table view yang dipake, semua simpan dan baca dari table flat"* — view tidak dibangun ulang (K7) | work owner | ✅ ditutup 02-10-2026 |

### Keputusan work owner 02-10-2026 malam — *"ikuti rekomendasi"*

| OQ | Keputusan | Dibangun |
| --- | --- | --- |
| OQ-FLAT-01 | FK lampiran → induk flat **tidak** ditambahkan | — |
| OQ-FLAT-02 | ekspor flat → JSON **tidak** dibangun | — |
| OQ-FLAT-07 | 2 normalisasi `RICOMM` (koma desimal, nol depan) **diterima** | `-terima-normalisasi="koma desimal,nol depan"` di panduan langkah 4 |
| OQ-FLAT-08 | 189 objek outward bukan-OR **dipindah** | empat kolom `OUTWARDNAMEID`, `OUTWARDNAME`, `OUTWARDRATEID`, `OUTWARDRATE` di migrasi 146; K4 hanya objek yang SEMUA kuncinya kosong *(Ralat 02-10-2026 malam: keempat kolom itu ditambah migrasi BARU `148_m_productname_life_outward_kolom` (`ALTER TABLE ... ADD`), bukan di 146 - 146 ternyata sudah dijalankan di DEV pukul 15:39 (`T_MIGRASI`), sehingga isinya dikembalikan ke bentuk yang dijalankan.)* |
| OQ-FLAT-09 | tanggal inward berbentuk lain **dikonversi** (dicatat per produk dan kolom) | alat pindah; `MATURE` produk 100175 (`dd-MM-yyyy`) |

Masih terbuka: **OQ-FLAT-03** (siapa dan kapan Pega berhenti menulis) dan **OQ-FLAT-04** (pembaca Claim Life atas view
`PRODUCTINWARD_LIFE`) — keduanya tanpa rekomendasi tunggal. *(Ralat 02-10-2026: OQ-FLAT-04 **ditunda** — keputusan work owner
*"OQ-FLAT-04 :  BIARKAN SAJA, NANTI PAS DEVELOP BAGIAN ITU AKAN DIPEERBAIKI!"*; yang masih terbuka tinggal OQ-FLAT-03.)*

### OQ-MPNL-05 sebagian terjawab 03-10-2026 - daftar `Document List`

Work owner mengirim XML rule `Rule-Obj-Property` `.Document` kelas `ASM-FW-GISFW-Data-UnderwritingLimit` (ruleset GISFW
01-01-91, `pyTableOption` PromptList, 16 nilai) dengan kalimat *"untuk document list productname pilih dari list ini"*.
`Document List` grid DOCUMENT CLAIM kini dipilih dari ke-16 nilai itu (`PILIHAN_DOKUMEN_KLAIM`, `bentuk.ts`); nilai lama di
luar daftar tetap tampil bertanda *(not in the reference list)* dan tidak dibuang. XML-nya tidak disimpan di repo (memuat nama
operator; korpus hanya-baca). `Product Name` dan `Birthday` tetap terbuka.

Data DEV 03-10-2026 (SELECT saja, 195 produk lama): 1.121 baris DOCUMENT CLAIM - 36 persis di daftar; 1.085 di luar daftar
(16 nilai): 10 nilai / 425 baris sama dengan teks Indonesia entri daftar tanpa terjemahan `(English)`, 6 nilai / 660 baris
tidak cocok entri mana pun. **Terbuka:** dibiarkan apa adanya, atau dikonversi ke teks daftar saat disalin (Copy Old / alat
pindah) - menunggu keputusan work owner.


### OQ-MPNL-10 dibalik 03-10-2026 — lampiran ke penyimpanan nyata seperti XML

Work owner 03-10-2026, sesudah dijelaskan bahwa berkas lampiran masih tersimpan di folder stub: *"untuk document masih belum
berfungsi"*, lalu *"untuk document masih belum berfungsi, ikuti dari XML nya aja"*.

Bukti DEV (SELECT saja, 03-10-2026): 2 lampiran produk yang ada di tabel flat adalah unggahan Pega 22-08-2025 —
`T_STORAGE_IMAGE` berisi `URLPUBLIC` https, `APPFOLDER` = awalan gs + App + `Contract/Doc/2025/08/` + `FILENAME`, `EXPDATE`
sudah lewat. Di stub, berkas itu tidak dapat diunduh (409 "not in the storage stub"). `M_LINK_SERVICE` memuat kunci `Google` /
`upload`, `geturl`, `delete`; `T_FOLDER_IMAGE` 1 baris; `GCP_IMAGE` ada.

| Aksi | XML diikuti | Dibangun (`backend/services/mpnl_storage.go`) |
| --- | --- | --- |
| Add attachment → Attach | `InsertGoogleStorage_Act` (Set Data b1339, SET JSON b1572, `ServiceGoogle` b1846, `Insert_T_Storage_SQL`) | POST JSON `App`, `Kodestring`, `Durasi` 1800, `Folder`, `Namafile`, `Image` base64, `ext`, `MimeType`; `T_STORAGE_IMAGE` dari jawaban (`URLImage`, `appfolder`, `exp`) |
| Tautan nama berkas, Download All | `DownloadAttProdName_Act` 6 b953 / `DownloadAll_Act` 5 b875 → `GetUrlGoogleStorage_Act` | URL tersimpan dipakai selama `EXPDATE` belum lewat; selain itu `geturl` (Folder = APPFOLDER tanpa Namafile tanpa awalan gs+App) lalu `Update_T_Storage_SQL` |
| Delete | `DeleteAttacProdName_act` 2 b411 → `DeleteGoogleStorage_Act` | POST `delete` (`Namafile` = jalur objek penuh); gagal = rekam tetap |

Cara menyalakan (di `.env` server aplikasi; bawaan tetap stub):

```
PELAKSANA_STORAGE=nyata
STORAGE_TOKEN_SALT=<garam procedure GET_TOKEN_STORAGE — dari DBA>
```

`UNGGAHAN_DIR` tetap wajib (antrean lokal sampai objek tercatat). Saklar ini bersama Treaty Contract Out (`inti/backend/config`).

Penyimpangan sadar: (1) token ditiru `inti/backend/layanan` (keputusan tim "jangan ada lagi pemanggilan procedure"), bukan
procedure `GET_TOKEN_STORAGE`; (2) pengiriman tetap efek keluar outbox (P5) — rekam dulu, kirim sesudahnya, dapat diulang;
(3) isi diunduh backend dari URL bertanda tangan (hanya https, tanpa pengalihan) lalu diteruskan ke peramban — Pega membuka URL
itu di jendela peramban; (4) jawaban geturl tanpa `appfolder` tidak mengosongkan `APPFOLDER`; (5) objek yang dicatat stub
(`URLPUBLIC` kosong) tetap dibaca dan dihapus di folder stub. `View Office Online` tetap 503 (OQ-MPNL-11): penampilnya
alamat literal di XML (b1103), dan alamat literal dilarang `TestMPNLNolAlamatLayanan`.

Uji: layanan tiruan `httptest` TLS (`mpnl_storage_test.go`), status HTTP 502/503/409 (`rute_storage_test.go`). Layanan
sungguhan **tidak** dipanggil dari uji atau dari sesi ini; data DEV hanya dibaca.

Sesudah tinjauan kode 03-10-2026: Namafile dari waktu REKAM (ID lampiran) sehingga kirim ulang menimpa objek yang sama;
efek kirim tidak ikut batal bila peramban memutus; token tidak dipakai ulang bila sisa < 15 detik; pengalihan 3xx tidak
diikuti; unduhan yang putus diputus sambungannya (bukan 200 terpotong); stub menolak menghapus objek yang ada di
penyimpanan nyata (503). **Batasan diketahui:** menghapus lampiran yang SEDANG dikirim dari tab/sesi lain dapat
meninggalkan objek tanpa rekam — jendelanya selama unggah berjalan; tidak ditangani.

### Lanjutan 03-10-2026 sore — "MASIH GA BISA": selalu nyata dan View Office Online

Work owner melaporkan dua pesan: `View Office Online is a stub …` dan `the attachment file is not in the storage stub …
until it is connected (OQ-MPNL-10)`. Sebabnya: backend yang berjalan (`go run`, mulai 10.33) masih kode sebelum
`af853685`, dan `.env` tidak memuat `PELAKSANA_STORAGE=nyata` / `STORAGE_TOKEN_SALT`. Data DEV (SELECT saja):
`GCP_IMAGE` memuat token yang masih berlaku ±58 hari; tidak ada job penyegar; `M_LINK_SERVICE` tidak punya alamat penampil.

Keputusan work owner 03-10-2026:

| Butir | Keputusan | Dibangun |
| --- | --- | --- |
| Aktivasi | *"Selalu nyata, ikut XML"* | Layanan Oracle SELALU memakai `penyimpananGoogle` — saklar `PELAKSANA_STORAGE` tidak lagi dibaca modul ini (bab di atas, "Cara menyalakan", **gugur**). Token berlaku di `GCP_IMAGE` (`INPUTDATE > SYSDATE`, seperti procedure; *ralat sore: margin 15 detik dicabut*) dipakai ulang TANPA garam; `STORAGE_TOKEN_SALT` hanya untuk token baru — tanpa token berlaku dan tanpa garam: 503 berkalimat. `UNGGAHAN_DIR` tetap wajib (antrean). Teks entri penjaga inti disesuaikan. |
| OQ-MPNL-11 `View Office Online` | *"Izinkan ditulis di kode"* | rute `GET …/lampiran/{lid}/office` menjawab `{url}` bertanda tangan (`GetUrlGoogleStorage_Act`, Durasi 1800); frontend `penampilOffice.ts` memuat alamat penampil b1103 — satu-satunya alamat literal modul, pengecualian bernama `alamatDiizinkan` di `TestMPNLNolAlamatLayanan`. |
| Cara membuka penampil | *"KENAPA HARUS NYENGGOL MODUL LAIN?"* — tidak menyentuh modul lain | penjaga lintas-modul `modul/claimlife/frontend/unduhdokumen.test.ts` melarang `window.open` / `location.*` (navigasi ke backend tanpa identitas). Penampil dibuka lewat FORM GET ber-`action` tetap ke alamat penampil, `target="_blank"`, input `src` = URL bertanda tangan (`application/x-www-form-urlencoded` = `@encodeURL`); dua langkah — klik mengambil URL, tombol di jendela kecil membuka tab (klik pengguna, tidak diblokir pemblokir pop-up). *(Ralat 03-10-2026: work owner "TIDAK UDAH TAMBAH TAMBAH LAGI BUAT DOCUMENT, PAKE APA YANG SUDAH ADA" — jendela kecil dan tombolnya dicabut; link `View Office Online` yang ada langsung mengirim form tersembunyi, satu klik.)* |

Permintaan work owner 03-10-2026 sesudahnya: *"SELAIN OFFICE, YANG BISA DIBUKA SECARA ONLINE SEPERTI PDF IMAGE BUATIN VIEW
AJA DARI POPUP ATAU WINDOWS BARU (BUKAN TAB BARU); YANG OFFICE JUGA SAMAIN"*. Dibangun `[tidak ada di korpus]`: link `View`
(label `GRID_MPNL.view` b74798) untuk pdf dan gambar raster (png, jpg, jpeg, gif, bmp, webp — svg/html tidak: objek URL
mewarisi asal aplikasi); isinya dari rute unduh yang ada (berheader identitas), ditampilkan di POPUP (`Modal` penuh) lewat
objek URL lokal. `View Office Online` kini juga di popup: form GET ke bingkai bernama di dalam popup, bukan tab baru.
Jendela peramban terpisah tidak dibangun: menuntut pembukaan jendela lewat skrip, yang dilarang penjaga lintas-modul
`unduhdokumen.test.ts` (modul lain tidak disentuh). Backend tidak berubah.

Permintaan work owner 03-10-2026: *"PERBAIKI UPLOAD DOCUMENT BISA BANYAK DAN BISA DRAG AND DROP"*. `Add attachment` kini
menerima banyak berkas (pemilih `multiple` dan kotak seret-lepas); berkas diunggah SATU PER SATU lewat rute unggah yang ada
(`POST …/lampiran`, satu berkas per permintaan - backend tidak berubah). Nama ganda dalam pilihan dibuang (backend menolak
nama ganda per produk); kegagalan per berkas ditampilkan dan berkasnya tinggal di pilihan, yang berhasil langsung ke grid.

### Unduh lampiran 502 (03-10-2026 sore) — alur `GetUrlGoogleStorage_Act` apa adanya

Gejala: `…/lampiran/20250822135941817/unduh` 502. Diagnosis baca-saja: URL tersimpan ber-EXPDATE `13:28:00` ditolak
Google `400 ExpiredToken` sesudah 13:28 WIB, sedangkan kode membaca EXPDATE sebagai GMT (dianggap berlaku sampai 20:28 WIB).

Perbaikan pertama (`0ab5e116`) menambah penyesuaian Go: margin 1 menit dan coba ulang bila URL ditolak. Work owner:
*"ITU HARUSNYA IKUTI ACTIVITY DARI XML NYA SEMUA, KAMU HANYA MELAKUKAN PENYESUAIAN KE GO!! LIHAT ACTIVITY VIEW, ADA
GETTOKEN, DAN KETIKA EXPIRED, HIT API LAGI UNTUK GENERATE LINK BARU!"* — penyesuaian itu dicabut; `Tautan`
(`backend/services/mpnl_storage.go`) kini mengikuti langkah XML apa adanya:

| Langkah | XML | Go |
| --- | --- | --- |
| 4 b671 | GET LINK `GetLinkStorage_SQL` | `AmbilObjek` (pemanggil) |
| 5 b855 | `Param.Url` = URLImage tersimpan | URL tersimpan |
| 6 b1022 | JIKA EXPDATE SUDAH EXPIRED — dilewati bila `@CompareDates(exp, @CurrentDateTime())` (b2610); `exp == ""` (b705) = dijalankan | `berlaku`: EXPDATE sesudah sekarang, tanpa margin |
| 6.1 b1056 | GET TOKEN `GetTokenStorage_SQL` | `tokenStorage` (`INPUTDATE > SYSDATE`, tanpa margin) |
| 6.2–6.5 | Set Data, SET JSON, GetLinkService Google/geturl, Connect-REST | `panggil` |
| 6.6 b1951 | Insert ke table — dilewati bila `Response.URLImage == ""` (b2558); UpdateDoc apa adanya → `Update_T_Storage_SQL` | URL tersimpan bila kosong; selain itu objek baru (appfolder, DateTime apa adanya) dicatat |
| 7 b2650 | Return `Param.Url` | URL |

Satu-satunya tafsir: EXPDATE dibaca jam Asia/Jakarta — jam yang Pega pakai saat membandingkan DATE basis data dengan
`@CurrentDateTime()` (bukti di atas). Verifikasi DEV baca-saja: enam lampiran produk 100003 terunduh — empat lewat
langkah 6 (token + geturl), dua memakai URL tersimpan.
