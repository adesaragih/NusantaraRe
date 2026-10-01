# Paritas layar dan aksi — Master Product Name Life

> Disusun 01-10-2026 (sesi implementasi modul, paket 0, brief `PROMPT-IMPLEMENTASI-MODUL-MASTER-PRODUCT-NAME-LIFE.md`).
> Satu baris per harness, section, tombol, dan aksi korpus → rule yang dipanggil → RDB/RD → tabel/kunci JSON → rute API → komponen React.
> Nomor baris `bNNN` = baris berkas korpus sesudah `sed -e 's/></>\n</g'` atas `D:\XML\RNM_BRD\Master Product Name Life\` (konvensi repo).
> Setiap langkah activity yang dikutip mencetak `pyStepsBlockName` — `·` = kosong (langkah hidup), `//` = ter-remark (tidak pernah jalan).
> Prakondisi dibaca dari urutan aksi: `WhenTrue 2` = lanjut, `3` = lewati langkah; `T=-` = kosong = lanjut; **`PRE=false` = prakondisi
> dimatikan → langkah SELALU jalan**, apa pun isi WHEN-nya.
> Keadaan: ✅ dibangun · ⏸️ dibangun, datanya menunggu OQ · ➖ sengaja tidak dibawa (bukti di kolom Catatan) · ⏳ menunggu paket N.
> Kunci JSON ditulis **peka huruf besar-kecil** persis Pega (`POLICYHODER`, `UnderwritingLimitList`, `Non_Employee`, ...).

## 0. Alat baca dan sensus

| Hal | Isi |
| --- | --- |
| Alat | pengurai token bernomor baris `pega.py` (scratchpad sesi, tidak di repo): `akt` pohon langkah + `pyStepsBlockName`, `tombol`/`sel` sel section + aksi + rantai visibilitas wadah + sumber grid, `dt`, `rd`, `rdb`, `fa`; ditambah peta pemanggil (`>Nama<` di seluruh 114 berkas) |
| Uji alat | nomor baris alat = nomor `sed` (sel tombol pertama `InboxProductName.xml` b5181 = baris `<pyFormat>pxButton` mentah); `ProteksiPlanListLife` = 2.1.1–2.1.3 bersarang (cocok XML) |
| Berkas | 114: Activity 37, RDBList 22, ReportDefinition 17, Section 12, FlowAction 11, DataTransform 10, Harness 1, DecisionTable 1, When 1, ConnectREST 1, SystemSettings 1 (+ `Struktur_InboxProductName.xlsx`, bukan rule) |
| Langkah ber-`//` | **27** — cara 1 alat (jumlah per activity), cara 2 `grep -c '<pyStepsBlockName>//'` mentah: sepakat. Sebaran: `ProductNameSaveAttachment` 7, `DownloadAttProdName_Act` 5, `CheckDuplicateOffer` 4, `DeleteAttacProdName_act` 3, `SaveInwardProductName_Act` 2, `SetProductName` 2, `LoadAttachment` 2, `LoadAttachmentProdName` 1, `SetCategoryAttachTreatyin` 1 |
| Tombol `pxButton` | **42** — cara 1 `grep -c '<pyFormat>pxButton</pyFormat>'`, cara 2 sel `Embed-Display-Table-Cell` ber-`pyFormat = pxButton`: sepakat per berkas (`InboxProductName` 33, tujuh pemilih masing-masing 1, `InwardProductName` section 1, harness 1) |
| Tautan/ikon beraksi | **4** di `InboxProductName`: `pxIcon` b39632, b45290 (salin baris), `pxLink` b68857 (nama berkas), b69247 (`View Office Online`) |
| Teks tombol | `pyModes[2].pyLabel` (keterangan yang dirender `pxButton`/`pxLink`); `pyLabelFieldValue` = label SEL (`Button`/`End Period`), bukan teks tombol |
| Medan visibilitas | sel: `pyUserData.pyVisible` (`ALWAYS` → syarat diabaikan; `OTHER`/`NOTBLANK` → syarat berlaku); wadah: `pyContainerVisibleWhen` berlaku hanya bila `pyIsVisibilityOption` = `CONDITION` atau `ExpressionCondition` (62 `ALWAYS`, 14 `CONDITION`, 2 `ExpressionCondition` di korpus) |

