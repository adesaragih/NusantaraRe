---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 3 (`.scratch/claim-life/grilling-ronde-3.md`), keputusan work owner
---

# Unit status Claim — Life adalah **baris `AdjustmentList`**, bukan klaim

Status klaim Life diputuskan **per baris `AdjustmentList`**. `STS_REJECT` pada
`PremiumListDetail` dan pada header klaim adalah **cerminan** baris adjustment terakhir, bukan unit
keputusan tersendiri. Revisi keputusan dilakukan dengan **menambah baris baru**, bukan mengubah
baris lama.

Ini menggantikan pembacaan awal saya bahwa Claim — Life memiliki satu status di tingkat klaim.

## Mesin status `[keputusan work owner 2026-09-14]`

| # | Langkah | Peran | Akibat pada `STS_REJECT` baris |
| --- | --- | --- | --- |
| 1 | Input baris `AdjustmentList` pertama, lalu **Save ke OS** | `ReasLifeAdmin` | `0` |
| 1b | **Reject Outstanding** atas adjustment yang ia input sendiri — tanpa Komite. **Membatalkan baris itu saja**; Admin lalu input baris baru. | `ReasLifeAdmin` | `2` |
| 2 | Submit ke Medical Check | `ReasLifeAdmin` → `ReasLifeMedicalAdvisor` | tetap `0` |
| 3 | Submit ke SPV | `ReasLifeMedicalAdvisor` → `ReasLifeSPV` | tetap `0` |
| 4 | **Send ke Komite**, untuk semua adjustment | `ReasLifeSPV` — **kecuali `Type` `TP`/`TR`**, lihat ralat di bawah | tetap `0` |
| 5 | Komite memutus | Komite Life (konteks luar) | aksep → `1`; tolak → `2` |
| 6 | Setelah tolak: **tambah baris `AdjustmentList` baru** → Save ke OS → send Komite. Berulang. | `ReasLifeSPV` | baris baru mulai dari `0` |

**Aksep final selalu lewat Komite.**

### Ralat langkah 4 — `Type` `TP`/`TR` bebas peran

`[keputusan work owner 2026-09-14, Ronde 3 lanjutan]` Aturan "send ke Komite hanya SPV" berlaku
**hanya untuk tipe selain `TP` dan `TR`**. Untuk klaim ber-`Type` **`TP`** atau **`TR`**, pengiriman
ke Komite **tidak dibatasi peran** — **`ReasLifeAdmin` pun dapat mengirim langsung**.

Ini **meralat** kalimat "Admin tidak mengirim ke Komite" pada versi pertama ADR ini.

`[terverifikasi]` Gerbang visibilitasnya di `Claim Life/Section/AdjustmentDetail_Section.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION`):

```
pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'
```

Seluruh operatornya `||`, tanpa `&&` — jadi **tidak ada ambiguitas presedensi**: kedua nilai `Type`
itu memang meloloskan gerbang tanpa memeriksa peran.

`[keputusan work owner 2026-09-14]` **`TP` = Payable, `TR` = Receivable.** Definisi ini **tidak ada
di korpus**; sumbernya work owner.

**Perilaku ini dibawa apa adanya (paritas), tidak diperketat** — keputusan itu beserta risikonya
dicatat terpisah di **ADR-0012**, karena ia menyangkut wewenang, bukan mesin status.

### Kefinalan

- **Baris terminal.** Sekali sebuah baris bernilai `1` atau `2`, nilainya tidak berubah lagi.
- **Klaim tidak terminal.** Setelah penolakan, baris baru dibuat; siklus berulang — oleh SPV bila
  penolaknya Komite, oleh Admin bila ia sendiri yang menolak.
- **Revisi = baris baru**, bukan pengubahan baris lama.

### Arti tunggal nilai `2` `[keputusan work owner 2026-09-14]`

**`STS_REJECT = 2` SELALU berarti "baris ini ditolak" — TIDAK PERNAH "klaim selesai".**

Kedua sumber penolakan menulis nilai yang **sama** dengan **makna operasional setara pada tingkat
baris**:

| Sumber | Rule penulis | Akibat pada klaim |
| --- | --- | --- |
| **Admin** — "Reject Outstanding", tanpa Komite | `Claim Life/Activity/RejectOSClaimLife_Act.xml` | klaim **tidak** tertutup; Admin input baris baru |
| **Komite** — lewat SPV | `Komite Claim Life/Activity/KomitePostAdjustment.xml` | klaim **tidak** tertutup; SPV input baris baru |

