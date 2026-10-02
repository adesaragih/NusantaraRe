# Rekap Inventaris Rule — STEP D1

**Status: FINAL — STEP D1 SELESAI. Keempat batch tuntas, 20 dari 20 modul, 9.369 dari 9.369 file
(100%).**

Rekap lintas korpus ada di **§22** di bagian bawah. Bagian §2–§21 adalah catatan per batch, ditulis
saat batch itu dikerjakan, dan **dipertahankan apa adanya** termasuk koreksi angka yang terjadi di
tengah jalan — jejak koreksinya sendiri adalah informasi.

Laporan penutup D1 beserta peta transisi ke STEP D2: **`../D1-CLOSING-REPORT.md`**.

| Batch | Domain | Modul | File | Status |
| ---: | --- | --- | ---: | --- |
| 1 | Treaty inward | NB Treaty In, Treaty In, Treaty In Adjustment, EDM Treaty In | 1.149 | **selesai** |
| 2 | Facultative inward | Komite Claim FacIn, Claim Fac In, RNW Fac In, Endorsment Fac In, NB FacIn | 6.667 | **selesai** |
| 3 | Claim non-fac | Komite Claim Life, Komite Claim Non Prop, Komite Claim Prop, Claim Life, Claim Prop, Claim Non Prop | 871 | **selesai** |
| 4 | Life, master & outward | Endorsement Life, Master Contract Retro Life, Master Product Name Life, PremiumList Life, Treaty Contract Out | 682 | **selesai** |

Register konflik isi rule: **`_oq011-konflik-isi.md`** — **533 identitas** (angka korpus penuh)
(angka ini **turun dari 976** setelah metode normalisasi diperbaiki saat batch 3; lihat register
§"Revisi metode").

Dibangkitkan 2026-09-12 dari korpus READ-ONLY `D:\XML\RNM_BRD\`.

## 1. Metode & cara audit ulang

Setiap file XML dibaca satu kali, mengambil tag berikut:

| Tag | Dipakai untuk |
| --- | --- |
| `<pxInsName>` | **identitas rule** — `CLASS!NAMA`, atau `CLASS!PREFIX!NAMA` untuk RDBList |
| `<pzOriginalInstanceKey>` | field pertama = **tipe rule**; field 2–3 = **asal salinan** (bukan identitas — lihat `../README.md` §6.4) |
| `<pyActivityType>`, `<pyStepsActivityName>` | jenis Activity + jumlah step |
| `<pyStartActivity>` | activity awal Flow |
| `<pyInclude>` | embed Section (dependency nyata) |
| `<pyBrowseSQL>` | SQL RDBList → operasi, tabel (termasuk `TABEL@DBLINK`), pemanggilan stored procedure |
| `<pyServiceName>`, `<pyBaseURLSelectionType>`, `<pyBaseURLSetting>` | ConnectREST → target integrasi sebagai konfigurasi |
| `<pyPageName>`, `<pyScope>`, `<pyStructure>` | DataPage |
| `<pyPurpose>` | SystemSettings |
| `<pxObjClass>` | tipe rule, bila `<pzOriginalInstanceKey>` tidak ada |

Fakta per tipe **digerbang oleh folder** agar tag satu tipe tidak bocor ke tipe lain.

**Catatan revisi (2026-09-12):** keempat file modul batch 1 **dibangkitkan ulang** memakai extractor
dan generator yang sudah diperbaiki di batch 2, sehingga kesembilan file modul kini berformat sama
(7 bagian, termasuk §4 tabrakan nama dan §7 integrasi eksternal). Kolom identitas batch 1 diperiksa
tidak berubah — 1.149 file, 729 identitas unik, sama persis dengan pengukuran sebelumnya. Angka
batch 1 di §2–§6 di bawah tetap sah.

**Catatan revisi kedua (2026-09-12, saat batch 3) — dua perbaikan metode:**

1. **Daftar tag normalisasi diperluas dari 7 menjadi 18.** Sebelas tag provenance ekspor
   (`pxCreateOperator`, `pxCreateOpName`, `pxCreateSystemID`, `pxHostId`, `pxMoveImportDateTime`,
   `pxMoveImportOperId`, `pxMoveImportOperName`, `pzIndexCount`, `pyJavaGenerateTime`,
   `pyJavaClassName`, `pxServerReqURI`) sebelumnya ikut dihitung sebagai isi, sehingga rule yang
   sebenarnya sama terhitung "berbeda". **Akibatnya jumlah konflik OQ-011 turun dari 976 menjadi
   543** — 44,4% adalah positif palsu. Angka di §3.1 dan §9 sudah dikoreksi; rincian di
   `_oq011-konflik-isi.md` §"Revisi metode".
2. **`<pxObjClass>` dipakai sebagai fallback tipe rule** bila `<pzOriginalInstanceKey>` tidak ada,
   ditandai akhiran `(pxObjClass)`. Diperlukan karena 2 rule `Flow` batch 3 dan 5 rule batch 2 tidak
   memiliki tag itu.

Batch 1 dan 2 **tidak** diekstraksi ulang untuk perbaikan (2) — jumlah dan pengelompokan
identitasnya tidak terpengaruh, hanya tampilan tipe untuk 5 file batch 2 yang tetap kosong.

Perintah audit satu file (dijalankan dari `D:\XML\RNM_BRD\`):

```
grep -o "<pxInsName>[^<]*" "<modul>/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "<modul>/<tipe>/<file>.xml" | head -1
awk '/<pyBrowseSQL>/,/<\/pyBrowseSQL>/' "<modul>/RDBList/<file>.xml"
```

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis. Perilaku
ditelusuri di STEP D2.

## 2. Batch 1 — angka

| Ukuran | Nilai |
| --- | ---: |
| File XML | 1.149 |
| Ukuran | 226,8 MB |
| **Identitas rule unik** (tipe + class + nama) | **729** |
| File yang merupakan salinan identitas lintas modul | 420 |
| Class Pega unik | 65 |

Per modul:

| Modul | File | Class unik |
| --- | ---: | ---: |
| NB Treaty In | 278 | 38 |
| Treaty In | 329 | 38 |
| Treaty In Adjustment | 379 | 37 |
| EDM Treaty In | 163 | 25 |

(Jumlah per modul tidak dapat dijumlahkan menjadi 65 — banyak class dipakai lintas modul.)

Rincian per tipe rule per modul ada di masing-masing `<modul>.md` §1, dan matriks lengkap 20 modul
ada di `../README.md` §5.1.

## 3. Batch 1 — 1.149 file hanya berisi 729 rule

**420 file (36,6%) adalah identitas yang sama muncul di lebih dari satu modul.** Ini penting untuk
perencanaan migrasi: jumlah pekerjaan sebenarnya jauh lebih kecil dari jumlah file — **tetapi hanya
bila salinannya memang identik.** Ternyata tidak semuanya.

### 3.1 Dari 387 identitas lintas-modul, 57 isinya benar-benar berbeda

> **ANGKA DIKOREKSI 2026-09-12 (saat batch 3).** Bagian ini semula melaporkan **130**. Daftar tag
> normalisasi ternyata belum lengkap — 11 tag provenance ekspor terlewat, sehingga sebagian rule
> yang isinya sama terhitung "berbeda". Angka batch 1 yang berlaku adalah **57**, bukan 130.
> Metode dan daftar tag lengkap ada di `_oq011-konflik-isi.md` §"Revisi metode".

| Kategori | Jumlah identitas |
| ---: | --- |
| Muncul di >1 modul, isi **identik** setelah normalisasi | **330** |
| Muncul di >1 modul, isi **benar-benar berbeda** | **57** |

**Cara pengukuran (penting — perbandingan mentah menyesatkan).** Perbandingan `md5sum` langsung
menyatakan **seluruh 387 berbeda**, dan itu **salah**. Ekspor Pega menuliskan elemen dalam urutan
berbeda dan menyisipkan timestamp per ekspor, sehingga dua salinan rule yang sama tidak pernah
identik byte-per-byte.

Contoh: `Treaty In/Activity/TreatyInCopy.xml` dan `Treaty In Adjustment/Activity/TreatyInCopy.xml`
sama-sama 62.810 byte, berbeda pada 144 baris — seluruhnya berupa urutan elemen yang bergeser
ditambah `<pyRuleFormStatusTime>` (`20260903T033314.669 GMT` vs `20260904T025337.511 GMT`) dan
`<pyShowJavaWindowName>` (`GeneratedJava<timestamp>`). Secara isi, keduanya sama persis.

Normalisasi yang dipakai — buang tag volatil, urutkan baris, lalu hash:

```
VOL='<(pyRuleFormStatusTime|pyShowJavaWindowName|pxUpdateDateTime|pxCreateDateTime|pxCommitDateTime|pxSaveDateTime|pyRuleAvailableTime)>'
grep -vE "$VOL" "<file>.xml" | sort | md5sum
```

Setelah normalisasi **18 tag** (versi terkoreksi), 330 identitas terbukti identik dan **57 tetap
berbeda**. Blok kode di atas menampilkan daftar 7 tag versi pertama; daftar lengkap yang berlaku
ada di `_oq011-konflik-isi.md`.

**Risiko migrasi `[terverifikasi]`:** untuk 57 identitas tersebut, `class + nama + tipe` yang sama
merujuk isi yang **berlainan** tergantung modul mana yang dibuka. Membaca salah satu saja lalu
menganggapnya mewakili yang lain akan salah. Contoh:

- `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYINPRE_ACT` — `NB Treaty In/Activity/` vs `EDM Treaty In/Activity/`
- `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / COUNTTOTALPCTSPEAD` — `Treaty In/DataTransform/` vs `Treaty In Adjustment/DataTransform/`
- `DATA-PORTAL / TREATYINDIFFERENCEFACSHARE` — `Treaty In/Activity/` vs `Treaty In Adjustment/Activity/`
- `ASM-FW-GISFW-DATA-TREATYINSHAREREINS / SETTREATYINRETRO_ACT` — `Treaty In/Activity/` vs `Treaty In Adjustment/Activity/`

Mana yang berlaku di production **belum terverifikasi** → **OQ-011**.

Catatan koreksi: tiga contoh yang tercantum di versi sebelumnya (`TREATYINFIXSPL`, `AKSEPTASI_ACT`,
`ISMULTICOB`) **bukan konflik** setelah normalisasi diperbaiki — perbedaannya hanya metadata ekspor.
Contoh di atas sudah diverifikasi ulang terhadap register terkoreksi.

### 3.2 `Treaty In` dan `Treaty In Adjustment` berbagi 320 identitas

| Pasangan modul | Identitas yang dibagi |
| --- | ---: |
| Treaty In ↔ Treaty In Adjustment | **320** |
| EDM Treaty In ↔ NB Treaty In | 69 |
| NB Treaty In ↔ Treaty In | 22 |
| EDM Treaty In ↔ Treaty In Adjustment | 21 |
| NB Treaty In ↔ Treaty In Adjustment | 19 |
| EDM Treaty In ↔ Treaty In | 16 |

`Treaty In` memiliki 329 file dan `Treaty In Adjustment` 379; keduanya berbagi 320 identitas.
`[dugaan]` keduanya adalah satu aplikasi yang sama yang diekspor dua kali dengan cakupan sedikit
berbeda, bukan dua modul terpisah. Konfirmasi memerlukan D2 — dan perlu memperhitungkan bahwa
sebagian dari yang dibagi itu justru **berbeda isi** (§3.1).

## 4. Batch 1 — tabrakan nama

### 4.1 Nama file dipakai di lebih dari satu tipe rule — 63 nama

Pola yang mendominasi adalah konvensi Pega, bukan anomali:

- **Section + FlowAction** — satu FlowAction merujuk Section bernama sama
  (mis. `Share.xml`, `Layers.xml`, `Installments.xml`, `MaxRetention.xml`, `CoBList.xml`).
- **Section + Harness** — mis. `BusinessAndSOBList.xml`, `InputTreatyInOffer.xml`,
  `ShowAttachmentTreaty.xml`.

Yang **perlu perhatian** karena tipe rulenya berlainan sifat:

| Nama file | Tipe yang bertabrakan |
| --- | --- |
| `FetchTreatyContractDetail.xml` | Activity + RDBList |
| `GetAchievement.xml` | Activity + RDBList |
| `GetTreatyName.xml` | Activity + RDBList |
| `BrowseTreatyOutDetail.xml` | ReportDefinition + RDBList |

Contoh tabrakan yang bahkan class-nya sama: `ASM-FW-GISFW-DATA-INSTALLMENT / INSTALLMENTLIST`
dipakai oleh `FlowAction/InstallmentList.xml` **dan** `Section/InstallmentList.xml` di
`NB Treaty In`. Ini menegaskan aturan `../README.md` §3.2: **rule harus selalu ditulis
`class / nama / tipe`** — nama saja, bahkan nama + class, tidak cukup.

### 4.2 Rule `When` bernama sama di lebih dari satu class — 2 nama

| Nama When | Class |
| --- | --- |
| `TREATYMASTERINEDM` | `ASM-FW-GISFW-DATA-POLICYTREATYIN`, `DATA-PORTAL` |
| `ISUW` | `ASM-FW-GISFW-WORK`, `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` |

Jauh lebih sedikit dari dugaan awal. Arti keduanya **belum terverifikasi** — `ISUW` khususnya
(`UW` kepanjangan belum terverifikasi) akan menentukan guard otorisasi di D2.

## 5. Batch 1 — silsilah salinan (`Save As`)

410 dari 1.149 file memiliki `<pzOriginalInstanceKey>` yang berbeda dari identitasnya sendiri —
artinya rule tersebut dibuat dengan menyalin rule lain. Ini **sinyal kandidat konsolidasi**, bukan
bukti perilaku identik.

Keluarga klon paling menonjol, seluruhnya berasal dari
`ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT`:

- `NB Treaty In/Activity/CheckSpreadingProtectFire_ACT.xml`
- `NB Treaty In/Activity/CheckSpreadingProtectAnekaGolf_ACT.xml`
- `NB Treaty In/Activity/CheckSpreadingProtectMCargoMBU_ACT.xml`
- `NB Treaty In/Activity/CheckSpreadingProtectPATravel_ACT.xml`

`[dugaan]` empat rule ini adalah satu logika yang sama dengan parameter lini bisnis berbeda
(Fire, Aneka Golf, Marine Cargo/MBU, PA Travel). **Belum terverifikasi** — D2 harus membandingkan
isinya sebelum menyimpulkan bahwa keempatnya dapat dikolaps menjadi satu fungsi.

**Contoh mengapa nama tidak boleh dipercaya:** `NB Treaty In/RDBList/GetCurrency.xml`
(`ASM-FW-GISFW-INT-CURRENCY / GETCURRENCY`, sebuah SELECT) adalah hasil salinan dari
`ASM-FW-GISFW-INT-CURRENCY / ASM!UPDATEMASTERCURRENCY`. Sebuah rule baca yang lahir dari rule tulis.

Daftar lengkap keluarga klon per modul ada di `<modul>.md` §4.

## 6. Batch 1 — objek database

Diambil dari `<pyBrowseSQL>` pada 161 rule `RDBList`. **Nama apa adanya dari SQL.**

| Ukuran | Nilai |
| --- | ---: |
| RDBList berupa query biasa | 130 |
| RDBList berupa blok PL/SQL anonim (`BEGIN…END;`) | 31 |
| Tabel/view unik yang terlihat | 54 |
| Stored procedure/function yang dipanggil | 13 (+1 built-in Oracle) |

### 6.1 Kualifikasi skema tidak konsisten

Objek yang sama muncul dengan dan tanpa prefix skema — mis. `POOLDATA.M_TREATY_IN` dan
`M_TREATY_IN`, `POOLDATA.TREATY_IN_EDM` dan `TREATY_IN_EDM`. Jumlah "54 tabel/view unik" karena itu
**dihitung apa adanya dari teks SQL** dan kemungkinan melebih-hitung objek yang sebenarnya sama.
Angka objek fisik yang sesungguhnya **belum terukur** dan bergantung pada OQ-001.

Dua skema terlihat: **`POOLDATA`** (mayoritas objek bisnis) dan **`DATAPEGA`**
(mis. `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` — tabel milik Pega sendiri).

### 6.2 Data disimpan sebagai JSON di dalam kolom

`NB Treaty In/RDBList/BrowseTreatyIn.xml` (`ASM-FW-GISFW-INT-TREATY_IN / ASM!BROWSETREATYIN`):

```sql
select * from (
select JSONDATA as ClassofBusiness from pooldata.M_TREATY_IN where ID={TreatyIn.ID}
union all
select JSONDATA as ClassofBusiness from pooldata.M_TREATY_IN_edm where ID={TreatyIn.ID}
)
```

Kolom **`JSONDATA`** `[terverifikasi]`. Isi struktur JSON di dalamnya **tidak ada di korpus** →
**OQ-012**. Selama itu terbuka, model data domain tidak dapat diturunkan dari SQL saja.

Terlihat juga pola tabel berpasangan `X` dan `X_EDM` (`M_TREATY_IN` / `M_TREATY_IN_EDM`,
`M_TREATY_IN_DETAIL` / `M_TREATY_IN_DETAIL_EDM`, `TREATY_IN` / `TREATY_IN_EDM`). Arti `EDM`
**kepanjangan belum terverifikasi**.

### 6.3 Stored procedure yang dipanggil

**Body seluruh prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML (OQ-002).

| Procedure | Dipanggil oleh (jumlah rule) |
| --- | ---: |
| `POOLDATA.PEGA_TREATY_IN` | 4 |
| `POOLDATA.TREATYINOFFER` | 2 |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | 2 |
| `POOLDATA.PEGA_M_TREATY_IN_EDM` | 2 |
| `POOLDATA.PEGA_M_TREATY_IN_DETAIL_EDM` | 2 |
| `POOLDATA.PEGA_M_TREATY_IN_DETAIL` | 2 |
| `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | 2 |
| `POOLDATA.LOG_ACHIEVEMENT` | 2 |
| `POOLDATA.GET_TOKEN_STORAGE` | 2 |
| `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` | 1 |
| `POOLDATA.INSERTUPDATEACHIEVMENT` | 1 |
| `POOLDATA.INSERTJSONPOLISMONITORING` | 1 |
| `POOLDATA.HISTORYAKSEPTASIPRODUCTION` | 1 |
| `DBMS_LOB.CREATETEMPORARY` | 1 (built-in Oracle) |

