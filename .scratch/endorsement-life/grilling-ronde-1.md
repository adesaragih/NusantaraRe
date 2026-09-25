# Grilling — Endorsement Life — Ronde 1

Tanggal: 2026-09-15
Konteks: `endorsement-life` — **konteks/menu terpisah** `[keputusan work owner]`
Modul: **Endorsement Life** (75 berkas), class work `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE`
Skill: `/mattpocock-skills:grilling` + `domain-modeling`
Sumber: korpus `D:\XML\RNM_BRD\` (READ-ONLY), `CONTEXT.md`, `docs/adr/`,
`.scratch/premiumlist-life/spec.md` (mesin bersama)

> **Konvensi.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[data DBA]`; `[dugaan]` = belum terbukti; `[terbuka]` = OQ.

---

## Bagian A — Alur inti yang belum pernah digrill: "pilih polis NB → muat data lama → endorse"

### A1. Gerbangnya: `SetErrorBatalEndorsement_Act` — **lima pemeriksaan sebelum endorsement boleh dibuat**

`[terverifikasi]` `Endorsement Life/Activity/SetErrorBatalEndorsement_Act.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETERRORBATALENDORSEMENT_ACT` / `RULE-OBJ-ACTIVITY`,
141.275 byte).

⚠️ **Namanya menyesatkan.** "SetErrorBatal" terbaca seperti pembatalan endorsement; isinya adalah
**gerbang validasi pembuatan** endorsement. Sub-langkah di bawah step 3 "Set Messages Error":

| Sub-step | Langkah | Precondition / rule | Arti |
| --- | --- | --- | --- |
| 3.1 | `Property-Set-Messages` | `param.Nopolis==""` | **nomor polis kosong** |
| 3.2 | `RDB-List` → `GetPL_NumberLife` (baris 895) | — | cari nomor polis di `JSON_POLIS` |
| 3.3 | `Property-Set-Messages` | `OutData.pxResults(1).CARI1==""` | **polis tidak ada di `JSON_POLIS`** |
| 3.4 | `Property-Set` | — | |
| 3.5 | `Call …pxRetrieveReportData` | `Param.pyReportName = "FilterProteksiEDMLife"` (baris 1211-1212) | cari EDM berjalan |
| 3.6 | `Property-Set-Messages` | `@LengthOfPageList(ListEdm.pxResults)>0` | **sudah ada EDM yang belum resolve** |
| 3.7 | `RDB-List` → `GetEdmTypeLife` (baris 1766) | — | baca `EdmType` polis |
| 3.8 | `Property-Set-Messages` | `@contains(OutData1.pxResults(1).CARI1,"3")` | **polis sudah pernah EDM batal** |
| 3.9 | `RDB-List` → `SearcStatusBayarArasaps_SQL` (baris 2090) | — | **"CEK PEMBAYARAN"** |
| 3.10 | `Property-Set-Messages` | `@LengthOfPageList(ListPembayaran.pxResults)>0 && TempWork.EdmType=="3"` | **sudah ada pembayaran** |
| 3.11 | `Property-Set` | — | |
| ~~4~~ | ~~`Obj-Save`~~ | — | **REMARK** (`<pyStepsBlockName>//`, baris 2562) |

`[terverifikasi]` **Satu endorsement terbuka per polis** — gerbang 3.6 menolak bila masih ada EDM
yang belum resolve. Report Definition-nya `FilterProteksiEDMLife`
(`Endorsement Life/ReportDefinition/FilterProteksiEDMLife.xml`).

`[terverifikasi]` **`EdmType == "3"` bermakna BATAL.** Dua bukti bebas dalam satu berkas:
gerbang 3.8 berdeskripsi *"Jika policy no sudah pernah edm batal"* menguji `@contains(…,"3")`, dan
gerbang 3.10 menyatukan cek pembayaran dengan `TempWork.EdmType=="3"`. **Nilai `EdmType` yang lain
tidak terbaca di korpus** → lihat **Q3**.

