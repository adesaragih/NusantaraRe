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
- ~~**Kolom tombol** (`pxButton`, Add/Delete) dibuang — jalur tulis.~~ ⛔ **Dibalik 7 Oktober 2026**
  (§8): tombol — juga tombol ikon tanpa `pyLabel` — dibawa beserta aksinya. **Butir bersyarat
  identitas operator** (blok dev) tetap dibuang; nol nama orang masuk kode.
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

Belum: isi Retro (§17, disengaja) · ~~rincian baris grid yang dibuka rule
`…!pyGridRowDetails` (Installment, Limits — rule itu TIDAK ada di ekspor)~~ — **RALAT §9: ada di
ekspor dan kini dibangkitkan** · isi
`TreatyInTabsAchievement` (Section-nya kosong di ekspor) · penyaring daftar · Activity penghitung
total (nilai tersimpan yang tampil). Medan di dalam tab baca-saja di kedua panel.

Pertanyaan yang belum terjawab: [`PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md`](PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md).

## 8 · Ronde 7 Oktober 2026 — panel New dapat disunting, tombol dan rumus hidup

Permintaan pemakai: *"bandingkan terlebih dahulu antara pega dan di aplikasi kemudian lakukan
pengerjaannya samakan mulai dari fungsi, rumus, button, posisi design juga"*.

| Bagian | Yang berubah | Bukti / ukuran |
|---|---|---|
| Rate of Exchange | Dibaca dari `TREATYEXCHANGEYEARLY` per Treaty Year **sisi masing-masing** (`backend/repository/kurs_tahunan.go`), bukan larik `CurrencyList` yang tidak pernah mendarat | terisi 279/280 kedua sisi (`TestKursPenyesuaianDariKursTahunan`) |
| Peta pendaratan | + `Installment`, `CoInScale`, `Share`, `ShareReins`, `FacultativeShareList`, `ShareFacultativeReinsurers` — disalin dari peta Treaty In, diadu `uji/lintasmodul` | setiap baris terbaca: Share 1.938, Installment 518, Co-Ins 68, … (`TestLarikTambahanPenyesuaianTerbaca`) |
| Visibilitas | `pyIsVisibilityOption = ALWAYS` MENIMPA `pyContainerVisibleWhen` — grid Co-Ins Scale Prop dipulihkan | satu-satunya wadah `ALWAYS` + `1=2` di Section Adjustment |
| Tombol | Pembangkit membawa tombol (label, ikon, Activity + parameter, `pyPreDataTransform`, syarat tampil, `pyDisabledWhen`) dan kolom tombol grid | 43 tombol; butir dibuang 109 → 98 |
| Baca-saja | Per medan/sel: `Read-only` tanpa syarat = selalu; `Read-only` + `pyReadOnlyCondition` = bila syarat benar; `pyDisabledWhen` menambah kunci (termasuk `EDMMaterialType = 2`) | pola diukur: 103 · 41 · 23(+19) · 18 sel |
| Mode Edit | Panel New (`ViewState 0`) dapat disunting; Add/Delete bekerja; panel Old tetap baca-saja. Suntingan di keadaan layar — Save mati | `penyesuaian.test.tsx` § mode sunting |
| Tambah baris | Disalin dari Activity: `TreatyInNonAddItem` (retention/egnpi/limits/sharereins/sharefacname), `TreatyInPropAdd`, `TreatyInAddCurrency`, DataTransform `TreatyInAddAccumulation` | `komponen/aksiTombol.ts` |
| Rumus tombol | Rute `/api/treaty-in/hitung/*` modul Treaty In — Activity-nya IDENTIK di kedua korpus (14/14 diadu). Tersambung: Update Total Retention, Update Total EGNPI, Update EGNPI Value, Apply Reporting Period | `komponen/rumus.ts`, `tombol.test.ts` |
| Rumus kepala | `TreatyInSetTreatyYear` (Commencement) dan `TreatyCalculateProratePct` (Effective Date / Is Pro Rate) | diukur atas 22 penyesuaian ber-Pro Rate: 6 terbaru cocok ekspor, 16 lama menyimpan Commencement → Effective (versi sebelumnya) — `TestRumusProRataTerukur` |

**Belum** — daftar persisnya dijaga `tombol.test.ts` dan berubah tiap ronde; keadaan terkini di §9.
(Paragraf asli ronde ini menyebut Limits, Share, dropdown, dan Choose Ceding/SoB "belum"; keempatnya
tersambung pada ronde yang sama sesudah pohon bersarang dibaca.)

