# Katalog tabel warisan yang dipakai tiket 02, 03, dan 06 — instance pengembangan

`[data DBA — dibaca sendiri 26 September 2026 malam, dilengkapi malam yang sama (RATE_LIFE, PRODUCT_LIFE, agregat treaty), belum dikonfirmasi DBA]`

Sumber: `ALL_TAB_COLUMNS`, `ALL_IND_COLUMNS`, `ALL_INDEXES`, dan **agregat** (`COUNT`, `GROUP BY`)
pada instance **pengembangan**, skema `POOLDATA`. Nol baris data dibaca. Produksi belum dibaca.

## `M_LIFE_PREMIUM_DETAIL` — peserta polis (tiket 02 pencarian peserta, 03, 06)

⚠️ **66.795.704 baris**, 85 kolom. Pencarian peserta tanpa index dan tanpa batas hasil akan
memindai seluruh tabel. Index yang ada dan berguna untuk pencarian:

| Index | Kolom |
| --- | --- |
| `..._INDEX4` | `PL_NUMBER` |
| `..._INDEX5` | `PL_NUMBER`, `CURRENCY` |
| `..._INDEX8` | `CERTIFICATE_NO` |
| `..._INDEX1` | `POLICY_NO` |
| `..._INDEX22` | `NAME_OF_INSURED`, `PL_NUMBER` |
| `..._INDEX3` | `ID`, `CURRENCY` |
| `..._INDEX21` | `PL_NUMBER_EDM` |

Tidak ada index yang berawalan `EDMSTATUS`. Penyaring `EDMSTATUS` karena itu harus datang
**sesudah** penyaring `PL_NUMBER`, bukan sendirian.

### `EDMSTATUS` — sebaran nilai (agregat)

| Nilai | Cacah |
| --- | ---: |
| `NULL` | 59.083.299 |
| `Batal` | 7.711.454 |
| `Old` | 661 |
| `New` | 148 |
| `Delete` | 142 |
| `''` (teks kosong) | **0** |

⭐ Menjawab **OQ-001 sisa** untuk instance pengembangan: baris new business ber-`EDMSTATUS`
**`NULL`**, bukan teks kosong. Penyaring yang benar tetap menangani keduanya
(`EDMSTATUS IS NULL OR EDMSTATUS NOT IN ('Batal','Delete')`) sampai DBA memastikan produksi.

### Kolom (nomor | nama | tipe | panjang)

```
 1 POLICY_NO VARCHAR2 255        30 WPC DATE                          59 OVR_COMM_RETRO NUMBER
 2 POLICY_HOLDER VARCHAR2 255    31 GROSS_VALUATION_BEGIN_DATE DATE   60 BROKERAGE_FEE_RETRO NUMBER
 3 CERTIFICATE_NO VARCHAR2 255   32 GROSS_VALUATION_EXPIRED_DATE DATE 61 NET_PREMIUM_RETRO NUMBER
 4 NAME_OF_INSURED VARCHAR2 255  33 RETRO_VALUATION_BEGIN_DATE DATE   62 GROSS_PREMIUM_REFUND_RETRO NUMBER
 5 SEX VARCHAR2 255              34 RETRO_VALUATION_EXPIRED_DATE DATE 63 DISCOUNT_PREMIUM_REFUND_RETRO NUMBER
 6 AGE NUMBER                    35 RI_ADMIN_FEE NUMBER               64 OVR_COMM_REFUND_RETRO NUMBER
 7 PLAN VARCHAR2 255             36 SHARE_NUSANTARA_RE NUMBER         65 BROKERAGE_FEE_REFUND_RETRO NUMBER
 8 PERIOD_YY NUMBER              37 GROSS_PREMIUM NUMBER              66 NET_PREMIUM_REFUND_RETRO NUMBER
 9 PERIOD_MM NUMBER              38 RI_COMM NUMBER                    67 CLAIM_AMOUNT NUMBER
10 DESCRIPTION VARCHAR2 1000     39 NET_PREMIUM NUMBER                68 FLEET_DISCOUNT NUMBER
11 ID VARCHAR2 50                40 CLAIM NUMBER                      69 RATE NUMBER
12 CURRENCY VARCHAR2 255         41 TAX NUMBER                        70 SUM_AT_RISK_GROSS NUMBER
13 MEDICAL_STATUS VARCHAR2 255   42 OVR_COMM NUMBER                   71 SUM_AT_RISK_RETRO NUMBER
14 PL_NUMBER VARCHAR2 255        43 BROKERAGE_FEE NUMBER              72 RETROCEDED_SHARE NUMBER
15 PASSED_PERIOD VARCHAR2 100    44 EM_PERCENT NUMBER                 73 SHARE_NUSANTARA_RE_GROSS NUMBER
16 COUNT NUMBER(*,0)             45 SUM_REASURED NUMBER               74 DEDUCTION NUMBER
17 CEDING_CO VARCHAR2 999        46 SUM_INSURED NUMBER                75 FACTOR NUMBER
18 TGL_INPUT DATE                47 CEDING_RETENTION NUMBER           76 DEDUCTION_REFUND NUMBER
19 PRORATETYPE VARCHAR2 999      48 PROF_COMM NUMBER                  77 RI_ADMIN_FEE_REFUND_RETRO NUMBER
20 ENTRY_AGE NUMBER              49 COMM NUMBER                       78 RI_ADMIN_FEE_RETRO NUMBER
21 CURRENT_AGE NUMBER            50 GROSS_PREMIUM_REFUND NUMBER       79 RI_ADMIN_FEE_REFUND NUMBER
22 EDMSTATUS VARCHAR2 25         51 NET_PREMIUM_REFUND NUMBER         80 IDPEGA VARCHAR2 50
23 DOB DATE                      52 COMM_REFUND NUMBER                81 KTP NUMBER
24 BEGIN_DATE DATE               53 OVR_COMM_REFUND NUMBER            82 STATUS VARCHAR2 10
25 EXPIRED_DATE DATE             54 BROKERAGE_FEE_REFUND NUMBER       83 STATUSOLD VARCHAR2 10
26 START_DATE DATE               55 TAX_REFUND NUMBER                 84 PL_NUMBER_EDM VARCHAR2 100
27 EFFECTIVE_DATE DATE           56 SHARE_RETRO NUMBER                85 RISK NUMBER
28 LAPSE_DATE DATE               57 GROSS_PREMIUM_RETRO NUMBER
29 STNC DATE                     58 DISCOUNT_PREMIUM_RETRO NUMBER
```

