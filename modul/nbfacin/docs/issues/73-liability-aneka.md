# 73: Class of Business LIABILITY (keluarga Aneka) — port layar & alur

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 09-10-2026 (case LIABILITY baru,
> tab detail kosong) "ini kenapa? saya baru buat cob liability" → pilihan "Analisis dulu, lalu port". Riset: subagen sesi 0f
> (laporan lengkap: transkrip sesi; skrip `an_*.py` di scratchpad); klaim inti diverifikasi ulang (`DecisionTable\
> BusinessType_DeT.xml` L360/L489–497/L674, `When\IsLiability.xml` pyLogic, keberadaan `WORK!INPUTDTLCOVERAGE_FACIN`).

**Status:** needs-info — P0 selesai (backend c3 + frontend 0f); spesifikasi per-lini (K73-2) selesai 09-10-2026 (bagian
"Spesifikasi per-lini" di bawah); P1 menunggu keputusan **K73-4…K73-9**. Sementara: tab case non-FIRE menampilkan "Tab … untuk
Class of Business … belum tersedia di sistem baru (saat ini baru FIRE)" (sesi 0f 09-10-2026).

## Temuan inti `[terverifikasi]`

- Rute bukan "Kode Bisnis" (teks `pyConditionString` usang): `BusinessType_DeT` (first match) — GroupPanel "003" + OldId ∈
  {15,24,89,90,91,E7,E9} → **"Liability"**; GroupPanel 003 lain → "Aneka" (+ Bonding, CAR/EAR, MBD, Burglary, Glass, …);
  GroupPanel 006 → Fire / FireStyle1/2. `IsLiability` = BusinessType Liability ∨ Products/Professional Liability (Aneka +
  10021/10023) ∨ Workmen ∨ 10048 ∨ 10184 ∨ Billboard/Neon. `IsAneka` mencakup IsLiability + ±25 lini lain.
- Layar: sakelar `IsFire` / `!IsFire` di `Section\InputInwardFacultativeDtl.xml`. Liability: tab Object
  `InputDtlObject_FacIn` (Lokasi → Occupation List → Object List `.AnekaList` → form `InputDtlObjectAneka_FacIn`), tab
  Coverage `InputDtlCoverage_FacIn` (Lokasi → Occupation → Aneka → CoverageList → `InputCoverageAneka_FacIn` + Deductible),
  Spreading & Commission lewat rantai Aneka yang memakai ulang `InputDtlSpreadingCoverage_FacIn` /
  `InputDtlSpreadingCommisionCoverage_FacIn` (sudah diport untuk FIRE); Payment sama; **Clauses & Scoring Risk tersembunyi**.
- Data: `LocationList.Property.RiskLocation.OccupationList.AnekaList.CoverageList.{SpreadingList, DeductibleList}` (FIRE:
  `…Property.PropertyItemList.CoverageList`). Tabel: T_OCCUPATIONLIST / T_COVERAGELIST / T_DEDUCTIBLELIST / T_SPREADINGLIST ada;
  **T_ANEKALIST belum ada** (skema tiket 22 `skema_gen.go` L42–69); T_COVERAGELIST kurang CALCULATE_METHOD_FAC_IN /
  PCT_SHORT_PERIOD / LOADING / SUB_LIMIT_NOTE / ANEKA_ID.
- Premi: `CountPremiCoverageAneka` (Pro Rata / Short Period / Fix Rate, Loading, LossLimit 100) — sebagian sudah diport tiket 18
  (`premium/lini_lain.go` `calculateAnekaGolf`, metode 3).
- Hilir: Submit `ProtectCoverage_Act` langkah 5 (Aneka) + `CheckSpreadingProtectAnekaGolf_ACT`; selesai Policy
  `SaveTreatyProduction_Act` langkah 19 (NB Aneka → FACINPRODUCTION); CSV `InsertUploadAneka_act`.
- **Rule hilang dari ekspor NB FacIn** (hanya di ekspor Endorsment Fac In atau tidak ada sama sekali): varian Work
  `InputDtlCoverage_FacIn`; flow action `ObjectDtlOccupation_FacIn`, `ObjectDtlAneka_FacIn`, `InputCoverageAneka_FacIn`
  (Data-Aneka), `InputDtlCoverage`; activity `AddCoverageAneka_act`, `CheckCoverageAneka_Act`, `CopyCoverageFormSection_Act`,
  `CopyDeductibleToAllLocationAneka_Act`; DT `setTSIAneka_preDT`; tidak ada di ekspor mana pun: `ChooseCov_Aneka`,
  `InputDtlDeductibleAneka_FacIn`, `addDeductibleAneka_ACT`, `CountPCTLOL`, `CalculateRate`, `SetTSI_DT`, `SetCurrency_DT`,
  When `IsObjectLimiliabilityAneka`.

## Rencana bertahap (usulan agent)

P0 gerbang COB (BusinessType dari BUSINESS GROUPPANEL + OLDID; IsFire / IsLiability / IsAneka dari registry predikat) →
P1 Object (migrasi T_ANEKALIST + ALTER T_COVERAGELIST; Lokasi → Occupation → Aneka) → P2 Coverage + premi + Payment →
P3 Spreading + Commission + panel cedant non-FIRE → P4 proteksi Submit → P5 selesai Policy (FACINPRODUCTION langkah 19) →
P6 CSV Aneka.

## Keputusan work owner 09-10-2026

- **K73-1** Rule yang hilang dari ekspor NB FacIn: **work owner mengekspor XML-nya** ke `D:\migrasi\RNM\DDL\` (bukan
  meminjam salinan Endorsment Fac In). Daftar diserahkan ke work owner 09-10-2026 (lihat Temuan inti).
- **K73-2** Cakupan: **semua lini Aneka sekaligus** (bukan Liability dulu) — Liability, Products/Professional Liability,
  Workmen, Billboard, Bonding/KBG/Custom Bond, CAR/EAR, MBD, Burglary, Glass, All Risk, Marine Hull, Aviation Hull, HE, Land
  Rig, Contractors P&M, Electronic Equipment, Fidelity, CIT/CIS, Boiler, Growing Trees, Yield Shortfall, Crime, Environmental,
  Golf (berbagi section). Analisis per-lini (kontainer form objek per lini, tabel anak T_AVIATIONHULL / T_MARINEHULL /
  T_VEHICLEHE) menyusul sebelum P1.
- **K73-3** No KTP / No NPWP form objek: **disimpan, ikut Pega** (kolom baru lewat migrasi; data pribadi — tanpa contoh nyata di
  uji / dokumen, CLAUDE.md §10).

## Langkah berikut

1. P0 gerbang COB di backend (tidak bergantung XML yang kurang) — diserahkan ke sesi c3.
2. Analisis per-lini Aneka (sesi 0f) + XML dari work owner.
3. P1–P6 sesudah 1–2.

## Backend P0 — gerbang Class of Business (sesi c3, 09-10-2026)

Berkas: `services/lini.go` (+ `BisnisLini` di `repository/tableoflimit.go` — semula `GroupPanelBisnis`, diganti K73-B5); kunci baru
KasusNB di `handlers/kasus.go`.

- **K73-B1 BusinessType** `[terverifikasi]`:
  - `Activity\SetBusinessType_Act.xml` langkah 2: Property-Map-DecisionTable `.Quotation.BusinessType` = `BusinessType_DeT` (langkah
    1 yang men-SET GroupPanel / BusinessCode di-REMARK).
  - `DecisionTable\BusinessType_DeT.xml`: kolom `.Quotation.GroupPanel` / `.Quotation.BusinessOldId` (operator `=`), 36 baris,
    first match (pyEvaluateAllRows "no", L238), `pyDefaultResult` "UNKNOWN". Diurai sesi c3: sel utama per REPEATINGINDEX +
    pyOrConditions per pyRowNum **0-based** (pyRowNum 7 = baris 8 "Liability", L489-497); sel kosong = cocok apa pun.
  - Sensus dua cara: tabel Go = urai XML (36 = 36 baris, urutan sama, A\B = B\A = ∅).
  - ⚠️ Baris yang tidak ada di kontrak: **006 + OldId 78 -> "OilGas"** (baris 29); 001 / 008 -> Medicare, 002+03 PA, 004 MarineCargo,
    005 Travel, 007 MBUCar / MBUMotorCycle, 009 Life.
  - Sumber sistem baru: baris BUSINESS Class of Business case (NOTE + BUSINESSGROUPID, cara kode bisnis tiket 40 / OLDID tiket 62).
    GROUPPANEL VARCHAR2(4000) `[terverifikasi]` `DDL\BUSINESS.txt` L5 (dibandingkan teks). 0 / > 1 baris -> "UNKNOWN".
  - Instrumen: sampel nyata 003/30 -> Aneka, 006/22 -> FireStyle1.
- **K73-B2 Predikat** lewat `rules.Eval` (registry korpus): IsFire, IsAneka, IsAnekaNotBonding, IsLiability, IsBondingAndCustomBonds,
  isGolfInsurance.
  - Jalur yang dibaca (dicacah dari `registry_gen.go`, rekursif; 210 predikat terurai = 210 `asal:`): BusinessType dalam empat
    bentuk (`pyWorkPage.Quotation`, `pyWorkPage.OfferFacIn.QuotationData`, `pyWorkPage.Policy.Quotation`, `.Quotation`) — satu
    nilai; `pyWorkPage.Quotation.BusinessCode` = BUSINESS.ID; `pyWorkPage.Quotation.BusinessName` = Class of Business
    (BUSINESS.NOTE).
  - Dasar `[dugaan kuat]`: `Activity\CopyToPolicyListPASSG.xml` L12143-12229 mengisi BusinessCode 10028, BusinessName "PERSONAL
    ACCIDENT", GroupPanel 002.
  - Hasil korpus (diuji):
    - Liability -> aneka / anekaNotBonding / liability benar;
    - Bonding -> aneka benar, anekaNotBonding salah;
    - GolfInsurance -> golf benar, **aneka salah**;
    - Fire / FireStyle1 -> fire;
    - Class of Business "BILLBOARD/NEON SIGN" (IsBillboardNeon) -> liability.
  - Lunak: galat predikat (ErrTafsirBerbeda, mis. BusinessCode kosong dibanding angka) -> medan false + log.
  - IsHE / IsMBD / IsEngineering ada di registry, belum dikabelkan (tidak diminta di kunci `lini`).
- **K73-B3 Kunci KasusNB** (GET kasus, PUT general, POST following - lewat `denganAksesTab`): `businessType: string`, `lini: {fire,
  aneka, anekaNotBonding, liability, bonding, golf}`.
- **K73-B4 BusinessCode 5 digit vs 2 karakter** (`[pertanyaan terbuka]`, tidak ditebak):
  - Ketiga predikat objek Aneka membandingkan `pyWorkPage.Quotation.BusinessCode` dengan nilai 2 karakter:
    - IsObjectSectionAneka: 24, 15, 32, 17, 18, 27, 07, SV, SX, 90, 91, T9;
    - isBisnisOtherObjAneka: 15, 32, SQ, 21;
    - IsInterestInsuredOtherObjAneka: 05, B7, 19, 15, 32, 06, 17, 18, SQ, SK, SL.
  - Predikat lain membandingkannya dengan 5 digit (10021, 10048, 10184, …).
  - Nilai 2 karakter itu berbentuk BUSINESS.OLDID (`VARCHAR2(5)`; 15 / 24 / 90 / 91 = OldId Liability di BusinessType_DeT). Di NB
    FacIn BusinessCode = 5 digit, jadi ketiga predikat itu `[dugaan]` tidak pernah benar di FacIn, kecuali BusinessCode di halaman
    lain diisi OLDID.
  - Tidak dipakai P0; perlu dijawab sebelum P1 (form objek Aneka).
- Uji: `TestBusinessTypeDari` (22 kasus + jumlah baris), `TestLiniKasus` (7 kasus perbandingan struct penuh + billboard + tanpa DB /
  404). Mutasi 8 / 9 mati; 1 ekuivalen (jalur `pyWorkPage.Policy.Quotation.BusinessType` di IsCustomBonds = OR dengan
  `pyWorkPage.Quotation.BusinessType` bernilai sama).
- ⚠️ `TestKunciKasusSamaDenganApiTs` MERAH sampai api.ts KasusNB memuat `businessType` + `lini` (sesudah `bolehFollowing`).
  Frontend P0 (sesi 0f): `api.ts` memuat kedua kunci; tab Clauses / Scoring Risk hanya `lini.fire`; Go + vitest hijau.
- **K73-B5 Saringan `GROUPPANEL IS NOT NULL`** (keputusan work owner 09-10-2026, lewat sesi 9d) `[terverifikasi]`:
  - Bukti: `D:\migrasi\RNM\NB FacIn\RDBList\GetOldIDBusiness_SQL.xml` — `select oldid as CARI1, GROUPPANEL as CARI2, ID as CARI3
    from business where (ID ={OldID.CARI2} OR NOTE ={OldID.CARI2}) AND GROUPPANEL IS NOT NULL`.
  - Kueri lini kini `repository.sqlBisnisLini`: `SELECT ID, OLDID, GROUPPANEL FROM <BUSINESS> WHERE NOTE = :1 AND BUSINESSGROUPID =
    :2 AND GROUPPANEL IS NOT NULL ORDER BY ID FETCH FIRST 2 ROWS ONLY`.
  - ⚠️ Kenapa satu kueri, bukan hanya menambah saringan di `GroupPanelBisnis`: lini membaca ID (BusinessCode), OLDID, dan
    GROUPPANEL. Bila hanya kueri GROUPPANEL yang disaring, dua baris NOTE + group (satu NULL) memberi 1 GROUPPANEL tetapi 2 OLDID →
    OldId kosong → baris DT yang salah (mis. 003 + "" → Aneka, bukan Liability). Ketiga kolom kini dari baris yang sama sesudah
    saringan. `GroupPanelBisnis` dihapus (pemakainya hanya lini).
  - TIDAK berubah: `KodeBisnis` (tiket 40 / Table of Limit) dan `OldIDBisnis` (tiket 62 / bentuk B tangga, butir 89) — teks SQL
    keduanya tetap dipatok `TestSQLTableOfLimit`; jalur FIRE tidak disentuh selain lewat lini yang sama.
  - Temuan, TIDAK dikerjakan (menunggu 9d / work owner): Pega mengenali baris dengan `ID = x OR NOTE = x` tanpa
    `BUSINESSGROUPID`; sistem baru NOTE + BUSINESSGROUPID (butir 89).
  - Uji: `TestLiniKasus` — dua baris NOTE + group sama, satu GROUPPANEL NULL → baris ber-GROUPPANEL (Liability; juga urutan NULL
    lebih dulu → FireStyle1); satu baris NULL saja → UNKNOWN; tanpa baris → UNKNOWN; BusinessCode dari baris ber-GROUPPANEL
    (IsLiability E `10048`); kasus lama 003/30 Aneka dan 006/22 FireStyle1 tetap. Tiruan menerapkan saringan yang teks SQL-nya
    dipatok `TestSQLTableOfLimit`. Mutasi: saringan dibuang / kolom ditukar / > 1 baris diterima / kode dari OLDID — semua mati.

## Backend P1 — `tampilanAneka` (sesi c3, 09-10-2026)

Berkas: `services/lini.go` (`TampilanAneka`, `kontainerAneka`, `gerbangAneka`, `topRiskAneka`, `tampilanAneka`), kunci KasusNB
`tampilanAneka` di `handlers/kasus.go` (lewat `denganAksesTab`, sesudah `lini`).

- Bentuk: `tampilanAneka: {kontainer: string[], kolomTahunJumlah, he, golf, tanamanHasil, engineering, bonding, briguna, topRisk,
  sectionLol, jumlahSatuan: boolean}`; `kontainer` selalu larik (`[]` bila tanpa data / di luar gerbang).
- Gerbang `[terverifikasi]` `NB\Section\InputDtlObject_FacIn.xml` L36463 (`IsAnekaNotBonding || IsHE || IsGolfInsurance ||
  IsGrowingTrees`) ATAU L22724 (`IsBonding || IsCustomBonds`); di luar gerbang semua false.
- Kontainer C1–C9 `[terverifikasi]` visible-when `NB\Section\InputDtlObjectAneka_FacIn.xml` L413 / L9678 / L17776 / L19883 / L23090 /
  L24543 / L27780 / L32104 / L34670 (nama When persis; 36 nama ada di registry, tidak ada yang panik).
- Sakelar: `kolomTahunJumlah` IsObjectWithQuantityYear; `he` IsHE; `golf` isGolfInsurance; `tanamanHasil` isGrowingTrees ||
  IsYieldShortfall; `engineering` IsEngineering; `bonding` IsBonding; `briguna` IsCreditBriguna (L4166 / L4465); `sectionLol`
  IsObjectSectionAneka (K73-4: kode 2 karakter ditiru); `jumlahSatuan` = BusinessCode ≠ 20 (L23869 `OTHER`).
- ⚠️ `topRisk`: kondisi sel `.Property.IsTopRisk` ada di `NB\Section\InputDtlObjectLocation_FacIn.xml` **L6301** (cell 17), bukan
  L6064 seperti di kontrak: `IsMBD || IsHE || IsContractorsPlantMachinery||IsElectronicEquipment||IsCMI || IsMarineHull||IsLiability`.
- Uji `TestTampilanAneka`: Liability C1, MarineHull C2, HE C5+C6, Boiler C1+C4+C5, Golf C1, FIRE []; Briguna (10235), Engineering
  (10032), SectionLol ("15"), BusinessCode 20. `TestTampilanAnekaKabel` (HTTP, Boiler / HE / Golf) dan `TestBacaKasus` (19 kunci,
  `kontainer: []`). Mutasi 11 + 5 mati.
- ⚠️ `TestKunciKasusSamaDenganApiTs` MERAH sampai `api.ts` KasusNB memuat `tampilanAneka` (sesudah `lini`).

### Migrasi 214 (sesi c3, 09-10-2026)

- `214_objek_aneka.sql` / `_down`: T_ANEKALIST (27 rancangan + 23 kolom form), T_MARINEHULL / T_AVIATIONHULL / T_VEHICLEHE
  (UNIQUE PARENT_ID, FK), T_MEMORANDUMACCEPT (baru, 16 medan), + 8 kolom T_COVERAGELIST. Ditulis, TIDAK dijalankan.
- Selisih dari kontrak `[keputusan agent]`: + `T_ANEKALIST.ALM_RISK_ID` dan `T_VEHICLEHE.OBJECT_ID_HE` (di-SET lookup kapal /
  objek HE); lebar sumber (OCCUPATION_ID 1000, FLAG + kolom MarineHull 1000 = POOLDATA.SHIP, OBJECT_NAME 4000);
  TYPE_OF_DISCOUNT VARCHAR2(10) (sampel "true" 95/95, dua cara; kolomTeksMenyimpang butir 68.1); TSI_OLD ditulis aplikasi
  (loader V-28b membuang medan *Old); penunjuk IDX_* tidak dibuat (pola 198).
- Amandemen loader (`amandemenAneka` / `Kapal` / `HE` / `Memorandum` + 2 jalur), STRUKTUR-TABEL (+ 5 bab, T_COVERAGELIST +8),
  `migrasi_sebagian_test` 457 kolom. Mutasi 7/7 mati.

### Backend P1 — `GET`/`PUT …/kasus/{caseId}/objek-aneka` (sesi c3, 09-10-2026)

Berkas: `models/aneka.go`, `repository/aneka.go`, `services/aneka.go`, `handlers/aneka.go`.

- Kabel `{lokasi: LokasiAneka[]}`; LokasiAneka = medan alamat ObjekFire + `okupasi[]` + `objekGolf[]`; ObjekAneka teks, desimal teks
  (ADR-0003), tanggal DD-MM-YYYY; anak `briguna` / `marineHull` / `aviationHull` / `vehicleHe` (+ `objectIdHe`, `brand`, `type`) /
  `memorandum` (`suretyType` baca-saja); + `almRiskId`.
- Simpan: lokasi = baris T_LOCATIONLIST / T_PROPERTY / T_RISKLOCATION yang SAMA dengan tab FIRE; PUT mengganti seluruh lokasi
  (sisa anak FIRE ikut dihapus lewat `sqlHapusObjek` yang tidak diubah). Okupasi Aneka `PARENT_TABLE` T_RISKLOCATION
  (jalur `LocationList/Property/RiskLocation/OccupationList`), objek T_ANEKALIST (T_OCCUPATIONLIST / Golf T_RISKLOCATION), anak
  hanya bila terisi. ⚠️ Kebalikannya: Save tab Object FIRE / CSV / Copy Following pada case Aneka meninggalkan baris Aneka
  yatim — frontend tidak boleh memanggilnya untuk case Aneka (diangkat ke 9d).
- Coverage Aneka (P2): T_COVERAGELIST ber-`PARENT_TABLE` T_ANEKALIST dipindah induk ke objek baru di posisi yang sama
  (lokasi, okupasi / Golf, objek); objek yang hilang -> spreading, deductible, coverage-nya dihapus (`rencanaCoverageAneka`).
- SuretyType `[terverifikasi]` `NB FacIn\Activity\InputDtlObjectAneka_PreAct.xml`: BusinessCode 10076 SB-BB, 10007 SB-PB,
  10077 SB-APB, 10083 SB-MB, 10011 SB-PayB, 10075 CB, 10126 KBG; kode lain -> kosong; isian klien diabaikan.
- Penjaga: identitas 401 -> hak ubah = tab Object FIRE (`wajibBolehUbah`, matriks TestPenjagaBacaSaja 15 metode) -> lebar /
  tanggal 400 -> 503 -> mata uang di CURRENCY (kosong boleh: Aviation Hull tanpa Currency; CURRENCY_CODE tidak disebut ->
  DEFAULT, dibaca kosong pola A169) -> 404.
- Proteksi 422 (`services/proteksianeka.go`) `[terverifikasi]` `NB FacIn\Activity\ProtectionObjectAneka_Act.xml` = dispatcher
  12 `Call ProtectionObject<Lini>_Act` (tanpa loop, tanpa pesan sendiri; post-processing flow action `InputDtlObjectAneka_FacIn`
  FromWhere "0"); 63 literal pesan (grep dua cara), 3 ber-label `//` tidak diport -> 60 aktif. Dijalankan untuk setiap objek
  saat Save, sebelum transaksi. Kabel 422 `{galat, pelanggaran:[{lokasi, okupasi (null = Golf), objek, medan, pesan}]}`;
  `medan` = Field Pega apa adanya; akhiran ` \n` literal (dua karakter) dipotong `[keputusan agent]`.
  - Gerbang: HE `IsHe` (rule IsHE), Electronic Eq `isElectronicEquipment` (pembanding bertepi spasi — dengan BusinessType
    "ElectronicEquipment" gerbang TIDAK terbuka; pemangkasan Pega `belum terverifikasi`), 10 lainnya BusinessCode 2 karakter
    (K73-4: tidak pernah terbuka di FacIn, 5 digit).
  - Anomali dibawa apa adanya: Golf / Fidelity pesan Section -> medan `quantity`; Electronic Eq pesan Quantity -> `noOfTree`;
    AllRisk `pyWorkPage.Aneka.TerritorialLimits` (tidak ada di model -> selalu kosong); Golf / Fidelity `.LimitOfLiability`.
  - `[pertanyaan terbuka]` GrowingTree langkah 7/8: teks di-SET bila `.NoOfTree` kosong, ditempel bila `.AreaHectar` kosong
    — AreaHectar kosong + NoOfTree terisi = pesan kosong di Pega; efeknya belum terverifikasi, TIDAK dikirim. Electronic Eq
    `Local.cekIdemMonth` vs `local.cekIdemMonth` dianggap variabel sama `[dugaan]`.
  - Uji `TestTabelProteksiAneka` (12 langkah urut, 60 pesan), `TestProteksiObjekAneka`, `TestGantiObjekAnekaProteksi`,
    `TestObjekAnekaProteksi422`. Mutasi 12/12 mati.
