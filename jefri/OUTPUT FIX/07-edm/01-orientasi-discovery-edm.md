# Discovery Endorsement — Orientasi dan rencana bertahap

> **Sumber:** `D:\migrasi\RNM\Endorsment Fac In\` — **2.061 berkas `.xml`**, READ-ONLY.
> Korpus Treaty tersendiri (`RNM_BRD\`) **tidak dibaca** (K-005). Fokus aktif ditetapkan **K-042**.
> Label mengikuti `CLAUDE.md` §3. Setiap angka disertai perintah audit yang menghasilkannya.
>
> ⚠️ **Ini orientasi, bukan discovery mendalam.** Isinya peta skala dan rencana; temuan perilaku
> menyusul per putaran.

---

## 1. Temuan yang menentukan pendekatan

### ⛔ EDM bukan RNW. Nol berkas identik.

`[terverifikasi]`

| Ukuran | NB | RNW | **EDM** |
| --- | ---: | ---: | ---: |
| Total `.xml` | 2.083 | 1.927 | **2.061** |
| Bernama sama dengan NB | — | 1.907 | **1.707** |
| — **identik byte-per-byte** | — | **1.907** | **0** |
| — berbeda isinya | — | **0** | **1.707** |
| Eksklusif (tidak ada di NB) | — | 20 | **354** |

**Konsekuensi yang mengikat cara kerja:** pada RNW, 1.907 berkas cukup dibuktikan identik lalu
diwarisi. Pada EDM **tidak ada satu pun** berkas yang boleh diperlakukan begitu. Setiap berkas
bernama sama **wajib dibandingkan isinya**, dan nol identik adalah buktinya.

Discovery EDM karena itu **lebih dekat ke discovery NB dari awal** daripada ke discovery RNW yang
ringan.

### ⚠️ Dua angka yang sah, tergantung kepekaan huruf

`[terverifikasi]`

| Perbandingan | Bernama sama | EDM-only |
| --- | ---: | ---: |
| **Peka huruf** (dipakai sebagai angka resmi) | **1.707** | **354** |
| Tidak peka huruf | 1.719 | 342 |

Selisih **12 berkas** yang namanya hanya berbeda kapitalisasi — 11 rule `When` + 1 Section
(`periode`↔`Periode`). Daftar lengkapnya di **K-042**.

⚠️ **Jebakan metodologis baru:** `Test-Path` di Windows **tidak peka huruf**, sehingga tidak dapat
membedakan varian kapitalisasi. Perbandingan nama wajib memakai `HashSet[string]`, bukan `Test-Path`
atau hashtable PowerShell (keduanya tidak peka huruf).

📌 Ini jebakan **kelima** dari keluarga "nama bukan bukti" di proyek ini — setelah beda tipe rule
(`isApproved`), nama identik tipe berbeda (`GetInsuredID`), posisi kata (`TSIOld` vs `.OldTSI`), dan
ada-tidaknya sufiks (`InputDtlObject`). Kali ini **kapitalisasi**.

---

## 2. Inventaris per tipe rule

```powershell
$nb="D:\migrasi\RNM\NB FacIn"; $rw="D:\migrasi\RNM\RNW Fac In"; $ed="D:\migrasi\RNM\Endorsment Fac In"
foreach ($t in (Get-ChildItem $ed -Directory).Name) {
  "{0,-22} {1,6} {2,6} {3,6}" -f $t,
    @(Get-ChildItem "$nb\$t" -Filter *.xml -File -EA SilentlyContinue).Count,
    @(Get-ChildItem "$rw\$t" -Filter *.xml -File -EA SilentlyContinue).Count,
    @(Get-ChildItem "$ed\$t" -Filter *.xml -File -EA SilentlyContinue).Count }
