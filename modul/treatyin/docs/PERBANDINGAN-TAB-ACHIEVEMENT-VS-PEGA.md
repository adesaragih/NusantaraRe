# Tab Achievement In IDR (prop) — aplikasi lawan ekspor Pega

Dibaca 7 Oktober 2026 dengan
[`alat-baca-ekspor`](../../../../../XML_NURE/_migration-docs/alat-baca-ekspor/),
bukan skrip sekali pakai.

```
Section/TreatyInTabsProportional.xml   nama tab @1581141
Section/TreatyInTabsAchievement.xml    cangkang tab (nyaris kosong)
Section/AchievementCombine.xml         grid + kaki        ← isinya
Activity/GetAchievement.xml            angkanya
```

---

## 1 · Hasil perbandingan

| | Pega | Aplikasi SEBELUM | Sekarang |
|---|---|---|---|
| Tab dirender | ✅ | ⛔ **nol cabang** — membukanya menampilkan tab pertama | ✅ |
| Grid 6 kolom | ✅ | ⛔ tidak ada | ✅ |
| Tiga sel kaki | ✅ | ⛔ tidak ada | ✅ |
| Baris ` Total In IDR` | ✅ | ⛔ tidak ada | ✅ |
| Rumus angkanya | ✅ | ✅ **sudah ada** (`hitung_achievement.go`) | ✅ dipakai ulang |

⭐ Rumusnya ternyata **sudah dibangun** dan dipakai sub-tab Achievement di
dalam Limits Prop. Yang hilang hanya **tampilan ringkasnya**, dan nama tabnya
sudah terdaftar di `TAB_PROPORSIONAL` tanpa ada yang merendernya.

### 1.1 Keenam kolom

| Judul (ejaan Pega) | Properti | Golongan |
|---|---|---|
| `Treaty Group` | `.TreatyGroup` | teks |
| `Reins Type` | `.TreatyType` | teks |
| `Gross Premium Before Claim in IDR` | `.TotalAchPremium` | uang |
| `Net Premium Before Claim in IDR` | `.TotalAchNetPremium` | uang |
| `Incured Claim in IDR` | `.TotalAchIncured` | uang |
| `Net Loss Ratio` | `.LossRatio` | persen |

⛔ **Keenamnya `pyReadOnly = true` TANPA `pyReadOnlyCondition`** — jadi
benar-benar hanya-baca. Diperiksa dengan `hanya_baca()`, bukan dengan
membaca tag telanjang; itu persis pembedaan yang salah dibaca di tab EGNPI
dan Maximum Retention.

### 1.2 Kaki — hanya TIGA sel

```
Treaty Group | Reins Type | Gross Premium | Net Premium | Incured | Loss Ratio
   (kosong)     (kosong)     (kosong)       SumNet       SumInc     SumLR
```

⛔ `Gross Premium` **nol punya total**. Ekspor menaruh `Spacers` di kolom
1–3, dan sel nilai hanya di bawah tiga kolom terakhir. Itu bunyi ekspornya,
bukan sel yang lupa dipasang.

---

## 2 · Asal penarikan datanya

```
Activity/GetAchievement.xml
  [3] RDB-List `GetAchievement`   InputParam.CARI1 = TreatyIn.ID
                                  InputParam.CARI2 = Treaty Group
  [4] RDB-List `GetQuarter`       daftar kuartal
  [5] RDB-List `GetQuarterYear`   daftar tahun kuartal
  [6] GetEstimasiClaim_act        HistoryClaimEstimasi.pxResults
  [7] GetHistoryClaim_act         HistoryClaim.pxResults
  [11.1.2.4.4] RDB-List kurs      SearchCurrencyValueOut.pxResults(1).CARI12
```

Pemetaan kolom RDB (langkah 11.1.2.2):

| | | | |
|---|---|---|---|
| CARI1 `IDPEGA` | CARI2 `NOOFFER` | CARI3 `SOBNAME` | CARI4 `TREATYGROUPNAME` |
| CARI5 `TREATYTYPE` | CARI6 `Quarter` (`"Q "+`) | CARI7 `QUARTERYEAR` | CARI8 `CurrencyID` |
| CARI9 `Currency` | CARI10 `PREMIUM` | CARI11 `RICOMM` | CARI12 `BROKERAGE` |
| CARI13 `NETPREMIUM` | CARI14 `PaidClaim` | CARI15 `OutstandingClaim` | CARI30 `NOPOLIS` |

### 2.1 Rumus per baris (11.1.2.4)

```
.Conversion = SearchCurrencyValueOut.pxResults(1).CARI12
              ⟶ 1 bila .Currency == "IDR"

bila .FlagTreaty != "1":
    .CASHCALL            = @if(ClaimIDR==0, 0, @divide(ClaimIDR, .Conversion, 20))
    .EstimationCASHCALL  = @if(EstIDR==0,   0, @divide(EstIDR,   .Conversion, 20))
    .OutstandingCASHCALL = .EstimationCASHCALL - .CASHCALL
bila .FlagTreaty == "1":  ketiganya 0

.IncuredClaim = .PaidClaim + .CASHCALL + .OutstandingClaim + .OutstandingCASHCALL
.Total        = .NETPREMIUM - .IncuredClaim
.LossRatio    = @if(.NETPREMIUM==0, 0, @divide(.IncuredClaim, .NETPREMIUM, 20)*100)
*toIDR        = * .Conversion
```

