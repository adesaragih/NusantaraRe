# C.2 — Peta Sumber Dropdown + Query Ekstraksi untuk DBA

> Menutup butir C.2 checklist (604 sel dropdown). Korpus hanya memuat **definisi** DataPage (`D_*`),
> yaitu dari tabel/view mana sebuah dropdown mengambil isinya — **bukan nilainya**. Nilai ada di
> baris tabel Oracle, jadi harus diekstrak DBA.
>
> Sumber: `NB FacIn\DataPage\D_*.xml` (`<pyClassName>`) — `[terverifikasi]`.
> Nama tabel diturunkan dari class Pega `ASM-FW-GISFW-Int-XXX` → tabel/view Oracle `XXX`.
> ⚠️ Nama tabel bertanda "(konfirmasi DBA)" belum diverifikasi terhadap DDL; DBA mengoreksi bila beda.
>
> **A sudah dicek dan TIDAK melengkapi:** `D_EnumerationList.xml` hanya pointer ke ReportDefinition
> `DataTableEditorReport` pada tabel class `Data-Enumeration` (param Type/Language/Status) — nilainya
> ada di DB, bukan korpus. Karena itu B ini dibuat.

---

## 1. Dropdown bersumber tabel/view Oracle — dapat diekstrak DBA

| # | DataPage | Dropdown untuk | Tabel/View sumber | Catatan |
| ---: | --- | --- | --- | --- |
| 1 | `D_CITY` | Kota | `CITY` (konfirmasi DBA) | master wilayah |
| 2 | `D_DISTRICT` | Kecamatan/Distrik | `DISTRICT` (konfirmasi DBA) | master wilayah |
| 3 | `D_RWNAME` | RW | `RW` (konfirmasi DBA) | master wilayah |
| 4 | `D_ZipCodeList` | Kode Pos | `V_ZIPCODE` (view) | |
| 5 | `D_BankGroup` | Grup Bank | `LST_BANK_GROUP` (konfirmasi DBA) | |
| 6 | `D_BrowseOccupationFacInFIRE` | Okupasi (Fire) | `OCCUPATION` (konfirmasi DBA) | difilter untuk Fire |
| 7 | `D_Occupation` | Okupasi (umum) | `OCCUPATION` (konfirmasi DBA) | |
| 8 | `D_CoinsList` | Ko-asuransi | `V_COINS` (view) | |
| 9 | `D_ModelList` | Model/Merk kendaraan | `BRANDDETAIL` (konfirmasi DBA) | untuk MBU |
| 10 | `D_BENEFIT_CLAUSE` | Klausa Benefit | `VIEW_BENEFIT_PROPERTY` (view) | |
| 11 | `D_BrowseTableOfLimit` | Table of Limit | `TABLEOFLIMIT` (konfirmasi DBA) | |
| 12 | `D_PlanClause` | Klausa Plan | `M_KLAUSUL_PLAN` (konfirmasi DBA) | |
| 13 | `D_PlanProperty` | Plan Property | `M_TYPE_PROPERTY_PLAN` (konfirmasi DBA) | |
| 14 | `D_J_PlanProperty` | Plan Property (varian J) | `VJ_M_TYPE_PROPERTY_PLAN` (view) | |

## 2. Dropdown enumerasi — satu tabel, banyak jenis

| # | DataPage | Dropdown untuk | Sumber |
| ---: | --- | --- | --- |
| 15 | `D_EnumerationList` | **Banyak enum kecil** (difilter parameter `Type`) | tabel class `ASM-FW-GISFW-Data-Enumeration` via ReportDefinition `DataTableEditorReport` |

⚠️ Satu tabel enum ini kemungkinan mengisi **banyak** dari 604 sel dropdown — tiap jenis dropdown
dibedakan nilai parameter `Type`. DBA perlu mengekstrak **seluruh isi** tabel enum lalu kelompokkan
per `Type`. Kolom yang dipakai (dari parameter DataPage): `Type`, `Language`, `Status`, `Name`,
`Name2`, `Name3`.

## 3. Bukan dropdown (data kerja/summary — abaikan untuk C.2)

`D_Coverage`, `D_CoverageFacOut`, `D_CoverageSummary`, `D_CommentCeding`, `D_FacOfferObjectSection`,
`D_FacOutFromVehicle`, `D_FilterObjectItem`, `D_GolfCoverageSummary`, `D_GolfLocationSummary`,
`D_GrowingTreesCoverageSummary`, `D_GrowingTreesLocationSummary`, `D_LocationSummary`,
`D_PersonListCoverageSummary`, `D_PersonListSummary`, `D_AnekaList` — class `Data-*` (agregat case
yang sedang dikerjakan), bukan daftar pilihan.

---

## 4. Query siap-jalan untuk DBA

⚠️ Ganti nama tabel bertanda "(konfirmasi DBA)" bila ejaan aslinya berbeda. Ambil kolom kode + label
tampilan (biasanya `ID`/`CODE` + `NAME`/`DESCRIPTION`). **Hanya nilai referensi, bukan data pelanggan.**

```sql
-- Enumerasi (paling penting — satu tabel banyak dropdown)
SELECT TYPE, STATUS, NAME, NAME2, NAME3
FROM   <skema>.<TABEL_ENUMERATION>      -- class ASM-FW-GISFW-Data-Enumeration
ORDER  BY TYPE, NAME;

-- Wilayah
SELECT * FROM <skema>.CITY;
SELECT * FROM <skema>.DISTRICT;
SELECT * FROM <skema>.RW;
SELECT * FROM <skema>.V_ZIPCODE;

-- Master bisnis
SELECT * FROM <skema>.LST_BANK_GROUP;
SELECT * FROM <skema>.OCCUPATION;
SELECT * FROM <skema>.V_COINS;
SELECT * FROM <skema>.BRANDDETAIL;
SELECT * FROM <skema>.VIEW_BENEFIT_PROPERTY;
SELECT * FROM <skema>.TABLEOFLIMIT;
SELECT * FROM <skema>.M_KLAUSUL_PLAN;
SELECT * FROM <skema>.M_TYPE_PROPERTY_PLAN;
SELECT * FROM <skema>.VJ_M_TYPE_PROPERTY_PLAN;
```

## 5. Batas yang jujur

- **Yang diberikan di sini:** 14 dropdown + 1 tabel enum = **sumber**-nya (dari korpus, terverifikasi).
- **Yang TIDAK diberikan:** nilai/isi tiap dropdown — ada di DB, diekstrak lewat query di §4.
- **604 vs 15:** 604 = jumlah sel dropdown di layar; 15 sumber ini mengisi sebagian besarnya (satu
  tabel dipakai banyak sel; tabel enum sendiri melayani banyak jenis via parameter `Type`).
- Nama tabel "(konfirmasi DBA)" diturunkan dari nama class Pega, belum dicek ke DDL — DBA wajib
  mengoreksi bila berbeda; itu sendiri sudah jawaban yang berguna.

*Tanpa nama orang, tanpa data pelanggan. Semua sumber terverifikasi dari `NB FacIn\DataPage\`.*
