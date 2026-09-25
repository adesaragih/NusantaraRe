# Struktur Data Klaim Sistem Lama — pohon `.ClaimData` sampai turunan terdalam

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop` — **279 berkas XML**, disapu seluruhnya (bukan cuplikan).
> Pembanding lintas modul diambil dari `D:\XML_NURE\Komite Claim Non Prop` hanya untuk mengonfirmasi bentuk kelas yang dipakai bersama; tidak ada properti modul Komite yang ikut masuk pohon ini.
> Dibangkitkan `alat/buat-pohon-claimdata.py` — sapuan **19 September 2026**.
>
> **Berkas ini memerikan SISTEM LAMA, bukan rancangan sistem baru.** Nama properti ditulis apa adanya seperti di XML, termasuk salah ejanya (`CurencyAdjustment`, `ComiteeClaim`, `CNPCircumtances`, `Comitee`). Penerjemahan ke istilah CONTEXT.md terjadi di berkas lain, bukan di sini.
>
> **TURUNAN.** Keempat CSV di §1.4 adalah pendamping mesin-baca berkas ini. Bila CSV berbeda dari sumber XML-nya, **alatnya yang salah**.

Pendamping mesin-baca: **empat CSV**, lihat §1.4.
Melengkapi [`BLUEPRINT.md`](../1-grilling/BLUEPRINT.md) §2.1 yang berhenti di **daftar PageList tingkat satu**. Berkas ini meneruskannya sampai **kedalaman 4** dan memberi kamus field per simpul.

---

## 1. Cara baca

### 1.1 Notasi

| Tulisan | Arti |
|---|---|
| `.Prop` | properti tunggal (skalar) |
| `Prop[]` | **Page List** — baris berulang, selalu dirujuk dengan indeks `(...)` di XML |
| `Prop{}` | **Page** — satu halaman bersarang, tanpa indeks |
| `→ KELAS` | kelas Pega yang menaungi isi simpul, diambil dari `pyStepsClassName` / `pyPagesAndClasses` di berkas Activity |

### 1.2 Kolom `BUKTI` di CSV

Tiga tingkat keyakinan. Bedanya penting — jangan diratakan.

| Nilai | Artinya | Cacah |
|---|---|---:|
| `jalur-penuh` | jalur lengkapnya terbaca utuh di XML, mis. `pyWorkPage.ClaimData.SpreadingRisk(1).ClaimEstimation` | 366 |
| `rujukan-relatif` | hanya terbaca sebagai `.Prop` di dalam repeat Section yang sumbernya jelas; induknya disimpulkan dari sumber repeat itu | 50 |
| `alias-primary` | hanya terbaca lewat alias step-page `Primary`, di rule yang applies-to `Data-Adjustment` — jadi `Primary` = page `Adjustment` | 48 |
| `identitas-kelas` | tidak pernah ditulis lewat jalur ini, tapi **pasti ada** karena simpulnya berkelas sama dengan simpul lain yang isinya terbaca — lihat §5 | 25 |

### 1.3 Tipe data

Pega tidak mengekspor `Rule-Obj-Property`, jadi **tidak ada satu pun tipe yang dinyatakan langsung oleh sistem lama.** Kolom `JENIS` di CSV adalah **dugaan**, disusun dari:
1. bentuk pemakaian di `Property-Set` (aritmetika → `Decimal`, `@CurrentDateTime()` → `DateTime`, isian `"true"`/`"0"` → `Boolean`/`Flag`),
2. daftar override eksplisit untuk nama yang menyesatkan heuristik (`NoClaim` itu teks, bukan angka).

Perlakukan `JENIS` sebagai **titik awal pemetaan, bukan kebenaran**. Yang bisa dipegang dari berkas ini adalah **bentuk pohon dan nama**, bukan presisi angkanya — presisi ada di ADR-0003 dan `datar-lama-aritmetika.csv`.

### 1.4 Empat CSV pendamping

Semua tabel di berkas ini punya bentuk mesin-baca. Dibangkitkan sekaligus oleh `alat/buat-pohon-claimdata.py`.

| Berkas | Baris | Isi | Kolom |
|---|---:|---|---|
| [`datar-claimdata-lama.csv`](datar-claimdata-lama.csv) | 489 | **pohon properti** — satu baris per simpul | `JALUR` · `KEDALAMAN` · `NAMA` · `JENIS` · `KELAS_PEGA` · `CACAH_ANAK` · `REF` · `CACAH_BERKAS` · `BUKTI` · `BERKAS_CONTOH` |
| [`datar-claimdata-kelas.csv`](datar-claimdata-kelas.csv) | 10 | **peta kelas** — kelas Pega dan simpul yang memakainya | `KELAS_PEGA` · `CACAH_PEMAKAIAN` · `DIPAKAI_SEBAGAI` |
| [`datar-claimdata-salinan.csv`](datar-claimdata-salinan.csv) | 6 | **salinan beku** ke dalam `AdjustmentList[]` (§6) | `DITULIS_KE` · `DISALIN_DARI` · `KELAS` · `GANTI_NAMA` · `BERKAS` |
| [`datar-claimdata-muatan-keluar.csv`](datar-claimdata-muatan-keluar.csv) | 94 | **muatan keluar** — tiap isian ke `InputParamOs` / `TempPNC` / `ParamKasir` (§7) | `WADAH` · `FIELD` · `DIISI_DARI` · `HALAMAN_LANGKAH` · `BERKAS` |

Catatan pemakaian:

- `JALUR` memakai **nama kanonik tanpa indeks**. `ClaimData.SpreadingRisk.ClaimEstimation` berarti `ClaimData.SpreadingRisk(n).ClaimEstimation` di XML; `JENIS = Page List` pada simpul induknya yang menandai perulangan, bukan tanda kurung di jalurnya.
- `REF` adalah cacah rujukan mentah, `CACAH_BERKAS` cacah berkas berbeda yang merujuknya. `REF = 0` wajar untuk baris ber-`BUKTI` `rujukan-relatif` atau `identitas-kelas` — jalur kanoniknya memang tidak pernah ditulis.
- `BERKAS_CONTOH` memuat maksimal **tiga** berkas sumber, dipisah ` | `. Untuk daftar penuh, jalankan ulang alatnya.
- `datar-claimdata-muatan-keluar.csv` sengaja **tidak dideduplikasi**: satu field bisa muncul beberapa kali karena diisi di beberapa Activity, dan sering dua kali dalam satu Activity (isi awal, lalu dikurangi nilai tercatat — pola selisih di §7.1). Bentuk itulah yang ingin terlihat.

---

## 2. Ringkasan angka

| | |
|---|---:|
| Simpul dalam pohon `.ClaimData` | **393** |
| Simpul dalam pohon `.Adjustment` (cabang sejajar, §5) | **96** |
| Kedalaman maksimum | **4** tingkat di bawah `ClaimData` |
| Properti skalar langsung di `.ClaimData` | **91** |
| Kontainer (Page/Page List) langsung di `.ClaimData` | **22** |
| Kelas Pega berbeda yang menyusun pohon | **9** (+7 tak dideklarasikan) |

---

## 3. Peta kelas

Sembilan kelas menyusun seluruh pohon. **Satu kelas dipakai ulang di banyak tempat** — inilah sebab pohonnya terlihat berulang-ulang.

| Kelas Pega | Dipakai sebagai | Catatan |
|---|---|---|
| `ASM-FW-GCNMFW-Data-ClaimData` | `.ClaimData` | akar |
| `ASM-FW-GISFW-Data-SpreadingRisk` | `SpreadingRisk`, `CNPSpreadLoss`, `SpreadingClaim`, `SpreadingBreakQS`, `SpreadingAdjustment`, `SpreadingAdjustmentQS`, `SpreadingQuotaShare`, `AlokasiXOLPaid`, `LossAllocation` | **9 pemakaian** — satu bentuk baris, sembilan peran berbeda |
| `ASM-FW-GISFW-Data-TreatyInTotal` | `ListClaimAmount`, `ListTotalEstimation`, `ReinstatementList`, `TotalInterestInsured`, `ListClaimAcceptation` | **5 pemakaian** |
| `ASM-FW-GCNMFW-Data-Adjustment` | `AdjustmentList`, `Adjustment`, `CNPLayerList`, `CNPCurrencyList`, `CurencyAdjustment` | **5 pemakaian, termasuk bersarang di dirinya sendiri** — lihat §5 |
| `ASM-FW-GCNMFW-Data-Comitee` | `ComiteeClaim`, `ClaimComitee` | `KomiteList` juga berkelas ini tapi menempel di work page, di luar pohon `.ClaimData` |
| `ASM-FW-GCNMFW-Data-ObjectItem` | `InterestList`, `ObjectList`, `ObjectItemList` | |
| `ASM-FW-GCNMFW-Data-Estimasi` | `EstimationList` | |
| `ASM-FW-GISFW-Data-OfferFacIn-SuggestList` | `SuggestList` | kelas pinjaman dari domain fakultatif |
| `ASM-FW-GISFW-Int-RETROCESSIONLIFE` | `RetroList` | kelas integrasi, bukan kelas data |

**Tujuh simpul tidak punya deklarasi kelas di seluruh ekspor**: `ReceiverClaim`, `Attachment`, `AttachmentPaid`, `PolicyData`, `QuotationData`, `MarketingData`, `DataCommitteeTreaty`. Simpulnya nyata dan isinya terbaca; hanya nama kelasnya yang tidak ada di 279 berkas ini. Jangan tebak nama kelasnya — pakai isinya.

> **Konsekuensi untuk pemetaan tabel.** Satu kelas Pega ≠ satu tabel. `ASM-FW-GISFW-Data-SpreadingRisk` muncul di sembilan tempat dengan arti bisnis berbeda (alokasi layer, pembagian per porsi, penyebaran ke retrosesi, salinan beku di dalam Adjustment). Memetakannya jadi satu tabel akan menggabungkan hal-hal yang tidak boleh digabung.

---

## 4. Pohon `.ClaimData`

Angka `n=` adalah cacah rujukan di ekspor — penanda seberapa hidup simpul itu, bukan penanda penting atau tidaknya.

```
ClaimData{}                                    → ASM-FW-GCNMFW-Data-ClaimData        n=1743
│
├─ 91 properti skalar (identitas, tanggal, parameter hitung)  ...................  §4.1
│
├─ PolicyData{}                                → (tidak dideklarasikan)              n=144
│    ├─ PolicyNo · StartDateTime · EndDateTime · SDateTime · TreatyGroup
│
├─ QuotationData{}                             → (tidak dideklarasikan)              n=15
│    ├─ BusinessOldId · BusinessName · CedingCoName · SobName
│
├─ MarketingData{}                             → (tidak dideklarasikan)              n=7
│    ├─ ID · ClientID · ClientName · TeamGroup · BranchDetailID · BranchDetailName
│
├─ AttachmentPaid{}                            → (tidak dideklarasikan)              n=4
│    └─ File
│
├─ SpreadingRisk[]                             → GISFW-Data-SpreadingRisk            n=215  ★ inti alokasi XOL
│    ├─ 29 field (lihat §4.2)
│    └─ RetroList[]                            → GISFW-Int-RETROCESSIONLIFE          n=4
│         └─ TREATYTYPENAME · PERCENTSHARE · OVR_COMM
│
├─ ReceiverClaim[]                             → (tidak dideklarasikan)              n=101
│    └─ 13 field — sepasang rekening (utama + varian `…2`)
│
├─ ListClaimAmount[]                           → GISFW-Data-TreatyInTotal            n=80
│    └─ 15 field — nilai kerugian per mata uang
│
├─ AdjustmentList[]                            → GCNMFW-Data-Adjustment              n=60   ★ unit transaksi
│    ├─ 47 field skalar (lihat §4.3)
│    ├─ ComiteeClaim[]                         → GCNMFW-Data-Comitee
│    │    └─ KomiteID · IDKomite · KomiteEmail · KomitePost · Initial · KomiteAproval
│    │       · KomiteComment · DateApproval · DateApprove
│    ├─ DataCommitteeTreaty{}                  → (tidak dideklarasikan)
│    │    └─ CircumCauseOfLoss · ExtentOfLoss · LegalLiability · Occupation · Remarks
│    ├─ CNPLayerList[]                         → GCNMFW-Data-Adjustment       [identitas-kelas]
│    │    ├─ XOL · XOLID
│    │    └─ CNPCurrencyList[]                 → GCNMFW-Data-Adjustment       ◄── KEDALAMAN 4
│    │         └─ 21 field (lihat §4.4)
│    ├─ CurencyAdjustment[]                    → GCNMFW-Data-Adjustment
│    │    └─ CurrencyID
│    ├─ LossAllocation[]          ═ salinan beku dari ClaimData.SpreadingRisk   (§6)
│    ├─ SpreadingRisk[]           ═ salinan beku dari ClaimData.SpreadingRisk   (§6)
│    ├─ SpreadingAdjustment[]     ═ salinan beku dari ClaimData.SpreadingClaim  (§6)
│    ├─ SpreadingQuotaShare[]     ═ salinan beku dari ClaimData.SpreadingBreakQS (§6)
│    ├─ CNPSpreadLoss[]                        → GISFW-Data-SpreadingRisk
│    │    └─ CNPFlagOuts
│    └─ ListClaimAcceptation[]                 → GISFW-Data-TreatyInTotal
│         └─ CNPFlagOuts
│
├─ SpreadingClaim[]                            → GISFW-Data-SpreadingRisk            n=41
├─ CNPSpreadLoss[]                             → GISFW-Data-SpreadingRisk            n=37
│    └─ RetroList[]                            → GISFW-Int-RETROCESSIONLIFE          n=5
│         └─ TREATYTYPENAME · COMMISION · COMMISION_AMOUNT · OVR_COMM · PREMIUM_SPREADED_NET
├─ SpreadingBreakQS[]                          → GISFW-Data-SpreadingRisk            n=28
├─ SpreadingAdjustment[]                       → GISFW-Data-SpreadingRisk            n=20
├─ SpreadingAdjustmentQS[]                     → GISFW-Data-SpreadingRisk            n=18
│
├─ ListTotalEstimation[]                       → GISFW-Data-TreatyInTotal            n=37
├─ ReinstatementList[]                         → GISFW-Data-TreatyInTotal            n=16
├─ TotalInterestInsured[]                      → GISFW-Data-TreatyInTotal            n=14
│
├─ InterestList[]                              → GCNMFW-Data-ObjectItem              n=39
│    └─ Currency · CurrencyID · KursObjectItem · TSIPerObject
├─ ObjectList[]                                → GCNMFW-Data-ObjectItem              n=7
│    └─ ObjectItemList[]                       → GCNMFW-Data-ObjectItem
│         ├─ ObjectName · TSIPerObject · Currency · CurrencyID · TPLAmount · Value
│         └─ Adjustment{}                                                     ◄── KEDALAMAN 4
│              └─ CurrencyID
│
├─ EstimationList[]                            → GCNMFW-Data-Estimasi                n=24
├─ Attachment[]                                → (tidak dideklarasikan)              n=8
│    └─ IMAGEID · URLPUBLIC · PNOTE · AttachStream
├─ ClaimComitee[]                              → GCNMFW-Data-Comitee                 n=4
│    └─ IDKomite · KomiteEmail
└─ SuggestList[]                               → GISFW-Data-OfferFacIn-SuggestList   n=5
     └─ CommentSuggest · PICSuggest · DateSuggest · IsCedingConfirm
