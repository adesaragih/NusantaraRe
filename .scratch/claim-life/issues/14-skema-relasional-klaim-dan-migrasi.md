# 14: Skema relasional klaim (6 tabel, berakar di `T_WORK_CLAIM`) + migrasi — **PREFACTOR**

**Status:** claimed

**Blocked by:** 01 (kerangka aplikasi + seam API)

⚠️ **Ini PREFACTOR dan harus dikerjakan LEBIH DULU daripada tiket 02–13.** Nomornya 14 karena
seri 01–13 sudah terbit dan tidak dinomori ulang — **urutan ditentukan tepi pemblokir, bukan
nomor**. Tiket **02, 03, 04, 05, 12, 13** kini memblokir pada tiket ini.

Alasannya: bentuk penyimpanan berubah total — dokumen JSON + satu tabel flat warisan → **enam tabel
relasional**. Tidak ada irisan lain yang dapat berdiri sebelum bentuk barunya ada.
*"Make the change easy, then make the easy change."*

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin setiap atribut klaim menjadi **kolom bernama** dengan tipe yang
benar, dan seluruh klaim lama pindah **tanpa kehilangan satu nilai pun** — termasuk **seluruh baris
adjustment**, bukan hanya keadaan terakhir. Sebagai **organisasi**, saya ingin bentuk klaim dapat
diperiksa, dicari, dan divalidasi, bukan tersembunyi di dalam satu dokumen. *(Spec §2b, §14)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL enam tabel klaim + `T_WORK_CLAIM` + sequence; skrip migrasi & rekonsiliasi |
| `internal/models` | Agregat klaim: header → peserta → adjustment & dokumen |
| — | Skrip rekonsiliasi nilai uang, tanggal, dan jumlah baris |

## Bentuk baru — spec §2b RALAT D/E (2026-09-18)

⚠️ **REVISI 2026-09-18** `[keputusan work owner]` — **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING`
DIHAPUS**, dan **`T_WORK_CLAIM` naik menjadi AKAR.** Skema klaim Life turun dari **8 tabel menjadi
6**, dan pohonnya kini **enam tingkat**. ⛔ Jangan membuat `T_CLAIMLF_POLICY` maupun
`T_CLAIMLF_MARKETING` — penghapusan ini **bukan** penggantian nama.

```
TINGKAT 1   T_WORK_CLAIM ─────────────────── akar · LINTAS-LINI · satu baris per work object
            PK  ID   TEKS BERFORMAT: CLM-xxxxxx (klaim) / KMT-xxxxxx (komite)
            FK  COVER_KEY → T_WORK_CLAIM.ID   (menunjuk dirinya sendiri; NULL bila tak punya induk)
            LINI · PY_POSITION · ACCEPT_STATUS · SENDTO_ADMIN · SENDTO_MEDICAL · TYPE
            CASEID · CREATE_OP · CREATE_OP_NAME · TGL_UPDATE  <- PINDAHAN dari header klaim
            |
   +--------+-------------------------------------------+
   | baris KLAIM  COVER_KEY = NULL                      | baris KOMITE  COVER_KEY = ID baris klaim
   |                                                    |
TINGKAT 2                                           TINGKAT 2
   +--1:1-- T_GENERAL_CLAIM                             +--1:1-- T_GENERAL_KOMITE
           PK ID = T_WORK_CLAIM.ID   <- SHARED PK               PK ID = T_WORK_CLAIM.ID baris komite
             (tidak ada kolom FK terpisah)                        <- SHARED PK, tidak ada kolom FK
           penunjuk polis (ke LUAR, bukan anak):                FK ADJUSTMENT_ID -> T_CLAIMLF_ADJUSTMENT.ID
             CASEID_POLICY   [terbuka]                            (tutup lingkar)
             POLICY_NO       [terbuka]                          KOMITE_LOOP · KOMITE_COUNT · ACCEPT_STATUS
             ENDORSMENT_NO   [terbuka]                          |
           |                                         TINGKAT 3  |
TINGKAT 3  |                                                    +--1:N-- T_KOMITE_KOMITELIST
           +--1:N-- T_CLAIMLF_PREMIUMLIST_DETAIL                        PK ID
           |        PK ID                                              FK DATA_KOMITE_ID
           |        FK CLAIM_ID -> T_GENERAL_CLAIM.ID  CASCADE             -> T_GENERAL_KOMITE.ID  CASCADE
           |        |                                                   KOMITE_URUT · KOMITE_ID · ID_KOMITE
TINGKAT 4  |        +--1:N-- T_CLAIMLF_ADJUSTMENT                       KOMITE_EMAIL · KOMITE_APROVAL
           |        |        PK ID                                      KOMITE_COMMENT · DATE_APPROVE
           |        |        FK PREMIUM_LIST_DETAIL_ID
           |        |           -> T_CLAIMLF_PREMIUMLIST_DETAIL.ID  CASCADE
           |        |        FK KOMITE_ID -> T_WORK_CLAIM.ID  <- penunjuk balik KE ATAS, nullable
           |        |        |
TINGKAT 5  |        |        +--1:N-- T_CLAIMLF_ADJUSTMENT_SPREADING
           |        |                 FK ADJUSTMENT_ID -> T_CLAIMLF_ADJUSTMENT.ID  CASCADE
           |        |                 |
TINGKAT 6  |        |                 +--1:N-- T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO
           |        |                          FK SPREADING_ID
           |        |                             -> T_CLAIMLF_ADJUSTMENT_SPREADING.ID  CASCADE
TINGKAT 4  |        +--1:N-- DOCUMENT_CLAIM                      <- LINTAS-LINI
           |                 FK PREMIUM_LIST_DETAIL_ID -> ...DETAIL.ID   (Life saja)
           |                 ON DELETE di Go -- induk beda tabel per lini
           |
           +-- DIBACA dari luar, BUKAN anak, TIDAK ikut cascade:
                 T_PREMIUM_LIST dkk (polis) · master marketing officer · EMAILKOMITE ·
                 master retro · security reinsurer · currency
```

**ENAM tabel klaim** — `T_GENERAL_CLAIM`, `T_CLAIMLF_PREMIUMLIST_DETAIL`, `T_CLAIMLF_ADJUSTMENT`,
`T_CLAIMLF_ADJUSTMENT_SPREADING`, `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`, `DOCUMENT_CLAIM` —
digantung pada akar `T_WORK_CLAIM`. Kedalaman **enam tingkat**.

✅ **Ejaan penaut induk-anak — DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.**
`[keputusan work owner]` **Ejaan final `COVER_KEY` (snake_case);** `CoverKey` warisan Pega
**tidak dipakai** sebagai nama kolom.
Alasan: seluruh kolom SQL baru proyek ini snake_case, dan Oracle melipat identifier tanpa kutip
menjadi huruf besar — `CoverKey` akan menjadi identifier **tanpa garis bawah**, berbeda dari
kolom-kolom sekitarnya. Kolomnya dibuat oleh **tiket 00 Komite Claim Life** lewat `ALTER`;
tiket ini hanya membuat tabel `T_WORK_CLAIM`.

### Sebelas relasi — seluruh kunci tamu ber-index

| # | Induk | Anak | Kunci tamu | Kard | Hapus |
| --- | --- | --- | --- | --- | --- |
| 1 | `T_WORK_CLAIM` | `T_WORK_CLAIM` | `COVER_KEY` | 1:N | di Go |
| 2 | `T_WORK_CLAIM` | `T_GENERAL_CLAIM` | **tidak ada kolom terpisah** — `T_GENERAL_CLAIM.ID` = `T_WORK_CLAIM.ID` (**shared PK**) | 1:1 | — |
| 3 | `T_GENERAL_CLAIM` | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `CLAIM_ID` | 1:N | CASCADE |
| 4 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `T_CLAIMLF_ADJUSTMENT` | `PREMIUM_LIST_DETAIL_ID` | 1:N | CASCADE |
| 5 | `T_CLAIMLF_ADJUSTMENT` | `T_CLAIMLF_ADJUSTMENT_SPREADING` | `ADJUSTMENT_ID` | 1:N | CASCADE |
| 6 | `T_CLAIMLF_ADJUSTMENT_SPREADING` | `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` | `SPREADING_ID` | 1:N | CASCADE |
| 7 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `DOCUMENT_CLAIM` | `PREMIUM_LIST_DETAIL_ID` | 1:N | **di Go** |
| 8 | `T_WORK_CLAIM` | `T_GENERAL_KOMITE` | **tidak ada kolom terpisah** — `T_GENERAL_KOMITE.ID` = `T_WORK_CLAIM.ID` baris komite (**shared PK**) | 1:1 | di Go |
| 9 | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | 1:N | CASCADE |
| 10 | `T_CLAIMLF_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` | 1:1 | di Go |
| 11 | `T_CLAIMLF_ADJUSTMENT` | `T_WORK_CLAIM` | `KOMITE_ID` | N:1 | penunjuk |

Relasi **8 · 9 · 10** dan kolom `COVER_KEY` dibuat oleh **tiket 00 Komite Claim Life**, bukan di
sini. Tiket ini membuat relasi **1 · 2 · 3 · 4 · 5 · 6 · 7 · 11** dan **tabel** `T_WORK_CLAIM`.

✅ **REVISI 2026-09-18 — relasi 2 dan 8 memakai SHARED PRIMARY KEY.** `[keputusan work owner]`
Keduanya **tidak punya kunci tamu**: `T_GENERAL_CLAIM.ID` **sama persis** dengan `T_WORK_CLAIM.ID`
baris klaim, dan `T_GENERAL_KOMITE.ID` **sama persis** dengan `T_WORK_CLAIM.ID` baris komite.
⛔ **Kolom `WORK_CLAIM_ID` DIBUANG** — bukan diganti nama, **tidak ada**. Test yang menemukan
kolom `WORK_CLAIM_ID` **gagal**. Karena keduanya PK, indexnya sudah ada dengan sendirinya.

✅ **`DOCUMENT_CLAIM` — `ON DELETE` DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.**
`[keputusan work owner]` Penghapusannya **ditangani di Go**, **bukan** cascade basis data. Alasan:
ia **LINTAS-LINI** dan induknya **tabel yang berbeda per lini**, sehingga satu `ON DELETE CASCADE`
tidak dapat seragam. Untuk Life induknya tetap `T_CLAIMLF_PREMIUMLIST_DETAIL`.

### Kolom yang bertambah dan yang pindah — 2026-09-18

| Tabel | Aksi | Kolom |
| --- | --- | --- |
| `T_GENERAL_CLAIM` | **TAMBAH** | `CASEID_POLICY` · `POLICY_NO` · `ENDORSMENT_NO` — ketiganya ⚠️ `[terbuka]` |
| `T_GENERAL_CLAIM` | **BUANG** | `PL_NUMBER` → diganti nama menjadi **`POLICY_NO`** |
| `T_GENERAL_CLAIM` | **BUANG** | `CREATE_OP` · `CREATE_OP_NAME` · `TGL_UPDATE` → **pindah** ke `T_WORK_CLAIM` |
| `T_GENERAL_CLAIM` | **BUANG** | `CASEID` → **pindah** ke `T_WORK_CLAIM` *(DIPUTUSKAN 2026-09-18)* |
| `T_WORK_CLAIM` | **TAMBAH** | `CREATE_OP` · `CREATE_OP_NAME` · `TGL_UPDATE` (pindahan) |
| `T_WORK_CLAIM` | **TAMBAH** | `CASEID` (pindahan) *(DIPUTUSKAN 2026-09-18)* |
| `T_WORK_CLAIM` | **TAMBAH** | `LINI` — kolom **ADA**; nilai Life = konstanta lini Life. `[terbuka — Non-Life]` hanya **daftar enum lintas-lini** |

⚠️ `[terbuka]` **Ketiga penunjuk polis belum menunjuk apa pun yang ada.** Ditulis apa adanya di atas
**bukan** karena sudah beres. **Jangan dijawab sendiri** — pemilik **work owner**:

| Penunjuk | Mengapa belum menunjuk |
| --- | --- |
| `CASEID_POLICY` | Menunjuk **kolom apa?** `T_PREMIUM_LIST` **sengaja membuang** `pyID`/`px*`/`py*` — catatannya berbunyi *"ID work Pega diganti ID sequence baru"*. **Case id Pega milik polis tidak disimpan di sisi sana.** |
| `POLICY_NO` | **Tidak ada di header polis.** Ia ada di `T_PREMIUM_LIST_DETAIL` — **per peserta**. Header polis punya `PL_NUMBER`. Join header-klaim → header-polis lewat `POLICY_NO` **tidak nyambung.** |
| `ENDORSMENT_NO` | **Belum menunjuk satu versi.** Di sisi polis namanya **`NOENDORS`** `[terverifikasi]` (`Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml`), dan **versi berjalan ditentukan `PRODKE` terbesar**. |

✅ **Nasib `CASEID` — DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.** `[keputusan work owner]`
`CASEID` **PINDAH** ke `T_WORK_CLAIM`. Alasannya **sama** dengan
`CREATE_OP`/`CREATE_OP_NAME`/`TGL_UPDATE`: ia identitas **work object**, bukan atribut klaim.
`T_GENERAL_CLAIM` **tidak lagi memuat** `CASEID`.

### Mengapa kedua tabel boleh dihapus `[terverifikasi]`

| Tabel | Bukti |
| --- | --- |
| `T_CLAIM_POLICY` | **27 dari 32 kolomnya tersedia di `T_PREMIUM_LIST`** (32 = 23 kolom dasar + 9 kolom tambahan hasil audit; `ID`/`CLAIM_ID` tidak dihitung). Lima sisanya kehilangan rumah — §Blocker. |
| `T_CLAIM_MARKETING` | **Terbukti turunan, nol nilai asli milik klaim.** `Claim Life/Activity/SetMOClaim_Act.xml` mengisi `MarketingData` dari `PolicyDataLife` (`MOID` · `MarketingCode` · `MarketingName` · `TeamGroup` · `BranchCode` · `BranchName`) **atau** dari master marketing officer (`ASM-FW-GISFW-Int-marketingofficer`, lewat `Claim Life/ReportDefinition/BrowseMarketingOfficer_RD.xml`). |

⚠️ **Penyimpangan sadar — potret berubah menjadi baca hidup.** `[keputusan work owner]` Keputusan
lama menyebut `PolicyDataLife` sebagai **snapshot** — *"disalin, bukan baca live"* — dan itu **kini
dibalik**. Akibatnya nyata dan diterima: **polis yang berubah sesudah klaim dibuat akan mengubah
tampilan klaim lama**; untuk polis ber-endorsement itu **pasti terjadi**. Migrasi karena itu
**tidak** menyalin atribut polis ke mana pun.

✅ **Spec §2b selaras** — `.scratch/claim-life/spec.md` §2b memuat RALAT 2026-09-18 A–E dengan pohon
dan tabel relasi yang sama. Teks lama di spec **tidak dihapus**, hanya diralat.

⚠️ **Dua tabel spreading tetap ada** — keduanya **terlewat di audit awal** dan ditemukan pada
verifikasi 2026-09-16. Bukti di bawah.

### Dua tabel spreading — `[terverifikasi]` induk = baris adjustment

`[terverifikasi]` `AdjustmentList` (class `ASM-FW-GISFW-Data-AdjustmentLife`) punya anak
`SpreadingList` → `RetroLifeList`. Bukti:

| Rule | Class / Nama / Tipe | Path | Isi |
| --- | --- | --- | --- |
| `SpreadingClaimLife_Act` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `SPREADINGCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/SpreadingClaimLife_Act.xml` | mengisi `.SpreadingList` **dan** `.RetroLifeList` beserta perhitungannya |
| `AdjustmentDetail_Section` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-HTML-SECTION` | `Claim Life/Section/AdjustmentDetail_Section.xml` | grid `.SpreadingList` |
| `RetroDetailClaimLife` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `RETRODETAILCLAIMLIFE` / `RULE-HTML-SECTION` | `Claim Life/Section/RetroDetailClaimLife.xml` | menampilkan `.REINSURERNAME`, `.PERCENTSHARE`, `.Amount` |

`[terverifikasi]` Kelas baris, dari deklarasi halaman di `SpreadingClaimLife_Act`:

```
.SpreadingList   = OutwardList.pxResults   → ASM-FW-GISFW-Int-TREATYYEAR_LIFE      (per treaty-year)
.RetroLifeList   = RetroLife.pxResults     → ASM-FW-GISFW-Int-RETROCESSIONLIFE     (per reinsurer)
```

**`T_CLAIMLF_ADJUSTMENT_SPREADING`** (per treaty-year) — `ID` (PK), **`ADJUSTMENT_ID`** (FK →
`T_CLAIMLF_ADJUSTMENT.ID`), `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`, `TREATY_YEAR_LIFE`,
`RETROCADED_SHARE` **NUMBER**, `RATE` **NUMBER**, `IDR` **NUMBER**, `USD` **NUMBER**, `CURRENCY`.

