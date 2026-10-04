# Spec — Claim Life: penyimpanan relasional (JSON dibuang)

> ⛔ **SUPERSEDED — nama tabel di berkas ini SUDAH TIDAK BERLAKU.** Yang berlaku ada di
> `revisi-penyimpanan-json-dibuang.md` dan `spec.md`. Jangan dipakai sebagai acuan.
> *(Ditambahkan 2026-09-18, saat tujuh tabel diganti nama. 14 kemunculan nama lama di berkas ini
> DIBIARKAN dengan sengaja — bukan pekerjaan yang terlewat.)*

> ⛔ **DIGANTIKAN 2026-09-16 — JANGAN DIPAKAI SEBAGAI SUMBER.** Dokumen ini adalah draf antara yang
> ditulis sebelum work owner menutup lima OQ audit. Isinya **berbeda** dari keputusan final dalam
> tiga hal: ia mengenal **enam tabel** tanpa `T_WORK_CLAIM`, menandai tiga OQ sebagai pemblokir, dan
> berstatus `needs-info`.
> **Sumber yang berlaku:** `revisi-penyimpanan-json-dibuang.md` (keputusan) dan `spec.md` §2b
> (spesifikasi). Disimpan hanya sebagai jejak.

Status: ~~needs-info~~ → **superseded**
Konteks: `claim-life` — **revisi bagian penyimpanan** dari `spec.md` (yang tetap berlaku untuk
perilaku, peran, mesin status, dan integrasi).
Tanggal: 2026-09-16
Sumber: `.scratch/claim-life/revisi-penyimpanan-json-dibuang.md` (`[keputusan work owner]`),
audit sensus property 2026-09-16 (sesi ini), `.scratch/claim-life/spec.md`,
`docs/adr/ADR-0001`–`ADR-0015`, korpus `Claim Life/`
Skill: `/mattpocock-skills:to-spec`

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[fakta bisnis — work owner]`; `[data DBA]`; `[terbuka]` = OQ.
> **Identitas rule wajib menyertakan class** — nama sama di class berbeda = rule berbeda.

> **Sumber tunggal.** Korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY) dan artefak di
> `OUTPUT_HASIL_RNM\`, sekaligus repo target tunggal. ⛔ `D:\XML\nusantara-re\` di-blacklist.

> **Hubungan dengan `spec.md`.** Dokumen ini **menggantikan** pernyataan penyimpanan di `spec.md`
> (§2 penempatan modul, §14 migrasi, dan seluruh rujukan ke blob JSON). Yang **tidak** disentuh dan
> tetap berlaku: mesin status per baris (§3, **ADR-0011**), peran (§4), `Type` (§5), validasi DOL
> (§6), kontrak Komite (§8), uang (§9), penomoran (§10), efek keluar (§12), jejak audit (§13), dan
> aturan **peserta hidup** (§16).

---

## Problem Statement

Klaim jiwa hari ini disimpan sebagai **satu dokumen JSON** berisi seluruh isi `ClaimData` —
polis, marketing, daftar peserta, dan daftar adjustment tiap peserta — di samping satu tabel flat
warisan. Empat hal membuat bentuk itu tidak dapat dipertahankan:

1. **Bentuknya tidak dapat diperiksa.** Tidak ada kolom, tidak ada tipe, tidak ada `NOT NULL`. Nilai
   uang dan tanggal masuk sebagai teks. Tidak ada cara memeriksa bahwa sebuah klaim utuh selain
   membukanya.

2. **Sistem hilir dipaksa ikut mengurai JSON.** Konversi ke Arasapas/produksi hari ini mengirim
   payload JSON, sehingga setiap pembaca harus tahu bentuk dalamnya.

3. **Tabel flat warisan mencampur tiga tingkat jadi satu baris.**
   `[terverifikasi]` `Claim Life/RDBList/InsertJsonKlaimLife_sql.xml`
   (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!INSERTJSONKLAIMLIFE_SQL` / `RULE-CONNECT-SQL`)
   menulis **51 kolom** ke `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`, satu baris per peserta, dengan
   seluruh atribut polis (ceding, source of business, retro, security reinsurer, produk) **diulang
   pada setiap peserta**. Mengubah satu atribut polis berarti menyentuh setiap baris peserta.

4. **Nama berbohong, dan itu sudah menyesatkan sekali.** Rule bernama `InsertJsonKlaimLife_sql`
   ternyata `INSERT` **flat**, bukan JSON. Sebaliknya
   `Claim Life/RDBList/GetJsonProductLife.xml` (`ASM-FW-GISFW-INT-TREATYYEAR_LIFE` /
   `ASM!GETJSONPRODUCTLIFE` / `RULE-CONNECT-SQL`) benar-benar membaca JSON — `json_table` atas
   `m_product_life.JSONDATA`, path `$.OutwardList[*]`.

## Solution

Menyimpan klaim jiwa dalam **enam tabel relasional baru, milik konteks klaim sendiri**. Tidak ada
blob JSON di sistem baru, dan sistem hilir membaca **langsung dari tabel**, bukan dari payload.

`[keputusan work owner]` Tabel existing (`OS_AKSEPTASI_KLAIM_LIFE`, `M_LIFE_PREMIUM_DETAIL`)
**tidak dipakai untuk menyimpan klaim baru**; keduanya tetap **dibaca** sebagai sumber snapshot.

### Delapan penyimpangan sadar

