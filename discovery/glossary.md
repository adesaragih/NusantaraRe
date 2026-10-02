# Glossary — Istilah Domain Nusantara Re (FASE A, final)

**STEP D3 selesai.** Disusun 2026-09-14 dengan mengonsolidasikan `glossary-draft.md` (diisi
sepanjang D2 Tahap 1–5) dan seluruh `flows/*.md` + `modules/*.md`.

> **Kandidat seed `CONTEXT.md` untuk FASE B.** Dokumen ini dimaksudkan menjadi bahan dasar
> `CONTEXT.md` yang dipakai FASE B. Ia **bukan** kamus resmi bisnis: sebagian besar entri masih
> `arti belum terverifikasi` dan menunggu jawaban pemilik peran (lihat `open-questions.md`).
> Menyalinnya ke FASE B **tanpa menjawab OQ-020, OQ-008, OQ-014, OQ-057** akan memindahkan
> ketidakpastian, bukan menghilangkannya.

## Aturan pengisian (mengikat, dipertahankan dari draft)

1. Istilah dicatat **apa adanya** dari korpus — Indonesia tetap Indonesia, Inggris tetap Inggris.
   **Jangan mengarang terjemahan.**
2. **Jangan mengarang kepanjangan singkatan.** Yang tidak dijabarkan korpus ditulis
   `kepanjangan belum terverifikasi`.
3. **Jangan mengambil arti field dari caption di sebelahnya.**
4. Setiap entri wajib mengisi kolom **Bukti** (`path + rule`).
5. Istilah yang maknanya tidak pasti tetap dicatat, dengan `arti belum terverifikasi`, dan
   didaftarkan ke `open-questions.md`.

## Cakupan terukur

| Bagian | Entri |
| --- | ---: |
| 1. Istilah domain | **52** |
| 2. Kode dan enumerasi | **51** |
| 3. Nama workbasket / antrean | **25** |
| 4. Singkatan belum terjabarkan | **41** |
| **Total** | **169** |

Perintah audit:
```
awk '/^## /{sec=$0} /^\| /{ if ($0 !~ /^\| *:?-+/ && $0 !~ /^\| (Istilah|Kode|Nama|Singkatan) \|/) n[sec]++ }
     END{for(s in n) print n[s]"\t"s}' glossary.md | sort -rn
```

`[terverifikasi]` Dari 174 entri, **yang artinya terverifikasi penuh adalah minoritas** — mayoritas
mekanismenya terbaca tetapi arti bisnisnya tidak. Itu bukan kekurangan telusur, melainkan
**batas korpus**: kode dan singkatan tidak dijelaskan di mana pun dalam 9.369 file.

---

## 1. Istilah domain

Diurutkan sesuai urutan penemuan (D2 Tahap 1 → Tahap 5).

