# PARITAS — Komite Claim Prop ↔ Pega

Dibuat 08-10-2026 (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7 butir 11). Setiap isian / tombol
Section `ShowTransfer` (wajah TT 2) dan setiap langkah `KomitePost*` dipasangkan dengan tiket dan statusnya.

Status: **dibangun** · **nonaktif-OQ** (tampil, tidak dapat dipakai, menunggu ekspor / keputusan) · **di luar
lingkup** (alasan disebut) · **mati** (blok `//` / prakondisi tak terjangkau, tidak dibangun).

Bukti = berkas korpus `D:\XML\RNM_BRD\Komite Claim Prop` + rule + langkah.

## 1. Alur dan rute

| Rule | Langkah | Sistem baru | Tiket | Status |
| --- | --- | --- | --- | --- |
| `Flow/KomiteTreaty_Flow.xml` | Start → assignment "KomiteRouter" → `ViewTransferDtl` | daftar kerja + layar kasus | 02 · 03 | dibangun |
| `Flow/KomiteTreaty_Flow.xml` | Decision `KomiteLoop` (`When/IsKomiteLoop`) | `models.MasihBerjalan`; salah → `STATUS_WORK = Resolved-Completed` | 06 | dibangun |
| `Activity/KomiteRouter.xml` | S6 / S6.1 `AssignTo` = baris `KomiteList` pertama ber-`KomiteAproval == 0` | `repository/tangga.go` `sqlSaringKerja` (MIN `KOMITE_URUT` menunggu), `models.Kasus.Giliran` | 02 | dibangun |
| `Activity/KomiteRouter.xml` | S1-S5, S7 | — | — | mati (`//`) |
| kelas `Work-KomiteTreaty` | `pyCanCreateWorkObject=false` | nol tombol "buat baru"; kasus dilahirkan Claim Prop (opsi B 07-10-2026) | 01 | dibangun |

Kasus komite tanpa baris tangga menunggu tidak tampil di daftar kerja siapa pun (KomiteRouter tanpa cadangan); cacahnya
`repository.CacahTanpaGiliran` — DEV 08-10-2026: **0** (DEV belum punya kasus `TKMT-`).

## 2. Layar `ShowTransfer` (flow action `ViewTransferDtl`)

Pra-proses `SetKomiteList_Act` (TT 2): total per mata uang tanpa baris `AcceptanceStatus` 2 + baris "Total in IDR" —
`models.TotalPenyesuaian` (tiket 09, dibangun).

