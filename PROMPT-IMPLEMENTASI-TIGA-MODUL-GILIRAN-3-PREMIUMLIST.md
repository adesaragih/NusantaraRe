# PROMPT — GILIRAN 3 · SESI B **PremiumList Life** *(folder `OUTPUT_HASIL_RNM\.worktrees\premiumlist-life`, cabang `modul/premiumlist-life` @ `4f40272`)*: **tiket 01 → 09 seluruhnya, dari XML, dalam satu giliran**

> Baca `PROMPT-INDUK-TIGA-MODUL.md` *(§0.4)*, `PROMPT-IMPLEMENTASI-MODUL-PREMIUMLIST-LIFE.md` *(§1–§3, pl1–pl5 tetap berlaku)*,
> tiket `.scratch/premiumlist-life/issues/`. Tiket **00 sudah tuntas dan menyatu** *(`83658a8` → `c1d3b8c`; migrasi `050`–`056`)*.
> Mekanisme giliran = lanjutan 8 §1; XML patokan; tiket = hasil grilling, diralat bertanggal bila XML membantahnya.
> Sesi ini **tidak** merge ke `main` — asisten menyatukan sesudah memverifikasi log. Bila sesi tunggal: kerjakan di `main`
> sesudah brief A.

---

## 0. KONTEKS YANG SUDAH PASTI

| Hal | Isi |
| --- | --- |
| Cabang dan port | `modul/premiumlist-life` @ `4f40272` *(fast-forward asisten 28-09-2026)*; backend `:8082`, Vite `5174`; `.env` sudah ada di worktree |
| Katalog DEV *(27-09-2026)* | **nol tabrakan nama** untuk `T_WORK_POLIS`, `T_PREMIUM_LIST`, `_DETAIL`, `_SPREADING`, `_SPREADING_RETRO`, `_SUMMARY`, `T_VIEW_SUGGEST`, `SEQ_WORK_POLIS`; `T_MIGRASI` DEV belum memuat `050`–`056` *(work owner menjalankan dari `main`)* |
| Keputusan | pl1–pl5 *(brief modul §2)*; decision **o** *(prosedur tidak dipanggil, logikanya ditiru dari `ALL_SOURCE`)*; ADR-U-0029 nol `COMMIT`; av = milik sesi A sesudah tiket 04 menyatu |
| Berkas milik sesi ini | `polis_*` di `internal/*`, `handlers/rute_premiumlist.go`, `frontend/src/pages/premiumlist/`, `labels.premiumlist.ts`; Shell hanya **ditambah** kelompok menu; berkas Claim Life/Komite **tidak** disunting *(perubahan bersama hanya aditif, dilaporkan)* |
| Migrasi | `057`+ hanya dari keputusan tercatat; `-migrate` **tidak** dijalankan executor |
| Menu `[dari bukti]` | kelompok **PremiumList Life** → satu butir **`PremiumList`** *(harness portal `PremiumLife_harness` → section `PremiumList`)*; tombol `Input Premium` b16472 dan `Input Offer` b16964 → `runActivity` `CreateInputLife` b3291/b3578/b3939/b4214. Tidak ada butir lain |

## 0.1 PAKET 0 — MIGRASI `056` GAGAL DI DEV: `INITIAL` KATA CADANGAN ORACLE `[keputusan pl6; veto work owner]`

Work owner menjalankan `-migrate` 28-09-2026 09.33: `017`–`020`, `030`, `050`–`055` **terpasang**; `056_t_view_suggest`
**gagal** dan `T_VIEW_SUGGEST` tidak ada di DEV. Sebabnya bukan `{skema}` *(sebelas langkah lain lolos dengan penanda itu)*,
melainkan kolom **`INITIAL`**: kata cadangan Oracle — `SELECT 1 AS INITIAL FROM DUAL` → `ORA-00923`; `"INITIAL"` berkutip
lolos; `NO` lolos. Tidak ada tabel warisan yang memakai nama itu *(katalog DEV: nol kolom `INITIAL`, nol tabel `%SUGGEST%`)*;
sumbernya properti `.Initial` *(`Activity/AddHistorySuggest.xml`)* dan STRUKTUR baris 393 *(spec §12 `Offer`/`Bind`)*.

