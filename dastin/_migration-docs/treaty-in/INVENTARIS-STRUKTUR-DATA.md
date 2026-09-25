# Inventaris struktur data — Treaty In

**Untuk:** sesi perancangan model data / ERD berikutnya
**Tanggal:** 23 September 2026

Ini **bahan**, bukan rancangan. ERD bukan pekerjaan sesi ini. Yang ada di sini adalah daftar benda,
atributnya, dan cabang mana yang memakainya.

---

## 0. Cara inventaris ini disusun, dan batasnya

**Sumbernya `JSONDATA`, bukan tabel relasional.** `M_TREATY_IN` menyimpan seluruh halaman clipboard
sebagai satu dokumen JSON; `TREATY_IN` hanya punya 20 kolom dan `TREATYINDETAIL` adalah proyeksi
datar dari sebagian isinya. Menyusun struktur dari tabel relasional akan kehilangan sebagian besar
model, dan yang lebih berbahaya — tidak akan tahu apa yang hilang.

Karena isi `JSONDATA` produksi belum tersedia, strukturnya direkonstruksi dari **ekspor aturan
Pega**: setiap jalur properti yang dibaca atau ditulis aturan mana pun.

**Angka hasil penyapuan:**

| | Jumlah |
|---|---|
| Jalur properti unik berakar `TreatyIn.` | 864 |
| Akar (properti tingkat pertama) | 199 |
| Jalur yang **pernah ditulis** oleh suatu aturan | 580 |

### Batas yang harus diketahui pembaca

1. **Ini batas bawah, bukan daftar lengkap.** Penyapuan hanya menangkap jalur yang ditulis dengan
   awalan literal `TreatyIn.`. Aturan yang bekerja pada *step page* menulis `.SpreadingList`,
   `.DeductionList`, `.RNMShareList` tanpa awalan itu, sehingga sebagian daftar anak tidak muncul
   di sini meski keberadaannya sudah dipastikan selama pembedahan. Yang diketahui hilang minimal:
   `SpreadingList`, `BreakDownSprdList`, `RNMShareList`, `RNMSpreadedList`, `RNMSpreadedListRI`,
   `DeductionList` pada tingkat layer, `Reinstatement_List` selain dua kolom jumlah, dan
   `MDPMinList`.
2. **Verifikasi terhadap `JSONDATA` produksi wajib dilakukan** sebelum model dikunci. Cara
   termurahnya sudah ada di berkas uji: kueri `N-0`, `P-0`, dan `R-1` masing-masing mencetak
   cuplikan `JSONDATA` mentah.
3. Tabel relasional **tetap berguna**, tetapi untuk hal lain: ia memberi tahu field mana yang dulu
   dianggap layak dilaporkan atau dibagikan ke sistem lain. Itu masukan tentang **prioritas**,
   bukan tentang **kelengkapan**.

---

## 1. Temuan struktural yang paling menentukan: pohon yang sama disimpan sampai empat kali

Halaman kontrak membawa **salinan paralel** dari bentuk yang sama:

| Cabang | Jalur anak yang ditulis | Yang bentuknya identik dengan pohon utama |
|---|---|---|
| `ValueDifference` | 161 | **90** |
| `ActualValue` | 107 | **77** |
| `OLDDATA` | disalin utuh sebagai halaman | — |
| `ValueBeforeProrate` | disalin utuh dari `ValueDifference` | — |

Artinya satu kontrak menyimpan struktur limit, share, potongan, premi, penyebaran, dan installment
**berkali-kali** dalam satu dokumen: sekali sebagai nilai berjalan, sekali sebagai nilai aktual,
sekali sebagai selisih, sekali sebagai nilai sebelum prorata, dan sekali lagi sebagai data lama.

**Implikasi untuk model baru:** ini bukan lima jenis benda. Ini **satu bentuk** yang muncul dalam
beberapa peran. Menyalin kelima cabang menjadi lima kelompok tabel akan mengabadikan duplikasi yang
seharusnya menjadi *versi* dan *turunan*. ADR-0036 dan ADR-0048 sudah menetapkan arah penyelesaian:
versi yang dibekukan, dan selisih sebagai turunan yang dibukukan dengan penunjuk ke versinya.

