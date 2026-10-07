# 50: Tab Premium Deduction / Brokerage Fee kasus FIRE

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 05-10-2026 "sekarang lanjut ke
> tab Premium Deduction / Brokerage Fee". Survei korpus oleh subagen (pengurai XML, bukan pencacah baris); butir kunci
> diperiksa ulang sesi 0f (`SetRIComm_Act` @1636–1637, @2246–2247, @2623–2624, @1307; Dtl @77982; `SetTAXNote` @655–656;
> `ProtekDeduction2_Act` @355–356).

**Status:** frontend selesai 05-10-2026 (sesi 0f); backend selesai 05-10-2026 (sesi c3, uji hijau). Migrasi 200 ditulis,
menunggu dijalankan work owner — sebelum biner baru.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

### Layar — `Section\InputInwardFacultativeDtl.xml` (tab baris 74726–90306, `!IsLife`, aktif `IsTABActive_Deduction`)

| Elemen | Properti | Tampil | Aksi |
|---|---|---|---|
| % Premium Deduction | `.Policy.Payment.RICommision` (pxNumber 4 dp) | ALWAYS | change/enter → `SetRIComm_Act` |
| % Premium Deduction 2 | `.Policy.Payment.PctDeduction2` | OTHER `never` (tidak ada When `never`) → tersembunyi `[dugaan]` | — |
| % Brokerage Fee | `.Policy.Payment.PctBrokerageFee` | ALWAYS | → `SetRIComm_Act` (Action=pct) |
| centang "Exclude PPN ?" (`pyCheckboxCaption`) | `.OfferFacIn.PPnCheck` | ALWAYS | → `SetRIComm_Act` |
| Total Premium Deduction / Total Premium RNM | `.OfferFacIn.TotalRIComNusaRe` / `.TotalPremiNusaRe` | OTHER `1=2` → tersembunyi | — |
| section `CoverageCommisionList` | FIRE: Risk Location → item → coverage → `CommisionItem` | ALWAYS | drill-down baca-saja, tanpa activity |
| section `SummarySpreading_Section` | ringkasan Spreading | ALWAYS | — |
| grid 1 `.OfferFacIn.CurrencyList` | Currency · Total Premium Deduction (`.Policy.Payment.Commision`, dapat diisi) · Total Premium Deduction 2 (`.Deduction2`, dapat diisi) | ALWAYS | → `SetTotalPremium_Act` / `ProtekDeduction2_Act` |
| grid 2 `.OfferFacIn.CurrencyList` | Currency · Total Brokerage Fee (`.BrokerageFee`, dapat diisi) · PPh · PPN (baca-saja; kepala = `.PPH_Note` / `.PPN_Note` = "PPH n %" / "PPN n %", `SetTAXNote`) | ALWAYS | → `TotalBrokerageFee_Act` |
| Save | `SaveFacIn_Act` (Obj-Save) | ALWAYS | — |

### Rumus per mata uang (Premium = `.Policy.Payment.Premium` = Σ PremiNusantaraRe mata uang itu, `CountPremiAndTSIRNMFireMBU_ACT` 6.2.5)

- `SetRIComm_Act`: Commision = RIComm × Premium / 100; BrokerageFee = PctBrokerageFee × Premium / 100; PPh = tarif PPH ×
  dasar / 100; PPN = tarif PPN × dasar / 100. Dasar = BrokerageFee bila `PPnCheck==true`, = BrokerageFee / ((PPN + 100) /
  100) bila `false` — tetapi langkah 5 **selalu** menyetel `PPnCheck = true` (@1636–1637), jadi cabang `false` tak pernah
  jalan. Lalu `SaveFillPaymentInstallment` bila installment terisi.
- `SetTotalPremium_Act`: % Premium Deduction = Commision ketikan ÷ Premium baris itu × 100; semua baris dihitung ulang.
- `ProtekDeduction2_Act`: Deduction2 > Premium → "Deduction 2 cannot be greater than the Premium" (berhenti); else
  PctDeduction2 = Deduction2 ÷ Premium × 100; semua baris Deduction2 = PctDeduction2 × Premium / 100.
