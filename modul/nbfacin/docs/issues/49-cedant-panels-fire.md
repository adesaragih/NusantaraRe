# 49: Tab Inw Fac Cedant Panels kasus FIRE

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 05-10-2026 "lanjut saja dlu
> ke menu Inw Fac Cedant Panels". Frontend sesi 0f; backend sesi c3.

**What to build:** tab Inw Fac Cedant Panels: Source of Business, Share Cedant Type, % Share RNM, Total TSI RNM, Total
Premi RNM, grid cedant (Ceding · % Share; Add / Delete; baris dibuka = popup `CedingCedant`), proteksi Save.

**Status:** frontend selesai 05-10-2026 (sesi 0f); backend selesai 05-10-2026 (sesi c3, uji hijau). Migrasi 199 ditulis,
menunggu dijalankan work owner — sebelum biner baru.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\InputInwardFacultativeDtl.xml` tab "Inw Fac Cedant Panels" (baris 60510–67588, `pyContainerVisibleWhen
  !IsLife`, active when `IsTABActive_Cedant`):
  - Source of Business — `.OfferFacIn.QuotationData.SobName`, pxDisplayText;
  - Share Cedant Type — `.OfferFacIn.ShareCedantType`, pxRadioButtons horizontal; wajib bila `.FlagOnGoingPolicy ==2 ||
    pyWorkPage.OfferFacIn.QuotationData.IsGroup = "Group"`; klik → refresh + `SetShareOfCeding`;
  - % Share RNM — `.OfferFacIn.PercentShare`, pxNumber (properti yang sama dengan tab Spreading);
  - Total TSI RNM — `.OfferFacIn.TotalTSINusaReSpreading`; Total Premi RNM — `.OfferFacIn.TotalPremiNusaRe` (pxNumber);
  - grid `.OfferFacIn.CedingCedantList` (kelas `ASM-FW-GISFW-Data-Quotation`): kepala "Ceding" · "% Share" · tombol
    "Add" (`AddCedantList_act`, nonaktif bila `.IsOldData = 'old' && IsEdmPerubahan`); sel `.CedingCoName`
    pxDisplayText baca-saja; `.ShareCeding` pxNumber; tombol "Delete" (deleteRow, tampil bila `pyWorkPage.IsOldData!=1
    || !IsUW`); klik baris → editItem flow action `CedingCedant`.
- `Activity\SetShareOfCeding.xml`: `ShareCedantType==0` → `QuotationData.ShareOfCeding = "100%"`; `==1` →
  `PercentShare + "%"`; disalin ke `pyWorkPage.Quotation.ShareOfCeding`.
- `Activity\AddCedantList_act.xml`: grid kosong + `QuotationData.CedingCoList` berisi → salin semua CedingCo /
  CedingCoName; grid kosong + daftar kosong → satu baris dari `QuotationData.CedingCo` / `CedingCoName`; grid berisi →
  tambah baris kosong.
- `Section\CedingCedant.xml` (popup baris): Ceding · % Share (change → `SetTSIPremiCedant_Act`) · "OF"
  `QuotationData.ShareOfCeding` · grid `CurrencyList`.
- `Activity\SetTSIPremiCedant_Act.xml` / `CountPremiCedant_Act.xml`: `CedingCedantList(i).CurrencyList` per mata uang —
  type 1 (SHARE RNM): TSI = ShareCeding × (TSI × PercentShare / 100) / 100, Premium = ShareCeding × `.Policy.Payment.Premium`
  / 100; type 0 (GROSS): TSI = ShareCeding × TSI / 100, Premium = ShareCeding × (Premium × 100 / PercentShare) / 100;
  ShareCeding > 100 → "Value can't be more than 100 or less than 0!".
- `Activity\ProtectShareCedant_Act.xml` (dari `Protection_Act`): "% Share cedant more than 100 or less than 0!"
  (`.ShareCeding<=0||.ShareCeding>100`), "Ceding on list inward facultative cedant can't be null!" (`.CedingCoName==""`),
  "Total premi share cedant must be equal total premi RNM!" (presisi 4), "List Inward facultative cedant can't be null!";
  gerbangnya dipetakan di bawah.
- Contoh kasus `DDL\CONTOH\*.xml` (nilai tidak disalin): `ShareCedantType` bernilai 0 (65×) dan 1 (46×); baris
  `CedingCedantList` = CedingCo, CedingCoName, ShareCeding, CurrencyList(Name, TSI, Premium, TSIOld, PremiumOld).
- `DDL\ShareCedantType.xml` (dikirim work owner 05-10-2026; `ASM-FW-GISFW-DATA-OFFERFACIN!SHARECEDANTTYPE`): "0" Gross,
  "1" Share RNM.
- Choose Ceding: `Section\CedingCedant.xml` tombol "Choose Ceding" → showHarness `CedingCedant` (popup, jendela "Choose
  Ceding", param Index = `.pxListSubscript`) → `CedingCedantHierarki` → Choose `SetDataSobCeding_Act` cabang
  `btnQuotation=="CedingCedant"`: `CedingCedantList(SearchSOB.CARI2).CedingCo = .ID`, `.CedingCoName = .ClientName` —
  sumber sama dengan popup Ceding Company tab General (gambar layar work owner 05-10-2026).
- Gerbang `ProtectShareCedant_Act` langkah 1: IsLife → 6, EmailTypeBinding "4" → 6, IsGroup "Group" → 5,
  FlagOnGoingPolicy 2 → 5, IsEdmAdjShareCedant true → 2 / false → 6. `[dugaan]` kode 5 = jalankan langkah (arti kode
  belum terverifikasi) → proteksi hanya untuk Group / on-going.
- Total TSI / Premi RNM diisi `CountPremiAndTSIRNMFireMBU_ACT` (`TotalTSINusaReSpreading = PercentShare × totaltsi / 100`
  atau `× Tabarufund / 10000`; `TotalPremiNusaRe = Local.totalpreminusare` / `Local.AllPremiRNM`).

## Frontend (sesi 0f)

- `components/TabCedant.tsx` (`tambahCedant`, `shareSah`, `galatCedant`) + uji; `api.ts` (`ambilCedant`, `simpanCedant`,
  `TampilanCedant`, `BarisCedant`); `labels.ts` (`CEDANT`, `OPSI_SHARE_CEDANT_TYPE`, `TEKS_CEDANT`); `nbfacin.css`;
  `pages/InwardFacultative.tsx` (tab kasus FIRE).

## Kontrak USULAN (menunggu sesi c3)

- `GET /api/nbfacin/kasus/{caseId}/cedant` → `TampilanCedant` { sobName, shareCedantType ("" | "0" | "1"),
  percentShare, totalTsiRnm, totalPremiRnm, wajib, cedant[{cedingCo, cedingCoName, shareCeding}], cedingUmum[…] }.
- `PUT …/cedant` { shareCedantType, cedant } → simpan + `SetShareOfCeding` → jawab bentuk GET; 400 pesan verbatim.

## Keputusan agent (menunggu konfirmasi)

- **C-1 RALAT 05-10-2026** (gambar layar Pega detail baris dari work owner): detail baris `CedingCedant` DIBUAT sebagai
  panel di bawah baris (pola tab Object): Ceding + Choose Ceding; % Share* "of" ShareOfCeding; grid Currency · TSI
  Cedant · Premium Cedant; kontrak `shareOfCeding` + `dasarMataUang` (sesi c3, K49-8).
  > RALAT dasar: sesi 0f semula menduga cabang `TotalSpreadAll.pxResults` — keliru; langkah itu ber-label `//`
  > (dimatikan). Yang berjalan langkah 4 atas `OfferFacIn.CurrencyList` (`.TSI` + `.Policy.Payment.Premium` diisi
  > `CountPremiAndTSIRNMFireMBU_ACT`), lihat bagian Backend. Proteksi total premi kini di-port sesi c3 (K49-9).
  > Teks lama C-1 di bawah dibatalkan.