## 1. Halaman awal — pintu masuk modul

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Section/InboxProductName.xml` mode daftar: wadah b71246 `OutputParam.DATASHOW != 1` (`CONDITION`) berisi grid produk; mode form: wadah b2573 / b30995 / b57867 / b61044 / b64133 `OutputParam.DATASHOW = 1` | tombol menu **Master Product Name Life** → halaman `mpnl-produk` (grid produk; form dibuka dari dalam) | ⏳ paket 10 (menu) / 2 (layar) |
| ⭐ **Bukti `InboxProductName` = pintu masuk**: tidak dibuka rule mana pun di korpus (peta pemanggil: hanya dirujuk sebagai sasaran `refresh otherSection` oleh enam section pemilih — `Ceding_Section.xml` b2285, `SOB_Section.xml` b2300, `PolicyHolder_Section.xml` b2318, `Currency_Section.xml` b2323, `RIRISK_Section.xml` b2307, `CauseOfLoss_Section.xml` b2242), jadi dibuka portal/navigasi di luar ekspor. Satu-satunya harness, `Harness/InwardProductName.xml`, dibuka **hanya** oleh tombol `Inward` b75322 yang bervisibilitas `OTHER 1=2` (mati) | halaman awal = `InboxProductName`; harness `InwardProductName` **tidak dibangun** | ralat spec §1 *"Titik masuk `InwardProductName` … satu-satunya Harness"* (RALAT R11) |

## 2. Grid daftar produk — `InboxProductName` mode daftar (wadah b71246)

Sumber: RD `BrowseProduct_Life` (b76658, kelas `ASM-FW-GISFW-Int-PRODUCT_LIFE` = view `PRODUCT_LIFE`), tanpa saringan (b681), urut `.ID ASC` (b1094), 10 baris/halaman (`pyRDLPageSize` b71371), `pyGridPaginator` b72105.
Sistem baru membaca `M_PRODUCT_LIFE.ID` + `JSONDATA` sendiri (P2) — kunci kolom grid = kunci JSON yang dibaca view.

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → tabel/kunci | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kolom `ID` b72406 · `Ceding` b72554 · `Treaty Number` b72711 · `Treaty Name` b72870 · `Create Operator` b73029 · `Last Updated Operator` b73188 → `.ID` b73628 · `.CEDING` b73813 · `.TREATYNUMBER` b73990 · `.INWARDNAME` b74182 · `.CREATEOP` b74374 · `.UPDATEOP` b74566 | — | `M_PRODUCT_LIFE.ID`; JSON `CEDING`, `TREATYNUMBER`, `INWARDNAME`, `CREATEOP`, `UPDATEOP` | `GET /api/master-product-name-life/produk` | ⏳ paket 1 (baca) / 2 (layar) |
| ⚠️ Wadah grid b71574 `pyContainerVisibleWhen = IsFire` | — | When `IsFire` tidak ada di korpus | — | ✅ ikut XML: `pyIsVisibilityOption = ALWAYS` → kondisi diabaikan, grid selalu tampil |
| Tombol b71816: label sel **`End Period`**, teks **`Add`**, tooltip **`Add New Data`**, vis `ALWAYS !IsSpreadingUW` (syarat diabaikan) | `NewProductLife` b71870/b71980: 1 b232 `·` Page-New `ProductName`; 2 b376 `·` `DATASHOW=1`, `ERRMSG=""`, `HASILD7=0`, `STSSAVE=0`; 3 b577 `·` Page-Remove `ProductName`, `ProductNameDisplay`, `ProductNameInward` | — | form baru (klien) | ⏳ paket 3 |
| Tombol baris **`View`** b74753 | `SetProductName` b74802 (param `ID ← .ID`): 1 b468 `·` hapus halaman; 2 b583 `·` `DATASHOW=1`, `ProductName.ID ← Param.ID`; 3 b780 `·` RDB-List `BrowseUnderwritingList`; 4 b957 `·` Java `adoptJSONObject(BENEFIT)`; 5 b1065 **`//`**; 6 b1204 **`//`** (salin param lama); 7 b1716 `·` call `SetProductNameInward(ProductID ← ProductName.ID)`; 8 b2145 `·` `IsView = true` — lalu refresh `LoadAttachmentProdName` b74954 | `BrowseUnderwritingList` b61 `select JSONDATA as BENEFIT from m_product_life where ID={ProductName.ID}`; `SetProductNameInward` 1 b650 `·`, 2 b826 `·` pxShowReport RD `BrowseProductInward` (saringan b834 `.PRODUCTID = Param.ProductID`), 3 b997 `·` ulang, 3.1 b1046 `·` salin **40** medan view | `GET …/produk/{id}` | ⏳ paket 1 / 4 |
| Tombol baris **`Inward`** b75322 (vis `OTHER 1=2`) | `showHarness` `InwardProductName` popup + `SetProductNameInward` | — | — | ➖ **mati** (`1=2`); harness dan `SaveInwardProductName_Act` tak terjangkau (RALAT R11) |

## 3. Form produk — `InboxProductName` mode form (wadah b2573, `DATASHOW = 1`)

Pesan: wadah b616 `STSSAVE==100` → `OutputParam.ERRMSG` b881; wadah b1269 `STSSAVE==99` → `ERRMSG` b1534 (dipakai `CopyProduct`); wadah b1921 `ProductName.ERRMSG!=''` → b2186.
Baca-saja: setiap medan ber-`ro = ProductName.IsView=='true'` (mode lihat sesudah `View`; `Edit` membukanya).

### 3.1 Sisi umum — halaman `ProductName` → `M_PRODUCT_LIFE.JSONDATA`

| Medan korpus (bNNN · label VERBATIM) | Kontrol · visibilitas · aksi | Kunci JSON | Keadaan |
| --- | --- | --- | --- |
| b3143 `ProductName.TYPE` · b3349 `TYPE_CEDING` · b6108 `GRUP` (`Group`) | radio, vis `OTHER 1=2` — **mati**; `TYPE_CEDING` onChange `CountMaxReasured_Act` b3385 ikut mati | `TYPE`, `TYPE_CEDING`, `GRUP` (dibaca view `PRODUCT_LIFE`) — ditulis apa adanya dari JSON lama, `""` untuk produk baru | ➖ medan; kunci ⏳ paket 3 |
| b3652 **`Product Name`** `ProductName.PRODUCTNAME` | dropdown `associated` (daftar lokal Property, **tidak ikut ekspor**); onChange `SetTreatyName_Act` b3686: 1 b237 `·`, 2 b434 `·` `INWARDNAME ← PRODUCTNAME + " " + POLICYHODERNAME` | `PRODUCTNAME` (+ kolom datar `PRODUCTNAME`) | ⏳ paket 3 — isian teks (OQ-MPNL-05) |
| b3928 **`Product Code`** `ProductName.ID` | teks, vis `ALWAYS` | `ID` | ⏳ paket 3 — **baca-saja** (ADR-0006; RALAT R14) |
| b4107 **`Ceding`** `ProductName.CEDING` + tombol **`Choose Ceding Name`** b5181 | autocomplete RD `BrowseCedingCoLife_RD` b4203; tombol `localAction ChooseCeding` b5352 (wadah b4945 `IsView!='true'`) | `CEDING`, `CEDINGID` | ⏳ paket 2 |
| b4495 **`SOB`** `ProductName.SOBNAME` + **`Choose SOB`** b5522 | autocomplete RD `BrowseCedingCoLife_RD` b4586; `localAction ChooseSOB` b5693 | `SOBNAME`, `SOBID` | ⏳ paket 2 |
| b6315 `Plan` (`PRODUCTTYPE`) · b6722 `R/I Rate` (`RIRATE`) · b9113 `Outward Name` · b9531 `Outward Rate` · b9927 `Outward Comm` · b10312 `Benefit` | autocomplete, vis `OTHER 1=2` — **mati** | `PRODUCTTYPE`, `RIRATE`, `OUTWARDNAME`, `OUTWARDRATE`, `OUTWARDCOMM`, `BENEFIT` (+ `*ID`) dibaca view — dipertahankan dari JSON lama / `""` | ➖ medan; kunci ⏳ paket 3 |
| b7137 **`Deduction (%)`** `ProductName.RICOMM` | angka | `RICOMM` | ⏳ paket 3 |
| b7430 **`R/I Risk Name`** `ProductName.RIRISK` + **`Choose R/I Risk`** b8164 | autocomplete RD `BrowseRIRiskSummary` b7520; `localAction ChooseRIRisk` b8335 | `RIRISK`, `RIRISKID` (+ kolom datar `RIRISK`, `RIRISKID`) | ⏳ paket 2 |
| b8753 **`Treaty Name`** `ProductName.INWARDNAME` · b8934 **`Treaty Number`** `ProductName.TREATYNUMBER` | teks | `INWARDNAME`, `TREATYNUMBER` | ⏳ paket 3 |
| b10693 **`Cause Of Loss`** `ProductName.CAUSE` + **`Choose Cause of Loss`** b11316 | dropdown RD `BrowseCauseofLossLife_RD` b10781 (nilai/tampil `.CauseofLoss`); `localAction ChooseCauseOfLoss` b11487 | `CAUSE`, `CAUSEID` | ⏳ paket 2 |
| b47303 checkbox **`On Retention`** (`pyCheckboxCaption` b47312) `ProductName.IsORS`, vis `ALWAYS` (di luar wadah `1==2`) | onChange `GetReinsTypeOR_Life` b47476: 1 b234 `·` **PRE=false** hapus `OutwardList`; 2 b381 `·` `Temp.CARI1/2 ← BEGIN/MATURE`; 3 b532 `·` **PRE=false** RDB-List `BrowseReinstypeOR_SQL`; 4 b723 `·` **PRE=false** ulang; 4.1 b768 `·` tambah baris `OutwardList` | `BrowseReinstypeOR_SQL` b84: `TREATYCONTRACT_LIFE tc JOIN TREATYYEAR_LIFE ty ON ty.id = tc.idtreatyyear WHERE TO_DATE(BEGIN,'DD/MM/YYYY') >= tc.treatystartdate AND TO_DATE(MATURE,'DD/MM/YYYY') <= tc.treatyenddate AND REINSTYPEID = '10200'` → `IsORS`, `OutwardList[*].{REINSTYPEID, REINSTYPENAME, TRANSACTIONYEAR←TREATYYEAR, TREATYCONTRACTID, UNDERWRITINGYEAR}` | ⏳ paket 9 — ⭐ **hidup** (RALAT R9, P3) |

