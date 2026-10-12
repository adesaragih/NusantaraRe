# Struktur Tabel — NB FacIn: PETA TABEL WARISAN yang dibaca + tabel yang dibuat

Modul ini **membuat tujuh belas tabel** — `T_NB_OPPORTUNITY` (migrasi 180, tiket 29 / butir 76.3), `T_GENERAL_POLIS` /
`T_QUOTATIONDATA` **sebagian** (migrasi 182/183, tiket 31 / butir 78.4), `T_CEDINGCOLIST` utuh (185, tiket 34), dan tabel tab Object FIRE (186, tiket 35: `T_LOCATIONLIST` /
`T_PROPERTY` sebagian, `T_RISKLOCATION` / `T_BUILDINGCONSTRUCTION` utuh; 187, tiket 38: `T_SURROUNDINGRISK` utuh; 188, tiket 39: `T_PROPERTYITEMLIST` sebagian; 189, tiket 40: `T_OCCUPATIONLIST` sebagian, `T_TABLEOFLIMIT` utuh; 190, tiket 41: `T_FEALIST` baru; 191, tiket 42: `T_LISTCAUSEOFLOSS`, `T_COINSDATA` utuh; 193, tiket 43: `T_COVERAGELIST` sebagian; 194, tiket 45: `T_DEDUCTIBLELIST` sebagian) — dan **menulis**
baris `T_WORK_POLIS` milik premiumlistlife (K-064; tidak dibuat, tidak dipetakan kolomnya di sini). Selebihnya ia hanya
**membaca** dua puluh satu tabel / view yang sudah ada (lima di antaranya — `ACCUMULATION`, `PROVINCE`, `ACCUMULATEDTYPE`, `CZONE` milik `masterprovince` (dulu Master Data) dan `NATION` yang ditulis `masternation` — digambarkan di `modul/masterprovince/docs/STRUKTUR-TABEL-MASTER-DATA.md`) — enam tabel limit
akseptasi (tiket 20), tabel akun `T_M_ACCOUNT` (tiket 27), tabel bisnis `BUSINESS` (tiket 28), `MARKETINGOFFICER` (31), `AGENT`
(33), `RISKADDRESS` (36; juga DISISIPI tiket 37), `RW` (36/37), `OCCUPATION` (38), `V_JN_OBJ_ITEM` dan `CURRENCY` (39), `TABLEOFLIMIT` (40), `COVERAGE` (43), view `CITY` / `DISTRICT`, tabel `JSON_POLIS` dan `NATION` (46); ditambah — di luar hitungan dua puluh satu itu — view `ZONES` / `CONSTRUCTION` dan tabel `ZONEOJK` (52) — sama dengan bab "Tabel warisan" `MODUL.md`. Enam tabel limit akseptasi yang sudah ada
di `POOLDATA`, dengan nama tabel dan kolom **verbatim**. Berkas ini **peta**, bukan DDL: hanya kolom yang dibaca
repository (`backend/repository/limit.go`). Ke-enamnya dinyatakan di `MODUL.md` bab "Tabel warisan: dibaca, tidak
dibuat".

Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` (diberikan work owner, 01-10-2026; butir 42).
Kolom lain di DDL (`ID`, `MAX_LIMIT_*`, `BATAS_WAKTU`, `TGL_UPDATE`, `EFFECTIVE_DATE`, `WORKBASKET`, `JABATAN_ATASAN`,
`LIMITTRADE_BOTTOM`) tidak dibaca. ⛔ `NAMA` dan `LOGIN` **tidak pernah** dibaca (CLAUDE.md §4 butir 10, K-025).

## M_LIMIT_PROPERTYY

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga (`LetterNo`) |
| `TEAM_GROUP` | VARCHAR2(20) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_PROPERTY_NON_PREFERREDD

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(20) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_ENGINEERINGG

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_NONPROPANDENGG

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_FINANCIALINS

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga bentuk B (berspasi, mis. `DIREKTUR TEKNIK`) |
| `LIMITBOND_BOTTOM` | NUMBER(*,0) | limit Bond |
| `LIMITCREDITCL_BOTTOM` | NUMBER(*,0) | limit Kredit CL (juga Trade Credit, A26) |
| `LIMITCREDITNCL_BOTTOM` | NUMBER(*,0) | limit Kredit NCL |

## T_M_ACCOUNT

Tiket 27 (popup ChooseAccount). Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\T_M_ACCOUNT.txt` (ditambahkan work
owner 02-10-2026); kelima kolom DDL itu yang dibaca.

Sejak merge `origin/dev` 04-10-2026 modul `accounts` mengubah tabel ini (migrasi 840-842: kolom `CREATEDATE`,
`CREATEOP`, `DESCRIPTION`, kunci, sequence) dan menggambarkannya di `modul/accounts/docs/STRUKTUR-TABEL-ACCOUNTS.md`
(satu tabel bersama = satu bentuk). Yang dibaca di sini:

- `ID` — VARCHAR2(255 CHAR); identitas baris akun.
- `GROUPBUSINESSID` — VARCHAR2(32 CHAR); id group business.
- `GROUPBUSINESS` — VARCHAR2(64 CHAR); kolom layar Group Business; dicari (A69).
- `INSUREDID` — VARCHAR2(255 CHAR); kolom layar Insured ID; dicari; urutan (A72).
- `INSUREDNAME` — VARCHAR2(64 CHAR); kolom layar Insured Name; dicari.

## BUSINESS

Tiket 28 (pilihan Class Of Business). Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\BUSINESS.txt`
(`POOLDATA.BUSINESS`, 18 kolom). **Hanya tiga kolom di bawah yang dibaca**; 15 kolom lain tidak disentuh dan tidak
didaftar di sini (tabel warisan — penjaga kolom tidak membandingkannya).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(4000 BYTE) | identitas pilihan; pemutus seri urutan (A75) |
| `NOTE` | VARCHAR2(4000 BYTE) | teks pilihan yang tampil (`pyDisplayProperty .Note` di section `InputLossRecord_Sec`); NULL dibuang (A78); urutan (A74) |
| `BUSINESSGROUPID` | VARCHAR2(4000 BYTE) | saringan `= :1` (RD filter C `.BusinessGroupID = Param.Group`) |

## T_NB_OPPORTUNITY

Tiket 29 (tombol Create opportunity), butir 76.3 — **dibuat** modul ini, migrasi `180_t_nb_opportunity.sql`. Satu baris
per case NB, **berbagi PK** dengan `T_WORK_POLIS.ID` (`NB-<n>`), tanpa constraint FK (pola `T_PREMIUM_LIST`
premiumlistlife). Dasar tabel sendiri `[terverifikasi]`: di Pega opportunity adalah kelas work tersendiri
(`NB FacIn\ReportDefinition\GetListOpportunity.xml`), dan rancangan tabel flat tidak punya tabel untuknya. Tipe =
keputusan agent A81.

| Kolom | Tipe | Tipe DDL (migrasi 180) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(32) NOT NULL | PK = `T_WORK_POLIS.ID` |
| `ESTIMATED_CLOSING_DATE` | DATE | DATE | Estimated Closing Date (kabel `DD-MM-YYYY`) |
| `BUSINESS_PROSPECT_NAME` | teks | VARCHAR2(255) | Business Prospect Name |
| `ACCOUNT_ID` | teks | VARCHAR2(255 CHAR) | `T_M_ACCOUNT.ID` akun terpilih |
| `INSURED_ID` | teks | VARCHAR2(255 CHAR) | `T_M_ACCOUNT.INSUREDID` |
| `GROUP_BUSINESS_ID` | teks | VARCHAR2(32 CHAR) | `T_M_ACCOUNT.GROUPBUSINESSID` |
| `GROUP_BUSINESS` | teks | VARCHAR2(64 CHAR) | `T_M_ACCOUNT.GROUPBUSINESS` |
| `CLASS_OF_BUSINESS` | teks | VARCHAR2(4000 BYTE) | Class Of Business (lebar `BUSINESS.NOTE`) |
| `TYPE_OF_INWARD` | teks | VARCHAR2(255) | Type Of Inward |
| `TYPE_OF_FACULTATIVE` | teks | VARCHAR2(255) | Type Of Facultative (kosong bila bukan Facultative) |
| `PHASE` | teks | VARCHAR2(255) | Phase |
| `STAGE` | teks | VARCHAR2(255) | Stage |
| `OPPORTUNITY_SOURCE` | teks | VARCHAR2(255) | Opportunity Source |
| `BUSINESS_STATUS` | teks | VARCHAR2(255) | Business Status |
| `DESCRIPTION` | teks | VARCHAR2(4000) | Description |

## T_GENERAL_POLIS

Tiket 31 (blok General layar Inward Facultative), butir 78.4 — **dibuat sebagian** modul ini, migrasi
`182_t_general_polis.sql`: kolom sistem + kolom yang dipakai layar; kolom rancangan lain (termasuk kolom uang yang
menunggu keputusan presisi tim inti) ditambah tiket 23 lewat `ALTER`. Nama/tipe = rancangan (`loader/skema_gen.go`,
dijaga `TestMigrasiFlatSebagianCocokRancangan`). Berbagi PK dengan `T_WORK_POLIS` (K-064), tanpa FK lintas modul.

| Kolom | Tipe | Tipe DDL (migrasi 182) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(32) NOT NULL | PK = `T_WORK_POLIS.ID` (butir 76.1) |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader untuk data lama) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `START_DATE_TIME` | teks | VARCHAR2(30) | Begin date — `.PolicyData.StartDateTime`, teks Pega `YYYYMMDDTHHMMSS.mmm GMT`, ditulis 12:00 WIB = `T050000.000 GMT` (butir 78.1) |
| `OFFERING_DATE` | teks | VARCHAR2(30) | Offering date — `.PolicyData.OfferingDate`, teks Pega `YYYYMMDD` |
| `END_DATE_TIME` | teks | VARCHAR2(30) | End date — `.PolicyData.EndDateTime`, teks Pega |
| `FOLLOWING` | teks | VARCHAR2(50) | Old Policy Number — `.Following` (sel 72), tampil saja |
| `PERCENT_SHARE` | angka desimal | NUMBER(38,8) | tiket 48 (migrasi 198) — % Share RNM `.OfferFacIn.PercentShare` (tab Spreading sel 80), 0–100 |
| `PPN_CHECK` | teks | VARCHAR2(10) | tiket 50 (migrasi 200 NUMBER(5), diganti teks oleh migrasi 202) — "Exclude PPN ?" `.OfferFacIn.PPnCheck`, teks Pega apa adanya "true" / "false" (butir 68.1; Pega menyetelnya true tiap hitung, `SetRIComm_Act` langkah 5) |
| `PAY_RI_COMMISION` | angka desimal | NUMBER(38,8) | tiket 50 (200) — % Premium Deduction `.PolicyData.Payment.RICommision` |
| `PAY_PCT_BROKERAGE_FEE` | angka desimal | NUMBER(38,8) | tiket 50 (200) — % Brokerage Fee `.PolicyData.Payment.PctBrokerageFee` |
| `PAY_INSTALLMENT` | angka bulat | NUMBER(5) | tiket 51 (201) — Installment `.PolicyData.Payment.Installment` (jumlah cicilan, 1–99) |
| `SHARE_CEDANT_TYPE` | angka bulat | NUMBER(5) | tiket 49 (migrasi 199) — Share Cedant Type `.OfferFacIn.ShareCedantType`: 0 Gross, 1 Share RNM (`DDL\ShareCedantType.xml`); kosong = belum dipilih |
| `NB_STATUS` | teks | VARCHAR2(1000) | tiket 57 (migrasi 208) — `pyWorkPage.NBStatus` (halaman KERJA, amandemen loader: data lama kosong), ditulis Submit tahap marketing: "NB IS IN TEAM LEADER'S INBOX" / "NB IS IN <MARKETING>'S OPPORTUNITY" / "NB IS IN <jabatan>'S INBOX" / "NB WAS DECLINED BY  <OPERATOR>" |
| `IS_DEDUCTIBLE_ACCEPTANCE` | teks | VARCHAR2(10) | tiket 58 (migrasi 210) — "Deductible Acceptance Case" `.OfferFacIn.IsDeductibleAcceptance`, teks Pega "true" / "false" (butir 68.1); ditulis Submit Team Leader |
| `COMMENT_DEDUCTIBLE` | teks | VARCHAR2(500) | tiket 58 (migrasi 210) — "Comment Deductible" `.OfferFacIn.CommentDeductible`; ditulis Submit Team Leader |
| `DATE_TO_UW` | teks | VARCHAR2(30) | tiket 58 (migrasi 210) — `.OfferFacIn.DateToUW` teks Pega DateTime; Accept Team Leader (`IsJuniorUW_B`); rancangan menaruhnya di T_WORK_POLIS (premiumlistlife) — amandemen loader |
| `SUBMIT_TO_UW` | teks | VARCHAR2(30) | tiket 58 (migrasi 210) — `.OfferFacIn.SubmitToUW` teks Pega DateTime; Accept Team Leader bila kosong (`IsJuniorUW_B`); amandemen loader |
| `DAYS_VALIDITY` | teks | VARCHAR2(50) | tiket 61 (migrasi 212) — "Validity Of Acceptance" `.OfferFacIn.DaysValidity` (hari, teks bulat); ditulis Submit tahap Underwriting |
| `DATE_VALIDITY` | teks | VARCHAR2(30) | tiket 61 (migrasi 212) — "days, or up to" `.OfferFacIn.DateValidity` teks Pega DateTime 00:00 WIB (`GetDateValidity_ACT`); ditulis Submit tahap Underwriting |
| `WAITING_BIND_DATE` | teks | VARCHAR2(30) | tiket 61 (migrasi 212) — "Waiting Bind" `.OfferFacIn.WaitingBindDate` teks Pega DateTime 00:00 WIB (`GetDateValidity_ACT` L701); ditulis Submit tahap Underwriting |
| `COMMENT_TEXT` | teks | VARCHAR2(500) | tiket 61 (migrasi 212) — Comment `.OfferFacIn.Comment` Submit tahap Underwriting (juga ke riwayat KETERANGAN) |
| `CONFIRM_BINDING` | teks | VARCHAR2(50) | tiket 64 (migrasi 213) — `.OfferFacIn.ConfirmBinding` "1" / "0"; ditulis tombol Binding (`POST …/submit/binding`), direset Revise |
| `RECEIVED_RI_SLIP` | teks | VARCHAR2(50) | tiket 64 (migrasi 213) — `.OfferFacIn.ReceivedRiSlip` "true" / "false"; ditulis Submit tahap Binding (Accept), direset Revise |
| `PRO_RATE_PERCENT` | angka desimal | NUMBER(38,8) | tiket 73 P2 (migrasi 255, rentang kedua) — `.OfferFacIn.ProRatePercent` % Pro Rate tingkat case tab Coverage Aneka; ditulis PUT `…/objek-aneka`; NULL = belum diisi (dihitung dari tanggal, CPCA langkah 4) |
| `OLD_POLIS_ID` | angka bulat | NUMBER(19) | migrasi 252 — kolom sistem rancangan — hanya RENEWAL (aturan loader K-071 J-5); NULL di NB |
| `ADDITIONAL_CAPITAL` | teks | VARCHAR2(50) | migrasi 252 — medan rancangan `AdditionalCapital` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `BINDER_RNM` | teks | VARCHAR2(50) | migrasi 252 — `.OfferFacIn.BinderRNM` — tanpa Property-Set di korpus NB FacIn (Section `Periode_IsUW` + InsertTreatyProduction_Sql); belum diisi |
| `CEDING_RETENTION` | angka desimal | NUMBER(38,8) | migrasi 252 — medan rancangan `CedingRetention` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `CEDING_RETENTION_NOMINAL` | angka desimal | NUMBER(38,8) | migrasi 252 — medan rancangan `CedingRetentionNominal` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `COMMENT_CEDING` | teks | VARCHAR2(500) | migrasi 252 — medan rancangan `CommentCeding` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `CONDITION` | teks | VARCHAR2(500) | migrasi 252 — medan rancangan `Condition` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `CURRENT_YEAR` | teks | VARCHAR2(50) | migrasi 252 — medan rancangan `CurrentYear` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `DATE_OFFER` | teks | VARCHAR2(30) | migrasi 252 — `.OfferFacIn.Date` (`ViewOfferFacInUW_PreDT` baris 2) teks "MM/dd/yyyy hh:mm a" GMT; ditulis Submit tahap marketing bila masih NULL (pendekatan, keputusan work owner 12-10-2026) |
| `FLAG_SAVE_FO` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `FlagSaveFO` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `FLAG_SPECIAL_ACCEPTANCE` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `FlagSpecialAcceptance` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `FLAG_SPREADING` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `FlagSpreading` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `FOLLOWING_NB` | teks | VARCHAR2(50) | migrasi 252 — medan rancangan `FollowingNB` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `GENERATE_DATE_POLICY_NO` | teks | VARCHAR2(30) | migrasi 252 — `.PolicyData.GenerateDatePolicyNo` (`GenerateNopolis_Act` langkah 10) teks Pega UTC; ditulis Selesai Policy saat nomor polis TERBIT (pendekatan, keputusan work owner 12-10-2026) |
| `ID_FOLLOWING_NB` | teks | VARCHAR2(50) | migrasi 252 — medan rancangan `IDFollowingNB` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `INPUT_COMMENT_DEDUCTIBLE` | teks | VARCHAR2(500) | migrasi 252 — medan rancangan `InputCommentDeductible` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_B2_B` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsB2B` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_BANDING` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsBanding` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_EDIT_REFF_NUMBER` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsEditReffNumber` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_FAC_RETRO` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsFacRetro` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_GROSS_NET_SHOW` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsGrossNetShow` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_INPUT_FAC_RETRO` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsInputFacRetro` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_LIFE_PERIOD` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsLifePeriod` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_OCCUP_EXCEPTION` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsOccupException` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_PREFERRED_RISK` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsPreferredRisk` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_PRO_RATE` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsProRate` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_RI_SLIP` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsRISlip` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_SPECIAL_ACCEPTANCE` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsSpecialAcceptance` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `IS_TOP_RISK` | teks | VARCHAR2(10) | migrasi 252 — medan rancangan `IsTopRisk` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `OFFICER_FAC_OUT` | teks | VARCHAR2(50) | migrasi 252 — medan rancangan `OfficerFacOut` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `PARAM_CONDITION` | teks | VARCHAR2(500) | migrasi 252 — medan rancangan `ParamCondition` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `POLICY_MASTER_ID_PEGA` | teks | VARCHAR2(50) | migrasi 252 — `.OfferFacIn.PolicyMasterIDPega` — tanpa Property-Set di korpus NB FacIn (Obj-Open-By-Handle saja); belum diisi |
| `POLICY_MASTER_NUMBER` | teks | VARCHAR2(50) | migrasi 252 — `.OfferFacIn.PolicyMasterNumber` — tanpa Property-Set di korpus NB FacIn (hanya dibaca); belum diisi |
| `POLICY_NO` | teks | VARCHAR2(50) | migrasi 252 — TAHAP 1 — `.PolicyData.PolicyNo` nomor polis (`SaveJsonPolicyFacIn_Act` langkah 12; = JSON_POLIS.NOPOLIS); ditulis Selesai Policy |
| `PRO_RATE_TYPE` | angka bulat | NUMBER(5) | migrasi 252 — medan rancangan `ProRateType` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `PROD_DATE_TIME` | teks | VARCHAR2(30) | migrasi 252 — TAHAP 1 — `.PolicyData.ProdDateTime` (`SaveJsonPolicyFacIn_Act` langkah 5), teks Pega `YYYYMMDDTHHMMSS.000 GMT` dari JSON_POLIS.TGL_PROD (jam dinding WIB); ditulis Selesai Policy |
| `PROD_KE` | angka bulat | NUMBER(5) | migrasi 252 — TAHAP 1 — versi produksi = JSON_POLIS.PRODKE (NB `'0'`, `INSERTJSON_JSONPOLIS_FACIN`; aturan loader K-071 J-3); ditulis Selesai Policy |
| `PRODUCT_BRIGUNA` | teks | VARCHAR2(50) | migrasi 252 — medan rancangan `ProductBriguna` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `SAVE_FAC_OUT` | teks | VARCHAR2(50) | migrasi 252 — medan rancangan `SaveFacOut` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `SHARE_RNML` | angka desimal | NUMBER(38,8) | migrasi 252 — medan rancangan `ShareRNML` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `SOURCE_ID` | teks | VARCHAR2(50) | migrasi 252 — kolom turunan rancangan (loader); belum diisi aplikasi |
| `SUM_TOTAL_TSI` | angka desimal | NUMBER(38,8) | migrasi 252 — `.OfferFacIn.SumTotalTSI` (SumTSIObjItem_act / SetTotalFormula_ACT) — Go belum punya hitungannya; NULL |
| `TOTAL_PREMI_NUSA_RE` | angka desimal | NUMBER(38,8) | migrasi 252 — `.OfferFacIn.TotalPremiNusaRe` = Total Premi RNM IDR tab Cedant (`ringkasanRNM` + `totalKeIDR`), MENYIMPANG K49-11; ditulis TIAP Submit; kurs tidak ada = NULL |
| `TOTAL_RI_COM_NUSA_RE` | teks | VARCHAR2(50) | migrasi 252 — medan rancangan `TotalRIComNusaRe` — belum diisi aplikasi (tahap 2 / 3, peta di docs/issues/76) |
| `TOTAL_TSI_NUSA_RE` | angka desimal | NUMBER(38,8) | migrasi 252 — `.OfferFacIn.TotalTSINusaRe` = hitungan Go `totalTSIUW` (CountTotalTSIPremiNusaRe_Act L14, kurs terbaru), 8 desimal; ditulis TIAP Submit (keputusan work owner 12-10-2026); kurs tidak ada = NULL |
| `TOTAL_TSI_NUSA_RE_SPREADING` | angka desimal | NUMBER(38,8) | migrasi 252 — `.OfferFacIn.TotalTSINusaReSpreading` = Total TSI RNM IDR tab Cedant (`ringkasanRNM` + `totalKeIDR`), MENYIMPANG K49-11 (per mata uang -> IDR); ditulis TIAP Submit; kurs tidak ada = NULL |
| `TOTAL_TSI_TOP_RISK` | angka desimal | NUMBER(38,8) | migrasi 252 — `.OfferFacIn.TotalTSITopRisk` = hitungan Go `totalTSIUW` (L14), 8 desimal; ditulis TIAP Submit; kurs tidak ada = NULL |
| `TREATY_CAPACITY` | angka desimal | NUMBER(38,8) | migrasi 252 — `.OfferFacIn.TreatyCapacity` (SumTreatyCapacity_Act L9) — Go belum punya hitungannya; NULL |

## T_QUOTATIONDATA

Tiket 31, butir 78.4 — **dibuat sebagian**, migrasi `183_t_quotationdata.sql` (+ `SEQ_T_QUOTATIONDATA`). Satu baris
per case (`UQ_T_QUOTATIONDATA_PARENT`, A87); `PARENT_ID` → `T_GENERAL_POLIS.ID` (FK tanpa `ON DELETE`).

| Kolom | Tipe | Tipe DDL (migrasi 183) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_QUOTATIONDATA` (A92) |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (butir 76.1) |
| `NO_OFFER_SLIP` | teks | VARCHAR2(50) | Reff. number — `.QuotationData.NoOfferSlip` (sel 9) |
| `QQ_NAME` | teks | VARCHAR2(500) | QQ name — `.QuotationData.QQName` (sel 20) |
| `POLICY_TYPE` | teks | VARCHAR2(50) | Policy Type — `.QuotationData.PolicyType` (sel 26); label layar apa adanya (butir 78.2) |
| `MOID` | teks | VARCHAR2(50) | Marketing Name — `.QuotationData.MOID` (sel 75) = `MARKETINGOFFICER.ID` |
| `EDM_DAY` | teks | VARCHAR2(50) | Day — `.QuotationData.EDMDay` (sel 78) |
| `TYPE_FACULTATIVE` | teks | VARCHAR2(50) | Type facultative — `.QuotationData.TypeFacultative` (sel 43); apa adanya (butir 78.3) |
| `SOURCE_OF_BUSINESS` | teks | VARCHAR2(50) | kode SOB — `.QuotationData.SourceOfBusiness`; migrasi **184** (`ALTER … ADD`, tiket 33) |
| `SOB_NAME` | teks | VARCHAR2(500) | Source of business — `.QuotationData.SobName` (sel 48); ditulis server dari `AGENT` menurut kode (tiket 33, E-4) |
| `CEDING_CO` | teks | VARCHAR2(1000) | gabungan `;` kode Ceding Co — `.QuotationData.CedingCo`; migrasi **185** (`ALTER … ADD`, butir 80) |
| `CEDING_CO_NAME` | teks | VARCHAR2(4000) | Ceding co name — `.QuotationData.CedingCoName` (sel 49), gabungan `;` nama dari `AGENT` (tiket 34); 183 VARCHAR2(500), dilebarkan migrasi **185** (butir 80) |
| `GROUP_NAME` | teks | VARCHAR2(500) | Group Name — `.QuotationData.GroupName` (sel 56), tampil saja |
| `SHARE_OF_CEDING` | teks | VARCHAR2(50) | tiket 49 (migrasi 199) — `.QuotationData.ShareOfCeding`, teks Pega apa adanya: "100%" (Gross) atau % Share RNM + "%" (Share RNM), `SetShareOfCeding`; rancangan NUMBER, contoh 94/94 berbentuk "N%" |
| `BRANCH_CODE` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.BranchCode` = MARKETINGOFFICER.BRANCHDETAILID per MOID (`CheckDataQuotation` langkah 6); ditulis Save General (tiket 76b); lebih lebar dari kolom = NULL |
| `BRANCH_NAME` | teks | VARCHAR2(500) | migrasi 253 — `.QuotationData.BranchName` = MARKETINGOFFICER.BRANCHDETAILNAME per MOID (`CheckDataQuotation` langkah 6); ditulis Save General (tiket 76b); lebih lebar dari kolom = NULL |
| `BUSINESS_CODE` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.BusinessCode` = BUSINESS.ID baris kasus (tepat satu, ber-GROUPPANEL); ditulis Save General |
| `BUSINESS_FAC` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.BusinessFac` "F" (`InputQuotation_PreDT`, workgroup ReasFacIn); ditulis Save General |
| `BUSINESS_NAME` | teks | VARCHAR2(500) | migrasi 253 — `.QuotationData.BusinessName` = Class of Business opportunity (= BUSINESS.NOTE); ditulis Save General |
| `BUSINESS_OLD_ID` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.BusinessOldId` = BUSINESS.OLDID (GetOldIDBusiness_SQL); ditulis Save General |
| `BUSINESS_TYPE` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.BusinessType` = BusinessType_DeT (`SetBusinessType_Act`, port businessTypeDari); ditulis Save General |
| `BUSINESS_TYPE2` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.BusinessType2` "LIFE" / "NONLIFE" (`InwardFacultative_PreDT`); ditulis Save General |
| `COMMENT_TEXT` | teks | VARCHAR2(500) | migrasi 253 — medan rancangan `.QuotationData.Comment` — belum diisi aplikasi (peta di docs/issues/76b) |
| `COPY_FROM` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.CopyFrom` — belum diisi aplikasi (peta di docs/issues/76b) |
| `EMAIL` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.Email` — belum diisi aplikasi (peta di docs/issues/76b) |
| `FAC_OUT_STATUS` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.FacOutStatus` — belum diisi aplikasi (peta di docs/issues/76b) |
| `FRONTING_STATUS` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.FrontingStatus` — belum diisi aplikasi (peta di docs/issues/76b) |
| `GROUP_PANEL` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.GroupPanel` = BUSINESS.GROUPPANEL; ditulis Save General |
| `INPUT_MONTH` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.InputMonth` — belum diisi aplikasi (peta di docs/issues/76b) |
| `INPUT_QUO_DATE` | teks | VARCHAR2(30) | migrasi 253 — medan rancangan `.QuotationData.InputQuoDate` — belum diisi aplikasi (peta di docs/issues/76b) |
| `INSURED_ID` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.InsuredID` = T_NB_OPPORTUNITY.INSURED_ID (rujukan sistem baru; Pega dari POOLDATA.CLIENT); ditulis Save General |
| `INSURED_NAME` | teks | VARCHAR2(500) | migrasi 253 — `.QuotationData.InsuredName` = T_M_ACCOUNT.INSUREDNAME per ACCOUNT_ID (seperti layar); ditulis Save General |
| `IS_EDM_SHARE` | teks | VARCHAR2(10) | migrasi 253 — medan rancangan `.QuotationData.IsEDMShare` — belum diisi aplikasi (peta di docs/issues/76b) |
| `IS_GROUP` | teks | VARCHAR2(10) | migrasi 253 — medan rancangan `.QuotationData.IsGroup` — belum diisi aplikasi (peta di docs/issues/76b) |
| `IS_MOP` | teks | VARCHAR2(10) | migrasi 253 — `.QuotationData.IsMOP` (`SetMOP_Act`): POLICY_TYPE 0 -> "NonMOP", 1 / 2 -> "MOP"; ditulis Save General |
| `IS_PROPOSAL` | teks | VARCHAR2(10) | migrasi 253 — medan rancangan `.QuotationData.IsProposal` — belum diisi aplikasi (peta di docs/issues/76b) |
| `IS_SURVEY_REPORT` | teks | VARCHAR2(10) | migrasi 253 — medan rancangan `.QuotationData.IsSurveyReport` — belum diisi aplikasi (peta di docs/issues/76b) |
| `IS_TBA` | teks | VARCHAR2(10) | migrasi 253 — medan rancangan `.QuotationData.IsTBA` — belum diisi aplikasi (peta di docs/issues/76b) |
| `IS_TEMPLATE` | teks | VARCHAR2(10) | migrasi 253 — medan rancangan `.QuotationData.IsTemplate` — belum diisi aplikasi (peta di docs/issues/76b) |
| `MARKETING_CODE` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.MarketingCode` = MARKETINGOFFICER.CLIENTID per MOID (`CheckDataQuotation` langkah 6); ditulis Save General (tiket 76b); lebih lebar dari kolom = NULL |
| `MARKETING_NAME` | teks | VARCHAR2(500) | migrasi 253 — `.QuotationData.MarketingName` = MARKETINGOFFICER.CLIENTNAME per MOID (`CheckDataQuotation` langkah 6); ditulis Save General (tiket 76b); lebih lebar dari kolom = NULL |
| `OFFERING_DATE` | teks | VARCHAR2(30) | migrasi 253 — medan rancangan `.QuotationData.OfferingDate` — belum diisi aplikasi (peta di docs/issues/76b) |
| `OLD_POLICY_NO` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.OldPolicyNo` — belum diisi aplikasi (peta di docs/issues/76b) |
| `OPERATOR_ID` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.OperatorID` — belum diisi aplikasi (peta di docs/issues/76b) |
| `PERIOD` | teks | VARCHAR2(30) | migrasi 253 — medan rancangan `.QuotationData.Period` — belum diisi aplikasi (peta di docs/issues/76b) |
| `PERIOD_MM` | teks | VARCHAR2(30) | migrasi 253 — medan rancangan `.QuotationData.PeriodMM` — belum diisi aplikasi (peta di docs/issues/76b) |
| `QQID` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.QQID` — belum diisi aplikasi (peta di docs/issues/76b) |
| `SOB_LEADER0` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.SobLeader0` — belum diisi aplikasi (peta di docs/issues/76b) |
| `SOB_LSG` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.SobLsg` — belum diisi aplikasi (peta di docs/issues/76b) |
| `STATUS_BUSINESS` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.StatusBusiness` "1" (NB); ditulis Save General; NULL = salinan belum pernah ditulis (penanda backfill) |
| `STATUS_SYARIAH` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.StatusSyariah` — belum diisi aplikasi (peta di docs/issues/76b) |
| `TEAM_GROUP` | teks | VARCHAR2(50) | migrasi 253 — `.QuotationData.TeamGroup` = MARKETINGOFFICER.TEAMGROUP per MOID (`CheckDataQuotation` langkah 6); ditulis Save General (tiket 76b); lebih lebar dari kolom = NULL |
| `TYPE` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.Type` — belum diisi aplikasi (peta di docs/issues/76b) |
| `BTN_QUOTATION` | teks | VARCHAR2(50) | migrasi 253 — medan rancangan `.QuotationData.btnQuotation` — belum diisi aplikasi (peta di docs/issues/76b) |

