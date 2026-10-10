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
> Cek layar 07-10-2026 (server tiruan `UJI-`, vite port 5177, Chrome tanpa kepala): halaman awal (tab Process /
> Resolve, switch Teknik - mati = worklist sendiri + Add Claim, nyala hanya bagi anggota ReasKlaimTeknik; rupa kelas
> inti seperti Kelola User),
> layar kasus 08-10-2026 berkulit Kelola User (kartu `panel`, `form-grid`, kotak hanya-baca; tema terang dan
> gelap; kasus baru tanpa daftar pesan galat; kulit neumorfik `.kelola-user` disalin ke token `--cp-*` untuk inbox,
> layar kasus, dan popup; susunan layout = format layout XML: kepala tengah, Claim Treaty Inline grid double 8|5,
> Claim Information dua kolom, Quarter/Year sebaris, tab Interest / Estimation / Spreading, Claim History paging 5),
> layar Outstanding Claim kasus baru (54 medan, nol galat konsol), popup "Data Master TreatyIn" — tombol tampak
> cocok dengan baris berstatus di bawah untuk keadaan kasus baru. Tahap Input Acceptation sampai Resolved-Completed
> dijaga uji HTTP `TestAlurPenuhSampaiResolved`.

## 1. Flow dan halaman awal

| XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| `Flow_TreatyIn` Start1 → Assignment2 "Outstanding Claim" (worklist pembuat) | pembuatan kasus | 01 | **dibangun** — tab *Process* bawaan (switch Teknik mati): worklist pembuat tanpa cek workbasket, XML apa adanya; tombol **Add Claim** hanya saat switch Teknik mati (keputusan 08-10-2026; OQ-CP-13: harness `New` tidak diekspor) |
| Assignment1 "Input Acceptation" (workbasket) | FlowAction `InputAcceptation` | 10 | **dibangun** — tab *Process*, switch *Teknik* nyala = workbasket `ReasKlaimTeknik` (keputusan 07-10 dan 08-10-2026); switch nonaktif bagi akun tanpa workbasket itu (`GET /hak`) |
| Decision3 `IsBackStage` → Assignment2 / End1 Resolved-Completed | — | 11 | **dibangun** — tab *Resolve*; jalur kembali lewat "Back" yang di XML `NEVER` → tidak dibangun |
| FlowAction `OutstandingClaim` pre-act `CheeckNoRNM_Act`, pre-DT `SetDateOutstanding`, post-act `ProteksiData_act` | — | 01, 02 | **dibangun** (`siapkan`; tanggal hanya diisi bila kosong — penyimpangan, OQ.md; pesan pra-proses termasuk ProteksiData langkah 12 TIDAK tampil saat kasus dibuat / dibuka — keputusan work owner 08-10-2026, pesan tampil sesudah aksi) |
| FlowAction `InputAcceptation` pre-act `GetPICAdjutment_Act`, pre-DT `SetDateAcceptation` | — | 08 | **dibangun** |

## 2. Section `OutstandingClaim` (Assignment2)

