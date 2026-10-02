# Laporan Penutup STEP D1 — Inventaris Rule Pega Nusantara Re

**Status: STEP D1 SELESAI.** 20 dari 20 modul, **9.369 dari 9.369 file XML (100%)**.
Tanggal: 2026-09-12. Korpus sumber: `D:\XML\RNM_BRD\` (READ-ONLY).

STEP D2 (telusur flow end-to-end) **belum dijalankan**. Dokumen ini menutup D1 dan menyiapkan D2.

---

## 1. Apa yang D1 hasilkan, dan apa yang TIDAK

**Dihasilkan:** inventaris lengkap setiap rule di korpus — identitas (`class / nama / tipe`),
silsilah salinan, objek database yang disentuh, stored procedure yang dipanggil, target integrasi
eksternal, tabrakan nama, dan konflik isi antar modul.

**TIDAK dihasilkan — dan tidak boleh dianggap ada:**

- **Perilaku rule.** D1 tidak membaca logika satu pun rule. Itu pekerjaan D2.
- **Arti kode/status apa pun.** Seluruh enumerasi yang ditemukan masih `belum terverifikasi`.
- **Model data.** Tidak ada DDL di korpus; struktur kolom JSON tidak diketahui.
- **Keputusan desain, ADR, BRD, spec, atau tiket.** Itu FASE B, setelah gate D4.

Artefak: `inventory/<modul>.md` (20), `inventory/_summary.md` (FINAL),
`inventory/_oq011-konflik-isi.md` (533 entri), `open-questions.md`, `README.md`.

---

## 2. Angka pokok korpus

| Ukuran | Nilai |
| --- | ---: |
| Modul | 20 |
| File XML | **9.369** |
| Ukuran | 1.319,3 MB |
| **Identitas rule unik** | **4.444** |
| File yang merupakan salinan identitas lintas modul | 4.925 (52,6%) |
| Tabel/view unik terlihat | 216 |
| Skema Oracle lokal | 11 |
| Stored procedure kustom **tanpa body** | 67 |
| Rule `ConnectREST` | 51 (15 nama service unik; **5** memuat URL literal) |
| Rule `Flow` | 24 |
| **Identitas berkonflik isi antar modul** | **533** |

**Lebih dari separuh file adalah duplikat identitas.** Volume migrasi sebenarnya diukur dari 4.444
rule — dengan syarat salinannya identik, yang **tidak selalu benar** (§4.1).

Perintah audit ada di `inventory/_summary.md` §22.

---

## 3. Temuan struktural per domain

| Domain | Modul | File | Identitas unik | Pola ruleset |
| --- | ---: | ---: | ---: | --- |
| Treaty inward | 4 | 1.149 | 729 | **sebagian besar satu** |
| Facultative inward | 3 | 6.071 | 2.552 | **satu ruleset** |
| Klaim (fac + non-fac + komite) | 8 | 1.467 | 1.018 | **terpisah** |
| Life, master & outward | 5 | 682 | 623 | **terpisah** |

`[terverifikasi]` Bukti pola:

- **Facultative inward = satu ruleset.** `RNW Fac In` berbagi **99,0%** identitasnya dengan
  `NB FacIn` (1.906 dari 1.926). **1.567 identitas (61,4%)** ada di ketiga modul
  (NB / RNW / Endorsment). 6.071 file hanya berisi 2.552 rule. → **OQ-015**
- **Treaty inward = sebagian besar satu.** `Treaty In` ↔ `Treaty In Adjustment` berbagi **320**
  identitas dari modul berukuran 329 dan 379. → **OQ-010**
- **Klaim = terpisah.** Overlap tertinggi trio Claim hanya **102** (~37%); `Claim Life` hanya
  12–14% terhadap keduanya. Modul Komite membawa **44,4% identitas miliknya sendiri** (59 dari 133)
  — bukan wrapper approval tipis. → **OQ-019**
- **Life/master/outward = terpisah.** Overlap tertinggi hanya **23** identitas. Tidak ada
  "ruleset Life bersama".

**Implikasi perencanaan:** penghematan dari konsolidasi rule **besar di facultative dan treaty,
kecil di klaim dan life**. Jangan menggeneralisasi rasio file→rule antar domain.

### 3.1 Dua framework Pega

`[terverifikasi]` Class aplikasi terbagi dua keluarga: `ASM-FW-GISFW-*` dan `ASM-FW-GCNMFW-*`.
`GCNMFW` **mendominasi domain klaim** (Claim Prop 190 vs 53; Claim Non Prop 169 vs 81) dan tidak
muncul sama sekali di batch 1. Prefix kunci RDBList sejalan: `ASM` 751, `RNM` 315, `GCNM` 80.
Kepanjangan kedua singkatan **belum terverifikasi** → **OQ-008**.

### 3.2 Nama tidak dapat dipercaya — terbukti berulang kali

- **268 nama file** dipakai di lebih dari satu tipe rule; **55 nama `When`** dipakai di lebih dari
  satu class (`ISMAINTENANCE` dan `ISCLAIM` masing-masing di 3 class).
- **Nama folder bukan tipe rule**: `Claim Fac In/Activity/InsertLogServiceClaim.xml` sebenarnya
  `RULE-CONNECT-SQL`; `Claim Fac In/Section/CedingCedant.xml` sebenarnya `RULE-HTML-HARNESS`.
- **Nama modul pun bisa menyesatkan**: `Treaty Contract Out` tidak merujuk satu pun objek treaty
  outward, sementara tujuh modul lain melakukannya → **OQ-022**.
- Sebuah rule `SELECT` bernama `GETCURRENCY` ternyata salinan dari `UPDATEMASTERCURRENCY`.

Karena itu aturan `README.md` §3.2 mengikat untuk D2: rule **selalu** ditulis `class / nama / tipe`,
diambil dari `<pxInsName>` + `<pzOriginalInstanceKey>`/`<pxObjClass>`.

---

## 4. Risiko halusinasi yang teridentifikasi

Empat risiko konkret. Masing-masing punya pemicu, dampak, dan cara menghindarinya di D2.

### 4.1 OQ-011 — 533 identitas dengan isi berbeda antar modul

**Risiko paling konkret di korpus.** `class + nama + tipe` yang sama merujuk isi **berlainan**
tergantung modul mana yang dibuka. Membaca satu varian lalu menganggapnya mewakili yang lain akan
menghasilkan pemahaman perilaku yang **salah, tetapi terlihat meyakinkan**.

- 2.545 identitas muncul di >1 modul; 2.012 identik, **533 berbeda**.
- 474 dari 533 menyentuh lebih dari 2 modul.
- Terparah: **`@BASECLASS / ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH`** — masing-masing di
  **6 modul dengan 6 isi berbeda**; tidak ada dua modul yang sepakat. Keempatnya rule `When`,
  yaitu **guard percabangan**.
- Satu rule `Flow` ikut berkonflik: `ASM-FW-GISFW-WORK / OFFERFACRETRO`.

**Cara menghindari di D2:** sebelum menelusur rule apa pun, **cek register
`inventory/_oq011-konflik-isi.md`**. Bila identitasnya ada di sana, telusur **setiap varian
terpisah** dan nyatakan perbedaannya — atau tunggu OQ-011 dijawab. **Jangan memilih satu versi.**

### 4.2 OQ-018 — korpus kemungkinan berasal dari lingkungan bukan-production

`[terverifikasi]` Sapuan 20 modul: URL literal ditemukan di **81 file di 15 modul**. Host yang
muncul mencakup IP internal (`192.168.105.112`, `192.168.105.116`), hostname internal
(`ssdecamwin02/03:7070`, `sdvpwin105:8383`), `10.100.10.75:7315`, `app.sinarmas.co.id`,
`pega.nusantarare.com`, SaaS pihak ketiga `view.officeapps.live.com` (27×, artinya **dokumen
dikirim ke layanan luar untuk ditampilkan**), dan — yang paling penting —
**`appdev.nusantarare.com` (7×)**.

Lima dari 51 rule `ConnectREST` memakai `<pyBaseURLSelectionType>URL` alih-alih konfigurasi
`LinkService`; empat di antaranya service pencatatan klaim Non-Proportional.

Di luar itu, metadata ekspor memuat dua `pxHostId` berbeda (`jboss1073`, `jboss122117`), yang
berarti korpus ini dirakit dari **lebih dari satu server Pega**.

**Risiko:** bila sebagian ekspor diambil dari lingkungan dev, maka sebagian rule mungkin **bukan
versi production** — yang juga menjadi penjelasan alternatif untuk sebagian dari 533 konflik OQ-011.

**Cara menghindari:** jangan memperlakukan korpus sebagai cerminan production sebelum OQ-018
dijawab. Setiap pernyataan perilaku di D2 menyebut **file mana** yang dibaca.

### 4.3 OQ-002 / OQ-017 — logika material berada di luar korpus

- **67 stored procedure kustom tanpa body**, termasuk `FIRE.PEGA_FIRE_SET_RATE` (penetapan rate)
  dan `FIRE.CEK_PRORATA_TANGGAL`.
- **Database link `@ASMD.SINARMAS.CO.ID`** — 8 objek berisi tabel rate per peril (EQ, terorisme,
  banjir, FLEXAS, RSMD, BI indemnity) dan kurs standar. Hanya dipakai 3 modul facultative inward,
  nol di 17 modul lain.
- **`COMMIT` dijalankan di dalam stored procedure**, dan `{OperatorID.pyUserIdentifier}` dikirim
  sebagai parameter → **OQ-013**.
- **Kolom JSON opaque** (`JSON_POLIS`, `JSON_KLAIM`, `JSONDATA`, `JSON_OFFER_LIFE`) → **OQ-012**.

**Risiko:** godaan untuk "menyimpulkan" apa yang dilakukan sebuah procedure dari namanya.
`PEGA_FIRE_SET_RATE` **tidak** memberi tahu rumusnya.

**Cara menghindari:** di D2, pemanggilan procedure dicatat sebagai **batas pengetahuan** —
"rule X memanggil `POOLDATA.Y` dengan parameter a, b, c; hasilnya dipakai untuk Z; **isi Y tidak
diketahui**". Bukan sebagai langkah yang dijelaskan.

### 4.4 OQ-020 / OQ-021 — kode dan identitas yang menggerbangi logika

- **Kode tanpa arti:** `PaymentType` (0–7), `TransferType` (1–4), `EdmType` (1, 3),
  `ProRateType` (1, 2, 3), `.Type` PremiumList (`QR`, `QP`, `TP`, `TR`). Pengelompokannya berbeda
  antar rule — mis. `{1,2,5}`, `{4,6}`, `{3}` dalam satu rule. Korpus tidak memuat tabel kode.
- **Identitas orang sebagai guard:** 66 file di 9 modul membandingkan `OperatorID.pyUserIdentifier`
  dengan nama orang literal; 35 file di 14 modul memakai `pyPosition`. Guard ini berada di
  `<pyStepsPreCondParamsWhen>` (193×) dan menggerbangi antara lain **perhitungan premi**
  (`COUNTGROSSPREMI_ACT`) dan **keputusan komite klaim** (`KOMITEPOST_REJECT`).

**Risiko:** menebak bahwa `PaymentType=1` berarti "tunai", atau menafsirkan guard identitas sebagai
aturan bisnis yang disengaja.

**Cara menghindari:** catat literalnya apa adanya; arti ditulis `belum terverifikasi`.

---

## 5. Pertanyaan terbuka per pemilik peran

**20 OQ terbuka, 2 terjawab** (OQ-004, OQ-006). Rincian di `open-questions.md`.

### 5.1 Yang MEMBLOKIR D2 (harus dijawab atau dikelola sebelum telusur flow)

| OQ | Judul | Pemilik | Mengapa memblokir D2 |
| --- | --- | --- | --- |
| **OQ-011** | 533 identitas berbeda isi | Product+UW + pemilik export | Menelusur varian yang salah = pemahaman perilaku yang salah |
| **OQ-018** | Host/URL hardcode, hostname DEV | IAM + DBA | Menentukan apakah korpus mencerminkan production |
| **OQ-020** | Arti kode yang menggerbangi logika | Finance + Product+UW | Percabangan tidak dapat dijelaskan tanpa arti kodenya |
| **OQ-005** | 5 modul tanpa rule `Flow` | Product+UW | Titik masuk telusur harus ditentukan lebih dulu |
| **OQ-022** | `Treaty Contract Out` — outward atau bukan | Product+UW | 303 rule berisiko ditelusur dengan premis yang salah |

### 5.2 Yang memblokir FASE B (boleh ditunda sampai setelah D2)

| OQ | Judul | Pemilik |
| --- | --- | --- |
| OQ-001 | Tidak ada DDL Oracle / definisi properti | DBA |
| OQ-002 | 67 stored procedure tanpa body | DBA |
| OQ-008 | Arti prefix `ASM` / `RNM` / `GCNM` | DBA |
| OQ-012 | Struktur JSON dalam kolom | DBA + Product+UW |
| OQ-013 | `COMMIT` & identitas pengguna di dalam procedure | DBA |
| OQ-016 | Peran 11 skema Oracle, mana yang dalam lingkup | DBA (+ Product+UW) |
| OQ-017 | Database link `ASMD.SINARMAS.CO.ID` | DBA + Actuarial + Product+UW |
| OQ-007 | Tidak ada rule identitas/otorisasi di korpus | IAM |
| OQ-021 | Identitas orang sebagai guard otorisasi | IAM (+ Product+UW) |
| OQ-003 | Folder `excludeXML` — aktif atau tidak | Product+UW |
| OQ-009 | Rule di class bawaan Pega — dalam lingkup? | Product+UW |
| OQ-010 | `Treaty In` vs `Treaty In Adjustment` | Product+UW |
| OQ-014 | Arti singkatan `EDM` | Product+UW |
| OQ-015 | Trio facultative: satu ruleset bersiklus? | Product+UW |
| OQ-019 | Trio Claim: mengapa terpisah? | Product+UW |

### 5.3 Ringkasan beban per pemilik

| Pemilik | Pemilik tunggal | Pemilik bersama | OQ |
| --- | ---: | ---: | --- |
| Product+Underwriting | 8 | 6 | tunggal: 003, 005, 009, 010, 014, 015, 019, 022 · bersama: 011, 012, 016, 017, 020, 021 |
| DBA | 4 | 4 | tunggal: 001, 002, 008, 013 · bersama: 012, 016, 017, 018 |
| IAM | 1 | 2 | tunggal: 007 · bersama: 018, 021 |
| Finance | 0 | 1 | bersama: 020 |
| Actuarial | 0 | 1 | bersama: 017 |

Total 20 OQ terbuka. **Product+Underwriting adalah penentu terbesar** — 14 dari 20 OQ menyentuh
peran ini, termasuk empat dari lima pemblokir D2.

---

## 6. Rekomendasi urutan D2 — berbasis bukti

D2 menelusur flow end-to-end. Urutan di bawah disusun dari **ketersediaan titik masuk** dan
**tingkat konflik**, bukan dari kepentingan bisnis (itu keputusan manusia).

### 6.1 Titik masuk yang tersedia

`[terverifikasi]` **24 rule `Flow` di 15 modul.** **5 modul tanpa `Flow`**: Treaty In,
Treaty In Adjustment, Treaty Contract Out, Master Product Name Life, Master Contract Retro Life.

### 6.2 Urutan yang direkomendasikan

**Tahap 1 — mulai dari Flow, konteks paling bersih (rendah konflik, satu Flow, satu domain).**

| # | Konteks | Titik masuk | Alasan |
| ---: | --- | --- | --- |
| 1 | **NB Treaty In** | `Flow/InputRealizationTreatyIn.xml` → `Start1` (`ASM-FW-GISFW-WORK / INPUTREALIZATIONTREATYIN`) | 278 file, satu Flow, identik dengan salinannya di NB FacIn. Konteks treaty inward terkecil yang utuh. |
| 2 | **EDM Treaty In** | `Flow/InputAddendumTreatyIn.xml` → `Start1` | 163 file, satu Flow; melengkapi siklus treaty (addendum). Menjawab sebagian OQ-014 (`EDM`). |
| 3 | **Endorsement Life** | `Flow/InputEDMLife.xml` → `Start1` | 75 file, modul terkecil ber-Flow. Menguji `EdmType` (OQ-020). |
| 4 | **PremiumList Life** | `InputPolicyHolder.xml` → `Start1` (di root modul, lihat OQ-004) | 124 file, terpisah dari modul lain; menguji kode `QR/QP/TP/TR`. |

**Tahap 2 — komite klaim: kecil, tetapi konflik guard paling parah.**

| # | Konteks | Titik masuk | Catatan |
| ---: | --- | --- | --- |
| 5 | **Komite Claim Life** | `Flow/KomiteLife_Flow.xml` → `Start2` | 47 file, modul terkecil. Tipe Flow dari `<pxObjClass>`. |
| 6 | **Komite Claim FacIn** | `Flow/Komite_Flow.xml` → `Start1` | 114 file. |
| 7 | **Komite Claim Prop / Non Prop** | `Flow/KomiteTreaty_Flow.xml` (class berbeda = **dua rule berbeda**) | Wajib menelusur **keduanya**; di sini juga terdapat guard identitas orang (OQ-021) di `KOMITEPOST_REJECT`. |

**Tahap 3 — klaim: tiga basis kode terpisah, konflik `When` 6 versi.**

| # | Konteks | Titik masuk | Catatan |
| ---: | --- | --- | --- |
| 8 | **Claim Life** | `Flow/Register_Flow.xml` → `Start2` | Class `...WORK-CLAIMLIFE`; berbeda dari `Register_Flow` Claim Fac In. |
| 9 | **Claim Prop** | `Flow/Flow_TreatyIn.xml` → `Start1` | Class `...WORK-CLAIMTREATY`. Kode `PaymentType` (OQ-020). |
| 10 | **Claim Non Prop** | `Flow/Flow_TreatyIn.xml` → `Start1` | Class `...WORK-CLAIMTREATYNONPROP` — **rule berbeda** dari Claim Prop. |
| 11 | **Claim Fac In** | `Flow/Register_Flow.xml` → `Start1` | Class `...WORK-PNC`. 482 file. |

**Peringatan tahap 3:** `ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` punya **6 isi berbeda di 6
modul**. Setiap kali salah satu muncul sebagai guard, **baca versi milik modul yang sedang
ditelusur** dan catat bahwa versinya berbeda dari modul lain.

**Tahap 4 — facultative inward: 6.071 file, tetapi kemungkinan satu ruleset.**

| # | Konteks | Titik masuk | Catatan |
| ---: | --- | --- | --- |
| 12 | **NB FacIn** | 6 Flow: `InputInwardFacultativeOffer`, `InputInwardFacultativeRISlip`, `InputQuotation`, `InputRealizationTreatyIn`, `OfferFacOut` (`Start62`), `OfferFacRetro` | Modul paling kaya; telusur di sini menutup sebagian besar domain. |
| 13 | **RNW Fac In** | 4 Flow, 3 di antaranya **identik** dengan NB FacIn | Fokus pada **`InputRenewalFacultativeIn`** (`Start2`) — satu-satunya yang khas renewal. |
| 14 | **Endorsment Fac In** | `InputAddendumFacultativeIn` (`Start2`) + **`OfferFacRetro` versi berbeda** | `OfferFacRetro` di sini **berbeda isi** dari NB/RNW → telusur terpisah (OQ-011). |

**Tahap 5 — modul tanpa `Flow`: metode telusur harus ditetapkan lebih dulu.**

| # | Konteks | Kandidat titik masuk | Prasyarat |
| ---: | --- | --- | --- |
| 15 | **Treaty In** + **Treaty In Adjustment** | `Harness` (3 / 6), `FlowAction` (32 / 43) | Jawab **OQ-010** dulu — 320 identitas dibagi; kemungkinan satu aplikasi. |
| 16 | **Master Product Name Life**, **Master Contract Retro Life** | `Harness` (1 / 4), `FlowAction` (11 / 2) | Modul master; kemungkinan CRUD tanpa workflow. |
| 17 | **Treaty Contract Out** | `Harness` (3), `FlowAction` (2) | **Jawab OQ-022 dulu.** 67% rule di `@BASECLASS`; premis "outward" tidak didukung bukti. |

### 6.3 Aturan versi rule saat menelusur (kaitan OQ-011)

Mengikat untuk seluruh D2:

1. **Sebelum membuka rule, cek `inventory/_oq011-konflik-isi.md`.** Bila identitasnya terdaftar,
   rule itu punya lebih dari satu isi di korpus.
2. **Selalu sebut path file yang dibaca**, bukan hanya nama rule. `path + class/nama/tipe`.
3. Untuk identitas berkonflik: telusur **varian milik modul konteks yang sedang dikerjakan**, dan
   catat eksplisit bahwa modul lain punya versi berbeda. **Jangan menyatakan "rule X melakukan Y"**
   tanpa menyebut varian mana.
4. Bila satu konteks bergantung pada rule berkonflik yang variannya **tidak ada** di modul itu
   sendiri, hentikan dan catat sebagai blocker — jangan meminjam varian dari modul lain.
5. Pemanggilan stored procedure dan objek db-link dicatat sebagai **batas pengetahuan**, bukan
   langkah yang dijelaskan.

---

## 7. Yang belum dikerjakan di D1

- **Rule tertanam** (`<pyIncludedRuleXML>`): diperiksa hanya untuk batch 2 — 208 file memuat rule
  tambahan, dan hanya **3 rule aplikasi** yang benar-benar belum masuk inventaris. Batch 1, 3, 4
  **belum diperiksa**.
- **Tabel yang disentuh lewat `ReportDefinition` dan `DataPage`** belum diekstraksi — baru `RDBList`.
- **Dependency antar-Section** (`<pyInclude>`) baru dihitung jumlahnya, belum dipetakan menjadi graf.
- **Isi rule `When`, `DecisionTable`, `DataTransform`, `DecisionTree`** belum dibaca — hanya didaftar.
- **Daftar 18 tag normalisasi mungkin belum lengkap**; bila ada tag provenance lain yang lolos,
  angka 533 dapat turun lagi.

---

## 8. Gate

STEP D2 boleh dimulai. **FASE B (ADR/BRD/spec/tiket) tetap tidak boleh dimulai** sebelum D4 selesai
dan di-review manusia, sesuai `README.md` §9.