## T_CEDINGCOLIST

Tiket 34 (daftar Ceding Co), butir 80 — **dibuat utuh** modul ini, migrasi `185_t_cedingcolist.sql` (+ `SEQ_T_CEDINGCOLIST`,
indeks `IX_CEDINGCOLIST_PARENT`). Tabel RANCANGAN (jalur `QuotationData/CedingCoList`, induk `T_QUOTATIONDATA`); nama/tipe =
rancangan, ID/PARENT_ID NUMBER(19) (A92). Satu baris per ceding, urut pilih; diganti utuh tiap *Save for later*.

| Kolom | Tipe | Tipe DDL (migrasi 185) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_CEDINGCOLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_QUOTATIONDATA.ID` (FK tanpa `ON DELETE`) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan pilih, mulai 1 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID — baris aplikasi v4 acak (A106), baris loader v5 |
| `CEDING_CO` | teks | VARCHAR2(50) | kode `AGENT.ID` — `.CedingCoList(n).CedingCo` |
| `CEDING_CO_NAME` | teks | VARCHAR2(500) | nama dari `AGENT.CLIENTNAME` — `.CedingCoList(n).CedingCoName` |

## T_CURRENCYLIST

Tiket 50 (tab Premium Deduction / Brokerage Fee) — **dibuat sebagian** modul ini, migrasi `200_premium_deduction.sql`
(+ `SEQ_T_CURRENCYLIST`, indeks `IX_CURRENCYLIST_PARENT`). Tabel RANCANGAN (jalur `CurrencyList`, induk `T_GENERAL_POLIS`);
`.Policy.Payment.*` dilipat ke awalan `PAY_` (V-22b). Satu baris per mata uang objek (urut kemunculan pertama), diperbarui
menurut nama mata uang tiap Save tab Premium Deduction / Payment (cicilan dipertahankan; mata uang yang hilang dibuang
beserta cicilannya, K50-11). Nilai uang = potret saat Save; layar menghitung ulang dari persen bila premi berubah
(K50-3).

| Kolom | Tipe | Tipe DDL (migrasi 200) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_CURRENCYLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (FK tanpa `ON DELETE`) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan mata uang, mulai 1 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak (A106) |
| `CURRENCY_CODE` | teks | VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL | = `NAME` |
| `NAME` | teks | VARCHAR2(500) | kode mata uang `.Name` (mis. IDR) |
| `PAY_PREMIUM` | angka desimal | NUMBER(38,8) | `.Policy.Payment.Premium` = Σ PremiNusantaraRe item bermata uang itu (FireMBU 6.2.5) saat Save |
| `PAY_COMMISION` | angka desimal | NUMBER(38,8) | Total Premium Deduction `.Policy.Payment.Commision` |
| `PAY_DEDUCTION2` | angka desimal | NUMBER(38,8) | Total Premium Deduction 2 `.Policy.Payment.Deduction2` — amandemen tiket 50 (K50-2) |
| `PAY_BROKERAGE_FEE` | angka desimal | NUMBER(38,8) | Total Brokerage Fee `.Policy.Payment.BrokerageFee` |
| `PAY_PPH` | angka desimal | NUMBER(38,8) | `.Policy.Payment.PPh` = tarif PPH × Brokerage Fee / 100 |
| `PAY_PPN` | angka desimal | NUMBER(38,8) | `.Policy.Payment.PPN` = tarif PPN × Brokerage Fee / 100 |
| `SUM_TOTAL_PAYMENT` | angka desimal | NUMBER(38,8) | tiket 51 (201) — Total Payment `.SumTotalPayment` = Net Premium (keputusan work owner 06-10-2026, K51-1) |
| `PAY_NET_PREMIUM` | angka desimal | NUMBER(38,8) | tiket 51 (201) — `.Policy.Payment.NetPremium` = Premium − Commision − Deduction2 − BrokerageFee + PPh + PPN |

## T_LISTINSTALLMENT

Tiket 51 (tab Payment) — **dibuat sebagian** modul ini, migrasi `201_payment.sql` (+ `SEQ_T_LISTINSTALLMENT`, indeks
`IX_LISTINSTALLMENT_PARENT`). Tabel RANCANGAN (jalur `CurrencyList/ListInstallment`, induk `T_CURRENCYLIST`). Satu baris
per cicilan per mata uang, diganti utuh tiap Save tab Payment; nilai uang = potret saat Save (layar menghitung ulang dari
persen dan nilai mata uang terkini, K51-4).

| Kolom | Tipe | Tipe DDL (migrasi 201) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_LISTINSTALLMENT` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_CURRENCYLIST.ID` (FK tanpa `ON DELETE`) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan cicilan, mulai 1 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak (A106) |
| `CURRENCY_CODE` | teks | VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL | = mata uang induk |
| `INSTALLMENT_NO` | angka bulat | NUMBER(5) | No `.InstallmentNo` |
| `DUE_DATE` | teks | VARCHAR2(30) | Due Date `.DueDate`, teks Pega `YYYYMMDD` |
| `INSTALLMENT_PERCENTAGE` | teks | VARCHAR2(50) | % Installment `.InstallmentPercentage`, teks desimal apa adanya (skala sampai 20; K51-3) |
| `PREMIUM` | angka desimal | NUMBER(38,8) | Premium RNM `.Premium` = Premium mata uang × % / 100 |
| `RI_COMMISION` | angka desimal | NUMBER(38,8) | Deduction `.RICommision` = Commision × % / 100 |
| `DEDUCTION2` | angka desimal | NUMBER(38,8) | Deduction 2 `.Deduction2` = Deduction2 × % / 100 (4 desimal) |
| `BROKERAGE_FEE` | angka desimal | NUMBER(38,8) | Total Brokerage Fee `.BrokerageFee` × % / 100 |
| `PPH` | angka desimal | NUMBER(38,8) | `.PPh` × % / 100 |
| `PPN` | angka desimal | NUMBER(38,8) | `.PPN` × % / 100 |
| `DISCOUNT` | angka desimal | NUMBER(38,8) | Discount `.Discount` (NB: 0) |
| `PAYMENT_TOTAL` | angka desimal | NUMBER(38,8) | Total Payment `.PaymentTotal` = Net Premium × % / 100 |

## T_CEDINGCEDANTLIST

Tiket 49 (tab Inw Fac Cedant Panels) — **dibuat utuh** modul ini, migrasi `199_cedant.sql` (+ `SEQ_T_CEDINGCEDANTLIST`,
indeks `IX_CEDINGCEDANTLIST_PARENT`). Tabel RANCANGAN (jalur `CedingCedantList`, induk `T_GENERAL_POLIS`, kelas
Data-Quotation). Satu baris per cedant, urut grid; diganti utuh tiap Save tab Cedant. `CurrencyList` per cedant
(`T_CEDING_CURRENCYLIST`) belum dibuat (C-1).

| Kolom | Tipe | Tipe DDL (migrasi 199) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_CEDINGCEDANTLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (FK tanpa `ON DELETE`) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan grid, mulai 1 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak (A106) |
| `CEDING_CO` | teks | VARCHAR2(50) | kode `AGENT.ID` — `.CedingCedantList(n).CedingCo` (Choose Ceding) |
| `CEDING_CO_NAME` | teks | VARCHAR2(500) | nama — `.CedingCedantList(n).CedingCoName` |
| `SHARE_CEDING` | angka desimal | NUMBER(38,8) | % Share cedant — `.CedingCedantList(n).ShareCeding`, 0–100 |

## T_SCORINGRISK

Tiket 54 (tab Scoring Risk) — **dibuat utuh** modul ini, migrasi `205_scoring_risk.sql` (+ `SEQ_T_SCORINGRISK`). Tabel
RANCANGAN (jalur `ScoringRisk`, induk `T_GENERAL_POLIS`, kelas Data-ScoringRisk). Satu baris per kasus (UNIQUE
`PARENT_ID`); kelima tabel Scoring Risk diganti utuh tiap Save tab Scoring Risk. GET menghitung ulang dari masukan tersimpan
(kolom hasil di bawah = hasil hitung saat Save, untuk laporan / tahap lain).

| Kolom | Tipe | Tipe DDL (migrasi 205) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SCORINGRISK` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (FK, UNIQUE) |
| `FINAL_SCORE` | angka desimal | NUMBER(38,8) | `.FinalScore` — Score Top Risk, atau rata-rata dua kolom 2 desimal |
| `NOTE_FINAL_SCORE` | teks | VARCHAR2(500) | `.NoteFinalScore` — "Good" / "Above Average" / "Average" / "Below Average" / "Poor" |
| `STATUS` | teks | VARCHAR2(50) | `.Status` — "New Business" (ID kasus memuat "NB") / "Renewal" |