### 3.2 Sisi inward — halaman `ProductNameInward` → `M_PRODUCTINWARD_LIFE.JSONDATA`

Diisi di form yang **sama** (bukan harness `InwardProductName`). Disimpan `SaveProductName_Act` langkah 13–16 (§7).

| Medan korpus (bNNN · label VERBATIM) | Kunci JSON | Keadaan |
| --- | --- | --- |
| b17129 **`Policy Holder`** autocomplete RD `BrowseClientNusaRe_RD` b17224 + **`Choose Policy Holder`** b17786 (`localAction ChoosePolicyHolder` b17957); onChange `SetTreatyName_Act` b17168 | `POLICYHODER` (id), `POLICYHODERNAME` (disalin juga ke JSON umum, `SaveProductName_Act` 1 b359) | ⏳ paket 2 / 4 |
| b18376 **`Insured`** · b18866 **`Addendum No.`** · b19466 **`Addendum`** · b20189 **`Amandement No.`** · b20794 **`Amandement`** · b21813 **`Max Notification Claim Expired`** | `INSURED`, `ADDENDUMNO`, `ADDENDUMWORD`, `AMANDEMENTNO`, `AMANDEMENTSCHD`, `MAXEXPIREDCLAIM` | ⏳ paket 4 |
| b22001 **`Begin Date`** · b22336 **`STNC`** · b27282 **`Expired Date`** | `BEGIN`, `STNC`, `MATURE` — teks `dd/MM/yyyy` (bentuk `TO_DATE(…,'DD/MM/YYYY')` `BrowseReinstypeOR_SQL` b84, contoh `[data DBA]`) | ⏳ paket 4 |
| b22526 **`Ceding Retention (%)`** · b22689 **`Ceding's Limit`** (onChange `CountMaxSumReasured_Act` b22725) · b22947 **`Brokerage Fee (%)`** · b23110 **`Minimum Age (Years)`** · b23297 **`Maximum Age (Years)`** · b23484 **`Expiry Age (Years)`** · b23667 **`Extra Premium`** · b23828 **`Min Sum Insured`** · b23988 **`Max Sum Insured`** (onChange `CountMaxSumReasured_Act` b24024) · b24246 **`Max Sum Reasured`** · b24705 **`Nusantara Re Share (%)`** + label b24882 **`OF SUM REASURED`** · b25268 **`Nusantara Re's Limit`** | `CEDINGRETENTIONNUM`, `CEDINGLIMIT`, `BROKERAGE`, `MINAGE`, `MAXAGE`, `EXPIRYAGE`, `EXTRAPREMI`, `MINSUMINSURED`, `MAXSUMINSURED`, `MAXSUMREASURED`, `RNMSHARE`, `RNMLIMITNUM` | ⏳ paket 4 |
| `CountMaxSumReasured_Act`: 1 b237 `·` `local ← MAXSUMINSURED, CEDINGLIMIT`; 2 b434 `·` `MAXSUMREASURED = MaxSumInsured − CedingLimit` | `MAXSUMREASURED` dihitung, tetap dapat ditimpa (medan b24246 dapat disunting) | ⏳ paket 4 (hitung desimal persis di klien) |
| b25430 **`Premium Factor (%)`** (vis `OTHER PAYMENT==3`) · b25642 **`Premium Payment Method`** (dropdown `associated`) · b25952 **`Subject To`** (area teks) · b26124 **`Annuity Interest (%)`** · b26338 **`Premium Refund Factor (%)`** (keduanya `ALWAYS`, syarat `PAYMENT==3` diabaikan) | `PREMIUMFACTOR`, `PAYMENT` (kode `1` Annual · `2` Semi Annual · `3` Quarterly · `4` Monthly · lainnya Single — `GenerateUpload_Act` 2 b330 `CARI37`), `SUBJECTTO`, `AnnuityInterest`, `PremiumRefundFactor` | ⏳ paket 4 (OQ-MPNL-05 untuk kode Single) |
| b27094 **`Max Production Data Receive`** · b27992 **`Birthday`** (radio `associated`) · b28173 **`Currency`** autocomplete RD `BrowseCurrencyFacIn_RD` b28204 + **`Choose Currency`** b28701 · b29288 **`Extra Mortality (%)`** · b29474 **`Max Contract (year)`** · b29658 **`Proportional Table`** | `MAXDATARECEIVE`, `BIRTHDAY`, `CURRENCY` + `CURRENCYID` (`setCurrency_DT`), `EXTRAMORTALITY`, `MAXCONTRACT`, `PROPORTIONALTABLE` | ⏳ paket 2 / 4 |
| b27617 `%` (`CEDINGRETENTIONPCT`) · b27805 `X + n` | vis `OTHER 1=2` — **mati** | `CEDINGRETENTIONPCT` dipertahankan / `""` | ➖ medan |
| Kunci view tanpa medan form: `CEDING`, `TREATYNUMBER`, `INWARDTREATYNM`, `CEDINGLIMITXPN`, `RNMLIMITPCT`, `LIENCLAUSE`, `MONTHS` | dimuat `SetProductNameInward` 3.1 b1046 dari view, ditulis ulang apa adanya; produk baru `""` | ⏳ paket 4 |

