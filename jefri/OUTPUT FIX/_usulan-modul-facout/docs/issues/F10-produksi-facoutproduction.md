# F10: Produksi Fac Out ke tabel `FACOUTPRODUCTION`

**What to build:** Penawaran retrosesi yang sudah disetujui tersimpan sebagai baris produksi Fac Out,
tertaut ke objek Fac In induknya, dan **aman diulang** tanpa menghasilkan baris ganda.

✅ **Tabel produksinya sudah ada dan sudah flat** — 65 kolom, multi-lini, dengan pasangan
`_MENJADI`/`_SELISIH`. **Tidak perlu tabel flat baru** (K-056). Ini membedakan F10 dari `edm\E21`, yang
masih menunggu rancangan tabel flat.

⛔ **Insert = Connect-SQL langsung, BUKAN stored procedure.** `[terverifikasi]`
`DDL\InsertTreatyProd_Sql.xml`, `<pxObjClass>` = `Rule-Connect-SQL`, tag `<pyBrowseSQL>`:
`BEGIN INSERT INTO FACOUTPRODUCTION (…65 kolom…) VALUES (…); COMMIT; END;`. Catatan lama "via SP
`INSERTFACOUTPRODUCTION`" sudah **dibatalkan** — itu nama request/connector, bukan SP basis data.

## Gerbang idempotensi

`[terverifikasi]` `Activity\InsertFacoutProduction` memanggil `RDBList\CekFacoutProd_Sql`:

```sql
select DISTINCT IDPEGA from FACoutPRODUCTION
 where RISLIPNO = {DataIN.CARI7} and IDPEGA = {pyWorkPage.pzInsKey}
```

Gerbangnya `@SizeOfPropertyList(FacoutList.pxResults)>0` dengan
`pyStepsPreCondParamsWhenTrue` = **6** / `pyStepsPreCondParamsWhenFalse` = **2**.

`[terverifikasi]` Gerbang itu berada di alamat **`RH_1.pySteps(7)`**, dan
`<pyStepsPreCondParamsWhen>`, `…WhenTrue`, `…WhenFalse` berada di **`rowdata` yang sama** — keanggotaan
struktural, bukan kedekatan teks.

`[dugaan kuat]` Penambatan kode transisi dari dua jangkar independen (`SaveEDMToJsonPolicy_Act`
`RH_1.pySteps(10)` T=6, dan `SaveTreatyProduction_Act` `RH_1.pySteps(4)` T=2/F=6): **`2` = jalankan
lalu lanjut · `3` = lewati langkah · `6` = keluar activity**. Maka: **baris sudah ada → keluar, tidak
menyisipkan; belum ada → lanjut dan sisipkan.**

⚠️ Ditandai **`[dugaan kuat]`, bukan `[terverifikasi]`**. Yang sudah terbukti adalah **keanggotaan
struktural** gerbang dan T/F pada langkah yang sama (ditunjukkan lewat `<pyStepPageReference>` untuk
kelima jangkar) serta pola distribusinya. Yang **belum** terbukti adalah **arti** kodenya — korpus
tidak memuat berkas yang memetakan kode→perilaku secara eksplisit. Kenaikan label memerlukan
konfirmasi UI Pega work owner.

Kode **1**, **4** dan **5** **belum tertambat** — `[pertanyaan terbuka]`, butuh UI Pega work owner.
Kode **5** yang paling sering di antaranya (908 kemunculan pada satu pasangan).

📌 **Sensus lengkap 31 pasangan (T,F), bukti keanggotaan kelima jangkar, dan perintah auditnya ada di
`..\..\10-audit\02-peta-kode-transisi.md`.**

## Pemetaan 65 kolom — disalin dari `09-facout\01` §4.3.1

`[terverifikasi]` **65 kolom seimbang dengan 65 nilai**, parameter dari **lima halaman**:
`DataIN` 35 · `DataIN1` 14 · `DataINCargo` 9 · `DataINMBU` 5 · `DataINPA` 2.

