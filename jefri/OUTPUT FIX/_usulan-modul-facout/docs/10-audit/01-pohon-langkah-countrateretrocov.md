# Audit 01 — Pohon langkah `CountRateRetroCov` — **DUA rule**

> **Diperluas 21 September 2026** setelah amandemen **K-060**: `CountRateRetroCov` adalah **dua rule
> berbeda di kelas berbeda**, dipilih menurut **COB**. Keduanya berlaku.
> Sumber angka bagi tiket `05-tickets\facout\F06-premi-retro.md`.

## 0. Dua rule, bukan dua salinan

⛔ **Identitas rule ditentukan `pzInsKey` + `pyClassName`, BUKAN nama berkas.** Dua berkas bernama
sama dengan basis `pzInsKey` berbeda adalah **DUA RULE**.

`[terverifikasi]`

| | Varian **FIRE** | Varian **ANEKA** |
| --- | --- | --- |
| `pyClassName` | `ASM-FW-GISFW-Data-PropertyItem` | `ASM-FW-GISFW-Data-Aneka` |
| Basis `pzInsKey` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-PROPERTYITEM COUNTRATERETROCOV` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-ANEKA COUNTRATERETROCOV` |
| **Berkas berlaku** | `DDL\CountRateRetroCov.xml` | `DDL\CountRateRetroCov(ANEKA).xml` |
| Salinan korpus | NB · RNW | Endorsment |
| Salinan korpus mutakhir? | ⛔ **BASI** (§4) | ✅ **MUTAKHIR** (§6.5) |

`[terverifikasi]` `pzOriginalInstanceKey` **keduanya sama** — varian Aneka lahir dari menyalin varian
PropertyItem lalu menyimpang.

⚠️ `pzInsKey` memuat **stempel waktu per-penyimpanan**; identitas dibaca dari **basisnya**. Kunci utuh
salinan `DDL\` FIRE bahkan berbeda dari salinan korpus FIRE (`#20260921…` vs `#20260619…`) meski
**rule-nya sama**.

---

**Sumber kebenaran varian FIRE:** `DDL\CountRateRetroCov.xml` (K-060, K-058).
**Pembanding:** `NB FacIn\Activity\CountRateRetroCov.xml` — salinan korpus, **BASI**.

---

## 1. Kedua salinan dan mengapa nomor versi tidak menolong

`[terverifikasi]`

| Salinan | Ukuran | `<pyRuleSetVersion>` | `<pxCommitDateTime>` | Status |
| --- | ---: | --- | --- | --- |
| `DDL\CountRateRetroCov.xml` | 131.407 B | **`01-01-01`** | **`20260921T095256 GMT`** | ✅ **berlaku** (K-060) |
| `NB FacIn\Activity\CountRateRetroCov.xml` | 131.835 B | `01-01-95` | `20260728T042513 GMT` | ⛔ basi |

⛔ **Nomor versi LEBIH RENDAH pada commit LEBIH BARU.** `<pyRuleSetVersion>` **tidak dapat dipakai
menilai kebaruan**. Ini bukti tambahan bagi pertanyaan terbuka di K-058 tentang metode deteksi drift.

---

## 2. Metode

Pohon `pySteps`/`rowdata` ditelusuri **rekursif** dengan `XmlDocument`. Alamat tiap langkah diambil
dari tag `<pyStepPageReference>`.

⛔ **Parser memisahkan anak LANGSUNG dari keturunan.** Tidak ada `SelectNodes('.//…')` di jalur
pembacaan langkah: `pyStepsPreCondParams`, `pyParamArray` dan `pySteps` dibaca **hanya** sebagai anak
langsung `rowdata` yang bersangkutan, dan isinya **hanya** sebagai anak langsung berikutnya. Penyapuan
keturunan akan mencampur langkah bersarang ke induknya.

`[terverifikasi]` **Arah tag**: `<pyStepPageReference>` selalu muncul **sesudah**
`<pyStepsActivityName>` (32/32 `rowdata` puncak pada tiga activity) — tetapi keanggotaan ditetapkan
**secara struktural**, bukan tekstual, sehingga parser tidak bergantung pada urutan sama sekali.
Rinciannya di `06-tinjau-ulang-suntingan-facout.md` §3.

---

## 3. Ukuran pohon — **identik di kedua salinan**

`[terverifikasi]`

| Ukuran | DDL | Korpus |
| --- | ---: | ---: |
| Total langkah | **14** | 14 |
| Puncak (kedalaman 0) | **10** | 10 |
| Bersarang | **4** | 4 |
| Indeks kedalaman maksimum | **1** | 1 |
| Jumlah **tingkat** | **2** | 2 |

