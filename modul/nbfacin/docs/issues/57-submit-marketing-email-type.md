# 57: Submit tahap marketing - Email Type, Comment, dan alur sesudahnya

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 06-10-2026 (gambar layar
> "Email type * Accept / Reject / Banding / Decline", "Comment: *") "selanjutnya buatin ini, dan saat di submit berjalan
> sesuai flow nya". Survei korpus oleh dua subagen (pengurai XML); butir kunci diperiksa ulang sesi 0f: syarat wajib /
> nonaktif `EmailSection`, konektor `Decision3`, `SendOfferEmail` langkah 1, `M_LIMIT_PROPERTYY` di `GetAksepBanding_SQL`,
> keberadaan DDL tabel tulis.

**Status:** frontend selesai 07-10-2026 (uji hijau); backend (proteksi, alur, tulis tabel lama) diserahkan ke sesi c3.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

### Layar
- `Flow\InputInwardFacultativeOffer.xml` Assignment12 MARKETING (`ReasFacInMarketing`) → flow action `InwardFacultative` →
  `Section\InputInwardFacultative.xml` → `ComfirmPolis_NonLife` → **`Section\EmailSection.xml`** (bila
  `.IsCedingConfirm != 'Binding' && .FlagOnGoingPolicy != 2` - tahap Offer).
- `.EmailType` radio wajib (`.FlagOnGoingPolicy==0 && PositionNote != 'ReasFacInTeamLeader'`); kode 1 Accept · 2 Reject ·
  3 Ask · 4 Banding · 7 Decline · 9 Revise (`SaveViewSuggest`); gambar layar work owner menampilkan Accept / Reject / Banding
  / Decline. Daftar opsi properti `EmailType` tidak di korpus.
- `.OfferFacIn.Comment` textarea wajib (`.IsCedingConfirm != 'Binding'`), nonaktif bila `.EmailType = 7` (diisi popup
  "Choose Cause of decline" `InputCauseOfDecline` → `SetCauseOfDecline_ACT`).
- "Banding to" `.BandingTo` (wajib bila `.EmailType = 4`): `RDBList\GetAksepBanding_SQL` = JABATAN / JABATAN_ATASAN
  `M_LIMIT_PROPERTYY` dari workbasket non-marketing TERAKHIR di `HISTORYAKSEPTASIPEGA` + TEAM_GROUP case → kosong untuk case
  yang belum pernah ke UW.
- Grid riwayat Date · PIC · Approval · Status · Comment (`OfferFacIn.ViewSuggest`).
- Tombol aksi berganti menurut Email Type (Send Email / Reject / Banding / Decline); Reject dan Decline meminta konfirmasi
  "Are you sure want to reject?" / "Are you sure want to decline?" (`Section\RejectNotificationSection.xml`).

### Proteksi (Accept) - `Activity\Protection_Act.xml`
Reject / Ask / Decline melewati seluruh proteksi (langkah 7–8). Accept diperiksa ±40 pesan verbatim, mis. "Scoring Risk
Cannot Empty !!" (IsFire), "Survey report cannot be empty", "Ceding Co cannot be empty", lokasi / okupasi / zip / akumulasi /
coverage / deductible (`ProtectFIREMBUPA_Act`), spreading 100 % dan premi (`CheckSpreadingProtectFire_ACT`), payment /
installment (`ProtectPremiPolicy_Act`), mata uang TSI (`ProtectCurrencyTSI_Act`).

### Alur sesudah Submit marketing
- `DecisionTable\IsUWAccepted.xml` atas `.ProposalAcceptStatus`: 1 confirm · 2 reject · 3 ask · 4 banding · 9 revise ·
  bawaan decline. **`Decision3` hanya punya 3 konektor**: confirm → `Decision36`, banding → `Utility3 SetBanding_ACT`,
  decline → `End5`. **Reject / Ask / Revise tidak punya jalur** dari titik marketing `[dugaan]` galat flow / tertahan.
- Accept (FIRE NB biasa): non-Life → `SaveJsonOfferFacIn_Act` (CALL `POOLDATA.PEGA_M_JSON_OFFER` + `FACINOFFER`) → bukan
  special case / group / PKS ASM / financial / declaration policy → **Assignment14 TEAM LEADER** (`ReasFacInTeamLeader`,
  NBStatus "<ID> IS IN TEAM LEADER'S INBOX"). Tangga akseptasi (`GetLimitAkseptasi_ActFlow`, `M_LIMIT_*`) baru berjalan di
  tahap UW sesudah Team Leader dan Junior UW B.
- Decline: `End5`; NBStatus "<ID> WAS DECLINED BY <operator>".
- Reject (post-processing saja): NBStatus "<ID> IS IN <MARKETING>'S OPPORTUNITY".
- Banding: `SetBanding_ACT` memetakan `LetterNo` (= Banding to) ke ticket assignment jabatan itu.
- **Tanpa email di tahap marketing**: `Activity\SendOfferEmail.xml` langkah 1 `Exit-Activity` tanpa syarat; kirim di
  `SendEmailtoCeding_ACT` di-`//`.
