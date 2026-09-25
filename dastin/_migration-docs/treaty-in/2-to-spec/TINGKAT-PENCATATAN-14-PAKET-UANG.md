# Tingkat pencatatan — 14 paket uang yang `SPEC-MODEL-DATA.md` §10.23a tinggalkan terbuka

**Tanggal:** 24 September 2026 · **Langkah 1 sesi to-spec**
**Perkakas:** `alat/sapu-tingkat-paket-uang.py` · **Sumber jalur:** kolom *Asal* §10, bukan ingatan

> ### BATAS KEMAMPUAN sapuan ini
>
> Ia dapat menunjukkan sebuah paket uang **DITULIS oleh ekspresi tertentu**, dan memperlihatkan
> ekspresinya apa adanya.
>
> Ia **tidak dapat menyatakan sebuah paket bertingkat `TREATY_100_PERSEN`.** Ketiadaan perkalian
> bukan bukti tingkat 100% — ia bisa berarti penulisnya **tidak terekspor** (`L-10`:
> `Declare Expression`, `Declare Trigger`), atau nilainya **diketik pengguna** dan tingkatnya adalah
> **kesepakatan**, bukan perhitungan. Yang begitu dilaporkan **BELUM**, bukan 100%.

## 0. Kalibrasi — dijalankan sebelum hasilnya dipakai

Satu kasus positif yang **sudah diketahui jawabannya** dimasukkan ke sapuan yang sama:
`NILAI_PENYEBARAN`, bertingkat `BAGIAN_NURE` menurut `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` §1.4.

```
FetchQSfromMaster.xml
   Primary.SpreadingList(<LAST>).Value = (.Pct/100) * primary.RNMShareList(1).Value
```

**LULUS** — perkalian terhadap `RNMShare` terdeteksi. Sapuan yang gagal pada kasus ini akan
melaporkan nol yang salah pada ketiga belas lainnya.

**Penolakan:** sasaran memuat operator (ekspresi `when`, bukan penugasan) — **0**;
berkas tidak terurai — **0**.

---

## 1. Dan satu koreksi terhadap aturan pemilahannya sendiri

Aturan langkah 1 berbunyi *"rumusnya mengalikan dengan `Share`/`Pct` maka `BAGIAN_NURE`"*.
**Diterapkan harfiah, ia menjawab salah pada tiga dari empat belas.**

| Dikalikan dengan | Contoh | Apakah itu bagian NuRe? |
|---|---|---|
| `RNMShare`, `ShareCalc`, `SplitRNMShare` | `.Value * @divide(Local.ShareCalc,100,4)` | **YA** — `Local.ShareCalc = TreatyIn.RNMShare` |
| **`RetentionPct`** | `.Value * @divide(Primary.RetentionPct,100,4)` | **TIDAK** — `RetentionPct = 100 - QSPct`, yaitu porsi **cedant**, pembagian di tingkat treaty |
| **`Surplus`** (jumlah *lines*) | `.Value * Primary.Surplus` | **TIDAK** — cacah, bukan porsi |
| **`MDPPct`, `MDPMinPct`, `PremiumReservePct`** | `.Value * local.mdprate` | **TIDAK** — tarif syarat kontrak |

> **Yang menentukan bukan ADANYA perkalian, melainkan APA yang dikalikan — dan basis gelungnya.**
> Setiap putusan di bawah menyebut keduanya.

---

## 2. Putusan — 8 dari 14 terbaca dari ekspor

### 2.1 `BAGIAN_NURE` — 3 paket, terbaca

| Paket | Entitas | Ekspresi penentunya | Basis gelung |
|---|---|---|---|
| **`PREMI_BRUTO`** | `BAGIAN` | `TreatyIn.Share(<LAST>).GrossPremiumList(<LAST>).Value = .Value * @divide(Local.ShareCalc,100,4)` — dan `Local.ShareCalc = TreatyIn.RNMShare` | `TreatyInNonAddItem` |
| **`PREMI_BRUTO_MINIMUM`** | `BAGIAN` | `…GrossPremiumMinList(<LAST>) = .Value * @divide(Local.ShareCalc,100,4)` — ekspresi sekeluarga, `ShareCalc` yang sama | `TreatyInNonAddItem` |
| **`NILAI_TERMIN`** | `TERMIN` | `.InstallmentList(<LAST>).Amount = @divide(Local.PctForCount,100,4) * Local.value` | **`TreatyInSetValueInstallment` langkah 6 ber-`pyStepsObjectName = TreatyIn.TotalShareNetNP`** — `Local.value` diambil dari **total share net**, yaitu bagian NuRe. Termin adalah yang ditagihkan NuRe |

> Varian `…GrossPremiumList(<LAST>).Value = .Value * @divide(Local.ShareCalc,100,4) * @divide(TreatyIn.ProRatePercent,100,4)`
> **tidak mengubah tingkat** — faktor kedua penyesuaian periode, bukan porsi.

### 2.2 `TREATY_100_PERSEN` — 5 paket, terbaca