---

## 2. Kepala kontrak — dipakai kedua cabang

| Kelompok | Atribut |
|---|---|
| Identitas | `ID`, `OLDID`, `TreatyContractName`, `TreatyYear` |
| Pihak | `Ceding` / `CedingID`, `LeadingReinsSource` / `LeadingReinsSourceID`, `LeadingReinsID` |
| Lingkup | `ProportionType`, `TeritorialScope`, `ClassOfBusiness`, `Information` |
| Periode | `Commencement`, `Termination` |
| Bagian dan biaya | `RNMShareP`, `RNMShareAcrossTheBoard`, `BrokeragePercent`, `FacultativeShare`, `FacultativeShareBrokerage`, `IsMultipleRetro` |
| Siklus hidup (**dilebur**, ADR-0046) | `StatusAkseptasi`, `ChooseStatusAkseptasi`, `Position`, `PositionUsername`, `ViewState`, `IsEditData`, `RevisionState` |
| Addendum | `EDMState`, `EDMMaterialType`, `EDMEffective`, `IsProRate`, `ProRateDays`, `ProRateTotalDays`, `ProRatePercent` |

### Daftar milik kepala

| Daftar | Atribut | Catatan |
|---|---|---|
| `CommentList[]` | `Date`, `OperatorName`, `IsApproved`, `Suggest` | **narasi manusia**; jangan digabung dengan jejak perubahan (ADR-0045) |
| `CurrencyList[]` | `CurrencyID`, `Currency`, `Conversion`, `PeriodStart`, `PeriodEnd` | **bentuknya sudah benar** — kurs beku per periode di dalam kontrak. **KOREKSI 23 Sep 2026:** klaim sebelumnya — "di aplikasi berjalan tidak ada yang mengisinya" — **salah**. Penulisnya dua, dan keduanya berlangkah HIDUP: `TreatyInMappingDataconvert` langkah 5 (konversi satu kali) **dan** `TreatyInAddCurrency` langkah 1, yang duduk di balik kontrol bersyarat `TreatyIn.ViewState != '1'` — terjangkau setiap pengguna selama kontrak dapat disunting. Klaim lama lahir dari menghitung rujukan, bukan langkah hidup (CONTEXT.md §2.0a). Arah koreksinya **menguatkan** ADR-0053: daftar mata uang memang diisi dan dipelihara pengguna |
| `Retention[]` | `TreatyGroupID`, `TreatyGroup`, `CurrencyID`, `Currency`, `Amount`, `ID` | retensi sebagai **jumlah uang** per treaty group — ini **masukan** untuk Surplus (ADR-0037) |
| `EGNPI[]` | `TreatyGroupID`, `TreatyGroup`, `CurrencyID`, `Currency`, `Amount`, `AmountIDR`, `AsDate`, `ID` | dasar perhitungan premi XOL |
| `Portfolio[]` | `Type`, `TypePortfolio`, `Description` | |
| `ReportingPeriodList[]` | `Period`, `InitialDate`, `SubmissionDue`, `ConfirmationDue`, `SettlementDue` | tiga jatuh tempo bertingkat; **tidak ada apa pun yang mengawasinya** |
| `AccumulationList[]` | `Period`, `ReportDate` | |
| `Installment[]` | `Currency`, `ID`, + `InstallmentList[]` (`InstallmentPct`, `Installment`, `DueDate`, `PaymentDate`, `WPC`, `Currency`) | |

---

## 3. Cabang NON-PROPORSIONAL (XOL)

Pembagian prop/non-prop di bawah **tidak ditebak dari nama berkas** — ia dibaca dari struktur:
`Limits[]` memuat layer, deductible, dan reinstatement; `Limits[].Detail[]` memuat quota share,
surplus, dan retensi.

### 3.1 `Limits[]` — layer

| Kelompok | Atribut |
|---|---|
| Identitas layer | `ID`, `Layer`, `LayerPart`, `LayerType`, `LayerPartType`, `Cover`, `TreatyType` |
| Batas | `Limit`, `Limit2`, `Deductible`, `Deductible2`, `Currency`, `Currency2`, `CurrencyRelation`, `AggregateLimit`, `AggregateLimit2` |
| Premi | `AdjRate`, `MDPPct`, `MDPMinPct`, `ROLPct` |
| Reinstatement | `ReinstatementValue`, `ReinstatementPct`, `ReinstatementNote` |

