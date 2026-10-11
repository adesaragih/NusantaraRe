# 54: Tab Scoring Risk kasus FIRE

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 06-10-2026 "lanjut ke tab
> Scoring Risk" (tab Loss Record dilewati). Survei korpus oleh subagen (pengurai XML); butir kunci diperiksa ulang sesi 0f
> (`SetScore_Act` FinalScore, `ScoringResult` Good, kondisi `CallDualScoringRisk`, `Protection_Act` pesan,
> `ValueScoringRisk_act` poin Occupation).

**Status:** frontend selesai 06-10-2026 (uji hijau); backend (`GET/PUT …/scoring`, `POST …/scoring/hitung`) diserahkan
ke sesi c3 — menunggu "stabil scoring" lalu diverifikasi sesi 0f.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- Dtl tab "Scoring Risk" (L99417, `IsFire`, aktif `IsTABActive_ScoringRisk`) → `Section\CallDualScoringRisk.xml`:
  - "Write" (`IsFire && pyWorkPage.PositionNote = 'ReasFacInMarketing' && pyWorkPage.FlagOnGoingPolicy = 0`):
    `ScoringRisk` dapat diisi + Save; pre-load `ValueScoringRisk_act`.
  - "Read" (`IsFire && (PositionNote != 'ReasFacInMarketing' || (… == 'ReasFacInMarketing' && FlagOnGoingPolicy == 2))`):
    `ScoringRisk` baca-saja.
- `Section\ScoringRisk.xml` + `ScoringRiskForm1..4`: kepala (The Insured, Top Risk, Dominant Risk autocomplete dari
  LocationList, Period, Risk Location 1/2, Status, Score); 8 faktor berbobot (Occupation 15, Risk Location 10, Object
  Conditions 12,5, FEA 10, Net Loss Ratio 22,5, Risk Improvement 10, Survey Report 10, Others 10 = 100) dengan 24 sub-faktor,
  tiap sub-faktor = centang saling-eksklusif per kolom (Top Risk / Dominant Risk); hasil (FinalScore, kelas Good 86–100,
  Above Average 76–85, Average 60–75, Below Average 30–59, Poor 0–29; rekomendasi; persetujuan Operational Div. Head /
  Operational Director / Technical Director / President Director karena Occupation / Object Conditions / Loss Ratio / Others).
- `ValueScoringRisk_act`: tabel poin per centang (mis. Occupation 15 / 10 / 5 / -5 / -15) + `SetDataScoringRisk_act`
  (isi otomatis Top Risk: okupasi dari `POOLDATA.OCCUPATION.KDRISKEXPOSURE` 01–05, TSI dalam IDR via
  `pooldata.getcurrencystandard` ≤ / > 30.000.000.000, zona gempa dari coverage Zone "Daerah Zona 1..5").
- `SetScore_Act`: skor sub-faktor = poin bila dicentang; Score kolom = Σ; `FinalScore =
  @if(.DataScoringRiskList(2).OccupationNote="",local.TotalScore,@divide(local.TotalScore,2,2))`; `SetLossRatio_Act`:
  % loss ratio → band; `ScoringResult`: kelas, rekomendasi, persetujuan.
- `Protection_Act` langkah 18: FIRE non-EDM tanpa FinalScore → "Scoring Risk Cannot Empty !!" (saat submit).
- Penyimpanan Pega: halaman `pyWorkPage.OfferFacIn.ScoringRisk` (BLOB / JSON kasus), bukan tabel relasional.

## Cacat Pega yang ditemukan `[terverifikasi]`

- `SetScore_Act` 4.4: ScorePerFactor Earthquake kolom 2 memakai `.DataScoringRiskList(1)`; 4.5 Tsunami menjumlah
  `.RiskLocation.Tsunami.Score4` dari halaman induk.
- Muat tab mencentang ulang okupasi / TSI / zona tetapi tidak menghitung ulang skor (skor bisa basi).
- Kurs TSI: `getcurrencystandard` mengembalikan 1 bila kurs tidak ada.

## Keputusan work owner 06-10-2026 ("1. itu gbr menu scoring risk 2. hanya marketing 3. dibetulkan 4. disamakan dengan cedant")

1. Gambar layar Pega diterima — tata letak mengikutinya.
2. **Hanya marketing** = mode "Write" Pega: posisi kasus tahap marketing (`T_WORK_POLIS.POSITION`, padanan `PositionNote =
   'ReasFacInMarketing'`) dan bukan on-going (`FlagOnGoingPolicy = 0`). Backend mengirim `bolehUbah`; selain itu baca-saja.
