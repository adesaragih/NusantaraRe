> ⚠️ **Disusun subagen baca-saja sesi c3 (09-10-2026), bukan hasil `/to-tickets`.** Diperiksa ulang sesi c3 (contoh, bukan
> seluruhnya): `START WITH 151028` + butir 74.2 (KEPUTUSAN L1066 / L1137) dan butir 77 (151028 < NB-184351); skema loader nol
> `T_NB_OPPORTUNITY`; nomor migrasi berikut nbfacin 214 (MODUL.md); berkas `services/loader` yang berstatus modified bertanggal
> 07–08-10-2026 (bukan ulah analisis ini). Skrip / keluaran mentah: scratchpad sesi c3 (`riset67/`), tidak masuk repo.

# Analisis tiket 72 — tahap A: migrasi polis lama Pega → tabel case Fac In

> Analisis BACA-SAJA, sesi riset 09-10-2026. Nol berkas repo diubah. Satu-satunya keluaran: berkas ini. Skrip bantu dan
> keluaran mentah ada di folder scratchpad yang sama (`gap*.py`, `presisi.py`, `baris_xml.py`, `ukur72_test.go`,
> `skema_efektif.tsv`, `ukur72*.tsv`). Kode Go dijalankan **lewat `go test -overlay`** (berkas uji tinggal di scratchpad,
> tidak masuk pohon repo). Salinan JSON hasil ubah dari 115 contoh XML korpus dibuat di scratchpad lalu **dihapus**.
> Nol akses basis data; hitungan yang butuh data produksi diberikan sebagai SQL (§3, §4). Nama orang tidak disalin.
>
> ⚠️ Pohon kerja `APP_RNM` punya banyak perubahan belum di-commit milik sesi lain (mis. `loader/aturan.go`,
> `amandemen.go`, `migrasi_sebagian_test.go`, migrasi 210–212 belum dilacak). Analisis ini membaca **pohon kerja**, bukan HEAD.

Label: `[terverifikasi]` = dibaca/diukur dengan bukti berkas+baris · `[dugaan]` · `belum terverifikasi`.

---

## 1. Status tiket 22, 23, 24

| Tiket | Isi | Status |
| --- | --- | --- |
| 22 `loader.Flatten` | seam murni dokumen → baris tabel flat; diuji fixture + 115 contoh | **dibangun**, `ready-for-human` |
| 23 migrasi tabel flat | 78 (kini 80) tabel dari DDL draf | **⏸ DITAHAN** — nol berkas ditulis (sebagai paket) |
| 24 pemuat repository + perkakas | baca JSON_POLIS → Flatten → tulis per penawaran | **needs-info**, nol kode |

- `[terverifikasi]` Tiket 22 (`docs/issues/22-loader-flatten.md` L29–30): "kelima kasus nyata kini **lolos mode ketat**;
  butir 69–72 diterapkan; uji mutasi 62/62". Uji: `flatten_test.go`, `kasus_test.go`, `skema_test.go`,
  `migrasi_sebagian_test.go` (`services/loader/`). Pengukuran 115 contoh: L88–131.
- `[terverifikasi]` Saya jalankan ulang Flatten **ketat** atas 115 contoh `D:\migrasi\RNM\DDL\CONTOH\*.xml` (diubah ke
  JSON, PageList → larik): **115/115 lolos**, nol galat. Instrumen diuji atas butir yang jawabannya sudah diketahui —
  ketiganya cocok persis dengan tiket 22: COB Aneka 81 · FIRE 12 · Life 5 · MBUCar 6 · MarineCargo 6 · PA 5 (L100);
  penampung `IsCedingConfirm` **805** (L203); konflik tipe 68.1 (pengurai Python terpisah): `TYPE_OF_DISCOUNT` 94,
  `SHARE_OF_CEDING` 87, `PPN_CHECK` 74, `WARR_PAYMENT` 15, `MASTER_RATE_COVERAGE` 5, `EDM_CHARGE_FEE` 1 (L108–116).
- `[terverifikasi]` Tiket 23 ditahan butir 66, `docs/KEPUTUSAN-30-09-2026.md` L760: *"**"Flatten dulu, migrasi
  ditahan."** `loader.Flatten` (murni) dikerjakan; berkas migrasi 180–219 ditulis setelah tim inti memutuskan presisi
  tabel flat. Penjaga inti tidak diubah"*. Butir 68.2 (L804): *"**Tetap keputusan tim inti.** Agent menulis usulan
  (`docs/USULAN-PRESISI-TABEL-FLAT.md`) … Tiket 23 tetap ditahan"*. Tiket 23 L13: *"⏸ ditahan; tidak satu berkas
  migrasi pun ditulis"*.
- `[terverifikasi]` **Tetapi** tabel flat kemudian dibuat **sebagian per tiket layar** (butir 78.4, migrasi 182–213)
  dengan uang `NUMBER(38,8)` atas dasar ADR-0016 (mis. komentar `193_t_coveragelist.sql` L3–4, L10). Jadi butir 66 untuk
  kolom yang sudah ada **sudah terlampaui secara de facto**; yang masih ditahan hanyalah sisa kolom/tabel (§2).
- `[terverifikasi]` Tiket 24 (`24-pemuat-repository.md` L18–20): *"Blocked by: 22, 23 … menunggu 23"*. Keputusan
  yang sudah mengikat: 68.8 gagal per penawaran → dilewati, dicatat, ringkasan akhir (L42–44); ralat butir 69:
  *"muat untuk UJI boleh dengan penampung berisi; muat produksi / fase 1 baru boleh dinyatakan selesai bila penampung
  kosong"* (L47–51); "Fase 2 mengosongkan 78 tabel … hanya atas perintah eksplisit work owner" (L37–38). Terbuka: J-16
  penyaring "endorsement sudah jadi" (L57).