| Label XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| Choose Master | harness `MasterTreatyIn` (dis `IsOutstanding==1`) | 01 | **dibangun** (popup, tombol Choose per baris = `SetValueToClaim_Act`); RD `BrowseCLAIM_MASTER_TREATY` + parameter section `TREATYTYPE = "Proportional"` → `PROPORTIONTYPE = 'Proportional'`; filter per kolom (AND) di server, terbaru dulu, batas 500, 50 per halaman — penyimpangan dari `pyMaxRecords 20` atas keputusan work owner 08-10-2026 |
| Summary Outstanding Claim | harness `SummaryOutSClaim`, pre-act `GetDataOustanding` (vis `IsOutstanding = 1`) | 07 | **dibangun** |
| teks `.Message` (klik) | harness `InputTreatyInOffer` (tidak diekspor) | 01 | **nonaktif-OQ** OQ-CP-02 — tampil hanya-baca |
| Choose Policy No | harness `ListPolicyNoTreaty_Harness`, pre-act `SetMasterID` | 01 | **dibangun** (pilih = klik ganda baris → `CheckNoPolicy`, seperti XML) |
| View (polis) | harness `DetailPolisRealization` (tidak diekspor); pre-act `GetDetailPolis_act` mengurai dokumen polis ke kelas `ASM-FW-GISFW-Work` | 01 | **dibangun lain** — berkas NB / EDM Treaty In (PRODKE terbesar, `GET /berkas-polis`) tampil di jendela (modal selebar layar) di atas layar klaim, tanpa menu dan tanpa tab baru; polis Pega lama tanpa berkas = pesan (keputusan work owner 08-10-2026) |
| View Status Payment Premi | harness `ListPaymentPremi_Harness`, pre-act `GetDtlPaymentPremi_act` | 02 | **nonaktif-OQ** OQ-CP-03 |
| Choose Cause of Loss | harness `CauseofLoss_Harness` (dis `IsAnyAcceptation =1`) | 03 | **dibangun** |
| ikon tambah Consultant / Appointed Adjuster | harness `MstAdjusterConsultant` | 08 | **dibangun** — popup tambah (Name, Telp No, Address; label XML) disimpan lewat rute pinjaman modul Adjuster Consultant `POST /api/adjuster-consultant`, ID baru diisi ke medannya; nonaktif selama ID hanya-baca (keputusan work owner 08-10-2026) |
| ikon edit RNM Share % | `GetRNMShareTreaty` (vis `EstimationList(1).PrintFaceClaim = ''`) | 06 | **dibangun** |
| ikon "Standard Icon" RNM Share % | tanpa aksi, vis never | — | tidak tampil di XML |
| Save | `SetOutstanding_Act` lalu save (vis `IsAcceptation!=1`) | 07 | **dibangun** |
| Save to issue RNM | `SaveOutstanding_Act` → LA `PrintFile` bila `IsPLA=1` (dis `IsCFS = ''`) | 07, 12 | **dibangun** (PrintFile = pesan "Please Print Pla") |
| PRINT PLA | LA `GeneratePLA` (dis `IsPLA!=1`) | 12 | **dibangun** — nomor PLA; berkas PDF OQ-CP-05 |
| Send to Acceptation | `CheckNopolicy_Act` (vis `IsOutstanding = 1`, dis `IsAcceptation==1`) | 01 | **dibangun lain** — CheckNopolicy_Act lalu Submit dalam satu aksi: berkas langsung pindah ke `ReasKlaimTeknik`; validasi Submit gagal = seluruhnya batal (perintah work owner 08-10-2026 "langsung kirim ke teknik, ga usah klik submit lagi") |
| Submit | finishAssignment `OutstandingClaim` bila PolicyNo dan NoClaim terisi (vis `IsAcceptation==1`) | 10 | **dibangun** — pindah ke workbasket `ReasKlaimTeknik`; kini hanya tampil untuk berkas yang sudah ber-IsAcceptation tetapi belum pindah tahap (dikirim sebelum 08-10-2026 / data lama) |

### Sub-section `Catastrope_Sec`, `OutstandingClaim_Intrs`, `OutstandingClaim_Est`, `OutstandingClaim_Sprd`

