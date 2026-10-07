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
- ~~**K-2** Hasil tanpa paging 10 baris~~ — diganti 04-10-2026 (work owner: "hasil nya buat page, 10 list per page"): hasil
  dipaging 10 baris di layar, sama dengan grid Pega; backend tetap mengirim maks. 500 sekaligus. Tampilan popup juga
  dirapikan (kepala + zip lokasi, kartu Conditions 4 kolom ringkas, tanda Zip sesuai / berbeda per baris).
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
- [x] Backend: `GET /api/nbfacin/akumulasi` dan saran city / district / area (sesi c3, 04-10-2026).
- [x] Saran nation / province / accumtype / czone (DDL 04-10-2026; sesi c3).
- [ ] ~~Migrasi 196 (view → tabel flat) dijalankan di DEV~~ → 196 DIPINDAH ke modul `masterdata` (760, + CITYINPUT /
  DISTRICTINPUT, + STS_AKTIF) — register butir 106; jalankan 904 → 760 → 761 → 990.
- [ ] CITY / DISTRICT / V_JN_OBJ_ITEM / COVERAGE_FACIN / BRANDDETAIL / VJ_M_TYPE_PROPERTY_PLAN → tabel flat (A179, menunggu).

## Tabel flat dan saran lengkap (sesi c3, 04-10-2026; butir 102–103)

- **Saran ketujuh jenis** (`repository/akumulasi.go` `daftarSaran`) `[terverifikasi]` `Section\SearchRiskAccumCov.xml` +
  RD: `nation` (tabel `NATION`; cari `.Note`, ekstra `.NationInitial` → SyariahStatus; `.ID OR .Note` dikirim kosong →
  dibuang), `province` (`PROVINCE`; cari `.Note`, id `.ID` → ProvinceID, induk `NATIONNAME` = nama negara), `accumtype`
  (`ACCUMULATEDTYPE`; cari `.AccumulationType`, tetap `NOTE IS NOT NULL`), `czone` (`CZONE`; cari `.Code`, tetap `GROUPOF
  IS NOT NULL`, urut `DESCRIPTION`). Keempat RD tanpa DISTINCT (pyGetDistinctRows false) — diikuti. Jalur 501 dicabut.
- **Migrasi 196** (ditulis, belum dijalankan; urutan … → 195 → **196**): view `PROVINCE`, `ACCUMULATEDTYPE`, `CZONE`,
  `ACCUMULATION` → tabel flat bernama sama. Per view: tabel `_SALIN` → salin isi view → `DROP VIEW` → `CREATE TABLE` nama
  view (kolom eksplisit, terbaca penjaga) → salin balik → `DROP TABLE _SALIN`; PROVINCE lebih dulu (view CITY me-JOIN-nya).
  Jalur mundur: `DROP TABLE` lalu `CREATE OR REPLACE FORCE VIEW` **persis** teks `DDL\<NAMA>.txt` (diperiksa: keempatnya
  sama setelah `POOLDATA.` → `{skema}.`). Kepemilikan: `ACCUMULATION` keluar dari bab "Tabel warisan" `MODUL.md`
  (keputusan work owner); `NATION` masuk sebagai tabel warisan.
- ⚠️ Konsekuensi: tabel flat = salinan saat migrasi; perubahan `M_*` (JSON, mis. prosedur `RDBMASTERACCUMULATION` dari
  Add New akumulasi Pega) tidak lagi terlihat. `DROP VIEW` berlaku untuk SIAPA PUN yang membaca skema sasaran.

### Keputusan agent butir 103 (menunggu konfirmasi)

