# Paritas layar dan aksi — Endorsement Life

> Disusun 01-10-2026 (gelombang 1, brief `PROMPT-IMPLEMENTASI-TIGA-MODUL-LIFE-GELOMBANG-2.md` §5). Satu baris per harness, section, tombol,
> dan aksi korpus → rule yang dipanggil → RDB/RD → tabel/kolom → rute API → komponen React.
> Nomor baris `bNNN` = baris yang **memuat nilai yang dikutip** pada berkas korpus hasil pecah standar, dapat diulang dengan
> `sed -e 's/></>\n</g' "<Tipe>/<Rule>.xml" | grep -n '<teks yang dikutip>'` atas `D:\XML\RNM_BRD\Endorsement Life\`
> *(berkas korpus sudah satu tag per baris, jadi nomor ini sama dengan nomor baris mentahnya)*. Langkah activity dirujuk lewat baris
> `pyStepsActivityName`-nya; prakondisi lewat `pyStepsPreCondParamsWhen`; saringan RD lewat `pyFilterName`; grid lewat `pyBodyType` `REPEATING`;
> wadah lewat `pyContainerVisibleWhen`; tombol lewat `pyLabel` di `pyModes`; aksi lewat `pyAction`. Baris di dalam teks SQL multi-baris
> (`pyBrowseSQL`) dikutip langsung.
> Setiap langkah activity mencetak `pyStepsBlockName` — `·` = kosong (langkah hidup), `//` = ter-remark (tidak pernah jalan).
> Precondition: `WhenTrue 2` = lanjut, `3` = lewati langkah, `6` = keluar; `PRE=false` = prakondisi dimatikan → langkah **selalu** jalan.
> Keadaan: ✅ dibangun · ⏸️ menunggu OQ · ➖ sengaja tidak dibawa (bukti di kolom Catatan) · ➕ tambahan tiket tanpa padanan XML.
> Rute API ber-prefix `/api/endorsement-life`. Komponen di `modul/endorsementlife/frontend/`. Temuan `Rnn` dan ralat `En`:
> `RALAT-DEV-01-10-2026.md`. OQ: `OQ-ENDORSEMENT-LIFE.md`.

## 0. Alat baca dan sensus

| Hal | Isi |
| --- | --- |
| Alat | pengurai expat bernomor baris (`pohon_activity.py`, `pohon_layar.py`, `baca_rd.py`, `sql_rdb.py` — scratchpad sesi, tidak di repo) |
| Berkas | 75 XML: Activity 16, RDBList 17, Section 14, Harness 8, FlowAction 6, ReportDefinition 6, DataTransform 2, When 2, Flow 1, DecisionTable 1, ConnectREST 1, SystemSettings 1 |
| Langkah ber-`//` | **20**: `CreateCaseEMDL` 7 (langkah 7–13), `UploadCSVEDMLifePremium_Act` 5 (4, 5, 5.2, 5.4, 6), `GetOldDetail_EDM` 4 (1, 3, 4, 5), `InsertJsonPolisLife_Act` 2 (13, 15), `SetErrorBatalEndorsement_Act` 1 (4), `serviceInsertArasapasLife_act` 1 (7). Sepuluh activity lain nol |
| Tombol — cara 1 | cacah baris `<pyFormat>pxButton</pyFormat>` di 14 Section: **22** (`ConfirmSubmitEDM` 1, `EndorsmentLife_Section` 2, `InboxEndorsementLife` 4, `InputEDMLife` 9, `ShowLifePremiumSummary_EDM` 5, `ViewCSVResult_LifeEDM` 1) |
| Tombol — cara 2 | pengurai sel `pyFormat = pxButton`: **22**, sepakat per berkas. Daftar lengkap bab 10 |
| Harness vs Section | salinan section di dalam kedelapan harness identik dengan berkas Section-nya; tabel ini berpijak pada berkas Section |
| Visibilitas | sel: `pyVisible` (`ALWAYS`/`OTHER`/`NOTBLANK`) + kondisinya; wadah: `pyContainerVisibleWhen` berlaku bila opsi `CONDITION` **atau** `ExpressionCondition` |
| Rahasia | `InsertJsonPolisLife_Act` langkah 15 (`//`) memuat kredensial dan alamat server — **tidak disalin ke berkas mana pun**. `SystemSettings/LinkService.xml` dan alamat `M_LINK_SERVICE` tidak dikutip |