⚠️ **Konvensi kedalaman.** Alamat terdalam `…pySteps(9).pySteps(2)` = dua ruas `pySteps` = **2
tingkat**, **indeks kedalaman maksimum 1**. Keterangan yang menyebut "kedalaman maksimum 2" memakai
konvensi *tingkat*; isinya sama. Dokumen ini memakai **indeks** dan menyebut keduanya.

---

## 4. ⛔ DUA perbedaan antara salinan berlaku dan salinan korpus — bukan satu

`[terverifikasi]` Diukur dengan parser anak-langsung, alamat dinormalkan (`RH_1`/`RH_2` → `RH`):

| Kelas | Jumlah |
| --- | ---: |
| Beda **metode** | **0** |
| Beda **gerbang** | **1** |
| Beda **WhenTrue/WhenFalse** | **1** |
| Beda **penugasan** | **1** |
| Alamat hilang / bertambah | **0** |

### 4.1 `pySteps(8)` — penugasan bertambah ✅ (sesuai keterangan work owner)

```
korpus : [1] Local.RateIN  := Local.RateIN + .Rate
         [2] Local.LengCov := .pxListSubscript

DDL    : [1] Local.RateIN  := Local.RateIN + .Rate
         [2] Local.LengCov := .pxListSubscript
         [3] Local.RIComIN := Local.RIComIN + .RIComm      <-- BERTAMBAH
```

Langkah 8: `Property-Set`, iterasi `.CoverageList`, **tanpa gerbang**, T/F = 2/2, **tanpa anak**.

### 4.2 ⛔ `pySteps(3)` — **gerbang `IsEDM` HILANG di ekspor baru**

**Perbedaan ini TIDAK disebut dalam keterangan yang diberikan**, dan ditemukan oleh pengukuran ini.

| | Korpus (basi) | **DDL (berlaku)** |
| --- | --- | --- |
| Baris precondition | **2** | **1** |
| Baris #1 | `IsEDM` — T=2 / F=3 | `IsEdmExtendPeriod` — T=2 / F=3 |
| Baris #2 | `IsEdmExtendPeriod` — T=2 / F=3 | *(tidak ada)* |
| Penugasan | `Local.prorate := 100` | `Local.prorate := 100` (sama) |

Hitungan seluruh pohon `[terverifikasi]`:

| | Korpus | DDL |
| --- | ---: | ---: |
| `rowdata` precondition | **15** | **14** |
| `<pyStepsPreCondParamsWhen>` tak kosong | **8** | **7** |
| Gerbang `IsEDM` persis | **3** | **2** |

⚠️ **`<pyStepsDescription>` langkah 3 tetap berbunyi `IsEDM IsEdmExtendPeriod` di KEDUA salinan.**
Label itu **basi** di salinan DDL — bukti lagi bahwa **label bukan bukti** (`PANDUAN-KERJA` §3).

⛔ **Akibatnya pada jalur uang:** di salinan berlaku, `Local.prorate := 100` dipaksakan pada kondisi
**lebih luas** — cukup `IsEdmExtendPeriod`, tanpa lagi menuntut `IsEDM`. Diport mengikuti DDL (K-060).

📌 **Awalan halaman harness berbeda** (`RH_1` korpus vs `RH_2` DDL). Itu **artefak ekspor**, bukan
perbedaan logika — dinormalkan sebelum dibandingkan.

---

## 5. Tabel langkah — **salinan DDL, berlaku**

