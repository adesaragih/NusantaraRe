# Paritas layar dan aksi — Master Contract Retro Life

> Disusun 30-09-2026 (sesi implementasi modul, paket 0, brief `PROMPT-IMPLEMENTASI-MODUL-MASTER-CONTRACT-RETRO-LIFE.md`).
> Satu baris per harness, section, tombol, dan aksi korpus → rule yang dipanggil → RDB/RD → tabel/kolom → rute API → komponen React.
> Nomor baris `bNNN` = baris berkas korpus sesudah `sed -e 's/></>\n</g'` atas `D:\XML\RNM_BRD\Master Contract Retro Life\` (konvensi repo).
> Setiap langkah activity yang dikutip mencetak `pyStepsBlockName` — `·` = kosong (langkah hidup), `//` = ter-remark (tidak pernah jalan).
> Precondition dibaca dari urutan aksi: `WhenTrue 2` = lanjut, `3` = lewati langkah; `PRE=false` = prakondisi dimatikan → langkah **selalu** jalan.
> Keadaan: ✅ dibangun · 🔜 paket N · ➖ sengaja tidak dibawa (bukti di kolom Catatan) · ➕ tambahan tiket tanpa padanan XML (OQ).

## 0. Alat baca dan sensus

| Hal | Isi |
| --- | --- |
| Alat | pengurai expat bernomor baris (`pohon_activity.py`, `pohon_layar.py`, `baca_rd.py`, `sql_rdb.py` — scratchpad sesi, tidak di repo) |
| Uji alat | `CountingPercentShare_Act` = 5 langkah termasuk 3.1 (cocok `TAMBAHAN-SPEC-ronde-2.md` T3); `TreatyLimit_TypeProtect` = 5 langkah termasuk 4.1 (cocok T3); langkah ber-`//` = 1 di seluruh modul (cocok `grilling-ronde-2.md` H3 butir 5) |
| Langkah ber-`//` | **1**: `NewInputBusinessLife_Act` langkah 1 b268 `Page-New InputBusinessLife` — 26 activity lain nol |
| Tombol — cara 1 | cacah baris `<pyFormat>pxButton</pyFormat>` di 8 Section: **33** |
| Tombol — cara 2 | pengurai: sel `Embed-Display-Table-Cell` ber-`pyFormat = pxButton`: **33** (sepakat per berkas) |
| Tombol unik | **31** — `Save` b5407 dan `Cancel` b5683 di `InputRetrocessionLife.xml` adalah salinan `InputDtlRetrocessionLife` yang disertakan lewat `pyIncludedRuleXML` b575 |
| Harness vs Section | salinan section di dalam keempat harness **identik** dengan berkas Section-nya (pohon kontrol dibandingkan tanpa nomor baris: beda 0, selain pembungkus include) — tabel ini berpijak pada berkas Section |
| Medan visibilitas | sel: `pyUserData.pyVisible` (`ALWAYS`/`OTHER`/`NOTBLANK`) + `pyCondition`; wadah: `Embed-Harness-Section.pyContainerVisibleWhen` **berlaku hanya bila** `pyIsVisibilityOption = CONDITION` |
| Medan RD | kolom `Embed-ReportListFields.pyFieldName`; saringan `Embed-ReportFilter.pyFilterName`/`pyFilterValue`/`pyFilterOperation`; logika `Embed-ReportFilters.pyFilterLogic` (ronde 2 menyisir `pyObjClass`/`pyColumnName`/`pyCriteriaValue` — medan yang salah) |