⚠️ `[terverifikasi]` **Pembacaan lintas skema.** `SearcStatusBayarArasaps_SQL`
(`ASM-FW-GISFW-INT-POLICYJSON` / `ASM!SEARCSTATUSBAYARARASAPS_SQL` / `RULE-CONNECT-SQL`):

```sql
select * from ARASAPAS.DETAIL_INVOICE where inv_inv_no = {InputData.CARI18} and IVD_JR_ID = '5'
```

Pega membaca **langsung ke skema `ARASAPAS`**, bukan lewat layanan. Arti `IVD_JR_ID = '5'` tidak
terbaca di korpus → **Q4**.

### A2. Mesinnya: `MappingEDMLife` — **14 langkah, nol yang ter-remark**

`[terverifikasi]` `Endorsement Life/Activity/MappingEDMLife.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` / `RULE-OBJ-ACTIVITY`, 194.433 byte).

| Step | Langkah | Catatan |
| ---: | --- | --- |
| 1 | *(gerbang)* | "Jika ada edm blm resolve exit act" — `@LengthOfPageList(ListEdm.pxResults)>0` (baris 456) |
| 2 | `Property-Set` | |
| **3** | **`Call svcAddWorkObject`** | **membuat case EDM baru** |
| 4 | `Property-Set` | |
| 5 | `Obj-Open-By-Handle` | |
| **6** | `RDB-List` → **`GetProdkeNopolis`** (baris 1172) | "Ambil idpega prodke terakhir dari tabel json_polis" |
| **7** | `Property-Set` | **"Set Old ID Pega dan Old Policy No"** |
| **8** | `Obj-Open-By-Handle` | **"Open IDPEGA NB Life"** ← **memuat polis new business** |
| **9** | `Property-Set` | **"Mapping detail dari Life"** |
| **10** | `Property-Set` | **"Copy page dari Life ke EDM"** |
| 11 | `Property-Set` | |
| 11.1 | `Property-Set` | **Set `"Old"` untuk detail lama** |
| 11.2 | `Property-Remove` | Remove `EDMStatus == "Delete"` |
| 12 | `Property-Set` | |
| 13 | `Obj-Save` | |
| **14** | **`Commit`** | **AKTIF** |

`[terverifikasi]` **Nol `<pyStepsBlockName>` di berkas ini.**

**Inilah alur inti yang tidak ada di PremiumList Life:** step 6→8 mengambil `IDPEGA` polis NB dari
`JSON_POLIS` berdasarkan nomor polis, membukanya, lalu step 9–10 menyalin detail dan halamannya ke
case EDM yang baru dibuat di step 3.

### A3. `EDMStatus` — mesin status **per baris detail**

`[terverifikasi]` Tiga nilai terbaca di korpus:

| Nilai | Terbaca di | Peran |
| --- | --- | --- |
| `"Old"` | `MappingEDMLife` step 11.1 "Set \"Old\" untuk detail lama" | baris warisan polis NB |
| `"New"` | `Endorsement Life/Activity/SaveCSVEDMLife.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE`) — "Remove EdmStatus \"New\"", "Set \"New\" untuk detail baru" | baris yang ditambahkan endorsement |
| `"Delete"` | `MappingEDMLife` step 11.2, precondition `.EDMStatus=="Delete"` | baris yang dihapus endorsement |

→ **Q7**: apakah ketiganya lengkap, dan apa yang terjadi pada baris `"Delete"` saat disimpan.

### A4. `GetOldDetail_EDM` — **4 dari 5 langkah ter-REMARK**

`[terverifikasi]` `Endorsement Life/Activity/GetOldDetail_EDM.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM` / `RULE-OBJ-ACTIVITY`, 53.517 byte):

| Step | Langkah | `<pyStepsBlockName>` |
| ---: | --- | --- |
| ~~1~~ | `Property-Set` | **`//`** (baris 244) |
| **2** | `Obj-Open-By-Handle` "Get data Life Old" (page `WorkLife`) | kosong — **AKTIF** |
| ~~3~~ | `Property-Set` | **`//`** (539) |
| ~~4~~ | `Property-Set` | **`//`** (841) |
| ~~5~~ | `Obj-Save` | **`//`** (975) |