| Alamat | Dlm | Metode | Iterasi | Gerbang | T | F | Penugasan |
| --- | ---: | --- | --- | --- | --- | --- | --- |
| `RH_2.pySteps(1)` | 0 | Property-Set | — | — | 2 | 2 | `[1]` `Local.prorate := EndPeriod − StartPeriod`<br>`[2]` `Local.prorate := @divide(Local.prorate, @if(EDMDay=="",365,@toDecimal(EDMDay)),20) * 100`<br>`[3]` `Local.prorate := 1` |
| `RH_2.pySteps(2)` | 0 | Property-Set | — | `IsEDM` | 2 | 3 | `[1]` `Local.prorate := (ProrateEDMEnd + ProrateStartEDM)*100` |
| `RH_2.pySteps(3)` | 0 | Property-Set | — | **`IsEdmExtendPeriod`** *(satu gerbang saja)* | 2 | 3 | `[1]` `Local.prorate := 100` |
| `RH_2.pySteps(4)` | 0 | Property-Set | — | `IsEDM` | **3** | **2** | `[1]` `Local.prorate := 100` |
| `RH_2.pySteps(5)` | 0 | Property-Set | — | — | 2 | 2 | `[1]` `Local.RateOut := .CoverageList(1).FacOutObjectList(1).Rate`<br>`[2]` `Local.RateIN := 0`<br>`[3]` `Local.ShareOffered := .CoverageList(1).FacOutObjectList(1).ShareOffered`<br>`[4]` `Local.TotalPremiCov := 0` |
| `RH_2.pySteps(6)` | 0 | Property-Set | — | `.Currency=="IDR" && …PctPremiAllObj != ""` | 2 | 3 | `[1]` `Local.ShareOffered := @divide((Local.ShareOffered*@toDecimal(PctPremiAllObj)),100,20)` |
| `RH_2.pySteps(7)` | 0 | Property-Set | — | `.Currency=="USD"&&…PctPremiAllObjUSD!=""` | 2 | 3 | `[1]` `Local.ShareOffered := @divide((Local.ShareOffered*@toDecimal(PctPremiAllObjUSD)),100,20)` |
| `RH_2.pySteps(8)` | 0 | Property-Set | `.CoverageList` | — | 2 | 2 | `[1]` `Local.RateIN := Local.RateIN + .Rate`<br>`[2]` `Local.LengCov := .pxListSubscript`<br>**`[3]` `Local.RIComIN := Local.RIComIN + .RIComm`** |
| `RH_2.pySteps(9)` | 0 | *(kosong)* | `.CoverageList` | — | 2 | 2 | — |
| `RH_2.pySteps(9).pySteps(1)` | 1 | Property-Set | — | `@Math.divide(Local.RateIN,1,10)==@Math.divide(Local.RateOut,1,10)` | **3** | **2** | `[1]` `.Rate := @Math.divide(Local.RateOut,Local.LengCov,20)` |
| `RH_2.pySteps(9).pySteps(2)` | 1 | Property-Set | — | — | 2 | 2 | `[1]` `.PremiumRetro := @Math.divide((Local.ShareOffered*.Rate*Local.prorate),100000,20)`<br>`[2]` `.PremiumRetro := .PremiumRetro − (@Math.divide((.PremiumRetro*.DiscountPercentage),100,20))`<br>`[3]` `Local.TotalPremiCov := Local.TotalPremiCov + .PremiumRetro` |
| `RH_2.pySteps(9).pySteps(3)` | 1 | Property-Set | — | — | 2 | 2 | `[1]` `.RICommPercentage := @Math.divide(Local.RIComIN,Local.LengCov,20)`<br>`[2]` `.RIComm := @Math.divide((.PremiumRetro*.RICommPercentage),100,20)` |
| `RH_2.pySteps(10)` | 0 | *(kosong)* | — | `@Math.divide(Local.TotalPremiCov,1,4)==@Math.divide(.CoverageList(1).FacOutObjectList(1).ObjectPremi,1,4)` | **3** | **2** | — |
| `RH_2.pySteps(10).pySteps(1)` | 1 | Property-Set | `.CoverageList` | — | 2 | 2 | `[1]` `.PremiumRetro := @Math.divide((@toDecimal(Primary.CoverageList(1).FacOutObjectList(1).ShareOffered)*.Rate*Local.prorate),100000,20)` |

⚠️ Dua langkah — `pySteps(9)` dan `pySteps(10)` — ber-`pyStepsActivityName` **kosong**. Pembacaan
regex polos atas tag itu menghasilkan **12**, bukan 14.

---

## 6. Jawaban atas empat pertanyaan

### 6.1 Urutan diskon terhadap `.RIComm` — **tidak berubah**

`[terverifikasi]` **Diskon dikurangkan lebih dulu; `.RIComm` dihitung sesudahnya dari `.PremiumRetro`
yang sudah terdiskon.** Diskon di `pySteps(9).pySteps(2)` `[2]`; komisi di `pySteps(9).pySteps(3)`
`[2]`. Keduanya **tanpa gerbang** (T/F = 2/2), jadi tidak ada percabangan yang dapat membalik urutan.

### 6.2 `Local.prorate` — **enam penugasan, tidak berubah**

> ⛔ **DIGANTIKAN §7D** (21 September 2026). Hitungan "enam penugasan di 4 alamat" tetap benar, tetapi
> bagian ini **tidak menyatakan nilai akhirnya**. Work owner menetapkan `Local.prorate = 1` keluar dari
> langkah 1 → dua rumus pertama **mati**. **Pakai §7D.** Bagian ini dipertahankan sebagai jejak.

`[terverifikasi]` Tetap **6**, di 4 alamat. Yang berubah hanya **pemilih langkah 3**:

| Alamat | Pemilih (DDL, berlaku) | Penugasan |
| --- | --- | --- |
| `pySteps(1)` `[1]` `[2]` `[3]` | **tanpa gerbang** | selisih tanggal → dibagi `EDMDay`/365 ×100 → **`1`** |
| `pySteps(2)` | `IsEDM`, T=2/F=3 | `(ProrateEDMEnd + ProrateStartEDM)*100` |
| `pySteps(3)` | **`IsEdmExtendPeriod` saja** ⛔ (korpus: `IsEDM` **dan** `IsEdmExtendPeriod`) | `100` |
| `pySteps(4)` | `IsEDM` **terbalik**, T=3/F=2 | `100` |

