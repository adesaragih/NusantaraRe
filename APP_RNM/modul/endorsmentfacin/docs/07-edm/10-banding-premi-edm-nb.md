# Banding rumus premi DASAR — EDM vs NB

> **Putaran verifikasi, bukan spec baru.** Dijalankan atas keputusan work owner: bandingkan
> **sekarang**, sebelum spec Seam 3 final, agar spec dan tiket tidak dibongkar ulang.
>
> **Sumber:** `Endorsment Fac In\` dan `NB FacIn\`. READ-ONLY. `RNM_BRD\` tidak dibaca (K-005).
> Label mengikuti `CLAUDE.md` §3.
>
> **Pertanyaan yang dijawab:** apakah perhitungan premi dasar di endorsement **sama** dengan NB
> (hanya dipecah per-lini) atau **berbeda secara fungsional** — dan apa dampaknya ke **Seam 3**.

---

## 0. Jawaban singkat

| Pertanyaan | Jawaban |
| --- | --- |
| Rumus premi **dasar** EDM sama dengan NB? | ✅ **SAMA — identik ekspresinya**, termasuk pembagi komposit `1.000.000.000` |
| EDM menambah sesuatu? | ✅ **Ya** — satu **varian dua-bagian** khusus endorsement |
| Seam 3 perlu seam baru? | ⛔ **Tidak** |
| Seam 3 perlu diperluas? | ✅ **Ya** — menerima **sepasang rasio prorata endorsement** dan menghasilkan penjumlahan dua bagian |

⚠️ **Tiga temuan yang mengubah gambaran** ditemukan sepanjang banding — §1.1, §4.2, §5.

---

## 1. Baseline NB — dan jebakan penamaan yang hampir menyesatkan

### 1.1 ⛔ `CalculatePremi_Act` (NB) BUKAN baseline premi

`[terverifikasi]` Penugasan menduga `NB FacIn\Activity\CalculatePremi_Act.xml` adalah pembanding
utama. **Dugaan itu tidak bertahan.**

| Ukuran | Nilai |
| --- | --- |
| `pxObjClass` | `Rule-Obj-Activity` |
| `pyClassName` | **`ASM-FW-GISFW-Data-PolicyTreatyIn`** ⚠️ — Treaty In, bukan Fac In |
| Cacah `.Rate` / `.TSI` / `.Premium` | **0 / 0 / 0** |

Isinya seluruhnya (3 langkah):

```
[1] SET .PremiOgp    = @divide(.GrossPremium * pyWorkPage.TreatyIn.RNMShareP, 100, 4)
    SET .ResultOgp1  = @divide(.PremiOgp * pyWorkPage.TreatyIn.BrokeragePercentP, 100, 4)
    SET .RiCommOgp   = @divide(.ResultOgp1, .PremiOgp, 4) * 100