## 9 · Ronde 7 Oktober 2026 (lanjutan) — rincian baris, perilaku `change`, rumus Installment

| Bagian | Yang berubah | Bukti |
|---|---|---|
| Rincian baris | Grid `pyGridProps/pyEditingMode = expandPane` menyebut flow action di `pyEditAction`; flow action itu menunjuk Section lewat `pySectionReference`. Pembangkit kini mengikutinya (`seksi_flowaction`) dan menulis **22 Section rincian** ke `KERANGKA_RINCIAN`, termasuk rincian di dalam rincian (CoBList di Layers, DetailLimits/DetailShare di Prop). Ikatan `.X` = halaman BARIS (`konteksBaris`). Grid menampilkan kolom ▸/▾, tertutup semula | gambar Pega 38 (Installment, baris IDR terbuka) dan 34 (Share, ▸ per layer); `penyesuaian.test.tsx` § rincian |
| Judul blok | Kepala blok tampil hanya bila `pyIncludeHeader` ≠ `false` — **bukan** `pyContainerFormat = NOHEADER` (Exclusions dan *Summarry of RNM Share* ber-NOHEADER, judulnya TAMPIL di gambar 34/39). 29 judul kini tersembunyi: `hiddden`, `hidden, reference`, `Testing`, salinan *Account Reporting Period* di tab Accumulation, *After Pro Rate Calculation* cabang tanpa Pro Rate | gambar 19 (Accumulation hanya *Accumulation Control*), 34, 39 |
| Perilaku `change` | Medan dan sel membawa `aksiUbah` (`pyBehaviors` ber-`pyEvent = change`) — 249 baris di kerangka. Dijalankan saat isian ditinggalkan atau Enter, hanya bila nilainya berubah; rantai yang salah satu Activity-nya belum punya rumus TIDAK dijalankan | `KerangkaTab.tsx` `PemicuUbah`, `rantaiUbah` |
| Rumus Installment | Rute baru `POST /api/treaty-in/hitung/angsuran` (modul Treaty In, Activity identik di kedua korpus): `TreatyInSetValueInstallment` (isian Installment tanpa status; Update Value `status=update`; langkah 10 memulihkan dari sisi Old bila `EDMState` 1/2/3), `SetTotalInstallment` (sel % Installment `editpercentage`, sel Amount), `TreatyInNPSetTotal(installment)` | `treatyin/backend/services/hitung_angsuran.go` + 10 uji — termasuk angka gambar 38: 4 × 25 % atas 644.674.819,59 = 161.168.704,90 |
| Hapus baris | `deleteRow` yang DISUSUL Activity (Remove + LimitCalculation / CalculateDeduction) kini mati sampai Activity itu berumus — dulu menghapus tanpa menghitung ulang (rantai separuh jalan) | `aksiTombol.ts` `rantaiSesudahHapus` |

⚠️ **Yang TIDAK berbukti di ekspor, diputuskan dan dinyatakan** (`hitung_angsuran.go`): langkah 9/10
membaca baris sumber yang bisa tidak ada (jumlah angsuran naik) — nilai hasil hitungan dipertahankan;
isian Installment kosong/bukan angka diperlakukan `< 1`; `DueDate` ditulis bentuk tersimpan
`YYYYMMDD` (panjang maksimum `DueDate` di data Pega = 8). Activity `TreatyInUpdatePaymentDate`
(sel Due Date) dan `TreatyInUpdatePaymentDate_Act` (sel WPC) **tidak ada di ekspor** — Payment Date
tidak dihitung ulang.

⚠️ **Disalin apa adanya walau janggal:** Update Value memulihkan persen dari grid sebelumnya (dan dari
sisi Old di Adjustment) tetapi TIDAK menghitung ulang Amount — `SetTotalInstallment` dipanggil tanpa
parameter, jadi melompat ke label `calc`. Isian Installment < 1 mengosongkan tab sebelum pesan.

⭐ **Sebagian besar yang tercantum di paragraf berikut sudah disambungkan — lihat §10.**

**Masih terbuka** (dijaga daftar persis di uji): 11 sumber dropdown rincian Layers/Share/DetailLimits
(`pilihan.test.tsx`), 36 kunci angka rincian tanpa golongan tampil (`penyesuaian.test.tsx`), tombol
rincian Prop Limits (Add*, Achievement, Remove) serta Remove Share/Delete Layers (`tombol.test.ts`),
dan perilaku `change` yang Activity-nya sudah punya rute Treaty In tetapi belum dipetakan
(LimitCalculation, CalculateDeduction, PremiumReserveCalculate, DetailCalculation,
SetReinstatementPct, SetAmountConversion, TotalEgnpi, TreatyInXOLAddSpreading, TreatyInSetBrokerage).
Value Difference (`TreatyEDMCalculateDifference`) belum.

