# PARITAS — Komite Claim Fac In

Setiap isian / tombol / langkah korpus `Komite Claim FacIn` (dan penyerahan TT3 / TT4 di `Claim Fac In`) ↔ status di
sistem baru. Dibangun 10-10-2026 atas prompt work owner "IMPLEMENTASI KOMITE CLAIM FAC IN, TAHAP 2 DARI 2" dengan
keputusan **KCF-01** (roster FACIN → workbasket, cadangan SPV B, pita SPV B), **KCF-02** (perluasan tangga ikut XML),
**KCF-03** (TT3 / TT4 untuk Fac In, `T_GENERAL_KOMITE` diubah), **KCF-04** (kasus komite lama tidak dimigrasi). Pola:
Komite Claim Prop / Non Prop (disalin, tidak diimpor). OQ: [`OQ.md`](OQ.md).

Status: **ada** = dibangun persis; **ubah** = dibangun dengan penyimpangan sadar (alasan); **tidak** = tidak dibangun.

## 1. Menu dan pintu masuk

| Korpus | Status | Catatan |
| --- | --- | --- |
| Menu `komiteclaimfacin` | **ubah** | tanpa menu (perintah work owner 09-10-2026, migrasi inti 949): modul `frontend/layar.ts`, dipasang bagi pemegang menu `claimfacin` (`MODUL_DIPINJAM` `frontend/App.tsx`), rute dipinjam (`ruteDipinjam` `cmd/api/rakit.go`) |
| Worklist assignment "KomiteRouter" (`Flow/Komite_Flow.xml` ASSIGNMENT63) | **ubah** | tabel "Committee" di bawah inbox Claim Fac In (`GET /api/komite-claim-fac-in/kasus`), hanya bila `GET /api/claim-fac-in/hak` → `komite` (pemegang `ReasClaimSPVB` atau akun / workbasket roster FACIN aktif); tanpa switch, tidak ikut tab; klik baris = layar komite di tempat (`onBukaModul`), Back / Submit kembali ke inbox |
| KomiteRouter S6 `.KomiteList` + S6.1 `AssignTo = .KomiteID` (`KomiteAproval == 0`) | **ada** | S6.1 transisi pasca `true` → keluar: Pega menugaskan baris menunggu **pertama**; di sini baris tingkat berjalan (`KOMITE_URUT = KOMITE_COUNT`) - setara pada alur berurutan. S6 gerbang `Primary.TransferType=='2'` pre=false → selalu berjalan: TT3 / TT4 ikut dirutekan |
| KomiteRouter S1–S5, S7 | **tidak** | ter-remark (`//`) |
| Decision IsKomiteLoop (`AcceptStatus == "1" && KomiteCount <= KomiteLoop`) | **ada** | benar → kembali ke assignment (POSITION = KomiteID tingkat berikut); salah → END (`T_WORK_CLAIM.STATUS_WORK` Resolved-Completed, POSITION kosong) |

## 2. Pra-proses `ViewTransferDtl`

| Langkah | Status | Catatan |
| --- | --- | --- |
| `SetValueKomite` S1 `.TransferType != 2` → keluar | **ada** | TT3 / TT4 tanpa pra-proses (tanpa grid adjustment / total / tangga perluasan) |
| S2–S7.1.1 `DataTempAdj`: adjustment ber-`KomiteNo` = KMT ini | **ada** | dibaca dari kasus klaim induk lewat `kontrak.KlaimFacInKomite` (`Baca`) |
| S8.2 / S8.3 inisial, buang adjustment yang sudah diakseptasi | **ada** | |
| S9–S11.2 total per mata uang (PT 4 / 6 → "Adjuster Fee" + VAT, PT 3 → "Salvage"), baris "Total in IDR" | **ada** | kurs = `POOLDATA.GETCURRENCYSTANDARD` bila `CurrencyDol` kosong — `CurrencyDol` tidak disimpan Claim Fac In (OQ-KCFI-05) |
| S12 nama jenis reas, S13 total | **ada** | `REINSURANCETYPE` |
| S14 `ApprovalKomite_Act` | **ada** | §3 |
| S15 `SetProteksiSubmiteKomite` | **ubah** | §5 butir 1: Submit hanya pemegang tingkat berjalan; pemetaan dua akun tertulis mati dan lolos "IT Developer" **dibuang** (gugur, prompt §3) |

