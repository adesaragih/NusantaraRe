# 01: Skema relasional + migrasi + tipe dirapikan — **PREFACTOR**

**Status:** selesai (28-09-2026)

**Blocked by:** **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks, tidak dibuat
di sini)

⚠️ **Ini PREFACTOR, dan ia tiket PERTAMA.** Bentuk skema berubah di tiga sumbu sekaligus —
tabel JSON dibuang, `MTREATYSECURITY` dibersihkan, dan tipe kolom diperbaiki — sehingga **tidak ada
irisan lain yang dapat berdiri** sebelum bentuk barunya ada. *"Make the change easy, then make the
easy change."*

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin seluruh data master arrangement treaty non-life pindah ke bentuk
yang dapat dipercaya — uang sebagai angka, tanggal sebagai tanggal, setiap baris punya identitas —
tanpa kehilangan satu nilai pun, dan **tanpa mengambil apa pun dari tabel JSON yang sudah mati**.
*(User story 37–40 di spec)*

Dan sebagai **konteks hilir** (`Claim Prop`, `Komite Claim Prop`, `Claim Fac In`), saya ingin tetap
dapat **membaca** master arrangement dalam bentuk yang saya kenal sampai saya ikut bermigrasi.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL delapan tabel + sequence; skrip migrasi & rekonsiliasi |
| `internal/models` | Bentuk tahun treaty, kontrak, reinsurer, security, business, klausul |
| — | Skrip rekonsiliasi nilai uang & tanggal |

## Keadaan lama `[data DBA]`

| Tabel | Temuan tipe |
| --- | --- |
| `PROPORTIONALARRG` | ⚠️ **CAMPUR**: `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`, `MORERP`, `MOREUSD` = **`NUMBER`**; tetapi `RP`, `USD`, `PCT`, `PCTME` = **`VARCHAR2(1000)`**. Ada `TGLUPDATE DATE`, `OBJECT VARCHAR2(50)`, `PROPORTIONALLIST VARCHAR2(1000)` |
| `MTREATYSECURITY` | ⚠️ `THN_TREATY VARCHAR2(4) DEFAULT '1' NOT NULL`, `TOP_ID VARCHAR2(9)`, `TP_TREATY CHAR(2)`, `REAS_ID CHAR(7) NOT NULL`, `PCT_SHARE VARCHAR2(99)`, `USER_ID CHAR(99)`, `REAS_SECURITY CHAR(10) NOT NULL` — **tanpa primary key** |
| `TREATYCONTRACT` | `TREATYSTARTDATE`/`TREATYENDDATE` sudah `DATE`; tetapi `TGLUPDATE VARCHAR2(1000)` |
| `TREATYYEAR` | **seluruh** kolom `VARCHAR2` — termasuk `STARTDATE`, `ENDDATE` |
| `TREATYBUSINESS` | **seluruh** kolom `VARCHAR2` |
| `TREATYREINSURER` | `RICOMM`, `PCTSHARE` = `NUMBER`; sisanya `VARCHAR2` |
| `TREATYEXCHANGEYEARLY` | **seluruh** kolom `VARCHAR2` |
| `TREATYDESC` | `ID, DESCNAME, ISXOL, STATUSAKTIF` — seluruhnya `VARCHAR2` |

## Bentuk baru

```
treaty_year   (ID + grup + underwriting year + proporsi + STARTDATE/ENDDATE sebagai DATE)
  └─ treaty_contract  (IDTREATYYEAR → treaty_year.ID, REINSTYPEID, tanggal sebagai DATE)

Menggantung pada kunci gabungan (TREATYYEAR, TREATYGROUPID, REINSTYPEID) — BUKAN FK ke kontrak:
  ├─ treaty_reinsurer      (PCTSHARE, RICOMM sebagai desimal)
  │    └─ mtreaty_security (PK surrogate, REAS_ID → treaty_reinsurer.ID, PCT_SHARE desimal)
  ├─ treaty_business       (BIZCODE, BIZNAME, ISACTIVE)
  └─ proportionalarrg      (SATU tabel untuk 25 jenis klausul, dibedakan TREATYDESCID)
```

⚠️ **Kunci gabungan dipertahankan** `[fakta bisnis — work owner]`. Klausul dan anak-anak lain
menggantung pada **(TreatyYear, TreatyGroupID, ReinsTypeID)**, **bukan** pada `ID` kontrak.
Menambahkan foreign key ke kontrak akan **mengubah arti data** — jangan lakukan.

### Kolom `proportionalarrg` — 35 kolom `[data DBA]`

