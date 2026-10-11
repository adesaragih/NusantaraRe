# Penampung halaman `TreatyIn` — cabang Proporsional

Ditulis 7 Oktober 2026, atas permintaan pemilik proses: *"membuat penampung data seperti page list
yang ada di Pega … pastikan saat saya menginput data saya pindah ke tab lain datanya tidak hilang,
jangan gunakan table dulu … dan pastikan juga contoh tiap tab yang saling berelasi yang bergantung
nilai dari tab lain itu berfungsi juga"*.

## 1 · Masalahnya

Di Pega seluruh tab satu kasus menulis ke SATU halaman clipboard `TreatyIn` — properti skalar
(`ReportingStart`, `RNMShareP`, …) dan page list (`Portfolio`, `Limits`, `AccumulationList`, …).
Berpindah tab tidak membuang apa pun, dan tab yang satu membaca isian tab lain dari halaman yang sama.

Layar ini menyimpan isian di `useState` TIAP TAB. Tab yang tidak tampil dilepas React, sehingga
isiannya hilang, dan tab lain tidak pernah melihatnya. Grid `Kind of Treaty` tab Share, misalnya,
membaca proyeksi lain (`M_TREATY_IN2` per layer), bukan `Limits` yang sedang diisi.

## 2 · Bentuknya — `frontend/halaman.tsx`

| Bagian | Isi |
|---|---|
| `usePenampungHalaman()` | Dipegang `FormKontrakTreatyIn`, DI ATAS strip tab. Dikosongkan saat kontrak lain dimuat dan saat data kontrak tiba |
| `PenyediaHalaman` | Membungkus isi tab |
| `useProperti(kunci, awal)` | Pengganti `useState` di tab. `kunci` = ejaan properti Pega. Properti yang belum ada DISEMAI sekali dari `awal` (data kontrak yang dimuat), sehingga semua tab yang membacanya mendapat nilai yang sama. Tanpa penyedia (uji satu tab) → `useState` biasa |

⛔ **Bukan tabel.** Nol tulisan ke basis data. Isinya masuk tabel saat Save/Submit, dan itu pekerjaan
berikutnya. Kunci dan isi barisnya **ejaan Pega** supaya pemetaan Save nanti satu lawan satu.
Tanggal di dalam page list disimpan dalam bentuk simpan `YYYYMMDD`.

## 3 · Properti per tab — dari `Section/TreatyInTabsProportional.xml`

| Tab | Properti di penampung |
|---|---|
| Reporting Period | `ReportingStart` `ReportingEnd` `ReportingPeriod` `ReportingInterval` `ReportingSubmission` `ReportingConfirmation` `ReportingSettlement` · `ReportingPeriodList[]{Period, AutoCalculate, InitialDate, SubmissionDue, ConfirmationDue, SettlementDue}` |
| Portfolio | `Portfolio[]{TypePortfolio, Type, Description}` |
| Limits | `Limits[]` (pohon `Detail[]` utuh, ejaan dokumen) · halaman sesi `SearchData.CARI1` (As At Quarter) `SearchData.CARI2` (Quarter Year) `TempQuarter` `TempQuarterYear` `FlagExcel.CARI1` |
| Share | `RNMShareP` `BrokeragePercentP` `OptionLimit` · **`Limits[]`** · `TotalShareRnmProp` `TotalSpreadedRnmProp` `TotalSpreadedRnmRIProp` |
| Co-Ins Scale | `CoInScale[]{CoInShare, PctLimit}` `MaxCoNonGroup` `MaxCoGroup` |
| Accumulation | `AccumulationPeriod` `AccumulationList[]{Period, ReportDate, SubDays, SubDueDate}` · **`ReportingStart` `ReportingEnd`** (dibaca) |
| Exclusions / Special Conditions | `ExclusionsP` / `SpecialConditionsP` (Non-Prop: `Exclusions` / `SpecialConditions`) |
| Information & Submit | `Information` `Comment` |
| Achievement In IDR | **`Limits[]`** (dibaca) |
| Retro | — keputusan §17, tidak dibangun |

Properti bercetak tebal adalah jalur ketergantungan antartab. Daftarnya juga ada di kode
(`PROPERTI_TAB_PROP`) dan dijaga `halaman.test.tsx`.

## 4 · Ketergantungan antartab yang kini berfungsi

