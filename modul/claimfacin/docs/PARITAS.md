# PARITAS — tombol dan aksi XML Claim Fac In ↔ status

> 10-10-2026, tahap 1 dari 2 (prompt "IMPLEMENTASI CLAIM FAC IN, TAHAP 1 DARI 2"). Sumber inventaris: flow
> `Register_Flow`, Section `InputRegister`, `InputRegisterDetail`, `Catastrope_Sec`, `ViewHistoryClaim`,
> `InputInwardFacultativeDtl`, `ProgresClaim_Sec`, `SubProgresClaim_Sec`, `ViewPolis` (harness ChoosePolis), `ProtectDOL`,
> `ClaimComiteeReject` + `SureRejectClaim_section`, `InputEstimasiAdmin`, `InputEstimasiDetail`,
> `PropertyItemListGridEstimation*`, `ShowItemPA`, `Estimasi`, `EstimasiPA`, `EstimasiMarine`, `Pla_Dtl`, `ClaimSurvey`,
> `InputEstimasi`, `ShowObjectAdj`, `ItemListEstimation*`, `Adjusment_SC`, `InputAdjustment`, `ViewCedantPanel`,
> `ClaimComite`, `PrintDLA_dtl`, `PreventRejectClaim` di `D:\XML\RNM_BRD\Claim Fac In` (setiap pxButton / pxLink dan setiap
> medan ber-`on change` activity). Sumber status: definisi layar server `backend/models/layar_register.go`,
> `layar_estimasi.go`, `layar_adjustment.go`, `komite.go`, `tutup.go`, `dla.go` — satu-satunya tempat tombol lahir;
> frontend merender pohon tata apa adanya (`frontend/components/TataView.tsx`). Peta aksi server:
> `backend/services/aksi.go` (`penangan`, 82 aksi) dan `services/adjustment.go` (`penanganAdjustment`).
>
> **Status:**
> - **dibangun** — tampil menurut kondisi XML yang sama (visible-when / disabled-when / read-only-when), aksinya port
>   activity-nya.
> - **nonaktif-OQ** — tampil sesuai section, tetapi nonaktif; keterangan `OQ-CFI-nn` di atribut `title` (`OQ.md`).
> - **tidak tampil di XML** — kondisi tampilnya statis `1=2` / `NEVER` / tidak pernah benar: tidak dibangun.
> - **tahap 2** — milik modul `komiteclaimfacin` (prompt §2 "Di luar lingkup").
>
> Penanda: `[penyimpangan sadar]` = perilaku sengaja berbeda dari XML (alasan + perilaku Pega lama dicatat);
> `[inferensi]` = disimpulkan, tidak tertulis di XML.

## 1. Flow dan halaman awal

| XML | Aksi XML | Status |
| --- | --- | --- |
| `Register_Flow` Start1 → Decision4 B2B (`When/IsSPK`: `OfferFacIn.IsB2B = "SPK"`) | jalur B2B langsung ke Choose Surveyor | **dibangun sebagian** — kasus lahir lewat **Add Claim** tanpa polis, sehingga `IsSPK` salah dan kasus selalu masuk Input Register. Kasus B2B lahir dari harness `New` / `NewSample` yang tidak diekspor (OQ-CFI-29) |
| Assignment1 "Input Register" (operator saat ini, FlowAction `InputRegister`) | `TAHAP = InputRegister` | **dibangun** — tab *Process* (worklist pembuat) |
| Assignment7 "Input Estimasi" (WorkList Custom = pembuat, FlowAction `InputEstimasi`) | `TAHAP = InputEstimasi` | **dibangun** |
| Assignment3 "Choose Surveyor" (WorkBasket Custom, FlowAction `InputSurveyor` berlabel "Input Adjustment") | `TAHAP = InputSurveyor` | **dibangun** — workbasket `ReasKlaimTeknik` (bawaan b: `ReasPNCTeknik` tidak ada di DEV); switch **Teknik** di halaman awal bagi anggota workbasket (`GET /hak`) |
| Decision5 / Decision8 `IsBackStage` (`.pyNote = "Back"`) | Back Input Estimasi → Input Register; Back Choose Surveyor → Input Estimasi | **dibangun** (§3, §4) |
| END52 Resolved-Completed | `CloseClaim` (`ASMForceCaseClose`) | **dibangun** — tab *Resolve* |
| Pra-proses FlowAction `InputRegister` (`CallActivityInputRegister` + DT `InsertObjects_dt`) | — | **dibangun** (`models.PraRegister`, `objek.go`) |
| Pra-proses `InputEstimasi` (`InputEstimationPre` + DT `SetEstimation_DT`) | — | **dibangun** (`models.PraEstimasi`) |
| Pra-proses `InputSurveyor` (`SetTypePDFAdjustment` + DT `SetStartDateAdjustment`) | — | **dibangun** (`models.PraAdjustment`; StartDateAdjustment hanya bila kosong, §8) |
| `InsertProgressClaim` (setiap assignment dibuka) | `PROGRESSCLAIM` / `SUBPROGRESSCLAIM` | **dibangun** — idempoten per (kasus, posisi), isi `PEGA_PROGRESSCLAIM` / `PEGA_SUBPROGRESSCLAIM` ditulis ulang tanpa memanggil prosedur |
| Validate FlowAction `ValidateDate` | — | **tidak dibangun** — rule tidak diekspor (OQ-CFI-13); pemeriksaan tanggal `CheckDate*` tetap jalan |
| Kotak masuk Beranda | — | **tidak dibangun** — bukan bagian Pega (prompt §2 "Di luar lingkup"; `menu.ts` tanpa `antreanBeranda`) |
| Tabel komite di bawah inbox | — | **tahap 2** |

## 2. Input Register — Section `InputRegister` / `InputRegisterDetail`

