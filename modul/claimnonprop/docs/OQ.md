# OQ — pertanyaan terbuka modul Claim Non Prop

> Dibuat 09-10-2026 dari pindai baca-saja ([`PINDAI.md`](PINDAI.md)). Tahap 1 dibangun 09-10-2026: OQ-CNP-01, 04, 05,
> 06 diputuskan work owner; OQ lain memakai kolom **Rekomendasi** sebagai bawaan (prompt tahap 1 §3). OQ baru dari tahap
> bangun ada di bagian D. Status setiap tombol: [`PARITAS.md`](PARITAS.md).
>
> Pengelompokan per pemilik jawaban mengikuti `modul/claimprop/docs/OQ.md`:
> - **work owner:** keputusan.
> - **pemilik ekspor Pega:** rule yang tidak ada di korpus.
> - **DBA:** objek basis data dan data warisan.
>
> Kolom **Rekomendasi** adalah usulan bawaan bila work owner setuju. Yang sudah diputuskan dan tidak dibahas lagi:
> komite ikut pola Claim Prop (tanpa menu, tabel komite di bawah inbox Claim Non Prop, roster `EMAILKOMITE` NONPROP →
> workbasket lewat UPDATE baris yang ada, penyetuju = anggota workbasket, email ke semua anggota, tanpa larangan rangkap).

## A. Work owner

