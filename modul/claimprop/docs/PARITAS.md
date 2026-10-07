# PARITAS — tombol dan aksi XML Claim Prop ↔ tiket ↔ status

> 07-10-2026. Sumber inventaris: sapuan 36 Section, 11 Harness, 12 FlowAction, 1 Flow di `D:\XML\RNM_BRD\Claim Prop`
> (103 kontrol pxButton / pxLink / pxIcon unik, aksi baris grid, ikon grid standar). Sumber status: definisi layar
> server `backend/models/layar.go` dan `layar_adjustment.go` — satu-satunya tempat tombol lahir; frontend merender
> pohon tata apa adanya.
>
> **Status:**
> - **dibangun** — tampil menurut kondisi XML yang sama (visible-when / disabled-when), aksinya port activity-nya.
> - **nonaktif-OQ** — tampil sesuai section, tetapi nonaktif; keterangan `OQ-CP-nn` di atribut `title` (lihat `OQ.md`).
> - **tidak tampil di XML** — kondisi tampilnya statis `NEVER` / `1=2` / `never`, atau tanpa aksi: tidak dibangun.
> - **di luar lingkup** — pemeliharaan master milik menu lain (prompt §2).
>
> Cek layar 07-10-2026 (server tiruan `UJI-`, vite port 5177, Chrome tanpa kepala): halaman awal (New + tiga tab),
> layar Outstanding Claim kasus baru (54 medan, nol galat konsol), popup "Data Master TreatyIn" — tombol tampak
> cocok dengan baris berstatus di bawah untuk keadaan kasus baru. Tahap Input Acceptation sampai Resolved-Completed
> dijaga uji HTTP `TestAlurPenuhSampaiResolved`.

## 1. Flow dan halaman awal

| XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| `Flow_TreatyIn` Start1 → Assignment2 "Outstanding Claim" (worklist pembuat) | pembuatan kasus | 01 | **dibangun** — tombol **New** + tab *Outstanding Claim* (OQ-CP-13: harness `New` tidak diekspor) |
| Assignment1 "Input Acceptation" (workbasket) | FlowAction `InputAcceptation` | 10 | **dibangun** — tab *Input Acceptation*, workbasket `ReasKlaimTeknik` (keputusan 07-10-2026) |
| Decision3 `IsBackStage` → Assignment2 / End1 Resolved-Completed | — | 11 | **dibangun** — tab *Resolved-Completed*; jalur kembali lewat "Back" yang di XML `NEVER` → tidak dibangun |
| FlowAction `OutstandingClaim` pre-act `CheeckNoRNM_Act`, pre-DT `SetDateOutstanding`, post-act `ProteksiData_act` | — | 01, 02 | **dibangun** (`siapkan`; tanggal hanya diisi bila kosong — penyimpangan, OQ.md) |
| FlowAction `InputAcceptation` pre-act `GetPICAdjutment_Act`, pre-DT `SetDateAcceptation` | — | 08 | **dibangun** |

## 2. Section `OutstandingClaim` (Assignment2)

| Label XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| Choose Master | harness `MasterTreatyIn` (dis `IsOutstanding==1`) | 01 | **dibangun** (popup, tombol Choose per baris = `SetValueToClaim_Act`) |
| Summary Outstanding Claim | harness `SummaryOutSClaim`, pre-act `GetDataOustanding` (vis `IsOutstanding = 1`) | 07 | **dibangun** |
| teks `.Message` (klik) | harness `InputTreatyInOffer` (tidak diekspor) | 01 | **nonaktif-OQ** OQ-CP-02 — tampil hanya-baca |
| Choose Policy No | harness `ListPolicyNoTreaty_Harness`, pre-act `SetMasterID` | 01 | **dibangun** (pilih = klik ganda baris → `CheckNoPolicy`, seperti XML) |
| View (polis) | harness `DetailPolisRealization` (tidak diekspor) | 01 | **nonaktif-OQ** OQ-CP-02 |
| View Status Payment Premi | harness `ListPaymentPremi_Harness`, pre-act `GetDtlPaymentPremi_act` | 02 | **nonaktif-OQ** OQ-CP-03 |
| Choose Cause of Loss | harness `CauseofLoss_Harness` (dis `IsAnyAcceptation =1`) | 03 | **dibangun** |
| ikon tambah Consultant / Appointed Adjuster | harness `MstAdjusterConsultant` | 08 | **di luar lingkup** — nonaktif OQ-CP-04 |
| ikon edit RNM Share % | `GetRNMShareTreaty` (vis `EstimationList(1).PrintFaceClaim = ''`) | 06 | **dibangun** |
| ikon "Standard Icon" RNM Share % | tanpa aksi, vis never | — | tidak tampil di XML |
| Save | `SetOutstanding_Act` lalu save (vis `IsAcceptation!=1`) | 07 | **dibangun** |
| Save to issue RNM | `SaveOutstanding_Act` → LA `PrintFile` bila `IsPLA=1` (dis `IsCFS = ''`) | 07, 12 | **dibangun** (PrintFile = pesan "Please Print Pla") |
| PRINT PLA | LA `GeneratePLA` (dis `IsPLA!=1`) | 12 | **dibangun** — nomor PLA; berkas PDF OQ-CP-05 |
| Send to Acceptation | `CheckNopolicy_Act` (vis `IsOutstanding = 1`, dis `IsAcceptation==1`) | 01 | **dibangun** |
| Submit | finishAssignment `OutstandingClaim` bila PolicyNo dan NoClaim terisi (vis `IsAcceptation==1`) | 10 | **dibangun** — pindah ke workbasket `ReasKlaimTeknik` |

