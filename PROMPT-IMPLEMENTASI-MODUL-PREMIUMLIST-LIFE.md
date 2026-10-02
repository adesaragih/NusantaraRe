# PROMPT — MODUL **PremiumList Life** *(sesi 2, cabang `modul/premiumlist-life`, worktree `.worktrees\premiumlist-life`)*: seluruh modul dalam satu giliran, dari XML

> Baca dulu `PROMPT-INDUK-TIGA-MODUL.md`. Aturan modul Claim Life berlaku di sini juga:
> `PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md` §1 *(per modul; XML menang, tiket diralat dengan bukti
> path + baris; pembacaan ulang XML per tiket)* dan §2 *(cara membaca XML: `sed -e 's/></>\n</g'`,
> `pyStepsPreCondParamsWhenTrue/False` **2 = lanjut, 3 = lewati**, pohon lewat PowerShell `[xml]`)*;
> mekanisme giliran lanjutan 8 §1; larangan keamanan; bab TELEMETRI EKSEKUSI.

---

## 0. KONTEKS YANG SUDAH ADA — baca sebelum menulis satu baris pun

| Sumber | Isi |
| --- | --- |
| `.scratch/premiumlist-life/spec.md` | Solution, kontrak hilir ke Claim Life *(`M_LIFE_PREMIUM_SUMMARY` + `LIFE_PREMIUM_DETAIL` berkunci `PL_NUMBER`)*, keputusan 1–7 *(batas konteks: **bukan** Endorsement Life; mesin alur; periode tutup buku dari `TANGGAL_CLOSING` **gagal terang** tanpa fallback `25`; uang teks di batas, desimal di dalam; batas transaksi **direvisi 16-09-2026: titik potong lenyap**; penomoran)* |
| `.scratch/premiumlist-life/revisi-penyimpanan-premiumlist.md` | `[keputusan work owner]` **JSON dibuang**; polis → **7 tabel baru** relasional; `T_WORK_POLIS` lintas-lini terpisah *(pola `T_WORK_CLAIM`)*; FK `ON DELETE CASCADE` + popup; uang `NUMBER(38,8)`, tanggal `DATE`, nullable, PK sequence; skala jutaan baris → index FK, pertimbangkan partisi/arsip |
| `.scratch/premiumlist-life/STRUKTUR-TABEL-PREMIUMLIST-LIFE.md` | `T_WORK_POLIS`, `T_PREMIUM_LIST`, `T_PREMIUM_LIST_DETAIL`, `T_PREMIUM_LIST_SPREADING`, `T_PREMIUM_LIST_SPREADING_RETRO`, `T_PREMIUM_LIST_SUMMARY`, `T_VIEW_SUGGEST`; di luar pohon `DOCUMENT_POLIS`, `M_TEMPUPLOADLIFE` |
| `.scratch/premiumlist-life/issues/` | `00` skema + migrasi *(PREFACTOR)* → `01` penawaran Confirm/Reject/Decline → `02` periode tutup buku → `03` PL detail + nomor PL → `04` unggah CSV → `05a` simpan summary uang → `05b` JSON polis *(periksa statusnya: JSON dibuang)* → `06` efek keluar Arasapas + alarm → `08` kontrak hilir ke Claim Life → `09` migrasi data. Status `ready-for-agent` kecuali yang `wontfix` |
| Yang sudah ada di `main` dan **dipakai ulang** | outbox `T_LOG_SERVICE_RNM` + pengirim *(aq)*, resolver `M_LINK_SERVICE` *(an/ADR-U-0013)*, `PenomorCounter` *(o1–o3, `SUMBER-PENOMORAN-DBA.md`)*, `apd.Decimal` + `desimalUang`, `PeriksaSQL`/`Qualify`, `KolomAlterTambah`, `T_MIGRASI`, Shell + `dasar.tsx` + pola `labels.ts` berbukti baris, identitas stub dari env |

## 1. XML — TITIK MASUK DAN MENU *(dibaca 27 September 2026; korpus `D:\XML\RNM_BRD\PremiumList Life\`)*

