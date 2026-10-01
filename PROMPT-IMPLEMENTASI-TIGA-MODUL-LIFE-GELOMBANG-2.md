# PROMPT — GELOMBANG 2 RUMPUN LIFE, TIGA MODUL SEKALIGUS: **Master Contract Retro Life *(penutup)* · Master Product Name Life · Endorsement Life** *(sesi baru, folder `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM`, cabang **`dev`**)*

> Permintaan work owner 01-10-2026: *"kerjakan 3 modul sekaligus, jangan hanya 1 tiket, jangan banyak yang di-skip, baca ulang XML, ambil logic, button,
> dan method yang benar dari XML, jangan membuat menu yang tidak ada di Pega; XML patokan dasar; tiket = hasil grilling; tiket yang salah diperbaiki di
> dokumennya."* Gelombang ini **menutup rumpun Life**: sesudahnya kelima modul Life *(Retro, Product Name, PremiumList, Endorsement, Claim Life + Komite)*
> berdiri semua.
>
> **Mengapa bukan NB Treaty In:** 92 activity, 75 when, rancangan penyimpanan 984 baris yang masih menimbang tabel baru, dan pertanyaan untuk DBA,
> Finance, IAM, dan Product belum dijawab. Modul itu butuh satu sesi tersendiri sesudah keputusan penyimpanannya diambil.

> ⛔ **Syarat mulai (01-10-2026):** `PROMPT-SERAGAM-KOLOM-T_WORK_POLIS-DAN-T_WORK_CLAIM.md` sudah selesai di `dev`. Sejak itu kolom
> `T_WORK_POLIS.STATUS` bernama **`STATUS_WORK`**, dan tabel itu punya `COVER_KEY`, `CREATE_OP`, `CREATE_OP_NAME`, `TGL_CREATE`, `TGL_UPDATE`.
> Endorsement yang membuat atau mengubah kasus polis wajib mengisi kolom pembuat dan waktu seperti PremiumList.

## 0. VERIFIKASI SESI RETRO LIFE *(asisten, 01-10-2026)*

| Klaim log | Hasil cek |
| --- | --- |
| 12 commit `99c2bac` → `5efc1a7` | ✅ ada, berurutan |
| Di luar folder modul hanya 2 berkas uji + berkas bangkitan | ✅ `frontend/daftar.datar.test.ts`, `frontend/daftar.sinkron.test.ts`, `inti/backend/daftar/modul_mastercontractretrolife_gen.go` |
| Nol DDL, slot 958 hanya `UPDATE` | ✅ `UPDATE … DIMIGRASI = '1'`, `_down` ke `'0'`; nol `CREATE/ALTER/DROP` |
| Nol alamat host di kode | ✅ |
| Gerbang | ✅ diulang asisten: vet bersih, Go **1.132** lulus · 0 gagal, tsc bersih, vitest **924** di **72** berkas, build **120** modul |
| Halaman awal `GridRetrocessionLife` | ✅ judul `MASTER CONTRACT RETRO LIFE` ada; tidak dirujuk rule lain di korpus modul; keempat harness dibuka `showHarness` dari dalam section |
| Nomor baris `bNNN` | ⚠️ **bergeser** dari cara pecah standar `sed -e 's/></>\n</g'`: judul di **b1082** *(log: b1017)*, `InputRetrocessionLife` **b2502** *(log: b2475)*, `showHarness` **b12725/b13005/b13590** *(log: b12720/b13000/b13585)*. Isinya benar, penomorannya tidak dapat diulang |
| OQ-MCRL-13 rate | DEV punya view **`RATE_LIFE`** dan **`RATE_LIFE_SUMMARY`** *(di atas `M_RATE_LIFE` 96.038 baris, `M_RATE_LIFE_SUMMARY` 339)*. Tertahan oleh penjaga Claim Life — lihat §2 |

## 1. ATURAN BACA XML — berlaku di setiap tiket ketiga modul