### Sub-section `Catastrope_Sec`, `OutstandingClaim_Intrs`, `OutstandingClaim_Est`, `OutstandingClaim_Sprd`

| Kontrol XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| ikon pensil / centang Catastrophe | `SetEditCatastrope` Edit / Save | 03 | **dibangun** |
| ikon gear Catastrophe | LA `CatastrofeList` | 03 | **dibangun** (popup: tambah = `InputCatastrope`, Save = `SaveCatasrtope_Act`, pilih = `SetCatastrope_act`, Cancel) |
| tambah / hapus baris Interest | `AddInterest_act` / `DeleteInterest_act` (dis `IsAdjVal=='Yes'`) | 05 | **dibangun** |
| tambah / hapus baris Total Claim Amount | `AddListClaimAmount` / `DeleteListClaim_Act` (dis `Note=='Yes'`) | 09 | **dibangun** |
| tambah / hapus baris Loss Allocation | `AddLossAllocation_act` / `RemoveLossAlloction_act` (dis `IsOldData=='Yes'`) | 06 | **dibangun** |
| tambah / hapus baris Estimation | `AddEstimation_Act` / `DeleteEstimation_Act` (dis `PrintFaceClaim==1`) | 07 | **dibangun** |
| tambah baris Spreading Claim | `AddSpreading_Act` (**tidak diekspor**) | 06 | **nonaktif-OQ** OQ-CP-01 |
| hapus baris Spreading Claim (ikon baris) | `DeleteSpreading_Act` (**tidak diekspor**, dis `IsOldData='Yes'`) | 06 | **nonaktif-OQ** OQ-CP-01 |
| ikon grid standar Spreading List (tambah / hapus bawaan) | `pzPegaDefaultGridIcons` ("always"), tanpa activity | 06 | **tidak dibangun** — tambah / hapus spreading dinonaktifkan (keputusan work owner 07-10-2026) |
| tombol tanpa label di samping RNMShareP | vis NEVER | — | tidak tampil di XML |
| grid "Claim History" (SuggestList) | — | 14 | **dibangun** — `T_VIEW_SUGGEST` (keputusan work owner 07-10-2026) |

## 3. Section `InputAcceptation` (Assignment1) dan sub-section-nya

