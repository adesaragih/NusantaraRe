# Sapuan 43 gambar `design treaty in.docx` — medan, tombol, dan sumbernya

**Dibaca 5 Oktober 2026**, ke-43 berkas di `D:\NUSANTARA RE APP\design-treaty-in-gambar\`.
Kontrak contohnya **`1001846`** (proporsional, `MACHINERY BREAKDOWN LINESLIP 2025`) dan
**`1001841`** (non-proporsional, `MOTOR VEHICLE COMBINED RISK & CATASTROPHE XOL 2025`).

⛔ Baris mana pun yang berbunyi "belum dibangun" menyebut **sebabnya**. Panel yang isinya
dikarang lebih buruk daripada panel yang dinyatakan belum punya sumber.

---

## 1 · Kepala kontrak — 16 medan, kedua cabang

| Medan | Gambar | Sumber | Status |
| --- | --- | --- | --- |
| `ID` | 01, 26 | kolom `TREATY_IN.ID` | ⭐ ada |
| `Reinsurance Type` | 01, 26 | kolom `PROPORTIONTYPE` | ⭐ ada |
| `Treaty Contract Name` | 01, 26 | kolom `TREATYCONTRACTNAME` | ⭐ ada |
| `Contract Ref No` | 01, 26 | `JSONDATA.ContractRefNo` | ⭐ ada, selalu tampil (ralat 7 Okt 2026) |
| `Teritorial Scope` | 01, 26 | kolom `TERITORIALSCOPE` | ⭐ ada |
| `Bordereaux` | 01 **saja** | `JSONDATA.Bordeaux` | ⭐ ada, **prop saja** |
| `Bordereaux Note` | 01, 26 | `JSONDATA.BordereauxNote` | ⭐ ada |
| `Commencement` | 01, 26 | kolom `COMMENCEMENT` | ⭐ ada, kini `dd/mm/yyyy` |
| `Termination` | 01, 26 | kolom `TERMINATION` | ⭐ ada, kini `dd/mm/yyyy` |
| `Treaty Year` | 01, 26 | kolom `TREATYYEAR` | ⭐ ada, **tidak diformat** |
| `Accounting Mode` | 01 (prop) | `JSONDATA.AccountingMode` | ⭐ ada, kini berlabel |
| `Accounting Mode` | 26 (non-prop) | `JSONDATA.AccountingModeNonProp` | ⭐ ada, kini berlabel |
| `Ceding` | 01, 26 | kolom `CEDING` + `CEDINGID` | ⭐ ada |
| `Source of Business` | 01, 26 | kolom `LEADINGREINSSOURCE` | ⭐ ada |
| `RNM as Treaty Leader` | 01, 26 | `JSONDATA.TreatyLeader` | ⭐ ada |
| `Existing Policy for Master ID` | 01, 26 | **`TREATYINPRODUCTION`** | ⭐ **BARU ronde ini** |

⭐ **Gambar 26 membuktikan `Bordereaux` memang proporsional-saja**: layar non-prop
memperlihatkan `Bordereaux Note` tanpa `Bordereaux` di atasnya. Itu padanan independen bagi
`pyCondition TreatyIn.ProportionType='Proportional'` @60.649 yang diukur ronde sebelumnya.

⭐ **Gambar 01 dan 26 juga menutup pertanyaan §6 ronde sebelumnya**: `RNM as Treaty Leader`
**TAMPIL** di keduanya (tak tercentang di 01, tercentang di 26), padahal rule `TreatyMasterInEDM`
salah pada setiap kontrak. Keputusan ronde itu untuk **tidak** menyembunyikannya terbukti benar.

---

## 2 · `Existing Policy for Master ID` — sumbernya KETEMU

| Langkah | Rule | Isi |
| --- | --- | --- |
| panel | `Section/InputTreatyInOffer.xml` @89.949 | `pyPageListProperty = PolisList.pxResults` @105.310 |
| pengisi | `Activity/FetchTreatyExistingProduction.xml` | `InputData.CARI1 := @if(TreatyIn.EDMState="", TreatyIn.ID, TreatyIn.OLDID)`; lalu `RDB-List`; lalu `.CARI2 := @substring(.CARI2,18)` |
| SQL | `RDBList/FetchTreatyInProductionUsingNooffer.xml` | lihat di bawah |

```sql
SELECT DISTINCT NOPOLIS AS CARI1, IDPEGA AS CARI2,
       QUARTER AS CARI3, QUARTER_YEAR AS CARI4
  FROM POOLDATA.TREATYINPRODUCTION
 WHERE SUBSTR(NOOFFER,1,7) = SUBSTR({InputData.CARI1},1,7)
 ORDER BY CARI4, CARI3 ASC
