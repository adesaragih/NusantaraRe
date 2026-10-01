# PROMPT — IMPLEMENTASI MODUL **MASTER PRODUCT NAME LIFE** *(sesi baru, folder `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM`, cabang **`dev`**)*

> Permintaan work owner 01-10-2026: *"aku mau menjalankan implement master product name life"*. Satu sesi, satu modul, satu folder.
> Isinya diambil dari bab 3 `PROMPT-IMPLEMENTASI-TIGA-MODUL-RUMPUN-LIFE.md` *(yang ditahan)* dan dilengkapi daftar rule korpus.
>
> **XML adalah patokan dasar. Tiket adalah hasil grilling.** Tiket yang tidak sesuai XML atau katalog DEV **diperbaiki di dokumennya** *(catatan
> bertanggal, kalimat lama dikutip, bukti disebut)*. Nol menu, layar, tombol, atau aksi tanpa bukti XML.

## 0. ATURAN BACA XML

1. **Setiap pembacaan activity mencetak `pyStepsBlockName`**. `//` = langkah ter-remark, **tidak pernah jalan** — jangan dibangun.
2. Pecah XML dengan `sed -e 's/></>\n</g'`; kutip nomor baris `bNNN` di setiap bukti.
3. WHEN dibaca dari urutan aksi. Precondition `WhenTrue 2` = lanjut, `3` = lewati.
4. Cari **penulis** sebuah kolom, bukan hanya pembacanya.
5. Setiap tombol = satu tombol di section XML *(nama section + `bNNN`)* yang memanggil activity yang sama. Label **VERBATIM** dari XML.
6. `grilling-ronde-*` tersegel, tidak disunting.

## 1. BATAS KERJA

| Hal | Isi |
| --- | --- |
| Folder repo | **`D:\XML\RNM_BRD\OUTPUT_HASIL_RNM`** — **bukan** klon lain seperti `D:\RNM SOURCE TREE` |
| Cabang | **`dev`**. Commit dengan jalur eksplisit `git commit -o -- <jalur>`, pesan berawalan `master-product-name-life:`. **Nol `git push`, nol `git pull`, nol `git pull --rebase`** |
| Folder boleh disunting | **`APP_RNM/modul/masterproductnamelife/**`** + `inti/backend/daftar/modul_masterproductnamelife_gen.go` hasil `go generate` |
| Pengecualian *(ralat 01-10-2026)* | penjaga Claim Life `TestMasterViewTidakDisentuh` (`modul/claimlife/backend/repository/migrasibatas_test.go`) memindai seluruh `APP_RNM` dan melarang kode menyebut `M_PRODUCT_LIFE`/`PRODUCTINWARD_LIFE` — modul ini, sebagai **penulis**, pasti merah. Persempit pindaiannya ke `modul/claimlife/` saja, aturannya tidak berubah; satu commit sebelum paket 1. Rinciannya di `PROMPT-IMPLEMENTASI-TIGA-MODUL-LIFE-GELOMBANG-2.md` §2 |
| Dilarang | `inti/`, modul lain *(termasuk `mastercontractretrolife` dan `claimlife`)*, `package.json`, `package-lock.json`, `go.mod`, konfigurasi root. Pola modul lain boleh ditiru, tidak boleh diimpor |
| Nomor | rentang migrasi **140–179**, slot menu **960–961** |
| Oracle | skema `POOLDATA`. `-migrate` oleh **work owner**. Uji `db` tanpa `ORACLE_DSN` = SKIP; **POOLDATA bukan target uji `db`**. Nol penulisan ke DEV |
| Prosedur | `PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE` **tidak dipanggil**; logikanya ditiru di Go. Nol `COMMIT` di teks SQL |
| Layanan luar | `ConnectREST/ServiceGoogle`, `SystemSettings/LinkService`, tombol **View Office Online**: **stub** lewat pola outbox; alamat nyata **tidak dipanggil dan tidak ditulis** ke berkas apa pun |
| Rahasia | nol nama orang, nomor polis, kredensial, host/IP, URL `M_LINK_SERVICE`. Fixture berawalan `UJI-` |
| CSS | kelas berawalan `mpnl-` di `modul/masterproductnamelife/frontend/masterproductnamelife.css`, diimpor dari `rute.tsx` |