| # | Penyimpangan | Alasan | Sumber |
| --- | --- | --- | --- |
| 1 | **Seluruh JSON dibuang** — `JSON_KLAIM`, serialisasi `ClaimData`, dan payload JSON ke Arasapas | bentuk JSON tidak dapat divalidasi, tidak dapat di-index, dan memaksa hilir mengurai | `[keputusan work owner]` |
| 2 | **Klaim punya tabelnya sendiri**, terpisah dari tabel existing | klaim adalah entitas tersendiri, bukan tempelan pada akseptasi | `[keputusan work owner]` |
| 3 | **Atribut polis disimpan sekali**, bukan diulang pada tiap peserta | existing mengulang 14 atribut polis di setiap baris peserta | `[keputusan work owner]` |
| 4 | ⚠️ **`AdjustmentList` digantungkan pada peserta**, bukan pada klaim | korpus membuktikan ia bersarang di bawah `PremiumListDetail` — skema kandidat semula salah induk | `[terverifikasi]` |
| 5 | **Dokumen per peserta jadi tabel sendiri** | ia **gerbang simpan**, bukan pelengkap | `[terverifikasi]` |
| 6 | **Kaskade hapus + popup konfirmasi** | pola yang sudah ditetapkan di Master Contract Retro Life & Treaty Contract Out | `[keputusan work owner]` |
| 7 | **Uang/share/premi desimal; tanggal `DATE`; seluruh kolom nullable** (wajib-isi di Go) | existing menyimpan keduanya sebagai teks | `[keputusan work owner]`, **ADR-0003** |
| 8 | **Nama jujur** — tidak ada komponen bernama "Json" yang bukan JSON | nama yang berbohong sudah menyesatkan sekali di konteks ini | `[keputusan work owner]` |

### Kontrak batas

| Arah | Isi |
| --- | --- |
| **Keluar** | Arasapas/produksi **`SELECT` langsung** dari tabel klaim baru — **bukan** payload JSON |
| **Masuk** | Snapshot dari data polis life & PremiumList NB/EDM saat klaim dibuat; master (retro, security reinsurer, currency, diagnosa) dibaca saat pengisian |
| **Tidak ada** | Blob JSON, di sisi mana pun |

---

## User Stories

### Menyimpan klaim

1. Sebagai **`ReasLifeAdmin`**, saya ingin klaim yang saya daftarkan tersimpan dalam bentuk yang
   dapat diperiksa, supaya kesalahan ketahuan saat menyimpan, bukan berbulan-bulan kemudian.
2. Sebagai **organisasi**, saya ingin setiap atribut klaim menjadi **kolom bernama**, supaya ia
   dapat dicari, divalidasi, dan dilaporkan.
3. Sebagai **organisasi**, saya ingin satu klaim tersimpan **utuh atau tidak sama sekali**.
4. Sebagai **`ReasLifeAdmin`**, saya ingin **diberi tahu bila penyimpanan gagal**, supaya saya tidak
   mengira klaim tersimpan padahal tidak.
5. Sebagai **admin**, saya ingin **tidak perlu mengetik nomor identitas** baris mana pun.

### Snapshot polis

6. Sebagai **organisasi**, saya ingin data polis yang melekat pada klaim adalah **snapshot saat
   klaim dibuat**, supaya perubahan polis di kemudian hari tidak mengubah klaim yang sudah berjalan.
7. Sebagai **organisasi**, saya ingin atribut polis tersimpan **sekali per klaim**, bukan diulang
   pada setiap peserta.
8. Sebagai **`ReasLifeSPV`**, saya ingin melihat tanggal respon, tanggal realisasi, dan tanggal
   konfirmasi balik polis pada layar klaim.

### Peserta

9. Sebagai **`ReasLifeAdmin`**, saya ingin memilih **peserta mana** dari premium list yang akan
   diklaim, supaya klaim hanya memuat yang relevan.
10. Sebagai **organisasi**, saya ingin pilihan itu **terekam**, supaya dapat dipertanggungjawabkan.
11. Sebagai **`ReasLifeAdmin`**, saya ingin peserta yang sudah dibatalkan atau dihapus **tidak
    muncul** sebagai kandidat klaim.
12. Sebagai **`ReasLifeMedicalAdvisor`**, saya ingin mencatat diagnosa dan kode ICD per peserta.
13. Sebagai **`ReasLifeSPV`**, saya ingin melihat tanggal diterima, tanggal konfirmasi, dan tanggal
    penyelesaian **per peserta**, supaya peserta yang tertahan terlihat.
14. Sebagai **`ReasLifeSPV`**, saya ingin mencatat **rekomendasi** dan **status** per peserta.

### Adjustment

15. Sebagai **`ReasLifeSPV`**, saya ingin setiap peserta punya **daftar adjustment**-nya sendiri,
    supaya putaran keputusan satu peserta tidak tercampur dengan peserta lain.
16. Sebagai **organisasi**, saya ingin seluruh baris adjustment **ikut tersimpan** — bukan hanya
    yang terakhir — supaya riwayat putaran Komite tidak hilang.
17. Sebagai **Finance**, saya ingin nilai share, sum insured, sum reasured, dan claim amount
    **tidak berubah** saat menyeberang batas penyimpanan.

### Dokumen

18. Sebagai **`ReasLifeAdmin`**, saya ingin mengunggah dokumen **per peserta**, supaya kelengkapan
    berkas terlihat per orang.