⚠️ **OQ-066 berlaku** — penanda `blockname` terbukti tidak dapat dipercaya sendirian di modul Life.
→ **Q6**.

### A5. `CreateCaseEMDL` — **enam langkah hidup, sisanya mati**

`[terverifikasi]` `Endorsement Life/Activity/CreateCaseEMDL.xml` (`DATA-PORTAL` / `CREATECASEEMDL` /
`RULE-OBJ-ACTIVITY`, 119.680 byte) — ⚠️ **class `DATA-PORTAL`**, bukan class work endorsement:

| Step | Langkah | Status |
| ---: | --- | --- |
| 1 | `Property-Set` | aktif |
| 2 | `Page-New` | aktif |
| 3 | `Call svcAddWorkObject` | aktif |
| **4** | **`Commit`** | aktif |
| 5 | `Property-Set` | aktif |
| 6 | `Obj-Open-By-Handle` | aktif |
| ~~7~~ | `RDB-List` "-- Set data dari json_polis" | **REMARK** (1084) |
| ~~8~~ | `Java` | **REMARK** (1273) |
| ~~9~~ | `Page-Copy` "-- Copy List old" | **REMARK** (1382) |
| ~~10~~ | `Page-Copy` "-- Copy Ke new" | **REMARK** (1532) |
| ~~11~~ | `Property-Set` | **REMARK** (1671) |

Penyalinan data lama di sini **mati**; yang hidup adalah jalur `MappingEDMLife` (§A2). Deskripsi
langkah mati diawali `--`, pola penandaan manual.

`[terverifikasi]` `Endorsement Life/Activity/CancelCreateCaseEDML.xml` (`DATA-PORTAL` /
`CANCELCREATECASEEDML`, 20.134 byte) hanya **satu langkah** `Property-Set`. → **Q9**.

---

## Bagian B — Penomoran & PRODKE

`[terverifikasi]` **Dua rule berbeda membaca `PRODKE`, dengan urutan yang berbeda:**

| Rule | Identitas | SQL |
| --- | --- | --- |
| `GetProdkeNopolis` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPRODKENOPOLIS` / `RULE-CONNECT-SQL` | `SELECT IDPEGA … WHERE NOPOLIS={TempWork.PolicyNo} AND PRODKE IS NOT NULL ORDER BY **PRODKE DESC**` |
| `GetProdKeOldData_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!GETPRODKEOLDDATA_SQL` / `RULE-CONNECT-SQL` | `select PRODKE … where NOPOLIS={TempPolis.CARI4} order by **TGL_INPUT desc**` |

⚠️ **Urutannya berbeda** — `PRODKE DESC` versus `TGL_INPUT desc`. Bila urutan produksi dan urutan
input pernah tidak sejalan, keduanya memberi baris berbeda. → **Q8**.

`[terverifikasi]` `GetEdmTypeLife` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETEDMTYPELIFE`)
membaca **atribut di dalam CLOB**: `SELECT A.DATA_JSON.EdmType … FROM POOLDATA.JSON_POLIS A`.
`EdmType` **bukan kolom** — ia field di dalam `DATA_JSON`.

`[terverifikasi]` Kunci masuk di seluruh rule endorsement adalah **`TempWork.PolicyNo`** (nomor
polis), **bukan** `PL_NUMBER`. → **Q1**.