**`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`** (per reinsurer) — `ID` (PK), **`SPREADING_ID`** (FK →
`T_CLAIMLF_ADJUSTMENT_SPREADING.ID`), `REINSURER_NAME`, `PERCENT_SHARE` **NUMBER**, `AMOUNT`
**NUMBER**, `RATE` **NUMBER**, `PREMIUM_SPREADED_GROSS` **NUMBER**, `PREMIUM_SPREADED_NET`
**NUMBER**, `COMMISION` (sic) **NUMBER**, `OVR_COMM` **NUMBER**, `TREATY_TYPE_ID`,
`TREATY_TYPE_NAME`.

`[terverifikasi]` Perhitungan yang membekukan nilainya, dari `SpreadingClaimLife_Act`:

```
.RetrocadedShare        = local.ClaimNet   (atau .IDR / .USD menurut mata uang)
.Amount                 = local.AmountRetroShare * @divide(.PERCENTSHARE, 100, 4)
.RATE                   = @divide(local.Rate, 1, 4)          ← local.Rate sudah dibagi 1000 (per-mil)
.PREMIUM_SPREADED_GROSS = local.Rate * (1 + local.EMPercent) * .Amount
.PREMIUM_SPREADED_NET   = .PREMIUM_SPREADED_GROSS - local.Comm
.PREMIUM_SPREADED_NET   = .PREMIUM_SPREADED_GROSS - local.Discount - local.Comm      ← cabang kedua
```

⚠️ **`PREMIUM_SPREADED_NET` punya DUA rumus** `[terverifikasi]` — satu mengurangi komisi saja, satu
mengurangi diskon **dan** komisi. Keduanya ada di rule yang sama. Cabang mana yang berlaku
**tidak terbaca dari korpus** → **`[terbuka]`**, dicatat di §Blocker. **Jangan tebak.**

⚠️ **Penyimpangan sadar (baru) — spreading adjustment DIBEKUKAN dan DISIMPAN.**
`[keputusan work owner]` Nilai spreading pada saat adjustment dibuat tersimpan apa adanya; perubahan
master treaty sesudahnya **tidak** mengubahnya. Konsisten dengan spreading di PremiumList Life.

### `KOMITE_ID` pada `T_CLAIMLF_ADJUSTMENT`

⚠️ **Penyimpangan sadar (baru) — roster & keputusan komite TIDAK disimpan di Claim Life.**
`[keputusan work owner]` `T_CLAIMLF_ADJUSTMENT` menyimpan **satu kolom rujukan**, `KOMITE_ID`
(**nullable** — `NULL` bila baris belum pernah dikirim ke Komite). Roster dan keputusan per anggota
milik konteks **Komite Claim Life**; Claim Life **melihat**nya lewat join. **Tidak ada
`T_CLAIMLF_ADJUSTMENT_KOMITE`.**

⚠️ **`KOMITE_ID` berisi identitas kasus komite — yaitu `T_WORK_CLAIM.ID` baris komite.**
`[keputusan work owner]` REVISI 2026-09-17. Ia **tidak** menunjuk `T_GENERAL_KOMITE.ID`. Kolomnya
**ber-index**. Alasan work owner: saat user melihat baris adjustment, ia langsung tahu baris itu
merujuk kasus komite yang mana, tanpa query terbalik.

✅ **Tipe `KOMITE_ID` — DIPUTUSKAN 2026-09-18, `[terbuka]` tipe DITUTUP.**
`[keputusan work owner]` Ia mengikuti `T_WORK_CLAIM.ID`, yang kini **teks berformat**: baris klaim
`CLM-xxxxxx` (contoh `CLM-123456`), baris komite `KMT-xxxxxx` (contoh `KMT-000789`). **Bukan angka
sequence.** `KOMITE_ID` karena itu berisi teks berformat `KMT-xxxxxx`.

⚠️ **Penyimpangan sadar dari ADR-0006 — dicatat, bukan dilanggar diam-diam.**
`[keputusan work owner]` **ADR-0006** menetapkan identitas dari **sequence**. `T_WORK_CLAIM` dan
kedua tabel ber-shared-PK memakai **nomor bisnis berformat** alih-alih sequence murni. ADR-0006
tetap berlaku untuk identitas tabel klaim lainnya (`T_CLAIMLF_*`, `DOCUMENT_CLAIM`).

⚠️ `[terbuka]` **Dua hal TETAP terbuka. Jangan tebak:** (a) **generator** nomor `CLM-`/`KMT-` —
siapa yang membuatnya, apakah ada sequence di belakang prefiks, apakah di-reset per tahun —
pemilik **DBA / work owner**; (b) apakah `KOMITE_ID` dan `COVER_KEY` dipasangi
`REFERENCES T_WORK_CLAIM(ID)` atau dibiarkan tanpa constraint — pemilik **DBA / work owner**.

Rantai penuh untuk membaca keputusan komite dari sebuah baris adjustment:

```
T_CLAIMLF_ADJUSTMENT.KOMITE_ID
  → T_WORK_CLAIM   (ID = KOMITE_ID ; COVER_KEY = ID baris klaim)
  → T_GENERAL_KOMITE  (ID = T_WORK_CLAIM.ID  <- SHARED PK ; ADJUSTMENT_ID = adjustment.ID)
  → T_KOMITE_KOMITELIST (keputusan per anggota, diurut KOMITE_URUT)
```

⚠️ **`T_WORK_CLAIM` TIDAK memuat `ADJUSTMENT_ID`** `[keputusan work owner]` 2026-09-17 — penutup
lingkar ke baris adjustment ada di **`T_GENERAL_KOMITE.ADJUSTMENT_ID`**. Kolom **`COVER_KEY`** pada
`T_WORK_CLAIM` **ditambahkan oleh tiket 00 Komite Claim Life** lewat `ALTER`, bukan oleh tiket ini —
tiket ini hanya membuat tabelnya.

⚠️ **REVISI 2026-09-17** `[keputusan work owner]` — **`KMT_NO` DIBUANG.** `T_WORK_CLAIM` adalah tabel
**satu baris per work object**: saat kirim komite, kasus komite mendapat **barisnya sendiri**
(`ID` = identitas kasus komite) dan menunjuk klaimnya lewat **`COVER_KEY`**. Tidak ada kolom nomor
terpisah. Test yang menemukan `KMT_NO` **gagal**.

#### ⚠️ Catatan korpus — bentuk penautan ini TIDAK ADA di Pega

`[keputusan work owner]` `COVER_KEY` dan **shared primary key** adalah **bentuk baru**, bukan
temuan korpus.
Jangan mencarinya di Pega. Begitu pula `KMT_NO` yang sempat dirancang lalu dibuang.

Sensus 2026-09-17 atas **81 berkas korpus** — seluruh modul `Komite Claim Life` (47 dari 53 berkas;
sisanya `.xlsx`) dan 34 berkas `Claim Life` (**seluruh** `RDBList/`, `Section/ClaimComite.xml`,
`Harness/Committe_Life.xml`, dan aktivitas kunci) menghasilkan:

| Yang dicari | Hasil |
| --- | --- |
| string `KMT-` | **0 berkas** |
| string `KMT_NO` | **0 berkas** — sempat dirancang, dibuang 2026-09-17 |
| `KMT` lainnya | hanya di **nama rule** — `CreateKMTLife_Act`, `Generate_NoAccept_KMT_Life`, `HitServiceToKasirKMTLife_Act` |

`[terverifikasi]` **Penaut nyata di Pega bukan nomor kasus.** `Claim Life/Activity/CreateKMTLife_Act.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) menautkan kasus
komite ke klaim lewat **`pxAddChildWork`** (3×, mekanisme parent-child bawaan Pega), `pzInsKey` (7×),
`pyID` (2×), ditambah muatan `childPageKomite.CLMNO`, `.IndexAdjustment`, `.IndexPremiumList`. **Nol**
nilai literal ber-`KMT`.

Jadi rantai `KOMITE_ID` → `T_WORK_CLAIM` (`ID`/`COVER_KEY`) → `T_GENERAL_KOMITE` → `T_KOMITE_KOMITELIST`
adalah **rancangan pengganti** mekanisme parent-child Pega — sah sebagai keputusan desain, tetapi
**bukan** `[terverifikasi]`. Siapa pun yang menuliskannya ulang wajib memakai tanda
`[keputusan work owner]`, bukan `[terverifikasi]`.

Catatan: `COVER_KEY` **lebih dekat** ke bentuk Pega daripada rancangan `KMT_NO` sebelumnya, karena
Pega memang menautkan induk-anak lewat mekanisme *cover* (`pxAddChildWork`, `pzInsKey`) dan bukan
lewat kolom nomor. Tetapi nama `COVER_KEY` sendiri tetap tidak muncul di korpus.

`[terbuka]` **Nasib `CLMNO`.** `[terverifikasi]` `CLMNO` adalah nomor klaim yang **benar-benar**
dibawa ke kasus komite di Pega — muncul di `Claim Life/Activity/CreateKMTLife_Act.xml` dan
`Komite Claim Life/Activity/SendEmailKlaimLife.xml` (2 berkas korpus), dan tercatat di tiket 10
baris 33 serta tiket 01 Komite sebagai bagian muatan penyerahan. Tetapi ia **tidak muncul sama
sekali** dalam rancangan rantai relasional di atas. Apakah `CLMNO` tetap dibawa (misalnya sebagai
kolom pada `T_GENERAL_KOMITE`), atau memang digantikan mekanisme `ID`/`COVER_KEY`, **belum pernah
dinyatakan**. **Jangan tebak** — keputusan work owner.

`[terverifikasi]` Di Pega justru sebaliknya — `KomiteList` adalah **anak `AdjustmentList`**:
`Komite Claim Life/Activity/KomitePostAdjustment.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` /
`KOMITEPOSTADJUSTMENT`) menulis `…AdjustmentList(idx).KomiteList(k).KomiteAproval` /
`.KomiteComment` / `.DateApprove`; `Claim Life/Activity/CreateKMTLife_Act.xml` mengisi
`childPageKomite.KomiteList`, `IndexPremiumList`, dan `IndexAdjustment = .pxListSubscript`.

⚠️ **Penyimpangan sadar (baru) — rujukan memakai ID stabil, bukan indeks posisi.**
`[keputusan work owner]` Pega merujuk baris lewat **subscript posisi** (`IndexPremiumList`,
`IndexAdjustment`) — rusak begitu urutan bergeser. Sistem baru memakai **`ADJUSTMENT_ID`** dan
**`KOMITE_ID`**. Kelas bug yang sama sudah diperbaiki di Endorsement Life dengan `PARENT_ID`.

Daftar kolom lengkap tiap tabel ada di **spec §2b** dan
`.scratch/claim-life/revisi-penyimpanan-json-dibuang.md`.

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InsertJsonKlaimLife_sql` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!INSERTJSONKLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/InsertJsonKlaimLife_sql.xml` | ⚠️ **INSERT flat 51 kolom** ke `OS_AKSEPTASI_KLAIM_LIFE` — **bukan** JSON; sumber migrasi |
| `UpdateOsAkseptasiClaimLife_sql` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | 55 nama kolom terbaca langsung |
| `SavePesertaClaim` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEPESERTACLAIM` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/SavePesertaClaim.xml` | ⚠️ **bukti induk adjustment = peserta** |
| `SaveOutStandingLife_Act` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/SaveOutStandingLife_Act.xml` | ⚠️ bukti dokumen **per peserta** |

`[terverifikasi]` `OS_AKSEPTASI_KLAIM_LIFE` menyimpan **satu baris per peserta** dan **mengulang 14
atribut polis** pada setiap baris peserta.

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0006** (identitas lewat sequence), **ADR-0009** (migrasi
penuh; koeksistensi ditolak), **ADR-0011** (seluruh **baris** adjustment ikut pindah).

## Acceptance criteria

- [x] ⚠️ Skema klaim **relasional penuh**: setiap atribut menjadi **kolom bernama**. Test yang
      menemukan kolom JSON menyimpan atribut klaim **gagal**. *(AC 31 spec; penyimpangan sadar 1)*
- [x] ⚠️ **Enam tabel klaim** ada dengan PK dan FK sesuai diagram — `T_GENERAL_CLAIM`,
      `T_CLAIMLF_PREMIUMLIST_DETAIL`, `T_CLAIMLF_ADJUSTMENT`, `T_CLAIMLF_ADJUSTMENT_SPREADING`,
      `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`, `DOCUMENT_CLAIM`. **REVISI 2026-09-18:** ⛔ **tidak
      ada** `T_CLAIMLF_POLICY` maupun `T_CLAIMLF_MARKETING` — keduanya **dihapus**, bukan diganti
      nama. Test yang menemukan salah satunya **gagal**. *(AC 48 spec; penyimpangan sadar 8)*
- [x] ⚠️ **REVISI 2026-09-18** — **tidak** seluruh FK `ON DELETE CASCADE`. Yang **CASCADE**: relasi
      **3 · 4 · 5 · 6** (klaim → peserta → adjustment → spreading → spreading retro). Yang **di Go**:
      relasi **1** (`COVER_KEY`) dan **11** (`KOMITE_ID`, penunjuk). Yang ⚠️ `[terbuka]`: relasi
      **7** (`DOCUMENT_CLAIM`, kini lintas-lini — induk beda tabel per lini) dan relasi **2**
      (belum punya kolom). Lihat tabel sebelas relasi di atas. **Jangan menyeragamkan sendiri.**
- [x] ⚠️ FK `T_CLAIMLF_ADJUSTMENT` menunjuk **`T_CLAIMLF_PREMIUMLIST_DETAIL.ID`**, **bukan**
      `T_GENERAL_CLAIM.ID`. Test yang menemukan adjustment menggantung pada header **gagal**.
      *(AC 33 spec; penyimpangan sadar 2)*
- [x] ⚠️ FK `DOCUMENT_CLAIM` menunjuk **`T_CLAIMLF_PREMIUMLIST_DETAIL.ID`**. *(AC 44 spec;
      penyimpangan sadar 5)*
- [x] ⚠️ **`T_WORK_CLAIM` ada sebagai tabel mandiri** — bukan anak `T_GENERAL_CLAIM`, dan keadaan tangga
      **tidak** menjadi kolom header klaim. *(AC 46 spec; penyimpangan sadar 6)*
