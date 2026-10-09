# PARITAS — tombol dan aksi XML Claim Non Prop ↔ status

> 09-10-2026, tahap 1 dari 2. Sumber inventaris: Section `OutstandingClaim(1)`, `InputAcceptation`, `AdjustmentDetailNP`,
> `AdjustmentDetailNP_Section`, `Subjectivity`, `KomiteCLMNP`, `CloseClaimMD`, `CloseClaimNP`, `PreviewPLA`,
> `InputDtlInterest`, `ReinstatementPremiumDetails`, `EditXOLAlokasi`, harness `ChooseMasterTNonProp`, `CauseofLoss_Harness`,
> `Hitung_Test`, `ViewOldAllocation`, `ViewAttachmentNP`, `ViewHistoryMasterID_NP`, `KomiteCNP`, flow action
> `ViewListPolicyCNP`, `CatastrofeList`, `GeneratePLACNP`, flow `Flow_TreatyIn` di `D:\XML\RNM_BRD\Claim Non Prop`
> (setiap pxButton / pxLink dan setiap medan ber-`on change` activity). Sumber status: definisi layar server
> `backend/models/layar.go` dan `layar_akseptasi.go` — satu-satunya tempat tombol lahir; frontend merender pohon tata
> apa adanya (`frontend/components/TataView.tsx`).
>
> **Status:**
> - **dibangun** — tampil menurut kondisi XML yang sama (visible-when / disabled-when / read-only-when), aksinya port
>   activity-nya.
> - **nonaktif-OQ** — tampil sesuai section, tetapi nonaktif; keterangan `OQ-CNP-nn` di atribut `title` (`OQ.md`).
> - **tidak tampil di XML** — kondisi tampilnya statis `1=2` / `NEVER`: tidak dibangun.
> - **tahap 2** — milik modul `komiteclaimnonprop` (prompt tahap 1 §1).
>
> Cek layar 09-10-2026 (server uji sementara di atas gudang tiruan `UJI-`, port 18098; vite `--port 5178 --strictPort`;
> Chrome tanpa kepala, nol DEV): halaman awal (Process / Resolve, switch Teknik), Outstanding Claim kasus baru dan kasus
> terisi (tab Claim Information / Interests / Estimation / Spreading), pop-up Choose Master Treaty XOL, ViewListPolicyCNP,
> modal GeneratePLACNP, Input Acceptation (tab Acceptation dengan panel AdjustmentDetailNP terbuka), modal Send To
> Commitee, CloseClaimMD, CloseClaimNP, lebar 390 px. Nol galat konsol sesudah perbaikan `key` Popup (lihat §9).

## 1. Flow dan halaman awal

| XML | Aksi XML | Status |
| --- | --- | --- |
| `Flow_TreatyIn` Start1 → Assignment2 "Outstanding Claim" (worklist pembuat) | pembuatan kasus | **dibangun** — tab *Process*, switch Teknik mati = worklist pembuat (`ToCurrentOperator`, tanpa cek workbasket); **Add Claim** hanya saat switch mati (OQ-CNP-23: harness `New` tidak diekspor, pola Claim Prop) |
| Assignment1 "Input Acceptation" (workbasket `TreatyinPNCTeknik`) | FlowAction `InputAcceptation` | **dibangun** — switch Teknik nyala = workbasket `ReasKlaimTeknik` (OQ-CNP-08 bawaan); switch nonaktif bagi akun tanpa workbasket itu (`GET /hak`) |
| End Resolved-Completed | `CloseClaimTNonProp` / `ASMForceCaseClose` | **dibangun** — tab *Resolve* |
| FlowAction `OutstandingClaim` pra-proses (`InputOutStandingClmTNP_PreAct`) | — | **dibangun** (`services.siapkan`, `models.PraOutstanding`) |
| FlowAction `InputAcceptation` pra-proses | — | **dibangun** (`models.PraAkseptasi`) |
| Kotak masuk Beranda | — | **tidak dibangun** — bukan bagian Pega (pola Claim Prop; `menu.ts` tanpa `antreanBeranda`) |
| Tabel komite di bawah inbox | — | **tahap 2** |