## T_SCORINGRESULT

Tiket 54 — **dibuat utuh**, migrasi 205 (+ `SEQ_T_SCORINGRESULT`). Jalur `ScoringRisk/ScoringResult`, induk
`T_SCORINGRISK`, satu baris (UNIQUE `PARENT_ID`). Bendera = teks "true" / "false" (butir 68.1; `Activity\ScoringResult.xml`).

| Kolom | Tipe | Tipe DDL (migrasi 205) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SCORINGRESULT` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_SCORINGRISK.ID` (FK, UNIQUE) |
| `GOOD` | teks | VARCHAR2(50) | FinalScore 86–100 |
| `ABOVE_AVERAGE` | teks | VARCHAR2(50) | FinalScore 76–<86 |
| `AVERAGE` | teks | VARCHAR2(50) | FinalScore 60–<76 |
| `BELOW_AVERAGE` | teks | VARCHAR2(50) | FinalScore 30–<60 |
| `POOR` | teks | VARCHAR2(50) | FinalScore 0–<30 |
| `ACCEPTABLE_NOT_APPROVAL` | teks | VARCHAR2(50) | FinalScore 30–100 tanpa persetujuan |
| `ACCEPTABLE_WITH_APPROVAL` | teks | VARCHAR2(50) | FinalScore 30–100 dengan persetujuan |
| `NOT_ACCEPTABLE` | teks | VARCHAR2(50) | FinalScore 0–<30 |
| `OCCUPATION` | teks | VARCHAR2(50) | "Due to" Occupation (Occupation ChechBox4 / 5) |
| `OBJECT_CONDITIONS` | teks | VARCHAR2(500) | "Due to" Object Conditions (Class Construction 3 / 4, Building 4, Machinery 4) |
| `LOSS_RATIO` | teks | VARCHAR2(50) | "Due to" Loss Ratio (okupasi × band loss ratio) |
| `OTHERS` | teks | VARCHAR2(50) | "Due to" Others (Total Sum Insured ChechBox1) |
| `OPERATIONAL_DIV_HEAD` | teks | VARCHAR2(50) | "Approval From" Operational Div. Head |
| `OPERATIONAL_DIRECTOR` | teks | VARCHAR2(50) | "Approval From" Operational Director (TotalRNMShare dibandingkan NUMERIK) |
| `TECHNICAL_DIRECTOR` | teks | VARCHAR2(50) | "Approval From" Technical Director (rumus Pega dijaga OperationalDirector — selalu sama) |
| `PRESIDENT_DIRECTOR` | teks | VARCHAR2(50) | "Approval From" President Director (Occupation ChechBox5) |

## T_DATASCORINGRISKLIST

Tiket 54 — **dibuat utuh** + satu kolom amandemen, migrasi 205 (+ `SEQ_T_DATASCORINGRISKLIST`, indeks
`IX_DATASCORING_PARENT`). Jalur `ScoringRisk/DataScoringRiskList`, induk `T_SCORINGRISK`: `SEQ_NO` 1 = Top Risk, 2 =
Dominant Risk (hanya bila lokasi dominan dipilih dan ber-okupasi). Lebar teks = lebar sumbernya (pola butir 80).

| Kolom | Tipe | Tipe DDL (migrasi 205) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_DATASCORINGRISKLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_SCORINGRISK.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | 1 Top Risk, 2 Dominant Risk |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak (A106) |
| `OCCUPATION_CODE` | teks | VARCHAR2(1000) | `.OccupationCode` = OccupationList(1).OccupationId lokasi (rancangan 50) |
| `OCCUPATION_NOTE` | teks | VARCHAR2(1000) | `.OccupationNote` = OccupationList(1).OccupationName (rancangan 500) |
| `SCORE` | angka desimal | NUMBER(38,8) | `.Score` = Σ ScorePerFactor kolom |
| `TOP_RISK_LOCATION` | teks | VARCHAR2(4000) | `.TopRiskLocation` = RiskLocation.ASMAddress lokasi (rancangan 50) |
| `INDEX_LOCATION` | teks | VARCHAR2(50) | `.IndexLocation` — ObjectNo lokasi Dominant Risk (amandemen tiket 54; kosong di kolom 1) |

## T_SCORING_FACTOR

Tiket 54 — **dibuat utuh** + satu kolom amandemen, migrasi 205 (+ `SEQ_T_SCORING_FACTOR`, indeks
`IX_SCORING_FACTOR_PARENT`). Jalur `…/ScoringFactor` (aturan loader V-30), induk `T_DATASCORINGRISKLIST`: 24 baris per kolom,
`SEQ_NO` = urutan `models.DaftarFaktorScoring` (urut `SetScore_Act` langkah 5).

| Kolom | Tipe | Tipe DDL (migrasi 205) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SCORING_FACTOR` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_DATASCORINGRISKLIST.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | 1–24 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak |
| `CHECKED` | teks | VARCHAR2(10) | `.Checked` "true" / "false" |
| `FACTOR_GROUP` | teks | VARCHAR2(50) | nama halaman kelompok (RiskLocation / ObjectConditions / FEA / LossRatio / Others); NULL = faktor langsung |
| `FACTOR_NAME` | teks | VARCHAR2(500) | nama halaman faktor Pega (mis. `Riot`, `Stockinthepolicycover`) |
| `REMARKS` | teks | VARCHAR2(500) | `.Remarks` — tidak ditulis tab ini (bukan medan kontrak tiket 54) |
| `SCORE_PER_FACTOR` | angka desimal | NUMBER(38,8) | `.ScorePerFactor` = poin opsi tercentang, 0 bila tidak |
| `PERCENT_LOSS_RATIO` | angka desimal | NUMBER(38,8) | `.LossRatio.LossRatio.PercentLossRatio` — hanya baris faktor LossRatio (amandemen tiket 54) |

## T_SCORING_OPTION

Tiket 54 — **dibuat utuh**, migrasi 205 (+ `SEQ_T_SCORING_OPTION`, indeks `IX_SCORING_OPTION_PARENT`). Jalur
`…/ScoringOption`, induk `T_SCORING_FACTOR`: satu baris per ChechBoxN (101 per kolom).

| Kolom | Tipe | Tipe DDL (migrasi 205) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SCORING_OPTION` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan (diisi loader) |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan (diisi loader) |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_SCORING_FACTOR.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | = nomor opsi |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak |
| `CHECK_BOX` | teks | VARCHAR2(10) | `.ChechBoxN` "true" / "false"; dibaca balik sebagai `pilihan` (nomor terkecil yang "true") |
| `OPTION_NO` | teks | VARCHAR2(50) | N |
| `SCORE` | angka desimal | NUMBER(38,8) | `.ScoreN` = poin `ValueScoringRisk_act` bila tercentang, 0 bila tidak |

## T_LOSSRECORD

Tiket 55 (tab utama Loss Record, tingkat CASE) — **tabel BARU** modul ini (bukan rancangan loader; keputusan work owner
06-10-2026 butir 4), migrasi `206_loss_record_kasus.sql` (+ `SEQ_T_LOSSRECORD`). Satu baris per kasus (UNIQUE `PARENT_ID`),
induk `T_GENERAL_POLIS`. Diganti utuh bersama `T_LOSSRECORDLIST` tiap Save tab Loss Record. GET menghitung ulang dari
masukan tersimpan (rasio di bawah = hasil hitung saat Save).

| Kolom | Tipe | Tipe DDL (migrasi 206) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_LOSSRECORD` |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (FK, UNIQUE) |
| `PREMI_PAR` | angka desimal | NUMBER(38,8) | `.OfferFacIn.PremiPAR`; bawaan bila kosong / 0 = Σ premi spread × kurs Begin date, 0 desimal |
| `PREMI_MB` | angka desimal | NUMBER(38,8) | `.OfferFacIn.PremiMB` |
| `PREMI_PL` | angka desimal | NUMBER(38,8) | `.OfferFacIn.PremiPL` |
| `PREMI_OTHER` | angka desimal | NUMBER(38,8) | `.OfferFacIn.PremiOther` |
| `LOSS_RATIO1_YEAR_PERCENT` | angka desimal | NUMBER(38,8) | `.LossRatio1YearPercent` = Y1 / Gross × 100 |
| `LOSS_RATIO3_YEAR_PERCENT` | angka desimal | NUMBER(38,8) | `.LossRatio3YearPercent` = Y3 / (Gross × 3) × 100 |
| `LOSS_RATIO5_YEAR_PERCENT` | angka desimal | NUMBER(38,8) | `.LossRatio5YearPercent` = Y5 / (Gross × 5) × 100 |
| `NET_LOSS_RATIO1_YEAR` | angka desimal | NUMBER(38,8) | `.NetLossRatio1Year` = Y1 / (Gross × (1 − RI Comm/100)) × 100 |
| `NET_LOSS_RATIO3_YEAR` | angka desimal | NUMBER(38,8) | `.NetLossRatio3Year` (× 3) |
| `NET_LOSS_RATIO5_YEAR` | angka desimal | NUMBER(38,8) | `.NetLossRatio5Year` (× 5) |

## T_LOSSRECORDLIST

Tiket 55 — **tabel BARU**, migrasi 206 (+ `SEQ_T_LOSSRECORDLIST`, indeks `IX_LOSSRECORDLIST_PARENT`). Baris
`.OfferFacIn.ListCauseOfLoss` (kelas Data-CauseOfLoss) tingkat case — BEDA dari `T_LISTCAUSEOFLOSS` (per lokasi, tiket 42).
Induk `T_LOSSRECORD`; urut `SEQ_NO` = Date of Loss terbaru dulu (`SaveLossRecord` langkah 9).

| Kolom | Tipe | Tipe DDL (migrasi 206) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_LOSSRECORDLIST` |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_LOSSRECORD.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan, mulai 1 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak |
| `UW_YEAR` | teks | VARCHAR2(10) | `.UWYEAR` = tahun ADD_MONTHS(Begin, −12L) jendela yang memuat DOL; kosong = di luar 5 jendela |
| `TAHUN_KE` | angka bulat | NUMBER(5) | `.TahunKe` = L (1..5); NULL = di luar jendela |
| `CLASS_OF_BUSINESS` | teks | VARCHAR2(4000) | `.ClassOfBusiness` = BUSINESS.NOTE (saran `GET /class-of-business/saran`) |
| `DATE_OF_LOSS` | teks | VARCHAR2(30) | `.DOL` teks Pega DateTime pukul 12:00 WIB (pola T_LISTCAUSEOFLOSS); kosong saat isi = Begin − 1 tahun |
| `CAUSE_OF_LOSS` | teks | VARCHAR2(500) | `.CauseOfLoss` huruf besar |
| `CURRENCY` | teks | VARCHAR2(10) | `.Currency` huruf besar, dipangkas |
| `LOSS_AMOUNT` | angka desimal | NUMBER(38,8) | `.LOSS_AMOUNT` |
| `LOSS_AMOUNT_IDR` | angka desimal | NUMBER(38,8) | `.LOSS_AMOUNT_IDR` = LOSS_AMOUNT (IDR / mata uang kosong) atau × KURS_IDR; NULL = kurs tidak ada |
| `KURS_IDR` | angka desimal | NUMBER(38,10) | `.KursIDR` = TOIDR TREATYEXCHANGEYEARLY pada Begin date (non-IDR saja) |

## T_CORRESPONDENCELIST