## 1. Halaman awal — pintu masuk modul

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Harness/InboxEndorsementLife.xml` b27 `DATA-PORTAL!INBOXENDORSEMENTLIFE`, b91 kelas `Data-Portal` (satu-satunya harness portal di modul ini), b142 label **`Endorsement Life`**, b843 `pyInclude` `InboxEndorsementLife` | tombol menu **Endorsement Life** di kelompok `TREATY` (`M_NAV_MENU`, slot 976) → `HALAMAN_AWAL` = halaman `edm-inbox`, `pages/InboxEndorsementLife.tsx` | ✅ |
| Tujuh harness lain ber-kelas work/int dan dibuka aksi: `EndorsmentLife_harnes` (b91 `ASM-FW-GISFW-Work-EndorsementLife`) oleh `Create Addendum` (`Section/InboxEndorsementLife.xml` b6984 `showHarness`); `ViewCSVResult_LifeEDM` oleh `View Upload` (`InputEDMLife.xml` b9844); `ViewOldPolicy_EDM`/`_QP`/`_TR`/`_TP` oleh `View Old Policy` (`ShowLifePremiumSummary_EDM.xml` b65274/b65806/b66364/b66922); `EditDetail_Harness` oleh `View Premium` (`InputEDMLife.xml` b36673, tombol mati) | bukan pintu masuk | — |
| Ikon kerangka harness `pxIconHistory`/`pxIconAttachments`/`pxIconExpandCollapse`/`pxIconCancel` (`Harness/InboxEndorsementLife.xml` `pyFormat` b508, b568, b626, b686) | tutup panel/modal | ➖ kontrol bawaan Pega tanpa aksi modul |

## 2. Kotak masuk — `Section/InboxEndorsementLife.xml`

Tabel `T_PREMIUM_LIST` (baris kasus EDM, R20). Rute `GET /inbox`. Komponen `pages/InboxEndorsementLife.tsx`.

| Korpus | Rule / data | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| wadah b4992 `InData.CARI1!=1` — judul b5505 **`Endorsement Life`** | — | judul halaman | ✅ |
| teks b5626 `VIS 1=2` | — | — | ➖ mati |
| tombol **`Create Case Endorsement`** b5879 `VIS=never` → `CreateCaseEMDL` b6067 | — | — | ➖ mati (R11) |
| tombol **`Create Addendum`** b6620 → `showHarness` b6984 `EndorsmentLife_harnes` (`pyTarget=current`) | — | buka halaman `edm-buat` | ✅ |
| grid b8284 `REPEATING`, RD `ReportDefinition/InboxEDMLife.xml` (`pyRDName` b12852): saring b659 `pyStatusWork != "Resolved-Completed"`, b673 `!= "Resolved-Rejected"`; urut b1040 `pxCreateDateTime` `DESC` | `STATUSS IS NULL`, `ORDER BY TGL_INPUT DESC, ID` (R27) | tabel berhalaman | ✅ |
| kepala grid: `Case ID` b8418 · `Endorsement No` b8556 · `Type` b8692 · `EDM Type` b8828 · `Policy No` b8964 · `SOB` b9100 · `Ceding` b9236 · `Policy Holder` b9372 · `Marketing Name` b9508 · `Create Date` b9644 · `Create Operator` b9780 · `EDM Type Batal` b9916 · `Status` b10052 | sel: `A.pyID` b10202 → `ID` · `A.PremiumListSummary.PL_NUMBER` b10516 → `PL_NUMBER` versi (R25) · `A.Type` b10700 → `TYPE` · `A.EdmType` b10869 → `EDM_TYPE` · `A.PolicyNo` b11053 → `OLD_POLICY_NO` · `A.SobName` b11222 → `SOB_NAME` · `A.CedingCoName` b11391 → `CEDING_CO_NAME` · `A.PolicyHolderName` b11560 → `POLICY_HOLDER_NAME` · `A.MarketingName` b11729 → `MARKETING_NAME` · `A.pxCreateDateTime` b11897 → `TGL_INPUT` · `A.pxCreateOpName` b12080 → `CREATE_OP_NAME` · `A.EdmTypeBatal` b12243 (`VIS A.EdmType==3`) · `A.pyStatusWork` b12457 → `STATUSS` | 13 kolom, label VERBATIM | ✅ (kolom `EDM Type Batal` ⏸️ OQ-EDM-005) |
| tautan `A.pyID` b10202 → `openAssignment` b10349 | — | buka halaman `edm-kasus` | ✅ |
| wadah b729 `InData.CARI1==1`: `Policy No` b1502, `EDM Note` b1726, `Effective Date` b2477, `Edm Type` b2715, tombol **`Create`** b3780 (`CreateCaseEMDL` b3966, `openAssignment` b4025), tombol **`Cancel`** b4228 (`CancelCreateCaseEDML` b4318) | `InData.CARI1` tak pernah bernilai 1 (R11) | — | ➖ tak terjangkau |

## 3. Buat endorsement — `Harness/EndorsmentLife_harnes.xml` → `Section/EndorsmentLife_Section.xml`

Rute `POST /kelayakan`, `POST /kasus`. Komponen `pages/BuatEndorsement.tsx`.

| Korpus | Rule / data | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| `Policy No` b1075 `TempWork.PolicyNo` → ubah b1249 `SetErrorBatalEndorsement_Act` (`Nopolis<-TempWork.PolicyNo`) | gerbang, bab 3a | isian + cek kelayakan saat berubah | ✅ |
| `EDM Type` b1347 `TempWork.EdmType` → ubah b1545 `SetErrorBatalEndorsement_Act` | opsi `1` Perubahan Data, `3` Batal (spec §5 `[keputusan work owner]`) | pilihan | ✅ |
| `EDM Type Perubahan Data` b1645 `VIS TempWork.EdmType=1` → ubah b1841 refresh | `pyListSource associated` | — | ⏸️ OQ-EDM-005 (R12) |
| `EDM Type Batal` b1937 `VIS TempWork.EdmType=3` | `pyListSource associated` | — | ⏸️ OQ-EDM-005 (R12) |
| `Description` b2140 `TempWork.DesBatal` | `MappingEDMLife` 7 b1444/b1445 | `EDM_NOTE` (R13) | ✅ |
| `EDM Date` b3160 `TempWork.EdmDate` → ubah b3347 `SetErrorBatalEndorsement_Act` | `MappingEDMLife` 7 b1465/b1466 | `EDM_DATE` | ✅ |
| tombol **`Process Policy No`** b2414 `VIS=never` → `MappingEDMLife` b2533 | — | — | ➖ mati |
| tombol **`Submit`** b4226 `VIS (Protect.CARI1 = '0' && (TempWork.EdmType='1' \|\| TempWork.EdmType='3'))` → `MappingEDMLife` b4415, `openAssignment` b4450 | bab 3b | tampil bila nol pesan kelayakan dan `EdmType` `1`/`3`; `POST /kasus` lalu buka `edm-kasus` | ✅ |

### 3a. Gerbang kelayakan — `Activity/SetErrorBatalEndorsement_Act.xml` (11 langkah + 1 `//`)

| Langkah | Isi | Sistem baru |
| --- | --- | --- |
| 1 b285 `·` | pesan VERBATIM: b312 `"Sudah Di endorsement Batal"`, b359 `"There's EDM with this policy no that haven't finish yet!"`, b380 `"Invalid Policy No !"`, b401 `"Policy no Not Found !"`, b422 `"There's already payment with this policy no"`; b443 `Protect.CARI1 = "0"`; b464 nomor invoice `@replaceAll(TempWork.PolicyNo,".","")` | konstanta pesan di `models`; ✅ |
| 2 b545 `·` | `Page-Clear-Messages` | pesan dihitung ulang setiap panggilan; ✅ |
| 3.1 b691 `·` prakondisi b801 `param.Nopolis==""` | pesan `Invalid Policy No !` | ✅ |
| 3.2 b841 `·` `GetPL_NumberLife` (b895); 3.3 b1015 `·` prakondisi b1125 `CARI1==""` | `RDBList/GetPL_NumberLife.xml` b84 `JSON_POLIS` `NOPOLIS`; pesan `Policy no Not Found !` | polis ada di salah satu sumber (RALAT bab 4); ✅ |
| 3.4 b1185 / 3.5 b1378 `·` RD `FilterProteksiEDMLife` (b1212); 3.6 b1561 `·` prakondisi b1671 `@LengthOfPageList(ListEdm.pxResults)>0` | RD b526 `pyStatusWork != Resolved-Completed`, b541 `.PolicyNo = Param.Nopolis` | kasus terbuka atas polis itu (`OLD_POLICY_NO`, `STATUSS IS NULL`, `IDX_PL_OLD_POLICY_NO`) → pesan `There's EDM …`; ✅ (R04, OQ-EDM-002) |
| 3.7 b1711 `·` `GetEdmTypeLife` (b1766); 3.8 b1885 `·` prakondisi b1995 `@contains(OutData1.pxResults(1).CARI1,"3")` | `RDBList/GetEdmTypeLife.xml` b84 `A.DATA_JSON.EdmType … PRODKE IS NOT NULL ORDER BY PRODKE DESC` | `EDM_TYPE` versi berjalan memuat `3` → `Sudah Di endorsement Batal`; ✅ |
| 3.9 b2035 `·` `SearcStatusBayarArasaps_SQL` (b2090); 3.10 b2209 `·` prakondisi b2319 `@LengthOfPageList(ListPembayaran.pxResults)>0&&TempWork.EdmType=="3"` | `RDBList/SearcStatusBayarArasaps_SQL.xml` b58 `ARASAPAS.DETAIL_INVOICE` … `IVD_JR_ID ='5'` | repository `arasapas` tunggal, dibaca hanya bila `EdmType=3` → `There's already payment …`; ✅ (R26, OQ-EDM-012) |
| 3.11 b2359 `·` prakondisi b2461 `@hasMessages(myStepPage)` → b2385 `Protect.CARI1 = 1` | — | `bolehBuat = false` bila ada pesan; ✅ |
| 4 b2550 `//` `Obj-Save` | — | ➖ ter-remark |

