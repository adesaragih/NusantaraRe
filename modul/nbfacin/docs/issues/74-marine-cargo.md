# 74: Class of Business MARINE CARGO — port layar & alur

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 10-10-2026 (case MARINE CARGO,
> tab Object berisi "Tab Object untuk Class of Business MARINE CARGO belum tersedia di sistem baru (saat ini baru FIRE)")
> "PERBAIKI JANGAN HALU". Riset: empat subagen sesi 0f (Object, Coverage, Spreading + tab lain, Activity), laporan lengkap di
> transkrip sesi. Klaim inti diverifikasi ulang sesi 0f (lihat "Verifikasi ulang").

`NB\` = `D:\migrasi\RNM\NB FacIn\`. Label: `[terverifikasi]` = dibaca di baris itu; `[dugaan]` = tafsiran.

**Status:** needs-info. Spesifikasi selesai 10-10-2026, tetapi pembangunan menunggu dua hal: keputusan **K74-1…K74-6**, dan
nama tabel lookup di DEV (bagian "Data", Q-DBA).

## Temuan inti `[terverifikasi]`

- **Gerbang.** `BusinessType_DeT`: GroupPanel `004` (OldId apa pun) menghasilkan **"MarineCargo"**
  (`backend/services/lini.go` baris tabel `{"004", nil, "MarineCargo"}`). `NB\When\IsMarineCargo.xml` L44 `G OR H`:
  `.OfferFacIn.QuotationData.BusinessType = "MarineCargo"` (L486) atau `pyWorkPage.Quotation.BusinessType = "MarineCargo"`
  (L605). Blok "Kode Bisnis 24/18/…" (L164–269) hanya `pyConditionViewer` usang.
- **Bukan keluarga Aneka.** Hierarkinya **Cargo → Coverage → Spreading**, tanpa lokasi atau okupasi. Lapisan daun (layar
  spreading `SpreadingItem` / `InputDtlSpreadingCoverage_FacIn`) sama dengan FIRE. Rantai komisi memakai ulang
  `CommisionItemAneka`.
- **Belum ada apa pun di sistem baru.** Tabel `T_CARGOLIST` / `T_SHIP` / `T_LC` baru ada di rancangan
  `backend/services/loader/skema_gen.go` (dibangkitkan dari sampel), dan belum ada migrasinya. Yang sudah ada:
  - rumus premi `calculateMarine` (`services/premium/lini_lain.go` L229–240, tiket 18; 104/104 fixture);
  - jumlah pembayaran `PremiPaymentMarine` cabang NB (`services/pembayaran/marine.go`, tiket 19);
  - halaman alat periksa `frontend/pages/CoverageCargo.tsx` (tiket 21), yang tidak terhubung ke kasus.
- **Tab** (`NB\Section\InputInwardFacultativeDtl.xml`):
  - Object (L6909, isinya `InputDtlObject_FacIn` karena `!IsFire`);
  - Coverage (L12978);
  - Spreading (`!IsLife`);
  - Inw Fac Cedant Panels;
  - Premium Deduction / Brokerage Fee;
  - Payment (L94007);
  - Correspondence.

  **Clauses tersembunyi** (L17833 `IsOilGas || IsFire`), begitu juga Scoring Risk (`IsFire`) dan Loss Record (`IsIT`).

## Spesifikasi

### S1. Tab Object — grid CargoList `[terverifikasi]`

- **Isi tab hanya grid** `.OfferFacIn.CargoList` (`NB\Section\InputDtlObject_FacIn.xml` kontainer L1994 `IsMarineCargo`,
  grid L2849–2868):
  - Kolom: `No.` · `Top Risk` (cekbox, selalu nonaktif di grid) · `SailDate` (`.PolicyData.SailDate`, tanggal) · `Goods`
    (`.GoodNote`) · `Ship Name` (`.PolicyData.Ship.NM_SHIP`) · `Conveyance` (`.ConveyanceNote`).
  - Baris dibuka dengan expand pane berisi form `InputDtlCargo_FacIn` (`pyEditAction` L6260). Ukuran halaman 10.
- **Tombol "+"** (L3813–4409): tambah baris di akhir, lalu catat log "Tambah Obyek". Nonaktif bila `IsEDMRiSlip`.
- **Tombol hapus** (L5403–6200):
  - tampil bila `.FlagDelete != 1`;
  - NB (`StatusBusiness != 3`) menghapus baris;
  - EDM (`= 3`) menjalankan `SetFlagDeleteEDM_Act`, yaitu hapus lunak.
- **Blok tersembunyi, tidak diport:** grid Trading / Goods / Conveyance / tombol Location Limit (L7081 `NEVER`) dan sisipan
  `MarineCargoDtl` langsung (L20337 `1=2`).
- **Tahap UW** (`InputDtlObject_FacIn_IsUW`, edit action `InputDtlCargo_FacIn_IsUW`):
  - kolom grid menjadi No., Top Risk, Trading, Goods, Packing, Conveyance (`*Note`);
  - tanpa tombol tambah dan hapus;
  - semua medan read-only **kecuali Top Risk** (`InputDtlCargo_FacIn_IsUW` L500 `Editable`).

### S2. Form cargo (`InputDtlCargo_FacIn` + `MarineCargoDtl`) `[terverifikasi]`

Aturan read-only: semua medan `pyEditOptions=Read-only`. Yang bersyarat dapat diubah di NB dan read-only saat syaratnya
benar. Pola kunci yang sama dengan yang sudah ditetapkan di tiket 73:

- `InputDtlCargo_FacIn`: `pyPortal.IsShowOpenPolicyMarine == true || IsSpreadingUW || IsGroupFac || IsGroupUWFac || IsUW`.
- `MarineCargoDtl`: per medan.

| Kelompok | Medan (label persis) | Properti | Kontrol / sumber | Catatan |
| --- | --- | --- | --- | --- |
| Kepala | Top Risk | `.IsTopRisk` | cekbox | refresh grid |
| Kiri | Trading *, Goods * | `.TradingID`, `.GoodID` | dropdown "Choose": `BrowseTrading_RD` (Int-TRADING), `BrowseGoodsType` (Int-GOODSTYPE), nilai `.ID` / tampil `.Description` | ubah lalu `ProtectObjectMarine` mengisi `.TradingNote` / `.GoodNote` |
| Kanan | Conveyance *, Packing * | `.ConveyanceID`, `.PackingID` | `BrowseConveyance_RD` (Int-CONVEYANCE), `BrowsePacking_RD` (Int-PACKING) | mengisi `.ConveyanceName` + `.ConveyanceNote` / `.PackingNote` |
| Survey Agent | tombol Select Agent; Full Name, Address, City/Town, Country Code (RO) | `.PolicyData.SurveyAgent.*` | popup `SelectAgent` di atas `BrowseLLOYDAGENT_RD` (Int-lloydagent) | Agent ID tersembunyi (`1=2`) |
| LC | LC Status (cekbox); LC Bank, LC Number, LC Condition, LC Remark, LC Expired Date | `.PolicyData.LC.*` | teks / tanggal | RO `IsUW \|\| IsSpreadingUW` |
| Ship Data | tombol Select Ship; Ship Name, Ship Type, GRT, NRT, DWT, Age, Flag Name, Construction, Year Make 1, Remark (selalu RO) | `.PolicyData.Ship.*` | popup `SelectShip` di atas `BrowseShip_RD` (Int-SHIP); `SelectShip_act` menyalin 17 medan, GRT/NRT/DWT `replaceAll(',', '.')` | tombol "Clear (For Test)" tersembunyi |
| Kanan bawah | Sail Date; From Rute *; To Rute *; No BL; Invoice Number; Invoice Date; Max Age ID; Max Age Detail (RO) | `.PolicyData.SailDate`, `.FromRute`, `.ToRute`, `.PolicyData.BLNumber/InvoiceNumber/InvoiceDate/MaxAge/MaxAgeDetail` | Max Age: autocomplete `BrowseMaxAge_RD` (Int-MAXAGE) mengisi Detail dari `.Warranty` | Sail Date berubah: `CountAgeShip` (`AGE = tahun SailDate − Y_MAKE1`) |
| Additional Ship | grid Register ID / Ship Name / Voyage No, tombol Add / Delete | `.PolicyData.Ship.AdditionalShip()` | autocomplete `BrowseShip_RD` | Add tampil `!IsUW` |

`*` = wajib (`pyRequiredNew always`). Sampel `D:\migrasi\RNM\DDL\CONTOH\NB-18*(MARINE CARGO).xml` (struktur saja, nilai
tidak disalin): 8, 1, 1, 1, dan 1 baris CargoList. NB-181622 identik byte dengan NB-184183. LC hanya berisi `LC` (bool);
SurveyAgent, InvoiceDate, MaxAge, dan AdditionalShip tidak pernah terisi di sampel.

### S3. Tab Coverage `[terverifikasi]`

- **Rantai layar:**
  - grid cargo (`WORK!INPUTDTLCOVERAGE_FACIN` S131 `IsMarineCargo`: No., Trading, Goods, Packing, Conveyance), dengan tombol
    "Copy Deductible To All" di baris 1 (`CopyDeductibleToAllCargo_Act` **tidak ada di korpus**);
  - expand pane flow action `DATA-CARGO!INPUTDTLCOVERAGE_FACIN`, dengan pre-activity `AppendCoverageMarineCargo_ACT`. Bila
    CoverageList kosong, satu coverage ditambahkan dengan `ProRatePercent = min(100, hari(Start..End)/365×100)`;
  - grid coverage (Coverage Note, Coverage Initial; ikon tambah/hapus **tersembunyi**, jadi satu coverage per cargo);
  - grid `TotalTSIPremiGrossList` (Currency, Grand Total Gross TSI / Premium / Rate);
  - expand pane flow action `InputDtlCoverage_Marine`, yang memuat section `InputCoverageCargo_FacIn`.
- **Form `InputCoverageCargo_FacIn`:**
  - Keterangan Jaminan (RO);
  - Coverage Initial;
  - tombol **Select Coverage** (popup `SelectCoverage`, `BrowseMarineCondition_RD` dengan `Type = "MARINE"`, kolom Coverage
    ID / Coverage Note; dari master polis bila `PolicyMasterIDPega` terisi);
  - Name (mata uang, `BrowseCurrency_RD`);
  - Rate (%);
  - LimitofLiability;
  - Diskon (%);
  - TSI * ;
  - Premi;
  - Min Premi;
  - Diskon (RO);
  - grid **Warranty Extension** `.AdditionalCoverage` (Coverage autocomplete `BrowseMarineCondition_RD`, Deductible RO, Add /
    Delete);
  - grid **Deductible** (editor `addDeductibleCargo`, sama polanya dengan deductible FIRE / Aneka) + Deskripsi deductible.
- **Rumus** (`NB\Activity\CountCoverageMarine.xml` L998–1073):
  - `Premium = TSI×Rate/100` (4 desimal);
  - `Discount = Premium×DiscountPercentage/100`;
  - `TSILiability = TSI`;
  - MOP master (`PolicyType==1 && IsMOP=="MOP"`) membuat `Premium = 0`;
  - Premi diketik (`CountTSIandRateMarineCargo_Act` L245) menghitung `Rate = Premium×100/TSI`;
  - CurrencyList per mata uang = Σ Premium / Σ TSI.
- **Bagian RNM** (`CountPremiAndTSIRNMCargo_ACT` langkah 7):
  - `TSINusantaraRe = PS×TSI/100`;
  - `PremiNusantaraRe = PS×Premium/100 − PS×Discount/100` (20 desimal);
  - total TSI RNM = PS × Σ `TSILiability` / 100.
- **Tahap UW** (`_IsUW`): semua medan read-only, tanpa Select Coverage / Get Currency Master, dan tanpa Add/Delete pada
  Warranty Extension.

### S4. Tab Spreading `[terverifikasi]`

- **Rantai layar:**
  - grid cargo `CoverageSpreadingList` (`IsMarineCargo` L14765: No., Trading, Goods, Packing, Conveyance);
  - flow action `InputSpreadingMarineCargo_FacIn`, yang memuat `InputCoverageSpreadingMarineCargo` (grid coverage: Coverage,
    Standard Rate (%), Gross Premium, Premium Nusantara Re);
  - **`SpreadingItem` + `InputDtlSpreadingCoverage_FacIn`, sama dengan FIRE** (TSI 100%, TSI Liability, TSI Nusantara Re,
    Gross Premium, Premium Nusantara Re; baris Type Treaty, % Share, TSI 100% Spreaded, TSI Spreaded (RNM), Limit Of
    Liability, Premium Spreaded).
- **Tanpa spreading otomatis berbasis kapasitas treaty.** `GetKapasitasTreaty` tidak punya cabang Marine; grep
  `IsMarineCargo` hasilnya 0. Satu-satunya baris bawaan: `SetSpreadingCoverage_DT` langkah 8, treaty `10001` 100%, lewat tombol
  Spreading di form coverage. Tombol itu tersembunyi (S20 `1=2`).
- **Copy To All Spreading:** `CopyToAllSpreadingCargoMBU_ACT` langkah 2. `TSI = PS×TSI/100`, `Premium = PremiNusantaraRe`.
- **Hitung ulang** `SetTSIPremiSpreaded_FacIn` langkah 12: blok FIRE `CoverageBasis!=5` tanpa level lokasi, ditambah empat
  langkah pembulatan.
- **Ringkasan per mata uang:** `SumTSIPremiSpreadedRNM_MARINECARGO_Act`.
  - LoL hanya dari cargo 1.
  - TSI dibagi top risk / non top risk.
  - Teks `"Total Limit of Liability(" + pct + "%)"`.

### S5. Cedant Panels, Premium Deduction, Payment `[terverifikasi]`

- **Cedant Panels:** tanpa cabang Marine. Tabel total memakai tata letak `!IsFire` (sama dengan Aneka).
- **Premium Deduction:** `CoverageCommisionList` cabang Marine (grid cargo, lalu `InputCommissionMC_FacIn`, lalu
  `CommisionItemAneka`) bersifat read-only. Medan potongan sama dengan FIRE.
- **Payment:** `PaymentCurrencyList`, tanpa cabang Marine. Angkanya datang dari `CountPremiAndTSIRNMCargo_ACT`.

Artinya ketiga tab ini dapat memakai komponen FIRE/Aneka yang sudah ada; yang perlu ditambahkan hanya sumber data cargo.
`PremiPaymentMarine` hanya dipanggil jalur unggah CSV anggota polis (`InputDtlPayment_PreAct`), **bukan** tab Payment FacIn.

### S6. Proteksi Submit `[terverifikasi]`

| Pesan (verbatim) | Syarat | Bukti |
| --- | --- | --- |
| `" Ship data CargoList " + IdxCargo + " can't  be empty!"` (dua spasi) | per cargo: `Ship.ID==""` atau `NM_SHIP==""`; hanya bila `IsMOP!="MOP"` | `ProtectShipData_Act` L676–1074; `Protection_Act` 27.1 L4948 |
| `"Multi Location Must Have Toprisk"` | CargoList > 1 baris dan tidak ada `IsTopRisk==true` | `Protection_Act` L5013–5519 |
| `" Coverage "+ IndexCargo +" can't be empty!"` | `.Coverage==""` atau `.CoverageNote==""` | `ProtectCoverage_Act` L1228–1525 |
| `"Currency on object " + IndexCargo + " coverage " + n + " must be choose!"` | `Currency.ID==""` atau `OldID==""` | `ProtectCoverage_Act` L1294–1673 |
| `"% Total spreading on cargo " + i + " coverage " + j + " must be 100%!"` | Σ share ≠ 100; bukan MOP | `CheckSpreadingProtectMCargoMBU_ACT` L1586–1987 |
| `"Total premium spreaded on cargo " + i + " coverage " + j + " is not equals premium RNM!"` | Σ PremiumSpreaded ≠ PremiNusantaraRe (2 desimal) | idem L2027–2158 |

Kedua pesan spreading **dilewati** bila `OperatorID == <literal satu operator>` atau `IsGroup == "Group"`. Cabang itu
memakai pesan lain: "Total premium/TSI spreaded isn't match with … RNM!" (K74-5).

### S7. Data

- **Rancangan loader (`skema_gen.go`):**
  - `T_CARGOLIST` (22 kolom, induk `T_GENERAL_POLIS`);
  - `T_SHIP` (19 kolom; induk `T_CARGOLIST`);
  - `T_LC` (hanya `LC`).

  Coverage memakai `T_COVERAGELIST` (induk `T_CARGOLIST`), dengan `T_SPREADINGLIST`, `T_DEDUCTIBLELIST`, dan `T_CURRENCY` di
  bawahnya (`skema_gen.go` L1499–1563).
- **Kurang dari rancangan**, karena rancangan dibangkitkan dari sampel yang tidak memuatnya:
  - InvoiceDate, MaxAge, MaxAgeDetail;
  - sub-medan LC (LCBANK, LCNumber, LCCondition, LCRemark, LCEndDate);
  - SurveyAgent (ASMAgentID, FullName, Address, City, Country);
  - grid AdditionalShip (ID, NM_SHIP, VOYAGE_NO, …);
  - medan kapal tambahan dari `SelectShip_act` (FormerShipName, Classification, Rebuilt);
  - di coverage: CoverageInitial, CurrencyMaster, IndexCoveregeMaster, `.AdditionalCoverage`. Sesi c3 perlu memeriksa
    kolom mana yang sudah ada di T_COVERAGELIST.
- **Lookup**, dari kelas Pega, karena nama tabel Oracle tidak ada di korpus. **Q-DBA:** kelas Int-TRADING, Int-GOODSTYPE,
  Int-CONVEYANCE, Int-PACKING, Int-TYPESHIP, Int-MAXAGE, Int-LLOYDAGENT, Int-MARINECONDITION / Int-CONDITION. `SHIP` sudah
  ada di `D:\migrasi\RNM\DDL\SHIP.txt`.
- **Hanya-baca:** `POOLDATA.MASTERCARGO`, `POOLDATA.json_polis` (`RDBList\ViewMasterCargo*.xml`).

### S8. Rule yang hilang / cacat di korpus `[terverifikasi]`

- **Tidak ada:** `CopyDeductibleToAllCargo_Act`; penyetel `Index.CARI1` / `Index.CARI10` yang dibaca `SelectCoverageMarine_Act`;
  `SetEmptyShip_DT`. `BrowseMarineCondition_RD` kelas Int-CONDITION hanya ada di folder Endorsment Fac In.
- **Cacat:**
  - `CalcultePersentageSpeading_Act_FacIn` langkah 8 memutar `pyWorkPage.CargoList`, bukan `.OfferFacIn.CargoList`.
  - `CountGrossPremiEDM_Act` menghitung Premium dari TSILiability sebelum TSILiability direset.
  - `SetSurveyAgent_Act` menulis Survey Agent di level kasus, bukan per cargo.

## Verifikasi ulang sesi 0f `[terverifikasi]`

Diperiksa langsung, 10-10-2026:

- `When\IsMarineCargo.xml` L44 `G OR H`;
- `InputDtlObject_FacIn.xml` L6260 `pyEditAction InputDtlCargo_FacIn` + grid CargoList L2849–2852;
- `Protection_Act.xml` memuat "Multi Location Must Have Toprisk" (1×);
- `ProtectShipData_Act.xml` memuat teks pesan ship data;
- `CountCoverageMarine.xml` L998–1002 `@Math.divide((.TSI*.Rate),100,4)`;
- `FlowAction\InputDtlCoverage_FacIn.xml` L26 pre-activity `AppendCoverageMarineCargo_ACT`;
- `SetSpreadingCoverage_DT.xml` memuat treaty `"10001"`;
- `InputInwardFacultativeDtl.xml` L17833 Clauses `IsOilGas || IsFire`.

## Rencana bertahap (usulan agent)

| Fase | Isi | Backend (c3) | Frontend (0f) |
| --- | --- | --- | --- |
| P0 | gerbang `lini.marineCargo` + `IsMarineCargo` di GET kasus | predikat lini | kasus Marine Cargo membuka tab Marine Cargo |
| P1 | Object: grid + form cargo + popup Select Ship / Select Agent + lookup | migrasi 256 (K74-1) + GET/PUT objek cargo | `TabObjectCargo` |
| P2 | Coverage: form coverage + Select Coverage + Warranty Extension + deductible + rumus S3 | coverage cargo + hitung | `TabCoverageCargo` |
| P3 | Spreading: rantai S4, memakai ulang daun FIRE | spreading cargo | `TabSpreadingCargo` |
| P4 | Cedant / Premium Deduction / Payment: sumber data cargo | ringkasan cargo | memakai ulang komponen FIRE/Aneka |
| P5 | Submit: proteksi S6 + simpan Policy | proteksi + polis | `bolehSubmit` untuk Marine Cargo |
| P6 | CSV objek Marine Cargo (`InsertUploadMarineCargo_act`) | unggah | tombol CSV |

## Keputusan work owner 10-10-2026

Seluruh usulan agent dipilih:

- **K74-1 — migrasi 256 boleh ditulis.** Work owner / DBA yang menjalankan.
- **K74-2 — P0–P2 dulu**, lalu P3–P5, dan P6 (CSV) terakhir.
- **K74-3 — MOP / Declaration Policy menjadi fase terpisah sesudah P5.** P1–P5 hanya polis individual.
- **K74-4 — Select Coverage** mengisi cargo/coverage yang sedang dibuka.
- **K74-6 — Survey Agent** disimpan per cargo.
- **K74-5:** cabang operator mengikuti A14 (tidak diport); cabang polis master literal tidak diport.

## Keputusan yang diminta (K74-1…K74-6)

- **K74-1 — migrasi 256** (rentang nbfacin 255–259): `T_CARGOLIST` / `T_SHIP` / `T_LC` dari rancangan, ditambah kolom yang
  kurang (S7) dan tabel anak Additional Ship + Survey Agent. Butuh persetujuan migrasi (CLAUDE.md §10).
- **K74-2 — urutan fase:** P0–P2 dulu (kasus bisa diisi objek dan coverage), lalu P3–P5, dan P6 terakhir.
- **K74-3 — Master / Declaration Policy (MOP):**
  - PolicyType Individual / Master / Declaration (`SetPolicyType_Act`);
  - View Policy Master Cargo;
  - coverage diambil dari master polis (`POOLDATA.json_polis`);
  - `CekLimitSpreading` metode 1–3 untuk Declaration.

  Usulan agent: fase terpisah sesudah P5. P1–P5 hanya polis individual.
- **K74-4 — Select Coverage:** indeks `Index.CARI1` / `CARI10` tidak ada di korpus. Usulan: tiru maksudnya, yaitu coverage
  terpilih masuk ke cargo/coverage yang sedang dibuka.
- **K74-5 — literal di korpus:**
  - cabang proteksi spreading untuk **satu operator ber-ID literal**;
  - satu **nomor polis master literal** (LimitofLiability / Get Currency Master `"USD"`);
  - satu **nama grup literal** (`IsSpecialCase`).

  Usulan: cabang polis master literal tidak diport (pola tiket 21).

  **Cabang operator sudah diputuskan:** `KEPUTUSAN-30-09-2026.md` butir A14 menolak gerbang login operator ber-literal, dan
  nilainya tidak ditulis; portnya kelak lewat model peran (butir 33). Untuk proteksi spreading S6, hanya cabang
  `IsGroup == "Group"` yang dipertimbangkan, menunggu model peran.
- **K74-6 — Survey Agent:** layar memasangnya per cargo, tetapi `SetSurveyAgent_Act` menulisnya di level kasus. Usulan:
  per cargo (mengikuti layar).

## Frontend P0–P1 — tab Object Marine Cargo (sesi 0f, 10-10-2026)

**Status:** dibangun, uji hijau. Menunggu backend `objek-cargo` + lookup (sesi c3), dan nama tabel lookup dari DEV.

- **Gerbang** (`pages/InwardFacultative.tsx`):
  - `kasusCargo = !kasusFire && !kasusAneka && lini?.marineCargo === true`;
  - tab Object memasang `TabObjectCargo`; tab lain tetap "belum tersedia" sampai P2+;
  - FIRE dan Aneka tidak berubah.
- **`components/TabObjectCargo.tsx`:**
  - grid CargoList (S1);
  - form cargo dua kolom mengikuti `MarineCargoDtl`: kiri Survey Agent / LC / Ship Data, kanan medan pelayaran / Additional
    Ship;
  - popup Select Ship / Select Agent;
  - dropdown master Trading / Goods / Conveyance / Packing. Note diisi seperti `ProtectObjectMarine`
    (`ConveyanceName = ConveyanceNote = Description`, `Activity\ProtectObjectMarine.xml` L1470–1517);
  - Max Age autocomplete mengisi Detail.
- **`components/cargo.ts`** (fungsi murni): baris baru, perapian jawaban, `terapkanKapal` (GRT/NRT/DWT koma → titik,
  `SelectShip_act` L665/L706/L727), `terapkanAgent`, `terapkanKapalTambahan`.
- **Label:** `pyLabelFieldValue` / `pyLabelPreview`, misalnya "Goods". `pyLabelFor` ("Barang", "Pengakutan") hanya sasaran label;
  diverifikasi di `InputDtlCargo_FacIn.xml` sel Goods L1886–2207.
- **Kabel tanggal `DD-MM-YYYY`** (ralat kontrak, konvensi nbfacin; c3 setuju).
- **Penyimpangan sadar:**
  1. Popup Select Ship / Select Agent memakai satu kotak cari (Pega: filter per kolom grid).
  2. Top Risk di tahap UW ikut baca-saja. Pega `InputDtlCargo_FacIn_IsUW.xml` L500 membiarkannya dapat diubah, tetapi sistem
     baru membuat case baca-saja di UW kecuali Spreading (tiket 59).
- **Uji:** `components/TabObjectCargo.test.tsx` (label, gerbang, fungsi murni, urutan render). Tampilan dicek di Chrome
  headless.

## Backend P0–P2 — dibangun (sesi c3, 10-10-2026)

**Status:** dibangun, uji hijau. Masih menunggu:
- nama tabel lookup dari DEV (Q-DBA): master catatan dan endpoint lookup belum terpasang;
- cek bahwa 256 belum terpakai di DEV (dikerjakan work owner).

- **P0:** `lini.marineCargo` di GET kasus = When IsMarineCargo.
  - `kasusBisnis` kini mengenal jalur relatif `.OfferFacIn.QuotationData.BusinessType`, sebab When berkelas WORK.
  - Uji `TestLiniKasus` "marine cargo" (GroupPanel 004).
- **Migrasi 256** (`migrations/256_marine_cargo.sql` + `_down`; DITULIS, tidak dijalankan):
  - tabel baru T_CARGOLIST (rancangan + INVOICE_DATE, MAX_AGE, MAX_AGE_DETAIL, FLAG_DELETE; BL_NUMBER 4000);
  - T_SHIP (lebar POOLDATA.SHIP + FORMER_SHIP_NAME / CLASSIFICATION / REBUILT);
  - T_LC (+ LCBANK / LC_NUMBER / LC_CONDITION / LC_REMARK / LC_END_DATE);
  - T_SURVEYAGENT (BARU, per cargo K74-6), T_ADDITIONALSHIP (BARU, induk T_SHIP), T_ADDITIONALCOVERAGE (sebagian rancangan);
  - kolom T_COVERAGELIST.COVERAGE_INITIAL dan T_DEDUCTIBLELIST.DESCRIPTIONS.
  - Kolom dicatat di `docs/STRUKTUR-TABEL-NB-FACIN.md`; penjaga `TestKolomDDLCocokDenganStruktur` hijau.
- **Endpoint:**
  - `GET` / `PUT /api/nbfacin/kasus/{caseId}/objek-cargo` → `{cargo: ObjekCargo[]}`;
  - `POST /api/nbfacin/kasus/{caseId}/hitung-coverage-cargo` → `{baris: CoverageCargo[]}`;
  - `GET /api/nbfacin/kapal` diperluas ADITIF: shipTypeId, shipTypeName (= SHIP.SHIPNAME, `SelectShip.xml` L6439), register, imo,
    yMake2, remark, formerShipName, rebuilt.
- **Beda dari kontrak 0f:**
  1. `deskripsiDeductible` menjadi `descriptions` PER BARIS deductible. Sebabnya, `.Descriptions` terikat ke grid
     `.DeductibleList` kedua (`InputCoverageCargo_FacIn.xml` L16997 / L17652), bukan ke coverage.
  2. `additionalShip` ada di tingkat cargo seperti kontrak; di model ia anak `ship` (`.PolicyData.Ship.AdditionalShip`).
- **Medan server dari server; badan PUT diabaikan untuk medan ini:**
  - catatan master (ProtectObjectMarine): ConveyanceName = ConveyanceNote = Description;
  - maxAgeDetail = `.Warranty`;
  - ship.age:
    - CountAgeShip `@substring(SailDate,0,4) − Y_MAKE1`;
    - Y_MAKE1 kosong atau tak terurai → nilai tersimpan;
    - dihitung setiap PUT `[keputusan agent]`;
  - flagDelete / flagOldData: nilai tersimpan;
  - currencyId / OldId (SetCurrencyCoverage_Act);
  - discount = r4(Premium×Disc%/100), tsiLiability = TSI;
  - tsiNusantaraRe / premiNusantaraRe (PS × …, 20 desimal);
  - proRatePercent = AppendCoverageMarineCargo_ACT min(100, hari/365 (8 desimal) ×100):
    - disintesis di GET untuk cargo tanpa coverage, TIDAK disimpan sampai PUT;
    - coverage lama di posisi yang sama mempertahankan nilainya;
  - totalGross = SumTotalTSIPremiGross_Act langkah 4: ×100 untuk Marine `[terverifikasi]` L7154, ×1000 hanya FIRE.
  - Angka hitung dan total dikirim dalam bentuk tersimpan (8 desimal) `[keputusan agent]`.
- **Pertahankan anak:**
  - cargo ber-`id` lama dengan `coverages` absen / null → coverage lama (beserta spreading) dipertahankan;
  - `coverages` dikirim → spreading coverage lama di posisi yang sama dipertahankan (untuk P3).
- **GRT / NRT / DWT kapal utama:** `replaceAll(',', '.')` (SelectShip_act). Disimpan TEKS seperti T_MARINEHULL (214), sebab
  perilaku `@toDecimal` atas teks tak terurai `belum terverifikasi`.
- **Wajib isi** (`pyRequiredNew always`): trading, goods, conveyance, packing, fromRute, toRute → 400.
- **Status galat:**
  - bukan MARINE CARGO → 409;
  - master belum terpasang dan ID master berubah / baru → 503 (ID yang sama mempertahankan catatan lama);
  - ID master tidak ada → 400.
- **Penulisan:** hapus pohon cargo case, lalu sisip PER BARIS bind skalar. Setiap baris lewat `periksaBindBaris`.
  - Pemeriksa itu menemukan `T_CARGOLIST.PARENT_ID` = ID case VARCHAR2(32); ditambahkan ke `indukTeks` bersama T_LOCATIONLIST.
- **Uji:**
  - repository: `TestKolomCargoAdaDiMigrasi` (dua arah, kolom kode ↔ seluruh migrasi), `TestSQLCargo`, `TestTulisCargo`;
  - services: `TestBacaObjekCargo`, `TestGantiObjekCargo`, `TestGantiObjekCargoPertahankan`, `TestHitungCoverageCargo`,
    `TestUmurKapalDanProRate`;
  - handlers: `TestObjekCargo` (kunci kabel persis), `TestHitungCoverageCargoKabel`, `TestLookupAneka` (kapal aditif);
  - mutasi 19: 18 mati, 1 setara (handler meneruskan proRatePercent klien; layanan mengabaikannya).

## Frontend P2 — tab Coverage Marine Cargo (sesi 0f, 10-10-2026)

**Status:** dibangun, uji hijau. Memakai bentuk final c3: `totalGross` per cargo dan `descriptions` per deductible.
`deskripsiDeductible` tidak dipakai.

- **`components/TabCoverageCargo.tsx`:**
  - grid cargo (No., Trading, Goods, Packing, Conveyance);
  - per cargo: grid coverage (Coverage Note, Coverage Initial), dengan form coverage di bawahnya, lalu grid total
    (Currency, Grand Total Gross TSI / Premium / Rate);
  - form `InputCoverageCargo_FacIn`: Coverage Note (baca-saja), Coverage Initial, Select Coverage (popup), Currency,
    Rate (%), Limit of Liability, Discount (%), TSI *, Premium, Min Premium, Discount (baca-saja);
  - Warranty Extension (autocomplete kondisi marine);
  - Deductible: `GridDeductible` FIRE + Deskripsi.
- **Simpan** lewat `PUT …/objek-cargo` (cargo utuh). **Hitung** lewat `POST …/hitung-coverage-cargo` saat isian selesai:
  - pemicu currency / rate / tsi / diskonPersen / premium;
  - jawaban usang dibuang;
  - Warranty Extension + Deductible di layar dipertahankan.
- **Label:** `pyLabelPreview` bila ada, kalau tidak `pyLabelFieldValue`; diurai per sel dengan parser XML. Hasilnya
  "Coverage Note", "Currency", "Premium", "Discount (%)", bukan `pyLabelFor` "Keterangan Jaminan" / "Name" / "Premi" /
  "Diskon (%)". Ini konsisten dengan form Aneka yang sudah disetujui. Halaman alat tiket 21 (`CoverageCargo.tsx`) tidak
  diubah.
- **Deskripsi deductible:** prop opsional `deskripsi` di `GridDeductible` (aditif). FIRE / Aneka tidak mengirimnya, jadi tidak
  berubah; ini dijaga uji. Labelnya "Deskripsi" dari `InputDeductible_FacIn.xml` S3, karena sel grid cargo hanya berlabel generik
  "Text Area".
- **Tidak diport:** "Copy Deductible To All" (`CopyDeductibleToAllCargo_Act` tidak ada di korpus); cabang polis master literal
  (K74-5); MOP (K74-3).
- **Menunggu backend:** `GET /api/nbfacin/marine/coverage` (lookup kondisi marine) — nama tabel DEV belum ada. Sampai terpasang,
  popup Select Coverage dan autocomplete Warranty Extension akan gagal memuat.

## Frontend P3–P5 — Spreading, Cedant, Premium Deduction, Payment, Submit (sesi 0f, 10-10-2026)

**Status:** dibangun, uji hijau. Perintah work owner 10-10-2026: "LANJUT SAJA KE TAB BERIKUTNYA SAMPAI AKHIR". Kontrak P3–P5 sudah
dikirim ke c3; backend P3–P5 belum ada, jadi tab-tab ini belum dapat dicoba.

- **`components/TabSpreadingCargo.tsx`** (kabel `…/spreading-cargo`, `…/hitung-share`, `…/salin`, `…/hitung-baris`, PUT):
  - grid cargo `CoverageSpreadingList` → coverage `InputCoverageSpreadingMarineCargo` (Coverage / Standard Rate (%) / Gross
    Premium / Premium Nusantara Re; label pyValue diverifikasi) → detail `SpreadingItem` FIRE (5 angka) + baris spreading FIRE;
  - % Share RNM;
  - Copy Spreading hanya UW;
  - ringkasan `RingkasanSpreading`.
  - Tanpa spreading otomatis dan tanpa salin per lokasi (Marine tidak punya cabangnya).
  - `components/spreadingCargo.ts`: fungsi murni pohon cargo.
- **`TabCedant` / `TabDeduction`:** prop `cargo` (aditif, pola prop `aneka`):
  - ringkasan dari `…/spreading-cargo`;
  - Cedant tata letak `!IsFire`;
  - `DaftarKomisiCargo` baca-saja: cargo → coverage (Coverage / TSI / Gross Premium; `InputCommissionMC_FacIn`, tanpa kolom Old
    K73-39) → detail `CommisionItemAneka` + baris spreading.

  FIRE / Aneka tidak berubah; dua uji teks-sumber lama disesuaikan dengan ekspresi baru, perilaku yang diuji sama.
- **Payment:** `TabPayment` apa adanya; backend memilih pohon cargo.
- **Submit:** `bolehSubmit = kasusFire || kasusAneka || kasusCargo`. Proteksi Marine (S6) di backend.
- **Uji:** `components/TabSpreadingCargo.test.tsx`.

## Backend P3–P5 — dibangun (sesi c3, 10-10-2026)

**Status:** P3, P4, dan proteksi + jenjang P5 sudah dibangun. Uji hijau (`go vet`, `go test -short ./...`). Belum dicoba di DEV.
**Yang menunggu:** ~~FACINOFFER / FACINPRODUCTION Marine~~ — dipasang 11-10-2026 (lihat "FACINOFFER / FACINPRODUCTION Marine" di
bawah).

**P3 — Spreading** (`services/spreadingcargo.go`, `handlers/spreadingcargo.go`). Endpoint yang dipakai frontend:

- `GET` / `PUT …/spreading-cargo`;
- `POST …/spreading-cargo/hitung-share`, `…/salin`, `…/hitung-baris`.

Pada Save, persen dihitung ulang — penyimpangan dari Pega yang sudah diputuskan (K74-7). Total mengikuti
`SumTSIPremiSpreadedRNM_MARINECARGO_Act`:

- LoL hanya dari cargo 1;
- ClaimSpreaded hanya dari cargo top risk;
- bila share antarbaris berbeda, totalnya kosong;
- label CARI15.

**P4 — Cedant / Premium Deduction / Payment.** Untuk case MARINE CARGO tanpa item FIRE dan tanpa objek Aneka, sumber ringkasan
spreading adalah pohon cargo (`spreadingRingkasCargo`). FIRE dan Aneka tidak berubah.

**P5 — Submit.** Implementasi di `services/proteksicargo.go`, dengan wiring di `proteksi.go`, `submituw.go`, dan `submitpolicy.go`.
Uji di `services/proteksicargo_test.go`; mutasi di wiring baru semuanya mati.

- **Proteksi S6** `[terverifikasi]` (port dari kepala berkas `proteksicargo.go`):
  - Coverage dan Currency: `ProtectCoverage_Act` langkah 4 MARINE, flag `ProtectMarine.CARI2`. Coverage kosong tidak diproteksi
    (K74-12).
  - Kapal: `ProtectShipData_Act`, flag `ProtectMarine.CARI1`. Pesannya verbatim dengan dua spasi.
  - Toprisk: "Multi Location Must Have Toprisk" muncul bila ada lebih dari satu cargo dan tidak satu pun top risk.
  - Spreading 100 % dan premi RNM (`CheckSpreadingProtectMCargoMBU_ACT`), dijalankan di tahap uw. Σshare tidak dibulatkan; premi
    dibandingkan pada 2 desimal (K74-13). Cabang IsGroup dan operator literal tidak diport.
- **Langkah IsFire tidak berjalan untuk case cargo:** Scoring Risk, mata uang TSI FIRE (langkah 21), dan proteksi lokasi FIRE.
  Deductible Acceptance TL juga tidak menahan, karena `addDeductibleCargo_ACT` L248 men-set `IsBedaCurr = 0`.
- **Jenjang limit** (K74-8): tabel `M_LIMIT_NONPROPANDENGG`, termasuk di ActFlow Senior lewat When COB `kb` case. TSI mengikuti
  `CountTotalTSIPremiNusaRe_Act` langkah 10 dan 14:
  - per cargo × coverage, `.TSI` × kurs mata uang coverage × PS;
  - cargo top risk memakai `.TSILiability`.
  - `[keputusan agent]` Cargo `FlagDelete "1"` dilewati.
- **IsFacRetro** `[keputusan agent]`: benar bila ada baris spreading FACOUT di pohon cargo (pola FIRE / Aneka).
- **IsFlagUW:** MARINE CARGO = 1, kecuali polis deklarasi (PolicyType 2) = 0. Jalurnya sama dengan Aneka.

**FACINOFFER / FACINPRODUCTION Marine — dipasang (sesi c3, 11-10-2026)** (K74-9/10 "ikut Pega apa adanya", K74-11):

- **Gerbang 409 `ErrProduksiCargoBelumTersedia` DICABUT.** Binding dan selesai Policy case MARINE CARGO kini menulis baris cargo.
  Kolom Marine terbukti di DDL (`DDL\FACINOFFER.txt` L17-50, `DDL\FACINPRODUCTION.txt` L18-63) `[terverifikasi]`.
- **FACINOFFER** = `SaveOfferProduction_Act` langkah 16 (`services/produksicargo.go` `barisFacInOfferCargo`). Satu baris per cargo ×
  coverage × baris spreading. OBJECTNO = PackingID, OBJECT_NAME = PackingNote, OBJECT_ITEM = GoodNote, TRADING_NAME = TradingNote,
  CONVEYANCE_NAME = ConveyanceNote, IDSHIP = Ship.ID, FROMRUTE / TORUTE, SAILDATE. Alamat / okupasi / TOP_RISK = NULL.
- **FACINPRODUCTION** = `SaveTreatyProduction_Act` langkah 22 (`barisFacInProductionCargo`): kolom yang sama + RICOMM / BROKERAGE,
  LOL = ClaimEstimation (kosong 0). MOID case didahulukan, OBJECT_NAME_ID / LAYERNO teks "NULL" (pola Aneka). IDX_OBJECT_ITEM /
  TOP_RISK / TSI100 NULL.
- **`M_RISK_LOSS_PROFILE`** = Aneka (`profilRisikoLossAneka`, K74-11 = K73-48).
- **SAILDATE** (K74-9b): `@FormatDateTime(SailDate, "dd/MM/yyyy hh:mm:ss", Asia/Jakarta)` apa adanya — `hh` 12-jam, jadi 13:00 WIB
  tertulis 01:00:00; 12:00 WIB tertulis 12:00:00. Ditulis `TO_DATE(?, 'DD/MM/YYYY HH24:MI:SS')`. Kosong = NULL.
- **IDSHIP VARCHAR2(20)** sedangkan `POOLDATA.SHIP.ID` VARCHAR2(100): ID kapal lebih dari 20 karakter (atau 20 bita) → **409
  `ErrIDKapalFacIn`** bernomor cargo. Tidak dipotong diam-diam.
- **K74-22 DEDUCTION2_MENJADI / _SELISIH = NULL** — **keputusan work owner 11-10-2026** (sama dengan FIRE / Aneka, K73-45). Pega mengisi
  `PremiumSpreaded × Payment.PctDeduction2 / 100`, tetapi `PctDeduction2` tidak tersimpan di sistem baru. Tab Premium Deduction
  menyimpan Deduction2 per mata uang. NULL ini sama dengan FIRE / Aneka (K73-45).
- **Catatan korpus `[terverifikasi]`:** langkah INSERT di SETIAP cabang `SaveOfferProduction_Act` / `SaveTreatyProduction_Act`
  (FIRE juga) berprakondisi `.pxListSubscript>1` (T=2 lanjut / F=3 lewati).
  - **Keputusan work owner 11-10-2026 (K74-23):** TETAP tulis SEMUA baris spreading (K74-9: satu baris per spreading). Port FIRE /
    Aneka yang sudah berjalan tidak diubah; Marine sama. Catatan silang di tiket 63 (FACINOFFER FIRE), 69 (FACINPRODUCTION
    FIRE), 73 (Aneka) - tanpa perubahan kode.
- **Uji:** `services/produksicargo_test.go`: SAILDATE 12-jam, batas IDSHIP, kolom offer / produksi, jalur Binding / Policy +
  profil risiko, FIRE tidak masuk cabang cargo. `repository/binding_test.go` / `selesaipolicy_test.go`: FACINOFFER 46 bind,
  FACINPRODUCTION 67 bind, SAILDATE TO_DATE. Enam mutasi tangkap semua.

**Belum dibangun — menunggu:**

- **Lookup master Marine** dan implementasi Oracle `PembacaMasterCargo`: menunggu nama tabel di DEV. Sampai itu ada, endpoint
  hitung-coverage mengembalikan 503 `ErrMasterCargoBelumAda`.

## P6 — Upload CSV Marine Cargo (sesi 0f, 11-10-2026)

**Pemicu:** work owner 11-10-2026 meminta tombol CSV di case MARINE CARGO "perbaiki". Sejak bug DEV 11-10-2026, CSV FIRE ditolak
untuk non-FIRE (tiket 75).

**Korpus `[terverifikasi]`:**

- **Pemilihan impor:** `ChooseActivityMappingCSV_Act` L539 / L599 menjalankan `InsertUploadMarineCargo_act` bila `IsMarineCargo`
  (FIRE L439 / L499).
- **Kolom yang dibaca** (nama judul; Pega `pxUploadCSVResults` mengikat per NAMA): TRADING, DEPARTURE, DESTINATION, GOODS,
  GOODS_NAME, PACKING, CONVEYENCE, CONVEYENCE_NAME, COVERAGE_IDX, COVERAGE_NAME, TSI, CURRENCY, AGENT, DATASHIP, INVOICE, ETD,
  BLNUMBER, TOP_RISK (18).
- **Template:** biner `TemplateCSVMarineCargo` (`SetDownloadMasterCSV` L1012–1091), TIDAK ada di ekspor.
- **Proteksi baris:** `ProtectUploadMarineCargo_act`, pesan "* Trading / Goods / PACKING / Conveyence / Agent / Ship … Not Available
  in Master" dan "Jangan gunakan ID kapal dikolom DATASHIP pada CSV". Kapal yang tidak ada di master didaftarkan lewat procedure
  `POOLDATA.RDBMASTERSHIP` (isi di `D:\migrasi\RNM\DDL\RDBMASTERSHIP.txt`).
- **Coverage:** disalin dari polis master (langkah 3.12 L5042–5841, `PolicyMasterNumber`); tanpa master hanya TSI / Currency /
  IndexCoveregeMaster yang terisi.

**Keputusan work owner 11-10-2026:**

- **K74-14 Template:** work owner MENGIRIM template Pega; urutan kolom mengikuti berkas itu. Pembacaan tetap per NAMA judul.
- **K74-15 Coverage:** seperti Pega tanpa master. Setiap cargo mendapat 1 coverage berisi TSI + Currency dari CSV; kode coverage
  dipilih lewat Select Coverage. Penyalinan coverage dari polis master menyusul bersama MOP (K74-3).
- **K74-16 Kapal baru:** procedure `RDBMASTERSHIP` DIPORT ke Go, tanpa CALL (ADR-0043), menulis `POOLDATA.SHIP` seperti Pega
  (persetujuan tulis tabel Pega).
- **K74-17 Jalur:** tanpa batas baris, memakai infrastruktur tiket 75 (unggah ke server, latar, kemajuan, Ulangi).

## Backend P6 — dibangun (sesi c3, 11-10-2026)

**Status:** dibangun, uji hijau (`go vet`, `go test -short ./...`); mutasi aturan inti mati (dua yang hidup ekuivalen: cache
kapal, nilai kosong). Migrasi **258** DITULIS, belum dijalankan. Belum dicoba di DEV. Validasi master Trading / Goods / Packing /
Conveyance / Agent masih menunggu nama tabel DEV: tanpa `PembacaMasterCargo`, unggah cargo dijawab 503 `ErrMasterCargoBelumAda`
sebelum berkas dibaca.

- **Pemilihan importer** (`services/csvfire.go` `jenisCsvKasus`), persis `ChooseActivityMappingCSV_Act`:
  - IsFire → `fire` (tiket 66 / 75, tidak berubah);
  - IsMarineCargo → `cargo`;
  - selain itu → 409 dengan pesan gerbang yang sama.
  Endpoint `…/objek/csv/berkas`, `/{id}`, `/aktif`, `/{id}/baris`, `/simpan`, `/ulangi` tetap. Pekerjaan menyimpan `JENIS`;
  simpan / ulangi / jalur latar menolak pekerjaan yang jenisnya tidak sesuai Class of Business case sekarang (409). Save Data JSON
  lama (tiket 66) tetap FIRE saja.
- **Migrasi 258** (rentang 255–259):
  - `T_NB_CSV_UNGGAH.JENIS` (DEFAULT `'fire'`) dan `KOLOM` (judul dikenali urut berkas);
  - 16 kolom penampung cargo `T_NB_CSV_BARIS.C_<judul>`. CURRENCY dan TOP_RISK memakai kolom FIRE yang sama.
  - STRUKTUR sudah diperbarui (penjaga golongan tipe hijau).
- **Kabel (aditif):**
  - status + `jenis`;
  - `GET …/{id}/baris` + `kolom: string[]` (urut berkas; FIRE = urutan template), dengan baris dipetakan per kolom jenisnya.
- **Port** (`services/csvcargo.go`, bukti per langkah di kepala berkas):
  - Judul dibaca per nama.
  - Langkah 3.1: TSI, RMB→CNY, huruf besar rute / DATASHIP.
  - Protect, dengan urutan dan transisinya:
    - Trading / Goods / Conveyence tidak ada → berhenti;
    - PACKING → lanjut;
    - ID kapal → pesan diganti;
    - agen tidak ada → kosong.
  - Prakondisi 3.9: cargo baru hanya bila TRADING / GOODS / CONVEYENCE / PACKING berbeda dari baris sebelumnya.
  - 3.10–3.12 diterapkan ke cargo terakhir.
  - Coverage tanpa polis master: satu coverage TSI + Currency (K74-15).
  - Kapal tak ada di master → port RDBMASTERSHIP (`repository/kapalmaster.go`), INSERT `POOLDATA.SHIP` tanpa CALL (K74-16).
- **Simpan per kelompok tanpa batas** (K74-17):
  - satu kelompok ≥ 1000 baris dan berakhir di batas cargo, jadi cargo tidak pernah terbelah dan Ulangi selalu mulai di cargo baru;
  - kelompok pertama mengganti seluruh cargo case dan menghapus CurrencyList + cicilan;
  - tulisan array-bind berkelompok (`repository/cargocsv.go`, mesin tiket 75).
- **Kunci K75-2:** PUT `…/objek-cargo` dan Save spreading cargo 409 selama unggah berjalan.

**Keputusan** (K74-19..21 dijawab work owner 11-10-2026; K74-18 masih `[keputusan agent]`):

- **K74-18** `[keputusan agent]`, menunggu template Pega dari work owner: judul wajib memuat ke-18 nama. Kolom lain diabaikan.
- **K74-19** Keputusan work owner 11-10-2026 "ikut Pega": baris berurutan berkunci sama MENIMPA kapal / rute / agen / coverage cargo
  yang sama, sehingga TSI baris sebelumnya hilang.
- **K74-20** Keputusan work owner 11-10-2026 "setuju", penyimpangan sadar dari Pega:
  - TRADING / GOODS / PACKING / CONVEYENCE kosong = tidak ada di master (Pega memilih baris master pertama sembarang);
  - AGENT kosong = tanpa agen (Pega memakai agen baris sebelumnya karena halaman AgentTemp tidak dihapus);
  - kapal baru didaftarkan saat simpan (Pega: di Protect, juga untuk baris yang gagal).
- **K74-21** Keputusan work owner 11-10-2026 "setuju", mengikuti CSV FIRE:
  - mata uang tak ada di master = galat (K66 D6);
  - TSI / ETD wajib terurai (K66 D7);
  - CurrencyList dihapus tanpa dibangun ulang (K66 D5; Pega langkah 5 membangunnya dari coverage);
  - SailDate disimpan dalam bentuk tanggal kita.

Langkah Insert 3.3–3.8 dan Protect 5.4 / 7.4 ber-label `//` (mati) dan tidak diport. Lompatan `KPL` pada Protect 7.5.2 menunjuk
label yang tidak ada di activity, sehingga perilaku runtime-nya `belum terverifikasi`; port melanjutkan dengan membaca ulang kapal.

