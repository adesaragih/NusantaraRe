# STRUKTUR TABEL — Claim Non Prop

> Dibuat 09-10-2026 bersama migrasi `600`–`611`. Keputusan yang mengikat: **OQ-CNP-06 "pola Claim Prop + tabel XoL"**
> (work owner 09-10-2026) dan izin work owner 09-10-2026 menyunting STRUKTUR Claim Life / Claim Prop untuk kolom yang
> ditambahkan ke tabel bersama. Katalog properti ↔ kolom: `backend/models/katalog.go` (satu-satunya pemetaan).

## Ringkasan susunan

- **Kasus:** `T_WORK_CLAIM` (tabel bersama Claim Life) `LINI = 'NONPROP'`, ID `CLMNP-nnnnnn` (klaim) dan
  `KMTNP-nnnnnn` (kasus komite, awalan Pega), `TAHAP` = nama FlowAction (`OutstandingClaim` / `InputAcceptation`) atau
  `KomiteTreaty_Flow`. Tidak ada kolom baru.
- **Kepala klaim:** `T_GENERAL_CLAIM` (shared PK dengan `T_WORK_CLAIM`). Kolom Claim Life / Claim Prop dipakai ulang bila
  maknanya sama; 25 kolom khas Non Prop ditambah migrasi `600` — baris kolomnya di bab `## T_GENERAL_CLAIM`
  `modul/claimlife/docs/STRUKTUR-TABEL-CLAIM-LIFE.md` (satu tabel, satu dokumen).
- **Tabel anak Claim Prop dipakai ulang** dengan kolom tambahan nullable (migrasi `601`–`607`, baris kolomnya di bab
  masing-masing `modul/claimprop/docs/STRUKTUR-TABEL-CLAIM-PROP.md`):

  | Tabel | Daftar halaman | Migrasi | Kolom tambahan |
  | --- | --- | --- | --- |
  | T_CLAIM_INTEREST | `ClaimData.InterestList` (+ TPL, Section InputDtlInterest) | `601` | 13 |
  | T_CLAIM_CLAIM_AMOUNT | `ClaimData.ListClaimAmount` | `602` | 9 |
  | T_CLAIM_SPREADING | `ClaimData.SpreadingClaim` | `603` | 5 |
  | T_CLAIM_BREAK_QS | `ClaimData.SpreadingBreakQS` | `604` | 5 |
  | T_CLAIM_ADJUSTMENT | `ClaimData.AdjustmentList` (akseptasi) | `605` | 8 |
  | T_CLAIM_ADJ_SPREADING | `AdjustmentList(n).SpreadingAdjustment` (Spreading In) | `606` | 5 |
  | T_CLAIM_ADJ_QUOTA_SHARE | `AdjustmentList(n).SpreadingQuotaShare` (Spreading Out) | `607` | 6 |

- **Tabel baru khas XoL** (bab pengikat di bawah): `T_CLAIM_NP_LOSS_ALLOC` (Loss Allocation + To XOL),
  `T_CLAIM_NP_XOL_ALLOC` (XOL Allocation per layer; `JENIS` membedakan alokasi berjalan, alokasi lama, dan Previously
  Calculated - satu tabel, tiga daftar halaman yang bentuknya sama persis), `T_CLAIM_NP_CLAIM_ACCEPT` (Claim
  Acceptation per mata uang). Tidak ada tabel lain: Previously Calculated (`AlokasiXOLPaid`) memakai
  `T_CLAIM_NP_XOL_ALLOC` `JENIS = 'DIBAYAR'` karena kolomnya identik dengan XOL Allocation.
- **ID baris** tabel anak / baru: `SEQ_T_CLAIM` (claimprop `533`), urutan bersama tabel `T_CLAIM_*`.
- **Riwayat:** `T_VIEW_SUGGEST` `CLAIM_ID` (claimprop `532`), hanya bertambah.
- **Tabel warisan yang ditulis:** `OS_AKSEPTASI_KLAIM` (INSERT 8 kolom procedure `PEGA_JSON_OS_AKSEP_KLAIMTNP`,
  `DATA_JSON` diisi), `JSON_KLAIM` (tanpa `DATA_JSON`), `CATASTROPHE`, `T_LOG_SERVICE_RNM` (outbox, hanya produksi),
  tangga komite `T_GENERAL_KOMITE` / `T_KOMITE_KOMITELIST` (hanya `backend/repository/komite.go`). Roster `EMAILKOMITE`
  NONPROP diganti di tempat (migrasi `611`).

### Medan TURUNAN - tanpa kolom

Dihitung ulang sesudah setiap aksi / saat halaman dimuat (`models.HitungTurunan`, PARITAS): `TotalInterestInsured`,
`TotalSumInsuredIDR`, `ListTotalEstimation` (Summary / Total XOL Allocation), `SpreadingAdjustment(QS)` tingkat
klaim, `AdjustmentList(n).CurencyAdjustment / .ComiteeClaim / .TotalKomite / .CommentLOD / .FlagCurrency`, rekening
Spreading In (`SetAccoutNo_Act`). `CNPLayerList` disusun saat baris OS / komite ditulis, tanpa tabel.