**pl6**: kolom diganti nama **`INITIAL_SUGGEST`** *(mengikuti pola saudaranya `DATE_SUGGEST`, `PIC_SUGGEST`, `COMMENT_SUGGEST`)*;
berkas `056` **disunting di tempat** *(belum pernah terpasang di mana pun: `T_MIGRASI` DEV tanpa `056`, skema uji belum ada)*;
ralat bertanggal di STRUKTUR baris 393, `spec.md` §12, tiket 00 dan 09; penjaga `TestKolomDDLCocokDenganStruktur` mengikuti.
Commit `premiumlist-life: tiket 00 — ralat 056, INITIAL kata cadangan Oracle (pl6)`. Sesudah menyatu ke `main`, work owner
menjalankan `-migrate` lagi; hanya `056` yang tersisa dijalankan.

⚠️ Pohon kerja `main` saat ini memuat **suntingan work owner yang belum di-commit** pada `016_*.sql` *(`{skema}` → `POOLDATA`)*
yang membuat dua penjaga merah *(`TestSetiapPernyataanSahDanBerskema`, `TestKolomDDLCocokDenganStruktur`)*. Executor **tidak**
meng-commit dan **tidak** membuangnya — itu keputusan work owner *(pesan asisten 28-09-2026)*; bekerja di cabang/worktree
modul, atau tunggu pohon `main` bersih.

## 1. XML — TITIK MASUK *(korpus `D:\XML\RNM_BRD\PremiumList Life\`, dibaca 28-09-2026; nomor baris = `sed -e 's/></>\n</g'`)*

