# 14: Skema relasional klaim (6 tabel, berakar di `T_WORK_CLAIM`) + migrasi — **PREFACTOR**

**Status:** ready-for-agent

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

- [ ] ⚠️ Skema klaim **relasional penuh**: setiap atribut menjadi **kolom bernama**. Test yang
      menemukan kolom JSON menyimpan atribut klaim **gagal**. *(AC 31 spec; penyimpangan sadar 1)*
- [ ] ⚠️ **Enam tabel klaim** ada dengan PK dan FK sesuai diagram — `T_GENERAL_CLAIM`,
      `T_CLAIMLF_PREMIUMLIST_DETAIL`, `T_CLAIMLF_ADJUSTMENT`, `T_CLAIMLF_ADJUSTMENT_SPREADING`,
      `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`, `DOCUMENT_CLAIM`. **REVISI 2026-09-18:** ⛔ **tidak
      ada** `T_CLAIMLF_POLICY` maupun `T_CLAIMLF_MARKETING` — keduanya **dihapus**, bukan diganti
      nama. Test yang menemukan salah satunya **gagal**. *(AC 48 spec; penyimpangan sadar 8)*
- [ ] ⚠️ **REVISI 2026-09-18** — **tidak** seluruh FK `ON DELETE CASCADE`. Yang **CASCADE**: relasi
      **3 · 4 · 5 · 6** (klaim → peserta → adjustment → spreading → spreading retro). Yang **di Go**:
      relasi **1** (`COVER_KEY`) dan **11** (`KOMITE_ID`, penunjuk). Yang ⚠️ `[terbuka]`: relasi
      **7** (`DOCUMENT_CLAIM`, kini lintas-lini — induk beda tabel per lini) dan relasi **2**
      (belum punya kolom). Lihat tabel sebelas relasi di atas. **Jangan menyeragamkan sendiri.**
- [ ] ⚠️ FK `T_CLAIMLF_ADJUSTMENT` menunjuk **`T_CLAIMLF_PREMIUMLIST_DETAIL.ID`**, **bukan**
      `T_GENERAL_CLAIM.ID`. Test yang menemukan adjustment menggantung pada header **gagal**.
      *(AC 33 spec; penyimpangan sadar 2)*
- [ ] ⚠️ FK `DOCUMENT_CLAIM` menunjuk **`T_CLAIMLF_PREMIUMLIST_DETAIL.ID`**. *(AC 44 spec;
      penyimpangan sadar 5)*
- [ ] ⚠️ **`T_WORK_CLAIM` ada sebagai tabel mandiri** — bukan anak `T_GENERAL_CLAIM`, dan keadaan tangga
      **tidak** menjadi kolom header klaim. *(AC 46 spec; penyimpangan sadar 6)*
- [ ] Header memuat keempat field `PremiumListSummary` (`CLAIM_NO`, `PL_NUMBER`, `RISLIPRNM`,
      `BUSINESS_NAME`) beserta `CASEID` dan `CLAIM_RETRO`. *(AC 36 spec)*
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS, MENUNGGU JAWABAN.**
      `T_CLAIM_POLICY` **dihapus**, jadi tidak ada tabel yang memuat kesembilan kolom itu. **27 dari
      32** kolomnya dibaca dari `T_PREMIUM_LIST`; **lima** ⚠️ `[terbuka]` **kehilangan rumah** —
      `TEAM_GROUP`, `BUSINESS_ID`, `TANGGAL_RESPON`, `TANGGAL_REALISASI`,
      `TANGGAL_KONFIRMASI_BALIK` (§Blocker). `PL_NUMBER` tetap tidak ada, tetapi karena **diganti
      nama menjadi `POLICY_NO`** di `T_GENERAL_CLAIM` — bukan karena tabel ini menolaknya.
      **Jangan tebak rumah kelima kolom itu.** *(AC 37 spec)*
- [ ] `T_CLAIMLF_PREMIUMLIST_DETAIL` memuat **kesembilan kolom tambahan** hasil audit — termasuk
      `IS_CHECK` dan ketiga tanggal per peserta. *(AC 39, 40, 41 spec)*
- [ ] `T_CLAIMLF_ADJUSTMENT` memuat `CLAIM_AMOUNT`, `STS_REJECT`, `ACCEPTEDNO`, `ACCEPTATION_DATE`.
      *(AC 42, 43 spec)*