**Fondasi yang sudah mapan (tidak digrill ulang):** `PL_NUMBER_EDM` = `<nomor polis>/<PRODKE+1>`
lewat `Generate_NoEndorsmentLife` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` /
`ASM!GENERATE_NOENDORSMENTLIFE`), dijaga idempoten oleh precondition `PL_NUMBER_EDM==""` pada
`GenerateNoEDM_Life` step 5 dan 6.

---

## Bagian C — Alur, simpan, dan efek keluar

`[terverifikasi]` Flow `Endorsement Life/Flow/InputEDMLife.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-FLOW`, 68.374 byte):
task `InputEDMLife_Detail`, `InputEDMLife_Summary` ("Input EDM Summary"), `Accept`, utility
`InsertJsonPolisLife_Act`, keputusan `IsLifeAccepted`; routing **`WorkList`**; status akhir
**`Resolved-Completed`** dan **`Resolved-Rejected`**.

⚠️ `[terverifikasi]` **`IsLifeAccepted` versi Endorsement hanya mengeluarkan `Confirm` dan
`Decline`** (`Endorsement Life/DecisionTable/IsLifeAccepted.xml`,
`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `ISLIFEACCEPTED`, 15.096 byte) — **tanpa `Reject`**,
berbeda dari versi PremiumList Life (`ASM-FW-GISFW-WORK-LIFE` / `ISLIFEACCEPTED`) yang mengeluarkan
`Confirm` / `Decline` / `Reject`. → **Q5**.

`[terverifikasi]` Rantai simpan `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT`, 444.634 byte, **16 langkah**):
`GetProdKeOldData_SQL` (911) → `Generate_NoEndorsmentLife` (1304) → `InsertJsonPolisEDM` (2543) →
`SaveLifeinProduction_SQL` (2720) → `SaveMasterLPDet` (5154) → `InsertPLSummary` (6226) →
`GetNopolisByIDPega` (6444). Dua langkah REMARK: step **13** "Cek sudah masuk atau blm datanya"
(6399) dan step **15** `SendEmailNotification` (6752). Step **16**
`serviceInsertArasapasLife_act` aktif, dijaga `IsPEGAPROD` (7230).

⚠️ **Endorsement hari ini berjalan tanpa alarm** — deteksi dan email keduanya mati. → **Q10**.

`[terverifikasi]` **Titik potong transaksi di jalur endorsement** (procedure/SQL yang commit
sendiri): `InsertJsonPolisEDM` (`COMMIT;` baris 107), `SaveLifeinProduction_SQL` (164),
`SaveMasterLPDet` (252), `InsertPLSummary` (125), ditambah `Commit` **aktif** di `MappingEDMLife`
step 14 dan `CreateCaseEMDL` step 4.

⚠️ `[terverifikasi]` **`InsertJsonPolisEDM` adalah `INSERT` langsung, bukan upsert.**
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM` / `RULE-CONNECT-SQL`):
`INSERT INTO POOLDATA.JSON_POLIS (IDPEGA, DATA_JSON, TGL_INPUT, NOPOLIS, NOENDORS, PRODKE, TGL_PROD,
USERNAME) VALUES (…)`. Berbeda dari jalur NB yang memakai procedure upsert `INSERTJSONPOLISLIFE`
berkunci `IDPEGA`. **Pengulangan di jalur endorsement TIDAK aman secara gratis.** → **Q11**.

---

## Bagian D — Layar

`[terverifikasi]` Empat varian layar polis lama — **satu per jenis transaksi**:

| Section | Identitas |
| --- | --- |
| `ViewOldPolicy_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM` |
| `ViewOldPolicy_EDM_QP` | `…` / `VIEWOLDPOLICY_EDM_QP` |
| `ViewOldPolicy_EDM_TP` | `…` / `VIEWOLDPOLICY_EDM_TP` |
| `ViewOldPolicy_EDM_TR` | `…` / `VIEWOLDPOLICY_EDM_TR` |

Harness dengan nama sama ada untuk keempatnya. `QP`/`TP`/`TR` adalah nilai `Type` yang sudah ada di
glossary (`TP` = Payable, `TR` = Receivable). ⚠️ **`QR` tidak punya varian** — hanya tiga dari empat.
→ **Q12**.

`[terverifikasi]` Report Definition `BrowsePremiumList_RD` di modul ini beridentitas
`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `BROWSEPREMIUMLIST_RD` — **class integrasi summary**, sama
seperti salinan di PremiumList Life.

