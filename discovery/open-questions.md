# Pertanyaan Terbuka — FASE A Discovery

Daftar hal yang **tidak dapat dipastikan dari korpus XML Pega**. Diisi terus sepanjang D0–D4.

Aturan: tidak ada pertanyaan di sini yang boleh ditutup dengan tebakan. Setiap entri wajib punya
**pemilik peran** dan keterangan **apa yang diblokir**. Penutupan hanya sah bila ada jawaban dari
pemilik atau bukti baru dari korpus (sebut path + rule).

Format entri:

```
### OQ-<nnn> — <judul singkat>
- **Pemilik:** DBA | Product+Underwriting | Actuarial | Finance | IAM
- **Ditemukan di:** STEP <Dx>, <path bukti>
- **Pertanyaan:** ...
- **Memblokir:** ...
- **Status:** terbuka | terjawab (<tanggal>, sumber)
```

Pemilik peran yang dipakai: **DBA**, **Product+Underwriting**, **Actuarial**, **Finance**, **IAM**.

---

## Terbuka

### OQ-001 — Tidak ada DDL Oracle maupun definisi properti di korpus

- **Pemilik:** DBA
- **Ditemukan di:** STEP D0, pemetaan struktur `D:\XML\RNM_BRD\` (17 folder tipe rule, tidak ada
  folder `Property/` dan tidak ada file DDL/SQL skema)
- **Pertanyaan:** Di mana definisi tabel Oracle existing (nama tabel, kolom, tipe, presisi,
  nullability, PK/FK, index, constraint) dapat diperoleh?
- **Memblokir:** setiap pernyataan tentang tipe data; keputusan representasi uang; pemetaan
  entitas; desain repository Go. Selama ini terbuka, D2 hanya boleh mencatat **nama tabel dan
  kolom apa adanya** dari `<pyBrowseSQL>` tanpa menyimpulkan tipenya.
- **Status:** terbuka — **TERTUTUP untuk Claim — Life** (2026-09-14, `[data DBA]`, DDL 12
  tabel/view). **Sisa untuk PremiumList Life + Endorsement Life** (2026-09-15): **kolom sudah
  terbaca, DDL fisik belum.**

  `[data DBA]` Body tiga procedure premium diserahkan 2026-09-15 (lihat **OQ-002**), sehingga
  **nama kolom** `M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS`, dan `JSON_OFFER_LIFE` kini diketahui —
  **cukup untuk menulis spec**.

  **Yang masih kurang:** DDL fisik ketiga tabel — **tipe, presisi, PK, index, nullability**.
  ⚠️ Khususnya: parameter procedure seluruhnya `VARCHAR2` termasuk kolom uang, sehingga **tipe
  kolom sebenarnya di tabel belum diketahui** — apakah `NUMBER` seperti pada
  `OS_AKSEPTASI_KLAIM_LIFE`, atau memang `VARCHAR2`.

  ### Tambahan 2026-09-15 — `EDMSTATUS` pada `M_LIFE_PREMIUM_DETAIL`

  `[terverifikasi]` **Nama kolomnya sudah terbaca dari korpus — tidak perlu ditebak.**
  `M_LIFE_PREMIUM_DETAIL` punya **tiga** kolom berakhiran status, dan hanya **satu** yang menandai
  peserta hidup/mati (bukti: daftar kolom `INSERT` + klausa `VALUES` di
  `Endorsement Life/RDBList/SaveMasterLPDet.xml` dan `PremiumList Life/RDBList/SaveMasterLPDet.xml`,
  `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL`):

  | Kolom | Diisi dari | Nilai | Penanda hidup/mati? |
  | --- | --- | --- | --- |
  | **`EDMSTATUS`** | `TempValue.EDMStatus` | `Old` / `New` / `Delete` / `Batal` | ✅ **ya** |
  | `STATUS` | `CARI48` | `0` untuk `QR`/`QP`, `1` untuk `TP`/`TR` | ❌ penanda **jenis transaksi** |
  | `STATUSOLD` | `CARI47` | `1`/`0` | ❌ |

  ⚠️ `[terverifikasi]` **Jalur new business tidak mengisi `EDMSTATUS` sama sekali** — sensus
  `PremiumList Life/Activity/InsertLifePremiumDetail_act.xml` (`ASM-FW-GISFW-WORK-LIFE` /
  `INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY`): **nol** kemunculan `EDMStatus`.

  **Yang perlu DBA:** **tipe dan nullability `EDMSTATUS`** — apakah baris new business menyimpan
  `NULL` atau **string kosong `''`**. Keduanya menuntut bentuk penyaring yang berbeda, dan salah
  pilih akan **membuang seluruh peserta new business** dari layar klaim (di Oracle, perbandingan
  dengan `NULL` tidak pernah bernilai benar).

  **Memblokir:** **tidak memblokir** aturan bisnisnya — hanya **bentuk akhir penyaring**. Sampai
  dipastikan, implementasi menangani **keduanya** (`IS NULL` *atau* `= ''`). Terkait **OQ-071**
  (tertutup) dan **AC 25–30 spec Claim — Life**.

  **Memblokir:** hanya **tiket migrasi** dan penetapan presisi target. **Tidak memblokir spec.**
  Pemilik: **DBA**.

  Cakupan modul lain di luar Life tetap terbuka sepenuhnya.
  DBA menyerahkan **DDL 12 tabel/view** yang mencakup **seluruh persistensi Claim — Life**.

### Daftar DDL yang diterima `[data DBA]`

| # | Objek | Catatan |
| --- | --- | --- |
| 1 | `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` | tabel akseptasi/outstanding klaim Life; **5 index** |
| 2 | `POOLDATA.M_LIFE_PREMIUM_DETAIL` | baris premi per peserta; banyak index; `PL_NUMBER`, `CERTIFICATE_NO`, `NAME_OF_INSURED`, `BEGIN_DATE`, `CURRENCY` |
| 3 | `POOLDATA.RATE_LIFE` | **VIEW** atas `M_RATE_LIFE.JSONDATA` |
| 4 | `POOLDATA.PRODUCTINWARD_LIFE` | **VIEW** atas `m_productinward_life.JSONDATA` |
| 5 | `POOLDATA.CURRENCY` | **VIEW** atas `m_currency.JSONDATA`, join `m_nation` |
| 6 | `POOLDATA.KODE_PRODUKSI` | `KODE VARCHAR2(5)`, `TYPE VARCHAR2(10)` |
| 7 | `POOLDATA.TANGGAL_CLOSING` | `TANGGAL VARCHAR2(10)` — **string, bukan `DATE`** |
| 8 | `POOLDATA.T_FOLDER_IMAGE` | `APPNAME VARCHAR2(20)` |
| 9 | `POOLDATA.MONITORING_KLAIM_LOG` | log service |
| 10 | `POOLDATA.EMAILKOMITE` | roster komite — rinci di **OQ-037** |
| 11 | `POOLDATA.GENERATE_SEQUENCE_NUMBER` | PK komposit `(CLASS, JENIS, TAHUN)`; `TAHUN VARCHAR2(5)`; `NO_SEQ NUMBER`; `MM_YYYY VARCHAR2(10)` |
| 12 | `POOLDATA.GCP_IMAGE` | cache token; `APPNAME`, `KODEAKSES`, `USERINPUT`, `INPUTDATE DATE DEFAULT sysdate` — semua `NOT NULL` |

### Tipe pada `OS_AKSEPTASI_KLAIM_LIFE` `[data DBA]`

| Kolom | Tipe | Catatan |
| --- | --- | --- |
| `STS_REJECT` | `NUMBER(38)` | ⚠️ **berbeda** dari `STS_REJECT VARCHAR2(15)` di `EMAILKOMITE` — nama sama, tipe berbeda, tabel berbeda |
| `CLAIM_AMOUNT`, `SUM_INSURED`, `SHARE_NUSANTARA_RE`, `SHARE_RETRO`, `CLAIM_RETRO`, `SUM_REASURED`, `CEDING_RETENTION`, `RETROCEDED_SHARE` | `NUMBER` **tanpa presisi** | delapan kolom uang **ADR-0003** |
| `CURRENCY` | `VARCHAR2(100)` | |
| `TYPE` | `VARCHAR2(10)` | kode `QP`/`QR`/`TP`/`TR` — **ADR-0012** |
| `LAYER_1` … `LAYER_4` | `VARCHAR2(10)` | ⚠️ **peran belum terverifikasi** — tidak muncul di jalur yang sudah ditelusur |

⚠️ **Konsekuensi mengikat:** kolom uang `NUMBER` **tanpa presisi** berarti Oracle menyimpan desimal
presisi arbitrer. Go **wajib** memakai tipe desimal presisi arbitrer; **`float64` dilarang**.
Diperkuat di **ADR-0003**.

### Temuan pola master JSON `[data DBA]`

Tiga master — `RATE_LIFE`, `PRODUCTINWARD_LIFE`, `CURRENCY` — **bukan tabel**, melainkan **view atas
kolom `JSONDATA`**. Data master yang sesungguhnya berbentuk **JSON**. Migrasi dan pembacaan master
harus sadar pola ini; membacanya sebagai tabel relasional biasa akan menyesatkan.

### Temuan: `AdjustmentList` tidak punya tabel fisik `[terverifikasi]`

**Tidak ada tabel tersendiri** untuk `AdjustmentList`. Bukti di
`Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`):
langkah `Property-Set` dengan `pyStepsObjectName = .AdjustmentList` dan
`pyStepsClassName = ASM-FW-GISFW-Data-AdjustmentLife` — kelas berawalan **`Data-`**, yaitu
**embedded, tanpa tabel** — dengan induk `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL`.

Kolom barisnya (`CLAIM_AMOUNT`, `CURRENCY`, `CURRENCYID`, `SUM_INSURED`, `SHARE_NUSANTARA_RE`,
`SHARE_RETRO`, `SUM_REASURED`, `CEDING_RETENTION`, `RETROCEDED_SHARE`, `STS_REJECT`) **identik**
dengan kolom `OS_AKSEPTASI_KLAIM_LIFE`.

**Artinya:** `AdjustmentList` hanyalah **page-list in-memory** sebelum disimpan; hasilnya
di-persist ke `OS_AKSEPTASI_KLAIM_LIFE`. Ini **tidak membatalkan ADR-0011** — unit *keputusan* tetap
baris adjustment; yang sekarang diketahui adalah **di mana baris itu mendarat**: sebagai baris pada
tabel akseptasi, bukan pada tabel adjustment tersendiri.

**Cakupan modul lain tetap terbuka** — DDL yang diserahkan hanya mencakup persistensi Claim — Life.

### OQ-002 — Body stored procedure/function Oracle tidak ada di korpus

- **Pemilik:** DBA
- **Ditemukan di:** STEP D0; diperluas di seluruh batch D1; angka final setelah sapuan 20 modul
- **Pertanyaan:** Di mana source body seluruh stored procedure/function yang dipanggil rule?
- **Memblokir:** pemahaman logika bisnis yang berada di sisi database; tidak bisa direkonstruksi
  dari XML.
- **Status:** terbuka — **TERTUTUP untuk Claim — Life** (2026-09-14) **dan untuk PremiumList Life +
  Endorsement Life** (2026-09-15). 66 stored procedure lain di korpus tetap terbuka.

### Penutupan untuk PremiumList Life + Endorsement Life (2026-09-15) `[data DBA]`

DBA menyerahkan **body tiga procedure premium**. Kolomnya kini terbaca; spec dapat ditulis tanpa
menebak.

**1. `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY`** — `INSERT` ke tabel fisik `M_LIFE_PREMIUM_SUMMARY`,
PK dari `M_LIFE_PREMIUM_SUMMARY_SEQ`. **37 kolom.**

| Kelompok | Kolom |
| --- | --- |
| Identitas | `ID`, `COB`, `PL_NUMBER`, `PL_NUMBER_EDM`, `CURRENCY`, `IDPEGA` |
| Uang gross | `PREMIUM`, `COMMISSION`, `BROKERAGE_FEE`, `OVR_COMM`, `TAX`, `PROF_COMM`, `CLAIM`, `CLAIM_AMOUNT`, `BALANCE`, `RI_ADMIN_FEE`, `DEDUCTION` |
| Turunan | kelompok `*_REFUND` dan `*_RETRO`, termasuk `*_REFUND_RETRO` |

⚠️ **Seluruh parameter bertipe `VARCHAR2` — termasuk kolom uang.** Ini **memperkuat ADR-0003**:
uang menyeberang batas sebagai **teks**, sehingga sistem baru wajib memakai desimal presisi
arbitrer dan **melakukan konversi eksplisit di satu batas** — bukan membiarkannya melewati `float`.

**Tidak commit sendiri** — `INSERT` + rollback saat error.

`[terverifikasi]` Pemanggilnya `PremiumList Life/RDBList/InsertPLSummary.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL`), mengirim
`pyWorkPage.BusinessName`, `PL_NUMBER`, `PL_NUMBER_EDM`, lalu `TempInputData.CARI2`…`CARI34`,
`pyWorkPage.pzInsKey`, dan `InputParam.ERRMSG out`.

**2. `POOLDATA.INSERTJSONPOLISLIFE`** — **upsert** ke `POOLDATA.JSON_POLIS` berkunci `IDPEGA`
(`UPDATE` bila ada, `INSERT` bila belum). Kolom: `IDPEGA`, `DATA_JSON` (**CLOB**), `TGL_INPUT`,
`NOPOLIS`, `PRODKE`, `TGL_PROD`, `USERNAME`. Fallback ke `JSON_POLIS_ERROR` bila gagal.
**COMMIT di dalam procedure.**

**3. `POOLDATA.INSERTJSONOFFERLIFE`** — **upsert** ke `POOLDATA.JSON_OFFER_LIFE` berkunci
(`IDPEGA`, `STATUS`), ID dari `JSON_OFFER_SEQ`. Kolom relasional: `BUSINESSCODE`/`NAME`,
`CEDINGCO`/`NAME`, `REINSURANCETYPE`, `POLICYHOLDER`/`NAME`, `TYPECEDING`,
`STATUS`/`STATUS_UPDATE`/`STATUS_FINAL`, `MARKETING_NOTE`, `OPERATORNAME`; **delapan tanggal**
(`RECEIVED`, `OFFERING`, `RESPONSE`, `CONFIRMATION`, `RECONFIRMATION`, `REALIZATION`,
`BINDING_DATE`, `MAX_TBC`); `OLDID`, `JSONDATA` (**CLOB**), `UW_NOTE`.
**COMMIT di dalam procedure.**

### ⚠️ Batas transaksi jalur Life adalah **CAMPURAN** `[terverifikasi + data DBA]`

| Procedure | Commit sendiri? |
| --- | --- |
| `PEGA_M_LIFE_PREMIUM_SUMMARY` | **tidak** |
| `PROC_GENERATE_SEQUENCE_NUMBER` | **tidak** |
| `INSERTJSONPOLISLIFE` | **ya** |
| `INSERTJSONOFFERLIFE` | **ya** |

**Konsekuensi desain mengikat:** Go harus **sadar batas ini**. Procedure yang commit sendiri
**memutus transaksi Go** di tengah jalan. **Jangan menempatkan data yang harus atomik pada dua sisi
procedure yang commit sendiri.** Terkait **OQ-013**.

Ini **berbeda dari Claim Life**, di mana kedua procedure jalur Life tidak commit sendiri sehingga
Go dapat memegang satu transaksi utuh.
  Menggantikan penyempitan sebelumnya (2 SQL penomoran lama tidak dipakai — tetap berlaku).

  `[data DBA]` **Format nomor klaim Life:**

  ```
  <prefix> + "K" + <kode bisnis L1/L2/…> + "." + MM.YYYY + "." + LPAD(seq,5)
  contoh:  RNML-KL1.08.2026.00936
  ```

  `[terverifikasi]` Prefix **di-LOOKUP, bukan hardcode**:
  `Claim Life/RDBList/GetKodeProdLife_SQL.xml` (`RULE-CONNECT-SQL`) berisi
  `SELECT KODE AS "ParamSeq.HASIL3" FROM POOLDATA.KODE_PRODUKSI WHERE TYPE ='LIFE'`.

  `[data DBA]` **Kontrak procedure:**
  `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER(p_class, p_jenis, p_proddate, OUT p_bulan, OUT p_seq_number)`

  | Aspek | Isi |
  | --- | --- |
  | Cakupan sequence | per **`(class, jenis, tahun)`** di tabel `GENERATE_SEQUENCE_NUMBER` |
  | Penguncian | `SELECT … FOR UPDATE` |
  | `p_jenis` | membedakan **retro / non-retro** |
  | Pengguliran periode | lewat `POOLDATA.TANGGAL_CLOSING` |
  | Aturan cutover | `TRUNC(now) <= 02/01/2026` → periode **`12.2025`** |
  | Bentuk keluaran | `p_seq_number = LPAD(seq, 5, '0')` |

  Ini menjawab pertanyaan lama "sequence di-reset per tahun / lini / global": **per `(class, jenis,
  tahun)`** — jadi per tahun **dan** per jenis.

  `[data DBA]` **Token penyimpanan:** `POOLDATA.GET_TOKEN_STORAGE` menghasilkan token MD5 lewat
  `STANDARD_HASH`, di-cache di tabel `GCP_IMAGE`, dengan **masa berlaku 1 menit**.

  Batas transaksi kedua procedure: lihat **OQ-013** (tertutup untuk Claim — Life).

  **66 stored procedure lain di korpus tetap terbuka.**

**ANGKA FINAL KORPUS (20/20 modul, STEP D1 selesai):**

**67 stored procedure/function kustom tanpa body**, tersebar di **7 skema**. (Ditambah 3 built-in
Oracle yang tidak perlu dimigrasi: `DBMS_LOB.CREATETEMPORARY` 41 rule,
`UTL_MATCH.EDIT_DISTANCE_SIMILARITY` 2, `DBMS_OUTPUT.PUT_LINE` 1.)

| Skema | Procedure kustom |
| --- | ---: |
| `POOLDATA` | 61 |
| `FIRE` | 2 |
| `GENERAL` | 1 |
| `GL` | 1 |
| `MBU` | 1 |
| `NEW_GENERAL` | 1 |
| **Total** | **67** |

Daftar lengkap 67 procedure beserta jumlah rule pemanggilnya ada di `inventory/_summary.md` §22.

> **Koreksi metode (batch 4).** Angka per-batch yang dilaporkan sebelumnya (13 / 50 / 19) **terlalu
> tinggi**: detektor procedure ikut menangkap `INSERT INTO skema.tabel (kolom…)` dan referensi page
> Pega `{Page.property(...)}` sebagai "procedure". Detektor sudah diperbaiki — nama yang sudah
> teridentifikasi sebagai tabel dan isi `{...}` kini dikecualikan — lalu **seluruh 1.146 rule
> ber-SQL di korpus diekstraksi ulang**. Angka 67 di atas yang berlaku.

Delapan procedure baru dari batch 3, seluruhnya seputar **akseptasi/outstanding klaim**:
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`, `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT`,
`POOLDATA.PEGA_JSON_OS_AKSEP_SUBJECTIVITY`, `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`,
`POOLDATA.XOL2_AKSEP_KLAIM`, `POOLDATA.PEGA_M_CAUSE_OF_LOSS`,
`POOLDATA.GENERATE_NOCLMTREATYIN`, `POOLDATA.GENERATE_NOCLMTRTYINTEMP`.

Akhiran `TNP`, `TRT`, `PNC`, dan `XOL2` **kepanjangan belum terverifikasi** — jangan ditebak.

Yang paling material — procedure di luar `POOLDATA`, karena sebagian memuat **logika rating dan
perhitungan**, bukan sekadar CRUD:

| Procedure | Dugaan peran (belum terverifikasi) |
| --- | --- |
| `FIRE.PEGA_FIRE_SET_RATE` | penetapan rate lini Fire |
| `FIRE.CEK_PRORATA_TANGGAL` | perhitungan pro-rata tanggal |
| `GENERAL.F_GET_NM_ASURADUR` | pencarian nama asuradur |
| `GL.F_GET_EMAIL` | pencarian email dari skema `GL` |
| `MBU.F_CEK_HURUF` | arti belum terverifikasi |
| `NEW_GENERAL.CEK_PLAT_NO` | validasi nomor plat kendaraan |

Ditambah `POOLDATA.RDBMASTER*` (10 procedure master data) dan keluarga `POOLDATA.FACIN*`.
Daftar lengkap di `inventory/_summary.md` §10.3 dan di §6 tiap file modul.

**DIPERSEMPIT untuk konteks Claim — Life — 2026-09-14, sumber: work owner** (grilling Ronde 1 Q7,
`.scratch/claim-life/grilling-ronde-1.md`).

