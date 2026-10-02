# Discovery Endorsement — E-5: jalur produksi endorsement

> **Sumber:** `D:\migrasi\RNM\Endorsment Fac In\` (banding `NB FacIn\`). READ-ONLY.
> Korpus Treaty (`RNM_BRD\`) **tidak dibaca** (K-005). Label mengikuti `CLAUDE.md` §3.
>
> ⚠️ **E-5 sebagian MEMBLOKIR.** Jalur produksi NB masih menunggu tabel flat + `ALL_SOURCE` dari DBA.
> Dokumen ini memetakan **apa yang ada di korpus EDM**; yang bergantung pada stored procedure atau
> tabel di luar korpus ditandai `[pertanyaan terbuka]` dan **tidak ditebak**.
>
> ⚠️ **Privasi.** Alamat email literal **dihitung, tidak disalin** (K-025 / E-Q15).
>
> Kontrak volatile yang berlaku: **23 tag + isi blok `pzIndexes` diabaikan** (amandemen K-043).

---

## 1. Rantai produksi — langkah demi langkah

### 1.1 `SaveEDMToJsonPolicy_Act` — 29 langkah

`[terverifikasi]` `Endorsment Fac In\Activity\SaveEDMToJsonPolicy_Act.xml`, langkah tingkat atas:

| # | Metode | Deskripsi / gerbang |
| ---: | --- | --- |
| 1–3 | `Obj-Refresh-And-Lock` · `Page-Remove` · `Property-Set` | *"Set ID pega"* |
| 4–5 | `Property-Set` · `Page-Set-Messages` | galat bila `QuotationData.Type` kosong — `IF Local.errEDMType=="4" && Local.cekTypeEDM==""` |
| 6–7 | `RDB-List` `GETTanggalClosing_SQL` · `Property-Set` | tanggal closing |
| **8** | — | *"SET PRODDATETIME"* · `IF …PolicyData.EndorsementNo==""` · ⛔ `pyStepsPreCondition` = **`false`** |
| 9 | `RDB-List` `GetPolicyNoByCaseId` | ambil nopolis |
| **10** | — | ⛔ **gerbang idempotensi JSON_POLIS** — `IF @LengthOfPageList(EndorsementFacIn.pxResults)>=1` → **T=6 (keluar activity)**, F=3. Deskripsi: *"kalau data udh ada d jsonpolis exit act"* |
| 11 | `Page-New` | |
| **12** | `RDB-List` | *"Generate No Endorsement"* · `RequestType` = **`GenerateEndorsementNo`** · `IF IsLife` **[T=3 F=2]** |
| **13** | `RDB-List` | *"Generate No Endorsement LIFE"* · `RequestType` = **`GenerateEDMNoLife`** · `IF IsLife` **[T=2 F=3]** |
| 14 | — | *"Nopolis bukan LIFE"* · `GetKodeProdNonLife_SQL` · `IF EndorsementNo==""` |
| 15 | — | *"Nopolis LIFE"* · `GetKodeProdLife_SQL` · `IF EndorsementNo==""` |
| 16 | `RDB-List` | *"generate MM.YYYY DAN SEQUENCE"* · **`GetSequenceNumber_SQL`** |
| 17–18 | `Property-Set` | cabang `EndorsementNo==""` vs `!=""` |
| 19 | `Property-Set` | *"Get IDPega \| InputData.CARI1"* · ⛔ `pyStepsPreCondition` = **`false`** |
| **20** | `RDB-List` | *"Get ProdKe dari JSON_POLIS"* · **`GetProdKeOldData_SQL`** · ⛔ `pyStepsPreCondition` = **`false`** |
| **21** | `Property-Set` | `InputData.CARI3 = @ASM.GetPageJSONString()` ← muatan JSON |
| **22** | `Property-Set` | `Local.Prodke = OldData.pxResults(1).HASIL2` · **`InputData.CARI4 = Local.Prodke+1`** · `InputData.CARI5 = OldPolicyNo` |
| **23** | `RDB-List` | **tulis JSON_POLIS** — `RequestType` = **`INSERTJSON_JSONPOLISEDM_FACIN`** |
| 24 | `Call InsertCedingProduction` | |
| 25–26 | `Property-Set` · `Obj-Save` | `OutputData.HASIL3 = "step akhir"` |
| **27** | **`Call SaveTreatyProduction_Act`** | *"To save data in trearty Production"* |
| 28–29 | `Call SaveFacinLive_Act` · `Call SaveFacinSpreadLife_Sql` | *"-- Untuk LIfe"* |

### 1.2 ⛔ Langkah 12 dan 13 — gerbang sama, **transisi terbalik**

`[terverifikasi]` Keduanya bergerbang `IsLife`, tetapi:

| Langkah | `RequestType` | WhenTrue | WhenFalse | Berjalan saat |
| ---: | --- | ---: | ---: | --- |
| 12 | `GenerateEndorsementNo` | **3** | 2 | `IsLife` **salah** (non-Life) |
| 13 | `GenerateEDMNoLife` | 2 | **3** | `IsLife` **benar** |

📌 **Ekspresi gerbang saja tidak cukup — arah ditentukan `WhenTrue`/`WhenFalse`.** Ini melengkapi
pelajaran §3.3 `CLAUDE.md`: bukan hanya nama dan label yang bukan bukti, **ekspresi kondisi tanpa
transisinya pun belum menentukan perilaku**.

✅ **Uji integritas E-3.** Seluruh gerbang `EdmType` di `CountEndorsementData` dan `CountDataEDMElse`
diperiksa ulang: **19 dari 19** bertransisi seragam (`WhenTrue` kosong, `WhenFalse=3`). Tidak ada
pembalikan. **Kesimpulan E-3 §5 tidak terpengaruh.**

### 1.3 `SaveTreatyProduction_Act` — gerbang idempotensi kedua

`[terverifikasi]`

```
[3] RDB-List  « get idpega from facinproduction »   RequestType = CekFacin_Sql
[4] « CEK APAKAH SUDAH ADA DI FACINPRODUCTION »
      IF: CekFacin.pxResults(1).CARI1==""   [T=2 F=6]      ← F=6 = KELUAR activity