| OQ | Pertanyaan | Rekomendasi | Bukti |
| --- | --- | --- | --- |
| **OQ-CNP-01** | ✅ **Diputuskan 09-10-2026: ikuti XML.** **Aturan tangga komite.** XML Non Prop **tidak memakai kolom batas roster** NONPROP (25 jt / 495 jt / 580 jt). Bila RNM Share ≤ 30% dan nilai akseptasi terakhir ≤ 30.000.000, komitenya hanya tingkat 1 (Dept Head). Selain itu semua tingkat aktif (4). Pilihan: **(a)** ikuti XML; **(b)** pakai batas roster per nilai seperti Claim Prop. | **(a)**. Batas 30 jt / 30% menjadi konstanta bernama di satu berkas model, seperti `komite.go` Claim Prop. | `CreateChildKomiteCNP_Act` s.10–12, s.26; riwayat DEV `HISTORYAKSEPTASIPEGA`: 1, 2, 3, dan 4 tingkat |
| **OQ-CNP-02** | Komite **Close Without Payment** di XML tertulis mati ke satu akun Claim Dept. Head. | Diganti tingkat 1 roster NONPROP, yaitu workbasket `ReasClaimDeptHead`. | `CreateChildKomiteCloseNP_Act` |
| **OQ-CNP-03** | **Hardcode per kasus / master / akun** (CLMNP-232/367/382/861/975, master 1000393 dan 1001130, akun uji pembuat kasus, peta nama orang ke jabatan komite). | Buang semuanya. Jabatan dan inisial diambil dari roster. | `AdjClaimCNP_Act`, `CountLossAllocation_act`, `SetKomiteList_Act`, `GetHistoryMasterID_NP` |
| **OQ-CNP-04** | ✅ **Diputuskan 09-10-2026: (a) kata sandi di konfigurasi.** **Edit XOL Allocation** dijaga kata sandi tetap di XML. Pilihan: **(a)** kata sandi pindah ke konfigurasi; **(b)** diganti hak akses (misalnya anggota `ReasKlaimTeknik`); **(c)** tombol dibuang. | **(a)**: paling dekat dengan XML, tanpa teks rahasia di kode. | `EditXOLAlokasi` |
| **OQ-CNP-05** | ✅ **Diputuskan 09-10-2026: perbaiki semua** (rekomendasi "persis XML" ditolak). **Kelainan perilaku XML** (lihat daftar di bawah tabel). | Setiap perbaikan dicatat di `PARITAS.md` sebagai penyimpangan sadar beserta perilaku Pega lamanya. | lihat butir |
| **OQ-CNP-06** | ✅ **Diputuskan 09-10-2026: pola Prop + tabel XoL.** **Penyimpanan.** Pola Claim Prop: tabel datar, tanpa JSON baru. Susunannya ada di bawah tabel. | Setuju susunan itu. | `modul/claimprop/backend/migrations/520–533` |
| **OQ-CNP-07** | **Data lama.** DEV punya 1.684 kasus `CLMNP-` di OS, tetapi hanya 6 `JSON_KLAIM` dan 8 kasus Pega. Pemuat seperti OQ-CP-18? | Tunda. Bangun alur kasus baru dulu; pemuat dibahas sesudahnya. | [DEV] |
| **OQ-CNP-08** | Workbasket Input Acceptation `TreatyinPNCTeknik` (tidak ada di DEV) diganti `ReasKlaimTeknik` seperti Claim Prop? | Ya. | `Flow_TreatyIn` Assignment1; claimprop `models/kasus.go` |
| **OQ-CNP-09** | Jalur **Back** (Input Acceptation → Outstanding) tidak terjangkau dari layar (container NEVER). | Tidak dibangun, sama dengan Claim Prop. | `InputAcceptation` S28 |
| **OQ-CNP-10** | **Akhir Input Acceptation.** Section tidak punya tombol Submit; kasus ditutup lewat Close Claim atau CWP (`ASMForceCaseClose`). Tombol Submit standar flow action dibuat? | Tanpa Submit. Kasus selesai lewat Close Claim / CWP. | `Section/InputAcceptation.xml`, `CloseClaimTNonProp` |
| **OQ-CNP-11** | Grid Spreading (Spreading List / Break QS / Spreading In) ada di container ALWAYS berisi NEVER. Ikut aturan R2 Claim Prop (tetap tampil)? | Ya, tampil hanya-baca. Tombol hitung share-nya nonaktif (OQ-CNP-21). | `grilling-ronde-2.md` §14a |
| **OQ-CNP-12** | **Reinstatement** punya dua rumus berbeda: (1) di akseptasi, dipakai komite / OS / Kasir; (2) di grid kasus yang NEVER. | Bangun rumus akseptasi saja. Grid kasus tidak dibangun. Popup rumus hanya varian `FlagProrate` 0/1 (tidak ada penulis varian 2/3). | `AdjClaimCNP_Act` s.11, `CountReinstatement_Act` |
| **OQ-CNP-13** | **View polis.** Tampilan polis Non Prop memuat hitungan polis Treaty In (komisi, PPN/PPh, installment). | Pakai tampilan polis NB Treaty In yang sudah ada (jendela Modal penuh seperti Claim Prop); hitungan polis tidak dibangun ulang di modul klaim. | `ViewDetailDeptHeadTreatyIn_UW`, `DetailPolisCNP` |
| **OQ-CNP-14** | **Prompt values `associated`** tidak diekspor: Payment Type 1–7, Type akseptasi, Acceptance Status, Min/Max, TPL Format/Type, Subjectivity Note, Bordereaux, Accounting Mode. | Minta screenshot Pega. Sementara tampilkan kode DB apa adanya. | sama dengan OQ-CP-07 |
| **OQ-CNP-15** | Komite menolak **satu akseptasi**, tetapi status seluruh klaim jadi "CLAIM REJECTED". | Ikuti XML. | `KomitePostAdjustment` step 20 |
| **OQ-CNP-16** | Cek pemutus komite Non Prop hanya memberi pesan; proses tetap jalan. | Diblokir, sesuai pola workbasket Komite Claim Prop (pemutus wajib anggota workbasket tingkat berjalan). | `KomitePostAdjustment` s.3–4 |
| **OQ-CNP-17** | **Waiting For Actual Premium** menimpa layer terakhir dengan Claim RNM 50.000 (tertulis mati). | Bangun persis XML, dengan nilai 50.000 sebagai konstanta bernama. | `SetActualPremium_ACT` |

**Butir OQ-CNP-05 (kelainan XML):**

1. `GenerateCFS_act`: penghapus list di-REMARK, sehingga agregasi CFS menjumlah ke list yang sudah berisi. Field Fee ↔
   Claim Amount tertukar di Loss Allocation.
2. `GetSelisihActual_Act`: kolom hasil `GetDataOS` tidak cocok dengan yang dibaca, sehingga popup selisih aktual selalu
   kosong.
3. `ProteksiSendKomiteCNP_Act`: proteksi "Error No Account" tidak pernah aktif. `IsError = 2` tidak pernah direset.
4. `CreateChildKomiteCNP_Act`: pesan "Nilai Gross Value tidak sesuai dengan Spreading In" muncul **setelah** kasus
   komite dibuat.