`InputRegisterDetail` juga tampil hanya-baca sebagai tab **View Registration** (Input Estimasi, `IsTreatyIn != 1`) dan
**Registration** (Choose Surveyor); medan terkunci menurut kondisi section yang sama (`isEstimation`, `IsAnyAccept`).

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| View Status Payment Premi | REST `getPremiumPaidOn` (M_LINK_SERVICE) | **nonaktif-OQ** OQ-CFI-14 |
| Reject Claim | local action `RejectSurveyClaim` (pra `SetRejectClaim_pre`) | **dibangun** — pop-up `ClaimComiteeReject` (`BukaRejectClaim`, modal `tolak`); tombol **Yes** = `SendRejectClaimToKomite2` **dibangun** tahap 2 (§5) |
| Tab Register / Policy Detail & Claims History / Progress Claim | — | **dibangun** |
| Save (VIS `PolicyNo` kosong) | click:save | **dibangun** (`Simpan`) |
| Submit (VIS `IsError = 1`) | local action `ProtectDOL` | **dibangun** (`BukaProtectDOL`, modal `protectDOL`) |
| Submit (VIS `IsError = 0`) | `CheckListEstimasi_Act` + `ProteksiDataRegister_Act` + finishAssignment | **dibangun** (`Submit`) → Input Estimasi |
| Choose Polis (NA `isEstimation = 1`) | harness `ChoosePolis` (Section `ViewPolis`) | **dibangun** (`BukaPolis`, modal `pilihPolis`) |
| Date of Loss / Report Date / Received Date (change) | `CheckDate_Act` / `CheckDateReport_Act` / `CheckDateReceived_Act` | **dibangun** — langkah 9/13 `+1` atas teks tanggal = OQ-CFI-12 (dibangun menurut hasil XML) |
| Choose Cause of Loss | harness `CauseofLoss_Harness` → `GetNameCauseofLoss` | **dibangun** — pop-up daftar `pilihan/sebab` (`BukaSebab`), pilih = `GetNameCauseofLoss`. Master Cause of Loss milik modul lain, tidak diubah |
| Catastrophe (radio) / Non-Catastrophe type (change) | `SetDefNonCatastrope` | **dibangun** |
| Catastrophe List | harness `CatastrofeList` | **dibangun** — pop-up daftar `pilihan/katastrofe` (`BukaKatastrofe`); pilih `SetCatastrope`, input `InputCatastrope`, simpan `SaveCatasrtope` (tabel `CATASTROPHE`) |
| ikon Edit / Save Catastrophe | `SetEditCatastrope` Edit / Save | **dibangun** |
| Reporter Status (change) | `GetReportStatus_Act` | **dibangun** |
| grid Ceding Co (LS132) — Choose | `GetCeding_act` | **dibangun** |
| grid pilih objek (LS137, per lini) — centang | `CheckListEstimasi_Act` | **dibangun** |
| Consultant ID / Adjuster ID (autocomplete) | `SetConsultant` / `SetAdjsuter` | **dibangun** |
| ikon + Consultant / Adjuster | harness `MstAdjusterConsultant` (Data-Portal) | **nonaktif-OQ** OQ-CFI-32 — Claim Prop menyimpan master barunya lewat rute pinjaman `POST /api/adjuster-consultant`, yang di `cmd/api/rakit.go` (`ruteDipinjam`) hanya dipinjamkan ke `claimprop`; memanggilnya dari modul ini juga gagal penjaga `TestPanggilanLintasModulTerdaftar`. Pemilihan adjuster / consultant yang sudah ada tetap jalan |
| Location of Loss (change) | `SetTempLocation` | **dibangun** |
| Country / Province / City / District / Region / Zip Code | `BrowseCountry_RD` … `InputKodePos1` | **dibangun** — saringan provinsi menurut negara diperbaiki (§8) |
| Currency (change) | `InputCurrencyValueAct_Register` | **dibangun** |
| Claim Estimate (change) | `CheckDoubleClaim_Act` | **dibangun** — nomor polis tertulis mati dibuang (bawaan c) |
| Policy Detail — View Retro List | harness `RetroList_Harnness` | **dibangun** — pop-up hanya-baca `FacRetroList` dari halaman polis yang dibaca ulang (`BukaRetroList`, OQ-CFI-16) |
| Claim Status (kronologi) / Claim History | `T_VIEW_SUGGEST` / klaim lain polis ini | **dibangun** |
| Progress Claim — Input Progres Claim | FlowAction `InputSubProgressClaim` (`AddDelProgress`) | **nonaktif-OQ** OQ-CFI-15 (hanya mengubah clipboard; sumber Progres1/2 tidak diekspor) |
| Field `ClaimData.UserTeknis` (LS177) | — | **tidak tampil di XML** (VIS `1=2`) |

### Pop-up `ViewPolis` (Choose Polis) dan `ProtectDOL`

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| Search Type / Search | `SearchPolis_act` + `GetPolisForClaim_SQL` | **dibangun** — teks cari **diikat** (§7 butir 1); hasil = `pilihan/polis` (`FACINPRODUCTION` berindeks, batas 500 baris) |
| tautan Policy No (grid hasil) | `CopyNB_Act` | **dibangun** (`CopyNB`, param `nopolis|prodke`) — salin `JSON_POLIS.DATA_JSON` → `OfferFacIn`; polis tanpa `JSON_POLIS` ditolak terang (OQ-CFI-11). Polis NB Fac In sistem baru tidak terbaca (OQ-CFI-30) |
| detail polis (`InputParam.CARI4 = 'true'`) | Section `InputInwardFacultativeDtl` | **dibangun** |
| Cancel / Submit | — / `SetInputParam_Act` | **dibangun** — `CheckPeriodePolicy` tidak dibangun: langkah 1 `Exit-Activity` tanpa syarat (OQ-CFI-34) |
| ProtectDOL: grid kasus kembar, Cancel / Submit | `InsertObjectItemList_DT` + `CheckListEstimasi_Act` + finishAssignment | **dibangun** (`Submit`, konteks `protectDOL`) |

