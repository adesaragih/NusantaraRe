# Paritas layar dan aksi — Master Product Name Life

> Disusun 01-10-2026 (sesi implementasi modul, paket 0, brief `PROMPT-IMPLEMENTASI-MODUL-MASTER-PRODUCT-NAME-LIFE.md`).
> Satu baris per harness, section, tombol, dan aksi korpus → rule yang dipanggil → RDB/RD → tabel/kunci JSON → rute API → komponen React.
> Nomor baris `bNNN` = baris berkas korpus sesudah `sed -e 's/></>\n</g'` atas `D:\XML\RNM_BRD\Master Product Name Life\` (konvensi repo).
> ⭐ **Ralat 01-10-2026 (lanjutan 1 L9) — perintah standar:** `sed -e 's/></>\n</g' <berkas> | grep -n '<tag>teks</tag>'`. Setiap `bNNN` di
> dokumen ini menunjuk baris TAG yang memuat teks yang dikutip: tombol/tautan = `<pyLabel>`, centang = `<pyCheckboxCaption>`, medan =
> `<pyLabelFieldValue>`, judul/kolom/properti sel = `<pyValue>`, ikon = rule yang dipanggilnya, aksi sel = `<pyActivity>`/`<pyLocalAction>`/
> `<pyHarnessName>`/`<pyName>`/`<pyAction>`, langkah activity = `<pyStepsActivityName>` (perulangan: `<pyStepsObjectName>`), langkah
> DataTransform = `<pyPropertiesName>`, kolom RD = `<pyFieldName>`, saringan RD = `<pyFilterValue>`/`<pyFilterOperation>`; WHEN, SQL,
> label FlowAction = baris tagnya sendiri. *Kalimat lama dikutip:* nomor sebelumnya menunjuk baris STRUKTURAL (`<pyFormat>` sel,
> `<rowdata>` langkah/aksi) — benar menurut penomoran `sed`, tetapi bergeser 2–65 baris dari baris yang `grep -n` atas teksnya temukan
> (mis. `Inward` b75322 → `<pyLabel>Inward</pyLabel>` b75368). 871 nomor di 60 berkas modul diganti (putaran kedua sesudah `/code-review`: 26 lagi, total 897); diverifikasi dua cara: peta
> pengurai (baris struktural → tag) dan pembacaan ulang setiap nomor baru dengan perintah standar (0 salah).
> Setiap langkah activity yang dikutip mencetak `pyStepsBlockName` — `·` = kosong (langkah hidup), `//` = ter-remark (tidak pernah jalan).
> Prakondisi dibaca dari urutan aksi: `WhenTrue 2` = lanjut, `3` = lewati langkah; `T=-` = kosong = lanjut; **`PRE=false` = prakondisi
> dimatikan → langkah SELALU jalan**, apa pun isi WHEN-nya.
> Keadaan: ✅ dibangun · ⏸️ dibangun, datanya menunggu OQ · ➖ sengaja tidak dibawa (bukti di kolom Catatan) · ⏳ menunggu paket N.
> Kunci JSON ditulis **peka huruf besar-kecil** persis Pega (`POLICYHODER`, `UnderwritingLimitList`, `Non_Employee`, ...).

## 0. Alat baca dan sensus

| Hal | Isi |
| --- | --- |
| Alat | pengurai token bernomor baris `pega.py` (scratchpad sesi, tidak di repo): `akt` pohon langkah + `pyStepsBlockName`, `tombol`/`sel` sel section + aksi + rantai visibilitas wadah + sumber grid, `dt`, `rd`, `rdb`, `fa`; ditambah peta pemanggil (`>Nama<` di seluruh 114 berkas) |
| Uji alat | nomor baris alat = nomor `sed`; *kalimat lama (paket 0) dikutip:* "sel tombol pertama `InboxProductName.xml` b5181 = baris `<pyFormat>pxButton` mentah" — **ralat 01-10-2026 (L9):** nomor kini menunjuk tag ber-teks, tombol pertama `Choose Ceding Name` = `<pyLabel>` b5222; `ProteksiPlanListLife` = 2.1.1–2.1.3 bersarang (cocok XML) |
| Berkas | 114: Activity 37, RDBList 22, ReportDefinition 17, Section 12, FlowAction 11, DataTransform 10, Harness 1, DecisionTable 1, When 1, ConnectREST 1, SystemSettings 1 (+ `Struktur_InboxProductName.xlsx`, bukan rule) |
| Langkah ber-`//` | **27** — cara 1 alat (jumlah per activity), cara 2 `grep -c '<pyStepsBlockName>//'` mentah: sepakat. Sebaran: `ProductNameSaveAttachment` 7, `DownloadAttProdName_Act` 5, `CheckDuplicateOffer` 4, `DeleteAttacProdName_act` 3, `SaveInwardProductName_Act` 2, `SetProductName` 2, `LoadAttachment` 2, `LoadAttachmentProdName` 1, `SetCategoryAttachTreatyin` 1 |
| Tombol `pxButton` | **42** — cara 1 `grep -c '<pyFormat>pxButton</pyFormat>'`, cara 2 sel `Embed-Display-Table-Cell` ber-`pyFormat = pxButton`: sepakat per berkas (`InboxProductName` 33, tujuh pemilih masing-masing 1, `InwardProductName` section 1, harness 1) |
| Tautan/ikon beraksi | **4** di `InboxProductName`: `pxIcon` b39697, b45355 (salin baris), `pxLink` b68903 (nama berkas), b69291 (`View Office Online`) |
| Teks tombol | `pyModes[2].pyLabel` (keterangan yang dirender `pxButton`/`pxLink`); `pyLabelFieldValue` = label SEL (`Button`/`End Period`), bukan teks tombol |
| Medan visibilitas | sel: `pyUserData.pyVisible` (`ALWAYS` → syarat diabaikan; `OTHER`/`NOTBLANK` → syarat berlaku); wadah: `pyContainerVisibleWhen` berlaku hanya bila `pyIsVisibilityOption` = `CONDITION` atau `ExpressionCondition` (62 `ALWAYS`, 14 `CONDITION`, 2 `ExpressionCondition` di korpus) |

## 1. Halaman awal — pintu masuk modul

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Section/InboxProductName.xml` mode daftar: wadah b71246 `OutputParam.DATASHOW != 1` (`CONDITION`) berisi grid produk; mode form: wadah b2573 / b30995 / b57867 / b61044 / b64133 `OutputParam.DATASHOW = 1` | tombol menu **Master Product Name Life** → halaman `mpnl-produk` (grid produk; form dibuka dari dalam) | ✅ paket 10 (menu slot 960 + layar `mpnl-produk`) |
| ⭐ **Bukti `InboxProductName` = pintu masuk**: tidak dibuka rule mana pun di korpus (peta pemanggil: hanya dirujuk sebagai sasaran `refresh otherSection` oleh enam section pemilih — `Ceding_Section.xml` b2285, `SOB_Section.xml` b2300, `PolicyHolder_Section.xml` b2318, `Currency_Section.xml` b2323, `RIRISK_Section.xml` b2307, `CauseOfLoss_Section.xml` b2242), jadi dibuka portal/navigasi di luar ekspor. Satu-satunya harness, `Harness/InwardProductName.xml`, dibuka **hanya** oleh tombol `Inward` b75368 yang bervisibilitas `OTHER 1=2` (mati) | halaman awal = `InboxProductName`; harness `InwardProductName` **tidak dibangun** | ralat spec §1 *"Titik masuk `InwardProductName` … satu-satunya Harness"* (RALAT R11) |

## 2. Grid daftar produk — `InboxProductName` mode daftar (wadah b71246)

Sumber: RD `BrowseProduct_Life` (b76658, kelas `ASM-FW-GISFW-Int-PRODUCT_LIFE` = view `PRODUCT_LIFE`), tanpa saringan (b682), urut `.ID ASC` (b1096), 10 baris/halaman (`pyRDLPageSize` b71371), `pyGridPaginator` b72103.
Sistem baru membaca `M_PRODUCT_LIFE.ID` + `JSONDATA` sendiri (P2) — kunci kolom grid = kunci JSON yang dibaca view.

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → tabel/kunci | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kolom `ID` b72403 · `Ceding` b72551 · `Treaty Number` b72707 · `Treaty Name` b72866 · `Create Operator` b73025 · `Last Updated Operator` b73184 → `.ID` b73625 · `.CEDING` b73810 · `.TREATYNUMBER` b73987 · `.INWARDNAME` b74179 · `.CREATEOP` b74371 · `.UPDATEOP` b74563 | — | `M_PRODUCT_LIFE.ID`; JSON `CEDING`, `TREATYNUMBER`, `INWARDNAME`, `CREATEOP`, `UPDATEOP` | `GET /api/master-product-name-life/produk` | ✅ paket 1 (baca) / 2 (layar) |
| ⚠️ Wadah grid b71574 `pyContainerVisibleWhen = IsFire` | — | When `IsFire` tidak ada di korpus | — | ✅ ikut XML: `pyIsVisibilityOption = ALWAYS` → kondisi diabaikan, grid selalu tampil |
| Tombol b71865: label sel **`End Period`**, teks **`Add`**, tooltip **`Add New Data`**, vis `ALWAYS !IsSpreadingUW` (syarat diabaikan) | `NewProductLife` b71889/b71992: 1 b234 `·` Page-New `ProductName`; 2 b378 `·` `DATASHOW=1`, `ERRMSG=""`, `HASILD7=0`, `STSSAVE=0`; 3 b579 `·` Page-Remove `ProductName`, `ProductNameDisplay`, `ProductNameInward` | — | form baru (klien) | ✅ paket 3 |
| Tombol baris **`View`** b74798 | `SetProductName` b74821 (param `ID ← .ID`): 1 b470 `·` hapus halaman; 2 b585 `·` `DATASHOW=1`, `ProductName.ID ← Param.ID`; 3 b782 `·` RDB-List `BrowseUnderwritingList`; 4 b959 `·` Java `adoptJSONObject(BENEFIT)`; 5 b1067 **`//`**; 6 b1206 **`//`** (salin param lama); 7 b1718 `·` call `SetProductNameInward(ProductID ← ProductName.ID)`; 8 b2147 `·` `IsView = true` — lalu refresh `LoadAttachmentProdName` b74973 | `BrowseUnderwritingList` b61 `select JSONDATA as BENEFIT from m_product_life where ID={ProductName.ID}`; `SetProductNameInward` 1 b652 `·`, 2 b828 `·` pxShowReport RD `BrowseProductInward` (saringan b840 `.PRODUCTID = Param.ProductID`), 3 b1006 `·` ulang, 3.1 b1048 `·` salin **40** medan view | `GET …/produk/{id}` | ✅ paket 1 / 4 |
| Tombol baris **`Inward`** b75368 (vis `OTHER 1=2`) | `showHarness` `InwardProductName` popup + `SetProductNameInward` | — | — | ➖ **mati** (`1=2`); harness dan `SaveInwardProductName_Act` tak terjangkau (RALAT R11) |