## 10 · Ronde 7 Oktober 2026 (ketiga) — ketergantungan antartab dan rumus Section rincian

Permintaan pemakai: isian tidak hilang saat pindah tab, dan tab yang nilainya saling bergantung
terhubung, sesuai alur XML dan Pega.

**Isian panel New memang tidak hilang saat pindah tab.** Keadaannya hidup di `SisiForm`, di atas
strip tab, dan tidak punya efek reset. Yang kurang adalah RUMUS yang menghubungkan tab: perilaku
`change` di Section rincian berjalan di halaman BARIS, sedangkan banyak Activity-nya menulis ke
AKAR panel.

### 10.1 Infrastruktur

| Bagian | Isi | Di mana |
|---|---|---|
| Jalur | Tiap konteks tahu jalurnya dari akar (`Limits(2).Detail(1)`); rantai berjalan atas akar + jalur, jadi rumus rincian dapat membaca dan menulis akar | `baris.ts` (`ambilHalaman`, `ubahHalaman`), `konteksBaris(…, langkah)` |
| Gabung tiga arah | Hasil rantai membawa akar sebelum (`awal`) dan sesudah (`akhir`). `SisiForm` menerapkannya atas keadaan TERKINI: yang rumus ubah menimpa, isian yang diketik selama rute menjawab bertahan | `gabungTigaArah`, `SisiPenyesuaian.tsx` |
| Syarat aksi | `pyActionConditions` kini dibangkitkan (`AksiTombol.syarat`) dan dinilai atas halaman PEMICU (akar + skalar bertitik baris sel: `.Note`, `.Layer`). Aksi yang syaratnya gagal dilewati | `ekstrak_kerangka.py` `syarat_aksi`, `syarat.ts` |
| Parameter DT | `pyDataTransformParams` kini dibangkitkan (`AksiTombol.paramDT`), mis. `AddLimitRetentionCession(type = limit/retention/cession)` | `ekstrak_kerangka.py` |
| Kotak centang | Perilaku `click` pada `pxCheckbox` = perubahan nilainya (Share Across The Board) | `aksi_ubah` |
| Parameter `.X` | Dibaca dari halaman pemicu; `.pxListSubscript` = sel/jalur + 1 | `nilaiParam` |
| Langkah bersyarat tanpa rumus | Rantai diterima; bila syaratnya terpenuhi, rantai BATAL dengan pesan dan nol perubahan | `rumus.ts` `jalankan` |

### 10.2 Yang kini tersambung (rute `/api/treaty-in/hitung/*`, nol tulisan basis data)