19. Sebagai **`ReasLifeAdmin`**, saya ingin **dicegah menyimpan ke Outstanding** bila masih ada
    peserta yang dokumennya belum lengkap, dengan pesan yang menyebut **peserta mana**.

### Menghapus

20. Sebagai **`ReasLifeAdmin`**, saya ingin menghapus klaim **beserta seluruh isinya**, supaya tidak
    ada sisa yang menggantung.
21. Sebagai **`ReasLifeAdmin`**, saya ingin **diberi peringatan berisi jumlah baris yang akan ikut
    terhapus** lebih dulu, supaya saya dapat membatalkan.

### Hilir

22. Sebagai **sistem produksi/Arasapas**, saya ingin **membaca langsung dari tabel klaim**, supaya
    saya tidak perlu mengurai dokumen JSON.
23. Sebagai **tim integrasi**, saya ingin bentuk yang saya baca **stabil dan bertipe**, supaya
    perubahan tampilan tidak memecahkan integrasi.

### Migrasi

24. Sebagai **tim migrasi**, saya ingin seluruh klaim lama pindah **tanpa kehilangan satu nilai
    pun**, termasuk **seluruh baris adjustment**.
25. Sebagai **tim migrasi**, saya ingin nilai uang dan tanggal yang hari ini berupa teks menjadi
    **tipe yang semestinya**, dan yang tidak dapat diurai **dilaporkan**, bukan didiamkan.
26. Sebagai **tim migrasi**, saya ingin atribut polis yang hari ini **diulang di setiap baris
    peserta** menjadi **satu baris per klaim**, dengan perbedaan antar baris **dilaporkan** bila ada.

---

## Implementation Decisions

### 1. Bentuk penyimpanan — enam tabel, tiga tingkat

```
T_CLAIMDATA (PK ID)                           ← header klaim (1 baris per klaim)
  ├─ T_CLAIM_POLICY               1:1         ← snapshot polis
  ├─ T_CLAIM_MARKETING            1:1         ← snapshot marketing
  └─ T_CLAIM_PREMIUM_LIST_DETAIL  1:N         ← peserta yang diklaim
        ├─ T_CLAIM_ADJUSTMENT     1:N         ← putaran keputusan per peserta
        └─ T_CLAIM_DOCUMENT       1:N         ← dokumen per peserta
```

Seluruh FK **`ON DELETE CASCADE`** `[keputusan work owner]`.

⚠️ **Koreksi terhadap skema kandidat.** Kandidat menggantungkan adjustment pada **header klaim**.
`[terverifikasi]` `Claim Life/Activity/SavePesertaClaim.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEPESERTACLAIM` / `RULE-OBJ-ACTIVITY`) menulis delapan field
ke path `…PremiumListSummary.PremiumListDetail(<LAST>).AdjustmentList(<LAST>).<field>` —
**adjustment bersarang di bawah peserta**. Menggantungkannya pada header akan **menghapus
informasi peserta mana** yang dimiliki tiap putaran keputusan, dan itu mematahkan mesin status
(**ADR-0011**) yang menjadikan baris adjustment sebagai unit.

### 2. `PremiumListSummary` dilipat ke header

`[terverifikasi]` Sensus: `PremiumListSummary.CLAIM_NO` 40×, `.PL_NUMBER` 34×, `.RISLIPRNM` 12×,
`.BusinessName` 2× — dan **tidak pernah ber-subscript** di seluruh modul, sehingga **1:1** dengan
klaim.

**Keputusan:** keempat field itu menjadi kolom **`T_CLAIMDATA`**, bersama. `PL_NUMBER` **tidak**
ditempatkan di tabel polis: sumbernya path bersarang `PolicyDataLife.PremiumListSummary.PL_NUMBER`,
bukan `PolicyDataLife.PL_NUMBER` — ia milik summary, bukan milik polis.

### 3. Kolom `T_CLAIMDATA`

Dari `ClaimData` skalar `[terverifikasi]`: `pyID`, `StsKatastrofe`, `KatastrofeNote`, `IsKPR`,
`STNCClaim`, `ACCEPTEDNO`, `STS_REJECT`.
Dari `PremiumListSummary` (§2): `CLAIM_NO`, `PL_NUMBER`, `RISLIPRNM`, `BusinessName`.
Dari tabel flat warisan `[terverifikasi]`: `CASEID` (kunci kasus Pega — jejak ke data lama),
`CLAIM_RETRO`.
Audit: operator pembuat, nama operator pembuat, waktu perubahan terakhir.

⚠️ **Aturan pemilihan kolom diperlebar.** Dokumen kandidat menulis *"kolom = property `ClaimData`
yang dipakai"*, tetapi kolom audit dan `CASEID` **bukan** property `ClaimData` — asalnya work object
dan tabel warisan. Aturan sebenarnya: **kolom = nilai yang dipakai dan perlu bertahan**, dari mana
pun asalnya.

`[terverifikasi]` **Tidak ada satu pun kolom kandidat yang mati.** Pemeriksaan per-field atas
`StsKatastrofe`, `KatastrofeNote`, dan `IsKPR` di `Section/InputRegisterClaimLife.xml`,
`Section/InputOSClaimLife.xml`, `Section/InputAkseptasiClaimLife.xml`, dan
`Section/MedicalCheckClaimLife.xml` menemukan **tidak ada ekspresi visibilitas** yang menggerbangi
ketiganya.