[5] Call SaveFacinRNWProd_Act        IF @contains(pyWorkIDPrefix,"RNW-")
[6] Call SaveFacinProdAllEDM_Act     IF @contains(pyWorkIDPrefix,"EDM-")
[7]…[24]  jalur NB, seluruhnya bergerbang @contains(pyWorkIDPrefix,"NB-")
```

⛔ **Dua gerbang idempotensi berbeda lapis:**

| Lapis | Di mana | Uji | Akibat bila sudah ada |
| --- | --- | --- | --- |
| JSON_POLIS | `SaveEDMToJsonPolicy_Act` langkah 10 | `EndorsementFacIn.pxResults >= 1` | keluar activity (T=6) |
| `facinproduction` | `SaveTreatyProduction_Act` langkah 4 | `CekFacin.pxResults(1).CARI1 == ""` | keluar activity (F=6) |

`[terverifikasi]` `RDBList\CekFacin_Sql.xml` (`pxObjClass` = `Rule-Connect-SQL`):

```sql
select distinct(IDPEGA) as CARI1 FROM pooldata.facinproduction
WHERE IDPEGA={DataSearch.CARI13}
```

📌 Satu case = satu `IDPEGA`; penulisan ulang dicegah di dua titik. **Ini yang membuat rantai aman
diulang** — penting untuk rekonsiliasi paralel run.

---

## 2. Tujuh cabang lini — dan satu yang hilang

### 2.1 `SaveFacinProdAllEDM_Act`

`[terverifikasi]` 8 langkah; langkah 1 menyiapkan parameter bersama, langkah 2–8 bercabang:

| # | Gerbang `When` | Memanggil | Berkas ada? |
| ---: | --- | --- | --- |
| 2 | `IsFire` | `SaveFacinProdEDMFire_Act` | ✅ |
| 3 | `IsAneka` | `SaveTreatyProductionEDMAneka_Act` | ✅ |
| **4** | `IsBonding` | **`SaveFacinProdEDMBonding_Act`** | ⛔ **NIHIL di korpus** — §2.2 |
| 5 | `isGolfInsurance` | `SaveFacinProdEDMGolf_Act` | ✅ |
| 6 | `IsMarineCargo` | `SaveFacinProdEDMMarineCargo_Act` | ✅ |
| 7 | `IsMBU` | `SaveFacinProdEDMMBU_Act` | ✅ |
| 8 | `IsPA` | `SaveFacinProdEDMPA_Act` | ✅ |

⛔ **Tidak ada cabang Life maupun Travel.** Life ditangani terpisah di `SaveEDMToJsonPolicy_Act`
langkah 28–29 (`SaveFacinLive_Act`, `SaveFacinSpreadLife_Sql`). **Travel tidak punya jalur produksi
sama sekali** di rantai ini — meskipun `SetOldData` menanganinya (E-3 §3.2). → **E-Q31**.

### 2.2 ⛔ `SaveFacinProdEDMBonding_Act` — hilang dari korpus, **ada di `DDL\`**

`[terverifikasi]` Status per lokasi:

| Lokasi | Status |
| --- | --- |
| `NB FacIn\` | NIHIL |
| `RNW Fac In\` | NIHIL |
| `Endorsment Fac In\` | **NIHIL** |
| **`D:\migrasi\RNM\DDL\`** | **ADA (670.777 B)** |

⚠️ **Karena itu ini BUKAN kasus `panic` seperti `GenerateEndorsementNo`.** Berkasnya ada — di folder
yang work owner pakai untuk melengkapi ekspor (pola yang sama dengan amandemen K-006 butir B.1).
Isinya **tidak saya baca** (di luar batas kerja putaran ini). → **E-Q28**: apakah dilipat masuk
sebagai amandemen, seperti B.1 dulu.

⛔ Sesuai **K-006**, cabang `IsBonding` **tidak dihapus dan tidak divonis usang**.

### 2.3 Ukuran dan cacah penulisan per cabang

`[terverifikasi]`

| Cabang | Ukuran | Langkah | `InsertTreatyProduction_Sql` | `…Backup_Sql` | `INSERTERRORFACIN_Sql` |
| --- | ---: | ---: | ---: | ---: | ---: |
| `SaveFacinProdEDMFire_Act` | 1.246.060 B | 16 | **6** | 4 | 4 |
| `SaveTreatyProductionEDMAneka_Act` | 806.152 B | 15 | 4 | 2 | 2 |
| `SaveFacinProdEDMGolf_Act` | 731.548 B | 15 | 4 | 2 | 2 |
| `SaveFacinProdEDMMarineCargo_Act` | 728.795 B | 15 | 4 | **4** | **4** |
| `SaveFacinProdEDMMBU_Act` | 723.069 B | 15 | 4 | 2 | 2 |
| `SaveFacinProdEDMPA_Act` | 714.920 B | 14 | 4 | 2 | 2 |

📌 FIRE paling besar dan punya paling banyak titik `INSERT`. Marine Cargo satu-satunya non-FIRE yang
punya 4 titik backup + 4 titik galat.

---

## 3. Kolom ← parameter ← sumber nilai

### 3.1 Tabel `facinproduction` — 82 kolom

`[terverifikasi]` `Endorsment Fac In\RDBList\InsertTreatyProduction_Sql.xml` ·
`pxObjClass` = **`Rule-Connect-SQL`** (bukan tipe RDBList — jebakan folder, E-4 §1) ·
`pyClassName` = `ASM-FW-GISFW-Int-policyjson`.

Pasangan nilai-sesudah / delta, beserta parameter pengisinya:

| Kolom "MENJADI" | Parameter | Kolom `*_SELISIH` | Parameter |
| --- | --- | --- | --- |
| `tsi_menjadi` | `Datain.CARI24` | `tsi_selisih` | `Datain.CARI26` |
| `premi_menjadi` | `Datain.CARI25` | `premi_selisih` | `Datain.CARI27` |
| `LOL_MENJADI` | `Datain.CARI50` | `LOL_SELISIH` | `Datain.CARI51` |
| `TSI100_MENJADI` | `Datain1.CARI48` | `TSI100_SELISIH` | `Datain1.CARI49` |
| `BROKERAGE_FEE_MENJADI` | `Datain1.CARI14` | `BROKERAGE_FEE_SELISIH` | `Datain1.CARI15` |
| `ricomm` *(tanpa sufiks)* | `Datain1.CARI10` | `ricomm_selisih` | `Datain1.CARI11` |
| `percent_ri_comm` *(ejaan beda)* | `Datain.CARI42` | `pct_ri_comm_selisih` | `Datain1.CARI9` |
| **`PCT_BROKERGARE_FEE`** *(ejaan asli)* | `Datain.CARI52` | `PCT_BROKERAGE_FEE_SELISIH` | `Datain.CARI53` |
| `prorate` | `Datain1.CARI8` | — | — |

⚠️ **5 kolom `*_MENJADI` vs 8 kolom `*_SELISIH`** — tiga pasangan memakai kolom "sesudah" **tanpa
sufiks** (E-2 §1.5, dikonfirmasi ulang). `PCT_BROKERGARE_FEE` **salah ejaan di skema**; skema tidak
berubah (`CLAUDE.md` §4.3) → **diport apa adanya**.

⛔ Seluruh kolom nilai dibungkus `To_number(Replace({…},',','.'))` — **koma desimal**, kontrak K-027.

### 3.2 Apakah ketujuh cabang mengisi parameter yang sama?

`[terverifikasi]` Cacah `<PropertiesName>` bertarget tiap parameter delta:

| Cabang | CARI24 | CARI25 | CARI26 | CARI27 | CARI30 | CARI42 | CARI50 | CARI51 | CARI52 | CARI53 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| Fire | 5 | 5 | **10** | **10** | 6 | 2 | 5 | 4 | 2 | 2 |
| Aneka | 3 | 3 | 7 | 7 | 3 | 2 | 3 | 3 | 2 | 2 |
| Golf | 3 | 3 | 5 | 5 | 3 | 2 | 3 | 2 | 2 | 2 |
| Marine Cargo | 3 | 3 | 5 | 5 | 3 | 2 | 3 | 2 | 2 | 2 |
| MBU | 3 | 3 | 5 | 5 | 3 | 2 | 2 | 2 | 2 | 2 |
| PA | 3 | 3 | 6 | 6 | 3 | 2 | 2 | 3 | 2 | 2 |
| **Bonding** | — | — | — | — | — | — | — | — | — | — |

✅ **Keenam cabang yang ada mengisi seluruh parameter delta.** Bedanya cacah penugasan (jumlah cabang
kondisi di dalamnya), bukan cakupan. FIRE punya dua kali lipat titik penugasan `CARI26`/`CARI27` —
konsisten dengan `Type`/`EdmType` bercabang lebih banyak di FIRE (E-3 §5.2).

### 3.3 Rumus delta baris spreading — dikonfirmasi ulang

`[terverifikasi]` `Endorsment Fac In\Activity\SaveFacinProdEDMFire_Act.xml`:

```
# nilai baru diprorata lebih dulu
SET Local.TsiSpreadEDM   = .TSISpreaded     * pyWorkPage.OfferFacIn.ProrateEDMEnd
SET Local.PremiSpreadEDM = .PremiumSpreaded * pyWorkPage.OfferFacIn.ProrateEDMEnd