5. `HitServiceToKasir_Act`: tahun Tgl Boleh Bayar naik menurut bulan **sekarang**, bukan bulan akseptasi.
6. `CloseClaimTNonProp`: kronologi diisi ke `CARI12`, padahal DT membaca `CARI1`.
7. `CountReinstatement_Act`: Adjuster Fee porsi RNM, tetapi Salvage digross-up. Hanya relevan bila grid Reinstatement
   kasus (NEVER) dibangun; lihat OQ-CNP-12.

**Butir OQ-CNP-06 (susunan penyimpanan):**

- **Kasus:** `T_WORK_CLAIM` `LINI = 'NONPROP'`, awalan `CLMNP-` (klaim) dan `KMTNP-` (komite, awalan Pega).
- **Kepala:** `T_GENERAL_CLAIM` ditambah kolom khas Non Prop: Reinsurance Slip, Claim No Ceding, DLA No Ceding,
  Circumstances, Supporting Document, Deductible Min/Max, Waiting For Actual Premium, layer polis.
- **Tabel anak Claim Prop** dipakai ulang bila kolomnya cocok (interest, akseptasi, Spreading In / QS akseptasi) dan
  ditambah kolom XoL.
- **Tabel baru** di rentang `600–639` hanya untuk daftar khas XoL: Loss Allocation (To XOL), XOL Allocation per layer
  (UR / limit / MDP / reinstatement), Claim Acceptation per mata uang, Previously Calculated.
- **OS:** `OS_AKSEPTASI_KLAIM.DATA_JSON` diisi seperti Claim Prop. `CNPLayerList` dihitung saat kirim, tanpa tabel.

**Catatan pelaksanaan (bukan pertanyaan).** Roster `EMAILKOMITE` NONPROP akan diganti workbasket lewat migrasi
claimnonprop 6xx: UPDATE `OPERATOR_ID` keempat baris menurut DEGREE, batas dan jabatan tetap, mengikuti pola migrasi
claimprop 537. Penulisan ke DEV menunggu izin work owner saat tahap bangun.

## B. Pemilik ekspor Pega

| OQ | Rule / makna yang tidak ada di korpus | Akibat bila tidak terjawab |
| --- | --- | --- |
| **OQ-CNP-18** | Arti `Limit` / `Limit2`, `Deductible` / `Deductible2`, `Currency` / `Currency2` di layer treaty. Rumus konversi ke mata uang klaim tidak simetris (IDR dibagi kurs, valas dikali), dan retensi bernilai 0 bila mata uang klaim = mata uang layer. | Dibangun persis XML (`CountLossAllocation_act` s.17), termasuk keanehannya. |
| **OQ-CNP-19** | Activity `GenerateDLACNP` (post-activity FA DLA). | Tombol Generate DLA nonaktif. |
| **OQ-CNP-20** | Harness `InputTreatyInOffer` (View Master). | Tombol View Master nonaktif, atau memakai tampilan master Treaty In yang ada (tanya work owner saat bangun). |
| **OQ-CNP-21** | `CountSpreadingCNP_Act` / `SetTreatyNameSpreading_Act` (ubah share grid Spreading). | Kolom share hanya-baca. |
| **OQ-CNP-22** | Stream HTML dokumen (CFS, PLA, Claim Analysis, persetujuan, CWP, email), When `IsFire`, dan message rule `ErrorMsg1` dkk. | PDF disusun ulang dengan go-pdf/fpdf seperti Claim Prop; tata letaknya perlu contoh dokumen. |
| **OQ-CNP-23** | Harness `New` / `NewSample` (pintu pembuatan kasus); siapa yang me-raise tiket `AcceptanceClaim`. | Tombol New seperti Claim Prop (OQ-CP-13). Tiket diabaikan. |
| **OQ-CNP-24** | Pemetaan kolom hasil `GetDataPolisNonProp_SQL`. Autocomplete membaca CARI5–8 sebagai layer, tetapi popup daftar polis membaca CARI5 = Ceding, CARI6/7 = tanggal. | Diturunkan dari teks SQL saat bangun; bila bertentangan, ikut SQL. |
| **OQ-CNP-25** | Kasus komite CWP tidak mengisi `TransferType`, sehingga router tidak menetapkan tujuan. | Terjawab oleh pola workbasket: tujuan = POSITION tingkat 1. |

## C. DBA