3. **Cacat Pega dibetulkan** (backend): Earthquake kolom 2 memakai kolom 2; Tsunami tanpa `Score4` induk; muat tab
   menghitung ulang skor; pembandingan angka, bukan string.
4. **Kurs TSI = Cedant (C-8)**: `TREATYEXCHANGEYEARLY` per Begin date; kurs tidak ada → tidak dianggap 1.

## Frontend (sesi 0f)

- `components/TabScoring.tsx` (+ `TabScoring.test.tsx`), dipasang di `pages/InwardFacultative.tsx` (tab 'Scoring Risk',
  kasus FIRE); `labels.ts` `SCORING` / `FAKTOR_SCORING` / `HASIL_SCORING` / `TEKS_SCORING` (teks verbatim, diuji sebagai
  teks elemen utuh di `ScoringRisk.xml` + `ScoringRiskForm1..4`); `api.ts` `KolomScoring` / `HasilScoring` /
  `TampilanScoring` / `ambilScoring` / `hitungScoring` / `simpanScoring`; `nbfacin.css` `.nbf-score*`.
- Seluruh hitungan (poin, Score kolom, FinalScore, band Loss Ratio, kelas, rekomendasi, persetujuan, isian otomatis)
  di backend; layar mengirim pilihan lalu menampilkan jawaban `…/scoring/hitung` (jeda 300 ms, jawaban basi dibuang).
- Isian % Loss Ratio terkunci bila ChechBox 1 / 2 / 3 / 10 `[terverifikasi]` `ScoringRiskForm2.xml` `pyDisabledWhen`
  (`.DataScoringRiskList(n).LossRatio.LossRatio.ChechBox1 = "true" || … ChechBox10 = "true"`).
- Kolom Dominant Risk terkunci selama Dominant Risk belum dipilih.
- **Dirapikan 06-10-2026** (work owner "itu tampilannya bisa km rapikan ga biar lebih user friendly"):
  - kartu ringkasan (Top Risk, Dominant Risk, Score + kelas berwarna, Recomendation);
  - kartu info kasus;
  - satu kartu per faktor dengan bobot dan subtotal Top / Dominant di kepalanya;
  - judul sub-faktor sebagai baris pemisah, baris pilihan terpilih disorot, skor tampil sebagai lencana di samping
    centang, tanda "Diisi sistem" pada pilihan terkunci;
  - kartu hasil (Risk Classifications / Recomendation / Required approval) dengan butir aktif disorot.

  Isi, urutan, dan teks tetap Pega. Subtotal per faktor = Σ skor sub-faktor (`jumlahDesimal`, eksak) — hanya tampilan,
  bukan hitungan baru. Uji render statis memakai prop `awal`.
- **Dirapikan ulang 06-10-2026** (work owner: atas & tengah "jadi makin kurang user friendly, yg bwh nya sudah sesuai"):
  - atas = satu blok kepala dua kolom berlabel kiri seperti Pega (kartu ringkasan atas dan kartu info dihapus);
  - tengah = SATU tabel menerus (bukan 8 kartu): baris pita per faktor dengan subtotal lurus di bawah kolom Top /
    Dominant, judul sub-faktor menjorok, pilihan menjorok lagi, centang + skor di slot tetap, baris padat tanpa
    selang-seling; lebar penuh sama dengan bagian bawah (permintaan work owner, batas 1040px dicabut);
  - bawah tidak diubah.

  Diperiksa visual lewat Chrome headless atas data sintetis (pratinjau statis, berkas uji sementara dihapus sesudahnya).
  Label pembungkus centang di-nol-kan margin-nya — aturan global `:where(label:not([class]):has(> input))` memberi
  `margin-bottom: 12px` sehingga centang tidak sebaris dengan teks.

## Keputusan agent

- **S54-1** Centang per sub-faktor saling-eksklusif: memilih centang lain langsung memindahkan pilihan; mengklik yang
  terpilih = mengosongkan. Pega mengunci centang lain sampai yang terpilih dilepas (`pyDisabledWhen` per ChechBox) — hasil
  akhirnya sama (paling banyak satu centang), langkahnya lebih sedikit.