| Paket | Entitas | Ekspresi penentunya | Kenapa 100% |
|---|---|---|---|
| **`NILAI_RETENSI`** | `RETENSI_CEDANT` | **dua penulis hidup, dua cabang:** QS → `Primary.RetentionList(<APPEND>).Value = .Value * @divide(Primary.RetentionPct,100,4)` dengan `.RetentionPct = 100 - .QSPct`; Surplus → `Primary.RetentionList(<LAST>).Value = .Value` (masukan manual, langkah 11 ber-obj `.RetentionList`) | keduanya membagi **risiko treaty**, bukan mengambil bagian NuRe. Sejalan `CONTEXT.md` §3.2: retensi di Surplus adalah **jumlah uang yang disepakati cedant** |
| **`KAPASITAS_SURPLUS`** | `DETAIL_PROPORSIONAL` | `Primary.IOOLimitList(<LAST>).Value = .Value * Primary.Surplus` | retensi × jumlah *lines* — persis definisi `CONTEXT.md` §3.2, dan basisnya retensi yang 100% |
| **`MDP`** | `LAYER` | `Primary.MDPList(<LAST>).Value = .Value * local.mdprate`, `local.mdprate = @divide(.MDPPct,100,20)` | **`DetailCalculation` langkah 5 ber-obj `TreatyIn.Limits`** — basisnya **limit layer**, yang §10.3 sudah tetapkan `TREATY_100_PERSEN`. Perkalian tarif tidak mengubah tingkat |
| **`MDP_MINIMUM`** | `LAYER` | `Primary.MDPMinList(<LAST>).Value = .Value * local.mdpminrate`, `local.mdpminrate = @if(.MDPMinPct==0,0,@divide(.MDPMinPct,100,20))` | basis gelung yang sama |
| **`CADANGAN_PREMI`** | `DETAIL_PROPORSIONAL` | `Primary.ReserveList(<LAST>).Value = .Value * Primary.PremiumReservePct / 100` | **`PremiumReserveCalculate` langkah 5 ber-obj `.CessionList`** — basisnya jumlah yang **diserahkan cedant** ke treaty, dan `CessionList` sendiri ditulis `= .Value * Primary.Surplus` pada cabang Surplus dan lewat `CessionPct = 100 - RetentionPct` pada QS. **Tidak ada perkalian bagian NuRe di sepanjang rantainya** |

> **Syarat pembalikan bersama untuk kelimanya:** rantai *"limit/retensi 100% maka turunannya 100%"*
> putus bila salah satu basisnya ternyata bukan 100%. Yang paling mungkin: **`MDP` punya penulis
> kedua**, `TreatyIn.Limits(<LAST>).MDPList(<APPEND>).Value = .MinimumDepositPremi` di
> `TreatyInMappingDataconvert` — **bawaan dari penawaran**, bukan hitungan. Ia **tidak** mengalikan
> bagian NuRe, jadi ia tidak membalik putusannya; tetapi bila besaran penawaran itu ternyata dicatat
> pada bagian NuRe, `MDP` menjadi **dua tingkat dalam satu daftar**. Ujinya `T-5` di §4.

---

## 3. Yang TIDAK terbaca — 6 paket, dilaporkan apa adanya

### 3.1 Dua yang ekspresinya ada tetapi tidak memisahkan apa pun

| Paket | Entitas | Yang terbaca | Kenapa tidak memutus |
|---|---|---|---|
| **`NILAI_EGNPI`** | `EGNPI` | `TreatyIn.EGNPI(<LAST>).Amount = .AmountEPI` (`TreatyInMappingDataconvert`) — **dibawa apa adanya dari penawaran**, tanpa perkalian apa pun | tingkatnya **diwarisi dari `EPI`**, dan `CONTEXT.md` §3.3 menyatakan tingkat `EPI` *"belum dipastikan"*. Pertanyaannya **bukan dua melainkan satu**: jawab `EPI`, dan `EGNPI` ikut terjawab |
| **`LIMIT_AGREGAT`** | `LAYER` | **tidak ada satu pun penulis yang menghitungnya.** Yang ada hanya penjumlah ringkasan: `.AggregateLimit = .AggregateLimit + local.aggregatelimit` di `TreatyInSummaryLimit`, `TreatyInSummaryLimitActual`, `TreatyInDifferenceLimitsSumary` | penjumlahan **lintas layer** adalah agregat, dan agregat tidak disimpan (§4.2). Nilai per layernya **diketik pengguna** |

> **Untuk `LIMIT_AGREGAT` ada alasan kuat menduga `TREATY_100_PERSEN`, dan ia tidak saya jadikan
> putusan:** ia duduk di entitas yang sama, di layar yang sama, dan dengan cara pengisian yang sama
> seperti `LIMIT` dan `DEDUCTIBLE` — keduanya sudah tertulis `TREATY_100_PERSEN` di §10.3. Tetapi
> kesamaan **cara pengisian** bukan kesamaan **tingkat**, dan menjadikannya putusan berarti menebak.
> **Ia ditanyakan**, dan pertanyaannya dibuat murah dengan menyertakan dugaan itu (`T-3`).