## 3. `ApprovalKomite_Act` — perluasan tangga (KCF-01, KCF-02)

| Langkah | Status | Catatan |
| --- | --- | --- |
| S1–S2.2 total adjustment IDR, penanda retro | **ada** | |
| S3 `Local.Retro == 1` → keluar | **ada** | Fac Retro melewati perluasan saja (tangganya dibentuk sisi klaim, `RosterKomiteCalon` tahap 1) |
| S4 Obj-Browse roster bila `KomiteCount == 1` | **ada** | `EMAILKOMITE` STS_KLAIM FACIN aktif (workbasket sesudah migrasi 640) |
| S5–S5.2 pita SPV B: `pyPosition == "SPV B"` **dan** 30.000.000 < total ≤ 57.750.000 → satu jenjang atas | **ubah** | "SPV B" = pemutus tingkat 1 anggota `ReasClaimSPVB` **saja** (KCF-01; anggota SPVA + SPVB sekaligus = SPV A, keputusan OQ-KCFI-06); konstanta bernama `BatasBawahPitaSPVB` / `BatasAtasPitaSPVB`; hanya DEGREE > 1 pertama |
| S6–S6.2 buang calon yang sama dengan tingkat 1 (`KomiteList(1).KomiteID == .OPERATOR_ID`) | **ada** | bukan larangan menyetujui klaim sendiri (K9 / tiket 01 gugur, RALAT) |
| S7 `KomiteCount == 1 && @LengthOfPageList(KomiteList) = 1` | **ada** | `=` tunggal dibaca pembandingan (§5 butir 2) |
| S7.1 tambahkan anggota DEGREE > 1 ber-`LIMIT_BOTTOM` < total (urut DEGREE) | **ubah** | **dihitung saat tingkat 1 membuka kasus** (tampil di List of Committee, nol tulisan saat GET) dan **disimpan saat tingkat 1 Submit** (`T_KOMITE_KOMITELIST` lewat `repository/tangga.go`, KCF-02) |
| S8 `KomiteLoop` := jumlah tangga | **ada** | ditulis ulang di Submit tingkat 1 |

## 4. Layar `ShowTransfer`

| Section | Status | Catatan |
| --- | --- | --- |
| LS1 "CLAIM COMMITTEE -" + ADJUSTMENT / REJECT / CLOSE | **ada** | dari `TRANSFER_TYPE` (2 / 3 / 4) |
| LS2 CLAIM No. (VIS `1=2`), LS22 / LS23 / LS51 / LS52 | **tidak** | tidak pernah tampil |
| LS3–LS9 Policy Detail(s) (`VIS IsTreatyIn == 0`) | **ada** | periode `StartDateTime "-" EndDateTime` dari `OfferFacIn.PolicyData` `[inferensi]` (`pyWorkCover.Policy` tanpa penulis lain) |
| LS9 Object Detail + `SpreadingDetail` (expandPane), Spreading Policy (VIS `!IsTravel`) | **ada** | Spreading Claim LS40 / LS44 (VIS `1=2`) **tidak** |
| LS12–LS21 Claim Details | **ada** | LS20 / LS21 bergerbang `StsKatastrofe` / `NonKatastrofeType` VERBATIM; LS50 "Claim No" = NoClaim " / " ClaimNo |
| LS24 grid adjustment KMT + `DetailAdjustmentFac` (expandPane) | **ada** | Claim Adjustment (PT 1 / 2 / 5) / Adjuster Fee (4 / 6) / Salvage (3), Spreading Adjustment (`!IsTravel`), BreakDown Quota Share (`SpreadingQuotaShare(1).TreatyName != ''`), Payable To / rekening (`SwiftCode != ''`) |
| `DetailAdjustmentFac` "View Retro" (spreading `TreatyType = 10015`) → local action `ShowRetro` (`PreShowRetro_Act`) + `ShowSecurityReinsurer` (`PreSecurityReas_Act`) | **ada** | jendela pop-up di layar; rincian Security Reinsurer dapat dibuka |
| LS28 History Adjustment, LS33 Total Adjustment (+ "Total in IDR") | **ada** | hanya TT2 |
| LS38–LS39 teks komite (Chronology / Extent / Legal / Remarks) | **ada** | TT2: `DataCommitteFacin.*` adjustment. TT3 / TT4: Legal / Chronology / Extent = teks pop-up tersimpan di kepala kasus komite (migrasi 643, keputusan OQ-KCFI-03), Remarks = `Remark` klaim (langkah 2); urutan LS39 |
| LS3–LS9 / LS13–LS14 untuk TT3 / TT4 | **ada** | grid berulang atas `TempObjectList` / `DataTempAdj` yang hanya diisi `SetValueKomite` (TT2) → Policy Detail, Object Detail, Claim Details tidak tampil (sama dengan Pega) |
| LS41 "List of Committee" (`.TransferType == 2`) | **ada** | termasuk anggota perluasan yang belum tersimpan; Status Waiting / Approved / Reject |
| LS45 isian: `AcceptStatus` (wajib), "Propose To Close Case" / "Propose To Reserved" (VIS TT2, NA `KomiteCount != '1'`), Note (wajib) | **ada** | `AcceptStatus` dua tombol pilihan (pola Komite Prop / Non Prop yang disetujui), label + nilai `1` / `2` |
| LS48 NA `... || pyWorkPage.Adjustment.AcceptedNo != ''` | **ada** | Submit nonaktif bila adjustment KMT sudah bernomor |
| Cancel / Submit | **ada** | Submit kosong → "Value cannot be blank" per isian; 403 bukan pemegang; 409 kasus / klaim tertutup atau berubah; 422 pesan isian |
| Tata letak (ubin ringkasan, kartu berjudul, pasangan label-nilai; tangga + kartu keputusan **di bawah** semua rincian) | **ubah** | `[tidak ada di korpus]`, keputusan work owner 09-10 (Komite Prop); kulit `--kcfi-*` / `komiteclaimfacin__` |
| "View more details" | **ubah** | klaim induk Claim Fac In dibuka hanya-baca (`onLihatBerkas('claimfacin', klaim)` → `GET /api/claim-fac-in/kasus/{id}?lihat=1`) |
| Ikon grid `pzPegaDefaultGridIcons` | **tidak** | |

