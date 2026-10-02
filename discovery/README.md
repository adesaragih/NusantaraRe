# FASE A — DISCOVERY / PEMAHAMAN APLIKASI

Status fase: **SELESAI — FASE A DITUTUP** (D0, D1, D2, D3, D4 seluruhnya selesai).
D1: 20/20 modul, **9.369/9.369 file** (100 %). D2: **20/20 konteks, 428 rule**, 5 tahap.
D3: **20 catatan modul + `glossary.md`** (169 entri). D4: **`context-map.md`**
(9 bounded context + 1 absent) + **`understanding-report.md`**.
Pertanyaan terbuka: **57 terbuka / 2 terjawab** — **38 memblokir FASE B** (`open-questions.md`).

> **MENUNGGU GATE MANUSIA.** FASE B (ADR/BRD/spec/tiket) **tidak boleh dimulai** sebelum gate.
> Paket review + urutan FASE B + cara menjalankan skill secara manual:
> **`D3-D4-CLOSING-REPORT.md`**.

Tanggal mulai: 2026-09-12. Tanggal tutup FASE A: 2026-09-14.

## 1. Tujuan fase

Menghasilkan **peta pemahaman lengkap** atas aplikasi reasuransi Nusantara Re yang berjalan di
Pega, cukup detail sehingga engineer/agent lain dapat memahami cara kerja sistem **tanpa membuka
Pega** — seluruhnya berbukti `path file + nama rule`.

**Yang BUKAN tujuan fase ini (non-goal):**

- Tidak membuat keputusan desain, ADR, BRD, spec, atau tiket. Itu FASE B.
- Tidak merancang skema Oracle baru, API Go, atau komponen React.
- Tidak menilai apakah perilaku existing "benar" atau "salah" — fase ini hanya merekam apa adanya.
- Tidak menulis kode.

## 2. Sumber kebenaran & batas tulis

| Hal | Lokasi | Sifat |
| --- | --- | --- |
| Korpus XML Pega | `D:\XML\RNM_BRD\` (20 modul) | **READ-ONLY** — dilarang tulis/ubah/hapus |
| Panduan metodologi | `D:\XML\RNM_BRD\1. Agentic Development Methodology\`, `D:\XML\RNM_BRD\.runbook.txt` | READ-ONLY |
| Seluruh output artefak | `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\` | satu-satunya folder yang boleh ditulis |

## 3. Aturan main FASE A (mengikat)

1. **Setiap klaim wajib bukti** berupa `path file + nama rule`.
   Contoh sah: `NB Treaty In/Flow/InputRealizationTreatyIn.xml` → rule `InputRealizationTreatyIn`.
   Klaim tanpa bukti tidak boleh ditulis sebagai fakta.
2. **Identitas rule ditulis sebagai `class / nama / tipe`**, tidak pernah nama saja. Nama rule
   tidak unik lintas tipe maupun lintas class.
   **Sumber identitas adalah `<pxInsName>`, BUKAN `<pzOriginalInstanceKey>`** — lihat §6.4.
   - `class` dan `nama` → dari `<pxInsName>` (format `CLASS!NAMA`; untuk RDBList
     `CLASS!PREFIX!NAMA`).
   - `tipe` → dari field pertama `<pzOriginalInstanceKey>` (mis. `RULE-OBJ-ACTIVITY`), atau dari
     `<pxObjClass>` bila tag itu tidak ada. **Bukan dari nama folder** — batch 2 membuktikan folder
     dapat berbeda dari tipe rule sebenarnya (§8.4.5).
   - `<pzOriginalInstanceKey>` selebihnya dicatat sebagai **asal salinan**, bukan identitas.
     Tag ini **tidak selalu ada**: 5 file batch 2 sama sekali tidak memilikinya (§8.4.5).
3. **Grep dulu, baca range belakangan.** Korpus berukuran **1.319,3 MB** dengan file terbesar
   **8,2 MB** (`Treaty In Adjustment/Harness/InputTreatyInOffer.xml`). Membaca file utuh akan
   menghabiskan konteks. Pola kerja: grep tag untuk menemukan posisi, lalu baca rentang baris
   seperlunya.
4. **Jangan menyimpulkan perilaku dari nama rule saja** — penamaan Pega di korpus ini terbukti
   menipu. Bukti terukur dari D1 batch 1 (`inventory/_summary.md` §4–§5):
   63 nama file dipakai di lebih dari satu tipe rule; satu identitas `class + nama` bahkan dipakai
   dua tipe sekaligus; dan sebuah rule `SELECT` bernama `GETCURRENCY` ternyata hasil salinan dari
   rule bernama `UPDATEMASTERCURRENCY`.
5. **Jangan menebak skema/tipe data.** Tidak ada DDL maupun folder `Property/` di export. Tipe
   kolom, panjang, nullability, dan constraint **tidak diketahui** dan tidak boleh dikarang.
6. **Jangan menebak arti singkatan** atau arti sebuah field dari caption di sebelahnya.
   Singkatan yang tidak dijabarkan korpus ditulis: `kepanjangan belum terverifikasi`.
7. **Jangan menebak arti kode/status literal.** Nilai seperti status atau kode dicatat apa adanya;
   artinya ditulis `arti belum terverifikasi` bila korpus tidak menjelaskannya.
8. **Yang tidak pasti masuk `open-questions.md`**, bukan ditutup dengan asumsi. Setiap pertanyaan
   wajib punya **pemilik peran** (DBA / Product+Underwriting / Actuarial / Finance / IAM) dan
   keterangan **apa yang diblokir** olehnya.
9. **Label wajib pada setiap pernyataan:**
   - `[terverifikasi]` — didukung tag XML yang dikutip langsung;
   - `[dugaan]` — hanya dari nama file/rule atau pola, belum dikonfirmasi tag;
   - `[pertanyaan terbuka]` — tidak dapat dipastikan dari korpus.
10. **Angka wajib dapat diaudit ulang.** Setiap angka disertai perintah yang menghasilkannya.
    Yang tidak terukur ditulis `belum terukur`.

## 4. Struktur folder discovery

```
discovery/
├── README.md               <- dokumen ini: tujuan fase, aturan bukti, status per modul
├── inventory/              <- STEP D1: inventaris rule per modul (per tipe)
│   ├── _summary.md              rekap lintas modul per batch
│   └── _oq011-konflik-isi.md    register identitas yang isinya berbeda antar modul (OQ-011)
├── flows/                  <- STEP D2: telusur flow end-to-end per konteks
│   └── _METHOD.md               konvensi telusur D2 (mengikat)
├── modules/                <- STEP D3: catatan pemahaman per modul
├── glossary-draft.md       <- STEP D3: istilah domain yang DITEMUKAN (draft, belum final)
├── open-questions.md       <- lintas step: hal yang tak bisa dipastikan, berpemilik peran
├── context-map.md          <- STEP D4: bounded context + cakupan bukti
└── understanding-report.md <- STEP D4: sintesis naratif cara kerja sistem
```

## 5. Fakta korpus terukur (STEP D0)

Semua angka di bawah dihitung pada 2026-09-12 dan dapat diaudit ulang dengan perintah yang
disertakan. Dijalankan dari `D:\XML\RNM_BRD\`.

- **Jumlah modul: 20.**
- **Jumlah file XML: 9.369.**
  `find . -type f -name "*.xml" -not -path "./OUTPUT_HASIL_RNM/*" | wc -l`
- **Ukuran korpus: 1.319,3 MB.**
  `find . -type f -name "*.xml" -not -path "./OUTPUT_HASIL_RNM/*" -printf "%s\n"` lalu dijumlahkan.
- **9.368 file berada di bawah folder tipe rule**; **1 file berada di luar struktur tipe**, yaitu
  `PremiumList Life/InputPolicyHolder.xml` (langsung di root modul, tanpa folder tipe rule).
  **Terjawab di D1 batch 4:** file itu rule `RULE-OBJ-FLOW`
  (`ASM-FW-GISFW-WORK-LIFE / INPUTPOLICYHOLDER`) — lihat OQ-004 pada `open-questions.md` §Terjawab.
  Konsekuensinya matriks §5.1 di bawah (yang menghitung **per folder tipe**) mencatat `PremiumList
  Life` = 123, sedangkan jumlah file modul itu sebenarnya **124**.

### 5.1 Matriks modul × tipe rule (status per modul)

Perintah audit:

```
find . -type f -name "*.xml" -not -path "./OUTPUT_HASIL_RNM/*" \
  | awk -F/ 'NF>=4{print $2"\t"$3}' | sort | uniq -c
