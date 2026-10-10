# 61: Submit tahap Underwriting (Email Type UW)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 08-10-2026 (gambar layar UW Pega:
> Email type Accept / Reject / Ask / Decline, Comment, Validity Of Acceptance … days, or up to … Waiting Bind) "pilihan aksep
> reject nya di uw mana. knpa tidak ada". Survei korpus oleh subagen (pengurai XML expat, atribut dipasangkan per `rowdata`;
> konektor flow diuji dua wadah `pyConnectors` vs `pyFromTasks/pyToTasks`: 176 = 176).

**Status:** frontend selesai 08-10-2026 (uji hijau); backend diserahkan ke sesi c3.

## Keputusan work owner 08-10-2026 (AskUserQuestion)

- **K61-1** Cakupan "Jenjang UW dulu": Reject / Ask / Decline penuh; Accept naik jenjang UW lewat limit
  (`GetLimitAkseptasi_JUW_UW`). Accept yang keluar dari jenjang UW (Kadiv / Direktur / Manager Teknik / FacRetro /
  SaveToProduction / Binding / Policy) ditolak "belum tersedia" sampai tiket berikutnya.
- **K61-2** Email Type dari `D:\migrasi\RNM\DDL\EmailTypeUW.xml` (`ASM-FW-GISFW-WORK!EMAILTYPEUW`): 1 Accept, 2 Reject,
  3 Ask, 7 Decline `[terverifikasi]` (uji `BlokEmailType.test.tsx`).
- **K61-3** Decline juga di JUW A — penyimpangan sadar (flow Pega tanpa konektor Decline dari Assignment21).
- **K61-4** NBStatus Senior / DepHead (korpus berisi nama orang, tidak disalin): "NB IS IN SENIOR UNDERWRITER'S INBOX" /
  "NB IS IN DEPHEAD UNDERWRITER'S INBOX" — work owner "setuju".
- **K61-5** Jenjang dari JUW_B: work owner "emg JUW_A lagi tidak ada, jadi dari JUW_B lansung ke underwriter atau ga senior
  underwriter, liat di GetLimitAkseptasi_JUW_UW" — sejalan `[terverifikasi]` dengan `RDBList\GetLimitAkseptasiJUWA_SQL.xml`
  (JABATAN berikutnya = LIMIT_BOTTOM terkecil di atas jabatan sendiri dan di bawah TSI Nusantara Re; tanpa baris = jalur akhir).
- **K61-6** Migrasi 212 (hanya ditulis): DAYS_VALIDITY, DATE_VALIDITY, WAITING_BIND_DATE + COMMENT_TEXT di T_GENERAL_POLIS.
- **K61-7** Proteksi Accept UW = `CheckSpreadingProtect_ACT` bagian spreading (baris kosong, total % ≠ 100, premium spreaded ≠
  premium RNM, TSI spreaded > kapasitas, treaty tidak valid); "Spreading ORS" / "special flow" / "Production Date" dilewati.
- **K61-8** FACINPERFORMANCE ikut Pega (Accept: ACCEPT/OFFER + OUTSTANDING/BINDING; Decline: DECLINE/OFFER), tanpa penggantian
  nama ter-hardcode.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `FlowAction\InwardFacultative_IsUW.xml` → `Section\InputInwardFacultative_2` → `Section\EmailSection_IsUW.xml`:
  `.EmailTypeUW` (required/visible `StatusBusiness!=3 && FlagOnGoingPolicy != 1`), Comment `.OfferFacIn.Comment` (required),
  "Validity Of Acceptance" `.OfferFacIn.DaysValidity`, "days, or up to" `.OfferFacIn.DateValidity` (disabled), "Waiting Bind"
  `.OfferFacIn.WaitingBindDate` (disabled). Tanpa "Cause of decline" (yang hanya ada di `EmailSection.xml`).
- `Activity\GetDateValidity_ACT.xml`: NB tanpa edit = hari ke StartDateTime + 1 (maks 30); edit = hari ini + DaysValidity;
  WaitingBind L701 (arti argumen addCalendar `[belum terverifikasi]`).
- Routing `Flow\InputInwardFacultativeOffer.xml`: Reject / Ask → marketing (Transition170); Decline → End; Accept JUW B / JUW A
  / UW → `GetLimitAkseptasi_JUW_UW` (`RDBList\GetLimitAkseptasiJUWA_SQL.xml`, `M_LIMIT_PROPERTYY`, satu jenjang per submit) →
  Decision42 (SENIORUW / DEPHEADUNDERWRITER / JUW_A / UNDERWRITER / Else); Senior / DepHead / Financial → Decision24 (jalur akhir).
