# Struktur Tabel — EDM Treaty In

Acuan bentuk tabel yang **dibuat** modul `edmtreatyin` (migrasi 360-363). Empat tabel **proyeksi selisih**
(spec-penyimpanan ID-6, ID-23..ID-27; diagram `Diagram-Skema-Tabel-NusantaraRe.xlsx` sheet *EDM Treaty In Prop*
J85 · R103 · R116 dan *EDM Treaty In NonProp* R139). Dicocokkan dengan katalog `backend/models/katalog_selisih.go`
oleh uji `backend/migrasi_test.go`.

Tipe ditulis sebagai kategori logis: teks · angka desimal · bilangan bulat · DATE. Uang dan persen **angka desimal**
`NUMBER(38,10)` — sama dengan tabel dasar NB 320-327 (RALAT `NUMBER(38,8)` spec-penyimpanan AC 58, keputusan work
owner 04-10-2026 untuk tabel dasar); tidak pernah float (ADR-0003).

Tabel generasi yang **dibaca dan ditulis, tidak dibuat** modul ini (`T_WORK_POLIS`, `T_GENERAL_POLIS_TREATY`,
`T_POLIS_*`) dideklarasikan di `MODUL.md`; bentuknya di `modul/nbtreatyin/docs/STRUKTUR-TABEL-NB-TREATY-IN.md`.

## T_POLIS_DIFFERENCE

Satu baris per generasi endorsemen (1:1, `POLIS_ID` unik): `PolicyTreatyIn.TreatyDifference` =
generasi ini dikurangi generasi yang ditunjuk `OLD_POLIS_ID` (`Activity/EDMTCalculateTreatyDifference` langkah 1).
Total* tidak disimpan — turunan baris spreading.

> **Membaca baris induk hasil rumus lama Pega** (keputusan work owner 07-10-2026: tanpa kolom penanda di tabel ini).
> Pada data migrasi, generasi ke-2 dan seterusnya dihitung Pega dengan varian berlapis (`EDMTCalculateTreatyDifference`
> blok `HasEDMNo`: baru − selisih generasi lama, dipilih bila `OldData.EDMNo` terisi). Baris induk itu dikenali tanpa
> kolom tambahan: `SUMBER = 'PEGA' AND PRODKE >= 2` — syarat yang sama dengan `RUMUS_BERLAPIS = 1` di tabel anaknya.

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | — |
| `POLIS_ID` | teks | tidak | UQ, FK T_GENERAL_POLIS_TREATY | kode | generasi endorsemen |
| `NOPOLIS` | teks | ya |  | kode | `PolicyTreatyIn.PolicyNo` |
| `PRODKE` | bilangan bulat | ya |  | cacah | `PolicyTreatyIn.ProdKe` |
| `EDM_NO` | teks | ya |  | kode | `PolicyTreatyIn.EDMNo` |
| `IDPEGA` | teks | ya |  | kode | `pyWorkPage.pzInsKey` (ID kasus) |
| `SUMBER` | teks | tidak | CK `PEGA` / `GO` | kode | `PEGA` = migrasi (beku), `GO` = dihitung sistem baru |
| `INSTALLMENT` | teks | ya |  | kode | `TreatyDifference.Installment` (disalin) |
| `GROSS_PREMIUM` | angka desimal | ya |  | uang | `TreatyDifference.GrossPremium` |
| `PREMI_OGP` | angka desimal | ya |  | uang | `TreatyDifference.PremiOgp` |
| `RI_COMM_OGP` | angka desimal | ya |  | persen | `TreatyDifference.RiCommOgp` (disalin) |
| `RESULT_OGP1` | angka desimal | ya |  | uang | `TreatyDifference.ResultOgp1` |
| `OVERIDDING_COMM_OGP` | angka desimal | ya |  | persen | `TreatyDifference.OveriddingCommOgp` (disalin) |
| `RESULT_OGP2` | angka desimal | ya |  | uang | `TreatyDifference.ResultOgp2` |
| `CLAIM` | angka desimal | ya |  | uang | `TreatyDifference.Claim` |
| `SALVAGE_VALUE` | angka desimal | ya |  | uang | `TreatyDifference.SalvageValue` |
| `EXCESS_LOSS` | angka desimal | ya |  | uang | `TreatyDifference.ExcessLoss` |
| `NET_PREMIUM` | angka desimal | ya |  | uang | `TreatyDifference.NetPremium` |
| `BALANCE_DUE_TO` | angka desimal | ya |  | uang | `TreatyDifference.BalanceDueTo` |
| `PREMI_ONP` | angka desimal | ya |  | uang | `TreatyDifference.PremiOnp` |
| `RI_COMM_ONP` | angka desimal | ya |  | persen | `TreatyDifference.RiCommOnp` (disalin) |
| `RESULT_ONP1` | angka desimal | ya |  | uang | `TreatyDifference.ResultOnp1` |
| `OVERIDDING_COMM_ONP` | angka desimal | ya |  | persen | `TreatyDifference.OveriddingCommOnp` (disalin) |
| `RESULT_ONP2` | angka desimal | ya |  | uang | `TreatyDifference.ResultOnp2` |
| `DEDUCTION1` | angka desimal | ya |  | uang | `TreatyDifference.Deduction1` |
| `DEDUCTION2` | angka desimal | ya |  | uang | `TreatyDifference.Deduction2` |
| `PPN_VALUE` | angka desimal | ya |  | uang | `TreatyDifference.PPNValue` |
| `PPH_VALUE` | angka desimal | ya |  | uang | `TreatyDifference.PPHValue` |
| `BALANCE_BEFORE_TAX` | angka desimal | ya |  | uang | `TreatyDifference.BalanceBeforeTax` |
| `BALANCE_BEFORE_PPH` | angka desimal | ya |  | uang | `TreatyDifference.BalanceBeforePPH` |

