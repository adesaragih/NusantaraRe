# 36: Choose / Clear Risk Address — popup alamat risiko dari tabel RISKADDRESS

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026:
> *"lanjut ke choose/clear risk address"*, disertai tangkapan layar popup (tidak disalin: memuat alamat).

**What to build:** di sub-tab Object Address (tab Object FIRE, tiket 35), **Choose Risk Address** membuka popup
*Choose Risk Location*: saringan di kiri, grid alamat di kanan, **Pilih** per baris, lalu **Search** / **Add** di kaki.
Pilih mengisi blok Risk Address objek. **Clear Risk Address** mengosongkannya.

**Blocked by:** tiket 35.

**Status:** frontend selesai 03-10-2026 (uji hijau); backend → sesi c3.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\ObjectDetails.xml` sel 16 `Choose Risk Address`: runActivity `SetDefaultSearchRiskLocation_Act` →
  showHarness `ChooseRiskAddress`, WindowName `Choose Risk Location`, Target popup. Sel 17 `Clear Risk Address` →
  `ClearRiskLocation_act`.
- `Section\ChooseRiskAddress.xml`: saringan sel 78–84 (`pyLabelFieldValue`) Address · Zip Code · Country · Province ·
  City · District · Territory, terikat ke `.Property.RoadName`, `.Property.RiskLocation.ASMZipCode`, `.Property.Country`,
  `.Property.Province`, `.Property.RiskLocation.ASMCity` / `ASMDistrict` / `ASMRW`. Sel 73 `Search` →
  `SearchRiskAddressAct`; sel 74 `Add` → `setOutputValue_DT`, `ClearPageRiskAddress`, harness `ChooseRiskLocation`.
- `Section\ChooseRiskAddress_ResultList.xml`: kolom sel 19–26 Type · Address · Country · Province · City · District ·
  Territory · Zip Code (`.Title`, `.Address`, `.NationName`, `.ProvinceName`, `.CityName`, `.DistrictName`,
  `.TerritoryName`, `.PostalCode`); `pyPageSize` 10, Numeric; Pilih (ikon `IconChoose.png`) → DT `SetRiskIdDT_FacIn` →
  `ObjSave` (tidak ada di korpus) → `opener.location.reload` → `window.close`.
- `Activity\SetDefaultSearchRiskLocation_Act.xml`: kunci cari palsu `RoadName = "z"`, `ZipCode = "123456"`, sehingga grid
  kosong sampai Search.
- `Activity\SearchRiskAddressAct.xml`: bila semua saringan kosong mencari dengan kunci palsu (`"z"` / `"123456"`) → hasil
  kosong; nilai saringan dijadikan huruf besar kecuali Zip Code.
- `DataTransform\SetRiskIdDT_FacIn.xml`: AlmRiskID = ID; RoadType = Title; RoadName = Address; ASMRW = TerritoryName;
  ASMDistrict; ASMCity; Province; Country = NationName; ASMZipCode = PostalCode;
  ASMAddress = `Title+" "+Address+","+Territory+","+District+","+City+","+Province+","+Nation`; BuildingNo = "".
- `Activity\ClearRiskLocation_act.xml`: mengosongkan RoadType, RoadName, ASMZipCode, Province, Country, ASMRW, ASMCity,
  ASMDistrict, ASMAddress, AlmRiskID (Building No. tidak ikut).
- `ReportDefinition\BrowseRisksAddress_RD.xml` (kelas `ASM-FW-GISFW-Int-RISKADDRESS`): parameter JALAN, KODEPOS,
  COUNTRY, PROVINCE, CITY, DISTRICT, TERRITORY; `pyMaxRecords` 100. DDL `D:\migrasi\RNM\DDL\RISKADDRESS.txt`.

## Keputusan agent

- **H-1** Saringan = salinan medan objek saat popup dibuka, bukan terikat ke objek: mengetik di saringan tidak
  mengubah objek sebelum Pilih. Di Pega saringan terikat langsung ke `.Property.*`.
- **H-2** Grid kosong sampai Search (padanan kunci palsu `SetDefaultSearchRiskLocation_Act`). Search tanpa saringan
  tidak memanggil server dan menampilkan pesan sistem baru.
- **H-3** Pilih mengisi objek di layar; tersimpan lewat Save tab Object (sejalan E-1). Di Pega `ObjSave` dijalankan
  saat Pilih.
- **H-4** `Add` (alamat baru lewat harness `ChooseRiskLocation`) = tahap berikut, sementara nonaktif.
- **H-5** Label XML ("Country", "City", "Zip Code") dipakai, bukan huruf besar di tangkapan layar ("COUNTRY", "CITY",
  "Zip code"), sejalan label blok Risk Address tiket 35. Teks tombol `Pilih` diambil dari tangkapan layar.

## Kontrak

- `GET /api/nbfacin/risk-address?address=&zipCode=&country=&province=&city=&district=&territory=&halaman=` →
  `{ baris: [{ id, title, address, nationName, provinceName, cityName, districtName, territoryName, postalCode }],
  total, halaman, ukuran: 10 }`. Hanya saringan terisi yang dikirim; tanpa satu pun saringan → 400.
- Penyimpanan: tidak ada perubahan. Medan Risk Address ikut `PUT …/objek` (tiket 35).

## Acceptance criteria

- [x] Choose Risk Address membuka popup; saringan diisi awal dari objek; label diuji ke korpus.
- [x] Search: grid 10 per halaman; Pilih mengisi Risk Address (`SetRiskIdDT_FacIn`); Clear mengosongkan.
- [x] Backend `GET /api/nbfacin/risk-address` sesuai filter `BrowseRisksAddress_RD`.
- [ ] Add (tahap berikut).

## Backend (sesi c3, 03-10-2026) — disusun agent

`GET /api/nbfacin/risk-address?address=&zipCode=&country=&province=&city=&district=&territory=&halaman=` →
`{"baris":[{"id","title","address","nationName","provinceName","cityName","districtName","territoryName","postalCode"}],
"total","halaman","ukuran":10}`. 400: tanpa satu pun saringan (spasi saja = kosong), halaman tak sah, saringan > 255
karakter; 503; 500. Tanpa identitas (pola lookup tiket 27/28/33). Tanpa migrasi.

**Laporan RD `[terverifikasi]`** `NB FacIn\ReportDefinition\BrowseRisksAddress_RD.xml` (kelas `ASM-FW-GISFW-Int-RISKADDRESS`), `pyContent`:

| Filter | Kolom | Relasi | pyCaseInsensitive | Parameter |
| --- | --- | --- | --- | --- |
| A | `.Address` | Contains | **true** | `Param.JALAN` (pyUseNullIfEmpty false) |
| B | `.PostalCode` | Contains | — | `Param.KODEPOS` |
| C | `.NationName` | Contains | — | `Param.COUNTRY` |
| D | `.ProvinceName` | Contains | — | `Param.PROVINCE` |
| E | `.CityName` | Contains | — | `Param.CITY` |
| F | `.DistrictName` | Contains | — | `Param.DISTRICT` |
| G | `.TerritoryName` | Contains | — | `Param.TERRITORY` |

Logika `A AND B AND C AND D AND E AND F AND G`; `pyGetDistinctRows` **true**; **tanpa urutan** (tak satu pun `pySortOrder`
selain 99999); `pyMaxRecords` **100**, paging aktif. **`RW.ZipCode`** = filter **JOIN**, bukan saringan: `pyJoinInfo` kelas
`ASM-FW-GISFW-Int-RW`, prefix `RW`, **INNER**, `.PostalCode = RW.ZipCode` — tabel `POOLDATA.RW` kolom `ZIPCODE` (DDL
`RW.txt`); hanya alamat yang kode posnya ada di RW yang tampil. **Dipakai** (dibangun sebagai INNER JOIN), tidak dibuang.

**Keputusan agent — DISETUJUI work owner 03-10-2026** (butir 91, diteruskan sesi `nusantarare-0f`: *"setuju
A116–A122"*):

| # | Keputusan | Dasar |
| --- | --- | --- |
| A116 | Semua tujuh saringan **tidak peka huruf** (`UPPER(kolom) LIKE` pola huruf besar, `ESCAPE`); di RD hanya A yang `pyCaseInsensitive` — B–G membandingkan kolom apa adanya. `SearchRiskAddressAct` meng-UPPER nilai A dan C–G, **Zip Code (B) tidak** (baris 600 disalin apa adanya) | permintaan sesi 0f; sama dengan Pega bila data tersimpan huruf besar `[dugaan]` |
| A117 | Batas **100 baris** RD diikuti: total ≤ 100, halaman 11 ke atas kosong | `pyMaxRecords` 100 |
| A118 | Urut **sembilan kolom hasil** (`ID` dulu; RD tanpa urutan) — `ID` tidak ber-PK di DDL dan DISTINCT atas sembilan kolom, jadi `ID` saja tidak menjamin deterministik | — |
| A119 | Tiap saringan ≤ 255 karakter → 400 | pola A73 (kolom VARCHAR2(4000)) |
| A120 | Tanpa identitas | pola lookup tiket 27/28/33 |
| A121 | Saringan kosong dilewati (WHERE hanya untuk yang terisi); **semua kosong = 400** | perilaku Pega untuk parameter kosong tanpa "use null if empty" `[dugaan]`. Pega bila semua kosong TETAP mencari dengan kunci palsu `RoadName = "z"`, `ZipCode = "123456"` (`SearchRiskAddressAct` langkah "Jika filter kosong", baris 930–989) → hasil kosong; 400 = penyimpangan sadar (ralat: bukan "tidak mencari") |
| A122 | Kelas `ASM-FW-GISFW-Int-RW` = tabel `POOLDATA.RW` | RDB-List kelas yang sama `NB FacIn\RDBList\BrowseRW_SQL.xml` membaca `FROM pooldata.rw`; aturan pemetaan kelas→tabel tidak ada di korpus `[dugaan]` |

**Bug DEV 03-10-2026 (butir 87):** Save tab Object menolak Risk Location > 50 bita - rangkaian alamat dari popup ini
selalu lebih panjang. Migrasi **192**: Risk Location (`ASM_ADDRESS`) dan Address (`ROAD_NAME`) VARCHAR2(4000) (=
RISKADDRESS); delapan kolom alamat lain VARCHAR2(100) (ralat butir 88); batas validasi ikut 4000 / 100.