### 3.3 Tombol bawah form (wadah b57867)

| Tombol (bNNN · teks VERBATIM) | Aksi korpus (bNNN · blok) | Rute API | Keadaan |
| --- | --- | --- | --- |
| **`Close`** b58728 | refresh + DT `HideCreateLife` b58865: 1 b139 `DATASHOW=0`, 2 b160 `STSSAVE=""`, 3 b190 `ERRMSG=""` | kembali ke grid (klien) | ⏳ paket 3 |
| **`Save`** b58998 (vis `OTHER ProductName.IsView!='true'`) | `localAction SaveProductName_Confirm` b59237 → FlowAction b34 submit **`Save`**, b33 **`Cancel`**, section `SaveProductName_Confirm` (b519 **`Do you want to save the data?`**, b1058 area teks **`Comment`** `ProductName.Comment`), aktivitas `SaveProductName_Act` b101 (§7) | `POST …/produk` · `PUT …/produk/{id}` | ⏳ paket 3 / 4 |
| **`Edit`** b59443 (vis `OTHER ProductName.IsView=='true'`) | `localAction EditProductName_Confirm` b59639 → FlowAction b19 submit **`Edit`**, b18 **`Cancel`**, section b496 **`Do you want to Edit the data?`**, transform `SetViewEdit` b71: 1 b145 `IsView = "false"` | membuka kunci form (klien) | ⏳ paket 3 |
| **`Copy`** b59812 | refresh + DT `CopyProduct` b59949: 1 b138 `ProductName.ID=""`, 2 b161 `ProductNameInward.ID=""`, 3 b191 `STSSAVE=99`, 4 b236 `ERRMSG="Data sudah dicopy, silakan melakukan perubahan dan tekan SAVE untuk menyimpan"` | form salinan (klien) → `Save` = `POST` | ⏳ paket 9 |
| **`Generate`** b60081 | refresh + `GenerateUpload_Act` b60216: 1 b224 `·`; 2 b330 `·` 40 nilai `CARI1..40`; 3 b1283 `·` `pxConvertResultsToCSV` `FileName=SeeDetail` b1335, 39 judul b1340, 40 properti b1343 | `POST …/produk/generate` → `SeeDetail.csv` | ⏳ paket 9 — OQ-MPNL-06 (judul bergeser) |

## 4. Tujuh pemilih master — `Choose*` → section → `Choose` → `set*_DT`

Setiap section pemilih: medan **`Search`** (`SearchPolicyHolder.CARI1`), Enter → `SearchPolicyHolder_act` 1 b234 `·` `CARI1 = @toUpperCase(CARI1)`; grid RD berparam `CARI1`; tombol baris **`Choose`** → DT/activity + `closeContainer`. FlowAction submit **`Submit`** / **`Cancel`** (bawaan, `pyCustomizeFALabels=false`).

| Pemilih | FlowAction → Section | RD (bNNN) · saringan · urut | Kolom grid (VERBATIM) | `Choose` → set | Rute API | Keadaan |
| --- | --- | --- | --- | --- | --- | --- |
| Ceding | `ChooseCeding` b166 → `Ceding_Section` | `BrowseCedingCoLife_RD` b2533 (kelas `AGENT`): `B AND A AND C` — `.ID Contains "L0"` b565, `.ClientName Contains Param.CedingCoLeader` b584 (= `CARI1` b1372), `.StatusActive = 1` b601; urut `.ClientName ASC` b745 | `ID` b1506 · `Name` b1647 → `.ID`, `.ClientName` | `Choose` b2226 → `setCeding_DT` b2388: `CEDING ← clientname`, `CEDINGID ← clientid` | `GET …/master/ceding?cari=` | ⏳ paket 2 |
| SOB | `ChooseSOB` b169 → `SOB_Section` | `BrowseCedingCoLife_RD` b2547 (sama) | `ID` b1521 · `Name` b1662 | `Choose` b2241 → `setSOB_DT` b2403: `SOBNAME`, `SOBID` | `GET …/master/sob?cari=` | ⏳ paket 2 |
| Policy Holder | `ChoosePolicyHolder` b91 → `PolicyHolder_Section` | `BrowseClientNusaRe_RD` b2561 (kelas `CLIENT`): `.Name Contains param.PolicyHolderName` b549, `.Name != "-"` b565, `.Name IS NOT NULL` b583, `.BU_Note = Param.Business` b599 (param tidak dikirim → abaikan); urut `.Name ASC` b745, `.BU_Note ASC` b758 | `ID` b1536 · `Name` b1677 | `Choose` b2256 → `setPolicyHolder_DT` b2416: `ProductNameInward.POLICYHODER ← id`, `POLICYHODERNAME ← name` | `GET …/master/pemegang-polis?cari=` | ⏳ paket 2 |
| Currency | `ChooseCurrency` b91 → `Currency_Section` | `BrowseCurrencyLIFE_RD` b2569 (kelas `CURRENCY`): `.Currency != "ITL"` b534, `.Currency Contains Param.Currency` b549; urut `.Currency ASC` b755 | `ID` b1544 · `Name` b1685 → `.ID`, `.Currency` | `Choose` b2264 → `setCurrency_DT` b2424: `ProductNameInward.CURRENCYID ← id`, `CURRENCY ← currency` | `GET …/master/mata-uang?cari=` | ⏳ paket 2 |
| R/I Risk | `ChooseRIRisk` b86 → `RIRISK_Section` | `BrowseRIRiskSummary` b2555 (kelas `RIRISK_LIFE_SUMMARY`): `.ID = param.id` b520, `.USEDBY = param.usedby` b533 (tidak dikirim), `.USEDBY Contains Param.SearchUsedby` b550 (= `CARI1`); urut `.ID ASC` b680 | `ID` b1528 · `Name` b1669 → `.ID`, `.USEDBY` | `Choose` b2248 → `setRIRISK_DT` b2410: `RIRISK ← usedby`, `RIRISKID ← id` | `GET …/master/ri-risk?cari=` | ⏳ paket 2 (OQ-MPNL-04) |
| R/I Rate | `ChooseRIRate` (kelas `ASM-FW-GISFW-Data-Plan`) b96 → `RIRate_Section` | `BrowseRateLifeSummary` b2707 (kelas `RATE_LIFE_SUMMARY`): `.ID = param.id`, `.USEDBY Contains param.idusedby` (= `CARI1`); urut `.ID ASC` b690 | `ID` b1594 · `RIRate Name` b1743 → `.ID`, `.USEDBY` | `Choose` b2383 → `SetRIRate` b2556: 1 b247 `·` baris plan `.RIRATEID ← id`, `.RIRATE ← usedby` | `GET …/master/ri-rate?cari=` | ⏸️ paket 2 — **OQ-MPNL-03** (view atas JSON rate; 503 berkalimat, preseden OQ-MCRL-13) |
| Cause Of Loss | `ChooseCauseOfLoss` b91 → `CauseOfLoss_Section` | `BrowseCauseofLossLife_RD` b2489 (kelas `CAUSEOFLOSS_LIFE`): `.CauseofLoss Contains Param.CauseOfLoss` b496; urut `.ID ASC` b596 | `ID` b1463 · `Name` b1604 → `.ID`, `.CauseofLoss` | `Choose` b2183 → `setCauseOfLoss_DT` b2345: `CAUSEID ← id`, `CAUSE ← cause` | `GET …/master/penyebab?cari=` | ⏳ paket 2 |