### 3b. Buat kasus dan salin polis lama — `Activity/MappingEDMLife.xml` (14 langkah, nol `//`) + `Activity/CreateCaseEMDL.xml`

Tabel tulis: `T_PREMIUM_LIST` (kasus = versi baru, `ID` `EDMLF-<n>`), `T_PREMIUM_LIST_DETAIL`, `T_PREMIUM_LIST_SPREADING`, `T_PREMIUM_LIST_SPREADING_RETRO`; jejak `T_CLAIMLF_JEJAK` (`inti/backend/jejak`).

| Langkah | Isi | Sistem baru |
| --- | --- | --- |
| 1 `·` b388 *"Jika ada edm blm resolve exit act"*, prakondisi b433 `@contains(OutData1.pxResults(1).CARI1,"3")` T=6, b456 `@LengthOfPageList(ListEdm.pxResults)>0` T=6 | keluar bila polis sudah Batal atau masih punya EDM terbuka | `POST /kasus` menjalankan **kelima** gerbang ulang di server, di dalam transaksi, dengan kunci baris versi berjalan (`FOR UPDATE`); ✅ |
| 2 b496, 3 b689 `·` `svcAddWorkObject` (kelas b519 `ASM-FW-GISFW-Work-EndorsementLife`) | buat work object | baris `T_PREMIUM_LIST` baru, `ID = ID_PEGA = EDMLF-<n>` dari `SEQ_WORK_EDM_LIFE` (R20); ✅ |
| 4 b837, 5 b971 `·` | catatan assignment `!InputEDMLife` (b864), buka kasus ber-kunci | — (kasus langsung di tahap `InputEDMLife`); ✅ |
| 6 b1114 `·` `GetProdkeNopolis` (b1172) | `RDBList/GetProdkeNopolis.xml` b84 `IDPEGA … ORDER BY PRODKE DESC` | versi berjalan polis lama (RALAT bab 4); ✅ |
| 7 b1291 `·` | `OldCaseID` b1313, `PolicyNo` b1361, `EdmType` b1382, `EdmTypePerubahanData` b1403, `EdmTypeBatal` b1424, `DesBatal` b1445, `EdmDate` b1466 | `OLD_POLICY_NO`, `EDM_TYPE`, `EDM_NOTE`, `EDM_DATE`; dua sub-jenis ⏸️ OQ-EDM-005; ✅ |
| 8 b1547 `·` `Obj-Open-By-Handle` work NB (b1599) | — | baca versi berjalan dari sumbernya; ✅ |
| 9 b1690 `·` PRE=false (selalu) | 25 properti kepala: `Type` b1717, `TypeCeding` b1764, `BusinessCode` b1785, `SobName` b1806, `PolicyHolderName` b1827, `CedingCoName` b1848, `MarketingName` b1869, `ProRateType` b1890, `BusinessName` b1911, `PolicyHolder` b1932, `SourceOfBusiness` b1953, `CedingCo` b1974, `MOID` b1995, `MarketingCode` b2016, `DateReceived` b2037, `NoOffer` b2058, `RetroID` b2079, `RetroName` b2100, `SecurityReinsurerID` b2121, `SecurityReinsurer` b2142, `RISLIPRNM` b2163, `TempWork.PolicyNo` b2184 → `PL_NUMBER` (b2183), `ProductName` b2205, `ProductNameID` b2226, `WPC` b2247 | kolom `T_PREMIUM_LIST` bernama sama (R09); ✅ |
| 10 b2339 `·` PRE=false | salin `PremiumListDetail` | salin peserta + spreading + spreading retro, `PARENT_ID` = peserta sumber (sumber warisan: kosong, OQ-EDM-007); ✅ |
| 11 b2490 PRE=false; 11.1 b2583 `·` prakondisi b2691 `.EDMStatus=="Delete"` T=3 → b2609 `"Old"`; 11.2 b2737 `·` prakondisi b2831 → hapus baris `Delete` | — | baris `Delete` versi lama tidak disalin, sisanya `EDM_STATUS = 'Old'` (R10); ✅ |
| 12 b2930 `·` b2957 `.EditInput = 1` | kunci field | kasus selalu terkunci sesudah dibuat (tiket 08); ✅ |
| 13 b3064 `·` `Obj-Save`, 14 b3196 `·` `Commit` | commit dini | satu transaksi `POST /kasus`, commit sekali — kasus terlihat pengguna lain begitu dibuat (spec §11); ✅ |
| `CreateCaseEMDL` 1–6 `·` (b244–b929) | pembuat kasus jalur panel Inbox | ➖ jalurnya tak terjangkau (R11) |
| `CreateCaseEMDL` 7–13 `//` (b1072 `GetClobJson_polis` … b2450 `Commit`) | penyalinan data lama versi lama | ➖ ter-remark |
| `CancelCreateCaseEDML` 1 b222 `·` b245 `InData.CARI1 = "0"` | tutup panel | ➖ panelnya tak terjangkau (R11) |

## 4. Rincian endorsement — FlowAction `InputEDMLife` (b166) → `Section/InputEDMLife.xml`

Rute `GET /kasus/{id}`, `POST /kasus/{id}/simpan`, `POST /kasus/{id}/putuskan`. Komponen `pages/InputEDMLife.tsx`.

### 4a. Kepala (seluruhnya baca-saja)

| Korpus | Kolom | Keadaan |
| --- | --- | --- |
| judul b817 **`Endorsement Life Detail`** | — | ✅ |
| `.PolicyNo` b3155 (label dari deskripsi properti — `pyLabelFieldValue` b3126 `Text Input`; `Section/EndorsmentLife_Section.xml` b1075 menamainya `Policy No`), `pyDisabledWhen` `.EditInput==1` b3166; ubah b3301 `SetErrorBatalEndorsement_Act` | `OLD_POLICY_NO` | ✅ terkunci + pesan penjelas (tiket 08) |
| `Product Name` b3399 · `Product Name ID` b3589 | `PRODUCT_NAME`, `PRODUCT_NAME_ID` | ✅ |
| `Type` b3812 · `Reinsurance System` b4143 · `Class of Business` b4473 · `SOB` b4658 (RD `BrowseCedingCoLife_RD` b4726) · `Policy Holder` b4897 (RD `BrowseClientNusaRe_RD` b4963) · `Premium Method` b5976 · `Ceding` b6339 (RD b6407) · `Marketing Officer` b6578 (RD `BrowseMarketingOfficer_RD` b6683) | `TYPE`, `TYPE_CEDING_NAME`, `BUSINESS_NAME`, `SOB_NAME`, `POLICY_HOLDER_NAME`, `PRO_RATE_TYPE`, `CEDING_CO_NAME`, `MARKETING_NAME` | ✅ baca-saja (`ro`); pemilih RD dan aksi ubahnya (`runDataTransform` b4016/b4347/b6197) tidak pernah jalan pada sel baca-saja ➖ |
| `EDM Type` b6910 `VIS .PolicyNo!=''` · `EDM Type Perubahan Data` b7203 `VIS .EdmType=1` · `EDM Type Batal` b7494 `VIS .EdmType=3` (`pyDisabledWhen` b7537) · `Description` b7697 | `EDM_TYPE`, `EDM_NOTE`; dua sub-jenis ⏸️ OQ-EDM-005 | ✅ |
| tombol **`Process Policy No`** b5233 `VIS=never` (`pyDisabledWhen` b5231) | — | ➖ mati |
| tiga `pxTextInput` `VIS=NEVER` (b587, b943, b1142, b1877) | — | ➖ mati |
| ➕ tombol **`View Old Policy`** (label VERBATIM `ShowLifePremiumSummary_EDM.xml` b64965) | bab 5 | ✅ ➕ dipindah dari layar yatim (R02, OQ-EDM-004) |