## 5. `KomitePostAct` → `KomitePost_Adjustment` (TT2)

| Langkah | Status | Catatan |
| --- | --- | --- |
| KomitePostAct S2 / S4 / S5 (S3 TT1 survey `//`) | **ada** | `TRANSFER_TYPE` 2 / 3 / 4 |
| S1 operator ≠ tingkat berjalan → galat | **ubah** | 403 (`ErrBukanPemegang`) di lapisan layanan |
| S2 Obj-Open-By-Handle + Lock klaim induk | **ada** | `kontrak.KlaimFacInKomite.Kunci` di transaksi Submit; klaim tertutup → 409 + rollback |
| S3–S5 kronologi "Accepted by " / "Rejected by " + IDKomite + " - " + KMT | **ada** | `T_VIEW_SUGGEST` lewat kontrak (`Riwayat`); IDKomite = jabatan tingkat |
| S6–S7.2.1.1 indeks adjustment ber-`KomiteNo` = KMT yang belum diakseptasi | **ada** | satu adjustment per KMT (OQ-CFI-28) |
| S7.2.1.2.4–S7.2.1.2.7 nomor akseptasi `KODE_PRODUKSI(NONLIFE) + "A" + BusinessOldId + "." + MM.YYYY + "." + urut5` | **ada** | `inti/backend/penomor` di transaksi Submit; tanpa `.TP` / `.TX`; 21–22 karakter (gerbang Kasir `CLM`). S7.2.1.2.8–14 `//` |
| S7.2.1.3 / S7.2.1.4 `ComiteeClaim` adjustment | **tidak** | tanpa kolom; tangga dibaca dari tabel komite (K7) |
| S7.2.1.5 tingkat akhir setuju: AcceptanceStatus 1, AcceptedNo, AcceptedDate, `Notes := Comment` | **ada** | kolom milik komite (`Kolom.MilikKomite`) lewat kontrak `UbahAdjustmentKomite` |
| S7.2.1.6 tolak: AcceptanceStatus 2 | **ada** | |
| S7.2.1.7 tingkat akhir: `IsApproved := AcceptStatus`, `AktifButton := 0` | **ada** | |
| S7.2.1.8 keputusan baris tangga (`KomiteAproval`, `KomiteComment`, `DateApprove`) | **ada** | `KOMITE_OPERATORID` ditimpa akun pemutus (pola Komite Prop 09-10-2026) |
| S7.2.1.9–S7.2.1.9.1.1 tolak: sisa baris menunggu → "2" (tanggal, tanpa komentar) | **ada** | S7.2.1.9.1 transisi ke label tak ada → tanpa lompatan (§5 butir 3) |
| S7.2.1.10–S7.2.1.14, S10–S11, S26–S27 Obj-Save / Refresh-And-Lock / UpdateWorkObject | **ubah** | satu transaksi Submit (§5 butir 6) |
| S7.2.1.15 `SendEmailKlaim_KMT` (IsPEGAPROD) | **ada** | outbox `email-komite`, hanya produksi: tingkat berikut disetujui → semua anggota workbasket KomiteID berikut; tingkat akhir → pembuat "(Approval)"; tolak → pembuat "(Reject)". BCC pribadi tidak disalin; badan = stream tidak diekspor (OQ-KCFI-01) |
| S7.2.1.16 `SaveAcceptation_KMT` (`IsPrintAccept == ""`) | **ada** | `IsFacretro` objek / item / adjustment dari SpreadingList `TreatyType = 10015`, `DLAStatus = 0` |
| S7.2.1.17 + S8 `SaveAccept_ACT` + `SaveOSClaim_SQL` → `OS_AKSEPTASI_KLAIM` | **ada** | STS **4** bila PaymentType 1, selain itu **1**; `STS_DLA` 8 retro / 7; `DATA_JSON` per lini (Fire / Aneka / Golf, MarineCargo + ObjectItemName, MBU, Travel, PA) format `GetPageJSONString` (DEV); **sekali per KMT** (§5 butir 5). Tanpa procedure: INSERT kolom procedure (`TGL_PROD` trigger); status konversi hanya ikut INSERT bila terisi (TT3 S17, keputusan OQ-KCFI-08, `repository/osakseptasi.go`) |
| S9 `PrintPDFAccep_MultiAksep_KMT` | **ubah** | gerbang `OutputData.START_DATE == ""` sesudah nomor terisi tak pernah benar → diterbitkan bila nomor terisi (§5 butir 4). Lini Fire / Aneka / Golf / Marine Cargo: stream `AcceptanceNotePDF` (diberikan work owner 10-10-2026) disusun `models/akseptasi_pdf.go`, digambar go-pdf/fpdf `models/pdf_akseptasi.go`, diunggah SESUDAH commit (`services/dokumen.go`: InsertGoogleStorage_Act Folder "Claim" Durasi 1800, lalu T_STORAGE_IMAGE + DOCUMENT_CLAIM satu transaksi, KATEGORI_1 "AcceptanceNote", MIME "pdf", nama "Persetujuan   <ClaimNo> AcceptNo <AcceptedNo>.pdf"). Gagal → keputusan tetap, `info` `PesanDokumenGagal`. Lini MBU / Travel / PA → stream tidak diekspor (OQ-KCFI-01, `info`). Rincian §8 |
| S12 tingkat akhir: objek `IsPrintAccept := @if(AcceptStatus == "2", "1", "")`; S12.1.1 adjustment `IsPrintAccept := 1` (`.AcceptanceStatus = "1"`) | **ada** | `=` tunggal dibaca pembandingan (§5 butir 2) |
| S13 (pre=false, setiap tingkat) `IsCloseFile` / `IsReservedClaim` ← Propose tingkat 1 | **ada** | disimpan di kepala kasus komite (`KOMITE_USUL_TUTUP` / `_CADANG`), ditulis balik ke kepala klaim. `IDObjectUpdate` **tidak** (tanpa kolom / pembaca) |
| S14 tolak: `KomiteCount := KomiteLoop` | **ada** | |
| S15 `FlagOnGoingCommitte := "Adjustment"` | **tidak** | tanpa kolom / pembaca |
| S16 `InsertJsonClaimNonMBU_act` → `JSON_KLAIM` | **ada** | INSERT bila IDPEGA belum ada, tanpa `DATA_JSON` |
| S17 `KonversiKlaim_Act` STSREJECT "1" (bukan Fac Retro) | **ada** | outbox `konversi-klaim`, hanya produksi |
| S18–S19 `InsertLogServiceClaim` "AKSEPATSI" (bukan Fac Retro) | **ada** | `MONITORING_KLAIM_LOG`; `JN_SERVICE` VERBATIM (literal data, tidak diperbaiki) |
| S20–S21 `HISTORYAKSEPTASIPEGA` ACCEPT / REJECT, Workbasket "KLAIM" | **ada** | setiap keputusan |
| S22–S23 `SUBPROGRESSCLAIM.POSITION2` "Accepted" / "Rejected" | **ada** | tingkat akhir atau tolak; `IDPEGA` = KMT |
| S24 (pre=false) `KomiteCount + 1` | **ada** | selalu (akhir = loop + 1) |
| S25 (pre=false) `HitServiceToKasirKMT_Act` (`DirectToKasir == "true" && StatusKasir == ""`) | **ada** | S3 `getStatusKonversi_Act` (COUNT `reinsurance.trloss_detail_t` nomor tanpa titik, hanya produksi) bertransisi pasca `.StatusKonversi=="1"`, selainnya keluar sebelum S12 IDOfBank; jalur `CLM`, jumlah per `AcceptedNo`; outbox `kasir`, hanya produksi; berjalan juga untuk Fac Retro (pre=false — keputusan OQ-KCFI-02: ikut XML); dedupe OQ-CFI-26. ⚠️ Konversi S17 = outbox (asinkron): saat Submit status konversi biasanya belum "1", kasir dikirim kemudian lewat tombol Acceptation Claim Fac In (pola Komite Claim Prop) |
| S28 ASMForceCaseClose | **tidak** | ter-remark; kasus komite ditutup di END flow |
| S29 `SetDataForInformation_Act` | **tidak** | halaman informasi tanpa pembaca |