| Section / tab | Pemicu | Activity / DT | Rute |
|---|---|---|---|
| Share Prop | Refresh, `% RNM Share`, Option, Fac Share | `TreatyInPropshare` | `share-prop` `share` |
| `DetailShare` | `% RNM Share` · Spreading Type · sel spreading | `TreatyInPropshareDetail` · `FetchQSfromMaster(.SpreadingTypeID)` · `SetSpreadName` | `share-prop` `detail`/`spreading`/`sebar-nama` |
| Accumulation | Period · sel Reporting Date / Submission Days | DT `TreatyInDeleteAccumulationLists` + `TreatyInSetAccountReport` · `TreatyInAccumulationSetSubDue` | `akumulasi` |
| Reporting Period | sel Initial Date | `TreatyInSetReport(.InitialDate, .AutoCalculate)` | `periode-pelaporan` (`awal`) |
| Exclusions | `ExclusionsP` / `SpecialConditionsP` | DT `TreatyInCopyConditions` | di layar |
| `LimitProportional` | Treaty Type · Add | `SetTreatyTypeName_Act` + DT `TreatyTypeSetIndex` · `AddClassofBusiness` | di layar |
| `DetailLimits` | Treaty Group | `SetTreatyGroupName_Act` → DT `SetDetailsID` + `FetchQSfromMaster("")` [QUOTA SHARE] → `LimitCalculation(surplus)` [SURPLUS] | `share-prop`, `limit` |
| `DetailLimits` | QS % / Surplus · sel 100% Limit / Retention (Value, Currency, Remove) | `LimitCalculation` [syarat `.Note`, autocalculate `.Layer`] | `limit` |
| `DetailLimits` | sel Deduction (CurrencyID/Deduction/%) · Remove | `CalculateDeduction(sts, index)` | `limit-deduksi` |
| `DetailLimits` | `% Premium Reserve` | `PremiumReserveCalculate` | `limit-cadangan` |
| `DetailLimits` | sel mata uang | `SetCurrName_Act` (lihat PERTANYAAN-TERBUKA §14) | di layar |
| `DetailLimits` | Add (Limit/Retention/Cession, Deduction, Reserve, PLA, Cash Loss, Claim Coop, EPI) | DT `AddLimitRetentionCession(type)`, `AddDeduction`, `AddValue(type)` | di layar |
| Share Non-Prop | `% RNM Share`, Fac Share, Brokerage From Other Retro · `% Brokerage` · Share Across The Board · Update Summary | `TreatyInXOLAddSpreading` · `SetBrokerage`→`XOLAddSpreading` · keduanya bersyarat · `setValue` + `NonAddItem(share)` (+XOL, SetBrokerage tercakup; langkah EDM `7897987` tak pernah jalan) | `share-np` `rnm`/`set-brokerage`/`summary` |
| rincian `Share` | `% RNM Share` · Spreading Type · sel spreading Pct · sel Deduction, Add, Remove | `TreatyInXOLAddSpreadingDetail` (+`…Actual` bila `EDMState = 3`, lihat §14) · `FetchQSfromMasterXOL` · `SetSpreadingXOL` · `CalculateDeduction`, `AddDeduction` | `share-np` `rnm-baris`/`spreading-type`/`spreading-pct`/`deduksi` |
| `Layers` | Adj Rate · MDP % / MDP Min % · Reinstatement Value · sel Reinstatement (3 DT) · Treaty Group (Add, Delete, sel) | `DetailCalculation(adj/mdp)` · `SetReinstatementPct` · `ReCalculateReinstatement`/`CalculateReinstatement`/`CalculateReinstatementPct` · `TreatyTypeSetIndex`, `SetIndexLayer_DT`, `TotalEgnpi` | `limit-np` |
| `CoBList` | Treaty Group | `TotalEgnpi` ×2 (satu panggilan) | `limit-np` `egnpi` |
| `DetailEGNPI` | Amount, Currency | `SetAmountConversion` | `egnpi` `konversi` |

Hanya medan yang Activity TULIS yang digabung kembali (`DITULIS_LIMIT_CALCULATION`,
`DITULIS_LIMIT_NP` — disalin dari Treaty In, nol impor). Uji: `rantai-jalur.test.ts` (32).

### 10.3 Ralat ikut ronde ini

- `pilihan.ts`: dropdown `BrowseTreatyGroup_RD` / `BrowseReinsuranceType_RD` bernilai `.ID`
  (`TreatyGroupID` rincian Limits, `TreatyTypeID`) dulu mengisi NAMA ke medan pengenal.
- Tanggal ke rute dikirim bentuk simpan `YYYYMMDD` (`keSimpanTanggal`): kotak tanggal sel
  mengirim `YYYY-MM-DD`, yang tidak dibaca rute Treaty In.
- `TreatyInSetReport` kini membaca parameternya (`autocalculate`, `startdate`).
- `PemicuUbah`: kotak tanggal ketik menyimpan isiannya di blur yang SAMA dengan pemicu `change`, jadi sel
  Initial Date / Reporting Date tidak pernah memicu rumusnya. Pemeriksaan kini diulang sekali sesudah
  render berikutnya. ⚠️ Tanpa jsdom di repo, jalur ini belum diuji otomatis; perlu dicoba di peramban.

### 10.4 Masih belum

Achievement rincian Limits (`GetAchievement`, `GenerateCSVTreaty`, `InsertToLogAchievement` —
halaman sesi `SearchData`), Retro (§17), Value Difference (`TreatyEDMCalculateDifference`),
`TreatyInUpdatePaymentDate`/`_Act` (tidak diekspor), `TreatyInShowHideFacShare` (halaman sesi),
Submit/Decline dan Save (jalur tulis).

## 11 · Ronde 7 Oktober 2026 (keempat) — sisa §10.4 dikerjakan