---

## Bagian E — Catatan metodologi

`[dugaan]` **Kode transisi `WhenTrue`.** Di `MappingEDMLife` dan `SetErrorBatalEndorsement_Act`,
gerbang yang berdeskripsi *"exit act"* memakai `<pyStepsPreCondParamsWhenTrue>6</…>` dan
`<pyStepsPreCondParamsWhenFalse>2</…>`. Ini **menguatkan** dugaan bahwa **`6` = Exit Activity** dan
`2` = lanjut. Belum cukup untuk menyatakan terbukti — pemetaan kode transisi Pega tidak ada di
korpus. Terkait catatan sebelumnya di Komite Claim Life yang membaca nilai `2` pada langkah
berdeskripsi EXIT. **Jangan dipakai sebagai fakta** sampai dikonfirmasi.

---

# Frontier Ronde 1 — **12 pertanyaan, 12 TERJAWAB** (work owner, 2026-09-15)

Frontier Ronde 1 **kosong.** Tiga fakta bisnis yang tidak saya tebak (Q3, Q4a, Q12) dijawab work
owner — **tidak ada OQ baru yang perlu dibuka untuk ketiganya.**

---

## ✅ V1 (Q1) — Titik masuk: **nomor polis** adalah kuncinya

`[keputusan work owner]` Pengguna memilih polis yang akan di-endorse dengan **nomor polis**
(`NOPOLIS` / `TempWork.PolicyNo`). `BrowsePremiumList_RD`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `BROWSEPREMIUMLIST_RD` / `RULE-OBJ-REPORT-DEFINITION`)
hanya **alat bantu pencarian**, bukan sumber kunci.

Sejalan dengan korpus `[terverifikasi]`: `GetPL_NumberLife`, `GetEdmTypeLife`, dan
`GetProdkeNopolis` seluruhnya `WHERE NOPOLIS = {TempWork.PolicyNo}`.

---

## ✅ V2 (Q2) — Lima gerbang kelayakan berlaku; **rule-nya salah nama**