### Master Marine + template Pega (sesi c3, 11-10-2026)

- **DDL dikirim work owner:** view POOLDATA.TRADING / GOODSTYPE / PACKING (ID, DESCRIPTION, TYPE) dan CONVEYANCE (+ TYPECONVEYANCE),
  semuanya atas JSONDATA m_*. `repository/mastercargo.go` membacanya; filter mengikuti RD (`.ID =` / `.Description =`, tanpa
  saringan TYPE).
  - `GET /api/nbfacin/marine/{trading|goods|conveyance|packing}` → `{baris: [{id, description}]}`, urut DESCRIPTION, paling
    banyak 500.
  - Pembacaan yang sama dipakai untuk catatan PUT / GET objek-cargo dan IDMaster unggah CSV. Bila deskripsi ganda, ID terkecil
    yang diambil `[keputusan agent]`.
- **Masih menunggu DDL:** MaxAge, LLOYDAGENT, MARINECONDITION, TYPESHIP.
  - `/marine/max-age`, `/agent`, `/coverage` → 503.
  - MaxAge yang tidak berubah pada PUT objek-cargo tetap memakai catatan lama; MaxAge baru → 503.
  - Unggah CSV dengan AGENT terisi → validasi berhenti di baris itu (gagal_validasi, galat jelas) `[keputusan agent]`. AGENT
    kosong tetap sah.