### 6.3 Pemilih dua rumus `.PremiumRetro` — **tidak berubah**

`[terverifikasi]` Gerbang `pySteps(10)`: kesetaraan **4 desimal** antara `Local.TotalPremiCov` dan
`FacOutObjectList(1).ObjectPremi`, T=3 / F=2. Rumus kedua di `pySteps(10).pySteps(1)` memakai
`Primary.…ShareOffered` **langsung**, sehingga **melewati** penskalaan mata uang `pySteps(6)`/`(7)`.

### 6.4 ⛔ `Local.RIComIN` — **keanehan GUGUR pada salinan yang berlaku**

`[terverifikasi]`

| Salinan | `<PropertiesName>` (menyetel) | `<PropertiesValue>` (membaca) |
| --- | ---: | ---: |
| Korpus (basi) | **0** | 1 |
| **DDL (berlaku)** | **1** | 2 |

> ⛔ **Temuan lama "dibaca tetapi tidak pernah disetel" DICABUT.** Itu **akibat salinan korpus yang
> basi**, **bukan** kode usang dan **bukan** kandidat perbaikan K-046.
>
> Pada salinan yang berlaku, `Local.RIComIN` **diakumulasi di `pySteps(8)`** bersama `Local.RateIN`,
> lalu dirata-ratakan di `pySteps(9).pySteps(3)`. Rumus
> `.RICommPercentage = @Math.divide(Local.RIComIN, Local.LengCov, 20)` menjadi **wajar dan sejajar**
> dengan `.Rate = @Math.divide(Local.RateOut, Local.LengCov, 20)`.

---

## 7. Yang TIDAK dapat ditentukan

1. ✅ **NILAI-nya DIJAWAB work owner 21 September 2026 — aturan umumnya TETAP TERBUKA.**
   ~~Apakah urutan `REPEATINGINDEX` di dalam satu `pyParamArray` adalah urutan eksekusi; bila ya,
   `Local.prorate := 1` di `pySteps(1)[3]` menimpa dua penugasan sebelumnya dan keduanya mati.~~
   Work owner menetapkan **`Local.prorate = 1` keluar dari langkah 1** → kedua rumus itu **memang
   mati** (§7D). ⚠️ Yang ditutup adalah **nilai pada rule ini**, **bukan** aturan umum
   `REPEATINGINDEX` = urutan eksekusi; aturan umum itu tetap **`[pertanyaan terbuka]` `[di luar
   korpus]`** dan **tidak ikut ditutup** di dokumen lain yang menyandarinya.
2. `[pertanyaan terbuka]` `[di luar korpus]` Apakah penomoran alamat `pySteps(n)` adalah urutan
   eksekusi. Seluruh kesimpulan urutan bersandar padanya.
3. `[pertanyaan terbuka]` **`Local.RIComIN` tidak pernah di-nol-kan.** `pySteps(5)` menginisialisasi
   `Local.RateIN := 0` dan `Local.TotalPremiCov := 0`, tetapi **tidak ada** `Local.RIComIN := 0`.
   Apakah variabel lokal Pega otomatis bernilai nol di awal activity adalah **`[di luar korpus]`**.
4. `[pertanyaan terbuka]` Arti kode transisi **1**, **4**, **5** — lihat `02-peta-kode-transisi.md`.
5. `[pertanyaan terbuka]` **Mengapa gerbang `IsEDM` hilang di `pySteps(3)`** — apakah penghapusan
   disengaja atau kekeliruan ekspor. Korpus tidak dapat menjawab; butuh UI Pega work owner.

---

## 7A. Varian **ANEKA** — pohon lengkap

**Sumber:** `D:\migrasi\RNM\DDL\CountRateRetroCov(ANEKA).xml` (115.823 B, ver `01-01-95`, commit
`20260831T101051 GMT`), kelas `ASM-FW-GISFW-Data-Aneka`.

`[terverifikasi]` **12 langkah — 9 puncak + 3 bersarang**, indeks kedalaman maksimum **1** (2 tingkat).