- ~~**C-1** popup baris `CedingCedant` tidak dibuat: Choose Ceding dan % Share ada di baris grid. Grid CurrencyList dan
  proteksi "Total premi share cedant must be equal total premi RNM!" menunggu tab Payment (disetujui work owner
  05-10-2026: "setuju").
- **C-2** % Share RNM baca-saja di tab ini (diubah di tab Spreading).
- **C-3** proteksi hanya bila `wajib` (gerbang di atas); Share Cedant Type bertanda wajib dalam kondisi yang sama.
- **C-8 (keputusan work owner 05-10-2026)** "kalau ada 2 currency, ditotalin jadi idr. ini yang sekarang cuma di total
  angka nya saja tanpa liat currency": Total TSI RNM / Total Premi RNM = Σ per mata uang × kurs TOIDR
  (`TREATYEXCHANGEYEARLY`, sumber yang sama dengan tab Spreading) — menyimpang sadar dari `CountPremiAndTSIRNMFireMBU_ACT`
  (Pega menjumlah angka tanpa kurs). Mata uang tanpa kurs → total "" + `kursTidakAda` (sesi c3).
  Sumber kurs dikonfirmasi work owner 05-10-2026: "kurs yang ada pada database berdasarkan begin date, nilai kursnya
  pakai tabel TREATYEXCHANGEYEARLY" → TOIDR yang periodenya memuat Begin date (sama dengan Spreading), tanpa jatuh ke
  kurs terbaru; `M_CURRENCYSTANDARD` (kurs harian) tidak dipakai.