| Bagian | Isi | Bukti / rute |
|---|---|---|
| Payment Date | Sel Due Date → `TreatyInUpdatePaymentDate` (halaman itu), sel WPC → `_Act` (semua halaman, baris ber-Due Date): `PaymentDate = @addCalendar(DueDate, …, WPC + 1 hari)`. Kedua Activity DIUNGGAH pemakai (`D:\NUSANTARA RE APP\TreatyInUpdatePaymentDate*.xml`) | `angsuran` `tanggal-bayar` / `tanggal-bayar-semua` |
| Retro | Tab `Retro` (ketiga cabang) dan `Actual Retro` DISEMBUNYIKAN dari strip — keputusan pemilik proses: *"untuk sementara retro di hide dari tab sampai ada perintah dari developer"*. Kerangka dan syaratnya tetap dibangkitkan | `TAB_DISEMBUNYIKAN` (`SisiPenyesuaian.tsx`) |
| Value Difference | `TreatyEDMCalculateDifference` + sembilan Activity panggilannya (identik di kedua korpus): selisih EGNPI, Limits, Share, Deduction, ringkasan, total, Installment, Pro Rate (`ValueBeforeProrate`) | `selisih` (`hitung_selisih.go`, 7 uji) |
| Achievement rincian Limits | Refresh dan Quarter Year → `GetAchievement`; As At → DT `Reset_DT`; Generate Excel → `GenerateCSVTreaty` (13 kolom, SEMUA baris, `CSVAchievementTreatyIn.xlsx`). Halaman SESI (`SearchData.CARI1/2`, `FlagExcel.CARI1`, `TempQuarter*`) kini dibangkitkan (`dari = 'sesi'`) dan hidup di akar panel | `achievement`, `unduhAchievement.ts` |
| Cabang Adjust Premium | `EDMState = 3` (Non-Prop) → `TreatyInTabsNonProportionalAdjustPremi` (@566980): Actual GNPI, Actual Limits, Actual Share, Premium Adjustment, Information & Submit. Halaman `ActualValue`; selisih = Actual − Treaty In | `aktual` (`hitung_aktual.go`, 6 uji) |
| Rincian Share `% RNM Share` (EDMState 3) | `TreatyInXOLAddSpreadingDetailActual(idx)` kini dibangun — dahulu rantainya batal | `aktual` `rnm-baris` |

**Rute `aktual` — Activity per tombol:** Actual GNPI · Update Total = `TreatyInActualUpdateValue`
(17 langkah); Actual Share · Update Summary = `setValue` + `TreatyInActualUpdateValueShare`;
Actual Limits / Premium Adjustment · Update Total = `TreatyInNPSetTotalActual(limits)` +
`TreatyInSummaryLimitActual` + `TreatyInSummaryLimitShareActual`; Add Actual GNPI =
`TreatyInNonAddItem(actualpremium)` (mata uang dari `Retention(1)`).

⚠️ **Disalin apa adanya walau janggal** (rinciannya di kepala `hitung_selisih.go` /
`hitung_aktual.go`): `ActualValue.Limits` ditimpa salinan Limits Treaty In tiap Update Total;
selisih premi/MDP/Gross/Net dinolkan bila Actual < Treaty In; beberapa Activity menulis atau
membuang total AKAR (`TotalSpreaded*`, `TotalFac*`, `LimitFacShareSummaryList`,
`TotalLimitPremiEarnNP/MDPNP`, `RnmShareDeducted`); langkah Deduction Value Difference membuang
baris Brokerage fee yang baru ditambahkannya.

**Masih belum:** Submit Achievement (`InsertToLogAchievement` → `POOLDATA.LOG_ACHIEVEMENT`;
migrasi 423 memutuskan log itu TIDAK menjadi tabel kedua — jalur tulisnya bersama Save, yang
dikerjakan sesi lain); `Show Facultative Share` (harness tanpa Activity); isi Retro (§17, kini
disembunyikan). Belum diukur atas data: tabel pendaratan sedang dikosongkan.

## 12 · Ralat 8 Oktober 2026 — penyesuaian tanpa pendaratan tidak dapat disunting

Laporan pemakai: `1001540/R01` dibuka lewat Edit, kepala kosong (ID Original, Reinsurance Type,
Adjustment/Material Type "tidak ada"), dan panel Old/New tidak muncul sama sekali.

**Sebab:** `BacaPenyesuaianPendaratan` membaca HANYA tabel pendaratan `T_TREATY_*`, yang sengaja
dikosongkan 7 Oktober. Grid daftar membaca `TREATY_IN_EDM`, jadi barisnya tampil; begitu dibuka,
sisi New tidak punya `ProportionType`, dan wadah panel bersyarat `Proportional`/`NonProportional`.

**Perbaikan:** kunci kepala yang tidak ada di pendaratan dilengkapi dari kolom `TREATY_IN_EDM`
(`ProportionType`, `OLDID`, `EDMState`, `EDMMaterialType`, nama kontrak, Ceding, SoB, tanggal, posisi,
status, Treaty Year) — cara yang sama dengan `BacaDokumenMaster` (picker Add). Nilai pendaratan tetap
menang bila ada.

