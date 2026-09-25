# TABLE EXTRACTION REQUEST — Claim Non Prop

> SELESAI 2026-09-18 — digantikan seluruhnya oleh PULL-LIST.csv dan ORACLE-REQUESTS.md; REQ-003 terjawab lewat ddl/

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Daftar belanja objek Oracle. Jalankan [`RECON.sql`](./RECON.sql), isi hasilnya ke [`SCHEMA-ACTUAL.csv`](./SCHEMA-ACTUAL.csv).

**Label keyakinan**
- **CONFIRMED** — nama objek muncul literal di `Rule-Connect-SQL` dalam folder `Claim Non Prop`
- **DERIVED** — disimpulkan dari konvensi penamaan Pega atas class `-Int-` yang nyata dipakai Report Definition; belum pernah terlihat literal
- **GUESS** — belum bisa dipastikan ada; **wajib diverifikasi, jangan diperlakukan sebagai ada**

---

## 0. Peringatan yang harus dibaca sebelum menarik DESC

Pega menyimpan properti work object di kolom BLOB `pzPVStream`, kecuali properti itu di-*expose* sebagai kolom nyata.

**Pemeriksaan 22 Report Definition di folder ini: tidak satu pun berjalan di atas class work** (`ASM-FW-GCNMFW-Work-ClaimTreatyNonProp` / `ASM-FW-GCNMFW-Work`). Semuanya berjalan di atas class `-Int-` (tabel Oracle bisnis), `Link-Attachment`, atau `Data-Admin-Operator-ID`.

**Konsekuensi**: tidak ada bukti XML bahwa properti klaim mana pun di-expose. Dugaan default untuk seluruh `.ClaimData.*` adalah **IN-BLOB**. Jangan kaget kalau `DESC` tabel work Pega nyaris kosong — itu memang yang diharapkan, bukan tanda data hilang.

**Karena itu urutan prioritasnya dibalik**: sumber tipe yang sahih adalah definisi property di schema PegaRULES (§1), bukan `ALL_TAB_COLUMNS` tabel work. `ALL_TAB_COLUMNS` tetap penting, tapi untuk **tabel bisnis** (§2) tempat nilai klaim disalin — dan justru di situlah presisi yang diterima akuntansi ditentukan.

---

## 1. Definisi Property Pega — PRIORITAS TERTINGGI

| Objek | Label | Untuk menjawab |
|---|---|---|
| Schema PegaRULES, tabel definisi rule (umumnya `PR4_RULE`, `PR4_BASE`, `PR4_RULE_VW`; sebagian versi memakai `pr4_rule_property`) | GUESS — struktur berbeda antar versi Pega | Tipe, max length, dan presisi setiap `Rule-Obj-Property` pada class klaim |

`RECON.sql` §1 memuat **tiga varian query** karena struktur tabel rule berbeda antar versi. Jalankan berurutan sampai ada yang mengembalikan baris.

---

## 2. Tabel bisnis — CONFIRMED

Semua muncul literal di `RDBList\*.xml`.

