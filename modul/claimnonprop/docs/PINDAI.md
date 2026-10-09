# Pindai Claim Non Prop (Klaim Treaty Non Proporsional / XoL)

> Pindai baca-saja 09-10-2026, perintah work owner "sekarang beralih ke klaim treaty non prop!" (pilihan
> "Pindai dulu, lalu laporan"). **Belum ada kode.** Pertanyaan yang perlu diputuskan ada di [`OQ.md`](OQ.md).
>
> Sumber: korpus `D:\XML\RNM_BRD\Claim Non Prop\` dan `D:\XML\RNM_BRD\Komite Claim Non Prop\` (termasuk
> `Struktur_Flow_TreatyIn.xlsx` dan `Struktur_KomiteTreaty_Flow.xlsx`), pembanding `Claim Prop` / `Komite Claim Prop`
> beserta port `modul/claimprop`, dan DEV POOLDATA lewat SELECT saja.
>
> Penanda: **[inferensi]** = disimpulkan, bukan tertulis di XML. **[DEV]** = fakta dari DEV 09-10-2026. **REMARK** =
> langkah activity berlabel `//`, tidak pernah jalan. Aturan baca layar R1–R5 sama dengan
> `modul/claimprop/docs/grilling-ronde-2.md` §13 (bab 10 di bawah).
>
> Akun orang, alamat email, kata sandi, dan alamat IP yang tertulis mati di XML **sengaja tidak dicetak** di dokumen
> ini. Penyebutannya diganti peran, misalnya "akun uji" atau "satu akun Claim Dept. Head".

## 0. Ringkasan

1. **Alurnya sama persis dengan Claim Prop.** `Flow_TreatyIn` punya dua assignment: Outstanding Claim (worklist
   pembuat) lalu Input Acceptation (workbasket `TreatyinPNCTeknik`), dengan jalur Back yang tidak terjangkau. Yang
   berbeda adalah isi layar dan hitungannya.
2. **Hitungan intinya waterfall layer XoL.**
   - Klaim 100% per mata uang dikurangi deductible dan ditambah TPL, lalu dikali share ceding.
   - Hasilnya dialokasikan per treaty (Loss Allocation, centang *To XOL*).
   - Porsi XoL dibagi: retensi (UR) dulu, lalu layer 1, 2, … sampai batas tiap layer.
   - Porsi RNM = alokasi layer × RNM Share.
   - Reinstatement premium = klaim layer / limit × MDP × % reinstatement.
3. **Tidak ada Estimation List seperti di Prop.** Layar Non Prop berisi grid Claim Amount per mata uang, Loss
   Allocation, XOL Allocation, dan Summary XOL.
4. **Satu klaim bisa punya banyak akseptasi.** Akseptasi berjenis Payment Type 1–7, termasuk interim (2) dan
   cancellation (7) [inferensi dari pesan validasi]. Tiap akseptasi membawa salinan alokasi XoL sendiri, nilai yang
   sudah dibayar ("Previously Calculated"), dan Net Claim = Total Claim − Reinstatement Premium.
5. **Penutupan klaim lewat dua tombol terpisah:**
   - **Close Claim:** langsung menutup kasus.
   - **Close Without Payment (CWP):** lewat komite satu tingkat yang di XML tertulis mati ke satu akun Claim Dept. Head.
6. **Tangga komite di XML tidak memakai kolom batas roster `EMAILKOMITE`.** Bila RNM Share ≤ 30% dan nilai
   ≤ 30.000.000, komite hanya Dept Head. Selain itu semua tingkat roster NONPROP aktif (4 tingkat). Riwayat DEV
   memperlihatkan kasus komite lama melewati 1–4 tingkat.
7. **Penyimpanan Pega hanya menambah baris**, tanpa UPDATE:
   - `OS_AKSEPTASI_KLAIM` (prosedur `PEGA_JSON_OS_AKSEP_KLAIMTNP`, satu baris per layer × mata uang per simpan).
   - `CLAIMXOL2` (prosedur `XOL2_AKSEP_KLAIM`, setiap persetujuan komite).
8. **Data warisan DEV:**
   - 1.684 kasus `CLMNP-` (4.178 baris OS, 2019–2026) dan 685 kasus di `CLAIMXOL2`.
   - Hanya 6 baris `JSON_KLAIM` dan 8 kasus di tabel kerja Pega.
9. **Ada tiga objek DEV yang bermasalah:**
   - Tabel `TREATY_OUT`, yang dibaca untuk layer PLA retro, **tidak ada**. Yang ada `TREATY_OUT2`.
   - View `CLAIMXOL`, yang dibaca "Claim History Master ID", **INVALID**.
   - `CLAIMREJECTED` berisi 102 kasus `CLMNP-`, padahal XML Non Prop tidak pernah menulisnya.
10. **Banyak hardcode di XML:** nomor kasus tertentu, nomor master tertentu, akun uji, kata sandi Edit XOL Allocation,
    premi aktual 50.000, serta nama orang komite. Semuanya perlu keputusan sebelum dibangun (`OQ.md`).

## 1. Identitas

| Unsur | Nilai | Bukti |
| --- | --- | --- |
| Work class klaim | `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp` (induk assignment `…-Work-PNC`) | `Flow/Flow_TreatyIn.xml` |
| Awalan kasus | `CLMNP-` (varian `CLMNPS-` [inferensi: syariah]) | `When/IsCLMNP.xml`, `CountLossAllocation_act` |
| Work class komite | `ASM-FW-GCNMFW-Work-KomiteTreatyNonProp`, flow `KomiteTreaty_Flow`, awalan `KMTNP-` [DEV] | `CreateChildKomiteCNP_Act`, `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` |
| Workbasket Input Acceptation | `TreatyinPNCTeknik` (Claim Prop memakai `ReasKlaimTeknik`, keputusan 07-10-2026) | `Flow_TreatyIn` Assignment1 |
| Ukuran korpus | 151 berkas hanya ada di Non Prop; 129 bernama sama dengan Prop, semuanya beda isi. Pohon xlsx: 281 rule unik. | pembanding berkas, `Struktur_Flow_TreatyIn.xlsx` |
| Kasus di Pega DEV | CLMNP: 8 (1 Open, 1 Resolved-Completed, 6 New). KMTNP: 23 (16 Resolved-Completed). | [DEV] |

