# PETA-TELUSUR-JSON — Treaty In

**Tanggal:** 23 September 2026
**Keadaan:** hasil langkah 9.
**Ini bukan lampiran. Ini bukti kelengkapan.**

---

## 0. Kriteria terima

> **Tidak ada satu jalur pun di luar empat kategori.**

Empat kategori, tidak ada yang kelima:

| Kategori | Artinya |
|---|---|
| **DIPETAKAN** | menjadi atribut, entitas, atau bagian paket uang di model baru |
| **DITURUNKAN** | dihitung saat dibaca; tidak disimpan |
| **DIBUANG** | tidak dibawa — **wajib menyebut alasannya**; "tidak dipakai" bukan alasan |
| **DITUNDA** | menunggu sesuatu — **wajib menyebut apa yang ditunggu dan siapa pemiliknya** |

Kategori kelima adalah tempat hal-hal bersembunyi. Ia tidak ada di berkas ini.

---

## 1. Hitungan

| Kategori | Jumlah | Bagian |
|---|---|---|
| **DIPETAKAN** | 485 | 72,7 % |
| **DITURUNKAN** | 166 | 24,9 % |
| **DIBUANG** | 12 | 1,8 % |
| **DITUNDA** | 4 | 0,6 % |
| **Total** | **667** | 100 % |

**Di luar keempat kategori: nol.** Kriteria terima terpenuhi.

### Kenapa DIPETAKAN jauh lebih besar daripada jumlah atribut

485 jalur dipetakan, tetapi atribut modelnya jauh lebih sedikit. Tiga sebab, dan ketiganya bukan
penyusutan informasi:

| Sebab | Contoh |
|---|---|
| **272 jalur cermin** memetakan ke **bentuk yang sama** dengan pohon utama | `ValueDifference.*` (162), `ActualValue.*` (108) memakai entitas yang sama, bukan entitas baru |
| **pasangan nilai + mata uang** menjadi **satu** paket uang | `EGNPI[].Amount` + `EGNPI[].Currency` + `.CurrencyID` → satu paket |
| **simpul wadah** dipetakan ke entitasnya, bukan ke atribut | `Limits`, `Limits[]`, `Limits[].Detail` → entitas `LAYER` dan `DETAIL_PROPORSIONAL` |

Cermin **tidak** menjadi kategori tersendiri. Ia dipetakan, dan sasarannya kebetulan sama dengan
sasaran pohon utama. Membuatkannya kategori sendiri akan menjadi kategori kelima.

---

## 2. Pengadilan ulang — yang dicari bukan transkripsi

Adjudikasi langkah 3 dijalankan oleh pengklasifikasi berbasis pola. Pola **mencocokkan**; ia tidak
**mengadili**. Sebuah jalur dapat tercocokkan dengan benar secara kebetulan, dan sebuah jalur dapat
tercocokkan dengan salah tanpa ada yang tahu.

Maka langkah 9 dijalankan sebagai **audit**, bukan penyalinan: untuk setiap jalur dicetak **pola
alternasi mana** yang mencocokkannya, lalu disaring jalur yang polanya cocok **di dalam kata yang
berbeda** — karena di situlah kecocokan kebetulan tinggal.

**Tiga jalur lolos saringan dan harus diadili manusia. Ketiganya menghasilkan sesuatu.**

### 2.1 `ProportionType` — terklasifikasi SALAH

Tercocokkan `Proportion` dan digolongkan **turunan**. Ia **masukan**, dan bukan masukan sembarangan:
ia anggota lapisan beku `KONTRAK` dan pembeda cabang seluruh model.

Bila daftar ini dipakai apa adanya ke DDL, **sifat proporsi kontrak hilang dari lapisan yang paling
tidak boleh kehilangan apa pun**. Dikoreksi menjadi **DIPETAKAN**.

### 2.2 `ChooseStatusAkseptasi` — benar, tetapi karena alasan yang salah

Tercocokkan `StatusAkseptasi` dan ikut dibuang bersama induknya. Hasilnya kebetulan benar — ia
memang bendera pemilihan di layar, bukan fakta kontrak — tetapi ia **tidak pernah diadili**, ia
hanya ikut terbawa.

Sekarang diadili dan **tetap DIBUANG**, dengan alasannya sendiri.

### 2.3 `AdditionalAmount1` / `AdditionalAmount2` — benar, dan menyembunyikan DUA cacat

Tercocokkan `AdditionalAmount` dan digolongkan turunan. Golongannya benar. Yang disembunyikannya
tidak pernah terlihat sampai ia diadili:

```
Reinstatement_List(<LAST>).AdditionalAmount1 = @if(.MDPList(1).Currency=="IDR", .MDPList(1).Value, 0)
Reinstatement_List(<LAST>).AdditionalAmount2 = @if(.MDPList(1).Currency=="USD", .MDPList(1).Value, 0)
Reinstatement_List(<LAST>).ReinstatementAmount1 = .Limit
Reinstatement_List(<LAST>).ReinstatementAmount2 = .Limit2
```

**Cacat pertama — pasangan kolom kembar untuk dua mata uang**, di tempat yang belum pernah kami
catat. Akhiran `1` adalah slot IDR dan akhiran `2` slot USD. Ia keluarga yang sama dengan
`Limit`/`Limit2` dan `Currency`/`Currency2`, dan ia **menguatkan asumsi yang Uji J dirancang untuk
menguji**: slot 1 = IDR, slot 2 = USD.

