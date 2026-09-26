# Katalog tabel warisan yang dipakai tiket 02, 03, dan 06 — instance pengembangan

`[data DBA — dibaca sendiri 26 September 2026 malam, belum dikonfirmasi DBA]`

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

### `RETROCESSIONLIFE` — **VIEW**, 13 kolom, **seluruhnya `VARCHAR2(4000)`**

`ID` (7), `IDTREATYYEAR_LIFE`, `PERCENTSHARE`, `RATE`, `COMMISION`, `OVR_COMM`, `TREATYTYPEID`,
`TREATYTYPENAME`, `TREATYSTARTDATE`, `TREATYENDDATE`, `USERID`, `TGLUPDATE`, `REINSURERNAME`.

⚠️ Angka spreading (`PERCENTSHARE`, `RATE`, `COMMISION`, `OVR_COMM`) tiba sebagai **teks** dari
view ini. Pembacanya wajib mengurai lewat `utils.ParseDecimal` dan **melaporkan** yang tidak
terurai, bukan menebak; pemisah desimal yang dipakai view belum diketahui — periksa agregat dulu.

### `TREATYYEAR_LIFE` — tabel, 7 kolom

`ID VARCHAR2(100)`, `TREATYYEAR VARCHAR2(100)`, `UNDERWRITINGYEAR VARCHAR2(100)`,
`USERID VARCHAR2(100)`, `TGLUPDATE DATE`, `STARTDATE DATE`, `ENDDATE DATE`.

### `RATE_LIFE` — tabel, enam kolom pertama terbaca

`ID VARCHAR2(10)`, `IDUSEDBY`, `USEDBY`, `TYPE`, `GENDER`, `CONTRACT` (semuanya `VARCHAR2(4000)`).
Sisanya belum dibaca; cacah baris belum dibaca. Executor tiket 03 melengkapinya dari katalog
sebelum menulis pembaca.

## Yang dokumen ini TIDAK putuskan

Kolom mana dari 85 yang menjadi "27 kolom polis dibaca hidup" (tiket 02/14), aturan baris negatif
jurnal balik (kolom mana yang bertanda negatif), dan rumus `PREMIUM_SPREADED_NET`. Ketiganya milik
work owner, Product+UW, atau spec Endorsement Life.