| # | Keputusan | Dasar |
| --- | --- | --- |
| A179 | Migrasi 196 hanya empat view atas JSON yang DIBACA nbfacin. **Ditunda:** `CITY` / `DISTRICT` — view atas `CITYINPUT` / `DISTRICTINPUT` / `RWINPUT` / `BRANCH` / `PROVINCE` yang DDL-nya tidak ada, sehingga tipe kolom tabel flat tidak dapat ditulis tanpa menebak; `V_JN_OBJ_ITEM` (DDL view dikirim 04-10-2026) — view atas tabel `POOLDATA.OBJECTITEMTYPE` yang DDL-nya tidak ada (`PCTADJUSTABLE1/2` mungkin NUMBER — tidak ditebak); `COVERAGE_FACIN` (tidak dibaca lagi sejak butir 96), `BRANDDETAIL`, `VJ_M_TYPE_PROPERTY_PLAN` (tidak dibaca modul mana pun di repo) — memasukkannya ke migrasi nbfacin berarti nbfacin memiliki objek yang tidak dipakainya | Penjaga kepemilikan (`TestTabelBukanMilikKitaTidakDibuat`); CLAUDE.md §4.5 (tipe tidak ditebak). Dibutuhkan: DDL tabel dasar `CITYINPUT`, `DISTRICTINPUT`, `RWINPUT`, `BRANCH`, `OBJECTITEMTYPE` (atau keluaran `ALL_TAB_COLUMNS` view CITY / DISTRICT / V_JN_OBJ_ITEM dari DEV), dan keputusan modul pemilik tiga view sisanya |
| A180 | Tipe kolom keempat tabel flat = `VARCHAR2(4000)` | Setiap kolom view = JSON dot-notation, yang dikembalikan Oracle sebagai VARCHAR2(4000) — sifat Oracle, bukan dari korpus `[dugaan]`. Sebelum menjalankan, periksa di DEV: `SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, DATA_LENGTH FROM ALL_TAB_COLUMNS WHERE TABLE_NAME IN ('PROVINCE','ACCUMULATEDTYPE','CZONE','ACCUMULATION')` |

## Backend butir 2–3 (sesi c3, 04-10-2026)

Tanpa migrasi. `repository/akumulasi.go`, `services/akumulasi.go`, `handlers/akumulasi.go`.

**`GET /api/nbfacin/akumulasi`** → `{ baris: [{ id, accumulationName, note }] }`; parameter dipangkas, > 255 karakter → 400;
503 tanpa basis data. Jalur `[terverifikasi]` `Activity\GetDataAccumulation_act.xml` (dibaca utuh, sub-langkah 6.1–6.7):
- `postalCode` terisi, atau `cityId` / `districtId` / `policyNo` semuanya kosong → **RD** `SearchRiskAccumulation_RD`:
  view `ACCUMULATION` INNER JOIN `RW` atas ZIPCODE, DISTINCT sebelas kolom laporan, saringan diisi saja: `ID =`, `NOTE`
  Contains tidak peka huruf, `CZONE =`, `KEYWORD =`, `ZIPCODE =` (= `.PostalCode`), `PROVINCEID =`; ≤ 500.
- selain itu → **SQL**: `districtId` → `GetAccumulationDistrict_SQL`; `cityId` → `GetAccumulationProvince_SQL` (persis:
  `accumulation where ZIPCODE in (select zipcode from rw where cityid|DISTRICTID = …)`, accumulationName = kolom
  `ACCUMULATIONTYPE`); `policyNo` → `GetSummaryRiskAccumPolis_Sql` + langkah 6.6 (unik menurut AccumulationCode;
  accumulationName / note dari `ACCUMULATION`).

**`GET /api/nbfacin/akumulasi/saran/{jenis}?q=&induk=`** → `{ baris: [{ id, label, ekstra? }] }`. Medan cari / nilai
pilih `[terverifikasi]` `Section\SearchRiskAccumCov.xml` (`pyUseForSearch`, `pyPropertyTarget`):
- `city` — view `CITY`: label `NOTE` (cari), id `ID` (→ `InputFilter.City`), induk `PROVINCEID`; DISTINCT (ID, NOTE).
- `district` — view `DISTRICT`: label `DISTRICTNAME` (cari), id `ID` (→ `InputFilter.District`), induk `CITYNAME` (nama kota).
- `area` — tabel `RW`: label `NOTE` (cari), ekstra `ZIPCODE` (→ `InputFilter.PostalCode`), `STS_AKTIF = '1'`; id kosong.
- `nation` / `province` / `accumtype` / `czone` → **501** (A178). 400 jenis tak dikenal; 503.

**Bukti yang dicari dan hasilnya:**
- `GetAccumulation_SQL` kelas Int-ACCUMULATION (langkah 6.5, SyariahStatus) **TIDAK ADA** di korpus `[terverifikasi]` —
  hanya `ASM-FW-GISFW-INT-ACCUMULATION_LIFE!ASM!GETACCUMULATION_SQL` (3 salinan).
- Tabel kelas Int-NATION / PROVINCE / ACCUMULATEDTYPE: hanya `M_NATION` / `M_PROVINCE` / `M_ACCUMULATEDTYPE` (ID, OLDID,
  JSONDATA) — kunci JSON tidak ditebak; view `CITY` me-LEFT JOIN objek `PROVINCE` yang DDL-nya tidak ada. Int-CZONE: tabel
  `czone` hanya terlihat di SQL (`id`, `code`), DDL tidak ada; RD menyaring `.GroupOf` yang kolomnya tidak terverifikasi.