**Cacat kedua — mata uang ketiga menjadi nol.** Kontrak yang MDP-nya bukan IDR dan bukan USD
menghasilkan **0** pada keduanya. Ini instans **ketiga** dari "kegagalan menjadi nilai" (ADR-0035),
di samping kurs yang tidak ditemukan dan `GetAchievement`. Arahnya **menghilangkan** — dan
kesalahan yang menghilangkan tidak pernah mengundang pertanyaan.

**Akibatnya pada model:** tidak ada. Keduanya tetap DITURUNKAN, dan pasangan kembar itu tidak
dibawa. **Akibatnya pada daftar temuan:** satu instans baru, dicatat.

### 2.4 Yang tidak lolos saringan, dan kenapa itu benar

Empat puluh jalur lain tercocokkan oleh pola yang juga berada "di dalam kata yang berbeda" —
`RNMSpreadedListDeductRIXOL` tercocokkan `RNMSpreadedList`, `LimitSummaryList` tercocokkan
`SummaryList`. Seluruhnya **keluarga berawalan**: kata yang lebih panjang adalah varian dari kata
yang lebih pendek, dan golongannya diwarisi dengan benar.

Pembedanya: pada ketiga jalur di §2.1–2.3, kata yang lebih panjang **berarti hal lain**
(`ProportionType` bukan varian dari `Proportion`); pada keempat puluh sisanya ia varian.

---

## 3. DIBUANG — 12, masing-masing dengan alasannya

| Jalur | Alasan |
|---|---|
| `ViewState` | turunan yang disimpan dan disetel oleh klik; dilebur ke keadaan siklus hidup (ADR-0046) |
| `IsEditData` | idem |
| `RevisionState` | idem; jalur pintas revisi dihapus (ADR-0052) |
| `Position` | idem; nilai kosongnya menandai dua keadaan yang berlawanan (ADR-0055) |
| `StatusAkseptasi` | idem |
| `ChooseStatusAkseptasi` | bendera pemilihan di layar, bukan fakta kontrak — diadili §2.2 |
| `ContractRefNo` | **tidak punya jalur pengisian sama sekali**: tak pernah ditulis, tidak ada di DDL mana pun, `pyReadOnly=true` di seluruh layar |
| `OLDID` | menampung dua arti; diganti dua relasi bernama (ADR-0040) |
| `FacultativeShareList[].Limit2` | pasangan kolom kembar untuk mata uang kedua |
| `ShareFacultativeReinsurers[].pxCreateDateTime` | fakta mesin Pega; digantikan jejak perubahan sendiri (ADR-0045) |
| `ShareFacultativeReinsurers[].pxCreateOpName` | idem |
| `BrokeragePct` | bukan properti kontrak — ia **nama parameter panggilan Oracle** di dua RDB List |

Tidak satu pun beralasan "tidak dipakai".

---

## 4. DITUNDA — 4, apa yang ditunggu dan siapa pemiliknya

| Jalur | Menunggu | Pemilik |
|---|---|---|
| `EDMMaterialType` | konfirmasi tidak ada pemakaian lain di luar yang sudah digantikan ADR-0049 | perancang model, bersama sesi Adjustment |
| `Limits[].Detail[].ProfitCommision` | satu slip treaty proporsional yang memuat klausul *profit commission* | bisnis (bagian teknik) |
| `Limits[].Detail[].ProfitME` | slip yang sama | bisnis |
| `Limits[].Detail[].ProfitYDCF` | slip yang sama | bisnis |

`CurrencyRelation` **tidak lagi ditunda**: ADR-0053 menetapkannya pasif — daftar mata uang yang
memutuskan, bukan bendera relasi. Ia kini DITURUNKAN.

---

## 4a. Pohon cermin — 279 jalur, nasibnya ditambahkan 24 September 2026 (diff `D-6`)

272 jalur cermin dan 7 skalar addendum selama ini **dikecualikan** dari adjudikasi: mereka milik
modul Adjustment. To-spec Adjustment memberi nasibnya, dan hasilnya ditambahkan di sini supaya
**satu berkas memuat seluruh nasib**.

| Nasib | Jumlah |
|---|---:|
| **DISIMPAN** | **154** |
| **TURUNAN** | 12 |
| **DIBUANG** | **113** |
| **DITUNDA** | **0** |
| **JUMLAH lingkup cermin** | **279** |

**Menutup: 279 + 388 = 667.** Perkakasnya
`../treaty-in-adjustment/tools/telusur-pohon-cermin.py`; rinciannya
`../treaty-in-adjustment/2-to-spec/PENELUSURAN-POHON-CERMIN.md`.

| Akar | Nasib | Dasar |
|---|---|---|
| `OLDDATA` | **DIBUANG** — tidak ada tabel; sisi lama di-`SELECT` ke versi dasar | ADR-0048 butir 2 · `GRL-03` · `GRL-10` |
| `ActualValue` | **DIBUANG seluruhnya** — 108 jalur, **nol kolom** | **`GRL-14`** |
| `ValueDifference` | **DISIMPAN** sebagai `NILAI_SELISIH` berkunci `KUNCI_PADANAN`; agregatnya **TURUNAN** | ADR-0048 butir 3 |
| `ValueBeforeProrate` | **DIBUANG** — nol penulis; atribut *"berlaku sejak"* dibawa, mesin pro rata **sengaja tidak dibangun** | **`GRL-15`** |
| `EDMMaterialType` | **DISIMPAN** — `SIFAT_MATERIAL_ADDENDUM`, atribut **masukan** | **`GRL-20`** |
| `EDMState` | **DISIMPAN** — `JENIS_ADDENDUM`, enum dua nilai | `GRL-13` · `GRL-18` |
| `EDMDate` | **DIBUANG** — **nol penulis di Pega**; diisi prosedur dengan `SYSDATE` pada sisip **dan** perbarui, sehingga ia tanggal **sentuh terakhir** | ADR-0006 |
| `RevisionState`, `ViewState`, `IsEditData` | **DIBUANG** — dilebur ke satu keadaan siklus hidup | ADR-0046 |
| `AddendumPremi` | **TURUNAN** — dapat dibaca dari `JENIS_ADDENDUM` | ADR-0041 · `GRL-13` |
| `OLDID` | **DISIMPAN** — `ID_VERSI_KONTRAK_DASAR`, **tepat di satu tempat** | `GRL-10` · `GRL-17` |