- Lookup (`repository/lookupaneka.go`): `GET /api/nbfacin/kapal?cari=` -> `{baris:[{id, nama, flag, tahun, nrt, grt, dwt,
  klasifikasi, konstruksi}]}` dan `GET /api/nbfacin/tipe-kendaraan` -> `{baris:[{id, typeName}]}`.
  - Kapal `[terverifikasi]` `NB FacIn\ReportDefinition\BrowseShip_RD.xml`: `(REMARK != 'Not Active' OR REMARK IS NULL)` (peka huruf),
    `UPPER(NAME) LIKE %cari%` (Contains pyCaseInsensitive) hanya bila diisi, `ID DESC` (VARCHAR2 -> urut teks), maks 500. Tabel
    `[dugaan kuat]` POOLDATA.SHIP (`DDL\SHIP.txt`); pyUseAlternateDB -> DB alternatif `[pertanyaan terbuka]`, dibaca skema biasa.
    Pemetaan saat dipilih (frontend; dicek sesi 9d `InputDtlObjectAneka_FacIn.xml` L10480-10570): id -> almRiskId, nama ->
    objectName, flag -> flag, tahun -> year, nrt / grt / dwt / klasifikasi / konstruksi -> marineHull.*.
  - Tipe `[terverifikasi]` `BrowseType_RD.xml`: view POOLDATA.BRANDDETAIL `TYPE = 'ANEKA'` (peka huruf), TIDAK disaring merk / objek
    (form mengirim parameter kosong), `ID DESC`, maks 500. id -> vehicleHe.type, typeName -> vehicleHe.typeName.
- Choose Occupation Aneka `GET /api/nbfacin/kasus/{caseId}/occupation-aneka?q=` -> `{baris:[{id, name}]}` (permintaan sesi 9d,
  temuan sesi 0f; ditelusur ulang sesi c3) `[terverifikasi]`:
  - Tombol `NB FacIn\Section\InputOkupasiAneka_FacIn.xml` L1224 -> HARNESS `Harness\ChooseOccupation.xml` (bukan section FIRE
    tiket 40): grid `D_Occupation.pxResults` L2307; parameter DP L3803-3806 / L3814-3817: BUSINESSNAME =
    `pyWorkPage.Quotation.BusinessName` (bergantung case -> endpoint di bawah case; = Class of Business case `[dugaan kuat]`),
    TYPE = `ANEKA` TANPA kutip (FIRE `"FIRE"` berkutip, `Section\ChooseOccupation.xml`; ditafsir literal `[dugaan]`), NAME =
    `InputParam.CARI5` (Search), ID kosong.
  - `DataPage\D_Occupation.xml` -> `ReportDefinition\BrowseOccupationANEKA_RD.xml` (Int-OCCUPATION): `A AND B AND C AND D` -
    `.ID =`, `.Type =`, `.BusinessName =` (peka huruf), `.Name Contains` (pyCaseInsensitive); urut `.ID` DESC; maks 500; filter
    kosong diabaikan `[dugaan]`. POOLDATA.OCCUPATION (`DDL\OCCUPATION.txt`, VARCHAR2 -> ID diurut teks).
  - Pilih: harness L3314 / L3478 `GetOccupation_act(ID=.ID, Name=.Name)` -> `Activity\GetOccupation_act.xml` Property-Set
    `.OccupationId = Param.ID`, `.OccupationName = Param.Name` (L357-404), lalu SetOccupation_DT (IdxLocation / IdxOccupation) dan
    Obj-Save pyWorkPage. Jadi **occupationId = OCCUPATION.ID** (bukan OLDID seperti FIRE). Sensus dua cara (urai ElementTree vs
    potong teks blok RiskLocation) `DDL\CONTOH`: 126 OccupationId okupasi Aneka - 122 angka 8 digit, 4 angka 5 digit (sama).
  - ⚠️ `[pertanyaan terbuka]` harness memanggil `GetOccupation_act` di kelas **Int-OCCUPATION**, sedangkan satu-satunya varian di
    korpus (NB / Endorsment / RNW) berkelas **Data-Occupation** - aturan yang benar-benar berjalan belum terverifikasi; pemetaan
    di atas memakai varian yang ada.
- `SetObjectDescHE_DT` (`NB FacIn\DataTransform\SetObjectDescHE_DT.xml` L155) dihitung FRONTEND on-change (keputusan sesi 9d);
  backend menyimpan `objectName` apa adanya.
- ⚠️ Risiko tercatat (keputusan 09-10-2026, jawaban (c)): Save tab Object FIRE / CSV objek / Copy Following atas case Aneka
  meninggalkan baris Aneka yatim. TIDAK dijaga di server (409 berbasis lini ditolak: jalur FIRE tidak boleh bergantung pada lini);
  penjagaannya frontend - tab tersebut hanya tampil di layar FIRE.
- Pesan 422 dipotong akhiran ` \n` literal di server — `[keputusan agent]`, diterima sesi 9d 09-10-2026.

## Backend P2 — coverage Aneka (sesi c3, 09-10-2026)

Riset korpus tiga subagen (laporan di scratchpad sesi c3 `riset67\p2_onchange_cpca.md`, `p2_activity_premi.md`,
`p2_agregat_lookup.md`); butir kunci dicek ulang sesi c3 ke baris XML.

- **Data** (`models/aneka.go` `CoverageAneka`, `repository/covaneka.go`): `ObjekAneka.coverages` di GET / PUT `…/objek-aneka`.
  T_COVERAGELIST ber-`PARENT_TABLE` T_ANEKALIST, SRC_PATH `…/OccupationList/AnekaList/CoverageList` atau Golf
  `…/RiskLocation/AnekaList/CoverageList`; deductible T_DEDUCTIBLELIST (pola tiket 45). PUT: kunci absen / null = coverage
  tersimpan dipertahankan per posisi (P1); larik = ganti utuh (coverage lama posisi itu dihapus beserta deductible / spreading).
  Mata uang coverage kosong -> DEFAULT 'UNKNOWN', dibaca kosong (A169). Medan server (minPremium, tsiLiability, tsiNusantaraRe,
  premiNusantaraRe) diabaikan saat PUT - rumusnya menunggu B2.
- **Total per mata uang** (`services/totalaneka.go`) `[terverifikasi]` `NB FacIn\Activity\SumTotalTSIPremiGross_Act.xml`:
  per lokasi x mata uang (okupasi dilebur), premi gross tanpa diskon, rate = premi / TSI x 100 (L3058 / L4583); Golf (langkah 2)
  dan Aneka biasa (3a) TSI objek + premi ke mata uang objek; 10 tipe khusus (3b, L4768-L4975) per coverage. Tidak disimpan.
- **Lookup** `GET /api/nbfacin/kasus/{caseId}/coverage-aneka?q=` (`repository/lookupcovaneka.go`): `BrowseCoverageFacIn_RD`
  Type ANEKA, BizCode = BusinessCode case; sumber tabel COVERAGE (butir 96; kelas RD = view COVERAGE_FACIN tanpa ACTIVESTATUS).
- **Hitung** `POST /api/nbfacin/kasus/{caseId}/hitung-coverage-aneka` (`services/hitunganeka.go`): badan `{proRatePercent, baris:
  [{coverage, pemicu}]}` -> `{proRatePercent, baris: CoverageAneka[]}`, tanpa menyimpan; ProRate case berantai antar baris. Pemicu
  (rantai on-change `InputCoverageAneka_FacIn.xml`, dua definisi aksi identik): metode / loading / indemnity / diskon -> CPCA
  tanpa Status; currency / tsi / rate -> CPCA Status="Coverage" (langkah 4 L964 menimpa ProRate case dari tanggal); proRate ->
  fillPremiBond (`.ProRatePercent` COVERAGE, kosong 100; minPremium = premium); shortPeriod / scale -> FillPremiGolf (L399 SELALU
  dari tanggal); day -> ChangeCoveragePremium (tanpa premi); premium -> CalculateRate (pecahan dijumlah dua kali, L718 / L1003);
  diskonPersen / diskonJumlah -> ChangeDiskon pre-DT lalu CPCA; lol / pctLol -> CountPCTLOL. 14.1 (ProRate case 0 / kosong, metode
  1) dihitung dari tahun polis tanpa menulis case; MBD (17.x) memakai ProRate case apa adanya.
  - Instrumen: premi 4 coverage nyata `DDL\CONTOH` (NB-143055 L134 11905.2442, NB-147721 L119 371698962.76865742, NB-183524 L101
    27128887.5, NB-183603 L120 850000) dan medan RNM NB-143055 / NB-183524 = nilai sampel.
- **Medan server saat PUT** (`medanServerAneka`, K73-13): TSINusantaraRe = Share x TSI / 100, PremiNusantaraRe = Share x Premium / 100
  - Share x Discount / 100 (skala 20 tiap suku, `CountPremiAndTSIRNMFireAnekaGolf_ACT` L3853 / L3905 / L6499 / L6551); TSILiability =
  TSI (okupasi) / kosong (Golf); minPremium = premium; Share = T_GENERAL_POLIS.PERCENT_SHARE. ⚠️ Hanya coverage yang DIKIRIM PUT;
  perubahan % Share sesudahnya (tab Spreading) tidak menghitung ulang coverage Aneka (FIRE: `sqlUbahNR`) - catatan untuk tiket
  spreading Aneka.
- **% Pro Rate case** (K73-12): migrasi `255_pro_rate_percent.sql` `T_GENERAL_POLIS.PRO_RATE_PERCENT NUMBER(38,8)` (rentang kedua
  nbfacin 255-259 dari rnwfacin, keputusan work owner 09-10-2026; 214 sudah jalan di DEV) - WAJIB sebelum biner. GET objek-aneka
  `proRatePercent` = tersimpan, NULL -> dari tanggal (CPCA langkah 4, Day 365); PUT menulis (kosong = NULL). Penjaga
  `inti/backend/penjaga` `rentangNomorJamak` menerima rentang dipisah koma (`TestRentangMigrasiJamak`).
- Uji: `hitunganeka_test.go`, `totalaneka_test.go`, `aneka_test.go` (repository / services / handlers), `lookupcovaneka_test.go`,
  `TestRentangMigrasiJamak`. Mutasi 25/25 mati.
- `[pertanyaan terbuka]` tercatat: arti "KBRU" / IsB2B (dianggap bukan KBRU); `Policy.Currency == 1` (MinPremium 100000 mati);
  penjumlahan ganda CalculateRate; `m_bi_indemnity` (SetIndemnityRate_ACT); satuan DateTime - DateTime (hari, dugaan kuat); Day
  kosong pada FillPremiGolf dari sel % Short Period (di sini .Day); Rate agregat sampel yang menyimpang (35 dari 101 baris).