### 4. Kolom `T_CLAIM_POLICY` — snapshot polis

`[keputusan work owner]` Isinya **snapshot** yang diambil dari data polis life / PremiumList NB/EDM
**saat klaim dibuat** — disalin, bukan dibaca live. Perubahan polis kemudian **tidak** mengubah
klaim berjalan.

Dari `PolicyDataLife` `[terverifikasi]`: `BusinessCode`, `BusinessID`, `BusinessName`, `CedingCo`,
`CedingCoName`, `Description`, `JenisAsuransi`, `MarketingCode`, `MarketingName`, `MOID`, `NoOffer`,
`PolicyHolder`, `PolicyHolderName`, `ProRateType`, `RetroID`, `RetroName`, `SecurityReinsurer`,
`SecurityReinsurerID`, `SobName`, `SourceOfBusiness`, `Type`, `TypeCeding`, `WPC`.

**Delapan tambahan dari audit** `[terverifikasi]`, tidak ada di kandidat:

| Field | Bukti |
| --- | --- |
| `TanggalRespon`, `TanggalRealisasi`, `TanggalKonfirmasiBalik` | tampil di **empat** section aktif: `InputRegisterClaimLife`, `InputOSClaimLife`, `InputAkseptasiClaimLife`, `MedicalCheckClaimLife` — tanpa gerbang visibilitas |
| `ProductNameID`, `ProductName` | `pyWorkPage.PolicyDataLife.ProductNameID`; `PRODUCTNAMEID`/`PRODUCTNAME` di INSERT warisan |
| `TeamGroup` | `pyWorkPage.PolicyDataLife.TeamGroup` |
| `DateReceived` | `pyWorkPage.PolicyDataLife.DateReceived` |
| `BranchName`, `BranchCode` | `pyWorkPage.PolicyDataLife.BranchName` / `.BranchCode` — **sama dengan keluaran** `PEGA_JSON_KLAIM_PNC` (§10) |

⚠️ `Type` di sini adalah **sumber otoritatif** — `spec.md` §5 sudah menetapkan bahwa `pyWorkPage.Type`
dan `PolicyDataLife.Type` menjadi **satu field**. Perubahannya masuk jejak audit (**ADR-0007**)
karena ia menggerbangi wewenang kirim Komite (**ADR-0012**).

### 5. Kolom `T_CLAIM_MARKETING`

Dari `MarketingData` `[terverifikasi]`: `ID`, `ClientID`, `ClientName`, `BranchDetailID`,
`BranchDetailName`, `TeamGroup`. Tidak berubah dari kandidat.

⚠️ `TeamGroup` muncul **di dua tempat** — `MarketingData.TeamGroup` dan
`PolicyDataLife.TeamGroup`. Keduanya dibawa apa adanya; **jangan digabung** sebelum terbukti
nilainya selalu sama.

### 6. Kolom `T_CLAIM_PREMIUM_LIST_DETAIL` — satu baris per peserta yang diklaim

Dari kandidat `[terverifikasi]`: `PL_NUMBER`, `POLICY_NO`, `POLICY_HOLDER`, `CERTIFICATE_NO`,
`NAME_OF_INSURED`, `DOB`, `AGE`, `SEX`, `PLAN`, `DISEASE`, `ICD_CODE`, `DESCRIPTION`, `NOTES`,
`KETERANGAN`, `DATE_OF_LOSS`, `RECEIVED_DATE`, `BEGIN_DATE`, `EFFECTIVE_DATE`, `EXPIRED_DATE`,
`LAPSE_DATE`, `GROSS_VALUATION_BEGIN_DATE`, `GROSS_VALUATION_EXPIRED_DATE`,
`RETROCESSION_VALUATION_BEGIN_DATE`, `RETROCESSION_VALUATION_EXPIRED_DATE`, `CURRENCY`,
`SUM_INSURED`, `SUM_REASURED`, `GROSS_PREMIUM`, `NET_PREMIUM`, `CLAIM_AMOUNT`, `EM_PERCENT`,
`SHARE_NUSANTARA_RE`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `WPC`, `STNCTreaty`.

**Tujuh tambahan dari audit** `[terverifikasi]`:

| Field | Bukti | Mengapa penting |
| --- | --- | --- |
| **`IsCheck`** | 3 rujukan | Penanda **peserta dipilih untuk diklaim**. Keputusan work owner *"hanya peserta yang mau diklaim"* tidak punya penyimpan tanpa ini |
| `STATUS` | `Activity/SaveAdjustment_Act.xml` → `TempInputDetail.CARI14 = .STATUS`; kolom `STATUS` di INSERT warisan | |
| `RECOMMENDATION` | 1 rujukan | rekomendasi per peserta |
| `CEDING_RETENTION` | 1 rujukan pada peserta; kolom `CEDING_RETENTION` di INSERT warisan | kandidat hanya menaruhnya di adjustment |
| `CONFIRMATION_DATE` | `…PremiumListDetail(N).CONFIRMATION_DATE` | ⚠️ **per peserta**, bukan per klaim |
| `COMPLETE_DATE` | `…PremiumListDetail(N).COMPLETE_DATE` | ⚠️ idem |
| `CLAIM_RECEIVED_DATE` | `…PremiumListDetail(N).CLAIM_RECEIVED_DATE` | ⚠️ idem |