Tabel fisik = nama kelas `ASM-FW-GISFW-Int-<X>` (konvensi; preseden Retro Life `AGENT`, `BUSINESS`, `REINSURANCETYPE`) — objek yang tidak ada dijawab 503 yang menyebut objeknya (OQ-MPNL-04).

## 5. Daftar bersarang di `ProductName` (wadah b30995 / b2573)

| Grid (bNNN · judul VERBATIM) | Kolom korpus → kunci baris | Tombol / aksi (bNNN · teks) | Keadaan |
| --- | --- | --- | --- |
| b12204 **`LIEN CLAUSE (Potongan Manfaat Klaim)`** — `ProductName.LienClause` | b12745 `Usia saat Klaim` → `.Usia` b13057; b12894 `% Manfaat yang dibayarkan` → `.Manfaat` b13293 | ikon grid bawaan `pzPegaDefaultGridIcons` b12339 (vis `ProductName.IsView!='true'`), Enter = `addRow` b13616 — teks ikon **tidak ada di korpus** | ⏳ paket 7 — ⭐ **daftar** (RALAT R10) |
| b14604 **`DOCUMENT CLAIM`** — `ProductName.DocumentClaim` | b15129 `Document List` → `.Document` b15289 (dropdown `associated`, daftar tak ikut ekspor) | ikon grid bawaan b14726 (vis `IsView!='true'`), Enter = `addRow` b15591 | ⏳ paket 7 (OQ-MPNL-05) |
| b31560 **`PLAN LIST`** — `ProductName.PlanList` | b31849 `Plan Name` → `.Plan` b33124 (autocomplete RD `BrowseProductTypeLife_RD` b33198: `.CoverName` → `.Plan`, `.ID` → `.PlanID`, `.Business` → `.Name`, `.Benefit` → `.Benefit`; onChange `ProteksiPlanListLife` b33163); b31998 `Bussines` → `.Name` b33573; b32147 `Benefit` → `.Benefit` b33727; b32300 `R/I Rate` → `.RIRATE` b33890 | **`Add`** b32778 (`addRow` b32944); **`View Rate`** b34067 (vis `OTHER .RIRATE!=''`: `SetParamRate` b34299 + `localAction ViewRate` b34332); **`Choose R/I Rate`** b34548 (vis `OTHER IsView!='true'`, `localAction ChooseRIRate` b34784); **`Delete`** b35026 (`deleteRow` b35187 + `ProteksiPlanListLife` b35202) | ⏳ paket 6 (`Choose R/I Rate`, `View Rate` ⏸️ OQ-MPNL-03) |
| `ProteksiPlanListLife`: 1 b256 `·` Page-Clear-Messages; 2 b345 `·` ulang `PlanList` (`errmsg = "Plan tidak boleh sama"` b421, `errmsg1 = "Plan tidak boleh kosong"` b442, `errmsg2 = "RI/RATE tidak boleh kosong"` b463); 2.1 b522 `·` ulang; 2.1.1 b563 `·` PRE=true `.Plan==""` T=2 F=3; 2.1.2 b713 `·` PRE=true `local.plan==.Plan && local.idx!=.pxListSubscript` T=2 F=3; 2.1.3 b874 `·` PRE=true `.RIRATE==""` T=2 F=3 | — | pesan VERBATIM menolak simpan | ⏳ paket 6 |
| b37151 **`FINANCIAL UNDERWRITING`** — `ProductName.FinancialUnderwritingList` | b37674 `Min Insured` → `.MinInsured` b38721; b37822 `Max Insured` → `.MaxInsured` b38937; b37970 `Employee` → `.Employee` b39158; b38119 `Non-Employee` → `.Non_Employee` b39397 | **`Add`** b38377; ikon salin b39632 → `CopyFinancialWriting` 1 b223 `·` tambah baris salinan; **`Delete`** b39975 (+ `ProteksiPlanListLife`) | ⏳ paket 7 |
| b42078 **`UNDERWRITING LIMIT`** — `ProductName.UnderwritingLimitList` | b42598 `Min Insured` · b42746 `Max Insured` · b42894 `Min Age` · b43042 `Max Age` · b43191 `Medical` · b43340 `Description` → `.MinInsured` b43942, `.MaxInsured` b44158, `.MinAge` b44376, `.MaxAge` b44595, `.Medical` b44817 (teks bebas), `.Description` b45055 | **`Add`** b43598; ikon salin b45290 → `CopyUnderWritingLimit` 1 b224 `·`; **`Delete`** b45633 (+ `ProteksiPlanListLife`) | ⏳ paket 6 |
| Komentar — `ProductName.CommentList` (wadah b61044) | b62099 `Date` → `.Date` b62564; b62250 `PIC` → `.OperatorName` b62758; b62403 `Comment` → `.Suggest` b62940 | baca-saja (aksi baris hanya `setFocus`) | ⏳ paket 7 |
| `AddCommentList_Act` 1 b233 `·` (dipanggil **setiap** simpan, `SaveProductName_Act` 7 b1513, tanpa prakondisi): `CommentList(<APPEND>).Date = @CurrentDateTime()`, `.OperatorName = OperatorID.pxInsName`, `.IsApproved = param.status` (tak dikirim), `.Suggest = param.comment ← ProductName.Comment` | `CommentList[*].{Date, OperatorName, IsApproved, Suggest}` | — | ⏳ paket 7 |
| b48120 / b52735 **`OUTWARD`** — `ProductName.OutwardList` (wadah b47876 / b52491 `1==2`, `ExpressionCondition`) | `Reins Type` · `Transaction Year` · `Underwriting Year` → autocomplete RD `BrowseTreatyContract_Life_RD` | `Add` b49073 / b53688, **`View Reinstype`** b50721 / b55336 (`SetParamOutward` — tidak ada di korpus — + `RetrocessionModalDialog`), `Delete` b51261 / b55876 | ➖ grid **mati** (`1==2`); isinya tetap ditulis `GetReinsTypeOR_Life` (§3.1) |

