# Discovery Renewal — Putaran 2: bedah 20 berkas delta dan jawaban R-1…R-6

> **Sumber:** `D:\migrasi\RNM\RNW Fac In\`, diurai dengan `[System.Xml.XmlDocument]`, langkah dihitung
> sampai kedalaman penuh. `Endorsment Fac In\` dan korpus Treaty **tidak dibaca** (K-030, K-005).
> Label mengikuti `CLAUDE.md` §3.

---

## Ringkasan temuan

### ⛔ T-6 — **Koreksi atas temuan T-3 saya sendiri: `isApproved` BUKAN rujukan menggantung**

Putaran 1 melaporkan `isApproved` sebagai rujukan menggantung karena DecisionTable-nya tidak ada di
RNW sementara dua activity menyebut namanya. **Itu keliru.**

`[terverifikasi]` Kedua rujukan itu menunjuk **properti**, bukan rule keputusan:

| Berkas perujuk | Kelas rule | `pxRuleObjClass` |
| --- | --- | --- |
| `Activity\InsertHistoryAkseptasiPega` | `ASM-FW-GISFW-Data-PolicyTreatyIn` | **`Rule-Obj-Property`** |
| `Activity\SaveViewSuggest` | `ASM-FW-GISFW-Data-SuggestList` | **`Rule-Obj-Property`** |

Pemakaiannya memperkuat: `pyWorkPage.PolicyTreatyIn.IsApproved=="1"` — perbandingan nilai properti,
dan `@if(.IsApproved="1","Accept", @if(.IsApproved="0","Reject",""))` — pembacaan nilai properti.
**Dua properti berbeda pada dua kelas berbeda, keduanya kebetulan bernama sama dengan sebuah
DecisionTable.**

⚠️ **Ini kekeliruan "nama bukan bukti" yang saya buat sendiri**, dan bentuknya sama persis dengan
jebakan yang sudah tiga kali muncul di proyek ini (`IsFacout` 11-vs-10, `UpdateErrorNoteJsonPolis`
awalan, `TotalSumInsured` filter kata kunci). Pencarian nama **tidak pernah** cukup; tipe rule wajib
diperiksa.

📌 Efek samping bernilai: `IsApproved` sebagai properti memakai pengkodean **`"1"` = Accept, `"0"` =
Reject** — **berbeda** dari domain `ProposalAcceptStatus` {1,2,3,4,7,9} yang dikunci K-021. Dua ruang
nilai berbeda yang tidak boleh disatukan.

### ✅ T-7 — RNW memuat **implementasi konversi produksi yang ekspor NB tidak punya**

`[terverifikasi]` Perbandingan langsung:

| Activity | Folder | Kelas | Langkah | Isi |
| --- | --- | --- | ---: | --- |
| `serviceInsertArasapas_act` | NB | `Data-PolicyTreatyIn` | **1** | hanya `Call serviceInsertArasapas_act` — stub pendelegasi |
| **`serviceInsertArasapasRNW_act`** | **RNW** | **`ASM-FW-GISFW-Work`** | **10** | implementasi penuh |

Kesepuluh langkahnya: `Call serviceInsertArasapas_act` → `Property-Set` → `RDB-List` →
`Property-Set` → **`Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService`** → **`Connect-REST`** →
`Property-Set` → `Call InsertLogServiceProd` → `Page-Remove` → `RDB-List`.

**Mengapa ini penting.** Discovery NB mencatat bahwa versi kelas `Work` dari activity konversi
produksi **tidak ada di korpus** — yang ada hanya stub kelas Treaty yang mendelegasikan ke sana.
Folder RNW memuat sebuah implementasi kelas `Work` yang lengkap, dan di dalamnya terlihat polanya:
**endpoint diambil dari `M_LINK_SERVICE` lewat `GetLinkService`, lalu dipanggil dengan `Connect-REST`,
lalu dicatat ke log layanan.**

Ini menguatkan `CLAUDE.md` §4.4 dari korpus: endpoint memang datang dari tabel, bukan literal.

⚠️ `[dugaan]` Apakah activity ini setara dengan versi `Work` yang hilang dari NB, atau varian khas
renewal, **belum dibuktikan** — namanya berakhiran `RNW`. Yang terbukti hanya bahwa **polanya
terbaca di sini**.

### ⚠️ T-8 — Flow renewal **tidak dirujuk berkas mana pun**; ia titik masuk mandiri

`[terverifikasi]` `Flow\InputRenewalFacultativeIn` dicari di seluruh korpus RNW dan **hanya ditemukan
di berkasnya sendiri**. Tidak ada rule lain yang memanggilnya.

`[dugaan]` Ia **starting flow** yang diluncurkan dari konfigurasi work type atau portal Pega, bukan
dipanggil dari rule. Konfigurasi itu tidak ada di korpus.

---

## R-1 · `isApproved` menggantung — **TERJAWAB: tidak menggantung**

**Jawaban `[terverifikasi]`: rujukannya aktif dan sah — tetapi ke PROPERTI bernama `IsApproved`,
bukan ke DecisionTable yang tidak ada.** Tidak ada rujukan menggantung, tidak ada indikasi ekspor
tidak lengkap dari butir ini. Lihat T-6.

**Status: DITUTUP.**

---

## R-2 · Gerbang masuk renewal — **TERJAWAB sebagian**

`[terverifikasi]` `IsOfferFacIn` **tidak dirujuk sama sekali** di korpus RNW, dan **tidak termasuk**
19 rule `When` yang dirujuk flow renewal.

Kesembilan belas rule `When` yang dirujuk `Flow\InputRenewalFacultativeIn`:

```
IsDeclarationPolicy · IsFacRetro · IsGroup · IsGroupCreate · IsInputFacRetro · IsNotPrintRISlip
IsPKSASM · IsTBonding · LetterNoNull · ToDepHeadUW · ToDirMarketing · ToDirTeknik · ToJUW_A
ToKadivFacultative · ToKadivFin · ToKadivTeknik · ToManagerTeknik · ToSeniorUW · ToUW
```

**Seluruhnya predikat routing tangga akseptasi dan percabangan yang sama dengan NB** — tidak ada satu
pun predikat gerbang-masuk di antaranya.

`[pertanyaan terbuka]` **Gerbangnya tidak ada di dalam korpus.** Flow tidak dirujuk rule mana pun
(T-8), sehingga penentu "kapan sebuah kasus masuk alur renewal" berada di **konfigurasi work type /
portal Pega yang tidak ikut terekspor**. Ini bukan celah pembacaan; ini memang di luar korpus rule.

📌 Konsekuensi rancangan: sistem baru **harus menetapkan gerbang masuk renewal secara eksplisit**,
karena tidak ada perilaku terekam yang dapat direproduksi. Itu keputusan bisnis, bukan porting.

---

## R-3 · Apakah `Work-Renewal` antrean tersendiri — **TERJAWAB: bukan, ia kelas pelaporan**

`[terverifikasi]` Dua fakta yang bertolak belakang bila dibaca terpisah, dan sejalan bila dibaca
bersama:

| | Kelas |
| --- | --- |
| `Flow\InputRenewalFacultativeIn` → `pyWorkClass` | **`ASM-FW-GISFW-Work`** — sama dengan NB |
| `ReportDefinition\RenewalList_RD` → `pyClassName` | **`ASM-FW-GISFW-Work-Renewal`** |

**Kasus renewal dibuat pada kelas kerja yang sama dengan NB.** `Work-Renewal` dipakai **hanya** oleh
satu ReportDefinition, dan tidak oleh flow, activity, maupun section mana pun.

`RenewalList_RD` adalah **daftar kandidat renewal** — 11 kolom `[terverifikasi]`:

```
.pyID · .pzInsKey · .pxCreateDateTime · .pxCreateOperator
.Quotation.OldPolicyNo · .Quotation.InsuredName · .Quotation.MarketingName
.OfferFacIn.PolicyData.StartDateTime · .OfferFacIn.PolicyData.EndDateTime
.NBStatus · .NBStatusNew
```

`[dugaan]` `Work-Renewal` berfungsi sebagai **irisan pelaporan** atas kelas kerja yang sama, bukan
antrean terpisah. Belum dibuktikan sepenuhnya — definisi kelasnya sendiri tidak ada di korpus.

---

## R-4 · Penentuan kelompok bisnis — **TERJAWAB**

`[terverifikasi]` `Activity\GetBusinessGroup_Act`, 4 langkah, kedalaman penuh:

| Langkah | Metode | Isi |
| ---: | --- | --- |
| 1 | `Property-Set` | `InputData.CARI2` ← `pyWorkPage.OfferFacIn.QuotationData.BusinessCode` |
| 2 | `RDB-List` | → `RDBList\CariBusinessGID` |
| 3 | `Property-Set` | `ParamBis.CARI10` ← `OutBis.pxResults(1).CARI3` |
| 4 | `Page-Remove` | bersihkan halaman |

`RDBList\CariBusinessGID` `[terverifikasi]`:

```sql
select ID, OLDID, NOTE, GROUPPANEL, BusinessGroupID as CARI3
from business where ID = {InputData.CARI2}
```

Jadi rantainya: **`BusinessCode` → tabel `business` → `BusinessGroupID`**.

⚠️ **Dua kehati-hatian:**

1. Ini menentukan **kelompok bisnis**, **bukan** lini bisnis (COB) yang menggerakkan skala rasio
   K-018. COB tetap ditentukan predikat `IsFire`/`IsPA`/`IsMBU`/… yang **ada** di RNW.
   `[terverifikasi]` DecisionTable `BusinessType_DeT` **dirujuk 0 kali** di RNW — renewal tidak
   memakainya sama sekali.
2. `BusinessCode` adalah enum **98 kode yang artinya masih belum terverifikasi** (pemilik: Product).
   Rantai ini karena itu terbaca mekanismenya, **bukan** artinya.

⛔ `[pertanyaan terbuka]` `pySaveSQL` pada `CariBusinessGID` memanggil **`POOLDATA.PROSESCOPY(...)`**,
sementara `PROSESCOPY` `[terverifikasi]` **tidak ada di basis data** (checklist §A.1, dikonfirmasi
DBA). Dicatat sebagai **rujukan ke prosedur yang tidak ada**; artinya **tidak disimpulkan**.

---

## R-5 · "Ambil polis lama & hitung ulang" — **TERJAWAB sebagian, dan hipotesisnya perlu dikoreksi**

`[terverifikasi]` `OldPolicyNo` disentuh **21 berkas** di korpus RNW — tetapi **17 di antaranya adalah
berkas bersama yang identik byte-per-byte dengan NB**:

| Kelompok | Berkas |
| --- | --- |
| **Bersama dengan NB** (17) | `GetLimitAkseptasi_Act` · `SaveJsonPolicyFacIn_Act` · `SaveEDMToJsonPolicy_Act` · `ViewOldDataEDM_Act` · `CheckSpreadingProtect_ACT` · `CountPaymentEdm_Act` · `ProtectFIREMBUPA_Act` · `SetEndorsementRISlipData` · `SumTSIPremiSpreadedRNM_Act` · `ViewEndorsementRISlipData` · `GetEDMOldIDPEGA` · `GetLastPPNCheckEDM` · `INSERTJSON_JSONPOLIS_FACIN` · `INSERTJSON_JSONPOLISMONITORING_FACIN` · `Periode_IsUW` · `ViewDataPolicy` · `IsErrorSpreading` |
| **Khas RNW** (4) | `PeriodeRenewal` · `PeriodeRenewal_IsUW` · `SFAPortal_Renewal` · `RenewalList_RD` |

**Koreksi terhadap hipotesis.** `OldPolicyNo` **bukan konsep yang ditambahkan renewal** — ia sudah ada
di basis kode bersama, dan dipakai juga oleh jalur endorsement. Yang **khas renewal** hanyalah
**cara menampilkan dan memasukkannya**: dua Section periode, portal, dan daftar kandidat.

`Section\PeriodeRenewal` mengikat **22 properti** `[terverifikasi]`, di antaranya:

```
.QuotationData.OldPolicyNo · .QuotationData.RNWDate · .QuotationData.StatusBusiness
.PolicyData.StartDateTime · .PolicyData.EndDateTime · .PolicyData.OfferingDate
.IsSpecialAcceptance · .IsDeductibleAcceptance · .QuotationData.TypeFacultative
```

Kehadiran `.QuotationData.StatusBusiness` di layar ini sejalan dengan temuan lama bahwa
**`StatusBusiness` adalah pembeda siklus**.

`Section\ChooseInsuredDtl` ternyata **sangat tipis** — hanya 3 properti terikat (`.Name`, `.pyID`,
template). Ia pemilih daftar, bukan penarik data polis lama.

⚠️ `[pertanyaan terbuka]` **Mekanisme "hitung ulang" belum ditemukan.** Tidak ada satu pun dari 20
berkas delta yang berisi perhitungan. Bila renewal memang menghitung ulang, perhitungannya terjadi di
**berkas bersama** — yang berarti dikendalikan `StatusBusiness`, bukan oleh kode khas renewal. Belum
dibuktikan; ini bahan putaran berikutnya.

---

## R-6 · Apakah 176 NB-only wajar tak dipakai — **TERJAWAB dengan pengukuran**

Uji yang dipakai: apakah nama berkas NB-only **muncul** di dalam korpus RNW (269,5 MB dimuat sekali,
pencocokan **batas kata**).

| Hasil | Jumlah |
| --- | ---: |
| **Tidak muncul sama sekali di RNW** — bersih | **132** |
| **Muncul di korpus RNW** — perlu ditelaah | **44** |

**132 berkas bersih** `[terverifikasi]`: tidak dirujuk, tidak disebut. Konsisten dengan "renewal
memang tidak memakai jalur ini". Temanya sejalan temuan putaran 1 — Treaty, Life, B2B/SFA/CRM,
perangkat/portal.

### ⚠️ Keempat puluh empat itu **belum boleh disebut menggantung**

Pelajaran T-6 berlaku langsung di sini: **kemunculan nama bukan bukti rujukan rule.** Dua dari 44
sudah terbukti bukan rujukan menggantung:

| Nama | Kenyataan |
| --- | --- |
| `isApproved` (DecisionTable + When) | **properti** bernama sama — T-6 |
| `IsUWAccepted` (When) | **DecisionTable** bernama sama **ADA** di RNW — putaran 1 |

Daftar 44, per tipe:

| Tipe | Nama |
| --- | --- |
| Section (20) | `addDeductible` · `InputCoverageAneka_FacIn` · `InputDtlCargo` · `InputEndorsement` · `InputEndorsementDtl` · `InputEndorsementDtl_IsUW` · `InputInwardFacultativeSuggest` · `InstallmentList` · `InwardFacIn` · `ObjectDtlAneka_FacIn` · `ObjectOccupation_FacIn` · `OfferFacIn_NusaRe` · `OfferFacIn_NusaRe_IsUW` · `PeriodeEndorsement` · `PeriodePolicy` · `ProtectCurrency` · `TableOfLimit` · `VehicleGrid` · `ViewDtlAdditional` · `WarrantyList` |
| Activity (7) | `CalcultePersentageSpeadingLife_Act` · `GetData_ACT` · **`GetLimitAkseptasi_ActFlow`** · `SaveJsonOfferFacIn_Act` · `SetToJsonOffer_ACT` · `ShowViewCheckListFacOut` · **`SumTreatyCapacity_Act`** |
| FlowAction (6) | `InputDtlDeductibleFire_FacIn` · `InputDtlWarranty` · `InstallmentList` · `ProtectCurrency` · `ReasViewAttachment` · `TableOfLimit` |
| RDBList (4) | `GenerateNoPolicy` · `GetInsuredID` · `GetKurs` · `SearchJobID` |
| When (4) | `isApproved` ✅ · `IsUWAccepted` ✅ · `pyIsMobile` · `StepStatusFail` |
| Flow (1) | `InputQuotation` |
| DecisionTable (1) | `isApproved` ✅ |
| ReportDefinition (1) | `BrowseConveyanceType` |

⛔ **Dua yang paling perlu diperiksa lebih dulu di putaran berikutnya:**
**`GetLimitAkseptasi_ActFlow`** dan **`SumTreatyCapacity_Act`** — keduanya adalah **activity pemanggil
cabang K-006**. Bila rujukannya di RNW ternyata rujukan rule yang sesungguhnya, itu berarti dua cabang
tangga akseptasi menggantung di siklus renewal.

**Status R-6: terjawab untuk 132; 44 sisanya berstatus `[pertanyaan terbuka]` sampai tipe rujukannya
diperiksa satu per satu.**

---

## Peta 20 berkas delta — terpetakan

| Kelompok | Berkas | Temuan |
| --- | --- | --- |
| **Alur masuk** | `Flow\InputRenewalFacultativeIn` (75 shape, 19 rule `When`) · `FlowAction\Renewal_FlowAct` → Section `InputRenewal` · `Renewal_FlowAct_IsUW` | Titik masuk mandiri; gerbangnya di luar korpus (R-2) |
| **Layar input** | `Section\InputRenewal` · `InputRenewal_IsUW` · `InputRenewalDtl` (2,25 MB) · `InputRenewalDtl_IsUW` (2,40 MB) | Dua terbesar; belum dibedah isinya |
| **Periode & polis lama** | `Section\PeriodeRenewal` (22 properti) · `PeriodeRenewal_IsUW` | Mengikat `OldPolicyNo`, `RNWDate`, `StatusBusiness` (R-5) |
| **Portal** | `Section\SFAPortal_Renewal` (kelas `Data-Portal`) | Menyentuh `OldPolicyNo` |
| **Pemilih tertanggung** | `Harness\ChooseInsured` · `Section\ChooseInsuredDtl` (3 properti) · `ReportDefinition\BrowseAccountInsuredEDM` · `Activity\SetDataInsuredEDM_Act` (5 langkah) | Pemilih tipis, kelas SFA |
| **Kelompok bisnis** | `Activity\GetBusinessGroup_Act` (4 langkah) · `RDBList\CariBusinessGID` | `BusinessCode` → `business` → `BusinessGroupID` (R-4) |
| **Sunting marketing** | `FlowAction\EditMarketing` · `Section\EditMarketing` | Belum dibedah |
| **Daftar renewal** | `ReportDefinition\RenewalList_RD` (11 kolom, kelas `Work-Renewal`) | Daftar kandidat (R-3) |
| **Konversi produksi** | `Activity\serviceInsertArasapasRNW_act` (10 langkah) | Implementasi yang NB tidak punya (T-7) |

---

## Pertanyaan terbuka baru

| # | Pertanyaan | Catatan |
| ---: | --- | --- |
| R-7 | Apakah `GetLimitAkseptasi_ActFlow` dan `SumTreatyCapacity_Act` benar-benar dirujuk sebagai rule di RNW, atau tabrakan nama? | Dua cabang K-006; paling berdampak dari 44 |
| R-8 | Gerbang masuk renewal ada di konfigurasi work type/portal yang **tidak terekspor** — sistem baru perlu menetapkannya eksplisit | Keputusan bisnis, bukan porting |
| R-9 | Di mana "hitung ulang" renewal terjadi? Dugaan: di berkas bersama, dikendalikan `StatusBusiness` | Belum dibuktikan |
| R-10 | `CariBusinessGID.pySaveSQL` memanggil `POOLDATA.PROSESCOPY` yang tidak ada di DB | Artinya tidak disimpulkan |
| R-11 | Apakah `serviceInsertArasapasRNW_act` setara versi `Work` yang hilang dari NB, atau varian renewal? | Berdampak pada spec jalur produksi NB |
| R-12 | `Work-Renewal`: irisan pelaporan atau kelas kerja nyata? Definisi kelasnya tidak ada di korpus | |

---

## Belum dikerjakan

- `Section\InputRenewalDtl` dan `InputRenewalDtl_IsUW` (4,65 MB gabungan) — belum dibedah
- `Struktur_InputRenewalFacultativeIn.xlsx` (7,2 MB)
- 42 dari 44 nama pada R-6 — pemeriksaan tipe rujukan per nama
- `FlowAction\EditMarketing` + `Section\EditMarketing`
- `Harness\ChooseInsured` (183 KB) — hanya kelasnya yang dibaca

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