- **C-7 (keputusan work owner 05-10-2026)** "ini masa bisa lebih dari 100%": Σ % Share cedant > 100 selalu ditolak (layar
  + diminta ke backend), tanpa gerbang `wajib`; pesan agent "Total % Share cedant tidak boleh lebih dari 100.". Kartu
  % Share RNM tanpa akhiran "%".
- **C-6 (keputusan work owner 05-10-2026, membatalkan C-5)** "saat di add salin aja dari data Ceding Co dari tab General,
  dan tombol choose ceding nya di hapus saja": Add menyalin `cedingUmum`; grid berisi → hanya ceding General yang
  belum ada (kode), tanpa baris kosong (tak dapat diisi tanpa Choose Ceding) — keputusan agent atas langkah 4
  `AddCedantList_act`. Tombol Choose Ceding tidak ada.
- ~~**C-5 (keputusan work owner 05-10-2026)** "saat klik add, ceding co dari tab general gausa di salin": Add selalu
  menambah satu baris kosong (menyimpang sadar dari `AddCedantList_act`); `cedingUmum` tidak dipakai layar.
- ~~**C-4** baris dari Add boleh berceding kosong sampai Choose Ceding; Save menolaknya bila `wajib`.~~ Gugur oleh C-6
  (Add hanya menyalin Ceding Co General; tidak ada baris kosong, tidak ada Choose Ceding).

## Backend (sesi c3, 05-10-2026)

Kontrak final = kontrak usulan di atas. `GET /api/nbfacin/kasus/{caseId}/cedant` → `TampilanCedant` persis `api.ts`
(larik selalu `[]`, angka teks desimal, `""` bila kosong; `cedingUmum[].shareCeding` = `""`). `PUT …/cedant`
{shareCedantType, cedant[{cedingCo, cedingCoName, shareCeding}]} (kunci asing → 400) → jawab bentuk GET. Galat: 400 isian /
proteksi (`ErrMasukanCedant`, pesan verbatim), 401 tanpa identitas, 404 case, 503 tanpa basis data.

- `backend/migrations/199_cedant.sql` (+ `_down`): `T_GENERAL_POLIS.SHARE_CEDANT_TYPE` NUMBER(5),
  `T_QUOTATIONDATA.SHARE_OF_CEDING` VARCHAR2(50), tabel `T_CEDINGCEDANTLIST` (rancangan, dibuat utuh) + indeks + sequence.
  **Ditulis, belum dijalankan** — wajib sebelum biner baru (GET tab Cedant membaca kolom baru). Penjaga rancangan
  `TestMigrasiFlatSebagianCocokRancangan` kini juga membaca 198 dan 199 (234 kolom); `docs/STRUKTUR-TABEL-NB-FACIN.md`.