- **K74-18 final** (template `templatecsvmarinecargo.csv`): judul wajib = 17 kolom template tanpa INVOICEDATE.
  - COVERAGE_NAME opsional.
  - INVOICEDATE diterima tetapi diabaikan: tidak dibaca Pega, tidak disimpan di penampung, dan tidak tampil di `kolom` View
    Upload.
  - Kolom lain juga diabaikan.
- **DDL kedua (11-10-2026)** `[terverifikasi DDL]`: view POOLDATA.LLOYDAGENT / MAXAGE / CONDITION.
  - `/marine/agent?q=`, `/marine/max-age?q=`, `/marine/coverage?q=` kini terisi.
  - PUT objek-cargo: MaxAgeDetail = WARRANTY master.
  - CSV P6: AGENT dicari ke master (`ID =`). Agen tidak ada → kosong TANPA pesan, karena Protect 5.4 berlabel `//` (mati)
    `[terverifikasi]`. Galat sementara "master agen belum tersedia" tidak lagi terjadi dengan master terpasang.
- **Rincian pencarian** (`repository/mastercargo.go`):
  - Agen (BrowseLLOYDAGENT_RD NB): `.Name Contains`, PEKA huruf (tanpa pyCaseInsensitive).
  - MaxAge (autocomplete MarineCargoDtl, RD NB tanpa filter, kolom cari .ID / .Warranty): memuat, tak peka huruf
    `[keputusan agent]` — mode cocok tidak tertulis di ekspor.
  - Kondisi marine: BrowseMarineCondition_RD **NB** = kelas Int-MARINECONDITION, TANPA filter (varian Endorsment = Int-CONDITION
    ber-filter Type "MARINE" / ConditionAdd = identitas lain, tidak dipinjam). Sumber = view CONDITION `TYPE = 'MARINE'` tanpa
    saringan CONDITIONADD; Int-MARINECONDITION = CONDITION ber-TYPE 'MARINE' — keputusan work owner 11-10-2026.
  AGENT CSV tak dikenal = kosong tanpa pesan (ikut Pega, Protect 5.4 mati) — keputusan work owner 11-10-2026.