## 2. Kesenjangan keluaran Flatten vs tabel `T_*` yang ADA

### 2a. Skema sumber — dihitung dua cara

| Besaran | Cara 1 | Cara 2 | Hasil |
| --- | --- | --- | --- |
| Skema loader efektif (bangkitan + amandemen) | dump `skemaTabel` sesudah `init()` lewat uji overlay: **80 tabel · 1.445 kolom · 150 jalur** | rekonsiliasi delta: 1.391 kolom (tiket 23 L47, butir 76) + `T_FEALIST` 15 + amandemen sesudahnya 39 (Deduction2 1, gempa 1, scoring 2, NBStatus 1, KeUW 2, bangunan 3, sekitar 18, item 6, kerugian 5; `amandemen.go` L200–312) = **1.445**; tabel 79+1 = **80**; jalur 149+1 = **150** | cocok |
| Tabel/kolom migrasi (nbfacin 180–213 + `T_WORK_POLIS` premiumlistlife 050/057/059/063) | pengurai pernyataan (CREATE/ALTER ADD/MODIFY/DROP): **32 tabel · 463 kolom** | pencacah baris regex: **463**; A\B = ∅, B\A = ∅ | cocok |

Asal kolom skema loader (dump `asalKolom`): medan dokumen 911 · diisi Flatten 331 · diisi repository 171 · sengaja
kosong 32 = 1.445.

### 2b. Tabel

- `[terverifikasi]` **26 tabel** ada di keduanya — seluruhnya dibaca kode `repository/*.go` layar case:
  `T_WORK_POLIS, T_GENERAL_POLIS, T_QUOTATIONDATA, T_CEDINGCOLIST, T_LOCATIONLIST, T_PROPERTY, T_RISKLOCATION,
  T_BUILDINGCONSTRUCTION, T_SURROUNDINGRISK, T_PROPERTYITEMLIST, T_OCCUPATIONLIST, T_TABLEOFLIMIT, T_FEALIST,
  T_LISTCAUSEOFLOSS, T_COINSDATA, T_COVERAGELIST, T_DEDUCTIBLELIST, T_SPREADINGLIST, T_CEDINGCEDANTLIST, T_CURRENCYLIST,
  T_LISTINSTALLMENT, T_SCORINGRISK, T_SCORINGRESULT, T_DATASCORINGRISKLIST, T_SCORING_FACTOR, T_SCORING_OPTION`.
- `[terverifikasi]` **54 tabel loader belum ada** (828 kolom): seluruh 36 kembaran `T_FR_*` (FacRetro) + `T_ANEKALIST,
  T_ADDITIONALCOVERAGE, T_CURRENCY, T_CEDING_CURRENCYLIST, T_RETROLIST, T_VEHICLELIST, T_VEHICLEHE, T_PERSONLIST,
  T_ASMHEIR, T_CARGOLIST, T_SHIP, T_ADDITIONALSHIP, T_LC, T_MARINEHULL, T_AVIATIONHULL, T_ACCESSORYLIST, T_QUOTATIONLIST,
  T_COVERAGEDATALIST, T_SURVEYREPORTLIST, T_FACRETROLIST, T_FACRETRODETAILS`. Rekonsiliasi: 617 + 828 = 1.445.
- `[terverifikasi]` **6 tabel migrasi di luar rancangan loader** (Flatten tidak menghasilkan barisnya): `T_NB_OPPORTUNITY`
  (180), `T_CLAUSELIST` (197; komentar L5: "ClauseList TIDAK ada di workbook rancangan … kosong di seluruh 108 contoh"),
  `T_LOSSRECORD`/`T_LOSSRECORDLIST` (206, tingkat case), `T_CORRESPONDENCELIST`/`T_RECIPIENTLIST` (207). Rekonsiliasi
  migrasi: 381 + 4 (kolom `T_WORK_POLIS` di luar rancangan) + 78 (enam tabel ini) = 463.
- ⛔ `[terverifikasi]` **`T_NB_OPPORTUNITY` adalah kesenjangan terpenting bagi layar**: `repository/kasus.go` L58–96
  `sqlBacaKasus` LEFT JOIN `o.CLASS_OF_BUSINESS, o.GROUP_BUSINESS_ID, o.ACCOUNT_ID, o.INSURED_ID …`, dan
  `services/lini.go` L151–166 menurunkan **BusinessType / lini** case dari `CLASS_OF_BUSINESS` + `GROUP_BUSINESS_ID`.
  Polis lama tanpa baris opportunity → lini "UNKNOWN" → tab FIRE tidak terbuka. Isian lamanya ada di dokumen
  (`QuotationData.BusinessName/BusinessCode/BusinessOldId/GroupPanel/InsuredID/InsuredName`) tetapi Flatten memetakannya
  ke kolom `T_QUOTATIONDATA` yang **belum ada**. Pemetaan dokumen → `T_NB_OPPORTUNITY` = keputusan (§6, D3).
  `[dugaan]` `GROUP_BUSINESS_ID` = BUSINESS.BUSINESSGROUPID, tidak sama dengan `GroupPanel` — belum terverifikasi.

### 2c. Kolom pada 26 tabel bersama — 617 kolom loader: **381 ada · 236 hilang**