⚠️ **`STS_REJECT` peserta TIDAK menjadi kolom.** `spec.md` §3 sudah menetapkan bahwa
`PremiumListDetail.STS_REJECT` adalah **cerminan** baris adjustment terakhir, **bukan** unit
keputusan tersendiri — di sistem baru ia **turunan**, bukan kolom yang ditulis mandiri. Keputusan
itu **tetap berlaku** dan tidak diubah di sini, meskipun korpus menulis kedua tingkat berbarengan.

⚠️ `idx` / `IndexPremium` adalah **variabel loop**, bukan kolom `[terverifikasi]`.

### 7. Kolom `T_CLAIM_ADJUSTMENT` — putaran keputusan per peserta

Dari kandidat `[terverifikasi]`: `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`, `SUM_REASURED`,
`SUM_INSURED`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `CURRENCYID`, `CURRENCY`.

**Empat tambahan dari audit** `[terverifikasi]` — seluruhnya di-set ke `.AdjustmentList(<LAST>)`
di `SavePesertaClaim.xml` dan `SaveInsuredClaim_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEINSUREDCLAIM_ACT`):

`CLAIM_AMOUNT`, **`STS_REJECT`**, `ACCEPTEDNO`, `ACCEPTATION_DATE`.

⚠️ **`STS_REJECT` di sini adalah status baris** — unit mesin status (**ADR-0011**). Nilai `0`
Outstanding, `1` **diaksep**, `2` ditolak. Baris bersifat **terminal**; revisi dilakukan dengan
**menambah baris baru**. Seluruh aturan itu ada di `spec.md` §3 dan **tidak diubah**.

### 8. `T_CLAIM_DOCUMENT` — tabel baru, dokumen per peserta