| Kontrol XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| ikon pensil / centang Catastrophe | `SetEditCatastrope` Edit / Save | 03 | **dibangun** |
| ikon gear Catastrophe | LA `CatastrofeList` | 03 | **dibangun** (popup: tambah = `InputCatastrope`, Save = `SaveCatasrtope_Act`, pilih = `SetCatastrope_act`, Cancel) |
| tambah / hapus baris Interest | `AddInterest_act` / `DeleteInterest_act` (dis `IsAdjVal=='Yes'`) | 05 | **dibangun** |
| tambah / hapus baris Total Claim Amount | `AddListClaimAmount` / `DeleteListClaim_Act` (dis `Note=='Yes'`) | 09 | **dibangun** |
| tambah / hapus baris Loss Allocation | `AddLossAllocation_act` / `RemoveLossAlloction_act` (dis `IsOldData=='Yes'`) | 06 | **dibangun** |
| tambah / hapus baris Estimation | `AddEstimation_Act` / `DeleteEstimation_Act` (dis `PrintFaceClaim==1`) | 07 | **dibangun** |
| tambah baris Spreading Claim | `AddSpreading_Act` (**tidak diekspor**) | 06 | **dibangun lain** — baris terisi dari polis saat CheckNoPolicy (TREATYINPRODUCTION JN_REAS / PCT_SHARE_PREMI / CURR_ID); Add aktif: satu klik = satu baris per mata uang yang belum ada untuk treaty spreading polis pertama yang masih kurang; Treaty Type LANGSUNG terisi dan hanya-baca; spreading sama (Treaty Type + Currency) ditolak, juga saat memilih di baris kosong; tanpa spreading polis baris kosong, Treaty Type dipilih (keputusan work owner 08-10-2026) |
| hapus baris Spreading Claim (ikon baris) | `DeleteSpreading_Act` (**tidak diekspor**, dis `IsOldData='Yes'`) | 06 | **dibangun lain** — baris dihapus, nonaktif bila `IsOldData='Yes'` (XML); tabel bawah disusun ulang (keputusan work owner 08-10-2026) |
| ikon grid standar Spreading List (tambah / hapus bawaan) | `pzPegaDefaultGridIcons` ("always"), tanpa activity | 06 | **tidak dibangun** — tambah / hapus spreading dinonaktifkan (keputusan work owner 07-10-2026) |
| SpreadingBreakQS (`SetTreatyNameSpreading_Act` langkah 11) | `TreatyInMaster.Limits(1).Detail(1).SpreadingList` | 06 | **dibangun** — master (XML) didahulukan; PROPORTIONALARRG sebagai pengganti dibatalkan sesudah uji data DEV (keputusan work owner 08-10-2026); SpreadingList master KOSONG -> anak PROPORTIONALARRG setiap TreatyType tabel atas (treaty group klaim, tahun arrangement terbaru <= tahun treaty; keputusan work owner 09-10-2026 "kalo tidak ada di master, ambil dari proportionalarrg"). DEV 09-10-2026: 786 / 1.856 M_TREATY_IN dan 237 / 281 M_TREATY_IN_EDM tanpa SpreadingList |
| dropdown Treaty Type spreading | `Spreading.pxResults` (penulis tidak diekspor); XML Outstanding: pxTextInput `.TreatyName` | 06 | **dibangun lain** — dropdown TreatyType di kedua layar = spreading polis + treaty baris yang sudah ada, terkunci bila terisi atau data lama; koreksi `[dugaan]` master SpreadingList berdasar data DEV (induk vs anak) (keputusan work owner 08-10-2026) |
| tombol tanpa label di samping RNMShareP | vis NEVER | — | tidak tampil di XML |
| grid "Claim History" (SuggestList) | — | 14 | **dibangun** — `T_VIEW_SUGGEST` (keputusan work owner 07-10-2026) |

## 2a. Tampilan medan (keputusan work owner 08-10-2026, berlaku juga di section 3)

| Medan | XML | Status |
| --- | --- | --- |
| Consultant ID / Adjuster / Professional ID | Section OutstandingClaim / InputAcceptation: medan ID + baris nama hanya-baca | **dibangun lain** — dropdown memilih dan menampilkan NAMA (label "Consultant Name" / "Adjuster / Professional Name", saringan NAME saja); ID tetap disimpan (CONSULTANT_ID / ADJUSTER_ID); baris nama hanya-baca XML hanya tampil saat `IsAnyAcceptation = 1` |
| ikon Add di samping Consultant / Adjuster | `pyWorkActionsAddWork.png` | **dibangun** — selalu menempel ke medannya, juga saat baris nama tersembunyi |
| isian angka (kendali angka, medan dan sel grid) | Currency / Decimal | **dibangun lain** — hanya angka, separator Indonesia (titik ribuan, koma desimal), paling banyak 4 desimal saat mengetik dan saat tampil; nilai tersimpan tetap mentah |
| isian tanggal / tanggal-waktu | DateTime (pemilih bawaan) | **dibangun lain** — diketik `dd-mm-yyyy` (+ `hh:mm`), pemisah otomatis, tanggal tidak sah ditandai merah, tombol kalender; aksi server (CheckDateDOL dll.) hanya saat isian lengkap dan sah atau dikosongkan |
| Reporter Address (`GetReportStatus_Act` -> `GetAddressCeding`) | `M_CLIENT.JSONDATA.AddressList.*` klien milik agen ceding | **dibangun lain** — dibaca dari tabel datar CLIENT_ADDRESS: baris Kantor (tipe 2) dulu, tanpa Email (tipe 7) ("ubah jangan dari json, ambil dari client address"); DEV: setiap bagian alamat yang terisi di kedua sumber sama, 2 agen tambahan mendapat alamat |
| Type (Acceptation List / panel, `ASM-FW-GCNMFW-Data-Adjustment.Type`), Deductible Type (`.IndividualRiskType`), Payable To (`Data-ClaimData.Payable`) | prompt values tidak ada di korpus | **dibangun** — dari screenshot work owner 08-10-2026: Type 1 Claim / 2 Adjuster Fee / 3 Salvage / 4 Consultant Fee; IndividualRiskType 0 Select.. (label saja) / 1 % From claims / 2 % FromTSI / 3 Other; Payable 1 Ceding Co Name / 2 Broker Name / 3 Others; nilai tersimpan kode |
| Type (grid Estimation List, `ASM-FW-GCNMFW-Data-Estimasi.Type`) | prompt values tidak ada di korpus | **dibangun** — dari screenshot work owner: 1 Claim, 2 Adjuster Fee, 3 Salvage, 4 Consultant Fee; nilai tersimpan kode |