### Relasi

```
T_WORK_CLAIM (LINI NONPROP) ── shared PK ── T_GENERAL_CLAIM
  ├─ T_CLAIM_INTEREST / T_CLAIM_CLAIM_AMOUNT / T_CLAIM_SPREADING / T_CLAIM_BREAK_QS     (CLAIM_ID, CASCADE)
  ├─ T_CLAIM_NP_LOSS_ALLOC / T_CLAIM_NP_XOL_ALLOC (JENIS ALOKASI)                       (CLAIM_ID, CASCADE)
  ├─ T_VIEW_SUGGEST                                                                      (CLAIM_ID, CASCADE)
  └─ T_CLAIM_ADJUSTMENT                                                                  (CLAIM_ID, CASCADE)
       ├─ T_CLAIM_NP_CLAIM_ACCEPT / T_CLAIM_NP_LOSS_ALLOC                                (ADJUSTMENT_ID, CASCADE)
       ├─ T_CLAIM_NP_XOL_ALLOC (JENIS ALOKASI / LAMA / DIBAYAR)                          (ADJUSTMENT_ID, CASCADE)
       ├─ T_CLAIM_ADJ_SPREADING / T_CLAIM_ADJ_QUOTA_SHARE                                (ADJUSTMENT_ID, CASCADE)
       └─ T_GENERAL_KOMITE.ADJUSTMENT_ID  (kasus komite KMTNP-; baris berkomite tidak pernah dihapus aplikasi)
```

`T_CLAIM_NP_LOSS_ALLOC` dan `T_CLAIM_NP_XOL_ALLOC` punya dua induk: tepat satu terisi (CHECK).

## Lampiran pengikat penjaga — kolom DDL tabel milik Claim Non Prop (09-10-2026)

Bab `## T_…` di bawah **mengikat** (`inti/backend/penjaga`: `TestKolomDDLCocokDenganStruktur`,
`TestGolonganTipeDDLCocokDenganStruktur`) dan dibangkitkan dari katalog yang sama dengan migrasinya. Kolom sumber =
properti anggota daftar halaman Pega.

## T_CLAIM_NP_LOSS_ALLOC

Migrasi `608`. Grid "Loss Allocation" + To XOL: `ClaimData.CNPSpreadLoss` (klaim, `CLAIM_ID`) dan `AdjustmentList(n).CNPSpreadLoss` (akseptasi, `ADJUSTMENT_ID`).

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `608` — SEQ_T_CLAIM |
| `CLAIM_ID` | teks | ya | FK | migrasi `608` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE; tepat satu induk (CHECK) |
| `ADJUSTMENT_ID` | teks | ya | FK | migrasi `608` — → `T_CLAIM_ADJUSTMENT.ID`, ON DELETE CASCADE; tepat satu induk (CHECK) |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `608` — urutan baris daftar halaman |
| `CURRENCY_ID` | teks | ya |  | `.CurrencyID` |
| `CURRENCY_NAME` | teks | ya |  | `.Currency` |
| `TREATY_NAME` | teks | ya |  | `.TreatyName` |
| `CLAIM_PERCENTAGE` | angka desimal | ya |  | `.ClaimPercentage` |
| `CLAIM_AMOUNT_ADJUST` | angka desimal | ya |  | `.ClaimAmountAdjust` |
| `ADJUSTER_FEE` | angka desimal | ya |  | `.AdjusterFee` |
| `SALVAGE` | angka desimal | ya |  | `.Salvage` |
| `OTHERS_FEE` | angka desimal | ya |  | `.CNPOthersFee` |
| `TO_XOL` | teks | ya |  | `.CNPFlagXOL` |
| `IS_LOCKED` | teks | ya |  | `.CNPFlagOuts` |

**Unik:** `(CLAIM_ID, ADJUSTMENT_ID, NOURUT)`.

## T_CLAIM_NP_XOL_ALLOC