## 3. Input Estimasi — Section `InputEstimasiAdmin` / `InputEstimasiDetail`

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| No Claim, View Status Payment Premi | — / REST | **dibangun** / **nonaktif-OQ** OQ-CFI-14 |
| Send to PIC Claim (NA `IsPicTransfer != 1 \|\| IsError > 1`) | `SetDisable_ACT` + finishAssignment (bila `ObjectList(1).IsFacretro != 1` atau `PlaStatus = 1`) | **dibangun** (`SetDisable`) → Choose Surveyor, `POSITION = ReasKlaimTeknik` |
| Tab Estimation / Policy Detail & Claims History / Progress Claim / View Registration | — | **dibangun** |
| Back (VIS `IsCFS != 1`) | DT `BackFromRegister` + `BackToRegister_act` + finishAssignment | **dibangun** — nilai uji langkah 8–11 tidak ditulis (§7 butir 3) |
| grid objek: Download Claim Face Sheet | `CLaimFaceSheet_Act` | **dibangun** — nomor klaim **K**, `OS_AKSEPTASI_KLAIM` STS 0 per estimasi, `JSON_KLAIM`, konversi (efek hanya produksi), `PrintFaceClaim` / `IsCFS`. Berkas PDF = OQ-CFI-20 |
| grid objek: Print PLA (VIS `ClaimData.ExGratia = 0`, NA `IsFacretro != 1 \|\| PlaStatus = 1 \|\| IsError > 1`) | `SetIndexObj_Act` + harness `Pla_Dtl` | **dibangun** (`BukaPLA`, modal `pla:o`) — `PlaStatus` hanya sesudah Submit berhasil (§7 butir 2) |
| Pla_Dtl: Send Email / Message / Share Retro / Remarks / Submit | `GeneratePLA` → `GeneratePLATreaty_Act` / `ObjectItem.GeneratePLA` | **dibangun sebagian** — PLA treaty **G** + revisi `/n` dibangun tetapi gerbangnya `CekLimit.CARI1 = 1` tidak pernah tercapai (OQ-CFI-24); PLA fac retro (`ObjectItem.GeneratePLA`) tidak diekspor (OQ-CFI-21); penerima surel tidak dibangun (OQ-CFI-22) |
| grid objek: Outstanding Summary | `GetAllData_Act` | **dibangun** — pop-up ringkasan `pilihan/outstanding` (dua kolom mata uang diperbaiki, §8) |
| panel objek: Add (item) | `addRow` | **dibangun** (`TambahItem`); nama item kembar diperiksa saat nama / coverage dipilih (`ProtectionObjectItem_Act`) |
| panel objek: Delete item | `HapusItem` | **dibangun** |
| panel objek: Object Name / Occupation / Coverage / Object Id / Coverage ID (change) | `SetObjectItem` / `GetOccupation` / `GetCoverage*` / `GetCoverageAneka_Act` / `GetCurencyCoverage_Act` | **dibangun** (`SetObjectItem`, `PilihOkupasi`, `PilihCoverageFire`, `PilihAneka`, `PilihCoverageAneka`, `PilihCoverageObjek`) |
| panel item MBU (EditAction `PropertyItemListGridEstimationMBU_FA`) | merender grid item itu sendiri | **dibangun** — `[penyimpangan sadar]` baris item MBU membuka section `Estimasi` (Pega tidak dapat mengisi estimasi MBU) |
| panel item: Add estimasi | `ValidateInputEstimate_act` | **dibangun** |
| Estimation Date (change) | `ProtectionDate_Act` (+ `Set7Hours`) | **dibangun** — `Set7Hours` tanpa padanan (tanggal sudah WIB, §8) |
| Currency (change) | `SetConvertValueKurs_Estimation` | **dibangun** |
| Estimation Gross / Estimasi > TSI (change) | `CheckEstimateValue` | **dibangun** — pzInsKey tertulis mati langkah 25 dibuang (bawaan c) |
| Delete estimasi | `DeleteValueEstimation` + `CountSpreadingClaim_ACT` | **dibangun** |
| Deductible (% / Type / Amount, change) | `CountTSI_Act` | **dibangun** |
| Spreading Claim: View Retro (TreatyType FAC retro) | local action `ShowRetro` | **dibangun** — pop-up hanya-baca dari halaman (`BukaRetro`) |
| Spreading Policy / Spreading Claim yang dapat disunting, **Save Spreading** (`LS42`) | `SetValueSaveSpreading` → `SaveSpreadingSP_act` (kata sandi) | **tidak tampil di XML** — bergerbang `pyWorkPage.ClaimData.ExGratia = 1`; penulis satu-satunya `InsertObjects_dt` 11 menulis `0`. Kata sandi tidak ditulis di mana pun (OQ-CFI-17) |
| Spreading: Share % (change) | `CheckTotalSpreadingPct_Act` | **tidak tampil di XML** — medan RO=ALWAYS, event tidak pernah terpicu (§7 butir 10) |
| Tambah / Delete baris spreading | `addRow` / `deleteRow` | **tidak tampil di XML** (VIS `… && 1=2`) |
| baris Total `InputParam.CARI41 / CARI42` | — | **tidak dibangun** — halaman requestor tanpa penulis |