## 2. Alur dan siklus hidup

```
Start → Assignment2 "Outstanding Claim"  (worklist pembuat; FA OutstandingClaim)
      → Assignment1 "Input Acceptation"  (workbasket TreatyinPNCTeknik; FA InputAcceptation; tiket AcceptanceClaim)
      → Decision IsBackStage (.pyNote = "Back" → kembali ke Assignment2) / Else → End (Resolved-Completed)
```

- **Pembuatan kasus.** Harness `New` / `NewSample` tidak diekspor, sama seperti di Prop (OQ-CP-13).
- **Outstanding Claim**
  - Pre-activity `InputOutStandingClmTNP_PreAct` memuat polis Non Prop dari `TREATYINPRODUCTION` dan daftar treaty
    (`PROPORTIONALARRG` 'TREATY LIMIT', hanya yang memuat TRT/ORS). Loss allocation diisi "OR".
  - Post-activity `InputOutStandingCTNP_PostAct` mengisi `CNPStatusCase = "INPUT ACCEPTATION CLAIM"`. Bila nomor polis
    kosong, muncul pesan "Policy no can't be empty!".
- **Input Acceptation**
  - Pre-activity `InputAkseptasi_PreAct`:
    - Akseptasi lama (≤ 14-12-2020 tanpa `ListClaimAcceptation`) ditandai `CommentLOD = 1`, sehingga detailnya
      disembunyikan.
    - `IsAcceptation = 1` bila ada akseptasi berstatus 1.
    - Mesin XoL dijalankan bila alokasinya masih kosong.
- **Jalan keluar dari Input Acceptation**
  - **Close Claim** (`CloseClaimMD` → `CloseClaimTNonProp`):
    - Ditolak bila masih ada akseptasi di komite (status 0).
    - Bila lolos: baris OS `STS_REJECT = 4`, JSON_KLAIM, REST `insertClaimFinalOrClosed_NP` (**tanpa** penjaga
      produksi), lalu `ASMForceCaseClose` Resolved-Completed.
  - **CWP** (`CloseClaimNP`):
    - Ditolak bila akseptasi pertama sudah berstatus 1 ("Claim tidak bisa di Close, Sudah ada akseptasi").
    - Bila lolos: kasus komite satu tingkat (bab 7.6). Bila disetujui, klaim induk ditutup Resolved-Completed (close)
      atau Resolved-Rejected (reject).
  - **Submit** standar flow action menuju End [inferensi]. Section Input Acceptation sendiri tidak punya tombol Submit.
  - **Back** ada di container `NEVER`, sehingga tidak terjangkau. DT `BackFromRegister` hanya ada di korpus Claim Fac In.
- **Nilai `CNPStatusCase`:** INPUT ACCEPTATION CLAIM → COMITEE ACCEPTANCE (DEPT. HEAD) → CLAIM ACCEPTED / CLAIM REJECTED.

## 3. Layar

Kedua layar utama adalah section tunggal yang besar. Prop memecahnya menjadi `_Intrs` / `_Est` / `_Adjs`.

### 3.1 Outstanding Claim (`Section/OutstandingClaim(1).xml`)

**Kepala.** Claim No (tampil bila terisi) dan blok "Claim Treaty". Semua isian kepala hanya-baca:

- Master Treaty ID, Treaty Name, Class of Business, Ceding Name / SOB Name (tampil bila `IDMasterTONP` kosong),
  Bordereaux (+ catatan), Treaty Year, Start/End Date, Accounting Mode, Territorial Scope, RNM Share %.

**Tombol kepala:**

- Choose Master In: popup `ChooseMasterTNonProp`, saringan `TREATYINDETAIL` ∪ `TREATYINDETAILEDM` dengan
  `PROPORTIONTYPE = 'NonProportional'`.
- View Master: harness `InputTreatyInOffer`, **tidak diekspor**.
- Choose Master Out dan CWP: `NEVER`.

**Tab Claim Information:**

- Policy No: autocomplete dari polis master terpilih, wajib. Memilih polis juga mengisi Layer / LayerType / LayerPart.
- Tombol View List Policy, View (polis), dan View Payment Status (REST `getPremiumPaidOnTreatyIn`).
- Isian lain: Policy No Ceding, Insured Name, Reinsurance Slip, Claim No Ceding, PLA No Ceding, DLA No Ceding,
  PLA No SOB, Policy Start/End Ceding (cek periode treaty), DOL / Report / Received date, Reporter, Catastrophe,
  Cause of Loss (popup), Report Type, Report Address, Adjuster, Consultant, Circumstances (wajib), Occupation,
  Location, Province, Postal Code.
- Isian yang menyangkut polis menjadi hanya-baca setelah `IsOutstanding = 1`.

**Tab Interests.** Grid Insured Interests 100% dengan Add/Delete; detailnya popup `InputDtlInterest` beserta aturan TPL
(format 1/2, % dari Claim/TSI, min/max). Di bawahnya total per mata uang dan total IDR.

**Tab Estimation:**

- Share Ceding.
- Deductible: Deductible Type → Form Type 1/2, mata uang, nilai, % of Claim/TSI, Min/Max. Nonaktif setelah CFS dicetak.
- Grid Claim Amount per mata uang: Currency, ROE, Claim Amount, TPL, Adjuster Fee, Salvage, Fee, Proportion %, Claim
  Amount IDR, Claim Amount Cedant.
- Grid Loss Allocation (`CNPSpreadLoss`): Currency, Treaty Name, Share %, Claim Amount, fee-fee, centang **To XOL**.
- Tombol **Edit XOL Allocation** (popup berkata sandi, bab 4.10).
- Grid XOL Allocation (`SpreadingRisk`, baris UR + layer): Claim Amount bisa diedit, sisanya hasil hitung.
- Grid Summary XOL Allocation.
- Grid Reinstatement: `NEVER`.