```

| Modul | Activity | ConnectREST | DataPage | DataTransform | DecisionTable | DecisionTree | Flow | FlowAction | HTMLRule | Harness | Menu | RDBList | ReportDefinition | Section | SystemSettings | When | excludeXML | TOTAL |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| Claim Fac In | 179 | 8 | 8 | 28 | 1 | - | 1 | 33 | - | 16 | - | 63 | 27 | 57 | 1 | 60 | - | 482 |
| Claim Life | 51 | 2 | - | 1 | 1 | - | 1 | 16 | - | 3 | - | 29 | 8 | 20 | 1 | 3 | - | 136 |
| Claim Non Prop | 114 | 7 | - | 10 | 1 | - | 1 | 16 | - | 14 | - | 50 | 22 | 36 | 1 | 7 | - | 279 |
| Claim Prop | 112 | 6 | - | 11 | 1 | - | 1 | 12 | - | 11 | - | 54 | 19 | 36 | 1 | 6 | - | 270 |
| EDM Treaty In | 66 | 1 | - | 10 | 1 | - | 1 | 6 | - | 3 | - | 36 | 9 | 18 | 1 | 11 | - | 163 |
| Endorsement Life | 16 | 1 | - | 2 | 1 | - | 1 | 6 | - | 8 | - | 17 | 6 | 14 | 1 | 2 | - | 75 |
| Endorsment Fac In | 582 | 3 | 41 | 157 | 10 | 1 | 2 | 265 | - | 37 | - | 188 | 138 | 434 | 1 | 202 | - | 2061 |
| Komite Claim FacIn | 28 | 3 | - | 1 | 1 | - | 1 | 5 | - | - | - | 17 | 3 | 5 | 1 | 49 | - | 114 |
| Komite Claim Life | 18 | 1 | - | - | 1 | - | 1 | 3 | - | - | - | 14 | 3 | 3 | 1 | 2 | - | 47 |
| Komite Claim Non Prop | 20 | 5 | - | 1 | 1 | - | 1 | 2 | - | - | - | 16 | 3 | 2 | 1 | 6 | 1 | 59 |
| Komite Claim Prop | 35 | 3 | - | 2 | 1 | - | 1 | 2 | - | - | - | 23 | 4 | 2 | 1 | 6 | - | 80 |
| Master Contract Retro Life | 27 | - | - | 1 | - | - | - | 2 | - | 4 | - | 12 | 12 | 8 | - | - | - | 66 |
| Master Product Name Life | 37 | 1 | - | 10 | 1 | - | - | 11 | - | 1 | - | 22 | 17 | 12 | 1 | 1 | - | 114 |
| NB FacIn | 609 | 3 | 30 | 148 | 12 | 1 | 6 | 250 | - | 41 | - | 217 | 123 | 432 | 1 | 210 | - | 2083 |
| NB Treaty In | 92 | - | - | 12 | 2 | - | 1 | 9 | - | 6 | - | 41 | 15 | 25 | - | 75 | - | 278 |
| PremiumList Life | 37 | 1 | - | 7 | 2 | - | - | 10 | 1 | 9 | - | 25 | 9 | 19 | 1 | 2 | - | 123 |
| RNW Fac In | 556 | 3 | 30 | 138 | 10 | 1 | 4 | 239 | - | 41 | - | 200 | 118 | 397 | 1 | 189 | - | 1927 |
| Treaty Contract Out | 168 | 1 | 2 | 10 | - | - | - | 2 | - | 3 | - | 37 | 30 | 49 | 1 | - | - | 303 |
| Treaty In | 149 | 1 | - | 35 | 1 | - | - | 32 | - | 3 | - | 40 | 16 | 50 | 1 | 1 | - | 329 |
| Treaty In Adjustment | 155 | 1 | - | 38 | 1 | - | - | 43 | - | 6 | 1 | 44 | 18 | 68 | 1 | 3 | - | 379 |
| **TOTAL** | **3051** | **51** | **111** | **622** | **49** | **3** | **23** | **964** | **1** | **206** | **1** | **1145** | **600** | **1687** | **18** | **835** | **1** | **9368** |

### 5.2 Status pengerjaan per modul

Semua modul berstatus **belum dikerjakan** pada akhir D0. Kolom diisi saat D1–D3 berjalan.

| Modul | D1 inventaris | D2 flow | D3 catatan modul |
| --- | --- | --- | --- |
| Claim Fac In | **selesai** | belum | belum |
| Claim Life | **selesai** | belum | belum |
| Claim Non Prop | **selesai** | belum | belum |
| Claim Prop | **selesai** | belum | belum |
| EDM Treaty In | **selesai** | belum | belum |
| Endorsement Life | **selesai** | belum | belum |
| Endorsment Fac In | **selesai** | belum | belum |
| Komite Claim FacIn | **selesai** | belum | belum |
| Komite Claim Life | **selesai** | belum | belum |
| Komite Claim Non Prop | **selesai** | belum | belum |
| Komite Claim Prop | **selesai** | belum | belum |
| Master Contract Retro Life | **selesai** | belum (tanpa Flow) | belum |
| Master Product Name Life | **selesai** | belum (tanpa Flow) | belum |
| NB FacIn | **selesai** | belum | belum |
| NB Treaty In | **selesai** | belum | belum |
| PremiumList Life | **selesai** | belum (punya Flow, lihat OQ-004) | belum |
| RNW Fac In | **selesai** | belum | belum |
| Treaty Contract Out | **selesai** | belum (tanpa Flow) | belum |
| Treaty In | **selesai** | belum (tanpa Flow) | belum |
| Treaty In Adjustment | **selesai** | belum (tanpa Flow) | belum |

## 6. Koreksi terhadap asumsi awal (temuan D0)

Dua asumsi dalam dokumen prompt awal (`../PROMPT-AWAL-CLAUDE.md`) tidak sesuai korpus nyata.
Dicatat di sini agar D1–D4 tidak bekerja di atas premis yang keliru.

### 6.1 Tipe rule lebih banyak dari yang didaftar

Prompt awal menyebut 10 tipe (`Activity`, `DataTransform`, `DecisionTable`, `Flow`, `FlowAction`,
`Harness`, `RDBList`, `ReportDefinition`, `Section`, `When`). Korpus nyata memiliki **17 folder
tipe**. Tujuh tipe berikut **tidak tercakup** prompt awal dan tetap harus diinventarisasi di D1:

| Tipe | Jumlah | Lokasi contoh |
| --- | ---: | --- |
| `DataPage` | 111 | `Claim Fac In/DataPage/`, `NB FacIn/DataPage/` |
| `ConnectREST` | 51 | tersebar di 17 modul — kandidat integrasi eksternal |
| `SystemSettings` | 18 | 18 modul, masing-masing 1 file |
| `DecisionTree` | 3 | `NB FacIn/DecisionTree/Tree_ShortPeriod.xml`, `RNW Fac In/DecisionTree/Tree_ShortPeriod.xml`, `Endorsment Fac In/DecisionTree/Tree_ShortPeriod.xml` |
| `Menu` | 1 | `Treaty In Adjustment/Menu/MasterNavTreaty.xml` |
| `HTMLRule` | 1 | `PremiumList Life/HTMLRule/GeneratePdfOfferLife.xml` |
| `excludeXML` | 1 | `Komite Claim Non Prop/excludeXML/GetBase64Attachment.xml` — nama folder mengindikasikan file sengaja dikecualikan `[pertanyaan terbuka]` |

`ConnectREST` (51 file) penting: prompt awal mengasumsikan integrasi eksternal hanya terlihat dari
dalam Activity. Korpus punya tipe rule khusus untuk itu.

### 6.2 Enam modul tidak punya satu pun rule `Flow`

Hanya **23 file Flow** di seluruh korpus, dan **6 dari 20 modul tidak memilikinya**:

- Master Contract Retro Life
- Master Product Name Life
- PremiumList Life
- Treaty Contract Out
- Treaty In
- Treaty In Adjustment

Konsekuensi untuk STEP D2: metode telusur "mulai dari `<pyStartActivity>` pada Flow" **tidak dapat
dipakai** pada enam modul tersebut. Modul-modul ini perlu titik masuk alternatif — kandidatnya
`Harness` dan `FlowAction` — dan metodenya harus ditetapkan di awal D2, bukan diimprovisasi.

Catatan tambahan: prompt awal menyarankan D2 dimulai dari `Treaty In` sebagai konteks "paling
utuh", padahal `Treaty In` justru tanpa Flow. Titik masuk Flow yang benar-benar ada untuk domain
treaty inward adalah `NB Treaty In/Flow/InputRealizationTreatyIn.xml` (1 Flow).

### 6.3 Daftar rule Flow per folder `Flow/` (23 file) — angka korpus sebenarnya 24

Perintah audit: `find . -type f -name "*.xml" -path "*/Flow/*" | sort`

| Modul | File Flow |
| --- | --- |
| Claim Fac In | `Flow/Register_Flow.xml` |
| Claim Life | `Flow/Register_Flow.xml` |
| Claim Non Prop | `Flow/Flow_TreatyIn.xml` |
| Claim Prop | `Flow/Flow_TreatyIn.xml` |
| EDM Treaty In | `Flow/InputAddendumTreatyIn.xml` |
| Endorsement Life | `Flow/InputEDMLife.xml` |
| Endorsment Fac In | `Flow/InputAddendumFacultativeIn.xml`, `Flow/OfferFacRetro.xml` |
| Komite Claim FacIn | `Flow/Komite_Flow.xml` |
| Komite Claim Life | `Flow/KomiteLife_Flow.xml` |
| Komite Claim Non Prop | `Flow/KomiteTreaty_Flow.xml` |
| Komite Claim Prop | `Flow/KomiteTreaty_Flow.xml` |
| NB FacIn | `Flow/InputInwardFacultativeOffer.xml`, `Flow/InputInwardFacultativeRISlip.xml`, `Flow/InputQuotation.xml`, `Flow/InputRealizationTreatyIn.xml`, `Flow/OfferFacOut.xml`, `Flow/OfferFacRetro.xml` |
| NB Treaty In | `Flow/InputRealizationTreatyIn.xml` |
| RNW Fac In | `Flow/InputInwardFacultativeRISlip.xml`, `Flow/InputRenewalFacultativeIn.xml`, `Flow/OfferFacOut.xml`, `Flow/OfferFacRetro.xml` |

Nama file Flow yang sama muncul di lebih dari satu modul (`Register_Flow`, `Flow_TreatyIn`,
`OfferFacRetro`, `OfferFacOut`, `InputInwardFacultativeRISlip`, `InputRealizationTreatyIn`,
`KomiteTreaty_Flow`). Apakah isinya identik atau bercabang **belum terukur** — ini pekerjaan D1.

### 6.4 `<pzOriginalInstanceKey>` BUKAN identitas rule — melainkan asal salinan

Ini koreksi paling penting dari D0 dan **membatalkan aturan nomor 4 pada prompt awal**
(`Identitas rule = class/nama/tipe, ambil dari <pzOriginalInstanceKey>`).

**Bukti terukur** atas 1.149 file batch 1 (NB Treaty In, Treaty In, Treaty In Adjustment,
EDM Treaty In). Seluruh 1.149 file memiliki kedua tag; tidak ada `<pxInsName>` yang kosong.

| Pengukuran | Hasil |
| --- | ---: |
| Nama pada `<pzOriginalInstanceKey>` **berbeda** dari nama pada `<pxInsName>` | **410 file (35,7%)** |
| Class pada `<pzOriginalInstanceKey>` **berbeda** dari class pada `<pxInsName>` | **62 file (5,4%)** |
| Nama file cocok dengan nama di `<pxInsName>` | **1.149 file (100%)** |
| Nama file cocok dengan nama di `<pzOriginalInstanceKey>` | 696 file (60,6%) |

Perintah audit: ekstrak `<pzOriginalInstanceKey>` dan `<pxInsName>` pertama tiap file, lalu
bandingkan bagian nama terhadap nama file (uppercase, tanpa `.xml`). Untuk RDBList, bagian nama
adalah segmen **setelah prefix** (lihat §6.5).

**Contoh nyata:**

- `NB Treaty In/Activity/AgentSourceBizTreatyIn_Act.xml`
  `<pxInsName>` = `ASM-FW-GISFW-DATA-AGENT!AGENTSOURCEBIZTREATYIN_ACT` ← identitas sebenarnya
  `<pzOriginalInstanceKey>` = `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-AGENT AGENTSOURCEBIZ_ACT` ← rule asal
- `NB Treaty In/Activity/TreatyRealizationCheckXOLList.xml`
  `<pxInsName>` class = `ASM-FW-GISFW-WORK`
  `<pzOriginalInstanceKey>` class = `ASM-SFAGIS-WORK` ← class yang bahkan berbeda framework

**Kesimpulan `[terverifikasi]`:** `<pzOriginalInstanceKey>` adalah kunci instance rule yang
di-"Save As" menjadi rule ini. Memakainya sebagai identitas akan **salah melabeli 410 dari 1.149
rule** di batch 1 saja.

**Nilai gunanya tetap besar, tapi untuk hal lain:** tag ini memperlihatkan **silsilah salinan**.
Contoh keluarga klon yang terbaca langsung — keempatnya berasal dari
`ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT`:

- `NB Treaty In/Activity/CheckSpreadingProtectFire_ACT.xml`
- `NB Treaty In/Activity/CheckSpreadingProtectAnekaGolf_ACT.xml`
- `NB Treaty In/Activity/CheckSpreadingProtectMCargoMBU_ACT.xml`
- `NB Treaty In/Activity/CheckSpreadingProtectPATravel_ACT.xml`

Pola ini adalah sinyal migrasi yang berharga (kandidat kolaps menjadi satu fungsi berparameter),
sehingga D1 **wajib tetap merekam** `<pzOriginalInstanceKey>` — dengan label kolom
**"asal salinan"**, bukan "identitas".

### 6.5 RDBList memakai kunci tiga bagian: `CLASS!PREFIX!NAMA`

Dari 1.149 file batch 1, **161 file** memiliki `<pxInsName>` yang tidak cocok langsung dengan nama
file. **Seluruh 161 file tersebut bertipe `RDBList`** (`RULE-CONNECT-SQL`), dan seluruhnya
mengikuti pola `CLASS!PREFIX!NAMA` — tidak ada satu pun yang merupakan rename sesungguhnya.

Dua prefix yang muncul:

| Prefix | Jumlah (batch 1) |
| --- | ---: |
| `ASM` | 89 |
| `RNM` | 72 |

Contoh: `NB Treaty In/RDBList/BrowseTreatyIn.xml` → `ASM-FW-GISFW-INT-TREATY_IN!ASM!BROWSETREATYIN`,
sedangkan `NB Treaty In/RDBList/BrowseTreatyOut.xml` → `...!RNM!BROWSETREATYOUT`.

Setelah prefix diperhitungkan, kecocokan nama file dengan `<pxInsName>` menjadi **100%**.

**Arti prefix `ASM` dan `RNM` belum terverifikasi** — dugaan kuat bahwa keduanya menunjuk
koneksi/database yang berbeda, tetapi korpus tidak menjelaskannya. Dicatat sebagai **OQ-008**
(pemilik: DBA). Bila benar ada dua database, ini berdampak langsung pada desain repository.

## 7. Tag XML kunci (terverifikasi pada file nyata)

Diverifikasi pada `NB Treaty In/Flow/InputRealizationTreatyIn.xml` (287.745 byte, 7.863 baris):

| Tag | Isi terverifikasi | Guna |
| --- | --- | --- |
| `<pxInsName>` | `ASM-FW-GISFW-WORK!INPUTREALIZATIONTREATYIN` | **sumber otoritatif class + nama** rule (lihat §6.4) |
| `<pzOriginalInstanceKey>` | `RULE-OBJ-FLOW ASM-FW-GISFW-WORK INPUTREALIZATIONTREATYIN #20221108T060345.907 GMT` | **tipe rule** (field pertama) + **asal salinan** (field 2–3). **Bukan identitas** |
| `<pyStartActivity>` | `Start1` | activity awal sebuah Flow |