> ### PERINGATAN SEMESTA — ikut wajib, dan bukan formalitas
>
> **Jalur yang ditambahkan di sini TIDAK BOLEH dihitung sebagai kelengkapan.**
>
> Berkas ini tetap berlabel **"semestanya kurang"**: titik buta `L-8` menyisakan **340 properti yang
> belum diperiksa siapa pun**. Yang dapat dinyatakan: *seluruh jalur **di dalam semesta yang
> diperiksa** punya nasib.* **Bukan:** *seluruh jalur punya nasib.*
>
> Bukti terbarunya `TDA-17`: **tiga kolom — `DEDUCTIBLE2`, `PREMIUM_EARNED`, `ROL_PCT` — muncul di
> berkas ini hanya di dalam pohon cermin**, tidak pernah di pohon utama, dan karena itu tidak punya
> rumah di §10. **Itu bukan pembuangan; itu semesta yang kurang.**

---

## 5. Lubang metode yang diketahui, dan penambalnya

Penyisiran ini hanya menangkap jalur berawalan `TreatyIn.`. Daftar yang ditulis lewat *step page*
**tidak tertangkap**, dan itu menghasilkan lima entitas yang atributnya kosong atau nyaris kosong.

Penambalnya sumber yang mandiri — inventaris atribut per kelas Pega dari indeks rujukan aturan — dan
hasilnya ada di `SPEC-MODEL-DATA.md` §12.4. Properti yang **hanya** terlihat lewat penambal itu
karena itu **tidak ada di daftar di bawah**, dan itu disebut di sini supaya pembaca berikutnya tidak
membaca daftar ini sebagai daftar yang lengkap atas dirinya sendiri.

Dua yang diketahui hilang dari sapuan dan hanya terbaca lewat inventaris kelas:
`SpreadingTypeID` pada `Data-TreatyInLimitsDetail`, dan `SpreadingTypeXOL` pada `Data-TreatyInShare`.

---

## 6. Daftar penuh — 667 jalur

### DIPETAKAN — 485 jalur

