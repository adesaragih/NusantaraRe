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

**Status:** frontend selesai 03-10-2026 (uji hijau); backend selesai 03-10-2026 (sesi c3; migrasi 191 ditulis, BELUM dijalankan).

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
- [x] Backend: simpan / baca ListCauseOfLoss (migrasi 191); loss ratio dihitung server (W-4); klaim internal selalu [] (N-5).
- [ ] Daftar Remarks; rumus LR; keputusan Loss Record Internal.

## Keputusan work owner — 03-10-2026 (diteruskan sesi `nusantarare-0f`)

- **W-4** (butir 84): Loss Ratio dihitung ulang server menurut `D:\migrasi\RNM\DDL\SetLossRatio_Act.xml` (kelas
  LocationReinsurance); bila ΣAmount = 0, LR dan %LR tetap 0. Total of Loss tetap baca-saja seperti Pega.
- **Loss Record Internal** (butir 83): *"biarkan saja kosong"* — N-5 dikonfirmasi; GET selalu `[]`, tanpa tabel, tanpa
  pengisian dari DATAKLAIM.
- Keputusan agent sesi 0f N-1…N-8 disetujui (butir 82). ⚠️ N-1 ("rumus LR tidak diport") **digantikan** W-4.

## Backend (sesi c3, 03-10-2026) — disusun agent

- `GET`/`PUT …/objek`: `lossRecords` (kunci persis `CatatanKerugian`) → tabel rancangan `T_LISTCAUSEOFLOSS` (+ lima kolom
  baru) dan `T_COINSDATA` (satu per catatan), urut `SEQ_NO`. `lossRatio` dan `internalLossRecords` BACA-SAJA (badan PUT
  diabaikan); `internalLossRecords` selalu `[]`.
- **Saat PUT, per objek** (`services/kerugian.go`, `SetLossRatio_Act` `[terverifikasi]`):
  - `coinsName` = nama tertanggung case (langkah 3.1, semua catatan);
  - `dateOfLoss` DD-MM-YYYY → teks Pega `YYYYMMDDTHHMMSS.mmm GMT` pukul 12:00 WIB (pola Begin date; kontrol pxDateTime);
  - Loss Ratio: tanggal acuan = hari ini Asia/Jakarta; ≤ 365 hari → Σ1, ≤ 1825 hari → Σ2 (kumulatif); hanya bila
    ΣClaim ≠ 0: LR = ΣClaim/ΣAmount, %LR = ΣClaim·100/ΣAmount; ΣAmount = 0 → 0 (W-4); disimpan ke `T_LOCATIONLIST.LOSS_RATIO*`.
- **Validasi PUT** (400 `baris[n].lossRecords[m].<medan>`): `currency` wajib dan ada di `CURRENCY` (≠ ITL);
  `amount`/`claim`/`preventionOfLoss` kosong atau desimal ≤ 8; `dateOfLoss` kosong atau DD-MM-YYYY sah; lebar kolom; `detail`
  tidak wajib (N-4); `remarks` tanpa enumerasi.
- Migrasi **191** (`191_t_listcauseofloss.sql` + `_down`); loader `amandemenKerugian` + pelebaran `DETAIL` (skema 80 tabel /
  **1.438** kolom). ⛔ Ditulis, **tidak dijalankan** agent.