## 6. Lampiran — `InboxProductName` wadah b64133, grid `TempData.AttachmentList`

Pesan statik VERBATIM: b66074 `Make sure the file name doesn't contain forbidden character such as , / | ' "`, b66233 `Recommended safe substitute should be . or _`. Kolom b68430 `File Name`.

| Unsur (bNNN · teks) | Rule dipanggil (bNNN · blok) | RDB → tabel | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| **`Add attachment`** b64698 (vis `ALWAYS TreatyIn.ViewState !='1'` → syarat diabaikan) | `SetCategory_act` b64979 (1 b229 `·` call `SetCategoryAttachTreatyin` — kategori treaty-in, ditimpa `"File"`); `localAction ProductNameAttachContent` b65012 → FlowAction: submit **`Attach`** b24, **`Cancel`** b22, pra `TreatyInitAttach` b31, pasca `ProductNameSaveAttachment` b179; refresh `LoadAttachmentProdName` b65086 | `ProductNameSaveAttachment`: 1 b290 `·` `Err = "Tidak ada file yg diattach"`; 2 b424 `·` ulang berkas; 2.1 b473 `·` `.pyCategory = "File"`; 2.2 b624 **`//`**; 2.3 b771 `·`; 2.4 b902 `·` `InsertGoogleStorage_Act` (`Durasi="1800"`, `Folder="Contract"`, gagal → keluar b935); 2.5 b1073 `·` `CARI50 = ""`; 2.6 b1199 `·` PRE=true `CARI51==""` T=3 F=2 RDB `InsertAttachProdName_Sql`; 3–8 **`//`** | `InsertAttachProdName_Sql` b84 `INSERT INTO M_ATTACHMENTPRODUCTNAME (ID = TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF3'), TREATYID = ProductName.ID, CATEGORY, FILENAME, FILEMIMETYPE, DATA_JSON = "", USERNAME, T_STORAGE_ID = ImageID)`; `Insert_T_Storage_SQL` b85 `T_STORAGE_IMAGE (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE='standard')`; `GenerateImageID_SQL` b79 MD5 `'ASMPP'‖stempel‖SYS_GUID()`; `GetAppName_SQL` b58 `T_FOLDER_IMAGE.APPNAME` | `POST …/produk/{id}/lampiran` | ✅ paket 8 — unggah = stub outbox (P5, OQ-MPNL-10) |
| **`Refresh`** b65223 (vis `ALWAYS !pyIsMobile`) | `LoadAttachmentProdName` b65389: 1 b251 `·`; 2 b370 `·` RDB `GetAttachmentProdName_Sql`; 3 b541 `·`; 4 b670 `·`; 5 b771 **`//`** | `GetAttachmentProdName_Sql` b84 `select id, FILENAME, CATEGORY, FILEMIMETYPE, T_STORAGE_ID as "type" from M_ATTACHMENTPRODUCTNAME where treatyid = ProductName.ID` | `GET …/produk/{id}/lampiran` | ✅ paket 8 |
| `Download` b67338 (vis `OTHER FALSE`) | `DownloadAll_Act` | — | — | ➖ **mati** |
| **`Download All`** b67619 | DT `TreatyInIDSetPyPortal` b67943 + `showHarness InputTreatyInOffer` (kelas `Data-Portal`, harness tidak ada di korpus) pra `TreatyInDownloadAll` (1 b223 `·` `SetTreatyIn_Act`; 3 b429 `·` `DownloadAll_Act` → `GetAllAttachment2_Sql` b84 `M_ATTACHMENTTREATY_2 where treatyid = TreatyIn.ID`) | **jalur treaty-in** (salah ekspor, OQ-056) | `GET …/produk/{id}/lampiran/unduh-semua` (zip lampiran produk) | ✅ paket 8 — perilaku dari AC tiket 08 (RALAT R15, OQ-MPNL-07) |
| Tautan nama berkas b68857 (`.pyFileName`) | `DownloadAttProdName_Act` b68926 (`ImageID ← .type`): 1–5 **`//`**; 6 b951 `·` `GetUrlGoogleStorage_Act` (`Durasi="1800"`); 7 b1078 `·` PRE=true `Param.ViewOffice` T=2 F=3 | `GetLinkStorage_SQL`, `Update_T_Storage_SQL` b85 | `GET …/lampiran/{id}/unduh` | ✅ paket 8 — berkas dari penyimpanan stub |
| Tautan **`View Office Online`** b69247 (vis `OTHER .pyFileMimeType = xls/xlsx/doc/docx/ppt/pptx`) | `DownloadAttProdName_Act` b69314 dengan `ViewOffice`: langkah 7 b1078 membungkus URL bertanda tangan ke penampil kantor **di luar** | — | `GET …/lampiran/{id}/office` | ⏸️ paket 8 — **stub**, alamat luar tidak dipanggil dan tidak ditulis (OQ-MPNL-11) |
| **`Delete`** b69663 (vis `ALWAYS TreatyIn.ViewState !='1'`) | `DeleteAttacProdName_act` b69838 (`ID ← .ID`, `ImageID ← .type`): 1 b279 `·`; 2 b409 `·` `DeleteGoogleStorage_Act` (gagal → keluar b444); 3 b526 `·` RDB `DeleteAttachProdName_Sql`; 4 b698 `·` `LoadAttachmentProdName`; 5–7 **`//`** | `DeleteAttachProdName_Sql` b85 `delete M_ATTACHMENTPRODUCTNAME where treatyid = … and id = …`; `DeleteStorage_SQL` b85 `delete T_STORAGE_IMAGE where imageid = …` | `DELETE …/lampiran/{id}` | ✅ paket 8 |
| `GetMimeType` (DecisionTable, dipanggil `InsertGoogleStorage_Act` b761) | 48 baris + `otherwise` | `inti/backend/unggah.MimeDariNamaFile` (salinan tabel yang sama) | — | ✅ paket 8 |
| `ServiceGoogle` (ConnectREST, `pyBaseURLSetting = LinkService!LinkService` b161) · `LinkService` (SystemSettings) · `GetLinkService` · `GetTokenStorage_SQL` (`GET_TOKEN_STORAGE`) | — | — | — | ⏸️ paket 8 — **stub outbox**; nol alamat di berkas (ADR-0013) |