**Tab Spreading.** Grid Spreading List dan Spreading Break QS. Container-nya ALWAYS, jadi menurut aturan R2 tetap
tampil walau isinya NEVER. Tombol hitungnya memanggil `CountSpreadingCNP_Act`, yang **tidak diekspor**.

**Tombol bawah:**

| Tombol | Aksi | Syarat |
| --- | --- | --- |
| Save | `SaveDataToJClaim_Act` → JSON_KLAIM | — |
| Save to issue RNM | `SaveDataToOSAksep_Act`: validasi → nomor klaim → baris OS selisih (STS 0) → `IsOutstanding = 1`, baris dikunci (`CNPFlagOuts = 1`) | — |
| Print CFS | `GenerateCFS_act` → PDF CFS → `IsCFS = 1` | nonaktif sampai `IsOutstanding = 1`; nonaktif lagi setelah dicetak |
| Print PLA | FA `GeneratePLACNP`: nomor PLA + PDF per reasuradur retro | — |
| Submit | finishAssignment → Input Acceptation | **nonaktif** sampai `IsCFS = 1` dan `IsOutstanding = 1` |

Grid Claim History (`SuggestList`, kronologi) ada di bawah.

**Beda dengan Prop.** Non Prop tidak punya tombol "Send to Acceptation". Di Prop, tombol itu langsung memindahkan
berkas ke Teknik (keputusan WO 08-10-2026). Di Non Prop, berkas pindah lewat Submit, setelah CFS dicetak.

### 3.2 Input Acceptation (`Section/InputAcceptation.xml`)

**Tombol kepala:**

- View Master (tidak diekspor).
- Close Claim, Close Without Payment (CWP).
- Choose Master: `IN` sebelum ada akseptasi; `INEDM` setelahnya.
- View Payment Attachment (tampil bila `IsAcceptation = 1`).
- Claim History Master ID: riwayat pemakaian layer per master dari `CLAIMXOL2`.

**Tab Claim Information.** Sama dengan Outstanding, tetapi hanya-baca setelah ada akseptasi. Tidak ada Province /
Postal Code. Ada tambahan isian Supporting Document.

**Tab Interests.** Hanya-baca.

**Tab Estimation:**

- Isinya grid-grid Outstanding, ditambah:
  - centang **Waiting For Actual Premium** dengan tombol View (popup selisih);
  - tombol View Old Allocation;
  - tombol **Save To OS** (`SaveToOS`: baris OS selisih + konversi + CFS ulang).
- Grid status komite tampil bila kasus sedang CWP (`IsCloseFile` / `IsReject`).

**Tab Acceptation:**

- Grid Acceptation List (`AdjustmentList`): Type (bisa diedit), Accepted No/Date, Acceptance Status. Membuka baris
  menampilkan detail akseptasi (bab 3.3).
- Add → `AddAkseptasiCNP_Act`, dengan validasi: BusinessName, Adjuster, dan Consultant wajib.
- Delete → `DeleteAkseptasi_Act`.
- Di bawahnya grid Spreading Claim Out (In + QS) dan Claim History.

### 3.3 Detail akseptasi (`AdjustmentDetailNP` + `Subjectivity`)

**Kepala dan grid:**

- **Payment Type**, wajib. Nilai `2` memunculkan nomor interim `CNPIndexInterim` = interim sebelumnya + 1. Kunci saat
  akseptasi sudah di komite.
- DLA No.
- Grid Claim Acceptation (per mata uang).
- Grid Loss Allocation.
- Grid **XOL Allocation**:
  - Kolom: Total Claim RNM, Reinstatement Premium RNM.
  - Membuka baris menampilkan popup rumus reinstatement (bab 4.4).
- Grid **Previously Calculated** (`AlokasiXOLPaid`, hanya-baca).
- Grid Spreading In / Spreading Out. Kolom yang ditambah Non Prop: Fee, Total Claim, Reinstatement Premium,
  **RNM Net Claim**.

**Pembayaran:**

- Payable To: 1 Ceding / 2 SOB / 3 Others [inferensi]. Bila 3, Specify dari daftar klien.
- Rekening pertama, dan **rekening kedua** bila multi mata uang (`FlagCurrency = 1`).
- Transfer Direct to Kasir. Nonaktif bila sudah di komite atau `FlagErrorKasir = 1` (ada Net Claim negatif).
- Status Kasir.

**Tombol:**

| Tombol | Aksi | Syarat |
| --- | --- | --- |
| Save | simpan | — |
| Send to Committe | proteksi → popup `KomiteCNP` → "Send Claim to Committee" → `CreateChildKomiteCNP_Act` | tampil bila `TotalKomite` terisi |
| Acceptation | `HitServiceToKasir_Act` (tanpa simpan akseptasi / cetak DLA seperti di Prop) | — |
| Generate DLA | FA `GenerateDLACNP`, post-activity **tidak diekspor** | — |
| Generate Claim Analysis | `GenerateCACNP_Act` → PDF Claim Analysis | — (Prop: tampil bila status 1) |
| Save Previously Paid | `SaveCNPLayerList_Act` → baris OS `STS_REJECT = 5` | — |

Seluruh body disembunyikan bila `CommentLOD = 1` (akseptasi lama).

### 3.4 Popup

- **KomiteCLMNP**
  - Isian: tanggal, initial, Circumstances, Remarks.
  - "Send Claim to Committee" nonaktif bila Payable kosong, ada error proteksi, atau Occupation kosong.
- **CloseClaimMD:** Remarks wajib, lalu Yes menutup kasus.
- **CloseClaimNP (CWP):** tanggal, user, Remarks wajib, lalu Yes membuat komite.
- **PreviewPLA:** Remarks. **GenerateDLACNP:** Remarks DLA.
- **ChooseMasterTNonProp:** cari Treaty ID, grid master, Choose.
- **ViewListPolicyCNP:** daftar polis.
- **ReinstatementPremiumDetails:** rumus per `FlagProrate`, hanya-baca. Tidak ada activity yang menulis `FlagProrate`,
  jadi praktis hanya varian 0/1 yang tampil [inferensi].