## 2. Section `OutstandingClaim` (Assignment2)

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| Choose Master In | harness `ChooseMasterTNonProp`, `BrowseDtlMasterTNP_Act` TreatyIN=IN | **dibangun** — pop-up, Choose = `SetValueClaimTNP_Act` (ID, BusinessName, BusinessCode, TreatyGroup dibaca ULANG di server); saringan per kolom di server + batas 500, 50 per halaman (pola Claim Prop, prompt §8 P2) — isian "Input Treaty ID" + Search section menjadi saringan kolom Treaty ID |
| Choose Master Out | TreatyIN=OUT | **tidak tampil di XML** (`1=2`) |
| View Master | harness `InputTreatyInOffer` | **nonaktif-OQ** OQ-CNP-20 (harness tidak diekspor) |
| Close Without Payment (CWP) | local action `CloseClaimNP` | **tidak tampil di XML** (`1=2`) di Outstanding |
| Policy No (autocomplete) | `CheckNoPolicy` NOPOLIS + COB | **dibangun** (+ `SetMOClaimTreaty`) |
| View List Policy | local action `ViewListPolicyCNP` | **dibangun** — pop-up, tautan Policy No = `CheckNoPolicy` lalu tutup |
| View | (tanpa aksi di section; harness polis `ViewDetailDeptHeadTreatyIn_UW`) | **dibangun** — jendela Modal berkas NB / EDM Treaty In (OQ-CNP-13 bawaan, `PropsRute.onLihatBerkas`) |
| View Payment Status | harness `DetailPaymentStsCNP`, `GetDetailPaymentStatus_Act` (REST luar) | **nonaktif-OQ** OQ-CNP-34 |
| Policy Start Ceding | `SetEndDate_Act` | **dibangun** |
| Date of Loss / Report Date / Received Date | `CheckDateDOL_Act` / `CheckReportDate_Act` / `CheckDateReceived_Act` | **dibangun** |
| Catastrophe (ikon ubah / simpan) | `SetEditCatastrope` Edit / Save, local action `CatastrofeList` | **dibangun** (pola Claim Prop) |
| Choose Cause of Loss | harness `CauseofLoss_Harness` | **dibangun** — Choose = `GetNameCauseofLoss_Act` |
| Reporter Status | `GetReportStatus_Act` | **dibangun** |
| Adjuster / Professional ID, Consultant ID | `SetAdjsuter_act`, `SetConsultant_Act` | **dibangun** — dropdown dapat dicari (BrowseAdjusterConsultant, ID + nama) |
| Report Description | `MakeLowercase_Act` | **tidak tampil di XML** (`1=2`) |
| Zip Code | `GetAdders_Act` | **dibangun** |
| Insured Interests: Add / Delete / baris | `AddInterestListCNP_Act` / deleteRow / expand pane `InputDtlInterest` | **dibangun** (expand pane di bawah baris) |
| Deductible / Format / Type Deductible / Min-Max | `SetFormat_Act`, `SetDataDeductible_act` (Posisi TypeDeductible) | **dibangun** |
| Claim Amount: Add / Currency / angka / Delete | `AddListClaimNP_Act` Table=Claim, `SetCurrency_Act` Note=Claim, `CountClaimTNP_Act`, deleteRow + `CountClaimTNP_Act` | **dibangun** (NA `.CNPFlagOuts==1`) |
| Loss Allocation: Add / Share / Claim Amount / To XOL / Delete | `AddLossAlocation_Act`, `CountLossAllocation_act` pct / amount / CountXOL, deleteRow | **dibangun** |
| Edit XOL Allocation | harness `EditXOLAlokasi` (kata sandi tetap) | **nonaktif-OQ** OQ-CNP-04 — keputusan "sandi ke konfigurasi"; modul dilarang membaca env dan nol sandi di repo, jadi belum ada jalur sah |
| XOL Allocation Claim Amount | `AdjClaimAmount_Act` | **dibangun** |
| Summary XOL Allocation | (turunan) | **dibangun** (hanya-baca, `SusunSummaryXOL`) |
| Spreading List: Add / Share / Delete | `AddListClaimNP_Act` Spreading, `CountSpreadingCNP_Act`, deleteRow | **nonaktif-OQ** OQ-CNP-11 / OQ-CNP-21 — grid tampil hanya-baca |
| Break QS: Treaty Name / Share | `SetTreatyNameSpreading_Act`, `CountSpreading_act` | hanya-baca di XML (`RO selalu`) — tidak beraksi |
| Reinstatement (grid kasus) | — | **tidak tampil di XML** (`1=2`, OQ-CNP-12) |
| Save | `SaveDataToJClaim_Act` ; save | **dibangun** |
| Save to issue RNM | `SaveDataToOSAksep_Act` (NA `IsOutstanding=1` ‖ ProtectDOL ‖ ProtectStartDate ‖ ProtectEndDate) | **dibangun** — penomor, OS per layer × mata uang, JSON_KLAIM, outbox `outstanding-np` (produksi) |
| Print CFS | `GenerateCFS_act` | **dibangun** — data CFS; berkas PDF OQ-CNP-22 (pemberitahuan `info`) |
| Print PLA | local action `GeneratePLACNP` → `GeneratePlaCNP_Act` | **dibangun** — modal `PreviewPLA`; nomor PLA `PLATNP_SEQ` (OQ-CNP-39: tidak ada di DEV → 409 terang); berkas OQ-CNP-22 |
| Submit | finishAssignment (NA `IsCFS!='1'` ‖ `IsOutstanding!='1'`) | **dibangun** — pindah ke Input Acceptation |
| Claim History | grid SuggestList | **dibangun** (paging 5, urut menurun) |