## 7. Simpan — `SaveProductName_Act` (FlowAction `SaveProductName_Confirm` b101)

| Langkah (bNNN · blok) | Isi | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| 1 b359 `·` **PRE=false** (WHEN `TYPE=="" ‖ GRUP==""` b630 tidak dievaluasi) | `IDPEGA ← ProductName.ID`; pesan `"Product Name Empty"` b431, `"Ceding Empty"` b452, `"Policy Holder Empty"` b473, `"SOB Empty"` b494; `ProductName.POLICYHODER/POLICYHODERNAME ← ProductNameInward.*`; `UPDATEOP ← OperatorID` | salinan pemegang polis ke JSON umum; `UPDATEOP` = pelaku. ⛔ **bukan** validasi tipe/grup (RALAT R7) | ⏳ paket 3 |
| 2 b668 `·` PRE=true `PRODUCTNAME==""` b805 F=3, trans `1==1`→6 b714 | `Property-Set-Messages` pesan `local.errMsg1` | 422 `Product Name Empty` | ⏳ paket 5 |
| 3 b843 `·` PRE=true `CEDING==""` b980 | `local.errMsg2` | 422 `Ceding Empty` | ⏳ paket 5 |
| 4 b1018 `·` PRE=true `ProductNameInward.POLICYHODERNAME==""` b1155 | `local.errMsg3` | 422 `Policy Holder Empty` | ⏳ paket 5 |
| 5 b1193 `·` PRE=true `SOBNAME==""` b1330 | `local.errMsg4` | 422 `SOB Empty` | ⏳ paket 5 — Pega berhenti di pesan pertama (trans →6); sistem baru melaporkan **semua** (AC tiket 05) |
| 6 b1368 `·` PRE=true `@PropertyHasValue(CREATEOP)` b1475 T=3 F=2 | `CREATEOP ← OperatorID` bila kosong | `CREATEOP` = pelaku untuk produk baru, tetap untuk ubah | ⏳ paket 3 (OQ-MPNL-13 salinan) |
| 7 b1513 `·` | call `AddCommentList_Act(comment ← ProductName.Comment)` | satu baris `CommentList` per simpan | ⏳ paket 7 |
| 8 b1623 `·` **PRE=false** (WHEN `PoductName.*` b1793 tidak dievaluasi) | `DATAPEGA ← @GetPageJSONString()` halaman `ProductName`; `CARI1..3 ← ID, RIRISKID, RIRISK` | JSON umum dirakit repository (kunci Pega) | ⏳ paket 3 |
| 9 b1831 `·` **PRE=false** | RDB `SaveProductNameLIfe` b58 → `PEGA_M_PRODUCT_LIFE` + `COMMIT` | ⛔ prosedur **tidak** dipanggil; upsert `M_PRODUCT_LIFE (ID, JSONDATA)` ditiru, ID baru `'1' ‖ LPAD(M_PRODUCT_LIFE_SEQ.NEXTVAL, 5, '0')` (`dba-procedures-and-ddl.md` §1) | ⏳ paket 3 |
| 10 b2019 `·` **PRE=false** | RDB `SaveProductNameLIfeFlat` b84 `UPDATE M_PRODUCT_LIFE SET RIRISKID, RIRISK WHERE ID` | kolom datar `RIRISKID`, `RIRISK` (+ `PRODUCTNAME`, `BEGIN_DATE`, P1) ditulis di transaksi yang sama | ⏳ paket 3 / 4 — ⭐ **hidup** (RALAT R8) |
| 11 b2207 `·` **PRE=false** | `ProductName.ID ← OutputParam.IDPEGA`; `DATASHOW = 0` | ID hasil simpan dikembalikan; form ditutup ke grid | ⏳ paket 3 |
| 12 b2373 `·` | kosongkan `IDPEGA`, `DATAPEGA` | — | — |
| 13 b2528 `·` **PRE=false** | `IDPEGA ← ProductNameInward.ID`; `ProductNameInward.PRODUCTID ← IDPEGAOUT` | `PRODUCTID` = ID produk | ⏳ paket 4 |
| 14 b2717 `·` **PRE=false** | `DATAPEGA ← @GetPageJSONString()` halaman `ProductNameInward` | JSON inward dirakit repository | ⏳ paket 4 |
| 15 b2862 `·` **PRE=false** | RDB `SaveProductNameInwardLIfe` b84 → `PEGA_M_PRODUCT_INWARD_LIFE` + `COMMIT` | ⛔ prosedur **tidak** dipanggil; upsert `M_PRODUCTINWARD_LIFE (ID, JSONDATA)` di transaksi yang **sama** dengan langkah 9 (P4); `ID` inward = `ID` produk (P6, OQ-MPNL-02) | ⏳ paket 4 |
| 16 b3050 `·` **PRE=false** | `ProductNameInward.ID ← IDPEGAOUT` | — | ⏳ paket 4 |
| `SaveInwardProductName_Act` 1 b238 `·` · 2 b400 **`//`** `error Type` · 3 b587 **`//`** `error Grup` · 4 b762 `·` · 5 b907 `·` · 6 b1095 `·` | dipanggil hanya `Save` b14930 harness / b14125 section `InwardProductName` | — | ➖ tak terjangkau (RALAT R11) |

### 7.1 Kunci yang dibaca view (dari `dba-view-produk-life.md`) — wajib ada di JSON hasil simpan

| View | Kunci JSON |
| --- | --- |
| `PRODUCT_LIFE` ← `M_PRODUCT_LIFE` | `TYPE`, `TYPE_CEDING`, `CEDING`, `CEDINGID`, `SOBNAME`, `SOBID`, `CAUSEID`, `GRUP`, `PRODUCTNAME`, `PRODUCTCODE`, `PRODUCTTYPEID`, `PRODUCTTYPE`, `RIRISKID`, `RIRISK`, `RIRATEID`, `RIRATE`, `RICOMMID`, `RICOMM`, `INWARDNAME`, `UnderwritingLimitList`, `OUTWARDNAMEID`, `OUTWARDNAME`, `OUTWARDRATEID`, `OUTWARDRATE`, `OUTWARDCOMMID`, `OUTWARDCOMM`, `BENEFITID`, `BENEFIT`, `CAUSE`, `OutwardList[0].OVR_COMM`, `POLICYHODER`, `POLICYHODERNAME`, `TREATYNUMBER`, `CREATEOP`, `UPDATEOP` |
| `DOCUMENTCLAIM_LIFE` ← `M_PRODUCT_LIFE` | `DocumentClaim[*].Document` |
| `PRODUCTINWARD_LIFE` ← `M_PRODUCTINWARD_LIFE` | `PRODUCTID`, `INSURED`, `CEDING`, `TREATYNUMBER`, `ADDENDUMWORD`, `AMANDEMENTSCHD`, `INWARDTREATYNM`, `BEGIN`, `MATURE`, `CEDINGRETENTIONNUM`, `CEDINGRETENTIONPCT`, `CEDINGLIMIT`, `CEDINGLIMITXPN`, `MINAGE`, `MAXAGE`, `BIRTHDAY`, `EXTRAPREMI`, `CURRENCY`, `RNMSHARE`, `EXTRAMORTALITY`, `RNMLIMITNUM`, `RNMLIMITPCT`, `LIENCLAUSE`, `MONTHS`, `MINSUMINSURED`, `MAXSUMINSURED`, `MAXSUMREASURED`, `MAXCONTRACT`, `PAYMENT`, `PROPORTIONALTABLE`, `SUBJECTTO`, `POLICYHODER`, `POLICYHODERNAME`, `BROKERAGE`, `ADDENDUMNO`, `AMANDEMENTNO`, `MAXDATARECEIVE`, `MAXEXPIREDCLAIM`, `STNC` |