- Post-processing `SetValidateDate_PostAct`: riwayat `ViewSuggest` → `POOLDATA.HISTORYAKSEPTASIPRODUCTION`
  (`InsertViewSuggest_SQL`), `HISTORYAKSEPTASIPEGA` (`InsertHistoryAkseptasiPega_Sql`), `POOLDATA.FACINPERFORMANCE`
  (Decline), NBStatus, Obj-Save; `AddToListSuggestOfferFacIn_DT` → `SetCorrespondenceReadOnly` (kunci korespondensi).
- DDL tersedia di `DDL\`: HISTORYAKSEPTASIPEGA, HISTORYAKSEPTASIPRODUCTION, FACINPERFORMANCE, JSON_OFFER,
  PEGA_M_JSON_OFFER, FACINOFFER, M_LIMIT_*.

## Batas pengetahuan
- Label opsi Email Type, HTML email, `ObjSave`, field value `ASMMessageQuotation` tidak ada di korpus.
- Popup "Choose Cause of decline" (`InputCauseOfDecline`: Term and Conditions (T/C) / Non - T/C, % RNM, % SOB, Remarks)
  belum ditelusur rinci.
- ADR-0043: aplikasi tidak memanggil stored procedure (`PEGA_M_JSON_OFFER`).

## Popup Cause of Decline `[terverifikasi]`
- `Section\InputCauseOfDecline.xml` (judul "Cause of Decline"): grid T/C (`RDBList\GetCauseOfDecline.xml`
  `select TERMCONDITION from POOLDATA.CAUSEOFDECLINE WHERE DESCRIPTION IS NULL`) dengan centang prioritas I / II, Option,
  Remarks, % RNM, % SOB; grid Non - T/C (`GetCauseOfDecline2`: `TERMCONDITION2 IS NOT NULL`); Check →
  `CekCauseOfDecline_ACT` (pesan verbatim, mis. " PctRNM or PctSOB Can't be NULL", " Remaks Can't be NULL"); Submit →
  `SetCauseOfDecline_ACT` (Comment = ReasonOfDecline + "\n" + ReasonOfDecline2); % RNM otomatis `SetCODRate`.
- Sumber pilihan Option (`Getdescription_ACT`, RDB `GetDescription1`) dan DDL `POOLDATA.CAUSEOFDECLINE` TIDAK ada di korpus
  `[pertanyaan terbuka]` — layar memakai isian teks bila backend tidak mengirim pilihan.

## Keputusan work owner 07-10-2026 ("setuju" + "setuju, jangan halu" atas 7 saran sesi 0f)

1. **Reject** — case tetap di marketing, riwayat REJECT, NBStatus "… IS IN <MARKETING>'S OPPORTUNITY".
2. **Banding** ikut Pega — daftar dari riwayat approver; kosong → pesan, Banding belum dapat dipilih.
3. **Semua proteksi Accept** Pega dibangun; yang datanya belum ada dilaporkan dulu.
4. **Riwayat / performance ke tabel lama** (HISTORYAKSEPTASIPEGA, HISTORYAKSEPTASIPRODUCTION, FACINPERFORMANCE; ADR-0042 butir 2).
5. **JSON_OFFER / FACINOFFER ditunda** (ADR-0043: tanpa CALL prosedur).
6. **Tahap marketing dulu**; Team Leader / UW / tangga akseptasi menyusul.
7. **Hanya 4 opsi** (Accept / Reject / Banding / Decline).

## Frontend (sesi 0f)

- `components/BlokEmailType.tsx` (hook `useSubmitMarketing` + blok), `components/PopupSebabDecline.tsx`,
  `components/sebabDecline.ts` (aturan kolom / validasi / teks alasan), `BlokEmailType.test.tsx`; `api.ts` kontrak; `labels.ts`
  `EMAIL_TYPE` / `TEKS_EMAIL_TYPE` / `SEBAB_DECLINE` / `TEKS_SEBAB_DECLINE`; `nbfacin.css` `.nbf-et*` / `.nbf-sebab*`;
  `pages/InwardFacultative.tsx` (blok di bawah detail = urutan Pega cell 118; tombol Submit kaki halaman aktif bila boleh).

## Keputusan agent

- **K57-1** Satu tombol Submit di kaki halaman menggantikan tombol per Email Type Pega (Send Email / Reject / Banding / Decline
  di `EmailSection`); Reject / Decline tetap lewat konfirmasi verbatim.
- **K57-2** Proteksi Accept dijalankan saat Accept dipilih (Pega: `Protection_Act` pada klik radio) dan diulang server saat Submit.
- **K57-3** Memilih Non - T/C mengosongkan pilihan T/C (Pega hanya menonaktifkan kontrol; centang lama ikut tersimpan lalu
  alasannya ditimpa).
- **K57-4** Semua pesan validasi Cause of Decline ditampilkan (Pega menimpa `Local.Errmsg` → hanya pesan terakhir).
- **K57-5** Decline wajib prioritas I atau satu Non - T/C (Pega: tombol Decline nonaktif selama `ReasonOfDecline` kosong).
- **K57-6** Isian baris T/C yang tidak dicentang dinonaktifkan.
- **K57-6b (sesi 0f)** Penolakan proteksi backend 422 `{"galat": pesan digabung "; "}`: klien inti hanya menyimpan `galat`, jadi
  layar memulihkan daftar pesan dengan memecah "; " (`pesanProteksiDari`); diverifikasi tidak ada pesan proteksi yang memuat "; ".
- **K57-7 (sesi 0f, menjawab pertanyaan c3 07-10-2026) NBStatus** `[terverifikasi]`: di `SetValidateDate_PostAct` langkah
  NextUser "TEAM LEADER" ber-`//` (nonaktif), yang aktif "JUNIOR UNDERWRITER"; tetapi konektor `Decision29 → Assignment14`
  (flow) men-set NBStatus `… IS IN TEAM LEADER'S INBOX` sesudah post-processing → nilai akhir Accept = "NB IS IN TEAM LEADER'S
  INBOX". Reject = langkah 6.9, Decline = langkah 6.8 (verbatim, dua spasi). Banding: `SetBanding_ACT` memuat nama orang literal
  → tidak disalin (CLAUDE.md §4.10); dipakai "NB IS IN <JABATAN>'S INBOX" (penyimpangan sadar). Notify `SetNBStatus_Act` tidak
  diterapkan (`param.AssignTo` tidak terbukti). Decline: POSITION "Decline", STATUS_WORK tidak diubah `[dugaan]` (End5 tanpa
  status kerja).