`[terverifikasi work owner]` Dua rule penomoran klaim life **sudah tidak dipakai (di-remark)**:
`Generate_NoKlaim_Life` dan `Generate_NoKlaim_LifeRetro`. Penomoran sekarang lewat
`GetSequenceNumber_SQL` → `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.
**Keduanya jangan dimigrasikan.**

`[terverifikasi]` **Dikuatkan bukti korpus** — asimetri indeks rujukan rule di
`Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`):

| RequestType | `<RequestType>` | `pxRuleReferences` (`<pyRuleName>`) |
| --- | ---: | ---: |
| `Generate_NoKlaim_Life` | 1 | **0** |
| `Generate_NoKlaim_LifeRetro` | 1 | **0** |
| `GetSequenceNumber_SQL` | 1 | **2** |

Keduanya tertinggal sebagai parameter langkah tetapi **tidak terindeks sebagai rujukan aktif**.

```
f="Claim Life/Activity/SaveOutStandingLife_Act.xml"
grep -c '<pyRuleName>.*Generate_NoKlaim_Life<' "$f"        # 0
grep -c '<pyRuleName>.*GetSequenceNumber_SQL<' "$f"        # 2
```

**Koreksi artefak:** `flows/Claim Life.md` §3 mendaftar kedua rule itu sebagai dipanggil — itu
berdasarkan kehadiran `<RequestType>` saja. Status sesungguhnya: **tidak aktif**.

**YANG MASIH TERBUKA untuk Claim — Life:** kontrak **`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`** —
format nomor klaim yang dihasilkan, dan apakah sequence di-reset per tahun / per lini / global.
**Pemilik: DBA.** Sampai dijawab, spec menyebut penomoran sebagai **dependensi database**, bukan
mereplikasinya.

**Cakupan korpus-wide TETAP TERBUKA:** 66 stored procedure lain di 10 skema aplikasi tidak
tersentuh jawaban ini.


### OQ-003 — Folder `excludeXML` — apa artinya dan apakah isinya aktif

- **Pemilik:** Product+Underwriting (atau pemilik export Pega)
- **Ditemukan di:** STEP D0, `Komite Claim Non Prop/excludeXML/GetBase64Attachment.xml` (1 file,
  satu-satunya di korpus)
- **Pertanyaan:** Apakah folder `excludeXML` berarti rule ini sengaja dikecualikan dari export /
  tidak dipakai di production, atau sekadar konvensi penamaan folder saat ekspor?
- **Memblokir:** apakah rule ini masuk ruang lingkup migrasi atau tidak.
- **Status:** terbuka — **pertanyaannya jadi lebih tajam setelah batch 3**

**Diperbarui STEP D1 batch 3 `[terverifikasi]`:** isi file itu bukan placeholder. Ia rule
**`RULE-OBJ-ACTIVITY`** beridentitas `ASM-FW-GCNMFW-WORK / GETBASE64ATTACHMENT`, hasil salinan dari
`WORK- / LOADATTACHMENTDATA`. Jadi ini rule Activity lengkap yang sengaja ditempatkan di folder
bernama "exclude". Pertanyaan aktif/tidaknya tetap terbuka.

### OQ-005 — LIMA modul tanpa rule `Flow`: apa titik masuk prosesnya

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D0 — modul tanpa file di `Flow/`
- **Pertanyaan:** Apakah modul-modul ini memang tidak berbasis workflow (mis. hanya layar master
  data / screen flow), ataukah Flow-nya tinggal di modul lain atau tidak ikut terekspor?
- **Memblokir:** metode telusur STEP D2 untuk modul tersebut; penentuan apakah ada proses
  approval/assignment di dalamnya.
- **Status:** **terjawab untuk keperluan D2** (2026-09-13, STEP D2 Tahap 5) — titik masuk ketetapan ada di `flows/_METHOD-noflow.md`; menyisakan pertanyaan sejarah (mengapa tanpa `Flow`), yang tidak memblokir D3/D4.

**Dikoreksi STEP D1 batch 4 `[terverifikasi]`.** Daftar asli menyebut **enam** modul. `PremiumList
Life` **keliru dimasukkan**: modul itu punya rule `Flow`
(`ASM-FW-GISFW-WORK-LIFE / INPUTPOLICYHOLDER`, `<pxObjClass>Rule-Obj-Flow`,
`<pyStartActivity>Start1`), hanya filenya berada di root modul, bukan di folder `Flow/` — lihat
OQ-004 (terjawab).

Daftar yang benar — **5 modul tanpa rule `Flow`**:

| Modul | File | Kandidat titik masuk D2 |
| --- | ---: | --- |
| Treaty In | 329 | `Harness` (3), `FlowAction` (32) |
| Treaty In Adjustment | 379 | `Harness` (6), `FlowAction` (43) |
| Treaty Contract Out | 303 | `Harness` (3), `FlowAction` (2) — juga terkait **OQ-022** |
| Master Product Name Life | 114 | `Harness` (1), `FlowAction` (11) |
| Master Contract Retro Life | 66 | `Harness` (4), `FlowAction` (2) |

Jumlah rule `Flow` di korpus juga terkoreksi: **24**, bukan 23.

**TERJAWAB untuk keperluan D2 — STEP D2 Tahap 5 (2026-09-13) `[terverifikasi]`.**

Metode titik masuk untuk modul tanpa `Flow` sudah **ditetapkan dan dipakai**, bukan diimprovisasi:
**`flows/_METHOD-noflow.md`**. Titik masuk kelima modul:

| Modul | Harness titik masuk | `pxInsName` |
| --- | --- | --- |
| Treaty In | `Harness/InputTreatyInOffer.xml` (8.225.095 B) | `DATA-PORTAL!INPUTTREATYINOFFER` |
| Treaty In Adjustment | `Harness/InputTreatyInAdjustment.xml` | `DATA-PORTAL!INPUTTREATYINADJUSTMENT` |
| Treaty Contract Out | `Harness/InboxTreatyContract.xml` | `DATA-PORTAL!INBOXTREATYCONTRACT` |
| Master Product Name Life | `Harness/InwardProductName.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE!INWARDPRODUCTNAME` |
| Master Contract Retro Life | `Harness/InboxRetroLifeReinsurersList.xml` | `DATA-PORTAL!INBOXRETROLIFEREINSURERSLIST` |

**Jawaban atas pertanyaan aslinya `[terverifikasi]`: keduanya benar, tergantung modul.**

- **`Treaty In` dan `Treaty In Adjustment` TETAP berbasis proses** — ada properti status
  (`StatusAkseptasi`, 4 nilai), tangga persetujuan 4 tingkat, dan tombol Submit/Decline.
  Prosesnya hanya **tidak ditulis sebagai graf `Flow`**, melainkan sebagai satu Data Transform
  bersarang (`Akseptasi_DT`: 10 `WHEN`, 18 `OTHERWISE_WHEN`, 64 `SET`).
- **`Treaty Contract Out`, `Master Product Name Life`, `Master Contract Retro Life` memang layar
  master murni** — nol rujukan `StatusAkseptasi`, tanpa `Akseptasi_DT`, tanpa tombol
  Submit/Decline; polanya `New…` / `Set…` / `Save…` / `Delete…` / `CancelActivity…`.

**Menyisakan:** apakah `Flow` untuk kedua modul treaty memang tidak pernah ada, atau pernah ada dan
digantikan — pertanyaan sejarah, bukan isi korpus. Tidak memblokir D3/D4.


### OQ-007 — Tidak ada rule identitas/otorisasi di korpus

- **Pemilik:** IAM
- **Ditemukan di:** STEP D0 — tidak ada folder tipe rule terkait access group, role, atau operator
  dalam 17 tipe yang ada
- **Pertanyaan:** Bagaimana peran, hak akses, dan segregation of duties didefinisikan di Pega
  existing, dan di mana artefaknya?
- **Memblokir:** pemahaman routing assignment/workbasket di D2; seluruh bahasan RBAC di FASE B.
- **Status:** terbuka

### OQ-008 — Arti prefix `ASM` dan `RNM` pada kunci rule RDBList

- **Pemilik:** DBA
- **Ditemukan di:** STEP D1 batch 1 — seluruh 161 rule `RDBList` memakai kunci tiga bagian
  `CLASS!PREFIX!NAMA`; prefix `ASM` pada 89 rule dan `RNM` pada 72 rule. Contoh:
  `NB Treaty In/RDBList/BrowseTreatyIn.xml` → `ASM-FW-GISFW-INT-TREATY_IN!ASM!BROWSETREATYIN`;
  `NB Treaty In/RDBList/BrowseTreatyOut.xml` → `...!RNM!BROWSETREATYOUT`
- **Pertanyaan:** Apakah `ASM`, `RNM`, dan `GCNM` menunjuk koneksi/database yang berbeda? Bila ya,
  mana yang mana, dan apakah semuanya masih aktif?
- **Memblokir:** desain lapisan repository Go (satu atau beberapa sumber data); pemahaman apakah
  ada data yang tersebar di beberapa database.
- **Status:** terbuka — **prefix ketiga ditemukan di batch 2**

**Diperbarui STEP D1 batch 2:** muncul prefix **`GCNM`** yang tidak ada di batch 1, khusus di modul
claim facultative. Sebaran per modul:

| Modul | `ASM` | `RNM` | `GCNM` |
| --- | ---: | ---: | ---: |
| Komite Claim FacIn | 2 | 13 | 2 |
| Claim Fac In | 13 | 24 | 27 |
| RNW Fac In | 170 | 30 | — |
| Endorsment Fac In | 163 | 25 | — |
| NB FacIn | 178 | 39 | — |

`GCNM` sejalan dengan keluarga class **`ASM-FW-GCNMFW-*`**, yang tidak muncul sama sekali di
batch 1 dan di batch 2 terpusat pada modul claim facultative `[terverifikasi]`:

| Modul | `ASM-FW-GISFW-*` | `ASM-FW-GCNMFW-*` |
| --- | ---: | ---: |
| Claim Fac In | 112 | **302** |
| Komite Claim FacIn | 42 | **43** |
| Endorsment Fac In | 1.744 | 2 |
| NB FacIn | 1.790 | — |
| RNW Fac In | 1.657 | — |
| seluruh modul batch 1 | 769 | — |

`[dugaan]` ada dua framework Pega berbeda, dan domain claim facultative dibangun di atas
`GCNMFW` sementara sisanya di atas `GISFW`. **Belum terverifikasi.** Kepanjangan `GISFW` dan
`GCNMFW` **belum terverifikasi** — jangan ditebak.

Dua rule `GCNMFW` yang menyeberang ke `Endorsment Fac In` perlu diperiksa di D2: apakah itu
ketergantungan lintas framework yang nyata atau sisa ekspor.

**Diperbarui STEP D1 batch 3 — dugaan di atas perlu diperluas.** `GCNMFW` bukan khusus claim
*facultative*; ia mendominasi **seluruh domain klaim** `[terverifikasi]`:

| Modul | `ASM-FW-GISFW-*` | `ASM-FW-GCNMFW-*` |
| --- | ---: | ---: |
| Claim Prop | 53 | **190** |
| Claim Non Prop | 81 | **169** |
| Claim Life | 74 | 55 |
| Komite Claim Prop | 24 | **47** |
| Komite Claim Non Prop | 15 | **37** |
| Komite Claim Life | 25 | 17 |

Prefix `GCNM` juga muncul pada 51 kunci RDBList batch 3 (`ASM`=47, `RNM`=88).

`[dugaan]` `GCNMFW` = framework klaim, `GISFW` = framework underwriting/produksi; modul klaim
memakai keduanya karena membaca data produksi. **Belum terverifikasi.**

### OQ-009 — Rule yang tinggal di class bawaan Pega, bukan class aplikasi

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D1 batch 1 — dari 65 class unik, sebagian bukan class aplikasi `ASM-FW-*`:
  `DATA-PORTAL` (341 rule di batch 1, mayoritas Section milik Treaty In / Treaty In Adjustment),
  `@BASECLASS` (23), `ASSIGN-WORKLIST` (4), `PEGACRM-PORTAL` (2),
  `PEGACRM-WORK-` (1, `NB Treaty In/When/crmCreateOpportunity.xml`),
  `LINK-ATTACHMENT`, `DATA-WORKATTACH-FILE`, `WORK-`, `DATA-PARTY`
- **Pertanyaan:** Mana di antara rule ini yang merupakan kustomisasi Nusantara Re (wajib
  dimigrasikan) dan mana yang bawaan produk Pega (tidak perlu dimigrasikan, cukup diganti
  padanannya)? Khususnya: mengapa ratusan Section kustom ditempatkan di `Data-Portal`?
- **Memblokir:** penentuan ruang lingkup migrasi yang sebenarnya — berpotensi memperkecil atau
  memperbesar cakupan secara signifikan.
- **Status:** terbuka

**Diperbarui STEP D2 Tahap 5 `[terverifikasi]` — proporsi tertinggi ditemukan.**

`Treaty Contract Out`: **202 dari 303 rule (66,7 %) berada di `@BASECLASS`** — proporsi tertinggi
di korpus. Perintah audit:
```
awk -F'\t' '$1 ~ "^Treaty Contract Out/"{print $3}' all-rules.tsv | sort | uniq -c | sort -rn
```
→ `@BASECLASS` 202, `ASM-FW-GISFW-INT-PROPORTIONALARRG` 40, `DATA-PORTAL` 8,
`ASM-FW-GISFW-INT-T_STORAGE_IMAGE` 7, `ASM-FW-GISFW-INT-TREATY_IN` 5, dan seterusnya.

Akibatnya konkret: **pemetaan rule → entitas domain tidak dapat diturunkan dari class** untuk dua
pertiga modul itu. Pola serupa di `Master Contract Retro Life`, yang activity perhitungannya
(`CountingPercentShare_Act`, `TreatyLimit_TypeProtect`, `SetValueRetroLimit_TreatyYearLife`)
seluruhnya di `@BASECLASS`.


### OQ-010 — `Treaty In` dan `Treaty In Adjustment`: dua modul atau satu

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D1 batch 1 — keduanya berbagi **320 identitas rule**
  (`Treaty In` 329 file, `Treaty In Adjustment` 379 file). Keduanya juga sama-sama tanpa rule `Flow`
  dan sama-sama menaruh Section di class `DATA-PORTAL`
- **Pertanyaan:** Apakah keduanya satu aplikasi yang diekspor dua kali dengan cakupan sedikit
  berbeda, atau memang dua modul terpisah yang berbagi basis kode?
- **Memblokir:** penentuan bounded context di D4; menghindari memigrasi pekerjaan yang sama dua kali.
- **Status:** terbuka — sebagian dapat diuji sendiri di D2.

**Diperbarui STEP D2 Tahap 5 (2026-09-13) — bukti pengukuran lengkap `[terverifikasi]`.**

| Ukuran | Hasil |
| --- | ---: |
| File bernama sama | **323** |
| — **identik** (normalisasi 18 tag) | **277** (85,8 %) |
| — berbeda | 46 |
| — berbeda setelah normalisasi **21 tag** | **43** |
| — **konflik semu** (3 tag metadata) | **3** — keluarga `T_STORAGE_IMAGE` |
| Hanya di `Treaty In` | 6 |
| Hanya di `Treaty In Adjustment` | **56** |

**Berbagi mesin status `[terverifikasi]`:** `DataTransform/Akseptasi_DT.xml` **identik**
(hash `58b8e650`), demikian pula `Activity/Akseptasi_Act.xml` (`027a42f4`),
`TreatyInAkseptasi_Act.xml` (`226e2ca4`), `Section/TreatyInActionButtons.xml`,
`SaveTreatyIn_Act.xml`, `TreatyInDeclineConfirmation_postact.xml`.

**Yang khas `Treaty In Adjustment` `[terverifikasi]`:** 23 layar *data lama*
(9 `FlowAction/*OldData*` + 14 `Section/*OldData*`), mesin penomoran revisi
(`TreatyInRevisi_post`, `GetTreatyRevisionID`, `TreatyInSetAddendumToHistory`,
`TreatyRevisionCopyAttachment`), `DataTransform/TreatyCreateEDM.xml`,
`PickerTreatyInMaster(Revisi)`, dan satu-satunya `RULE-OBJ-MENU` di kelompok ini.

**Pembedanya terbaca `[terverifikasi]`:** `Param.type` bernilai `"revision"` / `"adjustment"` →
`TreatyIn.EDMState` `"1"` / `"3"`; ditambah `TreatyIn.RevisionState` yang **memendekkan tangga
persetujuan** (dari Sec Head langsung `Resolve Complete`, melewati Dept Head dan Direktur).

**Untuk 43 rule yang tetap berbeda, struktur langkah dan Section yang di-include justru identik**
(diperiksa pada `Activity/TreatyInSubmit.xml` dan `Harness/InputTreatyInOffer.xml`). Selisihnya
ada pada tingkat yang tidak terbaca dari tag struktural — **tidak dikarakterisasi**, batas cakupan
telusur.

**Kesimpulan berbukti:** satu basis rule, dibedakan **saat runtime** oleh nilai properti — pola
yang **sama dengan OQ-015** (facultative NB/RNW/EDM). **Penetapan bounded context tetap D4.**


### OQ-011 — 533 identitas rule yang isinya berbeda antar modul: mana yang berlaku di production

- **Pemilik:** Product+Underwriting (+ pemilik export Pega)
- **Ditemukan di:** STEP D1 batch 1 — 387 identitas muncul di >1 modul; setelah normalisasi
  (buang tag timestamp, urutkan baris, hash) **257 identik** tetapi **130 benar-benar berbeda
  isinya**. Contoh: `ASM-FW-GISFW-INT-TREATY_IN / AKSEPTASI_ACT` berbeda antara
  `Treaty In/Activity/Akseptasi_Act.xml` dan `Treaty In Adjustment/Activity/Akseptasi_Act.xml`
- **Pertanyaan:** Untuk tiap identitas yang berbeda, versi mana yang benar-benar berjalan di
  production? Apakah ekspor diambil dari ruleset/versi yang berbeda?
- **Memblokir:** seluruh telusur D2 atas rule tersebut — membaca salinan yang salah akan
  menghasilkan pemahaman perilaku yang keliru. **Ini risiko halusinasi paling konkret di korpus.**
- **Status:** terbuka — angka dikoreksi turun di batch 3 (lihat di bawah)

**ANGKA FINAL KORPUS (20/20 modul, STEP D1 selesai).**
Register lengkap: **`inventory/_oq011-konflik-isi.md`** (533 entri).

| Ukuran | Nilai |
| --- | ---: |
| Identitas rule unik di korpus | 4.444 |
| Identitas yang muncul di **>1 modul** | **2.545** |
| — isi **identik** setelah normalisasi 18 tag | **2.012** |
| — isi **benar-benar berbeda** | **533** |
| File yang di-hash | 7.468 |

Sebaran versi isi: 520 identitas punya 2 versi, 8 punya 3, 1 punya 4, dan **4 punya 6 versi**.
**474 dari 533 konflik menyentuh lebih dari 2 modul.**

**Konflik paling parah** `[terverifikasi]` — empat rule `When` di `@BASECLASS` yang masing-masing
muncul di **6 modul dengan 6 isi berbeda**, tidak ada dua modul pun yang sepakat: `ISCLM`,
`ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` (Komite Claim FacIn, Claim Fac In, Komite Claim Non Prop,
Komite Claim Prop, Claim Prop, Claim Non Prop). Karena `When` berperan sebagai **guard
percabangan**, perbedaan ini berdampak langsung ke alur proses klaim. Arti dan kepanjangan keempat
nama itu **belum terverifikasi**.

Satu rule **`Flow`** juga berkonflik: `ASM-FW-GISFW-WORK / OFFERFACRETRO` — `NB FacIn` dan
`RNW Fac In` identik, `Endorsment Fac In` berbeda. D2 harus menelusur proses itu **dua kali**.

**Riwayat koreksi angka** (jejaknya dipertahankan agar metodenya dapat diaudit):

| Tahap | Angka | Sebab |
| --- | ---: | --- |
| Batch 1–3, normalisasi 7 tag | 976 | 11 tag provenance ekspor ikut terhitung sebagai isi |
| Batch 1–3, normalisasi 18 tag | 543 | 433 positif palsu dibuang |
| + batch 4, dijumlahkan per batch | 558 | — |
| **Korpus penuh, dikelompokkan lintas modul** | **533** | identitas lintas batch tadinya terhitung berulang |

Daftar 18 tag volatil dan cara auditnya ada di register. Daftar itu mungkin **masih belum
lengkap** — bila ditemukan tag provenance lain, angka 533 dapat turun lagi.

**Cara batch 2 memperlakukan ini (dan cara batch 3–4 harus memperlakukannya):** setiap identitas
berbeda-isi dicatat sebagai **konflik bernomor beserta seluruh path variannya**. Tidak ada satu
varian pun yang dipilih, ditandai "utama", atau dijadikan dasar ringkasan peran rule. Perbandingan
**wajib** memakai hash ternormalisasi; `md5sum` mentah dilarang menjadi dasar klaim "berbeda".

Contoh konflik 3 varian dari batch 2:
`ASM-FW-GISFW-INT-DOCUMENT_POLIS / GETALLDOCUMENT` berbeda di
`RNW Fac In/ReportDefinition/`, `Endorsment Fac In/ReportDefinition/`, dan `NB FacIn/ReportDefinition/`.

**Diperbarui STEP D2 Tahap 4 (2026-09-13) — pengukuran ulang untuk domain facultative
`[terverifikasi]`.** **438 dari 533 identitas (82,2%)** menyentuh `NB FacIn`, `RNW Fac In`, **dan**
`Endorsment Fac In` sekaligus. Telusur Tahap 4 menemukan tiga tag metadata Pega yang **luput** dari
daftar 18 tag normalisasi — `<pyDelete>`, `<pyVersionSecure>`, `<pzIsPrivateCheckOut>` — lalu
menghitung ulang keseluruhan 438 identitas dengan normalisasi **21 tag**:

| Ukuran | Hasil |
| --- | ---: |
| Masih >1 isi — normalisasi 18 tag | **438** |
| Masih >1 isi — normalisasi **21 tag** | **429** |
| **Konflik semu** (runtuh jadi identik) | **9** — 2,1% |

Kesembilan konflik semu: #12 `ASM!GETDATAKLAIM_SQL`, #13 `ASM!GETEDMOLDIDPEGA`,
**#324 `LINKSERVICE!LINKSERVICE`**, #367 `RNM!DELETESTORAGE_SQL`, #369 `RNM!GETAPPNAME_SQL`,
#371 `RNM!GETLINKSTORAGE_SQL`, #372 `RNM!GETMKTANDLEADER_SQL`, #376 `RNM!INSERT_T_STORAGE_SQL`,
#377 `RNM!UPDATE_T_STORAGE_SQL`.

**Register tetap berlaku: 97,9% konflik facultative adalah perbedaan isi yang nyata.** Dua yang
berubah secara material: **#324 `LINKSERVICE`** (konfigurasi base URL ternyata identik di ketiga
modul — lihat OQ-047) dan **#415 `SERVICEINSERTARASAPAS_ACT`** (lihat OQ-025).

Entri **#338 `OFFERFACRETRO`** dikonfirmasi **perbedaan perilaku nyata**: varian NB/RNW punya
tangga persetujuan fac out **2 anak tangga**, varian Endorsment punya **4**. Bukti:
`flows/_SUMMARY-facultative.md` §4.


**Koreksi STEP D2 Tahap 5 (2026-09-13) `[terverifikasi]` — contoh di narasi entri ini USANG.**

Narasi di atas mencontohkan `ASM-FW-GISFW-INT-TREATY_IN / AKSEPTASI_ACT` sebagai rule yang
**berbeda** antara `Treaty In` dan `Treaty In Adjustment`. Pengukuran ulang menunjukkan keduanya
**identik** (hash ternormalisasi 18 tag = `027a42f4`), dan identitas itu **tidak ada di register**
`inventory/_oq011-konflik-isi.md`.

Contoh tersebut berasal dari pengukuran batch 1 dengan normalisasi **7 tag**, sebelum perbaikan ke
18 tag (yang memangkas konflik 976 → 543). **Register tetap berlaku; contohnya yang perlu
diganti.** Contoh pengganti yang terverifikasi: `Harness/InputTreatyInOffer.xml` dan
`Activity/TreatyInSubmit.xml` — keduanya tetap berbeda bahkan setelah normalisasi 21 tag.

**Pengukuran tambahan Tahap 5:** pola **konflik semu** yang ditemukan di Tahap 4 muncul lagi dan
konsisten — keluarga rule penyimpanan berkas `T_STORAGE_IMAGE` / `LINKSERVICE`:
3 dari 46 di `Treaty In` ↔ `Treaty In Adjustment`, dan **7 dari 19** di
`Master Product Name Life` ↔ `Treaty In` (termasuk `SystemSettings/LinkService.xml`).


### OQ-012 — Struktur JSON di dalam kolom `JSONDATA`

- **Pemilik:** DBA + Product+Underwriting
- **Ditemukan di:** STEP D1 batch 1 — `NB Treaty In/RDBList/BrowseTreatyIn.xml`
  (`ASM-FW-GISFW-INT-TREATY_IN / ASM!BROWSETREATYIN`) menyeleksi kolom `JSONDATA` dari
  `pooldata.M_TREATY_IN` dan `pooldata.M_TREATY_IN_edm`
- **Pertanyaan:** Apa struktur/skema JSON yang disimpan di kolom `JSONDATA`? Adakah kontrak yang
  mengaturnya, atau bentuknya bergantung pada rule yang menulis?
- **Memblokir:** model data domain tidak dapat diturunkan dari SQL saja; setiap pemetaan entitas
  akan menjadi tebakan selama ini terbuka. Berkaitan erat dengan OQ-001.
- **Status:** terbuka — **meluas di batch 2**

**Diperbarui STEP D1 batch 2:** pola JSON opaque bukan kekhususan treaty inward. **115 rule batch 2**
menyentuh objek/kolom berlabel JSON, dan **`JSON_POLIS` adalah objek yang paling sering dirujuk di
seluruh batch 2 (58 rule)**. Objek lain: `JSON_OFFER`, `POOLDATA.JSON_POLIS`, serta procedure
`POOLDATA.INSERTJSONPOLIS`, `POOLDATA.PEGA_M_JSON_OFFER`, `POOLDATA.PEGA_JSON_KLAIM_PNC`,
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`, `POOLDATA.INSERTJSONPOLISMONITORING`.

Pertanyaan tambahan: apakah struktur JSON berbeda antar lini (treaty vs facultative vs claim),
atau satu kontrak bersama?

**Diperbarui STEP D1 batch 3:** pola yang sama muncul di domain klaim —
`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` (7 rule), `OS_AKSEPTASI_KLAIM` (6), `JSON_POLIS` (5),
`JSON_KLAIM` (3), `HISTORYAKSEPTASIPEGA` (2), ditambah procedure `POOLDATA.PEGA_JSON_KLAIM_PNC`,
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`, `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT`,
`POOLDATA.PEGA_JSON_OS_AKSEP_SUBJECTIVITY`. Struktur JSON di dalamnya tetap tidak ada di korpus.

**Diperbarui STEP D2 Tahap 4 `[terverifikasi]`.** Pola JSON opaque **memuncak** di domain
facultative: `JSON_POLIS` adalah objek Oracle **paling banyak dirujuk** di ketiga modul —
NB FacIn **18** rule, RNW Fac In **18**, Endorsment Fac In **17**; ditambah `POOLDATA.JSON_POLIS`
(4–5 rule) dan procedure `POOLDATA.INSERTJSONPOLIS`, `POOLDATA.INSERTJSONPOLISMONITORING`,
`POOLDATA.PEGA_M_JSON_OFFER`.

Sumber isinya kini teridentifikasi lebih jelas: fungsi Pega kustom **`@ASM.GetPageJSONString()`**
dipakai di `Endorsment Fac In/Activity/serviceInsertArasapas_act.xml` untuk membentuk JSON yang
dikirim keluar lewat `Connect-REST` — **source fungsi itu tidak ada di korpus** (sama seperti
temuan `flows/NB Treaty In.md` §5.3). Selama fungsi itu tidak tersedia, **struktur JSON tidak dapat
diturunkan dari korpus sama sekali**, bukan sekadar tidak terdokumentasi.

Perintah audit:
```
awk -F'\t' '$1 ~ "^NB FacIn/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
 | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn | head
grep -o "GetPageJSONString" "Endorsment Fac In/Activity/serviceInsertArasapas_act.xml" | wc -l
```


### OQ-013 — Batas transaksi dan penggunaan identitas pengguna berada di dalam stored procedure

- **Pemilik:** DBA
- **Ditemukan di:** STEP D1 batch 1 — `NB Treaty In/RDBList/SavePolisTreatyIn_SQL.xml`
  (`ASM-FW-GISFW-INT-POLISTREATYIN / ASM!SAVEPOLISTREATYIN_SQL`) memanggil
  `POOLDATA.PEGA_JSON_POLIS_TREATYIN(...)` dan melakukan `COMMIT` **di dalam blok PL/SQL**,
  serta mengirim `{OperatorID.pyUserIdentifier}` sebagai parameter
- **Pertanyaan:** Apakah seluruh prosedur `POOLDATA.*` melakukan commit sendiri? Untuk apa
  identitas pengguna dipakai di dalam prosedur — audit trail, otorisasi, atau keduanya?
- **Memblokir:** desain pengelolaan transaksi di backend Go; desain audit trail; pemahaman apakah
  otorisasi sebagian ditegakkan di database.
- **Status:** terbuka — **TERTUTUP untuk Claim — Life** (2026-09-14, sumber: **data DBA + korpus +
  keputusan work owner**).

  `[data DBA]` **Stored procedure jalur Life tidak memuat `COMMIT` sendiri.** DBA memeriksa badan
  `PROC_GENERATE_SEQUENCE_NUMBER` dan `GET_TOKEN_STORAGE`: keduanya bersih dari `COMMIT`.

  **Koreksi atas catatan lama.** Artefak sebelumnya menyebut "`COMMIT` berada di dalam blok PL/SQL".
  Itu benar tetapi **salah tempat**: `[terverifikasi]` `COMMIT` ada di **pembungkus Pega**, bukan di
  procedure. Contoh `Claim Life/RDBList/GetSequenceNumber_SQL.xml`
  (`ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL`):

  ```sql
  BEGIN
    POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER(…);
    COMMIT;          -- ← milik rule Pega, bukan milik procedure
  END;
  ```

  Jadi commit dilakukan **Pega di lapisan aplikasi**. Keduanya konsisten; yang berubah adalah letak
  tanggung jawabnya.

  `[keputusan work owner]` **Sistem baru: Go memegang batas transaksi.** Buka transaksi, panggil
  procedure, commit/rollback **di Go**. Commit **segera setelah nomor terbentuk**, agar lock
  `SELECT … FOR UPDATE` di `GENERATE_SEQUENCE_NUMBER` lekas lepas dan tidak menahan pendaftar lain.

  ⚠️ **Tetap terbuka untuk modul non-Life:** procedure modul lain — mis.
  `PEGA_JSON_POLIS_TREATYIN` — **commit di dalam** dirinya sendiri. Pola Go-memegang-transaksi
  **tidak dapat langsung disalin** ke sana.

**Diperbarui STEP D2 Tahap 4 `[terverifikasi]`.** Pola yang sama muncul di domain facultative:
26–29 stored procedure `POOLDATA.*` dipanggil per modul, termasuk
`POOLDATA.INSERTJSONPOLIS`, `POOLDATA.INSERTUPDATECEDINGPRODUCTION`,
`POOLDATA.INSERTUPDATERISKADDRESS`, `POOLDATA.GENERATE_FACRETRO_NO`,
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`, `POOLDATA.RDBINSERTCLIENT`.

Tiga skema **baru** ikut terlibat: **`GENERAL`** (`GENERAL.F_GET_NM_ASURADUR`), **`MBU`**
(`MBU.F_CEK_HURUF`), dan **`NEW_GENERAL`** (`NEW_GENERAL.CEK_PLAT_NO`, hanya Endorsment Fac In),
ditambah `FIRE.CEK_PRORATA_TANGGAL` dan `FIRE.PEGA_FIRE_SET_RATE` (juga hanya Endorsment Fac In)
→ juga dicatat di OQ-016.

**Pertanyaan tambahan:** `FIRE.CEK_PRORATA_TANGGAL` menunjukkan bahwa **perhitungan prorata tanggal
berada di dalam database**, bukan di rule. Berapa banyak logika bisnis lain yang seperti itu, dan
apakah seluruhnya ikut melakukan `COMMIT` sendiri?


### OQ-014 — Arti singkatan `EDM`

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D1 batch 1 — muncul sebagai nama modul (`EDM Treaty In`), sufiks tabel
  (`M_TREATY_IN_EDM`, `M_TREATY_IN_DETAIL_EDM`, `TREATY_IN_EDM`), dan sufiks rule
  (`CheckRISLIP_EDM`, `DetailPolicyTreatyInNonProportionalEDM`)
- **Pertanyaan:** Kepanjangan `EDM` apa, dan mengapa hampir setiap tabel utama punya pasangan
  `_EDM`? Apakah itu versi amandemen, draft, atau salinan historis?
- **Memblokir:** pemahaman model data inti treaty inward; korpus **tidak** menjabarkan singkatan ini
  sehingga tidak boleh ditebak.
- **Status:** terbuka

### OQ-015 — `NB FacIn`, `RNW Fac In`, `Endorsment Fac In`: tiga modul atau satu ruleset bersiklus

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D1 batch 2 — 6.071 file ketiga modul hanya berisi 2.552 identitas unik;
  **1.567 identitas (61,4%) ada di ketiganya**; **1.906 dari 1.926 identitas `RNW Fac In` (99,0%)
  juga ada di `NB FacIn`**
- **Pertanyaan:** Apakah ketiganya satu aplikasi yang sama dengan pembeda tahap siklus
  (penutupan baru / renewal / endorsement), atau tiga aplikasi terpisah yang berbagi basis kode?
  Bila satu, apa mekanisme pembedanya — parameter, class turunan, atau ruleset berbeda?
- **Memblokir:** penetapan bounded context di D4; perkiraan volume pekerjaan migrasi (selisihnya
  besar: 6.071 file vs 2.552 rule); risiko memigrasi pekerjaan yang sama tiga kali.
- **Status:** **terjawab sebagian** (2026-09-13, STEP D2 Tahap 4) — mekanismenya terbaca; arti nilai `StatusBusiness` dan penetapan bounded context tetap terbuka.
- **Catatan:** pola ini sama dengan OQ-010 (`Treaty In` ↔ `Treaty In Adjustment`). Kemungkinan
  keduanya satu pertanyaan arsitektur yang sama; **belum terverifikasi**.

**TERJAWAB SEBAGIAN — STEP D2 Tahap 4 (2026-09-13). Mekanismenya kini terbaca `[terverifikasi]`.**

Pertanyaan "satu ruleset bersiklus atau tiga aplikasi" terjawab ke arah **satu ruleset**, dengan
empat bukti langsung:

1. **Rule pembeda siklus ada, terbaca, dan ada di KETIGA modul:**

   | Rule `When` | Kondisi literal |
   | --- | --- |
   | `When/IsNB.xml` | `pyWorkPage.Quotation.StatusBusiness = 1` |
   | `When/IsRenewal.xml` | `pyWorkPage.Quotation.StatusBusiness = 2` |
   | `When/IsEDM.xml` | `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` |
   | `When/IsNotEDM.xml` | `pyWorkPage.Quotation.StatusBusiness != 3` |

   Keempatnya ada di `NB FacIn`, `RNW Fac In`, **dan** `Endorsment Fac In` — 12 file.

2. **Prefix ID case sebagai penanda kedua:** `"EDM-"` 262×, `"NB-"` 127×, `"RNW-"` 43× di korpus
   facultative; diperiksa antara lain oleh `RNW Fac In/Activity/ProtectRenewal_Act.xml`
   (`@contains(pyWorkPage.pxInsName,"NB-")`).

3. **Mesin bersama, hash identik lintas modul:** `Activity/SetOldData.xml` (`df291257aa`,
   450.829 byte, 59 langkah), `Activity/CountPremiNusareRetro_Act.xml` (`8f1268a660`),
   `When/IsGroup.xml` (`2cff7b9a67`).

4. **Kode siklus lain ikut terbawa di setiap ekspor:** modul "new business" memuat 35 activity
   ber-nama EDM dan 15 rule `IsEdm*`.

**Kesimpulan `[terverifikasi]`:** ketiga folder menjalankan **satu basis rule** yang membedakan
tahap siklus **saat runtime**; pemisahan folder adalah pemotongan administratif. Yang berbeda per
siklus adalah **subhimpunan rule yang aktif**, bukan basis kodenya.

**Yang MASIH terbuka:** arti nilai `StatusBusiness` 1/2/3 **belum terverifikasi** (korelasi nama
rule bukan bukti — OQ-020); dan penetapan bounded context tetap keputusan **D4**, bukan di sini.

**Kontras terukur dengan domain Claim (Tahap 3):** tumpang tindih identitas facultative **99,0%**
vs Claim **12–40%**; facultative punya rule pembeda siklus, Claim **tidak**; nama Flow sama di
facultative berarti rule sama atau varian isi, di Claim berarti **class berbeda → rule berbeda**.
Bukti lengkap: `flows/_SUMMARY-facultative.md` §7.


### OQ-016 — Sepuluh skema Oracle: apa peran masing-masing dan mana yang dalam ruang lingkup

- **Pemilik:** DBA (+ Product+Underwriting untuk peran bisnisnya)
- **Ditemukan di:** STEP D1 batch 2 — objek dan procedure merujuk **10 skema lokal**:
  `POOLDATA` (dominan), `DATAPEGA`, `NEW_UNDERWRITING`, `GENERAL`, `NEW_GENERAL`, `REINSURANCE`,
  `FIRE`, `ARASAPAS`, `MBU`, `GL`. Batch 1 hanya memperlihatkan dua yang pertama.
  Bukti antara lain `ARASAPAS.INVOICE` (`Claim Fac In/RDBList/CekLunasPremi_Sql.xml`),
  `ARASAPAS.DETAIL_INVOICE` (`Endorsment Fac In/RDBList/SearcStatusBayarArasaps_SQL.xml`),
  `NEW_UNDERWRITING.M_COVERAGE` (`Endorsment Fac In/RDBList/SearchTemplateMainCoverageSQL.xml`)
- **Pertanyaan:** Apa peran tiap skema? Mana yang milik aplikasi ini dan **mana yang milik sistem
  lain** yang kebetulan satu instance Oracle? Mana yang ikut dimigrasikan?
- **Memblokir:** batas ruang lingkup migrasi data; desain repository; keputusan apakah akses lintas
  skema menjadi panggilan antar-service atau tetap query langsung.
- **Status:** terbuka. `[dugaan]` `GL` menyerupai general ledger dan `ARASAPAS` menyerupai sistem
  invoice/pembayaran — **jangan dijadikan fakta sebelum dikonfirmasi**.

**Diperbarui STEP D2 Tahap 5 `[terverifikasi]` — daftar skema diukur ulang menyeluruh.**

Gabungan skema dari nama tabel **dan** nama procedure di seluruh korpus, di luar paket bawaan
Oracle (`DBMS_LOB`, `DBMS_OUTPUT`, `UTL_MATCH`):

```
ARASAPAS  DATAPEGA  FIRE  GENERAL  GL  MBU  NEW_GENERAL  NEW_UNDERWRITING  POOLDATA  REINSURANCE
```
→ **10 skema aplikasi.**

Perintah audit:
```
{ awk -F'\t' '$2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
    | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | grep "\." | grep -v "@" | sed 's/\..*//';
  awk -F'\t' '$2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
    | grep -oE "procs=[^;]*"  | sed 's/procs=//'  | tr ',' '\n' | sed 's/\..*//'; } \
| sort -u | grep -vE '^(DBMS_LOB|DBMS_OUTPUT|UTL_MATCH)$'
```

**Temuan baru Tahap 5:** skema **`DATAPEGA`** ternyata bukan sekadar penyimpanan Pega yang pasif —
ia **dibaca langsung lewat SQL** (`DATAPEGA.PC_ASM_FW_GCNMFW_WORK`, di `Master Product Name Life`)
→ **OQ-058**. Skema `NEW_GENERAL` (procedure `CEK_PLAT_NO`) dan `GL` juga tercatat.


### OQ-017 — Database link `ASMD.SINARMAS.CO.ID`: data rating berada di luar sistem

- **Pemilik:** DBA + Actuarial + Product+Underwriting
- **Ditemukan di:** STEP D1 batch 2 — 8 objek diakses lewat Oracle database link
  `@ASMD.SINARMAS.CO.ID`: `M_EQS_RATE`, `M_TERORISME_RATE`, `M_RSMD_RATE`, `M_FLEXAS_RATE`,
  `M_FLOOD_AREA`, `FIRE.M_FLOOD_RATE`, `FIRE.M_BI_INDEMNITY`, `LST_KURS_STANDARD`.
  Bukti: `NB FacIn/RDBList/SearchRatePolisEQS_SQL.xml`, `NB FacIn/RDBList/GetKurs.xml`,
  `NB FacIn/RDBList/SearchRatePolisFlood_SQL.xml`
- **Pertanyaan:** Database apa yang berada di ujung link ini, siapa pemiliknya, dan apa perjanjian
  aksesnya? Apakah tabel rate dan kurs standar akan tetap diambil dari sana setelah migrasi, atau
  direplikasi? Seberapa sering berubah?
- **Memblokir:** desain perhitungan premi — sebagian **input rating tidak berada di basis data yang
  dimigrasikan**. Juga memblokir rencana cutover: sistem baru tetap bergantung pada sistem luar.
- **Status:** terbuka — **cakupan FINAL, sapuan 100% korpus**

**Sapuan tuntas STEP D1 batch 4 `[terverifikasi]`.** Seluruh 9.369 file (20/20 modul) sudah
diperiksa. Objek `@ASMD.SINARMAS.CO.ID` **hanya muncul di 3 modul**, seluruhnya facultative inward:

| Modul | Objek berbeda | Total rujukan |
| --- | ---: | ---: |
| NB FacIn | 8 | 9 |
| RNW Fac In | 7 | 8 |
| Endorsment Fac In | 8 | 9 |

Objek yang diakses: `M_EQS_RATE`, `M_TERORISME_RATE`, `M_RSMD_RATE`, `M_FLEXAS_RATE`,
`M_FLOOD_AREA`, `FIRE.M_FLOOD_RATE`, `FIRE.M_BI_INDEMNITY`, `LST_KURS_STANDARD`.

**Nol rujukan** dari 17 modul lainnya — termasuk seluruh domain treaty, klaim, dan life.

Artinya ketergantungan pada database luar **terlokalisasi pada perhitungan premi facultative
inward**: tabel rate per peril dan kurs standar. Ini membatasi luas dampaknya, tetapi tidak
menguranginya — tanpa akses ke database itu, perhitungan premi facultative tidak dapat direplikasi.

Perintah audit: cari pola `@` pada daftar tabel hasil ekstraksi `<pyBrowseSQL>` di seluruh modul.

**Diperbarui STEP D2 Tahap 4 `[terverifikasi]` — konteks pemakaian kini diketahui.**

Telusur menemukan **di mana** tabel rate dipakai: rantai limit akseptasi. `GetLimitAkseptasi_ActFlow`
(NB FacIn, 539.687 byte) memulai dengan `Call CountTotalTSIPremiNusaRe_Act` lalu memanggil
**12 RDB-List** di class `ASM-FW-GISFW-Int-policyjson`, terpisah per lini bisnis:
`GetLimitAkseptasi_SQL`, `GetLimitAkseptasiBanding_SQL`, `GetLimitAkseptasiPreferedComm_SQL`,
`GetLimitAkseptasiNonPrefer_SQL`, `GetLimitAkseptasiNonPreferBanding_SQL`,
`GetLimitAkseptasiNonFire_SQL`, `GetLimitAkseptasiNonFireBanding_SQL`,
`GetLimitAkseptasiKreditNCL_SQL`, `GetLimitAkseptasiKreditCL_SQL`, `GetLimitAkseptasiBond_SQL`,
`GetLimitAccEngineeringUW_SQL`, `GetLimitAccEngineeringBanding_SQL`.

**Selisih per modul `[terverifikasi]`:** `RNW Fac In` mengakses **7** objek db-link, bukan 8 —
**`LST_KURS_STANDARD` tidak dirujuk dari modul renewal**. NB FacIn dan Endorsment Fac In mengakses
8. `[dugaan]` renewal memakai kurs yang tersimpan pada data lama; **belum terverifikasi**.

**Ketergantungan luar bertambah satu:** procedure **`FIRE.CEK_PRORATA_TANGGAL`** (perhitungan
prorata tanggal) dan **`FIRE.PEGA_FIRE_SET_RATE`** dipanggil dari
`Endorsment Fac In/RDBList/SearchSQLRateKPR.xml`, `SearchSQLRateNonKPR.xml`, dan
`SearchTemplateMainCoverageSQL.xml` — isinya juga **tidak ada di korpus** (OQ-002).


### OQ-018 — Host dan URL ter-hardcode, termasuk hostname lingkungan DEV

- **Pemilik:** IAM + DBA (pemilik infrastruktur/environment)
- **Ditemukan di:** STEP D1 batch 2 — URL literal di **48 file**, mayoritas **bukan** di
  `ConnectREST` melainkan di `Section`, `FlowAction`, `Activity`, `Harness`. Host yang muncul:
  `192.168.105.112`, `192.168.105.116`, `ssdecamwin02:7070`, `ssdecamwin03:7070`,
  `sdvpwin105:8383`, `10.100.10.75:7315`, `pega.nusantarare.com:80`, `app.sinarmas.co.id`,
  `view.officeapps.live.com`, dan **`appdev.nusantarare.com`**.
  Satu `ConnectREST` juga menyimpang dari pola konfigurasi:
  `Claim Fac In/ConnectREST/getPremiumPaidOnMarine.xml`
  (`ASM-FW-GCNMFW-WORK-PNC / GETPREMIUMPAIDONMARINE`) memakai `<pyBaseURLSelectionType>URL`
  dengan URL literal, sedangkan 19 rule lainnya memakai `SETTING` → `LinkService!LinkService`
- **Pertanyaan:** Alamat mana yang benar untuk production? Mengapa hostname DEV
  (`appdev.nusantarare.com`) ada di korpus — apakah ekspor diambil dari lingkungan dev?
  Sistem apa yang berada di balik tiap host? Apakah `view.officeapps.live.com` (SaaS pihak ketiga)
  memang disengaja, mengingat dokumen dikirim ke layanan luar untuk ditampilkan?
- **Memblokir:** seluruh alamat harus menjadi konfigurasi/env var di arsitektur Go, bukan literal.
  Pertanyaan tentang hostname DEV juga menyentuh **keandalan korpus itu sendiri** — bila ekspor
  berasal dari dev, sebagian rule mungkin bukan versi production (berkaitan dengan OQ-011).
- **Status:** terbuka — **TERTUTUP untuk Claim — Life** (2026-09-14, sumber: **korpus**).

  Penutupan ini **menggantikan** penyempitan sebelumnya (`jboss1073` = production, `jboss117` = dev,
  work owner 2026-09-14), yang tetap berlaku sebagai keterangan lingkungan.

  `[terverifikasi]` **Nol URL endpoint bisnis ter-hardcode di modul `Claim Life`.** Sapuan seluruh
  `http(s)://` di modul ini hanya menemukan **dokumentasi**, bukan runtime:

  | URL yang ditemukan | Jumlah | Sifat |
  | --- | ---: | --- |
  | `https://community.pega.com/help_v88/…` | 132+ | `pyHelpURI` — tautan bantuan Pega |
  | `https://view.officeapps.live.com/op/view.aspx?src=` | 3 | penampil dokumen Office, bukan endpoint bisnis |

  `[terverifikasi]` Pembedaan lingkungan **tidak** memakai hostname tertanam, melainkan flag
  tingkat produksi Pega: `Claim Life/When/IsPEGAPROD.xml` (`@BASECLASS` / `ISPEGAPROD` /
  `RULE-OBJ-WHEN`, 18.306 byte) berisi `pzProductionLevel = "5"` —
  `compareTwoValues(pxProcess.pzProductionLevel,"=","5")`.

  Ini **menjawab OQ-029 juga**: kondisi `IsPEGAPROD` yang semula "tidak terbaca" kini terbaca penuh.

  `[terverifikasi]` Endpoint keluar di-resolve lewat `M_LINK_SERVICE` — lihat **OQ-047**.

  **Aturan sistem baru:** pembedaan dev/prod memakai **profil environment + isi tabel per-database**,
  bukan `if hostname` dan bukan URL di env var (**ADR-0013**).

  URL literal di **15 modul lain** tetap terbuka pada gilirannya.