## 3. Section `InputAcceptation` (Assignment1)

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| View Master | (tanpa aksi di section) | **nonaktif-OQ** OQ-CNP-20 |
| View Outstanding Claim | harness `OutstandingClaim` | **tidak tampil di XML** (`1=2`) |
| Close Claim | local action `CloseClaimMD` | **dibangun** |
| Close Without Payment (CWP) | local action `CloseClaimNP` (NA `IsAcceptation = 1`) | **dibangun** — pra-proses `CloseClaimNP_preAct`; tombol Yes **nonaktif-OQ** OQ-CNP-36 |
| Choose Master (IN / INEDM) | harness `ChooseMasterTNonProp` IN (tampil `IsAcceptation != 1`) / INEDM (tampil `IsAcceptation = 1`) | **dibangun** |
| View Payment Attachment | harness `ViewAttachmentNP`, `GetPayAttachmentNP_Act` | **dibangun** — baca `DOCUMENT_CLAIM` (`[penyimpangan sadar]` §8) |
| Claim History Master ID | harness `ViewHistoryMasterID_NP` | **dibangun** — `GetHistoryMasterID_NP` (view `CLAIMXOL` INVALID di DEV, OQ-CNP-28) |
| View Payment Status | (tanpa aksi di section) | **nonaktif-OQ** OQ-CNP-34 |
| Policy Start / Policy End | `SetEndDate_Act` / `CheckPeriodPolicy_Act` | **dibangun** |
| Date of Loss, Location of Loss | `CheckDateDOL_Act`, `MakeLowercase_Act` | **dibangun** |
| Choose Cause of Loss (NA `IsAcceptation = 1`) | harness `CauseofLoss_Harness` | **dibangun** |
| Adjuster / Consultant ID (tampil bila kosong) | `SetAdjsuter_act` / `SetConsultant_Act` | **dibangun** |
| Insured Interests | hanya-baca | **dibangun** (tanpa Add / Delete / pane) |
| Share Ceding(%) | `IntIsSaveToOs` | **dibangun** |
| Deductible (NA `IsAcceptation = 1`) | `SetFormat_Act` | **dibangun** (Format / Type postValue saja) |
| Waiting For Actual Premium (centang) | `SetActualPremium_ACT` | **dibangun** — nilai 50.000 konstanta bernama (OQ-CNP-17) |
| View (tampil `FlagActualPremium==true`) | harness `Hitung_Test`, `GetSelisihActual_Act` | **dibangun** — perbaikan OQ-CNP-05 butir 2 |
| View Old Allocation | harness `ViewOldAllocation`, `GetDataOldAllocation` | **dibangun** |
| Claim Amount / Loss Allocation / XOL Allocation | sama dengan Outstanding | **dibangun** |
| Spreading grid | `CountSpreadingCNP_Act`, `SetTreatyNameSpreading_Act` | hanya-baca (OQ-CNP-11 / 21) |
| Save (dua) | `Objsave` | **tidak tampil di XML** (`1=2`) |
| Save To OS | `SaveToOS` (NA `IsSaveToOs==1`) | **dibangun** |
| Acceptation List: Add / Delete / baris | `AddAkseptasiCNP_Act` (NA `IsSaveToOs==0`), `DeleteAkseptasi_Act` (NA `AcceptanceStatus!=''`), expand pane `AdjustmentDetailNP` | **dibangun** (baris terbaru terbuka, pola Claim Prop) |
| Spreading Claim Out (dua Spreading List) | hanya-baca | **dibangun** |
| Committe Accept Status | grid `ClaimComitee` | **dibangun** (tampil bila `IsCloseFile` ‖ `IsReject`) |
| Back | `BackToRegister_act` | **tidak dibangun** — container NEVER (OQ-CNP-09) |
| Submit | — | tidak ada di section (OQ-CNP-10): kasus selesai lewat Close Claim / CWP |