## 3. Form produk — `InboxProductName` mode form (wadah b2573, `DATASHOW = 1`)

Pesan: wadah b616 `STSSAVE==100` → `OutputParam.ERRMSG` b878; wadah b1269 `STSSAVE==99` → `ERRMSG` b1531 (dipakai `CopyProduct`); wadah b1921 `ProductName.ERRMSG!=''` → b2183.
Baca-saja: medan ber-`ro = ProductName.IsView=='true'` (mode lihat sesudah `View`; `Edit` membukanya). *(Ralat audit 02-10-2026: tidak SETIAP medan - enam medan pemilih master SELALU baca-saja, `pyReadOnly` true + `pyEditOptions` Read-only + `pyReadOnlyCondition` kosong: Ceding b4040, SOB b4428, R/I Risk Name b7362, Cause Of Loss b10626, Policy Holder b17062, Currency b28105; nilainya hanya dari tombol `Choose*`. Sel `PLAN LIST` `Bussines` `.Name` b33504 dan `Benefit` b33658 juga selalu baca-saja. Lihat RALAT 02-10-2026.)*

### 3.1 Sisi umum — halaman `ProductName` → `M_PRODUCT_LIFE.JSONDATA`

| Medan korpus (bNNN · label VERBATIM) | Kontrol · visibilitas · aksi | Kunci JSON | Keadaan |
| --- | --- | --- | --- |
| b3140 `ProductName.TYPE` · b3346 `TYPE_CEDING` · b6076 `GRUP` (`Group`) | radio, vis `OTHER 1=2` — **mati**; `TYPE_CEDING` onChange `CountMaxReasured_Act` b3385 ikut mati | `TYPE`, `TYPE_CEDING`, `GRUP` (dibaca view `PRODUCT_LIFE`) — ditulis apa adanya dari JSON lama, `""` untuk produk baru | ➖ medan; kunci ✅ paket 3 |
| b3620 **`Product Name`** `ProductName.PRODUCTNAME` | dropdown `associated` (daftar lokal Property, **tidak ikut ekspor**); onChange `SetTreatyName_Act` b3686: 1 b239 `·`, 2 b436 `·` `INWARDNAME ← PRODUCTNAME + " " + POLICYHODERNAME` | `PRODUCTNAME` (+ kolom datar `PRODUCTNAME`) | ✅ paket 3 — isian teks (OQ-MPNL-05) |
| b3894 **`Product Code`** `ProductName.ID` | teks, vis `ALWAYS` | `ID` | ✅ paket 3 — **baca-saja** (ADR-0006; RALAT R14) |
| b4075 **`Ceding`** `ProductName.CEDING` + tombol **`Choose Ceding Name`** b5222 | autocomplete RD `BrowseCedingCoLife_RD` b4203; tombol `localAction ChooseCeding` b5374 (wadah b4945 `IsView!='true'`) | `CEDING`, `CEDINGID` | ✅ paket 2 |
| b4463 **`SOB`** `ProductName.SOBNAME` + **`Choose SOB`** b5563 | autocomplete RD `BrowseCedingCoLife_RD` b4586; `localAction ChooseSOB` b5715 | `SOBNAME`, `SOBID` | ✅ paket 2 |
| b6282 `Plan` (`PRODUCTTYPE`) · b6690 `R/I Rate` (`RIRATE`) · b9081 `Outward Name` · b9498 `Outward Rate` · b9895 `Outward Comm` · b10280 `Benefit` | autocomplete, vis `OTHER 1=2` — **mati** | `PRODUCTTYPE`, `RIRATE`, `OUTWARDNAME`, `OUTWARDRATE`, `OUTWARDCOMM`, `BENEFIT` (+ `*ID`) dibaca view — dipertahankan dari JSON lama / `""` | ➖ medan; kunci ✅ paket 3 |
| b7104 **`Deduction (%)`** `ProductName.RICOMM` | angka | `RICOMM` | ✅ paket 3 |
| b7398 **`R/I Risk Name`** `ProductName.RIRISK` + **`Choose R/I Risk`** b8205 | autocomplete RD `BrowseRIRiskSummary` b7520; `localAction ChooseRIRisk` b8357 | `RIRISK`, `RIRISKID` (+ kolom datar `RIRISK`, `RIRISKID`) | ✅ paket 2 |
| b8719 **`Treaty Name`** `ProductName.INWARDNAME` · b8900 **`Treaty Number`** `ProductName.TREATYNUMBER` | teks | `INWARDNAME`, `TREATYNUMBER` | ✅ paket 3 |
| b10661 **`Cause Of Loss`** `ProductName.CAUSE` + **`Choose Cause of Loss`** b11360 | dropdown RD `BrowseCauseofLossLife_RD` b10781 (nilai/tampil `.CauseofLoss`); `localAction ChooseCauseOfLoss` b11509 | `CAUSE`, `CAUSEID` | ✅ paket 2 |
| b47312 checkbox **`On Retention`** (`pyCheckboxCaption` b47312) `ProductName.IsORS`, vis `ALWAYS` (di luar wadah `1==2`) | onChange `GetReinsTypeOR_Life` b47488: 1 b236 `·` **PRE=false** hapus `OutwardList`; 2 b383 `·` `Temp.CARI1/2 ← BEGIN/MATURE`; 3 b534 `·` **PRE=false** RDB-List `BrowseReinstypeOR_SQL`; 4 b731 `·` **PRE=false** ulang; 4.1 b770 `·` tambah baris `OutwardList` | `BrowseReinstypeOR_SQL` b84: `TREATYCONTRACT_LIFE tc JOIN TREATYYEAR_LIFE ty ON ty.id = tc.idtreatyyear WHERE TO_DATE(BEGIN,'DD/MM/YYYY') >= tc.treatystartdate AND TO_DATE(MATURE,'DD/MM/YYYY') <= tc.treatyenddate AND REINSTYPEID = '10200'` → `IsORS`, `OutwardList[*].{REINSTYPEID, REINSTYPENAME, TRANSACTIONYEAR←TREATYYEAR, TREATYCONTRACTID, UNDERWRITINGYEAR}` | ✅ paket 9 — ⭐ **hidup** (RALAT R9, P3); `hitungOutward` di badan simpan; `TREATYCONTRACTID` `""` (OQ-MPNL-15) |

### 3.2 Sisi inward — halaman `ProductNameInward` → `M_PRODUCTINWARD_LIFE.JSONDATA`

Diisi di form yang **sama** (bukan harness `InwardProductName`). Disimpan `SaveProductName_Act` langkah 13–16 (§7).