- **Tambahan P2 (sesi c3, 09-10-2026)**:
  - `tampilanAneka.currencyCoverage` = visible-when sel Currency form coverage (`NB FacIn\Section\InputCoverageAneka_FacIn.xml` L4052,
    verbatim `IsPortRisk || IsBuilderRisk||IsObjectWithQuantityYear || isAllRisk || IsMarineHull || IsMarineHullOffshore || IsMBD ||
    IsHE || IsAviationHull || isBurglary`); `tampilanAneka.mbd` = IsMBD (L5826 / L6132 / L6322).
  - Pemicu `indemnity` kini port penuh `SetIndemnityRate_ACT` (DDL `D:\migrasi\RNM\DDL\M_BI_INDEMNITY.txt` ditaruh work owner):
    `SearchIndemnityRate_SQL` L84 (`INDEMNITY = .Indemnity AND replace(INDEMNITY_UNIT,' ','') = .Unit` - padding CHAR(3) dipangkas SQL
    Pega sendiri); baris terakhir menang; nol baris -> % Indemnity kosong; `.Unit == 0 || .Unit == 2` -> 100 (L1079); IsMBD -> CPCA.
    Indemnity bukan bilangan bulat <= 4 digit tidak dikueri `[keputusan agent]`. Menggantikan butir "hanya hitung ulang" K73-15.
  - Copy Deductible Aneka (frontend): `DDL\CopyDeductibleToAllLocationAneka_Act.xml` langkah 2 ber-halaman langkah
    `…AnekaList(Param.idxAneka).CoverageList` (L505) dengan repeat EMBEDDED (L834) - `CoverageList(<CURRENT>)` (L721) = coverage
    iterasi berjalan objek sumber `[terverifikasi struktur; semantik <CURRENT> = dugaan kuat]`.

### Keputusan work owner 09-10-2026 (P2, lewat sesi 9d)

- **K73-12** % Pro Rate tingkat case disimpan di `T_GENERAL_POLIS.PRO_RATE_PERCENT` (migrasi 255, rentang kedua 255-259); NULL =
  dihitung dari tanggal. IsB2B dianggap bukan "KBRU" (OQ).
- **K73-13** Medan RNM coverage Aneka / Golf memakai keluarga `CountPremiAndTSIRNM*` (pola FIRE tiket 48); TSILiability = TSI
  (Aneka), kosong (Golf).
- **K73-14** MinPremium: sakelar 100000 (`Policy.Currency == 1`) mati, minPremium = premium; arti kode = OQ.
- **K73-15** Tiru apa adanya: CalculateRate (prorate dua kali), fillPremiBond pada % Pro Rate untuk semua lini, FillPremiGolf selalu
  dari tanggal, Day -> ChangeCoveragePremium; pemicu indemnity hanya menghitung ulang premi (pencarian `m_bi_indemnity` = OQ menunggu
  DDL); uji Go Golf pro-rata diperbarui dengan bukti `FillPremiGolf.xml` L399. Lookup coverage-aneka dari tabel COVERAGE (butir 96).

### Keputusan work owner 09-10-2026 (lookup, lewat sesi 9d)

- **K73-10** Object Name HE dan Merk HE: TEKS BEBAS (Pega `pyAllowFreeFormInput=true`); tanpa endpoint, BRANDDETAIL tidak dipakai
  sebagai pengganti. Daftar pilihan menyusul setelah DDL OBJECT / BRAND ditaruh work owner di folder DDL.
- **K73-11** Draft Wording ID: TEKS BEBAS; dropdown setelah DDL tabel enumerasi (`D_EnumerationList` / Data-Enumeration) ada.
- Kapal dan tipe kendaraan: dibangun (di atas).
- Uji: `repository/aneka_test.go` (kolom = migrasi 214 dua arah, bind 1..n, traversal berkelompok, rencana coverage),
  `services/aneka_test.go`, `handlers/aneka_test.go` (kunci kabel tiap tingkat). Mutasi 21 mati; titik panggil `mataUangBaca`
  di jalur baca Oracle tidak teruji tanpa basis data (fungsinya teruji).

## Spesifikasi per-lini (K73-2) — sesi 0f, 09-10-2026