## 4. Panel `AdjustmentDetailNP` + `Subjectivity` (expand pane Acceptation List)

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| Payment Type (NA `IsKomite=="1"`) | `SetInterimXOL_Act` | **dibangun** (kode mentah 1–7, OQ-CNP-14) |
| No DLA : | postValue | **dibangun** |
| Claim Acceptation: Add | `AddListClaimNP_Act` (kelas Data-Adjustment) | **nonaktif-OQ** OQ-CNP-33 |
| Claim Acceptation: Rate of Exchange | `CountClaimTNP_Act` | **dibangun** (hitungan tingkat klaim, jalur absolut pyWorkPage) |
| Claim Acceptation: Currency / Claim Amount / TPL / Adjuster Fee / Salvage / Fee | `SetCurrency_Act` / `CountClaimTNP_Act` kelas Data-Adjustment | postValue saja — OQ-CNP-33 |
| Claim Acceptation: Delete | deleteRow + `CountClaimTNP_Act` | **dibangun** (`HapusKlaimAkseptasi`) |
| Loss Allocation: Add / angka / To XOL | `AddLossAlocation_Act`, `CountLossAllocation_act` kelas Data-Adjustment | **nonaktif-OQ** / postValue — OQ-CNP-33 |
| Loss Allocation: Delete | deleteRow | **dibangun** |
| XOL Allocation Claim Amount | `AdjClaimCNP_Act` | **dibangun** |
| XOL Allocation baris | expand pane `ShowDetailXOL` (`ReinstatementPremiumDetails`) | **dibangun** — rumus varian `FlagProrate` 0 / 1 (OQ-CNP-12) |
| Previously Calculated | hanya-baca | **dibangun** |
| Spreading In / Spreading Out | hanya-baca | **dibangun** |
| Payable To | `SetPayableTreatyNP_Act` Posisi=Acc + `SetAccoutNo_Act` | **dibangun** |
| Specify (teks / dropdown klien bila Payable = 3) | `SetPayableTreatyNP_Act` + `SetAccoutNo_Act` | **dibangun** |
| Name (Payable = 3) | refresh | **dibangun** |
| Name of Bank / Name of Bank 2 | `SetAccoutNo_Act` (autocomplete Result.pxResults) | **dibangun** — `[inferensi]` §8 (pilih baris rekening) |
| Currency / Swift Code / Branch / Account No (1 dan 2) | NA selalu | **dibangun** (rekening kedua tampil `FlagCurrency==1`, `[inferensi]`) |
| Transfer Direct to Kasir | postValue (NA komite ‖ `FlagErrorKasir`) | **dibangun** |
| Save | save | **dibangun** (`Simpan`) |
| Send to Committe | `SetAccoutNo_Act` + proteksi + harness `KomiteCNP` (tampil `TotalKomite!=''`, NA komite ‖ `IsSaveToOs==0`) | **dibangun** |
| Acceptation | `HitServiceToKasir_Act` | **dibangun** — outbox `kasir` (produksi); `konfigurasi/kasir.json` |
| Generate DLA | local action `GenerateDLACNP` (NA selalu di XML) | **nonaktif-OQ** OQ-CNP-19 |
| Generate Claim Analysis | `GenerateCACNP_Act` | **dibangun** — data; berkas OQ-CNP-22 |
| Save Previously Paid | `SaveCNPLayerList_Act` | **dibangun** — OS `STS_REJECT = 5` |
| Subjectivity Note | hanya-baca | **dibangun** |

