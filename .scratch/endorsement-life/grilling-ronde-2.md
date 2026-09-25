# Grilling — Endorsement Life — Ronde 2

Tanggal: 2026-09-15
Konteks: `endorsement-life` — konteks/menu terpisah `[keputusan work owner]`
Lanjutan dari: `.scratch/endorsement-life/grilling-ronde-1.md` (12 verdict, frontier Ronde 1 kosong)
Skill: `/mattpocock-skills:grilling` + `domain-modeling`

> **Konvensi.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[data DBA]`; `[dugaan]`; `[terbuka]` = OQ.

---

## Bagian A — Editabilitas per form: **bukan rule `When`, melainkan properti `.EditInput`**

`[terverifikasi]` **Modul Endorsement Life hanya punya dua rule `When`**, dan keduanya bukan soal
editabilitas:

| Rule | Identitas |
| --- | --- |
| `IsPEGAPROD` | `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` |
| `recordEvent` | `@BASECLASS` / `RECORDEVENT` / `RULE-OBJ-WHEN` |

Fondasi "field yang boleh diubah mengikuti kondisi `When` per form" **terwujud sebagai kondisi
sebaris di Section**, bukan sebagai rule `When` bernama. Mekanismenya `<pyDisabledWhen>`.

### A1. Tiga kondisi penggerak `[terverifikasi]`

`Endorsement Life/Section/InputEDMLife.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` /
`RULE-OBJ-HTML-SECTION`, 1.255.106 byte) — **tujuh** `<pyDisabledWhen>`:

| Baris | Kondisi | Kendali atas |
| ---: | --- | --- |
| 3166 | `.EditInput==1` | **`.PolicyNo`** (input box) |
| 5231 | `.EditInput==1` | tombol |
| 7537 | `.EditInput==1` | **`.EdmTypeBatal`** (dropdown) |
| 8965 | `.EditInput1=1` | input box |
| 10403 | `.EditInput1=1` | tombol |
| 15763 | `.EditInput==1` | **`.EdmBatal`** (checkbox) |
| 37200 | `.IsJsonPolis=1` | tombol |

`[terverifikasi]` `<pyReadOnlyCondition>` di berkas ini hanya berisi `1=1` (6×) dan `always` (3×) —
**konstanta**, bukan aturan bisnis. Jadi **satu-satunya mesin editabilitas adalah `pyDisabledWhen`**.

### A2. `.EditInput` adalah **kunci satu arah** `[terverifikasi]`

Sensus seluruh modul — hanya **tiga** tempat menulisnya, semuanya **ke nilai `1`**, tidak ada yang
mengembalikannya ke `0`:

| Penulis | Identitas | Baris | Nilai |
| --- | --- | ---: | --- |
| `MappingEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` / `RULE-OBJ-ACTIVITY` | 2956 | `.EditInput = 1` |
| `SetPremi_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY` | 1338 | `.EditInput = 1` |
| `SaveCSVEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE` / `RULE-OBJ-ACTIVITY` | 2899 | `.EditInput1 = 1` |

**Artinya:** begitu data polis lama dipetakan (`MappingEDMLife`), atau premi ditetapkan
(`SetPremi_EDM`), **nomor polis, jenis batal, dan centang batal terkunci permanen** pada case itu.
Bagian yang dikunci `.EditInput1` terkunci setelah CSV disimpan. → **Q16**.

⚠️ Dua properti berbeda (`.EditInput` dan `.EditInput1`) mengunci **dua kelompok field berbeda** —
bukan salah ketik.

---

## Bagian B — Empat layar polis lama: **variannya dipilih oleh `.Type`, dan induk = `QR`**

`[terverifikasi]` `Endorsement Life/Section/ShowLifePremiumSummary_EDM.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SHOWLIFEPREMIUMSUMMARY_EDM` / `RULE-OBJ-HTML-SECTION`,
2.194.694 byte) memasang **empat pasang tombol**, masing-masing diikuti `<pyCondition>` atas `.Type`:

| Harness (baris) | `<pyCondition>` (baris) | Jenis |
| --- | --- | --- |
| `ViewOldPolicy_EDM` (65043, 65282) | **`.Type=='QR'`** (65394) | **induk = QR** |
| `ViewOldPolicy_EDM_QP` (65575, 65816) | `.Type=='QP'` (65952) | varian |
| `ViewOldPolicy_EDM_TR` (66136, 66374) | `.Type=='TR'` (66510) | varian |
| `ViewOldPolicy_EDM_TP` (66694, 66932) | `.Type=='TP'` (67068) | varian |

✅ **Ini mengangkat V12 dari `[keputusan work owner]` menjadi `[terverifikasi]`.** Section induk
`ViewOldPolicy_EDM` memang melayani `QR`; keempat jenis didukung.

`[terverifikasi]` Section juga memuat empat kontainer berpenanda `<pyContainerVisibleWhen>`:
`.Type=='QR'` (4801), `.Type=='QP'` (16218), `.Type=='TP'` (32111), `.Type=='TR'` (47725), serta satu
`<pyCondition>.Type=='TR' || .Type=='TP'` (4354).

### ⚠️ B1. `GetOldDetail_EDM` **terpasang hidup pada keempat tombol**

`[terverifikasi]` Tiap tombol adalah **pasangan aksi**: `Run Activity` → **`GetOldDetail_EDM`**
(`pyActivityClass` = `ASM-FW-GISFW-Work-EndorsementLife`), lalu `showHarness` dengan
`pyTarget = popup`.

Ini **menyempurnakan V6**, tidak membatalkannya: `GetOldDetail_EDM` memang boleh tidak
dimigrasikan sebagai rule, **tetapi perilakunya masih terpakai** — satu-satunya langkahnya yang
hidup (step 2 `Obj-Open-By-Handle` "Get data Life Old" ke page `WorkLife`) adalah **yang memuat data
polis lama sebelum popup ditampilkan**. Membuangnya begitu saja **mengosongkan keempat popup**.
→ **Q15**.

---

## Bagian C — `SetPremi_EDM`: mesin **pembalikan nilai** untuk endorsement batal

`[terverifikasi]` `Endorsement Life/Activity/SetPremi_EDM.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` /
`SETPREMI_EDM` / `RULE-OBJ-ACTIVITY`, 345.689 byte, **9 langkah, nol `<pyStepsBlockName>`**):

| Step | Langkah | Precondition / catatan |
| --- | --- | --- |
| 1 | `Page-Remove` | |
| 2 | *(blok)* | |
| **2.1** | `Property-Set` **"Set 0 jika EDM Batal"** | **`.EdmBatal=="True" \|\| pyWorkPage.EdmType==3`** (baris 1273) |
| **2.2** | `Property-Set` "flag pengurangan" | `.EdmBatal=="True"` (1441) → `.EditInput = 1`, **`.EDMStatus = "Delete"`** |
| **2.3** | `Property-Set` "flag batal" | `pyWorkPage.EdmType==3` (1583) / `.EdmBatal=="True"` (1636) → **`.EDMStatus = "Batal"`** |
| 3 | *(blok)* "Insert to detail" | |
| 3.1 | `Property-Set` "insert nilai dari data-batch → int" | `TempError.CARIDESC==1` (3321) |
| 3.x | *(gerbang)* | **`@SizeOfPropertyList(TempWorkPage.ListLifePremiumDetailUpload)>50000`** (3373) |
| 4 | *(blok)* + `Java` + sub-blok | precondition `.EDMStatus==local.EDMStatus` (3954) |
| 5 | `Property-Set` "Set COB" | |
| 6 | `Property-Set` "Set property PremiumListSummary" | |
| 7 | *(blok)* "Insert to summary" → 7.1 "Set PremiumListSummary to Pega (Temporary)" | |
| 8 | `Property-Set` | |
| 9 | `Obj-Save` | |

### ⚠️ C1. Deskripsi berkata "Set 0"; kode mengalikan **`× -1`**

`[terverifikasi]` Step 2.1 **tidak** menulis nol. Ia menulis `<kolom> = <kolom> * -1` pada
**32 kolom uang**, antara lain:

```
SUM_INSURED, CEDING_RETENTION, SUM_REASURED, SHARE_NUSANTARA_RE, SHARE_NUSANTARA_RE_GROSS,
SUM_AT_RISK_GROSS, SUM_AT_RISK_RETRO, RETROCEDED_SHARE, SHARE_RETRO, RATE, FACTOR,
GROSS_PREMIUM, NET_PREMIUM, DEDUCTION, CLAIM_AMOUNT, RI_ADMIN_FEE, BROKERAGE_FEE,
…_REFUND, …_RETRO, …_REFUND_RETRO
```

**Endorsement batal menghasilkan jurnal balik (counter-entry) bernilai negatif**, bukan baris nol —
sehingga penjumlahan seluruh baris polis menjadi nol. Deskripsi langkah **berbohong** terhadap
kodenya. → **Q14**.

### ⚠️ C2. `EDMStatus` punya **nilai keempat: `"Batal"`**

`[terverifikasi]` Sensus penulisan `EDMStatus` di seluruh modul:

| Nilai | Penulis | Baris |
| --- | --- | ---: |
| `"Old"` | `MappingEDMLife` | 2608 |
| `"New"` | `SaveCSVEDMLife` (`…PremiumListDetail(<LAST>).EDMStatus`) | 2714 |
| `"Delete"` | `SetPremi_EDM` step 2.2 | 1384 |
| **`"Batal"`** | **`SetPremi_EDM` step 2.3** | **1506** |

**Verdict V7 Ronde 1 menyebut "tiga nilai lengkap" — korpus menunjukkan empat.** Saya tidak
menutupi selisih ini. → **Q13.**

`[terverifikasi]` `EDMStatus` juga mengalir ke rantai simpan: `InsertJsonPolisLife_Act` baris 4455
menyalin `.EDMStatus` → `TempValue.EDMStatus`.

`[terverifikasi]` ✅ **Bukti kode langsung untuk `EdmType == 3` = batal** — precondition step 2.1
`… || pyWorkPage.EdmType==3` di bawah deskripsi **"Set 0 jika EDM Batal"**. Ini mengangkat arti
nilai `3` pada V3 dari sandaran work owner menjadi **`[terverifikasi]`**. (Arti nilai `1` tetap
`[keputusan work owner]`.)

---

## Bagian D — Unggah CSV endorsement

### D1. `UploadCSVEDMLifePremium_Act` — **hanya tiga langkah pertama yang hidup**

`[terverifikasi]` `Endorsement Life/Activity/UploadCSVEDMLifePremium_Act.xml` (`@BASECLASS` /
`UPLOADCSVEDMLIFEPREMIUM_ACT` / `RULE-OBJ-ACTIVITY`, 119.366 byte) — **rule base-class**:

| Step | Langkah | Status |
| ---: | --- | --- |
| 1 | `Page-Remove` | aktif |
| 2 | `Page-New` "deklarasi pyWorkPage" | aktif |
| **3** | **`Call pxUploadCSVResults`** "import data life premium" | aktif |
| ~~4~~ | `Page-New` | **REMARK** (645) |
| ~~5~~ | `Property-Set` (+ 5.2, 5.4 juga REMARK) | **REMARK** (791, 2120, 2400) |
| ~~6~~ | `Obj-Save` | **REMARK** (2525) |

Ia **hanya mengimpor CSV ke halaman**; penyimpanan dikerjakan `SaveCSVEDMLife`.

### D2. `SaveCSVEDMLife` — **5 langkah, nol remark**

`[terverifikasi]` `Endorsement Life/Activity/SaveCSVEDMLife.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE`
/ `SAVECSVEDMLIFE` / `RULE-OBJ-ACTIVITY`, 230.581 byte):

| Step | Langkah | Precondition / catatan |
| --- | --- | --- |
| 1 | `Page-Clear-Messages` | |
| 2 | `Property-Set` | |
| **2.1** | `Property-Remove` **"Remove EdmStatus \"New\""** | **`.EDMStatus=="New"`** (547) — **membuang baris `New` sebelumnya** sebelum impor ulang |
| **3** | **`call ASM-FW-GISFW-Work-LIFE.Calculate1_Act`** | ⚠️ **lintas class** — mesin bersama |
| 4 | `Property-Set` | |
| **4.1** | `Property-Set` "Set property PremiumListDetail" | `.PLAN = …PremiumListDetail(1).PLAN` (2466); `.POLICY_HOLDER = …PremiumListDetail(1).POLICY_HOLDER` (2495) |
| 4.2 | `Page-Set-Messages` "Set message error" | `local.errmsg==""` (2649) |
| **4.3** | `Property-Set` **"Set \"New\" untuk detail baru"** | **`pyWorkPage.EdmType==1`** (2791) |
| 5 | `Property-Set` | → `.EditInput1 = 1` (2899) |

Tiga hal yang keluar dari sini:

1. `[terverifikasi]` **Unggah ulang bersifat mengganti, bukan menumpuk** — step 2.1 membuang seluruh
   baris `EDMStatus=="New"` lebih dulu. Ini penjaga idempotensi yang sudah ada.
2. `[terverifikasi]` **Baris baru hanya ditambahkan bila `EdmType==1`** (Perubahan Data). Pada
   endorsement **batal** (`EdmType==3`), CSV **tidak** menambah peserta — konsisten dengan V3.
3. `[terverifikasi]` **Validasi konsistensi**: `PLAN` dan `POLICY_HOLDER` tiap baris diuji terhadap
   **baris pertama** premium list. Seluruh peserta dalam satu endorsement wajib satu plan dan satu
   pemegang polis.

### D3. Mesin bersama yang belum terdaftar: `Calculate1_Act`

`[terverifikasi]` `Calculate1_Act` (`ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` /
`RULE-OBJ-ACTIVITY`, 302.597 byte) ada di **kedua folder modul** — `PremiumList Life/Activity/` dan
`Endorsement Life/Activity/` — dengan **identitas sama** dan **diff ternormalisasi nol baris**.

⚠️ Class-nya **`ASM-FW-GISFW-WORK-LIFE`** (class PremiumList Life), dipanggil eksplisit lintas class
oleh `SaveCSVEDMLife` step 3. **Ini mesin bersama keenam**, dan belum tercatat di daftar
`CONTEXT.md`. → **Q18**.

---

# Frontier Ronde 2 — **6 pertanyaan, 6 TERJAWAB** (work owner, 2026-09-15)

Frontier Ronde 2 **kosong.** Dua fakta bisnis yang tidak saya tebak (Q13, Q14) dijawab work owner;
**OQ-070, OQ-071, dan OQ-072 ditutup.**

---

## ✅ V13 (Q13) — `EDMStatus`: empat nilai, dan **`Delete` juga menghasilkan minus**

`[keputusan work owner]`

| Nilai | Kapan | Akibat pada nilai uang |
| --- | --- | --- |
| `"Old"` | baris warisan polis new business | tidak diubah |
| `"New"` | peserta **ditambah** — hanya pada **Perubahan Data** (`EdmType=1`) | positif, baris baru |
| `"Delete"` | peserta **dihapus** dalam Perubahan Data | **diminuskan — per peserta (selektif)** |
| `"Batal"` | lewat **EDM Batal** (`EdmType=3`) | **otomatis SELURUH peserta batal — semua nilai diminuskan (menyeluruh)** |

**Bedanya adalah cakupan, bukan jenis:** `Delete` = **minus selektif per peserta**; `Batal` =
**minus menyeluruh**. Keduanya menempuh jalur yang sama — jurnal balik `× -1` di `SetPremi_EDM`.

⚠️ **Ini mengoreksi verdict V7 Ronde 1** yang menyebut `"Delete"` sebagai "penanda/soft-delete saja".
**`Delete` JUGA menghasilkan nilai negatif.** Anggapan lama sudah diberi catatan koreksi di
`grilling-ronde-1.md`.

---

## ✅ V14 (Q14) — Baris negatif masuk **tabel yang sama**; Claim Life **menyaringnya keluar**

`[keputusan work owner]`

1. **Penyimpanan:** baris negatif (batal maupun delete) ditulis ke tabel yang **SAMA** —
   `POOLDATA.M_LIFE_PREMIUM_DETAIL` dan `POOLDATA.M_LIFE_PREMIUM_SUMMARY`. Baris positif asli dan
   baris negatif **hidup berdampingan** → **net akunting** diperoleh dari penjumlahan, bukan dari
   penghapusan.
2. **Kontrak hilir — penajaman:** peserta yang sudah **EDM Batal** atau **soft-delete**
   **TIDAK BOLEH MUNCUL** di Claim Life. Kueri klaim
   `Claim Life/RDBList/GetPesertaClaim_sql1.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` /
   `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL`, berkunci `PL_NUMBER`) **WAJIB menyaring keluar**
   peserta batal/delete.

⚠️ **Ini penajaman kontrak lintas konteks**, bukan detail internal endorsement. Menyentuh
**tiket 08 PremiumList Life** (kontrak hilir) dan **spec Claim Life**. Akuntansi melihat seluruh
baris; klaim hanya melihat peserta yang masih hidup.

---

## ✅ V15 (Q15) — Popup polis lama **tetap ada**; rule mati dibuang, **perilaku dipertahankan**

`[keputusan work owner]` Popup "lihat polis lama" **tetap ada dan dipakai**, dipicu dari
`ShowLifePremiumSummary_EDM` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SHOWLIFEPREMIUMSUMMARY_EDM` /
`RULE-OBJ-HTML-SECTION`) menurut `.Type`:

| `.Type` | Layar |
| --- | --- |
| `QR` | **section induk `ViewOldPolicy_EDM`** |
| `QP` / `TP` / `TR` | varian `ViewOldPolicy_EDM_QP` / `_TP` / `_TR` |

**Perilaku "muat data polis lama untuk mengisi popup" WAJIB ADA** di sistem baru. Rule Pega
`GetOldDetail_EDM` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM` / `RULE-OBJ-ACTIVITY`,
4 dari 5 langkah mati) **tidak ditiru apa adanya**, tetapi **fungsinya dipertahankan**.

⚠️ **Mengoreksi V6 Ronde 1**: bukan "buang total", melainkan **"buang rule mati, pertahankan
perilaku popup"**.

---

## ✅ V16 (Q16) — Kunci field **permanen, dengan pesan penjelas**

`[keputusan work owner + desain]` Editabilitas lewat `<pyDisabledWhen>` atas `.EditInput` adalah
**kunci satu arah** — di-set `1`, tidak pernah kembali `0`.

Field yang terkunci **permanen** setelah endorsement dibuat: **nomor polis (`.PolicyNo`)**, **jenis
endorsement/batal (`.EdmTypeBatal`, `.EdmBatal`)**.

Sistem baru: **pertahankan kunci permanen** + **tampilkan pesan penjelas** — di Pega field hanya
mati tanpa alasan. Jalan keluar bila salah pilih polis tetap `Decline` lalu buat baru (V9).

---

## ✅ V17 (Q17) — **Tidak ada batas jumlah baris CSV**

`[keputusan work owner]` Batas **50.000 baris** milik Pega (`SetPremi_EDM` step 3, gerbang
`@SizeOfPropertyList(TempWorkPage.ListLifePremiumDetailUpload)>50000`, baris 3373) **DIBUANG**.
Sistem baru **tidak membatasi** jumlah baris unggahan CSV.

⚠️ **Penyimpangan sadar.**

**Catatan teknis:** unggahan sangat besar diproses **bertahap** (streaming/batch internal) —
**tanpa menolak karena jumlah baris**.

---

## ✅ V18 (Q18) — `Calculate1_Act` **BUKAN** mesin bersama untuk endorsement

`[keputusan work owner]` Meski berkas `Calculate1_Act` (`ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` /
`RULE-OBJ-ACTIVITY`, 302.597 byte) **identik** di kedua folder modul, **di jalur ENDORSEMENT ia
TIDAK dijalankan.** Ia milik **PremiumList Life (new business)** saja.

**Perhitungan endorsement dikerjakan `SetPremi_EDM`** (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` /
`SETPREMI_EDM`), bukan `Calculate1_Act`.