Riset: subagen sesi 0f (tiga subagen A+F form objek, B jalur coverage, C+E premi & proteksi; laporan lengkap A–H di transkrip
sesi, skrip `an2*_*.py` di scratchpad). Ringkasan di bawah adalah bagian yang mengikat P1–P4. Klaim kunci diperiksa ulang sesi 0f:
`NB FacIn\Section\InputDtlObjectAneka_FacIn.xml` L413 (visible C1), L1708/L1711 (Start Date `ALWAYS` + `IsObjectLimiliabilityAneka`),
L23867–23869 (`OTHER` `BusinessCode!=20`), `grep -c` `pyContainerVisibleWhen` berisi = 11; `Activity\ProtectionObjectAneka_Act.xml`
L495 (`"B3"`); `When\IsAneka.xml` tanpa Golf; `DDL\CONTOH\*.xml` BusinessCode panjang 5 = 128/128 (dua cara: `grep -ho | awk
length` dan per berkas). Singkatan jalur: `NB\` = `D:\migrasi\RNM\NB FacIn\`, `DDL\` = `D:\migrasi\RNM\DDL\`.

### S1. Form objek `NB\Section\InputDtlObjectAneka_FacIn.xml` — 9 kontainer `[terverifikasi]`

Hitungan dua cara: 9 kontainer tingkat atas (urai) + 2 bersarang `IsCreditBriguna` = 11 `pyContainerVisibleWhen` (grep); 84 sel
FIELD + 1 SUB_SECTION (urai = grep).

| Kontainer | Visible-when (verbatim, ringkas) | Lini `[dugaan, evaluasi When]` | Medan khas |
| --- | --- | --- | --- |
| C1 L361 | `IsAsuransiKredit \|\| IsEnvironmental \|\| IsCrime \|\| isGolfInsurance \|\| isCitCis \|\| IsLiability \|\| IsWhetherIndex \|\| isFidelity \|\| isEAR \|\| isGlass \|\| IsBonding \|\| isCAR \|\| IsCivilEngineeringCompletedRisks \|\| isBillboardNeonSyariah \|\| IsCustomBonds \|\| IsBoiler \|\| IsComprehensiveMachineriesInsurance` (L413) | Liability dkk., Golf, Fidelity, EAR, Glass, Bonding/KBG/Custom, CAR, Boiler, Kredit, CitCis, Crime, Environmental, … | ObjectName (wajib), Section, NO_KONTRAK, START/END_DATE, NO_KTP, NO_NPWP, Currency, TSI, DraftWordingID (`IsBonding`), MaintenancePeriod (`IsEngineering`), blok Briguna (`IsCreditBriguna`), Others, sub-form Memorandum |
| C2 L9626 | `IsPortRisk \|\| IsBuilderRisk\|\|IsMarineHull \|\| IsMarineHullOffshore` | Marine Hull (+Offshore), Port/Builder Risk | Vessel Name (`BrowseShip_RD`), Currency, Year (wajib), TSI, BuiltIn, Flag, `.MarineHull.*` (9) |
| C3 L17725 | `isAllrisk \|\| isBurglary` | All Risk, Burglary | ObjectName, Currency, TSI, Quantity, InvoiceNo |
| C4 L19831 | `IsObjectWithQuantityYear` | MBD, Boiler, ContractorsPM, LandRig | ObjectName, Section, Currency, TSI, Chassis, Engine, SerialNo, MaintenancePeriod |
| C5 L23038 | `IsObjectWithQuantityYear \|\| IsHE` | C4 + HE | Year, Quantity, UnitQuantity (keduanya `OTHER BusinessCode!=20`) |
| C6 L24492 | `IsHE` | HE | ObjectName, `.VehicleHE.{ObjectNameHE, ChassisNumber, EngineNumber, BrandName, TypeName}`, Currency, TSI |
| C7 L27728 | `isAviationHull` | Aviation Hull | Model/Tipe, Year, GeographicalLimits, `.AviationHull.*` (6) — **tanpa Currency & TSI** |
| C8 L32052 | `isGrowingTrees \|\| IsYieldShortfall` | Yield Shortfall | ObjectName, Currency, TSI, Year of Planting, NoOfTree, AreaHectar |
| C9 L34618 | `isElectronicEquipment` | Electronic Equipment `[belum terverifikasi]` (pembanding bertepi spasi) | ObjectName, Section, Quantity, Currency, TSI, Year, LimitIdemnityDay/Month |

- Satu lini bisa mendapat > 1 kontainer: Boiler = C1+C4+C5; MBD/ContractorsPM/LandRig = C4+C5; HE = C5+C6. CIT, CIS dan "Aneka"
  polos di luar daftar: tanpa kontainer `[dugaan]`.
- ⭐ Sel bersyarat `IsObjectLimiliabilityAneka` / `IsObjectSectionAneka` / SUB_SECTION Memorandum ber-`pyVisible=ALWAYS`
  (`pxRuleReferences` tidak mendaftar kedua When itu) → **selalu tampil** di kontainernya `[terverifikasi nilai; efek dugaan kuat]`.
- Opsi `Section` (`DDL\Section.xml` L189-208): 1 "Section I", 2 "Section II", 3 "Warranty", 4 "No Section".
- Post-processing FA: `ProtectionObjectAneka_Act`; pre `InputDtlObjectAneka_PreAct` mengisi `.MemorandumAccept.SuretyType` dari
  BusinessCode 5 digit (10076 SB-BB, 10007 SB-PB, 10077 SB-APB, 10083 SB-MB, 10011 SB-PayB, 10075 CB, 10126 KBG).

### S2. Grid objek & okupasi `[terverifikasi]`

- Tab Object: grid lokasi `NB\Section\InputDtlObject_FacIn.xml` L36960 (`IsAnekaNotBonding || IsHE || IsGolfInsurance ||
  IsGrowingTrees`) / L29960 (`IsBonding || IsCustomBonds`) → `InputDtlObjectLocation_FacIn` (blok alamat = FIRE + `IsTopRisk`).
- Okupasi: grid `.Property.RiskLocation.OccupationList` (kolom Occupation Code / Name) → `InputOkupasiAneka_FacIn`: tombol Choose
  Occupation + grid `.AnekaList` (K2: Object Name, Year/Quantity/Unit bila `IsObjectWithQuantityYear`; kolom Section & Limit of
  Liability ber-`OTHER IsObjectSectionAneka` → `[dugaan]` tak pernah tampil di FacIn, lihat S6), K3 grid HE, K4 Growing/Yield.
- Golf: `.Property.RiskLocation.AnekaList` langsung di bawah lokasi (tanpa okupasi), kolom Object Name / Currency / TSI Old / TSI.

### S3. Coverage `[terverifikasi]`

- Rantai: `DDL\InputDtlCoverage_FacIn.xml` kontainer 7 (`IsAneka || IsGolfInsurance || IsBonding || IsHE || IsCustomBonds`, L25072)
  grid lokasi → `ObjectOccupation_FacIn` (okupasi + total per mata uang) → `ObjectDtlAneka_FacIn` (objek) → section Data-Aneka
  `DDL\InputCoverageAneka_FacIn.xml` (grid CoverageList: Coverage / Gross Premium, Add, Trash, tombol salin) → form Data-Coverage
  `NB\Section\InputCoverageAneka_FacIn.xml`. FA Data-Aneka `InputCoverageAneka_FacIn` hanya ada di ekspor Endorsment
  (ekspor ulang NB diminta).
- Form coverage (tanpa medan wajib): Choose Coverage, Coverage Note (`BrowseCoverageFacIn_RD` Type ANEKA), Calculate Method
  (`DDL\CalculateMethod_FacIn.xml`: 1 PRO RATA, 2 SHORT PERIOD, 3 FIX RATE; default 1), % Pro Rate (metode 1), % Short Period
  (metode 2), Day (365/366/360, default 365), Currency (lini tertentu), TSI, % Loading (0), % Scale `.LostLimit` (100), % Rate,
  Indemnity / Unit / % Indemnity (`IsMBD`), Premium, % Discount, Discount, Limit of Liability, % LoL, Sub Limit.
- Add coverage (`DDL\AddCoverageAneka_act.xml` L294-295): Calculate Method diambil dari `LocationList(1)…AnekaList(1).CoverageList(1)`
  (indeks keras); Currency & TSI disalin dari objek. TSI coverage didorong ke objek (`SetTSI_DT`), TSI objek ke semua coverage
  (`setTSIAneka_preDT`).
- Tombol salin (kontainer `IsAneka`): Copy Coverage From This Section (objek ber-`.Section` sama, semua lokasi), Copy Deductible To
  All / (This Section) (semua coverage berkode sama, tanpa filter Method/BI), Copy Coverage To All / For TPL (⛔ activity Data-Aneka
  tidak ada di ekspor mana pun).
- Deductible `addDeductibleAneka`: 10 medan = FIRE tanpa Descriptions; `GridDeductible.tsx` dipakai ulang.

### S4. Premi `[terverifikasi]` — `NB\Activity\CountPremiCoverageAneka.xml`

- Metode 1: `(.TSI*.Rate*Prorate)/10⁴ + (.TSI*.Rate*Prorate*.Loading)/10⁶` (skala 4, L3374) × LossLimit/100; metode 2 memakai
  `.PctShortPeriod`, metode 3 memakai 100; LossLimit 0/"" → 100 (L2480); MBD memakai `.IndemnityPercentage` dengan pembagi 10⁶/10⁸
  (L3985). Langkah 4: ProRatePercent dihitung ulang dari tanggal polis saat `Status="Coverage"`.
- Golf `FillPremiGolf` langkah 2 selalu menghitung ulang ProRatePercent; Bond `fillPremiBond` memakai ProRatePercent coverage +
  MinPremium 100000 bila `Policy.Currency == 1`.
- Sudah diport tiket 18 (`premium/lini_lain.go` `calculateAnekaGolf`): metode 1/2/3, skala 4, LossLimit, MBD. Belum: langkah 4 & 14.1,
  Golf langkah 2 (uji Go "Golf pro-rata 0 → premi 0" bertentangan), Bond, CalculateRate, CountPCTLOL, ChangeDiskon, agregat
  `SumTotalTSIPremiGross` / `CountGrossPremi_Act` 5.2 (Aneka, skala 20 minus diskon) & 5.3 (Golf, skala 4 tanpa diskon).

### S5. Data `[terverifikasi jalur; tipe kolom = usulan]`

- Jalur: `LocationList().Property.RiskLocation.OccupationList().AnekaList().CoverageList().{SpreadingList(), DeductibleList()}`;
  Golf: `…RiskLocation.AnekaList()`.
- Tabel baru: **T_ANEKALIST** (rancangan `skema_gen.go` L42-69, 27 kolom; PARENT_TABLE `T_OCCUPATIONLIST` / `T_RISKLOCATION`) +
  kolom form NO_KONTRAK, NO_KTP, NO_NPWP, OTHERS, MAINTENANCE_PERIOD, INVOICE_NO, NO_OF_TREE, AREA_HECTAR, LIMIT_IDEMNITY_DAY/MONTH,
  LIMIT_OF_LIABILITY, TSI_OLD, FLAG_DELETE; anak **T_MARINEHULL** (+LENGTHM/BREADTHM/DEPTHM), **T_AVIATIONHULL**, **T_VEHICLEHE**
  (+CHASSIS_NUMBER/ENGINE_NUMBER), **T_MEMORANDUMACCEPT** (16 kolom, tanpa rancangan).
- ALTER **T_COVERAGELIST**: CALCULATE_METHOD_FAC_IN, CALCULATE_METHOD, PCT_SHORT_PERIOD, LOADING, SUB_LIMIT_NOTE, MIN_PREMIUM,
  TYPE_OF_DISCOUNT, IS_OLD_DATA. `T_OCCUPATIONLIST` / `T_DEDUCTIBLELIST` / `T_SPREADINGLIST` dipakai tanpa perubahan struktur.
- Tidak pernah ada di 115 sampel `DDL\CONTOH` (grep 0 berkas): MemorandumAccept, NO_KTP, NO_NPWP, NO_KONTRAK, LimitOfLiability objek,
  MaintenancePeriod, LoanAccountNumber, NoOfTree, AreaHectar, InvoiceNo, LimitIdemnityDay.

### S6. BusinessCode 2 karakter (K73-B4) — sensus `[terverifikasi]`

- When `NB\When\*` + DDL: 235 perbandingan BusinessCode/BusinessOldId (grep = urai). BusinessCode vs **2 karakter: 78 kondisi di 8
  rule** (IsObjectLimiliabilityAneka, IsObjectSectionAneka, isBisnisOtherObjAneka, isInterestInsuredOtherObjAneka, IsDM,
  IsHealth/IsHealthIndividu, IsHealthCorporate); vs 5 digit: 70 kondisi di 25 rule.
- `ProtectionObjectAneka_Act`: 10 dari 12 langkah membandingkan BusinessCode dengan 2 karakter (`"B3"` L495, `"12"`, `"11"`, …);
  banyak literal itu = OldId lini yang sama di `BusinessType_DeT` ("B3" LandRig, "10" MBD, "06" Aviation, "27" Golf, …).
- Nilai FacIn: BusinessCode selalu 5 digit (128/128 sampel), BusinessOldId selalu 2 karakter. Akibat bila ditiru apa adanya:
  proteksi objek LandRig…Fidelity mati, Quantity/Unit selalu tampil, kolom Section/LoL grid okupasi tersembunyi.

### S7. Proteksi (P4) `[terverifikasi]`

- Objek: `ProtectionObjectAneka_Act` 12 langkah, 60 pesan aktif (63 − 3 dinonaktifkan; grep = urai), mis. "Tahun harus diisi \n".
- Submit: `ProtectCoverage_Act` langkah 5 ANEKA (" Coverage on location "+n+" can't be empty!" L3778, Currency L4228, Deductible
  L4529/L5029, "Double Input Case " L6146) dan langkah 6 GOLF; `CheckSpreadingProtectAnekaGolf_ACT` ("% Total spreading on location
  … must be 100%!" L2268, premi L2289; Golf skala 2). Belum ada satu pun di Go.

### S8. Rule yang hilang dari semua ekspor `D:\migrasi\RNM\` `[terverifikasi]`

Section `ChooseCov_Aneka` (FA-nya ada), `SetParam_Act`, `CopyCoverageToAllLocation_Act` kelas Data-Aneka, `SaveDeductibleAneka`,
`ProtectCurrencyAneka`, `ObjSave`, `InputCoverageAneka_FacIn_IsUW`, Rule-Obj-Class Data-Aneka/Data-Coverage. Hanya di ekspor
Endorsment: FA Data-Aneka `InputCoverageAneka_FacIn`, `InputCoverageAneka_PreAct`, `CopyCoverage_Act`, `CopyCoverageToALL_HIO_Act`,
DP `D_FilterAneka`. Uji instrumen: AddCoverageAneka_act (diketahui ada) → ketemu; SaveDeductibleAneka (diketahui hilang) → hilang.

### Jawaban work owner 09-10-2026

- **K73-4 → tiru apa adanya.** Perbandingan BusinessCode 2 karakter diport persis (predikat lewat `rules.Eval`, tanpa memetakan ke
  OldId): proteksi objek langkah 2–10/12 tidak aktif, Quantity/Unit selalu tampil, kolom Section/LoL grid okupasi tersembunyi.
- **K73-6 → NUMBER(38,8)** untuk MAINTENANCE_PERIOD, NO_OF_TREE, AREA_HECTAR, LIMIT_IDEMNITY_DAY/MONTH, PCT_SHORT_PERIOD (ADR-0016).
- **K73-7 → Briguna disimpan, ikut Pega** (seperti K73-3; tanpa contoh nyata di uji/dokumen).
- **K73-8 → work owner mengekspor sambil jalan.** Selama XML belum ada di `DDL\`: tombol Copy Coverage To All / For TPL
  disembunyikan; Choose Coverage memakai autocomplete Coverage Note (`BrowseCoverageFacIn_RD` Type ANEKA).
- K73-5 dan K73-9 tidak ditanyakan terpisah; agent memakai rekomendasinya (sel `ALWAYS` selalu tampil; ProRatePercent dihitung
  ulang dari tanggal) — dapat dikoreksi work owner.

### Keputusan yang diminta (K73-4…K73-9)

- **K73-4** BusinessCode 2 karakter (S6): tiru apa adanya (rekomendasi agent — migrasi merekam perilaku) atau baca sebagai OldId.
- **K73-5** Sel ber-`ALWAYS` dengan sisa kondisi (S1): ikut runtime = selalu tampil (rekomendasi).
- **K73-6** Tipe kolom: MAINTENANCE_PERIOD / NO_OF_TREE / AREA_HECTAR / LIMIT_IDEMNITY_* / PCT_SHORT_PERIOD; OCCUPATION_ID.
- **K73-7** Blok Briguna (data pribadi pinjaman: nama, no. KTP, tanggal lahir, gaji): disimpan seperti KTP/NPWP (K73-3) atau tidak.
- **K73-8** Rule yang hilang (S8): tombol Copy Coverage To All / For TPL disembunyikan sampai ekspornya ada; Choose Coverage
  memakai autocomplete Coverage Note.
- **K73-9** Premi: ProRatePercent dihitung ulang dari tanggal (CPCA langkah 4, Golf langkah 2) → uji Go Golf "pro-rata 0 → premi
  0" diubah; MinPremium Bond 100000 bila `Policy.Currency == 1` (arti kode = OQ-037/040/046, tidak ditebak).

## Frontend P1 — tab Object Aneka (sesi 0f, 09-10-2026)

Berkas: `components/TabObjectAneka.tsx` (grid lokasi → form lokasi → grid okupasi / Golf → grid objek), `components/FormObjekAneka.tsx`
(kontainer C1–C9), `components/aneka.ts` (fungsi murni), `api.ts` (`TampilanAneka`, `ObjekAneka` dkk., `ambilObjekAneka` /
`simpanObjekAneka`, `cariKapal`, `daftarTipeKendaraan`, `cariOkupasiAneka`), `labels.ts` (label sel Pega per blok), `nbfacin.css`
(`nbf-aneka-form`). Halaman case: `kasusAneka = !kasusFire && (lini.anekaNotBonding || lini.golf || lini.bonding)` — FIRE selalu
menang (regresi `lini.fire` 09-10-2026: layar FIRE = Group Business "FIRE" ATAU `lini.fire`).

- Label = label tampil sel Pega (urai `InputDtlObjectAneka_FacIn.xml`, `InputDtlObjectLocation_FacIn.xml`,
  `InputOkupasiAneka_FacIn.xml`, `InputDtlMemorandumAccept.xml`, HARNESS `ChooseOccupation.xml`). Sensus medan dua arah: spesifikasi
  `KONTAINER` 73 + Briguna 9 + Others 1 + tombol C2 `1=2` 1 = 84 = urai 84 sel FIELD (dikunci `FormObjekAneka.test.tsx`).
- Lookup kapal (pemetaan SET L10480-10570 diverifikasi sesi 0f), tipe kendaraan (disaring di layar), okupasi Aneka (id → OCCUPATION.ID).
  Object Name / Merk HE dan Draft Wording ID teks bebas (K73-10 / K73-11). `SetObjectDescHE_DT` L155 dihitung on-change di frontend.
- Keputusan agent: A-1 properti ganda lintas kontainer tampil sekali (kontainer pertama); A-2 sel `ALWAYS` selalu tampil (K73-5);
  A-3 Vessel Name (L10329) dan Object Description HE (L24736) `Read-only` TANPA syarat diperlakukan dapat diisi `[dugaan]` — **perlu
  konfirmasi layar Pega**; A-4 medan C2 `pyReadOnlyCondition=Always` tampil-saja (Type of Vessel / Length / Breadth / Depth / Place of
  Manufacture tidak di-SET lookup kapal → kosong, sama dengan Pega); A-5 tanda wajib = penanda, penolakan dari 422; B-1 kolom grid
  yang di Pega dapat diisi di baris ditampilkan sebagai teks, isiannya di form objek; B-2 Delete NB = buang baris.
- 422 proteksi ditampilkan sebagai daftar pesan (`pesanProteksiDari`, klien inti hanya membawa `galat`).
- Uji: `aneka.test.ts`, `FormObjekAneka.test.tsx`, `TabObjectAneka.test.tsx`; vitest 1019/1019, tsc bersih.
- Belum: P2 Coverage Aneka (kontrak menyusul ke sesi c3), P3 Spreading/Commission, P4 proteksi Submit.

## Frontend P2 — tab Coverage Aneka (sesi 0f, 09-10-2026)

Berkas: `components/TabCoverageAneka.tsx` (lokasi → okupasi / Golf → objek → daftar coverage → form; total per mata uang server),
`components/FormCoverageAneka.tsx` (sel `NB FacIn\Section\InputCoverageAneka_FacIn.xml`, kartu Coverage / Premium / Indemnity /
Limit + `GridDeductible`), `components/coverageAneka.ts` (operasi pohon murni), `api.ts` (`CoverageAneka`, `PemicuAneka`,
`cariCoverageAneka`, `hitungCoverageAneka`, `DaftarAneka.proRatePercent`, `TampilanAneka.currencyCoverage` / `mbd`). Halaman case:
tab Coverage `kasusAneka` di belakang FIRE.

- Hitung: setiap sel ber-aksi on-change memanggil `hitung-coverage-aneka` saat isian SELESAI (lepas fokus / Enter / pilih), membawa
  % Pro Rate case dan memasang nilai barunya; jawaban usang per coverage dibuang. Copy Coverage From This Section menghitung semua
  sasaran dalam satu permintaan (pemicu 'tsi').
- Operasi layar terverifikasi: Add (metode indeks keras L294-295, Currency/TSI objek, `setTSIAneka_preDT`), `SetTSI_DT` /
  `SetCurrency_DT`, Copy Coverage From This Section (`CopyCoverageFormSection_Act`: Section sama, TSI sasaran, Golf tidak ikut),
  Copy Deductible (`CopyDeductibleToAllLocationAneka_Act`: kode Coverage sama; `<CURRENT>` = iterasi coverage objek sumber,
  struktur diverifikasi sesi c3), `addDeductibleAneka_ACT` (currency deductible terakhir setiap coverage).
- Tab Object kini mengirim balik `proRatePercent` GET (PUT menulis kosong = NULL - tanpa ini Save tab Object menghapus % Pro Rate).
- Keputusan agent: C-1 % Pro Rate = nilai case, `fillPremiBond` membaca `.ProRatePercent` coverage (persis Pega); C-2 kartu tampilan;
  C-3 Choose Coverage = saran Coverage Note (K73-8); D-1 tombol ALL / TPL / Golf HIO / "Copy From" tidak dibuat (rule hanya di
  ekspor Endorsment / tidak ada); D-2 salin tersimpan lewat Save tab; D-3 total per mata uang = server, diperbarui sesudah Save.
- Uji: `coverageAneka.test.ts`, `FormCoverageAneka.test.tsx`, `TabObjectAneka.test.tsx` (bagian P2); vitest 1038/1038, tsc bersih,
  `go test ./modul/nbfacin/...` hijau (kunci kabel `currencyCoverage` / `mbd` dua sisi).
- Belum: P3 Spreading / Commission Aneka (catatan c3: ubah % Share belum menghitung ulang medan RNM coverage Aneka), P4 proteksi
  Submit (`ProtectCoverage_Act` langkah 5/6, `CheckSpreadingProtectAnekaGolf_ACT`).

### Ralat tata letak form coverage + arti `Read-only` (sesi 0f, 09-10-2026)

- Work owner 09-10-2026 (tangkapan layar Pega form coverage Liability, "SAMAIN KAYA GAMBAR2"): form coverage = SATU kisi 3 kolom
  urut sel Pega (kartu C-2 dibatalkan): Choose Coverage + Coverage Note · Calculate Method · % Pro Rate / % Short Period · Day ·
  [Currency] · TSI · % Loading · % Scale · % Rate · [Indemnity · Unit · % Indemnity] · Premium · % Discount · Discount · LoL · % LoL ·
  Sub Limit (lebar penuh) · garis · Deductible.
- ⭐ `[terverifikasi]` arti `pyEditOptions=Read-only` + `pyReadOnlyCondition`: syarat `IsClaim` = dapat diisi di NB; TANPA syarat =
  SELALU baca-saja (layar Pega: TSI L4069 tampil teks). Akibat: Coverage Note L1709 (diisi tombol Choose Coverage), Currency L3557,
  TSI L4069, % Indemnity L6149 = teks; Object Description HE form objek (L24736) = teks. Vessel Name (L10329) juga tanpa syarat,
  tetapi medan kapal C2 hanya terisi lewat pilihannya → tetap dapat dipilih, `[pertanyaan terbuka]` sampai ada layar Pega Marine Hull.
- Nilai awal sel (`pyDefaultValue`): % Short Period 0, Indemnity 0, % Discount 0 (selain Day 365, % Loading 0, % Scale 100,
  % Indemnity 100, Calculate Method 1).
- Choose Coverage: popup (kotak cari + ID / Coverage, sumber RD `BrowseCoverageFacIn_RD` Type ANEKA) - isi section
  `ChooseCov_Aneka` tidak ada di ekspor (K73-8), keputusan agent C-3.

### K73-16 Marine Hull: form DAN CSV (work owner 09-10-2026)

- Temuan `[terverifikasi]`: di form objek Pega seluruh medan kapal C2 baca-saja (Vessel Name L10329 tanpa syarat; Production
  Year / Place of Manufacture / Flag / `.MarineHull.*` `pyReadOnlyCondition=Always`). Penulis Type of Vessel / Place of
  Manufacture hanya `Activity\InsertUploadAneka_act.xml` + `UploadCSVAneka_PostAct.xml`; CSV mengisi ObjectName, Year, Flag,
  BuiltIn, TypeOfVessel, Classification, Construction, NRT/GRT/DWT, Currency, TSI, CoverageList - TIDAK Length/Breadth/Depth
  (cocok dengan tangkapan layar polis Marine Hull work owner: "--").
- Tata letak C2 = layout Pega (L9856-17613): kolom kiri Vessel Name / Currency / Production Year, kanan TSI / Place of
  Manufacture / Flag; blok MarineHull 3 kolom (Type of Vessel · NRT · Length | Classification · GRT · Breadth | Construction ·
  DWT · Depth) - kartu Vessel Details kini 3 kolom, urut per baris layar Pega.
- **K73-16** "HARUSNYA BISA DUA-DUANYA, YANG DI GAMBAR ITU YANG SUDAH JADI POLIS MAKANYA READ ONLY": Vessel Name tetap dapat dipilih
  dari daftar kapal (`GET /api/nbfacin/kapal`, BrowseShip_RD) di form, DAN Upload CSV objek Aneka dibangun di P6. Medan kapal lain
  tetap tampil-saja sesuai XML (terisi dari pilihan kapal / CSV).

## P3 — tab Spreading Aneka: spesifikasi (sesi 0f, 09-10-2026)

Riset subagen sesi 0f (laporan A–F di transkrip; skrip `sp3_*.py` di scratchpad). Diperiksa ulang sesi 0f:
`NB\Section\CoverageSpreadingList.xml` L3010 (`isGolfInsurance || IsAneka || IsBondingAndCustomBonds`); `DDL\GetKapasitasTreaty.xml`
cabang L10284 IsFire / L11453 isGolfInsurance / L12493 IsMBU / L21510 IsEngineering (grep IsAneka = 0);
`NB\Activity\CopyToAllSpreadingAnekaGolf_ACT.xml` L1098-1099 `local.TSI = @Math.divide(.TSI*Local.PercentShare,100,20)`;
`SumTSIPremiSpreadedRNM_ANEKA_Act.xml` L5624 / L6374 (`Local.IdxCoverage==1`); `services/deduction.go` / `cedant.go` memakai
pohon FIRE `KasusSpreading`.

- Layar: tab tidak bercabang per COB (% Share RNM, Copy Spreading, ringkasan sama dengan FIRE). Rantai Aneka: lokasi (No ·
  Location) → okupasi (Occupation Name) → objek (`SpreadingObjectAneka`: Object Name · Currency · TSI Old · TSI) → coverage
  (`SpreadingCoverageAneka`: Coverage · Gross Premium) → detail `SpreadingItemAneka` (Gross TSI · TSI Nusantara Re · Gross Premium
  · Premium Nusantara Re) + grid baris spreading = form FIRE `InputDtlSpreadingCoverage_FacIn(_IsUW)` apa adanya. Golf: lokasi →
  objek (Object Name) → coverage (`InputCoverageSpreadingAneka`: Coverage Note · % Loading · % Rate · TSI · Premi) → detail sama.
- Hitung: NR coverage `CountPremiAndTSIRNMFireAnekaGolf_ACT` langkah 7 (Aneka) / 8 (Golf) skala 20 (sudah ada `medanServerAneka`);
  DT `SetTSIPremiSpreaded_FacIn` percent (TSISpreaded / PremiumSpreaded dari TSI / Premi NR); Copy To All
  `CopyToAllSpreadingAnekaGolf_ACT` (`.TSI×PS`, tanpa TSIGross / FlagDelete); total `SumTSIPremiSpreadedRNM_ANEKA_Act` (TSI
  coverage ke-1 tiap objek, LoL → ClaimEstimation, top risk lini tertentu) / `_GOLF_Act`; ringkasan langkah 18–22 = FIRE.

### Keputusan work owner 09-10-2026

- **K73-17** Spreading otomatis ikut Pega: hanya Engineering (CAR / EAR / MBD / 10166 / 10032) dan Golf (1 baris ORS 10007 100%);
  lini Aneka lain lewat Copy / Add di tahap UW.
- **K73-18** TSI 100% Spreaded (TSIGrossSpreaded) Aneka kosong seperti Pega (kecuali spreading otomatis Golf / Engineering).
- **K73-19** Total per lokasi ikut Pega: tersembunyi untuk Aneka / Golf, tampil untuk MBD di tahap UW (`_IsUW` L17067).
- **K73-20** Kuirk Golf otomatis (TSI Spreaded dari TSILiability yang kosong → 0 sampai Save) ditiru apa adanya.

### Terbuka

- Commission Golf (`PropertyItemListCoverageCommision.xml` L9814: FA kelas Data-PropertyItem atas baris Data-Aneka) - perlu layar Pega.
- Payment / Premium Deduction / Cedant case Aneka: backend masih membaca pohon FIRE → perlu sumber per mata uang dari pohon Aneka
  (`.Policy.Payment.*` diisi `CountPremiAndTSIRNMFireAnekaGolf` 7.2.5 / 8.2.5) - dijadwalkan bersama P3.
- ~~Copy Spreading per lokasi `CopyToAllLocSpreading_ACT` (`_IsUW` L613): tombolnya ber-`pyCondition=1=2` → tidak diport.~~
  RALAT (K73-49): KELIRU - sel tombol ber-pyVisible ALWAYS (pyCondition tidak dievaluasi); digantikan K73-40.
- RALAT 09-10-2026 (work owner, gambar Choose Coverage: "ini balikin aja, sama seperti sblm nya"): tombol Choose Coverage + popup
  dibuang; Coverage Note kembali = isian dengan saran sambil mengetik (RD `BrowseCoverageFacIn_RD` Type ANEKA). Tata letak 3 kolom
  dan sel baca-saja lain (Currency / TSI / % Indemnity) tetap.

## Backend P3 — tab Spreading Aneka (sesi c3, 09-10-2026)

Riset korpus `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\` (NB = RNW byte-identik; Endorsment sama per langkah), baris kunci
dicocokkan ulang dengan `sed`. Catatan instrumen: `find` lintas korpus mengembalikan salinan `Endorsment Fac In` lebih dulu, jadi
rujukan baris selalu memakai path `NB FacIn\`.

**Endpoint** (kabel dikirim ke 9d sebelum uji dikunci):

- `GET`/`PUT /api/nbfacin/kasus/{caseId}/spreading-aneka`
- `POST …/spreading-aneka/hitung-share`
- `POST …/spreading-aneka/salin`
- `POST …/spreading-aneka/hitung-baris`
- `POST …/spreading/periksa-template` milik FIRE dipakai ulang apa adanya.

**Penyimpanan:**

- Tanpa migrasi baru.
- Baris spreading coverage Aneka = `T_SPREADINGLIST` ber-induk `T_COVERAGELIST` (migrasi 198), `SRC_PATH` jalur coverage +
  `/SpreadingList`.
- `CURRENCY_CODE` diambil dari coverage. Bila kosong dipakai mata uang objek, lalu `UNKNOWN`.
- Save tab Spreading Aneka menulis ulang pohon objek Aneka beserta coverage dan spreading lewat `GantiObjekAneka`, % Pro Rate
  tersimpan diteruskan, lalu `T_GENERAL_POLIS.PERCENT_SHARE` ditulis. Semuanya dalam satu transaksi dengan langkah 11 uw.
- Save tab Object Aneka mempertahankan spreading per posisi (jalur objek + indeks coverage, oldId sama; pola K48-7) dan
  menghitungnya ulang dari % Share RNM tersimpan.

**Rumus** (`services/spreadinganeka.go`):

- NR: `CountPremiAndTSIRNMFireAnekaGolf_ACT` L3854 / L3906.
  - TSI NR = PS × `.TSI` / 100, bukan TSILiability.
- DT `SetTSIPremiSpreaded_FacIn` percent (blok 8 / 9 / 10 / 13).
  - Hanya TSISpreaded dan PremiumSpreaded.
  - TSIGross dan ClaimEstimation tidak disentuh (grep = 0).
- Copy To All: `CopyToAllSpreadingAnekaGolf_ACT` L1098 / L3404 / L3450.
- Spreading otomatis, hanya di uw:
  - Golf: `GetKapasitasTreaty` DDL L10367–11168, ORS 10007 100 %.
  - Engineering: L12531–21510, kerangka FIRE dengan lima beda verbatim:
    1. template tidak dibuang per lokasi;
    2. TSI top risk tanpa kurs;
    3. TSIRNMTopRisk1 = 0 (`[dugaan]` Local tak diset);
    4. share QS dipotong = MaxLimitQSIDR / TSISpreaded sebelum dipotong (L19128 / L19216, CARID1);
    5. SPL ditambahkan tanpa syarat TSIRNM > maxtsi.
- Total per lokasi: `SumTSIPremiSpreadedRNM_ANEKA_Act` cabang 2.4 (L9275) / 2.5 (L15995) dan `_GOLF_Act`.
- Total case, ringkasan, batas treaty, dan jenis treaty: fungsi FIRE (dispatcher langkah 17–22 tidak bergerbang lini).
- Daftar When lini dicek mentah dengan `sed`:
  - ANEKA L4100–4360 / L10400–10700 / L12440–12600;
  - Count L2540–4620.

**Payment / Premium Deduction / Cedant:**

- `bacaSpreadingRingkas` memakai pohon FIRE apa adanya bila ada item.
- Bila tidak ada item, dipakai adaptor pohon Aneka (Count 7.2.5 / 8.2.5):
  - lini Count-khusus: satu item per coverage (mata uang coverage);
  - lini lain: satu item per objek.
- Kabel tidak berubah.

**Penyimpangan sadar `[keputusan agent]`** (dilaporkan ke work owner lewat 9d, menunggu jawaban):

1. Kuirk multi-mata-uang tidak ditiru. Kunci total = mata uang coverage untuk lini daftar (kosong → objek). Kuirk Pega yang
   dimaksud:
   - Golf CurrObj tanpa reset;
   - APPEND total Golf tanpa penjaga mata uang (L5797);
   - Aneka CurrObj lintas coverage;
   - `Local.totaltsi` kumulatif lintas mata uang.
2. ~~Langkah 24 dispatcher (IsMBD, Page-Clear-Messages per lokasi, L8926–10186) tidak di-port.~~ RALAT: di-port (K73-22).
3. Engineering non-IDR:
   - pembilang share QS dipotong = MaxLimitQSIDR tanpa kurs (verbatim), sehingga share bisa > 100;
   - QS hasil pecah tepat di batas dapat ikut dipotong karena sisa presisi skala 20 (dicatat dari uji; Pega `@divide` skala sama
     `[dugaan]` berperilaku sama).
4. BuilderRisk / PortRisk tidak menyumbang TSI ke Payment / Cedant — verbatim, ditiru.
5. Tanpa BUSINESS case terbaca (lini tidak diketahui), DT / total tidak berjalan (Pega pun bergantung pada When lini).

**Uji:**

- `services/spreadinganeka_test.go` dengan mutasi:
  - 14 mutasi rumus: 13 mati; yang hidup 1 setara logika (penjaga `top` berlebih);
  - 4 mutasi Engineering: 3 mati, 1 no-op.
- `TestLiniSpreadingAneka` = uji instrumen When atas lini yang jawabannya diketahui.
- `handlers/spreadinganeka_test.go` mengunci kunci kabel.
- `repository/aneka_test.go` `TestSisipSpreadingAneka`.
- `services/aneka_test.go` `TestPertahankanSpreadingAneka`.

**Tetap terbuka:** Commission Golf (`PropertyItemListCoverageCommision.xml` L9814). P4 (proteksi Submit) menunggu kontrak.

### Keputusan work owner 09-10-2026 (P3, lewat sesi 9d)

- **K73-21** Multi mata uang: SETUJU. Kuirk Pega tidak ditiru; total dikelompokkan per mata uang dengan benar (penyimpangan 1
  di atas).
- **K73-22** Langkah 24 SumTSIPremiSpreadedRNM_Act (IsMBD): PORT APA ADANYA. Sudah di-port sebagai `hapusPesanMBD`:
  - bila ada total per lokasi yang TSISpreaded-nya (mata uang asli, tanpa kurs) ≥ limit GetTreatyName berjenis sama, SELURUH
    pesan yang terkumpul dihapus;
  - limit ORS = 150 M, atau 181,5 M untuk Begin date ≥ 1-7-2026;
  - bukti `NB FacIn\Activity\SumTSIPremiSpreadedRNM_Act.xml` L9002 / L9049 / L9163 / L9307 / L9849 / L9895 / L10014 / L10069 /
    L10396 / L10448.

  Tidak di-port: 24.4 group / pyID hardcode (pola K48-6), dan `ProtectSpreading.CARI1 = 0` (L10186, penanda yang tidak
  dimodelkan). Limit yang tidak terurai sebagai angka tidak dibandingkan `[dugaan]`. Uji: `TestHapusPesanMBD`,
  `TestPesanMBDDihapus`.
- **K73-23** Engineering top risk (pembagi share QS = `Local.TSIRNMTopRisk1` yang hanya diisi cabang Fire): PERBAIKI, pilihan
  **A**, sudah diterapkan.
  - Fase A menghitung `tsitoprisk += TSILiability` (× kurs bila bukan IDR, padanan FIRE `DDL\GetKapasitasTreaty.xml` L2441) dan
    `TSIRNMTopRisk1 += TSILiability (× kurs) × PS / 100` (L2494), untuk coverage pertama lokasi top risk (WHEN L2595 / L13933).
  - Jalurnya OccupationList → AnekaList. Kurs = GetCurrencyToIDR_SQL yang sudah dipanggil, jadi tanpa sumber baru.
  - L16245 tetap. Verbatim Engineering L13793 / L13839 (tanpa kurs) diganti.
  - Uji: `TestSpreadingOtomatisEngineeringTopRisk`, `TestSpreadingOtomatisEngineeringTopRiskKurs` (TSI pencarian kapasitas 1000
    IDR).
- **K73-24** BuilderRisk / PortRisk tidak menyumbang TSI ke Payment / Cedant: SETUJU, ditiru (sudah begitu).
- **K73-25** Engineering non-IDR (pembilang share QS dipotong = MaxLimitQSIDR tanpa kurs, L19216 / L19237): PERBAIKI, pilihan
  **P1**, sudah diterapkan.
  - `share = Local.maxtsiQS / CARID1 × 100`, dengan maxtsiQS dari L16086 (IDR) / L16384 (÷ kurs). Hasil IDR tidak berubah.
  - Uji: `TestSpreadingOtomatisEngineeringKurs` (share 75, verbatim Pega 150).
- **K73-26** Sisa presisi skala 20 memicu pemotongan QS (IDR pun kena): pilihan **(c)**, sudah diterapkan.
  - Pemicu pemotongan .19.3 (`.TSISpreaded > Local.maxtsiQS`, L19342) membandingkan KEDUA nilai sesudah dibulatkan 8 desimal.
  - Pemicu top risk .20.3 (L20772) dan rumus lain tetap verbatim.
  - Uji: `TestSpreadingOtomatisEngineeringSisaPresisi`, `TestSpreadingOtomatisEngineeringBatasDibulatkan`.
  - Hasil kasus TSI RNM 700 / maxQS 300 / maxSPL 400 dengan K73-26 saja: QS 42,85714286 % / 300 tanpa pemotongan ulang.
    - Namun baris SPL template (57,142857 % → TSI 400 > maxQS 300) tetap dipotong oleh .19.3, karena WHEN L19342 tidak
      menyaring jenis treaty: hasilnya SPL 75 % / 300.
    - "Beda 5" lalu menambah SPL 25 % / 175. Total TSI 775 dari 700.
    - Itu sebabnya K73-27 diputuskan.
- **K73-27** Pemicu pemotongan .19.3 Engineering (WHEN L19342 `.TSISpreaded>Local.maxtsiQS`) hanya untuk baris berjenis treaty
  QS (`TreatyType = IDTreatyQS`). Pilihan **(ii)**, keputusan work owner 09-10-2026 lewat 9d. **MENYIMPANG dari Pega**:
  - Alasannya, Pega membandingkan baris SPL dengan limit QS lalu memotongnya.
  - RALAT dasar: usulan sesi c3 sempat mengutip keterangan langkah Pega "JIKA TSI SPREADED QS ..." sebagai dasar. Kutipan itu
    terpotong. Keterangan utuhnya "JIKA TSI SPREADED QS DAN SPL LEBIH BESAR DARI LIMIT SPREADING" (`DDL\GetKapasitasTreaty.xml`
    L6513 / L7226 / L17611 / L18327), dan milik langkah pecah .15 / .18, bukan .19.3 (keterangan L19120 kosong). Jadi dasar
    K73-27 adalah keputusan work owner, bukan keterangan Pega.
  - Cakupan: hanya `spreadingOtomatisEngineering`. FIRE `spreadingOtomatis` TIDAK diubah. Pemicu top risk .20.3 tetap verbatim.
  - Hasil kasus 700 / 300 / 400: QS 42,85714286 % / 300 + SPL 57,14285714 % / 400 = 700, tanpa baris tambahan.
  - Uji `TestSpreadingOtomatisEngineeringSisaPresisi`. Mutasi K73-23 / 25 / 26 / 27: 14/14 mati.
- Catatan (tanpa perubahan, keputusan work owner 09-10-2026): FIRE `spreadingOtomatis` (tiket 48, sudah berjalan) juga memotong
  baris SPL hasil pecah di maxQS (.19.3 berlaku untuk semua baris). Contoh 700 / 300 / 400 memberi SPL 42,857 % / 300 (total 600
  dari 700). Port verbatim; jalur FIRE TIDAK diubah.

## P4 — Submit case keluarga Aneka (sesi c3, 09-10-2026)

Riset `[terverifikasi]` (baris kunci dicek ulang dengan `sed`): laporan sesi c3 `riset67\p4_a.md` (proteksi) dan
`riset67\p4_b.md` (jenjang limit + tulisan), di scratchpad sesi, bukan di repo.

Temuan inti:

- Port Go memblokir Submit Aneka:
  - `proteksiBerflag` tanpa gerbang lini;
  - `submituw.go` meng-hardcode BusinessType "Fire".
- Proteksi spreading tahap uw hanya membaca pohon FIRE.
- Pega menjalankan CheckSpreadingProtect di Binding; port Go tidak.
- Jenjang limit Aneka / Golf memakai `M_LIMIT_NONPROPANDENGG`, Engineering `M_LIMIT_ENGINEERINGG`, Bond / Kredit (ActFlow)
  `M_LIMIT_FINANCIALINS`. Tidak ada `M_LIMIT_ANEKA`.
- FACINOFFER bercabang Aneka / Bonding / Golf (`SaveOfferProduction_Act`).

### Keputusan work owner 09-10-2026 (P4, lewat sesi 9d)

- **K73-28** Proteksi FIRE (lokasi / okupasi / Loss Record / Province-City-ZipCode, Scoring Risk, ProtectFIREMBUPA) hanya
  `IsFire`, sesuai Pega (`Protection_Act` L3873 dll.).
- **K73-29** CheckSpreadingProtect di tahap Binding untuk SEMUA COB, termasuk FIRE yang sudah berjalan (perubahan perilaku FIRE
  DISETUJUI; uji regresi FIRE wajib).
- **K73-30** `IsBedaCurr` Aneka = ada deductible yang Currency-nya beda dari Currency coverage-nya. `[dugaan]` tafsiran agent
  karena `SaveDeductibleAneka` tidak ada di ekspor; diganti bila XML-nya datang.
- **K73-31** Kasus ganda (GetDataDoubleCase): case sendiri DIKECUALIKAN (menyimpang dari Pega).
- **K73-32** Pesan Golf: nomor objek diisi dengan benar (Pega memakai `Local.idxobjitem` yang tidak di-set di cabang Golf;
  menyimpang); teks lain verbatim.
- **K73-33** TSI top risk untuk jenjang limit MEMAKAI KURS, untuk FIRE juga (ikut Pega `CountTotalTSIPremiNusaRe_Act`
  L9060 / L9254; port tiket 61 diubah; uji regresi FIRE).
- **K73-34** Ikut Pega: Bond / Kredit di JUW / UW memakai NONPROPANDENGG; Senior ke atas memakai NextFinancial
  (`M_LIMIT_FINANCIALINS`).
- Butir DBA: isi / ejaan JABATAN di NONPROPANDENGG / ENGINEERINGG / FINANCIALINS. Uji memakai data sintetis.
- **K73-35** FACINOFFER Bonding satu baris saja (lewat cabang Bonding).
- **K73-36** IsFacRetro Golf memeriksa SEMUA lokasi (Pega hanya lokasi terakhir; menyimpang).
- **K73-37** (P5) IsFlagUW Aneka = 1 seperti Pega: TL Policy Accept naik ke UW Policy.
- **K73-38** Pembulatan proteksi spreading: Aneka skala 4, Golf skala 2, HALF_UP.

Cakupan P4:

1. gerbang lini proteksi FIRE;
2. ProtectCoverage cabang Aneka / Golf;
3. proteksi spreading Aneka / Golf di uw, plus Binding untuk semua COB;
4. jenjang limit Aneka / Golf / Engineering / Bond;
5. FACINOFFER.

P5: FACINPRODUCTION, M_RISK_LOSS_PROFILE, IsFlagUW, varian 01-01-95.

### K73-39 TSI Old hanya EDM (work owner 09-10-2026)

- "tsi old hanya untuk EDM, di nb,rnw tidak perlu. jangan halu": kolom TSI Old dibuang dari grid objek tab Object (Golf), Spreading
  dan Premium Deduction Aneka. ⚠️ MENYIMPANG dari XML: `NB FacIn\Section\SpreadingObjectAneka.xml` L1584 dan
  `CommisionObjectAneka.xml` L1594 menampilkan kolomnya tanpa syarat (pyMemo "Tambah column TSI Old"); TSI Old hanya terisi
  endorsement, sehingga di NB / RNW selalu kosong. Kabel `tsiOld` tetap ada (dipakai EDM nanti), hanya tidak ditampilkan.
- Detail angka coverage (Gross TSI / TSI Liability, TSI Nusantara Re, Gross Premium, Premium Nusantara Re) di tab Spreading dan
  Premium Deduction Aneka = empat kartu kecil berlabel (perapian tampilan).

### Backend P4 — dibangun (sesi c3, 09-10-2026)

Kabel `…/submit`, `…/submit/periksa`, `…/submit/validitas` TIDAK berubah. Tanpa migrasi. Case keluarga Aneka = lini terbaca
dengan IsAneka (termasuk Bonding / Engineering) atau isGolfInsurance. Bila lini tak terbaca, perilaku FIRE tidak berubah.

1. **Gerbang lini (K73-28)** — `services/proteksi.go` `proteksiBerflag`: untuk case Aneka, Scoring Risk (L3873),
   ProtectCurrencyTSI dan ProtectFIREMBUPA (lokasi / okupasi / Loss Record / Province-City-ZipCode) tidak dijalankan. TL
   Deductible Acceptance memakai `deductibleBedaMataUangAneka` (K73-30 `[dugaan]`).
2. **ProtectCoverage cabang Aneka / Golf** — `services/proteksisubmitaneka.go` `proteksiCoverageAneka`: pesan verbatim
   L3779 / L4229 / L4529 / L5030 / L5482 / L6147 / L6920 / L7281, termasuk spasi depan dan dua spasi sebelum "DeductibleList".
   - Deductible hanya tahap Binding (flag 2, EmailTypeBinding bukan 4) dan bukan IsTBonding.
   - Kasus ganda: `repository/kasusganda.go` (GetDataDoubleCase verbatim + case sendiri dikecualikan, K73-31) per lokasi;
     ditemukan -> lompat label AA (lokasi berikut dan blok Golf dilewati).
3. **Proteksi spreading** — `pesanSpreadingAneka` (CheckSpreadingProtectAnekaGolf: Aneka skala 4, "object item" = nomor
   okupasi verbatim; Golf skala 2, nomor objek diisi K73-32) di `proteksiSpreadingUW` (tahap uw / UW Policy) sesudah cabang FIRE;
   M4 / M5 juga dari ringkasan spreading Aneka.
   - Tahap Binding (K73-29, SEMUA COB): `proteksiSpreadingBinding` di `proteksiTahapBinding` dan Banding Binding, pesan ber-flag
     ProtectSpreading.CARI1 (`[dugaan]` untuk M4 / M5).
   - Fixture uji Binding FIRE kini memakai spreading sah; regresi FIRE: `TestSpreadingBindingFire`.
4. **Jenjang limit** — `services/submituw.go`:
   - `kasusTangga` membaca BusinessType / BusinessCode / BusinessName case Aneka;
   - JUW B / JUW A / UW: `M_LIMIT_ENGINEERINGG` (IsEngineering) / `M_LIMIT_NONPROPANDENGG` (selain itu), tanpa IsPreferredRisk;
   - TSI dari pohon Aneka (`barisTSIAneka`, CountTotalTSIPremiNusaRe_Act 6 / 7);
   - TotalTSITopRisk BERKURS untuk semua COB (K73-33; `TestCarid2UW` diubah: top risk USD 7500200, kurs tidak ada -> ditahan);
   - IsFacRetro dari spreading Aneka / Golf semua lokasi (K73-36);
   - ErrBentukB -> `acceptance.NextFinancial` (`M_LIMIT_FINANCIALINS`, K73-34; uji "bentuk B -> ditahan" diubah: FINANCIALINS
     kosong -> Binding; `TestAcceptBentukBFinancial`).
5. **FACINOFFER** — `services/bindinganeka.go`: cabang 13 / 14 / 15 (lihat kepala berkas).
   - Bonding satu baris (K73-35).
   - NOSERIAL kini ber-bind (`repository/binding.go`; FIRE tetap NULL karena "" = NULL; `TestSQLBinding` 10 kolom NULL / 40 bind).

**Penyimpangan / tafsiran `[keputusan agent]`:**

- Currency coverage kosong = Currency objek (pesan L4229, IsBedaCurr, FACINOFFER MBD). Layar coverage Aneka hanya menampilkan
  Currency untuk lini daftar; Count L2969 mengisinya dari objek.
- Model Aneka hanya menyimpan CURRENCY_CODE, jadi "Currency.ID / OldID kosong" diperlakukan sebagai kode kosong.
- Golf objek IDR = kurs 1 (Pega hanya memeriksa Currency coverage "IDR", L4538).
- Langkah ORS-only (S7, "Can not proceed spreading ORS") dan cabang Group / user-ID tidak diport, sama dengan FIRE (M6 / K48-6).
- `.IdxLocation` / `.IndexLocation` pesan objek = nomor lokasi `[dugaan]` (SetIndexObject_DT).

**Uji** (`services/proteksisubmitaneka_test.go`, `services/bindinganeka_test.go`, `services/submituw_test.go`,
`repository/binding_test.go`). Mutasi:

- butir 1–4: 27 mutasi, 26 mati. Yang hidup: ProtectCurrencyTSI tetap jalan untuk Aneka; setara, karena tanpa item FIRE tidak
  ada pesan;
- FACINOFFER: 6/6 mati;
- K67-6: 2/2 mati.

**Persetujuan work owner 10-10-2026 (lewat 9d):** P4 disetujui (verifikasi sesi 0f: build, vet, go test, vitest 1051/1051).
Tafsiran `[keputusan agent]` berikut DITERIMA:

- (a) Currency coverage kosong = Currency objek.
- (b) M4 / M5 di Binding ber-flag ProtectSpreading.CARI1.
- (c) Golf objek IDR = kurs 1 untuk TSI limit.
- (d) ORS-only (S7) dan cabang Group / user-ID tidak diport.
- (e) Isi / ejaan JABATAN tabel limit Aneka = butir DBA; uji memakai data sintetis.

## P5 — selesai Policy case keluarga Aneka (sesi c3, 10-10-2026)

Riset: `riset67\p5_a.md` (FACINPRODUCTION, `SaveTreatyProduction_Act` cabang 19 Aneka / 21 Golf) dan `riset67\p5_b.md`
(M_RISK_LOSS_PROFILE, IsFlagUW, tulisan Policy lain), di scratchpad sesi c3.

### Keputusan work owner 10-10-2026 (P5, lewat sesi 9d)

- **K73-40** Tombol "Copy To All Spreading" per lokasi (`PropertyItemListCoverageSpreading_IsUW.xml`, pyVisible ALWAYS) DIPORT
  untuk SEMUA COB, termasuk FIRE (perubahan perilaku FIRE disetujui; uji regresi FIRE wajib). Termasuk efek
  `CopyToAllLocSpreading_ACT` L532 (IsFlagUW = 0 bila kosong) dan cabang Aneka langkah 6 / Golf langkah 7. Kabel baru
  `…/spreading/salin-lokasi` dan `…/spreading-aneka/salin-lokasi`; usulan kabel dikirim ke 9d sebelum dikunci.
- **K73-41** Polis deklarasi (PolicyType 2) keluarga Aneka: IsFlagUW 0, jadi selesai setelah TL Policy Accept (Pega
  CekLimitSpreading L2407).
- **K73-42** OCCUPATIONCODE FACINPRODUCTION Aneka = NAMA okupasi per baris (Pega L4614). FIRE tetap.
- **K73-43** TOP_RISK / IDX_OBJECT_ITEM / TSI100_* / TYPE_EDM Aneka / Golf = NULL (Pega tidak men-SET).
- **K73-44** Kota / kecamatan / kelurahan Golf = ASMCity / ASMDistrict / ASMRW seperti Aneka (Pega Golf: CityName / DistrictName
  / RWName, tidak disimpan). Tanpa migrasi.
- **K73-45** DEDUCTION2 Aneka = NULL seperti port FIRE.
- **K73-46** MARKETINGOFFICERCODE Aneka = MOID case (Pega S11 L2492). FIRE tetap GetMOID.
- **K73-47** Tulis-backup Golf (InsertFacinProductionBackup_Sql) TIDAK ditiru.
- **K73-48** M_RISK_LOSS_PROFILE Aneka: premi 0 ikut Pega (sumber per kolom: baris Local di SaveJsonPolicyFacIn_Act, lihat
  "Backend P5").
- **K73-49** Ralat: catatan "Copy Spreading per lokasi … tidak diport" keliru, digantikan K73-40.
- **K73-50** (salin-lokasi M1) Efek IsFlagUW tombol per lokasi (CopyToAllLocSpreading_ACT L532 "0 bila kosong") TIDAK diport:
  IsFlagUW tidak disimpan sistem baru, tanpa migrasi. Aneka tetap 1 (K73-37) kecuali polis deklarasi (K73-41); tombol hanya
  menyalin baris.
- **K73-51** (salin-lokasi M2) `flagSpreading` (Type Facultative) diterima dan divalidasi (wajib bila treatyType "10015":
  `PropertyItemListCoverageSpreading_IsUW.xml` sel `.FlagSpreading` L2934-3289, `pyRequiredWhen=.TreatyType=="10015"` L3020,
  `pyDisabledWhen=.TreatyType!="10015"` L3014 - bukti sesi 0f), TIDAK disimpan (sama dengan salin case-level). Opsi FlagSpreading
  ("associated as defined on property") tidak ada di ekspor - teks bebas, tidak ditebak.
- Kabel salin-lokasi (diterima 9d / 0f): `POST …/spreading/salin-lokasi`, `POST …/spreading-aneka/salin-lokasi`, badan
  `{percentShare, idxLokasi (1-based), template:[{treatyType, treatyName, sharePercentage, flagSpreading}]}`, jawaban = Tampilan
  salin dengan `template: []`.

Cakupan P5: IsFlagUW keluarga Aneka (K73-37 / K73-41), FACINPRODUCTION cabang 19 / 21, M_RISK_LOSS_PROFILE (K73-48), tombol
Copy To All per lokasi (K73-40). Tanpa migrasi.

### Backend P5 — dibangun (sesi c3, 10-10-2026)

Kabel `…/submit` tidak berubah. Kabel baru (K73-40): `…/spreading/salin-lokasi` dan `…/spreading-aneka/salin-lokasi`. Tanpa
migrasi.

1. **IsFlagUW** — `services/submitpolicy.go` `isFlagUWAneka`: keluarga Aneka (IsAneka ∨ IsBondingAndCustomBonds ∨
   isGolfInsurance) = 1, kecuali polis deklarasi (PolicyType "2" = When IsDeclarationPolicy L137) = 0 (K73-37 / K73-41). FIRE
   (lini tak terbaca / bukan Aneka) tetap K69-5.
2. **FACINPRODUCTION** — `services/selesaipolicyaneka.go` `barisFacInProductionAneka` (cabang 19 / 21, lihat kepala berkas).
   - `dasarBarisProduksi` dipakai bersama cabang FIRE.
   - Repository: MARKETINGOFFICERCODE = `NVL(MOID case, GetMOID)`. FIRE mengirim NULL, jadi hasilnya sama.
   - NOSERIAL ber-bind; kosong = teks "NULL" (FIRE tetap). `TestSQLSelesaiPolicy`: 61 bind, 2 teks 'NULL'.
3. **M_RISK_LOSS_PROFILE** Aneka — `profilRisikoLossAneka`: TSI NULL; premi / OR / QS / FACOUT 0 (K73-48).
   - Tipe Local `SaveJsonPolicyFacIn_Act`: TSINusare String L437 / L439; PremiNusare / PremiumQs / PremiumOR / HasilQs /
     PremiumFacout double L443-469.
   - FIRE tidak berubah.
4. **Copy To All per lokasi** — `services/salinlokasi.go`, `handlers/salinlokasi.go`:
   - FIRE = salinTemplate, Aneka / Golf = salinTemplateAneka, keduanya dibatasi lokasi idxLokasi;
   - penjaga / validasi = salin case-level + flagSpreading wajib bila "10015" (K73-51) + idxLokasi 1..n;
   - tidak menyimpan; tanpa pesan jenis treaty (Pega tidak memanggil cekSpreadingFactIn);
   - efek IsFlagUW tidak diport (K73-50).
   - `ResetFlagSpreading_ACT` (`ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE!RESETFLAGSPREADING_ACT`, dipanggil onchange
     TreatyType): Property-Remove `.FlagSpreading` SEMUA baris template per lokasi (L290) + `OfferFacIn.FlagSpreading` (L321).
     Ini status layar, jadi efeknya ada di frontend.
5. **Ralat P3** — `kasusSpreadingAneka` (sumber Payment / Deduction / Cedant) kini menggerbang objek okupasi dengan IsAneka (Count
   langkah 7 L5462) dan objek Golf dengan isGolfInsurance (langkah 8 L7763). Sebelumnya objek Golf selalu ikut.

**Uji:** `services/selesaipolicyaneka_test.go`, `services/salinlokasi_test.go`, `handlers/salinlokasi_test.go`; uji P3
`TestKasusSpreadingAneka` / `TestBacaSpreadingRingkas` disesuaikan (gerbang lini). Mutasi P5: 19 mutasi + 1 sesudah uji Golf
MBD ditambah, semua mati.

### Berkas work owner 10-10-2026 (`D:\migrasi\RNM\DDL\`)

- **FlagSpreading.xml** (`ASM-FW-GISFW-DATA-SPREADINGRISK!FLAGSPREADING`): nilai lokal "SF-MR 2023" dan "Open Fac Market" (L186-192;
  sensus dua cara grep 2 / urai XML 2). Validasi salin-lokasi kini membatasi flagSpreading terisi ke dua nilai itu (peka huruf;
  400 di luar daftar). Wajib bila treatyType "10015" (K73-51) tetap.
- **M_LIMIT_NONPROPANDENGG.xls / M_LIMIT_ENGINEERINGG.xls / M_LIMIT_FINANCIALINS.xls** (butir DBA P4). Dibaca lewat Excel COM
  read-only. Hanya JABATAN / TEAM_GROUP / LIMIT_BOTTOM* / WORKBASKET / JABATAN_ATASAN yang diambil (NAMA / LOGIN tidak disalin).
  Isi tetap dibaca dari DB, tidak di-hardcode.
  - Kolom kedua tabel bentuk A sama dengan DDL (14 kolom; 37 baris masing-masing). FINANCIALINS 7 kolom, 4 baris.
  - Ejaan JABATAN NONPROPANDENGG / ENGINEERINGG = ejaan `jabatanPengguna` / `jabatanPosisiUW` Go (JUW_B, UNDERWRITER, SENIORUW,
    DEPHEADUNDERWRITER, KADIVFACULTATIVE, KADIVTEKNIK, MANAGERTEKNIK, DIREKTURMARKETING, DIREKTURTEKNIK) + LEADER (JUW_UW langkah 7):
    tanpa selisih ejaan.
  - TIDAK ADA baris JUW_A dan KADIVFINANCIAL di kedua tabel. UNDERWRITER tidak ada untuk team 3 (sama dengan PROPERTYY DDL).
    DEPHEADUNDERWRITER hanya team 3. Akibatnya Accept JUW A / Kadiv Financial / UNDERWRITER team 3 = subquery tanpa baris ->
    LetterNo "" -> Else (perilaku Pega yang sama, bukan selisih).
  - SENIORUW team 1 dua baris (limit identik, LOGIN berbeda -> Pega tidak ORA-01427). `acceptance.Next` menerima limit identik.
  - FINANCIALINS JABATAN berspasi: "SENIOR UNDERWRITER", "KADIV KEUANGAN", "DIREKTUR MARKETING", "DIREKTUR TEKNIK". Ini cocok
    dengan gerbang `acceptance.NextFinancial` langkah 22 ("KADIV KEUANGAN", "DIREKTUR MARKETING", "DIREKTUR TEKNIK").

## Bug DEV 10-10-2026 — Save Spreading Aneka 500 `ORA-01722` di `T_SPREADINGLIST` (sesi c3)

- `[terverifikasi: go-ora v2.8.19, module cache]` INSERT array-bind (`command.go` `Stmt._exec` cabang `bulkExec`) menyandi
  tiap elemen kolom dengan SATU `ParameterInfo` yang dipakai ulang: `encodeValue` → `par.init()` (`parameter_encode.go` L516,
  CharsetID / CharsetForm / ContFlag = 0) → `setDataType` menurut tipe Go elemen (nil → NCHAR tanpa charset L16-18; string →
  NCHAR + charset server L98-102). DataType kolom = elemen non-kosong terakhir, tetapi CharsetID / CharsetForm di kepala bind
  (`parameter.go` `write` L373-374) = milik elemen TERAKHIR. Kolom teks / desimal (angkaMasuk) yang baris terakhirnya NULL
  dan baris sebelumnya berisi terkirim NCHAR ber-charset 0. `[dugaan]` dampaknya di server = ORA-01722 di `TO_NUMBER` (tidak
  dapat dibuktikan dari kode klien; tanpa percobaan ke DEV). Baris spreading Aneka lazim begitu (mis. baris OR berisi, baris
  terakhir tanpa TreatyName / ClaimEstimation). Penyebab format nilai (>8 desimal) disingkirkan: spreading lewat
  `bulatkanSimpanAneka` (`simpan8`).
- Perbaikan (opsi b): `repository.kolomArrayBertipe` — kolom teks yang memuat NULL diikat `sql.NullString` di SETIAP baris
  (satu tipe Go per kolom; cabang `tyNullString` = cabang string); kolom bercampur tipe / tipe bukan teks bercampur NULL =
  galat. Dipakai `sisipAnekaBerkelompok` (seluruh tabel jalur Aneka: lokasi, okupasi, objek, coverage, deductible, spreading).
  `[keputusan work owner 10-10-2026, lewat sesi 9d]` FIRE ikut: `sisipBerkelompok` (Save tab Object FIRE) memakai
  `kolomArrayBertipe` sebagai pencegahan — pengecualian sadar dari "jangan ubah yang sudah berhasil"; isi data yang disimpan
  tidak berubah (`bandingkanJalur` membandingkan dengan jalur lama per baris setelah `sql.NullString` dikembalikan ke nilai
  asal). `[terverifikasi]` fixture `objekKecilUji` (jalur lama) memuat 5 kolom NULL di baris terakhir setelah baris berisi,
  termasuk T_LOCATIONLIST kolom 5 (log `TestSisipBerkelompokFireNullTerakhir`).
- Uji: `TestKolomArrayBertipe` (instrumen `periksaKolomSeragam` diuji dulu atas jawaban yang diketahui), tiga uji
  `sisipAnekaBerkelompok` memakai `eksekutorTiruan{seragam: true}` (satu tipe per kolom, tanpa nil telanjang, teks bind
  angkaMasuk cocok `PolaAngkaMasuk`); mutasi `kolomArray` lama di jalur Aneka → ketiganya gagal.
- Uji FIRE: `TestSisipBerkelompokFireNullTerakhir` (bentuk lama T_LOCATIONLIST ditolak instrumen; jalur FIRE lolos
  `periksaKolomSeragam`), `bandingkanJalur` / urutan eksekusi memakai `seragam: true`; mutasi `kolomArray` lama di
  `sisipBerkelompok` → `TestSisipBerkelompokFireNullTerakhir`, `TestSisipBerkelompokSetara`, `TestSisipBerkelompokRoundTrip`
  gagal ("nil telanjang").
- **Lanjutan 10-10-2026 (biner DEV dibangun ulang 11:18, ORA-01722 tetap):** `kolomArrayBertipe` bukan (satu-satunya)
  penyebab. Perubahan:
  - Coverage / deductible / spreading jalur Aneka kini disisip PER BARIS dengan bind skalar, bentuk yang jalan di DEV untuk
    spreading FIRE (`tabelPerBarisAneka`, `jalankanKelompok`), tetap setelah induknya.
  - Setiap INSERT objek (Aneka, FIRE berkelompok, spreading FIRE per baris) diperiksa sebelum Oracle oleh `periksaBindBaris`
    (`repository/periksabind.go`):
    - bind angkaMasuk wajib cocok `PolaAngkaMasuk`;
    - ID / PARENT_ID / SEQ_NO wajib bilangan bulat (kecuali T_LOCATIONLIST.PARENT_ID VARCHAR2 ID case);
    - jumlah nilai = jumlah bind.
    Galat `ErrBindTidakSah` menyebut tabel, kolom, baris, dan BENTUK nilai (tanpa nilai mentah).
  - Galat Oracle dibungkus tabel + rentang baris + jumlah bind + cara (array-bind / per baris).
- Uji jalur PUT utuh (`TestSimpanSpreadingAnekaBindCaseDEV`, service → `repository.SimulasiTulisAneka`) atas rekonstruksi case
  work owner: 2 lokasi, PUBLIC LIABILITY, % Share 25, Copy To All "2024 QS 155M FAC" 100 %, TSI 500.000.000.000.000.
  `[terverifikasi di uji]` semua bind T_SPREADINGLIST lolos mask; total = ringkasan layar (TSI 137.500.000.000.000, premi
  1.375.000.000.000). Jadi isi bind bukan penyebab untuk bentuk ini `[dugaan]`.
- **Akar masalah (DEV 11:31, ditangkap periksaBindBaris: `T_SPREADINGLIST kolom CLAIM_ESTIMATION baris 1: panjang bulat 1,
  desimal 20`; ditelusur sesi 0f, diperbaiki sesi c3):**
  - Urutan di cabang uw `SimpanSpreadingAneka`: `hitungSaveSpreadingAneka` membulatkan 8 desimal, lalu `totalAneka` dipanggil
    untuk `periksaJenisTreaty`.
  - `totalLokasiAneka` di dalamnya menimpa `ClaimEstimation` baris dengan nilai skala 20 yang belum dibulatkan (LoL kosong →
    0E-20), dan pohon itulah yang ditulis.
  - Perbaikan: pemeriksaan memakai `salinDalamAneka(d.lokasi)`. Hasil hitung tidak berubah: pohon yang ditulis tetap hasil
    `hitungSaveSpreadingAneka`, dan total pemeriksaan dihitung dari LoL, bukan dari baris.
  - Uji `TestSimpanSpreadingAnekaBindCaseDEVTahapUW` (posisi ReasFacInJuniorUnderwritingB, LoL kosong, treaty 10015) MERAH
    sebelum perbaikan dengan pesan yang sama persis dengan log DEV, HIJAU sesudahnya.
  - Padanan FIRE `[terverifikasi]` aman: `hitungSaveSpreading` memanggil `totalLokasi` SEBELUM `bulatkanSimpan`, dan cabang uw
    memeriksa lewat `totalMataUang(k, adaTop)` tanpa memanggil ulang `totalLokasi`. `pertahankanSpreading` juga membulatkan
    sesudah `totalLokasi`; `lossrecord.go` / `scoring.go` hanya membaca. Di Aneka, `GantiObjekAneka` (aneka.go L419-423)
    membulatkan sesudah `totalLokasiAneka`; L688 (`tampilanSpreadingAneka`) bekerja atas salinan dan tidak ditulis.

## Catatan silang — prakondisi `.pxListSubscript>1` (11-10-2026, tiket 74 P5)

`[terverifikasi]` Langkah INSERT `SaveOfferProduction_Act` / `SaveTreatyProduction_Act` (FACINOFFER / FACINPRODUCTION) di cabang ini berprakondisi `.pxListSubscript>1` (T=2 lanjut / F=3 lewati), sama dengan
setiap cabang lain (FIRE / Aneka / Marine). Port sistem baru menulis SEMUA baris spreading. **Keputusan work owner 11-10-2026
(K74-23, `docs/issues/74-marine-cargo.md`):** tetap begitu — kode tiket ini tidak diubah.

## P6 — Upload CSV keluarga Aneka, mulai Marine Hull: riset (sesi c3, 11-10-2026)

**Pemicu:** work owner 11-10-2026 meminta tombol CSV case MARINE HULL diperbaiki ("PERBAIKI, SUDAH DI TAMBAH DI FOLDER"). Template
Pega: `D:\migrasi\RNM\DDL\templatecsvmarinehull.csv`, 31 judul. Nilai contohnya TIDAK disalin.

**Status:** riset selesai. `jenisCsv` di GET kasus sudah dipasang. Port importer menunggu K73-52..K73-62 (penampung butuh
keputusan migrasi).

### Korpus `[terverifikasi]` (`D:\migrasi\RNM\NB FacIn\`, varian NB)

- **Pemilihan:** `Activity\ChooseActivityMappingCSV_Act.xml` langkah 5 `IsAneka` → `InsertUploadAneka_act` (langkah 1 / 2 ber-label
  `//`). Pemicunya tombol "Save Data" `Section\Periode.xml` L28596 / L28873 (refresh, lalu `InputInwardFacultativeDtl`). Tidak ada
  hitung premi sesudahnya.
