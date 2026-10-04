# Perbandingan kolom — delapan tabel NB Treaty In lawan diagram grilling

*Disusun 03-10-2026 (putaran 2, paket penyimpanan). Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir
11–12 dan bab 2 K4, K16, K17; PESAN-KOREKSI-PUTARAN-2.*

**Acuan (urutan menang):** (1) diagram grilling `Diagram-Skema-Tabel-NusantaraRe.xlsx` sheet
*NB Treaty In Prop* dan *NB Treaty In NonProp* (sel dikutip dengan alamatnya, mis. `R43`);
(2) `docs/rancangan-tabel-datar-treaty-in.md` (§4.1 inti, §4.2 XOL, §4.3 angsuran, §4.4 spreading);
(3) XML korpus `NB Treaty In (Done)` — hanya untuk kolom **di luar** (1) dan (2).

**Aturan keputusan (bab 0 butir 12):**

| Keadaan | Keputusan |
| --- | --- |
| di diagram atau rancangan | **dipertahankan** (atau ditambah bila belum ada) |
| di luar keduanya, XML membuktikan medannya **DIBACA** rule NB terjangkau (syarat, rumus, sel Section) | **dipertahankan** — `RALAT` rancangan §4sexies + bukti XML |
| di luar keduanya, medannya hanya **ditulis**, atau pembacanya hanya langkah berlabel `//`, atau hanya jalur treaty keluar (K8 butir 4) | **dibuang** |
| dicoret / ditetapkan lain oleh diagram atau keputusan work owner | **tidak dibuat**, sebab dikutip |

> ⛔ **RALAT 04-10-2026** `[keputusan work owner]` — bunyi lama, dikutip: *"Hasil: **tepat delapan `CREATE TABLE`** di
> migrasi modul (320–327)"*; baris tabel 1 *"`T_GENERAL_POLIS` | 320 | … | 1:1 shared PK `T_WORK_POLIS`"*; bab 1a
> *"`ID` … dipertahankan | PK + FK `T_WORK_POLIS`"* dan *"`IDPEGA` · `TGL_INPUT` · `USERNAME` … dipertahankan"*. Bunyi
> baru: `T_GENERAL_POLIS` adalah SATU tabel **bersama FacIn + Treaty In**; tabel dasarnya (`ID VARCHAR2(32)` PK,
> `IDPEGA VARCHAR2(50)`, `COB_GROUP`, `START_DATE_TIME`, `OFFERING_DATE`, `END_DATE_TIME`, `FOLLOWING`) dibuat
> migrasi `nbfacin` `182_t_general_polis`; 320 hanya `ALTER TABLE … ADD (` kolom Treaty (semua kolom bab 1 kecuali
> `ID` dan `IDPEGA`, yang dipakai dari tabel dasar) + constraint/indeks Treaty. Daftar kolom Treaty **tidak berubah**.
> Dasar: keputusan WO 04-10-2026 (mengalahkan bab 0 butir 11 / K18 PROMPT putaran 3); PERMINTAAN-TIM-INTI C10.

Hasil: **tepat tujuh `CREATE TABLE`** (`T_POLIS_*`, 321–327) + **`T_GENERAL_POLIS` bersama** lewat `ALTER ADD` (320),
ditagih `repository/kolom_test.go` `TestTabelDanKolomMengikutiDiagramGrilling` (nama tabel, jumlah tabel, dan
kolom persis — untuk `T_GENERAL_POLIS`: kolom ALTER 320 + `ID`, `IDPEGA` dasar), `backend/migrasi_test.go`, dan
`TestSQLTidakMenyebutKolomYangDibuang`. Tabel putaran 1 di luar diagram — 328
`T_POLIS_SUGGEST`, 329 `T_POLIS_MEDAN_LAIN`, 330 `M_NBTRIN_PERAN_TEMPAT` — **dihapus** (K4, K17, K16).