- **Catatan keamanan:** URL literal pada `getPremiumPaidOnMarine.xml` memuat contoh nilai data
  (nomor polis dan case id) di query string. Nilai itu **sengaja tidak disalin** ke artefak
  discovery; rujuk file aslinya bila perlu.

**Diperbarui STEP D1 batch 3 dan 4:**

- Batch 3 menambah **24 file** ber-URL literal, dan **4 dari 24 `ConnectREST`** memakai
  `<pyBaseURLSelectionType>URL` (batch 2: 1 dari 20). Keempatnya service pencatatan klaim
  Non-Proportional: `InsertClaimOutstanding_NP`, `insertClaimFinalOrClosed_NP` (2 modul),
  `insertClaimReject_NP`. Korpus penuh: **5 dari 51 `ConnectREST`** memuat URL literal.
- `appdev.nusantarare.com` muncul 1× di batch 2 dan 3× di batch 3.
- **Bukti tambahan bahwa korpus dirakit dari lebih dari satu server:** metadata ekspor memuat dua
  nilai `<pxHostId>` berbeda — `jboss1073` dan `jboss122117`.

**ANGKA FINAL KORPUS (sapuan 20 modul, selesai setelah batch 4) `[terverifikasi]`:**
URL literal ditemukan di **81 file di 15 modul**. (Angka 72 yang sempat dilaporkan hanya mencakup
modul batch 2 dan 3; sapuan batch 1 dan 4 menambah 9 file.)

| Modul | File | | Modul | File |
| --- | ---: | --- | --- | ---: |
| NB FacIn | 14 | | Komite Claim Life | 3 |
| RNW Fac In | 13 | | Komite Claim Non Prop | 3 |
| Endorsment Fac In | 13 | | Treaty In | 2 |
| Claim Non Prop | 9 | | Treaty In Adjustment | 2 |
| Claim Fac In | 8 | | Master Product Name Life | 2 |
| Claim Prop | 6 | | EDM Treaty In | 1 |
| Claim Life | 3 | | PremiumList Life | 1 |
| | | | Treaty Contract Out | 1 |

Host yang muncul (host saja; path tidak disalin):

| Host | Kemunculan | Catatan |
| --- | ---: | --- |
| `view.officeapps.live.com` | 27 | **SaaS pihak ketiga** — dokumen dikirim ke layanan luar untuk ditampilkan |
| `ssdecamwin03:7070` | 18 | hostname internal |
| `192.168.105.116:80` | 18 | IP internal |
| `app.sinarmas.co.id` | 12 | domain grup Sinarmas |
| `192.168.105.112` (dan `:80`) | 24 | IP internal |
| `10.100.10.75:7315` | 10 | integration web service |
| `ssdecamwin02:7070` | 9 | hostname internal |
| `sdvpwin105:8383` | 9 | hostname internal |
| `pega.nusantarare.com:80` | 6 | server Pega |
| **`appdev.nusantarare.com`** (dan `:80`) | **7** | **hostname lingkungan DEV** |

Perintah audit:

```
grep -rlE "https?://" "<modul>" --include="*.xml" \
  | while read f; do grep -ohE "https?://[^< \"']+" "$f" | grep -qv "pega.com" && echo "$f"; done
```

**Diperbarui STEP D2 Tahap 4 — cakupan sapuan D1 terbukti KURANG LUAS `[terverifikasi]`.**

Ketiga modul facultative punya **3 rule `RULE-CONNECT-REST` yang sama** (`ServiceGoogle`,
`convertJsonNusareToProduction`, `getPremiumPaidOn`), seluruhnya `pyBaseURLSelectionType = SETTING`
dengan `pyBaseURLSetting = LinkService!LinkService`, dan **tanpa satu pun URL literal endpoint** —
satu-satunya URL di ketiga file adalah `<pyHelpURI>` (tautan dokumentasi Pega).

**Tetapi alamat sesungguhnya tidak ada di file rule sama sekali.** `GetLinkService`
(`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE`) melakukan `Obj-Browse` atas tabel Oracle
`M_LINK_SERVICE` dengan kunci `KATEGORI_1` / `KATEGORI_2`, dan hasilnya dipakai langkah
`Connect-REST`. Sapuan D1 hanya mencari URL literal **di dalam file rule**, sehingga **sumber
konfigurasi yang sebenarnya luput**. → **OQ-047**.