- [x] Header memuat keempat field `PremiumListSummary` (`CLAIM_NO`, `PL_NUMBER`, `RISLIPRNM`,
      `BUSINESS_NAME`) beserta `CASEID` dan `CLAIM_RETRO`. *(AC 36 spec)*
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS, MENUNGGU JAWABAN.**
      `T_CLAIM_POLICY` **dihapus**, jadi tidak ada tabel yang memuat kesembilan kolom itu. **27 dari
      32** kolomnya dibaca dari `T_PREMIUM_LIST`; **lima** ⚠️ `[terbuka]` **kehilangan rumah** —
      `TEAM_GROUP`, `BUSINESS_ID`, `TANGGAL_RESPON`, `TANGGAL_REALISASI`,
      `TANGGAL_KONFIRMASI_BALIK` (§Blocker). `PL_NUMBER` tetap tidak ada, tetapi karena **diganti
      nama menjadi `POLICY_NO`** di `T_GENERAL_CLAIM` — bukan karena tabel ini menolaknya.
      **Jangan tebak rumah kelima kolom itu.** *(AC 37 spec)*
      > `[keputusan work owner 26-09-2026, butir e′]` **Tiga dari lima kolom itu sudah ketemu
      > rumahnya, dan rumahnya bukan Claim Life:** `TANGGAL_RESPON`, `TANGGAL_KONFIRMASI_BALIK`, dan
      > `TANGGAL_REALISASI` adalah **tanggal PENAWARAN**, milik modul **PremiumList Life**. Penanda
      > mereka karena itu **`[terbuka — PremiumList Life]`**, bukan `[terbuka]` tiket ini, dan berkas
      > migrasi `002` **tidak disentuh** — nol kolom `TANGGAL_*` di DDL mana pun.
      >
      > `[terverifikasi]` Rantai buktinya, diperiksa **tiga kali dengan dua alat berbeda**:
      > diketik manusia di `PremiumList Life/Section/InputOfferLife.xml` (label *Tanggal Respon*,
      > *Tanggal Konfirmasi Balik*, *Tanggal Realisasi*) → disimpan `SaveOfferJsonLife_SQL.xml`
      > → mendarat di `POOLDATA.JSON_OFFER_LIFE` sebagai `RESPONSE_DATE`, `RECONFIRMATION_DATE`,
      > `REALIZATION_DATE` → dibaca hidup `GetOfferLife_sql.xml`.
      >
      > ⭐ Yang memutuskan: ketiga **nama kolom** itu muncul **nol kali** di seluruh modul Claim;
      > ketiga **propertinya** muncul di Claim Life hanya di **4 berkas `Section\`** (menampilkan),
      > dan **nol berkas `Activity\`** (menulis). Pega sendiri karena itu **membaca dan menampilkan**
      > ketiganya di layar Claim Life, tetapi **tidak pernah menyimpannya** sebagai data klaim.
      > Membuat kolomnya di skema klaim akan menyalin data milik modul lain.
      >
      > ⚠️ **Dua kolom sisanya — `TEAM_GROUP` dan `BUSINESS_ID` — tetap `[terbuka]` penuh.**
      > Keputusan ini tidak menyentuh keduanya, dan AC 8 karena itu **tetap tidak dicentang**.
- [x] `T_CLAIMLF_PREMIUMLIST_DETAIL` memuat **kesembilan kolom tambahan** hasil audit — termasuk
      `IS_CHECK` dan ketiga tanggal per peserta. *(AC 39, 40, 41 spec)*
- [x] `T_CLAIMLF_ADJUSTMENT` memuat `CLAIM_AMOUNT`, `STS_REJECT`, `ACCEPTEDNO`, `ACCEPTATION_DATE`.
      *(AC 42, 43 spec)*
- [x] `T_CLAIMLF_ADJUSTMENT` memuat **ketiga kolom bank** — `NAME_OF_BANK`, `ID_BANK`, `ACCOUNT_NO` —
      dan migrasi mengisinya dari kolom warisan `NAME_OF_BANK`, `IDBANK`, `ACCOUNTNO`.
      *(AC 56 spec; `[terverifikasi]`)*
- [x] ⚠️ Seluruh uang dan share bertipe **desimal presisi arbitrer**; seluruh tanggal **`DATE`**;
      seluruh kolom **nullable**; identitas dari **sequence**. Test yang menemukan kolom uang
      bertipe teks atau melewati `float` **gagal**. *(AC 50 spec; **ADR-0003**, **ADR-0006**;
      penyimpangan sadar 7)*
      > `[keputusan work owner 26-09-2026, butir t]` "Presisi arbitrer" dibaca sebagai **jenis data
      > desimal non-float, lawan dari `float`** — bacaan yang sejalan dengan rujukan AC ini sendiri
      > (ADR-U-0003) dan dengan klausa ujinya sendiri di kalimat di atas. `NUMBER(38,8)` ↔
      > `apd.Decimal` memenuhinya. Klausa **sequence** tercakup penyimpangan sadar AC 34 (nomor
      > bisnis berformat untuk `T_WORK_CLAIM` dan kedua tabel ber-shared-PK; ADR-U-0006 tetap
      > berlaku untuk `T_CLAIMLF_*` dan `DOCUMENT_CLAIM`). Ronde 3 membacanya sebagai "`NUMBER`
      > tanpa presisi" lalu menahannya — bacaan itu diralat di sini.
- [x] Seluruh klaim Life terbawa **beserta seluruh baris adjustment**-nya — bukan hanya keadaan
      terakhir. Jumlah baris per klaim setelah migrasi **sama** dengan sebelumnya. *(**ADR-0011**)*
- [x] Setiap baris adjustment hasil migrasi **menunjuk peserta yang benar**. *(AC 34 spec)*
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS.** Tidak ada
      `T_CLAIM_POLICY` untuk dinormalkan ke dalamnya; atribut polis **tidak dipindahkan ke mana
      pun** — ia **dibaca hidup** dari tabel polis ⚠️ **penyimpangan sadar**. Yang **tetap
      mengikat**: bila atribut polis **berbeda antar baris peserta** dalam satu klaim, migrasi
      **melaporkannya** dan **tidak** diam-diam memilih salah satu. *(AC 51 spec)*
- [x] ⚠️ Migrasi **melaporkan** bahwa data lama ber-`ACCEPTATION_DATE` = waktu insert dan
      `STS_REJECT` = `0` karena **di-hardcode** di sumbernya, bukan karena nilainya sebenarnya.
      Migrasi **tidak mengarang** tanggal akseptasi. *(AC 52 spec; penyimpangan sadar 4)*
- [x] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
      tepat**, bukan dengan toleransi. *(AC 53 spec; **ADR-0003**)*
- [x] Tanggal yang berupa teks menjadi `DATE` **tanpa pergeseran zona waktu**; yang **tidak dapat
      diurai dilaporkan**, bukan didiamkan. *(AC 54 spec)*
- [ ] ⚠️ Migrasi **tidak mengambil apa pun** dari `m_product_life.JSONDATA`; `PRODUCT_NAME` dan
      `PRODUCT_NAME_ID` berasal dari **`product_life` relasional**. *(AC 38 spec)*
- [ ] Penomoran klaim **tidak melompat dan tidak mengulang** setelah migrasi; sequence pindah dengan
      **nilai berjalan yang benar**. *(**ADR-0006**)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

### Kolom yang bertambah dan yang pindah ⚠️ BARU 2026-09-18

- [x] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `CREATE_OP`, `CREATE_OP_NAME`, `TGL_UPDATE` —
      ketiganya **pindah ke `T_WORK_CLAIM`**. Test yang menemukannya di header klaim **gagal**.
- [x] ⚠️ `T_WORK_CLAIM` memuat `CREATE_OP`, `CREATE_OP_NAME`, `TGL_UPDATE`, **`CASEID`**
      (seluruhnya pindahan) dan `LINI`.
- [x] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `CASEID` — pindah ke `T_WORK_CLAIM`. Test yang
      menemukannya di header klaim **gagal**. *(DIPUTUSKAN 2026-09-18; `[keputusan work owner]`)*
- [x] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `PL_NUMBER`; kolom itu **diganti nama** menjadi
      **`POLICY_NO`**. Test yang menemukan `PL_NUMBER` di header klaim **gagal**.
- [x] ⚠️ `T_GENERAL_CLAIM` memuat ketiga **penunjuk polis** — `CASEID_POLICY`, `POLICY_NO`,
      `ENDORSMENT_NO` — dan atribut polis **tidak** disalin ke klaim. Test yang menemukan salinan
      atribut polis di tabel klaim **gagal**.
- [ ] ⚠️ `[terbuka]` **Isi ketiga penunjuk polis belum dapat ditulis** — tidak satu pun menunjuk
      kolom yang ada di sisi polis (§Blocker). Kolomnya dibuat; **pengisiannya** menunggu keputusan
      work owner. **Jangan tebak.**
- [x] ⚠️ Kolom `LINI` **ADA** pada `T_WORK_CLAIM`, dan untuk Life isinya **konstanta lini Life**.
      `[terbuka — Non-Life]` hanya **daftar nilai enum lintas-lini**, yang ditetapkan saat konteks
      Non-Life digarap — **bukan** keberadaan kolomnya. *(DIPUTUSKAN 2026-09-18)*

### Identitas: shared PK dan nomor bisnis berformat ⚠️ BARU 2026-09-18

- [x] ⚠️ **`T_GENERAL_CLAIM.ID` sama persis dengan `T_WORK_CLAIM.ID` baris klaim** — **shared
      primary key**, 1:1, **tanpa kolom penyambung**. `T_GENERAL_CLAIM.ID` sekaligus PK **dan** FK
      ke `T_WORK_CLAIM.ID`. *(relasi 2; `[keputusan work owner]`)*
- [x] ⚠️ ⛔ **Tidak ada kolom `WORK_CLAIM_ID`** di mana pun — pada `T_GENERAL_KOMITE` maupun
      tabel lain. Hubungan ke kasus komite juga **shared PK**: `T_GENERAL_KOMITE.ID` =
      `T_WORK_CLAIM.ID` baris komite. Test yang menemukan kolom `WORK_CLAIM_ID` **gagal**.
      *(relasi 8; `[keputusan work owner]`)*
- [ ] ⚠️ **`T_GENERAL_KOMITE.ADJUSTMENT_ID` TETAP ADA** — penutup lingkar ke baris adjustment; ia
      **bukan** bagian shared PK dan **tidak** ikut dibuang. *(relasi 10)*
- [x] ⚠️ **`T_WORK_CLAIM.ID` bertipe teks berformat** — baris klaim `CLM-xxxxxx`, baris komite
      `KMT-xxxxxx`. **Bukan angka sequence.** Test yang menemukan tipe numerik **gagal**.
      *(`[keputusan work owner]`)*
- [ ] ⚠️ `T_WORK_CLAIM.COVER_KEY`, `T_GENERAL_CLAIM.ID`, `T_GENERAL_KOMITE.ID`, dan
      `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` **bertipe sama** dengan `T_WORK_CLAIM.ID` — teks berformat.
- [x] ⚠️ **Penyimpangan sadar dari ADR-0006 dicatat di artefak**, tidak dilanggar diam-diam:
      identitas `T_WORK_CLAIM` dan kedua tabel ber-shared-PK adalah **nomor bisnis berformat**,
      bukan sequence. ADR-0006 **tetap berlaku** untuk `T_CLAIMLF_*` dan `DOCUMENT_CLAIM`.
- [ ] ⚠️ `[terbuka]` **Generator nomor `CLM-`/`KMT-` belum ditetapkan** — siapa yang membuatnya,
      apakah ada sequence di belakang prefiks, apakah di-reset per tahun. Pemilik **DBA / work
      owner**. Kolomnya dibuat; **pembangkitannya** menunggu jawaban. **Jangan tebak.**

### Spreading adjustment + rujukan Komite ⚠️ BARU 2026-09-16

- [x] ⚠️ **`T_CLAIMLF_ADJUSTMENT_SPREADING` ada**, dengan FK **`ADJUSTMENT_ID`** → `T_CLAIMLF_ADJUSTMENT.ID`
      dan **`ON DELETE CASCADE`**. Test yang menemukannya menggantung pada peserta atau pada header
      klaim **gagal**. *(`[terverifikasi]` `SpreadingClaimLife_Act` mengisi `.SpreadingList` pada
      baris adjustment; **AC 58 spec**)*
- [x] ⚠️ **`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` ada**, dengan FK **`SPREADING_ID`** →
      `T_CLAIMLF_ADJUSTMENT_SPREADING.ID` dan **`ON DELETE CASCADE`**. *(AC 58 spec)*
- [ ] Pohon klaim berkedalaman **lima tingkat** — klaim → peserta → adjustment → spreading →
      spreading retro — dan menghapus klaim **mengkaskade sampai tingkat terdalam**. Test wajib
      memeriksa **cicit** (`_SPREADING_RETRO`) ikut hilang. *(AC 60 spec)*
- [x] `T_CLAIMLF_ADJUSTMENT_SPREADING` memuat `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`,
      `TREATY_YEAR_LIFE`, `RETROCADED_SHARE`, `RATE`, `IDR`, `USD`, `CURRENCY`.
- [x] `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` memuat `REINSURER_NAME`, `PERCENT_SHARE`, `AMOUNT`,
      `RATE`, `PREMIUM_SPREADED_GROSS`, `PREMIUM_SPREADED_NET`, `COMMISION`, `OVR_COMM`,
      `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`. *(`[terverifikasi]`
      `Claim Life/Section/RetroDetailClaimLife.xml` + `SpreadingClaimLife_Act`)*
- [x] ⚠️ Seluruh kolom uang dan persen pada kedua tabel bertipe **desimal presisi arbitrer**;
      **tidak** melewati `float`. *(**ADR-0003**)*
- [x] ⚠️ Nilai spreading **dibekukan**: perubahan master treaty setelah adjustment tersimpan
      **tidak mengubah** angka yang sudah ada. *(penyimpangan sadar — spreading disimpan)*
- [x] ⚠️ **`T_CLAIMLF_ADJUSTMENT` memuat kolom `KOMITE_ID`**, **nullable** dan **ber-index** — `NULL`
      bila baris belum pernah dikirim ke Komite. *(AC 61 spec; penyimpangan sadar — rujukan, bukan salinan)*
- [x] ⚠️ **Tidak ada tabel `T_CLAIMLF_ADJUSTMENT_KOMITE`.** Roster dan keputusan per anggota **tidak**
      disimpan di Claim Life. Test yang menemukan tabel itu **gagal**. *(AC 61 spec; penyimpangan sadar)*
- [x] ⚠️ Rujukan ke Komite memakai **`KOMITE_ID`**, bukan **indeks posisi**. Test yang menemukan
      padanan `IndexPremiumList`/`IndexAdjustment` sebagai kunci rujukan **gagal**.
      *(AC 62 spec; `[terverifikasi]` `CreateKMTLife_Act` memakai `.pxListSubscript`; penyimpangan sadar)*
- [ ] Migrasi **membongkar** `SpreadingList` dan `RetroLifeList` dari data lama ke kedua tabel,
      dan setiap baris dapat ditelusuri ke **baris adjustment yang benar**.
- [x] Setiap FK baru (`ADJUSTMENT_ID`, `SPREADING_ID`) **ber-index**.
- [x] ⚠️ `KOMITE_ID` **ber-index** juga, dan berisi **identitas kasus komite** = `T_WORK_CLAIM.ID`
      baris komite. Test yang menemukannya menunjuk `T_GENERAL_KOMITE.ID` **gagal**.
      *(§`KOMITE_ID` pada `T_CLAIMLF_ADJUSTMENT`; `[keputusan work owner]` REVISI 2026-09-17)*
- [x] ✅ Tipe `KOMITE_ID` **SUDAH DITETAPKAN 2026-09-18** — teks berformat `KMT-xxxxxx`, mengikuti
      `T_WORK_CLAIM.ID`. *(`[keputusan work owner]`)*
- [x] ✅ `[DIJAWAB 26-09-2026]` **Apakah `KOMITE_ID` dan `COVER_KEY` dipasangi
      `REFERENCES T_WORK_CLAIM(ID)`** atau dibiarkan tanpa constraint. Pemilik **DBA / work owner**.
      ~~`[terbuka]` — belum diputuskan; tiket ini tidak dinyatakan selesai sebelum jawabannya ada.
      Jangan tebak.~~ *(kalimat ini dicoret, bukan dihapus: ia benar sampai 26-09-2026 dan menjadi
      alasan tiket tertahan empat sesi)*
      > ⭐ `[keputusan work owner 26-09-2026, butir d]` sesudah **empat sesi** `[USULAN]`:
      > **DIPASANGI `REFERENCES`, tetap nullable.** Diterapkan `001` (`FK_WORK_COVER_KEY`, menunjuk
      > tabelnya sendiri) dan `004` (`FK_ADJ_KOMITE`). Keduanya **tanpa `ON DELETE`** — ini penunjuk
      > KE ATAS, bukan kepemilikan, sehingga kaskade akan terbalik arah. Dikunci
      > `TestAC50PenunjukKeAtasBerReferences`, yang juga menolak `ON DELETE` dan `NOT NULL`.
      >
      > ⚠️ **`[terbuka]` baru yang lahir dari keputusan ini — pemiliknya tiket 15, bukan tiket
      > ini.** Model di bab TINGKAT 1 menyatakan baris KOMITE memuat `COVER_KEY = ID baris klaim`.
      > Dengan FK tanpa `ON DELETE`, menghapus baris klaim yang masih ditunjuk baris komite akan
      > **ditolak `ORA-02292`**. `PohonKlaim.Hapus` menghapus baris work paling akhir, jadi urutannya
      > sendiri selamat, tetapi baris komite tidak ikut dibereskannya. Yang harus diputuskan: jalur
      > hapus melepas penunjuknya lebih dulu, atau penghapusan memang ditolak selama komite ada.
      > ⛔ **Belum teruji**: migrasi belum pernah berjalan di Oracle mana pun, jadi ini penalaran
      > atas semantik Oracle, bukan hasil pengamatan.

### Koreksi: `OS_AKSEPTASI_KLAIM_LIFE` TETAP ditulis ⚠️ 2026-09-16

⚠️ `[keputusan work owner]` **Keputusan sebelumnya dibalik.** Sistem baru menulis **dua tempat**:
tabel relasional baru **dan** `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE`, karena hilir (Arasapas /
produksi) masih membaca dari sana. **Yang dibuang hanya JSON** (`JSON_KLAIM` / serialisasi
`ClaimData`).

- [ ] ⚠️ Menyimpan klaim menulis **tabel relasional baru** *dan* `INSERT` flat ke
      `OS_AKSEPTASI_KLAIM_LIFE`. Test yang menuntut `OS_AKSEPTASI_KLAIM_LIFE` **tidak** ditulis
      adalah **keliru** dan harus dibalik.
- [x] ⚠️ Yang **tetap dibuang** hanya **blob JSON** — `JSON_KLAIM` dan serialisasi `ClaimData`.
      *(AC 31, 34 spec tetap berlaku untuk JSON saja)*
- [ ] Kedua penulisan berada dalam **satu transaksi**; kegagalan pada salah satunya **membatalkan
      keduanya**. *(AC 49 spec)*

## Blocker

✅ **PEMBLOKIR RELASI 2 DICABUT 2026-09-18.** `[keputusan work owner]` Relasi
`T_WORK_CLAIM` → `T_GENERAL_CLAIM` **dipecahkan lewat SHARED PRIMARY KEY**: `T_GENERAL_CLAIM.ID`
**sama persis** dengan `T_WORK_CLAIM.ID` baris klaim. Tidak perlu kolom penyambung, dan
**`WORK_CLAIM_ID` dibuang** dari seluruh rancangan. Relasi 8 dipecahkan dengan cara yang sama.

> Teks lama: *"⛔ PEMBLOKIR — relasi nomor 2 tidak punya kolom di sisi mana pun … Tiket ini tidak
> dinyatakan selesai sebelum dijawab."* — **dicabut.** Relasi nomor 2 **bukan lagi** alasan tiket
> ini belum selesai.

⚠️ **Tiket ini MASIH belum dapat dinyatakan selesai**, tetapi karena **satu** hal lain:
⚠️ `[terbuka]` **apakah `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)`**
atau dibiarkan tanpa constraint — pemilik **DBA / work owner**. Tipe kolomnya **sudah** ditetapkan
(teks berformat); yang belum hanya **constraint**-nya.

(Yang tetap benar: **OQ-001 dan OQ-002 (bagian penyimpanan) ditutup 2026-09-16** — tabel klaim
dirancang sendiri, dan jalur tulis klaim tidak lagi lewat stored procedure.)

### ⚠️ `[terbuka]` — lima kolom kehilangan rumah

**Status pemblokirnya BELUM ditetapkan work owner.** Keputusan 2026-09-18 menyatakan kelima kolom ini
`[terbuka]` dan **tidak** menyatakannya memblokir. Ini 5 dari 32 kolom `T_CLAIM_POLICY` yang
**tidak** tersedia di `T_PREMIUM_LIST`. **JANGAN DITEBAK:**

| Kolom | Keadaan |
| --- | --- |
| `TEAM_GROUP` | Tidak ada di `T_PREMIUM_LIST`; **ada di master marketing officer** — diambil lewat `MO_ID` |
| `BUSINESS_ID` | Tidak ada di `T_PREMIUM_LIST`; dipakai sebagai **parameter** `Claim Life/RDBList/Generate_NoAccept_Life.xml`. **Sumbernya wajib ditetapkan** |
| `TANGGAL_RESPON` · `TANGGAL_REALISASI` · `TANGGAL_KONFIRMASI_BALIK` | **Tanggal proses klaim, bukan atribut polis.** Tidak ada di `T_PREMIUM_LIST`, dan **tidak lagi ada di klaim**. Rumahnya `T_GENERAL_CLAIM` atau `T_WORK_CLAIM` — **belum diputuskan** |

⚠️ **Usulan — bukan keputusan.** Ketiga `TANGGAL_*` bukan perkara kosmetik: `[terverifikasi]`
ketiganya **tampil di empat section aktif** (`InputRegisterClaimLife`, `InputOSClaimLife`,
`InputAkseptasiClaimLife`, `MedicalCheckClaimLife`) **tanpa gerbang visibilitas** — layar akan
kehilangan isi bila rumahnya tidak ditetapkan. Atas dasar itu **saya usulkan** butir ini ikut
memblokir tiket 14, tetapi **penetapannya keputusan work owner** dan **belum diambil**.

### ⚠️ `[terbuka]` tambahan — tidak memblokir pembuatan tabel, tetapi wajib dijawab

- **Isi ketiga penunjuk polis** (`CASEID_POLICY`, `POLICY_NO`, `ENDORSMENT_NO`) — tidak satu pun
  menunjuk kolom yang ada di sisi polis. **TETAP terbuka.**
- **Generator nomor `CLM-`/`KMT-`** — siapa yang membuatnya, sequence di belakang prefiks atau
  tidak, reset per tahun atau tidak. Pemilik **DBA / work owner**.

✅ **Ditutup 2026-09-18, tidak lagi terbuka:** relasi 2 tanpa kolom (shared PK) · relasi 8
memakai `WORK_CLAIM_ID` (shared PK, kolom dibuang) · nasib `CASEID` (pindah ke `T_WORK_CLAIM`) ·
`ON DELETE` `DOCUMENT_CLAIM` (di Go) · tipe `T_WORK_CLAIM.ID`/`KOMITE_ID` (teks berformat) ·
keberadaan kolom `LINI` (ada; hanya daftar enum lintas-lini yang menunggu Non-Life).

⚠️ `[terbuka — Non-Life]` **tidak memblokir:** relasi & cascade formal `T_WORK_CLAIM` ditetapkan
saat konteks Non-Life digarap. Yang mengikat di sini hanya: klaim life dihapus → baris work-nya ikut.

⚠️ `[terbuka]` **tidak memblokir skema, memblokir perilaku hitung:** **`PREMIUM_SPREADED_NET` punya
dua rumus** di rule yang sama — `GROSS − Comm` dan `GROSS − Discount − Comm`
(`Claim Life/Activity/SpreadingClaimLife_Act.xml`). **Cabang mana yang berlaku tidak terbaca dari
korpus.** Kolomnya tetap dibuat di tiket ini; **rumusnya** ditetapkan Product + UW sebelum tiket
**03** dinyatakan selesai. **Jangan tebak.**

⚠️ `[terbuka — selaraskan Komite Claim Life]` **tidak memblokir tiket ini:** tabel komite di konteks
**Komite Claim Life** perlu menampung keputusan **per baris adjustment** — kolom `ADJUSTMENT_ID`
(FK → `T_CLAIMLF_ADJUSTMENT.ID`) beserta `KOMITE_APROVAL`, `KOMITE_COMMENT`, `DATE_APPROVE`, dan
roster. Apakah tabel di sana sudah punya kolom itu **belum diperiksa**. Penyelarasannya pekerjaan
konteks Komite, bukan di sini.

## Catatan

⚠️ **Sumber migrasi kedua bernama "Update" tetapi isinya INSERT.** `[terverifikasi]`
`Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL`) berisi
**`INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` dengan 55 kolom** — **tidak ada `UPDATE` di
dalamnya**. Dicatat supaya **tidak ada waktu terbuang mencari jalur `UPDATE` yang memang tidak
ada**. Ketiga kolom bank berasal dari rule ini (**OQ-066**).

⚠️ **Nama berbohong — sepasang, saling membalik.** `[terverifikasi]` `InsertJsonKlaimLife_sql`
bernama "Json" tetapi **INSERT flat 51 kolom**; sebaliknya `GetJsonProductLife`
(`ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!GETJSONPRODUCTLIFE`) bernama "Json" dan **memang** JSON
(`json_table` atas `m_product_life.JSONDATA`). Nama tidak memberi tahu apa-apa di kedua arah —
**baca kodenya** (**OQ-066**).

⚠️ **`DOCUMENT_CLAIM` kolomnya belum final.** Rinciannya (nama berkas / kategori / kelengkapan)
diturunkan dari sensus `.DocumentList` di `SaveOutStandingLife_Act` saat tiket **03** dikerjakan.
Tiket ini menyiapkan **tabel dan kuncinya**; kolom isian ditetapkan di sana — **jangan tebak dari
nama tabel**.

⚠️ **`T_WORK_CLAIM` kolomnya belum final** — bentuk lengkapnya menyusul saat penyelarasan
lintas-lini. Untuk Claim Life yang wajib ada: posisi tangga, status akseptasi, penanda pengembalian
ke Admin dan ke Medical, serta `Type`.

## Seam & perintah verifikasi

**Seam: API HTTP Claim — Life** terhadap **skema uji Oracle nyata** — kaskade tiga tingkat, presisi
desimal, konversi tanggal, dan normalisasi atribut polis **hanya berperilaku benar pada basis data
sungguhan**; memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Implementasi — 25 September 2026

**Status: `claimed`** — **25 dari 53 AC tertutup**, diverifikasi **34 test Go yang benar-benar
berjalan**. Sesi ini berjalan di **jalur B**: `ORACLE_DSN` kosong, sehingga seluruh test bertag
`db` **melewati dengan pesan**, bukan lulus diam-diam.

⛔ **Tiket ini tidak dapat `resolved`** walau seluruh kodenya jadi dan hijau, sebab AC nomor 50
masih `[terbuka]`: apakah `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)`
belum diputuskan. Bab Blocker tiket ini menyatakannya sendiri.

### Berkas yang dibuat

| Berkas | Total | Berisi |
| --- | ---: | ---: |
| `internal/repository/migrations/` — 16 berkas `.sql` | 425 | — |
| `internal/repository/migrasi.go` | 306 | 285 |
| `internal/repository/migrasi_test.go` | 304 | 285 |
| `internal/repository/migrasidata.go` | 324 | 301 |
| `internal/repository/migrasidata_test.go` | 287 | 261 |
| `internal/repository/pohonklaim.go` | 397 | 369 |
| `internal/repository/pohonklaim_db_test.go` | 259 | 236 |
| `internal/models/pohonklaim.go` | 143 | 132 |
| **jumlah Go** | **2.020** | **1.869** |

Diubah: `internal/repository/repository.go` (penanda `IS_PEGA_PROD` kini punya pembaca),
`internal/repository/skemauji/skemauji.go` (memakai migrasi, bukan DDL tangan),
`internal/models/klaimlife.go` (`BarisAdjustment` membawa `Spreading`),
`internal/services/klaimlife.go` dan `cmd/api/main.go` (`-migrate` menjalankan migrasi
sungguhan), `README-BACA-DULU.md` bab 1 dan 5.

### Yang dijalankan, dan hasilnya

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` | lulus |
| `go vet -tags=db ./...` | lulus — test Oracle ikut kompilasi |
| `gofmt -l` | nol berkas |
| `go test ./...` | lulus, **34 test** |
| `go test -tags=db ./internal/...` | **MELEWATI** dengan `ORACLE_DSN belum dikonfigurasi` |
| `npm run typecheck` · `npm test` | lulus, 5 test — frontend tidak diubah tiket ini |