## T_POLIS_DIFFERENCE_SPREADING

`TreatyDifference.SpreadingRiskList` (`EDMTCalculateTreatyDifference` langkah 2.1). Pasangan baris antar generasi
menurut posisi (`NOURUT` = `.pxListSubscript`).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | — |
| `DIFFERENCE_ID` | teks | tidak | FK T_POLIS_DIFFERENCE, UQ (DIFFERENCE_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (DIFFERENCE_ID, NOURUT) | cacah | `.pxListSubscript` |
| `TREATY_TYPE` | teks | ya |  | kode | `.TreatyType` (disalin) |
| `TREATY_NAME` | teks | ya |  | teks | `.TreatyName` (disalin) |
| `SHARE_PERCENTAGE` | angka desimal | ya |  | persen | `.SharePercentage` (disalin) |
| `CLAIM_PERCENTAGE` | angka desimal | ya |  | persen | `.ClaimPercentage` (disalin) |
| `PREMIUM_SPREADED` | angka desimal | ya |  | uang | `.PremiumSpreaded` (selisih) |
| `CLAIM_SPREADED` | angka desimal | ya |  | uang | `.ClaimSpreaded` (selisih) |
| `PASANGAN_BERGESER` | teks | ya |  | penanda | penanda migrasi (ID-35), hanya SUMBER `PEGA` |
| `RUMUS_BERLAPIS` | teks | ya |  | penanda | penanda migrasi (ID-36), hanya SUMBER `PEGA` |

## T_POLIS_DIFFERENCE_INSTALMENT

`TreatyDifference.ListInstallment` — satu tingkat (ID-42; `EDMTCalculateTreatyDifference` langkah 3.1).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | — |
| `DIFFERENCE_ID` | teks | tidak | FK T_POLIS_DIFFERENCE, UQ (DIFFERENCE_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (DIFFERENCE_ID, NOURUT) | cacah | `.pxListSubscript` |
| `INSTALLMENT_NO` | bilangan bulat | ya |  | cacah | `.InstallmentNo` (disalin) |
| `DUE_DATE` | DATE | ya |  | tanggal | `.DueDate` (disalin) |
| `INSTALLMENT_PERCENTAGE` | angka desimal | ya |  | persen | `.InstallmentPercentage` (disalin) |
| `PREMIUM` | angka desimal | ya |  | uang | `.Premium` (selisih) |
| `PAYMENT_TOTAL` | angka desimal | ya |  | uang | `.PaymentTotal` (selisih) |
| `PASANGAN_BERGESER` | teks | ya |  | penanda | penanda migrasi (ID-35), hanya SUMBER `PEGA` |
| `RUMUS_BERLAPIS` | teks | ya |  | penanda | penanda migrasi (ID-36), hanya SUMBER `PEGA` |

## T_POLIS_XOL_LAYER_DIFFERENCE

`TreatyXOLDifferenceList(c).ValueList(l)` — hanya NonProporsional (`Activity/CalculateDifferenceEDM_act`
langkah 1.2). Induk per mata uang tidak bertabel (ID-7): dibangun ulang dari baris ini menurut `ID_CURRENCY`.

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | — |
| `DIFFERENCE_ID` | teks | tidak | FK T_POLIS_DIFFERENCE, UQ (DIFFERENCE_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (DIFFERENCE_ID, NOURUT) | cacah | posisi datar (mata uang, lapisan) |
| `LAYER` | teks | ya |  | kode | `.Layer` |
| `LAYER_TYPE` | teks | ya |  | kode | `.LayerType` |
| `LAYER_PART` | teks | ya |  | kode | `.LayerPart` |
| `LAYER_PART_TYPE` | teks | ya |  | kode | `.LayerPartType` |
| `CURRENCY` | teks | ya |  | kode | `.Currency` |
| `ID_CURRENCY` | teks | ya |  | kode | `.IDCurrency` |
| `DUE_TO` | teks | ya |  | kode | `.DueTo` (`DUE TO US` / `DUE TO YOU`, dari tanda selisih) |
| `GROSS_PREMI` | angka desimal | ya |  | uang | `.GrossPremi` |
| `DEDUCTION` | angka desimal | ya |  | uang | `.Deduction` |
| `NET_PREMI` | angka desimal | ya |  | uang | `.NetPremi` |
| `DUE_TO_VALUE` | angka desimal | ya |  | uang | `.DueToValue` |
| `BROKERAGE_FEE_SEBENARNYA` | angka desimal | ya |  | uang | `.BrokerageFeeSebenarnya` |
| `PPH_VALUE` | angka desimal | ya |  | uang | `.PPHValue` |
| `PPN_VALUE` | angka desimal | ya |  | uang | `.PPNValue` |
| `NET_PREMI_AFTER_PPH` | angka desimal | ya |  | uang | `.NetPremiAfterPPH` |
| `NET_PREMI_AFTER_PPN` | angka desimal | ya |  | uang | `.NetPremiAfterPPN` |
| `NET_PREMI_AFTER_TAX` | angka desimal | ya |  | uang | `.NetPremiAfterTax` |
| `PASANGAN_BERGESER` | teks | ya |  | penanda | penanda migrasi (ID-35), hanya SUMBER `PEGA` |