| # | Tabel | Migrasi | Sheet / sel diagram | Bentuk |
| ---: | --- | --- | --- | --- |
| 1 | `T_GENERAL_POLIS` **bersama** (dasar `nbfacin` 182) | 320 `ALTER ADD` | Prop/NonProp F9–F33 | 1:1 shared PK `T_WORK_POLIS`, satu baris per generasi; baris FacIn (`LINI='FAC'`) di tabel yang sama, tidak dibaca modul ini |
| 2 | `T_POLIS_QUOTATION` | 321 | J35–J38 | 1:1 `POLIS_ID` |
| 3 | `T_POLIS_CEDING` | 322 | O39, R40–R50 | 1:N `QUOTATION_ID` |
| 4 | `T_POLIS_INSTALMENT` | 323 | Prop J52–J56 · NonProp J52–J55 | 1:N `POLIS_ID` |
| 5 | `T_POLIS_INSTALMENT_DETAIL` | 324 | NonProp O56, R57–R61 | 1:N `INSTALMENT_ID` · NonProp saja |
| 6 | `T_POLIS_SPREADING` | 325 | Prop J66–J69 · NonProp J71–J74 | 1:N `POLIS_ID` |
| 7 | `T_POLIS_XOL` | 326 | NonProp J76–J80 | 1:N `POLIS_ID` · NonProp saja |
| 8 | `T_POLIS_XOL_LAYER` | 327 | NonProp O81, R82–R87 | 1:N `XOL_ID` · NonProp saja |

`T_WORK_POLIS` (akar, B5) sudah ada — milik premiumlistlife, tidak dibuat.
`POOLDATA.HISTORYAKSEPTASIPRODUCTION` (Prop J74–J78, NonProp J89–J93) — tabel **warisan**, ditulis dan
dibaca, tidak dibuat (MODUL.md *Tabel warisan*), lihat bab 9.

---

## 1 · `T_GENERAL_POLIS` — 79 medan `PolicyTreatyIn` + 7 kolom `json_polis` (F10–F11)

### 1a · Kunci, generasi, `json_polis`

| Kolom | Asal | Diagram / rancangan | Keputusan | Bukti |
| --- | --- | --- | --- | --- |
| `ID` | `pyWorkPage.pzInsKey` → `T_WORK_POLIS.ID` | F9 *SHARED PK* | dipertahankan — **kolom dasar `nbfacin` 182** (RALAT 04-10-2026) | PK milik 182; FK `T_WORK_POLIS` ditambah 320 |
| `NOPOLIS` | `PolicyTreatyIn.PolicyNo` / json_polis | F11–F12 kunci alami | dipertahankan | `UNIQUE (NOPOLIS, PRODKE)` (F16) = indeks unik fungsi `CASE` — `UNIQUE` biasa menganggap dua draf `(NULL, 0)` kembar di Oracle |
| `PRODKE` | json_polis | F12, F22 cacah | dipertahankan | `NUMBER(10) DEFAULT 0`; NB = 0 |
| `NOENDORS` | json_polis | F11, rancangan §2 | dipertahankan | |
| `OLD_POLIS_ID` | penunjuk generasi | F13–F16 | dipertahankan | FK `T_WORK_POLIS`, `UNIQUE (OLD_POLIS_ID)` |
| `IDPEGA` · `TGL_INPUT` · `USERNAME` | json_polis | F11, rancangan §2 | dipertahankan; `IDPEGA` = **kolom dasar `nbfacin` 182** `VARCHAR2(50)` (RALAT 04-10-2026) | `USERNAME` = identitas login (P4); pelebaran `IDPEGA` diminta (C10) |
| `TGL_PROD` | `PolicyTreatyIn.ProductionDate` / json_polis | F11, rancangan §4.1 *PRODUCTION_DATE* | dipertahankan | satu kolom untuk keduanya (ID-21) |
| ~~`TGL_TUTUP`~~ | — (bukan medan Pega) | ✗ | **dibuang** | ID-10: generasi lampau tidak disunting, *"ditegakkan di services, bukan di tabel"*. Tertutup kini = ada baris penerus yang `OLD_POLIS_ID`-nya menunjuk generasi ini (`repository.syaratTerbuka`) |

### 1b · Medan `PolicyTreatyIn` (rancangan §4.1) — dipertahankan

Seluruhnya ada di rancangan §4.1 (penentu bentuk, penanda alur, tanggal, uang, persen, kode dan teks,
pihak) atau §4q5.7 (`REMARK`):