Contoh `[terverifikasi]` — `NB Treaty In/RDBList/SavePolisTreatyIn_SQL.xml`
(`ASM-FW-GISFW-INT-POLISTREATYIN / ASM!SAVEPOLISTREATYIN_SQL`):

```sql
BEGIN
  POOLDATA.PEGA_JSON_POLIS_TREATYIN(
      {pyWorkPage.pzInsKey}, {pyWorkPage.PolicyTreatyIn.PolicyNo}, NULL, '0',
      {InputData.CARI21}, {OperatorID.pyUserIdentifier}, {InputData.CARI3},
      {OutputData.HASIL1 out});
  COMMIT;
END;
```

Dua hal yang terbaca langsung dan berdampak ke desain:

- **`COMMIT` dilakukan di dalam prosedur**, bukan dikelola aplikasi → batas transaksi ada di sisi
  database.
- **`{OperatorID.pyUserIdentifier}` dikirim sebagai parameter** → identitas pengguna ikut masuk ke
  prosedur, kemungkinan untuk audit. Bagaimana dipakai di dalam prosedur **belum terverifikasi**.

### 6.4 Prefix `ASM` / `RNM` pada kunci RDBList

Seluruh 161 RDBList batch 1 memakai kunci tiga bagian `CLASS!PREFIX!NAMA`, dengan prefix `ASM` (89)
atau `RNM` (72). Arti prefix **belum terverifikasi** → **OQ-008**. Bila keduanya menunjuk koneksi
database berbeda, ini berdampak langsung pada desain repository.

## 7. Batch 2 — angka

Domain facultative inward. Lima modul, dikerjakan berurutan kecil → besar.

| Ukuran | Nilai |
| --- | ---: |
| File XML | 6.667 |
| **Identitas rule unik** (tipe + class + nama) | **2.970** |
| File yang merupakan salinan identitas lintas modul | 3.697 |
| Tabel/view unik terlihat | 158 |
| Stored procedure/function unik | 50 (42 di antaranya baru) |
| Rule `ConnectREST` | 20 (9 nama service unik) |

Per modul — seluruhnya **cocok** dengan `../README.md` §5.1:

| Modul | File | Identitas unik | Class unik | Nama file dipakai >1 tipe | Rule di class bukan-aplikasi |
| --- | ---: | ---: | ---: | ---: | ---: |
| Komite Claim FacIn | 114 | 114 | 22 | 4 | 29 |
| Claim Fac In | 482 | 481 | 55 | 11 | 60 |
| RNW Fac In | 1.927 | 1.926 | 163 | 111 | 238 |
| Endorsment Fac In | 2.061 | 2.060 | 183 | 127 | 240 |
| NB FacIn | 2.083 | 2.081 | 173 | 122 | 260 |