| Istilah | Muncul sebagai | Definisi / catatan | Label | Bukti |
| --- | --- | --- | --- | --- |
| **EDM** | nama modul, sufiks rule, nilai `param.isFOR`, properti `EdmType`, `EDMStatus` | Dipakai bergantian dengan **Addendum** dan **Endorsement** untuk konsep yang sama: perubahan atas polis yang sudah ada, dibedakan dari penutupan baru (`POLICY`). **Kepanjangan belum terverifikasi** | `[dugaan]` | `EDM Treaty In/Flow/InputAddendumTreatyIn.xml` (class `...WORK-ENDORSEMENTTREATY`); `NB Treaty In/Activity/SaveJsonPolisTreatyIn_Act.xml` (`param.isFOR` = `"EDM"` / `"POLICY"`) |
| **Addendum** | nama Flow, sufiks FlowAction | Proses perubahan polis treaty inward | `[terverifikasi]` | `EDM Treaty In/Flow/InputAddendumTreatyIn.xml`, `FlowAction/InboxPolicyTreatyInAddendum.xml` |
| **Akseptasi** | nama Activity, objek DB | Tahap persetujuan; muncul sebagai `Akseptasi_Act`, `OS_AKSEPTASI_KLAIM`, `HISTORYAKSEPTASIPEGA` | `[dugaan]` | `../inventory/_summary.md` §16.3; `Treaty In/Activity/Akseptasi_Act.xml` |
| **Realisasi / Realization** | nama Flow, label shape | Tahap setelah penawaran; `Input Realitation` (ejaan apa adanya di korpus) | `[terverifikasi]` | `NB Treaty In/Flow/InputRealizationTreatyIn.xml`, label shape `Assignment2` |
| **Quotation** | class, properti | `ASM-FW-GISFW-Data-Quotation`; `.Quotation.BusinessCode` | `[terverifikasi]` | `NB Treaty In/Activity/SaveJsonPolisTreatyIn_Act.xml` |
| **PolicyTreatyIn** | class, halaman kerja | Entitas polis treaty inward | `[terverifikasi]` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` (banyak rule) |
| **Nopolis / PolicyNo** | properti, kolom DB | Nomor polis. Format dibentuk `RNM-<CARI20>.T<BusinessOldId>.<MM.yyyy>.<5 digit urut>` | `[terverifikasi]` | `NB Treaty In/RDBList/GenerateNoPolicy.xml` |
| **ProdKe** | kolom DB, RequestType | Penanda urutan produksi; dipakai `ORDER BY PRODKE DESC` untuk mengambil versi terbaru | `[dugaan]` | `Endorsement Life/RDBList/GetEdmTypeLife.xml`; `NB Treaty In/RDBList/TreatyInSearchProdKe` |
| **Tanggal Closing / TglProd** | Local variable, RequestType | Tanggal tutup buku; default **25** bila query tidak mengembalikan nilai | `[dugaan]` | `EDM Treaty In/Activity/SaveJsonPolisTreatyInEDM_Act.xml`; `RDBList` `GETTanggalClosing_SQL` |
| **Premium List (PL)** | nama modul, layar, properti | `PL_NUMBER`, `PL_NUMBER_EDM`; layar Detail dan Summary | `[terverifikasi]` | `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` |
| **Spreading** | nama Activity, tabel | mis. `BreakDownSpreading_Act`, `CheckSpreadingProtect_ACT`, `SELECTSPREADINGTREATYINPRODUCTION` | `[dugaan]` | `NB Treaty In/Activity/BreakDownSpreading_Act.xml` |
| **Retro** | sufiks rule, nama modul | mis. `SetTreatyinRetro_Act`, `OfferFacRetro`, `Master Contract Retro Life` | `[dugaan]` | `Treaty In/Activity/SetTreatyinRetro_Act.xml` |
| **SOB** | sufiks rule/section | mis. `BusinessAndSOBList`, `RealizationPopulateCedingSOB`. **Kepanjangan belum terverifikasi** | `[pertanyaan terbuka]` | `NB Treaty In/Section/BusinessAndSOBList.xml` |
| **XOL** | nilai `.PolicyTreatyIn.ClaimType`, nama rule | mis. `InsertToTreatyXOLList`, `XOL2_AKSEP_KLAIM`. **Kepanjangan belum terverifikasi** | `[pertanyaan terbuka]` | `NB Treaty In/When/isClaimTreaty.xml` |
| **Ceding / Cedant** | nama rule, procedure | mis. `CedingCedant`, `INSERTUPDATECEDINGPRODUCTION`, `RealizationPopulateCedingSOB` | `[dugaan]` | `Claim Fac In/Section/CedingCedant.xml` |
| **Asuradur** | procedure | `GENERAL.F_GET_NM_ASURADUR` — istilah Indonesia untuk penanggung | `[dugaan]` | `../inventory/_summary.md` §22.7 |
| **Konversi** | properti, rule | `FlagErrorKonversi`, `CekSTSKonversiJson`, `PEGA_DELETE_ERROR_KONVERSI` — proses mengubah data Pega ke sistem produksi | `[dugaan]` | `EDM Treaty In/Activity/serviceInsertArasapas_act.xml` |
| **Kasir** | nama service | `SendAcceptationToKasir`, `DIRECTTOKASIR_LOG` — istilah Indonesia untuk fungsi pembayaran | `[dugaan]` | `../inventory/_summary.md` §17 |
| **Komite** | nama modul, class, properti | Badan/tahap persetujuan berjenjang. Tangga dimodelkan sebagai **loop satu assignment** yang dikendalikan data, bukan rangkaian shape | `[terverifikasi]` | `Komite Claim Life/Flow/KomiteLife_Flow.xml` |
| **KomiteList** | page list | Jejak persetujuan per tingkat: `KomiteAproval`, `KomiteComment`, `DateApprove`, `KomiteID` | `[terverifikasi]` | `Komite Claim Life/Activity/KomitePostAdjustment.xml` |
| **KomiteCount / KomiteLoop** | properti | Pencacah tingkat berjalan / batas jumlah tingkat. FacIn memakai `Local.TotalKomite` untuk konsep yang sama | `[terverifikasi]` | `Komite Claim Life/When/IsKomiteLoop.xml` |
| **KomiteRouter** | Activity | Penentu sasaran penugasan tiap tingkat (`param.AssignTo`) | `[terverifikasi]` | `Komite Claim Life/Activity/KomiteRouter.xml` |
| **Subjectivity** | properti, Activity | `IsSubjectivity`, `InsertOSSubjectivityCNP`, `ResetSubjectivityNote` | `[dugaan]` | `Komite Claim Non Prop/Activity/InsertOSSubjectivityCNP.xml` |
| **OS (Outstanding)** | prefix tabel/Activity | `OS_AKSEPTASI_KLAIM_LIFE`, `InsertOSKlaimCNP`, `SaveOSClaim_SQL`. **Kepanjangan belum terverifikasi** | `[pertanyaan terbuka]` | `Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` |
| **EMAILKOMITE** | class Pega | Sumber roster anggota komite (`Obj-Browse` → halaman `GetKomite`) | `[terverifikasi]` | `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` |
| **Register (klaim)** | tahap flow | Tahap pertama siklus kerugian; label shape `Input Register` | `[terverifikasi]` | `Claim Life/Flow/Register_Flow.xml` |
| **Estimasi** | tahap flow, Activity | Tahap penaksiran kerugian; `Input Estimasi`, `CountEstimation_Act` | `[terverifikasi]` | `Claim Fac In/Flow/Register_Flow.xml` |
| **Surveyor** | tahap flow | `Choose Surveyor` — penunjukan surveyor | `[terverifikasi]` | `Claim Fac In/Flow/Register_Flow.xml` |
| **Medical Check** | tahap flow | Hanya di Claim Life | `[terverifikasi]` | `Claim Life/Flow/Register_Flow.xml` |
| **Outstanding Claim** | tahap flow, tabel | `OS_AKSEPTASI_KLAIM`, `OS_AKSEPTASI_KLAIM_LIFE` | `[terverifikasi]` | `Claim Prop/Flow/Flow_TreatyIn.xml` |
| **Loss Allocation** | Activity | `CountLossAllocation_act` (635 KB, belum ditelusur) | `[dugaan]` | `Claim Non Prop/Activity/CountLossAllocation_act.xml` |
| **DLA** | Activity, service | `GenerateDLAFacin_Act`, `DLAFacintoTreaty_Act`, `HitDLAClaimFacin`. **Kepanjangan belum terverifikasi** | `[pertanyaan terbuka]` | `Claim Fac In/Activity/DLAFacintoTreaty_Act.xml` |
| **SPK** | guard `IsSPK` | Menggerbangi percabangan B2B di titik masuk flow. **Kepanjangan belum terverifikasi** | `[pertanyaan terbuka]` | `Claim Fac In/When/IsSPK.xml` |
| **Fac Retro / Fac Out** | nama Flow `OfferFacRetro`, `OfferFacOut`; flag `.OfferFacIn.IsFacRetro` | Jalur penempatan keluar atas risiko facultative inward yang diterima. Digerbangi `IsFacRetro = 1`; punya tangga persetujuan sendiri dan penomoran sendiri (`POOLDATA.GENERATE_FACRETRO_NO`) | `[terverifikasi]` keberadaan alur / `[dugaan]` arti bisnis | `NB FacIn/Flow/OfferFacRetro.xml`, `Endorsment Fac In/Flow/OfferFacRetro.xml`, `NB FacIn/When/IsFacRetro.xml` |
| **R/I Slip** | nama Flow `InputInwardFacultativeRISlip`, Harness `PrintRISlip`, flag `.OfferFacIn.IsRISlip` | Tahap setelah binding pada siklus facultative; punya rule `Flow` tersendiri dan jalur cetak | `[terverifikasi]` | `NB FacIn/Flow/InputInwardFacultativeRISlip.xml`, `NB FacIn/When/IsNotPrintRISlip.xml` |
| **Banding** | hasil `IsUWAccepted`; Data Transform `SetBandingProposal_DT`; Activity `SetBanding_ACT`; properti `.OfferFacIn.IsBanding`, `.BandingTo` | Salah satu dari enam hasil keputusan akseptasi. Menulis `ProposalAcceptStatus = 4` **dan** `.LetterNo = .BandingTo`, sehingga case dirutekan ke peran yang ditunjuk `.BandingTo` | `[terverifikasi]` mekanisme / `[dugaan]` arti bisnis | `NB FacIn/DataTransform/SetBandingProposal_DT.xml`; `NB FacIn/Flow/InputInwardFacultativeOffer.xml` shape `Decision33` |
| **Spreading** | 34–38 Activity per modul facultative; `CekLimitSpreading_Act`, `SumTSIPremiSpreadedRNM_*` | Pembagian/penyebaran nilai pertanggungan; punya varian per lini bisnis (FIRE, ANEKA, GOLF, MARINECARGO, MBU, PA, TRAVEL). **Rumusnya belum dibaca** — activity terbesar 996 KB | `[terverifikasi]` keberadaan | `NB FacIn/Activity/CekLimitSpreading_Act.xml`; `flows/_SUMMARY-facultative.md` §5 |
| **Limit Akseptasi** | label shape Decision; Activity `GetLimitAkseptasi_ActFlow`, `GetLimitAkseptasi_Act`, `GetLimitAkseptasi_JUW_UW`, `GetLimitAkseptasiLife_Act` | Ambang wewenang akseptasi yang menentukan ke peran mana case dirutekan. **Diambil dari database**: 12 RDB-List berbeda per lini bisnis, plus varian "Banding" | `[terverifikasi]` | `NB FacIn/Activity/GetLimitAkseptasi_ActFlow.xml`; tabel `POOLDATA.M_LIMIT_PROPERTYY`, `POOLDATA.M_LIMIT_NONPROPANDENGG` |
| **Arasapas** | label shape "HIT SERVICE ARASAPAS"; Activity `serviceInsertArasapas_act`, `…RNW_act`, `…EDM_act`; skema Oracle `ARASAPAS` | Integrasi keluar yang dipanggil setelah case disetujui. Rantainya: `GetLinkService` (tabel `M_LINK_SERVICE`) → `Connect-REST` ×2 → `Obj-Save` → tutup case → email | `[terverifikasi]` langkah — **varian Endorsment Fac In saja** | `Endorsment Fac In/Activity/serviceInsertArasapas_act.xml`; OQ-025 |
| **Konversi** | label shape "UPDATE STS KONVERSI"; properti `.StatusService.StsKonversiFacIn` / `StsKonversiFacOut`; `When/IsSuccessHitService.xml` | Status hasil pengiriman data ke sistem luar. Nilai `1` dan `""` menggerbangi apakah alur lanjut atau kembali ke inbox | `[terverifikasi]` mekanisme | `NB FacIn/When/IsSuccessHitService.xml`, `NB FacIn/Activity/UpdateStsKonversiFacOut_Act.xml` |
| **Prorata / ProRate** | properti `.OfferFacIn.ProRateType`; Activity `SetLocalNProrate`, `SetLocalNonMbuProrate`; procedure `FIRE.CEK_PRORATA_TANGGAL` | Perhitungan proporsi atas periode. `ProRateType != 3` mengubah cara premi Nusantara Re dihitung. **Perhitungan tanggalnya ada di database**, bukan di rule | `[terverifikasi]` keberadaan / `[dugaan]` arti nilai | `NB FacIn/Activity/CountPremiNusareRetro_Act.xml`; `Endorsment Fac In/RDBList/SearchSQLRateKPR.xml` |
| **Binding** | label shape `MARKETING (BINDING)`; properti `.OfferFacIn.ConfirmBinding`; Utility `SendEmailBind` | Tahap penegasan penawaran sebelum R/I Slip. `SetRejectProposal` dan `SetReviseProposal` men-set `ConfirmBinding = 0` | `[terverifikasi]` mekanisme | `NB FacIn/Flow/InputInwardFacultativeOffer.xml` shape `Assignment9`; `NB FacIn/DataTransform/SetRejectProposal.xml` |
| **Siklus facultative (NB / RNW / EDM)** | `When/IsNB.xml`, `When/IsRenewal.xml`, `When/IsEDM.xml`; prefix ID case `NB-`, `RNW-`, `EDM-` | Tiga tahap siklus yang dibedakan **saat runtime** oleh `Quotation.StatusBusiness` (1/2/3), bukan oleh basis kode terpisah. Ketiga rule ada di ketiga modul | `[terverifikasi]` mekanisme / arti nilai **belum terverifikasi** | `flows/_SUMMARY-facultative.md` §7; OQ-015 |
| **Akseptasi (treaty inward)** | `DataTransform/Akseptasi_DT.xml`; properti `TreatyIn.StatusAkseptasi`, `ChooseStatusAkseptasi`, `Position`, `PositionUsername` | Tangga persetujuan master treaty inward, **tanpa rule `Flow`**. Empat tingkat posisi: `ReasTreatyInAdmin` → `ReasTreatyInSecHead` → `ReasTreatyInDeptHead` → `ReasTreatyInDirector`; `Accept` di Direktur menghasilkan `Resolve Complete` | `[terverifikasi]` mekanisme / arti nilai **belum terverifikasi** | `Treaty In/DataTransform/Akseptasi_DT.xml` (hash `58b8e650`, identik di `Treaty In Adjustment`) |
| **Revisi treaty (`/R01`)** | `Activity/TreatyInRevisi_post.xml`; properti `TreatyIn.ID`, `OLDID`, `RevisionState` | Nomor revisi disimpan **di dalam string ID**, sufiks `/R01`, `/R02`, …, posisi karakter 10–12; ID lama disalin ke `OLDID`. Jalur revisi **memendekkan** tangga persetujuan | `[terverifikasi]` | `Treaty In Adjustment/Activity/TreatyInRevisi_post.xml`; `RDBList/GetTreatyRevisionID.xml` (`SUBSTR(ID,10,2)+1`) |
| **Data lama (*OldData*)** | 9 `FlowAction/*OldData*` + 14 `Section/*OldData*` | Pasangan layar "sebelum vs sesudah" untuk tiap tab data treaty (limit, layer, share, retensi, CoB, EGNPI, perhitungan share facultative). **Eksklusif `Treaty In Adjustment`** | `[terverifikasi]` | `Treaty In Adjustment/FlowAction/`, `Section/` |
| **Treaty Arrangement** | 81 activity `*TreatyArr*` di `Treaty Contract Out` | Master *term/klausul* kontrak treaty. 16 jenis klausul dapat diedit (BordereAux, CashLossLimit, ClaimCoorperation, EPI, ExGratia, FacIn, PLA, ProfitCommision, Ricomm, TreatyLimit, Portfolio, TerrLimit, …) | `[terverifikasi]` | `Treaty Contract Out/Activity/CancelActivity*.xml`, `SaveTreatyArr*_Act.xml` |
| **Proportional / Non-Proportional** | properti `TreatyIn.ProportionType` | Dua lini treaty inward yang memisahkan layar, perhitungan limit, dan tab. Nilai literal: `"Proportional"`, `"NonProportional"` | `[terverifikasi]` nilai / arti bisnis **belum terverifikasi** | `Treaty In/Section/TreatyInShareProp.xml`, `FlowAction/LimitProportional.xml`, `Activity/TreatyInNPSetTotal.xml` |
| **Position (treaty)** | properti `TreatyIn.Position`, `TreatyIn.PositionUsername` | Nama workbasket disimpan sebagai **nilai data**, bukan sebagai assignment Pega; `PositionUsername` menyimpan **identitas orang** pemilik tugas berikutnya | `[terverifikasi]` | `Treaty In/DataTransform/Akseptasi_DT.xml`; OQ-053 |
| **Retro Life** | 4 Harness `Inbox*` di `Master Contract Retro Life` | Master kontrak retrosesi lini life: business, retro reinsurers list, retro limit reinsurers, security reinsurer — di atas 4 tabel `POOLDATA.*_LIFE` | `[terverifikasi]` | `Master Contract Retro Life/Harness/`; `RDBList/SaveMasterTreaty*_Life_SQL.xml` |
| **Product Name (inward) Life** | `Harness/InwardProductName.xml`; `SaveProductName_Act`, `SaveInwardProductName_Act`, `SetProductNameInward` | Master nama produk life, dengan 7 pemilih master (ceding, currency, policy holder, SOB, RI rate, RI risk, cause of loss). Membawa jalur tulis ke master **treaty inward** | `[terverifikasi]` | `Master Product Name Life/`; OQ-056 |
| **Slot `CARI<n>`** | `InputData.CARI1`…`CARI30`, `TempInputData.CARI3/4/5`, `FlagBanding.CARI11`, `StatusDoc.CARI30` | Properti bernomor yang dipakai sebagai **slot parameter generik** menuju SQL; diisi tepat sebelum pemanggilan, dibaca di `<pyBrowseSQL>` sebagai `{…CARI<n>}` | `[terverifikasi]` pola / arti tiap slot **belum terverifikasi** | `Master Contract Retro Life` (`REINSTYPEID := {TempInputData.CARI3}`); OQ-059 |

---

## 2. Kode dan enumerasi — arti belum terverifikasi

Seluruhnya **dicatat literal apa adanya**. Lihat **OQ-020**.

| Kode | Nilai literal | Di mana dipakai | Bukti |
| --- | --- | --- | --- |
| `EdmType` | `1`, `3` | 13 rule di Endorsement Life; disimpan di kolom JSON `DATA_JSON.EdmType` | `Endorsement Life/RDBList/GetEdmTypeLife.xml` |
| `.Type` (PremiumList) | `QR`, `QP`, `TP`, `TR` | 11 rule di PremiumList Life + Endorsement Life; `TR`/`TP` sering berpasangan | `PremiumList Life/Activity/SubmitPremiumList_Act.xml` |
| Status connector | `Confirm`, `Decline`, `Reject`, `Offer`, `Premium` | transisi flow PremiumList Life | `PremiumList Life/InputPolicyHolder.xml` |
| `.EDMStatus` | `"Old"` | precondition step | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` |
| `.Quotation.BusinessCode` | `!= "40"` | precondition step 1 | `NB Treaty In/Activity/SaveJsonPolisTreatyIn_Act.xml` |
| `.ClaimPaymentType` | `"Claim"`, `"Salvage"` | guard `isClaimTreaty` | `NB Treaty In/When/isClaimTreaty.xml` |
| `.ClaimType` | `"XOL"` | guard `isClaimTreaty` | `NB Treaty In/When/isClaimTreaty.xml` |
| `pyWorkPage.LetterNo` | `"TREATYINDEPTHEAD"` | guard routing | `NB Treaty In/When/ToTREATYDEPTHEAD.xml` |
| `OperatorID.pyTelephone` | `"TREATY1"`, `"SPVTREATY1"` | guard peran — field telepon dipakai sebagai kode peran | `NB Treaty In/When/IsTreaty1.xml` |
| `.ID` | `"1000032"`–`"1000035"` | 4 precondition `Property-Set` | `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` |
| ambang tanggal | `25` | 3 konteks | `PremiumList Life`, `Endorsement Life`, `EDM Treaty In` |
| ambang baris | `50000` | precondition `SaveMasterLPDet` | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` |
| `AcceptStatus` | `"1"`, `"2"` | guard loop tangga komite + cabang pasca-keputusan (keempat modul Komite) | `Komite Claim Life/When/IsKomiteLoop.xml` |
| `TransferType` | `'1'`, `'2'` | gerbang routing komite | `Komite Claim Life/Activity/KomiteRouter.xml` |
| `PaymentType` | `1`–`6` | modul Komite FacIn/Prop/Non Prop | `Komite Claim FacIn/Activity/KomitePost_Adjustment.xml` |
| `TypeComentAnalysis` | `"5"` | gerbang `Obj-Open-By-Handle` | `Komite Claim Prop/Activity/KomitePost_Reject.xml` |
| `RetroID` | `"1000013"` | gerbang `Call InsertJsonClaimLife_Act` | `Komite Claim Life/Activity/KomitePostAdjustment.xml` |
| ambang nominal | `> 30.000.000,00` dan `<= 57.750.000,00` | menentukan komposisi roster komite; mata uang tidak disebut | `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` |
| tanggal cutover | `20250207T000000.000 GMT` | gerbang 2 `Property-Set` | `Komite Claim Life/Activity/KomitePostAdjustment.xml` |
| `PaymentType` | Life: **tidak ada**; Prop `1–6`; Non Prop & Fac In `1–7` | gate di Activity/Section domain klaim | `Claim Fac In/Activity/` (23× `=4`) |
| `STS_REJECT` | `'0'`, `'1'`, `'2'` | gate di 20+ tempat | `Claim Life/Activity/SaveOutStandingLife_Act.xml` |
| `.pyNote` | `"Back"` | gerbang `IsBackStage` (Prop, Non Prop, Fac In) | `Claim Prop/When/IsBackStage.xml` |
| `BusinessCode` (Life) | `L1`…`L11` | satu precondition menguji 11 kode berderet | `Claim Life/Activity/SaveOutStandingLife_Act.xml` |
| `ContentNote` | `"DEATH"` | gate `RDB-List` | `Claim Life/Activity/SaveOutStandingLife_Act.xml` |
| limit wewenang | `30000000.00`, `50000000.00`, `30.00` | gerbang pembuatan case Komite | `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` |
| `ProposalAcceptStatus` | `1`, `2`, `3`, `4`, `7`, `9` | 8 Data Transform + `When/IsUWAccepted` + `When/IsFacout`, di 3 modul facultative. Hasil `IsUWAccepted`: `confirm`, `reject`, `ask`, `banding`, `revise`, `decline`. **Pemetaan nilai → hasil tidak ada di korpus** (OQ-043); nilai `4` ambigu (OQ-044) | `NB FacIn/DataTransform/Set*Proposal*.xml`; `NB FacIn/DecisionTable/IsUWAccepted.xml` tag `<pyTaskStatusXml>` |
| `Quotation.StatusBusiness` | `1`, `2`, `3` | **Diskriminator siklus facultative**: `IsNB` = 1, `IsRenewal` = 2, `IsEDM` = 3. Keempat rule ada di **ketiga** modul | `NB FacIn/When/IsNB.xml`, `When/IsRenewal.xml`, `When/IsEDM.xml`, `When/IsNotEDM.xml` |
| `QuotationData.EdmType` (facultative) | `1`, `2`, `3`, `4` | Endorsment Fac In; `= 4` menjadi prasyarat seluruh 14 rule tipe endorsement | `Endorsment Fac In/When/IsEdmAdjTSI.xml` dkk |
| `QuotationData.Type` (tipe endorsement) | `0`, `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `11`, `12` | Di bawah `EdmType = 4`: 0 `IsEdmAdjRefNo`/`IsEDMRiSlip`, 1 `IsEdmExtendPeriod`, 2 `IsEdmAdjTSI`, 3 `IsEdmAdjRate`, 4 `IsEdmAdjSpreading`, 5 `IsEdmAddObject`, 6 `IsEdmAdjPeriod`, 7 `IsEdmAdjCurrency`, 8 `IsEdmAdjRIC`, 9 `IsEdmAdjInsured`, 11 `IsEdmAdjShareCedant`, 12 `IsEdmAdjCeding`/`IsEdmAdjRIC`/`IsEdmPPNPPH`. **Nilai `10` tidak dipakai** | `Endorsment Fac In/When/IsEdm*.xml` |
| `LetterNo` (token routing) | `KADIVTEKNIK`, `KADIVFINANCIAL`, `KADIVFACULTATIVE`, `MANAGERTEKNIK`, `DEPHEADUNDERWRITER`, `DIREKTURTEKNIK`, `DIREKTURMARKETING`, `SENIORUW`, `DEPTHEADUWLIFE`, `TREATYINDEPTHEAD` | Field nomor surat dipakai memilih Assignment pada shape "Limit Akseptasi". Ditulis `SetBandingProposal_DT` dari `.BandingTo` (OQ-045) | `NB FacIn/When/ToKadivTeknik.xml` dkk |
| `.OfferFacIn.ProRateType` | `3`, `4` | 3 modul facultative; `!= 3` mengubah rumus `PremiLifeNusantaraRe` | `NB FacIn/Activity/CountPremiNusareRetro_Act.xml` |
| `.OfferFacIn.IsFacRetro` / `IsInputFacRetro` / `IsRISlip` | `1` | Gerbang jalur fac out/retro dan cetak slip | `NB FacIn/When/IsFacRetro.xml`, `IsInputFacRetro.xml`, `IsNotPrintRISlip.xml` |
| `.StatusService.StsKonversiFacIn` / `StsKonversiFacOut` | `1`, `""` | Gerbang `IsSuccessHitService` setelah panggilan keluar | `NB FacIn/When/IsSuccessHitService.xml` |
| `Quotation.PolicyType` | `2` | `When/IsDeclarationPolicy.xml` | `NB FacIn/When/IsDeclarationPolicy.xml` |
| `QuotationData.TeamGroup` | `5` | `When/IsTBonding.xml`, bersama guard `.MarketingName` | `NB FacIn/When/IsTBonding.xml` |
| `.OfferFacIn.IsB2B` | `!= "ASM"` | Salah satu kondisi `When/ToDepHeadUW.xml` | `NB FacIn/When/ToDepHeadUW.xml` |
| `Quotation.BusinessOldId` (facultative) | `"L1"`, `"L2"`, `"L3"`, … | `When/IsLife.xml` — kode lini bisnis yang sama dengan OQ-038 (`L1`…`L11` di Komite Claim FacIn) | `NB FacIn/When/IsLife.xml` |
| Pangsa retro | `0.45`, `0.05`; ambang `"3000000000"` | `CountPremiNusareRetro_Act`, identik di 3 modul. **Mata uang tidak disebut**; ambang dibandingkan sebagai string (OQ-046) | `NB FacIn/Activity/CountPremiNusareRetro_Act.xml` |
| `TreatyIn.StatusAkseptasi` | `"Accept"`, `"Reject"`, `"Decline"`, **`"Resolve Complete"`** | Status akseptasi master treaty inward — **4 nilai, dihitung mesin**. `"Resolve Complete"` tidak pernah dipilih pengguna. Beda `Reject` vs `Decline` **tidak dijelaskan korpus** (OQ-052) | `Treaty In/DataTransform/Akseptasi_DT.xml`; `Activity/TreatyInDeclineConfirmation_postact.xml` |
| `TreatyIn.ChooseStatusAkseptasi` | `"Accept"`, `"Reject"`, `"Decline"` | **3 nilai, dipilih pengguna** — masukan bagi `Akseptasi_DT` | `Treaty In/DataTransform/Akseptasi_DT.xml` |
| `TreatyIn.Position` | `""`, `"ReasTreatyInAdmin"`, `"ReasTreatyInSecHead"`, `"ReasTreatyInDeptHead"`, `"ReasTreatyInDirector"`, `"ReasTreatyInGroupLeader"` | Nama workbasket disimpan sebagai **nilai data** | idem |
| `TreatyIn.ProportionType` | `"Proportional"`, `"NonProportional"` | Pembeda lini treaty inward | `Treaty In/Section/TreatyInShareProp.xml` dkk |
| `TreatyIn.RevisionState` | `""`, `"1"` | `"1"` memendekkan tangga persetujuan (Sec Head → `Resolve Complete`) | `Treaty In/DataTransform/Akseptasi_DT.xml` |
| `TreatyIn.EDMState` | `"1"` (dari `Param.type=="revision"`), `"3"` (dari `Param.type=="adjustment"`) | Pembeda dua proses di `Treaty In Adjustment`. Nilai `"2"` **tidak muncul** | `Treaty In Adjustment/DataTransform/TreatyCreateEDM.xml` |
| `Param.type` | `"revision"`, `"adjustment"` | Parameter pembeda proses | idem |
| `TreatyIn.ViewState` | `"1"`, `""` | Mode tampilan; menggerbangi visibilitas tombol | `Treaty In/Section/TreatyInActionButtons.xml` |
| `TreatyIn.IsEditData` | `"1"` | Menggerbangi visibilitas tombol | idem |
| `OperatorID.pyTelephone` (treaty) | `"TREATY1"`, **`"TREATY2"`**, `"SPVTREATY1"`, **`"SPVTREATY2"`** | Field telepon operator menyimpan **kode peran**; dua nilai terakhir baru ditemukan di Tahap 5 (OQ-027) | `Treaty In/DataTransform/Akseptasi_DT.xml` |
| `OperatorID.pyWorkBasketList(n)` | indeks `(1)`, `(2)` | Otorisasi tombol memakai **indeks tetap** pada daftar workbasket operator — 36 kemunculan pada `(2)` (OQ-051) | `Treaty In/Section/TreatyInActionButtons.xml` |
| `REINSTYPEID` | **tidak ada nilai literal di korpus** | Diisi dari `TempInputData.CARI3/4/5`, diteruskan sebagai `p_REINSTYPEID` ke 3 procedure. Daftar nilainya **tidak dapat dinyatakan** (OQ-057) | `Master Contract Retro Life/RDBList/*.xml` |
| Sufiks revisi ID | `"/R01"`, `"/R02"`, … pada posisi karakter **10–12** | Nomor revisi disimpan di dalam string `TreatyIn.ID`; ter-hardcode di rule **dan** di SQL (OQ-055) | `Treaty In Adjustment/Activity/TreatyInRevisi_post.xml` |

---

## 3. Nama workbasket / antrean yang terlihat

`[terverifikasi]` Bukti paling konkret untuk **OQ-007** (tidak ada rule identitas/otorisasi di
korpus). Di domain treaty inward nama-nama ini muncul sebagai **nilai properti `TreatyIn.Position`**,
bukan sebagai assignment Pega; di domain facultative sebagai `<pyRuleName>` pada shape Assignment
ber-`pyRouteTo=Custom` (pemetaannya tidak terbaca — **OQ-024**).

| Nama | Konteks | Bukti |
| --- | --- | --- |
| `ReasTreatyInAdmin` | treaty inward — titik masuk tangga & tujuan `Reject` | `NB Treaty In/Flow/InputRealizationTreatyIn.xml` `<pyRuleName>`; `Treaty In/DataTransform/Akseptasi_DT.xml` (`DATA-PORTAL!AKSEPTASI_DT`) |
| `ReasTreatyInSecHead` | treaty inward — tingkat 1 | idem |
| `ReasTreatyInGroupLeader` | treaty inward — **diuji sebagai posisi, tidak pernah di-`SET`** | idem |
| `ReasTreatyInDeptHead` | treaty inward — tingkat 2 | idem |
| `ReasTreatyInDirector` | treaty inward — tingkat akhir, `Accept` → `Resolve Complete` | idem |
| `komitepnc` | tingkat 1 komite | `Komite Claim FacIn/Activity/KomiteRouter.xml` |
| `komitepnc2` | tingkat 2 komite | idem (juga Prop, Non Prop) |
| `komitepnc3` | tingkat 3 komite | idem |
| `komitepnc4` | tingkat 4 komite | idem |
| `ReasFacInMarketing` (377) | facultative inward | sapuan korpus facultative |
| `ReasFacInAdmin` (143) | facultative inward | idem |
| `ReasFacInUnderwriting` (114) | facultative inward | idem |
| `ReasFacInTeamLeader` (109) | facultative inward | idem |
| `ReasFacInSeniorUnderwriting` (105) | facultative inward | idem |
| `ReasFacInUnderwritingFinancial` (91) | facultative inward | idem |
| `ReasFacInDepHeadUnderwriting` (91) | facultative inward | idem |
| `ReasFacInManagerTeknik` (80) | facultative inward | idem |
| `ReasFacInGroupLeader` (80) | facultative inward | idem |
| `ReasFacInJuniorUnderwritingA` (78) | facultative inward | idem |
| `ReasFacInTechnicalDirector` (68) | facultative inward | idem |
| `ReasFacOut` (132) | fac out / retro | idem |
| `ReasFacOutAdmin` (47) | fac out / retro | `NB FacIn/Flow/OfferFacRetro.xml` tag `<pyMOName>` |
| `ReasFacOutHead` (42) | fac out / retro | idem |
| `ReasFacOutGroupLeader` (13) | fac out / retro — **hanya varian Endorsment** | `Endorsment Fac In/Flow/OfferFacRetro.xml` |
| `ReasFacOutTechnicalDirector` (13) | fac out / retro — **hanya varian Endorsment** | idem |

---

## 4. Singkatan yang ditemukan tetapi belum terjabarkan

`[terverifikasi]` Tak satu pun singkatan di bawah dijabarkan di dalam korpus. Kolom **Bukti** memuat
path rule tempat singkatan itu benar-benar muncul — **bukan** tebakan kepanjangannya.

| Singkatan | Muncul di | Status | Bukti (path + rule) |
| --- | --- | --- | --- |
| `EDM` | nama modul, sufiks rule, properti `EdmType`, nilai `param.isFOR` | **kepanjangan belum terverifikasi** (OQ-014) | `EDM Treaty In/Flow/InputAddendumTreatyIn.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTTREATY!INPUTADDENDUMTREATYIN`) |
| `GISFW` / `GCNMFW` | prefix class framework | **kepanjangan belum terverifikasi** (OQ-008) | class `ASM-FW-GISFW-WORK`, `ASM-FW-GCNMFW-WORK-PNC` pada `<pxInsName>` di seluruh korpus |
| `SFAGIS` | prefix class `ASM-SFAGIS-WORK-ENDORSEMENT` (31 rule), `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` / `-ACCOUNT` (5 masing-masing) | **kepanjangan belum terverifikasi** (OQ-049) | `Endorsment Fac In/Activity/SetValueToEDMWork.xml` (`ASM-SFAGIS-WORK-ENDORSEMENT!SETVALUETOEDMWORK`) |
| `ASM`, `RNM`, `GCNM` | prefix kunci RDBList (`ASM!SAVETREATYIN`, `RNM!CekSTSKonversiJson`) | **arti belum terverifikasi** (OQ-008) | `Treaty In/RDBList/SaveTreatyIn.xml` (`ASM-FW-GISFW-INT-TREATY_IN!ASM!SAVETREATYIN`); `Endorsment Fac In/Activity/serviceInsertArasapas_act.xml` (rujukan `RNM CekSTSKonversiJson`) |
| `XOL` | nilai `.ClaimType`; nama rule spreading/limit | **kepanjangan belum terverifikasi** | `NB Treaty In/When/isClaimTreaty.xml` (`ASM-FW-GISFW-WORK!ISCLAIMTREATY`, `.ClaimType = "XOL"`) |
| `TP` / `TXOL` | FlowAction spreading | **kepanjangan belum terverifikasi** | `Treaty In/FlowAction/SpreadingTPDtl.xml`, `Treaty In/FlowAction/SpreadingTXOLDtl.xml` |
| `QS` | activity pengambil data master | **kepanjangan belum terverifikasi** | `Treaty In/Activity/FetchQSfromMasterXOL.xml` |
| `ROL` | activity + tombol `Save ROL Profile` | **kepanjangan belum terverifikasi** | `Treaty In/Activity/TreatyInSaveROL.xml`; `Treaty In/Section/TreatyInActionButtons.xml` (`<pyLabel>Save ROL Profile`) |
| `EGNPI` | **1.047 kemunculan** di `Treaty In`; properti `EGNPIAMOUNT`, `EGNPIPROPORTION`, `EGNPIAMOUNTNP` | **kepanjangan belum terverifikasi** | `Treaty In/Activity/TreatyInEGNPIListValue.xml`, `Treaty In/Activity/TotalEgnpi.xml`, `Treaty In/FlowAction/DetailEGNPI.xml` |
| `EPI` | klausul arrangement treaty | **kepanjangan belum terverifikasi** | `Treaty Contract Out/Activity/SaveTreatyArrEPI_Act.xml`, `Treaty Contract Out/Activity/CancelActivityEpi.xml` |
| `PLA` | klausul arrangement treaty; juga limit klaim non-prop | **kepanjangan belum terverifikasi** | `Treaty Contract Out/Activity/SaveTreatyArrPLA_Act.xml`; `Claim Non Prop/RDBList/GetLimitTONPPLA.xml` |
| `Ricomm` | klausul arrangement treaty; tombol `Fix Soa Upload RIComm` | `[dugaan]` reinsurance commission — **belum terverifikasi** | `Treaty Contract Out/Activity/SaveTreatyArrRicomm_Act.xml`; `Treaty In/Section/TreatyInActionButtons.xml` |
| `LOL` | klausul arrangement treaty | **kepanjangan belum terverifikasi** | `Treaty Contract Out/Activity/SaveTreatyArrMinLOL.xml`, `SaveTreatyArrMinLOLMB.xml` |
| `RIC` | tipe endorsement facultative | **kepanjangan belum terverifikasi** | `Endorsment Fac In/When/IsEdmAdjRIC.xml` |
| `SOB` | FlowAction pencarian; pemilih master life | **kepanjangan belum terverifikasi** | `Treaty In/FlowAction/TreatyInSearchSoB.xml`; `Master Product Name Life/FlowAction/ChooseSOB.xml` |
| `CoB` | FlowAction + DataTransform | **kepanjangan belum terverifikasi** | `Treaty In/FlowAction/CoBList.xml`, `Treaty In/DataTransform/AddCob.xml`, `Treaty In/DataTransform/SetCoBID.xml` |
| `TSI` | properti `.TSILiability`; nama activity spreading | **kepanjangan belum terverifikasi** | `NB FacIn/Activity/CountPremiNusareRetro_Act.xml` (`.TSILiability <= "3000000000"`); `NB FacIn/Activity/SumTSIPremiSpreadedRNM_FIRE_Act.xml` |
| `KPR` | rule rate facultative | **kepanjangan belum terverifikasi** | `Endorsment Fac In/RDBList/SearchSQLRateKPR.xml`, `SearchSQLRateNonKPR.xml` |
| `OGP` / `ONP` | activity perhitungan komisi/hasil | **kepanjangan belum terverifikasi** | `NB FacIn/Activity/CountOverridingCommOgp_Act.xml`, `NB FacIn/Activity/CountResult1Onp_Act.xml` |
| `NCL` / `CL` | varian limit akseptasi kredit | **kepanjangan belum terverifikasi** | dirujuk `NB FacIn/Activity/GetLimitAkseptasi_ActFlow.xml` sebagai `ASM GetLimitAkseptasiKreditNCL_SQL` / `…KreditCL_SQL` |
| `PPN` / `PPH` | gerbang tipe endorsement; activity pajak | `[dugaan]` jenis pajak Indonesia — **belum terverifikasi** | `Endorsment Fac In/When/IsEdmPPNPPH.xml`; `NB FacIn/Activity/SetPPNPPH.xml` |
| `PKS` (dalam `IsPKSASM`) | gerbang percabangan di 3 modul facultative; **kondisinya tidak terbaca** | **kepanjangan belum terverifikasi** (OQ-029) | `NB FacIn/When/IsPKSASM.xml` (`ASM-FW-GISFW-WORK!ISPKSASM`) |
| `SPK` | guard `IsSPK` — percabangan **B2B** di titik masuk flow klaim facultative; **kondisinya tidak terbaca** | **kepanjangan belum terverifikasi** (OQ-029) | `Claim Fac In/When/IsSPK.xml` (`ASM-FW-GCNMFW-WORK!ISSPK`) |
| `B2B` | properti `.OfferFacIn.IsB2B`; label shape "B2B" | **arti belum terverifikasi** | `NB FacIn/When/ToDepHeadUW.xml` (`.OfferFacIn.IsB2B != "ASM"`); `Claim Fac In/Flow/Register_Flow.xml` (`<pyMOName>B2B`) |
| `DLA` | activity generator dokumen klaim facultative; service `HitDLAClaimFacin` | **kepanjangan belum terverifikasi** | `Claim Fac In/Activity/DLAFacintoTreaty_Act.xml` (`ASM-FW-GCNMFW-DATA-OBJECT!DLAFACINTOTREATY_ACT`); `Claim Fac In/ConnectREST/HitDLAClaimFacin.xml` |
| `PNC` | class kerja `…WORK-PNC` (Claim Fac In); sufiks procedure `PEGA_JSON_KLAIM_PNC`; sasaran routing `komitepnc*` | **kepanjangan belum terverifikasi** (OQ-036) | `Claim Fac In/Flow/Register_Flow.xml` (`ASM-FW-GCNMFW-WORK-PNC!REGISTER_FLOW`) |
| `CNP` | sufiks activity klaim & komite non-prop | **kepanjangan belum terverifikasi** | `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml`; `Komite Claim Non Prop/Activity/InsertOSKlaimCNP.xml` |
| `TNP`, `TRT` | sufiks stored procedure akseptasi klaim | **kepanjangan belum terverifikasi** (OQ-002) | `Claim Prop/RDBList/SaveDataToOsAkseptasiNP.xml` (memanggil `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` / `…KLAIMTRT`) |
| `CA`, `CFS` | activity generator dokumen klaim non-prop | **kepanjangan belum terverifikasi** | `Claim Non Prop/Activity/GenerateCACNP_Act.xml`, `Claim Non Prop/Activity/GenerateCFS_act.xml` |
| `CWP` | activity komite non-prop | **kepanjangan belum terverifikasi** | `Komite Claim Non Prop/Activity/KomitePostAdjustmentCWP.xml` |
| `KMT` | sufiks activity komite | **kepanjangan belum terverifikasi** | `Komite Claim FacIn/Activity/HitServiceToKasirKMT_Act.xml`; `Komite Claim Life/…/Generate_NoAccept_KMT_Life` |
| `OS` | prefix tabel akseptasi (`OS_AKSEPTASI_KLAIM`, `OS_AKSEPTASI_KLAIM_LIFE`) | **kepanjangan belum terverifikasi** | `Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` (`INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`) |
| `WPC` | activity lini life | **kepanjangan belum terverifikasi** | `PremiumList Life/Activity/WPCLife_Act.xml` |
| `UW` | sufiks FlowAction akseptasi; label shape `UNDERWRITING` | **kepanjangan belum terverifikasi** | `NB Treaty In/FlowAction/DeptHeadTreatyIn_UW.xml` |
| `JUW` | `When/ToJUW_A.xml`; label shape `JUNIOR UNDERWITER A`/`B` (ejaan apa adanya di korpus) | `[dugaan]` Junior Underwriter — **belum terverifikasi** | `NB FacIn/When/ToJUW_A.xml`; `NB FacIn/Flow/InputInwardFacultativeOffer.xml` (`<pyMOName>`) |
| `DivHead` | variabel limit ter-hardcode | **kepanjangan belum terverifikasi** | `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` (`Local.LimitMaxDivHead = 50000000.00`) |
| `SPL` | tombol `fix Spl(dev)`; activity | **kepanjangan belum terverifikasi** | `Treaty In/Activity/TreatyInFixSpl.xml`; `Treaty In/Section/TreatyInActionButtons.xml` |
| `SOA` | tombol `Fix Soa Upload RIComm` | **kepanjangan belum terverifikasi** | `Treaty In/Section/TreatyInActionButtons.xml` (`<pyLabel>`) |
| `EQS`, `RSMD`, `FLEXAS` | tabel rate lewat db-link `@ASMD.SINARMAS.CO.ID` | **kepanjangan belum terverifikasi** (OQ-017) | `NB FacIn/RDBList/SearchRatePolisEQS_SQL.xml`; objek `M_EQS_RATE@…`, `M_RSMD_RATE@…`, `M_FLEXAS_RATE@…` |
| `MBU` | skema Oracle; sufiks activity spreading & prorate | **kepanjangan belum terverifikasi** (OQ-016) | procedure `MBU.F_CEK_HURUF` dipanggil dari `NB FacIn`; `Endorsment Fac In/Activity/SetLocalNonMbuProrate.xml` |
| `ASMD` | nama database link `@ASMD.SINARMAS.CO.ID` | **arti belum terverifikasi** (OQ-017) | `NB FacIn/RDBList/GetKurs.xml` (`LST_KURS_STANDARD@ASMD.SINARMAS.CO.ID`) |

---

## 5. Apa yang harus dijawab sebelum glossary ini menjadi `CONTEXT.md`

`[terverifikasi]` Glossary ini **siap dipakai sebagai peta istilah**, tetapi **belum siap menjadi
kamus normatif**. Empat OQ menentukan sebagian besar isinya:

| OQ | Yang diblokir di sini |
| --- | --- |
| **OQ-020** | Arti seluruh kode/enumerasi di Bagian 2 (51 entri) — `ProposalAcceptStatus`, `StatusAkseptasi`, `PaymentType`, `EdmType`, `QuotationData.Type`, `StatusBusiness`, `ProRateType`, `AcceptStatus`, `TransferType`, `Type` life |
| **OQ-008** | Arti prefix `ASM` / `RNM` / `GCNM` dan kepanjangan `GISFW` / `GCNMFW` |
| **OQ-014** | Kepanjangan `EDM`, dan apakah EDM/Addendum/Endorsement memang satu konsep |
| **OQ-057** | `REINSTYPEID` dan kode `OR` — **tidak ada satu pun nilai literal di korpus untuk dikutip** |
| **OQ-049** | Kepanjangan `SFAGIS` dan mengapa ia class terpisah |
| **OQ-043** | Arti hasil `confirm`/`reject`/`ask`/`banding`/`revise`/`decline` — baris tabel keputusannya tidak ikut terekspor |

`[terverifikasi]` Bagian 3 (nama workbasket) adalah **bukti paling konkret untuk OQ-007** — korpus
tidak memuat satu pun rule identitas/otorisasi, sehingga nama-nama antrean inilah satu-satunya
petunjuk peran yang tersedia.

**Jangan mengisi kolom mana pun di atas dengan tebakan saat menyusun `CONTEXT.md`.** Yang belum
terverifikasi harus tetap ditandai demikian sampai pemilik peran menjawab.