`[terverifikasi]` Rule SystemSettings `LINKSERVICE!LINKSERVICE` terdaftar berkonflik (OQ-011 #324)
tetapi **identik** di ketiga modul setelah normalisasi 21 tag — base URL **tidak** bercabang antar
siklus.


**TERJAWAB untuk konteks Claim — Life — 2026-09-14, sumber: work owner** (grilling Ronde 1 Q3).

`[terverifikasi work owner]` Pemetaan server:

| `pxHostId` | Lingkungan |
| --- | --- |
| **`jboss1073`** | **production** |
| **`jboss117`** | **dev** |

Sistem memakai **mirroring**, sehingga **perbedaan antar `pxHostId` TIDAK dianggap risiko
perilaku**. Untuk `Claim Life`, korpus **dianggap production**.

`[terverifikasi]` Bukti korpus pendukung — `Claim Life` **nol hostname DEV**
(`grep -rlio "appdev\.nusantarare" "Claim Life"` → kosong), dengan sebaran
`pxHostId`: `pega-nusre` 73, `jboss1073` 64, `jboss117` 1, `1d407e1106c501b92d737506992c6d06` 1.

**Status: terjawab untuk Claim — Life.**

**Cakupan korpus-wide TETAP TERBUKA:**
1. Lingkungan untuk `pxHostId` **`pega-nusre`** (73 berkas di Claim Life saja) dan
   `1d407e1106c501b92d737506992c6d06` **belum dinyatakan** work owner.
2. Hostname DEV `appdev.nusantarare.com` di **81 berkas / 15 modul lain** belum tersentuh —
   `Claim Life` bukan salah satunya.
3. **OQ-050** (10 tombol ber-label `(dev)` di `Treaty In`, termasuk `Force Resolve Complete`)
   berdiri sendiri dan tidak dijawab di sini.


### OQ-019 — Trio Claim tidak berbagi basis kode: apa sebenarnya pembedanya

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D1 batch 3 — overlap identitas di trio Claim **jauh lebih rendah** daripada
  trio facultative di batch 2:

| Pasangan | Identitas dibagi | Ukuran modul terkecil | ~% |
| --- | ---: | ---: | ---: |
| Claim Non Prop ↔ Claim Prop | 102 | 270 | ~37% |
| Claim Life ↔ Claim Prop | 19 | 136 | ~14% |
| Claim Life ↔ Claim Non Prop | 16 | 136 | ~12% |

  Bandingkan `RNW Fac In` ↔ `NB FacIn` di batch 2: **99,0%**.
- **Pertanyaan:** Mengapa domain klaim dibangun sebagai tiga basis kode terpisah sementara domain
  facultative satu ruleset bersama? Apakah perbedaan Life / Proportional / Non-Proportional memang
  menuntut logika yang berlainan, atau ini duplikasi historis yang tidak pernah dikonsolidasi?
- **Memblokir:** perkiraan volume pekerjaan migrasi domain klaim (tidak bisa diasumsikan hemat
  seperti facultative); penetapan bounded context di D4.
- **Status:** terbuka. `[dugaan]` pembedanya jenis bisnis, bukan tahap siklus — **belum
  terverifikasi**.
- **Terkait:** OQ-010, OQ-015 (pertanyaan sejenis untuk treaty dan facultative).

### OQ-020 — Arti kode `PaymentType` dan `TransferType` yang menggerbangi percabangan pembayaran klaim

- **Pemilik:** Finance + Product+Underwriting
- **Ditemukan di:** STEP D1 batch 3 — nilai literal yang diuji di dalam rule `[terverifikasi]`:
  - `PaymentType` → **0, 1, 2, 3, 4, 5, 6, 7**
  - `TransferType` → **1, 2, 3, 4**

  Yang penting: nilai dikelompokkan **berbeda-beda** untuk memilih jalur logika berbeda. Dari
  `Claim Prop/Activity/HitServiceToKasir_Act.xml` (`ASM-FW-GCNMFW-WORK / HITSERVICETOKASIR_ACT`),
  tiga kelompok dalam satu rule: `{1,2,5}`, `{4,6}`, `{3}`.
  Dari `Claim Prop/Activity/AttachmentProtect_ACT.xml`: `{1,2,3}` dan `{2,4}`.
  Dari `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml`: `==7` dan `!=7`.
  Rule lain yang memakainya sebagai guard: `Claim Non Prop/Activity/HitServiceToKasir_Act.xml`,
  `Claim Non Prop/Activity/SetInterimXOL_Act.xml`,
  `Komite Claim Prop/Activity/HitServiceToKasirKMT_Act.xml`,
  `Komite Claim Non Prop/Activity/HitServiceToKasirKMT_Act.xml`,
  `Komite Claim Prop/Section/ShowTransfer.xml`
- **Pertanyaan:** Apa arti tiap nilai `PaymentType` (0–7) dan `TransferType` (1–4)? Mengapa
  pengelompokannya berbeda antar rule — apakah tiap kelompok mewakili kategori pembayaran yang
  berbeda, dan apakah pengelompokan yang berbeda itu disengaja atau inkonsistensi?
- **Memblokir:** seluruh telusur D2 atas jalur pembayaran klaim; tidak mungkin menulis aturan bisnis
  pembayaran tanpa mengetahui arti kodenya. Korpus **tidak memuat** tabel kode atau label yang
  menjelaskannya, dan arti **tidak boleh ditebak** dari nama rule.
- **Status:** terbuka — **terjawab untuk Claim — Life** (2026-09-14, work owner): `STS_REJECT` 0/1/2, `AcceptStatus` 1/2, `SendtoAdmin`, `ContentNote`. Kode di konteks lain tetap terbuka.
- **Kode `.Type` — sebagian terjawab (2026-09-14, work owner, Ronde 4):**
  `[keputusan work owner]` **`TP` = Payable, `TR` = Receivable.** Kepanjangan ini **tidak ada di
  korpus**; sumbernya work owner. Dicatat di **ADR-0012** dan `CONTEXT.md`.
  **`QP` dan `QR` BELUM dijawab** — tetap terbuka di sini.

  `[terverifikasi]` `.Type` **menggerbangi wewenang**, bukan sekadar mencabangkan data: gerbang
  "send ke Komite" di `Claim Life/Section/AdjustmentDetail_Section.xml` bebas-peran untuk `TP`/`TR`.
  Ia juga menggerbangi jendela validasi Date of Loss di `Claim Life/Activity/ValidasiDOL_Act.xml`
  (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!VALIDASIDOL_ACT`, 59.747 byte):

  | Cabang | Jendela | Pergeseran tanggal |
  | --- | --- | --- |
  | `Type=="QR" \|\| Type=="QP"` | `GROSS_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | **nol** |
  | `Type=="TR" \|\| Type=="TP"` | `RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | **+1 hari** |

  `[dugaan]` — **jangan dipakai sebagai fakta.** Pola di atas memperlihatkan huruf **pertama**
  memisahkan gross (`Q*`) dari retrosesi (`T*`). Bila huruf **kedua** memang Payable/Receivable
  seperti pada `TP`/`TR`, maka `QP`/`QR` mengikuti pola yang sama. **Satu konfirmasi work owner akan
  menutup sisa butir ini** — belum ditanyakan.

  `[terverifikasi]` Pencarian kepanjangan di korpus **sudah tuntas dan nihil**: tidak ada label,
  caption, atau opsi mana pun; satu-satunya tempat nilainya ditetapkan adalah
  `PremiumList Life/Activity/SubmitPremiumList_Act.xml` (`ASM-FW-GISFW-WORK-LIFE!SUBMITPREMIUMLIST_ACT`),
  dan di sana ia hanya diteruskan ke slot generik `InputData.CARI20` dengan precondition
  `.Type=="<kode itu sendiri>"` — melingkar. Lihat **OQ-063** (tertutup) dan **OQ-059**.

**Diperbarui STEP D1 batch 4 `[terverifikasi]`** — tiga enumerasi lain yang juga menggerbangi logika:

| Properti | Nilai literal | Bukti |
| --- | --- | --- |
| `EdmType` | `1`, `3` | `Endorsement Life/Activity/SetPremi_EDM.xml`, `Endorsement Life/Activity/SaveCSVEDMLife.xml` |
| `ProRateType` | `1`, `2`, `3` | `PremiumList Life/Activity/Calculate1_Act.xml`, `PremiumList Life/Activity/SavePremiumList_Act.xml` |
| `.Type` (PremiumList) | `QR`, `QP`, `TP`, `TR` | `PremiumList Life/Activity/GetPLNumber_Act.xml`, `PremiumList Life/Activity/SubmitPremiumList_Act.xml` |

Kode `QR/QP/TP/TR` berada di dalam `<pyStepsPreCondParamsWhen>` — **menggerbangi eksekusi step**,
mis. `.Type=="QR"` dan `.Type = "TR" || .Type = "TP"`. `ProRateType` menggerbangi perhitungan di
`Calculate1_Act`.

Arti seluruh nilai di atas **belum terverifikasi**. Korpus tidak memuat tabel kode. Ini memblokir
D2 atas perhitungan premium list dan endorsement Life.

**D2 Tahap 3 — tabel `PaymentType` per modul Claim `[terverifikasi]`:** Claim Life **tidak memakainya sama sekali**; Claim Prop `1–6`; Claim Non Prop `1–7`; Claim Fac In `1–7` (pemakaian terpadat: `=4` 23×, `=6` 21×, `==3` 19×). `TransferType` = `2` di tiga modul treaty/fac, tidak ada di Life. Kode lain yang ditemukan: `STS_REJECT` (`'0'`,`'1'`,`'2'`), `.pyNote="Back"`, `BusinessCode` `L1`…`L11` (OQ-038), `ContentNote="DEATH"`. Arti seluruhnya **belum terverifikasi** — bukti Tahap 2 (`KomitePost_Reject` digerbangi `AcceptStatus=="1"`) menegaskan mapping tidak boleh ditebak.

**TERJAWAB untuk konteks Claim — Life — 2026-09-14, sumber: work owner** (grilling Ronde 1 Q5).

`[terverifikasi work owner]`

| Kode | Nilai | Arti |
| --- | --- | --- |
| **`STS_REJECT`** | `0` | **Outstanding** |
| | `1` | **Aksep** (diterima) |
| | `2` | **Reject** (ditolak) |
| **`AcceptStatus`** | `1` | diaksep |
| | `2` | reject |
| **`SendtoAdmin`** | `1` | dari `ReasLifeMedicalAdvisor` **atau** `ReasLifeSPV`, kasus dikembalikan ke `ReasLifeAdmin` |
| **`ContentNote`** | `DEATH`, `HEALTH`, `CI`, `TPD`, `TI` | **jenis klaim**, diturunkan dari `BusinessCode` |

⚠️ **PERINGATAN NAMA MENYESATKAN — direkam atas permintaan work owner.**
Field bernama **`STS_REJECT`**, tetapi nilai **`1` berarti *diaksep*, bukan ditolak**. Membaca nama
field sebagai artinya akan menghasilkan logika terbalik. Dicatat juga di `CONTEXT.md`.

`[terverifikasi]` Konsisten dengan pola korpus: `STS_REJECT=='1' || STS_REJECT=='2'` muncul 7× —
yaitu "sudah selesai diproses" (aksep **atau** reject), lawan dari `0` (outstanding).

**Cakupan korpus-wide TETAP TERBUKA:** `PaymentType` (1–7), `ProposalAcceptStatus` (1/2/3/4/7/9),
`EdmType`, `QuotationData.Type` (12 nilai), `StatusBusiness`, `ProRateType`, `TransferType`,
`StatusAkseptasi` treaty, `REINSTYPEID`, kode `OR` — **seluruhnya belum tersentuh**.
`PaymentType` memang **tidak dipakai** di `Claim Life`.


### OQ-022 — Modul `Treaty Contract Out`: isinya tidak menunjukkan treaty *outward*

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D1 batch 4 — pemeriksaan apakah nama folder menyesatkan
- **Bukti `[terverifikasi]`:**
  - Objek database yang disentuh seluruhnya bertema kontrak treaty (`M_PROPORTIONALARRG`,
    `MTREATYSECURITY`, `TREATYREINSURER`, `TREATYBUSINESS`, `TREATYCONTRACT`, `M_TREATYYEAR`,
    `TREATYEXCHANGE`) — **tidak satu pun bernama `*OUT*`**.
  - Objek treaty outward (`M_TREATY_OUT`, `TREATY_OUT2`, `M_TREATY_OUT_DETAIL`, `FACOUTPRODUCTION`)
    dirujuk oleh **tujuh modul lain** di korpus, dan **nol** dari modul ini.
  - Dari 303 rule, hanya 5 menyebut `TreatyOut` — seluruhnya soal lampiran
    (`LOADATTACHMENTTREATYOUT`, `TREATYOUTDOWNLOADALL_ACT`, `TREATYOUTDOWNLOADONE`,
    `TREATYOUTSAVEATTACHMENT`, `TREATYOUTATTACHCONTENT`).
  - Sebaliknya ada rujukan treaty **inward** eksplisit: `BROWSEDELETEROWTREATYINCONTRACT`,
    `SETCATEGORYATTACHTREATYIN`, dan class `ASM-FW-GISFW-INT-TREATY_IN` pada 5 rule.
  - 16 rule `CANCELACTIVITY*` menamai jenis klausul kontrak treaty (bordereaux, cash loss limit,
    claims cooperation, EPI, ex-gratia, profit commission, RI comm, treaty limit, dan lainnya).
- **Pertanyaan:** Apa sebenarnya cakupan modul ini — master/klausul kontrak treaty (inward?), atau
  memang ada sisi outward yang tidak terlihat dari rule yang diekspor? Mengapa dinamai "Out"?
- **Memblokir:** penetapan bounded context di D4; risiko salah menempatkan seluruh modul 303 rule
  ke konteks yang keliru. Juga memengaruhi OQ-005 (modul ini tanpa rule `Flow`, sehingga titik masuk
  D2-nya harus ditentukan lebih dulu).
- **Status:** **terjawab sebagian** (2026-09-13, STEP D2 Tahap 5) — pertanyaan "outward atau bukan" **terjawab: BUKAN outward**; menyisakan alasan penamaan dan penetapan konteks (D4). **Jangan menyimpulkan dari nama folder.**

**TERJAWAB untuk pertanyaan "outward atau bukan" — STEP D2 Tahap 5 (2026-09-13)
`[terverifikasi]`. Jawaban: BUKAN outward.**

Telusur penuh dari Harness/FlowAction (`flows/Treaty Contract Out.md` §2) menjalankan **lima uji
terpisah**, seluruhnya dapat diaudit ulang:

| Uji | Hasil |
| --- | --- |
| **Objek database** | **0** objek bernama `*_OUT*`. Yang ada: `M_PROPORTIONALARRG` (10 rule), `MTREATYSECURITY` (5), `TREATYREINSURER` (3), `TREATYBUSINESS` (3), `TREATYEXCHANGE`, `TREATYCONTRACT`, `PROPORTIONALARRG`, `M_TREATYYEAR` |
| **Class rule** | **0** class bernama `*OUT*`; justru **5 rule berclass `ASM-FW-GISFW-INT-TREATY_IN`** dan 11 file menyebut class itu |
| **Penamaan rule** | 81 activity `*TreatyArr*` vs **4** `*Out*` (seluruhnya lampiran) |
| **Arah tulis** | **9 Connect-SQL MENULIS master** lewat `PEGA_TREATYCONTRACT`, `PEGA_TREATYYEAR`, `PEGA_TREATYREINSURER`, `PEGA_TREATYBUSINESS`, `PEGA_PROPORTIONALARRG`, `PEGA_M_PROPORTIONALARRG_CHILD`, `PROSESCOPY` |
| **Pembanding korpus** | Modul yang benar-benar menyentuh objek treaty outward: `NB Treaty In` **8** file, `Claim Non Prop` **4**, `EDM Treaty In` **4**, `Treaty In Adjustment` 2, `Claim Prop` 1, **`Treaty Contract Out` 0** |

**`Treaty Contract Out` adalah editor master *term / arrangement* kontrak treaty**, dengan 16 jenis
klausul yang dapat diedit (`CancelActivity*`: BordereAux, CashLossLimit, ClaimCoorperation, EPI,
ExGratia, FacIn, PLA, ProfitCommision, Ricomm, TreatyContract, TreatyLimit, Portfolio, TerrLimit,
dan lainnya). Kata "Out" hanya bertahan di 5 rule lampiran.

`[terverifikasi]` Satu-satunya class treaty-outward di **seluruh korpus** adalah
`ASM-FW-GISFW-INT-TREATYOUTDETAIL` (5 rule), letaknya di `NB Treaty In` (3),
`Treaty In Adjustment` (1), `EDM Treaty In` (1). Ini juga melengkapi **OQ-042**: klaim
non-proporsional membaca master treaty outward, dan sumbernya **bukan** modul ini.

**Menyisakan (tetap terbuka):** *mengapa* modul ini dinamai "Out", dan ke bounded context mana
303 rule ini ditempatkan. Penetapan konteks adalah **D4**.


### OQ-021 — Identitas orang ter-hardcode sebagai guard otorisasi di dalam rule

- **Pemilik:** IAM (+ Product+Underwriting untuk aturan bisnisnya)
- **Ditemukan di:** STEP D1 batch 3, tetapi **berlaku lintas batch** — pola ini luput dari
  pemeriksaan batch 1 dan 2 dan baru dicari sekarang.
- **Temuan `[terverifikasi]` — angka korpus penuh (sapuan 20 modul, STEP D1 batch 4):**
  **66 file di 9 modul** membandingkan identitas operator dengan **nilai literal** — bukan dengan
  peran atau hak akses. (Angka batch 3 sebelumnya, 45 file / 6 modul, memakai pola pencarian yang
  hanya mengenali tanda kutip ganda; sapuan ini mengenali kutip tunggal juga.)

| Modul | File |
| --- | ---: |
| NB FacIn | 17 |
| Endorsment Fac In | 14 |
| RNW Fac In | 13 |
| NB Treaty In | 8 |
| Komite Claim FacIn | 4 |
| Treaty In | 3 |
| Treaty In Adjustment | 3 |
| Komite Claim Prop | 3 |
| EDM Treaty In | 1 |
| **Total** | **66** |

  **18 varian identifier** dengan total ±463 kemunculan. Yang terbanyak satu identifier muncul
  **221×**. Nilai literalnya tidak diulang di sini — daftar lengkap dapat dihasilkan ulang dengan
  perintah audit di bawah, dan sengaja **tidak disebarkan ke artefak lain**.

  **Temuan tambahan `[terverifikasi]`:** beberapa orang yang sama ditulis dalam **beberapa ejaan
  berbeda** (huruf besar, huruf campur, dengan garis bawah, dengan angka di belakang) — satu orang
  bisa punya 4 ejaan. Artinya pencocokan string ini rapuh, bukan sekadar tidak rapi.

  **Di mana guard itu berada** (per tag XML, seluruh korpus):

| Tag | Kemunculan | Pengaruh |
| --- | ---: | --- |
| `<pyStepsPreCondParamsWhen>` | 193 | precondition **eksekusi step** di Activity |
| `<pyCondition>` | 58 | kondisi umum |
| `<pyOtherOperandLeft>` | 24 | operand ekspresi |
| `<pyReadOnlyCondition>` | 18 | field dapat diedit atau tidak |
| `<pyExpression>` | 12 | ekspresi |
| `<pyDisabledWhen>` | 12 | kontrol UI aktif/nonaktif |
| `<pyParametersParamValue>` | 9 | nilai parameter |
| `<pyLabel>` | 9 | label |
| `<pyContainerVisibleWhen>` | 8 | **visibilitas** bagian layar |
| `<pyUnmodifiedPath>` / `<pyRequiredWhen>` / `<pyMemo>` | 6 | lain-lain |

  **Rule apa yang digerbangi** — sebaran tipe: `Activity` 42, `Section` 19, `When` 3, `Harness` 2.
  Yang material, guard ini berada di dalam **perhitungan premi** dan **keputusan komite klaim**,
  bukan sekadar kosmetik UI:

| Identitas rule | Modul |
| --- | --- |
| `ASM-FW-GISFW-WORK / COUNTGROSSPREMI_ACT` | NB FacIn, Endorsment Fac In |
| `ASM-FW-GISFW-WORK / COUNTGROSSPREMIEDM_ACT` | NB FacIn, Endorsment Fac In |
| `ASM-FW-GISFW-WORK / COUNTGPWMARINEPAMBU_ACT` | NB FacIn, Endorsment Fac In |
| `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` (+4 varian lini) | NB FacIn, Endorsment Fac In |
| `ASM-FW-GCNMFW-WORK-KOMITE / KOMITEPOST_REJECT` | Komite Claim FacIn |
| `ASM-FW-GCNMFW-WORK-KOMITE / KOMITEPOST_CLOSECLAIM` | Komite Claim FacIn |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY / KOMITEPOST_REJECT` | Komite Claim Prop |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY / KOMITEPOSTADJUSTMENT` | Komite Claim Prop |
| `ASM-FW-GISFW-WORK / ISGROUP` | Endorsment Fac In |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINADDENDUM` | EDM Treaty In |

  **Guard berbasis jabatan** — terpisah dan lebih menjanjikan: **35 file di 14 modul** memakai
  `OperatorID.pyPosition` dibandingkan dengan string literal. Nilainya **terlihat seperti nama peran
  yang sebenarnya**: `IT Developer` (42×), `ReasLifeAdmin` (13×), `ReasLifeSPV` (8×),
  `ReasLifeMedicalAdvisor` (8×), `Admin` (6×), `SPV A`, `SPV B`.

  `[dugaan]` konsep peran **sudah ada** di sistem, hanya ditegakkan lewat perbandingan string, bukan
  lewat model otorisasi. Ini titik awal paling konkret untuk memetakan RBAC — **belum terverifikasi**.

  **Perintah audit:**

```
grep -rlE "OperatorID\.pyUser(Identifier|Name)[ ]*[=!]+[ ]*[\"']" "<modul>" --include="*.xml"
grep -rhoE "<py[A-Za-z]+>[^<]*OperatorID\.pyUser(Identifier|Name)[^<]*" "<modul>" --include="*.xml"
grep -rlE "pyPosition[ ]*[=!]+[ ]*[\"']" "<modul>" --include="*.xml"
```

  **Diperbarui STEP D2 Tahap 2 — domain Komite `[terverifikasi]`:** guard identitas berada
  **di dalam activity keputusan komite**, bukan di lapisan UI. Sebarannya tidak merata:

| Modul Komite | File ber-guard nama orang | Pola alternatif |
| --- | ---: | --- |
| Komite Claim FacIn | 4 | — |
| Komite Claim Prop | 3 | — |
| Komite Claim Life | 0 | sasaran routing dari `.KomiteID` |
| Komite Claim Non Prop | 0 | **mencocokkan pengguna dengan roster** |

  `Komite Claim Non Prop/Activity/KomitePostAdjustment.xml` memakai
  `@contains(@toUpperCase(KomiteList(IdxKomite).KomiteID), @toUpperCase(OperatorID.pyUserIdentifier))`
  — pengguna dicocokkan dengan **entri roster**, bukan daftar nama ter-hardcode.

  **Korpus memuat dua pendekatan otorisasi untuk proses yang setara**, dan salah satunya
  sudah berbasis data. Ini titik awal paling konkret untuk RBAC di domain Komite.

- **Pertanyaan:** Apa maksud tiap perbandingan itu — hak akses khusus, pengecualian sementara, atau
  workaround? Peran (role) apa yang sebenarnya diwakili tiap orang tersebut? Apakah orang-orang itu
  masih menjabat? Apa arti `"IT Developer"` sebagai guard — apakah layar tertentu memang disembunyikan
  dari developer, atau sebaliknya?
- **Memblokir:** seluruh desain RBAC. Perilaku ini **tidak boleh dipindahkan apa adanya** ke sistem
  Go — identitas orang harus menjadi peran/hak akses. Tetapi memetakannya ke peran mustahil tanpa
  mengetahui maksud aslinya. Juga berkaitan dengan **OQ-007** (tidak ada rule identitas/otorisasi
  di korpus).
- **Status:** terbuka — **terjawab untuk Claim — Life** (2026-09-14, work owner): 3 kode peran lengkap, pemetaan tahap→peran dikonfirmasi, nol identitas orang ter-hardcode di modul ini. Identitas ter-hardcode di 6 modul lain + OQ-053 tetap terbuka.
- **Catatan data pribadi:** identifier yang dipakai sebagai guard adalah **nama/ID pengguna nyata**.
  Yang dicatat di artefak discovery adalah **nama rule, path, tag, dan jumlah** — bukan daftar
  namanya. Nilai literalnya dapat dihasilkan ulang kapan saja dengan perintah audit di atas, dan
  sebaiknya tidak disalin ke dokumen yang beredar lebih luas dari tim proyek.

**Diperbarui STEP D2 Tahap 4 — pola ketiga `[terverifikasi]`.**

Domain facultative memakai **tiga** rule ber-guard identitas orang, dan ketiganya **identik isinya
di NB FacIn, RNW Fac In, dan Endorsment Fac In** (bukan variasi per modul):

| Rule | Tag yang memuat identitas | Jumlah identitas berbeda |
| --- | --- | ---: |
| `When/IsGroup.xml` (hash `2cff7b9a67`) | `OperatorID.pyUserIdentifier` | **3** |
| `When/IsGroupCreate.xml` | `pyWorkPage.pxCreateOperator` | **2** |
| `When/IsTBonding.xml` | `.OfferFacIn.QuotationData.MarketingName` | **1** |

**Pola ketiga yang baru:** `.MarketingName` adalah **field data bisnis**, bukan field operator —
berbeda dari OQ-021 biasa (`pyUserIdentifier`/`pyPosition`) dan dari OQ-027 (`pyTelephone`).

Ketiganya menggerbangi percabangan nyata: `Decision1` "Is it group?", `Decision30`
"pxCreateOperator …", `Decision25` "UW FINANCIAL?" — di **ketiga** rule `Flow` utama facultative.
Label `<pyMOName>` shape `Decision30` juga memuat nama orang.

**Nilai nama orang tidak disalin** ke artefak mana pun (`flows/_METHOD.md` §1.4). Perintah audit
yang menghitung tanpa menampilkan nilai:
```
for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do
  grep -oE "pyUserIdentifier\]\[&amp;#61;\]\[&amp;quot;[A-Z]+" "$m/When/IsGroup.xml" | sort -u | wc -l
done
```


**Diperbarui STEP D2 Tahap 5 `[terverifikasi]` — pola keempat, dan yang paling berat.**

Sampai Tahap 4, seluruh temuan guard identitas bersifat **menguji**: `OperatorID.pyUserIdentifier`,
`pyPosition` (OQ-021), `OperatorID.pyTelephone` (OQ-027), `.MarketingName` (Tahap 4).

Tahap 5 menemukan pola yang berbeda jenisnya: `Treaty In/DataTransform/Akseptasi_DT.xml`
**menetapkan** `TreatyIn.PositionUsername` dengan **5 identitas orang ter-hardcode** sebagai
**pemilik tugas berikutnya** pada setiap cabang `Accept`. Ini bukan pengujian, melainkan
**perutean ke orang tertentu**. Dilacak tersendiri sebagai **OQ-053**.

Sebaran guard identitas di kelima modul Tahap 5 `[terverifikasi]`:

| Modul | `pyUserIdentifier` | `pyTelephone` | `pxCreateOperator` (guard) |
| --- | ---: | ---: | ---: |
| Treaty In | 3 file | 1 file | 0 |
| Treaty In Adjustment | 3 | 1 | 0 |
| Treaty Contract Out | 2 | 0 | 0 |
| Master Product Name Life | 3 | 0 | 0 |
| Master Contract Retro Life | **0** | 0 | 0 |


**TERJAWAB untuk konteks Claim — Life — 2026-09-14, sumber: work owner** (grilling Ronde 1 Q4).

`[terverifikasi work owner]` **Tiga kode peran adalah daftar lengkap** untuk siklus klaim Life,
dengan pemetaan tahap → peran (**bukan lagi dugaan**):

| Tahap | Peran |
| --- | --- |
| Register + Outstanding | **`ReasLifeAdmin`** |
| Medical Check | **`ReasLifeMedicalAdvisor`** |
| Claim Analis | **`ReasLifeSPV`** |

**Rangkap peran tidak diperbolehkan**, kecuali akses ditambahkan eksplisit pada role akun.

`[terverifikasi]` Temuan korpus pendukung — **`Claim Life` tidak memuat identitas orang
ter-hardcode sama sekali**. `pyPosition` dibandingkan terhadap tiga **kode peran** di 17 berkas;
`OperatorID.pyUserIdentifier` (2 berkas) dan `pyUserName` (4 berkas) dipakai sebagai **data jejak**,
bukan guard terhadap literal nama orang:

```
grep -rhoE "(pyUserIdentifier|pyUserName)[^<]{0,45}" "Claim Life" --include="*.xml" \
  | sed 's/&amp;#61;/=/g;s/\]\[/ /g;s/[][]//g' | grep -E "[=!]"      # -> kosong
```

**Koreksi artefak:** nilai `IT Developer` yang tercatat di D1 untuk `pyPosition` berasal dari sapuan
korpus-wide — **tidak ada di `Claim Life`**.

**Cakupan korpus-wide TETAP TERBUKA:** identitas orang ter-hardcode di `Komite Claim FacIn`
(4 berkas), `Komite Claim Prop` (3), `NB FacIn` (43), `RNW Fac In` (36), `Endorsment Fac In` (33),
`NB Treaty In` (12) — **tidak tersentuh**. Demikian pula **OQ-053** (5 identitas orang ditetapkan
sebagai pemilik tugas berikutnya di `Akseptasi_DT`), **OQ-027**, **OQ-045**, **OQ-051**.


### OQ-023 — Shape flow tanpa connector masuk: jalur mati atau ekspor tidak lengkap

- **Pemilik:** Product+Underwriting (+ pemilik export Pega)
- **Ditemukan di:** STEP D2 Tahap 1 — **tiga dari empat konteks** yang ditelusur punya shape yang
  hanya memiliki connector **keluar**, tidak ada satu pun `<pyTo>` yang menujunya:

| Konteks | Shape | Isi jalur |
| --- | --- | --- |
| NB Treaty In | `Decision5` → `Assignment5` → `Decision7` → `Assignment1` → `Decision1` | jalur **komite klaim + persetujuan Direktur** (5 shape) |
| Endorsement Life | `Assignment1` "Input EDM Summary" | layar ringkasan endorsement |
| PremiumList Life | `Assignment1` "Input Premium List Summary" | layar ringkasan premium list |

- **Pertanyaan:** Apakah shape ini benar-benar tidak terjangkau (jalur mati), ataukah dicapai lewat
  jalan lain (menu, harness, `Flow-Jump`), ataukah connector masuknya hilang saat ekspor?
  Untuk NB Treaty In: apakah alur persetujuan Direktur memang tidak aktif?
- **Memblokir:** kelengkapan pemahaman proses di D2. Bila jalur itu aktif lewat mekanisme lain,
  telusur D2 melewatkan bagian nyata dari proses; bila mati, ia tidak boleh ikut dimigrasikan.
- **Status:** terbuka — **dipersempit** (2026-09-15, sumber: **korpus**).

  `[terverifikasi]` Untuk `PremiumList Life/InputPolicyHolder.xml`
  (`ASM-FW-GISFW-WORK-LIFE` / `INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW`), shape **`Assignment1` =
  "Input Premium List Summary"** (label baris 1234, subscript baris 1248) **tidak punya jalur masuk
  yang terbaca**.

  Sapuan seluruh kemunculan `Assignment1` di berkas flow:

  | Baris | Tag | Artinya |
  | ---: | --- | --- |
  | 328 | `<pyTaskID>` | id tugasnya sendiri |
  | 585 | `<pyFromTaskName>` | ia **sumber** connector |
  | 1215 | `<pyMOId>` | id shape |
  | 1248 | `<pxSubscript>` | subscript shape |
  | 1900 | `<pyFrom>` | ia **sumber** connector |

  **Tidak pernah muncul sebagai `<pyTo>` maupun `<pyToTaskName>`** — nol connector masuk.

  `[terverifikasi]` Dua jalur alternatif juga nihil:
  - **Harness:** nol dari 9 harness modul merujuk `ShowLifePremiumSummary` maupun `Assignment1`.
  - **Ticket:** `Ticket1` berjenis `Data-MO-Event-Exception` dan menuju `Decision2`.

  Shape itu **punya** FlowAction (`PremiumList Life/FlowAction/ShowLifePremiumSummary.xml`) dan
  **punya** connector keluar — dapat **ditinggalkan**, tidak dapat **dicapai**.

  **Yang tersisa terbuka:** apakah ia dicapai lewat mekanisme **di luar 17 tipe rule yang
  diekspor** (navigasi/menu Pega), atau memang mati. **Tidak dapat ditentukan dari ekspor.**
  Pemilik: **Arsitektur Pega**.
- **Audit:** bandingkan daftar `<pyTo>` dengan daftar `<pxSubscript>` shape pada file Flow.

**D2 Tahap 2 `[terverifikasi]`:** keempat flow Komite **tertutup** — tidak ada shape tanpa connector masuk. Pola ini sejauh ini hanya muncul di NB Treaty In, Endorsement Life, dan PremiumList Life (Tahap 1).

### OQ-024 — Pemetaan Assignment → workbasket tidak terbaca

- **Pemilik:** IAM + Product+Underwriting
- **Ditemukan di:** STEP D2 — seluruh Assignment di NB Treaty In dan EDM Treaty In memakai
  `<pyImplementation>WorkBasket` dengan `<pyRouteTo>Custom`, sedangkan `<pyWorkBasket>` **kosong**.
  Nama workbasket yang muncul di flow: `ReasTreatyInAdmin`, `ReasTreatyInSecHead`,
  `ReasTreatyInGroupLeader`, `ReasTreatyInDeptHead`, `ReasTreatyInDirector`
  (`NB Treaty In/Flow/InputRealizationTreatyIn.xml`, tag `<pyRuleName>`)
- **Pertanyaan:** Assignment mana dirutekan ke workbasket mana? Routing `Custom` dievaluasi oleh
  rule apa? Siapa anggota tiap workbasket?
- **Memblokir:** pemetaan peran → langkah proses; seluruh desain RBAC. Ini bukti paling konkret
  untuk **OQ-007**, tetapi belum cukup untuk memetakan otorisasi.
- **Status:** terbuka — **maju sebagian** (2026-09-14, work owner): pemetaan tahap→peran Claim Life dikonfirmasi. Pemetaan Assignment→nama workbasket di facultative - **Status:** terbuka treaty tetap terbuka.

**D2 Tahap 2 `[terverifikasi]`:** sasaran routing kini **sebagian terbaca**. Di tiga modul Komite sasarannya string ter-hardcode `komitepnc`…`komitepnc4` (OQ-036); di Komite Claim Life sasarannya nilai data `.KomiteID`. Roster anggota berasal dari class `ASM-FW-GCNMFW-Int-EMAILKOMITE` (OQ-037). Jadi untuk domain Komite pertanyaannya bergeser: bukan lagi "tidak terbaca", melainkan "apa isi `komitepnc*` dan tabel EMAILKOMITE".

**MAJU SEBAGIAN untuk konteks Claim — Life — 2026-09-14, sumber: work owner** (grilling Ronde 1 Q4).

Pemetaan **tahap → peran** untuk `Claim Life` kini dikonfirmasi (lihat OQ-021). Ini menjawab
"siapa yang mengerjakan assignment mana" **untuk konteks ini**, meski graf tetap tidak menyatakannya
(`<pyWorkBasket>` kosong, `<pyRouteTo>` = `Custom`).

**TETAP TERBUKA:** pemetaan Assignment → **nama workbasket** di konteks lain — khususnya
facultative inward (16 nama `ReasFacIn*`/`ReasFacOut*`) dan treaty inward (`ReasTreatyIn*`).
`Claim Life` memakai `pyPosition`, bukan nama workbasket, sehingga jawabannya tidak dapat
digeneralisasi.


### OQ-025 — `SERVICEINSERTARASAPAS_ACT` dipanggil NB Treaty In tetapi variannya tidak ada di modul itu

- **Pemilik:** Product+Underwriting (+ pemilik export Pega)
- **Ditemukan di:** STEP D2 — `NB Treaty In/Activity/serviceInsertArasapas_act.xml`
  (`ASM-FW-GISFW-DATA-POLICYTREATYIN` / `SERVICEINSERTARASAPAS_ACT`) berisi **satu langkah**:
  `Call serviceInsertArasapas_act` dengan `<pyStepsObjectName>pyWorkPage` (class
  `ASM-FW-GISFW-WORK`). Identitas target `ASM-FW-GISFW-WORK / SERVICEINSERTARASAPAS_ACT` terdaftar
  di **register OQ-011 entri #415** dengan **2 isi berbeda**, dan variannya hanya ada di
  `EDM Treaty In` dan `Endorsment Fac In` — **tidak di NB Treaty In**
- **Pertanyaan:** Varian mana yang dieksekusi saat flow NB Treaty In berjalan?
- **Memblokir:** langkah terakhir flow NB Treaty In ("HIT SERVICE ARASAPAS") **tidak dapat
  dijelaskan**. Telusur D2 dihentikan di titik itu sesuai `flows/_METHOD.md` §1.1 aturan 4 —
  varian dari modul lain **tidak dipinjam**.
- **Status:** terbuka

**Diperbarui STEP D2 Tahap 4 — naik ke `[terverifikasi]`, cakupan meluas, risiko turun.**

**(a) Class sasaran panggilan kini terbukti, bukan lagi dugaan `[terverifikasi]`.**
`NB FacIn/Activity/serviceInsertArasapas_act.xml` memuat satu langkah `Call
serviceInsertArasapas_act` dengan **`<pyStepsClassName>ASM-FW-GISFW-Work`** — class sasaran
tertulis eksplisit di tag langkah. Di `NB Treaty In.md` §5.2 hal itu hanya dapat **diduga** dari
class `pyWorkPage`.

**(b) Ada DUA identitas berbeda dengan nama file yang sama `[terverifikasi]`:**

| `pxInsName` | Ukuran | Hash | Ada di |
| --- | ---: | --- | --- |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN!SERVICEINSERTARASAPAS_ACT` | 18.013 | `b013027e0c` | NB Treaty In, NB FacIn, RNW Fac In |
| `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` | 232.000 | `ab3ae3c9ba` | EDM Treaty In |
| `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` | 231.975 | `7314b6c49c` | Endorsment Fac In |

Yang pertama adalah **pembungkus 18 KB**; implementasinya 232 KB dan hanya ada di dua modul.

**(c) Blocker meluas ke NB FacIn dan RNW Fac In.** Keduanya memanggil implementasi itu
(lewat `SetToJsonOffer_ACT`, `SaveToProduction_ACT`, `serviceInsertArasapasRNW_act`, dan shape
`Utility2 "HIT SERVICE ARASAPAS"`) tetapi **tidak memuatnya**. Telusur dihentikan di sana; varian
modul lain **tidak dipinjam**.

**(d) Risiko salah tafsir jauh lebih kecil dari perkiraan D1 `[terverifikasi]`.** Kedua varian
implementasi yang ada di korpus berbeda **hanya pada 3 baris**, seluruhnya metadata pengelolaan
rule Pega: `<pyDelete>` (false/true), `<pyVersionSecure>` (true/false), `<pzIsPrivateCheckOut>`
(true/kosong). **Urutan 20 langkahnya identik** dan **daftar rule yang dirujuk identik**. Jadi
konflik OQ-011 #415 adalah **artefak check-out/versi, bukan percabangan perilaku**.

**(e) Perilakunya kini terbaca — untuk varian Endorsment Fac In saja.** Dua puluh langkah:
`Obj-Refresh-And-Lock` → `Property-Set` → `RDB-List` → `Property-Set` →
`Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService` → `Connect-REST` ×2 → `Property-Set` →
`Obj-Save` → `Page-Set-Messages` → `RDB-List` ×2 → `Property-Set` ×2 → `RDB-List` ×2 →
`Page-Remove` → `Call ASMForceCaseClose` → `Property-Set` → `Call SendEmailWithAttachments`.
Rule yang dirujuk: `ASM GetPolicyNoByCaseId`, `ASM INSERTJSON_JSONPOLISMONITORING_FACIN`,
`ASM UpdateErrorNoteJsonPolisMonitoring`, **`RNM CekSTSKonversiJson`**,
**`RNM DeleteDataProduction`**, dan `@ASM.GetPageJSONString()`.

**Status: tetap terbuka.** Yang tersisa adalah keputusan **versi mana yang aktif di production** —
milik pemilik export, bukan sesuatu yang dapat disimpulkan dari korpus.


### OQ-026 — Satu nama rule ada sebagai dua tipe berbeda; mana yang dipakai flow

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2, dua bentuk berbeda:

| Bentuk | Contoh | Bukti |
| --- | --- | --- |
| Dua tipe **dalam satu modul** | `ASM-FW-GISFW-WORK / ISAPPROVED` ada sebagai `RULE-OBJ-WHEN` **dan** `RULE-DECLARE-DECISIONTABLE` | `NB Treaty In/When/isApproved.xml` + `NB Treaty In/DecisionTable/isApproved.xml` |
| Tipe berbeda **antar modul** | `IsFlagOnGoingPolicy` = `DecisionTable` di PremiumList Life, `When` di NB FacIn / RNW Fac In / Endorsment Fac In | `PremiumList Life/DecisionTable/IsFlagOnGoingPolicy.xml` vs `NB FacIn/When/IsFlagOnGoingPolicy.xml` |

- **Pertanyaan:** Ketika shape Decision merujuk nama itu, tipe mana yang di-resolve Pega? Apakah
  keduanya aktif, atau salah satu peninggalan?
- **Memblokir:** guard percabangan tidak dapat dinyatakan pasti. Di NB Treaty In, **enam** shape
  Decision merujuk `isApproved`.
- **Status:** terbuka
- **Catatan:** isi rule `DecisionTable` **tidak terekstraksi** oleh pembacaan tag sederhana —
  perlu metode baca tersendiri di D2 lanjutan.

**Diperbarui STEP D2 Tahap 4 — dipersempit `[terverifikasi]`.**

Pola yang sama ditemukan lagi, kali ini pada `ISUWACCEPTED`: `NB FacIn/When/IsUWAccepted.xml`
(`RULE-OBJ-WHEN`, kondisi `[.ProposalAcceptStatus][=]["1"]`) **dan**
`NB FacIn/DecisionTable/IsUWAccepted.xml` (`RULE-DECLARE-DECISIONTABLE`) — identitas sama,
`ASM-FW-GISFW-WORK!ISUWACCEPTED`.

**Dua bukti yang mempersempit pertanyaan:**

1. Connector flow membawa **enam** nilai hasil (`confirm`, `reject`, `ask`, `banding`, `revise`,
   `decline`) yang persis sama dengan `<pyTaskStatusXml>` milik **`DecisionTable`**. Sebuah
   `RULE-OBJ-WHEN` hanya dapat mengembalikan benar/salah, sehingga `[dugaan]` kuat bahwa yang
   dipakai flow adalah varian `DecisionTable`.
2. **`RNW Fac In` hanya punya `DecisionTable/IsUWAccepted.xml`**, tanpa varian `When` — namun
   flow-nya memakai keenam hasil yang sama.

**Belum tertutup:** flow tetap tidak menyebut tipe rule, dan isi baris tabelnya tidak ada di ekspor
(**OQ-043**), sehingga hasil yang dikembalikan untuk tiap nilai `ProposalAcceptStatus` tetap
**tidak dapat dinyatakan**.


### OQ-027 — `OperatorID.pyTelephone` dipakai menyimpan kode peran

- **Pemilik:** IAM
- **Ditemukan di:** STEP D2 — `NB Treaty In/When/IsTreaty1.xml` menguji
  `OperatorID.pyTelephone = "TREATY1"`; `When/IsSPVTreaty1.xml` menguji `= "SPVTREATY1"`.
  Sapuan korpus: **14 file di 7 modul**, hanya dua nilai (`TREATY1` 6×, `SPVTREATY1` 6×)
- **Pertanyaan:** Apakah field telepon operator memang dipakai sebagai penyimpan kode peran? Siapa
  yang mengisinya, dan apa daftar kode yang sah?
- **Memblokir:** desain RBAC. Ini **berbeda** dari OQ-021 (hardcode nama orang): di sini ada konsep
  peran, tetapi disimpan di field yang bukan peruntukannya.
- **Status:** terbuka
- **Audit:** `grep -rl "pyTelephone" "<modul>" --include="*.xml"`
- **Terkait:** OQ-007, OQ-021, OQ-024.

**Diperbarui STEP D2 Tahap 5 `[terverifikasi]` — dua nilai kode peran baru.**

`Treaty In/DataTransform/Akseptasi_DT.xml` menguji `OperatorID.pyTelephone` dengan **empat** nilai:

```
"TREATY1"   "TREATY2"   "SPVTREATY1"   "SPVTREATY2"
```

D1 dan D2 Tahap 1 hanya menemukan `TREATY1` dan `SPVTREATY1`. **`TREATY2` dan `SPVTREATY2` baru
ditemukan di sini**, dan keduanya menggerbangi cabang masuk tangga persetujuan treaty inward
(seluruhnya menuju `Position := "ReasTreatyInSecHead"`, `StatusAkseptasi := "Accept"`).

Jadi field telepon operator menyimpan **setidaknya 4 kode peran**, bukan 2.


### OQ-028 — `WorkList` vs `WorkBasket`: dua model penugasan

- **Pemilik:** Product+Underwriting + IAM
- **Ditemukan di:** STEP D2 — seluruh Assignment di NB Treaty In (6) dan EDM Treaty In (4) memakai
  `<pyImplementation>WorkBasket`; seluruh Assignment di Endorsement Life (2) dan PremiumList Life (3)
  memakai `<pyImplementation>WorkList`
- **Pertanyaan:** Apakah perbedaan ini disengaja — treaty memakai antrean bersama sedangkan Life
  memakai antrean per-pengguna? Bila ya, siapa yang menerima tugas di jalur Life?
- **Memblokir:** model penugasan dan notifikasi di sistem baru.
- **Status:** terbuka

**D2 Tahap 2 `[terverifikasi]`:** keempat modul Komite memakai **`WorkList`**, tidak satu pun `WorkBasket`. Sejauh D2: treaty inward = `WorkBasket`; Life, PremiumList, dan seluruh Komite = `WorkList`.

### OQ-029 — `IsPEGAPROD` menggerbangi integrasi keluar, kondisinya tidak terbaca, isinya berkonflik

- **Pemilik:** IAM + DBA + Product+Underwriting
- **Ditemukan di:** STEP D2 — `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` ada di **14 modul** dan
  terdaftar di **register OQ-011 entri #305 dengan 2 isi berbeda**.
  Di `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` ia menjadi precondition **tiga** langkah
  sekaligus: `call @baseclass.SendEmailNotification`, `call serviceInsertArasapasLife_act`, dan
  `Connect-REST`. Di `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` ia menggerbangi
  pemanggilan `serviceInsertArasapasLife_act`
- **Pertanyaan:** Apa yang diuji `IsPEGAPROD`? Mengapa isinya berbeda di dua kelompok modul? Apakah
  benar ia memeriksa lingkungan (production vs non-production)?
- **Memblokir:** seluruh pemahaman perilaku integrasi keluar — **email, service eksternal, dan
  REST semuanya di belakang gerbang ini**. Bila ia memang gerbang lingkungan, maka perilaku sistem
  **berbeda antar lingkungan**, dan dua isi berbeda berarti gerbang itu **tidak konsisten**.
  Menyentuh langsung **OQ-018** (korpus mungkin dari dev).
- **Status:** terbuka — **terjawab untuk Claim — Life** (2026-09-14, work owner): `IsPEGAPROD` diganti flag `ENV`, `IsSendtoMedical` artinya dikonfirmasi. Empat keluarga `When` lain + varian `IsPEGAPROD` di 13 modul tetap terbuka.
- **Catatan:** kondisi rule **tidak terbaca** dari tag — `<pyLabel>` hanya berisi template kosong
  `[first value][relation][second value]`. **Tidak boleh disimpulkan dari namanya.**

**D2 Tahap 2 `[terverifikasi]`:** `IsPEGAPROD` ada di keempat modul Komite. Di `Komite Claim Life/Activity/KomitePostAdjustment.xml` ia menggerbangi **tiga** pemanggilan keluar berturut-turut: `serviceInsertArasapasClaimLife_act`, `SendEmailKlaimLife`, `HitServiceToKasirKMTLife_Act` — pola identik dengan PremiumList Life di Tahap 1. Kondisinya tetap **tidak terbaca**; cabang aktif **tidak ditebak**.

**D2 Tahap 3 `[terverifikasi]`:** pola When tak terbaca bertambah **dua kasus**, keduanya menggerbangi alur: `Claim Life/When/IsSendtoMedical.xml` (pengembalian ke tahap medis) dan `Claim Fac In/When/IsSPK.xml` (**percabangan di titik masuk flow**, jalur B2B). Ditambah keduabelas varian `ISCLM*` (OQ-041). Jadi ada **tiga keluarga rule** yang kondisinya tidak dapat dibaca dari tag: `IsPEGAPROD`, `IsSendtoMedical`/`IsSPK`, dan `ISCLM*`.

**Diperbarui STEP D2 Tahap 5 `[terverifikasi]` — keluarga kelima, dan catatan cakupan.**

`Treaty In Adjustment` memuat dua rule `When` yang `<pyLabel>`-nya hanya berisi template kosong
`[first value][relation][second value]`:

| Rule | Kondisi kosong | Catatan |
| --- | ---: | --- |
| `When/IsTreatyUser.xml` | **6** | namanya menyangkut otorisasi pengguna treaty; **apa yang diujinya tidak dapat dinyatakan** |
| `When/recordEvent.xml` | 1 | `<pyLabel>` = "Record User Waiting Time Events" |

Keduanya **hanya ada di modul itu**.

`[terverifikasi]` **`IsPEGAPROD` tidak ada di kelima modul Tahap 5**
(`find "<modul>" -iname 'IsPEGAPROD.xml' | wc -l` → 0 untuk kelimanya), sehingga cabang
production/non-production tidak menjadi persoalan di kelompok ini.


**TERJAWAB untuk konteks Claim — Life — 2026-09-14, sumber: work owner** (grilling Ronde 1 Q6).

`[terverifikasi work owner]`

| Rule | Keputusan / arti |
| --- | --- |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` (`Claim Life/When/IsPEGAPROD.xml`) | **Jangan ditiru sebagai rule.** Diganti flag lingkungan eksplisit (`ENV=production`) yang menggerbangi tiga efek keluar yang sama: unggah berkas, simpan utama, email. → **ADR** |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` (`Claim Life/When/IsSendtoMedical.xml`) | `SendtoMedical = 1`: dari **`ReasLifeSPV`**, kasus dikembalikan ke **`ReasLifeMedicalAdvisor`** |

`[terverifikasi]` Tiga efek keluar yang digerbangi `IsPEGAPROD` di modul ini:
`Claim Life/Activity/InsertGoogleStorage_Act.xml`,
`Claim Life/Activity/SaveOutStandingLife_Act.xml` (642.787 byte),
`Claim Life/Activity/SendEmailKlaimLF.xml`.

**Cakupan korpus-wide TETAP TERBUKA:** empat keluarga `When` tak terbaca lainnya —
`IsSPK` (Claim Fac In), `ISCLM*` 12 varian (**OQ-041**), `IsPKSASM`/`LetterNoNull`/`ToJUW_A`/`ToUW`/
`IsEdmInternalRetro` (facultative), `IsTreatyUser`/`recordEvent` (Treaty In Adjustment).
Varian `IsPEGAPROD` di 13 modul lain juga belum dinyatakan.


### OQ-030 — Ambang tanggal 25 muncul di tiga konteks

- **Pemilik:** Finance + Product+Underwriting
- **Ditemukan di:** STEP D2:

| Konteks | Bentuk | Path |
| --- | --- | --- |
| EDM Treaty In | `@if(Local.TglProd=="",25,Local.TglProd)` — **default** bila query kosong | `Activity/SaveJsonPolisTreatyInEDM_Act.xml` |
| Endorsement Life | `@toDecimal(Local.currentdate) > 25` | `Activity/InsertJsonPolisLife_Act.xml` |
| PremiumList Life | `@toDecimal(Local.currentdate) > 25` | `Activity/InsertJsonPolisLife_Act.xml` |

  Di EDM Treaty In nilainya seharusnya berasal dari `RDBList` `GETTanggalClosing_SQL`, tetapi
  **25 dipakai sebagai fallback hardcoded**.
- **Pertanyaan:** Apakah 25 adalah hari tutup buku bulanan? Mengapa ia di-hardcode sebagai default
  padahal ada tabel sumbernya? Apakah nilai itu pernah berubah?
- **Memblokir:** aturan penentuan periode produksi/akuntansi. Angka ini menentukan ke bulan mana
  sebuah transaksi dibukukan.
- **Status:** terbuka — **TERTUTUP untuk PremiumList Life** (2026-09-15, sumber: **korpus +
  keputusan work owner**).

  `[terverifikasi]` **Ambang tutup buku BUKAN hardcode `25`.**
  `PremiumList Life/Activity/SubmitPremiumList_Act.xml`
  (`ASM-FW-GISFW-WORK-LIFE` / `SUBMITPREMIUMLIST_ACT` / `RULE-OBJ-ACTIVITY`, 214.154 byte,
  tersimpan `20260122T072213`):

  | Baris | Isi |
  | ---: | --- |
  | 773 | memanggil `GETTanggalClosing_SQL` — `SELECT * FROM POOLDATA.TANGGAL_CLOSING` |
  | 918-919 | `Local.TglProd = TglProd.pxResults(1).TANGGAL` |
  | **966** | `@if(Local.TglProd=="", 25, Local.TglProd)` — **`25` hanya FALLBACK bila tabel kosong** |
  | 1211 | `@if(@toDecimal(Local.currentdate)>Local.TglProd, @toDecimal(Local.CurrentMonth)+1, ...)` |
  | 3491 | `@toDecimal(Local.currentdate) > Local.TglProd` |

  **Aturan:** transaksi setelah tanggal closing dibukukan ke periode **bulan berikutnya** —
  tanggal 1, `05:00 GMT` (= 12:00 WIB).

  ⚠️ `[terverifikasi]` **Dua versi hidup berdampingan:**

  | Rule | Ambang | Tersimpan |
  | --- | --- | --- |
  | `SubmitPremiumList_Act` | dari `TANGGAL_CLOSING` | `20260122T072213` |
  | `InsertJsonPolisLife_Act` | **hardcode `>25`** (baris ~1170) | `20260728T024422` |

  Yang lebih **baru** justru yang hardcode. `[keputusan work owner]` **Ikuti yang dari DB** —
  sistem baru membaca `TANGGAL_CLOSING`/konfigurasi; **jangan tanam `25`**.

  `[data DBA]` `POOLDATA.TANGGAL_CLOSING` — `TANGGAL VARCHAR2(10)`, **string bukan `DATE`**
  (OQ-001 untuk Claim Life).

  Cakupan modul lain yang memakai pola ambang tanggal serupa tetap terbuka.

### OQ-031 — Empat ID literal `1000032`–`1000035` menggerbangi logika

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 — `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` step 7–10
  memakai `@contains(.ID,"1000032")`, `"1000033"`, `"1000034"`, `"1000035"` sebagai precondition
  untuk empat `Property-Set` berbeda
- **Pertanyaan:** Entitas apa yang diwakili keempat ID itu (produk? mata uang? jenis premi?) dan
  dari tabel mana asalnya?
- **Memblokir:** empat cabang logika di jalur simpan premium list tidak dapat dijelaskan. Nilai
  ter-hardcode seperti ini juga tidak boleh dipindahkan apa adanya ke sistem baru.
- **Status:** terbuka — **TERTUTUP untuk PremiumList Life** (2026-09-15, sumber: **keputusan work
  owner + korpus**).

  `[keputusan work owner]` Keempat gerbang `@contains(.ID,"1000032")` … `"1000035"` di step 6-7
  `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` adalah **logika POLIS LAMA yang sudah
  tidak dipakai** — kemungkinan kelupaan di-remark. **JANGAN direplikasi di sistem baru.**

  ⚠️ **Perhatian:** di korpus keempat langkah itu **`blockname` masih kosong = AKTIF**. Ia
  **dead code yang masih hidup** — lihat **OQ-066**.

  `[terverifikasi]` Class sumbernya `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`; label langkahnya
  **QS / 2nd QS / SURPLUS / 2nd SURPLUS**.

  `[terverifikasi]` **Arti kode jenis treaty terbaca penuh** di berkas yang sama (baris ~3787):

  ```
  @if(pyWorkPage.TypeCeding="1","QS",
   @if(pyWorkPage.TypeCeding="2","SURPLUS",
    @if(pyWorkPage.TypeCeding="3","QS + SURPLUS",
     @if(pyWorkPage.TypeCeding="4","XOL",""))))
  ```

  → **`TypeCeding`: 1 = QS (Quota Share), 2 = SURPLUS, 3 = QS + SURPLUS, 4 = XOL.**
  Bila sistem baru butuh jenis treaty, ambil dari **data `TypeCeding`**, bukan konstanta ID.

  Ini juga meringankan glossary: `QS` yang semula "kepanjangan belum terverifikasi"
  (`discovery/glossary.md` baris 216) kini punya padanan di jalur Life.

  Pola ID literal di modul lain (mis. `1000013` di Komite Claim Life — sudah dibuang, OQ-064) tetap
  terpisah.

**D2 Tahap 2 `[terverifikasi]`:** pola ID literal muncul lagi — `Komite Claim Life/Activity/KomitePostAdjustment.xml` menggerbangi `Call InsertJsonClaimLife_Act` dengan `pyWorkCover.ClaimData.PolicyDataLife.RetroID=="1000013"`. Rentang `1000013` dan `1000032`–`1000035` menyiratkan satu tabel referensi bersama; **belum terverifikasi**.

### OQ-032 — Sumber nilai `.KomiteLoop` (jumlah tingkat tangga komite) tidak ditemukan

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 2 — `.KomiteLoop` menentukan **berapa tingkat** tangga persetujuan
  komite (guard `IsKomiteLoop`: `.AcceptStatus = "1"` DAN `.KomiteCount <= .KomiteLoop`), tetapi
  tidak ada rule di keempat modul Komite yang **mengisinya** — seluruhnya hanya membaca
- **Pertanyaan:** Di mana `.KomiteLoop` di-set, dan apa yang menentukan nilainya (nominal klaim,
  jenis produk, lini bisnis)?
- **Memblokir:** jumlah tingkat persetujuan tidak dapat dinyatakan. Tanpa ini, tangga komite tidak
  dapat direplikasi.
- **Catatan `[terverifikasi]`:** Komite Claim FacIn memakai properti **berbeda** untuk konsep yang
  sama — `Local.TotalKomite`, bukan `.KomiteLoop`.
- **Status:** **TERTUTUP** (2026-09-14, sumber: **korpus**). Tangga komite **sepenuhnya
  data-driven** — tidak ada konstanta di mana pun.

  `[terverifikasi]` **Tiga bukti berantai:**

  1. `Claim Life/Activity/GetListKomiteLife.xml`
     (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) menetapkan
     `Local.IsADj = .CLAIM_AMOUNT`, lalu memanggil report `FilterEmailKomiteWithLimit`
     (`Param.pyReportClass = "ASM-FW-GCNMFW-Int-EMAILKOMITE"`) dan menampung hasilnya ke
     `Primary.KomiteList(<APPEND>)`.
  2. `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml`
     (`ASM-FW-GCNMFW-INT-EMAILKOMITE` / `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`)
     memfilter `EMAILKOMITE` dengan
     `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM AND .STS_KLAIM = Param.STS_KLAIM AND .STS_AKTIF = "1"`.
  3. `Claim Life/Activity/CreateKMTLife_Act.xml`
     (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) menetapkan
     `childPageKomite.KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)`.

  **Kesimpulan:** **jumlah tingkat komite = COUNT baris roster aktif yang `LIMIT_BOTTOM <=
  CLAIM_AMOUNT`.** Ladder mencapai final saat `KomiteCount == KomiteLoop` — `[terverifikasi]`
  seluruh penulisan status di `Komite Claim Life/Activity/KomitePostAdjustment.xml`
  (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`) digerbangi
  perbandingan itu.

  `[terverifikasi]` **Detail yang perlu direplikasi:** ambang yang dikirim ke report adalah **nilai
  mutlak** klaim —
  `Param.LIMIT_BOTTOM = @if(Local.IsADj<0, Local.IsADj * -1, Local.IsADj)`. Nilai klaim negatif
  dihilangkan tandanya sebelum pencarian roster.

  Dicatat di **ADR-0001** (kontrak penyerahan) dan **ADR-0012**.
  Claim — Life**, bukan oleh Komite; nilainya menyeberang sebagai `childPageKomite.KomiteLoop` di
  `Claim Life/Activity/CreateKMTLife_Act.xml`. Yang **masih terbuka**: apa yang menentukan nilainya
  (nominal klaim, produk, lini bisnis).