`NO_OFFER` · `MASTER_ID` · `IS_APPROVED` · `SUGGEST` · `SUGGEST_DATE` · `OPERATOR_NAME` ·
`IS_NEW_POLICY_NON_PROP` · `HAS_FAC_OUT` · `FLAG_PPH` · `FLAG_RETRO_TREATY` · `DUE_TO` · `TYPE_TAX` ·
`STATEMENT_TYPE` · `TREATY_GROUP_ID` · `TREATY_GROUP_NAME` · `TREATY_GROUP_OLD_ID` · `OJK_BUSINESS_ID` ·
`BIZ_CODE` · `BIZ_NAME` · `SOB` · `SOB_NAME` · `INSURED_ID` · `INSURED_NAME` · `MARKETING_OFFICER` ·
`TREATY_TYPE` · `TREATY_YEAR` · `CURRENCY` · `ID_CURRENCY` · `QUARTAL` · `YEAR_OF_QUARTAL` · `CLAIM_TYPE` ·
`CLAIM_PAYMENT_TYPE` · `INSTALLMENT` · `REMARK` · `START_DATE` · `END_DATE` · `STATEMENT_DATE` ·
`GROSS_PREMIUM` · `PREMI_OGP` · `RESULT_OGP1` · `RESULT_OGP2` · `PREMI_ONP` · `RESULT_ONP1` · `RESULT_ONP2` ·
`CLAIM` · `OUTSTANDING_CLAIM` · `SALVAGE_VALUE` · `EXCESS_LOSS` · `NET_PREMIUM` · `BALANCE_DUE_TO` ·
`BALANCE_BEFORE_TAX` · `BALANCE_BEFORE_PPH` · `DEDUCTION1` · `DEDUCTION2` · `PPH_VALUE` · `PPN_VALUE` ·
`SHARE_VALUE` · `RI_COMM_OGP` · `OVERIDDING_COMM_OGP` · `RI_COMM_ONP` · `OVERIDDING_COMM_ONP`.

| Kolom | Catatan |
| --- | --- |
| `CEDING_CO` · `CEDING_CO_NAME` | ⭐ **di sini, tingkat polis** (R47): disalin apa adanya dari `PolicyTreatyIn.CedingCo/CedingCoName` (preACT langkah 3), termasuk ekor `"; "` — tidak dirangkai ulang dari `T_POLIS_CEDING`. Salinan `Quotation.CedingCo/CedingCoName` **dibuang** dari `T_POLIS_QUOTATION` (bab 2) |
| `IS_EDM_INPUT_ON_NB` | **ditambah** — rancangan §4.1 *penanda alur*; ditulis dan dibaca `InputPolicyTreatyInDetail_NonProp` langkah 10–11 (jalur XOL, K8) |
| `ID_NEW_BISNIS` | **ditambah** — rancangan §4.1 *kode dan teks*; data guide `$.IDNewBisnis`. Nol rule NB menyentuhnya — kolom untuk medan dokumen lama (pemuat tiket 22) |
| `DEDUCTION1` · `DEDUCTION2` | golongannya urusan paket layar (K3) — tidak diubah paket ini |

### 1c · Di luar rancangan — **dipertahankan** dengan `RALAT` (dibaca rule terjangkau)

| Kolom | Medan | Bukti XML (rule terjangkau, langkah / sel) |
| --- | --- | --- |
| `POSITION_NOTE` | `pyWorkPage.PositionNote` | halaman kerja, bukan `PolicyTreatyIn`; `T_WORK_POLIS` (milik premiumlistlife) tidak punya kolomnya. Dibaca connector `Flow/InputRealizationTreatyIn`, `InboxPolicyTreatyIn_postDT` langkah 4, `DeptHeadTreatyIn_UW_postDT`, syarat sel `ListSuggest` (`pyWorkPage.PositionNote != 'ReasTreatyInSecHead'`) |
| `NB_STATUS` | `pyWorkPage.NBStatus` | tampil `Section/SFAPortal_OpportunitiesList` sel `A.NBStatus`; kolom RD `GetListOpportunity` |
| `TREATY_IN_ID` | `pyWorkPage.TreatyIn.ID` | dibaca `RDBList/BrowseTreatyIn` (`where ID={TreatyIn.ID}`, jalur XOL K8) dan `CheckDuplicateOffer` (`.ID != TreatyIn.ID`); master kontrak dan `LAYER*` dimuat ulang dari view lewat nilai ini |
| `SHARE_CURRENCY` | `PolicyTreatyIn.ShareCurrency` | tampil `DetailPolicyTreatyIn` dan `DetailDeptHeadTreatyIn_UW` (label *RNM Share*) — AC 64 |
| `GROSS_CLAIM` | `PolicyTreatyIn.GrossClaim` | tampil `DetailPolicyTreatyIn` (*Claim 100%*); dibaca `CalculatePremi_Act` (`@divide(.GrossClaim*…RNMShareP,100,4)`) |
| `BROKERAGE_FEE_SEBENARNYA` | `PolicyTreatyIn.BrokerageFeeSebenarnya` | dibaca rumus `SetPPNPPH` langkah 3 (`.PPHValue = .BrokerageFeeSebenarnya * @divide(2,100,8)`, `.PPNValue = … 2.2 …`) |
| `EDM_TYPE` ⭐ *(putaran 3, F3)* | `PolicyTreatyIn.EDMType` | medan dokumen lama (data guide `$.EDMType`) **dibaca** prasyarat `Activity\InputPolicyTreatyInPre_Act.xml` langkah 10 (`.PolicyTreatyIn.EDMType=="3"` → lewati `Call TreatyRealizationCheckXOLList`; terjangkau, `models.PerluCekDaftarXOL`). Rancangan §4.1 *penentu bentuk* memuat `EDM_TYPE`; properti lain dari `QuotationData.EdmType` (`T_POLIS_QUOTATION.EDM_TYPE`, J37). RALAT rancangan §4sexies |