⚠️ **Pelajaran, seiring OQ-066:** **berkas identik ≠ dipakai.** Sebagaimana `<pyStepsBlockName>`
tidak boleh dipakai sendirian untuk menyimpulkan hidup/mati, **kesamaan berkas antar folder tidak
boleh dipakai sendirian untuk menyimpulkan "mesin bersama"**. Rujukan `<RequestType>` / `Call` pun
hanya menunjukkan **kemungkinan** pemanggilan — konfirmasi work owner tetap diperlukan.

**Akibat:** daftar mesin bersama di `CONTEXT.md` **dikoreksi** — `Calculate1_Act` dikeluarkan.

---

## Ringkasan penyimpangan sadar Endorsement Life — **daftar final**

| # | Penyimpangan | Alasan | Sumber |
| --- | --- | --- | --- |
| 1 | `ORDER BY PRODKE DESC` di **kedua** pembaca | Pega punya dua urutan berbeda → nomor EDM bisa salah | `[keputusan desain]` |
| 2 | **Penjaga anti-dobel** `(NOPOLIS, PRODKE)` sebelum insert `JSON_POLIS` | `InsertJsonPolisEDM` = `INSERT` polos + commit sendiri | `[keputusan desain]` |
| 3 | **Alarm dihidupkan** (deteksi separuh + email) | Pega menjalankan endorsement tanpa alarm; NB punya | `[keputusan work owner]` |
| 4 | Baca Arasapas **terkurung di satu repository** bertanda batas lintas sistem | mencegah kueri lintas skema tersebar | `[keputusan desain]` |
| 5 | Rule gerbang **diganti nama** (bukan "SetErrorBatal…") | nama Pega menyesatkan | `[keputusan work owner]` |
| 6 | **Tanpa batas jumlah baris CSV** (50.000 dibuang) | batas teknis, bukan aturan bisnis; unggahan besar diproses bertahap | `[keputusan work owner]` |
| 7 | **Peserta batal/delete disaring keluar** dari jalur baca klaim | akuntansi melihat seluruh baris; klaim hanya peserta hidup | `[keputusan work owner]` |