⚠️ Hanya sisi **New**. Sisi Old (`OLDDATA`) adalah salinan saat penyesuaian dibuat; kepala master hari
ini bukan salinan itu, jadi panel Old tetap kosong sampai pendaratannya terisi. Isi tab (Limits, Share,
dst.) juga kosong sampai pendaratannya terisi — kolom `TREATY_IN_EDM` hanya kepala.

## 13 · 8 Oktober 2026 — halaman `ActualValue` tersimpan ber-MASTERID sendiri

Keputusan pemakai: *"masukkan seperti yang ada di modul treaty in asal ada master id nya pasti
dapat"*. `ActualValue` berbentuk dokumen `TreatyIn`, jadi ia mendarat di tabel `T_TREATY_*` yang SAMA
lewat peta yang sama — nol tabel baru.

| Sisi | MASTERID |
|---|---|
| New (akar) | `1000080/R01` |
| Old (`OLDDATA`) | `1000080/R01#LAMA` |
| Actual (`ActualValue`) | `1000080/R01#AKTUAL` |

- **Tulis** (`treatyin/services/simpan_penyesuaian.go` `halamanAktual`): `SaveTreatyIn_EDM_Act` [2] —
  EDMState `3` memakai `ActualValue` isian layar apa adanya; selainnya [2]–[4] menyalin `TreatyIn` tanpa
  `OLDDATA`, `ValueDifference`, dan `ActualValue` lamanya. Decline offer ikut mengosongkan `#AKTUAL`.
- **Baca** (`repository/pendaratan_aktual.go`): `#AKTUAL` kembali ke sisi New sebagai kunci bertitik
  (`ActualValue.EGNPI`, `ActualValue.TotalEgnpiAmount`, pohon `ActualValue.Share` …).
- Akhiran dijaga sama di kedua modul oleh `uji/lintasmodul` (`TestAkhiranSisiAktualSama`).
- Kunci Actual yang tidak punya kolom dilaporkan `ActualValue.<kunci>` — hanya untuk EDMState 3 (salinan
  EDM 1/2 identik dengan sisi New, kuncinya sudah dilaporkan sekali).

## 14 · 8 Oktober 2026 — Submit Achievement dan Show Facultative Share HIDUP

**Submit Achievement** (`InsertToLogAchievement`, sub-tab Achievement `DetailLimits`; Treaty In dan
Adjustment). Keputusan pemakai: sasarannya `POOLDATA.LOG_ACHIEVEMENT` apa adanya (tabel sudah ada,
17 kolom sama dengan `InsertToLogAchievement_SQL`); setiap klik menyisipkan (Pega tanpa pencegah
duplikat), KECUALI baris ber-Quarter kosong — di Pega baris itu tersisip dengan nilai baris
sebelumnya; di sini dilewati dan cacahnya dilaporkan.
- Rute `POST /api/treaty-in/achievement/log` (`treatyin/services/log_achievement.go`), satu transaksi;
  angka `TO_NUMBER(koef)/POWER(10,skala)`; angka tak terbaca DITOLAK, bukan dicatat 0.
- Operator: akun yang menekan (Pega: `OperatorID.pxUpdateOperator`).
- Pega tidak menampilkan pesan; layar menambah "n baris Achievement tercatat".

**Show Facultative Share** (`showHarness` → `TreatyInFacultativeShareCalculation`, panel New;
`…OldData`, panel Old; tampil bila `FacultativeShare > 0`). Pembangkit kini mencatat `harness` dan
`jendela` aksi, membangkitkan Section isi harness (`TreatyInActualFacultativeShareCalculation`,
`TreatyInFacultativeShareCalculationOldData`) dan peta `KERANGKA_HARNESS`. Tombol membuka jendela
"Facultative Calculation" berisi Section itu atas halaman yang SAMA (`pySubmitData = Yes`); nol Activity
saat dibuka/ditutup. Penyertaan inline Section itu di tab Retro TETAP tersembunyi (keputusan Retro).

## 15 · Ralat 8 Oktober 2026 — mode baris grid (`pyRowEditing`) kini dibaca

Laporan pemakai: sel `Kind of Treaty` tab **Achievement In IDR** panel New tampil sebagai dropdown,
padahal di Pega tidak dapat disunting. Sebabnya: pembangkit hanya membaca kunci PER SEL
(`pyEditOptions`, `pyReadOnlyCondition`, `pyDisabledWhen`), dan sel itu `Auto` — padahal Pega
menentukan juga dari MODE BARIS GRID di `pyGridProps/pyRowEditing`. Terukur di korpus
`Treaty In Adjustment/Section`:

| `pyRowEditing` | grid | Arti di Pega | Di sini |
|---|---|---|---|
| `row` | 87 | sel disunting di tempat | kunci per sel berlaku seperti semula |
| `readOnly` | 183 | grid seluruhnya baca-saja | setiap kolom data `baca = "selalu"` |
| `masterDetail` (`expandPane`) | 45 | baris TAMPIL saja; disunting di panel rincian | setiap kolom data `baca = "selalu"`; rincian tetap mengikuti kuncinya sendiri |

153 grid di kerangka berubah (kebanyakan grid total/ringkasan `readOnly`). Kolom TOMBOL baris
(Delete) tidak tersentuh. Setiap grid `masterDetail` di panel New punya rincian berisi medan
yang dapat disunting (mis. `MaxRetention`: TreatyGroup · Currency · Amount · Note; `LimitProportional`:
`TreatyTypeID` — Kind of Treaty baris baru dipilih di sana). `pyRODetails = true`
(`LimitFacRetro_Sec` `.FacultativeLimits`) dibawa sebagai `rincianBaca`: rinciannya ikut baca-saja.
Dijaga `frontend/mode-baris-grid.test.tsx`.

## 16 · 8 Oktober 2026 — lampiran master ikut ke penyesuaian baru (`TreatyRevisionCopyAttachment` KINI hidup)

Di Pega, `Choose` picker menjalankan `TreatyInEDMSetValue`: [3] `OLDID = ID` (master yang dipilih),
[6] `TreatyInRevisi_post` (pengenal revisi baru), **[8] `TreatyRevisionCopyAttachment`** tanpa prasyarat,
[9] `SaveTreatyIn_EDM_Act`. Rantai salinannya (diverifikasi dari XML korpus `Treaty In Adjustment`):

| Langkah | Isi |
|---|---|
| [1] `CopyAllAttachment2_Sql` | `select ID, filename, CATEGORY, FILEMIMETYPE, T_STORAGE_ID, CATEGORY_ID from M_ATTACHMENTTREATY_2 where treatyid = {TreatyIn.OLDID}` |
| [2.3] `GetUrlGoogleStorage_Act` | URL bertanda tangan objek SUMBER (`ImageID = .type`, Durasi 1800; URL baru bila EXPDATE lewat) |
| [2.5] Java | unduh isi dari URL → Base64; galat ditelan (base64 kosong) |
| [2.7] `InsertGoogleStorage_Act` | dilewati bila base64 kosong; Ext = `.pyFileMimeType`, Folder `Contract`, Namafile = nama asli → objek BARU + `Insert_T_Storage_SQL` |
| [2.8] `InsertAttachment2_Sql` | dilewati bila ImageID kosong; TREATYID = pengenal BARU, CATEGORY/CATEGORY_ID/FILENAME/FILEMIMETYPE dari baris sumber apa adanya, USERNAME operator, T_STORAGE_ID = objek BARU, DATA_JSON tidak diisi |

Pega tidak punya pagar duplikat — setiap `Choose` menyalin lagi.

**Di aplikasi ini** (`modul/treatyin/backend/services/salin_lampiran_penyesuaian.go`):

- `Choose` hanya menyusun draf (keputusan pemilik proses), jadi salinan dibuat pada **tulisan pertama draf**
  (Save, Submit, atau Actions dengan `draf = true`), **sesudah** kepala `TREATY_IN_EDM` + pendaratan mengikat.
  Karena `SimpanPenyesuaian` menolak draf yang pengenalnya sudah ada, salinan berjalan **tepat sekali** per pengenal.
- Objek **diunggah ulang** seperti Pega, bukan baris yang menunjuk objek yang sama: tombol Delete menghapus objek
  Google Storage + baris `T_STORAGE_IMAGE` ber-IMAGEID itu, jadi berbagi objek akan merusak lampiran master.
- Tabel yang ditulis SAMA dengan tombol Upload file (`CatatLampiran`: `T_STORAGE_IMAGE` + `M_ATTACHMENTTREATY_2`,
  satu transaksi per berkas). Nol tabel/kolom baru.
- Penyimpanan tidak transaksional: gagalnya salinan **tidak** menggagalkan Save. Nasib tiap berkas dilaporkan di
  `salinanLampiran` hasil Save dan ditulis layar di bawah pesan prosedur (`komponen/salinanLampiran.ts`). INSERT yang gagal
  sesudah unggah berhasil meninggalkan objek yatim — sama dengan Upload file.
