# PARITAS — Komite Claim Non Prop

Setiap isian / tombol / langkah korpus `Komite Claim Non Prop` ↔ status di sistem baru. Dibangun 09-10-2026 atas
perintah work owner ("untuk menu ikuti menu klaim prop, dimana tidak ada menu komite klaim non prop. inbox nya di
gabung di menu klaim non prop … untuk jenjang komite juga mengikuti klaim prop, dimana sudah menggunakan workbasket"),
pilihan "Bangun sekarang": modul dibangun seperti `komiteclaimprop`; aturan tangga tetap XML (OQ-CNP-01), roster ke
workbasket lewat migrasi claimnonprop 611.

Status: **ada** = dibangun persis; **ubah** = dibangun dengan penyimpangan sadar (alasan); **tidak** = tidak dibangun.

## 1. Menu dan pintu masuk

| Korpus | Status | Catatan |
| --- | --- | --- |
| Menu `komiteclaimnonprop` (M_NAV_MENU, migrasi inti 900) | **ubah** | Menu DIHAPUS (perintah work owner 09-10-2026: "anggap menu itu tidak pernah ada"): baris M_NAV_MENU + hak akunnya dihapus migrasi inti 949; modul tanpa menu (`frontend/layar.ts`), dipasang bagi pemegang menu `claimnonprop` (`MODUL_DIPINJAM`) |
| Worklist assignment "KomiteRouter" | **ubah** | Tabel "Committee" di bawah inbox Claim Non Prop (rute pinjaman `GET /api/komite-claim-non-prop/kasus`), hanya bagi pemegang workbasket roster NONPROP aktif (`GET /api/claim-non-prop/hak` `komite`); klik baris = layar komite di tempat (`onBukaModul`) |
| KomiteRouter S6.1 `AssignTo = .KomiteID` baris pertama ber-approval 0 | **ada** | KomiteID = akun ATAU workbasket aktif pelaku (roster ke workbasket, 611) |
| KomiteRouter S6 gerbang `Primary.TransferType=='2'` | **ubah** | `TransferType` = `.Type` akseptasi (CreateChild 9; 2 = Adjuster Fee): XML hanya menetapkan tujuan untuk Type 2. Tidak ditegakkan - tujuan setiap kasus = POSITION / KomiteID tingkat berjalan (jawaban OQ-CNP-25, pola workbasket) |
| KomiteRouter S1-S5, S7 | **tidak** | ter-remark (`//`) |

## 2. Layar `ShowTransfer`

| Section | Status | Catatan |
| --- | --- | --- |
| S1 "CLAIM COMMITTEE -" | **ada** | "CLOSE" / "REJECT" (IsCloseFile / IsReject kasus komite) **tidak**: kasus komite CWP tidak dilahirkan Claim Non Prop (OQ-CNP-36) |
| S4 tombol "View Claim" (harness `ViewDtlClaimKmt`) | **ubah** | Harness tidak diekspor: berkas Claim Non Prop klaim induk dibuka hanya-baca di jendela (`onLihatBerkas`) |
| S4 "Claim Analysis", "sample text" (NEVER) | **ada** / **tidak** | |
| S6 kolom kiri (Class of Business … Type) | **ada** | `Komite.Occupation` = `ClaimData.Occupation` (CreateChild 20); Type berlabel `AdjustmentType` Claim Non Prop |
| S7 inline "Insurance Period" `StartDateTime "-" EndDateTime` | **ada** | satu sel berlabel (`pyLabelFieldValue`) |
| S8-S9 kolom kanan (PIC Name … RNM Share (%)) | **ada** | S9 inline "Claim No / Claim ID": NoClaim (NOTBLANK) "/" CLMNO; `.Adjustment.PersenRNM` kosong → RNM Share master (`TreatyInMaster.RNMShare`, penulis baris akseptasi Claim Non Prop) |
| S10 Circumstances | **ada** | `Komite.CircumtansesCouseOfLoss` = akseptasi `DataCommitteeTreaty.CircumCauseOfLoss` (CreateChild 12) |
| S11 grid ListClaimAcceptation + "Total Claim Amount" | **ada** | kaki total tanpa penulis di korpus (OQ-CNP-37) → kosong persis XML |
| S15 Loss Allocation, S19 XOL Allocation | **ada** | |
| S23 Previously Calculated | **ada** | `IsPrevious` = SetKomiteList_Act S2 (`AlokasiXOLPaid` terisi) |
| S27 Spreading In / Spreading Out | **ada** | |
| S36-S40 Payable To / rekening / rekening kedua | **ada** | `.Adjustment.Payable` / `.PayableTo` / rekening akseptasi; syarat tampil VERBATIM `IsPrevious==1` (S36) dan `FlagCurrency==1` (S40) |
| S41 "Committe Accept Status" | **ada** | kolom `Initial` kosong: SetKomiteList_Act S1 menulis inisial / jabatan per NAMA ORANG tertulis mati — dibuang, jabatan dari roster |
| S45 Remarks (ExtentOfLoss / LegalLiability NEVER) | **ada** | `Komite.Remarks` = akseptasi `DataCommitteeTreaty.Remarks` |
| S46 isian keputusan | **ada** | AcceptStatus wajib (1 Approve / 2 Reject `[inferensi]` prompt values Komite Claim Prop); Subjectivity / catatannya / dua Propose nonaktif bila KomiteCount != 1; Note wajib |
| S47-S50 Cancel / Submit (kontainer NEVER, `pyShowFAButtons` false) | **ubah** | tombol layar seperti Komite Claim Prop |
| Tata letak pola Komite Claim Prop (`[tidak ada di korpus]`, perintah work owner "ikuti tampilan klaim prop") | **ubah** | ringkasan kepala (Claim No, Policy No, Insured, Committee Level) dan lencana status; judul kartu (Claim Analysis, Loss Details, Claim Acceptation, Allocation, Spreading, Payment & Bank, Committee Notes, Your Decision); grid S41 dirender sebagai langkah tangga (jabatan, status, tanggal, komentar); kolom tabel Committee inbox / daftar kerja (Case ID, Update Date/Time, Claim, Level, Committe Name); pemegang menu Claim Non Prop dapat membuka kasus komite hanya-baca (Submit hanya pemegang tingkat berjalan) |

## 3. `KomitePostAdjustment`

| Langkah | Status | Catatan |
| --- | --- | --- |
| S1 / S31 cabang CWP (`KomitePostAdjustmentCWP`) | **tidak** | OQ-CNP-36 |
| S3-S4 pemutus bukan KomiteID | **ubah** | XML hanya memberi pesan; di sini 403 (OQ-CNP-16 bawaan, pola workbasket Komite Claim Prop) |
| S5 Obj-Open-By-Handle + Lock | **ada** | kunci kasus klaim lewat kontrak, satu transaksi |
| S6 / S7 keputusan tingkat (S7: akseptasi sudah subjectivity → baris TERAKHIR) | **ada** | |
| S8-S10 kronologi "Accepted by " / "Rejected by " + pyUserName | **ada** | |
| S11.3 tolak: baris tingkat berjalan dan sesudahnya | **ada** | komentar ikut ditulis; akun pemutus hanya di baris yang diputus |
| S11.4 AcceptanceStatus 2, S12 → EXIT | **ada** | tanpa nomor, OS, email, Kasir |
| S13 / S14.1-S14.2 flag CNPAccNo* (halaman kasus komite) | **tidak** | dibaca hanya S30 info layar |
| S14.3-S14.7, S14.14, S19.1 | **tidak** | ter-remark |
| S14.8-S14.11 nomor `kode NONLIFE + "A" + BusinessOldId + "." + MM.YYYY + ".TX" + urut5` | **ada** | penomor bersama (`inti/backend/penomor`), tanpa procedure; `OfferFacIn.QuotationData.BusinessOldId` tanpa penulis → kosong persis XML (OQ-CNP-37), akibatnya panjang nomor bukan 23 / 24 dan Kasir keluar di S9 sampai penulisnya diketahui |
| S14.12 IsCloseFile ← IsProposeClose | **ada** | "1" / "0" |
| S14.13 AcceptedNo / AcceptedDate / AcceptanceStatus 1 | **ubah** | `[inferensi]` bila AcceptedNo sudah ada (S14.9-S14.11 dilewati) nomor lama dipertahankan — XML menulis `OutNOAcc` kosong |
| S14.17 CNPStatusCase "CLAIM ACCEPTED" | **ada** | `ClaimData.IsFInalAccXOL` **tidak**: tanpa pembaca di korpus |
| S14.18 InsertOSKlaimCNP (OS + CLAIMXOL2 + JSON_KLAIM) | **ada** | S1 `GenerateAccCNP_act` (PDF) **tidak**: stream `AccClaimKomite_HTML` tidak diekspor (OQ-CNP-22); S4 CLAIMOLD akun tertulis mati dibuang |
| S14.19 `stsReject` | **tidak** | tanpa pembaca (KonversiKlaim_Act menerima STSREJECT "1" tertulis mati) |
| S14.20 GetBase64Attachment | **tidak** | |
| S14.21 / S24 JSON_KLAIM | **ada** | sekali per IDPEGA |
| S14.22 KonversiKlaim_Act | **ada** | outbox, hanya produksi; pengecualian satu akun tertulis mati dibuang |
| S16 IsKomite 0 (subjectivity) | **ada** | |
| S17 InsertOSSubjectivityCNP | **ada** | S6 InsertXOLKlaimCNP saat subjectivity = keluar di S1 (nol baris) |
| S18 IsSubjectivity / SubjectivityNote | **ada** | baris akseptasi; `TempMainWork.ClaimData.IsSubjectivity` (header) **tidak**: tanpa pembaca di korpus Non Prop, tanpa kolom di katalog Claim Non Prop |
| S19.2 SendEmailKlaim_KMT | **ada** | outbox, hanya produksi; BCC pribadi tidak disalin; penyetuju berikut ber-KomiteID workbasket = semua anggota aktifnya |
| S19.3 HitServiceToKasirKMT_Act cabang IsCLMNP | **ubah** | S7 NoAccount angka saja (ditulis balik); S9 / S10 panjang AcceptedNo bukan 23 / 24 = Exit Activity (tanpa S12-S13); S9-S10 muatan per mata uang Spreading In (Nett = TotalClaim − PremiumSpreaded); **S10.3 ber-WHEN `.TreatyName=="UR"` atas baris TempSpreadingRisk yang tidak pernah membawa TreatyName (XML: CARI1-18 tidak pernah terisi) — muatan tetap disusun, sama dengan tahap 1 Claim Non Prop `HitServiceToKasir_Act` 9.3**; S12-S13 IDOfBank; S10.8 StatusKasir / S11 IsPrintAccept bergantung jawaban REST — tidak ditulis di Submit (pola Komite Claim Prop); IsPEGASyariah `false` (OQ-CNP-40) |
| S20 KomiteCount := KomiteLoop, CNPStatusCase "CLAIM REJECTED" | **ada** | |
| S22-S23 HISTORYAKSEPTASIPEGA | **ada** | |
| S25 UpdateWorkObject, S26 Obj-Save, S27 Commit | **ubah** | satu transaksi aplikasi; efek keluar lewat outbox sesudah commit |
| S28 ASMForceCaseClose (count ≥ loop), S29 count + 1 | **ada** | POSITION dikosongkan / maju ke KomiteID berikut |
| S30 info layar | **tidak** | |

## 4. Penyimpangan lain

- Tanggal Boleh Bayar Kasir: tahun bergulir bersama bulan (perbaikan sama dengan Claim Non Prop OQ-CNP-05 butir 5).
- Daftar kerja tanpa kolom nilai / mata uang (akseptasi Non Prop bernilai per layer; `ValueAdjustment` tanpa penulis).
- Email `SendEmailKlaim_KMT`: `TotalSharePersen` / `TotalSpreadAdjustment` tanpa penulis di korpus Non Prop - `[inferensi]`
  jumlah baris Spreading In akseptasi (pola Komite Claim Prop); penulis aslinya ikut OQ-CNP-37.
- `SusunMuatanKasir` IsPEGASyariah selalu `false` (OQ-CNP-40, pola Claim Non Prop).