## 5. Local action, harness, expand pane

| Rule | Tombol / medan | Status |
| --- | --- | --- |
| `KomiteCLMNP` (harness KomiteCNP "Send To Commitee") | Date, Initial, Circumstanses (NA `Type==3`), Remarks; **Send Claim to Committee** (NA `Payable==''` ‖ CekError ‖ `Occupation==''`) → `CreateChildKomiteCNP_Act`; Cancel | **dibangun** — kasus `KMTNP-` + tangga + `POSITION` = workbasket tingkat 1 + outbox `email-komite`; perbaikan OQ-CNP-05 butir 4 |
| `CloseClaimMD` | Remarks (wajib), No, Yes → `CloseClaimTNonProp` | **dibangun** — OS STS 4, JSON_KLAIM, outbox `tutup-klaim-np`; perbaikan butir 6 |
| `CloseClaimNP` | Date, Initial, Remarks, No, Yes → `CreateChildKomiteCloseNP_Act` | No **dibangun**; Yes **nonaktif-OQ** OQ-CNP-36 |
| `PreviewPLA` (FlowAction GeneratePLACNP, `pyShowFAButtons=false`) | Remarks (`PreviewPLA.CARI14`), Submit / Cancel | **dibangun** — tombol Submit / Cancel bawaan modal Pega `[inferensi]` |
| `InputDtlInterest` (expand pane) | Interest Insured, Value In IDR, Currency → `SetCurrency_Act` Note=Interest, Value → `CountTotalInterest_Act`, Deductible / Format, TPL → `SetTPLNote_Act`; Submit / Cancel bawaan flow action (`pyShowFAButtons=true`) | **dibangun** — `view.CARI21` tanpa penulis dianggap kosong (OQ-CNP-35) |
| `ReinstatementPremiumDetails` (expand pane ShowDetailXOL) | rumus, Cancel bawaan | **dibangun** (hanya-baca) |
| `EditXOLAlokasi` | Submit → `EditXOLAlokasi` | **nonaktif-OQ** OQ-CNP-04 (tombol pembukanya nonaktif) |
| `GenerateDLACNP` | — | **nonaktif-OQ** OQ-CNP-19 |
| `DetailPaymentStsCNP` | — | **nonaktif-OQ** OQ-CNP-34 |
| `ChooseMasterTNonProp` | Note, saringan, Choose | **dibangun** (judul pyLabel "Choose Master Treaty XOL") |
| `ViewListPolicyCNP` | Policy No (pxLink), Source of Business, Ceding Co, Begin Date, End Date | **dibangun** (OQ-CNP-24: pemetaan CARI menurut teks SQL) |
| `CauseofLoss_Harness` | Choose; **Add New Cause of Loss** → harness `TambahCauseofLoss` | Choose **dibangun**; Add New **tidak dibangun** — pemeliharaan master Cause of Loss milik menu lain (pola Claim Prop) |
| `CatastrofeList` | Add New, Choose, Save, Cancel | **dibangun** (pola Claim Prop) |
| `Hitung_Test` | Nilai Estimasi / Nilai Akseptasi / NIlai Selisih Actual Premi / 100% Reserve Updated | **dibangun** (label verbatim, termasuk "NIlai") |
| `ViewOldAllocation` | grid Old XOL Allocation | **dibangun** |
| `ViewAttachmentNP` | Kategori, Nama File, No Akseptasi, No Prekas, Tanggal/Waktu | **dibangun** |
| `ViewHistoryMasterID_NP` | label `InputSpreading.CARI40`, grid klaim + grid total | **dibangun** — judul kolom keempat "Loss to Layer (100%)" dua kali verbatim (isinya CNPReinstatement; kelainan XML dipertahankan) |