| # | Kolom | Tipe Oracle | Parameter |
| ---: | --- | --- | --- |
| 1 | `IDPEGA` | `VARCHAR2(150)` | `{DataIN.CARI1}` |
| 2 | `POLICYNO` | `VARCHAR2(50)` | `{DataIN.CARI2}` |
| 3 | `GROUPPANEL` | `VARCHAR2(10)` | `{DataIN.CARI3}` |
| 4 | `REINSURER_ID` | `VARCHAR2(15)` | `{DataIN.CARI4}` |
| 5 | `REINSURER_NAME` | `VARCHAR2(150)` | `{DataIN.CARI5}` |
| 6 | `TGL_PRINT` | `DATE` | `To_date({DataIN.CARI6}, 'DD/MM/YYYY HH24:MI:SS')` |
| 7 | `RISTARTPERIOD` | `DATE` | `To_date({DataIN.CARI8}, …)` |
| 8 | `RIENDPERIOD` | `DATE` | `To_date({DataIN.CARI9}, …)` |
| 9 | `START_DATE` | `DATE` | `To_date({DataIN.CARI11}, …)` |
| 10 | `END_DATE` | `DATE` | `To_date({DataIN.CARI12}, …)` |
| 11 | `RISLIPNO` | `VARCHAR2(600)` | `{DataIN.CARI7}` |
| 12 | `NOENDORS` | `VARCHAR2(100)` | `{DataIN.CARI13}` |
| 13 | `PACKINGID` | `VARCHAR2(200)` | `{DataINCargo.CARI1}` |
| 14 | `PACKINGNOTE` | `VARCHAR2(900)` | `{DataINCargo.CARI2}` |
| 15 | `GOODNOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI3}` |
| 16 | `TRADINGNOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI4}` |
| 17 | `SHIPID` | `VARCHAR2(100)` | `{DataINCargo.CARI6}` |
| 18 | `FROMRUTE` | `VARCHAR2(500)` | `{DataINCargo.CARI7}` |
| 19 | `TORUTE` | `VARCHAR2(500)` | `{DataINCargo.CARI8}` |
| 20 | `SAILDATE` | `DATE` | `To_date({DataINCargo.CARI9}, …)` |
| 21 | `CONVEYANCENOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI5}` |
| 22 | `OBJECTNO` | `VARCHAR2(50)` | `{DataIN1.CARI1}` |
| 23 | `ZIPCODE` | `VARCHAR2(50)` | `{DataIN1.CARI2}` |
| 24 | `PROVINCE` | `VARCHAR2(100)` | `{DataIN1.CARI3}` |
| 25 | `CITY` | `VARCHAR2(100)` | `{DataIN1.CARI4}` |
| 26 | `DISTRICT` | `VARCHAR2(100)` | `{DataIN1.CARI5}` |
| 27 | `RW` | `VARCHAR2(1000)` | `{DataIN1.CARI6}` |
| 28 | `ADDRESS` | `VARCHAR2(4000)` | `{DataIN1.CARI7}` |
| 29 | `BUILDINGNO` | `VARCHAR2(30)` | `{DataIN1.CARI8}` |
| 30 | `ROADNAME` | `VARCHAR2(4000)` | `{DataIN1.CARI9}` |
| 31 | `OBJECTNAME` | `VARCHAR2(4000 CHAR)` | `{DataIN1.CARI10}` |
| 32 | `OCCUPATION` | `VARCHAR2(4000)` | `{DataIN1.CARI11}` |
| 33 | `OBJECTITEM` | `VARCHAR2(4000)` | `{DataIN1.CARI12}` |
| 34 | `OBJECTITEMID` | `VARCHAR2(100)` | `{DataIN1.CARI13}` |
| 35 | `BRANDNAME` | `VARCHAR2(500 CHAR)` | `{DataINMBU.CARI2}` |
| 36 | `LICENSEPLATE` | `VARCHAR2(50)` | `{DataINMBU.CARI3}` |
| 37 | `TYPENAME` | `VARCHAR2(50)` | `{DataINMBU.CARI4}` |
| 38 | `MODELNAME` | `VARCHAR2(500)` | `{DataINMBU.CARI5}` |
| 39 | `ENGINENUMBER` | `VARCHAR2(100)` | `{DataINMBU.CARI6}` |
| 40 | `CLASSPA` | `VARCHAR2(100)` | `{DataINPA.CARI1}` |
| 41 | `DOB` | `VARCHAR2(100)` | `{DataINPA.CARI2}` |
| 42 | `CURRENCY` | `VARCHAR2(10)` | `{DataIN.CARI23}` |
| 43 | `CURRENCYID` | `VARCHAR2(10)` | `{DataIN.CARI29}` |
| 44 | `TSIRNM` | `NUMBER(20,4)` | `{DataIN.CARI30}` |
| 45 | `TSISPREADED` | `NUMBER(20,4)` | `{DataIN.CARI31}` |
| 46 | `PCTOFFERED` | `NUMBER(20,4)` | `{DataIN.CARI32}` |
| 47 | `SHAREOFFERED` | `NUMBER(20,4)` | `{DataIN.CARI24}` |
| 48 | `OBJECTPREMI` | `NUMBER(20,4)` | `{DataIN.CARI25}` |
| 49 | `RICOMM` | `NUMBER(20,4)` | `{DataIN.CARI26}` |
| 50 | `COMMISION` | `NUMBER(20,4)` | `{DataIN.CARI27}` |
| 51 | `RATE` | `NUMBER(25,20)` | `{DataIN.CARI28}` |
| 52 | `SHAREOFFERED_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI33}` |
| 53 | `OBJECTPREMI_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI34}` |
| 54 | `COMMISION_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI35}` |
| 55 | `TYPEFACULTATIVE` | `VARCHAR2(50)` | `{DataIN.CARI39}` |
| 56 | `PRORATE` | `NUMBER(25,20)` | `{DataIN.CARI40}` |
| 57 | `RATE_COVERAGE` | `NUMBER(20,4)` | `{DataIN.CARI36}` |
| 58 | `PREMI_COVERAGE_MENJADI` | `NUMBER(20,8)` | `{DataIN.CARI37}` |
| 59 | `PREMI_COVERAGE_SELISIH` | `NUMBER(20,8)` | `{DataIN.CARI38}` |
| 60 | `COVERAGE_NAME` | `VARCHAR2(1000)` | `{DataIN.CARI42}` |
| 61 | `COVERAGE_ID` | `VARCHAR2(100)` | `{DataIN.CARI41}` |
| 62 | `COMMISION_COVERAGE_PCT` | `NUMBER(20,8)` | `{DataIN.CARI43}` |
| 63 | `COMMISION_COVERAGE_MENJADI` | `NUMBER(20,8)` | `{DataIN.CARI44}` |
| 64 | `COMMISION_COVERAGE_SELISIH` | `NUMBER(20,8)` | `{DataIN.CARI45}` |
| 65 | `OBJECTNO_FACIN` | `VARCHAR2(10)` | `{DataIN1.CARI14}` |