```

**Diadu dengan gambarnya, dan cocok persis:**

| Kontrak | Gambar | Oracle |
| --- | --- | --- |
| `1001841` | 1 baris, `RNM-QR.T02.05.2025.11987` / `NB-147044` | 1 baris, `NOPOLIS` sama, `IDPEGA = ASM-FW-GISFW-WORK NB-147044` |
| `1001846` | `No items` | 0 baris |

⭐ Prefiks `ASM-FW-GISFW-WORK ` tepat **18 aksara** — persis `@substring(.CARI2,18)`. Kecocokan
itu yang mengubah rantai ini dari dugaan menjadi bukti.

⚠️ `TREATYINPRODUCTION` **41.936 baris**, tabel warisan, **baca-saja**. Nol tulisan.

---

## 3 · Tab proporsional — 10 tampil dari 11

Gambar 01 memperlihatkan: Reporting Period · Portfolio · Limits · Share · Co-Ins Scale ·
Accumulation · Exclusions · Special Conditions · Information & Submit · Achievement In IDR.

⭐ **`Retro` tidak tampil**, dan syaratnya `TreatyIn.IsMultipleRetro` — nol dari 1.079 kontrak
proporsional bernilai `"true"`. Konsisten.

| Tab | Gambar | Isi yang terlihat | Status kita |
| --- | --- | --- | --- |
| Reporting Period | 01 | Start/End Date · Period · Submission/Confirmation/Settlement · grid 5 kolom | ⭐ ada |
| Portfolio | 02 | `Portfolio Type` · `Premium / Loss Type` · `Description` | ⭐ ada, **judul diralat** |
| Limits | 03–15 | **bersarang tiga tingkat + 11 sub-tab** | ⚠️ grid datar — §6 |
| Share | 16, 17 | panel `Total Share` + sub-tab `RNM Share` | ⭐ sub-tab dibangun |
| Co-Ins Scale | 18 | `Co-Insurance Share` · `% Treaty Limit` · 2 medan Max Panel | ⭐ ada |
| Accumulation | 19 | `Period` dropdown + grid 4 kolom | ⭐ ada, **judul diralat** |
| Exclusions | 20 | satu blok teks | ⭐ ada |
| Special Conditions | 21 | satu blok teks | ⭐ ada |
| Information & Submit | 22 | `Additional Information` · `Comment` · tombol `Decline offer` | ⭐ ada |
| Achievement In IDR | 23 | pohon `Kind of Treaty` → grid 6 kolom | ⚠️ belum dibangun |

---

## 4 · Tab non-proporsional — 10 tampil dari 11

Gambar 26 memperlihatkan: Maximum Retention · Event Limits · EGNPI · Limits · Share · Retro ·
Installment · Exclusions · Special Conditions · Information & Submit.

⭐ **Dua pembuktian sekaligus:**
1. `Value Difference` **tidak tampil** — syaratnya `EDMState != 3 && EDMMaterialType == 1`, dan
   `EDMMaterialType` nol di seluruh 1.854 dokumen. Konsisten.
2. `Retro` **TAMPIL** — di cabang ini ia memang tidak bersyarat. Pembetulan bercabang ronde
   sebelumnya (§20.1) terbukti benar; peta berkunci nama tab akan menyembunyikannya di sini.
3. `RNM Share` **bukan tab** — ia sub-tab di dalam `Share` (gambar 34).

| Tab | Gambar | Isi yang terlihat | Status kita |
| --- | --- | --- | --- |
| Maximum Retention | 26, 27 | grid 3 kolom + rincian + panel `Total Retention Amount` | ⭐ ada, **judul diralat** |
| Event Limits | 28 | RSMD · Earthquake · Flood (Jabodetabek) · Flood (Nationwide) | ⚠️ grid layer, bukan 4 medan |
| EGNPI | 29 | grid 6 kolom + rincian + 3 panel total + `Update Total` | ⭐ grid ada, **judul diralat**; panel total ⚠️ |
| Limits | 30–33 | grid layer + `Summary of Limit` + `Total All Layers` + `Total ROL` | ⭐ grid ada; panel ⚠️ |
| Share | 34, 35 | panel Share + sub-tab `RNM Share` + `Summarry of RNM Share` + `Total All Layers RNM Share` | ⭐ sub-tab dibangun; panel ⚠️ |
| Retro | 36, 37 | centang `Has Share To Retro` → 2 grid + 4 panel, **semuanya `No items`** | ⚠️ §17 — lihat §7 |
| Installment | 38 | `Installment` + `Update Value` + grid 6 kolom + `% Total` + `Total` | ⭐ grid ada; panel ⚠️ |
| Exclusions | 39 | satu blok teks | ⭐ ada |
| Special Conditions | 40 | satu blok teks | ⭐ ada |
| Information & Submit | 41 | sama dengan prop | ⭐ ada |

---

## 5 · Tab `Limits` proporsional — bersarang TIGA tingkat

⛔ **Empat belas gambar (03–15) untuk satu tab, dan sebabnya struktural.** Yang kita render hari
ini grid datar dari `M_TREATY_IN2`; yang Pega render:

```
Limits
└─ Kind of Treaty                    [Add]        ← grid
   └─ SPECIAL SURPLUS                [Delete]     ← baris yang dapat dibuka
      ├─ Treaty Type                 (dropdown)
      └─ Treaty Group                [Add]        ← grid di dalam grid
         └─ ENGINEERING              [Delete]     ← baris yang dapat dibuka
            ├─ Treaty Group          (dropdown)
            ├─ Class of Business     (grid, "No items")
            ├─ 100% Limit            [Add] {mata uang, nilai} [Remove]
            ├─ Retention             [Add] {mata uang, nilai} [Remove] + ☑ Auto calculate
            ├─ Cession to R/I        [Add] {mata uang, nilai} [Remove]
            └─ 11 SUB-TAB ↓