1. **Setiap pembacaan activity mencetak `pyStepsBlockName`**. `//` = langkah ter-remark, tidak pernah jalan, jangan dibangun.
2. Nomor baris **hanya** dari `sed -e 's/></>\n</g' <berkas> | grep -n …` — tulis perintahnya sekali di kepala PARITAS, supaya setiap `bNNN` dapat diulang orang lain.
3. WHEN dibaca dari urutan aksi. Precondition `WhenTrue 2` = lanjut, `3` = lewati.
4. Cari **penulis** sebuah kolom, bukan hanya pembacanya.
5. Setiap tombol = satu tombol di section XML yang memanggil activity yang sama. Label **VERBATIM**. Nol menu, layar, tombol, aksi tanpa bukti XML.
6. Aksi XML yang perilakunya jelas tetapi tidak ada di tiket = **dibangun** dan dicatat sebagai tambahan tiket. Yang tidak jelas = OQ dengan bukti.
7. Tiket yang berbeda dari XML atau katalog DEV diperbaiki di dokumennya *(catatan bertanggal, kalimat lama dikutip)*. `grilling-ronde-*` tersegel.

## 2. BATAS KERJA

| Hal | Isi |
| --- | --- |
| Folder repo | **`D:\XML\RNM_BRD\OUTPUT_HASIL_RNM`**, bukan klon lain |
| Cabang | **`dev`**. `git commit -o -- <jalur>`, pesan berawalan nama modul. **Nol `git push`, nol `git pull`, nol `git pull --rebase`** *(GitHub `dev` memuat kenaikan Vite 8 yang membuat `npm ci` gagal)* |
| Folder boleh disunting | `APP_RNM/modul/mastercontractretrolife/**`, `APP_RNM/modul/masterproductnamelife/**`, `APP_RNM/modul/endorsementlife/**`, ketiga `inti/backend/daftar/modul_<nama>_gen.go` |
| **Satu pengecualian, wajib paket 0** | `modul/claimlife/backend/repository/migrasibatas_test.go`: `TestMasterViewTidakDisentuh` kini memindai **seluruh `APP_RNM`** *(`akarModul = "../../../.."`)* dan melarang setiap kode menyebut `M_PRODUCT_LIFE`, `PRODUCTINWARD_LIFE`, `RATE_LIFE`. Larangan itu milik **Claim Life sebagai pembaca** *(AC 38, OQ-M7)*. Persempit pindaiannya ke **`modul/claimlife/`** saja; aturan dan daftar izin Claim Life **tidak berubah**. Uji gigit: kode Claim Life yang menyebut `M_PRODUCT_LIFE` tetap merah. Tanpa ini, Product Name *(penulis tabel produk)* dan Retro Life *(pembaca rate)* tidak dapat hijau |
| Dilarang | `inti/`, modul lain selain pengecualian di atas, `package.json`, `package-lock.json`, `go.mod`, konfigurasi root |
| Nomor | Retro **100–139** / slot **958** · Product Name **140–179** / slot **960** · Endorsement **480–519** / slot **976** |
| Oracle | skema `POOLDATA`. `-migrate` oleh **work owner**. Uji `db` tanpa `ORACLE_DSN` = SKIP; **POOLDATA bukan target uji `db`**; jalankan uji `db` dengan `-p 1` *(skema uji bersama, nama fixture dapat bertabrakan)*. Nol penulisan ke DEV |
| Prosedur | tidak dipanggil, logikanya ditiru di Go. Nol `COMMIT` di teks SQL |
| Layanan luar | `ConnectREST` *(Product Name `ServiceGoogle`, Endorsement)*, `LinkService`, **View Office Online**: stub outbox; alamat nyata tidak dipanggil dan tidak ditulis |
| Rahasia | nol nama orang, nomor polis, kredensial, host/IP, URL `M_LINK_SERVICE`. Fixture `UJI-` |
| CSS | kelas berawalan modul di berkas CSS folder frontend modul itu |