- [ ] `T_CLAIMLF_ADJUSTMENT` memuat **ketiga kolom bank** — `NAME_OF_BANK`, `ID_BANK`, `ACCOUNT_NO` —
      dan migrasi mengisinya dari kolom warisan `NAME_OF_BANK`, `IDBANK`, `ACCOUNTNO`.
      *(AC 56 spec; `[terverifikasi]`)*
- [ ] ⚠️ Seluruh uang dan share bertipe **desimal presisi arbitrer**; seluruh tanggal **`DATE`**;
      seluruh kolom **nullable**; identitas dari **sequence**. Test yang menemukan kolom uang
      bertipe teks atau melewati `float` **gagal**. *(AC 50 spec; **ADR-0003**, **ADR-0006**;
      penyimpangan sadar 7)*
- [ ] Seluruh klaim Life terbawa **beserta seluruh baris adjustment**-nya — bukan hanya keadaan
      terakhir. Jumlah baris per klaim setelah migrasi **sama** dengan sebelumnya. *(**ADR-0011**)*
- [ ] Setiap baris adjustment hasil migrasi **menunjuk peserta yang benar**. *(AC 34 spec)*
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS.** Tidak ada
      `T_CLAIM_POLICY` untuk dinormalkan ke dalamnya; atribut polis **tidak dipindahkan ke mana
      pun** — ia **dibaca hidup** dari tabel polis ⚠️ **penyimpangan sadar**. Yang **tetap
      mengikat**: bila atribut polis **berbeda antar baris peserta** dalam satu klaim, migrasi
      **melaporkannya** dan **tidak** diam-diam memilih salah satu. *(AC 51 spec)*
- [ ] ⚠️ Migrasi **melaporkan** bahwa data lama ber-`ACCEPTATION_DATE` = waktu insert dan
      `STS_REJECT` = `0` karena **di-hardcode** di sumbernya, bukan karena nilainya sebenarnya.
      Migrasi **tidak mengarang** tanggal akseptasi. *(AC 52 spec; penyimpangan sadar 4)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
      tepat**, bukan dengan toleransi. *(AC 53 spec; **ADR-0003**)*
- [ ] Tanggal yang berupa teks menjadi `DATE` **tanpa pergeseran zona waktu**; yang **tidak dapat
      diurai dilaporkan**, bukan didiamkan. *(AC 54 spec)*
- [ ] ⚠️ Migrasi **tidak mengambil apa pun** dari `m_product_life.JSONDATA`; `PRODUCT_NAME` dan
      `PRODUCT_NAME_ID` berasal dari **`product_life` relasional**. *(AC 38 spec)*
- [ ] Penomoran klaim **tidak melompat dan tidak mengulang** setelah migrasi; sequence pindah dengan
      **nilai berjalan yang benar**. *(**ADR-0006**)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

### Kolom yang bertambah dan yang pindah ⚠️ BARU 2026-09-18

- [ ] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `CREATE_OP`, `CREATE_OP_NAME`, `TGL_UPDATE` —
      ketiganya **pindah ke `T_WORK_CLAIM`**. Test yang menemukannya di header klaim **gagal**.
- [ ] ⚠️ `T_WORK_CLAIM` memuat `CREATE_OP`, `CREATE_OP_NAME`, `TGL_UPDATE`, **`CASEID`**
      (seluruhnya pindahan) dan `LINI`.
- [ ] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `CASEID` — pindah ke `T_WORK_CLAIM`. Test yang
      menemukannya di header klaim **gagal**. *(DIPUTUSKAN 2026-09-18; `[keputusan work owner]`)*
- [ ] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `PL_NUMBER`; kolom itu **diganti nama** menjadi
      **`POLICY_NO`**. Test yang menemukan `PL_NUMBER` di header klaim **gagal**.
- [ ] ⚠️ `T_GENERAL_CLAIM` memuat ketiga **penunjuk polis** — `CASEID_POLICY`, `POLICY_NO`,
      `ENDORSMENT_NO` — dan atribut polis **tidak** disalin ke klaim. Test yang menemukan salinan
      atribut polis di tabel klaim **gagal**.
- [ ] ⚠️ `[terbuka]` **Isi ketiga penunjuk polis belum dapat ditulis** — tidak satu pun menunjuk
      kolom yang ada di sisi polis (§Blocker). Kolomnya dibuat; **pengisiannya** menunggu keputusan
      work owner. **Jangan tebak.**