## 4. Choose Surveyor / Input Adjustment — Section `ClaimSurvey` / `InputAdjustment`

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| Claim No / View Status Payment Premi / Remark from Inputor | — / REST / — | **dibangun** / **nonaktif-OQ** OQ-CFI-14 / **dibangun** |
| Tab Claim Details (Registration / Estimation / Adjustment & Acceptation) / Policy Details and Claim History / Progress Claim | Section `InputEstimasi` / `ShowObjectAdj` | **dibangun** |
| Back | `BackToEstimasi` + finishAssignment | **dibangun** → Input Estimasi |
| Close Claim | local action `PreventRejectClaim` | **dibangun** (modal `tutup`, §5) |
| grid objek: Print DLA (VIS `IsFacretro = 1`, NA `DLAStatus` bukan 0/kosong) | `SetIndexObject` + harness `PrintDLA_dtl` | **dibangun** (`BukaDLA`, modal `dla:o`) |
| PrintDLA_dtl: Send Email / Message / Share Retro / RemarksDLA / Submit | `ChooseDla_Act` → `GenerateDLAFacin_Act` (P) / `DLAFacintoTreaty_Act` (S) | **dibangun** — DLA fac **P**; DLA treaty **S** bergerbang `CekLimit.CARI1` (OQ-CFI-24); tanpa lompatan mundur, `HitDLAClaimFacin` hanya produksi (§7 butir 4–5); `OS_AKSEPTASI_KLAIM` `UPDATE DLA_NO, DLA_DATE`; surel lewat outbox produksi; Draft DLA saat Send to Committe |
| grid item PA / Travel: Send / Sent Claim to committee | harness `ClaimCommittee` | **nonaktif-OQ** OQ-CFI-25 |
| grid item: View Status Payment Claim / View Payment Attachment | REST `getPaymentClaim` / `getPayAttachment` | **nonaktif-OQ** OQ-CFI-14 |
| grid Adjustment: + (VIS item `IsKomite` kosong/0) | `SetIndexObject_ACT` + `CountTotalEstimasi_Act` | **dibangun** (`CountTotalEstimasi`) — `[inferensi]` baris baru (OQ-CFI-23) |
| grid Adjustment: ikon hapus (VIS `IsKomite != 1`, NA `IsApproved = 1`) | `DisableSendComite` | **dibangun** — baris yang punya `KOMITE_ID` tidak terhapus (server) |
| grid `.KomiteList` ItemListEstimation | — | **tidak tampil di XML** (VIS `1=2`); daftar komite tampil di panel adjustment |
| footer `CountSpread.CARI24–27` | — | **tidak dibangun** — halaman requestor tanpa penulis |
| View Inc Fac Cedant Panel | harness `ViewCedantpanels` (`ViewCedantPanel`) | **dibangun** (`BukaCedant`, modal `cedant:…`); Choose / Choose All = `SetCedant_act` 1 / 2 |
| DLA No Ceding / DLA No SOB (change) | `SetDLACedingSOB` | **dibangun** |
| Payment Type (change) | `SetAdjTypePayment` | **dibangun** — kode 1–7 tampil apa adanya (label OQ-CFI-18) |
| Choose Currency (change) | `CheckCurrency` | **dibangun** |
| Gross Adjustment (change) | `SetGrossAdjustment` | **dibangun** — langkah 6 (pre=false) selalu menolkan deductible % (literal XML) |
| Deductible Type / % / Value (change) | `SetNilaiResikoSendiri` | **dibangun** |
| VAT (change, PT 4 / 6) | `SetValueAdjusterFee` | **dibangun** |
| Payable To / Specify (change) | `SetPayable_Act` | **dibangun** |
| grid Ceding Co: Choose (Payable) | `SetPayableTo` | **dibangun** — langkah 1–4 ber-REMARK; tanpa efek tersimpan |
| Name of Bank (autocomplete) | `Result.pxResults` `SetPayable_Act` | **dibangun** (`PilihRekening`) |
| View KMT NO | `ViewKomite_act` | **nonaktif** — `KomiteNo` sudah tampil sebagai medan (dibaca ulang saja) |
| Ex Gratia (change) | `CekExGratia` | **dibangun** |
| Spreading Adjustment: Add (VIS adjustment `ExGratia = 1`) / Share % (change) / Delete | `addRow` / `CountSpreadingAdjustment` / `deleteRow` | **dibangun** (`TambahSpreadAdj`, `CountSpreadingAdjustment`, `HapusSpreadAdj`) |
| Spreading Adjustment: View Retro | local action `ShowRetro` kelas SpreadingRisk | **nonaktif-OQ** OQ-CFI-25 |
| BreakDown Spreading Quota Share | — | **dibangun** (hanya-baca) |
| Save | click:save | **dibangun** (`SimpanAdjustment`) |
| Send to Committe (NA `IsKomite = 1` atau spreading kosong) | `SendPICProtect_Act` → harness `Comittee` | **dibangun** (`SendPICProtect`, modal `komite:…`, §5) |
| Acceptation (VIS `AcceptanceStatus = 1`) | `SaveAcceptation` + `HitServiceToKasir_Act` | **dibangun** — aktif sesudah komite tahap 2 menulis `AcceptanceStatus`; kasir lewat outbox produksi, sekali per klik (§8). RALAT 10-10-2026: `HitServiceToKasir_Act` 3 `getStatusKonversi_Act` bertransisi pasca `.StatusKonversi=="1"` (COUNT `reinsurance.trloss_detail_t`, hanya produksi) - selainnya keluar; kini ditegakkan (`Acuan.StatusKonversi`) |