## Backend (sesi c3, 07-10-2026)

**Status:** backend selesai; migrasi 208 DITULIS, belum dijalankan. Kontrak = `frontend/api.ts`. Tiga tabel lama POOLDATA
DITULIS (keputusan work owner butir 4) - tidak ada satu pun yang dijalankan atas basis data selama pengembangan; skema dan
hak tulis di DEV **belum diperiksa** (agent tidak menyentuh DEV) - DBA / work owner perlu memastikan akun aplikasi boleh
INSERT ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION`, `POOLDATA.HISTORYAKSEPTASIPEGA`, `POOLDATA.FACINPERFORMANCE`.

- **K57-8 Endpoint** - `GET /api/nbfacin/kasus/{caseId}/submit`, `POST …/submit/periksa {emailType}`, `POST …/submit`.
  Urutan galat Submit: 401 (handler, SEBELUM badan diurai) → 400 badan bukan JSON ketat → 503 → 404 → 403 (akun tanpa
  workbasket ReasFacInMarketing) → 409 (posisi bukan tahap marketing Offer atau FlagOnGoingPolicy ≠ 0) → 400 isian → 422 proteksi Accept (tanpa tulis) → satu transaksi (409 bila posisi berubah
  serentak: `UPDATE … WHERE POSITION = <lama>` 0 baris). 422 = `{galat, pesan, position, nbStatus}` - `pesan` larik verbatim,
  `galat` gabungannya dengan "; " (klien `minta` membaca `galat` untuk setiap jawaban tidak ok; `BlokEmailType` kini menampilkan
  422 lewat pita galat, bukan daftar proteksi - `[pertanyaan terbuka]` 0f: baca `pesan` dari 422?).
- **K57-9 bolehSubmit** - `posisiMarketing[POSITION]` ("Offer" / "ReasFacInMarketing") && FLAG_ONGOING_POLICY "0" && pelaku
  memegang workbasket ReasFacInMarketing (`inti.Pelaku.Peran` = workbasket login). `alasan` teks Indonesia.
- **K57-10 Riwayat** - `IDPEGA` = pzInsKey = `"ASM-FW-GISFW-WORK " + caseId` (`InsertHistoryAkseptasiPega` langkah 8 + indeks
  DBA `HISTOAKSPEGA_INDEX4`). Grid: TGL_INP `TO_CHAR 'DD-MM-YYYY HH24:MI'` (jam dinding WIB, tanpa konversi), PIC, APPROVAL,
  POSISI, KETERANGAN; urut TGL_INP lalu NOURUT turun. Satu baris per Submit `[terverifikasi]`: blok APPEND ViewSuggest
  `AddToListSuggestOfferFacIn_DT` ber-`pyDisabled=true` (baris 3, 3.1-3.4.2); yang aktif hanya APPLY_MODEL
  `SetCorrespondenceReadOnly`. Baris ditulis (`SaveViewSuggest` + `InsertViewSuggest_SQL`): TYPE_POLIS "NB", NOURUT =
  MAX+1 per IDPEGA, POSISI "Offer", PIC = M_LOGIN_GO.NAME (`OperatorID.pyUserName`; kosong → LOGIN_ID `[dugaan]`), TGL_INP =
  sekarang WIB jam 24 (Pega `hh` tanpa AM/PM - dibetulkan), DIV = DIVISION_CODE, TYPE "F" (SaveViewSuggest hanya untuk
  BusinessFac "F" - kasus Fac In `[dugaan]`), PUTARAN "1", APPROVAL Accept / Reject / Banding / Decline, KETERANGAN ≤ 3990
  karakter dan ≤ 4000 bita, AKSES_LOGIN = LOGIN_ID, **B2B dan BUSINESS_CODE NULL** (IsB2B / BusinessCode tidak ada di sistem
  baru), PERCENT_RNM = T_GENERAL_POLIS.PERCENT_SHARE. HISTORYAKSEPTASIPEGA: STATUS ACCEPT / REJECT / BANDING / DECLINE,
  USERNAME = nama tampilan, WORKBASKET "ReasFacInMarketing" (posisi SEBELUM Submit), ID_KOMITE / OPERATORID NULL (jalur Pega
  tidak mengisinya). FACINPERFORMANCE hanya Decline: POSITIONDOC "OFFER", STATUSDOC "DECLINE", USERNAME = LOGIN_ID
  (`pyUserIdentifier`), INSUREDNAME (T_M_ACCOUNT), SOB = SOB_NAME, MOLEADER, COB = T_NB_OPPORTUNITY.CLASS_OF_BUSINESS
  (`Quotation.BusinessName` `[dugaan]` padanan), RISLIPRECEIVED NULL; dilewati (dicatat log) bila kolom NOT NULL kosong -
  penjaga Pega langkah 12. Accept di tahap Offer TIDAK menulis FACINPERFORMANCE `[terverifikasi]` (STATUSDOC hanya diisi
  EmailTypeUW / Binding; langkah 13 butuh STATUSDOC ACCEPT).
- **K57-11 Marketing** - MarketingName / TeamGroup / MOLeader dari MARKETINGOFFICER **ber-ID = MOID** (Pega: `CLIENTNAME LIKE
  <MarketingName>%` + saringan MOSTATUS '1', BRANCHSTATUS kosong, CLIENTID2 ≠ 'LEADER' untuk MOLEADER; GetTeamGroup_RD).
  ID lebih tepat daripada LIKE nama `[keputusan agent]`. SetValidateDate langkah 8 (TeamGroup 4 → 3) ber-`//` - tidak diport.