### 2.2 Rumus per Treaty Group (12.2.5) — yang MASUK grid

```
.TotalAchPremium    = Σ .PREMIUMtoIDR
.TotalAchNetPremium = Σ .NETPREMIUMtoIDR
.TotalAchIncured    = Σ .IncuredClaimtoIDR
.LossRatio          = @if(.TotalAchPremium=0, 0,
                          @divide(.TotalAchIncured, .TotalAchPremium, 20)*100)
```

### 2.3 Kaki (12.4)

```
SumTotalAchievNetPremium = Σ .TotalAchNetPremium
SumTotalAchievIncured    = Σ .TotalAchIncured
SumLossRatio             = @if(SumNet=0, 0, @divide(SumInc, SumNet, 20)*100)
```

---

## 3 · ⛔ Temuan

### 3.1 Kolom `Net Loss Ratio` memakai pembagi **GROSS**

Judulnya `Net Loss Ratio`, tetapi langkah 12.2.5 membagi dengan
`.TotalAchPremium` — **gross**, bukan `.TotalAchNetPremium`.

⭐ Rasio ber-pembagi Net **memang dihitung** di langkah yang sama
(`Local.LossRatioNet`), tetapi ia masuk ke `.CurrencyList` (panel
`Based on Nett`), bukan ke kolom ini.

⛔ **Tidak diam-diam diperbaiki.** Memperbaikinya membuat angka kita berbeda
dari layar lama tanpa seorang pun memutuskannya. Ini pertanyaan untuk pemilik
proses: apakah kolomnya salah judul, atau rumusnya salah pembagi?

### 3.2 Kaki memakai pembagi **NET** — tidak konsisten dengan barisnya

`SumLossRatio` membagi dengan `SumTotalAchievNetPremium`, sementara
`.LossRatio` tiap baris membagi dengan premi **gross**. Jadi angka kaki
**bukan** rekap dari kolom di atasnya — keduanya rasio dengan pembagi
berbeda.

### 3.3 `SumLossRatio` DITIMPA, bukan dijumlah

Langkah 12.2.7 menumpuk `Local.SumLossRatio + .LossRatio`, lalu 12.4
**menimpanya** dengan pembagian. Jadi penjumlahan itu terbuang.

⭐ Benar secara matematis — menjumlahkan persen tidak bermakna — tetapi
pembaca yang berhenti di 12.2.7 akan membangun jumlah yang salah. Dipaku
uji: `TestTotalLossRatioDihitungUlangBukanDijumlah`.

### 3.4 Label kaki `Total in IDR` TAMPIL walau ber-`pyCondition 1=2`

Selnya `pyVisible = ALWAYS`, dan **`pyCondition` hanya berlaku ketika
`pyVisible = OTHER`**. `Spacers` di sekitarnya ber-`pyVisible = OTHER` dan
karena itu memang tersembunyi; labelnya tidak.

⚠️ Aturan baca baru, belum ada di `alat-baca-ekspor`. Dicatat di §5.

### 3.5 Baris total memakai kolom `Reins Type`

Langkah 12.4 menulis `.Detail(<LAST>).TreatyType = " Total In IDR"` —
kolom **kedua**, bukan `Treaty Group`. Perhatikan spasi di depannya.

### 3.6 Tiga hal yang Activity lakukan dan aplikasi SENGAJA tidak

Tercatat sejak `hitung_achievement.go` ditulis, dan tetap berlaku:

1. **Klaim** (`GetHistoryClaim_act` / `GetEstimasiClaim_act`) — sumbernya
   `OS_AKSEPTASI_KLAIM.DATA_JSON`; nilai dari JSON **dilarang pemilik
   proses**. Cash Call = Estimation = 0.
2. **[12.3]/[12.4] mengubah pohon `Limits`.** Di aplikasi keduanya dibangun
   sebagai **proyeksi** (`RingkasAchievement`), bukan penulisan: menirunya
   apa adanya membuat tab Limits ikut berubah hanya karena tab Achievement
   dibuka.
3. **Achievement %** memerlukan `TreatyIn.RNMShareP`, yang belum punya tabel
   pendaratan. Kolom itu tidak ada di tab ini, jadi tidak menghalangi.

---

## 4 · Uji

`services/hitung_achievement_ringkas_test.go` — 8 uji memaku: keenam kolom,
Detail bergrup kosong dibuang, total Loss Ratio dihitung ulang (bukan
dijumlah), total Net/Incured dijumlah, baris total memakai `Reins Type`,
grid kosong nol baris total, Net nol tidak membagi nol, nol larik nil.

`frontend/achievement-tab.test.tsx` — 9 uji memaku bentuk layarnya,
termasuk kaki tiga sel dan sifat hanya-bacanya.

---

## 5 · Aturan baca baru untuk `alat-baca-ekspor`

⚠️ **`pyCondition` hanya berlaku ketika `pyVisible == "OTHER"`.** Sel
ber-`pyVisible = ALWAYS` tampil walau `pyCondition`-nya `1=2`.

Terukur di `AchievementCombine.xml`: `Spacers` (`OTHER` + `1=2`)
tersembunyi, `Total in IDR` (`ALWAYS` + `1=2`) tampil.

⛔ Pembaca yang menyaring semata-mata dengan `pyCondition` akan
menyembunyikan label yang Pega tampilkan — kekeliruan sejenis dengan
`pyReadOnly` tanpa `pyReadOnlyCondition`.