- Riwayat: `SetValidateDateUW_PostAct` step 5 (ViewSuggest) dan step 8 (`HISTORYAKSEPTASIPEGA`, Status ACCEPT / REJECT / ASK /
  DECLINE); Decline NBStatus "NB WAS DECLINED BY  " + nama pelaku (dua spasi).
- Proteksi `Activity\CheckSpreadingProtect_ACT.xml` (pesan verbatim di kontrak backend).

## Frontend (sesi 0f)

- `labels.ts`: `EMAIL_TYPE_UW` (opsi + label isian, verbatim korpus), `TEKS_EMAIL_TYPE.hariAngka` (teks sistem baru).
- `api.ts`: `TampilanSubmit.tahap` + `'uw'`, `daysValidity` / `dateValidity` / `waitingBindDate`; `IsiSubmit.daysValidity`,
  `IsiPeriksaSubmit.daysValidity`; `hitungValiditas` (`POST …/submit/validitas`).
- `components/BlokEmailType.tsx`: opsi UW; Decline UW tanpa Cause of decline dan Comment tetap diisi; baris Validity (hari +
  dua tanggal dari server, dihitung ulang saat lepas fokus); konfirmasi UW hanya Decline (tombol Reject UW Pega tanpa modal).
- `nbfacin.css`: `.nbf-et__validitas`. Uji: `BlokEmailType.test.tsx` (blok "tahap Underwriting (tiket 61)").

## Backend (diminta ke sesi c3, 08-10-2026)

Lihat pesan kontrak tiket 61: tahap `uw`, `POST …/submit/validitas`, routing, limit, proteksi, riwayat.

## Terbuka

- Mata uang `TotalTSINusaRe` pembanding limit; arti argumen addCalendar WaitingBind.
- Data DEV `M_LIMIT_PROPERTYY`: tidak ada baris JUW_A; SENIORUW team 1 ganda.

## Backend sesi c3 — selesai 08-10-2026

Kode: `backend/services/submituw.go` (baru), `services/submit.go`, `services/proteksi.go` (`langkahPreferred`),
`repository/submit.go`, `models/submit.go`, `handlers/submit.go`, `handlers/handlers.go`.

- **K61-B1 Tahap uw.**
  - Syarat: When IsUWSpread (`tahapUW`) + workbasket POSITION + FlagOnGoingPolicy "0" (`tahapKonteks`, `alasanTidakBoleh`).
  - Pega menampilkan EmailTypeUW juga pada flag 2 (Binding); di sini hanya 0 (D7) — beda dari Pega, di luar cakupan.
  - Email Type uw 1 / 2 / 3 / 7; selain itu 400. Endpoint tab lain tetap seperti tiket 59 / 60.
- **K61-B2 Kontrak.**
  - `TampilanSubmit` + `daysValidity` / `dateValidity` / `waitingBindDate` (14 kunci, dua arah dengan api.ts).
  - `IsiSubmit` / `IsiPeriksaSubmit` + `daysValidity`. Handler `DisallowUnknownFields`, jadi tanpa ini badan frontend baru akan
    400.
  - Baru `POST …/submit/validitas {daysValidity}` -> `{dateValidity, waitingBindDate}`: mode edit, tidak menyimpan, penjaga
    Submit uw. "" / rusak -> 400.