- **K57-12 Opsi Banding** - `GetAksepBanding_SQL` verbatim (bind :1 pzInsKey, :2 'ReasFacInMarketing', :3 TeamGroup); jabatan
  kembar (M_LIMIT_PROPERTYY memuat SENIORUW dua kali untuk TEAM_GROUP 1) ditampilkan sekali `[keputusan agent]`. Tujuan Banding
  `SetBanding_ACT` langkah 3: JUW_A → ReasFacInJuniorUnderwritingA, JUW_B → ReasFacInJuniorUnderwriting, UNDERWRITER →
  ReasFacInUnderwriting, SENIORUW → ReasFacInSeniorUnderwriting, DEPHEADUNDERWRITER → ReasFacInDepHeadUnderwriting,
  KADIVFACULTATIVE → ReasFacInFacultativeDivHead, KADIVTEKNIK → ReasFacInGroupLeader, DIREKTURMARKETING →
  ReasFacInMarketingDirector, DIREKTURTEKNIK → ReasFacInTechnicalDirector. **MANAGERTEKNIK (3.7) dan Financial (3.6, kunci PositionNote) ber-`//`** → 400 "tidak punya
  tujuan Banding" walau muncul di opsi (M_LIMIT_PROPERTYY memuatnya) `[pertanyaan terbuka]` work owner. NB_STATUS Banding: teks
  Pega dipakai verbatim bila tanpa nama ("JUNIOR UNDERWRITER (A)" / "(B)", "UNDERWRITER"); yang memuat nama orang diganti
  JABATAN (K57-7). Posisi baru = workbasket tujuan; FlagOnGoingPolicy tetap "0".
- **K57-13 Posisi / NB_STATUS** - Accept → POSITION "ReasFacInTeamLeader", "NB IS IN TEAM LEADER'S INBOX"; Reject → POSITION
  tetap, "NB IS IN " + UPPER(CLIENTNAME tanpa " -") + "'S OPPORTUNITY"; Decline → POSITION "Decline", "NB WAS DECLINED BY  " +
  UPPER(M_LOGIN_GO.NAME), STATUS_WORK "Resolved-Decline" (K57-19); Accept / Reject / Banding: STATUS_WORK tidak diubah. NB_STATUS di `T_GENERAL_POLIS` (migrasi 208, MERGE: baris General belum
  ada → dibuat ber-ID saja). Setiap Submit mengunci seluruh korespondensi kasus (`IS_EDITABLE '0'`, juga penerimanya).
- **K57-14 Cause of Decline** - port `sebabDecline.ts` (`validasiSebab`, `komentarDecline` di services/submit.go); instrumen: 20 masukan yang
  jawabannya dihasilkan sebabDecline.ts sendiri (Node 24) - `TestSebabDeclineParitas`. Server menolak T/C / Non - T/C di luar
  daftar CAUSEOFDECLINE, prioritas ganda, persen bukan desimal (pertahanan server - BUKAN bagian sebabDecline.ts); Comment
  Decline disusun ulang server dari teks baris APA ADANYA (tanpa dipangkas, sama dengan popup). `sebabTc[].opsi` selalu `[]`
  (sumber `Getdescription_ACT` / RDB GetDescription1 tidak di korpus). `pctRnmOtomatis.rate` selalu "": `TempTotal` DIISI
  `SumTSIPremiSpreadedRNM_FIRE_Act` langkah 4 (per mata uang, lokasi Top Risk, Premium / TSI × 1000) yang dipanggil pemeriksa
  spreading - pemeriksa itu keluar untuk tahap marketing, jadi isinya saat popup marketing belum terverifikasi
  `[pertanyaan terbuka]` (ralat tinjauan: semula ditulis "tidak di korpus"). `riComm` = PAY_RI_COMMISION.