- `TotalBrokerageFee_Act`: % Brokerage Fee = BrokerageFee ketikan ÷ Premium × 100; semua baris BrokerageFee, PPh, PPN
  dihitung ulang (cabang PPN membaca `pyWorkPage.PPnCheck`).
- Tarif PPH / PPN = `Obj-Browse` kelas `ASM-FW-GISFW-Int-TAX`, `TANGGAL` turun, 1 baris, tanpa saringan. Tabel fisiknya
  `POOLDATA.TAX` (DDL dikirim work owner 05-10-2026, lihat Jawaban butir 2).
- `SaveFillPaymentInstallment`: NetPremium = Premium − Commision − Deduction2 − BrokerageFee + PPh + PPN (tab Payment).

## Jawaban work owner 05-10-2026

1. Gambar layar Pega dikirim (cocok: Σ Total Premium RNM 469.550 × 35% = 164.342,5).
2. DDL `D:\migrasi\RNM\DDL\TAX.txt`: `POOLDATA.TAX (TANGGAL DATE, PPH NUMBER, PPN NUMBER)`.
3. "Exclude PPN ?" — "iya buatin kaya dipega": tampil, dan backend menyetel `ppnCheck = true` tiap hitung (PPh / PPN dari
   Brokerage Fee penuh).
4. Daftar komisi — "dibuatin seperti di pega": lokasi → item → coverage → detail, baca-saja.

## Frontend (sesi 0f)

- `components/TabDeduction.tsx` (`kepalaPajak`, `lebihBesar`, `deduction2Lebih`) + uji; `TabSpreading.tsx` mengekspor
  `RingkasanSpreading` (dipakai ulang); `api.ts` (`ambilDeduction`, `hitungDeduction`, `simpanDeduction`,
  `TampilanDeduction`, `BarisDeduction`, `AksiDeduction`); `labels.ts` (`DEDUCTION`, `TEKS_DEDUCTION`); `nbfacin.css`;
  `pages/InwardFacultative.tsx`.
- Daftar komisi + ringkasan memakai `GET …/spreading` (pohon lokasi → item → coverage yang sama).

## Kontrak USULAN (menunggu sesi c3)

- `GET /api/nbfacin/kasus/{caseId}/deduction` → `TampilanDeduction` { pctPremiumDeduction, pctBrokerageFee, ppnCheck,
  tarifPph, tarifPpn, mataUang[{currency, premium, commision, deduction2, brokerageFee, pph, ppn}], pesan[] }.
- `POST …/deduction/hitung` { aksi: persen | commision | deduction2 | brokerage, idx, pctPremiumDeduction,
  pctBrokerageFee, ppnCheck, mataUang[{currency, commision, deduction2, brokerageFee}] } → tampilan (tanpa simpan).
- `PUT …/deduction` (badan sama tanpa aksi / idx) → hitung ulang + simpan.

## Backend (sesi c3, 05-10-2026)

Kontrak final = kontrak usulan di atas. `GET /api/nbfacin/kasus/{caseId}/deduction` → `TampilanDeduction` persis `api.ts`
(larik selalu `[]`, angka teks desimal 8 angka di belakang koma, `""` bila kosong). `POST …/deduction/hitung` {aksi, idx,
pctPremiumDeduction, pctBrokerageFee, ppnCheck, mataUang[{currency, commision, deduction2, brokerageFee}]} → tampilan, TANPA
simpan. `PUT …/deduction` (badan sama TANPA aksi / idx) → simpan → jawab bentuk GET. Kunci asing → 400 (badan dibangun
`badanDeduction`, tepat). Galat: 400 isian / aksi / idx / "Deduction 2 cannot be greater than the Premium"
(`ErrMasukanDeduction`), 401 PUT tanpa identitas, 404 case, 409 mata uang badan ≠ mata uang objek, 503 tanpa basis data.