- `models/cedant.go`, `repository/cedant.go` (`PenyimpanCedant`), `services/cedant.go`, `handlers/cedant.go` (+ uji ketiganya).

## Keputusan agent backend (menunggu konfirmasi)

- **K49-1** (c3) Penyimpanan = workbook rancangan `[terverifikasi]` (`loader/skema_gen.go`): `T_GENERAL_POLIS.SHARE_CEDANT_TYPE`
  (`ShareCedantType`, rancangan NUMBER → NUMBER(5), kode A170), `T_QUOTATIONDATA.SHARE_OF_CEDING` (`ShareOfCeding`,
  rancangan NUMBER, tetapi Pega menulis teks dan 94 dari 94 contoh `DDL\CONTOH` berbentuk "N%" → VARCHAR2(50) apa adanya,
  pola `T_TABLEOFLIMIT.PCT_LIMIT`), `T_CEDINGCEDANTLIST` (jalur `CedingCedantList`, induk `T_GENERAL_POLIS`; PARENT_ID
  VARCHAR2(32) = PK teks General, SHARE_CEDING NUMBER(38,8)). `IS_GROUP` (rancangan T_QUOTATIONDATA) dan
  `T_CEDING_CURRENCYLIST` tidak dibuat.
- **K49-2** (c3) Gerbang proteksi `ProtectShareCedant_Act` langkah 1 `[terverifikasi]` (When verbatim: IsLife T6/F2,
  `EmailTypeBinding=="4"` T6/F2, `QuotationData.IsGroup=="Group"` T5/F2, `FlagOnGoingPolicy==2` T5/F2, IsEdmAdjShareCedant
  T2/F6) di-port sebagai fungsi murni (`proteksiCedantBerlaku`); `[dugaan]` 6 = keluar activity, 5 = lewati When sisanya.
  Sumber di sistem baru:
  - `pyWorkPage.FlagOnGoingPolicy` = **`T_WORK_POLIS.FLAG_ONGOING_POLICY`** (VARCHAR2(1), milik premiumlistlife migrasi
    057; case NB ditulis "0" oleh `repository/casenb.go` `FlagAwalCaseNB`) — gerbang `== "2"`. Tidak ada di workbook
    rancangan, tetapi tersimpan (temuan review spec 05-10-2026; versi pertama keliru menyatakan tidak ada sumber);
  - `QuotationData.IsGroup` di Pega = "Group" bila **login pembuat case** termasuk daftar login literal
    (`Activity\GetTeamGroup_Act.xml`, `InputOfferFacInEngineer_preACT.xml`; nilai login tidak disalin, A14) — **belum ada
    sumber**, false sampai model peran (butir 33, A20 `AnggotaGrup`) (sisa);
  - `EmailTypeBinding` dan IsEdmAdjShareCedant (EDM) tidak berlaku untuk NB (false).

  Case NB baru (FLAG_ONGOING_POLICY "0") → proteksi tidak jalan dan Share Cedant Type tidak wajib, sama dengan Pega untuk
  case NonGroup. Saat gerbang bernilai true: Share Cedant Type kosong → 400 "Share Cedant Type wajib diisi" (pesan agent, mengikuti
  tanda wajib Section), lalu pesan verbatim urut Page-Set-Messages: "List Inward facultative cedant can't be null!" (grid
  kosong), "% Share cedant more than 100 or less than 0!" (share ≤ 0, kosong, atau > 100), "Ceding on list inward facultative
  cedant can't be null!" (Pega `.CedingCoName==""`; di sini kode ATAU nama kosong — nama yang disimpan tetap dari AGENT). ⚠️ Di Pega langkah 9 "% SHARE NULL" memasang `Local.ErrMsgNull` (bukan
  `Local.ErrMsg`), sehingga pesan share tidak pernah tampil walau Save tetap ditahan (`ProtectSpreading.CARI1 = 1`); di sini
  pesan yang dimaksud ditampilkan. "Total premi share cedant must be equal total premi RNM!" di-port kemudian (K49-9).
