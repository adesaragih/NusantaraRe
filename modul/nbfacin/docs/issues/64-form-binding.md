# 64: Form tahap MARKETING (BINDING) — Accept / Banding / Decline / Revise + tombol Binding

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 08-10-2026 (gambar layar Pega
> tahap Binding vs sistem baru) "INI SAAT DI MARKETING(BINDING) KENAPA TIDAK ADA YA". Riset korpus subagen: scratchpad sesi 0f
> `riset-E-binding.md` (aksi sel dua wadah 14 = 14; konektor flow 176 = 176, RISlip 126 = 126).

**Status:** frontend selesai 08-10-2026 (uji hijau); backend selesai sesi c3 08-10-2026 (lihat "Backend sesi c3").

## Keputusan work owner 08-10-2026 (AskUserQuestion)

- **K64-1** Accept di Binding → case pindah ke **Team Leader tahap Policy** seperti Pega (SubProcess1 POLICY →
  `InputInwardFacultativeRISlip`: flag 1, IsCedingConfirm "Policy", lalu Assignment6 TEAM LEADER); layar Policy baca-saja sampai
  tiket jalur Policy.
- **K64-2** Keanehan Pega diperbaiki (penyimpangan sadar): status riwayat sesuai pilihan (Pega selalu "Accepted"), NB Status sesuai
  tujuan (Pega: Decline tidak terpasang), FACINPERFORMANCE BINDING/ACCEPT satu kali (Pega: dua — tombol Binding + Submit).