- **K57-15 Proteksi Accept** - tabel di bawah. Kasus diperlakukan bukan Group / bukan B2B ASM / bukan MOP / bukan EDM / FIRE
  (data penentunya tidak ada di sistem baru). Pesan kembar ditampilkan sekali; urutan = urutan langkah Pega. Langkah ber-`//`
  tidak diport (Protection_Act 35, ProtectFIREMBUPA 9.2.27.2 / 9.2.28.2, ProtectPremiPolicy 12).

### Tabel proteksi Accept (tahap marketing, FIRE NB)

| # | Pesan (verbatim) | Rule | Status |
| --- | --- | --- | --- |
| 1 | (keluar tanpa proteksi) | Protection_Act 2: InsuredID ORG-7889 / ORG-2344 atau > 400 lokasi | terpenuhi (INSURED_ID dibandingkan utuh atau akhiran " ORG-…" `[dugaan]` bentuknya) |
| 2 | "Deductible Acceptance cannot be empty" / "Comment Deductible cannot be empty" | Protection_Act 9-12 | tidak berlaku (hanya PositionNote TeamLeader) |
| 3 | "pyWorkPage.OfferFacIn.QuotationData.Type Tidak boleh kosong" | Protection_Act 13-15 | tidak berlaku (EDM) |
| 4 | "Harus memilih dari daftar nama marketing yang tersedia" | Protection_Act 16-17 (MOID kosong; field value `ASMMessageQuotation`) | terpenuhi - teks = **keputusan work owner 07-10-2026**, bukan teks Pega (field value tidak ada di korpus) |
| 5 | "Scoring Risk Cannot Empty !!" | Protection_Act 18 | terpenuhi (T_SCORINGRISK.FINAL_SCORE) |
| 6 | "Currency Pada Lokasi n Object Item m Bermasalah, Coba Pilih Ulang Currency" | ProtectCurrencyTSI_Act | **padanan**: Pega membandingkan CURRENCY(id) dengan nama; sistem baru hanya menyimpan nama → nama wajib ada di CURRENCY; hanya item gagal terakhir; Currency kosong tidak memicunya (Pega: "" = "" `[dugaan]`) - dilaporkan #21 |
| 7 | spreading 100 % / premium spreaded / ORS / special acceptance | CheckSpreadingProtect_ACT | **tidak berlaku di tahap ini** `[terverifikasi]`: langkah 2 keluar untuk PositionNote ReasFacInMarketing kecuali IsGroup "Group"; IsGroup tidak ada di sistem baru |
| 8 | "List Payment can't be null!" | ProtectPremiPolicy_Act 3-4 | terpenuhi (mata uang tab Payment) |
| 9 | "List Installment can't be null!" | ProtectPremiPolicy_Act 5.7 | terpenuhi |
| 10 | "Total premium RNM on installment didn't match with premium RNM!" | ProtectPremiPolicy_Act 5.5 / 10 | terpenuhi (Σ cicilan vs Premium / Commision / NetPremium 4 desimal, ditekan bila salah satu pasangan selisih < 0,001; bukan MOP / EDM; tanpa cicilan ikut muncul karena jumlah Pega mulai dari 0) |
| 11 | "Premi on paymnet can't be 0" | ProtectPremiPolicy_Act 11 | terpenuhi (PolicyType 1 tidak dapat diperiksa - tersimpan sebagai label, `[dugaan]` bukan 1) |
| 12 | "Premium Master Policy must be 0" | ProtectPremiPolicy_Act 8-9 | **belum** (IsMOP tidak ada) |
| 13 | "Province on Location ncan't be null!" / "City on Location ncan't be null!" / "ZipCode on Location ncan't be null!" | ProtectFIREMBUPA 9.2.3-8 | terpenuhi |
| 14 | "Occupation on location n can't be null!" | 9.2.9 | terpenuhi |
| 15 | "Loss record on location n can't be empty !" | 9.2.10 (kecuali NB-83801) | terpenuhi (T_LISTCAUSEOFLOSS per lokasi) |
| 16 | "ZipCode on Location n doesn’t exist in Master!" | 9.2.11 / 9.2.22 (CheckZipCode_SQL) | terpenuhi (RW) |
| 17 | "TERITORY  x on Location n doesn’t exist in Master!" | 9.2.13 / 9.2.23 | **belum**: kolom TERITORYNAME tidak ada di DDL RW (tabel kelas Int-RW belum terverifikasi) |
| 18 | "CITYNAME x …" / "DISTRICT x …" / "PROVINCE x … doesn’t exist in Master!" | 9.2.14-16 / 9.2.24-26 | terpenuhi (RW ZIPCODE + kolom; kode pos tak ada → ketiganya muncul `[dugaan]`, halaman Cek* kosong) |
| 19 | "Occupation ID on locationncan't be null!" | 9.2.27.4 | terpenuhi |
| 20 | "Kategori Okupasi ID x lokasin tidak boleh kosong" | 9.2.27.5-8 / 9.3 (GetCategoryOccupation) | terpenuhi (OCCUPATION.CATEGORY per OCCUPATIONCODE; aturan IsPreferredRisk / Top Risk) |
| 21 | "Currency on location n object Item m must be choose!" / "Choose Object Item on location n !" | 9.2.30.4-5 | terpenuhi |
| 22 | " Coverage c on location n object item m can't be empty!" / " Accumulation on location n object item m can't be empty!" | 9.2.30.6.3-4 | terpenuhi (IndexProperty / IndexPropertyItem = nomor lokasi / item `[dugaan]`) |
| 23 | "coverage on location n object item m can't be Net Rate!" | 9.2.30.7-8 | terpenuhi (TotalNetRate > 0 wajib memuat FLEXAS, 4.1A CC, 4.3, 4.2 PRGBI, OTHERS) |
| 24 | deductible "… should be added!" / "… DeductibleList i can't be empty!" | 9.2.30.6.6-7, 13 | tidak berlaku (FlagOnGoingPolicy 2) |
| 25 | "Double Input Case <IDPEGA>" | 9.2.32-35 (GetDataDoubleCase, FACINPRODUCTION) | **belum**: butuh BusinessCode (tidak ada di sistem baru) |
| 26 | "% Share cedant …" / "Ceding on list …" / "Total premi share cedant …" / "List Inward facultative cedant can't be null!" | ProtectShareCedant_Act | **belum**: hanya untuk IsGroup "Group" (tidak ada) |
| 27 | "Invalid Start Date!" | ProtectRenewal_Act | tidak berlaku (NB) |
| 28 | "Ceding doesn't exist in master" | ProtectCedingCo (IsB2B "ASM") | **belum**: IsB2B tidak ada |
| 29 | "Survey report cannot be empty" | Protection_Act 31-32 | **belum**: IsSurveyReport / SurveyReportList tidak ada |
| 30 | "Ceding Co cannot be empty" | Protection_Act 33-34 | **belum** untuk Group (IsGroup tidak ada); FlagOnGoingPolicy 2 tidak berlaku |
| 31 | marine / MBU / PA / Aneka / Golf | Protection_Act 27, ProtectCoverage 4-6, ProtectFIREMBUPA 7-8 | tidak berlaku (FIRE) |

