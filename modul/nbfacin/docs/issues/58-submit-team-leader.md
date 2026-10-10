# 58: Submit tahap Team Leader

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 07-10-2026 (gambar layar case
> di posisi Team Leader baca-saja) "team leader kan bisa aksep dan submit". Survei korpus oleh subagen (pengurai XML); butir
> kunci diperiksa ulang sesi 0f: `DataTransform\IsJuniorUW_B.xml`, konektor `Decision9`, syarat tampil `.EmailTypeTL` /
> tombol Accept TL, sel Deductible Acceptance di `Section\Periode.xml`.

**Status:** frontend selesai 07-10-2026 (uji hijau); backend diserahkan ke sesi c3.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Flow\InputInwardFacultativeOffer.xml` Assignment14 TEAM LEADER (workbasket `ReasFacInTeamLeader`): flow action
  `InwardFacultative` (SAMA dengan marketing: section `InputInwardFacultative`, post-processing `SetValidateDate_PostAct`).
- `Section\EmailSection.xml`: radio `.EmailTypeTL` (label "Email Type") tampil bila `pyWorkPage.PositionNote ==
  'ReasFacInTeamLeader'`, wajib bila `FlagOnGoingPolicy==0`; radio marketing `.EmailType` disembunyikan untuk TL. Daftar opsi
  properti tidak di korpus; nilai yang dibandingkan di seluruh korpus hanya **1 / 2 / 3 / 7** (Accept / Reject / Ask / Decline
  menurut pemetaan `SaveViewSuggest`); TIDAK ada Banding / Revise untuk TL.
- Tombol TL: Accept (`.EmailTypeTL = 1`, proteksi seperti marketing), Reject (konfirmasi `RejectNotificationFlow`), Ask,
  Decline (Choose Cause of decline + konfirmasi).
- `DecisionTable\IsUWAccepted` → `Decision9`: confirm → **Assignment13 JUNIOR UNDERWITER B** via `DataTransform\IsJuniorUW_B`
  (PositionNote "ReasFacInJuniorUnderwriting", `OfferFacIn.DateToUW` = sekarang, NBStatus "<prefix> IS IN JUNIOR UNDERWRITER
  (B)'S INBOX", `SubmitToUW` bila kosong); reject / ask → `Utility8 SendEmailPolicy` → kembali ke **Assignment12 MARKETING**
  (PositionNote "ReasFacInMarketing", IsCedingConfirm "Offer", FlagOnGoingPolicy 0, NBStatus "<prefix> IS IN <MarketingName>'S
  INBOX"); decline → `End5`.
- Isian khusus TL (`Section\Periode.xml`): centang "Deductible Acceptance Case" (`.IsDeductibleAcceptance`, tampil bila TL dan
  `IsBedaCurr = 1` = ada deductible bermata uang beda dari item) dan "Comment Deductible" (wajib bila dicentang); proteksi
  `Protection_Act` langkah 9–12: "Deductible Acceptance cannot be empty" / "Comment Deductible cannot be empty".
- Riwayat: `SetValidateDate_PostAct` 3.3 (Approval = EmailTypeTL), `InsertHistoryAkseptasiPega` (ACCEPT / REJECT / ASK /
  DECLINE, workbasket ReasFacInTeamLeader). FACINPERFORMANCE tidak ditulis untuk TL. Tab detail dapat diubah TL (tanpa syarat
  baca-saja khusus TL).
- `SendEmailPolicy` (reject / ask) mengirim email hanya bila `IsPEGAPROD`.

## Keputusan work owner 07-10-2026 ("setuju" + gambar layar TL + `DDL\EmailTypeTL.xml`)

`[terverifikasi]` `D:\migrasi\RNM\DDL\EmailTypeTL.xml` (`ASM-FW-GISFW-WORK!EMAILTYPETL`, PromptList): 1 Accept, 2 Reject, 3 Ask,
7 Decline = gambar layar.

1. Empat opsi TL (Accept / Reject / Ask / Decline).
2. Tanpa email saat Reject / Ask.
3. TL boleh mengubah isi tab selama case di posisinya (pemegang workbasket posisi saat itu).
4. Isian Deductible Acceptance Case + Comment Deductible dan proteksinya dibangun.
5. Decline TL = Resolved-Decline (seperti marketing).
6. Tahap UW menyusul (case Accept TL muncul di kotak masuk Junior Underwriter (B), baca-saja).

## Frontend (sesi 0f)

- `components/BlokEmailType.tsx`: opsi menurut `TampilanSubmit.tahap` (`opsiEmailType`), isian TL Deductible Acceptance Case /
  Comment Deductible (tampil bila `deductibleBedaMataUang`; centang / lepas fokus memeriksa ulang proteksi Accept), Ask tanpa
  konfirmasi; `api.ts` (`KodeEmailType` + '3', `TampilanSubmit` + tahap / deductible*, `IsiSubmit` + deductible*,
  `IsiPeriksaSubmit`); `labels.ts` `EMAIL_TYPE_TL`; `nbfacin.css` `.nbf-et__deductible`; uji `BlokEmailType.test.tsx`.
- **K58-1** Isian Deductible ditaruh di blok Email Type (Pega: `Section\Periode.xml` = blok General) - hanya relevan untuk Submit TL.

## Backend (sesi c3, 07-10-2026)

**Status:** backend selesai; migrasi 210 DITULIS, belum dijalankan. Kontrak = `frontend/api.ts`.

- **K58-2 Satu aturan tahap** (`services/bacasaja.go`): `tahapDari` = "marketing" (POSITION Offer / ReasFacInMarketing, flag 0) /
  "teamLeader" (ReasFacInTeamLeader, flag 0) / "" (baca-saja). `alasanUbah`: tahap "" → 409; akun tanpa workbasket POSISI saat
  itu (marketing ReasFacInMarketing, TL ReasFacInTeamLeader) → 403. Dipakai GET kasus `bolehUbah`, `bolehSubmit`, dan
  penjaga 13 endpoint tulis (K57-20) - TL kini dapat mengubah tab dan Submit di posisinya; marketing di posisi TL 403. Tab IT
  tetap syarat tahap saja. Scoring Risk tetap "Write" HANYA tahap marketing (`CallDualScoringRisk` `[terverifikasi]`) - di
  posisi TL PUT scoring lolos penjaga workbasket lalu 409 ErrScoringBacaSaja.
- **K58-3 Email Type per tahap** - marketing 1 / 2 / 4 / 7, TL 1 / 2 / 3 / 7 (`DDL\EmailTypeTL.xml`); lainnya 400. `POST
  …/submit/periksa` kini `{emailType, deductibleAcceptance, commentDeductible}` (ketat); kasus baca-saja → 409.
- **K58-4 Proteksi Accept TL** = proteksi marketing + Protection_Act 9-12 verbatim, di DEPAN (langkah 10 / 12 sebelum 17-18):
  `IsBedaCurr` && tidak dicentang → "Deductible Acceptance cannot be empty"; dicentang && Comment kosong → "Comment Deductible
  cannot be empty". `IsBedaCurr` padanan: Pega menyetelnya di `Activity\SaveDeductible.xml` 2-3 per Save SATU deductible
  (nilai = deductible terakhir disimpan); di sini = ADA deductible tersimpan ber-Currency terisi (bukan kosong / 'UNKNOWN',
  A169) yang beda dari Currency item induknya `[keputusan agent]`. Tidak dicentang = "false" `[dugaan]` (checkbox Pega
  mengirim "false"; prasyarat Pega `IsDeductibleAcceptance=="false"`). Reject / Ask / Decline tanpa proteksi.
- **K58-5 Hasil TL** - Accept → POSITION "ReasFacInJuniorUnderwriting", NB_STATUS "NB IS IN JUNIOR UNDERWRITER (B)'S INBOX",
  DATE_TO_UW = sekarang, SUBMIT_TO_UW = sekarang bila kosong (teks Pega DateTime; `IsJuniorUW_B` `[terverifikasi]`). Reject /
  Ask → POSITION "ReasFacInMarketing", NB_STATUS "NB IS IN " + MARKETINGOFFICER.CLIENTNAME APA ADANYA + "'S INBOX"
  (Transition170 `[terverifikasi]`; FLAG_ONGOING_POLICY sudah "0", tidak ditulis ulang), TANPA email. Decline → POSITION
  "Decline", STATUS_WORK "Resolved-Decline"; NB_STATUS memakai teks Decline marketing "NB WAS DECLINED BY  <UPPER nama>"
  `[keputusan agent]` - SetValidateDate 6.8 hanya untuk PositionNote Admin / Marketing, sehingga Pega menulis "NB IS IN  'S
  INBOX" (NextUser kosong) untuk Decline TL `[terverifikasi]` - dilaporkan. Isian Deductible (IS_DEDUCTIBLE_ACCEPTANCE "true" /
  "false", COMMENT_DEDUCTIBLE) disimpan SETIAP Submit TL (Pega menyimpan saat centang berubah; PUT general TIDAK
  menyimpannya `[keputusan agent]`) - GET submit menampilkan nilai tersimpan.
- **K58-6 Riwayat** - HISTORYAKSEPTASIPRODUCTION (Approval Accept / Reject / Ask / Decline, POSISI "Offer", PUTARAN 1) +
  HISTORYAKSEPTASIPEGA (STATUS ACCEPT / REJECT / ASK / DECLINE, WORKBASKET "ReasFacInTeamLeader"). FACINPERFORMANCE TIDAK untuk
  TL (`InsertFacInPerformance_Act` tidak membaca EmailTypeTL `[terverifikasi]`). ⚠️ `SetValidateDate_PostAct` 3.3 memakai
  `@if(EmailType=="", EmailTypeTL, EmailType)` - bila `.EmailType` marketing masih terisi saat TL Submit, Pega menulis Approval
  marketing, bukan TL `[dugaan]` (tidak terlihat ada yang mengosongkannya) - sistem baru menulis Approval TL.
- **K58-7 Kotak masuk** - Accept TL tampil di antrean ReasFacInJuniorUnderwriting (sudah di `antreanKasusNB`); Reject / Ask
  kembali ke kotak masuk marketing PEMBUATNYA (POSITION ReasFacInMarketing + CREATE_OP).
- **K58-8 Migrasi 210** (DITULIS): `ALTER TABLE T_GENERAL_POLIS ADD (IS_DEDUCTIBLE_ACCEPTANCE VARCHAR2(10), COMMENT_DEDUCTIBLE
  VARCHAR2(500), DATE_TO_UW VARCHAR2(30), SUBMIT_TO_UW VARCHAR2(30))` - dua pertama kolom rancangan T_GENERAL_POLIS; dua
  terakhir rancangan T_WORK_POLIS (premiumlistlife) → amandemen loader T_GENERAL_POLIS (asal = repository, tanpa FIELD ASLI:
  V-48 melarang medan akar ganda). WAJIB sebelum biner (GET / POST submit membaca / menulis kolomnya).
- **Instrumen**: `TestTampilanSubmitTL`, `TestProteksiDeductibleTL` (verbatim, urutan sebelum Scoring, Email Type per tahap),
  `TestAjukanSubmitTL` (empat hasil, riwayat, Deductible, KeUW, 422 / 400 / 403), `TestPenjagaBacaSaja` / `TestAksesUbahKasus`
  (TL boleh, marketing di TL 403, UW / flag 1 baca-saja), `TestBacaSajaHandler`, `TestSQLSubmit`; mutasi 11 titik → semuanya
  tertangkap.

## Pemetaan workbasket (work owner 07-10-2026, ADR-0032)

Keputusan: "jika ada yang menggunakan ReasFacInJuniorUnderwriting ganti menjadi ReasFacInJuniorUnderwritingB, karena kedepannya akan
menggunakan itu". Penyimpangan sadar dari nama Pega.

| Nama Pega (korpus) | Nama sistem baru | Dipakai |
|---|---|---|
| `ReasFacInJuniorUnderwriting` | `ReasFacInJuniorUnderwritingB` | posisi sesudah Accept TL (`IsJuniorUW_B`), Banding JUW_B, antrean kotak masuk |

Layar (frontend) tidak memakai kode ini. Backend + migrasi data (posisi case FAC yang sudah ada) diminta ke c3. Perlu tindakan
data (bukan kode): entri workbasket baru di `M_WORKBASKET` + penugasan akun; `POOLDATA.M_LIMIT_PROPERTYY.WORKBASKET` masih
memakai nama lama (dipakai daftar Banding).

## Pemetaan workbasket (ADR-0032) - keputusan work owner 07-10-2026

Verbatim: "jika ada yang menggunakan ReasFacInJuniorUnderwriting ganti menjadi ReasFacInJuniorUnderwritingB, karena kedepannya
akan menggunakan itu". `[penyimpangan sadar]` dari korpus (ADR-0032 butir 3).

| Nama lama (korpus Pega) | Nama baru (sistem baru) | Dipakai di |
| --- | --- | --- |
| `ReasFacInJuniorUnderwriting` | `ReasFacInJuniorUnderwritingB` | posisi Accept Team Leader (`IsJuniorUW_B`), tujuan Banding `JUW_B` (SetBanding_ACT 3.2), antrean kotak masuk Beranda, predikat bangkitan `ISUWSPREAD` kondisi D (`registry_gen.go`) |

`ReasFacInJuniorUnderwritingA` (JUW_A) TIDAK berubah.

Backend sesi c3:
- `services/submit.go` (`posisiJuniorUWB`, `tujuanBanding["JUW_B"]`), `services/kotakmasuk.go` (`antreanKasusNB`), uji ikut.
- `registry_gen.go` TIDAK disunting tangan: generator `services/rules/bangkit` memegang `PetaWorkbasket` (lama → baru) dan
  menerapkannya pada literal teks kondisi profil NB SAJA (berkutip dan kata telanjang); sidik isi (deteksi identitas
  berkonflik CLAUDE.md §4.6) tetap atas teks korpus; profil folder lain (EDM) verbatim. Dibangkitkan ulang dari
  `D:\migrasi\RNM` - selisih tepat SATU baris (ISUWSPREAD D); `TestRegistrySamaDenganKeluaranBangkit` (dengan
  KORPUS_RNM_AKAR) lulus; `TestPetaWorkbasket` menegaskan pemetaan, A tak tersentuh, tanpa peta = verbatim, profil lain tanpa
  peta. Dipilih pemetaan di generator (bukan saat evaluasi) supaya registry yang dibaca orang = perilaku yang berjalan; sistem
  baru menulis POSITION kode B sehingga predikat verbatim tidak akan pernah benar.
- Migrasi **211** (DITULIS, data saja): `UPDATE T_WORK_POLIS SET POSITION = 'ReasFacInJuniorUnderwritingB' WHERE LINI = 'FAC'
  AND POSITION = 'ReasFacInJuniorUnderwriting'` (+ jalur mundur kebalikannya; baris yang sejak awal ditulis kode B ikut
  dikembalikan - hanya bersama biner lama). Riwayat HISTORYAKSEPTASIPEGA lama TIDAK diubah. Urutan: 211 SEBELUM biner (biner
  baru tidak lagi mengenal kode lama: case lama di posisi itu tak muncul di kotak masuk dan dianggap baca-saja).

⚠️ Dilaporkan, TIDAK diubah:
- (a) `POOLDATA.M_LIMIT_PROPERTYY.WORKBASKET` memuat `ReasFacInJuniorUnderwriting` untuk JABATAN JUW_B (`DDL\M_LIMIT_PROPERTYY.csv`,
  4 baris TEAM_GROUP 1-4). `GetAksepBanding_SQL` mencocokkan WORKBASKET master dengan workbasket non-marketing TERAKHIR di
  HISTORYAKSEPTASIPEGA. Kelak, Submit tahap Junior UW B (belum dibangun) menulis WORKBASKET kode B → opsi Banding case yang
  terakhir dipegang JUW B KOSONG kecuali master diperbarui. Saat ini tidak berdampak (Submit hanya menulis ReasFacInMarketing /
  ReasFacInTeamLeader). Usul: bila tabel POOLDATA dipakai bersama Pega yang masih berjalan, JANGAN ubah master - pemetaan di
  kueri (`WORKBASKET IN (:lama, :baru)`) saat tahap UW dibangun; bila Pega sudah pensiun, DBA memperbarui 4 baris master.
  `[pertanyaan terbuka]` work owner / DBA.
- (b) Data login: `M_WORKBASKET` (nama tampil kotak masuk) dan `M_LOGIN_GO_WORKBASKET` (penugasan akun) wajib memuat
  `ReasFacInJuniorUnderwritingB` - tanpa itu tidak ada akun yang memegang antrean JUW B (kotak masuk kosong, case baca-saja).
  Data, bukan kode; DBA / admin login.
- (c) Modul `endorsmentfacin` (`registry_gen.go` ±1855) masih memakai nama lama - bukan modul nbfacin, dicatat saja.

## Kolom inbox "Time in Underwriting" (work owner 07-10-2026, gambar layar inbox underwriter)

- `[terverifikasi]` (gambar layar) Inbox workbasket `ReasFacInJuniorUnderwritingB`, `ReasFacInJuniorUnderwritingA`,
  `ReasFacInUnderwriting`, `ReasFacInSeniorUnderwriting`, `ReasFacInDepHeadUnderwriting` punya kolom **Time in Underwriting**
  ("1d 1h ago") di antara Validity Date dan Status Inbox. Report inbox itu tidak ada di korpus; label dari gambar layar.
- `[keputusan work owner]` 07-10-2026: isinya dari **`SUBMIT_TO_UW`** (`.OfferFacIn.SubmitToUW`, migrasi 210; Accept Team
  Leader `IsJuniorUW_B` mengisinya hanya bila kosong = waktu pertama kali ke UW). ⚠️ RALAT: rancangan awal memakai
  `DATE_TO_UW` sebagai `[dugaan]` (ditimpa tiap Accept TL) - diganti; kunci kabel ikut berganti `dateToUw` -> `submitToUw`.
- Frontend: `labels.ts` `KOLOM_INBOX_UW_NBFACIN` + `WORKBASKET_UW_NBFACIN`; `beranda.ts` `lamaDiUw` (>= 1 hari: "Xd Yh ago",
  di bawah sehari gaya Time Since Last Update - keputusan agent, gambar hanya memperlihatkan "1d 1h ago") dan `tampilKolomUw`
  (workbasket UW; daftar semua workbasket = tampil bila ada baris ber-`submitToUw`); uji `beranda.test.ts`.
- Backend (sesi c3): field `submitToUw` di `GET /api/nbfacin/kotak-masuk/kasus` - teks Pega DateTime GMT ->
  `YYYY-MM-DD HH:MM:SS` WIB, kosong / tak terbaca = "".

## Kolom inbox UW "Time in Underwriting" - `dateToUw` (permintaan work owner 07-10-2026, kontrak dari sesi frontend)

- GET `/api/nbfacin/kotak-masuk/kasus` menambah kunci `dateToUw` (string, selalu ada) sesudah `tglUpdate`; 16 kunci = api.ts
  `BarisKotakMasukNbFacIn` (uji dua arah + urutan `TestKunciBarisInboxSamaDenganApiTs`, membaca api.ts).
- Sumber `T_GENERAL_POLIS.DATE_TO_UW` (migrasi 210; ditulis Accept Team Leader sebagai BentukWaktuPega UTC). Dibaca mentah di
  SQL (tanpa TO_DATE atas teks), dikonversi di Go (`services.waktuUWKeKabel`) ke `YYYY-MM-DD HH:MM:SS` WIB (bentuk tglUpdate).
  NULL / kosong / teks yang tidak cocok BentukWaktuPega -> "" (tidak mengirim teks mentah). Tanpa migrasi baru; 210 tidak diubah.

> **RALAT 07-10-2026 (keputusan work owner, verbatim: "isinya dari SUBMIT_TO_UW")** - dua butir di atas tidak berlaku lagi
> pada sumber dan nama kunci. Sumber kolom inbox UW "Time in Underwriting" diganti dari `T_GENERAL_POLIS.DATE_TO_UW` ke
> `T_GENERAL_POLIS.SUBMIT_TO_UW` (juga migrasi 210), dan kunci JSON `dateToUw` diganti `submitToUw` (posisi sama: terakhir,
> sesudah `tglUpdate`; api.ts sudah diralat sesi frontend). Konversi tetap: BentukWaktuPega UTC -> `YYYY-MM-DD HH:MM:SS` WIB di
> Go, NULL / kosong / tak cocok -> "". Tanpa migrasi baru; 210 tidak diubah.
> Beda arti yang perlu diketahui `[terverifikasi]` (`repository.sqlKeUW`, `DataTransform\IsJuniorUW_B.xml`): DATE_TO_UW ditimpa
> SETIAP Accept Team Leader, SUBMIT_TO_UW hanya diisi bila kosong (`NVL(SUBMIT_TO_UW, :2)`) = waktu PERTAMA case ke UW. Case yang
> kembali ke marketing lalu di-Accept TL lagi menghitung lama di UW sejak Accept pertama.

## Kotak masuk per team group (keputusan work owner 08-10-2026)

Kutipan work owner: "yg muncul sesuai workbasketnya dan team groupnya". Keputusan diambil lewat AskUserQuestion dan diteruskan
sesi 9d. Ini **PENYIMPANGAN SADAR** dari Pega.

Bukti Pega `[terverifikasi]`:
- Di Pega kotak masuk hanya per workbasket. Ke-21 assignment `Flow\InputInwardFacultativeOffer.xml` memakai router
  `ToWorkBasket` ke workbasket tetap, tanpa workbasket per team.
- Tidak ada aturan inbox NB FacIn yang menyaring TeamGroup.
- Satu-satunya saringan TeamGroup ada di daftar opportunity SFA (`ReportDefinition\GetListOpportunityF.xml`).

Aturan (`services/kotakmasuk.go` `saringanKotakMasuk`, `repository/kotakmasuk.go` `saringKotakMasuk`): satu saringan untuk
hitungan dan daftar.

| Workbasket | Saringan | Team akun |
| --- | --- | --- |
| Marketing (POSITION Offer / ReasFacInMarketing) | team case IN team akun — **menggantikan CREATE_OP** di Beranda | `POOLDATA.MARKETINGOFFICER.TEAMGROUP` ber-`AKSES_LOGIN` = login |
| Team Leader (ReasFacInTeamLeader) | team case IN team akun (bisa lebih dari satu) | `POOLDATA.M_LIMIT_PROPERTYY.TEAM_GROUP` ber-`LOGIN` = login |
| JUW B / A, UW, Senior, DepHead, Financial, lainnya | tanpa saringan team (seperti sebelumnya) | tidak dibaca |

- **Team case** `[keputusan agent]`: dibaca saat itu juga lewat `T_QUOTATIONDATA.MOID` -> `MARKETINGOFFICER.TEAMGROUP`
  (sama dengan tiket 57 / 61), tanpa migrasi. Pega menyimpan `QuotationData.TeamGroup` (GetTeamGroup_Act); bedanya, case
  ikut berpindah bila team MO berubah. Rancangan `T_QUOTATIONDATA.TEAM_GROUP` belum dimigrasi.
- **Akun tanpa team** (tidak ada baris di sumbernya) -> antrean itu kosong + log (butir 3d).
- **Pencocokan:** login dibandingkan `UPPER(TRIM)`, team di-TRIM.
- ⚠️ Bahwa `AKSES_LOGIN` / `LOGIN` berisi nilai yang sama dengan `M_LOGIN_GO.LOGIN_ID` `belum terverifikasi` — perlu
  pemeriksaan DBA.
- **Portal** opportunity (`GET /api/nbfacin/opportunity`) TETAP CREATE_OP.
- **Tanpa migrasi.** Tabel POOLDATA hanya dibaca.
- **Uji:**
  - `TestKotakMasukPerTeam`: matriks marketing / TL / UW × punya team / tanpa team × team case sama / beda; hitungan =
    daftar; UW tidak membaca team akun.
  - `TestSQLKotakMasuk`: SQL, bind berurutan, `kosong()`.
  - Mutasi 8, semua tertangkap.