## 6. `KomitePost_Reject` (TT3) dan `KomitePost_CloseClaim` (TT4)

Satu tingkat `ReasClaimDeptHead` (KCF-03).

| Langkah (Reject / Close) | Status | Catatan |
| --- | --- | --- |
| S1 / S1 Obj-Open-By-Handle + Lock (pre=false) | **ada** | kontrak `Kunci` |
| S2 Quotation | **tidak** | halaman salinan tanpa pembaca |
| S3–S5 / S2–S4 kronologi "Accepted by " / "Rejected by " | **ada** | |
| S6 / S5 keputusan baris tangga | **ada** | `ClaimComitee(<LAST>)` klaim induk **tidak** (tanpa kolom) |
| S7 / S6 `ApprovalCommite` | **tidak** | bahan stream PDF yang tidak diekspor |
| S8–S10 / S7–S9 Obj-Save, UpdateWorkObject | **ubah** | satu transaksi Submit |
| S11 / S10 `InsertJsonClaimNonMBU_act` | **ada** | tanpa syarat |
| S12 (setuju) `SaveReject_ACT_KMT` per estimasi `PrintFaceClaim == 1` → OS STS 2 | **ada** | Reject saja; `STS_KONVERSI` "1" bila Value kosong (S17), `STS_DLA` 8 / 7 dari `PlaStatus` |
| — / S12.1–S12.3 (setuju) OS STS 4 (`CauseOfLoss`, `CauseOfLossID`, `NoClaim`, `IDMasterTreaty`) | **ada** | Close saja |
| S13.4–S13.9 / S11.4–S11.9 PDF "Reject Claim " + NoClaim / "Close Claim " + NoClaim + NoClaim, kategori CloseClaim | **ubah** | nama berkas VERBATIM; **berkas tidak dibuat** (stream `CommitteReject_CC` / `CommitteCloseClaim` tidak diekspor, OQ-KCFI-01, `info` `InfoTanpaStream`) |
| S13.10 / S11.10 email ke pembuat "(Approval)" / "(Reject)" Pengajuan Reject / Close Klaim | **ada** | outbox, hanya produksi; ejaan bulan VERBATIM ("Febuari", "July"); BCC pribadi tidak disalin |
| S14 (setuju) `CLAIMREJECTED` | **ada** | Reject saja; LABEL "CLM", OBJCLASS kelas klaim, STATUSWORK "New" (status sebelum ditutup), REMARK = `Remark` klaim |
| S15 / S12.4 (setuju) `KonversiKlaim_Act` STSREJECT "2" / "4" | **ada** | outbox, hanya produksi |
| S16 / S13 `KomiteCount + 1` | **ada** | |
| S17 / S14 (setuju) `pxForceCaseClose` klaim induk | **ada** | **Resolved-Rejected** / **Resolved-Completed** lewat kontrak `Tutup` (status daftar putih); `CloseAllSubCases=true`: kasus komite KMT- lain klaim itu yang masih terbuka ikut ditutup berstatus sama (`[inferensi]` status sub-kasus = status induk), kasus yang memutus ditutup modul ini di END |