## 3. MODUL A — MASTER CONTRACT RETRO LIFE *(penutup, ±85% → 100% kecuali OQ)*

| Butir | Kerjakan |
| --- | --- |
| A1 | **OQ-MCRL-13 rate.** Bila `OQ-MASTER-CONTRACT-RETRO-LIFE.md` mencatat **izin work owner** *(lihat §10)*: baca `RATE_LIFE_SUMMARY` untuk autocomplete `R/I RATE` dan `RATE_LIFE` untuk section `Rate List`, **baca-saja**, kolom yang dibaca RD XML saja, nol `SELECT *`, nol `JSONDATA`. Rute yang kini 503 menjawab 200; simpan business baru jalan. Tanpa izin: tetap 503 berkalimat |
| A2 | Nomor `bNNN` di `PARITAS-LAYAR-DAN-AKSI.md` dan `RALAT-DEV-30-09-2026.md` dihitung ulang dengan perintah §1.2 |
| A3 | OQ-MCRL-07 tetap **rute API tanpa layar** — Pega tidak punya layarnya |
| A4 | uji manual: backend + `npm run dev`, buka kelima layar, cocokkan setiap tombol dengan PARITAS; catat hasilnya di `LAPORAN-UJI-MANUAL.md` modul |

## 4. MODUL B — MASTER PRODUCT NAME LIFE

**Seluruh isi `PROMPT-IMPLEMENTASI-MODUL-MASTER-PRODUCT-NAME-LIFE.md` berlaku** — korpus, fakta DEV, ralat **P1–P6**, tiket, paket 0–11, OQ-MPNL-01.
Ringkas: produk disimpan **seperti Pega** sebagai `JSONDATA` di `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE` *(tiga view, dua prosedur, dan Claim Life
membacanya)*, kunci JSON persis daftar `docs/dba-view-produk-life.md`, satu transaksi untuk kedua tabel, nol tabel baru, tiket 01 ditangguhkan.

## 5. MODUL C — ENDORSEMENT LIFE

### 5.1 Fakta DEV dan kode *(asisten 30-09-2026)*

| Hal | Isi |
| --- | --- |
| Tabel aplikasi milik PremiumList | `T_PREMIUM_LIST`, `T_PREMIUM_LIST_DETAIL` *(migrasi 051, 052)* — ada di DEV, **0 baris** |
| Tabel warisan | `JSON_POLIS` 148 baris, `M_TEMPUPLOADLIFE` 1.350, **`M_LIFE_PREMIUM_DETAIL` ±66,8 juta baris** |
| Penulis `M_LIFE_PREMIUM_DETAIL` | PremiumList *(`polis_warisan.go`, cermin `SaveMasterLPDet`)* |
| Pembaca peserta | Claim Life membaca `M_LIFE_PREMIUM_DETAIL` dan **sudah** menyaring `EDMSTATUS` `Batal`/`Delete` *(`pesertapolis.go`)* |
| Korpus | 75 XML: Harness 8 *(`EditDetail_Harness`, `EndorsmentLife_harnes`, `InboxEndorsementLife`, `ViewCSVResult_LifeEDM`, `ViewOldPolicy_EDM`, `_QP`, `_TP`, `_TR`)*, Section 14, Activity 16, RDBList 17, RD 6, FlowAction 6, Flow 1, When 2, DecisionTable 1, ConnectREST 1 |
| Dokumen | `modul/endorsementlife/docs/`: `spec.md`, `grilling-ronde-1.md`, `grilling-ronde-2.md`, `keputusan-struktur-edm.md`, `STRUKTUR-TABEL-ENDORSEMENT-LIFE.md`, `BAHAN-tiket-pl-number-edm.md`, `issues/00`–`12` |

### 5.2 Ralat