- **K49-3** (c3) Total TSI / Premi RNM tidak disimpan, dihitung tiap baca dari objek tersimpan, dengan NR dihitung ulang dari
  % Share RNM seperti tiap Save Pega: `TotalTSINusaReSpreading` = PercentShare × Σ `CoverageList(1).TSILiability` tiap item /
  100, `TotalPremiNusaRe` = Σ `PremiNusantaraRe` seluruh coverage (`CountPremiAndTSIRNMFireMBU_ACT` 6.2.2.1.2, 6.2.2.1.3.3,
  6.2.3; tanpa konversi kurs dan tidak dinolkan per mata uang, verbatim). Kosong bila % Share RNM kosong. Cabang Tabarufund
  (syariah, `× Tabarufund / 10000`, langkah 6.2.4 `Local.IsSyariah==1`) di luar cakupan (seperti K48-6). Syarat FireMBU yang
  tidak dibawa: `.FlagDelete!=1` (item / coverage; FlagDelete tidak dimodelkan, K48-6), `Local.Currency==Local.ObjectCurrency`
  (hanya item bermata uang yang ada di `OfferFacIn.CurrencyList` — di sini semua item, `[dugaan]` sama karena CurrencyList
  disusun dari mata uang item), dan rumus layering `CoverageBasis` 5 (dihitung dengan rumus non-layer, K48-6).
- **K49-4** (c3) SetShareOfCeding dijalankan saat Save (Pega: saat radio diklik): tipe "0" → "100%", "1" → % Share RNM +
  "%" (8 desimal, nol di belakang dibuang; % Share RNM kosong → "%", verbatim); tipe kosong → Share of Ceding tidak diubah.
  Salinan ke `pyWorkPage.Quotation.ShareOfCeding` tidak dibawa (halaman Quotation work bukan milik modul ini).
- **K49-5** (c3) `% Share` di atas 100 → 400 "Value can't be more than 100 or less than 0!" untuk semua case
  (`SetTSIPremiCedant_Act`, tanpa gerbang); angka negatif sudah ditolak pengurai desimal. Nama ceding ditulis dari AGENT
  menurut kode (pola tiket 33/34, E-4) — `cedingCoName` dari layar diabaikan; kode tak lolos AGENT → 400. Baris tanpa kode
  boleh disimpan bila proteksi tidak berlaku (C-4). Paling banyak 100 baris. ⚠️ Kode yang kini tidak lolos syarat AGENT
  (status / tipe, tiket 33) ditolak walau proteksi tidak berlaku — Pega `AddCedantList_act` tidak memeriksa (pola tiket 34).
- **K49-6** (c3) `cedingUmum` saat `CedingCoList` kosong: T_QUOTATIONDATA.CEDING_CO / CEDING_CO_NAME di sistem baru menyimpan
  gabungan `;` (tiket 34), bukan `QuotationData.CedingCo` tunggal Pega → dipecah per kode; nama dipasang menurut posisi hanya
  bila jumlah potongan sama (nama ber-`;` → dikosongkan, diisi dari AGENT saat Save). Share of Ceding ditulis dari % Share
  RNM tersimpan (NUMBER(38,8), paling banyak 8 desimal; contoh Pega memuat hingga ±17 desimal). Endpoint tidak memeriksa
  `!IsLife` / jenis FIRE — sama dengan tab lain modul ini (Spreading, Clauses).
- **K49-7** (c3) `cedingUmum` **tetap** di jawaban GET. C-5 (Add tidak menyalin Ceding Co General) dibatalkan C-6 work owner
  05-10-2026 ("saat di add salin aja dari data Ceding Co dari tab General, dan tombol choose ceding nya di hapus saja"): Add
  kembali menyalin `cedingUmum`, sehingga baris cedant berasal dari Ceding Co tab General (T_CEDINGCOLIST / gabungan `;`,
  K49-6). Backend tidak memaksa asal baris (layar yang menyusun); kode tetap divalidasi ke AGENT dan nama ditulis dari AGENT.