Cara 2 untuk 381: konstanta `TestMigrasiFlatSebagianCocokRancangan` `diperiksa != 344` (23 tabel,
`migrasi_sebagian_test.go` L196) + `T_SPREADINGLIST` 16 + `T_FEALIST` 15 + `T_WORK_POLIS` 6 = **381** — cocok.

Hilang per tabel (asal): `T_COVERAGELIST` 54 · `T_GENERAL_POLIS` 54 · `T_QUOTATIONDATA` 51 · `T_DEDUCTIBLELIST` 19 ·
`T_WORK_POLIS` 13 · `T_PROPERTYITEMLIST` 11 · `T_SPREADINGLIST` 11 · `T_CURRENCYLIST` 7 · `T_LISTINSTALLMENT` 7 ·
`T_OCCUPATIONLIST` 5 · `T_LOCATIONLIST` 3 · `T_PROPERTY` 1 (`TOTAL_TSI`); 14 tabel lain utuh. Menurut asal: medan
dokumen 209 · sengaja kosong (FK V-47) 16 · repository 6 (`OLD_POLIS_ID`, `PROD_KE`, `NOURUT`, `PUTARAN`,
`STS_KONVERSI`, `TGL_KONVERSI`) · Flatten 5 (`T_WORK_POLIS.IDPEGA/JENIS_WORK/NO_WORK`, `SOURCE_ID`, `CURRENCY_REF_ID`).

**Yang sungguh berisi data** (Flatten atas 115 contoh; nilai tidak dicetak): **211 kolom hilang terisi, 45.068 nilai**.
Terbanyak: 48 kolom penunjuk `INDEX_*`/`IDX_*` (butir 72; mis. `T_DEDUCTIBLELIST.INDEX_DEDUCTIBLE` 2.930,
`T_COVERAGELIST.INDEX_COVERAGE` 2.397), `T_COVERAGELIST.IS_SHOW_FORMULA` 2.323, `RISK_ZIP_CODE` 2.313,
`T_SPREADINGLIST.IS_OLD_DATA` 2.304, `MAX/MIN_STANDARD_RATE` 928, `T_PROPERTYITEMLIST.TOTAL_PREMIUM_NUSANTARA_RE` 470,
`T_PROPERTY.TOTAL_TSI` 194; per dokumen (115/115): `T_GENERAL_POLIS.POLICY_NO`, `TOTAL_TSI_NUSA_RE`,
`TOTAL_PREMI_NUSA_RE`, `T_QUOTATIONDATA.BUSINESS_TYPE/BUSINESS_CODE/BUSINESS_NAME/INSURED_ID/INSURED_NAME/MARKETING_CODE`,
`T_WORK_POLIS.PIC`; endorsement (13/115): `ENDORSEMENT_NO`, `ENDORSMENT_REASON`, `EDM_*`, `OLD_POLICY_NO`,
`PAY_EDM_*`. Jendela FIRE saja (12 dokumen): **126 kolom hilang terisi, 36.168 nilai**. ⚠️ Satu cara saja (Flatten);
cara kedua belum dibuat — label: `[terverifikasi satu cara]`.

Nilai terisi total 253.264; yang mendarat di kolom yang ada 191.871 (FIRE: 208.404 / 171.026).

### 2d. Baris di tabel yang belum ada

| Jendela | Cara 1 (Flatten) | Cara 2 (pengurai XML berjalur, `baris_xml.py`) | Hasil |
| --- | --- | --- | --- |
| 115 contoh | 1.443 baris / 53 tabel (terbanyak `T_CURRENCY` 368, `T_ANEKALIST` 127, `T_ADDITIONALCOVERAGE` 125, `T_CEDING_CURRENCYLIST` 103, `T_RETROLIST` 80) | 1.411 / 51 | **belum punya data** (berselisih; cara 2 tidak meniru V-27/V-30 dkk.) |
| 12 FIRE | **77 / 23** (FacRetro `T_FR_*` + `T_CEDING_CURRENCYLIST` 14) | **77 / 23** | cocok |

### 2e. Tipe & presisi (di luar beda `NOT NULL`)

| Kelompok | Jumlah | Cara 2: peta `selisih` `migrasi_sebagian_test.go` L72–184 |
| --- | ---: | --- |
| identitas `ID`/`PARENT_ID` `NUMBER` → `NUMBER(19)` (A92) | 43 | 39 + 4 tabel di luar daftar uji (`T_SPREADINGLIST`, `T_FEALIST`) = 43 ✓ |
| uang/rate `NUMBER` → `NUMBER(38,8)` (ADR-0016) | 58 | 54 = 58 − 5 `T_SPREADINGLIST` + 1 `T_LISTINSTALLMENT.DEDUCTION2` ✓ |
| kode `NUMBER` → `NUMBER(5)` | 5 | 5 ✓ |
| angka → teks (68.1 / K51-3) | 4 (`PPN_CHECK`, `SHARE_OF_CEDING`, `PCT_LIMIT`, `INSTALLMENT_PERCENTAGE`) | 4 ✓ |
| ⚠️ teks → angka: `T_LISTINSTALLMENT.DEDUCTION2` rancangan `VARCHAR2(50)`, migrasi 201 `NUMBER(38,8)` | 1 | — |
| ⚠️ `T_LISTCAUSEOFLOSS.CURRENCY`: 209 membuatnya NULL; Flatten mengisi `'UNKNOWN'` (K-069) | 1 | — |

