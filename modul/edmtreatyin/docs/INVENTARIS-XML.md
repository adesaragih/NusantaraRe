# Inventaris XML — EDM Treaty In

> Disusun 06-10-2026 dari korpus `EDM Treaty In` (163 berkas `.xml`, 12 tipe; dibaca saja) lewat laporan pindai
> sesi ini, lalu **setiap** baris dicocokkan ke kode modul: nama rule dicari di komentar kode yang mengutip langkah
> XML (mis. `EDMChooseBusiness_Act 5`), fungsinya dibaca, dan pemanggilnya diperiksa. Bukan berkas bangkitan — ubah
> bersama kodenya. Jalur relatif terhadap `modul/edmtreatyin/`. Nol nama orang, ID operator, nomor polis, atau
> alamat layanan; nilai tertanam di XML disebut jenisnya saja.
>
> Status kolom *Dibangun?*: ✅ dibangun · 🟡 sebagian (bagian yang tidak dibangun disebut) · ⛔ tidak dibangun.
> Hasil per acceptance criteria: `docs/HASIL-IMPLEMENTASI.md`. Nomor `#n` = baris log `docs/KOREKSI-DOKUMEN-2026-10-06.md`.

## Isi

1. Ringkasan — cacah, titik masuk, penyimpangan sadar
2. Tabel seluruh rule (163)
3. Rule dan temuan yang diragukan

## 1 · Ringkasan

### 1.1 Cacah per tipe

| Tipe | Berkas | ✅ dibangun | 🟡 sebagian | ⛔ tidak |
| --- | ---: | ---: | ---: | ---: |
| `Activity` | 66 | 45 | 5 | 16 |
| `RDBList` | 36 | 22 | 1 | 13 |
| `Section` | 18 | 18 | 0 | 0 |
| `When` | 11 | 0 | 0 | 11 |
| `DataTransform` | 10 | 8 | 1 | 1 |
| `ReportDefinition` | 9 | 6 | 1 | 2 |
| `FlowAction` | 6 | 6 | 0 | 0 |
| `DecisionTable` | 1 | 1 | 0 | 0 |
| `Harness` | 3 | 3 | 0 | 0 |
| `Flow` | 1 | 0 | 1 | 0 |
| `ConnectREST` | 1 | 0 | 1 | 0 |
| `SystemSettings` | 1 | 0 | 0 | 1 |
| **Jumlah** | **163** | **109** | **10** | **44** |

Sebab rule ⛔ (menurut kolom *Sumber alasan*, sebab pertama):

| Sebab | Rule |
| --- | ---: |
| tak terjangkau / tak terpicu | 21 |
| efek dibaca nol rule | 4 |
| langkah XML `//` | 3 |
| penyimpangan sadar | 8 |
| keputusan WO | 8 |
| **Jumlah** | **44** |

*Tak terjangkau* termasuk pemicu yang tidak pernah menyala: sel ber-readOnly / disabled, wadah `1=2`, dan syarat yang
tidak pernah benar (mis. `PolicyNo = ""` sesudah `SetEDMTNoPolis`). Setiap ⛔ menyebut pemanggil atau langkahnya.

### 1.2 Titik masuk

| Urutan | Rule Pega | Padanan |
| ---: | --- | --- |
| 1 | `Harness/SFAPortal_Endorsement_Treaty` + `Section/SFAPortal_Endorsement_Treaty` — grid RD `InboxEDM_RD2`, tombol *Create New Addendum Treaty* | `frontend/pages/PortalEDMTreatyIn.tsx` (satu halaman menu, `frontend/menu.ts`); `GET /kasus` → `backend/repository/kasus.go` `sqlDaftarKasus` |
| 2 | `Harness/TreatyCreateEdm` + `Section/TreatyCreateEdm` — No Polis Treaty (`TrtEdmCheckPolicyError`, `CheckNopolisAvailability`), No Master Treaty, Source of Change (EDMType) | `frontend/components/BuatEDM.tsx`; `GET /periksa-polis` → `backend/services/layanan.go` `PeriksaPolis` |
| 3 | `Activity/CreateEDMT` (tombol Create) | `POST /kasus` → `backend/services/layanan.go` `BuatKasus` (satu transaksi: generasi terakhir → ID `EDMT-<n>` → `T_WORK_POLIS` + generasi `PRODKE` + 1, `OLD_POLIS_ID`, `NOENDORS`, `EDM_TYPE` → halaman lahir `RakitHalamanBaru`) |
| 4 | `Flow/InputAddendumTreatyIn` — Assignment2 Admin (FA `InboxPolicyTreatyInAddendum`), Assignment4/7 Sec Head dan Assignment1 Dept Head (FA `DeptHeadTreatyIn_UWAddendum`), Utility1 `SaveJsonPolisTreatyInEDM_Act`, Utility2 `serviceInsertArasapas_act` | `GET/PUT /kasus/{id}`, `POST /kasus/{id}/hitung`, `/bisnis`, `/pilih-bisnis`, `/kirim` → `backend/services/tindakan.go`; tangga `backend/models/tangga.go` `Langkah` |

✅ **Kotak masuk Beranda** (keputusan WO 07-10-2026 butir 1 *"YA"*) — `frontend/menu.ts` mendaftarkan `antreanBeranda`
(`GET /kotak-masuk`) dan `daftarBeranda` (`GET /kotak-masuk/kasus`): Sec Head / Dept Head membuka berkas EDM dari
Beranda (portal XML hanya menautkan kasus berposisi Admin). Terbukti cek Chrome 07-10-2026: Beranda → ReasTreatyInSecHead →
EDM Treaty In → kasus → Accept → Submit.

⚠️ **Penyimpangan atas keputusan WO 07-10-2026** (bukan XML): portal ber-switch *In Progress* / *Resolved* (aturan portal NB;
XML satu grid ber-filter C), dan sel header With Tax / Type Tax / Overiding Commision / Marketing Officer terkunci bagi
atasan (XML `pyEditOptions` Auto, tidak dikunci per posisi). Label EDMType = DT `TreatyEDMListType` dari screenshot WO
(DT itu tidak ada di korpus 163 berkas).

🔜 Migrasi 360-363 (tabel proyeksi selisih) dan slot menu 970 **belum dijalankan** di DEV (`-migrate` oleh WO). Setiap
penyimpanan sesudah Create (Save, Choose Business, Submit) menulis `T_POLIS_DIFFERENCE*`, jadi di DEV baru dapat berjalan
sesudah migrasi itu.

### 1.3 Penyimpangan sadar