## 7. Arah perbaikan prompt §5 — `[penyimpangan sadar]`

| # | Rule | Dibangun | Uji |
| ---: | --- | --- | --- |
| 1 | `SetProteksiSubmiteKomite` L2.1 (baris menunggu mana saja, pemetaan akun, "IT Developer") | pemutus = anggota workbasket / akun **tingkat berjalan**; tingkat 1 juga anggota `ReasClaimSPVB` | `handlers/alur_test.go` `TestWewenangDanValidasi`, `TestTT2PitaSPVB`; `models/models_test.go` `TestPemegangTingkatBerjalan` |
| 2 | `ApprovalKomite_Act` L7, `KomitePost_Adjustment` L12.1.1, L25.2.1.1 `=` tunggal | pembandingan | `TestTT2SPVATingkatAkhirTanpaPerluasan` (`IsPrintAccept`), `TestTT2PerluasanTanggaTampilLaluTersimpan` |
| 3 | `KomitePost_Adjustment` L7.2.1.9.1 lompat ke label tak ada | tanpa lompatan | `TestTT2TolakMenutupSisaTangga` (sisa tangga "2", HISTORY / SUBPROGRESS / kronologi sesudahnya tetap berjalan) |
| 4 | `KomitePost_Adjustment` L9 gerbang PDF tak pernah benar | diterbitkan bila nomor terisi | `TestTT2SPVATingkatAkhirTanpaPerluasan` (DOCUMENT_CLAIM AcceptanceNote + isi PDF), `TestDokumenAkseptasiLiniTanpaStream`, `TestDokumenAkseptasiGagalKeputusanTetapTersimpan` |
| 5 | `KomitePost_Adjustment` L8 OS dengan adjustment terakhir | sekali per KMT, adjustment KMT itu | `TestTT2SPVATingkatAkhirTanpaPerluasan` (tepat satu baris OS) |
| 6 | Obj-Save per iterasi | satu transaksi Submit | `TestKlaimTertutupMenolakSubmitUtuh` |
| 7 | `SendRejectClaimToKomite2` / `SendCloseClaimToKomite` anggota tunggal akun + email tertulis mati | `ReasClaimDeptHead` | claimfacin `handlers/komite_test.go` `TestRejectClaimTT3MelahirkanKomiteTanpaAdjustment`, `TestCloseWithoutPaymentTT4` |