- **EditXOLAlokasi:** kata sandi + Note.
- **ViewHistoryMasterID_NP**, **ViewOldAllocation**, **Hitung_Test** (selisih aktual).
- **ViewAttachmentNP:** lampiran invoice. **DetailPaymentCNP:** status bayar premi.
- **View polis:** `ViewDetailDeptHeadTreatyIn_UW` + `DetailPolisCNP`, tampilan polis NB Treaty In NonProportional
  (limit / MDP / share / installment).

### 3.5 Yang tersembunyi permanen (NEVER)

- Grid Reinstatement dan grid Limit Layer, beserta popup turunannya `ViewClaimLayerDetail` / `CoBList`.
- Choose Master Out, CWP di Outstanding, View Outstanding Claim, Back, Save (Objsave).
- Report Description dan Insured Interest di Outstanding.

## 4. Hitungan

Singkatan: CD = `pyWorkPage.ClaimData`, RNM% = `TreatyInMaster.RNMShare`, p(x) = x/100.

### 4.1 Klaim 100% per mata uang (`CountClaimTNP_Act`)

**Per baris Claim Amount:**

- Deductible hanya untuk baris pertama:
  - Form Type 1: nilai tetap atau % × (Claim | TSI).
  - Form Type 2: max(nilai, %) bila Min/Max = 1.
  - Konversi mata uang memakai ROE baris.
- **Nilai bersih** = Value − Deductible + TPL. TPL diambil max(TPL baris, TPL interest) bila `IsTPL`.
- **NetGross** = nilai bersih + Fee + Adjuster Fee − Salvage.
- **USD** (IDR) = NetGross × ROE.
- **Claim Amount Cedant** = NetGross × p(Share Ceding).
- **Proportion %** = IDR mata uang itu / Σ IDR × 100. Ini menjadi bobot prorata limit per mata uang.

**Ke Loss Allocation.** Untuk baris berkunci `CNPFlagOuts ≠ 1`, Claim Amount = (nilai bersih × Share Ceding) × Share %
baris. Fee-fee dikali share yang sama.

**Penutup:** memanggil mesin XoL (4.2).

### 4.2 Mesin XoL (`CountLossAllocation_act`, parameter `Calculation`)

1. **Mode `pct` / `amount`.** Menghitung satu baris Loss Allocation (Claim Amount ↔ Share %). Lompatan keluar sebelum
   bagian XoL berprakondisi **nonaktif**, sehingga XoL selalu dihitung ulang [inferensi semantik flag, bab 10].
2. **Hanya baris Loss Allocation bercentang *To XOL*** yang dijumlahkan per mata uang:
   ClaimValue = Claim + Adjuster Fee + Fee − Salvage.
3. **Layer master.** Layer diambil dari `TreatyInMaster.Limits` yang treaty group-nya cocok. Mata uang layer dan kurs
   konversi diambil dari `TreatyInMaster.CurrencyList`.
4. **Per mata uang klaim, per layer:**
   - `LayerLimit` = Limit layer dikonversi ke mata uang klaim × Proportion %.
   - `URLimit` = Deductible layer (retensi) dikonversi × Proportion %.
   - Rumus konversinya **tidak simetris** antara IDR dan valas, memakai `Limit2` / `Deductible2` (OQ-CNP-18).
5. **Layer pertama** menambah baris **UR** (retensi):
   - Claim Estimation = min(ClaimValue, URLimit).
   - Share 0, Claim RNM 0.
6. **Tiap layer** menambah baris:
   - **Nama layer** "XL 1ST/2ND/3RD/nTH LAYER" (atau SUB LAYER), dipetakan ke `REINSURANCETYPE` (type 4, active). Bila
     nama tidak memuat ST/ND/RD/TH, muncul "Error TreatyName, Hubungi IT".
   - **Claim Estimation** = min(ClaimValue − sudah teralokasi, LayerLimit). Berhenti bila sudah habis.
   - **Claim Spreaded (porsi RNM)** = Claim Estimation × RNM%. Bila Claim Estimation = LayerLimit, dikurangi fee bersih.
   - **Fee-fee** = total fee × RNM%. Nol bila layer penuh.
   - **MDP** = `MDPList` layer menurut mata uang, dikonversi. Bila `IsCombineMDP`, dijumlah lintas layer sejenis. Nol
     bila `NoRIPCalculation`.
   - Juga disimpan: CNPLimit, CNPMDP, CNPPctReinstate (% reinstatement layer), ClaimAmountIDR.
7. **Rekap Summary XOL** per mata uang, lalu Spreading:
   - Share XoL polis (`Share(1).SpreadingTypeXOL`, `SpreadingTotalPctXOL`).
   - Break QS per `SpreadingListXOL`. Bila tipe kosong: "QS (OR)" → 10028, "QS (R/I)" → 10004.
8. **Penutup:** memanggil `CountReinstatement_Act`.

### 4.3 Penyesuaian alokasi (`AdjClaimAmount_Act` di klaim, `AdjClaimCNP_Act` di akseptasi)

- Mengubah Claim Amount satu layer akan menghitung ulang AdjClaimValue, Claim Spreaded, dan Claim Amount Adjust.
- **Layer terakhir menyerap selisih** terhadap total Loss Allocation.
- Di akseptasi ada tambahan:
  - **Total Claim RNM** = Claim Spreaded + fee bersih × RNM%.
  - **Net per layer** = alokasi sekarang − `AlokasiXOLPaid` (yang sudah dibayar).
  - **Spreading In / QS:** Claim = Net × Share%; Premium Spreaded = Reinstatement RNM × Share%;
    **Net Claim = Total Claim − Premium Spreaded**.
  - Akhirnya memanggil `CountSpreadingXOL` (rekap ke level kasus) dan `ProtectNilaiClaim`. Bila ada Net Claim < 0,
    Direct to Kasir dimatikan (`FlagErrorKasir = 1`).

### 4.4 Reinstatement premium