| Alamat | Dlm | Metode | Iterasi | Gerbang | T | F | Penugasan |
| --- | ---: | --- | --- | --- | --- | --- | --- |
| `RH_1.pySteps(1)` | 0 | Property-Set | — | — | 2 | 2 | `[1]` `Local.prorate := EndPeriod − StartPeriod`<br>`[2]` `Local.prorate := @divide(Local.prorate, @if(EDMDay=="",365,@toDecimal(EDMDay)),20) * 100`<br>`[3]` `Local.prorate := 1` |
| `RH_1.pySteps(2)` | 0 | Property-Set | — | `IsEDM` | 2 | 3 | `[1]` `Local.prorate := (ProrateEDMEnd + ProrateStartEDM)*100` |
| `RH_1.pySteps(3)` | 0 | Property-Set | — | `IsEdmExtendPeriod` | 2 | 3 | `[1]` `Local.prorate := 100` |
| `RH_1.pySteps(4)` | 0 | Property-Set | — | `IsEDM` | **3** | **2** | `[1]` `Local.prorate := 100` |
| `RH_1.pySteps(5)` | 0 | Property-Set | — | — | 2 | 2 | `[1]` `Local.RateOut := .CoverageList(1).FacOutObjectList(1).Rate`<br>`[2]` `Local.RateIN := 0`<br>`[3]` `Local.ShareOffered := .CoverageList(1).FacOutObjectList(1).ShareOffered`<br>⛔ **tanpa** `Local.TotalPremiCov := 0` |
| `RH_1.pySteps(6)` | 0 | Property-Set | — | ⛔ **`.Currency.Name=="IDR"`** `&& …PctPremiAllObj != ""` | 2 | 3 | `[1]` `Local.ShareOffered := @divide((Local.ShareOffered*@toDecimal(PctPremiAllObj)),100,20)` |
| `RH_1.pySteps(7)` | 0 | Property-Set | — | ⛔ **`.Currency.Name=="USD"`** `&&…PctPremiAllObjUSD!=""` | 2 | 3 | `[1]` `Local.ShareOffered := @divide((Local.ShareOffered*@toDecimal(PctPremiAllObjUSD)),100,20)` |
| `RH_1.pySteps(8)` | 0 | Property-Set | `.CoverageList` | — | 2 | 2 | `[1]` `Local.RateIN := Local.RateIN + .Rate`<br>`[2]` `Local.LengCov := .pxListSubscript`<br>⛔ **tanpa** akumulasi `Local.RIComIN` |
| `RH_1.pySteps(9)` | 0 | *(kosong)* | `.CoverageList` | — | 2 | 2 | — *(3 langkah anak)* |
| `RH_1.pySteps(9).pySteps(1)` | 1 | Property-Set | — | `@Math.divide(Local.RateIN,1,10)==@Math.divide(Local.RateOut,1,10)` | **3** | **2** | `[1]` `.Rate := @Math.divide(Local.RateOut,Local.LengCov,20)` |
| `RH_1.pySteps(9).pySteps(2)` | 1 | Property-Set | — | — | 2 | 2 | `[1]` `.PremiumRetro := @Math.divide((Local.ShareOffered*.Rate * Local.prorate),`**`10000`**`,20)`<br>`[2]` `.PremiumRetro := .PremiumRetro − (@Math.divide((.PremiumRetro*.DiscountPercentage),100,20))`<br>⛔ **tanpa** akumulasi `Local.TotalPremiCov` |
| `RH_1.pySteps(9).pySteps(3)` | 1 | Property-Set | — | — | 2 | 2 | `[1]` `.RICommPercentage := @Math.divide(Local.RIComIN,Local.LengCov,20)`<br>`[2]` `.RIComm := @Math.divide((.PremiumRetro*.RICommPercentage),100,20)` |

⛔ **Tidak ada `pySteps(10)`** — varian Aneka **tidak punya rumus koreksi**.

---

## 7B. Banding berdampingan kedua varian

`[terverifikasi]` H2 **tereproduksi seluruhnya**, dengan dua rincian tambahan (baris bertanda ✚):

| Aspek | **FIRE** (`Data-PropertyItem`) | **ANEKA** (`Data-Aneka`) |
| --- | --- | --- |
| Langkah | **14** | **12** |
| ✚ Puncak / bersarang | 10 / 4 | **9 / 3** |
| **Pembagi premi** | **100.000** (2 kemunculan) | ⛔ **10.000** (1 kemunculan) |
| **Uji mata uang** | `.Currency=="IDR"` · `=="USD"` | ⛔ `.Currency.Name=="IDR"` · `=="USD"` |
| `Local.RIComIN` **diakumulasi** | ✅ `pySteps(8)[3]` | ⛔ **tidak ada** |
| `Local.RIComIN` **dibaca** | ✅ `pySteps(9).pySteps(3)[1]` | ✅ **ada** — ⛔ **tanpa pernah diisi** |
| `Local.TotalPremiCov` | ada (2 penugasan) | ⛔ **0 penugasan** |
| Rumus koreksi `pySteps(10)`/`(10).pySteps(1)` | ✅ ada | ⛔ **tidak ada** |
| ✚ Penugasan `.PremiumRetro` | **3** | **2** |
| ✚ `pySteps(5)` jumlah penugasan | **4** | **3** |
| Gerbang `pySteps(2)` | `IsEDM` T=2/F=3 | **identik** |
| Gerbang `pySteps(3)` | `IsEdmExtendPeriod` T=2/F=3 | **identik** |
| Gerbang `pySteps(4)` | `IsEDM` T=3/F=2 | **identik** |