Akibat bagi pemuat `[terverifikasi]` dari kode: Flatten mengeluarkan `DEDUCTION2` sebagai **teks** (rancangan
`VARCHAR2`) yang harus diubah ke angka saat tulis; `NUMBER(5)` membulatkan diam-diam bila nilai berpecahan `[dugaan]`
(perilaku Oracle, tidak diuji).

**Presisi — sumber angka 10.876 / 82** `[terverifikasi]`: `docs/USULAN-PRESISI-TABEL-FLAT.md` L23–26 (diukur
02-10-2026 atas keluaran Flatten, 115 contoh diubah ke JSON di luar repo; putaran kedua, nol ujung tidak dihitung).
Saya ukur ulang dua cara:

| Ukuran (> 8 desimal bermakna) | USULAN | Cara 1: Flatten (overlay) | Cara 2: pengurai XML (`presisi.py`) | Hasil |
| --- | ---: | ---: | ---: | --- |
| semua kolom NUMBER, 115 contoh | 10.876 / 82 kol. | **10.866 / 82** | 10.714 / 72 | **belum punya data** untuk angka pasti (ketiganya berselisih; orde ~10,7–10,9 ribu) |
| 5 kolom teratas | 2.249 · 2.241 · 2.224 · 1.813 · 440 | 2.249 · 2.241 · 2.224 · 1.813 · 440 | 2.248 · 2.240 · 2.223 · 1.809 · 440 | cara 1 = USULAN persis |
| **hanya kolom yang SUDAH ADA `NUMBER(38,8)`** | — | **9.969 / 29 kol.** | 9.841 / 23 kol. | **belum punya data** (dua cara berselisih) |
| > 38 digit bermakna | 3 (`PAY_NET_PREMIUM`) | — | 3 (`T_CURRENCYLIST.PAY_NET_PREMIUM`, kolom yang ada) | cocok dengan USULAN B |

Makna untuk tiket 72: **± 9,8–10 ribu nilai di 115 contoh akan DIBULATKAN Oracle ke 8 desimal** bila dimuat ke kolom
`NUMBER(38,8)` yang sudah ada (terbanyak `T_SPREADINGLIST.PREMIUM_SPREADED`, `T_COVERAGELIST.PREMI_NUSANTARA_RE` /
`PREMIUM` / `RATE`), dan **3 nilai `PAY_NET_PREMIUM` (> 38 digit) akan ditolak/dibulatkan** `[dugaan]` (perilaku Oracle
tidak diuji). Jendela FIRE (12 dokumen, Flatten): 10.246 nilai / 47 kolom, 9.698 / 21 di kolom yang ada. Ini 115 contoh,
**bukan** bukti untuk produksi (USULAN L65–66).

### 2f. Kolom sistem & kunci

| Kolom | Flatten | Repository (aplikasi) | Tabel bersama yang punya kolomnya |
| --- | --- | --- | --- |
| `ID` | **tidak** — `Baris.Kunci` 1..N lokal (`flatten.go` L60–68) | `SEQ_T_<TABEL>.NEXTVAL` `NUMBER(19)`; case = `NB-<n>` dari `SEQ_WORK_POLIS_NB` (`casenb.go` L99–106) | 26 |
| `PARENT_ID` | **tidak** — `Baris.Induk` (0 = `T_WORK_POLIS`) | ID induk; tabel berinduk case = ID case `VARCHAR2(32)` (butir 76.1) | 24 |
| `IDPEGA`, `COB_GROUP` | diisi (`lahir`, L293–322) | **tidak diisi** (mis. `cedant.go` L88, `ceding.go` L50 — INSERT tanpa keduanya) | 25 / 25 (`T_WORK_POLIS` tanpa `IDPEGA`) |
| `SEQ_NO` | posisi larik (1 untuk halaman tunggal) | urutan grid | 15 |
| `ROW_UID` | UUID v5 **deterministik** dari IDPEGA+tabel+posisi (`uidBaris` L324–335) | UUID v4 acak (`ceding.go` L110) | 15 |
| `PARENT_TABLE`, `SRC_PATH` | diisi | diisi (coverage/deductible) | 4 |

Hitungan cara 1 = himpunan skema loader∩migrasi; cara 2 = daftar `HILANG` per tabel — nol kolom sistem
(`ID/PARENT_ID/SEQ_NO/ROW_UID/PARENT_TABLE/SRC_PATH`) muncul di daftar hilang, kecuali `T_WORK_POLIS.IDPEGA`.
Peluang: **`IDPEGA` terisi = baris hasil migrasi** (aplikasi membiarkannya NULL) — penanda alami tanpa kolom baru.

## 3. Cakupan sumber `JSON_POLIS`

- `[terverifikasi]` DDL `D:\migrasi\RNM\DDL\JSON_POLIS.txt` L2–24: `IDPEGA VARCHAR2(50)` PK, `DATA_JSON CLOB`,
  `NOPOLIS` NOT NULL, `NOENDORS`, `PRODKE VARCHAR2(5)`, `STS_KONVERSI`, `TGL_KONVERSI`, `USERNAME`, `TGL_PROD`;
  `UNIQUE (NOPOLIS, IDPEGA)`; **`VALID_JSON … DISABLE`** (L17) → dokumen tak sah mungkin ada.
- `[terverifikasi]` `JSON_POLIS` **bercampur lini/modul**: pemuat `nbtreatyin` memilah dengan `pxObjClass`
  `ASM-FW-GISFW-Data-PolicyTreatyIn` (`modul/nbtreatyin/backend/models/dokumenlama.go` L40, L533–534) dan generasi
  `PRODKE` 0 (`repository/lama.go` L38–41). Ke-115 contoh Fac In ber-`pxObjClass` akar
  `ASM-FW-GISFW-Data-OfferFacIn` (dua cara: pengurai XML 115; grep tag 115).