- **Bukti tambahan `[terverifikasi]` (2026-09-14):** `STS_REJECT` hanya berubah pada **tingkat
  terakhir** — `Komite Claim Life/Activity/KomitePostAdjustment.xml`
  (`ASM-FW-GCNMFW-WORK-KOMITELIFE!KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`) menggerbangi seluruh
  penulisannya dengan `pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop`. Jadi nilai `KomiteLoop`
  langsung menentukan **kapan** keputusan menjadi final, bukan sekadar berapa tingkat.
- **Audit:** `grep -rl "KomiteLoop" "<modul Komite>" --include="*.xml"`

### OQ-033 — `KomiteRouter` hanya menetapkan sasaran bila `TransferType == '2'`

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 2 — `Komite Claim Life/Activity/KomiteRouter.xml` berisi **satu**
  langkah: `param.AssignTo = .KomiteID`, precondition `Primary.TransferType == '2'`
- **Pertanyaan:** Apa yang terjadi bila `TransferType != '2'` — `param.AssignTo` tidak di-set sama
  sekali. Ke mana tugas dirutekan? Apa arti nilai `'2'`?
- **Memblokir:** perilaku routing di jalur non-`'2'` tidak dapat dinyatakan.
- **Catatan:** FacIn, Prop, dan Non Prop **punya** cabang eksplisit untuk `TransferType != '2'`;
  hanya Komite Claim Life yang tidak.
- **Status:** **TERTUTUP untuk Komite Claim Life** (2026-09-15, sumber: **korpus + keputusan work owner**)
  — **`TransferType` adalah dead code.**

  `[terverifikasi]` `TransferType` **tidak pernah diisi** di korpus: **nol `Property-Set`** terhadapnya
  di `Claim Life` maupun `Komite Claim Life`. Ia hanya **dibaca** di dua tempat:

  | Tempat | Bentuk | Baris |
  | --- | --- | ---: |
  | `Komite Claim Life/Activity/KomiteRouter.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`) | precondition `Primary.TransferType=='2'` | ~442 |
  | `Komite Claim Life/Section/ShowTransfer.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION`) | visible-when `.TransferType==2` | ~29177 |

  Karena tak pernah di-set, `TransferType=='2'` **selalu FALSE** → cabang mati, sisa salinan dari
  modul Komite lain.

  `[terverifikasi]` **Routing tingkat sesungguhnya digerakkan precondition pertama** `KomiteRouter`:
  `.KomiteAproval == 0` (baris ~382), dengan sasaran `param.AssignTo = .KomiteID` (baris ~294).
  Itulah sebabnya graf cukup satu Assignment yang di-loop: tugas selalu jatuh ke **baris roster
  pertama yang belum menyetujui**.

  `[keputusan work owner]` **Buang `TransferType`, dan buang pula UI `ShowTransfer` yang bergantung
  padanya.** Routing sistem baru = cari baris roster pertama ber-`KomiteAproval == 0`.

  Cakupan modul Komite lain (FacIn / Prop / Non Prop) tetap terbuka — mereka punya cabang eksplisit
  yang Komite Life tidak punya.

### OQ-034 — Tanggal cutover ter-hardcode `20250207T000000.000 GMT`

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 2 — `Komite Claim Life/Activity/KomitePostAdjustment.xml`,
  dua `Property-Set` digerbangi
  `TempOpenPage.PolicyDataLife.OfferFacIn.PolicyData.ProdDateTime < "20250207T000000.000 GMT"`
- **Pertanyaan:** Aturan bisnis apa yang berubah pada 7 Februari 2025? Apakah data lama tetap harus
  diperlakukan dengan aturan lama setelah migrasi?
- **Memblokir:** perilaku berbeda untuk data sebelum/sesudah tanggal itu harus dipertahankan atau
  sengaja dihapus — keputusan yang tidak bisa diambil tanpa tahu sebabnya.
- **Status:** **TERTUTUP untuk Komite Claim Life** (2026-09-15, sumber: **keputusan work owner**).

  `[keputusan work owner]` Blok precondition
  `TempOpenPage.PolicyDataLife.OfferFacIn.PolicyData.ProdDateTime < "20250207T000000.000 GMT"`
  di `Komite Claim Life/Activity/KomitePostAdjustment.xml` (baris ~5135 dan ~7861) **sudah
  di-remark / tidak dipakai lagi**. **Jangan direplikasi** di sistem baru.

  `[terverifikasi]` Gerbang itu adalah **alternatif ketiga dari tiga** yang mengisi identitas retro
  (`TempInputDetail.CARI38/CARI39` = `RetroID` / `RetroName`) sebelum panggilan SQL. Dua alternatif
  yang **tetap aktif**:

  1. `TempOpenPage.PolicyDataLife.Type=="TP" || =="TR"`
  2. `TempOpenPage.PolicyDataLife.SecurityReinsurerID!="" && SecurityReinsurer!=""`

  Jadi pengisian identitas retro cukup dua precondition itu. Tanggal absolut 7 Feb 2025 tidak ikut
  pindah.

  Cakupan modul Komite lain tetap terbuka bila polanya ada di sana.

### OQ-035 — Activity dipanggil tetapi salinannya hanya ada di modul lain

- **Pemilik:** pemilik export Pega
- **Ditemukan di:** STEP D2 Tahap 2:

| Dipanggil dari | Activity | Satu-satunya salinan di korpus |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `UpdateWorkObject` | `Endorsment Fac In/Activity/UpdateWorkObject.xml` |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `serviceInsertArasapasClaimLife_act` | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` |
| `Komite Claim Prop/Activity/KomitePost_Reject.xml` | `UpdateWorkObject` | idem |

- **Pertanyaan:** Apakah rule ini di-resolve lewat pewarisan class Pega (sehingga salinan di modul
  lain memang rule yang sama), atau ekspor per-modul memang tidak lengkap?
- **Memblokir:** perilaku langkah tersebut tidak dinyatakan. **Berbeda dari OQ-025**: rule ini
  **tidak** berkonflik, jadi bukan blocker isi — hanya ketidakpastian resolusi.
- **Status:** terbuka — **bukan pemblokir isi** untuk Komite Claim Life (2026-09-15). Grilling Ronde 1
  menelusuri `KomitePostAdjustment` step demi step; `Call UpdateWorkObject` (step 5) dan
  `Call serviceInsertArasapasClaimLife_act` (step 9) **terpanggil dan terindeks aktif**. Yang belum
  diketahui hanyalah **mengapa salinannya ada di modul lain** — pewarisan class atau ekspor tidak
  lengkap. Perilakunya sudah terbaca; kepemilikan rule-nya belum. Pemilik: Arsitektur Pega.

### OQ-036 — Empat sasaran routing komite ter-hardcode (`komitepnc`…`komitepnc4`)

- **Pemilik:** IAM + Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 2 — `Activity/KomiteRouter.xml` di **tiga** modul
  (Komite Claim FacIn, Komite Claim Prop, Komite Claim Non Prop) menetapkan
  `param.AssignTo` ke string ter-hardcode menurut `.KomiteCount`:

| `.KomiteCount` | `param.AssignTo` |
| ---: | --- |
| 1 | `"komitepnc"` |
| 2 | `"komitepnc2"` |
| 3 | `"komitepnc3"` |
| 4 | `"komitepnc4"` |

  `Komite Claim Life` **tidak** memakai pola ini — sasarannya murni data (`.KomiteID`).
- **Pertanyaan:** Apakah `komitepnc`…`komitepnc4` adalah workbasket, operator, atau access group?
  Siapa anggotanya? Apa kepanjangan `PNC`? Mengapa maksimal empat tingkat?
- **Memblokir:** desain RBAC dan model penugasan. Nilai ter-hardcode ini **tidak boleh dipindahkan
  apa adanya** ke sistem baru.
- **Status:** terbuka
- **Terkait:** OQ-007, OQ-024, OQ-027.

### OQ-037 — Ambang nominal ter-hardcode menentukan komposisi roster komite

- **Pemilik:** Finance + Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 2 — `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml`
  (146.689 byte) memuat precondition:

```
Local.TotalAdj > 30000000.00 && Local.TotalAdj <= 57750000.00
```

  yang menggerbangi `Obj-Browse` terhadap class `ASM-FW-GCNMFW-Int-EMAILKOMITE` — artinya
  **komposisi roster komite bergantung pada besaran nilai**. `Local.TotalAdj` diakumulasi dari
  `Local.TotalAdj + .ValueAdjustment`.
- **Pertanyaan:** Apa mata uangnya (tidak disebut di rule)? Berapa pita wewenang seluruhnya?
  Siapa yang berwenang di tiap pita? Kapan terakhir diperbarui?
- **Memblokir:** batas wewenang persetujuan — bagian paling material dari tangga komite. Angka
  ter-hardcode juga berarti perubahan limit memerlukan perubahan kode.
- **Status:** terbuka — **TERTUTUP untuk Claim — Life** (2026-09-14, sumber: **korpus + data DBA**).

  `[terverifikasi]` Berbeda dari Komite Claim FacIn yang memakai pita **ter-hardcode**, jalur Life
  membaca **tabel**. `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml`
  (`ASM-FW-GCNMFW-INT-EMAILKOMITE` / `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`)
  memfilter tabel `EMAILKOMITE` dengan
  `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM AND .STS_KLAIM = Param.STS_KLAIM AND .STS_AKTIF = "1"`.

  `[data DBA]` Kolom pita pada `POOLDATA.EMAILKOMITE`: `LIMIT_BOTTOM INTEGER`, `LIMIT_TOP INTEGER`,
  `DEGREE VARCHAR2(150)`, `TYPE_KOMITE INTEGER`, `STS_AKTIF VARCHAR2(150)`, `STS_KLAIM VARCHAR2(10)`,
  `STS_REJECT VARCHAR2(15)`, PK `ID INTEGER`.

  ⚠️ `[data DBA]` **Roster tidak punya kolom mata uang** — pita nilai karena itu berlaku atas
  **satu mata uang implisit**. Bila kelak ada klaim bermata uang lain, pembandingan pita menjadi
  tidak sahih. Dicatat sebagai risiko, bukan sebagai pertanyaan terbuka.

  ⚠️ `[data DBA]` `STS_REJECT` di tabel roster bertipe **`VARCHAR2(15)`** — **berbeda** dari
  `STS_REJECT NUMBER(38)` di tabel klaim. Nama sama, tipe berbeda, tabel berbeda.

  **Aturan sistem baru:** ambang adalah **data operasional di tabel**; query roster saat dibutuhkan,
  **jangan tanam ambang di kode**.

  Cakupan modul lain (mis. Komite Claim FacIn yang hardcode) tetap terbuka.

### OQ-038 — Sebelas kode lini bisnis `L1`…`L11` ter-hardcode dalam satu precondition

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 3 — `Claim Life/Activity/SaveOutStandingLife_Act.xml` memuat satu
  precondition yang menguji `pyWorkPage.PolicyDataLife.BusinessCode` terhadap `"L1"`, `"L2"`, …
  `"L11"` berderet dengan `||`
- **Pertanyaan:** Apa arti tiap kode `L1`–`L11`? Adakah tabel referensinya? Apakah daftar ini pernah
  bertambah?
- **Memblokir:** cabang logika di jalur outstanding Life. Daftar ter-hardcode berarti penambahan
  lini bisnis memerlukan perubahan rule.
- **Status:** **TERJAWAB** (2026-09-14, sumber: work owner) — daftar `L1`–`L21` lengkap di `CONTEXT.md`. Menyisakan catatan migrasi (daftar ter-hardcode), bukan pertanyaan.

**TERJAWAB — 2026-09-14, sumber: work owner** (grilling Ronde 1 Q5).

`[terverifikasi work owner]` Daftar lengkap **`L1`–`L21`**: kode produk life → nama produk →
jenis klaim (`ContentNote`). Tabel penuh ada di **`CONTEXT.md`** §"Daftar BusinessCode".

Ringkas: `L1`–`L11` → `DEATH`; `L12`–`L15`, `L18` → `HEALTH`; `L16`, `L19` → `CI`;
`L17`, `L20` → `TPD`; `L21` → `TI`.

`[terverifikasi]` Modul `Claim Life` **hanya menguji `L1`–`L11`** (seluruhnya `DEATH`) dalam satu
precondition di `Claim Life/Activity/SaveOutStandingLife_Act.xml`. `L12`–`L21` berada di master
produk, bukan di modul klaim.

**Menyisakan (catatan migrasi, bukan OQ):** daftar ter-hardcode dalam satu precondition berarti
penambahan lini produk memerlukan perubahan rule. Perlakuan di sistem baru adalah keputusan
desain, bukan pertanyaan terbuka.


### OQ-039 — Adjustment, Close, dan Reject tidak muncul di graf flow

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 3 — keempat modul Claim tidak memuat shape untuk Adjustment,
  Close, maupun Reject, padahal Activity-nya ada (`SaveAdjustment_Act`, `ProtectCloseClaim_act`,
  `RejectOSClaimLife_Act`, `AdjClaimCNP_Act`, `SaveRejectOSKomiteCNP`)
- **Pertanyaan:** Bagaimana ketiga tahap itu dipicu — dari menu, harness, atau FlowAction di luar
  graf? Apakah ada flow lain yang tidak terekspor?
- **Memblokir:** siklus kerugian tidak dapat digambarkan utuh dari flow saja.
- **Status:** terbuka — **TERJAWAB untuk Claim — Life** (2026-09-14, sumber: work owner, Ronde 4).
  Tetap terbuka untuk **3 modul Claim lain** (`Claim Prop`, `Claim Non Prop`, `Claim Fac In`), yang
  pemicu Adjustment/Close/Reject-nya belum dinyatakan.

  Riwayat penyempitan:
  **Ronde 2** — lokasi penyerahan ke Komite terverifikasi; tersisa aturan bisnis kapan dipicu.
  **Ronde 3** — `[keputusan work owner]` ada **dua jalur reject**: (a) Admin menolak adjustment yang
  ia input sendiri, tanpa Komite; (b) Komite menolak lewat SPV. Aksep final selalu lewat Komite.
  `Claim Life/Activity/RejectOSClaimLife_Act.xml` **bukan** dead rule — ia jalur reject milik Admin.
  **Ronde 4** — tiga butir sisa dijawab:

  1. **Akibat reject oleh Admin** `[keputusan work owner]`: **membatalkan baris `AdjustmentList` itu
     saja**. Klaim **tidak** tertutup; Admin lalu menginput baris adjustment baru.
  2. **Mekanismenya** `[keputusan work owner]`: memakai **`STS_REJECT = 2` yang sama** dengan reject
     Komite. Konsekuensinya dicatat di **ADR-0011** §"Arti tunggal nilai `2`" — nilai `2` **selalu**
     berarti "baris ini ditolak", **tidak pernah** "klaim selesai"; konsisten dengan klaim yang tidak
     terminal dan terminal hanya per baris.
  3. **`SaveAdjustment_Act`** `[keputusan work owner]`: **dead rule, jangan dimigrasikan.** Penulis
     `STS_REJECT = 1` yang berlaku adalah `KomitePostAdjustment` di sisi Komite.
     ⚠️ Status "dead" ini **tidak dapat diverifikasi dari korpus** — yang terbaca justru sebaliknya
     (`[terverifikasi]` rule itu terpasang di UI, dirujuk 2× dari
     `Claim Life/Section/ClaimLifeDetailGCNM.xml`). Perbedaan itu dicatat sengaja di **ADR-0011**
     §`SaveAdjustment_Act` sebagai titik mula bila keputusan ini perlu ditinjau ulang.

  **Kapan penyerahan dipicu** juga terjawab untuk Claim — Life: Ronde 3 menyatakan send-ke-Komite
  berlaku **untuk semua adjustment** — bukan hanya di atas ambang nilai tertentu.
- **Bukti tambahan `[terverifikasi]` (2026-09-14) — ada DUA jalur akseptasi berdampingan:**
  selain jalur Komite, `Claim Life/Activity/SaveAdjustment_Act.xml`
  (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SAVEADJUSTMENT_ACT` / `RULE-OBJ-ACTIVITY`, 179.221 byte)
  punya satu langkah `Property-Set` **tanpa precondition** yang menulis pada baris yang baru
  ditambahkan: `.AdjustmentList(<LAST>).ACCEPTEDNO = InputData.pxResults(1).HASIL1`,
  `.AdjustmentList(<LAST>).STS_REJECT = 1`, `.AdjustmentList(<LAST>).ACCEPTATION_DATE = @CurrentDateTime()`.
  Artinya Claim — Life **dapat mengaksep tanpa melalui Komite**.
- **Pertanyaan tambahan:** Kapan akseptasi boleh lewat `SaveAdjustment_Act` langsung, dan kapan
  wajib lewat Komite? Apakah ambang nilai yang memisahkan (kaitan OQ-037, OQ-040)?
- **Audit tambahan:**
  ```
  f="Claim Life/Activity/SaveAdjustment_Act.xml"
  grep -n "AdjustmentList(&lt;LAST&gt;).STS_REJECT" "$f"   # nilai ada di baris berikutnya
  ```

**D2 Tahap 3 — bagian yang terjawab `[terverifikasi]`:** penyerahan ke **Komite** dilakukan
**Activity**, bukan shape flow:

| Modul Claim | Activity pembuat case anak Komite | Ukuran |
| --- | --- | ---: |
| Claim Non Prop | `Activity/CreateChildKomiteCNP_Act.xml` | 756.836 byte |
| Claim Prop | `Activity/AddKomiteTreatyChild_ACT.xml` | 420.247 byte |
| Claim Life | `Activity/CreateKMTLife_Act.xml`, `Activity/GetListKomiteLife.xml` | belum diukur |

Keduanya di class **`ASM-FW-GCNMFW-DATA-ADJUSTMENT`** — bukan class modulnya sendiri — dan
digerbangi perbandingan nilai terhadap limit (OQ-040).

**Masih terbuka:** pemicu Adjustment, Close, dan Reject.

**DIPERSEMPIT — 2026-09-14, sumber: work owner** (grilling Ronde 1 Q1).

**Lokasi penyerahan ke Komite kini TERVERIFIKASI** untuk `Claim Life`, dan alasan ia tidak terlihat
di graf flow sudah jelas: **ia dipicu dari UI, bukan dari Flow.**

`[terverifikasi]` Diperiksa ulang ke korpus dan **cocok dengan pernyataan work owner**:

| Fakta | Bukti |
| --- | --- |
| Activity pembuat komite | `Claim Life/Activity/CreateKMTLife_Act.xml` → `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!CREATEKMTLIFE_ACT`, 121.652 byte |
| Membuat *child work* | langkah: `Property-Set`, `Call pxRetrieveReportData`, `Property-Set` ×3, **`Call pxAddChildWork`**, `Obj-Refresh-And-Lock`, `Property-Set`, **`Obj-Save`**, `Call SendEmailKlaimLF` |
| Kelas anak | `ASM-FW-GCNMFW-Work-KomiteLife` (15 kemunculan); `pyWorkPage` = `ASM-FW-GCNMFW-Work-ClaimLife` |
| Dipicu dari UI | `<pyActivity>CreateKMTLife_Act</pyActivity>` di `Claim Life/Section/ClaimComite.xml` (2×) dan `Claim Life/Harness/Committe_Life.xml` (2×) |
| Tambahan yang ditemukan | activity ini juga merujuk class roster **`ASM-FW-GCNMFW-Int-EMAILKOMITE`** dan memanggil `SendEmailKlaimLF` |

**Kontrak keluar Claim Life → Komite** (dicatat di `CONTEXT.md` dan kandidat ADR-0001):
> Claim Life menyerahkan kasus ke Komite dengan membuat *child work* berkelas
> `ASM-FW-GCNMFW-Work-KomiteLife` lewat `CreateKMTLife_Act`, dipicu dari layar
> `ClaimComite` / `Committe_Life`.

**YANG MASIH TERBUKA:** **kondisi/aturan bisnis KAPAN penyerahan dipicu** — otomatis di atas nilai
tertentu, atau keputusan manual? **Tidak ada di korpus.** **Pemilik: Product+UW.**

**Cakupan modul lain TETAP TERBUKA:** Adjustment, Close, dan Reject yang juga tidak muncul di graf
(`SaveAdjustment_Act`, `ProtectCloseClaim_act`, `RejectOSClaimLife_Act` di `Claim Life`; pola sama
di Claim Prop, Claim Non Prop, Claim Fac In) belum dijawab.


### OQ-040 — Limit wewenang: ter-hardcode di satu modul, dari database di modul lain

- **Pemilik:** Finance + Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 2 dan 3 — dua mekanisme hidup berdampingan:

| Modul | Mekanisme | Bukti |
| --- | --- | --- |
| Claim Non Prop | **ter-hardcode** | `Activity/CreateChildKomiteCNP_Act.xml`: `Local.LimitMax = 30000000.00`, `Local.LimitMaxDivHead = 50000000.00`, `Local.LimitPersenMax = 30.00` |
| Komite Claim FacIn | **ter-hardcode** | `Activity/ApprovalKomite_Act.xml`: `> 30000000.00 && <= 57750000.00` (OQ-037) |
| Claim Prop | **dari database** | `RDBList/GetLimitDirekturUtama_SQL.xml`, `GetLimitPLATreatyin.xml`, `GetLimitsTreatyIn_SQL.xml` |
| Claim Non Prop | **juga** dari database | `RDBList/GetLimitsTreatyIn_SQL.xml`, `GetLimitTONPPLA.xml` |

- **Pertanyaan:** Mana yang berlaku bila keduanya ada di modul yang sama (Claim Non Prop)? Mengapa
  batas bawah **30.000.000,00 sama persis** di dua tempat tetapi batas atasnya berbeda
  (50.000.000,00 vs 57.750.000,00)? Apa mata uangnya — tidak disebut di mana pun. Apa arti
  `LimitPersenMax = 30.00` (persen dari apa)? Apa kepanjangan `DivHead`?
- **Memblokir:** batas wewenang persetujuan — aturan paling material di siklus klaim. Nilai
  ter-hardcode juga berarti perubahan limit menuntut perubahan kode.
- **Status:** terbuka

### OQ-041 — `ISCLM`/`ISCLMP`/`ISCLMNP`/`ISPEGASYARIAH`: berkonflik 6-versi **dan** tidak terbaca