- BrowseRW_RD: pyParameters `City, District, Province, Teritory, ZipCode`; layar mengirim `DistrictName` → saringan D
  dibuang di Pega (A177). `.ID` bukan kolom laporan RD.

### Keputusan agent — ✅ A172–A178 DISETUJUI work owner 04-10-2026 (butir 102; A178 kemudian diganti butir 103)

| # | Keputusan | Dasar |
| --- | --- | --- |
| A172 | Jalur: kecamatan > kota (RDB-List kemudian menimpa `pyReportContentPage` `[dugaan]`); `policyNo` bersama kota / kecamatan → hasil kosong; `syariahStatus` diterima tetapi diabaikan | Langkah 6.2–6.4 menulis halaman yang sama berurutan; langkah 6.6 lalu memindai baris kota / kecamatan tanpa CARI10 → satu baris kosong; langkah 6.5 memanggil aturan yang tidak ada, RD tidak menyaring syariah |
| A173 | Saringan RD C `.AccumulationName = Param.Coverage` tidak dipakai | Sumber param `SearchAccumulation.AccumulationType` tidak pernah terisi (K-8 sesi 0f) |
| A174 | Urutan hasil `ID, NOTE` (RD / SQL tanpa urutan); jalur SQL ikut dibatasi 500 | Hasil deterministik; RDB-List tanpa batas, grid Pega berpaging |
| A175 | Jalur nomor polis: SQL korpus dipersempit ke AccumulationCode (JSON_TABLE jalur sama), subkueri nama / note memakai `MAX` (SQL asal ORA-01427 bila ID ganda), kode kosong dibuang, urut kode | Langkah 6.6 hanya memakai CARI10 / 12 / 13; SELECT biasa (ADR-0043) |
| A176 | Saran: `q` dicocokkan **Contains tidak peka huruf** atas kolom label (medan `pyUseForSearch`); paling banyak 50; urut label | Mode cocok autocomplete tidak tertulis di korpus; pola A124 (saran RW 50) |
| A177 | Saran `area` **tanpa** saringan induk (`induk` diabaikan), seperti Pega | Param layar `DistrictName` bukan param RD (`District`) → saringan dibuang di Pega `[terverifikasi]` |
| ~~A178~~ | ~~Saran `nation` / `province` / `accumtype` / `czone` → 501 sampai DDL sumbernya ada~~ | Tabel / view kelas Int-NATION / PROVINCE / ACCUMULATEDTYPE / CZONE tidak ada di `DDL\` | ⛔ **DIGANTI** butir 103: DDL dikirim 04-10-2026, ketujuh saran kini dilayani.
| — | Kelas Int-CITY / Int-DISTRICT = view `CITY` / `DISTRICT` `[dugaan]` | Nama dan kolom sama persis; tidak ada SQL kelas itu yang menyebut tabelnya |

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

## Add New accumulation (04-10-2026)

Permintaan work owner: "tombol add accumulationnya mana? itu contoh saat di klik tombol add" (gambar layar Pega DEV).

Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`:
- `SearchRiskAccumCov` sel 39 "Add New": pre-DT `InputAccumulationAdd_PreDT` (mengosongkan isian, `CARI30 = 1`) +
  `GetCzone_Act` (Cresta Zone dari zip saringan: RDBList `GetAccumulationZipcode_SQL`).
- `Section\InputAccumulationCov.xml`: Accumulation Type* (sel 10) + Add (12), Province* (13), Scope Area* (16), Cresta
  Zone* (17), Primary Zip Code* baca-saja (20) + Choose Zip Code (21), Key Word* (24), Accumulation Description* (25),
  Save (28, `SaveAccumulation_Act` → RDBList `UpdateMasterAccumulation` → PROSEDUR `RDBMASTERACCUMULATION`), Close (29,
  `CancelInputAccumulation_Act`).
- `FlowAction\ChooseZipCode` → `Section\ChooseZipCodeDtl` (RD `BrowseRiskAddressZipCode_RD`; Choose → `SetPostalCode_Act`:
  PostalCode, NationInitial lewat `SetCountryID_Act`, lalu `GetCzone_Act`).