- `[terverifikasi]` Contoh korpus (bukan sensus produksi): nama berkas NB 102 · EDM 13; `BusinessType` (dua cara:
  regex kemunculan pertama = urai `QuotationData`, cocok): Bonding 28 · FireStyle1 10 · Aneka 10 · Liability 9 ·
  BondingKBG 8 · MBUCar 6 · MarineCargo 6 · Life 5 · PA 5 · 4 masing-masing AviationHull/CustomBond/GolfInsurance/HE/
  MarineHull · MBD 3 · FireStyle2 2 · LandRig 2 · ElectronicEquipment 1 → FIRE = 12 / 115.
- `[terverifikasi]` Sistem baru **juga menulis** `JSON_POLIS` saat selesai Policy: `IDPEGA = 'ASM-FW-GISFW-WORK ' + caseId`
  (`repository/submit.go` L34, `services/submit.go` L169), `PRODKE '0'`, `DATA_JSON = {"IDNewBisnis": caseId}`
  (`docs/issues/69-jalur-policy.md` L345–348). ⚠️ `IDNewBisnis` juga ada di 105/115 dokumen Pega lama — penyaring
  "baris sistem baru" harus memakai **ketiadaan `QuotationData`** dan/atau keberadaan baris `T_WORK_POLIS` LINI FAC.
- **Status case lama**: `[terverifikasi]` `pyStatusWork` **tidak ada** di 115 contoh (grep 0); medan `py*` memang dibuang
  Flatten (`aturan.go` alasanPy). Yang ada di akar: `PolicyStatus` — dua cara cocok: urai akar 51 (Complated 50,
  Incomplated 1), Flatten `T_WORK_POLIS.POLICY_STATUS` terisi 51. Arti kodenya **belum terverifikasi** (OQ-020).
  Posisi terakhir per case ada di `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (`IDPEGA`, `NOURUT`, `POSISI`, `PUTARAN`; DDL L3–).
- Endorsement: `[terverifikasi]` J-16 masih terbuka (`_DAFTAR-ISSUE-TERBUKA.md` L316, butir 69.6); View Policy menampilkan
  versi dari kolom relasional JSON_POLIS (K71-2) dan menandai `sistemBaru` bila `T_WORK_POLIS.ID = SUBSTR(IDPEGA,19)`
  LINI FAC (`repository/polis.go` L59–62) — **memuat EDM ke `T_WORK_POLIS` LINI FAC otomatis membuat versi EDM tampil
  sebagai layar case NB**.

**SQL untuk work owner / DBA** (baca saja; `{skema}` = skema aplikasi; Oracle 12c+ untuk `JSON_VALUE`/`IS JSON`):

```sql
-- Q1 komposisi: kelas dokumen x jenis case x generasi x NOENDORS x JSON sah
SELECT JSON_VALUE(j.DATA_JSON, '$.pxObjClass' NULL ON ERROR)                        AS KELAS_DOK,
       REGEXP_SUBSTR(SUBSTR(j.IDPEGA, INSTR(j.IDPEGA, ' ', -1) + 1), '^[A-Z]+')      AS JENIS,
       CASE WHEN j.PRODKE IS NULL OR TRIM(j.PRODKE) = '0' THEN '0' ELSE '>=1' END   AS GEN,
       NVL2(j.NOENDORS, 'ISI', 'NULL')                                               AS NOENDORS,
       COUNT(*) AS N,
       SUM(CASE WHEN j.DATA_JSON IS JSON THEN 1 ELSE 0 END)                          AS N_JSON_SAH,
       SUM(CASE WHEN JSON_EXISTS(j.DATA_JSON, '$.QuotationData') THEN 0 ELSE 1 END)  AS N_TANPA_QUOTATION
  FROM POOLDATA.JSON_POLIS j
 GROUP BY JSON_VALUE(j.DATA_JSON, '$.pxObjClass' NULL ON ERROR),
          REGEXP_SUBSTR(SUBSTR(j.IDPEGA, INSTR(j.IDPEGA, ' ', -1) + 1), '^[A-Z]+'),
          CASE WHEN j.PRODKE IS NULL OR TRIM(j.PRODKE) = '0' THEN '0' ELSE '>=1' END,
          NVL2(j.NOENDORS, 'ISI', 'NULL')
 ORDER BY N DESC;

-- Q2 Fac In per BusinessType / GroupPanel / generasi
SELECT JSON_VALUE(DATA_JSON, '$.QuotationData.BusinessType' NULL ON ERROR) AS BUSINESS_TYPE,
       JSON_VALUE(DATA_JSON, '$.QuotationData.GroupPanel'   NULL ON ERROR) AS GROUP_PANEL,
       CASE WHEN PRODKE IS NULL OR TRIM(PRODKE) = '0' THEN '0' ELSE '>=1' END AS GEN, COUNT(*) AS N
  FROM POOLDATA.JSON_POLIS
 WHERE JSON_VALUE(DATA_JSON, '$.pxObjClass' NULL ON ERROR) = 'ASM-FW-GISFW-Data-OfferFacIn'
 GROUP BY ROLLUP (JSON_VALUE(DATA_JSON, '$.QuotationData.BusinessType' NULL ON ERROR),
                  JSON_VALUE(DATA_JSON, '$.QuotationData.GroupPanel' NULL ON ERROR),
                  CASE WHEN PRODKE IS NULL OR TRIM(PRODKE) = '0' THEN '0' ELSE '>=1' END);