# lalu dikurangi nilai lama dengan TreatyType yang sama
SET Local.TsiSpreadEDM   = Local.TsiSpreadEDM   - Local.TsiSpreadNB
SET Local.PremiSpreadEDM = Local.PremiSpreadEDM - Local.PremiSpreadNB
SET Datain1.CARI11       = Local.NewRIComm        - Local.OldRIComm
SET Datain1.CARI15       = Local.NewBrokerageFree - Local.OldBrokerageFree
```

✅ Rumus kanoniknya:

> **`SELISIH = (nilai_baru × ProrateEDMEnd) − nilai_lama_dengan_TreatyType_sama`**

`nilai_lama` diambil dari **`OfferFacIn.OldData`** (lapis A), **bukan** dari properti `*Old` per baris
(lapis B). E-3 §6.2 dikonfirmasi ulang — **57 penugasan delta berbentuk pengurangan** di berkas ini.

---

## 4. Penomoran versi polis — `PRODKE`

### 4.1 Mekanisme

`[terverifikasi]` `SaveEDMToJsonPolicy_Act` langkah 20 → 22 → 23:

```
[20] RDB-List GetProdKeOldData_SQL → page OldData      (pyStepsPreCondition = false)
[22] SET Local.Prodke      = OldData.pxResults(1).HASIL2
     SET InputData.CARI4   = Local.Prodke + 1               ← versi baru
     SET InputData.CARI5   = …QuotationData.OldPolicyNo