- **S54-2 — diputuskan work owner 06-10-2026: "ditiru persis seperti Pega".** `[terverifikasi]` `ScoringRiskForm2.xml`
  memuat aturan silang: pilihan "Not Available (Growing Trees)" Portable Fire Extinguisher (pilihan 3, kolom 1 dan 2)
  serta Hydrant / Sprinkler (pilihan 2, **hanya kolom 1**) **dikunci bila Class Construction = Growing Trees**
  (`.ObjectConditions.ClassConstruction.ChechBox4 != ""`). Logikanya tampak terbalik dan tidak simetris antar kolom
  `[dugaan]`, tetapi ditiru apa adanya: `TabScoring.tsx` `terkunciSilang` (hanya mengunci, tidak mengosongkan centang
  yang sudah terpilih — seperti `pyDisabledWhen`). Frontend saja; Pega tidak memvalidasinya di server.
  Catatan: kondisi Sprinkler kolom 1 di korpus kehilangan `||` (`ChechBox1 != "".DataScoringRiskList(1)…`); cara Pega
  mengevaluasinya `belum terverifikasi` — diterapkan sesuai maksud bagian Growing Trees-nya.
  Audit: `tiket54_silang.py` (mengurai sel ChechBox + `pyDisabledWhen` terdekat di `ScoringRiskForm1..4`; 174 sel) dan,
  cara kedua, uji `TabScoring.test.tsx` (mencacah seluruh `pyDisabledWhen` yang memuat `ClassConstruction.ChechBox4`) —
  keduanya: PFE 1, PFE 2, Hydrant 1, Sprinkler 1.

## Backend (sesi c3, 06-10-2026)

**Status:** backend selesai; migrasi 205 DITULIS, belum dijalankan. Kontrak = `frontend/api.ts` `TampilanScoring` (sama
persis). Rule dibaca ulang sendiri dengan pengurai XML (`Activity\ValueScoringRisk_act`, `SetDataScoringRisk_act`,
`SetScore_Act`, `SetLossRatio_Act`, `ScoringResult`, `SetDataOccupationScoringRisk_Act`, `RDBList\GetRiskExposureOccupationFire`,
`ReportDefinition`/`Section\ScoringRisk.xml` autocomplete Dominant).

- **K54-1 Penyimpanan** — workbook rancangan SUDAH punya tabelnya (aturan loader V-30): `T_SCORINGRISK` (induk
  `T_GENERAL_POLIS`, satu per kasus) → `T_SCORINGRESULT` (16 bendera teks) dan `T_DATASCORINGRISKLIST` (SEQ_NO 1 Top / 2
  Dominant) → `T_SCORING_FACTOR` (24 sub-faktor, FACTOR_GROUP / FACTOR_NAME = nama halaman Pega) → `T_SCORING_OPTION`
  (ChechBoxN / ScoreN). Migrasi `205_scoring_risk.sql` (+ `_down`) membuat kelimanya UTUH + sequence. Amandemen loader:
  `INDEX_LOCATION` (0 sampel) dan `PERCENT_LOSS_RATIO` (2 sampel, sebelumnya penampung); lebar TOP_RISK_LOCATION 4000,
  OCCUPATION_CODE / OCCUPATION_NOTE 1000 (lebar sumbernya, pola butir 80). Simpan = ganti utuh lima tabel di satu transaksi
  (≈ 260 sisipan per Save; ID opsi langsung `SEQ.NEXTVAL`). Pilihan dibaca balik dari opsi `CHECK_BOX = 'true'`.
- **K54-2 bolehUbah** — padanan `PositionNote = 'ReasFacInMarketing'`: `T_WORK_POLIS.POSITION` = **"Offer"** (posisi awal
  case NB butir 77 = tahap marketing sistem baru; belum ada kode yang memindahkannya) ATAU "ReasFacInMarketing" apa adanya
  (bila data lama memuatnya `[dugaan]`), dan `FLAG_ONGOING_POLICY = "0"`. Selain itu baca-saja (termasuk celah Pega marketing
  + flag 1). GET selalu boleh; `POST …/hitung` dan `PUT` saat baca-saja → **409**.
