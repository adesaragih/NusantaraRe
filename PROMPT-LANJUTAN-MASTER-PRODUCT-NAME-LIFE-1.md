# PROMPT — LANJUTAN 1 **MASTER PRODUCT NAME LIFE**: perbaikan kolom datar yang tidak ada di DEV, OQ yang ditutup data DEV, penutup modul *(sesi baru, folder `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM`, cabang **`dev`**)*

> Permintaan work owner 01-10-2026: langkah berikutnya **untuk modul Master Product Name Life** sesudah laporan paket 0–11 + review *(`71c35b1` →
> `c47c6c9`)*. XML patokan dasar; tiket = hasil grilling; tiket dan dokumen yang salah **diperbaiki**. Jangan skip; baca ulang XML.
>
> Folder boleh disunting: **`APP_RNM/modul/masterproductnamelife/**`** saja *(+ berkas bangkitan bila berubah)*. Commit `git commit -o -- <jalur>`.
> **Nol `git push`, nol `git pull`, nol `-migrate`, nol penulisan ke DEV.** Setiap pembacaan activity mencetak `pyStepsBlockName`; nomor baris hanya dari
> `sed -e 's/></>\n</g' <berkas> | grep -n …` *(tulis perintahnya di kepala PARITAS)*.

## 0. VERIFIKASI LAPORAN PAKET 0–11 *(asisten, 01-10-2026)*

| Klaim | Hasil cek |
| --- | --- |
| 13 commit, di luar folder modul hanya `modul_masterproductnamelife_gen.go` | ✅ |
| Migrasi hanya slot 960, `UPDATE` saja | ✅ |
| ID `'1'‖LPAD(seq,5)` | ✅ sesuai `dba-procedures-and-ddl.md` §1; kolom `ID` di DEV `VARCHAR2(6)` |
| Tombol `Inward` mati | ✅ showHarness `InwardProductName`, sel `pyVisible OTHER` + `pyCondition 1=2`; nomor baris laporan bergeser dari perintah standar |
| Uji modul | ✅ lulus di tiruan; 9 uji `db` SKIP |
| **Kolom datar `PRODUCTNAME`, `BEGIN_DATE`** | ⛔ **TIDAK ADA di DEV** — lihat L1 |

## 1. KATALOG DEV *(asisten, `ALL_TAB_COLUMNS`, baca-saja)*

| Tabel | Kolom *(urut `COLUMN_ID`)* |
| --- | --- |
| `M_PRODUCT_LIFE` | `ID` VARCHAR2(6) · `JSONDATA` CLOB · `RIRISKID` VARCHAR2(10) · `RIRISK` VARCHAR2(100) — **hanya empat** |
| `M_PRODUCTINWARD_LIFE` | `ID` VARCHAR2(6) · `JSONDATA` CLOB |

## 2. PERBAIKAN DAN OQ YANG DITUTUP DATA DEV

| # | Hal | Bukti DEV | Kerjakan |
| ---: | --- | --- | --- |
| **L1** ⛔ | kode menulis `PRODUCTNAME` dan `BEGIN_DATE` *(`repository/mpnl_tabel.go:37` `KolomProduk`, `repository/mpnl_tulis.go:30` INSERT, `:35-36` UPDATE, `tiruan/gudang.go:41`)* | kolom itu **tidak ada** — setiap simpan di DEV akan gagal `ORA-00904`. Penulis datar di korpus, `SaveProductNameLIfeFlat`, memang hanya `RIRISKID`/`RIRISK` *(OQ-MPNL-08)* | buang keduanya dari INSERT, UPDATE, daftar kolom, tiruan, dan uji; tulis hanya `RIRISKID`, `RIRISK` seperti XML. Uji baru **katalog**: setiap kolom yang ditulis repository ⊆ daftar kolom §1 *(disimpan sebagai data uji di folder modul)*; uji panjang: `RIRISKID` > 10 atau `RIRISK` > 100 karakter = galat berkata-kata sebelum SQL. Uji gigit: menambah `PRODUCTNAME` ke INSERT = merah |
| L2 | **OQ-MPNL-02** identitas inward | 196 dari 196 produk punya baris inward ber-`ID` sama; dari 199 baris inward, 197 ber-`JSONDATA.PRODUCTID` = `ID`, 196 ber-`ID` yang ada di produk | **ditutup**: inward memakai **`ID` yang sama dengan produk**, `PRODUCTID` = `ID` *(kode sudah begitu, `mpnl_identitas.go:7`)*. Catat buktinya; uji: produk baru → inward ber-`ID` sama |
| L3 | **OQ-MPNL-04** objek master | ada di DEV: tabel `AGENT`, `CLIENT`; view `CURRENCY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE`, `RIRISK_LIFE_SUMMARY` *(di atas `M_CURRENCY`, `M_CAUSEOFLOSS_LIFE`, `M_PRODUCT_TYPE_LIFE`, `M_RIRISK_LIFE_SUMMARY`)* | **ditutup**: pastikan setiap pemilih membaca objek bernama persis itu; label `[dugaan]` di inventori dicabut |
| L4 | **OQ-MPNL-08** kolom datar | lihat L1 | **ditutup** oleh L1 |
| L5 | **OQ-MPNL-09** `OVR_COMM` | 153 produk punya `OutwardList[0]`, **0** ber-`OVR_COMM` terisi | **ditutup**: `OVR_COMM` ditulis kosong, seperti Pega |
| L6 | **OQ-MPNL-12** `Medical` | lebih dari 60 nilai berbeda di `UnderwritingLimitList[*].Medical` *(NM, FCL, A–H, kombinasi ME/MU/ECG, dll.)* | **ditutup**: **teks bebas** *(R16 benar)*; **tiket 06 diperbaiki** — tuntutan hanya `FCL`/`NM`/`MEDIS` keliru |
| L7 | **OQ-MPNL-15** `TREATYCONTRACTID` | kolom bernama itu tidak ada di `TREATYCONTRACT_LIFE` maupun `TREATYYEAR_LIFE` | **ditutup**: tetap kosong, seperti Pega |
| L8 | **OQ-MPNL-03** R/I Rate dan View Rate | view `RATE_LIFE_SUMMARY` dan `RATE_LIFE` ada *(di atas `M_RATE_LIFE_SUMMARY` 339 baris, `M_RATE_LIFE` 96.038)* | bila `OQ-MASTER-PRODUCT-NAME-LIFE.md` mencatat **izin work owner bertanggal**: baca keduanya baca-saja, kolom yang dibaca RD XML saja *(`BrowseRateLifeSummary`, `BrowseRateLife_RD`)*, nol `SELECT *`; rute yang kini 503 menjawab 200; baris PLAN LIST baru dapat diberi R/I Rate. Tanpa izin: tetap 503 |
| L9 | nomor `bNNN` | bergeser | hitung ulang di `PARITAS-LAYAR-DAN-AKSI.md` dan `RALAT-DEV-30-09-2026.md` dengan perintah standar |