-- Q3 PRODKE / NOENDORS / versi per polis (Fac In)
SELECT TRIM(PRODKE) AS PRODKE, COUNT(*) AS N, SUM(NVL2(NOENDORS, 1, 0)) AS N_NOENDORS_ISI
  FROM POOLDATA.JSON_POLIS
 WHERE JSON_VALUE(DATA_JSON, '$.pxObjClass' NULL ON ERROR) = 'ASM-FW-GISFW-Data-OfferFacIn'
 GROUP BY TRIM(PRODKE) ORDER BY 1;
SELECT COUNT(*) AS PRODKE_BUKAN_ANGKA FROM POOLDATA.JSON_POLIS
 WHERE PRODKE IS NOT NULL AND NOT REGEXP_LIKE(TRIM(PRODKE), '^[0-9]+$');
SELECT JML_VERSI, COUNT(*) AS N_POLIS FROM (SELECT NOPOLIS, COUNT(*) AS JML_VERSI FROM POOLDATA.JSON_POLIS
 WHERE JSON_VALUE(DATA_JSON, '$.pxObjClass' NULL ON ERROR) = 'ASM-FW-GISFW-Data-OfferFacIn' GROUP BY NOPOLIS)
 GROUP BY JML_VERSI ORDER BY 1;
SELECT COUNT(*) AS PASANGAN_GANDA FROM (SELECT NOPOLIS, TRIM(PRODKE) FROM POOLDATA.JSON_POLIS
 GROUP BY NOPOLIS, TRIM(PRODKE) HAVING COUNT(*) > 1);

-- Q4 IDPEGA yang SUDAH punya baris T_WORK_POLIS (per LINI) + baris tulisan sistem baru
SELECT w.LINI, COUNT(*) AS N FROM POOLDATA.JSON_POLIS j
  JOIN {skema}.T_WORK_POLIS w ON w.ID = SUBSTR(j.IDPEGA, INSTR(j.IDPEGA, ' ', -1) + 1)
 GROUP BY w.LINI;
SELECT COUNT(*) AS N_SISTEM_BARU FROM POOLDATA.JSON_POLIS
 WHERE JSON_EXISTS(DATA_JSON, '$.IDNewBisnis') AND NOT JSON_EXISTS(DATA_JSON, '$.QuotationData');

-- Q5 tabrakan nomor case
SELECT LAST_NUMBER FROM ALL_SEQUENCES WHERE SEQUENCE_NAME = 'SEQ_WORK_POLIS_NB';
SELECT MAX(TO_NUMBER(REGEXP_SUBSTR(IDPEGA, '[0-9]+$'))) AS NB_MAKS,
       SUM(CASE WHEN TO_NUMBER(REGEXP_SUBSTR(IDPEGA, '[0-9]+$')) >= 151028 THEN 1 ELSE 0 END) AS N_DI_ATAS_AWAL_SEQ
  FROM POOLDATA.JSON_POLIS WHERE REGEXP_LIKE(IDPEGA, ' NB-[0-9]+$');
SELECT COUNT(*) AS ID_CASE_LINTAS_KELAS FROM (
  SELECT SUBSTR(IDPEGA, INSTR(IDPEGA, ' ', -1) + 1) AS CID
    FROM POOLDATA.JSON_POLIS
   GROUP BY SUBSTR(IDPEGA, INSTR(IDPEGA, ' ', -1) + 1)
  HAVING COUNT(DISTINCT JSON_VALUE(DATA_JSON, '$.pxObjClass' NULL ON ERROR)) > 1);
SELECT COUNT(*) AS BENTROK_LINI_LAIN FROM {skema}.T_WORK_POLIS w
 WHERE w.LINI <> 'FAC' AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j
   WHERE SUBSTR(j.IDPEGA, INSTR(j.IDPEGA, ' ', -1) + 1) = w.ID
     AND JSON_VALUE(j.DATA_JSON, '$.pxObjClass' NULL ON ERROR) = 'ASM-FW-GISFW-Data-OfferFacIn');

-- Q6 volume & polis tanpa JSON_POLIS (hanya FACINPRODUCTION)
SELECT TO_CHAR(TGL_INPUT, 'YYYY') AS THN, COUNT(*) AS N, ROUND(SUM(DBMS_LOB.GETLENGTH(DATA_JSON)) / 1048576) AS MB,
       MAX(DBMS_LOB.GETLENGTH(DATA_JSON)) AS MAKS_KAR
  FROM POOLDATA.JSON_POLIS
 WHERE JSON_VALUE(DATA_JSON, '$.pxObjClass' NULL ON ERROR) = 'ASM-FW-GISFW-Data-OfferFacIn'
 GROUP BY TO_CHAR(TGL_INPUT, 'YYYY') ORDER BY 1;
SELECT COUNT(DISTINCT p.NOPOLIS) AS POLIS_TANPA_JSON FROM POOLDATA.FACINPRODUCTION p
 WHERE NOT EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.NOPOLIS = p.NOPOLIS);

-- Q7 status lama: PolicyStatus akar, STS_KONVERSI, posisi terakhir riwayat
SELECT JSON_VALUE(DATA_JSON, '$.PolicyStatus' NULL ON ERROR) AS POLICY_STATUS, STS_KONVERSI, COUNT(*) AS N
  FROM POOLDATA.JSON_POLIS
 WHERE JSON_VALUE(DATA_JSON, '$.pxObjClass' NULL ON ERROR) = 'ASM-FW-GISFW-Data-OfferFacIn'
 GROUP BY JSON_VALUE(DATA_JSON, '$.PolicyStatus' NULL ON ERROR), STS_KONVERSI;