- [ ] ⚠️ Kolom `LINI` **ADA** pada `T_WORK_CLAIM`, dan untuk Life isinya **konstanta lini Life**.
      `[terbuka — Non-Life]` hanya **daftar nilai enum lintas-lini**, yang ditetapkan saat konteks
      Non-Life digarap — **bukan** keberadaan kolomnya. *(DIPUTUSKAN 2026-09-18)*

### Identitas: shared PK dan nomor bisnis berformat ⚠️ BARU 2026-09-18

- [ ] ⚠️ **`T_GENERAL_CLAIM.ID` sama persis dengan `T_WORK_CLAIM.ID` baris klaim** — **shared
      primary key**, 1:1, **tanpa kolom penyambung**. `T_GENERAL_CLAIM.ID` sekaligus PK **dan** FK
      ke `T_WORK_CLAIM.ID`. *(relasi 2; `[keputusan work owner]`)*
- [ ] ⚠️ ⛔ **Tidak ada kolom `WORK_CLAIM_ID`** di mana pun — pada `T_GENERAL_KOMITE` maupun
      tabel lain. Hubungan ke kasus komite juga **shared PK**: `T_GENERAL_KOMITE.ID` =
      `T_WORK_CLAIM.ID` baris komite. Test yang menemukan kolom `WORK_CLAIM_ID` **gagal**.
      *(relasi 8; `[keputusan work owner]`)*
- [ ] ⚠️ **`T_GENERAL_KOMITE.ADJUSTMENT_ID` TETAP ADA** — penutup lingkar ke baris adjustment; ia
      **bukan** bagian shared PK dan **tidak** ikut dibuang. *(relasi 10)*
- [ ] ⚠️ **`T_WORK_CLAIM.ID` bertipe teks berformat** — baris klaim `CLM-xxxxxx`, baris komite
      `KMT-xxxxxx`. **Bukan angka sequence.** Test yang menemukan tipe numerik **gagal**.
      *(`[keputusan work owner]`)*
- [ ] ⚠️ `T_WORK_CLAIM.COVER_KEY`, `T_GENERAL_CLAIM.ID`, `T_GENERAL_KOMITE.ID`, dan
      `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` **bertipe sama** dengan `T_WORK_CLAIM.ID` — teks berformat.
- [ ] ⚠️ **Penyimpangan sadar dari ADR-0006 dicatat di artefak**, tidak dilanggar diam-diam:
      identitas `T_WORK_CLAIM` dan kedua tabel ber-shared-PK adalah **nomor bisnis berformat**,
      bukan sequence. ADR-0006 **tetap berlaku** untuk `T_CLAIMLF_*` dan `DOCUMENT_CLAIM`.
- [ ] ⚠️ `[terbuka]` **Generator nomor `CLM-`/`KMT-` belum ditetapkan** — siapa yang membuatnya,
      apakah ada sequence di belakang prefiks, apakah di-reset per tahun. Pemilik **DBA / work
      owner**. Kolomnya dibuat; **pembangkitannya** menunggu jawaban. **Jangan tebak.**

### Spreading adjustment + rujukan Komite ⚠️ BARU 2026-09-16

- [ ] ⚠️ **`T_CLAIMLF_ADJUSTMENT_SPREADING` ada**, dengan FK **`ADJUSTMENT_ID`** → `T_CLAIMLF_ADJUSTMENT.ID`
      dan **`ON DELETE CASCADE`**. Test yang menemukannya menggantung pada peserta atau pada header
      klaim **gagal**. *(`[terverifikasi]` `SpreadingClaimLife_Act` mengisi `.SpreadingList` pada
      baris adjustment; **AC 58 spec**)*
- [ ] ⚠️ **`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` ada**, dengan FK **`SPREADING_ID`** →
      `T_CLAIMLF_ADJUSTMENT_SPREADING.ID` dan **`ON DELETE CASCADE`**. *(AC 58 spec)*
- [ ] Pohon klaim berkedalaman **lima tingkat** — klaim → peserta → adjustment → spreading →
      spreading retro — dan menghapus klaim **mengkaskade sampai tingkat terdalam**. Test wajib
      memeriksa **cicit** (`_SPREADING_RETRO`) ikut hilang. *(AC 60 spec)*