- **TYPESHIP** tanpa DDL dan tanpa route.

### Audit lebar sesudah ORA-12899 — keputusan work owner 11-10-2026 (lewat sesi 9d; sesi c3)

- **(b) Tabel cargo 256, tanpa migrasi.** CSV cargo menolak dengan galat jelas per baris bila nilai MASTER melebihi kolom tujuan:
  - T_CARGOLIST.TRADING_ID / GOOD_ID / PACKING_ID / CONVEYANCE_ID 50;
  - T_SURVEYAGENT.ADDRESS / CITY / COUNTRY 500.
  - Batasnya sama dengan validasi layar. Lebar diukur dalam BYTE (len UTF-8).
  - Master TRADING / GOODSTYPE / PACKING / CONVEYANCE / LLOYDAGENT = view atas JSON; lebarnya tidak tertulis di DDL.
  - Galat dikumpulkan seperti galat cargo lain (P6), dan pekerjaan menjadi gagal_validasi.
- **Warranty Extension (500)** tidak ditulis CSV cargo (K74-15: coverage hanya TSI + Currency), jadi tidak ada nilai master untuk
  diperiksa. Di PUT objek-cargo, kolom itu sudah divalidasi 500 → 400.
- `services/csvcargo.go` `lebarMasterCargoCsv` (dipanggil `periksaBarisCargoCsv`). Uji `TestLebarMasterCargoCsv`: tepat di
  batas sah, tujuh kolom, multibyte, dan lewat periksaBarisCargoCsv.
- **Mutasi:** batas `>=` dan ukur rune — tertangkap.

## Catatan silang — IDPEGA = case id polos (11-10-2026)

Tulisan IDPEGA tiket ini kini memakai case id polos; pembaca menerima dua bentuk (polos / berawalan). Rincian: `docs/issues/69-jalur-policy.md`
bagian "IDPEGA = case id polos".