| Medan korpus (bNNN · label VERBATIM) | Kunci JSON | Keadaan |
| --- | --- | --- |
| b17097 **`Policy Holder`** autocomplete RD `BrowseClientNusaRe_RD` b17224 + **`Choose Policy Holder`** b17827 (`localAction ChoosePolicyHolder` b17979); onChange `SetTreatyName_Act` b17168 | `POLICYHODER` (id), `POLICYHODERNAME` (disalin juga ke JSON umum, `SaveProductName_Act` 1 b361) | ✅ paket 2 / 4 |
| b18341 **`Insured`** · b18834 **`Addendum No.`** · b19432 **`Addendum`** · b20155 **`Amandement No.`** · b20760 **`Amandement`** · b21781 **`Max Notification Claim Expired`** | `INSURED`, `ADDENDUMNO`, `ADDENDUMWORD`, `AMANDEMENTNO`, `AMANDEMENTSCHD`, `MAXEXPIREDCLAIM` | ✅ paket 4 |
| b21969 **`Begin Date`** · b22304 **`STNC`** · b27250 **`Expired Date`** | `BEGIN`, `STNC`, `MATURE` — teks `dd/MM/yyyy` (bentuk `TO_DATE(…,'DD/MM/YYYY')` `BrowseReinstypeOR_SQL` b84, contoh `[data DBA]`) | ✅ paket 4 |
| b22494 **`Ceding Retention (%)`** · b22657 **`Ceding's Limit`** (onChange `CountMaxSumReasured_Act` b22725) · b22915 **`Brokerage Fee (%)`** · b23078 **`Minimum Age (Years)`** · b23265 **`Maximum Age (Years)`** · b23452 **`Expiry Age (Years)`** · b23635 **`Extra Premium`** · b23796 **`Min Sum Insured`** · b23956 **`Max Sum Insured`** (onChange `CountMaxSumReasured_Act` b24024) · b24214 **`Max Sum Reasured`** · b24674 **`Nusantara Re Share (%)`** + label b24880 **`OF SUM REASURED`** · b25236 **`Nusantara Re's Limit`** | `CEDINGRETENTIONNUM`, `CEDINGLIMIT`, `BROKERAGE`, `MINAGE`, `MAXAGE`, `EXPIRYAGE`, `EXTRAPREMI`, `MINSUMINSURED`, `MAXSUMINSURED`, `MAXSUMREASURED`, `RNMSHARE`, `RNMLIMITNUM` | ✅ paket 4 |
| `CountMaxSumReasured_Act`: 1 b239 `·` `local ← MAXSUMINSURED, CEDINGLIMIT`; 2 b436 `·` `MAXSUMREASURED = MaxSumInsured − CedingLimit` | `MAXSUMREASURED` dihitung, tetap dapat ditimpa (medan b24214 dapat disunting) | ✅ paket 4 (hitung desimal persis di klien) |
| b25398 **`Premium Factor (%)`** (vis `OTHER PAYMENT==3`) · b25611 **`Premium Payment Method`** (dropdown `associated`) · b25924 **`Subject To`** (area teks) · b26092 **`Annuity Interest (%)`** · b26306 **`Premium Refund Factor (%)`** (keduanya `ALWAYS`, syarat `PAYMENT==3` diabaikan) | `PREMIUMFACTOR`, `PAYMENT` (kode `1` Annual · `2` Semi Annual · `3` Quarterly · `4` Monthly · lainnya Single — `GenerateUpload_Act` 2 b332 `CARI37`), `SUBJECTTO`, `AnnuityInterest`, `PremiumRefundFactor` | ✅ paket 4 (OQ-MPNL-05 untuk kode Single) |
| b27062 **`Max Production Data Receive`** · b27960 **`Birthday`** (radio `associated`) · b28140 **`Currency`** autocomplete RD `BrowseCurrencyFacIn_RD` b28204 + **`Choose Currency`** b28742 · b29256 **`Extra Mortality (%)`** · b29442 **`Max Contract (year)`** · b29626 **`Proportional Table`** | `MAXDATARECEIVE`, `BIRTHDAY`, `CURRENCY` + `CURRENCYID` (`setCurrency_DT`), `EXTRAMORTALITY`, `MAXCONTRACT`, `PROPORTIONALTABLE` | ✅ paket 2 / 4 |
| b27585 `%` (`CEDINGRETENTIONPCT`) · b27773 `X + n` | vis `OTHER 1=2` — **mati** | `CEDINGRETENTIONPCT` dipertahankan / `""` | ➖ medan |
| Kunci view tanpa medan form: `CEDING`, `TREATYNUMBER`, `INWARDTREATYNM`, `CEDINGLIMITXPN`, `RNMLIMITPCT`, `LIENCLAUSE`, `MONTHS` | dimuat `SetProductNameInward` 3.1 b1048 dari view, ditulis ulang apa adanya; produk baru `""` | ✅ paket 4 |

### 3.3 Tombol bawah form (wadah b57867)

| Tombol (bNNN · teks VERBATIM) | Aksi korpus (bNNN · blok) | Rute API | Keadaan |
| --- | --- | --- | --- |
| **`Close`** b58770 | refresh + DT `HideCreateLife` b58883: 1 b143 `DATASHOW=0`, 2 b172 `STSSAVE=""`, 3 b202 `ERRMSG=""` | kembali ke grid (klien) | ✅ paket 3 |
| **`Save`** b59041 (vis `OTHER ProductName.IsView!='true'`) | `localAction SaveProductName_Confirm` b59259 → FlowAction b34 submit **`Save`**, b33 **`Cancel`**, section `SaveProductName_Confirm` (b519 **`Do you want to save the data?`**, b1025 area teks **`Comment`** `ProductName.Comment`), aktivitas `SaveProductName_Act` b101 (§7) | `POST …/produk` · `PUT …/produk/{id}` | ✅ paket 3 / 4 |
| **`Edit`** b59489 (vis `OTHER ProductName.IsView=='true'`) | `localAction EditProductName_Confirm` b59663 → FlowAction b19 submit **`Edit`**, b18 **`Cancel`**, section b496 **`Do you want to Edit the data?`**, transform `SetViewEdit` b71: 1 b151 `IsView = "false"` | membuka kunci form (klien) | ✅ paket 3 |
| **`Copy`** b59854 | refresh + DT `CopyProduct` b59965: 1 b144 `ProductName.ID=""`, 2 b173 `ProductNameInward.ID=""`, 3 b203 `STSSAVE=99`, 4 b236 `ERRMSG="Data sudah dicopy, silakan melakukan perubahan dan tekan SAVE untuk menyimpan"` | form salinan (klien) → `Save` = `POST /produk` + `salinanDari` | ✅ paket 9 (backend; layar paket 10) — OQ-MPNL-13 |
| **`Generate`** b60122 | refresh + `GenerateUpload_Act` b60228: 1 b226 `·`; 2 b332 `·` 40 nilai `CARI1..40`; 3 b1285 `·` `pxConvertResultsToCSV` `FileName=SeeDetail` b1335, 39 judul b1340, 40 properti b1343 | `POST …/produk/generate` → `SeeDetail.csv` | ✅ paket 9 — OQ-MPNL-06 (judul bergeser) |

## 4. Tujuh pemilih master — `Choose*` → section → `Choose` → `set*_DT`

Setiap section pemilih: medan **`Search`** (`SearchPolicyHolder.CARI1`), Enter → `SearchPolicyHolder_act` 1 b236 `·` `CARI1 = @toUpperCase(CARI1)`; grid RD berparam `CARI1`; tombol baris **`Choose`** → DT/activity + `closeContainer`. FlowAction submit **`Submit`** / **`Cancel`** (bawaan, `pyCustomizeFALabels=false`).

