# Layar Adjustment — catatan ukur, 5 Oktober 2026

Layar `InputTreatyInAdjustment` sistem lama, dibangun di
`frontend/pages/PenyesuaianKontrak.tsx` + `frontend/komponen/SisiPenyesuaian.tsx`,
spesifikasinya di `frontend/labelsPenyesuaian.ts`. Dokumen ini menyimpan **dari mana**
setiap keputusan datang, supaya ronde berikutnya mengukur dari titik yang sama.

Korpus: `D:\XML_NURE\Treaty In Adjustment`. Offset `@n` di kode adalah posisi bita
di berkas Section-nya **sesudah** seluruh `<pyIncludedRuleXML>` dibuang.

## 1 · `<pyIncludedRuleXML>` BERSARANG — membuangnya harus menghitung kedalaman

Pembuangan naif (cari pembuka → cari penutup pertama) pada
`Harness/InputTreatyInAdjustment.xml` membuang 112.918 bita dari 943.918 dan
**menghasilkan XML rusak**: pratinjau memuat pratinjau. Dengan kedalaman dihitung:

| Berkas | Ukuran | Terbuang | Sisa |
|---|---:|---:|---:|
| `Harness/InputTreatyInAdjustment.xml` | 943.918 | **909.282** | 34.636 |
| `Section/InputTreatyInAdjustment.xml` | 915.124 | 228 | 914.896 |

Pembuka @24376 harness baru tertutup @933638: **seluruh isi harness adalah pratinjau
Section `InputTreatyInAdjustment`**. Sumber sah tata letaknya Section itu, bukan
harness — dan angka "113 KB" yang beredar berasal dari pembuangan naif.

## 2 · Bentuk layar — `Section/InputTreatyInAdjustment.xml`

| @ | Syarat | Isi |
|---|---|---|
| 26602 | — | judul `Adjustment Treaty In` |
| 40422 | `DATASHOW=1` | kepala: ID Original · ID Revision · Reinsurance Type · Adjustment Type / `Adjustment Premium` · Material Type |
| 181268 | `(NonProportional && DATASHOW=1) \|\| (Proportional && DATASHOW=1)` | **Old Data** @192816 (`TreatyInNONProportionalOldData`) ‖ **New Data** @233311 (`TreatyInNONProportional`) |
| 276129 | `DATASHOW=1` | Attachment (`WorkAttachments`) |
| 317776 | `TreatyWarning.CARI1 != '' && DATASHOW=1` | Warning — **tidak dibangun**: halaman `TreatyWarning` diisi saat jalan, tidak tersimpan |
| 336222 | `DATASHOW=1` | `TreatyInActionButtons` |
| 380295 | `DATASHOW=1` | History — `TreatyIn.CommentList` @407001 |
| 481804 | `DATASHOW != 1` | mode daftar: Add Revision · Add Adjustment Premium · Show/Hide filter · grid 12 kolom + Edit/View |

⭐ Kedua cabang memakai Section **NON-proporsional yang sama** di tingkat harness.
Cabangnya dipilih **di dalam** Section: `TreatyInNONProportionalOldData.xml` @278100
(`ProportionType='Proportional'` → `TreatyInTabsProportionalOldData`) dan @292774
(`NonProportional` → `TreatyInTabsNonProportionalOldData`).

## 3 · Sumber data — `TREATY_IN_EDM` + `M_TREATY_IN_EDM`, BUKAN `VERSI_KONTRAK`

| Tabel | Baris | Catatan |
|---|---:|---|
| `KONTRAK` / `VERSI_KONTRAK` | **0** / **0** | jalur versi-dasar migrasi 440 belum punya data |
| `TREATY_IN_EDM` | 280 | `ID` `1000080/R01`…, `OLDID` → 263 ke `TREATY_IN`, 17 ke penyesuaian lain |
| `M_TREATY_IN_EDM` | 280 | `JSONDATA` = halaman `TreatyIn` utuh; **280/280 membawa `OLDDATA`** |

Section Old mengikat `TreatyIn.OLDDATA.*`; Section New mengikat `TreatyIn.*`. Satu
dokumen, dua halaman. Ke-280 dokumen terurai; seluruh 24.054 nilai skalarnya teks.