Migrasi `609`. Grid "XOL Allocation" per layer: `JENIS` ALOKASI = `ClaimData.SpreadingRisk` / `AdjustmentList(n).SpreadingRisk`; LAMA = `AdjustmentList(n).LossAllocation` (salinan saat akseptasi lahir, View Old Allocation); DIBAYAR = `AdjustmentList(n).AlokasiXOLPaid` (grid "Previously Calculated").

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `609` — SEQ_T_CLAIM |
| `CLAIM_ID` | teks | ya | FK | migrasi `609` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE; tepat satu induk (CHECK) |
| `ADJUSTMENT_ID` | teks | ya | FK | migrasi `609` — → `T_CLAIM_ADJUSTMENT.ID`, ON DELETE CASCADE; tepat satu induk (CHECK) |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `609` — urutan baris daftar halaman |
| `JENIS` | teks | tidak |  | migrasi `609` — `ALOKASI` / `LAMA` / `DIBAYAR` (CHECK), satu daftar halaman per nilai |
| `TREATY_TYPE_ID` | teks | ya |  | `.TreatyType` |
| `TREATY_NAME` | teks | ya |  | `.TreatyName` |
| `CURRENCY_ID` | teks | ya |  | `.CurrencyID` |
| `CURRENCY_NAME` | teks | ya |  | `.Currency` |
| `KURS` | angka desimal | ya |  | `.Kurs` |
| `KURS_IDR` | angka desimal | ya |  | `.KursIDR` |
| `CLAIM_ESTIMATION` | angka desimal | ya |  | `.ClaimEstimation` |
| `CLAIM_AMOUNT_ADJUST` | angka desimal | ya |  | `.ClaimAmountAdjust` |
| `ADJ_CLAIM_VALUE` | angka desimal | ya |  | `.AdjClaimValue` |
| `TOTAL_CLAIM` | angka desimal | ya |  | `.TotalClaim` |
| `CLAIM_PERCENTAGE` | angka desimal | ya |  | `.ClaimPercentage` |
| `CLAIM_SPREADED` | angka desimal | ya |  | `.ClaimSpreaded` |
| `ADJUSTER_FEE` | angka desimal | ya |  | `.AdjusterFee` |
| `SALVAGE` | angka desimal | ya |  | `.Salvage` |
| `OTHERS_FEE` | angka desimal | ya |  | `.CNPOthersFee` |
| `TOTAL_SPREAD` | angka desimal | ya |  | `.TotalSpread` |
| `CLAIM_AMOUNT_IDR` | angka desimal | ya |  | `.ClaimAmountIDR` |
| `PRORATE_PCT` | angka desimal | ya |  | `.CNPProrateClaim` |
| `LIMIT_VALUE` | angka desimal | ya |  | `.CNPLimit` |
| `LIMIT_FULL` | angka desimal | ya |  | `.CNPLimitFull` |
| `LAYER_CURRENCY` | teks | ya |  | `.CurrLayerOri` |
| `MDP_VALUE` | angka desimal | ya |  | `.CNPMDP` |
| `REINSTATE_PCT` | angka desimal | ya |  | `.CNPPctReinstate` |
| `LAYER` | teks | ya |  | `.Layer` |
| `LAYER_TYPE` | teks | ya |  | `.LayerType` |
| `LAYER_PART` | teks | ya |  | `.LayerPart` |
| `LAYER_PART_TYPE` | teks | ya |  | `.LayerPartType` |
| `IS_EDIT_CLAIM` | teks | ya |  | `.IsEditClaim` |
| `IS_LOCKED` | teks | ya |  | `.CNPFlagOuts` |
| `TOTAL_CLAIM_RNM` | angka desimal | ya |  | `.TotalClaimRNM` |
| `REINSTATEMENT` | angka desimal | ya |  | `.CNPReinstatement` |
| `REINSTATEMENT_RNM` | angka desimal | ya |  | `.CNPReinstatementRNM` |

**Unik:** `(CLAIM_ID, ADJUSTMENT_ID, JENIS, NOURUT)`.

## T_CLAIM_NP_CLAIM_ACCEPT

Migrasi `610`. Claim Acceptation per mata uang: `AdjustmentList(n).ListClaimAcceptation` (AddAkseptasiCNP_Act langkah 16).

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `610` — SEQ_T_CLAIM |
| `ADJUSTMENT_ID` | teks | tidak | FK | migrasi `610` — → `T_CLAIM_ADJUSTMENT.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `610` — urutan baris daftar halaman |
| `CURRENCY_ID` | teks | ya |  | `.CurrencyID` |
| `CURRENCY_NAME` | teks | ya |  | `.Currency` |
| `KURS` | angka desimal | ya |  | `.AltValue` |
| `VALUE` | angka desimal | ya |  | `.Value` |
| `VALUE_IDR` | angka desimal | ya |  | `.USD` |
| `NET_DEDUCTIBLE_VALUE` | angka desimal | ya |  | `.CNPDeductible` |
| `NOTE` | teks | ya |  | `.Note` |
| `TPL` | angka desimal | ya |  | `.TPL` |
| `ADJUSTER_FEE` | angka desimal | ya |  | `.AdjusterFee` |
| `SALVAGE` | angka desimal | ya |  | `.Salvage` |
| `OTHERS_FEE` | angka desimal | ya |  | `.CNPOthersFee` |
| `PROPORTION_PCT` | angka desimal | ya |  | `.PctProrateClaim` |
| `CLAIM_AMOUNT_CEDANT` | angka desimal | ya |  | `.ClaimAmountCedant` |
| `CLAIM_AMOUNT_ADJUST` | angka desimal | ya |  | `.ClaimAmountAdjust` |
| `TSI_VALUE` | angka desimal | ya |  | `.CNPTSI` |
| `IS_LOCKED` | teks | ya |  | `.CNPFlagOuts` |

**Unik:** `(ADJUSTMENT_ID, NOURUT)`.