`[keputusan work owner + terverifikasi]` `SetErrorBatalEndorsement_Act`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETERRORBATALENDORSEMENT_ACT` / `RULE-OBJ-ACTIVITY`) adalah
**gerbang kelayakan endorsement**, bukan pembatalan. **Namanya diganti di sistem baru.**

Kelima penolakan **berlaku seluruhnya**:

| # | Penolakan | Precondition korpus |
| --- | --- | --- |
| 1 | nomor polis kosong | `param.Nopolis==""` |
| 2 | polis tidak ada di `JSON_POLIS` | `OutData.pxResults(1).CARI1==""` (lewat `GetPL_NumberLife`) |
| 3 | **sudah ada EDM belum resolve** | `@LengthOfPageList(ListEdm.pxResults)>0` (lewat RD `FilterProteksiEDMLife`) |
| 4 | polis sudah pernah EDM batal | `@contains(OutData1.pxResults(1).CARI1,"3")` (lewat `GetEdmTypeLife`) |
| 5 | sudah ada pembayaran | `@LengthOfPageList(ListPembayaran.pxResults)>0 && TempWork.EdmType=="3"` |

⚠️ Gerbang **(3) — satu EDM terbuka per polis — adalah aturan integritas inti** konteks ini.

---

## ✅ V3 (Q3) — `EdmType`: **1 = Perubahan Data, 3 = Batal**

`[keputusan work owner + terverifikasi]`

| Nilai | Arti | Sumber |
| --- | --- | --- |
| `1` | **Perubahan Data** | `[keputusan work owner]` |
| `3` | **Batal** | `[terverifikasi]` — dua bukti bebas di `SetErrorBatalEndorsement_Act` (deskripsi gerbang 3.8 *"sudah pernah edm batal"* menguji `@contains(…,"3")`; gerbang 3.10 menyatukan cek pembayaran dengan `EdmType=="3"`) |
| `2`, `4` | **tidak dipakai** | `[keputusan work owner]` |

**Enum efektif = {1, 3}.** Endorsement punya **dua maksud**: mengubah data polis, atau
membatalkannya. `[terverifikasi]` `EdmType` dibaca dari **dalam CLOB** —
`SELECT A.DATA_JSON.EdmType … FROM POOLDATA.JSON_POLIS A` (`GetEdmTypeLife`,
`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETEDMTYPELIFE`) — ia **bukan kolom**.

---

> ✅ **DIKUATKAN RONDE 2 (2026-09-15).** Arti nilai **`3` kini `[terverifikasi]` dari kode**, bukan
> hanya dari deskripsi langkah: `SetPremi_EDM` step 2.1 berpresyarat
> `.EdmBatal=="True" || pyWorkPage.EdmType==3` di bawah deskripsi **"Set 0 jika EDM Batal"**
> (baris 1273). Arti nilai `1` tetap `[keputusan work owner]`.


## ✅ V4 (Q4) — `IVD_JR_ID = '5'` = **pembayaran/pelunasan**; baca Arasapas **tetap langsung, tapi terkurung**

`[keputusan work owner]` `IVD_JR_ID = '5'` berarti **pembayaran/pelunasan**.

`[keputusan desain]` Sistem baru **tetap membaca `ARASAPAS.DETAIL_INVOICE` langsung** — tidak diubah
menjadi pemanggilan layanan, karena itu perubahan kontrak dengan pihak lain, di luar cakupan
migrasi. Syaratnya: pembacaan itu **dikurung dalam SATU repository yang ditandai sebagai batas
lintas sistem**, tidak tersebar.

`[terverifikasi]` `SearcStatusBayarArasaps_SQL` (`ASM-FW-GISFW-INT-POLICYJSON` /
`ASM!SEARCSTATUSBAYARARASAPS_SQL` / `RULE-CONNECT-SQL`):
`select * from ARASAPAS.DETAIL_INVOICE where inv_inv_no = {InputData.CARI18} and IVD_JR_ID = '5'`.

---

## ✅ V5 (Q5) — Endorsement **tidak punya `Reject`**

`[keputusan work owner + terverifikasi]` `IsLifeAccepted` versi Endorsement
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `ISLIFEACCEPTED` / `RULE-OBJ-DECISIONTABLE`, 15.096 byte)
hanya mengeluarkan **`Confirm`** dan **`Decline`** — **tanpa `Reject`**, berbeda dari versi
PremiumList Life (`ASM-FW-GISFW-WORK-LIFE` / `ISLIFEACCEPTED`).

**Endorsement tidak punya jalur "kembali ke input".** Ini perbedaan nyata antar konteks, bukan
kelalaian.

---

> ⚠️ **DISEMPURNAKAN RONDE 2 (2026-09-15).** `[terverifikasi]` `GetOldDetail_EDM` **terpasang hidup
> pada keempat tombol** "lihat polis lama" di `ShowLifePremiumSummary_EDM` (`Run Activity` lalu
> `showHarness` popup). Rule-nya boleh tidak dimigrasikan, tetapi **perilakunya masih terpakai** —
> ia yang memuat data polis lama ke page `WorkLife` sebelum popup tampil. Lihat
> `grilling-ronde-2.md` §B1 dan **Q15**.


## ✅ V6 (Q6) — `GetOldDetail_EDM` **dibuang**

`[keputusan work owner]` `Endorsement Life/Activity/GetOldDetail_EDM.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM` / `RULE-OBJ-ACTIVITY`) **tidak dipakai —
jangan dimigrasikan.** Pemuatan data polis lama dikerjakan `MappingEDMLife` step 6–10.

⚠️ **Penerapan OQ-066:** keputusan "mati" ini datang dari **work owner**, bukan dari pembacaan
`<pyStepsBlockName>`. Empat dari lima langkahnya memang ber-`//`, tetapi penanda itu sudah terbukti
tidak dapat dipercaya sendirian di modul Life — jadi ia **bukan** dasar keputusan, hanya pendukung.

