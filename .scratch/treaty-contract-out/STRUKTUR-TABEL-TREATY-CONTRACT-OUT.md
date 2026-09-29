# Struktur Tabel — Treaty Contract Out: PETA TABEL WARISAN yang dipakai

**Keputusan tco4 `[DIPUTUSKAN work owner 29-09-2026]`** — menggantikan tco1. Kutipan: *"khusus modul treaty contract out
tidak ada tabel baru sama sekali!!!"*. Modul ini **tidak membuat satu tabel atau sequence pun**: ia menulis dan membaca
tabel yang **sudah ada** di `POOLDATA`, dengan **nama tabel dan kolom VERBATIM**, persis seperti RDB XML. Migrasi
`300`–`307` (tco1: `T_TREATYYEAR` … `T_TREATYYEAR_LAMPIRAN`) **dibuang** 29-09-2026.

Berkas ini **peta**, bukan DDL: tabel → kolom VERBATIM → tipe katalog → RDB penulis/pembaca → bentuk nilai yang ditulis.
Penjaga yang memakainya:

- `TestKolomDDLCocokDenganStruktur` / `TestTabelBukanMilikKitaTidakDibuat` (`strukturkolom_test.go`): setiap tabel di sini
  terdaftar `tabelBukanMilikKita` — migrasi mana pun yang **membuatnya** gagal;
- `TestTCONolTabelBaru` (`tco_warisan_test.go`): nol berkas migrasi di rentang `300`–`319`;
- `TestKolomWarisanTCOSesuaiProcedure`: daftar kolom `KolomWarisanTCO` = urutan parameter procedure penulisnya.

Sumber tipe: `dba-procedures.md` bab "DDL + 4 procedure master lain" `[data DBA]`. Kolom yang DBA tidak sebut tipenya
tertulis **VARCHAR2 `[tidak disebut DBA]`** dan diperlakukan teks.

⛔ Procedure **tidak dipanggil** (keputusan **o**): isinya ditiru di Go — UPSERT dikunci `ID`, ID baru `'1' || lpad(seq, n)`
dengan sequence **warisan**, nol `COMMIT` di teks SQL (ADR-U-0029; `DalamTransaksi` yang commit).

⛔ Nol tabel jejak modul: `T_TREATYCO_JEJAK` dibuang, Pega tidak mencatat jejak modul ini. Efek keluar lampiran memakai
outbox bersama `T_LOG_SERVICE_RNM` milik aplikasi (bukan tabel modul ini).