## 3. Section `InputAcceptation` (Assignment1) dan sub-section-nya

| Label XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| View Status Payment Premi (#1) | aksi kosong, vis NEVER | — | tidak tampil di XML |
| Summary Outstanding Claim | harness `SummaryOutSClaim`, pre-act `GetDataOutsClaim_act` | 07 | **dibangun** |
| View Status Payment Premi (#2, #3) | harness `ListPaymentPremi_Harness` | 02 | **nonaktif-OQ** OQ-CP-03 |
| Close Claim | LA `PreventRejectClaimProp` | 11 | **dibangun** |
| Choose Policy No | vis NEVER | — | tidak tampil di XML |
| View (polis) | harness `DetailPolisRealization` (tidak diekspor); pre-act `GetDetailPolis_act` mengurai dokumen polis ke kelas `ASM-FW-GISFW-Work` | 01 | **dibangun lain** — berkas NB / EDM Treaty In (PRODKE terbesar, `GET /berkas-polis`) tampil di jendela (modal selebar layar) di atas layar klaim, tanpa menu dan tanpa tab baru; polis Pega lama tanpa berkas = pesan (keputusan work owner 08-10-2026) |
| Choose Cause of Loss | harness `CauseofLoss_Harness` | 03 | **dibangun** |
| ikon tambah Consultant / Appointed Adjuster | harness `MstAdjusterConsultant` | 08 | **dibangun** — popup tambah (Name, Telp No, Address; label XML) disimpan lewat rute pinjaman modul Adjuster Consultant `POST /api/adjuster-consultant`, ID baru diisi ke medannya; nonaktif selama ID hanya-baca (keputusan work owner 08-10-2026) |
| Save (di bawah Est) / Save (di bawah Adjs) | save | 07, 08 | **dibangun** |
| Back | `BackToRegister_act` (tidak diekspor) — layout NEVER | — | tidak tampil di XML — prompt §6 butir 9 menyebutnya, tetapi "tampil sesuai section" = tidak tampil; activity-nya tetap OQ-CP-01 |
| blok "Claim Estimation" tab Acceptation (`InputAcceptation_Adjs`: Count Claim Amount, Loss Allocation, Estimation List, total estimasi, Spreading Claim - hanya-baca) | tampilan saja | — | **tidak dibangun** — dibuang dari tab Acceptation (perintah work owner 08-10-2026 "semua data estimasi yang muncul di tab akseptasi di hapus"); data tetap di tab Estimation. Acceptation List dan Spreading Adjustment Total tetap |
| Update Estimation / Save (`InputAcceptation_Adjs`) | `UpdateEstimasi_Act` (tidak diekspor) — layout NEVER | — | tidak tampil di XML — idem; OQ-CP-01 |
| View Status Payment Claim | harness `ListPaymentClaim_Harness` | 13 | **nonaktif-OQ** OQ-CP-03 |
| View Payment Attachment | harness `ViewAttachment` | 13 | **nonaktif-OQ** OQ-CP-03 |
| tambah baris Adjustment | `AddAdjustment_Act` (vis `AktifButton = '0' \|\| ''`) | 08 | **dibangun** |
| hapus baris Adjustment | `DeleteAjsutment_Act` (dis `IsKomite=1 \|\| IsSubjectivity`) | 08 | **dibangun** (+ tolak baris ber-`KOMITE_ID`) |
| Save to issue RNM / PRINT PLA (`InputAcceptation_Est`) | sama dengan Outstanding (dis `ReCFS != 1` / `IsPLA!=1`) | 07, 12 | **dibangun** |
| klik baris Adjustment | expand pane `AdjustmentDetail` | 08 | **dibangun** (panel rinci baris; HANYA grid Adjustment - grid lain tidak punya expand pane, laporan work owner 08-10-2026). Baris terbaru terbuka tanpa klik (sejak awal dan sesudah Add), klik di mana pun pada baris membuka / menutup ("tidak harus klik angka sebelah kiri"). Isi panel mengikuti layout XML + screenshot Pega: tab Detail Chronology to Committe / Adjusment List Detail / Subjectivity, Inline grid double bersarang, checkbox "Transfer Direct to Kasir" (caption XML), tabel bebas 6 x 4, Payable \| bank berdampingan; grid Spreading In tampil di atas Spreading Out walau XML ContainerVisibleWhen NEVER (Pega hidup menampilkannya, screenshot work owner 08-10-2026) |
| klik baris Total Claim Amount | expand pane `ViewDetailInterest` (`InterestListDtl` — tidak diimpor, J1) | 05 | tidak dibangun: halaman kerja bukan data (STRUKTUR J1) |

## 4. Section `AdjustmentDetail` (panel rinci baris adjustment)

| Label XML | Aksi XML | Tiket | Status |
| --- | --- | --- | --- |
| Save | save (layout `view.CARI21!=1`) | 08 | **dibangun** |
| Send to Committe | `AttachmentProtect_ACT` → harness `CommitteeTreaty` bila `Protect.CARI1/2 = 1` (vis `TotalKomite!=''`, dis `IsKomite==1 \|\| IsError>1`) | 11 | **dibangun** — gerbang lampiran membaca cacah dokumen klaim per kategori master PROP (OQ-CP-12 ditutup 08-10-2026) |
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
| Send Claim to Committee | `AddKomiteTreatyChild_ACT` (dis `Payable = ''`) | 11 | **dibangun** (opsi "b" 07-10-2026) — kasus TKMT- + tangga + email komite; terjangkau sesudah gerbang lampiran (uji `TestKirimKomiteSesudahLampiranLengkap`: LOD + DLA + SPGR -> TKMT-). Penyerahan ULANG baris subjectivity (S16-S19 / S31, roster jenjang terbawah) dibangun 08-10-2026 (OQ-KCP-06 "a") |
| Cancel | tutup | 11 | **dibangun** |

## 6. Local action modal

| FlowAction | Isi | Tiket | Status |
| --- | --- | --- | --- |
| `PrintFile` / `PrintFileDLA` | label "Please Print Pla" / "Please Print DLA" | 12 | **dibangun** sebagai pesan info |
| `GeneratePLA` | `.ClaimData.Remark`, post-act `TryMakePLA_Act`; tombol `pySubmitLabel` "Submit" / `pyCancelLabel` "Cancel" | 12 | **dibangun** |
| `GenerateDLATreaty` | `.RemarksDLA`, post-act `PrintDLATreatyIn`; Submit / Cancel | 12 | **dibangun** |
| `PreventRejectClaimProp` | Yes (tanpa bayar → `SendCloseClaimToKomite`) / No / Yes (`CloseClaimProp`) / No / Close | 11 | **dibangun** — Yes tanpa bayar dibangun 10-10-2026 (OQ-CP-06 selesai): `SendCloseClaimToKomite` melahirkan kasus komite TT 4 satu tingkat (roster PROP "Claim Dept. Head"), kronologi + `Remark_Close` + riwayat "Request close claim without payment", klaim tetap terbuka sampai komite memutuskan. Penyimpangan sadar: (1) CWP kedua selagi kasus TT 4 klaim itu masih terbuka ditolak "Close without payment for this claim is already waiting for committee decision" `[tidak ada di korpus]`; (2) `CloseClaimProp` ditolak server bila centang tanpa bayar menyala (XML hanya menyembunyikan tombolnya); (3) Close biasa ikut menutup kasus komite TT 4 yang masih terbuka (`CloseAllSubCases`) |
| `CatastrofeList` | lihat §2 | 03 | **dibangun** |
| `MessageBeforeDeleteTreatyGroup` | dibuka tombol Delete yang `1==2` | — | tidak tampil di XML |

## 6a. Lampiran dan efek keluar (bukan tombol section, tetapi aksi XML)

| Rule XML | Dipicu | Tiket | Status |
| --- | --- | --- | --- |
| section `pyCaseAttachmentsWrapper` (direfresh sesudah PRINT PLA) | lampiran kasus | 11, 12 | **dibangun sebagai tab "Lampiran"** — susunan VERBATIM screenshot layar Pega work owner 09-10-2026: judul + (jumlah lampiran), tombol Add attachment / Refresh / Save (Save = tombol Save layar, aksi `Simpan`, aktif hanya bila layar menawarkannya), grid Category / Count Attach / Upload File / View File, 5 baris + "Show All" (judul tanpa jumlah - work owner 09-10-2026). Add attachment = tabel File / Category per berkas (`.pyCategory`, bawaan "File") + hapus pilihan, berkas dipilih atau diseret-lepas (work owner 09-10-2026); Upload File baris = kategori baris untuk semua berkas (TempInputParam.pyCategory). Jendela View File = screenshot layar Pega work owner 09-10-2026 (pilih Category "Label - (jumlah)", File Name / Category / Create Date `DD/MM/YYYY HH24:MI` / centang Action + centang-semua / Delete, tombol Download Selected, Delete Selected, Change Category = pindah dokumen terpilih ke kategori lain dalam satu transaksi) ditambah Attached By (PXCREATEOPERATOR; kolom No Acceptation / No Prekas dibuang atas perintah work owner) dan View (pdf / gambar) / View Office Online (xls...pptx) di popup penampil layar penuh pola NB Treaty In; Delete / Delete Selected / Change Category hanya pemegang assignment (Delete: objek penyimpanan lalu catatan + baris dokumen satu transaksi); "Move to" `[tidak ada di korpus]`; isi lewat `GetBase64Attachment` -> `GetUrlGoogleStorage_Act` (Durasi 1800). Section tidak diekspor (widget bawaan Pega); isi mengikuti `GCNMSaveAttachments` (diekspor 08-10-2026): kategori sekali untuk semua berkas (TempInputParam.pyCategory), satu `InsertDocument_Act` per berkas, IDPEGA = klaim; S1.1 / S1.3-S1.6 ter-remark (Link-Attachment Pega) tidak dibuat. Unggah hanya pemegang assignment; teks panel `[tidak ada di korpus]` |
| `GCNMSaveAttachments` -> `InsertDocument_Act` | S1.2 kategori / InsKey, S1.7; InsertDocument S2-S3 (ID `yyyyMMddhhmmssSSS` jam 12-an, MIME huruf kecil), S4 InsertGoogleStorage_Act (Folder Claim, Durasi 1800), S5 Obj-Save bila T_STORAGE_ID terisi | 11 | **dibangun** — `services.UnggahLampiran`: unggah lalu T_STORAGE_IMAGE + baris dokumen klaim satu transaksi (pola Bordereaux); jenis berkas dijaga `penyimpanan` (daftar lama GetMimeType) |
| REST `SendAcceptationToKasir` | Acceptation (`HitServiceToKasir_Act`) | 13 | **dibangun** — efek outbox "kasir" hanya bila `IS_PEGA_PROD`; panggilan nyata OQ-CP-03 |
| REST `KonversiKlaimNonLife` | `KonversiKlaim_Act` (Save to issue RNM, Close Claim) | 13 | **dibangun** — efek outbox "konversi-klaim" hanya bila `IS_PEGA_PROD`; OQ-CP-03 |
| REST `GetDtlPaymentClaim` | View Status Payment Claim | 13 | **nonaktif-OQ** OQ-CP-03 |
| REST `getPayAttachment` | View Payment Attachment (`GetPayAttachment_Act` → `GetPayAttachmentAdj_Act`) | 13 | **nonaktif-OQ** OQ-CP-03 |
| REST `getPremiumPaidOnTreatyIn` | View Status Payment Premi | 02 | **nonaktif-OQ** OQ-CP-03 |
| REST `ServiceGoogle` | `InsertGoogleStorage_Act` / `GetUrlGoogleStorage_Act` (berkas) | 11, 12 | InsertGoogleStorage_Act **dibangun** lewat `inti/backend/penyimpanan` (lampiran); GetUrl (unduh) belum - OQ-CP-05 |
| email `SendEmailKlaim` | `AddKomiteTreatyChild_ACT` 34 | 11, 13 | **dibangun** — efek outbox "email-komite" hanya bila `IS_PEGA_PROD`; CC/BCC OQ-CP-15 |
| email `SendEmailKlaimRejectClose` | `SendCloseClaimToKomite` (tutup tanpa bayar) | 11 | **dibangun 10-10-2026** — efek outbox "email-komite" (MUATAN `jenis` = `tutup`) hanya bila `IS_PEGA_PROD`; pengiriman nyata sama dengan `SendEmailKlaim` (OQ-CP-03 / OQ-CP-15) |

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