Konsekuensi untuk spesifikasi: **tidak ada nilai `STS_REJECT` yang menandakan klaim selesai.**
Selesainya sebuah klaim bukan status tersimpan — ia keadaan turunan dari kumpulan barisnya.
`[pertanyaan terbuka]` **Aturan turunannya belum ditetapkan** dan bukan fakta yang hilang dari
korpus, melainkan keputusan desain — dicatat di `.scratch/claim-life/kesiapan-to-spec.md`, bukan
sebagai OQ.

`[terverifikasi]` Korpus memang tidak membedakan keduanya: kedua rule menulis `2` ke tingkat baris
**dan** ke `PremiumListDetail(idx)` dalam bentuk yang sama persis. Perbedaannya memang tidak ada —
sesuai keputusan di atas, karena memang tidak seharusnya ada.

## Bukti korpus yang menguatkan `[terverifikasi]`

### 1. Gerbang peran memang ada, dan bekerja pada `.STS_REJECT` tingkat baris

`Claim Life/Section/AdjustmentDetail_Section.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION`, 557.563 byte):

| Kontrol | `<pyCondition>` |
| --- | --- |
| **"Reject Outstanding"** | `pyWorkPage.pyPosition =='ReasLifeAdmin' && pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO !='' && .STS_REJECT == 0` |
| Jalur Komite (`GetListKomiteLife`) | `pyWorkPage.pyPosition =='ReasLifeSPV' \|\| pyWorkPage.Type = 'TP' \|\| pyWorkPage.Type = 'TR'` |
| **"Save to Outstanding"** (`SaveOutstandingLife_Act`) | `pyWorkPage.pyPosition =='ReasLifeSPV'` |

Gerbang Reject berbunyi **`.STS_REJECT == 0`** — relatif, yaitu **tingkat baris**. Inilah bukti
terkuat bahwa unit keputusan adalah baris: wewenang diuji terhadap status baris, bukan status klaim.

Perintah audit:
```
grep -n "pyPosition" "Claim Life/Section/AdjustmentDetail_Section.xml"
```

### 2. `RejectOSClaimLife_Act` **bukan** dead rule

Dirujuk 2× sebagai `<pyActivity>` dari `Claim Life/Section/RejectOSClaimLife_Sec.xml`, dan
kontrol "Reject Outstanding" di atas digerbangi peran `ReasLifeAdmin`. Ia menulis:

```
.STS_REJECT                                                    = 2
pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail(idx).STS_REJECT = 2
```

**Koreksi:** pada laporan sebelumnya saya menyebut adanya "dua jalur aksep/tolak" tanpa dapat
menjelaskan yang mana milik siapa. Sekarang jelas — jalur reject langsung adalah **milik Admin**.

### 3. Pencerminan ke `PremiumListDetail` ditulis berbarengan, bukan disalin belakangan

Setiap rule yang memutus menulis **dua tingkat sekaligus** dengan nilai sama:

| Rule | Tingkat baris adjustment | Tingkat `PremiumListDetail` |
| --- | --- | --- |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `.STS_REJECT = 2` | `= 2` |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `1` (6×) / `2` (2×) | nilai sama, pasangan langsung |

Jadi "cerminan" bukan tafsir — ia terbaca sebagai penulisan berpasangan.

`[keputusan work owner]` Bahwa yang tercermin adalah baris **terakhir** berasal dari work owner;
korpus hanya menunjukkan yang tercermin adalah **baris yang sedang diputus**.

### 4. Baris baru diwarisi dari baris pertama

`Claim Life/Activity/SetIndexAdjustmentList.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SETINDEXADJUSTMENTLIST`, 47.748 byte) menyalin **delapan**
kolom dari `AdjustmentList(1)` ke `AdjustmentList(<LAST>)`:

```
SHARE_NUSANTARA_RE  CEDING_RETENTION  SUM_REASURED  SUM_INSURED
SHARE_RETRO         RETROCEDED_SHARE  CURRENCYID    CURRENCY
```

`STS_REJECT` **tidak** termasuk — baris baru tidak mewarisi status. Ini persis yang dibutuhkan
langkah 6 (baris baru mulai dari `0`).

### 5. Tidak ada jalur buka-ulang pada baris yang sama