## Kode mati yang tidak dimigrasikan — **daftar final** `[keputusan work owner]`

| Bagian | Identitas | Penanda korpus | Catatan |
| --- | --- | --- | --- |
| `GetOldDetail_EDM` **sebagai rule** | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM` / `RULE-OBJ-ACTIVITY` | 4 dari 5 langkah `//` (244, 539, 841, 975) | ⚠️ **perilakunya tetap ada** — popup polis lama (V15) |
| `CreateCaseEMDL` step **7–11** | `DATA-PORTAL` / `CREATECASEEMDL` / `RULE-OBJ-ACTIVITY` | `//` di 1084, 1273, 1382, 1532, 1671 | deskripsi diawali `--` |
| `UploadCSVEDMLifePremium_Act` step **4–6** | `@BASECLASS` / `UPLOADCSVEDMLIFEPREMIUM_ACT` / `RULE-OBJ-ACTIVITY` | `//` di 645, 791, 2120, 2400, 2525 | hanya step 1–3 hidup |
| `InsertJsonPolisLife_Act` step **13** dan **15** | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` | `//` di 6399, 6752 | ⚠️ **justru DIHIDUPKAN** di sistem baru (V10) |
| **`Calculate1_Act` di jalur endorsement** | `ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` / `RULE-OBJ-ACTIVITY` | — berkas identik, **tetapi tidak dijalankan** | V18 — milik PremiumList Life |
| Batas **50.000 baris** | `SetPremi_EDM` step 3, baris 3373 | aktif di Pega | dibuang (V17) |

⚠️ **OQ-066 diterapkan penuh:** setiap penetapan "mati" bersandar pada **keputusan work owner**.
`<pyStepsBlockName>` dan kesamaan berkas dipakai sebagai **pendukung**, tidak pernah sebagai dasar
tunggal — V18 adalah contoh tepatnya.

---

## Status OQ setelah Ronde 2

| OQ | Status |
| --- | --- |
| **OQ-070** (`EDMStatus` empat nilai) | **TERTUTUP** — dijelaskan V13 |
| **OQ-071** (batal = jurnal balik negatif) | **TERTUTUP** — dijelaskan V13 + V14 |
| **OQ-072** (batas 50.000 baris) | **TERTUTUP** — dibuang, V17 |

**Endorsement Life MATANG — nol OQ pemblokir.** Sisa yang menyentuh konteks ini hanyalah
**OQ-001 (sisa)**: DDL fisik tabel premium, yang memblokir **tiket migrasi** saja.
