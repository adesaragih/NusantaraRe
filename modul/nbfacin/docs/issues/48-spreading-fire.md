# 48: Tab Spreading kasus FIRE

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 05-10-2026 "lanjut ke tab
> Spreading", lalu (dengan gambar layar Pega DEV) "iya sama seperti dipega. awalnya pembuatan otomatis saat di isi % Share
> RNM ( seprti dipega skrng), namun bisa juga dibuat manual, dengan klik add." Frontend sesi 0f; backend sesi c3.

**What to build:** tab Spreading: % Share RNM (menghitung TSI / premi Nusantara Re dan membuat spreading otomatis QS / SPL),
Copy Spreading manual (template → Copy To All Spreading), tampilan lokasi → item → coverage → baris spreading, total per
lokasi, ringkasan, Save.

**Status:** frontend selesai 05-10-2026 (sesi 0f); backend selesai 05-10-2026 (sesi c3, uji hijau). Migrasi 198 ditulis, menunggu dijalankan work owner — sebelum biner baru.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\` dan `DDL\`

- `Section\InputInwardFacultativeDtl.xml` tab "Spreading" (`!IsLife`): sel 80 "% Share RNM" (`.OfferFacIn.PercentShare`,
  wajib) → `CountPremiAndTSINusantaraRe_ACT(InSpreading=Spreading)`; blok "Copy Spreading" (Type Treaty dari
  `SpreadingList1`, % Share → `cekSpreadingFactIn_Act`, "Copy To All Spreading" → `CopyToAllSpreading_ACT`; Pega hanya
  menampilkannya bila `IsGroup` / user tertentu); `CoverageSpreadingList` + `SummarySpreading_Section`; tombol Total
  Accumulation dan Save (`IsThereAnyObjectLocation_Act`).
- Param `spread=Spreading` (pemicu `GetKapasitasTreaty`) hanya ada di `InputInwardFacultativeDtl_IsUW.xml` (2×), 0× di layar
  NB; work owner memutuskan pembuatan otomatis berjalan saat % Share RNM diisi.
- Rumus (`CountPremiAndTSIRNMFireMBU_ACT`, `GetKapasitasTreaty`, DT `SetTSIPremiSpreaded_FacIn`,
  `SumTSIPremiSpreadedRNM_FIRE_Act`) dan sumber treaty (`GetTreatyName_SQL`, `GetCurrencyToIDR_SQL`,
  `GetKursLimitSpreading_SQL`) dirinci di pesan kontrak ke sesi c3 (05-10-2026).
- DDL dikirim work owner 05-10-2026: `KAPASITAS_TREATY`, `TREATYEXCHANGEYEARLY`, `REINSURANCETYPE`, `TREATYCONTRACT`,
  `TREATYBUSINESS`, `PROPORTIONALARRG`.
- Label / kolom layar: `CoverageSpreadingList`, `PropertyItemListCoverageSpreading`, `InputCoverageSpreadingFire`,
  `SpreadingItem`, `InputDtlSpreadingCoverage_FacIn`, `SummarySpreading_Section` (diuji `TabSpreading.test.tsx`).
- Tidak ada di korpus: activity `CountCoverageSummary`, section `CoverageSpreading_Section`, flow action
  `SpreadingToltal_FA`, When `IsSpreadingDepan`.

## Frontend (sesi 0f)

- `components/TabSpreading.tsx` (`lebihDari100`, `templateLebih`) + uji; `pages/InwardFacultative.tsx` (tab Spreading kasus
  FIRE); `api.ts` (`ambilSpreading`, `hitungShareSpreading`, `salinSpreading`, `simpanSpreading`, `TampilanSpreading`, …);
  `labels.ts` (`SPREADING`, `TEKS_SPREADING`); `nbfacin.css` (total per lokasi kuning, ringkasan hijau seperti Pega).

## Kontrak USULAN (menunggu "kontrak spreading" sesi c3)

- `GET /api/nbfacin/kasus/{caseId}/spreading` → `TampilanSpreading` { percentShare, treaty[], template[], lokasi[ items[
  coverages[ spreading[] ] ], total[] ], ringkasanTreaty[], ringkasanMataUang[], pesan[] } — angka teks desimal.
- `POST …/spreading/hitung-share` { percentShare } → tampilan (tanpa simpan); `POST …/spreading/salin` { percentShare,
  template } → tampilan (tanpa simpan); `PUT …/spreading` { percentShare, template, lokasi } → hitung ulang + simpan.

## Keputusan agent (menunggu konfirmasi)

- **S-1** Copy Spreading tampil untuk semua pengguna (Pega: IsGroup / user tertentu) — "bisa juga dibuat manual".
- **S-2** Baris spreading per coverage baca-saja (grid editable Pega hanya untuk user ID tertentu).
- **S-3** Total Accumulation = tahap berikut.
- **S-4** % Share RNM menghitung ulang sesudah jeda 500 ms; jawaban lama dibuang.

## Backend (sesi c3, 05-10-2026)

Kontrak final = kontrak usulan di atas (dikonfirmasi sesi 9d 05-10-2026), dengan klarifikasi K48-2. Berkas:

- `backend/migrations/198_spreading.sql` (+ `_down`): `T_GENERAL_POLIS.PERCENT_SHARE`, `T_COVERAGELIST.TSI_NUSANTARA_RE` /
  `PREMI_NUSANTARA_RE`, tabel `T_SPREADINGLIST` (induk `T_COVERAGELIST`, FK tanpa `ON DELETE`) + indeks + sequence.
  **Ditulis, belum dijalankan.** ⚠️ Urutan deploy: **198 sebelum biner baru** — Save tab Object kini menghapus
  `T_SPREADINGLIST` sebelum coverage (`repository/objek.go` `sqlHapusObjek`). `docs/STRUKTUR-TABEL-NB-FACIN.md` diperbarui.
- `models/spreading.go`; `repository/spreading.go` (`PenyimpanSpreading`: pohon lokasi → item → coverage urut `SEQ_NO`,
  baris spreading, tulis ganti-utuh; master `PROPORTIONALARRG` / `REINSURANCETYPE` / `TREATYBUSINESS` / `TREATYCONTRACT`
  (`GetTreatyName_SQL`), `KAPASITAS_TREATY`, `TREATYEXCHANGEYEARLY` (`GetCurrencyToIDR_SQL`, `GetKursLimitSpreading_SQL`),
  `CURRENCY`); `services/spreading.go` (mesin + alur); `handlers/spreading.go`; pengait Save Object di `services/objek.go`.
- Rute: `GET`/`PUT /api/nbfacin/kasus/{caseId}/spreading`, `POST …/spreading/hitung-share`, `POST …/spreading/salin`.
  Galat: 400 isian (`ErrMasukanSpreading`), 401 tanpa identitas (PUT), 404 case, 409 bentuk berubah, 503 tanpa basis data.
- Port rule (berkas yang dibaca, `D:\migrasi\RNM\NB FacIn\` + `DDL\`): `Activity\CountPremiAndTSIRNMFireMBU_ACT.xml` langkah 6
  (NR), `Activity\GetKapasitasTreaty.xml` cabang IsFire (DDL 6 = NB 7), `Activity\GetTreatyName.xml`,
  `DataTransform\SetTSIPremiSpreaded_FacIn.xml` mode "percent", `Activity\SumTSIPremiSpreadedRNM_FIRE_Act.xml` cabang 3.5 / 3.6,
  `Activity\SumTSIPremiSpreadedRNM_Act.xml` langkah 18–22, `Activity\CopyToAllSpreadingFire_ACT.xml`,
  `Activity\cekSpreadingFactIn.xml` langkah 3–9, `Activity\IsThereAnyObjectLocation_Act.xml` (urutan Save).
- Uji: `services/spreading_test.go` (NR, QS saja, QS + SPL, melebihi kapasitas, tanpa kapasitas, salin, total lokasi / mata
  uang + pesan batas treaty, terorisme dijumlah dua kali, daftar treaty SF-HRE / TRT / ganda, jenis treaty, alur simpan / salin /
  hitung-share, 409, 400, 503, pertahankan saat Save Object), `handlers/spreading_test.go` (bentuk JSON, PUT dengan kunci
  tambahan, kode galat), `repository/spreading_test.go` (SQL). Nilai harapan dihitung tangan dari rumus di atas.

## Keputusan agent backend (menunggu konfirmasi)

- **K48-1** (c3) Template Copy Spreading **tidak disimpan** (Pega: `SpreadingList.pxResults` halaman sementara); GET menjawab
  `template: []`, `salin` mengembalikan template yang dikirim.
- **K48-2** (c3) Kontrak: PUT menerima `TampilanSpreading` utuh — kunci tambahan (`oldId`, `rate`, `tsi`, `total[]`,
  `treaty[]`, `ringkasan…`, kunci tak dikenal) **diterima dan diabaikan** (tanpa `DisallowUnknownFields`, permintaan sesi 9d).
  Case tanpa lokasi → 400 "Object can't be empty!" (IsThereAnyObjectLocation_Act 7 / 16). "Invalid Treaty Type Spreading!"
  hanya pada `salin` (cekSpreadingFactIn dipanggil CopyToAll, bukan Save). `totalPremiumRnm` = Σ premi NR coverage, dihitung
  tiap jawaban (tidak disimpan).
  Yang dipakai hanya `percentShare`, `oldId` coverage (pemeriksa bentuk) dan `treatyType` / `treatyName` / `sharePercentage`
  baris spreading; seluruh angka dihitung ulang server. Jumlah lokasi / item / coverage dan `oldId` per posisi wajib sama dengan
  objek tersimpan, kalau tidak 409 "Objek berubah sejak spreading dimuat - muat ulang tab Spreading". `percentShare` wajib
  0–100; Σ share per coverage / template > 100 → 400 "Share Percentage can't be more than 100" (Pega: pesan, bukan tolak);
  paling banyak 50 baris per coverage; `treatyType` wajib, lebar 50 / 500 bita. `claimAmountIdr` hanya di `ringkasanTreaty`.
  Angka keluar teks desimal 8 angka di belakang koma (hitungan skala 20, `@Math.divide`).
- **K48-3** (c3) Urutan hitung-share: NR dihitung **sebelum** spreading otomatis (Pega memanggil `GetKapasitasTreaty` sebelum
  `CountPremiAndTSIRNMFireMBU_ACT`, sehingga premi spreading memakai `PremiNusantaraRe` lama). Keanehan Pega dibawa verbatim:
  `TSIRNM100` dibagi kurs ulang di tiap coverage non-IDR; `TSIRNMTopRisk` terbawa antar-lokasi; kapasitas dicari dengan TSI
  lokasi tanpa konversi kurs (non-top risk); share SPL top risk tidak ditulis ulang sesudah dipotong `maxSPL`; coverage pertama
  bercatatan TERRORISM & SABOTAGE dijumlah dua kali ke total lokasi (langkah .4.8.2 dan .4.8.3); pesan kapasitas menghentikan
  proses (label END) dengan hasil sebagian tetap tampil. `[dugaan]` baris SPL yang ditambahkan di tengah loop tidak ikut
  diiterasi (bila ikut, Pega berputar tanpa henti saat TSI RNM > kapasitas).
- **K48-4** (c3) Pesan agent (bukan Pega): kurs `TREATYEXCHANGEYEARLY` tidak ada → spreading otomatis tidak dibuat (kurs
  seluruh item diperiksa sebelum apa pun diubah, spreading tersimpan utuh) / batas treaty mata uang itu tidak diperiksa; `KAPASITAS_TREATY` tanpa baris → pesan (Pega: Exit-Activity diam).
- **K48-5** (c3) Save: NR → DT "percent" (`TSISpreaded` = share × TSI NR / 100, `PremiumSpreaded` = share × premi NR / 100)
  → `TSIGrossSpreaded` = share × TSI × % Share / 10000 (rumus `cekSpreadingFactIn`, DT tidak mengisinya) → total → simpan.
  Pesan (batas treaty, jenis treaty) **tidak** menghalangi simpan (Pega: Page-Set-Messages). Total / ringkasan tidak
  disimpan, dihitung ulang tiap GET.
- **K48-6** (c3) Di luar cakupan: cabang EDM / Group / OR 10001, syariah, layering `CoverageBasis` 5 (di Pega pun label
  `//`), daftar ID case / OldPolicyNo hardcode (NB-80462, EDM-3704, 25 pyID langkah 18, pengecualian langkah 20),
  `PositionNote` (Pega MELEWATI pesan batas / jenis treaty untuk ReasFacInMarketing — SumTSI 20.6 / 20.7, cekSpreadingFactIn
  8.6 / 9 — dan `GetKapasitasTreaty` langkah 2 keluar untuk Marketing / Team Leader sehingga spreading otomatis tidak dibuat
  bagi mereka; di sini pesan dan spreading otomatis untuk semua pengguna sesuai arahan work owner "pembuatan otomatis saat di
  isi % Share RNM"), `QuotationData.Type == 11`,
  `FlagDelete`, pesan SF-HRE "cannot be processed" (tidak pernah menyala untuk NB), "Different (%) Share" (mati di Pega),
  Total Accumulation; langkah Save Pega sesudah total (`SetRIComm_Act`, `SaveFillPaymentInstallment`,
  `CountPaymentInstallment_Act`, `CountPaymentEdm_Act`, `CountPremiCedant_Act`, IsThereAnyObjectLocation_Act 11–15) belum
  di-port — RI Comm / installment / premi cedant milik tab lain (tiket tersendiri). Urutan mata uang ringkasan dan baris total
  per lokasi = kemunculan pertama di lokasi (Pega: `OfferFacIn.CurrencyList`, `[dugaan]` sama).
- **K48-7** (c3) Save tab Object **mempertahankan** spreading: dalam transaksi yang sama spreading lama dibaca, objek ditulis
  ulang, lalu baris spreading dipasang ke coverage di posisi (lokasi, item, coverage) yang `oldId`-nya sama dan NR + DT
  "percent" dihitung ulang dari % Share RNM tersimpan (Pega menghitung ulang di setiap Save). Coverage baru / bergeser →
  spreading kosong.
- **K48-8** (c3) `[dugaan]` Pemetaan kolom DDL `KAPASITAS_TREATY` (`LIMIT_BOTTOM_IDR`, `MAX_LIMIT_IDR`, `MAX_LIMIT_QS_IDR`,
  `MAX_LIMIT_SPL_IDR`, `TREATYNAME_QS`, `IDTREATY_QS`, `TREATYNAME_SPL`, `IDTREATY_SPL`, `STARTDATE`, `ENDDATE`) ke properti
  Obj-Browse (`.LimitBtmIDR`, `.MaxLimitIDR`, …) `belum terverifikasi` (Obj-Browse memakai nama properti, bukan kolom).
  `GetKursLimitSpreading_SQL` membandingkan teks Pega `YYYYMMDDTHHMMSS.mmm GMT` apa adanya dengan `STARTDATE` / `ENDDATE`
  seperti Pega — perilakunya di DEV `belum terverifikasi`. Bizcode treaty = `BUSINESS.ID` Class of Business + Group Business
  case (cara Table of Limit); tanpa bizcode tunggal hanya ORS + FACOUT. Batas ORS 150 M / 181,5 M dibandingkan dengan total
  sesudah kurs (mata uang batas `belum terverifikasi` — sejenis OQ-046 "nilai hardcode tanpa mata uang", belum ada OQ
  sendiri; perlu OQ baru dari register bila work owner menghendaki).

## Proteksi Copy Spreading — QS dulu (work owner 05-10-2026)

Pemicu: "di pega kn ada itu proteksi kalo spreadingnya langsung spl tidak boleh, harus ada QS nya". Bukti
`[terverifikasi]` `NB FacIn\Activity\CalcultePersentageSpeading_Act.xml` (dipanggil `change` Type Treaty grid Copy
Spreading, `Section\InputInwardFacultativeDtl.xml` ±baris 29959, `SpreadingList.pxResults`):

- langkah "utk cek harus input QS dulu": `.pxListSubscript==1` dan `!@contains(InputData.CARI15,"QS")` dan
  `InputData.CARI5!="TRUE"` (CARI5 = TRUE hanya bila `@equals(.TreatyType,"10007")||@equals(.TreatyType,"10015")`,
  langkah "If ORS and FAC-OUT") → `"Can not proceed spreading without QS"`;
- langkah "check same treaty type": `Local.Counter>=2` → `"Treaty Type can't be same"`;
- langkah "set error" (`Hasil==0`, kueri treatyGroupID) → `"jenis treaty ini tidak terdapat untuk bisnis ini "` — butuh
  kueri, belum diport.
- `[dugaan]` `.pxListSubscript==1` = baris pertama template (nesting langkah tidak terbaca pasti dari ekspor).
- **S-5 (keputusan agent, diralat 05-10-2026)**: layar hanya memeriksa Type Treaty kembar (ID, sama dengan backend) dan
  menahan Copy To All selama pesannya tampil. Syarat QS diperiksa backend saja: sesi c3 menemukan `[terverifikasi]`
  `InputData.CARI15` = `REINSURANCETYPE.NOTE` (`RDBList\GetTreatyName.xml`: `SELECT ID AS CARI1, NOTE AS CARI2 FROM
  REINSURANCETYPE WHERE ID = {InputData.CARI3}`; langkah berikut `InputData.CARI15 = .CARI2`), bukan nama dropdown.
  > RALAT: versi pertama S-5 memeriksa QS pada nama dropdown di layar — sumber keliru, dapat menahan treaty yang sah.
  Galat 400 `…/salin` tampil apa adanya.

Backend (sesi c3, 05-10-2026) — `services/spreading.go` `proteksiTemplate`, dipanggil `SalinSpreading` sesudah 404 / 503:

- `[terverifikasi]` struktur langkah: seluruh pemeriksaan = anak langkah 2 (`RH_1.pySteps(2).pySteps(n)`, kelas
  Data-SpreadingRisk, berulang EMBEDDED atas baris template); `Local.Counter` hanya dinolkan di langkah 1, ditambah di 2.13
  bila `@String.equals(Param.TreatyType,.TreatyType)`; pesan QS = 2.17 (`Local.errmsgtype`), pesan kembar = 2.14.
- `[terverifikasi]` `InputData.CARI15` = `TempTreatyName.CARI2` dari RDB-List kelas `ASM-FW-GISFW-Int-PROPORTIONALARRG`
  RequestType `GetTreatyName` (`RDBList\GetTreatyName.xml`, pxInsName `ASM-FW-GISFW-INT-PROPORTIONALARRG!ASM!GETTREATYNAME`):
  `SELECT ID AS CARI1, NOTE AS CARI2 FROM REINSURANCETYPE WHERE ID = {InputData.CARI3}`, CARI3 = `.TreatyType` baris.
  Jadi "QS" dicari di **`REINSURANCETYPE.NOTE`**, BUKAN di nama dropdown (`PROPORTIONALARRG.REINSTYPENAME`). Backend membaca
  NOTE per ID (`CatatanJenisTreaty`, `TO_CHAR(ID) = :1`).
- **K48-9** (c3) `salin` → 400 `ErrMasukanSpreading` dengan pesan verbatim, digabung "; " bila keduanya: "Treaty Type can't
  be same" (jenis treaty ber-spasi-tepi dipangkas; satu jenis muncul ≥ 2 kali — Pega memeriksa jenis yang baru diganti,
  backend memeriksa seluruh template, hasilnya sama bila setiap penggantian diperiksa) lalu "Can not proceed spreading
  without QS" (baris pertama bukan 10007 / 10015 dan NOTE-nya tidak memuat "QS", peka huruf; NOTE tidak ada = tanpa QS —
  Pega membawa CARI15 lama, `[dugaan]` kosong untuk baris pertama). Langkah 2.12 "set treaty name" (`.TreatyName` = NOTE)
  TIDAK dibawa: nama baris tetap nama dropdown. Langkah 2.3–2.9 "jenis treaty ini tidak terdapat untuk bisnis ini " **belum
  di-port**: memakai `M_TREATYBUSINESS` kolom JSON (`ac.JSONDATA.TreatyGroupID`, BizCode = `QuotationData.BusinessCode`,
  TreatyYear = `@DateTime.CurrentDate("yyyy")` — tahun berjalan, bukan tahun polis) lalu `PROPORTIONALARRG`
  `treatydescid = '10001'` (`RDBList\BrowseMaxValueTreatyType.xml`, `BrowseMaxValueTreatyType2.xml`); keberadaan
  `M_TREATYBUSINESS` (JSON) di DEV `belum terverifikasi` (DDL yang dikirim = `TREATYBUSINESS` datar).
- **K48-10** (c3) PUT **tidak** memeriksa proteksi ini per coverage: grid spreading per coverage baca-saja (S-2), baris
  coverage hanya berasal dari `GetKapasitasTreaty` (QS lalu SPL dari `KAPASITAS_TREATY`) atau dari template yang sudah lolos
  `salin`; NOTE `IDTREATY_QS` memuat "QS" `belum terverifikasi`, dan SPL tambahan di tengah loop dapat kembar, sehingga
  pemeriksaan di PUT berisiko menolak hasil otomatis. Ditinjau ulang bila grid per coverage dibuat dapat diubah.

## Acceptance criteria

- [x] Layar sesuai gambar Pega: % Share RNM, Copy Spreading, lokasi → item → coverage → spreading, total, ringkasan.
- [x] Backend: hitung share + spreading otomatis, Copy To All, simpan (sesi c3, 05-10-2026; migrasi 198 ditulis, belum dijalankan).