Sensus penulis lengkap (register **OQ-061**): **tidak ada satu pun rule** di `Claim Life` maupun
`Komite Claim Life` yang menulis `0` setelah `1` atau `2`. Satu-satunya penulis nilai `0` adalah
`SaveOutStandingLife_Act`, yang dipakai saat baris **baru** disimpan.

## Considered Options

- **Unit status = baris `AdjustmentList`; header mencerminkan baris terakhir** — dipilih
- Satu status di tingkat klaim, baris hanya rincian — ditolak: gerbang wewenang di korpus menguji
  `.STS_REJECT` **tingkat baris**, dan penolakan Komite menghasilkan baris baru, bukan perubahan
  status klaim
- Status ganda (klaim dan baris masing-masing otoritatif) — ditolak: menghasilkan dua sumber
  kebenaran untuk pertanyaan yang sama

## Consequences

- **Model data**: entitas keputusan adalah baris adjustment. Status klaim adalah **turunan**
  (*derived*), bukan kolom yang ditulis sendiri — meskipun sistem lama menyimpannya sebagai kolom.
- **Riwayat**: menolak lalu mengajukan ulang menghasilkan **beberapa baris** pada satu klaim.
  Jumlah baris = jumlah putaran Komite. Ini yang membuat klaim tidak terminal.
- **Jejak audit** (**ADR-0007**) direkam **per baris**, bukan per klaim — transisi yang perlu
  siapa+kapan adalah transisi baris.
- **Kontrak Komite** (**ADR-0001** Kontrak 2) memang bekerja pada tingkat baris; `IndexAdjustment`
  dan `IndexPremiumList` yang menyeberang kini masuk akal sebagai penunjuk baris.
- **Migrasi** (**ADR-0009**) harus membawa **seluruh baris**, bukan hanya keadaan terakhir —
  kalau tidak, riwayat putaran Komite hilang.
- **Uang** (**ADR-0003**) melekat pada baris; `CURRENCY` dan `CURRENCYID` pun ada di tingkat baris
  dan disalin dari baris pertama.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-060** | Apakah keseragaman mata uang per klaim adalah **aturan** atau kebetulan implementasi — `CURRENCY` ada di tingkat baris dan disalin dari baris 1 |
| **OQ-032** | Apa yang menentukan `KomiteLoop` — berapa putaran sebelum keputusan final |
| **OQ-001** | Tidak ada DDL — bentuk penyimpanan baris adjustment di Oracle tidak diketahui |

**Sudah tertutup untuk Claim — Life** (2026-09-14, work owner): **OQ-039** (akibat reject Admin),
**OQ-061** (unit keputusan), **OQ-062** (kefinalan), **OQ-063** (gerbang `TP`/`TR` — lihat
**ADR-0012**).

## `SaveAdjustment_Act` — dead rule, tidak dimigrasikan

`[keputusan work owner 2026-09-14]` `Claim Life/Activity/SaveAdjustment_Act.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SAVEADJUSTMENT_ACT` / `RULE-OBJ-ACTIVITY`, 179.221 byte)
**sudah tidak dipakai**. **Jangan dimigrasikan.** Penulis `STS_REJECT = 1` yang berlaku adalah rule
sisi Komite, `KomitePostAdjustment`.

⚠️ **Status "dead" ini TIDAK dapat diverifikasi dari korpus — ia keputusan work owner.** Yang
terbaca justru sebaliknya:

`[terverifikasi]` Rule itu **terpasang di UI** — dirujuk 2× sebagai `<pyActivity>` dari
`Claim Life/Section/ClaimLifeDetailGCNM.xml` (824.562 byte) — dan menulis
`.AdjustmentList(<LAST>).STS_REJECT = 1` bersama `ACCEPTEDNO` dan `ACCEPTATION_DATE`, **tanpa
precondition Komite**. Ia juga hanya ada di `Claim Life`, tidak di modul lain.

Perbedaan ini dicatat dengan sengaja: bila kelak ditemukan bahwa jalur itu masih dipakai di
produksi, keputusan "jangan dimigrasikan" perlu ditinjau ulang, dan bukti di atas adalah titik
mulanya.

Perintah audit:
```
grep -rn "<pyActivity>SaveAdjustment_Act</pyActivity>" "Claim Life" --include="*.xml"
grep -n "AdjustmentList(&lt;LAST&gt;).STS_REJECT" "Claim Life/Activity/SaveAdjustment_Act.xml"
```