| Sel / blok (VERBATIM) | Sumber | Status |
| --- | --- | --- |
| Judul "CLAIM COMMITTEE -" + "ADJUSTMENT" | label ber-`pyCondition .TransferType =='2'` | dibangun |
| Judul "CLOSE" (TT 4) / "REJECT" (TT 3) | `pyCondition` TT 4 / TT 3 | di luar lingkup — TT 4 OQ-CP-06; TT 3 tanpa penulis (Claim Prop hanya menulis TT 2 `AddKomiteTreatyChild_ACT` dan TT 4 `SendCloseClaimToKomite`; `3` hanya dibaca `SendEmailKlaimRejectClose`) |
| No Claim di kepala "Claim Treaty", "sample text" | `pyVisible` NEVER | mati |
| "Claim Analysis" (label Heading 2) | — | dibangun |
| Class of Business · Ceding · Source of Business · Policy No · Policy No Ceding · Insured Name | `pyWorkCover.*` | dibangun |
| `.StartDateTime` - `.EndDateTime` | label = deskripsi properti (aturan properti tidak diekspor) | dibangun — label tampil nama properti (OQ ekspor) |
| Cause of Loss · Report Type | Report Type dropdown `associated` → label work owner (Claim Prop `LabelKode`) | dibangun |
| Adjuster / Professional · Consultant | `pyCondition` tidak kosong | dibangun |
| Location of Loss | — | dibangun |
| "Policy No" ber-nilai `.Adjustment.Type` (`pyVisible ALWAYS`) | dropdown → label `SetDataAcceptationTreaty_Act` S4-S7 | dibangun (VERBATIM, walau labelnya ganjil) |
| Catastrophe · `.NonKatastrofeType` (bila Non-Catastrophe) · Catastrophe Note (bila Catastrophe / Big Claim) | `pyCondition` | dibangun |
| Tombol "View more details" | `showHarness` `ViewClaimFormKomite` (diekspor 08-10-2026: tab Registration / Estimation / "Adjustment & Acceptation" / "Policy Detail & Claims History" atas `pyWorkCover`; section-nya tidak diekspor) | dibangun — berkas Claim Prop klaim induk dibuka hanya-baca di jendela di atas layar komite (`PropsRute.onLihatBerkas`); menu Claim Prop tidak dipegang akun → pesan |
| PIC Name · Date | `pyWorkPage.pxCreateOpName` / `pxCreateDateTime` (kasus komite) | dibangun |
| `.NoClaim` / `.CLMNO` | NOTBLANK | dibangun |
| PLA / DLA Number · Date of Loss · Report Date · Received Date · Reporter Name · Reporter Phone Number · Reporter Status · Specify... · Reporter Address · Report Description | — | dibangun |
| RNM Share (%) | `.Adjustment.PersenRNM`, 4 desimal | dibangun |
| Grid "Insured Interests 100 %" · "Count Claim Amount" · "Loss Allocation" · "Estimation List" | `ClaimData.*` | dibangun |
| "Total Original Currency Estimation" (grid) | TT 2 | dibangun — kolom Currency di XML dapat disunting tetapi tidak pernah disimpan → hanya-baca |
| "Total Estimation In IDR": Total Gross Estimate(100%) in IDR · Total Estimation in IDR | 2 sel lain `pyVisible` NEVER | dibangun |
| "History Adjustment" | `TempClaimData.AdjustmentList` | dibangun — Payment Type / Komite No dapat disunting di XML, tidak disimpan → hanya-baca; Payment Type tampil kode (sumber `BrowseReinsuranceType_RD` .Note = OQ); Status label `AcceptanceStatus.xml` (0 Transfer to committee / 1 Approve / 2 Reject), juga kolom Status daftar kerja |
| "Total Adjustment" | `TempTotalAdj.pxResults` | dibangun |
| Grid "Spreading Claim" | `ContainerVisibleWhen 1=2` | mati |
| Blok deductible ("Dedutible Type", 6 × 4) | `.Adjustment.*` | dibangun |
| "Spreading In" · "Spreading Out" | `.Komite.SpreadingAdjustment` / `SpreadingQuotaShare` = baris adjustment | dibangun — Share / Claim Spreaded Spreading Out dapat disunting di XML, tidak disimpan → hanya-baca |
| Payable To · Specify · Name of Bank · Swift Code (bila terisi) · Branch of Bank · Account No | `.Adjustment.*` | dibangun (Payable label `Payable.xml`: Ceding Co Name / Broker Name / Others) |
| Circumstances · Occupation · Salvage · Adjuster / Consultant Fee (NOTBLANK) · Remarks | `.Komite.*` ← `DataCommitteeTreaty` baris adjustment + `ClaimData.Occupation` | dibangun |
| Grid "Committe Accept Status" (Committe Name, Status, Date Approve, Comment) | tangga | dibangun (Status label `KomiteAproval.xml`: Waiting / Approved / Reject) |
| "Are you sure to accept this document?" `.AcceptStatus` | wajib; `pyNoSelectionText` "Choose"; pilihan `AcceptStatus.xml` 1 "Approve" / 2 "Reject" | dibangun |
| "Subjectivity ?" `.IsSubjectivity` | tampil `.AcceptStatus = 1 && TT 2`; nonaktif `.KomiteCount!='1'` | dibangun — isian tingkat 1 disimpan di header (`KOMITE_SUBJECTIVITY` / `_NOTE`, migrasi 682; OQ-KCP-01 dijawab "a" 08-10-2026) dan dipakai tingkat akhir |
| "Subjectivity Note" | tampil + wajib `.IsSubjectivity = true`; nonaktif `KomiteCount!='1'` | dibangun — dropdown `SubjectivityNote.xml` (1 Treaty Leader Approval … 7 Others), kode disimpan; server menolak kode di luar daftar (teks pesan bukan VERBATIM: validasi tabel bawaan Pega tidak diekspor) |
| "Propose To Close Case" · "Propose To Reserved" | tampil TT 2; nonaktif `KomiteCount!='1'` | dibangun |
| "Note" `.Comment` | wajib | dibangun |
| Tombol "Cancel" | `cancel` | dibangun (kembali ke daftar kerja) |
| Tombol "Submit" | `finishAssignment`; `pyDisabledWhen pyWorkPage.Adjustment.AcceptedNo != ''` | dibangun — juga ditegakkan server (409) |
| Flow action `ViewDetailInterest` | pintu masuk tak terlacak | di luar lingkup — butir terbuka 9 |