Instrumen: `TestProteksiAcceptLolos` (kasus sah tanpa pesan) + `TestProteksiAcceptGagal` (16 pesan verbatim berurutan);
mutasi 17 titik (bind RW, penjaga posisi, NB_STATUS, jabatan Banding, workbasket, jam 24, performance, T/C, scoring, teks
lokasi, Top Risk, net rate, currency, pengecualian langkah 2, installment, teks alasan) - semuanya menggagalkan uji.

- **K57-16 Belum / ditunda** - JSON_OFFER / FACINOFFER (`SaveJsonOfferFacIn_Act`, CALL `PEGA_M_JSON_OFFER`) DITUNDA (butir 5,
  ADR-0043). Tanpa email (SendOfferEmail langkah 1 Exit-Activity). Notify `SetNBStatus_Act` tidak diterapkan (K57-7). Cabang
  Accept special case / group / PKS ASM / financial / declaration policy (Decision36 → … → Decision29) tidak diport: kasus
  FIRE NB biasa = konektor Decision29 Else → Assignment14 `[dugaan]` bahwa penentunya (IsGroup, IsSpecialAcceptance, …) tidak
  berlaku - datanya tidak ada di sistem baru. GetLimitAkseptasi_Act (SetValidateDate 4/5) ber-`//` - tidak diport.
- **K57-17 Urutan deploy** - 207 (kunci korespondensi) dan 208 (NB_STATUS) sebelum Submit dipakai; GET kasus tidak membaca
  keduanya. Hak INSERT tiga tabel lama POOLDATA wajib ada (K57-8).
- **K57-18 Tinjauan dua sumbu (standar + spec), 07-10-2026** - standar: nol pelanggaran aturan tertulis; dibetulkan galat
  baris dibungkus nama tabel, FACINPERFORMANCE lewat KosongJadiNil, label approval bernama, helper Cause of Decline tidak
  diekspor, cara membangkitkan ulang uji paritas ditulis. Spec: 22 pesan proteksi dicek huruf demi huruf terhadap korpus -
  cocok; dibetulkan: Comment Decline tanpa pemangkasan, Currency kosong tidak memunculkan pesan ProtectCurrencyTSI, 401
  sebelum badan diurai dan 403 sebelum 409, ralat sumber `TempTotal` (K57-14), 3.6 ber-`//` (K57-12). Dipertahankan dengan
  sadar: WORKBASKET HISTORYAKSEPTASIPEGA = "ReasFacInMarketing" (nilai PositionNote Pega; POSITION "Offer" sistem baru adalah
  padanannya - pembaca GetAksepBanding / GetFlagReject menyaring nilai Pega ini); NOURUT = MAX+1 per IDPEGA (Pega
  `.pxListSubscript` - sama selama baris riwayat tidak dihapus `[dugaan]`); baris General kosong dibuat ber-ID saja bila belum
  ada (182: hanya ID NOT NULL); opsi Banding kembar sekali (K57-12); pertahanan server Cause of Decline (K57-14).