## 6. Tidak tampil di XML (tidak dibangun)

Choose Master Out, CWP di Outstanding, View Outstanding Claim, Report Description, grid Limit Layer + `ViewClaimLayerDetail`
/ `CoBList` (S11–S13 `1=2`), grid Reinstatement kasus, tombol "Button" (`1=2`) di Outstanding / Input Acceptation /
InputDtlInterest, dua tombol Save Input Acceptation, Back (container NEVER).

## 7. Tujuh perbaikan OQ-CNP-05 — `[penyimpangan sadar]`

| # | Rule XML | Perilaku Pega lama | Perilaku baru | Uji |
| --- | --- | --- | --- | --- |
| 1 | `GenerateCFS_act` | penghapus list ber-REMARK → agregasi CFS menjumlah ke list berisi; Fee ↔ Claim Amount tertukar di Loss Allocation | list dikosongkan dulu; kolom sesuai nama | `models.TestCFSAgregasiTanpaGandaDanFeeTidakTertukar` |
| 2 | `GetSelisihActual_Act` | membaca `OutOSAcc.pxResults(1).CARI1/2`, padahal SQL `GetDataOS` beralias Value / GrossValue → popup kosong | membaca alias Value / GrossValue | `models.TestSelisihAktualMembacaAliasGetDataOS` |
| 3 | `ProteksiSendKomiteCNP_Act` | proteksi "Error No Account" tak pernah aktif; `IsError = 2` tak pernah direset | rekening salah ditolak (maksud diturunkan dari 7.3 / 7.4, konfirmasi OQ-CNP-44); `IsError` direset tiap pemeriksaan | `models.TestProteksiKomiteRekeningSalahDanIsErrorDireset` |
| 4 | `CreateChildKomiteCNP_Act` | pesan "Nilai Gross Value tidak sesuai dengan Spreading In" muncul SESUDAH kasus komite dibuat | divalidasi SEBELUM; gagal = 422, nol tulisan | `models.TestValidasiGrossSpreadingIn`, `handlers.TestAkseptasiSampaiKasusKomite` |
| 5 | `HitServiceToKasir_Act` | tahun Tgl Boleh Bayar naik menurut bulan SEKARANG | menurut bulan akseptasi | `models.TestTanggalBolehBayarTahunBergulir` |
| 6 | `CloseClaimTNonProp` | kronologi ditulis ke `CARI12`, DT membaca `CARI1` | kronologi tercatat | `models.TestTutupKlaimMencatatKronologi` |
| 7 | `CountReinstatement_Act` | Adjuster Fee porsi RNM, Salvage digross-up | tidak berlaku: grid Reinstatement kasus `1=2` tidak dibangun (OQ-CNP-12); rumus akseptasi (`AdjClaimCNP_Act` s.11) tidak memuat kelainan ini | — |