- **K61-B3 Validity Of Acceptance** — `[dugaan]` diuji empiris (D2, `[keputusan agent, ikut Pega]`; perintah audit
  `PYTHONIOENCODING=utf-8 py docs/issues/61-audit/cekvaliditas.py` dan `cekhari.py` — hanya mencetak hitungan):
  - Rumus GetDateValidity_ACT langkah 2 / 4.
  - addCalendar: arg1 tahun, arg2 bulan, arg4 hari. DateValidity cocok 69/69 ekspor `DDL\CONTOH`.
  - Selisih hari dipotong (2/2 contoh yang membedakan; contoh lain 0 / 30).
  - ⚠️ Batas pengetahuan: WaitingBind L701 cocok 40/56 (ekspor mayoritas EDM).
  - LastAksep = TGL_INP riwayat HISTORYAKSEPTASIPRODUCTION terbaru sebelum Submit. IsTBA tidak ada -> bukan "true".
  - Tanggal disimpan teks Pega DateTime 00:00 WIB ("yyyyMMddT170000.000 GMT", seperti ekspor). Kabel DD-MM-YYYY.
  - GET: tersimpan bila ada, selain itu mode bawaan.
  - > **RALAT 08-10-2026 (perubahan kontrak, permintaan work owner lewat sesi 9d: "di pega bukannya saat uw pilih accept
    > bukannya nilainya lngsng ke isi ya otomatis?")** — GET kini mengembalikan nilai TERSIMPAN saja ("" bila belum). Mode
    > bawaan dihitung `POST …/submit/validitas` dengan `daysValidity` "" (Pega: klik radio `.EmailTypeUW` -> GetDateValidity_ACT
    > edit="" `[terverifikasi]` EmailSection_IsUW); "" tidak lagi 400. Jawaban validitas + `daysValidity`
    > ({daysValidity, dateValidity, waitingBindDate}). Submit dengan `daysValidity` "" sudah menghitung mode bawaan sebelum
    > menyimpan (tidak berubah).
- **K61-B4 Accept: jenjang UW** (D3, `[keputusan agent, ikut Pega]`): GetLimitAkseptasi_JUW_UW + Decision42 (`acceptUW`).
  - Posisi -> jabatan: JUW B -> JUW_B, JUW A -> JUW_A, UW -> UNDERWRITER.
  - Tabel limit menurut IsPreferredRisk: derivasi bersama proteksi (`langkahPreferred` / `preferredKasus`).
  - TeamGroup lewat MOID.
  - CARID2 = round0(TotalTSITopRisk) bila ada lokasi top risk (Σ TSILiability coverage pertama × PS / 100, TANPA kurs —
    mata uang campur, ikut Pega). Selain itu round0(TotalTSINusaRe): kurs TERBARU (KursTerbaru); BEDA dari tab Cedant yang
    memakai kurs tanggal polis. Syariah di luar cakupan.
  - Jenjang dihitung di Go atas baris `MuatLimit` (`jabatanBerikutUW`), tanpa SQL baru:
    - jabatan sendiri tidak ada / ganda ber-bottom beda -> Else (+ log);
    - seri -> urutan tabel;
    - batas TSI eksklusif.
  - Decision42 -> posisi / NB status (`tujuanAcceptUW`). Senior / DepHead memakai teks konstanta D7.
  - Else / Senior / DepHead / UW Financial / kurs tidak ada -> 422 `pesanJalurAkhir` (teks sistem baru) atau pesan kurs, tanpa
    tulisan.
  - Data DEV: tanpa baris JUW_A -> Accept JUW A selalu ditahan. Team tanpa UNDERWRITER -> SENIORUW (diuji).
- **K61-B5 Proteksi Accept uw** (D4, `proteksiSpreadingUW`): CheckSpreadingProtect(Fire)_ACT, pesan verbatim, indeks
  `.pxListSubscript` mulai 1.
  - M1 baris share / treaty kosong; M2 Σ share ≠ 100; M3 Σ PremiumSpreaded ≠ PremiNusantaraRe (4 desimal).
  - M4 "Total TSI Spreaded for …" (ringkasan spreading yang ada) dan M5 "Invalid Treaty Type Spreading!" (teks K59-8; korpus
    kosong di jalur ini).
  - Layering (CoverageBasis 5) dilewati. M6 / M7 dilewati (data tidak ada); M8 (EDM) dilewati.
  - Dijalankan di periksa (emailType 1) dan Submit Accept. Pesan proteksi lalu pesan jenjang dalam satu 422.
- **K61-B6 Tulisan** (satu transaksi, `TulisSubmit`):
  - HISTORYAKSEPTASIPRODUCTION (5.3) dan HISTORYAKSEPTASIPEGA (8, Workbasket = PositionNote).
  - Reject / Ask -> ReasFacInMarketing + "NB IS IN <MarketingName>'S INBOX" (Transition170, sama dengan TL).
  - Decline -> Decline / Resolved-Decline / "NB WAS DECLINED BY  <NAMA>" (dua spasi L3474), JUGA di JUW A (penyimpangan sadar
    butir c).
  - Migrasi 212: DAYS_VALIDITY / DATE_VALIDITY / WAITING_BIND_DATE / COMMENT_TEXT ditulis setiap Submit uw. Comment 500 bita
    -> 400; tetap ke riwayat.
  - FACINPERFORMANCE (D5):
    - Accept -> OFFER/ACCEPT oleh pelaku + OFFER→BINDING/OUTSTANDING oleh pembuat case (T_WORK_POLIS.CREATE_OP). Langkah 13.4
      (penggantian nama ter-hardcode) TIDAK diport — penyimpangan.
    - Decline -> OFFER/DECLINE. Reject / Ask -> tidak ada.
    - Baris berkolom wajib kosong dilewati.
  - Di luar cakupan (D6): SetValidateDateUW_PostAct 3 (SetDataFacOut_Act), 4 (GetTeamGroup_Act), 10 (UpdateJSONOffer_ACT),
    11-14.
- **K61-B7 Migrasi 212** (DITULIS): `ALTER T_GENERAL_POLIS ADD DAYS_VALIDITY VARCHAR2(50), DATE_VALIDITY VARCHAR2(30),
  WAITING_BIND_DATE VARCHAR2(30), COMMENT_TEXT VARCHAR2(500)` (rancangan) + down. **WAJIB sebelum biner**: GET / POST …/submit
  membaca kolomnya (ORA-00904 tanpa 212). Loader putar-ulang 342 kolom.
- **Uji:**
  - `services/submituw_test.go`: validitas, jenjang tabel, matriks posisi × Email Type, ditahan / 400 / 403 / 409, proteksi
    M1-M5, CARID2, periksa / GET / POST validitas.
  - `handlers/submit_test.go` `TestSubmitUWHandler` (HTTP, kunci dua arah TampilanSubmit / ValiditasSubmit).
  - `repository/submit_test.go` (SQL 212, konteks).
  - Mutasi: 20 valid, 19 tertangkap (3 sesudah uji diperkuat: batas flag di GET, DateToUW = Begindate, Comment 501 bita).
    `sort.Slice` vs `SliceStable` = mutan ekuivalen untuk n ≤ 12 (insertion sort stabil).
- **K61-B8 Catatan review dua sumbu (08-10-2026):**
  - Periksa kini memakai penjaga yang sama dengan Submit (403 workbasket / 409 tahap) untuk SEMUA tahap.
  - TeamGroup kosong -> Else (tidak mencocokkan baris TEAM_GROUP NULL; SQL Pega `team_group = {TeamGroup}` tidak mengembalikan
    apa pun).
  - Urutan pesan 422: proteksi spreading lalu jenjang (urutan Pega: validasi radio sebelum flow), digabung.
  - IsPKSASM (IsB2B "ASM") tidak pernah benar: IsB2B tidak ada di sistem baru.
  - Kurs terbaru tidak ada / item tanpa mata uang -> Accept ditahan dengan pesan agent "Kurs <X> ke IDR tidak ada di
    TREATYEXCHANGEYEARLY - limit akseptasi tidak dapat dihitung" `[keputusan agent]`.
  - daysValidity di luar tahap uw diabaikan (tidak 400).
  - Hari bawaan: Local.days bertipe Integer (GetDateValidity_ACT L228-270), jadi selisih hari dikonversi ke bulat SAAT
    ditetapkan (langkah 2.2), sebelum +1 / batas 30. Pemotongan `[dugaan]` (2/2 contoh) — sisa risiko bila Pega membulatkan.

## RALAT 08-10-2026 — IsGroup di GetLimitAkseptasi_JUW_UW (sesi c3)

- **Keputusan work owner 08-10-2026** (lewat sesi 9d, atas pertanyaan terbuka tiket 69a-1): *"Selaraskan, ikut Pega"*.
- `[terverifikasi]` `GetLimitAkseptasi_JUW_UW` langkah 8 (When IsGroup -> exit, LetterNo "") berlaku juga untuk Utility8 di UW
  tahap Offer. Port tiket 61 belum menerapkannya: case group tetap menghitung tabel limit lalu naik jenjang.
- **Kini**: `services/submituw.go` `letterJUWUW` dipakai bersama oleh `acceptUW` (Offer: JUW B / JUW A / UW) dan
  `acceptPolicy` (tiket 69a). Case group -> LetterNo "" -> Decision42 Else -> Decision24:
  - IsFacRetro -> ditahan `pesanFacRetro`;
  - selain itu ActFlow (langkah 3 exit group) -> Decision23 Else -> Decision16 IsGroup -> ditahan `pesanPolicyGroup` (tiket 64).
- **Akibat**: Accept JUW B / JUW A / UW atas case group (When IsGroup A pembuat / B-D akun / E kode marketing) tidak lagi naik
  ke UW / Senior / DepHead. Hasilnya 422 tanpa tulisan sampai jalur Policy group tersedia. Urutan cek: kurs (`pesanKurs`) lebih
  dulu, lalu IsGroup.
- **Uji** `TestAcceptUWGroup` (JUW B / UW bukan group naik jenjang; JUW B / UW group dan group pembuat case -> `pesanPolicyGroup`;
  group + FACOUT -> `pesanFacRetro`). Mutasi 3 / 4 mati. Yang hidup ekuivalen: tabel IsPreferredRisk "" tidak cocok baris
  mana pun, jadi hanya teks log yang berbeda.