| Dari → ke | Activity ekspor | Di layar |
|---|---|---|
| Limits → Share | Grid `Kind of Treaty` (`TreatyInShareProp.xml`) ADALAH `TreatyIn.Limits`; rinciannya `TotalLimits` → `DetailShare` adalah `.Detail[]` yang sama | Layer/Treaty Group yang diisi di tab Limits langsung tampil di Share |
| Limits ↔ Share | `TreatyInPropshare` (Refresh, % RNM Share, Option) membaca `CessionList`/`IOOLimitList` tiap Detail dan MENULIS `RNMShare`, `ShareNote`, `RNMShareList`, spreading, dan tiga total | Rute baru `POST /hitung/share-prop`; hasilnya ditulis balik ke `Limits` di penampung, sehingga tab Limits dan Achievement melihatnya |
| Limits → Achievement In IDR | `AchievementCombine` atas `TreatyIn.Limits` | `TabAchievement` membaca `Limits` dari penampung |
| Reporting Period → Accumulation | `TreatyInSetAccountReport` membaca `TreatyIn.ReportingStart`/`ReportingEnd` | Rute baru `POST /hitung/akumulasi`; Start/End dibaca dari penampung walau tab Reporting Period sedang tidak tampil |
| Exclusions / Special Conditions (Prop) | DT `TreatyInCopyConditions` (perilaku `change` `ExclusionsP`/`SpecialConditionsP`): `Exclusions = ExclusionsP`, `SpecialConditions = SpecialConditionsP` | Disalin saat isian ditinggalkan; sumber yang belum pernah dibuka tidak disalin |
| Maximum Retention → EGNPI (Non-Prop) | `TreatyInNonAddItem(egnpi)` langkah 3: mata uang baris baru = `Retention(1)` | Form membaca `Retention` TERKINI dari penampung |
| EGNPI → Limits (Non-Prop) | `TotalEgnpi` (`EgnpiTotalList` tiap layer) | Form membaca `EGNPI` TERKINI dari penampung |
| Share → Installment (Non-Prop) | `TreatyInSetValueInstallment` membaca `TreatyIn.TotalShareNetNP` | Form membaca total Share TERKINI |

## 5 · Rumus baru (services, diuji)

- **`hitung_akumulasi.go`**: `TreatyInSetAccountReport` dan `TreatyInAccumulationSetSubDue`. Memakai
  pembantu tanggal yang sama dengan `TreatyInSetReport`, yang terukur 99,5%. Add =
  DataTransform `TreatyInAddAccumulation` (di layar).
- **`hitung_share_prop.go`**: `TreatyInPropshare`, `TreatyInPropshareDetail`, `CalculateShareList`,
  `FetchQSfromMaster`, `SetSpreadName`, dan DT `CountTotalPctSpead`.
  - RD anak `BrowseTreatyArrangement_Limit_MstTrt_RD` dibangun dengan kelima filternya
    (`BacaAnakSpreadingProp`; parameter kosong dilewati).
  - **Terukur di master** `PROPORTIONALARRG`: induk `10241` = "2024 QS 155M TRT", yaitu Spreading Type
    gambar Pega 17. Anaknya dibagi per Treaty Group (22 × 2). Filter grup pada jalur Refresh
    menghasilkan tepat QS (OR) 40 / QS (R/I) 60, dan uji mereproduksi angka gambar 17:
    3.000.000.000 · 1.200.000.000 · 1.800.000.000.

## 6 · Disalin apa adanya walau janggal (rumus, bukan alamat)

- `SetSpreadName` langkah 5 menjumlah total seluruh Detail setiap kali dipanggil, yaitu sekali per
  Detail TANPA Spreading Type. Lalu `TreatyInPropshare` langkah 6 menjumlah lagi, sehingga total
  cabang itu berlipat. Gambar 17 (cabang ber-Spreading Type) tidak berlipat. Cabang tanpa Spreading
  Type belum punya gambar pembanding.
- `SetSpreadName` langkah 6 memasang "Total share must equal RNM share.!!" bila total spreading
  ≠ % RNM Share, termasuk pada Detail tanpa baris spreading.
- `TreatyInSetAccountReport`: `none` keluar tanpa mengubah daftar; `other` mengosongkan daftar;
  `TempDate` berantai (31 Jan → 28 Feb → 28 Mar).

## 7 · Tidak berbukti — diputuskan dan dinyatakan

- Start/End kosong: selisih 0 bulan, sehingga satu baris tanpa tanggal.
- `Page-Clear-Messages` dianggap membuang semua pesan sebelumnya.
- Label `AccumulationPeriod`: hanya "None" yang terlihat (gambar 19); sisanya tampil apa adanya
  (konvensi `prompt_value.go`).
- Tanggal Accumulation disimpan `YYYYMMDD` (tanggal WIB). Pega menyimpan stempel DateTime GMT;
  bentuk simpannya urusan Save.

## 8 · Belum