⭐ **Instrumennya sendiri diuji.** Nama terlarang `WORK_CLAIM_ID` sengaja disisipkan ke berkas
migrasi; `TestNamaYangDibuangTidakAda` **gagal** seperti seharusnya, lalu berkasnya dipulihkan.
Tanpa langkah itu, "seluruh test lulus" tidak membuktikan apa-apa.

### AC yang DITUTUP — 25, dengan test yang berjalan

**Bentuk tabel** (nomor 1, 2, 3, 4, 6, 22, 24, 25, 26, 29, 30, 32, 34, 36, 37, 43, 44, 47, 49, 52):
tujuh tabel ada dengan PK dan FK sesuai diagram; kaskade **hanya** pada relasi 3·4·5·6 dan
`DOCUMENT_CLAIM` tanpa kaskade; FK adjustment menunjuk **peserta**, bukan header; shared primary
key terbentuk (`T_GENERAL_CLAIM.ID` sekaligus PK dan FK); `T_WORK_CLAIM.ID` bertipe teks; setiap
kunci tamu ber-index dan `KOMITE_ID` ber-index **UNIK**; header sudah tidak memuat `CASEID`,
`CREATE_OP`, `CREATE_OP_NAME`, `TGL_UPDATE`, maupun `PL_NUMBER`, dan memuat ketiga penunjuk polis;
nol kolom JSON, CLOB, maupun float; nol `T_CLAIMLF_POLICY`, `T_CLAIMLF_MARKETING`,
`WORK_CLAIM_ID`, `KMT_NO`, dan `T_CLAIMLF_ADJUSTMENT_KOMITE`; penyimpangan sadar dari ADR-U-0006
dicatat di komentar berkas migrasi.

**Pembongkaran data lama** (nomor 13, 14, 16, 17, 18, dan bagian mengikat nomor 15): seluruh baris
adjustment ikut pindah — dua baris untuk peserta yang sama menghasilkan **dua** baris, bukan satu
yang menimpa; jumlah baris sesudah sama dengan sebelumnya; setiap baris menunjuk peserta yang
benar; uang pindah **tanpa berubah satu digit** dan dibandingkan **tepat**, termasuk
`1234567890.12345678` dan `0.00000001`; tanggal teks menjadi `DATE` tanpa pergeseran zona, dan yang
tidak terurai **dilaporkan** lalu dibiarkan kosong — bukan ditebak, bukan menjadi nol; atribut
tingkat klaim yang **berbeda** antar baris peserta dilaporkan **beserta kedua nilainya**, tidak
diam-diam dipilih salah satu; nilai yang berasal dari **hardcode** sumber (`ACCEPTATION_DATE` =
waktu insert, `STS_REJECT` = `0`) dilaporkan apa adanya dan tanggal akseptasi **tidak dikarang**.

### AC yang BELUM ditutup — 28, dan sebabnya

**Menunggu instance Oracle** (nomor 20, 21, 38, 51, 53): migrasi idempoten, jalur mundur,
kaskade sampai **cicit**, penulisan ganda ke tabel relasional **dan** baris datar warisan, serta
keduanya dalam **satu transaksi**. Kodenya lengkap dan kompilasi; test-nya ditulis persis menguji
hal itu di `pohonklaim_db_test.go`, dan ia **melewati** karena tidak ada instance.

**Menunggu keputusan pemilik** (nomor 8, 15, 27, 35, 50): lima AC bertanda `[terbuka]` yang
executor **tidak menutupnya**. Kolomnya dibuat; isi, pembangkit, dan constraint-nya menunggu.

**Tertulis tetapi belum punya test** (nomor 5, 7, 9, 10, 23, 28, 31, 33, 39, 40, 41, 42, 45, 48):
kolom dan relasinya ada di berkas migrasi, tetapi belum ada pernyataan test yang menguncinya.
Dilaporkan sebagai belum ditutup, bukan dihitung lulus.

**Belum dikerjakan** (nomor 11, 19, 46):
- **11** — ketiga kolom bank ada di tabel, tetapi **pembongkar belum mengisinya** dari nama warisan
  `NAME_OF_BANK`, `IDBANK`, `ACCOUNTNO`.
- **19** — `PRODUCT_NAME` dan `PRODUCT_NAME_ID` dari `product_life` relasional belum diambil.
  Bagian larangannya terpenuhi secara pasif: migrasi tidak menyentuh `m_product_life` sama sekali.
- **46** — pembongkaran `SpreadingList` dan `RetroLifeList` dari data lama **tidak dapat
  dikerjakan dari sumber ini**: tabel datar warisan `OS_AKSEPTASI_KLAIM_LIFE` tidak memuat
  keduanya. ⚠️ Sumbernya adalah blob JSON yang justru dibuang, sehingga jalur bacanya perlu
  ditetapkan lebih dulu — dicatat sebagai temuan, bukan ditebak.

### Keputusan brief sesi yang diikuti, dan yang sengaja tidak

- **§2 b berlaku** — kolom uang `NUMBER(38,8)`. `NUMBER(38,20)` yang ditemukan sesi sebelumnya
  memang berasal dari dokumen `dastin\` modul Claim Non-Prop, bukan dari korpus Pega. Butir ini
  **menutup** pertanyaan terbuka yang diangkat sesi tiket 01.
- **§2 c masih `[USULAN]`** — karena itu kolom share, persen, dan rate **tidak** diberi
  `NUMBER(38,8)`; ia memakai `NUMBER` berpresisi arbitrer. Kolomnya dibuat, keputusannya
  dibiarkan terbuka, persis seperti §1 butir 4 memerintahkan.
- **§2 d masih `[USULAN]`** — constraint `REFERENCES` untuk `KOMITE_ID` dan `COVER_KEY`
  **tidak dipasang**, dan itulah sebab tiket ini berakhir `claimed`.
- **§2 f, g, h berlaku** — migrasi tinggal di `internal/repository/migrations/`, ditanam lewat
  `//go:embed`, dijalankan `go run ./cmd/api -migrate`; skema uji memanggil migrasi itu, bukan DDL
  tangan; data lama diuji lewat tabel **tiruan** `OS_AKSEPTASI_KLAIM_LIFE` 55 kolom.

### Code review — lima cacat keras diperbaiki sebelum commit

Poros standar atas titik tetap `a60d3ba`. Seluruhnya diperbaiki; cacah test naik 32 menjadi 34.

1. ⛔ **BOM UTF-8 di `004_t_claimlf_adjustment.sql`.** Byte `EF BB BF` di awal berkas membuat baris
   komentar pertama **lolos menjadi bagian pernyataan SQL**, dan Oracle akan menolaknya - sementara
   seluruh test tanpa basis data tetap hijau. ⚠️ Saya sendiri yang menanamnya: berkas itu dipulihkan
   dengan `Set-Content -Encoding utf8` PowerShell saat menguji instrumen, dan PowerShell 5.1 menulis
   BOM. BOM-nya dibuang, `gabung()` dibuat kebal terhadapnya, dan **dua test penjaga** ditambahkan -
   nol BOM di berkas migrasi, dan setiap pernyataan harus mulai dengan kata perintah SQL.