## 5. Komite, penutupan, penolakan

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| ClaimComite: Date / Initial / Chronology / Extent / Legal Liability / Salvage / Adjuster Fee / Remarks | harness `Comittee` | **dibangun** |
| Send Claim to Committee (VIS `Protect.CARI1 = 1 && CARI2 = 1` + isian wajib) | `DraftGenerateDLAFacin_Act` + `CreateKMTNo_Act` + `SetListKomite_act` | **dibangun** (`KirimKomite`) — kasus `KMT-` TT2 lahir (`T_WORK_CLAIM` LINI FACIN, `TAHAP Komite_Flow`, `POSITION` = `OPERATOR_ID` tingkat 1, `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`), `KomiteNo`, progres "Waiting Committee", kronologi, email outbox produksi. Bundel PT 2 = OQ-CFI-28. Tanggal / circumstances tidak ditimpa (§7 butir 6–7) |
| `SendPICProtect_Act` proteksi lampiran LOD / DLA / SPGR / ADU / Invoice / Salvage, rekening, premi | — | **dibangun** — cacah `DOCUMENT_CLAIM` × master `T_KATEGORI_DOC_KLAIM` FAC; premi lunas pola Claim Prop (saldo premi non-produksi dilewati) |
| PreventRejectClaim: Close Without Payment (centang) | postValue | **dibangun** (`SetelAlokasiSalvage`) |
| PreventRejectClaim: Yes (CWP) | `SendCloseClaimToKomite` (TT4) | **dibangun** tahap 2 (`SendCloseClaimToKomite`, KCF-03) — 2 `Remark` / `Remark_Close` := Remarks; 4–6 adjustment belum diputus → pesan VERBATIM, kasus tidak lahir; 7.1–7.4 kasus `KMT-` tanpa adjustment, `TRANSFER_TYPE` 4, tangga satu tingkat `ReasClaimDeptHead` (prompt §5 butir 7: akun + email tertulis mati diganti workbasket); 7.5–7.7 `Komite.CARI1`, kronologi "Request close claim without payment " + KMT; 7.9 email outbox produksi. Penjaga permintaan ganda (`PesanTutupKomiteGanda`) `[penyimpangan sadar]`. Sesudah lahir: LS23 (isian + Yes) tersembunyi, LS28 label sukses, checkbox LS21 tetap tampil (tanpa syarat). Chronology / Extent / Policy Liability pop-up disimpan di kepala kasus komite (migrasi komiteclaimfacin 643, OQ-KCFI-03) |
| PreventRejectClaim: Yes (close) | `CloseClaim` | **dibangun** — `ValidationAdjustmentKomite`, setiap pesan menghentikan (§7 butir 8), `OS_AKSEPTASI_KLAIM` STS 4, konversi STS 4 (produksi), kronologi, `JSON_KLAIM`, Resolved-Completed; 12 `ASMForceCaseClose` `CloseAllSubCases=true`: kasus komite KMT- klaim ini yang masih menunggu ikut ditutup (10-10-2026, `TutupKomiteAnak`) |
| `UpdateTotalJob_sql` (`mst_user_teknis.TOTAL_JOB`) | CloseClaim 7 | **tidak dibangun** — `ClaimData.UserTeknis` tanpa penulis (LS177 VIS `1=2`); langkah tidak pernah berjalan (OQ-CFI-33) |
| ClaimComiteeReject: Yes (naJika `IsReject = 1`, tampil selama `Komite.CARI1` kosong) | `SendRejectClaimToKomite2` (TT3) | **dibangun** tahap 2 (`SendRejectClaimToKomite2`, KCF-03) — 2 `Remark` / `Remark_Close` := Remarks; 4–6 sudah ada akseptasi / estimasi belum Face Claim → pesan VERBATIM (pesan terakhir menang); 7.1–7.4 kasus `KMT-` tanpa adjustment, `TRANSFER_TYPE` 3, tangga `ReasClaimDeptHead`; 7.5–7.7 kronologi "Request Reject claim " + KMT; 7.9 email outbox produksi. Sesudah lahir: isian LS4 tetap tampil berisi ketikan (clipboard Pega; `models.IsianPopUpTutup` dibawa `BawaSementara`), Yes diganti label sukses. Chronology / Extent / Policy Liability pop-up TT3 / TT4 disimpan di kepala kasus komite (`repository.TeksKomite`, migrasi komiteclaimfacin 643, keputusan OQ-KCFI-03) dan tampil di layar komite. Keputusan komite (CLAIMREJECTED, OS STS 2, Resolved-Rejected) = modul `komiteclaimfacin` |
| Perluasan tangga `ApprovalKomite_Act`, keputusan komite, `SaveAccept_ACT` (OS STS 1) | — | **tahap 2 dibangun** — modul `komiteclaimfacin` (`modul/komiteclaimfacin/docs/PARITAS.md`), tulis balik lewat kontrak `kontrak.KlaimFacInKomite` (`services/kontrak_komite.go`) |

## 6. Tidak tampil di XML (tidak dibangun)

- Grid / tombol bergerbang `pyWorkPage.ClaimData.ExGratia = 1` (Spreading Policy / Claim yang dapat disunting, Save
  Spreading): `ExGratia` tingkat klaim selalu `0`.
- Tombol Add / Delete spreading estimasi (`1=2`), medan `ClaimData.UserTeknis` (`1=2`), grid `.KomiteList`
  ItemListEstimation (`1=2`).
- Grid estimasi Travel di section `Estimasi` (Travel memakai `EstimasiPA`).
- ~~Bagian "Success Create Request to Committee" PreventRejectClaim (`Komite.CARI1` hanya terisi oleh TT4, OQ-CFI-27).~~
  Tahap 2 (10-10-2026): `Komite.CARI1` kini terisi TT3 / TT4 — label tampil di pop-up Reject / Close sesudah kasus
  komite lahir, tombol Yes tersembunyi.

## 7. Sebelas arah perbaikan prompt §5 — `[penyimpangan sadar]`

| # | Rule | Perilaku Pega lama | Dibangun | Uji |
| ---: | --- | --- | --- | --- |
| 1 | `SearchPolis_act` / `GetPolisForClaim_SQL` | teks cari disisipkan ke SQL (`{ASIS:InputData.CARI1}`) | bind `:1` | `repository/sql_test.go` `TestCariPolisTeksTerikat` |
| 2 | `SetIndexObj_Act` 1 | `PlaStatus = 1` saat pop-up Print PLA dibuka | `PlaStatus = 1` hanya sesudah `GeneratePLA` berhasil | `models/perbaikan_test.go` `TestPlaStatusHanyaSesudahGeneratePLA` |
| 3 | `BackFromRegister` 8–11 | medan kosong diisi "Test", "0812", "IDR", tanggal sekarang | tidak mengisi apa pun | `TestBackFromRegisterTanpaNilaiUji` |
| 4 | `ChooseDla_Act` | `RetroList.TotalClaim = "123"`; transisi kosong #9 lompat mundur ke `JMP1` | tanpa lompatan: log `AKSEPATSI` dan `DLA` masing-masing sekali; `DLAList` / `TotalClaim` tidak dibangun — nol pembaca di korpus (OQ-CFI-25) | `handlers/perbaikan_test.go` `TestDLATanpaLompatanMundurDanHitDLAHanyaProduksi` |
| 5 | `ChooseDla_Act` 10 | REST `HitDLAClaimFacin` tanpa penjaga produksi | efek `dla-klaim` di outbox hanya bila produksi | uji yang sama (non-produksi: nol efek) |
| 6 | `CreateKMTNo_Act` L4 | memotong `pxCreateDateTime`, `DateOfLoss`, tanggal polis ke 8 karakter lalu menyimpannya | tanggal tersimpan tidak diubah | `models/komite_test.go` `TestTandaiKirimKomiteTanpaMengubahTanggalDanKronologi` |
| 7 | `CreateKMTNo_Act` L11 | `TempCommiteClaim.CircumtansesCouseOfLoss` ditimpa pzInsKey | circumstances tidak ditimpa | uji yang sama |
| 8 | `CloseClaim` 6 | `Param.ErrorMsg` dibaca tanpa halaman parameter: pesan tampil, klaim tetap tertutup | setiap pesan `ValidationAdjustmentKomite` / proteksi DLA menghentikan penutupan | `handlers/perbaikan_test.go` `TestCloseClaimBerhentiPadaProteksiDLA`, `models/komite_test.go` `TestCekDLATutup` |
| 9 | `InsertProgressClaim` 2.1 | `ProdKe := REGISTER` tanpa kutip → posisi kosong; dipanggil setiap assignment dibuka | posisi `REGISTER`; sisip hanya bila (kasus, posisi) belum ada, posisi lain `Done` (isi `PEGA_PROGRESSCLAIM`, upsert) | `TestProgresPosisiRegister` |
| 10 | `CheckTotalSpreadingPct_Act` | nilai tanpa kutip di #1 | **tidak berlaku** — pemicunya `change` medan Share % yang RO=ALWAYS di tiga section estimasi; activity tak terjangkau, tidak dibangun | — (bukti: `Section/Estimasi.xml`, `EstimasiPA.xml`, `EstimasiMarine.xml`) |
| 11 | `GetSequenceNumber_SQL` | `ParamSeq.CARI3` (tanggal procedure) tidak pernah diisi | periode dari tanggal tutup buku lewat `inti/backend/penomor` (`HitungPeriodeNomor`), pola Claim Prop / Claim Non Prop | uji milik `inti/backend/penomor`; `models/perbaikan_test.go` `TestRakitNomor` (bentuk nomor) |