```

Kesebelas sub-tabnya — dan mereka **persis** kesebelas wadah `TABBED` di
`Section/DetailLimits.xml` yang diukur ronde sebelumnya:

| Sub-tab | Gambar | Medan yang terlihat |
| --- | --- | --- |
| Event Limits | 05 | RSMD Limit · Earthquake Limit · Flood Limit (Jabodetabek) · Flood Limit (Nationwide) |
| Deduction In A | 06 | `% OGR` · `% ONR` |
| Deduction | 07 | grid `Description` · `Currency` · `Deduction` · *or* · `Deduction %` + panel `Deductions` |
| Reserve | 08 | `% Premium Reserve` + `Premium Reserve` {mata uang, nilai} |
| Experience Premium Refund | 09 | `% Commision` · `% ME` · `YDCF` + catatan *ME = Management Expense*, *YDCF = Years Deficit Carried Forward* |
| PLA | 10 | `PLA` {Choose, nilai} |
| Cash Loss Limit | 11 | `Cash Loss Limit` {Choose, nilai} |
| Claim Cooperation | 12 | `Claim Cooperation` {Choose, nilai} |
| LPC | 13 | `Lower Band` · `Upper Band` · `Reisured Participant` · `Period (Month)` |
| EPI | 14 | `EPI` {mata uang, nilai} |
| Achievement | 15 | grid 11 kolom + grid `Parameter`/`Based on Gross`/`Based on Nett` + `As At Quarter` · `Quarter Year` + galat `EPI Must Not be Empty` + `Refresh` |

⛔ **Belum dibangun, dan bukan karena terlewat.** Struktur ini memerlukan entitas yang jalur baca
kita belum punya: `Kind of Treaty`, `Treaty Type`, dan `Class of Business` per treaty group.
`M_TREATY_IN2` memuat satu baris **datar** per layer; ia tidak memuat pohon ini. Membangun
bentuknya di atas data datar akan memberi pohon bercabang satu di mana Pega punya banyak.

---

## 6 · Tombol — sapuan ke-43 gambar

| Tombol | Gambar | Di layar kita |
| --- | --- | --- |
| `Add` / `Delete` (Rate of Exchange) | 16, 19, 20, 21, 23, 26, 28 | ada, **mati** — jalur tulis tiket 20 |
| `Apply` (Reporting Period) | 01 | ada, **mati**; galat `Start Date, Due. Must Not Be Empty` disalin utuh |
| `Choose Ceding` / `Choose Source of Business` | 01 | ada, **hidup** — hanya memilih |
| `add Layer` / `Delete` (Limits non-prop) | 30, 33 | ⚠️ belum — grid kita baca-saja |
| `Add` / `Delete` (Kind of Treaty, Treaty Group) | 03–15 | ⚠️ belum — §5 |
| `Add` / `Remove` (100% Limit, Retention, Cession to R/I) | 05–15 | ⚠️ belum — §5 |
| `Add Treaty Group` | 31 | ⚠️ belum — §5 |
| `Add MDP` / `Remove` | 32 | ⚠️ belum — §5 |
| `Update Total` (Maximum Retention) | — | ⭐ ada, **mati** |
| `Update Total` (EGNPI) | **29** | ⚠️ belum |
| `Update Total` (Limits) | **30, 33** | ⚠️ belum |
| `Update Total` (Share RNM) | **35** | ⚠️ belum |
| `Update Total` (Installment) | **38** | ⚠️ belum |
| `Update Value` (Installment) | **38** | ⚠️ belum — tombol yang ronde sebelumnya tidak temukan |
| `Update Summary` (Share non-prop) | **34** | ⚠️ belum |
| `Refresh` (Total Share prop) | **16** | ⚠️ belum |
| `Refresh` (Achievement) | **15** | ⚠️ belum |
| `Decline offer` | 22, 41 | ⚠️ belum |
| `Download All` (Attachment) | 24, 42 | ada, **mati** — `T_STORAGE_IMAGE` belum berjalur |
| `Refresh` (Attachment) | 24, 42 | ada, **hidup** — membaca ulang |
| `Upload file` (per kategori) | 24, 42 | ⭐ **BARU**, mati — modal `ASM Attach Content` di baliknya |
| `View File` (per kategori) | 25, 42 | ⭐ **BARU**, **hidup sebagai modal** — §4 |
| `Select file(s)` / `Attach` / `Cancel` (modal unggah) | 24 | ⚠️ belum — bagian jalur unggah |
| `Close` · `Actions` | 24, 42, 43 | ada, **mati** |
| `Save(dev)` · `Force Edit (dev)` | 24, 42, 43 | ⛔ **sengaja tidak dibangun** — tombol pengembang Pega, bukan fungsi bisnis |

⭐ **Sembilan tombol baru yang ke-43 gambar perlihatkan dan ekspor tidak**: `Update Value`,
`Update Summary`, `Refresh` (Total Share), `Refresh` (Achievement), `add Layer`,
`Add Treaty Group`, `Add MDP`, `Select file(s)`, `Attach`.

---

## 7 · ⚠️ Gambar 36/37 dan keputusan §17 — jawaban sebelum Retro disentuh

**§6 briefing meminta ini dijawab lebih dulu, dan jawabannya: TIDAK, gambar 36/37 tidak mengubah
dasar §17.**

§17 menyatakan Retro **tidak dibangun karena jarang**, bukan karena tidak ada spesifikasinya.
Yang kedua gambar itu tambahkan adalah **bentuknya**, bukan bukti pemakaiannya:

| Gambar | Isi |
| --- | --- |
| 36 | centang `Has Share To Retro:` **tidak tercentang** → nol isi dirender |
| 37 | centang **tercentang** → `Other Treaty Retro` grid, `Summarry of Other Treaty Retro`, `Total All Layers Other Treaty Retro` (4 panel) |

⛔ **Dan di gambar 37 SETIAP SATUNYA berbunyi `No items`.** Kontrak `1001841` — kontrak contoh
non-proporsional milik pemilik proses sendiri — punya nol baris Retro. Itu memperkuat §17,
bukan membantahnya.

⭐ Yang gambar 37 **tambahkan** ke pengetahuan kita, dan layak dicatat:

1. Isinya digerbangi satu centang `Has Share To Retro`, bukan oleh syarat tab.
2. Strukturnya sejajar dengan `Share`: grid layer → `Summarry of …` → `Total All Layers …`.
3. Ejaan ekspornya `Summarry` dengan **dua r** — sama dengan `Summarry of RNM Share`.

**Nol baris Retro ditulis ronde ini.**

---

## 8 · Angka — yang gambarnya perlihatkan berbeda dari satu aturan tunggal

⚠️ **Ini temuan yang membantah §7 briefing, dan ia terukur.** Desimal **tidak** satu aturan per
jenis; ia **per medan**, dan dua medan bernilai sama di baris yang sama dapat berbeda:

| Gambar | Medan | Tampil | Desimal |
| --- | --- | --- | ---: |
| 29 | EGNPI grid `Amount` | `137.849.315.068,00` | 2 |
| 29 | EGNPI grid `Amount in IDR` | `137.849.315.068` | **0** |
| 29 | EGNPI grid `Proportion %` | `100,00` | 2 |
| 29 | EGNPI rincian `Proportion %` | `100,0000000000` | **10** |
| 30 | Layers grid `100% Limits ( IDR )` | `750.000.000,00` | 2 |
| 30 | `Summary of Limit` `100% Limit (IDR)` | `750.000.000` | **0** |
| 32 | `Adjustment Rate %` | `0,367` | **3** |
| 32 | `ROL %` | `67,424` | **3** |
| 38 | `% Installment` | `25,00` | 2 |
| 38 | `% Total` | `100,0000` | **4** |
| 01 | `Value to IDR` | `1,00` dan `16.000,00` | 2 |

⛔ **Dan "nol di ekor dibuang" tidak berlaku di layar lama**: `1,00` tetap `1,00`, `5.345,00`
tetap `5.345,00`. Aturan kita membuangnya, jadi nilai bulat kita tampil `1` di tempat Pega
menulis `1,00`.

⚠️ **Nol baris aturan angka diubah ronde ini.** Mengubahnya menyentuh setiap grid di modul,
dan yang menentukan jumlah desimal di Pega adalah `pyDecimalPlaces` tiap kontrol — yang
**tidak ikut diekspor**. Yang dapat diputuskan dari gambar hanya bahwa aturan tunggal itu
salah, bukan aturan benarnya apa. Pertanyaannya di `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §9.