[2] SET .Claim       = @divide(.GrossClaim * pyWorkPage.TreatyIn.RNMShareP, 100, 4)
[3] Call CountNetPremi_act
```

⛔ Ini **pembagian share treaty**, bukan perhitungan premi dari rate × TSI. Kelasnya
`…Data-PolicyTreatyIn` — domain **Treaty**, bukan Facultative Inward.

📌 **Jebakan penamaan ke-delapan di proyek ini:** nama berkas menjanjikan "hitung premi", isinya
tidak menyentuh rate maupun TSI sama sekali. Ditemukan hanya karena isinya dibaca.

### 1.2 Baseline yang sebenarnya — ditemukan dari isi

`[terverifikasi]` Activity NB dengan sinyal `.Rate` + `.TSI` tertinggi:

| Berkas NB | `.Rate` | `.TSI` | Ukuran |
| --- | ---: | ---: | ---: |
| `ChangeUpKendaraan` | 78 | 62 | 235 KB |
| `CountGrossPremiEDM_Act` | 56 | 70 | 704 KB |
| `FillPremiGolf` | 56 | 56 | 152 KB |
| `CountPremiCoverageAneka` | 54 | 54 | 260 KB |
| `fillPremiAneka` | 36 | 36 | 260 KB |
| `SumTotalTSIPremiGross_Act` | 31 | 67 | 436 KB |
| `CountPremi_ACT` | 30 | 93 | 522 KB |

⚠️ `CountGrossPremiEDM_Act` **bernama "EDM" tetapi ada di NB** — penegasan lagi bahwa nama bukan
bukti milik-siklus.

---

## 2. ⛔ Temuan pokok: EDM memakai ulang mesin premi NB

`[terverifikasi]` Activity ber-nama `premi`/`premium` di kedua korpus:

| | Jumlah |
| --- | ---: |
| NB | **44** |
| EDM | **47** |
| **Ada di KEDUANYA (nama sama)** | **41** |
| Hanya di EDM | **7** |
| Hanya di NB | **4** |

⛔ **Empat puluh satu activity perhitungan premi ada di kedua siklus** — termasuk **seluruh
kalkulator per-lini**:

```
CountPremi_ACT · FillPremiGolf · CountPremiCoverageAneka · FillPremiMBU_FacIn
CountGrossPremi_Act · CountGrossPremiEDM_Act · SumTotalTSIPremiGross_Act
CountPremiAndTSIRNM{Cargo,FireAnekaGolf,FireMBU,PALifeCargo}_ACT
SumTSIPremiSpreadedRNM_{FIRE,ANEKA,GOLF,MARINECARGO,MBU,PA,TRAVEL}_Act
CountPremiumNet · CountPremiumNetElse · CalculatePremiPA_FacIn · HitungPremiAndTSINusantaraRe_ACT
```

📌 **Hipotesis "endorsement punya mesin premi sendiri" TIDAK bertahan.** EDM **memakai ulang** mesin
premi NB. Yang EDM-only adalah **tambahan**, bukan pengganti.

### 2.1 Tujuh yang hanya di EDM

`[terverifikasi]`

| Berkas | `.Rate` | `.TSI` | Ukuran | Sifat |
| --- | ---: | ---: | ---: | --- |
| **`CalculatePremiFire`** | 17 | 18 | 153 KB | Ⓟ perhitungan premi dasar |
| **`HitungPremiOnChange`** | 14 | 14 | 56 KB | Ⓟ perhitungan premi dasar |
| **`calculatePremiPA`** | 7 | 3 | 234 KB | Ⓟ |
| **`CalculatePremiumTravel`** | 0 | 8 | 222 KB | Ⓟ |
| `ChangeCoveragePremium_FacIn` | 0 | 0 | 217 KB | pendukung layar |
| `fillActPremiPA` | 0 | 0 | 60 KB | pendukung layar |
| `ReCountPremiLifeEDM` | 4 | 11 | 171 KB | → **modul 6 (Life)** |

### 2.2 Keempatnya dipanggil dari **Section**, bukan activity

`[terverifikasi]`

| Berkas | Pemanggil |
| --- | --- |
| `CalculatePremiFire` | `Section\InputCoverageFire_IsUW` · `Section\InputPerCoverageFire_IsUW` · `Section\ViewCoverageFire` |
| `HitungPremiOnChange` | `Section\InputPerCoverageFire_IsUW` · `Section\ViewCoverageFire` |
| `calculatePremiPA` | `Section\ViewCoveragePA` |
| `CalculatePremiumTravel` | `Section\ViewCoverageTravel` |

⛔ **Keempatnya dipanggil dari lapisan tampilan** — hitung-ulang saat pengguna mengubah nilai di
layar endorsement. Bukan bagian rantai produksi.

---

## 3. Rumus premi dasar — SAMA

`[terverifikasi]` `Endorsment Fac In\Activity\CalculatePremiFire.xml` **L2003** dan
`HitungPremiOnChange.xml` **L645** — bentuk **identik**:

```
.Premium = @Math.divide(
             .TSI × .Rate × .ProRatePercent
                  × @if(.FirstLossScale=="", 100, .FirstLossScale)
                  × @if(.IndemnityPercentage=="", 100, .IndemnityPercentage),
             1000000000, 4)