---

## ✅ V7 (Q7) — `EDMStatus`: tiga nilai, dan `"Delete"` adalah **soft-delete**

`[keputusan work owner + terverifikasi]` Tiga nilai **lengkap**:

| Nilai | Arti | Bukti |
| --- | --- | --- |
| `"Old"` | baris warisan polis NB | `MappingEDMLife` step 11.1 "Set \"Old\" untuk detail lama" |
| `"New"` | baris ditambahkan endorsement | `SaveCSVEDMLife` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE`) — "Set \"New\" untuk detail baru" |
| `"Delete"` | baris dihapus endorsement | `MappingEDMLife` step 11.2, precondition `.EDMStatus=="Delete"` |

⚠️ **`"Delete"` adalah penanda / soft-delete — BUKAN penghapusan fisik.** Baris polis NB di
`M_LIFE_PREMIUM_DETAIL` **tetap utuh**. Konsisten dengan prinsip **"endorsement = entri baru, bukan
mutasi in-place"**.

---

> ⚠️ **DIKOREKSI RONDE 2 (2026-09-15).** `[terverifikasi]` Korpus menunjukkan **nilai keempat:
> `"Batal"`**, ditulis `SetPremi_EDM` step 2.3 (baris 1506). Pernyataan "tiga nilai lengkap" di atas
> **tidak akurat**. Lihat `grilling-ronde-2.md` §C2 dan **V13** — beda `"Delete"` versus `"Batal"`
> adalah **cakupan minus**: `Delete` = minus **selektif per peserta**, `Batal` = minus **menyeluruh**.
> ⚠️ Pernyataan "soft-delete/penanda saja" juga **dikoreksi** — `Delete` **JUGA** menghasilkan nilai
> negatif. Yang tetap benar: baris polis NB **tidak dihapus fisik**.


## ✅ V8 (Q8) — Pembacaan `PRODKE` **disatukan ke `ORDER BY PRODKE DESC`**

`[keputusan desain]` **Penyimpangan sadar.** Di Pega ada **dua pembaca dengan urutan berbeda**:

| Rule | Identitas | `ORDER BY` |
| --- | --- | --- |
| `GetProdkeNopolis` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPRODKENOPOLIS` / `RULE-CONNECT-SQL` | **`PRODKE DESC`** |
| `GetProdKeOldData_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!GETPRODKEOLDDATA_SQL` / `RULE-CONNECT-SQL` | `TGL_INPUT desc` |

Bila urutan produksi dan urutan input pernah tidak sejalan, keduanya memberi `PRODKE` berbeda — dan
**nomor EDM ikut salah**. Sistem baru memakai **`ORDER BY PRODKE DESC` di kedua tempat**.

---

## ✅ V9 (Q9) — Case dibuat **setelah gerbang lolos**; **membatalkan endorsement = `Decline`**

`[keputusan work owner]`

1. Case EDM dibuat **hanya setelah kelima gerbang kelayakan lolos** (V2). Commit dini di
   `CreateCaseEMDL` step 4 dan `MappingEDMLife` step 14 **wajar dan dipertahankan** — gerbang "satu
   EDM terbuka per polis" justru **membutuhkan** case itu terlihat oleh pengguna lain.
2. **Tidak ada tombol batal terpisah.** Membatalkan endorsement = **`Decline`**, yang menutup case
   (`Resolved-Rejected`); polis lalu **bebas di-endorse ulang**.
3. `CancelCreateCaseEDML` (`DATA-PORTAL` / `CANCELCREATECASEEDML` / `RULE-OBJ-ACTIVITY`, satu langkah
   `Property-Set`) **bukan** mekanisme pembatalan case.

---

## ✅ V10 (Q10) — Alarm endorsement **dihidupkan**

`[keputusan work owner]` Di `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT`), step **13** (deteksi keadaan
separuh, "Cek sudah masuk atau blm datanya", baris 6399) dan step **15** `SendEmailNotification`
(baris 6752) yang di Pega **REMARK** → **dihidupkan di sistem baru**, seragam dengan jalur new
business.

