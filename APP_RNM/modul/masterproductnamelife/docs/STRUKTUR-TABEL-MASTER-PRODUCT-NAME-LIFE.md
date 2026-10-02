# Struktur Tabel — Master Product Name Life: PETA TABEL WARISAN yang ditulis

> Disusun 01-10-2026 (paket 10, brief `PROMPT-IMPLEMENTASI-MODUL-MASTER-PRODUCT-NAME-LIFE.md`).
>
> ⭐ **Ralat 02-10-2026 (keputusan work owner, K5 — OQ-MPNL-01 flat):** kalimat P1 di bawah (*"modul ini **tidak membuat
> satu tabel, sequence, maupun constraint pun**"*, *"Rentang migrasi `140`–`179` sengaja kosong"*) **tidak berlaku lagi**.
> Migrasi 140–147 membuat tabel flat `M_PRODUCTNAME_LIFE` + tujuh anak — bab *"Tabel FLAT produk"* di ekor dokumen ini.
> Kedua tabel JSON di bawah tetap **dibaca, tidak dibuat** (`MODUL.md`), kini sebagai cadangan dan sumber alat pindah.

**P1 (`RALAT-DEV-30-09-2026.md`)** — modul ini **tidak membuat satu tabel, sequence, maupun constraint pun**. Produk
tetap JSON di dua tabel lama seperti Pega (OQ-MPNL-01, bawaan "JSON seperti Pega"); lampiran direkam di tabel lama
(P5). Rentang migrasi `140`–`179` sengaja kosong (`TestMPNLNolMigrasiDiRentang`); satu-satunya migrasi modul ini
adalah slot menu `960` (`UPDATE DIMIGRASI`).

Berkas ini **peta**, bukan DDL: tabel → kolom VERBATIM → tipe katalog → penulis Pega → cara kolom itu ditulis modul
ini. Penjaga yang memakainya: `TestKolomDDLCocokDenganStruktur` / `TestTabelBukanMilikKitaTidakDibuat`
(`inti/backend/penjaga/strukturkolom_test.go`) — ketiga tabel di bawah dinyatakan "Tabel warisan" di `MODUL.md`, jadi
migrasi mana pun yang **membuatnya** merah.

Sumber tipe: `docs/dba-procedures-and-ddl.md` `[data DBA]`. ⛔ Procedure **tidak dipanggil**
(`PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE`): isinya ditiru di Go — upsert dikunci `ID`, ID baru
`'1' || LPAD(M_PRODUCT_LIFE_SEQ.NEXTVAL, 5, '0')`, kedua tabel di SATU transaksi (P4), nol `COMMIT` di teks SQL.

## M_PRODUCT_LIFE

Sisi umum produk — halaman Pega `ProductName`. Penulis Pega `SaveProductName_Act` 9 b1833 (`SaveProductNameLIfe` →
`PEGA_M_PRODUCT_LIFE`) dan 10 b2021 (`SaveProductNameLIfeFlat`). Dibaca view `PRODUCT_LIFE` dan `DOCUMENTCLAIM_LIFE`.

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(6) | `'1' || LPAD(M_PRODUCT_LIFE_SEQ.NEXTVAL, 5, '0')` — dari server, tidak pernah dari klien; > 5 digit atau ID terpakai di salah satu tabel = gagal terang (P6) |
| `JSONDATA` | CLOB, `IS JSON` | halaman `ProductName` berkunci Pega (`POLICYHODER`, `UnderwritingLimitList`, `OutwardList`, …); kunci yang dibaca view dijamin ada; kunci lama yang tidak dikelola layar dipertahankan |
| `RIRISKID` | VARCHAR2(10) | = JSON `RIRISKID` (`SaveProductNameLIfeFlat` b84) |
| `RIRISK` | VARCHAR2(100) | = JSON `RIRISK` |

⭐ **Ralat 01-10-2026 (lanjutan 1 L1):** dua baris lama dikutip — *"`PRODUCTNAME` | VARCHAR2(1000) | = JSON `PRODUCTNAME` (OQ-MPNL-08)"*,
*"`BEGIN_DATE` | DATE | = `BEGIN` inward `dd/MM/yyyy` → `DATE`; NULL bila kosong (OQ-MPNL-08)"* — dicabut: katalog DEV `ALL_TAB_COLUMNS`
`M_PRODUCT_LIFE` hanya empat kolom di atas (`testdata/katalog-dev.json`, uji `TestKolomDitulisAdaDiKatalogDEV`).

## M_PRODUCTINWARD_LIFE

Sisi inward — halaman Pega `ProductNameInward`. Penulis Pega `SaveProductName_Act` 15 b2864
(`SaveProductNameInwardLIfe` → `PEGA_M_PRODUCT_INWARD_LIFE`). Dibaca view `PRODUCTINWARD_LIFE` (Claim Life).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(6) | = ID produk untuk baris baru (R14, OQ-MPNL-02); baris lama dipertahankan ID-nya |
| `JSONDATA` | CLOB, `IS JSON` | halaman `ProductNameInward`; `PRODUCTID` = ID produk; tanggal `dd/MM/yyyy` |

## M_ATTACHMENTPRODUCTNAME

Rekam lampiran produk. Penulis Pega `InsertAttachProdName_Sql` b84, penghapus `DeleteAttachProdName_Sql` b85, pembaca
`GetAttachmentProdName_Sql` b84. Tipe kolom belum ada di katalog modul ini `[data DBA]`.

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | — | `TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3')` (dicoba ulang bila bentrok) |
| `TREATYID` | — | ID produk (`ProductName.ID`) |
| `CATEGORY` | — | `File` (`ProductNameSaveAttachment` 2.1 b475) |
| `FILENAME` | — | nama berkas asli |
| `FILEMIMETYPE` | — | ekstensi huruf kecil (`InsertGoogleStorage_Act` 3 b566) |
| `DATA_JSON` | — | NULL (`""` di Pega) |
| `USERNAME` | — | akun pelaku |
| `T_STORAGE_ID` | — | `ImageID` objek penyimpanan |

## Tabel bersama dan master yang disentuh (digambarkan dokumen pemiliknya)

Tabel di bawah TIDAK digambarkan di sini — pemiliknya yang menggambarkannya, dan dua gambaran yang berbeda akan
saling membatalkan (`TestDokumenSTRUKTURSepakatAtasTabelBersama`).

| Tabel | Sentuhan modul ini |
| --- | --- |
| `T_STORAGE_IMAGE` | ditulis pelaksana stub lampiran (`Insert_T_Storage_SQL` b85: `IMAGEID`, `URLPUBLIC` NULL, `APPFOLDER`, `EXPDATE`, `FILENAME`, `APPNAME`, `STORAGE = 'standard'`), dihapus `DeleteStorage_SQL` b85 |
| `T_LOG_SERVICE_RNM` | outbox bersama (Claim Life 015): `MODUL = 'MASTERPRODUCTNAMELIFE'`, `JENIS_EFEK = 'unggah-lampiran'` |
| `T_FOLDER_IMAGE` | dibaca `APPNAME` (`GetAppName_SQL` b58) |
| `AGENT`, `CLIENT`, `CURRENCY`, `RIRISK_LIFE_SUMMARY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE` | dibaca saja — pemilih master dan `PLAN LIST` (OQ-MPNL-04) |
| `TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE` | dibaca saja — `On Retention` → `BrowseReinstypeOR_SQL` b84 (milik Master Contract Retro Life) |

---

# Tabel FLAT produk — dibuat migrasi 140–147 *(keputusan work owner 02-10-2026, tiket 01 bab bertanggal)*

Kolom dan tipe di bawah dibaca penjaga inti `TestKolomDDLCocokDenganStruktur` / `TestGolonganTipeDDLCocokDenganStruktur`
(kolom **Tipe**: golongan; kolom **DDL**: tipe fisik). Seluruh kolom NULLABLE kecuali kunci. Sumber = kunci halaman
Pega lama yang dipindah alat `pindahflat` (tiket 01, T2) dan medan form yang menulisnya (`PARITAS-LAYAR-DAN-AKSI.md` §3–§5).

## M_PRODUCTNAME_LIFE

Induk — satu baris = satu produk (sisi umum + sisi inward, grilling Q1b). PK `ID`; `CHECK (IS_ORS IN (0, 1))`.

| Kolom | Tipe | DDL | Sumber / isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(6) NOT NULL | `'1' || LPAD(M_PRODUCT_LIFE_SEQ.NEXTVAL, 5, '0')` dari server (P6); PK |
| `PRODUCTNAME` | teks | VARCHAR2(1000) | `PRODUCTNAME` - `Product Name` b3620 |
| `CEDINGID` | teks | VARCHAR2(100) | `CEDINGID` - `setCeding_DT` b191 (`AGENT.ID`) |
| `CEDING` | teks | VARCHAR2(1000) | `CEDING` - `Ceding` b4075 |
| `SOBID` | teks | VARCHAR2(100) | `SOBID` - `setSOB_DT` b191 (`AGENT.ID`) |
| `SOBNAME` | teks | VARCHAR2(1000) | `SOBNAME` - `SOB` b4463 |
| `RICOMM` | angka desimal | NUMBER(38,8) | `RICOMM` - `Deduction (%)` b7104 |
| `RIRISKID` | teks | VARCHAR2(10) | `RIRISKID` - `setRIRISK_DT` b191 (`RIRISK_LIFE_SUMMARY.ID` VARCHAR2(10)) |
| `RIRISK` | teks | VARCHAR2(100) | `RIRISK` - `R/I Risk Name` b7398 |
| `RIRATEID` | teks | VARCHAR2(10) | `RIRATEID` - milik server (medan layar mati; DEV 158 terisi) |
| `RIRATE` | teks | VARCHAR2(500) | `RIRATE` - milik server (DEV 153 terisi, maks 99) |
| `INWARDNAME` | teks | VARCHAR2(1000) | `INWARDNAME` - `Treaty Name` b8719 |
| `TREATYNUMBER` | teks | VARCHAR2(100) | `TREATYNUMBER` - `Treaty Number` b8900 |
| `CAUSEID` | teks | VARCHAR2(10) | `CAUSEID` - `setCauseOfLoss_DT` b172 (`CAUSEOFLOSS_LIFE.ID` VARCHAR2(10)) |
| `CAUSE` | teks | VARCHAR2(100) | `CAUSE` - `Cause Of Loss` b10661 |
| `IS_ORS` | bilangan bulat | NUMBER(5) | `IsORS` - checkbox `On Retention` b47312 (R9); 0/1, `CHECK` |
| `POLICYHOLDER` | teks | VARCHAR2(100) | `POLICYHODER` (ejaan Pega) - `setPolicyHolder_DT` b172 (`CLIENT.ID`); satu kolom untuk kedua sisi (sama di 196 produk DEV) |
| `POLICYHOLDERNAME` | teks | VARCHAR2(1000) | `POLICYHODERNAME` - `Policy Holder` b17097 |
| `INSURED` | teks | VARCHAR2(4000) | `INSURED` - `Insured` b18341 |
| `ADDENDUMNO` | bilangan bulat | NUMBER(5) | `ADDENDUMNO` - `Addendum No.` b18834 |
| `ADDENDUMWORD` | teks | VARCHAR2(4000) | `ADDENDUMWORD` - `Addendum` b19432 |
| `AMANDEMENTNO` | bilangan bulat | NUMBER(5) | `AMANDEMENTNO` - `Amandement No.` b20155 |
| `AMANDEMENTSCHD` | teks | VARCHAR2(4000) | `AMANDEMENTSCHD` - `Amandement` b20760 |
| `MAXEXPIREDCLAIM` | bilangan bulat | NUMBER(5) | `MAXEXPIREDCLAIM` - `Max Notification Claim Expired` b21781 |
| `BEGIN_DATE` | date | DATE | `BEGIN` (`dd/MM/yyyy`) - `Begin Date` b21969; `BEGIN` kata terlarang Oracle |
| `STNC` | date | DATE | `STNC` - `STNC` b22304 |
| `MATURE` | date | DATE | `MATURE` - `Expired Date` b27250; nilai lama bukan tanggal = NULL (K3) |
| `CEDINGRETENTIONNUM` | angka desimal | NUMBER(38,8) | `CEDINGRETENTIONNUM` - `Ceding Retention (%)` b22494 |
| `CEDINGLIMIT` | angka desimal | NUMBER(38,8) | `CEDINGLIMIT` - `Ceding's Limit` b22657 |
| `BROKERAGE` | angka desimal | NUMBER(38,8) | `BROKERAGE` - `Brokerage Fee (%)` b22915 |
| `MINAGE` | bilangan bulat | NUMBER(5) | `MINAGE` - `Minimum Age (Years)` b23078 |
| `MAXAGE` | bilangan bulat | NUMBER(5) | `MAXAGE` - `Maximum Age (Years)` b23265 |
| `EXPIRYAGE` | bilangan bulat | NUMBER(5) | `EXPIRYAGE` - `Expiry Age (Years)` b23452 |
| `EXTRAPREMI` | angka desimal | NUMBER(38,8) | `EXTRAPREMI` - `Extra Premium` b23635 |
| `MINSUMINSURED` | angka desimal | NUMBER(38,8) | `MINSUMINSURED` - `Min Sum Insured` b23796 |
| `MAXSUMINSURED` | angka desimal | NUMBER(38,8) | `MAXSUMINSURED` - `Max Sum Insured` b23956 |
| `MAXSUMREASURED` | angka desimal | NUMBER(38,8) | `MAXSUMREASURED` - `Max Sum Reasured` b24214 |
| `RNMSHARE` | angka desimal | NUMBER(38,8) | `RNMSHARE` - `Nusantara Re Share (%)` b24674 |
| `RNMLIMITNUM` | angka desimal | NUMBER(38,8) | `RNMLIMITNUM` - `Nusantara Re's Limit` b25236 |
| `PREMIUMFACTOR` | angka desimal | NUMBER(38,8) | `PREMIUMFACTOR` - `Premium Factor (%)` b25398 |
| `PAYMENT` | teks | VARCHAR2(1) | `PAYMENT` - `Premium Payment Method` b25611 (kode 1-5, OQ-MPNL-05) |
| `SUBJECTTO` | teks | VARCHAR2(4000) | `SUBJECTTO` - `Subject To` b25924 (teks) |
| `ANNUITYINTEREST` | angka desimal | NUMBER(38,8) | `AnnuityInterest` - `Annuity Interest (%)` b26092 |
| `PREMIUMREFUNDFACTOR` | angka desimal | NUMBER(38,8) | `PremiumRefundFactor` - `Premium Refund Factor (%)` b26306 |
| `MAXDATARECEIVE` | bilangan bulat | NUMBER(5) | `MAXDATARECEIVE` - `Max Production Data Receive` b27062 |
| `BIRTHDAY` | teks | VARCHAR2(1) | `BIRTHDAY` - `Birthday` b27960 (bukan tanggal; DEV panjang 1, OQ-MPNL-05) |
| `CURRENCYID` | teks | VARCHAR2(100) | `CURRENCYID` - `setCurrency_DT` b172 (`CURRENCY.ID`) |
| `CURRENCY` | teks | VARCHAR2(20) | `CURRENCY` - `Currency` b28140 |
| `EXTRAMORTALITY` | angka desimal | NUMBER(38,8) | `EXTRAMORTALITY` - `Extra Mortality (%)` b29256 |
| `MAXCONTRACT` | bilangan bulat | NUMBER(5) | `MAXCONTRACT` - `Max Contract (year)` b29442 |
| `PROPORTIONALTABLE` | angka desimal | NUMBER(38,8) | `PROPORTIONALTABLE` - `Proportional Table` b29626 |
| `CREATEOP` | teks | VARCHAR2(100) | `CREATEOP` - `SaveProductName_Act` 6 b1370 (pelaku produk baru) |
| `UPDATEOP` | teks | VARCHAR2(100) | `UPDATEOP` - `SaveProductName_Act` 1 b577 (pelaku) |

## M_PRODUCTNAME_LIFE_LIEN

Anak (migrasi 141) — `LienClause[*]` - grid b12201 (RALAT R10). PK (`PRODUCTID`, `URUT`); FK `PRODUCTID` → `M_PRODUCTNAME_LIFE(ID)` `ON DELETE CASCADE`.

| Kolom | Tipe | DDL | Sumber / isi |
| --- | --- | --- | --- |
| `PRODUCTID` | teks | VARCHAR2(6) NOT NULL | ID produk induk; FK, `NOT NULL` |
| `URUT` | bilangan bulat | NUMBER(5) NOT NULL | urutan baris di grid, mulai 1; `NOT NULL` |
| `USIA` | teks | VARCHAR2(200) | `Usia` - `Usia saat Klaim` b12741 (teks) |
| `MANFAAT` | teks | VARCHAR2(200) | `Manfaat` - `% Manfaat yang dibayarkan` b12890 (teks) |

## M_PRODUCTNAME_LIFE_DOCCLAIM

Anak (migrasi 142) — `DocumentClaim[*]` - grid b14601. PK (`PRODUCTID`, `URUT`); FK `PRODUCTID` → `M_PRODUCTNAME_LIFE(ID)` `ON DELETE CASCADE`.

| Kolom | Tipe | DDL | Sumber / isi |
| --- | --- | --- | --- |
| `PRODUCTID` | teks | VARCHAR2(6) NOT NULL | ID produk induk; FK, `NOT NULL` |
| `URUT` | bilangan bulat | NUMBER(5) NOT NULL | urutan baris di grid, mulai 1; `NOT NULL` |
| `DOCUMENT` | teks | VARCHAR2(500) | `Document` - `Document List` b15125 |

## M_PRODUCTNAME_LIFE_PLAN

Anak (migrasi 143) — `PlanList[*]` - grid b31557 (`OUTWARDRATEID` 0 terisi - tidak diambil). PK (`PRODUCTID`, `URUT`); FK `PRODUCTID` → `M_PRODUCTNAME_LIFE(ID)` `ON DELETE CASCADE`.

| Kolom | Tipe | DDL | Sumber / isi |
| --- | --- | --- | --- |
| `PRODUCTID` | teks | VARCHAR2(6) NOT NULL | ID produk induk; FK, `NOT NULL` |
| `URUT` | bilangan bulat | NUMBER(5) NOT NULL | urutan baris di grid, mulai 1; `NOT NULL` |
| `PLANID` | teks | VARCHAR2(10) | `PlanID` - autocomplete `.ID` (`PRODUCT_TYPE_LIFE.ID`) |
| `PLAN` | teks | VARCHAR2(200) | `Plan` - `Plan Name` b31845 |
| `NAME` | teks | VARCHAR2(200) | `Name` - `Bussines` b31994 |
| `BENEFIT` | teks | VARCHAR2(500) | `Benefit` - `Benefit` b32143 |
| `RIRATEID` | teks | VARCHAR2(10) | `RIRATEID` - `SetRIRate` 1 b249 (`RATE_LIFE_SUMMARY.ID` VARCHAR2(10)) |
| `RIRATE` | teks | VARCHAR2(500) | `RIRATE` - `R/I Rate` b32296 |

## M_PRODUCTNAME_LIFE_FINUW

Anak (migrasi 144) — `FinancialUnderwritingList[*]` - grid b37148. PK (`PRODUCTID`, `URUT`); FK `PRODUCTID` → `M_PRODUCTNAME_LIFE(ID)` `ON DELETE CASCADE`.

| Kolom | Tipe | DDL | Sumber / isi |
| --- | --- | --- | --- |
| `PRODUCTID` | teks | VARCHAR2(6) NOT NULL | ID produk induk; FK, `NOT NULL` |
| `URUT` | bilangan bulat | NUMBER(5) NOT NULL | urutan baris di grid, mulai 1; `NOT NULL` |
| `MININSURED` | angka desimal | NUMBER(38,8) | `MinInsured` - `Min Insured` b37670 |
| `MAXINSURED` | angka desimal | NUMBER(38,8) | `MaxInsured` - `Max Insured` b37818 |
| `EMPLOYEE` | teks | VARCHAR2(1000) | `Employee` - `Employee` b37966 (teks) |
| `NON_EMPLOYEE` | teks | VARCHAR2(1000) | `Non_Employee` - `Non-Employee` b38115 (teks) |

## M_PRODUCTNAME_LIFE_UWLIMIT

Anak (migrasi 145) — `UnderwritingLimitList[*]` - grid b42075. PK (`PRODUCTID`, `URUT`); FK `PRODUCTID` → `M_PRODUCTNAME_LIFE(ID)` `ON DELETE CASCADE`.

| Kolom | Tipe | DDL | Sumber / isi |
| --- | --- | --- | --- |
| `PRODUCTID` | teks | VARCHAR2(6) NOT NULL | ID produk induk; FK, `NOT NULL` |
| `URUT` | bilangan bulat | NUMBER(5) NOT NULL | urutan baris di grid, mulai 1; `NOT NULL` |
| `MININSURED` | angka desimal | NUMBER(38,8) | `MinInsured` - `Min Insured` b42594 |
| `MAXINSURED` | angka desimal | NUMBER(38,8) | `MaxInsured` - `Max Insured` b42742; nilai lama bukan angka = NULL (K3) |
| `MINAGE` | bilangan bulat | NUMBER(5) | `MinAge` - `Min Age` b42890 |
| `MAXAGE` | bilangan bulat | NUMBER(5) | `MaxAge` - `Max Age` b43038 |
| `MEDICAL` | teks | VARCHAR2(200) | `Medical` - `Medical` b43187 (teks bebas, R16) |
| `DESCRIPTION` | teks | VARCHAR2(1000) | `Description` - `Description` b43336 |

## M_PRODUCTNAME_LIFE_OUTWARD

Anak (migrasi 146) — `OutwardList[*]` berisi - `GetReinsTypeOR_Life` 4.1 b770 (objek kosong tidak dipindah, K4). PK (`PRODUCTID`, `URUT`); FK `PRODUCTID` → `M_PRODUCTNAME_LIFE(ID)` `ON DELETE CASCADE`.

| Kolom | Tipe | DDL | Sumber / isi |
| --- | --- | --- | --- |
| `PRODUCTID` | teks | VARCHAR2(6) NOT NULL | ID produk induk; FK, `NOT NULL` |
| `URUT` | bilangan bulat | NUMBER(5) NOT NULL | urutan baris di grid, mulai 1; `NOT NULL` |
| `REINSTYPEID` | teks | VARCHAR2(100) | `REINSTYPEID` |
| `REINSTYPENAME` | teks | VARCHAR2(100) | `REINSTYPENAME` |
| `TRANSACTIONYEAR` | bilangan bulat | NUMBER(5) | `TRANSACTIONYEAR` (= `TREATYYEAR`) |
| `TREATYCONTRACTID` | teks | VARCHAR2(100) | `TREATYCONTRACTID` |
| `UNDERWRITINGYEAR` | bilangan bulat | NUMBER(5) | `UNDERWRITINGYEAR` |
| `OVR_COMM` | angka desimal | NUMBER(38,8) | `OVR_COMM` - tanpa penulis di korpus, kosong (OQ-MPNL-09) |

## M_PRODUCTNAME_LIFE_COMMENT

Anak (migrasi 147) — `CommentList[*]` - `AddCommentList_Act` 1 b235 (`IsApproved` 0 terisi - tidak diambil). PK (`PRODUCTID`, `URUT`); FK `PRODUCTID` → `M_PRODUCTNAME_LIFE(ID)` `ON DELETE CASCADE`.

| Kolom | Tipe | DDL | Sumber / isi |
| --- | --- | --- | --- |
| `PRODUCTID` | teks | VARCHAR2(6) NOT NULL | ID produk induk; FK, `NOT NULL` |
| `URUT` | bilangan bulat | NUMBER(5) NOT NULL | urutan baris di grid, mulai 1; `NOT NULL` |
| `TANGGAL` | timestamp | TIMESTAMP | `Date` - stempel Pega `YYYYMMDDTHHMMSS.SSS GMT` (`@CurrentDateTime()`), disimpan GMT |
| `OPERATORNAME` | teks | VARCHAR2(100) | `OperatorName` - `PIC` b62246 |
| `SUGGEST` | teks | VARCHAR2(4000) | `Suggest` - `Comment` b62399 |