### 4b. Unggah CSV — wadah b8698 `.EdmType=1 && .EditInput=1`

| Korpus | Rule | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| tombol **`Upload CSV`** b8973 → `localAction` b9128 `UploadCSV_LifeEndorsement`; `pyDisabledWhen` b8965 `.EditInput1=1` | bab 6 | modal unggah | ✅ |
| tombol **`View Upload`** b9340 → `showHarness` b9844 `ViewCSVResult_LifeEDM` (popup) | bab 6 | popup hasil urai | ✅ |
| tombol **`Add CSV Data`** b10405 → `SaveCSVEDMLife` b10496; `pyDisabledWhen` b10403 `.EditInput1=1` | bab 6 | `POST /kasus/{id}/csv` | ✅ |

### 4c. Grid peserta

| Korpus | Data | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| wadah b11431 `pyWorkPage.EdmType=1 && .EditInput=1`; grid b11899 `.PremiumListSummary.PremiumListDetail` (`pyPageListProperty` b11860) | kepala `POLICY NO` b12035 · `POLICY HOLDER` b12177 · `CERTIFICATE NO` b12320 · `NAME OF INSURED` b12463 · `SEX` b12602 · `DATE OF BIRTH` b12738 · `ENTRY AGE` b12874 · `PLAN` b13010 · `BEGIN DATE` b13146 · `EFFECTIVE DATE` b13282 · `EXPIRED DATE` b13417; sel `.POLICY_NO` b13836 … `.EXPIRED_DATE` b15587 | `T_PREMIUM_LIST_DETAIL` versi kasus | ✅ |
| kotak centang `.EdmBatal` b15753, `pyDisabledWhen` b15763 `.EditInput==1` | tanda hapus per peserta | centang dikirim bersama `Save`; hanya peserta `Old` yang dapat dicentang — baris `Delete`/`New` terkunci (R30, OQ-EDM-017); mati sesudah simpan | ✅ |
| tombol **`DELETE ALL`** b13607 → `SelectAllEdmLife_act` b13688 | bab 6 | sakelar: seluruh peserta `Old` kasus tercentang (juga di halaman grid lain), centang yang lalu dilepas menjadi pengecualian | ✅ |
| `pyEditAction` `PL_DetailAction` (`InputEDMLife.xml` b16187, b20938) → FlowAction `PL_DetailAction` b170 `PL_Detail_Sec` | rincian peserta per `Type` (wadah `pyWorkPage.Type = 'QR'` b320, `'QP'` b12819, `'TR'` b21495, `'TP'` b33772) + grid retro (`Treaty Type` b11001, `Retroceded Share` b11147) | `GET /kasus/{id}/peserta/{pesertaId}`, `components/RincianPeserta.tsx` | ✅ |
| `PL_Detail_Sec.xml` `pyEditAction` `RetroLife` (b11653, b32606, b44911) → FlowAction `RetroLife` b94 `RetroDetailLife` | `Reinsurer Name` b1333 · `Percent Share (%)` b1484 · `Total Share` b1637 · `Rate` b1790 · `Gross Premium` b1943 · `1st Year Discount (%)` b2096 · `RI Admin Fee (%)` b2249 · `Net Premium` b2402 | `T_PREMIUM_LIST_SPREADING` + `_SPREADING_RETRO` versi kasus | ✅ |
| wadah b17034 `pyWorkPage.EdmType=3 && .EditInput=1`; grid b17500 | kepala `POLICY NO` b17636 · `CERTIFICATE NO` b17774 · `NAME OF INSURED` b17910 · `SEX` b18046 · `DATE OF BIRTH` b18182 · `ENTRY AGE` b18318 · `PLAN` b18454 · `BEGIN DATE` b18590 · `EFFECTIVE DATE` b18726 · `EXPIRED DATE` b18861 — **tanpa** kotak centang | seluruh peserta Batal | ✅ |

### 4d. Rekap mata uang — `.PremiumListSummary.CurrencyList` (b23025, b26038, b29037, b32036)

| Wadah | Kepala → sel | Keadaan |
| --- | --- | --- |
| b22598 `.Type=='QR'`, grid b23064 | `CURRENCY` b23197 → `.CURRENCY` · `GROSS PREMIUM` b23311 → `.PREMIUM` · `DEDUCTION` b23447 → `.DEDUCTION` · `BROKERAGE FEE` b23583 → `.BROKERAGE_FEE` · `NET PREMIUM` b23719 → `.BALANCE` | ✅ |
| b25611 `.Type=='QP'`, grid b26077 | `CURRENCY` b26210 · `GROSS PREMIUM REFUND` b26324 → `.PREMIUM` · `DEDUCTION REFUND` b26460 · `BROKERAGE FEE REFUND` b26596 · `NET PREMIUM REFUND` b26732 → `.BALANCE` | ✅ |
| b28610 `.Type=='TP'`, grid b29076 | `CURRENCY` b29209 · `GROSS PREMIUM RETRO` b29323 → `.PREMIUM` · `DISCOUNT PREMIUM RETRO` b29459 · `SHARE RETRO` b29595 · `NET PREMIUM RETRO` b29731 → `.BALANCE` | ✅ |
| b31609 `.Type=='TR'`, grid b32075 | `CURRENCY` b32208 · `GROSS PREMIUM REFUND RETRO` b32322 → `.PREMIUM` · `DISCOUNT PREMIUM REFUND RETRO` b32458 · `SHARE RETRO` b32594 · `NET PREMIUM REFUND RETRO` b32730 → `.BALANCE` | ✅ |

Tabel `T_PREMIUM_LIST_SUMMARY` (dihitung ulang dari peserta versi kasus, R08).

### 4e. Simpan — tombol **`Save`** b37202 `VIS .EditInput=1`, `pyDisabledWhen` b37200 `.IsJsonPolis=1` → `SetPremi_EDM` b37322

`Activity/SetPremi_EDM.xml` (9 langkah, nol `//`). Rute `POST /kasus/{id}/simpan`.

