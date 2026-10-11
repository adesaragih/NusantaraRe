# 62: Jenjang approval lanjutan sesudah UW (Kadiv / Direktur / Manager Teknik / DepHead)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 08-10-2026 (gambar layar pesan
> "Accept dari posisi ini meneruskan case ke jalur akhir … belum tersedia") "kenapa muncul gini, perbaiki" → "lengkapin sesuai
> flow nya InputInwardFacultativeOffer". Riset korpus tiga subagen (pengurai XML; konektor flow dua wadah 176 = 176; laporan
> lengkap di scratchpad sesi 0f `riset-A-approval.md`, `riset-B-facretro-group.md`, `riset-C-json.md`).

**Status:** backend selesai sesi c3 08-10-2026 (lihat "Backend sesi c3"); frontend memakai form UW tiket 61 apa adanya.

## Keputusan work owner 08-10-2026 (AskUserQuestion)

- **K62-1** Bertahap: tiket 62 jenjang approval lanjutan; 63 SAVE JSON_OFFER + FACINOFFER → marketing tahap Binding; 64 form
  Binding + Policy (RISlip); 65 FacRetro (FacOut). Jalur yang belum dibangun ditahan dengan pesan.
- **K62-2** NB Status posisi baru = nama jabatan (korpus berisi nama orang — tidak disalin): "NB IS IN KADIV TEKNIK'S INBOX",
  "… KADIV FACULTATIVE'S …", "… KADIV FINANCIAL'S …", "… DIREKTUR MARKETING'S …", "… DIREKTUR TEKNIK'S …", "… MANAGER TEKNIK'S …"
  (DepHead tetap K61-4).