## 8. Penyimpangan sadar lain dan `[inferensi]`

**Penyimpangan sadar** (perilaku Pega lama → di sini):

1. **Write-through** (pola Claim Prop): Pega menyimpan clipboard hanya pada Save / Obj-Save; di sini setiap aksi
   menyimpan halaman. Aksi yang keluar karena pesan validasi dibatalkan seluruhnya.
2. **Pra-proses setiap aksi**: Pega menjalankan pra-proses FlowAction saat assignment dibuka; di sini setiap muat. Akibat
   yang terlihat: `SetStartDateAdjustment` 1 menimpa `StartDateAdjustment` setiap kali dibuka → di sini hanya bila kosong;
   item `IsKomite` di-reset `0` (bila `AktifButton` kosong) sehingga tombol + adjustment muncul lagi sesudah aksi
   berikutnya — Pega sesudah assignment dibuka ulang.
3. `SetKomiteList_ACT` 2 menghitung ulang `ValueAdjustment` SEMUA baris adjustment item (kondisi
   `PaymentType != "4" || PaymentType != "3"` selalu benar) → baris fee / salvage menjadi 0. Tidak ditiru.
4. `CLaimFaceSheet_Act` 23.4 `TotalEstimasiConvertIDR` dari akumulator yang tidak dinolkan antar item → di sini
   `TotalEstimasiValue` adjustment = Σ ConvertValue item itu sendiri.
5. `GenerateDLAFacin_Act` dipanggil sekali untuk SETIAP baris spreading 10015 setiap item (objek ber-dua item retro
   menerbitkan dua nomor DLA) → di sini setiap jenis nomor (P / S) sekali per Submit.
6. `SaveAcceptation` menjalankan `HitServiceToKasir_Act` dua kali pada klik pertama → sekali per klik, lewat outbox.
7. `SetValueAdjusterFee` 13: kondisi `IsError < 2 || …` dibangun menurut teks XML.
8. `InsertObjects_dt` 8.6.25 menulis `OccupationList` hanya ke calon terakhir → setiap calon Aneka menerima
   `OccupationList` lokasinya sendiri.
9. `ProtectDownloadFaceClaim` 2 membaca klaim tersimpan hanya di produksi (pola Claim Prop OQ-CP-11): di luar produksi
   CFS objek kedua tidak tertahan.
10. `GetAllData_Act` (Outstanding Summary): kurs Travel / PA dicocokkan pada `.CARI15` yang tidak dipilih SQL, dan kolom
    Currency Marine menampilkan `.CARI10` yang tidak dipilih → di sini mata uang baris.
11. `GeneratePLATreaty_Act` 14–15 menolkan `TotalEstimasi` / `TotalGrossEstimasi` item (indeks salah) → tidak dibawa.
12. `BrowseProvince_RD`: section mengirim `ID2 = CountryID` dan membaca `.Note` / `.ID` yang tidak ada → disaring nama
    negara, `ProvinceID` ikut terisi.
13. `GetAddressCeding` membaca `M_CLIENT.JSONDATA` → tabel datar `CLIENT_ADDRESS` (keputusan work owner Claim Prop
    08-10-2026).
14. `SetchronologyKlaimFacIn` 1.1 memetakan tiga nama orang tertulis mati ke jabatan untuk setiap baris kronologi pelaku
    itu → jabatan pelaku dari roster `EMAILKOMITE` FACIN untuk setiap baris, selainnya "Claim Admin" (bawaan c, pola
    Claim Prop). Jabatan roster yang tidak dipetakan XML (mis. Claim Supervisor) kini tampil sebagai jabatannya.
15. `Set7Hours` (Estimation Date): pergeseran GMT → WIB tanpa padanan; tanggal disimpan dalam zona Jakarta.
16. Satu perlakuan presisi (pola Claim Prop): `@divide` berskala 4 / 8 di beberapa rumus → presisi penuh, pembulatan
    hanya di penyimpanan `NUMBER(38,10)` dan tampilan.
17. `IsBondingAndCustomBonds` C merujuk `IsCustomBonds` yang tidak ada (yang ada `IsCustomBond`) → bernilai salah.
18. Medan tanpa pembaca di korpus tidak disimpan: `ShareNusare`, `PUCLStatus`, `ReceiverClaim`, `TFAllObj`, `OurRef`,
    `StartPeriod` / `EndPeriod`, `DLAList`, `PLAList`.