⚠️ Penyimpangan sadar layar:
- Nilai `pyWorkPage.Adjustment` / `.Komite` di Pega = SALINAN saat penyerahan; di sini dibaca dari baris adjustment
  yang sama (beku selama diserahkan, ditegakkan Claim Prop).
- Isian awal setiap tingkat = nilai `pyWorkPage` tersimpan (keputusan terakhir, komentar tingkat sebelumnya, dua usul),
  seperti Pega yang tidak mengosongkannya.

## 3. `KomitePost` → `KomitePostAdjustment` (Submit, satu transaksi)

| Langkah | Isi Pega | Sistem baru | Tiket | Status |
| --- | --- | --- | --- | --- |
| KomitePost S1 | TT 2 → `KomitePostAdjustment` | `services.Putuskan` | 04 | dibangun |
| KomitePost cabang TT 3 / TT 4 | `KomitePost_Reject` / `KomitePost_Close` | — | — | di luar lingkup (§1 tabel judul) |
| S1 | `SetDataAcceptationTreaty_Act` (teks Approve / Reject, Type) | label pilihan + data dokumen | 03 · 11 | dibangun (teks) |
| S2 | Page-Remove | — | — | tidak perlu |
| S3 | prakondisi NONAKTIF (`KomiteCount != local.operator`) - langkah selalu jalan dan hanya menyetel variabel lokal | — ; wewenang pemegang tingkat berjalan ditegakkan layanan (403) atas keputusan 30 / ADR-0014, bukan atas S3 | 04 | dibangun |
| S4 | Obj-Open-By-Handle klaim induk, Lock | `KunciKlaimTreaty` (kontrak) + `FOR UPDATE` kepala komite | 04 | dibangun |
| S5 | variabel lokal | — | — | tidak perlu |
| S6 | `ComiteeClaim(count)` induk + `KomiteList(count)` komite | SATU baris `T_KOMITE_KOMITELIST` (`tangga.go` `TulisAnggota`) | 04 | dibangun |
| S6 / S7 | S6 dilewati dan S7 menulis `ComiteeClaim(<LAST>)` bila adjustment induk SUDAH subjectivity | `Rencanakan` selalu S6 | 04 · 05 | dibangun (OQ-KCP-06 dijawab "a" 08-10-2026): Claim Prop menyerahkan ulang baris subjectivity (`BolehSerahKomite`, `KOMITE_ID` ditimpa dari kasus lama ke kasus baru; grid ComiteeClaim memuat semua putaran - S17 tidak menghapus saat subjectivity); komite menulis baris TERAKHIR tangga (S7) dan `.Comment` awal = komentar anggota pertama putaran pertama (AddKomiteTreatyChild_ACT S16) |
| S8-S10 | `InsertChronology_DT` "Accepted by" / "Rejected by" + jabatan | kontrak `Riwayat` → `T_VIEW_SUGGEST` | 04 | dibangun — nama yang di-hardcode dibuang (keputusan 17-09); gerbang baris 1 `OperatorID.pyPosition != "IT Developer"` TIDAK dibangun: posisi operator Pega tidak punya padanan di akun sistem baru (sama dengan Claim Prop) — penyimpangan sadar, OQ-KCP-04 |
| S11 | IsCloseFile / IsReservedClaim ← dua usul | kontrak header | 05 | dibangun |
| S12 | tolak → EXT (S25) | `models.Rencanakan` | 07 | dibangun |
| S13-S14 | TotalKomite; `IsAnyAcceptation := 1` | kontrak header | 06 · 08 | dibangun |
| S15 | `IsApproved` ← `ComiteeClaim(<LAST>).KomiteAproval` | kontrak adjustment | 08 | dibangun |
| S16.1-S16.4 | `GETTanggalClosing_SQL`, `GenerateNoAcceptTreaty` | — | — | mati (`//`) |
| S16.5-S16.8 | kode NONLIFE + "A", `GetSequenceNumber_SQL`, rakit nomor (`OutputData.START_DATE`) | `repository/nomor.go` (penomor) + `models.RakitNomorAkseptasi` | 08 | dibangun — celah XML: `ParamSeq.CARI3` (tanggal) tidak pernah diisi; ALL_SOURCE `PROC_GENERATE_SEQUENCE_NUMBER` dengan tanggal NULL = SYSDATE Jakarta + pergeseran hari tutup buku = `penomor.HitungPeriodeNomor`, SAMA dengan `UrutNomor` Claim Prop (nol beda perilaku) |
| S16.9 | AcceptedNo / AcceptanceStatus 1 / AcceptedDate | kontrak adjustment | 08 | dibangun |
| S17 | `SaveAcceptation_Act` | `OS_AKSEPTASI_KLAIM` (+ DATA_JSON) + IsOutstanding 1 / IsCFS "" | 10 | dibangun — baris `TempOpenPage.AkseptasiList(<APPEND>)` (S1, 11 properti) TIDAK ditulis: Claim Prop tidak punya tabel untuk daftar itu (keputusan sesi Claim Prop, daftar hanya sasaran pesan); isinya ada di baris OS yang sama |
| S18-S20 | Obj-Save / `UpdateWorkObject` | transaksi Submit | — | dibangun (`UpdateWorkObject` tidak diekspor) |
| S21 | `SaveAcceptationTreaty_TKMT` | tahun treaty, batas RP, retro, FacRetroList, IsFacRetro, IsPrintAccept | 10 · 11 | dibangun (PDF = OQ) — S7 `TempOpenPage.Message` "Please Print DLA" milik halaman klaim dan transien di Claim Prop (tanpa kolom); Claim Prop menampilkannya sendiri dari IsFacRetro di local action PrintFileDLA, tidak dikembalikan ke penyetuju |
| S22 | `AktifButton := 0` | kontrak header | 06 | dibangun |
| S23 | subjectivity → `IsKomite := 0` | kontrak adjustment | 05 | dibangun |
| S24 | IsSubjectivity / SubjectivityNote ke adjustment + header | kontrak | 05 | dibangun |
| S25 | EXT: KomiteCount := Loop, AcceptanceStatus 2, AktifButton 0 | kepala komite + kontrak | 07 | dibangun |
| S26 / S26.1 | sisa tangga menunggu → 2 (+ tanggal) | `TulisAnggota` tanpa komentar | 07 | dibangun |
| S26.1 / S27 `FlagOnGoingCommitte` | — | — | dibuang (keputusan 26-28) |
| S27 | Notes ← Comment | kontrak adjustment | 06 · 07 | dibangun |
| S28 | `InsertJsonClaimTreaty_act` | `JSON_KLAIM` INSERT bila IDPEGA belum ada, tanpa DATA_JSON | 10 | dibangun |
| S29 | `KonversiKlaim_Act` (Connect-REST, IsPEGAPROD) | outbox `konversi-klaim` (hanya produksi) | 12 | dibangun — pelaksana berhenti `…BelumDisetujui` |
| S30-S31 | log "AKSEPTASI" | `MONITORING_KLAIM_LOG`: NO_AKSEPTASI = `OutputData.START_DATE` (nomor S16.8); status / respons REST kosong (panggilan kini asinkron) | 10 | dibangun |
| S32-S33 | `InsertHistoryAkseptasiPega_Sql` | `HISTORYAKSEPTASIPEGA` (ACCEPT / REJECT, KLAIM, ID_KOMITE; OPERATORID NULL seperti Pega) | 10 | dibangun |
| S34 | `HitServiceToKasirKMT_Act` | §4 | 12 | dibangun sebagian |
| S35 | `SendEmailKlaim_KMT` (IsPEGAPROD) | outbox `email-komite`; MUATAN hanya pengenal (jenis, ID akun penerima, ID baris tangga yang diputuskan - claimlife/015 melarang nama / alamat). Isi dirakit saat dikirim (`services.SusunEmailKomite`) | 12 | dibangun — badan `EmailKlaim_HTML_KMT` VERBATIM (`models/templat`), subjek, penerima S12 `KOMITE_EMAIL` / S13-S15 `M_LOGIN_GO.EMAIL`, akun NUSARE / NUSARESYARIAH (S17-S18), CC S4 (`konfigurasi/email.json`); BCC pribadi S3 tidak disalin; kirim SMTP berhenti `ErrEmailBelumDisetujui` |
| S36-S37 | Obj-Save + `UpdateWorkObject` | `TulisBalikKlaimTreaty` (kontrak) | 04 | dibangun |
| S38 | Page-Remove | — | — | tidak perlu |
| S39 | `ASMForceCaseClose` | — | — | mati (`//`) |
| S40 | `KomiteCount + 1` (prakondisi nonaktif) | kepala komite | 06 | dibangun |
| S41 | Commit | commit transaksi aplikasi | — | dibangun |

