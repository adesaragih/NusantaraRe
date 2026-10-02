# Paritas layar dan aksi — Master Contract Retro Life

> Disusun 30-09-2026 (sesi implementasi modul, paket 0, brief `PROMPT-IMPLEMENTASI-MODUL-MASTER-CONTRACT-RETRO-LIFE.md`).
> Satu baris per harness, section, tombol, dan aksi korpus → rule yang dipanggil → RDB/RD → tabel/kolom → rute API → komponen React.
> Nomor baris `bNNN` = baris yang **memuat nilai yang dikutip** pada berkas korpus hasil pecah standar, dapat diulang dengan
> `sed -e 's/></>\n</g' "<Folder>/<Tipe>/<Rule>.xml" | grep -n '<teks yang dikutip>'` atas `D:\XML\RNM_BRD\Master Contract Retro Life\`
> *(berkas korpus sudah satu tag per baris, jadi nomor ini sama dengan nomor baris mentahnya)*. Langkah activity dirujuk lewat baris
> `pyStepsActivityName`-nya; prakondisi lewat `pyStepsPreCondParamsWhen`; saringan RD lewat `pyFilterName`; grid lewat `pyBodyType` `REPEATING`.
> ⚠️ Dihitung ulang 01-10-2026 (gelombang 2 butir A2): paket 0–11 memakai baris AWAL elemen (`<rowdata>`), yang tidak dapat di-grep.
> Setiap langkah activity yang dikutip mencetak `pyStepsBlockName` — `·` = kosong (langkah hidup), `//` = ter-remark (tidak pernah jalan).
> Precondition dibaca dari urutan aksi: `WhenTrue 2` = lanjut, `3` = lewati langkah; `PRE=false` = prakondisi dimatikan → langkah **selalu** jalan.
> Keadaan: ✅ dibangun · ⏸️ dibangun, datanya menunggu OQ · ➖ sengaja tidak dibawa (bukti di kolom Catatan) · ➕ tambahan tiket tanpa padanan XML (OQ).
> Keadaan per 01-10-2026 (paket 11): seluruh baris yang di paket 0 bertanda "menunggu paket N" kini ✅ atau ⏸️. Commit: paket 0 `99c2bac` · paket 1 `db7f636` · paket 2 `fd42b86` · paket 3 `38c3d42` · paket 4 `c10d50f` · paket 5 `484f9d4` · paket 6 `b3e097d` · paket 7 `81b78bd` · paket 8 `422bd37` · paket 9+10 `a3c07bc`.

## 0. Alat baca dan sensus

| Hal | Isi |
| --- | --- |
| Alat | pengurai expat bernomor baris (`pohon_activity.py`, `pohon_layar.py`, `baca_rd.py`, `sql_rdb.py` — scratchpad sesi, tidak di repo) |
| Uji alat | `CountingPercentShare_Act` = 5 langkah termasuk 3.1 (cocok `TAMBAHAN-SPEC-ronde-2.md` T3); `TreatyLimit_TypeProtect` = 5 langkah termasuk 4.1 (cocok T3); langkah ber-`//` = 1 di seluruh modul (cocok `grilling-ronde-2.md` H3 butir 5) |
| Langkah ber-`//` | **1**: `NewInputBusinessLife_Act` langkah 1 b270 `Page-New InputBusinessLife` — 26 activity lain nol |
| Tombol — cara 1 | cacah baris `<pyFormat>pxButton</pyFormat>` di 8 Section: **33** |
| Tombol — cara 2 | pengurai: sel `Embed-Display-Table-Cell` ber-`pyFormat = pxButton`: **33** (sepakat per berkas) |
| Tombol unik | **31** — `Save` b5521 dan `Cancel` b5800 di `InputRetrocessionLife.xml` adalah salinan `InputDtlRetrocessionLife` yang disertakan lewat `pyInclude` b459 |
| Harness vs Section | salinan section di dalam keempat harness **identik** dengan berkas Section-nya (pohon kontrol dibandingkan tanpa nomor baris: beda 0, selain pembungkus include) — tabel ini berpijak pada berkas Section |
| Medan visibilitas | sel: `pyUserData.pyVisible` (`ALWAYS`/`OTHER`/`NOTBLANK`) + `pyCondition`; wadah: `Embed-Harness-Section.pyContainerVisibleWhen` **berlaku hanya bila** `pyIsVisibilityOption = CONDITION` |
| Medan RD | kolom `Embed-ReportListFields.pyFieldName`; saringan `Embed-ReportFilter.pyFilterName`/`pyFilterValue`/`pyFilterOperation`; logika `Embed-ReportFilters.pyFilterLogic` (ronde 2 menyisir `pyObjClass`/`pyColumnName`/`pyCriteriaValue` — medan yang salah) |