| Label XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| View Status Payment Premi (#1) | aksi kosong, vis NEVER | — | tidak tampil di XML |
| Summary Outstanding Claim | harness `SummaryOutSClaim`, pre-act `GetDataOutsClaim_act` | 07 | **dibangun** |
| View Status Payment Premi (#2, #3) | harness `ListPaymentPremi_Harness` | 02 | **nonaktif-OQ** OQ-CP-03 |
| Close Claim | LA `PreventRejectClaimProp` | 11 | **dibangun** |
| Choose Policy No | vis NEVER | — | tidak tampil di XML |
| View (polis) | harness `DetailPolisRealization` (tidak diekspor) | 01 | **nonaktif-OQ** OQ-CP-02 |
| Choose Cause of Loss | harness `CauseofLoss_Harness` | 03 | **dibangun** |
| ikon tambah Consultant / Appointed Adjuster | harness `MstAdjusterConsultant` | 08 | **di luar lingkup** — nonaktif OQ-CP-04 |
| Save (di bawah Est) / Save (di bawah Adjs) | save | 07, 08 | **dibangun** |
| Back | `BackToRegister_act` (tidak diekspor) — layout NEVER | — | tidak tampil di XML — prompt §6 butir 9 menyebutnya, tetapi "tampil sesuai section" = tidak tampil; activity-nya tetap OQ-CP-01 |
| Update Estimation / Save (`InputAcceptation_Adjs`) | `UpdateEstimasi_Act` (tidak diekspor) — layout NEVER | — | tidak tampil di XML — idem; OQ-CP-01 |
| View Status Payment Claim | harness `ListPaymentClaim_Harness` | 13 | **nonaktif-OQ** OQ-CP-03 |
| View Payment Attachment | harness `ViewAttachment` | 13 | **nonaktif-OQ** OQ-CP-03 |
| tambah baris Adjustment | `AddAdjustment_Act` (vis `AktifButton = '0' \|\| ''`) | 08 | **dibangun** |
| hapus baris Adjustment | `DeleteAjsutment_Act` (dis `IsKomite=1 \|\| IsSubjectivity`) | 08 | **dibangun** (+ tolak baris ber-`KOMITE_ID`) |
| Save to issue RNM / PRINT PLA (`InputAcceptation_Est`) | sama dengan Outstanding (dis `ReCFS != 1` / `IsPLA!=1`) | 07, 12 | **dibangun** |
| klik baris Adjustment | expand pane `AdjustmentDetail` | 08 | **dibangun** (panel rinci baris) |
| klik baris Total Claim Amount | expand pane `ViewDetailInterest` (`InterestListDtl` — tidak diimpor, J1) | 05 | tidak dibangun: halaman kerja bukan data (STRUKTUR J1) |

## 4. Section `AdjustmentDetail` (panel rinci baris adjustment)

| Label XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| Save | save (layout `view.CARI21!=1`) | 08 | **dibangun** |
| Send to Committe | `AttachmentProtect_ACT` → harness `CommitteeTreaty` bila `Protect.CARI1/2 = 1` (vis `TotalKomite!=''`, dis `IsKomite==1 \|\| IsError>1`) | 11 | **dibangun** — gerbang lampiran selalu menolak (OQ-CP-12) |
| Acceptation | `SaveAcceptationTreaty_Act` (when NEVER — tidak pernah jalan), `HitServiceToKasir_Act`, LA `PrintFileDLA` (vis `AcceptanceStatus = 1`) | 13 | **dibangun** — di luar produksi activity keluar di langkah 3 (status konversi hanya di produksi) |
| Generate DLA | LA `GenerateDLATreaty` (dis `IsFacRetro != 1 \|\| DLA_No != ''`) | 12 | **dibangun** — nomor DLA; berkas OQ-CP-05 |
| Print Claim Analysis | `CreatClaimAnalysis_Act` | 12 | **nonaktif-OQ** OQ-CP-05 |
| View Komite No | `SetKomiteNo_Act` (vis `KomiteNo==''`, layout `IsKomite == 1`) | 11 | **dibangun** |
| centang Direct To Kasir | postValue | 13 | **dibangun** |
| grid Spreading In / Spreading Out / Loss Allocation | ikon grid standar | 06 | **dibangun** (hanya-baca seperti XML) |
| grid "Committe Accept Status" (`.ComiteeClaim`) | — | 11 | **dibangun** — roster calon, lalu keputusan tangga (Committee Name / Status / Date Approve / Comment) |

## 5. Harness `CommitteeTreaty` (`ComiteeClaimTreaty`)

| Label XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| Send Claim to Committee | `AddKomiteTreatyChild_ACT` (dis `Payable = ''`) | 11 | **dibangun** (opsi "b" 07-10-2026) — kasus TKMT- + tangga + email komite; terjangkau sesudah gerbang lampiran (OQ-CP-12) |
| Cancel | tutup | 11 | **dibangun** |

## 6. Local action modal

| FlowAction | Isi | Tiket | Status |
| --- | --- | --- | --- |
| `PrintFile` / `PrintFileDLA` | label "Please Print Pla" / "Please Print DLA" | 12 | **dibangun** sebagai pesan info |
| `GeneratePLA` | `.ClaimData.Remark`, post-act `TryMakePLA_Act`; tombol `pySubmitLabel` "Submit" / `pyCancelLabel` "Cancel" | 12 | **dibangun** |
| `GenerateDLATreaty` | `.RemarksDLA`, post-act `PrintDLATreatyIn`; Submit / Cancel | 12 | **dibangun** |
| `PreventRejectClaimProp` | Yes (tanpa bayar → `SendCloseClaimToKomite`) / No / Yes (`CloseClaimProp`) / No / Close | 11 | **dibangun**, kecuali Yes tanpa bayar = **nonaktif-OQ** OQ-CP-06 |
| `CatastrofeList` | lihat §2 | 03 | **dibangun** |
| `MessageBeforeDeleteTreatyGroup` | dibuka tombol Delete yang `1==2` | — | tidak tampil di XML |

## 6a. Lampiran dan efek keluar (bukan tombol section, tetapi aksi XML)

| Rule XML | Dipicu | Tiket | Status |
| --- | --- | --- | --- |
| section `pyCaseAttachmentsWrapper` (direfresh sesudah PRINT PLA) | lampiran kasus | 11, 12 | **tidak dibangun** — section tidak diekspor; lampiran OQ-CP-12 |
| REST `SendAcceptationToKasir` | Acceptation (`HitServiceToKasir_Act`) | 13 | **dibangun** — efek outbox "kasir" hanya bila `IS_PEGA_PROD`; panggilan nyata OQ-CP-03 |
| REST `KonversiKlaimNonLife` | `KonversiKlaim_Act` (Save to issue RNM, Close Claim) | 13 | **dibangun** — efek outbox "konversi-klaim" hanya bila `IS_PEGA_PROD`; OQ-CP-03 |
| REST `GetDtlPaymentClaim` | View Status Payment Claim | 13 | **nonaktif-OQ** OQ-CP-03 |
| REST `getPayAttachment` | View Payment Attachment (`GetPayAttachment_Act` → `GetPayAttachmentAdj_Act`) | 13 | **nonaktif-OQ** OQ-CP-03 |
| REST `getPremiumPaidOnTreatyIn` | View Status Payment Premi | 02 | **nonaktif-OQ** OQ-CP-03 |
| REST `ServiceGoogle` | `InsertGoogleStorage_Act` / `GetUrlGoogleStorage_Act` (berkas) | 11, 12 | **tidak dibangun** — OQ-CP-12 / OQ-CP-05 |
| email `SendEmailKlaim` | `AddKomiteTreatyChild_ACT` 34 | 11, 13 | **dibangun** — efek outbox "email-komite" hanya bila `IS_PEGA_PROD`; CC/BCC OQ-CP-15 |
| email `SendEmailKlaimRejectClose` | `SendCloseClaimToKomite` (tutup tanpa bayar) | 11 | **tidak dibangun** — OQ-CP-06 |

## 7. Pemeliharaan master (di luar lingkup, prompt §2)

`BrowseCauseOfLoss` (Change, Save), `BrowseDetailCauseOfLoss` (Add Master, Save, Change), `CauseofLoss_Section`
"Add New Cause of Loss", `GridCauseOfLoss`, `ListDetailCauseOfLoss`, `MstAdjusterConsultant` (Save, Cancel, Add, Edit,
Delete), `IsDeleteTreatyGroup` (Yes, Delete) — **di luar lingkup**, OQ-CP-04. Pemilih yang dipakai Claim Prop
(Choose Cause of Loss, Choose pada popup) dibangun.

## 8. Rule yang dirujuk tetapi tidak diekspor

Activity `BackToRegister_act`, `UpdateEstimasi_Act`, `AddSpreading_Act`, `DeleteSpreading_Act`; data transform
`BackFromRegister`; harness `DetailPolisRealization`, `InputTreatyInOffer`; section `CaseContent_Bottom`,
`InputRegisterDetail`, `ShowObjectAdj`, `ItemListEstimation`, `pyCaseAttachmentsWrapper`, `InputTreatyGroup`,
`InputDtlTreatyGroup`, `Adjusment_SC`, `InputAcceptation_Sprd`; flow action `AddFacOfferList`. Yang dipanggil kontrol
yang tampil → nonaktif-OQ (OQ-CP-01 / 02); yang dipanggil kontrol `NEVER` → tidak dibangun.

## 9. Ringkasan

Dihitung dari 73 baris tabel §1–§6a (baris berstatus campuran dihitung menurut status utamanya):

| Status | Cacah baris |
| --- | --- |
| dibangun | 45 |
| nonaktif-OQ | 13 |
| tidak tampil di XML / bukan data / tidak dibangun (OQ atau keputusan) | 13 |
| di luar lingkup | 2 baris + seluruh §7 |