## 8. Penyimpangan lain

- Penjaga permintaan ganda TT3 / TT4 (`models.PesanTutupKomiteGanda`, claimfacin): `[tidak ada di korpus]`, menolak
  kasus komite TT3 / TT4 kedua selagi yang pertama menunggu.
- `SalinCatatanTutup` (claimfacin) mengisi `Remark` dan `Remark_Close` dari pop-up Remarks: `Remark` = REMARK
  `CLAIMREJECTED` dan Remarks layar komite.
- Tabel komite inbox Claim Fac In: kolom Committee No, Claim ID, Claim No, Type (ADJUSTMENT / REJECT / CLOSE),
  Committee, Level, Status, Updated `[tidak ada di korpus]` (pola Non Prop, bahasa Inggris).
- Tanggal Boleh Bayar Kasir (`HitServiceToKasirKMT_Act` S10.5 / S14.1.5) **ditiru VERBATIM**, sama dengan Claim Fac In
  tahap 1: hari > 25 → tanggal 01 dua bulan sesudahnya; bulan 13 / 14 → "1"; tahun naik hanya bila bulan hasil 01 dan
  bulan berjalan Desember.

- **Dokumen akseptasi (S9, stream `AcceptanceNotePDF`)** — pemetaan medan stream → halaman `TempAcceptedNo`
  (`PrintPDFAccep_MultiAksep_KMT` S1-S22) di `models/akseptasi_pdf.go`; uji `models/pdf_test.go`:
  - `[penyimpangan sadar]` urutan (prompt tahap 2 §6 butir 7: "seperti pola Komite Prop sesudah commit"): Pega
    membuat PDF di tengah Submit; di sini SESUDAH keputusan tersimpan (pemutusan klien tidak membatalkannya), dan
    kegagalan PDF - termasuk panik penggambar - tidak membatalkan keputusan. Akibatnya saringan S10.3.11
    `.IsPrintAccept == ""` tidak dipakai (S12 sudah menyetelnya 1): adjustment dipilih menurut status 1 + nomor
    akseptasi.
  - `[penyimpangan sadar]` `HTMLToPDF` (mesin HTML Pega) → go-pdf/fpdf yang MENGGAMBAR susunan stream (judul, baris
    label : nilai, tabel di sel nilai, penutup; kepala "PT. REASURANSI NUSANTARA MAKMUR", kaki "Page n of N pages"):
    teks label VERBATIM, ukuran dari CSS stream; lebar kolom tabel tanpa `width` = isi terpanjang; penutup rata kanan;
    nilai lebih panjang dari satu halaman berlanjut di halaman berikut; huruf cp1252 (é, –, ’) dicetak.
  - `[penyimpangan sadar]` angka `DecimalFormat("#,###.####")` id_ID atas nilai DESIMAL tersimpan, bukan `double` Java
    (nol float untuk uang): beda hanya pada seri tepat ...5 di desimal ke-5 dan nilai di atas 15 digit bermakna.
  - PERBAIKAN kelainan XML (maksud pasti): S10.3.6 / S10.3.9 `6 || 4 && U` dibaca `(6 || 4) && U`; `TemporaryUang`
    segar per cetak (Pega tak pernah membersihkannya), bertambah selama iterasi seperti S10.3.1; Property-Remove di
    dalam For Each = saringan biasa; lini dinilai atas halaman KLAIM (Pega menilai kasus komite tanpa `OfferFacIn`);
    "Swift Code" hanya bila terisi.
  - VERBATIM (maksud tidak pasti): "Location of Loss" berisi CauseOfLoss (S14); nominal, Payable To, dan tabel
    spreading hanya untuk item TERAKHIR yang lolos, AdjustmentGross kumulatif; tabel "BreakDown Spreading (QS)" tanpa
    penulis tingkat item - tidak tercetak; "Premium Paid On" kosong (PaymentData tidak disimpan, OQ-CFI-14); ID
    DOCUMENT_CLAIM `yyyyMMddhhmmssSSS` jam 12-an.