- **K49-8** (c3) Detail baris cedant (`Section\CedingCedant.xml`): GET/PUT menambah `shareOfCeding` (T_QUOTATIONDATA
  .SHARE_OF_CEDING tersimpan apa adanya, "" bila kosong) dan `dasarMataUang[{currency, tsiRnm, premiumRnm, tsiGross,
  premiumGross}]` = dasar `CedingCedantList(i).CurrencyList` pada ShareCeding = 100 (layar menghitung ShareCeding × dasar /
  100). `[terverifikasi]` `Activity\SetTSIPremiCedant_Act.xml` (pxInsName `ASM-FW-GISFW-DATA-QUOTATION!SETTSIPREMICEDANT_ACT`)
  punya dua blok berulang EMBEDDED: langkah 4 atas `pyWorkPage.OfferFacIn.CurrencyList` dan langkah 5 "Count TSI & Premi
  cedant" atas `TotalSpreadAll.pxResults` — langkah 5 ber-`pyStepsBlockName` `//` (konvensi label `//` = dimatikan, seperti
  SumTSIPremiSpreadedRNM_FIRE_Act 3.5.2.3.2.5), jadi yang berjalan **langkah 4**:
  - Share RNM (`ShareCedantType==1`): TSI = ShareCeding × (`.TSI` × PercentShare / 100) / 100, Premium = ShareCeding ×
    `.Policy.Payment.Premium` / 100;
  - Gross (`==0`): TSI = ShareCeding × `.TSI` / 100, Premium = ShareCeding × (`.Policy.Payment.Premium` × 100 / PercentShare) / 100.

  `CurrencyList(n).TSI` / `.Policy.Payment.Premium` diisi `CountPremiAndTSIRNMFireMBU_ACT` 6.2.5 (`.Name==Local.Currency`) =
  `Local.lTSI` (Σ `TSIObjectItem` item bermata uang itu, 6.2.2.1.2) / `Local.lPremi` (Σ `PremiNusantaraRe` coverage-nya,
  6.2.2.1.3.3), keduanya dinolkan per mata uang (6.2.1). Maka dasar: `tsiRnm` = Σ TSI objek × PercentShare / 100,
  `premiumRnm` = Σ premi NR, `tsiGross` = Σ TSI objek, `premiumGross` = Σ premi NR × 100 / PercentShare (kosong bila
  PercentShare 0). ⚠️ `[dugaan]` `.Policy.Payment.Premium` per mata uang dapat diubah lagi oleh activity tab Payment
  (installment / diskon) sesudah FireMBU — belum terverifikasi; di sini = nilai FireMBU. Mata uang urut kemunculan pertama;
  item tanpa mata uang tidak masuk dasar. Kosong (`[]`) bila % Share RNM kosong. Angka 8 desimal.
- **K49-9** (c3) Proteksi "Total premi share cedant must be equal total premi RNM!" (ProtectShareCedant_Act 3.3 / 4 / 5)
  **di-port** (di balik gerbang `wajib`, K49-2): `@Math.divide(TotalPremiNusaRe,1,4)` wajib sama dengan
  `@Math.divide(Σ CedingCedantList(i).CurrencyList(n).Premium,1,4)` dengan Premium menurut langkah 4 di atas; urutan pesan
  langkah 7 (paling depan, sesudah pesan agent Share Cedant Type). Akibatnya: Share RNM lolos bila Σ ShareCeding = 100; Gross
  lolos bila Σ ShareCeding = % Share RNM. `CedingCedantList(i).CurrencyList` **tidak disimpan** (nilai turunan, dihitung
  tiap baca; pola K-067) — rancangan punya `T_CEDING_CURRENCYLIST` (NAME, PREMIUM, TSI, CURRENCY_CODE; induk
  T_CEDINGCEDANTLIST) untuk data lama; migrasi 200 baru ditulis bila ada pembaca hilir (Payment / slip / FacOut).