| Langkah | Isi | Sistem baru |
| --- | --- | --- |
| 1 b343 `·` | `Page-Remove` | — |
| 2 b551 PRE=false; 2.1 b551 `·` desc *"Set 0 jika EDM Batal"*, prakondisi b1273 `.EdmBatal=="True" \|\| pyWorkPage.EdmType==3` | 32 kolom uang `× -1` | jurnal balik **idempoten** dari nilai sumber (R06); 32 kolom VERBATIM di `models`; ✅ |
| 2.2 b1313 `·` prakondisi b1441 `.EdmBatal=="True"` → b1339 `.EditInput = 1`, b1385 `"Delete"` | — | `EDM_STATUS = 'Delete'`, baris terkunci; ✅ |
| 2.3 b1481 `·` prakondisi b1583 `pyWorkPage.EdmType==3` → b1507 `"Batal"` | — | seluruh peserta `Batal`; ✅ |
| 3 b1725 PRE=false (`>50000` b3373 diabaikan) | perulangan salin | ➖ batas 50.000 dibuang (spec penyimpangan 6) |
| 4 b3462 `·` (4.2 b3593 `Java`) | total `NET_PREMIUM` per `EDMStatus` ke `TempDetail2` | ➖ pembacanya yatim (R07) |
| 5 b4440 `·` b4467 `BusinessName` → `PremiumListSummary.BusinessName` | COB rekap | `BUSINESS_NAME` kasus; ✅ |
| 6 b4574 `·` `AppendCurrencySummary_DT` (b4623) | rekap mata uang | `T_PREMIUM_LIST_SUMMARY` dihitung ulang (R08); ✅ |
| 7 b4746 `·` | salin rekap sementara | — |
| 8 b5573 `·` b5600 `.IsJsonPolis = 1` | tanda sudah simpan | ada rekap kasus = sudah simpan; simpan kedua ditolak 409; ✅ |
| 9 b5707 `·` `Obj-Save` | — | satu transaksi; jejak `InputEDMLife` → `InputEDMLife` berkomentar `Save` (siapa/kapan jurnal balik); ✅ |

### 4f. Keputusan — `Section/ConfirmSection.xml` (wadah `InputEDMLife.xml` b35518 `.IsJsonPolis=1`)

| Korpus | Rule / data | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| radio `Status` b496 `pyWorkPage.EmailTypePL` → `postValue` b713, refresh `InputEDMLife` b727 | `IsLifeAccepted` b272 | radio `1` Accept / `2` Decline | ✅ (R14, OQ-EDM-006) |
| `Comment` b829 `.Description` | `AddHistorySuggest` b374 `CommentSuggest ← .Description` | isian komentar | ✅ |
| grid riwayat b1856 `.OfferFacIn.ViewSuggest` (b1818): `Date` b1984 · `PIC` b2125 · `Status` b2263 · `Comment` b2399; sel `.DateSuggest` b2550, `.PICSuggest` b2732, `.IsCedingConfirm` b2903, `.CommentSuggest` b3074 | `AddHistorySuggest` | `T_VIEW_SUGGEST` ber-`PREMIUM_LIST_ID` = kasus (riwayat penawaran NB **tidak** disalin, spec §16) | ✅ |
| tombol **`View Premium`** b36494 `VIS 1=2` → `EditDetail_Harness` b36673 → `EditDetail_Section` | — | — | ➖ mati |

### 4g. Submit — dua tombol berlabel sama

| Korpus | Rantai aksi | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| tombol **`Submit`** b37494 `VIS .IsJsonPolis=1 && pyWorkPage.EmailTypePL =1` | `AddHistorySuggest` b37776 → `GenerateNoEDM_Life` b37809 → `localAction` `SetJsonPolisEDMLife_Confirm` b37844 → `finishAssignment` b37918 → flow `Confirm` → `InsertJsonPolisLife_Act` | `POST /kasus/{id}/putuskan {status:"1"}`: riwayat, nomor EDM, versi resmi, salinan warisan, outbox — **satu transaksi** (bab 7) → modal terima kasih | ✅ |
| tombol **`Submit`** b38109 `VIS .IsJsonPolis=1 && pyWorkPage.EmailTypePL = 2 \|\|pyWorkPage.EmailTypePL = 7` | `AddHistorySuggest` b38314 → refresh b38349 → `finishAssignment` b38385 → flow `Decline` → `Resolved-Rejected` | `POST /kasus/{id}/putuskan {status:"2"}` → `STATUSS = 'Resolved-Rejected'` | ✅ |
| `Section/ConfirmSubmitEDM.xml` (FlowAction `SetJsonPolisEDMLife_Confirm` b91): teks b526 **`Thank you for Submit !`**, `No. Endorsement` b654 `.PremiumListSummary.PL_NUMBER_EDM`, tombol **`Close`** b1366 → `finishAssignment` b1487 | — | `components/TerimaKasih.tsx`; `Close` kembali ke kotak masuk | ✅ |

`Activity/AddHistorySuggest.xml` (2 langkah, nol `//`): 1 b233 PRE=false (prakondisi b426 `.FlagOnGoingPolicy==1` diabaikan) — `DateSuggest` b255 `@CurrentDateTime()` b256, `No` b310, `PICSuggest` b331 `OperatorID.pyUserName` b332, `IsCedingConfirm` b352 `@if(.EmailTypePL=1,"Accept","Decline")` b353, `CommentSuggest` b373 `.Description` b374; 2 b466 `·` `Obj-Sort` `.No` b493 `Descending` b491 → `T_VIEW_SUGGEST` (`NO`, `DATE_SUGGEST`, `PIC_SUGGEST` = akun pelaku, `IS_CEDING_CONFIRM`, `COMMENT_SUGGEST`), dibaca `NO DESC`. ✅

`Activity/GenerateNoEDM_Life.xml` (7 langkah, nol `//`): 1 b290 `·` (b448 `TempPolis.CARI4 ← PolicyNo`); 2 b529 `·` `ProdDateTime` (R16 ➖); 3 b680 PRE=false `GetProdKeOldData_SQL` (b733, `RDBList/GetProdKeOldData_SQL.xml` b84 `order by TGL_INPUT desc` → `PRODKE DESC`, R19); 4 b872 `·` `CARI4 = Prodke+1` (b946), `CARI14` pad dua digit bila satu digit (b967); 5 b1068 `·` prakondisi b1216 `PL_NUMBER_EDM==""` → `Generate_NoEndorsmentLife` (b1126, `RDBList/Generate_NoEndorsmentLife.xml` b84 `NOPOLIS||'/'||{InputData.CARI14}`); 6 b1264 `·` prakondisi b1390 → `PL_NUMBER_EDM` b1337; 7 b1430 `·` `Obj-Save`. → `NO_ENDORS` = `PL_NUMBER_EDM` = `<polis>/<NN>`, lahir sekali. ✅ (R23, OQ-EDM-008) *(ralat 01-10-2026, K3 keputusan work owner 01-10-2026: akhiran nomor = `PRODKE` Pega + 1 (kosong = 0) → NB warisan `/01`; versi tetap `NVL(PRODKE, 1)` — rumus lengkap di tiket 04 bab status 01-10-2026)*

## 5. Lihat polis lama — `View Old Policy` → `ViewOldPolicy_EDM` (`_QP`, `_TP`, `_TR`)

Rute `GET /kasus/{id}/polis-lama`. Komponen `components/PolisLama.tsx`.