- **Template:** `Activity\SetDownloadMasterCSV.xml` langkah 2–6, Property-Set `InputData.CARI44` berurutan, sehingga yang terakhir
  benar menang:
  - IsAneka → TemplateCSVAneka;
  - IsMarineHull → TemplateCSVMarineHull;
  - IsContractorsPlantMachinery → **TemplateCSVAneka_CPM** (template ketiga, tidak disebut di permintaan);
  - IsFire → TemplateCSVFire;
  - IsMarineCargo → TemplateCSVMarineCargo.
  `When\IsAneka.xml` memuat IsMarineHull (B), IsContractorsPlantMachinery (W), IsBondingAndCustomBonds (AA), dan lainnya, tetapi
  TIDAK memuat isGolfInsurance. Akibatnya Golf tanpa template dan tanpa importer.
- **Tidak bercabang per lini.** `InsertUploadAneka_act` mengisi SEMUA medan untuk setiap lini Aneka: medan Marine Hull, VehicleHE,
  SerialNo, Year, Flag, dan BuiltIn. Kolom yang tidak ada di berkas tetap kosong.
- **Kolom yang dibaca** (per NAMA judul): 35. Sensus dua cara cocok: (1) urai langkah hidup / mati, (2) regex mentah atas kedua berkas;
  selisih ∅ dua arah. Uji instrumen: `RISKADDRESS_ID` tergolong mati (Insert 4.8 + Protect 4 ber-`//`), `KELURAHAN` hidup (4.6).
  - Template Marine Hull (31) ⊂ dibaca. Template \ dibaca = ∅.
  - Dibaca \ template = SERIALNO, CHASSISNUMBER, ENGINENUMBER, KELURAHAN (untuk template Aneka umum / CPM yang belum ada).