| # | Penyimpangan | Dasar | Rule (no.) |
| ---: | --- | --- | --- |
| 1 | **Tangga selalu tiga tingkat.** `CekLimitTreatyAcc_Act`, `LetterNo`, `ToTREATYDEPTHEAD`, cabang Decision9 Else (Sec Head menyelesaikan sendiri bila premi ≤ 200 juta) tidak dibangun; Sec Head menyetujui selalu naik ke Dept Head. | prompt eksekusi WO; #42; NB K2 | 2, 116, 130, 161 |
| 2 | **Identitas operator tertanam diganti posisi dan data.** `IsSPVCreate` (dua ID operator), `IsSPVTreaty1` / `IsTreaty1` (kolom telepon), syarat tombol Submit S18 ber-ID operator, nama orang di NBStatus (DT pasca-admin 3, connector) → posisi / workbasket, nama dari `M_LOGIN_GO` (`NamaKotakMasuk`). | spec AC 6; pola NB K5, P40 | 29, 108, 124, 125, 127, 135, 161 |
| 3 | **Selisih proporsional satu rumus** *baru − lama* untuk semua generasi; varian `HasEDMNo` (`EDMTCalculateTreatyDifference` 4-6) tidak ditiru. XOL: `CalculateDifferenceEDM_act` apa adanya. `Total*` diturunkan, tidak disimpan. Tabel 360-363 = proyeksi (`SUMBER` `'PEGA'` beku / `'GO'`). | keputusan WO 23-09-2026 (ID-28/30) | 1, 22 |
| 4 | **Selisih proporsional dihitung ulang setiap generasi ditulis** (`HitungSelisihGenerasi`), juga sebelum baris produksi Utility1 — Pega hanya menghitungnya lewat tombol / sel. | tinjauan kode 06-10-2026; ID-26 | 22, 45 |
| 5 | **Admin menolak → generasi dilepas**: `OLD_POLIS_ID` NULL (`LepasGenerasi`) dan proyeksi selisih `'GO'` dibuang (`BuangSelisihGenerasi`); polis dapat diendorse ulang dengan `PRODKE` yang sama. | ketetapan modul (MODUL.md); tinjauan kode 06-10-2026 | 161 |
| 6 | **Tanpa JSON.** OldData = generasi `OLD_POLIS_ID` di `T_GENERAL_POLIS_TREATY` (`FetchPolisJsonPolis`, `SelectProdKe` diganti; `PRODKE` bilangan, bukan teks terpotong). Utility1 menulis json_polis **tanpa** `DATA_JSON` + `ACHIEVEMENT` + `TREATYINPRODUCTION`, nol prosedur. | keputusan WO 06-10-2026 | 45, 75, 95, 97 |
| 7 | **Create ditolak bila json_polis lebih maju** dari rantai relasional (generasi Pega belum dimuat pemuat) — pesan bukan teks XML (`PesanGenerasiBelumDimuat`). | tinjauan kode 06-10-2026 | 20 |
| 8 | **SuggestList → `HISTORYAKSEPTASIPRODUCTION` `TYPE_POLIS 'EDMT'`** — korpus EDM menyimpannya di blob kasus (nol rujukan, #53). | pola NB K4 | 132 |
| 9 | **Konversi Arasapas tidak disambung** (sama NB): muatan `noPolis / caseId / tglInput` disusun, pengirim bawaan gagal terang, pesan `FlagErrorKonversi`; langkah 7 (FacOut) dan 11-20 tidak dibangun. | ketetapan NB; butir WO #49, #50 | 31, 46, 66, 71, 72, 91, 102, 122, 123, 126, 128, 162, 163 |
| 10 | **Status dan NBStatus dari data**: status tutup `Resolved-Rejected` / `Resolved-Completed` (XML tidak menulis `pyStatusWork`); NBStatus kembali ke Admin = nama pembuat (XML tidak menulis); `OperatorName` = nama tampilan. | NB P24, P33; keputusan WO NB 06-10-2026 | 133, 137, 161 |
| 11 | **Wewenang = keanggotaan antrean**: `TempEmail.CARI28` (urutan workbasket) dan `IsUW` tidak dibangun. | spec AC 51; pola NB AC 14 | 129, 138 |
| 12 | **Baris spreading endorsemen tidak dapat dihapus**: grid tanpa Delete; generasi yang kehilangan baris spreading generasi lama ditolak saat Submit (`BarisSpreadingHilang`, pesan bukan teks XML). | keputusan WO spec-penyimpanan ID-15, ID-16 | 110 |
| 13 | **Label EDMType 1-4 tidak ada di korpus** (DT `TreatyEDMListType` tidak terekspor) — kode ditampilkan apa adanya. | korpus; menunggu screenshot WO | 120 |

### 1.4 Arti kolom tabel bab 2

- **Padanan / alasan** — untuk ✅ / 🟡: berkas dan fungsi / komponen pemegang perilakunya (langkah XML yang tidak
  dibangun disebut); untuk ⛔: pemanggil dan sebab rule tidak pernah berefek, atau keputusan yang menahannya.
- **Sumber alasan** — `//` (langkah XML ber-remark), *tak terjangkau* (pemicu mati / syarat tak pernah benar),
  *tanpa pembaca* (efeknya dibaca nol rule), *keputusan WO*, *penyimpangan sadar*, *menunggu WO*; `—` = tidak ada
  bagian yang ditinggalkan.

## 2 · Tabel seluruh rule (163)

Urutan = urutan tipe lalu nama (sama dengan tabel banding laporan pindai).

| No | Tipe | Nama | Dibangun? | Padanan Go/TS atau alasan | Sumber alasan |
| ---: | --- | --- | :---: | --- | --- |
| 1 | Activity | `CalculateDifferenceEDM_act` | ✅ | `backend/models/edm_xol_selisih.go` `CalculateDifferenceEDM` (EDMChooseBusiness_Act 6); proyeksi lapisan `backend/models/katalog_selisih.go` `DatarSelisihLapisan`, induk turunan `BangunIndukSelisihXOL` (tidak disimpan). Rumus XML apa adanya: batas bawah 0, prorata, pajak, pembatalan mentah; langkah 1.2.2 `//` | XOL apa adanya (MODUL.md); `//` 1.2.2 |
| 2 | Activity | `CekLimitTreatyAcc_Act` | ⛔ | Tangga SELALU tiga tingkat: LetterNo dan batas 200 juta tidak dibangun. Action set radio Approval `ListSuggestEDM` dijalankan tanpa langkah ini (`frontend/components/Usulan.tsx`, `backend/services/aksiposisi.go`) | penyimpangan sadar (prompt WO; koreksi #42; NB K2) |
| 3 | Activity | `CheckDataMkt` | ✅ | `backend/models/pilihbisnis.go` `TerapkanMO` + aksi `CheckDataMkt` di `backend/services/tindakan.go` `aksiHitung` (Obj-Save langkah 5 ditiru) + `backend/repository/acuan.go` `MO` | `//` 2 langkah |
| 4 | Activity | `CheckDuplicateOffer` | ⛔ | Pemanggil tunggal SetTreatyIn_Act 14; langkah 1-4 `//`; pesannya di halaman `TreatyIn` tingkat atas yang tidak disimpan maupun ditampilkan (komentar `pbSetTreatyIn`, `backend/models/edm_pilih.go`) | `//` + tanpa pembaca |
| 5 | Activity | `CheckNopolisAvailability` | ✅ | `backend/repository/generasi.go` `CacahGenerasiJSONPolis` (ada = cacah > 0) + `backend/services/layanan.go` `PeriksaPolis`; layar `frontend/components/BuatEDM.tsx`. Panggilan CreateEDMT 4 `//` - pemeriksaan hidup hanya di layar Create | — |
| 6 | Activity | `ConvertHistoryDate` | ⛔ | SetTreatyIn_Act 12 mengubah `TreatyIn.CommentList` halaman master tingkat atas; nol pembaca di rantai EDM | tanpa pembaca |
| 7 | Activity | `CopyGeneralDataEDM_act` | ✅ | `backend/models/edm_pilih.go` `CopyGeneralDataEDM` (EDMChooseBusiness_Act 4; `IDCurrency = .Currency` ditiru) | — |
| 8 | Activity | `CountNetPremi_act` | ✅ | `backend/models/hitung.go` `CountNetPremi`; dihitung ulang server di `backend/services/tindakan.go` `turunkan` | — |
| 9 | Activity | `CountOGPONP_Act` | ✅ | `backend/models/hitung.go` `CountOGPONP` - aksi admin tab New Data (`backend/services/aksiposisi.go`) dan `validasiKirim` | — |
| 10 | Activity | `CountOverridingCommOgp_Act` | ⛔ | Pemicunya hanya sel PropNewData / PropOldData / PropOldData2 / PropValueDifference yang seluruhnya hanya-baca; fungsi salinan NB `backend/models/hitung.go` `CountOverridingCommOgp` tidak didaftarkan `aksiHitung` | tak terjangkau (sel hanya-baca) |
| 11 | Activity | `CountOverridingCommOnp_Act` | ⛔ | Idem no. 10 (`CountOverridingCommOnp` tanpa pemanggil) | tak terjangkau (sel hanya-baca) |
| 12 | Activity | `CountPctInstallment_Act` | ⛔ | Grid `.ListInstallment` ber-edit mode readOnly di semua tab; `CountPctInstallment` (`backend/models/angsuran.go`) salinan NB tanpa pemanggil | tak terjangkau (grid hanya-baca) |
| 13 | Activity | `CountResult1Onp_Act` | ✅ | `backend/models/hitung.go` `CountResult1Onp` (aksi admin; CountOGPONP_Act 5) | — |
| 14 | Activity | `CountResult1_Act` | ✅ | `backend/models/hitung.go` `CountResult1` (aksi admin; CountOGPONP_Act 3) | — |
| 15 | Activity | `CountResult2Ogp_act` | ✅ | `backend/models/hitung.go` `CountResult2Ogp` (aksi admin; CountOGPONP_Act 4) | — |
| 16 | Activity | `CountResult2Onp_act` | ✅ | `backend/models/hitung.go` `CountResult2Onp` (aksi admin; CountOGPONP_Act 6) | — |
| 17 | Activity | `CountRiCommOgp_act` | ⛔ | Idem no. 10 (`CountRiCommOgp` tanpa pemanggil) | tak terjangkau (sel hanya-baca) |
| 18 | Activity | `CountRiCommOnp_act` | ⛔ | Idem no. 10 (`CountRiCommOnp` tanpa pemanggil) | tak terjangkau (sel hanya-baca) |
| 19 | Activity | `CountSpreading_Act` | ✅ | Versi EDM (lebih baru dari NB): `backend/models/angsuran.go` `CountSpreading`, `HitungTotalSpreading` - langkah 4 bawaan 100, 5.1 `SplitRNMSharePct / RNMShare` presisi 20; langkah 1-3, 7 `//` | `//` 1-3, 7 |
| 20 | Activity | `CreateEDMT` | ✅ | `backend/services/layanan.go` `BuatKasus` + `backend/models/edm_buat.go` `RakitHalamanBaru`, `PasangOldData`; generasi lama `backend/repository/generasi.go` `GenerasiTerakhir` (pengganti FetchPolisJsonPolis); ditolak bila json_polis memuat generasi lebih banyak dari rantai relasional (`CacahGenerasiJSONPolis`, pesan bukan XML `PesanGenerasiBelumDimuat`). Langkah 16 (SetEDMTCancel atas halaman portal) [dugaan] tanpa efek - pembatalan berlaku lewat EDMChooseBusiness_Act 5 | `//` 4-6 |
| 21 | Activity | `EDMChooseBusiness_Act` | ✅ | `backend/models/edm_pilih.go` `PilihBisnisEDM` + `backend/services/tindakan.go` `PilihBisnis` (langkah 13-14 Obj-Save di transaksi). Langkah 8 DT ExpandAllExpandables = keadaan layar | keadaan layar (8) |
| 22 | Activity | `EDMTCalculateTreatyDifference` | 🟡 | `backend/models/edm_selisih.go` `EDMTCalculateTreatyDifference` langkah 1-3 (baru - lama, pasangan baris per posisi): aksi tombol Calculate Value Difference DAN dihitung ulang setiap generasi ditulis (`backend/models/edm_bentuk.go` `HitungSelisihGenerasi` dari `simpanGenerasi` dan sebelum produksi Utility1). Langkah 4-6 (label HasEDMNo, pengurang `OldData.TreatyDifference`) TIDAK - satu rumus untuk semua generasi | keputusan WO 23-09-2026 (ID-28/30) |
| 23 | Activity | `FetchTreatyGroupOJK` | ⛔ | Dipanggil hanya GeneratePolicyNoTreaty_Act 13 (no. 30, tak terjangkau); kelas Work-NB bukan leluhur kasus EDM. `SetelOJK` (`backend/models/pilihbisnis.go`) salinan NB tanpa pemanggil | tak terjangkau |
| 24 | Activity | `FetchTreatyGroupOldID` | ⛔ | Dipanggil hanya GeneratePolicyNoTreaty_Act 11 (no. 30) | tak terjangkau |
| 25 | Activity | `FillMasterInstallment` | ✅ | `backend/models/edm_pilih.go` `FillMasterInstallment` (EDMChooseBusiness_Act 7) | — |
| 26 | Activity | `FillPaymentInstallment` | ✅ | `backend/models/angsuran.go` `FillPaymentInstallment` - sel `.Installment` PropNewData2 (aksi admin) | — |
| 27 | Activity | `FillPaymentInstallmentEDMT` | ✅ | `backend/models/edm_angsuran.go` `FillPaymentInstallmentEDMT` (EDMChooseBusiness_Act 12; sel `.Installment` S12 NonProp baru). Langkah 3 pyExpanded = keadaan layar | keadaan layar (3) |
| 28 | Activity | `FillSpreading` | ✅ | `backend/models/edm_pilih.go` `FillSpreading` (EDMChooseBusiness_Act 9; pesan VERBATIM `PesanSpreadingAsalKosong`) | — |
| 29 | Activity | `GeneratePolicyNoTreatyAddendum_Act` | ⛔ | Tombol Submit S18 menjalankannya hanya bila `PolicyNo = ""` dan operator = ID tertanam; PolicyNo terisi sejak SetEDMTNoPolis 3, jadi tidak pernah jalan. Dept Head setuju = popup ShowPolicyNoTreaty (`backend/models/layar.go` `TombolUntuk`) | tak terjangkau + penyimpangan sadar (ID operator diganti posisi); `[terbuka]` S 9.2 #3 (#45) |
| 30 | Activity | `GeneratePolicyNoTreaty_Act` | ⛔ | Tombol C S18: hanya bila `PolicyNo = ""` dan PositionNote Sec Head - tidak pernah benar. Nomor polis generasi = nomor polis induk (`backend/repository/generasi.go` `SetelNomorPolisSelesai`) | tak terjangkau |
| 31 | Activity | `GetLinkService` | ⛔ | Konversi Arasapas tidak disambung (sama NB); `backend/models/konversi.go` hanya menyusun muatan | penyimpangan sadar (sambungan menunggu persetujuan WO) |
| 32 | Activity | `InputPolicyTreatyEDMDetail_NP` | ✅ | `backend/models/edm_pilih.go` `InputPolicyTreatyEDMDetailNP` + `pbInputDetailNP` (SetValueEDM_Act 9; dijalankan juga untuk Prop - ditiru) | `//` 1, 5.4, 7 |
| 33 | Activity | `InputPolicyTreatyEDMDetail_NP_AdjPremi` | ✅ | `backend/models/edm_pilih.go` `InputPolicyTreatyEDMDetailNPAdjPremi` (SetValueEDM_Act 10, EDMType 3) | `//` 1, 5.4, 7 |
| 34 | Activity | `InputPolicyTreatyInPre_Act` | ✅ | `backend/services/layanan.go` `siapkan` (langkah 2 `TerapkanBisnisPra`, 3-4 dan 9 `PraprosesTanggal`, 5-8 acuan spreading) + `backend/services/nonprop.go` `siapkanNonProp` (10). Langkah 1 SetCategoryAttach (lampiran FacIn) tidak | tanpa pembaca (langkah 1) |
| 35 | Activity | `InsertHistoryAkseptasiPega` | ✅ | `backend/repository/riwayat.go` `CatatRiwayat` dari `backend/services/tindakan.go` `kirim` (pasca-activity kedua flow action; OPERATORID = akses login) | — |
| 36 | Activity | `InsertToTreatyOutXOLList` | ✅ | `backend/models/edm_xol_keluar.go` `InsertToTreatyOutXOLList` (InputPolicyTreatyEDMDetail_NP 9, XOL Retro; master M_TREATY_OUT baca-saja) | `//` 2 langkah |
| 37 | Activity | `InsertToTreatyOutXOLListEDMOldData` | ✅ | `backend/models/edm_xol_keluar.go` `InsertToTreatyOutXOLListEDMOldData` (TreatyRealizationCheckXOLListEDM 7) | — |
| 38 | Activity | `InsertToTreatyXOLList` | ✅ | `backend/models/nonprop.go` `InsertToTreatyXOLList` (InputPolicyTreatyEDMDetail_NP 8; TreatyRealizationCheckXOLList 6 lewat `LengkapiDaftarXOL`) | — |
| 39 | Activity | `InsertToTreatyXOLListEDM` | ✅ | `backend/models/edm_xol.go` `InsertToTreatyXOLListEDM` (_NP_AdjPremi 8) | — |
| 40 | Activity | `InsertToTreatyXOLListEDMOldData` | ✅ | `backend/models/edm_xol.go` `InsertToTreatyXOLListEDMOldData` (TreatyRealizationCheckXOLListEDM 6) | — |
| 41 | Activity | `InsetTreatyInProdAddendum_Act` | ✅ | `backend/models/produksi.go` `BarisProduksi`, `NotaJenisReasXOL` + `backend/repository/produksi.go` `SimpanPolisProduksi` (penjaga IDPEGA langkah 6-8; nilai = SELISIH). Langkah 3 (GetTglInputCaseId) nol pembaca | tanpa pembaca (3) |
| 42 | Activity | `ProtectionNonProp_Act` | ✅ | `backend/models/edm_proteksi.go` `protectionNonProp` (Protection_Act 15) | — |
| 43 | Activity | `Protection_Act` | ✅ | `backend/models/edm_proteksi.go` `ProtectionAct` - aksi `Protection` radio Approval dan `validasiKirim`. Kalang 8-10 `PolicyTreatyInDetail` nol putaran (daftar tidak ditulis rule EDM mana pun) | `//` 6-7 |
| 44 | Activity | `RemoveTypeTax_ACT` | ✅ | `backend/models/pilihbisnis.go` `RemoveTypeTax` (aksi `RemoveTypeTax`, sel With Tax) | — |
| 45 | Activity | `SaveJsonPolisTreatyInEDM_Act` | 🟡 | Utility1 di transaksi Dept Head menyetujui: `backend/models/produksi.go` `PrasimpanPolis`, `PrasimpanMedan`, `SusunSimpananPolis`, `JSONPolis` + `backend/repository/produksi.go` `SimpanPolisProduksi` (json_polis, ACHIEVEMENT, TREATYINPRODUCTION). Langkah 10 (DATA_JSON) TIDAK; 1 dan 5 tidak (pola NB) | keputusan WO 06-10-2026 (tanpa JSON); `//` 9, 16-18 |
| 46 | Activity | `SendEmailWithAttachments` | ⛔ | Konversi langkah 19-20 (gagal dan IsPEGAPROD); rule bawaan Pega, alamat tertanam | keputusan WO 07-10-2026: konversi tetap tidak disambung (konversi 11-20; #50) |
| 47 | Activity | `SetCategoryAttach` | 🟡 | Versi Work (pra-activity DeptHeadTreatyIn_UWAddendum): langkah 2-4 daftar spreading = `backend/repository/acuan.go` `DaftarJenisSpreading` lewat `DaftarAcuan`. Langkah 1 (Call Data-OfferFacIn.SetCategoryAttach, lampiran) tidak | tanpa pembaca (lampiran tidak ada di layar EDM) |
| 48 | Activity | `SetCurrency_act` | ⛔ | Pemicunya sel Currency header ber-readOnly + disabled - tak pernah terpicu; `SetelIDMataUang` salinan NB tanpa pemanggil | tak terjangkau (sel disabled) |
| 49 | Activity | `SetDueTo_act` | ✅ | `backend/models/hitung.go` `SetDueTo` (aksi `SetDueTo` radio Approval; `kirim`; CountNetPremi_act 7) | — |
| 50 | Activity | `SetEDMAchivementValue` | ✅ | `backend/models/produksi.go` `BarisCapaian` + `sqlSisipCapaian` (`backend/repository/produksi.go`); nilai = selisih | — |
| 51 | Activity | `SetEDMTCancel` | ✅ | `backend/models/edm_pilih.go` `SetEDMTCancel` (EDMChooseBusiness_Act 5, EDMType 4) - 16 medan + angsuran + spreading seperti XML | cakupan nol = butir WO #47 |
| 52 | Activity | `SetEDMTNoPolis` | ✅ | `backend/models/edm_buat.go` `NomorEDM` dari `RakitHalamanBaru`; NOENDORS `backend/repository/kasus.go` `sqlSisipGenerasi` | — |
| 53 | Activity | `SetPPNPPH` | ✅ | `backend/models/hitung.go` `SetPPNPPH` + `backend/repository/acuan.go` `StsPKPAgen` | — |
| 54 | Activity | `SetReinstatementPct` | ✅ | `backend/models/nonprop_detail.go` `SetReinstatementPct` (lewat TreatySetReinstatement; halaman master, tidak disimpan) | — |
| 55 | Activity | `SetTreatyInEDM_Act` | ✅ | `FillMasterInstallment` langkah 3 (`MasterEDMMenurutID`, `MasterOutMenurutID`). Langkah 1 TreatyInInputVis nol efek; 6 tidak jalan (Param.viewstate tidak ada) | tanpa pembaca (1, 6) |
| 56 | Activity | `SetTreatyIn_Act` | ✅ | `backend/models/edm_pilih.go` `pbSetTreatyIn` (langkah 3-5) dan `siapkanNonProp` + `LengkapiDaftarXOL` (13). Langkah 1-2, 6-12, 14-15 nol pembaca / revisionstate tidak pernah 1 | tanpa pembaca |
| 57 | Activity | `SetValidateInstallment_Act` | ✅ | `backend/models/angsuran.go` `SetValidateInstallment` (CountOGPONP_Act 10) | — |
| 58 | Activity | `SetValueEDM_Act` | ✅ | `backend/models/edm_pilih.go` `SetValueEDM` + `AdopsiMasterEDM` (EDMChooseBusiness_Act 1; RDB terakhir menang ditiru) | — |
| 59 | Activity | `SetValueOldTax` | ✅ | `backend/models/edm_pilih.go` `SetValueOldTax` (EDMChooseBusiness_Act 2, FlagPPH) | — |
| 60 | Activity | `TreatyInInputVis` | ⛔ | SetTreatyInEDM_Act 1 / SetTreatyIn_Act 2 (Add=1): `TreatyIn.ID = "UnknownId"` lalu ditimpa langkah sesudahnya - nol efek | tanpa pembaca |
| 61 | Activity | `TreatyLoadMasterJoinEdmChooseBusiness` | ✅ | `backend/services/tindakan.go` `DaftarBisnis` -> `backend/repository/bisnis_edm.go` `DaftarBisnisEDM` (pre-load grid popup S6) | — |
| 62 | Activity | `TreatyRealizationCheckXOLList` | 🟡 | `backend/models/nonprop_detail.go` `PerluCekDaftarXOL`, `AwalCekDaftarXOL`, `LengkapiDaftarXOL` lewat `siapkanNonProp` (InputPolicyTreatyInPre_Act 10). Langkah 7 (`OldData.TreatyXOLList = TreatyXOLList`) TIDAK diport - salinan NB yang menganggap OldData kosong | diragukan - bab 3 butir 1 |
| 63 | Activity | `TreatyRealizationCheckXOLListEDM` | ✅ | `backend/models/edm_pilih.go` `TreatyRealizationCheckXOLListEDM` (EDMChooseBusiness_Act 3) | — |
| 64 | Activity | `TreatySetReinstatement` | ✅ | `backend/models/nonprop_detail.go` `TreatySetReinstatement` (SetTreatyIn_Act 13 di jalur `LengkapiDaftarXOL`; di jalur `pbSetTreatyIn` nol pembaca) | — |
| 65 | Activity | `TrtEdmCheckPolicyError` | ✅ | `backend/services/layanan.go` `PeriksaPolis` + `backend/repository/generasi.go` `AdaEDMBerjalan` + `NoMasterDariNoPolis`; pesan VERBATIM `PesanEDMBelumSelesai` | — |
| 66 | Activity | `serviceInsertArasapas_act` | 🟡 | Utility2 sesudah transaksi: `backend/services/konversi.go` `konversikan` + `backend/models/konversi.go` `RakitMuatanKonversi` (langkah 3-6 muatan, 8 `PesanGagalKonversi`); pengirim bawaan gagal terang. Langkah 7 (FacOut) dan 11-20 tidak | penyimpangan sadar (sama NB) + keputusan WO 07-10-2026: konversi tetap tidak disambung (#49, #50); `//` 10, 18 |
| 67 | RDBList | `BrowseTreatyIn` | ✅ | `backend/repository/master_edm.go` `MasterMenurutID` + `backend/repository/masterxol.go` `MasterXOLDariJSON` (JSONDATA baca-saja, daftar medan tertutup, K8) | — |
| 68 | RDBList | `BrowseTreatyInEDM` | ✅ | `backend/repository/master_edm.go` `MasterMenurutOldID` (SetValueEDM_Act 5) | — |
| 69 | RDBList | `BrowseTreatyInEDM_Int_treaty_in_edm` | ✅ | `backend/repository/master_edm.go` `MasterEDMMenurutID` (SetTreatyInEDM_Act 3) | — |
| 70 | RDBList | `BrowseTreatyOutDetailEDM` | ✅ | `backend/repository/master_edm.go` `MasterOutMenurutID` (XOL Retro) | — |
| 71 | RDBList | `CekSTSKonversiJson` | ⛔ | Konversi langkah 11 (STS_KONVERSI) | keputusan WO 07-10-2026: konversi tetap tidak disambung (konversi 11-20) |
| 72 | RDBList | `DeleteDataProduction` | ⛔ | Konversi langkah 12 (gagal dan IsPEGAPROD; prosedur PEGA_DELETE_ERROR_KONVERSI) | keputusan WO 07-10-2026: konversi tetap tidak disambung (#50) |
| 73 | RDBList | `FetchNoOfferFromNoPolis` | ✅ | `backend/repository/generasi.go` `NoMasterDariNoPolis` | — |
| 74 | RDBList | `FetchNopolisCount` | ✅ | `backend/repository/generasi.go` `CacahGenerasiJSONPolis` (COUNT DISTINCT PRODKE, DATA_JSON tidak dibaca; juga penjaga Create) | — |
| 75 | RDBList | `FetchPolisJsonPolis` | ✅ | Diganti `backend/repository/generasi.go` `GenerasiTerakhir` + `BacaGenerasi` (`backend/repository/polis.go`): OldData dari T_GENERAL_POLIS_TREATY, bukan DATA_JSON | keputusan WO (tanpa JSON) |
| 76 | RDBList | `FetchTreatyGroupOLDID` | ⛔ | Dipanggil hanya FetchTreatyGroupOldID (no. 24) | tak terjangkau |
| 77 | RDBList | `GETTanggalClosing_SQL` | ✅ | `backend/repository/riwayat.go` `HariClosing` (penomor inti; bawaan 25) - Utility1 2-3 dan pra-proses | — |
| 78 | RDBList | `GenerateNoEDMTreaty` | ⛔ | GeneratePolicyNoTreatyAddendum_Act 10-12 berpenjaga `EDMNo == ""`; EDMNo selalu terisi - satu jalur nomor (S AC 16) | tak terjangkau (#45) |
| 79 | RDBList | `GetCountClaim` | ⛔ | CheckDuplicateOffer 6 (no. 4) | `//` + tanpa pembaca |
| 80 | RDBList | `GetCurrency` | ⛔ | SetCurrency_act 2 (no. 48) | tak terjangkau (sel disabled) |
| 81 | RDBList | `GetCurrentDate` | ⛔ | SetTreatyIn_Act 8, hanya bila revisionstate == 1 - tidak pernah dari rantai EDM | tak terjangkau |
| 82 | RDBList | `GetDataCurrencyByName_SQL` | ✅ | `IDMataUang` (`backend/repository/master_edm.go`) dan `backend/repository/acuan.go` `IDMataUangDariNama` (`idMataUangMaster`) | — |
| 83 | RDBList | `GetDataTreatyInProd_SQL` | ✅ | `backend/repository/produksi.go` `prBaca` (cacah IDPEGA TREATYINPRODUCTION; penjaga InsetTreatyInProdAddendum_Act 6-8) | — |
| 84 | RDBList | `GetKodeProdNonLife_SQL` | ⛔ | Dipanggil hanya GeneratePolicyNoTreaty_Act (no. 30) | tak terjangkau |
| 85 | RDBList | `GetOldIDBusiness_SQL` | ✅ | `backend/repository/acuan.go` `BisnisDariKunci` (InputPolicyTreatyInPre_Act 2) | — |
| 86 | RDBList | `GetPolicyNoByCaseId` | ✅ | `backend/repository/produksi.go` `prBaca` (`sqlNoPolisJSON`, SetEDMAchivementValue 2). Konversi langkah 3: PolicyNo halaman (`RakitMuatanKonversi`) | — |
| 87 | RDBList | `GetReinstypeIDbyName_SQL` | ✅ | `backend/repository/produksi.go` `prBaca` (`prSQLIDJenisReas`) + `TerapkanJenisReasXOL` | — |
| 88 | RDBList | `GetSQLDate` | ✅ | Jam aplikasi `Layanan.jam` (`backend/services/layanan.go`): ProductionDate `ProtectionAct` / `protectionNonProp`, `PraprosesTanggal` | — |
| 89 | RDBList | `GetSequenceNumber_SQL` | ⛔ | Dipanggil hanya GeneratePolicyNoTreaty_Act (no. 30) | tak terjangkau |
| 90 | RDBList | `GetTglInputCaseId` | ⛔ | InsetTreatyInProdAddendum_Act 3 -> `InputTreaty.CARI2`, tidak dipakai INSERT | tanpa pembaca |
| 91 | RDBList | `INSERTJSON_JSONPOLISMONITORING_FACIN` | ⛔ | Konversi langkah 15-16 (JSON_POLIS_MONITORING) | keputusan WO 07-10-2026: konversi tetap tidak disambung (konversi 11-20) |
| 92 | RDBList | `InsertHistoryAkseptasiPega_Sql` | ✅ | `backend/repository/riwayat.go` `CatatRiwayat` (INSERT langsung) | — |
| 93 | RDBList | `InsertTreatyInProdEDMT_SQL` | ✅ | `sqlSisipProduksi` (`backend/repository/produksi.go`) + `backend/models/produksi.go` `KolomProduksi` (58 kolom, NOENDORS = EDMNo) | — |
| 94 | RDBList | `SaveAchievementSQL` | ✅ | `sqlSisipCapaian` (`backend/repository/produksi.go`): prosedur InsertUpdateAchievment = satu INSERT, nol prosedur | — |
| 95 | RDBList | `SavePolisTreatyInEDM_SQL` | 🟡 | `sqlSisipJSONPolis` (`backend/repository/produksi.go`): INSERT json_polis TANPA DATA_JSON (pengganti prosedur PEGA_JSON_POLIS_TREATYIN) | keputusan WO 06-10-2026 |
| 96 | RDBList | `SaveTreatyIn` | ⛔ | SetTreatyIn_Act 7-11 hanya bila revisionstate == 1 - tidak pernah dari rantai EDM; prosedur PEGA_TREATY_IN menulis master | tak terjangkau |
| 97 | RDBList | `SelectProdKe` | ✅ | Diganti `backend/repository/generasi.go` `GenerasiTerakhir`: PRODKE bilangan dari T_GENERAL_POLIS_TREATY; urut teks dan potongan 24 karakter tidak ditiru | penyimpangan sadar (bug lama) |
| 98 | RDBList | `SelectSpreadingTreatyInProduction` | ✅ | `backend/repository/acuan.go` `DaftarJenisSpreading` | — |
| 99 | RDBList | `TreatyInSearchProdKe` | ✅ | `prSQLCacahNoPolis` (`backend/repository/produksi.go`): PRODKE json_polis = COUNT per NOPOLIS (#56) | — |
| 100 | RDBList | `TreatyLoadMasterJoinEdmChooseBusiness` | ✅ | `backend/repository/bisnis_edm.go` `DaftarBisnisEDM` (TREATY_IN ∪ TREATY_IN_EDM, awalan 7 karakter `OldData.NoOffer`) | — |
| 101 | RDBList | `TreatyLoadMasterJoinEdmChooseBusinessRetro` | ✅ | `backend/repository/bisnis_edm.go` `DaftarBisnisEDM` (TREATY_OUT2, `pilihTanpaOldID`) | — |
| 102 | RDBList | `UpdateErrorNoteJsonPolisMonitoring` | ⛔ | Konversi langkah gagal (ERR_NOTE JSON_POLIS_MONITORING) | keputusan WO 07-10-2026: konversi tetap tidak disambung (konversi 11-20) |
| 103 | Section | `BusinessAndSOBListEDM` | ✅ | `frontend/components/PilihBisnis.tsx`, `frontend/gridbisnis.ts`; isi `DaftarBisnisEDM`. Baris yang diklik tidak dibaca rantai Choose (XML) - ditiru | — |
| 104 | Section | `DetailPolicyAddPremiDetail` | ✅ | `frontend/xol.ts` `KOLOM_XOL`, `ANAK_XOL` + `frontend/components/Grid.tsx` | — |
| 105 | Section | `DetailPolicyTreatyInAddGeneral` | ✅ | `frontend/medan.ts` `tampilTabData`, `nonPropBaru` + `frontend/pages/LayarKasus.tsx` (cabang AddPremi lawan tab) | — |
| 106 | Section | `DetailPolicyTreatyInAddGeneralEditable` | ✅ | `frontend/components/TabData.tsx` + strip tab `frontend/pages/LayarKasus.tsx`; `varianTab` | — |
| 107 | Section | `DetailPolicyTreatyInAddPremi` | ✅ | `frontend/components/AddPremi.tsx`, `frontend/xol.ts` | — |
| 108 | Section | `DetailPolicyTreatyInAddendum` | ✅ | Header `frontend/medan.ts` `MEDAN_KIRI`, `MEDAN_KANAN`, `medanRemark` + `frontend/pages/LayarKasus.tsx`; aturan `backend/models/layar.go` `MedanWajibKosong`, `GabungMasukanLayar`, `TombolUntuk`. Syarat tombol S18 ber-ID operator diganti posisi | penyimpangan sadar (ID operator); sel `1=2` / NEVER tidak dirender |
| 109 | Section | `DetailPolicyTreatyInPropNewData` | ✅ | `frontend/medan.ts` `tataTab` varian `baru` + `TabData.tsx` (atasan, hanya-baca) | — |
| 110 | Section | `DetailPolicyTreatyInPropNewData2` | ✅ | `tataTab` varian `baruAdmin` (Save, Calculate Value Difference, Add tanpa Delete); aksi `backend/services/aksiposisi.go` | keputusan WO (tanpa Delete, ID-16) |
| 111 | Section | `DetailPolicyTreatyInPropOldData` | ✅ | `tataTab` varian `lama` (awalan `OldData.`) | — |
| 112 | Section | `DetailPolicyTreatyInPropOldData2` | ✅ | `tataTab` varian `lama2` (awalan `OldData.TreatyDifference.`) | — |
| 113 | Section | `DetailPolicyTreatyInPropValueDifference` | ✅ | `tataTab` varian `selisih`; defer-load = aksi `EDMTCalculateTreatyDifference` | — |
| 114 | Section | `GeneralPolicyTreatyInAddendum` | ✅ | `frontend/pages/LayarKasus.tsx` (satu section semua posisi) | — |
| 115 | Section | `InstallmentList` | ✅ | `frontend/xol.ts` (`ANAK_ANGSURAN`, `KOLOM_RINCI_ANGSURAN`) + `Grid.tsx` expand pane | — |
| 116 | Section | `ListSuggestEDM` | ✅ | `frontend/components/Usulan.tsx`; wajib Approval/Suggest `backend/models/layar.go`; action set tanpa CekLimitTreatyAcc_Act (no. 2) | penyimpangan sadar (no. 2) |
| 117 | Section | `PolicyTreatyInDeclineConfirm` | ✅ | `frontend/pages/LayarKasus.tsx` (popup Yes kirim / No tutup) | — |
| 118 | Section | `SFAPortal_Endorsement_Treaty` | ✅ | `frontend/pages/PortalEDMTreatyIn.tsx` + `backend/repository/kasus.go` `sqlDaftarKasus` (tautan hanya kasus posisi Admin, seperti XML) | — |
| 119 | Section | `ShowPolicyNoTreaty_SC` | ✅ | `frontend/pages/LayarKasus.tsx` (Modal tanpa tutup, satu aksi OK) | — |
| 120 | Section | `TreatyCreateEdm` | ✅ | `frontend/components/BuatEDM.tsx` + `backend/services/layanan.go` `PeriksaPolis`, `BuatKasus` | — |
| 121 | When | `IsClaim` | ⛔ | `pyWorkIDPrefix = "CLM-"`; kasus EDM selalu `EDMT-` - tidak pernah benar, cabang benarnya tidak dibangun | tak terjangkau |
| 122 | When | `IsFacRetro` | ⛔ | Satu-satunya evaluasi: konversi langkah 7 (pengiriman FacOut) - tidak dibangun. Properti `OfferFacIn.IsFacRetro` tetap diisi `pbInputDetailNP` langkah 3 | keputusan WO P50 (S AC 36); butir WO #49 |
| 123 | When | `IsPEGAPROD` | ⛔ | Konversi 12 dan 20 (6 tak dicentang) tidak dibangun; konversi sistem baru hanya di produksi (`Lingkungan().AdalahProduksi()`, `backend/services/gudang.go` `DariDasar`) | keputusan WO 07-10-2026: konversi tetap tidak disambung (#50) |
| 124 | When | `IsSPVCreate` | ⛔ | Membandingkan pembuat kasus dengan dua ID operator tertanam; kedua cabang menunggu workbasket yang sama - diganti posisi | penyimpangan sadar (S AC 6; NB K5) |
| 125 | When | `IsSPVTreaty1` | ⛔ | Peran dibaca dari kolom telepon operator; diganti posisi | penyimpangan sadar (S AC 6; NB K5) |
| 126 | When | `IsSuccessHitService` | ⛔ | Status respons konversi; sambungan konversi tidak dibangun | penyimpangan sadar (sama NB) |
| 127 | When | `IsTreaty1` | ⛔ | Idem no. 125 | penyimpangan sadar (S AC 6; NB K5) |
| 128 | When | `IsTreatyIn` | ⛔ | Konversi 8-9 (pesan gagal dilewati untuk treaty); sistem baru selalu mengembalikan `PesanGagalKonversi` (niat P52, S AC 42) | keputusan WO 07-10-2026: konversi tetap tidak disambung (#50) |
| 129 | When | `IsUW` | ⛔ | Pembanding nama workbasket FacIn tanpa kutip ([dugaan] selalu salah) pada sel ber-wadah `1=2`; wewenang = keanggotaan antrean | tak terjangkau + penyimpangan sadar (pola NB) |
| 130 | When | `ToTREATYDEPTHEAD` | ⛔ | Decision9 (`LetterNo == "TREATYINDEPTHEAD"`): Sec Head menyetujui SELALU naik ke Dept Head (`backend/models/tangga.go` `Langkah`) | penyimpangan sadar (prompt WO; koreksi #42) |
| 131 | When | `recordEvent` | ⛔ | Rule standar Pega (timer klien harness TreatyCreateEdm) - bukan aturan dagang | tanpa pembaca |
| 132 | DataTransform | `AddToListCommentsPolicyTreatyIn_DT` | ✅ | `backend/models/layar.go` `TambahCatatan`; baris baru disimpan `CatatUsulan` ke HISTORYAKSEPTASIPRODUCTION TYPE_POLIS 'EDMT' | penyimpangan sadar (pola NB K4; #53) |
| 133 | DataTransform | `DeptHeadTreatyInAddendum_PreDT` | ✅ | `backend/models/layar.go` `PraprosesAtasan` (ViewState keadaan layar; OperatorName = nama tampilan) | ketetapan NB P33 |
| 134 | DataTransform | `ExpandAllExpandables` | 🟡 | pyExpanded tidak disimpan (keadaan layar); padanan tampilan `frontend/components/Grid.tsx` - rincian terbuka sejak awal | keadaan layar (SP AC 34) |
| 135 | DataTransform | `InboxPolicyTreatyInAddendum_postDT` | ✅ | `backend/models/layar.go` `PascaAdmin`; langkah 3 (nama orang tertanam) diganti data `TeksNBStatusKotakMasuk` sesudah connector | penyimpangan sadar (pola NB P40) |
| 136 | DataTransform | `InboxPolicyTreatyIn_UW_postDT` | ✅ | `backend/models/layar.go` `PascaAtasan` (langkah 1); 2-5 hanya untuk posisi yang tidak ada di flow EDM | tak terjangkau (2-5) |
| 137 | DataTransform | `InputPolicyTreatyInAddendum_preAddDT` | ✅ | `backend/models/layar.go` `PraprosesAdmin` | — |
| 138 | DataTransform | `InputPolicyTreatyIn_preAddDT` | ✅ | `backend/models/layar.go` `PraprosesAdmin` langkah 1-12; 13 (`TempEmail.CARI28` = urutan workbasket) diganti keanggotaan antrean | penyimpangan sadar (S AC 51) |
| 139 | DataTransform | `SetInstallmentValue` | ✅ | `backend/models/edm_pilih.go` `SetInstallmentValue` (FillMasterInstallment 4) | — |
| 140 | DataTransform | `SetOldMasterNo` | ✅ | `backend/repository/bisnis_edm.go` `awalanKarakter` (7 karakter `OldData.NoOffer`, saringan A RD BrowseTREATY_IN_EDM) | — |
| 141 | DataTransform | `SystemSetOneYear_DT` | ⛔ | Pemicunya sel Statement Period `.StartDate` ber-disabled - tak pernah terpicu; `SystemSetOneYear` (`backend/models/polisaturan.go`) salinan NB tanpa pemanggil | tak terjangkau (sel disabled) |
| 142 | ReportDefinition | `BrowseClientName_RD` | ✅ | `backend/repository/acuan.go` `StsPKPAgen` (SetPPNPPH 3; `pasangStsPKP`) | — |
| 143 | ReportDefinition | `BrowseCurrencyTreatyIn_RD` | ✅ | `backend/repository/acuan.go` `DaftarMataUang` (acuan `mataUang`, dropdown Currency header) | — |
| 144 | ReportDefinition | `BrowseMarketingOfficer_RD` | ✅ | `backend/repository/acuan.go` `DaftarMO` | — |
| 145 | ReportDefinition | `BrowseReinsuranceType_RD` | ✅ | `backend/repository/acuan.go` `DaftarJenisReas` (Type Treaty grid spreading) | — |
| 146 | ReportDefinition | `BrowseTREATY_IN` | ⛔ | CheckDuplicateOffer 1-2 (no. 4) | `//` |
| 147 | ReportDefinition | `BrowseTREATY_IN_EDM` | ✅ | `backend/repository/bisnis_edm.go` `DaftarBisnisEDM` S1 (EDMType 3: A Contains, J StatusAkseptasi, L EDMState 3) | — |
| 148 | ReportDefinition | `BrowseTreatyGroup_RD` | ⛔ | FetchTreatyGroupOJK 3 (no. 23) | tak terjangkau |
| 149 | ReportDefinition | `GetListEdmTreaty` | ✅ | `backend/repository/generasi.go` `AdaEDMBerjalan` | — |
| 150 | ReportDefinition | `InboxEDM_RD2` | 🟡 | `backend/repository/kasus.go` `sqlDaftarKasus` (A pembuat, B Contains, C status, urut update, 500). Filter D `Quotation.TeamGroup` (pengisi parameter tidak ada di korpus) dan E (tidak dikirim grid) tidak dibangun | tanpa pengisi (butir terbuka) |
| 151 | FlowAction | `DeptHeadTreatyIn_UWAddendum` | ✅ | `backend/services/layanan.go` `siapkan` (atasan) + `backend/services/tindakan.go` `kirim`; layar `frontend/pages/LayarKasus.tsx` | — |
| 152 | FlowAction | `DetailPolicyAddPremiDetail` | ✅ | Expand pane `frontend/components/Grid.tsx` (`frontend/xol.ts`) | — |
| 153 | FlowAction | `InboxPolicyTreatyInAddendum` | ✅ | `backend/services/layanan.go` `siapkan` (admin) + `kirim`; layar `LayarKasus.tsx` | — |
| 154 | FlowAction | `InstallmentList` | ✅ | Expand pane `frontend/components/Grid.tsx` (`frontend/xol.ts`) | — |
| 155 | FlowAction | `PolicyTreatyInDeclineConfirm` | ✅ | `frontend/pages/LayarKasus.tsx`; `backend/models/layar.go` `TombolKonfirmasiTolak` | — |
| 156 | FlowAction | `ShowPolicyNoTreaty` | ✅ | `frontend/pages/LayarKasus.tsx`; `backend/models/layar.go` `TombolNomorPolis` | — |
| 157 | DecisionTable | `isApproved` | ✅ | `backend/models/tangga.go` `Disetujui` (teks `"0"` = No, selain itu YES) | — |
| 158 | Harness | `BusinessAndSOBListEDM` | ✅ | `frontend/components/PilihBisnis.tsx` (popup) | — |
| 159 | Harness | `SFAPortal_Endorsement_Treaty` | ✅ | `frontend/pages/PortalEDMTreatyIn.tsx` + `frontend/menu.ts` (satu halaman); rute `GET /kasus` | — |
| 160 | Harness | `TreatyCreateEdm` | ✅ | `frontend/components/BuatEDM.tsx` (popup) | — |
| 161 | Flow | `InputAddendumTreatyIn` | 🟡 | `backend/models/tangga.go` `Langkah` + `backend/services/tindakan.go` `kirim`. Decision9 Else (Sec Head selesai sendiri) tidak; Decision5/10/6 diganti posisi; status tutup Resolved-Rejected / Resolved-Completed | penyimpangan sadar (prompt WO; NB P24) |
| 162 | ConnectREST | `convertJsonNusareToProduction` | 🟡 | `backend/models/konversi.go` `RakitMuatanKonversi` (noPolis, caseId, tglInput); POST tidak disambung - pengirim bawaan gagal terang | penyimpangan sadar (sama NB) |
| 163 | SystemSettings | `LinkService` | ⛔ | URL dari M_LINK_SERVICE tidak dibaca; sambungan konversi tidak dibangun | penyimpangan sadar (sama NB); S AC 45 |

## 3 · Rule dan temuan yang diragukan

1. **`TreatyRealizationCheckXOLList` langkah 7 (no. 62).** Pega menimpa `OldData.TreatyXOLList` dengan daftar XOL baru
   dari master saat pra-proses Admin (`IsNewPolicyNonProp` 1, EDMType ≠ 3, daftar kosong); nilainya ikut blob kasus dan
   dibaca `TreatyRealizationCheckXOLListEDM` 3 serta `CalculateDifferenceEDM_act`. Port NB yang disalin tidak memuat
   langkah ini (NB: OldData selalu kosong), dan sistem baru membaca OldData ulang dari generasi `OLD_POLIS_ID`, sehingga
   efek langkah 7 tidak dapat bertahan. Perilaku Pega itu tidak ditiru dan **belum diputuskan** — perlu ditinjau.
2. **`CreateEDMT` langkah 16 (no. 20)** — `SetEDMTCancel` atas halaman portal sebelum kasus lahir: `[dugaan]` tanpa efek;
   pembatalan yang berlaku `EDMChooseBusiness_Act` 5.
3. **`IsUW` (no. 129)** — pembanding nama workbasket tanpa kutip; `[dugaan]` laporan pindai: selalu salah.
4. **`InputPolicyTreatyEDMDetail_NP` untuk polis proporsional (no. 32)** — `SetValueEDM_Act` 9 tidak memeriksa jenis
   proporsi (ditiru); daftar XOL hasil antara generasi proporsional dibuang `RapikanBentukSimpan` sebelum disimpan.
5. **Kode salinan NB tanpa pemanggil non-uji** di `backend/models` — `CountOverridingCommOgp`, `CountOverridingCommOnp`,
   `CountRiCommOgp`, `CountRiCommOnp`, `CountPctInstallment`, `SetelOJK`, `SetelGrupLama`, `GrupTreatyTakTerbaca`,
   `SetelIDMataUang`, `SetelNamaMataUang`, `SystemSetOneYear`, `InputDetailNonProp`, `TetapkanHasFacOut`,
   `CalculatePremi`, `RakitNomorPolis`, `TerapkanDetailKontrak`, `TerapkanMasterKontrak`: diuji, tetapi tidak terjangkau
   dari HTTP (rule pemicunya ⛔ di tabel). Bukan pelanggaran; dibuang atau dibiarkan — keputusan di luar dokumen ini.