Seluruh kolom nullable. Yang dipakai tiket 06 ada di sini, bukan di tabel klaim: keempat tanggal
valuasi (`GROSS_VALUATION_*`, `RETRO_VALUATION_*`) dan `WPC`, semuanya `DATE`. ⚠️ `KTP` bertipe
`NUMBER` — nomor identitas yang disimpan sebagai angka menghilangkan nol depan; **jangan disalin ke
skema baru dalam bentuk itu** (ADR-U-0022), dan jangan pernah masuk fixture.

## Tabel treaty untuk spreading (tiket 03)

`SpreadingClaimLife_Act` merujuk empat rule: `GetRetroLife_SQL` → `RETROCESSIONLIFE`,
`GetRateRetro` → `POOLDATA.RATE_LIFE`, `GetProductLife` → view `PRODUCT_LIFE`, dan
`GetJsonProductLife` → `m_product_life.JSONDATA` ⛔ *(dilarang AC 38; jangan ditiru)*.

> **Ralat 26 September 2026 malam** `[terverifikasi — baris pertama tag `pyBrowseSQL` tiap berkas]`:
> `Claim Life/RDBList/GetJsonProductLife.xml` dimulai `select * from treatyyear_life`, sedangkan
> `Claim Life/RDBList/GetProductLife.xml` dimulai `SELECT M_PRODUCT_LIFE.JSONDATA AS CARI1,
> PRODUCT_LIFE.RICOMM FROM PRODUCT_LIFE …`. Pembaca JSON produk di Pega adalah **`GetProductLife`**;
> baris kedua dan seterusnya kedua SQL itu **belum dibaca** — executor tiket 03 membacanya utuh.
> SQL `GetRateRetro` `[terverifikasi]`: `SELECT AGE AS CARI1, CONTRACT AS CARI2, GENDER AS CARI3,
> RATE AS CARI4 FROM POOLDATA.RATE_LIFE WHERE IDUSEDBY = {InputData.CARI3}` *(tanpa `ORDER BY`)*;
> `GetRetroLife_SQL`: `select * from retrocessionlife where idtreatyyear_life ={ParamData.CARI2}
> order by id asc`. Larangan AC 38 berlaku pada **pembacaan JSON produk**, apa pun nama rule-nya.

### `RETROCESSIONLIFE` — **VIEW**, 13 kolom, **seluruhnya `VARCHAR2(4000)`**

`ID` (7), `IDTREATYYEAR_LIFE`, `PERCENTSHARE`, `RATE`, `COMMISION`, `OVR_COMM`, `TREATYTYPEID`,
`TREATYTYPENAME`, `TREATYSTARTDATE`, `TREATYENDDATE`, `USERID`, `TGLUPDATE`, `REINSURERNAME`.

⚠️ Angka spreading (`PERCENTSHARE`, `RATE`, `COMMISION`, `OVR_COMM`) tiba sebagai **teks** dari
view ini. Pembacanya wajib mengurai lewat `utils.ParseDecimal` dan **melaporkan** yang tidak
terurai, bukan menebak; pemisah desimal yang dipakai view belum diketahui — periksa agregat dulu.

