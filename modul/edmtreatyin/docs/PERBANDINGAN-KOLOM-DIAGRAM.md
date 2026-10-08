# Perbandingan kolom — empat tabel proyeksi selisih EDM Treaty In lawan diagram grilling

*Disusun 06-10-2026 (eksekusi satu modul utuh). Pola: `modul/nbtreatyin/docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.*

**Acuan (urutan menang):** (1) XML korpus `EDM Treaty In` — baseline setiap medan (prompt eksekusi: *"XML adalah
baseline"*); (2) diagram grilling `Diagram-Skema-Tabel-NusantaraRe.xlsx` sheet *EDM Treaty In Prop* dan *EDM Treaty In
NonProp* (sel dikutip dengan alamatnya); (3) `docs/spec-penyimpanan-relasional.md` (ID-n, AC n). Bila diagram memuat
kolom yang **tidak ditulis satu aturan XML pun**, kolom itu **tidak dibuat** dan koreksinya dicatat di sini.

Hasil: **tepat empat `CREATE TABLE`** (360-363), ditagih `backend/migrasi_test.go`
(`TestMigrasiHanyaEmpatTabelProyeksiSelisih`, `TestKatalogSelisihSepakatDenganDDL`, `TestStrukturMemuatSetiapKolomDDL`,
`TestUangSelisihBertipeNumber3810TanpaFloat`, `TestTabelGenerasiNBTidakDibuatUlang`). Tabel generasi
(`T_WORK_POLIS`, `T_GENERAL_POLIS_TREATY`, `T_POLIS_*`) milik `nbtreatyin` — dibaca dan ditulis, **tidak dibuat**
(diagram Prop K146–K153: golongan *dasar*).

| # | Tabel | Migrasi | Sheet / sel diagram | Bentuk |
| ---: | --- | --- | --- | --- |
| 1 | `T_POLIS_DIFFERENCE` | 360 | Prop J85–J101 · NonProp J100–J116 · ringkasan Prop K154 | 1:1 `POLIS_ID` unik (G84 *"1:1 POLIS_ID UNIK"*) |
| 2 | `T_POLIS_DIFFERENCE_SPREADING` | 361 | Prop O102, R103–R114 · NonProp O117, R118–R129 · K155 | 1:N `DIFFERENCE_ID` |
| 3 | `T_POLIS_DIFFERENCE_INSTALMENT` | 362 | Prop O115, R116–R122 · NonProp O130, R131–R137 · K156 | 1:N `DIFFERENCE_ID`, **satu** tingkat (R122) |
| 4 | `T_POLIS_XOL_LAYER_DIFFERENCE` | 363 | NonProp O138, R139–R147 | 1:N `DIFFERENCE_ID` · NonProp saja |

**Tidak dibuat (sebab dikutip):**

| Bentuk | Sebab |
| --- | --- |
| `V_POLIS_DIFFERENCE` (view selisih) | diagram Prop J99 / NonProp J114 *"V_POLIS_DIFFERENCE DIBATALKAN"* — rumus milik lapisan aplikasi |
| `T_POLIS_BREAKDOWN_SPREAD` | diagram Prop K151 *"DIBATALKAN 23-09-2026"* |
| tabel induk `TreatyXOLDifferenceList` (per mata uang) | spec ID-7: kuncinya salinan kunci lapisan, nilainya jumlah lapisan (`CalculateDifferenceEDM_act` langkah 2) — dibangun ulang dari 363 (`models.BangunIndukSelisihXOL`) |
| tabel `OldData` | AC 9: `OldData` = baris generasi yang ditunjuk `OLD_POLIS_ID`, dibaca lewat repository (`halaman_edm.go`) |
| tingkat kedua angsuran selisih (`TreatyDifference.ListInstallment(n).InstallmentList`) | diagram R122 / R137 + ID-42: nol rujukan korpus EDM |

> ⛔ **RALAT tipe uang/persen** — sama dengan NB (RALAT 04-10-2026, diagram NB Prop F20): **`NUMBER(38,10)`**, bukan
> `NUMBER(38,8)` bunyi spec AC 58. Ditagih `TestUangSelisihBertipeNumber3810TanpaFloat`.

---

## 1 · `T_POLIS_DIFFERENCE` — 30 kolom

Asal medan: `Activity/EDMTCalculateTreatyDifference` langkah 1 (blok pertama; blok `HasEDMNo` = varian berlapis
**tidak** ditiru — diagram J91–J92, keputusan WO 23-09-2026, spec ID-30). XML langkah 1 menulis **27** medan
`TreatyDifference.*`; 23 disimpan, 4 `Total*` diturunkan.

### 1a · Kunci dan penanda proyeksi

| Kolom | Asal | Diagram | Keputusan | Bukti |
| --- | --- | --- | --- | --- |
| `ID` | dibangkitkan Go | — | dipertahankan | PK |
| `POLIS_ID` | generasi endorsemen (`T_GENERAL_POLIS_TREATY.ID`) | G84 *1:1 POLIS_ID UNIK*, J87 | dipertahankan | FK + `UQ_POLIS_DIFFERENCE_POLIS` |
| `NOPOLIS` · `PRODKE` · `EDM_NO` · `IDPEGA` | disalin dari generasinya | **J101** *"sertakan kunci yang mereka saring — NOPOLIS · PRODKE · EDM_NO · IDPEGA"* | dipertahankan | pembaca SQL tanpa join balik |
| `SUMBER` | `'PEGA'` pemuat migrasi · `'GO'` layanan | **J95–J97** | dipertahankan | `CK_POLIS_DIFFERENCE_SUMBER`; baris `PEGA` beku (`repository/selisih.go` `ErrSelisihBeku`) |

### 1b · Medan `TreatyDifference.*` (langkah 1)

| Kolom | Medan XML | Aturan XML | Keputusan |
| --- | --- | --- | --- |
| `INSTALLMENT` | `.Installment` | disalin | dipertahankan |
| `GROSS_PREMIUM` · `PREMI_OGP` · `RESULT_OGP1` · `RESULT_OGP2` | `.GrossPremium` · `.PremiOgp` · `.ResultOgp1` · `.ResultOgp2` | baru − `OldData.X` | dipertahankan (uang, J89) |
| `RI_COMM_OGP` · `OVERIDDING_COMM_OGP` · `RI_COMM_ONP` · `OVERIDDING_COMM_ONP` | `.RiCommOgp` · `.OveriddingCommOgp` · `.RiCommOnp` · `.OveriddingCommOnp` | **disalin** (persen) | dipertahankan (J89 *"persen = baru saja"*) |
| `CLAIM` · `SALVAGE_VALUE` · `EXCESS_LOSS` · `NET_PREMIUM` · `BALANCE_DUE_TO` | `.Claim` · `.SalvageValue` · `.ExcessLoss` · `.NetPremium` · `.BalanceDueTo` | baru − lama | dipertahankan |
| `PREMI_ONP` · `RESULT_ONP1` · `RESULT_ONP2` | `.PremiOnp` · `.ResultOnp1` · `.ResultOnp2` | baru − lama | dipertahankan |
| `DEDUCTION1` · `DEDUCTION2` | `.Deduction1` · `.Deduction2` | baru − lama | dipertahankan (XML mengurangi; satuan persen di jalur prop — diagram F133 — tidak mengubah rumus) |
| `PPN_VALUE` · `PPH_VALUE` · `BALANCE_BEFORE_TAX` · `BALANCE_BEFORE_PPH` | `.PPNValue` · `.PPHValue` · `.BalanceBeforeTax` · `.BalanceBeforePPH` | baru − lama | dipertahankan |
| *(tidak bertabel)* `TotalPremium` · `TotalClaim` | baru − lama | | **diturunkan** = jumlah baris 361 (`models.HitungTotalSelisih`; spec ID-15, ID-16) |
| *(tidak bertabel)* `TotalSharePercentagePremium` · `TotalSharePercentageClaim` | disalin | | **diturunkan** dari data baru |
| *(tidak dibuat)* `DUE_TO` | — | nol aturan menulis `TreatyDifference.DueTo` | **tidak dibuat** — spec ID-29 KOREKSI 06-10-2026; diagram J94 *"DUE_TO = dari TANDA selisih"* hanya berlaku untuk lapisan XOL (363) |

---

## 2 · `T_POLIS_DIFFERENCE_SPREADING` — 11 kolom

Asal: `EDMTCalculateTreatyDifference` langkah 2.1 — **enam** medan `TreatyDifference.SpreadingRiskList(.pxListSubscript).*`
(diagram R106 *"hanya 2 dari 6 medan"*).

| Kolom | Medan XML | Aturan XML | Diagram | Keputusan |
| --- | --- | --- | --- | --- |
| `ID` · `DIFFERENCE_ID` | — | — | O102 *1:N DIFFERENCE_ID* | dipertahankan |
| `NOURUT` | `.pxListSubscript` | posisi pasangan | R107, R110 | dipertahankan — kunci pasangan, bukan kunci dagang (ID-34) |
| `TREATY_TYPE` · `TREATY_NAME` | `.TreatyType` · `.TreatyName` | disalin | R104 | dipertahankan |
| `SHARE_PERCENTAGE` · `CLAIM_PERCENTAGE` | `.SharePercentage` · `.ClaimPercentage` | disalin | R105 | dipertahankan |
| `PREMIUM_SPREADED` · `CLAIM_SPREADED` | `.PremiumSpreaded` · `.ClaimSpreaded` | baru(n) − `OldData…(n)` | R106 | dipertahankan |
| `PASANGAN_BERGESER` · `RUMUS_BERLAPIS` | penanda migrasi | — | R108–R109, R113 | dipertahankan — hanya baris `SUMBER = 'PEGA'` (ID-35..ID-37, AC 40-43) |
| ~~`CURRENCY_ID`~~ | **tidak ada** | — | R104 / NonProp R119 *"kunci disalin TREATY_NAME · TREATY_TYPE · CURRENCY_ID"* | **tidak dibuat** — lihat koreksi K1 |

> ⛔ **KOREKSI K1 (06-10-2026, XML menang atas diagram).** Diagram *EDM Treaty In Prop* **R104** dan *NonProp*
> **R119** menyebut `CURRENCY_ID` sebagai kunci yang disalin. `Activity/EDMTCalculateTreatyDifference` langkah 2.1
> (dan blok `HasEDMNo` padanannya) hanya menulis `ClaimPercentage`, `SharePercentage`, `ClaimSpreaded`,
> `PremiumSpreaded`, `TreatyType`, `TreatyName` — nol `CurrencyID`. Diagram sendiri menghitung **enam** medan (R106
> *"hanya 2 dari 6 medan"*) = keenam medan XML tanpa `CURRENCY_ID`. Kolom tidak dibuat; diagram tidak disunting
> (berkas milik WO).

---

## 3 · `T_POLIS_DIFFERENCE_INSTALMENT` — 10 kolom

Asal: `EDMTCalculateTreatyDifference` langkah 3.1 — **lima** medan `TreatyDifference.ListInstallment(.pxListSubscript).*`
(diagram R119 *"hanya 2 dari 5 medan"*).

| Kolom | Medan XML | Aturan XML | Diagram | Keputusan |
| --- | --- | --- | --- | --- |
| `ID` · `DIFFERENCE_ID` | — | — | O115 | dipertahankan |
| `NOURUT` | `.pxListSubscript` | posisi pasangan | R121 | dipertahankan |
| `INSTALLMENT_NO` · `DUE_DATE` | `.InstallmentNo` · `.DueDate` | disalin | R117, R120 (`INSTALLMENT_NO` kunci periksa) | dipertahankan |
| `INSTALLMENT_PERCENTAGE` | `.InstallmentPercentage` | disalin | R118 | dipertahankan |
| `PREMIUM` · `PAYMENT_TOTAL` | `.Premium` · `.PaymentTotal` | baru(n) − lama(n) | R119 | dipertahankan |
| `PASANGAN_BERGESER` · `RUMUS_BERLAPIS` | penanda migrasi | — | R121 | dipertahankan |

Satu tingkat saja: diagram R122 / NonProp R137, spec ID-42.

---

## 4 · `T_POLIS_XOL_LAYER_DIFFERENCE` — 21 kolom · NonProp saja

Asal: `Activity/CalculateDifferenceEDM_act` langkah 1.2.1 / 1.2.3 / 1.2.4 / 1.2.5 — **17** medan
`TreatyXOLDifferenceList(c).ValueList(l).*` (sapuan `PropertiesName` 06-10-2026).

| Kolom | Medan XML | Aturan XML | Diagram NonProp | Keputusan |
| --- | --- | --- | --- | --- |
| `ID` · `DIFFERENCE_ID` | — | — | O138 | dipertahankan |
| `NOURUT` | posisi lapisan dalam daftar datar (mata uang ke-c, lapisan ke-l) | — | R145 | dipertahankan |
| `LAYER` · `LAYER_TYPE` · `LAYER_PART` · `LAYER_PART_TYPE` · `CURRENCY` · `ID_CURRENCY` | `.Layer` · `.LayerType` · `.LayerPart` · `.LayerPartType` · `.Currency` · `.IDCurrency` | disalin | R140 | dipertahankan |
| `DUE_TO` | `.DueTo` | dari **tanda** `DueToValue` (`"DUE TO US"` / `"DUE TO YOU"`) | J109 | dipertahankan |
| `GROSS_PREMI` · `DEDUCTION` · `NET_PREMI` · `DUE_TO_VALUE` | `.GrossPremi` · `.Deduction` · `.NetPremi` · `.DueToValue` | selisih (batas bawah 0, prorata, pembatalan mentah) | R141 | dipertahankan |
| `BROKERAGE_FEE_SEBENARNYA` · `PPH_VALUE` · `PPN_VALUE` | `.BrokerageFeeSebenarnya` · `.PPHValue` · `.PPNValue` | selisih | R142 | dipertahankan |
| `NET_PREMI_AFTER_PPH` · `NET_PREMI_AFTER_PPN` · `NET_PREMI_AFTER_TAX` | `.NetPremiAfterPPH` · `.NetPremiAfterPPN` · `.NetPremiAfterTax` | selisih | R143 | dipertahankan |
| `PASANGAN_BERGESER` | penanda migrasi | — | R145 (kunci periksa `LAYER + LAYER_PART + ID_CURRENCY`, R144) | dipertahankan |
| *(tidak dibuat)* `RUMUS_BERLAPIS` | — | sisi NonProp tak punya varian kedua | **R147** *"tanpa RUMUS_BERLAPIS"* | **tidak dibuat** |

Rumus XOL dipakai **apa adanya** (keputusan WO ID-28/30 06-10-2026); diagram J108 *"sisi nonprop memang sudah
seragam"* sejalan.