```
ID, TREATYYEAR, TREATYYEARID, TREATYGROUPID, TREATYGROUPNAME, TREATYDESCID, TREATYDESCNAME,
REINSTYPEID, REINSTYPENAME, LAYER, LAYERPART, LAYERPARTTYPE, LAYERTYPE, KURS, TGLUPDATE, USERID,
LINE, PCT, PCTME, YDCF, METHOD, TERRITORIALLIMIT, PARENTREINSTYPEID, SPREADINGORDER, RP, USD,
ID_OCCUPATION, OCCUPATION, ID_CLAUSE, CLAUSE, TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD
```

Sembilan kolom terakhir yang bercetak — `ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`,
`TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`, `MORERP`, `MOREUSD` — hanya diisi baris **induk**; baris
"anak" mengisinya **NULL**.

### Sequence `[data DBA]`

| Tabel | Sequence | Format identitas |
| --- | --- | --- |
| `PROPORTIONALARRG` | `PROPORTIONALARRG_SEQ` | `'1' + lpad(seq, 7, '0')` — **7 digit** |
| `TREATYCONTRACT` | `treatycontract_seq` | `'1' + lpad(seq, 6, '0')` |
| `TREATYYEAR` | `TreatyYear_seq` | `'1' + lpad(seq, 6, '0')` |
| `TREATYREINSURER` | `M_TREATYREINSURER_SEQ` | `'1' + lpad(seq, 6, '0')` |
| `TREATYBUSINESS` | `TREATY_BUSINESS_SEQ` | `'1' + lpad(seq, 6, '0')` |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterProportionalArrg` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `ASM!SAVEMASTERPROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterProportionalArrg.xml` | 35 kolom induk |
| `SaveMasterProportionalArrgChild` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterProportionalArrgChild.xml` | 26 kolom anak |
| `GetMasterDescriptionLimitParentList` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterDescriptionLimitParentList.xml` | baca **relasional** |
| `GetMasterDescriptionEPIParentList` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterDescriptionEPIParentList.xml` | ⚠️ baca **JSON** — mati |
| `DeleteRowBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `ASM!DELETEROWBUSINESSLIST` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteRowBusinessList.xml` | ⚠️ memegang **kedua** tabel kembar |
| `InsertToMTreatySecurity` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `ASM!INSERTTOMTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/InsertToMTreatySecurity.xml` | ⚠️ INSERT posisional |

⚠️ **Penyimpangan sadar 1 — dualitas JSON dibuang.** `[keputusan work owner]`
`[data DBA]` **Tidak satu pun dari enam procedure penulis menyentuh tabel JSON.**
`M_PROPORTIONALARRG` sudah lama tidak dipakai; kueri yang masih membacanya adalah sisa yang lupa
dihapus. Sumber kebenaran tunggal = **`PROPORTIONALARRG` relasional**.

⚠️ **Penyimpangan sadar 5 — `MTREATYSECURITY` dibersihkan.** `[keputusan work owner]`
PK surrogate, kolom bernama dan diisi eksplisit, `PCT_SHARE` desimal, `REAS_SECURITY` atribut biasa.

⚠️ **Penyimpangan sadar 6 — seluruh uang/persen jadi desimal, seluruh tanggal jadi `DATE`.**
`[keputusan work owner]`

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0006** (identitas lewat sequence basis data),
**ADR-0009** (migrasi penuh; koeksistensi ditolak untuk data — **bukan** untuk kontrak baca hilir).

## Acceptance criteria

- [ ] Modul ini adalah **satu-satunya penulis** `TREATYYEAR`, `TREATYCONTRACT`, `TREATYREINSURER`,
      `MTREATYSECURITY`, `TREATYBUSINESS`, `PROPORTIONALARRG`. *(AC 1 spec; OQ-042)*
- [ ] Konteks hilir (`Claim Prop`, `Komite Claim Prop`, `Claim Fac In`) tetap dapat **membaca**
      bentuk relasional yang mereka pakai hari ini. *(AC 2 spec)*
- [ ] Skema **tidak memuat** satu pun objek treaty **outward**. *(AC 3 spec)*
- [ ] ⚠️ **Seluruh** nilai uang — `RP`, `USD`, `MORERP`, `MOREUSD`, `TREATYLIMIT`, `COINS_MIN`,
      `COINS_MAX`, `PCT_SHARE`, `RICOMM` — bertipe **desimal presisi arbitrer**; **tidak** melewati
      `float`. Test yang menemukan kolom uang bertipe teks **gagal**. *(AC 51 spec; **ADR-0003**;
      penyimpangan sadar 6)*
- [ ] ⚠️ **Seluruh** persentase — `PCT`, `PCTME` — desimal, tidak dibulatkan ke bilangan bulat dan
      tidak disimpan sebagai teks. *(AC 52 spec; penyimpangan sadar 6)*
- [ ] ⚠️ **Seluruh** tanggal — `STARTDATE`, `ENDDATE`, `TREATYSTARTDATE`, `TREATYENDDATE`,
      `TGLUPDATE` — bertipe **`DATE`**. *(AC 53 spec; penyimpangan sadar 6)*