Rule korpus (folder `D:\XML\RNM_BRD\Treaty Contract Out\`), dibaca 29-09-2026; `pyStepsBlockName` RDB: tidak ada
(Rule-Connect-SQL tanpa langkah):

| Singkatan | Rule | Isi |
| --- | --- | --- |
| **SaveYear** | `RDBList/SaveMasterTreatyYear_SQL.xml` b84 | `POOLDATA.PEGA_TREATYYEAR` (10 param + 2 out) |
| **SaveContract** | `RDBList/SaveMasterTreatyContract_SQL.xml` b84 | `POOLDATA.PEGA_TREATYCONTRACT` (8 + 2 out) |
| **SaveReins** | `RDBList/SaveMasterTreatyReinsurer_SQL.xml` b79 | `POOLDATA.PEGA_TREATYREINSURER` (19 + 2 out) |
| **SaveBiz** | `RDBList/SaveMasterTreatyBusiness_SQL.xml` b85 | `POOLDATA.PEGA_TREATYBUSINESS` (12 + 2 out) |
| **SaveArrg** | `RDBList/SaveMasterProportionalArrg.xml` b85 | `POOLDATA.PEGA_PROPORTIONALARRG` (35 + 2 out) |
| **SaveArrgChild** | `RDBList/SaveMasterProportionalArrgChild.xml` b85 | `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD` (26 + 2 out) — tabel yang SAMA `[data DBA]` |
| **InsSec** | `RDBList/InsertToMTreatySecurity.xml` b60 | `insert into mtreatysecurity values (…)` posisional 7 nilai |
| **UpdSec** | `RDBList/UpdateMTreatySecurity.xml` b85 | `where REAS_ID = … and trim(REAS_SECURITY) = trim(…)` |
| **DelSec** | `RDBList/DeleteSecurityReinsurer.xml` b85 | idem kunci |
| **DelReins** | `RDBList/DeleteFromTreatyReinsurer_Act.xml` b59 | security `REAS_ID` lalu reinsurer `ID` |
| **DelBiz** | `RDBList/DeleteRowBusinessList.xml` b85 | `delete from treatybusiness where id` *(tab Delete `m_treatybusiness` b88 tidak dipakai `RDB-List`)* |
| **DelContract** | `RDBList/DeleteFromTREATYCONTRACT_SQL.xml` b79 | empat DELETE: kontrak, business, security, reinsurer |
| **InsAtt** | `RDBList/InsertAtatchment_Sql.xml` b60 | `POOLDATA.PEGA_M_ATTACHMENT(IDPEGA, DATAPEGA CLOB)` — badan **`[terbuka — DBA]`** |
| **GetAtt** | `RDBList/GetAllAttachment2_Sql.xml` b84, `GetAttachment2_Sql.xml` b85 | `M_ATTACHMENTTREATY_2 where treatyid = {TreatyIn.ID}` |
| **DelAtt** | `RDBList/DeleteAttachment2_Sql.xml` b84 | `delete M_ATTACHMENTTREATY_2 where treatyid = … and id = …` |
| **Storage** | `RDBList/GetLinkStorage_SQL.xml`, `Update_T_Storage_SQL.xml`, `DeleteStorage_SQL.xml`, `GetTokenStorage_SQL.xml` | `T_STORAGE_IMAGE`, `GET_TOKEN_STORAGE` |
| **Saudara** | `Treaty In/RDBList/InsertAttachment2_Sql.xml` b84; `Claim Fac In/RDBList/Insert_T_Storage_SQL.xml` b85, `GenerateImageID_SQL.xml` b79, `GetAppName_SQL.xml` b59 | penulis tabel warisan YANG SAMA di modul saudara — bukti kolom |

**Bentuk nilai yang ditulis** (kolom teks warisan yang memuat tanggal/angka):

| Bentuk | Kolom | Bukti |
| --- | --- | --- |
| stempel Pega `YYYYMMDDTHHMMSS.SSS GMT` (UTC) | `TREATYYEAR.TGLUPDATE`, `TREATYCONTRACT.TGLUPDATE` | `@getCurrentTimeStamp()` — `SaveTreatyYear_Act` b328, `SaveTreatyContract_Act` b1458 `[terverifikasi]` |
| **`YYYYMMDD`** (delapan angka, tanpa jam/zona) | `TREATYYEAR.STARTDATE`, `TREATYYEAR.ENDDATE` | data DEV 182/182 baris (brief lanjutan 4) — **OQ-TCO-01 ditutup**; dugaan kuat "stempel 00:00 WIB" lanjutan 3 dibantah data. Dibaca HANYA bentuk ini; bentuk lain → galat berkata-kata (`tanggalTahunWarisanTeks`) |
| DATE lewat `to_date(…,'DD/MM/YYYY')` | `TREATYCONTRACT.TREATYSTARTDATE/ENDDATE` | `dba-procedures.md`; `SaveTreatyContract_Act` b1479 `[terverifikasi]` |
| desimal titik tanpa pemisah ribuan | `PROPORTIONALARRG.RP/USD/PCT/PCTME`, `MTREATYSECURITY.PCT_SHARE` | hasil `@toDecimal` (`HitungRpUsd_depan` b293–b294, b382–b383); data DEV `RP` 930 bertitik / 0 berkoma, `PCT` 520 / 0 — **OQ-TCO-23 ditutup: titik**. Dibaca: titik ATAU koma (`UraiDesimalWarisanTCO`) |
| waktu simpan (DATE) | `PROPORTIONALARRG.TGLUPDATE` | procedure memakai `SYSDATE` `[data DBA]`; layanan mengikat jam simpannya (setara) |
| NULL (Pega tidak mengisinya) | `TREATYREINSURER.IUDATE/STATUSON/STDRATING` (dari baris lama bila diubah) | `NewTreatyReinsurerDetail_Act` b917–b1086 mengosongkan `[terverifikasi]` |
| **selalu NULL** (sisip dan ubah) | `TREATYREINSURER.STARTDATE/ENDDATE` | kontrol `ViewDetailTreatyReinsurerGrid1.xml` b10311/b10516 bersyarat `1=2` (b10431/b10636); `NewTreatyReinsurerDetail_Act` b938/b959; data DEV 430/430 kosong — **OQ-TCO-01** (lanjutan 4) |
| **NULL seperti Pega** (`USERID` reinsurer lama dipertahankan saat ubah) | `TREATYREINSURER.USERID/TGLUPDATE`, `TREATYBUSINESS.USERID/TGLUPDATE` | data DEV 0/430 dan 2/4.621 terisi — **OQ-TCO-25 ditutup** (lanjutan 4); pelaku di log aplikasi |

---

## TREATYYEAR

Tahun treaty. `[data DBA]` seluruh kolom VARCHAR2. Penulis **SaveYear**; ID `'1' || lpad(TreatyYear_seq.nextval, 6, '0')`.

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2 | param 1; UPSERT kunci |
| `TREATYYEAR` | VARCHAR2 | param 2 — kode tahun (ADR-U-0022) |
| `UNDERWRITINGYEAR` | VARCHAR2 | param 3 |
| `TREATYGROUPID` | VARCHAR2 | param 4 |
| `TREATYGROUPNAME` | VARCHAR2 | param 5 |
| `USERID` | VARCHAR2 | param 6 — `OperatorID.pyUserName` (b281) |
| `TGLUPDATE` | VARCHAR2 | param 7 — stempel Pega |
| `PROPORTION` | VARCHAR2 | param 8 — `[terbuka]` arti |
| `STARTDATE` | VARCHAR2 | param 9 — `YYYYMMDD` (OQ-TCO-01 ditutup, lanjutan 4) |
| `ENDDATE` | VARCHAR2 | param 10 — idem |

## TREATYCONTRACT

Kontrak (jenis reasuransi) di satu tahun. Penulis **SaveContract**; ID `'1' || lpad(treatycontract_seq.nextval, 6, '0')`.
Penghapus **DelContract** (langkah 1).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2 `[tidak disebut DBA]` | param 1 |
| `IDTREATYYEAR` | VARCHAR2 `[tidak disebut DBA]` | param 2 → `TREATYYEAR.ID` (tanpa FK) |
| `REINSTYPEID` | VARCHAR2 `[tidak disebut DBA]` | param 3 — master `REINSURANCETYPE` |
| `REINSTYPENAME` | VARCHAR2 `[tidak disebut DBA]` | param 4 |
| `TREATYSTARTDATE` | DATE | param 5 — `to_date(…,'DD/MM/YYYY')` |
| `TREATYENDDATE` | DATE | param 6 — idem |
| `USERID` | VARCHAR2 `[tidak disebut DBA]` | param 7 — `OperatorID.pyUserName` |
| `TGLUPDATE` | VARCHAR2(1000) | param 8 — stempel Pega |

## TREATYREINSURER

Reinsurer pada kombinasi **(TREATYYEAR, TREATYGROUPID, REINSTYPEID)** — bukan FK ke kontrak. Penulis **SaveReins**; ID
`'1' || lpad(M_TREATYREINSURER_SEQ.nextval, 6, '0')`. Pembaca `GetMasterReinsurerList` b85 dan hilir `GetListRetro_Sql`.

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2 | param 1 |
| `TREATYYEAR` | VARCHAR2 | param 2 — kunci gabungan |
| `TREATYGROUPID` | VARCHAR2 | param 3 — kunci gabungan |
| `TREATYGROUPNAME` | VARCHAR2 | param 4 |
| `REINSTYPEID` | VARCHAR2 | param 5 — kunci gabungan |
| `REINSTYPENAME` | VARCHAR2 | param 6 |
| `REINSURERID` | VARCHAR2 | param 7 |
| `CLIENTID` | VARCHAR2 | param 8 |
| `NAME` | VARCHAR2 | param 9 |
| `RICOMM` | NUMBER | param 10 — bind angka |
| `PCTSHARE` | NUMBER | param 11 — bind angka |
| `IUDATE` | VARCHAR2 | param 12 — NULL (Pega mengosongkan) |
| `USERID` | VARCHAR2 | param 13 — NULL untuk baris baru, nilai baris lama saat ubah (OQ-TCO-25, lanjutan 4) |
| `STARTDATE` | VARCHAR2 | param 14 — selalu NULL (OQ-TCO-01, lanjutan 4) |
| `ENDDATE` | VARCHAR2 | param 15 — idem |
| `STATUSON` | VARCHAR2 | param 16 — NULL |
| `STDRATING` | VARCHAR2 | param 17 — NULL |
| `OPERATORNAME` | VARCHAR2 | param 18 — `OperatorID.pyUserName` (`SaveTreatyReinsurerDetail1_Act` b378) |
| `TGLUPDATE` | VARCHAR2 | param 19 — NULL seperti Pega (OQ-TCO-25, lanjutan 4) |

## MTREATYSECURITY

Security di bawah seorang reinsurer. `[data DBA]` **tanpa PK**. Kunci baris = **(`REAS_ID`, `trim(REAS_SECURITY)`)** —
UpdSec b85 dan DelSec b85; OQ-TCO-17 (security dobel ditolak) menjamin keunikannya. Penulis **InsSec** (posisional; Go
menulis daftar kolom bernama dengan nilai yang sama).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `THN_TREATY` | VARCHAR2(4) DEFAULT '1' NOT NULL | posisi 1 — `InputData.CARI8` = `.TreatyYear` baris reinsurer |
| `TOP_ID` | VARCHAR2(9) | posisi 2 — `''` (NULL) |
| `TP_TREATY` | CHAR(2) | posisi 3 — `''` (NULL) |
| `REAS_ID` | CHAR(7) NOT NULL | posisi 4 — `TREATYREINSURER.ID` |
| `PCT_SHARE` | VARCHAR2(99) | posisi 5 — desimal teks (OQ-TCO-23) |
| `USER_ID` | CHAR(99) | posisi 6 — `''` (NULL) |
| `REAS_SECURITY` | CHAR(10) NOT NULL | posisi 7 — nama security; CHAR berekor spasi → dibandingkan `trim` |

## TREATYBUSINESS

Bisnis pada kombinasi tahun + grup + jenis. `[data DBA]` seluruh kolom VARCHAR2. Penulis **SaveBiz**; ID
`'1' || lpad(TREATY_BUSINESS_SEQ.nextval, 6, '0')`. ⚠️ UPDATE procedure hanya men-set `ISACTIVE, BIZCODE, BIZNAME, USERID,
TGLUPDATE` `[data DBA]`. Penghapus **DelBiz**, **DelContract** (langkah 2: `TREATYYEARID = … OR TREATYYEARID IS NULL`).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2 | param 1 |
| `ISACTIVE` | VARCHAR2 | param 2 — `"0"` nonaktif (OQ-TCO-13) |
| `TREATYYEAR` | VARCHAR2 | param 3 |
| `TREATYYEARID` | VARCHAR2 | param 4 — `InputTreatyContract.IDTreatyYear` (`SaveTreatyBusinessDetail_Act` b355) |
| `TREATYGROUPID` | VARCHAR2 | param 5 |
| `TREATYGROUPNAME` | VARCHAR2 | param 6 |
| `REINSTYPEID` | VARCHAR2 | param 7 |
| `REINSTYPENAME` | VARCHAR2 | param 8 |
| `BIZCODE` | VARCHAR2 | param 9 |
| `BIZNAME` | VARCHAR2 | param 10 |
| `USERID` | VARCHAR2 | param 11 — NULL seperti Pega (OQ-TCO-25, lanjutan 4) |
| `TGLUPDATE` | VARCHAR2 | param 12 — NULL seperti Pega (OQ-TCO-25, lanjutan 4) |

## PROPORTIONALARRG

Klausul arrangement — induk dan "anak" di tabel yang SAMA `[data DBA]`. Penulis **SaveArrg** (35) / **SaveArrgChild** (26,
tanpa sembilan kolom terakhir); ID `'1' || lpad(PROPORTIONALARRG_SEQ.nextval, 7, '0')`. Pembaca hidup relasional:
`GetMasterDescriptionLimitParentList` b85 dan hilir. RDB `GetMasterDescription*`/`GetMasterPanggilID`/`GetMasterPortfolioListDetail`
membaca JSON `m_PROPORTIONALARRG` — **mati** `[keputusan work owner]`, tidak ditulis procedure mana pun; tidak dibaca.

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2 `[tidak disebut DBA]` | param 1 |
| `TREATYYEAR` | VARCHAR2 `[tidak disebut DBA]` | param 2 |
| `TREATYYEARID` | VARCHAR2 `[tidak disebut DBA]` | param 3 |
| `TREATYGROUPID` | VARCHAR2 `[tidak disebut DBA]` | param 4 |
| `TREATYGROUPNAME` | VARCHAR2 `[tidak disebut DBA]` | param 5 |
| `TREATYDESCID` | VARCHAR2 `[tidak disebut DBA]` | param 6 — master `TREATYDESC` |
| `TREATYDESCNAME` | VARCHAR2 `[tidak disebut DBA]` | param 7 |
| `REINSTYPEID` | VARCHAR2 `[tidak disebut DBA]` | param 8 |
| `REINSTYPENAME` | VARCHAR2 `[tidak disebut DBA]` | param 9 |
| `LAYER` | VARCHAR2 `[tidak disebut DBA]` | param 10 |
| `LAYERPART` | VARCHAR2 `[tidak disebut DBA]` | param 11 |
| `LAYERPARTTYPE` | VARCHAR2 `[tidak disebut DBA]` | param 12 |
| `LAYERTYPE` | VARCHAR2 `[tidak disebut DBA]` | param 13 |
| `KURS` | VARCHAR2 `[tidak disebut DBA]` | param 14 — teks `TREATYEXCHANGEYEARLY.TOIDR` apa adanya |
| `TGLUPDATE` | DATE | param 15 diabaikan procedure (`SYSDATE`); layanan mengikat jam simpan |
| `USERID` | VARCHAR2 `[tidak disebut DBA]` | param 16 — `OperatorID.pyUserName` |
| `LINE` | VARCHAR2 `[tidak disebut DBA]` | param 17 |
| `PCT` | VARCHAR2(1000) | param 18 — desimal teks |
| `PCTME` | VARCHAR2(1000) | param 19 — desimal teks |
| `YDCF` | VARCHAR2 `[tidak disebut DBA]` | param 20 |
| `METHOD` | VARCHAR2 `[tidak disebut DBA]` | param 21 |
| `TERRITORIALLIMIT` | VARCHAR2 `[tidak disebut DBA]` | param 22 |
| `PARENTREINSTYPEID` | VARCHAR2 `[tidak disebut DBA]` | param 23 — `"00"` = induk (OQ-TCO-20) |
| `SPREADINGORDER` | VARCHAR2 `[tidak disebut DBA]` | param 24 |
| `RP` | VARCHAR2(1000) | param 25 — desimal teks |
| `USD` | VARCHAR2(1000) | param 26 — desimal teks |
| `ID_OCCUPATION` | VARCHAR2 `[tidak disebut DBA]` | param 27 — induk saja |
| `OCCUPATION` | VARCHAR2 `[tidak disebut DBA]` | param 28 — induk saja |
| `ID_CLAUSE` | VARCHAR2 `[tidak disebut DBA]` | param 29 — induk saja |
| `CLAUSE` | VARCHAR2 `[tidak disebut DBA]` | param 30 — induk saja |
| `TREATYLIMIT` | NUMBER | param 31 — induk saja; bind angka |
| `COINS_MIN` | NUMBER | param 32 — induk saja |
| `COINS_MAX` | NUMBER | param 33 — induk saja |
| `MORERP` | NUMBER | param 34 — induk saja |
| `MOREUSD` | NUMBER | param 35 — induk saja |

Kolom katalog yang procedure **tidak** set: `OBJECT VARCHAR2(50)`, `PROPORTIONALLIST VARCHAR2(1000)` — tidak ditulis, tidak dibaca.

## M_ATTACHMENTTREATY_2

Lampiran tahun treaty. Kunci pemilik **`TREATYID = TreatyYear + TreatyYearID`** (teks disambung) — `TreatyOutSaveAttachment`
b1402 dan `DeleteAttachmentTreaty` b252 `[terverifikasi]`; **ralat** tiket 12 yang menyebutnya kunci treaty inward. Penulis
Treaty Contract Out: procedure `PEGA_M_ATTACHMENT(IDPEGA, DATAPEGA)` — badannya **`[terbuka — DBA]`** (OQ-TCO-24); kolom
dan bentuk ID ditiru dari penulis langsung tabel yang SAMA di `Treaty In/RDBList/InsertAttachment2_Sql.xml` b84.
Tipe: `[terbuka — DBA]`, diperlakukan VARCHAR2.

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | `[terbuka — DBA]` | `TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3')` (Saudara b84) |
| `TREATYID` | `[terbuka — DBA]` | `TreatyYear + TreatyYearID` |
| `CATEGORY` | `[terbuka — DBA]` | `.pyCategory` = `"File"` (`TreatyOutSaveAttachment` b558) |
| `FILENAME` | `[terbuka — DBA]` | nama berkas unggahan |
| `FILEMIMETYPE` | `[terbuka — DBA]` | MIME |
| `DATA_JSON` | `[terbuka — DBA]` | Saudara: `Datain.CARI50`; modul ini: NULL (nol kolom dokumen) |
| `USERNAME` | `[terbuka — DBA]` | `OperatorID.pyUserIdentifier` |
| `CATEGORY_ID` | `[terbuka — DBA]` | kategori pilihan (master `CATEGORY_ATTACH_REAS`); `GetAttachment2_Sql` menyaringnya |
| `T_STORAGE_ID` | `[terbuka — DBA]` | `T_STORAGE_IMAGE.IMAGEID` objek berkas |

## T_STORAGE_IMAGE

Objek berkas di layanan penyimpanan. Penulis: `Insert_T_Storage_SQL` (Saudara b85) sesudah unggah berhasil;
`Update_T_Storage_SQL` b85; `DeleteStorage_SQL` b85 sesudah hapus berhasil. Pembaca `GetLinkStorage_SQL` b85.
`IMAGEID` = `STANDARD_HASH('ASMPP' || SYSTIMESTAMP || SYS_GUID(), 'MD5')` (`GenerateImageID_SQL` b79; `models.ImageIDBaru`).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `IMAGEID` | `[terbuka — DBA]` | kunci |
| `URLPUBLIC` | `[terbuka — DBA]` | `URLImage` jawaban unggah; disegarkan jawaban geturl (`Update_T_Storage_SQL`, OQ-TCO-26) |
| `APPFOLDER` | `[terbuka — DBA]` | `appfolder` jawaban unggah, disegarkan jawaban geturl (OQ-TCO-26); delete mengirim `Namafile` = `APPFOLDER` dikurangi `gs://` + App + `/` (`DeleteGoogleStorage_Act` b1091) |
| `EXPDATE` | DATE | `To_date(exp, 'DD/MM/YYYY HH24:MI:SS')`; `exp` jawaban diubah seperti Pega dulu (`models.ExpStorageTCO`, b2146/b2211, unggah b2366/b2431) |
| `FILENAME` | `[terbuka — DBA]` | `Namafile` |
| `APPNAME` | `[terbuka — DBA]` | `T_FOLDER_IMAGE.APPNAME` (`GetAppName_SQL` b59) |
| `STORAGE` | `[terbuka — DBA]` | `'standard'` |
| `TANGGAL_UPLOAD` | DATE | `To_date(DateTime, 'MM/DD/YYYY HH24:MI:SS')` (Update b85) — diisi saat geturl (OQ-TCO-26); NULL sampai geturl pertama |

---

## Master yang dibaca saja

Tidak ditulis modul ini (`masterDibacaSajaTCO`): `REINSURANCETYPE`, `TREATYGROUP`, `TREATYDESC`, `TREATYEXCHANGEYEARLY`,
`CURRENCY`, `CATEGORY_ATTACH_REAS`, `AGENT`, `BUSINESS`, `OCCUPATION`, `CLAUSE`, `T_FOLDER_IMAGE`, `M_LINK_SERVICE`.
`GCP_IMAGE` ditulis lewat kode token bersama (`GET_TOKEN_STORAGE` ditiru). Outbox `T_LOG_SERVICE_RNM` milik aplikasi.