| # | Bukti | Isi |
| ---: | --- | --- |
| 1 | `Struktur_PremiumLife_harness.xlsx` *(baca dengan `.scratch/alat/baca-xlsx.ps1`)* | **root 1** Harness `PremiumLife_harness` *(class `Data-Portal`)* "(root / entry point)"; anaknya Section `PremiumList` → RD `InboxPremiumList` *(dan `InboxPremiumListOLD`, `Assign-Worklist`)* + Activity `CreateInputLife`; **root 2** Section `InputOfferLife` *(dari flow action `InputDataOfferLife`)*; **root 3** Section `ShowLifePremiumDetail` |
| 2 | `Harness\PremiumLife_harness.xml` | grid `pyGridRDName InboxPremiumList` ×48; tombol `pyButtonLabel Input Offer` dan `Input Premium` → `CreateInputLife` *(membuat kasus)*; class `Data-Portal` → **inilah harness portal modul** *(padanan `SFAPortalOpportunities` NB Treaty In)* |
| 3 | `InputPolicyHolder.xml` *(di akar folder modul, bukan `Flow\`)* | `pyWorkTypeName LIFE`; assignment `Input Offer` 351, `Input Premium List Detail` 306, `Input Premium List Summary` 321 — **semua `WorkList`, `Current operator`** *(1052, 1078, `ToCurrentOperator` 1121)*; decision `IsLifeAccepted` 1004, `IsFlagOnGoingPolicy` 909; utility `InsertJsonPolisLife` 790; task `Accept` 291, `FlagOnGoingPolicy` 276 |
| 4 | `FlowAction\` | `InputDataOfferLife`, `ShowPLNumber`, `ShowLifePremiumDetail`, `ShowLifePremiumSummary`, `PL_DetailAction`, `UploadCSV_LifePremium`, `RetroLife`, `ChooseProdName`, `ChooseRetroName`, `ChooseSecurityReinsurer` *(10)* |
| 5 | `Harness\` | popup dari `InputOfferLife`: `Ceding_Harness` *(sheet 2.10)*, `FollowingOffer_Harness` *(2.11)*, `PolicyHolder_Harness` *(2.12)*, `SOB_Harness`; tinjau CSV: `ViewCSVResult_LifePremium` + varian `_QP`, `_TP`, `_TR` |
| 6 | `ReportDefinition\InboxPremiumList.xml` *(salinannya di Claim Life sudah dibaca)* | kolom VERBATIM `Case ID`, `Create Date/Time` *(urut DESC)*, `Create Operator Name`, `Work Status`, `CedingCoName`, `PolicyHolderName`, `PL_NUMBER`, `RISLIPRNM`, `Type`, `MarketingName`, `SobName`, `DateReceived`, `Ketentuan Underwriting`; `pyPageSize 50`, `pyMaxRecords 500` |

**Menu `[dari bukti 1–2]`**: kelompok **PremiumList Life** *(nama modul korpus)* → satu butir
**`PremiumList`** *(VERBATIM nama section root, sheet r3)*, halamannya = inbox grid `InboxPremiumList`
+ dua tombol VERBATIM `Input Offer` dan `Input Premium`. Tidak ada butir lain. Layar lain dibuka dari
dalam kasus *(flow action)* atau popup. Tab/status inbox: dari `InboxType` *(baca rule-nya)*.

## 2. KEPUTUSAN YANG MENGIKAT GILIRAN INI

| Butir | Isi |
| --- | --- |
| **pl1** `[DIPUTUSKAN — konsisten dengan keputusan o Claim Life; veto work owner]` | prosedur Oracle *(`PROC_GENERATE_SEQUENCE_NUMBER`, `PEGA_M_LIFE_PREMIUM_SUMMARY`, `SaveMasterLPDet`, `INSERTJSONPOLISLIFE`)* **tidak dipanggil**; logikanya **ditiru** di Go dari sumber `ALL_SOURCE` *(dibaca sebagai katalog, teks prosedur — bukan data)* dan RDBList-nya; nol `COMMIT` di SQL *(ADR-U-0029)*; tiket 03/05a/08 yang menyebut "dipakai apa adanya"/"titik potong transaksi" diralat dengan bukti; `INSERTJSONPOLISLIFE` **tidak** ditiru sama sekali *(JSON dibuang)* |
| **pl2** | penulisan `M_LIFE_PREMIUM_SUMMARY` + `M_LIFE_PREMIUM_DETAIL` *(kontrak hilir spec; Claim Life membaca `M_LIFE_PREMIUM_DETAIL` lewat `GET /api/peserta-life`)* ditiru dari `SaveMasterLPDet`/`InsertPLSummary` **dalam transaksi Go yang sama** dengan tabel relasional *(titik potong lenyap)*; idempoten per `PL_NUMBER` |
| **pl3** | `T_WORK_POLIS` terpisah dari `T_WORK_CLAIM`; pengenal berformat dirakit di repository *(pola `PengenalWorkBerikut`)* dari `SEQ_WORK_POLIS`; awalannya **dari XML/penomoran** *(`GetPLandNopolis_sql`: `RNML-Q…`/`RNML-F…` adalah nomor PL, bukan pengenal work — jangan tertukar)* |
| **pl4** | pembaca untuk Claim Life: `repository.PolisRingkas(ctx, plNumber)` mengembalikan medan `PolicyDataLife` *(`Type`, `CedingCoName`, `PolicyHolderName`, `BusinessCode`, `BusinessName`, `MarketingName`, `SobName`, `DateReceived`, `TanggalRespon`, `TanggalKonfirmasi`, `TanggalRealisasi`, `Status`, `StatusUpdate`, `KetentuanUnderwriting`)* dari `T_PREMIUM_LIST`; kolom yang XML modul ini tidak tulis → `[tidak ada di korpus]` di tiket 08 |
| **pl5** | tutup buku: `TANGGAL_CLOSING` kosong/tak terbaca → transaksi ditolak dengan pesan yang menyebut tabelnya *(AC 8)*; geser ke tanggal 1 bulan berikutnya `05:00 GMT`; `>25` tertanam **tidak** ditiru; dua versi activity yang berbeda dicatat |

## 3. URUTAN KERJA — satu giliran, satu commit per tiket `premiumlist-life: tiket NN — <judul>`

`00` *(migrasi `050`+: tujuh tabel + `T_WORK_POLIS` + sequence + index FK; cek tabrakan nama di
katalog sebelum `-migrate` pertama)* → `01` → `02` → `03` → `04` → `05a` → *(`05b` sesuai status)* →
`06` → `08` → `09`. Tiap tiket: bab **Pembacaan ulang XML** dengan path + baris, blok ralat bila tiket
menyimpang dari XML, uji di seam services *(murni)* + handler + `db` *(SKIP tanpa skema uji)*.

Frontend per tiket, di dalam Shell yang sama *(kelompok menu bab 1)*: **PremiumList** *(inbox + dua
tombol)*; **Input Offer** *(section `InputOfferLife` sebagai pohon: field, `<pyCondition>`, tombol
Confirm/Reject/Decline `ConfirmSection` → `ProtectAccept`; popup Ceding/PolicyHolder/SOB/
FollowingOffer)*; **Premium List Detail** *(`ShowLifePremiumDetail` + `PL_DetailAction` +
`Calculate1_Act` → `SelectComm_SQL` `RI_COMM_LIFE`)*; **Summary** *(`ShowLifePremiumSummary`)*;
**Unggah CSV** *(`UploadCSV_LifePremium` + empat `ViewCSVResult_*`)*; label VERBATIM + baris di
`labels.premiumlist.ts` dengan uji bukti. `PARITAS-LAYAR-DAN-AKSI.md` milik modul ini: harness 9,
section, flow action 10, aksi layar — dengan kolom *ada / menunggu / tidak ditiru + bukti*.

## 4. LANGKAH 0 · LAPORAN · TELEMETRI

Langkah 0: pastikan worktree + `.env` sesuai induk §1; `git log -1` = `cfc4824` atau lebih baru;
`go test ./...`, `npm test` hijau; `-migrate` hanya sesudah cek tabrakan. Laporan akhir *(satu
pesan)*: tiket ter-commit, ralat tiket, tabel **menu → bukti → halaman → rute**, `PolisRingkas`
untuk Claim Life, kode bersama yang disentuh *(aditif)* dan sebabnya, bab TELEMETRI *(per tiket:
SHA · berkas · uji · XML dibaca · taksiran token · waktu)*.

---

*Disusun 27 September 2026 dari spec, revisi penyimpanan, struktur tabel, sepuluh tiket, sheet
struktur harness, `PremiumLife_harness.xml`, `InputPolicyHolder.xml`, daftar flow action/harness, tiga
RDBList polis, dan katalog DEV (`LIFE_PREMIUM_SUMMARY` 36 kolom, `JSON_POLIS` 151 baris, `TANGGAL_CLOSING`
— definisi dan cacah saja).*