- Retro (§17): `Share to Other Retro`, `Facultative Reinsurers`, `AddFacRetroProp`, `CountRetroShare_Act`.
- Tambah/hapus baris spreading manual di `DetailShare`: tombolnya tidak terbaca di ekspor
  (sel `.pyTemplateInputBox`).
- Limits/Share Non-Prop masih memakai keadaan form `limitsNP`/`shareNP` yang sudah ada (keadaan
  FORM, jadi tetap bertahan saat pindah tab).
- ~~`TreatyInUpdatePaymentDate` / `_Act`~~ — DIUNGGAH pemakai 7 Oktober 2026 dan dibangun (§10).

## 9 · Ronde 7 Oktober 2026 (kedua) — "input tiba-tiba hilang" dan cabang Non-Prop

**Dua penyebab isian hilang, keduanya dicabut:**

1. **Rumus asinkron menimpa simpul utuh.** Tab Limits Prop memicu `LimitCalculation`,
   `CalculateDeduction`, dan `PremiumReserveCalculate` saat isian ditinggalkan. Jawabannya dulu
   ditimpakan sebagai Detail UTUH yang disalin saat rumus dipicu, sehingga isian yang diketik selama
   rute menjawab kembali ke nilai lamanya. Kini jawaban diterapkan sebagai fungsi atas Detail TERKINI
   (`terapkan`), dan hanya medan yang Activity tulis yang diganti.
2. **Efek reset tiap render.** `TabRetensi`/`TabEgnpi` memasang
   `useEffect(() => setRows([...baris]), [baris])`, padahal `baris` dibuat ulang setiap kali form
   dirender. Apa pun yang berubah di form mengembalikan baris ke data kontrak. Efek itu dicabut.
   Kontrak lain = penampung dikosongkan dan tab dilahirkan ulang.

**Non-Prop di penampung:** `Retention`, `TotalRetentionAmountNP`; `EGNPI`, `TotalEgnpiAmount`,
`TotalEgnpiProportion`, `TotalEgnpiAmountNP`; kedelapan properti Event Limits (`CurrencyRSMD`,
`RSMDLimit`, `CurrencyEarthquake`, `Earthquake`, `CurrencyFloodJab`, `FloodJab`,
`CurrencyFloodNat`, `FloodNation`); tab Installment baru (`TabAngsuran`: `InstallmentNo`,
`Installment`, `TotalInstallmentNP`, rute `POST /hitung/angsuran`).

**Alur XML yang tadinya terlewat di Treaty In:**

| Medan | Ekspor | Sebelumnya | Kini |
|---|---|---|---|
| Accumulation · Period | Perilaku `change` = DT pra-refresh `TreatyInDeleteAccumulationLists`, LALU `TreatyInSetAccountReport` (kedua korpus) | Daftar lama ikut terkirim; `none` (langkah 2: keluar) membiarkannya | Daftar dikirim KOSONG; `none` meninggalkan daftar kosong |
| Reporting Period · sel Initial Date | `TreatyInSetReport(startdate = .InitialDate, autocalculate = .AutoCalculate)` | Tidak ada pemicu | Bila Auto Calculate baris dicentang, daftar disusun ulang mulai tanggal itu. Rute menerima `awal` (`param.startdate`); durasi tetap Start→End (langkah 8) |
| Exclusions / Special Conditions (Prop) | DT `TreatyInCopyConditions` | Tidak disalin | Disalin saat isian ditinggalkan |

Rute `share-np` mendapat aksi `set-brokerage` (`TreatyInSetBrokerage` SAJA). Layar Adjustment
memanggil satu aksi per langkah bersyarat kotak centang Share Across The Board; `brokerage` (SetBrokerage
lalu XOLAddSpreading) tetap untuk isian `% Brokerage` Treaty In. Uji: `set-brokerage` disusul `rnm` =
`brokerage`.

## 10 · Ronde 7 Oktober 2026 (ketiga)

- **Payment Date** (tab Installment): sel Due Date → `TreatyInUpdatePaymentDate` (halaman itu),
  sel WPC → `_Act` (semua halaman). `PaymentDate = DueDate + WPC + 1 hari`. Rute `angsuran`
  `tanggal-bayar` / `tanggal-bayar-semua`; hanya `PaymentDate` yang diganti atas halaman TERKINI.
- **Retro disembunyikan** dari strip tab (keputusan pemilik proses) — `TAB_DISEMBUNYIKAN`;
  `tabUntuk` tetap menilai syaratnya.
- Rute baru untuk layar Adjustment (dipakai juga bila kelak dibutuhkan di sini): `selisih`
  (Value Difference) dan `aktual` (cabang Adjust Premium).