⛔ **Tiga ikatan panel Old yang menunjuk halaman AKAR**, disalin apa adanya:
`TreatyIn.ContractRefNo` @40302, `TreatyIn.Exclusions` @1777411,
`TreatyIn.SpecialConditions` @1815373 (keduanya hanya cabang NP).

## 4 · Ukuran yang menjadi dasar keputusan — dijaga `warisan_penyesuaian_db_test.go`

| Ukuran | Nilai | Keputusan |
|---|---|---|
| `M_TREATYIN_COMMENT` berpengenal penyesuaian | **0** | History dari `CommentList` dokumen (berisi 280/280) |
| `M_ATTACHMENTTREATY_2` berpengenal penyesuaian | **0** | Attachment tetap dibaca dengan `TreatyIn.ID`; panelnya `No items` |
| `Retention[].TreatyGroupName` terisi | **0** dari 1.580 | sel `Treaty Group` menampilkan `.TreatyGroup` |
| `CurrencyID` ↔ `Currency` | `10026`↔IDR 265, `10001`↔USD 240, nol ganda | sel Currency New menampilkan `.Currency` |
| `EGNPI[].Proportion` per sisi | 468 sisi: **365 tepat 100**, 97 meleset ≤ 10⁻¹⁸, 6 lain | golongan persen share — tetapi desimalnya kini **2 dari ekspor** (@572297), bukan 8 (§6.2) |
| Domain `Bordeaux` / `AccountingMode` / `AccountingModeNonProp` | {reporting, nonreporting} / {accounting, underwriting} / {loss, risk} | `PILIHAN_MEDAN` |
| Tab Value Difference tampil (`EDMState != 3 && EDMMaterialType == 1`, NP) | 52 dokumen; 18 di antaranya `IsProRate = true` | dua blok Before/After |

## 5 · Deret tombol — `Section/TreatyInActionButtons.xml`

`When/TreatyMasterInEDM.xml`: `EDMState = 1 / 2 / 3` — terpenuhi 280/280, jadi yang
berlaku blok @103200:

| Tombol | Syarat tampil | Di layar ini |
|---|---|---|
| Save @113510 | `ViewState !='1' && StatusAkseptasi != 'Resolve Complete'` @119719 | tampil sesuai syarat, **mati** |
| Close @130374 | ALWAYS; mati bila `DATASHOW3=='CLM'` | **hidup** — menutup, tidak menulis |
| Actions @148044 | peran operator @155994 | tampil, **mati** — syarat peran belum dievaluasi |

Blok dev `OperatorID.pyOrgDivision = 'IT'` @167976 tidak dibangun.

## 6 · Isi tab — DIBANGKITKAN dari ekspor, bukan disalin tangan (ronde 5 Oktober 2026, kedua)

Spesifikasi tab tulisan tangan ronde pertama **keliru** pada grid Limits Old: ia menyatakan
4 kolom, padahal Section-nya (@589778) berisi **8 sel** — kepala `Layers` hanya di sel
pertama, empat sel kepala kosong sesudahnya, dan sel `.pyTemplateRichTextEditor` yang
dulu terleburkan. Karena itu isi tab kini **dibangkitkan**:

```
python modul/treatyinadjustment/alat/ekstrak_kerangka.py "D:/XML_NURE/Treaty In Adjustment"
```

menulis `frontend/ekspor/kerangka.gen.ts` (37 tab, 7 include, kurs) dan
`backend/repository/kunci_kerangka_gen.go` (daftar-izin kunci). Pembangkit membaca
POHON Section — grid `pyTable/pyRows/rowdata` → `pyCells/rowdata`, desimal dari
`pyModes/pyDecimalPlaces`, syarat dari `pyContainerVisibleWhen` dan `pyUserData` — sesudah
`<pyIncludedRuleXML>` dibuang dengan menghitung sarang (§1).

- **Butir mati dibuang dan dicatat** (`DIBUANG`, 109 butir): penjaga `1=2`/`3=4`/`Never`/`false`
  diperiksa pada BUTIR yang dijaganya, bukan wilayah sekitarnya. Penjaga di sekitar
  `Total Amount in IDR`, `Total Proportion %`, `Total ROL` membungkus tombol; ketiga nilai itu hidup.