Pada file contoh ini kedua tag kebetulan sepakat. Pada 410 dari 1.149 file batch 1 **tidak** —
itulah sebabnya `<pxInsName>` yang dipakai sebagai identitas.

Tag lain yang dipakai menelusur (akan diverifikasi per pemakaian di D1–D2):

- Activity: `<pyActivityType>`, `<pyActivityName>`, `<pyStepsActivityName>` (method tiap step),
  `<pyStepsPreCondParamsWhen>` (precondition), `<pyStepsJavaSource>` (Java tertanam).
- RDBList: `<pyBrowseSQL>` (SQL asli).
- Section/Harness: `<pyInclude>` (embed = dependency nyata), `<pyType>` bernilai
  `FIELD` / `LAYOUT` / `SUB_SECTION`.
- **`<pySection>` di aksi Refresh BUKAN dependency** — jangan dihitung sebagai ketergantungan.

## 8. Keputusan operasional (ditetapkan manusia, 2026-09-12; §8.1 diganti 2026-09-14)

### 8.1 `D:\XML\nusantara-re\` — ⛔ DI-BLACKLIST (keputusan work owner, 2026-09-14)

> **Status berlaku: DILARANG TOTAL.** Jangan dibaca, jangan dijadikan pembanding, jangan dijadikan
> repo target implementasi. Perlakukan seolah tidak ada.