```

### 3.1 Pembagi komposit — konfirmasi K-018

`[terverifikasi]` Pembagi **`1.000.000.000`** adalah **hasil kali kontribusi satuan per faktor**:

| Faktor | Satuan | Kontribusi |
| --- | --- | ---: |
| `.Rate` | per-mil (‰) | **1.000** |
| `.ProRatePercent` | persen | 100 |
| `.FirstLossScale` | persen | 100 |
| `.IndemnityPercentage` | persen | 100 |
| | **hasil kali** | **10⁹** |

⛔ **Ini persis aturan pembagi komposit K-018** — dan berlaku **sama** di EDM. Hanya faktor yang
melekat pada `.Rate` yang menentukan satuan rate; sisanya persen.

⛔ **Pembulatan 4 desimal literal per langkah** (`@Math.divide(…, 4)`) — konsisten dengan K-011.

### 3.2 Varian rupiah

`[terverifikasi]` `.PremiRp` = rumus yang sama **dikali kurs**:

```
.PremiRp = @Math.divide(…, 1000000000, 4) × @if(pyWorkPage.Policy.CurrencyValue=="", 1, …CurrencyValue)
```

⚠️ Perhatikan: di `.PremiRp`, `FirstLossScale` **tidak** dibungkus `@if(…=="",100,…)` seperti di
`.Premium` — bila kosong, hasilnya nol, bukan 100. **Diport apa adanya (K-046).**

---

## 4. Yang EDM TAMBAHKAN — varian dua-bagian

### 4.1 Bentuknya

`[terverifikasi]` `CalculatePremiFire` **L2563** dan **L2810**, `HitungPremiOnChange` **L383**:

```
.Premium = ( TSI × Rate × ProratePercentEDMEnd × FLS × IP  / 1e9 )      ← bagian BARU
         + ( OldCoverage(1).TSI × OldCoverage(1).Rate
                  × ProratePercentStartEDM × FLS_lama × IP_lama / 1e9 ) ← bagian LAMA
```

⛔ **Premi endorsement = porsi coverage BARU untuk sisa periode + porsi coverage LAMA untuk periode
yang sudah berjalan.** Kedua bagian memakai **rumus dasar yang sama**; yang berbeda hanya **rasio
prorata** dan **sumber nilainya** (`OldCoverage(1)` = before-image tingkat coverage).

`[terverifikasi]` `CalculatePremiFire` L282 menunjukkan bagian lama dapat berdiri sendiri:

```
.Premium = @Math.divide(.OldCoverage(1).TSI × .OldCoverage(1).Rate × .OldCoverage(1).ProRatePercent
             × FLS_lama × IP_lama, 1000000000, 4)