- Penyimpangan yang dinyatakan: berkas gagal dilaporkan (Pega menelannya); pagar `OLDID` (harus asal pengenal revisi
  menurut `IDRevisiBaru`, sebab draf datang dari layar); batas 25 MB seperti Upload file; konteks permintaan dilepas
  dari pembatalan supaya salinan tidak terpotong.
- Panel Attachment draf menampilkan lampiran master (`p.idAsal`); sesudah tersimpan layar dibuka ulang sebagai
  penyesuaian tersimpan dan panel membaca `p.id` — salinannya langsung tampil.
- ⚠️ Terbuka: panel Attachment modul ini baca-saja — berkas yang gagal disalin belum dapat diulang dari layar.

Dijaga `services/salin_lampiran_penyesuaian_test.go` dan `frontend/salinan-lampiran.test.ts`.

## 17 · Audit 8 Oktober 2026 — tombol, dropdown, grid panel New

Permintaan pemakai: *"pastikan tiap dropdown, table hingga button berfungsi semua dan kondisinya
diperhatikan … jangan semisal isi data table kosong tp bisa ditambah … samakan dengan yang ada di
pega … isi dari dropdown … sama seperti yang ada di treaty in"*. Inventaris dibuat dengan fungsi
layar itu sendiri atas setiap tab, include, rincian, dan harness panel New, dalam lima keadaan
(Revisi Material 1/2, EDMState 2, Adjust Premium, View).

**Tombol.** Setiap tombol yang TAMPIL berfungsi (tambah baris, `addRow`/`deleteRow`, atau rantai
rumus). Empat yang belum punya rumus (`Hide Facultative Share`, Submit/Decline non-EDM, Add Retro
Prop) tidak pernah tampil di Adjustment — syarat Pega-nya menyembunyikannya (`!TreatyMasterInEDM`,
`facsharedisp`, tab Retro disembunyikan keputusan pemilik proses). Material Type 2 mematikan
Add/Delete/Update persis `pyDisabledWhen` ekspor; tidak ada include yang dikunci di tingkat atas
(`pyEditOptions = Auto`).

**Dropdown — isinya kini sama dengan Treaty In:**

| Medan | Sumber | Sebelumnya |
|---|---|---|
| Layer · Part Of · Cover · Currency Relation · Note Reinstatement (rincian Layers) | `opsi-limits` Treaty In (`jenisLayer`, `cover`, `relasiMataUang`, `catatanReinstatement`) | kotak teks |
| Class of Business (CoBList, DetailLimits) | `kelas-bisnis?treatyGroupId=` — parameter `pTreatyGroupId` | kotak teks |
| Spreading Type (DetailShare) · Spreading Type XOL (Share) | `spreading-induk` — `TreatyGroupID` (`.TreatyGroupID` / `.TreatyGroupList(1).TreatyGroupID`) + `StartDate` | kotak teks |
| Bordereaux · Accounting Mode (kepala) | `opsi-kepala` Treaty In (label tampil, urutan Treaty In) | nilai mentah |

Pembangkit kini membawa parameter RD (`param`) dan medan ikut-isi (`setel`,
`pyAdditionalFields`/`pySetValueOnSelect`): memilih Currency mengisi `CurrencyID`, Treaty Group
`TreatyGroupID`, Treaty Type `TreatyTypeID`, Reinsurer `ReinsID`, Class of Business
`ClassOfBusinessID` — 76 sel/medan. Rincian menulis keduanya SEKALIGUS (`ubahMedanBanyak`).

**Grid kosong karena datanya tidak ada.** Penyesuaian (atau draf dari master) yang baris akarnya
tidak ada di `T_TREATY_REVISION` dibawa sebagai `terdarat = false`; panel New dan deret tombolnya
dibuka dalam mode LIHAT dengan keterangan — grid kosong tidak dapat ditambah lalu disimpan sebagai
data separuh. Di Pega keadaan ini tidak mungkin: penyesuaian selalu membawa salinan dokumen utuh.
Terukur: `1002305` terdarat, `1001540/R01` tidak.

**Masih terbuka (tab Retro disembunyikan):** perilaku `change` `CountRetroShare_Act`,
`Del/GetSpreadingRetro`, dan Add `AddFacRetroProp`.

Dijaga `frontend/pilihan-berparameter.test.ts`, `frontend/terkunci-pendaratan.test.ts`,
`backend/services/terdarat_test.go`, `backend/repository/terdarat_db_test.go`.