- **Insert, per baris `TempWorkPage.ListB2BHostUpload`:**
  - 1 Property-Remove LocationList + CurrencyList (MENGGANTI).
  - 4.1 `.NO == ""` → lompat (baris dilewati). `Local.IdxLocation = .NO`: LocationList DIINDEKS NILAI NO. Perbandingan
    "baris sebelumnya" memakai `ListB2BHostUpload(CountLocation)`, dan CountLocation hanya naik untuk baris ber-NO, sehingga bergeser
    bila ada baris NO kosong di tengah.
  - 4.2 `ProtectUploadAneka_Act`. TRANS `ERRMSG == ""` F=6: galat PERTAMA menghentikan impor dan `Work-.Save` tidak tercapai.
  - Lokasi baru (4.7) bila NO ≠ baris sebelumnya:
    - ObjectNo = NO;
    - RoadType / RoadName / BuildingNo / Country, ASMAddress = TYPE_ADDRESS + " " + ADDRESS, ZipCode = POSTAL_CODE;
    - kota / provinsi / district / RW = baris PERTAMA Obj-Browse `Int-RW` (ZipCode = POSTAL_CODE, Note = KELURAHAN).
    - Catatan: 4.7 memakai `(<LAST>)` untuk semua medan kecuali ObjectNo.
  - Okupasi (4.9): HANYA `OccupationList(1)`, ditimpa bila NO atau OCCUPATION_NAME berubah. ID / Name = baris pertama
    `BrowseOccupationANEKA_RD`: TYPE "ANEKA", BusinessName case, Name **Contains** tak peka huruf, urut ID DESC; sama dengan lookup
    tiket 73 P1.
  - Mata uang (4.10–4.11): `GetCurrencyIDByName` per baris → CurrencyMaster (CURR, OLDID).
  - Objek baru (4.12) bila NO / OBJECT_DESCRIPTION / SECTION berubah:
    - ObjectName = OBJECT_DESCRIPTION, Currency (Name / ID / OldID);
    - Section = SECTION I → 1, SECTION II → 2, WARRANTY → 3, selain itu 4;
    - TSI;
    - SerialNo, VehicleHE.ChassisNumber / EngineNumber;
    - Year, MarineHull.NRTRT / GRTRT / DWTRT / Construction / TypeOfVessel (= OBJECT_TYPE) / Classification;
    - Flag = FLAGNAME, BuiltIn = PLACE_OF_MANUFACTURE.
    - Medan objek diambil dari baris PERTAMA objek.
  - Coverage baru (4.13) bila objek atau COVERAGE_NAME berubah:
    - Coverage / CoverageNote = baris pertama `BrowseCoverageFacIn_RD`: Type "ANEKA", BizCode case, NamaCoverage **Contains**;
      sama dengan lookup tiket 73 P2;
    - CalculateMethod_FacIn = PRO RATA → 1, SHORT PERIOD → 2, selain itu 3;
    - Day: kosong → "365";
    - Loading = PCT_LOADING, LostLimit = PCT_SCALE, Rate = PCT_RATE, Premium = PREMIUM, DiscountPercentage / Discount,
      PctLoL / LimitofLiability, TSI = TSI baris;
    - Currency = mata uang objek.
    - Baris berurutan dengan coverage yang sama diabaikan.
  - Angka: TSI / PREMIUM / DISCOUNT / LOL buang "." lalu "," → "."; PCT_* hanya "," → ".".
  - 6 Java: buang duplikat CurrencyMaster dengan kunci `.Currency`, padahal halaman itu hanya ber-medan CURR / OLDID. `[dugaan]`
    Kuncinya selalu kosong, sehingga hanya SATU entri yang tersisa.
  - 7: CurrencyList (Premium / TSI per mata uang dari coverage) bila StatusBusiness ≠ "3".
  - 8 `Work-.Save`.
