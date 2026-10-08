# Permintaan ekspor rule Property — nilai Prompt List

**Dibuat 7 Oktober 2026** atas laporan pemilik proses:

> *"coba cek semua property yang ada prompt listnya seharusnya yg diambil
> itu nilai dari prompt value bukan standar value, krn di aplikasi ada yang
> salah"*

---

## Mengapa ini perlu diminta

⛔ **`Rule-Obj-Property` NOL diekspor.** Sapuan keempat korpus
`D:\XML_NURE` (`Treaty In`, `Treaty In Adjustment`, `Claim Non Prop`,
`Komite Claim Non Prop`) menemukan nol berkas berjenis itu. Yang ada hanya
dua yang pernah dikirim terpisah ke
`_migration-docs/treaty-in-adjustment/ekspor-tambahan/`.

Rule Property menyimpan DUA kolom per pilihan:

| | |
|---|---|
| **Standard value** | yang TERSIMPAN di basis data — `risk`, `OR`, `1` |
| **Prompt value** | yang DIBACA pemakai — `Of Cession to R/I`, … |

Tanpa ekspornya, aplikasi hanya punya standard value, dan itulah yang
tampil di dropdown.

⛔ **Labelnya tidak ditebak.** `OptionLimit` membuktikan tebakan akan
meleset: pasangan yang masuk akal bagi `Of Cession to R/I` adalah sesuatu
tentang *cession*; yang sebenarnya `Of 100% Limit`. Pilihan berlabel
karangan terbaca BENAR sampai seseorang memilihnya, lalu angkanya salah
tanpa satu pun tanda di layar.

---

## Yang DIBUTUHKAN — enam rule, sebelas nilai

Cukup tangkapan layar tab **General → Display and validation → Prompt
values**, sama seperti `OptionLimit` kemarin.

### Kelas `ASM-FW-GISFW-Data-TreatyInLimits`

Keempatnya terbaca di `Section/Layers.xml` (dan `LayersEDM.xml`).

| Property | Standard value | Cacah terukur | Tampil di |
|---|---|---|---|
| `Cover` | `risk` · `cat` · `riskcat` | 1.089 · 1.124 · 2.564 | Limits Non-Prop → kotak **Cover** |
| `LayerType` | `layer` · `sublayer` | 4.457 · 323 | Limits Non-Prop → **Layers** dan **Part of** |
| `CurrencyRelation` | `OR` · `AND` | 3.342 · 138 | Limits Non-Prop → **Currency Relation** |
| `ReinstatementNote` | `asamount` · `astime` | 1.456 · 16 | Limits Non-Prop → catatan **Reinstatement** |

### Kelas `ASM-FW-GISFW-Int-TREATY_IN`

Kelas yang sama dengan `OptionLimit`. Hanya SATU nilai masing-masing yang
kurang — sisanya sudah punya label.

| Property | Standard value yang kurang | Cacah | Yang sudah ada |
|---|---|---|---|
| `AccountingModeNonProp` | `risk` | 21 | `loss` → `Loss Occuring` |
| `Bordeaux` | `nonreporting` | 687 | `reporting` → `Reporting` |

---

## Yang SUDAH lengkap — tidak perlu diminta

| Property | Label |
|---|---|
| `OptionLimit` | 1 → `Of Cession to R/I` · 2 → `Of 100% Limit` |
| `EDMState` | 1 → `Internal` · 2 → `External` |
| `EDMMaterialType` | 1 → `Material` · 2 → `Non Material` |
| `ReportingPeriod` | `quarter` → `Quarter Year` · `half` → `Half Year` · `month` → `Monthly` · `other` → `Others` |
| `AccountingMode` | `underwriting` → `Underwriting Year` · `accounting` → `Accounting Year` |

---

## ⭐ PERKEMBANGAN 7 Oktober 2026 — sebagian sudah terisi

Pemilik proses memerintahkan seluruh dropdown diperbaiki, dengan SATU
contoh: *"riskcat itu seharusnya yang diambil risk & cat"*.

Satu contoh itu MENGUNCI `Cover`. Sisanya **disimpulkan**, dan kesimpulan
disimpan TERPISAH dari yang terbukti (`promptDisimpulkan` lawan
`promptValue`) — `AsalLabel()` dapat ditanyai mana yang mana.

### Yang masih DISIMPULKAN — lima, dan inilah yang perlu diralat

| Nilai | Label sekarang | Dasar kesimpulan |
|---|---|---|
| `LayerType.sublayer` | `sub layer` | pola pemecahan yang sama dengan `riskcat` → `risk & cat` |
| `ReinstatementNote.asamount` | `as amount` | pola yang sama |
| `ReinstatementNote.astime` | `as time` | pola yang sama |
| `AccountingModeNonProp.risk` | `Risk Attaching` | ⚠️ **PALING BERANI** — bukan pemecahan kode, melainkan istilah baku reasuransi sebagai lawan `Loss Occuring` |
| `Bordeaux.nonreporting` | `Non Reporting` | pasangannya `reporting` → `Reporting` sudah ada |

⭐ `CurrencyRelation` (`OR`/`AND`) dan `LayerType.layer` TIDAK perlu
diralat: labelnya memang sama dengan nilainya.

⛔ `TestBerapaLabelMasihDisimpulkan` memaku angka **5**. Tiap ralat yang
datang memindahkan satu baris ke `promptValue` dan membuat uji itu merah
sampai angkanya diperbarui.

---

## Ke mana jawabannya dimasukkan

`modul/treatyin/backend/services/prompt_value.go`, peta `promptValue` —
satu tempat, sebagai DATA. Nol perubahan kode di titik panggilnya.

⭐ `TestBerapaLabelMasihDisimpulkan` **memaku angka 5**. Menambahkan label
yang datang akan membuatnya merah dan menagih angkanya diperbarui; menambah
dropdown baru tanpa label juga merah. Jadi kesenjangan ini menyusut
terukur, bukan terlupakan.

⚠️ Tiga property (`ReportingPeriod`, `AccountingMode`, `Bordeaux`) labelnya
hidup di penerjemah yang sudah ada dan **tidak** disalin ke peta: dua sumber
untuk satu label adalah cara termudah keduanya berbeda diam-diam.