- **K54-3 Hitung** — GET dan setiap hitung SELALU menghitung ulang dari masukan (pilihan, persen, indexLocation). Poin =
  `ValueScoringRisk_act`; ScorePerFactor = poin centang; Score = Σ 24; FinalScore = Score Top, atau (Top + Dominant) / 2
  dibulatkan 2 desimal setengah ke atas `[dugaan]` (Pega `@divide(x,2,2)`). Kolom Dominant hanya bila lokasi
  `indexLocation` ada DAN ber-okupasi (Pega membuang `DataScoringRiskList(2)` ber-OccupationNote kosong); indexLocation bukan
  ObjectNo lokasi kasus → 400. **Instrumen:** pola centang 5 kasus `DDL\CONTOH` ber-FinalScore (nomor centang saja) diuji
  terhadap Score / persetujuan tersimpan Pega — 98.5 / 61.0 / 35.5 / 69.5 cocok; NB-181231 Pega 69.0 = skor basi (TSI
  dicentang ulang saat muat, skor tidak) → dihitung ulang 67.0 (cacat yang dibetulkan). Sensus dua cara: 119 / 120
  ScorePerFactor sampel cocok dengan tabel poin (pengurai XML), sama dengan jumlah Score per kasus.
- **K54-4 Isi otomatis per kolom dari lokasinya** (Pega: Top dari lokasi IsTopRisk — satu-satunya lokasi bila hanya satu —
  Dominant dari ObjectNo pilihan):
  - Occupation (dikunci) = KDRISKEXPOSURE "01".."05" → 1..5 (`GetRiskExposureOccupationFire`, OLDID = OccupationList(1)
    .OccupationId, TYPE 'FIRE'); kode lain / tanpa baris → 0.
  - Total Sum Insured (dikunci) = Σ per mata uang TSI lokasi × kurs TREATYEXCHANGEYEARLY pada Begin date (pembaca Cedant
    K49-11; Pega menjumlah TSI LINTAS mata uang mentah lalu dikali kurs mata uang terakhir — dibetulkan): ≤ 30.000.000.000 →
    1, > → 2 (numerik). Kurs tidak ada → pesan "Kurs … pada Begin date tidak ada di TREATYEXCHANGEYEARLY - Total Sum Insured
    <kolom> tidak diisi" dan TSI 0.
  - Earthquake = Zone coverage TERAKHIR berzona di lokasi kolom itu: "Daerah Zona 1..5" → ChechBox2..6, zona lain → 0 —
    hanya bila pengguna belum memilih (pilihan 0).
  - Loss Ratio = band `SetLossRatio_Act` dari persen: (0,10] 4, (10,25] 5, (25,50] 6, (50,75] 7, (75,100] 8, > 100 9 —
    menggantikan pilihan lain; persen kosong / ≤ 0 → pilihan 4–9 dikosongkan, 1–3 / 10 (manual) tetap.
  - OccupationNote / OccupationCode / TopRiskLocation = OccupationList(1).OccupationName / Id, RiskLocation.ASMAddress.
- **K54-5 ScoringResult VERBATIM** — kelas (Good 86–100 … Poor 0–<30; FinalScore di luar 0–100 → tanpa kelas, catatan
  "Poor"), "Due to" dan "Approval From" dievaluasi bergiliran kolom 1 lalu 2 (urutan Pega berpengaruh: persetujuan kolom 1
  terlihat kolom 2), TechnicalDirector memakai penjaga OperationalDirector (salin-tempel Pega; hasilnya SELALU =
  OperationalDirector). Tidak dibetulkan karena kontrak meminta verbatim — `[pertanyaan terbuka]` work owner bila
  Technical Director memang dimaksudkan berbeda.
  ⚠️ `TotalRNMShare` (syarat Operational Director bersama Occupation 4) = Σ **PremiumSpreaded** `TotalTSIPremiSpreadRNM`
  lokasi kolom itu × kurs (Pega: premi, bukan TSI, meski ambangnya 30 M seperti TSI — `[pertanyaan terbuka]`), dibandingkan
  NUMERIK. Akibat nyata: Pega membandingkan TEKS (`"4500000000" > "30000000000"` benar), sehingga kasus Occupation 4 seperti
  NB-181227 berpersetujuan Operational / Technical Director di Pega tetapi TIDAK di sini bila premi spread < 30 M IDR.
- **K54-6 Okupasi sama** — kolom Dominant ber-OccupationNote sama dengan Top → pesan "Can't Choose The Same Occupation"
  (verbatim) di `hitung`; `PUT` → 400 dengan pesan itu (Pega memasangnya sebagai pesan medan yang menghalangi Save).