Tiket 56 (tab Correspondence) — **tabel BARU** modul ini (bukan rancangan loader; tidak ada di `DDL\`; sampel CONTOH semua
kosong), migrasi `207_korespondensi.sql` (+ `SEQ_T_CORRESPONDENCELIST`, indeks `IX_CORRESPONDENCE_PARENT`). Baris
`.OfferFacIn.CorrespondenceList` (kelas Data-Correspondence), induk `T_GENERAL_POLIS`, urut `SEQ_NO` = urutan Add.

| Kolom | Tipe | Tipe DDL (migrasi 207) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_CORRESPONDENCELIST`; = `no` kontrak |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan Add (terbesar + 1) |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak |
| `TANGGAL` | teks | VARCHAR2(30) | `.Date` teks Pega DateTime saat Add (kabel `DD-MM-YYYY HH:MM` WIB) |
| `SENDER` | teks | VARCHAR2(150) | `.Sender` = M_LOGIN_GO.NAME akun yang menambah |
| `SENDER_MAIL` | teks | VARCHAR2(254) | `.SenderMail` = M_LOGIN_GO.EMAIL akun itu (Cc + Reply-To saat kirim) |
| `SUBJECT` | teks | VARCHAR2(1000) | `.Subject` (tanpa CR / LF) |
| `MESSAGE` | teks | VARCHAR2(4000) | `.Message` HTML yang SUDAH dibersihkan server (CLOB dilarang penjaga inti → ≤ 4000 bita) |
| `IS_EDITABLE` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | `.IsEditable` "1" draf / "0" terkirim (terkunci) |
| `SENT_AT` | teks | VARCHAR2(30) | waktu kirim berhasil, teks Pega DateTime |
| `CREATE_OP` | teks | VARCHAR2(64) | LOGIN_ID akun yang menambah |
| `TGL_CREATE` | date | DATE | SYSDATE saat Add |

## T_RECIPIENTLIST

Tiket 56 — **tabel BARU**, migrasi 207 (+ `SEQ_T_RECIPIENTLIST`, indeks `IX_RECIPIENTLIST_PARENT`). Baris `.RecipientList`
induk `T_CORRESPONDENCELIST`, diganti utuh tiap Save / Send draf.

| Kolom | Tipe | Tipe DDL (migrasi 207) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_RECIPIENTLIST` |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_CORRESPONDENCELIST.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan grid |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak |
| `RECIPIENT` | teks | VARCHAR2(150) | `.Recipient` nama (akun M_LOGIN_GO terpilih atau teks bebas) |
| `RECIPIENT_MAIL` | teks | VARCHAR2(254) | `.RecipientMail` alamat (To saat kirim) |
| `IS_EDITABLE` | teks | VARCHAR2(1) DEFAULT '1' NOT NULL | ikut induk ("0" sesudah terkirim) |

## T_CLAUSELIST

Tiket 47 (tab Clauses kasus FIRE) — **dibuat utuh** modul ini, migrasi `197_t_clauselist.sql` (+ `SEQ_T_CLAUSELIST`, indeks
`IX_CLAUSELIST_PARENT`). `pyWorkPage.ClauseList` (kelas Data-Clause) di akar kasus, induk `T_GENERAL_POLIS`, BESERTA
`.ClauseList(n).ArgumentList` dalam satu tabel (K47-4, work owner 05-10-2026 "bisa jadi 1 tabel aja"): **satu baris per
argumen**, kolom klausa berulang; klausa tanpa argumen = satu baris ber-`ARGUMENT_SEQ_NO` kosong. TIDAK ada di workbook
rancangan — rancangan agent berpola rancangan (K47-1). Lebar teks 4000 (K47-2). Diganti utuh tiap Save tab.

| Kolom | Tipe | Tipe DDL (migrasi 197) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_CLAUSELIST` (satu per baris argumen) |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (ID kasus; FK tanpa `ON DELETE`) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan klausa, mulai 1 — sama di semua baris satu klausa |
| `ARGUMENT_SEQ_NO` | angka bulat | NUMBER(5) | urutan argumen di klausa itu, mulai 1; kosong = klausa tanpa argumen |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak (A106) |
| `CLAUSE_CODE` | teks | VARCHAR2(4000) | `.ClauseCode` = `.ID` view `CLAUSE` (JSON `ID` `M_CLAUSE`) hasil Choose Clause |
| `CLAUSE_TITLE` | teks | VARCHAR2(4000) | `.ClauseTitle` |
| `CLAUSE_DESCRIPTION` | teks | VARCHAR2(4000) | `.ClauseDescription` (= `.Info` hasil cari) |
| `CLAUSE_LANGUAGE` | teks | VARCHAR2(10) | `.ClauseLanguage` "0" Indonesia / "1" Inggris / "2" lain (PostAct) |
| `CLAUSE_LANGUAGE_ID` | teks | VARCHAR2(4000) | `.ClauseLanguageID` (= `.Language` hasil cari) |
| `CLAUSE_CONTENT` | teks | VARCHAR2(4000) | `.ClauseContent` — isi sesudah argumen diganti; > 4000 bita ditolak (400) |
| `CLAUSE_CONTENT_TEMP` | teks | VARCHAR2(4000) | `.ClauseContentTemp` — isi asli (`.Text`), dasar penggantian |
| `ARGUMENT_COUNT` | teks | VARCHAR2(50) | `.ArgumentCount` — teks angka (RetrieveArgumentNumberClauseSQL `count(*)`) |
| `ARGUMENT_NUMBER` | teks | VARCHAR2(4000) | `.ArgumentList(n).ArgumentNumber` (`M_ARGCLAUSEFIRE` JSON `ArgumentNumber`) |
| `ARGUMENT_DESCRIPTION` | teks | VARCHAR2(4000) | `.ArgumentList(n).ArgumentDescription` (JSON `Description`) |
| `ARGUMENT_VALUE` | teks | VARCHAR2(4000) | `.ArgumentList(n).ArgumentValue` — Isi Argumen (awal JSON `DefaultValue`) |

## T_LOCATIONLIST

Tiket 35 (tab Object FIRE) — migrasi `186_t_objek_fire.sql`, **sebagian** (pola butir 78.4): kolom sistem saja; kolom
uang `LOSS_RATIO*_AMOUNT` (presisi tim inti) dan kolom lain rancangan ditambah tiket 23. Jalur rancangan `LocationList`,
induk `T_GENERAL_POLIS`. Satu baris per objek, urut `SEQ_NO`; diganti utuh tiap Save.

Tiket 42 — migrasi 191 menambah empat kolom loss ratio rancangan (`ALTER TABLE … ADD`); diisi **hasil hitung server** saat
Save (W-4 / `SetLossRatio_Act`), tidak dari layar.

| Kolom | Tipe | Tipe DDL (migrasi 186/191) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_LOCATIONLIST` (A92) |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (butir 76.1; FK) |
| `SEQ_NO` | angka bulat | NUMBER(10) NOT NULL | urutan objek 1..n — NUMBER(5) di 186, diperlebar migrasi 257 (tiket 75: 1 juta lokasi) |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID (aplikasi v4, A106) |
| `LOSS_RATIO1_YEAR_AMOUNT` | angka desimal | NUMBER(38,8) | LR 1 Year = ΣClaim/ΣAmount (≤ 365 hari), 8 desimal (A148) — **uang** |
| `LOSS_RATIO1_YEAR_PERCENT` | teks | VARCHAR2(50) | %LR 1 Years = ΣClaim·100/ΣAmount (8 desimal, A148), teks desimal |
| `LOSS_RATIO35_YEAR_AMOUNT` | angka desimal | NUMBER(38,8) | LR 3 - 5 Years (≤ 1825 hari, kumulatif) — **uang** |
| `LOSS_RATIO35_YEAR_PERCENT` | teks | VARCHAR2(50) | %LR 3 - 5 Years, teks desimal |

## T_PROPERTY

Tiket 35 — migrasi 186, **sebagian**: tanpa `TOTAL_TSI` (uang) — ditambah tiket 23. Tiket 38 — migrasi 187 menambah
empat kolom rancangan `OWNERSHIP`, `IS_PRODUCTION_PROCESS_FLAG`, `IS_HOT_WORK_PROCESS_FLAG`, `IS_FLAMMABLE_ITEM_FLAG`
(`ALTER TABLE … ADD`, tipe rancangan).
Satu baris per lokasi (`UQ_T_PROPERTY_PARENT`). Boolean disimpan teks `true`/`false` (bentuk data fixture).

| Kolom | Tipe | Tipe DDL (migrasi 186/187/192) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_PROPERTY` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_LOCATIONLIST.ID` |
| `ALM_RISK_ID` | teks | VARCHAR2(100) | Risk Address ID — `.Property.AlmRiskID` (sel 39) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `BUILDING_NO` | teks | VARCHAR2(50) | Building No. (sel 30) |
| `COUNTRY` | teks | VARCHAR2(100) | Country (sel 32) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `IS_MATERIAL_DAMAGE` | teks | VARCHAR2(10) | Material Damage (sel 11) `true`/`false` |
| `IS_TOP_RISK` | teks | VARCHAR2(10) | Top Risk (sel 12) `true`/`false` |
| `OBJECT_NAME` | teks | VARCHAR2(500) | Object Name (sel 13) |
| `OBJECT_NO` | teks | VARCHAR2(50) | Object No. (sel 5) |
| `OBJECT_TYPE` | teks | VARCHAR2(50) | Object Type (sel 6, wajib) |
| `PROVINCE` | teks | VARCHAR2(100) | Province (sel 38) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `ROAD_NAME` | teks | VARCHAR2(4000) | Address — `.Property.RoadName` (sel 29) — lebar 4000 sejak migrasi 192 (butir 87/88; semula 500) |
| `ROAD_TYPE` | teks | VARCHAR2(100) | Type — `.Property.RoadType` (sel 28) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `CURRENCY_CODE` | teks | VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL | K-069 (bawaan; tidak ditulis layar) |
| `OWNERSHIP` | teks | VARCHAR2(50) | Ownership — `.Property.Ownership` (`RiskAround` sel 65); kode apa adanya (187) |
| `IS_PRODUCTION_PROCESS_FLAG` | teks | VARCHAR2(10) | `.Property.IsProductionProcessFlag` (sel 77) `true`/`false` (187) |
| `IS_HOT_WORK_PROCESS_FLAG` | teks | VARCHAR2(10) | `.Property.IsHotWorkProcessFlag` (sel 78) `true`/`false` (187) |
| `IS_FLAMMABLE_ITEM_FLAG` | teks | VARCHAR2(10) | `.Property.IsFlammableItemFlag` (sel 79) `true`/`false` (187) |

## T_RISKLOCATION

Tiket 35 — migrasi 186, **utuh**. Satu baris per property (`UQ_T_RISKLOCATION_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 186/192) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_RISKLOCATION` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `ASM_ADDRESS` | teks | VARCHAR2(4000) | Risk Location — `.Property.RiskLocation.ASMAddress` (sel 40; kolom grid "Location") — lebar 4000 sejak migrasi 192 (butir 87/88; semula 50) |
| `ASM_CITY` | teks | VARCHAR2(100) | City (sel 36) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `ASM_DISTRICT` | teks | VARCHAR2(100) | District (sel 37) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `ASMRW` | teks | VARCHAR2(100) | Territory — `.Property.RiskLocation.ASMRW` (sel 35) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |
| `ASM_ZIP_CODE` | teks | VARCHAR2(100) | Zip Code (sel 31) — lebar 100 sejak migrasi 192 (butir 87/88; semula 50) |

## T_BUILDINGCONSTRUCTION

Tiket 35 — migrasi 186, **utuh** + tiga kolom **baru** (bukan di rancangan; A110, amandemen loader `amandemenBangunan`).
Satu baris per property (`UQ_T_BUILDINGCONSTR_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 186) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_BUILDINGCONSTRUCTION` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `FLOOR_TYPE` | teks | VARCHAR2(50) | Floor Type (sel 52) — kode apa adanya (G-4) |
| `NUMBER_OF_FLOOR` | teks | VARCHAR2(50) | Number of Floor (sel 49) — rancangan teks; diperiksa bilangan bulat ≥ 0 |
| `ROOF_TYPE` | teks | VARCHAR2(50) | Roof Type (sel 50) — kode apa adanya |
| `WALL_TYPE` | teks | VARCHAR2(50) | Wall Type (sel 51) — kode apa adanya |
| `PARTITION_TYPE` | teks | VARCHAR2(50) | **baru** — `.Property.BuildingConstruction.PartitionType` (sel 56) |
| `SUPPORT_WALL_TYPE` | teks | VARCHAR2(50) | **baru** — `.SupportWallType` (sel 57) |
| `OTHERS_TYPE` | teks | VARCHAR2(50) | **baru** — `.OthersType` (sel 58) |

## T_SURROUNDINGRISK

Tiket 38 — migrasi 187, **utuh** (rancangan: jalur `LocationList/Property/SurroundingRisk`, induk `T_PROPERTY`) + 18 kolom
**baru** (bukan di rancangan; A130, amandemen loader `amandemenSekitar`) dari `NB FacIn\Section\RiskAround.xml`.
Satu baris per property (`UQ_T_SURROUNDINGRISK_PARENT`). Kode disimpan apa adanya — server tidak mencocokkan ke daftar
pilihan (aturan properti `DDL\FrontConstruction.xml` / `Ownership.xml` / `FloodAreaStatus.xml` /
`HousekeepingStatus.xml` / `FloodArea.xml` dipakai frontend).

| Kolom | Tipe | Tipe DDL (migrasi 187) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SURROUNDINGRISK` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `FLOOD_AREA_STATUS` | teks | VARCHAR2(50) | Flood Area Status (sel 67) — kode apa adanya |
| `HOUSEKEEPING_STATUS` | teks | VARCHAR2(50) | Housekeeping Status (sel 66) — kode apa adanya |
| `FRONT_OCCUPATION` | teks | VARCHAR2(1000) | **baru** — `.Property.SurroundingRisk.FrontOccupation` (sel 9) = `OCCUPATION.OLDID` (lebar = sumber, butir 80) |
| `FRONT_CONSTRUCTION` | teks | VARCHAR2(500) | **baru** — `.FrontConstruction`; nilai standar terpanjang 215 bita |
| `FRONT_DISTANCE` | teks | VARCHAR2(50) | **baru** — `.FrontDistance`; teks angka 0..100000, ≤ 2 desimal (A129) |
| `FRONT_NOTE` | teks | VARCHAR2(1000) | **baru** — `.FrontNote` (diisi `OCCUPATION.NAME`; lebar = sumber, butir 80) |
| `LEFT_OCCUPATION` | teks | VARCHAR2(1000) | **baru** — sisi Left, pola sama (sel 23–28) |
| `LEFT_CONSTRUCTION` | teks | VARCHAR2(500) | **baru** |
| `LEFT_DISTANCE` | teks | VARCHAR2(50) | **baru** |
| `LEFT_NOTE` | teks | VARCHAR2(1000) | **baru** |
| `BACK_OCCUPATION` | teks | VARCHAR2(1000) | **baru** — sisi Back (sel 37–42) |
| `BACK_CONSTRUCTION` | teks | VARCHAR2(500) | **baru** |
| `BACK_DISTANCE` | teks | VARCHAR2(50) | **baru** |
| `BACK_NOTE` | teks | VARCHAR2(1000) | **baru** |
| `RIGHT_OCCUPATION` | teks | VARCHAR2(1000) | **baru** — sisi Right (sel 51–56) |
| `RIGHT_CONSTRUCTION` | teks | VARCHAR2(500) | **baru** |
| `RIGHT_DISTANCE` | teks | VARCHAR2(50) | **baru** |
| `RIGHT_NOTE` | teks | VARCHAR2(1000) | **baru** |
| `FLOOD_AREA` | teks | VARCHAR2(50) | **baru** — `.FloodArea` (sel 68); kode apa adanya (aturan properti `DDL\FloodArea.xml`: 4 nilai, 1 bita) |
| `HOUSEKEEPING_REMARK` | teks | VARCHAR2(500) | **baru** — `.HousekeepingRemark` (sel 69) |

## T_PROPERTYITEMLIST

Tiket 39 — migrasi 188, **sebagian** (pola A109): rancangan jalur `LocationList/Property/PropertyItemList`, induk
`T_PROPERTY`, **banyak** baris per property urut `SEQ_NO` (= `PROPERTY_ITEM_NO` 1..n). Tanpa kolom yang tidak dipakai layar
ini (`CURRENCY_ID`, `CURRENCY_OLD_ID`, `FLAG_NET_RATE`, `IS_OLD_DATA`, `PERCENTAGE_ADJUSTMENT`, `PROPERTY_ID`,
`SELECTED_LOCATION_ADDRESS`, `SELECTED_OBJECT_ITEM`, `TOTAL_PREMIUM_NUSANTARA_RE`) — ditambah tiket 23; `TOTAL_GROSS_PREMI` /
`TOTAL_NET_RATE` ditambah tiket 43 (migrasi 193 `ALTER ADD`). Enam kolom **baru** (A132, amandemen loader `amandemenItem`). Uang dan
persen `NUMBER(38,8)` (rancangan `NUMBER` polos; ADR-0016), ditulis/dibaca sebagai teks desimal bertitik — nol float.

| Kolom | Tipe | Tipe DDL (migrasi 188) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_PROPERTYITEMLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan item 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `CURRENCY` | teks | VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL | Currency (`CURRENCY.CURRENCY`, tanpa ITL); **wajib** diisi layar (A133; aplikasi tidak menulis `UNKNOWN`, K-012) |
| `IS_ADJUSTABLE_FLAG` | teks | VARCHAR2(10) | Adjustable `true`/`false` |
| `ITEM_TYPE` | teks | VARCHAR2(50) | Object Item Type — `V_JN_OBJ_ITEM.JN_OBJ_ITEM` |
| `ITEM_TYPE_ID` | teks | VARCHAR2(50) | `V_JN_OBJ_ITEM.MJOI_KODE` (nilai dropdown) |
| `PCT_ADJUST2` | angka desimal | NUMBER(38,8) | Adjustment Pct. (dropdown, tidak Adjustable) — persen |
| `PCT_ADJUST_OTHER` | angka desimal | NUMBER(38,8) | Adjustment Pct. (Adjustable: 60..100, `ValidateAdjustPct`) — persen |
| `PROPERTI_ITEM_NOTE` | teks | VARCHAR2(500) | Object Item Note — `V_JN_OBJ_ITEM.KETERANGAN` |
| `PROPERTY_ITEM_NO` | teks | VARCHAR2(50) | nomor item = `SEQ_NO` (K-6) |
| `REMARK` | teks | VARCHAR2(500) | Remark of `<ItemType>` |
| `TSI_OBJECT_ITEM` | angka desimal | NUMBER(38,8) | **UANG** — TSI Object Item (All Unit), ≥ 0, ≤ 8 desimal |
| `PROPERTY_YEAR` | teks | VARCHAR2(50) | **baru** — Year (`.PropertyYear`) |
| `UNIT` | teks | VARCHAR2(50) | **baru** — Unit(s) (`.Unit`), bilangan bulat > 0 |
| `CONDITION` | teks | VARCHAR2(500) | **baru** — Condition (`.Condition`), kode apa adanya |
| `YEAR` | teks | VARCHAR2(50) | **baru** — Year of Planting (`.Year`) |
| `NO_OF_TREE` | teks | VARCHAR2(50) | **baru** — No of Trees (`.NoOfTree`) |
| `AREA_HECTAR` | teks | VARCHAR2(50) | **baru** — Area ( Hectar ) (`.AreaHectar`) |
| `TOTAL_GROSS_PREMI` | angka desimal | NUMBER(38,8) | tiket 43 (193) — **uang**, Σ Premium coverage item, dihitung server saat PUT (`CountPremi_ACT` langkah 55-56); item tanpa coverage = kosong |
| `TOTAL_NET_RATE` | angka desimal | NUMBER(38,8) | tiket 43 (193) — ‰ Total Net Rate dari PUT; tiket 44: item ber-coverage tanpa kelima coverage net rate → 0 (`CekNetRate_ACT` langkah 4); item tanpa coverage → apa adanya (A166) |

## T_COVERAGELIST

Tiket 43 — migrasi 193, **sebagian** (pola A109). Rancangan memberi tabel ini **induk jamak** (CargoList / PropertyItemList /
AnekaList / PersonList / VehicleList), dibedakan `PARENT_TABLE` / `SRC_PATH`; NB menulis hanya jalur
`LocationList/Property/PropertyItemList/CoverageList` (`PARENT_TABLE` `T_PROPERTYITEMLIST`). `PARENT_ID` **tanpa FK**
(pola A140), indeks `IX_T_COVERAGELIST_PARENT (PARENT_TABLE, PARENT_ID)`. Banyak baris per item, urut `SEQ_NO`. Seluruh
medan kontrak ADA di rancangan — tanpa kolom baru. Medan uang / rate / persen dihitung server (`CountPremi_ACT`, basis
1–4) dan disimpan setengah-ke-atas 8 desimal (A155).

