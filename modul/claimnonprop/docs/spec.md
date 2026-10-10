# Spec — Claim Non Prop, tahap 1 dari 2

> 09-10-2026. Klaim treaty inward non proporsional (XoL), kasus `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp`
> (`Flow/Flow_TreatyIn.xml`): Outstanding Claim → Input Acceptation → Resolved-Completed. Sumber: korpus
> `D:\XML\RNM_BRD\Claim Non Prop` (baca-saja), pindai [`PINDAI.md`](PINDAI.md), keputusan work owner di
> [`OQ.md`](OQ.md) dan `MODUL.md`. Status setiap tombol: [`PARITAS.md`](PARITAS.md). Tabel:
> [`STRUKTUR-TABEL-CLAIM-NON-PROP.md`](STRUKTUR-TABEL-CLAIM-NON-PROP.md). Tahap 2 (keputusan komite `KMTNP-`,
> `ShowTransfer`, `KomitePostAdjustment(CWP)`) = modul `komiteclaimnonprop`.

## Masalah

Klaim XoL masih ditangani di Pega lama. Petugas klaim mencatat estimasi per layer XoL, menerbitkan nomor klaim dan OS,
lalu tim teknik memasukkan akseptasi, rekening penerima, Kasir, dan menyerahkan akseptasi ke komite berjenjang.

## Solusi

Modul `claimnonprop` meniru layar dan activity XML apa adanya dengan pola Claim Prop (layar dievaluasi server dari
section XML, satu aksi = satu activity, write-through ke tabel datar, efek luar lewat outbox), ditambah tabel khas XoL.

## Kriteria penerimaan

Setiap butir diuji di seam yang disebut (`models` = uji tabel murni, `handlers` = uji HTTP di atas gudang tiruan
`UJI-`, `ujidev` = SQL baca-saja di DEV).

### Halaman awal

1. Tab Process bawaan = worklist pembuat (Assignment2 "Outstanding Claim"); switch Teknik nyala = workbasket
   `ReasKlaimTeknik` (Assignment1) dan dapat dinyalakan hanya anggota workbasket itu; tab Resolve = kasus selesai.
   (`handlers.TestBuatKasusMasukWorklistPembuat`, `frontend/pages/inbox.test.ts`)
2. Add Claim hanya di tab Process saat switch mati; kasus baru `CLMNP-` di `T_WORK_CLAIM` `LINI = 'NONPROP'`.

### Outstanding Claim

3. Choose Master In membuka pop-up master `PROPORTIONTYPE = 'NonProportional'` (TREATYINDETAIL ∪ TREATYINDETAILEDM),
   saringan per kolom di server, batas 500; Choose membaca ulang baris di server lalu menjalankan `SetValueClaimTNP_Act`.
4. Policy No hanya dari polis master (`GetDataPolisNonProp_SQL`); View List Policy memilih polis yang sama.
5. Pemeriksaan tanggal (DOL, Report, Received, Policy Start / End) memberi pesan XML apa adanya.
6. Claim Amount per mata uang dihitung `CountClaimTNP_Act`; Loss Allocation + XOL Allocation per layer dihitung
   `CountLossAllocation_act` — cocok dengan kasus DEV CLMNP-3998 / master 1001789 (UR 2,4 M; XL1 3,1 M → RNM 930 jt;
   XL2 4,5 M → RNM 1,35 M; Break QS 1,368 M / 912 jt). (`models.TestHitungKlaimSamaDenganKasusDEV3998`,
   `ujidev TestHitungMasterDEV`)
7. Save menyimpan halaman; Save to issue RNM (terbuka hanya bila `IsOutstanding != 1` dan tanpa proteksi tanggal)
   menerbitkan nomor klaim, menulis satu baris OS per layer × mata uang bernilai, JSON_KLAIM, dan efek outbox
   `outstanding-np` (produksi saja). (`handlers.TestAlurOutstandingSampaiInputAcceptation`)