⚠️ **Penyimpangan sadar** — Pega hari ini menjalankan endorsement **tanpa alarm**.

---

## ✅ V11 (Q11) — **Penjaga anti-dobel** pada `JSON_POLIS` jalur endorsement

`[keputusan desain]` **Penyimpangan sadar.** `[terverifikasi]` `InsertJsonPolisEDM`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM` / `RULE-CONNECT-SQL`) adalah
**`INSERT` polos**, bukan upsert — dan ia **commit sendiri** (`COMMIT;` baris 107). Pengulangan
menggandakan rekam polis. Jalur NB tidak punya masalah ini karena memakai procedure upsert
`INSERTJSONPOLISLIFE` berkunci `IDPEGA`.

Sistem baru: **penjaga idempotensi** — periksa **`(NOPOLIS, PRODKE)`** sebelum insert; **jangan**
tulis ulang bila sudah ada.

---

## ✅ V12 (Q12) — **QR tetap didukung**, lewat section induk

`[keputusan work owner + terverifikasi]` Ketiadaan `ViewOldPolicy_EDM_QR` **bukan** berarti `QR`
tidak didukung. **`QR` ditangani section induk `ViewOldPolicy_EDM`**
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM`); `_QP`, `_TP`, `_TR` adalah **varian
khusus**.

**Keempat jenis (`QP`, `QR`, `TP`, `TR`) didukung endorsement.**

---

> ✅ **DIKUATKAN RONDE 2 (2026-09-15) — kini `[terverifikasi]`.** `ShowLifePremiumSummary_EDM`
> memasangkan harness `ViewOldPolicy_EDM` (induk) dengan `<pyCondition>.Type=='QR'` (baris 65394),
> `_QP` dengan `.Type=='QP'` (65952), `_TR` dengan `.Type=='TR'` (66510), `_TP` dengan `.Type=='TP'`
> (67068). Lihat `grilling-ronde-2.md` §B.


## Ringkasan penyimpangan sadar Endorsement Life

| # | Penyimpangan | Alasan | Sumber |
| --- | --- | --- | --- |
| 1 | `ORDER BY PRODKE DESC` di **kedua** pembaca | Pega punya dua urutan berbeda → nomor EDM bisa salah | `[keputusan desain]` |
| 2 | **Penjaga anti-dobel** `(NOPOLIS, PRODKE)` sebelum insert `JSON_POLIS` | Pega `INSERT` polos + commit sendiri → pengulangan menggandakan | `[keputusan desain]` |
| 3 | **Alarm dihidupkan** (deteksi separuh + email) | Pega menjalankan endorsement tanpa alarm; NB punya | `[keputusan work owner]` |
| 4 | Baca Arasapas **terkurung di satu repository** bertanda batas lintas sistem | mencegah kueri lintas skema tersebar | `[keputusan desain]` |
| 5 | Rule gerbang **diganti nama** (bukan "SetErrorBatal…") | nama Pega menyesatkan: ia gerbang kelayakan, bukan pembatalan | `[keputusan work owner]` |

## Kode mati yang tidak dimigrasikan `[keputusan work owner]`

| Bagian | Identitas | Penanda korpus |
| --- | --- | --- |
| `GetOldDetail_EDM` **seluruhnya** | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM` / `RULE-OBJ-ACTIVITY` | 4 dari 5 langkah `//` (244, 539, 841, 975) |
| `CreateCaseEMDL` step **7–11** | `DATA-PORTAL` / `CREATECASEEMDL` / `RULE-OBJ-ACTIVITY` | `//` di 1084, 1273, 1382, 1532, 1671; deskripsi diawali `--` |

⚠️ **OQ-066 tetap berlaku**: keputusan mati di atas datang dari **work owner**. `<pyStepsBlockName>`
dipakai sebagai **pendukung**, tidak pernah sebagai dasar tunggal.
