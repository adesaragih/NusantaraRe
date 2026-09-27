# Audit ulang XML atas seluruh yang sudah dibangun — Claim Life tiket 01–15

`[terverifikasi]` Dibaca 27 September 2026 dari korpus `D:\XML\RNM_BRD\Claim Life\` dan
`D:\XML\RNM_BRD\Komite Claim Life\`, keduanya READ-ONLY. Titik tetap **A1 = `c7ce88e`**.

Cara baca kolom **verdict**: **tepat** = kode cocok dengan XML, dengan bukti baris ·
**kurang tepat** = ada selisih yang dapat ditunjuk · **belum ditiru** = rule dirujuk komentar tetapi
perilakunya belum ada di kode · **tidak ditiru, beralasan** = sengaja, dengan bukti.

⛔ Verdict "tepat" pun membawa nomor baris. Tanpa bukti ia bukan audit, hanya pernyataan.

---

## Inventaris — rule yang dirujuk kode

27 rule disebut komentar kode *(sensus: pemindaian `internal/`, `cmd/`, `pkg/`, `frontend/src`
atas nama berakhiran `_Act`/`_act`/`_SQL`/`_sql`/`_Section`/`_Flow`/`_Harness` dan atas pola
`` `<Modul>/<Jenis>/<Nama>.xml` ``)*.

---

## Hasil audit — putaran pertama

| # | Rule *(path pecahan)* | Kode | Tiket/AC | Verdict | Bukti |
| ---: | --- | --- | --- | --- | --- |
| 1 | `Activity/SaveAdjustment_Act.xml` 1833 · 1879 · 1899 | — | 04, 08 | ⛔ **belum ditiru** | menulis `.AdjustmentList(<LAST>).ACCEPTEDNO`, `.STS_REJECT = 1`, `.ACCEPTATION_DATE = @CurrentDateTime()`. **Akseptasi di Claim Life sendiri** |
| 1a | idem, gerbang 854 dan 1048 | — | — | — | `.ACCEPTEDNO=="" && .IsCheck=true && .STS_REJECT=="0"` **dan** `Type=="QP"\|\|"QR"` (854) atau `Type=="TP"\|\|"TR"` (1048); `WhenTrue=2`/`WhenFalse=3` |
| 1b | pemicu: `Section/ClaimLifeDetailGCNM.xml` 22641 → 22665 | — | — | — | tombol berlabel **"Save Adjustment"**. Kemunculan kedua di 22773 bertetangga `pyCondition 1=2` (selalu palsu) — **mana yang hidup belum dapat dipastikan dari kedekatan baris saja**, `[terbuka]` |
| 2 | `Activity/ValidasiDOL_Act.xml` 458 · 698 | `services/dol.go` | 06 | ⛔ **kurang tepat → DIPERBAIKI** | 698 menggeser DOL **satu HARI**, bukan satu jam. Lihat §"Ralat pergeseran DOL" |
| 3 | `Activity/LoadDataPeserta_Act.xml` 1097 · 1143 · 1163 · 1183 · 1203 · 1223 · 1243 · 1263 | — | — | ⚠️ **tidak ditiru, dinyatakan** | kedelapan tanggal peserta dimuat lewat `@addCalendar(...,0,0,0,0,7,0,0)` = **+7 jam (WIB)**. Sistem kita menyimpan dan menampilkan konsisten, jadi tidak menggeser; `TestTanggalDiuraiTanpaGeserZona` mengatur jalur MIGRASI, hal yang berbeda. `[terbuka — work owner]` bila tampilan harus sama persis dengan Pega |
| 4 | `Activity/SendtoAdmin_Act.xml` 259 · `SendtoAdmin_Act1.xml` 282 · `SendtoMedical_Act.xml` 260 | `services/tahap.go` | 08 | ✅ **tepat** | ketiganya satu `Property-Set`: `"1"`, **`"0"`**, `"1"`. `JalurBalik{}.NilaiSendto()` menulis `""` untuk yang bukan jalur balik — setara pembersihan `_Act1` |
| 5 | `Activity/UpdateDateClaimLife_Act.xml` 252 · 299 · 320 · 341 | `repository.PerbaruiTanggalKejadian` | 06 | ⛔ **kurang tepat** | penulisnya menulis **empat** tanggal — `DATE_OF_LOSS`, `CLAIM_RECEIVED_DATE`, `COMPLETE_DATE`, `CONFIRMATION_DATE` — beserta `NAME_OF_INSURED` (362) dan `CERTIFICATE_NO` (383). Kode kita menulis **satu**: `UPDATE … SET DATE_OF_LOSS = …`. Diperbaiki di kelompok **Register** (A3) |
| 6 | `Activity/ValidasiDOL_Act.xml` *(seluruhnya)* | `services/dol.go` | 06 | ✅ **tepat** | hanya validasi: `Page-Clear-Message`, dua `Property-Set`, `Property-Set-Message`; nol penulisan kolom |
| 7 | `Activity/SaveOutStandingLife_Act.xml` 3446 · 3679 · 3823 · 4054 · 4199 | `models/businesscode.go` | 06 | ⛔ **kurang tepat** | tiga langkah bergerbang `ContentNote = "DEATH"`, **dua** bergerbang `!= "DEATH"`. Sensus kode kita: dari 21 kode, **11 DEATH · 5 HEALTH · 2 CI · 2 TPD · 1 TI** — jadi cabang non-DEATH terjangkau **10 dari 21**, dan kita tidak bercabang sama sekali. Selisihnya: cabang non-DEATH (4262) membandingkan DOL dengan jendela **GROSS** memakai `@addCalendar(...,0,0,0,0,0,0,0)` |
| 8 | `Activity/KomitePostAdjustment.xml` *(Komite)* | `services/hasilkomite.go` | 11 | ✅ **tepat** | disensus ulang tiket 11: 4 tulisan ke BARIS (3× `1`, 1× `2`), 4 ke PESERTA; empat gerbang `KomiteCount == KomiteLoop` |
| 9 | `Activity/GetLinkService.xml` 371 · 393 · 491 · 517 · 701 · 705 | `services/efekkeluar.go` | 12 | ✅ **tepat** | disensus tiket 12; penyimpangan hasil-kosong dinyatakan |
| 10 | `When/IsPEGAPROD.xml` 164 · 318 | `services/efekkeluar.go` | 12 | ✅ **tepat** | `pxProcess.pzProductionLevel = "5"` |

---

## Ralat pergeseran DOL — `[dugaan]` tiket 06 ditutup oleh korpus

⛔ **Ronde pertama menebak tanda tangan `@addCalendar` dari CACAH ARGUMEN saja** dan menulisnya
terus terang sebagai `[dugaan — Product+UW]`: *"Tujuh angka COCOK dengan tanda tangan yang
berargumen milidetik … dan pada tanda tangan itu posisi keempat adalah jam."*

**Korpus sendiri menentukannya, dan itu terlewat.** `LoadDataPeserta_Act` memuat kedelapan tanggal
peserta dengan `@addCalendar(<tanggal>,0,0,0,0,7,0,0)` — angka **7** di posisi **kelima**:

| Tafsir | posisi 5 | Masuk akal? |
| --- | --- | --- |
| lama *(…, hari, jam, menit, detik, milidetik)* | menit | ⛔ tidak ada yang menggeser tanggal LAHIR tujuh menit |
| — | minggu | ⛔ tujuh minggu memindahkan orang ke bulan lain |
| **baku** *(tahun, bulan, minggu, hari, **jam**, menit, detik)* | **jam** | ✅ **+7 jam = WIB (UTC+7)** — normalisasi zona untuk data Indonesia |

Maka posisi keempat = **HARI**, dan `ValidasiDOL_Act` baris 698 `(.DATE_OF_LOSS,0,0,0,1,0,0,0)`
menggeser DOL **satu hari** sebelum dibandingkan dengan jendela `RETROCESSION_VALUATION_*`
(baris 719, 746). Cabang jendela `GROSS` di baris 458 berargumen nol seluruhnya — tidak bergeser
(baris 486, 526).

**Diperbaiki:** `pergeseranTPTR = time.Hour` → `pergeseranDOLRetro = 24 * time.Hour`; namanya ikut
berubah sebab ia jendela **retro**, bukan Type TP/TR. Labelnya naik dari `[dugaan — Product+UW]`
menjadi `[terverifikasi — turunan]`, dengan jangkarnya tertulis di kodenya.

**Kasus yang berbalik** — persis "daftar kasus yang berbalik" yang komentar tiket 06 janjikan:
sepuluh baris kasus `TestBatasJendelaTiapType` untuk TP/TR diganti dari batas per-jam menjadi
batas per-hari. Dibuktikan: mengembalikan konstanta ke `time.Hour` memerahkan
`TestPergeseranDOLSatuHariBukanSatuJam` **dan** `TestBatasJendelaTiapType`.

---

## Yang masih harus diaudit *(putaran berikutnya)*

Titik brief §A0 yang belum dijawab: 1 *(kaskade `TO_NUMBER(IDR)`)*, 3 *(tulisan datar sesudah
reject)*, 6 *(field pendaftaran)*, 7 *(pemilihan dan hapus peserta)*, 8 *(peran per tahap)*,
10 *(7 dari 12 kolom dokumen)*, 11 *(muatan Komite dan roster)*, 12 *(efek keluar)*,
13 *(hapus klaim)*, 14 *(`IsCheck` di seluruh jalur)*.

⚠️ Titik 15 *(`@CompareDates`, `@addCalendar`)* **sebagian tertutup**: `@addCalendar` tidak lagi
`[dugaan]` — tanda tangannya ditentukan jangkar WIB di atas. `@CompareDates` tetap `[dugaan]`;
korpus tidak memuat definisinya, dan nilai bandingnya tidak dapat diturunkan dari contoh yang ada
(`@CompareDates("20240103","20240102")` di `SavePesertaClaim` 674 dan
`@CompareDates("02/02/2024","14/11/2022")` di `ValidasiDOL_Act` 737 keduanya "a sesudah b",
sehingga tidak membedakan ketat dari tidak ketat).