## Keputusan work owner 07-10-2026 (lanjutan, atas pertanyaan terbuka backend)

1. Pesan proteksi MOID kosong (Pega field value `ASMMessageQuotation`, isinya tidak di korpus) = **"Harus memilih dari daftar
   nama marketing yang tersedia"** (teks keputusan work owner).
2. **MANAGERTEKNIK disembunyikan** dari pilihan Banding to (pemetaan `SetBanding_ACT` 3.7 ber-`//`); saran c3: saring semua
   jabatan tanpa pemetaan aktif.
3. Decline: **STATUS_WORK = "Resolved-Decline"**.
4. Data `POOLDATA.CAUSEOFDECLINE` diberikan (`D:\migrasi\RNM\DDL\CAUSEOFDECLINE.csv`, 41 baris): pilihan Option T/C X = TERMCONDITION baris
   ber-DESCRIPTION = X (POI / Extension of coverage / Risk Exposure) `[dugaan kuat dari bentuk data]` - `Getdescription_ACT`
   tetap tidak di korpus. Keanehan data: ID 10 "Clauses" berakhiran baris baru (dipangkas = maksud syarat Pega); ID 24
   DESCRIPTION1 tidak sama dengan TERMCONDITION (tidak dipakai - Option memakai TERMCONDITION).