[23] RDB-List INSERTJSON_JSONPOLISEDM_FACIN               ← tulis baris baru
```

⛔ **Riwayat = baris bertambah**, bukan update in-place (`CLAUDE.md` §4.3). `PRODKE` **0-based**,
versi baru = versi terakhir + 1.

### 4.2 ⛔ Dua semantik "versi terakhir" yang berbeda di rantai yang sama

`[terverifikasi]`

| Rule | SQL | Semantik |
| --- | --- | --- |
| `GetProdKeOldData_SQL` *(dipakai langkah 20)* | `select PRODKE as HASIL2 from pooldata.json_polis where NOPOLIS={TempPolis.CARI4} **order by TGL_INPUT desc**` | baris **terbaru menurut waktu input** |
| `GetProdKeEDM_SQL` *(dipakai `GetEdmProdKe_Act`)* | `select PRODKE as HASIL2 from json_polis where NOPOLIS={InputData.CARI17} and **PRODKE=(SELECT COUNT(NOPOLIS)-1 …)**` | baris dengan **PRODKE = cacah − 1** |
| `GetEDMOldData_SQL` *(lapis A, E-3 §2)* | `… and **PRODKE=(SELECT COUNT(NOPOLIS)-1 …)**` | idem |

📌 `GetEdmProdKe_Act` sendiri hanya **1 langkah**: `RDB-List` → `GetProdKeEDM_SQL`.

⛔ **Keduanya memberi jawaban berbeda bila `PRODKE` berlubang.** Yang berbasis `COUNT-1` salah; yang
berbasis `TGL_INPUT desc` tetap benar. Jadi **arsip §10 #11 (kerapatan `PRODKE`) tetap MEMBLOKIR** —
dan kini diketahui **tidak seluruh jalur** rentan terhadapnya. → **E-Q30**.

⚠️ `[pertanyaan terbuka]` **tidak ditebak:** adakah constraint yang menjamin `PRODKE` rapat dari 0
tanpa lubang. Milik **DBA**.

### 4.3 JSON_POLIS ditulis oleh stored procedure

`[terverifikasi]` `RDBList\INSERTJSON_JSONPOLISEDM_FACIN.xml` (`Rule-Connect-SQL`):

```sql
BEGIN POOLDATA.INSERTJSONPOLIS( {pyWorkPage.pzInsKey},
                                {pyWorkPage.OfferFacIn.PolicyData.PolicyNo},
                                {pyWorkPage.OfferFacIn.PolicyData.EndorsementNo}, NULL, … ); END;
