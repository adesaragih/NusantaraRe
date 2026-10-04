# 46: Form Coverage FIRE — tata letak sesuai Pega, Days radio, Choose / Copy Accumulation

> ⚠️ **Disusun agent atas perintah sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: gambar layar Pega dari
> work owner 03-10-2026 ("choose acumulationnya mana? … days itu dibuat sesuai pega, jangan diubah jadi dropdown.
> rapikan juga kolom2nya, jangan halu"). Frontend sesi 0f; backend diserahkan ke sesi c3.

**What to build:** form satu coverage (`Section\CoverageItem.xml`) ditata seperti Pega, Days berupa radio, Accumulation
Code / Address dengan tombol Choose Accumulation Code (popup pencarian) dan Copy Accumulation.

**Status:** frontend selesai 04-10-2026 (popup lengkap sesuai gambar layar kedua; uji hijau). Backend butir 1 selesai
(migrasi 195). ⛔ Backend butir 2 (pencarian + saran autocomplete) belum — DDL view `ACCUMULATION` kini ada
(`DDL\ACCUMULATION.txt`, 04-10-2026).

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

### Tata letak `Section\CoverageItem.xml`

- Sel 1 "Stacked with labels left": 3 Coverage Basis · 4 Choose Coverage · 5 Coverage (`.OLDID`, `pyReadOnly` tanpa
  syarat) · 6 Zone / 13 4.2 Construction (OTHER, kode coverage tertentu — tahap berikut) · 14 Accumulation Code · 15
  Accumulation Address (`.AccumulationDescription`; keduanya baca-saja) · 16 (Inline) 18 Choose Accumulation Code + 19
  Copy Accumulation · 20 Conditions · 21 Layer (basis 5).
- Sel 22 "Inline grid double": kiri sel 24 = 26 Days · 27 TSI · 28 Indemnity · 29 ‰ Gross Rate · 30 % First Loss · 31
  % Discount · 32 First Loss · 33 ‰ Net Rate · 34 Limit of Liability · 35 % Limit of Liability; kanan sel 36 = 38/39
  % Pro Rate · 40 Indemnity Unit · 41 % Indemnity · 42 % First Scale · 43 % Sub Limit · 44 % Loss Limit · 45 % EML/PML ·
  46 Discount · 47 Gross Premium; sel 48 = 50 View Indemnity Table.
- Syarat tampil (`pyUserData`): ALWAYS dengan syarat sisa (diabaikan) = 28 Indemnity, 40 Indemnity Unit, 44 % Loss Limit,
  33 ‰ Net Rate; OTHER = 30 / 32 / 42 (`.CoverageBasis==2`), 43 (4 / 5), 45 (3), 38 / 39 (ProRatePercent).
- Baca-saja tanpa syarat: 5, 14, 15, 27, 41, 42 (cocok dengan gambar layar: % Indemnity tampil sebagai teks).
- Sel 26 Days = `pxRadioButtons`, `pyOrientation` horizontal; ubah → `CountPremi_ACT` (percent / percent). Daftar =
  PromptList `ASM-FW-GISFW-DATA-COVERAGE!DAY` (`DDL\Day.xml`, 04-10-2026): 365 / 366 / 360 (sama dengan gambar layar).
- Sel 40 Indemnity Unit = `pxDropdown`, `pyHasNoSelection` false; PromptList `ASM-FW-GISFW-DATA-COVERAGE!UNIT`
  (`DDL\Unit.xml`): 0 Day (s) / 1 Month (s) / 2 Year (s). Berkas yang sama memuat LocalList lama (Day(s) / Month(s) /
  Year(s) tanpa kode); yang dipakai = PromptList (gambar layar "Day (s)", data contoh kode "0" / "1").

### Accumulation

- Sel 18 Choose Accumulation Code: `SetFilterAccum_DT` → `SearchAccumAct` → harness `ChooseAccumulation_FacIn`
  (`pyWindowName` "Choose Accumulation") → section `SearchRiskAccumCov.xml`.
- `Activity\SearchAccumAct.xml` langkah 4: `SearchAccumulation.PostalCode = .RiskZipCode` (saringan awal = zip lokasi).
- `SearchRiskAccumCov.xml`: sel 40 Filter (`GetDataAccumulation_act`), sel 41 Clear Column (`SearchAccumulationPre_Act`
  = Page-Remove saringan & hasil); grid sel 90 = RD `SearchRiskAccumulation_RD` (kelas `ASM-FW-GISFW-Int-ACCUMULATION`,
  kolom `.ID` / `.AccumulationName` / `.Note`, judul Accumulation Code / Type / Description, tombol Choose sel 119).
- `ReportDefinition\SearchRiskAccumulation_RD.xml`: `.ID = Param.AccumulationCode`, `.Note` **Contains**
  `Param.AccumulationDescription` (satu-satunya `pyCaseInsensitive` true), `.AccumulationName = Param.Coverage`,
  `.CZone = Param.Czone`, `.Keyword = Param.Keyword`, `.PostalCode = Param.PostalCode`, `.ProvinceID = Param.Province`;
  INNER JOIN `Int-RW` pada `.PostalCode = RW.ZipCode`; DISTINCT; maks. 500.
- `Activity\SetDataAccum_Act.xml` (Choose): zip = `@substring(ID,4,9)` tanpa "-", atau `@substring(ID,3,8)` bila
  `@substring(ID,3,4)!="-"`; beda dengan `RiskLocation.ASMZipCode` → pesan "ZIpCode Harus Sama dengan ZIpCode  Yang di
  Object Item >>>> {zip lokasi} != {zip ID}", nilai tidak ditulis; sama → `AccumulationCode = ID`,
  `AccumulationDescription = Note`.
- `Activity\CopyAccumulationCode_Act.xml` (sel 19, tampil bila `.pxListSubscript==1`): Code / Description coverage sumber
  disalin ke setiap coverage di setiap PropertyItem **lokasi yang sama** (`LocationList(Param.idxlocation)`).
- SQL jalur grid sel 48 (`RDBList\GetAccumulationProvince_SQL.xml` / `GetAccumulationDistrict_SQL.xml`) membaca tabel /
  view **`accumulation`** (kolom `id`, `note`, `accumulationtype`, `ZIPCODE`) dan `rw`.
- `DDL\ACCUMULATION.txt` (04-10-2026): `VIEW POOLDATA.ACCUMULATION` = `SELECT DISTINCT` dari `m_accumulation` JSON
  (ID, Accumulation, AccumulationName dari `m_accumulatedtype`, `UPPER(Note)`, Keyword, ScopeArea, CZone, CZoneID,
  ProvinceName → PROVINCE, ProvinceID, PostalCode → ZIPCODE, AccumulationType, SyariahStatus) `WHERE IsActive IS NULL`.
- `Activity\GetDataAccumulation_act.xml` (Filter): `InputFilter.MapHeight = 2` (grid RD); `InputFilter.PostalCode`
  (dari Area) → `SearchAccumulation.PostalCode`; (City / District / Policy No terisi) DAN Zip Code kosong →
  `MapHeight = 1` (grid sel 48: RDB-List `GetSummaryRiskAccumPolis_Sql` / `GetAccumulationProvince_SQL` /
  `GetAccumulationDistrict_SQL` / `GetAccumulation_SQL`).
- Autocomplete `SearchRiskAccumCov` (RD, `pySetValueOnSelect`): 27 Country `BrowseNation_RD` (→ `.Note`, NationInitial
  → `SearchAccumulation.SyariahStatus`); 28 Province `BrowseProvince2_RD` (NationName = Country; ID → ProvinceID);
  29 Accum. Type `BrowseAccumulatedType_RD`; 32 City `BrowseCityInput_RD` (PROVINCEID; ID → `InputFilter.City`);
  33 District `BrowseDistrictInputC_RD` (CityName; ID → `InputFilter.District`); 34 Area `BrowseRW_RD` (DistrictName;
  ID → `InputFilter.Area`, ZipCode → `InputFilter.PostalCode`); 35 CZone `BrowseCZoneIsNotNull_RD` (`.Code`). 36 Key
  Word = PromptList `ASM-FW-GISFW-INT-ACCUMULATION!KEYWORD` (`DDL\Keyword.xml`, 19 pilihan), pilihan kosong bertulisan
  kosong.
- Label XML "Zip Code / Country / City / Area / Add New"; gambar layar DEV menulis "Zip code / COUNTRY / CITY / AREA /
  Add new" — versi rule berbeda; frontend memakai ejaan korpus.

## Frontend (sesi 0f)

- `components/FormCoverage.tsx`: blok atas + dua kolom berlabel kiri (kelas `nbf-labelkiri`, anak langsung — tanpa
  menyentuh kelas inti); Days radio (`PILIHAN_DAY_COVERAGE`); Indemnity / Indemnity Unit / % Loss Limit selalu tampil;
  % Indemnity, % First Scale, TSI, % Pro Rate = teks rata kanan; Coverage = OLDID; Accumulation kosong = "---".
- `components/PopupAkumulasi.tsx` (`zipDariId`, `keSaring`; seluruh medan saringan), `components/IsianSaran.tsx`
  (autocomplete), `components/TabCoverage.tsx` (`salinAkumulasi`, `zipRisiko`). Indemnity Unit = `<select>` tanpa
  pilihan kosong; coverage baru `unit: '0'`.
- `api.ts`: `CoverageObjek.unit? / accumulationCode? / accumulationDescription?`, `cariAkumulasi`, `SaringAkumulasi`,
  `BarisAkumulasi`, `saranAkumulasi`, `JenisSaranAkumulasi`, `SaranAkumulasi`. `labels.ts`: `FORM_COV` (+ sel 14 / 15 /
  18 / 19 / 40), `PILIHAN_DAY_COVERAGE`, `OPSI_INDEMNITY_UNIT`, `OPSI_KEYWORD`, `AKUMULASI_KOSONG`, `POPUP_AKUMULASI`,
  `TEKS_AKUMULASI`.
- Uji: `labels.test.ts` (kontrol & syarat CoverageItem, label popup, PromptList Unit / Day / Keyword menurut pxInsName),
  `TabCoverage.test.tsx`, `PopupAkumulasi.test.tsx`.

## Kontrak backend (untuk sesi c3)

1. Coverage membawa `unit`, `accumulationCode`, `accumulationDescription` (teks; kosong = NULL) di GET / PUT objek dan
   diteruskan apa adanya oleh hitung-coverage / hitung-net-rate.
2. `GET /api/nbfacin/akumulasi?id=&policyNo=&note=&postalCode=&syariahStatus=&provinceId=&cityId=&districtId=&czone=&keyword=`
   → `{ baris: [{ id, accumulationName, note }] }`, logika `GetDataAccumulation_act`: (cityId / districtId / policyNo
   terisi) DAN postalCode kosong → jalur SQL grid sel 48 (CARI1 → id, CARI2 → accumulationName, CARI3 → note); selain
   itu `SearchRiskAccumulation_RD` atas view `ACCUMULATION` (kosong = filter diabaikan; `note` Contains tanpa peka
   huruf; lainnya `=`; join RW; DISTINCT; maks. 500). Frontend sudah menimpa postalCode dengan zip Area (langkah 4).
3. `GET /api/nbfacin/akumulasi/saran/{jenis}?q=&induk=` → `{ baris: [{ id, label, ekstra? }] }`, `jenis` = `nation`
   (ekstra = NationInitial) · `province` (induk = nama Country) · `accumtype` · `city` (induk = ProvinceID) · `district`
   (induk = nama City) · `area` (induk = nama District; ekstra = ZipCode) · `czone`, masing-masing mengikuti RD sel
   27-35. `label` = nilai yang ditulis Pega ke kotak (`pyValue` autocomplete).

## Keputusan agent (menunggu konfirmasi)

- **K-1** (ralat 04-10-2026) Popup memuat seluruh saringan seperti gambar layar kedua. Tahap berikut: Add New, Summary
  (.Note), Accumulation Risk, Risk Accumulation report, Adjustment.
- **K-2** Hasil tanpa paging 10 baris (grid Pega berpaging 10; hasil maks. 500).
- **K-3** Copy Accumulation menyalin di layar (belum tersimpan sampai Save), sama dengan perubahan lain di tab Coverage.
- **K-4** Perubahan Indemnity Unit (`SetIndemnityRate_ACT`) dan View Indemnity Table (`ViewIndemnity_LA`) = tahap
  berikut.
- **K-6** Indemnity Unit kosong tampil dan dikirim sebagai "0" (dropdown tanpa pilihan kosong); coverage baru `unit "0"`.
- **K-7** Mengetik di autocomplete mengosongkan nilai hasil pilihan sebelumnya (ID / NationInitial / ZipCode); Pega
  membiarkan nilai lama.
- **K-8** Accum. Type hanya tampilan (lihat catatan `[dugaan]` di bawah) — tidak dikirim ke pencarian.
- **K-9** Ejaan label mengikuti korpus ("Zip Code", "Country", "City", "Area"), bukan gambar layar DEV.
- **K-5** Choose tidak menyimpan otomatis (Pega `Obj-Save` di `SetDataAccum_Act` langkah 9); tersimpan lewat Save tab.
- Catatan `[dugaan]`: field Accum. Type (sel 29) terikat `SearchAccumulation.Type`, sedangkan RD menerima
  `SearchAccumulation.AccumulationType` — di Pega saringan ini kemungkinan selalu kosong.

## Pertanyaan untuk work owner

- ~~DDL `ACCUMULATION`, aturan `.Unit` / `.Day` / `.Keyword`~~ — dikirim 04-10-2026 (`DDL\ACCUMULATION.txt`,
  `Unit.xml`, `Day.xml`, `Keyword.xml`).

## Acceptance criteria

- [x] Urutan dan susunan medan sama dengan gambar layar / `CoverageItem.xml`; label di kiri.
- [x] Days radio mendatar 365 / 366 / 360.
- [x] Accumulation Code / Address tampil ("---" bila kosong); Copy Accumulation hanya di coverage pertama.
- [x] Choose menolak accumulation yang zip-nya beda dengan pesan Pega.
- [x] Backend: simpan tiga medan (sesi c3, 03-10-2026; migrasi 195 ditulis, belum dijalankan).
- [x] Popup memuat seluruh saringan (gambar layar kedua); Indemnity Unit / Key Word dari PromptList korpus.
- [ ] Backend: `GET /api/nbfacin/akumulasi` + `…/akumulasi/saran/{jenis}` (sesi c3).

## Backend butir 1 (sesi c3, 03-10-2026)

⛔ **Migrasi 195 WAJIB dijalankan di DEV sebelum backend baru** — urutan … → 194 → **195**. Baca / tulis coverage kini
memakai `T_COVERAGELIST.UNIT` / `ACCUMULATION_CODE` / `ACCUMULATION_DESCRIPTION`; tanpa 195 seluruh tab Object gagal
ORA-00904.

- `coverages[k].unit` / `accumulationCode` / `accumulationDescription` (teks apa adanya, kosong = NULL) di `GET` / `PUT
  …/objek`; diteruskan apa adanya oleh `POST hitung-coverage` / `hitung-net-rate` (tanpa disimpan). 400 ber-jalur bila
  melebihi lebar kolom: `unit` / `accumulationCode` 50, `accumulationDescription` 500 byte. Tanpa pemeriksaan isi kode.
- Migrasi **195**: `ALTER TABLE T_COVERAGELIST ADD` tiga kolom RANCANGAN (skema loader `T_COVERAGELIST`: `UNIT`
  VARCHAR2(50), `ACCUMULATION_CODE` VARCHAR2(50), `ACCUMULATION_DESCRIPTION` VARCHAR2(500)) — tipe / lebar rancangan apa
  adanya, tanpa kolom baru. Sampel `DDL\P-5 *.txt` `[terverifikasi]`: 107 coverage ber-Unit (satu digit) dan 107
  ber-AccumulationCode (bentuk `aaa-99999-999999`); AccumulationDescription paling panjang 149 byte (≤ 500).
- Butir 2 (`GET /api/nbfacin/akumulasi`) **tidak dikerjakan**: sumber kelas `Int-ACCUMULATION` tetap belum terverifikasi.
  `DDL\RDBMASTERACCUMULATION.txt` `[terverifikasi]` adalah PROSEDUR penulis `M_ACCUMULATION` (ID, JSONDATA; ID =
  negara-zip-urutan) — bukan sumber baca RD; `DDL\M_ACCUMULATION.txt` ber-JSON. Kunci JSON tidak ditebak.