| OQ | Pertanyaan | Akibat hari ini |
| --- | --- | --- |
| **OQ-CNP-26** | Tabel `TREATY_OUT` (layer treaty out untuk PLA retro, `GetLimitTONPPLA`) **tidak ada** di DEV. Yang ada `TREATY_OUT2`, `M_TREATY_OUT`, `M_TREATY_OUT_DETAIL`. Penggantinya yang mana? | Hitung layer PLA retro tidak bisa diuji di DEV. |
| **OQ-CNP-27** | `PROC_GENERATE_SEQUENCE_NUMBER` menerima `TO_DATE(CARI3)`, tetapi tanggal itu tidak pernah diisi di korpus. Isinya tanggal apa (hari ini, DOL, tanggal akseptasi)? | Nomor klaim / akseptasi belum pasti bentuk periodenya. |
| **OQ-CNP-28** | View `CLAIMXOL` **INVALID** di DEV. SQL "Claim History Master ID" membaca `max(claimxol.TANGGAL)`. | Popup riwayat master gagal di DEV. |
| **OQ-CNP-29** | `CLAIMREJECTED` berisi 102 kasus `CLMNP-`, padahal XML Non Prop tidak menulisnya. Siapa penulisnya (trigger, versi lama)? | Reject Non Prop dibangun tanpa menulis `CLAIMREJECTED` (ikut XML). |
| **OQ-CNP-30** | `STS_REJECT = 5` (Save Previously Paid) 0 baris di DEV. Tombol itu pernah dipakai? | Dibangun persis XML. |
| **OQ-CNP-31** | URL `insertClaimFinalOrClosed` / `insertClaimReject` tertulis mati di connector; `M_LINK_SERVICE` DEV tidak punya kategorinya. Nama kunci konfigurasinya? | Efek outbox berhenti terang, seperti Claim Prop. |
| **OQ-CNP-32** | `REINSURANCE.TRLOSS_DETAIL_T` (syarat Kasir: akseptasi sudah terkonversi) tidak terbaca dari akun DEV. | Sama dengan OQ-CP-11: Kasir hanya jalan di produksi. |

## D. Ditemukan saat bangun (09-10-2026)