- [ ] `T_CLAIMLF_ADJUSTMENT_SPREADING` memuat `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`,
      `TREATY_YEAR_LIFE`, `RETROCADED_SHARE`, `RATE`, `IDR`, `USD`, `CURRENCY`.
- [ ] `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` memuat `REINSURER_NAME`, `PERCENT_SHARE`, `AMOUNT`,
      `RATE`, `PREMIUM_SPREADED_GROSS`, `PREMIUM_SPREADED_NET`, `COMMISION`, `OVR_COMM`,
      `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`. *(`[terverifikasi]`
      `Claim Life/Section/RetroDetailClaimLife.xml` + `SpreadingClaimLife_Act`)*
- [ ] ⚠️ Seluruh kolom uang dan persen pada kedua tabel bertipe **desimal presisi arbitrer**;
      **tidak** melewati `float`. *(**ADR-0003**)*
- [ ] ⚠️ Nilai spreading **dibekukan**: perubahan master treaty setelah adjustment tersimpan
      **tidak mengubah** angka yang sudah ada. *(penyimpangan sadar — spreading disimpan)*
- [ ] ⚠️ **`T_CLAIMLF_ADJUSTMENT` memuat kolom `KOMITE_ID`**, **nullable** dan **ber-index** — `NULL`
      bila baris belum pernah dikirim ke Komite. *(AC 61 spec; penyimpangan sadar — rujukan, bukan salinan)*
- [ ] ⚠️ **Tidak ada tabel `T_CLAIMLF_ADJUSTMENT_KOMITE`.** Roster dan keputusan per anggota **tidak**
      disimpan di Claim Life. Test yang menemukan tabel itu **gagal**. *(AC 61 spec; penyimpangan sadar)*
- [ ] ⚠️ Rujukan ke Komite memakai **`KOMITE_ID`**, bukan **indeks posisi**. Test yang menemukan
      padanan `IndexPremiumList`/`IndexAdjustment` sebagai kunci rujukan **gagal**.
      *(AC 62 spec; `[terverifikasi]` `CreateKMTLife_Act` memakai `.pxListSubscript`; penyimpangan sadar)*
- [ ] Migrasi **membongkar** `SpreadingList` dan `RetroLifeList` dari data lama ke kedua tabel,
      dan setiap baris dapat ditelusuri ke **baris adjustment yang benar**.
- [ ] Setiap FK baru (`ADJUSTMENT_ID`, `SPREADING_ID`) **ber-index**.
- [ ] ⚠️ `KOMITE_ID` **ber-index** juga, dan berisi **identitas kasus komite** = `T_WORK_CLAIM.ID`
      baris komite. Test yang menemukannya menunjuk `T_GENERAL_KOMITE.ID` **gagal**.
      *(§`KOMITE_ID` pada `T_CLAIMLF_ADJUSTMENT`; `[keputusan work owner]` REVISI 2026-09-17)*
- [ ] ✅ Tipe `KOMITE_ID` **SUDAH DITETAPKAN 2026-09-18** — teks berformat `KMT-xxxxxx`, mengikuti
      `T_WORK_CLAIM.ID`. *(`[keputusan work owner]`)*
- [ ] ⚠️ `[terbuka]` **Apakah `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)`**
      atau dibiarkan tanpa constraint — **belum diputuskan**. Pemilik **DBA / work owner**. Tiket
      ini **tidak dinyatakan selesai** sebelum jawabannya ada. **Jangan tebak.**

### Koreksi: `OS_AKSEPTASI_KLAIM_LIFE` TETAP ditulis ⚠️ 2026-09-16

⚠️ `[keputusan work owner]` **Keputusan sebelumnya dibalik.** Sistem baru menulis **dua tempat**:
tabel relasional baru **dan** `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE`, karena hilir (Arasapas /
produksi) masih membaca dari sana. **Yang dibuang hanya JSON** (`JSON_KLAIM` / serialisasi
`ClaimData`).

- [ ] ⚠️ Menyimpan klaim menulis **tabel relasional baru** *dan* `INSERT` flat ke
      `OS_AKSEPTASI_KLAIM_LIFE`. Test yang menuntut `OS_AKSEPTASI_KLAIM_LIFE` **tidak** ditulis
      adalah **keliru** dan harus dibalik.
- [ ] ⚠️ Yang **tetap dibuang** hanya **blob JSON** — `JSON_KLAIM` dan serialisasi `ClaimData`.
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