## 9. Sambungan Claim Fac In (`claimfacin`)

| Hal | Isi |
| --- | --- |
| Kontrak | `inti/backend/kontrak/klaimfacin.go` `KlaimFacInKomite` (Baca / Kunci / TulisBalik / Tutup), disediakan `services.KlaimUntukKomite` |
| Daftar putih kepala | `AktifButton` (S7.2.1.7), `ClaimData.IsCloseFile` / `ClaimData.IsReservedClaim` (S13) |
| Daftar putih adjustment | `AcceptanceStatus` (S7.2.1.5 / S7.2.1.6), `AcceptedNo` / `AcceptedDate` / `Notes` (S7.2.1.5), `IsApproved` (S7.2.1.7), `IsFacRetro` (SaveAcceptation_KMT), `IsPrintAccept` (S12.1.1), `StatusKasir` / `IDOfBank` (HitServiceToKasirKMT_Act S12–S13) |
| Daftar putih objek / item | objek `IsPrintAccept` (S12), `DLAStatus` / `IsFacretro` (SaveAcceptation_KMT); item `IsFacretro` |
| Status penutupan | `Resolved-Rejected` (TT3 S17), `Resolved-Completed` (TT4 S14) |
| Kolom katalog baru | adjustment `NOTES` (milik komite), kepala `IS_CLOSE_FILE`, `IS_RESERVED_CLAIM` — kolom tabel bersama yang sudah ada |
| Yes `SureRejectClaim` → `SendRejectClaimToKomite2` (TT3) | **ada** — galat VERBATIM (akseptasi sudah ada / estimasi belum Face Claim), kasus `KMT-` tanpa adjustment `TRANSFER_TYPE` 3, kronologi "Request Reject claim " + KMT, label "Success Create Request to Committee", email outbox |
| Yes `PreventRejectClaim` (CWP) → `SendCloseClaimToKomite` (TT4) | **ada** — galat VERBATIM (adjustment belum diputus), `TRANSFER_TYPE` 4, kronologi "Request close claim without payment " + KMT |
| Hak | `GET /api/claim-fac-in/hak` → `komite` |
| Tabel komite inbox | `frontend/components/TabelKomite.tsx` |