2. ⛔ **`Hapus` mengklaim membuang `DOCUMENT_CLAIM`, tetapi tidak melakukannya.** Kunci tamunya
   sengaja tanpa `ON DELETE`, sehingga menghapus header selagi ada baris dokumen akan ditolak
   Oracle (ORA-02292). Kini dokumen milik seluruh peserta klaim dibuang lebih dulu, dan urutannya
   dijelaskan di komentar.
3. **`PeriksaSQL` bolong** pada empat pernyataan `migrasi.go` yang benar-benar sampai ke Oracle,
   padahal kepala berkasnya menyatakan penjagaan itu. Ditambahkan.
4. **`BongkarMigrasi` tidak pernah membaca `T_MIGRASI`**, sehingga laporannya menyebut langkah yang
   tidak berbuat apa-apa sebagai "dijalankan". Kini ia membaca catatan dan melewati langkah yang
   memang belum pernah dijalankan.
5. ⛔ **Jalur migrasi data lama buntu, dan ADR-U-0006 belum hidup di kode.** `BongkarBarisLama` tidak
   pernah mengisi identitas, sedangkan `Simpan` mem-bind-nya ke kolom `NOT NULL`; sementara kelima
   sequence dibuat tetapi **nol pembaca** - `NEXTVAL` tidak muncul di mana pun. Kini `Simpan`
   mengambil nomor dari sequence untuk tingkat `T_CLAIMLF_*`, memakai shared PK untuk header, dan
   **menolak dengan terang** lewat `ErrIdentitasBelumAda` bila nomor akar `CLM-` belum ada -
   sebab pembangkitnya masih `[terbuka]` dan mengarangnya berarti menetapkan yang belum diputuskan.

Temuan penilaian yang **tidak** diubah: `BarisLama` 55 medan teks, tujuh blok `Qualify` kembar,
duplikasi pencocokan kode ORA, dan heuristik `BarisHardcode` yang mencampur "0" hardcode dengan "0"
yang sah. Ketiganya dicatat di sini, bukan didiamkan; yang terakhir perlu keputusan pemilik sebelum
dipertajam.

### Catatan

1. ⚠️ **`DOCUMENT_CLAIM` dibuat tanpa kolom isi.** `[data DBA]` Daftar kolomnya tidak dapat
   diturunkan — kelas Pega-nya `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`, SQL-nya dibuat Pega sendiri,
   nol kemunculan di rule SQL mana pun. Yang dibuat hanya kunci utama dan kunci tamu, yang memang
   sudah ditetapkan. Menuliskan kolom isinya berarti mengarang.
2. ⚠️ **Selisih nama kolom.** AC nomor 10 menulis `ACCEPTEDNO`; `STRUKTUR-TABEL-CLAIM-LIFE.md`
   menulis `ACCEPTED_NO`. Migrasi memakai **`ACCEPTED_NO`** mengikuti dokumen struktur, sebab
   seluruh kolom baru proyek ini snake_case. Selisihnya dicatat, bukan didiamkan.
3. ⚠️ Kelima puluh lima kolom tabel warisan `[terverifikasi]` dari rule
   `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` bertipe
   `Rule-Connect-SQL` — dibaca langsung dari korpus, bukan dari ingatan.
4. ⚠️ **Daftar "14 atribut polis berulang" tidak pernah diurutkan** di dokumen mana pun. Yang
   dibandingkan pembongkar adalah atribut yang bentuk barunya memang menyimpan di tingkat klaim
   (`NO_CLAIM`, `POLICY_NO`, `BUSINESSNAME`, `CLAIM_RETRO`) — bukan tebakan atas daftar yang tidak
   ada. Bila daftar sebenarnya ditetapkan, satu variabel di `migrasidata.go` yang berubah.
5. ⚠️ **Tabel pencatat migrasi `T_MIGRASI` `[usulan]`** — nama dan bentuknya belum pernah
   diputuskan lewat ADR.
6. ⚠️ Test bentuk migrasi tinggal di paket `repository`, bukan di seam `services` murni yang
   disebut brief induk §5. Ia tidak menyentuh basis data sama sekali — yang diuji adalah isi
   berkas `.sql` yang ditanam ke biner. Perbedaan penempatan dicatat, bukan didiamkan.

---

## Implementasi — ronde 2, 26 September 2026

**Status: `claimed`** — **38 dari 53 AC tertutup** (dari 25), diverifikasi **73 test Go yang
benar-benar berjalan** (dari 34). Sesi ini kembali berjalan di **jalur B**: `ORACLE_DSN` kosong,
sehingga Langkah A dilewati dan seluruh test bertag `db` **melewati dengan pesan**.

⛔ **Tiket ini tetap tidak dapat `resolved`.** Sebabnya tunggal dan tidak berubah: butir 2d brief
ronde 2 — apakah `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)` — **ditegaskan
work owner 26 September 2026 untuk tetap `[USULAN]`**, sehingga constraint-nya tidak dipasang dan
AC 50 tetap terbuka. Migrasi juga **belum pernah dijalankan di Oracle mana pun**.

### Sebelas temuan verifikasi — sepuluh tertutup penuh, satu sebagian

| # | Temuan | Keadaan |
| ---: | --- | --- |
| 1 | Lima kolom STRUKTUR hilang dari `T_CLAIMLF_PREMIUMLIST_DETAIL` | ✅ ditambahkan; **test pembanding DDL lawan STRUKTUR** dibuat, tabel demi tabel |
| 2 | `ACCEPT_STATUS` ada padahal STRUKTUR membuangnya | ✅ dibuang dari DDL, model, dan `INSERT` |
| 3 | Pembaca tabel warisan dari Oracle tidak ada | ✅ `AmbilBarisLama` ditulis; test db kini **membaca kembali dari tabel** |
| 4 | Pembacaan cicit hanya sebagian, galat parse ditelan | ✅ seluruh kolom uang/rasio dibaca; galat **dikembalikan** |
| 5 | Idempoten hanya untuk jalur mulus | ✅ toleransi ORA-00955/ORA-02264 per pernyataan + laporan |
| 6 | AC 11 kolom bank belum dikerjakan | ✅ model, pembongkar, `Simpan`, dan jalur baliknya |
| 7 | 14 AC tertulis tanpa test | ⚠️ **12 tertutup**. Nomor 33 **punya test** (`TestAC33IdentitasBertipeSama`) tetapi tetap tidak dicentang: satu dari empat kolomnya, `T_GENERAL_KOMITE.ID`, ada di tabel yang tidak dibuat tiket ini. Nomor 31 seluruhnya tentang `T_GENERAL_KOMITE` — nol test, nol centang |
| 8 | Kotak AC tidak dicentang; cacah test berselisih | ✅ 38 dicentang; "32 test" diralat menjadi 34 |
| 9 | Baris datar warisan ditulis sebagian | ✅ dicatat `[terbuka — tiket 02/03]`: **18 dari 55** kolom ditulis, 37 sisanya disebut namanya satu per satu |
| 10 | Penjaga akhiran baris belum ada | ✅ test nol byte CR + `.gitattributes` dua pola |
| 11 | `Hapus` adalah DELETE fisik | ✅ komentar batas + test statik nol pemanggil di luar `_test.go` |

### Keputusan work owner 26 September 2026

| | Keputusan | Yang dikerjakan |
| --- | --- | --- |
| **c** | `[DIPUTUSKAN]` kolom share/persen/rate `NUMBER(38,8)`, di Go `Ratio` | lima kolom `003` diubah dari `NUMBER` arbitrer |
| **j** | `[DIPUTUSKAN]` perpendek sekarang | tabel `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` (36 byte) menjadi **`T_CLAIMLF_ADJ_SPREADING_RETRO`** (29); kedua kolom valuasi retro memakai nama korpus `RETRO_VALUATION_*` (26 dan 28 byte) |
| **d** | tetap `[USULAN]` | constraint **tidak** dipasang; AC 50 tetap terbuka; tiket tetap `claimed` |
| **e** | **ditahan** | `002_t_general_claim.sql` **tidak disentuh** — lihat temuan korpus di bawah |
| **i** | — | diselesaikan lewat temuan 2: STRUKTUR sendiri memuat `[keputusan work owner]` 2026-09-18 yang membuang kolomnya |

⭐ **Temuan korpus baru yang membalik usulan 2e.** Ketiga `TANGGAL_*` **bukan tanggal proses klaim**.
`[terverifikasi]` rantainya utuh: `TanggalRespon`←`CARI12`←**`RESPONSE_DATE`**,
`TanggalKonfirmasiBalik`←`CARI15`←**`RECONFIRMATION_DATE`**,
`TanggalRealisasi`←`CARI16`←**`REALIZATION_DATE`**, seluruhnya dari
`SELECT … FROM POOLDATA.JSON_OFFER_LIFE WHERE STATUS = 'Bind' AND OLDID IS NULL`
(`PremiumList Life/RDBList/GetOfferLife_sql.xml`), dipasang oleh `setNoOffer_Act.xml` berkelas
`ASM-FW-GISFW-Work-LIFE` — kelas **Offer**, bukan kelas Claim. Itu **tanggal penawaran**, dan
ketiganya selesai sebelum klaim pertama ada.

Tiga bukti Claim Life hanya menumpang menampilkan: ketiganya muncul **hanya di 4 berkas `Section\`**
dan **nol `Activity\`**; ketiga kolomnya **nol kemunculan** di `UpdateOsAkseptasiClaimLife_sql.xml`;
dan nama `TANGGAL_RESPON` **nol kemunculan** di 9.430 berkas korpus — itu nama karangan spec.

⚠️ `[terbuka]` Bentuk yang konsisten adalah **baca hidup** lewat tautan ke offer, seperti
`TEAM_GROUP` lewat `MO_ID` — bukan kolom baru di header klaim. Yang menghalangi: sumbernya bernama
`JSON_OFFER_LIFE`, dan bentuk relasional penggantinya milik **modul PremiumList Life**.

### Berkas

| Berkas baru | Total | Berisi |
| --- | ---: | ---: |
| `internal/repository/barislamakolom.go` | 226 | 210 |
| `internal/repository/strukturkolom_test.go` | 231 | 213 |
| `internal/repository/acbentuk_test.go` | 252 | 235 |
| `internal/repository/barislamakolom_test.go` | 184 | 173 |
| `internal/repository/batasanpemakaian_test.go` | 150 | 141 |
| `internal/repository/uraidesimal_test.go` | 106 | 99 |
| `.gitattributes` (akar repositori) | 5 | 5 |
| **jumlah berkas baru** | **1.279** | **1.201** |

Diubah: 16 berkas. Seluruh diff sesi ini **+2.247 / −136** atas 24 berkas — `migrasi.go`, `migrasidata.go`, `pohonklaim.go`,
`skemauji.go`, `main.go`, `models/klaimlife.go`, `models/pohonklaim.go`, tiga berkas `.sql`,
dan tiga berkas test.

### Yang dijalankan, dan hasilnya

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` | lulus |
| `gofmt -l` | nol berkas |
| `go build ./...` | lulus |
| `go test ./...` | lulus, **73 test** |
| `go test -tags=db ./internal/...` | **18 MELEWATI** dengan `ORACLE_DSN belum dikonfigurasi`, 1,4 detik |
| `npm run typecheck` · `npm test` · `npm run build` | lulus, 5 test — frontend tidak diubah ronde ini |

⭐ **Setiap instrumen baru diuji lebih dulu pada kasus yang sudah diketahui salah**, dan dipulihkan
sesudahnya: kolom STRUKTUR dicabut dari DDL → gagal; kolom karangan ditambahkan → gagal; pengenal
40 byte → gagal; `.Hapus(` disisipkan ke `services` → gagal; CR ditanam ke berkas migrasi → gagal;
satu kolom dicabut dari daftar 55 → gagal; kolom uang dijadikan `VARCHAR2` → gagal; FK dokumen
dialihkan ke header → gagal. Tanpa langkah itu, "66 test lulus" tidak membuktikan apa pun.

⚠️ Satu kasus uji saya **salah hitung**: `PK_T_CLAIMLF_PLD_YANG_PANJANGX` saya kira 31 byte,
ternyata tepat 30 — jadi test yang lulus memang benar lulus. Diulang dengan nama 40 byte.

### AC yang masih terbuka — 15

- **Menunggu Oracle** (20, 21, 38, 51, 53): kodenya lengkap, test-nya ditulis, dan ia melewati.
- **Menunggu work owner** (8, 15, 27, 35, 50): lima butir `[terbuka]` yang executor tidak menutupnya.
- **Di luar tiket ini** (31, 33): keduanya menyebut `T_GENERAL_KOMITE`, milik konteks Komite.
- **Belum dikerjakan** (19, 46): `PRODUCT_NAME` dari `product_life` relasional; dan AC 46 yang
  sumbernya **tidak ada di tabel datar warisan** — perlu ekspor dari pemilik Pega (brief 2k).

### Catatan

1. ⚠️ `[data DBA]` **Tipe fisik kolom `OS_AKSEPTASI_KLAIM_LIFE` sungguhan belum dikonfirmasi.**
   Yang `[terverifikasi]` hanya NAMA kelima puluh lima kolomnya. Penggolongan 9 kolom angka dan
   8 kolom tanggal di `barislamakolom.go` adalah pembacaan atas **nama**, bukan atas katalog
   Oracle. Bila sebuah kolom ternyata `VARCHAR2` di produksi, `TO_CHAR` berformat angka menjawab
   ORA-01722 — dan yang berubah dua daftar itu, bukan kode pembacanya.
2. ⚠️ Tabel tiruan warisan kini memakai `DATE` dan `NUMBER` untuk ketujuh belas kolom itu. Kalau
   tiruannya berisi teks semua, test pulang-pergi tidak menguji apa pun tentang `TO_CHAR` dan
   jebakan NLS justru lolos.
3. ⚠️ Nama berkas `006_t_claimlf_adjustment_spreading_retro.sql` **tidak** ikut diperpendek; yang
   dilihat Oracle adalah nama tabel di dalamnya. Berkas `001`–`008` masih boleh disunting langsung
   selama `T_MIGRASI` belum pernah ada di instance mana pun (brief 2l).
   > **Ralat 26-09-2026 (ronde 4):** sudah tidak berlaku. `[keputusan work owner butir n]` — kedua
   > berkas `006` diganti nama menjadi `006_t_claimlf_adj_spreading_retro.sql` (+`_down`). Aman
   > sebab syarat di kalimat di atas masih terpenuhi: `T_MIGRASI` belum pernah ada di mana pun.
4. ⚠️ AC mengeja tiga nama berbeda dari STRUKTUR — `ACCEPTEDNO`/`ACCEPTED_NO`,
   `RISLIPRNM`/`RI_SLIP_RNM`, `CASEID`/`CASEID_POLICY`. Ejaan STRUKTUR yang dipakai, dan
   selisihnya ditulis di komentar test AC yang bersangkutan.

### Hasil `/code-review` atas titik tetap `76dcda4` — dan yang diperbaiki sesudahnya

Dua sumbu ditinjau terpisah. Keduanya menghitung ulang angka bab ini sendiri dan **cocok**:
66 test (sebelum perbaikan di bawah), 38 `[x]` + 15 `[ ]` = 53, `Simpan` menulis tepat 18 dari 55
kolom warisan, dan butir **d** serta **e** memang tidak disentuh.

**Cacat yang diperbaiki sesudah tinjauan:**

| Sumbu | Temuan | Perbaikan |
| --- | --- | --- |
| standards | ⛔ Komentar kepala `CacahBarisLama` terlepas dan menempel di atas `AmbilBarisLama`; godoc keduanya salah | dikembalikan ke fungsinya masing-masing |
| spec | ⛔ Test idempoten membuat `T_WORK_CLAIM (ID)` **satu kolom**, lalu menuntut migrasi LULUS — skema cacat resmi menjadi "migrasi sukses" | test kini menjalankan **pernyataan pertama langkah 001 yang sesungguhnya** lewat `PernyataanLangkah`, seperti brief §3-5 memang memerintahkan |
| spec | ⛔ `TestAC45Dan48` memakai assertion tautologis `!strings.Contains(sql, "CREATE")` — berkas 004 memuat tiga `CREATE`, jadi ia tidak pernah bisa gagal | diganti pemeriksaan `CREATE INDEX` atas `T_CLAIMLF_ADJUSTMENT (KOMITE_ID)`; diuji gagal saat index-nya dicabut |
| spec | Temuan 4 baru separuh: tujuh kolom uang/rasio dibaca, **nol** diuji | `contohPohon` mengisi ketujuhnya; `TestBacaSampaiCicit` memeriksanya **digit demi digit**, plus satu kolom sengaja kosong untuk membuktikan kosong ≠ nol |
| spec | `uraiUang`/`uraiRasio` nol test padahal itu inti temuan 4 | `uraidesimal_test.go` — 5 test murni: galat dikembalikan, NULL/kosong bukan galat, digit tidak berubah, galat menyebut baris dan kolomnya |
| spec | Temuan 1 membandingkan nama kolom saja; selisih tipe tetap lolos | `TestGolonganTipeDDLCocokDenganStruktur` — **113 kolom** dibandingkan golongannya (teks · tanggal · angka) |
| standards | Duplikasi `uraiUang`/`uraiRasio`; parameter `sumber` sebenarnya ID baris; literal nama tabel retro dua kali | diangkat ke `uraiDesimal`; diganti `idBaris`; menjadi konstanta `namaTabelRetro` |
| standards | Komentar 37-kolom dapat basi diam-diam | dikunci `TestCacahKolomWarisanYangDitulisSimpan` |