- **K57-19 Jawaban work owner 07-10-2026** (verbatim: "1. ganti aja pesan text nya jadi \"Harus memilih dari daftar nama
  marketing yang tersedia\" 2. setuju 4. ganti jadi \"Resolved-Decline\" 5. sudah ditambah csv nya dengan nama
  CAUSEOFDECLINE"; diteruskan sesi 9d):
  1. Pesan MOID kosong = "Harus memilih dari daftar nama marketing yang tersedia" (bukan teks Pega).
  2. Opsi Banding: jabatan TANPA pemetaan aktif SetBanding_ACT disembunyikan, disaring generik (kini MANAGERTEKNIK); POST
     tetap 400 bila dikirim.
  3. Decline: STATUS_WORK "Resolved-Decline" (UPDATE T_WORK_POLIS `STATUS_WORK = NVL(:2, STATUS_WORK)`; tipe Submit lain
     mengirim NULL = tetap). Menggantikan `[dugaan]` "tetap" (K57-13 lama, K57-7).
  4. Cause of Decline dari data `D:\migrasi\RNM\DDL\CAUSEOFDECLINE.csv` (ID;TERMCONDITION;DESCRIPTION;TERMCONDITION2;
     DESCRIPTION1, 41 baris; DDL tetap tidak ada): satu kueri seluruh baris, urut ID angka; T/C = `DESCRIPTION` kosong
     (GetCauseOfDecline, 14 baris); opsi T/C X = TERMCONDITION baris ber-DESCRIPTION X dan TERMCONDITION2 kosong (POI 4,
     Extension of coverage 4, Risk Exposure 12 termasuk "Others" ID 41) - sumber Pega `Getdescription_ACT` tetap tidak di
     korpus, `[dugaan]` kuat dari bentuk data; teks = TERMCONDITION (DESCRIPTION1 ID 24 "Earthquake prone" berisi "Power
     plant" = salah data); Non - T/C = TERMCONDITION2 terisi (GetCauseOfDecline2, 7 baris). Seluruh teks dipangkas -
     PEMBETULAN: ID 10 "Clauses" tersimpan berakhiran baris baru sehingga syarat Pega `.TERMCONDITION != 'Clauses'` tak pernah
     cocok. POST Decline: description1 T/C ber-opsi wajib salah satu opsinya (400). Ini MENGGANTIKAN "`sebabTc[].opsi` selalu
     `[]`" di K57-14. Instrumen: `TestSusunSebabDecline` atas 41 baris CSV (`services/sebabdecline_data_test.go`, bangkitan) -
     jawaban diketahui dari data; mutasi 8 titik (saring Banding, status Decline, pangkas, urut ID, opsi asing, syarat
     TERMCONDITION2, NVL status, teks MOID) - semuanya menggagalkan uji.
  Tanpa migrasi baru.

## Submit berhasil = case ditutup (work owner 07-10-2026: "harusnya saat berhasil di submit, case nya langsung tertutup")

Sebelumnya layar memuat ulang case di tempat (blok menampilkan "Submit tidak tersedia" karena posisi sudah Team Leader).
Kini (sesi 0f): `InwardFacultative` `onSelesai` → `rute.tsx` menutup case (`setKasus(null)`) dan kembali ke portal dengan
pesan hijau "Case berhasil di-submit. <NBStatus>" (dapat ditutup) - padanan `finishAssignment` Pega. Daftar portal tetap
memuat semua case NB (belum disaring per posisi / workbasket).

## Case yang sudah di-submit dibuka lagi (work owner 07-10-2026, gambar layar "INI KENAPA YA")

Case di posisi Team Leader dibuka dari portal: blok menampilkan form Email Type / Comment yang tidak bisa dipakai + pesan
teknis. Kini (sesi 0f): bila `bolehSubmit` false, form disembunyikan, keterangan `TEKS_EMAIL_TYPE.sudahDiproses` + alasan
backend (teks kecil), grid riwayat tetap. Catatan terbuka: tab detail dan Save for later masih dapat diubah walau case sudah
di Team Leader (Pega: marketing tidak memegang assignment-nya lagi) - menunggu keputusan work owner.

## Case baca-saja di luar tahap marketing (work owner 07-10-2026: "setuju, buatkan jangan halu")

- Layar (sesi 0f): `components/BacaSajaKasus.tsx` membungkus blok General dan isi tab detail; saat `KasusNB.bolehUbah === false`
  semua isian / tombol dinonaktifkan kecuali navigasi untuk melihat data (tombol buka-baris `aria-expanded`, strip tab
  `role="tab"`, paginasi `.pager`, `data-navigasi`); MutationObserver menerapkan ulang pada isian yang baru dirender. Save
  for later nonaktif; banner `TEKS_INWARD.bacaSaja` + alasan backend.
- **K57-21** Satu pembungkus generik, bukan prop baca-saja di tiap komponen tab (tidak ada isian yang terlewat).
- Backend (diminta ke c3): `KasusNB.bolehUbah` / `alasanBacaSaja` dengan aturan yang sama dengan `bolehSubmit`; semua endpoint
  tulis case → 409 bila baca-saja.
- **K57-20 Case baca-saja sesudah keluar dari tahap marketing** (keputusan work owner 07-10-2026 "setuju, buatkan jangan
  halu" atas saran sesi 9d; Pega `finishAssignment`: sesudah Submit marketing tidak lagi memegang assignment). Backend sesi c3:
  - **Satu sumber aturan** `services/bacasaja.go`: `tahapMarketing` (POSITION "Offer" / "ReasFacInMarketing" &&
    FLAG_ONGOING_POLICY "0") + `alasanUbah` (workbasket ReasFacInMarketing lebih dulu → 403, lalu tahap → 409). Dipakai
    `bolehUbah` / `alasanBacaSaja` GET kasus dan jawaban PUT …/general (BARU, lunak: galat baca tahap → baca-saja + log),
    `bolehSubmit` / `alasan` Submit (`alasanTidakBoleh` kini memanggil `alasanUbah`; `ErrSubmitBukanTahapMarketing` =
    `ErrKasusBacaSaja`), dan BolehUbah Scoring (`bolehUbahScoring` = `tahapMarketing`; GET scoring tanpa pelaku → tahap saja).
    Pembaca tahap: `KasusOracle.TahapKasus` (`SELECT POSITION, FLAG_ONGOING_POLICY … WHERE ID = :1 AND LINI = :2`).
  - **Endpoint tulis → 409 / 403**, dicek di layanan sesudah 401 (dan akses IT), SEBELUM isian / basis data / tulis; kasus
    asing 404: PUT general, PUT objek, PUT klausa, PUT spreading, PUT cedant, PUT deduction, PUT payment, PUT scoring, PUT
    loss-record, korespondensi POST / PUT / DELETE / kirim, Submit (sudah). Tab IT-only (Loss Record, Correspondence): syarat
    TAHAP saja (409), tanpa workbasket marketing - aksesnya sudah JOB_POSITION "IT"; mewajibkan keduanya membuat tab hanya
    dapat diubah akun yang memegang IT DAN workbasket marketing `[keputusan agent]` (korpus tidak membedakan; laporkan bila
    work owner mau keduanya).
  - **Tidak dijaga (tidak menyimpan data kasus)**: POST premi, POST akseptasi/langkah, POST …/hitung-coverage,
    …/hitung-net-rate, …/spreading/hitung-share, …/spreading/salin (Copy To All - "TIDAK menyimpan"), …/deduction/hitung,
    …/scoring/hitung, …/loss-record/hitung, …/payment/hitung, …/submit/periksa. Juga tidak: POST opportunity (membuat case
    BARU), POST risk-address dan POST akumulasi (data induk RISKADDRESS / ACCUMULATION, bukan milik satu case).
  - Penjaga aktif bila pembaca tahap terpasang; `DariDasar` selalu memasangnya (`TestDariDasarMemasangTahap`, uji teks
    sumber). Uji: `TestPenjagaBacaSaja` (13 metode tulis × posisi TL / Decline / flag 1 → 409, kasus asing 404, tanpa sesi
    401 lebih dulu, tanpa workbasket 403 kecuali tab IT, tahap marketing lolos), `TestAksesUbahKasus`, `TestBacaSajaHandler`
    (GET 200 + bolehUbah false, PUT general 409 / 403), `TestSQLTahapKasus`; mutasi 5 titik → semuanya tertangkap.