⛔ **Pemetaannya TIDAK berurutan** — lihat kolom 6–11, 17–21, dan 42–57. **Diport apa adanya; jangan
"dirapikan".** Satu geseran kolom = korupsi data produksi yang senyap.

`[terverifikasi]` **Sebelas parameter tidak terpakai:** `DataIN.CARI10` dan `CARI14`…`CARI22`
(sepuluh), plus `DataINMBU.CARI1`.

⛔ **Interpolasi string diganti parameter terikat.** `{Halaman.CARIn}` di sistem lama adalah
**substitusi teks** ke dalam SQL. Di sistem baru ia **wajib** menjadi parameter terikat. Itu perubahan
mekanisme penyampaian nilai, **bukan** perubahan perilaku — nilai yang tertulis ke kolom tetap sama,
sehingga tidak melanggar `CLAUDE.md` §1.

⚠️ **Kolom ber-PII:** `DOB`, `LICENSEPLATE`, `ENGINENUMBER`, `ADDRESS`, `OBJECTNAME`,
`REINSURER_NAME`. Nilainya **tidak pernah** disalin ke test, fixture, log, atau dokumen (K-025).

**Asal (Pega).** `DDL\InsertTreatyProd_Sql.xml` · `DDL\FACOUTPRODUCTION.txt` ·
`Activity\InsertFacoutProd` (pengirim) · `InsertFacoutProduction` · `RDBList\CekFacoutProd_Sql` ·
`Activity\UpdateStsKonversiFacOut_Act`

**Keputusan.** **K-056** · K-053 · K-027 · K-010/K-012 · `CLAUDE.md` §4.1, §4.3, §4.4, §4.6

**Blocked by:** F06 · F08 · F09

**Status:** blocked

- [ ] Baris produksi tertulis dengan **65 kolom** sesuai tabel di atas, **urutan pemetaan apa adanya**
- [ ] ⛔ Pemetaan diambil dari **tabel**, bukan dari urutan kolom — tidak ada penyesuaian "supaya rapi"
- [ ] Sebelas parameter tak terpakai **tetap tak terpakai**
- [ ] **Parameter terikat**, bukan interpolasi string; alasannya dicatat dalam komentar
- [ ] Gerbang idempotensi: baris sudah ada → **tidak menyisipkan**; rantai **aman diulang**
- [ ] Penaut `OBJECTNO_FACIN` terisi dari nilai yang lahir di **F02**, bukan dihitung ulang
- [ ] Uang bertipe `Money` dengan mata uang; ditulis ke `NUMBER` Oracle — **tidak pernah `float`**
- [ ] Skema Oracle **tidak berubah** (§4.3)
- [ ] ⛔ Nilai `DOB`/`LICENSEPLATE`/`ENGINENUMBER`/`ADDRESS` **tidak pernah** masuk test atau log
- [ ] `[pertanyaan terbuka]` kode transisi **1**, **4**, **5** belum tertambat — bila tercapai, **`panic`**
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
