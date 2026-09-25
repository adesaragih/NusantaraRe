# FINDING-007 — Dua rumus premi reinstatement bekerja atas nilai yang berbeda

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); `pengetahuan/DDL_Script_ClaimNonProp.xls` versi 2026-09-18 10:36 (48 objek); dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).

**Jenis**: laporan kondisi sistem lama
**Status**: TERBUKTI DARI XML
**Tanggal**: 18 September 2026

> **Satu pola, tiga temuan.** Ini kelas cacat yang sama dengan `CheckDateDOL_Act` (FINDING-005) dan `.IsEditClaim` (FINDING-004): nilai dibaca dalam keadaan yang tidak dimaksudkan, dan tidak ada yang menjaganya.

---

## 1. Dua rumus

| Berkas & baris | Langkah | Ekspresi |
|---|---|---|
| `Activity\CountReinstatement_Act.xml` **624** | — | `((((.ClaimEstimation + .AdjusterFee) − (.Salvage×100/RNMShare)) / .CNPLimit) × .CNPMDP) × (.CNPPctReinstate/100)` |
| `Activity\AdjClaimCNP_Act.xml` **3109** | `RH_1.pySteps(11).pySteps(1)` | `@if(.TreatyType=="UR", 0, @divide(.TotalClaim, .CNPLimit, 20) × .CNPMDP × @divide(.CNPPctReinstate, 100, 20))` |

Keduanya setara **hanya bila** `.TotalClaim` bernilai sama dengan `ClaimEstimation + AdjusterFee − Salvage×100/RNMShare` pada saat baris 3109 dijalankan.

## 2. Tidak setara — dan sebabnya lebih tajam dari dugaan

`MEMORI_PEMAHAMAN.MD` §6.2 Langkah 7 menyatakan `TotalClaim = ClaimEstimation`. **Penetapan seperti itu tidak ada di XML**; yang ada lebih rumit.

### 2.1 Di `CountLossAllocation_act`, keduanya memang menjadi sama

| Baris | Properti | Nilai |
|---|---|---|
| 9258 | `SpreadingRisk(<LAST>).ClaimEstimation` | `@if(Local.ClaimValue − Local.TotalAllocation < Local.LayerLimit, …)` |
| 9279 | `SpreadingRisk(<LAST>).TotalClaim` | `@if(Local.ClaimValue − Local.TotalAllocation < Local.LayerLimit, …)` |

Bentuk yang sama, hasil alokasi yang sama. Jadi saat `CountLossAllocation_act` selesai, **keduanya bernilai sama** — memori benar untuk titik itu.

### 2.2 Tetapi `AdjClaimCNP_Act` menimpanya **sebelum** baris 3109 dibaca

Urutan langkah di dalam rule yang sama, dari `pyStepPageReference`:

| Baris | Langkah | Yang terjadi |
|---|---|---|
| 2348 | `RH_1.pySteps(8)` | `.SpreadingRisk(IdxLastLayer).TotalClaim` **ditimpa** |
| 2369 | `RH_1.pySteps(8)` | `.SpreadingRisk(IdxLastLayer).TotalClaim` **ditimpa lagi** — `@if(Local.TotalClaim==0, Local.TotalValue − …, …)` |
| **3109** | **`RH_1.pySteps(11).pySteps(1)`** | **rumus reinstatement membaca `.TotalClaim`** |
| 6441 | `RH_1.pySteps(15).pySteps(4).pySteps(1)` | `.TotalClaim = .ClaimSpreaded + .AdjusterFee + .CNPOthersFee − .Salvage` |
| 7834 | `RH_1.pySteps(15).pySteps(5).pySteps(5)` | idem |

**Langkah 8 menimpa `TotalClaim` sebelum langkah 11.1 membacanya.** Nilai yang masuk rumus reinstatement bukan hasil alokasi, melainkan hasil langkah 8 — yang bersumber dari `Local.TotalClaim`, yaitu akumulasi `.ClaimAmountAdjust` (baris 1561).

Dan penetapan yang memang memuat `AdjusterFee` serta `Salvage` — baris 6441 dan 7834 — berjalan **sesudah** 3109, jadi tidak terbaca olehnya.

> **Kedua rumus bekerja atas nilai yang berbeda.** Bukan karena salah satunya menghilangkan komponen dengan sengaja, melainkan karena `TotalClaim` ditimpa tiga kali di dalam satu rule, dan rumus reinstatement membacanya di antara dua penimpaan.

## 3. Yang menguatkan bahwa keduanya *dimaksudkan* sama

`Activity\SetActualPremium_ACT.xml` menyetelnya secara eksplisit:

| Baris | Penetapan |
|---|---|
| 702 | `SpreadingRisk(<LAST>).ClaimEstimation = SpreadingRisk(<LAST>).TotalClaim` |
| 951 | `.ClaimEstimation = .TotalClaim` |

Jadi di tempat lain kesamaan itu **ditegakkan lewat penetapan**, bukan berlaku dengan sendirinya. Itu memperkuat pembacaan memori tentang maksudnya, sekaligus memperlihatkan bahwa maksud itu tidak dijaga di semua jalur.

## 4. Batas klaim

1. **Urutan langkah disimpulkan dari `pyStepPageReference` dan urutan dokumen.** Untuk langkah tingkat atas di dalam satu activity, keduanya sejalan. Percabangan `pyStepsPreCondParamsWhen` dapat membuat sebagian langkah dilewati, dan itu belum ditelusuri satu per satu.
2. **Berapa besar selisih angkanya belum diukur** — itu pertanyaan data, bukan pertanyaan kode.
3. Rule mana yang benar secara bisnis bukan kesimpulan dokumen ini. Acuan yang dipilih adalah `CountReinstatement_Act`.

## 5. Pengamatan terpisah, tidak ditarik lebih jauh

`Activity\CopyOldataCurr_act.xml` baris 1862 dan `Activity\CountTotalInsterest_Act.xml` baris 4769 memuat:

```
.ClaimEstimation = .ClaimSpreaded * .PremiumSpreaded
```

Nilai klaim dikalikan nilai premi. Dicatat apa adanya; tidak ditelusuri lebih jauh.