**Flow `InputPolicyHolder.xml`** *(di akar korpus; tidak ada folder `Flow\`)*: task `Input Premium List Detail` b306, `Input
Premium List Summary` b321, `Input Offer` b351; flow action `InputDataOfferLife` b366; utility `InsertJsonPolisLife` b261/b782,
`FlagOnGoingPolicy` b276/b900; decision `Accept` b291/b336/b996/b1180; **status kerja** `Resolved-Rejected` b848,
`Resolved-Completed` b947; `CloseCase` b960/b1470; konektor `Input Premium Detail` b1069, `Input Premium Summary` b1232,
`Input Offer Life` b1340. Peta konektor **utuh** *(pyTo mendahului pyFrom di DOM; cara Claim Life lanjutan 12 §1)* ditulis di
tiket 01 sebelum kode.

| Flow action | Section yang ditampilkannya | Activity |
| --- | --- | --- |
| `InputDataOfferLife` | `InputOfferLife` | `GetOfferLife_Act` b2172 → `FollowingOffer_Harness` b2209; `Ceding_Harness` b3632; `PolicyHolder_Harness` b4135; `SetMaxTBCLife_Act` b9103/b9245; `ProtectAccept` b11420/b11652 |
| *(konfirmasi)* | `ConfirmSection` | `ProtectAccept` b639/b868 |
| `PL_DetailAction` | `PL_Detail_Sec` *(47.283 baris — grid; tanpa tombol activity; dibaca sebagai pohon: kolom, sunting baris, syarat)* | — |
| `ShowLifePremiumDetail` | `ShowLifePremiumDetail` | `ChooseProdName` b3897; `ProtectProductName_Act` b4205; `GetPLNumber_Act` b4694; `SOB_Harness` b9225; `setSecurityReinsurer_act` b12225; `ChooseRetroName` b12624; `ChooseSecurityReinsurer` b13895; `UploadCSV_LifePremium` b19177/b19433 |
| `ShowLifePremiumSummary` | `ShowLifePremiumSummary` | `InsertJsonPolisLife_Act` b26414/b26578 + `finishAssignment` b26442/b26606; tombol `Submit` b27471 |
| `ShowPLNumber` | `ShowPolis_sc` | tombol `OK` b2751 *(closeContainer)* |
| `ChooseProdName` | `ChooseProdName` | — |
| `ChooseRetroName` | `Retro_Section` | `SearchPolicyHolder_act` b609/b732; tombol `Choose` b3820 |
| `ChooseSecurityReinsurer` | `SecurityReinsurer_Section` | — |
| `RetroLife` | `RetroDetailLife` | — |
| `UploadCSV_LifePremium` | — | `UploadCSVLifePremium_Act` *(hasil ke harness `ViewCSVResult_LifePremium` + `_QP`/`_TP`/`_TR`)* |

Harness: `PremiumLife_harness` → `PremiumList` *(portal)*; `PolicyHolder_Harness` → `PolicyHolder_Section`; `Ceding_Harness`,
`FollowingOffer_Harness`, `SOB_Harness`, `ViewCSVResult_*` = popup dari section. Activity *(38)*: `AddHistorySuggest`,
`AttachRISlipToWork`, `Calculate1_Act`, `CreateInputLife`, `GenerateDataDtlLife_act`, `GenerateDetailPeserta`, `GeneratePDFOffer`,
`GetLinkService`, `GetOfferLife_Act`, `GetPLNumber_Act`, `HTMLRISlipToPDF`, `InputOfferLife_ACT`, `InputOfferLife_preAct`,
`InputParamUploadReas_act`, `InsertJsonPolisLife_Act`, `InsertLifePremiumDetail_act`, `InsertLogServiceProd`, `LoadAttachmentData`,
`ProtectAccept`, `ProtectProductName_Act`, `SavePremiumList_Act`, `SearchPolicyHolder_act`, `SendEmailNotification`,
`SetCategoryAttach`, `SetMaxTBCLife_Act`, `SetProdNametoPolis`, `SetProductInwardLife_Act`, `SubmitPremiumList_Act`,
`UploadCSVLifePremium_Act`, `ValidasiUploadPL_act`, `WPCLife_Act`, `countCategoryAttachment_act`, `serviceInsertArasapasLife_act`,
`setCeding_act`, `setNoOffer_Act`, `setPolicyHolder_act`, `setSOB_act`, `setSecurityReinsurer_act`. Yang **menulis** ditelusuri
dari metode langkahnya *(`Obj-Save`, `RDB-Save`, `Call …_SQL`)*; `Struktur_PremiumLife_harness.xlsx` dibaca dengan
`.scratch/alat/baca-xlsx.ps1` untuk penyarangan.

## 2. URUTAN — satu commit per tiket `premiumlist-life: tiket NN — <judul>`, tanpa pesan di antaranya

| Tiket | Layar/rule XML yang menjadi pohonnya |
| --- | --- |
| **01** penawaran Confirm/Reject/Decline | `InputOfferLife` + `ConfirmSection` + `ProtectAccept`; decision `Accept`, utility `FlagOnGoingPolicy`; status `Resolved-Rejected`/`Resolved-Completed` VERBATIM; `T_WORK_POLIS.POSITION` = nama konektor VERBATIM *(050)*; inbox `PremiumList` + kedua tombol |
| **02** periode tutup buku | pl5; `SetMaxTBCLife_Act`, `WPCLife_Act`, `InputOfferLife_preAct` — gagal terang, bukan diam |
| **03** detail + nomor PL | `ShowLifePremiumDetail`, `PL_Detail_Sec`, `GetPLNumber_Act` *(penomoran ditiru; prosedur tidak dipanggil)*, `ProtectProductName_Act`, `ChooseProdName`, `SOB_Harness`, `Retro_Section`/`SecurityReinsurer_Section` |
| **04** unggah CSV | `UploadCSV_LifePremium` → `UploadCSVLifePremium_Act`, `ValidasiUploadPL_act`, `InsertLifePremiumDetail_act`, `GenerateDetailPeserta`, harness `ViewCSVResult_*`; fixture `UJI-*`; **pl4** `repository.PolisRingkas` untuk Claim Life |
| **05a** summary uang | `ShowLifePremiumSummary`, `SavePremiumList_Act`, `SubmitPremiumList_Act`, `Calculate1_Act`; **pl2** `M_LIFE_PREMIUM_SUMMARY`/`M_LIFE_PREMIUM_DETAIL` satu transaksi |
| **05b** JSON polis | `InsertJsonPolisLife_Act` + `finishAssignment` b26442 — sesuai keputusan tiket/spec *(bila spec menyatakan dibuang, tulis ralat berbukti; jangan diam)* |
| **06** efek keluar Arasapas | `serviceInsertArasapasLife_act`, `InsertLogServiceProd`, `SendEmailNotification` → outbox `T_LOG_SERVICE_RNM` + pelaksana stub *(DEV tidak memanggil endpoint nyata)* |
| **08** kontrak hilir | pl4 dikunci uji kontrak dua sisi *(pelajaran envelope `galat`)* |
| **09** migrasi skema | `050`–`056` sudah ada → verifikasi bentuk lawan STRUKTUR + penjaga; ralat tiket |

*(Tidak ada tiket 07 di folder — nomor lompat; catat, jangan mengarang tiketnya.)*

## 3. LAPORAN · TELEMETRI

Satu pesan di akhir *(atau batas tiket bila konteks menipis; pohon bersih; SHA)*: tabel **tiket → commit → menu/tombol XML →
rute/kontrol**; ralat tiket bertanggal; OQ dibuka/ditutup; kode bersama aditif yang disentuh; angka uji tiap commit; bab
**TELEMETRI EKSEKUSI** per tiket.

---

*Disusun 28 September 2026 dari `InputPolicyHolder.xml`, sepuluh flow action, sembilan harness, sembilan belas section (tombol →
activity), 38 activity, tiket 00–09, dan katalog DEV (nol tabrakan nama).*