| Pemilih | FlowAction → Section | RD (bNNN) · saringan · urut | Kolom grid (VERBATIM) | `Choose` → set | Rute API | Keadaan |
| --- | --- | --- | --- | --- | --- | --- |
| Ceding | `ChooseCeding` b166 → `Ceding_Section` | `BrowseCedingCoLife_RD` b2533 (kelas `AGENT`): `B AND A AND C` — `.ID Contains "L0"` b570, `.ClientName Contains Param.CedingCoLeader` b587 (= `CARI1` b1372), `.StatusActive = 1` b607; urut `.ClientName ASC` b747 | `ID` b1502 · `Name` b1643 → `.ID`, `.ClientName` | `Choose` b2265 → `setCeding_DT` b2388: `CEDING ← clientname`, `CEDINGID ← clientid` | `GET …/master/ceding?cari=` | ✅ paket 2 |
| SOB | `ChooseSOB` b169 → `SOB_Section` | `BrowseCedingCoLife_RD` b2547 (sama) | `ID` b1517 · `Name` b1658 | `Choose` b2280 → `setSOB_DT` b2403: `SOBNAME`, `SOBID` | `GET …/master/sob?cari=` | ✅ paket 2 |
| Policy Holder | `ChoosePolicyHolder` b91 → `PolicyHolder_Section` | `BrowseClientNusaRe_RD` b2561 (kelas `CLIENT`): `.Name Contains param.PolicyHolderName` b556, `.Name != "-"` b570, `.Name IS NOT NULL` b597, `.BU_Note = Param.Business` b603 (param tidak dikirim → abaikan); urut `.Name ASC` b747, `.BU_Note ASC` b760 | `ID` b1532 · `Name` b1673 | `Choose` b2298 → `setPolicyHolder_DT` b2416: `ProductNameInward.POLICYHODER ← id`, `POLICYHODERNAME ← name` | `GET …/master/pemegang-polis?cari=` | ✅ paket 2 |
| Currency | `ChooseCurrency` b91 → `Currency_Section` | `BrowseCurrencyLIFE_RD` b2569 (kelas `CURRENCY`): `.Currency != "ITL"` b541, `.Currency Contains Param.Currency` b553; urut `.Currency ASC` b757 | `ID` b1540 · `Name` b1681 → `.ID`, `.Currency` | `Choose` b2303 → `setCurrency_DT` b2424: `ProductNameInward.CURRENCYID ← id`, `CURRENCY ← currency` | `GET …/master/mata-uang?cari=` | ✅ paket 2 |
| R/I Risk | `ChooseRIRisk` b86 → `RIRISK_Section` | `BrowseRIRiskSummary` b2555 (kelas `RIRISK_LIFE_SUMMARY`): `.ID = param.id` b524, `.USEDBY = param.usedby` b537 (tidak dikirim), `.USEDBY Contains Param.SearchUsedby` b555 (= `CARI1`); urut `.ID ASC` b682 | `ID` b1524 · `Name` b1665 → `.ID`, `.USEDBY` | `Choose` b2287 → `setRIRISK_DT` b2410: `RIRISK ← usedby`, `RIRISKID ← id` | `GET …/master/ri-risk?cari=` | ✅ paket 2 (OQ-MPNL-04) |
| R/I Rate | `ChooseRIRate` (kelas `ASM-FW-GISFW-Data-Plan`) b96 → `RIRate_Section` | `BrowseRateLifeSummary` b2707 (kelas `RATE_LIFE_SUMMARY`): `.ID = param.id`, `.USEDBY Contains param.idusedby` (= `CARI1`); urut `.ID ASC` b692 | `ID` b1590 · `RIRate Name` b1739 → `.ID`, `.USEDBY` | `Choose` b2424 → `SetRIRate` b2556: 1 b249 `·` baris plan `.RIRATEID ← id`, `.RIRATE ← usedby` | `GET …/master/ri-rate?cari=` | ⏸️ paket 2 — **OQ-MPNL-03** (view atas JSON rate; 503 berkalimat, preseden OQ-MCRL-13) *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: view DEV `RATE_LIFE_SUMMARY` dibaca saja `SELECT ID, USEDBY`, `UPPER(USEDBY) LIKE`, urut `ID ASC`, 500 baris b674; `Choose` → `.RIRATEID ← .ID`, `.RIRATE ← .USEDBY` master; ✅ OQ-MPNL-03 ditutup)* |
| Cause Of Loss | `ChooseCauseOfLoss` b91 → `CauseOfLoss_Section` | `BrowseCauseofLossLife_RD` b2489 (kelas `CAUSEOFLOSS_LIFE`): `.CauseofLoss Contains Param.CauseOfLoss` b502; urut `.ID ASC` b598 | `ID` b1459 · `Name` b1600 → `.ID`, `.CauseofLoss` | `Choose` b2222 → `setCauseOfLoss_DT` b2345: `CAUSEID ← id`, `CAUSE ← cause` | `GET …/master/penyebab?cari=` | ✅ paket 2 |

Tabel fisik = nama kelas `ASM-FW-GISFW-Int-<X>` (konvensi; preseden Retro Life `AGENT`, `BUSINESS`, `REINSURANCETYPE`) — objek yang tidak ada dijawab 503 yang menyebut objeknya (OQ-MPNL-04).

## 5. Daftar bersarang di `ProductName` (wadah b30995 / b2573)

| Grid (bNNN · judul VERBATIM) | Kolom korpus → kunci baris | Tombol / aksi (bNNN · teks) | Keadaan |
| --- | --- | --- | --- |
| b12201 **`LIEN CLAUSE (Potongan Manfaat Klaim)`** — `ProductName.LienClause` | b12741 `Usia saat Klaim` → `.Usia` b13054; b12890 `% Manfaat yang dibayarkan` → `.Manfaat` b13290 | ikon grid bawaan `pzPegaDefaultGridIcons` b12373 (vis `ProductName.IsView!='true'`), Enter = `addRow` b13620 — teks ikon **tidak ada di korpus** | ✅ paket 7 — ⭐ **daftar** (RALAT R10) |
| b14601 **`DOCUMENT CLAIM`** — `ProductName.DocumentClaim` | b15125 `Document List` → `.Document` b15286 (dropdown `associated`, daftar tak ikut ekspor) | ikon grid bawaan b14760 (vis `IsView!='true'`), Enter = `addRow` b15595 | ✅ paket 7 (OQ-MPNL-05) |
| b31557 **`PLAN LIST`** — `ProductName.PlanList` | b31845 `Plan Name` → `.Plan` b33121 (autocomplete RD `BrowseProductTypeLife_RD` b33198: `.CoverName` → `.Plan`, `.ID` → `.PlanID`, `.Business` → `.Name`, `.Benefit` → `.Benefit`; onChange `ProteksiPlanListLife` b33163); b31994 `Bussines` → `.Name` b33570; b32143 `Benefit` → `.Benefit` b33724; b32296 `R/I Rate` → `.RIRATE` b33887 | **`Add`** b32823 (`addRow` b32949); **`View Rate`** b34113 (vis `OTHER .RIRATE!=''`: `SetParamRate` b34310 + `localAction ViewRate` b34354); **`Choose R/I Rate`** b34589 (vis `OTHER IsView!='true'`, `localAction ChooseRIRate` b34806); **`Delete`** b35075 (`deleteRow` b35192 + `ProteksiPlanListLife` b35214) | ✅ paket 6 (`Choose R/I Rate`, `View Rate` ⏸️ OQ-MPNL-03) |
| `ProteksiPlanListLife`: 1 b258 `·` Page-Clear-Messages; 2 b347 `·` ulang `PlanList` (`errmsg = "Plan tidak boleh sama"` b421, `errmsg1 = "Plan tidak boleh kosong"` b442, `errmsg2 = "RI/RATE tidak boleh kosong"` b463); 2.1 b530 `·` ulang; 2.1.1 b565 `·` PRE=true `.Plan==""` T=2 F=3; 2.1.2 b715 `·` PRE=true `local.plan==.Plan && local.idx!=.pxListSubscript` T=2 F=3; 2.1.3 b876 `·` PRE=true `.RIRATE==""` T=2 F=3 | — | pesan VERBATIM menolak simpan | ✅ paket 6 |
| b37148 **`FINANCIAL UNDERWRITING`** — `ProductName.FinancialUnderwritingList` | b37670 `Min Insured` → `.MinInsured` b38718; b37818 `Max Insured` → `.MaxInsured` b38934; b37966 `Employee` → `.Employee` b39155; b38115 `Non-Employee` → `.Non_Employee` b39394 | **`Add`** b38425; ikon salin b39697 → `CopyFinancialWriting` 1 b225 `·` tambah baris salinan; **`Delete`** b40024 (+ `ProteksiPlanListLife`) | ✅ paket 7 |
| b42075 **`UNDERWRITING LIMIT`** — `ProductName.UnderwritingLimitList` | b42594 `Min Insured` · b42742 `Max Insured` · b42890 `Min Age` · b43038 `Max Age` · b43187 `Medical` · b43336 `Description` → `.MinInsured` b43939, `.MaxInsured` b44155, `.MinAge` b44373, `.MaxAge` b44592, `.Medical` b44814 (teks bebas), `.Description` b45052 | **`Add`** b43643; ikon salin b45355 → `CopyUnderWritingLimit` 1 b226 `·`; **`Delete`** b45682 (+ `ProteksiPlanListLife`) | ✅ paket 6 |
| Komentar — `ProductName.CommentList` (wadah b61044) | b62095 `Date` → `.Date` b62561; b62246 `PIC` → `.OperatorName` b62755; b62399 `Comment` → `.Suggest` b62937 | baca-saja (aksi baris hanya `setFocus`) | ✅ paket 7 |
| `AddCommentList_Act` 1 b235 `·` (dipanggil **setiap** simpan, `SaveProductName_Act` 7 b1515, tanpa prakondisi): `CommentList(<APPEND>).Date = @CurrentDateTime()`, `.OperatorName = OperatorID.pxInsName`, `.IsApproved = param.status` (tak dikirim), `.Suggest = param.comment ← ProductName.Comment` | `CommentList[*].{Date, OperatorName, IsApproved, Suggest}` | — | ✅ paket 7 |
| b48117 / b52732 **`OUTWARD`** — `ProductName.OutwardList` (wadah b47876 / b52491 `1==2`, `ExpressionCondition`) | `Reins Type` · `Transaction Year` · `Underwriting Year` → autocomplete RD `BrowseTreatyContract_Life_RD` | `Add` b49118 / b53733, **`View Reinstype`** b50764 / b55382 (`SetParamOutward` — tidak ada di korpus — + `RetrocessionModalDialog`), `Delete` b51305 / b55925 | ➖ grid **mati** (`1==2`); isinya tetap ditulis `GetReinsTypeOR_Life` (§3.1) |