**Bukti `[terverifikasi]`:** `SetLossRatio_Act.xml` — enam syarat (`pyStepsPreCondParamsWhen`: IsB2B, ≤ 365, ≤ 1825, CARI8,
Claim1, Claim2) dan delapan penanda `pyStepsPreCondition` (5 `true`, 1 `false`, 2 kosong); dihitung dua cara (baris dan
kemunculan tag). ⚠️ Ralat: tulisan awal "enam penanda untuk tujuh syarat" salah hitung. Urutan elemen ekspor diacak, jadi
penanda tidak dapat dipasangkan pasti ke langkahnya (jebakan sensus 6) — "langkah CARI8 mati" bersandar pada uraian sesi
0f `[dugaan]`, tidak terbantah korpus. Kontrol
`Section\InputCauseOfLoss_FacIn.xml`: DateOfLoss pxDateTime; LossObject / CauseOfLoss pxTextInput; PreventionOfLoss pxNumber 4
desimal; Detail pxTextArea; Remarks pxDropdown. `.Amount` (Total of Loss) TIDAK ada di form itu — hanya di grid
`Section\CauseOfLoss_FacIn.xml` (2 kemunculan). Fixture: 6 catatan, hanya Claim / Currency / Detail / Remarks / CoinsData
berkunci (kolom rancangan diturunkan dari itu).

**Keputusan agent — DISETUJUI work owner 03-10-2026** (butir 90, diteruskan sesi `nusantarare-0f`: *"setuju sesuai
rekomendasi agent A145–A152"*):

| # | Keputusan | Dasar |
| --- | --- | --- |
| A145 | Lima kolom baru `T_LISTCAUSEOFLOSS`: `DATE_OF_LOSS` VARCHAR2(30), `LOSS_OBJECT`/`CAUSE_OF_LOSS` VARCHAR2(500), `AMOUNT`/`PREVENTION_OF_LOSS` NUMBER(38,8) | medan ada di layar, tidak di rancangan (contoh tidak memuatnya); pola A110 / START_DATE_TIME / V-6 / ADR-0016 |
| A146 | `DETAIL` dilebarkan 50 → 500 | pxTextArea "Loss Detail"; lebar rancangan diturunkan dari contoh kosong; pola V-6 catatan |
| A147 | Catatan tanpa `dateOfLoss` tidak ikut dijumlah LR | permintaan sesi 0f; selisih tanggal tidak terdefinisi |
| A148 | LR dan %LR dibulatkan setengah-ke-atas pada desimal ke-8; %LR = ΣClaim·100/ΣAmount (bukan (ΣClaim/ΣAmount)·100 yang sudah dibulatkan) | ADR-0016 (8 desimal); usul sesi 0f; tanpa float |
| A149 | `dateOfLoss` disimpan teks Pega DateTime 12:00 WIB; dibaca balik juga bentuk Date 8 digit | kontrol pxDateTime; pola Begin date butir 78.1 |
| A150 | `amount` (Total of Loss, baca-saja di layar, N-2) **diterima apa adanya** dari badan PUT (diperiksa desimal) dan disimpan | spesifikasi membolehkan "terima apa adanya atau abaikan"; mengabaikan akan menghapus nilai lama karena baris diganti utuh |
| A151 | `amount` / `claim` / `preventionOfLoss` negatif → 400 | pola uang ≥ 0 tiket 39 (`tsi`, K-3); Pega `SetErrorMessageTSIObjectItem_Act` menolak minus untuk TSI — untuk kerugian `belum terverifikasi` |
| A152 | Hasil LR yang tidak muat NUMBER(38,8) (ΣClaim jauh melebihi ΣAmount) → 400 `baris[n].lossRatio` | mencegah galat Oracle 500 saat simpan; temuan code review |

⚠️ **Tidak diport:** langkah 1 `SetLossRatio_Act` (`pyWorkPage.IsB2B=="ASM"` → Detail catatan terakhir "No Info").
⚠️ Catatan bertanggal sesudah hari ini (selisih negatif) ikut dijumlah — sama dengan perbandingan `<=` Pega.

**Aturan properti (03-10-2026, diteruskan sesi 0f; butir 85):** `DDL\Remarks.xml` (ASM-FW-GISFW-DATA-CAUSEOFLOSS!REMARKS,
PromptList) = Settled / Ex Gratia Payment / Withdraw / Others / -- (nilai = label) dipakai frontend; terpanjang 17 bita, muat
di VARCHAR2(500). Backend tetap tanpa validasi enumerasi.