### `TREATYYEAR_LIFE` — tabel, 7 kolom

`ID VARCHAR2(100)`, `TREATYYEAR VARCHAR2(100)`, `UNDERWRITINGYEAR VARCHAR2(100)`,
`USERID VARCHAR2(100)`, `TGLUPDATE DATE`, `STARTDATE DATE`, `ENDDATE DATE`.

### `RATE_LIFE` — **VIEW**, 8 kolom — dilengkapi 26 September 2026 malam

> Ralat: catatan sebelumnya menyebutnya "tabel, enam kolom pertama terbaca". Ia **view** atas
> `M_RATE_LIFE` *(`ID VARCHAR2(10)`, `JSONDATA CLOB`)*: `SELECT a.ID, a.JSONDATA.IDUSEDBY,
> a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER, a.JSONDATA.CONTRACT, a.JSONDATA.AGE,
> a.JSONDATA.RATE FROM M_RATE_LIFE a`.

| # | Kolom | Tipe | Isi (agregat) |
| ---: | --- | --- | --- |
| 1 | `ID` | `VARCHAR2(10)` NOT NULL | 98.305 baris, seluruhnya unik |
| 2 | `IDUSEDBY` | `VARCHAR2(4000)` | teks angka; 348 nilai berbeda; nol NULL — **kunci pencarian** `GetRateRetro` |
| 3 | `USEDBY` | `VARCHAR2(4000)` | 347 nilai berbeda *(nama; tidak disalin)* |
| 4 | `TYPE` | `VARCHAR2(4000)` | **NULL di seluruh baris** |
| 5 | `GENDER` | `VARCHAR2(4000)` | `U` 98.004 · `M` 201 · `F` 99 · NULL 1 |
| 6 | `CONTRACT` | `VARCHAR2(4000)` | teks angka 0–120, 88 nilai berbeda; NULL 10.203 |
| 7 | `AGE` | `VARCHAR2(4000)` | teks angka 0–120, 102 nilai berbeda; NULL 3 |
| 8 | `RATE` | `VARCHAR2(4000)` | nol NULL; **90.436 baris memuat koma** *(desimal Indonesia)*, 3.226 memuat titik, 4.113 bernilai `0`; panjang 1–20; nilai bertitik/bulat terbesar 2.182. Bentuk: `angka,angka` 89.872 · `angka.angka` 3.224 · bulat 4.645 · **tidak polos 564** *(556 berspasi tepi; 2 memuat koma dan titik sekaligus)* |

⚠️ Tidak ada index *(view atas CLOB — tiap pembacaan mengurai JSON; query agregat memakan
puluhan detik)*. Kombinasi `(IDUSEDBY, GENDER, AGE, CONTRACT)` = 86.981, **2.693 di antaranya ganda**
→ dengan `GetRateRetro` tanpa `ORDER BY`, baris yang "menang" di Pega tidak ditentukan; pembaca Go
mengurutkan dan **melaporkan** ambiguitas *(brief modul §5 ah)*. Pega mengurai `RATE` lewat
`@toDecimal(@replaceAll(.CARI4, ",", "."))` `[terverifikasi]` — koma memang bentuk yang diharapkan.

### `RETROCESSIONLIFE` — definisi dan agregat (26 September 2026 malam)

Definisi view: `SELECT a.ID, a.JSONDATA.IDTREATYYEAR_LIFE, a.JSONDATA.PERCENTSHARE, a.JSONDATA.RATE,
a.JSONDATA.COMMISION, a.JSONDATA.OVR_COMM, a.JSONDATA.TREATYTYPEID, b.JSONDATA.Note AS
TREATYTYPENAME, a.JSONDATA.TREATYSTARTDATE, a.JSONDATA.TREATYENDDATE, a.JSONDATA.USERID,
a.JSONDATA.TGLUPDATE, a.JSONDATA.REINSURERNAME FROM M_RETROCESSIONLIFE a, M_REINSURANCETYPE b WHERE
a.JSONDATA.TREATYTYPEID = b.ID`. `M_RETROCESSIONLIFE` **juga** punya kolom bertipe di samping
`JSONDATA`: `IDTREATYYEAR_LIFE VARCHAR2(10)`, `PERCENTSHARE`/`RATE`/`COMMISION`/`OVR_COMM` `NUMBER`,
`TREATYTYPEID VARCHAR2(10)`, `TREATYSTARTDATE`/`TREATYENDDATE VARCHAR2(10)`, `USERID VARCHAR2(100)`,
`REINSURERNAME VARCHAR2(1000)`, `TREATYTYPENAME VARCHAR2(100)`. Di DEV **13 dari 13 baris**: nilai
kolom bertipe **sama persis** dengan nilai JSON-nya *(`PERCENTSHARE`, `RATE`, `COMMISION`,
`OVR_COMM`, `IDTREATYYEAR_LIFE`, `TREATYSTARTDATE` dibandingkan sebagai teks)*, dan seluruh
`TREATYTYPEID` punya pasangan di `M_REINSURANCETYPE`. `M_REINSURANCETYPE`: `ID VARCHAR2(5)`, `OLD_LJR_ID CHAR(2)`, `OLD_LJT_ID
CHAR(3)`, `JSONDATA CLOB`.