```

| Tipe rule | NB | RNW | **EDM** | EDM − NB |
| --- | ---: | ---: | ---: | ---: |
| Activity | 609 | 556 | **582** | −27 |
| Section | 432 | 397 | **434** | **+2** |
| FlowAction | 250 | 239 | **265** | **+15** |
| When | 210 | 189 | **202** | −8 |
| RDBList | 217 | 200 | **188** | −29 |
| DataTransform | 148 | 138 | **157** | **+9** |
| ReportDefinition | 123 | 118 | **138** | **+15** |
| DataPage | 30 | 30 | **41** | **+11** |
| Harness | 41 | 41 | **37** | −4 |
| DecisionTable | 12 | 10 | **10** | −2 |
| ConnectREST | 3 | 3 | **3** | 0 |
| **Flow** | 6 | 4 | **2** | **−4** |
| DecisionTree | 1 | 1 | **1** | 0 |
| SystemSettings | 1 | 1 | **1** | 0 |
| **TOTAL** | 2.083 | 1.927 | **2.061** | −22 |

📌 **Dua angka yang menarik perhatian:**

- **`Flow` hanya 2** — paling sedikit di antara ketiga siklus (NB 6, RNW 4). `[dugaan]` alur
  endorsement mungkin terpusat pada sedikit flow besar; perlu dikonfirmasi.
- **`DataPage` +11, `ReportDefinition` +15, `FlowAction` +15** — EDM **menambah** lapisan tampilan dan
  pengambilan data dibanding NB. `[dugaan]` konsisten dengan kebutuhan endorsement membaca data polis
  yang sudah jadi.

---

## 3. Tiga ratus lima puluh empat berkas EDM-only

### Per tipe rule `[terverifikasi]`

| Tipe | Jumlah | | Tipe | Jumlah |
| --- | ---: | --- | --- | ---: |
| Activity | **94** | | RDBList | 22 |
| Section | **81** | | DataPage | 14 |
| FlowAction | **54** | | Harness | 4 |
| DataTransform | 29 | | Flow | 1 |
| ReportDefinition | 28 | | | |
| When | 27 | | **TOTAL** | **354** |

### Tema nama — dan mengapa hasilnya sendiri adalah temuan

`[terverifikasi]` Pengelompokan kasar menurut kata kunci pada nama:

| Tema | Jumlah |
| --- | ---: |
| Life | 15 |
| Old / before-image | 12 |
| Endorsement / EDM | 9 |
| Batal / cancel | 2 |
| Treaty · Retro | 1 · 1 |
| **lain-lain** | **314** |

⚠️ **314 dari 354 tidak masuk tema apa pun.** Itu **bukan kegagalan pengelompokan** — itu temuan:
delta EDM **tidak terkonsentrasi pada tema yang dapat dikenali dari nama**. Konsekuensinya, pemetaan
EDM-only **tidak dapat dipercepat lewat nama**; ia harus dibaca isinya.

📌 Hanya **9 dari 354** bernama mengandung "endors"/"edm". Nama **bukan** penanda milik-siklus.

### Enam berkas yang ada di EDM **dan** RNW tetapi bukan NB

> 🔁 **DIKOREKSI di E-4** — lihat `05-e4-edm-only.md` §6. Klaim "keenamnya berbeda" diperoleh dari
> **hash mentah**, yang selalu berbeda karena stempel waktu **dan urutan tag ekspor**. Dengan metode
> 23-tag: **4 dari 6 IDENTIK**; dua sisanya (`ChooseInsured`, `ChooseInsuredDtl`) berbeda **hanya di
> penomoran indeks turunan** `pzIndexes`. Premis "K-038 tidak otomatis berlaku untuk EDM karena
> isinya beda" karena itu **tidak bertahan** — keputusannya tetap milik work owner (**E-Q25**).

`[terverifikasi]` ~~Keenamnya **berbeda isinya** antara EDM dan RNW~~ (dicabut — lihat kotak di atas):

| Berkas | Catatan |
| --- | --- |
| `Activity\GetBusinessGroup_Act` · `RDBList\CariBusinessGID` | Kelompok bisnis — di RNW jadi tiket **R06** |
| `Harness\ChooseInsured` · `Section\ChooseInsuredDtl` · `ReportDefinition\BrowseAccountInsuredEDM` · `Activity\SetDataInsuredEDM_Act` | Fitur pilih-tertanggung — **dibuang untuk RNW** (K-038) |

⚠️ **Empat di antaranya adalah fitur yang sudah diputuskan dibuang untuk renewal.** Keberadaannya di
EDM dengan **isi berbeda** berarti keputusan K-038 **tidak otomatis berlaku untuk endorsement** — ia
perlu keputusan tersendiri di fase ini.

⚠️ Jangan memakai `Test-Path` untuk pemeriksaan semacam ini: pengukuran pertama memberi **18** karena
ikut menghitung 12 varian kapitalisasi. Angka yang benar **6**.

---

## 4. Yang diwarisi sebagai titik awal — wajib dikonfirmasi ulang

⛔ Nol berkas identik membuat pewarisan temuan **berisiko**. Keempat butir berikut adalah **titik awal
penyelidikan**, bukan fakta yang sudah berlaku untuk EDM.

| Warisan | Status |
| --- | --- |
| **39 rule `When` bercabang**, seluruhnya melibatkan EDM sebagai yang berbeda | dari arsip — **konfirmasi ulang** |
| **Pola P1**: EDM membaca agregat tersimpan `.OfferFacIn.QuotationData.*`; NB/RNW membaca halaman aktif `pyWorkPage.Quotation.*` | dari arsip — **konfirmasi ulang**; bila benar, ini pembeda struktural paling penting |
| Alur endorsement + **before/after image** sebagian terpetakan | dari arsip — **titik awal** |
| `EdmType` `1`/`2`/`4` terkunci; kode `3` usang tetapi cabangnya **diport apa adanya** | **K-029** — sudah keputusan, tidak ditanya ulang |
| `Activity\SetOldData` + properti sufiks `TSIOld` | **K-039** — ditunda dari RNW, **ditelusuri di fase ini** |

---

## 5. Rencana discovery bertahap

Urutannya dipilih agar setiap putaran **menghasilkan keputusan**, bukan sekadar data.

| Putaran | Area | Mengapa urutan ini |
| ---: | --- | --- |
| **E-1** | **Pola perbedaan pada 1.707 berkas bersama** — ukur *bagaimana* mereka berbeda, bukan hanya *bahwa* berbeda. Uji pola P1 (agregat tersimpan vs halaman aktif) sebagai hipotesis utama | Bila satu pola menjelaskan mayoritas perbedaan, seluruh 1.707 dapat ditangani sebagai **satu keputusan rancangan** alih-alih 1.707 pemeriksaan. Ini putaran dengan daya ungkit tertinggi |
| **E-2** | **Alur endorsement** — 2 Flow (paling sedikit dari ketiga siklus) + FlowAction. Titik masuk, transisi, gerbang `EdmType` | Alur menentukan bentuk seluruh spec. Hanya 2 flow, jadi murah dibaca tuntas |
| **E-3** | **Before/after image dan `SetOldData`** — mekanisme selisih, properti `*Old`, K-039 | Ini inti endorsement: nilai dasar akseptasi endorsement memakai **selisih**, bukan nilai penuh. Tanpa ini, tangga akseptasi endorsement tidak dapat dispec |
| **E-4** | **354 EDM-only** — Activity (94) dan Section (81) lebih dulu | Terbesar, tetapi tidak dapat dipercepat lewat nama (314/354 tanpa tema). Ditunda sampai E-1 memberi kerangka pembacaan |
| **E-5** | **Jalur produksi endorsement** — pasangan nilai-sesudah/selisih (`_MENJADI`/`_SELISIH`), penomoran versi polis (`PRODKE`) | Bergantung E-3; juga menunggu jalur produksi NB yang masih Out of Scope |
| **E-6** | **Celah dan penutup** — rujukan menggantung, kode usang, pertanyaan terbuka | Pola dari NB dan RNW: putaran penutup selalu memunculkan sisa |

### Yang **tidak** dikerjakan di fase ini

`Endorsment Fac In\` adalah fokus aktif; **NB dan RNW tetap final** dan tidak dibuka kembali tanpa
alasan baru. Korpus Treaty tersendiri tetap tidak dibaca.

---

## 6. Pertanyaan terbuka yang sudah terlihat di orientasi

| # | Pertanyaan | Dampak |
| ---: | --- | --- |
| E-Q1 | **Apakah resolusi rule Pega peka huruf?** Pertanyaan lama, dan di EDM ia berdampak pada **12 berkas** — apakah `IsCar` dan `IsCAR` satu rule atau dua | Menentukan angka resmi 1.707/354 vs 1.719/342 · pemilik: **IT** |
| E-Q2 | Mengapa EDM hanya punya **2 Flow** sementara NB punya 6? Apakah alur endorsement terpusat, atau sebagian flow tidak terekspor | Menentukan kelengkapan pemetaan alur (E-2) |
| E-Q3 | Apakah keputusan **K-038** (pilih-tertanggung dibuang) berlaku juga untuk endorsement? `[terverifikasi]` keempat berkasnya **ada di EDM dengan isi berbeda** | Menentukan lingkup delta EDM · pemilik: **work owner** |

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