```

⛔ **Penulisan JSON_POLIS bukan `INSERT` biasa — ia panggilan prosedur `POOLDATA.INSERTJSONPOLIS`
yang isinya tidak ada di korpus.** Arsip §10 #8 **MEMBLOKIR**, kini terletak persis. Lihat §6.

---

## 5. Jalur fac out EDM

### 5.1 Rantai

`[terverifikasi]` Titik masuk **`Flow\OfferFacRetro`** — salah satu dari hanya 2 Flow di EDM:

```
Flow\OfferFacRetro
  └─ Activity\InsertFacoutProd
       └─ Activity\InsertFacoutProductionEDM            (pxObjClass = Rule-Obj-Activity)
            ├─ [1] Call InsertFacoutProductionEDMCurr   IF QuotationData.Type=="7"
            │        « for adjusemnt curenncy masukin yg old datanya »
            ├─ [4] RDB-List GetNopolis1_Sql             IF @contains(pyWorkIDPrefix,"EDM-")
            ├─ [8] RDB-List CekFacoutProd_Sql           « Cek Ditable sudah ada atau belum »
            ├─ [9] « Exit Jika sudah ada di table »     IF QuotationData.Type=="7"
            └─ [10]–[15] per lini → RequestType = InsertTreatyProd_Sql
                 IsMarineCargo · IsAneka · IsFire · IsMBU · isGolfInsurance · IsPA