| Objek | Owner | Dipakai oleh | Catatan |
|---|---|---|---|
| `OS_AKSEPTASI_KLAIM` | POOLDATA + sinonim | `GetDataOS`, `GetDataCNPOS`, `CekOSClaimNonProp_SQL`, `SaveDataToOsAkseptasiNP` | **Tempat nilai akseptasi mendarat.** Kolom `DATA_JSON`, `CASEID`, `NOCLAIM`, `STS_REJECT`, `TANGGAL`, `MASTERID` |
| `json_klaim` | — | `CariHistoryClaim_SQL`, `CekHistoryClaimNonProp_SQL`, `InsertClaimPNC` | Dokumen JSON klaim; `DATA_JSON`, `NOPOLIS`, `IDPEGA`, `TGL_INPUT` |
| `JSON_POLIS` | — | `GetPolicyData`, `GetCaseIDNBTretyIn` | `DATA_JSON`, `NOPOLIS`, `prodke`, `IDPEGA` |
| `claimxol2` / `claimxol` | POOLDATA | `GetHistoryMasterID` | Riwayat alokasi XOL per case |
| `treatyinproduction` | — | `GetDataPolisNonProp_SQL`, `GetNopolis_SQL`, `GetMObyNopol_SQL` | Polis treaty-in produksi |
| `treaty_out` | — | `GetLimitTONPPLA` | `layer`, `layertype`, `layerpart`, `LIMIT`, `limit2`, `currency`, `currency2`, `deductible`, `deductible2`, `conversion` — **sumber Limit Layer & Retensi Cedant** |
| `m_treaty_in`, `m_treaty_in_edm` | — | `GetLimitsTreatyIn_SQL` | `JSONDATA` |
| `m_treaty_out`, `m_treaty_out_detail` | — | `GetDataMasterTOutNP`, `BrowseDtlTreatyOutNP` | `jsondata` |
| `TREATYINDETAIL`, `TREATYINDETAILEDM` | POOLDATA | `BrowseTreatyNP`, `BrowseDtlTreatyNP` | |
| `treatybusiness`, `TREATYYEAR`, `treaty_in`, `treatycontract` | POOLDATA | `BrowseTreatyNP`, `GetTreatyName_SQL` | |
| `PROPORTIONALARRG` | POOLDATA | `GetTreatyName_SQL` | Susunan layer |
| `bankaccount` | — | `GetDataBankAccount_sql`, `GetDatabyClientName1/2/3`, `GetClientName` | |
| `m_client`, `agent`, `marketingofficer` | — | `GetAddressCeding`, `GetLeaderReport`, `GetMObyNopol_SQL` | |
| `currency` | — | `GetCurrency` | |
| `business` | — | `GetDataBusiness_SQL` | |
| `REINSURANCETYPE` | POOLDATA | `GetReinsuranceTypeBYName_SQL` | |
| `rw`, `city` | — | `BrowseRW_SQL` | |
| `claimrejected` | — | `CekHistoryClaimNonProp_SQL` | `inskey` |
| `KODE_PRODUKSI` | POOLDATA | `GetKodeProdNonLife_SQL` | `KODE`, `TYPE` |
| `T_STORAGE_IMAGE`, `T_FOLDER_IMAGE` | POOLDATA | `Insert_T_Storage_SQL`, `GetLinkStorage_SQL` | |
| `DIRECTTOKASIR_LOG` | POOLDATA | `InsertLOGDirectKasir_SQL` | |
| `V_M_CAUSE_OF_LOSS`, `V_D_CAUSE_OF_LOSS`, `V_D_CAUSE_OF_LOSS_BUSINESS` | — | `UpdateMCauseOfLoss`, `GetLBUID_SQL` | View |
| `trloss_detail_t` | reinsurance | `getStatusKonversi_SQL` | `NO_AKSEP` |

## 3. Tabel — DERIVED (wajib diverifikasi keberadaannya)

Disimpulkan dari nama class `-Int-` yang dipakai Report Definition, belum pernah terlihat literal di SQL.

| Dugaan objek | Dari class | Dipakai oleh |
|---|---|---|
| `EMAILKOMITE` | `ASM-FW-GCNMFW-Int-EMAILKOMITE` | `FilterEmailKomiteWithLimit` — **penentu jenjang Komite** |
| `ADJUSTERCONSULTANT` | `ASM-FW-GISFW-Int-ADJUSTERCONSULTANT` | `BrowseAdjusterConsultant` |
| `LST_BANK_GROUP` | `ASM-FW-GISFW-Int-LST_BANK_GROUP` | `BrowseBankGroup` |
| `CATASTROPHE` | `ASM-FW-GCNMFW-Int-CATASTROPHE` | `GetCatastrope_RD` |
| `V_MST_USER_TEKNIS` | `ASM-FW-GCNMFW-Int-V_MST_USER_TEKNIS` | `BrowseVMstUserTeknis_RD` |
| `PROVINCE` | `ASM-FW-GISFW-Int-PROVINCE` | `BrowseProvince_RD` |
| `TREATYGROUP` | `ASM-FW-GISFW-Int-TREATYGROUP` | `BrowseTreatyGroup_RD` |
| `CURRENCYSTANDARD` | `ASM-FW-GISFW-Int-CURRENCYSTANDARD` | `CurrencyStandard` (lewat fungsi) |

## 4. Procedure / Function / Sequence — CONFIRMED

| Objek | Owner | Dipanggil dari | Kenapa perlu source-nya |
|---|---|---|---|
| `PEGA_JSON_OS_AKSEP_KLAIMTNP` | POOLDATA | `SaveDataToOsAkseptasiNP` | 11 parameter bernama `CARIn` — pemetaan ke kolom hanya ada di dalam procedure |
| `PEGA_JSON_KLAIM_PNC` | POOLDATA | `InsertClaimPNC` | idem, 6 parameter |
| `PROC_GENERATE_SEQUENCE_NUMBER` | POOLDATA | `GetSequenceNumber_SQL` | Pembentukan Nomor Akseptasi |
| `GET_TOKEN_STORAGE` | POOLDATA | `GetTokenStorage_SQL` | Token penyimpanan dokumen |
| `PEGA_M_CAUSE_OF_LOSS`, `PEGA_D_CAUSE_OF_LOSS` | POOLDATA | `UpdateMCauseOfLoss`, `UpdateDCauseOfLoss` | Upsert master penyebab kerugian |
| `getcurrencystandard` | POOLDATA | `CurrencyStandard` | **Sumber kurs** — presisi kurs menentukan seluruh konversi |
| `f_get_email` | gl | `GetEmailCeding_SQL` | |
| `PLATNP_SEQ` | — | `GenerateNoPLATNP` | Sequence nomor nota |
| `STANDARD_HASH` | (built-in) | `GenerateImageID_SQL` | — |