⚠️ **Penyimpangan sadar — satu transaksi** (AC 56 diralat): di Pega sembilan rule menyimpan sendiri (`COMMIT` di RDB)
lalu S41 Commit. Di sini keputusan + tangga + nomor + tulis balik + tabel warisan + log + antrean outbox = SATU
transaksi aplikasi; Submit yang gagal di tengah batal utuh. Panggilan keluar baru berjalan sesudah commit (outbox) —
urutan efek keluar sengaja diubah (keputusan 14, 20).

## 4. Activity yang dipanggil

| Activity | Langkah | Sistem baru | Status |
| --- | --- | --- | --- |
| `SaveAcceptation_Act` | S1 TempOSAkseptasi (18 kunci + pxObjClass), S3 parameter, S4 DATA_JSON, S5 `SaveOSClaim_SQL`, S6 IsOutstanding / IsCFS | `models.SusunOSAkseptasi`, `repository.SisipOS` | dibangun — DATA_JSON diisi (keputusan work owner 08-10-2026 sore); kolom datar = kunci bernama sama; `CARI16/17/20` NULL |
| `SaveAcceptation_Act` S2 | Type 4 "final" | — | mati (`//`) |
| `SaveAcceptationTreaty_TKMT` | S1-S9 | `services.jalan.retro` + `models.SusunRetro` | dibangun (S4: putaran TERAKHIR `SpreadingAdjustment` yang tersisa di LimitDla / ListInsurerDla) |
| `PrintFileAcceptance_TKMT` | S1-S4 IsPrintAccept, BusinessOldId | penanda ditulis | dibangun |
| `PrintFileAcceptance_TKMT` | S5-S10 halaman `TempAcceptedNo` + stream `FILEAcceptanceNote` (diekspor 08-10-2026); S9 nama berkas "Persetujuan Klaim   AcceptNo <no>.pdf", kategori AcceptanceNote | outbox `dokumen-akseptasi` (MUATAN: ID akun penyetuju + saat cetak); `services.SusunDokumenAkseptasi` + `models.SusunAcceptanceNote` | dibangun — markup VERBATIM; S8 Location ditimpa CauseOfLoss (XML apa adanya); tanggal `format="date"` = dd/MM/yyyy `[penyimpangan sadar]` |
| `PrintFileAcceptance_TKMT` | S11 `HTMLToPDF`, S12 Base64, S13 `InsertDocument_Act` (Google Storage folder Claim + `DOCUMENT_CLAIM`) | pelaksana outbox | **OQ-KCP-07** — mesin PDF Pega tidak punya padanan (go.mod tanpa pustaka PDF); pelaksana berhenti `ErrPenyimpananBelumDisetujui` |
| `InsertJsonClaimTreaty_act` | S1-S4 | `SalinJSONKlaim` | dibangun (`GetBase64Attachment` S1 tidak dipakai: DATA_JSON kosong) |
| `KonversiKlaim_Act` | S1-S3 | outbox | dibangun |
| `InsertLogServiceClaim` | S1-S2 | `CatatLogLayanan` | dibangun |
| `HitServiceToKasirKMT_Act` | S2 gerbang DirectToKasir & StatusKasir kosong; S3 `getStatusKonversi_Act` | S2 + S3 (hanya IsPEGAPROD) | dibangun |
| `HitServiceToKasirKMT_Act` | S12-S13 IDOfBank (BANKACCOUNT) | kontrak adjustment | dibangun (hanya produksi: di luar produksi S3 keluar) |
| `HitServiceToKasirKMT_Act` | S14.1 jalur CLMP: panjang AcceptedNo 23/24, email ceding, muatan | `models.SusunMuatanKasir`, kode tetap `konfigurasi/kasir.json` | dibangun; IsPEGASyariah = OQ (false) |
| `HitServiceToKasirKMT_Act` | S14.4-S14.6 REST, StatusKasir, `DIRECTTOKASIR_LOG` | pelaksana outbox | OQ-CP-03 (berhenti `ErrKasirBelumDisetujui`) |
| `HitServiceToKasirKMT_Act` | S9-S11, S14.2 (CLMNP / CLM) | — | di luar lingkup (lini lain) |
| `HitServiceToKasirKMT_Act` S1 | salvage | — | mati (`//`) |
| `SendEmailKlaim_KMT` | S1 tanpa SpreadingAdjustment keluar; S5 tanggal ("Febuari", "July" VERBATIM); S9-S15 Temp.CARI*; S11 baris spreading (`@divide(SharePercentage,1,0)`, ClaimSpreaded `#,###.####` id_ID); S10 total (= `CountSpreadingADJ_Act` S6.3 / S10, Claim Prop tidak menyimpannya); S16 stream; S17-S18 akun | `models.SusunDataEmail` / `RenderEmailKomite`, `services.SusunEmailKomite` | dibangun — Ceding / SOB dibaca dari `TreatyInMaster` (S7 `ClaimData.QuotationData.CedingCoName` / `SobName` diisi dari master yang sama dan tidak disimpan) |
| `getStatusKonversi_Act` | IsPEGAPROD → REINSURANCE.TRLOSS_DETAIL_T | `repository.Acuan.StatusKonversi` di S3 Submit | dibangun — hanya IsPEGAPROD, DEV tidak dijalankan (keputusan work owner 08-10-2026); hak baca diminta DBA. Konversi S29 asinkron: di produksi status bisa belum "1" saat Submit → Kasir lewat tombol "Acceptation" Claim Prop |
| `GetEmailCeding_SQL` | `gl.f_get_email` | `repository.Acuan.EmailCeding` | dibangun — DEV ORA-00904 (OQ DBA) |
| `CountEstimation_Act`, `CountSpreading_act`, `CurencyEstimation_Act`, `SetCurencyList_act`, `AddEstimation_Act`, … (berkelas Work-ClaimTreaty) | sel hanya-baca ShowTransfer | — | di luar lingkup (spec Out of Scope 2) |
| `SaveRejectTreatyIn_Act_KMT`, `KomitePost_Close`, `KomitePost_Reject` | TT 3 / TT 4 | — | di luar lingkup |