| # | Klaim | Fakta | Menjadi |
| ---: | --- | --- | --- |
| E1 | tiket **00** `ALTER` dua tabel, diblok PremiumList tiket 00 | kedua tabel sudah ada, blok terpenuhi | dikerjakan di rentang **480–519**, folder migrasi Endorsement. **`PRODKE NUMBER DEFAULT 1`** supaya kode PremiumList tidak diubah dan baris new business otomatis versi 1; kolom EDM lain nullable; `PARENT_ID` self-FK nullable ber-index. Tabel aplikasi milik kita, bukan warisan |
| E2 | tiket **11** kontrak hilir ke Claim Life | Claim Life sudah menyaring `EDMSTATUS` | Endorsement **menulis `M_LIFE_PREMIUM_DETAIL`** dengan `EDMSTATUS` `Old`/`New`/`Delete`/`Batal` persis `Endorsement Life/RDBList/SaveMasterLPDet.xml`, dan summary seperti `InsertPLSummary`. **Nol perubahan kode Claim Life dan PremiumList**. Uji: baris `Delete`/`Batal` tidak lolos penyaring yang sama *(ditiru di uji, tidak diimpor)* |
| E3 | tiket **12** migrasi data lama, `needs-info` | keputusan memindah `JSON_POLIS` belum ada | **di luar lingkup**, tetap `needs-info` |
| E4 | `M_LIFE_PREMIUM_DETAIL` | ±66,8 juta baris | setiap kueri memakai kunci ber-index *(sebut indeksnya)*; **nol pemindaian penuh**; uji tidak menyentuhnya |
| E5 | efek keluar dan alarm *(tiket 10)* | `ConnectREST` | stub outbox; alarm = pesan di layar + log tanpa alamat |

### 5.3 Tiket

| Tiket | Isi | Rule XML utama yang wajib dibaca |
| --- | --- | --- |
| 00 | kolom EDM + `PARENT_ID` *(E1)* | `pyFields DATA_JSON` EDM, `Generate_NoEndorsmentLife` |
| 01 | pilih polis dan gerbang kelayakan | section Inbox, activity pemilih polis |
| 02 | buat case dan muat polis lama | `EndorsmentLife_harnes`, flow |
| 03 | popup lihat polis lama per Type | `ViewOldPolicy_EDM`, `_QP`, `_TP`, `_TR` |
| 04 | penomoran `PL_NUMBER_EDM` | `BAHAN-tiket-pl-number-edm.md` |
| 05 | maksud endorsement dan `EDMSTATUS` | `EdmType` `1` Perubahan Data / `3` Batal |
| 06 | jurnal balik delete dan batal | `SaveMasterLPDet` |
| 07 | unggah CSV | `ViewCSVResult_LifeEDM`, `M_TEMPUPLOADLIFE` |
| 08 | alur keputusan, kunci field, audit | flow action keputusan |
| 09 | simpan transaksi campuran, anti-dobel | RDB simpan |
| 10 | efek keluar *(E5)* | `ConnectREST` |
| 11 | kontrak hilir *(E2)* | `SaveMasterLPDet`, `InsertPLSummary`, `Claim Life/RDBList/GetPesertaClaim_sql1` |
| 12 | **tidak** *(E3)* | — |

Halaman awal = `InboxEndorsementLife` bila XML menunjukkan harness itu pintu masuk portal *(sebut buktinya)*. Tombol menu **"Endorsement Life"** di `TREATY`.

## 6. MENU

Menu datar *(901 di `dev`)*. Satu tombol per modul dari baris `M_NAV_MENU` yang sudah ada. Slot tiap modul hanya:

```sql
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE WHERE KODE = '<modul>'
/
```

ditambah `_down` ke `'0'`. Nol `INSERT`. `menu.ts` menyatakan `HALAMAN_AWAL`.

## 7. URUTAN — tiga modul bergantian, satu commit per paket, vertikal