### 3.2 Empat yang NOL penulis — masukan layar, tingkatnya kesepakatan

| Paket | Jalur sumber | Penulis | Ada di layar? |
|---|---|---|---|
| **`BATAS_MAKSIMUM_KELOMPOK`** | `MaxCoGroup` | **0** | **ya** — 18 kemunculan di `Section`, 6 di `Harness` |
| **`BATAS_MAKSIMUM_NON_KELOMPOK`** | `MaxCoNonGroup` | **0** | **ya** — 18 `Section`, 6 `Harness` |
| **`BATAS_PILIHAN`** | `OptionLimit` | **0** | **ya** — 18 `Section`, 6 `Harness`, 6 `Activity` |
| **`NILAI_BATAS`** | `Earthquake` · `FloodJab` · `FloodNation` · `RSMDLimit` | **0** | **ya** — 79 · 54 · 27 · 27 di `Section` |

> **"Nol penulis" dibedakan dari "tidak ada sama sekali", dan pembedaan itu diperiksa.** Keempatnya
> **ada** dan ada **di layar** — jadi nilainya **diketik orang**, dan tingkatnya adalah apa yang
> disepakati di slip, bukan apa yang dihitung mesin. **Tidak ada ekspor yang dapat menjawabnya.**
>
> Satu kehati-hatian: ketiadaan penulis di sini juga konsisten dengan `L-10` — `Declare Expression`
> tidak terekspor sama sekali. **Kedua sebab menghasilkan nol yang sama**, dan sapuan ini tidak
> dapat memisahkannya. Bila jawaban `L-10` kelak **bukan nol**, keempat baris ini dibaca ulang.

---

## 4. Yang naik ke teknik treaty — lima pertanyaan untuk enam paket

| # | Pertanyaan, dalam bahasa bisnis | Menjawab |
|---|---|---|
| **T-1** | *"Batas maksimum per kelompok, batas maksimum non-kelompok, dan batas pilihan yang diketik di layar kontrak — angkanya untuk **seluruh treaty**, atau sudah bagian NuRe saja?"* | 3 paket sekaligus |
| **T-2** | *"Batas per bahaya — gempa, banjir Jabodetabek, banjir nasional, dan RSMD — untuk seluruh treaty atau bagian NuRe?"* | `NILAI_BATAS` |
| **T-3** | *"Batas agregat sebuah layer. Kami menduga ia setingkat dengan limit layer, yaitu seluruh treaty, karena diisi di kolom sebelahnya oleh orang yang sama. Benar?"* | `LIMIT_AGREGAT` — **berbentuk dugaan untuk dibantah**, supaya murah dijawab |
| **T-4** | *"Perkiraan pendapatan premi (`EPI`) yang dibawa dari penawaran — dicatat untuk seluruh treaty atau bagian NuRe?"* | `NILAI_EGNPI` **dan** butir §7 `EPI` sekaligus |
| **T-5** | *"Deposit premium yang dibawa dari penawaran — apakah angkanya setingkat dengan limit layer?"* | **uji pembalikan** putusan `MDP`, bukan pertanyaan baru |

> **Butir eskalasi `EGNPI` turun golongan — usulan, menunggu pemilik proses.** §7 menaikkannya ke
> eskalasi manajemen. Sesudah penulisnya terbaca — `= .AmountEPI`, dibawa apa adanya — ia bukan lagi
> pertanyaan tentang kewenangan, melainkan **satu pertanyaan faktual yang sama dengan pertanyaan
> `EPI` yang sudah ada di daftar**. **Usulan: dicabut dari daftar eskalasi, digabung ke `T-4`.**

---

## 5. Akibat pada constraint

| | Sebelum langkah ini | Sesudah |
|---|---:|---:|
| Paket uang bertingkat tertulis | 29 dari 43 | **37 dari 43** |
| Belum ditentukan | 14 | **6** |
| Di antaranya yang **tidak ada ekspor apa pun dapat menjawabnya** | — | **4** (§3.2) |

`INV-39` dan `INV-40` kini dapat ditulis untuk 37 paket. Untuk keenam sisanya, `ddl-usulan/`
**tidak memasang constraint** dan menulis **pernyataan keputusan** — bukan `TODO` — sesuai
larangan ketiga langkah 6.

## 6. Usulan diff ke `SPEC-MODEL-DATA.md` §10.23a — menunggu persetujuan

Berkas induk **tidak disunting oleh langkah ini.** Yang diusulkan: tabel §10.23a diganti sehingga
kolom **Tingkat** berbunyi `TREATY_100_PERSEN` untuk `NILAI_RETENSI`, `KAPASITAS_SURPLUS`, `MDP`,
`MDP_MINIMUM`, `CADANGAN_PREMI`; `BAGIAN_NURE` untuk `PREMI_BRUTO`, `PREMI_BRUTO_MINIMUM`,
`NILAI_TERMIN`; dan **`BELUM — pertanyaan T-n`** untuk keenam sisanya, masing-masing menyebut nomor
pertanyaannya supaya pembaca tahu di mana jawabannya akan datang.