Perintah audit jumlah file per modul (dari `D:\XML\RNM_BRD\`):

```
find "<modul>" -type f -name "*.xml" | wc -l
find "<modul>" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
```

`When` bernama sama di >1 class: **0** di kelima modul batch 2.

## 8. Batch 2 — NB FacIn, RNW Fac In dan Endorsment Fac In hampir satu ruleset

Ini fokus khusus batch 2. Overlap identitas diukur setelah normalisasi isi.

| Pasangan modul | Identitas yang dibagi |
| --- | ---: |
| **NB FacIn ↔ RNW Fac In** | **1.906** |
| **Endorsment Fac In ↔ NB FacIn** | **1.603** |
| **Endorsment Fac In ↔ RNW Fac In** | **1.573** |
| Claim Fac In ↔ RNW Fac In | 92 |
| Claim Fac In ↔ NB FacIn | 92 |
| Claim Fac In ↔ Komite Claim FacIn | 82 |
| Komite Claim FacIn ↔ RNW Fac In | 59 |
| Komite Claim FacIn ↔ NB FacIn | 59 |
| Claim Fac In ↔ Endorsment Fac In | 54 |
| Endorsment Fac In ↔ Komite Claim FacIn | 25 |

Sebaran identitas di antara ketiga modul besar:

| Sebaran | Identitas |
| --- | ---: |
| Ada di **ketiga** modul | **1.567 (61,4%)** |
| Ada di tepat 2 dari 3 | 381 |
| Hanya di 1 dari 3 | 604 |
| **Total identitas gabungan 3 modul** | **2.552** (dari 6.071 file) |

`[terverifikasi]` 1.906 dari 1.926 identitas `RNW Fac In` (99,0%) juga ada di `NB FacIn`.

`[dugaan]` ketiganya adalah **satu ruleset dengan diskriminator siklus** (baru / renewal /
endorsement), bukan tiga aplikasi terpisah — pola yang sama dengan `Treaty In` ↔
`Treaty In Adjustment` di batch 1 (§3.2).

**Ini BUKAN penetapan bounded context.** Penetapan konteks adalah STEP D4 dan memerlukan D2.
Yang dicatat di sini hanya angkanya. Ditambahkan sebagai **OQ-015**.

## 9. Batch 2 — 461 identitas berbeda isi (OQ-011)

> **ANGKA DIKOREKSI 2026-09-12 (saat batch 3).** Bagian ini semula melaporkan **807**. Lihat
> `_oq011-konflik-isi.md` §"Revisi metode". Angka batch 2 yang berlaku adalah **461**.

| Kategori | Jumlah identitas |
| ---: | --- |
| Muncul di >1 modul, isi **identik** setelah normalisasi 18 tag | **1.513** |
| Muncul di >1 modul, isi **benar-benar berbeda** | **461** |

Metode sama persis dengan §3.1 — hash ternormalisasi, bukan `md5sum` mentah. 5.669 file dihitung
hash ternormalisasinya.

**Tidak ada satu varian pun yang dipilih sebagai "benar".** Contoh yang **terverifikasi** masih
punya 3 versi isi berbeda:

- `ASM-FW-GISFW-INT-REINSURANCETYPE / BROWSEREINSURANCETYPE_RD` — 5 file di Komite Claim FacIn,
  Claim Fac In, RNW Fac In, Endorsment Fac In, NB FacIn
- `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTNUSARE_RD` — 4 file di Claim Fac In, RNW Fac In,
  Endorsment Fac In, NB FacIn

Contoh **2 versi** dari 3 varian: `ASM-FW-GISFW-INT-DOCUMENT_POLIS / GETALLDOCUMENT` —
`RNW Fac In/ReportDefinition/`, `Endorsment Fac In/ReportDefinition/`, `NB FacIn/ReportDefinition/`
(dua di antaranya isinya sama, satu berbeda).

Contoh `ASM-FW-GISFW-WORK / ISMULTICOB` yang tercantum di versi sebelumnya **bukan konflik** —
perbedaannya hanya metadata ekspor.

## 10. Batch 2 — objek database

### 10.1 Sepuluh skema Oracle, bukan satu

Batch 1 hanya memperlihatkan `POOLDATA` dan `DATAPEGA`. Batch 2 memperlihatkan **10 skema lokal**:

| Skema | Rujukan (tabel) | Rujukan (procedure) |
| --- | ---: | ---: |
| `POOLDATA` | 211 | 123 |
| `NEW_UNDERWRITING` | 11 | — |
| `GENERAL` | 8 | 3 |
| `DATAPEGA` | 7 | — |
| `REINSURANCE` | 5 | — |
| `FIRE` | 3 | 2 |
| `ARASAPAS` | 2 | — |
| `MBU` | — | 3 |
| `GL` | — | 2 |
| `NEW_GENERAL` | — | 1 |

Ditambah built-in Oracle: `DBMS_LOB` (31), `UTL_MATCH` (2).

Arti tiap skema **belum terverifikasi** — korpus tidak menjelaskannya. `[dugaan]` beberapa di
antaranya milik sistem lain (`GL` menyerupai general ledger, `ARASAPAS` menyerupai sistem
invoice/pembayaran — rujukannya `ARASAPAS.INVOICE` dan `ARASAPAS.DETAIL_INVOICE`). **Jangan
menebak** — ditambahkan sebagai **OQ-016**.

### 10.2 Database link ke sistem luar: `ASMD.SINARMAS.CO.ID`

`[terverifikasi]` Delapan objek diakses lewat **Oracle database link** bernama
`ASMD.SINARMAS.CO.ID` — bukan database Nusantara Re sendiri:

| Objek | Dirujuk oleh |
| --- | ---: |
| `FIRE.M_BI_INDEMNITY@ASMD.SINARMAS.CO.ID` | 6 rule |
| `M_TERORISME_RATE@ASMD.SINARMAS.CO.ID` | 3 rule |
| `M_RSMD_RATE@ASMD.SINARMAS.CO.ID` | 3 rule |
| `M_FLOOD_AREA@ASMD.SINARMAS.CO.ID` | 3 rule |
| `M_FLEXAS_RATE@ASMD.SINARMAS.CO.ID` | 3 rule |
| `M_EQS_RATE@ASMD.SINARMAS.CO.ID` | 3 rule |
| `FIRE.M_FLOOD_RATE@ASMD.SINARMAS.CO.ID` | 3 rule |
| `LST_KURS_STANDARD@ASMD.SINARMAS.CO.ID` | 2 rule |

Bukti: `NB FacIn/RDBList/SearchRatePolisEQS_SQL.xml`
(`ASM-FW-GISFW-INT-OFFERJSON / SEARCHRATEPOLISEQS_SQL`),
`NB FacIn/RDBList/GetKurs.xml` (`ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / GETKURS`),
`NB FacIn/RDBList/SearchRatePolisFlood_SQL.xml`.

Yang diambil lewat link ini adalah **tabel rate** (gempa/EQ, terorisme, banjir, FLEXAS, RSMD) dan
**kurs standar** — data master yang menentukan perhitungan premi. Artinya sebagian input rating
berada di luar basis data yang dimigrasikan. Ditambahkan sebagai **OQ-017**.

**Catatan metode:** parser awal membuang karakter `@` sehingga objek ini terbaca sebagai skema
palsu (`M_EQS_RATEASMD.SINARMAS.CO.ID`). Parser sudah diperbaiki dan seluruh RDBList batch 2
diekstraksi ulang. Bila menemukan nama objek berpola `...ASMD.SINARMAS.CO.ID` tanpa `@` di artefak
lain, itu sisa bug — bukan objek nyata.

### 10.3 Stored procedure: 50 unik, 42 baru

42 procedure yang tidak muncul di batch 1, di antaranya:

`POOLDATA.FACINLIFE`, `POOLDATA.FACINSPREADLIFE`, `POOLDATA.FACINOFFER`,
`POOLDATA.FACINPERFORMANCE`, `POOLDATA.FACINFORBACKUP`, `POOLDATA.ERRORFACINPROD`,
`POOLDATA.INSERTJSONPOLIS`, `POOLDATA.PEGA_M_ACCUMULATION_LIFE`, `POOLDATA.GENERATE_FACRETRO_NO`,
`POOLDATA.RDBMASTER*` (10 procedure master data: ACCUMULATION, ACCUMULATEDTYPE, BRANCH, CITY,
DISTRICT, NATION, PROVINCE, RW, SHIP, dan `RDBINSERTCLIENT`),
`POOLDATA.PEGA_PROGRESSCLAIM`, `POOLDATA.PEGA_SUBPROGRESSCLAIM`, `POOLDATA.CLAIMREJECTED`,
`POOLDATA.MONITORING_KLAIM_LOG`, `POOLDATA.MONITORING_PROD_LOG`, `POOLDATA.DIRECTTOKASIR_LOG`.

Dan yang berada di **skema selain POOLDATA** — ini memperluas OQ-002 secara material:

| Procedure | Catatan |
| --- | --- |
| `FIRE.PEGA_FIRE_SET_RATE` | `[dugaan]` penetapan rate lini Fire — logika rating di database |
| `FIRE.CEK_PRORATA_TANGGAL` | `[dugaan]` perhitungan pro-rata tanggal |
| `GENERAL.F_GET_NM_ASURADUR` | `[dugaan]` pencarian nama asuradur |
| `GL.F_GET_EMAIL` | `[dugaan]` pencarian email dari skema `GL` |
| `MBU.F_CEK_HURUF` | arti belum terverifikasi |
| `NEW_GENERAL.CEK_PLAT_NO` | `[dugaan]` validasi nomor plat kendaraan |

**Body seluruh procedure ini tidak ada di korpus** (OQ-002). Sebagian di antaranya menyimpan
**logika rating dan perhitungan premi** — bagian paling material dari sistem.

### 10.4 JSON opaque meluas

115 rule batch 2 menyentuh objek/kolom berlabel JSON (`JSON_POLIS`, `JSON_OFFER`, `JSONDATA`).
`JSON_POLIS` adalah objek paling sering dirujuk di seluruh batch 2 (**58 rule**). Menguatkan
OQ-012: model data domain tidak dapat diturunkan dari SQL saja.

## 11. Batch 2 — integrasi eksternal (`ConnectREST`)

20 rule `ConnectREST`, **9 nama service unik**. Daftar apa adanya dari
`<pyServiceName>`; **tidak ada nama sistem yang disimpulkan dari dokumen lain**.

| pyServiceName | Muncul di modul |
| --- | --- |
| `ServiceGoogle` | Komite Claim FacIn, Claim Fac In, RNW Fac In, Endorsment Fac In, NB FacIn |
| `getPremiumPaidOn` | Claim Fac In, RNW Fac In, Endorsment Fac In, NB FacIn |
| `convertJsonNusareToProduction` | RNW Fac In, Endorsment Fac In, NB FacIn |
| `KonversiKlaimNonLife` | Komite Claim FacIn, Claim Fac In |
| `SendAcceptationToKasir` | Komite Claim FacIn, Claim Fac In |
| `HitDLAClaimFacin` | Claim Fac In |
| `getPayAttachment` | Claim Fac In |
| `getPaymentClaim` | Claim Fac In |
| `getPremiumPaidOnMarine` | Claim Fac In |

### 11.1 Alamat tujuan = konfigurasi, kecuali satu

`[terverifikasi]` **19 dari 20** memakai `<pyBaseURLSelectionType>SETTING` dengan
`<pyBaseURLSetting>LinkService!LinkService`. Rule `SystemSettings` `LinkService` ada di kelima
modul, seluruhnya identik: `<pyPurpose>LinkService`, `<pySetting>=ResponLink.URL`.
**Nilai URL sebenarnya tidak ada di korpus** — ia berada di Dynamic System Setting yang tidak ikut
diekspor. Ini pola yang benar dan harus dipertahankan sebagai env var di arsitektur Go.

`[terverifikasi]` **Satu rule menyimpang**: `Claim Fac In/ConnectREST/getPremiumPaidOnMarine.xml`
(`ASM-FW-GCNMFW-WORK-PNC / GETPREMIUMPAIDONMARINE`) memakai
`<pyBaseURLSelectionType>URL` dan memuat URL literal berisi host, port, path, **serta contoh nilai
data** (nomor polis dan case id) di query string.

Dua hal yang perlu dicatat, bukan diperlakukan sebagai fakta arsitektur:

1. Path pada URL literal itu (`getPaymentDataClaim`) **tidak sama** dengan nama rule
   (`getPremiumPaidOnMarine`) — sekali lagi, nama rule tidak boleh dipercaya.
2. URL literal tersebut memuat contoh data yang menyerupai data nyata. Nilainya **tidak disalin**
   ke dokumen ini dengan sengaja; rujuk file aslinya bila perlu.

### 11.2 Host ter-hardcode tersebar jauh di luar ConnectREST

`[terverifikasi]` URL literal ditemukan di **48 file** batch 2, dan mayoritas **bukan** di
`ConnectREST` melainkan di `Section`, `FlowAction`, `Activity`, dan `Harness`:

| Modul / tipe | File |
| --- | ---: |
| RNW Fac In / Section | 6 |
| NB FacIn / Section | 6 |
| Endorsment Fac In / Section | 6 |
| RNW Fac In / FlowAction | 4 |
| NB FacIn / FlowAction | 4 |
| NB FacIn / Activity | 4 |
| Endorsment Fac In / FlowAction | 4 |
| RNW Fac In / Activity | 3 |
| Endorsment Fac In / Activity | 3 |
| Claim Fac In / Activity | 3 |
| Claim Fac In / Section | 2 |
| Claim Fac In / Harness | 2 |
| Claim Fac In / ConnectREST | 1 |

Host yang muncul (host saja; path tidak disalin):

| Host | Kemunculan | Catatan |
| --- | ---: | --- |
| `view.officeapps.live.com` | 12 | SaaS pihak ketiga (penampil dokumen Office) |
| `app.sinarmas.co.id` | 12 | domain grup Sinarmas |
| `192.168.105.112` (dan `:80`) | 24 | IP internal |
| `ssdecamwin03:7070` | 10 | hostname internal |
| `ssdecamwin02:7070` | 9 | hostname internal |
| `sdvpwin105:8383` | 3 | hostname internal |
| `192.168.105.116:80` | 3 | IP internal |
| `pega.nusantarare.com:80` | 2 | server Pega |
| `10.100.10.75:7315` | 2 | integration web service |
| `appdev.nusantarare.com` | 1 | **hostname lingkungan DEV** |

Perintah audit:

```
grep -rhoE "https?://[^/< \"']+" "<modul>" --include="*.xml" | grep -v "pega.com" | sort | uniq -c
```

Kehadiran `appdev.nusantarare.com` berarti korpus memuat campuran alamat dev dan non-dev.
Ditambahkan sebagai **OQ-018**.

## 12. Batch 2 — dua temuan metode yang mengoreksi asumsi inventaris

### 12.1 Nama folder BUKAN tipe rule

`[terverifikasi]` Dua file terbukti bertipe lain dari foldernya:

| File | Folder | Tipe sebenarnya | Bukti |
| --- | --- | --- | --- |
| `Claim Fac In/Activity/InsertLogServiceClaim.xml` | `Activity` | `RULE-CONNECT-SQL` | `<pxInsName>ASM-FW-GCNMFW-WORK!RNM!INSERTLOGSERVICECLAIM` |
| `Claim Fac In/Section/CedingCedant.xml` | `Section` | `RULE-HTML-HARNESS` | `<pxObjClass>Rule-HTML-Harness` |

Karena itu tipe rule di seluruh artefak diambil dari `<pzOriginalInstanceKey>` / `<pxObjClass>`,
**bukan** dari nama folder. Matriks `../README.md` §5.1 menghitung **file per folder**, bukan rule
per tipe — perbedaan yang kecil di batch ini (2 file) tetapi harus disebut agar tidak disalahartikan.

### 12.2 Lima file tanpa `<pzOriginalInstanceKey>`

`[terverifikasi]` Berbeda dari batch 1 (yang 1.149 filenya semua punya tag itu), batch 2 memuat
5 file tanpa `<pzOriginalInstanceKey>` sama sekali:

| File | `<pxInsName>` | `<pxObjClass>` |
| --- | --- | --- |
| `NB FacIn/Section/ViewObjectSavior.xml` | `ASM-FW-GISFW-WORK!VIEWOBJECTSAVIOR` | `Rule-HTML-Section` |
| `RNW Fac In/Section/ViewObjectSavior.xml` | sama | `Rule-HTML-Section` |
| `Endorsment Fac In/Section/ViewObjectSavior.xml` | sama | `Rule-HTML-Section` |
| `NB FacIn/When/hasPrimaryPage.xml` | `@BASECLASS!HASPRIMARYPAGE` | `Rule-Obj-When` |
| `RNW Fac In/When/hasPrimaryPage.xml` | sama | `Rule-Obj-When` |

Untuk kelima file ini, silsilah salinan tidak dapat ditelusuri. Identitas tetap terbaca dari
`<pxInsName>`. `[dugaan]` keduanya rule bawaan Pega (`hasPrimaryPage` adalah When standar), tetapi
**belum terverifikasi** — terkait OQ-009.

## 14. Batch 3 — angka

Domain claim non-facultative. Enam modul, dikerjakan berurutan kecil → besar.

| Ukuran | Nilai |
| --- | ---: |
| File XML | 871 |
| **Identitas rule unik** (tipe + class + nama) | **623** |
| File yang merupakan salinan identitas lintas modul | 248 |
| Tabel/view unik terlihat | 56 |
| Stored procedure/function unik | 19 (8 di antaranya baru) |
| Rule `ConnectREST` | 24 (10 nama service unik) |
| RDBList menurut jenis SQL | `QUERY`=132, `PLSQL`=54 |
| Prefix kunci RDBList | `RNM`=88, `GCNM`=51, `ASM`=47 |

Per modul — seluruhnya **cocok** dengan `../README.md` §5.1:

| Modul | File | Identitas unik | Class unik | Nama file >1 tipe | `When` sama di >1 class | Class bukan-aplikasi |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Komite Claim Life | 47 | 47 | 18 | 0 | 0 | 5 |
| Komite Claim Non Prop | 59 | 59 | 16 | 0 | 0 | 7 |
| Komite Claim Prop | 80 | 80 | 19 | 2 | 0 | 9 |
| Claim Life | 136 | 136 | 23 | 3 | 0 | 7 |
| Claim Prop | 270 | 270 | 38 | 11 | 0 | 27 |
| Claim Non Prop | 279 | 279 | 48 | 17 | 0 | 29 |

`[terverifikasi]` **Tidak ada identitas ganda di dalam satu modul** pada keenam modul — berbeda dari
batch 2 yang punya 5 kasus. Duplikasi di batch 3 seluruhnya lintas modul.

Catatan: folder `DataPage` ada di `Komite Claim Life` tetapi **kosong** (0 file XML) — konsisten
dengan §5.1 yang mencatat `-`. Folder `excludeXML` di `Komite Claim Non Prop` berisi 1 file yang
**terverifikasi bertipe `RULE-OBJ-ACTIVITY`** (lihat §19.2).

Perintah audit sama dengan batch 2:

```
find "<modul>" -type f -name "*.xml" | wc -l
find "<modul>" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
```

## 15. Batch 3 — trio Claim TIDAK seperti trio facultative

Ini fokus khusus batch 3. Hasilnya **berlawanan** dengan pola batch 2.

### 15.1 Overlap antar-modul

| Pasangan modul | Identitas dibagi |
| --- | ---: |
| Claim Non Prop ↔ Claim Prop | **102** |
| Claim Prop ↔ Komite Claim Prop | 55 |
| Claim Non Prop ↔ Komite Claim Prop | 36 |
| Komite Claim Non Prop ↔ Komite Claim Prop | 35 |
| Claim Non Prop ↔ Komite Claim Non Prop | 34 |
| Claim Prop ↔ Komite Claim Non Prop | 29 |
| Claim Life ↔ Komite Claim Life | 26 |
| Claim Life ↔ Claim Prop | 19 |
| Komite Claim Life ↔ Komite Claim Prop | 18 |
| Claim Prop ↔ Komite Claim Life | 18 |
| Claim Life ↔ Komite Claim Prop | 18 |
| Claim Non Prop ↔ Komite Claim Life | 17 |
| Claim Life ↔ Claim Non Prop | 16 |
| Komite Claim Life ↔ Komite Claim Non Prop | 15 |
| Claim Life ↔ Komite Claim Non Prop | 12 |

### 15.2 Trio Claim: bukan satu ruleset

`[terverifikasi]` Overlap tertinggi di trio Claim hanya **102 identitas** (Claim Non Prop ↔ Claim
Prop), dari modul berukuran 279 dan 270 identitas — sekitar **37%**. `Claim Life` jauh lebih
terpisah lagi: hanya 19 identitas dibagi dengan `Claim Prop` dan 16 dengan `Claim Non Prop`, dari
136 identitas miliknya (**12–14%**).

Bandingkan dengan batch 2, di mana `RNW Fac In` berbagi **99,0%** identitasnya dengan `NB FacIn`.

**Kesimpulan terukur:** pola "satu ruleset berdiskriminator" yang terlihat di domain facultative
**tidak berlaku** di domain claim. Ketiga modul Claim adalah basis kode yang sebagian besar
terpisah. `[dugaan]` pembedanya bukan tahap siklus melainkan **jenis bisnis** (Life vs
Proportional vs Non-Proportional) yang memang punya perhitungan berbeda — **belum terverifikasi**.
Ditambahkan sebagai **OQ-019**.

### 15.3 Trio Komite: juga terpisah

Overlap antar-Komite rendah: Komite NP ↔ Komite Prop 35, Komite Life ↔ Komite Prop 18, Komite Life
↔ Komite NP 15 — dari modul berukuran 59, 80, dan 47 identitas.

### 15.4 Komite BUKAN wrapper approval tanpa data sendiri

Pertanyaan yang diminta diukur. Jawabannya terukur, dan jawabannya **tidak**.

| Modul Komite | Identitas | Dibagi dengan modul Claim pasangannya | % |
| --- | ---: | ---: | ---: |
| Komite Claim Life | 47 | 26 (Claim Life) | 55,3% |
| Komite Claim Non Prop | 59 | 34 (Claim Non Prop) | 57,6% |
| Komite Claim Prop | 80 | 55 (Claim Prop) | 68,8% |

`[terverifikasi]` Gabungan ketiga modul Komite memiliki **133 identitas**, dan **59 di antaranya
(44,4%) tidak muncul di modul Claim mana pun**. Modul Komite membawa rule miliknya sendiri dalam
jumlah besar — termasuk class khusus seperti `ASM-FW-GCNMFW-WORK-KOMITELIFE`
(`Komite Claim Life/Flow/KomiteLife_Flow.xml`).

`[dugaan]` Komite adalah tahap approval dengan model data dan layar sendiri, bukan lapisan tipis di
atas modul Claim. **Belum terverifikasi** — pembuktian butuh D2. Ini **bukan** penetapan bounded
context; itu STEP D4.

## 16. Batch 3 — objek database

### 16.1 Tidak ada database link di domain claim

`[terverifikasi]` **Nol objek `@ASMD.SINARMAS.CO.ID`** di seluruh 871 file batch 3. Sejauh
pengukuran D1, ketergantungan pada database luar (OQ-017) **terbatas pada domain facultative**
(tabel rate dan kurs). Temuan negatif ini mempersempit cakupan OQ-017.

Perintah audit: cari pola `@` pada daftar tabel hasil ekstraksi `<pyBrowseSQL>`.

### 16.2 Skema: tidak ada yang baru

| Skema | Rujukan tabel | Rujukan procedure |
| --- | ---: | ---: |
| `POOLDATA` | 103 | dominan |
| `GL` | 5 | ya (`GL.F_GET_EMAIL`) |
| `REINSURANCE` | 4 | — |
| `ARASAPAS` | 1 | — |

Ditambah built-in `DBMS_LOB`. Seluruhnya sudah tercatat di OQ-016; batch 3 **tidak menambah skema
baru**. Objek baru yang menonjol: `REINSURANCE.TRLOSS_DETAIL_T` (4 rule), `BANKACCOUNT` (8 rule).

### 16.3 Delapan stored procedure baru — seluruhnya seputar akseptasi klaim

`[terverifikasi]` 19 procedure unik di batch 3, **8 belum pernah muncul** di batch 1/2:

| Procedure baru | Dipanggil (rule) |
| --- | ---: |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` | 4 |
| `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` | 3 |
| `POOLDATA.PEGA_M_CAUSE_OF_LOSS` | 2 |
| `POOLDATA.XOL2_AKSEP_KLAIM` | 1 |
| `POOLDATA.PEGA_JSON_OS_AKSEP_SUBJECTIVITY` | 1 |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT` | 1 |
| `POOLDATA.GENERATE_NOCLMTREATYIN` | 1 |
| `POOLDATA.GENERATE_NOCLMTRTYINTEMP` | 1 |

Yang paling sering dipanggil di batch 3 adalah `POOLDATA.PEGA_JSON_KLAIM_PNC` (6),
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` (6), `POOLDATA.GET_TOKEN_STORAGE` (6),
`POOLDATA.DIRECTTOKASIR_LOG` (5), `GL.F_GET_EMAIL` (5).

**Body seluruhnya tetap tidak ada di korpus** (OQ-002). Peran yang tertulis di atas adalah
**nama apa adanya** — tidak ada yang ditebak. Akhiran `TNP`, `TRT`, `PNC`, `XOL2`
**kepanjangan belum terverifikasi**.

### 16.4 JSON opaque juga di domain klaim

Objek/kolom berlabel JSON dan akseptasi yang muncul: `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` (7 rule),
`OS_AKSEPTASI_KLAIM` (6), `JSON_POLIS` (5), `JSON_KLAIM` (3), `HISTORYAKSEPTASIPEGA` (2).
Menguatkan **OQ-012** dan memperluasnya ke domain klaim. Struktur JSON di dalamnya tetap tidak ada
di korpus.

## 17. Batch 3 — integrasi eksternal (`ConnectREST`)

24 rule, **10 nama service unik** (apa adanya dari `<pyServiceName>`):

| pyServiceName | Modul |
| --- | --- |
| `ServiceGoogle` | keenam modul |
| `KonversiKlaimNonLife` | Claim Non Prop, Claim Prop, Komite Claim Non Prop, Komite Claim Prop |
| `SendAcceptationToKasir` | Claim Non Prop, Claim Prop, Komite Claim Non Prop, Komite Claim Prop |
| `getPayAttachment` | Claim Non Prop, Claim Prop |
| `getPremiumPaidOnTreatyIn` | Claim Non Prop, Claim Prop |
| `insertClaimFinalOrClosed_NP` | Claim Non Prop, Komite Claim Non Prop |
| `convertJsonNusareToProductionClaimLife` | Claim Life |
| `GetDtlPaymentClaim` | Claim Prop |
| `InsertClaimOutstanding_NP` | Claim Non Prop |
| `insertClaimReject_NP` | Komite Claim Non Prop |

### 17.1 Empat rule hardcode URL — naik tajam dari batch 2

`[terverifikasi]` **4 dari 24** memakai `<pyBaseURLSelectionType>URL` dengan URL literal, bukan
`SETTING`. Di batch 2 hanya 1 dari 20.

| Rule ber-URL literal | Modul |
| --- | --- |
| `InsertClaimOutstanding_NP` | Claim Non Prop |
| `insertClaimFinalOrClosed_NP` | Claim Non Prop |
| `insertClaimFinalOrClosed_NP` | Komite Claim Non Prop |
| `insertClaimReject_NP` | Komite Claim Non Prop |

Keempatnya adalah service pencatatan klaim Non-Proportional (`_NP`). 20 rule sisanya memakai
`SETTING` → `LinkService!LinkService`, pola yang benar.

URL literal ditemukan di **24 file** batch 3 secara keseluruhan (bukan hanya ConnectREST), dengan
host berikut — host saja, path tidak disalin:

| Host | Kemunculan |
| --- | ---: |
| `192.168.105.116:80` | 12 |
| `ssdecamwin03:7070` | 8 |
| `10.100.10.75:7315` | 8 |
| `view.officeapps.live.com` | 6 |
| `pega.nusantarare.com:80` | 4 |
| **`appdev.nusantarare.com`** | **3** |
| `sdvpwin105:8383` | 2 |

`appdev.nusantarare.com` muncul 3× di batch 3 (batch 2: 1×) — memperkuat **OQ-018**.

## 18. Batch 3 — 25 identitas berbeda isi (OQ-011)

| Kategori | Jumlah identitas |
| ---: | --- |
| Muncul di >1 modul, isi **identik** setelah normalisasi 18 tag | **122** |
| Muncul di >1 modul, isi **benar-benar berbeda** | **25** |

395 file dihitung hash ternormalisasinya. Metode sama dengan §3.1 dan §9.

**Angka korpus penuh (§22.8): 533 identitas berbeda isi.** Penjumlahan per batch (558) menghitung
identitas lintas batch lebih dari sekali. Seluruhnya terdaftar
beserta **semua** path variannya di **`_oq011-konflik-isi.md`**. Tidak ada satu varian pun yang
dipilih sebagai benar.

### 18.1 Empat rule `When` di `@BASECLASS` punya empat isi berbeda

`[terverifikasi]` Konflik paling tajam di seluruh D1 sejauh ini. Empat rule `When` — yang berperan
sebagai **guard percabangan** — masing-masing muncul di 4 file dengan **4 isi yang semuanya
berbeda**:

| Identitas | Varian | Versi isi |
| --- | ---: | ---: |
| `@BASECLASS / ISCLM` | 4 | **4** |
| `@BASECLASS / ISCLMP` | 4 | **4** |
| `@BASECLASS / ISCLMNP` | 4 | **4** |
| `@BASECLASS / ISPEGASYARIAH` | 4 | **4** |

Keempatnya tersebar di `Komite Claim Non Prop/When/`, `Komite Claim Prop/When/`,
`Claim Prop/When/`, dan `Claim Non Prop/When/`.

Artinya: satu nama `When` yang sama, di class yang sama, berperilaku **berlainan di setiap modul**
— dan tidak ada dua modul pun yang sepakat. Untuk rule guard, ini konsekuensinya langsung ke
percabangan proses.

`ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` — **arti dan kepanjangan belum terverifikasi**.
Jangan ditebak dari nama.

### 18.2 Konflik lain di batch 3

Satu identitas dengan 3 versi isi: `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL`,
tersebar di **keenam** modul batch 3 (`RDBList/GetTokenStorage_SQL.xml`).

Contoh konflik 2 versi:

- `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS / BROWSEVMCAUSEOFLOSS_RD` — `Claim Prop/ReportDefinition/`
  vs `Claim Non Prop/ReportDefinition/`
- `ASM-FW-GCNMFW-WORK-CLAIMTREATY / COUNTLISTCLAIMAMOUNTIDR` — `Komite Claim Prop/Activity/`
  vs `Claim Prop/Activity/`
- `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTHISTORYAKSEPTASIPEGA_SQL` —
  `Komite Claim Non Prop/RDBList/` vs `Komite Claim Prop/RDBList/`

## 19. Batch 3 — temuan khusus domain klaim

### 19.1 Kode `PaymentType` dan `TransferType` menggerbangi logika, artinya tidak diketahui

`[terverifikasi]` Nilai literal yang diuji di dalam rule, **dicatat apa adanya**:

- `PaymentType` — nilai **0, 1, 2, 3, 4, 5, 6, 7**
- `TransferType` — nilai **1, 2, 3, 4**

Yang penting: nilai-nilai ini **dikelompokkan berbeda-beda** untuk mengambil jalur logika berbeda.
Contoh dari `Claim Prop/Activity/HitServiceToKasir_Act.xml`
(`ASM-FW-GCNMFW-WORK / HITSERVICETOKASIR_ACT`) — tiga kelompok dalam satu rule:

```
.PaymentType = 1 || .PaymentType = 2 || .PaymentType = 5
.PaymentType = 4 || .PaymentType = 6
.PaymentType =  3
```

Pengelompokan lain, dari `Claim Prop/Activity/AttachmentProtect_ACT.xml`:

```
Local.PaymentType=1||Local.PaymentType=2||Local.PaymentType=3
Local.PaymentType=2||Local.PaymentType=4
```

Dan `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` menguji `PaymentType==7` serta
`PaymentType!=7`.

Rule lain yang memakai kode ini sebagai guard: `Claim Non Prop/Activity/HitServiceToKasir_Act.xml`,
`Claim Non Prop/Activity/SetInterimXOL_Act.xml`,
`Komite Claim Prop/Activity/HitServiceToKasirKMT_Act.xml`,
`Komite Claim Non Prop/Activity/HitServiceToKasirKMT_Act.xml`,
`Komite Claim Prop/Section/ShowTransfer.xml`.

**Arti tiap nilai belum terverifikasi.** Korpus tidak memuat tabel kode maupun label yang
menjelaskannya, dan **tidak boleh ditebak** dari nama rule. Karena kode ini menentukan percabangan
pembayaran klaim, ini menjadi **OQ-020**.

### 19.2 `excludeXML` berisi rule Activity yang nyata

`[terverifikasi]` `Komite Claim Non Prop/excludeXML/GetBase64Attachment.xml` bertipe
**`RULE-OBJ-ACTIVITY`**, beridentitas `ASM-FW-GCNMFW-WORK / GETBASE64ATTACHMENT`, hasil salinan
dari `WORK- / LOADATTACHMENTDATA`.

Jadi folder `excludeXML` **bukan** berisi berkas kosong atau placeholder — isinya rule Activity
lengkap. Apakah rule itu aktif di production tetap **belum terverifikasi** (**OQ-003** tetap
terbuka), tetapi pertanyaannya kini lebih tajam: ini rule nyata yang sengaja ditempatkan di folder
bernama "exclude".

### 19.3 `GCNMFW` adalah framework domain klaim, bukan khusus facultative

Di batch 2 saya mencatat `[dugaan]` bahwa `ASM-FW-GCNMFW-*` khusus claim facultative. Batch 3
memperluasnya: framework itu **mendominasi seluruh domain klaim**.

| Modul | `ASM-FW-GISFW-*` | `ASM-FW-GCNMFW-*` |
| --- | ---: | ---: |
| Claim Prop | 53 | **190** |
| Claim Non Prop | 81 | **169** |
| Claim Life | 74 | 55 |
| Komite Claim Prop | 24 | **47** |
| Komite Claim Non Prop | 15 | **37** |
| Komite Claim Life | 25 | 17 |

`[dugaan]` `GCNMFW` adalah framework klaim dan `GISFW` framework underwriting/produksi; modul klaim
memakai keduanya karena membaca data produksi. **Belum terverifikasi** — kepanjangan kedua
singkatan itu tetap tidak dijelaskan korpus (**OQ-008**).

### 19.4 Dua rule `Flow` tanpa `<pzOriginalInstanceKey>`

`[terverifikasi]` `Komite Claim Life/Flow/KomiteLife_Flow.xml` dan
`Claim Life/Flow/Register_Flow.xml` tidak memiliki tag itu sama sekali. Keduanya terkonfirmasi
`Rule-Obj-Flow` lewat `<pxObjClass>`, dengan identitas
`ASM-FW-GCNMFW-WORK-KOMITELIFE / KOMITELIFE_FLOW` dan
`ASM-FW-GCNMFW-WORK-CLAIMLIFE / REGISTER_FLOW`.

Extractor kini memakai `<pxObjClass>` sebagai **fallback** tipe rule dan menandainya
`(pxObjClass)` agar asal informasinya jelas. Silsilah salinan kedua rule ini tidak dapat ditelusuri.

### 19.5 Identitas orang ter-hardcode sebagai guard — temuan lintas batch

Pola ini ditemukan saat memeriksa batch 3, lalu dicari mundur ke batch 1 dan 2. **Pemeriksaan ini
tidak dilakukan pada batch 1 dan 2 sebelumnya** — jadi ini temuan baru atas modul lama, bukan
temuan batch 3 semata.

`[terverifikasi]` Angka batch 3 adalah 45 file di 6 modul. **Sapuan korpus penuh di batch 4
menaikkannya menjadi 66 file di 9 modul** (pola pencarian batch 3 hanya mengenali kutip ganda).
Angka final dan rinciannya ada di **OQ-021** pada `../open-questions.md`.

Guard ini berada terutama di `<pyStepsPreCondParamsWhen>` (precondition eksekusi step, 193×),
`<pyCondition>` (58×), `<pyReadOnlyCondition>` (18×), dan `<pyContainerVisibleWhen>` (8×) — jadi
menggerbangi **eksekusi langkah, keterubahan field, dan visibilitas layar**.

Rule yang digerbangi bukan kosmetik: termasuk `ASM-FW-GISFW-WORK / COUNTGROSSPREMI_ACT`
(perhitungan premi) dan `ASM-FW-GCNMFW-WORK-KOMITE / KOMITEPOST_REJECT` (keputusan komite klaim).

Terpisah dan lebih menjanjikan: **35 file di 14 modul** memakai `OperatorID.pyPosition` dengan nilai
yang **terlihat seperti nama peran** (`IT Developer`, `ReasLifeAdmin`, `ReasLifeSPV`,
`ReasLifeMedicalAdvisor`, `Admin`, `SPV A`, `SPV B`). `[dugaan]` konsep peran sudah ada, hanya
ditegakkan lewat perbandingan string.

**Nilai nama orang sengaja tidak dicatat di berkas ini** — lihat catatan data pribadi di OQ-021.

Perintah audit:

```
grep -rlE 'OperatorID\.pyUser(Identifier|Name)[ ]*[=!]+[ ]*"' "<modul>" --include="*.xml"
grep -rhoE 'OperatorID\.pyUser(Identifier|Name)[ ]*[=!]+[ ]*"[^"]{1,40}"' "<modul>" --include="*.xml"
```

Arti tiap perbandingan **belum terverifikasi** dan tidak boleh ditebak. → **OQ-021**, berkaitan
dengan **OQ-007**.

## 20. Yang belum dikerjakan — **catatan per akhir batch 3**

> Bagian ini ditulis saat batch 3 selesai dan **dipertahankan apa adanya**. Butir pertama sudah
> terselesaikan di batch 4. Daftar yang berlaku sekarang ada di **§22.10**.

- ~~**Batch 4 (681 file)** — inventaris belum ada.~~ **Selesai di batch 4** (682 file; dan ternyata
  hanya **empat** dari lima modul yang tanpa rule `Flow` — `PremiumList Life` punya Flow, lihat
  §21.2). Dari lima modul batch 4, yang tanpa `Flow`: Master Contract Retro Life, Master Product
  Name Life, Treaty Contract Out (3 modul), ditambah Treaty In dan Treaty In Adjustment dari batch 1
  → **5 modul korpus** tanpa `Flow`.
- **Rule tertanam belum diperiksa untuk batch 1, 3, dan 4.** Pemeriksaan `<pyIncludedRuleXML>` baru
  dilakukan untuk batch 2. **Masih belum terukur** per akhir D1.
- **Rule tertanam di dalam file — celah kecil, sudah terukur.** 208 file batch 2 memuat lebih dari
  satu `<pxInsName>` (7.310 kemunculan pada 6.667 file, selisih **643**), bersumber dari
  `<pyIncludedRuleXML>`. Extractor hanya mengambil rule pertama tiap file, jadi rule tertanam tidak
  masuk tabel inventaris. File dengan kemunculan terbanyak:
  `NB FacIn/Activity/SetAutoAccept_Act.xml` (54), `NB FacIn/Activity/GetPolicyData_ACT.xml` (49),
  `NB FacIn/Activity/GetData_ACT.xml` (49).

  Dampaknya sudah diukur: dari **176 identitas tertanam unik**, **97 sudah ada sebagai file
  tersendiri** (jadi sudah terinventarisasi), dan **79 tidak**. Dari 79 itu, **76 berada di class
  bawaan Pega** — `@BASECLASS` (54, mis. `ROUND`, `UPPERCASE`, `ABSOLUTEVALUEINTEGER`), `WORK-`
  (14), `PEGAGADGET-PULSE` (5), `CODE-PEGA-LIST`, `DATA-WORKATTACH-FILE`, `ASSIGN-`.

  **Hanya 3 rule aplikasi yang benar-benar hilang dari inventaris:**

  | Identitas | Tertanam di |
  | --- | --- |
  | `ASM-FW-GISFW-DATA-OFFERTREATYIN / SOURCEHIERARKI` | `RNW Fac In/Harness/SOB.xml` |
  | `ASM-FW-GISFW-DATA-OFFERTREATYIN / CEDINGCOHIERARKI` | `RNW Fac In/Harness/CedingCompany.xml` |
  | `ASM-FW-GISFW-WORK / INPDMPACKAGE` | `Endorsment Fac In/Section/InputDtlCoverage_FacIn_IsUW.xml` |

  Perintah audit: kumpulkan **semua** `<pxInsName>` per file (bukan hanya yang pertama), tandai
  kemunculan ke-1 sebagai identitas utama dan sisanya sebagai tertanam, lalu selisihkan.
  Hal yang sama **belum diperiksa untuk batch 1**.
- Isi rule `When`, `DecisionTable`, `DataTransform`, `DecisionTree` belum dibaca (hanya didaftar).
- Tabel yang disentuh lewat `ReportDefinition` dan `DataPage` belum diekstraksi — baru `RDBList`.
- Dependency antar-Section lewat `<pyInclude>` baru dihitung jumlahnya, belum dipetakan.
- Perbandingan isi untuk **533** identitas berbeda — pekerjaan D2, dan diblokir OQ-011.
- **Daftar tag normalisasi mungkin masih belum lengkap.** Diperluas sekali dari 7 ke 18 tag dan
  angka konflik turun 44%. Bila ditemukan tag provenance lain, angka 533 dapat turun lagi. Cara
  memeriksa: ambil sejumlah pasangan konflik, diff setelah normalisasi, lihat nama tag yang tersisa.

## 21. Batch 4 — angka dan temuan

Domain life, master & outward. Lima modul, 682 file, **623 identitas rule unik**.

| Modul | File | Identitas unik | Class unik | Nama file >1 tipe | `When` sama di >1 class |
| --- | ---: | ---: | ---: | ---: | ---: |
| Endorsement Life | 75 | 75 | 20 | 3 | 0 |
| Master Contract Retro Life | 66 | 66 | 14 | 0 | 0 |
| Master Product Name Life | 114 | 114 | 23 | 3 | 0 |
| PremiumList Life | 124 | 124 | 22 | 4 | 0 |
| Treaty Contract Out | 303 | 303 | 26 | 8 | 0 |

`PremiumList Life` tercatat 124 sedangkan `../README.md` §5.1 mencatat 123 — matriks itu menghitung
**per folder tipe**, dan satu file berada di luar folder tipe. Lihat §21.2.

### 21.1 Tidak ada "ruleset Life bersama"

Overlap identitas antar modul batch 4 **rendah di semua pasangan**:

| Pasangan | Identitas dibagi |
| --- | ---: |
| Endorsement Life ↔ PremiumList Life | 23 |
| Master Product Name Life ↔ Treaty Contract Out | 17 |
| Master Contract Retro Life ↔ Master Product Name Life | 8 |
| Master Product Name Life ↔ PremiumList Life | 6 |
| Master Contract Retro Life ↔ Treaty Contract Out | 6 |
| Endorsement Life ↔ Master Product Name Life | 5 |
| (sisanya) | ≤ 3 |

`[terverifikasi]` Tidak ada pasangan yang mendekati pola facultative (99,0%) maupun treaty
(320 identitas). **Keempat modul Life berdiri sendiri-sendiri.**

73 dari 623 identitas batch 4 juga muncul di batch 1–3 (Master Product Name Life 40,
PremiumList Life 26, Treaty Contract Out 25, Endorsement Life 11, Master Contract Retro Life 5) —
pemakaian ulang lintas domain yang sederhana, bukan basis kode bersama.

### 21.2 OQ-004 terjawab, dan OQ-005 ikut terkoreksi

`[terverifikasi]` `PremiumList Life/InputPolicyHolder.xml` — satu-satunya file XML korpus yang
berada di luar folder tipe — adalah rule **`RULE-OBJ-FLOW`**
(`ASM-FW-GISFW-WORK-LIFE / INPUTPOLICYHOLDER`, `<pxObjClass>Rule-Obj-Flow`,
`<pyStartActivity>Start1`).

Konsekuensinya, daftar "6 modul tanpa rule `Flow`" dari STEP D0 **keliru**: `PremiumList Life`
punya Flow, hanya tidak di folder `Flow/`. Angka korpus yang benar: **24 rule `Flow`**
(22 `RULE-OBJ-FLOW` + 2 yang tipenya dibaca dari `<pxObjClass>`), dan **5 modul tanpa `Flow`**:
Master Contract Retro Life, Master Product Name Life, Treaty Contract Out, Treaty In,
Treaty In Adjustment.

Rincian di `PremiumList Life.md` §8.

### 21.3 `Treaty Contract Out` — isinya tidak menunjukkan treaty *outward*

Ini fokus khusus batch 4: memeriksa apakah nama folder menyesatkan. **Buktinya mengatakan ya.**

`[terverifikasi]` Objek database yang disentuh modul ini seluruhnya bertema **kontrak treaty**,
tidak satu pun bernama `*OUT*`:

| Objek | Rule |
| --- | ---: |
| `M_PROPORTIONALARRG` | 10 |
| `MTREATYSECURITY` | 5 |
| `TREATYREINSURER` | 3 |
| `TREATYBUSINESS` | 3 |
| `TREATYEXCHANGE`, `TREATYCONTRACT`, `M_TREATYYEAR`, `PROPORTIONALARRG` | 1 masing-masing |

Empat objek **hanya** dirujuk modul ini di seluruh korpus: `M_TREATYYEAR`, `TREATYEXCHANGE`,
`MTREATYSECURITY`, `M_PROPORTIONALARRG`.

`[terverifikasi]` **Objek treaty outward justru dirujuk modul lain, bukan modul ini.** Di seluruh
korpus, `M_TREATY_OUT`, `TREATY_OUT2`, `M_TREATY_OUT_DETAIL`, dan `FACOUTPRODUCTION` dirujuk oleh
NB Treaty In, Treaty In Adjustment, EDM Treaty In, Claim Non Prop, NB FacIn, RNW Fac In, dan
Endorsment Fac In — **nol rujukan dari `Treaty Contract Out`**.

`[terverifikasi]` Dari 303 rule, hanya **5** menyebut `TreatyOut` di namanya, dan kelimanya soal
**lampiran**: `LOADATTACHMENTTREATYOUT`, `TREATYOUTDOWNLOADALL_ACT`, `TREATYOUTDOWNLOADONE`,
`TREATYOUTSAVEATTACHMENT`, `TREATYOUTATTACHCONTENT`. (Dua rule lain mengandung "OUT" hanya karena
`SETOUTPUTPARAM` — bukan outward.) Sebaliknya, dua rule menyebut treaty **inward** secara eksplisit:
`BROWSEDELETEROWTREATYINCONTRACT` dan `SETCATEGORYATTACHTREATYIN`, dan class
`ASM-FW-GISFW-INT-TREATY_IN` dipakai 5 rule di modul ini.

`[terverifikasi]` 16 rule `CANCELACTIVITY*` menamai **jenis klausul kontrak treaty**: `BORDEREAUX`,
`CASHLOSSLIMIT`, `CLAIMCOORPERATION`, `EPI`, `EXGRATIA`, `EXGRATIALIMITCHILD`, `FACIN`, `FACINLIST`,
`PLA`, `PPORTFOLIO`, `PROFITCOMMISION`, `PTERRLIMIT`, `RICOMM`, `TREATYCONTRACT`, `TREATYLIMIT`,
`TREATYLIMITCHILD`.

**Kesimpulan yang boleh ditarik sekarang:** isi modul ini adalah **master/klausul kontrak treaty**,
dan tidak ada bukti bahwa ia menangani treaty outward selain penanganan lampiran.
**Apakah "Out" dalam nama folder merujuk sesuatu yang tidak terlihat di rule ini — belum
terverifikasi.** → **OQ-022**. Penetapan konteks tetap STEP D4, bukan di sini.

Catatan lain `[terverifikasi]`: **202 dari 303 rule (67%) berada di class `@BASECLASS`** — proporsi
tertinggi di seluruh korpus. Kaitkan **OQ-009**.

### 21.4 Enumerasi baru yang menggerbangi logika

Nilai literal dicatat apa adanya; **artinya belum terverifikasi** dan tidak boleh ditebak.

| Properti | Nilai literal yang diuji | Bukti |
| --- | --- | --- |
| `EdmType` | `1`, `3` | `Endorsement Life/Activity/SetPremi_EDM.xml`, `Endorsement Life/Activity/SaveCSVEDMLife.xml`, `Endorsement Life/Harness/InboxEndorsementLife.xml` |
| `ProRateType` | `1`, `2`, `3` | `PremiumList Life/Activity/Calculate1_Act.xml`, `PremiumList Life/Activity/SavePremiumList_Act.xml`, `Endorsement Life/Activity/Calculate1_Act.xml` |
| `.Type` (PremiumList) | `QR`, `QP`, `TP`, `TR` | `PremiumList Life/Activity/Calculate1_Act.xml`, `PremiumList Life/Activity/GetPLNumber_Act.xml`, `PremiumList Life/Activity/SubmitPremiumList_Act.xml` |

Kode `QR/QP/TP/TR` muncul di dalam `<pyStepsPreCondParamsWhen>` — jadi **menggerbangi eksekusi step**,
contoh: `.Type=="QR"` dan `.Type = "TR" || .Type = "TP"`. Bersama `EdmType` dan `ProRateType`,
ketiganya ditambahkan ke **OQ-020** (yang semula hanya tentang `PaymentType`/`TransferType`).

### 21.5 Database batch 4

- **Nol objek database link** `@ASMD.SINARMAS.CO.ID` `[terverifikasi]` — melengkapi sapuan 100%
  korpus untuk **OQ-017**.
- **Tidak ada skema baru**: `POOLDATA` (67 rujukan), `DBMS_LOB`, `DATAPEGA`, `ARASAPAS`.
- **27 stored procedure baru**, mayoritas bertema Life dan master kontrak:
  `POOLDATA.INSERTTREATYCONTRACT_LIFE`, `POOLDATA.INSERTREINSURER_LIFE`,
  `POOLDATA.INSERTSECURITYREINSURER_LIFE`, `POOLDATA.INSERTTREATYYEAR_LIFE`,
  `POOLDATA.INSERTBUSINESS_LIFE`, `POOLDATA.PEGA_M_PRODUCT_LIFE`,
  `POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE`, `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY`,
  `POOLDATA.PEGA_TREATYCONTRACT`, `POOLDATA.PEGA_TREATYREINSURER`, `POOLDATA.PEGA_TREATYYEAR`,
  `POOLDATA.PEGA_TREATYBUSINESS`, `POOLDATA.PEGA_PROPORTIONALARRG`,
  `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD`, `POOLDATA.GETQUARTER`, `POOLDATA.GETQUARTERRETRO`,
  `POOLDATA.PROSESCOPY`, dan lainnya → **OQ-002**.
- Objek JSON: `POOLDATA.JSON_POLIS` (7 rule), `JSON_POLIS` (3), `POOLDATA.JSON_OFFER_LIFE` (2)
  → **OQ-012**.

## 22. REKAP FINAL SELURUH KORPUS (20 modul) — D1 SELESAI

Seluruh angka di bagian ini dihitung dari dataset ekstraksi 9.369 baris yang mencakup **setiap file
XML di korpus**, pada 2026-09-12.

### 22.1 Angka pokok

| Ukuran | Nilai |
| --- | ---: |
| Modul | **20** |
| File XML | **9.369** |
| Ukuran korpus | 1.319,3 MB |
| **Identitas rule unik** (tipe + class + nama) | **4.444** |
| File yang merupakan salinan identitas lintas modul | 4.925 |
| Tabel/view unik yang terlihat | **216** |
| Stored procedure/function kustom tanpa body | **67** |
| Rule `ConnectREST` | 51 |
| Rule `RDBList` (`RULE-CONNECT-SQL`) | 1.146 |

**52,6% file adalah duplikat identitas lintas modul.** Volume pekerjaan migrasi sebenarnya diukur
dari 4.444 rule, bukan 9.369 file — **dengan syarat** salinannya identik, yang tidak selalu benar
(§22.5).

Perintah audit:

```
find . -type f -name "*.xml" -not -path "./OUTPUT_HASIL_RNM/*" | wc -l
# identitas unik: ekstrak <pxInsName> + tipe tiap file, lalu hitung kombinasi unik
```

### 22.2 Per domain

| Domain | Modul | File | Identitas unik |
| --- | ---: | ---: | ---: |
| Treaty inward | 4 | 1.149 | 729 |
| Facultative inward | 3 | 6.071 | 2.552 |
| Klaim (fac + non-fac + komite) | 8 | 1.467 | 1.018 |
| Life, master & outward | 5 | 682 | 623 |

Jumlah identitas per domain (4.922) melebihi identitas korpus (4.444) karena sebagian identitas
dipakai lintas domain.

**Pola ruleset bersama sangat berbeda antar domain** `[terverifikasi]`:

| Domain | Pola | Bukti |
| --- | --- | --- |
| Facultative inward | **satu ruleset** | `RNW Fac In` berbagi **99,0%** identitasnya dengan `NB FacIn`; 1.567 identitas ada di ketiga modul |
| Treaty inward | **sebagian besar satu** | `Treaty In` ↔ `Treaty In Adjustment` berbagi **320** identitas |
| Klaim | **terpisah** | overlap tertinggi trio Claim hanya 102 (~37%); `Claim Life` 12–14% |
| Life/master/outward | **terpisah** | overlap tertinggi hanya 23 identitas |

### 22.3 Tipe rule di seluruh korpus

| Tipe rule | Jumlah |
| --- | ---: |
| `RULE-OBJ-ACTIVITY` | 3.051 |
| `RULE-HTML-SECTION` | 1.683 (+3 dari `<pxObjClass>`) |
| `RULE-CONNECT-SQL` | 1.146 |
| `RULE-OBJ-FLOWACTION` | 964 |
| `RULE-OBJ-WHEN` | 833 (+2 dari `<pxObjClass>`) |
| `RULE-OBJ-MODEL` (DataTransform) | 622 |
| `RULE-OBJ-REPORT-DEFINITION` | 600 |
| `RULE-HTML-HARNESS` | 207 |
| `RULE-DECLARE-PAGES` (DataPage) | 111 |
| `RULE-CONNECT-REST` | 51 |
| `RULE-DECLARE-DECISIONTABLE` | 49 |
| **`RULE-OBJ-FLOW`** | **24** (22 + 2 dari `<pxObjClass>`) |
| `RULE-ADMIN-SYSTEM-SETTINGS` | 18 |
| `RULE-DECLARE-DECISIONTREE` | 3 |
| `RULE-OBJ-HTML` | 1 |
| `RULE-NAVIGATION` | 1 |
| **TOTAL** | **9.369** |

Tipe diambil dari `<pzOriginalInstanceKey>`, atau `<pxObjClass>` bila tag itu tidak ada (5 file).
**Bukan dari nama folder** — 2 file terbukti berbeda antara folder dan tipe sebenarnya
(`Claim Fac In/Activity/InsertLogServiceClaim.xml` = `RULE-CONNECT-SQL`;
`Claim Fac In/Section/CedingCedant.xml` = `RULE-HTML-HARNESS`).

### 22.4 Tabrakan nama di seluruh korpus

- **268 nama file dipakai di lebih dari satu tipe rule.** Mayoritas konvensi Pega
  (Section+FlowAction, Section+Harness), tetapi ada yang tipenya berlainan sifat:
  `FetchTreatyContractDetail.xml`, `GetAchievement.xml`, `GetTreatyName.xml`, `GetInsuredID.xml`,
  `InsertLogServiceClaim.xml`, `DeleteSecurityReinsurer.xml` (Activity + RDBList);
  `BrowseTreatyOutDetail.xml` (RDBList + ReportDefinition);
  `IsFlagOnGoingPolicy.xml` (When + DecisionTable);
  `ChangeUpKendaraan.xml` (Activity + DataTransform);
  `GetBase64Attachment.xml` (Activity + excludeXML);
  `InputEDMLife.xml` (Flow + FlowAction + Section).
- **55 nama rule `When` dipakai di lebih dari satu class** — jauh lebih banyak dari 2 yang terlihat
  di batch 1. Contoh: `ISMAINTENANCE` dan `ISCLAIM` masing-masing di **3 class** berbeda;
  `ISMBU` di `ASM-FW-GISFW-DATA` dan `ASM-SFAGIS-WORK-ENDORSEMENT`.

Menegaskan aturan `../README.md` §3.2: rule **wajib** ditulis `class / nama / tipe`.

### 22.5 Prefix kunci RDBList

| Prefix | Rule |
| --- | ---: |
| `ASM` | 751 |
| `RNM` | 315 |
| `GCNM` | 80 |
| **Total** | **1.146** |

Jenis SQL: `QUERY` 888, `PLSQL` 258. Arti ketiga prefix **belum terverifikasi** → **OQ-008**.

### 22.6 Skema Oracle dan objek lintas database

**11 skema lokal** terlihat (di luar built-in Oracle `DBMS_LOB`, `UTL_MATCH`, `DBMS_OUTPUT`):

| Skema | Rujukan |
| --- | ---: |
| `POOLDATA` | 539 |
| `DATAPEGA` | 12 |
| `NEW_UNDERWRITING` | 11 |
| `GENERAL` | 11 |
| `REINSURANCE` | 9 |
| `GL` | 7 |
| `FIRE` | 5 |
| `ARASAPAS` | 4 |
| `MBU` | 3 |
| `NEW_GENERAL` | 1 |

→ **OQ-016** (peran tiap skema, mana yang milik sistem lain).

**Database link `@ASMD.SINARMAS.CO.ID`** — 8 objek, **hanya di 3 modul facultative inward**
(NB FacIn, RNW Fac In, Endorsment Fac In), nol di 17 modul lain. Isinya tabel rate per peril
(EQ, terorisme, banjir, FLEXAS, RSMD, BI indemnity) dan kurs standar → **OQ-017**.

### 22.7 Stored procedure tanpa body — 67 kustom

| Skema | Jumlah | Yang menonjol |
| --- | ---: | --- |
| `POOLDATA` | 61 | `GET_TOKEN_STORAGE` (15 rule), `PROC_GENERATE_SEQUENCE_NUMBER` (14), `PEGA_JSON_KLAIM_PNC` (8), `GETCURRENCYSTANDARD` (7), 10× `RDBMASTER*` |
| `FIRE` | 2 | `PEGA_FIRE_SET_RATE`, `CEK_PRORATA_TANGGAL` — **logika rating di database** |
| `GENERAL` | 1 | `F_GET_NM_ASURADUR` |
| `GL` | 1 | `F_GET_EMAIL` (7 rule) |
| `MBU` | 1 | `F_CEK_HURUF` |
| `NEW_GENERAL` | 1 | `CEK_PLAT_NO` |

**Tidak satu pun body-nya ada di korpus** → **OQ-002**. Sebagian memuat perhitungan rating dan
penomoran — bagian paling material dari sistem.

### 22.8 OQ-011 — angka final korpus penuh

| Ukuran | Nilai |
| --- | ---: |
| Identitas yang muncul di **>1 modul** | **2.545** |
| — isi **identik** setelah normalisasi 18 tag | **2.012** |
| — isi **benar-benar berbeda** | **533** |
| File yang di-hash | 7.468 |

Sebaran versi isi: 520 identitas punya 2 versi, 8 punya 3, 1 punya 4, **4 punya 6 versi**.
**474 dari 533 konflik menyentuh lebih dari 2 modul.**

Register lengkap dengan **seluruh** path varian: **`_oq011-konflik-isi.md`**.
Tidak ada satu varian pun yang dipilih sebagai benar.

**Empat rule `When` dengan enam isi berbeda** — konflik terparah di korpus. `@BASECLASS / ISCLM`,
`ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` masing-masing muncul di 6 modul klaim dengan 6 isi yang
semuanya berlainan. `When` adalah guard percabangan, jadi ini berdampak langsung ke alur proses.

Satu rule `Flow` berkonflik: `ASM-FW-GISFW-WORK / OFFERFACRETRO` (NB FacIn = RNW Fac In,
Endorsment Fac In berbeda) → D2 harus menelusurnya dua kali.

**Catatan angka:** penjumlahan per batch menghasilkan 558, sedangkan pengelompokan korpus penuh
menghasilkan **533** — selisihnya karena identitas yang muncul di lebih dari satu batch tadinya
terhitung berulang. **533 yang berlaku.**

### 22.9 Guard otorisasi berbasis identitas orang (OQ-021)

| Ukuran | Nilai |
| --- | ---: |
| File memakai identitas operator sebagai literal | **66** (9 modul) |
| Varian identifier orang | 18 |
| File memakai `pyPosition` sebagai literal | **35** (14 modul) |

Guard identitas berada di `<pyStepsPreCondParamsWhen>` (193×), `<pyCondition>` (58×),
`<pyReadOnlyCondition>` (18×), `<pyContainerVisibleWhen>` (8×) — menggerbangi **eksekusi step,
keterubahan field, dan visibilitas layar**. Rule yang digerbangi termasuk
`ASM-FW-GISFW-WORK / COUNTGROSSPREMI_ACT` (perhitungan premi) dan
`ASM-FW-GCNMFW-WORK-KOMITE / KOMITEPOST_REJECT` (keputusan komite klaim).

Nilai `pyPosition` menyerupai **nama peran** (`ReasLifeAdmin`, `ReasLifeSPV`,
`ReasLifeMedicalAdvisor`, `IT Developer`, `Admin`, `SPV A`, `SPV B`) — titik awal paling konkret
untuk memetakan RBAC. Nilai nama orang **tidak dicatat di sini**; lihat OQ-021.

### 22.10 Batas kepercayaan artefak D1

Yang **terverifikasi** di D1: struktur, identitas, tipe, silsilah salinan, objek database yang
disentuh, procedure yang dipanggil, target integrasi, dan konflik isi antar modul.

Yang **belum** dan tidak boleh diperlakukan sebagai fakta:

- **Perilaku rule.** D1 tidak membaca logika. Kolom "Fakta dari tag" hanya metadata.
- **Arti kode/status.** `PaymentType`, `TransferType`, `EdmType`, `ProRateType`, `QR/QP/TP/TR`,
  `ISCLM`/`ISCLMP`/`ISCLMNP`/`ISPEGASYARIAH` — seluruhnya **belum terverifikasi**.
- **Tipe dan struktur data.** Tidak ada DDL; isi kolom JSON tidak diketahui.
- **Isi stored procedure dan tabel di ujung database link.**
- **Versi rule mana yang berjalan di production** untuk 533 identitas berkonflik.
- **Rule tertanam** (`<pyIncludedRuleXML>`) baru diperiksa untuk batch 2; batch 1, 3, 4 belum.
- **Tabel yang disentuh lewat `ReportDefinition` dan `DataPage`** belum diekstraksi — baru `RDBList`.
- **Dependency antar-Section** (`<pyInclude>`) baru dihitung jumlahnya, belum dipetakan.