**Dokumen yang diperbaiki** *(catatan bertanggal 01-10-2026, kalimat lama dikutip)*: `dba-procedures-and-ddl.md` *(empat kolom datar → dua)*, `spec.md` §
yang menyebut `PRODUCTNAME`/`BEGIN_DATE`, `issues/01`, `issues/06` *(Medical)*, `RALAT-DEV-30-09-2026.md` *(P1 diralat)*, `STRUKTUR-TABEL-MASTER-PRODUCT-NAME-LIFE.md`,
dan register OQ *(02, 04, 08, 09, 12, 15 ditutup; 03 sesuai izin)*. Brief asisten `PROMPT-IMPLEMENTASI-MODUL-MASTER-PRODUCT-NAME-LIFE.md` bab 4 P1 ikut
menyebut empat kolom — **itu keliru**, jangan diikuti.

## 3. CEK PERAMBAN BACA-SAJA

Sesudah L1: backend dengan env DEV *(dibaca saja)* + `npm run dev` *(port **5174**)*. Buka **View** di baris grid dan **Edit** untuk produk DEV yang ada,
panel lampiran, dan *(bila L8 diizinkan)* **View Rate** — empat tombol yang laporan sebelumnya tidak dapat amati. **Jangan tekan Save** *(nol penulisan
ke DEV)*. Catat hasilnya di `docs/LAPORAN-UJI-MANUAL.md`: tombol, layar yang terbuka, nilai yang tampil sesuai `JSONDATA` *(tanpa data orang)*.

## 4. PAKET — satu commit per paket

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | L1 kode + uji katalog + uji panjang | `master-product-name-life: tulis hanya RIRISKID dan RIRISK - PRODUCTNAME dan BEGIN_DATE tidak ada di DEV` |
| 2 | L2, L3, L5, L7 uji + kode bila perlu | `master-product-name-life: identitas inward, objek master, OVR_COMM, TREATYCONTRACTID dari data DEV` |
| 3 | L8 bila diizinkan | `master-product-name-life: R/I Rate dan View Rate dari view rate (OQ-MPNL-03)` |
| 4 | cek peramban §3 | `master-product-name-life: uji manual baca-saja` |
| 5 | dokumen §2 + L6 + L9 | `docs: master-product-name-life lanjutan 1` |

Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`, `db` dengan `-p 1`)*, `gofmt`, `tsc`, `vitest`, `vite build`.

## 5. LAPORAN

Satu pesan: tabel **paket → commit → OQ/temuan → bukti → angka uji**; daftar SQL yang berubah di L1; isi `LAPORAN-UJI-MANUAL.md`; OQ yang masih terbuka
beserta pemiliknya; persentase kesiapan modul dan total migrasi; bab **TELEMETRI EKSEKUSI**.

## 6. OQ UNTUK WORK OWNER *(sisa sesudah lanjutan ini)*

| OQ | Pertanyaan | Bawaan |
| --- | --- | --- |
| OQ-MPNL-03 | Boleh membaca view `RATE_LIFE_SUMMARY` dan `RATE_LIFE`, baca-saja? Izin yang sama membuka OQ-MCRL-13 Retro Life | tetap 503 |
| OQ-MPNL-07 | `Download All` di XML menjalankan jalur lampiran Treaty In. Dibuat mengunduh lampiran produk ini? | unduh lampiran produk ini |
| OQ-MPNL-13 | Salinan produk (`Copy`) ber-pembuat pelaku, bukan pembuat produk asal? | pelaku |
| OQ-MPNL-14 | Setiap simpan menambah satu baris komentar walau komentarnya kosong, seperti XML? | ikut XML |
| OQ-MPNL-05, 06, 10, 11 | daftar pilihan `associated`, judul CSV 39/40, penyimpanan berkas nyata, penampil Office | bawaan laporan sebelumnya |

---

*Disusun 1 Oktober 2026 dari katalog DEV (`ALL_TAB_COLUMNS`, `ALL_OBJECTS`), agregat `JSONDATA` kedua tabel produk (identitas inward, `OutwardList`,
`UnderwritingLimitList.Medical`), kode `mpnl_tabel.go`, `mpnl_tulis.go`, `mpnl_identitas.go`, dan register OQ modul.*