- **K64-3** Revise ikut Pega: ConfirmBinding = 0 dan Received RI Slip = false; TGL_BIND JSON_OFFER dibiarkan.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\` (+ `D:\migrasi\RNM\DDL\`)

- `DDL\EmailTypeBinding.xml` (`ASM-FW-GISFW-WORK!EMAILTYPEBINDING`): 1 Accept, 4 Banding, 7 Decline, 9 Revise.
- `Section\EmailSectionCeding.xml` (tampil `IsCedingConfirm='Binding' || FlagOnGoingPolicy=='2'`): `.EmailTypeBinding`, tombol
  "Binding" (Accept; disabled bila ConfirmBinding = 1 / proteksi; aksi `UpdateTglBindJsonOffer_Act`: TGL_BIND, ConfirmBinding = 1,
  riwayat "Offer has been bind.", FACINPERFORMANCE BindStatus "Binding"), "Received RI Slip" (`.OfferFacIn.ReceivedRiSlip`, disabled
  bila ConfirmBinding != 1), Comment `.OfferFacIn.CommentCeding`, "Banding to" (`GetAksepBanding_SQL`), "Choose Cause of decline".
- Konfirmasi Revise `Section\ReviseNotification.xml` "Are you sure you want to revise?".
- Decision13: confirm → SubProcess1 POLICY; banding → `SetBanding_ACT` (kembali ke jenjang UW tahap Offer, sama dengan Banding
  marketing); revise → Assignment12 (Offer, flag 0); decline → End3. Post `SetValidateDate_PostAct` (riwayat PUTARAN "2",
  HISTORYAKSEPTASIPEGA, FACINPERFORMANCE Accept: BINDING/ACCEPT + POLICY/OUTSTANDING; Decline: BINDING/DECLINE). Tanpa CALL.
- `SendEmailtoCeding_ACT`: langkah kirim nonaktif (tidak ada email) — ⚠️ memuat kredensial SMTP ter-hardcode (temuan keamanan,
  tidak disalin).

## Frontend (sesi 0f)

- `labels.ts` `EMAIL_TYPE_BINDING`, `TEKS_EMAIL_TYPE.wajibBinding` / `membinding`; `api.ts` `tahap: 'binding'`,
  `TampilanSubmit.confirmBinding` / `receivedRiSlip`, `IsiSubmit.receivedRiSlip`, `KodeEmailType` + '9', `tekanBinding`
  (`POST …/submit/binding`); `components/BlokEmailType.tsx` (opsi Binding, tombol Binding + Received RI Slip, Accept wajib sudah
  Binding, konfirmasi Revise); `nbfacin.css` `.nbf-et__binding`. Uji: `BlokEmailType.test.tsx` blok "tahap Binding (tiket 64)".

## Backend (diminta ke sesi c3, 08-10-2026)

Lihat pesan kontrak tiket 64.

## Keputusan backend 08-10-2026 (dikumpulkan sesi 9d; work owner lewat AskUserQuestion kecuali ditandai)

- **D1** `[work owner]` "2 kolom": migrasi 213 `T_GENERAL_POLIS` ADD CONFIRM_BINDING VARCHAR2(50) + RECEIVED_RI_SLIP VARCHAR2(50) (nama /
  tipe = rancangan `loader/skema_gen.go`), teks Pega "1"/"0" dan "true"/"false"; COMMENT_CEDING tidak ditambah (komentar ke riwayat).
- **D2** `[work owner]` "Ke UW Financial": IsTBonding (`When\IsTBonding.xml` "A OR B OR C OR D": TeamGroup = 5 ATAU MarketingName = 3
  literal nama) → POSITION ReasFacInUnderwritingFinancial, flag 1; nama dari env `NBFACIN_BONDING_MARKETINGNAME` (literal tidak di kode).
- **D3** `[keputusan agent]` CekLimitSpreading_Act (hanya men-SET `IsFlagUW`, dibaca Decision21 sesudah TL Policy confirm) ditunda ke tiket
  Policy.
- **D4** `[work owner]` "Saat Submit Accept": tombol Binding TIDAK menulis FACINPERFORMANCE; BINDING/ACCEPT ditulis SEKALI saat Submit
  Accept dengan RISLIPRECEIVED = Received RI Slip akhir. POLICY/OUTSTANDING tidak ditulis (ikut Pega: USERNAME =
  `LimitAkseptasi.pxResults(1).CARI5`, halaman itu tidak diisi di jalur Binding → `[dugaan]` InsertFacInPerformance_Act langkah 12
  melewati insert).
- **D5** port Protection_Act tiket 57 dipakai ulang (Accept dan Banding; Revise / Decline lewat langkah 8) + langkah 33 CedingCoName kosong
  "Ceding Co cannot be empty" (verbatim) + langkah 34 SOB kosong `[work owner]` "Source Of Business cannot be empty" (penyimpangan sadar:
  Pega memakai pesan Ceding Co); CheckRISLIP (langkah 25, hanya PolicyStatus) dan survey report tidak diport.
- **D6** `[keputusan agent]` NB status Accept "NB IS IN TEAM LEADER'S INBOX" / Bonding "NB IS IN UNDERWRITER FINANCIAL'S INBOX"; riwayat
  APPROVAL = pilihan, POSISI = IsCedingConfirm tahap tujuan (Accept "Policy", Revise / Banding "Offer", Decline "Binding"); JSON_OFFER
  tombol WHERE IDPEGA, baris tidak ada → 409; Revise: CONFIRM_BINDING "0", RECEIVED_RI_SLIP "false", flag 0, POSITION ReasFacInMarketing;
  email ke ceding tidak.

## Backend sesi c3 (08-10-2026)

- **K64-B1 Tahap** `binding` (`services/submitbinding.go`): POSITION ReasFacInMarketing + FLAG_ONGOING_POLICY 2; Submit wajib pemegang
  workbasket ReasFacInMarketing (403); Email Type 1 / 4 / 7 / 9 (lainnya 400). Tab case TETAP baca-saja (`tahapDari` tidak berubah).
  **RALAT 64c (08-10-2026):** tab case tahap Binding DAPAT DIUBAH sampai tombol Binding ditekan - lihat "Tahap Binding dapat diubah".
  `TampilanSubmit` + `confirmBinding` / `receivedRiSlip` (false di luar Binding); `IsiSubmit.receivedRiSlip`.
- **K64-B2 Tombol Binding** `POST …/submit/binding` (badan {}): 401 → 503 → 404 → 403 → 409 bukan tahap Binding → 409 sudah Binding →
  422 proteksi Accept → satu transaksi: penjaga posisi + flag, CONFIRM_BINDING = '1' hanya bila belum, JSON_OFFER TGL_BIND = SYSDATE +
  JSONDATA minimal WHERE IDPEGA (0 baris → 409), riwayat "Offer has been bind." (APPROVAL Accept, POSISI Binding, PUTARAN 2) →
  TampilanSubmit terbaru.
- **K64-B3 Submit** per pilihan: Accept (wajib ConfirmBinding, selain itu 422 "Tekan tombol Binding terlebih dahulu sebelum Accept." —
  teks sistem baru) → TL / UW Financial (bonding) flag 1, CONFIRM_BINDING "1", RECEIVED_RI_SLIP dari isian, FACINPERFORMANCE BINDING/ACCEPT
  (RISLIPRECEIVED); Banding → SetBanding_ACT (tujuanBanding tiket 57, flag 0); Revise → marketing Offer flag 0 + reset; Decline → End3
  (POSITION Decline, Resolved-Decline, "NB WAS DECLINED BY  <NAMA>" dua spasi, BINDING/DECLINE). HISTORYAKSEPTASIPEGA Workbasket
  ReasFacInMarketing, STATUS ACCEPT / BANDING / REVISE / DECLINE.
- **K64-B4 Repository**: `sqlKonteksSubmit` + CONFIRM_BINDING / RECEIVED_RI_SLIP; `TulisSubmit` + FLAG_ONGOING_POLICY generik dan
  CONFIRM_BINDING / RECEIVED_RI_SLIP; FACINPERFORMANCE RISLIPRECEIVED kini diikat (sebelumnya selalu NULL; "" = NULL);
  `TulisTekanBinding`.
- **K64-B5 Config**: `NBFACIN_BONDING_MARKETINGNAME` (inti/backend/config `GrupNBFacIn.BondingMarketingName`, `.env.example`).

**Uji** — `TestSubmitBindingMatriks` (Accept TL / bonding TeamGroup 5 / bonding nama, Banding, Revise, Decline), `TestSubmitBindingDitolak`
(belum Binding, Ceding / SOB kosong, Email Type 2 / 3, workbasket, Banding diproteksi / Revise tidak, periksa), `TestTampilanDanTabBinding`,
`TestTekanBinding`, `TestTekanBindingHandler` (401 / 400 / 409, kunci IsiSubmit api.ts dua arah), `TestSubmitHandler` (TampilanSubmit 16
kunci), `TestSQLFormBinding`, replay migrasi loader 344 kolom, `TestGrupNBFacIn`. **Mutasi 24 / 24 mati.** gofmt / vet bersih;
`go test ./modul/nbfacin/... ./inti/backend/config/ ./cmd/...` hijau (handlers kembali hijau); penjaga hanya kegagalan lama premiumlistlife;
index git 0.

## Deploy

- **Migrasi 213 WAJIB sebelum biner** (GET / POST …/submit membaca CONFIRM_BINDING / RECEIVED_RI_SLIP; tanpa 213 ORA-00904). Ditulis,
  tidak dijalankan agent.
- Env opsional `NBFACIN_BONDING_MARKETINGNAME`. Tombol Binding MENULIS POOLDATA.JSON_OFFER (TGL_BIND, JSONDATA) dan
  HISTORYAKSEPTASIPRODUCTION.

## Terbuka

- CekLimitSpreading_Act / IsFlagUW dan layar Policy (TL / UW Financial flag 1) = tiket jalur Policy.
- `[pertanyaan terbuka]` Reject (2) di Binding: ada di flow / section tetapi tidak di prompt list EmailTypeBinding (tidak diport).
- ⚠️ Temuan keamanan (pemilik Pega): `SendEmailtoCeding_ACT` langkah 2.9 (nonaktif) memuat kredensial SMTP — tidak disalin.
- Inkonsistensi yang dicatat: Protection_Act langkah 8 tidak melewati Email Type 4 (Banding) — di Binding Banding diproteksi; port tiket 57
  untuk Banding tahap marketing tidak memproteksi (perlu ditinjau terpisah).
  **RALAT 57b (08-10-2026):** diselesaikan - Banding tahap marketing kini diproteksi; penahan Banding Binding = flag
  EmailSectionCeding (ProtectSpreading.CARI1 / CARI23 saja, bukan seluruh pesan Accept). Lihat issues/57 "RALAT 57b".

## 64b - penahan tombol Binding dan Submit Accept menurut flag Pega (08-10-2026)

> Teks lama yang diralat (K64-B2 / K64-B3): tombol Binding dan Submit Accept tahap Binding ditahan *"422 proteksi Accept"* =
> SELURUH pesan port Protection_Act tiket 57 + langkah 33 / 34.

Pemicu (work owner 08-10-2026, gambar layar tombol Binding nonaktif): *"tombol binding nya kenapa blm bisa ditekan"*. Bukti sesi 0f
`[terverifikasi]` `Section\EmailSectionCeding.xml` pyDisabledWhen:
- tombol **Binding**: `TempEmail.CARI28 != .PositionNote || .OfferFacIn.ConfirmBinding = 1 || ProtectSpreading.CARI23 = 1 ||
  ProtectDeductible.CARI2=='1' || ProtectDeductible.CARI4=='1'||ProtectSpreading.CARI1 == 1` -> `flagTombolBinding`;
- tombol **Send Email Accept**: `… || .OfferFacIn.ConfirmBinding != 1 || ProtectSpreading.CARI23 = 1 ||ProtectSpreading.CARI1 ==
  1||ProtectZipCode.CARI1==1` -> `flagAcceptBinding` (+ pesanBelumBinding untuk ConfirmBinding != 1);
- Banding: `flagBandingBinding` (CARI1 / CARI23, tiket 57b) tetap.

**Kontrak** - `HasilSubmit.pesanBinding: string[]` (SELALU ada; `[]` di luar periksa Accept tahap Binding dan pada 422): periksa
Accept Binding -> `pesan` = pesanBelumBinding (bila belum) + penahan Accept; `pesanBinding` = penahan tombol Binding. Submit
Accept Binding 422 hanya oleh penahan Accept (+ belum Binding); `POST …/submit/binding` 422 hanya oleh penahan tombol Binding.
Accept tahap marketing / TL / UW TIDAK berubah. ⚠️ Komentar api.ts menyebut pesanBinding "bagian `pesan`" - keliru: pesan deductible
hanya ada di pesanBinding (bukan penahan Accept); kunci kabel sama (uji dua arah `TestHasilSubmitKunciTS`).

**Pemetaan flag -> pesan port** `[terverifikasi]` `NB FacIn\Activity\ProtectFIREMBUPA_Act.xml` / `Protection_Act.xml`:

| Flag | Pesan | Langkah | Status port |
| --- | --- | --- | --- |
| ProtectZipCode.CARI1 | "ZipCode / CITYNAME / DISTRICT / PROVINCE … doesn’t exist in Master!" | FIREMBUPA 9.2.18-21 | sudah (tiket 57) |
| ProtectDeductible.CARI4 | " Deductible on location n object item m Coverage c should be added!" | FIREMBUPA 9.2.30.6.6-7 (`@SizeOfPropertyList(.DeductibleList)=0`) | **BARU 64b** |
| ProtectDeductible.CARI2 | " Deductible on location n object item m Coverage c  DeductibleList d can't be empty!" (dua spasi) | FIREMBUPA 9.2.30.6.13.2-3 (TypeDeductible, Condition, MinMax, Amount, TimeExcess semua kosong) | **BARU 64b** |
| ProtectDeductible.CARI2 | "Comment Deductible cannot be empty" | Protection_Act 11 | tidak berlaku (PositionNote TeamLeader) |
| ProtectSpreading.CARI1 | "Ceding Co cannot be empty" / "Source Of Business cannot be empty" | Protection_Act 31 (`(CedingCoName=="" \|\| SourceOfBusiness=="") && FlagOnGoingPolicy==2`, dilewati EmailTypeBinding 4) + 33 / 34 | sudah, kini ber-flag CARI1 |

Kedua langkah deductible berjalan hanya bila `pyWorkPage.FlagOnGoingPolicy==2` dan `EmailTypeBinding != "4"` (proteksiBerflag
`deductibleBinding`: tahap Binding Accept / tombol Binding; Banding dan tahap lain tidak). `[dugaan]` Pega tidak mereset CARI4 di
Protection_Act langkah 4 (yang direset hanya Deductible CARI1 / CARI2 / CARI11 / CARI21) - CARI4 yang pernah 1 tetap menonaktifkan
tombol sampai halaman dimuat ulang; port menghitung ulang tiap permintaan (tidak ditiru). Konsekuensi: Scoring kosong (tanpa flag)
tidak lagi menahan Accept Binding.

**Uji** - `TestProteksiDeductibleBinding`, `TestPenahanTahapBinding` (matriks 7 pesan + Premi 0 × periksa / tombol / Accept / Banding,
Accept marketing tanpa deductible), `TestHasilSubmitKunciTS`; `TestSubmitBindingDitolak` diralat (Scoring kosong -> Accept tidak
ditahan); fixture `lokasiBinding` (deductible terisi). **Mutasi 64b + 64c 28 / 28 mati**; mutasi 57b 45 / 46 (1 setara, sama).

## 64c - Tahap Binding dapat diubah sampai tombol Binding (08-10-2026)

Kutipan work owner 08-10-2026 (gambar layar NB-151029 tahap Binding, banner baca-saja): *"inikan lagi dimarketing harusnya msh
bisa di edit, kecuali sudah di klik binding, baru smua tidak bisa diedit lagi"*.

- **K64c-1** `services/bacasaja.go` `tahapDari`: tahap `binding` (POSITION ReasFacInMarketing, flag 2, CONFIRM_BINDING != "1") = dapat
  diubah pemegang workbasket ReasFacInMarketing - GET bolehUbah true, penjaga tulis (general, objek, CSV objek, klausa, spreading,
  cedant, deduction, payment, akumulasi, loss record / korespondensi = tab IT) lolos. Sesudah CONFIRM_BINDING "1": 409, alasan
  "Case baca-saja: tombol Binding sudah ditekan." `[keputusan agent]`. Submit tahap Binding tetap boleh sebelum DAN sesudah tombol.
- **K64c-2** `repository/kasus.go` `sqlTahapKasus` + `LEFT JOIN T_GENERAL_POLIS` CONFIRM_BINDING (migrasi 213 wajib, sudah).
- **Scoring** tetap baca-saja di flag 2: `Section\CallDualScoringRisk.xml` "Write" hanya `PositionNote = 'ReasFacInMarketing' &&
  FlagOnGoingPolicy = 0` (varian read `… == 'ReasFacInMarketing' && FlagOnGoingPolicy == 2`) - PUT scoring 409 ErrScoringBacaSaja.
- **Uji** `TestPenjagaBacaSajaBinding` (14 metode tulis × sebelum / sesudah tombol × workbasket; CSV objek ditambahkan ke
  `metodeTulis`), `TestTampilanDanTabBinding` diralat, `TestSQLTahapKasus`, invarian `TestAksesUbahKasus` (Submit Binding).

### Butir 2 - kontrol yang terkunci SEJAK tahap Binding `[terverifikasi]` (dua jalan sepakat: pengurai ElementTree vs cacah baris)

Pindai section NB FacIn ber-`FlagOnGoingPolicy` 2 pada kondisi: 5 berkas (Periode, InputInwardFacultativeDtl, CallDualScoringRisk,
ComfirmPolis_NonLife, EmailSectionCeding).
- `Section\Periode.xml` - 18 kondisi: 14 `.ConfirmBinding = 1` (terkunci SESUDAH Binding - sudah tercakup bolehUbah), **3
  `.ConfirmBinding = 1||pyWorkPage.FlagOnGoingPolicy=2`** = `.PolicyData.StartDateTime` (pyDisabledWhen + pyReadOnlyCondition) dan
  `.PolicyData.EndDateTime` (pyReadOnlyCondition) -> Begin / End Date General TERKUNCI sejak tahap Binding; 1 wadah ber-`.ConfirmBinding
  != 1`.
- `Section\InputInwardFacultativeDtl.xml` - `.OfferFacIn.ShareCedantType` **pyRequiredWhen** `.FlagOnGoingPolicy ==2 || IsGroup` (WAJIB
  diisi di tahap Binding - bukan kunci).
- `CallDualScoringRisk` (scoring baca-saja, di atas); `ComfirmPolis_NonLife` (pergantian sub-section konfirmasi, bukan data).
- `Section\InputInwardFacultative.xml` - wadah detail `ConfirmBinding != 1` / `= 1` (varian terkunci sesudah Binding).

`[usulan, menunggu keputusan]` kontrak: GET kasus `kunciPeriode: boolean` (tahap Binding) -> frontend mengunci Begin / End Date;
PUT …/general di tahap Binding menolak (409) bila StartDateTime / EndDateTime BERBEDA dari tersimpan. ShareCedantType: PUT …/cedant
di tahap Binding tanpa ShareCedantType -> 400 (teks pesan Pega platform "Value cannot be blank" `[dugaan]` - perlu keputusan).

### Butir 3 - efek hilir JSON_OFFER / FACINOFFER `[terverifikasi]`

- Save di tahap Binding (`IsThereAnyObjectLocation_Act` -> `SaveFacIn_Act`) TIDAK menulis JSON_OFFER / FACINOFFER: penulisnya
  (`SaveOfferProduction_Act` <- `SaveJsonOfferFacIn_Act` / `UpdateJSONOffer_ACT`) hanya dipanggil flow Decision23 Else
  (`SetToJsonOffer_ACT`), `SetValidateDateUW_PostAct`, `InputEndorsement`, dan tombol "Save to Json Offer"
  `InputInwardFacultative_2` yang TERSEMBUNYI (`pyContainerVisibleWhen 1==2`); `SaveFacIn_Act` nol rujukan.
- Tombol **Binding** (`UpdateTglBindJsonOffer_Act` langkah 1-2): `UpdateTgl_BindJsonOffer_SQL` `UPDATE JSON_OFFER a SET a.TGL_BIND =
  sysdate, jsondata = {InputaData.CARI2} WHERE a.JSONDATA.IDNewBisnis = …` dengan CARI2 = JSON halaman OfferFacIn UTUH -> perubahan
  di tahap Binding ikut terekam di JSON_OFFER Pega saat Binding; FACINOFFER TIDAK ditulis ulang (tetap isi saat Accept UW).
  Port tiket 64: JSONDATA minimal `{"IDNewBisnis":id}` (ADR-0028 - pembangun payload JSON tidak dimigrasi) - perubahan tahap
  Binding TIDAK masuk JSON_OFFER sistem baru `[pertanyaan terbuka]` work owner: cukup, atau JSONDATA perlu isi (melanggar
  ADR-0028)?
- Langkah 5 `InsertFacInPerformance_Act` (BindStatus "Binding") - tidak diport (keputusan D4).


## 64d - kunci periode, Share Cedant Type wajib, JSON_OFFER (keputusan work owner 08-10-2026 lewat sesi 9d)

- **K1 "Ikut Pega"** - Begin / End date terkunci di tahap Binding (`Section\Periode.xml` `.PolicyData.StartDateTime` /
  `.EndDateTime` `.ConfirmBinding = 1||pyWorkPage.FlagOnGoingPolicy=2`): GET kasus `kunciPeriode: boolean` = tahap Binding
  (POSITION ReasFacInMarketing + flag 2, sebelum maupun sesudah tombol; = `TampilanSubmit.tahap === "binding"`), `KunciPeriode` LUNAK
  (galat baca tahap -> false). PUT …/general di tahap itu: beginDate / endDate isian dibandingkan dengan tersimpan DALAM BENTUK
  KABEL GET (`waktuKeKabel` - jam tersimpan Pega tidak memicu tolakan palsu); berbeda -> 409 `ErrPeriodeTerkunci` "Begin Date dan End
  Date tidak dapat diubah di tahap Binding" (membungkus ErrKasusBacaSaja); medan General lain tetap tersimpan.
- **K2** - TampilanCedant `wajibType: boolean` = `wajib` (ProtectShareCedant: Group / FlagOnGoingPolicy 2 data cedant) ATAU tahap
  Binding; menggerbang Share Cedant Type SAJA (proteksi C-3 lain tetap `wajib`). PUT …/cedant `wajibType` + tipe kosong -> 400
  "Share Cedant Type wajib diisi." (pesan pilihan work owner; Pega memakai pesan bawaan platform - pyRequiredWhen
  `InputInwardFacultativeDtl` `.FlagOnGoingPolicy ==2 || IsGroup`). RALAT teks lama "Share Cedant Type wajib diisi" (tanpa titik,
  K49-2) -> dengan titik.
- **K3** - JSON_OFFER.JSONDATA tombol Binding tetap minimal (ADR-0028); tidak ada perubahan.

**Uji** `TestKunciPeriode`, `TestSimpanGeneralPeriodeBinding` (sama / Begin / End / kosong / marketing bebas), `TestWajibTypeCedant`,
handler `TestKunciPeriodeHandler`, `TestWajibTypeCedantHandler`; kunci KasusNB api.ts dua arah 15 kunci (`TestKunciKasusSamaDenganApiTs`);
uji JSON persis cedant / kasus diperbarui. **Mutasi 11 / 12 mati**; 1 setara (`err == nil` dibuang di KunciPeriode: galat baca tahap
memberi TahapKasus kosong -> tahapSubmitBinding false; mutan tidak terkompilasi).

**Frontend** (sesi 0f): `TanggalDMY disabled={kunciPeriode}` untuk Begin / End date (nilai tersimpan dikirim ulang apa adanya); TabCedant `typeWajib = v.wajib || v.wajibType === true` menggerbang tanda * dan pesan; label `typeWajib` -> "Share Cedant Type wajib diisi.". tsc bersih, vitest 969 / 969.