- **Protect** (`Activity\ProtectUploadAneka_Act.xml`), pesan verbatim:
  - "Occupation Cannot NULL" / "Coverage Cannot NULL". Keduanya Property-Set yang MENIMPA, jadi bila keduanya kosong yang tersisa
    pesan Coverage.
  - `\n * Occupation <NAMA> Not Available in Master` / `\n * Coverage <NAMA> Not Available in Master`.
  - Pencarian hanya bila nama berubah dari baris sebelumnya.
  - TRANS langkah 5.4 / 6.4 menguji `TradingTemp.pxResults`, halaman salinan dari cargo. `[dugaan]` Halaman itu tidak ada di alur
    Aneka, sehingga Protect keluar sesudah galat okupasi dan coverage tidak diperiksa. Hasilnya sama: impor berhenti di baris itu.
  - Langkah 4 (RiskAddress) ber-`//`.
- **Master:** POOLDATA.OCCUPATION (`DDL\OCCUPATION.txt`), COVERAGE (butir 96), RW, CURRENCY. Semuanya sudah dibaca sistem baru;
  tidak butuh DDL baru.
- **Tujuan:** semua medan sudah ada di `models.LokasiAneka` / `ObjekAneka` / `CoverageAneka` / `MarineHull` / `VehicleHE` (tiket 73
  P1/P2). Tabel objek TIDAK butuh migrasi.

### Penampung unggah — kebutuhan (rentang 255–259 HABIS; nomor TIDAK diambil)

Kolom fisik `T_NB_CSV_BARIS` = 77. Dua cara cocok: kolom `C_` di DDL 257 + 258 = daftar Go `KolomCsvObjek` ∪ `KolomCsvCargo`, selisih
∅. Dari 35 kolom Aneka, **12 sudah ada** (NO, TYPE_ADDRESS, ADDRESS, BUILDING_NO, POSTAL_CODE, KELURAHAN, OCCUPATION_NAME,
CURRENCY, TSI, OBJECT_TYPE, COVERAGE_NAME, LOL) dan **23 perlu kolom baru**:

> COUNTRY, OBJECT_DESCRIPTION, SECTION, YEAR, PLACE_OF_MANUFACTURE, NRT, GRT, DWT, FLAGNAME, CONSTRUCTION, CLASSIFICATION,
> CALCULATE_METHOD, DAY, PCT_LOADING, PCT_SCALE, PCT_RATE, PREMIUM, PCT_DISCOUNT, DISCOUNT, PCT_LOL, SERIALNO, CHASSISNUMBER,
> ENGINENUMBER

`T_NB_CSV_UNGGAH.JENIS` VARCHAR2(10) sudah cukup untuk "aneka". Tidak ada CHECK, sehingga jenis baru tanpa migrasi.

### Keputusan yang diminta (korpus tidak menjawab / Pega cacat) — JANGAN ditebak

- **K73-52 Penampung.**
  - (a) **Rekomendasi:** satu migrasi baru, `ALTER TABLE T_NB_CSV_BARIS ADD` 23 kolom `C_<judul>` VARCHAR2(4000), pola 258. Butuh
    rentang nomor baru.
  - (b) Tanpa migrasi: peta slot per jenis ke 65 kolom FIRE / cargo yang tidak dipakai Aneka (mis. `C_ITEM_TYPE_1` menampung
    OBJECT_DESCRIPTION). Bisa dikerjakan, tetapi nama kolom di DB menyesatkan saat DBA / IT memeriksa data unggah.
- **K73-53 NO.**
  - Pega mengindeks LocationList dengan NILAI NO dan memakai `<LAST>`. Akibatnya NO yang muncul lagi tidak berurutan (1, 2, 1)
    menimpa alamat lokasi 2, dan NO bukan angka 1..n tidak terdefinisi `[belum terverifikasi]`.
  - Usul: lokasi baru setiap NO berbeda dari baris sebelumnya, ObjectNo = NO apa adanya. NO yang muncul lagi tidak berurutan →
    galat jelas (bukan menimpa).