| Korpus | Rule / data | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| `Section/ShowLifePremiumSummary_EDM.xml` tombol **`View Old Policy`** b64965 `VIS .Type=='QR'` → `GetOldDetail_EDM` b65218 → popup b65274 `ViewOldPolicy_EDM`; b65522 `'QP'` → b65774/b65806 `_QP`; b66083 `'TR'` → b66332/b66364 `_TR`; b66640 `'TP'` → b66890/b66922 `_TP` | layar inangnya yatim (R02) | satu tombol berlabel sama di kepala `InputEDMLife`, popup sesuai `Type` | ✅ ➕ (OQ-EDM-004) |
| `Activity/GetOldDetail_EDM.xml`: 1 b233 `//`; 2 b384 `·` `Obj-Open-By-Handle` `OldCaseID` (b436) ke `WorkLife`; 3 b527 `//`; 4 b829 `//`; 5 b963 `//` | buka versi lama | baca versi **sebelumnya** (sumber salinan kasus) — rule tidak ditiru, perilakunya ada (spec §4) | ✅ |
| `ViewOldPolicy_EDM.xml` wadah b929 `.Type=='QR'`, grid b1396 `WorkLife.PremiumListSummary.PremiumListDetail` (b1357): `NAME_OF_INSURED` b1536 · `CURRENCY` b1681 · `SUM_INSURED` b1826 · `CEDING_RETENTION` b1971 · `SUM_REASURED` b2116 · `SHARE_NUSANTARA_RE_GROSS` b2261 · `SUM_AT_RISK_GROSS` b2406 · `GROSS_PREMIUM` b2551 · `DEDUCTION` b2696 · `NET_PREMIUM` b2841 · `FACTOR` b2986 · `CLAIM_AMOUNT` b3131 | peserta versi lama | ✅ |
| `ViewOldPolicy_EDM.xml` grid b6999 (RD `BrowsePremiumList_RD` b10402; RD b800 `.PL_NUMBER = Param.PL_NUMBER`, urut b1435 `PL_NUMBER DESC`): `COB` b7144 · `PL NUMBER` b7282 · `CURRENCY` b7409 · `PREMIUM` b7510 · `DEDUCTION` b7620 · `BROKERAGE FEE` b7747 · `RI ADMIN FEE` b7857 · `TAX` b7993 · `PROF COMM` b8129 · `CLAIM` b8256 · `BALANCE` b8357 | rekap versi lama | rekap mata uang versi sebelumnya (`T_PREMIUM_LIST_SUMMARY`, atau `M_LIFE_PREMIUM_SUMMARY` `PL_NUMBER` + `IDPEGA` untuk sumber warisan) | ✅ |
| `ViewOldPolicy_EDM_QP.xml` grid b1435 (11 kolom: `NAME_OF_INSURED` b1575 … `CLAIM_AMOUNT` b3025) + rekap b6689 (16 kolom: `COB` b6834 … `BALANCE` b8762); wadah luar b674 `.Type = 'QR'` | — | grid dalam dibangun (R03) | ✅ |
| `ViewOldPolicy_EDM_TP.xml` grid b1465 (15 kolom: `NAME_OF_INSURED` b1605 … `NET_PREMIUM_RETRO` b3635) + rekap b8115 (13 kolom: `COB` b8260 … `BALANCE` b9815); wadah luar b704 | — | R03 | ✅ |
| `ViewOldPolicy_EDM_TR.xml` grid b1450 (19 kolom: `NAME_OF_INSURED` b1590 … `CLAIM_AMOUNT` b4200) + rekap b9496 (13 kolom: `COB` b9641 … `BALANCE` b11196); wadah luar b689 | — | R03 | ✅ |
| `ShowLifePremiumSummary_EDM.xml` kepala `Endorsement Life Summary` b989, grid detail + rekap per Type (b5268 … b59639), `Production Date` b64353, tombol **`Submit`** b67733 → `InsertJsonPolisLife_Act` b67924 + `SetJsonPolisEDMLife_Confirm` b67957 | layar yatim (R02) | — | ➖ tak terjangkau |

## 6. CSV — `UploadCSV_LifeEndorsement`, `ViewCSVResult_LifeEDM`, `SaveCSVEDMLife`, `SelectAllEdmLife_act`

Rute `POST /kasus/{id}/unggah` (tinjau, nol tulis), `POST /kasus/{id}/csv` (multipart `berkas`, tanpa batas baris, dua lintasan atas
berkas sementara). Komponen `components/UnggahCSV.tsx`, `components/HasilCSV.tsx`, pengurai tampilan `csv.ts`. Baris acuan
`PremiumListDetail(1)` = peserta pertama di urutan grid, di luar baris `New`.

| Korpus | Isi | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| `FlowAction/UploadCSV_LifeEndorsement.xml` b191 judul **`Endorsement Life - Upload CSV`**, `pyLocalActionActivity` b144 `UploadCSVEDMLifePremium_Act` | wizard `pxUploadCSVResults` | modal unggah berjudul sama | ✅ |
| `Activity/UploadCSVEDMLifePremium_Act.xml`: 1 b252 `·`, 2 b342 `·`, 3 b490 `·` `pxUploadCSVResults` (b537 `ResultsClass` `ASM-FW-GISFW-Data-BatchLifePremiumDetail`, b539 `.ListLifePremiumDetailUpload`) | urai CSV ke daftar sementara | `POST /unggah` mengurai berkas dan **mengembalikan** barisnya — tidak disimpan (daftar sementara Pega hanya di clipboard) | ✅ |
| idem 4 b633 `//`, 5 b779 `//`, 5.2 b2109 `//`, 5.4 b2389 `//`, 6 b2513 `//`; 5.1 b872 dan 5.3 b2260 `·` di bawah induk `//` | simpan ke tabel batch | ➖ ter-remark (spec §15) |
| `Section/ViewCSVResult_LifeEDM.xml` wadah b616 `SELALU` (kondisi `1 = 2` diabaikan), grid b1082 `TempWorkPage.ListLifePremiumDetailUpload` (b1043), 38 kepala `POLICY NO` b1218 … | hasil urai | popup tabel 38 kolom, label VERBATIM | ✅ |
| tombol **`Generate Data Detail`** b14322 → `GenerateDataDtlLife_act` b14415 (`Activity/GenerateDataDtlLife_act.xml` 1 b223 `·` `pxConvertResultsToCSV`, b270 `FileName = DetailUpload`, b272 stempel waktu, b276 `CSVPropHeaders`) | unduh CSV hasil urai | unduhan CSV di peramban, kepala persis b276 | ✅ |
| `Activity/SaveCSVEDMLife.xml` 1 b264 `·`; 2 b361 `·` + 2.1 b454 `·` prakondisi b547 `.EDMStatus=="New"` → hapus baris `New` | unggah ulang mengganti | baris `New` kasus dihapus lebih dulu (AC 35) | ✅ |
| idem 3 b629 `·` `call ASM-FW-GISFW-Work-LIFE.Calculate1_Act` | mesin hitung NB | ➖ tidak dijalankan (R05, OQ-EDM-003) *(01-10-2026 K2 keputusan work owner 01-10-2026: tetap tidak dihitung — gerbang langkah 1 b416 `IsCalculationSystem` tak terbukti, OQ-EDM-021)* |
| idem 4 b726 `·` (b800 pesan `"Plan di CSV tidak sesuai, mohon di cek kembali"`); 4.1 b840 `·` prakondisi b2466 `.PLAN=…PremiumListDetail(1).PLAN`, b2495 `.POLICY_HOLDER=…(1).POLICY_HOLDER`, 72 penetapan `@divide(@toDecimal(@replaceAll(.X,",",".")),1,20)`; 4.2 b2541 `·` pesan; 4.3 b2689 `·` prakondisi b2791 `pyWorkPage.EdmType==1` → b2715 `"New"` | validasi + pemetaan | seluruh baris divalidasi lebih dulu, pesan VERBATIM + nomor baris/kolom (R24); desimal tanpa `float`, koma = titik desimal; `EdmType=3` ditolak (AC 13) | ✅ |
| idem 5 b2873 `·` b2899 `.EditInput1 = 1` | kunci unggah | ada baris `New` = unggah terkunci di layar (`InputEDMLife.xml` b8965/b10403) **dan** di server (409); sesudah `Save` rekap dihitung ulang (R31) | ✅ |
| `Activity/SelectAllEdmLife_act.xml` 1 b232 `·` b259 `@if(Select.CARI1=="","true",@if(Select.CARI1=="true","false","true"))`; 2 b431 + 2.1 b431 PRE=false (prakondisi b533 `.STS_REJECT !=""` diabaikan) → `.EdmBatal` b456 | balik centang semua | sakelar di layar atas seluruh baris yang belum terkunci | ✅ |

