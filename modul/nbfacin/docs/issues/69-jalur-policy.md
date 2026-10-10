# 69: Jalur Policy — Team Leader & jenjang UW tahap Policy (69a), email Policy + Arasapas (69b)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 08-10-2026 (gambar layar case
> status Policy di `ReasFacInTeamLeader`) "ini pilihan accept reject nya mana? posisinya di ReasFacInTeamLeader". Lanjutan
> K64-1 ("layar Policy baca-saja sampai tiket jalur Policy"). Analisis `CekLimitSpreading_Act`: subagen sesi 0f, klaim inti
> diverifikasi ulang (L2039, L2105, L2195, L2218, L2318, L2333, L2407).

**Status:** 69a kontrak dikirim ke sesi c3 (backend) 08-10-2026; frontend sesi 0f. 69b menyusul. FacRetro = tiket 65.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\Flow\InputInwardFacultativeRISlip.xml`

- Start (Transition1 L18430): FlagOnGoingPolicy 1, IsCedingConfirm "Policy". Decision8 IsLife → Decision9 IsGroup →
  Assignment3 MARKETING; Else → Decision27 IsTBonding → Assignment1 UNDERWRITING FINANCIAL; Else → Utility5
  `CekLimitSpreading_Act` → **Assignment6 TEAM LEADER** (Transition2 L18332, `ReasFacInTeamLeader`).
- Assignment6 flow action `InwardFacultative` (sama dengan TL tahap Offer) → Decision24 `IsUWAccepted`:
  **confirm** → Decision21 `IsFlagUW` (When `pyWorkPage.IsFlagUW = 1`): ya → Utility8 `GetLimitAkseptasi_JUW_UW` →
  Decision33 (ToJUW_A / ToUW / ToSeniorUW / ToDepHeadUW; Else → Decision4); tidak → Utility1 `SaveJsonPolicyFacIn_Act`.
  **reject** → End2. Ask / Decline: **tanpa konektor** (Pega tertahan).
- End2 → flow induk `InputInwardFacultativeOffer.xml` Decision14 (L11806): reject → Assignment9 MARKETING (BINDING)
  (Transition69 L18179: PositionNote ReasFacInMarketing, IsCedingConfirm "Binding", FlagOnGoingPolicy 2); confirm /
  decline → End3.
- Layar: `Section\EmailSection.xml` L3126 TL → radio `.EmailTypeTL` (`DDL\EmailTypeTL.xml` 1 Accept / 2 Reject / 3 Ask /
  7 Decline). UW tahap Policy: `Section\EmailSection_IsUW.xml` L2026 radio `.EmailTypeUWPolicy` (L2865 `FlagOnGoingPolicy
  == 1`); `DDL\EmailTypeUWPolicy.xml` (`ASM-FW-GISFW-WORK!EMAILTYPEUWPOLICY`) 1 Accept / 4 Banding / 7 Decline.
- Keputusan UW tahap Policy: Assignment8 UNDERWRITING → Decision16 confirm → Utility8 (ladder lagi); Assignment14 SENIOR →
  Decision22 confirm → Decision4; Assignment10 DEP HEAD → Decision29 confirm → Decision4; Assignment1 UW FINANCIAL →
  Decision2 confirm → Decision4; decline / reject → End2. Tidak ada konektor `banding`.
- Decision4 FAC OUT? `IsFacRetro` → Decision3 `IsInputFacRetro` (`OfferFacIn.IsInputFacRetro = 1`) → Utility1; Else →
  SubProcess3 OfferFacRetro (tiket 65). Decision4 Else → Utility1 `SaveJsonPolicyFacIn_Act` → Utility7 `SendEmailPolicy` →
  Utility2 `serviceInsertArasapas_act` → Decision30 `IsSuccessHitService` … → END52.

### `IsFlagUW` — `Activity\CekLimitSpreading_Act.xml` (FIRE, bukan polis deklarasi)

- Langkah 4 (syarat `IsFlagUW=="" || IsFlagUW==0`, L2039): coverage ber-AccumulationCode → RDB `GetCountAccumulation`
  (`RDBList\GetCountAccumulation.xml` L86–89, `pooldata.facinproduction` + `pooldata.facinoffer`, case sendiri dikecualikan)
  > 0 → 1 (L1561); spreading `@contains(.TreatyType,"10007")` (L1841) → 1 (L1764).
- Langkah 5 (bukan `IsDeclarationPolicy`, `IsFlagUW==""`) → 1 (L2105) lalu loncat `save`; langkah 6 `IsFlagUW!=""` → `save`
  (L2333). Langkah 7 (`IsFlagUW = 0`, L2407) tidak tercapai untuk non-deklarasi.
- Nilai 0 hanya bila sebelumnya diisi 0: `CopyToAllSpreading_ACT` L666 (syarat `IsFire`), lalu 1 bila `10015` / `SPL` /
  `10007` (lihat `docs/03-celah/01-spreading-capacity-scoring-nb.md` §9.5).

## Keputusan work owner 08-10-2026

- **K69-1** Bertahap: 69a = form TL Policy + IsFlagUW + jenjang UW Policy + selesai Policy (JSON minimal); 69b = email
  Policy (ikut Pega, lewat outbox) + service Arasapas (outbox, endpoint M_LINK_SERVICE / env — koneksi eksternal perlu
  persetujuan). FacRetro = tiket 65.
- **K69-2** TL tahap Policy: **empat opsi seperti Pega** (Accept / Reject / Ask / Decline).
- **K69-3** Ask TL Policy → kembali ke Marketing tahap Binding (sama dengan Reject). *Penyimpangan sadar: Pega tanpa jalur.*
- **K69-4** Decline TL Policy → Resolved-Decline (pola Decline TL tahap Offer). *Penyimpangan sadar: Pega tanpa jalur.*
- **K69-5** `IsFlagUW` dihitung saat TL Accept: 1 bila ada akumulasi berjalan (`GetCountAccumulation` > 0) atau spreading
  ber-TreatyType `10007` / `10015` atau TreatyName memuat `SPL`; selain itu 0. (Setara Pega bila Copy To All spreading sudah
  berjalan; Pega dengan flag kosong = 1.)
- **K69-6** "uw mana ada pilihan banding, banding hanya ada di marketing, dan banding tidak boleh di tahap policy" → UW tahap
  Policy hanya **Accept / Decline** (Banding pada `DDL\EmailTypeUWPolicy.xml` tidak ditampilkan).
- **K69-7** Status akhir Policy (Accept terakhir, tanpa FacRetro) = **Resolved-Completed**.

- **K69-8 (09-10-2026, RALAT aturan tab tahap Policy)** "jika flagongoingpolicy = 1, ini data2nya tidak bisa dirubah lagi.
  karena itu sudah di bind, jadi data gaboleh berubah" (gambar layar tab General case tahap Policy). Di tahap
  `teamLeaderPolicy` dan `uwPolicy` SELURUH data case baca-saja — General, semua tab detail, termasuk Spreading / Accumulation
  (IsUWSpread) dan Following Policy; yang tetap dapat diisi hanya form Email Type + Submit (termasuk isian Deductible
  Acceptance TL dan Validity UW, yang bagian dari form Submit). Menggantikan kalimat 69a-1 "tab boleh diubah pemegang
  workbasket TL seperti TL Offer".
  - **RALAT K69-8 (09-10-2026)** "kalau untuk uw, spreading masih bisa dirubah. add, copy,delete spreading msh bisa di uw":
    di `uwPolicy` tab **Spreading** tetap dapat diubah seperti UW tahap Offer (IsUWSpread: % Share RNM, Add / Copy /
    Delete / ubah baris). **Accumulation TIDAK** (jawaban work owner "Tidak, hanya Spreading"). `teamLeaderPolicy` tetap
    seluruhnya baca-saja.

## Keputusan agent (sesi 0f)

- **A69-1** Decline UW tahap Policy → Resolved-Decline (End2 → Decision14 decline → End3; pola K69-4).
- **A69-2** Case yang di Pega masuk SubProcess3 OfferFacRetro (IsFacRetro dan bukan IsInputFacRetro): Accept ditolak 409
  "Case FacRetro belum dapat diselesaikan di sistem baru (tiket 65)." sampai tiket 65 — mudah dibalik, tanpa status setengah
  jadi. Varian When `IsFacRetro` yang dipakai Decision4 wajib dicek per kelas (OQ-011).
- **A69-3** Tahap baru kontrak `TampilanSubmit.tahap`: `teamLeaderPolicy` (form = TL tahap Offer, opsi K69-2) dan `uwPolicy`
  (opsi K69-6). **RALAT 08-10-2026:** Decline UW Policy TANPA Cause of decline, sama dengan UW tahap Offer (K61: Cause of
  decline hanya ada di `Section\EmailSection.xml`; layar UW = `EmailSection_IsUW.xml`). Versi awal kalimat ini keliru
  mengutip `EmailSection.xml` L5477 — section non-UW.

## Frontend (sesi 0f, 08-10-2026)

- `api.ts` `TampilanSubmit.tahap` + `'teamLeaderPolicy' | 'uwPolicy'`; `labels.ts` `EMAIL_TYPE_UW_POLICY` (1 Accept / 7 Decline,
  Banding dibuang K69-6).
- `components/BlokEmailType.tsx`: `tahapTL` (Offer + Policy: opsi ETL, Deductible Acceptance, Cause of decline wajib) dan
  `tahapUW` (Offer + Policy: Decline tanpa sebab, Reject tanpa modal, Validity Of Acceptance — tampil juga di flag 1 menurut
  verifikasi c3 `EmailSection_IsUW.xml` L9297–L10935).
- Uji `BlokEmailType.test.tsx` "jalur Policy (tiket 69)" (opsi = DDL, isian). tsc bersih, vitest 976 / 976.

## Backend 69a-1 (sesi c3, 08-10-2026)

- **K69a-B1 Tahap** (`services/submitpolicy.go`, `bacasaja.go`, `submit.go`): `teamLeaderPolicy` = POSITION ReasFacInTeamLeader + flag 1
  (tahapDari - tab dapat diubah pemegang workbasket TL, seperti TL Offer; Scoring tetap baca-saja: CallDualScoringRisk "Write" hanya
  flag 0); `uwPolicy` = Underwriting / Senior / DepHead / UW Financial / JUW A + flag 1 (aturan tab = UW Offer: Spreading / Akumulasi
  When IsUWSpread hanya membaca PositionNote, sudah terbuka di flag 1). JUW B flag 1 tetap baca-saja (Decision33 tanpa ToJUW_B).
  `kunciPeriode` true di kedua tahap. Kotak masuk tidak menyaring FLAG (case flag 1 muncul di workbasket posisinya).
- **K69a-B2 Email Type**: TL 1 / 2 / 3 / 7; UW 1 / 7 (lainnya 400). Form TL = TL Offer (Deductible Acceptance, Cause of decline wajib
  untuk Decline); form UW = UW Offer (Comment wajib, tanpa Cause of decline - RALAT 9d atas A69-3, Validity Of Acceptance).
- **K69a-B3 Validity Of Acceptance di flag 1** `[terverifikasi]`: `Section\EmailSection_IsUW.xml` L9297-L10935 (DaysValidity,
  DateValidity, WaitingBindDate) - nol syarat Visible / ReadOnly pada leluhurnya (hanya radio `.EmailTypeUWPolicy` yang bersyarat
  flag 1) -> TAMPIL di flag 1; `POST …/submit/validitas` dan Submit uwPolicy menghitung / menyimpannya seperti UW Offer.
- **K69a-B4 Accept TL** (Decision24 confirm -> Decision21): IsFlagUW (K69-5) = spreading TreatyType memuat "10007" / "10015" atau
  TreatyName memuat "SPL", ATAU GetCountAccumulation > 0 per ACCUMULATIONCODE coverage (`repository/policy.go`, SQL verbatim
  RDBList L86-89; bind = kode, StartDateTime `dd/MM/yyyy` WIB, pzInsKey untuk FACINPRODUCTION, pyID untuk FACINOFFER -
  CekLimitSpreading_Act langkah 3 / 4.1.1.1). 1 -> GetLimitAkseptasi_JUW_UW dengan jabatan **"LEADER"** (langkah 7
  `[terverifikasi]`; langkah 8 IsGroup -> keluar, LetterNo "") -> Decision33; 0 -> selesai Policy (409).
- **K69a-B5 Accept UW**: Underwriting / JUW A -> Utility8 (jabatan UNDERWRITER / JUW_A) -> Decision33; Senior / DepHead / UW
  Financial -> Decision4. Decision33 ToUW / ToJUW_A NB status verbatim; ToSeniorUW / ToDepHeadUW men-SET NB status berisi NAMA ORANG
  di korpus - diganti teks jabatan statusInboxPosisi (pola tiket 61 D7). Flag tetap 1.
- **K69a-B6 Decision4** (FAC OUT?): `adaFacRetro` (spreading FACOUT 10015) -> 409 "Case FacRetro belum dapat diselesaikan di sistem
  baru (tiket 65)."; Else / IsFlagUW 0 / Decision33 Else -> 409 "Penyelesaian Policy belum tersedia di sistem baru (tiket 69a-2)."
  (juga di POST …/submit/periksa). Tanpa tulisan (nol status setengah jadi). Kurs tidak ada -> 422 seperti UW Offer.
- **IsFacRetro (OQ-011)** `[terverifikasi]`: varian kelas Work `ASM-FW-GISFW-WORK!ISFACRETRO` hanya ada di ekspor Endorsment Fac In
  (`pyWorkPage.OfferFacIn.IsFacRetro = 1`). NB FacIn hanya memuat `ASM-FW-GISFW-DATA-OFFERFACIN!ISFACRETRO` (juga RNW Fac In) yang
  `pyConditionString`-nya "pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInGroupLeader" - teks label USANG;
  ekspresi yang dievaluasi `pyConditionValue1` = `compareTwoValues(.IsFacRetro, "=", 1)`. Semantik sama -> `adaFacRetro`.
- **K69a-B7 Reject / Ask TL** (K69-3, Transition69): POSITION ReasFacInMarketing, FLAG 2, NB "NB IS IN " + MarketingName + "'S INBOX"
  (nama dari data); CONFIRM_BINDING / RECEIVED_RI_SLIP tidak direset. **Decline** TL / UW (K69-4 / A69-1): Decline, Resolved-Decline.
- **K69a-B8 Riwayat** (SaveViewSuggest L793; SetValidateDate_PostAct 3.5 / SetValidateDateUW_PostAct 5.4): POSISI "Policy",
  PUTARAN "2"; APPROVAL = label pilihan. ⚠️ Penyimpangan sadar: UW tahap Policy Pega mengisi Approval dari `EmailTypeUWEDM` (bukan
  radio `EmailTypeUWPolicy`) - kemungkinan kosong; di sini label pilihan. TL Policy TIDAK menulis DATE_TO_UW / SUBMIT_TO_UW (hanya
  TL Offer -> JUW B); isian Deductible TL tetap disimpan.
- **K69a-B9 FACINPERFORMANCE** (InsertFacInPerformance_Act langkah 5 / 9 / 10 / 13): TL - tidak ada (aktivitas tidak membaca
  EmailTypeTL); UW - POSITION_DOC "POLICY", Accept ACCEPT `[dugaan kuat]` (langkah 9 membaca `EmailTypeUW` = sisa Accept UW tahap
  Offer, bukan radio Policy), Decline DECLINE (`EmailTypeUWPolicy == "7"`); tanpa OUTSTANDING (langkah 13 dilewati flag 1).

**Temuan (tidak diubah):** `SetValidateDateUW_PostAct` langkah 10 `CALL UpdateJSONOffer_ACT` (-> SaveOfferProduction_Act: tulis
JSON_OFFER / FACINOFFER) berjalan di setiap Submit UW (Offer dan Policy) - tidak diport (ADR-0028 / ADR-0043), dilaporkan untuk
69a-2. GetLimitAkseptasi_JUW_UW langkah 8 (IsGroup -> keluar) juga berlaku di UW tahap Offer (tiket 61) - port tiket 61 belum
menerapkannya `[pertanyaan terbuka]` (69a-1 menerapkannya untuk jalur Policy).

**Uji** `TestTahapPolicy`, `TestIsFlagUW`, `TestAjukanSubmitTLPolicy` (jenjang LEADER, IsFlagUW 0, tanpa LEADER, Else FacRetro, TSI
di bawah jenjang, group, Reject / Ask, Decline + sebab wajib, proteksi 422, marketing 403, Banding 400), `TestAjukanSubmitUWPolicy`,
`TestDeductibleTLPolicy`, repository `TestSQLHitungAkumulasi`, handler `TestGalatPolicy409`; diralat `TestPenjagaBacaSaja` /
`TestAksesUbahKasus` (TL flag 1 = teamLeaderPolicy), `TestPeriksaDanTampilanSubmitUW` (UW flag 1 = uwPolicy). **Mutasi 27 / 27 mati.**

## Analisis 69a-2 — tulisan DB SaveJsonPolicyFacIn_Act jalur NB FIRE (sesi c3, 08-10-2026; ANALISIS SAJA, belum diimplementasi)

Sumber: `D:\migrasi\RNM\NB FacIn\` (Activity / RDBList / When) dan `D:\migrasi\RNM\DDL\`. Inventaris disusun subagen read-only. Klaim
yang ditandai ✔ dicek ulang sesi c3 ke XML/DDL mentah. Cakupannya hanya **NB FIRE**: cabang Aneka / Golf / Marine Cargo / MBU / PA di
SaveTreatyProduction_Act langkah 19-24 hanya dipindai pada tingkat langkah.

### Urutan tulisan (NB FIRE)

| # | Langkah | Target | CALL / COMMIT | Kondisi |
| --- | --- | --- | --- | --- |
| 1 | SaveJsonPolicy 11, `GetSequenceNumber_SQL` | `GENERATE_SEQUENCE_NUMBER` (lewat proc) | CALL `PROC_GENERATE_SEQUENCE_NUMBER` + COMMIT | `PolicyNo == ""` |
| 2 | SaveJsonPolicy 15, `INSERTJSON_JSONPOLIS_FACIN` | `JSON_POLIS` (PK IDPEGA ✔, UNIQUE NOPOLIS+IDPEGA ✔) | CALL `INSERTJSONPOLIS` + COMMIT (proc menelan galat) | selalu |
| 3 | InsertCedingProduction 1 | `CEDING_FACINPRODUCTION` DELETE per IDPEGA | RDB-Delete | selalu |
| 4 | InsertCedingProduction 2.6, per ceding | `CEDING_FACINPRODUCTION` INSERT 6 kolom + **UPDATE tanpa WHERE** ✔ (`INSERTUPDATECEDINGPRODUCTION.txt` 44-47) | CALL + COMMIT | per baris CedingCedantList |
| 5 | SaveJsonPolicy 17.2 | `JSON_POLIS_ERROR` | PL/SQL **tanpa `END;`** ✔ | `HASIL1` memuat "ERR" (HASIL1 sudah ditimpa langkah 4) |
| 6 | SaveJsonPolicy 21 | `M_RISK_LOSS_PROFILE` INSERT | tanpa COMMIT | selalu (simpan ulang -> ganda) |
| 7 | SaveTreatyProduction 3-4, `CekFacin_Sql` ✔ | (baca FACINPRODUCTION per IDPEGA) | — | IDPEGA sudah ada -> ExitActivity |
| 8 | SaveFacinProdFireNB 13.3.2.2.1.3.2 / 13.3.2.3.1.2 ✔ | **FACINPRODUCTION** INSERT 79 kolom per lokasi × item × coverage × [layer basis 5] × spreading | COMMIT per baris ✔ | penjaga master polis (langkah 23 `PolicyType!="1"`) **nonaktif** ✔ |
| 9 | SaveJsonPolicy 27, `UpdateFacProdObjItemID_SQL` | FACINPRODUCTION UPDATE OBJECT_NAME_ID | tanpa COMMIT | `where NOPOLIS = {InputData.CARI1}` = pzInsKey ✔ -> `[dugaan]` 0 baris |

- `[terverifikasi]` **FACINPRODUCTION DITULIS** (langkah 8). Ia dibaca GetCountAccumulation (IsFlagUW, tiket 69a-1), jadi tanpa tulisan
  ini akumulasi case lain tidak melihat polis yang selesai. DDL tanpa PK / UNIQUE / INDEX ✔; `LAYERNO NOT NULL` ✔ diisi teks "NULL".
- `[terverifikasi]` TREATYPRODUCTION tidak ditulis. Namanya menyesatkan: `InsertTreatyProduction_Sql` menulis FACINPRODUCTION.
- Bacaan: `GETTanggalClosing_SQL`, `GetOldIDBusiness_SQL`, `GetKodeProdNonLife_SQL` (`KODE_PRODUKSI TYPE='NONLIFE'`), `GetNopolis_Sql`
  (json_polis), `GetMOID_Sql`.

### Nomor polis dan ProdDateTime

- ProdDateTime (SaveJsonPolicy 5): kosong -> sekarang. Bila Begin date lebih besar dari hari ini -> StartDateTime (`[dugaan kuat]`
  arah `DateTimeDifference` dari deskripsi langkah). Hari > TANGGAL_CLOSING (kosong -> 25) -> tanggal 01 bulan berikut, 05:00 GMT.
- PolicyNo = `KODE(NONLIFE) + "F" + BusinessOldId + "." + MM.YYYY + "." + NNNNN`. Penghitung (CLASS "ASM-FW-GISFW-Work-NB/RNW", JENIS
  KODE+"F", TAHUN) disetel ulang per tahun. Cabang cutover `<= 02/01/2026` -> `12.2025` ✔.
- `inti/backend/penomor` (ADR-0043) sudah menyediakan HariClosing, HitungPeriodeNomor (ADD_MONTHS + tahun, sama dengan proc),
  UrutNomorBerikut (FOR UPDATE, tanpa COMMIT) dan AwalanProduksi(tipe). Dipakai dengan tipe "NONLIFE", tanpa mengubah inti.

### Kejanggalan korpus

Kunci UPDATE salah (9); UPDATE ceding tanpa WHERE (4); pemeriksaan "ERR" membaca HASIL1 yang sudah ditimpa dan blok tanpa `END;` (5) —
JSON_POLIS_ERROR `[dugaan]` tidak pernah terisi; Begin / End / spread date diformat `hh` 12 jam ✔ lalu diurai `HH24`; FireNB menimpa
override MOID (QuotationData.MOID); NOENDORS tidak pernah sampai; SAILDATE / DEDUCTION2 / TYPE_EDM tidak diisi di jalur FIRE; galat
INSERT FACINPRODUCTION ditelan; COMMIT per baris + penjaga per IDPEGA -> baris yang gagal tidak pernah dilengkapi; JSON_POLIS PK
IDPEGA -> simpan ulang gagal diam-diam. DDL GENERATE_SEQUENCE_NUMBER / TANGGAL_CLOSING / KODE_PRODUKSI / M_RISK_LOSS_PROFILE /
JSON_POLIS_ERROR tidak ada di korpus.

### Usulan `[keputusan agent]` — menunggu persetujuan work owner

- **U1 Satu transaksi**: semua tulisan di satu `db.Tx`, nol CALL (ADR-0043), nol COMMIT (ADR-U-0029); gagal -> rollback + 500.
  Menyimpang dari COMMIT per langkah di Pega.
- **U2 Nomor polis di aplikasi** lewat `inti/backend/penomor`, hanya bila PolicyNo kosong. Menulis GENERATE_SEQUENCE_NUMBER
  (ADR-0043 Akibat 2). ProdDateTime diport seperti di atas.
- **U3 JSON_POLIS**: INSERT langsung 9 kolom proc (PRODKE '0', NOENDORS NULL, TGL_PROD ProdDateTime). DATA_JSON minimal (ADR-0028,
  seperti keputusan B JSON_OFFER). Wajib, karena GetNopolis dan PRODKE ceding membacanya. Bila IDPEGA sudah ada: lewati, idempoten.
- **U4 CEDING_FACINPRODUCTION**: DELETE per IDPEGA, lalu INSERT per ceding dengan PRODKE / NOPOLIS_ARAS / NOENDORS_ARAS langsung dihitung.
  ⚠️ Menyimpang: tanpa UPDATE seluruh tabel.
- **U5 JSON_POLIS_ERROR**: tidak diport (blok tidak sah, kondisinya praktis tak pernah benar); galat = rollback.
- **U6 M_RISK_LOSS_PROFILE**: INSERT seperti Pega. Sumber CARI14-18 belum dirinci, tipe kolom belum terverifikasi (DDL tidak ada).
- **U7 FACINPRODUCTION**: INSERT 79 kolom per baris spreading seperti FireNB, dengan penjaga per IDPEGA. Uang apd (ricomm / brokerage =
  PremiumSpreaded × pct / 100; pembulatan belum terverifikasi).
- **U8 UpdateFacProdObjItemID**: tidak diport (kunci pzInsKey -> 0 baris); OBJECT_NAME_ID tetap ItemTypeID seperti hasil nyata Pega.
- **U9 Status akhir** Resolved-Completed (K69-7) di transaksi yang sama. Email Policy = 69b.

### Pertanyaan terbuka untuk work owner

- **P1** Jam Begin / End / spread date: tiru `hh` (siang tersimpan pagi) atau benar `HH24`?
- **P2** MOID: ikut FireNB (selalu GetMOID) atau override QuotationData.MOID (SaveTreatyProduction 11)?
- **P3** Teks "NULL" di LICENSEPLATE / ENGINENO / NOSERIAL / LAYERNO: tiru apa adanya?
- **P4** Master polis (PolicyType "1"): tetap masuk FACINPRODUCTION (perilaku nyata; penjaga nonaktif) atau tidak (REMARK langkah 22)?
- **P5** Cakupan bisnis: nbfacin 69a-2 hanya FIRE? Cabang Aneka / Golf / Marine / MBU / PA belum diinventaris.
- **P6** Simpan ulang (FACINPRODUCTION sudah ada): lewati seperti Pega?
- **P7** M_RISK_LOSS_PROFILE: tulis setiap kali (ganda) atau sekali per polis?
- **P8** DBA: DDL / tipe kolom GENERATE_SEQUENCE_NUMBER, TANGGAL_CLOSING, KODE_PRODUKSI, M_RISK_LOSS_PROFILE, PK FACINPRODUCTION.

### Keputusan work owner 08-10-2026 atas analisis 69a-2 (lewat sesi 9d) — "SILAKAN IMPLEMENTASI"

- **K69a2-1** U1-U9 disetujui apa adanya: satu transaksi, nol CALL, nol COMMIT; nomor polis lewat `inti/backend/penomor`; JSON_POLIS
  minimal; CEDING_FACINPRODUCTION tanpa UPDATE seluruh tabel; M_RISK_LOSS_PROFILE; FACINPRODUCTION 79 kolom; UpdateFacProdObjItemID dan
  JSON_POLIS_ERROR tidak diport; Resolved-Completed di transaksi yang sama.
- **K69a2-2** (P1) Jam Begin / End / spread date ditulis BENAR 24 jam (HH24). Ini penyimpangan sadar dari cacat `hh` Pega
  (SaveFacinProdFireNB_Act L1976 / L2022 `"dd/MM/yyyy hh:mm:ss"`, dicek juga oleh 9d).
- **K69a2-3** (P2) MOID ikut FIRE NB Pega: GetMOID, override QuotationData.MOID diabaikan.
- **K69a2-4** (P3) Teks "NULL" di LICENSEPLATE / ENGINENO / NOSERIAL / LAYERNO ditiru dari Pega.
- **K69a2-5** (P4) Polis master (PolicyType "1") ikut perilaku nyata Pega: tetap masuk FACINPRODUCTION.
- **K69a2-6** (P5) Cakupan FIRE dulu. Case non-FIRE yang sampai penyelesaian Policy -> 409 "Penyelesaian Policy untuk jenis bisnis ini
  belum tersedia di sistem baru (tiket 69a-2).".
- **K69a2-7** (P6) Simpan ulang dilewati seperti Pega (CekFacin): IDPEGA sudah ada -> FACINPRODUCTION tidak ditulis ulang.
- **K69a2-8** (P7) M_RISK_LOSS_PROFILE sekali per polis. Ini penyimpangan sadar dari baris ganda Pega.
- **K69a2-9** (P8) Item DBA, belum dijawab: DDL GENERATE_SEQUENCE_NUMBER / TANGGAL_CLOSING / KODE_PRODUKSI / M_RISK_LOSS_PROFILE,
  PK / indeks FACINPRODUCTION, dan grant yang dibutuhkan skema aplikasi:
  - SELECT: TANGGAL_CLOSING, KODE_PRODUKSI, JSON_POLIS, FACINPRODUCTION, sumber OldIDBusiness / MOID;
  - SELECT FOR UPDATE / INSERT / UPDATE: GENERATE_SEQUENCE_NUMBER;
  - INSERT: JSON_POLIS, FACINPRODUCTION, M_RISK_LOSS_PROFILE;
  - DELETE / INSERT: CEDING_FACINPRODUCTION.

## Backend 69a-2 — selesai Policy (sesi c3, 08-10-2026)

Berkas: `services/selesaipolicy.go`, `repository/selesaipolicy.go`, `models/policy.go`. Penyambungan di `services/submitpolicy.go`
(`decision4Policy` / IsFlagUW 0 -> `tujuanUW.selesai`, `lengkapiTulisanPolicy`) dan `repository/submit.go` (`TulisSubmit` ->
`tulisPolicy`). Semuanya satu transaksi Submit, nol CALL, nol COMMIT.

- **K69a2-B1 Jalur**: TL Accept dengan IsFlagUW 0 (Decision21 Else), dan Decision4 Else (TL / UW / JUW A sesudah Decision33 Else;
  Senior / DepHead / UW Financial), menjadi selesai Policy. FacRetro tetap 409 (tiket 65). Tujuan Utility1 tidak men-SET
  PositionNote / NBStatus `[terverifikasi]` (rislip.json: Transition3 / 91 / 76 / 87 tanpa SET kolom). Karena itu POSITION dan
  NB_STATUS tetap (`NBStatusTetap`), STATUS_WORK menjadi "Resolved-Completed" (K69-7). Riwayat / FACINPERFORMANCE UW tetap seperti
  69a-1. Email Policy + Arasapas = 69b.
- **K69a2-B2 Case selesai baca-saja** `[keputusan agent]`: karena POSITION tidak berubah, `TahapKasus` / `KonteksSubmit` membaca
  STATUS_WORK. `kasusSelesai` ("Resolved-…") menutup tahapDari / tahapUW / tahap Policy; alasan "Case baca-saja: case sudah
  selesai (Resolved-Completed)." (409). Penjaga serentak `sqlPosisiSubmit` menolak case Resolved. Kotak masuk sudah menyaring
  Resolved-.
- **K69a2-B3 Nomor polis**: JSON_POLIS case sudah ada -> NOPOLIS / TGL_PROD dipakai ulang. Belum ada -> TANGGAL_CLOSING lewat
  `penomor.HariClosing` (kosong -> 25, langkah 4) -> `models.ProdDateTime` (langkah 5) -> `penomor.HitungPeriodeNomor` ->
  `AwalanProduksi(tx, "NONLIFE")` -> `UrutNomorBerikut(CLASS "ASM-FW-GISFW-Work-NB/RNW", JENIS KODE+"F", TAHUN)` ->
  `KODE F + BusinessOldId + "." + MM.YYYY + "." + 5 digit`. Hari closing yang sama dipakai prosedur. ProdDateTime sudah bergeser
  ke tanggal 01 bila lewat closing, sehingga pergeseran ADD_MONTHS prosedur tidak menyala lagi (setara dengan prosedur yang
  membaca closing NULL). `[keputusan agent]`:
  - urut > 5 digit tidak dipotong (seperti RakitNomorKlaim);
  - BusinessOldId tak diketahui (BUSINESS 0 / > 1 baris / kosong) -> 409 `ErrBusinessOldIDTidakAda` (Pega memakai pxResults(1)
    / kosong);
  - "Begin date lebih besar" = Begin date >= 24 jam di depan (arah DateTimeDifference `[dugaan kuat]` dari deskripsi langkah
    5.2, @toInt memotong pecahan hari).
- **K69a2-B4 JSON_POLIS**: INSERT 9 kolom prosedur. PRODKE '0', NOENDORS / OLDNOPOLIS NULL (OldPolicyNo = jalur renewal, tidak ada
  di NB; BUKAN .Following), TGL_PROD = ProdDateTime 24 jam, USERNAME = akun. DATA_JSON minimal = kunci yang DIBACA SQL korpus
  atas JSON_POLIS (pola keputusan D1 JSON_OFFER). Sensus sesi c3: pindai seluruh RDBList yang menyebut json_polis -> 7 jalur
  titik (`PPnCheck`, `PolicyData.StartDateTime` / `EndDateTime`, `QuotationData.BusinessType` / `CedingCoName` / `EdmDate` /
  `EdmType`) dan 11 jalur JSON_TABLE (`LocationList[*]` > `Property.RiskLocation.ASMAddress`, `Property.PropertyItemList[*]` >
  `Currency`, `CoverageList[*]` > `AccumulationCode` / `CoverageNote`, `SpreadingList[*]` > `TreatyType` / `TSISpreaded` /
  `PremiumSpreaded`); susunan sarang dicek ke `GetSummaryRiskAccumPolis_Sql`. ⚠️ Temuan: saran akumulasi modul ini sendiri
  (`repository/akumulasi.go sqlAkumulasiPolis`) membaca `$.LocationList…AccumulationCode` JSON_POLIS. Dokumen yang hanya berisi
  IDNewBisnis akan mengosongkan saran itu untuk polis sistem baru. Pembaca lain ada di korpus Endorsment Fac In / RNW Fac In
  (modul lain).
- **K69a2-B5 CEDING_FACINPRODUCTION**: DELETE per IDPEGA, INSERT per CedingCedantList. CEDINGID / CEDINGNAME tanpa ";". PCT_SHARE
  mengikuti langkah 2.2-2.5 berurutan (GROSS ShareCeding / PercentShare × 100; NET / kosong ShareCeding; kosong + ShareCeding
  kosong 100 / n skala 4; kosong + satu baris 100). PRODKE dari JSON_POLIS, NOPOLIS_ARAS = REPLACE(NOPOLIS, '.', ''),
  NOENDORS_ARAS NULL. Baris kembar IDPEGA + CEDINGID + PCT_SHARE dilewati (prosedur). Tanpa UPDATE seluruh tabel (K69a2-1).
- **K69a2-B6 M_RISK_LOSS_PROFILE**: langkah 19-20 apa adanya — TSI / Premi Nusantara Re dari coverage TERAKHIR; PremiumQs
  dikosongkan per coverage; PREMIUMQS = QS terakhir × 0.8; PREMIUMOR = PremiumQs coverage itu (kosong = 0 `[dugaan]`) × 0.2 + OR
  terakhir; PREMIUMFACOUT = FACOUT terakhir; COB = Class of Business. Ditulis sekali per NOPOLIS (K69a2-8).
- **K69a2-B7 FACINPRODUCTION**: sensus dua cara cocok (daftar kolom korpus: pisah koma 79, token 79; VALUES 79, placeholder 79;
  daftar Go 79; A\B = B\A = ∅, urutan sama). 59 bind = 79 - 16 NULL - 3 'NULL' - PRORATE 1. Satu baris per lokasi × item ×
  coverage × spreading (basis != 5). Coverage basis "5" mengulang LayerList yang tidak ada di sistem baru (butir 30) -> tanpa baris,
  dicatat log. CekFacin: IDPEGA sudah ada -> dilewati (K69a2-7). Pemetaan:
  - BEGIN / END / SPREAD_DATE 24 jam WIB (K69a2-2; SPREAD_DATE = ProdDateTime);
  - MARKETINGOFFICERCODE = GetMOID_Sql (`CLIENTID = MarketingCode AND MOSTATUS '1'`, MIN(ID) `[keputusan agent]`, K69a2-3);
  - teks "NULL" di LICENSEPLATE / ENGINENO / NOSERIAL / LAYERNO (K69a2-4); PRORATE 1;
  - GROUP_TYPE "1" / "2" (pembuat group, seperti FACINOFFER);
  - OCCUPATIONCODE `baru;lama` (pemisah ";", berbeda dengan FACINOFFER " || ");
  - IDX_OBJECT_ITEM = urutan item 1-based;
  - RICOMM / BROKERAGE = PremiumSpreaded × persen / 100 (persen kosong = 0 `[dugaan]`); PERCENT_RI_COMM kosong -> 0; LOL =
    ClaimEstimation (kosong 0); TSI100 = TSIGrossSpreaded; STATUS_BUSINESS "1" (NB, seperti kasusTangga);
  - CEDINGID = kode Ceding Co gabungan ";";
  - NULL: NOENDORS (tak diteruskan), SAILDATE / TYPE_EDM / DEDUCTION2_* (tak di-set di jalur FIRE), BINDER / ISB2B / ASKRED /
    NO_MASTER_POLIS (halaman tidak ada di sistem baru).

  Semua angka dibulatkan 8 desimal sebelum bind (mask `angkaMasuk`). ⚠️ PCT_SHARE / PERCENT_* NUMBER(23,20) kehilangan digit
  ke-9 sampai ke-20 dibanding Pega.
- **K69a2-B8 Tidak diport**: JSON_POLIS_ERROR, UpdateFacProdObjItemID_SQL (K69a2-1). K69a2-6 (non-FIRE -> 409) tanpa gerbang
  kode: BusinessType tidak tersimpan di sistem baru dan modul hanya memodelkan objek FIRE (kasusTangga "Fire", tiket 61). Case
  non-FIRE belum dapat lahir.
- **Item DBA (K69a2-9)**: lihat daftar di atas. Grant tambahan yang dibutuhkan: SELECT MARKETINGOFFICER (sudah dipakai).

**Uji**:
- models: `TestProdDateTime` (11 kasus: closing, zona Jakarta, Desember, Begin date < / = / > 24 jam), `TestRakitNomorPolis`;
- services: `TestSelesaiPolicyTulisan` (33 kolom teks + 12 angka satu baris, ceding, profil, dokumen JSON),
  `TestSelesaiPolicyDitolak` (OldID 0 / ganda / kosong 409, tidak sejalan), `TestKasusSelesaiBacaSaja` (+ instrumen),
  `TestBarisCedingProduksi` (8 kasus), `TestProfilRisikoLoss`, `TestKodeOkupasiProduksi`, `TestBarisFacInProductionBasisLayer`;
- repository: `TestSQLSelesaiPolicy` (79 kolom, 59 bind, teks tetap, GetMOID, JSON_POLIS / CEDING tanpa UPDATE / profil,
  penjaga Resolved);
- diralat: `TestAjukanSubmitTLPolicy` / `TestAjukanSubmitUWPolicy` (409 selesai -> selesai Policy), handler `TestGalatPolicy409`
  (ErrBusinessOldIDTidakAda), `TestSQLTahapKasus` / `TestSQLBinding` / `TestSQLSubmit` (STATUS_WORK).

**Mutasi 39 / 39 mati** (empat ditulis ulang supaya terkompilasi; satu hidup "LPAD tepat 4 digit" -> kasus uji 1234 ditambah).
⚠️ `tulisPolicy` (urutan baca / tulis dalam transaksi) hanya teruji lewat teks SQL — sama dengan `tulisBinding` tiket 63; uji
Oracle nyata menunggu DEV.

### Keputusan user 09-10-2026 — service Arasapas DI-HOLD

- **K69b-1** Pertanyaan user (sesi c3): *"jadi kalau di accept di TL, ini jadi polis trus servisarasapas dijalankan ga ya? utk
  servis ini harusnya di hold dlu"*. Utility2 "HIT SERVICE ARASAPAS" (sesudah Utility7 Send Email Policy) **TIDAK dijalankan**
  sampai ada keputusan lanjutan. `[terverifikasi]` per 09-10-2026: jalur selesai Policy 69a-2 tidak memanggil Arasapas dan tidak
  menulis outbox (grep `arasapas|outbox|LOG_SERVICE` di nbfacin hanya mengenai komentar dan lampiran tiket 67). Tiket 69b tidak
  boleh menyambungkannya tanpa keputusan baru.

## Backend K69-8 + RALAT (sesi c3, 09-10-2026)

- `services/bacasaja.go`:
  - `tahapDari` tidak lagi membuka `teamLeaderPolicy`;
  - `alasanUbah` -> flag 1 = "Case baca-saja: tahap Policy (sudah di-bind)." / `ErrTahapPolicyBacaSaja` (409, membungkus
    ErrKasusBacaSaja). Case Resolved-* tetap memakai alasan "case sudah selesai" (diperiksa lebih dulu).
- `services/submit.go`: `alasanTidakBoleh` / `tahapKonteks` memakai `tahapTLPolicy` langsung. Form Submit TL Policy tetap
  (`TampilanSubmit.tahap` "teamLeaderPolicy", bolehSubmit pemegang workbasket); isian Deductible / Validity tersimpan lewat Submit.
- `services/akumulasiuw.go`: flag 1 -> `ErrTahapPolicyBacaSaja` sebelum IsUWSpread (RALAT: "Tidak, hanya Spreading").
- Spreading uwPolicy TETAP terbuka seperti UW Offer (`alasanUbahSpreading` / `wajibTahapUW` memakai IsUWSpread lebih dulu).
  TL Policy spreading 409.
- GET kasus flag 1: bolehUbah false + alasan di atas; ubahAkumulasi false; bolehFollowing false; kunciPeriode true;
  bolehUbahSpreading / ubahBarisSpreading = true hanya uwPolicy pemegang workbasket.
- Rute tulis yang terkunci di flag 1 (dicek dari daftar rute handler):
  - general, objek, objek/csv, klausa, cedant, deduction, payment, scoring, loss-record, korespondensi (POST / PUT / DELETE /
    kirim) — `TestPenjagaBacaSaja` (14 metode, UJI-TL1);
  - akumulasi, following (2) — `TestK698TahapPolicyBacaSaja` / tiket 70;
  - spreading (+ salin) hanya TL Policy;
  - endpoint hitung* (tanpa tulisan) tidak dijaga, seperti sebelumnya.
- **Lampiran (tiket 67) — dilaporkan, TIDAK diubah**:
  - `bolehUnggahLampiran` = K67-2 "Upload di semua tahap oleh pemegang workbasket posisi case" -> unggah dan ulangi TETAP
    boleh di flag 1;
  - Delete (`bolehHapusLampiran`) = FlagOnGoingPolicy 0 atau IT (Pega b4605) -> di flag 1 hanya IT;
  - ⚠️ case Resolved-Completed (posisi tetap) juga masih dapat diunggahi oleh pemegang workbasket posisi itu.
- Uji: `TestK698TahapPolicyBacaSaja` (TL / UW flag 1 + instrumen flag 0, TampilanSubmit TL Policy); diralat `TestPenjagaBacaSaja`,
  `TestAksesUbahKasus`, `TestTampilanSubmit` (flag 1), `TestTahapPolicy`, `TestKasusSelesaiBacaSaja`. Mutasi 6 / 6 mati.

## RALAT 69a-2 (sesi c3, 09-10-2026) — keputusan work owner lewat sesi 9d

- **DATA_JSON JSON_POLIS** — work owner: *"sistem baru tidak pakai data_json lgi."* K69a2-B4 (isi 18 jalur kunci SQL korpus)
  **DIBATALKAN**. DATA_JSON = `{"IDNewBisnis": caseId}` (pola JSON_OFFER tiket 63 D1; kolom CLOB nullable `[terverifikasi]` DDL,
  tetapi dokumen minimal dipilih supaya baris tetap dapat dikenali). Baris JSON_POLIS (NOPOLIS / PRODKE / TGL_PROD / USERNAME)
  tetap ditulis (nomor polis dipakai ulang, PRODKE ceding).
- **Saran akumulasi per nomor polis** (`repository/akumulasi.go sqlAkumulasiPolis`, layar Accumulation langkah 6.6) pindah dari
  JSON_TABLE atas JSON_POLIS.DATA_JSON ke TABEL `POOLDATA.FACINPRODUCTION.ACCUMULATIONCODE` per NOPOLIS (sumber yang sama dengan
  GetCountAccumulation; berlaku juga untuk polis lama Pega). ⚠️ Beda: FACINPRODUCTION hanya memuat coverage yang punya baris
  spreading (coverage tanpa spreading / basis 5 tanpa LayerList tidak ada); JSON_TABLE asal (NESTED PATH) juga memungut coverage
  tanpa spreading. Baris endorsement (NOENDORS terisi) ikut, seperti semua baris JSON_POLIS ber-NOPOLIS itu.
- Pembaca DATA_JSON lain di modul ini: hanya `services/loader` — alat migrasi data yang membaca JSON polis LAMA Pega menjadi tabel
  flat. Bukan pembaca runtime; tidak diubah.
- **nomorPolis di jawaban Submit** — work owner: YA. `HasilSubmit.nomorPolis` (kabel selalu ada): terisi hanya bila Submit
  menyelesaikan Policy (nomor baru atau NOPOLIS JSON_POLIS yang dipakai ulang; diisi repository di dalam transaksi lewat
  `TulisanPolicy.NomorPolis`), selain itu "".
- **Lampiran di flag 1 / Resolved** — work owner atas laporan K69-8: *"Biarkan seperti sekarang"*. Unggah / Ulangi tetap boleh
  pemegang workbasket posisi (K67-2) di flag 1 dan di case Resolved; Hapus tetap flag 0 / IT.
- Uji:
  - `TestSQLAkumulasiSQL` (FACINPRODUCTION, nol JSON, 2 bind), `TestSelesaiPolicyTulisan` (dokumen minimal + nomorPolis),
    `TestHasilSubmitNomorPolis`, jawaban Submit lain memuat `"nomorPolis":""`;
  - ⚠️ `TestHasilSubmitKunciTS` MERAH sampai api.ts `HasilSubmit` memuat `nomorPolis` (sesi 0f);
  - mutasi 5 / 5 sesudah uji kabel ditambah, KECUALI nama tabel di `AkumulasiOracle.AkumulasiPolis` (hanya teruji lewat Oracle
    nyata, seperti pembaca repository lain).

### BUG DEV 09-10-2026 — CEDING_FACINPRODUCTION.NOENDORS_ARAS (sesi c3)

- Accept TL Policy NB-151029 di DEV: `repository: CEDING_FACINPRODUCTION: ORA-00904: "NOENDORS_ARAS": invalid identifier`. Transaksi
  rollback; tidak ada tulisan setengah jadi.
- Sebab: DDL korpus `DDL\CEDING_FACINPRODUCTION.txt` baris 11 memuat `NOENDORS_ARAS`, skema DEV tidak.
- Perbaikan `[keputusan agent 0f]`: INSERT `sqlSisipCedingProduksi` tidak lagi menulis NOENDORS_ARAS. Nilainya selalu NULL
  (NOENDORS NULL -> REPLACE(NULL) = NULL), jadi setara Pega dan aman di skema yang punya maupun tidak punya kolom itu. Meralat
  K69a2-B5 ("NOENDORS_ARAS NULL").
- Uji `TestSQLSelesaiPolicy`: daftar kolom ceding persis 8 kolom, nol `NOENDORS_ARAS`, 7 bind.
- Selisih kolom lain (CEDING_FACINPRODUCTION / JSON_POLIS / FACINPRODUCTION / M_RISK_LOSS_PROFILE / GENERATE_SEQUENCE_NUMBER)
  menunggu hasil `all_tab_columns` DEV dari work owner; tidak ditebak.
- **Susulan 09-10-2026 10:10** (lewat 9d, tanpa perubahan kode): uji ulang gagal ORA-00904 `NOPOLIS_ARAS`. Work owner MENAMBAH kolom
  itu di DEV dan memperbarui `DDL\CEDING_FACINPRODUCTION.txt`. Skema DEV kini: IDPEGA, NOPOLIS, CEDINGID VARCHAR2(100),
  CEDINGNAME VARCHAR2(100), PCT_SHARE, PRODKE, NOENDORS VARCHAR2(50), CURRENCY_NAME, TSI_RNM, PREMIUM, NOPOLIS_ARAS,
  NOENDORS_ARAS `[terverifikasi]` (DDL baru dibaca sesi c3: CEDINGID / CEDINGNAME VARCHAR2(100), NOPOLIS_ARAS / NOENDORS_ARAS
  VARCHAR2(100)). INSERT 8 kolom cocok; NOENDORS_ARAS tetap tidak ditulis (selalu NULL). CURRENCY_NAME / TSI_RNM / PREMIUM tidak
  ditulis prosedur Pega -> NULL. Lebar CEDINGID / CEDINGNAME DEV 100; DDL lama 500.
  - `[pertanyaan terbuka]` Sumber CEDINGNAME = `T_CEDINGCEDANTLIST.CEDING_CO_NAME` VARCHAR2(**500**) (migrasi 199). Nama ceding
    > 100 bita -> ORA-12899 -> Accept Policy gagal (rollback). Kode CEDINGID aman (sumber `CEDING_CO` VARCHAR2(50),
    `lebarKodeCeding`). Pilihan: potong ke 100, tolak lebih awal dengan pesan, atau perlebar kolom DEV — belum diputuskan.
  - **Keputusan work owner 09-10-2026 (lewat 9d)**: pilihan (c) — PERLEBAR kolom DEV, disamakan dengan DDL produksi. Kode TIDAK
    diubah (tanpa potong, tanpa tolak). **Item DBA (K69a2-9)**, dijalankan DBA / work owner, bukan aplikasi:
    `ALTER TABLE POOLDATA.CEDING_FACINPRODUCTION MODIFY (CEDINGNAME VARCHAR2(500), CEDINGID VARCHAR2(500));`. Sampai itu
    dijalankan di sebuah lingkungan, nama ceding > 100 bita di sana tetap gagal ORA-12899 (rollback).

### BUG DEV 09-10-2026 10:13 — FACINPRODUCTION.POLICYTYPE ORA-12899 (sesi c3)

- Accept TL Policy NB-151029: `ORA-12899: value too large for column "POOLDATA"."FACINPRODUCTION"."POLICYTYPE" (actual: 13, maximum:
  10)`. Sebab: T_QUOTATIONDATA.POLICY_TYPE berisi LABEL "Master Policy" (butir 78.2 lama); Pega menyimpan standard value
  (SaveFacinProdFireNB_Act L1406 `.QuotationData.PolicyType`).
- **Keputusan work owner 09-10-2026**: *"COBA CEK PolicyType, YANG DISIMPAN ITU STANDARD VALUENYA. ( 0,1,2)"* — RALAT butir 78.2
  (KEPUTUSAN-30-09-2026.md). `[terverifikasi]` `DDL\PolicyType.xml` dipasangkan per rowdata: 0 Individual Policy, 1 Master Policy,
  2 Declaration Policy.
- Backend (`services/policytype.go` `kodePolicyType`):
  - GET kasus `general.policyType` = kode (label lama dipetakan);
  - FACINPRODUCTION.POLICYTYPE = kode;
  - PUT general hanya `""` / `0` / `1` / `2` (400 selain itu, termasuk label);
  - CSV objek membawa nilai GET (kode), lolos;
  - JSON_POLIS tidak memuat Policy Type (DATA_JSON minimal).
- **Proteksi** (`proteksi.go proteksiPayment`) `[terverifikasi]` ProtectPremiPolicy_Act langkah 7 / 11 dilewati bila
  `PolicyType==1`: "Premi on paymnet can't be 0" kini TIDAK menahan Master Policy (dulu tanpa gerbang karena nilainya label).
  Langkah 8-9 (Master Policy + MOP, premi harus 0) tidak berlaku: IsMOP tidak ada di sistem baru.
- `premium.MasterPolicy` (Marine Cargo `PolicyType == 1 && IsMOP`) hanya diisi rekonsiliasi (JSON lama, sudah kode) — jalur FIRE
  hidup tidak membacanya; tidak berubah.
- ⚠️ Temuan, TIDAK diubah: Save tab Payment (`payment.go proteksiPayment`, port ProtectPremiPolicy_Act yang sama) juga memunculkan
  "Premi on paymnet can't be 0" tanpa gerbang Master Policy. Menyamakannya perlu keputusan.
- Data lama berlabel: ditangani saat dibaca; migrasi UPDATE TIDAK dibuat. Usulan bila ingin rapi: satu migrasi data (persetujuan
  dulu) `UPDATE T_QUOTATIONDATA SET POLICY_TYPE = DECODE(POLICY_TYPE, 'Individual Policy','0','Master Policy','1','Declaration
  Policy','2', POLICY_TYPE)`.
- Uji: `TestKodePolicyType`, `TestPolicyTypeLabelLama` (GET "1", FACINPRODUCTION "1" untuk label "Master Policy"),
  `TestProteksiPremiNolMasterPolicy` (+ instrumen Individual / kosong tetap ditahan); diralat `TestSimpanGeneral` (services /
  handler: kirim kode; label / "3" -> 400). Mutasi 8 / 8 mati.
- **Keputusan work owner 09-10-2026 atas temuan Save tab Payment**: *"Samakan dengan Pega"*. `payment.go proteksiPayment` kini
  menerima kode Policy Type: langkah 11 "Premi on paymnet can't be 0" dilewati bila `PolicyType == 1`, sama dengan proteksi Accept.
  Kode dibaca lewat `policyTypeKasus` (`policytype.go`; kodePolicyType atas T_QUOTATIONDATA.POLICY_TYPE, label lama ikut
  dipetakan; tanpa pembaca kasus -> "" = tanpa gerbang; case tidak ada -> 404). Migrasi UPDATE label -> kode: **TIDAK** (keputusan
  yang sama). Uji: `TestProteksiPayment` (+ kode 1 / 0 / 2), `TestPolicyTypeKasusPayment`, `TestSimpanPaymentMasterPolicy` (Save
  Payment premi 0: Individual / "0" memuat pesan, Master / "1" tidak). Mutasi 3 / 3 mati.