- **K62-3** Email "Waiting for Bind" (tiket 63) ikut Pega: hanya terkirim di production (`IsPEGAPROD` → env var).
- **K62-4** JSON_OFFER dibuat (tiket 63) — meralat penundaan tiket 57.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Flow\InputInwardFacultativeOffer.xml`: Accept Senior (A7) / DepHead (A17) / UW Financial (A10) / Else jenjang UW (Decision42)
  → Decision24 (`IsFacRetro`) → Else → Utility2 `GetLimitAkseptasi_ActFlow` → Decision23 → A1 Kadiv Teknik
  (ReasFacInGroupLeader) / A17 DepHead (ReasFacInDepHeadUnderwriting, syarat non-ASM) / A8 Kadiv Facultative
  (ReasFacInFacultativeDivHead — konektor masuk TANPA SET PositionNote) / A11 Kadiv Financial (ReasFacInFinDivHead) / A6 Direktur
  Marketing (ReasFacInMarketingDirector) / A4 Direktur Teknik (ReasFacInTechnicalDirector) / A5 Manager Teknik
  (ReasFacInManagerTeknik); Else → Utility16 SAVE JSON (tiket 63).
- Semua level memakai flow action `InwardFacultative_IsUW` (form tiket 61). Accept tiap level → Utility2 lagi (LOOP; A17 →
  Decision24); Reject / Ask → marketing (Transition170); Decline → End4; A11 tanpa konektor Decline.
- `Activity\GetLimitAkseptasi_ActFlow.xml` (56 langkah, 11 nonaktif): CARID2 = TotalTSINusaRe (top risk: TotalTSITopRisk); keluar
  bila IsGroup; SQL per `IsPreferredRisk` (`GetLimitAkseptasi_SQL` / `…NonPrefer_SQL` / `…PreferedComm_SQL`): `team_group =
  TeamGroup AND LIMIT_BOTTOM > (LIMIT_BOTTOM baris LOGIN = OperatorID) AND LIMIT_BOTTOM < CARID2 ORDER BY LIMIT_BOTTOM`; langkah 18
  membuang baris KADIVFACULTATIVE (kecuali banding+reject) dan MANAGERTEKNIK; LetterNo = JABATAN baris pertama.

## Terbuka

- ⚠️ RALAT 08-10-2026: rancangan awal mencari baris limit pengguna lewat `LOGIN = OperatorID` (seperti SQL Pega). Itu
  bertentangan dengan KEPUTUSAN-30-09-2026 butir 33 (work owner: limit pengguna = masukan tangga dari MODEL PERAN, bukan
  `LOGIN`; pemetaan login → jabatan dipegang IAM) — ditangkap sesi c3. Yang dipakai: jabatan = peran POSITION case
  (`acceptance.Next`, tiket 11): JUW B→JUW_B, JUW A→JUW_A, UW→UNDERWRITER, Senior→SENIORUW, DepHead→DEPHEADUNDERWRITER,
  GroupLeader→KADIVTEKNIK, FacultativeDivHead→KADIVFACULTATIVE, FinDivHead→KADIVFINANCIAL, MarketingDirector→DIREKTURMARKETING,
  TechnicalDirector→DIREKTURTEKNIK, ManagerTeknik→MANAGERTEKNIK. Tidak ada ketergantungan pada kolom LOGIN.
- UW Financial tidak punya jabatan di tabel limit → Else → ditahan.
- A5 / A8 praktis tidak terpilih untuk FIRE karena langkah 18 — diikuti apa adanya.
- A11 tanpa Decline di Pega.

## Backend sesi c3 (08-10-2026)

Tanpa migrasi, tanpa CALL, tanpa endpoint baru — `POST …/submit/periksa` dan `POST …/submit` tiket 61 dipakai apa adanya.

- **K62-B1 Tahap Submit uw** — `services/submit.go` `tahapSubmitUW`: (IsUWSpread ATAU enam posisi lanjutan
  `posisiJenjangLanjut`) DAN `FLAG_ONGOING_POLICY` "0"; wajib memegang workbasket POSITION (403). Spreading / Akumulasi TETAP
  hanya `tahapUW` (IsUWSpread) — `TestSubmitJenjangLanjut`.
- **K62-B2 Rute Accept** — `services/submituw.go` `acceptUW`: JUW B / JUW A / UW = `GetLimitAkseptasi_JUW_UW` + Decision42,
  Else → `decision24`; Senior / UW Financial / DepHead → `decision24`; GroupLeader / FacultativeDivHead / FinDivHead /
  MarketingDirector / TechnicalDirector / ManagerTeknik → `actFlow` langsung (konektor Accept tiap level ke Utility2, bukan
  Decision24 — bukti di atas). `decision24`: spreading memuat baris TreatyType `10015` (FACOUT) = `IsFacRetro`
  (`Activity\SetDataFacOutFire_Act.xml` langkah 5 men-SET IsFacRetro = 1 untuk TreatyType 10015; langkah 6 IsInputFacRetro = 0)
  → ditahan "Accept meneruskan case ke jalur FacRetro yang belum tersedia." (tiket 65); selain itu `actFlow`. Galat program
  `actFlow` DITERUSKAN (bukan ditelan jadi pesan).
- **K62-B3 Jabatan pengguna = peran POSITION case** (butir 33, pengganti rancangan `LOGIN = OperatorID` — RALAT di bawah
  "Terbuka"): `jabatanPengguna`. `[pertanyaan terbuka]` padanan JABATAN untuk UW Financial tidak diketahui:
  `DDL\M_LIMIT_PROPERTYY.txt` hanya DDL (22 baris, 0 INSERT — `wc -l`, `grep -ci insert`), dan Decision23 tidak punya konektor
  ke UW Financial. `[keputusan agent]` tanpa jabatan → Else → ditahan + log.
- **K62-B4 `actFlow` = `acceptance.Next`** (tiket 11) atas tabel bentuk A `susunTangga` + `kasusTangga` (properti clipboard
  yang dibaca ActFlow dan When-nya). Decision23 → posisi + NB status K62-2 (`statusInboxPosisi`); Selesai (Else) → ditahan
  "Accept meneruskan case ke jalur akhir (simpan penawaran & Binding) yang belum tersedia." (tiket 63). Galat tangga yang
  berarti "tujuan tidak pasti" (`ErrDaftarGanda`, `ErrTabelTakTerpilih`, `ErrUrutanTakPasti`, `ErrLimitKosong`,
  `ErrBarisPenggunaGanda`, `ErrBentukB`, `rules.ErrTafsirBerbeda`) → Else + log, bukan 500.
- **K62-B5 `BusinessOldId`** — predikat bentuk B (`When\IsLimitSBondKBG.xml`, `IsLimitCreditCL.xml`, dst.) membandingkan
  `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId` dengan literal. `[terverifikasi]` Pega mengisinya dari
  `RDBList\GetOldIDBusiness_SQL.xml` (`select oldid as CARI1 … from business where (ID ={OldID.CARI2} OR NOTE ={OldID.CARI2})
  AND GROUPPANEL IS NOT NULL`, pxResults(1); `Activity\SaveJsonOfferFacIn_Act.xml` langkah "Search OldID business" / "Set OldID
  business to Pega Property"; juga `GenerateNopolis_Act.xml`). `[dugaan]` saat Accept UW properti itu sudah terisi karena
  SaveJsonOfferFacIn_Act juga dipanggil tombol `Section\InputInwardFacultative_2.xml`. `[keputusan agent]` Di sini dibaca
  segar: `repository.sqlOldIDBisnis` = `SELECT OLDID FROM BUSINESS WHERE NOTE = :1 AND BUSINESSGROUPID = :2` (pengenal baris sama
  dengan kode bisnis Table of Limit, butir 89 — bukan `ID OR NOTE` + `GROUPPANEL IS NOT NULL` Pega). 0 / > 1 baris / OLDID
  kosong → "" + log → `rules.ErrTafsirBerbeda` → Else (tidak menebak "bukan bond").
- **K62-B6 Penyimpangan sadar** — (a) A8 Kadiv Facultative: konektor Pega tidak men-SET PositionNote; di sini POSITION tetap
  ditulis (satu-satunya penanda posisi di sistem baru). (b) A11 Kadiv Financial: Pega tanpa konektor Decline; di sini Decline
  diizinkan seperti level lain (form tiket 61 sama untuk semua level).
- **K62-B7 Kotak masuk** — `antreanKasusNB` + ReasFacInUnderwritingFinancial, ReasFacInFinDivHead, ReasFacInManagerTeknik;
  `TestAntreanKasusNB` kini menuntut setiap tujuan Accept (`statusInboxPosisi`) dan posisi Decision24 / lanjutan ada di antrean.
- Riwayat / FACINPERFORMANCE / Validity / proteksi Accept (M1-M5) = tiket 61 tanpa perubahan.

**Uji** — `TestAcceptJenjangLanjut` 17 kasus (dua cara: blok struct `grep -c '^\t\t{"'` = 17 dan `grep -o` nama = 17, atas
rentang fungsi), `TestSubmitJenjangLanjut`, `TestAntreanKasusNB`, `TestTableOfLimit` (repository: SQL `OLDID`). Kasus lama
Senior / DepHead / Financial / JUW A kini mengharapkan pesan FacRetro (spreading uji memuat FACOUT). "409 posisi" di
`TestAjukanSubmitDitolak` memakai posisi fiktif (FacultativeDivHead kini tahap uw → 403).
gofmt / vet bersih; `go test ./modul/nbfacin/...` hijau; rules dengan `KORPUS_RNM_AKAR` hijau; penjaga hanya kegagalan lama
premiumlistlife (`handlers/rute_datapolis.go` mengimpor repository); index git 0.

**Mutasi** 12 / 13 mati: tanpa BusinessOldId, OLDID ganda diterima, ErrTafsirBerbeda bukan Else, FacRetro dibalik, lanjutan
tanpa ActFlow, tahap Submit tanpa lanjutan, jabatan Kadiv Teknik salah, jabatan Manager Teknik hilang, kolom OLDID, inbox tanpa
FinDivHead, inbox tanpa ManagerTeknik, OLDID tanpa trim. HIDUP: `kasusTangga` tanpa `pyWorkPage.PositionNote` — mutan
ekuivalen di jalur ini (dibaca hanya `bandingReject`, IsBanding selalu kosong, dan gerbang bentuk B yang sudah ditolak
`ErrBentukB`).


> ⚠️ **RALAT tiket 63 (08-10-2026, sesi c3):** pesan "Accept meneruskan case ke jalur akhir (simpan penawaran & Binding) yang belum
> tersedia." DIHAPUS. Decision23 Else dan posisi tanpa jabatan (UW Financial) kini → jalur akhir Binding (issues/63 K63-B1).
> Tangga tidak pasti (seri, bentuk B, BusinessOldId tak diketahui, …) tetap ditahan dengan pesan "Jenjang approval berikutnya
> tidak dapat ditentukan dari tabel limit; hubungi administrator." ActFlow langkah 3 "exit kalau group" kini memakai When
> IsGroup (K63-B2).