- **Pemilik:** Product+Underwriting + pemilik export Pega
- **Ditemukan di:** STEP D2 Tahap 3 — keempat `When` di `@BASECLASS` terdaftar di register OQ-011
  (#291, #292, #293, #306), masing-masing **6 varian dengan 6 isi berbeda**, di modul:
  Claim Fac In, Claim Non Prop, Claim Prop, Komite Claim FacIn, Komite Claim Non Prop,
  Komite Claim Prop. **Claim Life tidak memilikinya.**

  Hash varian yang diukur di Tahap 3:

| When | Claim Prop | Claim Non Prop | Claim Fac In |
| --- | --- | --- | --- |
| `ISCLM` | `2245b132` | `541179a0` | `4ee647f1` |
| `ISCLMP` | `4919199b` | `9b10dc97` | `a44d0c98` |
| `ISCLMNP` | `50c584af` | `d3255d2c` | `91c9b037` |
| `ISPEGASYARIAH` | `0a7077fd` | `1bbf8269` | `cac64576` |

- **Pertanyaan:** Apa yang diuji keempat rule ini, dan mengapa tiap modul punya versinya sendiri?
- **Memblokir:** **berlapis.** Kita tahu keenam varian berbeda (OQ-011), tetapi kondisinya
  **tidak terbaca dari tag** — `<pyLabel>` hanya berisi template kosong
  `[first value][relation][second value]`. Jadi perbedaannya **tidak dapat dilihat**, apalagi
  dinilai. Ini kombinasi terburuk antara OQ-011 dan OQ-029.
- **Status:** terbuka
- **Audit:** `grep -oE "<pyLabel>\[[^<]{1,80}" "<modul>/When/IsCLM.xml" | sort -u`

### OQ-042 — Modul klaim membaca master domain treaty (recovery/retrosesi?)

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 3:

| Modul | Objek treaty yang dibaca | Rule |
| --- | --- | --- |
| Claim Non Prop | `M_TREATY_OUT`, `M_TREATY_OUT_DETAIL`, `TREATY_OUT` | `RDBList/GetDataMasterTOutNP.xml`, `BrowseDtlTreatyOutNP.xml`, `GetLimitTONPPLA.xml` |
| Claim Fac In | `TREATYCONTRACT`, `TREATYBUSINESS`, `PROPORTIONALARRG` | 2–3 rule masing-masing; + `Activity/DLAFacintoTreaty_Act.xml` (644.505 byte) |

- **Pertanyaan:** Apakah ini jalur recovery/retrosesi (klaim memulihkan nilai lewat treaty outward)?
  Apa yang dilakukan `DLAFacintoTreaty_Act` — apa kepanjangan `DLA`?
- **Memblokir:** penetapan bounded context di D4 — domain klaim dan treaty tidak terpisah bersih.
- **Status:** terbuka
- **Catatan silang:** D1 §21.3 menemukan modul bernama `Treaty Contract Out` **tidak** merujuk satu
  pun objek treaty outward, sementara `Claim Non Prop` **merujuk** → memperkuat **OQ-022**.

---

### OQ-043 — Baris tabel keputusan tidak ikut terekspor (49 `DecisionTable`)

- **Pemilik:** pemilik export Pega + Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 4 — `NB FacIn/DecisionTable/IsUWAccepted.xml`
  (`ASM-FW-GISFW-WORK / ISUWACCEPTED / RULE-DECLARE-DECISIONTABLE`)
- **Temuan `[terverifikasi]`:** file `DecisionTable` di korpus memuat **nama properti masukan**
  (`<pyRuleName>ProposalAcceptStatus` di `pxRuleReferences`), **daftar hasil yang diizinkan**
  (`<pyTaskStatusXml>` → `confirm`, `reject`, `ask`, `banding`, `revise`, `decline`), **lebar kolom**
  (`<pyColumnPreferences>`, 5 kolom), dan metadata versi — **tetapi tidak memuat baris keputusannya**
  (kondisi → hasil). Berlaku untuk **seluruh 49 file** `DecisionTable` korpus.
- **Pertanyaan:** Apakah baris tabel keputusan memang tidak ikut dalam format ekspor ini, atau
  hilang pada proses ekspor? Di mana isi tabel yang berlaku dapat diperoleh?
- **Memblokir:** `IsUWAccepted` menggerbangi **22 shape Decision** di `InputInwardFacultativeOffer`
  — flow terbesar korpus. Tanpa isi tabelnya, **pemetaan nilai `ProposalAcceptStatus`
  (1/2/3/4/7/9) ke hasil (confirm/reject/ask/banding/revise/decline) tidak dapat dinyatakan**.
  Sama berlaku untuk `isApproved` (OQ-026), `IsLifeAccepted`, `MappingCoverage`, `BusinessType_DeT`,
  `LicensePlatRegion_DeT`, `MappingOutgoIndex`, `MappingAdditionalCoverageIndex`,
  `SetUploadHubAW1`–`3`, dan 40 tabel lain. Memblokir spesifikasi aturan persetujuan di FASE B.
- **Diperberat oleh OQ-011:** `DecisionTable/IsUWAccepted.xml` **berbeda isinya** di
  `Endorsment Fac In` (hash `9226a7db77`) dibanding `NB FacIn`/`RNW Fac In` (`f99bc43c45`).
  Aturan persetujuan **memang bercabang antar siklus**, tetapi apa yang bercabang tidak terbaca.
- **Status:** terbuka — **TERTUTUP untuk PremiumList Life + Endorsement Life** (2026-09-15,
  sumber: **keputusan work owner**).

  `[keputusan work owner]` **Keputusan accept/reject/decline dibuat MANUAL oleh inputor/admin** —
  bukan formula otomatis. Karena itu **kriteria masukan DecisionTable bukan aturan yang perlu
  direplikasi**; yang perlu direplikasi adalah **akibat** tiap keluaran.

  | DecisionTable | Keluaran | Akibat |
  | --- | --- | --- |
  | `ASM-FW-GISFW-WORK-LIFE` / `ISLIFEACCEPTED` | **Confirm** | naik ke tahap berikutnya |
  | | **Reject** | **selalu balik ke input** — inputor/admin isi ulang |
  | | **Decline** | case **ditutup/dibuang**, berhenti |
  | `ASM-FW-GISFW-WORK-LIFE` / `ISFLAGONGOINGPOLICY` | **1 = offer** | berhenti di tahap penawaran |
  | | **2 = premium** | lanjut ke Input Premium List Detail |

  `[terverifikasi]` Baris keputusan **tidak terekspor** (nol `pyRowData`), tetapi **nilai
  keluarannya terbaca** di kedua berkas — `Confirm`/`Decline`/`Reject` dan `Decline`/`Offer`/`Premium`.
  Konsisten dengan graf: `Decision1` dan `Decision3` memakai `IsLifeAccepted`.

  **Tetap terbuka untuk modul lain** yang memakai DecisionTable berbeda.
- **Perintah audit:**
  ```
  find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -path "*/DecisionTable/*.xml" -print | wc -l
  # -> 49
  grep -o "<pyTaskStatusXml>[^<]*" "NB FacIn/DecisionTable/IsUWAccepted.xml"
  ```

### OQ-044 — `ProposalAcceptStatus = 4`: ditulis rule "banding", diuji rule "facout"

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 4 — `NB FacIn/DataTransform/SetBandingProposal_DT.xml`
  dan `NB FacIn/When/IsFacout.xml`
- **Temuan `[terverifikasi]`:**
  - `SetBandingProposal_DT` (`RULE-OBJ-MODEL`) men-`SET` `.ProposalAcceptStatus = 4`, bersama
    `.OfferFacIn.IsBanding = "true"`, `.LetterNo = .BandingTo`, `FlagBanding.CARI11 = .BandingTo`.
  - `When/IsFacout.xml` (`RULE-OBJ-WHEN`) menguji `pyWorkPage.ProposalAcceptStatus = 4`.
- **Pertanyaan:** Apakah nilai `4` berarti "banding", "fac out", atau keduanya sekaligus? Bila
  keduanya, apa yang membedakan keduanya saat runtime?
- **Memblokir:** setiap pemetaan status penawaran ke keadaan bisnis. Nilai ini juga salah satu dari
  enam hasil `IsUWAccepted` yang pemetaannya tidak terbaca (OQ-043). Arti tetap **belum
  terverifikasi** (OQ-020).
- **Status:** terbuka
- **Perintah audit:**
  ```
  awk '/<pyPropertiesName>/{n=$0;gsub(/<[^>]*>/,"",n)} /<pyPropertiesValue>/{v=$0;gsub(/<[^>]*>/,"",v);
       if(n!=""){print n" := "v;n=""}}' "NB FacIn/DataTransform/SetBandingProposal_DT.xml"
  grep -oE "<pyLabel>[^<]{1,200}" "NB FacIn/When/IsFacout.xml" | grep -v "Available)"
  ```

### OQ-045 — `pyWorkPage.LetterNo` dipakai sebagai token routing persetujuan

- **Pemilik:** IAM + Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 4 — `NB FacIn/When/ToKadivTeknik.xml` dan 8 rule `To*` sejenis;
  penulisnya `NB FacIn/DataTransform/SetBandingProposal_DT.xml`
- **Temuan `[terverifikasi]`:** rantai lengkap terbaca — `SetBandingProposal_DT` men-set
  `.LetterNo = .BandingTo`, lalu shape Decision "Limit Akseptasi" memilih Assignment berdasarkan
  nilai literal `LetterNo`. Sepuluh token: `KADIVTEKNIK` (6×), `KADIVFACULTATIVE` (6×),
  `DIREKTURTEKNIK` (6×), `DIREKTURMARKETING` (6×), `DEPHEADUNDERWRITER` (6×), `SENIORUW` (3×),
  `MANAGERTEKNIK` (3×), `KADIVFINANCIAL` (3×), `DEPTHEADUWLIFE` (2×), `TREATYINDEPTHEAD` (2×).
- **Pertanyaan:** Apakah "nomor surat" memang dirancang sebagai penyimpan kode peran, atau ini
  workaround? Apakah nilai-nilai itu daftar tertutup? Siapa yang menetapkan `.BandingTo`?
- **Memblokir:** desain RBAC di sistem baru; pemetaan peran ke workbasket (bersama OQ-024 dan
  OQ-007). Pola sejajar dengan OQ-027 (`OperatorID.pyTelephone` menyimpan `TREATY1`/`SPVTREATY1`) —
  keduanya field bisnis yang dipakai menyimpan kode peran.
- **Status:** terbuka
- **Perintah audit:**
  ```
  grep -rhoE "LetterNo\]\[&amp;#61;\]\[&amp;quot;[A-Z0-9]+" "NB FacIn" "RNW Fac In" "Endorsment Fac In" \
    --include="*.xml" | sed 's/.*quot;//' | sort | uniq -c | sort -rn
  ```

### OQ-046 — Pangsa retro 0,45 / 0,05 dan ambang 3 miliar ter-hardcode, tanpa mata uang

- **Pemilik:** Product+Underwriting + Actuarial
- **Ditemukan di:** STEP D2 Tahap 4 — `NB FacIn/Activity/CountPremiNusareRetro_Act.xml`
  (`ASM-FW-GISFW-WORK / COUNTPREMINUSARERETRO_ACT`, 64.707 byte, 5 langkah `Property-Set`).
  Hash ternormalisasi `8f1268a660` — **identik di NB FacIn, RNW Fac In, dan Endorsment Fac In**.
- **Temuan `[terverifikasi]`** (rumus terbaca penuh, langka di korpus ini):
  ```
  local.sharernm1 := 0.45     local.sharernm2 := 0.05     local.tampungpreminet1 := 0
  bila .TSILiability <= "3000000000":
      .PremiLifeNusantaraRe := .GrossPremiumRetro * 0.45
  bila .TSILiability >  "3000000000":
      local.tampungpreminet1 := .GrossPremiumRetro * 0.45
      .PremiLifeNusantaraRe  := (.GrossPremiumRetro * 0.05) + local.tampungpreminet1
  .PremiumRetro := @toDecimal(.GrossPremiumRetro) - @toDecimal(.PremiLifeNusantaraRe)
  bila pyWorkPage.OfferFacIn.ProRateType != 3:
      .PremiLifeNusantaraRe := .Premium - .PremiumRetro
  ```
- **Pertanyaan:** Apa arti pangsa 0,45 dan 0,05 — dan terhadap siapa? Dalam mata uang apa ambang
  3.000.000.000 berlaku? Apakah nilai-nilai ini pernah berubah, dan di mana seharusnya dikelola?
  Apa arti `ProRateType = 3`?
- **Memblokir:** perhitungan premi retro di sistem baru; keputusan apakah parameter ini menjadi
  konfigurasi atau data referensi. **Ambang dibandingkan sebagai string berkutip**
  (`<="3000000000"`), bukan sebagai angka — perilakunya pada nilai di luar rentang perlu
  dikonfirmasi.
- **Status:** terbuka
- **Perintah audit:**
  ```
  awk '/<PropertiesName>/{n=$0;gsub(/<[^>]*>/,"",n)} /<PropertiesValue>/{v=$0;gsub(/<[^>]*>/,"",v);
       if(n!=""){print n" := "v;n=""}}' "NB FacIn/Activity/CountPremiNusareRetro_Act.xml"
  grep -oE "<pyStepsPreCondParamsWhen>[^<]*" "NB FacIn/Activity/CountPremiNusareRetro_Act.xml"
  ```

### OQ-047 — Daftar endpoint berada di tabel Oracle `M_LINK_SERVICE`, bukan hanya SystemSettings

- **Pemilik:** DBA + Platform
- **Ditemukan di:** STEP D2 Tahap 4 — `RNW Fac In/Activity/GetLinkService.xml`
  (`ASM-FW-GISFW-INT-M_LINK_SERVICE / GETLINKSERVICE / RULE-OBJ-ACTIVITY`)
- **Temuan `[terverifikasi]`:** activity ini melakukan `Obj-Browse` atas class
  `ASM-FW-GISFW-Int-M_LINK_SERVICE` dengan kunci **`KATEGORI_1`** dan **`KATEGORI_2`**, lalu
  hasilnya dipakai langkah `Connect-REST` berikutnya. Ia dipanggil dari activity Arasapas di ketiga
  siklus facultative dan dari seluruh activity Google Storage (9 rule di RNW Fac In saja).
- **Pertanyaan:** Apa isi tabel `M_LINK_SERVICE`, berapa baris, dan endpoint apa saja yang terdaftar?
  Apa arti `KATEGORI_1` / `KATEGORI_2`? Adakah endpoint production di dalamnya?
- **Memblokir:** inventaris integrasi keluar tidak lengkap tanpa isi tabel ini. Sapuan OQ-018 di D1
  hanya mencari URL literal **di dalam file rule**, sehingga sumber konfigurasi yang sebenarnya
  luput. Memblokir rencana konfigurasi/env var di sistem baru.
- **Catatan `[terverifikasi]`:** rule SystemSettings `LINKSERVICE!LINKSERVICE` terdaftar berkonflik
  (OQ-011 #324) tetapi **identik** di ketiga modul setelah normalisasi 21 tag — base URL **tidak**
  bercabang antar siklus.
- **Status:** **TERTUTUP** (2026-09-14, sumber: **korpus + data DBA**).

  `[terverifikasi]` Kontrak resolusi endpoint: `Claim Life/Activity/GetLinkService.xml`
  (`ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY`) melakukan
  `Obj-Browse` atas tabel `M_LINK_SERVICE` dengan
  `.KATEGORI_1 = Param.Kategori_1 AND .KATEGORI_2 = Param.Kategori_2`, mengambil kolom `.URL`, lalu
  memanggil `Connect-REST`.

  `[terverifikasi]` **Kunci kategori Claim — Life terbaca langsung di korpus:**
  `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` mengisi
  `Kategori_1 = "Klaim"` dan `Kategori_2 = "insertClaimLife"`.

  `[terverifikasi]` Lima berkas `Claim Life` menyentuh jalur ini: `GetLinkService.xml`,
  `InsertGoogleStorage_Act.xml`, `GetUrlGoogleStorage_Act.xml`, `DeleteGoogleStorage_Act.xml`,
  `serviceInsertArasapasClaimLife_act.xml`.

  `[data DBA]` Isi tabel **19 baris** diserahkan DBA. Endpoint Life:
  `(KATEGORI_1='Klaim', KATEGORI_2='insertClaimLife')` →
  `http://10.100.10.75:7315/Nusare-Integration-WS/resources1/restws/NusareClaim/insertClaimLife`.

  **Aturan mengikat sistem baru** — dicatat sebagai **ADR-0013**:
  resolusi URL keluar **harus runtime lookup** ke `M_LINK_SERVICE` lewat `(KATEGORI_1, KATEGORI_2)`
  **setiap kali**. **DILARANG** menanam URL sebagai literal, konstanta, **maupun env var**.
  Pemisahan dev–prod ditangani **isi tabel per-database**, bukan oleh kode. Yang boleh menjadi
  konstanta di kode hanyalah **kunci kategori**, bukan alamat.

  ⚠️ Aturan ini **meralat ADR-0004** (endpoint sebagai env var) untuk Claim — Life; lihat
  **ADR-0013**.
- **Perintah audit:**
  ```
  grep -oE "<pyStepsActivityName>[^<]*|<pyRuleName>[^<]*" "RNW Fac In/Activity/GetLinkService.xml" \
    | sed 's/<[^>]*>//' | sort -u
  grep -rl "M_LINK_SERVICE" "RNW Fac In" --include="*.xml" | wc -l
  ```

### OQ-048 — `GeminiAIGoogle_Act`: model AI pihak ketiga di dalam alur underwriting

- **Pemilik:** Product+Underwriting + Security
- **Ditemukan di:** STEP D2 Tahap 4 — `NB FacIn/Activity/GeminiAIGoogle_Act.xml` dan
  `RNW Fac In/Activity/GeminiAIGoogle_Act.xml`
- **Temuan `[terverifikasi]`:** activity bernama demikian ada di **NB FacIn** dan **RNW Fac In**,
  **tidak ada** di `Endorsment Fac In`. Ia termasuk 9 rule yang memakai `M_LINK_SERVICE` (OQ-047).
  **Isinya belum dibaca** — ini batas cakupan telusur, bukan batas korpus.
- **Pertanyaan:** Data apa yang dikirim ke layanan AI itu (apakah termasuk data tertanggung /
  data pribadi)? Apa yang dilakukan terhadap hasilnya — apakah memengaruhi keputusan underwriting?
  Apakah ada perjanjian pemrosesan data?
- **Memblokir:** penilaian kepatuhan dan privasi sebelum migrasi; keputusan apakah jalur ini ikut
  dimigrasi. Juga memblokir OQ-018/OQ-047 (endpoint mana yang dipakai).
- **Status:** terbuka
- **Perintah audit:**
  ```
  for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do ls "$m/Activity/" | grep -iE "google|gemini"; done
  ```

### OQ-049 — Empat class bersufiks `ENDORSEMENT`, satu di antaranya berprefix `ASM-SFAGIS-`

- **Pemilik:** Arsitektur Pega + pemilik export
- **Ditemukan di:** STEP D2 Tahap 4 — `Endorsment Fac In/Activity/SetValueToEDMWork.xml`
  (`ASM-SFAGIS-WORK-ENDORSEMENT / SETVALUETOEDMWORK`, 445.703 byte, 41 langkah)
- **Temuan `[terverifikasi]`:** empat class bersufiks `ENDORSEMENT` hidup berdampingan di korpus —
  `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` (38 rule), **`ASM-SFAGIS-WORK-ENDORSEMENT`** (31),
  `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` (10), `ASM-FW-GISFW-WORK-ENDORSEMENT` (5).
  Prefix `ASM-SFAGIS-` **tidak mengikuti pola `ASM-FW-GISFW-`** yang dipakai sisa korpus; dua class
  lain berprefix serupa juga ada: `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` (5),
  `ASM-FW-SFAGISFW-WORK-ACCOUNT` (5).
  Seluruh mesin salinan nilai lama endorsement facultative (`SetValueToEDMWork`,
  `SetErrorBatalEndorsement_Act`, `CountEndorsementData`, `SetEdmType`, 7× `SetOLDValueToEDMWork_*`)
  tinggal di class `ASM-SFAGIS-WORK-ENDORSEMENT`.
- **Pertanyaan:** Apakah `SFAGIS` aplikasi/ruleset yang berbeda dari `GISFW`? Kepanjangannya apa?
  Mengapa rule inti endorsement facultative tinggal di sana, bukan di `ASM-FW-GISFW-*`?
- **Memblokir:** penetapan batas aplikasi dan ruleset di D4; juga berkaitan dengan OQ-008
  (arti prefix `ASM`/`RNM`) dan OQ-009 (rule di class bawaan).
- **Status:** terbuka
- **Perintah audit:**
  ```
  awk -F'\t' 'toupper($3) ~ /SFAGIS/ {print $3}'      all-rules.tsv | sort | uniq -c | sort -rn
  awk -F'\t' 'toupper($3) ~ /ENDORSEMENT/ {print $3}' all-rules.tsv | sort | uniq -c | sort -rn
  ```


---

### OQ-050 — Sepuluh tombol ber-label `(dev)` di panel aksi, termasuk pemaksa status akhir

- **Pemilik:** Product+Underwriting + pemilik export Pega
- **Ditemukan di:** STEP D2 Tahap 5 — `Treaty In/Section/TreatyInActionButtons.xml`
  (`RULE-HTML-SECTION` / `DATA-PORTAL` / `TREATYINACTIONBUTTONS`, 369.789 byte)
- **Bukti `[terverifikasi]`:** dari `<pyLabel>`, **10 label memuat penanda `(dev)`** —
  `Force Edit (dev)`, `Force Resolve Complete(dev)`, `ReturnToInputor(dev)`,
  `RemoveLastComment (dev)`, `Populate treatyindetail (dev)`, `Save(dev)`, `Save EDM(dev)`,
  `Test Error (dev)`, `TestCopyDifference(dev)`, `fix Spl(dev)`. Jumlah yang sama di
  `Treaty In Adjustment` (section-nya identik). Tiga di antaranya terikat rule nyata:
  `TreatyInForceEdit`, `TreatyInForceResolveComplete`, `TreatyInReturntoInputor`,
  `TreatyInRemoveLastComment`.
- **Pertanyaan:** Apakah tombol-tombol ini aktif di production? Siapa yang dapat melihatnya?
  `Force Resolve Complete` melompati seluruh tangga persetujuan — apakah itu disengaja sebagai
  jalan darurat, atau sisa pengembangan yang tertinggal?
- **Memblokir:** pemahaman jalur status yang sebenarnya (tangga persetujuan dapat dilewati);
  desain RBAC dan audit trail di sistem baru. Berkaitan dengan OQ-018 (korpus belum tentu
  production).
- **Status:** terbuka
- **Perintah audit:**
  ```
  grep -oE "<pyLabel>[^<]{1,45}" "Treaty In/Section/TreatyInActionButtons.xml" \
    | sed 's/<pyLabel>//' | grep -i "(dev)" | sort -u
  ```

### OQ-051 — Otorisasi tombol memakai **indeks tetap** pada daftar workbasket operator

- **Pemilik:** IAM
- **Ditemukan di:** STEP D2 Tahap 5 — ekspresi visibilitas di
  `Treaty In/Section/TreatyInActionButtons.xml` dan rule terkait
- **Bukti `[terverifikasi]`:** hak akses diuji dengan
  `OperatorID.pyWorkBasketList(2).pyWorkBasketName = 'ReasTreatyInAdmin'` (36 kemunculan),
  `pyWorkBasketList(2).pyWorkBasketName = TreatyIn.Position` (8),
  `pyWorkBasketList(1).pyWorkBasketName = 'ReasTreatyInSecHead'` (4),
  `pyWorkBasketList(2).pyWorkBasketName = 'ReasTreatyInSecHead'` (1).
- **Pertanyaan:** Mengapa indeks tetap, bukan pencarian berdasarkan nama? Apa yang menjamin
  urutan workbasket dalam profil operator? Apa yang terjadi bila operator punya jumlah workbasket
  berbeda?
- **Memblokir:** desain RBAC di sistem baru; pemetaan peran (bersama OQ-007, OQ-024, OQ-027,
  OQ-045, OQ-053). Juga risiko migrasi: urutan daftar workbasket adalah **data**, bukan konfigurasi
  peran.
- **Status:** terbuka
- **Perintah audit:**
  ```
  grep -rhoE "pyWorkBasketList\([0-9]+\)\.pyWorkBasketName[^<]{0,34}" "Treaty In" "Treaty In Adjustment" \
    --include="*.xml" | sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g' | sort | uniq -c | sort -rn
  ```

### OQ-052 — `Reject` dan `Decline` hidup berdampingan; bedanya tidak dijelaskan korpus

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 5 — `Treaty In/DataTransform/Akseptasi_DT.xml`
  (`DATA-PORTAL!AKSEPTASI_DT`)
- **Bukti `[terverifikasi]`:** `TreatyIn.ChooseStatusAkseptasi` punya **3 nilai** yang dipilih
  pengguna: `"Accept"`, `"Reject"`, `"Decline"`. Akibatnya berbeda:
  - `Reject`  → `Position := "ReasTreatyInAdmin"`, `StatusAkseptasi := "Reject"`,
    `PositionUsername := TreatyIn.CommentList(1).OperatorName` (atau `(<LAST>)` di jalur revisi)
  - `Decline` → `Position := ""`, `StatusAkseptasi := "Decline"`, `PositionUsername := ""`
  Sementara `TreatyIn.StatusAkseptasi` punya **4 nilai** — yang keempat, `"Resolve Complete"`,
  tidak pernah dipilih pengguna, hanya dihasilkan mesin di ujung tangga.
- **Pertanyaan:** Apa beda bisnis antara "reject" dan "decline" pada akseptasi treaty inward?
  Apakah `Decline` final dan `Reject` dapat diperbaiki? Apa arti `Resolve Complete`?
- **Memblokir:** pemetaan status ke keadaan bisnis; desain state machine di sistem baru. Arti
  keempat nilai tetap **belum terverifikasi** (OQ-020).
- **Status:** terbuka
- **Perintah audit:** lihat awk berpenghitung kedalaman di `flows/_METHOD-noflow.md` §3.3.

### OQ-053 — Identitas orang **ditetapkan** sebagai pemilik tugas berikutnya (bukan sekadar diuji)

- **Pemilik:** IAM + Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 5 — `Treaty In/DataTransform/Akseptasi_DT.xml`
  (identik di `Treaty In Adjustment`, hash `58b8e650`)
- **Bukti `[terverifikasi]`:** setiap cabang `Accept` men-`SET` tripel
  `TreatyIn.Position` (nama workbasket), `TreatyIn.StatusAkseptasi`, dan
  **`TreatyIn.PositionUsername`** — yang diisi **identitas orang ter-hardcode**.
  **5 identitas berbeda** ditemukan. Cabang `Reject` mengisinya dari
  `TreatyIn.CommentList(1).OperatorName` (data), bukan konstanta.
  **Nilai nama orang tidak disalin** ke artefak D2 (`flows/_METHOD.md` §1.4).
- **Pertanyaan:** Apakah tangga persetujuan treaty inward memang menunjuk **orang tertentu**
  alih-alih peran? Apa yang dilakukan bila orang itu berganti jabatan, cuti, atau keluar —
  apakah rule harus diubah dan di-deploy ulang?
- **Memblokir:** desain RBAC di sistem baru; ini **pola paling berat** dari keluarga guard
  identitas. Berbeda dari OQ-021 (identitas dipakai sebagai *guard*, menguji siapa yang bekerja) —
  di sini identitas **ditetapkan sebagai tujuan** rute berikutnya.
- **Status:** terbuka
- **Perintah audit (menghitung tanpa menampilkan nilai):**
  ```
  awk '/<pyPropertiesName>/{n=$0;gsub(/<[^>]*>/,"",n)}
       /<pyPropertiesValue>/{v=$0;gsub(/<[^>]*>/,"",v);
          if(n=="TreatyIn.PositionUsername" && v!="" && v !~ /\./){gsub(/"/,"",v); print v}; n=""}' \
    "Treaty In/DataTransform/Akseptasi_DT.xml" | sort -u | wc -l
  # -> 5
  ```

### OQ-054 — Folder `Claude outputs` berisi berkas non-Pega di dalam korpus ekspor

- **Pemilik:** pemilik export Pega
- **Ditemukan di:** STEP D2 Tahap 5 — sapuan direktori korpus
- **Bukti `[terverifikasi]`:** empat modul memuat folder `Claude outputs` yang **bukan** folder
  tipe rule Pega:
  | Lokasi | Isi |
  | --- | --- |
  | `NB Treaty In/Claude outputs/` | `Struktur_MenuNBTreatyIn.xlsx` (275.860 B) |
  | `Treaty Contract Out/Claude outputs/` | `Struktur_InboxTreatyContract.xlsx` (306.979 B) |
  | `Master Contract Retro Life/Claude outputs/` | `Struktur_GridRetrocessionLife.xlsx` (11.944 B), **`perubahan_skill.diff`** (10.278 B) |
  D1 menginventarisasi hanya `*.xml`, sehingga berkas ini **tidak pernah terhitung**. Isinya
  **tidak dibaca** (bukan rule Pega; korpus READ-ONLY). Selain itu ditemukan berkas `.xlsx`
  sejenis di `RNW Fac In` (`Struktur_InputRenewalFacultativeIn.xlsx` dan berkas kunci sementara
  `~$…`).
- **Pertanyaan:** Siapa yang menaruh berkas ini, apakah isinya relevan sebagai dokumentasi, dan
  apakah ekspor korpus mencerminkan isi ruleset Pega yang sesungguhnya atau sudah tercampur
  artefak analisis?
- **Memblokir:** kepercayaan terhadap kelengkapan korpus (bersama OQ-003 `excludeXML` dan OQ-018).
  Tidak memblokir telusur rule.
- **Status:** terbuka
- **Perintah audit:**
  ```
  find . -path ./OUTPUT_HASIL_RNM -prune -o -type d -name "Claude outputs" -print
  find . -path ./OUTPUT_HASIL_RNM -prune -o -type f ! -name "*.xml" -print | head -20
  ```

### OQ-055 — Nomor revisi treaty disimpan di dalam string `ID`, posisi karakter ter-hardcode

- **Pemilik:** DBA + Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 5 — `Treaty In Adjustment/Activity/TreatyInRevisi_post.xml`
  (`DATA-PORTAL!TREATYINREVISI_POST`) dan `Treaty In Adjustment/RDBList/GetTreatyRevisionID.xml`
- **Bukti `[terverifikasi]`:**
  ```
  TreatyIn.ID    := TreatyIn.ID + "/R01"
  TreatyIn.OLDID := TreatyIn.ID
  TreatyIn.ID    := @If(@substring(TreatyIn.ID,10,12) < 10,
                        @substring(TreatyIn.ID,0,7) + "/R0" + (@toInt(@substring(TreatyIn.ID,10,12))+1),
                        @substring(TreatyIn.ID,0,7) + "/R"  + (@toInt(@substring(TreatyIn.ID,10,12))+1))
  ```
  dan di sisi database: `select TO_NUMBER(SUBSTR(ID,10,2))+1 as HASIL1`.
- **Pertanyaan:** Apa format lengkap `TreatyIn.ID` dan apakah posisi karakter 10–12 dijamin stabil?
  Apa yang terjadi pada revisi ke-100? Apakah `OLDID` dipakai sebagai tautan riwayat?
- **Memblokir:** desain kunci dan versi entitas treaty di sistem baru; migrasi data historis.
  Logika yang sama ter-hardcode **di dua tempat** (rule Pega dan SQL), sehingga perubahan format
  ID harus dilakukan serempak.
- **Status:** terbuka
- **Perintah audit:**
  ```
  awk '/<PropertiesName>/{n=$0;gsub(/<[^>]*>/,"",n)} /<PropertiesValue>/{v=$0;gsub(/<[^>]*>/,"",v);
       if(n!=""){print n" := "v;n=""}}' "Treaty In Adjustment/Activity/TreatyInRevisi_post.xml"
  grep -oE "<pyBrowseSQL>[^<]{0,300}" "Treaty In Adjustment/RDBList/GetTreatyRevisionID.xml"
  ```

### OQ-056 — Modul master produk life membawa jalur tulis ke master **treaty inward**

- **Pemilik:** Product+Underwriting + DBA
- **Ditemukan di:** STEP D2 Tahap 5 — `Master Product Name Life/RDBList/SaveTreatyIn.xml` dan
  `Activity/SetTreatyIn_Act.xml`
- **Bukti `[terverifikasi]`:**
  - `RDBList/SaveTreatyIn.xml` = `ASM-FW-GISFW-INT-TREATY_IN!ASM!SAVETREATYIN`, memanggil
    `POOLDATA.PEGA_TREATY_IN` — procedure yang sama yang dipakai modul `Treaty In`.
  - `RDBList/BrowseTreatyIn.xml` membaca `POOLDATA.M_TREATY_IN` dan `POOLDATA.M_TREATY_IN_EDM`.
  - Pemanggilnya di modul ini: `Activity/SetTreatyIn_Act.xml` (`DATA-PORTAL!SETTREATYIN_ACT`,
    159.737 B, 15 langkah — **urutan langkah identik** dengan versi `Treaty In`, isi berbeda
    47 byte).
  - **Nuansa penting:** `ASM!SAVETREATYIN` (hash `e371c194`) dan `ASM!BROWSETREATYIN`
    (`196d49b5`) **identik di 6 modul** (Treaty In, Treaty In Adjustment, NB Treaty In,
    EDM Treaty In, NB FacIn, Master Product Name Life) dan **tidak terdaftar di register OQ-011** —
    jadi ini rule tulis **bersama**, bukan salinan yang menyimpang.
  - **`Akseptasi_DT`, `TreatyInSubmit`, `TreatyInActionButtons` tidak ada** di modul ini — tangga
    persetujuan treaty **tidak** ikut tersalin.
- **Pertanyaan:** Apakah jalur tulis `M_TREATY_IN` benar-benar dieksekusi dari modul master produk
  life, dan dalam keadaan apa? Bila ya, apakah ia melewati persetujuan? Bila tidak, mengapa rule
  itu ikut diekspor ke sini?
- **Memblokir:** penetapan kepemilikan data master treaty inward di D4; integritas data (siapa
  boleh menulis `M_TREATY_IN`).
- **Status:** terbuka
- **Perintah audit:**
  ```
  awk -F'\t' '$1 ~ "^Master Product Name Life/" && $8 ~ /M_TREATY_IN|PEGA_TREATY_IN/{print $5"  "$1"  "$8}' all-rules.tsv
  grep -rl "SaveTreatyIn" "Master Product Name Life" --include="*.xml"
  ```

### OQ-057 — Kode `OR` dan `REINSTYPEID`: tidak pernah muncul sebagai nilai literal

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 5 — `Master Product Name Life/Activity/GetReinsTypeOR_Life.xml`,
  `RDBList/BrowseReinstypeOR_SQL.xml`; `Master Contract Retro Life` (12 Connect-SQL)
- **Bukti `[terverifikasi]`:**
  - Kode **`OR`** hanya ada **sebagai bagian nama rule**.
    `grep -rhoE "\"OR[A-Z0-9]*\"|'OR[A-Z0-9]*'" "Master Product Name Life" --include="*.xml"`
    → **kosong**. `flows/_METHOD.md` §2.2 melarang memakai nama rule sebagai bukti perilaku,
    sehingga arti `OR` **tidak dapat dinyatakan**.
  - **`REINSTYPEID`** muncul **129×** di `Master Contract Retro Life`, tetapi selalu diisi dari
    slot generik `TempInputData.CARI3/4/5` lalu diteruskan sebagai `p_REINSTYPEID` ke tiga
    procedure berbeda — **tidak ada satu pun nilai literal untuk dikutip**.
- **Pertanyaan:** Apa kepanjangan dan arti kode `OR` pada konteks produk life? Apa daftar nilai
  `REINSTYPEID` yang berlaku, dan di tabel mana ia didefinisikan?
- **Memblokir:** model data produk & kontrak retro life; enumerasi jenis reasuransi. Tidak dapat
  dijawab dari korpus sama sekali — perlu jawaban manusia atau isi tabel database.
- **Status:** terbuka

### OQ-058 — Tabel internal Pega diakses langsung lewat SQL

- **Pemilik:** DBA + Arsitektur Pega
- **Ditemukan di:** STEP D2 Tahap 5 — `Master Product Name Life`
- **Bukti `[terverifikasi]`:** objek **`DATAPEGA.PC_ASM_FW_GCNMFW_WORK`** dirujuk oleh rule
  `RULE-CONNECT-SQL` di modul ini. `DATAPEGA` adalah skema penyimpanan Pega, dan
  `PC_ASM_FW_GCNMFW_WORK` adalah tabel work object untuk class `ASM-FW-GCNMFW-WORK` —
  diakses **langsung lewat SQL**, bukan lewat API/rule Pega.
- **Pertanyaan:** Mengapa tabel internal Pega dibaca langsung? Apakah ada tempat lain yang
  melakukan hal serupa? Apakah ini menciptakan ketergantungan pada skema internal Pega yang akan
  hilang setelah migrasi?
- **Memblokir:** rencana migrasi — data yang dibaca dari tabel internal Pega tidak punya padanan
  otomatis di sistem baru. Berkaitan dengan OQ-016 (peran skema).
- **Status:** terbuka
- **Perintah audit:** `awk -F'\t' '$8 ~ /DATAPEGA/{print $1"  "$5"  "$8}' all-rules.tsv`

### OQ-059 — Slot parameter generik `CARI1`…`CARI30` menuju SQL

- **Pemilik:** DBA + Product+Underwriting
- **Ditemukan di:** STEP D2 Tahap 5 — `Master Contract Retro Life` (`TempInputData.CARI3/4/5`),
  `Treaty In Adjustment` (`InputData.CARI1`), `Treaty In` (`StatusDoc.CARI30`);
  sebelumnya juga di `NB Treaty In` (`InputData.CARI20`, `CARI21`, `CARI3` — `flows/NB Treaty In.md`
  §3.1 & §5.1) dan `NB FacIn` (`FlagBanding.CARI11` — `flows/NB FacIn.md` §2.1)
- **Bukti `[terverifikasi]`:** properti bernomor `CARI<n>` dipakai sebagai **slot parameter
  generik** yang diisi tepat sebelum pemanggilan SQL, lalu dibaca di dalam `<pyBrowseSQL>` sebagai
  `{InputData.CARI<n>}` / `{TempInputData.CARI<n>}`. Contoh:
  `REINSTYPEID := {TempInputData.CARI3}`, `InputData.CARI1 := TreatyIn.ID + "%"`,
  `.StatusDoc.CARI30 != 1`.
- **Pertanyaan:** Apa arti tiap slot `CARI<n>` per pemanggilan? Apakah ada konvensi, atau nomornya
  dipakai ulang untuk maksud berbeda di tempat berbeda?
- **Memblokir:** pembacaan kontrak antarmuka rule → SQL; setiap pemetaan parameter ke kolom akan
  menjadi tebakan selama ini terbuka. Berkaitan erat dengan OQ-001 (tidak ada DDL) dan OQ-002
  (isi procedure tidak diketahui).
- **Status:** terbuka
- **Perintah audit:**
  ```
  grep -rhoE "(Input|TempInput|StatusDoc|FlagBanding)[A-Za-z]*\.CARI[0-9]+" . --include="*.xml" \
    | sort | uniq -c | sort -rn | head -20
  ```

### OQ-060 — Cakupan kolom `CURRENCY` pada rekam akseptasi klaim Life

- **Pemilik:** Product+Underwriting (+ DBA untuk konfirmasi kolom)
- **Ditemukan di:** grilling Ronde 2 Q10 (`.scratch/claim-life/`), 2026-09-14 — work owner
  mengklasifikasikan 8 kolom sebagai **uang**, `EM_PERCENT` sebagai **persen**, dan menyebut
  `CURRENCY`; cakupannya tidak dinyatakan
- **Pertanyaan:** Apakah satu nilai `CURRENCY` berlaku untuk **seluruh** rekam akseptasi, atau
  sebuah rekam dapat memuat kolom uang dengan mata uang berbeda-beda?
- **Memblokir:** bentuk tipe uang di sistem baru — apakah cukup satu pasang `(amount, currency)`
  per rekam, atau tiap kolom uang membawa mata uangnya sendiri. **Tidak memblokir** keputusan
  non-float itu sendiri (**ADR-0003** sudah `accepted`); memblokir bentuknya.
- **Terkait:** **ADR-0003**, **ADR-0001** (nilai uang menyeberang ke Komite), OQ-001 (tidak ada DDL),
  OQ-046 (nilai ter-hardcode tanpa mata uang di konteks lain)
- **Bukti tambahan `[terverifikasi]` (2026-09-14, Ronde 3):** `CURRENCY` **ada di tingkat baris**,
  bukan hanya di header — muncul sebagai `…PremiumListDetail(idx).CURRENCY` dan
  `.AdjustmentList(<LAST>).CURRENCY`. `Claim Life/Activity/SetIndexAdjustmentList.xml`
  (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SETINDEXADJUSTMENTLIST`) menyalin `CURRENCY` **dan**
  `CURRENCYID` dari `AdjustmentList(1)` ke `AdjustmentList(<LAST>)`. Jadi dalam praktiknya seluruh
  baris satu klaim bermata uang **seragam**, meskipun strukturnya membolehkan campur.
  Pertanyaannya menyempit menjadi: apakah keseragaman itu **aturan** atau kebetulan implementasi.
- **Status:** terbuka — **TERTUTUP untuk Claim — Life** (2026-09-14, sumber: **korpus + keputusan work owner**).

  `[terverifikasi]` `Claim Life/Activity/SetIndexAdjustmentList.xml`
  (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY`) menyalin
  `CURRENCY` **dan** `CURRENCYID` dari `AdjustmentList(1)` ke `AdjustmentList(<LAST>)`. Dalam praktik
  seluruh baris satu klaim karena itu bermata uang **seragam**, meskipun struktur membolehkan
  per-baris.

  `[keputusan work owner]` **Bentuk tipe uang sistem baru: `(amount, currency)` per baris, dengan
  invariant — semua baris satu klaim WAJIB bermata uang sama.** Konsisten dengan **ADR-0003**.

  `[data DBA]` Kolom uang di Oracle bertipe **`NUMBER` tanpa presisi**, dan `CURRENCY VARCHAR2(100)`
  pada `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`. Konsekuensi mengikat: Go **wajib** memakai desimal
  presisi arbitrer — `float64` dilarang. Diperkuat di **ADR-0003**.

  Cakupan modul lain tetap terbuka.
- **Audit:**
  ```
  grep -rn "CURRENCY" "Claim Life" --include="*.xml" | head
  awk '/<pyBrowseSQL>/,/<\/pyBrowseSQL>/' "Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml"
  ```

### OQ-061 — `STS_REJECT` hidup di tiga tingkat; "satu status klaim" belum menyatakan tingkat baris

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** grilling Ronde 2 Q8/Q9 (`.scratch/claim-life/`), 2026-09-14 — work owner
  memutuskan sistem baru memakai **satu status klaim**; korpus memperlihatkan tiga tingkat
- **Bukti `[terverifikasi]`:** di `Claim Life`, `STS_REJECT` muncul sebagai
  `pyWorkPage.ClaimData.STS_REJECT` (tingkat klaim),
  `pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail(local.IndexPremium).STS_REJECT`
  (tingkat baris premi), dan `.AdjustmentList(<LAST>).STS_REJECT` (tingkat baris penyesuaian,
  `Claim Life/Activity/SaveAdjustment_Act.xml`). Activity
  `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SETSTS_REJECT` (`Claim Life/Activity/SetSTS_Reject.xml`)
  berisi **satu** `Property-Set`: `.STS_REJECT = Primary.STS_REJECT` — **menurunkan** status klaim
  ke tingkat baris. Di sisi Komite, `KomitePostAdjustment.xml` menulis **kedua** tingkat baris.
- **Pertanyaan:** Apakah baris-baris dalam satu klaim dapat berakhir dengan status **berbeda**
  (sebagian aksep, sebagian reject)? Bila ya, apa status klaimnya? Bila tidak, mengapa Pega
  menyimpannya per baris?
- **Memblokir:** model data dan mesin status di sistem baru; kontrak jalur balik Komite
  (**ADR-0001** Kontrak 2) bekerja pada tingkat baris, bukan tingkat klaim.
- **Terkait:** **ADR-0001**, **ADR-0007** (jejak audit per transisi — per klaim atau per baris?),
  OQ-039
- **Status:** terbuka — **TERJAWAB untuk Claim — Life** (2026-09-14, sumber: work owner).
  `[keputusan work owner]` **Unit keputusan = baris `AdjustmentList`.** `STS_REJECT` pada
  `PremiumListDetail` adalah **cerminan** hasil baris `AdjustmentList` **terakhir**, bukan unit
  keputusan tersendiri; header klaim mengikuti adjustment terakhir. Dicatat di **ADR-0011**.
  `[terverifikasi]` Didukung gerbang wewenang yang menguji `.STS_REJECT` **tingkat baris** —
  `Claim Life/Section/AdjustmentDetail_Section.xml`
  (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION`), tombol "Reject Outstanding":
  `pyWorkPage.pyPosition =='ReasLifeAdmin' && …CLAIM_NO !='' && .STS_REJECT == 0` — dan oleh
  penulisan berpasangan dua tingkat pada `RejectOSClaimLife_Act` dan `KomitePostAdjustment`.
  **Tetap terbuka untuk 3 modul Claim lain** — lihat catatan di OQ-062.
- **Sensus penulis `[terverifikasi]` (2026-09-14)** — setiap `Property-Set` yang menulis
  `STS_REJECT` di kedua modul, beserta sasaran dan nilainya:

  | Nilai | Rule penulis | Sasaran |
  | ---: | --- | --- |
  | `0` | `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `.STS_REJECT` |
  | `1` | `Claim Life/Activity/SaveAdjustment_Act.xml` | `.AdjustmentList(<LAST>).STS_REJECT` |
  | `1` (6×) | `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `…PremiumListDetail(idx).STS_REJECT` dan `…AdjustmentList(idx).STS_REJECT` |
  | `2` (2×) | `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `.STS_REJECT`, `…PremiumListDetail(idx).STS_REJECT` |
  | `2` (2×) | `Komite Claim Life/Activity/KomitePostAdjustment.xml` | dua tingkat baris yang sama |
  | salin **turun** | `Claim Life/Activity/SetSTS_Reject.xml` | `.STS_REJECT = Primary.STS_REJECT` |
  | salin **naik** | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` | `pyWorkPage.ClaimData.STS_REJECT = .STS_REJECT` |

  Dua baris terakhir penting: penyalinan terjadi **dua arah** — klaim→baris dan baris→klaim.
  Jadi tingkat mana yang otoritatif **tidak terbaca dari korpus**.
- **Audit:**
  ```
  grep -roh "[A-Za-z.()A-Za-z_]*STS_REJECT" "Claim Life" --include="*.xml" \
    | sed 's/(Local\.[A-Za-z]*)/(IDX)/g' | sort | uniq -c | sort -rn
  ```

### OQ-062 — Apakah `STS_REJECT` `1`/`2` final, dan apakah klaim dapat dibuka ulang?

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** grilling Ronde 2 Q9 (`.scratch/claim-life/grilling-ronde-2.md`), 2026-09-14 —
  work owner menjelaskan **bagaimana** `0` menjadi `1`/`2`, tetapi **tidak** menyatakan apakah
  keduanya final
- **Bukti `[terverifikasi]`:** sensus penulis lengkap (lihat **OQ-061**) memperlihatkan
  **tidak ada satu pun rule** di `Claim Life` maupun `Komite Claim Life` yang menulis `0` setelah
  `1` atau `2`. Satu-satunya penulis nilai `0` adalah
  `Claim Life/Activity/SaveOutStandingLife_Act.xml`. Jadi di sistem lama **tidak ada jalur buka
  ulang yang terbaca**.
- **Pertanyaan:** Apakah ketiadaan jalur buka-ulang itu **disengaja** (klaim yang sudah aksep/reject
  memang final) atau **kekurangan** yang selama ini ditangani di luar sistem? Bila klaim perlu
  dibuka ulang di sistem baru, itu perilaku **baru**, bukan paritas.
- **Memblokir:** keadaan akhir mesin status; apakah transisi `1`/`2` → `0` perlu ada, dan siapa yang
  boleh melakukannya (**ADR-0002** hanya menetapkan tiga peran, tanpa wewenang buka ulang).
  Ikut memengaruhi **ADR-0007** — sebuah buka-ulang wajib terekam siapa + kapan.
- **Terkait:** **OQ-061** (tingkat mana yang dibuka ulang), **OQ-039** (dua jalur akseptasi/penolakan
  langsung di Claim Life)
- **Status:** terbuka — **TERJAWAB untuk Claim — Life** (2026-09-14, sumber: work owner).
  `[keputusan work owner]` **Terminal per baris `AdjustmentList`**: sekali sebuah baris bernilai
  `1` atau `2`, nilainya tidak berubah lagi. **Klaim tidak terminal** — setelah penolakan, SPV
  membuat baris `AdjustmentList` **baru** dan siklus berulang. **Revisi = baris baru.**
  Dicatat di **ADR-0011**.
  `[terverifikasi]` Didukung sensus penulis: tidak ada rule yang menulis `0` setelah `1`/`2`, dan
  `Claim Life/Activity/SetIndexAdjustmentList.xml` menyalin 8 kolom dari `AdjustmentList(1)` ke
  `AdjustmentList(<LAST>)` **tanpa** `STS_REJECT` — baris baru memang mulai dari kosong.
  **Tetap terbuka untuk 3 modul Claim lain** (`Claim Prop` 15 kemunculan `STS_REJECT` /
  126 `.AdjustmentList`; `Claim Non Prop` 17 / 66; `Claim Fac In` 41 / 12) — pola yang sama ada di
  sana dan **belum** dinyatakan work owner.
- **Audit:**
  ```
  for m in "Claim Life" "Komite Claim Life"; do
    grep -rn "<PropertiesName>[^<]*STS_REJECT</PropertiesName>" "$m" --include="*.xml"
  done   # nilai ada di baris berikutnya
  ```

### OQ-063 — Gerbang "Send ke Komite" dapat dilewati oleh `Type = 'TP'` / `'TR'`, tanpa cek peran

- **Pemilik:** Product+Underwriting (aturan bisnis) + IAM (wewenang)
- **Ditemukan di:** verifikasi jawaban Ronde 3 (`.scratch/claim-life/grilling-ronde-3.md`),
  2026-09-14
- **Bukti `[terverifikasi]`:** `Claim Life/Section/AdjustmentDetail_Section.xml`
  (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION`,
  557.563 byte) menggerbangi kontrol jalur Komite (`GetListKomiteLife`) dengan:

  ```
  pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'
  ```

  Karena **disjungsi**, kedua nilai `Type` itu membuat kontrol tampil **tanpa memeriksa peran sama
  sekali**. Bandingkan dua gerbang tetangganya di berkas yang sama, yang murni peran:
  `pyWorkPage.pyPosition =='ReasLifeSPV'` (tombol "Save to Outstanding") dan
  `pyWorkPage.pyPosition =='ReasLifeAdmin' && … && .STS_REJECT == 0` (tombol "Reject Outstanding").
- **Pertanyaan asli:** Apakah ini **disengaja** — mis. `TP`/`TR` adalah jenis polis yang boleh
  dikirim ke Komite oleh peran mana pun — atau **cacat** pada ekspresi (kurang tanda kurung,
  sehingga maksud sebenarnya `SPV && (TP || TR)`)?

**DIJAWAB SEBAGIAN (2026-09-14, sumber: work owner).**

`[keputusan work owner]` **Bacaan gerbang ditetapkan:** untuk klaim ber-`Type` **`TP`** atau
**`TR`**, pengiriman ke Komite **tidak dibatasi peran SPV** — **`ReasLifeAdmin` pun dapat mengirim
langsung**. Untuk tipe lain, **hanya SPV**. Ini **meralat** aturan "send ke Komite hanya SPV" di
**ADR-0011**, yang kini berlaku untuk non-`TP`/`TR` saja.

`[terverifikasi]` **Bagian "cacat ekspresi" dari pertanyaan ini TERTUTUP.** Dua alasan, keduanya
terbaca dari korpus:

1. **Tidak ada ambiguitas presedensi.** Ekspresinya memakai `||` saja, tanpa `&&`:
   `pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'`.
   Pembacaan `SPV && (TP || TR)` menuntut sebuah `&&` yang **tidak ada** di ekspresi itu.

2. **`=` di korpus ini adalah pembanding, bukan assignment.** Sensus korpus-wide:

   | Tag | memakai `==` | memakai `=` tanpa `==` |
   | --- | ---: | ---: |
   | `<pyCondition>` | 954 | **2.282** |
   | `<pyStepsPreCondParamsWhen>` | 16.390 | 273 |

   Pada `<pyCondition>` — jenis tag yang dipakai gerbang ini — bentuk `=` justru **mayoritas**.
   Lebih tegas lagi, banyak ekspresi **mencampur keduanya dalam satu baris**, pada tempat yang akan
   rusak kalau `=` berarti assignment:

   - `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"` — bila `=` assignment, **setiap**
     baris akan terpilih (`Komite Claim Life/Activity/KomitePostAdjustment.xml`)
   - `(.Type== 1 && …) || ((.Type = 2 ||.Type = 4) && …) || (.Type = 3 && …)` — `==` dan `=` pada
     **properti yang sama** di ekspresi yang sama
   - `pyWorkPage.AcceptStatus = 1 && pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop` — bila
     assignment, cabang aksep **dan** cabang tolak akan sama-sama jalan
   - `pyWorkPage.pxFlow(InputOfferFacultativeIn).pxRouteTo="ReasFacInAdmin"` — assignment ke
     properti sistem yang bersifat baca

   Kesimpulan: `=` dan `==` dipakai **bergantian sebagai pembanding**. Gerbang itu memang
   bebas-peran untuk `TP`/`TR` — **bukan** karena bug operator.

**SISA BUTIR — SELURUHNYA TERJAWAB (2026-09-14, Ronde 4, sumber: work owner):**

1. **Arti `TP` dan `TR`** → `[keputusan work owner]` **`TP` = Payable, `TR` = Receivable.**
   ⚠️ Definisi ini **tidak ada di korpus**. `[terverifikasi]` Pencarian sudah tuntas dan nihil:
   tidak ada satu pun label, caption, atau opsi yang memasangkan kode itu dengan teks; satu-satunya
   tempat nilainya **ditetapkan** adalah `PremiumList Life/Activity/SubmitPremiumList_Act.xml`
   (`ASM-FW-GISFW-WORK-LIFE!SUBMITPREMIUMLIST_ACT`, 214.154 byte), dan di sana keempat kode
   (`QP`/`QR`/`TP`/`TR`) hanya **diteruskan** ke slot parameter generik `InputData.CARI20` dengan
   precondition `.Type=="<kode itu sendiri>"` — melingkar. Definisinya berada di luar korpus.
   **`QP` dan `QR` belum dijawab** dan tetap terbuka di **OQ-020**.
   ⚠️ **KOREKSI 2026-09-16:** sempat ditandai "tertutup dari dropdown EDM Life (QR=Receivable/
   QP=Payable)" — itu **KELIRU**. Sensus korpus ulang: `grep Receivable|Payable|Receiveable` di
   `Endorsement Life/` = **NOL berkas**. Label itu hanya muncul di `pyLocalizedValue` payload
   `DATA_JSON` runtime satu polis (instance data), **BUKAN rule korpus**. Arti QP/QR **tetap belum
   terverifikasi**. OQ-020 **TETAP TERBUKA**. Pelajaran: instance JSON ≠ bukti korpus.
2. **Apakah perilaku "TP/TR bebas peran" diinginkan** → **PARITAS.** `[keputusan work owner]`
   Dibawa apa adanya, **tidak diperketat** pada migrasi ini; risikonya dicatat sebagai **risiko
   RBAC** di **ADR-0012**, untuk ditinjau saat konteks **Komite Life** / **IAM** digarap.

`[terverifikasi]` `Type` menggerbangi **dua** hal, bukan satu: (a) wewenang send-ke-Komite di
`Claim Life/Section/AdjustmentDetail_Section.xml`, dan (b) jendela validasi Date of Loss di
`Claim Life/Activity/ValidasiDOL_Act.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!VALIDASIDOL_ACT`,
59.747 byte) — `QP`/`QR` → `GROSS_VALUATION_*` dengan pergeseran **nol**; `TP`/`TR` →
`RETROCESSION_VALUATION_*` dengan pergeseran **+1 hari**.