```

### 4.2 ⛔ Temuan: DUA keluarga properti prorata yang namanya nyaris sama

`[terverifikasi]` Sebaran di korpus EDM:

| Properti | Kemunculan | Berkas | Dipakai di |
| --- | ---: | ---: | --- |
| **`ProratePercentEDMEnd`** | 35 | **5** | perhitungan premi **dasar** (modul ini) |
| **`ProratePercentStartEDM`** | 37 | **5** | idem |
| `ProrateEDMEnd` | 101 | 15 | **delta spreading** (modul 2 / E-5) |
| `ProrateStartEDM` | 36 | 14 | idem |
| `ProRatePercent` | 511 | 59 | rumus dasar non-endorsement |

⛔ **`ProratePercentEDMEnd` ≠ `ProrateEDMEnd`.** Dua properti berbeda, dipakai di dua tingkat
perhitungan berbeda, dengan nama yang berbeda hanya pada sisipan `Percent`.

📌 **Jebakan penamaan ke-sembilan.** Menyatukannya akan mencampur premi dasar dengan delta spreading
— dua hal yang §5 tegaskan harus dipisah.

⚠️ `[pertanyaan terbuka]` apakah nilainya selalu sama atau bisa berbeda. Keduanya **tidak** diisi oleh
rule yang sama: `ProrateEDMEnd` diisi `SetValueToEDMWork` langkah 15 dan `CountPremiEDM_DT` (K-048
§8.1); pengisi `ProratePercentEDMEnd` **belum ditelusuri**.

### 4.3 ⚠️ Dugaan pembalikan pada `.PremiRp` — diport apa adanya

`[terverifikasi]` Pada varian dua-bagian, `.Premium` dan `.PremiRp` memakai rasio yang **tertukar**:

| Baris | Properti | Bagian BARU pakai | Bagian LAMA pakai |
| --- | --- | --- | --- |
| `CalculatePremiFire` L2810 | `.Premium` | `ProratePercentEDMEnd` | `ProratePercentStartEDM` |
| `CalculatePremiFire` **L2856** | `.PremiRp` | ⚠️ **`ProratePercentStartEDM`** | ⚠️ **`ProratePercentEDMEnd`** |
| `HitungPremiOnChange` L383 | `.Premium` | `ProratePercentEDMEnd` | `ProratePercentStartEDM` |
| `HitungPremiOnChange` **L435** | `.PremiRp` | ⚠️ **`ProratePercentStartEDM`** | ⚠️ **`ProratePercentEDMEnd`** |

⛔ **Tampak seperti pembalikan, konsisten di dua berkas.** Karena konsisten, bisa jadi disengaja —
tetapi korpus tidak menjelaskannya.

⛔ **Diport apa adanya (K-046).** Kandidat perbaikan milik bisnis. Wajib jadi **kasus uji bernama**:
`K046_PremiRp_RasioProrataTertukar_TerhadapPremium`.

⚠️ `[pertanyaan terbuka]` — milik work owner + Aktuaria.

---

## 5. Pemisahan tegas: premi DASAR vs delta SPREADING

⛔ **Dua perhitungan berbeda, jangan dicampur:**

| | **Premi DASAR** | **Delta SPREADING** |
| --- | --- | --- |
| Berkas | `CalculatePremiFire` · `HitungPremiOnChange` · `calculatePremiPA` · `CalculatePremiumTravel` | `SaveFacinProdEDM*_Act` |
| Rumus | `TSI × Rate × ProRatePercent × FLS × IP / 1e9` | `(nilai_baru × ProrateEDMEnd) − nilai_lama_TreatyType_sama` |
| Masukan | `.TSI`, `.Rate` per coverage | `.TSISpreaded`, `.PremiumSpreaded` per baris spreading |
| Rasio | `ProratePercentEDMEnd` / `ProratePercentStartEDM` | `ProrateEDMEnd` / `ProrateStartEDM` |
| Keluaran | `.Premium`, `.PremiRp` per coverage | kolom `_MENJADI` / `_SELISIH` `facinproduction` |
| Dipanggil dari | **Section** (layar) | rantai produksi |
| Modul | **modul premi** (Seam 3) | **modul 2 / E-5** |

`[terverifikasi]` `CalculatePremiFire` memuat **0 kemunculan** `TSISpreaded`, `PremiumSpreaded`,
`ProrateEDMEnd`, dan `OldData`. Konfirmasi: ia **tidak menyentuh** jalur delta spreading sama sekali.

---

## 6. ⚠️ Temuan tambahan: salinan EDM dari mesin premi bersama **berbeda isinya**

`[terverifikasi]` Banding **metode 23-tag + `pzIndexes`** (kontrak K-043) atas enam activity premi
bernama sama:

| Berkas | NB vs EDM |
| --- | --- |
| `CountPremi_ACT` | **BERBEDA** |
| `FillPremiGolf` | **BERBEDA** |
| `CountPremiCoverageAneka` | **BERBEDA** |
| `CountGrossPremiEDM_Act` | **BERBEDA** |
| `SumTotalTSIPremiGross_Act` | **BERBEDA** |
| `FillPremiMBU_FacIn` | ✅ identik fungsional |

⛔ **"Nama sama" ≠ "dapat dipakai ulang begitu saja".** Lima dari enam salinan EDM sudah dimodifikasi.
Isi perbedaannya **belum dibedah** — di luar lingkup putaran ini.

📌 Ini **tidak membatalkan** kesimpulan §3: rumus **dasarnya** tetap sama (terbukti dari ekspresi
yang dikutip). Yang berbeda ada di tempat lain dalam berkas-berkas itu.

⚠️ `[pertanyaan terbuka]` — perbedaan isi kelima berkas itu perlu dibedah **sebelum implementasi
modul premi**, bukan sebelum spec.

---

## 7. Dampak ke Seam 3 — diperluas, bukan ditambah

### 7.1 Kesimpulan per lini

| Lini | Rumus dasar | Varian endorsement | Status |
| --- | --- | --- | --- |
| **FIRE** | ✅ **SAMA** — L2003 identik bentuk NB | ✅ ada — L2563/L2810 dua-bagian | terverifikasi penuh |
| **PA** | `[dugaan]` sama — `calculatePremiPA` rate 7 / TSI 3, belum dikutip rumusnya | belum diperiksa | **perlu banding lanjutan** |
| **TRAVEL** | `[dugaan]` — `CalculatePremiumTravel` **`.Rate` = 0**, jadi rumusnya **bukan** rate × TSI | belum diperiksa | **perlu banding lanjutan** |
| **Lini lain** (Aneka, Golf, MBU, Marine Cargo) | ✅ memakai activity **bersama** dengan NB | — | tidak ada activity premi EDM-only |

⚠️ **Hanya FIRE yang terverifikasi penuh** di putaran ini. PA dan TRAVEL `[dugaan]` — ⚠️ `.Rate` = 0
pada Travel adalah sinyal bahwa rumusnya **berbeda bentuk**, bukan sekadar pecahan per-lini.

### 7.2 Bentuk perluasan Seam 3

⛔ **Tidak ada seam baru.** Total tetap **5**.

`services/premium.Calculate` diperluas menerima **konteks endorsement opsional**:

| Masukan tambahan | Isi |
| --- | --- |
| Sepasang rasio prorata premi | `ProratePercentEDMEnd` dan `ProratePercentStartEDM` — ⚠️ **bukan** `ProrateEDMEnd`/`ProrateStartEDM` |
| Nilai coverage lama | `OldCoverage(1).{TSI, Rate, ProRatePercent, FirstLossScale, IndemnityPercentage}` |

Perilakunya:

- **Tanpa konteks endorsement** → rumus dasar tunggal (jalur NB/RNW) — tidak berubah.
- **Dengan konteks endorsement** → penjumlahan dua bagian §4.1, keduanya memakai rumus dasar yang
  sama.

⛔ **Rumus dasar, pembagi komposit K-018, dan pembulatan 4 desimal tidak berubah.** Yang ditambahkan
hanya **cara menyusun dua panggilan rumus dasar menjadi satu hasil**.

📌 Keuntungannya: rumus dasar tetap **satu implementasi**, diuji sekali, dipakai NB, RNW, dan EDM.
Varian endorsement menjadi **komposisi**, bukan cabang di dalam rumus.

---

## 8. Kejanggalan — diport apa adanya (K-046)

| # | Butir | Nama test yang disarankan |
| ---: | --- | --- |
| 1 | **`.PremiRp` memakai rasio prorata tertukar** terhadap `.Premium`, konsisten di 2 berkas (§4.3) | `K046_PremiRp_RasioProrataTertukar` |
| 2 | **`FirstLossScale` tidak dibungkus `@if(…=="",100,…)`** di `.PremiRp`, padahal dibungkus di `.Premium` — kosong menghasilkan nol, bukan 100 (§3.2) | `K046_PremiRp_FirstLossScaleTanpaGuard` |
| 3 | **`CalculatePremi_Act` (NB) bernama "hitung premi" tetapi domain Treaty In** (§1.1) | — catatan penamaan |
| 4 | **`CountGrossPremiEDM_Act` bernama "EDM" tetapi ada di NB** | — catatan penamaan |

---

## 9. Pertanyaan terbuka

| # | Pertanyaan | Pemilik | Memblokir spec? |
| --- | --- | --- | --- |
| **1** | **Rumus PA dan TRAVEL belum dikutip** — Travel `.Rate` = 0, kemungkinan bentuk berbeda (§7.1) | — pekerjaan verifikasi lanjutan | ⚠️ **ya, untuk PA & Travel** |
| **2** | **Siapa mengisi `ProratePercentEDMEnd` / `ProratePercentStartEDM`** — pengisinya belum ditelusuri; apakah nilainya sama dengan `ProrateEDMEnd`/`ProrateStartEDM` (§4.2) | — verifikasi lanjutan | ⚠️ **ya** |
| **3** | **Pembalikan rasio pada `.PremiRp`** — disengaja? (§4.3) | work owner + Aktuaria | tidak — diport apa adanya |
| **4** | **Isi perbedaan 5 salinan EDM** dari mesin premi bersama (§6) | — verifikasi lanjutan | tidak — sebelum implementasi, bukan sebelum spec |

⛔ **Butir 1 dan 2 sebaiknya diselesaikan sebelum spec Seam 3 final** — keduanya menyangkut bentuk
kontrak perluasan, bukan detail implementasi.

---

## 10. Ringkasan untuk keputusan

✅ **Yang sudah pasti:**

1. EDM **memakai ulang** mesin premi NB — 41 activity bersama.
2. Rumus premi **dasar SAMA**, termasuk pembagi komposit `1e9` (K-018) dan pembulatan 4 desimal.
3. EDM **menambah** varian dua-bagian: porsi baru + porsi lama, masing-masing memakai rumus dasar.
4. **Tidak perlu seam baru.** Seam 3 **diperluas** menerima konteks endorsement.
5. Premi **dasar** dan delta **spreading** adalah dua perhitungan berbeda dengan **dua keluarga
   properti prorata yang berbeda** — jangan dicampur.

⚠️ **Yang belum:** rumus PA dan Travel, serta pengisi `ProratePercentEDMEnd`/`ProratePercentStartEDM`.

⛔ **`/to-spec` dan `/to-tickets` TIDAK dijalankan.**

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