- **Kolom tombol** (`pxButton`, Add/Delete) dibuang — jalur tulis. **Butir bersyarat identitas
  operator** (blok dev) dibuang; nol nama orang masuk kode.
- **Syarat** dinilai `frontend/ekspor/syarat.ts`: satu penilai per TEKS syarat, semua atas halaman
  akar (`TreatyIn.*`, juga di panel Old). Teks yang tidak dikenal → butir tidak dirender, dan uji gagal.
- **Isi Retro** (`TreatyInFacultativeShareCalculation*`, `TreatyInFacultativeRetro`) tidak dibangkitkan —
  §17. Wadah tab Retro dan syaratnya dibangun; isinya kotak "jarang" (`IsMultipleRetro='true'`
  2 dari 280 penyesuaian, diukur ulang `TestRetroJarangDiPenyesuaian`).

### 6.1 Perbedaan Old lawan New yang DIPERTAHANKAN

| Tab | Old | New |
|---|---|---|
| strip NP | 8 tab | 11 tab (+ Event Limits, Value Difference, Information & Submit) |
| strip P | 7 tab | 11 tab |
| Limits NP | 8 kolom: Layers…, `Currency`, `100% Limits`, `Deductible` | 9 kolom: Layers…, pasangan IDR/USD |
| Share NP, label | `Facultative Share`, `Fakultative Brokerage` | `Share to Other Retro`, `Brokerage From Other Retro` (bersyarat `FacultativeShare >0`) |
| Share NP, grid RNM Share | 13 kolom | 14 kolom (+ `% Share`) |
| Share NP, `ShareFacultativeReinsurers` | tidak ada | ada, bersyarat |
| Exclusions/Special Conditions NP | mengikat AKAR (@1777411, @1815373) | halaman sendiri |
| Accumulation P, `Period` | mengikat AKAR (`TreatyIn.AccumulationPeriod`) | halaman sendiri |

### 6.2 Desimal per kolom — dari ekspor (keputusan pemilik proses 5 Oktober 2026)

Desimal yang dinyatakan ekspor dipakai dan **nol di ekornya dipertahankan**
(`padankan`, lapis modul, sesudah `formatNumber` inti dipanggil — `format.ts` tidak disentuh).
Yang tidak dinyatakan memakai aturan lama dan DIDAFTARKAN di sini. Tanda `%` tidak lagi
ditempel ke nilai: di sistem lama ia berdiri di label kolom atau sel terpisah.