- [ ] ⚠️ Skema target **relasional penuh**. Test yang menemukan kolom JSON sebagai penyimpan
      atribut arrangement **gagal**. *(AC 63 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Migrasi **tidak mengambil apa pun** dari `M_PROPORTIONALARRG`, `M_TREATYCONTRACT`,
      `M_TREATYBUSINESS`, atau tabel `M_*` lain. Sumbernya **hanya** tabel relasional.
      *(AC 64 spec; penyimpangan sadar 1)*
- [ ] Seluruh baris keenam tabel pindah **tanpa kehilangan satu nilai pun**. *(AC 65 spec;
      **ADR-0009**)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
      tepat**, bukan dengan toleransi. *(AC 66 spec; **ADR-0003**)*
- [ ] ⚠️ Nilai uang & persen yang hari ini berupa **teks berkoma desimal** (mis. `"12,5"`) terurai
      **benar** menjadi desimal; teks yang **tidak dapat diurai** dilaporkan, **tidak** didiamkan
      dan **tidak** diam-diam jadi nol. *(AC 66 spec; `[terverifikasi]` — existing memanggil
      `@replaceAll(.Pct,",",".")` di `Activity/TreatyTestChildTotal_Act.xml`,
      `@BASECLASS!TREATYTESTCHILDTOTAL_ACT`, dan `Activity/SetErrorMessageReinsurer.xml`)*
- [ ] Tanggal yang hari ini berupa **teks** menjadi `DATE` **tanpa pergeseran zona waktu**; teks
      yang tidak dapat diurai **dilaporkan**. *(AC 67 spec)*
- [ ] ⚠️ Baris `MTREATYSECURITY` mendapat **primary key surrogate** saat migrasi, dan rujukannya ke
      reinsurer tetap utuh. Test yang menemukan kunci berbasis nama security **gagal**.
      *(AC 68 spec; penyimpangan sadar 5)*
- [ ] Kelima sequence pindah dengan **nilai berjalan yang benar**, sehingga identitas baru **tidak
      bertabrakan** dengan yang lama. *(AC 69 spec; **ADR-0006**)*
- [ ] `[terbuka]` Kolom `PROPORTIONALLIST` dan `OBJECT` **tidak dibawa** ke skema baru kecuali
      migrasi membuktikan ada data hidup di sana; temuannya **dilaporkan**. *(AC 70 spec)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.
      *(AC 71 spec)*

## Blocker

**Tidak ada pemblokir.** `[data DBA]` DDL delapan tabel dan body enam procedure sudah diterima —
**OQ-001 dan OQ-002 ditutup**.

## Catatan

⚠️ **Kejanggalan yang wajib dicatat, bukan ditiru.**

- `RDBList/DeleteRowBusinessList.xml` menghapus `treatybusiness` di `pyBrowseSQL` **dan**
  `m_treatybusiness` di `pyDeleteSQL` — satu rule, dua tabel kembar. Sebaliknya
  `RDBList/DeleteFromTREATYCONTRACT_SQL.xml` menghapus `treatycontract` tetapi **bukan**
  `M_TREATYCONTRACT`. Inkonsistensi ini **lenyap dengan sendirinya** begitu JSON dibuang.
- `[data DBA]` Procedure bernama **`PEGA_M_PROPORTIONALARRG_CHILD`** **tidak** menulis ke tabel
  child — ia menulis ke tabel yang **sama**. Nama menipu (**OQ-066**).
- `[data DBA]` Teks galat sebagian procedure menyebut **`"JSON_KLAIM"`** — sisa salin-tempel
  template, **bukan** petunjuk bahwa data berbentuk JSON.
- `[terverifikasi]` `Activity/GetPeriode.xml` (`@BASECLASS!GETPERIODE`) memakai
  `ParentReinsTypeID = "00"` sebagai penanda baris tanpa induk. Nilai sentinel ini perlu dibawa
  utuh saat migrasi.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — presisi desimal, konversi tanggal, dan
keutuhan rujukan setelah pemberian PK **hanya berperilaku benar pada basis data sungguhan**;
memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```


---

## Pembacaan ulang XML — 28-09-2026 (sesi modul, sebelum kode)

Nomor baris = `sed -e 's/></>\n</g'` atas berkas `D:\XML\RNM_BRD\Treaty Contract Out\`. Pohon penyarangan
dibaca dari `Struktur_InboxTreatyContract.xlsx` (lewat `.scratch/alat/baca-xlsx.ps1`; ⚠️ berkas itu
hanya me-root **satu** harness, `InboxTreatyContract` r2 — dua harness lain dibaca langsung dari XML).

### Peta harness → section → tombol → activity → RDB/SQL

**1. `Harness/InboxTreatyContract.xml`** (`DATA-PORTAL!INBOXTREATYCONTRACT` b30; `<pyLabel>InboxTreatyContract</pyLabel>` b151)
→ `Section/GridTreatyContract.xml` (xlsx r3; judul `TREATY CONTRACT OUT` b1057) → `Section/InputTreatyContract.xml` (b2482; xlsx r4):

| Unsur | Baris | Activity → RDB/SQL | Tiket |
| --- | --- | --- | --- |
| grid tahun treaty `pyGridRDName BrowseTreatyYear_RD` | b17355; kolom `Underwriting Year` b17378, `Treaty Group` b17990, `.UnderwritingYear` b18964, `.TreatyYear` b19125, `.StartDate` b19286, `.EndDate` b19446, `.TreatyGroupName` b19605, `.Proportion` b19724 | `ReportDefinition/BrowseTreatyYear_RD.xml` (kelas `ASM-FW-GISFW-Int-TREATYYEAR`; 10 field b594–b735; sort `.ID DESC` b672; max 500 b757) | 03 |
| tombol baris `Edit` | b19939 → b20030 | `Activity/SetTreatyYear_Act.xml` — memuat `BrowseTreatyYear_RD` b449 lalu menyalin Param ke `InputTreatyYear.*` b1110–b1353 | 03 |
| tombol `Add` | b16387 | `Activity/NewInputTreatyYear_Act.xml` — mengosongkan `InputTreatyYear.ID/TreatyYear/TreatyGroupID` b469–b512 | 03 |
| form `Input New Data` | b6351; `ID` b7282 (`InputTreatyYear.ID`), `Treaty Year` b7455 (+ `CheckYear` b7520/b7646), `Treaty Group` b7732, `Reinsurance Type` b7925 (`InputTreatyYear.ReinsuranceType`), `Modified Date` b8661, `By User` b9021 (`OperatorID.pyUserName`) | — | 03 |
| tombol `Save` | b10059 → b10082/b10166 | `Activity/SaveTreatyYear_Act.xml`: `UserID←OperatorID.pyUserName` b280, `TglUpdate←@getCurrentTimeStamp()` b327; prasyarat `TreatyGroupID==""` b388, `TreatyYear==""` b411, `@Default.isNumber(TreatyYear)` b434; RDB-List `SaveMasterTreatyYear_SQL` b1100 → `POOLDATA.PEGA_TREATYYEAR` (10 param + `HASIL1/HASIL2 out`) | 03 |
| tombol `Cancel` | b10336 → b10354 | `Activity/CancelActivityTreatyContract.xml` — mengosongkan `DATASHOW3/HASILD7/DATASHOW` b272–b341 | 03 |
| varian form kedua `Section/InputDtlTreatyContact.xml` | disertakan b1206; `Treaty Group` b6560, `Reinsurance Type` b6800 → **`InputTreatyYear.Proportion`**, `Start Date` b7532, `End Date` b7816, ⚠️ `Underwriting Year` b8004 → `InputTreatyYear.TreatyYear`, ⚠️ `Transaction Year` b8284 → `InputTreatyYear.UnderwritingYear`, `Save` b10332 → `SaveTreatyYear_Act` b10356 | idem | 03 (OQ-TCO-05: label bersilang) |
| form salin `From`/`To` + `Proces` | b2374/b3818; b5104 → `BrowseCopyData` b5128/b5211 | ➖ MATI — fitur salin dibuang (penyimpangan sadar 3) | — |
| tombol baris `ReinsType` | b20778 (membawa `.TreatyGroupName/.TreatyYear/.StartDate/.EndDate/.UnderwritingYear` b21310–b21346) | membuka konteks kontrak (harness 2) | 04 |
| tombol baris `List Description` | b22196 (membawa `.TreatyYear` b22743) | membuka konteks klausul (harness 3) | 08 |
| `Attachment for` + `Section/GridTreatyArrangementAttachment.xml` | b11721; xlsx r5–r46: `DeleteAttachmentTreaty` → `Delete_act` (`ASM-FW-GISFW-Int-TREATY_IN`) → `LoadAttachment` → `GetAttachment2_Sql`; `DeleteAttachment2_Sql`; `DeleteGoogleStorage_Act` → `GetLinkService` | rantai teknis lampiran (belum rampung di Pega) | 12 |
| `Section/ViewDetailDescription.xml` | xlsx r47 (1.1.1.2) — lihat harness 3 | | 08 |

**2. `Harness/InboxTreatyContractReinsType.xml`** (`DATA-PORTAL!INBOXTREATYCONTRACTREINSTYPE` b30; `<pyLabel>InboxTreatyContractReinsType</pyLabel>` b151; judul `ReinsType` b1670)
→ `Section/PanggilReinsType.xml` (b1825) → `Section/InputTreatyContractReinsType.xml` (b1300):

| Unsur | Baris | Activity → RDB/SQL | Tiket |
| --- | --- | --- | --- |
| kepala: `Underwriting Year` b1145, `ReinsType`(= nama grup) b1358, `.IDTreatyYear` b1553 (`InputTreatyContract.IDTreatyYear`), `.TreatyGroupID` b1757 | | | 04 |
| pemilih `ReinsType` | b2652 (`InputTreatyContract.ReinsTypeName`, tampil `.Note` b2724, saringan `"active"` b2768) | **`ReportDefinition/BrowseReinsuranceType_RD.xml`** (non-Old) — satu-satunya pemakainya | 02/04 |
| `Start Date` | b2905 (`InputData.CARIDATETIME`) → b3007/b3155 | `Activity/SetTanggalTreatyContract.xml` — `TreatyEndDate` = start + 365/366 hari b620–b1274 (tahun kabisat b1191); validasi tahun start = `TreatyYear` b847 | 04 |
| `End Date` | b3244 (`InputData.CARIENDDATE`) | | 04 |
| tombol `Save` | b3618 → b3642/b3727 | `Activity/SaveTreatyContract_Act.xml`: `ReinsTypeID==""` b585 → `"Data Reins Masih Kosong!!!"` b1266; `TreatyStartDate==""` b760; tahun start ≠ `TreatyYear` b1199; `ID="UnknownId"` bila kosong b827; tanggal diformat `dd/MM/yyyy` b1040/b1069 dan b1479/b1508; `UserID←pyUserName` b971, `TglUpdate←getCurrentTimeStamp` b1018; `ERRMSG3="Data sudah pernah di Input"` b2270 bila `STSSAVE==3` b2348; RDB-List `SaveMasterTreatyContract_SQL` b2649 → `PEGA_TREATYCONTRACT` (8 param) | 04 |
| tombol `Undo` b5343 → `UndoOperation` b5366; `Add` b8528 → `NewInputTreatyContract_Act` b8552 (mengosongkan `InputTreatyContract.*` b502–b587) | | | 04 |
| grid kontrak: `Reins Type` b9164, `Treaty Start` b9304, `Treaty End` b9444; `Edit` b10519 → `SetUbahTreatyContract` b10543 (Param `.ID/.ReinsTypeID/.ReinsTypeName/.TreatyStartDate/.TreatyEndDate` b10559–b10583) | | `DataPage/D_TreatyContract.xml` (kelas `ASM-FW-GISFW-Int-TREATYCONTRACT`) | 04 |
| `Section/ViewDetailTreatyReinsurerGrid1.xml` | disertakan b13364 | `Add`/`Tambah` b1730/b1996 → `NewTreatyReinsurerDetail_Act`; kolom `ReinsID` b2418, `Reinsurer` b2560, `%Share` b2702, `%Comm` b2844, `Rating` b2990, `Operator Name` b3138; `Edit` b4491 → `SetUbahTreatyReinsurerList_Act` b4543; `Delete` b4936 → `DeleteTreatyReins_Act` b4953 → `RDBList/DeleteFromTreatyReinsurer_Act.xml` (hapus `MTREATYSECURITY` lalu `TREATYREINSURER`, tanpa COMMIT); `Security Reinsurer` b5277 (membawa `.TreatyYear` b5301, `.ID` b5307); `Total Share -->>` b6186 (`InputTreatyReinsurer.TotalShare` b6330); form `Reins.ID` b8042, `Reinsurer` b8226, `%Share` b8522 (+ `SetErrorMessageReinsurer` b8586), `%Comm` b8800, `Rating` b9076; `Save` b11405 → `SaveTreatyReinsurerDetail1_Act` b11429 | 05/06 |
| `Activity/SaveTreatyReinsurerDetail1_Act.xml` | `ReinsTypeID/TreatyGroupID/TreatyYear ← InputData.CARI6/5/4` b398–b441; `ID="UnknownId"` b549; `ReinsurerID` kosong → `"Data tidak boleh kosong...!!!"` b694; loop `pTotalShare += toDecimal(replaceAll(.PctShare,",","."))` b1016; ⚠️ gerbang `TotalShare + (PctShare − PctShare1) <= 100.000` b1452/b1645/b1875; pesan `"Persentase tidak boleh lebih dari 100!"` b1357, `"Total Share tidak boleh lebih dari 100%"` b1561; RDB `SaveMasterTreatyReinsurer_SQL` → `PEGA_TREATYREINSURER` (19 param) | | 05 |
| `Activity/SaveSecurityReinsurer_Act.xml` | report `SelectSecurityReinsurer` b289 (`THN_TREATY←CARI8` b357, `REAS_ID←CARI9` b378); `IsUpdate` bila `REAS_SECURITY` cocok b837; `THN_TREATYID←IDTreatyYear` b999 (**tidak dipersist** — INSERT posisional 7 nilai); INSERT b1095 (`HASILD3==0`), UPDATE b1279 (`HASILD3==1`) → `InsertToMTreatySecurity.xml` / `UpdateMTreatySecurity.xml` | | 06 |
| `Section/ViewDetailTreatyBusinessGrid.xml` | disertakan b14117 | `Business List` b1988; `Add` b2785 → `NewTreatyBusinessDetail_Act`; grid `Treaty Group` b3422, `Business ID` b3566, `Business Name` b3710 (saringan `InputData.CARI1/2/3` b3288–b3300 → `RDBList/GetMasterBusinessList.xml`); `Edit` b4547 → `SetUbahTreatyBusinessList_Act`; `Delete` b4826 → `DeleteRowBusiness` b4850 → `RDBList/DeleteRowBusinessList.xml` (`treatybusiness` + ⚠️ `m_treatybusiness`); form `Business Name` b6241 (`.Note`), `Active` b6499, `Business Code` b6680; `Save` b6966 → `SaveTreatyBusinessDetail_Act` (`TreatyYearID←InputTreatyContract.IDTreatyYear` b354; `"Data sudah pernah di Input"` b1103) → `PEGA_TREATYBUSINESS`; `Close List` b10889 → `CancelActivity` | 07 |
| `Activity/BrowseDeleteRowTreatyInContract.xml` | `CARI17←TreatyYear` b372, `CARI18←TreatyGroupID` b419, `CARI16←IDTreatyYear` b440, `CARI19←ReinsTypeID` b461, `CARI20←ID` b482 → `RDBList/DeleteFromTREATYCONTRACT_SQL.xml` (4 DELETE + COMMIT); `"Data Berhasil di Hapus"` b762 | pemanggil tombolnya ditelusuri di tiket 10 | 10 |

**3. `Harness/InboxTreatyContractDescription.xml`** (`DATA-PORTAL!INBOXTREATYCONTRACTDESCRIPTION` b26; `<pyLabel>InboxTreatyContractDescription</pyLabel>` b359)
→ `NitipKurs` b3882, `SubViewDetailDescription` b11366, `ViewDetailDescriptionProp` b12583, `ViewDetailDescriptionNonProp` b13547 — sama dengan `Section/ViewDetailDescription.xml` (xlsx 1.1.1.2):

| Unsur | Baris | Activity → RDB/SQL | Tiket |
| --- | --- | --- | --- |
| daftar jenis klausul `For Non XOL` / `For XOL` (kolom `ID`, `Description Name`) | b3847/b4410/b4520; b6952/b7515/b7625 | `ReportDefinition/BrowseTreatyDesc_RD.xml` (`.ID`, `.DescName`, `.IsXOL`, `.StatusAktif`; filter `.IsXOL`) — master `TREATYDESC` dibaca saja | 08 |
| tombol `Show` | b5031/b8132 → `BrowseDescriptionLimit` b5049/b5457/b8150/b8589, `testingKurs` b5110/b5513/b8245/b8679, `SetKirimIDDesc` b5154, `PanggilID` b5228, `RefreshErrorProportionalarrg` b5348/b5730 | `Activity/testingKurs.xml`: `CARI1←Param.StartDate` b273 → `RDBList/GetMasterKursList.xml` (`TOIDR … IDCURRENCY='10001' … Quarter='0'`) → `InputTreatyArrangement.Kurs←.HASIL1` b580 | 11 |
| `NitipKurs` | `Section/NitipKurs.xml` b505 `.Kurs` (`InputTreatyArrangement.Kurs`) | `Activity/RefreshKurs.xml` mengosongkan `Kurs` b245 | 11 |
| 25 grid klausul (xlsx r48–r139+: `GridTreatyArrangementProfitCommision`, `…ExGratia` + `gridTreatyArrangementExGratiaList` + `GridTreatyArrTreatyExGratiaChildList`, `…CashLossLimit` + `GridTreatyArrTreatyCashLossLimitList`, `…Portfolio` + `…PortfolioList`, `…TreatyLimit` + `GridTreatyArrTreatyLimitList`, …) | tiap grid: `New…`, `Save…_Act` → `SaveMasterProportionalArrg` (induk) / `SaveMasterProportionalArrgChild` (anak), `Set…_Act`, `SetErrorMessage`, `CancelActivity…`, `HitungRpUsd`, RD `BrowseTreatyArrangement_*`, RD `BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull` (pemilih jenis reasuransi di 11 grid) | `Activity/HitungRpUsd.xml`: anak `Usd = Pct × Usd_induk / 100`, `Rp = Pct × Rp_induk / 100` b251–b607 (7 jenis); `Activity/TreatyTestChildTotal_Act.xml`: Σ`Pct` anak lewat `GetMasterDescriptionLimitParentList` b1310, tolak bila `TotalPCt != 100.00` b331/b825 → `" Please make sure spreading is 100%"` b892; `Activity/SetErrorMessage.xml`: `Pct>100 || Pct<0` per jenis b346–b1957 | 08 |
| `GetMasterPortfolioListDetail`, `GetMasterPanggilID`, `GetMasterDescriptionEPIParentList`, … (`FROM m_PROPORTIONALARRG`) | xlsx r122–r136 | ➖ MATI — dibaca ulang dari `PROPORTIONALARRG` relasional | 08 |

**WHEN**: modul ini tanpa rule `When` (folder tidak ada); seluruh syarat adalah `pyStepsPreCondParamsWhen` per
langkah, dibaca dari urutan langkah di atas. **Penulis medan** enam tabel = enam procedure lewat RDB-List
`SaveMaster*` + tiga SQL mentah (`InsertToMTreatySecurity`, `UpdateMTreatySecurity`, `DeleteSecurityReinsurer`)
+ tiga penghapus (`DeleteFromTREATYCONTRACT_SQL`, `DeleteFromTreatyReinsurer_Act`, `DeleteRowBusinessList`).

### Kontrak hilir yang dibaca (spec b224–b232, tco3) — kolom VERBATIM

| Hilir | Rule | Tabel | Kolom dibaca | Saringan |
| --- | --- | --- | --- | --- |
| Claim Prop, Claim Fac In | `GetLimitPLATreatyin.xml`, `GetLimitPLADLA_Sql.xml` | `PROPORTIONALARRG` | `TREATYYEAR, TREATYGROUPID, TREATYDESCID, REINSTYPEID, REINSTYPENAME, PCT, RP, USD` | `treatydescid='10001' and treatyyear and treatygroupid and REINSTYPEID` |
| Claim Fac In | `GetQuotaShare.xml` | `PROPORTIONALARRG` | idem | `treatyyear, treatygroupid, TREATYDESCID='10001', PARENTREINSTYPEID` |
| Claim Prop, Komite Claim Prop, Claim Fac In | `GetListRetro_Sql.xml` (teks identik) | `TREATYREINSURER` | `id, reinsurerid, clientid, name, ricomm, pctshare` | `reinstypeid, treatyyear, treatygroupid` |
| Claim Prop | `GetTreatyGroupID.xml` | `TREATYBUSINESS` | `distinct TREATYGROUPID` | `BIZCODE, treatyyear, isactive='1'` |
| Claim Fac In | `GetTreatyGroup_Sql.xml` | `TREATYBUSINESS` | `ID, TREATYYEAR, TREATYGROUPID, REINSTYPEID` | `BIZCODE, treatyyear, reinstypeid` |
| Claim Fac In | `GetDataTreatyLimit_Sql.xml` | `treatybusiness JOIN proportionalarrg` (`treatyyearid`, `reinstypeid`) + `EXISTS treatycontract` (`idtreatyyear`, `reinstypeid`, tanggal `BETWEEN treatystartdate AND treatyenddate`) | idem klausul | `bizcode, reinstypeid, treatydescid='10001'` |
| Claim Prop | `GetTreatyInMasterProp_SQL.xml` | `TREATYBUSINESS` (`TREATYGROUPID, BIZCODE, BIZNAME`), `TREATYYEAR` (`TREATYGROUPID, STARTDATE, ENDDATE`) + `TREATYINDETAIL`, `treaty_in` | — | dibaca lewat gabungan; pembaca khusus ditulis saat Claim Prop dibangun |

Pembacanya: `repository/tco_kontrak_hilir.go` — enam fungsi read-only; dikunci `tco_kontrak_hilir_test.go`
tiga sisi (DDL, korpus, nol kata kerja tulis).

### Ralat bertanggal 28-09-2026 — XML membantah tiket/brief

1. **Tiket 02, bab Catatan** (*"`Harness/InboxTreatyContractReinsType.xml` … adalah perilaku layar master
   [jenis reasuransi], konteks/menu tersendiri"*) — **DIBANTAH**: harness itu → `PanggilReinsType` b1825 →
   `InputTreatyContractReinsType` b1300 adalah **editor kontrak per tahun** (`Save` → `SaveTreatyContract_Act`
   b3642) beserta grid business (b14117) dan reinsurer/security (b13364). RD non-Old
   `BrowseReinsuranceType_RD` dipakai **pemilih** `ReinsType` di form kontrak itu (b2652), bukan layar master
   terpisah. Harness ini adalah butir menu kedua modul (tiket 04), bukan di luar konteks.
2. **Brief §2** (*"kelompok Treaty Contract Out sudah ada di sidebar sebagai belum dimigrasi"*) —
   **DIBANTAH**: `frontend/src/assets/labels.ts` `MODUL` memuat 17 kelompok; `Treaty Contract Out`,
   `Treaty In`, `Treaty In Adjustment` tidak ada (folder korpus 20; ralat §8 `PROMPT-EKSEKUSI-HULU-HILIR.md`).
   Kelompok ditambahkan ADITIF bersama layar pertama (tiket 03); dua kelompok lain bukan lingkup sesi ini.
3. **Brief §2** (*"baca `Struktur_InboxTreatyContract.xlsx` untuk pohon penyarangan"*) — berkas itu hanya
   me-root `InboxTreatyContract` (r2); dua harness lain tidak ada di dalamnya, dibaca dari XML langsung.
4. **Tiket 05, bab Blocker** (*"Existing tidak menolak kombinasi yang totalnya ≠ 100 pada tingkat
   reinsurer"*) — **sebagian dibantah**: `SaveTreatyReinsurerDetail1_Act.xml` b1452/b1645/b1875 mensyaratkan
   `TotalShare + (PctShare − PctShare1) <= 100.000` dan menulis `"Total Share tidak boleh lebih dari 100%"`
   b1561 — total **> 100 DITOLAK**; yang benar: total **= 100 tidak diwajibkan**. Ditegakkan di tiket 05.
5. **tco1** (*"+ dua tabel lain dari tiket 01"*) — ditafsirkan: `T_TREATYCO_JEJAK` (jejak audit ADR-0007,
   dibuat di 306) dan tabel lampiran tahun treaty (fitur baru tiket 12, migrasi 307). `TREATYEXCHANGEYEARLY`
   dan `TREATYDESC` adalah master dibaca saja (spec b107) — **tidak** dibuat ulang ber-`T_`.

### Keputusan tiket 01 yang mengikat tiket berikut

- Nama tabel `T_TREATYYEAR`, `T_TREATYCONTRACT`, `T_TREATYREINSURER`, `T_MTREATYSECURITY`, `T_TREATYBUSINESS`,
  `T_PROPORTIONALARRG`, `T_TREATYCO_JEJAK`; sequence `SEQ_T_*`; konstanta di `repository/tco_migrasidata.go`.
- Identitas lewat `DB.IdentitasBerikutTCO(ctx, tx, seq)` — `'1'+lpad(6)` (klausul 7); nomor yang tidak muat
  **gagal terang** (`ErrIdentitasMelampauiLebar`), bukan dipotong seperti `LPAD`.
- `IUDATE` dan `PROPORTION` tetap teks (`[terbuka]`); `TOP_ID`, `TP_TREATY`, `USER_ID` security dibawa bernama.
- FK hanya dua: kontrak → tahun (tanpa kaskade), security → reinsurer (`ON DELETE CASCADE`). Reinsurer,
  business, klausul menggantung pada kunci gabungan (index gabungan, tanpa FK).
- Anti-dobel tahun (AC 73) ditegakkan di Go, bukan unique index (data warisan boleh sudah berduplikat).
- Migrasi data: `-migrate-data-treaty-contract-out` (tco2) — satu transaksi, temuan memblokir membatalkan
  seluruhnya, rekonsiliasi tepat, sequence diselaraskan sesudah commit; **tidak dijalankan di DEV oleh sesi ini**.

### OQ dibuka tiket 01

| OQ | Isi | Pemilik |
| --- | --- | --- |
| OQ-TCO-01 | bentuk teks tanggal warisan `TREATYYEAR.STARTDATE/ENDDATE/TGLUPDATE`, `TREATYREINSURER.STARTDATE/ENDDATE/TGLUPDATE`, `TREATYBUSINESS.TGLUPDATE`, `TREATYCONTRACT.TGLUPDATE` — pengurai mengenal 7 bentuk (`bentukTanggalWarisanTCO`), sisanya dilaporkan; dipastikan dari laporan pemindahan di DEV | DBA / work owner |
| OQ-TCO-02 | arti dan bentuk `TREATYREINSURER.IUDATE` | Product + UW |
| OQ-TCO-03 | isi hidup `PROPORTIONALARRG.PROPORTIONALLIST` / `OBJECT` — dicacah `KolomMatiBerisi` (AC 70) | work owner |
| OQ-TCO-04 | arti `TREATYYEAR.PROPORTION` (diisi pilihan "Reinsurance Type" `.ID` di `InputDtlTreatyContact.xml` b6829) | Product + UW |
| OQ-TCO-05 | `InputDtlTreatyContact.xml`: label `Underwriting Year` b8004 terikat `InputTreatyYear.TreatyYear`, label `Transaction Year` b8284 terikat `InputTreatyYear.UnderwritingYear` — bersilang dengan `InputTreatyContract.xml` b7455/b12020 | Product + UW (dipakai tiket 03) |

**Status:** selesai 28-09-2026 — commit `treaty-contract-out: tiket 01 — skema relasional + migrasi + tipe dirapikan`.