- **K54-7 Di luar kontrak / tidak dibangun** — `.Remarks` sub-faktor dan teks keterangan ScoringResult (OccupationR1/R2,
  ObjConditionR1, LossRatioR1/R2, RateR1/R2, …) tidak ada di kontrak; kolom `REMARKS` ditulis NULL. `Protection_Act` langkah
  18 "Scoring Risk Cannot Empty !!" milik submit (tahap lain) — catatan: FinalScore kini SELALU terisi sesudah Save pertama.
- **K54-8 Urutan deploy** — 205 sebelum biner baru (GET tab memanggil `T_SCORINGRISK`; tanpa 205 → ORA-00942).
- **K54-9** ✅ *Terjawab 06-10-2026 — pengecualian ADR-0035 (bab "Keputusan work owner 06-10-2026 — menjawab K54-5, K54-9,
  K54-10").* **ADR-0035** ("Nol tanggal, ambang, atau nama tetap di dalam kode"): tabel
  poin `ValueScoringRisk_act`, ambang 30.000.000.000, batas kelas 86 / 76 / 60 / 30, dan band loss ratio 10 / 25 / 50 / 75 /
  100 kini literal di `services/scoring.go`, persis literal di activity Pega. Pilihan: catat sebagai pengecualian ADR-0035
  (preseden tiket 52 K52-11) atau pindahkan ke tabel tetapan (butuh migrasi + pemilik isi). Belum ada keputusan.
- **K54-10** (c3, telaah spec) Dicatat tanpa diubah:
  - Earthquake: Pega `SetDataScoringRisk_act` langkah 5 berjalan atas SELURUH lokasi / coverage, hanya menulis kolom 1, dan
    menimpa ChechBox2..6 setiap tab dimuat; di sini per lokasi kolomnya sendiri dan hanya bila belum dipilih (kontrak).
  - `TotalRNMShare.CARI1` Pega = satu halaman global berisi nilai TERAKHIR yang ditulis (lokasi Top saat muat, lokasi
    Dominant sesudah memilih); di sini share lokasi kolom masing-masing — ⛔ `[pertanyaan terbuka]` work owner.
  - Loss ratio: Pega hanya menyetel ChechBox4..9, sehingga centang manual 1–3 / 10 BISA berdampingan dengan band (skor
    keduanya dijumlah); di sini satu centang per faktor — band menggantikan pilihan manual (layar mengunci % bila 1/2/3/10
    dicentang, jadi hanya terlihat lewat API).
  - Kasus lama yang `POSITION`-nya berisi nilai marketing lain (selain "Offer" / "ReasFacInMarketing") menjadi baca-saja —
    `[pertanyaan terbuka]` (isi POSITION data lama tidak ada di korpus).
- **K54-11** (c3, perbaikan sesudah telaah) TSI dan premi spread dikonversi ke IDR TERPISAH: kurs mata uang yang hanya ada
  di premi spread tidak lagi mengosongkan centang TSI (dan sebaliknya); share tanpa kurs = kosong → syarat Operational
  Director "Occupation 4 && share > 30 M" bernilai salah. Lokasi Dominant tanpa Occupation → kolom 2 tidak dihitung (seperti
  Pega) DAN pesan "Dominant Risk lokasi <ObjectNo> belum punya Occupation - kolom Dominant Risk tidak dihitung" (teks
  buatan c3, Pega diam).

## Keputusan work owner 06-10-2026 — menjawab K54-5, K54-9, K54-10 ("setuju" atas rekomendasi sesi 0f)

Verifikasi sesi 0f sebelum keputusan: `go build` / `go vet` / `go test ./modul/nbfacin/... ./cmd/...` hijau; kunci JSON
`handlers/scoring.go` ↔ `api.ts` cocok dua arah (kolom 7, hasil 11, tampilan 11, lokasiDominan 2, badan 3); handler selalu
mengirim 2 kolom × 24 pilihan; migrasi 205 hanya CREATE 5 tabel + sequence/indeks (FK ke `T_GENERAL_POLIS`), `_down`
membalik semuanya. `[terverifikasi]` `Activity\ScoringResult.xml`: `TotalRNMShare.CARI1 > "30000000000"` (banding teks) dan
`TechnicalDirector = @if(Primary.ScoringResult.OperationalDirector=false, …)`.

1. **Ambang 30 M dibandingkan NUMERIK** (bukan teks seperti Pega) — tetap; sesuai keputusan "dibetulkan". Akibat yang
   diterima: kasus Occupation 4 dengan premi spread < 30 M IDR tidak lagi meminta Operational / Technical Director.
2. **`TotalRNMShare` = Σ premi spread (bukan TSI)** — ikut Pega; maksud sebenarnya tetap `[pertanyaan terbuka]` untuk
   Product+Underwriting (ambangnya sama dengan ambang TSI). Per lokasi kolom masing-masing (K54-10) — diterima.
3. **Technical Director = Operational Director** (penjaga salin-tempel Pega) — dibiarkan seperti Pega.
4. **Pengecualian ADR-0035** — tabel poin `ValueScoringRisk_act`, ambang 30.000.000.000, batas kelas 86 / 76 / 60 / 30,
   dan band loss ratio 10 / 25 / 50 / 75 / 100 tetap literal di `services/scoring.go`, tanpa tabel tetapan dan tanpa
   migrasi (preseden tiket 52 K52-11). Alasannya: nilai-nilai ini bagian dari aturan penilaian Pega, bukan tetapan
   operasional yang berubah menurut jadwal dagang.
5. **Kasus lama dengan `POSITION` marketing selain "Offer" / "ReasFacInMarketing" baca-saja** — diterima; diperbaiki bila
   nanti ditemukan nilai lain di data.

## Temuan 06-10-2026 — "Others" (Due to) menyala untuk TSI KECIL `[pertanyaan terbuka]`

Pemicu: work owner menanyakan kapan Required approval menyala. `[terverifikasi]` `Activity\ScoringResult.xml` langkah 3:
`Others = .Others.TotalSumInsured.ChechBox1=true`; `Section\ScoringRiskForm3.xml`: ChechBox1 = "Less than IDR
30,000,000,000"; `Activity\SetDataScoringRisk_act.xml`: ChechBox1 bila `local.TSIinIDR<="30000000000"`, ChechBox2 bila
`>"30000000000"`; `Activity\ValueScoringRisk_act.xml`: Score1 = -1, Score2 = 1. Akibat: SETIAP kasus dengan TSI ≤ 30 M IDR
mendapat "Due to: Others" → Operational Div. Head + Operational Director + Technical Director, dan rekomendasi
"Acceptable, required approval" (bila skor ≥ 30); TSI > 30 M tidak. Poin dan persetujuan konsisten satu sama lain, jadi
bisa disengaja — atau label / ambangnya terbalik `[dugaan]`. Backend c3 meniru apa adanya. Perlu jawaban
Product+Underwriting.

## Risk Scoring di General (work owner 06-10-2026: "hasil score di scoring risk, munculin di general")

`[terverifikasi]` `Section\Periode.xml` (`ASM-FW-GISFW-DATA-OFFERFACIN!PERIODE`, juga `Periode_IsUW.xml`): sel ber-label
"Risk Scoring" = `.ScoringRisk.FinalScore` (Integer, `pyDecimalPlaces` 2, baca-saja) + `.ScoringRisk.NoteFinalScore`
(baca-saja). Frontend (sesi 0f): `KasusNB.riskScoring` (`api.ts`), General menampilkan `angka4(finalScore) noteFinalScore`
(4 desimal, seragam dengan tab Scoring); Save di tab Scoring memperbarui General lewat `saatTersimpan` (jawaban PUT).
Isi = nilai TERSIMPAN `T_SCORINGRISK`, kosong bila belum pernah di-Save. Backend `GET …/kasus/{caseId}` diserahkan ke c3.
- **K54-12** (c3, 06-10-2026, permintaan work owner) `GET /api/nbfacin/kasus/{caseId}` (dan jawaban `PUT …/general`) membawa
  `riskScoring: {finalScore, noteFinalScore}` = `T_SCORINGRISK.FINAL_SCORE` / `NOTE_FINAL_SCORE` TERSIMPAN (subkueri ber-
  `PARENT_ID` = ID kasus di `sqlBacaKasus`; TIDAK dihitung ulang) untuk kolom "Risk Scoring" General (`Section\Periode.xml`
  `.ScoringRisk.FinalScore` / `.NoteFinalScore`, baca-saja). Belum pernah Save / bukan FIRE (tak ada baris) → "" dan "".
  finalScore teks desimal TM9 apa adanya (mis. "98.5"). ⛔ Akibat deploy: GET kasus kini MEMBACA `T_SCORINGRISK` — 205
  WAJIB jalan sebelum biner baru, kalau tidak SELURUH layar kasus gagal ORA-00942 (bukan hanya tab Scoring).