- **K73-54 Baris NO kosong.** Pega melewatinya diam-diam, tetapi perbandingan "baris sebelumnya" ikut bergeser. Usul: lewati baris
  itu dan bandingkan dengan baris terproses sebelumnya.
- **K73-55 Mata uang tak ada di master.** Pega diam (Currency.ID kosong). Usul: galat, seperti K74-21 / K66 D6.
- **K73-56 CurrencyList.** Pega membangun ulang lewat langkah 7, tetapi duplikatnya dibuang Java dengan kunci yang salah. Usul:
  CurrencyList + cicilan dihapus tanpa dibangun ulang, seperti K66 D5 / K74-21.
- **K73-57 Premi.**
  - (a) Ikut Pega: PREMIUM CSV disimpan apa adanya, tanpa hitung.
  - (b) Seperti FIRE K66 D5: dihitung mesin `hitung-coverage-aneka` dari Rate / TSI / Day / Loading / Discount.
- **K73-58 RW.**
  - Pega: Obj-Browse ZipCode + Note = KELURAHAN. Template Marine Hull TIDAK punya KELURAHAN, dan Obj-Browse dengan nilai kosong
    `[belum terverifikasi]`.
  - Usul: pakai `cariRWCsv` / `pilihRWCsv` CSV FIRE (K66-B8 / B9): kode pos dulu; kelurahan bila ada.
- **K73-59 Okupasi / coverage "Contains" + baris pertama.** "ICC" bisa terpilih ke coverage lain yang memuat teks itu.
  - (a) Ikut Pega.
  - (b) Cocok persis, tak peka huruf. Tidak ada yang persis → galat "Not Available in Master".
- **K73-60 Template Aneka umum / CPM belum ada.**
  - Usul: unggah `marinehull` dengan judul wajib = 31 judul template, kolom lain diabaikan (pola K74-18).
  - `aneka` / `cpm` tetap 409 sampai template dikirim. Saat itu judul wajibnya = judul template.
- **K73-61 Bonding / CustomBonds** termasuk IsAneka. Pega mengimpornya ke AnekaList generik, tanpa MemorandumAccept. Ikut Pega, atau
  tolak CSV untuk Bonding?
- **K73-62 Galat pertama menghentikan impor** (Pega 4.2 F=6), sama dengan cargo P6. Konfirmasi saja. Alternatifnya: kumpulkan
  semua galat seperti CSV FIRE D6 / D7.

### Jawaban work owner (11-10-2026)

- **K73-52 → (a)** migrasi. Rentang ketiga nbfacin = **250-254**, diambil dari rnwfacin (yang belum memakai satu nomor pun);
  rnwfacin menjadi 220-249. Polanya sama dengan 255-259. `MODUL.md` kedua modul diperbarui.
  - **Ralat bentuk (work owner 11-10-2026):** ALTER `T_NB_CSV_BARIS` bernomor 250 akan jalan sebelum 257 yang membuat tabelnya
    (`inti/backend/migrasi/migrasi.go` mengurutkan menurut nama berkas), sehingga gagal di DB baru. Bentuknya diganti:
    250 = `CREATE TABLE T_NB_CSV_BARIS_ANEKA` (UNGGAH_ID, NO_BARIS, PK keduanya, kolom `C_<judul>` VARCHAR2(4000)), tanpa FK,
    pembersihan lewat kode.
- **K73-57 → (a)** ikut Pega: PREMIUM / DISCOUNT dari CSV disimpan apa adanya, tanpa hitung ulang.
- **K73-59 → (a)** ikut Pega: Contains + baris pertama.
- **K73-53 / 54 / 55 / 56 / 58 / 60 / 61 / 62 →** usulan disetujui: NO berulang tidak berurutan = galat; baris NO kosong
  dilewati; mata uang tak ada = galat; CurrencyList dihapus tanpa dibangun ulang; RW lewat `cariRWCsv` / `pilihRWCsv`; judul wajib
  marinehull = 31, aneka / cpm 409 sampai template dikirim; Bonding / CustomBonds diimpor seperti Pega; galat pertama menghentikan
  impor.

- **K73-63 (work owner 11-10-2026, sesudah bug DEV premi > TSI):** impor CSV Aneka menolak baris dengan PREMIUM > TSI pada
  baris yang sama. Ini pengaman tambahan sistem baru, penyimpangan sadar: Pega `InsertUploadAneka_act` L4363 menyimpan PREMIUM apa
  adanya, dan tidak ada proteksi premi > TSI di korpus. Baris contoh `templatecsvmarinehull.csv` sendiri ber-PREMIUM jauh di
  atas TSI. K73-57 tidak berubah.

**Frontend P6 (11-10-2026, selesai):** `components/csvMarineHull.ts` (31 judul VERBATIM `templatecsvmarinehull.csv`, semua wajib,
urutan bebas; template unduhan = baris judul saja); `BilahCsvObjek` jenis `'marinehull'` lewat tabel `CSV_PER_JENIS`;
`InwardFacultative` memasang bilah bila `jenisCsv` server = `'marinehull'` (FIRE / cargo tetap dari penanda lama). Aneka umum / CPM
tetap pesan "belum tersedia". Sampai importer c3 terpasang, unggah Marine Hull dijawab 409 oleh server.

Mengikuti preseden tanpa pertanyaan:

- spreading Aneka case terbuang karena objek diganti (K66 D5);
- unggah tanpa batas lewat infrastruktur tiket 75 (K74-17);
- kelompok ≥ 1000 baris tidak membelah lokasi;
- kunci K75-2 berlaku untuk PUT `…/objek-aneka` dan Save spreading Aneka.

### Dibangun sekarang (aditif)

- **GET kasus `jenisCsv`** (`services/lini.go` `templateCsvKasus`, `handlers/kasus.go`, `api.ts` `KasusNB.jenisCsv?`):
  - nilai 'fire' | 'cargo' | 'aneka' | 'marinehull' | 'cpm' | '', menurut SetDownloadMasterCSV langkah 2–6 berurutan;
  - 'cpm' DITAMBAHKAN dari korpus (TemplateCSVAneka_CPM).
  - Importer ketiga jenis Aneka sama.
  - ~~Selama port belum ada, unggah untuk case Aneka tetap 409 (`ErrCsvBukanFire`).~~ Importer terpasang (lihat bagian berikut);
    Aneka umum / CPM → 409 `ErrCsvTemplateAnekaBelumAda`.
- **Uji:**
  - `TestJenisCsvKasus`: Marine Hull 05 / 55, CPM 12, Liability, HE, Golf "", FIRE, cargo, tanpa pembaca;
  - `TestBacaKasus` 21 kunci;
  - `TestKunciKasusSamaDenganApiTs` 21 kunci, `jenisCsv` terakhir.

### Backend P6 Aneka — importer dibangun (sesi c3, 11-10-2026)

**Status:** dibangun, uji hijau (`go vet`, `go test -short ./...`). Migrasi **250** DITULIS, belum dijalankan. Belum dicoba di
DEV.

**Migrasi 250 — WAJIB sebelum biner bila unggah CSV case Aneka dipakai:**

- `CREATE TABLE T_NB_CSV_BARIS_ANEKA` (UNGGAH_ID, NO_BARIS, 35 kolom `C_<judul>`), tanpa FK (K73-52 opsi 1).
- Tanpa 250, hanya unggah CSV case Aneka yang gagal (ORA-00942). FIRE / cargo tidak menyentuh tabel ini.
- Tidak bergantung urutan dengan 257 / 258.
- Penyimpangan BENTUK dari teks awal K73-52 ("ALTER `T_NB_CSV_BARIS`, pola 258") sudah diputuskan work owner (lihat ralat di atas).
- `repository/csvunggah.go` memilih penampung per jenis (`tabelUnggah.penampung`).
- Tidak ada kode yang menghapus baris penampung, untuk jenis mana pun — pembersihan baris FIRE / cargo juga belum ada.

**Pemilihan importer** (`services/csvfire.go` `jenisCsvKasus`):

- IsFire → fire;
- IsMarineCargo → cargo;
- IsAneka dengan `jenisCsv` = marinehull → **aneka**;
- IsAneka lainnya → 409 `ErrCsvTemplateAnekaBelumAda` (K73-60);
- selain itu → `ErrCsvBukanFire`.
- Save Data JSON lama tetap FIRE saja.

**Port** (`services/csvaneka.go`, bukti per langkah di kepala berkas):

- Judul per nama: 31 judul template wajib, 4 opsional.
- Angka: ribuan / persen seperti Pega. TSI "1.5" → 15 (titik = ribuan) ditiru apa adanya. Angka wajib terurai (preseden K66 D7).
- Protect berurutan:
  - kosong menimpa;
  - okupasi tidak ada → keluar;
  - Contains + baris pertama (K73-59);
  - BusinessName = Class of Business case;
  - BizCode = BUSINESS.ID;
  - mata uang tidak ada / kosong = galat (K73-55).
- Validasi berhenti di baris PERTAMA bermasalah (K73-62). NO berulang tidak berurutan = galat (K73-53). Baris NO kosong dilewati,
  isinya tidak diperiksa (K73-54).
- Susun pohon (Insert 4.7 / 4.9 / 4.12 / 4.13):
  - lokasi per NO;
  - okupasi tunggal yang DITIMPA;
  - objek per (deskripsi, section) dengan medan dari baris pertama;
  - coverage per nama; baris coverage sama yang berurutan diabaikan;
  - Premium / Discount apa adanya; medan server RNM tidak dihitung (K73-57).
- RW (K73-58):
  - KELURAHAN terisi → `RWKelurahan`;
  - kosong → `RWKodePos` baru: baris pertama urut NOTE / DISTRICTNAME / CITYNAME `[keputusan agent]`;
  - kode pos tidak ada → kota / provinsi kosong.
- Simpan berkelompok:
  - ≥ 1000 baris dan berakhir di batas lokasi;
  - kelompok pertama mengganti seluruh objek case (coverage Aneka lama + anak + okupasi + lokasi) dan menghapus CurrencyList +
    cicilan (K73-56);
  - SEQ_NO lokasi menyambung lokasi terbesar (`repository/anekacsv.go`, `sisipAnekaBerkelompokDari`; `GantiObjekAneka` tetap awal 0).
- **Kunci K75-2:** PUT `…/objek-aneka` dan Save spreading Aneka 409 selama unggah berjalan (gerbang baru, hanya menambah).

**Uji** (`services/csvaneka_test.go`, `repository/anekacsv_test.go`, `repository/csvunggah_test.go`):

- peta judul, urai angka / pemetaan, Protect (8 kasus pesan), susun pohon;
- alur 2500 baris;
- kelompok tidak memotong lokasi tiga-baris;
- baris NO kosong berisi rusak dilewati;
- berhenti di galat pertama;
- NO berulang;
- kunci K75-2;
- offset SEQ_NO;
- kolom DDL 250 = kode, dua arah.

**Mutasi:** 18 diuji, semuanya tertangkap (satu lewat timeout `go test`: `akhir` tidak maju untuk baris NO kosong terakhir →
simpan berputar). Tiga awalnya lolos dan menghasilkan uji tambahan: batas kelompok, baris NO kosong, offset SEQ_NO.

**Nama berkas template unduhan:** frontend memakai `templatecsvmarinehull.csv` verbatim (judul saja). Backend tidak menyajikan
template.

### Bug DEV 11-10-2026 — unggah CSV 3 case NB-151034 gagal simpan ORA-12899 `T_COVERAGELIST.COVERAGE_NOTE` (sesi c3)

- **Sebab `[terverifikasi]`:**
  - `services/csvaneka.go` mengisi CoverageNote = COVERAGE.NAME, yang di DDL berlebar VARCHAR2(4000 BYTE) (`DDL\COVERAGE.txt` L4).
  - Kolom tujuannya VARCHAR2(500), lebar RANCANGAN (193).
- **Perbaikan — migrasi 251** `ALTER TABLE T_COVERAGELIST MODIFY (COVERAGE_NOTE VARCHAR2(4000))`, satu pernyataan per blok;
  `_down` kembali ke 500 dan gagal ORA-01441 bila sudah ada isi > 500.
  - `services/loader/amandemen.go` mencatat `amandemenLebar` (berkas bangkitan `skema_gen.go` tidak disentuh).
  - Validasi coverageNote FIRE / Aneka / cargo dinaikkan ke 4000.
  - STRUKTUR diperbarui.
  - Penjaga migrasi-sebagian membaca 251.
- **T_ADDITIONALCOVERAGE.COVERAGE_NOTE (Warranty Extension cargo) TIDAK ikut.** Tabelnya lahir di 256; pelari migrasi mengurutkan
  menurut nama berkas, sehingga ALTER di 251 gagal ORA-00942 di basis data baru. Nomor nbfacin > 256 sudah habis.
  - Kolom ini tetap 500 dan divalidasi 400 di PUT objek-cargo, jadi tidak ada crash ORA.
  - Lebar sumbernya (CONDITION.CONDITIONNAME, view JSON) tidak tertulis di DDL.
- **Ulangi:** pekerjaan gagal di fase simpan melanjutkan simpan dari BARIS_TERSIMPAN. Kemajuan ikut rollback bersama kelompok yang
  gagal, jadi kegagalan di kelompok 1 berarti ganti lagi dari baris 1.
  - Uji: `TestUlangiSimpanAnekaSesudahGagal`, gagal di kelompok 1 dan 2, tanpa lokasi ganda.
  - Sesudah 251 dijalankan DBA: klik Ulangi.

### Audit lebar sesudah ORA-12899 — keputusan work owner 11-10-2026 (lewat sesi 9d; sesi c3)

- **(a) RW, tanpa migrasi.** CSV Aneka menolak dengan galat jelas per baris bila RW.NOTE / DISTRICTNAME / CITYNAME / PROVINCENAME
  (master VARCHAR2(4000), `DDL\RW.txt`) melebihi kolom tujuan ASMRW / ASM_DISTRICT / ASM_CITY / PROVINCE VARCHAR2(100) (192).
  - Lebar diukur dalam BYTE (len UTF-8), bukan jumlah rune. Teks pesan buatan sistem baru, contoh: `baris 3: master RW kode pos
    "99999" - CITYNAME 102 byte melebihi kolom T_RISKLOCATION.ASM_CITY (paling banyak 100 byte)`.
  - Diperiksa saat validasi untuk setiap lokasi baru. RW kini dicari juga di fase validasi, sehingga galat pertama menghentikan
    impor (K73-62).
  - Diperiksa ulang saat simpan: RW yang memanjang sesudah validasi → pekerjaan gagal "data master berubah sejak validasi", objek
    tidak ditulis.
  - `services/csvaneka.go` `periksaRWAnekaCsv`; uji `TestLebarRWAnekaCsv` (tepat 100 byte sah; empat medan; validasi; simpan).
- **(c) COVERAGE.ID / OLDID:** dibiarkan, tanpa bukti data.
- **Mutasi:** batas `>=`, ukur rune, tanpa cek di validasi, tanpa cek di simpan — semuanya tertangkap.

### K73-63 — PREMIUM > TSI ditolak di CSV Aneka (keputusan work owner 11-10-2026, lewat sesi 9d; sesi c3)

- **Pemicu:** case Marine Hull di DEV ber-Total Premium Spreaded 496,8 M dengan TSI 1,6 M. InsertUploadAneka_act L4363 mengisi
  Premium = PREMIUM CSV apa adanya (K73-57), dan baris contoh `templatecsvmarinehull.csv` memang ber-PREMIUM jauh di atas TSI.
  Pega TIDAK memproteksi premi > TSI.
- **Keputusan: pengaman sistem baru, PENYIMPANGAN SADAR dari Pega.**
  - Baris dengan PREMIUM (sesudah urai) > TSI pada baris yang sama → galat `baris N: PREMIUM <nilai> melebihi TSI <nilai>`
    (nilai bentuk desimal titik).
  - Premi = TSI sah. Premi yang lolos tetap dari CSV (K73-57 tidak berubah).
  - Galat pertama menghentikan impor (K73-62).
  - Tidak menyentuh FIRE / cargo / Save layar.
- **Port:** `services/csvaneka.go` `uraiBarisAnekaCsv`, dipakai validasi DAN simpan. Di simpan, baris penampung tidak berubah
  sejak validasi, jadi cek ulang itu praktis tidak pernah memicu.
  - `[keputusan agent]` dicek hanya bila PREMIUM dan TSI keduanya terisi.
- **Uji:** `TestPremiMelebihiTSIAneka`: sama sah; lebih 0,01 galat; "43.200.000.002,00" vs "100.000.001,00" galat; kurang sah;
  salah satu kosong tanpa cek; validasi berhenti di baris 3.
  - Fixtur `TestSusunAnekaCsv` baris "UJI MESIN" diberi PREMIUM 5 (≤ TSI 7).
  - Mutasi (batas `>=`, tanpa cek) tertangkap.

## Catatan silang — IDPEGA = case id polos (11-10-2026)

Tulisan IDPEGA tiket ini kini memakai case id polos; pembaca menerima dua bentuk (polos / berawalan). Rincian: `docs/issues/69-jalur-policy.md`
bagian "IDPEGA = case id polos".