---

## 5. Properti yang WAJIB dikonfirmasi tipenya (MUST-CONFIRM)

Semua nilai uang, kurs, persentase share, dan rate. Kolom "Dugaan simpan" adalah tebakan saya yang **harus dibantah atau dibenarkan** oleh hasil recon.

| Properti Pega | Makna | Dugaan simpan | Status |
|---|---|---|---|
| `.ClaimData.SpreadingRisk[].ClaimEstimation` | Kerugian yang diserap satu Layer | IN-BLOB + salinan di `claimxol2` | MUST-CONFIRM |
| `.ClaimData.SpreadingRisk[].ClaimSpreaded` | Porsi reasuradur atas Layer | IN-BLOB + `OS_AKSEPTASI_KLAIM.DATA_JSON` | MUST-CONFIRM |
| `.ClaimData.SpreadingRisk[].TotalClaim` / `.TotalSpread` | Total per Layer | IN-BLOB | MUST-CONFIRM |
| `.ClaimData.SpreadingRisk[].CNPLimit` / `.CNPLimitFull` | Limit Layer | dari `treaty_out.LIMIT` / `.limit2` | MUST-CONFIRM |
| `.ClaimData.SpreadingRisk[].CNPMDP` | Premi Deposit | dari `treaty_out` / `LimitSummaryList` | MUST-CONFIRM |
| `.ClaimData.SpreadingRisk[].CNPPctReinstate` | % Premi Pemulihan | IN-BLOB | MUST-CONFIRM |
| `.ClaimData.ReinstatementList[].CNPReinstatePremium` / `.CNPReinsPremiRNM` | Premi Pemulihan | IN-BLOB | MUST-CONFIRM |
| `.ClaimData.SpreadingRisk[].Kurs` / `.KursIDR` | Kurs konversi | `getcurrencystandard`, `TreatyInMaster.CurrencyList.Conversion` | MUST-CONFIRM |
| `.TreatyInMaster.RNMShare` | Porsi Reasuradur (%) | IN-BLOB | MUST-CONFIRM |
| `.ClaimData.ShareCeding` | Share cedant (%) | IN-BLOB | MUST-CONFIRM |
| `.ClaimData.ListClaimAmount[].AltValue` | Kurs per mata uang | IN-BLOB | MUST-CONFIRM |
| `.ClaimData.ListClaimAmount[].PctProrateClaim` | Prorata Klaim | IN-BLOB (dihitung skala 5) | MUST-CONFIRM |
| `.ClaimData.DeductibleValue` / `.Amount` / `.TSIDeductible` | Retensi Cedant | IN-BLOB | MUST-CONFIRM |
| `AdjusterFee`, `Salvage`, `CNPOthersFee`, `TPL` (di beberapa PageList) | Komponen nilai | IN-BLOB + `OS_AKSEPTASI_KLAIM` | MUST-CONFIRM |
| `.AdjustmentList[].AcceptedNo` | Nomor Akseptasi | `OS_AKSEPTASI_KLAIM.NOCLAIM`? | MUST-CONFIRM (panjang string: `HitServiceToKasir_Act` memvalidasi panjang **23 atau 24** untuk non-prop) |
| `.ClaimData.NoClaim`, `.IDMaster`, `.PolicyData.PolicyNo` | Identitas | kemungkinan EXPOSED | MUST-CONFIRM |
| `EMAILKOMITE.LIMIT_BOTTOM` / `.LIMIT_TOP` | Ambang kewenangan Komite | kolom nyata | MUST-CONFIRM |
| `.ClaimData.DateOfLoss`, `.ReportDate`, `.DateReceived` | Tanggal | kemungkinan EXPOSED | SAFE-DEFAULT (DATE) |
| `.ClaimData.InsuredName`, `.ReporterName`, `.Location` | Teks bebas | IN-BLOB | SAFE-DEFAULT |
| `.ClaimData.SuggestList[]` (Kronologi) | Audit trail | IN-BLOB | SAFE-DEFAULT |

---

## 6. Yang TIDAK saya minta

- Data nasabah dalam bentuk apa pun. Seluruh query di `RECON.sql` dibatasi 20 baris dan kolom teks bebas di-masking.
- Tidak ada DML.
- Objek milik modul Komite — DEFERRED-TO-KOMITE-SESSION.