### ⛔ H3 tereproduksi — keanehan `RIComIN` hidup di varian ANEKA

`[terverifikasi]` Varian Aneka **membaca** `Local.RIComIN` tetapi **tidak pernah mengakumulasinya**:

```
RH_1.pySteps(9).pySteps(3)
  [1] .RICommPercentage := @Math.divide(Local.RIComIN,Local.LengCov,20)
  [2] .RIComm           := @Math.divide((.PremiumRetro*.RICommPercentage),100,20)
```

Di jalur **ANEKA**, komisi RI tetap dibagi dari variabel yang **tidak pernah diisi**. Keanehan itu
**MASIH ADA**, hanya **terbatas pada varian Aneka**. Pencabutan menyeluruh sebelumnya **dibatasi
ulang** ke varian FIRE.

### ⚠️ Satu implementasi tidak cukup

Jalur properti mata uang berbeda: `.Currency` vs `.Currency.Name`. Implementasi yang membaca mata
uang lewat **satu** jalur akan **salah pada salah satu varian** — gerbangnya tidak akan pernah benar,
sehingga penskalaan mata uang terlewat diam-diam.

---

## 7C. Salinan korpus varian ANEKA — **MUTAKHIR**

`[terverifikasi]` `DDL\CountRateRetroCov(ANEKA).xml` vs
`Endorsment Fac In\Activity\CountRateRetroCov.xml`, parser anak-langsung, alamat dinormalkan:

| Kelas perbedaan | Jumlah |
| --- | ---: |
| Metode | **0** |
| Gerbang | **0** |
| WhenTrue/WhenFalse | **0** |
| Penugasan | **0** |
| Alamat hilang/bertambah | **0** |

**Salinan korpus Endorsment sudah mutakhir** — berbeda dari varian FIRE, yang salinan korpusnya basi.

---

## 7D. `Local.prorate` — gerbang keempat langkah dan **nilai akhir per jalur**

> **Ditambahkan 21 September 2026** setelah keputusan work owner. **Menggantikan §6.2**, yang
> menyajikan keenam penugasan sebagai sama-sama mungkin. §6.2 **tidak dihapus** (`PANDUAN-KERJA` §7).

### 7D.1 Keputusan work owner

⛔ **21 September 2026:** pada `CountRateRetroCov`, `Local.prorate` **bernilai `1`** keluar dari
langkah 1. **Dua penugasan sebelumnya di langkah yang sama TIDAK berlaku.**

### 7D.2 Keempat langkah prorata — **IDENTIK di kedua varian**

