# 46: Form Coverage FIRE — tata letak sesuai Pega, Days radio, Choose / Copy Accumulation

> ⚠️ **Disusun agent atas perintah sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: gambar layar Pega dari
> work owner 03-10-2026 ("choose acumulationnya mana? … days itu dibuat sesuai pega, jangan diubah jadi dropdown.
> rapikan juga kolom2nya, jangan halu"). Frontend sesi 0f; backend diserahkan ke sesi c3.

**What to build:** form satu coverage (`Section\CoverageItem.xml`) ditata seperti Pega, Days berupa radio, Accumulation
Code / Address dengan tombol Choose Accumulation Code (popup pencarian) dan Copy Accumulation.

**Status:** frontend selesai 03-10-2026 (uji hijau). ⛔ Backend belum: endpoint `GET /api/nbfacin/akumulasi` dan
penyimpanan `unit` / `accumulationCode` / `accumulationDescription` (sampai itu ada, ketiga medan dibuang diam-diam saat
PUT — decoder `simpanObjek` tanpa `DisallowUnknownFields`). ⛔ Sumber data popup menunggu DDL `ACCUMULATION` (lihat
Pertanyaan).

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
- Sel 26 Days = `pxRadioButtons`, `pyOrientation` horizontal; ubah → `CountPremi_ACT` (percent / percent). Daftar dari
  aturan `.Day` — **tidak ada di korpus**; 365 / 366 / 360 diambil dari gambar layar.
- Sel 40 Indemnity Unit = `pxDropdown`, `pyHasNoSelection` false, daftar dari aturan `.Unit` — **tidak ada di korpus**
  (data contoh `DDL\P-5 *.txt` hanya kode "0" / "1").

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

## Frontend (sesi 0f)

- `components/FormCoverage.tsx`: blok atas + dua kolom berlabel kiri (kelas `nbf-labelkiri`, anak langsung — tanpa
  menyentuh kelas inti); Days radio (`PILIHAN_DAY_COVERAGE`); Indemnity / Indemnity Unit / % Loss Limit selalu tampil;
  % Indemnity, % First Scale, TSI, % Pro Rate = teks rata kanan; Coverage = OLDID; Accumulation kosong = "---".
- `components/PopupAkumulasi.tsx` (`zipDariId`), `components/TabCoverage.tsx` (`salinAkumulasi`, `zipRisiko`).
- `api.ts`: `CoverageObjek.unit? / accumulationCode? / accumulationDescription?`, `cariAkumulasi`, `SaringAkumulasi`,
  `BarisAkumulasi`. `labels.ts`: `FORM_COV` (+ sel 14 / 15 / 18 / 19 / 40), `PILIHAN_DAY_COVERAGE`,
  `OPSI_INDEMNITY_UNIT` (kosong), `AKUMULASI_KOSONG`, `POPUP_AKUMULASI`, `TEKS_AKUMULASI`.
- Uji: `labels.test.ts` (kontrol & syarat CoverageItem, label popup), `TabCoverage.test.tsx`, `PopupAkumulasi.test.tsx`.

## Kontrak backend (untuk sesi c3)

1. Coverage membawa `unit`, `accumulationCode`, `accumulationDescription` (teks; kosong = NULL) di GET / PUT objek dan
   diteruskan apa adanya oleh hitung-coverage / hitung-net-rate.
2. `GET /api/nbfacin/akumulasi?id=&note=&postalCode=&czone=` → `{ baris: [{ id, accumulationName, note }] }` mengikuti
   `SearchRiskAccumulation_RD` (kosong = filter diabaikan; `note` Contains tanpa peka huruf; lainnya `=`; join RW;
   DISTINCT; maks. 500). ⛔ Tabel / view sumber kelas `Int-ACCUMULATION` belum terverifikasi — jangan ditebak.

## Keputusan agent (menunggu konfirmasi)

- **K-1** Popup tahap 1 = jalur RD dengan medan teks Accumulation Code, Road, Zip Code, CZone. Tahap berikut: Policy No /
  City / District / Area (jalur SQL grid sel 48), Country / Province / Accum. Type (autocomplete), Key Word (aturan
  `.Keyword` belum ada), Add New, Summary (.Note), Accumulation Risk, Risk Accumulation report, Adjustment.
- **K-2** Hasil tanpa paging 10 baris (grid Pega berpaging 10; hasil maks. 500).
- **K-3** Copy Accumulation menyalin di layar (belum tersimpan sampai Save), sama dengan perubahan lain di tab Coverage.
- **K-4** Indemnity Unit tampil tanpa pilihan sampai aturan `.Unit` dikirim; perubahan Unit (`SetIndemnityRate_ACT`) dan
  View Indemnity Table (`ViewIndemnity_LA`) = tahap berikut.
- **K-5** Choose tidak menyimpan otomatis (Pega `Obj-Save` di `SetDataAccum_Act` langkah 9); tersimpan lewat Save tab.
- Catatan `[dugaan]`: field Accum. Type (sel 29) terikat `SearchAccumulation.Type`, sedangkan RD menerima
  `SearchAccumulation.AccumulationType` — di Pega saringan ini kemungkinan selalu kosong.

## Pertanyaan untuk work owner

- DDL tabel / view `ACCUMULATION` (POOLDATA) — sumber RD dan SQL; tidak ada di `DDL\` (hanya `M_ACCUMULATION` ber-JSON).
- Aturan properti `ASM-FW-GISFW-DATA-COVERAGE!UNIT` dan `!DAY`, serta aturan `.Keyword` (Data-Portal / Int-ACCUMULATION).

## Acceptance criteria

- [x] Urutan dan susunan medan sama dengan gambar layar / `CoverageItem.xml`; label di kiri.
- [x] Days radio mendatar 365 / 366 / 360.
- [x] Accumulation Code / Address tampil ("---" bila kosong); Copy Accumulation hanya di coverage pertama.
- [x] Choose menolak accumulation yang zip-nya beda dengan pesan Pega.
- [x] Backend: simpan tiga medan (sesi c3, 03-10-2026; migrasi 195 ditulis, belum dijalankan).
- [ ] Backend: endpoint `GET /api/nbfacin/akumulasi` (menunggu DDL `ACCUMULATION`).

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