```

### 4.1 Skalar langsung di `.ClaimData` (91)

**Identitas & nomor** (13)
`NoClaim` · `ClaimNo` · `PolicyNo` · `IDMaster` · `IDMasterTONP` · `InsuredName` · `NoPla` · `PlaNoCeding` · `PlaNoSOB` · `DLANoCeding` · `CNPClmNoCedant` · `CASEDB` · `YearofAccount`

**Tanggal & periode** (8)
`DateOfLoss` · `ReportDate` · `DateReceived` · `StartDateTreaty` · `StartDateTreaty1` · `EndDateTreaty` · `TreatyYear` · `PeriodPolicyTBA`

**Konteks layer aktif** (6) — layer yang sedang dikerjakan, terpisah dari baris-baris di `SpreadingRisk[]`
`TreatyName` · `Layer` · `LayerType` · `LayerPart` · `LayerPartType` · `TotalLayerList`

**Sebab & lokasi kerugian** (14)
`CauseOfLoss` · `CauseOfLossID` · `Occupation` · `Location` · `ReportAddress` · `Province` · `ProvinceID` · `City` · `CityID` · `District` · `DistrictID` · `RW` · `RWID` · `PostalCode`

**Pelapor & konteks laporan** (11)
`ReporterName` · `ReporterTelp` · `ReporterStatus` · `ReportType` · `ReportDescription` · `InsuredRelationship` · `InsuredRelationshipOthers` · `InsuredInterest` · `CNPCircumtances` · `CNPSupportDoc` · `CNPReinsuranceSlip`

**Parameter perhitungan** (10) — yang menentukan hasil alokasi
`ShareCeding` · `DeductibleType` · `DeductibleValue` · `CurrencyDeductible` · `TypeDeductible` · `TSIDeductible` · `CNPDeducMinMax` · `Amount` · `FormType` · `TotalSumInsuredIDR`

**TPL (tanggung jawab pihak ketiga)** (6)
`IsTPL` · `TPLAmount` · `CurrencyTPL` · `PctTPL` · `TPLType` · `TPLFormat`

**Katastrofe** (5)
`StsKatastrofe` · `KatastrofeID` · `KatastrofeNote` · `NonKatastrofeType` · `EditCatastrope`

**Penilai & pembayaran** (6)
`AppointedADJ` · `AppointedADJID` · `ConsultantID` · `ConsultantName` · `Payable` · `PayableTo`

**Rekap terhitung** (4) — turunan, bukan masukan
`TotalListClaimAmount` · `TotalListClaimAmountIDR` · `TotalEstimasiIDR` · `TotalGrossEstimateIDR`

**Bendera proses & sisa** (8)
`IsPLA` · `IsAnalisTransfer` · `IsFlagError` · `Remark_Close` · `Quater` · `YearofQuartal` · `ListClaimAcceptation` · `SpreadingQuotaShare`

> Dua nama terakhir tampak seperti Page List dan memang berkelas list, tapi **di modul ini tidak pernah ada satu pun anaknya yang dirujuk** — hanya seluruh listnya yang disalin sekaligus. Bentuk isinya baru terlihat saat list yang sama muncul di dalam `AdjustmentList[]`.

### 4.2 `SpreadingRisk[]` — baris alokasi per layer (29 field + `RetroList[]`)

Simpul paling hidup di seluruh pohon (n=215, 20 berkas). **Satu baris = satu layer**, dan baris Retensi Cedant ikut di sini bertanda `TreatyName="UR"` — bukan entitas terpisah (BLUEPRINT §2.4, ADR-0010).

| Kelompok | Field |
|---|---|
| Identitas layer | `TreatyName` · `TreatyType` · `Layer` · `LayerType` · `LayerPart` · `LayerPartType` · `CurrLayerOri` |
| Nilai klaim | `ClaimEstimation` (n=33, tertinggi) · `ClaimAmountAdjust` · `ClaimAmountIDR` · `AdjClaimValue` · `TotalClaim` · `TotalSpread` · `ClaimSpreaded` · `ClaimPercentage` |
| Biaya | `AdjusterFee` · `Salvage` · `CNPOthersFee` |
| Parameter layer | `CNPLimit` · `CNPLimitFull` · `CNPMDP` · `CNPPctReinstate` · `CNPProrateClaim` |
| Mata uang | `Currency` · `CurrencyID` · `Kurs` · `KursIDR` |
| Lain | `IsEditClaim` (ditulis, tidak pernah dibaca — BLUEPRINT §2.5) · `AdjustmentList` |

### 4.3 `AdjustmentList[]` — unit transaksi pembayaran (47 skalar + 10 kontainer)

Satu baris = satu Adjustment = satu case Komite.

| Kelompok | Field |
|---|---|
| Jenis & keadaan | `PaymentType` · `Type` · `AcceptanceStatus` · `StatusKasir` · `StatusKonversi` · `DirectToKasir` |
| Nomor & tanggal | `AcceptedNo` · `AcceptedDate` · `KomiteNo` · `DLANoCeding` |
| Nilai | `TotalEstimasiValue` · `NetClaim` · `TotalClaimRNM` · `PersenRNM` |
| Komite | `IsKomite` · `TotalKomite` · `IsSubjectivity` · `SubjectivityNote` |
| Usulan penutupan | `IsProposeClose` · `IsPropReserved` |
| Konteks layer | `XOLID` · `CNPIndexInterim` · `CNPFlagXOL` · `CNPFlagReinstate` · `IndexObject` — pasangannya `XOL` hanya terbaca di page `Adjustment` (§5), tidak di `AdjustmentList[]` |
| Bendera nomor akseptasi | `CNPAccNoAdjustF` · `CNPAccNoOtherF` · `CNPAccNoSalvage` · `CNPAccNoReinstate` |
| Rekening | `Payable` · `PayableTo` · `NameOfBank` · `BranchOfBank` · `NoAccount` · `IDOfBank` · `SwiftCode` · `Currency` — plus varian `…2` lengkap |
| Risiko individual | `IndividualRiskType` · `IndividualRiskPercentage` · `IndividualRiskValue` |
| Catatan | `RemarksDLA` |

`IndexObject` adalah **indeks balik** ke posisi baris ini di dalam `AdjustmentList` — begitulah case Komite menemukan induknya kembali saat menulis hasil approval (BLUEPRINT §8.4).

### 4.4 `CNPLayerList[].CNPCurrencyList[]` — simpul terdalam (21 field)

Ini **rincian per layer per mata uang untuk satu Adjustment** — tingkat terdalam yang ada di sistem lama.

| Kelompok | Field |
|---|---|
| Kunci | `Currency` · `CurrencyID` · `KursIDR` · `PersenRNM` |
| Parameter layer | `CNPLimit` · `CNPMindep` · `CNPPctReinstate` |
| Nilai bruto | `GrossValue` · `GrossAdjustment` |
| Komponen | `AdjusterFeeValue` / `AdjusterFeeRNM` · `SalvageValue` / `SalvageRNM` · `CNPOthersFee` / `CNPOthersFeeRNM` · `CNPReinstatement` / `CNPReinstatementRNM` |
| Total | `TotalXOL` · `TotalXOLGross` · `TotalXOLRNM` · `TotalXOLLayerRNM` |

Polanya konsisten: **setiap komponen biaya punya sepasang field** — nilai penuh dan bagian reasuradur (`…RNM`, dari `TreatyInMaster.RNMShare`). Pasangan ini sejalan dengan ADR-0007 soal nilai uang berpasangan.

### 4.5 `ReceiverClaim[]` — selalu satu baris

Page List, tapi **di seluruh ekspor hanya pernah dirujuk sebagai `ReceiverClaim(1)`** — tidak ada satu pun perulangan atau `<APPEND>`. Isinya sepasang rekening dalam satu baris (`NoAccount` + `NoAccount2`, `NameOfBank` + `NameOfBank2`, dst.), bukan dua baris.

Artinya: bentuknya list, pemakaiannya tunggal, dan "rekening kedua" diwakili sufiks kolom — bukan baris kedua.

---

## 5. Cabang sejajar: page `Adjustment`, dan kenapa `CNPLayerList` ada dua kali

Ada satu hal yang harus dibaca dengan teliti karena mudah salah kutip.

**Yang terbaca langsung di XML:**
- `pyWorkPage.Adjustment.CNPLayerList(...).CNPCurrencyList(...).<field>` — **ada**, lengkap, 21 field.
- `pyWorkPage.ClaimData.AdjustmentList(...).CNPLayerList(...)` — **tidak pernah ditulis** dalam bentuk jalur penuh di 279 berkas ini.

**Yang tetap bisa dipastikan:** `Adjustment` dan `AdjustmentList` **berkelas sama** (`ASM-FW-GCNMFW-Data-Adjustment`, dikonfirmasi `pyStepsClassName` pada sepuluh langkah Activity), dan `Adjustment` adalah halaman kerja yang isinya masuk ke `AdjustmentList`. Karena itu `CNPLayerList[]` → `CNPCurrencyList[]` **berlaku juga di dalam setiap baris `AdjustmentList`** — di CSV ditandai `BUKTI = identitas-kelas`, bukan `jalur-penuh`.

Kelas `Data-Adjustment` **bersarang di dalam dirinya sendiri tiga tingkat**: `Adjustment` → `CNPLayerList` → `CNPCurrencyList`. Struktur rekursif dengan arti berbeda di tiap tingkat — transaksi, layer, mata uang. Satu tabel untuk tiga tingkat ini akan salah.

```
Adjustment{}                     → GCNMFW-Data-Adjustment            96 simpul
├─ PaymentType · Type · XOL · XOLID · PersenRNM · IndexObject · Occupation
├─ CNPIndexInterim · CNPFlagReinstate
├─ AcceptedNo · AcceptedDate · PayableTo · NoAccount
├─ StatusKasir · pyStatusMessage
├─ StatusServiceKasir{}          → (tidak dideklarasikan)
│    └─ ReponseCode · ResponseMsg        ← salah eja ada di sumbernya
├─ ComiteeClaim[]                → GCNMFW-Data-Comitee
├─ SpreadingRisk[] · SpreadingAdjustment[] · SpreadingQuotaShare[]
├─ CNPSpreadLoss[] · AlokasiXOLPaid[] · ListTotalEstimation[] · ListClaimAcceptation[]
└─ CNPLayerList[]                → GCNMFW-Data-Adjustment      (kelas sama, tingkat berbeda)
     ├─ XOL · XOLID
     └─ CNPCurrencyList[]        → GCNMFW-Data-Adjustment      (kelas sama lagi)
          └─ 21 field