`[terverifikasi]` Keduanya membaca **salinan berbeda** dari nilai yang sama: gerbang Komite membaca
`pyWorkPage.Type`, `ValidasiDOL_Act` membaca `pyWorkPage.PolicyDataLife.Type`. Penyalinannya di
`Claim Life/Activity/LoadDataPeserta_Act.xml` — `pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type`.
Sumber otoritatif = `PolicyDataLife.Type`.
- **Memblokir:** tidak ada lagi. Seluruh butir terjawab; keputusan paritas tercatat di ADR-0012.
- **Terkait:** **ADR-0011** (mesin status), **ADR-0002** (RBAC tiga peran), **OQ-020** (arti kode
  `.Type` `QP`/`QR` **TETAP TERBUKA** — koreksi 2026-09-16: penutupan dari "dropdown EDM Life" dibatalkan,
  label Receivable/Payable NOL di korpus, hanya ada di instance JSON runtime; bukan bukti).
- **Catatan `[terverifikasi]`:** penegakan peran di modul ini memang tidak seragam —
  `pyPosition` hanya ada di 17 dari berkas `Claim Life`, dan **tidak ada** di
  `Section/InputOSClaimLife.xml`, `Section/RejectOSClaimLife_Sec.xml`, `Section/ClaimComite.xml`,
  maupun `Harness/Committe_Life.xml`. Penegakan sesungguhnya bertumpu pada penugasan tahap di
  `Claim Life/Flow/Register_Flow.xml`.
- **Status:** **TERJAWAB** (2026-09-14, sumber: work owner, Ronde 4). Bacaan gerbang ditetapkan;
  dugaan "cacat ekspresi" ditutup dengan bukti korpus; arti `TP`/`TR` dijawab; kelonggarannya
  diputuskan **paritas** dan dicatat sebagai risiko RBAC di **ADR-0012**.
  `[terverifikasi]` Gerbang seperti ini **hanya ada di `Claim Life`** — tidak ditemukan di 19 modul
  lain, sehingga penutupan ini tidak menyisakan cakupan modul lain.
- **Audit:**
  ```
  grep -n "pyPosition" "Claim Life/Section/AdjustmentDetail_Section.xml"
  grep -rc "pyPosition" "Claim Life" --include="*.xml" | grep -v ":0$"
  # sensus '=' vs '==' korpus-wide
  grep -rhoE "<pyCondition>[^<]*" . --include="*.xml" | grep -c "=="
  grep -rhoE "<pyCondition>[^<]*" . --include="*.xml" | grep -v "==" | grep -cE "[A-Za-z0-9_) ]=[ ']"
  ```

---

### OQ-064 — Tiga identitas retro ter-hardcode memicu EXIT; artinya tidak diketahui

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** grilling Komite Claim Life Ronde 1 (`.scratch/komite-claim-life/grilling-ronde-1.md`),
  2026-09-15 — saat mengoreksi klaim ronde yang salah atribusi