## 7. Penyimpanan resmi — flow `Confirm` → `Activity/InsertJsonPolisLife_Act.xml` (16 langkah, 2 `//`)

Satu transaksi di `POST /kasus/{id}/putuskan` (spec §11), `services/edm_putusan.go`; urutan langkah dikunci uji (AC 42). Tabel: `T_PREMIUM_LIST` (kolom EDM + `NO_POLIS`, `PROD_KE`, `NO_ENDORS`, `PL_NUMBER_EDM`,
`STATUSS`), `T_PREMIUM_LIST_DETAIL` (`PL_NUMBER_EDM`, `STATUS_OLD`, `STATUS`), **`M_LIFE_PREMIUM_DETAIL`**, **`M_LIFE_PREMIUM_SUMMARY`** (E2), outbox
`T_LOG_SERVICE_RNM` (E5), jejak `T_CLAIMLF_JEJAK`.

| Langkah | Isi | Sistem baru |
| --- | --- | --- |
| 1 b468 `·`, 2 b707 `·` (prakondisi b818 tanggal `>25`) | `ProdDateTime` | ➖ R16 |
| 3 b858 PRE=false `GetProdKeOldData_SQL` (b911); 4 b1050 `·` `Prodke+1` (b1124), pad (b1145) | `PRODKE` | sama dengan `GenerateNoEDM_Life`, `PRODKE DESC`; ✅ |
| 5 b1246 `·` prakondisi b1394 `PL_NUMBER_EDM==""` `Generate_NoEndorsmentLife` (b1304); 6 b1442 `·` prakondisi b1568 | nomor lahir sekali | ✅ |
| 7 b1608 `·` (b1635 `@ASM.GetPageJSONString()`, b1682 `IsJsonPolis = 0`); 9 b2486 `·` `InsertJsonPolisEDM` (b2543) | blob JSON ke `JSON_POLIS` | ➖ dibuang (spec §16, AC 54; R18) |
| 8 b1805 `·`; 10 b2663 `·` `SaveLifeinProduction_SQL` (b2720) | `LIFEINPRODUCTION` | ➖ R17, OQ-EDM-010 |
| 11 b2889 PRE=false (`>50000` b5296 diabaikan) atas `pyWorkPage.PremiumListSummary.PremiumListDetail` (b2847); 11.1 b2889 PRE=false (b4456 `EDMStatus ← .EDMStatus`); 11.2 b4546 `·` prakondisi b4644 `Type=="QR" \|\| Type=="QP"` → `CARI48` b4567 `0`; 11.3 b4684 `·` b4782 `TR`/`TP` → b4706 `1`; 11.4 b4822 `·` b4920 `.EDMStatus=="Old"` → `CARI47` b4843 `1`; 11.5 b4960 `·` b5058 `!="Old"` → b4982 `0`; 11.6 b5098 PRE=false `SaveMasterLPDet` (b5154) | seluruh peserta kasus, termasuk `Old`/`New`/`Delete`/`Batal` | baris `M_LIFE_PREMIUM_DETAIL` per peserta, 80 kolom `RDBList/SaveMasterLPDet.xml` b87, `EDMSTATUS` ← `{TempValue.EDMStatus}` b246, `STATUSOLD` ← `{TempInputDetail.CARI47}` b247 (`1` bila `Old`), `STATUS` ← `{TempInputDetail.CARI48}` b248 (`0` QR/QP, `1` TR/TP), `PL_NUMBER` b222, `PL_NUMBER_EDM` b223, `IDPEGA` ← `{pyWorkPage.pzInsKey}` b245 = `EDMLF-<n>`; hapus-sebelum-sisip `PL_NUMBER`+`IDPEGA` (`_INDEX4`) — ⏸️ **belum ditulis**: benturan penjaga Claim Life (R29), menunggu OQ-EDM-016; peserta versi resmi sudah di `T_PREMIUM_LIST_DETAIL` (`PL_NUMBER_EDM`, `STATUS_OLD`, `STATUS`) ✅ |
| 12 b5385 `·` atas `CurrencyList` (b5343); 12.1 b5385 `·`; 12.2 b6171 `·` `InsertPLSummary` (b6226) | rekap per mata uang | baris `M_LIFE_PREMIUM_SUMMARY` per mata uang (37 kolom prosedur VERBATIM, `CARIn` ← properti bernama sama, prosedur tidak dipanggil), hapus-sebelum-sisip `PL_NUMBER`+`IDPEGA`; ✅ E2 |
| 13 b6387 `//` `GetNopolisByIDPega` (b6444) | cek sudah masuk | ➖ ter-remark; alarm **dihidupkan** (spec penyimpangan 3) dengan pemicu bergeser (R33): simpan atomik — keadaan separuh mustahil — alarm menyala pada kegagalan efek keluar, pesan di layar + log tanpa alamat; baca-balik `_INDEX21` menunggu peserta warisan (OQ-EDM-016); ✅ |
| 14 b6564 `·` b6591 pesan galat `… Error Insert Data to Json_polis (InsertJsonPolisLife_Act)` | pesan galat JSON | ➖ penulis JSON dibuang |
| 15 `//` `SendEmailNotification` | email | ➖ ter-remark, **isinya rahasia — tidak dikutip**; alarm E5 menggantikannya |
| 16 b7162 `·` prakondisi b7230 `IsPEGAPROD` → `serviceInsertArasapasLife_act` | efek keluar | `services/edm_efekkeluar.go`: sesudah commit, gerbang lingkungan `outbox.Penyalur` (satu tempat), kunci `M_LINK_SERVICE` (`Activity/serviceInsertArasapasLife_act.xml` b793 `"Production"`, b795 `"convertJsonNusareToProduction"`) di-resolve sungguhan, Arasapas lalu alarm **hanya bila gagal**, kegagalan ke outbox `T_LOG_SERVICE_RNM` `MODUL = ENDORSEMENTLIFE`; kedua efek stub gagal terang (OQ-EDM-013); ✅ E5 |

`Activity/serviceInsertArasapasLife_act.xml` (7 langkah, 1 `//`): 1 b253 `·` (b280 `pzInsKey`); 2 b387 `·` `GetPolicyNoByCaseId` (b445,
`RDBList/GetPolicyNoByCaseId.xml` b85 `json_polis`); 3 b564 `·`; 4 b748 `·` `GetLinkService` (b793 `Kategori_1 = "Production"`, b795
`Kategori_2 = "convertJsonNusareToProduction"`); 5 b864 `·` `Connect-REST` (b912 `convertJsonNusareToProduction`, b916 `POST`); 6 b1015 `·`;
7 b1130 `//` (prakondisi b1278 `IsSuccessHitService`). → stub `EfekArasapasEDM` me-resolve alamat lewat `inti/backend/layanan` lalu gagal
terang; nomor polis dari kasus (`OLD_POLICY_NO`), bukan `JSON_POLIS`. ✅ ⏸️ panggilan nyata OQ-EDM-013