- `backend/migrations/200_premium_deduction.sql` (+ `_down`): `T_GENERAL_POLIS.PPN_CHECK` VARCHAR2(10), `PAY_RI_COMMISION` /
  `PAY_PCT_BROKERAGE_FEE` NUMBER(38,8); tabel `T_CURRENCYLIST` (rancangan, sebagian) + indeks + sequence. **Ditulis, belum
  dijalankan** — wajib sebelum biner baru. Amandemen rancangan `loader/amandemen.go` `T_CURRENCYLIST.PAY_DEDUCTION2` (K50-2).
  Penjaga rancangan membaca 200 (251 kolom); `docs/STRUKTUR-TABEL-NB-FACIN.md`.
- `models/deduction.go`, `repository/deduction.go` (`PenyimpanDeduction`, tarif `POOLDATA.TAX`), `services/deduction.go`,
  `handlers/deduction.go` (+ uji ketiganya).
- Port `[terverifikasi]` (dibaca ulang sesi c3, `Activity\`): `SetRIComm_Act` langkah 4 (Obj-Browse TAX, Sort D, MaxRecords 1),
  5 (`PPnCheck = true`), 8.1 / 8.3 / 8.4–8.6 (rumus); 8.2 dan 9 ber-`pyStepsBlockName` `//`. `SetTotalPremium_Act` 1–2,
  `ProtekDeduction2_Act` 2–5, `TotalBrokerageFee_Act` 2–3, 6, 8.1, 8.3–8.5 (8.2 IsEDM).

## Keputusan agent backend (menunggu konfirmasi)

- **K50-1** (c3) Penyimpanan = workbook rancangan `[terverifikasi]`: persen di `T_GENERAL_POLIS` (`PAY_RI_COMMISION`,
  `PAY_PCT_BROKERAGE_FEE` = `.PolicyData.Payment.*`, lipatan V-22b; contoh `DDL\CONTOH` menaruhnya di `/PolicyData/Payment`
  115 / 81 kali), centang `PPN_CHECK` (contoh: `/PPnCheck` 74 kali; teks Pega apa adanya "true" / "false", VARCHAR2(10) — butir 68.1 work owner
  02-10-2026 sudah menggolongkan PPN_CHECK kolom teks menyimpang, `loader/aturan.go`; versi pertama keliru menulis NUMBER(5) 1 / 0,
  diralat 06-10-2026 sebelum dijalankan), potret per mata uang di
  `T_CURRENCYLIST` (`PAY_PREMIUM` / `PAY_COMMISION` / `PAY_BROKERAGE_FEE` / `PAY_PPH` / `PAY_PPN`, jalur `CurrencyList`,
  induk T_GENERAL_POLIS). `PctDeduction2` tidak disimpan (tidak di rancangan, 0 contoh; diturunkan dari Deduction2 / Premium
  saat aksi).
- **K50-2** (c3) Amandemen rancangan `T_CURRENCYLIST.PAY_DEDUCTION2` NUMBER (`loader/amandemen.go`, uji ukuran loader +1):
  `.CurrencyList(n).Policy.Payment.Deduction2` ada di layar (grid "Total Premium Deduction 2", dapat diisi) tetapi tidak di
  rancangan — `[terverifikasi]` 0 dari 115 contoh memuatnya di tingkat CurrencyList.Policy.Payment; seluruh 123 tag
  `Deduction2` ada di ListInstallment (101 CurrencyList, 21 OldData, 1 FacRetroList), dihitung dua cara (pengurai XML = grep
  tag, 05-10-2026).
- **K50-3** (c3) Premium per mata uang **dihitung tiap baca** dari objek tersimpan (Σ premi NR item bermata uang itu, % Share
  RNM tersimpan; `ringkasanRNM`), tidak dari kolom. GET memakai Commision / BrokerageFee / PPh / PPN tersimpan bila
  `PAY_PREMIUM` = premi sekarang (8 desimal); bila premi berubah (Save tab Object / Coverage / Spreading sesudahnya) nilainya
  dihitung ulang dari persen tersimpan — setara Pega yang menjalankan SetRIComm_Act di setiap Save tab lain
  (IsThereAnyObjectLocation_Act langkah 11). Deduction2 = tersimpan apa adanya (SetRIComm tidak menyentuhnya); bila kini >
  premi → pesan verbatim (tidak menghalangi tampil). `[pertanyaan terbuka]` **Keanggotaan dan urutan CurrencyList**: di sini
  satu baris per mata uang item objek, urut kemunculan pertama; item tanpa mata uang tidak punya baris. Di Pega
  `CountPremiAndTSIRNMFireMBU_ACT` TIDAK membangun CurrencyList (langkah 4–5 membaca daftar yang sudah ada, 6.2.5 hanya
  mengisi Premium baris ber-`.Name` sama) — daftarnya dibangun di tempat lain yang belum dibaca, sehingga mata uang
  terdaftar tanpa item (premi 0) di Pega tidak muncul di sini. Premi = Σ PremiNusantaraRe seluruh coverage (`ringkasanRNM`):
  syarat Pega `FlagDelete != "1"` dan premi basis 5 dari LayerList tidak dibawa (K48-6).
- **K50-4** (c3) `@Math.divide(a, b)` tanpa skala (SetTotalPremium_Act langkah 1, ProtekDeduction2_Act langkah 4) ditulis
  **skala 20** seperti saudaranya (TotalBrokerageFee_Act langkah 6 `@Math.divide(a, b, 20) × 100`) — skala bawaan Pega
  `belum terverifikasi`.
- **K50-5** (c3) Premium baris 0: Pega tanpa penjaga (bagi nol); di sini aksi commision / deduction2 / brokerage atas baris
  itu → 400 "Premium <mata uang> 0 - <medan> tidak dapat dihitung menjadi persen" (pesan agent). Aksi persen / PUT tetap
  jalan (hasil 0).
- **K50-6** (c3) `ppnCheck` badan diterima tetapi diabaikan: setiap hitung dan PUT memakai `true` (SetRIComm_Act langkah 5,
  jawaban work owner butir 3); GET menampilkan nilai tersimpan (`false` sebelum pernah disimpan). Cabang `PPnCheck == false`
  (dasar = Brokerage / ((PPN + 100) / 100), 8 desimal) diport dan diuji tetapi tak terjangkau.
- **K50-7** (c3) Aksi commision / deduction2 juga menghitung ulang PPh / PPN dari Brokerage Fee badan (Pega tidak menyentuh
  PPh / PPN di dua activity itu; nilainya sama karena Brokerage Fee tidak berubah). Aksi deduction2 tidak mengubah persen
  atas. Tarif TAX kosong → PPh / PPN 0, `tarifPph` / `tarifPpn` "" + pesan agent. Baris TAX dipilih `ORDER BY TANGGAL DESC`
  apa adanya (baris TANGGAL NULL lebih dulu, perilaku bawaan Oracle seperti SQL Pega `[dugaan]`).
- **K50-9** (c3, temuan review spec) Persen tampil / tersimpan 8 desimal, padahal Pega menyimpan persen hasil aksi baris
  berpresisi penuh di clipboard. Tanpa perbaikan, Save menghitung ulang dari persen 8 desimal dan nilai yang diketik
  bergeser (Brokerage Fee 100000 atas premi 469550 tersimpan 100000.00001684). Kini hitung dan PUT membangun ulang persen
  berpresisi penuh dari baris grid pertama berpremi bukan nol bila persen kiriman = pembulatan 8 desimalnya (`persenPenuh`);
  bila tidak (persen diketik ulang), persen kiriman dipakai apa adanya. Uji `TestDeductionTidakBergeser`. GET sesudah premi
  berubah tetap menghitung dari persen tersimpan 8 desimal (geseran ≤ premi × 5·10⁻¹¹).
- **K50-10** (c3) Koma desimal: `% Premium Deduction` menerima "10,5" = 10.5 (`@replaceAll(RICommision, ",", ".")`,
  SetRIComm_Act langkah 5 / TotalBrokerageFee_Act langkah 3); medan lain tetap titik (Pega tidak mengganti koma di sana).
- **K50-11** (c3, tiket 51) Save tab ini kini memperbarui baris `T_CURRENCYLIST` menurut nama mata uang (bukan hapus-sisip)
  supaya cicilan tab Payment (`T_LISTINSTALLMENT`, migrasi 201) tidak ikut terhapus; mata uang yang hilang dibuang beserta
  cicilannya. ⚠️ Karena itu migrasi 201 wajib jalan bersama 200 sebelum biner baru.
- **K50-12** (c3, temuan review spec tiket 51) Save tab ini juga menghitung ulang cicilan tab Payment yang tersimpan di
  transaksi yang sama (persen dan Due Date tetap; nilai uang cicilan, Total Payment, Net Premium dari nilai Payment baru) —
  Pega menjalankan SaveFillPaymentInstallment + CountPaymentInstallment_Act di setiap Save, jadi kolom tersimpan tidak basi
  bagi pembaca Oracle langsung. Tanpa Installment / cicilan, atau bila Net Premium negatif, tidak ditulis (GET tetap
  menghitung ulang). Proteksi Payment tidak dijalankan di sini.
- **K50-13** (c3, 06-10-2026) **Migrasi 202** `202_ppn_check_teks.sql` (+ `_down`): `200_premium_deduction` SUDAH jalan di DEV
  dengan `PPN_CHECK NUMBER(5)` (T_MIGRASI + ALL_TAB_COLUMNS diperiksa work owner) sebelum ralat teks; berkas 200
  dikembalikan ke bentuk yang jalan (pra-terbang tidak memeriksa langkah tercatat — mengubahnya membuat lingkungan baru
  berbeda dari DEV), dan perubahan tipe lewat 202: tambah `PPN_CHECK_TEKS` → salin (1 → 'true', 0 → 'false', NULL tetap) →
  buang `PPN_CHECK` (berpelindung katalog) → tambah `PPN_CHECK VARCHAR2(10)` → salin balik → buang `PPN_CHECK_TEKS`. Bukan
  MODIFY (ORA-01439 pada kolom berisi) dan bukan RENAME (tak terlihat penjaga STRUKTUR, pola 059 premiumlistlife). Penjaga
  rancangan nbfacin kini membaca migrasi pernyataan demi pernyataan termasuk DROP COLUMN. ⛔ Sampai 202 jalan, Save tab
  Premium Deduction di DEV gagal ORA-01722 (kode menulis "true" ke kolom NUMBER).

  Langkah DEV (work owner / DBA, BUKAN agent):
  1. Sebelum 202 — isi dan tipe sekarang:
     `SELECT COLUMN_NAME, DATA_TYPE, DATA_PRECISION, DATA_SCALE FROM ALL_TAB_COLUMNS WHERE OWNER = '<SKEMA>' AND TABLE_NAME =
     'T_GENERAL_POLIS' AND COLUMN_NAME IN ('PPN_CHECK', 'PPN_CHECK_TEKS');` (harap: PPN_CHECK NUMBER 5 0, tanpa PPN_CHECK_TEKS)
     `SELECT ID, PPN_CHECK FROM <SKEMA>.T_GENERAL_POLIS;` (NULL / 1 / 0; nilai lain selain 1 / 0 menjadi NULL oleh 202).
  2. `-migrate` (202 jalan sendiri; urut nama sesudah 201).
  3. Sesudahnya: kueri 1 lagi → PPN_CHECK VARCHAR2 10, tanpa PPN_CHECK_TEKS; `SELECT ID, PPN_CHECK FROM <SKEMA>.T_GENERAL_POLIS`
     → 'true' / 'false' / NULL sesuai isi sebelumnya; `SELECT NAMA FROM <SKEMA>.T_MIGRASI WHERE NAMA = '202_ppn_check_teks'`.
  Bila 202 gagal di tengah (berhenti di ORA-01430 saat diulang): lihat kueri 1 — (a) ada PPN_CHECK NUMBER + PPN_CHECK_TEKS:
  `ALTER TABLE <SKEMA>.T_GENERAL_POLIS DROP COLUMN PPN_CHECK_TEKS` lalu ulangi; (b) tanpa PPN_CHECK, ada PPN_CHECK_TEKS:
  `ALTER TABLE <SKEMA>.T_GENERAL_POLIS ADD (PPN_CHECK VARCHAR2(10))`, `UPDATE <SKEMA>.T_GENERAL_POLIS SET PPN_CHECK =
  PPN_CHECK_TEKS`, `ALTER TABLE <SKEMA>.T_GENERAL_POLIS DROP COLUMN PPN_CHECK_TEKS`, `COMMIT`, lalu catat
  `INSERT INTO <SKEMA>.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('202_ppn_check_teks', SYSDATE)`; (c) PPN_CHECK VARCHAR2 +
  PPN_CHECK_TEKS: salin balik bila perlu, buang PPN_CHECK_TEKS, catat T_MIGRASI seperti (b).

  Verifikasi kolom 201 di DEV (berkas 201 tidak diubah sejak ditulis 06-10-2026 10:17):
  `SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, DATA_LENGTH, DATA_PRECISION, DATA_SCALE, NULLABLE FROM ALL_TAB_COLUMNS WHERE
  OWNER = '<SKEMA>' AND ((TABLE_NAME = 'T_GENERAL_POLIS' AND COLUMN_NAME = 'PAY_INSTALLMENT') OR (TABLE_NAME =
  'T_CURRENCYLIST' AND COLUMN_NAME IN ('SUM_TOTAL_PAYMENT', 'PAY_NET_PREMIUM')) OR TABLE_NAME = 'T_LISTINSTALLMENT') ORDER BY
  TABLE_NAME, COLUMN_ID;` — harap: PAY_INSTALLMENT NUMBER 5 0; SUM_TOTAL_PAYMENT / PAY_NET_PREMIUM NUMBER 38 8;
  T_LISTINSTALLMENT 18 kolom: ID / PARENT_ID NUMBER 19 0 NOT NULL, SEQ_NO NUMBER 5 0 NOT NULL, ROW_UID VARCHAR2 36 NOT NULL,
  CURRENCY_CODE VARCHAR2 10 NOT NULL, INSTALLMENT_NO NUMBER 5 0, DUE_DATE VARCHAR2 30, INSTALLMENT_PERCENTAGE VARCHAR2 50,
  PREMIUM / RI_COMMISION / DEDUCTION2 / BROKERAGE_FEE / PPH / PPN / DISCOUNT / PAYMENT_TOTAL NUMBER 38 8, IDPEGA VARCHAR2 50,
  COB_GROUP VARCHAR2 20.
- **K50-8** (c3) Di luar cakupan: cabang EDM (SetRIComm 1, 6; TotalBrokerage 7, 8.2), IsMultiCOB, `TotalRIComNusaRe`
  (langkah 9 `//`), `SaveFillPaymentInstallment` / NetPremium (tab Payment), PctDeduction2 di layar (tersembunyi, D-3).

## Keputusan agent (menunggu konfirmasi)

- **D-1** isian dihitung saat lepas fokus / Enter (event change / enter Pega), bukan tiap ketikan.
- **D-2** "Deduction 2 cannot be greater than the Premium" diperiksa juga di layar sebelum memanggil server.
- **D-3** bagian tersembunyi Pega (% Premium Deduction 2, Total Premium Deduction, Total Premium RNM) tidak ditampilkan;
  `SaveFillPaymentInstallment` (installment) = tab Payment.

## Pertanyaan terbuka (sisa)

1. Arti `pyStepsPreCondition=false` dengan When tersimpan (dianggap mati) dan `@Math.divide(a,b)` tanpa skala —
   `belum terverifikasi`.