- **Bukti `[terverifikasi]`:** `Komite Claim Life/Activity/KomitePostAdjustment.xml`
  (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`) **step 9**
  — nomor langkah dibaca dari `<pyStepPageReference>RH_1.pySteps(9)` — membawa
  `<pyStepsDescription>EXIT JIKA RETROID "L0000141"</pyStepsDescription>` (baris 8474),
  `<pyStepsPreCondition>true`, dan dua baris precondition:

  ```
  baris 8516 : …PolicyDataLife.RetroID=="L0000141" || …SecurityReinsurerID=="L0000134"
  baris 8539 : …PolicyDataLife.RetroID=="1000013"
  ```

- **Akibat `[terverifikasi]`:** karena gerbang ini **EXIT** (menghentikan activity), bukan memilih
  cabang, klaim yang cocok **melewati step 9–14**:

  | Step | Efek | Nasib |
  | ---: | --- | --- |
  | 7 | `InsertJsonClaimLife_Act` | tetap berjalan |
  | 9 | `serviceInsertArasapasClaimLife_act` | **DI-SKIP** |
  | 10 | `SendEmailKlaimLife` | **DI-SKIP** |
  | 11 | `HitServiceToKasirKMTLife_Act` (Kasir) | **DI-SKIP** |
  | 13 | `SetInformationData` | **DI-SKIP** |

- **Pertanyaan:** Apa arti ketiga identitas `1000013`, `L0000141`, `L0000134`, dan **mengapa** klaim
  milik mereka sengaja tidak dikirim ke Arasapas, email, dan Kasir?
- **Status:** **TERTUTUP** (2026-09-15, sumber: **keputusan work owner**).

  `[keputusan work owner]` Ketiga identitas retro ter-hardcode — `1000013`, `L0000141`, `L0000134` —
  **beserta gerbang EXIT step 9** dibuang dari sistem baru. Tidak direplikasi sebagai konstanta.

  **Konsekuensi yang diterima work owner secara sadar:** klaim ber-retro tersebut yang selama ini
  di-EXIT — sehingga **tidak menjalankan efek keluar sama sekali** — kini menjalankan **seluruh**
  efek keluar, **termasuk Kasir (pembayaran)**. Setelah migrasi, **semua klaim menjalankan keempat
  efek keluar**, tanpa pengecualian.

  Arti ketiga identitas itu **tetap tidak diketahui** dan tidak dicari lagi: keputusannya membuang,
  bukan mereplikasi. Bila kelak diperlukan → jadikan data/konfigurasi seperti roster `EMAILKOMITE`.

  Syarat pengaman yang menyertainya ada di **ADR-0015**: ID idempoten unik per kiriman, dan
  **cek status "sudah terkirim sukses" sebelum kirim ulang** khusus Email dan Kasir.
  `[keputusan work owner 2026-09-15]` Ketiganya **dibuang** dari logika sistem baru; tidak
  direplikasi sebagai konstanta. Bila kelak diperlukan → jadikan data/konfigurasi seperti roster
  `EMAILKOMITE`.
- ⚠️ **Risiko yang menyertai keputusan itu:** membuang gerbang EXIT berarti klaim yang selama ini
  **dikecualikan** akan **mulai menerima keempat efek keluar** — termasuk **Kasir** (pembayaran).
  Ini **perubahan perilaku**, bukan pembersihan kode, dan berinteraksi langsung dengan **ADR-0015**
  (semua efek wajib berhasil). Layak dikonfirmasi ke Product+UW sebelum rilis.
- **Terkait:** **ADR-0015**, **ADR-0013**, OQ-031 (pola identitas ter-hardcode serupa di modul lain)
- **Audit:**
  ```
  f="Komite Claim Life/Activity/KomitePostAdjustment.xml"
  grep -n "pyStepPageReference>RH_1.pySteps(" "$f" | sed 's/<[^>]*>//g'
  grep -n "L0000141\|L0000134\|1000013" "$f"
  ```

### OQ-065 — Isi kolom retro pada rekam akseptasi setelah langkah "Tukar SecurityReinsurer dengan RetroName" dimatikan

- **Pemilik:** Product+Underwriting (+ DBA untuk kolom sasaran)
- **Ditemukan di:** grilling Komite Claim Life Ronde 2
  (`.scratch/komite-claim-life/grilling-ronde-2.md`), 2026-09-15 — saat memverifikasi status
  enable/disable dari berkas yang diperbarui work owner
- **Bukti `[terverifikasi]`:** `Komite Claim Life/Activity/KomitePostAdjustment.xml`
  (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, 515.675 byte,
  versi **2026-09-15**). Dua langkah berjudul **"Tukar SecurityReinsurer dengan RetroName"**
  ber-`<pyStepsBlockName>//` → **REMARK**:

  | Step | Rentang baris | Blockname `//` | Precondition di dalamnya |
  | --- | --- | ---: | --- |
  | **4.14** (jalur aksep) | 4937 – 5185 | 4947 | 5098 `Type=="TP"\|\|"TR"` · 5121 `SecurityReinsurerID!="" && SecurityReinsurer!=""` · 5144 `ProdDateTime<"20250207T000000.000 GMT"` |
  | **5.5** (jalur reject) | 7663 – 7911 | 7673 | 7824 · 7847 · 7870 (ketiganya sama) |

  Langkah inilah yang dahulu mengisi `TempInputDetail.CARI38` = `RetroID` dan `CARI39` = `RetroName`
  sebelum `Insert ke OS`.

- **Pertanyaan:** Karena kedua langkah itu mati, **apa yang kini masuk ke kolom retro pada
  `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`** — nilai mentah tanpa penukaran, kosong, atau diisi jalur lain
  yang belum ditemukan?
- **Memblokir:** pemetaan kolom retro pada rekam akseptasi Komite. **Tidak memblokir** mesin status,
  penomoran, wewenang, maupun efek keluar — ketiganya sudah tertutup.
- **Catatan:** ini **menguatkan** penutupan **OQ-034** (cutover 7 Feb 2025 tidak dipakai) dan
  memperluasnya: yang mati bukan hanya cabang cutover, melainkan **seluruh** langkah penukaran di
  kedua jalur. Ketiga precondition — termasuk dua yang sempat dianggap masih aktif — berada di dalam
  langkah yang ter-remark.
- **Status:** **TERTUTUP** (2026-09-15, sumber: **keputusan work owner**).

  `[keputusan work owner]` Kelima langkah ter-remark — **4.4** `Generate No Akseptasi (QP,QR)`,
  **4.5** `Generate No Akseptasi (TP,TR)`, **4.6** `Set Nilai Akseptasi`, **4.14** dan **5.5**
  `Tukar SecurityReinsurer dengan RetroName` — **sengaja dimatikan** karena nilai `RetroID` /
  `RetroName` **sudah di-set di langkah sebelumnya**, mentah dari data policy.

  **Jawaban atas pertanyaan:** nilai retro yang masuk ke `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` adalah
  **nilai apa adanya dari data policy, tanpa penukaran**.

  **Aturan sistem baru:** abaikan langkah tukar yang mati; pakai nilai retro **apa adanya**;
  **tidak ada logika penukaran yang direplikasi**.

  Ini sekaligus menutup sisa keraguan pada **OQ-034** — matinya cabang cutover 7 Feb 2025 bukan
  kelalaian, melainkan bagian dari penyederhanaan yang sama.
- **Audit:**
  ```
  f="Komite Claim Life/Activity/KomitePostAdjustment.xml"
  grep -n "<pyStepsBlockName>" "$f"                       # 5 langkah ter-remark
  grep -n "<pyStepPageReference>RH_1.pySteps(4).pySteps(1[4-5])<" "$f"
  grep -n "CARI38\|CARI39" "$f" | head
  ```

### OQ-066 — Dead code di `InsertJsonPolisLife_Act`: yang mati tertulis aktif, yang hidup tertulis mati

- **Pemilik:** Arsitektur Pega + Product+Underwriting
- **Ditemukan di:** grilling PremiumList Life Ronde 1
  (`.scratch/premiumlist-life/grilling-ronde-1.md`), 2026-09-15
- **Berkas:** `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml`
  (`ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 327.332 byte,
  tersimpan `20260728T024422`)

**Penanda enable/disable di berkas ini TIDAK dapat dipercaya sendirian.** Ada dua arah kekeliruan:

| Bagian | Penanda korpus | Kenyataan | Sumber |
| --- | --- | --- | --- |
| Step 6-7, empat gerbang `@contains(.ID,"1000032")` … `"1000035"` (QS / 2nd QS / SURPLUS / 2nd SURPLUS) | `<pyStepsBlockName/>` **kosong = AKTIF** | **tidak dipakai** — logika polis lama | `[keputusan work owner]` |
| Step 16 `Commit` (baris `//` 5294) | `<pyStepsBlockName>//` = **REMARK** | memang mati | `[terverifikasi]` |
| Step 17 `Connect-REST` `ConvertJsonNusareToProduction` (baris `//` 5383) | **REMARK** | memang mati | `[terverifikasi]` |

`[terverifikasi]` **17 langkah**, bukan 21 seperti catatan D2. Hanya **dua** `<pyStepsBlockName>`
berisi `//` di seluruh berkas.

- **Pertanyaan:** Mengapa keempat gerbang polis lama **tidak di-remark** padahal tidak dipakai — dan
  adakah langkah lain di modul Life yang aktif tetapi sebenarnya mati? Penanda `blockname` sudah
  terbukti **tidak lengkap** di berkas ini.
- **Memblokir:** **tidak memblokir isi** spec — work owner sudah menyatakan mana yang dipakai. Yang
  diblokir adalah **keandalan metode**: di modul lain, `blockname` tidak boleh dipakai sendirian
  untuk menyimpulkan langkah hidup/mati; perlu konfirmasi work owner.
- **Aturan yang berlaku sistem baru `[keputusan work owner]`:** JANGAN replikasi —
  (1) keempat gerbang ID polis lama; (2) `Commit` step 16 (Go memegang transaksi);
  (3) `ConvertJsonNusareToProduction` step 17.
- **Terkait:** **OQ-031** (arti keempat ID — tertutup), **OQ-013** (batas transaksi),
  **ADR-0013** (endpoint), **ADR-0005** (flag lingkungan)
- **Catatan versi:** `SubmitPremiumList_Act` (`20260122T072213`) membaca ambang tutup buku dari
  `POOLDATA.TANGGAL_CLOSING`, sedangkan `InsertJsonPolisLife_Act` (`20260728T024422`, **lebih baru**)
  masih memakai hardcode `>25`. Yang lebih baru justru yang lebih usang aturannya.
  `[keputusan work owner]` **ikuti yang dari DB** (OQ-030).
- **Status:** terbuka
- **Audit:**
  ```
  f="PremiumList Life/Activity/InsertJsonPolisLife_Act.xml"
  grep -n "<pyStepsBlockName>" "$f"                    # hanya 2 hasil: 5294, 5383
  grep -oE "<pyStepPageReference>RH_1\.pySteps\([0-9]+\)<" "$f" | sort -u   # 1..17
  grep -n "100003[2-5]" "$f"
  ```

### OQ-067 — `INSERTJSONOFFERLIFE` dipanggil di tahap **penawaran**, bukan di rantai simpan premium list

- **Pemilik:** Product+Underwriting + Arsitektur Pega
- **Ditemukan di:** penulisan tiket PremiumList Life, 2026-09-15

`[keputusan desain]` Aturan urutan procedure yang ditetapkan work owner berbunyi: satu transaksi Go
memuat `PROC_GENERATE_SEQUENCE_NUMBER` + `PEGA_M_LIFE_PREMIUM_SUMMARY` lalu commit, **kemudian**
`INSERTJSONPOLISLIFE`, **kemudian** `INSERTJSONOFFERLIFE` paling akhir.

⚠️ `[terverifikasi]` Di korpus, `INSERTJSONOFFERLIFE` **tidak berada di rantai itu sama sekali**. Ia
dipanggil satu tahap lebih awal:

| Pemanggil | Identitas | Langkah |
| --- | --- | --- |
| `PremiumList Life/Activity/InputOfferLife_ACT.xml` | `ASM-FW-GISFW-WORK-LIFE` / `INPUTOFFERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | step **4** `RDB-List` "Insert to table json_offer_life" → `SaveOfferJsonLife_SQL` |
| `PremiumList Life/RDBList/SaveOfferJsonLife_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!SAVEOFFERJSONLIFE_SQL` / `RULE-CONNECT-SQL` | `POOLDATA.INSERTJSONOFFERLIFE(...)` lalu `COMMIT;` (baris 141) |

`[terverifikasi]` Pencarian seluruh dua modul: **hanya** `InputOfferLife_ACT` yang merujuk
`SaveOfferJsonLife_SQL`. Rantai `InsertJsonPolisLife_Act` (PremiumList) merujuk `GetJsonProductLife`,
`InsertPLSummary`, `InsertJsonPolis`, `SaveLifeinProduction_SQL`, `GetNopolisByIDPega` — **tidak**
`SaveOfferJsonLife_SQL`.

- **Pertanyaan:** Apakah `INSERTJSONOFFERLIFE` tetap dipanggil di tahap penawaran (mengikuti korpus),
  atau sengaja dipindahkan ke akhir rantai simpan premium list (mengikuti aturan urutan)? Keduanya
  punya konsekuensi berbeda: di korpus, rekam offer sudah ter-commit sebelum premium list disentuh.
- **Memblokir:** **tidak memblokir** — tiket **01** menulisnya di tahap penawaran mengikuti korpus,
  dan aturan urutan diterapkan pada rantai simpan premium list (tiket **05a**/**05b**). Yang perlu
  dikonfirmasi hanyalah bahwa pemisahan itu memang yang dimaksud.
- **Terkait:** **OQ-013** (batas transaksi), **OQ-002** (body procedure), **ADR-0015**
- **Status:** terbuka
- **Audit:**
  ```
  grep -rl "SaveOfferJsonLife_SQL" "PremiumList Life/" "Endorsement Life/"
  grep -n "<RequestType>" "PremiumList Life/Activity/InsertJsonPolisLife_Act.xml"
  grep -n "COMMIT;" "PremiumList Life/RDBList/SaveOfferJsonLife_SQL.xml"
  ```

### OQ-068 — Tidak ada penulis `M_LIFE_PREMIUM_DETAIL` di PremiumList Life, padahal Claim Life membacanya lewat `PL_NUMBER`

- **Pemilik:** Arsitektur Pega + Product+Underwriting
- **Ditemukan di:** penulisan tiket PremiumList Life, 2026-09-15

> ⚠️ **Catatan versi.** Sensus di bawah ini adalah keadaan korpus **sebelum** work owner menambahkan
> `InsertLifePremiumDetail_act` dan salinan `SaveMasterLPDet` ke `PremiumList Life/` pada 2026-09-15.
> Ia dipertahankan apa adanya sebagai jejak audit. Keadaan terkini ada di **Penutupan** di bawah.

`[terverifikasi]` Sensus seluruh korpus atas tabel `POOLDATA.M_LIFE_PREMIUM_DETAIL` — **dua** berkas
saja:

| Berkas | Identitas | Peran |
| --- | --- | --- |
| `Endorsement Life/RDBList/SaveMasterLPDet.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | **satu-satunya penulis** — `INSERT INTO POOLDATA.M_LIFE_PREMIUM_DETAIL`, PK `M_LIFE_PREMIUM_DETAIL_SEQ.nextval`, lalu `COMMIT;` (baris 252) |
| `Claim Life/RDBList/GetPesertaClaim_sql1.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL` | **pembaca** — `SELECT * FROM POOLDATA.M_LIFE_PREMIUM_DETAIL WHERE PL_NUMBER = {…PremiumListSummary.PL_NUMBER}` |

`[terverifikasi]` Kolom `INSERT` memuat **`PL_NUMBER` dan `PL_NUMBER_EDM`** — jadi tabel yang sama
melayani kedua jalur. Tetapi `SaveMasterLPDet` hanya dirujuk `Endorsement Life/Activity/
InsertJsonPolisLife_Act.xml` (step **11.6**).

`[terverifikasi]` Di PremiumList Life, `PremiumListDetail` disusun sebagai properti work object
(`SavePremiumList_Act`, `ASM-FW-GISFW-WORK-LIFE` / `SAVEPREMIUMLIST_ACT`, step **8.4** "Set property
PremiumListDetail") lalu dipersistensikan `Obj-Save` (step **15**) — ke penyimpanan work Pega,
**bukan** ke `M_LIFE_PREMIUM_DETAIL`.

- **Pertanyaan:** Bagaimana baris detail **new business** sampai ke `M_LIFE_PREMIUM_DETAIL` sehingga
  Claim Life dapat menemukannya lewat `PL_NUMBER`?
- **Memblokir:** ~~tiket 08~~ — **tidak lagi memblokir.**
- **Terkait:** **OQ-023**, **OQ-013** (batas transaksi), **ADR-0001**, **ADR-0011**, **ADR-0003**
- **Status:** **TERTUTUP** (2026-09-15) — `[terverifikasi]` + `[keputusan work owner]`

---

### Penutupan OQ-068 (2026-09-15)

#### Bukti baru di korpus `[terverifikasi]`

Work owner menambahkan **dua** berkas ke jalur NB Life (keduanya bercap `20260211`):

| Berkas | Identitas | Tipe | Ukuran | `pxUpdateDateTime` |
| --- | --- | --- | ---: | --- |
| `PremiumList Life/Activity/InsertLifePremiumDetail_act.xml` | `ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` | `RULE-OBJ-ACTIVITY` | 286.827 byte | `20260211T064342.645 GMT` (ruleset `01-01-91`) |
| `PremiumList Life/RDBList/SaveMasterLPDet.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` | `RULE-CONNECT-SQL` | 15.537 byte | `20260211T064412.717 GMT` |

`[terverifikasi]` **`SaveMasterLPDet` adalah satu rule yang sama** dengan yang dipakai Endorsement,
bukan salinan yang menyimpang: identitas identik, `pxUpdateDateTime` identik, ukuran identik, dan
diff atas dua ekspor yang dinormalisasi menyisakan **8 baris** — seluruhnya cap waktu ekspor
(`pyRuleFormStatusTime`, `pyShowJavaWindowName`). **NB dan EDM sudah menulis lewat penulis yang sama;
yang berbeda hanyalah pemicunya.**

#### Peta langkah `InsertLifePremiumDetail_act` `[terverifikasi]`

Penomoran dari `<pyStepPageReference>` — perhatikan akarnya **`RH_2`**, bukan `RH_1`:

| Step | Langkah | Catatan |
| --- | --- | --- |
| 1 | `Property-Set` | set `Param.pyReportName = "SelectNoJsonPolis_RD"`, `Param.pyReportClass = "ASM-FW-GISFW-Work-LIFE"` (baris 531 dst.) |
| **2** | `Call Rule-Obj-Report-Definition.pxRetrieveReportData` | ⚠️ **inilah pemicu batch** — mengambil **daftar** kasus, lalu step 3 melooping `hasil.pxResults` |
| 3 | *(loop per hasil)* | |
| 3.1 | `Property-Set` | |
| 3.2 | `Obj-Open-By-Handle` | buka work object per baris hasil |
| 3.3 | "Insert to table detail" | |
| 3.3.1 | `Page-Remove` | |
| 3.3.2 | `Property-Set` "get ceding co name" | |
| 3.3.3 | `Property-Set` "insert nilai dari data-batch → int" | precondition `TempError.CARIDESC==1` (3634) |
| **3.3.4** | `RDB-List` "insert ke tabel detail premium list" → **`SaveMasterLPDet`** (3730) | precondition **`hasilDetail.pxResults(1).PL_NUMBER==""`** (3820) — **penjaga idempotensi** |
| 3.4 | `Property-Set` "Pega to jsondata" | |
| 3.5 | `Obj-Save` | |
| **3.6** | `Commit` | **AKTIF** |

`[terverifikasi]` **Tidak ada satu pun `<pyStepsBlockName>` di berkas ini** — nol langkah ter-remark.

`[terverifikasi]` Report Definition `SelectNoJsonPolis_RD` (class `ASM-FW-GISFW-Work-LIFE`) **dirujuk
tetapi tidak ada berkasnya di korpus** — asimetri rujukan. Pemilih batch itu sendiri tidak terekspor.

#### Keputusan `[keputusan work owner]`

1. **Logika penulisan TETAP `InsertLifePremiumDetail_act`** — tidak diganti, tidak ditulis ulang. Ia
   yang menulis peserta NB ke `M_LIFE_PREMIUM_DETAIL`.
2. **Yang disamakan dengan EDM adalah TIMING/pemicunya, bukan logikanya.** Di Pega existing, NB
   dipicu **job** (batch terjadwal — step 2 `pxRetrieveReportData` + loop step 3). Di sistem baru:
   **JANGAN pakai job.** Panggil logika itu **INLINE saat proses insert/simpan polis**, meniru cara
   **EDM/Endorsement** memicu penulisan detailnya (langsung di alur simpan, step 11.6
   `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml`).
3. **Penyimpangan sadar:** pemicu job dihilangkan; diganti pemanggilan inline di alur simpan polis.
   **Konsekuensi positif:** data peserta langsung tersedia untuk klaim **tanpa jeda job**, dan NB
   seragam dengan EDM.

**Akibat teknis yang mengikat:** step 2 dan loop step 3 **tidak dimigrasikan** — keduanya semata
mesin batch. Yang dimigrasikan adalah **badan per-kasus** (3.1–3.6), dipanggil sekali untuk kasus
yang sedang disimpan. Penjaga `PL_NUMBER==""` (3820) **wajib dipertahankan** agar pemanggilan ulang
tidak menggandakan baris. `Commit` step 3.6 tunduk pada aturan urutan transaksi campuran yang sudah
ditetapkan (spec §6, AC 20–24) — ⚠️ `SaveMasterLPDet` **commit sendiri** (`COMMIT;` baris 252),
sehingga ia **titik potong**, bukan bagian transaksi summary.

#### Audit

```
f="PremiumList Life/Activity/InsertLifePremiumDetail_act.xml"
grep -o "<pxInsName>[^<]*" "$f" | head -1
grep -n "<pyStepsBlockName>" "$f"                       # nol hasil
grep -n "<RequestType>" "$f"                            # 3730: SaveMasterLPDet
grep -n "pyStepsPreCondParamsWhen" "$f"                 # 3634, 3820
grep -rl "SelectNoJsonPolis" --include=*.xml .          # hanya berkas ini
find . -name "SaveMasterLPDet.xml" -not -path "./OUTPUT_HASIL_RNM/*"
```
### OQ-069 — Pesan validasi menjanjikan aturan "NET PREMIUM > GROSS PREMIUM" yang tidak pernah diperiksa

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** penulisan tiket PremiumList Life, 2026-09-15
- **Berkas:** `PremiumList Life/Activity/ValidasiUploadPL_act.xml`
  (`ASM-FW-GISFW-WORK-LIFE` / `VALIDASIUPLOADPL_ACT` / `RULE-OBJ-ACTIVITY`, **43 langkah**)

`[terverifikasi]` Pesan kesalahan berbunyi:

```
"NET PREMIUM HARUS ADA DAN LEBIH BESAR DARI GROSS PREMIUM"
"NET PREMIUM HARUS ADA DAN LEBIH BESAR DARI GROSS PREMIUM DAN SEPARATOR MENGGUNAKAN TITIK"
```

`[terverifikasi]` Namun **satu-satunya** precondition atas `NET_PREMIUM` di seluruh berkas adalah
`@PropertyHasValue(.NET_PREMIUM)`. **Tidak ada** perbandingan terhadap `GROSS_PREMIUM` di mana pun.

- **Pertanyaan:** Apakah aturan `NET_PREMIUM > GROSS_PREMIUM` memang aturan bisnis yang sah (dan
  implementasinya hilang), atau pesan itu sisa dari aturan yang dicabut? ⚠️ Arah perbandingannya pun
  janggal — lazimnya net **lebih kecil** dari gross.
- **Memblokir:** **tiket 04** (unggah CSV) — AC validasi tidak boleh menebak. Selama terbuka, tiket
  04 memeriksa **keberadaan** saja, persis seperti korpus, dan mencatat pesan Pega apa adanya.
- **Terkait:** **ADR-0003** (uang non-float)
- **Status:** terbuka
- **Audit:**
  ```
  f="PremiumList Life/Activity/ValidasiUploadPL_act.xml"
  grep -oE "NET PREMIUM[^<]{0,160}" "$f" | sort -u
  grep -oE "<pyStepsPreCondParamsWhen>[^<]*NET_PREMIUM[^<]*" "$f" | sort -u
  ```

### OQ-070 — `EDMStatus` punya **empat** nilai, dan beda `"Delete"` versus `"Batal"` belum diketahui

- **Pemilik:** Product+Underwriting
- **Ditemukan di:** grilling Endorsement Life Ronde 2
  (`.scratch/endorsement-life/grilling-ronde-2.md` §C2), 2026-09-15

`[terverifikasi]` Sensus penulisan `EDMStatus` di seluruh modul `Endorsement Life/`:

| Nilai | Penulis | Identitas | Baris |
| --- | --- | --- | ---: |
| `"Old"` | `MappingEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` / `RULE-OBJ-ACTIVITY` | 2608 |
| `"New"` | `SaveCSVEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE` / `RULE-OBJ-ACTIVITY` | 2714 |
| `"Delete"` | `SetPremi_EDM` step 2.2 "flag pengurangan" | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY` | 1384 |
| **`"Batal"`** | **`SetPremi_EDM` step 2.3 "flag batal"** | idem | **1506** |

⚠️ Verdict **V7 Ronde 1** `[keputusan work owner]` menyatakan **tiga** nilai lengkap. Korpus
menunjukkan **empat**. Selisih ini **tidak ditutupi**; verdict V7 sudah diberi catatan koreksi.

`[terverifikasi]` Kedua nilai batal ditulis dengan precondition **berbeda**: `"Delete"` saat
`.EdmBatal=="True"` (baris 1441); `"Batal"` saat `pyWorkPage.EdmType==3` (1583) atau
`.EdmBatal=="True"` (1636).

`[terverifikasi]` `EDMStatus` mengalir ke rantai simpan — `Endorsement Life/Activity/
InsertJsonPolisLife_Act.xml` baris 4455 menyalin `.EDMStatus` → `TempValue.EDMStatus`.

- **Pertanyaan:** Apa beda **`"Delete"`** (baris dihapus pada endorsement perubahan data) dan
  **`"Batal"`** (baris dari polis yang dibatalkan)? Apakah keduanya diperlakukan berbeda di hilir,
  khususnya oleh Claim Life yang membaca `M_LIFE_PREMIUM_DETAIL`?
- **Memblokir:** **spec Endorsement Life** — mesin status baris detail tidak dapat ditulis tanpa ini.
- **Terkait:** **OQ-071** (nilai negatif), **ADR-0011** (mesin status), **ADR-0001** (kontrak hilir)
- **Status:** **TERTUTUP** (2026-09-15) — `[keputusan work owner]`, verdict **V13** Ronde 2

  **Empat nilai, dan bedanya adalah CAKUPAN MINUS — bukan jenis:**

  | Nilai | Kapan | Akibat pada nilai uang |
  | --- | --- | --- |
  | `"Old"` | baris warisan polis new business | tidak diubah |
  | `"New"` | peserta **ditambah** — hanya pada Perubahan Data (`EdmType=1`) | positif, baris baru |
  | `"Delete"` | peserta **dihapus** dalam Perubahan Data | **diminuskan — selektif PER PESERTA** |
  | `"Batal"` | lewat EDM Batal (`EdmType=3`) | **seluruh peserta otomatis batal — diminuskan MENYELURUH** |

  ⚠️ **Mengoreksi anggapan sebelumnya** bahwa `"Delete"` sekadar penanda: **`Delete` JUGA
  menghasilkan nilai negatif.** Keduanya menempuh jalur yang sama — jurnal balik `× -1` di
  `SetPremi_EDM`. Yang tetap benar: baris polis NB **tidak dihapus fisik**.
- **Audit:**
  ```
  grep -rn "PropertiesName>\.EDMStatus" "Endorsement Life/"
  sed -n '1380,1390p;1502,1510p' "Endorsement Life/Activity/SetPremi_EDM.xml"
  grep -n "pyStepsPreCondParamsWhen" "Endorsement Life/Activity/SetPremi_EDM.xml"
  ```

### OQ-071 — Endorsement batal menulis nilai **negatif** (32 kolom uang `× -1`), bukan nol

- **Pemilik:** Finance + Product+Underwriting
- **Ditemukan di:** grilling Endorsement Life Ronde 2
  (`.scratch/endorsement-life/grilling-ronde-2.md` §C1), 2026-09-15
- **Berkas:** `Endorsement Life/Activity/SetPremi_EDM.xml`
  (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY`, 345.689 byte,
  **9 langkah, nol `<pyStepsBlockName>`**)

⚠️ `[terverifikasi]` **Deskripsi langkah bertentangan dengan kodenya.** Step 2.1 berdeskripsi
**"Set 0 jika EDM Batal"** dengan precondition
`.EdmBatal=="True" || pyWorkPage.EdmType==3` (baris 1273) — tetapi kodenya **tidak menulis nol**.
Ia menulis `<kolom> = <kolom> * -1` pada **32 kolom uang**:

```
SUM_INSURED, CEDING_RETENTION, SUM_REASURED, SHARE_NUSANTARA_RE, SHARE_NUSANTARA_RE_GROSS,
SUM_AT_RISK_GROSS, SUM_AT_RISK_RETRO, RETROCEDED_SHARE, SHARE_RETRO, RATE, FACTOR,
GROSS_PREMIUM, NET_PREMIUM, DEDUCTION, CLAIM_AMOUNT, RI_ADMIN_FEE, BROKERAGE_FEE,
dan kelompok *_REFUND / *_RETRO / *_REFUND_RETRO
```

Hasilnya **jurnal balik (counter-entry)**: baris bernilai negatif yang meniadakan baris asli saat
dijumlahkan — bukan penimpaan dengan nol.

- **Pertanyaan:** (a) Apakah pembalikan tanda memang perilaku akuntansi yang dikehendaki?
  (b) Apakah nilai negatif itu memang **masuk** ke `M_LIFE_PREMIUM_DETAIL` dan
  `M_LIFE_PREMIUM_SUMMARY`, dan terbaca Claim Life?
- **Memblokir:** **spec Endorsement Life** dan **kontrak hilir**. Menebak di sini menentukan apakah
  saldo premium list setelah pembatalan menjadi **nol** atau justru **terhitung dua kali**.
- **Terkait:** **OQ-070**, **ADR-0003** (uang non-float — nilai negatif wajib tetap presisi penuh),
  **ADR-0001**, **ADR-0011**
- **Status:** **TERTUTUP** (2026-09-15) — `[keputusan work owner]`, verdict **V13 + V14** Ronde 2

  (a) **Ya — pembalikan tanda memang perilaku yang dikehendaki.** Endorsement batal (dan hapus
  peserta) menghasilkan **jurnal balik negatif**, bukan penimpaan dengan nol. Deskripsi langkah
  Pega **"Set 0 jika EDM Batal"** **menyesatkan**; kodenya (`× -1` pada 32 kolom uang) yang berlaku.

  (b) **Ya — baris negatif masuk ke tabel yang SAMA**: `POOLDATA.M_LIFE_PREMIUM_DETAIL` dan
  `POOLDATA.M_LIFE_PREMIUM_SUMMARY`. Baris positif asli dan baris negatif **hidup berdampingan**;
  **net akunting** diperoleh dari penjumlahan, bukan dari penghapusan.

  ⚠️ **Penajaman kontrak hilir:** peserta yang sudah **EDM Batal** atau **soft-delete**
  **TIDAK BOLEH MUNCUL** di Claim Life. Kueri klaim
  `Claim Life/RDBList/GetPesertaClaim_sql1.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` /
  `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL`, berkunci `PL_NUMBER`) **wajib menyaring keluar**
  peserta batal/delete. Akuntansi melihat seluruh baris; klaim hanya peserta yang masih hidup.
  Menyentuh **tiket 08 PremiumList Life** dan **spec Claim Life**.
- **Audit:**
  ```
  f="Endorsement Life/Activity/SetPremi_EDM.xml"
  grep -n "Set 0 jika EDM Batal" "$f"
  sed -n '565,1275p' "$f" | grep -c '\* -1'        # 32
  sed -n '1273p' "$f"
  ```

### OQ-072 — Perilaku Pega saat unggahan melewati **50.000 baris** tidak terbaca

- **Pemilik:** Product+Underwriting + Arsitektur Pega
- **Ditemukan di:** grilling Endorsement Life Ronde 2, 2026-09-15
- **Berkas:** `Endorsement Life/Activity/SetPremi_EDM.xml` (identitas di **OQ-071**), step 3

`[terverifikasi]` Gerbang baris 3373:
`@SizeOfPropertyList(TempWorkPage.ListLifePremiumDetailUpload)>50000`, dengan
`<pyStepsPreCondParamsWhenTrue>3</…>` dan `<pyStepsPreCondParamsWhenFalse>2</…>`.

⚠️ **Yang tidak diketahui:** apa yang **terjadi** saat batas terlewati — kode transisi `3` tidak
dapat diterjemahkan karena **pemetaan kode transisi Pega tidak ada di korpus**. Terkait catatan
`[dugaan]` bahwa `6` = Exit Activity (grilling Ronde 1 §Bagian E) yang juga belum terbukti.

- **Pertanyaan:** Apakah `50000` adalah batas bisnis yang sah, dan apa yang seharusnya terjadi saat
  dilewati — ditolak, dipotong, atau diproses bertahap?
- **Memblokir:** **tidak memblokir spec.** Rekomendasi desain sudah ada: pertahankan batas, jadikan
  **konfigurasi** (jangan ditanam), dan **tolak terang-terangan** dengan pesan yang menyebut batas
  serta jumlah baris — sejalan dengan keputusan "gagal terang-terangan" pada ambang tutup buku
  PremiumList Life.
- **Terkait:** **OQ-066** (penanda korpus tidak dapat dipercaya sendirian)
- **Status:** **TERTUTUP** (2026-09-15) — `[keputusan work owner]`, verdict **V17** Ronde 2

  **Batas 50.000 baris DIBUANG.** Sistem baru **tidak membatasi** jumlah baris unggahan CSV
  endorsement. ⚠️ **Penyimpangan sadar.**

  **Catatan teknis:** unggahan sangat besar diproses **bertahap** (streaming/batch internal) —
  **tanpa menolak karena jumlah baris**. Pertanyaan "apa yang terjadi saat batas terlewati" karena
  itu **gugur**: tidak ada lagi batas yang bisa dilewati.

---

## Temuan Endorsement Life yang **tidak** menjadi OQ (sudah diputus) — 2026-09-15

Dicatat di sini agar register tidak ditanyai ulang. Rincian di
`.scratch/endorsement-life/grilling-ronde-1.md` dan `grilling-ronde-2.md`.

| Hal | Putusan | Tanda |
| --- | --- | --- |
| Arti `EdmType` | **`1` = Perubahan Data, `3` = Batal**; `2`/`4` tidak dipakai. Nilai `3` **terbukti kode** (`SetPremi_EDM` precondition `pyWorkPage.EdmType==3`, baris 1273) | `[keputusan work owner]` + `[terverifikasi]` |
| Arti `IVD_JR_ID = '5'` (`ARASAPAS.DETAIL_INVOICE`) | **pembayaran/pelunasan**; sistem baru tetap membaca skema Arasapas langsung, **dikurung di satu repository** bertanda batas lintas sistem | `[keputusan work owner]` + `[keputusan desain]` |
| Tidak ada `ViewOldPolicy_EDM_QR` | **`QR` didukung lewat section induk `ViewOldPolicy_EDM`** — kini **terbukti korpus**: `ShowLifePremiumSummary_EDM` memasangkan harness induk dengan `<pyCondition>.Type=='QR'` (baris 65394), `_QP`→`'QP'` (65952), `_TR`→`'TR'` (66510), `_TP`→`'TP'` (67068) | `[terverifikasi]` |
| `GetOldDetail_EDM` | **rule tidak dimigrasikan.** ⚠️ Tetapi **perilakunya masih terpakai** — ia dijalankan (`Run Activity`) oleh keempat tombol "lihat polis lama" sebelum popup `ViewOldPolicy_EDM*` tampil. Pemuatan data polis lama wajib tetap ada | `[keputusan work owner]` + `[terverifikasi]` |
| `CreateCaseEMDL` step 7–11 | **mati** — tidak dimigrasikan (`//` di 1084, 1273, 1382, 1532, 1671; deskripsi diawali `--`) | `[keputusan work owner]` |
| `UploadCSVEDMLifePremium_Act` step 4–6 | **mati** (`//` di 645, 791, 2120, 2400, 2525); hanya step 1–3 (`Call pxUploadCSVResults`) yang hidup | `[terverifikasi]` |
| Alarm jalur endorsement | **dihidupkan** di sistem baru (deteksi separuh step 13 + `SendEmailNotification` step 15 yang di Pega REMARK), seragam dengan new business | `[keputusan work owner]` |
| Pembacaan `PRODKE` | **disatukan ke `ORDER BY PRODKE DESC`** di kedua pembaca (`GetProdkeNopolis` sudah DESC; `GetProdKeOldData_SQL` memakai `TGL_INPUT desc`) — penyimpangan sadar | `[keputusan desain]` |
| Anti-dobel `JSON_POLIS` | `InsertJsonPolisEDM` adalah **`INSERT` polos + commit sendiri**; sistem baru memberi **penjaga idempotensi** `(NOPOLIS, PRODKE)` — penyimpangan sadar | `[keputusan desain]` |
| Pembatalan endorsement | **tidak ada tombol batal** — membatalkan = **`Decline`** (`Resolved-Rejected`); polis bebas di-endorse ulang. `CancelCreateCaseEDML` (1 langkah) bukan mekanisme batal | `[keputusan work owner]` |
| Editabilitas per form | **bukan rule `When`** — kondisi sebaris `<pyDisabledWhen>` atas `.EditInput` / `.EditInput1` / `.IsJsonPolis`; `.EditInput` adalah **kunci satu arah** (tiga penulis, semuanya ke `1`, tidak ada yang mengembalikan ke `0`) | `[terverifikasi]` |
| `Calculate1_Act` | ⚠️ **BUKAN mesin bersama.** Berkasnya identik di kedua folder (`ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` / `RULE-OBJ-ACTIVITY`, 302.597 byte, diff ternormalisasi nol baris), **tetapi di jalur ENDORSEMENT tidak dijalankan** — milik PremiumList Life saja. Perhitungan endorsement lewat `SetPremi_EDM`. **Pelajaran: berkas identik ≠ dipakai** (seiring OQ-066) | `[keputusan work owner]` |
| `EDMStatus` empat nilai | **Cakupan minus**: `Delete` = minus selektif **per peserta** (Perubahan Data); `Batal` = minus **menyeluruh** (`EdmType=3`). ⚠️ `Delete` **juga** negatif, bukan sekadar penanda | `[keputusan work owner]` — OQ-070 |
| Baris negatif & kontrak hilir | Baris negatif masuk **tabel yang sama**, berdampingan dengan baris positif → net akunting. ⚠️ **Peserta batal/delete TIDAK muncul di Claim Life** — jalur baca klaim wajib menyaring | `[keputusan work owner]` — OQ-071 |
| Batas 50.000 baris CSV | **Dibuang** — tanpa batas; unggahan besar diproses bertahap | `[keputusan work owner]` — OQ-072 |
| Kunci field permanen | `.PolicyNo`, `.EdmTypeBatal`, `.EdmBatal` terkunci permanen setelah endorsement dibuat; sistem baru **pertahankan + beri pesan penjelas** | `[keputusan work owner + desain]` |
| Popup polis lama | **Tetap ada**; rule `GetOldDetail_EDM` tidak ditiru, **perilakunya wajib ada** | `[keputusan work owner]` |

⚠️ **Penerapan OQ-066 di seluruh temuan di atas:** setiap penetapan "mati" bersandar pada
**keputusan work owner**; `<pyStepsBlockName>` dipakai sebagai **pendukung**, tidak pernah sebagai
dasar tunggal.

## Terjawab

### OQ-004 — `PremiumList Life/InputPolicyHolder.xml` berada di luar struktur folder tipe rule

- **TERJAWAB (2026-09-12, STEP D1 batch 4).** Pemilik semula: Product+Underwriting.
- **Pertanyaan asli:** tipe rule apa file ini, dan mengapa tidak di folder tipenya?
- **Jawaban `[terverifikasi]`:** file itu rule **`Flow`**. Dua sumber tipe sepakat —
  `<pzOriginalInstanceKey>` = `RULE-OBJ-FLOW ASM-FW-GISFW-WORK-LIFE INPUTPOLICYHOLDER
  #20180807T032319.287 GMT` dan `<pxObjClass>` = `Rule-Obj-Flow`.
  Identitas: `ASM-FW-GISFW-WORK-LIFE / INPUTPOLICYHOLDER`. `<pyStartActivity>` = `Start1`.
  Ukuran 106.880 byte.
- **Dampak:** **mengoreksi OQ-005** — `PremiumList Life` ternyata **punya** rule `Flow`, sehingga
  modul tanpa `Flow` turun dari 6 menjadi 5, dan jumlah rule `Flow` korpus naik dari 23 menjadi 24.
  Untuk D2, modul ini dapat ditelusur dengan metode standar.
- **Menyisakan:** **mengapa** file ditempatkan di luar struktur folder — pertanyaan proses ekspor,
  bukan isi rule. Tidak memblokir apa pun, tidak dilacak sebagai OQ terpisah.
- **Perintah audit:** `find . -maxdepth 2 -name "*.xml" -not -path "./OUTPUT_HASIL_RNM/*"` lalu
  baca `<pzOriginalInstanceKey>` dan `<pxObjClass>` file tersebut.

### OQ-006 — Nama Flow identik di beberapa modul: satu proses atau bercabang

**TERJAWAB PENUH (2026-09-12, STEP D1 batch 4).** Seluruh 24 rule `Flow` korpus sudah diinventarisasi
dan dibandingkan dengan hash ternormalisasi 18 tag.

**Jawaban `[terverifikasi]`** — nama file yang sama ternyata **campuran** antara rule yang sama dan
rule yang benar-benar berbeda; pembedanya adalah **class**, bukan nama:

| Nama file Flow | Class | Kesimpulan |
| --- | --- | --- |
| `OfferFacOut` | `ASM-FW-GISFW-WORK` (NB FacIn, RNW Fac In) | **satu rule**, isi identik |
| `InputInwardFacultativeRISlip` | `ASM-FW-GISFW-WORK` (NB FacIn, RNW Fac In) | **satu rule**, isi identik |
| `InputRealizationTreatyIn` | `ASM-FW-GISFW-WORK` (NB FacIn, NB Treaty In) | **satu rule**, isi identik |
| `OfferFacRetro` | `ASM-FW-GISFW-WORK` (3 modul) | **satu identitas, DUA isi**: NB FacIn = RNW Fac In, tetapi **Endorsment Fac In berbeda** → masuk OQ-011 |
| `Register_Flow` | `...WORK-PNC` (Claim Fac In) vs `...WORK-CLAIMLIFE` (Claim Life) | **dua rule berbeda** |
| `Flow_TreatyIn` | `...WORK-CLAIMTREATYNONPROP` vs `...WORK-CLAIMTREATY` | **dua rule berbeda** |
| `KomiteTreaty_Flow` | `...WORK-KOMITETREATYNONPROP` vs `...WORK-KOMITETREATY` | **dua rule berbeda** |

**Konsekuensi untuk D2:** dari 24 file `Flow`, jumlah proses berbeda yang perlu ditelusur lebih
sedikit dari 24 — tetapi `OfferFacRetro` **harus ditelusur dua kali** karena isinya bercabang.

- **Terjawab sebagian (2026-09-12, STEP D1 batch 1)** — pengukuran awal, untuk batch 1 saja.
- **Temuan:** dari 387 identitas rule yang muncul di >1 modul, **257 identik isinya** (satu rule
  yang sama terekspor berulang) dan **130 berbeda isinya** (varian nyata).
- **Catatan metode yang penting:** perbandingan `md5sum` mentah menyatakan seluruh 387 berbeda dan
  itu **keliru** — ekspor Pega mengacak urutan elemen dan menyisipkan timestamp. Perbandingan wajib
  dilakukan setelah normalisasi (buang `pyRuleFormStatusTime`, `pyShowJavaWindowName`,
  `pxUpdateDateTime`, `pxCreateDateTime`, `pxCommitDateTime`, `pxSaveDateTime`,
  `pyRuleAvailableTime`; urutkan baris; lalu hash).
- **Menyisakan:** OQ-011 (versi mana yang berlaku di production untuk 130 yang berbeda) dan
  pertanyaan yang sama untuk batch 2–4.

**Konfirmasi STEP D2 Tahap 4 (2026-09-13) `[terverifikasi]`.** Prediksi entri ini terbukti tepat di
lapangan. Pengukuran ulang seluruh 12 file `Flow` ketiga modul facultative + `NB Treaty In`:

| Hash | File | Ada di | Perlakuan D2 |
| --- | --- | --- | --- |
| `1306f68d56` | `InputRealizationTreatyIn` | NB FacIn = NB Treaty In | **dipakai ulang**, tidak ditelusur ulang |
| `c099bf4ebb` | `InputInwardFacultativeRISlip` | NB FacIn = RNW Fac In | ditelusur **sekali** |
| `832fb12b9d` | `OfferFacOut` | NB FacIn = RNW Fac In | ditelusur **sekali** |
| `bfd6252070` | `OfferFacRetro` | NB FacIn = RNW Fac In | varian A |
| **`b370146c63`** | `OfferFacRetro` | **Endorsment Fac In** | varian B — **ditelusur terpisah** |

`OfferFacRetro` memang harus ditelusur dua kali, dan selisihnya **nyata, bukan kosmetik**: varian A
punya tangga persetujuan fac out **2 anak tangga** (`ReasFacOutAdmin` → `ReasFacOutHead`); varian B
punya **4** (+ `ReasFacOutGroupLeader` → `ReasFacOutTechnicalDirector`), dan seluruh `reject` di
tangga ≥2 kembali ke Admin. Bukti: `flows/_SUMMARY-facultative.md` §4,
`flows/Endorsment Fac In.md` §2.1.