---

## 9 · Tanggal — dua bentuk, dan polanya lebih halus daripada "kepala lawan grid"

| Gambar | Tempat | Bentuk |
| --- | --- | --- |
| 01 | `Commencement`, `Termination` | `01/01/2025` |
| 01 | grid Rate of Exchange `Valid From` | `01/01/25` |
| 01 | grid Reporting Period `Initial Date` | `01/01/25` |
| 29 | grid EGNPI `As Date` | `18/01/25` |
| 29 | **rincian** EGNPI `As At` | `18/01/2025` |
| 38 | grid Installment `Due Date` (isian) | `18/01/2024` |
| 38 | grid Installment `Payment Date` (baca-saja) | `18/03/2024` |
| 43 | History `Date` | `09/04/2025 11:42` |

⭐ **Aturannya bukan "kepala lawan grid"** — Installment membantah itu. Yang terbaca:
**baris grid yang TERLIPAT** memakai dua digit tahun; **medan** — kepala maupun rincian yang
terbuka, dan kolom yang isinya medan isian — memakai empat.

Ronde ini memasang bentuk panjang pada `Commencement` dan `Termination`; grid tetap pendek.
Installment dan EGNPI rincian **belum** — keduanya memerlukan pembedaan baris terlipat dari
baris terbuka, dan grid kita belum punya baris yang dapat dibuka.