| Gelombang | Retro Life | Product Name | Endorsement |
| ---: | --- | --- | --- |
| 1 | pengecualian penjaga §2 *(satu commit, sebelum yang lain)* + A2 | paket 0 PARITAS + RALAT P1–P6 + OQ | PARITAS + RALAT E1–E5 + OQ |
| 2 | A1 bila diizinkan | paket 1–2 kerangka, Inbox, tujuh pemilih | tiket 00 migrasi + kerangka + baca polis |
| 3 | — | paket 3–4 sisi umum, inward, simpan atomik + uji view | 01, 02, 03 |
| 4 | — | paket 5–7 validasi dan daftar | 04, 05 |
| 5 | — | paket 8–9 lampiran, aksi tambahan XML | 06, 07, 09 |
| 6 | A4 uji manual | paket 10 menu + frontend lengkap | 08, 10, 11 + frontend lengkap + menu |
| 7 | dokumen | paket 11 dokumen | status tiket, PARITAS, OQ, `MODUL.md` |

Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`; tag `db` dengan `-p 1`)*, `gofmt`, `tsc`, `vitest`, `vite build`, di pohon kerja
utama. Satu tiket yang terhenti tidak menghentikan dua modul lainnya.

## 8. BERHENTI HANYA BILA

- Aturan hanya bisa dibangun dengan DDL pada **tabel warisan** — OQ, jangan tulis migrasi.
- Perlu menyunting `inti/` atau modul lain di luar pengecualian §2 — catat, lanjutkan.
- Perilaku aksi XML tidak dapat dipastikan — OQ dengan bukti `bNNN`, lanjutkan.

## 9. LAPORAN

Satu pesan:

1. Tabel **modul → paket → commit → tiket → bukti XML → angka uji**.
2. Ralat yang diterapkan *(A, P, E)* dan tiket yang diperbaiki dokumennya.
3. Tombol dan aksi XML di luar tiket beserta statusnya.
4. OQ baru per modul.
5. Halaman awal Product Name dan Endorsement beserta buktinya.
6. Contoh `JSONDATA` produk `UJI-` dan contoh baris `M_LIFE_PREMIUM_DETAIL` endorsement `UJI-`, tanpa data orang.
7. Pengingat `-migrate`: Endorsement 480-an, slot 958, 960, 976.
8. Persentase kesiapan per modul dan total migrasi *(modul dimigrasi dari 20)*.
9. Bab **TELEMETRI EKSEKUSI**.

## 10. OQ UNTUK WORK OWNER

| OQ | Pertanyaan | Bawaan sampai dijawab |
| --- | --- | --- |
| OQ-MCRL-13 | Retro Life boleh **membaca** view `RATE_LIFE_SUMMARY` dan `RATE_LIFE` *(baca-saja, kolom RD saja)*? Tanpa ini, business baru tidak bisa disimpan | tetap 503 |
| OQ-MCRL-01 | Gerbang tahun Retro Life mati di Pega. Tetap ditegakkan? | ikut Pega |
| OQ-MPNL-01 | Produk tetap JSON seperti Pega, atau tabel relasional baru? | JSON seperti Pega |
| OQ-EDM-001 | Data endorsement lama di `JSON_POLIS` dipindah? | tidak |

Jawaban work owner dicatat sesi ini di berkas OQ modul masing-masing, bertanggal, sebelum butir yang bergantung padanya dikerjakan.

---

*Disusun 1 Oktober 2026 dari verifikasi ulang sesi Retro Life (commit, batas folder, slot, XML, gerbang), katalog DEV (view dan tabel rate), penjaga
`TestMasterViewTidakDisentuh`, `PROMPT-IMPLEMENTASI-MODUL-MASTER-PRODUCT-NAME-LIFE.md`, dan bab 4 `PROMPT-IMPLEMENTASI-TIGA-MODUL-RUMPUN-LIFE.md`.*