## 2. KORPUS `D:\XML\RNM_BRD\Master Product Name Life\` *(READ-ONLY, 114 berkas XML)*

| Jenis | Rule |
| --- | --- |
| Harness | `InwardProductName` |
| Section | `InboxProductName`, `InwardProductName`, `EditProductName_Confirm`, `SaveProductName_Confirm`, `ViewRate`, tujuh pemilih: `Ceding_Section`, `Currency_Section`, `PolicyHolder_Section`, `SOB_Section`, `RIRISK_Section`, `RIRate_Section`, `CauseOfLoss_Section` |
| FlowAction | `ChooseCeding`, `ChooseCurrency`, `ChoosePolicyHolder`, `ChooseSOB`, `ChooseRIRisk`, `ChooseRIRate`, `ChooseCauseOfLoss`, `EditProductName_Confirm`, `SaveProductName_Confirm`, `ProductNameAttachContent`, `ViewRate` |
| DataTransform | `CopyProduct`, `HideCreateLife`, `SetViewEdit`, `TreatyInIDSetPyPortal`, enam `set*_DT` |
| Lain | Activity 37, RDBList 22, ReportDefinition 17, `DecisionTable/GetMimeType`, `When/recordEvent`, `ConnectREST/ServiceGoogle`, `SystemSettings/LinkService` |
| Label tombol di section | `Add`, `Delete`, `Save`, `Choose`, `View`, `View Rate`, `View Reinstype`, `View Office Online`, `Refresh`, `Generate`, `Inward` — setiap satu dicari pemanggilnya |
| Tabel Retro Life yang dibaca | `TREATYYEAR_LIFE`, `TREATYCONTRACT_LIFE` *(`Activity/NewProductLife`, `ReportDefinition/BrowseTreatyYear_Life_RD`, `Section/InboxProductName`)* — dibaca **langsung dengan SQL modul ini**, bukan lewat kode modul Retro Life |

Hasil grilling: `modul/masterproductnamelife/docs/` — `spec.md`, `grilling-ronde-1*.md`, `dba-procedures-and-ddl.md`, `issues/01`–`09`, dan
**`dba-view-produk-life.md`** *(definisi tiga view dari DEV, dibaca asisten 30-09-2026)*.

## 3. FAKTA DEV *(asisten 30-09-2026, agregat, baca-saja)*

| Hal | Isi |
| --- | --- |
| Tabel produk | `M_PRODUCT_LIFE` 196 baris, `M_PRODUCTINWARD_LIFE` 197 — `JSONDATA` + satu constraint `IS JSON`, **nol PK** |
| Pembaca | view `PRODUCT_LIFE`, `PRODUCTINWARD_LIFE`, `DOCUMENTCLAIM_LIFE`; prosedur `PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE`; **kode Claim Life** membaca view `PRODUCTINWARD_LIFE` *(ambang produk, kategori dokumen, rate, validasi tanggal)* |
| Sequence | `M_PRODUCT_LIFE_SEQ`, `M_PRODUCT_INWARD_LIFE_SEQ`, `M_PRODUCT_TYPE_LIFE_SEQ` |
| Lampiran | `M_ATTACHMENTPRODUCTNAME` 2 baris, `T_STORAGE_IMAGE` 517, `T_FOLDER_IMAGE` 1 — tabel warisan |

> ⛔ **RALAT ASISTEN 01-10-2026:** P1 di bawah menyebut empat kolom datar `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE`. **Keliru.** Di DEV
> `M_PRODUCT_LIFE` hanya punya `ID`, `JSONDATA`, `RIRISKID`, `RIRISK`. Perbaikannya di `PROMPT-LANJUTAN-MASTER-PRODUCT-NAME-LIFE-1.md` L1.