**Keputusan work owner, 2026-09-14 — "Opsi A": repo target tunggal adalah
`D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\`.** Seluruh keluaran — dokumen (discovery, spec, ADR, tiket)
**dan** kode Go/React nanti — ditulis di sana. Korpus `D:\XML\RNM_BRD\` (20 modul XML) tetap
READ-ONLY sebagai **satu-satunya** sumber kebenaran Pega.

**Alasan:** dua sumber dokumen yang sama-sama mengaku otoritatif menimbulkan dua bahaya nyata yang
sudah terbukti muncul: (a) dua set ADR dengan **nomor bentrok** — `ADR-000`…`ADR-017` di sana versus
`ADR-0001`…`ADR-0012` di sini, sebagian bertopik sama; dan (b) klaim yang bersandar pada dokumen di
luar korpus sehingga **tidak dapat diverifikasi** ke `path + rule`. Satu repo target menghapus
keduanya.

#### Jejak aturan lama (untuk audit — sudah TIDAK berlaku)

Aturan yang berlaku **2026-09-12 s/d 2026-09-14**, dicatat agar artefak lama dapat dibaca dalam
konteksnya:

> Terdapat hasil kerja sebelumnya di luar korpus: `CLAUDE.md`, `CONTEXT.md`, `docs/context-map.md`,
> `docs/adr/ADR-000` s/d `ADR-017`, `docs/fase1`–`fase4`,
> `docs/migration/name-collisions.json`, serta export `.docx`.
>
> **Perlakuan yang ditetapkan (lama):** discovery tetap dikerjakan dari nol terhadap korpus XML.
> `nusantara-re` boleh dibuka hanya untuk **membandingkan** temuan. Setiap klaim dari sana wajib
> diverifikasi ulang ke korpus (`path + rule`) sebelum boleh masuk artefak discovery. Dilarang
> menyalin pernyataan dari sana sebagai fakta. Bila temuan discovery berbeda, perbedaan dicatat —
> korpus yang menang. `nusantara-re` juga READ-ONLY.

**Yang berubah:** izin "boleh dibuka untuk membandingkan" **dicabut**. Kewajiban verifikasi ke
korpus tidak berubah — ia justru kini menjadi **satu-satunya** jalan.

#### Akibat pada artefak yang sudah ada

- Artefak FASE A (D0–D4) **tidak terpengaruh**: seluruhnya dibangun dari korpus XML, dan tidak satu
  pun klaimnya bersandar pada `nusantara-re`. Diperiksa 2026-09-14: `discovery/open-questions.md`
  dan `docs/adr/` **nol** rujukan.
- `.scratch/claim-life/spec.md` **sudah dibersihkan** pada 2026-09-14. Klaim yang tadinya bersandar
  pada dokumen di sana diverifikasi ulang ke korpus bila memungkinkan; yang tidak, diturunkan
  menjadi pertanyaan terbuka berpemilik. Rinciannya di kepala spec.
- Berkas **historis** — `PROMPT-*.md`, `D2-CLOSING-REPORT.md`, `D3-D4-CLOSING-REPORT.md`,
  `.scratch/claim-life/grilling-ronde-1.md` — **sengaja tidak diubah**. Ia merekam instruksi dan
  verifikasi apa adanya pada tanggalnya; menulis ulang akan merusak jejak audit. Bacalah di bawah
  aturan yang berlaku saat itu.

### 8.2 D1 dikerjakan bertahap per domain

Bukan sekaligus 20 modul. Urutan batch yang ditetapkan:

| Batch | Domain | Modul | File | Status |
| ---: | --- | --- | ---: | --- |
| 1 | Treaty inward | NB Treaty In, Treaty In, Treaty In Adjustment, EDM Treaty In | 1.149 | **selesai** |
| 2 | Facultative inward | Komite Claim FacIn, Claim Fac In, RNW Fac In, Endorsment Fac In, NB FacIn | 6.667 | **selesai** |
| 3 | Claim non-fac | Komite Claim Life, Komite Claim Non Prop, Komite Claim Prop, Claim Life, Claim Prop, Claim Non Prop | 871 | **selesai** |
| 4 | Life & master & outward | Endorsement Life, Master Contract Retro Life, Master Product Name Life, PremiumList Life, Treaty Contract Out | 681 | belum |

Batch 2 dikerjakan berurutan kecil → besar (Komite Claim FacIn → Claim Fac In → RNW Fac In →
Endorsment Fac In → NB FacIn), satu modul tuntas sebelum lanjut.

Tiap batch diselesaikan utuh sebelum batch berikutnya, agar ada hasil yang dapat direview lebih
awal. `inventory/_summary.md` baru final setelah keempat batch selesai.

## 8.3 Temuan utama D1 batch 1 (ringkas)

Rincian dan cara audit ada di `inventory/_summary.md`.

- 1.149 file batch 1 hanya berisi **729 identitas rule unik** — 420 file adalah salinan lintas modul.
- Dari 387 identitas yang muncul di >1 modul: **257 identik**, **130 benar-benar berbeda isinya**
  (→ OQ-011, risiko halusinasi paling konkret di korpus).
- `Treaty In` dan `Treaty In Adjustment` berbagi **320 identitas** (→ OQ-010).
- **13 stored procedure `POOLDATA.*`** dipanggil, **body-nya tidak ada di korpus** (→ OQ-002).
- Data bisnis disimpan sebagai **JSON dalam kolom `JSONDATA`**, bukan kolom ternormalisasi
  (→ OQ-012). Struktur JSON-nya tidak ada di korpus.
- `COMMIT` dilakukan **di dalam** stored procedure, dan `{OperatorID.pyUserIdentifier}` dikirim
  sebagai parameter (→ OQ-013).

## 8.4 Temuan utama D1 batch 2 (ringkas)

Rincian, angka, dan perintah audit ada di `inventory/_summary.md` §7–§13.
Cakupan: 6.667 file di Komite Claim FacIn, Claim Fac In, RNW Fac In, Endorsment Fac In, NB FacIn.
Jumlah file kelima modul **cocok** dengan §5.1 di atas.

### 8.4.1 6.667 file hanya berisi 2.970 identitas rule unik

`NB FacIn`, `RNW Fac In`, dan `Endorsment Fac In` nyaris satu ruleset: **1.567 identitas (61,4%)
ada di ketiganya**, dan **1.906 dari 1.926 identitas `RNW Fac In` (99,0%) juga ada di `NB FacIn`**.
`[dugaan]` satu ruleset dengan diskriminator siklus (baru / renewal / endorsement) — pola sama
dengan `Treaty In` ↔ `Treaty In Adjustment` di batch 1. **Bukan** penetapan bounded context; itu
STEP D4. → **OQ-015**.

### 8.4.2 OQ-011 membesar: 461 konflik baru

> **ANGKA DIKOREKSI saat batch 3** — semula dilaporkan 807. Lihat §8.5.5.

Dari 1.974 identitas lintas-modul batch 2: **1.513 identik**, **461 berbeda isi** (metode hash
ternormalisasi, bukan `md5sum` mentah).

Seluruhnya terdaftar beserta **semua** path variannya di
**`inventory/_oq011-konflik-isi.md`**. **Tidak ada satu varian pun yang dipilih sebagai benar.**
D2 tidak boleh menelusur rule yang ada di register itu sebelum OQ-011 dijawab.

### 8.4.3 Basis data jauh lebih luas dari dugaan batch 1

- **10 skema Oracle lokal**, bukan hanya `POOLDATA`/`DATAPEGA`: ditambah `NEW_UNDERWRITING`,
  `GENERAL`, `NEW_GENERAL`, `REINSURANCE`, `FIRE`, `ARASAPAS`, `MBU`, `GL` → **OQ-016**.
- **Database link ke sistem luar `ASMD.SINARMAS.CO.ID`** — 8 objek, berisi **tabel rate**
  (EQ, terorisme, banjir, FLEXAS, RSMD) dan **kurs standar**. Sebagian input perhitungan premi
  berada di luar basis data yang dimigrasikan → **OQ-017**.
- **50 stored procedure unik, 42 baru**, termasuk `FIRE.PEGA_FIRE_SET_RATE` dan
  `FIRE.CEK_PRORATA_TANGGAL` — logika rating di dalam database. Body tetap tidak ada di korpus
  → memperluas **OQ-002**.
- `JSON_POLIS` adalah objek paling sering dirujuk di batch 2 (58 rule) → menguatkan **OQ-012**.

### 8.4.4 Integrasi eksternal terkonfigurasi rapi, dengan satu penyimpangan

19 dari 20 `ConnectREST` mengambil base URL dari `SystemSettings` `LinkService`
(`=ResponLink.URL`); nilainya tidak ada di korpus. Satu rule —
`Claim Fac In/ConnectREST/getPremiumPaidOnMarine.xml` — memuat URL literal.

Namun URL literal **tersebar jauh di luar ConnectREST**: 48 file batch 2, mayoritas di `Section`,
`FlowAction`, `Activity`, `Harness`. Host yang muncul termasuk IP internal, hostname internal,
`app.sinarmas.co.id`, SaaS `view.officeapps.live.com`, dan **hostname DEV**
(`appdev.nusantarare.com`) → **OQ-018**. Sapuan korpus penuh (batch 4): **81 file di 15 modul**, `appdev` muncul 7×.

### 8.4.5 Dua koreksi metode inventaris

- **Nama folder bukan tipe rule.** `Claim Fac In/Activity/InsertLogServiceClaim.xml` sebenarnya
  `RULE-CONNECT-SQL`; `Claim Fac In/Section/CedingCedant.xml` sebenarnya `RULE-HTML-HARNESS`
  (`<pxObjClass>Rule-HTML-Harness`). Matriks §5.1 menghitung **file per folder**, bukan rule per tipe.
- **`<pzOriginalInstanceKey>` tidak selalu ada.** 5 file batch 2 tidak memilikinya sama sekali
  (`ViewObjectSavior` di 3 modul, `hasPrimaryPage` di 2 modul). Identitas tetap terbaca dari
  `<pxInsName>`; silsilah salinannya tidak dapat ditelusuri.
- **Sebagian file memuat lebih dari satu rule.** 208 file batch 2 memuat rule tertanam lewat
  `<pyIncludedRuleXML>` (643 kemunculan tambahan). Sudah diukur: hanya **3 rule aplikasi** yang
  belum masuk inventaris; 76 sisanya rule bawaan Pega. Belum diperiksa untuk batch 1.

## 8.5 Temuan utama D1 batch 3 (ringkas)

Rincian, angka, dan perintah audit ada di `inventory/_summary.md` §14–§20.
Cakupan: 871 file di Komite Claim Life, Komite Claim Non Prop, Komite Claim Prop, Claim Life,
Claim Prop, Claim Non Prop. Keenam modul **cocok** dengan §5.1.

### 8.5.1 Domain klaim BUKAN satu ruleset — kebalikan dari facultative

871 file berisi **623 identitas unik**. Overlap tertinggi di trio Claim hanya **102 identitas**
(Claim Non Prop ↔ Claim Prop, ~37%); `Claim Life` hanya berbagi 12–14% dengan keduanya.
Bandingkan batch 2, di mana `RNW Fac In` berbagi **99,0%** dengan `NB FacIn`.

`[dugaan]` pembedanya jenis bisnis (Life / Proportional / Non-Proportional), bukan tahap siklus.
**Belum terverifikasi** → **OQ-019**. Ini **bukan** penetapan bounded context; itu D4.

### 8.5.2 Komite bukan wrapper approval tanpa data sendiri

Modul Komite berbagi 55–69% identitasnya dengan modul Claim pasangannya, tetapi **59 dari 133
identitas Komite (44,4%) tidak muncul di modul Claim mana pun** — termasuk class khusus seperti
`ASM-FW-GCNMFW-WORK-KOMITELIFE`. `[dugaan]` Komite adalah tahap approval dengan model data dan
layar sendiri. **Belum terverifikasi.**

### 8.5.3 Kode pembayaran klaim menggerbangi logika, artinya tidak diketahui

`PaymentType` (nilai 0–7) dan `TransferType` (nilai 1–4) dipakai sebagai guard, dengan
**pengelompokan berbeda-beda** antar rule — mis. `{1,2,5}`, `{4,6}`, `{3}` dalam satu rule
(`Claim Prop/Activity/HitServiceToKasir_Act.xml`). Korpus tidak memuat tabel kode → **OQ-020**.

### 8.5.4 Identitas orang ter-hardcode sebagai guard otorisasi — temuan lintas batch

Dicari saat batch 3, ditemukan di **45 file di 6 modul** termasuk modul batch 1 dan 2:
perbandingan `OperatorID.pyUserIdentifier` dengan **nama orang sebagai literal**, ditambah guard
berbasis `pyPosition`. Sapuan korpus penuh di batch 4 menaikkan angkanya menjadi **66 file di 9
modul** (lihat §8.6.4). Memblokir seluruh desain RBAC → **OQ-021**, berkaitan dengan **OQ-007**.
Pemeriksaan ini **tidak dilakukan** pada batch 1 dan 2 sebelumnya.

### 8.5.5 Metode normalisasi dikoreksi — 44% konflik OQ-011 ternyata palsu

Rule `LINKSERVICE / LINKSERVICE` terdeteksi "berbeda" padahal isinya sama; diff menunjukkan
perbedaannya murni **metadata provenance ekspor** (`pxHostId`, `pxMoveImport*`, `pxCreateOperator`,
dll). Daftar tag volatil diperluas dari **7 menjadi 18**.

| Batch | Angka lama | **Terkoreksi** |
| ---: | ---: | ---: |
| 1 | 130 | **57** |
| 2 | 807 | **461** |
| 3 | 39 | **25** |
| **Total** | **976** | **543** |

Angka yang berlaku adalah kolom terkoreksi. Daftar tag mungkin masih belum lengkap.

### 8.5.6 Temuan negatif dan anomali lain

- **Nol objek database link** `@ASMD.SINARMAS.CO.ID` di domain klaim — mempersempit **OQ-017**
  ke domain facultative.
- **Tidak ada skema Oracle baru**; 8 stored procedure baru, seluruhnya seputar akseptasi klaim
  → **OQ-002**.
- **4 dari 24 `ConnectREST` hardcode URL** (batch 2: 1 dari 20), keempatnya service klaim
  Non-Proportional. `appdev.nusantarare.com` muncul 3× → **OQ-018**.
- `excludeXML/GetBase64Attachment.xml` terverifikasi **rule Activity nyata**, bukan placeholder
  → mempertajam **OQ-003**.
- **4 rule `When` di `@BASECLASS` punya 4 isi berbeda** masing-masing (`ISCLM`, `ISCLMP`,
  `ISCLMNP`, `ISPEGASYARIAH`) — konflik guard paling tajam sejauh ini.
- `GCNMFW` terbukti framework **seluruh domain klaim**, bukan khusus facultative → **OQ-008**.

## 8.6 Temuan utama D1 batch 4 (ringkas) — batch penutup

Rincian di `inventory/_summary.md` §21–§22. Cakupan: 682 file di Endorsement Life, Master Contract
Retro Life, Master Product Name Life, PremiumList Life, Treaty Contract Out → **623 identitas unik**.

### 8.6.1 Tidak ada "ruleset Life bersama"

Overlap tertinggi antar modul batch 4 hanya **23 identitas** (Endorsement Life ↔ PremiumList Life).
Keempat modul Life berdiri sendiri-sendiri — berbeda tajam dari facultative (99,0%) dan treaty (320).

### 8.6.2 OQ-004 terjawab, OQ-005 terkoreksi

`PremiumList Life/InputPolicyHolder.xml` = rule **`Flow`**. Modul tanpa `Flow` turun dari 6 menjadi
**5**; jumlah rule `Flow` korpus naik dari 23 menjadi **24**.

### 8.6.3 `Treaty Contract Out` — nama folder tidak didukung isinya

`[terverifikasi]` Modul ini **tidak merujuk satu pun objek treaty outward**, sementara **tujuh modul
lain** merujuk `M_TREATY_OUT` / `TREATY_OUT2` / `FACOUTPRODUCTION`. Objek yang disentuhnya bertema
**kontrak treaty** (`TREATYCONTRACT`, `M_PROPORTIONALARRG`, `MTREATYSECURITY`, `TREATYREINSURER`,
`M_TREATYYEAR`, `TREATYEXCHANGE`), dan 16 rule `CANCELACTIVITY*` menamai **klausul kontrak**
(bordereaux, cash loss limit, EPI, ex-gratia, profit commission, RI comm, …). Hanya 5 dari 303 rule
menyebut `TreatyOut`, seluruhnya soal lampiran. → **OQ-022**. Penetapan konteks tetap D4.

Catatan: **202 dari 303 rule (67%) di class `@BASECLASS`** — proporsi tertinggi di korpus → OQ-009.

### 8.6.4 Sapuan lintas-batch yang dituntaskan di batch 4

- **OQ-017 tuntas 100% korpus:** database link `@ASMD.SINARMAS.CO.ID` **hanya di 3 modul
  facultative inward**, nol di 17 modul lain.
- **OQ-021 disapu 20 modul:** **66 file di 9 modul** memakai identitas orang sebagai guard
  (sebelumnya tercatat 45/6 — pola pencarian batch 3 melewatkan kutip tunggal). Guard ini berada di
  `<pyStepsPreCondParamsWhen>` (193×), `<pyCondition>` (58×), `<pyReadOnlyCondition>` (18×),
  `<pyContainerVisibleWhen>` (8×), dan menggerbangi antara lain
  `ASM-FW-GISFW-WORK / COUNTGROSSPREMI_ACT` (perhitungan premi) dan
  `ASM-FW-GCNMFW-WORK-KOMITE / KOMITEPOST_REJECT` (keputusan komite klaim). Terpisah: **35 file di
  14 modul** memakai `pyPosition` dengan nilai yang **menyerupai nama peran** (`ReasLifeAdmin`,
  `ReasLifeSPV`, `ReasLifeMedicalAdvisor`, `IT Developer`, `Admin`, `SPV A/B`) — titik awal paling
  konkret untuk RBAC.
- **OQ-002 dikoreksi turun:** detektor procedure sebelumnya ikut menangkap
  `INSERT INTO skema.tabel (…)` dan referensi page Pega `{Page.property(...)}`. Setelah diperbaiki
  dan seluruh 1.146 rule ber-SQL diekstraksi ulang: **67 procedure kustom tanpa body** di 6 skema.
- **OQ-020 diperluas:** tiga enumerasi baru yang menggerbangi logika — `EdmType` (1, 3),
  `ProRateType` (1, 2, 3), dan `.Type` PremiumList (`QR`, `QP`, `TP`, `TR`, di dalam
  `<pyStepsPreCondParamsWhen>`). Arti seluruhnya belum terverifikasi.

## 9. Urutan kerja & gate

| Step | Keluaran | Status |
| --- | --- | --- |
| D0 | kerangka discovery + fakta korpus terukur | **selesai** |
| D1 | `inventory/<modul>.md` (20) + `inventory/_summary.md` + `inventory/_oq011-konflik-isi.md` | **SELESAI** — 20/20 modul, 9.369/9.369 file (100%) |
| D2 | `flows/<konteks>.md` per konteks | **SELESAI** — Tahap 1–5, **20 konteks / 428 rule**: T1 treaty/life offer (4), T2 Komite (4), T3 Claim (4), T4 Facultative Inward (3), **T5 Treaty & Master tanpa `Flow` (5)**. Empat sintesis domain + 2 dokumen metode (`flows/_METHOD.md`, `flows/_METHOD-noflow.md`). Penutup: **`D2-CLOSING-REPORT.md`** |
| D3 | `modules/<modul>.md` (20) + **`glossary.md`** | **SELESAI** — 20/20 catatan modul (7 bagian: peran, proses, entitas & tabel, integrasi, ketergantungan antar-modul, batas pengetahuan, OQ). `glossary.md` **169 entri** berkolom Bukti (52 istilah + 51 kode/enumerasi + 25 workbasket + 41 singkatan) — **kandidat seed `CONTEXT.md`**. `glossary-draft.md` dipertahankan sebagai jejak |
| D4 | `context-map.md` + `understanding-report.md` | **SELESAI** — **9 bounded context** (5 `full`, 4 `partial`) + **Identity & Access = `absent`**; 7 kebocoran batas berbukti; narasi per konteks & antar-konteks; peta integrasi (51 `ConnectREST`, 15 layanan); batas pengetahuan total; 7 risiko migrasi. Penutup: **`D3-D4-CLOSING-REPORT.md`** |

**Penutup D1:** `D1-CLOSING-REPORT.md` — ringkasan lintas domain, daftar OQ terbuka per pemilik
peran, risiko halusinasi teridentifikasi, dan peta transisi ke D2.

**Penutup D2:** `D2-CLOSING-REPORT.md` — 20 konteks + titik masuk, peta integrasi eksternal,
rekap batas pengetahuan (batas korpus vs batas cakupan telusur), 57 OQ per pemilik peran dengan
tanda pemblokir D3/D4/FASE B, dan kesiapan D3–D4.

**Penutup D3 + D4:** `D3-D4-CLOSING-REPORT.md` — ringkas FASE A, isi GATE (4 dokumen yang
di-review + 7 keputusan manusia + 7 risiko), 57 OQ per pemilik peran dengan penanda **⛔ pemblokir
FASE B** (38 dari 57), usulan **urutan FASE B berbasis bukti**, dan **cara menjalankan skill FASE B
secara manual**.

**GATE — BERLAKU SEKARANG.** Manusia me-review `understanding-report.md`, `context-map.md`,
`glossary.md`, dan `open-questions.md`. **FASE B (ADR/BRD/spec/tiket) tidak boleh dimulai sebelum
gate ini lewat.**

`[terverifikasi]` Skill FASE B (`/mattpocock-skills:grill-with-docs`, `:to-spec`, `:to-tickets`)
ber-frontmatter `disable-model-invocation: true` — **hanya dapat dipicu manusia**, tidak dapat
dipanggil agent. Command persis dan prasyaratnya ada di `D3-D4-CLOSING-REPORT.md` §6.

**Urutan FASE B yang diusulkan** (berdasarkan jumlah OQ pemblokir + cakupan bukti):
**1. Claim — Life** (5 pemblokir, cakupan `full`, 136 file) → **2. Life — Penawaran & Premium List**
(9, `full`) → **3. Treaty Arrangement** (5, `partial`). Yang **harus menunggu jawaban bisnis**:
**Facultative Inward** (20 pemblokir, 6.071 file) dan **Treaty Inward — Master & Akseptasi** (18).