**Temuan yang ditolak, dengan alasan:**

- *"Rename tabel 006 tidak punya jalur migrasi"* — brief §2 l: berkas `001`–`008` masih boleh
  disunting langsung selama `T_MIGRASI` belum pernah ada di instance mana pun. Tidak ada skema yang
  memuat nama lama.
- *"`lap.Pernyataan++` dilewati `continue` sehingga cacahnya menyimpang"* — medan itu memang
  mencacah pernyataan yang **dieksekusi**; yang dilewati muncul di `ObjekSudahAda`. Ditegaskan di
  doc-nya, bukan diubah.

**⚠️ `[terbuka]` Kelemahan yang TIDAK tertutup ronde ini, dinyatakan apa adanya:**

1. ⛔ **Uji yang mengonfirmasi dirinya sendiri.** `AmbilBarisLama` membungkus 9 kolom angka dan
   8 kolom tanggal dengan `TO_CHAR` berformat, dan tabel **tiruan** dibangun dari penggolongan yang
   sama. Kedua sisi test berasal dari satu tebakan, sehingga `ORA-01722` / `ORA-01861` pada tabel
   **sungguhan** mustahil tertangkap di sini. Yang dapat menutupnya hanya **daftar tipe kolom
   `OS_AKSEPTASI_KLAIM_LIFE` dari DBA** — sudah masuk permintaan.
2. ⚠️ **Toleransi `ORA-00955` tidak memeriksa bentuk objek.** Oracle hanya menjawab "nama sudah
   dipakai". Objek lama berbentuk berbeda karena itu ikut diterima. Yang menjaga hal ini bukan kode
   melainkan jalur mundur — bongkar dulu, baru pasang lagi.
3. ⚠️ **`UX_ADJ_KOMITE_ID` bersifat UNIQUE** (sejak ronde 1), yang berarti **satu kasus komite =
   satu baris adjustment**. AC 48 hanya meminta ber-index, tidak menyebut unik. Konsekuensinya belum
   pernah dibahas dan tidak diubah ronde ini.
4. ⚠️ Tiga test melampaui kesebelas temuan: `TestHandlersTidakMengimporRepository`,
   `TestTabelDikecualikanTetapDibuat`, `TestRingkasPernyataanPendek`. Ketiganya penjaga, bukan
   perubahan perilaku, dan menegakkan aturan yang memang sudah tertulis di `README-BACA-DULU.md`.

---

## Implementasi — ronde 3, 26 September 2026

**Status: `claimed`** — **38 dari 53 AC tertutup** (tidak berubah dari ronde 2), diverifikasi
**77 test Go yang benar-benar berjalan** (dari 73). Sesi **keempat** berturut-turut di jalur B:
`ORACLE_DSN` kosong, Langkah A dilewati, seluruh test bertag `db` **melewati dengan pesan**.

⛔ **Tiket tetap tidak dapat `resolved`**, dan sebabnya tidak berubah: butir 2d ditegaskan work
owner 26 September 2026 untuk **tetap `[USULAN]` untuk kali keempat**, sehingga `REFERENCES` tidak
dipasang dan AC 50 tetap terbuka. Migrasi juga masih belum pernah dijalankan di Oracle mana pun.

### Sembilan temuan verifikasi ronde 2

| # | Temuan | Keadaan |
| ---: | --- | --- |
| 1 | ⛔ Keputusan c baru diterapkan pada `003`; delapan kolom share/rate masih `NUMBER` polos | ✅ kedelapannya menjadi `NUMBER(38,8)`; `TestAC41…` diperketat dari "memuat kata NUMBER" menjadi `NUMBER(38,8)`; **`TestNolNumberTanpaPresisi`** baru menjaga seluruh 30 kolom `NUMBER` |
| 2 | ⛔ Toleransi `ORA-02264` memalsukan sukses — dan itu usulan brief ronde 2 yang saya jalankan | ✅ `sudahAda` kini **hanya** `ORA-00955`; sesudah sebuah `CREATE` dilewati, keberadaan objeknya **dibuktikan** lewat `ALL_OBJECTS`; test db tabrakan constraint menuntut migrasi **GAGAL** |
| 3 | ⚠️ AC 12 tidak dicentang dan tidak ada di daftar terbuka | ✅ alasannya ditulis di bawah; daftar terbuka kini **berjumlah 15**, sama dengan kotak `[ ]` |
| 4 | ⚠️ "Kelemahan terbuka 3" ronde 2 keliru | ✅ ditarik — lihat ralat di bawah |
| 5 | ⚠️ AC 40 dicentang padahal teksnya menyebut nama tabel lama | ✅ satu baris jejak ditulis di bawah |
| 6 | ⚠️ Angka bab ronde 2 basi | ✅ diralat di bawah |
| 7 | ⚠️ Uji tabel warisan mengonfirmasi dirinya sendiri | ⛔ **TIDAK dikerjakan** — ia menuntut berkas `[data DBA]` daftar tipe kolom yang belum ada. Tetap `[terbuka]` |
| 8 | ⚠️ Kalimat "nama karangan spec" terlalu jauh | ✅ diralat di bawah |
| 9 | ℹ️ Test membaca STRUKTUR di luar `APP_RNM\` | catatan saja; sengaja `t.Fatalf` bila hilang — **tidak diubah** |

### Ralat atas laporan ronde 2 sendiri

1. **Angka diff basi.** Ronde 2 menulis "24 berkas +2.247/−136". Yang benar, diukur sesudah
   perbaikan `/code-review`: **`76dcda4`→`8f5453b` = 24 berkas, +2.293 / −136** (termasuk brief);
   tanpa brief `8b18b00`→`8f5453b` = 23 berkas, +2.123 / −136. `strukturkolom_test.go` bukan
   231/213 melainkan **368/344**. Sebabnya satu: angkanya diambil **sebelum** perbaikan review.
   Ronde ini angkanya diambil **sesudah** commit terakhir.
2. **"`TANGGAL_RESPON` nama karangan spec" terlalu jauh.** Yang nol kemunculan hanya **ejaan**
   `TANGGAL_RESPON`; properti Pega **`TanggalRespon` memang ada** — 2 kemunculan di
   `setNoOffer_Act.xml` dan 4 berkas `Section\` Claim Life. Spec menurunkan nama kolomnya dari
   properti sungguhan. Kesimpulan pokoknya **tetap sahih**: ketiganya tanggal penawaran milik Offer.
3. **"Kelemahan terbuka 3" ditarik.** Ronde 2 menulis `UX_ADJ_KOMITE_ID` UNIQUE "konsekuensinya
   belum pernah dibahas". Keliru: `STRUKTUR-TABEL-CLAIM-LIFE.md` bab `T_CLAIMLF_ADJUSTMENT` memuat
   `[keputusan work owner]` — *"Satu baris `AdjustmentList` = TEPAT satu kasus komite, dan
   sebaliknya. Karena itu `KOMITE_ID` ber-index UNIK meski nullable"*. Sudah diputuskan, dan index
   unik Oracle memang mengizinkan banyak `NULL`. **Nol perubahan DDL.**
4. **Jejak ganti nama.** AC 2, 38, dan 40 menulis `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`; **nama
   fisiknya sejak 26-09-2026 adalah `T_CLAIMLF_ADJ_SPREADING_RETRO`** (keputusan work owner j, batas
   30 byte pengenal Oracle). Teks AC tidak diubah. ⚠️ Butir m (blok ralat di STRUKTUR) dan n (ganti
   nama berkas migrasi `006`) **tetap `[USULAN]`** dan **tidak dikerjakan**, sehingga dokumen dan
   nama berkas masih menulis nama 36 byte.

### Yang dikerjakan pada kode

**Temuan 1 — presisi.** Delapan kolom: `004` `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`,
`SHARE_RETRO`, `RETROCEDED_SHARE`; `005` `RETROCADED_SHARE`, `RATE`; `006` `PERCENT_SHARE`, `RATE`.
Sesudahnya seluruh DDL memuat **28 `NUMBER(38,8)`**, satu `NUMBER(5)` (`AGE`), satu `NUMBER(19)`
(`DOCUMENT_CLAIM.ID`) — **nol `NUMBER` tanpa presisi**. Komentar kepala ketiga berkas menyebut
keputusan c.

**Temuan 2 — penjaga `ORA-02264`.** `ORA-02264` berarti nama constraint sudah dipakai, dan Oracle
baru memeriksanya **ketika tabelnya belum ada**; bila tabelnya ada ia menjawab `ORA-00955` lebih
dulu. Jadi `ORA-02264` pada sebuah `CREATE TABLE` berarti tabel itu **justru tidak terbuat**.
Sekarang: `sudahAda` hanya `ORA-00955`; setiap `CREATE` yang dilewati **dibuktikan keberadaannya**
lewat `SELECT COUNT(*) FROM ALL_OBJECTS WHERE OWNER = :1 AND OBJECT_NAME = :2`; nol → galat
*"dilaporkan sudah ada, tidak ditemukan di katalog"* dan langkahnya **tidak dicatat** di `T_MIGRASI`.

⚠️ Ini menutup **separuh** kelemahan terbuka 2 ronde 2: **keberadaan** kini diperiksa,
**bentuk** masih tidak. Kolom yang berbeda tetap lolos; yang menjaganya jalur mundur.

### AC 12 — kenapa tetap terbuka

AC 12 menuntut *"seluruh uang dan share bertipe **desimal presisi arbitrer**"*. Tiga klausa lainnya
terpenuhi dan teruji: tanggal `DATE`, kolom nullable, dan identitas dari sequence — yang terakhir
tercakup penyimpangan sadar **AC 34** (nomor bisnis berformat untuk `T_WORK_CLAIM` dan kedua tabel
ber-shared-PK; ADR-U-0006 tetap berlaku untuk `T_CLAIMLF_*` dan `DOCUMENT_CLAIM`).

⛔ Klausa pertamanya **dibalik oleh keputusan work owner c**: `NUMBER(38,8)` adalah presisi
**tetap**, bukan arbitrer. Executor tidak mencentang AC yang teksnya sendiri sudah tidak berlaku;
yang diperlukan satu baris dari work owner yang menyatakan c menggantikan klausa itu.

### AC yang masih terbuka — 15, sama dengan jumlah kotak `[ ]`

| Sebab | Nomor |
| --- | --- |
| Menunggu instance Oracle | 20, 21, 38, 51, 53 |
| Menunggu work owner | 8, 15, 27, 35, 50 |
| **Teks AC dibalik keputusan c** | **12** |
| Di luar tiket ini — konteks Komite | 31, 33 |
| Pemilik export Pega · modul PremiumList Life | 46 · 19 |

### Keputusan §2 sesi ini

| | Keadaan | Yang dikerjakan |
| --- | --- | --- |
| **c** | `[DIPUTUSKAN]` | **dituntaskan** — delapan kolom sisa |
| **o** | ⭐ `[DIPUTUSKAN]` 26-09: *"jangan ada lagi pemanggilan procedure, segala procedure hardcode dalam skrip"* | dicatat di tiket 02; **tidak** dikerjakan di sini |
| **d** | `[USULAN]` kali keempat | tidak disentuh; tiket tetap `claimed` |
| **e′** | `[USULAN]` | `002` tidak disentuh; nol kolom `TANGGAL_*` di DDL mana pun |
| **m, n** | `[USULAN]` | STRUKTUR tidak disunting; berkas `006` tidak diganti nama |

### Yang dijalankan, dan hasilnya

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l` · `go build ./...` | lulus, nol berkas |
| `go test ./...` | lulus, **77 test** |
| `go test -tags=db ./internal/...` | **19 MELEWATI** dengan `ORACLE_DSN belum dikonfigurasi` |
| `npm run typecheck` · `npm test` · `npm run build` | lulus, 5 test — frontend tidak diubah ronde ini |

⭐ **Tiap penjaga baru diuji gagal dulu pada kasus buruknya**, lalu dipulihkan: satu kolom
dikembalikan ke `NUMBER` polos → gagal; `ORA-02264` dikembalikan ke `sudahAda` → gagal. Yang kedua
penting: test ronde 2 justru **mengunci** perilaku yang salah, dan sekarang ia menolaknya.

### Hasil `/code-review` atas titik tetap `62ca5cc`

Dua sumbu ditinjau terpisah. Angka yang keduanya hitung ulang **cocok**: nol `NUMBER` tanpa presisi
(28 + 1 + 1 = 30), 38 `[x]` + 15 `[ ]` = 53, daftar terbuka berjumlah 15 = jumlah `[ ]`, dan
butir **d, e′, m, n** memang tidak disentuh.

**Cacat nyata yang diperbaiki sesudah tinjauan:**

| Temuan | Perbaikan |
| --- | --- |
| ⛔ **Penjaga palsu.** `TestSeluruhCreateDapatDibacaNamanya` memecah ulang teks yang sudah disambung `seluruhSQL`, sehingga pemisahnya **tidak pernah memisah apa pun** — yang diperiksa hanya **8 dari 19** pernyataan `CREATE`, dan seluruh `CREATE INDEX` lolos | pernyataan diambil dari `daftarMigrasi` apa adanya; cacahnya **dikunci 19** supaya kelemahan yang sama gagal, bukan diam |
| ⛔ **`ALL_OBJECTS` telanjang** melanggar ADR-U-0033 Akibat 1, dan Akibat 3 menuntut test yang menemukannya gagal — test itu tidak pernah ada | menjadi **`SYS.ALL_OBJECTS`**; **`TestNolNamaTabelTelanjangDiQuery`** baru memeriksa 33 rujukan tabel di paket repository |
| ⛔ **Komentar `005` bertentangan dengan DDL-nya sendiri** — masih berbunyi *"tipe pastinya masih USULAN … NUMBER berpresisi arbitrer"* enam baris di atas `NUMBER(38,8)` | paragraf basi dicabut |
| ⚠️ Slot `harusYa` terbuang menjadi duplikat sesudah `ORA-02264` pindah | diisi bentuk pembungkus ketiga yang benar-benar berbeda |
| ⚠️ `objekAda` membandingkan `OWNER` bergantung huruf besar-kecil; skema huruf kecil akan **menggagalkan migrasi yang sehat** | `UPPER()` di kedua sisi |
| ⚠️ Galat pemeriksaan tidak membungkus nama langkah, berbeda dari cabang tetangganya | dibungkus |

**⛔ Penyimpangan sadar yang ditemukan tinjauan dan kini dicatat:**

ADR-U-0016 Akibat 2 berbunyi *"Kolom **persen** tidak termasuk — ia bukan uang, dan tetap mengikuti
ketetapan modulnya"*. Keputusan work owner **c** menyamakan presisi kolom persen dan rate dengan
kolom uang, sehingga ia **menyimpang dari ADR-U-0016**. ⚠️ `[terbuka]` ADR-U-0016 **belum
diamandemen** dan tidak ada ADR baru bertanggal 2026-09-26; sampai work owner menuliskannya,
catatan itu dipikul komentar kepala berkas `004`, `005`, dan `006`.

> ⭐ **RALAT 26 September 2026 (ronde 4) — paragraf di atas KELIRU, dan `[terbuka]`-nya DICABUT.**
> `[keputusan work owner butir p1]`: Akibat 2 tidak melarang apa pun di sini — kalimat *"tetap
> mengikuti ketetapan modulnya"* justru **menyerahkan** kolom persen kepada modul masing-masing.
> Keputusan **c** adalah ketetapan modul Claim Life, persis yang diserahkan kepadanya. Jadi c
> **sesuai** ADR-U-0016, bukan menyimpang darinya, dan **nol amandemen ADR diperlukan**. Label di
> kepala `004`, `005`, `006` sudah diganti menjadi *"ketetapan modul Claim Life"*.

**Tuduhan yang saya tolak sesudah mengukur ulang:**

- *"Blob index memuat CRLF — 004 = 55, 005 = 43, 006 = 47"*. **Salah.** Diukur byte demi byte
  lewat `git show`: HEAD, index, dan working tree **semuanya `CR = 0`**. Angka 55/43/47 itu cacah
  **baris**, bukan cacah CR. ⚠️ Instrumen saya sendiri sempat menghasilkan angka yang sama
  kelirunya — `grep -c` dengan pola CR di dalam loop shell tidak mengembang seperti dugaan — dan
  itulah sebabnya angkanya diulang dengan alat kedua sebelum dipercaya.