Ada dua rumus di XML yang tidak sama persis (OQ-CNP-12):

| Tempat | Rumus |
| --- | --- |
| Akseptasi (`AdjClaimCNP_Act` s.11), dipakai komite dan OS | `CNPReinstatement` = (UR ? 0 : Total Claim / CNPLimit × CNPMDP × p(CNPPctReinstate)); `CNPReinstatementRNM` = × p(RNM Share baris) |
| Kasus (`CountReinstatement_Act`), mengisi grid Reinstatement yang `NEVER` | ((Claim Estimation + Adjuster Fee RNM) − Salvage × 100/RNM%) / CNPLimit × CNPMDP × p(Pct) |
| Popup `ReinstatementPremiumDetails` | 0/1: Claim Layer / Limit × MDP × Pct. 2: (Sisa Limit − Limit) / Limit × … 3: Old Gross / Limit × … |

### 4.5 Payment Type, interim, cancellation

- **Payment Type 2 (interim):** nomor interim naik per akseptasi interim yang sudah disetujui (`SetInterimXOL_Act`).
- **Payment Type 7 (cancellation)** [inferensi]:
  - Proteksi Send to Committe: bila semua baris XoL adalah UR (tanpa recovery) dan Payment Type ≠ 7, muncul "Payment
    type must be Cancellation".
  - Saat komite dibuat dengan Payment Type 7 dan satu baris UR, disusun baris pembalik bernilai negatif dari akseptasi
    sebelumnya.
- **Payment Type 1** saat disetujui komite menghasilkan baris OS `STS_REJECT = 4` (final) dan `IsFInalAccXOL = 4`.
  Tipe lain menghasilkan STS 1.
- Label 1–7 lainnya ada di prompt values `associated` yang tidak diekspor (OQ-CNP-14).

### 4.6 Proteksi sebelum komite (`ProteksiSendKomiteCNP_Act`)

- Pesan yang bisa muncul:
  - Data Bank Account Can't NULL / Occupation cannot be empty / Please choose Payment Type first / Please Input Policy
    No / Currency cant be NULL / ID of Bank Can't be NULL / Email Profile Can't NULL / Please choose payable first.
  - **"Adjustment value should not be more than estimation value"**: total adjustment per mata uang > Summary XoL.
- Proteksi "Error No Account" tidak pernah aktif karena kondisinya bertentangan.
- Di `CreateChildKomiteCNP_Act`, "Nilai Gross Value tidak sesuai dengan Spreading In" baru dievaluasi **setelah**
  kasus komite dibuat.

### 4.7 Premi aktual

Centang Waiting For Actual Premium (`SetActualPremium_ACT`) menimpa layer terakhir dengan Claim Spreaded = **50.000**
(tertulis mati), lalu total digross-up ÷ RNM%. Popup `Hitung_Test` menampilkan selisih terhadap OS (`GetSelisihActual_Act`).
Kolom hasil SQL `GetDataOS` tidak cocok dengan yang dibaca (CARI1/CARI2), jadi nilainya selalu kosong [inferensi].

### 4.8 Penomoran

| Nomor | Format | Sumber |
| --- | --- | --- |
| Nomor klaim | `KODE + "K" + BusinessOldId + "." + HASIL1 + ".TX" + HASIL2` | `KODE_PRODUKSI` NONLIFE + `PROC_GENERATE_SEQUENCE_NUMBER` |
| Nomor akseptasi (komite) | `KODE + "A" + BusinessOldId + "." + HASIL1 + ".TX" + HASIL2` (Prop: `.TP`) | prosedur yang sama |
| Nomor PLA | `'RNM-M' + BusinessOldId + '.' + MM + '.' + yyyy + '.TX' + lpad(PLATNP_SEQ, 5)`; cetak ulang = `/n+1` | `GenerateNoPLATNP` |

Parameter tanggal `TO_DATE(CARI3)` prosedur nomor tidak pernah diisi di korpus (OQ-CNP-27).

### 4.9 Dokumen

- **CFS** (`GenerateCFS_act`): PDF laporan XoL.
  - Langkah penghapus list REMARK, sehingga agregasi menjumlah ke list yang sudah terisi.
  - Field Loss Allocation tertukar (Fee ↔ Claim Amount).
  - Field komisi dipakai sebagai slot nilai. Lihat OQ-CNP-05.
- **PLA** (`GeneratePlaCNP_Act` + `GeneratePlaCNP2_Act`):
  - Recovery retro dihitung dari share "R/I" × total klaim IDR.
  - Layer treaty out dibaca dari `TREATY_OUT` (tidak ada di DEV).
  - Satu PDF per reasuradur.
- **Claim Analysis** (`GenerateCACNP_Act`): alokasi berjalan (+) dan lama (−) per layer, net setelah reinstatement.
- **DLA:** activity tidak diekspor.
- **Persetujuan klaim** (komite, `GenerateAccCNP_act`): nama PDF hanya diisi bila When `IsFire` (tidak diekspor).
- Semua stream HTML dokumen **tidak diekspor**. Sama seperti Prop, Claim Prop memakai go-pdf/fpdf.

### 4.10 Edit XOL Allocation

Kata sandi dibandingkan dengan teks tetap di XML. Bila benar, semua baris XoL ditandai `IsEditClaim = 1`, Note masuk
kronologi, lalu disimpan. Bila salah, muncul "Invalid Password".

### 4.11 Hitungan polis (bukan klaim)

`CountResult*`, `CountNetPremi_act`, `SetPPNPPH` (brokerage 2,5%, PPh 2%, PPN 2,2%), `FillPaymentInstallment`,
`SetDueTo_act`, `DetailCalculation`, `TotalEgnpi`, dan `TreatyInNonSetTotal` berkelas polis/master Treaty In. Semuanya
dipanggil dari tampilan polis (`ViewDetailDeptHeadTreatyIn_UW`), bukan dari hitungan klaim.

## 5. Penyimpanan Pega dan fakta DEV