Daftar anak:

| Daftar | Atribut |
|---|---|
| `EgnpiTotalList[]` | `CurrencyID`, `Currency`, `Value` |
| `PremiumEarnedList[]` | `Currency`, `Value` |
| `MDPList[]` | `Currency`, `Value` |
| `Reinstatement_List[]` | `ReinstatementValue`, `ReinstatementPct`, `AdditionalPct`, `ReinstatementAmount1/2`, `AdditionalAmount1/2`, `ReinstatementNote`, `ID` |
| `TreatyGroupList[]` | `TreatyGroupID`, `TreatyGroup`, + `ClassOfBusinessList[]` (`ClassOfBusinessID`, `ClassOfBusiness`) |

**Catatan model yang mengikat:**

- `Currency`/`Currency2` dan `Limit`/`Limit2` adalah **kolom kembar yang tidak dibawa**. Yang
  bertahan: mata uang sebagai atribut bernilai, `CurrencyRelation` sebagai atribut bernama, dan
  "paling banyak dua mata uang" sebagai **aturan validasi**, bukan sebagai bentuk data.
- `MDPPct` (deposit) dan `MDPMinPct` (minimum) adalah **dua besaran berbeda** dan tidak wajib sama.
- `ReinstatementPct` dan `AdditionalPct` adalah **dua masukan yang dinegosiasikan terpisah**;
  `ReinstatementAmount*` dan `AdditionalAmount*` adalah turunannya. Penamaan di sistem lama
  tertukar — lihat `RINGKASAN-GRILLING.md` area E.

### 3.2 `Share[]` — bagian NuRe atas layer

| Kelompok | Atribut |
|---|---|
| Identitas | `Layer`, `LayerPart`, `LayerType`, `LayerPartType`, `Cover`, `ClassofBusinessList`, `TreatyGroupList` |
| Premi | `GrossPremiumList[]`, `GrossPremiumMinList[]`, `NetPremiumList`, `RnmLimitList[]` |
| Penyebaran | `RNMSpreadedListNetXOL[]`, `RNMSpreadedListNetRIXOL[]`, `RNMSpreadedListDeductXOL[]`, `RNMSpreadedListDeductRIXOL[]` |
| Penyebaran per baris | `SpreadingListXOL[]` |

`SpreadingListXOL[]` sendiri membawa **sepuluh** daftar nilai paralel: `RNMSpreadedListXOL`,
`…RIXOL`, `…GrossXOL`, `…GrossRIXOL`, `…GrossMinXOL`, `…GrossRIMinXOL`, `…NetXOL`, `…NetRIXOL`,
`…DeductXOL`, `…DeductRIXOL`, masing-masing dengan `Currency` dan `Value`.

**Catatan model:** kesepuluh daftar itu adalah **satu daftar nilai dengan dua sumbu** — jenis
besaran (bruto / bruto minimum / potongan / neto) dan pihak (OR NuRe / RI retro). Membawanya
sebagai sepuluh daftar terpisah akan mengabadikan penamaan, bukan modelnya.

---

## 4. Cabang PROPORSIONAL

### 4.1 `Limits[].Detail[]` — detail proporsional per treaty group

| Kelompok | Atribut |
|---|---|
| Identitas | `TreatyGroupID`, `TreatyGroup`, `TreatyType` |
| Quota Share | `QSPct` |
| Surplus | `Surplus` (jumlah lines) |
| Turunan kapasitas | `RetentionPct`, `CessionPct`, `IOOPct` |
| Bagian NuRe | `RNMShare`, `ShareNote` |
| Retro | `RIOGR`, `RIONR` |
| Komisi | `ProfitCommision`, `ProfitME`, `ProfitYDCF`, `PremiumReservePct` |
| Akumulasi khusus | `CurrencyRSMD` / `RSMDLimit`, `CurrencyEarthquake` / `Earthquake`, `CurrencyFloodJab` / `FloodJab`, `CurrencyFloodNat` / `FloodNation` |
| Penyebaran | `SpreadingTypeID`, `SpreadingType`, `SpreadingTotalPct` |