- *"`PANDUAN-MENJALANKAN.txt` 376 baris adalah scope creep ronde ini"*. **Bukan milik sesi ini.**
  Berkas itu dibuat sesi lain dan hanya ikut ter-*stage* oleh `git add -A`; ia sudah dikeluarkan
  dari commit dan dibiarkan utuh di working tree.

**Yang dicatat tetapi TIDAK diubah:** jalur `nama == ""` menggagalkan `CREATE` yang bentuknya tidak
dikenali regex (`CREATE OR REPLACE`, `GLOBAL TEMPORARY`, pengenal berkutip) — tidak ada bentuk
begitu hari ini, dan `TestSeluruhCreateDapatDibacaNamanya` mengunci kesembilan belasnya;
`migrasi.go` kini juga mengurai DDL dan menanyai katalog (**Divergent Change**), diterima demi
menjaga pemeriksaan tetap satu tempat dengan pemakainya; dan kode galat Oracle masih berupa teks
yang tersebar (**Primitive Obsession**).

---

## Implementasi — ronde 4, 26 September 2026

**Gerbang yang terbuka: G2 saja.** G1 tertutup (`ORACLE_DSN` kosong — sesi **kelima** berturut-turut
di jalur B); G3 tertutup sebab empat objek DBA §5 brief belum diserahkan, sehingga o1–o3 tidak
ditanyakan: keputusannya sendiri tidak cukup membukanya.

**Status: `claimed`** — **40 dari 53 AC tertutup** (dari 38), diverifikasi **83 test Go yang
benar-benar berjalan** (dari 77). Angka ini diambil **sesudah** perbaikan `/code-review`. Sesuai brief §4, ini ditulis di awal, bukan di akhir.

⛔ **Kenapa belum `resolved` meski butir d akhirnya disahkan.** Syarat `resolved` brief ada empat,
dan yang terpenuhi baru satu: **d `[DIPUTUSKAN]` dan diterapkan ✅**; migrasi sungguhan 0/8/kosong
**belum**; seluruh test db PASS **belum**; AC `[terbuka]` tinggal milik luar-executor **belum** —
lima di antaranya (20, 21, 38, 51, 53) menunggu Oracle, dan itu ada di dalam jangkauan executor
begitu G1 terbuka.

### Lima temuan verifikasi ronde 3

| # | Temuan | Keadaan |
| ---: | --- | --- |
| 1 | ⚠️ Label "PENYIMPANGAN SADAR dari ADR-U-0016 Akibat 2" keliru | ✅ butir **p1**: label diganti di **tiga** `.sql`; `[terbuka]` amandemen ADR **dicabut** lewat blok ralat bertanggal di bab ronde 3 di atas. ⚠️ Brief menyebut label itu ada juga di `presisiSah`; **tidak** — isinya sudah "keputusan work owner c". Komentarnya tetap diperjelas. ⛔ Tempat keempat ini **terlewat** di lintasan pertama dan ditemukan `/code-review` |
| 2 | ⚠️ Konsekuensi o di tiket 02 kurang satu AC | ✅ ralat ditulis di tiket 02: **AC 2, 3, dan 10** |
| 3 | ⚠️ AC 12 dibaca terlalu sempit | ✅ butir **t**: dicentang; daftar terbuka turun |
| 4 | ℹ️ `TestNolNamaTabelTelanjangDiQuery` memindai komentar juga | catatan saja — **tidak diubah**, nol gagal palsu sampai hari ini |
| 5 | ⛔ Skema uji tidak punya pagar | ✅ ditutup, dan **lebih luas dari yang diminta** — lihat di bawah |

**4 dari 4 temuan yang dapat dikerjakan ditutup**; yang kelima memang hanya catatan.

### Temuan 5 — pagar skema uji, dan satu hal yang brief sendiri tidak lihat

`skemauji.Buka` dan `BukaRepositori` kini **menolak dengan galat** kecuali **dua** syarat terpenuhi:
env `ORACLE_SKEMA_UJI=true` **dan** `ORACLE_SCHEMA` bukan `POOLDATA`. Pesannya menyebut sebabnya,
termasuk bahwa test ini **menghapus** tabel. Bawaan `POOLDATA` dicabut dari `Makefile`.

⭐ **Yang tidak tertulis di brief, ditemukan saat mengerjakannya:** memasang pagar di `skemauji`
saja **tidak cukup** — keenam pemanggilnya menjawab **setiap** galat `Buka()` dengan
`t.Skipf("lewati: %v", err)`. Pagar itu akan melewati diam-diam, persis yang brief larang. Karena
itu `BolehDilewati(err)` ditambahkan, dan keenam titik kini `t.Fatalf` untuk apa pun selain
`ErrTanpaOracle`.

⛔ **Akibat sampingannya menyingkap lubang yang sudah ada sejak dulu:** dengan pola lama,
`ErrProduksi` — penolakan `IS_PEGA_PROD=true` — **juga** dijawab SKIP. Menunjuk instance produksi
karena itu menghasilkan lari **hijau yang tidak menguji apa pun**. Sekarang ia menggagalkan test.

⚠️ **Tempat keempat yang tidak disebut brief**: `.env.example` juga memberi `ORACLE_SCHEMA=POOLDATA`.
Ikut dicabut, dan `ORACLE_SKEMA_UJI` didaftarkan di sana.

`-migrate` **tidak** ikut dipagari, sesuai brief: ia hanya `CREATE`, tidak pernah `DROP`.

### Keputusan §2 sesi ini

| | Keadaan | Yang dikerjakan |
| --- | --- | --- |
| **d** | ⭐ `[DIPUTUSKAN]` 26-09 sesudah **empat sesi** — *"pasang `REFERENCES`, nullable"* | `FK_WORK_COVER_KEY` di `001`, `FK_ADJ_KOMITE` di `004`, keduanya **tanpa `ON DELETE`**; fixture diberi baris komite yang sungguh ada; AC 50 dicentang; dikunci `TestAC50PenunjukKeAtasBerReferences` |
| **p** | `[DIPUTUSKAN]` — **p1** | label di tiga `.sql` diralat; nol amandemen ADR |
| **t** | `[DIPUTUSKAN]` — centang | AC 12 dicentang, dengan alasannya ditulis di bawah AC-nya |
| **e′, m, n** | ⭐ didelegasikan ke executor 26-09: *"rekomendasikan saja jawaban paling baik, dan samakan dengan XML logic sistem existing (Pega)"* | **ketiganya dikerjakan**, masing-masing dengan bukti korpusnya — lihat di bawah |
| **o1, o2, o3** | tidak ditanyakan | G3 tertutup: empat objek DBA belum ada, jadi keputusannya saja tidak membukanya |

### e′, m, n — rekomendasi executor, disandarkan ke korpus Pega

**e′ — ketiga `TANGGAL_*` bukan milik Claim Life.** Diperiksa **ketiga kalinya, dengan alat
berbeda**: ketiga **nama kolom** (`RESPONSE_DATE`, `RECONFIRMATION_DATE`, `REALIZATION_DATE`) muncul
**nol kali** di seluruh modul Claim; ketiga **propertinya** muncul di Claim Life hanya di **4 berkas
`Section\`** dan **nol berkas `Activity\`**. Logic Pega-nya jelas: Claim Life **menampilkan**
ketiganya, **tidak pernah menyimpannya**. AC 8 diberi catatan `[terbuka — PremiumList Life]` untuk
tiga kolom itu; `002` tidak disentuh. ⚠️ Dua kolom sisanya (`TEAM_GROUP`, `BUSINESS_ID`) tetap
`[terbuka]` penuh, sehingga **AC 8 tetap tidak dicentang**.

**m dan n — nama fisik.** ⭐ Bukti korpus yang menentukan: **seluruh nama tabel fisik yang sungguh
dipakai Pega di modul Claim Life berukuran ≤ 23 byte**, terpanjang `OS_AKSEPTASI_KLAIM_LIFE`. Nama
36 byte `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` **tidak pernah ada di korpus** — ia lahir di dokumen
spec, bukan di Pega. (`SpreadingRetro` sendiri muncul di korpus hanya di modul **Treaty In**, nol di
Claim Life.) Memendekkannya karena itu **mengikuti** sistem berjalan, bukan menyimpang darinya.
**m** → satu blok ralat bertanggal di STRUKTUR bab itu, tanpa menulis ulang judul/diagram/tabel
relasinya. **n** → `git mv` dua berkas `006` menjadi `006_t_claimlf_adj_spreading_retro.sql`; aman
sebab `T_MIGRASI` belum pernah ada di mana pun.

### ⚠️ `[terbuka]` baru yang lahir dari keputusan d — pemiliknya tiket 15

Model TINGKAT 1 menyatakan baris KOMITE memuat `COVER_KEY = ID baris klaim`. Dengan FK tanpa
`ON DELETE`, menghapus baris klaim yang masih ditunjuk baris komite akan **ditolak `ORA-02292`**.
`PohonKlaim.Hapus` menghapus baris work paling akhir, jadi urutannya sendiri selamat, tetapi baris
komite tidak dibereskannya. ⛔ **Belum teruji** — penalaran atas semantik Oracle, bukan pengamatan.
Fixture sengaja membiarkan `COVER_KEY` baris komite NULL supaya AC 38 tidak digagalkan hal yang
pemiliknya tiket lain; itu **keputusan sadar**, dan dicatat di `skemauji.go`.

### AC yang masih terbuka — 13, sama dengan jumlah kotak `[ ]`

| Sebab | Nomor |
| --- | --- |
| Menunggu instance Oracle | 20, 21, 38, 51, 53 |
| Menunggu work owner | 8, 15, 27, 35 |
| Di luar tiket ini — konteks Komite | 31, 33 |
| Pemilik export Pega · modul PremiumList Life | 46 · 19 |

### Yang dijalankan, dan hasilnya

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l` · `go build ./...` | lulus, nol berkas |
| `go test ./...` | lulus, **83 test** (dari 77), nol FAIL |
| `go test -tags=db ./internal/...` | **19 MELEWATI** — ⛔ SKIP, bukan PASS |
| `npm run typecheck` · `npm test` · `npm run build` | lulus, 5 test, 87 modul — frontend tidak diubah |

⭐ **Tiap penjaga baru diuji gagal dulu pada kasus buruknya**, lalu dipulihkan: syarat `POOLDATA`
dicabut → gagal; syarat env dicabut → gagal; `FK_ADJ_KOMITE` dicabut → gagal; `ON DELETE CASCADE`
ditambahkan ke FK itu → gagal; `COVER_KEY` dijadikan `NOT NULL` → gagal; satu pemanggil `Buka()`
dikembalikan ke pola lama → gagal. Ditambah satu uji **tuduhan palsu**: `ON DELETE` dipasang pada
constraint **tetangga** → test tetap **lulus**, tidak menuduh yang salah.

### ⛔ Yang ditemukan `/code-review`, dan diperbaiki sebelum commit

| # | Temuan | Perbaikan |
| ---: | --- | --- |
| 1 | ⛔ **Penjaga nullable AC 50 MATI.** Polanya `KOLOM VARCHAR2(32) NOT NULL` berspasi tunggal, sedangkan DDL **meratakan kolom** dengan banyak spasi — tidak akan pernah cocok. Bab ini sempat mengklaim ia menjaga | teks DDL **dirapatkan** dulu; dibuktikan dengan menjadikan `COVER_KEY` `NOT NULL` → kini gagal |
| 2 | ⛔ **Jendela pemeriksaan `ON DELETE` salah.** Ia memotong 120 byte **sesudah baris**, dan karena seluruh pernyataan satu berkas disambung menjadi satu teks, jendelanya menembus `CREATE INDEX` di bawahnya — bisa menuduh constraint yang benar | diganti `klausaConstraint`, yang berhenti di `CONSTRAINT` berikutnya atau akhir daftar kolom |
| 3 | ⛔ **`[terbuka]` amandemen ADR belum dicabut** di bab ronde 3, padahal bab ini sudah mengklaim sebaliknya — **overclaim saya sendiri** | blok ralat bertanggal ditulis di bab ronde 3; klaimnya kini benar |
| 4 | ⚠️ **AC 50 bercentang sambil menyangkal dirinya** — teksnya masih berbunyi `[terbuka]`, *"belum diputuskan"*, *"tidak dinyatakan selesai sebelum jawabannya ada"* | penanda diganti `✅ [DIJAWAB 26-09-2026]`; kalimat lamanya **dicoret, bukan dihapus** |
| 5 | ⚠️ **`make db-up` rusak** oleh pencabutan bawaan: ia masih mengirim `APP_USER="$(ORACLE_SCHEMA)"`, kini kosong | `db-up` menolak bila `ORACLE_SCHEMA` kosong atau `POOLDATA` |
| 6 | ⚠️ **Daftar hitam lebih sempit daripada janji README**: hanya `POOLDATA` persis, sehingga `POOLDATA_DEV` lolos | diperluas menjadi **memuat** `POOLDATA`; dua kasus turunan ditambahkan ke test |
| 7 | ⚠️ **Pola enam baris disalin ke enam tempat** — test db ketujuh yang lupa akan membuka lubangnya lagi, diam-diam | `TestSetiapPemanggilBukaMemeriksaBolehDilewati` mengunci keenamnya beserta cacahnya |
| 8 | ⚠️ Duplikasi tiga pemeriksaan di `Buka` dan `BukaRepositori` | diekstrak menjadi `pastikanAman(cfg)` |
| 9 | ⚠️ `PANDUAN-MENJALANKAN.txt` menulis 81 test | → 83 |

⚠️ **Penjaga baru itu sempat mengukur dirinya sendiri**: ia menemukan **10** pemanggil `Buka()`,
bukan 6, sebab empat kemunculan di komentar dan literalnya sendiri ikut terhitung. Berkasnya
dikecualikan dan polanya dipersempit ke pemanggilan yang hasilnya ditampung.

### Tuduhan tinjauan yang saya TOLAK sesudah memeriksa premisnya

1. ⛔ *"Rename 006 memutus ledger `T_MIGRASI`: skema yang sudah menjalankan nama lama akan
   menjalankannya lagi."* **Premisnya tidak ada.** Migrasi belum pernah berjalan di Oracle mana pun
   — lima sesi berturut-turut, dan katalog instance pengembangan yang dibaca siang ini tidak memuat
   `T_MIGRASI`. Brief §2 n justru mensyaratkan rename dilakukan **hanya sebelum** migrasi pertama,
   dan itulah keadaannya. Tuduhan lanjutannya — *"tanpa satu pun komentar, tiket, atau ADR"* — juga
   keliru: rename tercatat di blok ralat STRUKTUR dan di bab ini.
2. ⚠️ *"`BolehDilewati` dan keenam `t.Fatalf` adalah scope creep; Langkah B menuntut nol perubahan
   perilaku."* **Ditolak sebagian.** Tanpa keduanya pagar §3-5 **tidak terpasang sama sekali** —
   galatnya akan dijawab `t.Skipf` dan lewat diam-diam, persis yang brief larang. Ia bukan tambahan
   di luar temuan 5; ia syarat agar temuan 5 berlaku. Langkah B memang menuntut nol perubahan
   perilaku, tetapi temuan 5 ada di Langkah B dan **isinya sendiri** perubahan perilaku.
3. ⚠️ *"Blok ralat STRUKTUR melampaui butir m; ia menurunkan dasar keputusan j menjadi `[dugaan]`."*
   **Diterima sebagai catatan, tidak dibatalkan.** Blok itu awalnya saya tulis dengan alasan
   "Oracle di bawah 12.2 menolak" — lalu katalog siang ini menunjukkan instance pengembangan
   **12.2.0.1**. Membiarkan kalimat yang sudah diketahui keliru lebih buruk daripada melampaui m
   satu paragraf. ⚠️ **Work owner dipersilakan mencabut paragraf itu bila menganggapnya di luar m.**
   Keputusan **j** sendiri **tidak** dibatalkan executor.

### ⭐ Katalog instance pengembangan dibaca di tengah sesi — brief §9

Sesi verifikasi menyambung ke instance pengembangan (**katalog saja, nol baris data**) dan
menghasilkan dua dokumen `[data DBA — dibaca sendiri, belum dikonfirmasi DBA]`:
`TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md` dan `SUMBER-PENOMORAN-DBA.md`. Keduanya **tidak ditulis
sesi ini** dan **tidak ikut commit ronde 4**; keduanya dibaca dan diperiksa nol kredensial, nol
alamat host, nol baris data.