```

Empat puluh delapan simpul di cabang ini **hanya terbaca lewat alias step-page `Primary`** — di CSV bertanda `BUKTI = alias-primary`. Penyamaannya dengan page `Adjustment` bukan tebakan: kedelapan Activity yang memakai `Primary.*` — `AdjClaimCNP_Act`, `CreateChildKomiteCNP_Act`, `GenerateCACNP_Act`, `HitServiceToKasir_Act`, `ProteksiSendKomiteCNP_Act`, `SaveAdjustmentToOSAksep_Act_Tes`, `SaveCNPLayerList_Act`, `SetAccoutNo_Act` — semuanya **applies-to `ASM-FW-GCNMFW-Data-Adjustment`**, jadi `Primary` di situ memang halaman itu sendiri. Alatnya menerapkan syarat applies-to ini; `Primary` di rule berkelas lain tidak ikut disamakan.

`StatusServiceKasir{}` dan `pyStatusMessage` adalah jejak integrasi kasir yang tidak muncul di tingkat `ClaimData` — hanya ada di sini.

---

## 6. Salinan beku di dalam `AdjustmentList`

`AddAkseptasiCNP_Act` menulis satu baris Adjustment baru dengan `AdjustmentList(<APPEND>)`, dan **menyalin utuh empat Page List tingkat klaim ke dalam baris itu**:

| Ditulis ke | Disalin dari | Kelas |
|---|---|---|
| `AdjustmentList[].LossAllocation` | `ClaimData.SpreadingRisk` | GISFW-Data-SpreadingRisk |
| `AdjustmentList[].SpreadingRisk` | `ClaimData.SpreadingRisk` | GISFW-Data-SpreadingRisk |
| `AdjustmentList[].SpreadingAdjustment` | `ClaimData.SpreadingClaim` | GISFW-Data-SpreadingRisk |
| `AdjustmentList[].SpreadingQuotaShare` | `ClaimData.SpreadingBreakQS` | GISFW-Data-SpreadingRisk |
| `AdjustmentList[].CurencyAdjustment` | `Kurs.pxResults` (hasil query kurs) | GCNMFW-Data-Adjustment |

Perhatikan barisan kedua dan ketiga: **nama tujuan tidak sama dengan nama sumber.** `SpreadingClaim` di tingkat klaim menjadi `SpreadingAdjustment` di dalam Adjustment, dan `SpreadingBreakQS` menjadi `SpreadingQuotaShare`. Nama yang sama (`SpreadingAdjustment`) karena itu berarti **dua hal berbeda** tergantung induknya. Ini jebakan pemetaan yang nyata, bukan kerapian penamaan yang kurang.

Maknanya untuk rancangan: setiap Adjustment membawa **potret alokasi pada saat ia dibuat**. Perhitungan berikutnya di tingkat klaim tidak mengubah potret yang sudah beku. Perilaku ini yang harus dipertahankan, apa pun bentuk tabelnya.

---

## 7. Turunan di luar clipboard

Pohon di atas adalah bentuk **di memori**. Saat keluar dari Pega, bentuknya berubah tiga kali.

### 7.1 `OS_AKSEPTASI_KLAIM.data_json` — akseptasi outstanding

Ditulis lewat `InputParamOs` (kelas `ASM-FW-GCNMFW-Data-osAkseptasi`) oleh `SaveDataToOSAksep_Act`, `SaveToOS`, `GetSelisihActual_Act`, dan `SaveRejectOSKomiteCNP` (modul Komite). **17 field**, jauh lebih sempit dari sumbernya:

| Field JSON | Diisi dari |
|---|---|
| `NoClaim` | `ClaimData.NoClaim` |
| `IDMasterTreaty` | `TreatyInMaster.ID` |
| `CauseOfLoss` / `CauseOfLossID` | `ClaimData.CauseOfLoss` / `…ID` |
| `TypeLoss` / `TypeLossID` | `SpreadingRisk[].TreatyName` / `.TreatyType` — **nama layer, bukan jenis kerugian** |
| `Currency` / `CurrencyID` | `SpreadingRisk[].Currency` / `…ID` |
| `KursValue` | `SpreadingRisk[].Kurs` |
| `GrossValue` | `SpreadingRisk[].ClaimEstimation + .AdjClaimValue` |
| `Value` | `SpreadingRisk[].ClaimSpreaded` |
| `Adjusterfee` · `Salvage` · `CNPOthersFee` | field senama di `SpreadingRisk[]` |
| `PersenRNM` | `TreatyInMaster.RNMShare` |
| `Type` · `EstimationDate` | konstanta (`0`, string kosong) |

Dua hal yang harus dicatat:
- **`TypeLoss` berisi nama layer**, bukan jenis kerugian. Namanya menyesatkan; isinya yang benar.
- Penulisannya **selisih, bukan nilai penuh**: `InputParamOs.X = X − OutOSAcc.pxResults(1).X`. Baris baru berisi delta terhadap yang sudah tercatat. Pembacaan balik (`GetDataCNPOS`) memang `sum(...)` dengan `STS_REJECT = 0`.

### 7.2 JSON klaim penuh — `InsertJsonClaimTreaty_act`

Empat field saja, dan salah satunya memuat seluruh sisanya:

| Kolom | Isi |
|---|---|
| `TempPNC.BUSINESS_CODE` | `pyWorkPage.pzInsKey` |
| `TempPNC.POLICY_NO` | `ClaimData.PolicyData.PolicyNo` |
| `TempPNC.No_Klaim` | `ClaimData.NoClaim` |
| `TempPNC.TSI` | `@GCNM.GetPageJSONString()` — **seluruh halaman diserialkan jadi satu string JSON** |

`GetPageJSONString()` adalah fungsi Java kustom yang **tidak ada di ekspor**. Bentuk JSON hasilnya karena itu tidak bisa dipastikan dari berkas-berkas ini; yang pasti hanya bahwa isinya adalah serialisasi `pyWorkPage.ClaimData`. Ini gap nyata — dicatat, tidak ditebak.

### 7.3 Connect-REST

| Rule | Endpoint | Yang dikirim |
|---|---|---|
| `InsertClaimOutstanding_NP` | `…/NusareClaim/insertClaimAccept` | `.pzInsKey` · `.ClaimData.PolicyData.PolicyNo` |
| `insertClaimFinalOrClosed_NP` | `…/NusareClaim/insertClaimFinalOrClosed` | `.pzInsKey` · `.ClaimData.PolicyData.PolicyNo` |
| `SendAcceptationToKasir` | (kasir) | `ParamKasir.CARI1` |

Muatannya sangat tipis — **dua field** untuk dua endpoint klaim, satu untuk kasir. Sistem penerima menarik sendiri sisanya dari basis data bersama lewat `pzInsKey`. Kopling ini lewat data, bukan lewat muatan pesan; sejalan dengan temuan BLUEPRINT §8.4b.

---

## 8. Batas bukti

Hal-hal yang **tidak** bisa dijawab berkas ini, dan tidak boleh dikarang dari isinya:

1. **Tipe dan panjang properti.** Tidak ada `Rule-Obj-Property` di ekspor. Kolom `JENIS` adalah dugaan (§1.3).
2. **Bentuk JSON `TempPNC.TSI`.** Bergantung pada fungsi Java `GCNM.GetPageJSONString()` yang tidak diekspor (§7.2).
3. **Kelas untuk tujuh simpul** di §3. Isinya terbaca, namanya tidak ada.
4. **Nilai enum.** `AcceptanceStatus`, `PaymentType`, `StatusKasir`, `StatusKonversi`, `IndividualRiskType` dipakai sebagai kode tanpa daftar nilainya di mana pun. Sebagiannya terbaca dari perbandingan di rule When — lihat BLUEPRINT §3.
5. **Wajib/opsional.** Tidak ada satu pun aturan validasi tingkat properti di ekspor; validasi yang ada hidup di dalam Activity sebagai `Property-Set` bersyarat.
6. **Kardinalitas sebenarnya.** `ReceiverClaim[]` berbentuk list tapi selalu satu baris (§4.5); kebalikannya juga mungkin ada dan tidak terlihat dari struktur.

Enam butir ini sejalan dengan keputusan ruang lingkup 17 September 2026: **ambil perilaku yang terlihat di XML sebagai sumber kebenaran, tandai asumsinya, jangan berhenti menunggu konfirmasi.**