| OQ | Pemilik | Pertanyaan | Bawaan hari ini | Bukti |
| --- | --- | --- | --- | --- |
| **OQ-CNP-33** | pemilik ekspor Pega | Sel grid panel `AdjustmentDetailNP` (Add Claim Acceptation, Add Loss Allocation, Currency / angka / To XOL) memanggil `AddListClaimNP_Act`, `SetCurrency_Act`, `CountClaimTNP_Act`, `AddLossAlocation_Act`, `CountLossAllocation_act` dengan kelas `ASM-FW-GCNMFW-Data-Adjustment`; activity itu hanya ada di kelas Work. Versi kelas Data-Adjustment ada di Pega? | Tombol Add nonaktif, sel postValue; Rate of Exchange (kelas Work) tetap menghitung tingkat klaim. | `Section/AdjustmentDetailNP.xml` r1 / r2 |
| **OQ-CNP-34** | work owner / DBA | View Payment Status (`DetailPaymentStsCNP`, `GetDetailPaymentStatus_Act`) memanggil layanan REST luar (`M_LINK_SERVICE`). Boleh dipanggil dari aplikasi, dan kunci konfigurasinya apa? Di Input Acceptation tombolnya tanpa aksi. | Tombol nonaktif. | `Section/OutstandingClaim(1).xml` b109 |
| **OQ-CNP-35** | pemilik ekspor Pega | `view.CARI21` (RO medan Currency / Value pane InputDtlInterest) tidak punya penulis di korpus. | Dianggap kosong (medan dapat diisi). | `Section/InputDtlInterest.xml` |
| **OQ-CNP-36** | work owner / DBA | CWP Yes (`CreateChildKomiteCloseNP_Act`) membuat kasus komite tanpa baris akseptasi, padahal `T_GENERAL_KOMITE.ADJUSTMENT_ID` NOT NULL. Kasus komite CWP ditulis bagaimana (kolom nullable lewat DBA, atau rujukan lain)? | Tombol Yes nonaktif. | `Activity/CreateChildKomiteCloseNP_Act.xml`, katalog DEV |
| **OQ-CNP-37** | pemilik ekspor Pega | Properti tanpa penulis di korpus: `OfferFacIn.QuotationData.BusinessOldId` (nomor PLA, LbuID Kasir), `ClaimData.QuotationData.BusinessName`, `ValueAdjustment` (tangga komite), `FlagProrate` (popup reinstatement), `TotalListClaimAmount(IDR)`, `TotalEstimasi`, `ProtectEndDate`, `IsTreatyIn`. Penulisnya rule apa? | Kosong persis XML (tangga memakai nilai 0, nomor PLA tanpa kode bisnis lama). | `PARITAS.md` §8 |
| **OQ-CNP-38** | pemilik ekspor Pega | Baris DEV `OS_AKSEPTASI_KLAIM` STS 0 memuat kunci DATA_JSON `EstimationDate`, `PolicyNo`, `pzInsKey` yang tidak ditulis halaman `SaveDataToOSAksep_Act` ekspor. Rule Pega hidup lebih baru dari ekspor? | Ikut XML ekspor. | `repository/sqldev_test.go`, DEV CLMNP-3998 |
| **OQ-CNP-39** | DBA | Sequence `PLATNP_SEQ` (`GenerateNoPLATNP`) tidak terlihat dari akun DEV (`ALL_SEQUENCES`, `ALL_SYNONYMS`, `ALL_OBJECTS` 09-10-2026; yang ada `PLATREATYIN_SEQ`, `PLA_SEQ`), padahal prompt tahap 1 §4 mencatatnya "ada". Ada di skema lain tanpa hak baca, atau belum dibuat? | Print PLA berhenti terang (409) tanpa nomor. | katalog DEV |
| **OQ-CNP-40** | work owner | When `IsPEGASyariah` (Kasir langkah 9.4: LdcId syariah) memeriksa nama node server Pega, bukan data kasus. Aplikasi baru satu instans: kapan muatan Kasir memakai LdcId syariah? | Selalu konvensional (`ldcId`); parameter `syariah` sudah ada di `SusunMuatanKasir`. | `When/IsPEGASyariah.xml`, `HitServiceToKasir_Act` 9.4 |
| **OQ-CNP-41** | work owner / tim inti | Lanjutan OQ-CNP-04: kata sandi Edit XOL Allocation "ke konfigurasi", tetapi modul dilarang membaca env (ADR-U-0013) dan sandi tidak boleh di repo. Jalur konfigurasi rahasia modul dari inti, atau diganti hak akses? | Tombol Edit XOL Allocation nonaktif. | `EditXOLAlokasi` |
| **OQ-CNP-42** | work owner | `TotalKomite` hanya ditulis `CreateChildKomiteCNP_Act` 26.11 SESUDAH penyerahan, padahal "Send to Committe" tampil hanya bila `TotalKomite != ''` - persis XML tombol itu tak pernah tampil untuk akseptasi baru (Pega hidup jelas menyerahkan: DEV `HISTORYAKSEPTASIPEGA`). | `[penyimpangan sadar]`: TotalKomite = cacah calon tangga saat layar disusun (pola Claim Prop). Konfirmasi atau sebut penulisnya. | `Section/AdjustmentDetailNP.xml` b261, `models/komite.go` |
| **OQ-CNP-43** | work owner | `SaveDataToOSAksep_Act` langkah 8 menolak bila `ClaimData.QuotationData.BusinessOldId` kosong, padahal pengisinya (15.2-15.3) jalan SESUDAHNYA - persis XML Save to issue RNM tidak pernah lolos untuk kasus baru (DEV punya 1.684 kasus ber-OS). | `[penyimpangan sadar]`: OLDID dibaca lebih dulu (GetDataBusiness menurut nama bisnis). Konfirmasi. | `models/outstanding.go` `PeriksaSimpanOS` |
| **OQ-CNP-44** | work owner | Perbaikan OQ-CNP-05 butir 3 diturunkan dari langkah 7.3 / 7.4: konjungsi `.Currency != Curr1` yang meniadakan `Curr1 == .Currency` dibuang, sisanya dipakai (baris Spreading In bermata uang rekening 1 / 2 yang nomor rekeningnya bukan rekening 1 maupun 2 = "Error No Account"). Maksud ini benar? | Dibangun seperti itu (uji `TestProteksiKomiteRekeningSalahDanIsErrorDireset`). | `ProteksiSendKomiteCNP_Act` 7.3 / 7.4 |