```

⚠️ **`InsertFacoutProductionEDM` dan `…Curr` adalah `Rule-Obj-Activity`**, bukan rule SQL — namanya
menyesatkan. Penulisannya lewat **`InsertTreatyProd_Sql`**, yang **berbeda** dari
`InsertTreatyProduction_Sql` (jalur fac in). Dua nama nyaris sama, dua rule berbeda, dua tabel
berbeda.

### 5.2 Tabel `FACOUTPRODUCTION`

`[terverifikasi]` `RDBList\InsertTreatyProd_Sql.xml` (`Rule-Connect-SQL`) → **`FACOUTPRODUCTION`**,
**65 kolom**:

| | Jumlah | Kolom |
| --- | ---: | --- |
| `*_MENJADI` | **2** | `PREMI_COVERAGE_MENJADI` · `COMMISION_COVERAGE_MENJADI` |
| `*_SELISIH` | **5** | `SHAREOFFERED_SELISIH` · `OBJECTPREMI_SELISIH` · `COMMISION_SELISIH` · `PREMI_COVERAGE_SELISIH` · `COMMISION_COVERAGE_SELISIH` |

📌 **Ketidaksetangkupan yang sama seperti `facinproduction`** (5 vs 8 di sana, 2 vs 5 di sini): tiga
`*_SELISIH` tanpa pasangan `*_MENJADI`. Polanya konsisten antar-tabel — nilai "sesudah" sebagian
disimpan di kolom tanpa sufiks.

Gerbang idempotensi: `CekFacoutProd_Sql` →
`select DISTINCT IDPEGA from FACoutPRODUCTION where RISLIPNO={DataIN.CARI7} and IDPEGA={pyWorkPage.pzInsKey}`

⚠️ **Tetapi langkah "Exit Jika sudah ada di table" bergerbang `Type=="7"` saja.** Untuk `Type` lain,
hasil `CekFacoutProd_Sql` **tidak dipakai untuk keluar**. `[pertanyaan terbuka]` apakah disengaja →
**E-Q32**.

📌 Kaitan `IsFacRetro` / `EndorsementInternalRetro`: `IsFacRetro` menggerbangi langkah **7**
`serviceInsertArasapasEDM_act` (§6). Gerbang `EndorsementInternalRetro==1` ada di
`SetErrorBatalEndorsement_Act` (arsip §1.3). Keduanya **belum ditelusuri tuntas** — wilayah arsip §6,
bukan E-5.

---

## 6. `serviceInsertArasapasEDM_act` — konversi ke produksi

`[terverifikasi]` 21 langkah:

| # | Isi |
| ---: | --- |
| 1–4 | `Obj-Refresh-And-Lock` · `Property-Set` · `RDB-List GetPolicyNoByCaseId` · set bila `PolicyNo==""` |
| **5** | `Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService` — *"GET LINK SERVICE"* |
| **6** | `Connect-REST` — *"Hit service Arasapas 1 - FACIN"* · `IF IsPEGAPROD` |
| **7** | `Connect-REST` — *"Hit service Arasapas 2 - FACOUT"* · `IF IsFacRetro` |
| 8–9 | set param log · `Call InsertLogServiceProd` · `IF IsPEGAPROD` |
| 10–12 | *"SET ERROR GAGAL KONVERSI"* · `Obj-Save` · `Page-Set-Messages` · `IF IsSuccessHitService` |
| 13 | `RDB-List` `CekSTSKonversiJson` |
| **14** | `RDB-List` — *"DELETE PRODUKSI JIKA ERROR (PROCEDURE)"* · `RequestType` = **`DeleteDataProduction`** · `IF IsSuccessHitService` |
| 15–16 | siapkan data `JSON_POLIS_MONITORING` |
| **17** | `RDB-List` `INSERTJSON_JSONPOLISMONITORING_FACIN` · `IF pyWorkIDPrefix=="EDMT-"` |
| **18** | `RDB-List` `UpdateErrorNoteJsonPolisMonitoring` · `IF pyWorkIDPrefix=="EDMT-"` |
| 19–21 | `Page-Remove` · set param email · `Call SendEmailWithAttachments` · `IF IsSuccessHitService` |

⛔ **Endpoint diambil dari tabel `M_LINK_SERVICE`** (langkah 5) — isinya **tidak ada di korpus**
(`CLAUDE.md` §4.4). Endpoint wajib jadi konfigurasi, tidak pernah literal.

⚠️ **Langkah 17–18 bergerbang `pyWorkIDPrefix=="EDMT-"`** — prefiks yang tidak muncul di tempat lain
di korpus Fac In (arsip §10 #23, `[pertanyaan terbuka]`, kemungkinan endorsement Treaty). Artinya
untuk case `EDM-` biasa, `JSON_POLIS_MONITORING` **tidak ditulis**.

### 6.1 Alamat email literal — dihitung, tidak disalin

`[terverifikasi]` Di `serviceInsertArasapasEDM_act.xml`:

| Ukuran | Jumlah |
| --- | ---: |
| Kemunculan pola alamat email | **6** |
| **Alamat unik** | **3** |
| Domain unik | **2** |

✅ **Mengonfirmasi arsip §10 #19 ("3 alamat email tertanam literal").** Nilainya **tidak disalin**
(K-025 / E-Q15). Harus jadi konfigurasi di sistem baru (`CLAUDE.md` §4.4).

---

## 7. ⛔ Blocker — bergantung pada yang tidak ada di korpus

`[terverifikasi]` Stored procedure yang **dipanggil** rantai produksi EDM, **isinya tidak ada di
korpus**:

| Rule pemanggil | Panggilan | Menentukan |
| --- | --- | --- |
| `RDBList\INSERTJSON_JSONPOLISEDM_FACIN` | `BEGIN POOLDATA.INSERTJSONPOLIS( pzInsKey, PolicyNo, EndorsementNo, NULL, … ); END;` | **penulisan JSON_POLIS + penomoran versi** |
| `RDBList\DeleteDataProduction` | `BEGIN POOLDATA.PEGA_DELETE_ERROR_KONVERSI( pzInsKey, Quotation.BusinessFac, HASIL10 out ); END;` | **rollback produksi saat gagal konversi** |
| `RDBList\GetSequenceNumber_SQL` | `BEGIN POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER( CARI1, CARI2, TO_DATE(CARI3,'DD/MM/YYYY'), HASIL1 out, HASIL2 out ); END;` | **penomoran endorsement** |
| `RDBList\MachingDataFacin_Sql` | blok `DECLARE … ` PL/SQL anonim | pencocokan data saat galat |

Tabel yang disentuh tetapi isinya tidak ada di korpus:
`pooldata.treatyproduction_backup` (via `InsertFacinProductionBackup_Sql`) ·
`POOLDATA.ERRORFACINPROD` (via `INSERTERRORFACIN_Sql`) · `JSON_POLIS_MONITORING` ·
`M_LINK_SERVICE` · `OPENPROTEKSI_EDM`.

⛔ **Arsip §10 #8 tetap MEMBLOKIR** — kini dengan lokasi panggilan yang persis. Milik **DBA**.
Menunggu tabel flat + `ALL_SOURCE`.

### 7.1 ⛔ Kualifikasi skema tidak konsisten

`[terverifikasi]` Tabel yang **sama** dirujuk dua cara berbeda dalam satu rantai:

| Rule | Tabel seperti tertulis |
| --- | --- |
| `InsertTreatyProduction_Sql` | **`facinproduction`** *(tanpa skema)* |
| `CekFacin_Sql` | **`pooldata.facinproduction`** |
| `GetProdKeEDM_SQL` | **`json_polis`** *(tanpa skema)* |
| `GetProdKeOldData_SQL` · `CekSTSKonversiJson` | **`pooldata.json_polis`** |
| `GetEDMOldData_SQL` | **`JSON_POLIS`** *(tanpa skema)* |
| `InsertTreatyProd_Sql` · `CekFacoutProd_Sql` | `FACOUTPRODUCTION` / `FACoutPRODUCTION` *(tanpa skema)* |

⚠️ Di Oracle, nama tanpa skema diselesaikan lewat skema sesi atau sinonim. **Gerbang idempotensi
membaca `pooldata.facinproduction` sementara INSERT-nya menulis ke `facinproduction`** — bila
keduanya tidak menunjuk objek yang sama, guard idempotensi tidak berfungsi. `[pertanyaan terbuka]`
→ **E-Q29**, milik **DBA**.

---

## 8. Kejanggalan yang diport apa adanya

Sesuai `CLAUDE.md` §1 — kandidat perbaikan milik bisnis, **tidak diperbaiki di sini**.

| Kejanggalan | Bukti | Sikap |
| --- | --- | --- |
| **`GenerateEndorsementNo` & `GenerateEDMNoLife`** dirujuk langkah 12–13 sebagai `RequestType`, berkasnya **NIHIL di seluruh `RNM\`** | E-2 §2 · dikonfirmasi ulang di §1.1 | **`panic`** (`CLAUDE.md` §4.5) |
| **`SaveFacinProdEDMBonding_Act`** dipanggil tapi tidak ada di korpus — **ada di `DDL\`** | §2.2 | **bukan `panic`** — kandidat amandemen K-006 · **E-Q28** |
| **`PCT_BROKERGARE_FEE`** salah ejaan di skema (`BROKERGARE`, bukan `BROKERAGE`) | §3.1 | **diport apa adanya** — skema tidak berubah (§4.3) |
| **Langkah 8, 19, 20** `SaveEDMToJsonPolicy_Act` ber-`pyStepsPreCondition = false` — termasuk langkah 20 yang **mengambil `PRODKE`** | §1.1 | tergantung **E-Q16** (semantik `false`) |
| **Dua semantik `PRODKE`** (`COUNT-1` vs `TGL_INPUT desc`) dalam satu rantai | §4.2 | **diport apa adanya** · E-Q30 |
| **Ketidaksetangkupan `*_MENJADI`/`*_SELISIH`** di kedua tabel produksi | §3.1, §5.2 | **diport apa adanya** |
| Gerbang idempotensi fac out hanya aktif untuk `Type=="7"` | §5.2 | **diport apa adanya** · E-Q32 |
| `JSON_POLIS_MONITORING` hanya ditulis untuk prefiks `EDMT-` | §6 | **diport apa adanya** · arsip #23 |

---

## 9. Pertanyaan terbuka baru dari E-5

| # | Pertanyaan | Pemilik | Dampak |
| ---: | --- | --- | --- |
| **E-Q28** | **`SaveFacinProdEDMBonding_Act` ada di `DDL\` tetapi tidak di korpus.** Dilipat masuk sebagai amandemen (pola K-006 B.1)? Isinya belum dibaca — di luar batas kerja putaran ini | **work owner** | Menentukan apakah lini Bonding punya jalur produksi |
| **E-Q29** | **Kualifikasi skema tidak konsisten** — `facinproduction` vs `pooldata.facinproduction` dalam satu rantai; guard idempotensi dan INSERT bisa menunjuk objek berbeda | **DBA** | Bila berbeda objek, guard idempotensi tidak berfungsi |
| **E-Q30** | **Dua semantik "versi terakhir `PRODKE`"** — `COUNT-1` (rentan lubang) vs `order by TGL_INPUT desc` (tidak rentan). Mana yang benar? | **DBA + work owner** | Menentukan before-image dan penomoran versi mana yang sah |
| **E-Q31** | **Tidak ada cabang Life maupun Travel** di `SaveFacinProdAllEDM_Act`. Life lewat jalur sendiri; **Travel tidak punya jalur produksi sama sekali** meski `SetOldData` menanganinya | work owner | Apakah endorsement Travel memang tidak masuk produksi |
| **E-Q32** | **Gerbang idempotensi fac out** (`CekFacoutProd_Sql`) hasilnya hanya dipakai keluar bila `Type=="7"`. Untuk `Type` lain, penulisan ganda tidak dicegah | work owner | Risiko baris ganda di `FACOUTPRODUCTION` |

Yang **tetap** dan belum terjawab: **E-Q13** · **E-Q15** (dikonfirmasi ulang di §6.1) · **E-Q16**
(menyentuh langkah 8/19/20) · **E-Q17** · **E-Q18** · **E-Q19** · **E-Q21** · **E-Q22** ·
**E-Q24** · **E-Q26** · **E-Q27**.

---

## 10. Belum dikerjakan (bukan untuk prompt ini)

- **E-6** — celah, rujukan menggantung, penutup. **Tidak dimulai tanpa perintah.**
- Isi `SaveFacinProdEDMBonding_Act` di `DDL\` — menunggu E-Q28
- Enam varian `SaveFacinProdCurr*_Act` (jalur Adjustment Currency, `Type=="7"`) — belum ditelusuri
- `SaveFacinLive_Act` dan `SaveFacinSpreadLife_Sql` (jalur produksi Life) — belum ditelusuri
- `InsertCedingProduction` (langkah 24) — belum ditelusuri
- Arsip §4 (gerbang akseptasi EDM) termasuk klaim "nilai dasar akseptasi = SELISIH TSI" — **belum diuji**
- Jalur produksi NB sebagai pembanding — **menunggu DBA** (tabel flat + `ALL_SOURCE`)

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
