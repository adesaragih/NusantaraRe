# Pemetaan medan kepala kontrak Treaty In — sapuan lengkap, kedua cabang

**Diukur 5 Oktober 2026** atas `Section/TreatyInNONProportional.xml`, **satu-satunya** Section
yang memuat kepala kontrak — ia melayani **kedua** cabang, dan strip tab-nyalah yang berbeda
(`TreatyInTabsProportional.xml` lawan `TreatyInTabsNonProportional.xml`).

⛔ `<pyIncludedRuleXML>` dibuang lebih dulu; seluruh offset di bawah terhadap teks **bersih**
(6.904.904 bita).

---

## 1 · Keenam belas medan kepala, berurut layar

| # | `pyValue` | Label ekspor | Kontrol | Syarat tampil | Sumber data kita | Di layar kita? |
| ---: | --- | --- | --- | --- | --- | --- |
| 1 | `TreatyIn.ID` @31.441 | *ID-hidden testing* | `pxTextInput` | ⛔ `Never` | — | tidak — **kode mati** |
| 2 | `TreatyIn.TreatyContractName` @38.115 | Treaty Contract Name | `pxTextInput` | `ALWAYS` | kolom `TREATY_IN.TREATYCONTRACTNAME` | ⭐ ya |
| 3 | `TreatyIn.ContractRefNo` @44.674 | Contract Ref No | `pxTextInput` | `ALWAYS` | `JSONDATA.ContractRefNo` (742/1.854) | ⭐ ya — **selalu** (ralat 7 Okt 2026, §3a) |
| 4 | `TreatyIn.TeritorialScope` @50.841 | Teritorial Scope | `pxTextArea` | `ALWAYS` | kolom `TERITORIALSCOPE` | ⭐ ya |
| 5 | `TreatyIn.Bordeaux` @56.974 | Bordereaux | `pxDropdown` | **`ProportionType='Proportional'`** @60.649 | `JSONDATA.Bordeaux` (1.851) | ⭐ ya — **kini prop saja** |
| 6 | `TreatyIn.BordereauxNote` @63.375 | Bordereaux Note | `pxTextArea` | `ALWAYS` | `JSONDATA.BordereauxNote` (1.018) | ⭐ ya — **selalu** (ralat 7 Okt 2026, §3a) |
| 7 | `TreatyIn.Commencement` @84.930 | Commencement | `pxDateTime` | — | kolom `COMMENCEMENT` | ⭐ ya |
| 8 | `TreatyIn.Termination` @95.371 | Termination | `pxDateTime` | `ALWAYS` | kolom `TERMINATION` | ⭐ ya |
| 9 | `TreatyIn.TreatyYear` @102.091 | Treaty Year | `pxTextInput` | `ALWAYS` | kolom `TREATYYEAR` | ⭐ ya, baca-saja |
| 10 | `TreatyIn.AccountingMode` @108.010 | Accounting Mode | `pxDropdown` | **`ProportionType='Proportional'`** @111.777 | `JSONDATA.AccountingMode` (1.851) | ⭐ ya — **kini prop saja** |
| 11 | `TreatyIn.AccountingModeNonProp` @114.491 | Accounting Mode | `pxDropdown` | **`ProportionType='NonProportional'`** @118.211 | `JSONDATA.AccountingModeNonProp` (1.851) | ⭐ **BARU ronde ini** |
| 12 | `TreatyIn.Ceding` @120.957 | Ceding | `pxAutoComplete` | — | kolom `CEDING` + `CEDINGID` | ⭐ ya, + tombol `Choose Ceding` |
| 13 | `TreatyIn.LeadingReinsSource` @132.888 | Business Source / Source of Business | `pxAutoComplete` | — | kolom `LEADINGREINSSOURCE` + `…ID` | ⭐ ya |
| 14 | `TreatyIn.LeadingReinsName` @144.641 | Leading Reinsurer | `pxAutoComplete` | ⛔ `1=2` @146.051 | — | tidak — **kode mati** |
| 15 | `TreatyIn.TreatyLeader` @191.054 | RNM as Treaty Leader @191.624 | `pxCheckbox` | `ALWAYS` (ralat §3a — `TreatyMasterInEDM` milik sel tetangga) | `JSONDATA.TreatyLeader` (659) | ⭐ ya — **selalu** |
| 16 | `TreatyIn.EDMEffective` @240.791 | Effective Date @239.791 | `pxDateTime` | `TreatyIn.ViewState != 1` @238.233 | ⛔ nol | **tidak** — lihat §2 |
| 17 | `TreatyIn.IsProRate` @250.948 | *(kotak centang "Pro Rate:")* | `pxCheckbox` | ⚠️ `TreatyMasterInEDM` | ⛔ nol | **tidak** — lihat §2 |
| 18 | `TreatyIn.ProRateDays` @289.738 | *(angka di samping "Pro Rate:" @280.747, satuan `%` @284.933)* | `pxTextInput` | `ALWAYS` | ⛔ nol | **tidak** — lihat §2 |