| Tabel | Penulis | Kapan | DEV, kasus `CLMNP-` |
| --- | --- | --- | --- |
| `JSON_KLAIM` | `PEGA_JSON_KLAIM_PNC` (`InsertClaimPNC`) | Save, Save to issue RNM, close, komite | 6 baris |
| `OS_AKSEPTASI_KLAIM` | `PEGA_JSON_OS_AKSEP_KLAIMTNP` (klaim), `PEGA_JSON_OS_AKSEP_KLAIM` (komite) | lihat kode di bawah | 1.684 kasus / 4.178 baris, 2019-07-25 s.d. 2026-07-06 |
| `OS_AKSEPTASI_SUBJECTIVITY` | `PEGA_JSON_OS_AKSEP_SUBJECTIVITY` (UPDATE bila ada, INSERT bila belum) | persetujuan komite tingkat akhir | 549 baris |
| `CLAIMXOL2` | `XOL2_AKSEP_KLAIM`: **INSERT saja**, per baris XoL non-UR | persetujuan komite tingkat akhir | 685 kasus / 1.066 baris; satu kasus 41 kali |
| `HISTORYAKSEPTASIPEGA` | `InsertHistoryAkseptasiPega_Sql` | tiap keputusan komite jalur akseptasi | 652 ACCEPT / 42 REJECT |
| `DIRECTTOKASIR_LOG` | `InsertLOGDirectKasir_SQL` | kirim Kasir | 397 baris |
| `CLAIMREJECTED` | **tidak ada penulis** di XML Non Prop | — | 102 baris Resolved-Rejected |
| `DOCUMENT_CLAIM` + `T_STORAGE_IMAGE` | `InsertDocument_Act` / `InsertGoogleStorage_Act` | PDF CFS / PLA / CA / persetujuan | — |

**Kode `STS_REJECT` di `OS_AKSEPTASI_KLAIM`.** Parameter `CARI10` = kolom `STS_REJECT` (terbukti dari sumber prosedur
di DEV).

| Kode | Arti | DEV: kasus / baris |
| --- | --- | --- |
| 0 | Outstanding (Save to issue RNM, Save To OS); nilai = **selisih** terhadap Σ baris STS 0 sebelumnya | 1.678 / 2.274 |
| 1 | Akseptasi disetujui komite (Payment Type ≠ 1) | 659 / 904 |
| 2 | Reject lewat CWP: baris **negatif** pembalik outstanding per mata uang × layer | 52 / 71 |
| 4 | Final / close (Close Claim, CWP close, akseptasi Payment Type 1) | 837 / 929 |
| 5 | Save Previously Paid | **0 baris** |

**Bentuk baris.** Kunci baris adalah `CASEID` (pzInsKey Pega `ASM-FW-GCNMFW-WORK CLMNP-n`) + `data_json.TypeLoss`
(layer) + `data_json.Currency`. Prosedur klaim selalu INSERT; cabang UPDATE lamanya dikomentari. `DATA_JSON` berisi
`CNPLayerList` / kolom XoL. `OS_AKSEPTASI_KLAIM` sudah punya kolom datar XoL: `XOL`, `XOLID`, `CNPLIMIT`, `CNPMINDEP`,
`CNPPCTREINSTATE`, `CNPREINSTATEMENT(RNM)`, `GROSSADJUSTMENT`, `KURSIDR`, `TOTALXOL*`, `CNPOTHERSFEE(RNM)`,
`ADJUSTERFEE*`, `SALVAGE*`, `PERSENRNM`, `KURSVALUE`. Kolom-kolom ini fisik, bukan virtual, tetapi **kosong di semua
4.178 baris `CLMNP-`** [DEV]. Isi XoL hanya ada di `DATA_JSON`. Yang terisi pada baris Non Prop hanya kolom yang ditulis
prosedur (`CASEID`, `NOCLAIM`, `MASTERID`, `DATA_JSON`, `TANGGAL`, `NOPOLIS`, `STS_REJECT`, `STS_KONVERSI`) dan
`TGL_PROD` [inferensi dari nama trigger `TRG_TLG_PROD_OS_AKSEPTASI`, BEFORE INSERT, aktif di DEV].

**Objek DEV yang dirujuk:**

- **Ada:** semua prosedur di atas, `PLATNP_SEQ`, `M_TREATY_IN(_EDM)`, `M_TREATY_OUT(_DETAIL)`, `M_LINK_SERVICE`,
  `KODE_PRODUKSI`.
- **Tidak ada:** `TREATY_OUT`. Yang ada `TREATY_OUT2`.
- **INVALID:** view `CLAIMXOL`.
- **Kategori `M_LINK_SERVICE`:** Klaim hanya punya `insertClaimAccept`, `insertClaimDLA`, `insertClaimLife`. URL
  `insertClaimFinalOrClosed` dan `insertClaimReject` tertulis mati di connector.

**Roster `EMAILKOMITE` STS_KLAIM NONPROP.** Empat baris aktif. Kolom `OPERATOR_ID`-nya masih akun orang.

| ID | DEGREE | JABATAN | LIMIT_BOTTOM | LIMIT_TOP |
| --- | --- | --- | --- | --- |
| 9 | 1 | Claim Dept. Head | −9.999.999.999.999 | 25.000.000 |
| 14 | 2 | Technic Div. Head | 25.000.001 | 495.000.000 |
| 11 | 3 | Operational Director | 495.000.001 | 580.000.000 |
| 12 | 4 | Technical Director | 580.000.001 | 999.999.999.999.999 |

**Master treaty.** `TREATYINDETAIL` NonProportional 18.401 baris; `TREATYINDETAILEDM` 11.482 baris. `T_WORK_CLAIM` belum
punya baris Non Prop.

## 6. Integrasi