19. `CheckDate_Act` 9 / 13 `+1` atas teks tanggal (penggabungan teks) → pesan "Start Date Policy" tidak pernah terpasang;
    dibangun menurut hasil XML (OQ-CFI-12).
20. `GetSpreadingMarine_Act` (pra-proses `EstimasiMarine_FA`, Pega: setiap pane item Marine dibuka) → pra-proses Input
    Estimasi hanya bila Spreading Policy item masih kosong (`models.LengkapiMarine`): pra-proses di sini berjalan setiap
    aksi, dan menjalankannya ulang menimpa Spreading Claim hasil estimasi.
21. Pesan validasi membatalkan seluruh aksi (write-through); Pega membiarkan perubahan clipboard sebelum pesan (mis.
    baris spreading yang dihapus tetap hilang walau total < 100%). Layar validasi = halaman sebelum penangan + isian
    layar + pesan; contoh: 50 / 30 / 20 → nolkan baris yang akan dihapus dan naikkan baris lain dalam satu kiriman, lalu
    hapus.
22. Kolom milik komite (`ACCEPTED_NO`, `ACCEPTED_DATE`, `ACCEPTANCE_STATUS`, `IS_APPROVED` di `T_CLAIM_ADJUSTMENT`) tidak
    ikut UPDATE simpan halaman klaim (`Kolom.MilikKomite`): klaim tidak pernah menulisnya, dan simpan yang membaca
    halaman sebelum komite memutuskan tidak boleh menimpa keputusan itu.
23. Submit pop-up (Print PLA `GeneratePLA`, Comittee `KirimKomite`) memeriksa ulang tombol pembukanya di server (Pega:
    pop-up hanya dapat dibuka lewat tombol yang aktif); hapus adjustment bertaut komite = 409.

**`[inferensi]`:**

1. Tombol + grid Adjustment tidak ber-`addRow`; `CountTotalEstimasi_Act` menulis `.Adjustment(<LAST>)` → baris baru
   ditambahkan lalu diisi; pemeriksaan 20–21 atas baris terakhir yang sudah ada (OQ-CFI-23).
2. Langkah Java "hapus spreading / treaty / value yang sama" tidak diekspor → baris berkunci sama (TreatyType; atau
   CurrencyID + TreatyType) dibuang kecuali yang pertama, menurut agregasi sesudahnya.
3. `.CurencyAdjustment` (pilihan Choose Currency) disusun ulang dari estimasi item ber-`PrintFaceClaim = 1`, bukan kolom.
4. Pilihan Cedant Panel disimpan sebagai `CEDANT_CHOICE` (`*` = Choose All); daftar disusun ulang dari polis.
5. `KunciCedingKasir` `@substring(pyWorkPage.Quotation.CedingCo, 0, 8)`: `pyWorkPage.Quotation` tanpa penulis →
   `OfferFacIn.QuotationData.CedingCo`.
6. Pembagi nol (`@divide` / `/`) bernilai 0 di rumus yang membagi properti yang dapat kosong (`hasilAkhirDLA`, dll.).
7. `IsError` kosong sama dengan 0 (tombol Submit Input Register).
8. Label ReportType / IndividualRiskType / Payable dari screenshot Claim Prop (kelas properti sama).
9. Perbandingan `CheckLimit_Act1` 9 numerik (`FindData.HASILD2` desimal); `@divide(x,1,2)` = ROUND_HALF_UP.

### RALAT 10-10-2026 — transisi pasca-langkah yang terlewat alat dump tahap 1

Alat dump activity tahap 1 hanya mencetak prasyarat langkah (`pyStepsPreCondParams`), tidak transisi pasca-langkah
(`pyStepsTransition` + `pyStepsTransParams`). Sensus ulang seluruh `Claim Fac In/Activity` (tahap 2, review kode):

| Rule / langkah | Transisi pasca | Status |
| --- | --- | --- |
| `HitServiceToKasir_Act` 3 | `.StatusKonversi=="1"` T=2 F=6 | **diperbaiki** — gerbang `Acuan.StatusKonversi` (hanya produksi) sebelum IDOfBank / muatan |
| `ValidateInputEstimate_act` 6 (pre: tanpa spreading) | `true` T=6 | **diperbaiki** — `TambahEstimasi` keluar sesudah baris dibuang |
| `SetNilaiResikoSendiri` 15 (pesan "can not be filled by Zero") | `true` T=6 | **diperbaiki** — keluar; pesan langkah 17 tidak ikut tampil |
| `CLaimFaceSheet_Act` 1 (`ProtectDownloadFaceClaim`) | `Protect.HASIL1==1` T=6 | sudah ditegakkan (`validasi`) |
| `CLaimFaceSheet_Act` 6 | `1==1` T=6 | jalur mati (`ExGratia` selalu 0, langkah 4–6 tidak dibangun) |
| `CloseClaim` 5.2 | `.IsFacretro==1 && .RemarksDLA==""` T=6 | sudah ditegakkan (§7 butir 8) |
| `SetProtectionEstimation` 23 | F=1 tanpa label | tidak dibangun (pemanggil = tombol OQ-CFI-25) |

## 9. Frontend

`frontend/` merender pohon tata server apa adanya (pola Claim Non Prop, disalin, tidak diimpor); bahasa layar Inggris,
kulit token `--cfi-*` salinan `.kelola-user`, kelas `claimfacin__*`.

- **Halaman awal** (`pages/ClaimFacIn.tsx`): tab *Process* (worklist pembuat) / *Resolve*, switch **Teknik** (workbasket
  `ReasKlaimTeknik`, nonaktif bagi akun bukan anggota, `GET /hak`), **Add Claim** (`POST /kasus`). Tanpa tabel komite
  (tahap 2).
  > **RALAT 10-10-2026** (tahap 2). Kalimat lamanya dikutip utuh, tidak dihapus: *"Tanpa tabel komite (tahap 2)."* → tabel
  > **Committee** (`components/TabelKomite.tsx`) di bawah inbox, tanpa switch, tidak ikut tab, tampil bila `GET /hak` →
  > `komite` (pemegang `ReasClaimSPVB` atau akun / workbasket roster FACIN aktif); isinya `GET
  > /api/komite-claim-fac-in/kasus` (rute pinjaman); klik baris = layar komite di tempat (`onBukaModul`,
  > `MODUL_DIPINJAM`), Back / Submit kembali ke inbox. "View more details" di layar komite membuka berkas ini hanya-baca
  > (`?lihat=1`).