## 8. Penyimpangan sadar lain dan `[inferensi]`

- **`[penyimpangan sadar]` ViewAttachmentNP**: Report Definition `GCNMGetInvoiceAttachments` membaca Link-Attachment,
  tetapi pembuat Link-Attachment ber-REMARK (`GetPayAttachmentAdj_Act` 3.3.3–3.3.6); yang hidup menulis tabel warisan
  `DOCUMENT_CLAIM` (3.3.10) → dibaca dari sana (KATEGORI_1 Invoice, NOAKSEP terisi).
- **`[penyimpangan sadar]` alamat agen** dibaca dari tabel datar `CLIENT_ADDRESS` (pola Claim Prop `AlamatKlien`).
- **`[penyimpangan sadar]` popup master**: saringan per kolom + batas 500 + 50 per halaman, menggantikan isian "Input
  Treaty ID" + Search (pola Claim Prop, prompt §8 P2).
- **`[inferensi]` TreatyName** saat Choose master = `.TREATYCONTRACTNAME` baris pop-up (langkah 2 membaca
  `Param.TreatyName` yang tidak dikirim tombol Choose).
- **`[inferensi]` label kode** hanya untuk properti sekelas dengan label yang diberikan work owner di Claim Prop
  (ReportType, ReporterStatus, Payable, Adjustment.Type, Comitee.KomiteAproval); kode lain tampil mentah (OQ-CNP-14).
- **`[inferensi]` FlagCurrency** akseptasi = mata uang akseptasi lebih dari satu (tanpa penulis di korpus); menyalakan
  rekening kedua.
- **`[penyimpangan sadar]` TotalKomite** dihitung saat layar disusun (calon tangga), padahal XML hanya menulisnya
  SESUDAH penyerahan — persis XML tombol Send to Committe tidak pernah tampil (OQ-CNP-42, konfirmasi work owner).
- **`[penyimpangan sadar]` Save to issue RNM langkah 8**: `BusinessOldId` dibaca (GetDataBusiness menurut nama bisnis)
  SEBELUM pemeriksaan langkah 8; XML mengisinya baru di 15.2–15.3 sehingga persis XML tombol itu tidak pernah lolos untuk
  kasus baru (OQ-CNP-43).
- **Kasir LdcId syariah** (When `IsPEGASyariah` = pemeriksaan node server Pega): selalu konvensional (OQ-CNP-40).
- **Roster tingkat 1** = RD `FilterEmailKomiteWithLimit` `LIMIT_BOTTOM <= 0` persis (baris ber-LIMIT_BOTTOM NULL tidak
  lolos); DEV DEGREE 1 = -9.999.999.999.999.
- **`[inferensi]` SetAccoutNo_Act** disederhanakan: memilih satu baris rekening dari `Result.pxResults` (nilai
  `AccountNo|NameOfBank`), server membaca ulang barisnya.
- **`[inferensi]` PreviewPLA** Submit / Cancel bawaan modal; `Remark_Close` disimpan dari `Message` CloseClaimMD.
- **`[inferensi]` pembagi nol** `@divide` / `/` bernilai 0 (`models/hitung.go`).
- Properti tanpa penulis di korpus dibiarkan kosong persis XML (OQ-CNP-37): `OfferFacIn.QuotationData.BusinessOldId`
  (nomor PLA, LbuID Kasir), `ClaimData.QuotationData.BusinessName`, `ValueAdjustment` (tangga = 0),
  `FlagProrate` (varian 0), `TotalListClaimAmount(IDR)`, `TotalEstimasi`, `ProtectEndDate`, `IsTreatyIn`.
