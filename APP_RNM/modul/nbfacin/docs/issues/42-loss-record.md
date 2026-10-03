# 42: Sub-tab Loss Record dan Loss Record Internal

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026 (urutan
> sub-tab objek, lihat tiket 39). Belum ada tangkapan layar; tampilan diturunkan dari XML.

**What to build:** dua sub-tab terakhir baris objek FIRE.
- **Loss Record:** grid catatan kerugian (Date of Loss · Insured Name · Loss Object · Currency · Total of Loss · Total
  Claim) dengan Tambah / Hapus dan form detail, plus empat angka loss ratio baca-saja.
- **Loss Record Internal:** grid klaim internal baca-saja.

Data Loss Record ikut Save tab Object.

**Blocked by:** daftar Remarks (aturan properti tidak ada di korpus); rumus Loss Ratio (`SetLossRatio_Act` versi baris
objek tidak ada); keputusan work owner atas pengisian Loss Record Internal.

**Status:** frontend selesai 03-10-2026 (uji hijau); backend → sesi c3 (sesudah tiket 41).

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\CauseOfLoss_FacIn.xml`: grid `.Property.ListCauseOfLoss` (kelas `Data-CauseOfLoss`, master-detail, flow
  action `InputCauseOfLoss_FacIn`); kepala sel 17–22; Add / Delete = ikon (+ `SetLossRatio_Act`).
- `Section\InputCauseOfLoss_FacIn.xml`: sel 3 Insured Name (`QuotationData.InsuredName`, tampil) · 4 Date of Loss · 5
  Loss Object · 6 Currency (RD `BrowseCurrency_RD`) · 7 Total Claim (100%) (`.Claim`, 4 desimal, awal 0) · 8 Prevention
  Of Loss · 9 Cause of Loss · 10 `.Remarks` (dropdown associated; label dari deskripsi properti) · 11 Loss Detail
  (wajib). `.Amount` (Total of Loss) tanpa isian.
- `Section\InputOfferFacInLossRatio.xml`: LR 1 Year (3 desimal) · %LR 1 Years (2) · LR 3 - 5 Years (3) · %LR 3 - 5 Years
  (2), baca-saja, atas `.LossRatio1Year*` / `.LossRatio35Year*` baris objek.
- `SetLossRatio_Act` di korpus hanya versi `Data-ScoringRisk`; versi baris objek yang dipanggil grid TIDAK ada → rumus
  LR tidak diketahui. Data contoh: semua LR = 0.
- `Section\CauseOfLossClaim_FacIn.xml`: tombol "Get Loss Record Internal" tersembunyi (`1=0`); grid baca-saja
  `.Property.ListCauseOfLossClaim`, sel 24–34.
- `Activity\MappingKlaimToLossRecord_Act.xml`: mengisi grid dari `RDBList\GetDataKlaim_SQL.xml` (POOLDATA.DATAKLAIM)
  dengan NOPOLIS ter-hardcode; semua pemanggilnya dikomentari / tersembunyi → di Pega grid ini tidak pernah terisi
  (data contoh: kosong).

## Keputusan agent

- **N-1** Rumus Loss Ratio tidak diport; nilai LR tampil apa adanya dari server.
- **N-2** Total of Loss (`.Amount`) baca-saja.
- **N-3** Insured Name grid (`.CoinsData.CoinsName`) diisi nama tertanggung case saat Add `[dugaan]` — pengisi Pega
  tidak ada di korpus; data contoh selalu terisi.
- **N-4** Loss Detail bertanda wajib tanpa menahan Save (pola L-2).
- **N-5** Loss Record Internal = grid baca-saja dari server; pengisian otomatis menunggu keputusan work owner.
- **N-6** Total Claim awal "0".
- **N-7** Uang (Total of Loss, Total Claim, Prevention Of Loss, Premium, klaim) = teks desimal; tampil 4 desimal lewat
  `formatNumber`; isian harus desimal bertitik.
- **N-8** Label `.Remarks` = "Remarks" (sistem baru).

## Kontrak

- `ObjekFire.lossRecords: CatatanKerugian[]` (`dateOfLoss` DD-MM-YYYY, `coinsName`, `lossObject`, `currency`,
  `amount`, `claim`, `preventionOfLoss`, `causeOfLoss`, `remarks`, `detail`) — ikut `GET`/`PUT …/objek`.
- `ObjekFire.lossRatio` (`oneYearAmount`, `oneYearPercent`, `threeFiveYearAmount`, `threeFiveYearPercent`) dan
  `ObjekFire.internalLossRecords: KlaimInternal[]` — BACA-SAJA; dikirim `GET`, diabaikan `PUT`.

## Acceptance criteria

- [x] Loss Record: grid + Tambah / Hapus + form + loss ratio; label diuji ke korpus.
- [x] Loss Record Internal: grid baca-saja; label diuji.
- [ ] Backend: simpan / baca ListCauseOfLoss; kirim loss ratio dan klaim internal.
- [ ] Daftar Remarks; rumus LR; keputusan Loss Record Internal.