## 6. Lampiran — `InboxProductName` wadah b64133, grid `TempData.AttachmentList`

Pesan statik VERBATIM: b66071 `Make sure the file name doesn't contain forbidden character such as , / | ' "`, b66230 `Recommended safe substitute should be . or _`. Kolom b68426 `File Name`.

| Unsur (bNNN · teks) | Rule dipanggil (bNNN · blok) | RDB → tabel | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| **`Add attachment`** b64747 (vis `ALWAYS TreatyIn.ViewState !='1'` → syarat diabaikan) | `SetCategory_act` b64990 (1 b231 `·` call `SetCategoryAttachTreatyin` — kategori treaty-in, ditimpa `"File"`); `localAction ProductNameAttachContent` b65034 → FlowAction: submit **`Attach`** b24, **`Cancel`** b22, pra `TreatyInitAttach` b31, pasca `ProductNameSaveAttachment` b179; refresh `LoadAttachmentProdName` b65098 | `ProductNameSaveAttachment`: 1 b292 `·` `Err = "Tidak ada file yg diattach"`; 2 b433 `·` ulang berkas; 2.1 b475 `·` `.pyCategory = "File"`; 2.2 b626 **`//`**; 2.3 b773 `·`; 2.4 b904 `·` `InsertGoogleStorage_Act` (`Durasi="1800"`, `Folder="Contract"`; *ralat audit 02-10-2026: transisi gagal b935 ber-`pyStepsTransition=false`, jadi kegagalan TIDAK keluar - penjaga sebenarnya 2.6 b1201 `·` PRE b1347 `Datain.CARI51==""` T=3: rekam `M_ATTACHMENTPRODUCTNAME` dilewati bila ImageID tidak kembali b960*); 2.5 b1075 `·` `CARI50 = ""`; 2.6 b1201 `·` PRE=true `CARI51==""` T=3 F=2 RDB `InsertAttachProdName_Sql`; 3–8 **`//`** | `InsertAttachProdName_Sql` b84 `INSERT INTO M_ATTACHMENTPRODUCTNAME (ID = TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF3'), TREATYID = ProductName.ID, CATEGORY, FILENAME, FILEMIMETYPE, DATA_JSON = "", USERNAME, T_STORAGE_ID = ImageID)`; `Insert_T_Storage_SQL` b85 `T_STORAGE_IMAGE (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE='standard')`; `GenerateImageID_SQL` b79 MD5 `'ASMPP'‖stempel‖SYS_GUID()`; `GetAppName_SQL` b58 `T_FOLDER_IMAGE.APPNAME` | `POST …/produk/{id}/lampiran` | ✅ paket 8 — unggah = stub outbox (P5, OQ-MPNL-10) |
| **`Refresh`** b65270 (vis `ALWAYS !pyIsMobile`) | `LoadAttachmentProdName` b65401: 1 b253 `·`; 2 b372 `·` RDB `GetAttachmentProdName_Sql`; 3 b543 `·`; 4 b672 `·`; 5 b773 **`//`** | `GetAttachmentProdName_Sql` b84 `select id, FILENAME, CATEGORY, FILEMIMETYPE, T_STORAGE_ID as "type" from M_ATTACHMENTPRODUCTNAME where treatyid = ProductName.ID` | `GET …/produk/{id}/lampiran` | ✅ paket 8 |
| `Download` b67376 (vis `OTHER FALSE`) | `DownloadAll_Act` | — | — | ➖ **mati** |
| **`Download All`** b67657 | DT `TreatyInIDSetPyPortal` b67959 + `showHarness InputTreatyInOffer` (kelas `Data-Portal`, harness tidak ada di korpus) pra `TreatyInDownloadAll` (1 b225 `·` `SetTreatyIn_Act`; 3 b431 `·` `DownloadAll_Act` → `GetAllAttachment2_Sql` b84 `M_ATTACHMENTTREATY_2 where treatyid = TreatyIn.ID`) | **jalur treaty-in** (salah ekspor, OQ-056) | `GET …/produk/{id}/lampiran/unduh-semua` (zip lampiran produk) | ✅ paket 8 — perilaku dari AC tiket 08 (RALAT R15, OQ-MPNL-07) |
| Tautan nama berkas b68903 (`.pyFileName`) | `DownloadAttProdName_Act` b68926 (`ImageID ← .type`): 1–5 **`//`**; 6 b953 `·` `GetUrlGoogleStorage_Act` (`Durasi="1800"`); 7 b1080 `·` PRE=true `Param.ViewOffice` T=2 F=3 | `GetLinkStorage_SQL`, `Update_T_Storage_SQL` b85 | `GET …/lampiran/{id}/unduh` | ✅ paket 8 — berkas dari penyimpanan stub |
| Tautan **`View Office Online`** b69291 (vis `OTHER .pyFileMimeType = xls/xlsx/doc/docx/ppt/pptx`) | `DownloadAttProdName_Act` b69314 dengan `ViewOffice`: langkah 7 b1080 membungkus URL bertanda tangan ke penampil kantor **di luar** | — | `GET …/lampiran/{id}/office` | ⏸️ paket 8 — **stub**, alamat luar tidak dipanggil dan tidak ditulis (OQ-MPNL-11) |
| **`Delete`** b69714 (vis `ALWAYS TreatyIn.ViewState !='1'`) | `DeleteAttacProdName_act` b69850 (`ID ← .ID`, `ImageID ← .type`): 1 b281 `·`; 2 b411 `·` `DeleteGoogleStorage_Act` (gagal → keluar b444); 3 b528 `·` RDB `DeleteAttachProdName_Sql`; 4 b700 `·` `LoadAttachmentProdName`; 5–7 **`//`** | `DeleteAttachProdName_Sql` b85 `delete M_ATTACHMENTPRODUCTNAME where treatyid = … and id = …`; `DeleteStorage_SQL` b85 `delete T_STORAGE_IMAGE where imageid = …` | `DELETE …/lampiran/{id}` | ✅ paket 8 |
| `GetMimeType` (DecisionTable, dipanggil `InsertGoogleStorage_Act` b761) | 48 baris + `otherwise` | `inti/backend/unggah.MimeDariNamaFile` (salinan tabel yang sama) | — | ✅ paket 8 |
| `ServiceGoogle` (ConnectREST, `pyBaseURLSetting = LinkService!LinkService` b161) · `LinkService` (SystemSettings) · `GetLinkService` · `GetTokenStorage_SQL` (`GET_TOKEN_STORAGE`) | — | — | — | ⏸️ paket 8 — **stub outbox**; nol alamat di berkas (ADR-0013) |

## 7. Simpan — `SaveProductName_Act` (FlowAction `SaveProductName_Confirm` b101)