## 5. Data lama (tiket 13)

Kasus komite warisan **tidak dimigrasi** (keputusan work owner 08-10-2026, OQ-KCP-02 "b"): klaimnya dimuat pemuat Claim
Prop - baris berlaku `STS_REJECT = 1` kini → Input Acceptation (uji-kering DEV 08-10-2026: 2.451 kasus siap, 0 gagal,
0 ditunda). Sensus sebelumnya: 419 klaim ber-STS 1 = 2 ber-BLOB KomiteTreaty, 316 hanya riwayat akseptasi, 101 tanpa
jejak komite. Alat pemuat komite dibuang; AC 83 diralat.

## 6. Butir terbuka modul ini

| Kode | Butir | Pemilik |
| --- | --- | --- |
| ~~OQ-KCP-01~~ | DIJAWAB 08-10-2026 "a": migrasi 682 menyimpan isian Subjectivity tingkat 1 | — |
| ~~OQ-KCP-02~~ | DIJAWAB 08-10-2026 "b": kasus komite lama tidak dimigrasi; klaim dimuat Claim Prop | — |
| OQ-KCP-03 | Hak baca `REINSURANCE.TRLOSS_DETAIL_T` dan eksekusi `GL.F_GET_EMAIL` untuk akun aplikasi produksi (work owner 08-10-2026: "ikuti"; bacaan bergerbang IsPEGAPROD, DEV tidak dijalankan) | DBA |
| OQ-KCP-04 | Ekspor yang masih kurang: deskripsi properti `StartDateTime` / `EndDateTime` / `NoClaim` / `CLMNO` / `NonKatastrofeType`, padanan `OperatorID.pyPosition`. Diterima 08-10-2026: harness `ViewClaimFormKomite`, stream `FILEAcceptanceNote` / `EmailKlaim_HTML_KMT`, prompt values `AcceptStatus` / `KomiteAproval` / `AcceptanceStatus` / `Payable` / `SubjectivityNote`. Harness `Confirm` DITUTUP: harness bawaan platform Pega (bukan aturan aplikasi, maka tidak ada di ekspor) - layar sesudah Submit kembali ke daftar kerja | pemilik ekspor Pega |
| ~~OQ-KCP-05~~ | DIABAIKAN work owner 08-10-2026: 166 baris OS status 1 DEV berbentuk lain | — |
| OQ-KCP-07 | PDF dokumen akseptasi: `HTMLToPDF` memakai mesin PDF platform Pega. Padanannya butuh pustaka Go baru (go.mod) atau layanan konversi; markup sudah dirakit, unggah lewat `inti/backend/penyimpanan` + `DOCUMENT_CLAIM` menunggu keputusan | work owner |
| ~~OQ-KCP-06~~ | DIJAWAB 08-10-2026 "a": penyerahan ulang subjectivity dibangun di Claim Prop dan komite | — |
| OQ-CP-06 | Jalur Close tanpa pembayaran (TT 4) | work owner (ditunda 08-10-2026) |