> ⛔ **RALAT 04-10-2026 (F3)** — baris `BROKERAGE_FEE_SEBENARNYA` di atas, bunyi lama: *"dibaca rumus `SetPPNPPH`
> langkah 3"*. Bunyi baru: langkah **4** (`Activity\SetPPNPPH.xml` `REPEATINGINDEX="4"`, prasyarat
> `.FlagPPH=="true"` dan `ListAgent.pxResults(1).STS_PKP == 1`); langkah 3 = `pxRetrieveReportData`.

### 1d · Di luar rancangan — **dibuang**

| Kolom lama | Medan | Bukti |
| --- | --- | --- |
| `IS_OJK_NOPOLIS` | `PolicyTreatyIn.IsOJKNopolis` | hanya **ditulis** `GeneratePolicyNoTreaty_Act` langkah 23; nol pembaca terjangkau. ⚠️ ada di dokumen lama (data guide) → masuk laporan CSV pemuat (butir terbuka) |
| `BROKERAGE_FEE` | `PolicyTreatyIn.BrokerageFee` | hanya **ditulis** `SetPPNPPH` langkah 3; nol pembaca |

> ⛔ **RALAT 04-10-2026 (F3)** — dua baris di atas. Bunyi lama: *"hanya **ditulis** `GeneratePolicyNoTreaty_Act`
> langkah 23 … ⚠️ ada di dokumen lama (data guide) → masuk laporan CSV pemuat (butir terbuka)"* dan *"hanya
> **ditulis** `SetPPNPPH` langkah 3"*. Bunyi baru: langkah 23 `GeneratePolicyNoTreaty_Act` **berlabel `//`**
> (`pyStepsBlockName`), jadi XML kini bahkan tidak menulisnya; `.BrokerageFee` ditulis `SetPPNPPH` langkah **4**.
> Keduanya tetap tanpa kolom dan, sebagai medan dokumen lama, **dibuang** dengan bukti di
> `backend/models/medan_abaikan_lama.json` bagian `pola` (keputusan WO F3) — tidak lagi butir terbuka.
> (Di `docs/dataguide-json-polis.json` `IsOJKNopolis` hanya muncul di bawah `OldData` — `$.OldData.IsOJKNopolis`,
> salinan generasi sebelumnya yang seluruhnya dibuang; keputusan tingkat polis berlaku bila muncul di data nyata.)

### 1e · Medan `PolicyTreatyIn` diagram yang **tidak** menjadi kolom di sini

| Medan | Sebab |
| --- | --- |
| `PolicyNo` | menjadi kunci `NOPOLIS` (ditulis `SetelNomorPolis`, sekali — AC 74) |
| `ProportionalType`, `EdmType` (rancangan §4.1 *penentu bentuk*) | **dipindah** ke `T_POLIS_QUOTATION` — diagram J37 memuatnya di sana (`QuotationData`) |
| `Layer` · `LayerType` · `LayerPart` · `LayerPartType` | **dicoret** diagram F26: pantulan `pxResults(1)`; dibaca balik dari view lewat `TREATY_IN_ID` setiap layar dibuka (`models.TerapkanMasterKontrak`, ID-22) |
| `TotalPremium` · `TotalClaim` · `TotalSharePercentagePremium/Claim` | turunan baris spreading (`HitungTotalSpreading`); penjaga repo melarang `TOTAL_` di migrasi (RALAT AC 38 lama) |
| `isApprovedtoDeptHead` | P36, spec AC 64 — tidak dibangun |
| `EDMNo` | P55, terhitung di EDM |