`[terverifikasi]` Diukur ulang dengan parser anak-langsung atas **kedua** berkas `DDL\`. Gerbang,
T/F dan penugasan **sama persis**; satu-satunya perbedaan adalah awalan halaman harness
(`RH_2` FIRE vs `RH_1` ANEKA) — **artefak ekspor**, bukan logika.

| Langkah | Gerbang | T | F | Penugasan |
| --- | --- | ---: | ---: | --- |
| `pySteps(1)` | *(tanpa gerbang)* | 2 | 2 | `[1]` `Local.prorate := @String.toDate(…FacRetroList(Param.IdxFac).EndPeriod) − @String.toDate(…StartPeriod)` ⛔ **mati**<br>`[2]` `Local.prorate := @divide(Local.prorate, @if(…QuotationData.EDMDay=="",365,@toDecimal(EDMDay)), 20) * 100` ⛔ **mati**<br>`[3]` `Local.prorate := 1` ✅ **berlaku** |
| `pySteps(2)` | `IsEDM` | 2 | 3 | `Local.prorate := (…ProrateEDMEnd + …ProrateStartEDM) * 100` |
| `pySteps(3)` | `IsEdmExtendPeriod` | 2 | 3 | `Local.prorate := 100` |
| `pySteps(4)` | `IsEDM` | **3** | **2** | `Local.prorate := 100` — ⛔ **terbalik** |

`[terverifikasi]` Atribut `REPEATINGINDEX` pada ketiga `rowdata` `pyParamArray` langkah 1 = **1, 2, 3**;
yang ber-indeks **3** adalah `Local.prorate := 1`. ⚠️ `REPEATINGINDEX` adalah **atribut `rowdata`**,
bukan tag anak — membacanya sebagai tag anak menghasilkan nilai kosong.

### 7D.3 Nilai akhir per jalur

`[terverifikasi]` dari 7D.2, dengan penambatan kode transisi (`2` = jalankan lalu lanjut, `3` =
lewati langkah — **`[dugaan kuat]`**, `02-peta-kode-transisi.md` §2.2):

| Jalur | Langkah yang berjalan | **Nilai akhir `Local.prorate`** |
| --- | --- | --- |
| **NB / RNW** (non-EDM) | 1 → 4 | **`100`** |
| **EDM biasa** | 1 → 2 | **`(ProrateEDMEnd + ProrateStartEDM) × 100`** |
| **EDM extend period** | 1 → 2 → 3 | **`100`** |

📌 `[dugaan]` `<pyStepsDescription>` langkah 4 berbunyi **`IsNB`** di **kedua** varian — sejalan dengan
"jalan saat bukan EDM". **Label bukan bukti** (`PANDUAN-KERJA` §3): dicatat sebagai **penguat**
penambatan kode transisi, **tidak** dijadikan dasar. Statusnya tetap `[dugaan kuat]`.

⚠️ Deskripsi langkah **3** tetap berbunyi `IsEDM IsEdmExtendPeriod` padahal gerbangnya tinggal satu
(§4.2). Dua contoh label basi pada satu rule yang sama.

### 7D.4 ⛔ Dua rumus mati di langkah 1 — **tetap diport**

Kedua rumus pertama — selisih periode, lalu pembagian terhadap `EDMDay`/365 — **tidak berpengaruh
pada jalur mana pun**, karena penugasan ketiga menimpanya dengan `1` sebelum langkah berakhir.

**Sikap:** diport **apa adanya** (`CLAUDE.md` §1), ditandai **KODE MATI**, didaftarkan sebagai
**kandidat perbaikan K-046** — `K046_Prorate_DuaRumusMati_Langkah1` — bersama keanehan
`Local.RIComIN` varian Aneka (§7B). **Bukan** dihapus: menghapusnya adalah "perbaikan diam-diam"
yang dilarang `CLAUDE.md` §1.

⚠️ **Yang dijawab work owner adalah NILAI `Local.prorate` pada rule ini**, **bukan** aturan umum
bahwa urutan `REPEATINGINDEX` = urutan eksekusi. Aturan umum itu tetap **`[di luar korpus]`** dan
tetap terbuka di dokumen lain yang menyandarinya — lihat §7 butir 1.

### 7D.5 Pemeriksaan silang K-018 — **kedua pembagi benar**

`[terverifikasi]` Pada jalur non-EDM, dengan `Local.prorate = 100`:

| Varian | Kelas | Rumus efektif | Skala `.Rate` | Sesuai K-018? |
| --- | --- | --- | :-: | :-: |
| **FIRE** | `Data-PropertyItem` | `Share × Rate × 100 / 100000` = `Share × Rate / 1000` | **‰** | ✅ |
| **ANEKA** | `Data-Aneka` | `Share × Rate × 100 / 10000` = `Share × Rate / 100` | **%** | ✅ |

⛔ **Pembagi yang berbeda BUKAN kekeliruan ekspor.** Keduanya konsekuensi langsung **aturan pembagi
komposit K-018**: FIRE per mille (1.000 × 100), ANEKA persen (100 × 100). Ini menguatkan §0 — dua rule
kembar dengan pembagi berbeda **memang benar keduanya**.

⛔ **Pada jalur non-EDM, angka `100` itu BUKAN prorata** melainkan **konstanta skala** yang menyatu
dengan pembagi. Kosakatanya dikunci di `..\steering\GLOSARIUM.md` bagian *Kosakata yang dihindari*.

### 7D.6 Perintah audit

```powershell
# gerbang + penugasan empat langkah prorata, parser anak-langsung, kedua varian
function DC($n,$m){ foreach($c in $n.ChildNodes){ if($c.Name -eq $m){ return $c } }; return $null }
function DT($n,$m){ $c=DC $n $m; if($c -eq $null){ return '' }; return $c.InnerText.Trim() }
function RootSteps($p){ $x=New-Object Xml.XmlDocument; $x.Load($p)
  foreach($s in $x.SelectNodes('//pySteps')){ $d=0;$n=$s.ParentNode
    while($n -ne $null){ if($n.Name -eq 'pySteps'){$d++}; $n=$n.ParentNode }
    if($d -eq 0){ return $s } } }