```
AccountingMode
AccountingModeNonProp
AccumulationList
AccumulationList[]
AccumulationList[].Period
AccumulationList[].ReportDate
AccumulationPeriod
ActualValue
ActualValue.BrokeragePercent
ActualValue.FacultativeShare
ActualValue.FacultativeShareBrokerage
ActualValue.FacultativeShareList[].ClassofBusinessList
ActualValue.FacultativeShareList[].Cover
ActualValue.FacultativeShareList[].GrossPremiumList[].Currency
ActualValue.FacultativeShareList[].GrossPremiumList[].Value
ActualValue.FacultativeShareList[].Layer
ActualValue.FacultativeShareList[].LayerPart
ActualValue.FacultativeShareList[].LayerPartType
ActualValue.FacultativeShareList[].LayerType
ActualValue.FacultativeShareList[].Limit
ActualValue.FacultativeShareList[].Limit2
ActualValue.FacultativeShareList[].NetPremiumList
ActualValue.FacultativeShareList[].RnmLimitList[].Currency
ActualValue.FacultativeShareList[].RnmLimitList[].Value
ActualValue.FacultativeShareList[].TreatyGroupList
ActualValue.LimitShareSummaryList[].AggregateLimit
ActualValue.LimitShareSummaryList[].AggregateLimit2
ActualValue.LimitShareSummaryList[].Deductible
ActualValue.LimitShareSummaryList[].Deductible2
ActualValue.LimitShareSummaryList[].Limit
ActualValue.LimitShareSummaryList[].Limit2
ActualValue.LimitShareSummaryList[].MDP
ActualValue.LimitShareSummaryList[].MDP2
ActualValue.LimitShareSummaryList[].Note
ActualValue.Limits
ActualValue.Limits[].EgnpiTotalList[].Currency
ActualValue.Limits[].EgnpiTotalList[].CurrencyID
ActualValue.Limits[].EgnpiTotalList[].Value
ActualValue.Limits[].MDPList[].Currency
ActualValue.Limits[].MDPList[].Value
ActualValue.Limits[].PremiumEarnedList[].Currency
ActualValue.Limits[].PremiumEarnedList[].Value
ActualValue.RNMShare
ActualValue.Share[].AddendumStatus
ActualValue.Share[].ClassofBusinessList
ActualValue.Share[].Cover
ActualValue.Share[].GrossPremiumList[].Currency
ActualValue.Share[].GrossPremiumList[].Value
ActualValue.Share[].Layer
ActualValue.Share[].LayerPart
ActualValue.Share[].LayerPartType
ActualValue.Share[].LayerType
ActualValue.Share[].NetPremiumList
ActualValue.Share[].RNMSpreadedListDeductRIXOL[].Currency
ActualValue.Share[].RNMSpreadedListDeductRIXOL[].Value
ActualValue.Share[].RNMSpreadedListDeductXOL[].Currency
ActualValue.Share[].RNMSpreadedListDeductXOL[].Value
ActualValue.Share[].RNMSpreadedListGrossRIXOL[].Currency
ActualValue.Share[].RNMSpreadedListGrossRIXOL[].Value
ActualValue.Share[].RNMSpreadedListGrossXOL[].Currency
ActualValue.Share[].RNMSpreadedListGrossXOL[].Value
ActualValue.Share[].RNMSpreadedListNetRIXOL[].Currency
ActualValue.Share[].RNMSpreadedListNetRIXOL[].Value
ActualValue.Share[].RNMSpreadedListNetXOL[].Currency
ActualValue.Share[].RNMSpreadedListNetXOL[].Value
ActualValue.Share[].RNMSpreadedListRIXOL[].Currency
ActualValue.Share[].RNMSpreadedListRIXOL[].Value
ActualValue.Share[].RNMSpreadedListXOL[].Currency
ActualValue.Share[].RNMSpreadedListXOL[].Value
ActualValue.Share[].RnmLimitList[].Currency
ActualValue.Share[].RnmLimitList[].Value
ActualValue.Share[].SpreadingListXOL[].Pct
ActualValue.Share[].SpreadingListXOL[].ReinsTypeName
ActualValue.Share[].SpreadingTotalPctXOL
ActualValue.Share[].SpreadingTypeXOL
ActualValue.Share[].TreatyGroupList
ActualValue.TotalEgnpiAmount
ActualValue.TotalEgnpiAmountNP[].AltValue
ActualValue.TotalEgnpiAmountNP[].Currency
ActualValue.TotalEgnpiAmountNP[].ID
ActualValue.TotalEgnpiAmountNP[].Value
ActualValue.TotalEgnpiProportion
ActualValue.TotalFacShareDeductionNP[].Currency
ActualValue.TotalFacShareDeductionNP[].Value
ActualValue.TotalFacShareGrossNP[].Currency
ActualValue.TotalFacShareGrossNP[].Value
ActualValue.TotalFacShareNetNP[].Currency
ActualValue.TotalFacShareNetNP[].Value
ActualValue.TotalFacShareRnmNP[].Currency
ActualValue.TotalFacShareRnmNP[].Value
ActualValue.TotalLimitDeductblNP[].Currency
ActualValue.TotalLimitDeductblNP[].Value
ActualValue.TotalLimitIOONP[].Currency
ActualValue.TotalLimitIOONP[].Value
ActualValue.TotalLimitMDPNP[].Currency
ActualValue.TotalLimitMDPNP[].Value
ActualValue.TotalLimitPremiEarnNP[].Currency
ActualValue.TotalLimitPremiEarnNP[].Value
ActualValue.TotalLimitsROL
ActualValue.TotalShareDeductionNP[].Currency
ActualValue.TotalShareDeductionNP[].Value
ActualValue.TotalShareGrossNP[].Currency
ActualValue.TotalShareGrossNP[].Value
ActualValue.TotalShareNetNP[].Currency
ActualValue.TotalShareNetNP[].Value
ActualValue.TotalShareRnmNP[].Currency
ActualValue.TotalShareRnmNP[].Value
ActualValue.TotalSpreadedNetPremiRI[].Currency
ActualValue.TotalSpreadedNetPremiRI[].Value
ActualValue.TotalSpreadedNetPremi[].Currency
ActualValue.TotalSpreadedNetPremi[].Value
ActualValue.TotalSpreadedRnmProp[].Currency
ActualValue.TotalSpreadedRnmProp[].Value
ActualValue.TotalSpreadedRnmRIProp[].Currency
ActualValue.TotalSpreadedRnmRIProp[].Value
Bordeaux
BordereauxNote
BrokeragePercent
BrokeragePercentP
Ceding
CedingID
CedingStatusActive
ClassofBusiness
CoInScale
Commencement
Comment
CommentList
CommentList[]
CommentList[].Date
CommentList[].IsApproved
CommentList[].OperatorName
CommentList[].Suggest
Currency
CurrencyEarthquake
CurrencyEgnpiAmount
CurrencyFloodJab
CurrencyFloodNat
CurrencyInstallmentAmount
CurrencyList
CurrencyList[]
CurrencyList[].Conversion
CurrencyList[].Currency
CurrencyList[].CurrencyID
CurrencyList[].PeriodEnd
CurrencyList[].PeriodStart
CurrencyRSMD
EDMEffective
EDMState
EGNPI
EGNPI[]
EGNPI[].Amount
EGNPI[].AsDate
EGNPI[].Currency
EGNPI[].CurrencyID
EGNPI[].ID
EGNPI[].TreatyGroup
EGNPI[].TreatyGroupID
Earthquake
Exclusions
ExclusionsP
FacShare
FacShareBrokerage
FacultativeShare
FacultativeShareBrokerage
FacultativeShareList
FacultativeShareList[]
FacultativeShareList[].ClassofBusinessList
FacultativeShareList[].Cover
FacultativeShareList[].GrossPremiumList[].Currency
FacultativeShareList[].GrossPremiumList[].Value
FacultativeShareList[].Layer
FacultativeShareList[].LayerPart
FacultativeShareList[].LayerPartType
FacultativeShareList[].LayerType
FacultativeShareList[].Limit
FacultativeShareList[].ShareFacultativeReinsurers[].Amount
FacultativeShareList[].ShareFacultativeReinsurers[].Amount2
FacultativeShareList[].ShareFacultativeReinsurers[].ReinsID
FacultativeShareList[].ShareFacultativeReinsurers[].ReinsName
FacultativeShareList[].ShareFacultativeReinsurers[].SharePct
FacultativeShareList[].TreatyGroupList
FloodJab
FloodNation
ID
Information
Installment
InstallmentNo
Installment[]
Installment[].Currency
Installment[].ID
IsMultipleRetro
IsProRate
LeadingReinsID
LeadingReinsName
LeadingReinsSource
LeadingReinsSourceID
Limits
Limits[]
Limits[].AdjRate
Limits[].Cover
Limits[].Currency
Limits[].Deductible
Limits[].Detail[].COBList[].ClassOfBusiness
Limits[].Detail[].COBList[].ClassOfBusinessID
Limits[].Detail[].CashLossList[].Currency
Limits[].Detail[].CashLossList[].Value
Limits[].Detail[].ClaimCoopList[].Currency
Limits[].Detail[].ClaimCoopList[].Value
Limits[].Detail[].EPIList[].Currency
Limits[].Detail[].EPIList[].Value
Limits[].Detail[].IOOLimitList[].Currency
Limits[].Detail[].IOOLimitList[].CurrencyID
Limits[].Detail[].IOOLimitList[].ID
Limits[].Detail[].IOOLimitList[].Value
Limits[].Detail[].PLAList[].Currency
Limits[].Detail[].PLAList[].Value
Limits[].Detail[].QSPct
Limits[].Detail[].RIOGR
Limits[].Detail[].RIONR
Limits[].Detail[].RetentionList[].Currency
Limits[].Detail[].RetentionList[].Value
Limits[].Detail[].Surplus
Limits[].Detail[].TreatyGroup
Limits[].Detail[].TreatyGroupID
Limits[].Detail[].TreatyType
Limits[].ID
Limits[].Layer
Limits[].LayerPart
Limits[].LayerPartType
Limits[].LayerType
Limits[].Limit
Limits[].MDPPct
Limits[].ReinstatementPct
Limits[].ReinstatementValue
Limits[].TreatyGroupList[].ClassOfBusinessList[].ClassOfBusiness
Limits[].TreatyGroupList[].ClassOfBusinessList[].ClassOfBusinessID
Limits[].TreatyGroupList[].TreatyGroup
Limits[].TreatyGroupList[].TreatyGroupID
Limits[].TreatyType
MaxCoGroup
MaxCoNonGroup
NusareSharePct
OLDDATA
OptionLimit
Portfolio
Portfolio[]
Portfolio[].Description
Portfolio[].Type
Portfolio[].TypePortfolio
PositionUsername
ProportionType
RNMShare
RNMShareAcrossTheBoard
RNMShareP
RSMDLimit
ReminderDays
ReportingConfirmation
ReportingEnd
ReportingInterval
ReportingPeriod
ReportingPeriodList
ReportingPeriodList[]
ReportingPeriodList[].ConfirmationDue
ReportingPeriodList[].InitialDate
ReportingPeriodList[].Period
ReportingPeriodList[].SettlementDue
ReportingPeriodList[].SubmissionDue
ReportingSettlement
ReportingStart
ReportingSubmission
Retention
Retention[]
Retention[].Amount
Retention[].Currency
Retention[].CurrencyID
Retention[].ID
Retention[].TreatyGroup
Retention[].TreatyGroupID
RetroList
RnmShareDeducted
Share
ShareFacultativeReinsurers
ShareFacultativeReinsurers[]
ShareFacultativeReinsurers[].BrokerName
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].CessionList
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].IOOLimitList
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].QSPct
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].RNMShare
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].ShareNote
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].TreatyGroup
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].TreatyGroupID
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].TreatyType
ShareFacultativeReinsurers[].FacultativeLimits[].TreatyType
ShareFacultativeReinsurers[].FacultativeLimits[].TreatyTypeID
ShareFacultativeReinsurers[].ID
ShareFacultativeReinsurers[].Layer
ShareFacultativeReinsurers[].ReinsID
ShareFacultativeReinsurers[].ReinsName
ShareFacultativeReinsurers[].SharePct
ShareReins
ShareReins[]
ShareReins[].ID
Share[]
Share[].ClassofBusinessList
Share[].Cover
Share[].GrossPremiumList[].Currency
Share[].GrossPremiumList[].Value
Share[].Layer
Share[].LayerPart
Share[].LayerPartType
Share[].LayerType
Share[].SpreadingListXOL[].GrossPremiumList[].Currency
Share[].SpreadingListXOL[].GrossPremiumList[].Value
Share[].TreatyGroupList
SourceStatusActive
SpecialConditions
SpecialConditionsP
TeritorialScope
Termination
TreatyContractName
TreatyLeader
TreatyYear
ValueBeforeProrate
ValueDifference
ValueDifference.BrokeragePercent
ValueDifference.EGNPI[].AddendumStatus
ValueDifference.EGNPI[].Amount
ValueDifference.EGNPI[].AmountIDR
ValueDifference.EGNPI[].Currency
ValueDifference.EGNPI[].Proportion
ValueDifference.EGNPI[].TreatyGroup
ValueDifference.FacultativeShare
ValueDifference.FacultativeShareBrokerage
ValueDifference.FacultativeShareList[].Cover
ValueDifference.FacultativeShareList[].DeductionList[].Comment
ValueDifference.FacultativeShareList[].DeductionList[].Currency
ValueDifference.FacultativeShareList[].DeductionList[].Deduction
ValueDifference.FacultativeShareList[].DeductionList[].DeductionPct
ValueDifference.FacultativeShareList[].GrossPremiumList[].Currency
ValueDifference.FacultativeShareList[].GrossPremiumList[].Value
ValueDifference.FacultativeShareList[].Layer
ValueDifference.FacultativeShareList[].LayerPart
ValueDifference.FacultativeShareList[].LayerPartType
ValueDifference.FacultativeShareList[].LayerType
ValueDifference.FacultativeShareList[].Limit
ValueDifference.FacultativeShareList[].Limit2
ValueDifference.FacultativeShareList[].NetPremiumList[].Currency
ValueDifference.FacultativeShareList[].NetPremiumList[].Value
ValueDifference.FacultativeShareList[].RnmGrossPremiDisplay[].Currency
ValueDifference.FacultativeShareList[].RnmGrossPremiDisplay[].Value
ValueDifference.FacultativeShareList[].RnmLimitListDisplay[].Currency
ValueDifference.FacultativeShareList[].RnmLimitListDisplay[].Value
ValueDifference.FacultativeShareList[].RnmLimitList[].Currency
ValueDifference.FacultativeShareList[].RnmLimitList[].Value
ValueDifference.Installment[].Currency
ValueDifference.Installment[].InstallmentList[].Currency
ValueDifference.Installment[].InstallmentList[].CurrencyID
ValueDifference.Installment[].InstallmentList[].DueDate
ValueDifference.Installment[].InstallmentList[].Installment
ValueDifference.Installment[].InstallmentList[].InstallmentPct
ValueDifference.Installment[].InstallmentList[].PaymentDate
ValueDifference.Installment[].InstallmentList[].WPC
ValueDifference.LimitSummaryList[].AggregateLimit
ValueDifference.LimitSummaryList[].AggregateLimit2
ValueDifference.LimitSummaryList[].Deductible
ValueDifference.LimitSummaryList[].Deductible2
ValueDifference.LimitSummaryList[].Limit
ValueDifference.LimitSummaryList[].Limit2
ValueDifference.LimitSummaryList[].MDP
ValueDifference.LimitSummaryList[].MDP2
ValueDifference.LimitSummaryList[].Note
ValueDifference.Limits[].AdjRate
ValueDifference.Limits[].Cover
ValueDifference.Limits[].Currency
ValueDifference.Limits[].Currency2
ValueDifference.Limits[].Deductible
ValueDifference.Limits[].Deductible2
ValueDifference.Limits[].EgnpiTotalList[].Currency
ValueDifference.Limits[].EgnpiTotalList[].Value
ValueDifference.Limits[].Layer
ValueDifference.Limits[].LayerPart
ValueDifference.Limits[].LayerPartType
ValueDifference.Limits[].LayerType
ValueDifference.Limits[].Limit
ValueDifference.Limits[].Limit2
ValueDifference.Limits[].MDPList[].Currency
ValueDifference.Limits[].MDPList[].Value
ValueDifference.Limits[].MDPPct
ValueDifference.Limits[].PremiumEarnedList[].Currency
ValueDifference.Limits[].PremiumEarnedList[].Value
ValueDifference.Limits[].ROLPct
ValueDifference.Limits[].ReinstatementPct
ValueDifference.Limits[].TreatyGroupList
ValueDifference.RNMShare
ValueDifference.Share[].AddendumStatus
ValueDifference.Share[].Cover
ValueDifference.Share[].DeductionList[].Comment
ValueDifference.Share[].DeductionList[].Currency
ValueDifference.Share[].DeductionList[].Deduction
ValueDifference.Share[].DeductionList[].DeductionPct
ValueDifference.Share[].GrossPremiumList[].Currency
ValueDifference.Share[].GrossPremiumList[].Value
ValueDifference.Share[].Layer
ValueDifference.Share[].LayerPart
ValueDifference.Share[].LayerPartType
ValueDifference.Share[].LayerType
ValueDifference.Share[].NetPremiumList[].Currency
ValueDifference.Share[].NetPremiumList[].Value
ValueDifference.Share[].RNMSpreadedListDeductRIXOL[].Currency
ValueDifference.Share[].RNMSpreadedListDeductRIXOL[].Value
ValueDifference.Share[].RNMSpreadedListDeductXOL[].Currency
ValueDifference.Share[].RNMSpreadedListDeductXOL[].Value
ValueDifference.Share[].RNMSpreadedListGrossMinXOL[].Currency
ValueDifference.Share[].RNMSpreadedListGrossMinXOL[].Value
ValueDifference.Share[].RNMSpreadedListGrossRIMinXOL[].Currency
ValueDifference.Share[].RNMSpreadedListGrossRIMinXOL[].Value
ValueDifference.Share[].RNMSpreadedListGrossRIXOL[].Currency
ValueDifference.Share[].RNMSpreadedListGrossRIXOL[].Value
ValueDifference.Share[].RNMSpreadedListGrossXOL[].Currency
ValueDifference.Share[].RNMSpreadedListGrossXOL[].Value
ValueDifference.Share[].RNMSpreadedListNetRIXOL[].Currency
ValueDifference.Share[].RNMSpreadedListNetRIXOL[].Value
ValueDifference.Share[].RNMSpreadedListNetXOL[].Currency
ValueDifference.Share[].RNMSpreadedListNetXOL[].Value
ValueDifference.Share[].RNMSpreadedListRIXOL[].Currency
ValueDifference.Share[].RNMSpreadedListRIXOL[].Value
ValueDifference.Share[].RNMSpreadedListXOL[].Currency
ValueDifference.Share[].RNMSpreadedListXOL[].Value
ValueDifference.Share[].RnmGrossPremiDisplay[].Currency
ValueDifference.Share[].RnmGrossPremiDisplay[].Value
ValueDifference.Share[].RnmLimitListDisplay[].Currency
ValueDifference.Share[].RnmLimitListDisplay[].Value
ValueDifference.Share[].RnmLimitList[].Currency
ValueDifference.Share[].RnmLimitList[].Value
ValueDifference.Share[].SpreadingListXOL[].ParentReinsTypeID
ValueDifference.Share[].SpreadingListXOL[].Pct
ValueDifference.Share[].SpreadingListXOL[].ReinsTypeID
ValueDifference.Share[].SpreadingListXOL[].ReinsTypeName
ValueDifference.Share[].SpreadingListXOL[].Rp
ValueDifference.Share[].SpreadingListXOL[].Usd
ValueDifference.Share[].SpreadingTotalPctXOL
ValueDifference.Share[].SpreadingTypeIDXOL
ValueDifference.Share[].SpreadingTypeXOL
ValueDifference.Share[].TreatyGroupList
ValueDifference.TotalEgnpiAmount
ValueDifference.TotalEgnpiAmountNP[].AltValue
ValueDifference.TotalEgnpiAmountNP[].Currency
ValueDifference.TotalEgnpiAmountNP[].ID
ValueDifference.TotalEgnpiAmountNP[].Value
ValueDifference.TotalEgnpiProportion
ValueDifference.TotalFacShareDeductionNP[].Currency
ValueDifference.TotalFacShareDeductionNP[].Value
ValueDifference.TotalFacShareGrossNP[].Currency
ValueDifference.TotalFacShareGrossNP[].Value
ValueDifference.TotalFacShareNetNP[].Currency
ValueDifference.TotalFacShareNetNP[].Value
ValueDifference.TotalFacShareRnmNP[].Currency
ValueDifference.TotalFacShareRnmNP[].Value
ValueDifference.TotalInstallmentNP[].Currency
ValueDifference.TotalInstallmentNP[].Value
ValueDifference.TotalLimitDeductblNP[].Currency
ValueDifference.TotalLimitDeductblNP[].Value
ValueDifference.TotalLimitIOONP[].Currency
ValueDifference.TotalLimitIOONP[].Value
ValueDifference.TotalLimitMDPNP[].Currency
ValueDifference.TotalLimitMDPNP[].Value
ValueDifference.TotalLimitPremiEarnNP[].Currency
ValueDifference.TotalLimitPremiEarnNP[].Value
ValueDifference.TotalLimitsROL
ValueDifference.TotalShareDeductionNP[].Currency
ValueDifference.TotalShareDeductionNP[].Value
ValueDifference.TotalShareGrossNP[].Currency
ValueDifference.TotalShareGrossNP[].Value
ValueDifference.TotalShareNetNP[].Currency
ValueDifference.TotalShareNetNP[].Value
ValueDifference.TotalShareRnmNP[].Currency
ValueDifference.TotalShareRnmNP[].Value
ValueDifference.TotalSpreadedNetPremiRI[].Currency
ValueDifference.TotalSpreadedNetPremiRI[].Value
ValueDifference.TotalSpreadedNetPremi[].Currency
ValueDifference.TotalSpreadedNetPremi[].Value
ValueDifference.TotalSpreadedRnmProp[].Currency
ValueDifference.TotalSpreadedRnmProp[].Value
ValueDifference.TotalSpreadedRnmRIProp[].Currency
ValueDifference.TotalSpreadedRnmRIProp[].Value
```