## 1. Halaman awal — pintu masuk modul

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Section/GridRetrocessionLife.xml` b1082 teks **`MASTER CONTRACT RETRO LIFE`** (format `Heading 1`), b2502 menyertakan section `InputRetrocessionLife` | tombol menu **Master Contract Retro Life** (`M_NAV_MENU.LABEL` isi awal 900) → halaman `mcrl-tahun` berjudul `MASTER CONTRACT RETRO LIFE` | ✅ paket 9+10 |
| ⭐ **Bukti bahwa tidak satu pun dari keempat harness adalah pintu masuk**: `InboxRetroLimitReinsurers` dibuka tombol `ReinsType` (`InputRetrocessionLife.xml` b12725 `showHarness` popup, b12734); `InboxBusinessLifeReinsurers` dan `InboxRetroLifeReinsurersList` dibuka tombol `Business List`/`Reinsurer List` (`InputRetroLimitReinsurers.xml` b13590/b13599, b14672/b14681); `InboxSecurityReinsurerLife` dibuka `Security Reinsurer` (`InputSecurityLifeReinsurers.xml` b13005/b13014). `GridRetrocessionLife` tidak dirujuk rule mana pun di korpus modul (grep nama: hanya berkasnya sendiri) — ia dibuka portal/navigasi di luar ekspor | halaman awal = `GridRetrocessionLife`; keempat harness = panel yang dibuka **dari dalam** | ralat spec §4 *"InboxRetroLifeReinsurersList … titik masuk"* (RALAT R1) |
| Harness `pxIconHistory` b519 · `pxIconAttachments` b576 · `pxIconExpandCollapse` b634 (`InboxRetroLifeReinsurersList.xml`; sama di ketiga harness lain, bergeser beberapa baris) | — | ➖ ikon kerangka harness Pega tanpa aksi modul (nol `pyActionSets`) |
| Harness `pxIconCancel` (`<pyFormat>` `InboxRetroLimitReinsurers.xml` b694, `InboxRetroLifeReinsurersList.xml` b692, `InboxSecurityReinsurerLife.xml` b697, `InboxBusinessLifeReinsurers.xml` b691) — kontrol bawaan penutup harness popup | ikon tutup panel (`components/Bingkai.tsx` `KepalaPanel`), keterangan `Close` `[tidak ada di korpus]` | ✅ paket 9 |

## 2. Tahun treaty — `GridRetrocessionLife` → `InputRetrocessionLife` (+ form `InputDtlRetrocessionLife`)

Tabel `POOLDATA.TREATYYEAR_LIFE` (`ID`, `TREATYYEAR`, `UNDERWRITINGYEAR`, `USERID`, `TGLUPDATE`, `STARTDATE`, `ENDDATE`). Komponen `pages/MasterContractRetroLife.tsx`.

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Grid `REPEATING` b9443, sumber RD `BrowseTreatyYear_Life_RD` (b13483), 10 baris/halaman b9441; tanpa parameter RD | — | RD kelas `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`, saringan `.ID = Param.ID` b589 (param kosong), urut `.ID ASC` b768 | `GET /api/master-contract-retro-life/tahun` | ✅ paket 1 (baca) / 9 (layar) |
| Kolom grid `ID` b9588 · `UNDERWRITING YEAR` b9750 · `TRANSACTION YEAR` b9922 · `START DATE` b10094 · `END DATE` b10266 → `.ID` b10724 · `.UNDERWRITINGYEAR` b10921 · `.TREATYYEAR` b11121 · `.STARTDATE` b11321 · `.ENDDATE` b11521 | — | `TREATYYEAR` = label **TRANSACTION YEAR** | idem | ✅ paket 9 |
| ⚠️ Wadah grid + `End Period`: `pyContainerVisibleWhen` `IsFire` b8686 | — | When `IsFire` tidak ada di korpus modul | — | ✅ ikut XML: `pyIsVisibilityOption = ALWAYS` b8727 → kondisi **diabaikan**, grid selalu tampil (RALAT R6) |
| Tombol (baris baru): label sel **`End Period`** b8927 + teks tombol **`Add`** b9007, tooltip `Add New Data` b9005 (RALAT R14) | `NewInputTreatyYear_Life_Act` b9137: 1 b251 `·` Page-New `RetrocessionLife`; 2 b395 `·` kosongkan `InputTreatyYear.*`, `DATASHOW3=1`, `DATASHOW2=0` | — | form baru (klien) | ✅ paket 9 |
| Tombol **`Edit`** b11774 | `SetTreatyYearLife_Act` b11905: 1 b360 `·` `DATASHOW2=1`, `DATASHOW3=1`, `TempEdit.CARI1=1`; 2 b558 `·` salin `ID, TREATYYEAR, UNDERWRITINGYEAR, STARTDATE, ENDDATE` + `USERID=OperatorID`, `TGLUPDATE=@CurrentDateTime()` | — | form dari baris (klien) | ✅ paket 9 |
| Tombol **`ReinsType`** b12117 | `showHarness` popup `InboxRetroLimitReinsurers` b12734 + `SetValueRetroLimit_TreatyYearLife` 1 b415 `·` (salin konteks `IDTREATYYEAR`, `TREATYYEAR`, `TREATYSTARTDATE←.STARTDATE`, `TREATYENDDATE←.ENDDATE`), 2 b809 `·` `STSSAVE=""`, `DATASHOW=""` | — | `GET …/tahun/{id}/kontrak` | ✅ paket 3 / 9 (panel kontrak) |
| Form `Input New Data` (`InputDtlRetrocessionLife.xml` b892), wadah b348 `DATASHOW3 = 1` (`CONDITION`) | — | — | — | ✅ paket 9 |
| Medan `ID` b1834 (ro, `NOTBLANK`) · `UNDERWRITING YEAR` b2018 · `TRANSACTION YEAR` b2277 · `START DATE` b2522 · `END DATE` b2812 · `Modified Date` b3633 (ro) · `Inputor` b3819 `OperatorID.pyUserName` bila `DATASHOW2 = 0` / b4035 `InputTreatyYear.USERID` bila `DATASHOW2 = 1` | — | — | — | ✅ paket 9 |
| Tombol **`Save`** b5059 | `SaveTreatyYearLife_Act` b5167 (refresh `InputRetrocessionLife`): 1 b291 `·` errmsg `"Value cannot be empty."`; 2 b420 `·` Page-Clear-Messages; 3 b512 `·` PRE=true b619 `UNDERWRITINGYEAR==""‖TREATYYEAR==""‖STARTDATE==""‖ENDDATE==""` T=2 F=3, trans b540 `1==1`→6; 4 b666 `·` CARI1..7 (tanggal `dd/MM/yyyy`); 5 b934 `·` PRE=true `TempEdit.CARI1==1` (format ulang saat Edit); 6 b1126 `·` RDB-List `SaveMasterTreatyYear_Life_SQL`; 7 b1292 `·` `DATASHOW3=""` | `SaveMasterTreatyYear_Life_SQL` b102 → `POOLDATA.INSERTTREATYYEAR_LIFE` (upsert `ID`, `ID` baru `'1'‖LPAD(TREATYYEAR_LIFE_SEQ,6)`, `TGLUPDATE=SYSDATE` `[data DBA]`) — **tidak dipanggil** (keputusan o), ditiru di Go | `POST …/tahun` · `PUT …/tahun/{id}` | ✅ paket 2 |
| Tombol **`Cancel`** b5336 | `CancelActivityTreatyContract` b5479: 1 b250 `·` `DATASHOW3=""`, `HASILD7=""`, `DATASHOW=""` + refresh | — | tutup form (klien) | ✅ paket 9 |
| `OutputParam.ERRMSG` b7582 (wadah b7320 `STSSAVE = 1`) · `ERRMSG3` b8243 (wadah b8001 `STSSAVE = 2`) | — | tidak ada activity modul yang mengeset `STSSAVE` 1/2 | — | ➖ tak terjangkau (nol penulis `STSSAVE=1/2`) |
| Hapus tahun treaty | — | **nol** penghapus `TREATYYEAR_LIFE` (5 penulis, 4 penghapus) | **nol** rute hapus | ✅ ikut XML (tahun abadi) — dijaga uji |
| Grid kolom-filter `BrowseTreatyYear_RD` b10415 (kelas non-life `TREATYYEAR`) | — | metadata saring kolom, bukan sumber grid | — | ➖ |

## 3. Kontrak dan batas proteksi — `InboxRetroLimitReinsurers` → `InputRetroLimitReinsurers`

Tabel `POOLDATA.TREATYCONTRACT_LIFE`. Komponen `components/PanelKontrak.tsx` (judul VERBATIM `Reins Type` b666).

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kepala `Reins Type` b666 · `ID Treaty Year` b1093 (ro `InputTreatyContract.IDTREATYYEAR`) | — | — | — | ✅ paket 9 |
| Grid `REPEATING` b9231, RD `BrowseTreatyContract_Life_RD` (b15718), param `IDTREATYYEAR = InputTreatyContract.IDTREATYYEAR` b9252 | — | RD b633 `.IDTREATYYEAR = Param.IDTREATYYEAR`, b647 `.REINSTYPEID = Param.REINSTYPEID` (tidak dikirim), urut `.ID ASC` b943 | `GET …/tahun/{id}/kontrak` | ✅ paket 1 / 3 |
| Kolom `ID` b9370 · `REINS TYPE` b9519 · `TREATY START` b9681 · `TREATY END` b9843 · `MINIMUM LIMIT (IDR)` b10009 · `MAXIMUM LIMIT (IDR)` b10175 · `MINIMUM LIMIT (USD)` b10341 · `MAXIMUM LIMIT (USD)` b10507 | — | `.REINSTYPENAME`, `.TREATYSTARTDATE`, `.TREATYENDDATE`, `.B_IDR`, `.IDR`, `.B_USD`, `.USD` | idem | ✅ paket 9 |
| Tombol **`Add`** b8711 | `NewInputTreatyLimit_Life` b8841: 1 b270 `·` `ProtectTreatyType.CARI1=0`; 2 b404 `·` `DATASHOW=1`, `DATASHOW2=0`, kosongkan `ID, REINSTYPEID, REINSTYPENAME, IDR, USD, B_IDR, B_USD`, `USERID=OperatorID` | — | form baru (klien) | ✅ paket 9 |
| Tombol **`Edit`** b12738 | `SetRetroListLife_Act` b12893: 1 b423 `·` `DATASHOW=1`, `DATASHOW2=1`; 2 b599 `·` salin `REINSTYPEID, REINSTYPENAME, IDR, USD, B_IDR, B_USD, ID` (`TREATYSTARTDATE`/`ENDDATE` param **kosong**) | — | form dari baris (klien) | ✅ paket 9 |
| Tombol **`Business List`** b13113 | `showHarness` `InboxBusinessLifeReinsurers` b13599 + `SetValueRetroLimit_TreatyYearLife` (`IDTREATYCONTRACT←.ID`, `REINSTYPEID/NAME`) | — | `GET …/kontrak/{id}/business` | ✅ paket 6 / 9 |
| Tombol **`Reinsurer List`** b14148 | `showHarness` `InboxRetroLifeReinsurersList` b14681 + `runActivity CountingPercentShare_Act` b15080 (`TREATYYEARID←.IDTREATYYEAR`, `TREATYCONTRACTID←.ID`) | `GetMasterReinsurerLifeList_SQl` b85 `SELECT PCTSHARE … WHERE TREATYYEARID = … AND TREATYCONTRACTID = …` | `GET …/kontrak/{id}/reinsurer` (memuat total) | ✅ paket 4 / 9 |
| Tombol **`Delete`** b15275 (ikon `toolbar_close.gif`; `pyLabelFieldValue` `Button` b15200 adalah label sel, bukan teks tombol) | `DeleteTreatyLimit_Act` b15386: 1 b253 `·` kosongkan halaman; 2 b629 `·` PRE=true `Param.IDTREATYCONTRACT!=""` F=3 RDB-List `DeleteTreatyLimit_SQL`; 3 b817 `·` bila kosong `STSSAVE=3`, `ERRMSG3="Data Gagal di Hapus"`; 4 b983 `·` bila terisi `STSSAVE=3`, `ERRMSG3="Data Berhasil di Hapus"` | `DeleteTreatyLimit_SQL` b84 `delete from POOLDATA.treatycontract_life where ID = …` (datar) | `GET …/kontrak/{id}/dampak-hapus` · `DELETE …/kontrak/{id}` | ✅ paket 7 — kaskade K2 + popup (penyimpangan 3) |
| Form (wadah b1647 `DATASHOW =1`, `CONDITION`): `REINS TYPE` dropdown b3278 | onChange `TreatyLimit_TypeProtect` b3532: 1 b255 `·` `CARI1←REINSTYPEID`; 2 b384 `·` RD `BrowseReinsuranceTypeLimit_RD`; 3 b557 `·` pxRetrieveReportData; 4 b782 `·` ulang `ReinsuranceType.pxResults`; 4.1 b782 `·` PRE=true `CARI1==.ID` T=2 F=3 → `REINSTYPENAME ← .Note` | RD b581 `.Flag = 1`, urut `.ID DESC` b768; dropdown nilai `.ID`, tampil `.Note` (`pyPrompt`) | `GET …/jenis-reasuransi` | ✅ paket 1 (baca) / 3 |
| `TREATY START` b3642 / `TREATY END` b3932 — **read-only** | nilai dari `SetValueRetroLimit_TreatyYearLife` b485/b507 (`← .STARTDATE/.ENDDATE` tahun) | — | diisi server dari tahun induk | ✅ paket 3 (salinan K4) |
| `MINIMUM LIMIT (IDR)` b4222 (`B_IDR`) · `MAXIMUM LIMIT (IDR)` b4486 (`IDR`) · `MINIMUM LIMIT (USD)` b4709 (`B_USD`) · `MAXIMUM LIMIT (USD)` b4974 (`USD`) · `Inputor` b2108/b2308 · `Modified Date` b2537 | — | — | — | ✅ paket 9 |
| Tombol **`Save`** b5889 | `SaveTreatyLimit_Act` b6009: 1 b291 `·` errmsg `"All value cannot be empty."`; 3 b512 `·` PRE=true b619 `REINSTYPEID‖TREATYSTARTDATE‖TREATYENDDATE‖B_IDR‖IDR‖B_USD == ""` T=2 F=3, trans →6; 4 b676 `·` CARI1..14 (**`CARI13 ← IDR_SELISIH`** b1001, **`CARI14 ← USD_SELISIH`** b1023); 5 b1105 `·` RDB-List `SaveMasterTreatyContract_Life_SQL`; 6 b1277 `·` `DATASHOW3=""`, `DATASHOW=""` | `SaveMasterTreatyContract_Life_SQL` b116 → `INSERTTREATYCONTRACT_LIFE`; **nol penulis** `IDR_SELISIH`/`USD_SELISIH` di korpus modul (grep `SELISIH`: hanya pembaca) | `POST …/tahun/{id}/kontrak` · `PUT …/kontrak/{id}` | ✅ paket 3 — selisih dihitung Go (K5) |
| Tombol **`Cancel`** b6176 (ikon `toolbar_close.gif`, aksi `pyBehaviors` lama b6200) | `CancelActivityTreatyContract` | — | tutup form (klien) | ✅ paket 9 |
| `ERRMSG` b7131 (wadah b6959 `STSSAVE = 4`) | — | nol penulis `STSSAVE=4` | — | ➖ tak terjangkau |
| `ERRMSG3` b7876 (wadah b7705 `STSSAVE = 3`) | pesan hapus kontrak | — | pesan sesudah hapus | ✅ paket 7 / 9 |
| Kolom-filter `BrowseRetrocessionLife_RD` (`InputRetroLimitReinsurers.xml` b9348/b10976/b11138) (kelas `RETROCESSIONLIFE`) | — | metadata saring kolom | — | ➖ |

## 4. Reinsurer — `InboxRetroLifeReinsurersList` → `InputSecurityLifeReinsurers`

Tabel `POOLDATA.TREATYREINSURER_LIFE`. Komponen `components/PanelReinsurer.tsx` (judul VERBATIM `Reinsurer List` b637).

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kepala `Reinsurer List` b637 · `ID Treaty Year` b1066 · `ID Reins Type` b1280 (**nilainya `InputBusinessLife.TREATYCONTRACTID`** — label VERBATIM, isinya ID kontrak) · `Reins Type` b1483 | — | — | — | ✅ paket 9 |
| Grid b9334, RD `BrowseDetailTreatyReisurerLife_RD` (b13920), param `TREATYCONTRACTID` b9366, `TREATYYEARID` b9367 | — | RD b614/b632 saring dua kolom, urut `.ID DESC` b899 | `GET …/kontrak/{id}/reinsurer` | ✅ paket 1 / 4 |
| Kolom `ID` · `REINSURER NAME` · `(%) SHARE` · `(%) DISCOUNT` · `(%) OVR COMM` · `INPUTOR` · `UPDATE DATE` (b9481–b10408) → `.ID, .REINSURERNAME, .PCTSHARE, .COMMISION, .OVR_COMM, .USERID, .TGLUPDATE` | — | `COMMISION` = label **(%) DISCOUNT** | idem | ✅ paket 9 |
| `Total Share -->>` b14404 + `InputSecurityLife.STDRATING` b14547 (ro) | `CountingPercentShare_Act`: 3.1 b727 `·` `TotalShare += @toDecimal(.PCTSHARE)`; 4 b900 `·` `STDRATING ← TotalShare` — nol perbandingan dengan 100 | — | `totalShare` + `totalBukan100` di jawaban daftar | ✅ paket 4 — nama jujur (bukan `STDRATING`); tanda mencolok ≠ 100 (tiket 05) |
| Tombol **`Add`** b8797 | `NewInputSecurityLife_Act` b8902: 1 b235 `·` Page-New `InputSecurityLife`; 2 b379 `·` kosongkan `NAME, REINSURERID_LIFE, PCTSHARE`, `DATASHOW=1`, `DATASHOW2=0`; 3 b660 `·` Call `CountingPercentShare_Act` | — | form baru (klien) | ✅ paket 9 |
| Sel paginator `pyGridPaginator` b9022 berlabel mode `Tambah` | `NewTreatyReinsurerDetail_Act` b9173 + pre-DT `SetOutputParam_DT` b9094 (`DATASHOW="tambah"`, `ERRMSG=""`): mengosongkan halaman **non-life** `InputTreatyReinsurer*`, `STSSAVE=101` | — | — | ➖ paginator dibangun sebagai penomoran halaman; label sisa `Tambah` dan aksinya tidak dibawa (efek bersih: menyembunyikan form/pesan) — OQ-MCRL-11 |
| Tombol **`Edit`** b12216 | `SetSecurityLife_Act` b12353: 1 b339 `·` `DATASHOW=1`, `DATASHOW2=1`; 2 b515 `·` salin `NAME←.REINSURERNAME, REINSURERID_LIFE←.REINSURERID, PCTSHARE, COMMISION, OVR_COMM, ID` | — | form dari baris (klien) | ✅ paket 9 |
| Tombol **`Security Reinsurer`** b12555 | `showHarness` `InboxSecurityReinsurerLife` b13014 + `SetValueRetroLimit_TreatyYearLife` (`IDREINSURER←.ID`, `REINSURERNAME`, `PCTSHARE`) | — | `GET …/reinsurer/{id}/security` | ✅ paket 5 / 9 |
| Tombol **`Delete`** b13532 | `DeleteSecurityLife_Act` b13639: 2 b514 `·` PRE=true `Param.ID!=""` F=3 RDB-List `DeleteSecurityReinsurer_SQL`; 3 b702 `·` `STSSAVE=99`, `ERRMSG6="Data Gagal di Hapus"`; 4 b868 `·` `STSSAVE=100`, `ERRMSG="Data Berhasil di Hapus"`; 5 b1034 `·` Call `CountingPercentShare_Act` | `DeleteSecurityReinsurer_SQL` b85 `DELETE FROM POOLDATA.TREATYREINSURER_LIFE WHERE ID = …` | `GET …/reinsurer/{id}/dampak-hapus` · `DELETE …/reinsurer/{id}` | ✅ paket 7 — kaskade ke security |
| Form (wadah b2018 `DATASHOW =1`): `REINSURER NAME` autocomplete b3782 | RD `BrowseCedingCoLife_RD` b3846 kelas `ASM-FW-GISFW-Int-AGENT`, tampil `.ClientName`, **isi `REINSURERID_LIFE ← .ID`** | RD b558 `B AND A AND C`: `.ID Contains "L0"` b582, `.ClientName Contains Param.CedingCoLeader` b587, `.StatusActive = 1` b603; urut `.ClientName ASC` b747 | `GET …/master-reinsurer?cari=` | ✅ paket 1 / 4 |
| `REINS ID` b4607 | — | — | — | ➖ sel `OTHER` `1 = 2` b4764 — **mati**; nilainya diisi autocomplete |
| `(%) SHARE` b4821 · `(%) DISCOUNT` b5111 · `(%) OVR COMM` b5403 | onChange `SetErrorMessageReinsurer` b5017/b5308/b5600: 1 b243 `·` koma→titik; 2 b398 `·` PRE=true `@toDecimal(InputTreatyReinsurer.PctShare)>100 ‖ <0` → `SetErrorMessageBetween`; 3 b577 `·` idem `Ricomm` | ⚠️ memeriksa halaman **`InputTreatyReinsurer`** (kelas non-life `…-TREATYREINSURER`), bukan `InputSecurityLife` yang disunting | gerbang 0..100 di Go | ✅ paket 4 — maksud gerbang ditegakkan untuk `PCTSHARE`, `COMMISION` (preseden OQ-TCO-17); `OVR_COMM` tidak diperiksa Pega (RALAT R3, OQ-MCRL-03) |
| `Inputor` b2539 (`DATASHOW2 = 0`) / b2892 (`= 1`) | — | — | — | ✅ paket 9 |
| Tombol **`Save`** b5770 | `SaveSecurityLife_Act` b5883: 1 b277 `·` `"All value cannot be empty."`; 3 b498 `·` PRE=true b605 `NAME‖REINSURERID_LIFE‖PCTSHARE‖COMMISION‖OVR_COMM == ""` T=2 F=3; 4 b652 `·` **PRE=false** (gerbang `InputRetrocessionLife.*` mati → langkah selalu jalan) CARI1..12 (`TREATYYEARID/CONTRACTID/REINSTYPEID/NAME ← InputBusinessLife.*`); 5 b1085 `·` PRE=false RDB-List `SaveMasterTreatyReinsurer_Life_SQL`; 6 b1337 `·` PRE=false `DATASHOW=""` | `SaveMasterTreatyReinsurer_Life_SQL` b112 → `INSERTREINSURER_LIFE` | `POST …/kontrak/{id}/reinsurer` · `PUT …/reinsurer/{id}` | ✅ paket 4 |
| Tombol **`Cancel`** b6050 (aksi lama b6074) | `CancelActivityTreatyContract` | — | tutup form | ✅ paket 9 |
| `ERRMSG` b7221 (wadah b7049 `STSSAVE = 100`) · `ERRMSG6` b7969 (wadah b7797 `STSSAVE = 99`) | pesan hapus | — | pesan sesudah hapus | ✅ paket 7 / 9 |

## 5. Security reinsurer — `InboxSecurityReinsurerLife` → `InputSecurityReinsurerLife`

Tabel `POOLDATA.TREATYSECURITYREINSURER_LIFE`. Komponen `components/PanelSecurity.tsx` (judul VERBATIM `Security Reinsurer` b659).

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kepala: `ID Reinsurer` b1088 · `Reinsurer Name` b1306 · `PCT Share (%)` b1522 (ro, reinsurer induk) | — | — | `induk` di jawaban daftar | ✅ paket 5 / 9 |
| Grid b8824, RD `BrowseSecurityReinsurer_Life_RD` (b11538), param `TREATYCONTRACTID` b8862 · `TREATYREINSURER` b8863 · `TREATYYEARID` b8864 | — | RD b581/b594/b610 tiga saringan; tanpa urut (`pySortOrder` 99999 semua) | `GET …/reinsurer/{id}/security` | ✅ paket 1 / 5 |
| Kolom `ID` · `REINSURER NAME` · `(%) SHARE` · `INPUTOR` · `UPDATE DATE` (b8978–b9589) | — | `.ID, .REINSURERNAME, .PCTSHARE, .USERID, .TGLUPDATE` | idem | ✅ paket 9 |
| ➕ eksposur efektif (share anak × share induk) | — | tidak ada di korpus | `eksposur` di jawaban daftar | ✅ paket 5 — tiket 06 `[fakta bisnis — work owner]`, kolom tanpa padanan XML (OQ-MCRL-08) |
| Tombol **`Add`** b8287 | `NewInputSecurityLife_Act` b8392 — mengosongkan halaman **reinsurer** (`InputSecurityLife`), bukan `InputSecurityReinsurerLife`; `DATASHOW=1` menampilkan form security berisi nilai terakhir | — | form baru **kosong** (klien) | ✅ paket 9 — OQ-MCRL-12 |
| Sel paginator `Tambah` b8506 | seperti §4 | — | — | ➖ OQ-MCRL-11 |
| Tombol **`Edit`** b10893 | `SetSecurityReinsurerLife_Act` b11021: 1 b290 `·` `DATASHOW=1`, `DATASHOW2=1`; 2 b466 `·` salin `REINSURERNAME, REINSURERID, PCTSHARE, ID` | — | form dari baris | ✅ paket 9 |
| Tombol **`Delete`** b11214 | `DeleteSecurityReinsurerLife_Act` b11321: 2 b453 `·` PRE=true `Param.ID!=""` RDB-List `DeleteSecurityReinsurerLife_SQL`; 3 b641 `·` `STSSAVE=99` `"Data Gagal di Hapus"`; 4 b807 `·` `STSSAVE=100` `"Data Berhasil di Hapus"`; 5 b973 `·` Call `CountingPercentShare_Act` | `DeleteSecurityReinsurerLife_SQL` b85 `DELETE FROM POOLDATA.TREATYSECURITYREINSURER_LIFE WHERE ID = …` | `DELETE …/security/{id}` (+ dampak-hapus) | ✅ paket 7 |
| Form (wadah b2084 `DATASHOW =1`): `SECURITY REINSURER NAME` autocomplete b3856 | RD `BrowseCedingCoLife_RD` b3928, isi `REINSURERID ← .ID` | seperti §4 | `GET …/master-reinsurer?cari=` | ✅ paket 5 |
| `REINS ID` b4683 | — | — | — | ➖ sel `OTHER` `1 = 2` — mati |
| `(%) SHARE` b4898 | onChange `SetErrorMessageReinsurer` b5094 (halaman non-life, lihat §4) | — | gerbang 0..100 di Go | ✅ paket 5 (OQ-MCRL-03) |
| Tombol **`Save`** b5266 | `SaveSecurityReinsurerLife_Act` b5377: 1 b262 `·` `"All value cannot be empty."`; 3 b483 `·` PRE=true b590 `REINSURERNAME‖REINSURERID‖PCTSHARE == ""` T=2 F=3; 4 b637 `·` PRE=false CARI1..8 (`TREATYYEARID/CONTRACTID/REINSURERID ← InputSecurityReinsurer.*`); 5 b982 `·` PRE=false RDB-List `SaveMasterTreatySecurityReinsurer_Life_SQL`; 6 b1234 `·` PRE=false `DATASHOW=""` | `SaveMasterTreatySecurityReinsurer_Life_SQL` b106 → `INSERTSECURITYREINSURER_LIFE` | `POST …/reinsurer/{id}/security` · `PUT …/security/{id}` | ✅ paket 5 |
| Tombol **`Cancel`** b5537 (aksi lama b5561) | `CancelActivityTreatyContract` | — | tutup form | ✅ paket 9 |
| `ERRMSG` b6711 (`STSSAVE = 100`) · `ERRMSG6` b7459 (`STSSAVE = 99`) | pesan hapus | — | — | ✅ paket 7 / 9 |

## 6. Business — `InboxBusinessLifeReinsurers` → `InputBusinessLifeReinsurers`, dan tampilan rate

Tabel `POOLDATA.TREATYBUSINESS_LIFE`. Komponen `components/PanelBusiness.tsx` (judul VERBATIM `Business List` b684), `components/ModalRate.tsx` (judul `Rate List`).

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kepala `Business List` b684 · `ID Treaty Year` b1103 · `ID Reins Type` b1320 (isi ID kontrak) · `Reins Type` b1535 | — | — | — | ✅ paket 9 |
| ⚠️ Wadah `Add` + grid b8980 `InputData.HASILD10 = 2` | — | nol penulis `HASILD10` di korpus modul | — | ✅ ikut XML: `pyIsVisibilityOption = ALWAYS` b8344 → kondisi diabaikan, selalu tampil (RALAT R6) |
| Grid b8980, RD `BrowseTreatyBusiness_Life_RD` (b12529), param `TREATYCONTRACTID` b9012, `TREATYYEARID` b9013 | — | RD b617/b629; urut `.TGLUPDATE ASC` b941 | `GET …/kontrak/{id}/business` | ✅ paket 1 / 6 |
| Kolom `BUSINESS NAME` b9128 · `R/I RATE` b9286 · `INPUTOR` b9440 · `UPDATE DATE` b9595 | — | `.BIZNAME, .RIRATE, .USERID, .TGLUPDATE` | idem | ✅ paket 9 |
| Tombol **`Add`** b8599 | `NewInputBusinessLife_Act` b8705: 1 b270 **`//`** Page-New (ter-remark, tidak jalan); 2 b414 `·` kosongkan `ID, BIZCODE, BIZNAME, RIRATEID, RIRATE`, `TREATYYEARID ← InputTreatyContract.IDTREATYYEAR`, `USERID`, `TGLUPDATE`, `DATASHOW=1`, `DATASHOW2=0` | — | form baru | ✅ paket 9 |
| Tombol **`Edit`** b10938 | `SetBusinessListLife_Act` b11072: 1 b327 `·` `DATASHOW=1`, `DATASHOW2=1`; 2 b503 `·` salin `ID←IDPEGA, BIZCODE, BIZNAME, RIRATEID, RIRATE` | — | form dari baris | ✅ paket 9 |
| Tombol **`Delete`** b11264 | `setValue` b11386 `STSSAVE=""` lalu `DeleteRowBusiness` b11414: 1 b284 `·` `InputData.HASIL12 ← Param.ID`; 2 b414 `·` RDB-List `DeleteRowBusinessList`; 3 b591 `·` `STSSAVE=100`, `ERRMSG = "Data Dengan ID" + " " + HASIL12 + " " + "Berhasil di Hapus"` | `DeleteRowBusinessList` b84 `delete from pooldata.treatybusiness_life where id={InputData.HASIL12}` — `HASIL12` = penampung **ID**, bukan penampung galat (Pertanyaan D terjawab, RALAT R8) | `DELETE …/business/{id}` | ✅ paket 7 |
| Tombol **`View Rate`** (baris) b11581 | `runActivity SetParamRateTable` b11835 (`ParamID.RIRATEID ← .RIRATEID`) + `localAction ViewRateTable` b11909 | FlowAction `ViewRateTable` b96 menampilkan section `ViewRate` | `GET …/rate?idusedby=` | ✅ paket 6 / 9 — data view `RATE_LIFE` baca saja sejak K1 keputusan work owner 01-10-2026 (OQ-MCRL-13 ditutup) |
| Tombol **`Copy to all Reinstype`** b12135 | `SaveBusinessToAllLife_Act` b12264 (param `BIZCODE, BIZNAME, RIRATEID, RIRATE, REINSTYPEID ← .REINSTYPEID`): 1 b302 `·` `CARI1 ← InputBusinessLife.TREATYYEARID`; 2 b431 `·` RDB-List `GetTreatyContract_life`; 3 b648 `·` PRE=false ulang `TempTreatyContract.pxResults`; 3.1 b648 `·` PRE=true b915 `.REINSTYPEID==Param.REINSTYPEID` **T=3 F=2** (kontrak berjenis SAMA dilewati); 3.2 b955 `·` PRE=true b1097 idem → RDB-List `SaveTreatyBusinessAll_Life_SQL`; 4 b1189 `·` `STSSAVE=100`, `ERRMSG="Copied to all reins types."` | `GetTreatyContract_life` b84 `select * from pooldata.treatycontract_life where idtreatyyear = …`; `SaveTreatyBusinessAll_Life_SQL` b85 `INSERT` baris baru `'1'‖lpad(TREATYBUSINESS_LIFE_SEQ,6)`, `TGLUPDATE=SYSDATE`, **tanpa kolom `TREATYYEAR`** | `GET …/business/{id}/salin-semua` (pratinjau) · `POST …/business/{id}/salin-semua` | ✅ paket 6 — sasaran = kontrak **lain** setahun berjenis **berbeda** (RALAT R2) |
| Form (wadah b2070 `DATASHOW =1`): `BUSINESS CODE` b3712 (ro) · `BUSINESS NAME` autocomplete b3898 | RD `BrowseBusinessLife_RD` b4028 kelas `ASM-FW-GISFW-Int-BUSINESS`, tampil `.Note`, **isi `BIZCODE ← .ID`**, kolom tampil `.OLDID` | RD b664 `.OLDID StartsWith "L"`; urut `.ID ASC` b1060 | `GET …/master-business?cari=` | ✅ paket 1 / 6 |
| `R/I RATE` autocomplete b4351 | RD `BrowseRateLifeSummary` b4477 kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`, tampil `.USEDBY`, **isi `RIRATEID ← .ID`** | RD b534 `.ID = param.id`, b547 `.USEDBY Contains param.idusedby`; urut `.ID ASC` b693; kedua param autocomplete kosong (b4534 `id`, b4542 `idusedby`) → saringan gugur, kata ketikan dicari `Contains` di `.USEDBY` (medan cari b4494); isi `RIRATEID ← .ID` b4527 — objek fisik = view DEV `RATE_LIFE_SUMMARY` (OQ-MCRL-05 ditutup K1) | `GET …/ringkasan-rate?cari=` | ✅ paket 1 / 6 / 9 — ⭐ `RIRATE` = **nama tabel rate** (teks, RALAT R7) — K1 keputusan work owner 01-10-2026: `SELECT ID, USEDBY` saja, urut `ID ASC`, 100 saran; RIRATEID pilihan **baru** wajib ada di view (penyimpangan sadar — Pega tidak memeriksa) |
| Tombol **`View Rate`** (form) b4849, tampil bila `InputBusinessLife.RIRATEID!=''` | `localAction ViewRate` b5031; FlowAction `ViewRate` b27 pre-processing `SetParamRate` (1 b238 `·` `ParamID.RIRATEID ← InputBusinessLife.RIRATEID`) | FlowAction b91 section `ViewRate` | `GET …/rate?idusedby=` | ✅ paket 6 / 9 — K1 keputusan work owner 01-10-2026 (OQ-MCRL-13 ditutup) |
| Tombol **`Save`** b5846 (aksi hanya di `pyBehaviors` lama b5872; refresh `thisSection`, `pySection = GridTreatyBusiness` — section tidak ada di korpus, diabaikan target `thisSection`) | `SaveBusinessLife_Act`: 1 b271 `·` `"All value cannot be empty."`; 3 b492 `·` PRE=true b599 `BIZCODE‖BIZNAME‖RIRATEID‖RIRATE == ""` T=2 F=3; 4 b646 `·` CARI1..12 (`TREATYYEAR ← InputBusinessLife.TREATYYEAR`); 5 b1017 `·` RDB-List `SaveMasterTreatyBusiness_Life_SQL`; 6 b1189 `·` PRE=false `DATASHOW=""` | `SaveMasterTreatyBusiness_Life_SQL` b112 → `INSERTBUSINESS_LIFE` | `POST …/kontrak/{id}/business` · `PUT …/business/{id}` | ✅ paket 6 |
| Tombol **`Cancel`** b6064 (aksi lama b6088) | `CancelActivityTreatyContract` | — | tutup form | ✅ paket 9 |
| `ERRMSG` b7008 (`STSSAVE = 100`) · `ERRMSG4` b7759 (`STSSAVE = 99`, nol penulis `ERRMSG4`) | pesan hapus/salin | — | — | ✅ paket 7 / 9; `ERRMSG4` ➖ |
| Section `ViewRate` b861 `Rate List`, grid b1040 RD `BrowseRateLife_RD` (b3063), param `idusedby = ParamID.RIRATEID` b1055 *(ralat 01-10-2026: tertulis b1061 — baris struktural)* | — | RD kelas `ASM-FW-GISFW-Int-M_RATE_LIFE` = view `RATE_LIFE` (`Endorsment Fac In/RDBList/BrowseLifeRate_SQL.xml` b85 `… FROM RATE_LIFE WHERE IDUSEDBY = …`); saring b557 `.IDUSEDBY = Param.idusedby`; urut `.ID DESC` b749, `.RATE ASC` b784; `pyMaxRecords` 500 b729 | `GET …/rate?idusedby=` | ✅ paket 6 / 9 — K1 keputusan work owner 01-10-2026: `SELECT ID, USEDBY, GENDER, CONTRACT, AGE, RATE … WHERE IDUSEDBY = :1 ORDER BY ID DESC, RATE ASC`, paling banyak 500 baris + `terpotong` (Pega memotong diam-diam); nol `SELECT *`, nol `JSONDATA`, tulisan ke view ditolak `periksaBacaSaja` |
| Kolom `ID` · `USEDBY` · `GENDER` · `CONTRACT` · `AGE` · `RATE` (b1178–b1908) | — | teks apa adanya | idem | ✅ paket 9 |
| Tombol dialog FlowAction `ViewRate`/`ViewRateTable`: `pyShowFAButtons` true, `Submit` (`ViewRate.xml` b19, `ViewRateTable.xml` b21) · `Cancel` (b20, b22); tanpa post-processing | — | — | tutup dialog (klien) | ✅ paket 9 — keduanya menutup `Rate List` |

## 7. Rule tanpa pemanggil di layar

| Rule | Keadaan |
| --- | --- |
| `Activity/SetParamRate.xml` | ✅ dipanggil FlowAction `ViewRate` b27 (pre-processing) — bukan yatim |
| `ReportDefinition/BrowseTreatyYear_RD.xml` (kelas `TREATYYEAR`), `BrowseRetrocessionLife_RD.xml` (kelas `RETROCESSIONLIFE`) | ➖ hanya metadata saring kolom grid; bukan sumber data |
| `DataTransform/SetOutputParam_DT.xml` (memo `not used`) | ➖ dirujuk hanya dari sel paginator `Tambah` (§4, §5) — memo tidak dipakai sebagai alasan; alasannya aksi paginator itu (OQ-MCRL-11) |
| `Section/GridTreatyBusiness` | ➖ dirujuk `InputBusinessLifeReinsurers.xml` b5874 tetapi tidak ada di korpus; target refresh `thisSection` membuatnya tak berpengaruh |
| Laporan kontrak total share ≠ 100 (tiket 11) | ➕ **tidak ada di Pega** (nol rule mencari kontrak bercelah) → rute API baca saja, tanpa layar/tombol (OQ-MCRL-07) |

## 8. Cocokkan tiket 01–11 dengan tabel ini — tiket yang meleset

| Tiket | Meleset | Ralat |
| --- | --- | --- |
| 01 | grid/tombol tahun tidak disebut; label tombol baru `End Period` | ralat bertanggal di tiket (bukan pemblokir) |
| 02 | tanggal kontrak **read-only** dari tahun; selisih tanpa penulis | ralat bertanggal |
| 03 | gerbang mati (T1) — K7 | ralat bertanggal `[menunggu OQ-MCRL-01]` |
| 05 | validasi 0..100 memeriksa halaman non-life | ralat bertanggal (OQ-MCRL-03) |
| 06 | `Add` security mengosongkan halaman lain; kolom eksposur tanpa XML | ralat bertanggal |
| 07 | Pertanyaan A terjawab dari korpus (`RIRATE` = nama tabel rate); tiket tidak lagi macet | ralat bertanggal |
| 08 | ⭐ arah penerapan terbalik; "tidak atomik" gugur oleh keputusan o | ralat bertanggal |
| 09 | FK tidak ada (K1/K2) | ralat bertanggal |
| 10 | procedure tidak dipanggil → `HASIL1`/`o_message` tidak ada (K8) | ralat bertanggal |
| 11 | tidak ada di Pega → rute API saja | ralat bertanggal |
| 12 | dicabut (K1) | ralat bertanggal |

## 9. Catatan layar paket 9 — medan baca-saja yang di Pega membaca halaman lain

| Unsur korpus | Di Pega | Layar | Alasan |
| --- | --- | --- | --- |
| `Inputor` form security saat `Edit` (`InputSecurityReinsurerLife.xml` b2962, `DATASHOW2 = 1`) | `InputSecurityReinsurer.USERID` — halaman kepala reinsurer induk, tidak diisi `SetSecurityReinsurerLife_Act` | `USERID` baris security yang disunting | tampilan saja; `USERID` ditulis server dari pelaku |
| `Modified Date` form business saat `Edit` (`InputBusinessLifeReinsurers.xml` b2993) | `SetBusinessListLife_Act` tidak mengisinya → nilai sisa halaman | `TGLUPDATE` baris yang disunting | tampilan saja; `TGLUPDATE = SYSDATE` saat simpan |
| `Inputor` form baru (`OperatorID.pyUserName`, mis. `InputDtlRetrocessionLife.xml` b3819) | nama tampilan operator | akun pelaku (`X-Pelaku`) | aplikasi belum punya sumber nama tampilan — OQ-MCRL-14 |
| `Modified Date` form tahun/kontrak/business baru dan `Edit` tahun/kontrak | `@CurrentDateTime()` (`NewInputTreatyYear_Life_Act` b512, `SetTreatyYearLife_Act` b693, `NewInputTreatyLimit_Life` b607, `SetRetroListLife_Act` b673, `NewInputBusinessLife_Act` b748) | waktu saat form dibuka | ikut XML |