⚠️ **Dua label untuk medan 13.** Ekspor memakai `Business Source` di selnya (@131.896) dan
`Source of Business` di teks tampil baca-saja (@223.691). Layar kita memakai yang kedua, dan itu
yang tangkapan layar perlihatkan.

---

## 2 · ⛔ Tiga medan EDM / Pro Rate yang BELUM ada sumbernya

`EDMEffective`, `IsProRate`, dan `ProRateDays` **tidak ada** di `M_TREATY_IN.JSONDATA` maupun di
kolom `TREATY_IN` mana pun. Mereka milik mesin **pro rata**, dan `GRL-15` menyatakan mesin itu
**sengaja tidak dibangun kembali**.

⚠️ Dua blok di sekitar `ProRateDays` berpenjaga `NEVER` (@286.629, @296.888) — tetapi medannya
sendiri `ALWAYS`. Jadi yang mati hiasan di sekitarnya, bukan medannya.

⛔ **Tidak dibangun, dan bukan karena terlewat:** medan yang dirender tanpa sumber akan selalu
kosong, dan kotak yang selalu kosong di sebelah kotak berisi membuat pemakai mengira datanya
hilang. Begitu jalur pro rata punya rumah, ketiganya masuk di sini.

---

## 3a · ⭐ RALAT 7 Oktober 2026 — ketiga medan kepala SELALU tampil

Dibaca ulang dengan alat baca bersama (`pyIncludedRuleXML` dibuang):

| Sel | Kontrol | Tampil | Baca-saja / mati |
| --- | --- | --- | --- |
| `TreatyIn.ContractRefNo` | `pxTextInput`, `pyWidth` 130 | `ALWAYS` | `pyReadOnlyCondition TreatyIn.IsEditData= 1` · `pyDisabledWhen TreatyIn.EDMMaterialType = 1` |
| `TreatyIn.BordereauxNote` | `pxTextArea`, `pyWidth` 0 | `ALWAYS` | `pyReadOnlyCondition TreatyIn.IsEditData= 1 \|\| TreatyIn.EDMMaterialType = 1` |
| `TreatyIn.TreatyLeader` | `pxCheckbox`, `pyCheckboxCaption` "RNM as Treaty Leader", `pyIncludeLabel=false` | **tanpa syarat** | `pyDisabledWhen` / `pyReadOnlyCondition TreatyIn.ViewState = 1` |

- ⛔ **§3 di bawah keliru membaca sel.** `VIS[TreatyMasterInEDM]` milik `TreatyIn.EDMEffective` dan
  `TreatyIn.IsProRate` — sel tetangganya — bukan `TreatyLeader`. Medan ini tidak bersyarat.
- ⛔ **`MedanTakAda` dicabut.** Layar dahulu mengganti ketiga medan dengan medan mati
  "Tidak ada di dokumen sistem lama" bila kuncinya tak ada di dokumen lama. Ekspor tidak mengenal
  syarat itu, dan sejak tabel `T_TREATY_*` sengaja dikosongkan ketiganya mati pada SETIAP kontrak.
  Kini `FormKontrakTreatyIn` selalu merender `Field` · `textarea` · kotak centang ber-keterangan.
- `IsEditData` didekati mode lihat (`<fieldset disabled>`), preseden `TabPortofolio.tsx` /
  `TabRetensi.tsx`; `EDMMaterialType = 1` dibaca dari `edmJenisMaterial` kontrak.