| Connector | Dipanggil | Penjaga produksi | Catatan |
| --- | --- | --- | --- |
| `InsertClaimOutstanding_NP` | Save to issue RNM | ya, kecuali akun uji | URL tertulis mati; stsReject = 0 |
| `KonversiKlaimNonLife` | Save To OS; komite final (STS 1) | ya | LinkService Klaim/insertClaimAccept; Non Prop **tanpa** log `MONITORING_KLAIM_LOG` |
| `insertClaimFinalOrClosed_NP` | Close Claim, CWP close | **tidak** di Close Claim | URL tertulis mati |
| `insertClaimReject_NP` (Komite) | CWP reject | ya | URL tertulis mati |
| `SendAcceptationToKasir` | tombol Acceptation; komite final | ya | Nett = Total Claim − Premium Spreaded per Spreading In; wajib sudah terkonversi (`REINSURANCE.TRLOSS_DETAIL_T`) |
| `getPremiumPaidOnTreatyIn` | View Payment Status | — | sama dengan Prop |
| `getPayAttachment` | View Payment Attachment | — | sama dengan Prop |
| `ServiceGoogle` | unggah / URL dokumen | — | sama dengan Prop |

Pola Claim Prop yang akan diikuti: efek luar dicatat ke outbox hanya bila produksi, dan pelaksananya berhenti terang
(OQ-CP-03).

## 7. Komite Claim Non Prop

### 7.1 Pembentukan dari akseptasi (`CreateChildKomiteCNP_Act`)

- **Batas tertulis mati.** Batas nilai 30.000.000 dan batas persen 30. `Flagkomite = 1` bila RNM Share ≤ 30 **dan**
  `ValueAdjustment` akseptasi terakhir ≤ 30.000.000. Batas Div Head 50.000.000 di-REMARK.
- **Roster** dibaca dari RD `FilterEmailKomiteWithLimit` dengan STS_KLAIM = NONPROP, aktif, `LIMIT_BOTTOM ≤ param`:
  - `Flagkomite = 1` atau subjectivity: param = 0, sehingga **hanya tingkat 1** (Dept Head).
  - Selain itu param tidak diisi, sehingga saringan batas diabaikan dan **semua tingkat aktif** terambil (4)
    [inferensi]. Ini didukung riwayat DEV: kasus komite lama melewati 1, 2, 3, atau 4 tingkat.
  - Akibatnya **kolom LIMIT roster NONPROP tidak dipakai** XML. Claim Prop memakainya (`LIMIT_BOTTOM = TotalAdjustment`).
- **Jabatan, inisial, dan penggantian akun** tertulis mati per nama orang. Akun uji punya jalur komite dirinya sendiri.
- **Penyalinan ke kasus komite:** data akseptasi, `CNPLayerList`, Spreading In / QS, dan PDF Claim Analysis.
  `KomiteNo` dari pzInsKey anak.
- **Status:** `CNPStatusCase = "COMITEE ACCEPTANCE (DEPT. HEAD)"`, `AcceptanceStatus = 0`, lalu email ke tingkat 1.

### 7.2 Flow dan router

- `KomiteTreaty_Flow`: assignment `KomiteRouter` → FA `ViewTransferDtl` (pre `SetKomiteList_Act`, post
  `KomitePostAdjustment`) → loop selama `AcceptStatus = 1` dan `KomiteCount ≤ KomiteLoop`. Shape End tanpa status
  (Prop: Resolved-Completed).
- **Router.** Hanya step 6 yang aktif: bila `TransferType == 2`, kasus ke **akun** KomiteID pertama yang belum
  menyetujui. Routing ke workbasket `komitepnc..4` ter-REMARK, sama seperti Prop.
- Kasus CWP tidak mengisi TransferType, sehingga tujuan routing-nya tidak jelas (OQ-CNP-25).

### 7.3 Layar keputusan (`ShowTransfer`)

- **Judul:** "CLAIM COMMITTEE -" + CLOSE / REJECT.
- **Blok Claim Treaty** (data klaim) dan grid Claim Acceptation.
- **Grid alokasi:** Loss Allocation (+ To XOL), XOL Allocation (+ Reinstatement Premium, popup rumus), Previously
  Calculated, Spreading In / Out.
- **Pembayaran:** Payable / bank, dan rekening kedua.
- **Grid status komite:** jabatan / inisial, keputusan, tanggal, komentar.
- **Isian:** Remarks, keputusan ("Are you sure to accept this document?"), Subjectivity + note, Propose Close,
  Propose Reserved, Comment.

### 7.4 Keputusan (`KomitePostAdjustment`)

- **Cek pemutus.** Bila operator bukan KomiteID tingkat berjalan, hanya muncul pesan "Invalid User Acceptance!";
  proses **tetap jalan**. Prop keluar dari activity.
- **Tolak:**
  - Semua tingkat sisa ditandai tolak, `AcceptanceStatus = 2`, status klaim "CLAIM REJECTED".
  - Tanpa email.
- **Setuju di tingkat akhir (bukan subjectivity):**
  1. Nomor akseptasi `…A….TX…`.
  2. Status "CLAIM ACCEPTED".
  3. `InsertOSKlaimCNP`:
     - PDF persetujuan.
     - OS lewat `PEGA_JSON_OS_AKSEP_KLAIM` (STS 4 bila Payment Type 1, selain itu 1; JSON memuat `CNPLayerList`).
     - `CLAIMXOL2` per layer.
     - JSON_KLAIM.
  4. Konversi (STS 1).
- **Tingkat akhir, termasuk non-subjectivity:** `OS_AKSEPTASI_SUBJECTIVITY`.
- **Produksi saja:** email (ke tingkat berikut, atau ke pembuat bila final) dan Kasir (bila Direct to Kasir).
- **Setiap keputusan jalur akseptasi:** `HISTORYAKSEPTASIPEGA` (Workbasket "KLAIM").
- **Kasus komite ditutup** bila semua tingkat selesai.

### 7.5 Beda dengan Komite Claim Prop

| Unsur | Non Prop | Prop |
| --- | --- | --- |
| Pembagi jalur | flag `IsCloseFile` / `IsReject` | `TransferType` 2 / 3 / 4 |
| Cek pemutus | pesan saja | keluar |
| Nomor akseptasi | `.TX` | `.TP` |
| Tulis OS | + `CLAIMXOL2`, + OS subjectivity, + PDF persetujuan | `SaveAcceptationTreaty_TKMT` |
| Konversi | tanpa transisi gagal, tanpa `MONITORING_KLAIM_LOG` | dengan log |
| Email saat tolak | tidak | ya |
| Reject klaim | baris OS negatif (STS 2) + REST reject; tanpa `CLAIMREJECTED` | `CLAIMREJECTED`, konversi "2" |
| Tangga | 1 tingkat (≤ 30 jt & ≤ 30%) atau semua tingkat | batas roster per nilai |