| Langkah (bNNN · blok) | Isi | Sistem baru | Keadaan |
| --- | --- | --- | --- |
| 1 b361 `·` **PRE=false** (WHEN `TYPE=="" ‖ GRUP==""` b630 tidak dievaluasi) | `IDPEGA ← ProductName.ID`; pesan `"Product Name Empty"` b431, `"Ceding Empty"` b452, `"Policy Holder Empty"` b473, `"SOB Empty"` b494; `ProductName.POLICYHODER/POLICYHODERNAME ← ProductNameInward.*`; `UPDATEOP ← OperatorID` | salinan pemegang polis ke JSON umum; `UPDATEOP` = pelaku. ⛔ **bukan** validasi tipe/grup (RALAT R7) | ✅ paket 3 |
| 2 b670 `·` PRE=true `PRODUCTNAME==""` b805 F=3, trans `1==1`→6 b714 | `Property-Set-Messages` pesan `local.errMsg1` | 422 `Product Name Empty` | ✅ paket 5 |
| 3 b845 `·` PRE=true `CEDING==""` b980 | `local.errMsg2` | 422 `Ceding Empty` | ✅ paket 5 |
| 4 b1020 `·` PRE=true `ProductNameInward.POLICYHODERNAME==""` b1155 | `local.errMsg3` | 422 `Policy Holder Empty` | ✅ paket 5 |
| 5 b1195 `·` PRE=true `SOBNAME==""` b1330 | `local.errMsg4` | 422 `SOB Empty` | ✅ paket 5 — Pega berhenti di pesan pertama (trans →6); sistem baru melaporkan **semua** (AC tiket 05) |
| 6 b1370 `·` PRE=true `@PropertyHasValue(CREATEOP)` b1475 T=3 F=2 | `CREATEOP ← OperatorID` bila kosong | `CREATEOP` = pelaku untuk produk baru, tetap untuk ubah | ✅ paket 3 (OQ-MPNL-13 salinan) |
| 7 b1515 `·` | call `AddCommentList_Act(comment ← ProductName.Comment)` | satu baris `CommentList` per simpan | ✅ paket 7 |
| 8 b1625 `·` **PRE=false** (WHEN `PoductName.*` b1793 tidak dievaluasi) | `DATAPEGA ← @GetPageJSONString()` halaman `ProductName`; `CARI1..3 ← ID, RIRISKID, RIRISK` | JSON umum dirakit repository (kunci Pega) | ✅ paket 3 |
| 9 b1833 `·` **PRE=false** | RDB `SaveProductNameLIfe` b58 → `PEGA_M_PRODUCT_LIFE` + `COMMIT` | ⛔ prosedur **tidak** dipanggil; upsert `M_PRODUCT_LIFE (ID, JSONDATA)` ditiru, ID baru `'1' ‖ LPAD(M_PRODUCT_LIFE_SEQ.NEXTVAL, 5, '0')` (`dba-procedures-and-ddl.md` §1) | ✅ paket 3 |
| 10 b2021 `·` **PRE=false** | RDB `SaveProductNameLIfeFlat` b84 `UPDATE M_PRODUCT_LIFE SET RIRISKID, RIRISK WHERE ID` | kolom datar `RIRISKID`, `RIRISK` ditulis di transaksi yang sama — *kalimat lama dikutip:* "(+ `PRODUCTNAME`, `BEGIN_DATE`, P1)"; **ralat 01-10-2026 (L1):** kedua kolom itu tidak ada di DEV, tidak ditulis | ✅ paket 3 / 4 — ⭐ **hidup** (RALAT R8) |
| 11 b2209 `·` **PRE=false** | `ProductName.ID ← OutputParam.IDPEGA`; `DATASHOW = 0` | ID hasil simpan dikembalikan; form ditutup ke grid | ✅ paket 3 |
| 12 b2375 `·` | kosongkan `IDPEGA`, `DATAPEGA` | — | — |
| 13 b2530 `·` **PRE=false** | `IDPEGA ← ProductNameInward.ID`; `ProductNameInward.PRODUCTID ← IDPEGAOUT` | `PRODUCTID` = ID produk | ✅ paket 4 |
| 14 b2719 `·` **PRE=false** | `DATAPEGA ← @GetPageJSONString()` halaman `ProductNameInward` | JSON inward dirakit repository | ✅ paket 4 |
| 15 b2864 `·` **PRE=false** | RDB `SaveProductNameInwardLIfe` b84 → `PEGA_M_PRODUCT_INWARD_LIFE` + `COMMIT` | ⛔ prosedur **tidak** dipanggil; upsert `M_PRODUCTINWARD_LIFE (ID, JSONDATA)` di transaksi yang **sama** dengan langkah 9 (P4); `ID` inward = `ID` produk (P6, OQ-MPNL-02) | ✅ paket 4 |
| 16 b3052 `·` **PRE=false** | `ProductNameInward.ID ← IDPEGAOUT` | — | ✅ paket 4 |
| `SaveInwardProductName_Act` 1 b240 `·` · 2 b402 **`//`** `error Type` · 3 b589 **`//`** `error Grup` · 4 b764 `·` · 5 b909 `·` · 6 b1097 `·` | dipanggil hanya `Save` b14930 harness / b14125 section `InwardProductName` | — | ➖ tak terjangkau (RALAT R11) |

### 7.1 Kunci yang dibaca view (dari `dba-view-produk-life.md`) — wajib ada di JSON hasil simpan

| View | Kunci JSON |
| --- | --- |
| `PRODUCT_LIFE` ← `M_PRODUCT_LIFE` | `TYPE`, `TYPE_CEDING`, `CEDING`, `CEDINGID`, `SOBNAME`, `SOBID`, `CAUSEID`, `GRUP`, `PRODUCTNAME`, `PRODUCTCODE`, `PRODUCTTYPEID`, `PRODUCTTYPE`, `RIRISKID`, `RIRISK`, `RIRATEID`, `RIRATE`, `RICOMMID`, `RICOMM`, `INWARDNAME`, `UnderwritingLimitList`, `OUTWARDNAMEID`, `OUTWARDNAME`, `OUTWARDRATEID`, `OUTWARDRATE`, `OUTWARDCOMMID`, `OUTWARDCOMM`, `BENEFITID`, `BENEFIT`, `CAUSE`, `OutwardList[0].OVR_COMM`, `POLICYHODER`, `POLICYHODERNAME`, `TREATYNUMBER`, `CREATEOP`, `UPDATEOP` |
| `DOCUMENTCLAIM_LIFE` ← `M_PRODUCT_LIFE` | `DocumentClaim[*].Document` |
| `PRODUCTINWARD_LIFE` ← `M_PRODUCTINWARD_LIFE` | `PRODUCTID`, `INSURED`, `CEDING`, `TREATYNUMBER`, `ADDENDUMWORD`, `AMANDEMENTSCHD`, `INWARDTREATYNM`, `BEGIN`, `MATURE`, `CEDINGRETENTIONNUM`, `CEDINGRETENTIONPCT`, `CEDINGLIMIT`, `CEDINGLIMITXPN`, `MINAGE`, `MAXAGE`, `BIRTHDAY`, `EXTRAPREMI`, `CURRENCY`, `RNMSHARE`, `EXTRAMORTALITY`, `RNMLIMITNUM`, `RNMLIMITPCT`, `LIENCLAUSE`, `MONTHS`, `MINSUMINSURED`, `MAXSUMINSURED`, `MAXSUMREASURED`, `MAXCONTRACT`, `PAYMENT`, `PROPORTIONALTABLE`, `SUBJECTTO`, `POLICYHODER`, `POLICYHODERNAME`, `BROKERAGE`, `ADDENDUMNO`, `AMANDEMENTNO`, `MAXDATARECEIVE`, `MAXEXPIREDCLAIM`, `STNC` |

Kunci skalar tanpa nilai ditulis `""` (Oracle membacanya NULL, sama dengan kunci absen); kunci daftar ditulis `[]`. Setiap baris
`OutwardList` membawa `OVR_COMM` (`""` — korpus tidak punya penulisnya), setiap baris `DocumentClaim` membawa `Document`.
Kunci yang sudah ada di JSON lama dan tidak dikelola layar **dipertahankan** apa adanya (Pega memuat halaman utuh lewat
`adoptJSONObject`, `SetProductName` 4 b959).

## 8. Rule yang sengaja tidak dibangun

| Rule | Bukti | Keadaan |
| --- | --- | --- |
| `SetTreatyIn_Act`, `TreatyInInputVis`, `CheckDuplicateOffer`, `GetCountClaim`, `GetCurrentDate`, `BrowseTreatyIn`, `SaveTreatyIn`, `BrowseTREATY_IN`, `TreatySetReinstatement`, `SetReinstatementPct`, `ConvertHistoryDate`, `LoadAttachment`, `GetAttachment2_Sql`, `TreatyInDownloadAll`, `DownloadAll_Act`, `GetAllAttachment2_Sql`, `TreatyInIDSetPyPortal` | jalur treaty inward (kelas `Data-Portal` / `TREATY_IN`) — **salah ekspor** (OQ-056 `[keputusan work owner]`); satu-satunya pemicu di layar produk adalah `Download All` b67657 (§6) | ➖ |
| `SaveInwardProductName_Act`, harness + section `InwardProductName` | hanya dibuka `Inward` b75368 `OTHER 1=2` | ➖ |
| `CountMaxReasured_Act` | dipanggil hanya `TYPE_CEDING` b3385 (`OTHER 1=2`) | ➖ |
| `SetParamRate` (`ParamID.RIRATEID ← InputBusinessLife.RIRATEID`, kelas `@baseclass`) | halaman `InputBusinessLife` milik Retro Life — di layar ini `ParamID` diisi baris plan (`View Rate` b34310) | ⏸️ `View Rate` OQ-MPNL-03 *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: ✅ `View Rate` membaca view `RATE_LIFE` disaring `RIRATEID` baris plan. ⚠️ Penyimpangan sadar: grid `ViewRate.xml` b1024 menyaring `ParamID.OUTWARDRATEID`, yang tidak diisi rule korpus mana pun — `SetParamRate` b259 hanya mengisi `ParamID.RIRATEID` dari halaman Retro — sehingga di Pega dialog ini tidak pernah menampilkan rate baris plan)* |
| `When/recordEvent`, `SetCategoryAttachTreatyin` + `CategoryAttach_SQL` | `recordEvent` nol pemanggil; kategori lampiran ditimpa `"File"` (`ProductNameSaveAttachment` 2.1 b475) sesudah `SetCategory_act` mengisi daftar kategori treaty-in | ➖ |