8. Print CFS hanya sesudah Save to issue RNM; Submit hanya sesudah CFS dan memindah kasus ke Input Acceptation
   workbasket `ReasKlaimTeknik`.
9. Print PLA membuka PreviewPLA; Submit menerbitkan nomor PLA dari `PLATNP_SEQ` atau berhenti terang (409) bila sequence
   tidak ada.

### Input Acceptation

10. Save To OS menyalakan Add akseptasi; Add akseptasi menyalin XOL Allocation, Claim Acceptation, Spreading In / Out
    kasus ke akseptasi baru (`AddAkseptasiCNP_Act`); tanpa Adjuster / Consultant ditolak pesan XML.
11. Payable To / Specify mengisi rekening penerima dari `BANKACCOUNT`; Name of Bank memilih rekening yang sah untuk
    akseptasi itu saja.
12. Send to Committe tampil bila tangga calon ada; Send Claim to Committee menolak (422, nol tulisan) bila Gross Value
    layer ≠ Spreading In; bila lolos membuat kasus `KMTNP-` + `T_GENERAL_KOMITE` + `T_KOMITE_KOMITELIST`, tangga = hanya
    DEGREE 1 bila RNM Share ≤ 30 dan ValueAdjustment ≤ 30.000.000 atau subjectivity, selain itu semua tingkat roster
    NONPROP; `POSITION` = workbasket tingkat 1; email ke anggota lewat outbox. Akseptasi di komite terkunci (409).
    (`handlers.TestAkseptasiSampaiKasusKomite`, `models.TestTanggaKomiteHanyaTingkat1`)
13. Acceptation (Kasir) mengantre efek `kasir` hanya di produksi dan hanya untuk akseptasi berstatus disetujui.
14. Save Previously Paid menulis OS `STS_REJECT = 5`; Waiting For Actual Premium menimpa layer terakhir dengan Claim RNM
    50.000; View menampilkan selisih aktual; View Old Allocation / Claim History Master ID / View Payment Attachment
    menampilkan data baca-saja.
15. Close Claim (Remarks wajib) menulis OS STS 4 dan JSON_KLAIM lalu menutup kasus; ditolak selama akseptasi di komite.
    (`handlers.TestCloseClaimMenutupKasusDanMenulisOS4`)
15a. Close Without Payment (NA bila `IsAcceptation = 1`) membuka CloseClaimNP dengan pra-proses `CloseClaimNP_preAct`
    (tanggal, inisial, riwayat). Tombol Yes (kasus komite satu tingkat `ReasClaimDeptHead`, OQ-CNP-02) **nonaktif**:
    kasus komite tanpa baris akseptasi tidak dapat ditulis karena `T_GENERAL_KOMITE.ADJUSTMENT_ID` NOT NULL (OQ-CNP-36,
    prompt §9: tidak mengubah kolom yang ada). (`models.TestCWPYesNonaktif`)
16. Tombol tanpa rule / tanpa jalur sah tampil nonaktif dengan OQ di `title` dan aksinya ditolak server (409).
    (`models.TestTombolOQNonaktif`, `handlers.TestTombolOQTidakDapatDijalankan`)

### Batas

17. Nol stored procedure, nol `COMMIT` di teks SQL; uang / persen `NUMBER(38,10)` dengan `apd` (nol float).
18. Nol penulisan ke DEV dari uji; semua SQL baru dicoba baca-saja di DEV (`ujidev TestSQLDiDEV`).
19. Baris NONPROP tidak masuk kotak masuk Claim Life / Komite Claim Life / Claim Prop / Komite Claim Prop
    (`PARITAS.md` §10).

## Di luar lingkup tahap 1

Layar keputusan komite, tabel komite di inbox, efek persetujuan / penolakan (`komiteclaimnonprop`); pemuat data lama
(OQ-CNP-07); berkas PDF dokumen (OQ-CNP-22); Edit XOL Allocation (OQ-CNP-04 / 41); DLA (OQ-CNP-19); View Master
(OQ-CNP-20); kasus komite CWP (OQ-CNP-36).