| Rule pendukung | Isi | Keadaan |
| --- | --- | --- |
| `Activity/GetLinkService.xml` (4 langkah `·`, `Obj-Browse` b371 `.KATEGORI_1` b492, `.KATEGORI_2` b518) | baca `M_LINK_SERVICE` | ✅ lewat `inti/backend/layanan` |
| `ConnectREST/ConvertJsonNusareToProduction.xml` peta `.OfferFacIn.PolicyData.PolicyNo` b199, `.pzInsKey` b207, `.pxCreateDateTime` b218; respons `.StatusService` b333 | muatan tiga pengenal; respons ke `STATUS_SERVICE` | ⏸️ stub (E5) |
| `When/IsPEGAPROD.xml` b174 `pxProcess.pzProductionLevel = "5"` | gerbang produksi | ✅ `inti.Lingkungan.AdalahProduksi()` |
| `When/recordEvent.xml` (dirujuk tiga harness) | peristiwa UI | ➖ tanpa padanan |

## 8. Alur — `Flow/InputEDMLife.xml` dan `DecisionTable/IsLifeAccepted.xml`

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Start1` → `ASSIGNMENT63` (b1313 `pyTo`), assignment `InputEDMLife` b731, `pyRouteTo` `Current operator` b759 | kasus terbuka di layar `InputEDMLife` | ✅ |
| konektor b1045 `pyExpression` `InputEDMLife` → `Decision1` `Accept` b960 (`IsLifeAccepted`) | `POST /putuskan` | ✅ |
| `IsLifeAccepted` b272 `pyWorkPage.EmailTypePL`, kondisi b284 `1` → b317 `Confirm`; bawaan b86 `Decline` — **tanpa `Reject`** | `status = "1"` → Confirm, selain itu Decline; nol jalur kembali ke input (AC 30) | ✅ |
| `Confirm` b1192 → `Utility1` `InsertJsonPolisLife_Act` b576 → `END52` `Resolved-Completed` b686 | bab 7, `STATUSS = 'Resolved-Completed'` | ✅ |
| `Decline` b1408 → `End1` `Resolved-Rejected` b642 | `STATUSS = 'Resolved-Rejected'` | ✅ |
| `Assignment1` `Input EDM Summary` b859, konektor keluar b1273 (`pyExpression` `InputEDMLife_Summary` b1254) — nol konektor masuk | — | ➖ yatim (R02) |

## 9. Rule korpus yang tidak dibawa

| Rule | Alasan |
| --- | --- |
| `Activity/Calculate1_Act.xml` beserta `RDBList` yang hanya dirujuknya: `CekDoubleInsured`, `InsertDataUploadLife`, `DeleteTempUploadDataLife`, `GetProductDtlPL`, `GetRateLifePM`, `SelectComm_SQL` (`M_TEMPUPLOADLIFE`, produk, rate, komisi) | R05 — mesin hitung NB tidak dijalankan di jalur endorsement (keputusan work owner, AC 53) |
| `RDBList/InsertJsonPolisEDM.xml`, `RDBList/SaveLifeinProduction_SQL.xml` | R17, R18 |
| `DataTransform/SetCoBName_Act.xml` (dirujuk sel ubah `Type`/`Reinsurance System`/`Premium Method` yang baca-saja) | aksi ubah tidak pernah jalan pada sel `ro` |
| RD `BrowseCedingCoLife_RD`, `BrowseClientNusaRe_RD`, `BrowseMarketingOfficer_RD` | pemilih di sel baca-saja |
| `Section/EditDetail_Section.xml` + `Harness/EditDetail_Harness.xml` | hanya dibuka `View Premium` (`InputEDMLife.xml` b36494, `1=2`) |
| `Section/ShowLifePremiumSummary_EDM.xml`, FlowAction `InputEDMLife_Summary` | R02 (kecuali tombol `View Old Policy`, bab 5) |

## 10. Sensus tombol — 22 sel `pxButton`, label VERBATIM

| # | Section | Label | `bNNN` | Keadaan |
| ---: | --- | --- | --- | --- |
| 1 | `ConfirmSubmitEDM` | `Close` | b1366 | ✅ |
| 2 | `EndorsmentLife_Section` | `Process Policy No` | b2414 | ➖ `VIS=never` |
| 3 | `EndorsmentLife_Section` | `Submit` | b4226 | ✅ |
| 4 | `InboxEndorsementLife` | `Create` | b3780 | ➖ R11 |
| 5 | `InboxEndorsementLife` | `Cancel` | b4228 | ➖ R11 |
| 6 | `InboxEndorsementLife` | `Create Case Endorsement` | b5879 | ➖ `VIS=never` |
| 7 | `InboxEndorsementLife` | `Create Addendum` | b6620 | ✅ |
| 8 | `InputEDMLife` | `Process Policy No` | b5233 | ➖ `VIS=never` |
| 9 | `InputEDMLife` | `Upload CSV` | b8973 | ✅ |
| 10 | `InputEDMLife` | `View Upload` | b9340 | ✅ |
| 11 | `InputEDMLife` | `Add CSV Data` | b10405 | ✅ |
| 12 | `InputEDMLife` | `DELETE ALL` | b13607 | ✅ |
| 13 | `InputEDMLife` | `View Premium` | b36494 | ➖ `VIS 1=2` |
| 14 | `InputEDMLife` | `Save` | b37202 | ✅ |
| 15 | `InputEDMLife` | `Submit` (`EmailTypePL =1`) | b37494 | ✅ |
| 16 | `InputEDMLife` | `Submit` (`EmailTypePL = 2 \|\| 7`) | b38109 | ✅ |
| 17 | `ShowLifePremiumSummary_EDM` | `View Old Policy` (`QR`) | b64965 | ✅ ➕ dipindah (R02) |
| 18 | `ShowLifePremiumSummary_EDM` | `View Old Policy` (`QP`) | b65522 | ✅ ➕ dipindah (R02) |
| 19 | `ShowLifePremiumSummary_EDM` | `View Old Policy` (`TR`) | b66083 | ✅ ➕ dipindah (R02) |
| 20 | `ShowLifePremiumSummary_EDM` | `View Old Policy` (`TP`) | b66640 | ✅ ➕ dipindah (R02) |
| 21 | `ShowLifePremiumSummary_EDM` | `Submit` | b67733 | ➖ R02 |
| 22 | `ViewCSVResult_LifeEDM` | `Generate Data Detail` | b14322 | ✅ |

Aksi XML di luar tiket yang dibangun (tambahan tiket): rincian peserta `PL_DetailAction` + retro `RetroLife` (bab 4c), riwayat keputusan
`T_VIEW_SUGGEST` (bab 4f), unduh `Generate Data Detail` (bab 6), alarm efek keluar (bab 7, R33). Sensus ini dikunci `frontend/sensus.test.ts`.