`[terverifikasi]` `Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) beriterasi atas
`.DocumentList` per peserta dan **menolak simpan** dengan pesan
`"The document hasn't been uploaded person number "+<nomor>` serta
`"Documents are incomplete, please complete the documents"`.

Jadi dokumen **bersifat per peserta** dan menjadi **gerbang simpan ke Outstanding** — bukan
pelengkap. Ia **tidak ada** di skema kandidat.

⚠️ `[terbuka]` **Bentuk penyimpanannya belum ditetapkan** — lihat §Pertanyaan terbuka OQ-A.

### 9. Tipe dan nullability

`[keputusan work owner]`, **ADR-0003**:

| Kelompok | Tipe |
| --- | --- |
| Uang: `SUM_INSURED`, `SUM_REASURED`, `GROSS_PREMIUM`, `NET_PREMIUM`, `CLAIM_AMOUNT`, `CEDING_RETENTION` | **desimal presisi arbitrer**, tidak pernah `float` |
| Persen/share: `SHARE_NUSANTARA_RE`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `EM_PERCENT` | desimal |
| Tanggal: seluruh `*_DATE`, `DOB`, `Tanggal*` | **`DATE`** |
| Usia | bilangan bulat |
| Sisanya | teks |

**Seluruh kolom nullable**; wajib-isi ditegakkan **di Go**. Identitas via **sequence**
(**ADR-0006**) — pengguna tidak pernah mengetiknya.

### 10. Yang dibaca, bukan dimiliki

- **Data polis life / PremiumList NB/EDM** → sumber snapshot `T_CLAIM_POLICY` dan kandidat peserta.
- **Master** retro, security reinsurer, currency → dibaca saat pengisian.
- **Master diagnosa** `[terverifikasi]` class `ASM-FW-GISFW-DATA-DIAGNOSELIFE` (132 rujukan;
  `Harness/Diagnose_Harness.xml`, `Section/Diagnose_Section.xml`,
  `Activity/SearchDiagnose_act.xml`, `Activity/SetDisease.xml`) adalah **master pencarian**, bukan
  sub-object klaim: hasilnya mendarat sebagai **nilai tunggal** `DISEASE` / `ICD_CODE` pada peserta.
  **Tidak** menjadi tabel.

### 11. Aturan peserta hidup tetap berlaku — dan kini punya penyimpan

`spec.md` §16 menetapkan: **peserta yang dibatalkan (EDM Batal) atau di-soft-delete tidak boleh
muncul di Claim Life**. Penandanya `EDMSTATUS` di `M_LIFE_PREMIUM_DETAIL`.

`[terverifikasi]` Sensus ulang: **nol** kemunculan `EDMSTATUS` di seluruh modul `Claim Life/`.
Aturan itu **tidak rusak** oleh pindah tabel — ia berlaku pada **pembacaan kandidat peserta** dari
tabel premium list, bukan pada penyimpanan klaim.

⚠️ Justru sebaliknya: karena klaim kini menyimpan **snapshot** peserta, `IsCheck` (§6) menjadi
**wajib** — ia yang merekam siapa yang dipilih pada saat pembacaan itu.

### 12. Hilir membaca tabel, bukan JSON

`[keputusan work owner]` Arasapas/produksi **`SELECT` langsung** dari tabel klaim baru.
`convertJsonNusareToProductionClaimLife` dan seluruh serialisasi JSON keluar **tidak
dimigrasikan**. Kontrak efek keluar (**ADR-0008**, **ADR-0015**) dan flag lingkungan
(**ADR-0005**) di `spec.md` §12 **tidak berubah** — yang berubah hanya **bentuk** yang dibaca.

### 13. Migrasi

**ADR-0009** migrasi penuh; **ADR-0011** seluruh **baris** adjustment ikut pindah, bukan hanya
keadaan terakhir.

Sumber: `OS_AKSEPTASI_KLAIM_LIFE` (51 kolom, satu baris per peserta) + dokumen JSON klaim.
Bentuk tujuan: enam tabel di atas — sehingga migrasi harus **menormalkan**: 14 atribut polis yang
hari ini diulang di setiap baris peserta menjadi **satu baris `T_CLAIM_POLICY`** per klaim.

⚠️ `[terverifikasi]` **Dua nilai di-hardcode** di INSERT warisan, bukan di-bind:
`ACCEPTATION_DATE = SYSDATE` dan `STS_REJECT = '0'`. Artinya data lama **selalu** berstempel tanggal
akseptasi pada saat insert dan **selalu** ber-`STS_REJECT` nol pada baris pertama. Migrasi
**melaporkan** ini, tidak menafsirkannya. Perlakuan sistem baru → §Pertanyaan terbuka OQ-C.

---

## Acceptance Criteria

### Bentuk penyimpanan

1. [ ] ⚠️ **Tidak ada blob JSON** sebagai penyimpan atribut klaim. Test yang menemukan kolom JSON
   menyimpan isi klaim **gagal**. *(penyimpangan sadar 1)*
2. [ ] ⚠️ Klaim disimpan di **tabel milik konteks klaim sendiri**; `OS_AKSEPTASI_KLAIM_LIFE` dan
   `M_LIFE_PREMIUM_DETAIL` **tidak ditulis** untuk klaim baru. *(penyimpangan sadar 2)*
3. [ ] ⚠️ **`T_CLAIM_ADJUSTMENT` menggantung pada peserta**, bukan pada header klaim. Setiap baris
   adjustment dapat ditelusuri ke **satu** peserta. Test yang menemukan FK adjustment menunjuk
   header **gagal**. *(§1; penyimpangan sadar 4)*
4. [ ] ⚠️ Atribut polis tersimpan **sekali per klaim**. Test yang menemukan atribut polis berulang
   per peserta **gagal**. *(§4; penyimpangan sadar 3)*
5. [ ] Keempat field `PremiumListSummary` — nomor klaim, nomor premium list, RISLIPRNM, nama
   bisnis — berada di **header**, tidak terpecah ke tabel lain. *(§2)*
6. [ ] Identitas setiap baris berasal dari **sequence**; pengguna tidak pernah mengetiknya.
   *(**ADR-0006**)*

### Kelengkapan kolom

7. [ ] ⚠️ Peserta menyimpan **penanda dipilih-untuk-diklaim**. Tanpa kolom ini fitur "hanya peserta
   yang diklaim" **tidak dapat dinyatakan selesai**. *(§6; `IsCheck`)*
8. [ ] ⚠️ Tanggal **diterima**, **konfirmasi**, dan **penyelesaian** tersimpan **per peserta**.
   *(§6)*
9. [ ] Peserta menyimpan **status** dan **rekomendasi**. *(§6)*
10. [ ] Peserta menyimpan **retensi ceding**. *(§6)*
11. [ ] ⚠️ **`STS_REJECT` peserta TIDAK menjadi kolom tersimpan** — ia **turunan** dari baris
    adjustment terakhir. Test yang menemukannya ditulis mandiri **gagal**. *(§6; `spec.md` §3)*
12. [ ] Baris adjustment menyimpan **claim amount**, **status baris**, **nomor akseptasi**, dan
    **tanggal akseptasi**. *(§7)*
13. [ ] Snapshot polis menyimpan **tanggal respon, tanggal realisasi, tanggal konfirmasi balik**.
    *(§4)*
14. [ ] Snapshot polis menyimpan **produk** (id + nama), **team group**, **tanggal diterima**, dan
    **cabang** (kode + nama). *(§4)*
15. [ ] Header menyimpan **kunci kasus lama** dan **claim retro**. *(§3)*
16. [ ] Tidak ada kolom yang hanya hidup di balik **visibilitas mati**. *(§3)*

### Atomisitas dan kegagalan

17. [ ] Satu klaim beserta seluruh anaknya ditulis dalam **satu transaksi**; kegagalan di mana pun
    **membatalkan seluruhnya**. *(User story 3)*
18. [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan** dengan pesan yang menyebut
    apa yang gagal. *(User story 4; **ADR-0015**)*
19. [ ] Setiap penyimpanan mencatat **jejak audit**. *(**ADR-0007**)*

### Dokumen

20. [ ] ⚠️ Dokumen tersimpan **per peserta**. *(§8; penyimpangan sadar 5)*
21. [ ] ⚠️ Menyimpan ke Outstanding **ditolak** bila ada peserta yang dokumennya belum lengkap,
    dengan pesan yang **menyebut peserta mana**. *(§8; User story 19)*

### Menghapus

22. [ ] ⚠️ Menghapus klaim **mengkaskade** ke polis, marketing, peserta, adjustment, dan dokumen.
    *(§1; penyimpangan sadar 6)*
23. [ ] ⚠️ Penghapusan didahului **popup konfirmasi Ya/Batal** yang menyebut **jumlah baris tiap
    jenis**. **Batal** tidak mengubah apa pun. *(penyimpangan sadar 6)*
24. [ ] Kaskade berjalan dalam **satu transaksi**.

### Tipe

25. [ ] ⚠️ Seluruh nilai uang dan share bertipe **desimal presisi arbitrer**; **tidak** melewati
    `float` dan **tidak** disimpan sebagai teks. *(§9; **ADR-0003**; penyimpangan sadar 7)*
26. [ ] ⚠️ Seluruh tanggal bertipe **`DATE`**. *(§9; penyimpangan sadar 7)*
27. [ ] Seluruh kolom **nullable**; wajib-isi ditegakkan **di Go**. *(§9)*

### Hilir

28. [ ] ⚠️ Sistem hilir membaca **langsung dari tabel**; **tidak ada** payload JSON yang dikirim.
    *(§12; penyimpangan sadar 1)*
29. [ ] ⚠️ **Tidak ada komponen bernama "Json"** yang bukan JSON di sistem baru. *(penyimpangan
    sadar 8)*

### Migrasi

30. [ ] Seluruh klaim lama pindah **tanpa kehilangan satu nilai pun**, termasuk **seluruh baris
    adjustment** — bukan hanya yang terakhir. *(**ADR-0009**, **ADR-0011**)*
31. [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
    tepat**, bukan dengan toleransi. *(**ADR-0003**)*
32. [ ] Tanggal yang berupa teks menjadi `DATE` **tanpa pergeseran zona waktu**; yang **tidak dapat
    diurai dilaporkan**, bukan didiamkan. *(§13)*
33. [ ] ⚠️ Atribut polis yang hari ini **berulang di setiap baris peserta** menjadi **satu baris**
    per klaim; **perbedaan nilai antar baris dalam satu klaim dilaporkan**, tidak diam-diam diambil
    salah satu. *(§13; penyimpangan sadar 3)*
34. [ ] ⚠️ Migrasi **melaporkan** bahwa data lama ber-`ACCEPTATION_DATE` = waktu insert dan
    `STS_REJECT` = `0` karena **di-hardcode**, bukan karena nilainya sebenarnya. *(§13)*
35. [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

---

## Testing Decisions

**Apa yang membuat test bagus di sini.** Test menguji **perilaku yang terlihat dari luar** — apa
yang tersimpan, apa yang ditolak, apa yang dibaca hilir — bukan bentuk internal.

**Seam: API HTTP Claim — Life** — seam **yang sama** dengan `spec.md` §Testing. **Tidak ada seam
baru.** Revisi ini mengubah bentuk penyimpanan, bukan titik pengamatan.

**Oracle tidak pernah dipalsukan.** Yang diuji di sini — atomisitas lintas enam tabel, kaskade
hapus tiga tingkat, presisi desimal, konversi tanggal, dan normalisasi atribut polis saat migrasi —
**hanya berperilaku benar pada basis data sungguhan**.

**Yang wajib punya test tersendiri**, karena di sinilah bug akan muncul:

- **Induk adjustment** — dua peserta, masing-masing dua putaran adjustment; pastikan keempat baris
  dapat ditelusuri ke peserta yang benar. *(AC 3)*
- **Atomisitas** — suntikkan kegagalan pada baris dokumen ke-N; pastikan **tidak ada** klaim
  tersimpan. *(AC 17)*
- **Kaskade tiga tingkat** — hapus klaim; pastikan adjustment dan dokumen **cucu** ikut hilang.
  *(AC 22)*
- **Gerbang dokumen** — peserta tanpa dokumen menahan simpan ke Outstanding, dan pesannya menyebut
  peserta itu. *(AC 21)*
- **`STS_REJECT` peserta sebagai turunan** — ubah baris adjustment terakhir, pastikan nilai peserta
  ikut tanpa pernah ditulis mandiri. *(AC 11)*
- **Normalisasi migrasi** — data lama dengan atribut polis **berbeda** antar baris peserta dalam
  satu klaim harus **dilaporkan**, bukan diam-diam diambil salah satu. *(AC 33)*
- **Presisi uang** — nilai berangka banyak di belakang koma melewati simpan-baca tanpa berubah.
  *(AC 25)*

**Prior art**: pola kaskade+popup mengikuti `.scratch/master-contract-retro-life/` (tiket 09);
pola migrasi JSON→relasional mengikuti `.scratch/master-product-name-life/` (tiket 01); pola satu
transaksi lintas banyak tabel mengikuti `.scratch/treaty-contract-out/` (tiket 09).

---

## Out of Scope

**Tidak disentuh revisi ini** — tetap berlaku sebagaimana di `spec.md`: mesin status per baris
(§3, **ADR-0011**), peran dan wewenang (§4, **ADR-0002**, **ADR-0012**), `Type` sebagai satu field
(§5), validasi Date of Loss (§6), jenis klaim dari kode produk (§7), kontrak Komite (§8,
**ADR-0014**), penomoran (§10), environment (§11), efek keluar asinkron (§12, **ADR-0008**), jejak
audit (§13, **ADR-0007**).

**Tidak menjadi tabel** `[terverifikasi]`:

| Yang dikecualikan | Alasan |
| --- | --- |
| **Master diagnosa** (`ASM-FW-GISFW-DATA-DIAGNOSELIFE`) | master pencarian; hasilnya nilai tunggal pada peserta (§10) |
| `ASM-FW-GCNMFW-Data-Object`, `-ObjectCoverage`, `-ObjectItem` | hanya muncul di berkas **Komite** (`CreateKMTLife_Act`, `GetListKomiteLife`, `Committe_Life`, `ClaimComite`, `RejectOSClaimLife_Sec`, `SendEmailKlaimLF`) dan **nol path data** — milik konteks Komite, bukan klaim |
| `ASM-FW-GCNMFW-Data-Comitee` | konteks **Komite Claim Life**, spec tersendiri |

**Tidak dimigrasikan**: `convertJsonNusareToProductionClaimLife` dan seluruh serialisasi JSON keluar
(§12); rule yang sudah tercatat di `spec.md` §15.

**Di luar konteks**: `Claim Life/RDBList/GetJsonProductLife.xml`
(`ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!GETJSONPRODUCTLIFE`) membaca JSON milik **Master Product
Name Life**, yang sedang dimigrasikan ke relasional di konteksnya sendiri. Penyelarasannya
pekerjaan lintas konteks, bukan bagian revisi ini.

---

## Pertanyaan terbuka — **spec ini `needs-info` karena tiga hal berikut**

| # | Pertanyaan | Pemilik | Yang menunggu |
| --- | --- | --- | --- |
| **OQ-A** | **Dokumen per peserta**: tabel sendiri, atau lampiran bawaan? `[terverifikasi]` `Section/AttachDocScreenLife.xml` berclass **`DATA-WORKATTACH-FILE`** — mekanisme lampiran bawaan Pega — sedangkan `.DocumentList` adalah daftar periksa per peserta. Keduanya mungkin dua hal berbeda. Kolomnya **tidak dapat ditetapkan** tanpa jawaban ini. | work owner | §8, AC 20–21 |
| **OQ-B** | **Tanggal per peserta atau per klaim?** `CONFIRMATION_DATE`, `COMPLETE_DATE`, `CLAIM_RECEIVED_DATE` ada di **kedua** tingkat di Pega. Spec ini menempatkannya **di peserta** (§6) karena di situlah korpus menulisnya. Bila header juga perlu, itu kolom tambahan — bukan pemindahan. | work owner | §3, §6, AC 8 |
| **OQ-C** | **`ACCEPTATION_DATE = SYSDATE` dan `STS_REJECT = '0'` yang di-hardcode**: ditiru, atau diganti nilai sebenarnya? | work owner | §13, AC 34 |

**Terbuka tetapi tidak memblokir:**

| OQ | Pertanyaan | Pemilik |
| --- | --- | --- |
| **OQ-D** | Body `POOLDATA.PEGA_JSON_KLAIM_PNC`. `[terverifikasi]` `Claim Life/RDBList/InsertJsonClaimLifeGCNM.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONCLAIMLIFEGCNM`) memanggilnya dengan **4 masukan + `BRANCH_NAME`/`BRANCH_CODE` keluar**, lalu **`COMMIT`**. Adanya `COMMIT` menunjukkan ia **menulis**, bukan sekadar lookup — dugaan "lookup cabang" **belum terbukti**. | DBA |
| **OQ-E** | Apakah `MarketingData.TeamGroup` dan `PolicyDataLife.TeamGroup` selalu bernilai sama (§5) | work owner |
| **OQ-001** (sisa) | DDL tabel **existing** yang dibaca sebagai sumber snapshot — tipe, presisi, skala | DBA |

**Apa yang membuat spec ini `ready-for-agent`:** jawaban atas **OQ-A, OQ-B, dan OQ-C**. OQ-A
menentukan satu tabel penuh; OQ-B dan OQ-C masing-masing satu AC. Selebihnya sudah mengendap.

---

## Further Notes

### Urutan irisan yang disarankan untuk `/to-tickets`

Revisi ini **menyentuh tiket yang sudah terbit**, bukan hanya menambah:

| Tiket `spec.md` | Dampak |
| --- | --- |
| `03-baris-adjustmentlist-dan-save-ke-outstanding` | ⚠️ **induk adjustment berubah** — dari klaim ke peserta; + gerbang dokumen |
| `02-register-klaim-dan-penomoran` | + snapshot polis & marketing; + `IsCheck` |
| `04-mesin-status-per-baris` | `STS_REJECT` peserta jadi **turunan**, bukan kolom |
| `13-migrasi-data-penuh` | ⚠️ **berubah total** — JSON→relasional + normalisasi atribut polis |
| `12-efek-keluar-asinkron-dan-flag-lingkungan` | payload JSON → `SELECT` dari tabel |

Irisan **baru** yang diperlukan: **skema + migrasi sebagai PREFACTOR** (pola yang sama dengan
Master Product Name Life tiket 01 dan Treaty Contract Out tiket 01) — tidak ada irisan lain yang
berdiri sebelum bentuk barunya ada. Dan satu irisan **dokumen per peserta**, setelah OQ-A terjawab.

### Motif yang berulang

**"Baca kodenya, jangan namanya."** Konteks ini menyumbang **sepasang** contoh yang saling
membalik: `InsertJsonKlaimLife_sql` bernama "Json" tetapi **INSERT flat 51 kolom**; sementara
`GetJsonProductLife` bernama "Json" dan **memang** JSON. Nama tidak memberi tahu apa-apa di kedua
arah.

Dan satu motif baru: **"periksa induk daftarnya, jangan namanya."** `AdjustmentList` terdengar
seperti daftar milik klaim — dan itulah yang diasumsikan skema kandidat. Path lengkapnya
membuktikan ia milik **peserta**. Satu kesalahan induk akan menghapus informasi yang tidak dapat
dipulihkan setelah migrasi berjalan.