### 7.6 Close Without Payment (`CreateChildKomiteCloseNP_Act` → `KomitePostAdjustmentCWP`)

- **Pembentukan:**
  - Komite **satu tingkat**, tertulis mati ke satu akun Claim Dept. Head.
  - `IsReject = 1` bila `TypeComentAnalysis == "5"`. Satu-satunya pengisinya di-REMARK, sehingga praktis selalu
    `IsCloseFile = 1` [inferensi].
- **Disetujui:**
  - PDF "Close Claim / Reject Claim".
  - Close: baris OS STS 4 + REST final/closed.
  - Reject: baris OS negatif STS 2 per mata uang × layer + REST reject.
  - Klaim induk ditutup Resolved-Completed / Resolved-Rejected.
- **Ditolak:** status klaim "CLAIM REJECTED".

### 7.7 Rencana yang sudah diputuskan work owner (09-10-2026)

"Ikut pola Claim Prop":

- Tanpa menu Komite Claim Non Prop.
- Tabel komite di bawah inbox Claim Non Prop, dibuka di tempat.
- View more details = klaim hanya-baca.
- Roster `EMAILKOMITE` NONPROP diganti workbasket (`ReasClaimDeptHead` / `ReasClaimTechDivHead` / `ReasClaimOpsDir` /
  `ReasClaimTechDir`) dengan **UPDATE baris yang ada**. Penyetuju = anggota workbasket. Email ke semua anggota. Tanpa
  larangan rangkap.

Yang masih terbuka adalah aturan tangganya: XML tidak memakai batas roster (OQ-CNP-01).

## 8. Beda utama dengan Claim Prop

| Bidang | Claim Prop | Claim Non Prop |
| --- | --- | --- |
| Master | `PROPORTIONTYPE = 'Proportional'`, pilih RNM share | `'NonProportional'`, Choose Master In / INEDM, master `m_treaty_in(_edm)` JSON |
| Estimasi | Estimation List (TypeLoss, tanggal, %) | Claim Amount per mata uang + Loss Allocation + XOL Allocation |
| Alokasi | Loss Allocation per treaty (share %) | waterfall layer XoL (UR + layer), reinstatement, MDP |
| Akseptasi | Adjustment (Individual Risk, Propose Adjustment) | akseptasi XoL (Payment Type 1–7, interim, previously paid, Net Claim, rekening kedua) |
| Kirim ke Teknik | "Send to Acceptation" (langsung pindah, WO 08-10) | Submit setelah CFS + Save to issue RNM |
| Close | satu tombol, dua jawaban | Close Claim (langsung) + CWP (komite) |
| OS | per NOCLAIM; `PEGA_JSON_OS_AKSEP_KLAIM(TRT)` | per CASEID + layer + mata uang; INSERT selisih `…KLAIMTNP` |
| Dokumen | PLA, DLA, Claim Analysis, Acceptance Note | CFS, PLA retro per reasuradur, Claim Analysis, persetujuan; DLA tidak diekspor |
| Nomor | `GENERATE_NOCLMTREATYIN` dkk. | `PROC_GENERATE_SEQUENCE_NUMBER` (`.TX`), `PLATNP_SEQ` |
| Kasir | Nett = Adjustment Value | Nett = Total Claim − Premium Spreaded per Spreading In |

## 9. Rule dirujuk tetapi tidak diekspor

- **Activity:** `GenerateDLACNP` (post FA DLA), `CountSpreadingCNP_Act`, `SetTreatyNameSpreading_Act` (hanya di Prop),
  `Objsave`, `BackToRegister_act`.
- **Data transform:** `BackFromRegister` (hanya Fac In).
- **Harness:** `New`, `NewSample`, `Perform`, `InputTreatyInOffer`.
- **Flow action:** `CobDetails`.
- **RD:** `BrowseTreatyInDetail`.
- **Stream HTML:** `CFSClaimXOL`, `PLACNP_Html`, `ClaimAnalysisHTML`, `AccClaimKomite_HTML`, `CWPCaseClaim`,
  `EmailKlaim_HTML_KMT`, `EmailRejectCloseKlaim_KMT_HTML`.
- **Message rule:** `ErrorMsg1`, `StatusNotAccepted`, `FailedToOpenInstance`, `AddAttachmentFailed`.
- **When:** `IsFire`.
- **Section target refresh:** `CaseContent_Bottom`, `ItemListEstimation`, `ShowObjectAdj`, `Layers`, `LayersTONP`.
- **Prompt values `associated`:** Payment Type, Type akseptasi, Acceptance Status, Min/Max, TPL Format/Type,
  Subjectivity Note, Bordereaux, Accounting Mode, Layer Type, Komite Approval.
- **Tiket `AcceptanceClaim`:** tidak pernah di-raise di XML mana pun.

## 10. Aturan baca

- **Layar R1–R5** (`grilling-ronde-2.md` §13):
  - `pyCondition` berlaku hanya bila `pyVisible = OTHER`.
  - Container ber-opsi ALWAYS tetap tampil walau kondisinya NEVER.
  - Disable / wajib berlaku hanya bila flag `…New` = `true` / `truewhn`.
  - Hasil sensus Non Prop: 0 dari 36 section memakai `pyActionConditions`.
- **Activity:**
  - Label `//` = REMARK.
  - Prakondisi dengan `pyStepsPreCondition = false` tidak dievaluasi, jadi step selalu jalan (74 step di Non Prop;
    konvensi yang sama di port claimprop `turunan.go`, `kerugian.go`).
  - Kode aksi WHEN 1 lompat, 3 lewati, 4 keluar iterasi, 6 keluar activity [inferensi pola].