- `OS_AKSEPTASI_KLAIM.DATA_JSON` ditulis menurut halaman XML; baris DEV STS 0 memuat kunci tambahan (EstimationDate,
  PolicyNo, pzInsKey) yang tidak ditulis rule ekspor (OQ-CNP-38).

## 9. Frontend

- **`[penyimpangan sadar]` tata letak ikut Claim Prop** (perintah work owner 09-10-2026 "ikuti tampilan klaim prop"):
  Claim Information menjadi kartu sendiri di atas layout group (di XML tab pertama), tab tinggal Interests / Estimation /
  Spreading (Outstanding) dan Interests / Estimation / Acceptation (Input Acceptation); tombol layar satu baris aksi tanpa
  kartu. Isi medan, kondisi tampil / hanya-baca / nonaktif, dan aksinya tidak berubah.
- **Lampiran klaim ikut Claim Prop** (perintah work owner 09-10-2026 "untuk attachment juga mengikuti dari klaim prop"):
  tab Lampiran (Add attachment / Refresh / Save, Category - Count Attach - Upload File - View File, jendela View File
  dengan View / View Office Online / Delete / Change Category), master `T_KATEGORI_DOC_KLAIM` TYPE_KLAIM NONPROP, tabel
  warisan dokumen klaim + penyimpanan bersama inti. Korpus Non Prop sendiri tanpa section unggahan; Save tab = Save
  Outstanding (`SaveDataToJClaim`).

- Renderer tata, kulit, isian angka / tanggal disalin dari `modul/claimprop/frontend` (bukan impor); kelas berawalan
  `claimnonprop__`, token `--cnp-*`.
- Expand pane umum (`components/rincian.ts`): Acceptation List (nomor akseptasi untuk aksi panel), Insured Interests, XOL
  Allocation panel; Cancel / Submit flow action menutup pane.
- Aksi panel dikirim `{indeks: nomor akseptasi, baris: baris grid panel}` (`TataView.alamatAksi`); server membaca gerbang
  sel pada `baris` untuk semua aksi.
- Baris tombol `sebaris` di akhir modal pindah ke kaki Modal (tanpa Cancel ganda).
- Popup diberi `key` per jenis — temuan cek layar: berganti jenis tanpa dilepas merender data jenis lama (galat
  `tampilTanggal` pada data master).

## 10. Isolasi kotak masuk (prompt §6 butir 1)

Kasus Claim Non Prop: `T_WORK_CLAIM.LINI = 'NONPROP'`, `TAHAP` selalu terisi (`OutstandingClaim` / `InputAcceptation` /
`KomiteTreaty_Flow`; penutupan hanya mengisi `STATUS_WORK`, `repository/gudang.go` `sqlTutupKasus`).

| Kotak masuk | Saringan (dibaca dari SQL modul itu, tidak disunting) | Baris NONPROP |
| --- | --- | --- |
| Claim Life | `NVL(w.TAHAP, :tahapCadangan) = :tahap` dengan tahap "Input Register" / "Outstanding Claim" / "Medical Check" / "Claim Analis" (`claimlife/backend/repository/inbox.go` `sqlInboxWhere`, `models/tahap.go`) | tidak cocok: TAHAP NONPROP tanpa spasi dan tidak pernah NULL |
| Komite Claim Life | `AND (w.LINI = :lini OR w.LINI IS NULL)` (`komiteclaimlife/backend/repository/komite_inbox.go`) | tersaring (`NONPROP` ≠ Life, tidak NULL) |
| Claim Prop | `w.LINI = 'PROP' AND w.ID LIKE 'CLMP-%'` (`claimprop/backend/repository/gudang.go`) | tersaring |
| Komite Claim Prop | `w.LINI = 'PROP'` ketat + awalan `TKMT-` (`komiteclaimprop/backend/repository/gudang.go`, `tangga.go`) | tersaring |