### DITURUNKAN — 166 jalur

```
EGNPI[].AmountIDR
FacultativeShareList[].NetPremiumList
FacultativeShareList[].RnmLimitList[].Currency
FacultativeShareList[].RnmLimitList[].Value
LimitFacShareSummaryList
LimitFacShareSummaryList[]
LimitShareSummaryList
LimitShareSummaryList[]
LimitSummaryList
LimitSummaryList[]
Limits[].EgnpiTotalList[].Currency
Limits[].EgnpiTotalList[].CurrencyID
Limits[].EgnpiTotalList[].Value
Limits[].MDPList[].Currency
Limits[].MDPList[].Value
Limits[].PremiumEarnedList[].Currency
Limits[].PremiumEarnedList[].Value
Limits[].ROLPct
Limits[].Reinstatement_List[].AdditionalAmount1
Limits[].Reinstatement_List[].AdditionalAmount2
MDPSummaryList
MDPSummaryList[]
ProRateDays
ProRatePercent
ProRateTotalDays
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].CessionPct
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].RNMShareList
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].RetentionPct
Share[].GrossPremiumMinList[].Currency
Share[].GrossPremiumMinList[].Value
Share[].NetPremiumList
Share[].RNMSpreadedListDeductRIXOL[].Currency
Share[].RNMSpreadedListDeductRIXOL[].Value
Share[].RNMSpreadedListDeductXOL[].Currency
Share[].RNMSpreadedListDeductXOL[].Value
Share[].RNMSpreadedListNetRIXOL[].Currency
Share[].RNMSpreadedListNetRIXOL[].Value
Share[].RNMSpreadedListNetXOL[].Currency
Share[].RNMSpreadedListNetXOL[].Value
Share[].RnmLimitList[].Currency
Share[].RnmLimitList[].Value
Share[].SpreadingListXOL[].RNMSpreadedListDeductRIXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListDeductRIXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListDeductXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListDeductXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListGrossMinXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListGrossMinXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListGrossRIMinXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListGrossRIMinXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListGrossRIXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListGrossRIXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListGrossXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListGrossXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListNetRIXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListNetRIXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListNetXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListNetXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListRIXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListRIXOL[].Value
Share[].SpreadingListXOL[].RNMSpreadedListXOL[].Currency
Share[].SpreadingListXOL[].RNMSpreadedListXOL[].Value
TotalEgnpiAmount
TotalEgnpiAmountNP
TotalEgnpiAmountNP[]
TotalEgnpiAmountNP[].AltValue
TotalEgnpiAmountNP[].Currency
TotalEgnpiAmountNP[].ID
TotalEgnpiAmountNP[].Value
TotalEgnpiProportion
TotalFacShareDeductionNP
TotalFacShareDeductionNP[]
TotalFacShareDeductionNP[].Currency
TotalFacShareDeductionNP[].Value
TotalFacShareGrossNP
TotalFacShareGrossNP[]
TotalFacShareGrossNP[].Currency
TotalFacShareGrossNP[].Value
TotalFacShareNetNP
TotalFacShareNetNP[]
TotalFacShareNetNP[].Currency
TotalFacShareNetNP[].Value
TotalFacShareRnmNP
TotalFacShareRnmNP[]
TotalFacShareRnmNP[].Currency
TotalFacShareRnmNP[].Value
TotalInstallmentAmount
TotalInstallmentNP
TotalInstallmentNP[]
TotalInstallmentNP[].Currency
TotalInstallmentNP[].Value
TotalInstallmentPct
TotalLimitDeductblNP
TotalLimitDeductblNP[]
TotalLimitDeductblNP[].Currency
TotalLimitDeductblNP[].Value
TotalLimitIOONP
TotalLimitIOONP[]
TotalLimitIOONP[].Currency
TotalLimitIOONP[].Value
TotalLimitMDPMinNP
TotalLimitMDPMinNP[]
TotalLimitMDPMinNP[].Currency
TotalLimitMDPMinNP[].Value
TotalLimitMDPNP
TotalLimitMDPNP[]
TotalLimitMDPNP[].Currency
TotalLimitMDPNP[].Value
TotalLimitPremiEarnNP
TotalLimitPremiEarnNP[]
TotalLimitPremiEarnNP[].Currency
TotalLimitPremiEarnNP[].Value
TotalLimitsAdjPct
TotalLimitsDeductible
TotalLimitsIOOLimit
TotalLimitsMdp
TotalLimitsPremiumEarned
TotalLimitsROL
TotalRetentionAmount
TotalRetentionAmountNP
TotalRetentionAmountNP[]
TotalRetentionAmountNP[].Currency
TotalRetentionAmountNP[].ID
TotalRetentionAmountNP[].Value
TotalShareDeductionNP
TotalShareDeductionNP[]
TotalShareDeductionNP[].Currency
TotalShareDeductionNP[].Value
TotalShareGross
TotalShareGrossMinNP
TotalShareGrossMinNP[]
TotalShareGrossMinNP[].Currency
TotalShareGrossMinNP[].Value
TotalShareGrossNP
TotalShareGrossNP[]
TotalShareGrossNP[].Currency
TotalShareGrossNP[].Value
TotalShareNet
TotalShareNetNP
TotalShareNetNP[]
TotalShareNetNP[].Currency
TotalShareNetNP[].Value
TotalShareRnmLimit
TotalShareRnmNP
TotalShareRnmNP[]
TotalShareRnmNP[].Currency
TotalShareRnmNP[].Value
TotalShareRnmProp
TotalShareRnmProp[]
TotalShareRnmProp[].Currency
TotalShareRnmProp[].Value
TotalSpreadedNetPremi
TotalSpreadedNetPremiRI
TotalSpreadedNetPremiRI[]
TotalSpreadedNetPremiRI[].Currency
TotalSpreadedNetPremiRI[].Value
TotalSpreadedNetPremi[]
TotalSpreadedNetPremi[].Currency
TotalSpreadedNetPremi[].Value
TotalSpreadedRnmProp
TotalSpreadedRnmProp[]
TotalSpreadedRnmProp[].Currency
TotalSpreadedRnmProp[].Value
TotalSpreadedRnmRIProp
TotalSpreadedRnmRIProp[]
TotalSpreadedRnmRIProp[].Currency
TotalSpreadedRnmRIProp[].Value
```

### DIBUANG — 12 jalur

```
BrokeragePct
ChooseStatusAkseptasi
ContractRefNo
FacultativeShareList[].Limit2
IsEditData
OLDID
Position
RevisionState
ShareFacultativeReinsurers[].pxCreateDateTime
ShareFacultativeReinsurers[].pxCreateOpName
StatusAkseptasi
ViewState
```

### DITUNDA — 4 jalur

```
EDMMaterialType
Limits[].Detail[].ProfitCommision
Limits[].Detail[].ProfitME
Limits[].Detail[].ProfitYDCF
```