## 4. RALAT — mengalahkan spec dan tiket bila bertentangan

| # | Klaim di dokumen | Fakta | Menjadi |
| ---: | --- | --- | --- |
| P1 | tiket **01** dan spec penyimpangan 2: skema relasional penuh, `product_life` + lima tabel anak, dua tabel lama digabung, JSON dibuang | tiga view, dua prosedur, dan Claim Life membaca `JSONDATA` kedua tabel lama; produk di tabel baru **tidak terlihat** oleh semuanya | **tiket 01 ditangguhkan → OQ-MPNL-01**. Bawaan: **ikut Pega** — simpan ke `M_PRODUCT_LIFE` dan `M_PRODUCTINWARD_LIFE` dengan `JSONDATA` berbentuk persis yang ditulis activity Pega, ditambah kolom flatten `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE` seperti `PEGA_M_PRODUCT_LIFE`. **Nol tabel baru, nol DDL**. Model Go tetap struct bernama; JSON hanya di repository |
| P2 | batasan AC 38 *"tidak mengurai JSONDATA m_product_life"* | batasan itu milik Claim Life sebagai **pembaca** | modul ini **penulis**, jadi boleh membaca dan menulis `JSONDATA` miliknya sendiri |
| P3 | `OutwardList` bukan tabel, jalur OR mati | view `PRODUCT_LIFE` membaca `OutwardList[0].OVR_COMM` | kunci `OutwardList` **tetap ditulis** dengan bentuk yang sama seperti Pega menulisnya |
| P4 | simpan atomik dua sisi *(tiket 03)* | tetap benar | kedua tabel ditulis dalam **satu transaksi** |
| P5 | lampiran dan efek keluar *(08, 09)* | tabel warisan ada; `ServiceGoogle` layanan luar | tulis ke tabel warisan seperti XML; kirim berkas lewat stub outbox |
| P6 | ID produk | sequence ada di DEV | format ID **dari XML/prosedur**, sebut buktinya *(`dba-procedures-and-ddl.md` dan RDB `Save*`)* |

**Uji wajib**: produk `UJI-` yang disimpan layanan menghasilkan `JSONDATA` yang, bila dibaca dengan **teks SQL ketiga view** di skema uji tiruan, memberi
nilai sama dengan masukannya. Uji murni: setiap kunci yang dibaca view *(daftar di `dba-view-produk-life.md`, peka huruf besar-kecil, termasuk ejaan
Pega `POLICYHODER`)* ada di JSON hasil simpan.

## 5. TIKET

| Tiket | Kerjakan |
| --- | --- |
| 01 | **ditangguhkan** *(P1)* — ralat bertanggal di tiket dan spec |
| 02 | produk sisi umum, identitas dari sequence, gagal terang-terangan |
| 03 | sisi inward *(harness `InwardProductName`)*, simpan atomik |
| 04 | tujuh pemilih master: `Choose*` + `set*_DT` + RD/RDB sumber masing-masing |
| 05 | validasi wajib isi seragam kedua sisi — hanya aturan dari XML atau keputusan tertulis di tiket |
| 06 | daftar plan dan batas underwriting |
| 07 | daftar komentar, dokumen, underwriting finansial |
| 08, 09 | lampiran *(`ProductNameAttachContent`, `GetMimeType`)*, efek keluar *(P5)* |
| tambahan | `CopyProduct`, `ViewRate`, `View Reinstype`, `Generate`, `Refresh`, `EditProductName_Confirm`, `SaveProductName_Confirm` — bila XML menunjukkan perilakunya dan tiket belum mencakupnya, **dibangun** dan dicatat sebagai tambahan tiket |

## 6. MENU