| Kolom | Tipe | Tipe DDL (migrasi 193) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_COVERAGELIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTYITEMLIST.ID` |
| `PARENT_TABLE` | teks | VARCHAR2(30) | `T_PROPERTYITEMLIST` |
| `SRC_PATH` | teks | VARCHAR2(200) | `LocationList/Property/PropertyItemList/CoverageList` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan coverage 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `COVERAGE` | teks | VARCHAR2(50) | `.Coverage` — ID coverage (`COVERAGE.ID`) |
| `OLDID` | teks | VARCHAR2(50) | `.OLDID` — kode tampil |
| `COVERAGE_NOTE` | teks | VARCHAR2(4000) | `.CoverageNote` — nama coverage (= COVERAGE.NAME 4000; 500 rancangan, dilebarkan migrasi 251) |
| `COVERAGE_BASIS` | teks | VARCHAR2(50) | `.CoverageBasis` — 1..4 (5 Layering ditolak 400, tahap C2–C4) |
| `DAY` | teks | VARCHAR2(50) | `.Day` |
| `TSI` | angka desimal | NUMBER(38,8) | **uang** — server: = `TSI_OBJECT_ITEM` item |
| `INDEMNITY` | teks | VARCHAR2(50) | `.Indemnity` |
| `RATE` | angka desimal | NUMBER(38,8) | ‰ Gross Rate (dihitung balik pada mode amount) |
| `RATE_OJK` | angka desimal | NUMBER(38,8) | ‰ Standard Rate (lookup = tahap C4) |
| `FIRST_LOSS` | teks | VARCHAR2(50) | % First Loss — teks desimal (rancangan VARCHAR2) |
| `DISCOUNT_PERCENTAGE` | angka desimal | NUMBER(38,8) | % Discount |
| `TSI_LIABILITY` | angka desimal | NUMBER(38,8) | **uang** — server, per basis |
| `NET_RATE` | angka desimal | NUMBER(38,8) | ‰ Net Rate (terisi → menggantikan Rate dalam rumus); tiket 44: item ber-flag + Total Net Rate → Rate / ΣRate × Total Net Rate (`CalculateNetRate_ACT`), disimpan 8 desimal, premi dari nilai 20 desimal |
| `LIMITOF_LIABILITY` | angka desimal | NUMBER(38,8) | **uang** — Limit of Liability |
| `PCT_LO_L` | angka desimal | NUMBER(38,8) | % LoL |
| `PRO_RATE_PERCENT` | angka desimal | NUMBER(38,8) | server — Prorate periode polis case |
| `INDEMNITY_PERCENTAGE` | angka desimal | NUMBER(38,8) | % Indemnity |
| `FIRST_SCALE` | teks | VARCHAR2(50) | % First Scale — teks desimal |
| `SUBLIMIT` | angka desimal | NUMBER(38,8) | % Sub Limit |
| `LOST_LIMIT` | teks | VARCHAR2(50) | % Loss Limit — teks desimal (ejaan rancangan / Pega `.LostLimit`) |
| `EML_PML` | teks | VARCHAR2(50) | % EML / PML — teks desimal |
| `DISCOUNT` | angka desimal | NUMBER(38,8) | **uang** — Discount |
| `PREMIUM` | angka desimal | NUMBER(38,8) | **uang** — Gross Premium (server pada mode percent) |
| `CONDITIONS` | teks | VARCHAR2(500) | `.Conditions` |
| `PCT_ADJUSTMENT` | angka desimal | NUMBER(38,8) | server — `CountPremi_ACT` langkah 11-15 (A156); tidak dikirim JSON |
| `CURRENCY_CODE` | teks | VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL | = `CURRENCY` item (K-069 / K-012); item ber-coverage dengan mata uang > 10 byte → 400 |
| `UNIT` | teks | VARCHAR2(50) | tiket 46 (migrasi 195 `ALTER ADD`) — `.Unit` Indemnity Unit, kode apa adanya |
| `ACCUMULATION_CODE` | teks | VARCHAR2(50) | tiket 46 (195) — `.AccumulationCode` (Choose / Copy Accumulation) |
| `ACCUMULATION_DESCRIPTION` | teks | VARCHAR2(500) | tiket 46 (195) — `.AccumulationDescription` = Accumulation Address |
| `ZONE` | teks | VARCHAR2(4000) | tiket 52 (203 VARCHAR2(1000), dilebarkan 204) — `.Zone` = `ZONES.DESCRIPTION` terpilih; rancangan VARCHAR2(50); 203 memakai lebar `ZONEOJK.DESCRIPTION` (sampel `DDL\CONTOH` sampai 111 byte), tetapi Zone di DEV > 1000 byte → 4000 (perintah work owner 06-10-2026) |
| `CONSTRUCTION_EARTHQUAKE` | teks | VARCHAR2(500) | tiket 52 (203) — `.ConstructionEarthquake` 4.2 Construction = `CONSTRUCTION.CONSTRUCTION`; TIDAK di rancangan (amandemen loader); lebar = kolom konstruksi saudaranya, lebar sumber belum terverifikasi |
| `TSI_NUSANTARA_RE` | angka desimal | NUMBER(38,8) | tiket 48 (198) — `.TSINusantaraRe` = PercentShare × TSILiability / 100 (CountPremiAndTSIRNMFireMBU_ACT langkah 6) |
| `PREMI_NUSANTARA_RE` | angka desimal | NUMBER(38,8) | tiket 48 (198) — `.PremiNusantaraRe` = PercentShare × Premium / 100 − PercentShare × Discount / 100 |
| `CALCULATE_METHOD` | teks | VARCHAR2(50) | tiket 73 P2 (214) — `.CalculateMethod` coverage Aneka |
| `CALCULATE_METHOD_FAC_IN` | teks | VARCHAR2(50) | tiket 73 P2 (214) — `.CalculateMethod_FacIn` |
| `PCT_SHORT_PERIOD` | angka desimal | NUMBER(38,8) | tiket 73 P2 (214) — `.PctShortPeriod` (K73-6; rancangan VARCHAR2(30)) |
| `LOADING` | angka desimal | NUMBER(38,8) | tiket 73 P2 (214) — `.Loading` |
| `SUB_LIMIT_NOTE` | teks | VARCHAR2(500) | tiket 73 P2 (214) — `.SubLimitNote` |
| `MIN_PREMIUM` | angka desimal | NUMBER(38,8) | tiket 73 P2 (214) — **uang**, `.MinPremium` |
| `TYPE_OF_DISCOUNT` | teks | VARCHAR2(10) | tiket 73 P2 (214) — `.TypeOfDiscount` teks apa adanya (sampel "true" 95/95; kolomTeksMenyimpang butir 68.1) |
| `IS_OLD_DATA` | teks | VARCHAR2(10) | tiket 73 P2 (214) — `.IsOldData` |
| `COVERAGE_INITIAL` | teks | VARCHAR2(500) | tiket 74 P2 (256 `ALTER ADD`) — `.CoverageInitial` coverage cargo |

## T_ANEKALIST

Tiket 73 P1 — migrasi 214. Objek keluarga Aneka (form `NB FacIn\Section\InputDtlObjectAneka_FacIn.xml`, kontainer C1–C9). Rancangan 27 kolom (tanpa penunjuk IDX_* milik loader) + kolom form **baru** (amandemen loader `amandemenAneka`). Induk **jamak**: okupasi Aneka (`T_OCCUPATIONLIST` ber-`PARENT_TABLE` `T_RISKLOCATION`) atau Golf (`T_RISKLOCATION`); `PARENT_ID` tanpa FK, indeks `IX_T_ANEKALIST_PARENT (PARENT_TABLE, PARENT_ID)`.

| Kolom | Tipe | Tipe DDL (migrasi 214) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_ANEKALIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_OCCUPATIONLIST.ID` (non-Golf) atau `T_RISKLOCATION.ID` (Golf); tanpa FK (induk jamak) |
| `PARENT_TABLE` | teks | VARCHAR2(30) | `T_OCCUPATIONLIST` / `T_RISKLOCATION` |
| `SRC_PATH` | teks | VARCHAR2(200) | `LocationList/Property/RiskLocation/OccupationList/AnekaList` / `LocationList/Property/RiskLocation/AnekaList` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan objek 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `BUILT_IN` | teks | VARCHAR2(50) | `.BuiltIn` Place of Manufacture (C2) |
| `DRAFT_WORDING_ID` | teks | VARCHAR2(50) | `.DraftWordingID` (C1, `IsBonding`) |
| `END_DATE` | teks | VARCHAR2(30) | `.END_DATE` (C1) |
| `FLAG` | teks | VARCHAR2(1000) | `.Flag` (C2; SHIP.FLAGNAME, lebar sumber) |
| `GEOGRAPHICAL_LIMITS` | teks | VARCHAR2(50) | `.GeographicalLimits` (C7) |
| `IS_OLD_DATA` | teks | VARCHAR2(10) | `.IsOldData` |
| `LOCATION_ID` | teks | VARCHAR2(50) | `.LocationID` (turunan rancangan) |
| `OBJECT_NAME` | teks | VARCHAR2(4000) | `.ObjectName`; HE = SetObjectDescHE_DT; lebar 4000 (SHIP.NAME 1000 + sambungan HE) |
| `OCCUPATION_ID` | teks | VARCHAR2(1000) | `.OccupationID`; lebar = 189 |
| `QUANTITY` | teks | VARCHAR2(50) | `.Quantity` (C3 / C5 / C9) |
| `START_DATE` | teks | VARCHAR2(30) | `.START_DATE` (C1) |
| `SECTION` | teks | VARCHAR2(50) | `.Section` 1..4 (`DDL\Section.xml`) |
| `SELECTED_LOCATION_ADDRESS` | teks | VARCHAR2(50) | `.SelectedLocationAddress` (rancangan) |
| `SELECTED_OBJECT_ITEM` | teks | VARCHAR2(50) | `.SelectedObjectItem` (rancangan) |
| `SERIAL_NO` | teks | VARCHAR2(50) | `.SerialNo` (C4) |
| `TSI` | angka desimal | NUMBER(38,8) | **uang** — `.TSI` |
| `UNIT_QUANTITY` | teks | VARCHAR2(50) | `.UnitQuantity` (C5) |
| `YEAR` | teks | VARCHAR2(50) | `.Year` (C2 / C5 / C7 / C8 / C9) |
| `CURRENCY_CODE` | teks | VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL | = `.Currency.Name` objek (K-069 / K-012) |
| `NO_KONTRAK` | teks | VARCHAR2(50) | **baru** — `.NO_KONTRAK` (C1) |
| `NO_KTP` | teks | VARCHAR2(50) | **baru** — `.NO_KTP` (C1) |
| `NO_NPWP` | teks | VARCHAR2(50) | **baru** — `.NO_NPWP` (C1) |
| `OTHERS` | teks | VARCHAR2(500) | **baru** — `.Others` (C1, blok Briguna luar) |
| `INVOICE_NO` | teks | VARCHAR2(50) | **baru** — `.InvoiceNo` (C3) |
| `MAINTENANCE_PERIOD` | angka desimal | NUMBER(38,8) | **baru** — `.MaintenancePeriod` (C1 / C4, `IsEngineering`; K73-6) |
| `NO_OF_TREE` | angka desimal | NUMBER(38,8) | **baru** — `.NoOfTree` (C8; K73-6) |
| `AREA_HECTAR` | angka desimal | NUMBER(38,8) | **baru** — `.AreaHectar` (C8; K73-6) |
| `LIMIT_IDEMNITY_DAY` | angka desimal | NUMBER(38,8) | **baru** — `.LimitIdemnityDay` (C9; K73-6) |
| `LIMIT_IDEMNITY_MONTH` | angka desimal | NUMBER(38,8) | **baru** — `.LimitIdemnityMonth` (C9; K73-6) |
| `LIMIT_OF_LIABILITY` | angka desimal | NUMBER(38,8) | **baru**, **uang** — `.LimitOfLiability` (grid okupasi) |
| `TSI_OLD` | angka desimal | NUMBER(38,8) | **baru**, **uang** — `.TSIOld` (grid Golf); ditulis aplikasi (loader V-28b) |
| `FLAG_DELETE` | teks | VARCHAR2(10) | **baru** — `.FlagDelete` (ADR-0031) |
| `ALM_RISK_ID` | teks | VARCHAR2(100) | **baru** — `.AlmRiskID` (lookup kapal = SHIP.ID) |
| `LOAN_ACCOUNT_NUMBER` | teks | VARCHAR2(50) | **baru** — Briguna `.LoanAccountNumber` (K73-7) |
| `LOAN_NAME` | teks | VARCHAR2(500) | **baru** — Briguna `.LoanName` |
| `JOB` | teks | VARCHAR2(500) | **baru** — Briguna `.Job` |
| `DEPOSIT_ACCOUNT_NO` | teks | VARCHAR2(50) | **baru** — Briguna `.DepositAccountNo` |
| `ID_CARD_NO` | teks | VARCHAR2(50) | **baru** — Briguna `.IDCardNo` |
| `INSTITUTION` | teks | VARCHAR2(500) | **baru** — Briguna `.Institution` |
| `LOAN_TYPE` | teks | VARCHAR2(50) | **baru** — Briguna `.LoanType` |
| `DATE_OF_BIRTH` | teks | VARCHAR2(30) | **baru** — Briguna `.DateOfBirth` (teks Pega) |
| `SALARY` | angka desimal | NUMBER(38,8) | **baru**, **uang** — Briguna `.Salary` |

## T_MARINEHULL

Tiket 73 P1 — migrasi 214, satu halaman per objek (`UQ_T_MARINEHULL_PARENT`). `.MarineHull` objek C2 (lookup kapal men-SET CLASSIFICATION / CONSTRUCTION / DWTRT / GRTRT / NRTRT; lebar = kolom POOLDATA.SHIP); LENGTHM / BREADTHM / DEPTHM **baru**.

| Kolom | Tipe | Tipe DDL (migrasi 214) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_MARINEHULL` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_ANEKALIST.ID` (UNIQUE, FK) |
| `CLASSIFICATION` | teks | VARCHAR2(1000) | `.CLASSIFICATION` |
| `CONSTRUCTION` | teks | VARCHAR2(1000) | `.CONSTRUCTION` |
| `DWTRT` | teks | VARCHAR2(1000) | `.DWTRT` |
| `GRTRT` | teks | VARCHAR2(1000) | `.GRTRT` |
| `NRTRT` | teks | VARCHAR2(1000) | `.NRTRT` |
| `TYPE_OF_VESSEL` | teks | VARCHAR2(50) | `.TYPE_OF_VESSEL` |
| `LENGTHM` | teks | VARCHAR2(50) | `.LENGTHM` |
| `BREADTHM` | teks | VARCHAR2(50) | `.BREADTHM` |
| `DEPTHM` | teks | VARCHAR2(50) | `.DEPTHM` |

## T_AVIATIONHULL

Tiket 73 P1 — migrasi 214, satu halaman per objek (`UQ_T_AVIATIONHULL_PARENT`). `.AviationHull` objek C7 — rancangan apa adanya.

| Kolom | Tipe | Tipe DDL (migrasi 214) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_AVIATIONHULL` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_ANEKALIST.ID` (UNIQUE, FK) |
| `MAX_PASSENGERS` | teks | VARCHAR2(50) | `.MAX_PASSENGERS` |
| `PILOTS` | teks | VARCHAR2(50) | `.PILOTS` |
| `REGISTRATION_MARKS` | teks | VARCHAR2(50) | `.REGISTRATION_MARKS` |
| `SPECIAL_RENTAL_USES` | teks | VARCHAR2(50) | `.SPECIAL_RENTAL_USES` |
| `SPECIAL_USES` | teks | VARCHAR2(50) | `.SPECIAL_USES` |
| `STANDARD_USES` | teks | VARCHAR2(50) | `.STANDARD_USES` |

## T_VEHICLEHE

Tiket 73 P1 — migrasi 214, satu halaman per objek (`UQ_T_VEHICLEHE_PARENT`). `.VehicleHE` objek C4 / C6 — BRAND / TYPE = ID terpilih lookup merk / tipe; CHASSIS_NUMBER / ENGINE_NUMBER / OBJECT_ID_HE **baru**.

| Kolom | Tipe | Tipe DDL (migrasi 214) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_VEHICLEHE` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_ANEKALIST.ID` (UNIQUE, FK) |
| `BRAND` | teks | VARCHAR2(50) | `.BRAND` |
| `BRAND_NAME` | teks | VARCHAR2(500) | `.BRAND_NAME` |
| `OBJECT_NAME_HE` | teks | VARCHAR2(500) | `.OBJECT_NAME_HE` |
| `TYPE` | teks | VARCHAR2(50) | `.TYPE` |
| `TYPE_NAME` | teks | VARCHAR2(500) | `.TYPE_NAME` |
| `CHASSIS_NUMBER` | teks | VARCHAR2(50) | `.CHASSIS_NUMBER` |
| `ENGINE_NUMBER` | teks | VARCHAR2(50) | `.ENGINE_NUMBER` |
| `OBJECT_ID_HE` | teks | VARCHAR2(50) | `.OBJECT_ID_HE` |

## T_MEMORANDUMACCEPT

Tiket 73 P1 — migrasi 214, satu halaman per objek (`UQ_T_MEMORANDUMACCEPT_PARENT`). `.MemorandumAccept` objek C1 — **tabel baru** (tidak di rancangan): SURETY_TYPE diisi server dari BusinessCode (`InputDtlObjectAneka_PreAct`), 15 medan `Section\InputDtlMemorandumAccept.xml`.

| Kolom | Tipe | Tipe DDL (migrasi 214) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_MEMORANDUMACCEPT` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_ANEKALIST.ID` (UNIQUE, FK) |
| `SURETY_TYPE` | teks | VARCHAR2(50) | `.SURETY_TYPE` |
| `APPLICATION_NUMBER` | teks | VARCHAR2(500) | `.APPLICATION_NUMBER` |
| `ACCUMULATION` | teks | VARCHAR2(500) | `.ACCUMULATION` |
| `COLLATERAL` | teks | VARCHAR2(500) | `.COLLATERAL` |
| `INSURED_SIUP` | teks | VARCHAR2(500) | `.INSURED_SIUP` |
| `KEY_PERSON` | teks | VARCHAR2(500) | `.KEY_PERSON` |
| `ORIGINAL_RATE` | teks | VARCHAR2(500) | `.ORIGINAL_RATE` |
| `PRINCIPALS_OUT_GO` | teks | VARCHAR2(500) | `.PRINCIPALS_OUT_GO` |
| `TYPE_OF_WORDING` | teks | VARCHAR2(500) | `.TYPE_OF_WORDING` |
| `BUILDING_STATUS` | teks | VARCHAR2(4000) | `.BUILDING_STATUS` |
| `SUPPORTING_DOCUMENT` | teks | VARCHAR2(4000) | `.SUPPORTING_DOCUMENT` |
| `PRINCIPALS_DATA` | teks | VARCHAR2(4000) | `.PRINCIPALS_DATA` |
| `CONCLUSIONS_PROPOSALS` | teks | VARCHAR2(4000) | `.CONCLUSIONS_PROPOSALS` |
| `PRINCIPALS_BACKGROUND` | teks | VARCHAR2(4000) | `.PRINCIPALS_BACKGROUND` |
| `SUMMARY_FINANCIAL` | teks | VARCHAR2(4000) | `.SUMMARY_FINANCIAL` |
| `PRINCIPALS_GUARANTEE` | teks | VARCHAR2(4000) | `.PRINCIPALS_GUARANTEE` |

## T_SPREADINGLIST

Tiket 48 (tab Spreading FIRE) — **dibuat sebagian** modul ini, migrasi `198_spreading.sql` (+ `SEQ_T_SPREADINGLIST`, indeks
`IX_T_SPREADINGLIST_PARENT`). Tabel RANCANGAN (jalur `LocationList/Property/PropertyItemList/CoverageList/SpreadingList`,
induk `T_COVERAGELIST`, kelas Data-SpreadingRisk); kolom = yang ditulis tab Spreading NB FIRE. Diganti utuh tiap Save tab
Spreading; dipertahankan menurut posisi saat Save tab Object (K48-7). Total spreading TIDAK disimpan (K-067).

| Kolom | Tipe | Tipe DDL (migrasi 198) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SPREADINGLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_COVERAGELIST.ID` (FK tanpa `ON DELETE`) |
| `PARENT_TABLE` | teks | VARCHAR2(30) | `T_COVERAGELIST` (pola loader) |
| `SRC_PATH` | teks | VARCHAR2(200) | jalur rancangan di atas |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan baris di coverage, mulai 1 |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID v4 acak (A106) |
| `TREATY_TYPE` | teks | VARCHAR2(50) | `.TreatyType` — ID jenis reasuransi (REINSURANCETYPE / PROPORTIONALARRG.REINSTYPEID; "10007" ORS, "10015" FACOUT) |
| `TREATY_NAME` | teks | VARCHAR2(500) | `.TreatyName` |
| `SHARE_PERCENTAGE` | angka desimal | NUMBER(38,8) | `.SharePercentage` (%) |
| `TSI_SPREADED` | angka desimal | NUMBER(38,8) | `.TSISpreaded` = Share × TSINusantaraRe / 100 (DT SetTSIPremiSpreaded_FacIn "percent") |
| `TSI_GROSS_SPREADED` | angka desimal | NUMBER(38,8) | `.TSIGrossSpreaded` = Share × TSI × PercentShare / 10000 |
| `PREMIUM_SPREADED` | angka desimal | NUMBER(38,8) | `.PremiumSpreaded` = Share × PremiNusantaraRe / 100 |
| `CLAIM_ESTIMATION` | angka desimal | NUMBER(38,8) | `.ClaimEstimation` = LimitofLiability × PercentShare / 100 × Share / 100 (SumTSIPremiSpreadedRNM_FIRE_Act) |
| `CURRENCY_CODE` | teks | VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL | = `CURRENCY` item (pola T_COVERAGELIST) |

## T_OCCUPATIONLIST

Tiket 40 — migrasi 189, **sebagian** (pola A109; tanpa `CATEGORY` okupasi, `IS_OLD_DATA`, `LOCATION_ID` — tiket 23).
Rancangan memberi tabel ini **tiga induk** (`LocationList/Property/OccupationList`, `.../RiskLocation/OccupationList`,
`VehicleList/OccupationList`); NB menulis jalur pertama saja: `PARENT_TABLE = 'T_PROPERTY'`, `SRC_PATH` = jalur itu (isi sama
dengan loader). Karena berinduk jamak, `PARENT_ID` **tanpa FK** — indeks `IX_T_OCCUPATIONLIST_PARENT (PARENT_TABLE, PARENT_ID)`.

| Kolom | Tipe | Tipe DDL (migrasi 189) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_OCCUPATIONLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` (bila `PARENT_TABLE = 'T_PROPERTY'`) |
| `PARENT_TABLE` | teks | VARCHAR2(30) | tabel induk — `T_PROPERTY` |
| `SRC_PATH` | teks | VARCHAR2(200) | jalur sumber — `LocationList/Property/OccupationList` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan okupasi 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `OCCUPATION_ID` | teks | VARCHAR2(1000) | Occupation ID (`.OccupationId` = `OCCUPATION.OLDID`); lebar = sumber (A138, rancangan 50) |
| `OCCUPATION_NAME` | teks | VARCHAR2(1000) | Occupation Name (`OCCUPATION.NAME`); lebar = sumber (A138, rancangan 500) |

## T_TABLEOFLIMIT