⇒ Hitungan: 69 kolom katalog medan `PolicyTreatyIn` (termasuk `TGL_PROD`) + `NOPOLIS` + 9 medan tak
berkolom (4 `LAYER*`, 4 `Total*`, `isApprovedtoDeptHead`) = **79** — angka diagram F10 (batas bawah, B119).
⭐ Putaran 3 (F3): + `EDM_TYPE` (bab 1c) ⇒ 70 kolom katalog medan `PolicyTreatyIn`; 79 tetap batas bawah.

---

## 2 · `T_POLIS_QUOTATION` — `QuotationData`, 10 medan (J36–J37)

| Kolom | Medan | Diagram | Keputusan | Bukti |
| --- | --- | --- | --- | --- |
| `POLIS_ID` | kunci | G34 *1:1 POLIS_ID UNIK* | dipertahankan | PK + FK `T_GENERAL_POLIS` |
| `PROPORTIONAL_TYPE` · `MO_ID` · `BUSINESS_CODE` · `BUSINESS_OLD_ID` · `GROUP_PANEL` · `SOURCE_OF_BUSINESS` · `TYPE` · `EDM_TYPE` · `OLD_POLICY_NO` · `MARKETING_NAME` | 10 medan | ✓ J37 | dipertahankan | `SOURCE_OF_BUSINESS` dibaca `SetPPNPPH` langkah 1–3 (status PKP agen) dan ditulis pemilih SOB (CATATAN §2c) |
| `BUSINESS_NAME` | `Quotation.BusinessName` | ✗ | **dipertahankan — RALAT** | `InputPolicyTreatyInPre_Act` langkah 2 (syarat `@contains(…BusinessName,"MBU")` dst. + `OldID.CARI2`, berjalan setiap layar dibuka bila `BizCode` kosong); tampil `SFAPortal_OpportunitiesList` sel `A.Quotation.BusinessName` |
| `BUSINESS_FAC` | `Quotation.BusinessFac` | ✗ | **dipertahankan — RALAT** | `SaveViewSuggest` langkah 2 `CARI7` (→ `HISTORYAKSEPTASIPRODUCTION.TYPE`); `GetListOpportunity` filter E `A.Quotation.BusinessFac = T`. Bernilai `T` sejak kasus lahir |
| `INSURED_ID` | `Quotation.InsuredID` | ✗ | **dipertahankan — RALAT** | `InputPolicyTreatyInDetail_preACT` langkah 3 (`PolicyTreatyIn.InsuredID = Quotation.InsuredID`) dan 14.3 (syarat `Quotation.InsuredID==""` — nilai pilihan bisnis sebelumnya bertahan) |
| `INSURED_NAME` | `Quotation.InsuredName` | ✗ | **dipertahankan — RALAT** | preACT langkah 3 dan 14.1 (`SearchClient.CARI1`); tampil `SFAPortal_OpportunitiesList` `A.Quotation.InsuredName` |
| `NO_OFFER_SLIP` | `QuotationData.NoOfferSlip` | ✗ | **dipertahankan — RALAT** | tampil `DetailPolicyTreatyIn` (pxTextArea, diisi admin) dan `DetailDeptHeadTreatyIn_UW` — AC 64 |
| `IS_SURVEY_REPORT` | `QuotationData.IsSurveyReport` | ✗ | **dipertahankan — RALAT** | tampil + **wajib** `DetailPolicyTreatyIn`; syarat nonaktif tombol *Survey Report* |
| ~~`BUSINESS_TYPE`~~ | `Quotation.BusinessType` | ✗ | **dibuang** | hanya ditulis DT `BusinessType_DeT` (preACT 14.8); nol pembaca. Turunan masukan yang tersimpan (`GROUP_PANEL` + `BUSINESS_OLD_ID`, J38; `models.GolongkanJenisUsaha`) |
| ~~`SOB_NAME`~~ · ~~`SOB_LEADER0`~~ · ~~`SOB_LEADER1`~~ | `Quotation.SobName/SobLeader0/1` | ✗ | **dibuang** | ditulis preACT langkah 3 dan `SearchHierarkiSourceBizAgent_PostDT` (pemilih SOB); pembacanya hanya `CheckDataMkt` langkah 7 (label `//`) dan `InputPolicyTreatyOutDetail_preACT` (treaty keluar, K8 butir 4). Dikonfirmasi paket pemilih SOB: nol pembaca |
| ~~`CEDING_CO`~~ · ~~`CEDING_CO_NAME`~~ | `Quotation.CedingCo/CedingCoName` | ✗ (R47: di `T_GENERAL_POLIS`) | **dibuang** (tingkat polis tetap di `T_GENERAL_POLIS`) | ditulis preACT langkah 3; nol pembaca NB (`TreatyRealizationCheckDuplicate` membaca `PolicyTreatyIn.CedingCoName`) |
| ~~`MARKETING_CODE`~~ | `Quotation.MarketingCode` | ✗ | **dibuang** | ditulis `CheckDataMkt` langkah 4; pembacanya hanya langkah 6 (label `//`) |
| ~~`TEAM_GROUP`~~ · ~~`BRANCH_CODE`~~ · ~~`BRANCH_NAME`~~ | `Quotation.TeamGroup/BranchCode/BranchName` | ✗ | **dibuang** | hanya ditulis `CheckDataMkt` langkah 4 |