## 9. Rute API dan komponen

Prefix `/api/master-product-name-life` (`backend/handlers/rute_mpnl.go` `Prefix`). Setiap rute menuntut identitas pelaku (401 tanpa),
menjawab 503 bila Oracle tidak dikonfigurasi, galat berbadan `{"galat": "..."}`.

| Metode dan jalur | Layar / tombol Pega | Paket |
| --- | --- | --- |
| `GET /produk` | grid `InboxProductName` mode daftar (b71246, RD `BrowseProduct_Life`) | 1 |
| `GET /produk/{id}` | tombol `View` b74798 (`SetProductName` + `SetProductNameInward`) | 1 |
| `GET /master/{jenis}?cari=` — `ceding`, `sob`, `pemegang-polis`, `mata-uang`, `ri-risk`, `penyebab`; `ri-rate` = 503 (OQ-MPNL-03) | tombol `Choose*` → section pemilih → grid RD (§4); autocomplete medan form | 2 *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: `ri-rate` = 200, view `RATE_LIFE_SUMMARY`)* |
| `POST /produk` · `PUT /produk/{id}` | `Add` b71865 / `View` → `Edit` b59489 → `Save` b59041 → `SaveProductName_Confirm` → `SaveProductName_Act` (§7); nol rute hapus (korpus tanpa hapus produk) | 3 |
| (rute yang sama) | sisi inward `SaveProductName_Act` 13–16 — `M_PRODUCTINWARD_LIFE` di transaksi yang sama (P4); uji tiga view (`db`) | 4 |
| (rute yang sama) | wajib-isi `SaveProductName_Act` 2–5 — `Product Name Empty`, `Ceding Empty`, `Policy Holder Empty`, `SOB Empty` (422, semua sekaligus) | 5 |
| `GET /master-plan?cari=` · `GET /rate?riRateId=` (503, OQ-MPNL-03) | autocomplete `Plan Name` b33121 · tombol `View Rate` b34113; gerbang `ProteksiPlanListLife` dan `UNDERWRITING LIMIT` di simpan | 6 *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: `GET /rate?riRateId=` = 200, view `RATE_LIFE`, `{daftar, total, terpotong}`, 500 baris `BrowseRateLife_RD` b730)* |
| (rute yang sama) | `AddCommentList_Act` setiap simpan; `LIEN CLAUSE`, `DOCUMENT CLAIM`, `FINANCIAL UNDERWRITING` disimpan bersama produk | 7 |
| `GET /produk/{id}/lampiran` · `POST /produk/{id}/lampiran` (multipart `berkas`) | `Refresh` b65270 / `View` · `Add attachment` b64747 → `Attach` b24 (`ProductNameSaveAttachment`) | 8 |
| `POST /produk/{id}/lampiran/{lid}/ulangi` | (tiket 09) kirim ulang efek outbox yang gagal | 8 |
| `GET /produk/{id}/lampiran/{lid}/unduh` · `GET /produk/{id}/lampiran/unduh-semua` | tautan nama berkas b68903 · `Download All` b67657 (zip, R15) | 8 |
| `GET /produk/{id}/lampiran/{lid}/office` (503 stub, OQ-MPNL-11) · `DELETE /produk/{id}/lampiran/{lid}` | `View Office Online` b69291 · `Delete` b69714 (`DeleteAttacProdName_act`) | 8 |
| `POST /produk` + `salinanDari` | `Copy` b59854 (`CopyProduct` b144/b173) — medan milik server diwarisi produk asal | 9 |
| `POST /produk` · `PUT /produk/{id}` + `hitungOutward` | checkbox `On Retention` b47312 → `GetReinsTypeOR_Life` b47488 (`BrowseReinstypeOR_SQL` b84) | 9 |
| `POST /produk/generate` → `SeeDetail.csv` | `Generate` b60122 → `GenerateUpload_Act` b60228 | 9 |

## 10. Komponen layar (paket 10) dan cek peramban

| Komponen (`frontend/`) | Unsur korpus |
| --- | --- |
| `menu.ts` `HALAMAN_AWAL_MPNL = 'mpnl-produk'` | tombol menu **Master Product Name Life** (GROUPMENU `MASTER`, slot `960_menu_masterproductnamelife.sql`) → `InboxProductName` (§1) |
| `pages/MasterProductNameLife.tsx` | mode daftar: `End Period` / `Add` / `Add New Data` b71865, grid `BrowseProduct_Life` 10 baris/halaman, `View` b74798 (§2) |
| `components/FormProduk.tsx` | mode form: §3.1, §3.2, §5, tombol bawah §3.3; `SetTreatyName_Act`, `CountMaxSumReasured_Act`, `CopyProduct`, `CopyFinancialWriting`, `CopyUnderWritingLimit` di klien (`bentuk.ts`) |
| `components/Saran.tsx` | `pxAutoComplete` `Ceding`, `SOB`, `R/I Risk Name`, `Policy Holder`, `Currency`, `Plan Name` |
| `components/PemilihMaster.tsx` | ketujuh `Choose*` → `Search` / grid `ID` · `Name` (`RIRate Name`) / `Choose` / `Submit` · `Cancel` (§4) |
| `components/Dialog.tsx` | `SaveProductName_Confirm` (`Do you want to save the data?`, `Comment`, `Save` · `Cancel`), `EditProductName_Confirm` (`Do you want to Edit the data?`, `Edit` · `Cancel`) |
| `components/ModalRate.tsx` | `View Rate` → `ViewRate` (`Outward List`, `ID` · `USEDBY` · `GENDER` · `CONTRACT` · `AGE` · `RATE`, `Submit` · `Cancel`) — data ⏸️ OQ-MPNL-03 *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: data view `RATE_LIFE` ✅)* |
| `components/PanelLampiran.tsx` | §6 — hanya untuk produk ber-ID (tiket 08: produk dulu, lampiran menyusul) |

⛔ Penjaga: `labels.test.ts` membuka korpus pada 131 baris tag label + 5 baris ekspresi/nilai berkutip, memastikan setiap
tombol/tautan yang dirender bertekskan label berbukti dan ke-37 tombol korpus hidup (25 `pxButton` `InboxProductName`, tautan
`View Office Online`, `Choose` pemilih, 10 tombol FlowAction) seluruhnya dirender. Delapan `pxButton` mati tidak dirender *(ralat audit 02-10-2026: dulu tertulis "Tujuh"; 33 − 25 = 8)*:
`Inward` b75368 (`1=2`), `Add`/`View Reinstype`/`Delete` grid `OUTWARD` ×2 (wadah `1==2`), `Download` b67376 (`OTHER FALSE`).

**Cek peramban 01-10-2026** (Chrome tanpa kepala lewat CDP, backend lokal TANPA `ORACLE_DSN` di port sendiri, Vite
`--port 5199`): sidebar memuat tombol `Master Product Name Life`; grid menampilkan `End Period` + `Add` dan galat
`database is not configured` dinyatakan; `Add` membuka form dengan 47 label medan/judul VERBATIM dan tombol `Choose Ceding
Name`, `Choose SOB`, `Choose R/I Risk`, `Choose Cause of Loss`, `Choose Policy Holder`, `Choose Currency`, ikon grid,
`Add` ×3, `Choose R/I Rate`, `Delete`, `Copy row`, `Close`, `Save`, `Copy`, `Generate` (`Edit` tersembunyi — `IsView`
false); `Save` membuka dialog berpertanyaan + `Comment` (`Cancel`/`Save`); `Choose Ceding Name` membuka pemilih (`Search`,
`Cancel`/`Submit`, galat basis data dinyatakan); `Copy` menampilkan `Data sudah dicopy, silakan melakukan perubahan dan
tekan SAVE untuk menyimpan`. Tanpa basis data tidak teramati: `View` baris grid, `Edit` (mode lihat), `View Rate`
(butuh `RIRATE`), panel lampiran (butuh produk tersimpan) — keempatnya dijaga `labels.test.ts`.

## 11. Commit per paket (cabang `dev`)