Tiket 40 — migrasi 189, **utuh**. Halaman `.TableOfLimit` okupasi: satu baris per okupasi (`UQ_T_TABLEOFLIMIT_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 189) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_TABLEOFLIMIT` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_OCCUPATIONLIST.ID` |
| `CATEGORY` | teks | VARCHAR2(50) | Category I/II/III (`SetDataOccupation` dari KDRiskExposure 01/02/03) |
| `DESCRIPTION` | teks | VARCHAR2(500) | Class of Construction |
| `PCT_LIMIT` | teks | VARCHAR2(50) | PctLimit — **teks apa adanya** (butir 68.1: rancangan NUMBER; koma desimal dan spasi ujung dipertahankan); lebar A139 |

## T_FEALIST

Tiket 41 — migrasi 190, tabel **BARU** (rancangan flat tidak punya tabel FEA; A142, amandemen loader `amandemenFEA` +
jalur `LocationList/FEAList` + lipatan `.DataFEA`). `.FEAList` milik **baris lokasi** (kelas
`Data-OfferFacIn-LocationReinsurance`), jadi induknya `T_LOCATIONLIST`; banyak baris per lokasi, urut `SEQ_NO`. Kolom
sistem = pola tabel berulang rancangan (T_ADDITIONALSHIP).

| Kolom | Tipe | Tipe DDL (migrasi 190) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_FEALIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_LOCATIONLIST.ID` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan FEA 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `APAR` | teks | VARCHAR2(50) | APAR (`.APAR`) — jumlah unit, bilangan bulat ≥ 0 (M-2) |
| `SPRINKLER` | teks | VARCHAR2(50) | Sprinkler (`.Sprinkler`) — jumlah unit |
| `SMOKE_DETECTOR` | teks | VARCHAR2(50) | Smoke Detector & Alarm (`.SmokeDetector`) — jumlah unit |
| `HYDRANT` | teks | VARCHAR2(50) | Hydrant (`.Hydrant`) — jumlah unit |
| `PRIVATE_TRUCK_BRIGADE` | teks | VARCHAR2(50) | Private Truck Brigade (Unit) (`.DataFEA.PrivateTruckBrigade`) — jumlah unit |
| `PRIVATE_FIRE_BRIGADE` | teks | VARCHAR2(50) | Private Team Fire Brigade (`.DataFEA.PrivateFireBrigade`) — kode apa adanya |
| `TEAM_SOP_SAFETY` | teks | VARCHAR2(50) | Team & SOP Safety (`.DataFEA.TeamSOPSafety`) — kode apa adanya |
| `TEAM_SOP_RISK_MANAGEMENT` | teks | VARCHAR2(50) | Team & SOP Risk Management (`.DataFEA.TeamSOPRiskManagement`) — kode apa adanya |
| `INFO_FEA` | teks | VARCHAR2(500) | Others Info (`.InfoFEA`) |

## T_LISTCAUSEOFLOSS

Tiket 42 — migrasi 191, **utuh** (rancangan jalur `LocationList/Property/ListCauseOfLoss`, induk `T_PROPERTY`) + lima kolom
**baru** (A145, amandemen loader `amandemenKerugian`). Banyak baris per property, urut `SEQ_NO`.

| Kolom | Tipe | Tipe DDL (migrasi 191) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_LISTCAUSEOFLOSS` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_PROPERTY.ID` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan catatan 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `CLAIM` | angka desimal | NUMBER(38,8) | Total Claim (100%) — **uang** |
| `CURRENCY` | teks | VARCHAR2(50) DEFAULT 'UNKNOWN' (NOT NULL di 191, boleh NULL sejak migrasi 209) | Currency (`CURRENCY`, ≠ ITL); **tidak wajib** sejak tiket 42b (keputusan work owner 07-10-2026, Pega pyRequired=false; meralat A133) - kosong = NULL, aplikasi tidak menulis 'UNKNOWN' (K-012) |
| `DETAIL` | teks | VARCHAR2(500) | Loss Detail — dilebarkan dari rancangan 50 (A146) |
| `REMARKS` | teks | VARCHAR2(500) | Remarks — apa adanya (`DDL\Remarks.xml`: Settled / Ex Gratia Payment / Withdraw / Others / --, dipakai frontend) |
| `DATE_OF_LOSS` | teks | VARCHAR2(30) | **baru** — Date of Loss, teks Pega `YYYYMMDDTHHMMSS.mmm GMT` pukul 12:00 WIB (pola Begin date) |
| `LOSS_OBJECT` | teks | VARCHAR2(500) | **baru** — Loss Object |
| `AMOUNT` | angka desimal | NUMBER(38,8) | **baru** — Total of Loss (baca-saja di layar, N-2) — **uang** |
| `PREVENTION_OF_LOSS` | angka desimal | NUMBER(38,8) | **baru** — Prevention Of Loss — **uang** |
| `CAUSE_OF_LOSS` | teks | VARCHAR2(500) | **baru** — Cause of Loss |

## T_COINSDATA

Tiket 42 — migrasi 191, **utuh** (rancangan `.../ListCauseOfLoss/CoinsData`, satu halaman per catatan:
`UQ_T_COINSDATA_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 191) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_COINSDATA` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_LISTCAUSEOFLOSS.ID` |
| `COINS_NAME` | teks | VARCHAR2(500) | Insured Name grid — **selalu** nama tertanggung case (`SetLossRatio_Act` langkah 3.1) |

## T_DEDUCTIBLELIST

Tiket 45 — migrasi 194, **sebagian** (pola A109): rancangan jalur `.../CoverageList/DeductibleList` (lima jalur, semuanya
berinduk `T_COVERAGELIST`) → `PARENT_ID` ber-FK ke `T_COVERAGELIST.ID`, indeks `IX_T_DEDUCTIBLELIST_PARENT`. Banyak baris
per coverage, urut `SEQ_NO`. Tanpa kolom yang belum dipakai layar (`FLAG_CURRENCY` R-3, `DESCRIPTIONS` - ditambah 256 tiket 74, `BASIS_TYPE`,
`CLAIM_CATEGORY`, `DEDUCTIBLE_TYPE`, `TYPE`, `VALUE`, `COVERAGE`, kolom turunan `*_ID`).

| Kolom | Tipe | Tipe DDL (migrasi 194) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_DEDUCTIBLELIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_COVERAGELIST.ID` (FK) |
| `PARENT_TABLE` | teks | VARCHAR2(30) | `T_COVERAGELIST` |
| `SRC_PATH` | teks | VARCHAR2(200) | `LocationList/Property/PropertyItemList/CoverageList/DeductibleList` |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan deductible 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `AMOUNT` | angka desimal | NUMBER(38,8) | **uang** — Amount |
| `CONDITION` | teks | VARCHAR2(500) | kode Condition 1–5 (`DDL\Condition.xml` Deductible) |
| `CURRENCY` | teks | VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL | Currency (`CURRENCY`, ≠ ITL); boleh kosong → tidak disisipkan, DEFAULT mengisi, dibaca kosong (A169) |
| `INPUT_CONDITION` | teks | VARCHAR2(500) | teks kondisi bila Condition 5 Other |
| `MIN_MAX` | teks | VARCHAR2(50) | kode 1 Min / 2 Max / 3 Or (`DDL\MinMax.xml`) |
| `PCT_DEDUCTIBLE` | angka desimal | NUMBER(38,8) | persen |
| `PCT_DEDUCTIBLE2` | angka desimal | NUMBER(38,8) | persen (bila MinMax 3) |
| `TIME_EXCESS` | teks | VARCHAR2(30) | Time Excess (hari) — teks desimal (A171) |
| `TYPE_DEDUCTIBLE` | angka bulat | NUMBER(5) | kode pyStandardValue 0–7 (`DDL\TypeDeductible.xml`, A170) |
| `TYPE_DEDUCTIBLE2` | angka bulat | NUMBER(5) | kode (`DDL\TypeDeductible2.xml`, bila MinMax 3) |
| `DESCRIPTIONS` | teks | VARCHAR2(4000) | tiket 74 P2 (256 `ALTER ADD`) — `.Descriptions` deskripsi deductible coverage cargo (pxTextArea) |

## T_CARGOLIST

Tiket 74 P1 — migrasi 256: grid `.OfferFacIn.CargoList` MARINE CARGO (rancangan 22 kolom + 4 baru). Satu baris per cargo, urut `SEQ_NO`.

| Kolom | Tipe | Tipe DDL (migrasi 256) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_CARGOLIST` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | teks | VARCHAR2(32) NOT NULL | = `T_GENERAL_POLIS.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan baris 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `BL_NUMBER` | teks | VARCHAR2(4000) | `.PolicyData.BLNumber` "No BL" (pxTextArea; rancangan VARCHAR2(50)) |
| `CONVEYANCE_ID` | teks | VARCHAR2(50) | `.ConveyanceID` (BrowseConveyance_RD .ID) |
| `CONVEYANCE_NAME` | teks | VARCHAR2(500) | server — `.ConveyanceName` = Description master (ProtectObjectMarine) |
| `CONVEYANCE_NOTE` | teks | VARCHAR2(500) | server — `.ConveyanceNote` = Description master |
| `FLAG_OLD_DATA` | teks | VARCHAR2(10) | `.FlagOldData` rancangan |
| `FROM_RUTE` | teks | VARCHAR2(50) | `.FromRute` |
| `GOOD_ID` | teks | VARCHAR2(50) | `.GoodID` (BrowseGoodsType .ID) |
| `GOOD_NOTE` | teks | VARCHAR2(500) | server — `.GoodNote` = Description master |
| `INVOICE_NUMBER` | teks | VARCHAR2(50) | `.PolicyData.InvoiceNumber` |
| `IS_TOP_RISK` | teks | VARCHAR2(10) | `.IsTopRisk` teks "true" / "false" |
| `PACKING_ID` | teks | VARCHAR2(50) | `.PackingID` (BrowsePacking_RD .ID) |
| `PACKING_NOTE` | teks | VARCHAR2(500) | server — `.PackingNote` = Description master |
| `SAIL_DATE` | teks | VARCHAR2(30) | `.PolicyData.SailDate` teks tanggal Pega (pola START_DATE) |
| `TO_RUTE` | teks | VARCHAR2(50) | `.ToRute` |
| `TRADING_ID` | teks | VARCHAR2(50) | `.TradingID` (BrowseTrading_RD .ID) |
| `TRADING_NOTE` | teks | VARCHAR2(500) | server — `.TradingNote` = Description master |
| `INVOICE_DATE` | teks | VARCHAR2(30) | **baru** — `.PolicyData.InvoiceDate` (pola START_DATE) |
| `MAX_AGE` | teks | VARCHAR2(50) | **baru** — `.PolicyData.MaxAge` (BrowseMaxAge_RD .ID) |
| `MAX_AGE_DETAIL` | teks | VARCHAR2(4000) | **baru** — server, `.PolicyData.MaxAgeDetail` = `.Warranty` master |
| `FLAG_DELETE` | teks | VARCHAR2(10) | **baru** — `.FlagDelete` (hapus lunak EDM; NB tidak menulis) |

## T_SHIP

Tiket 74 P1 — migrasi 256: `.PolicyData.Ship` satu halaman per cargo (`UQ_T_SHIP_PARENT`), diisi popup Select Ship (`SelectShip_act` L470-873). Lebar = POOLDATA.SHIP.

| Kolom | Tipe | Tipe DDL (migrasi 256) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SHIP` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_CARGOLIST.ID` (UNIQUE, FK) |
| `AGE` | teks | VARCHAR2(50) | server — CountAgeShip: tahun SailDate − `.Y_MAKE1` |
| `CONST` | teks | VARCHAR2(1000) | `.CONST` (SelectShip Construction) |
| `DWT` | teks | VARCHAR2(1000) | `.DWT` teks, koma → titik (SelectShip) |
| `FLAG_NAME` | teks | VARCHAR2(1000) | `.FLAG_NAME` |
| `GRT` | teks | VARCHAR2(1000) | `.GRT` teks, koma → titik |
| `IMO` | teks | VARCHAR2(1000) | `.Imo` |
| `NM_SHIP` | teks | VARCHAR2(1000) | `.NM_SHIP` (SelectShip Name) |
| `NRT` | teks | VARCHAR2(1000) | `.NRT` teks, koma → titik |
| `REGISTER` | teks | VARCHAR2(1000) | `.REGISTER` |
| `REMARK` | teks | VARCHAR2(1000) | `.REMARK` |
| `SHIPTYPE_ID` | teks | VARCHAR2(1000) | `.SHIPTYPE_ID` |
| `SHIPTYPE_NAME` | teks | VARCHAR2(1000) | `.SHIPTYPE_NAME` (= SHIP.SHIPNAME, SelectShip.xml L6439) |
| `SHIP_REF_ID` | teks | VARCHAR2(100) | `.Ship.ID` (= SHIP.ID) |
| `Y_MAKE1` | teks | VARCHAR2(1000) | `.Y_MAKE1` |
| `Y_MAKE2` | teks | VARCHAR2(1000) | `.Y_MAKE2` |
| `FORMER_SHIP_NAME` | teks | VARCHAR2(1000) | **baru** — `.FormerShipName` |
| `CLASSIFICATION` | teks | VARCHAR2(1000) | **baru** — `.Classification` |
| `REBUILT` | teks | VARCHAR2(1000) | **baru** — `.Rebuilt` |

## T_LC

Tiket 74 P1 — migrasi 256: `.PolicyData.LC` satu halaman per cargo (`UQ_T_LC_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 256) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_LC` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_CARGOLIST.ID` (UNIQUE, FK) |
| `LC` | teks | VARCHAR2(50) | `.LC` cekbox "LC Status" teks "true" / "false" |
| `LCBANK` | teks | VARCHAR2(500) | **baru** — `.LCBANK` |
| `LC_NUMBER` | teks | VARCHAR2(500) | **baru** — `.LCNumber` |
| `LC_CONDITION` | teks | VARCHAR2(500) | **baru** — `.LCCondition` |
| `LC_REMARK` | teks | VARCHAR2(500) | **baru** — `.LCRemark` |
| `LC_END_DATE` | teks | VARCHAR2(30) | **baru** — `.LCEndDate` "LC Expired Date" (pola START_DATE) |

## T_SURVEYAGENT

Tiket 74 P1 — migrasi 256 (**baru**, K74-6 per cargo): `.PolicyData.SurveyAgent` satu halaman per cargo (`UQ_T_SURVEYAGENT_PARENT`).

| Kolom | Tipe | Tipe DDL (migrasi 256) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_SURVEYAGENT` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_CARGOLIST.ID` (UNIQUE, FK) |
| `ASM_AGENT_ID` | teks | VARCHAR2(100) | `.ASMAgentID` (BrowseLLOYDAGENT_RD .ID) |
| `FULL_NAME` | teks | VARCHAR2(4000) | `.pyFullName` (pxTextArea) |
| `ADDRESS` | teks | VARCHAR2(500) | `.pyAddress` |
| `CITY` | teks | VARCHAR2(500) | `.pyCity` |
| `COUNTRY` | teks | VARCHAR2(500) | `.pyCountry` |

## T_ADDITIONALSHIP

Tiket 74 P1 — migrasi 256 (**baru**): grid `.PolicyData.Ship.AdditionalShip`, banyak baris per kapal, urut `SEQ_NO`.

| Kolom | Tipe | Tipe DDL (migrasi 256) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_ADDITIONALSHIP` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_SHIP.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan baris 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `SHIP_REF_ID` | teks | VARCHAR2(100) | `.ID` (autocomplete Register ID) |
| `REGISTER` | teks | VARCHAR2(1000) | `.REGISTER` |
| `NM_SHIP` | teks | VARCHAR2(1000) | `.NM_SHIP` |
| `VOYAGE_NO` | teks | VARCHAR2(500) | `.VOYAGE_NO` |
| `SHIPTYPE_ID` | teks | VARCHAR2(1000) | `.SHIPTYPE_ID` |
| `FLAG_NAME` | teks | VARCHAR2(1000) | `.FLAG_NAME` |
| `NRT` | teks | VARCHAR2(1000) | `.NRT` |
| `DWT` | teks | VARCHAR2(1000) | `.DWT` |
| `GRT` | teks | VARCHAR2(1000) | `.GRT` |
| `Y_MAKE1` | teks | VARCHAR2(1000) | `.Y_MAKE1` |
| `CONST` | teks | VARCHAR2(1000) | `.CONST` |

## T_ADDITIONALCOVERAGE

Tiket 74 P2 — migrasi 256: grid Warranty Extension `.AdditionalCoverage` coverage cargo (sebagian rancangan), banyak baris per coverage.

| Kolom | Tipe | Tipe DDL (migrasi 256) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_ADDITIONALCOVERAGE` |
| `IDPEGA` | teks | VARCHAR2(50) | kolom sistem rancangan |
| `COB_GROUP` | teks | VARCHAR2(20) | kolom sistem rancangan |
| `PARENT_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_COVERAGELIST.ID` (FK) |
| `SEQ_NO` | angka bulat | NUMBER(5) NOT NULL | urutan baris 1..n |
| `ROW_UID` | teks | VARCHAR2(36) NOT NULL | UUID baris |
| `COVERAGE` | teks | VARCHAR2(50) | `.Coverage` (BrowseMarineCondition_RD .ID) |
| `COVERAGE_NOTE` | teks | VARCHAR2(500) | `.CoverageNote` (= ConditionName) |

## T_NB_CSV_UNGGAH

Tiket 75 — migrasi 257: satu pekerjaan unggah CSV objek tanpa batas (status, kemajuan, header, galat proses).

| Kolom | Tipe | Tipe DDL (migrasi 257) | Isi |
| --- | --- | --- | --- |
| `ID` | angka bulat | NUMBER(19) NOT NULL | surrogate dari `SEQ_T_NB_CSV_UNGGAH` |
| `CASE_ID` | teks | VARCHAR2(32) NOT NULL | ID case (T_GENERAL_POLIS.ID / T_WORK_POLIS.ID), tanpa FK |
| `STATUS` | teks | VARCHAR2(20) NOT NULL | memvalidasi / valid / gagal_validasi / menyimpan / selesai / gagal |
| `FASE` | teks | VARCHAR2(20) NOT NULL | validasi / simpan (Ulangi melanjutkan fase ini) |
| `NAMA_BERKAS` | teks | VARCHAR2(500) | nama berkas bagian multipart |
| `UKURAN` | angka bulat | NUMBER(19) | byte diterima |
| `JUMLAH_BARIS` | angka bulat | NUMBER(10) DEFAULT 0 NOT NULL | baris data berkas |
| `BARIS_DIPROSES` | angka bulat | NUMBER(10) DEFAULT 0 NOT NULL | baris tervalidasi |
| `BARIS_TERSIMPAN` | angka bulat | NUMBER(10) DEFAULT 0 NOT NULL | baris tersimpan (lokasi) |
| `KELOMPOK_TERAKHIR` | angka bulat | NUMBER(10) DEFAULT 0 NOT NULL | kelompok simpan terakhir yang COMMIT (1000 baris) |
| `JUMLAH_GALAT` | angka bulat | NUMBER(10) DEFAULT 0 NOT NULL | jumlah galat baris |
| `JUMLAH_PESAN` | angka bulat | NUMBER(10) DEFAULT 0 NOT NULL | jumlah pesan (KodePos …) |
| `NO_OFFER_SLIP` | teks | VARCHAR2(4000) | header langkah 5.3 (teks CSV, D1 (b)) |
| `BEGIN_DATE` | teks | VARCHAR2(4000) | header langkah 5.3 (teks CSV dd/MM/yyyy) |
| `END_DATE` | teks | VARCHAR2(4000) | header langkah 5.3 (teks CSV dd/MM/yyyy) |
| `GALAT_PROSES` | teks | VARCHAR2(4000) | sebab status gagal ("gagal di baris N: …") |
| `GALAT_DI_BARIS` | angka bulat | NUMBER(10) | baris awal kelompok yang gagal |
| `DIBUAT_OLEH` | teks | VARCHAR2(100) | akun pengunggah |
| `DIBUAT` | timestamp | TIMESTAMP NOT NULL | SYSTIMESTAMP saat dibuat |
| `DIUBAH` | timestamp | TIMESTAMP NOT NULL | detak (SYSTIMESTAMP tiap langkah; basi > 15 menit = terputus) |
| `JENIS` | teks | VARCHAR2(10) DEFAULT 'fire' NOT NULL | importer: fire / cargo (migrasi 258, tiket 74 P6; pekerjaan lama = fire) / aneka (tiket 73 P6, tanpa migrasi kolom ini) |
| `KOLOM` | teks | VARCHAR2(4000) | judul dikenali urut berkas, dipisah ";" (migrasi 258; NULL = urutan template FIRE) |

## T_NB_CSV_BARIS

Tiket 75 — migrasi 257: penampung baris CSV apa adanya; satu kolom `C_<kolom>` VARCHAR2(4000) per kolom template (61, urut KolomCsvObjek); kunci (UNGGAH_ID, NO_BARIS).

| Kolom | Tipe | Tipe DDL (migrasi 257) | Isi |
| --- | --- | --- | --- |
| `UNGGAH_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_NB_CSV_UNGGAH.ID` (FK) |
| `NO_BARIS` | angka bulat | NUMBER(10) NOT NULL | nomor baris data 1..n |
| `C_NO` | teks | VARCHAR2(4000) | kolom CSV `NO` apa adanya |
| `C_NO_OFFER_SLIP` | teks | VARCHAR2(4000) | kolom CSV `NO_OFFER_SLIP` apa adanya |
| `C_CLASS_OF_BUSINESS` | teks | VARCHAR2(4000) | kolom CSV `CLASS_OF_BUSINESS` apa adanya |
| `C_BEGIN_DATE` | teks | VARCHAR2(4000) | kolom CSV `BEGIN_DATE` apa adanya |
| `C_END_DATE` | teks | VARCHAR2(4000) | kolom CSV `END_DATE` apa adanya |
| `C_INSURED_NAME` | teks | VARCHAR2(4000) | kolom CSV `INSURED_NAME` apa adanya |
| `C_SOB_NAME` | teks | VARCHAR2(4000) | kolom CSV `SOB_NAME` apa adanya |
| `C_CEDING_CO_NAME` | teks | VARCHAR2(4000) | kolom CSV `CEDING_CO_NAME` apa adanya |
| `C_MARKETING_NAME` | teks | VARCHAR2(4000) | kolom CSV `MARKETING_NAME` apa adanya |
| `C_OCCUPATION_CODE` | teks | VARCHAR2(4000) | kolom CSV `OCCUPATION_CODE` apa adanya |
| `C_OCCUPATION_NAME` | teks | VARCHAR2(4000) | kolom CSV `OCCUPATION_NAME` apa adanya |
| `C_CLASS_OF_CONSTRUCTION` | teks | VARCHAR2(4000) | kolom CSV `CLASS_OF_CONSTRUCTION` apa adanya |
| `C_TYPE_ADDRESS` | teks | VARCHAR2(4000) | kolom CSV `TYPE_ADDRESS` apa adanya |
| `C_ADDRESS` | teks | VARCHAR2(4000) | kolom CSV `ADDRESS` apa adanya |
| `C_BUILDING_NO` | teks | VARCHAR2(4000) | kolom CSV `BUILDING_NO` apa adanya |
| `C_POSTAL_CODE` | teks | VARCHAR2(4000) | kolom CSV `POSTAL_CODE` apa adanya |
| `C_KELURAHAN` | teks | VARCHAR2(4000) | kolom CSV `KELURAHAN` apa adanya |
| `C_OBJECT_TYPE` | teks | VARCHAR2(4000) | kolom CSV `OBJECT_TYPE` apa adanya |
| `C_INTEREST` | teks | VARCHAR2(4000) | kolom CSV `INTEREST` apa adanya |
| `C_TOP_RISK` | teks | VARCHAR2(4000) | kolom CSV `TOP_RISK` apa adanya |
| `C_ITEM_TYPE` | teks | VARCHAR2(4000) | kolom CSV `ITEM_TYPE` apa adanya |
| `C_ITEM_TYPE_1` | teks | VARCHAR2(4000) | kolom CSV `ITEM_TYPE_1` apa adanya |
| `C_ITEM_TYPE_2` | teks | VARCHAR2(4000) | kolom CSV `ITEM_TYPE_2` apa adanya |
| `C_ITEM_TYPE_3` | teks | VARCHAR2(4000) | kolom CSV `ITEM_TYPE_3` apa adanya |
| `C_ITEM_TYPE_4` | teks | VARCHAR2(4000) | kolom CSV `ITEM_TYPE_4` apa adanya |
| `C_ITEM_TYPE_5` | teks | VARCHAR2(4000) | kolom CSV `ITEM_TYPE_5` apa adanya |
| `C_ITEM_TYPE_6` | teks | VARCHAR2(4000) | kolom CSV `ITEM_TYPE_6` apa adanya |
| `C_TSI_OBJECT_ITEM` | teks | VARCHAR2(4000) | kolom CSV `TSI_OBJECT_ITEM` apa adanya |
| `C_TSI_OBJECT_ITEM_1` | teks | VARCHAR2(4000) | kolom CSV `TSI_OBJECT_ITEM_1` apa adanya |
| `C_TSI_OBJECT_ITEM_2` | teks | VARCHAR2(4000) | kolom CSV `TSI_OBJECT_ITEM_2` apa adanya |
| `C_TSI_OBJECT_ITEM_3` | teks | VARCHAR2(4000) | kolom CSV `TSI_OBJECT_ITEM_3` apa adanya |
| `C_TSI_OBJECT_ITEM_4` | teks | VARCHAR2(4000) | kolom CSV `TSI_OBJECT_ITEM_4` apa adanya |
| `C_TSI_OBJECT_ITEM_5` | teks | VARCHAR2(4000) | kolom CSV `TSI_OBJECT_ITEM_5` apa adanya |
| `C_TSI_OBJECT_ITEM_6` | teks | VARCHAR2(4000) | kolom CSV `TSI_OBJECT_ITEM_6` apa adanya |
| `C_ACCUMULATION_CODE` | teks | VARCHAR2(4000) | kolom CSV `ACCUMULATION_CODE` apa adanya |
| `C_ACCUMULATION_NAME` | teks | VARCHAR2(4000) | kolom CSV `ACCUMULATION_NAME` apa adanya |
| `C_ADJUSTABLE_PCT` | teks | VARCHAR2(4000) | kolom CSV `ADJUSTABLE_PCT` apa adanya |
| `C_CURRENCY` | teks | VARCHAR2(4000) | kolom CSV `CURRENCY` apa adanya |
| `C_INDEMNITY_MONTH` | teks | VARCHAR2(4000) | kolom CSV `INDEMNITY_MONTH` apa adanya |
| `C_LOL` | teks | VARCHAR2(4000) | kolom CSV `LOL` apa adanya |
| `C_RATE_FLEXAS` | teks | VARCHAR2(4000) | kolom CSV `RATE_FLEXAS` apa adanya |
| `C_NET_RATE_FLEXAS` | teks | VARCHAR2(4000) | kolom CSV `NET_RATE_FLEXAS` apa adanya |
| `C_DISCOUNT_RATE_FLEXAS` | teks | VARCHAR2(4000) | kolom CSV `DISCOUNT_RATE_FLEXAS` apa adanya |
| `C_RATE_RSMDCC` | teks | VARCHAR2(4000) | kolom CSV `RATE_RSMDCC` apa adanya |
| `C_NET_RATE_RSMDCC` | teks | VARCHAR2(4000) | kolom CSV `NET_RATE_RSMDCC` apa adanya |
| `C_DISCOUNT_RATE_RSMDCC` | teks | VARCHAR2(4000) | kolom CSV `DISCOUNT_RATE_RSMDCC` apa adanya |
| `C_RATE_FWTWD` | teks | VARCHAR2(4000) | kolom CSV `RATE_FWTWD` apa adanya |
| `C_NET_RATE_FWTWD` | teks | VARCHAR2(4000) | kolom CSV `NET_RATE_FWTWD` apa adanya |
| `C_DISCOUNT_RATE_FWTWD` | teks | VARCHAR2(4000) | kolom CSV `DISCOUNT_RATE_FWTWD` apa adanya |
| `C_RATE_EQ` | teks | VARCHAR2(4000) | kolom CSV `RATE_EQ` apa adanya |
| `C_NET_RATE_EQ` | teks | VARCHAR2(4000) | kolom CSV `NET_RATE_EQ` apa adanya |
| `C_DISCOUNT_RATE_EQ` | teks | VARCHAR2(4000) | kolom CSV `DISCOUNT_RATE_EQ` apa adanya |
| `C_RATE_OTHERS` | teks | VARCHAR2(4000) | kolom CSV `RATE_OTHERS` apa adanya |
| `C_NET_RATE_OTHERS` | teks | VARCHAR2(4000) | kolom CSV `NET_RATE_OTHERS` apa adanya |
| `C_DISCOUNT_RATE_OTHERS` | teks | VARCHAR2(4000) | kolom CSV `DISCOUNT_RATE_OTHERS` apa adanya |
| `C_KODE_POS` | teks | VARCHAR2(4000) | kolom CSV `KODE_POS` apa adanya |
| `C_KECAMATAN` | teks | VARCHAR2(4000) | kolom CSV `KECAMATAN` apa adanya |
| `C_KABUPATEN_KOTA` | teks | VARCHAR2(4000) | kolom CSV `KABUPATEN_KOTA` apa adanya |
| `C_PROVINSI` | teks | VARCHAR2(4000) | kolom CSV `PROVINSI` apa adanya |
| `C_SUMBER_KODE_POS` | teks | VARCHAR2(4000) | kolom CSV `SUMBER_KODE_POS` apa adanya |
| `C_TINGKAT_KEYAKINAN` | teks | VARCHAR2(4000) | kolom CSV `TINGKAT_KEYAKINAN` apa adanya |
| `C_TRADING` | teks | VARCHAR2(4000) | kolom CSV cargo `TRADING` apa adanya (migrasi 258, tiket 74 P6) |
| `C_DEPARTURE` | teks | VARCHAR2(4000) | kolom CSV cargo `DEPARTURE` apa adanya (migrasi 258, tiket 74 P6) |
| `C_DESTINATION` | teks | VARCHAR2(4000) | kolom CSV cargo `DESTINATION` apa adanya (migrasi 258, tiket 74 P6) |
| `C_GOODS` | teks | VARCHAR2(4000) | kolom CSV cargo `GOODS` apa adanya (migrasi 258, tiket 74 P6) |
| `C_GOODS_NAME` | teks | VARCHAR2(4000) | kolom CSV cargo `GOODS_NAME` apa adanya (migrasi 258, tiket 74 P6) |
| `C_PACKING` | teks | VARCHAR2(4000) | kolom CSV cargo `PACKING` apa adanya (migrasi 258, tiket 74 P6) |
| `C_CONVEYENCE` | teks | VARCHAR2(4000) | kolom CSV cargo `CONVEYENCE` apa adanya (migrasi 258, tiket 74 P6) |
| `C_CONVEYENCE_NAME` | teks | VARCHAR2(4000) | kolom CSV cargo `CONVEYENCE_NAME` apa adanya (migrasi 258, tiket 74 P6) |
| `C_COVERAGE_IDX` | teks | VARCHAR2(4000) | kolom CSV cargo `COVERAGE_IDX` apa adanya (migrasi 258, tiket 74 P6) |
| `C_COVERAGE_NAME` | teks | VARCHAR2(4000) | kolom CSV cargo `COVERAGE_NAME` apa adanya (migrasi 258, tiket 74 P6) |
| `C_TSI` | teks | VARCHAR2(4000) | kolom CSV cargo `TSI` apa adanya (migrasi 258, tiket 74 P6) |
| `C_AGENT` | teks | VARCHAR2(4000) | kolom CSV cargo `AGENT` apa adanya (migrasi 258, tiket 74 P6) |
| `C_DATASHIP` | teks | VARCHAR2(4000) | kolom CSV cargo `DATASHIP` apa adanya (migrasi 258, tiket 74 P6) |
| `C_INVOICE` | teks | VARCHAR2(4000) | kolom CSV cargo `INVOICE` apa adanya (migrasi 258, tiket 74 P6) |
| `C_ETD` | teks | VARCHAR2(4000) | kolom CSV cargo `ETD` apa adanya (migrasi 258, tiket 74 P6) |
| `C_BLNUMBER` | teks | VARCHAR2(4000) | kolom CSV cargo `BLNUMBER` apa adanya (migrasi 258, tiket 74 P6) |

## T_NB_CSV_BARIS_ANEKA

Tiket 73 P6 — migrasi 250 (K73-52 opsi 1): penampung baris pekerjaan unggah berjenis `aneka`; satu kolom `C_<judul>` VARCHAR2(4000) per judul yang dibaca InsertUploadAneka_act (35, urut KolomCsvAneka); kunci (UNGGAH_ID, NO_BARIS). Tabel sendiri TANPA FK: 250 berjalan sebelum 257 yang membuat T_NB_CSV_UNGGAH.

| Kolom | Tipe | Tipe DDL (migrasi 250) | Isi |
| --- | --- | --- | --- |
| `UNGGAH_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_NB_CSV_UNGGAH.ID` (tanpa FK) |
| `NO_BARIS` | angka bulat | NUMBER(10) NOT NULL | nomor baris data 1..n |
| `C_NO` | teks | VARCHAR2(4000) | kolom CSV Aneka `NO` apa adanya |
| `C_TYPE_ADDRESS` | teks | VARCHAR2(4000) | kolom CSV Aneka `TYPE_ADDRESS` apa adanya |
| `C_ADDRESS` | teks | VARCHAR2(4000) | kolom CSV Aneka `ADDRESS` apa adanya |
| `C_BUILDING_NO` | teks | VARCHAR2(4000) | kolom CSV Aneka `BUILDING_NO` apa adanya |
| `C_POSTAL_CODE` | teks | VARCHAR2(4000) | kolom CSV Aneka `POSTAL_CODE` apa adanya |
| `C_COUNTRY` | teks | VARCHAR2(4000) | kolom CSV Aneka `COUNTRY` apa adanya |
| `C_OCCUPATION_NAME` | teks | VARCHAR2(4000) | kolom CSV Aneka `OCCUPATION_NAME` apa adanya |
| `C_OBJECT_DESCRIPTION` | teks | VARCHAR2(4000) | kolom CSV Aneka `OBJECT_DESCRIPTION` apa adanya |
| `C_SECTION` | teks | VARCHAR2(4000) | kolom CSV Aneka `SECTION` apa adanya |
| `C_CURRENCY` | teks | VARCHAR2(4000) | kolom CSV Aneka `CURRENCY` apa adanya |
| `C_TSI` | teks | VARCHAR2(4000) | kolom CSV Aneka `TSI` apa adanya |
| `C_YEAR` | teks | VARCHAR2(4000) | kolom CSV Aneka `YEAR` apa adanya |
| `C_PLACE_OF_MANUFACTURE` | teks | VARCHAR2(4000) | kolom CSV Aneka `PLACE_OF_MANUFACTURE` apa adanya |
| `C_OBJECT_TYPE` | teks | VARCHAR2(4000) | kolom CSV Aneka `OBJECT_TYPE` apa adanya |
| `C_NRT` | teks | VARCHAR2(4000) | kolom CSV Aneka `NRT` apa adanya |
| `C_GRT` | teks | VARCHAR2(4000) | kolom CSV Aneka `GRT` apa adanya |
| `C_DWT` | teks | VARCHAR2(4000) | kolom CSV Aneka `DWT` apa adanya |
| `C_FLAGNAME` | teks | VARCHAR2(4000) | kolom CSV Aneka `FLAGNAME` apa adanya |
| `C_CONSTRUCTION` | teks | VARCHAR2(4000) | kolom CSV Aneka `CONSTRUCTION` apa adanya |
| `C_CLASSIFICATION` | teks | VARCHAR2(4000) | kolom CSV Aneka `CLASSIFICATION` apa adanya |
| `C_COVERAGE_NAME` | teks | VARCHAR2(4000) | kolom CSV Aneka `COVERAGE_NAME` apa adanya |
| `C_CALCULATE_METHOD` | teks | VARCHAR2(4000) | kolom CSV Aneka `CALCULATE_METHOD` apa adanya |
| `C_DAY` | teks | VARCHAR2(4000) | kolom CSV Aneka `DAY` apa adanya |
| `C_PCT_LOADING` | teks | VARCHAR2(4000) | kolom CSV Aneka `PCT_LOADING` apa adanya |
| `C_PCT_SCALE` | teks | VARCHAR2(4000) | kolom CSV Aneka `PCT_SCALE` apa adanya |
| `C_PCT_RATE` | teks | VARCHAR2(4000) | kolom CSV Aneka `PCT_RATE` apa adanya |
| `C_PREMIUM` | teks | VARCHAR2(4000) | kolom CSV Aneka `PREMIUM` apa adanya |
| `C_PCT_DISCOUNT` | teks | VARCHAR2(4000) | kolom CSV Aneka `PCT_DISCOUNT` apa adanya |
| `C_DISCOUNT` | teks | VARCHAR2(4000) | kolom CSV Aneka `DISCOUNT` apa adanya |
| `C_PCT_LOL` | teks | VARCHAR2(4000) | kolom CSV Aneka `PCT_LOL` apa adanya |
| `C_LOL` | teks | VARCHAR2(4000) | kolom CSV Aneka `LOL` apa adanya |
| `C_SERIALNO` | teks | VARCHAR2(4000) | kolom CSV Aneka `SERIALNO` apa adanya |
| `C_CHASSISNUMBER` | teks | VARCHAR2(4000) | kolom CSV Aneka `CHASSISNUMBER` apa adanya |
| `C_ENGINENUMBER` | teks | VARCHAR2(4000) | kolom CSV Aneka `ENGINENUMBER` apa adanya |
| `C_KELURAHAN` | teks | VARCHAR2(4000) | kolom CSV Aneka `KELURAHAN` apa adanya |

## T_NB_CSV_PESAN

Tiket 75 — migrasi 257: galat dan pesan per baris pekerjaan unggah; kunci (UNGGAH_ID, URUT).

| Kolom | Tipe | Tipe DDL (migrasi 257) | Isi |
| --- | --- | --- | --- |
| `UNGGAH_ID` | angka bulat | NUMBER(19) NOT NULL | = `T_NB_CSV_UNGGAH.ID` (FK) |
| `URUT` | angka bulat | NUMBER(10) NOT NULL | urutan per pekerjaan |
| `JENIS` | teks | VARCHAR2(10) NOT NULL | galat / pesan |
| `NO_BARIS` | angka bulat | NUMBER(10) | nomor baris data (NULL = tingkat berkas) |
| `TEKS` | teks | VARCHAR2(4000) NOT NULL | teks galat / pesan |

## Kolom status aktif objek master (bahan menu master)

Permintaan work owner 04-10-2026 (diteruskan sesi 0f): menu master insert / update / aktif-nonaktif untuk tabel flat —
penempatan belum diputuskan. Catatan ini **hanya** mendaftar status yang SUDAH ada di DDL `D:\migrasi\RNM\DDL\<NAMA>.txt`
`[terverifikasi]`, dihitung dua cara (grep pola nama kolom + urai daftar kolom DDL; keduanya sama). Arti nilainya
`belum terverifikasi` kecuali disebut.

| Objek | Kolom / saringan status di DDL | Catatan |
| --- | --- | --- |
| `ACCUMULATION` | tidak ada kolom; view menyaring `JSONDATA.IsActive IS NULL` | baris nonaktif TIDAK ikut tersalin ke tabel flat (Master Data 760, dulu 196); tabel flat mendapat `STS_AKTIF` |
| `PROVINCE`, `ACCUMULATEDTYPE`, `CZONE`, `NATION`, `DISTRICT`, `COVERAGE_FACIN` | tidak ada | — |
| `CITY` | `JABODETABEKSTATUS` | bukan status aktif `[dugaan]` (nama) |
| `COVERAGE` | `ACTIVESTATUS`, `STATUS`, `STATUSMAINCOVER`, `TRANSWORKSHOPSTATUS` | `ACTIVESTATUS = 1 / NULL` = aktif menurut RD BrowseCoverageFacIn_RD (butir 96) |
| `RW` | `STS_AKTIF` VARCHAR2(10) | `"1"` = aktif (saringan RD BrowseRW_RD) |
| `V_JN_OBJ_ITEM` (tabel `OBJECTITEMTYPE`) | `ISACTIVE` | nilai `belum terverifikasi` |
| `BRANDDETAIL` | `ISACTIVE`, `ACTIVESTATUS` (JSON `IsActive`, `ActiveStatus`) | tidak dibaca modul mana pun |
| `VJ_M_TYPE_PROPERTY_PLAN` | `STS_AKTIF` (JSON `STS_AKTIF`) | tidak dibaca modul mana pun |

## MARKETINGOFFICER

Tiket 31 (pilihan Marketing Name). Tabel warisan `POOLDATA`, **baca saja**. Sumber tipe `[terverifikasi]`: DDL
`D:\migrasi\RNM\DDL\MARKETINGOFFICER.txt` (15 kolom; `PEGA_MARKETINGOFFICER.txt` adalah prosedur penulis tabel
ini, bukan tabel). Hanya tiga kolom di bawah yang dibaca.

Struktur lengkapnya digambarkan `modul/marketingofficer/docs/STRUKTUR-TABEL-MARKETINGOFFICER.md` — modul itu yang
menyatakannya tabel warisan sejak merge `origin/dev` 04-10-2026 (satu tabel bersama = satu bentuk; baris warisannya di
`MODUL.md` modul ini dicabut). Yang dibaca di sini:

- `ID` — VARCHAR2(100); nilai pilihan (`.ID`) → disimpan ke `T_QUOTATIONDATA.MOID`; urutan DESC.
- `CLIENTNAME` — VARCHAR2(100); teks pilihan (pyPrompt `.ClientName`).
- `MOSTATUS` — VARCHAR2(100); saringan RD `.MOStatus = "1"`.

## AGENT

Tiket 33 (popup Change SOB). Tabel warisan `POOLDATA.AGENT`, **baca saja**. Sumber tipe `[terverifikasi]`: DDL
`D:\migrasi\RNM\DDL\AGENT.txt` (ditambahkan work owner 03-10-2026; 28 kolom dihitung dua cara; tanpa PK, dua indeks).
⚠️ Kolom tipe-2 dieja **`AGENTTPYE2`** di DDL (salah eja di basis data) — dipetakan dari properti `.AgentType2` `[dugaan]`.
Hanya lima kolom di bawah yang dibaca.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(1000 BYTE) | kode SOB (`.ID`) → `T_QUOTATIONDATA.SOURCE_OF_BUSINESS` (≤ 50, butir 79.2); dicari; urut DESC |
| `CLIENTID` | VARCHAR2(1000 BYTE) | Client ID; syarat `IS NOT NULL`; dicari |
| `CLIENTNAME` | VARCHAR2(1000 BYTE) | nama SOB (`.ClientName`) → `T_QUOTATIONDATA.SOB_NAME` (500); dicari |
| `STATUSACTIVE` | VARCHAR2(1000 BYTE) | syarat `= '1'` (teks) |
| `AGENTTPYE2` | VARCHAR2(1000 BYTE) | syarat `IS NULL OR <> 'LIFE INSURANCE'` (butir 79.1) |

## RISKADDRESS

Tiket 36 (popup Choose Risk Address) dan 37 (Add). Tabel warisan `POOLDATA.RISKADDRESS`, **dibaca dan DISISIPI** (tiket 37:
alamat baru, port Go prosedur `InsertUpdateRISKADDRESS`, butir 81 — ID `GETCURRENTSITE || LPAD(TO_CHAR(RISKADDRESS_SEQ.NEXTVAL),12,'0')`, ≤ 15 karakter
memakai fungsi dan sequence POOLDATA yang ada; tidak dibuat modul ini). Sumber tipe `[terverifikasi]`:
DDL `D:\migrasi\RNM\DDL\RISKADDRESS.txt` (16 kolom, semua VARCHAR2(4000), tanpa PK). Dibaca sembilan kolom (kolom hasil RD
`BrowseRisksAddress_RD`); kesembilan kolom yang sama **ditulis** saat alamat baru disisipkan (tiket 37).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(4000) | nilai Pilih → `.Property.AlmRiskID`; urutan (A118) |
| `TITLE` | VARCHAR2(4000) | Type → `.Property.RoadType` |
| `ADDRESS` | VARCHAR2(4000) | Address → `.Property.RoadName`; saringan A (JALAN) |
| `NATIONNAME` | VARCHAR2(4000) | Country; saringan C |
| `PROVINCENAME` | VARCHAR2(4000) | Province; saringan D |
| `CITYNAME` | VARCHAR2(4000) | City; saringan E |
| `DISTRICTNAME` | VARCHAR2(4000) | District; saringan F |
| `TERRITORYNAME` | VARCHAR2(4000) | Territory; saringan G |
| `POSTALCODE` | VARCHAR2(4000) | Zip Code; saringan B; kunci JOIN ke `RW.ZIPCODE` |

## RW

Tiket 36/37. Tabel warisan `POOLDATA.RW`, **baca saja** — JOIN RD tiket 36 dan saran Zip Code tiket 37; JOIN RD (`ASM-FW-GISFW-Int-RW`, prefix `RW`, INNER JOIN
`.PostalCode = RW.ZipCode`; kelas → tabel `[terverifikasi]` sejak tiket 46: `RDBList\CheckZipCode_SQL.xml` kelas Int-RW
"from rw"); tiket 46: popup akumulasi (JOIN RD, subkueri jalur kota / kecamatan, saran Area). Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\RW.txt` (19 kolom;
`RDBMASTERRW.txt` adalah prosedur penulisnya).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ZIPCODE` | VARCHAR2(4000) | syarat JOIN `= RISKADDRESS.POSTALCODE` (tiket 36); saran Zip Code "diawali" (tiket 37) |
| `NOTE` | VARCHAR2(4000) | Territory saran Zip Code (RD `.Note`, tiket 37) |
| `DISTRICTNAME` | VARCHAR2(4000) | District saran Zip Code |
| `CITYNAME` | VARCHAR2(4000) | City saran Zip Code |
| `PROVINCENAME` | VARCHAR2(4000) | Province saran Zip Code |
| `NATION` | VARCHAR2(100) | Country saran Zip Code — RD `.NATIONNAME`; alias `NATION as "NATIONNAME"` `[terverifikasi]` `RDBList\BrowseRW2_SQL.xml` |
| `STS_AKTIF` | VARCHAR2(10) | saringan RD `= "1"` (tiket 37; saran Area tiket 46) |
| `CITYID` | VARCHAR2(4000) | tiket 46 — `GetAccumulationProvince_SQL` `where cityid = ...` |
| `DISTRICTID` | VARCHAR2(4000) | tiket 46 — `GetAccumulationDistrict_SQL` `where DISTRICTID = ...` |
| `CZONE` | VARCHAR2(4000) | tiket 46 — kolom laporan BrowseRW_RD (ikut DISTINCT saran Area) |

## Tabel master akumulasi (ACCUMULATION, PROVINCE, ACCUMULATEDTYPE, CZONE)

Tiket 46 — **dipindah ke modul `masterdata`** (04-10-2026, MD-1), kini **`masterprovince`** (pemecahan delapan modul master, 04-10-2026): tabel flat migrasi 880 (dulu masterdata 760, dulu nbfacin 196, ditahan sebelum dijalankan), kolom dan status `STS_AKTIF` di `modul/masterprovince/docs/STRUKTUR-TABEL-MASTER-DATA.md`. nbfacin hanya MEMBACA (saran / popup akumulasi; baris aktif saja). Penyaring aktif ikut membaca tabel Master Data
`CITYINPUT` / `DISTRICTINPUT` (`STS_AKTIF`, untuk view CITY / DISTRICT) dan `T_MASTER_STATUS` (NATION) — MD-5.

## CITY

Tiket 46 (saran City popup akumulasi, `BrowseCityInput_RD`). **View** warisan `POOLDATA.CITY`, **baca saja**. Sumber
`[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\CITY.txt` (view atas CITYINPUT + DISTRICTINPUT + RWINPUT + PROVINCE + BRANCH;
tipe kolom tidak tertulis). Kelas Int-CITY → view CITY `[dugaan]` (nama dan kolom sama persis; tidak ada SQL kelas itu).

| Kolom | Dibaca untuk |
| --- | --- |
| `ID` | `id` (→ `cityId` pencarian) |
| `NOTE` | `label`; medan cari |
| `PROVINCEID` | saringan induk (ProvinceID terpilih) |

## DISTRICT

Tiket 46 (saran District, `BrowseDistrictInputC_RD`). **View** warisan `POOLDATA.DISTRICT`, **baca saja**. Sumber
`[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\DISTRICT.txt` (view atas DISTRICTINPUT + RWINPUT + CityInput). Kelas → view
`[dugaan]` (seperti CITY).

| Kolom | Dibaca untuk |
| --- | --- |
| `ID` | `id` (→ `districtId` pencarian) |
| `DISTRICTNAME` | `label`; medan cari |
| `CITYNAME` | saringan induk (NAMA kota terpilih) |
| `CITYID` | kolom laporan RD (ikut DISTINCT) |

## JSON_POLIS

Tiket 46 (pencarian akumulasi menurut nomor polis, `RDBList\GetSummaryRiskAccumPolis_Sql.xml`). Tabel warisan
`POOLDATA.JSON_POLIS`, **baca saja**. Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\JSON_POLIS.txt`.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `NOPOLIS` | VARCHAR2(100) NOT NULL | saringan `= policyNo` |
| `DATA_JSON` | CLOB (JSON) | `JSON_TABLE` `$.LocationList[*].Property.PropertyItemList[*].CoverageList[*].AccumulationCode` |

## OCCUPATION

Tiket 38 (saran Occupation Surrounding Risk). Tabel warisan `POOLDATA.OCCUPATION`, **baca saja**. Sumber tipe
`[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\OCCUPATION.txt`; kelas → tabel `[terverifikasi]`
`RDBList\SearchOccupationIDSQL.xml` (kelas `ASM-FW-GISFW-INT-OCCUPATION`, `FROM OCCUPATION`). RD
`BrowseOccupationFacInFIRE_RD`. Hanya tujuh kolom di bawah yang dibaca (tiga terakhir: rate FLEXAS tiket 52,
`RDBList\RateFlexas_Sql.xml`).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `OLDID` | VARCHAR2(1000) | nilai saran (`.OldID`) → `T_SURROUNDINGRISK.*_OCCUPATION`; dicari; urutan kedua |
| `NAME` | VARCHAR2(1000) | nama (`.Name`) → `*_NOTE` di layar; dicari; urutan pertama |
| `TYPE` | VARCHAR2(1000) | saringan RD `.Type = Param.TYPE` = `'FIRE'` |
| `KDRISKEXPOSURE` | VARCHAR2(1000) | `kdRiskExposure` → Category okupasi (`SetDataOccupation`, tiket 40); tiket 54 Scoring Risk Occupation "01".."05" → ChechBox1..5 (`GetRiskExposureOccupationFire`: `TYPE = 'FIRE'`, OLDID, baris pertama urut ID) |
| `OCCUPATIONCODE` | VARCHAR2(1000) | tiket 52 FLEXAS: saringan `IN` kode `occupationIds` (`= .OccupationId` SQL asal), urutan pertama |
| `RATELOWER` | NUMBER(20,6) | tiket 52 FLEXAS: terbesar atas okupasi objek → `.Rate`; NULL dilewati |
| `ID` | VARCHAR2(1000) | tiket 52: urutan kedua — baris pertama per kode bila ganda (`pxResults(1)` kueri tanpa urutan, `[dugaan]`) |

## V_JN_OBJ_ITEM

Tiket 39 (pilihan Object Item Type). View warisan `POOLDATA`, **baca saja**. Kelas → tabel `[terverifikasi]`
`RDBList\GetObjectItembyName_SQL.xml` (kelas `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM`, `FROM V_JN_OBJ_ITEM ... AND ISACTIVE ='1'`).
RD `BrowseV_JN_OBJ_ITEM` (DISTINCT, urut JN_OBJ_ITEM, maks 10000) dan `GetObjectItem` (KETERANGAN). ⚠️ DDL view **tidak ada**
di `DDL\` — tipe kolom `belum terverifikasi`; MJOI_KODE dibaca lewat `TO_CHAR`.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `MJOI_KODE` | belum terverifikasi | `kode` → `T_PROPERTYITEMLIST.ITEM_TYPE_ID` |
| `JN_OBJ_ITEM` | belum terverifikasi | `nama` → `ITEM_TYPE`; urutan |
| `KETERANGAN` | belum terverifikasi | `keterangan` → `PROPERTI_ITEM_NOTE` (layar) |
| `ISACTIVE` | belum terverifikasi | saringan RD `.ISACTIVE = 1` (dibandingkan `'1'` seperti SQL Pega) |

## CURRENCY

Tiket 39 (pilihan Currency + pemeriksaan mata uang item). Tabel warisan `POOLDATA`, **baca saja**. Kelas → tabel
`[terverifikasi]` `RDBList\GetAllCurrency.xml` (kelas `ASM-FW-GISFW-INT-CURRENCY`, `select Currency as CURR, OldID as OldID
from currency`). RD `BrowseCurrency_RD` (`.Currency != "ITL"`, maks 500, tanpa urutan). ⚠️ DDL tabel **tidak ada** di
`DDL\` — tipe kolom `belum terverifikasi`.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `CURRENCY` | belum terverifikasi | pilihan Currency (≠ `ITL`), urut CURRENCY (A135); anggota sah `T_PROPERTYITEMLIST.CURRENCY` |

## TABLEOFLIMIT

Tiket 40 (popup Choose Class of Construction). Tabel warisan `POOLDATA.TABLEOFLIMIT`, **baca saja**. Sumber tipe
`[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\TABLEOFLIMIT.txt` (03-10-2026, seluruh kolom VARCHAR2(4000 BYTE)). RD
`BrowseTableOfLimit_RD`. Hanya lima kolom di bawah yang dibaca. `TAHUN` **tidak** disaring (A161, menggantikan A153 —
popup tombol Pega mengirim Tahun kosong).

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `BIZCODE` | VARCHAR2(4000 BYTE) | saringan = `BUSINESS.ID` Class of Business case (butir 89); DISTINCT |
| `CATEGORY` | VARCHAR2(4000 BYTE) | saringan `category` (bila diisi); urutan pertama |
| `DESCRIPTION` | VARCHAR2(4000 BYTE) | `description` → Class of Construction; urutan kedua |
| `PCTLIMIT` | VARCHAR2(4000 BYTE) | `pctLimit` teks apa adanya |
| `NOTE` | VARCHAR2(4000 BYTE) | kolom laporan RD (ikut DISTINCT), tidak dikirim |

## COVERAGE

Tiket 43 (lima coverage otomatis dan popup Choose Coverage). Tabel warisan `POOLDATA.COVERAGE`, **baca saja**. Sumber
tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\COVERAGE.txt`. Coverage otomatis: `Activity\AddCoverageAutoFire.xml`
Obj-Browse `ASM-FW-GISFW-Int-COVERAGE`, `RDBList\SearchCoverageIDSQL.xml` ("FROM COVERAGE") `[terverifikasi]`. Popup:
keputusan work owner butir 96 — RD `BrowseCoverageFacIn_RD` (`DDL\BrowseCoverageFacIn_RD.xml`) dibaca dari tabel ini,
**bukan** view `COVERAGE_FACIN` (tidak dibaca lagi). Enam kolom dibaca.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `ID` | VARCHAR2(4000 BYTE) | otomatis: saringan = lima kode `AddCoverageAutoFire` (100815, 100828, 100829, 100825, 100840; urut korpus); popup: `id` → `T_COVERAGELIST.COVERAGE`, ikut DISTINCT |
| `NAME` | VARCHAR2(4000 BYTE) | `nama` → `.CoverageNote`; popup: saringan `cari` (Contains, tidak peka huruf — `pyCaseInsensitive`), urutan pertama |
| `OLDID` | VARCHAR2(23 BYTE) | `oldId` → `.OLDID`; popup: urutan kedua |
| `TYPE` | VARCHAR2(4000 BYTE) | popup: saringan `= 'FIRE'` (peka huruf, `=` RD) |
| `ACTIVESTATUS` | VARCHAR2(4000 BYTE) | popup: saringan `= '1'` atau kosong (D OR E RD) |
| `BUSINESSCODE` | VARCHAR2(4000 BYTE) | popup: kolom laporan RD (ikut DISTINCT), tidak dikirim; saringan B dikirim kosong → dibuang |

## ZONES

Tiket 52 (autocomplete Zone form coverage FIRE, RD `BrowseZones_RD`). **View** warisan `POOLDATA.ZONES`, **baca saja**.
Sumber `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\ZONES.txt` (view atas JSON `m_zones`; tipe kolom tidak tertulis). Kelas
`ASM-FW-GISFW-Int-ZONES` → view `ZONES` `[dugaan]` (nama dan kolom laporan RD sama persis). Saringan RD kosong dibuang
(`pyUseNullIfEmpty` false), tidak peka huruf (`pyCaseInsensitive` true); urut `DESCRIPTION`, `ID`; batas saran 50.

| Kolom | Dibaca untuk |
| --- | --- |
| `ID` | `id` saran |
| `DESCRIPTION` | `label` saran → `.Zone`; medan cari `q` (Contains) |
| `OLDID` | saringan I `= oldId` (OLDID coverage) |
| `ZONE` | saringan C..H: X / Y / Z (coverage 100835 / 100828 / 100817 / 100816) atau 1–6 (100829 / 100821), `SetPropertyCoverageDT_FacIn` |

## CONSTRUCTION

Tiket 52 (autocomplete 4.2 Construction). **View** warisan `POOLDATA.CONSTRUCTION`, **baca saja**. Sumber
`[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\CONSTRUCTION.txt` (view atas JSON `m_construction`; tipe kolom tidak tertulis);
`RDBList\SearchClassConstructionReas_SQL.xml` = `select distinct construction from construction order by CONSTRUCTION`.

| Kolom | Dibaca untuk |
| --- | --- |
| `CONSTRUCTION` | `id` = `label` saran → `.ConstructionEarthquake`; DISTINCT, urut, medan cari `q` (Contains, tidak peka huruf) |

## ZONEOJK

Tiket 52 (rate 4.2 PRGBI, `CountPremiRate_ACT` + `RateEarthquake_Sql`). **Tabel** warisan `POOLDATA.ZONEOJK`, **baca saja**.
Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\ZONEOJK.txt`.

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `DESCRIPTION` | VARCHAR2(1000) | saringan `= .Zone` (persis, peka huruf seperti SQL asal) |
| `RATE` | NUMBER(20,3) | → `.Rate` (‰ Gross Rate) coverage OLDID "4.2 PRGBI" bila `callAct` "rate"; tanpa baris / NULL → Rate tetap |
| `ID` | VARCHAR2(20) | urutan saja — baris pertama bila DESCRIPTION ganda (`pxResults(1)` kueri tanpa urutan, `[dugaan]`) |