Daftar anak:

| Daftar | Atribut |
|---|---|
| `COBList[]` | `ClassOfBusinessID`, `ClassOfBusiness` |
| `EPIList[]` | `Currency`, `Value` |
| `RetentionList[]` | `Currency`, `Value` |
| `CessionList[]` | `Currency`, `CurrencyID`, `Value` |
| `IOOLimitList[]` | `CurrencyID`, `Currency`, `Value`, `ID` |
| `RNMShareList[]` | `Currency`, `Value` |
| `ClaimCoopList[]` | `Currency`, `Value` |
| `CashLossList[]` | `Currency`, `Value` |
| `PLAList[]` | `Currency`, `Value` |
| `ReserveList[]` | `Currency`, `Value` |
| `DeductionList[]` | `Comment`, `DeductionPct`, `DeductionPctCalculate`, `Deduction`, `Currency`, `CurrencyID` |
| `SpreadingList[]` | `ReinsTypeID`, `ReinsTypeName`, `ParentReinsTypeID`, `Pct`, `Value`, `Rp`, `Usd`, + `BreakDownSprdList[]` (`ReinsID`, `ReinsName`, `SharePct`, `Amount`, `Currency`) |
| `RNMSpreadedList[]`, `RNMSpreadedListRI[]` | `Currency`, `Value` |

**Catatan model yang mengikat:**

- `CessionPct` **pecah**: di Quota Share ia persentase sesi yang bermakna; di Surplus isinya
  jumlah lines × 100 — bukan persentase apa pun, dan persentase sesi memang tidak pernah menjadi
  syarat kontrak di Surplus.
- `RetentionPct` **tidak dibawa**: di Quota Share ia definisi (`100 − QSPct`), bukan kesepakatan.
- `DeductionList` perlu dimensi baru: **dasar perhitungan** tiap potongan (premi bruto / laba /
  jumlah yang disepakati). `DeductionPctCalculate` adalah sisa dari "turunan yang boleh disunting"
  dan tidak dibawa.
- `SpreadingList` membawa **asal-usul**: disemai dari acuan baku, diisi tangan, atau disemai lalu
  disunting (ADR-0036).

### 4.2 Fakultatif

| Benda | Atribut |
|---|---|
| `FacultativeShareList[]` | `Layer`, `LayerPart`, `LayerType`, `LayerPartType`, `Cover`, `Limit`, `Limit2`, `ClassofBusinessList`, `TreatyGroupList`, `GrossPremiumList[]`, `NetPremiumList`, `RnmLimitList[]`, `DeductionList[]` |
| `ShareFacultativeReinsurers[]` | `ReinsID`, `ReinsName`, `BrokerName`, `SharePct`, `Amount`, `Amount2`, `Layer`, `ID` |
| `…FacultativeLimits[].Detail[]` | `TreatyType`, `TreatyGroup`, `QSPct`, `RetentionPct`, `CessionPct`, `CessionList`, `IOOLimitList`, `RNMShare`, `RNMShareList`, `ShareNote` |

Struktur fakultatif **mengulang bentuk proporsional** satu tingkat lebih dalam — indikasi kuat
bahwa "bagian atas sebuah kapasitas" adalah satu bentuk yang berulang, bukan tiga bentuk berbeda.

---

## 5. Apa yang harus dikerjakan sesi model data berikutnya

1. Verifikasi inventaris ini terhadap `JSONDATA` produksi (lihat batas di bagian 0).
2. Putuskan bentuk **versi kontrak** — itu yang menggantikan empat pohon paralel di bagian 1, dan
   itu pula yang dituntut ADR-0048.
3. Petakan sepuluh daftar `RNMSpreadedList*XOL` menjadi dua sumbu.
4. Pecah `CessionPct`, buang `RetentionPct`, dan naikkan `Retention[]` menjadi masukan.
5. Beri `DeductionList` dimensi dasar perhitungan.
6. Beri setiap nilai uang **tingkat** dan **bagian yang dipakai** (ADR-0039), dan kurs beserta
   tanggal atau sumbernya (ADR-0007).