Kunci skalar tanpa nilai ditulis `""` (Oracle membacanya NULL, sama dengan kunci absen); kunci daftar ditulis `[]`. Setiap baris
`OutwardList` membawa `OVR_COMM` (`""` — korpus tidak punya penulisnya), setiap baris `DocumentClaim` membawa `Document`.
Kunci yang sudah ada di JSON lama dan tidak dikelola layar **dipertahankan** apa adanya (Pega memuat halaman utuh lewat
`adoptJSONObject`, `SetProductName` 4 b957).

## 8. Rule yang sengaja tidak dibangun

| Rule | Bukti | Keadaan |
| --- | --- | --- |
| `SetTreatyIn_Act`, `TreatyInInputVis`, `CheckDuplicateOffer`, `GetCountClaim`, `GetCurrentDate`, `BrowseTreatyIn`, `SaveTreatyIn`, `BrowseTREATY_IN`, `TreatySetReinstatement`, `SetReinstatementPct`, `ConvertHistoryDate`, `LoadAttachment`, `GetAttachment2_Sql`, `TreatyInDownloadAll`, `DownloadAll_Act`, `GetAllAttachment2_Sql`, `TreatyInIDSetPyPortal` | jalur treaty inward (kelas `Data-Portal` / `TREATY_IN`) — **salah ekspor** (OQ-056 `[keputusan work owner]`); satu-satunya pemicu di layar produk adalah `Download All` b67619 (§6) | ➖ |
| `SaveInwardProductName_Act`, harness + section `InwardProductName` | hanya dibuka `Inward` b75322 `OTHER 1=2` | ➖ |
| `CountMaxReasured_Act` | dipanggil hanya `TYPE_CEDING` b3385 (`OTHER 1=2`) | ➖ |
| `SetParamRate` (`ParamID.RIRATEID ← InputBusinessLife.RIRATEID`, kelas `@baseclass`) | halaman `InputBusinessLife` milik Retro Life — di layar ini `ParamID` diisi baris plan (`View Rate` b34299) | ⏸️ `View Rate` OQ-MPNL-03 |
| `When/recordEvent`, `SetCategoryAttachTreatyin` + `CategoryAttach_SQL` | `recordEvent` nol pemanggil; kategori lampiran ditimpa `"File"` (`ProductNameSaveAttachment` 2.1 b473) sesudah `SetCategory_act` mengisi daftar kategori treaty-in | ➖ |

## 9. Rute API dan komponen

Prefix `/api/master-product-name-life` (`backend/handlers/rute_mpnl.go` `Prefix`). Setiap rute menuntut identitas pelaku (401 tanpa),
menjawab 503 bila Oracle tidak dikonfigurasi, galat berbadan `{"galat": "..."}`.

| Metode dan jalur | Layar / tombol Pega | Paket |
| --- | --- | --- |
| `GET /produk` | grid `InboxProductName` mode daftar (b71246, RD `BrowseProduct_Life`) | 1 |
| `GET /produk/{id}` | tombol `View` b74753 (`SetProductName` + `SetProductNameInward`) | 1 |
| `GET /master/{jenis}?cari=` — `ceding`, `sob`, `pemegang-polis`, `mata-uang`, `ri-risk`, `penyebab`; `ri-rate` = 503 (OQ-MPNL-03) | tombol `Choose*` → section pemilih → grid RD (§4); autocomplete medan form | 2 |
| `POST /produk` · `PUT /produk/{id}` | `Add` b71816 / `View` → `Edit` b59443 → `Save` b58998 → `SaveProductName_Confirm` → `SaveProductName_Act` (§7); nol rute hapus (korpus tanpa hapus produk) | 3 |
| (rute yang sama) | sisi inward `SaveProductName_Act` 13–16 — `M_PRODUCTINWARD_LIFE` di transaksi yang sama (P4); uji tiga view (`db`) | 4 |
| (rute yang sama) | wajib-isi `SaveProductName_Act` 2–5 — `Product Name Empty`, `Ceding Empty`, `Policy Holder Empty`, `SOB Empty` (422, semua sekaligus) | 5 |
| `GET /master-plan?cari=` · `GET /rate?riRateId=` (503, OQ-MPNL-03) | autocomplete `Plan Name` b33124 · tombol `View Rate` b34067; gerbang `ProteksiPlanListLife` dan `UNDERWRITING LIMIT` di simpan | 6 |
| (rute yang sama) | `AddCommentList_Act` setiap simpan; `LIEN CLAUSE`, `DOCUMENT CLAIM`, `FINANCIAL UNDERWRITING` disimpan bersama produk | 7 |
| `GET /produk/{id}/lampiran` · `POST /produk/{id}/lampiran` (multipart `berkas`) | `Refresh` b65223 / `View` · `Add attachment` b64698 → `Attach` b24 (`ProductNameSaveAttachment`) | 8 |
| `POST /produk/{id}/lampiran/{lid}/ulangi` | (tiket 09) kirim ulang efek outbox yang gagal | 8 |
| `GET /produk/{id}/lampiran/{lid}/unduh` · `GET /produk/{id}/lampiran/unduh-semua` | tautan nama berkas b68857 · `Download All` b67619 (zip, R15) | 8 |
| `GET /produk/{id}/lampiran/{lid}/office` (503 stub, OQ-MPNL-11) · `DELETE /produk/{id}/lampiran/{lid}` | `View Office Online` b69247 · `Delete` b69663 (`DeleteAttacProdName_act`) | 8 |