⭐ **Yang paling penting bagi ronde ini: kredensial itu menunjuk `POOLDATA` — persis nama yang
sampai pagi ini menjadi bawaan `ORACLE_SCHEMA` di `Makefile` dan `.env.example`.** Skema itu memuat
**760 tabel, 383 sequence, 198 procedure**, dan `OS_AKSEPTASI_KLAIM_LIFE` di dalamnya berisi
**13.694 baris**. Menjalankan `go test -tags=db` dengan bawaan itu akan mengirim
`DROP TABLE POOLDATA.OS_AKSEPTASI_KLAIM_LIFE CASCADE CONSTRAINTS`. Pagar temuan 5 yang dibangun
sesi ini menolaknya **dua kali**: tanpa `ORACLE_SKEMA_UJI=true`, dan karena skemanya `POOLDATA`.
Temuan itu ternyata bukan kemungkinan teoretis.

⛔ **G1 tetap tertutup.** `POOLDATA` pengembangan **bukan** skema uji. DBA diminta membuat user
kosong baru di instance yang sama.

**Akibat lain, dicatat tetapi TIDAK dikerjakan ronde ini** — semuanya `[USULAN]` dan pemiliknya work
owner:

| Fakta baru | Akibat |
| --- | --- |
| ⛔ **`DOCUMENT_CLAIM` sudah ada** di `POOLDATA` dengan **14 kolom berbeda**, 295 baris | Langkah `007` **bertabrakan nama**. Penjaga ronde 3 akan melewatinya sebagai "sudah ada" — keberadaan benar, **bentuk berbeda** — dan mencatat langkahnya sukses. Inilah **separuh kelemahan yang sengaja dibiarkan terbuka ronde 3**, kini terbukti nyata, bukan hipotetis. Keputusan **v** |
| `OS_AKSEPTASI_KLAIM_LIFE` punya **62** kolom, bukan 55; `STS_REJECT NUMBER(38,0)`, `WPC DATE`, `CLAIM_RETRO NUMBER` — ketiganya ditebak **teks** | Ronde 3 §3-7 kini dapat dikerjakan. Keputusan **w** |
| Oracle pengembangan **12.2.0.1**; `COMPATIBLE` tidak terbaca | Batas 30 byte yang mendasari keputusan **j** turun menjadi `[dugaan]`. Ralat sudah ditulis di STRUKTUR. **j tidak dibatalkan sendiri oleh executor** |
| `JSON_OFFER_LIFE` memuat ketiga `*_DATE` bertipe `DATE`, plus `STATUS` dan `OLDID` | **menguatkan e′** — rantai buktinya kini sampai ke katalog, bukan berhenti di XML |
| Sumber procedure penomoran terbaca utuh: **nol `COMMIT`**, kunci `FOR UPDATE` atas `(CLASS, JENIS, TAHUN)`, format dirakit **pemanggil** | §5-4 terpenuhi; **o1–o3** kini dapat diputuskan tanpa menebak. ⚠️ AC 3 tiket 02 harus mencakup **perakitan format**, yang ternyata ada di rule Pega, bukan di procedure |

---

## Implementasi — ronde 5, 26 September 2026

**Gerbang: G0 dan G2 terbuka; G1 dan G3 tertutup.** Sesi ini yang pertama punya pekerjaan kode
nyata tanpa menunggu Oracle. G1 tetap tertutup — DBA belum membuat user kosong, dan `POOLDATA`
bukan skema uji. Karena itu **tiket tetap `claimed`**, dan kalimat ini ditulis di awal.

**40 dari 53 AC** (tidak berubah — yang dikerjakan ronde ini bukan AC baru melainkan pembetulan
dasar yang dipijak AC lain), diverifikasi **91 test Go yang benar-benar berjalan** (dari 83) dan
**20 test bertag `db`** (dari 19). Angka diambil **sesudah** perbaikan `/code-review`.

### Empat temuan verifikasi ronde 4

| # | Temuan | Keadaan |
| ---: | --- | --- |
| 1 | ⭐ Tipe kolom warisan kini diketahui; tiga di antaranya ditebak keliru | ✅ ditutup — lihat di bawah |
| 2 | ⭐ Pra-terbang bentuk (butir x) | ✅ ditutup penuh, termasuk penolakan di `JalankanMigrasi` |
| 3 | ⚠️ `[dugaan]` lock `FOR UPDATE` di tiket 02 | ✅ naik ke `[terverifikasi]`, dengan dua fakta tambahan |
| 4 | ℹ️ `PANDUAN-MENJALANKAN.txt` bab 5 | ✅ sudah diralat sesi verifikasi; isinya cocok dengan `vite.config.ts` yang kini punya proxy |
| 5 | ⚠️ Jalur mundur tidak punya pintu masuk manusia | ✅ flag **`-migrate-down`** — lihat di bawah |

### Temuan 1 — tebakan yang dibongkar katalog

Sampai hari ini tipe kolom `OS_AKSEPTASI_KLAIM_LIFE` **ditebak dari nama kolom**, dan tabel TIRUAN
di skema uji dibuat dari tebakan yang sama. ⛔ Akibatnya test pulang-pergi **menguji tiruan
terhadap dirinya sendiri** — itulah temuan ronde 3 §3-7 yang menunggu data DBA selama dua ronde.

Sekarang `kolomWarisan` adalah **satu tabel 62 kolom** yang diturunkan dari katalog, dan
`kolomAngkaLama`, `kolomTanggalLama`, serta `TipeKolomBarisLama` **semuanya diturunkan darinya** —
bukan lagi tiga daftar yang bisa berselisih. Tiga tebakan yang meleset: `STS_REJECT` `NUMBER(38,0)`,
`CLAIM_RETRO` `NUMBER`, `WPC` `DATE`. Ketiganya dulu teks, sehingga dibaca **tanpa `TO_CHAR`** dan
bentuknya bergantung setelan NLS sesi.

Tabel tiruan kini **62 kolom** (dari 55): tujuh kolom yang tidak pernah ditulis rule ikut dibuat,
supaya `INSERT` 55 kolom diuji terhadap tabel yang berbentuk sama dengan produksi.

⛔ **Yang ikut tersingkap, dan tidak disebut brief:** fixture `contohBaris` mengisi `CLAIM_RETRO`
dengan teks `"UJI-RETRO"`. Selama tiruan bertipe `VARCHAR2` semuanya itu lolos; terhadap tiruan
berbentuk benar, Oracle menjawab **`ORA-01722`**. Fixture diperbaiki, dan
`TestNilaiKolomAngkaSelaluAngkaAtauKosong` kini mengunci bahwa setiap nilai yang menuju kolom angka
berupa angka atau kosong.

### Temuan 2 — pra-terbang bentuk

`[keputusan work owner butir x]`. `JalankanMigrasi` kini memeriksa **seluruh langkah sebelum satu
pernyataan pun dikirim**: tiap `CREATE TABLE` yang tabelnya sudah ada dibandingkan kolomnya dengan
`SYS.ALL_TAB_COLUMNS`. Berbeda → migrasi **berhenti**, pesannya menyebut tabel dan kedua arah
selisihnya, dan **nol langkah tercatat** di `T_MIGRASI`. `ObjekSudahAda` hanya memuat tabel yang
bentuknya cocok.

⭐ Pemecah kolom DDL **diangkat dari test ke kode** (`KolomCreateTable`), dan test struktur kini
memakai pemecah yang sama — sehingga yang diuji adalah alat yang sungguh dipakai saat migrasi
berjalan, bukan kembarannya.

Ini menutup **sisa** kelemahan ronde 3, yang sengaja dibiarkan separuh: keberadaan objek dibuktikan,
bentuknya tidak.

### Keputusan §2 sesi ini

| | Keadaan | Yang dikerjakan |
| --- | --- | --- |
| **x** | ⭐ `[DIPUTUSKAN]` — sahkan penuh | pra-terbang lengkap + test murni + test db |
| **butir 5 brief** | ⛔ **terlewat di lintasan pertama** | baris ini **ditambahkan ke brief sesudah commit Langkah 0**, jadi tidak pernah terbaca saat merencanakan sesi. Ditemukan `/code-review` dan **dikerjakan penuh** sebelum commit |
| **w** | `[DIPUTUSKAN]` — ikuti warisan | `003` `WPC DATE`; STRUKTUR diberi ralat bertanggal. ⚠️ `CLAIM_RETRO` **tidak disentuh**: tipenya diketahui (`NUMBER`), **artinya tidak** — uang atau perbandingan belum dijawab, jadi `002` tetap `VARCHAR2(64)` dan itu `[terbuka]` |
| **v** | `[DIPUTUSKAN]` — **v1** | `DOCUMENT_CLAIM` → **`T_CLAIMLF_DOCUMENT`**; berkas `007` dan `SEQ_` ikut; tabel warisan 295 baris **tidak disentuh sama sekali**; STRUKTUR diberi blok ralat |
| **y, j** | `[terbuka]` / `[USULAN]` | tidak disentuh — pemiliknya tiket 15 dan work owner |
| **o1–o3** | `[USULAN]` | G3 tertutup. Dua fakta baru dicatat di tiket 02 |

### Yang dijalankan, dan hasilnya

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l` · `go build ./...` | lulus, nol berkas |
| `go test ./...` | lulus, **89 test** (dari 83), nol FAIL |
| `go test -tags=db ./internal/...` | **20 MELEWATI** (dari 19) — ⛔ SKIP, bukan PASS |
| `npm run typecheck` · `npm test` · `npm run build` | lulus, 5 test, 87 modul |

⭐ **Tiap penjaga baru diuji gagal dulu pada kasus buruknya**, lalu dipulihkan: `WPC` dikembalikan ke
teks → gagal; satu baris katalog dihapus dari tabel Go → gagal; kolom angka diisi teks → gagal;
arah "kolom lebih" pada pembanding bentuk dilumpuhkan → gagal.

⚠️ **Dua kali penjaga menangkap asumsi saya sendiri:** cacah `CREATE TABLE` yang saya tulis 8
ternyata **7** (yang kedelapan, `T_MIGRASI`, dibuat di luar berkas migrasi), dan cacah pemanggil
`skemauji.Buka()` bertambah menjadi 7 karena test db baru. Keduanya dikunci angkanya.


### ⛔ Temuan 5 — jalur mundur akhirnya punya pintu masuk manusia

Sampai hari ini `BongkarMigrasi` **hanya** punya satu pemanggil: skema uji. Siapa pun yang ingin
membongkar skema uji sendiri terpaksa menyalin isi berkas `*_down.sql` ke sqlplus — dan itu
**melewati pengaman `T_MIGRASI`**, yang membongkar hanya langkah yang benar-benar tercatat selesai.
Catatan migrasi karena itu bisa berbohong tanpa ada yang tahu.

Sekarang ada **`go run ./cmd/api -migrate-down`**, dipagari **sama persis** dengan test bertag `db`:
menolak `IS_PEGA_PROD=true`, menolak tanpa `ORACLE_SKEMA_UJI=true`, menolak skema yang memuat
`POOLDATA`. ⭐ Pagarnya **diangkat ke `internal/config`** supaya hanya ada **satu** — `cmd/api` tidak
boleh mengimpor paket penunjang test, sehingga menaruhnya di `skemauji` akan memaksa jalur kedua
menyalin logikanya, dan salinan yang berselisih adalah cara paling mudah kehilangan tabel.

### ⛔ Yang ditemukan `/code-review`, dan diperbaiki sebelum commit

| # | Temuan | Perbaikan |
| ---: | --- | --- |
| 1 | ⛔ **Butir 5 brief tidak dikerjakan sama sekali** | dikerjakan penuh: flag, pagar terangkat, test murni, paragraf panduan |
| 2 | ⛔ **Jaminan `CLAIM_RETRO` KOSONG** — `BarisLamaDari` tidak pernah mengisinya, jadi test melewatinya lewat cabang "kosong itu sah". Bab ini sempat mengklaim "mengunci setiap nilai yang menuju kolom angka" | klaimnya diralat di sini dan di kepala test; yang terkunci memang **fixture**, bukan kode. `[terbuka]` baru ditulis — lihat di bawah |
| 3 | ⛔ **Pra-terbang memeriksa langkah yang SUDAH tercatat**, sehingga satu kolom audit yang ditambahkan DBA pada tabel lama akan **menggagalkan seluruh migrasi berikutnya**, termasuk langkah baru yang tidak ada hubungannya | langkah yang `selesai` dilewati pra-terbang |
| 4 | ⛔ **Lubang: `CREATE TABLE` yang tidak terurai lolos tanpa pemeriksaan bentuk** — regex menuntut `
)`, sehingga tabel yang ditulis satu baris jatuh ke pemeriksaan keberadaan saja | `tampakCreateTabel` menolak pernyataan yang terlihat `CREATE TABLE` tetapi tidak terurai |
| 5 | ⚠️ `README-BACA-DULU.md` dan komentar kepala `barislamakolom.go` masih menulis **55 kolom** | → 62, dengan 55 disebut sebagai yang ditulis rule |
| 6 | ⚠️ Kepala `tipewarisan_test.go` mengklaim menutup "menguji tiruan terhadap dirinya sendiri" | diralat: yang dibuktikan hanya **dokumen == kode**; bahwa dokumen == Oracle bersandar pada satu pembacaan katalog DEV yang belum dikonfirmasi DBA |

### ⛔ `[terbuka]` baru — `STS_REJECT`: ADR-U-0022 lawan katalog

ADR-U-0022 menetapkan **kode tetap teks**, dan `models` mengunci `KodeStatus` sebagai teks bebas
lewat `TestKodeStatusTidakPernahJadiBilangan`. Katalog menyebut `OS_AKSEPTASI_KLAIM_LIFE.STS_REJECT`
bertipe **`NUMBER(38,0)`**. Keduanya tidak dapat benar sekaligus: `Simpan` mem-bind teks itu apa
adanya, sehingga kode status non-angka akan dijawab **`ORA-01722`** di lapangan.

⛔ **Executor tidak memutuskannya.** Menambah penolakan di jalur tulis berarti diam-diam memihak
katalog dan melanggar ADR; membiarkannya berarti menunggu galat di lapangan. Pemilik: **work owner**.
Sampai dijawab, yang dilakukan test hanyalah mencegah fixture menambah kasus baru yang pasti gagal.

### Tuduhan tinjauan yang saya TOLAK sesudah memeriksa premisnya

⛔ *"`003` diubah di tempat tanpa jalur naik; skema uji lama tetap memegang `WPC VARCHAR2(32)`."*
**Premisnya tidak ada.** Migrasi belum pernah berjalan di Oracle mana pun — enam sesi berturut-turut
— dan katalog instance pengembangan yang dibaca 26-09 tidak memuat `T_MIGRASI`. Brief §2 l
menetapkan berkas `001`–`008` boleh disunting langsung **selama** keadaan itu berlaku, dan Langkah C
brief ronde 5 mengulanginya kata demi kata. Yang benar dari tuduhan itu hanya catatannya: pra-terbang
membandingkan **nama kolom**, bukan tipe — dan itu memang tertulis di komentarnya sendiri.

### ⚠️ Sisa rename v1 yang BUKAN milik tiket ini

`T_CLAIMLF_DOCUMENT` sudah konsisten di seluruh SQL, kode, dan test. Nama lama masih tersisa sebagai
**tabel baru** di `spec.md`, `issues/03`, `issues/15`, dan `revisi-penyimpanan-json-dibuang.md`.
Keempatnya **bukan milik executor tiket 14** — executor tidak menyunting tiket lain maupun berkas
tersegel. Diserahkan ke work owner. ✅ Yang memang **harus** tetap `DOCUMENT_CLAIM`: nama kelas Pega
`ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`, dan seluruh rujukan ke tabel **warisan**.

### TELEMETRI EKSEKUSI

| Besaran | Nilai |
| --- | --- |
| Gerbang terbuka | **G0, G2** · G1 tertutup (DBA belum membuat user kosong) · G3 tertutup |
| Keputusan diterapkan | **x** (penuh), **w** (`WPC DATE`; `CLAIM_RETRO` ditahan), **v1** |
| Temuan §3 ditutup | **5 / 5** |
| AC tiket 14 | **40 / 53**; daftar terbuka **13** = jumlah kotak `[ ]` |
| AC tiket 01 | **4 / 7** — tidak disentuh |
| **Kolom tabel tiruan warisan** | **62** *(harus 62)* |
| **Selisih peta tipe lawan katalog** | **0** *(harus 0)* |
| `-migrate` sungguhan | ⛔ **tidak dijalankan** — G1 tertutup |
| Objek DB per `OBJECT_TYPE` | ⛔ **kosong** — tidak ada migrasi yang berjalan |
| Test | **91 PASS · 0 FAIL** (tanpa tag) · **20 SKIP** bertag `db` — ⛔ **SKIP, bukan PASS** · 5 test JS · 87 modul |
| Sub-agen review | standards **118.845 token / 19 panggilan / 264 detik**; spec **125.788 token / 29 panggilan / 313 detik** |
| Token sesi utama · biaya · jam dinding | ⛔ **tidak diukur** |