| Paket | Commit | Tiket | Isi |
| ---: | --- | --- | --- |
| 0 | `71c35b1` | 01 (ditangguhkan) | PARITAS, `RALAT-DEV-30-09-2026.md` P1–P6 + R7–R18, OQ-MPNL-01..14, ralat bertanggal spec dan tiket 01 |
| 1 | `4ce0771` | 02 (baca) | kerangka backend, baca `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE` lewat `JSONDATA`, rute baca, penjaga modul |
| 2 | `a94905e` | 04 | tujuh pemilih master — saringan/urutan RD; R/I Rate ⏸️ OQ-MPNL-03 |
| 3 | `2f39341` | 02 | sisi umum, identitas dari sequence, upsert `JSONDATA` + kolom datar, gagal terang |
| 4 | `b179aa2` | 03 | sisi inward, simpan atomik dua tabel, uji tiga view DEV di skema uji |
| 5 | `6e49cc6` | 05 | wajib-isi — empat pesan VERBATIM, semua sekaligus |
| 6 | `05ca173` | 06 | `PLAN LIST`, `UNDERWRITING LIMIT`, `ProteksiPlanListLife` VERBATIM |
| 7 | `e124f18` | 07 | komentar setiap simpan, `DOCUMENT CLAIM`, `LIEN CLAUSE`, `FINANCIAL UNDERWRITING` |
| 8 | `1ff3091` | 08, 09 | lampiran tabel warisan, efek keluar stub lewat outbox `T_LOG_SERVICE_RNM` |
| 9 | `14d3d48` | 10 (tambahan) | `Copy`, `On Retention` → `OutwardList`, `Generate` → `SeeDetail.csv` |
| 10 | `1ada8d7` | semua (layar) | menu slot 960, `modul.go`, layar `InboxProductName` lengkap, `MODUL.md` dimigrasi |
| 11 | `bb1c96e` | — | status tiket, daftar commit, register OQ |

## 12. Perbaikan sesudah `/code-review` (01-10-2026)

| # | Temuan | Keputusan |
| ---: | --- | --- |
| 1 | baris inward ber-`ID` = produk tetapi `PRODUCTID`-nya produk lain dipilih lalu ditimpa | **diperbaiki** — dipilih hanya bila `PRODUCTID` kosong/sama; sisip baris inward ber-ID itu = `ErrIdentitasBentrok` berkalimat (menyebut pemiliknya) |
| 2 | lampiran `belum` tanpa tombol kirim ulang (modul tanpa pekerja) | **diperbaiki** — `Retry` untuk setiap status selain terunggah |
| 3 | `Copy` membuang kunci halaman yang tidak dikelola layar | **diperbaiki** — JSON tersimpan produk asal menjadi dasar kedua sisi (`CopyProduct` menyalin halaman utuh) |
| 4 | `ProteksiPlanListLife` ditegakkan setiap simpan atas baris lama yang tidak disentuh | **diperbaiki** — berjalan bila `PLAN LIST` berubah atau baris `FINANCIAL UNDERWRITING`/`UNDERWRITING LIMIT` berkurang (padanan onChange `.Plan` b33163 dan `Delete` b35214/b40024/b45682) |
| 5 | konversi `@replaceAll` `SetProductNameInward` b1048 | **tidak diubah** — hasilnya hanya variabel `Local.*` yang tidak dipakai lagi; halaman menerima nilai mentah; contoh `[data DBA]` berbentuk angka polos |
| 6 | `Download All` gagal seluruhnya bila satu lampiran lama tidak ada di stub | **diperbaiki** — berkas yang ada di-zip, yang tidak ada dicantumkan di `_not-available.txt`; nol berkas = 409 `ErrBerkasTidakDiStub` (bukan anjuran menghapus) |
| 7 | satu `ulangi` = dua percobaan | **diperbaiki** — efek antre dipakai ulang tanpa dipungut dua kali |
| 8 | galat multipart ditelan menjadi "Tidak ada file yg diattach" | **diperbaiki** — 413 berkalimat / 400 multipart rusak |
| 9 | `Premium Payment Method` dapat diubah di mode lihat | **diperbaiki** — baca-saja (`ro` b25611) |
| 10 | autocomplete/pemilih membaca hingga 100.000 baris | **diperbaiki** — `?batas=` (autocomplete 20, kosong tidak membaca), grid pemilih berhalaman 10 |
| 11 | grid membaca CLOB utuh 500 produk | **diperbaiki** — kelima kolom grid lewat `JSON_VALUE` |
| 12 | kirim lampiran tanpa jejak "menyerah" | **tidak diubah** — jalur jejak bersama menulis tabel jejak Claim Life; kegagalan tercatat di outbox (`gagal-permanen` + `GALAT_TERAKHIR`), tampil di layar, dapat diulang (tiket 09) |
| 13 | duplikat plan diperiksa sebelum nama diseragamkan master | **diperbaiki** — diperiksa sesudahnya |
| 14 | ubah membaca kedua tabel tiga kali | **tidak diubah** — performa saja, tabel ±200 baris, semuanya di satu transaksi |
| 15 | badan JSON tanpa batas; `asli` baris dikarang klien | **diperbaiki** — 4 MiB / 413; `asli` harus sama dengan `asli` salah satu baris tersimpan (produk yang diubah atau produk asal `Copy`), selain itu 422 |

## 13. Lanjutan 1 (01-10-2026, `PROMPT-LANJUTAN-MASTER-PRODUCT-NAME-LIFE-1.md`)

| Paket | Commit | Isi |
| ---: | --- | --- |
| 1 | `6fd539c` | L1 — INSERT/UPDATE `M_PRODUCT_LIFE` hanya `ID`, `JSONDATA`, `RIRISKID`, `RIRISK` (`SaveProductNameLIfeFlat` b84); uji katalog `testdata/katalog-dev.json` |
| 2 | `6209638` | L2/L3/L5/L7 — inward ber-ID produk, objek master DEV, `OVR_COMM` dan `TREATYCONTRACTID` kosong |
| 3 | — | L8 tidak dijalankan: register OQ tidak memuat izin bertanggal work owner untuk OQ-MPNL-03 — `R/I Rate`/`View Rate` tetap 503 |
| 4 | `e1ff8b8` | uji manual baca-saja terhadap DEV (`LAPORAN-UJI-MANUAL.md`) |
| 5 | `7b0fa8e` | dokumen §2, L6 (`Medical` teks bebas), L9 (nomor `bNNN` perintah standar) |
| tinjauan | `43816c2` | perbaikan `/code-review`: panjang atas nilai akhir, batas 4000 byte kunci view, katalog bertipe, DDL tiruan = katalog, uji `db` kolom datar + inward; L9 putaran kedua **26** nomor langkah yang tertinggal (pesan commit-nya keliru menulis 23) |

## Keputusan OQ work owner 01-10-2026 — K1 (`PROMPT-LANJUTAN-TIGA-MODUL-LIFE-KEPUTUSAN-OQ.md` §1)

Nomor `bNNN` di bab ini dari `sed -e 's/></>\n</g' "Master Product Name Life/<Tipe>/<Rule>.xml" | grep -n '<teks>'`.

| Unsur | Rule / baris | Dibangun |
| --- | --- | --- |
| `Choose R/I Rate` b34589 → `RIRate_Section` | grid RD `BrowseRateLifeSummary` b2707; param `id` kosong b1461, `idusedby = SearchPolicyHolder.CARI1` b1466; RD `.USEDBY Contains` b807, urut `.ID ASC` b694, `pyMaxRecords` 500 b674; `Choose` b2424 → `SetRIRate` b2448 (`.ID` b2463, `.USEDBY` b2469) | `GET /master/ri-rate?cari=` — view `RATE_LIFE_SUMMARY`, `SELECT ID, USEDBY` |
| `View Rate` b34113 (`SetParamRate` b34310, `localAction ViewRate` b34354) | section `ViewRate` judul b843, grid RD `BrowseRateLife_RD` b3032; param b1024 `ParamID.OUTWARDRATEID`; RD `.IDUSEDBY = Param.idusedby` b859/b868, urut `.ID DESC` b748, `.RATE ASC` b786, `pyMaxRecords` 500 b730; kolom grid b2035–b2809 | `GET /rate?riRateId=` — view `RATE_LIFE`, enam kolom, `IDUSEDBY = :1`, 500 baris + `terpotong`; disaring `RIRATEID` baris plan (penyimpangan sadar, bab 6 baris `SetParamRate`) |
| Simpan baris `PLAN LIST` | `ProteksiPlanListLife` (`RI/RATE tidak boleh kosong`) | pasangan R/I Rate **baru** wajib ada di `RATE_LIFE_SUMMARY`; nama = `.USEDBY` master; ID di luar view / nama ketikan tanpa pilihan = 422 berkalimat |
| Penjaga | — | `periksaBacaSaja` (runtime, setiap objek `DaftarMasterDibacaSaja`), `TestMPNLSetiapSQLMasterAdalahSelect`, `TestMPNLRateDibacaKolomRDSaja`, gigit `TestMPNLPeriksaBacaSajaMenolakTulisanKeView`; mutasi `SELECT *` dan `UPDATE` = merah |