- **Panel bersarang** (`components/rincian.ts`, `TataView.tsx`): grid ber-`rincian` membuka panel baris
  `layar.panel["<prefiks>:<jalur>(n)"]`; panel dirender ulang lewat `TataView`, sehingga grid di dalam panel membuka
  panelnya sendiri (objek → item → estimasi / adjustment → detail adjustment). Satu baris terbuka per grid; status buka
  dikunci per tab dan per grid (`kunciGrid`) — temuan cek layar: tanpa itu baris terbuka bocor antar-tab.
- **Alamat aksi**: setiap unsur membawa `konteks` (`''` layar utama, kunci panel, atau kunci modal) dan `indeks` (baris
  grid di konteks itu; 0 = bukan sel grid / tombol Add kepala grid). Permintaan memuat `masukan` (layar utama + panel +
  modal), `mode` (`layar.mode` ditimpa isian lokal 14 jalur mode), dan `tahap`.
- **422**: pesan + pesan per medan tampil, layar yang dikirim server dirender ulang tanpa membuang isian lokal (juga di
  dalam modal / pop-up).
- **Pop-up sisi klien** sesudah aksi server tanpa efek: Cause of Loss (`pilihan/sebab` → `GetNameCauseofLoss`),
  Catastrophe (`pilihan/katastrofe` → `SetCatastrope` / `InputCatastrope` / `SaveCatasrtope`), Outstanding Summary
  (`pilihan/outstanding`), View Retro / View Retro List (hanya-baca dari halaman), hasil Choose Polis (`pilihan/polis`,
  tautan nomor polis → `CopyNB` `nopolis|prodke`).
- **Lampiran** (`components/PanelLampiran.tsx`, pola Claim Prop): kategori dari server (`T_KATEGORI_DOC_KLAIM` FAC).
- Isian angka: hanya angka, ribuan otomatis, rata kanan; tanggal `dd-mm-yyyy`.

## 10. Isolasi kotak masuk dan pembaca per-ID (prompt §6 butir 1)

Kasus Claim Fac In: `T_WORK_CLAIM.LINI = 'FACIN'`, `TAHAP` selalu terisi (`InputRegister` / `InputEstimasi` /
`InputSurveyor` / `Komite_Flow`; penutupan hanya mengisi `STATUS_WORK`, `repository/gudang.go` `sqlTutupKasus`),
awalan `CLM-` / `KMT-` (OQ-CFI-02), nomor `SEQ_WORK_CLAIM` (unik lintas lini). Dibaca dari SQL modul-modul itu, tidak
disunting.

| Modul | Saringan kotak masuk | Pembaca / penulis per-ID | Baris FACIN |
| --- | --- | --- | --- |
| Claim Life | `NVL(w.TAHAP, :tahapCadangan) = :tahap`, tahap "Input Register" / "Outstanding Claim" / "Medical Check" / "Claim Analis" (`claimlife/backend/repository/inbox.go` `sqlInboxWhere`, `models/tahap.go`) | **`GET /api/klaim-life/{id}`** → `KlaimLife.AmbilHeader` `SELECT … FROM T_GENERAL_CLAIM WHERE ID = :1` **tanpa LINI**. Penulis per-ID memeriksa tahap: `TahapBerlaku` → `TahapDariNama("InputRegister")` tidak dikenal → `TahapDariPeran(POSITION)` (pembuat / `ReasKlaimTeknik` ≠ `ReasLife*`) tidak dikenal → ditolak; UPDATE tahap / tutup bersyarat `NVL(TAHAP, …) = <label Life>`; `DELETE /api/klaim-life/{id}` gagal terang (ADR-U-0031) | **kotak masuk: tidak cocok** (TAHAP tanpa spasi, tidak pernah NULL). **Detail per-ID: kepala klaim FACIN terbaca** bila ID-nya diketik di URL Claim Life (hanya-baca; keadaan yang sama sudah berlaku untuk `CLMP-` / `CLMNP-`) → **OQ-CFI-31** (perbaikan = saringan LINI di claimlife, tidak disunting) |
| Komite Claim Life | `AND (w.LINI = :lini OR w.LINI IS NULL)` (`komiteclaimlife/backend/repository/komite_inbox.go`) | `sqlKasusKomite` `WHERE g.ID = :2 AND (w.LINI = :3 OR w.LINI IS NULL)` — membuka, memutuskan, riwayat semuanya lewat `Kasus` | tersaring (`FACIN` ≠ `LIFE`, tidak NULL) |
| Claim Prop | `w.LINI = 'PROP' AND w.ID LIKE 'CLMP-%'` (`claimprop/backend/repository/gudang.go`) | `Keadaan` `WHERE w.ID = :1 AND w.LINI = :2`; lampiran lewat `Keadaan` | tersaring |
| Komite Claim Prop | `w.LINI = 'PROP'` ketat + awalan `TKMT-` (`komiteclaimprop/backend/repository/gudang.go`, `tangga.go`) | `WHERE g.ID = :1 AND w.LINI = :2 AND w.ID LIKE :3` | tersaring |
| Claim Non Prop | `w.LINI = 'NONPROP' AND w.ID LIKE 'CLMNP-%'` (`claimnonprop/backend/repository/gudang.go`) | `WHERE w.ID = :1 AND w.LINI = :2` | tersaring |
| Komite Claim Non Prop | `w.LINI = 'NONPROP'` ketat + awalan `KMTNP-` | `WHERE g.ID = :1 AND w.LINI = :2 AND w.ID LIKE :3` | tersaring |

Arah sebaliknya: setiap kueri `T_WORK_CLAIM` modul ini menyaring `w.LINI = 'FACIN'` (kotak masuk juga `w.TAHAP <> 'Komite_Flow'`,
`repository/gudang.go`), sehingga kasus Life `CLM-` tidak terbaca di Claim Fac In.