- Aturan properti `.ScopeArea` TIDAK ada di korpus (dicari menurut pxInsName di NB FacIn dan `DDL\`).

Frontend (sesi 0f): `components/FormTambahAkumulasi.tsx` (+ uji), `PopupAkumulasi.tsx` (tombol Add New, mode tambah),
`api.ts` (`tambahAkumulasi`, `czoneDariZip`, `cariZipAkumulasi`), `labels.ts` (`TAMBAH_AKUMULASI`, `OPSI_SCOPE_AREA` kosong,
`TEKS_TAMBAH_AKUMULASI`). Backend = sesi c3 (`POST …/akumulasi`, `GET …/akumulasi/czone`, `GET …/akumulasi/zipcode`;
simpan lewat mesin `inti/backend/master`, ADR-0043 tanpa CALL).

Keputusan agent (menunggu konfirmasi):
- ~~**T-1** Choose Zip Code = panel di dalam form~~ — DICABUT 04-10-2026 (work owner: "saat klik choose zipcode,
  dibuatin pop up aja dong, jgn dibwh gini"): popup di atas popup Choose Accumulation lewat portal ke `document.body`
  (kelas modal inti, akar `.nbfacin`, Escape ditangkap di fase capture `window` supaya hanya popup zip yang tertutup).
- **T-2** Save berhasil → accumulation baru langsung dipakai coverage dan popup ditutup (Pega menulisnya ke coverage lalu
  kembali ke mode cari tanpa menutup popup).
- **T-3** Save ganda (Note + zip sudah ada) menampilkan pesan prosedur dan TIDAK menulis apa pun ke coverage (Pega menulis
  "UNKNOWNID").
- **T-4** Tombol "Add" tipe akumulasi baru (`ModalAccumulation_FacIn` → `RDBMASTERACCUMULATEDTYPE`) = tahap berikut; menu
  Master Accumulated Type menggantikannya.
- **T-5** Province dari tabel PROVINCE (RD `BrowseProvince_RD` kelas PROVINCE yang dipakai section tidak ada di korpus).
- ⛔ Scope Area (wajib) belum berpilihan sampai XML aturan `.ScopeArea` dikirim work owner.

### Backend Add New (sesi c3, 04-10-2026)

- `POST /api/nbfacin/akumulasi` badan `{accumulation, accumulationType, note, keyword, scopeArea, cZone, cZoneId,
  provinceId, zipCode, negara}` (teks; kunci asing 400) → 201 `{id, note}` (note huruf besar). Identitas wajib (401).
  Urutan periksa = `SaveAccumulation_Act` `[terverifikasi]`: zip / negara (CARI34) / CZone / Note kosong → 400 pesan
  verbatim "Postal code, Nation, CZone and Accumulation Description can't be null!"; lalu wajib section → 400
  "<medan> wajib diisi"; lalu mesin master inti (`services.Tambah` master "accumulation": rujukan Accumulation Type /
  Province harus ada, lebar, ganda Note + zip → 409 pesan verbatim prosedur "Error master item Akumulasi Sudah Ada",
  ID negara-zip-urutan, Note huruf besar, CREATE_OP = pelaku). Tanpa basis data → 503.
- `GET /api/nbfacin/akumulasi/czone?zip=` → `{cZoneId, cZone}`; zip kosong / tanpa CZone = string kosong.
- `GET /api/nbfacin/akumulasi/zipcode?provinceName=&q=&halaman=` → `{baris: [{zipCode, city, province, nation,
  nationInitial}], total, halaman, ukuran: 15}` (A185; halaman ≥ 1, bawaan 1, bukan angka / < 1 → 400, di luar jangkauan
  → baris kosong).
- Saran `accumtype` kini ber-`ekstra` = `NOTE` (kolom tambahan `.Note` → AccumulationName di section).
- Berkas: `repository/akumulasibaru.go` (+ uji), `services/akumulasibaru.go`, `services/layanan.go` (penambah = mesin
  `inti/backend/master/services`), `handlers/akumulasi.go`, `handlers.go` (rute, peta galat), `handlers/akumulasibaru_test.go`,
  `models/objek.go`, `repository/akumulasi.go` (+ uji saran).

| # | Keputusan agent (menunggu konfirmasi) | Alasan |
| --- | --- | --- |
| A181 | ~~CZone dari zip: baris pertama urut ID; hanya CZONE aktif; RW tidak disaring~~ → **diganti 04-10-2026** (work owner: "pada saat add new, ada kolom Cresta Zone, itu otomatis keisi dari tabel rw, diambil dari zcone zipcode yang dipilih"): `cZone` = `RW.CZONE` ber-`ZIPCODE` zip itu — beberapa nilai → `MIN` sesudah NULL / spasi dibuang; ~~`RW.STS_AKTIF` TIDAK disaring~~ → **HANYA RW AKTIF** (`STS_AKTIF = '1'`; work owner 04-10-2026: "yang aktif saja, dan di dalam master datanya tidak ada yg czone nya kosong"); `cZoneId` = `MIN(CZONE.ID)` ber-`CODE` itu bila ada, tanpa syarat aktif, selain itu ""; zip kosong / tanpa CZONE di RW → keduanya "" | MENYIMPANG dari `GetAccumulationZipcode_SQL` `[terverifikasi]` (CZONE yang CODE-nya ada di RW zip itu, baris 1) — kueri itu kosong bila CZONE RW tidak ada / tidak aktif di tabel CZONE `[dugaan]` sebab Cresta Zone kosong di DEV. Saringan aktif = saringan RD BrowseRW / saran area. Save (POST) TIDAK memeriksa cZone / cZoneId ke tabel CZONE (mesin master: kolom biasa; rujukan hanya Accumulation Type dan Province) — nilai RW yang tak ada di CZONE tetap tersimpan, dengan cZoneId kosong |
| A182 | Zip Code: `q` = Contains tidak peka huruf atas zip (tambahan); hasil **DISTINCT** atas lima kolom tampil; urut zip, kota, provinsi, negara; NationInitial = `MAX` NATION ber-`ID = IDNATION OR NOTE = NATIONNAME`; maks 100 | RD tanpa DISTINCT atas RISKADDRESS (satu baris per alamat) → 100 baris bisa satu zip berulang; SetCountryID_Act pxResults(1) tanpa urutan; CARI33 / CARI38 = IDNation / NationName `[dugaan kuat]` (`ChooseZipCodeDtl`) |
| A183 | Wajib section (Accumulation Type, Province, Key Word, Scope Area) diperiksa backend juga, sesudah wajib Pega | Layar menandainya wajib; korpus tanpa pesan untuk itu |
| A184 | `ACCUMULATIONNAME` tetap TURUNAN `ACCUMULATEDTYPE.ACCUMULATIONTYPE` (mesin master), bukan `.Note` seperti section | View warisan `DDL\ACCUMULATION.txt` (yang dibaca semua pembaca) menurunkan AccumulationName dari AccumulationType tipe — JSON AccumulationName dari form tidak pernah terlihat lewat view. `accumulationType` (nama tipe) disimpan ke kolom ACCUMULATIONTYPE seperti JSON asal |
| A185 | Choose Zip Code **berhalaman di server**, 15 baris (work owner 04-10-2026: "ada page nya, 15 list perpage, datanya diambil dari kolom riskaddress berdasarkan ProvinceName yang dipilih"); `total` = COUNT atas kueri yang sama; urut kelima kolom (halaman stabil); batas pyMaxRecords 100 DICABUT. Saringan provinsi TETAP **Contains** tidak peka huruf seperti RD (bukan "="): "berdasarkan ProvinceName yang dipilih" tidak menyebut cara cocok, dan RD `.ProvinceName Contains Param.ProvinceName` satu-satunya bukti — akibatnya nama provinsi yang menjadi bagian nama lain ikut (mis. "PAPUA" juga "PAPUA BARAT"), sama dengan Pega | Gambar Pega DEV "Page 1 of 2421" (10 baris) = RD dipakai berhalaman; ⚠️ total di sini menghitung baris DISTINCT (A182), Pega menghitung baris RISKADDRESS mentah — jumlah halaman berbeda |
| A187 | Choose Zip Code hanya zip yang punya baris **RW AKTIF ber-CZONE terisi**: `EXISTS (RW r WHERE r.ZIPCODE = a.POSTALCODE AND r.STS_AKTIF = '1' AND TRIM(r.CZONE) IS NOT NULL)` — untuk hitungan dan halaman sekaligus (satu kueri dasar); zip terpilih selalu ber-Cresta Zone | Work owner 04-10-2026: "yang aktif saja, dan di dalam master datanya tidak ada yg czone nya kosong" (tafsiran sesi 9d `[dugaan]`, ditawarkan untuk dibatalkan). MENYIMPANG dari RD BrowseRiskAddressZipCode_RD (tidak menyaring RW). Kolom zip `[terverifikasi]` DDL: `RISKADDRESS.POSTALCODE`, `RW.ZIPCODE` |
| A186 | `scopeArea` wajib salah satu dari `AREA, DISTRICT, CITY, PROVINCE, COUNTRY` (selain itu 400) | `DDL\ScopeArea.xml` (ASM-FW-GISFW-INT-ACCUMULATION!SCOPEAREA, pyStandardValue) `[terverifikasi]`, dikirim work owner 04-10-2026; layar Pega = daftar pilihan. Menu Master Accumulation (mesin inti) belum memeriksa daftar ini |
