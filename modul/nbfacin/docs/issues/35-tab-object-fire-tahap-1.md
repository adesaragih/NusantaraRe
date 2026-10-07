# 35: Tab Object (FIRE) tahap 1 — grid objek, sub-tab Object Address, Save

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026:
> *"setelah itu masuk ke tab object"*, disertai tangkapan layar tab Object kasus FIRE (tidak disalin: memuat data
> nasabah).

**What to build:** isi tab **Object** layar Inward Facultative untuk kasus FIRE. Isinya grid objek (Top Risk · No. ·
Object Name · Location, tombol Tambah / Hapus, baris dapat dibuka) dengan tujuh sub-tab per baris. Sub-tab **Object
Address** berisi medan objek, blok Risk Address, dan blok Building Construction. Tombol **Save** menyimpan seluruh
daftar.

**Blocked by:** —

**Status:** frontend selesai 03-10-2026 (uji hijau); backend → sesi c3; migrasi **ditulis, tidak dijalankan**.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\InputInwardFacultativeDtl.xml`: tab `Object` → sel 21 include `ObjectList` (visible `IsFire`); sel 22
  `InputDtlObject_FacIn` (`!IsFire`); sel 25 `Save` (pyLabel) → `SaveFacIn_Act`.
- `Section\ObjectList.xml`: grid `.LocationList`, `pyRowEditing=masterDetail`, `pyPageSize` 10; kepala sel 13–16
  `Top Risk` / `No.` / `Object Name` / `Location` (`pyValue`); baris `.Property.IsTopRisk`, `.Property.ObjectNo`,
  `.Property.ObjectName`, `.Property.RiskLocation.ASMAddress`; tambah = ikon → `AddLocation_act`; hapus = ikon.
- `Activity\AddLocation_act.xml`: `LocationList(<APPEND>)`, `ObjectNo = pxListSubscript`.
- `Section\Property.xml`: sub-tab (`pyTitle`) Object Address · Surrounding Risk · Object Item · Occupation · FEA ·
  Loss Record · Loss Record Internal.
- `Section\ObjectDetails.xml`: sel 5 Object No. · 6 Object Type (wajib, "Please Select") · 11 Material Damage · 12 Top
  Risk · 13 Object Name (visible `ObjectType = 'Others'`) · 16 Choose Risk Address · 17 Clear Risk Address · 28 Type ·
  29 Address · 30 Building No. · 31 Zip Code · 32 Country · 35 Territory · 36 City · 37 District · 38 Province · 39 Risk
  Address ID · 40 Risk Location · 49 Number of Floor · 50 Roof Type · 51 Wall Type · 52 Floor Type · 56 Partition Type ·
  57 Support Wall Type · 58 Other Type. Properti per sel diuji `labels.test.ts`.
- `Activity\SetValueOnObjectName_Act.xml`: ObjectName = ObjectType; dikosongkan bila `'Lainnya'`.
- `Activity\SetErrorMessageFloorNumber_Act.xml`: `NumberOfFloor < 0` → "Floor number can't be minus".
- `Activity\InsertUploadFire_act.xml`: himpunan Object Type (9 nilai + "Others").
- Rancangan flat `loader/skema_gen.go`: `T_LOCATIONLIST` (induk T_GENERAL_POLIS) → `T_PROPERTY` → `T_RISKLOCATION` /
  `T_BUILDINGCONSTRUCTION`. `T_BUILDINGCONSTRUCTION` belum punya PartitionType / SupportWallType / OthersType.

## Keputusan agent

- **G-1** Choose / Clear Risk Address nonaktif sampai tahap 2 (popup `ChooseRiskAddress`, RD `BrowseRisksAddress_RD`
  atas tabel RISKADDRESS). Medan Risk Address tampil-saja, kecuali Building No.
- **G-2** `SetValueOnObjectName_Act` memeriksa `'Lainnya'`, sedangkan daftar dan syarat tampil sel 13 memakai `'Others'`.
  Diikuti `'Others'`: Object Name = Object Type, atau kosong lalu diisi manual bila Others.
- **G-3** Sub-tab selain Object Address = tahap berikut (`BelumTersedia`).
- **G-4** Daftar Roof / Wall / Floor Type tidak ada di korpus (`pyListSource=associated`, aturan properti hilang; data
  contoh berisi kode `14`, `9`, `KELAS I`). Dropdown hanya memuat nilai tersimpan; nilai awal Pega belum diterapkan.
  **Menunggu daftar dari work owner.**
- **G-5** Kasus FIRE = Group Business "FIRE". When `IsFire` Pega lebih luas (IsKPR, IsOilGas, IsFireStyle1/2,
  BusinessType "Fire"). Selain FIRE, tab Object tetap `BelumTersedia`.
- **G-6** Hapus = buang baris (jalur NB `deleteRow`); jalur EDM `SetFlagDeleteEDM_Act` bukan di sini. Object No. tidak
  dinomori ulang.
- **G-7** Urutan pilihan Object Type = urutan ekspresi `InsertUploadFire_act` `[dugaan]`.
- **G-8** Save menahan simpan bila ada Object Type kosong atau lantai minus, lalu membuka baris yang salah.

## Kontrak

- `GET /api/nbfacin/kasus/{caseId}/objek` → `{ baris: ObjekFire[] }`, urut `.LocationList`.
- `PUT /api/nbfacin/kasus/{caseId}/objek` badan `{ baris: ObjekFire[] }`: daftar diganti utuh dalam satu transaksi,
  dijawab `{ baris }` hasil baca ulang.
- `ObjekFire` (lihat `frontend/api.ts`): objectNo, objectType, objectName, isMaterialDamage, isTopRisk, roadType,
  roadName, buildingNo, zipCode, country, riskLocation, territory, city, district, province, riskAddressId,
  numberOfFloor, roofType, wallType, floorType, partitionType, supportWallType, otherType — semua teks kecuali dua
  boolean.

## Acceptance criteria

- [x] Tab Object kasus FIRE: grid, Tambah / Hapus, baris dibuka → tujuh sub-tab; label diuji ke korpus.
- [x] Object Address: medan, Risk Address (tampil), Building Construction; Object Name ikut Object Type.
- [x] Lantai minus dan Object Type kosong menahan Save.
- [x] Backend GET / PUT `…/objek` ke tabel rancangan; migrasi ditulis (tidak dijalankan) — 186.
- [ ] Daftar Roof / Wall / Floor Type (G-4).
- [ ] Tahap 2: Choose / Clear Risk Address. Tahap 3: sub-tab lain.

## Backend (sesi c3, 03-10-2026) — disusun agent

- `GET /api/nbfacin/kasus/{caseId}/objek` → `{"baris": ObjekFire[]}` urut `T_LOCATIONLIST.SEQ_NO` (selalu larik); 404 case
  tidak ada / bukan `FAC`; 503; 500. Tanpa identitas (pola baca).
- `PUT …/objek` badan `{"baris"}` → daftar **diganti utuh** dalam satu transaksi (sentuh `T_WORK_POLIS` → pastikan baris
  `T_GENERAL_POLIS` → hapus anak→induk → sisip urut `SEQ_NO` 1..n) → `{"baris"}` baca ulang. 400: `objectType` kosong,
  `numberOfFloor` bukan kosong / bilangan bulat ≥ 0, teks melebihi lebar kolom — pesan menyebut `baris[n].medan`; 401; 404;
  503. Kosong = daftar dikosongkan. Server tidak mencocokkan Object Type / Roof / Wall / Floor ke daftar pilihan —
  daftar layar sudah ada (G-9, urutan Object Type = tangkapan layar; G-4/G-7 selesai di sisi frontend); nilai disimpan
  apa adanya sesuai lebar kolom.
- **G-9 `[terverifikasi]`** (03-10-2026, berkas properti Pega yang ditambahkan work owner `D:\migrasi\RNM\DDL\RoofType.xml`
  / `WallType.xml` / `FloorType.xml`, PromptList `ASM-FW-GISFW-DATA-BUILDINGCONSTRUCTION!ROOFTYPE` / `!WALLTYPE` /
  `!FLOORTYPE`; nilai–label dipasangkan **per rowdata**, bukan per urutan): Roof `"1"`..`"14"` (`14` = Lain-lain), Wall
  `"1"`..`"9"` (`9` = Lain-lain), Floor **`KELAS III` = Keramik, `KELAS II` = Kayu, `KELAS I` = Lain-lain** (ralat atas
  "Keramik"/"Kayu" sebagai nilai). Ketiganya tanpa `pyDefaultValue` — nilai awal 14 / 9 / KELAS I tetap dari data contoh +
  tangkapan layar. Kolom 50 bita memuat semua nilai; server tetap tanpa validasi enumerasi.
- Migrasi **186** (ditulis, tidak dijalankan; butuh 182): lihat `STRUKTUR-TABEL-NB-FACIN.md`.

**Pemetaan `[terverifikasi]`** `Section\ObjectDetails.xml` (properti per sel) → kolom rancangan: Territory = `ASMRW` (sel
35), Risk Location = `ASMAddress` (sel 40), Zip Code = `ASMZipCode` (31), City/District = `ASMCity`/`ASMDistrict`, Province /
Country / Building No. / Road Type / Road Name / Risk Address ID (`AlmRiskID`) di `T_PROPERTY`. Boolean IS_* = teks
`"true"`/`"false"` — dihitung **per wadah halaman `Property`** (ralat: angka semula "4×true / 109×false lintas wadah"
menjumlahkan `CargoList`, jebakan sensus 5): `IsTopRisk` 2×`"true"`, 5×`"false"`, 1× kunci tidak ada; `IsMaterialDamage`
6×`"true"`, **0×`"false"`**, 2× kunci tidak ada → menulis `"false"` untuk Material Damage tak dicentang `[dugaan]` (A115).

**Keputusan agent (menunggu konfirmasi):**

| # | Keputusan | Dasar |
| --- | --- | --- |
| A109 | `T_LOCATIONLIST` / `T_PROPERTY` dibuat **sebagian** (tanpa kolom uang `LOSS_RATIO*_AMOUNT`, `TOTAL_TSI` dan kolom tak dipakai layar) | pola keputusan work owner butir 78.4 |
| A110 | Tiga kolom baru `PARTITION_TYPE`, `SUPPORT_WALL_TYPE`, `OTHERS_TYPE` VARCHAR2(50) + amandemen loader (`amandemenBangunan`) | properti ada di layar (sel 56–58), tidak di rancangan; tipe = kolom saudaranya |
| A111 | `numberOfFloor` sah = kosong atau `^[0-9]+$` (tanpa tanda, tanpa desimal) | `SetErrorMessageFloorNumber_Act` (`< 0` ditolak) + permintaan "bilangan bulat ≥ 0" |
| A112 | Badan PUT ≤ 1 MiB; paling banyak 99.999 baris (`SEQ_NO` NUMBER(5)) → 400 | lebar kolom |
| A113 | Satu baris Property / RiskLocation / BuildingConstruction per induk (`UNIQUE PARENT_ID`); PUT tanpa medan anak tetap membuat baris anak (kolom kosong) | halaman tunggal di rancangan (pola A87) |
| A114 | Baris `T_GENERAL_POLIS` kosong dibuat bila General belum pernah disimpan | FK `T_LOCATIONLIST.PARENT_ID` → `T_GENERAL_POLIS` |
| A115 | Boolean `false` ditulis teks `"false"` (juga `IsMaterialDamage`, yang di fixture tidak pernah `"false"` — kuncinya hilang); baca: hanya `"true"` = benar | fixture `Property` (lihat atas) `[dugaan]` |