SELECT POSISI, COUNT(*) AS N FROM (
  SELECT h.POSISI, ROW_NUMBER() OVER (PARTITION BY h.IDPEGA ORDER BY h.TGL_INP DESC, h.NOURUT DESC) AS RN
    FROM POOLDATA.HISTORYAKSEPTASIPRODUCTION h
   WHERE h.IDPEGA IN (SELECT IDPEGA FROM POOLDATA.JSON_POLIS)) WHERE RN = 1
 GROUP BY POSISI ORDER BY N DESC;
```

⚠️ Q7 kedua memakai `TGL_INP`/`NOURUT` sebagai urutan — urutan "baris terakhir" yang sah (V-48) **belum terverifikasi**.

## 4. Pemetaan identitas `IDPEGA` → case `T_WORK_POLIS`

- `[terverifikasi]` ID case = `SUBSTR(IDPEGA, 19)` (awalan `ASM-FW-GISFW-WORK ` 18 karakter; `repository/submit.go` L34,
  `following_test.go` L44–45); Flatten `uraiIDPega` mengambil ruas sesudah spasi terakhir, jenis NB/RNW/EDM
  (`flatten.go` L210–224). Kolom tulis case (`casenb.go` L76–78): `ID, LINI, POSITION, FLAG_ONGOING_POLICY, STATUS_WORK,
  CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE`, `LINI = 'FAC'` (L35–37, butir 76.2). Kolom `T_WORK_POLIS` yang
  ada = 10 (050/057/059/063); 13 kolom rancangan (`IDPEGA`, `NO_WORK`, `PIC`, `POLICY_STATUS`, `NOURUT`, …) **belum ada**.
- **Tabrakan nomor.** Butir 74.2 (`KEPUTUSAN-30-09-2026.md` L1066), dikutip: *"**"Lanjut dari nomor terakhir Pega"** —
  penghitung mulai dari MAX nomor NB yang ada + 1; angka pastinya diisi saat migrasi dijalankan oleh work owner/DBA"*.
  `[terverifikasi]` `181_seq_work_polis_nb.sql` L13 kini `START WITH 151028` (dulu `{NB_MULAI}`, butir 76.5; MODUL.md:
  "sudah dijalankan di DEV 03-10-2026").
  - ⛔ `[terverifikasi]` Di contoh korpus `DDL\CONTOH`, **60 dari 102** nomor NB ≥ 151028 (maks NB-184351) — dua cara
    cocok: nama berkas 60/102, medan `IDNewBisnis` 60/105. `[dugaan]` contoh berasal dari lingkungan Pega lain (produksi)
    sedangkan 151028 diambil dari DEV. **Bila produksi memang melampaui 151028, migrasi ke T_WORK_POLIS akan bertabrakan
    PK dengan case baru, dan sequence yang sama harus di-reset sebelum go-live** → jalankan Q5 di lingkungan sasaran.
  - ⚠️ `[terverifikasi]` Pemuat `nbtreatyin` juga menyisip `T_WORK_POLIS` ber-ID `NB-<n>` (LINI `NONLIFE`,
    `nbtreatyin/backend/models/kasus.go` L18, L27) dan menganggap "ID sudah ada = dokumen sudah dimuat"
    (`services/pemuat.go` L147–156; cek di tabel treaty-nya sendiri, `repository/lama.go` L55–57). PK `T_WORK_POLIS` = `ID`
    saja (050). Bila Pega memakai deret NB yang sama untuk Fac In dan Treaty In, tidak ada tabrakan; bila tidak, ada —
    **belum terverifikasi** (Q5 `ID_CASE_LINTAS_KELAS`, `BENTROK_LINI_LAIN`). premiumlistlife memakai `NBLF-` (aman).
- **Status case selesai.** `[terverifikasi]` Tiket 69a-2 menulis `STATUS_WORK = 'Resolved-Completed'` dan **POSITION
  tidak diubah** (`services/selesaipolicy.go` L34; `69-jalur-policy.md` L234–238 K69a2-B1: *"POSITION dan NB_STATUS
  tetap … STATUS_WORK menjadi "Resolved-Completed" (K69-7)"*); `kasusSelesai` = awalan `Resolved-` → baca-saja di tahap
  mana pun (`services/bacasaja.go` L56–59). Kotak masuk menyaring `Resolved-`. Usulan untuk polis lama (keputusan WO):

  | Kolom | Usulan | Dasar |
  | --- | --- | --- |
  | `STATUS_WORK` | `Resolved-Completed` | K69-7; preseden pemuat nbtreatyin `TutupKasus` (`pemuat.go` L21–24, L170) |
  | `POSITION` | posisi terakhir `HISTORYAKSEPTASIPRODUCTION` (V-48, tiket 24 L25–26) **atau** NULL | 69a-2 tidak mengubahnya; data lama belum terverifikasi (Q7) |
  | `FLAG_ONGOING_POLICY` | `'1'` `[dugaan]` (tahap Policy, K69-8 baca-saja) | arti 0/1/2 di kode: 0 Offer, 1 Policy, 2 Accept/Waiting Bind (tiket 63) — pemetaan untuk data lama keputusan WO |
  | `CREATE_OP`/`CREATE_OP_NAME` | NULL (pola nbtreatyin "pembuat tidak dikarang") atau `USERNAME` JSON_POLIS | ⛔ `USERNAME` = login orang: disalin ke kolom, tidak pernah ke log (tiket 24 L26) |
  | `TGL_CREATE` | `JSON_POLIS.TGL_INPUT` | butir 76.4 `TGL_INPUT → TGL_CREATE` |

## 5. Idempotensi, transaksi, uji kering, laporan — yang sudah ada

| Kebutuhan | Di nbfacin sekarang | Preseden lain |
| --- | --- | --- |
| Isi deterministik | ✅ `ROW_UID` v5 deterministik; Flatten murni (tanpa jam/DB) | — |
| Satu polis satu transaksi | ✅ pola `db.Tx` milik pemanggil (`casenb.go` L3–4); penulis objek berkelompok + ID dari sequence (`tulisberkelompok.go` L1–18) | nbtreatyin `muat` satu transaksi (`pemuat.go` L147–177) |
| Lewati yang sudah dimuat | ❌ belum ada | nbtreatyin: ID sudah ada → `SudahDimuat` |
| Uji kering | ❌ | nbtreatyin `alat/pemuatlama` bawaan uji-kering, `-jalankan` untuk tulis (`main.go` L5–6, L46–47) |
| Laporan per polis tanpa nilai | sebagian: `Hasil.Diagnostik` (hitungan saja), `Penampung` | nbtreatyin `LaporanPemuat` CSV |
| Gagal per polis → lewati | keputusan 68.8 ada, kode belum | nbtreatyin `lap.Gagal` |
| Penampung kosong sebelum "selesai" | kebijakan ADR-0023/butir 69; korpus masih 805 entri `IsCedingConfirm` (P4, menunggu DBA) | — |

Rancangan diusulkan: kunci idempotensi = `T_WORK_POLIS.ID` LINI FAC **dan** `T_GENERAL_POLIS.IDPEGA` (aplikasi
menulis NULL) — ada + IDPEGA sama → lewati; ada + IDPEGA lain/NULL → galat "bentrok case sistem baru". Muat ulang paksa
(hapus lalu muat) = operasi destruktif, hanya atas perintah WO (tiket 24 L37–38).

## 6. Urutan kerja & keputusan yang dibutuhkan

**Urutan usulan**

1. Jalankan Q1–Q7 di lingkungan sasaran (DBA / WO) → angka cakupan, volume, tabrakan.
2. Putuskan D1–D8 di bawah.
3. Perkakas uji-kering dulu (pola nbtreatyin): baca JSON_POLIS → Flatten → **proyeksi ke kolom yang ada** → laporan
   per polis (kolom/tabel tak tertampung, nilai yang akan dibulatkan > 8 desimal, penampung) — nol tulis.
4. Migrasi tambahan (bila D2/D3 disetujui) — hanya **satu nomor tersisa (214)**: MODUL.md "Rentang migrasi `180-214`
   … nomor berikut nbfacin = 214" `[terverifikasi]`; satu berkas boleh memuat banyak pernyataan, atau minta rentang baru.
5. Penulis per polis satu transaksi (`-jalankan`), uji di skema uji, lalu DEV; produksi oleh DBA/WO.
6. Verifikasi: rekonsiliasi baris per tabel laporan vs `COUNT(*) … WHERE IDPEGA IS NOT NULL`; View Policy polis lama
   beralih ke layar case (`sistemBaru`).

**Keputusan work owner / DBA**

| # | Pertanyaan | Bahan |
| --- | --- | --- |
| D1 | Cakupan lini: FIRE saja (sejalan K69a2-6, layar non-FIRE belum ada — tiket 73) atau semua lini (tabel `T_ANEKALIST` dkk. belum ada)? | §2b, §3 Q2 |
| D2 | 211 kolom berisi data tanpa kolom (FIRE 126): tambah kolom (ALTER, rentang 214) / dibuang eksplisit / biarkan (ADR-0023: penampung wajib kosong) | §2c |
| D3 | Pemetaan `QuotationData` → `T_NB_OPPORTUNITY` (`CLASS_OF_BUSINESS`, `GROUP_BUSINESS_ID`, `ACCOUNT_ID`, `INSURED_ID`) — tanpa ini lini case = UNKNOWN | §2b |
| D4 | Presisi: terima pembulatan `NUMBER(38,8)` untuk ~9,8–10 rb nilai (115 contoh) atau kolom teks mentah (USULAN 68.2) — **tim inti** | §2e |
| D5 | Versi: NB saja (PRODKE 0) atau juga EDM/RNW? Bila EDM dimuat ke LINI FAC, View Policy membuka layar case NB untuknya; J-16 masih terbuka | §3 |
| D6 | Nomor case: verifikasi Q5; bila Pega > `LAST_NUMBER` → DBA reset `SEQ_WORK_POLIS_NB` (migrasi 181 tidak diubah — ralat lewat migrasi/DBA) | §4 |
| D7 | Nilai `STATUS_WORK` / `POSITION` / `FLAG_ONGOING_POLICY` / `CREATE_OP` polis lama | §4 tabel |
| D8 | DBA: kolom `IsCedingConfirm` (P4) di `HISTORYAKSEPTASIPRODUCTION`; `ALTER T_WORK_POLIS` milik premiumlistlife (koordinasi, risiko R1 tiket 23 L101–102); `T_LISTINSTALLMENT.DEDUCTION2` teks→angka; `T_LISTCAUSEOFLOSS.CURRENCY` `UNKNOWN` vs NULL | §2e |

*Tanpa nama orang, tanpa nilai data kasus; hanya nama kolom, kode lini, dan hitungan.*