---

## 3 · `T_POLIS_CEDING` — `QuotationData.CedingCoList()`, 2 medan (R41–R43)

| Kolom | Diagram | Keputusan | Bukti |
| --- | --- | --- | --- |
| `ID` | — | dipertahankan | kunci baris |
| `QUOTATION_ID` | O39 *1:N QUOTATION_ID* | **dipindah** dari `POLIS_ID → T_GENERAL_POLIS` ke `QUOTATION_ID → T_POLIS_QUOTATION (POLIS_ID)` | diagram; repository menghapus anak **sebelum** quotation ditulis ulang |
| `NOURUT` | R42 | dipertahankan | `UNIQUE (QUOTATION_ID, NOURUT)` |
| `CEDING_CO_ID` | R43 *CEDING_CO_ID ← .CedingCo* | **diganti nama** dari `CEDING_CO` | diagram |
| `CEDING_CO_NAME` | R43 | dipertahankan | |

---

## 4 · `T_POLIS_INSTALMENT` — `ListInstallment()`, 10 medan (J53)

Rancangan §4.3 (sebelas kolom) tanpa `PAYMENT_DATE` (diagram R61: *"punya PAYMENT_DATE yang tidak ada di
induk"*) = sepuluh.

| Kolom | Keputusan | Bukti |
| --- | --- | --- |
| `ID` · `POLIS_ID` · `NOURUT` | dipertahankan | J54 |
| `INSTALLMENT_NO` · `DUE_DATE` · `INSTALLMENT_PERCENTAGE` · `PREMIUM` · `PAYMENT_TOTAL` · `CURRENCY` · `ID_CURRENCY` | dipertahankan | rancangan §3.1/§4.3; tampil grid `DetailPolicyTreatyIn` |
| `PREMIUM_AFTER_PPH` · `PREMIUM_AFTER_PPN` · `PREMIUM_AFTER_TAX` | **ditambah** | rancangan §4.3; data guide `$.ListInstallment.PremiumAfter*`. Nol penulis di korpus NB — kolom untuk dokumen lama |
| `PPN` · `PPH` · `PAYMENT_TOTAL_AFTER_PPN` · `PAYMENT_TOTAL_AFTER_TAX` | **dipertahankan — RALAT** | ditulis lalu **dibaca rumus** `InputPolicyTreatyInDetail_preACT` langkah 18.3.4.1 (`.PaymentTotalAfterPPN = .PaymentTotal+.PPN`; `Local.PPNins = .PPN`; `Local.PPHins = .PPh`; `Local.PaymentTotalAfterPPN/Tax = …`); ada di dokumen lama (data guide) |
| ~~`PAYMENT_DATE`~~ | **dibuang** | diagram R61; nol rujukan `ListInstallment().PaymentDate` di korpus NB |

## 5 · `T_POLIS_INSTALMENT_DETAIL` — `ListInstallment().InstallmentList()`, 11 medan (R58)

| Kolom | Keputusan | Bukti |
| --- | --- | --- |
| `ID` · `INSTALMENT_ID` · `NOURUT` | dipertahankan | O56, R59 |
| `INSTALLMENT_NO` · `DUE_DATE` · `PAYMENT_DATE` · `INSTALLMENT_PERCENTAGE` · `PREMIUM` · `PAYMENT_TOTAL` · `PREMIUM_AFTER_PPH` · `PREMIUM_AFTER_PPN` · `PREMIUM_AFTER_TAX` · `CURRENCY` · `ID_CURRENCY` | dipertahankan / `PREMIUM_AFTER_PPH` **ditambah** | rancangan §4.3, diagram R58/R61; `PremiumAfterPPN/Tax` tampil `Section/InstallmentList`; ditulis `InputPolicyTreatyInDetail_NonProp` 20.4.1 dan preACT 18.3.4.2.1 |
| ~~`PPN`~~ · ~~`PPH`~~ | **dibuang** | hanya ditulis preACT 18.3.4.2.1; nol pembaca. ⚠️ ada di dokumen lama (data guide) → laporan CSV pemuat (butir terbuka) |

> ⛔ **RALAT 04-10-2026 (F3)** — baris di atas, bunyi lama: *"⚠️ ada di dokumen lama (data guide) → laporan CSV
> pemuat (butir terbuka)"*. Bunyi baru: medan dokumen lama `ListInstallment().InstallmentList().PPN/PPh`
> **dibuang** dengan bukti (`backend/models/medan_abaikan_lama.json` bagian `pola`: pembaca `.PPN/.PPh` di preACT
> 18.3.4.1 adalah baris `ListInstallment` — kolom bab 4); tercatat di arsip CSV pemuat, bukan butir terbuka.
> (Di `docs/dataguide-json-polis.json` keduanya hanya muncul di bawah `OldData` —
> `$.OldData.ListInstallment.InstallmentList.PPN/PPh`; keputusan tingkat polis berlaku bila muncul di data nyata.)

## 6 · `T_POLIS_SPREADING` — `SpreadingRiskList()`, 9 medan (J67)

`ID` · `POLIS_ID` · `NOURUT` + `TREATY_TYPE` · `TREATY_NAME` · `CURRENCY` · `CURRENCY_ID` ·
`SHARE_PERCENTAGE` · `SPLIT_RNM_SHARE_PCT` · `CLAIM_PERCENTAGE` · `PREMIUM_SPREADED` · `CLAIM_SPREADED` —
**sama persis** dengan rancangan §4.4 (sembilan). Tidak berubah.

## 7 · `T_POLIS_XOL` — `TreatyXOLList()`, 13 medan (J77)

`ID` · `POLIS_ID` · `NOURUT` + `CURRENCY` · `ID_CURRENCY` · `GROSS_PREMI` · `NET_PREMI` · `DEDUCTION` ·
`DUE_TO` · `DUE_TO_VALUE` · `BROKERAGE_FEE_SEBENARNYA` · `PPH_VALUE` · `PPN_VALUE` · `NET_PREMI_AFTER_PPH` ·
`NET_PREMI_AFTER_PPN` · `NET_PREMI_AFTER_TAX` — rancangan §4.2 tanpa `LAYER*` (J79: *induk tidak punya
penanda layer*). Tidak berubah. Logika jalurnya milik paket XOL.

## 8 · `T_POLIS_XOL_LAYER` — `TreatyXOLList().ValueList()`, 17 medan (R83)

`ID` · `XOL_ID` · `NOURUT` + `LAYER` · `LAYER_TYPE` · `LAYER_PART` · `LAYER_PART_TYPE` (R85: *ADA DI SINI
SAJA*) + 13 kolom `T_POLIS_XOL`. Tidak berubah.

---

## 9 · `POOLDATA.HISTORYAKSEPTASIPRODUCTION` — tabel WARISAN (J74–J78, J89–J93)

Tidak dibuat, tidak diubah strukturnya. Catatan `PolicyTreatyIn.SuggestList` ditulis ke sini
(`repository/usulan.go`, pengganti `RDBList/InsertViewSuggest_SQL`) dan dibaca balik untuk layar
`Section/ListSuggest` (K4). ⭐ Putaran 3 (F3): `SuggestList` **dokumen lama** juga disalin ke sini oleh pemuat
(tiket 22; `repository.SalinUsulanLama`, penjaga dobel menurut `IDPEGA`). Pemetaan 15 kolom: `docs/STRUKTUR-TABEL-NB-TREATY-IN.md` bab tabel ini dan
`backend/models/usulan.go`.

---

## 10 · AC 64 — setiap medan yang tampil di Section NB terjangkau

Sapuan sel `Section` terjangkau (korpus, `graf.Graf.terjangkau()`), dicocokkan ke kolom:

| Medan tampil | Section | Tempat simpan |
| --- | --- | --- |
| skalar `PolicyTreatyIn` layar admin/atasan (`BalanceBeforePPH` … `YearOfQuartal`) | `DetailPolicyTreatyIn`, `DetailDeptHeadTreatyIn_UW` | `T_GENERAL_POLIS` (bab 1b/1c) |
| `PolicyNo` · `ProductionDate` | `DetailDeptHeadTreatyIn_UW`, `ShowPolicyNoTreaty_SC`, `ListSuggest` | `NOPOLIS` · `TGL_PROD` |
| `.QuotationData.IsSurveyReport` · `.NoOfferSlip` · `.MOID` · `.ProportionalType` | admin, atasan | `T_POLIS_QUOTATION` |
| `A.Quotation.BusinessName` · `InsuredName` · `MarketingName` · `A.NBStatus` | `SFAPortal_OpportunitiesList` | `T_POLIS_QUOTATION` · `NB_STATUS` |
| `Layer` · `LayerType` · `LayerPart` · `LayerPartType` | admin, atasan | ⛔ dicoret F26 — dibaca balik dari view saat dibuka |
| `TotalPremium` · `TotalClaim` · `TotalSharePercentage*` | admin, atasan, `SpreadingRiskList` | turunan baris spreading saat dimuat |
| `TreatyType` · `SharePercentage` · `PremiumSpreaded` · `ClaimPercentage` · `ClaimSpreaded` | `SpreadingRiskList` | `T_POLIS_SPREADING` |
| `InstallmentNo` · `DueDate` · `InstallmentPercentage` · `Premium` · `PaymentTotal` | grid angsuran admin/atasan | `T_POLIS_INSTALMENT` |
| `DueDate` · `InstallmentPercentage` · `Currency` · `Premium` · `PremiumAfterPPN` · `PremiumAfterTax` | `InstallmentList` | `T_POLIS_INSTALMENT_DETAIL` |
| `Date` · `IsApproved` · `OperatorName` · `Suggest` | `ListSuggest` (grid) | `HISTORYAKSEPTASIPRODUCTION` (`TGL_INP` · `APPROVAL` · `PIC` · `KETERANGAN`) |
| `Limit*` · `MDP*` · `Deductible*` · `NetPremi*` · `Note` · `Value` · `Total*Value` | `DetailPolicyTreatyInNonProportional` | halaman master `pyWorkPage.TreatyIn.*SummaryList` (K8, baca-saja JSON master) — bukan medan polis |
| `pyWorkPage.TreatyIn.*`, `Installments_ReadOnly`, `BusinessAndSOBList` | | master kontrak / view — bukan medan polis |
| `DateofSurvey` · `SurveyedBy` · `Remarks` · `LossPrevention` | `HistoricalSurveyReport*` | ⛔ K7 — tidak dibangun, tidak ada tabel diagram |

⇒ Nol medan polis tampil tanpa tempat simpan, kecuali yang dicoret diagram (F26) atau turunan.

---

## 11 · Butir terbuka

1. Medan dokumen lama tanpa kolom (mis. `IsOJKNopolis`, `InstallmentList.PPN/PPh`, dan `QuotationData`
   di luar sepuluh + enam RALAT: `BranchCode`, `BusinessType2`, `SobLsg`, `StatusBusiness`, …, data guide)
   akan muncul di **laporan CSV** pemuat (K17) — wajib 0 sebelum selesai, jadi butuh keputusan work owner
   (tambah kolom lewat RALAT diagram, atau nyatakan dibuang).
   ✅ **Diputuskan 04-10-2026 (F3)** — butir ini selesai: (a) `PolicyTreatyIn.EDMType` dibaca rule NB terjangkau
   → kolom `T_GENERAL_POLIS.EDM_TYPE` (bab 1c); (b) 25 pola lain dibuang dengan alasan + bukti XML per medan
   (`backend/models/medan_abaikan_lama.json` bagian `pola`; daftar dan bukti: tiket 22 bab *Putaran 3*);
   `SuggestList` disalin ke `HISTORYAKSEPTASIPRODUCTION` (bab 9). CSV pemuat = arsip audit, bukan penampung;
   yang wajib 0: medan **belum diputuskan** (RALAT spec-penyimpanan AC 59, ID-27).
2. Tipe fisik kolom `HISTORYAKSEPTASIPRODUCTION` (terutama `NOURUT`, `PERCENT_RNM`) belum dicek ke
   katalog Oracle; pembacaan memakai `TO_NUMBER(NOURUT)`.
3. `HISTORYAKSEPTASIPRODUCTION.DIV` (`OperatorID.pyOrgDivision`) tanpa sumber di `inti.Pelaku` — NULL.