| Agregat DEV | Nilai |
| --- | --- |
| Baris | **13** |
| `RATE` | 1 NULL; sisanya angka satu karakter, tanpa koma |
| `PERCENTSHARE` | angka bertitik, 5–100; **jumlah per treaty-year = 100** *(total 500 untuk 5 treaty-year)* |
| `COMMISION`, `OVR_COMM` | seluruhnya angka bertitik/bulat |
| `TREATYSTARTDATE`/`TREATYENDDATE` | 10 karakter, pola `DD/MM/YYYY` *(bagian pertama tanggal akhir > 12 di seluruh baris)* |
| `IDTREATYYEAR_LIFE` | `1000032` (4) · `1000033` (4) · `1000034` (3) · `1000035` (1) · `1000036` (1) — **nol** yang cocok dengan `TREATYYEAR_LIFE.ID` |
| Jenis treaty (`TREATYTYPEID`/`TREATYTYPENAME`) | `10196` QS (4) · `10197` 2ND QS (4) · `10198` SURPLUS (3) · `10199` 2ND SURPLUS (1) · `10200` OR (1) |
| Reinsurer berbeda | 4 *(nama tidak disalin)* |

### `TREATYYEAR_LIFE` — isi DEV (2 baris konfigurasi; bukan data orang)

| `ID` | `TREATYYEAR` | `UNDERWRITINGYEAR` | `STARTDATE` | `ENDDATE` |
| --- | --- | --- | --- | --- |
| `1000078` | `2018` | `2018` | 2018-01-01 | 2024-12-31 |
| `1000079` | `2025` | `2025` | 2025-01-01 | 2036-12-31 |

⚠️ **Tidak ada kolom `IDR`/`USD`** di sini, padahal `SpreadingClaimLife_Act` membaca `.IDR`/`.USD`
per baris treaty-year *(kaskade kapasitas)* dan `GetJsonProductLife` dimulai `select * from
treatyyear_life` *(tanpa skema)*. `[data DBA]` objek `treatyyear_life` yang dilihat koneksi Pega dan
sumber kapasitas `IDR`/`USD`.

### `PRODUCT_LIFE` — VIEW 36 kolom atas `M_PRODUCT_LIFE.JSONDATA`

`M_PRODUCT_LIFE`: `ID VARCHAR2(6)`, `JSONDATA CLOB`, `RIRISKID VARCHAR2(10)`, `RIRISK VARCHAR2(100)`.
Kolom view *(semua `VARCHAR2(4000)` kecuali `ID VARCHAR2(6)`)*: `TYPE`, `TYPE_CEDING`, `CEDING`,
`CEDINGID`, `SOBNAME`, `SOBID`, `CAUSEID`, `GRUP`, `PRODUCTNAME`, `PRODUCTCODE`, `PRODUCTTYPEID`,
`PRODUCTTYPE`, `RIRISKID`, `RIRISK`, `RIRATEID`, `RIRATE`, `RICOMMID`, `RICOMM`, `INWARDNAME`,
`UNDERWRITINGLIMITLIST`, `OUTWARDNAMEID`, `OUTWARDNAME`, `OUTWARDRATEID`, `OUTWARDRATE`,
`OUTWARDCOMMID`, `OUTWARDCOMM`, `BENEFITID`, `BENEFIT`, `CAUSE`, `OVR_COMM` *(dari
`OutwardList[0]`)*, `POLICYHODER`, `POLICYHODERNAME`, `TREATYNUMBER`, `CREATEOP`, `UPDATEOP`.
`OUTWARDRATEID` di sini **tingkat produk**; yang dipakai `SpreadingClaimLife_Act` adalah
`OUTWARDRATEID` **per plan** di dalam `PlanList` JSON *(plan yang `Name`-nya = `BusinessName`
klaim)* — tidak tersedia lewat view ini. ⛔ JSON produk tidak dibaca *(AC 38)*; lihat brief modul
§5 **ag**.

## Yang dokumen ini TIDAK putuskan

Kolom mana dari 85 yang menjadi "27 kolom polis dibaca hidup" (tiket 02/14), aturan baris negatif
jurnal balik (kolom mana yang bertanda negatif), dan rumus `PREMIUM_SPREADED_NET`. Ketiganya milik
work owner, Product+UW, atau spec Endorsement Life.