foreach($f in @('D:\migrasi\RNM\DDL\CountRateRetroCov.xml',
                'D:\migrasi\RNM\DDL\CountRateRetroCov(ANEKA).xml')){
  $root=RootSteps $f; $i=0
  foreach($row in $root.ChildNodes){
    if($row.Name -ne 'rowdata'){ continue }
    $i++; if($i -gt 4){ break }
    $pc=DC $row 'pyStepsPreCondParams'
    foreach($pr in $pc.ChildNodes){ if($pr.Name -ne 'rowdata'){ continue }
      "$([IO.Path]::GetFileName($f)) L$i When=$(DT $pr 'pyStepsPreCondParamsWhen')" +
      " T=$(DT $pr 'pyStepsPreCondParamsWhenTrue') F=$(DT $pr 'pyStepsPreCondParamsWhenFalse')" }
    $pa=DC $row 'pyParamArray'
    foreach($pr in $pa.ChildNodes){ if($pr.Name -ne 'rowdata'){ continue }
      "   ri=$($pr.GetAttribute('REPEATINGINDEX')) $(DT $pr 'PropertiesName') := $(DT $pr 'PropertiesValue')" } } }
# -> L1 When=(kosong) T=2 F=2 ; ri=1,2,3 dengan ri=3 -> Local.prorate := 1
# -> L2 IsEDM T=2 F=3 ; L3 IsEdmExtendPeriod T=2 F=3 ; L4 IsEDM T=3 F=2   (identik kedua berkas)
```

---

## 8. ⛔ Tabel lama — **BASI, jangan dipakai**

> Versi dokumen ini sebelum 21 September 2026 memuat tabel langkah yang disusun di atas **salinan
> korpus** `NB FacIn\Activity\CountRateRetroCov.xml`, dengan:
>
> - ~~`RH_1.pySteps(3)` bergerbang `IsEDM` **dan** `IsEdmExtendPeriod` (T=2;2 / F=3;3)~~
> - ~~`RH_1.pySteps(8)` hanya dua penugasan, tanpa `Local.RIComIN`~~
> - ~~§4.4 "`Local.RIComIN` tidak pernah disetel — diport apa adanya, kandidat perbaikan"~~
>
> **Tidak dihapus** (`PANDUAN-KERJA` §7). **Penggantinya §4 dan §5 di atas**, di atas salinan `DDL\`
> yang ditetapkan berlaku oleh **K-060**. Salinan korpus tetap berguna **sebagai pembanding** untuk
> mengukur apa yang berubah — lihat §4.

---

## 9. Perintah audit

```powershell
# banding dua salinan, parser anak-langsung, alamat dinormalkan
function DC($n,$m){ foreach($c in $n.ChildNodes){ if($c.Name -eq $m){ return $c } }; return $null }
function DT($n,$m){ $c=DC $n $m; if($c -eq $null){ return '' }; return $c.InnerText.Trim() }
function Root($p){ $x=New-Object Xml.XmlDocument; $x.Load($p)
  foreach($s in $x.SelectNodes('//pySteps')){ $d=0;$n=$s.ParentNode
    while($n -ne $null){ if($n.Name -eq 'pySteps'){$d++}; $n=$n.ParentNode }
    if($d -eq 0){ return $s } } }
# -> beda metode=0  beda gerbang=1  beda T/F=1  beda penugasan=1  alamat hilang=0
```

```powershell
# jumlah baris precondition dan gerbang IsEDM di kedua salinan
foreach($p in @('D:\migrasi\RNM\NB FacIn\Activity\CountRateRetroCov.xml',
                'D:\migrasi\RNM\DDL\CountRateRetroCov.xml')){
  $t=[IO.File]::ReadAllText($p)
  $d=New-Object Xml.XmlDocument; $d.LoadXml($t)
  $n=($d.SelectNodes('//pyStepsPreCondParams/rowdata')).Count
  $w=([regex]::Matches($t,'<pyStepsPreCondParamsWhen>[^<]+</pyStepsPreCondParamsWhen>')).Count
  $i=([regex]::Matches($t,'<pyStepsPreCondParamsWhen>IsEDM</pyStepsPreCondParamsWhen>')).Count
  "$([IO.Path]::GetFileName($p)) : precond=$n When=$w IsEDM=$i" }
# korpus -> precond=15 When=8 IsEDM=3
# DDL    -> precond=14 When=7 IsEDM=2
```

```powershell
# Local.RIComIN : menyetel vs membaca, kedua salinan
foreach($p in @('D:\migrasi\RNM\NB FacIn\Activity\CountRateRetroCov.xml',
                'D:\migrasi\RNM\DDL\CountRateRetroCov.xml')){
  $t=[IO.File]::ReadAllText($p); $s=0;$b=0
  foreach($m in [regex]::Matches($t,'<PropertiesName>([^<]*)</PropertiesName>')){ if(([regex]'RIComIN').IsMatch($m.Groups[1].Value)){$s++} }
  foreach($m in [regex]::Matches($t,'<PropertiesValue>([^<]*)</PropertiesValue>')){ if(([regex]'RIComIN').IsMatch($m.Groups[1].Value)){$b++} }
  "$([IO.Path]::GetFileName($p)) : setel=$s baca=$b" }
# korpus -> setel=0 baca=1        DDL -> setel=1 baca=2
```

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