## 1. Halaman awal — pintu masuk modul

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Section/GridRetrocessionLife.xml` b1017 teks **`MASTER CONTRACT RETRO LIFE`** (format `Heading 1`), b2475 menyertakan section `InputRetrocessionLife` | tombol menu **Master Contract Retro Life** (`M_NAV_MENU.LABEL` isi awal 900) → halaman `mcrl-tahun` berjudul `MASTER CONTRACT RETRO LIFE` | 🔜 paket 9/10 |
| ⭐ **Bukti bahwa tidak satu pun dari keempat harness adalah pintu masuk**: `InboxRetroLimitReinsurers` dibuka tombol `ReinsType` (`InputRetrocessionLife.xml` b12720 `showHarness` popup, b12734); `InboxBusinessLifeReinsurers` dan `InboxRetroLifeReinsurersList` dibuka tombol `Business List`/`Reinsurer List` (`InputRetroLimitReinsurers.xml` b13585/b13599, b14667/b14681); `InboxSecurityReinsurerLife` dibuka `Security Reinsurer` (`InputSecurityLifeReinsurers.xml` b13000/b13014). `GridRetrocessionLife` tidak dirujuk rule mana pun di korpus modul (grep nama: hanya berkasnya sendiri) — ia dibuka portal/navigasi di luar ekspor | halaman awal = `GridRetrocessionLife`; keempat harness = panel yang dibuka **dari dalam** | ralat spec §4 *"InboxRetroLifeReinsurersList … titik masuk"* (RALAT R1) |
| Harness `pxIconHistory` b487 · `pxIconAttachments` b545 · `pxIconExpandCollapse` b603 · `pxIconCancel` b661 (sama di keempat harness) | — | ➖ ikon kerangka harness Pega tanpa aksi modul (nol `pyActionSets`) |

## 2. Tahun treaty — `GridRetrocessionLife` → `InputRetrocessionLife` (+ form `InputDtlRetrocessionLife`)

Tabel `POOLDATA.TREATYYEAR_LIFE` (`ID`, `TREATYYEAR`, `UNDERWRITINGYEAR`, `USERID`, `TGLUPDATE`, `STARTDATE`, `ENDDATE`). Komponen `pages/MasterContractRetroLife.tsx`.

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Grid `REPEATING` b9375, sumber RD `BrowseTreatyYear_Life_RD` (b13483), 10 baris/halaman b9441; tanpa parameter RD | — | RD kelas `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`, saringan `.ID = Param.ID` b582 (param kosong), urut `.ID ASC` b765 | `GET /api/master-contract-retro-life/tahun` | 🔜 paket 1 (baca) / 9 (layar) |
| Kolom grid `ID` b9513 · `UNDERWRITING YEAR` b9676 · `TRANSACTION YEAR` b9848 · `START DATE` b10020 · `END DATE` b10192 → `.ID` b10651 · `.UNDERWRITINGYEAR` b10845 · `.TREATYYEAR` b11045 · `.STARTDATE` b11245 · `.ENDDATE` b11445 | — | `TREATYYEAR` = label **TRANSACTION YEAR** | idem | 🔜 paket 9 |
| ⚠️ Wadah grid + `End Period` b8631 `pyContainerVisibleWhen = IsFire` | — | When `IsFire` tidak ada di korpus modul | — | ✅ ikut XML: `pyIsVisibilityOption = ALWAYS` b8727 → kondisi **diabaikan**, grid selalu tampil (RALAT R6) |
| Tombol **`End Period`** b8888 (baris baru) | `NewInputTreatyYear_Life_Act` b9123: 1 b249 `·` Page-New `RetrocessionLife`; 2 b393 `·` kosongkan `InputTreatyYear.*`, `DATASHOW3=1`, `DATASHOW2=0` | — | form baru (klien) | 🔜 paket 9 — label VERBATIM `End Period` |
| Tombol **`Edit`** b11645 | `SetTreatyYearLife_Act` b11893: 1 b358 `·` `DATASHOW2=1`, `DATASHOW3=1`, `TempEdit.CARI1=1`; 2 b556 `·` salin `ID, TREATYYEAR, UNDERWRITINGYEAR, STARTDATE, ENDDATE` + `USERID=OperatorID`, `TGLUPDATE=@CurrentDateTime()` | — | form dari baris (klien) | 🔜 paket 9 |
| Tombol **`ReinsType`** b11988 | `showHarness` popup `InboxRetroLimitReinsurers` b12720 + `SetValueRetroLimit_TreatyYearLife` 1 b413 `·` (salin konteks `IDTREATYYEAR`, `TREATYYEAR`, `TREATYSTARTDATE←.STARTDATE`, `TREATYENDDATE←.ENDDATE`), 2 b807 `·` `STSSAVE=""`, `DATASHOW=""` | — | `GET …/tahun/{id}/kontrak` | 🔜 paket 3 / 9 (panel kontrak) |
| Form `Input New Data` (`InputDtlRetrocessionLife.xml` b829), wadah b296 `DATASHOW3 = 1` (`CONDITION`) | — | — | — | 🔜 paket 9 |
| Medan `ID` b1798 (ro, `NOTBLANK`) · `UNDERWRITING YEAR` b1981 · `TRANSACTION YEAR` b2238 · `START DATE` b2484 · `END DATE` b2774 · `Modified Date` b3598 (ro) · `Inputor` b3783 `OperatorID.pyUserName` bila `DATASHOW2 = 0` / b3967 `InputTreatyYear.USERID` bila `DATASHOW2 = 1` | — | — | — | 🔜 paket 9 |
| Tombol **`Save`** b4944 | `SaveTreatyYearLife_Act` b5153 (refresh `InputRetrocessionLife`): 1 b289 `·` errmsg `"Value cannot be empty."`; 2 b418 `·` Page-Clear-Messages; 3 b510 `·` PRE=true b609 `UNDERWRITINGYEAR==""‖TREATYYEAR==""‖STARTDATE==""‖ENDDATE==""` T=2 F=3, trans b538 `1==1`→6; 4 b664 `·` CARI1..7 (tanggal `dd/MM/yyyy`); 5 b932 `·` PRE=true `TempEdit.CARI1==1` (format ulang saat Edit); 6 b1124 `·` RDB-List `SaveMasterTreatyYear_Life_SQL`; 7 b1290 `·` `DATASHOW3=""` | `SaveMasterTreatyYear_Life_SQL` b84 → `POOLDATA.INSERTTREATYYEAR_LIFE` (upsert `ID`, `ID` baru `'1'‖LPAD(TREATYYEAR_LIFE_SEQ,6)`, `TGLUPDATE=SYSDATE` `[data DBA]`) — **tidak dipanggil** (keputusan o), ditiru di Go | `POST …/tahun` · `PUT …/tahun/{id}` | 🔜 paket 2 |
| Tombol **`Cancel`** b5221 | `CancelActivityTreatyContract` b5468: 1 b248 `·` `DATASHOW3=""`, `HASILD7=""`, `DATASHOW=""` + refresh | — | tutup form (klien) | 🔜 paket 9 |
| `OutputParam.ERRMSG` b7513 (wadah b7268 `STSSAVE = 1`) · `ERRMSG3` b8194 (wadah b7949 `STSSAVE = 2`) | — | tidak ada activity modul yang mengeset `STSSAVE` 1/2 | — | ➖ tak terjangkau (nol penulis `STSSAVE=1/2`) |
| Hapus tahun treaty | — | **nol** penghapus `TREATYYEAR_LIFE` (5 penulis, 4 penghapus) | **nol** rute hapus | ✅ ikut XML (tahun abadi) — dijaga uji |
| Grid kolom-filter `BrowseTreatyYear_RD` b10415 (kelas non-life `TREATYYEAR`) | — | metadata saring kolom, bukan sumber grid | — | ➖ |

## 3. Kontrak dan batas proteksi — `InboxRetroLimitReinsurers` → `InputRetroLimitReinsurers`

Tabel `POOLDATA.TREATYCONTRACT_LIFE`. Komponen `components/PanelKontrak.tsx` (judul VERBATIM `Reins Type` b601).

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kepala `Reins Type` b601 · `ID Treaty Year` b1058 (ro `InputTreatyContract.IDTREATYYEAR`) | — | — | — | 🔜 paket 9 |
| Grid `REPEATING` b9178, RD `BrowseTreatyContract_Life_RD` (b15718), param `IDTREATYYEAR = InputTreatyContract.IDTREATYYEAR` b9252 | — | RD b628 `.IDTREATYYEAR = Param.IDTREATYYEAR`, b642 `.REINSTYPEID = Param.REINSTYPEID` (tidak dikirim), urut `.ID ASC` b940 | `GET …/tahun/{id}/kontrak` | 🔜 paket 1 / 3 |
| Kolom `ID` b9300 · `REINS TYPE` b9450 · `TREATY START` b9611 · `TREATY END` b9773 · `MINIMUM LIMIT (IDR)` b9935 · `MAXIMUM LIMIT (IDR)` b10101 · `MINIMUM LIMIT (USD)` b10267 · `MAXIMUM LIMIT (USD)` b10433 | — | `.REINSTYPENAME`, `.TREATYSTARTDATE`, `.TREATYENDDATE`, `.B_IDR`, `.IDR`, `.B_USD`, `.USD` | idem | 🔜 paket 9 |
| Tombol **`Add`** b8596 | `NewInputTreatyLimit_Life` b8829: 1 b268 `·` `ProtectTreatyType.CARI1=0`; 2 b402 `·` `DATASHOW=1`, `DATASHOW2=0`, kosongkan `ID, REINSTYPEID, REINSTYPENAME, IDR, USD, B_IDR, B_USD`, `USERID=OperatorID` | — | form baru (klien) | 🔜 paket 9 |
| Tombol **`Edit`** b12621 | `SetRetroListLife_Act` b12881: 1 b421 `·` `DATASHOW=1`, `DATASHOW2=1`; 2 b597 `·` salin `REINSTYPEID, REINSTYPENAME, IDR, USD, B_IDR, B_USD, ID` (`TREATYSTARTDATE`/`ENDDATE` param **kosong**) | — | form dari baris (klien) | 🔜 paket 9 |
| Tombol **`Business List`** b12999 | `showHarness` `InboxBusinessLifeReinsurers` b13585 + `SetValueRetroLimit_TreatyYearLife` (`IDTREATYCONTRACT←.ID`, `REINSTYPEID/NAME`) | — | `GET …/kontrak/{id}/business` | 🔜 paket 6 / 9 |
| Tombol **`Reinsurer List`** b14034 | `showHarness` `InboxRetroLifeReinsurersList` b14667 + `runActivity CountingPercentShare_Act` b15069 (`TREATYYEARID←.IDTREATYYEAR`, `TREATYCONTRACTID←.ID`) | `GetMasterReinsurerLifeList_SQl` b85 `SELECT PCTSHARE … WHERE TREATYYEARID = … AND TREATYCONTRACTID = …` | `GET …/kontrak/{id}/reinsurer` (memuat total) | 🔜 paket 4 / 9 |
| Tombol **`Delete`** b15161 (ikon `toolbar_close.gif`; `pyLabelFieldValue` `Button` b15200 adalah label sel, bukan teks tombol) | `DeleteTreatyLimit_Act` b15372: 1 b251 `·` kosongkan halaman; 2 b627 `·` PRE=true `Param.IDTREATYCONTRACT!=""` F=3 RDB-List `DeleteTreatyLimit_SQL`; 3 b815 `·` bila kosong `STSSAVE=3`, `ERRMSG3="Data Gagal di Hapus"`; 4 b981 `·` bila terisi `STSSAVE=3`, `ERRMSG3="Data Berhasil di Hapus"` | `DeleteTreatyLimit_SQL` b84 `delete from POOLDATA.treatycontract_life where ID = …` (datar) | `GET …/kontrak/{id}/dampak-hapus` · `DELETE …/kontrak/{id}` | 🔜 paket 7 — kaskade K2 + popup (penyimpangan 3) |
| Form (wadah b1566 `DATASHOW =1`, `CONDITION`): `REINS TYPE` dropdown b3242 | onChange `TreatyLimit_TypeProtect` b3520: 1 b253 `·` `CARI1←REINSTYPEID`; 2 b382 `·` RD `BrowseReinsuranceTypeLimit_RD`; 3 b555 `·` pxRetrieveReportData; 4 b735 `·` ulang `ReinsuranceType.pxResults`; 4.1 b780 `·` PRE=true `CARI1==.ID` T=2 F=3 → `REINSTYPENAME ← .Note` | RD b578 `.Flag = 1`, urut `.ID DESC` b765; dropdown nilai `.ID`, tampil `.Note` (`pyPrompt`) | `GET …/jenis-reasuransi` | 🔜 paket 1 (baca) / 3 |
| `TREATY START` b3604 / `TREATY END` b3894 — **read-only** | nilai dari `SetValueRetroLimit_TreatyYearLife` b485/b507 (`← .STARTDATE/.ENDDATE` tahun) | — | diisi server dari tahun induk | 🔜 paket 3 (salinan K4) |
| `MINIMUM LIMIT (IDR)` b4184 (`B_IDR`) · `MAXIMUM LIMIT (IDR)` b4451 (`IDR`) · `MINIMUM LIMIT (USD)` b4671 (`B_USD`) · `MAXIMUM LIMIT (USD)` b4936 (`USD`) · `Inputor` b2070/b2290 · `Modified Date` b2499 | — | — | — | 🔜 paket 9 |
| Tombol **`Save`** b5773 | `SaveTreatyLimit_Act` b5990: 1 b289 `·` errmsg `"All value cannot be empty."`; 3 b510 `·` PRE=true b609 `REINSTYPEID‖TREATYSTARTDATE‖TREATYENDDATE‖B_IDR‖IDR‖B_USD == ""` T=2 F=3, trans →6; 4 b674 `·` CARI1..14 (**`CARI13 ← IDR_SELISIH`** b1001, **`CARI14 ← USD_SELISIH`** b1023); 5 b1103 `·` RDB-List `SaveMasterTreatyContract_Life_SQL`; 6 b1275 `·` `DATASHOW3=""`, `DATASHOW=""` | `SaveMasterTreatyContract_Life_SQL` b84 → `INSERTTREATYCONTRACT_LIFE`; **nol penulis** `IDR_SELISIH`/`USD_SELISIH` di korpus modul (grep `SELISIH`: hanya pembaca) | `POST …/tahun/{id}/kontrak` · `PUT …/kontrak/{id}` | 🔜 paket 3 — selisih dihitung Go (K5) |
| Tombol **`Cancel`** b6066 (ikon `toolbar_close.gif`, aksi `pyBehaviors` lama b6181) | `CancelActivityTreatyContract` | — | tutup form (klien) | 🔜 paket 9 |
| `ERRMSG` b7082 (wadah b6880 `STSSAVE = 4`) | — | nol penulis `STSSAVE=4` | — | ➖ tak terjangkau |
| `ERRMSG3` b7828 (wadah b7626 `STSSAVE = 3`) | pesan hapus kontrak | — | pesan sesudah hapus | 🔜 paket 7 / 9 |
| Kolom-filter `BrowseRetrocessionLife_RD` b10198/b11815 (kelas `RETROCESSIONLIFE`) | — | metadata saring kolom | — | ➖ |

## 4. Reinsurer — `InboxRetroLifeReinsurersList` → `InputSecurityLifeReinsurers`

Tabel `POOLDATA.TREATYREINSURER_LIFE`. Komponen `components/PanelReinsurer.tsx` (judul VERBATIM `Reinsurer List` b572).

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kepala `Reinsurer List` b572 · `ID Treaty Year` b1029 · `ID Reins Type` b1245 (**nilainya `InputBusinessLife.TREATYCONTRACTID`** — label VERBATIM, isinya ID kontrak) · `Reins Type` b1446 | — | — | — | 🔜 paket 9 |
| Grid b9281, RD `BrowseDetailTreatyReisurerLife_RD` (b13920), param `TREATYCONTRACTID` b9366, `TREATYYEARID` b9367 | — | RD b613/b631 saring dua kolom, urut `.ID DESC` b896 | `GET …/kontrak/{id}/reinsurer` | 🔜 paket 1 / 4 |
| Kolom `ID` · `REINSURER NAME` · `(%) SHARE` · `(%) DISCOUNT` · `(%) OVR COMM` · `INPUTOR` · `UPDATE DATE` (b9415–b10343) → `.ID, .REINSURERNAME, .PCTSHARE, .COMMISION, .OVR_COMM, .USERID, .TGLUPDATE` | — | `COMMISION` = label **(%) DISCOUNT** | idem | 🔜 paket 9 |
| `Total Share -->>` b14339 + `InputSecurityLife.STDRATING` b14482 (ro) | `CountingPercentShare_Act`: 3.1 b725 `·` `TotalShare += @toDecimal(.PCTSHARE)`; 4 b898 `·` `STDRATING ← TotalShare` — nol perbandingan dengan 100 | — | `totalShare` + `totalBukan100` di jawaban daftar | 🔜 paket 4 — nama jujur (bukan `STDRATING`); tanda mencolok ≠ 100 (tiket 05) |
| Tombol **`Add`** b8689 | `NewInputSecurityLife_Act` b8890: 1 b233 `·` Page-New `InputSecurityLife`; 2 b377 `·` kosongkan `NAME, REINSURERID_LIFE, PCTSHARE`, `DATASHOW=1`, `DATASHOW2=0`; 3 b658 `·` Call `CountingPercentShare_Act` | — | form baru (klien) | 🔜 paket 9 |
| Sel paginator `pyGridPaginator` b8956 berlabel mode `Tambah` | `NewTreatyReinsurerDetail_Act` b9159 + pre-DT `SetOutputParam_DT` b9094 (`DATASHOW="tambah"`, `ERRMSG=""`): mengosongkan halaman **non-life** `InputTreatyReinsurer*`, `STSSAVE=101` | — | — | ➖ paginator dibangun sebagai penomoran halaman; label sisa `Tambah` dan aksinya tidak dibawa (efek bersih: menyembunyikan form/pesan) — OQ-MCRL-11 |
| Tombol **`Edit`** b12099 | `SetSecurityLife_Act` b12341: 1 b337 `·` `DATASHOW=1`, `DATASHOW2=1`; 2 b513 `·` salin `NAME←.REINSURERNAME, REINSURERID_LIFE←.REINSURERID, PCTSHARE, COMMISION, OVR_COMM, ID` | — | form dari baris (klien) | 🔜 paket 9 |
| Tombol **`Security Reinsurer`** b12441 | `showHarness` `InboxSecurityReinsurerLife` b13000 + `SetValueRetroLimit_TreatyYearLife` (`IDREINSURER←.ID`, `REINSURERNAME`, `PCTSHARE`) | — | `GET …/reinsurer/{id}/security` | 🔜 paket 5 / 9 |
| Tombol **`Delete`** b13422 | `DeleteSecurityLife_Act` b13625: 2 b512 `·` PRE=true `Param.ID!=""` F=3 RDB-List `DeleteSecurityReinsurer_SQL`; 3 b700 `·` `STSSAVE=99`, `ERRMSG6="Data Gagal di Hapus"`; 4 b866 `·` `STSSAVE=100`, `ERRMSG="Data Berhasil di Hapus"`; 5 b1032 `·` Call `CountingPercentShare_Act` | `DeleteSecurityReinsurer_SQL` b85 `DELETE FROM POOLDATA.TREATYREINSURER_LIFE WHERE ID = …` | `GET …/reinsurer/{id}/dampak-hapus` · `DELETE …/reinsurer/{id}` | 🔜 paket 7 — kaskade ke security |
| Form (wadah b1972 `DATASHOW =1`): `REINSURER NAME` autocomplete b3746 | RD `BrowseCedingCoLife_RD` b3840 kelas `ASM-FW-GISFW-Int-AGENT`, tampil `.ClientName`, **isi `REINSURERID_LIFE ← .ID`** | RD b556 `B AND A AND C`: `.ID Contains "L0"` b565, `.ClientName Contains Param.CedingCoLeader` b584, `.StatusActive = 1` b601; urut `.ClientName ASC` b745 | `GET …/master-reinsurer?cari=` | 🔜 paket 1 / 4 |
| `REINS ID` b4570 | — | — | — | ➖ sel `OTHER` `1 = 2` b4764 — **mati**; nilainya diisi autocomplete |
| `(%) SHARE` b4784 · `(%) DISCOUNT` b5072 · `(%) OVR COMM` b5364 | onChange `SetErrorMessageReinsurer` b4998/b5289/b5581: 1 b241 `·` koma→titik; 2 b396 `·` PRE=true `@toDecimal(InputTreatyReinsurer.PctShare)>100 ‖ <0` → `SetErrorMessageBetween`; 3 b575 `·` idem `Ricomm` | ⚠️ memeriksa halaman **`InputTreatyReinsurer`** (kelas non-life `…-TREATYREINSURER`), bukan `InputSecurityLife` yang disunting | gerbang 0..100 di Go | 🔜 paket 4 — maksud gerbang ditegakkan untuk `PCTSHARE`, `COMMISION` (preseden OQ-TCO-17); `OVR_COMM` tidak diperiksa Pega (RALAT R3, OQ-MCRL-03) |
| `Inputor` b2503 (`DATASHOW2 = 0`) / b2709 (`= 1`) | — | — | — | 🔜 paket 9 |
| Tombol **`Save`** b5656 | `SaveSecurityLife_Act` b5869: 1 b275 `·` `"All value cannot be empty."`; 3 b496 `·` PRE=true b595 `NAME‖REINSURERID_LIFE‖PCTSHARE‖COMMISION‖OVR_COMM == ""` T=2 F=3; 4 b650 `·` **PRE=false** (gerbang `InputRetrocessionLife.*` mati → langkah selalu jalan) CARI1..12 (`TREATYYEARID/CONTRACTID/REINSTYPEID/NAME ← InputBusinessLife.*`); 5 b1083 `·` PRE=false RDB-List `SaveMasterTreatyReinsurer_Life_SQL`; 6 b1335 `·` PRE=false `DATASHOW=""` | `SaveMasterTreatyReinsurer_Life_SQL` b84 → `INSERTREINSURER_LIFE` | `POST …/kontrak/{id}/reinsurer` · `PUT …/reinsurer/{id}` | 🔜 paket 4 |
| Tombol **`Cancel`** b5942 (aksi lama b6055) | `CancelActivityTreatyContract` | — | tutup form | 🔜 paket 9 |
| `ERRMSG` b7172 (wadah b6969 `STSSAVE = 100`) · `ERRMSG6` b7920 (wadah b7717 `STSSAVE = 99`) | pesan hapus | — | pesan sesudah hapus | 🔜 paket 7 / 9 |

## 5. Security reinsurer — `InboxSecurityReinsurerLife` → `InputSecurityReinsurerLife`

Tabel `POOLDATA.TREATYSECURITYREINSURER_LIFE`. Komponen `components/PanelSecurity.tsx` (judul VERBATIM `Security Reinsurer` b594).

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kepala: `ID Reinsurer` b1051 · `Reinsurer Name` b1267 · `PCT Share (%)` b1483 (ro, reinsurer induk) | — | — | `induk` di jawaban daftar | 🔜 paket 5 / 9 |
| Grid b8771, RD `BrowseSecurityReinsurer_Life_RD` (b11538), param `TREATYCONTRACTID` b8862 · `TREATYREINSURER` b8863 · `TREATYYEARID` b8864 | — | RD b579/b592/b609 tiga saringan; tanpa urut (`pySortOrder` 99999 semua) | `GET …/reinsurer/{id}/security` | 🔜 paket 1 / 5 |
| Kolom `ID` · `REINSURER NAME` · `(%) SHARE` · `INPUTOR` · `UPDATE DATE` (b8912–b9524) | — | `.ID, .REINSURERNAME, .PCTSHARE, .USERID, .TGLUPDATE` | idem | 🔜 paket 9 |
| ➕ eksposur efektif (share anak × share induk) | — | tidak ada di korpus | `eksposur` di jawaban daftar | 🔜 paket 5 — tiket 06 `[fakta bisnis — work owner]`, kolom tanpa padanan XML (OQ-MCRL-08) |
| Tombol **`Add`** b8179 | `NewInputSecurityLife_Act` b8380 — mengosongkan halaman **reinsurer** (`InputSecurityLife`), bukan `InputSecurityReinsurerLife`; `DATASHOW=1` menampilkan form security berisi nilai terakhir | — | form baru **kosong** (klien) | 🔜 paket 9 — OQ-MCRL-12 |
| Sel paginator `Tambah` b8446 | seperti §4 | — | — | ➖ OQ-MCRL-11 |
| Tombol **`Edit`** b10779 | `SetSecurityReinsurerLife_Act` b11009: 1 b288 `·` `DATASHOW=1`, `DATASHOW2=1`; 2 b464 `·` salin `REINSURERNAME, REINSURERID, PCTSHARE, ID` | — | form dari baris | 🔜 paket 9 |
| Tombol **`Delete`** b11097 | `DeleteSecurityReinsurerLife_Act` b11307: 2 b451 `·` PRE=true `Param.ID!=""` RDB-List `DeleteSecurityReinsurerLife_SQL`; 3 b639 `·` `STSSAVE=99` `"Data Gagal di Hapus"`; 4 b805 `·` `STSSAVE=100` `"Data Berhasil di Hapus"`; 5 b971 `·` Call `CountingPercentShare_Act` | `DeleteSecurityReinsurerLife_SQL` b85 `DELETE FROM POOLDATA.TREATYSECURITYREINSURER_LIFE WHERE ID = …` | `DELETE …/security/{id}` (+ dampak-hapus) | 🔜 paket 7 |
| Form (wadah b2038 `DATASHOW =1`): `SECURITY REINSURER NAME` autocomplete b3820 | RD `BrowseCedingCoLife_RD` b3922, isi `REINSURERID ← .ID` | seperti §4 | `GET …/master-reinsurer?cari=` | 🔜 paket 5 |
| `REINS ID` b4646 | — | — | — | ➖ sel `OTHER` `1 = 2` — mati |
| `(%) SHARE` b4861 | onChange `SetErrorMessageReinsurer` b5075 (halaman non-life, lihat §4) | — | gerbang 0..100 di Go | 🔜 paket 5 (OQ-MCRL-03) |
| Tombol **`Save`** b5151 | `SaveSecurityReinsurerLife_Act` b5363: 1 b260 `·` `"All value cannot be empty."`; 3 b481 `·` PRE=true b580 `REINSURERNAME‖REINSURERID‖PCTSHARE == ""` T=2 F=3; 4 b635 `·` PRE=false CARI1..8 (`TREATYYEARID/CONTRACTID/REINSURERID ← InputSecurityReinsurer.*`); 5 b980 `·` PRE=false RDB-List `SaveMasterTreatySecurityReinsurer_Life_SQL`; 6 b1232 `·` PRE=false `DATASHOW=""` | `SaveMasterTreatySecurityReinsurer_Life_SQL` b84 → `INSERTSECURITYREINSURER_LIFE` | `POST …/reinsurer/{id}/security` · `PUT …/security/{id}` | 🔜 paket 5 |
| Tombol **`Cancel`** b5432 (aksi lama b5542) | `CancelActivityTreatyContract` | — | tutup form | 🔜 paket 9 |
| `ERRMSG` b6662 (`STSSAVE = 100`) · `ERRMSG6` b7410 (`STSSAVE = 99`) | pesan hapus | — | — | 🔜 paket 7 / 9 |

## 6. Business — `InboxBusinessLifeReinsurers` → `InputBusinessLifeReinsurers`, dan tampilan rate

Tabel `POOLDATA.TREATYBUSINESS_LIFE`. Komponen `components/PanelBusiness.tsx` (judul VERBATIM `Business List` b619), `components/ModalRate.tsx` (judul `Rate List`).

| Unsur korpus | Rule dipanggil (bNNN · blok) | RDB/RD → kolom | Rute API | Keadaan |
| --- | --- | --- | --- | --- |
| Kepala `Business List` b619 · `ID Treaty Year` b1066 · `ID Reins Type` b1282 (isi ID kontrak) · `Reins Type` b1497 | — | — | — | 🔜 paket 9 |
| ⚠️ Wadah `Add` + grid b8260 `InputData.HASILD10 = 2` | — | nol penulis `HASILD10` di korpus modul | — | ✅ ikut XML: `pyIsVisibilityOption = ALWAYS` b8344 → kondisi diabaikan, selalu tampil (RALAT R6) |
| Grid b8927, RD `BrowseTreatyBusiness_Life_RD` (b12529), param `TREATYCONTRACTID` b9012, `TREATYYEARID` b9013 | — | RD b615/b627; urut `.TGLUPDATE ASC` b939 | `GET …/kontrak/{id}/business` | 🔜 paket 1 / 6 |
| Kolom `BUSINESS NAME` b9061 · `R/I RATE` b9220 · `INPUTOR` b9375 · `UPDATE DATE` b9530 | — | `.BIZNAME, .RIRATE, .USERID, .TGLUPDATE` | idem | 🔜 paket 9 |
| Tombol **`Add`** b8484 | `NewInputBusinessLife_Act` b8691: 1 b268 **`//`** Page-New (ter-remark, tidak jalan); 2 b412 `·` kosongkan `ID, BIZCODE, BIZNAME, RIRATEID, RIRATE`, `TREATYYEARID ← InputTreatyContract.IDTREATYYEAR`, `USERID`, `TGLUPDATE`, `DATASHOW=1`, `DATASHOW2=0` | — | form baru | 🔜 paket 9 |
| Tombol **`Edit`** b10824 | `SetBusinessListLife_Act` b11060: 1 b325 `·` `DATASHOW=1`, `DATASHOW2=1`; 2 b501 `·` salin `ID←IDPEGA, BIZCODE, BIZNAME, RIRATEID, RIRATE` | — | form dari baris | 🔜 paket 9 |
| Tombol **`Delete`** b11154 | `setValue` b11381 `STSSAVE=""` lalu `DeleteRowBusiness` b11402: 1 b282 `·` `InputData.HASIL12 ← Param.ID`; 2 b412 `·` RDB-List `DeleteRowBusinessList`; 3 b589 `·` `STSSAVE=100`, `ERRMSG = "Data Dengan ID" + " " + HASIL12 + " " + "Berhasil di Hapus"` | `DeleteRowBusinessList` b84 `delete from pooldata.treatybusiness_life where id={InputData.HASIL12}` — `HASIL12` = penampung **ID**, bukan penampung galat (Pertanyaan D terjawab, RALAT R8) | `DELETE …/business/{id}` | 🔜 paket 7 |
| Tombol **`View Rate`** (baris) b11469 | `runActivity SetParamRateTable` b11824 (`ParamID.RIRATEID ← .RIRATEID`) + `localAction ViewRateTable` b11887 | FlowAction `ViewRateTable` b96 menampilkan section `ViewRate` | `GET …/rate?idusedby=` | 🔜 paket 6 / 9 |
| Tombol **`Copy to all Reinstype`** b12020 | `SaveBusinessToAllLife_Act` b12252 (param `BIZCODE, BIZNAME, RIRATEID, RIRATE, REINSTYPEID ← .REINSTYPEID`): 1 b300 `·` `CARI1 ← InputBusinessLife.TREATYYEARID`; 2 b429 `·` RDB-List `GetTreatyContract_life`; 3 b601 `·` PRE=false ulang `TempTreatyContract.pxResults`; 3.1 b646 `·` PRE=true b905 `.REINSTYPEID==Param.REINSTYPEID` **T=3 F=2** (kontrak berjenis SAMA dilewati); 3.2 b953 `·` PRE=true b1087 idem → RDB-List `SaveTreatyBusinessAll_Life_SQL`; 4 b1187 `·` `STSSAVE=100`, `ERRMSG="Copied to all reins types."` | `GetTreatyContract_life` b84 `select * from pooldata.treatycontract_life where idtreatyyear = …`; `SaveTreatyBusinessAll_Life_SQL` b84 `INSERT` baris baru `'1'‖lpad(TREATYBUSINESS_LIFE_SEQ,6)`, `TGLUPDATE=SYSDATE`, **tanpa kolom `TREATYYEAR`** | `GET …/business/{id}/salin-semua` (pratinjau) · `POST …/business/{id}/salin-semua` | 🔜 paket 6 — sasaran = kontrak **lain** setahun berjenis **berbeda** (RALAT R2) |
| Form (wadah b2022 `DATASHOW =1`): `BUSINESS CODE` b3675 (ro) · `BUSINESS NAME` autocomplete b3862 | RD `BrowseBusinessLife_RD` b4022 kelas `ASM-FW-GISFW-Int-BUSINESS`, tampil `.Note`, **isi `BIZCODE ← .ID`**, kolom tampil `.OLDID` | RD b651 `.OLDID StartsWith "L"`; urut `.ID ASC` b1057 | `GET …/master-business?cari=` | 🔜 paket 1 / 6 |
| `R/I RATE` autocomplete b4314 | RD `BrowseRateLifeSummary` b4471 kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`, tampil `.USEDBY`, **isi `RIRATEID ← .ID`** | RD b529 `.ID = param.id`, b542 `.USEDBY Contains param.idusedby`; urut `.ID ASC` b690 — objek fisik kelasnya tidak terbukti (OQ-MCRL-05) | `GET …/ringkasan-rate?cari=` | 🔜 paket 1 / 6 — ⭐ `RIRATE` = **nama tabel rate** (teks), Pertanyaan A terjawab (RALAT R7) |
| Tombol **`View Rate`** (form) b4732, tampil bila `InputBusinessLife.RIRATEID!=''` | `localAction ViewRate` b5026; FlowAction `ViewRate` b27 pre-processing `SetParamRate` (1 b236 `·` `ParamID.RIRATEID ← InputBusinessLife.RIRATEID`) | FlowAction b91 section `ViewRate` | `GET …/rate?idusedby=` | 🔜 paket 6 / 9 |
| Tombol **`Save`** b5739 (aksi hanya di `pyBehaviors` lama b5854; refresh `thisSection`, `pySection = GridTreatyBusiness` — section tidak ada di korpus, diabaikan target `thisSection`) | `SaveBusinessLife_Act`: 1 b269 `·` `"All value cannot be empty."`; 3 b490 `·` PRE=true b589 `BIZCODE‖BIZNAME‖RIRATEID‖RIRATE == ""` T=2 F=3; 4 b644 `·` CARI1..12 (`TREATYYEAR ← InputBusinessLife.TREATYYEAR`); 5 b1015 `·` RDB-List `SaveMasterTreatyBusiness_Life_SQL`; 6 b1187 `·` PRE=false `DATASHOW=""` | `SaveMasterTreatyBusiness_Life_SQL` b84 → `INSERTBUSINESS_LIFE` | `POST …/kontrak/{id}/business` · `PUT …/business/{id}` | 🔜 paket 6 |
| Tombol **`Cancel`** b5954 (aksi lama b6069) | `CancelActivityTreatyContract` | — | tutup form | 🔜 paket 9 |
| `ERRMSG` b6960 (`STSSAVE = 100`) · `ERRMSG4` b7711 (`STSSAVE = 99`, nol penulis `ERRMSG4`) | pesan hapus/salin | — | — | 🔜 paket 7 / 9; `ERRMSG4` ➖ |
| Section `ViewRate` b796 `Rate List`, grid b987 RD `BrowseRateLife_RD` (b3063), param `idusedby = ParamID.RIRATEID` b1061 | — | RD kelas `ASM-FW-GISFW-Int-M_RATE_LIFE` = view `RATE_LIFE` (`Endorsment Fac In/RDBList/BrowseLifeRate_SQL.xml` b85 `… FROM RATE_LIFE WHERE IDUSEDBY = …`); saring b552 `.IDUSEDBY = Param.idusedby`; urut `.ID DESC` b745, `.RATE ASC` b782 | `GET …/rate?idusedby=` | 🔜 paket 6 / 9 |
| Kolom `ID` · `USEDBY` · `GENDER` · `CONTRACT` · `AGE` · `RATE` (b1109–b1840) | — | teks apa adanya | idem | 🔜 paket 9 |

## 7. Rule tanpa pemanggil di layar

| Rule | Keadaan |
| --- | --- |
| `Activity/SetParamRate.xml` | ✅ dipanggil FlowAction `ViewRate` b27 (pre-processing) — bukan yatim |
| `ReportDefinition/BrowseTreatyYear_RD.xml` (kelas `TREATYYEAR`), `BrowseRetrocessionLife_RD.xml` (kelas `RETROCESSIONLIFE`) | ➖ hanya metadata saring kolom grid; bukan sumber data |
| `DataTransform/SetOutputParam_DT.xml` (memo `not used`) | ➖ dirujuk hanya dari sel paginator `Tambah` (§4, §5) — memo tidak dipakai sebagai alasan; alasannya aksi paginator itu (OQ-MCRL-11) |
| `Section/GridTreatyBusiness` | ➖ dirujuk b5874 tetapi tidak ada di korpus; target refresh `thisSection` membuatnya tak berpengaruh |
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