Satu tombol **"Master Product Name Life"** di `MASTER` *(baris `M_NAV_MENU` sudah ada; menu datar 901)*. Halaman awal = layar yang di Pega menjadi pintu
masuk *(kemungkinan `InboxProductName` — pastikan dari XML, sebut buktinya)*. Layar lain dibuka dari dalam.
Slot `960_menu_masterproductnamelife.sql`: hanya `UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE WHERE KODE = 'masterproductnamelife'`,
`_down` ke `'0'`. **Nol `INSERT`**. `menu.ts` menyatakan `HALAMAN_AWAL`.

## 7. PAKET — satu commit per paket, vertikal *(backend + frontend + uji)*

| # | Paket |
| ---: | --- |
| 0 | `docs/PARITAS-LAYAR-DAN-AKSI.md` *(setiap section, tombol, aksi → activity `bNNN` + `pyStepsBlockName` → RDB/RD → tabel/kunci JSON → rute → komponen)* + `docs/RALAT-DEV-30-09-2026.md` *(P1–P6)* + ralat bertanggal di tiket 01 dan spec + `docs/OQ-MASTER-PRODUCT-NAME-LIFE.md` |
| 1 | kerangka: `backend/modul.go`, repository baca kedua tabel lewat JSON, rute baca, penjaga modul *(nol migrasi DDL, nol tabel baru, kata cadangan Oracle, nol URL layanan)* |
| 2 | Inbox produk + tujuh pemilih *(04)* |
| 3 | produk sisi umum + identitas *(02)* |
| 4 | sisi inward + simpan atomik dua tabel + uji view *(03)* |
| 5 | validasi wajib isi *(05)* |
| 6 | plan dan batas underwriting *(06)* |
| 7 | komentar, dokumen, underwriting finansial *(07)* |
| 8 | lampiran dan efek keluar, stub *(08, 09)* |
| 9 | aksi tambahan dari XML *(bab 5 baris "tambahan")* |
| 10 | menu *(bab 6)* + `MODUL.md` Status → dimigrasi |
| 11 | dokumen: status tiket, PARITAS, OQ, `MODUL.md` *(prefix rute API, kontrak)* |

Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`)*, `gofmt`, `tsc`, `vitest`, `vite build`, di pohon kerja utama. Sesudah
paket 10: jalankan backend dan `npm run dev` dari `APP_RNM`, buka layar di peramban, pastikan tombol yang tampil = tombol di tabel PARITAS.

## 8. BERHENTI HANYA BILA

- Sebuah aturan hanya bisa dibangun dengan DDL pada tabel warisan — OQ, jangan tulis migrasi, lanjutkan paket lain.
- Perlu menyunting `inti/` atau modul lain — catat, lanjutkan paket lain.
- Perilaku aksi XML tidak dapat dipastikan — OQ dengan bukti `bNNN`, lanjutkan paket lain.

## 9. LAPORAN

Satu pesan: tabel **paket → commit → tiket → bukti XML → angka uji**; ralat P1–P6 yang diterapkan; tombol dan aksi XML di luar tiket beserta statusnya;
OQ baru; halaman awal beserta buktinya; contoh `JSONDATA` hasil simpan produk `UJI-` *(tanpa data orang)*; pengingat `-migrate` slot 960; persentase
kesiapan modul dan total migrasi *(modul dimigrasi dari 20)*; bab **TELEMETRI EKSEKUSI**.

## 10. OQ UNTUK WORK OWNER

| OQ | Pertanyaan | Bawaan sampai dijawab |
| --- | --- | --- |
| OQ-MPNL-01 | Produk tetap JSON di dua tabel lama seperti Pega, atau dipindah ke tabel relasional baru? Tabel baru membuat tiga view, dua prosedur, dan Claim Life tidak melihat produk baru, kecuali kedua bentuk ditulis bersamaan | JSON seperti Pega |

---

*Disusun 1 Oktober 2026 dari bab 3 `PROMPT-IMPLEMENTASI-TIGA-MODUL-RUMPUN-LIFE.md`, katalog DEV (30-09-2026), dan daftar rule korpus
`Master Product Name Life`.*