| Kunci | Golongan | Desimal | Sel |
|---|---|---|---|
| `AggregateLimit` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 2 sel, mis. `TreatyInTabsNonProportional@1147598` |
| `AggregateLimit2` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 2 sel, mis. `TreatyInTabsNonProportional@1152044` |
| `Amount` | uang | **2** (ekspor) | 2 sel, mis. `TreatyInTabsNonProportional@581815` |
| `Amount` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 2 sel, mis. `TreatyInTabsNonProportional@79118` |
| `AmountIDR` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 2 sel, mis. `TreatyInTabsNonProportional@586894` |
| `BrokeragePercent` | persen | tidak dinyatakan → aturan lama (2, nol ekor dibuang) | 4 sel, mis. `TreatyInTabsNonProportional@1750505` |
| `BrokeragePercentP` | persen | **2** (ekspor) | 2 sel, mis. `TreatyInTabsProportional@717944` |
| `Conversion` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 2 sel, mis. `kurs-lama@206033` |
| `Deductible` | uang | **2** (ekspor) | 2 sel, mis. `TreatyInTabsNonProportional@1002988` |
| `Deductible` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 7 sel, mis. `TreatyInTabsNonProportional@1156492` |
| `Deductible2` | uang | **2** (ekspor) | 1 sel, mis. `TreatyInTabsNonProportional@1013782` |
| `Deductible2` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 7 sel, mis. `TreatyInTabsNonProportional@1160935` |
| `Earthquake` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 1 sel, mis. `TreatyInTabsNonProportional@343882` |
| `FacShare` | persen | tidak dinyatakan → aturan lama (2, nol ekor dibuang) | 1 sel, mis. `TreatyInTabsProportional@780884` |
| `FacShareBrokerage` | persen | tidak dinyatakan → aturan lama (2, nol ekor dibuang) | 1 sel, mis. `TreatyInTabsProportional@790889` |
| `FacultativeShare` | persen | tidak dinyatakan → aturan lama (2, nol ekor dibuang) | 2 sel, mis. `TreatyInTabsNonProportional@1797115` |
| `FacultativeShareBrokerage` | persen | tidak dinyatakan → aturan lama (2, nol ekor dibuang) | 4 sel, mis. `TreatyInTabsNonProportional@1806086` |
| `FloodJab` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 1 sel, mis. `TreatyInTabsNonProportional@400573` |
| `FloodNation` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 1 sel, mis. `TreatyInTabsNonProportional@457287` |
| `Limit` | uang | **2** (ekspor) | 2 sel, mis. `TreatyInTabsNonProportional@997908` |
| `Limit` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 7 sel, mis. `TreatyInTabsNonProportional@1129850` |
| `Limit2` | uang | **2** (ekspor) | 1 sel, mis. `TreatyInTabsNonProportional@1008073` |
| `Limit2` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 7 sel, mis. `TreatyInTabsNonProportional@1134288` |
| `MDP` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 7 sel, mis. `TreatyInTabsNonProportional@1138727` |
| `MDP2` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 7 sel, mis. `TreatyInTabsNonProportional@1143162` |
| `NetPremi` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 5 sel, mis. `TreatyInTabsNonProportional@2415059` |
| `NetPremi2` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 5 sel, mis. `TreatyInTabsNonProportional@2420122` |
| `ProRatePercent` | persen | **4** (ekspor) | 1 sel, mis. `TreatyInTabsNPValueDifferenceProRate@81824` |
| `Proportion` | persen | **2** (ekspor) | 2 sel, mis. `TreatyInTabsNonProportional@572297` |
| `RNMShare` | persen | tidak dinyatakan → aturan lama (2, nol ekor dibuang) | 5 sel, mis. `TreatyInTabsNonProportional@1741623` |
| `RNMShareP` | persen | **2** (ekspor) | 2 sel, mis. `TreatyInTabsProportional@711242` |
| `RSMDLimit` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 1 sel, mis. `TreatyInTabsNonProportional@286560` |
| `RnmShareDeducted` | persen | **2** (ekspor) | 1 sel, mis. `TreatyInShareProp@35412` |
| `RnmShareDeducted` | persen | tidak dinyatakan → aturan lama (2, nol ekor dibuang) | 4 sel, mis. `TreatyInTabsNonProportional@2107030` |
| `SharePct` | persen | **2** (ekspor) | 4 sel, mis. `TreatyInTabsNonProportional@1918452` |
| `TotalEgnpiAmount` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 2 sel, mis. `TreatyInTabsNonProportional@757649` |
| `TotalEgnpiProportion` | persen | **2** (ekspor) | 2 sel, mis. `TreatyInTabsNonProportional@801254` |
| `TotalLimitsROL` | persen | **2** (ekspor) | 2 sel, mis. `TreatyInTabsNonProportional@1599054` |
| `Value` | uang | **2** (ekspor) | 59 sel, mis. `TreatyInTabsNonProportional@164168` |
| `Value` | uang | tidak dinyatakan → aturan lama (4, nol ekor dibuang) | 20 sel, mis. `TreatyInTabsNonProportional@2238497` |

Kepala form New: `ProRatePercent` **4** (@327220). Kolom tanggal dan teks tidak bergolongan angka.

## 7 · Yang dibangun dan yang belum

Dibangun: mode daftar · kepala · kedua panel (kepala form, Rate of Exchange, strip tab per sisi ×
cabang) · **seluruh 37 tab dari kerangka** · Value Difference (Before = `ValueBeforeProrate.*`,
After = `ValueDifference.*`) · Attachment · deret tombol · History.

Belum: isi Retro (§17, disengaja) · rincian baris grid yang dibuka rule
`…!pyGridRowDetails` (Installment, Limits — rule itu TIDAK ada di ekspor) · isi
`TreatyInTabsAchievement` (Section-nya kosong di ekspor) · penyaring daftar · Activity penghitung
total (nilai tersimpan yang tampil). Medan di dalam tab baca-saja di kedua panel.

Pertanyaan yang belum terjawab: [`PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md`](PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md).