- Dijaga `layar.test.ts` (`selalu dirender — nol MedanTakAda`).

## 3 · ⚠️ `RNM as Treaty Leader` bersyarat, dan syaratnya tidak terpenuhi hari ini

> ⛔ **Dibatalkan oleh §3a** — syarat di bawah milik sel tetangga. Dibiarkan sebagai jejak.

`TreatyMasterInEDM` adalah rule `When` (`When/TreatyMasterInEDM.xml`, `pyRuleAvailable = Yes`),
dan isinya tiga perbandingan:

```
TreatyIn.EDMState = "1"   atau   = "2"   atau   = "3"
```

**Diadu dengan POOLDATA:** `EDMState` ada di **2** dari 1.854 dokumen, dan keduanya bernilai
`0` — jadi `TreatyMasterInEDM` **salah pada setiap kontrak**.

⛔ **Medannya TIDAK disembunyikan ronde ini, dan itu keputusan sadar.** Tiga sebabnya:

1. `JSONDATA.TreatyLeader` **ada di 659 kontrak** dengan nilai sungguhan (`true` 103,
   `false` 556). Menyembunyikan medan yang punya data nyata di sepertiga kontrak adalah
   perubahan besar untuk dijalankan tanpa persetujuan.
2. `EDMState` adalah properti yang sama yang §4 `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` tanyakan:
   ia mungkin diisi Activity saat jalan dan tidak pernah tersimpan. Menyembunyikan medan atas
   dasar properti yang belum jelas asalnya berarti menebak.
3. Ia bukan tab, dan §0 ronde ini tentang tab.

**Yang ditanyakan:** apakah `RNM as Treaty Leader` memang hanya tampil pada kontrak ber-EDM — dan
kalau ya, bagaimana 659 kontrak itu memperoleh nilainya?

---

## 4 · Grid `Rate of Exchange` dan tombol-tombolnya

Panel `Rate of Exchange` @390.574. Kolomnya sudah terpasang (`Currency` · `Value to IDR` ·
`Valid From` · `Valid Until`).

| Tombol | Di layar kita |
| --- | --- |
| `Add` | ada, **mati** — baris baru milik jalur tulis (tiket 20) |
| `Delete` | ada, **mati** — idem |
| `Apply` (Reporting Period) | ada, **mati**; galatnya `Start Date, Due. Must Not Be Empty` disalin utuh |
| `Choose Ceding` | ada, **hidup** — ia hanya memilih |
| `Choose Source of Business` | ada, **hidup** — idem |
| `Save` · `Close` · `Actions` | ada, **mati** |
| `Update Total` (Maximum Retention) | ⭐ **BARU ronde ini**, mati |
| `Download All` (Attachment) | ada, **mati** — `T_STORAGE_IMAGE` belum berjalur |
| `Refresh` (Attachment) | ada, **hidup** — membaca ulang itu yang panel baca-saja boleh lakukan |

⛔ **Aturannya satu:** tombol tanpa jalur tulis **dinonaktifkan, bukan dihilangkan**. Tombol
hidup yang tidak melakukan apa pun berbohong; tombol hilang menyembunyikan bahwa layar lama
punya langkah itu.

---

## 5 · Sumber data per tab — tetap seperti §7 briefing menetapkan

| Sumber | Yang dilayani |
| --- | --- |
| kolom `TREATY_IN` | 8 medan kepala (#2, #4, #7, #8, #9, #12, #13, + `PROPORTIONTYPE`) |
| `M_TREATY_IN.JSONDATA` | 5 medan kepala (#3, #5, #6, #10, #11), 2 tab teks, 3 kunci syarat tab, grid `Rate of Exchange` |
| 9 tabel pendaratan `M_TREATYIN_*` | Reporting Period · Portfolio · Accumulation · EGNPI · Maximum Retention · Installment · Information & Submit · Co-Ins Scale · History |
| `M_TREATY_IN2` | Limits · Share · Event Limits · RNM Share (41 kolom, 31 terpakai) |
| `M_ATTACHMENTTREATY_2` | Attachment |
| **diturunkan** | `Total Retention Amount` (dari Maximum Retention) |

⛔ Kesembilan tabel pendaratan **tidak dibongkar** — 27.238 baris, dan tab yang sudah membacanya
tetap membacanya.