- **K49-10** (c3, C-7 work owner 05-10-2026 "ini masa bisa lebih dari 100%", pemicu case NB 50 + 60 = 110 lolos Save) —
  PUT menolak Σ ShareCeding seluruh baris > 100 dengan 400 "Total % Share cedant tidak boleh lebih dari 100." (pesan agent,
  Pega tidak memeriksa jumlah), **tanpa gerbang**, sejajar dengan "Value can't be more than 100 or less than 0!". Dihitung
  eksak (apd); tepat 100 lolos. Di bawah gerbang `wajib` proteksi total premi (K49-9) tetap berjalan.
- **K49-11** (c3, C-8 work owner 05-10-2026 "kalau ada 2 currency, ditotalin jadi idr") — MENYIMPANG sadar dari
  `CountPremiAndTSIRNMFireMBU_ACT` (yang menjumlah angka lintas mata uang tanpa kurs): `totalTsiRnm` / `totalPremiRnm` =
  Σ per mata uang asli (TSI = PercentShare × Σ TSILiability coverage pertama tiap item / 100; Premi = Σ premi NR) × kurs ke
  IDR. Pembaca kurs = **ringkasan Spreading** (`SumTSIPremiSpreadedRNM_Act` langkah 19, `services/spreading.go`
  `ringkasan`) `[terverifikasi]` kode sendiri: ID mata uang dari `CURRENCY` (`IDMataUang`); ID `"10026"` → kurs 1 (langkah
  19.4 "SET KURS = 1, WHEN IDR"; IDR TIDAK membaca TREATYEXCHANGEYEARLY, tetapi ID-nya tetap dibaca dari CURRENCY seperti
  Spreading); selain itu `KursPada` (GetKursLimitSpreading_SQL pada Begin date) **tanpa** jatuh ke kurs terbaru — usul
  fallback `KursTerbaru` tidak dipakai karena Spreading tidak memakainya. Kunci baru `kursTidakAda: string[]` (selalu
  larik): kode mata uang tanpa kurs (item tanpa mata uang = `"(kosong)"`); bila berisi, kedua total `""` (tidak menjumlah
  sebagian). `dasarMataUang` dan CurrencyList cedant tetap dalam mata uang asli. Proteksi total premi (K49-9) tetap
  membandingkan angka **mentah tanpa kurs** di kedua sisi seperti Pega (TotalPremiNusaRe mentah vs Σ premi cedant mentah):
  premi cedant per mata uang = rasio share yang sama atas premi RNM mata uang itu, sehingga kesamaan jumlah mentah setara
  dengan kesamaan per mata uang, dan proteksi tidak bergantung pada ada-tidaknya kurs.

## Acceptance criteria

- [x] Layar: Source of Business, Share Cedant Type, % Share RNM, Total TSI / Premi RNM, grid Ceding / % Share / Add /
  Delete / Choose Ceding.
- [x] Backend GET / PUT + SetShareOfCeding + proteksi (sesi c3, 05-10-2026; migrasi 199 ditulis, belum dijalankan;
  gerbang: FLAG_ONGOING_POLICY tersimpan, IsGroup menunggu model peran, K49-2).

## Total per mata uang (work owner 06-10-2026, gambar layar Pega)

> RALAT batas tab: tab "Inw Fac Cedant Panels" di `Section\InputInwardFacultativeDtl.xml` membentang 60510–±74880 (bukan
> 60510–67588 seperti tertulis di atas). Bagian 67588–74880 terlewat saat analisis pertama.

- `[terverifikasi]` layout S62 (`pyContainerVisibleWhen IsFire`, format Notifications) @67554: grid
  `TotalSpreadAll.pxResults` (@68107) — Currency (`.Name`) · Total TSI Top Risk Spreaded (`.TSIAddCap`) · Total TSI Spreaded
  (`.TSI`) · Total Limit of Liability (`.OldTSI`) · Total Premium Spreaded (`.Premium`); varian `!IsFire` @70737 tanpa Top
  Risk; tombol Save `SaveFacIn_Act` @73883.
- Frontend: `TabCedant.tsx` memuat `GET …/spreading` dan menampilkan `RingkasanSpreading … hanyaMataUang` (data sama dengan
  ringkasan per mata uang tab Spreading). Tanpa backend baru.
