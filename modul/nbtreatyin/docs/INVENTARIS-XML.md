# Inventaris XML — NB Treaty In

> **Berkas bangkitan.** Ditulis `docs/alat/inventaris.py` dari korpus XML Pega (READ-ONLY). Jangan disunting tangan selain lewat `docs/alat/status.json`; jalankan ulang:
>
> ```
> cd modul/nbtreatyin/docs/alat
> python korpus.py korpus.json && python inventaris.py korpus.json
> ```
>
> Jalur korpus: env var `NBTREATYIN_KORPUS` (bawaan: `D:/NUSARE DEV/NusantaraRe/NB Treaty In (Done)`).
> Nama orang, surel, dan alamat **disamarkan berbasis pola** (AC 60, 61) — lihat kepala `inventaris.py`.

## Isi

1. Jangkauan dan bukti lingkup
2. Tabel seluruh berkas (278)
3. Menu dan titik masuk
4. Alur `Flow/InputRealizationTreatyIn`
5. Layar — Harness, FlowAction, Section
6. Activity terjangkau — langkah demi langkah
7. Data Transform terjangkau
8. RDBList — naskah SQL
9. Report Definition
10. Decision Table
11. When — syarat yang DIJALANKAN
12. Rujukan ke rule yang tidak ada di korpus, dan selisih kelas
13. Padanan tombol → activity

## 1 · Jangkauan dan bukti lingkup

### 1.1 Cacah

| Jenis | Berkas | Terjangkau | Tidak terjangkau |
| --- | ---: | ---: | ---: |
| `Activity` | 92 | 63 | 29 |
| `DataTransform` | 12 | 12 | 0 |
| `DecisionTable` | 2 | 2 | 0 |
| `Flow` | 1 | 1 | 0 |
| `FlowAction` | 9 | 9 | 0 |
| `Harness` | 6 | 6 | 0 |
| `RDBList` | 41 | 28 | 13 |
| `ReportDefinition` | 15 | 14 | 1 |
| `Section` | 25 | 25 | 0 |
| `When` | 75 | 16 | 59 |
| **Jumlah** | **278** | **176** | **102** |

**Titik masuk nyata:** `Harness/SFAPortalOpportunities` · `Flow/InputRealizationTreatyIn` — harness portal (akar `Struktur_MenuNBTreatyIn.xlsx`) dan flow yang dibuat tombol *Create* portal itu. Jangkauan diikuti lewat rujukan terstruktur dan indeks `pxRuleReferences` Pega (`graf.py`).

### 1.2 Satu sambungan diputus — dengan bukti

- `Section/ListSuggest` → `Activity/Protection_Act`: Protection_Act di korpus ini md5 8d6fd933b1e10bbbdab2945ffa681ff6 = varian NB FacIn/RNW Fac In (sidik-jari-Protection_Act-sebelum.json), kelas ASM-FW-GISFW-Work. ListSuggest memanggilnya dengan pyActivityClass ASM-FW-GISFW-Data-PolicyTreatyIn - Work bukan leluhur kelas itu, sehingga Pega memilih varian Data-PolicyTreatyIn (157.304 B, 16 langkah) yang TIDAK ada di korpus ini. Keputusan work owner 23-09-2026 butir 6: rantai spreading milik Fac In.

Akibatnya rantai spreading/premi **Fac In** tidak terjangkau. Rule yang hanya hidup lewat sambungan ini tercatat di bab 2 dengan status *varian Fac In* atau *tidak terjangkau*.

### 1.3 Uji silang dengan `jangkauan.json` (tim migrasi)

- Activity terjangkau menurut inventaris ini: **63**; menurut `jangkauan.json`: **64**.
- Ada di `jangkauan.json`, tidak di sini: `Protection_Act` — `Protection_Act` varian Treaty, yang berkasnya tidak ada di korpus ini (bab 1.2).
- Ada di sini, tidak di `jangkauan.json`: nihil.

## 2 · Tabel seluruh berkas

Kolom *Status*: **dibangun di** `<berkas>` · **diganti** (keputusan work owner) · **tidak dibangun** + alasan + bukti. Sumber: `alat/status.json` untuk rule terjangkau; otomatis untuk yang tidak terjangkau.

| # | Jenis | Nama | Kelas | Dipanggil dari | Memanggil | Status |
| ---: | --- | --- | --- | --- | --- | --- |
| 1 | Activity | `AgentSourceBizTreatyIn_Act` | ASM-FW-GISFW-Data-Agent | SourceHierarki | BrowseAgentHierarkiList_RD, BrowseCedingCo_RD | dibangun — sebagian: langkah 3 (`pxShowReport` RD `BrowseAgentHierarkiList_RD`, pyStepsPreCondition=false = selalu jalan) dan langkah 4 (Page-Copy baris RD) = `repository.DaftarAgenHierarki` -> `GET /sumber-bisnis` (dan dijalankan ulang saat Save/Submit untuk mencocokkan pilihan yang dipegang layar - F4, `services.terimaSumberBisnis`); langkah 1-2 tidak dibangun (a): prasyaratnya `pyWorkPage.OfferTreatyIn.QuotationData.btnQuotation=="SOB"/"CedingCo"` — properti yang ditulis NOL rule korpus (btnSOB_DT menulis `pyWorkPage.Quotation.btnQuotation`) ⇒ `Param.Leader` tak pernah diisi di NB, sehingga perluasan simpul pohon mengembalikan daftar akar yang sama |
| 2 | Activity | `BreakDownSpreading_Act` | ASM-FW-GISFW-Work | CountSpreading_Act, TreatyInputPctCommSpreading | GetBreakDownSpread_SQL | tidak dibangun — (b) K9 + KEPUTUSAN-RONDE-12 butir 3/3b (*"tidak usah di migrasi perhitungannya"*: nol tabel, nol rumus, nol layar; diagram Prop J58 `T_POLIS_BREAKDOWN_SPREAD` DIBATALKAN): dipanggil `CountSpreading_Act` langkah 6 dan `TreatyInputPctCommSpreading` langkah 3; mengisi `PolicyTreatyIn.BreakDownSpreadList` (langkah 1.3 `GetBreakDownSpread_SQL`, langkah 2.1 PremiumSpreaded/ClaimSpreaded) — grid `.BreakDownSpreadList` `Section/DetailPolicyTreatyIn` ikut tidak dibangun |
| 3 | Activity | `CalculatePremi_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn | CountNetPremi_act | dibangun — `models.CalculatePremi`; langkah 1-2 membaca `TreatyIn.RNMShareP` / `BrokeragePercentP` (master JSON; aksi refresh sel di wadah `.IsNewPolicyNonProp != 1` = jalur Proporsional) dilewati bila master kosong (`models.MasterTersedia`) — (c): jalur Proporsional di luar ukuran K8 tunggal (P9; disetujui WO 04-10-2026 (F1)): medan master yang DIBACA rule terjangkau jalur NB NonProp, setiap medan di `models/masterxol.go` berkutip langkah XML; langkah 3 (`CountNetPremi_act`) berjalan |
| 4 | Activity | `CekLimitTreatyAcc_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | ListSuggest |  | tidak dibangun — (b) K2 [keputusan WO 03-10-2026: WO lama berlaku, batas 200 juta XML tidak dibangun]: dipanggil aksi refresh `.IsApproved` `Section/ListSuggest`; langkah 2 (bila `PositionNote=="ReasTreatyInSecHead"`) langkah 2.4 `pyWorkPage.LetterNo = @if(Local.NETPREMI>200000000.00,"TREATYINDEPTHEAD","")` (kurs dari `pyWorkPage.TreatyIn.CurrencyList`, JSON master) — dibaca Decision13 `ToTREATYDEPTHEAD` dan syarat tampil tombol Submit atasan; Sec Head selalu naik ke Dept Head (AC 8), pertentangan dicatat tiket 03 |
| 5 | Activity | `cekSpreadingFactIn` ⛔ | ASM-FW-GISFW-Data-Coverage | CheckSpreadingProtect_ACT | SumTSIPremiSpreadedRNM_Act, GetDataByID_SQL, IsEDM | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 6 | Activity | `CheckDataMkt` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn |  | dibangun — `models.TerapkanMO` + aksi `CheckDataMkt` di services (Obj-Save ditiru); langkah ke `OfferFacIn` tidak (halaman Fac) |
| 7 | Activity | `CheckDuplicateOffer` | Data-Portal | SetTreatyIn_Act | GetCountClaim, BrowseTREATY_IN | tidak dibangun — (a) langkah 1-4 berlabel `//`; pemanggil tunggal SetTreatyIn_Act langkah 14 tanpa parameter (`pyPassCurrentParameterPage=false`) ⇒ `TreatyWarning.CAIREINSFACIN` kosong ⇒ GetCountClaim `masterid = NULL` = 0 ⇒ pesan langkah 7-8 tak pernah benar — RALAT AC 59 (tiket 02 putaran 2) |
| 8 | Activity | `CheckRISLIP` ⛔ | ASM-FW-GISFW-Work | Protection_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 9 | Activity | `CheckRISLIP_EDM` ⛔ | ASM-FW-GISFW-Work | Protection_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 10 | Activity | `CheckSpreadingProtect_ACT` ⛔ | ASM-FW-GISFW-Work | Protection_Act | CheckSpreadingProtectAnekaGolf_ACT, CheckSpreadingProtectFire_ACT, CheckSpreadingProtectMCargoMBU_ACT, CheckSpreadingProtectPATravel_ACT, GetDateValidity_ACT, G … | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 11 | Activity | `CheckSpreadingProtectAnekaGolf_ACT` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT | IsAneka, IsBondingAndCustomBonds, IsGolfInsurance | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 12 | Activity | `CheckSpreadingProtectFire_ACT` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT | IsFire | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 13 | Activity | `CheckSpreadingProtectMCargoMBU_ACT` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT | IsMBU, IsMarineCargo | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 14 | Activity | `CheckSpreadingProtectPATravel_ACT` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT | IsPA, IsTravel | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 15 | Activity | `ConcatSlipOfferNo_Act` | ASM-FW-GISFW-Data-Quotation | HistoricalSurveyReportDtl, HistoricalSurveyReportDtlUW |  | tidak dibangun — hanya dari aksi refresh `change` medan `.DateofSurvey` di `Section/HistoricalSurveyReportDtl` dan `HistoricalSurveyReportDtlUW` (popup survei); langkah 1 `Quotation.NoOfferSlip = ""`, langkah 2.1 kalang `.QuotationList` merangkai `.NoOfferSlip + "; "`. (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 16 | Activity | `ConvertHistoryDate` | Data-Portal | SetTreatyIn_Act |  | tidak dibangun — (a) efeknya dibaca NOL rule NB. RALAT bunyi putaran 1 (*"hanya dipanggil bagian survei historis"*): pemanggil tunggal `Activity/SetTreatyIn_Act` langkah 12, yang hanya dipanggil `TreatyRealizationCheckXOLList` langkah 4 (jalur NonProp: `InputPolicyTreatyInPre_Act` langkah 10 `IsNewPolicyNonProp==1`) — terjangkau. Isinya hanya kalang `TreatyIn.CommentList` (`.Date` +7 jam lalu `.ConvertDate=1`, bila `.ConvertDate!=1`) — riwayat komentar MASTER. Bukti nol pembaca: `CommentList` hanya muncul di dua berkas korpus, `Activity/SetTreatyIn_Act.xml` (langkah 9 menulisnya, hanya `param.revisionstate==1` — tidak pernah dari NB) dan `Activity/ConvertHistoryDate.xml` sendiri; nol Section/Activity/When NB membacanya. RALAT P9 atas alasan lama *"di luar daftar medan pengecualian K8 (...) ⇒ (c)"* — ukuran K8 tunggal (P9; disetujui WO 04-10-2026 (F1)): medan master yang DIBACA rule terjangkau jalur NB NonProp, setiap medan di `models/masterxol.go` berkutip langkah XML; CommentList tidak dibaca rule NB mana pun, maka (a), bukan (c) |
| 17 | Activity | `CountNetPremi_act` | ASM-FW-GISFW-Data-PolicyTreatyIn | CalculatePremi_Act, CountOGPONP_Act, CountOverridingCommOgp_Act, CountOverridingCommOnp_Act, CountResult1Onp_Act, CountResult1_Act, CountResult2Ogp_act, CountRe … | SetDueTo_act, SetPPNPPH | dibangun — `models.CountNetPremi` |
| 18 | Activity | `CountOGPONP_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn | CountNetPremi_act, CountResult1Onp_Act, CountResult1_Act, CountResult2Ogp_act, CountResult2Onp_act, CountSpreading_Act, SetValidateInstallment_Act | dibangun — `models.CountOGPONP` — RALAT R6 audit silang P3 (04-10-2026): langkah 9 (CountSpreading_Act per baris) dan 10 (SetValidateInstallment_Act) diputar ulang server atas daftar tersimpan bila sel uang ber-action set ini berubah (`models.TerapkanPemicu`) - W5 |
| 19 | Activity | `CountOverridingCommOgp_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW | CountNetPremi_act | tidak dibangun — (a) tidak terpicu: sel pemicu ber-pyReadOnly / grid readOnly (audit silang P3 7.4, kategori ditetapkan konsolidasi P3K 04-10-2026). Satu-satunya pemanggil: aksi refresh sel `.ResultOgp2` `Section/DetailDeptHeadTreatyIn_UW`, `pyReadOnly=true` ⇒ tidak pernah terpicu. Fungsi port `models.CountOverridingCommOgp` (langkah 3, 5; langkah 1/2/4 `//`) dipertahankan berujian (`TestAktivitasLayarDeptHeadMembulatkanEmpatDesimal`); `POST /hitung` menolaknya 409 (`services/aksiposisi.go`) |
| 20 | Activity | `CountOverridingCommOnp_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW | CountNetPremi_act | tidak dibangun — (a) tidak terpicu: sel pemicu ber-pyReadOnly / grid readOnly (audit silang P3 7.4, kategori ditetapkan konsolidasi P3K 04-10-2026). Satu-satunya pemanggil: aksi refresh sel `.ResultOnp2` `Section/DetailDeptHeadTreatyIn_UW`, `pyReadOnly=true` ⇒ tidak pernah terpicu. Fungsi port `models.CountOverridingCommOnp` (langkah 3-4; langkah 1/2 `//`) dipertahankan berujian (`TestAktivitasLayarDeptHeadMembulatkanEmpatDesimal`); `POST /hitung` menolaknya 409 (`services/aksiposisi.go`; `TestHitungHanyaAksiSelTerbukaDiLayarPosisi`) |
| 21 | Activity | `CountPctInstallment_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |  | tidak dibangun — (a) tidak terpicu: sel pemicu ber-pyReadOnly / grid readOnly (audit silang P3 7.4, kategori ditetapkan konsolidasi P3K 04-10-2026). Satu-satunya pemanggil: aksi refresh sel `.Premium` grid S45 `.ListInstallment` di `Section/DetailPolicyTreatyIn` dan `Section/DetailDeptHeadTreatyIn_UW` — sel `pyReadOnly=true`, grid `pyEditingMode`/`pyRowEditing` `readOnly` ⇒ tidak pernah terpicu. Fungsi port `models.CountPctInstallment` dipertahankan berujian (`TestCountPctInstallmentBarisIdx`); `POST /hitung` menolaknya 409 (`services/aksiposisi.go`; `TestHitungHanyaAksiSelTerbukaDiLayarPosisi`) |
| 22 | Activity | `CountResult1_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | CountOGPONP_Act, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn | CountNetPremi_act | dibangun — `models.CountResult1` (langkah 3-4 ber-master dilewati bila master kosong) |
| 23 | Activity | `CountResult1Onp_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | CountOGPONP_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn | CountNetPremi_act | dibangun — `models.CountResult1Onp` |
| 24 | Activity | `CountResult2Ogp_act` | ASM-FW-GISFW-Data-PolicyTreatyIn | CountOGPONP_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn | CountNetPremi_act | dibangun — `models.CountResult2Ogp` |
| 25 | Activity | `CountResult2Onp_act` | ASM-FW-GISFW-Data-PolicyTreatyIn | CountOGPONP_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn | CountNetPremi_act | dibangun — `models.CountResult2Onp` |
| 26 | Activity | `CountRiCommOgp_act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW | CountNetPremi_act | tidak dibangun — (a) tidak terpicu: sel pemicu ber-pyReadOnly / grid readOnly (audit silang P3 7.4, kategori ditetapkan konsolidasi P3K 04-10-2026). Satu-satunya pemanggil: aksi refresh (change) sel `.ResultOgp1` `Section/DetailDeptHeadTreatyIn_UW`, sel itu `pyReadOnly=true` / `pyEditOptions=Read-only` ⇒ event tidak pernah terpicu. Fungsi port `models.CountRiCommOgp` (langkah 3; langkah 1, 2, 4 `//`) dipertahankan berujian (`TestAktivitasLayarDeptHeadMembulatkanEmpatDesimal`), tetapi `POST /hitung` menolaknya 409 di setiap posisi (`services/aksiposisi.go`; `TestHitungHanyaAksiSelTerbukaDiLayarPosisi`) |
| 27 | Activity | `CountRiCommOnp_act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW | CountNetPremi_act | tidak dibangun — (a) tidak terpicu: sel pemicu ber-pyReadOnly / grid readOnly (audit silang P3 7.4, kategori ditetapkan konsolidasi P3K 04-10-2026). Satu-satunya pemanggil: aksi refresh sel `.ResultOnp1` `Section/DetailDeptHeadTreatyIn_UW`, `pyReadOnly=true` ⇒ tidak pernah terpicu. Fungsi port `models.CountRiCommOnp` (langkah 2-3) dipertahankan berujian (`TestAktivitasLayarDeptHeadMembulatkanEmpatDesimal`); `POST /hitung` menolaknya 409 (`services/aksiposisi.go`) |
| 28 | Activity | `CountSpreading_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | CountOGPONP_Act, TreatyInputPctCommSpreading, TreatyNonPropOutSetSpreading, TreatyNonPropSetSpreading, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, Spreadin … | BreakDownSpreading_Act | dibangun — `models.CountSpreading` (langkah 6 BreakDownSpreading tidak — KEPUTUSAN 3b) — RALAT R6 audit silang P3 (04-10-2026): langkah 1-3 berlabel `//` tidak lagi diport (port lama menjalankannya untuk baris Param.Index; nol beda hasil, 4.1 menulis keempat medan setiap baris); kolom PremiumSpreaded/ClaimSpreaded hanya dari langkah 4.1 di server (`models.TerapkanPemicu`, terpicu sel %Share / baris dihapus / CountOGPONP_Act 9) - W5; aksi sel terbuka di layar admin DAN atasan (subsection NonProp) bila grid terbuka (`models.SpreadingDariLayar`) - W2; R7 (PERMINTAAN H2): CountSpreading_Act tidak menulis TreatyName/Currency/CurrencyID/SplitRNMSharePct, jadi anggota baris itu tidak diterima dari layar - nilai baris server (`models.gabungSpreading`, `TestAnggotaBarisSpreadingDariServer`) |
| 29 | Activity | `FetchMasterTreatyIn` | ASM-FW-GISFW-Work | TreatyInputPctCommSpreading | BrowseTreatyInJoinEDM | tidak dibangun sebagai pembongkar JSON — satu-satunya pemanggil `TreatyInputPctCommSpreading` langkah 1 (preACT langkah 17, HANYA bila BUKAN NonProportional = jalur Proporsional): RiCommOgp dibaca dari kolom view TREATYINDETAILJOINEDM (`repository.KomisiKontrak`: TREATYID = NoOffer, TREATYTYPE/TREATYGROUP/RIOGR/RIONR); medan master lain (SpreadingTotalPct, SpreadingTypeID, SpreadingType) hanya ada di JSON — (c): jalur Proporsional di luar ukuran K8 tunggal (P9; disetujui WO 04-10-2026 (F1)): medan master yang DIBACA rule terjangkau jalur NB NonProp, setiap medan di `models/masterxol.go` berkutip langkah XML (K8 terbatas jalur NonProp/XOL) |
| 30 | Activity | `FetchTreatyGroupOJK` | ASM-FW-GISFW-Work-NB | GeneratePolicyNoTreaty_Act, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT | BrowseTreatyGroup_RD | dibangun — `repository.OJKGrupTreaty` + `models.SetelOJK` |
| 31 | Activity | `FetchTreatyGroupOldID` | ASM-FW-GISFW-Work | GeneratePolicyNoTreaty_Act, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT | FetchTreatyGroupOLDID | dibangun — `repository.OldIDGrupTreaty` + `models.SetelGrupLama` (dari GeneratePolicyNoTreaty_Act; langkah 7 preACT berlabel `//`) — RALAT R6 audit silang P3 (04-10-2026): langkah 2-3 (Param.Errmsg "Cannot fetch Treaty Group ID, Contact IT", Page-Set-Messages bila TreatyGroupID=="") DIBANGUN - `models.GrupTreatyTakTerbaca`, `services.validasiKirim` menahan submit Dept Head (OK nomor polis) 422; nomor tetap dibentuk (GeneratePolicyNoTreaty_Act 30 Obj-Save WithErrors) |
| 32 | Activity | `FillPaymentInstallment` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |  | dibangun — `models.FillPaymentInstallment` (langkah 4-5 `//` tidak) — RALAT R6 audit silang P3 (04-10-2026): juga diputar ulang SERVER saat Save/Submit bila `.Installment` berubah (`models.TerapkanPemicu`); ListInstallment tidak pernah diterima dari layar (grid readOnly) - W5 |
| 33 | Activity | `GeneratePolicyNoTreaty_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW | FetchTreatyGroupOJK, FetchTreatyGroupOldID, GETTanggalClosing_SQL, GetKodeProdNonLife_SQL, GetSQLDate, GetSequenceNumber_SQL | dibangun — `repository.TerbitkanNomorPolis` + `services.TerbitkanNomor` lewat `inti/backend/penomor` (langkah 14-24 `//` tidak); tombol hanya di posisi Dept Head (AC 8) — RALAT R6 audit silang P3 (04-10-2026): langkah 11 -> FetchTreatyGroupOldID langkah 3 (pesan menahan submit) dibangun - 7.4 |
| 34 | Activity | `GetCurrencyMaster` ⛔ | ASM-FW-GISFW-Data-PropertyItem | SumTSIPremiSpreadedRNM_Act | GetAllCurrency, BrowseCurrency_RD | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 35 | Activity | `GetDateValidity_ACT` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 36 | Activity | `GetTreatyName` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT | GetTreatyName_SQL, GetTreatyName_SQL2, IsEDM, IsErrorSpreading, IsFire, IsMarineCargo | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 37 | Activity | `InputParamUploadReas_act` | ASM-FW-GISFW-Int-OFFERJSON | SetCategoryAttach |  | tidak dibangun — (a) hanya dari `SetCategoryAttach` langkah 4.2; Obj-Browse `ASM-FW-GISFW-Int-DOCUMENT_POLIS` (IDPEGA, KATEGORI_2) → `AttachShowList`, yang hanya dipakai `.CountAttach` (langkah 4.3 pemanggil) — pembacanya di korpus hanya `CheckRISLIP`, `CheckRISLIP_EDM`, `CheckSpreadingProtect_ACT`, ketiganya hanya dipanggil `Protection_Act` — sambungan DIPUTUS (varian Fac In: `graf.py` PUTUS, keputusan work owner 23-09-2026 butir 6); nol Section menampilkannya; nol penulisan |
| 38 | Activity | `InputPolicyTreatyEDMDetail_NP` | ASM-FW-GISFW-Work | InputPolicyTreatyInDetail_NonProp | InsertToTreatyOutXOLList, InsertToTreatyXOLList, GetDataCurrencyByName_SQL | tidak dibangun — (a) tidak terjangkau: satu-satunya pemanggil InputPolicyTreatyInDetail_NonProp langkah 11 bersyarat `TreatyMasterInEDM` DAN `IsEDMInputOnNB != true`, padahal langkah 10 mengisi `IsEDMInputOnNB = true` tepat ketika `TreatyMasterInEDM` |
| 39 | Activity | `InputPolicyTreatyInDetail_NonProp` | ASM-FW-GISFW-Work | InputPolicyTreatyInDetail_preACT | InputPolicyTreatyEDMDetail_NP, InsertToTreatyXOLList, InsertToTreatyXOLListRetroShare, TreatyNonPropSetSpreading, BrowseTreatyInJoinEDM, GetDataCurrencyByName_S … | dibangun — [keputusan work owner] K8 (03-10-2026): `models.InputDetailNonProp` (1, 7-10, 13, 15-23) + `services.pilihBisnisNonProp` (2-6: master baca-saja `repository.MasterXOLDariJSON`); 11 tidak terjangkau (syarat mustahil sesudah 10), 12 `pyExpanded`, 14 halaman OfferFacIn (kasus Fac, nol kolom diagram), 24 `//` |
| 40 | Activity | `InputPolicyTreatyInDetail_preACT` | ASM-FW-GISFW-Work | SetValue_Act | CountResult1_Act, FetchTreatyGroupOJK, FetchTreatyGroupOldID, InputPolicyTreatyInDetail_NonProp, SetTreatyCurrencyID, TreatyInputPctCommSpreading, BusinessType_ … | dibangun — sebagian: langkah 1-6 dan 8 (RD `BrowseTreatyJoinEDM` = `repository.DetailKontrak`, view; mata uang, grup, OJK), 11, 14, 15, 19 di `models/pilihbisnis.go` + `services.PilihBisnis`; 16 (NonProp) dan 18 (PPN/PPH layer XOL) [keputusan work owner] K8 — `services.pilihBisnisNonProp`; 17 (`TreatyInputPctCommSpreading`, bukan NonProportional) — dari kolom view, lihat entrinya; 7 dan 12 berlabel `//`. Tidak: langkah 9-10 (RDB `BrowseTreatyInDetailJoinEDM` `select JSONDATA as CLASSOFBUSINESS from pooldata.M_TREATY_IN_DETAIL_EDM where ID={TreatyIn.ID}` + Java `adoptJSONObject` ke `TreatyIn`) dan 13 (kalang `TreatyIn.INSTALLMENT().InstallmentList` hasil 9-10, `Currency == PolicyTreatyIn.Currency` → `ListInstallment`) — (c): JSON tabel `M_TREATY_IN_DETAIL_EDM`, di luar dua tabel pengecualian K8 butir 2 (`M_TREATY_IN` / `M_TREATY_IN_EDM` lewat `BrowseTreatyIn`) |
| 41 | Activity | `InputPolicyTreatyInPost_Act` | ASM-FW-GISFW-Work | InboxPolicyTreatyIn | InsertHistoryAkseptasiPega, SaveViewSuggest, TreatyRealizationCheckDuplicate | dibangun — `services.Kirim` (langkah 2 riwayat, langkah 3 cek duplikat bila IsApproved 1, langkah 4 SaveViewSuggest → HISTORYAKSEPTASIPRODUCTION); Page-Clear-Messages langkah 1 = halaman tidak membawa pesan antar permintaan |
| 42 | Activity | `InputPolicyTreatyInPre_Act` | ASM-FW-GISFW-Work | DeptHeadTreatyIn_UW, InboxPolicyTreatyIn | SetCategoryAttach, TreatyRealizationCheckXOLList, GetOldIDBusiness_SQL, GetSQLDate, SelectSpreadingTreatyInProduction | dibangun — `services.siapkan`: langkah 2, 3-4, 9 (hari tutup buku dari TANGGAL_CLOSING — preseden WO, penjaga ambang), 10 (TreatyRealizationCheckXOLList — K8, `services.siapkanNonProp`); langkah 1 tidak (lampiran); prasyarat langkah 10 `.PolicyTreatyIn.EDMType=="3"` membaca medan dokumen lama — RALAT F3 (04-10-2026): kini berkolom `T_GENERAL_POLIS.EDM_TYPE` (katalog `kKode(pt+"EDMType")`, migrasi 320), jadi nilai dokumen lama yang dimuat tidak hilang — RALAT R6 audit silang P3 (04-10-2026): syarat ketiga `models.PerluCekDaftarXOL` (QuotationData.ProportionalType != "Proportional", di luar XML langkah 10) BERDASAR spec-penyimpanan AC 33 (ID-35): penanda IsNewPolicyNonProp tak pernah diturunkan (preDT 4-5), baris XOL pada polis Proportional ditolak penyimpanan; perilaku tetap, dasar dicatat tiket 19 |
| 43 | Activity | `InputPolicyTreatyOutDetail_NonProp` | ASM-FW-GISFW-Work | InputPolicyTreatyOutDetail_preACT | InsertToTreatyOutXOLList, TreatyNonPropOutSetSpreading, BrowseTreatyOut, GetDataCurrencyByName_SQL, TreatyMasterInEDM | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — pemanggil tunggal `Activity/InputPolicyTreatyOutDetail_preACT` langkah 17 (Call); langkah 5 RDB-List `RDBList/BrowseTreatyOut` membaca JSON master treaty keluar: `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`; tidak dibaca, tidak ditulis (AC 57 ⛔, AC 58 ✅) |
| 44 | Activity | `InputPolicyTreatyOutDetail_preACT` | ASM-FW-GISFW-Work | SetValueRetro_Act | CountResult1_Act, FetchTreatyGroupOJK, FetchTreatyGroupOldID, InputPolicyTreatyOutDetail_NonProp, SetTreatyCurrencyID, BusinessType_DeT, BrowseTreatyOutDetail,  … | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — pemanggil tunggal `Activity/SetValueRetro_Act` langkah 1 (Call); langkah 2 `Param.pyReportName` RD `BrowseTreatyOutDetail`, langkah 10 RDB-List `RDBList/BrowseTreatyOutDetail` (`select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`), langkah 17 Call `InputPolicyTreatyOutDetail_NonProp`; tidak dibaca, tidak ditulis (AC 57 ⛔, AC 58 ✅) |
| 45 | Activity | `InputQuotation_PreAct` | ASM-FW-GISFW-Data-Quotation | DetailPolicyTreatyIn | btnCedingCO_DT, btnSOB_DT | dibangun — langkah 1 (`Param.Acton=="SOB"` -> `btnSOB_DT`) = `models.TombolSOB` di `models.HasilPostDT` (klik baris, `services.PilihSumberBisnis`); langkah 2 (`Acton=="Ceding"` -> `btnCedingCO_DT`) tidak dibangun (a): satu-satunya pengirim Acton di korpus (tombol `Select Source Of Business` `DetailPolicyTreatyIn`/`GeneralPolicyTreatyIn`) mengirim `Acton=SOB` |
| 46 | Activity | `InsertHistoryAkseptasiPega` | ASM-FW-GISFW-Work | InputPolicyTreatyInPost_Act, DeptHeadTreatyIn_UW | InsertHistoryAkseptasiPega_Sql | dibangun — `repository.CatatRiwayat` di transaksi submit (`services.Kirim`); OPERATORID = login, USERNAME = nama tampilan |
| 47 | Activity | `InsertToTreatyOutXOLList` | ASM-FW-GISFW-Work | InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyOutDetail_NonProp | GetDataCurrencyByName_SQL | tidak dibangun — (a) pemanggilnya InputPolicyTreatyEDMDetail_NP (tidak terjangkau) dan InputPolicyTreatyOutDetail_NonProp (treaty KELUAR, K8 butir 4 ⛔) |
| 48 | Activity | `InsertToTreatyXOLList` | ASM-FW-GISFW-Work | InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, TreatyRealizationCheckXOLList | GetDataCurrencyByName_SQL | dibangun — [keputusan work owner] K8 (03-10-2026): `models.InsertToTreatyXOLList` (seluruh langkah; DueTo kosong karena `Page-New InputXOL` ditiru apa adanya) -> T_POLIS_XOL, T_POLIS_XOL_LAYER |
| 49 | Activity | `InsertToTreatyXOLListRetroShare` | ASM-FW-GISFW-Work | InputPolicyTreatyInDetail_NonProp | GetDataCurrencyByName_SQL | dibangun — [keputusan work owner] K8 (03-10-2026): `models.InsertToTreatyXOLListRetroShare` (2.2 nol pembaca, 2.3.3.3.2-3 `//`) -> T_POLIS_XOL, T_POLIS_XOL_LAYER |
| 50 | Activity | `ProtectCedingCo` ⛔ | ASM-FW-GISFW-Work | Protection_Act | GetDataAgentByNameNonLife_SQL | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 51 | Activity | `ProtectCoverage_Act` ⛔ | ASM-FW-GISFW-Work | Protection_Act | ProtectFIREMBUPA_Act, CheckZipCode_SQL, GetDataDoubleCase, IsAneka, IsEDM, IsGolfInsurance, IsMarineCargo, IsRenewal, IsTBonding | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 52 | Activity | `ProtectCurrencyTSI_Act` ⛔ | ASM-FW-GISFW-Work | Protection_Act | GetCurrency, IsFire | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 53 | Activity | `ProtectDate` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn |  | dibangun — `models.ProtectDate` (refresh EndDate + validasi kirim) |
| 54 | Activity | `ProtectFIREMBUPA_Act` ⛔ | ASM-FW-GISFW-Work | ProtectCoverage_Act | CheckZipCode_SQL, GetCategoryOccupation, GetDataDoubleCase, GetObjectItembyName_SQL, IsEDM, IsFire, IsMBU, IsPA, IsRenewal | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectCoverage_Act` — seluruhnya tidak terjangkau |
| 55 | Activity | `Protection_Act` ⛔ | ASM-FW-GISFW-Work | ListSuggest | CheckRISLIP, CheckRISLIP_EDM, CheckSpreadingProtect_ACT, ProtectCedingCo, ProtectCoverage_Act, ProtectCurrencyTSI_Act, ProtectPremiPolicy_Act, ProtectRenewal_Ac … | tidak dibangun — **varian Fac In**: satu-satunya pemanggil hidup diputus dengan bukti (bab 1.2) |
| 56 | Activity | `ProtectPremiPolicy_Act` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT, Protection_Act | IsEDM, IsLife | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT`, `Protection_Act` — seluruhnya tidak terjangkau |
| 57 | Activity | `ProtectRenewal_Act` ⛔ | ASM-FW-GISFW-Work | Protection_Act | IsRenewal | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 58 | Activity | `ProtectShareCedant_Act` ⛔ | ASM-FW-GISFW-Work | Protection_Act | IsEdmAdjShareCedant, IsLife | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 59 | Activity | `ProtectShipData_Act` ⛔ | ASM-FW-GISFW-Work | Protection_Act | IsEDM | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 60 | Activity | `RemoveTypeTax_ACT` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn |  | dibangun — `models.RemoveTypeTax` (langkah 1 TempWorkPage tidak ada) — RALAT R6 audit silang P3 (04-10-2026): juga dijalankan server saat Save/Submit bila `.FlagPPH` kiriman berubah (`models.GabungMasukanLayar`; TypeTax tersembunyi tak diterima, W4) |
| 61 | Activity | `SaveJsonPolisTreatyIn_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | InputRealizationTreatyIn | GenerateNoPolicy, SavePolisTreatyIn_SQL | tidak dibangun — (b) AC 16 [keputusan WO: sistem baru tidak menulis JSON] + diagram Prop F9–F11 (`T_GENERAL_POLIS` menyerap 7 kolom datar `json_polis`): Flow `Utility1` ('Save json policy', sesudah Decision8 Else; parameter `isFOR` dikirim kosong). Langkah 1 (`CopyToPolicy`) dan 3 berlabel `//`; langkah 4 `GenerateNoPolicy` hasilnya dibuang; langkah 6–7 `@ASM.GetPageJSONString()` → `SavePolisTreatyIn_SQL`. Padanan: penyimpanan relasional `repository.SimpanHalaman` (tabel diagram); parameter `isFOR` (AC 66) nol efek: pemanggil tunggal Flow Utility1 mengirim kosong, langkah 3 (EDM) berlabel `//`, langkah 5 (POLICY) kotak When tak dicentang — RALAT AC 66 (tiket 06 putaran 2); dokumen lama yang ditulisnya (langkah 6: halaman PolicyTreatyIn, kelas ASM-FW-GISFW-Data-PolicyTreatyIn) DIBACA pemuat tiket 22 `backend/alat/pemuatlama` (`models.PecahDokumenLama`); F3 (04-10-2026): setiap medan halaman itu berkolom, disalin (SuggestList → HISTORYAKSEPTASIPRODUCTION), atau dibuang berbukti (`models/medan_abaikan_lama.json` bagian `pola`); F6: baris hasil pemuat dikenali dari IDPEGA `<kelas> <pyID>` + status Resolved-Completed, tanpa kolom SUMBER |
| 62 | Activity | `SaveViewSuggest` | ASM-FW-GISFW-Work | InputPolicyTreatyInPost_Act | InsertViewSuggest_SQL | dibangun — langkah 2 "UNTUK TREATY" → `models.UsulanBelumTersimpan` + `repository.CatatUsulan` (POOLDATA.HISTORYAKSEPTASIPRODUCTION, K4; baris IsSave kosong, lalu IsSave=Yes). [penyimpangan sadar] K4: syarat BusinessFac=="F" tidak ditiru; DIV NULL (tanpa sumber). [penyimpangan sadar — disetujui WO 04-10-2026] (F2): ditulis di submit KETIGA jenjang (XML hanya `InputPolicyTreatyInPost_Act` admin; akibat langsung K4 — SuggestList tak punya tempat simpan lain di antara langkah), TGL_INP jam 24 (XML `hh`), NOURUT dari repository (XML `.pxListSubscript`) — terbukti SAMA dengan `.pxListSubscript` baris SuggestList yang dibangun ulang dari tabel (uji `models.TestCatatanBaruBerposisiNourutBerikut`, `handlers.TestNourutUsulanSamaDenganSubskripSuggestList`). Langkah 1 "UNTUK FACIN" (OfferFacIn.ViewSuggest) milik kasus fakultatif — tidak terjangkau kasus treaty |
| 63 | Activity | `SearchHierarkiSourceBizAgentTreatyIn_Act` | ASM-FW-GISFW-Data-Agent | SourceHierarki | BrowseClientEmail_SQL | tidak dibangun — (a) aktivitas tombol `Choose` `SourceHierarki`: langkah 1-4 hanya menulis `pyWorkPage.OfferTreatyIn.QuotationData.SourceOfBusiness/SobName/SobLeader0/SobLeader1/Email` (+ RDB `BrowseClientEmail_SQL`), dibaca NOL rule NB (pembaca `OfferTreatyIn.QuotationData.*` satu-satunya = prasyarat `btnQuotation` `AgentSourceBizTreatyIn_Act`, yang tidak ditulisnya); `window.close` + `opener.location.reload` = popup ditutup di `components/PilihSumberBisnis.tsx`; nol Obj-Save dan nol refresh berhitung sesudahnya (F4) |
| 64 | Activity | `serviceInsertArasapas_act` | ASM-FW-GISFW-Data-PolicyTreatyIn | InputRealizationTreatyIn |  | dibangun — efek keluar sesudah realisasi selesai — `services/konversi.go` + `models.RakitMuatanKonversi` (KEPUTUSAN-RONDE-12 butir 7: bentuk 4 medan dari kelas Work; sambungan M_LINK_SERVICE `[terbuka]`, di luar produksi dilewati, gagal tidak membatalkan simpan) |
| 65 | Activity | `SetCategoryAttach` | ASM-FW-GISFW-Data-OfferFacIn | InputPolicyTreatyInPre_Act | InputParamUploadReas_act, setCategoryAttachment_DT, AttachmentLife, CategoryAttach_SQL | tidak dibangun — (a) efeknya dibaca nol rule terjangkau: dipanggil `InputPolicyTreatyInPre_Act` langkah 1 (Call `ASM-FW-GISFW-Data-OfferFacIn.SetCategoryAttach`). Untuk treaty (`BusinessFac` "T") hanya langkah 1 (`CategoryAttach_SQL` → halaman `Attachment`) dan langkah 4 (per kategori: `InputParamUploadReas_act`, `.CountAttach`) yang berjalan, nol penulisan; pembacanya di korpus hanya `CheckRISLIP`, `CheckRISLIP_EDM`, `CheckSpreadingProtect_ACT`, ketiganya hanya dipanggil `Protection_Act` — sambungan DIPUTUS (varian Fac In: `graf.py` PUTUS, keputusan work owner 23-09-2026 butir 6); nol Section menampilkannya |
| 66 | Activity | `SetCurrency_act` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn | GetCurrency | dibangun — `repository.NamaMataUang` + `models.SetelNamaMataUang` (aksi `SetCurrency`) |
| 67 | Activity | `SetDueTo_act` | ASM-FW-GISFW-Data-PolicyTreatyIn | CountNetPremi_act, PolicyTreatyInDeclineConfirm, DetailPolicyTreatyIn, ListSuggest |  | dibangun — `models.SetDueTo` |
| 68 | Activity | `SetFlagOccupation_ACT` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT | GetCategoryOccupation, IsFire | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 69 | Activity | `SetPPNPPH` | ASM-FW-GISFW-Data-PolicyTreatyIn | CountNetPremi_act | BrowseClientName_RD | dibangun — `models.SetPPNPPH` + `repository.StsPKPAgen` (syarat FlagPPH/STS_PKP langkah 4); diuji nilai XML hitung tangan (TestSetPPNPPHNilaiXML, TestSetPPNPPHMembacaStatusPKPAgen); langkah 1 membaca `QuotationData.SourceOfBusiness` = pilihan SOB yang dipegang layar saat refresh/Save/Submit (F4, TestPilihanDipegangMenggerakkanPPNPPH) |
| 70 | Activity | `SetReinstatementPct` | ASM-FW-GISFW-Data-TreatyInLimits | TreatySetReinstatement |  | dibangun — [keputusan work owner] K8 (03-10-2026): `models.SetReinstatementPct` (1-3.3) |
| 71 | Activity | `SetSurveyReport_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | HistoricalSurveyReportDtl |  | tidak dibangun — `Activity/SetSurveyReport_Act.xml` langkah 1 hanya `Property-Set pyWorkPage.Quotation = pyWorkPage.PolicyTreatyIn.QuotationData` (nol RDB-List, nol Obj-Save); dipanggil tombol Submit `Section/HistoricalSurveyReportDtl` (runActivity lalu `window.close`). (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 72 | Activity | `SetTreatyCurrencyID` | ASM-FW-GISFW-Work | InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT | GetCurrencyIDByName | dibangun — `repository.IDMataUangDariNama` + `models.SetelIDMataUang` |
| 73 | Activity | `SetTreatyIn_Act` | Data-Portal | TreatyRealizationCheckXOLList | CheckDuplicateOffer, ConvertHistoryDate, TreatyInInputVis, TreatySetReinstatement, BrowseTreatyIn, GetCurrentDate, SaveTreatyIn | dibangun — sebagian, [keputusan work owner] K8 (03-10-2026): 3-5 (`TreatyIn.ID` + RDB `BrowseTreatyIn` + adoptJSONObject -> `repository.MasterXOLDariJSON`, baca-saja) dan 13 (`TreatySetReinstatement`); 1-2, 6-7 penanda tampilan layar master Data-Portal tanpa pembaca NB; 8-11 revisi master yang MENULIS JSON (`SaveTreatyIn`) hanya bila `revisionstate==1` — tidak pernah dari NB; 12 `ConvertHistoryDate` (komentar master `TreatyIn.CommentList`, dibaca nol rule NB — (a), lihat entrinya); 14 `CheckDuplicateOffer` (1-4 `//`); 15 `FlagExcel` tanpa pembaca |
| 74 | Activity | `SetValidateInstallment_Act` | ASM-FW-GISFW-Data-PolicyTreatyIn | CountOGPONP_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |  | dibangun — `models.SetValidateInstallment` (langkah 1 Page-Clear-Messages SESUDAH validasi tidak — AC 78) — RALAT R6 audit silang P3 (04-10-2026): sel grid angsuran pemicunya readOnly (tidak terpicu); berjalan sebagai CountOGPONP_Act langkah 10 - juga diputar ulang server saat sel uang berubah (`models.TerapkanPemicu`) - W5 |
| 75 | Activity | `SetValue_Act` | ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM | BusinessAndSOBList, BusinessAndSOBListRetro | InputPolicyTreatyInDetail_preACT | dibangun — `services.PilihBisnis` (Obj-Save ditiru); R7 pola F4: `Param.ID` hanya dari baris grid aktif `BrowseTreatyJoinEDM` tersaring kasus (`models.PeriksaPilihanBisnis`, di luar -> 422, `TestPilihBisnisDiLuarDaftarPopupDitolak`) |
| 76 | Activity | `SetValueRetro_Act` | ASM-FW-GISFW-Int-TREATYOUTDETAIL | BusinessAndSOBListRetro | InputPolicyTreatyOutDetail_preACT | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — pemanggil tunggal aksi `runActivity` medan `.pyTemplateInputBox` di grid `Section/BusinessAndSOBListRetro` (popup treaty keluar); langkah 1 Call `InputPolicyTreatyOutDetail_preACT` → JSON `M_TREATY_OUT` (`select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`); tidak dibaca, tidak ditulis (AC 57 ⛔, AC 58 ✅) |
| 77 | Activity | `SumTSIPremiSpreadedRNM_Act` ⛔ | ASM-FW-GISFW-Data-Coverage | cekSpreadingFactIn | GetCurrencyMaster, SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_FIRE_Act, SumTSIPremiSpreadedRNM_GOLF_Act, SumTSIPremiSpreadedRNM_MARINECARGO_Act, S … | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `cekSpreadingFactIn` — seluruhnya tidak terjangkau |
| 78 | Activity | `SumTSIPremiSpreadedRNM_ANEKA_Act` ⛔ | ASM-FW-GISFW-Data-Coverage | SumTSIPremiSpreadedRNM_Act | GetDataCurrencyByName_SQL, IsAneka, IsBondingAndCustomBonds, IsBuilderRisk, IsCMI, IsEdmAdjSpreading, IsHE, IsLiability, IsMBD, IsMarineHull, IsMarineHullOffsho … | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 79 | Activity | `SumTSIPremiSpreadedRNM_FIRE_Act` ⛔ | ASM-FW-GISFW-Data-Coverage | SumTSIPremiSpreadedRNM_Act | GetDataCurrencyByName_SQL, IsEdmAdjSpreading, IsFire | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 80 | Activity | `SumTSIPremiSpreadedRNM_GOLF_Act` ⛔ | ASM-FW-GISFW-Data-Coverage | SumTSIPremiSpreadedRNM_Act | GetDataCurrencyByName_SQL, IsEdmAdjSpreading, IsGolfInsurance | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 81 | Activity | `SumTSIPremiSpreadedRNM_MARINECARGO_Act` ⛔ | ASM-FW-GISFW-Data-Coverage | SumTSIPremiSpreadedRNM_Act | GetDataCurrencyByName_SQL, IsMarineCargo | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 82 | Activity | `SumTSIPremiSpreadedRNM_MBU_Act` ⛔ | ASM-FW-GISFW-Data-Coverage | SumTSIPremiSpreadedRNM_Act | GetDataCurrencyByName_SQL, IsEdmAdjSpreading, IsMBU | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 83 | Activity | `SumTSIPremiSpreadedRNM_PA_Act` ⛔ | ASM-FW-GISFW-Data-Coverage | SumTSIPremiSpreadedRNM_Act | GetDataCurrencyByName_SQL, IsPA | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 84 | Activity | `SumTSIPremiSpreadedRNM_TRAVEL_Act` ⛔ | ASM-FW-GISFW-Data-Coverage | SumTSIPremiSpreadedRNM_Act | IsTravel | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 85 | Activity | `TreatyInInputVis` | Data-Portal | SetTreatyIn_Act |  | tidak dibangun — (a) nol efek: dipanggil HANYA `Activity/SetTreatyIn_Act` langkah 2 (`Add=1`, `bypasscondition=1`) — langkah 1 `OutputParam.DATASHOW=1`, `OutputParam.ERRMSG=""`, `TreatyIn.ID="UnknownId"`, lalu transisi langkah 1 `1==1` → 6 (keluar activity), sehingga langkah 3 `Page-Remove TreatyIn` TIDAK berjalan (RALAT audit silang P3 bab 8; bunyi lama: *"langkah 3 `Page-Remove TreatyIn`"* seolah berjalan); `SetTreatyIn_Act` langkah 3 menulis `TreatyIn.ID = Param.ID` dan mengosongkan `OutputParam.ERRMSG`; `OutputParam.DATASHOW` dibaca NOL rule korpus ⇒ tak satu pun hasilnya sampai ke NB |
| 86 | Activity | `TreatyInNonSetTotal` | Data-Portal | DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM |  | tidak dibangun — (a) satu-satunya pemicu tombol `.pyTemplateButton` di DetailPolicyTreatyInNonProportional / ...EDM bersyarat tampil `1=2` |
| 87 | Activity | `TreatyInputPctCommSpreading` | ASM-FW-GISFW-Work | InputPolicyTreatyInDetail_preACT | BreakDownSpreading_Act, CountSpreading_Act, FetchMasterTreatyIn | dibangun sebagian — `models.TreatyInputPctCommSpreading` + `repository.KomisiKontrak` dari `services.PilihBisnis` (preACT langkah 17, bukan NonProportional): langkah 2-2.1.1.1 RiCommOgp = @replaceAll(RIOGR) lalu ditimpa @replaceAll(RIONR) dari kolom view TREATYINDETAILJOINEDM (syarat TREATYTYPE/TREATYGROUP); tiga baris SpreadingRiskList(1) <- SpreadingTotalPct/SpreadingTypeID/SpreadingType tidak — bukan kolom view, hanya JSON master, (c): jalur Proporsional di luar ukuran K8 tunggal (P9; disetujui WO 04-10-2026 (F1)): medan master yang DIBACA rule terjangkau jalur NB NonProp, setiap medan di `models/masterxol.go` berkutip langkah XML; langkah 3 BreakDownSpreading (K9) dan 4 (`//`) tidak |
| 88 | Activity | `TreatyNonPropOutSetSpreading` | ASM-FW-GISFW-Work | InputPolicyTreatyOutDetail_NonProp | CountSpreading_Act, GetCurrencyIDByName | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — pemanggil tunggal `Activity/InputPolicyTreatyOutDetail_NonProp` langkah 28 (Call) — spreading treaty keluar dari halaman master JSON `M_TREATY_OUT` (`select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`); tidak dibaca, tidak ditulis (AC 57 ⛔, AC 58 ✅) |
| 89 | Activity | `TreatyNonPropSetSpreading` | ASM-FW-GISFW-Work | InputPolicyTreatyInDetail_NonProp | CountSpreading_Act, GetCurrencyIDByName | dibangun — [keputusan work owner] K8 (03-10-2026): `models.TreatyNonPropSetSpreading` (1, 3-4, 6-7, 10 -> CountSpreading_Act); 2 nol pembaca; 5 (`ClaimType=="XOL Retro"` -> `ShareReins` treaty KELUAR) ⛔ K8 butir 4; 8-9 `//` (RetroList) |
| 90 | Activity | `TreatyRealizationCheckDuplicate` | ASM-FW-GISFW-Work | InputPolicyTreatyInPost_Act | TreatyRealizationCheckDuplicate | dibangun — `repository.PolisSerupa` + `services.cekDuplikat` (langkah 1-6, pesan verbatim; syarat pemanggil IsApproved==1 hanya pasca admin; pemetaan parameter apa adanya; Totaltsi kosong = nol baris) — AC 59 bunyi baru (putaran 2); uji TestPolisSerupaMenahan, TestCekPolisSerupaHanyaSaatAdminMenyetujui |
| 91 | Activity | `TreatyRealizationCheckXOLList` | ASM-FW-GISFW-Work | InputPolicyTreatyInPre_Act | InsertToTreatyXOLList, SetTreatyIn_Act | dibangun — [keputusan work owner] K8 (03-10-2026): `services.siapkanNonProp` + `models.AwalCekDaftarXOL`/`LengkapiDaftarXOL` (1-6, 8; pesan VERBATIM "Error fetching XolList"); 7 `OldData.TreatyXOLList` tidak (OldData digantikan OLD_POLIS_ID, NB selalu kosong); ID harfiah 3000004 tidak dipakai (`pyPassCurrentParameterPage=true`) |
| 92 | Activity | `TreatySetReinstatement` | Data-Portal | SetTreatyIn_Act | SetReinstatementPct | dibangun — [keputusan work owner] K8 (03-10-2026): `models.TreatySetReinstatement` (dipanggil SetTreatyIn_Act 13 di rantai TreatyRealizationCheckXOLList); hasilnya halaman master baca-saja, nol pembaca NB, tidak disimpan |
| 93 | DataTransform | `AddToListCommentsPolicyTreatyIn_DT` | ASM-FW-GISFW-Work | DeptHeadTreatyIn_UW_postDT, InboxPolicyTreatyIn_postDT |  | dibangun — `models.TambahCatatan` (+ OperatorID login → AKSES_LOGIN HISTORYAKSEPTASIPRODUCTION, AC 39, 71) |
| 94 | DataTransform | `btnCedingCO_DT` | ASM-FW-GISFW-Data-Quotation | InputQuotation_PreAct |  | tidak dibangun — (a) hanya dipanggil `InputQuotation_PreAct` langkah 2 bersyarat `Param.Acton=="Ceding"`; satu-satunya pengirim Acton di korpus mengirim `Acton=SOB` |
| 95 | DataTransform | `btnSOB_DT` | ASM-FW-GISFW-Data-Quotation | InputQuotation_PreAct |  | dibangun — langkah 1 `Quotation.btnQuotation = "SOB"` = `models.TombolSOB` (dijalankan `models.HasilPostDT` sebelum PostDT); langkah 2 `SearchSOB.CARI1 = ""` tidak dibangun (a): halaman di luar pyWorkPage, nilai autocomplete yang dibaca nol rule |
| 96 | DataTransform | `DeptHeadTreatyIn_UW_postDT` | ASM-FW-GISFW-Work | DeptHeadTreatyIn_UW | AddToListCommentsPolicyTreatyIn_DT | dibangun — `models.PascaAtasan` (cabang Director/GroupLeader tidak — P13, AC 10) |
| 97 | DataTransform | `DeptHeadTreatyInUW_preDT` | ASM-FW-GISFW-Work | DeptHeadTreatyIn_UW |  | dibangun — `models.PraprosesAtasan` (OperatorName = nama tampilan, P33; isApprovedtoDeptHead tidak, AC 64) |
| 98 | DataTransform | `InboxPolicyTreatyIn_postDT` | ASM-FW-GISFW-Work | InboxPolicyTreatyIn | AddToListCommentsPolicyTreatyIn_DT, TestTreatyToFacStatus | dibangun — `models.PascaAdmin` |
| 99 | DataTransform | `InputPolicyTreatyIn_preDT` | ASM-FW-GISFW-Work | InboxPolicyTreatyIn |  | dibangun — `models.PraprosesAdmin` (langkah 12-13 nomor urut antrean diganti keanggotaan, AC 14) |
| 100 | DataTransform | `SearchHierarkiSourceBizAgent_PostDT` | ASM-FW-GISFW-Data-Agent | AgentSourceBizDetails |  | dibangun — `models.TerapkanSumberBisnis` (langkah 1 WHEN btnQuotation=="SOB" 1.1-1.2, langkah 4-5; `ChildCount > 0` -> kosong; `.Leader1` bukan kolom RD -> "") lewat `models.HasilPostDT`. ⭐ F4 (keputusan WO 04-10-2026: ikuti XML; RALAT tiket 11 R2): PostDT hanya menulis clipboard (nol Obj-Save) - `services.PilihSumberBisnis` (`POST /kasus/{id}/pilih-sumber-bisnis`, admin, ClaimType 'XOL Retro') = pencarian TANPA simpan yang menjawab keempat nilai; layar memegangnya (`pegangSumberBisnis`) dan mengirimnya bersama Save/Submit/refresh; `services.terimaSumberBisnis` menerimanya hanya bila sama persis dengan PostDT baris `BrowseAgentHierarkiList_RD` yang dijalankan ulang (R7: satu pola kiriman terkunci `models.KirimanSumberBisnis` + `models.TerimaKirimanTerkunci`, empat medan = struct `models.SumberBisnisPostDT`; tidak cocok -> 422; ClaimType bukan XOL Retro -> terkunci/diabaikan), lalu menyalinnya ke QuotationData sebelum SetPPNPPH (`[penyesuaian sadar]`, `models.SalinKeQuotationData`); hanya SOURCE_OF_BUSINESS berkolom. Uji: `handlers/sumberbisnis_test.go`, `models/sumberbisnis_test.go`. Cabang 2 (`CedingCo`) dan 3 (`CedingCedant`) tidak dibangun (a): di NB btnQuotation hanya pernah "SOB" (`btnCedingCO_DT` tak terjangkau; "CedingCedant" ditulis nol rule) |
| 101 | DataTransform | `setCategoryAttachment_DT` | ASM-FW-GISFW-Work-LIFE | SetCategoryAttach |  | tidak dibangun — (a) `SetCategoryAttach` langkah 3 berjalan hanya bila `pyWorkPage.Quotation.BusinessFac` BUKAN "F"/"T" (prasyarat benar → lewati langkah, kode 3); kasus NB treaty = "T" (`ReportDefinition/GetListOpportunity` filter E `A.Quotation.BusinessFac = T`; `services.BuatKasus`) ⇒ tidak pernah berjalan; kelas `ASM-FW-GISFW-Work-LIFE` |
| 102 | DataTransform | `SystemSetOneYear_DT` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |  | dibangun — `models.SystemSetOneYear` + aksi hitung `SystemSetOneYear` (`services/tindakan.go`) dipicu `.StartDate` layar admin (`frontend/medan.ts`). XML: `pyPreDataTransform` aksi refresh `change` sel `.StartDate` `Section/DetailPolicyTreatyIn` (dapat disunting; di `DetailDeptHeadTreatyIn_UW` `pyReadOnly 1==1` ⇒ tak pernah berubah); langkah 1 `.EndDate = @DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)` (Calendar.add: 29 Feb → 28 Feb). RALAT bunyi putaran 1 (*"[keputusan work owner] P35/AC 34: EndDate kosong diisi HARI INI, bukan +1 tahun"*): P35 hanya memutus pengisian tanggal-akhir-KOSONG saat layar dibuka dan menulis peran rule ini `[terbuka]` (spec §9.2 butir 12) — tidak melarangnya |
| 103 | DataTransform | `TestTreatyToFacStatus` | ASM-FW-GISFW-Work | InboxPolicyTreatyIn_postDT |  | dibangun — `models.TetapkanHasFacOut` (AC 76) |
| 104 | DataTransform | `TreatyEnableDisableInput` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn |  | dibangun — `models.TreatyEnableDisableInput` (urutan WHEN ditiru) — RALAT R6 audit silang P3 (04-10-2026): hasilnya (IsNewPolicyNonProp, QuotationData.ProportionalType) diterima dari layar hanya bila TreatyType XOL dan sama dengan hasil DT atas halaman server (R7: satu pola kiriman terkunci `models.KirimanEnableDisable` + `models.TerimaKirimanTerkunci`, kedua medan dicocokkan sekaligus; bukan XOL diabaikan, lain 422) - W3 |
| 105 | DecisionTable | `BusinessType_DeT` | ASM-FW-GISFW-Work | InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT |  | dibangun — `models.GolongkanJenisUsaha` (128 nilai diuji, AC 19-21) |
| 106 | DecisionTable | `isApproved` | ASM-FW-GISFW-Work | InputRealizationTreatyIn |  | dibangun — `models.Disetujui` (AC 1-4) |
| 107 | Flow | `InputRealizationTreatyIn` | ASM-FW-GISFW-Work |  | SaveJsonPolisTreatyIn_Act, serviceInsertArasapas_act, isApproved, DeptHeadTreatyIn_UW, InboxPolicyTreatyIn, IsSPVCreate, IsSPVTreaty1, IsTreaty1, NopolisEmpty,  … | dibangun — `models.Langkah` + `services.Kirim` — tiga posisi (P13); Sec Head menyetujui -> Dept Head selalu (AC 8, WO) |
| 108 | FlowAction | `AgentSourceBizDetails` | ASM-FW-GISFW-Data-Agent | SourceHierarki | SearchHierarkiSourceBizAgent_PostDT | dibangun — pyPreProcessingTransformRule (dan pyActionTransformRule) `SearchHierarkiSourceBizAgent_PostDT` (klik baris TreeGrid masterDetail) = `services.PilihSumberBisnis` - pencarian tanpa simpan, hasilnya dipegang layar sampai Save/Submit (F4); flow action tanpa section/medan |
| 109 | FlowAction | `DeptHeadTreatyIn_UW` | ASM-FW-GISFW-Work | InputRealizationTreatyIn | InputPolicyTreatyInPre_Act, InsertHistoryAkseptasiPega, DeptHeadTreatyInUW_preDT, DeptHeadTreatyIn_UW_postDT, GeneralDeptHeadTreatyIn_UW | dibangun — layar atasan `LayarKasus.tsx` + services (pra/pasca) |
| 110 | FlowAction | `InboxPolicyTreatyIn` | ASM-FW-GISFW-Work | InputRealizationTreatyIn | InputPolicyTreatyInPost_Act, InputPolicyTreatyInPre_Act, InboxPolicyTreatyIn_postDT, InputPolicyTreatyIn_preDT, GeneralPolicyTreatyIn | dibangun — layar admin `LayarKasus.tsx` + services (pra/pasca) |
| 111 | FlowAction | `InputHistoricalSurveyReport` | ASM-FW-GISFW-Data-Quotation | HistoricalSurveyReportDtl | InputHistoricalSurveyReportDtl | tidak dibangun — `pyEditAction` baris grid `Section/HistoricalSurveyReportDtl`; section `InputHistoricalSurveyReportDtl`. (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 112 | FlowAction | `InputHistoricalSurveyReportUW` | ASM-FW-GISFW-Data-Quotation | HistoricalSurveyReportDtlUW | InputHistoricalSurveyReportDtlUW | tidak dibangun — `pyEditAction` baris grid `Section/HistoricalSurveyReportDtlUW`; section `InputHistoricalSurveyReportDtlUW`. (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 113 | FlowAction | `InstallmentList` | ASM-FW-GISFW-Data-Installment | DetailPolicyTreatyInNonProportional, DetailPolicyTreatyOutNonProportional | InstallmentList | dibangun — `DetailNonProp.tsx` (masterDetail ListInstallment, hanya-baca); disimpan T_POLIS_INSTALMENT_DETAIL |
| 114 | FlowAction | `Installments_ReadOnly` | ASM-FW-GISFW-Data-TreatyInInstallment | DetailPolicyTreatyInNonProportionalEDM | Installments_ReadOnly | tidak dibangun — (a) hanya di DetailPolicyTreatyInNonProportionalEDM (tidak terjangkau) |
| 115 | FlowAction | `PolicyTreatyInDeclineConfirm` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn | SetDueTo_act, PolicyTreatyInDeclineConfirm | dibangun — modal konfirmasi `LayarKasus.tsx` (Yes = kirim) |
| 116 | FlowAction | `ShowPolicyNoTreaty` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW | ShowPolicyNoTreaty_SC | dibangun — modal nomor polis `LayarKasus.tsx` (OK = kirim) |
| 117 | Harness | `BusinessAndSOBList` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn | BusinessAndSOBList | dibangun — `components/PilihBisnis.tsx` (judul jendela showHarness `pyWindowName` "Business And SOB List") + `POST /api/nb-treaty-in/kasus/{id}/bisnis` (isian layar `pySubmitData=Yes` digabung tanpa simpan; posisi admin dan wadah `.ClaimType != 'XOL Retro'`, selain itu 409). RALAT audit silang P3 W1: bunyi lama *"`GET /api/nb-treaty-in/bisnis`"* (tanpa kasus, tanpa saringan) |
| 118 | Harness | `BusinessAndSOBListRetro` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn | BusinessAndSOBListRetro | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — popup dari tombol `Choose Business R` (`.pyTemplateInputBox`, aksi `showHarness` `pyHarnessName=BusinessAndSOBListRetro`, tampil `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy'`) di `Section/DetailPolicyTreatyIn`; memuat `Section/BusinessAndSOBListRetro` (grid RD `BrowseTreatyOutDetail`) → `SetValueRetro_Act` → JSON `M_TREATY_OUT` (`select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`); tidak dibaca, tidak ditulis (AC 57 ⛔, AC 58 ✅) |
| 119 | Harness | `HistoricalSurveyReport` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyIn | HistoricalSurveyReportDtl | tidak dibangun — popup dari tombol `Survey Report` (showHarness popup, tampil `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`, nonaktif `.QuotationData.IsSurveyReport=='No' \|\| ==''`) `Section/DetailPolicyTreatyIn`; menyertakan `HistoricalSurveyReportDtl` (grid `.QuotationData.SurveyReportList`). (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 120 | Harness | `HistoricalSurveyReportUW` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW | HistoricalSurveyReportDtlUW | tidak dibangun — popup dari tombol `Survey Report` (showHarness popup, tampil `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`, nonaktif `.QuotationData.IsSurveyReport=='No' \|\| ==''`) `Section/DetailDeptHeadTreatyIn_UW`; menyertakan `HistoricalSurveyReportDtlUW`. (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 121 | Harness | `SFAPortalOpportunities` | PegaCRM-Portal |  | SFAPortalOpportunitiesHeader, SFAPortal_Opportunities, SFAPortal_OpportunitiesList, SFAPortal_OpportunitiesList_Header | dibangun — `pages/PortalNBTreatyIn.tsx` (daftar + Create) — titik masuk satu-satunya |
| 122 | Harness | `SOB` | ASM-FW-GISFW-Data-OfferTreatyIn | DetailPolicyTreatyIn | SourceHierarki | dibangun — popup `components/PilihSumberBisnis.tsx` (judul pyWindowName 'Source Of Business'), dibuka tombol `Select Source Of Business` layar admin bila ClaimType 'XOL Retro' |
| 123 | RDBList | `AttachmentLife` | ASM-FW-GISFW-Int-OFFERJSON | SetCategoryAttach |  | tidak dibangun — (a) `SetCategoryAttach` langkah 2 berjalan hanya bila `pyWorkPage.Quotation.BusinessFac` BUKAN "F"/"T" (prasyarat benar → lewati langkah, kode 3); kasus NB treaty = "T" (`ReportDefinition/GetListOpportunity` filter E `A.Quotation.BusinessFac = T`; `services.BuatKasus`) ⇒ tidak pernah berjalan; SQL `select * from CATEGORY_ATTACH_REAS where note in (...)` |
| 124 | RDBList | `BrowseClientEmail_SQL` | ASM-FW-GISFW-Int-CLIENT | SearchHierarkiSourceBizAgentTreatyIn_Act |  | tidak dibangun — (a) hanya dipanggil `SearchHierarkiSourceBizAgentTreatyIn_Act` langkah 1 untuk `pyWorkPage.OfferTreatyIn.QuotationData.Email`, yang dibaca nol rule NB |
| 125 | RDBList | `BrowseTreatyIn` | ASM-FW-GISFW-Int-TREATY_IN | SetTreatyIn_Act |  | dibangun — [keputusan work owner] K8 (03-10-2026): `repository.MasterXOLDariJSON` (SQL yang sama: M_TREATY_IN UNION ALL M_TREATY_IN_EDM berkunci nomor kontrak; baca-saja, hanya medan master K8) |
| 126 | RDBList | `BrowseTreatyInDetailJoinEDM` | ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM | InputPolicyTreatyInDetail_preACT |  | tidak dibangun — (c): `select JSONDATA as CLASSOFBUSINESS from pooldata.M_TREATY_IN_DETAIL_EDM where ID={TreatyIn.ID}` (preACT langkah 9; hasilnya di-`adoptJSONObject` langkah 10 dan dibaca langkah 13 `TreatyIn.INSTALLMENT`) — JSON tabel di luar dua tabel pengecualian K8 butir 2 (`M_TREATY_IN` / `M_TREATY_IN_EDM`); nol penulisan JSON (AC 16) |
| 127 | RDBList | `BrowseTreatyInJoinEDM` | ASM-FW-GISFW-Data-PolicyTreatyIn | FetchMasterTreatyIn, InputPolicyTreatyInDetail_NonProp |  | dibangun — [keputusan work owner] K8 (03-10-2026): `repository.MasterXOLDariJSON` (NonProp langkah 5; SQL sama dengan BrowseTreatyIn) |
| 128 | RDBList | `BrowseTreatyOut` | ASM-FW-GISFW-Data-PolicyTreatyIn | InputPolicyTreatyOutDetail_NonProp |  | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — dipanggil HANYA `Activity/InputPolicyTreatyOutDetail_NonProp` langkah 5 (RDB-List); SQL `RDBList/BrowseTreatyOut.xml` `pyBrowseSQL`: `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}` — JSON master treaty keluar; tidak dibaca, tidak ditulis (AC 57 ⛔, AC 58 ✅) |
| 129 | RDBList | `BrowseTreatyOutDetail` | ASM-FW-GISFW-Int-TREATYOUTDETAIL | InputPolicyTreatyOutDetail_preACT |  | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — dipanggil HANYA `Activity/InputPolicyTreatyOutDetail_preACT` langkah 10 (RDB-List); SQL `RDBList/BrowseTreatyOutDetail.xml` `pyBrowseSQL` sama persis dengan `BrowseTreatyOut`: `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}` (AC 57 ⛔) |
| 130 | RDBList | `CategoryAttach_SQL` | ASM-FW-GISFW-Int-OFFERJSON | SetCategoryAttach |  | tidak dibangun — (a) berjalan (`SetCategoryAttach` langkah 1, BusinessFac "T"), SQL `select * from CATEGORY_ATTACH_REAS order by note` (baca saja), tetapi halaman hasilnya (`Attachment`) — pembacanya di korpus hanya `CheckRISLIP`, `CheckRISLIP_EDM`, `CheckSpreadingProtect_ACT`, ketiganya hanya dipanggil `Protection_Act` — sambungan DIPUTUS (varian Fac In: `graf.py` PUTUS, keputusan work owner 23-09-2026 butir 6); nol Section menampilkannya |
| 131 | RDBList | `CheckZipCode_SQL` ⛔ | ASM-FW-GISFW-Int-RW | ProtectCoverage_Act, ProtectFIREMBUPA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectCoverage_Act`, `ProtectFIREMBUPA_Act` — seluruhnya tidak terjangkau |
| 132 | RDBList | `FetchTreatyGroupOLDID` | ASM-FW-GISFW-Int-TREATYGROUP | FetchTreatyGroupOldID |  | dibangun — `repository.OldIDGrupTreaty` |
| 133 | RDBList | `GenerateNoPolicy` | ASM-FW-GISFW-Int-POLISTREATYIN | SaveJsonPolisTreatyIn_Act |  | tidak dibangun — (a) hasil tidak dipakai: `SaveJsonPolisTreatyIn_Act` langkah 4 (tanpa BrowsePage) `SELECT 'RNM-'\|\|…\|\|LPAD(POOLDATA.JSON_POLIS_TREATYIN_SEQ.NEXTVAL,5,'0') AS "HASIL"`; langkah 7 memakai `PolicyTreatyIn.PolicyNo` dari `GeneratePolicyNoTreaty_Act`; deret `JSON_POLIS_TREATYIN_SEQ` dibaca nol rule lain di korpus ⇒ efek NEXTVAL tak teramati (RALAT tiket 02) |
| 134 | RDBList | `GetAllCurrency` ⛔ | ASM-FW-GISFW-Int-CURRENCY | GetCurrencyMaster |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `GetCurrencyMaster` — seluruhnya tidak terjangkau |
| 135 | RDBList | `GetBreakDownSpread_SQL` | ASM-FW-GISFW-Work | BreakDownSpreading_Act |  | tidak dibangun — (b) K9 + KEPUTUSAN-RONDE-12 butir 3b: hanya dari `BreakDownSpreading_Act` langkah 1.3; SQL `select DISTINCT Reinstypeid AS "TreatyType", PCT AS "SharePercentage" from pooldata.proportionalarrg ... and TreatyDescID ='10001'` |
| 136 | RDBList | `GetCategoryOccupation` ⛔ | ASM-FW-GISFW-Int-OCCUPATION | ProtectFIREMBUPA_Act, SetFlagOccupation_ACT |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectFIREMBUPA_Act`, `SetFlagOccupation_ACT` — seluruhnya tidak terjangkau |
| 137 | RDBList | `GetClientID_SQL` | ASM-FW-GISFW-Int-CLIENT | InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT |  | dibangun — `repository.KlienDariNama` |
| 138 | RDBList | `GetCountClaim` | Assign-Worklist | CheckDuplicateOffer |  | tidak dibangun — (a) tidak terjangkau efeknya: dipanggil HANYA `Activity/CheckDuplicateOffer` langkah 6 (RDB-List) — `select count(1) as "CARI1" from DATAPEGA.PC_ASM_FW_GCNMFW_WORK where masterid={TreatyWarning.CAIREINSFACIN}`; `TreatyWarning.CAIREINSFACIN = Param.ID` (langkah 5), dan pemanggil tunggal `CheckDuplicateOffer` = `Activity/SetTreatyIn_Act` langkah 14 TANPA parameter (`pyPassCurrentParameterPage=false`, nol `pyStepsParam`) ⇒ `masterid = NULL` ⇒ hitungan selalu 0 ⇒ langkah 7–8 (`HasilClaim.pxResults(1).CARI1>0`, pesan "Sudah Ada Claim…") tidak pernah benar — RALAT AC 59 |
| 139 | RDBList | `GetCurrency` | ASM-FW-GISFW-Int-CURRENCY | ProtectCurrencyTSI_Act, SetCurrency_act |  | dibangun — `repository.NamaMataUang` |
| 140 | RDBList | `GetCurrencyIDByName` | ASM-FW-GISFW-Int-TREATY_IN | SetTreatyCurrencyID, TreatyNonPropOutSetSpreading, TreatyNonPropSetSpreading |  | dibangun — `repository.IDMataUangDariNama` |
| 141 | RDBList | `GetCurrentDate` | ASM-FW-GISFW-Int-TREATY_IN | SetTreatyIn_Act |  | tidak dibangun — (a) tidak terjangkau: dipanggil HANYA `Activity/SetTreatyIn_Act` langkah 8 (RDB-List `select sysdate as CARI1 from dual`) berprasyarat `param.revisionstate==1`; satu-satunya pemanggil `SetTreatyIn_Act` yang terjangkau dari NB = `Activity/TreatyRealizationCheckXOLList` langkah 4 dengan `pyPassCurrentParameterPage=true`, sehingga isian parameter langkah (`viewstate=1`, `ID`) diabaikan dan halaman parameter hanya membawa `Param.ID` (= NoOffer, langkah 3) — `revisionstate` kosong ⇒ langkah 7–11 (revisi master) tidak pernah jalan (RALAT audit silang P3 bab 8; bunyi lama: *"langkah 4 dengan parameter `viewstate` dan `ID` saja"*) |
| 142 | RDBList | `GetDataAgentByNameNonLife_SQL` ⛔ | ASM-FW-GISFW-Int-AGENT | ProtectCedingCo |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectCedingCo` — seluruhnya tidak terjangkau |
| 143 | RDBList | `GetDataByID_SQL` ⛔ | ASM-FW-GISFW-Int-REINSURANCETYPE | cekSpreadingFactIn |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `cekSpreadingFactIn` — seluruhnya tidak terjangkau |
| 144 | RDBList | `GetDataCurrencyByName_SQL` | ASM-FW-GISFW-Int-CURRENCY | InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp, InsertToTreatyOutXOLList, InsertToTreatyXOLList, InsertToT … |  | dibangun — `repository.IDMataUangDariNama` lewat `services.idMataUangMaster` (tabel/kolom sama dengan GetCurrencyIDByName) |
| 145 | RDBList | `GetDataDoubleCase` ⛔ | ASM-FW-GISFW-Int-policyjson | ProtectCoverage_Act, ProtectFIREMBUPA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectCoverage_Act`, `ProtectFIREMBUPA_Act` — seluruhnya tidak terjangkau |
| 146 | RDBList | `GetFlagReject_SQL` ⛔ | ASM-FW-GISFW-Int-policyjson | Protection_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 147 | RDBList | `GetKodeProdNonLife_SQL` | ASM-FW-GISFW-Int-policyjson | GeneratePolicyNoTreaty_Act |  | dibangun — `inti/backend/penomor.AwalanProduksi("NONLIFE")` |
| 148 | RDBList | `GetKursLimitSpreading_SQL` ⛔ | ASM-FW-GISFW-Int-CURRENCY | SumTSIPremiSpreadedRNM_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 149 | RDBList | `GetObjectItembyName_SQL` ⛔ | ASM-FW-GISFW-Int-V_JN_OBJ_ITEM | ProtectFIREMBUPA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectFIREMBUPA_Act` — seluruhnya tidak terjangkau |
| 150 | RDBList | `GetOldIDBusiness_SQL` | ASM-FW-GISFW-Int-OFFERJSON | InputPolicyTreatyInDetail_preACT, InputPolicyTreatyInPre_Act, InputPolicyTreatyOutDetail_preACT |  | dibangun — `repository.BisnisDariKunci` (+ ORDER BY ID supaya tetap) |
| 151 | RDBList | `GetSequenceNumber_SQL` | ASM-FW-GISFW-Int-policyjson | GeneratePolicyNoTreaty_Act |  | dibangun — `inti/backend/penomor.UrutNomorBerikut` (procedure tidak dipanggil — prinsip o); penomor bersama MENULIS `GENERATE_SEQUENCE_NUMBER` seperti `PROC_GENERATE_SEQUENCE_NUMBER` (RDB ini, dipanggil `GeneratePolicyNoTreaty_Act` pySteps(27)) — disetujui WO 04-10-2026 (F8); catatan diagram F103/F118 "dibaca saja" diralat di `rancangan-tabel-datar-treaty-in.md` §4bis.4 |
| 152 | RDBList | `GetSQLDate` | ASM-FW-GISFW-Work | GeneratePolicyNoTreaty_Act, InputPolicyTreatyInPre_Act |  | dibangun — jam aplikasi `Layanan.jam` (InputPolicyTreatyInPre_Act langkah 3-4) |
| 153 | RDBList | `GETTanggalClosing_SQL` | ASM-FW-GISFW-Int-policyjson | GeneratePolicyNoTreaty_Act |  | dibangun — `inti/backend/penomor.HariClosing` |
| 154 | RDBList | `GetTgl_InputJsonPolis_SQL` ⛔ | ASM-FW-GISFW-Int-policyjson | CheckSpreadingProtect_ACT |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 155 | RDBList | `GetTreatyName` ⛔ | ASM-FW-GISFW-Int-PROPORTIONALARRG |  |  | tidak dibangun — **yatim**: nol pemanggil di seluruh korpus (graf.py) |
| 156 | RDBList | `GetTreatyName_SQL` ⛔ | ASM-FW-GISFW-Int-PROPORTIONALARRG | GetTreatyName |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `GetTreatyName` — seluruhnya tidak terjangkau |
| 157 | RDBList | `GetTreatyName_SQL2` ⛔ | ASM-FW-GISFW-Int-PROPORTIONALARRG | GetTreatyName |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `GetTreatyName` — seluruhnya tidak terjangkau |
| 158 | RDBList | `InsertHistoryAkseptasiPega_Sql` | ASM-FW-GISFW-int-policyjson | InsertHistoryAkseptasiPega |  | dibangun — `repository.CatatRiwayat` (skema eksplisit, tanpa COMMIT, + OPERATORID) |
| 159 | RDBList | `InsertViewSuggest_SQL` | ASM-FW-GISFW-Work | SaveViewSuggest |  | dibangun — `repository/usulan.go` sqlSisipUsulan: 15 kolom, KETERANGAN SUBSTR 3990, TGL_INP TO_DATE, nol COMMIT (transaksi submit), skema lewat Qualify; dibaca balik `bacaUsulan` untuk Section ListSuggest; F3 (04-10-2026): juga penulis salinan SuggestList dokumen lama oleh pemuat tiket 22 — `models.UsulanDokumenLama` (langkah 2.1-2.1.2 SaveViewSuggest; CARI7/BUSINESS_CODE dari QuotationData dokumen) + `repository.SalinUsulanLama` (penjaga dobel menurut IDPEGA, lalu `CatatUsulan`) |
| 160 | RDBList | `SavePolisTreatyIn_SQL` | ASM-FW-GISFW-Int-POLISTREATYIN | SaveJsonPolisTreatyIn_Act |  | tidak dibangun — (b) AC 16 + spec penyimpanan AC 48 (nol stored procedure) + diagram F11: `BEGIN POOLDATA.PEGA_JSON_POLIS_TREATYIN(pzInsKey, PolicyNo, NULL, '0', InputData.CARI21, OperatorID.pyUserIdentifier, InputData.CARI3 /*JSON halaman*/, out); COMMIT; END;` — `JSON_POLIS` tidak termasuk tabel lama yang ditulis menurut diagram (F98); kolom JSON_POLIS yang ditulisnya (IDPEGA=pzInsKey, NOPOLIS=PolicyNo, NOENDORS NULL, PRODKE '0', USERNAME) DIBACA pemuat tiket 22 (`repository.BacaJSONPolis`, `SetelKolomDatarLama`) |
| 161 | RDBList | `SaveTreatyIn` | ASM-FW-GISFW-Int-TREATY_IN | SetTreatyIn_Act |  | tidak dibangun — (b) AC 16 (sistem baru nol tulis JSON) + AC 48 (nol stored procedure) + K8 butir 3: `pyBrowseSQL` = `BEGIN POOLDATA.PEGA_TREATY_IN({InputParam.DATAPEGA}, {TreatyIn.ID}, …)` — PENULIS JSON master lewat procedure; juga (a): dipanggil HANYA `Activity/SetTreatyIn_Act` langkah 11 berprasyarat `param.revisionstate==1`, dan pemanggil NB `TreatyRealizationCheckXOLList` langkah 4 tidak mengirim `revisionstate` ⇒ tidak pernah jalan |
| 162 | RDBList | `SelectSpreadingTreatyInProduction` | ASM-FW-GISFW-Int-REINSURANCETYPE | InputPolicyTreatyInPre_Act |  | dibangun — `repository.DaftarJenisSpreading` |
| 163 | RDBList | `TreatyRealizationCheckDuplicate` | ASM-FW-GISFW-Work | TreatyRealizationCheckDuplicate |  | dibangun — `repository.PolisSerupa` |
| 164 | ReportDefinition | `BrowseAgentHierarkiList_RD` | ASM-FW-GISFW-Int-AGENT | AgentSourceBizTreatyIn_Act |  | dibangun — `repository.DaftarAgenHierarki` (grid; dijalankan ulang saat Save/Submit untuk mencocokkan pilihan yang dipegang layar - F4) / `AgenHierarki` (baris yang diklik) sebagaimana efektif di NB: filter A `.Leader0 = Param.Leader` (Param kosong + pyUseNullIfEmpty) = `LEADER0 IS NULL`, B `STATUSACTIVE IS NULL`, urut ClientName, maks 10000; tabel AGENT lewat db.Qualify; kolom CLIENTNAME/LEADER0/CHILDCOUNT/CLIENTID [data DBA - belum dikonfirmasi] |
| 165 | ReportDefinition | `BrowseAgentNusaRe_RD` | ASM-FW-GISFW-Int-AGENT | SourceHierarki |  | tidak dibangun — (a) sumber autocomplete `SearchSOB.CARI1` di `SourceHierarki`; nilai terpilihnya dibaca NOL rule (`btnSOB_DT`/`btnCedingCO_DT` hanya mengosongkannya) dan grid TreeGrid tidak memakainya |
| 166 | ReportDefinition | `BrowseCedingCo_RD` | ASM-FW-GISFW-Int-AGENT | AgentSourceBizTreatyIn_Act |  | tidak dibangun — (a) hanya disebut `AgentSourceBizTreatyIn_Act` langkah 2 (Param.pyReportName) bersyarat `pyWorkPage.OfferTreatyIn.QuotationData.btnQuotation=="CedingCo"` — ditulis nol rule; langkah 3 pun memanggil `BrowseAgentHierarkiList_RD` secara tetap |
| 167 | ReportDefinition | `BrowseClientName_RD` | ASM-FW-GISFW-Int-AGENT | SetPPNPPH |  | dibangun — `repository.StsPKPAgen` (SetPPNPPH) |
| 168 | ReportDefinition | `BrowseCurrency_RD` ⛔ | ASM-FW-GISFW-Int-CURRENCY | GetCurrencyMaster |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `GetCurrencyMaster` — seluruhnya tidak terjangkau |
| 169 | ReportDefinition | `BrowseCurrencyTreatyIn_RD` | ASM-FW-GISFW-Int-CURRENCY | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |  | dibangun — `repository.DaftarMataUang` (tanpa ITL, AC 54) |
| 170 | ReportDefinition | `BrowseMarketingOfficer_RD` | ASM-FW-GISFW-Int-marketingofficer | DetailPolicyTreatyIn |  | dibangun — `repository.DaftarMO` (MOStatus 1) |
| 171 | ReportDefinition | `BrowseReinsuranceType_RD` | ASM-FW-GISFW-Int-REINSURANCETYPE | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |  | dibangun — `repository.DaftarJenisReas` |
| 172 | ReportDefinition | `BrowseTREATY_IN` | ASM-FW-GISFW-Int-TREATY_IN | CheckDuplicateOffer |  | tidak dibangun — (a) tidak terjangkau efeknya: dirujuk HANYA `Activity/CheckDuplicateOffer` langkah 1 (`Param.pyReportName="BrowseTREATY_IN"`), dan langkah 1–4 aktivitas itu berlabel `//` (dinonaktifkan: Property-Set, `pxShowReport`, kalang pembanding, pesan) ⇒ RD tidak pernah dijalankan — RALAT AC 59 |
| 173 | ReportDefinition | `BrowseTreatyGroup_RD` | ASM-FW-GISFW-Int-TREATYGROUP | FetchTreatyGroupOJK |  | dibangun — `repository.OJKGrupTreaty` |
| 174 | ReportDefinition | `BrowseTreatyInDetail` | ASM-FW-GISFW-Int-TREATYINDETAIL | BusinessAndSOBList, BusinessAndSOBListRetro |  | tidak dibangun — (a) tidak terjangkau: satu-satunya grid NB yang memakainya, S11 `Section/BusinessAndSOBList` (`pgRepPgSubSectionBusinessAndSOBListBB`, kelas `ASM-FW-GISFW-Int-TREATYINDETAIL`), ber-`pyContainerVisibleWhen 1=2` (memo rule *"Hidden the old one, now use treatyindetail join edm"*); `InputPolicyTreatyInDetail_preACT` langkah 1 hanya ber-`pyStepsDescription` "RD BrowseTreatyInDetail", nilainya `Param.pyReportName = "BrowseTreatyJoinEDM"`; sebutan lain hanya `pyGridRDName` sel salinan (bukan sumber grid). 33 `pyUIFields`-nya sama persis dengan `BrowseTreatyJoinEDM`. RALAT audit silang P3 W1: bunyi lama *"dibangun — `repository.DaftarDetailKontrak` (grid popup)"* (tabel TREATYINDETAIL; kode dibuang) |
| 175 | ReportDefinition | `BrowseTreatyJoinEDM` | ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM | InputPolicyTreatyInDetail_preACT, BusinessAndSOBList |  | dibangun — `repository.DetailKontrak` (view TREATYINDETAILJOINEDM, filter G `.ID`; preACT langkah 1-3) dan `repository.DaftarBisnis` (grid AKTIF popup `Section/BusinessAndSOBList` S16/S17: filter H `.PROPORTIONTYPE = Param.PROPORTIONALTYPE` = `PolicyTreatyIn.QuotationData.ProportionalType`, kosong diabaikan - tanpa `pyUseNullIfEmpty`; urut `.TREATYID` ASC; `pyMaxRecords` 500). Uji: `handlers/daftarbisnis_test.go`, `repository/kontrak_db_test.go` `TestDaftarBisnisRDBrowseTreatyJoinEDM` (db, K11). RALAT audit silang P3 W1: bunyi lama hanya *"`repository.DetailKontrak` (view TREATYINDETAILJOINEDM, filter ID)"* |
| 176 | ReportDefinition | `BrowseTreatyOutDetail` | ASM-FW-GISFW-Int-TREATYOUTDETAIL | InputPolicyTreatyOutDetail_preACT, BusinessAndSOBListRetro |  | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — RD kelas `ASM-FW-GISFW-Int-TREATYOUTDETAIL` (33 kolom sejajar RD treaty masuk, filter `A: .SOBID = Param.SOB` …) — namanya disebut harfiah K8 butir 4 (`BrowseTreatyOut*`); dirujuk HANYA `Activity/InputPolicyTreatyOutDetail_preACT` langkah 2 (`Param.pyReportName`) dan grid `Section/BusinessAndSOBListRetro` (`pyGridRDName`), keduanya jalur treaty keluar di atas; tidak dibaca (AC 57 ⛔, AC 58 ✅) |
| 177 | ReportDefinition | `crmOpportunitiesList` | ASM-FW-SFAGISFW-Work-Opportunity | SFAPortal_OpportunitiesList |  | tidak dibangun — (a) + (b) P44/AC 53: dua rujukan saja: (1) `pyGridRDName` pada sel label kolom 'Offer No' `Section/SFAPortal_OpportunitiesList` — metadata kolom; sumber grid sebenarnya `pyGridProps/pyRDName = GetListOpportunity` (dibangun, `repository.Gudang.DaftarKasus`); (2) parameter `RDName` tombol `Stage view` `Section/SFAPortalOpportunitiesHeader` (activity `SFAPopulateListHeader`, tidak ada di korpus) yang berada di wadah `pyContainerVisibleWhen = 1=2` ⇒ elemen mati. Kolomnya data peluang CRM (OpportunityStage, OpportunityAmount, CloseDate, …), bukan medan NB |
| 178 | ReportDefinition | `GetListOpportunity` | ASM-FW-SFAGISFW-Work-Opportunity | SFAPortal_OpportunitiesList |  | dibangun — sebagian: saringan status terbuka (D, F) + cari filter G `.TextNoQuotation Contains Param.Search` (pyCaseInsensitive; = pengenal kasus `NB-<n>`, filter B) + batas pyMaxRecords 500 (`models.BatasDaftarPortal`, semula 200 — P9) di `repository.DaftarKasus`, kini di balik gerbang wadah `ReasTreatyInAdmin` (`services.DaftarKasus`, P8). XML: kelas `ASM-FW-SFAGISFW-Work-Opportunity` INNER JOIN `A` = `ASM-FW-GISFW-Work-NB` (`A.pzInsKey = .NBHandle`; `B` Assign-Worklist hanya di Pages&Classes, tidak di-join — BUKAN daftar assignment); logika filter `B AND D AND E AND F AND A AND (C OR G) AND F1`; pyMaxRecords 500. Filter `A` (`A.pxCreateOperator = Param.UserIdentifier`, dari `InputParam.CARI38` yang ditulis NOL rule korpus; tanpa pyUseNullIfEmpty) tidak dibangun — antrean bersama AC 11 (penyimpangan putaran 1); filter C `.Name Contains Param.Search` tidak dibangun — `.Name` milik kelas CRM opportunity, ditulis NOL rule korpus dan tak berkolom di diagram grilling (butir terbuka, bab 0 butir 11); nama bisnis/tertanggung BUKAN medan pencarian RD (RALAT P9 atas pencarian putaran 1 yang lebih luas). Uji `handlers/portal_test.go` `TestDaftarPortalSesuaiGetListOpportunity` |
| 179 | Section | `BusinessAndSOBList` | ASM-FW-GISFW-Data-PolicyTreatyIn | BusinessAndSOBList, DetailPolicyTreatyIn | SetValue_Act, BrowseTreatyInDetail, BrowseTreatyJoinEDM | dibangun — `components/PilihBisnis.tsx`: satu grid tampil (S16/S17 RD `BrowseTreatyJoinEDM`; S11 RD `BrowseTreatyInDetail` dan LABEL "For Treatyindetail join EDM" ber-`1=2`), judul 25 kolom VERBATIM (`labels.ts` `KOLOM_BISNIS`), urut/saring per kolom di klien (`pyGridSorting`/`pyGridFiltering`, `pyColumnFilteringDropDown`; urut awal `.TREATYID` ASC), 50 baris per halaman (`pyPageMode Numeric`, `pyPageSize 50`) - `gridbisnis.ts`; tombol Choose `SetValue_Act` -> `pilih-bisnis`. RALAT audit silang P3 W1: bunyi lama *"dibangun — `components/PilihBisnis.tsx`"* membaca tabel TREATYINDETAIL, judul = nama kolom view, kotak cari TREATYID + tombol Filter buatan |
| 180 | Section | `BusinessAndSOBListRetro` | ASM-FW-GISFW-Data-PolicyTreatyIn | BusinessAndSOBListRetro, DetailPolicyTreatyIn | SetValueRetro_Act, SetValue_Act, BrowseTreatyInDetail, BrowseTreatyOutDetail | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — isi `Harness/BusinessAndSOBListRetro` (popup tombol `Choose Business R` `Section/DetailPolicyTreatyIn`): grid `pyGridRDName` RD `BrowseTreatyOutDetail`, aksi baris `runActivity SetValueRetro_Act` → `InputPolicyTreatyOutDetail_preACT` → JSON `M_TREATY_OUT` (`select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`); tidak dibaca, tidak ditulis (AC 57 ⛔, AC 58 ✅) |
| 181 | Section | `DetailDeptHeadTreatyIn_UW` | ASM-FW-GISFW-Data-PolicyTreatyIn | GeneralDeptHeadTreatyIn_UW, ListSuggest | CountNetPremi_act, CountOverridingCommOgp_Act, CountOverridingCommOnp_Act, CountPctInstallment_Act, CountResult1Onp_Act, CountResult1_Act, CountResult2Ogp_act,  … | dibangun — `LayarKasus.tsx` + `medan.ts` MEDAN_ATASAN_* (seluruhnya hanya-baca: `pyReadOnly`, atau `pyDisabled=always` untuk DueTo/FlagPPH/NoOfferSlip); yang dapat diisi hanya `ListSuggest` (AC 52, `models.medanAtasan`); wadah dicocokkan ulang 2026-10-03; tiga tombol Submit bersyarat identitas (`<ID-operator-1>`, `==`/`!=` + `LetterNo`) = tiga tempat berperan tiket 05 (`models.TempatDetailDHSubmit*`), DIGANTI posisi kasus (`models.TombolUntuk`; AC 8, K2, P13) — RALAT R6 audit silang P3 (04-10-2026): selain ListSuggest, grid SpreadingRiskList subsection NonProp dapat diisi bila FacultativeShare 0/'' (W2, RALAT AC 52; grid Proporsional S96 tetap hanya-baca); bagian uang dikelompokkan per wadah S91/S93/S94 dengan LABEL Heading 4 "OGP"/"ONP", judul panel buatan dibuang (W6) |
| 182 | Section | `DetailPoliciesNonProportional` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, GeneralDeptHeadTreatyIn_UW, GeneralPolicyTreatyIn | DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, TreatyMasterInEDM | dibangun — [keputusan work owner] K8 (03-10-2026): `frontend/components/DetailNonProp.tsx` (tampil `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'`; subsection non-EDM). LABEL "NON EDM" / "EDM" (`OperatorID.pxInsName = '<ID-operator-2>'`) = dua tempat berperan tiket 05 (`models.TempatLabelNonEDM`/`TempatLabelEDM`, `frontend/tempat.ts`) — digerbang pemetaan tempat (tertunda selama IAM belum menjawab), bukan disembunyikan keras; subsection EDM tidak terjangkau (lihat `DetailPolicyTreatyInNonProportionalEDM`) |
| 183 | Section | `DetailPolicyTreatyIn` | ASM-FW-GISFW-Data-PolicyTreatyIn | GeneralPolicyTreatyIn, ListSuggest | CalculatePremi_Act, CheckDataMkt, CountOGPONP_Act, CountPctInstallment_Act, CountResult1Onp_Act, CountResult1_Act, CountResult2Ogp_act, CountResult2Onp_act, Cou … | dibangun — `LayarKasus.tsx` + `medan.ts` MEDAN_ADMIN_*; sel + wadah (`pyContainerVisibleWhen`) dicocokkan ulang 2026-10-03 (paket P2), aksi sel berurutan (`hitung` `urutan`), format sajian K14 (`sajian.ts`); padanan tombol `alat/tombol.json`. Tombol `Select Source Of Business` — dibangun paket P3 (`components/PilihSumberBisnis.tsx`); `Choose Business R` + subsection NonProp — paket P5; `Survey Report` — K7 tidak dibangun; grid `Breakdown Spreading` — K9 — RALAT R6 audit silang P3 (04-10-2026): kiriman layar diterima hanya dari sel/wadah TAMPIL (S19, S7, S14, IDCurrency, IsSurveyReport - W4); hasil Enable/Disable hanya hasil DT (W3); ListInstallment / PremiumSpreaded dari server (W5); pyDefaultValue IsSurveyReport "No" / TypeTax "Inclusive" diterapkan (W5); FlagPPH berubah -> RemoveTypeTax_ACT di server; spanduk NBStatus dan judul panel buatan dibuang, LABEL "OGP"/"ONP" S24/S25 (W6) |
| 184 | Section | `DetailPolicyTreatyInNonProportional` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPoliciesNonProportional | TreatyInNonSetTotal, InstallmentList, SpreadingRiskList | dibangun — [keputusan work owner] K8 (03-10-2026): `DetailNonProp.tsx` (grid Limits/Share/Total Spreaded/Share Facultative baca-saja, SpreadingRiskList, ListInstallment + InstallmentList); tombol TreatyInNonSetTotal tampil `1=2` tidak — RALAT R6 audit silang P3 (04-10-2026): format sel: LimitShareSummaryList `.Limit` tanpa pyDecimalPlaces, LimitFacShareSummaryList tanpa Format (apa adanya); kolom pertama kosong TotalLimitIOONP (S7) / TotalFacShareRnmNP (S57); S73/S74 tanpa judul panel (W6) |
| 185 | Section | `DetailPolicyTreatyInNonProportionalEDM` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPoliciesNonProportional | TreatyInNonSetTotal, Installments_ReadOnly | tidak dibangun — (a) tampil bila `TreatyMasterInEDM && IsEDMInputOnNB != true`: mustahil sesudah NonProp langkah 10 |
| 186 | Section | `DetailPolicyTreatyOutNonProportional` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, GeneralDeptHeadTreatyIn_UW, GeneralPolicyTreatyIn | InstallmentList, SpreadingRiskList | tidak dibangun — (b) [keputusan work owner] **K8 butir 4** (03-10-2026): jalur treaty KELUAR (`BrowseTreatyOut*`, `M_TREATY_OUT`) tetap ⛔ — subsection yang disertakan `Section/DetailPolicyTreatyIn`, `Section/GeneralPolicyTreatyIn`, `Section/DetailDeptHeadTreatyIn_UW`, `Section/GeneralDeptHeadTreatyIn_UW`; tampil bila `.IsNewPolicyNonProp = 1 && .ClaimType = 'XOL Retro'`; isinya halaman hasil `InputPolicyTreatyOutDetail_NonProp` dari JSON `M_TREATY_OUT` (`select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`); tidak dibaca, tidak ditulis (AC 57 ⛔, AC 58 ✅) |
| 187 | Section | `GeneralDeptHeadTreatyIn_UW` | ASM-FW-GISFW-Work | DeptHeadTreatyIn_UW, DetailDeptHeadTreatyIn_UW | DetailDeptHeadTreatyIn_UW, DetailPoliciesNonProportional, DetailPolicyTreatyOutNonProportional, ListSuggest | dibangun — `LayarKasus.tsx`; nol sel sendiri (hanya menyertakan `DetailDeptHeadTreatyIn_UW` atas `.PolicyTreatyIn`) — AC 52. Berkas XML-nya memuat salinan sertakan tiga tombol Submit bersyarat identitas = tiga tempat berperan tiket 05 (`models.TempatGeneralDHSubmit*`), nasibnya sama dengan `DetailDeptHeadTreatyIn_UW` (posisi kasus) |
| 188 | Section | `GeneralPolicyTreatyIn` | ASM-FW-GISFW-Work | InboxPolicyTreatyIn | DetailPoliciesNonProportional, DetailPolicyTreatyIn, DetailPolicyTreatyOutNonProportional, ListSuggest | dibangun — `LayarKasus.tsx` |
| 189 | Section | `HistoricalSurveyReportDtl` | ASM-FW-GISFW-Data-PolicyTreatyIn | HistoricalSurveyReport | ConcatSlipOfferNo_Act, SetSurveyReport_Act, InputHistoricalSurveyReport | tidak dibangun — grid `.QuotationData.SurveyReportList`: tombol Add/Delete baris, `.DateofSurvey` (refresh `ConcatSlipOfferNo_Act`), `.SurveyedBy`, `.LossPrevention`, `.Remarks`, Submit (`SetSurveyReport_Act` + `window.close`). (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 190 | Section | `HistoricalSurveyReportDtlUW` | ASM-FW-GISFW-Data-PolicyTreatyIn | HistoricalSurveyReportUW | ConcatSlipOfferNo_Act, InputHistoricalSurveyReportUW | tidak dibangun — grid `.QuotationData.SurveyReportList` tanpa Add/Delete/Submit (`.DateofSurvey`, `.SurveyedBy`, `.Remarks`). (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 191 | Section | `InputHistoricalSurveyReportDtl` | ASM-FW-GISFW-Data-Quotation | InputHistoricalSurveyReport |  | tidak dibangun — `.DateofSurvey` (wajib, satu-satunya medan wajib di luar dua layar realisasi dan ListSuggest — AC 45), `.Remarks`, `.SurveyedBy`. (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 192 | Section | `InputHistoricalSurveyReportDtlUW` | ASM-FW-GISFW-Data-Quotation | InputHistoricalSurveyReportUW |  | tidak dibangun — `.DateofSurvey`, `.Remarks`, `.SurveyedBy`, ketiganya `pyReadOnly`. (b) K7 [keputusan work owner 03-10-2026] + bab 0 butir 11: data survei (`PolicyTreatyIn.QuotationData.SurveyReportList`: DateofSurvey, SurveyedBy, LossPrevention, Remarks) hanya hidup di halaman kerja (dokumen JSON); ⛔ tidak ada tabel di diagram grilling (sheet `NB Treaty In Prop`/`NonProp`) — butir terbuka untuk grilling berikutnya (AC 45 `DateofSurvey`) |
| 193 | Section | `InstallmentList` | ASM-FW-GISFW-Data-Installment | InstallmentList |  | dibangun — `DetailNonProp.tsx` (rincian per mata uang, hanya-baca — pyReadOnly) |
| 194 | Section | `Installments_ReadOnly` | ASM-FW-GISFW-Data-TreatyInInstallment | Installments_ReadOnly |  | tidak dibangun — (a) hanya di DetailPolicyTreatyInNonProportionalEDM (tidak terjangkau) |
| 195 | Section | `ListSuggest` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, GeneralDeptHeadTreatyIn_UW, GeneralPolicyTreatyIn | CekLimitTreatyAcc_Act, Protection_Act, SetDueTo_act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn | dibangun — panel Suggest `LayarKasus.tsx`; `.ProductionDate`: syarat tampil (`pyVisible`) dan wajib (`pyRequiredWhen`) `.IsApproved == 1 && (<ID-operator-3> \|\| <ID-operator-4>)` = EMPAT tempat berperan tiket 05 (`models.TempatProduksiTampil*`/`TempatProduksiWajib*`, `TanggalProduksiTampil`/`TanggalProduksiWajib`); diterima server hanya bila tampil (`models.GabungMasukanLayar`), wajibnya lewat SATU sumber `models.MedanWajibBerlaku` (layar, Save, submit); tertunda selama pemetaan IAM kosong (K12, K16); `.IsApproved` -> `SetDueTo` (CekLimitTreatyAcc_Act K2, Protection_Act tak ada di korpus) |
| 196 | Section | `PolicyTreatyInDeclineConfirm` | ASM-FW-GISFW-Data-PolicyTreatyIn | PolicyTreatyInDeclineConfirm |  | dibangun — modal `LayarKasus.tsx` |
| 197 | Section | `SFAPortal_Opportunities` | Data-Portal | SFAPortalOpportunities, SFAPortalOpportunitiesHeader | SFAPortal_OpportunitiesList, SFAPortal_OpportunitiesList_Header | dibangun — `PortalNBTreatyIn.tsx` |
| 198 | Section | `SFAPortal_OpportunitiesList` | Data-Portal | SFAPortalOpportunities, SFAPortal_Opportunities | GetListOpportunity, crmOpportunitiesList, IsOperatorLife | dibangun — `PortalNBTreatyIn.tsx` (saringan + daftar) + gerbang `services.DaftarKasus` (P8). XML: baris saringan (`.FilterTermForOpportunity`, tombol Filter) tampil selalu; SATU grid (badan REPEATING `pgRepPgSubSectionSFAPortal_OpportunitiesListBBBBBBBB.pxResults`, `pyGridProps/pyRDName = GetListOpportunity`, param `UserIdentifier = InputParam.CARI38`, `Search = .FilterTermForOpportunity`; kolom `.TextNoQuotation`, `.Name` (pxLink `openWorkByHandle` `.pzInsKey`), `A.Quotation.BusinessName`, `A.Quotation.InsuredName`, `A.Quotation.MarketingName`, `A.NBStatus`) bersarang di wadah `pyContainerVisibleWhen` `OperatorID.pyWorkGroup!='ReasLife' && OperatorID.pyWorkBasketList(2).pyWorkBasketName=='ReasTreatyInAdmin'` lalu wadah `!IsOperatorLife`. Padanan: anggota `ReasTreatyInAdmin` menurut NAMA (AC 14, bukan urutan ke-2) melihat semua kasus terbuka; pelaku lain tidak melihat grid itu — tugas Sec Head/Dept Head dirutekan `Flow/InputRealizationTreatyIn` ke workbasket (Assignment4/6 `ReasTreatyInSecHead`, Assignment3 `ReasTreatyInDeptHead`, ToWorkBasket); daftar kerja workbasket Pega TIDAK ada di korpus ⇒ di portal tunggal (bab 0 butir 7) mereka melihat HANYA kasus di posisi yang dipegang (AC 11, 92); tanpa posisi tangga ⇒ 403. Saringan `.FilterTermForOpportunity` → `Param.Search` = filter G RD (pengenal kasus, tanpa beda huruf), batas 500 baris (pyMaxRecords). Uji `handlers/portal_test.go` `TestGerbangDaftarPortal`, `TestDaftarPortalSesuaiGetListOpportunity`. ⛔ Tidak dibangun: klausa `pyWorkGroup!='ReasLife'` (work group Pega tak berpadanan di `inti.Pelaku`/M_LOGIN_GO, K12 kosong — tiket 05); nonaktif tautan `.Name` (`.pxPages(A).QuotationData.MarketingName != InputParam.CARI34 && InputParam.CARI38 != .pxCreateOperator && pyWorkBasketList(2) != 'ReasFacInTeamLeader'`) — `InputParam.CARI34/CARI38` ditulis NOL rule korpus (hanya dibaca di section ini), (a) nilainya tak terbaca ⇒ tautan buka kasus tetap aktif — RALAT R6 audit silang P3 (04-10-2026): kolom "Position"/"No Polis" (tanpa sel XML) dibuang; kolom "Name" (`.Name` ditulis nol rule) tidak dirender, tautan openWorkByHandle di sel "Offer No" - penyimpangan sadar (tiket 11 RALAT W6); ikon pengosong C[1.2] (setValue "" tanpa refresh) dibangun |
| 199 | Section | `SFAPortal_OpportunitiesList_Header` | Data-Portal | SFAPortalOpportunities, SFAPortal_Opportunities |  | tidak dibangun — (b) P44/AC 53 elemen mati: satu-satunya penyertanya sel SUB_SECTION `Section/SFAPortal_Opportunities` bertampil `1=2` (Harness `SFAPortalOpportunities` memuatnya hanya lewat salinan tertanam section itu); isinya ringkasan `SFAResults` (Prospek Bisnis, Estimasi Total Premi, Estimasi Tanggal Penutupan) yang diisi `SFAPopulateListHeader` (tidak ada di korpus) |
| 200 | Section | `SFAPortalOpportunitiesHeader` | PegaCRM-Portal | SFAPortalOpportunities | SFAPortal_Opportunities, IsNotAdmin, crmCreateOpportunity, isSellingModeB2B, isSellingModeB2BB2C, isSellingModeB2C, pyIsIpadOrDesktop | dibangun — sebagian: tombol Create; Stage view/List view milik CRM — portal CRM (opportunity/SellingMode) — milik modul CRM, bukan alur realisasi; portal NB Treaty In dibangun di `frontend/pages/PortalNBTreatyIn.tsx` |
| 201 | Section | `ShowPolicyNoTreaty_SC` | ASM-FW-GISFW-Data-PolicyTreatyIn | ShowPolicyNoTreaty |  | dibangun — modal nomor polis `LayarKasus.tsx` |
| 202 | Section | `SourceHierarki` | ASM-FW-GISFW-Data-OfferTreatyIn | SOB | AgentSourceBizTreatyIn_Act, SearchHierarkiSourceBizAgentTreatyIn_Act, AgentSourceBizDetails, BrowseAgentNusaRe_RD | dibangun — TreeGrid (kolom 'Source of Business Name' `.ClientName`, klik baris = PostDT tanpa simpan, hasilnya dipegang layar sampai Save/Submit (F4), tombol `Choose` pyVisible `.ChildCount = 0`) di `components/PilihSumberBisnis.tsx`; autocomplete `Search` dan perluasan simpul tidak dibangun (a) — lihat `BrowseAgentNusaRe_RD`, `AgentSourceBizTreatyIn_Act` |
| 203 | Section | `SpreadingRiskList` | ASM-FW-GISFW-Data-PolicyTreatyIn | DetailPolicyTreatyInNonProportional, DetailPolicyTreatyOutNonProportional | CountSpreading_Act | dibangun (K8) — section ini hanya disertakan subsection NonProp (`DetailPolicyTreatyInNonProportional`, `DetailPolicyTreatyOutNonProportional`): `frontend/components/DetailNonProp.tsx` (Add/Delete dan TreatyType/%Share terbuka bila TreatyIn.FacultativeShare = 0/'', pyDecimalPlaces 2; CountSpreading_Act per baris); varian treaty keluar tidak (K8 butir 4). Grid spreading layar Prop = sel milik `DetailPolicyTreatyIn` / `DetailDeptHeadTreatyIn_UW` (`LayarKasus.tsx`, 4 desimal) — RALAT R6 audit silang P3 (04-10-2026): disertakan juga layar ATASAN (`DetailDeptHeadTreatyIn_UW` S88 -> DetailPoliciesNonProportional, pyEditOptions=Auto): Add/Delete dan %Share terbuka bagi atasan dengan syarat SAMA (W2, RALAT spec AC 52); PremiumSpreaded/ClaimSpreaded Read-only dihitung server (W5); judul panel buatan "Spreading Risk" dibuang (S73 NOHEADER, W6) |
| 204 | When | `crmCreateOpportunity` | PegaCRM-Work- | SFAPortalOpportunitiesHeader |  | tidak dibangun — (b) AC 81 + K12: syarat tampil tiga tombol `Create opportunity` (`crmCreateOpportunity && isSellingMode* && IsNotAdmin`, wadah `pyIsIpadOrDesktop`) di `Section/SFAPortalOpportunitiesHeader`; `(crmBypassOperatorAccessChecks \|\| Declare_crmOperatorAccess.canCreateOpportunity) && crmIsOpen` — kedua When dan halaman Declare TIDAK ada di korpus (graf: rujukan hilang) ⇒ izin CRM tidak dapat diport dari XML; padanannya peran, tertunda IAM |
| 205 | When | `isAllRisk` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act, IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 206 | When | `IsAneka` ⛔ | ASM-FW-GISFW-Data | CheckSpreadingProtectAnekaGolf_ACT, ProtectCoverage_Act, SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_Act | IsBoiler, IsBondingAndCustomBonds, IsCIS, IsCIT, IsCar, IsContractorsPlantMachinery, IsCrime, IsEar, IsEnvironmental, IsExclusion, IsGlass, IsGrowingTrees, IsHE … | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtectAnekaGolf_ACT`, `ProtectCoverage_Act`, `SumTSIPremiSpreadedRNM_ANEKA_Act`, `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 207 | When | `isApproved` ⛔ | ASM-FW-GISFW-Work |  |  | tidak dibangun — **yatim**: nol pemanggil di seluruh korpus (graf.py) |
| 208 | When | `isAviationHull` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act, IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 209 | When | `IsBillboardNeon` ⛔ | ASM-FW-GISFW-Work | IsLiability |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsLiability` — seluruhnya tidak terjangkau |
| 210 | When | `isBillboardNeonSyariah` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 211 | When | `IsBoiler` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 212 | When | `IsBonding` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT | IsBondingKBG | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 213 | When | `IsBondingAndCustomBonds` ⛔ | ASM-FW-GISFW-Data | CheckSpreadingProtectAnekaGolf_ACT, SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_Act, IsAneka | IsBondingKBG, IsCustomBonds | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtectAnekaGolf_ACT`, `SumTSIPremiSpreadedRNM_ANEKA_Act`, `SumTSIPremiSpreadedRNM_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 214 | When | `IsBondingKBG` ⛔ | ASM-FW-GISFW-Work | IsBonding, IsBondingAndCustomBonds |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsBonding`, `IsBondingAndCustomBonds` — seluruhnya tidak terjangkau |
| 215 | When | `IsBuilderRisk` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act` — seluruhnya tidak terjangkau |
| 216 | When | `isBurglary` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act, IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 217 | When | `IsCar` ⛔ | ASM-FW-GISFW-Work | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 218 | When | `IsCIS` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 219 | When | `IsCIT` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 220 | When | `IsClaim` | ASM-FW-GISFW-Data | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |  | tidak dibangun — (a) tidak pernah benar: `pyWorkPage.pyWorkIDPrefix = "CLM-"`, sedangkan kasus NB berawalan `NB-` (filter B `ReportDefinition/GetListOpportunity` `.TextNoQuotation` "NB-"; kelas `ASM-FW-GISFW-Work-NB`; `models.AwalanKasus`). Pemakai: kunci `.DueDate` dan nonaktif `.InstallmentPercentage` di `Section/DetailPolicyTreatyIn` dan `DetailDeptHeadTreatyIn_UW` — selalu salah ⇒ tidak mengunci |
| 221 | When | `isClaimTreaty` | ASM-FW-GISFW-Work | InputRealizationTreatyIn |  | tidak dibangun — (a) `Flow/InputRealizationTreatyIn` Decision5 tidak punya connector MASUK (nol connector `ke: Decision5`); satu-satunya keluarannya isClaimTreaty → Assignment5 (Claim Manager) → Decision7 → Assignment1 (Director) — sub-graf tanpa jalan masuk (P5, AC 10) |
| 222 | When | `IsCMI` ⛔ | ASM-FW-GISFW-Work | SumTSIPremiSpreadedRNM_ANEKA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act` — seluruhnya tidak terjangkau |
| 223 | When | `IsContractorsPlantMachinery` ⛔ | ASM-FW-GISFW-Work | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 224 | When | `IsCrime` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 225 | When | `IsCustomBonds` ⛔ | ASM-FW-GISFW-Data | IsBondingAndCustomBonds |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsBondingAndCustomBonds` — seluruhnya tidak terjangkau |
| 226 | When | `IsEar` ⛔ | @baseclass | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 227 | When | `IsEDM` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtect_ACT, GetTreatyName, ProtectCoverage_Act, ProtectFIREMBUPA_Act, ProtectPremiPolicy_Act, ProtectShipData_Act, Protection_Act, cekSpreadingFa … |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT`, `GetTreatyName`, `ProtectCoverage_Act`, `ProtectFIREMBUPA_Act`, `ProtectPremiPolicy_Act`, `ProtectShipData_Act` … — seluruhnya tidak terjangkau |
| 228 | When | `IsEdmAdjShareCedant` ⛔ | ASM-FW-GISFW-Work | ProtectShareCedant_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectShareCedant_Act` — seluruhnya tidak terjangkau |
| 229 | When | `IsEdmAdjSpreading` ⛔ | ASM-FW-GISFW-Work | SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_FIRE_Act, SumTSIPremiSpreadedRNM_GOLF_Act, SumTSIPremiSpreadedRNM_MBU_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `SumTSIPremiSpreadedRNM_FIRE_Act`, `SumTSIPremiSpreadedRNM_GOLF_Act`, `SumTSIPremiSpreadedRNM_MBU_Act` — seluruhnya tidak terjangkau |
| 230 | When | `IsEDMRiSlip` ⛔ | ASM-FW-GISFW-Work | Protection_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `Protection_Act` — seluruhnya tidak terjangkau |
| 231 | When | `isElectronicEquipment` ⛔ | ASM-FW-GISFW-Work | SumTSIPremiSpreadedRNM_ANEKA_Act, IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 232 | When | `IsEngineering` ⛔ | ASM-FW-GISFW-Work | SumTSIPremiSpreadedRNM_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_Act` — seluruhnya tidak terjangkau |
| 233 | When | `IsEnvironmental` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 234 | When | `IsErrorSpreading` ⛔ | @baseclass | GetTreatyName |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `GetTreatyName` — seluruhnya tidak terjangkau |
| 235 | When | `IsExclusion` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 236 | When | `isFidelity` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 237 | When | `IsFire` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtectFire_ACT, GetTreatyName, ProtectCurrencyTSI_Act, ProtectFIREMBUPA_Act, Protection_Act, SetFlagOccupation_ACT, SumTSIPremiSpreadedRNM_Act, S … | IsFireStyle1, IsFireStyle2, IsKPR, IsOilGas | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtectFire_ACT`, `GetTreatyName`, `ProtectCurrencyTSI_Act`, `ProtectFIREMBUPA_Act`, `Protection_Act`, `SetFlagOccupation_ACT` … — seluruhnya tidak terjangkau |
| 238 | When | `IsFireStyle1` ⛔ | ASM-FW-GISFW-Work | IsFire |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsFire` — seluruhnya tidak terjangkau |
| 239 | When | `IsFireStyle2` ⛔ | ASM-FW-GISFW-Work | IsFire |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsFire` — seluruhnya tidak terjangkau |
| 240 | When | `IsGlass` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 241 | When | `IsGolfInsurance` ⛔ | ASM-FW-GISFW-Data | CheckSpreadingProtectAnekaGolf_ACT, CheckSpreadingProtect_ACT, ProtectCoverage_Act, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_GOLF_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtectAnekaGolf_ACT`, `CheckSpreadingProtect_ACT`, `ProtectCoverage_Act`, `SumTSIPremiSpreadedRNM_Act`, `SumTSIPremiSpreadedRNM_GOLF_Act` — seluruhnya tidak terjangkau |
| 242 | When | `IsGrowingTrees` ⛔ | @baseclass | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 243 | When | `IsHE` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act, IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 244 | When | `IsKPR` ⛔ | ASM-FW-GISFW-Work | IsFire |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsFire` — seluruhnya tidak terjangkau |
| 245 | When | `IsLandRig` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 246 | When | `IsLiability` ⛔ | ASM-FW-GISFW-Work | SumTSIPremiSpreadedRNM_ANEKA_Act, IsAneka | IsBillboardNeon, IsProductsLiability, IsProfessionalLiability, IsWorkmenCompensation | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 247 | When | `IsLife` ⛔ | ASM-FW-GISFW-Data | ProtectPremiPolicy_Act, ProtectShareCedant_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectPremiPolicy_Act`, `ProtectShareCedant_Act` — seluruhnya tidak terjangkau |
| 248 | When | `IsMaintenance` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 249 | When | `IsMarineCargo` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtectMCargoMBU_ACT, GetTreatyName, ProtectCoverage_Act, Protection_Act, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_MARINECARGO_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtectMCargoMBU_ACT`, `GetTreatyName`, `ProtectCoverage_Act`, `Protection_Act`, `SumTSIPremiSpreadedRNM_Act`, `SumTSIPremiSpreadedRNM_MARINECARGO_Act` — seluruhnya tidak terjangkau |
| 250 | When | `IsMarineHull` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act, IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 251 | When | `IsMarineHullOffshore` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act` — seluruhnya tidak terjangkau |
| 252 | When | `IsMBD` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_Act, IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act`, `SumTSIPremiSpreadedRNM_Act`, `IsAneka` — seluruhnya tidak terjangkau |
| 253 | When | `IsMBU` ⛔ | ASM-FW-GISFW-Data | CheckSpreadingProtectMCargoMBU_ACT, ProtectFIREMBUPA_Act, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_MBU_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtectMCargoMBU_ACT`, `ProtectFIREMBUPA_Act`, `SumTSIPremiSpreadedRNM_Act`, `SumTSIPremiSpreadedRNM_MBU_Act` — seluruhnya tidak terjangkau |
| 254 | When | `IsMBUCar` ⛔ | ASM-FW-GISFW-Data | CheckSpreadingProtect_ACT |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtect_ACT` — seluruhnya tidak terjangkau |
| 255 | When | `IsNotAdmin` | @baseclass | SFAPortalOpportunitiesHeader |  | tidak dibangun — (b) AC 12/81/91 + K12 (pemetaan IAM kosong): hanya syarat tampil tiga tombol `Create opportunity` (`crmCreateOpportunity && isSellingMode* && IsNotAdmin`, wadah `pyIsIpadOrDesktop`) di `Section/SFAPortalOpportunitiesHeader`; `OperatorID.pyPosition != "Admin"` = atribut operator Pega tanpa padanan di `inti.Pelaku` (AkunID + Peran) — menerjemahkannya ke peran = menebak siapa boleh membuat realisasi. Tombol Create dibangun tanpa gerbang (`services.BuatKasus`); butir terbuka IAM |
| 256 | When | `IsObjectWithQuantityYear` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act` — seluruhnya tidak terjangkau |
| 257 | When | `IsOilGas` ⛔ | ASM-FW-GISFW-Work | IsFire |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsFire` — seluruhnya tidak terjangkau |
| 258 | When | `IsOperatorLife` | Data-Portal | SFAPortal_OpportunitiesList |  | tidak dibangun — (b) K12 kosong: dipakai sebagai `!IsOperatorLife` (wadah dalam grid `GetListOpportunity`, `Section/SFAPortal_OpportunitiesList`) — syarat yang SAMA dengan klausa pertama wadah luar `OperatorID.pyWorkGroup!='ReasLife' && OperatorID.pyWorkBasketList(2).pyWorkBasketName=='ReasTreatyInAdmin'`; `OperatorID.pyWorkGroup` = atribut operator Pega tanpa padanan di `inti.Pelaku` (AkunID + workbasket) maupun M_LOGIN_GO (ORGANIZATION/DIVISION/UNIT CODE, pemetaannya tidak ada) ⇒ butir terbuka IAM (tiket 05). Klausa workbasket wadah luar DIBANGUN putaran 2 (P8) di `services.DaftarKasus` menurut nama (AC 14) |
| 259 | When | `IsPA` ⛔ | Data-Party-Person | CheckSpreadingProtectPATravel_ACT, ProtectFIREMBUPA_Act, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_PA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtectPATravel_ACT`, `ProtectFIREMBUPA_Act`, `SumTSIPremiSpreadedRNM_Act`, `SumTSIPremiSpreadedRNM_PA_Act` — seluruhnya tidak terjangkau |
| 260 | When | `IsPortRisk` ⛔ | ASM-FW-GISFW-Data | SumTSIPremiSpreadedRNM_ANEKA_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `SumTSIPremiSpreadedRNM_ANEKA_Act` — seluruhnya tidak terjangkau |
| 261 | When | `IsProductsLiability` ⛔ | ASM-FW-GISFW-Work | IsLiability |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsLiability` — seluruhnya tidak terjangkau |
| 262 | When | `IsProfessionalLiability` ⛔ | ASM-FW-GISFW-Work | IsLiability |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsLiability` — seluruhnya tidak terjangkau |
| 263 | When | `IsRenewal` ⛔ | ASM-FW-GISFW-Work | ProtectCoverage_Act, ProtectFIREMBUPA_Act, ProtectRenewal_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectCoverage_Act`, `ProtectFIREMBUPA_Act`, `ProtectRenewal_Act` — seluruhnya tidak terjangkau |
| 264 | When | `isSellingModeB2B` | @baseclass | SFAPortalOpportunitiesHeader |  | tidak dibangun — (a) tidak mengubah perilaku: tiga When `isSellingMode*` saling meniadakan (`getDataSystemSetting("PegaCRM-","SellingMode")` = B2B / B2C / B2B_B2C; nilainya tidak ada di korpus) dan hanya memilih SATU dari tiga tombol `Create opportunity` berlabel dan beraksi sama (`createWork` `D_crmAppExtPage.CreateWork_Flow`) di `Section/SFAPortalOpportunitiesHeader` ⇒ sistem baru satu tombol Create (`services.BuatKasus`) |
| 265 | When | `isSellingModeB2BB2C` | @baseclass | SFAPortalOpportunitiesHeader |  | tidak dibangun — (a) lihat `When/isSellingModeB2B`: memilih satu dari tiga tombol Create yang setara |
| 266 | When | `isSellingModeB2C` | @baseclass | SFAPortalOpportunitiesHeader |  | tidak dibangun — (a) lihat `When/isSellingModeB2B`: memilih satu dari tiga tombol Create yang setara (varian B2C ber-`WorkClass_OpportunityInd`, nilai `D_crmAppExtPage` tidak ada di korpus) |
| 267 | When | `IsSPVCreate` | ASM-FW-GISFW-Work | InputRealizationTreatyIn |  | tidak dibangun — (b) K5 + AC 12/81 (P12, P28; K12 kosong): Decision6 membandingkan `pyWorkPage.pxCreateOperator` dengan DUA ID operator (identitas orang); keempat cabang Decision6/Decision12/Decision10 `Flow/InputRealizationTreatyIn` berakhir di Assignment4 atau Assignment6 yang identik (workbasket `ReasTreatyInSecHead`, FlowAction `DeptHeadTreatyIn_UW`, `Position "5"`, `PositionNote "ReasTreatyInSecHead"`; Decision4 = Decision11); bedanya hanya teks `NBStatus` bernama orang, yang menurut K5 diganti nama posisi ⇒ nol perilaku hilang (models/tangga.go) |
| 268 | When | `IsSPVTreaty1` | ASM-FW-GISFW-Work | InputRealizationTreatyIn |  | tidak dibangun — (b) K5 + AC 13/81: Decision12 `OperatorID.pyTelephone = "SPVTREATY1"` (peran di kolom telepon) memilih Assignment6 atau Assignment4; keempat cabang Decision6/Decision12/Decision10 `Flow/InputRealizationTreatyIn` berakhir di Assignment4 atau Assignment6 yang identik (workbasket `ReasTreatyInSecHead`, FlowAction `DeptHeadTreatyIn_UW`, `Position "5"`, `PositionNote "ReasTreatyInSecHead"`; Decision4 = Decision11); bedanya hanya teks `NBStatus` bernama orang, yang menurut K5 diganti nama posisi ⇒ nol perilaku hilang (models/tangga.go) |
| 269 | When | `IsTBonding` ⛔ | ASM-FW-GISFW-Work | ProtectCoverage_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `ProtectCoverage_Act` — seluruhnya tidak terjangkau |
| 270 | When | `IsTravel` ⛔ | ASM-FW-GISFW-Work | CheckSpreadingProtectPATravel_ACT, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_TRAVEL_Act |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `CheckSpreadingProtectPATravel_ACT`, `SumTSIPremiSpreadedRNM_Act`, `SumTSIPremiSpreadedRNM_TRAVEL_Act` — seluruhnya tidak terjangkau |
| 271 | When | `IsTreaty1` | ASM-FW-GISFW-Work | InputRealizationTreatyIn |  | tidak dibangun — (b) K5 + AC 13/81: Decision10 `OperatorID.pyTelephone = "TREATY1"` memilih Assignment4 atau Assignment6; keempat cabang Decision6/Decision12/Decision10 `Flow/InputRealizationTreatyIn` berakhir di Assignment4 atau Assignment6 yang identik (workbasket `ReasTreatyInSecHead`, FlowAction `DeptHeadTreatyIn_UW`, `Position "5"`, `PositionNote "ReasTreatyInSecHead"`; Decision4 = Decision11); bedanya hanya teks `NBStatus` bernama orang, yang menurut K5 diganti nama posisi ⇒ nol perilaku hilang (models/tangga.go) |
| 272 | When | `IsUW` | ASM-FW-GISFW-Work | DetailPolicyTreatyIn |  | tidak dibangun — (b) AC 14 [keputusan WO] + AC 81/K12: `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName` = `ReasFacInDirector` / `ReasFacInGroupLeader` / `ReasFacInUnderwriting` (penunjukan menurut NOMOR URUT 1 — dilarang AC 14; operand kanan tanpa tanda kutip). Pemakai: kunci `.DueDate`, `.InstallmentPercentage`, `.Premium`, `.QuotationData.IsSurveyReport`, `.StatementType` di `Section/DetailPolicyTreatyIn`. Ketiga workbasket milik Fac In; Flow NB hanya menugaskan ke `ReasTreatyIn*`. Padanan menurut nama ('anggota ReasFacIn*') memperluas kunci ke setiap pemegang ganda — siapa memegang keduanya = pemetaan IAM ⇒ tertunda, butir terbuka |
| 273 | When | `IsWorkmenCompensation` ⛔ | ASM-FW-GISFW-Work | IsLiability |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsLiability` — seluruhnya tidak terjangkau |
| 274 | When | `IsYieldShortfall` ⛔ | ASM-FW-GISFW-Data | IsAneka |  | tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh `IsAneka` — seluruhnya tidak terjangkau |
| 275 | When | `NopolisEmpty` | ASM-FW-GISFW-Work | InputRealizationTreatyIn |  | dibangun — `models.Langkah` (adaNomorPolis) |
| 276 | When | `pyIsIpadOrDesktop` | @baseclass | SFAPortalOpportunitiesHeader |  | tidak dibangun — (a) syarat wadah tombol Create `Section/SFAPortalOpportunitiesHeader` (`pyContainerVisibleWhen`): `pyIsIPad` (tidak ada di korpus) \|\| `pxRequestor.pxDeviceType = desktop` — deteksi perangkat runtime Pega, benar pada peramban desktop; wadahnya (tombol Create) dibangun tampil-selalu di portal web |
| 277 | When | `ToTREATYDEPTHEAD` | ASM-FW-GISFW-Work | InputRealizationTreatyIn |  | tidak dibangun — (b) K2 [WO lama berlaku]: Decision13 `pyWorkPage.LetterNo = "TREATYINDEPTHEAD"` (diisi `CekLimitTreatyAcc_Act` langkah 2.4 bila premi > 200 juta) memilih Assignment3 (Dept Head) atau Decision8; K2: Sec Head menyetujui SELALU naik ke Dept Head ⇒ syarat ini tidak dievaluasi (`models.Langkah`, tiket 03 RALAT) |
| 278 | When | `TreatyMasterInEDM` | ASM-FW-GISFW-Data-PolicyTreatyIn | InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp, DetailPoliciesNonProportional |  | dibangun — `models.TreatyMasterInEDM` (EDMState 1/2/3, [keputusan work owner] K8 (03-10-2026)) |

⛔ = tidak terjangkau dari titik masuk.

## 3 · Menu dan titik masuk

Dari `Struktur_MenuNBTreatyIn.xlsx` (sheet `Struktur`, 3748 baris, kedalaman 14). Tingkat 1–3 pohon pemanggilan:

| No | Jenis rule | Nama | Dipanggil dari |
| --- | --- | --- | --- |
| 1 | Harness | `SFAPortalOpportunities` | menu portal (entry point) |
| 1.1 | Flow | `InputRealizationTreatyIn` | createWork dari tombol Create di SFAPortalOpportunitiesHeader |
| 1.1.1 | Flow Action | `InboxPolicyTreatyIn` | Assignment2 |
| 1.1.2 | Activity | `SaveJsonPolisTreatyIn_Act` | InputRealizationTreatyIn |
| 1.1.3 | Decision Table | `isApproved` | Decision3 |
| 1.1.4 | Decision Table | `isApproved` | InputRealizationTreatyIn |
| 1.1.5 | When | `IsSPVCreate` | Decision6 (pxCreateOperator SPV?) |
| 1.1.6 | When | `IsSPVCreate` | InputRealizationTreatyIn |
| 1.1.7 | When | `IsSPVTreaty1` | Decision12 (IS SPV TREATY 1) |
| 1.1.8 | When | `ToTREATYDEPTHEAD` | InputRealizationTreatyIn |
| 1.1.9 | When | `IsTreaty1` | Decision10 (is Treaty 1) |
| 1.1.10 | When | `IsSPVTreaty1` | InputRealizationTreatyIn |
| 1.1.11 | Flow Action | `DeptHeadTreatyIn_UW` | Assignment6 |
| 1.1.12 | When | `IsTreaty1` | InputRealizationTreatyIn |
| 1.1.13 | Flow Action | `DeptHeadTreatyIn_UW` | Assignment4 |
| 1.1.14 | When | `NopolisEmpty` | InputRealizationTreatyIn |
| 1.1.15 | Decision Table | `isApproved` | Decision4 |
| 1.1.16 | When | `isClaimTreaty` | InputRealizationTreatyIn |
| 1.1.17 | Decision Table | `isApproved` | Decision11 |
| 1.1.18 | Activity | `serviceInsertArasapas_act` | InputRealizationTreatyIn |
| 1.1.19 | When | `ToTREATYDEPTHEAD` | Decision13 (TO TREATY DEPT HEAD?) |
| 1.1.20 | Flow Action | `DeptHeadTreatyIn_UW` | InputRealizationTreatyIn |
| 1.1.21 | When | `NopolisEmpty` | Decision8 ([Decision]) |
| 1.1.22 | Flow Action | `InboxPolicyTreatyIn` | InputRealizationTreatyIn |
| 1.1.23 | Flow Action | `DeptHeadTreatyIn_UW` | Assignment3 |
| 1.1.24 | Activity | `SaveJsonPolisTreatyIn_Act` | Utility1 |
| 1.1.25 | Decision Table | `isApproved` | Decision2 |
| 1.1.26 | Activity | `serviceInsertArasapas_act` | Utility2 |
| 1.1.27 | Decision Table | `isApproved` | Decision7 |
| 1.1.28 | When | `isClaimTreaty` | Decision5 (Claim) |
| 1.1.29 | Flow Action | `DeptHeadTreatyIn_UW` | Assignment5 |
| 1.1.30 | Decision Table | `isApproved` | Decision1 |
| 1.1.31 | Flow Action | `DeptHeadTreatyIn_UW` | Assignment1 |
| 1.2 | Section | `SFAPortal_Opportunities` | SFAPortalOpportunities |
| 1.2.1 | Section | `SFAPortal_OpportunitiesList` | SFAPortal_Opportunities |
| 1.2.2 | Section | `SFAPortal_OpportunitiesList_Header` | SFAPortal_Opportunities |
| 1.3 | Section | `SFAPortalOpportunitiesHeader` | SFAPortalOpportunities |
| 1.3.1 | When | `crmCreateOpportunity` | SFAPortalOpportunitiesHeader |
| 1.3.2 | When | `isSellingModeB2C` | SFAPortalOpportunitiesHeader |
| 1.3.3 | When | `isSellingModeB2BB2C` | SFAPortalOpportunitiesHeader |
| 1.3.4 | When | `IsNotAdmin` | SFAPortalOpportunitiesHeader |
| 1.3.5 | When | `isSellingModeB2B` | SFAPortalOpportunitiesHeader |
| 1.3.6 | When | `pyIsIpadOrDesktop` | SFAPortalOpportunitiesHeader |

Akar tunggal: Harness `SFAPortalOpportunities` (menu portal). Satu-satunya jalan ke alur realisasi adalah tombol *Create* portal (`createWork` → `Flow/InputRealizationTreatyIn`). Tidak ada butir menu lain di pohon ini — sistem baru tidak boleh menambah menu di luar itu.

## 4 · Alur `Flow/InputRealizationTreatyIn`

Kelas `ASM-FW-GISFW-Work`. **22 shape**, **32 connector**.

### 4.1 Shape

| Id | Jenis | Nama | Implementasi | Workbasket | Dimasuki connector |
| --- | --- | --- | --- | --- | ---: |
| `Decision12` | Gateway-Decision | IS SPV TREATY 1 |  DataXOR |  | 1 |
| `Utility2` | Activity-Utility | HIT SERVICE ARASAPAS | `serviceInsertArasapas_act`  |  | 1 |
| `Decision11` | Gateway-Decision | Is Correct? | `isApproved` Rule-Declare-DecisionTable |  | 1 |
| `Utility1` | Activity-Utility | Save json policy | `SaveJsonPolisTreatyIn_Act`  |  | 1 |
| `Decision13` | Gateway-Decision | TO TREATY DEPT HEAD? |  DataXOR |  | 2 |
| `Decision8` | Gateway-Decision | [Decision] |  DataXOR |  | 3 |
| `Decision7` | Gateway-Decision | Is Correct? | `isApproved` Rule-Declare-DecisionTable |  | 1 |
| `Decision6` | Gateway-Decision | pxCreateOperator SPV? |  DataXOR |  | 1 |
| `Decision5` | Gateway-Decision | Claim |  DataXOR |  | 0 |
| `Decision10` | Gateway-Decision | is Treaty 1 |  DataXOR |  | 1 |
| `Assignment5` | Activity-Assignment | Claim Manager | `WorkBasket`  | `ReasTreatyInGroupLeader` | 1 |
| `Decision4` | Gateway-Decision | Is Correct? | `isApproved` Rule-Declare-DecisionTable |  | 1 |
| `Assignment6` | Activity-Assignment | Acceptance by Head. Treaty | `WorkBasket`  | `ReasTreatyInSecHead` | 2 |
| `Assignment3` | Activity-Assignment | Acceptance by Dept. Head | `WorkBasket`  | `ReasTreatyInDeptHead` | 2 |
| `Decision3` | Gateway-Decision | Is Correct? | `isApproved` Rule-Declare-DecisionTable |  | 1 |
| `Decision2` | Gateway-Decision | Is Correct? | `isApproved` Rule-Declare-DecisionTable |  | 1 |
| `Assignment4` | Activity-Assignment | Acceptance by Head. Treaty | `WorkBasket`  | `ReasTreatyInSecHead` | 2 |
| `Decision1` | Gateway-Decision | Is Correct? | `isApproved` Rule-Declare-DecisionTable |  | 1 |
| `Assignment1` | Activity-Assignment | Acceptance by Dir. | `WorkBasket`  | `ReasTreatyInDirector` | 1 |
| `Assignment2` | Activity-Assignment | Input Realitation | `WorkBasket`  | `ReasTreatyInAdmin` | 6 |
| `End3` |  |  |   |  | 2 |
| `Start1` |  |  |   |  | 0 |

### 4.2 Connector

| Dari | Ke | Nama | Syarat | Ekspresi | Tugas properti |
| --- | --- | --- | --- | --- | --- |
| `Assignment1` | `Decision1` | DeptHeadTreatyIn_UW | Action | `DeptHeadTreatyIn_UW` | .PolicyTreatyIn.Show = False |
| `Assignment2` | `Decision3` | InboxPolicyTreatyIn | Action | `InboxPolicyTreatyIn` | .PolicyTreatyIn.Show = True |
| `Assignment3` | `Decision2` | DeptHeadTreatyIn_UW | Action | `DeptHeadTreatyIn_UW` | .PolicyTreatyIn.Show = False |
| `Assignment4` | `Decision11` | DeptHeadTreatyIn_UW | Action | `DeptHeadTreatyIn_UW` |  |
| `Assignment5` | `Decision7` | DeptHeadTreatyIn_UW | Action | `DeptHeadTreatyIn_UW` | .PolicyTreatyIn.Show = False |
| `Assignment6` | `Decision4` | DeptHeadTreatyIn_UW | Action | `DeptHeadTreatyIn_UW` |  |
| `Decision1` | `Assignment2` | No | Status | `No` | pyWorkPage.Position = "4"; pyWorkPage.PositionNote = "ReasTreatyInAdmin" |
| `Decision1` | `Decision8` | Yes | Status | `Yes` |  |
| `Decision10` | `Assignment4` | IsTreaty1 | When | `IsTreaty1` | pyWorkPage.Position = "5"; pyWorkPage.NBStatus = "NB IS IN <nama-1>'S INBOX"; pyWorkPage.PositionNote = "ReasTreatyInSecHead" |
| `Decision10` | `Assignment6` | [Else] | Else |  | pyWorkPage.Position = "5"; pyWorkPage.PositionNote = "ReasTreatyInSecHead"; pyWorkPage.NBStatus = "NB IS IN <nama-2>'S INBOX" |
| `Decision11` | `Assignment2` | No | Status | `No` | pyWorkPage.Position = "4"; pyWorkPage.PositionNote = "ReasTreatyInAdmin"; pyWorkPage.NBStatus = "" |
| `Decision11` | `Decision13` | YES | Status | `YES` |  |
| `Decision12` | `Assignment4` | [Else] | Else |  | pyWorkPage.Position = "5"; pyWorkPage.PositionNote = "ReasTreatyInSecHead"; pyWorkPage.NBStatus = "NB IS IN <nama-1>'S INBOX" |
| `Decision12` | `Assignment6` | IsSPVTreaty1 | When | `IsSPVTreaty1` | pyWorkPage.Position = "5"; pyWorkPage.NBStatus = "NB IS IN <nama-2>'S INBOX"; pyWorkPage.PositionNote = "ReasTreatyInSecHead" |
| `Decision13` | `Assignment3` | ToTREATYDEPTHEAD | When | `ToTREATYDEPTHEAD` | pyWorkPage.Position = "5"; pyWorkPage.NBStatus = "NB IS IN <nama-3>'S INBOX"; pyWorkPage.PositionNote = "ReasTreatyInDeptHead" |
| `Decision13` | `Decision8` | [Else] | Else |  |  |
| `Decision2` | `Assignment2` | No | Status | `No` | pyWorkPage.Position = "4"; pyWorkPage.PositionNote = "ReasTreatyInAdmin"; pyWorkPage.NBStatus = "" |
| `Decision2` | `Decision8` | YES | Status | `YES` | pyWorkPage.Position = "6"; pyWorkPage.PositionNote = "ReasTreatyInDirector"; pyWorkPage.NBStatus = "NB IS IN <nama-4>'S INBOX" |
| `Decision3` | `Decision6` | YES | Status | `YES` |  |
| `Decision3` | `End3` | No | Status | `No` |  |
| `Decision4` | `Assignment2` | No | Status | `No` | pyWorkPage.Position = "4"; pyWorkPage.PositionNote = "ReasTreatyInAdmin"; pyWorkPage.NBStatus = "" |
| `Decision4` | `Decision13` | YES | Status | `YES` |  |
| `Decision5` | `Assignment5` | isClaimTreaty | When | `isClaimTreaty` | pyWorkPage.Position = "5"; pyWorkPage.PositionNote = "ReasTreatyInGroupLeader"; pyWorkPage.NBStatus = "NB IS IN <nama-5>'S INBOX" |
| `Decision6` | `Decision10` | [Else] | Else |  |  |
| `Decision6` | `Decision12` | IsSPVCreate | When | `IsSPVCreate` | pyWorkPage.Position = 5; pyWorkPage.PositionNote = "ReasTreatyInDeptHead"; pyWorkPage.NBStatus = "NB IS IN <nama-3>'S INBOX" |
| `Decision7` | `Assignment1` | YES | Status | `YES` |  |
| `Decision7` | `Assignment2` | No | Status | `No` | pyWorkPage.Position = "4"; pyWorkPage.PositionNote = "ReasTreatyInAdmin" |
| `Decision8` | `Assignment3` | NopolisEmpty | When | `NopolisEmpty` | pyWorkPage.Position = "5"; pyWorkPage.PositionNote = "ReasTreatyInGroupLeader"; pyWorkPage.NBStatus = "NB IS IN <nama-3>'S INBOX" |
| `Decision8` | `Utility1` | Nopolis not empty | Else |  |  |
| `Start1` | `Assignment2` | [Always] | Always |  | .Position = 4; .FlagOnGoingPolicy = 1; .PositionNote = "ReasTreatyInAdmin" |
| `Utility1` | `Utility2` | [Always] | Always |  |  |
| `Utility2` | `End3` | [Always] | Always |  |  |

**Shape tanpa connector masuk (tidak pernah dicapai):** `Decision5`.

## 5 · Layar — Harness, FlowAction, Section

Setiap medan milik rule itu sendiri (salinan section tertanam `pyIncludedRuleXML` dilompati). *Wajib*: `pyRequired`/mode; *Kunci*: `pyReadOnly` atau syarat baca-saja; *Tampil*: `pyVisible` + syarat; *Aksi*: event → aksi → activity/harness/local action yang dipanggil.

### 5.H `Harness/BusinessAndSOBList`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `BusinessAndSOBList`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | pxIconHistory |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconAttachments |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconExpandCollapse |  |  |  |  |  |
|  | pxIconCancel |  |  |  |  |  |

### 5.H `Harness/BusinessAndSOBListRetro`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `BusinessAndSOBListRetro`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | pxIconHistory |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconAttachments |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconExpandCollapse |  |  |  |  |  |
|  | pxIconCancel |  |  |  |  |  |

### 5.H `Harness/HistoricalSurveyReport`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `HistoricalSurveyReportDtl`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | pxIconHistory |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconAttachments |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconExpandCollapse |  |  |  |  |  |
|  | pxIconCancel |  |  |  |  |  |

### 5.H `Harness/HistoricalSurveyReportUW`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `HistoricalSurveyReportDtlUW`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | pxIconHistory |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconAttachments |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconExpandCollapse |  |  |  |  |  |
|  | pxIconCancel |  |  |  |  |  |

### 5.H `Harness/SFAPortalOpportunities`

Kelas `PegaCRM-Portal`
- menyertakan: `SFAPortalOpportunitiesHeader`, `SFAPortal_Opportunities`, `SFAPortal_OpportunitiesList`, `SFAPortal_OpportunitiesList_Header`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | pxIconHistory |  |  |  | OTHER `.pyID!=""` |  |
|  | pxIconAttachments |  |  |  | OTHER `.pyID!=""` |  |
|  | pxIconExpandCollapse |  |  |  |  |  |
|  | pxIconCancel |  |  |  |  |  |

### 5.H `Harness/SOB`

Kelas `ASM-FW-GISFW-Data-OfferTreatyIn`
- menyertakan: `SourceHierarki`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
|  | pxIconHistory |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconAttachments |  |  |  | OTHER `.pyID != ''` |  |
|  | pxIconExpandCollapse |  |  |  |  |  |
|  | pxIconCancel |  |  |  |  |  |

### 5.F `FlowAction/AgentSourceBizDetails`

Kelas `ASM-FW-GISFW-Data-Agent`
- `pyPreProcessingTransformRule` = `SearchHierarkiSourceBizAgent_PostDT`
- `pyActionTransformRule` = `SearchHierarkiSourceBizAgent_PostDT`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`

### 5.F `FlowAction/DeptHeadTreatyIn_UW`

Kelas `ASM-FW-GISFW-Work`
- `pySectionReference` = `GeneralDeptHeadTreatyIn_UW`
- `pyPreProcessingActivity` = `InputPolicyTreatyInPre_Act`
- `pyPreProcessingTransformRule` = `DeptHeadTreatyInUW_preDT`
- `pyActionTransformRule` = `DeptHeadTreatyIn_UW_postDT`
- `pyLocalActionActivity` = `InsertHistoryAkseptasiPega`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`
- menyertakan: `GeneralDeptHeadTreatyIn_UW`

### 5.F `FlowAction/InboxPolicyTreatyIn`

Kelas `ASM-FW-GISFW-Work`
- `pySectionReference` = `GeneralPolicyTreatyIn`
- `pyPreProcessingActivity` = `InputPolicyTreatyInPre_Act`
- `pyPreProcessingTransformRule` = `InputPolicyTreatyIn_preDT`
- `pyActionTransformRule` = `InboxPolicyTreatyIn_postDT`
- `pyLocalActionActivity` = `InputPolicyTreatyInPost_Act`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`
- menyertakan: `GeneralPolicyTreatyIn`

### 5.F `FlowAction/InputHistoricalSurveyReport`

Kelas `ASM-FW-GISFW-Data-Quotation`
- `pySectionReference` = `InputHistoricalSurveyReportDtl`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`
- menyertakan: `InputHistoricalSurveyReportDtl`

### 5.F `FlowAction/InputHistoricalSurveyReportUW`

Kelas `ASM-FW-GISFW-Data-Quotation`
- `pySectionReference` = `InputHistoricalSurveyReportDtlUW`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`
- menyertakan: `InputHistoricalSurveyReportDtlUW`

### 5.F `FlowAction/InstallmentList`

Kelas `ASM-FW-GISFW-Data-Installment`
- `pySectionReference` = `InstallmentList`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`
- menyertakan: `InstallmentList`

### 5.F `FlowAction/Installments_ReadOnly`

Kelas `ASM-FW-GISFW-Data-TreatyInInstallment`
- `pySectionReference` = `Installments_ReadOnly`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`
- menyertakan: `Installments_ReadOnly`

### 5.F `FlowAction/PolicyTreatyInDeclineConfirm`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- `pySectionReference` = `PolicyTreatyInDeclineConfirm`
- `pyLocalActionActivity` = `SetDueTo_act`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`
- menyertakan: `PolicyTreatyInDeclineConfirm`

### 5.F `FlowAction/ShowPolicyNoTreaty`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- `pySectionReference` = `ShowPolicyNoTreaty_SC`
- `pySubmitLabel` = `Submit`
- `pyconfirmchoice` = `ShowHarness`
- menyertakan: `ShowPolicyNoTreaty_SC`

### 5.S `Section/BusinessAndSOBList`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.pyTemplateInputBox` | pxButton | Choose |  |  |  | click→runActivity SetValue_Act(ID=.ID, SOB=, TREATYYEAR=, LIMITCURRENCY=, CEDING=, TREATYGROUP=, TREATYTYPE=, CLASSOFBUSINESSID=); click→runScript script:opener.location.reload; click→runScript script:window.close; click→refresh; -→runActivity SetValue_Act(ID=.ID, SOB=, TREATYYEAR=, LIMITCURRENCY=, CEDING=, TREATYGROUP=, TREATYTYPE=, CLASSOFBUSINESSID=) |
| `.TREATYID` | FIELD |  |  | true |  |  |
| `.TREATYCONTRACTNAME` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.CLASSOFBUSINESS` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.SOB` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.CEDING` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.PROPORTIONTYPE` | pxTextInput | Text Input |  | true |  |  |
| `.TREATYTYPE` | pxTextInput | Text Input |  | true |  |  |
| `.TREATYGROUP` | pxDisplayText |  |  | true |  |  |
| `.TREATYYEAR` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LIMITCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.LIMITVALUE` | pxCurrency |  |  | true |  |  |
| `.RETENTIONCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.RETENTIONVALUE` | pxCurrency |  |  | true |  |  |
| `.EPICURRENCY` | pxDisplayText |  |  | true |  |  |
| `.EPIVALUE` | pxCurrency |  |  | true |  |  |
| `.LAYERTYPE` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYER` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYERPARTTYPE` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYERPART` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.MDPCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.MDPVALUE` | pxCurrency |  |  | true |  |  |
| `.NETPREMICURRENCY` | pxDisplayText |  |  | true |  |  |
| `.NETPREMIVALUE` | pxCurrency |  |  | true |  |  |
| `.SHARECURRENCY` | pxDisplayText |  |  | true |  |  |
| `.SHAREVALUE` | pxCurrency |  |  | true |  |  |
| `For Treatyindetail join EDM` | LABEL |  |  |  | OTHER `1=2` |  |
| `.pyTemplateInputBox` | pxButton | Choose |  |  |  | click→runActivity SetValue_Act(ID=.ID, SOB=, TREATYYEAR=, LIMITCURRENCY=, CEDING=, TREATYGROUP=, TREATYTYPE=, CLASSOFBUSINESSID=); click→runScript script:opener.location.reload; click→runScript script:window.close; click→refresh; -→runActivity SetValue_Act(ID=.ID, SOB=, TREATYYEAR=, LIMITCURRENCY=, CEDING=, TREATYGROUP=, TREATYTYPE=, CLASSOFBUSINESSID=) |
| `.TREATYID` | FIELD |  |  | true |  |  |
| `.TREATYCONTRACTNAME` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.CLASSOFBUSINESS` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.SOB` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.CEDING` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.PROPORTIONTYPE` | pxTextInput | Text Input |  | true |  |  |
| `.TREATYTYPE` | pxTextInput | Text Input |  | true |  |  |
| `.TREATYGROUP` | pxDisplayText |  |  | true |  |  |
| `.TREATYYEAR` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LIMITCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.LIMITVALUE` | pxCurrency |  |  | true |  |  |
| `.RETENTIONCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.RETENTIONVALUE` | pxCurrency |  |  | true |  |  |
| `.EPICURRENCY` | pxDisplayText |  |  | true |  |  |
| `.EPIVALUE` | pxCurrency |  |  | true |  |  |
| `.LAYERTYPE` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYER` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYERPARTTYPE` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYERPART` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.MDPCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.MDPVALUE` | pxCurrency |  |  | true |  |  |
| `.NETPREMICURRENCY` | pxDisplayText |  |  | true |  |  |
| `.NETPREMIVALUE` | pxCurrency |  |  | true |  |  |
| `.SHARECURRENCY` | pxDisplayText |  |  | true |  |  |
| `.SHAREVALUE` | pxCurrency |  |  | true |  |  |

### 5.S `Section/BusinessAndSOBListRetro`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.pyTemplateInputBox` | pxButton | Choose |  |  |  | click→runActivity SetValue_Act(ID=.ID, SOB=, TREATYYEAR=, LIMITCURRENCY=, CEDING=, TREATYGROUP=, TREATYTYPE=, CLASSOFBUSINESSID=); click→runScript script:opener.location.reload; click→runScript script:window.close; click→refresh; -→runActivity SetValue_Act(ID=.ID, SOB=, TREATYYEAR=, LIMITCURRENCY=, CEDING=, TREATYGROUP=, TREATYTYPE=, CLASSOFBUSINESSID=) |
| `.TREATYID` | FIELD |  |  | true |  |  |
| `.TREATYCONTRACTNAME` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.CLASSOFBUSINESS` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.SOB` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.CEDING` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.PROPORTIONTYPE` | pxTextInput | Text Input |  | true |  |  |
| `.TREATYTYPE` | pxTextInput | Text Input |  | true |  |  |
| `.TREATYGROUP` | pxDisplayText |  |  | true |  |  |
| `.TREATYYEAR` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LIMITCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.LIMITVALUE` | pxCurrency |  |  | true |  |  |
| `.RETENTIONCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.RETENTIONVALUE` | pxCurrency |  |  | true |  |  |
| `.EPICURRENCY` | pxDisplayText |  |  | true |  |  |
| `.EPIVALUE` | pxCurrency |  |  | true |  |  |
| `.LAYERTYPE` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYER` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYERPARTTYPE` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYERPART` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.MDPCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.MDPVALUE` | pxCurrency |  |  | true |  |  |
| `.NETPREMICURRENCY` | pxDisplayText |  |  | true |  |  |
| `.NETPREMIVALUE` | pxCurrency |  |  | true |  |  |
| `.SHARECURRENCY` | pxDisplayText |  |  | true |  |  |
| `.SHAREVALUE` | pxCurrency |  |  | true |  |  |
| `For Treatyindetail join EDM` | LABEL |  |  |  | OTHER `1=2` |  |
| `.pyTemplateInputBox` | pxButton | Choose |  |  |  | click→runActivity SetValueRetro_Act(SOB=, CLASSOFBUSINESSID=, TREATYYEAR=, LIMITCURRENCY=, CEDING=, TREATYGROUP=, TREATYTYPE=, ID=.ID); click→runScript script:opener.location.reload; click→runScript script:window.close; click→refresh; -→runActivity SetValueRetro_Act(SOB=, CLASSOFBUSINESSID=, TREATYYEAR=, LIMITCURRENCY=, CEDING=, TREATYGROUP=, TREATYTYPE=, ID=.ID) |
| `.TREATYID` | FIELD |  |  | true |  |  |
| `.TREATYCONTRACTNAME` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.CLASSOFBUSINESS` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.SOB` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.CEDING` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.PROPORTIONTYPE` | pxTextInput | Text Input |  | true |  |  |
| `.TREATYTYPE` | pxTextInput | Text Input |  | true |  |  |
| `.TREATYGROUP` | pxDisplayText |  |  | true |  |  |
| `.TREATYYEAR` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LIMITCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.LIMITVALUE` | pxCurrency |  |  | true |  |  |
| `.RETENTIONCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.RETENTIONVALUE` | pxCurrency |  |  | true |  |  |
| `.EPICURRENCY` | pxDisplayText |  |  | true |  |  |
| `.EPIVALUE` | pxCurrency |  |  | true |  |  |
| `.LAYERTYPE` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYER` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYERPARTTYPE` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.LAYERPART` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.MDPCURRENCY` | pxDisplayText |  |  | true |  |  |
| `.MDPVALUE` | pxCurrency |  |  | true |  |  |
| `.NETPREMICURRENCY` | pxDisplayText |  |  | true |  |  |
| `.NETPREMIVALUE` | pxCurrency |  |  | true |  |  |
| `.SHARECURRENCY` | pxDisplayText |  |  | true |  |  |
| `.SHAREVALUE` | pxCurrency |  |  | true |  |  |

### 5.S `Section/DetailDeptHeadTreatyIn_UW`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `DetailPoliciesNonProportional`, `DetailPolicyTreatyOutNonProportional`, `ListSuggest`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.NoOffer` | pxDisplayText | Master ID |  | true |  |  |
| `pyWorkPage.TreatyIn.Commencement` | pxDateTime | Commencement |  | true |  |  |
| `.StartDate` | pxDateTime | Statement Period | true | true `1==1` |  | change→refresh; -→refresh |
| `.SOBName` | pxDisplayText | Source Of Business |  | true | selalu (sisa syarat tak berlaku: `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`) |  |
| `.TreatyGroupName` | pxDisplayText | Treaty Group |  | true | OTHER `.IsNewPolicyNonProp != 1` |  |
| `.BizName` | pxDisplayText | Class Of Business |  | true | OTHER `.IsNewPolicyNonProp != 1 && NEVER` |  |
| `pyWorkPage.PolicyTreatyIn.PolicyNo` | pxTextInput | No Polis |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.ProductionDate` | pxDateTime | Production Date |  | true |  |  |
| `.DueTo` | pxRadioButtons | Due To Us / You |  | nonaktif `` | selalu (sisa syarat tak berlaku: `1=2`) |  |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | Survey Report |  | true `1=1` | OTHER `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` | change→postValue; change→refresh |
| `.StatementType` | pxDropdown | Statement Type |  | true nonaktif `` | selalu (sisa syarat tak berlaku: `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`) | change→postValue |
| `.pyTemplateButton` | pxButton | Survey Report |  | nonaktif `.QuotationData.IsSurveyReport=='No' \\|\\| .QuotationData.IsSurveyReport==''` | OTHER `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` | click→showHarness HistoricalSurveyReportUW |
| `.FlagPPH` | pxCheckbox | Checkbox |  | nonaktif `` |  | click→postValue; -→postValue |
| `.TypeTax` | pxDisplayText | Type Tax |  | true | NOTBLANK | click→postValue |
| `.QuotationData.NoOfferSlip` | pxTextArea | No Offer Slip |  | nonaktif `` |  |  |
| `.ShareCurrency` | pxTextInput | RNM Share |  | true |  |  |
| `.ShareValue` | pxCurrency | ShareValue |  | true |  | change→refresh CountResult1Onp_Act |
| `.StatementDate` | pxDateTime | Statement Date | true | true `1==1` |  |  |
| `pyWorkPage.TreatyIn.Termination` | pxDateTime | Termination |  | true |  |  |
| `.EndDate` | pxDateTime | To | true | true `1==1` |  |  |
| `.CedingCoName` | pxDisplayText | Ceding Company |  | true | NOTBLANK |  |
| `.InsuredName` | pxTextInput | Formatted Text |  | true |  |  |
| `.TreatyType` | pxTextInput | Treaty Type |  | true |  |  |
| `.TreatyYear` | pxTextInput | UW Year |  | true | OTHER `.IsNewPolicyNonProp != 1` |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `.Quartal` | pxTextInput | .Quartal |  | true |  |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `/` | LABEL |  |  |  | selalu (sisa syarat tak berlaku: `1=2`) |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `.YearOfQuartal` | pxTextInput | Text Input |  | true |  |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `U/Y` | LABEL |  |  |  | selalu (sisa syarat tak berlaku: `1=2`) |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `.TreatyYear` | pxTextInput | TreatyYear |  | true |  |  |
| `.MarketingOfficer` | pxDisplayText | Marketing Officer |  | true | NOTBLANK |  |
| `.QuotationData.ProportionalType` | pxDropdown ⟵ `BrowseCurrencyTreatyIn_RD` | Proportional Type |  | true | OTHER `.IsNewPolicyNonProp != 1` | change→refresh SetCurrency_act(CURR=.IDCurrency) |
| `.Currency` | pxAutoComplete ⟵ `BrowseCurrencyTreatyIn_RD` | Text Input |  | true | OTHER `.IsNewPolicyNonProp != 1` |  |
| `.ClaimType` | pxDropdown ⟵ `BrowseCurrencyTreatyIn_RD` | Claim Type |  | true |  | change→refresh SetCurrency_act(CURR=.IDCurrency) |
| `.ClaimPaymentType` | pxDropdown ⟵ `BrowseCurrencyTreatyIn_RD` | Claim Payment Type |  | true |  | change→refresh SetCurrency_act(CURR=.IDCurrency) |
| `.LayerType` | pxTextInput | LayerType |  | true |  |  |
| `spc` | LABEL |  |  |  | OTHER `1=2` |  |
| `.Layer` | pxTextInput | Layer |  | true |  |  |
| `space` | LABEL |  |  |  | OTHER `1=2` |  |
| `Of` | LABEL |  |  |  | selalu (sisa syarat tak berlaku: `1=2`) |  |
| `space` | LABEL |  |  |  | OTHER `1=2` |  |
| `.LayerPartType` | pxTextInput | LayerPartType |  | true |  |  |
| `spc` | LABEL |  |  |  | OTHER `1=2` |  |
| `.LayerPart` | pxTextInput | LayerPart |  | true |  |  |
| `.Remark` | pxTextArea | Text Area |  | true |  |  |
| `Section for Non Proportional` | LABEL |  |  |  | OTHER `NEVER` |  |
| `DetailPoliciesNonProportional` | SUB_SECTION ⟵ `DetailPoliciesNonProportional` | DetailPoliciesNonProportional |  |  |  |  |
| `Section for when Non Proportional Retro` | LABEL |  |  |  | OTHER `NEVER` |  |
| `DetailPolicyTreatyOutNonProportional` | SUB_SECTION ⟵ `DetailPolicyTreatyOutNonProportional` | DetailPolicyTreatyOutNonProportional |  |  |  |  |
| `Section for Others Than Non Proportional` | LABEL |  |  |  | OTHER `NEVER` |  |
| `.GrossPremium` | pxNumber | Gross Premium 100% |  | true |  |  |
| `.PremiOgp` | pxCurrency | Premi Ogp | true | true `1==1` |  | change→refresh CountResult1_Act; -→refresh CountResult1_Act |
| `.RiCommOgp` | pxNumber | Deduction In A (OGP) | true | true `1==1` |  | change→refresh CountResult1_Act |
| `.ResultOgp1` | pxCurrency | ResultOgp1 |  | true `1 = 1` | selalu (sisa syarat tak berlaku: `1==2`) | change→refresh CountRiCommOgp_act(DiscountType="Amount") |
| `.OveriddingCommOgp` | pxNumber | Deduction In B (OGP) | true | true `1==1` |  | change→refresh CountResult2Ogp_act |
| `.ResultOgp2` | pxCurrency | ResultOgp2 |  | true |  | change→refresh CountOverridingCommOgp_Act(Overidding="Amount") |
| `.Claim` | pxCurrency | Claim | true | true `1==1` |  | change→refresh CountNetPremi_act; -→refresh CountNetPremi_act |
| `.OutstandingClaim` | pxCurrency | Outstanding Claim | true | true `1==1` |  |  |
| `.SalvageValue` | pxCurrency | Salvage | true | true `1==1` |  | change→refresh CountNetPremi_act |
| `.ExcessLoss` | pxCurrency | Excess Loss | true | true `1==1` |  | change→refresh CountNetPremi_act; -→refresh CountNetPremi_act |
| `.NetPremium` | pxCurrency | Net Premium |  | true |  |  |
| `.BalanceDueTo` | pxTextInput | Balance Due To You |  | true | OTHER `.BalanceDueTo < 0` |  |
| `.BalanceBeforeTax` | pxTextInput | Balance Before Tax |  | true | OTHER `.BalanceDueTo >= 0` |  |
| `.BalanceBeforePPH` | pxTextInput | Balance Before Withholding Tax (PPH 2.2) |  | true | OTHER `.BalanceDueTo >= 0` |  |
| `.BalanceDueTo` | pxTextInput | Balance Due To Us |  | true | OTHER `.BalanceDueTo >=0` |  |
| `.PremiOnp` | pxCurrency | Premi Onp | true | true `1==1` |  | change→refresh CountResult1Onp_Act |
| `.RiCommOnp` | pxNumber | Deduction In A (ONP) | true | true `1==1` |  | change→refresh CountResult1Onp_Act |
| `.ResultOnp1` | pxCurrency | ResultOnp1 | true | true `1==1` |  | change→refresh CountRiCommOnp_act(Result="Amount") |
| `.OveriddingCommOnp` | pxNumber | Deduction In B (ONP) | true | true `1==1` |  | change→refresh CountResult2Onp_act |
| `.ResultOnp2` | pxCurrency | ResultOnp2 |  | true |  | change→refresh CountOverridingCommOnp_Act(overidding="Amount") |
| `.Deduction1` | pxCurrency | Deduction1 | true | true `1==1` |  | change→refresh CountNetPremi_act; -→refresh CountNetPremi_act |
| `.Deduction2` | pxCurrency | Deduction2 | true | true `1==1` |  | change→refresh CountNetPremi_act; -→refresh CountNetPremi_act |
| `.PPHValue` | pxCurrency | PPH 2% |  | true |  |  |
| `.PPNValue` | pxCurrency | PPN 2.2% |  | true |  |  |
| `.TreatyType` | pxDropdown ⟵ `BrowseReinsuranceType_RD` |  |  | true |  |  |
| `.SharePercentage` | pxNumber |  |  | true |  | change→refresh CountSpreading_Act(SharePct=.SharePercentage, Index=.pxListSubscript, Type=Premium); -→refresh CountSpreading_Act(SharePct=.SharePercentage, Index=.pxListSubscript, Type=Premium) |
| `.PremiumSpreaded` | pxNumber |  |  | true |  |  |
| `.ClaimPercentage` | pxNumber |  |  | true |  | change→refresh CountSpreading_Act(SharePct=.ClaimPercentage, Index=.pxListSubscript, Type=Claim); -→refresh CountSpreading_Act(SharePct=.ClaimPercentage, Index=.pxListSubscript, Type=Claim) |
| `.ClaimSpreaded` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber |  |  | true |  |  |
| `.TreatyType` | pxDropdown ⟵ `BrowseReinsuranceType_RD` |  |  | true |  |  |
| `.SharePercentage` | pxNumber |  |  | true |  | change→refresh CountSpreading_Act(SharePct=.SharePercentage, Index=.pxListSubscript, Type=Premium); -→refresh CountSpreading_Act(SharePct=.SharePercentage, Index=.pxListSubscript, Type=Premium) |
| `.PremiumSpreaded` | pxNumber |  |  | true |  |  |
| `.ClaimSpreaded` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber |  |  | true |  |  |
| `.TotalSharePercentagePremium` | pxNumber | Total %Share |  | true |  |  |
| `.TotalPremium` | pxNumber | Total Premium |  | true |  |  |
| `.TotalSharePercentageClaim` | pxNumber | Total %Share Claim |  | true |  |  |
| `.TotalClaim` | pxNumber | Total Claim |  | true |  |  |
| `.Installment` | pxTextInput | Text Input |  | true |  | change→refresh FillPaymentInstallment(Installment=.Installment); -→refresh FillPaymentInstallment(Installment=.Installment) |
| `.InstallmentNo` | pxInteger |  |  | true |  |  |
| `.DueDate` | pxDateTime |  |  | true |  |  |
| `.InstallmentPercentage` | pxNumber |  |  | true nonaktif `IsClaim` |  | change→refresh SetValidateInstallment_Act; -→refresh SetValidateInstallment_Act |
| `.Premium` | pxNumber |  |  | true |  | change→refresh CountPctInstallment_Act(idx=.InstallmentNo); -→refresh CountPctInstallment_Act(idx=.InstallmentNo) |
| `.PaymentTotal` | pxNumber |  |  | true |  |  |
| `.pyTemplateButton` | pxButton | Submit |  |  | OTHER `.IsApproved = '0'` | click→finishAssignment |
| `.pyTemplateButton` | pxButton | Submit |  |  | OTHER `.IsApproved == 1 && OperatorID.pyUserIdentifier!='<ID-operator-1>' && pyWorkPage.LetterNo=='TREATYINDEPTHEAD'` | click→finishAssignment; -→finishAssignment |
| `.pyTemplateButton` | pxButton | Submit |  |  | OTHER `.IsApproved == 1 && OperatorID.pyUserIdentifier!='<ID-operator-1>' && pyWorkPage.LetterNo==''` | click→runActivity GeneratePolicyNoTreaty_Act; click→localAction ShowPolicyNoTreaty; click→refresh GeneralDeptHeadTreatyIn_UW |
| `.pyTemplateButton` | pxButton | Submit |  |  | OTHER `.IsApproved == 1 && OperatorID.pyUserIdentifier=='<ID-operator-1>'` | click→runActivity GeneratePolicyNoTreaty_Act; click→localAction ShowPolicyNoTreaty; click→refresh GeneralDeptHeadTreatyIn_UW |

### 5.S `Section/DetailPoliciesNonProportional`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `DetailPolicyTreatyInNonProportional`, `DetailPolicyTreatyInNonProportionalEDM`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `NON EDM` | LABEL |  |  |  | OTHER `OperatorID.pxInsName = '<ID-operator-2>'` |  |
| `DetailPolicyTreatyInNonProportional` | SUB_SECTION ⟵ `DetailPolicyTreatyInNonProportional` |  |  |  |  |  |
| `EDM` | LABEL |  |  |  | OTHER `OperatorID.pxInsName = '<ID-operator-2>'` |  |
| `DetailPolicyTreatyInNonProportionalEDM` | SUB_SECTION ⟵ `DetailPolicyTreatyInNonProportionalEDM` | DetailPolicyTreatyInNonProportionalEDM |  |  |  |  |

### 5.S `Section/DetailPolicyTreatyIn`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `DetailPoliciesNonProportional`, `DetailPolicyTreatyOutNonProportional`, `ListSuggest`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.NoOffer` | pxDisplayText | Master ID |  | true |  |  |
| `pyWorkPage.TreatyIn.Commencement` | pxDateTime | Commencement |  | true |  |  |
| `.StartDate` | pxDateTime | Statement Period | true |  |  | change→refresh |
| `pyWorkPage.PolicyTreatyIn.SOBName` | pxTextInput | Source Of Business |  | true |  |  |
| `.pyTemplateInputBox` | pxButton | Select Source Of Business |  |  | OTHER `.ClaimType = 'XOL Retro'` | click→showHarness InputQuotation_PreAct(Acton=SOB) |
| `.TreatyGroupName` | pxTextInput | Treaty Group |  | true | OTHER `.IsNewPolicyNonProp != 1` |  |
| `.BizName` | pxTextInput | Class Of Business |  | true | OTHER `1=2` |  |
| `.DueTo` | pxRadioButtons | Due To Us / You |  |  | OTHER `NEVER` |  |
| `.QuotationData.IsSurveyReport` | pxRadioButtons | Survey Report | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` | `IsUW` | OTHER `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` | change→postValue; change→refresh |
| `.StatementType` | pxDropdown | Statement Type |  | `IsUW` | selalu (sisa syarat tak berlaku: `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`) | change→postValue |
| `.QuotationData.NoOfferSlip` | pxTextArea | No Offer Slip |  |  |  | change→postValue; -→postValue |
| `.pyTemplateButton` | pxButton | Survey Report |  | nonaktif `.QuotationData.IsSurveyReport=='No' \\|\\| .QuotationData.IsSurveyReport==''` | OTHER `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` | click→showHarness HistoricalSurveyReport |
| `.FlagRetroTreaty` | pxCheckbox | Checkbox |  |  | OTHER `.ClaimType != 'XOL Retro'` | change→postValue; click→postValue |
| `.FlagPPH` | pxCheckbox | Checkbox |  |  |  | change→postValue; change→runActivity RemoveTypeTax_ACT; click→postValue; click→runActivity RemoveTypeTax_ACT |
| `.TypeTax` | pxRadioButtons | Type Tax | true `.FlagPPH = true` |  | OTHER `.FlagPPH = true` | change→postValue |
| `.pyTemplateInputBox` | pxButton | Choose Business |  |  | OTHER `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy'` | click→showHarness BusinessAndSOBList; click→refresh |
| `.pyTemplateInputBox` | pxButton | Choose Business R |  |  | OTHER `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy'` | click→showHarness BusinessAndSOBListRetro; click→refresh |
| `.ShareCurrency` | pxTextInput | RNM Share |  | true |  |  |
| `.ShareValue` | pxCurrency | ShareValue |  | true |  | change→refresh CountResult1Onp_Act |
| `.StatementDate` | pxDateTime | Statement Date | true | true `ALWAYS` nonaktif `` |  |  |
| `pyWorkPage.TreatyIn.Termination` | pxDateTime | Termination |  | true |  |  |
| `.EndDate` | pxDateTime | To | true |  |  | change→refresh ProtectDate; -→refresh ProtectDate |
| `.CedingCoName` | pxTextInput | Ceding Company |  | true |  |  |
| `.InsuredName` | pxTextInput | Formatted Text |  | true |  |  |
| `.TreatyType` | pxTextInput | Treaty Type |  | true |  |  |
| `.TreatyYear` | pxTextInput | UW Year |  | true | OTHER `.IsNewPolicyNonProp != 1` |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `.Quartal` | pxTextInput | .Quartal | true `.QuotationData.ProportionalType = 'Proportional'` |  |  |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `/` | LABEL |  |  |  | selalu (sisa syarat tak berlaku: `1=2`) |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `.YearOfQuartal` | pxTextInput | Text Input | true `.QuotationData.ProportionalType = 'Proportional'` |  |  |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `U/Y` | LABEL |  |  |  | selalu (sisa syarat tak berlaku: `1=2`) |  |
| `s` | LABEL |  |  |  | OTHER `1=2` |  |
| `.TreatyYear` | pxTextInput | TreatyYear |  | true |  |  |
| `.QuotationData.MOID` | pxDropdown ⟵ `BrowseMarketingOfficer_RD` | Marketing Officer | true |  |  | change→postValue; change→refresh CheckDataMkt; -→postValue |
| `.QuotationData.ProportionalType` | pxDropdown ⟵ `BrowseCurrencyTreatyIn_RD` | Proportional Type |  | true | OTHER `.IsNewPolicyNonProp != 1` | change→refresh; -→refresh |
| `.IDCurrency` | pxDropdown ⟵ `BrowseCurrencyTreatyIn_RD` | Currency |  |  | OTHER `.IsNewPolicyNonProp != 1` | change→refresh SetCurrency_act(CURR=.IDCurrency) |
| `.ClaimType` | pxDropdown ⟵ `BrowseCurrencyTreatyIn_RD` | Claim Type | true `.Claim != '' && .Claim != 0` |  |  |  |
| `.ClaimPaymentType` | pxDropdown ⟵ `BrowseCurrencyTreatyIn_RD` | Payment Type | true `.Claim != '' && .Claim != 0` |  |  |  |
| `.pyTemplateButton` | pxButton | Enable / Disable Input Type |  |  | OTHER `.TreatyType='XOL'` | click→refresh; -→refresh |
| `.LayerType` | pxTextInput | LayerType |  | true |  |  |
| `spc` | LABEL |  |  |  | OTHER `1=2` |  |
| `.Layer` | pxTextInput | Layer |  | true |  |  |
| `space` | LABEL |  |  |  | OTHER `1=2` |  |
| `Of` | LABEL |  |  |  | selalu (sisa syarat tak berlaku: `1=2`) |  |
| `space` | LABEL |  |  |  | OTHER `1=2` |  |
| `.LayerPartType` | pxTextInput | LayerPartType |  | true |  |  |
| `spc` | LABEL |  |  |  | OTHER `1=2` |  |
| `.LayerPart` | pxTextInput | LayerPart |  | true |  |  |
| `.Remark` | pxTextArea | Text Area |  |  |  |  |
| `Section for when Non Proportional` | LABEL |  |  |  | OTHER `NEVER` |  |
| `DetailPoliciesNonProportional` | SUB_SECTION ⟵ `DetailPoliciesNonProportional` | DetailPoliciesNonProportional |  |  |  |  |
| `Section for when Non Proportional Retro` | LABEL |  |  |  | OTHER `NEVER` |  |
| `DetailPolicyTreatyOutNonProportional` | SUB_SECTION ⟵ `DetailPolicyTreatyOutNonProportional` | DetailPolicyTreatyOutNonProportional |  |  |  |  |
| `Section for Old Soa Input Format (individual input)` | LABEL |  |  |  | OTHER `NEVER` |  |
| `.GrossPremium` | pxNumber | Gross Premium 100% |  |  |  | change→refresh CalculatePremi_Act(Action="PREMIUM") |
| `.GrossClaim` | pxNumber | Claim 100% |  |  |  | change→refresh CalculatePremi_Act(Action="CLAIM"); -→refresh CalculatePremi_Act(Action="CLAIM") |
| `.PremiOgp` | pxCurrency | Premi Ogp | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountOGPONP_Act; keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.RiCommOgp` | pxNumber | (%) Deduction In A (OGP) | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountResult1_Act(Data="Pct"); change→refresh CountOGPONP_Act; keyboard→refresh CountResult1_Act(Data="Pct"); keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.ResultOgp1` | pxCurrency | Deduction In A (OGP) |  |  |  | change→refresh CountResult1_Act(Data="Amount"); change→refresh CountOGPONP_Act; keyboard→refresh CountResult1_Act(Data="Amount"); keyboard→refresh CountOGPONP_Act; -→refresh CountResult1_Act(Data="Amount"); -→refresh CountOGPONP_Act |
| `.OveriddingCommOgp` | pxNumber | (%) Deduction In B (OGP) | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountResult2Ogp_act(Data="Pct"); change→refresh CountOGPONP_Act; keyboard→refresh CountResult2Ogp_act(Data="Pct"); keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.ResultOgp2` | pxCurrency | Deduction In B (OGP) |  |  |  | change→refresh CountResult2Ogp_act(Data="Amount"); change→refresh CountOGPONP_Act; keyboard→refresh CountResult2Ogp_act(Data="Amount"); keyboard→refresh CountOGPONP_Act; -→refresh CountResult2Ogp_act(Data="Amount"); -→refresh CountOGPONP_Act |
| `.PremiOnp` | pxCurrency | Premi Onp | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountOGPONP_Act; keyboard→refresh CountOGPONP_Act |
| `.RiCommOnp` | pxNumber | (%) Deduction In A (ONP) | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountResult1Onp_Act(Data="Pct"); change→refresh CountOGPONP_Act; keyboard→refresh CountResult1Onp_Act(Data="Pct"); keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.ResultOnp1` | pxCurrency | Deduction In A (ONP) |  |  |  | change→refresh CountResult1Onp_Act(Data="Amount"); change→refresh CountOGPONP_Act; keyboard→refresh CountResult1Onp_Act(Data="Amount"); keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.OveriddingCommOnp` | pxNumber | (%) Deduction In B (ONP) | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountResult2Onp_act(Data="Pct"); change→refresh CountOGPONP_Act; keyboard→refresh CountResult2Onp_act(Data="Pct"); keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.ResultOnp2` | pxCurrency | Deduction In B (ONP) |  |  |  | change→refresh CountResult2Onp_act(Data="Amount"); change→refresh CountOGPONP_Act; keyboard→refresh CountResult2Onp_act(Data="Amount"); keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.Claim` | pxCurrency | Claim | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountOGPONP_Act; keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.OutstandingClaim` | pxCurrency | Outstanding Claim | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  |  |
| `.SalvageValue` | pxCurrency | Salvage | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountOGPONP_Act; keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.ExcessLoss` | pxCurrency | Excess Loss | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountOGPONP_Act; keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.NetPremium` | pxCurrency | Total Premium Before Claim |  | true |  |  |
| `.BalanceDueTo` | pxTextInput | Balance Due To You |  | true | OTHER `.BalanceDueTo < 0` |  |
| `.BalanceBeforeTax` | pxTextInput | Balance Before Tax |  | true | OTHER `.BalanceDueTo >= 0` |  |
| `.BalanceBeforePPH` | pxTextInput | Balance Before Withholding Tax (PPH 2.2) |  | true | OTHER `.BalanceDueTo >= 0` |  |
| `.BalanceDueTo` | pxTextInput | Balance Due To Us |  | true | OTHER `.BalanceDueTo >= 0` |  |
| `.Deduction1` | pxCurrency | Deduction1 | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountOGPONP_Act; keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.Deduction2` | pxCurrency | Deduction2 | true `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |  |  | change→refresh CountOGPONP_Act; keyboard→refresh CountOGPONP_Act; -→refresh CountOGPONP_Act |
| `.pyTemplateButton` | pxButton | Save |  |  | OTHER `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy'` | click→save |
| `.PPHValue` | pxCurrency | PPH 2% |  | true |  |  |
| `.PPNValue` | pxCurrency | PPN 2.2% |  | true |  |  |
| `.pyTemplateButton` | pxButton | Add |  |  | OTHER `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy'` | click→addRow; -→addRow |
| `.TreatyType` | pxDropdown ⟵ `ListSpreading.pxResults` |  |  |  |  |  |
| `.SharePercentage` | pxNumber |  |  |  |  | change→refresh CountSpreading_Act(Index=.pxListSubscript); -→refresh CountSpreading_Act(Index=.pxListSubscript) |
| `.PremiumSpreaded` | pxNumber |  |  | true |  |  |
| `.ClaimPercentage` | pxNumber |  |  |  |  | change→refresh CountSpreading_Act(Index=.pxListSubscript); -→refresh CountSpreading_Act(Index=.pxListSubscript) |
| `.ClaimSpreaded` | pxNumber |  |  | true |  |  |
| `.pyTemplateButton` | pxButton | Delete |  |  | OTHER `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy'` | click→deleteRow; -→deleteRow |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber |  |  | true |  |  |
| `.TreatyType` | pxDropdown ⟵ `BrowseReinsuranceType_RD` |  |  | true |  |  |
| `.SharePercentage` | pxNumber |  |  | true |  | change→refresh CountSpreading_Act(SharePct=.SharePercentage, Index=.pxListSubscript, Type=Premium); -→refresh CountSpreading_Act(SharePct=.SharePercentage, Index=.pxListSubscript, Type=Premium) |
| `.PremiumSpreaded` | pxNumber |  |  | true |  |  |
| `.ClaimSpreaded` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber |  |  | true |  |  |
| `.TotalSharePercentagePremium` | pxNumber | Total %Share |  | true |  |  |
| `.TotalPremium` | pxNumber | Total Premium |  | true |  |  |
| `.TotalSharePercentageClaim` | pxNumber | Total %Share Claim |  | true |  |  |
| `.TotalClaim` | pxNumber | Total Claim |  | true |  |  |
| `.Installment` | pxTextInput | Text Input |  |  |  | change→refresh FillPaymentInstallment(Installment=.Installment); -→refresh FillPaymentInstallment(Installment=.Installment) |
| `.InstallmentNo` | pxInteger |  |  | true |  |  |
| `.DueDate` | pxDateTime |  |  | true `IsUW \\|\\| IsClaim \\|\\| pyPortal.IsShowOpenPolicyMarine == true \\|\\| pyWorkPage.IsOldData = 1` |  |  |
| `.InstallmentPercentage` | pxNumber |  |  | true `IsUW \\|\\| pyPortal.IsShowOpenPolicyMarine == true \\|\\| pyWorkPage.IsOldData = 1` nonaktif `IsClaim` |  | change→refresh SetValidateInstallment_Act; -→refresh SetValidateInstallment_Act |
| `.Premium` | pxNumber |  |  | true `IsUW \\|\\| pyPortal.IsShowOpenPolicyMarine == true \\|\\| pyWorkPage.IsOldData = 1` |  | change→refresh CountPctInstallment_Act(idx=.InstallmentNo); -→refresh CountPctInstallment_Act(idx=.InstallmentNo) |
| `.PaymentTotal` | pxNumber |  |  | true |  |  |
| `ListSuggest` | SUB_SECTION ⟵ `ListSuggest` |  |  |  |  |  |
| `.pyTemplateButton` | pxButton | Save |  |  |  | click→save; -→save |
| `.pyTemplateButton` | pxButton | Submit |  | nonaktif `TempEmail.CARI28 != pyWorkPage.PositionNote` | OTHER `.IsApproved==1` | click→runActivity SetDueTo_act; click→finishAssignment; -→runActivity SetDueTo_act; -→finishAssignment |
| `.pyTemplateButton` | pxButton | Submit |  | nonaktif `TempEmail.CARI28 != pyWorkPage.PositionNote` | OTHER `.IsApproved==0` | click→localAction PolicyTreatyInDeclineConfirm; -→localAction PolicyTreatyInDeclineConfirm |

### 5.S `Section/DetailPolicyTreatyInNonProportional`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `SpreadingRiskList`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.Note` | FIELD |  |  |  |  |  |
| `.Limit` | pxNumber |  |  | true |  |  |
| `.Limit2` | pxNumber |  |  | true |  |  |
| `.Deductible` | pxNumber |  |  | true |  |  |
| `.Deductible2` | pxNumber |  |  | true |  |  |
| `.MDP` | pxNumber |  |  | true |  |  |
| `.MDP2` | pxNumber |  |  | true |  |  |
|  | FIELD |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `pyWorkPage.TreatyIn.RNMShare` | pxTextInput | % RNM Share |  | true | OTHER `pyWorkPage.TreatyIn.FacultativeShare = 0` |  |
| `pyWorkPage.TreatyIn.RnmShareDeducted` | pxTextInput | % RNM Share |  | true nonaktif `TreatyIn.ViewState = 1` | OTHER `pyWorkPage.TreatyIn.FacultativeShare != 0` |  |
| `.Note` | FIELD |  |  |  |  |  |
| `.Limit` | pxNumber |  |  | true |  |  |
| `.Limit2` | pxNumber |  |  | true |  |  |
| `.MDP` | pxNumber |  |  | true |  |  |
| `.MDP2` | pxNumber |  |  | true |  |  |
| `.Deductible` | pxNumber |  |  | true |  |  |
| `.Deductible2` | pxNumber |  |  | true |  |  |
| `.NetPremi` | pxNumber |  |  | true |  |  |
| `.NetPremiAfterPPN` | pxNumber |  |  |  |  |  |
| `.NetPremiAfterPPH` | pxNumber |  |  |  |  |  |
| `.NetPremi2` | pxNumber |  |  | true |  |  |
| `.NetPremiAfterPPN2` | pxNumber |  |  |  |  |  |
| `.NetPremiAfterPPH2` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.TotalPPNValue` | pxNumber |  |  |  |  |  |
| `.TotalPPHValue` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.TotalNetPremiAfterPPN` | pxNumber |  |  |  |  |  |
| `.TotalNetPremiAfterTax` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.pyTemplateButton` | pxButton | Button |  |  | OTHER `1=2` | click→refresh TreatyInNonSetTotal |
| `pyWorkPage.TreatyIn.FacultativeShare` | pxTextInput | % Share Facultative |  | true nonaktif `TreatyIn.ViewState = 1` |  |  |
| `.Note` | FIELD |  |  |  |  |  |
| `.Limit` | FIELD |  |  |  |  |  |
| `.Limit2` | FIELD |  |  |  |  |  |
| `.MDP` | FIELD |  |  |  |  |  |
| `.MDP2` | FIELD |  |  |  |  |  |
| `.Deductible` | FIELD |  |  |  |  |  |
| `.Deductible2` | FIELD |  |  |  |  |  |
| `.NetPremi` | FIELD |  |  |  |  |  |
| `.NetPremi2` | FIELD |  |  |  |  |  |
|  | FIELD |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `SpreadingRiskList` | SUB_SECTION ⟵ `SpreadingRiskList` |  |  |  |  |  |
| `pyWorkPage.PolicyTreatyIn.Installment` | pxTextInput | Installment |  | true |  |  |
| `.Currency` | FIELD |  |  |  |  |  |

### 5.S `Section/DetailPolicyTreatyInNonProportionalEDM`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.Note` | FIELD |  |  |  |  |  |
| `.MDP` | FIELD |  |  |  |  |  |
| `.MDP2` | FIELD |  |  |  |  |  |
| `.Deductible` | FIELD |  |  |  |  |  |
| `.Deductible2` | FIELD |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `pyWorkPage.TreatyIn.RNMShare` | pxTextInput | % RNM Share |  | true nonaktif `TreatyIn.ViewState = 1` |  |  |
| `.Note` | FIELD |  |  |  |  |  |
| `.MDP` | FIELD |  |  |  |  |  |
| `.MDP2` | FIELD |  |  |  |  |  |
| `.Deductible` | FIELD |  |  |  |  |  |
| `.Deductible2` | FIELD |  |  |  |  |  |
| `.NetPremi` | FIELD |  |  |  |  |  |
| `.NetPremi2` | FIELD |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.pyTemplateButton` | pxButton | Button |  |  | OTHER `1=2` | click→refresh TreatyInNonSetTotal |
| `pyWorkPage.TreatyIn.InstallmentNo` | pxTextInput | Installment |  | true |  |  |
| `.Currency` | FIELD |  |  |  |  |  |

### 5.S `Section/DetailPolicyTreatyOutNonProportional`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`
- menyertakan: `SpreadingRiskList`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.Note` | FIELD |  |  |  |  |  |
| `.Limit` | pxNumber |  |  | true |  |  |
| `.Limit2` | pxNumber |  |  | true |  |  |
| `.Deductible` | pxNumber |  |  | true |  |  |
| `.Deductible2` | pxNumber |  |  | true |  |  |
| `.MDP` | pxNumber |  |  | true |  |  |
| `.MDP2` | pxNumber |  |  | true |  |  |
|  | FIELD |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `pyWorkPage.TreatyIn.ReinsurerShare` | pxTextInput | % Reinsurer Share |  | true | selalu (sisa syarat tak berlaku: `pyWorkPage.TreatyIn.FacultativeShare = 0`) |  |
| `pyWorkPage.TreatyIn.RnmShareDeducted` | pxTextInput | % RNM Share |  | true nonaktif `TreatyIn.ViewState = 1` | OTHER `pyWorkPage.TreatyIn.FacultativeShare != 0` |  |
| `.Note` | FIELD |  |  |  |  |  |
| `.Limit` | pxNumber |  |  | true |  |  |
| `.Limit2` | pxNumber |  |  | true |  |  |
| `.MDP` | pxNumber |  |  | true |  |  |
| `.MDP2` | pxNumber |  |  | true |  |  |
| `.Deductible` | pxNumber |  |  | true |  |  |
| `.Deductible2` | pxNumber |  |  | true |  |  |
| `.NetPremi` | pxNumber |  |  | true |  |  |
| `.NetPremi2` | pxNumber |  |  | true |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `.Currency` | pxNumber |  |  |  |  |  |
| `.Value` | pxNumber |  |  |  |  |  |
| `SpreadingRiskList` | SUB_SECTION ⟵ `SpreadingRiskList` |  |  |  |  |  |
| `pyWorkPage.PolicyTreatyIn.Installment` | pxTextInput | Installment |  | true |  |  |
| `.Currency` | FIELD |  |  |  |  |  |

### 5.S `Section/GeneralDeptHeadTreatyIn_UW`

Kelas `ASM-FW-GISFW-Work` · halaman `.PolicyTreatyIn`
- menyertakan: `DetailDeptHeadTreatyIn_UW`, `DetailPoliciesNonProportional`, `DetailPolicyTreatyOutNonProportional`, `ListSuggest`

### 5.S `Section/GeneralPolicyTreatyIn`

Kelas `ASM-FW-GISFW-Work` · halaman `.PolicyTreatyIn`
- menyertakan: `DetailPoliciesNonProportional`, `DetailPolicyTreatyIn`, `DetailPolicyTreatyOutNonProportional`, `ListSuggest`

### 5.S `Section/HistoricalSurveyReportDtl`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.InsuredName` | pxTextInput | Text Input |  | true |  |  |
| `.pyTemplateButton` | pxButton | Add |  |  |  | click→addRow; -→addRow |
| `.DateofSurvey` | pxDateTime | Text Input |  |  |  | change→refresh ConcatSlipOfferNo_Act; -→refresh ConcatSlipOfferNo_Act |
| `.SurveyedBy` | FIELD |  |  |  |  |  |
| `.LossPrevention` | pxNumber |  |  |  |  |  |
| `.Remarks` | pxDropdown |  |  |  |  |  |
| `.pyTemplateButton` | pxButton | Delete |  |  |  | click→deleteRow; -→deleteRow |
| `.pyTemplateButton` | pxButton | Submit |  |  |  | click→runActivity SetSurveyReport_Act; click→runScript script:window.close; -→runActivity SetSurveyReport_Act; -→runScript script:window.close |

### 5.S `Section/HistoricalSurveyReportDtlUW`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.InsuredName` | pxTextInput | Text Input |  | true |  |  |
| `.DateofSurvey` | pxDateTime | Text Input |  |  |  | change→refresh ConcatSlipOfferNo_Act; -→refresh ConcatSlipOfferNo_Act |
| `.SurveyedBy` | FIELD |  |  |  |  |  |
| `.Remarks` | pxDropdown |  |  |  |  |  |

### 5.S `Section/InputHistoricalSurveyReportDtl`

Kelas `ASM-FW-GISFW-Data-Quotation`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.DateofSurvey` | pxDateTime | Date of Survey | true |  |  | change→postValue; -→postValue |
| `.Remarks` | pxDropdown | Text Area |  |  |  | change→postValue |
| `.SurveyedBy` | pxTextInput | Surveyed by (Ceding Co) |  |  |  | change→postValue |

### 5.S `Section/InputHistoricalSurveyReportDtlUW`

Kelas `ASM-FW-GISFW-Data-Quotation`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.DateofSurvey` | pxDateTime | Date of Survey |  | true |  | change→postValue; -→postValue |
| `.Remarks` | pxDropdown | Text Area |  | true |  | change→postValue |
| `.SurveyedBy` | pxTextInput | Surveyed by (Ceding Co) |  | true |  | change→postValue |

### 5.S `Section/InstallmentList`

Kelas `ASM-FW-GISFW-Data-Installment`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.DueDate` | FIELD |  |  |  |  |  |
| `.InstallmentPercentage` | FIELD |  |  |  |  |  |
| `.Currency` | pxTextInput | Text Input |  |  |  |  |
| `.Premium` | pxTextInput | Text Input |  |  |  |  |
| `.PremiumAfterPPN` | pxTextInput | Text Input |  |  |  |  |
| `.PremiumAfterTax` | pxTextInput | Text Input |  |  |  |  |

### 5.S `Section/Installments_ReadOnly`

Kelas `ASM-FW-GISFW-Data-TreatyInInstallment`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.Installment` | pxTextInput | Text Input |  | true |  |  |
| `.DueDate` | pxDateTime |  |  | true |  |  |
| `.WPC` | pxTextInput | Text Input |  | true |  |  |
| `.PaymentDate` | pxDateTime |  |  | true |  |  |
| `.InstallmentPct` | pxNumber | Text Input |  | true |  |  |
| `.Amount` | pxNumber | Text Input |  | true |  |  |
| `.PctTotal` | pxNumber | % Total |  | nonaktif `` |  |  |
| `.AmountTotal` | pxNumber | Total |  | nonaktif `` |  |  |

### 5.S `Section/ListSuggest`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.IsApproved` | pxRadioButtons | Approval | true |  | selalu (sisa syarat tak berlaku: `pyWorkPage.PositionNote !='ReasTreatyInSecHead'`) | change→runActivity SetDueTo_act; change→refresh CekLimitTreatyAcc_Act; change→refresh Protection_Act; change→refresh DetailPolicyTreatyIn; change→refresh DetailDeptHeadTreatyIn_UW; change→refresh DetailPolicyTreatyInAddendum; click→runActivity SetDueTo_act; click→refresh CekLimitTreatyAcc_Act; click→refresh Protection_Act; click→refresh DetailPolicyTreatyIn; click→refresh DetailDeptHeadTreatyIn_UW; click→refresh DetailPolicyTreatyInAddendum; -→runActivity SetDueTo_act; -→refresh CekLimitTreatyAcc_Act; -→refresh Protection_Act; -→refresh DetailPolicyTreatyIn; -→refresh DetailDeptHeadTreatyIn_UW; -→refresh DetailPolicyTreatyInAddendum |
| `.ProductionDate` | pxDateTime | Production Date | true `.IsApproved == 1 && (OperatorID.pyUserIdentifier=='<ID-operator-3>' \\|\\| OperatorID.pyUserIdentifier=='<ID-operator-4>')` |  | OTHER `.IsApproved == 1 && (OperatorID.pyUserIdentifier=='<ID-operator-3>' \\|\\| OperatorID.pyUserIdentifier=='<ID-operator-4>')` |  |
| `.Suggest` | pxTextArea | Suggest | true |  |  | change→refresh Protection_Act; change→refresh DetailPolicyTreatyIn; change→refresh DetailDeptHeadTreatyIn_UW; -→refresh Protection_Act; -→refresh DetailPolicyTreatyIn; -→refresh DetailDeptHeadTreatyIn_UW |
| `.Date` | pxDateTime |  |  | true | NOTBLANK |  |
| `.OperatorName` | pxDisplayText |  |  | true | NOTBLANK |  |
| `.IsApproved` | pxDropdown |  |  | true |  |  |
| `.Suggest` | pxDisplayText |  |  | true | NOTBLANK |  |

### 5.S `Section/PolicyTreatyInDeclineConfirm`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.pyTemplateButton` | pxButton | Yes |  |  |  | click→finishAssignment |
| `.pyTemplateButton` | pxButton | No |  |  |  | click→closeContainer |

### 5.S `Section/SFAPortal_Opportunities`

Kelas `Data-Portal`
- menyertakan: `SFAPortal_OpportunitiesList`, `SFAPortal_OpportunitiesList_Header`, `undefined`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `SFAPortal_OpportunitiesList_Header` | SUB_SECTION ⟵ `SFAPortal_OpportunitiesList_Header` | SFAPortal_OpportunitiesList_Header |  |  | OTHER `1=2` |  |
| `SFAPortal_OpportunitiesList` | SUB_SECTION ⟵ `undefined` | Shows a list of opportunities in portal |  |  |  |  |

### 5.S `Section/SFAPortal_OpportunitiesList`

Kelas `Data-Portal`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.FilterTermForOpportunity` | pxTextInput | Filter Term for Opportunity |  |  |  | keyboard→refresh; keyboard→setValue; -→refresh |
| `.pyTemplateInputBox` | pxIcon | Image |  |  |  | click→setValue; click→postValue; -→setValue; -→postValue |
| `.pyTemplateInputBox` | pxButton | Filter |  |  |  | keyboard→setValue; keyboard→refresh; click→refresh; -→refresh |
| `.TextNoQuotation` | pxDisplayText |  |  | true |  |  |
| `.Name` | pxLink | .Name |  | nonaktif `.pxPages(A).QuotationData.MarketingName != InputParam.CARI34 && InputParam.CARI38 != .pxCreateOperator && OperatorID.pyWorkBasketList(2).pyWorkBasketName != 'ReasFacInTeamLeader'` |  | click→openWorkByHandle |
| `A.Quotation.BusinessName` | pxDisplayText |  |  | true |  |  |
| `A.Quotation.InsuredName` | pxDisplayText |  |  | true |  |  |
| `A.Quotation.MarketingName` | FIELD |  |  |  |  |  |
| `A.NBStatus` | pxDisplayText |  |  | true |  |  |

### 5.S `Section/SFAPortal_OpportunitiesList_Header`

Kelas `Data-Portal`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `SFAResults.pxResults(1).pySummaryCount(1)` | pxNumber | Prospek Bisnis |  | true |  |  |
| `.pyTemplateInputBox` | pxDisplayText | Estimasi Total Premi |  | true | OTHER `SFAResults.pxResults(1).pySummaryValue(1) == ''` |  |
| `SFAResults.pxResults(1).pySummaryValue(1)` | pxCurrency | Estimasi Total Premi |  | true | NOTBLANK |  |
| `SFAResults.pxResults(1).pySummaryDateTime(1)` | pxDateTime | Estimasi Tanggal Penutupan |  | true | NOTBLANK `crmIsReview` |  |
| `.pyTemplateInputBox` | pxDisplayText | Estimasi Tanggal Penutupan |  | true | OTHER `SFAResults.pxResults(1).pySummaryDateTime(1) == ''` |  |
| `.pyTemplateDisplayText` | pxDisplayText | . |  |  |  |  |
| `SFAResults.pxResults(1).pySummaryDateTime(2)` | pxDateTime | . |  | true | NOTBLANK `crmIsReview` |  |
| `.pyTemplateInputBox` | pxDisplayText | . |  | true | OTHER `SFAResults.pxResults(1).pySummaryDateTime(2) == ''` |  |

### 5.S `Section/SFAPortalOpportunitiesHeader`

Kelas `PegaCRM-Portal`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.pyTemplateInputBox` | pxButton | Stage view |  | nonaktif `.CurrentMode = 'StageView' \\|\\| .CurrentMode = ''` |  | click→refresh SFAPopulateListHeader(FromFilters="false", RDName="crmOpportunitiesList", SellingMode=.FilterForOpportunityInd); click→refresh SFAPortal_Opportunities; -→refresh SFAPortal_Opportunities |
| `.pyTemplateInputBox` | pxButton | List view |  | nonaktif `.CurrentMode = 'List'` |  | click→setValue; click→refresh; click→refresh SFAPortal_Opportunities; -→refresh SFAPortal_Opportunities |
| `.pyTemplateInputBox` | pxButton | Create Opportunity |  |  | OTHER `1=2` | click→showMenu |
| `.pyTemplateInputBox` | pxButton | Create opportunity |  |  | OTHER `crmCreateOpportunity && isSellingModeB2BB2C && IsNotAdmin` | click→createWork |
| `.CurrentMode` | pxHidden | CurrentMode |  |  |  |  |
| `.pyTemplateInputBox` | pxButton | Create opportunity |  |  | OTHER `crmCreateOpportunity && isSellingModeB2B&&IsNotAdmin` | click→createWork |
| `.pyTemplateInputBox` | pxButton | Create opportunity |  |  | OTHER `crmCreateOpportunity && isSellingModeB2C&&IsNotAdmin` | click→createWork |

### 5.S `Section/ShowPolicyNoTreaty_SC`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `pyWorkPage.pyID` | FIELD | .pyID |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.PolicyNo` | FIELD | .CARI2 |  | true |  |  |
| `.pyTemplateButton` | pxButton | OK |  |  |  | click→closeContainer; click→finishAssignment; -→closeContainer; -→finishAssignment |

### 5.S `Section/SourceHierarki`

Kelas `ASM-FW-GISFW-Data-OfferTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `SearchSOB.CARI1` | pxAutoComplete ⟵ `BrowseAgentNusaRe_RD` | Search |  |  |  |  |
| `.ClientName` | pxTextInput | .NAME |  |  |  |  |
| `.pyTemplateButton` | pxButton | Choose |  |  | OTHER `.ChildCount = 0` | click→runActivity SearchHierarkiSourceBizAgentTreatyIn_Act; click→runScript script:opener.location.reload; click→runScript script:window.close; -→runActivity SearchHierarkiSourceBizAgentTreatyIn_Act; -→runScript script:opener.location.reload; -→runScript script:window.close |

### 5.S `Section/SpreadingRiskList`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |
| --- | --- | --- | --- | --- | --- | --- |
| `.pyTemplateButton` | pxButton | Add |  |  | OTHER `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy' && (pyWorkPage.TreatyIn.FacultativeShare = 0 \\|\\| pyWorkPage.TreatyIn.FacultativeShare = '')` | click→addRow; -→addRow |
| `.TreatyType` | pxDropdown ⟵ `ListSpreading.pxResults` |  |  | true `pyWorkPage.TreatyIn.FacultativeShare >0` |  | change→postValue; change→refresh; -→postValue; -→refresh |
| `.SharePercentage` | pxNumber |  |  | true `pyWorkPage.TreatyIn.FacultativeShare >0` |  | change→refresh CountSpreading_Act(SharePct=.SharePercentage, Index=.pxListSubscript, Type=Premium); -→refresh CountSpreading_Act(SharePct=.SharePercentage, Index=.pxListSubscript, Type=Premium) |
| `.PremiumSpreaded` | pxNumber |  |  | true |  |  |
| `.ClaimPercentage` | pxNumber |  |  | true `pyWorkPage.TreatyIn.FacultativeShare >0` |  | change→refresh CountSpreading_Act(SharePct=.ClaimPercentage, Index=.pxListSubscript, Type=Claim); -→refresh CountSpreading_Act(SharePct=.ClaimPercentage, Index=.pxListSubscript, Type=Claim) |
| `.ClaimSpreaded` | pxNumber |  |  | true |  |  |
| `.pyTemplateButton` | pxButton | Delete |  |  | OTHER `InputParam.CARI12 != 'treaty' && InputParam.CARI12 != 'TreatyPolicy' && (pyWorkPage.TreatyIn.FacultativeShare = 0 \\|\\| pyWorkPage.TreatyIn.FacultativeShare = '')` | click→deleteRow; -→deleteRow |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentagePremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalPremium` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalSharePercentageClaim` | pxNumber |  |  | true |  |  |
| `pyWorkPage.PolicyTreatyIn.TotalClaim` | pxNumber | .TotalClaim |  | true |  |  |
| `.TotalSharePercentagePremium` | pxNumber | Total %Share |  | true |  |  |
| `.TotalPremium` | pxNumber | Total Premium |  | true |  |  |
| `.TotalSharePercentageClaim` | pxNumber | Total %Share Claim |  | true |  |  |
| `.TotalClaim` | pxNumber | Total Claim |  | true |  |  |

## 6 · Activity terjangkau — langkah demi langkah

Kode aksi syarat/transisi Pega dicantumkan apa adanya beserta tafsirannya: `1` lompat ke label · `2` lanjut · `3` lewati langkah · `4` keluar iterasi · `5` lewati syarat berikut · `6` keluar activity. Syarat langkah (`when`) hanya berlaku bila kotak *When* langkah itu aktif (`whenAktif=true`).

### 6 · `AgentSourceBizTreatyIn_Act`

Kelas `ASM-FW-GISFW-Data-Agent` · 65.261 B · parameter: `Leader`, `CedingCoLeader`, `ID`

- **1** `Property-Set`
  - syarat `pyWorkPage.OfferTreatyIn.QuotationData.btnQuotation=="SOB"`: benar → ; salah → lewati langkah [kode -/3]
  - `Param.pyReportName` = `"BrowseAgentHierarkiList_RD"`
  - `Param.pyReportClass` = `"ASM-FW-GISFW-Int-AGENT"`
  - `Param.Leader` = `Primary.ID`
- **2** `Property-Set`
  - syarat `pyWorkPage.OfferTreatyIn.QuotationData.btnQuotation=="CedingCo"`: benar → ; salah → lewati langkah [kode -/3]
  - `Param.pyReportName` = `"BrowseCedingCo_RD"`
  - `Param.pyReportClass` = `"ASM-FW-GISFW-Int-AGENT"`
  - `Param.ID` = `"A"`
  - `Param.CedingCoLeader` = `Primary.ID`
- **3** `Call pxShowReport` · halaman `Report`
  - syarat `pyWorkPage.Quotation.btnQuotation=="SOB"`: benar → ; salah → lewati langkah [kode -/3] *(When langkah tidak dicentang)*
  - parameter: pyReportName=`BrowseAgentHierarkiList_RD`; pyReportClass=`ASM-FW-GISFW-Int-AGENT`
- **4** `Property-Set` · halaman `pyReportContentPage.pxResults`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `.pxObjClass` = `"ASM-FW-GISFW-Data-Agent"`
  - **4.1** `Page-Copy`
    - parameter: CopyFrom=`pyReportContentPage.pxResults(<CURRENT>)`; CopyInto=`Primary.pxResults(<APPEND>)`

### 6 · `BreakDownSpreading_Act`

Kelas `ASM-FW-GISFW-Work` · 74.825 B

- **1** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.SpreadingRiskList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **1.1** `Property-Set`
    - `ParamData.CARI10` = `pyWorkPage.PolicyTreatyIn.TreatyGroupID`
    - `ParamData.CARI11` = `.TreatyType`
    - `ParamData.CARI12` = `pyWorkPage.PolicyTreatyIn.TreatyYear`
  - **1.2** `Page-New` · halaman `BreakdownSpreadList`
  - **1.3** `RDB-List` · halaman `BreakdownSpreadList`
    - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Work`; RequestType=`GetBreakDownSpread_SQL`
  - **1.4** `Property-Set`
    - `pyWorkPage.PolicyTreatyIn.BreakDownSpreadList` = `BreakdownSpreadList.pxResults`
- **2** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.BreakDownSpreadList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **2.1** `Property-Set`
    - `.PremiumSpreaded` = `pyWorkPage.PolicyTreatyIn.TotalPremium * @divide(.SharePercentage,100,8)`
    - `.ClaimSpreaded` = `@toDecimal(pyWorkPage.PolicyTreatyIn.TotalClaim) * @divide(.SharePercentage,100,8)`

### 6 · `CalculatePremi_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 39.683 B · parameter: `Action`

- **1** `Property-Set`
  - syarat `Param.Action=="PREMIUM"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.PremiOgp` = `@divide(.GrossPremium*pyWorkPage.TreatyIn.RNMShareP,100,4)`
  - `.ResultOgp1` = `@divide(.PremiOgp*pyWorkPage.TreatyIn.BrokeragePercentP,100,4)`
  - `.RiCommOgp` = `@divide(.ResultOgp1,.PremiOgp,4) * 100`
- **2** `Property-Set`
  - syarat `Param.Action=="CLAIM"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.Claim` = `@divide(.GrossClaim*pyWorkPage.TreatyIn.RNMShareP,100,4)`
- **3** `Call CountNetPremi_act`

### 6 · `CekLimitTreatyAcc_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 68.006 B

- **1** `Property-Set` — pyWorkIDPrefix EDMT-
  - `pyWorkPage.LetterNo` = `""`
- **2** `(kosong)`
  - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefLimit=`1`, pyStepsRepeatDefStart=`1`
  - syarat `pyWorkPage.PositionNote=="ReasTreatyInSecHead"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **2.1** `Property-Set`
    - `Local.IDCurr` = `.IDCurrency`
    - `Local.NETPREMI` = `.TotalPremium`
    - `Local.CurrName` = `.Currency`
    - `Local.NETPREMI` = `@if(Local.NETPREMI<0.0,Local.NETPREMI*-1,Local.NETPREMI)`
  - **2.2** `Property-Set` · label `//` — pyWorkPage.pyWorkIDPrefix=="EDMT-"
    - syarat `pyWorkPage.pyWorkIDPrefix=="EDMT-"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `Local.NETPREMI` = `pyWorkPage.PolicyTreatyIn.TreatyDifference.TotalPremium`
    - `Local.NETPREMI` = `@if(Local.NETPREMI<0.0,Local.NETPREMI*-1,Local.NETPREMI)`
  - **2.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.CurrencyList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - syarat : benar → lanjut; salah → lewati langkah [kode 2/3] *(When langkah tidak dicentang)*
    - **2.3.1** `Property-Set` — @if(Local.NETPREMI>200000000.00,"TREATYINDEPTHEAD","")
      - syarat `Local.IDCurr==.CurrencyID \|\| Local.CurrName==.Currency`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `Local.NETPREMI` = `Local.NETPREMI * .Conversion`
  - **2.4** `Property-Set` — kalau lebih 200jt set flah untuk naik keatasan
    - `pyWorkPage.LetterNo` = `@if(Local.NETPREMI>200000000.00,"TREATYINDEPTHEAD","")`

### 6 · `CheckDataMkt`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 112.503 B

- **1** `Page-New` · halaman `DataMarketing`
- **2** `Property-Set`
  - `.MarketingOfficer` = `""`
- **3** `Obj-Browse`
  - syarat `pyWorkPage.PolicyTreatyIn.QuotationData.MOID==""`: benar → lewati langkah; salah → lanjut [kode 3/2]
  - `Condition==; Field=.ID; Value=pyWorkPage.PolicyTreatyIn.QuotationData.MOID; Select=true`
  - `Field=.ClientID; Select=true`
  - `Field=.BranchDetailID; Select=true`
  - `Field=.BranchDetailName; Select=true`
  - `Field=.ClientName; Select=true`
  - `Field=.TeamGroup; Select=true`
  - parameter: ObjClass=`ASM-FW-GISFW-Int-marketingofficer`; RowKey=`ID`; PageName=`DataMarketing`; UseLightWeightList=`true`; GetRowKey=`true`
- **4** `Property-Set`
  - `pyWorkPage.Quotation.MOID` = `DataMarketing.pxResults(1).ID`
  - `pyWorkPage.Quotation.MarketingCode` = `DataMarketing.pxResults(1).ClientID`
  - `pyWorkPage.Quotation.MarketingName` = `DataMarketing.pxResults(1).ClientName`
  - `pyWorkPage.Quotation.TeamGroup` = `DataMarketing.pxResults(1).TeamGroup`
  - `pyWorkPage.Quotation.BranchCode` = `DataMarketing.pxResults(1).BranchDetailID`
  - `pyWorkPage.Quotation.BranchName` = `DataMarketing.pxResults(1).BranchDetailName`
  - `pyWorkPage.OfferFacIn.QuotationData.MOID` = `DataMarketing.pxResults(1).ID`
  - `pyWorkPage.OfferFacIn.QuotationData.MarketingCode` = `DataMarketing.pxResults(1).ClientID`
  - `pyWorkPage.OfferFacIn.QuotationData.MarketingName` = `DataMarketing.pxResults(1).ClientName`
  - `pyWorkPage.OfferFacIn.QuotationData.TeamGroup` = `DataMarketing.pxResults(1).TeamGroup`
  - `.QuotationData.MOID` = `DataMarketing.pxResults(1).ID`
  - `.QuotationData.MarketingCode` = `DataMarketing.pxResults(1).ClientID`
  - `.QuotationData.MarketingName` = `DataMarketing.pxResults(1).ClientName`
  - `.QuotationData.TeamGroup` = `DataMarketing.pxResults(1).TeamGroup`
  - `.MarketingOfficer` = `DataMarketing.pxResults(1).ClientName`
- **5** `Obj-Save` · halaman `pyWorkPage`
  - parameter: WithErrors=`true`; OnlyIfNew=`false`; WriteNow=`true`
- **6** `Property-Set-Messages` · label `//`
  - syarat `.MarketingCode=="" && .MarketingName!=""`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.MarketingOfficer`; Category=`pyMessageLabel`; Message=`ASMMessageQuotation`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **7** `Property-Set-Messages` · label `//`
  - syarat `.SobName!="" && .SourceOfBusiness==""`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.SobName`; Category=`pyMessageLabel`; Message=`ASMMessageQuotationSumbis`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`

### 6 · `CheckDuplicateOffer`

Kelas `Data-Portal` · 92.574 B · parameter: `sob`, `ceding`, `startdate`, `enddate`, `ID`

- **1** `Property-Set` · label `//` — Init RD
  - `Param.pyReportName` = `"BrowseTREATY_IN"`
  - `Param.pyReportClass` = `"ASM-FW-GISFW-Int-TREATY_IN"`
  - `Param.CedingID` = `TreatyIn.CedingID`
  - `Param.LeadingReinsSourceID` = `TreatyIn.LeadingReinsSourceID`
  - `Param.Commencement` = `TreatyIn.Commencement`
  - `Param.Termination` = `TreatyIn.Termination`
  - `TreatyWarning.CARI1` = `""`
  - `Local.Flag` = `0`
- **2** `call pxShowReport` · halaman `ReportPage` · label `//` — Call RD
- **3** `(kosong)` · halaman `pyReportContentPage.pxResults` · label `//` — Set warning when similarity found on Ceding, SoB, Period
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **3.1** `Property-Set` — When Similarity Found in any
    - syarat `.ID != TreatyIn.ID && .ProportionType == TreatyIn.ProportionType`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - syarat `.TreatyYear==TreatyIn.TreatyYear`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - syarat `.Commencement==TreatyIn.Commencement`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - syarat `.Termination==.Termination`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - syarat `.CedingID==TreatyIn.CedingID`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - syarat `.LeadingReinsSourceID==TreatyIn.LeadingReinsSourceID`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `TreatyWarning.CARI1` = `TreatyWarning.CARI1 +"ID " + .ID + ", "`
    - `Local.Flag` = `1`
- **4** `Property-Set` · label `//`
  - syarat `Local.Flag==1`: benar → ; salah → lewati langkah [kode -/3]
  - `TreatyWarning.CARI1` = `TreatyWarning.CARI1 + "Have similarities in SoB, Ceding, and Period, please check for duplicates"`
- **5** `Property-Set` — set Id master
  - `TreatyWarning.CAIREINSFACIN` = `Param.ID`
- **6** `RDB-List` — get data claim by idmaster
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`RNM`; ClassName=`Assign-Worklist`; BrowsePage=`HasilClaim`; RequestType=`GetCountClaim`
- **7** `Property-Set` — set Error message
  - syarat `HasilClaim.pxResults(1).CARI1>0`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `Local.ErrMSG` = `"Sudah Ada Claim Untuk IDMASTER ini , Tidak bisa Revisi"`
- **8** `Page-Set-Messages` — make error jika ada klaim
  - syarat `HasilClaim.pxResults(1).CARI1>0`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: Category=`pyMessageLabel`; ClassOfPage=`ASM-FW-GISFW-Int-TREATY_IN`; Message=`Local.ErrMSG`; ContainingClassOfProperty=`@baseclass`; Page=`TreatyIn`

### 6 · `ConcatSlipOfferNo_Act`

Kelas `ASM-FW-GISFW-Data-Quotation` · 30.516 B

- **1** `Property-Set`
  - `pyWorkPage.Quotation.NoOfferSlip` = `""`
- **2** `(kosong)` · halaman `.QuotationList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **2.1** `Property-Set`
    - `pyWorkPage.Quotation.NoOfferSlip` = `pyWorkPage.Quotation.NoOfferSlip + .NoOfferSlip + "; "`

### 6 · `ConvertHistoryDate`

Kelas `Data-Portal` · 26.973 B

- **1** `(kosong)` · halaman `TreatyIn.CommentList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **1.1** `Property-Set`
    - syarat `.ConvertDate==1`: benar → lewati langkah; salah → lanjut [kode 3/2]
    - `.Date` = `@DateTime.addCalendar(.Date,0,0,0,0,7,0,0)`
    - `.ConvertDate` = `1`

### 6 · `CountNetPremi_act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 84.931 B

- **1** `Call SetPPNPPH`
- **2** `Property-Set` — prop
  - syarat `pyWorkPage.PolicyTreatyIn.QuotationData.ProportionalType=="NonProportional"`: benar → lewati langkah; salah → lanjut [kode 3/2]
  - `.GrossClaim` = `@divide(.Claim,pyWorkPage.TreatyIn.RNMShareP,4) * 100`
- **3** `Property-Set` — "NonProportional"
  - syarat `pyWorkPage.PolicyTreatyIn.QuotationData.ProportionalType=="NonProportional"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.GrossClaim` = `@divide(.Claim,pyWorkPage.TreatyIn.RNMShare,4) * 100`
- **4** `Property-Set`
  - `.NetPremium` = `((.PremiOgp-.ResultOgp1)+(.PremiOnp-.ResultOnp1))-(.OveriddingCommOgp*.PremiOgp/100)-(.OveriddingCommOnp*.PremiOnp/100)-.Deduction1-.Deduction2`
- **5** `Property-Set`
  - `.BalanceDueTo` = `.NetPremium-.Claim-.ExcessLoss+.SalvageValue`
- **6** `Property-Set`
  - syarat `.Deduction1!="" \|\| .Deduction1!=0`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.BalanceDueTo` = `.NetPremium +.PPNValue+.PPHValue-.Claim-.ExcessLoss+.SalvageValue`
  - `.BalanceBeforePPH` = `.NetPremium +.PPNValue-.Claim-.ExcessLoss+.SalvageValue`
  - `.BalanceBeforeTax` = `.NetPremium -.Claim-.ExcessLoss+.SalvageValue`
- **7** `call SetDueTo_act`

### 6 · `CountOGPONP_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 87.729 B

- **1** `Property-Set` — .PremiOgp null or 0
  - syarat `.PremiOgp == "" \|\|.PremiOgp ==  "0"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.PremiOgp` = `0`
  - `.RiCommOgp` = `0`
  - `.ResultOgp1` = `0`
  - `.OveriddingCommOgp` = `0`
  - `.ResultOgp2` = `0`
- **2** `Property-Set` — PremiOnp null or 0
  - syarat `.PremiOnp=="" \|\| .PremiOnp=="0"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.PremiOnp` = `0`
  - `.RiCommOnp` = `0`
  - `.ResultOnp1` = `0`
  - `.OveriddingCommOnp` = `0`
  - `.ResultOnp2` = `0`
- **3** `Call CountResult1_Act`
- **4** `Call CountResult2Ogp_act`
- **5** `Call CountResult1Onp_Act`
- **6** `Call CountResult2Onp_act`
- **7** `Call CountNetPremi_act`
- **8** `Property-Set` — CEK ULANG JIKA ClaimType = NULL ATAU ClaimPaymentType = NULL
  - syarat `.ClaimType = "" \|\|.ClaimPaymentType = ""`: benar → lanjut; salah → lewati langkah [kode 2/3] *(When langkah tidak dicentang)*
  - `.ClaimType` = `@if((.Claim!=0&&.Claim!="")\|\|(.SalvageValue!=0&&.SalvageValue!=""),"SOA","")`
  - `.ClaimPaymentType` = `@if((.Claim!=0&&.Claim!="")\|\|(.SalvageValue!=0&&.SalvageValue!=""),"Claim","")`
- **9** `(kosong)` · halaman `.SpreadingRiskList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `@LengthOfPageList(Primary.SpreadingRiskList)>0`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **9.1** `Call CountSpreading_Act` · halaman `Primary`
    - parameter: Index=`.SpreadingRiskList(<CURRENT>).pxListSubscript`
- **10** `Call SetValidateInstallment_Act`

### 6 · `CountOverridingCommOgp_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 50.085 B · parameter: `Overidding` · keterangan: THIS ACTIVITY MADE TO COUNT OVERRIDING DIFFERENCE NOT BIGGER THEN 0,01 IF RESULT OGP 2 CHANGED

- **1** `Property-Set-Messages` · label `//`
  - syarat `((.ResultOgp2/.PremiOgp)-.OveriddingCommOgp)<=0.01`: benar → lewati langkah; salah →  [kode 3/-]
  - parameter: Field=`.OveriddingCommOgp`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **2** `Property-Set` · label `//`
  - syarat `((.ResultOgp2/.PremiOgp)-.OveriddingCommOgp)<=0.01`: benar → lewati langkah; salah →  [kode 3/-]
  - `.ResultOgp2` = `.PremiOgp*.OveriddingCommOgp`
- **3** `Property-Set`
  - syarat `Param.Overidding=="Amount"`: benar → ; salah → lewati langkah [kode -/3]
  - `.OveriddingCommOgp` = `@Math.divide(.ResultOgp2,.PremiOgp,4)*100`
- **4** `Property-Set-Messages` · label `//`
  - syarat `((.ResultOgp2/.PremiOgp)-.OveriddingCommOgp)<=0.01`: benar → lewati langkah; salah →  [kode 3/-]
  - parameter: Field=`.OveriddingCommOgp`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **5** `call CountNetPremi_act`

### 6 · `CountOverridingCommOnp_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 42.832 B · parameter: `overidding` · keterangan: THIS ACTIVITY MADE TO COUNTING OVERRIDING COMM ONP DIFFERENCE NOT BIGGER THEN 0,01 IF RESULT ONP 2 CHANGED

- **1** `Property-Set-Messages` · label `//`
  - syarat `((.ResultOnp2/.PremiOnp)-.OveriddingCommOnp)<=0.01`: benar → lewati langkah; salah →  [kode 3/-]
  - parameter: Field=`.OveriddingCommOnp`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **2** `Property-Set` · label `//`
  - syarat `((.ResultOnp2/.PremiOnp)-.OveriddingCommOnp)<=0.01`: benar → lewati langkah; salah →  [kode 3/-]
  - `.ResultOnp2` = `.PremiOnp*.OveriddingCommOnp`
- **3** `Property-Set`
  - syarat `Param.overidding=="Amount"`: benar → ; salah → lewati langkah [kode -/3]
  - `.OveriddingCommOnp` = `@Math.divide(.ResultOnp2,.PremiOnp,4)*100`
- **4** `call CountNetPremi_act`

### 6 · `CountPctInstallment_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 26.105 B · parameter: `idx`

- **1** `Property-Set`
  - `.ListInstallment(Param.idx).InstallmentPercentage` = `@Math.divide(@toDecimal(.ListInstallment(Param.idx).Premium),@toDecimal(.BalanceDueTo),4)*100`
  - `.ListInstallment(Param.idx).PaymentTotal` = `.ListInstallment(Param.idx).Premium`

### 6 · `CountResult1_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 79.689 B · parameter: `Data` · keterangan: THIS ACTIVITY USED FOR COUNTING OGP RESULT 1

- **1** `(kosong)` — .PremiOgp == "" \|\|.PremiOgp ==  "0"
  - syarat `.PremiOgp == "" \|\|.PremiOgp ==  "0"`: benar → keluar activity; salah → lanjut [kode 6/2]
- **2** `Property-Set-Messages`
  - syarat `@String.isDouble(.RiCommOgp)\|\|.RiCommOgp==""`: benar → lanjut; salah → keluar activity [kode 2/6]
  - syarat `.RiCommOgp>100`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.RiCommOgp`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **3** `Property-Set` — hitung GrossPremium utk prop
  - syarat `pyWorkPage.PolicyTreatyIn.QuotationData.ProportionalType=="NonProportional"`: benar → lewati langkah; salah → lanjut [kode 3/2]
  - syarat `Param.Data==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.GrossPremium` = `@divide(.PremiOgp*100 ,pyWorkPage.TreatyIn.RNMShareP ,4)`
- **4** `Property-Set` — hitung GrossPremium utk nonprop
  - syarat `pyWorkPage.PolicyTreatyIn.QuotationData.ProportionalType=="NonProportional"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - syarat `Param.Data==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.GrossPremium` = `@divide(.PremiOgp*100 ,pyWorkPage.TreatyIn.RNMShare ,4)`
- **5** `Property-Set`
  - syarat `(.RiCommOgp<100 && @String.isDouble(.RiCommOgp)) \|\| .RiCommOgp==""`: benar → ; salah → keluar activity [kode -/6]
  - syarat `Param.Data=="Pct"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.ResultOgp1` = `.PremiOgp*(.RiCommOgp/@String.toDecimal("100.00000"))`
- **6** `Property-Set`
  - syarat `Param.Data=="Amount"`: benar → ; salah → lewati langkah [kode -/3]
  - `.RiCommOgp` = `(.ResultOgp1/.PremiOgp)*@String.toDecimal("100.00000")`
- **7** `call CountNetPremi_act`

### 6 · `CountResult1Onp_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 50.025 B · parameter: `Data`

- **1** `(kosong)`
  - syarat `.PremiOnp=="" \|\| .PremiOnp=="0"`: benar → keluar activity; salah → lanjut [kode 6/2]
- **2** `Property-Set-Messages`
  - syarat `@String.isDouble(.RiCommOnp)\|\|.RiCommOnp==""`: benar → lanjut; salah → keluar activity [kode 2/6]
  - syarat `.RiCommOnp>100`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.RiCommOnp`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **3** `Property-Set`
  - syarat `(.RiCommOnp<100 && @String.isDouble(.RiCommOnp))\|\|.RiCommOnp=""`: benar → ; salah → keluar activity [kode -/6]
  - syarat `Param.Data=="Pct"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.ResultOnp1` = `.PremiOnp*(.RiCommOnp/@String.toDecimal("100.00000"))`
- **4** `Property-Set`
  - syarat `Param.Data=="Amount"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.RiCommOnp` = `(.ResultOnp1/.PremiOnp)*@String.toDecimal("100.00000")`
- **5** `call CountNetPremi_act`

### 6 · `CountResult2Ogp_act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 48.162 B · parameter: `Data`

- **1** `(kosong)`
  - syarat `.PremiOgp == "" \|\|.PremiOgp ==  "0"`: benar → keluar activity; salah → lanjut [kode 6/2]
- **2** `Property-Set-Messages`
  - syarat `@String.isDouble(.OveriddingCommOgp)`: benar → lanjut; salah → keluar activity [kode 2/6]
  - syarat `.OveriddingCommOgp>100`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.OveriddingCommOgp`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **3** `Property-Set`
  - syarat `(.OveriddingCommOgp<100 && @String.isDouble(.OveriddingCommOgp)) \|\| .OveriddingCommOgp==""`: benar → ; salah → keluar activity [kode -/6]
  - syarat `Param.Data=="Pct"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.ResultOgp2` = `(.OveriddingCommOgp/@String.toDecimal("100.00000"))*.PremiOgp`
- **4** `Property-Set`
  - syarat `Param.Data=="Amount"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.OveriddingCommOgp` = `(.ResultOgp2/.PremiOgp)*@String.toDecimal("100.00000")`
- **5** `call CountNetPremi_act`

### 6 · `CountResult2Onp_act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 48.368 B · parameter: `Data`

- **1** `(kosong)`
  - syarat `.PremiOnp=="" \|\| .PremiOnp=="0"`: benar → keluar activity; salah → lanjut [kode 6/2]
- **2** `Property-Set-Messages`
  - syarat `@String.isDouble(.OveriddingCommOnp)`: benar → lanjut; salah → keluar activity [kode 2/6]
  - syarat `.OveriddingCommOnp>100`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.OveriddingCommOnp`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **3** `Property-Set`
  - syarat `.OveriddingCommOnp<100 && @String.isDouble(.OveriddingCommOnp)`: benar → ; salah → keluar activity [kode -/6]
  - syarat `Param.Data=="Pct"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.ResultOnp2` = `(.OveriddingCommOnp/@String.toDecimal("100.00000"))*.PremiOnp`
- **4** `Property-Set`
  - syarat `Param.Data=="Amount"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.OveriddingCommOnp` = `(.ResultOnp2/.PremiOnp) *@String.toDecimal("100.00000")`
- **5** `call CountNetPremi_act`

### 6 · `CountRiCommOgp_act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 48.839 B · parameter: `DiscountType` · keterangan: THIS ACTIVITY MADE TO COUNT RI COMM OGP DIFFERENCE NOT BIGGER THEN 0,01 IF RESULT OGP 1 CHANGED

- **1** `Property-Set-Messages` · label `//`
  - syarat `@String.equals(Local.STS,"False")`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.ResultOgp1`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **2** `Property-Set` · label `//`
  - syarat `@String.equals(Local.STS,"False")`: benar → ; salah → lewati langkah [kode -/3]
  - `.ResultOgp1` = `.PremiOgp*.RiCommOgp`
- **3** `Property-Set`
  - syarat `Param.DiscountType=="Amount"`: benar → ; salah → lewati langkah [kode -/3]
  - `.RiCommOgp` = `@Math.divide(.ResultOgp1,.PremiOgp,4)*100`
- **4** `Property-Set-Messages` · label `//`
  - syarat `((.ResultOgp1/.PremiOgp)-.RiCommOgp)>=0.01`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.RiCommOgp`; Category=`pyMessageLabel`; Message=`ErrorMsg1`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`
- **5** `call CountNetPremi_act`

### 6 · `CountRiCommOnp_act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 33.794 B · parameter: `Result`

- **1** `Property-Set` · label `//`
  - syarat `((.ResultOnp1/.PremiOnp)-.RiCommOnp)<=0.01`: benar → lewati langkah; salah →  [kode 3/-]
  - `.ResultOnp1` = `.PremiOnp*.RiCommOnp`
- **2** `Property-Set`
  - syarat `Param.Result=="Amount"`: benar → ; salah → lewati langkah [kode -/3]
  - `.RiCommOnp` = `@Math.divide(.ResultOnp1,.PremiOnp,4)*100`
- **3** `call CountNetPremi_act`

### 6 · `CountSpreading_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 86.087 B

- **1** `Property-Set` · label `//`
  - syarat `.SpreadingRiskList(Param.Index).SharePercentage == ""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.SpreadingRiskList(Param.Index).SharePercentage` = `100/@toDecimal(@LengthOfPageList(Primary.SpreadingRiskList))`
- **2** `Property-Set` · label `//`
  - syarat `.SpreadingRiskList(Param.Index).ClaimPercentage == ""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `.SpreadingRiskList(Param.Index).ClaimPercentage` = `.SpreadingRiskList(Param.Index).SharePercentage`
- **3** `Property-Set` · label `//`
  - `.SpreadingRiskList(Param.Index).PremiumSpreaded` = `.NetPremium * @divide(.SpreadingRiskList(Param.Index).SharePercentage,100,10)`
  - `.SpreadingRiskList(Param.Index).ClaimSpreaded` = `(.ExcessLoss + .Claim - .SalvageValue)* @divide(.SpreadingRiskList(Param.Index).ClaimPercentage ,100,10)`
- **4** `(kosong)` · halaman `.SpreadingRiskList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **4.1** `Property-Set`
    - `.SharePercentage` = `@if(.SharePercentage=="",(100/@toDecimal(@LengthOfPageList(Primary.SpreadingRiskList))),.SharePercentage)`
    - `.ClaimPercentage` = `@if(.ClaimPercentage == "",.SharePercentage,.ClaimPercentage)`
    - `.PremiumSpreaded` = `Primary.NetPremium * @divide(.SharePercentage,100,10)`
    - `.ClaimSpreaded` = `(Primary.ExcessLoss + Primary.Claim - Primary.SalvageValue)* @divide(.ClaimPercentage ,100,10)`
  - **4.2** `Property-Set`
    - `Local.shareclaim` = `.ClaimPercentage + (Local.shareclaim)`
    - `Local.claim` = `.ClaimSpreaded + (Local.claim)`
    - `Local.sharepremium` = `.SharePercentage + (Local.sharepremium)`
    - `Local.premium` = `.PremiumSpreaded + (Local.premium)`
- **5** `Property-Set`
  - `.TotalSharePercentagePremium` = `Local.sharepremium`
  - `.TotalPremium` = `Local.premium`
  - `.TotalSharePercentageClaim` = `Local.shareclaim`
  - `.TotalClaim` = `Local.claim`
- **6** `Call ASM-FW-GISFW-Work.BreakDownSpreading_Act`

### 6 · `FetchMasterTreatyIn`

Kelas `ASM-FW-GISFW-Work` · 58.236 B

- **1** `Page-New` · halaman `TempResult` — init
- **2** `Page-Remove` · halaman `TreatyIn` — init
- **3** `Page-New` · halaman `TreatyIn` — init
  - parameter: NewClass=`ASM-FW-GISFW-Int-TREATY_IN`
- **4** `RDB-List` · halaman `TempResult` — Run Querry
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Data-PolicyTreatyIn`; Access=`RNM`; BrowsePage=`TempResult`; RequestType=`BrowseTreatyInJoinEDM`
- **5** `Java` — Map to Page pyWorkPage.TreatyIn
  - Java:

```java
try{
  ClipboardPage tempPage2 = tools.findPage("pyWorkPage.TreatyIn");
  ClipboardPage DataJSON = tools.findPage("TempResult");
  ClipboardProperty DataJSONList  = DataJSON.getProperty(".pxResults");
  java.util.Iterator DataJSONListIter = DataJSONList.iterator();
  while (DataJSONListIter.hasNext())
  {
    ClipboardProperty DataJSONData = (ClipboardProperty)DataJSONListIter.next();
    ClipboardPage DataJSONPage = DataJSONData.getPageValue();
    String IsiDataJson = DataJSONPage.getString("CLASSOFBUSINESS");
    tempPage2.adoptJSONObject(IsiDataJson);
  }
} catch(InvalidStreamError e){
  oLog.error("ReloadSection:Invalid JSON Stream for data page params : "+e.getMessage());
} catch(Exception e){
  oLog.error("ReloadSection:Expection : "+e.getMessage());
}
```


### 6 · `FetchTreatyGroupOJK`

Kelas `ASM-FW-GISFW-Work-NB` · 52.382 B · parameter: `ID`, `CobojkOUT`

- **1** `Page-New` · halaman `pyReportContentPage`
- **2** `Property-Set` — Set Report definition parameter
  - `Param.pyReportName` = `"BrowseTreatyGroup_RD"`
  - `Param.pyReportClass` = `"ASM-FW-GISFW-Int-TREATYGROUP"`
  - `Param.ID` = `pyWorkPage.PolicyTreatyIn.TreatyGroupID`
- **3** `Call pxRetrieveReportData` · halaman `ReportPage` — Call Report Definition
  - parameter: pyUseAlternateDB=`false`; pySkipSummaryProcessing=`false`
- **4** `Property-Set` — Set result(1) into property & param out
  - `pyWorkPage.PolicyTreatyIn.OJKBusinessID` = `pyReportContentPage.pxResults(1).OJKBusinessID`
  - `Param.CobojkOUT` = `pyReportContentPage.pxResults(1).OJKBusinessID`
- **5** `Page-Remove` — Cleanup
  - `Page=pyReportContentPage`

### 6 · `FetchTreatyGroupOldID`

Kelas `ASM-FW-GISFW-Work` · 57.316 B

- **1** `Page-Clear-Messages` · halaman `pyWorkPage`
- **2** `Property-Set`
  - `Param.Errmsg` = `"Cannot fetch Treaty Group ID, Contact IT"`
- **3** `Page-Set-Messages`
  - syarat `pyWorkPage.PolicyTreatyIn.TreatyGroupID==""`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Category=`pyMessageLabel`; Message=`Param.Errmsg`; ClassOfPage=`ASM-FW-GISFW-Work`; Page=`pyWorkPage`; ContainingClassOfProperty=`@baseclass`
- **4** `Property-Set`
  - `InputData.CARI1` = `pyWorkPage.PolicyTreatyIn.TreatyGroupID`
- **5** `RDB-List`
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Int-TREATYGROUP`; BrowsePage=`OutputData`; RequestType=`FetchTreatyGroupOLDID`
- **6** `Property-Set`
  - `pyWorkPage.PolicyTreatyIn.TreatyGroupOldID` = `OutputData.pxResults(1).CARI1`

### 6 · `FillPaymentInstallment`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 148.174 B · parameter: `Installment`

- **1** `Property-Remove`
  - `Property=.ListInstallment`
- **2** `Property-Set`
  - `Local.InstallmentNo` = `0`
  - `Local.InstallmentPercentage` = `0`
  - `Local.premi` = `.BalanceDueTo`
  - `Param.Installment` = `.Installment`
- **3** `(kosong)`
  - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefLimit=`Param.Installment-1`, pyStepsRepeatDefStart=`0`
  - **3.1** `Property-Set`
    - `Local.InstallmentNo` = `Local.InstallmentNo + 1`
  - **3.2** `Property-Set`
    - `Local.InstallmentPercentage` = `Local.InstallmentPercentage + @Math.divide(100, Param.Installment, 4)`
    - `.ListInstallment(<APPEND>).pxObjClass` = `"ASM-FW-GISFW-Data-Installment"`
  - **3.3** `Property-Set`
    - syarat `Local.InstallmentNo == Param.Installment && (Local.InstallmentPercentage + @Math.divide(100, Param.Installment, 4)) != 100`: benar → ; salah → lewati langkah [kode -/3]
    - `.ListInstallment(<LAST>).InstallmentPercentage` = `@Math.divide(100, Param.Installment, 4) + (100 - Local.InstallmentPercentage)`
  - **3.4** `Property-Set`
    - syarat `Local.InstallmentNo<Param.Installment`: benar → ; salah → lewati langkah [kode -/3]
    - `.ListInstallment(<LAST>).InstallmentPercentage` = `@Math.divide(100, Param.Installment, 4)`
  - **3.5** `Property-Set`
    - `.ListInstallment(<LAST>).InstallmentNo` = `Local.InstallmentNo`
    - `.ListInstallment(<LAST>).Premium` = `@Math.divide(((Local.premi)*.ListInstallment(<last>).InstallmentPercentage),100,4)`
    - `.ListInstallment(<LAST>).DueDate` = `@CurrentDateTime()`
    - `.ListInstallment(<LAST>).PaymentTotal` = `@Math.divide(((Local.premi)*.ListInstallment(<LAST>).InstallmentPercentage),100,4)`
- **4** `(kosong)` · halaman `.CurrencyList` · label `//`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **4.1** `Property-Set` · label `//`
    - `Local.InstallmentNo` = `0`
    - `Local.InstallmentPercentage` = `0`
  - **4.2** `Property-Remove`
    - `Property=.Policy.Payment.ListInstallment`
  - **4.3** `(kosong)` · label `//`
    - ulang: pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefStart=`0`, pyStepsRepeatDefLimit=`Param.installment-1`
    - **4.3.1** `Property-Set`
      - `Local.InstallmentNo` = `0`
      - `Local.InstallmentPercentage` = `0`
    - **4.3.2** `Property-Set`
- **5** `(kosong)` · halaman `.CurrencyList` · label `//`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **5.1** `Property-Set`
    - `Local.InstallmentNo` = `0`
    - `Local.InstallmentPercentage` = `0`
    - `.SumTotalPayment` = `.Premium`
    - `.Policy.Payment.Premium` = `.Premium`
  - **5.2** `(kosong)`
    - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefStart=`0`, pyStepsRepeatDefLimit=`Param.Installment-1`
    - **5.2.1** `Property-Set`
      - `Local.InstallmentNo` = `Local.InstallmentNo + 1`
    - **5.2.2** `Property-Set`
      - `Local.InstallmentPercentage` = `Local.InstallmentPercentage + @Math.divide(100, Param.Installment, 4)`
      - `.Policy.Payment.ListInstallment(<APPEND>).pxObjClass` = `"ASM-FW-GISFW-Data-Installment"`
    - **5.2.3** `Property-Set`
      - syarat `Local.InstallmentNo == Param.Installment && (Local.InstallmentPercentage + @Math.divide(100, Param.Installment, 4)) != 100`: benar → ; salah → lewati langkah [kode -/3]
      - `.Policy.Payment.ListInstallment(<LAST>).InstallmentPercentage` = `@Math.divide(100, Param.Installment, 4) + (100 - Local.InstallmentPercentage)`
    - **5.2.4** `Property-Set`
      - syarat `Local.InstallmentNo<Param.Installment`: benar → ; salah → lewati langkah [kode -/3]
      - `.Policy.Payment.ListInstallment(<LAST>).InstallmentPercentage` = `@Math.divide(100, Param.Installment, 4)`
    - **5.2.5** `Property-Set`
      - `.Policy.Payment.ListInstallment(<LAST>).InstallmentNo` = `Local.InstallmentNo`
      - `.Policy.Payment.ListInstallment(<LAST>).Premium` = `@Math.divide(((.Policy.Payment.Premium)*.Policy.Payment.ListInstallment(<LAST>).InstallmentPercentage),100,4)`
      - `.Policy.Payment.ListInstallment(<LAST>).DueDate` = `@DateTime.addCalendar(@DateTime.CurrentDate("dd/MM/yyyy","GMT"),0,1,0,0,0,0,0)`
      - `.Policy.Payment.ListInstallment(<LAST>).PaymentTotal` = `@Math.divide(((.Policy.Payment.Premium)*.Policy.Payment.ListInstallment(<LAST>).InstallmentPercentage),100,4)`

### 6 · `GeneratePolicyNoTreaty_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 285.758 B

- **1** `Page-Clear-Messages` · halaman `pyWorkPage`
- **2** `Page-Remove`
  - `Page=ParamSeq`
  - `Page=InputData`
- **3** `RDB-List` · halaman `TglProd` — get tanggal produksi
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Int-policyjson`; RequestType=`GETTanggalClosing_SQL`
- **4** `Property-Set`
  - `Local.TglProd` = `TglProd.pxResults(1).TANGGAL`
- **5** `(kosong)` — SET PRODDATETIME
  - ulang: pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
  - syarat `pyWorkPage.OfferFacIn.PolicyData.PolicyNo==""`: benar → lanjut; salah → lewati langkah [kode 2/3] *(When langkah tidak dicentang)*
  - **5.1** `Property-Set`
    - syarat `@PropertyHasValue(pyWorkPage.OfferFacIn.PolicyData.ProdDateTime)`: benar → lewati langkah; salah → lanjut [kode 3/2] *(When langkah tidak dicentang)*
    - `pyWorkPage.PolicyTreatyIn.ProductionDate` = `@CurrentDateTime()`
  - **5.2** `Property-Set` — jika start date lebih besar dari currentdate
    - syarat `@toInt(@DateTime.DateTimeDifference(pyWorkPage.PolicyTreatyIn.StatementDate,@CurrentDateTime(),D))<0`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `pyWorkPage.PolicyTreatyIn.ProductionDate` = `pyWorkPage.PolicyTreatyIn.StatementDate`
  - **5.3** `Property-Set` — Jika Tanggal Proddate > tanggal closing
    - syarat `@toInt(@FormatDateTime(pyWorkPage.PolicyTreatyIn.ProductionDate,"dd","Asia/Jakarta","in_ID"))>Local.TglProd`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `pyWorkPage.PolicyTreatyIn.ProductionDate` = `@FormatDateTime(pyWorkPage.PolicyTreatyIn.ProductionDate,"yyyyMM","Asia/Jakarta","in_ID")+"01T050000.000 GMT"`
    - `pyWorkPage.PolicyTreatyIn.ProductionDate` = `@addCalendar(pyWorkPage.PolicyTreatyIn.ProductionDate,0,1,0,0,0,0,0)`
- **6** `Property-Set` — set error messagee
  - `Local.errmsg` = `"Unable to generate policy no, please click button submit to proceed"`
  - `Param.errmsg` = `"Error in generating No Polis, Please contact IT."`
- **7** `Property-Set` — when premium (QR)
  - syarat `.DueTo==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `InputData.CARI20` = `"QR"`
- **8** `Property-Set` — when pay (QP)
  - syarat `.DueTo==0`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `InputData.CARI20` = `"QP"`
- **9** `Property-Set` — pyWorkPage.PolicyTreatyIn.ClaimType=="XOL Retro"
  - syarat `pyWorkPage.PolicyTreatyIn.ClaimType=="XOL Retro"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `InputData.CARI20` = `"TP"`
- **10** `Page-Copy` — copy quotation to inside the policytreatyin
  - parameter: CopyFrom=`pyWorkPage.Quotation`; CopyInto=`pyWorkPage.PolicyTreatyIn.QuotationData`
- **11** `Call FetchTreatyGroupOldID` · halaman `pyWorkPage`
  - syarat `pyWorkPage.PolicyTreatyIn.TreatyGroupOldID==""`: benar → ; salah → lewati langkah [kode -/3]
- **12** `Property-Set` · halaman `pyWorkPage.PolicyTreatyIn`
  - `Param.ID` = `.OJKBusinessID`
- **13** `Call FetchTreatyGroupOJK` · halaman `pyWorkPage`
  - syarat `pyWorkPage.PolicyTreatyIn.OJKBusinessID==""`: benar → ; salah → lewati langkah [kode -/3]
- **14** `RDB-List` · label `//` — Fetch Current Date
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; ClassName=`ASM-FW-GISFW-Work`; Access=`RNM`; BrowsePage=`TempDate`; RequestType=`GetSQLDate`
- **15** `Page-Set-Messages` · label `//` — When the fetched date format is dd-MMM-yyyy, (01-JAN-2020) set error message then abort the activity
  - syarat `@inString(TempDate.pxResults(1).CARI1,"-")>=0`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Category=`pyMessageLabel`; ClassOfPage=`ASM-FW-GISFW-Work-NB`; Message=`Param.errmsg`; Page=`pyWorkPage`; ContainingClassOfProperty=`@baseclass`
- **16** `Property-Set` · label `//`
  - `TempDate.pxResults(1).CARI1` = `@addCalendar(@CurrentDateTime(),'0','0','0','0','7','0','0')`
  - `local.CurrDate` = `@substring(TempDate.pxResults(1).CARI1,6,8)`
  - `TempDate.pxResults(1).CARI2` = `local.CurrDate`
- **17** `Property-Set` · label `//` — when not above 25 in that month
  - syarat `local.CurrDate>Local.TglProd`: benar → lewati langkah; salah →  [kode 3/-]
  - `InputData.CARI23` = `TempDate.pxResults(1).CARI1`
  - `InputData.CARI24` = `@substring(InputData.CARI23,4,6) + "." + @substring(InputData.CARI23,0,4)`
- **18** `Property-Set` · label `//` — Determine if it's past 25 on current month, if true move to next month
  - syarat `local.CurrDate>Local.TglProd`: benar → ; salah → lewati langkah [kode -/3]
  - `InputData.CARI23` = `@addCalendar(TempDate.pxResults(1).CARI1,'0','1','0','0','0','0','0')`
  - `InputData.CARI24` = `@substring(InputData.CARI23,4,6) + "." + @substring(InputData.CARI23,0,4)`
- **19** `Property-Set` · label `//`
  - `InputData.CARI24` = `@FormatDateTime(pyWorkPage.PolicyTreatyIn.ProductionDate,"MM","Asia/Jakarta","in_ID")+"."+@FormatDateTime(pyWorkPage.PolicyTreatyIn.ProductionDate,"yyyy","Asia/Jakarta","in_ID")`
  - `TempDate.pxResults(1).CARI3` = `InputData.CARI24`
- **20** `RDB-List` · label `//` — Generate nopolis using OldBusinessID
  - syarat `pyWorkPage.PolicyTreatyIn.PolicyNo==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-POLISTREATYIN`; BrowsePage=`PolicyTreatyIn`; RequestType=`GenerateNoPolicyTreaty`
- **21** `RDB-List` · label `//` — Generate nopolis using OldTreatyGroupID
  - syarat `pyWorkPage.PolicyTreatyIn.PolicyNo==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Int-POLISTREATYIN`; BrowsePage=`PolicyTreatyIn`; RequestType=`GenerateNoPolicyTreatyGroup`
- **22** `RDB-List` · label `//` — Generate nopolis using OJK ID
  - syarat `pyWorkPage.PolicyTreatyIn.PolicyNo==""`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Int-POLISTREATYIN`; BrowsePage=`PolicyTreatyIn`; RequestType=`GenerateNoPolicyOJKID`
- **23** `Property-Set` · halaman `pyWorkPage.PolicyTreatyIn` · label `//` — Set flag new nopolis
  - `.IsOJKNopolis` = `"1"`
- **24** `(kosong)` · halaman `PolicyTreatyIn.pxResults` · label `//`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **24.1** `Property-Set`
    - `InputData.CARI2` = `.HASIL`
- **25** `RDB-List` — AMBIL KODE PROD
  - syarat `pyWorkPage.PolicyTreatyIn.PolicyNo==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Int-policyjson`; RequestType=`GetKodeProdNonLife_SQL`
- **26** `Property-Set`
  - syarat `pyWorkPage.PolicyTreatyIn.PolicyNo==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `ParamSeq.CARI1` = `pyWorkPage.pxObjClass`
  - `ParamSeq.CARI2` = `ParamSeq.HASIL3+"QR/QP/TP"`
  - `ParamSeq.CARI3` = `@FormatDateTime(pyWorkPage.PolicyTreatyIn.ProductionDate,"dd/MM/yyyy","Asia/Jakarta","in_ID")`
  - `ParamSeq.CARI4` = `ParamSeq.HASIL3+InputData.CARI20`
- **27** `RDB-List` — generate MM.YYYY DAN SEQUENCE
  - syarat `pyWorkPage.PolicyTreatyIn.PolicyNo==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Int-policyjson`; Access=`RNM`; RequestType=`GetSequenceNumber_SQL`
- **28** `Property-Set`
  - syarat `pyWorkPage.PolicyTreatyIn.PolicyNo==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `pyWorkPage.PolicyTreatyIn.PolicyNo` = `ParamSeq.CARI4+".T"+pyWorkPage.PolicyTreatyIn.OJKBusinessID+"."+ParamSeq.HASIL1+"."+ParamSeq.HASIL2`
- **29** `Property-Set-Messages` — Set error blm generate nopolis
  - syarat `pyWorkPage.PolicyTreatyIn.PolicyNo==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: Field=`pyWorkPage.Test`; Category=`pyMessageLabel`; Message=`Local.errmsg`; ContainingClassOfProperty=`ASM-FW-GISFW-Work-NB`
- **30** `Obj-Save` · halaman `pyWorkPage`
  - parameter: WithErrors=`true`; OnlyIfNew=`false`; WriteNow=`true`

### 6 · `InputParamUploadReas_act`

Kelas `ASM-FW-GISFW-Int-OFFERJSON` · 85.872 B · parameter: `pzInsKey`, `category`, `GCNMCategory`, `GCNMType`

- **1** `Page-New` · halaman `AttachList`
- **2** `Page-New` · halaman `AttachShowList`
- **3** `Property-Set`
  - `TempInputParam.pyCategory` = `Param.category`
  - `TempInputParam.pzInsKey` = `Param.pzInsKey`
  - `TempInputParam.PNOTE` = `Param.GCNMCategory`
- **4** `(kosong)` — When Param.pzInsKey=="" exit activity
  - syarat `Param.pzInsKey==""`: benar → keluar activity; salah →  [kode 6/-]
- **5** `Obj-Browse`
  - `Condition==; Field=.IDPEGA; Value=Param.pzInsKey; Select=true`
  - `Condition==; Field=.KATEGORI_2; Value=Param.GCNMCategory; Select=true`
  - parameter: ObjClass=`ASM-FW-GISFW-Int-DOCUMENT_POLIS`; RowKey=`ID`; PageName=`AttachmentList`; GetRowKey=`true`; UseLightWeightList=`true`
- **6** `(kosong)` · halaman `AttachmentList.pxResults`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `.PNOTE==SearchReinsurer.CARI1`: benar → ; salah → lewati langkah [kode -/3] *(When langkah tidak dicentang)*
  - **6.1** `Page-Copy`
    - syarat `Param.GCNMCategory==.KATEGORI_2`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - parameter: CopyFrom=`AttachmentList.pxResults(<CURRENT>)`; CopyInto=`AttachShowList.pxResults(<APPEND>)`

### 6 · `InputPolicyTreatyEDMDetail_NP`

Kelas `ASM-FW-GISFW-Work` · 168.677 B

- **1** `(kosong)` · halaman `pyWorkPage.TreatyIn.Installment` · label `//` — Expand All Installments (adjust premi nggak pake ini)
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **1.1** `Property-Set`
    - `.pyExpanded` = `true`
- **2** `Property-Set` — Non Nested Value Set always DUE TO US
  - `pyWorkPage.PolicyTreatyIn.DueTo` = `"1"`
  - `pyWorkPage.PolicyTreatyIn.StartDate` = `pyWorkPage.TreatyIn.Commencement`
  - `pyWorkPage.PolicyTreatyIn.EndDate` = `pyWorkPage.TreatyIn.Termination`
- **3** `Property-Set` — pyWorkPage.OfferFacIn.IsFacRetro = 1 when pyWorkPage.TreatyIn.FacultativeShare > 0 (This flag is to ease Arasapas team on spreading)
  - syarat `pyWorkPage.TreatyIn.FacultativeShare>0`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.OfferFacIn.IsFacRetro` = `"1"`
- **4** `Property-Set` — init total value (copy currency to local)
  - `local.currency` = `pyWorkPage.PolicyTreatyIn.Currency`
- **5** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share` — Total all share value
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **5.1** `(kosong)` · halaman `.GrossPremiumList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **5.1.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.grosspremi` = `local.grosspremi + .Value`
  - **5.2** `(kosong)` · halaman `.DeductionTotalList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **5.2.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.deduction` = `local.deduction + .Value`
  - **5.3** `(kosong)` · halaman `.NetPremiumList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **5.3.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.netpremi` = `local.netpremi + .Value`
  - **5.4** `(kosong)` · halaman `.RnmLimitList` · label `//` — There is no limits in adj
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **5.4.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.share` = `local.share + .Value`
- **6** `Property-Set` — Set all share value
  - `pyWorkPage.PolicyTreatyIn.PremiOgp` = `local.grosspremi`
  - `pyWorkPage.PolicyTreatyIn.ShareValue` = `local.share`
  - `pyWorkPage.PolicyTreatyIn.Deduction1` = `local.deduction`
  - `pyWorkPage.PolicyTreatyIn.Deduction2` = `0`
  - `pyWorkPage.PolicyTreatyIn.NetPremium` = `local.netpremi`
  - `pyWorkPage.PolicyTreatyIn.BalanceDueTo` = `local.netpremi`
  - `pyWorkPage.PolicyTreatyIn.ListInstallment` = `""`
- **7** `(kosong)` · halaman `pyWorkPage.TreatyIn.Installment` · label `//` — set installment
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **7.1** `Property-Set` — init fetch currency ID
    - `InputDataCredit.CARI8` = `.Currency`
  - **7.2** `RDB-List` — fetch currency ID
    - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Int-CURRENCY`; Access=`ASM`; BrowsePage=`CurrencySearch`; RequestType=`GetDataCurrencyByName_SQL`
  - **7.3** `Property-Set`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).PaymentTotal` = `.AmountTotal`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).InstallmentPercentage` = `.PctTotal`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).Currency` = `.Currency`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
    - `local.subscript` = `.pxListSubscript`
  - **7.4** `(kosong)` · halaman `.InstallmentList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3] *(When langkah tidak dicentang)*
    - **7.4.1** `Property-Set`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Premium` = `.Amount`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Premium` = `.Amount`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Currency` = `.Currency`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).InstallmentPercentage` = `.InstallmentPct`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).InstallmentNo` = `.Installment`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).DueDate` = `.PaymentDate`
- **8** `Call InsertToTreatyXOLList` · halaman `pyWorkPage` — Insert Values to Treaty XOL List
  - syarat `pyWorkPage.PolicyTreatyIn.ClaimType=="XOL Retro"`: benar → lewati langkah; salah → lanjut [kode 3/2]
- **9** `Call InsertToTreatyOutXOLList` · halaman `pyWorkPage` — Insert Values to Treaty Out XOL List
  - syarat `pyWorkPage.PolicyTreatyIn.ClaimType=="XOL Retro"`: benar → lanjut; salah → lewati langkah [kode 2/3]

### 6 · `InputPolicyTreatyInDetail_NonProp`

Kelas `ASM-FW-GISFW-Work` · 330.784 B

- **1** `Page-Clear-Messages` · halaman `pyWorkPage` — Clear spreading error message
- **2** `Page-New` · halaman `TempResult` — init
- **3** `Page-Remove` · halaman `TreatyIn` — init
- **4** `Page-New` · halaman `TreatyIn` — init
  - parameter: NewClass=`ASM-FW-GISFW-Int-TREATY_IN`
- **5** `RDB-List` · halaman `TempResult` — Run Querry
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Data-PolicyTreatyIn`; Access=`RNM`; BrowsePage=`TempResult`; RequestType=`BrowseTreatyInJoinEDM`
- **6** `Java` — Map to Page
  - Java:

```java
try{
  ClipboardPage tempPage2 = tools.findPage("pyWorkPage.TreatyIn");
  ClipboardPage DataJSON = tools.findPage("TempResult");
  ClipboardProperty DataJSONList  = DataJSON.getProperty(".pxResults");
  java.util.Iterator DataJSONListIter = DataJSONList.iterator();
  while (DataJSONListIter.hasNext())
  {
    ClipboardProperty DataJSONData = (ClipboardProperty)DataJSONListIter.next();
    ClipboardPage DataJSONPage = DataJSONData.getPageValue();
    String IsiDataJson = DataJSONPage.getString("CLASSOFBUSINESS");
    tempPage2.adoptJSONObject(IsiDataJson);
  }
} catch(InvalidStreamError e){
  oLog.error("ReloadSection:Invalid JSON Stream for data page params : "+e.getMessage());
} catch(Exception e){
  oLog.error("ReloadSection:Expection : "+e.getMessage());
}
```

- **7** `Property-Set`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `pyWorkPage.TreatyIn.TotalShareNetNP` = `pyWorkPage.TreatyIn.TotalFacShareDeductionNP`
  - `pyWorkPage.TreatyIn.TotalShareGrossNP` = `pyWorkPage.TreatyIn.TotalShareNetNP`
  - `pyWorkPage.TreatyIn.LimitShareSummaryList` = `pyWorkPage.TreatyIn.LimitFacShareSummaryList`
- **8** `(kosong)` · halaman `pyWorkPage.TreatyIn.LimitShareSummaryList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **8.1** `Property-Set`
    - `.NetPremi` = `.Deductible`
    - `.NetPremi2` = `.Deductible2`
    - `.Deductible` = `0`
    - `.Deductible2` = `0`
- **9** `Property-Set`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `Local.NetValue` = `pyWorkPage.TreatyIn.TotalShareNetNP(1).Value`
  - `Local.isOR` = `@substring(pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(1).Pct,4,6)`
  - `Local.isRI` = `@substring(pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(2).Pct,4,7)`
  - `pyWorkPage.TreatyIn.TotalSpreadedNetPremi(1).Value` = `@if(Local.isOR=="OR",Local.NetValue*pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(1).Pct/100,Local.NetValue*pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(2).Pct/100)`
  - `pyWorkPage.TreatyIn.TotalSpreadedNetPremiRI(1).Value` = `@if(Local.isOR=="R/I",Local.NetValue*pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(2).Pct/100,Local.NetValue*pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(1).Pct/100)`
- **10** `Property-Set` — When input using master EDM, init
  - syarat `TreatyMasterInEDM`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.IsEDMInputOnNB` = `true`
- **11** `Call InputPolicyTreatyEDMDetail_NP` — When EDM, Run InputPolicyTreatyEDMDetail_NP instead, skipping the below steps
  - syarat `TreatyMasterInEDM`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - syarat `pyWorkPage.PolicyTreatyIn.IsEDMInputOnNB==true`: benar → lewati langkah; salah →  [kode 3/-]
  - transisi `ALWAYS`: benar → keluar activity; salah → lanjut [kode 6/2]
- **12** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.ListInstallment` — Expand All Installments
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **12.1** `Property-Set`
    - `.pyExpanded` = `true`
- **13** `Property-Set` — Set always DUE TO US
  - `pyWorkPage.PolicyTreatyIn.DueTo` = `"1"`
  - `pyWorkPage.PolicyTreatyIn.StartDate` = `pyWorkPage.TreatyIn.Commencement`
  - `pyWorkPage.PolicyTreatyIn.EndDate` = `pyWorkPage.TreatyIn.Termination`
- **14** `Property-Set` — pyWorkPage.OfferFacIn.IsFacRetro = 1 when pyWorkPage.TreatyIn.FacultativeShare > 0 (This flag is to ease Arasapas team on spreading)
  - syarat `pyWorkPage.TreatyIn.FacultativeShare>0`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.OfferFacIn.IsFacRetro` = `"1"`
- **15** `Property-Set` — init total value (copy currency to local)
  - `local.currency` = `pyWorkPage.PolicyTreatyIn.Currency`
- **16** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share` — Total all share value
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **16.1** `(kosong)` · halaman `.GrossPremiumList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **16.1.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.grosspremi` = `local.grosspremi + .Value`
  - **16.2** `(kosong)` · halaman `.DeductionTotalList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **16.2.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.deduction` = `local.deduction + .Value`
    - **16.2.2** `Property-Set` — pyWorkPage.PolicyTreatyIn.FlagPPH=="true"
      - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `Local.BrokerageFeeSebenarnya` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(local.deduction,@divide(102.2,100,4),4),local.deduction)`
      - `Local.PPNValue` = `Local.BrokerageFeeSebenarnya* @divide(2.2,100,4)`
      - `Local.PPHValue` = `Local.BrokerageFeeSebenarnya* @divide(2,100,4)`
  - **16.3** `(kosong)` · halaman `.NetPremiumList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **16.3.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.netpremi` = `local.netpremi + .Value`
  - **16.4** `(kosong)` · halaman `.RnmLimitList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **16.4.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.share` = `local.share + .Value`
- **17** `(kosong)` · halaman `pyWorkPage.TreatyIn.FacultativeShareList` — pyWorkPage.TreatyIn.FacultativeShareList
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **17.1** `(kosong)` · halaman `.DeductionTotalList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - syarat `.Currency==local.currency`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - **17.1.1** `Property-Set`
      - `local.deduction` = `local.deduction + .Value`
- **18** `Property-Set` — Set all share value
  - `pyWorkPage.PolicyTreatyIn.PremiOgp` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction,local.grosspremi)`
  - `pyWorkPage.PolicyTreatyIn.ShareValue` = `local.share`
  - `pyWorkPage.PolicyTreatyIn.Deduction1` = `local.deduction`
  - `pyWorkPage.PolicyTreatyIn.Deduction2` = `0`
  - `pyWorkPage.PolicyTreatyIn.NetPremium` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction,local.netpremi)`
  - `pyWorkPage.PolicyTreatyIn.BalanceDueTo` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction,local.netpremi)`
  - `pyWorkPage.PolicyTreatyIn.ListInstallment` = `""`
- **19** `Property-Set` — pyWorkPage.PolicyTreatyIn.FlagPPH=="true"
  - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `pyWorkPage.PolicyTreatyIn.PPNValue` = `Local.PPNValue`
  - `pyWorkPage.PolicyTreatyIn.PPHValue` = `Local.PPHValue`
  - `pyWorkPage.PolicyTreatyIn.BalanceDueTo` = `local.netpremi+Local.PPNValue+Local.PPHValue`
  - `pyWorkPage.PolicyTreatyIn.BalanceBeforePPH` = `local.netpremi+Local.PPNValue`
  - `pyWorkPage.PolicyTreatyIn.BalanceBeforeTax` = `local.netpremi`
- **20** `(kosong)` · halaman `pyWorkPage.TreatyIn.Installment` — set installment
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **20.1** `Property-Set` — init fetch currency ID
    - `InputDataCredit.CARI8` = `.Currency`
  - **20.2** `RDB-List` — fetch currency ID
    - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-CURRENCY`; BrowsePage=`CurrencySearch`; RequestType=`GetDataCurrencyByName_SQL`
  - **20.3** `Property-Set`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).PaymentTotal` = `.AmountTotal`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).InstallmentPercentage` = `.PctTotal`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).Currency` = `.Currency`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
    - `local.subscript` = `.pxListSubscript`
  - **20.4** `(kosong)` · halaman `.InstallmentList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3] *(When langkah tidak dicentang)*
    - **20.4.1** `Property-Set`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Premium` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction*@divide(.InstallmentPct,100,10),.Amount)`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Premium` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction*@divide(.InstallmentPct,100,10),.Amount)`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Currency` = `.Currency`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).InstallmentPercentage` = `.InstallmentPct`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).InstallmentNo` = `.Installment`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).DueDate` = `.PaymentDate`
- **21** `Call InsertToTreatyXOLListRetroShare` — Insert Values to Treaty XOL List Kalau ada Retro Share
  - syarat `pyWorkPage.TreatyIn.FacultativeShare != 0`: benar → ; salah → lewati langkah [kode -/3]
- **22** `Call InsertToTreatyXOLList` — Insert Values to Treaty XOL List
  - syarat `pyWorkPage.TreatyIn.FacultativeShare==0`: benar → ; salah → lewati langkah [kode -/3]
- **23** `Call TreatyNonPropSetSpreading` — Set Spreading For Non Prop (both with Additional Retro or with only one retro)
- **24** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.TreatyXOLDifferenceList` · label `//` — Set spreading based on master data. Loop per currency (not sure if this is getting used or not, but the older step (step 14 is alright for single currency.)
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **24.1** `Property-Set`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<APPEND>).TreatyType` = `pyWorkPage.TreatyIn.Share(1).SpreadingTypeIDXOL`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<APPEND>).Currency` = `.Currency`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<APPEND>).CurrencyID` = `.IDCurrency`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).PremiumSpreaded` = `.NetPremi`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).SharePercentage` = `"100"`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
    - `local.sprderror` = `"Spreading in Master data is incomplete"`

### 6 · `InputPolicyTreatyInDetail_preACT`

Kelas `ASM-FW-GISFW-Work` · 488.903 B · parameter: `SOB`, `CLASSOFBUSINESSID`, `TREATYYEAR`, `LIMITCURRENCY`, `CEDING`, `TREATYGROUP`, `TREATYTYPE`, `ID`

- **1** `Property-Set` — RD BrowseTreatyInDetail
  - `Param.pyReportName` = `"BrowseTreatyJoinEDM"`
  - `Param.pyReportClass` = `"ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM"`
- **2** `Call pxRetrieveReportData` · halaman `ReportPage`
  - parameter: pyUseAlternateDB=`false`; pySkipSummaryProcessing=`false`
- **3** `Property-Set` — Set property
  - `pyWorkPage.PolicyTreatyIn.IDCurrency` = `pyReportContentPage.pxResults(1).CURRENCYID`
  - `pyWorkPage.PolicyTreatyIn.Currency` = `pyReportContentPage.pxResults(1).LIMITCURRENCY`
  - `pyWorkPage.PolicyTreatyIn.TreatyYear` = `pyReportContentPage.pxResults(1).TREATYYEAR`
  - `pyWorkPage.PolicyTreatyIn.TreatyGroupID` = `pyReportContentPage.pxResults(1).TREATYGROUPID`
  - `pyWorkPage.PolicyTreatyIn.TreatyGroupName` = `pyReportContentPage.pxResults(1).TREATYGROUP`
  - `pyWorkPage.PolicyTreatyIn.BizName` = `pyReportContentPage.pxResults(1).CLASSOFBUSINESS`
  - `pyWorkPage.PolicyTreatyIn.BizCode` = `pyReportContentPage.pxResults(1).CLASSOFBUSINESSID`
  - `pyWorkPage.Quotation.BusinessName` = `pyReportContentPage.pxResults(1).CLASSOFBUSINESS`
  - `pyWorkPage.PolicyTreatyIn.SOBName` = `pyReportContentPage.pxResults(1).SOB`
  - `pyWorkPage.PolicyTreatyIn.SOB` = `pyReportContentPage.pxResults(1).SOBID`
  - `pyWorkPage.PolicyTreatyIn.CedingCoName` = `pyReportContentPage.pxResults(1).CEDING`
  - `pyWorkPage.PolicyTreatyIn.CedingCo` = `pyReportContentPage.pxResults(1).CEDINGID`
  - `pyWorkPage.PolicyTreatyIn.InsuredID` = `pyWorkPage.Quotation.InsuredID`
  - `pyWorkPage.PolicyTreatyIn.InsuredName` = `pyWorkPage.Quotation.InsuredName`
  - `pyWorkPage.PolicyTreatyIn.NoOffer` = `pyReportContentPage.pxResults(1).TREATYID`
  - `pyWorkPage.PolicyTreatyIn.PremiOgp` = `pyReportContentPage.pxResults(1).MDPVALUE`
  - `pyWorkPage.PolicyTreatyIn.Installment` = `pyReportContentPage.pxResults(1).INSTALLMENTNO`
  - `param.Installment` = `pyReportContentPage.pxResults(1).INSTALLMENTNO`
  - `pyWorkPage.PolicyTreatyIn.Deduction1` = `pyReportContentPage.pxResults(1).DEDUCTION1`
  - `pyWorkPage.PolicyTreatyIn.Deduction2` = `pyReportContentPage.pxResults(1).DEDUCTION2`
  - `pyWorkPage.PolicyTreatyIn.TreatyType` = `pyReportContentPage.pxResults(1).TREATYTYPE`
  - `pyWorkPage.PolicyTreatyIn.LayerType` = `pyReportContentPage.pxResults(1).LAYERTYPE`
  - `pyWorkPage.PolicyTreatyIn.Layer` = `pyReportContentPage.pxResults(1).LAYER`
  - `pyWorkPage.PolicyTreatyIn.LayerPartType` = `pyReportContentPage.pxResults(1).LAYERPARTTYPE`
  - `pyWorkPage.PolicyTreatyIn.LayerPart` = `pyReportContentPage.pxResults(1).LAYERPART`
  - `pyWorkPage.PolicyTreatyIn.ShareCurrency` = `pyReportContentPage.pxResults(1).SHARECURRENCY`
  - `pyWorkPage.PolicyTreatyIn.ShareValue` = `pyReportContentPage.pxResults(1).SHAREVALUE`
  - `pyWorkPage.Quotation.ProportionalType` = `pyReportContentPage.pxResults(1).PROPORTIONTYPE`
  - `pyWorkPage.Quotation.SourceOfBusiness` = `pyReportContentPage.pxResults(1).SOBID`
  - `pyWorkPage.Quotation.SobName` = `pyReportContentPage.pxResults(1).SOB`
  - `pyWorkPage.Quotation.CedingCoName` = `pyReportContentPage.pxResults(1).CEDING`
  - `pyWorkPage.Quotation.CedingCo` = `pyReportContentPage.pxResults(1).CEDINGID`
  - `pyWorkPage.Quotation.InsuredID` = `pyWorkPage.Quotation.InsuredID`
  - `pyWorkPage.Quotation.InsuredName` = `pyWorkPage.Quotation.InsuredName`
- **4** `Call SetTreatyCurrencyID` · halaman `pyWorkPage` — When Currency ID is null, set currency ID with this activity
  - syarat `pyWorkPage.PolicyTreatyIn.IDCurrency==""`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Currency=`pyWorkPage.PolicyTreatyIn.Currency`
- **5** `Property-Set` — when pyReportContentPage.pxResults(1).TREATYTYPE == ""
  - syarat `pyReportContentPage.pxResults(1).TREATYTYPE == ""`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.TreatyType` = `"XOL"`
- **6** `Property-Set` — Copy ID, set treatygroupid to parameter page
  - `TreatyIn.ID` = `pyReportContentPage.pxResults(1).ID`
  - `Param.ID` = `.PolicyTreatyIn.TreatyGroupID`
- **7** `Call FetchTreatyGroupOldID` · halaman `pyWorkPage` · label `//` — Fetch OLDID into PolicyTreatyIn.TreatyGroupOldID
- **8** `Call FetchTreatyGroupOJK` · halaman `pyWorkPage` — Fetch Treaty Group OJK
- **9** `RDB-List` · halaman `TempResult` — SQL BrowseTreatyInDetailJoinEDM
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM`; BrowsePage=`TempResult`; RequestType=`BrowseTreatyInDetailJoinEDM`
- **10** `Java` — Mapping
  - Java:

```java
try{
  ClipboardPage tempPage2 = tools.findPage("TreatyIn");
  ClipboardPage DataJSON = tools.findPage("TempResult");
  ClipboardProperty DataJSONList  = DataJSON.getProperty(".pxResults");
  java.util.Iterator DataJSONListIter = DataJSONList.iterator();
  while (DataJSONListIter.hasNext())
  {
    ClipboardProperty DataJSONData = (ClipboardProperty)DataJSONListIter.next();
    ClipboardPage DataJSONPage = DataJSONData.getPageValue();
    String IsiDataJson = DataJSONPage.getString("CLASSOFBUSINESS");
    tempPage2.adoptJSONObject(IsiDataJson);
  }
} catch(InvalidStreamError e){
  oLog.error("ReloadSection:Invalid JSON Stream for data page params : "+e.getMessage());
} catch(Exception e){
  oLog.error("ReloadSection:Expection : "+e.getMessage());
}
```

- **11** `Property-Remove` — Delete old installmentList
  - `Property=pyWorkPage.PolicyTreatyIn.ListInstallment`
- **12** `call CountResult1_Act` · halaman `pyWorkPage.PolicyTreatyIn` · label `//`
- **13** `(kosong)` · halaman `TreatyIn.INSTALLMENT` — Insert new installmentList
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **13.1** `(kosong)` · halaman `.InstallmentList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - syarat `pyWorkPage.PolicyTreatyIn.Currency==.Currency`: benar → ; salah → lewati langkah [kode -/3]
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<APPEND>).InstallmentNo` = `.AcceptStatus`
    - **13.1.1** `Property-Set`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<APPEND>).InstallmentNo` = `.Installment`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).DueDate` = `.PaymentDate`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).InstallmentPercentage` = `.InstallmentPct`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).Premium` = `@Math.divide(((pyWorkPage.PolicyTreatyIn.BalanceDueTo)*.InstallmentPct),100,4)`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).PaymentTotal` = `@Math.divide(((pyWorkPage.PolicyTreatyIn.BalanceDueTo)*.InstallmentPct),100,4)`
    - **13.1.2** `Property-Set` · label `//` — pyWorkPage.PolicyTreatyIn.FlagPPH=="true"
      - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).Premium` = `@Math.divide(((pyWorkPage.PolicyTreatyIn.BalanceDueTo)*.InstallmentPct),100,4)`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).PaymentTotal` = `@Math.divide(((pyWorkPage.PolicyTreatyIn.BalanceDueTo)*.InstallmentPct),100,4)`
- **14** `(kosong)` — Set Bisnis & Insured Quotation
  - ulang: pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefLimit=`1`, pyStepsRepeatDefStart=`1`
  - **14.1** `Property-Set` — Copy InsuredName
    - `SearchClient.CARI1` = `pyWorkPage.Quotation.InsuredName`
  - **14.2** `RDB-List` — SQL GetClientID_SQL
    - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; ClassName=`ASM-FW-GISFW-Int-CLIENT`; Access=`ASM`; BrowsePage=`DataClient`; RequestType=`GetClientID_SQL`
  - **14.3** `Property-Set` — Set InsuredID
    - syarat `pyWorkPage.Quotation.InsuredID==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `pyWorkPage.Quotation.InsuredID` = `DataClient.pxResults(1).CARI1`
    - `pyWorkPage.PolicyTreatyIn.InsuredID` = `DataClient.pxResults(1).CARI1`
    - `pyWorkPage.Quotation.InsuredID` = `DataClient.pxResults(1).CARI1`
  - **14.4** `Property-Set` — Copy BusinessName
    - `OldID.CARI2` = `pyWorkPage.Quotation.BusinessName`
  - **14.5** `(kosong)` — SET BUSINESS CODE KALAU KOSONG
    - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
    - syarat `pyWorkPage.PolicyTreatyIn.BizCode==""`: benar → lanjut; salah → lewati langkah [kode 2/3] *(When langkah tidak dicentang)*
    - **14.5.1** `Property-Set`
      - `OldID.CARI2` = `pyWorkPage.Quotation.BusinessName`
    - **14.5.2** `Property-Set` — MBU
      - syarat `@contains(pyWorkPage.Quotation.BusinessName,"MBU")`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"MOTOR VEHICLE"`
    - **14.5.3** `Property-Set` — pyWorkPage.Quotation.BusinessName=="ADVANCE PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="BID OR TENDER BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PAYMENT BONDS"\|\|pyWorkPage.Quotation.Bus …
      - syarat `pyWorkPage.Quotation.BusinessName=="ADVANCE PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="BID OR TENDER BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PERFORMANCE BONDS"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `@replaceAll(pyWorkPage.Quotation.BusinessName,"S","")`
    - **14.5.4** `Property-Set` — pyWorkPage.Quotation.BusinessName=="CUSTOMS BOND">>OTHERS CUSTOMS BOND
      - syarat `pyWorkPage.Quotation.BusinessName=="CUSTOMS BOND"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"OTHERS CUSTOMS BOND"`
    - **14.5.5** `Property-Set` — pyWorkPage.Quotation.BusinessName=="ASURANSI KREDIT">>ASURANSI KREDIT (CASH LOAN)
      - syarat `pyWorkPage.Quotation.BusinessName=="ASURANSI KREDIT"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"ASURANSI KREDIT (CASH LOAN)"`
    - **14.5.6** `Property-Set` — pyWorkPage.Quotation.BusinessName=="BOILER & PRESSURE VESSEL">>BOILER & EXCAVATOR
      - syarat `pyWorkPage.Quotation.BusinessName=="BOILER & PRESSURE VESSEL"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"BOILER & EXCAVATOR"`
    - **14.5.7** `Property-Set` — pyWorkPage.Quotation.BusinessName=="GOLF INSURANCE">>HOLE IN ONE
      - syarat `pyWorkPage.Quotation.BusinessName=="GOLF INSURANCE"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"HOLE IN ONE"`
  - **14.6** `RDB-List` — SQL GetOldIDBusiness_SQL
    - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Int-OFFERJSON`; Access=`ASM`; BrowsePage=`DataClient`; RequestType=`GetOldIDBusiness_SQL`
  - **14.7** `Property-Set` — Set BusinessOldId, GroupPanel, & BusinessCode
    - `pyWorkPage.Quotation.BusinessOldId` = `DataClient.pxResults(1).CARI1`
    - `pyWorkPage.Quotation.GroupPanel` = `DataClient.pxResults(1).CARI2`
    - `pyWorkPage.Quotation.BusinessCode` = `DataClient.pxResults(1).CARI3`
    - `pyWorkPage.PolicyTreatyIn.BizCode` = `pyWorkPage.Quotation.BusinessCode`
  - **14.8** `Property-Map-DecisionTable` · halaman `pyWorkPage` — DT BusinessType_DeT
    - parameter: PropertyName=`pyWorkPage.Quotation.BusinessType`; AllowMissingProperties=`false`; DecisionTableName=`BusinessType_DeT`
  - **14.9** `Property-Set` — Set BusinessType
    - `pyWorkPage.PolicyTreatyIn.QuotationData` = `pyWorkPage.Quotation`
- **15** `Page-Remove` — Remove page pyReportContentPage
  - `Page=pyReportContentPage`
- **16** `Call InputPolicyTreatyInDetail_NonProp` · halaman `pyWorkPage` — When Non Proportional
  - syarat `pyWorkPage.Quotation.ProportionalType=="NonProportional"`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.`
- **17** `Call TreatyInputPctCommSpreading` · halaman `pyWorkPage` — When Proportional, Copy OGR & Spreading
  - syarat `pyWorkPage.Quotation.ProportionalType=="NonProportional"`: benar → lewati langkah; salah →  [kode 3/-]
  - `pyWorkPage.PolicyTreatyIn.`
- **18** `(kosong)` · halaman `pyWorkPage.TreatyIn.LimitShareSummaryList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **18.1** `Property-Set`
    - `Local.BrokerageFeeSebenarnya` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(.Deductible,@divide(102.2,100,8),8),.Deductible)`
    - `.PPNValue` = `Local.BrokerageFeeSebenarnya* @divide(2.2,100,8)`
    - `.NetPremiAfterPPN` = `.NetPremi + .PPNValue`
    - `.PPHValue` = `Local.BrokerageFeeSebenarnya* @divide(2,100,8)`
    - `.NetPremiAfterPPH` = `.NetPremi + .PPHValue`
    - `Local.BrokerageFeeSebenarnya2` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(.Deductible2,@divide(102.2,100,8),8),.Deductible2)`
    - `.PPNValue2` = `Local.BrokerageFeeSebenarnya2* @divide(2.2,100,8)`
    - `.NetPremiAfterPPN2` = `.NetPremi2 + .PPNValue2`
    - `.PPHValue2` = `Local.BrokerageFeeSebenarnya2* @divide(2,100,8)`
    - `.NetPremiAfterPPH2` = `.NetPremi2 + .PPHValue2`
    - `Local.TotalPPNValue` = `Local.TotalPPNValue + .PPNValue`
    - `Local.TotalPPNValue2` = `Local.TotalPPNValue2 + .PPNValue2`
    - `Local.TotalPPHValue` = `Local.TotalPPHValue + .PPHValue`
    - `Local.TotalPPHValue2` = `Local.TotalPPHValue2 + .PPHValue2`
    - `Local.TotalNetPremiAfterPPN` = `Local.TotalNetPremiAfterPPN + .NetPremiAfterPPN`
    - `Local.TotalNetPremiAfterPPN2` = `Local.TotalNetPremiAfterPPN2 + .NetPremiAfterPPN2`
    - `Local.TotalNetPremiAfterPPH` = `Local.TotalNetPremiAfterPPH + .NetPremiAfterPPH`
    - `Local.TotalNetPremiAfterPPH2` = `Local.TotalNetPremiAfterPPH2 + .NetPremiAfterPPH2`
    - `Local.TotalNetPremiAfterTax` = `Local.TotalNetPremiAfterTax+.NetPremiAfterPPN+  .PPHValue`
    - `Local.TotalNetPremiAfterTax2` = `Local.TotalNetPremiAfterTax2+.NetPremiAfterPPN2+  .PPHValue2`
    - `Local.curr` = `@if(.Limit>0,"IDR","")`
    - `Local.curr2` = `@if(.Limit2>0,"USD","")`
  - **18.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.TotalShareNetNP`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **18.2.1** `Property-Set` — IDR
      - syarat `.Currency==Local.curr`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `.TotalNetPremiAfterPPN` = `Local.TotalNetPremiAfterPPN`
      - `.TotalNetPremiAfterPPH` = `Local.TotalNetPremiAfterPPH`
      - `.TotalNetPremiAfterTax` = `Local.TotalNetPremiAfterTax`
    - **18.2.2** `Property-Set` — USD
      - syarat `.Currency==Local.curr2`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `.TotalNetPremiAfterPPN` = `Local.TotalNetPremiAfterPPN2`
      - `.TotalNetPremiAfterPPH` = `Local.TotalNetPremiAfterPPH2`
      - `.TotalNetPremiAfterTax` = `Local.TotalNetPremiAfterTax2`
  - **18.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.TotalShareDeductionNP` — pyWorkPage.TreatyIn.TotalShareDeductionNP
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **18.3.1** `Property-Set` — IDR
      - syarat `.Currency==Local.curr`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `.TotalPPNValue` = `Local.TotalPPNValue`
      - `.TotalPPHValue` = `Local.TotalPPHValue`
      - `Local.TotalPPNValue` = `.TotalPPNValue`
      - `Local.TotalPPHValue` = `.TotalPPHValue`
    - **18.3.2** `Property-Set` — USD
      - syarat `.Currency==Local.curr2`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `.TotalPPNValue` = `Local.TotalPPNValue2`
      - `.TotalPPHValue` = `Local.TotalPPHValue2`
      - `Local.TotalPPNValue2` = `.TotalPPNValue`
      - `Local.TotalPPHValue2` = `.TotalPPHValue`
    - **18.3.3** `Property-Set`
      - `Local.Currency` = `.Currency`
    - **18.3.4** `Property-Set` · halaman `pyWorkPage.PolicyTreatyIn.ListInstallment`
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - `Local.CurrListInstallment` = `.Currency`
      - **18.3.4.1** `Property-Set`
        - syarat `Local.Currency==.Currency`: benar → lanjut; salah → lewati langkah [kode 2/3]
        - `.PPN` = `@if(.Currency=="IDR",Local.TotalPPNValue,Local.TotalPPNValue2)`
        - `.PPh` = `@if(.Currency=="IDR",Local.TotalPPHValue,Local.TotalPPHValue2)`
        - `Local.PPNins` = `.PPN`
        - `Local.PPHins` = `.PPh`
        - `.PaymentTotalAfterPPN` = `.PaymentTotal+.PPN`
        - `.PaymentTotalAfterTax` = `@if(.Currency=="IDR",Local.TotalNetPremiAfterTax,Local.TotalNetPremiAfterTax2)`
        - `Local.PaymentTotalAfterPPN` = `.PaymentTotalAfterPPN`
        - `Local.PaymentTotalAfterTax` = `.PaymentTotalAfterTax`
      - **18.3.4.2** `(kosong)` · halaman `.InstallmentList`
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **18.3.4.2.1** `Property-Set`
          - syarat `Local.CurrListInstallment==.Currency`: benar → lanjut; salah → lewati langkah [kode 2/3]
          - `.PPN` = `Local.PPNins * @divide(@toDecimal(.InstallmentPercentage),100,8)`
          - `.PPh` = `Local.PPHins * @divide(@toDecimal(.InstallmentPercentage),100,8)`
          - `.PremiumAfterPPN` = `Local.PaymentTotalAfterPPN * @divide(@toDecimal(.InstallmentPercentage),100,8)`
          - `.PremiumAfterTax` = `Local.PaymentTotalAfterTax* @divide(@toDecimal(.InstallmentPercentage),100,8)`
- **19** `Obj-Save` · halaman `pyWorkPage` — Save
  - parameter: WithErrors=`true`; OnlyIfNew=`false`; WriteNow=`true`

### 6 · `InputPolicyTreatyInPost_Act`

Kelas `ASM-FW-GISFW-Work` · 32.110 B

- **1** `Page-Clear-Messages` · halaman `pyWorkPage`
- **2** `Call InsertHistoryAkseptasiPega`
- **3** `Call TreatyRealizationCheckDuplicate`
  - syarat `.PolicyTreatyIn.IsApproved==1`: benar → ; salah → lewati langkah [kode -/3]
- **4** `Call SaveViewSuggest`

### 6 · `InputPolicyTreatyInPre_Act`

Kelas `ASM-FW-GISFW-Work` · 181.346 B

- **1** `Call ASM-FW-GISFW-Data-OfferFacIn.SetCategoryAttach`
- **2** `(kosong)` — SET BUSINESS CODE KALAU KOSONG
  - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
  - syarat `pyWorkPage.PolicyTreatyIn.BizCode==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **2.1** `Property-Set`
    - `OldID.CARI2` = `pyWorkPage.Quotation.BusinessName`
  - **2.2** `Property-Set` — MBU
    - syarat `@contains(pyWorkPage.Quotation.BusinessName,"MBU")`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `OldID.CARI2` = `"MOTOR VEHICLE"`
  - **2.3** `Property-Set` — pyWorkPage.Quotation.BusinessName=="ADVANCE PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="BID OR TENDER BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PAYMENT BONDS"\|\|pyWorkPage.Quotation.Bus …
    - syarat `pyWorkPage.Quotation.BusinessName=="ADVANCE PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="BID OR TENDER BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PERFORMANCE BONDS"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `OldID.CARI2` = `@replaceAll(pyWorkPage.Quotation.BusinessName,"S","")`
  - **2.4** `Property-Set` — pyWorkPage.Quotation.BusinessName=="CUSTOMS BOND">>OTHERS CUSTOMS BOND
    - syarat `pyWorkPage.Quotation.BusinessName=="CUSTOMS BOND"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `OldID.CARI2` = `"OTHERS CUSTOMS BOND"`
  - **2.5** `Property-Set` — pyWorkPage.Quotation.BusinessName=="ASURANSI KREDIT">>ASURANSI KREDIT (CASH LOAN)
    - syarat `pyWorkPage.Quotation.BusinessName=="ASURANSI KREDIT"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `OldID.CARI2` = `"ASURANSI KREDIT (CASH LOAN)"`
  - **2.6** `Property-Set` — pyWorkPage.Quotation.BusinessName=="BOILER & PRESSURE VESSEL">>BOILER & EXCAVATOR
    - syarat `pyWorkPage.Quotation.BusinessName=="BOILER & PRESSURE VESSEL"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `OldID.CARI2` = `"BOILER & EXCAVATOR"`
  - **2.7** `Property-Set` — pyWorkPage.Quotation.BusinessName=="GOLF INSURANCE">>HOLE IN ONE
    - syarat `pyWorkPage.Quotation.BusinessName=="GOLF INSURANCE"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `OldID.CARI2` = `"HOLE IN ONE"`
  - **2.8** `Property-Set` — pyWorkPage.Quotation.BusinessName=="BID OR TENDER BONDS"
    - syarat `pyWorkPage.Quotation.BusinessName=="BID OR TENDER BONDS"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `OldID.CARI2` = `"BID BOND"`
  - **2.9** `RDB-List`
    - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-OFFERJSON`; BrowsePage=`DataClient`; RequestType=`GetOldIDBusiness_SQL`
  - **2.10** `Property-Set`
    - `pyWorkPage.Quotation.BusinessOldId` = `DataClient.pxResults(1).CARI1`
    - `pyWorkPage.Quotation.GroupPanel` = `DataClient.pxResults(1).CARI2`
    - `pyWorkPage.Quotation.BusinessCode` = `DataClient.pxResults(1).CARI3`
    - `pyWorkPage.PolicyTreatyIn.BizCode` = `pyWorkPage.Quotation.BusinessCode`
    - `pyWorkPage.PolicyTreatyIn.QuotationData.BusinessCode` = `pyWorkPage.Quotation.BusinessCode`
    - `pyWorkPage.PolicyTreatyIn.QuotationData.BusinessOldId` = `pyWorkPage.Quotation.BusinessOldId`
    - `pyWorkPage.PolicyTreatyIn.QuotationData.GroupPanel` = `pyWorkPage.Quotation.GroupPanel`
- **3** `RDB-List` — Get Statement Date from SQL Date only when ReasTreatyInAdmin
  - syarat `pyWorkPage.PositionNote=="ReasTreatyInAdmin"&&@SizeOfPropertyList(pyWorkPage.PolicyTreatyIn.SuggestList)==0`: benar → ; salah → lompat ke label `SkipDate` [kode -/1] *(When langkah tidak dicentang)*
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Work`; Access=`RNM`; BrowsePage=`GetDate`; RequestType=`GetSQLDate`
- **4** `Property-Set` — Set Statement Date from SQL Date only when ReasTreatyInAdmin
  - `.PolicyTreatyIn.StatementDate` = `@toDateTime(GetDate.pxResults(1).CARI1)`
  - `.PolicyTreatyIn.ProductionDate` = `@toDateTime(GetDate.pxResults(1).CARI1)`
- **5** `Page-Copy` · label `SkipDate`
  - parameter: CopyFrom=`newAssignPage`; CopyInto=`OldAssignPage`
- **6** `Page-Remove` — Init List spreading
  - `Page=ListSpreading`
  - `Page=GetSpreadName`
- **7** `RDB-List` · halaman `GetSpreadName` — Get list for spreading
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Int-REINSURANCETYPE`; BrowsePage=`GetSpreadName`; RequestType=`SelectSpreadingTreatyInProduction`
- **8** `(kosong)` · halaman `GetSpreadName.pxResults` — Set Spreading into Temp Page
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **8.1** `Property-Set`
    - `ListSpreading.pxResults(<APPEND>).TreatyType` = `.CARI1`
    - `ListSpreading.pxResults(<LAST>).TreatyName` = `.CARI2`
- **9** `Property-Set` — Set Production Date kalau diatas tanggal 25
  - syarat `@substring(pyWorkPage.PolicyTreatyIn.StatementDate,6,2)>25`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.ProductionDate` = `@addCalendar(pyWorkPage.PolicyTreatyIn.StatementDate,'0','1','0','0','0','0','0')`
  - `pyWorkPage.PolicyTreatyIn.ProductionDate` = `@substring(pyWorkPage.PolicyTreatyIn.ProductionDate,0,6) + "01" + @substring(pyWorkPage.PolicyTreatyIn.ProductionDate,8)`
- **10** `Call TreatyRealizationCheckXOLList` — Call Treaty Realization check (check TreatyXOLList)
  - syarat `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - syarat `.PolicyTreatyIn.EDMType=="3"`: benar → lewati langkah; salah →  [kode 3/-]

### 6 · `InputPolicyTreatyOutDetail_NonProp`

Kelas `ASM-FW-GISFW-Work` · 476.848 B

- **1** `Page-Clear-Messages` · halaman `pyWorkPage` — Clear spreading error message
- **2** `Page-New` · halaman `TempResult` — init
- **3** `Page-Remove` · halaman `TreatyIn` — init
- **4** `Page-New` · halaman `TreatyIn` — init
  - parameter: NewClass=`ASM-FW-GISFW-Int-TREATY_IN`
- **5** `RDB-List` · halaman `TempResult` — Run Querry
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; ClassName=`ASM-FW-GISFW-Data-PolicyTreatyIn`; Access=`RNM`; BrowsePage=`TempResult`; RequestType=`BrowseTreatyOut`
- **6** `Java` — Map to Page
  - Java:

```java
try{
  ClipboardPage tempPage2 = tools.findPage("pyWorkPage.TreatyIn");
  ClipboardPage DataJSON = tools.findPage("TempResult");
  ClipboardProperty DataJSONList  = DataJSON.getProperty(".pxResults");
  java.util.Iterator DataJSONListIter = DataJSONList.iterator();
  while (DataJSONListIter.hasNext())
  {
    ClipboardProperty DataJSONData = (ClipboardProperty)DataJSONListIter.next();
    ClipboardPage DataJSONPage = DataJSONData.getPageValue();
    String IsiDataJson = DataJSONPage.getString("CLASSOFBUSINESS");
    tempPage2.adoptJSONObject(IsiDataJson);
  }
} catch(InvalidStreamError e){
  oLog.error("ReloadSection:Invalid JSON Stream for data page params : "+e.getMessage());
} catch(Exception e){
  oLog.error("ReloadSection:Expection : "+e.getMessage());
}
```

- **7** `Property-Set`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `pyWorkPage.TreatyIn.TotalShareNetNP` = `pyWorkPage.TreatyIn.TotalFacShareDeductionNP`
  - `pyWorkPage.TreatyIn.TotalShareGrossNP` = `pyWorkPage.TreatyIn.TotalShareNetNP`
  - `pyWorkPage.TreatyIn.LimitShareSummaryList` = `pyWorkPage.TreatyIn.LimitFacShareSummaryList`
- **8** `Property-Set`
  - `Local.SOBName` = `pyWorkPage.PolicyTreatyIn.SOBName`
  - `Local.TotalSummary` = `0`
  - `Local.TotalSummary2` = `0`
  - `Local.TotalSummaryMDP` = `0`
  - `Local.TotalSummaryMDP2` = `0`
  - `Local.TotalSummaryBrokerage` = `0`
  - `Local.TotalSummaryBrokerage2` = `0`
  - `Local.TotalSummaryNetPremi` = `0`
  - `Local.TotalSummaryNetPremi2` = `0`
- **9** `(kosong)` · halaman `pyWorkPage.TreatyIn.LimitShareSummaryList`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **9.1** `Property-Set`
    - `.NetPremi` = `.Deductible`
    - `.NetPremi2` = `.Deductible2`
    - `.Deductible` = `0`
    - `.Deductible2` = `0`
- **10** `Property-Set` · halaman `pyWorkPage.TreatyIn.LimitShareSummaryList` — \\\pyWorkPage.TreatyIn.LimitShareSummaryList
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `Local.idxsummary` = `.pxListSubscript`
  - **10.1** `Property-Remove` — pyWorkPage.TreatyIn.LimitShareSummaryList(Local.idxsummary)
    - syarat `Local.SOBName==.ReinsName`: benar → lewati langkah; salah → lanjut [kode 3/2]
    - `Property=pyWorkPage.TreatyIn.LimitShareSummaryList(Local.idxsummary)`
- **11** `Property-Set` · halaman `pyWorkPage.TreatyIn.LimitShareSummaryList` — \\\pyWorkPage.TreatyIn.LimitShareSummaryList
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `Local.idxsummary` = `.pxListSubscript`
  - **11.1** `Property-Set` — pyWorkPage.TreatyIn.TotalShareRnmNP
    - `Local.TotalSummary` = `Local.TotalSummary + .Limit`
    - `Local.TotalSummary2` = `Local.TotalSummary2 + .Limit2`
    - `Local.TotalSummaryMDP` = `Local.TotalSummaryMDP +.MDP`
    - `Local.TotalSummaryMDP2` = `Local.TotalSummaryMDP2 + .MDP2`
    - `Local.TotalSummaryBrokerage` = `Local.TotalSummaryBrokerage + .Deductible`
    - `Local.TotalSummaryBrokerage2` = `Local.TotalSummaryBrokerage2 + .Deductible2`
    - `Local.TotalSummaryNetPremi` = `Local.TotalSummaryNetPremi + .NetPremi`
    - `Local.TotalSummaryNetPremi2` = `Local.TotalSummaryNetPremi2 + .NetPremi2`
- **12** `(kosong)` · halaman `pyWorkPage.TreatyIn.TotalShareRnmNP` — Total LIMIT
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **12.1** `Property-Set`
    - syarat `.Currency=="IDR"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Value` = `Local.TotalSummary`
  - **12.2** `Property-Set`
    - syarat `.Currency=="USD"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Value` = `Local.TotalSummary2`
- **13** `(kosong)` · halaman `pyWorkPage.TreatyIn.TotalShareGrossNP` — Total MDP
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **13.1** `Property-Set`
    - syarat `.Currency=="IDR"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Value` = `Local.TotalSummaryMDP`
  - **13.2** `Property-Set`
    - syarat `.Currency=="USD"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Value` = `Local.TotalSummaryMDP2`
- **14** `(kosong)` · halaman `pyWorkPage.TreatyIn.TotalShareDeductionNP` — Total BROKERAGE
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **14.1** `Property-Set`
    - syarat `.Currency=="IDR"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Value` = `Local.TotalSummaryBrokerage`
  - **14.2** `Property-Set`
    - syarat `.Currency=="USD"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Value` = `Local.TotalSummaryBrokerage2`
- **15** `(kosong)` · halaman `pyWorkPage.TreatyIn.TotalShareNetNP` — Total NETPREMI
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **15.1** `Property-Set`
    - syarat `.Currency=="IDR"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Value` = `Local.TotalSummaryNetPremi`
  - **15.2** `Property-Set`
    - syarat `.Currency=="USD"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Value` = `Local.TotalSummaryNetPremi2`
- **16** `Property-Set`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `Local.NetValue` = `pyWorkPage.TreatyIn.TotalShareNetNP(1).Value`
  - `Local.isOR` = `@substring(pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(1).Pct,4,6)`
  - `Local.isRI` = `@substring(pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(2).Pct,4,7)`
  - `pyWorkPage.TreatyIn.TotalSpreadedNetPremi(1).Value` = `@if(Local.isOR=="OR",Local.NetValue*pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(1).Pct/100,Local.NetValue*pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(2).Pct/100)`
  - `pyWorkPage.TreatyIn.TotalSpreadedNetPremiRI(1).Value` = `@if(Local.isOR=="R/I",Local.NetValue*pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(2).Pct/100,Local.NetValue*pyWorkPage.TreatyIn.Share(1).SpreadingListXOL(1).Pct/100)`
- **17** `Property-Set` — When input using master EDM, init
  - syarat `TreatyMasterInEDM`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.IsEDMInputOnNB` = `true`
- **18** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.ListInstallment` — Expand All Installments
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **18.1** `Property-Set`
    - `.pyExpanded` = `true`
- **19** `Property-Set` — Set always DUE TO US
  - `pyWorkPage.PolicyTreatyIn.DueTo` = `"1"`
  - `pyWorkPage.PolicyTreatyIn.StartDate` = `pyWorkPage.TreatyIn.Commencement`
  - `pyWorkPage.PolicyTreatyIn.EndDate` = `pyWorkPage.TreatyIn.Termination`
- **20** `Property-Set` — pyWorkPage.OfferFacIn.IsFacRetro = 1 when pyWorkPage.TreatyIn.FacultativeShare > 0 (This flag is to ease Arasapas team on spreading)
  - syarat `pyWorkPage.TreatyIn.FacultativeShare>0`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.OfferFacIn.IsFacRetro` = `"1"`
- **21** `Property-Set` — init total value (copy currency to local)
  - `local.currency` = `pyWorkPage.PolicyTreatyIn.Currency`
- **22** `Property-Set` · halaman `pyWorkPage.TreatyIn.ShareReins` — HAPUS SHARE REINS YANG BEDA SOB
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `Local.idxShareReins` = `.pxListSubscript`
  - **22.1** `Property-Set` · halaman `.ReinsuranceListTONP`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - `Local.idxReinsurer` = `.pxListSubscript`
    - **22.1.1** `Property-Remove` — pyWorkPage.TreatyIn.ShareReins(Local.idxShareReins).ReinsuranceListTONP(Local.idxReinsurer)
      - syarat `Local.SOBName==.ReinsName`: benar → lewati langkah; salah → lanjut [kode 3/2]
      - `Property=pyWorkPage.TreatyIn.ShareReins(Local.idxShareReins).ReinsuranceListTONP(Local.idxReinsurer)`
- **23** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins` — Total all share value
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **23.1** `Property-Set` · halaman `.ReinsuranceListTONP` — spreading
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - `Local.spreading` = `.SpreadingTypeIDXOL`
    - **23.1.1** `(kosong)` · halaman `.GrossPremiumList`
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **23.1.1.1** `Property-Set`
        - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
        - `local.grosspremi` = `local.grosspremi + .Value`
    - **23.1.2** `(kosong)` · halaman `.DeductionTotalList`
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **23.1.2.1** `Property-Set`
        - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
        - `local.deduction` = `local.deduction + .Value`
      - **23.1.2.2** `Property-Set` — pyWorkPage.PolicyTreatyIn.FlagPPH=="true"
        - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
        - `Local.BrokerageFeeSebenarnya` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(local.deduction,@divide(102.2,100,4),4),local.deduction)`
        - `Local.PPNValue` = `Local.BrokerageFeeSebenarnya* @divide(2.2,100,4)`
        - `Local.PPHValue` = `Local.BrokerageFeeSebenarnya* @divide(2,100,4)`
    - **23.1.3** `(kosong)` · halaman `.NetPremiumList`
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **23.1.3.1** `Property-Set`
        - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
        - `local.netpremi` = `local.netpremi + .Value`
    - **23.1.4** `(kosong)` · halaman `.RnmLimitList`
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **23.1.4.1** `Property-Set`
        - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
        - `local.share` = `local.share + .Value`
- **24** `(kosong)` · halaman `pyWorkPage.TreatyIn.SpreadingTONP(1).ReinsuranceListTONP(1).LayerList` · label `//` — Total all share value
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **24.1** `(kosong)` · halaman `.GrossPremiumList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **24.1.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.grosspremi` = `local.grosspremi + .Value`
  - **24.2** `(kosong)` · halaman `.DeductionTotalList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **24.2.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.deduction` = `local.deduction + .Value`
    - **24.2.2** `Property-Set` — pyWorkPage.PolicyTreatyIn.FlagPPH=="true"
      - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `Local.BrokerageFeeSebenarnya` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(local.deduction,@divide(102.2,100,4),4),local.deduction)`
      - `Local.PPNValue` = `Local.BrokerageFeeSebenarnya* @divide(2.2,100,4)`
      - `Local.PPHValue` = `Local.BrokerageFeeSebenarnya* @divide(2,100,4)`
  - **24.3** `(kosong)` · halaman `.NetPremiumList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **24.3.1** `Property-Set`
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.netpremi` = `local.netpremi + .Value`
  - **24.4** `(kosong)` · halaman `.RnmLimitList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **24.4.1** `Property-Set` — pyWorkPage.TreatyIn.Installment
      - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3]
      - `local.share` = `local.share + .Value`
- **25** `Property-Set` — Set all share value
  - `pyWorkPage.PolicyTreatyIn.PremiOgp` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction,local.grosspremi)`
  - `pyWorkPage.PolicyTreatyIn.ShareValue` = `local.share`
  - `pyWorkPage.PolicyTreatyIn.Deduction1` = `local.deduction`
  - `pyWorkPage.PolicyTreatyIn.Deduction2` = `0`
  - `pyWorkPage.PolicyTreatyIn.NetPremium` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction,local.netpremi)`
  - `pyWorkPage.PolicyTreatyIn.BalanceDueTo` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction,local.netpremi)`
  - `pyWorkPage.PolicyTreatyIn.ListInstallment` = `""`
- **26** `(kosong)` · halaman `pyWorkPage.TreatyIn.Installment` — set installment
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **26.1** `Property-Set` — init fetch currency ID
    - `InputDataCredit.CARI8` = `.Currency`
  - **26.2** `RDB-List` — fetch currency ID
    - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-CURRENCY`; BrowsePage=`CurrencySearch`; RequestType=`GetDataCurrencyByName_SQL`
  - **26.3** `Property-Set`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).PaymentTotal` = `local.netpremi`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).InstallmentPercentage` = `.PctTotal`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).Currency` = `.Currency`
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
    - `local.subscript` = `.pxListSubscript`
  - **26.4** `(kosong)` · halaman `.InstallmentList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - syarat `.Currency==local.currency`: benar → ; salah → lewati langkah [kode -/3] *(When langkah tidak dicentang)*
    - **26.4.1** `Property-Set`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Premium` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction*@divide(.InstallmentPct,100,10),local.netpremi*@divide(.InstallmentPct,100,10))`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Premium` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,local.deduction*@divide(.InstallmentPct,100,10),local.netpremi*@divide(.InstallmentPct,100,10))`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).Currency` = `.Currency`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).InstallmentPercentage` = `.InstallmentPct`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).InstallmentNo` = `.Installment`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(local.subscript).InstallmentList(<CURRENT>).DueDate` = `.PaymentDate`
- **27** `Call InsertToTreatyOutXOLList` — Insert Values to Treaty XOL List
- **28** `Call TreatyNonPropOutSetSpreading` — Set Spreading For Non Prop (both with Additional Retro or with only one retro)
  - parameter: spreading=`Local.spreading`
- **29** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.TreatyXOLDifferenceList` · label `//` — Set spreading based on master data. Loop per currency (not sure if this is getting used or not, but the older step (step 14 is alright for single currency.)
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **29.1** `Property-Set`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<APPEND>).TreatyType` = `pyWorkPage.TreatyIn.Share(1).SpreadingTypeIDXOL`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<APPEND>).Currency` = `.Currency`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<APPEND>).CurrencyID` = `.IDCurrency`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).PremiumSpreaded` = `.NetPremi`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).SharePercentage` = `"100"`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
    - `local.sprderror` = `"Spreading in Master data is incomplete"`

### 6 · `InputPolicyTreatyOutDetail_preACT`

Kelas `ASM-FW-GISFW-Work` · 351.844 B · parameter: `SOB`, `CLASSOFBUSINESSID`, `TREATYYEAR`, `LIMITCURRENCY`, `CEDING`, `TREATYGROUP`, `TREATYTYPE`, `ID`

- **1** `Page-New` · halaman `pyWorkPage.TreatyIn`
- **2** `Property-Set` — RD BrowseTreatyInDetail
  - `Param.pyReportName` = `"BrowseTreatyOutDetail"`
  - `Param.pyReportClass` = `"ASM-FW-GISFW-Int-TREATYOUTDETAIL"`
- **3** `Call pxRetrieveReportData` · halaman `ReportPage`
  - parameter: pyUseAlternateDB=`false`; pySkipSummaryProcessing=`false`
- **4** `Property-Set` — Set property
  - `pyWorkPage.PolicyTreatyIn.IDCurrency` = `pyReportContentPage.pxResults(1).CURRENCYID`
  - `pyWorkPage.PolicyTreatyIn.Currency` = `pyReportContentPage.pxResults(1).LIMITCURRENCY`
  - `pyWorkPage.PolicyTreatyIn.TreatyYear` = `pyReportContentPage.pxResults(1).TREATYYEAR`
  - `pyWorkPage.PolicyTreatyIn.TreatyGroupID` = `pyReportContentPage.pxResults(1).TREATYGROUPID`
  - `pyWorkPage.PolicyTreatyIn.TreatyGroupName` = `pyReportContentPage.pxResults(1).TREATYGROUP`
  - `pyWorkPage.PolicyTreatyIn.BizName` = `pyReportContentPage.pxResults(1).CLASSOFBUSINESS`
  - `pyWorkPage.PolicyTreatyIn.BizCode` = `pyReportContentPage.pxResults(1).CLASSOFBUSINESSID`
  - `pyWorkPage.Quotation.BusinessName` = `pyReportContentPage.pxResults(1).CLASSOFBUSINESS`
  - `pyWorkPage.PolicyTreatyIn.SOBName` = `pyWorkPage.Quotation.SobName`
  - `pyWorkPage.PolicyTreatyIn.SOB` = `pyWorkPage.Quotation.SourceOfBusiness`
  - `pyWorkPage.PolicyTreatyIn.CedingCoName` = `pyReportContentPage.pxResults(1).CEDING`
  - `pyWorkPage.PolicyTreatyIn.CedingCo` = `pyReportContentPage.pxResults(1).CEDINGID`
  - `pyWorkPage.PolicyTreatyIn.InsuredID` = `pyWorkPage.Quotation.InsuredID`
  - `pyWorkPage.PolicyTreatyIn.InsuredName` = `pyWorkPage.Quotation.InsuredName`
  - `pyWorkPage.PolicyTreatyIn.NoOffer` = `pyReportContentPage.pxResults(1).TREATYID`
  - `pyWorkPage.PolicyTreatyIn.PremiOgp` = `pyReportContentPage.pxResults(1).MDPVALUE`
  - `pyWorkPage.PolicyTreatyIn.Installment` = `pyReportContentPage.pxResults(1).INSTALLMENTNO`
  - `param.Installment` = `pyReportContentPage.pxResults(1).INSTALLMENTNO`
  - `pyWorkPage.PolicyTreatyIn.Deduction1` = `pyReportContentPage.pxResults(1).DEDUCTION1`
  - `pyWorkPage.PolicyTreatyIn.Deduction2` = `pyReportContentPage.pxResults(1).DEDUCTION2`
  - `pyWorkPage.PolicyTreatyIn.TreatyType` = `pyReportContentPage.pxResults(1).TREATYTYPE`
  - `pyWorkPage.PolicyTreatyIn.LayerType` = `pyReportContentPage.pxResults(1).LAYERTYPE`
  - `pyWorkPage.PolicyTreatyIn.Layer` = `pyReportContentPage.pxResults(1).LAYER`
  - `pyWorkPage.PolicyTreatyIn.LayerPartType` = `pyReportContentPage.pxResults(1).LAYERPARTTYPE`
  - `pyWorkPage.PolicyTreatyIn.LayerPart` = `pyReportContentPage.pxResults(1).LAYERPART`
  - `pyWorkPage.PolicyTreatyIn.ShareCurrency` = `pyReportContentPage.pxResults(1).SHARECURRENCY`
  - `pyWorkPage.PolicyTreatyIn.ShareValue` = `pyReportContentPage.pxResults(1).SHAREVALUE`
  - `pyWorkPage.Quotation.ProportionalType` = `pyReportContentPage.pxResults(1).PROPORTIONTYPE`
  - `pyWorkPage.PolicyTreatyin.QuotationData.SourceOfBusiness` = `pyWorkPage.Quotation.SourceOfBusiness`
  - `pyWorkPage.PolicyTreatyin.QuotationData.SobName` = `pyWorkPage.Quotation.SobName`
  - `pyWorkPage.Quotation.CedingCoName` = `pyReportContentPage.pxResults(1).CEDING`
  - `pyWorkPage.Quotation.CedingCo` = `pyReportContentPage.pxResults(1).CEDINGID`
  - `pyWorkPage.Quotation.ProportionalType` = `pyReportContentPage.pxResults(1).PROPORTIONTYPE`
- **5** `Call SetTreatyCurrencyID` · halaman `pyWorkPage` — When Currency ID is null, set currency ID with this activity
  - syarat `pyWorkPage.PolicyTreatyIn.IDCurrency==""`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Currency=`pyWorkPage.PolicyTreatyIn.Currency`
- **6** `Property-Set` — when pyReportContentPage.pxResults(1).TREATYTYPE == ""
  - syarat `pyReportContentPage.pxResults(1).TREATYTYPE == ""`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.TreatyType` = `"XOL"`
- **7** `Property-Set` — Copy ID, set treatygroupid to parameter page
  - `TreatyIn.ID` = `pyReportContentPage.pxResults(1).ID`
  - `Param.ID` = `.PolicyTreatyIn.TreatyGroupID`
- **8** `Call FetchTreatyGroupOldID` · halaman `pyWorkPage` · label `//` — Fetch OLDID into PolicyTreatyIn.TreatyGroupOldID
- **9** `Call FetchTreatyGroupOJK` · halaman `pyWorkPage` — Fetch Treaty Group OJK
- **10** `RDB-List` · halaman `TempResult` — SQL BrowseTreatyInDetailJoinEDM
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; ClassName=`ASM-FW-GISFW-Int-TREATYOUTDETAIL`; Access=`RNM`; BrowsePage=`TempResult`; RequestType=`BrowseTreatyOutDetail`
- **11** `Java` — Mapping
  - Java:

```java
try{
  ClipboardPage tempPage2 = tools.findPage("TreatyIn");
  ClipboardPage DataJSON = tools.findPage("TempResult");
  ClipboardProperty DataJSONList  = DataJSON.getProperty(".pxResults");
  java.util.Iterator DataJSONListIter = DataJSONList.iterator();
  while (DataJSONListIter.hasNext())
  {
    ClipboardProperty DataJSONData = (ClipboardProperty)DataJSONListIter.next();
    ClipboardPage DataJSONPage = DataJSONData.getPageValue();
    String IsiDataJson = DataJSONPage.getString("CLASSOFBUSINESS");
    tempPage2.adoptJSONObject(IsiDataJson);
  }
} catch(InvalidStreamError e){
  oLog.error("ReloadSection:Invalid JSON Stream for data page params : "+e.getMessage());
} catch(Exception e){
  oLog.error("ReloadSection:Expection : "+e.getMessage());
}
```

- **12** `Property-Remove` — Delete old installmentList
  - `Property=pyWorkPage.PolicyTreatyIn.ListInstallment`
- **13** `call CountResult1_Act` · halaman `pyWorkPage.PolicyTreatyIn` · label `//`
- **14** `(kosong)` · halaman `TreatyIn.INSTALLMENT` — Insert new installmentList
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **14.1** `(kosong)` · halaman `.InstallmentList`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - syarat `pyWorkPage.PolicyTreatyIn.Currency==.Currency`: benar → ; salah → lewati langkah [kode -/3]
    - `pyWorkPage.PolicyTreatyIn.ListInstallment(<APPEND>).InstallmentNo` = `.AcceptStatus`
    - **14.1.1** `Property-Set`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<APPEND>).InstallmentNo` = `.Installment`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).DueDate` = `.PaymentDate`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).InstallmentPercentage` = `.InstallmentPct`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).Premium` = `@Math.divide(((pyWorkPage.PolicyTreatyIn.BalanceDueTo)*.InstallmentPct),100,4)`
      - `pyWorkPage.PolicyTreatyIn.ListInstallment(<LAST>).PaymentTotal` = `@Math.divide(((pyWorkPage.PolicyTreatyIn.BalanceDueTo)*.InstallmentPct),100,4)`
- **15** `(kosong)` — Set Bisnis & Insured Quotation
  - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefLimit=`1`, pyStepsRepeatDefStart=`1`
  - **15.1** `Property-Set` — Copy InsuredName
    - `SearchClient.CARI1` = `pyWorkPage.Quotation.InsuredName`
  - **15.2** `RDB-List` — SQL GetClientID_SQL
    - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-CLIENT`; BrowsePage=`DataClient`; RequestType=`GetClientID_SQL`
  - **15.3** `Property-Set` — Set InsuredID
    - syarat `pyWorkPage.Quotation.InsuredID==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `pyWorkPage.Quotation.InsuredID` = `DataClient.pxResults(1).CARI1`
    - `pyWorkPage.PolicyTreatyIn.InsuredID` = `DataClient.pxResults(1).CARI1`
    - `pyWorkPage.Quotation.InsuredID` = `DataClient.pxResults(1).CARI1`
  - **15.4** `Property-Set` — Copy BusinessName
    - `OldID.CARI2` = `pyWorkPage.Quotation.BusinessName`
  - **15.5** `(kosong)` — SET BUSINESS CODE KALAU KOSONG
    - ulang: pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
    - syarat `pyWorkPage.PolicyTreatyIn.BizCode==""`: benar → lanjut; salah → lewati langkah [kode 2/3] *(When langkah tidak dicentang)*
    - **15.5.1** `Property-Set`
      - `OldID.CARI2` = `pyWorkPage.Quotation.BusinessName`
    - **15.5.2** `Property-Set` — MBU
      - syarat `@contains(pyWorkPage.Quotation.BusinessName,"MBU")`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"MOTOR VEHICLE"`
    - **15.5.3** `Property-Set` — pyWorkPage.Quotation.BusinessName=="ADVANCE PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="BID OR TENDER BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PAYMENT BONDS"\|\|pyWorkPage.Quotation.Bus …
      - syarat `pyWorkPage.Quotation.BusinessName=="ADVANCE PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="BID OR TENDER BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PAYMENT BONDS"\|\|pyWorkPage.Quotation.BusinessName=="PERFORMANCE BONDS"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `@replaceAll(pyWorkPage.Quotation.BusinessName,"S","")`
    - **15.5.4** `Property-Set` — pyWorkPage.Quotation.BusinessName=="CUSTOMS BOND">>OTHERS CUSTOMS BOND
      - syarat `pyWorkPage.Quotation.BusinessName=="CUSTOMS BOND"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"OTHERS CUSTOMS BOND"`
    - **15.5.5** `Property-Set` — pyWorkPage.Quotation.BusinessName=="ASURANSI KREDIT">>ASURANSI KREDIT (CASH LOAN)
      - syarat `pyWorkPage.Quotation.BusinessName=="ASURANSI KREDIT"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"ASURANSI KREDIT (CASH LOAN)"`
    - **15.5.6** `Property-Set` — pyWorkPage.Quotation.BusinessName=="BOILER & PRESSURE VESSEL">>BOILER & EXCAVATOR
      - syarat `pyWorkPage.Quotation.BusinessName=="BOILER & PRESSURE VESSEL"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"BOILER & EXCAVATOR"`
    - **15.5.7** `Property-Set` — pyWorkPage.Quotation.BusinessName=="GOLF INSURANCE">>HOLE IN ONE
      - syarat `pyWorkPage.Quotation.BusinessName=="GOLF INSURANCE"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `OldID.CARI2` = `"HOLE IN ONE"`
  - **15.6** `RDB-List` — SQL GetOldIDBusiness_SQL
    - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; ClassName=`ASM-FW-GISFW-Int-OFFERJSON`; Access=`ASM`; BrowsePage=`DataClient`; RequestType=`GetOldIDBusiness_SQL`
  - **15.7** `Property-Set` — Set BusinessOldId, GroupPanel, & BusinessCode
    - `pyWorkPage.Quotation.BusinessOldId` = `DataClient.pxResults(1).CARI1`
    - `pyWorkPage.Quotation.GroupPanel` = `DataClient.pxResults(1).CARI2`
    - `pyWorkPage.Quotation.BusinessCode` = `DataClient.pxResults(1).CARI3`
    - `pyWorkPage.PolicyTreatyIn.BizCode` = `pyWorkPage.Quotation.BusinessCode`
  - **15.8** `Property-Map-DecisionTable` · halaman `pyWorkPage` — DT BusinessType_DeT
    - parameter: PropertyName=`pyWorkPage.Quotation.BusinessType`; AllowMissingProperties=`false`; DecisionTableName=`BusinessType_DeT`
  - **15.9** `Property-Set` — Set BusinessType
    - `pyWorkPage.PolicyTreatyIn.QuotationData` = `pyWorkPage.Quotation`
- **16** `Page-Remove` — Remove page pyReportContentPage
  - `Page=pyReportContentPage`
- **17** `Call InputPolicyTreatyOutDetail_NonProp` · halaman `pyWorkPage` — When Non Proportional
  - syarat `pyWorkPage.Quotation.ProportionalType=="NonProportional"`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.`
- **18** `Obj-Save` · halaman `pyWorkPage` — Save
  - parameter: WithErrors=`true`; OnlyIfNew=`false`; WriteNow=`true`

### 6 · `InputQuotation_PreAct`

Kelas `ASM-FW-GISFW-Data-Quotation` · 37.013 B · parameter: `Acton`

- **1** `Apply-DataTransform`
  - syarat `Param.Acton=="SOB"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: DataTransform=`btnSOB_DT`; PassParameterPage=`false`
- **2** `Apply-DataTransform`
  - syarat `Param.Acton=="Ceding"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: DataTransform=`btnCedingCO_DT`; PassParameterPage=`false`

### 6 · `InsertHistoryAkseptasiPega`

Kelas `ASM-FW-GISFW-Work` · 100.009 B · parameter: `Status`, `idpega`

- **1** `Property-Set`
  - `InsertHistory.CARI1` = `pyWorkPage.pzInsKey`
  - `InsertHistory.CARI2` = `pyWorkPage.PositionNote`
  - `InsertHistory.CARI3` = `@CurrentDateTime()`
  - `InsertHistory.CARI4` = `OperatorID.pyUserName`
- **2** `Property-Set` — Accept
  - syarat `pyWorkPage.PolicyTreatyIn.IsApproved=="1"\|\|pyWorkPage.EmailTypeUWEDM==1\|\|pyWorkPage.EmailType==1\|\|pyWorkPage.EmailTypeBinding==1\|\|pyWorkPage.EmailTypeRetro==1\|\|pyWorkPage.EmailTypeRetroSlip==1\|\|pyWorkPage.EmailTypeUW==1\|\|pyWorkPage.EmailTypeUWPolicy==1\|\|pyWorkPage.EmailTypeTL==1`: benar → ; salah → lewati langkah [kode -/3]
  - `InsertHistory.CARI5` = `"ACCEPT"`
- **3** `Property-Set` — Reject
  - syarat `pyWorkPage.PolicyTreatyIn.IsApproved=="0"\|\|pyWorkPage.EmailTypeUWEDM==2\|\|pyWorkPage.EmailType==2\|\|pyWorkPage.EmailTypeBinding==2\|\|pyWorkPage.EmailTypeRetro==2\|\|pyWorkPage.EmailTypeRetroSlip==2\|\|pyWorkPage.EmailTypeUW==2\|\|pyWorkPage.EmailTypeTL==2`: benar → ; salah → lewati langkah [kode -/3]
  - `InsertHistory.CARI5` = `"REJECT"`
- **4** `Property-Set` — Ask
  - syarat `pyWorkPage.EmailType==3\|\|pyWorkPage.EmailTypeRetro==3\|\|pyWorkPage.EmailTypeUW==3\|\|pyWorkPage.EmailTypeUWPolicy==3\|\|pyWorkPage.EmailTypeTL==3`: benar → ; salah → lewati langkah [kode -/3]
  - `InsertHistory.CARI5` = `"ASK"`
- **5** `Property-Set` — Revise
  - syarat `pyWorkPage.EmailTypeBinding==9`: benar → ; salah → lewati langkah [kode -/3]
  - `InsertHistory.CARI5` = `"REVISE"`
- **6** `Property-Set` — Decline
  - syarat `pyWorkPage.EmailTypeQuotation=="2"\|\|pyWorkPage.EmailTypeUWEDM==7\|\|pyWorkPage.EmailType==7\|\|pyWorkPage.EmailTypeBinding==7\|\|pyWorkPage.EmailTypeRetro==7\|\|pyWorkPage.EmailTypeRetroSlip==7\|\|pyWorkPage.EmailTypeUW==7\|\|pyWorkPage.EmailTypeUWPolicy==7\|\|pyWorkPage.EmailTypeTL==7`: benar → ; salah → lewati langkah [kode -/3]
  - `InsertHistory.CARI5` = `"DECLINE"`
- **7** `Property-Set` — BANDING
  - syarat `pyWorkPage.EmailType==4\|\|pyWorkPage.EmailTypeBinding==4`: benar → ; salah → lewati langkah [kode -/3]
  - `InsertHistory.CARI5` = `"BANDING"`
- **8** `Property-Set` — Khusus buat admin klik tombol generate
  - syarat `Param.Status=="InboxAdmin"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `InsertHistory.CARI5` = `"INPUT"`
  - `InsertHistory.CARI2` = `"ReasFacInAdmin"`
  - `InsertHistory.CARI1` = `"ASM-FW-GISFW-WORK "+Param.idpega`
- **9** `RDB-List` — RDB Insert ke table
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-int-policyjson`; RequestType=`InsertHistoryAkseptasiPega_Sql`

### 6 · `InsertToTreatyOutXOLList`

Kelas `ASM-FW-GISFW-Work` · 409.888 B

- **1** `Property-Remove` — pyWorkPage.TreatyIn.SpreadingTONP(1).ReinsuranceListTONP(1).LayerList(Local.ShareIndex).DeductionList
  - `Property=pyWorkPage.PolicyTreatyIn.TreatyXOLList`
- **2** `Property-Set` — Set Due To
  - `InputXOL.CARI1` = `"DUE TO US"`
- **3** `Property-Set` · halaman `pyWorkPage.TreatyIn.ShareReins` — Remove SOBname yang beda
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `Local.ShareIndex` = `.pxListSubscript`
  - **3.1** `(kosong)` · halaman `.ReinsuranceListTONP`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **3.1.1** `Property-Set`
      - `Local.Index` = `.pxListSubscript`
    - **3.1.2** `Page-Remove` — pyWorkPage.TreatyIn.ShareReins(Local.ShareIndex).ReinsuranceListTONP(Local.index)
      - syarat `.ReinsName==pyWorkPage.PolicyTreatyIn.SOBName`: benar → lewati langkah; salah → lanjut [kode 3/2]
      - `Page=pyWorkPage.TreatyIn.ShareReins(Local.ShareIndex).ReinsuranceListTONP(Local.Index)`
- **4** `(kosong)` · halaman `pyWorkPage.TreatyIn.Installment` — pyWorkPage.TreatyIn.Share(Local.ShareIndex).DeductionList
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **4.1** `Property-Set` — Init value (Begin XOL List Parent)
    - `local.tempcurrency` = `.Currency`
    - `local.subscript` = `.pxListSubscript`
    - `local.currency` = `""`
    - `local.grosspremi` = `0`
    - `local.netpremi` = `0`
    - `local.duetovalue` = `0`
    - `local.deduction` = `0`
    - `local.dueto` = `""`
    - `Local.PPHValue` = `0`
    - `Local.PPNValue` = `0`
  - **4.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **4.2.1** `Property-Set` — set Local.Index
      - `Local.Index` = `.pxListSubscript`
    - **4.2.2** `(kosong)` · halaman `.ReinsuranceListTONP`
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **4.2.2.1** `Page-New` · halaman `InputXOL` — Init Page
      - **4.2.2.2** `Property-Set` — Set per layer
        - `Local.ShareIndex` = `.pxListSubscript`
      - **4.2.2.3** `(kosong)` · halaman `.GrossPremiumList` — Gross Premi (CARI32, CARI31[CURRENCY])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **4.2.2.3.1** `Property-Set` — pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI32` = `.Value`
          - `InputXOL.CARI31` = `.Currency`
        - **4.2.2.3.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **4.2.2.3.2.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI47` = `.Value`
            - `InputXOL.CARI13` = `.Value`
            - `InputXOL.CARI31` = `.Currency`
        - **4.2.2.3.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **4.2.2.3.3.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI44` = `.Deduction`
      - **4.2.2.4** `Property-Set` — Insert Every Value to the property
        - `local.currency` = `InputXOL.CARI31`
        - `local.grosspremi` = `local.grosspremi + @toDecimal(InputXOL.CARI32)`
        - `local.netpremi` = `local.netpremi + @toDecimal(InputXOL.CARI13)`
        - `local.duetovalue` = `local.duetovalue + @toDecimal(InputXOL.CARI47)`
        - `local.deduction` = `local.deduction + @toDecimal(InputXOL.CARI44)`
        - `local.dueto` = `InputXOL.CARI1`
        - `Local.BrokerageFeeSebenarnya` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(local.deduction,@divide(102.2,100,8),8),local.deduction)`
        - `Local.PPHValue` = `Local.BrokerageFeeSebenarnya* @divide(2,100,8)`
        - `Local.PPNValue` = `Local.BrokerageFeeSebenarnya* @divide(2.2,100,8)`
        - `local.netpremi` = `local.netpremi+ @toDecimal(InputXOL.CARI13)`
        - `local.dueto` = `InputXOL.CARI1`
  - **4.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.SpreadingTONP(1).ReinsuranceListTONP(1).LayerList` · label `//` — Total all layers per currencypyWorkPage.TreatyIn.SpreadingTONP(1).ReinsuranceListTONP(1).LayerList
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **4.3.1** `Property-Set`
      - `Local.Index` = `.pxListSubscript`
    - **4.3.2** `Page-New` · halaman `InputXOL` — Init Page
    - **4.3.3** `Property-Set` — Set per layer
      - `Local.ShareIndex` = `.pxListSubscript`
    - **4.3.4** `(kosong)` · halaman `.GrossPremiumList` — Gross Premi (CARI32, CARI31[CURRENCY])
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **4.3.4.1** `Property-Set` — pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList
        - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
        - `InputXOL.CARI32` = `.Value`
        - `InputXOL.CARI31` = `.Currency`
      - **4.3.4.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **4.3.4.2.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI47` = `.Value`
          - `InputXOL.CARI13` = `.Value`
          - `InputXOL.CARI31` = `.Currency`
      - **4.3.4.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **4.3.4.3.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI44` = `.Deduction`
    - **4.3.5** `Property-Set` — Insert Every Value to the property
      - `local.currency` = `InputXOL.CARI31`
      - `local.grosspremi` = `local.grosspremi + @toDecimal(InputXOL.CARI32)`
      - `local.netpremi` = `local.netpremi + @toDecimal(InputXOL.CARI13)`
      - `local.duetovalue` = `local.duetovalue + @toDecimal(InputXOL.CARI47)`
      - `local.deduction` = `local.deduction + @toDecimal(InputXOL.CARI44)`
      - `local.dueto` = `InputXOL.CARI1`
      - `Local.BrokerageFeeSebenarnya` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(local.deduction,@divide(102.2,100,8),8),local.deduction)`
      - `Local.PPHValue` = `Local.BrokerageFeeSebenarnya* @divide(2,100,8)`
      - `Local.PPNValue` = `Local.BrokerageFeeSebenarnya* @divide(2.2,100,8)`
      - `local.netpremi` = `local.netpremi+ @toDecimal(InputXOL.CARI13)`
      - `local.dueto` = `InputXOL.CARI1`
  - **4.4** `Property-Set` — Init RD cari currency ID
    - `InputDataCredit.CARI8` = `.Currency`
  - **4.5** `RDB-List` — RD cari currency ID
    - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-CURRENCY`; BrowsePage=`CurrencySearch`; RequestType=`GetDataCurrencyByName_SQL`
  - **4.6** `Property-Set` — Set value to total (End XOL List Parent)
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).Currency` = `local.currency`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).GrossPremi` = `local.grosspremi`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).NetPremi` = `local.netpremi`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).DueToValue` = `local.duetovalue`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).Deduction` = `local.deduction`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).DueTo` = `local.dueto`
  - **4.7** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **4.7.1** `Property-Set` — set Local.Index
      - `Local.Index` = `.pxListSubscript`
    - **4.7.2** `(kosong)` · halaman `.ReinsuranceListTONP`
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **4.7.2.1** `Page-New` · halaman `InputXOL` — Init Page
      - **4.7.2.2** `Property-Set` — Set per layer
        - `InputXOL.CARI50` = `.LayerType`
        - `InputXOL.CARI51` = `.Layer`
        - `InputXOL.CARI52` = `.LayerPartType`
        - `InputXOL.CARI53` = `.LayerPart`
        - `Local.ShareIndex` = `.pxListSubscript`
      - **4.7.2.3** `(kosong)` · halaman `.GrossPremiumList` — Gross Premi (CARI32, CARI31[CURRENCY])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **4.7.2.3.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI32` = `.Value`
          - `InputXOL.CARI31` = `.Currency`
        - **4.7.2.3.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **4.7.2.3.2.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI47` = `.Value`
            - `InputXOL.CARI13` = `.Value`
            - `InputXOL.CARI31` = `.Currency`
        - **4.7.2.3.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **4.7.2.3.3.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI44` = `.Deduction`
            - `InputXOL.CARI45` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(@toDecimal(InputXOL.CARI44),@divide(102.2,100,8),8),@toDecimal(InputXOL.CARI44))`
      - **4.7.2.4** `Property-Set` — Insert Every Value to the property
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<APPEND>).Currency` = `InputXOL.CARI31`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).GrossPremi` = `InputXOL.CARI32`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).NetPremi` = `InputXOL.CARI13`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).DueToValue` = `InputXOL.CARI47`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).Deduction` = `InputXOL.CARI44`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerType` = `InputXOL.CARI50`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).Layer` = `InputXOL.CARI51`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerPartType` = `InputXOL.CARI52`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerPart` = `InputXOL.CARI53`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).DueTo` = `InputXOL.CARI1`
  - **4.8** `(kosong)` · halaman `pyWorkPage.TreatyIn.SpreadingTONP(1).ReinsuranceListTONP(1).LayerList` · label `//` — (XOL List Child)
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **4.8.1** `Property-Set` — set Local.Index
      - `Local.Index` = `.pxListSubscript`
    - **4.8.2** `Page-New` · halaman `InputXOL` — Init Page
    - **4.8.3** `Property-Set` — Set per layer
      - `InputXOL.CARI50` = `.LayerType`
      - `InputXOL.CARI51` = `.Layer`
      - `InputXOL.CARI52` = `.LayerPartType`
      - `InputXOL.CARI53` = `.LayerPart`
      - `Local.ShareIndex` = `.pxListSubscript`
    - **4.8.4** `(kosong)` · halaman `.GrossPremiumList` — Gross Premi (CARI32, CARI31[CURRENCY])
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **4.8.4.1** `Property-Set`
        - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
        - `InputXOL.CARI32` = `.Value`
        - `InputXOL.CARI31` = `.Currency`
      - **4.8.4.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **4.8.4.2.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI47` = `.Value`
          - `InputXOL.CARI13` = `.Value`
          - `InputXOL.CARI31` = `.Currency`
      - **4.8.4.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **4.8.4.3.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI44` = `.Deduction`
          - `InputXOL.CARI45` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(@toDecimal(InputXOL.CARI44),@divide(102.2,100,8),8),@toDecimal(InputXOL.CARI44))`
    - **4.8.5** `Property-Set` — Insert Every Value to the property
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<APPEND>).Currency` = `InputXOL.CARI31`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).GrossPremi` = `InputXOL.CARI32`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).NetPremi` = `InputXOL.CARI13`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).DueToValue` = `InputXOL.CARI47`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).Deduction` = `InputXOL.CARI44`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerType` = `InputXOL.CARI50`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).Layer` = `InputXOL.CARI51`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerPartType` = `InputXOL.CARI52`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerPart` = `InputXOL.CARI53`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).DueTo` = `InputXOL.CARI1`
  - **4.9** `(kosong)`

### 6 · `InsertToTreatyXOLList`

Kelas `ASM-FW-GISFW-Work` · 255.410 B

- **1** `Property-Remove`
  - `Property=pyWorkPage.PolicyTreatyIn.TreatyXOLList`
- **2** `Property-Set` — Set Due To
  - `InputXOL.CARI1` = `"DUE TO US"`
- **3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Installment`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **3.1** `Property-Set` — Init value (Begin XOL List Parent)
    - `local.tempcurrency` = `.Currency`
    - `local.subscript` = `.pxListSubscript`
    - `local.currency` = `""`
    - `local.grosspremi` = `0`
    - `local.netpremi` = `0`
    - `local.duetovalue` = `0`
    - `local.deduction` = `0`
    - `local.dueto` = `""`
    - `Local.PPHValue` = `0`
    - `Local.PPNValue` = `0`
  - **3.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share` — Total all layers per currency
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **3.2.1** `Page-New` · halaman `InputXOL` — Init Page
    - **3.2.2** `Property-Set` — Set per layer
      - `Local.ShareIndex` = `.pxListSubscript`
    - **3.2.3** `(kosong)` · halaman `.GrossPremiumList` — Gross Premi (CARI32, CARI31[CURRENCY])
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **3.2.3.1** `Property-Set` — pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList
        - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
        - `InputXOL.CARI32` = `.Value`
        - `InputXOL.CARI31` = `.Currency`
      - **3.2.3.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **3.2.3.2.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI47` = `.Value`
          - `InputXOL.CARI13` = `.Value`
          - `InputXOL.CARI31` = `.Currency`
      - **3.2.3.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **3.2.3.3.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI44` = `.Deduction`
    - **3.2.4** `Property-Set` — Insert Every Value to the property
      - `local.currency` = `InputXOL.CARI31`
      - `local.grosspremi` = `local.grosspremi + @toDecimal(InputXOL.CARI32)`
      - `local.netpremi` = `local.netpremi + @toDecimal(InputXOL.CARI13)`
      - `local.duetovalue` = `local.duetovalue + @toDecimal(InputXOL.CARI47)`
      - `local.deduction` = `local.deduction + @toDecimal(InputXOL.CARI44)`
      - `local.dueto` = `InputXOL.CARI1`
      - `Local.BrokerageFeeSebenarnya` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(local.deduction,@divide(102.2,100,8),8),local.deduction)`
      - `Local.PPHValue` = `Local.BrokerageFeeSebenarnya* @divide(2,100,8)`
      - `Local.PPNValue` = `Local.BrokerageFeeSebenarnya* @divide(2.2,100,8)`
  - **3.3** `Property-Set` — Init RD cari currency ID
    - `InputDataCredit.CARI8` = `.Currency`
  - **3.4** `RDB-List` — RD cari currency ID
    - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-CURRENCY`; BrowsePage=`CurrencySearch`; RequestType=`GetDataCurrencyByName_SQL`
  - **3.5** `Property-Set` — Set value to total (End XOL List Parent)
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).Currency` = `local.currency`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).GrossPremi` = `local.grosspremi`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).NetPremi` = `local.netpremi`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).DueToValue` = `local.duetovalue`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).Deduction` = `local.deduction`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).DueTo` = `local.dueto`
  - **3.6** `Property-Set` — pyWorkPage.PolicyTreatyIn.FlagPPH=="true"
    - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).BrokerageFeeSebenarnya` = `Local.BrokerageFeeSebenarnya`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).PPHValue` = `Local.PPHValue`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).PPNValue` = `Local.PPNValue`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).NetPremiAfterPPH` = `local.netpremi + Local.PPHValue`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).NetPremiAfterPPN` = `local.netpremi + Local.PPNValue`
    - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).NetPremiAfterTax` = `local.netpremi + Local.PPHValue + Local.PPNValue`
  - **3.7** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share` — (XOL List Child)
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **3.7.1** `Page-New` · halaman `InputXOL` — Init Page
    - **3.7.2** `Property-Set` — Set per layer
      - `InputXOL.CARI50` = `.LayerType`
      - `InputXOL.CARI51` = `.Layer`
      - `InputXOL.CARI52` = `.LayerPartType`
      - `InputXOL.CARI53` = `.LayerPart`
      - `Local.ShareIndex` = `.pxListSubscript`
    - **3.7.3** `(kosong)` · halaman `.GrossPremiumList` — Gross Premi (CARI32, CARI31[CURRENCY])
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **3.7.3.1** `Property-Set`
        - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
        - `InputXOL.CARI32` = `.Value`
        - `InputXOL.CARI31` = `.Currency`
      - **3.7.3.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **3.7.3.2.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI47` = `.Value`
          - `InputXOL.CARI13` = `.Value`
          - `InputXOL.CARI31` = `.Currency`
      - **3.7.3.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **3.7.3.3.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI44` = `.Deduction`
          - `InputXOL.CARI45` = `@if(pyWorkPage.PolicyTreatyIn.TypeTax="Inclusive",@divide(@toDecimal(InputXOL.CARI44),@divide(102.2,100,8),8),@toDecimal(InputXOL.CARI44))`
    - **3.7.4** `Property-Set` — Insert Every Value to the property
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<APPEND>).Currency` = `InputXOL.CARI31`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).GrossPremi` = `InputXOL.CARI32`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).NetPremi` = `InputXOL.CARI13`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).DueToValue` = `InputXOL.CARI47`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).Deduction` = `InputXOL.CARI44`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerType` = `InputXOL.CARI50`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).Layer` = `InputXOL.CARI51`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerPartType` = `InputXOL.CARI52`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerPart` = `InputXOL.CARI53`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).DueTo` = `InputXOL.CARI1`
    - **3.7.5** `Property-Set` — pyWorkPage.PolicyTreatyIn.FlagPPH=="true"
      - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).BrokerageFeeSebenarnya` = `InputXOL.CARI45`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).PPHValue` = `@toDecimal(InputXOL.CARI45)* @divide(2,100,8)`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).PPNValue` = `@toDecimal(InputXOL.CARI45)* @divide(2.2,100,8)`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).NetPremiAfterPPH` = `@toDecimal(InputXOL.CARI13)+@toDecimal(InputXOL.CARI45)*@divide(2,100,8)`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).NetPremiAfterPPN` = `@toDecimal(InputXOL.CARI13)+@toDecimal(InputXOL.CARI45)* @divide(2.2,100,8)`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).NetPremiAfterTax` = `@toDecimal(InputXOL.CARI13)+(@toDecimal(InputXOL.CARI45)*@divide(2,100,8)+@toDecimal(InputXOL.CARI45)* @divide(2.2,100,8))`

### 6 · `InsertToTreatyXOLListRetroShare`

Kelas `ASM-FW-GISFW-Work` · 328.838 B

- **1** `Property-Remove`
  - `Property=pyWorkPage.PolicyTreatyIn.TreatyXOLList`
- **2** `(kosong)` — For Non Prop
  - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
  - syarat `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp==1`: benar → ; salah → lewati langkah [kode -/3]
  - **2.1** `Property-Set` — Set Due To
    - `InputXOL.CARI1` = `"DUE TO US"`
  - **2.2** `Property-Set` — Fetch Spreading
    - `InputXOL.CARI14` = `pyWorkPage.TreatyIn.Share(1).SpreadingTypeIDXOL`
  - **2.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Installment` — Per installment treaty in
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **2.3.1** `Property-Set` — Init value
      - `local.tempcurrency` = `.Currency`
      - `local.subscript` = `.pxListSubscript`
      - `local.currency` = `""`
      - `local.grosspremi` = `0`
      - `local.netpremi` = `0`
      - `local.duetovalue` = `0`
      - `local.deduction` = `0`
      - `local.dueto` = `""`
    - **2.3.2** `Page-New` · halaman `InputXOL` — Init Page
    - **2.3.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share` — Total all layers per currency
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **2.3.3.1** `Property-Set` — Set per layer
        - `Local.ShareIndex` = `.pxListSubscript`
      - **2.3.3.2** `(kosong)` · halaman `.GrossPremiumList` — Gross Premi (CARI32, CARI31[CURRENCY])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **2.3.3.2.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI32` = `.Value`
          - `InputXOL.CARI31` = `.Currency`
        - **2.3.3.2.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **2.3.3.2.2.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI47` = `.Value`
            - `InputXOL.CARI13` = `.Value`
            - `InputXOL.CARI31` = `.Currency`
        - **2.3.3.2.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **2.3.3.2.3.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI44` = `.Deduction`
      - **2.3.3.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.FacultativeShareList(Local.ShareIndex).DeductionList` — Gross Premi Retro (CARI32, CARI31[CURRENCY])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
        - **2.3.3.3.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI44` = `.Deduction+ @toDecimal(InputXOL.CARI44)`
          - `InputXOL.CARI31` = `.Currency`
        - **2.3.3.3.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList` · label `//` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **2.3.3.3.2.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI47` = `.Value + @toDecimal(InputXOL.CARI47)`
            - `InputXOL.CARI13` = `.Value + @toDecimal(InputXOL.CARI13)`
            - `InputXOL.CARI31` = `.Currency`
        - **2.3.3.3.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).DeductionList` · label `//` — Deduction (CARI44[deduction])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **2.3.3.3.3.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI44` = `.Deduction + @toDecimal(InputXOL.CARI44)`
      - **2.3.3.4** `Property-Set` — Insert Every Value to the property
        - `local.currency` = `InputXOL.CARI31`
        - `local.grosspremi` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,@toDecimal(InputXOL.CARI44),local.grosspremi + @toDecimal(InputXOL.CARI32))`
        - `local.netpremi` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,@toDecimal(InputXOL.CARI44),local.netpremi + @toDecimal(InputXOL.CARI13))`
        - `local.duetovalue` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,@toDecimal(InputXOL.CARI44),local.duetovalue + @toDecimal(InputXOL.CARI47))`
        - `local.deduction` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,0,local.deduction + @toDecimal(InputXOL.CARI44))`
        - `Local.BrokerageFeeSebenarnya` = `@divide(local.deduction,@divide(102.2,100,8),8)`
        - `Local.PPHValue` = `Local.BrokerageFeeSebenarnya* @divide(2,100,8)`
        - `Local.PPNValue` = `Local.BrokerageFeeSebenarnya* @divide(2.2,100,8)`
        - `local.dueto` = `InputXOL.CARI1`
    - **2.3.4** `Property-Set` — Init RD cari currency ID
      - `InputDataCredit.CARI8` = `.Currency`
    - **2.3.5** `RDB-List` — RD cari currency ID
      - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-CURRENCY`; BrowsePage=`CurrencySearch`; RequestType=`GetDataCurrencyByName_SQL`
    - **2.3.6** `Property-Set` — Set value to total
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).Currency` = `local.currency`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).GrossPremi` = `local.grosspremi`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).NetPremi` = `local.netpremi`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).DueToValue` = `local.duetovalue`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).Deduction` = `local.deduction`
      - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(<CURRENT>).DueTo` = `local.dueto`
    - **2.3.7** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share` — Per Layer
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **2.3.7.1** `Page-New` · halaman `InputXOL` — Init Page
      - **2.3.7.2** `Property-Set` — Set per layer
        - `InputXOL.CARI50` = `.LayerType`
        - `InputXOL.CARI51` = `.Layer`
        - `InputXOL.CARI52` = `.LayerPartType`
        - `InputXOL.CARI53` = `.LayerPart`
        - `Local.ShareIndex` = `.pxListSubscript`
      - **2.3.7.3** `(kosong)` · halaman `.GrossPremiumList` — Gross Premi (CARI32, CARI31[CURRENCY])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **2.3.7.3.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI32` = `.Value`
          - `InputXOL.CARI31` = `.Currency`
        - **2.3.7.3.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **2.3.7.3.2.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI47` = `.Value`
            - `InputXOL.CARI13` = `.Value`
            - `InputXOL.CARI31` = `.Currency`
        - **2.3.7.3.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **2.3.7.3.3.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI44` = `.Deduction`
      - **2.3.7.4** `(kosong)` · halaman `pyWorkPage.TreatyIn.FacultativeShareList(Local.ShareIndex).GrossPremiumList` — Gross Premi Retro (CARI32, CARI31[CURRENCY])
        - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
        - **2.3.7.4.1** `Property-Set`
          - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
          - `InputXOL.CARI32` = `.Value + @toDecimal(InputXOL.CARI32)`
          - `InputXOL.CARI31` = `.Currency`
        - **2.3.7.4.2** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).NetPremiumList` — Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **2.3.7.4.2.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI47` = `.Value + @toDecimal(InputXOL.CARI47)`
            - `InputXOL.CARI13` = `.Value + @toDecimal(InputXOL.CARI13)`
            - `InputXOL.CARI31` = `.Currency`
        - **2.3.7.4.3** `(kosong)` · halaman `pyWorkPage.TreatyIn.Share(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - **2.3.7.4.3.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI44` = `.Deduction + @toDecimal(InputXOL.CARI44)`
        - **2.3.7.4.4** `(kosong)` · halaman `pyWorkPage.TreatyIn.FacultativeShareList(Local.ShareIndex).DeductionList` — Deduction (CARI44[deduction])
          - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
          - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → lanjut; salah → lewati langkah [kode 2/3]
          - **2.3.7.4.4.1** `Property-Set`
            - syarat `.Currency==local.tempcurrency`: benar → ; salah → lewati langkah [kode -/3]
            - `InputXOL.CARI19` = `.Deduction + @toDecimal(InputXOL.CARI44)`
      - **2.3.7.5** `Property-Set` — Insert Every Value to the property
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<APPEND>).Currency` = `InputXOL.CARI31`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).IDCurrency` = `CurrencySearch.pxResults(1).ID`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).GrossPremi` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,InputXOL.CARI19,InputXOL.CARI32)`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).NetPremi` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,InputXOL.CARI19,InputXOL.CARI13)`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).DueToValue` = `@if(pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true,InputXOL.CARI19,InputXOL.CARI47)`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).Deduction` = `InputXOL.CARI44`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerType` = `InputXOL.CARI50`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).Layer` = `InputXOL.CARI51`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerPartType` = `InputXOL.CARI52`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).LayerPart` = `InputXOL.CARI53`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).DueTo` = `InputXOL.CARI1`
      - **2.3.7.6** `Property-Set`
        - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH=="true"`: benar → lanjut; salah → lewati langkah [kode 2/3]
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).BrokerageFeeSebenarnya` = `InputXOL.CARI45`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).PPHValue` = `@toDecimal(InputXOL.CARI45)* @divide(2,100,8)`
        - `pyWorkPage.PolicyTreatyIn.TreatyXOLList(local.subscript).ValueList(<LAST>).PPNValue` = `@toDecimal(InputXOL.CARI45)* @divide(2.2,100,8)`

### 6 · `ProtectDate`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 31.891 B

- **1** `Property-Set` — set Error message
  - `local.Errmsg` = `"End Date Cannot be less than Start Date"`
  - `.pyMessageLabel` = `""`
- **2** `Property-Set-Messages` — Set Error message when end date > start date
  - syarat `@CompareDates(.StartDate,.EndDate)`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Field=`.EndDate`; Category=`pyMessageLabel`; Message=`local.Errmsg`; ContainingClassOfProperty=`ASM-FW-GISFW-Data-PolicyTreatyIn`

### 6 · `RemoveTypeTax_ACT`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 27.393 B

- **1** `Property-Remove`
  - syarat `TempWorkPage.PolicyTreatyIn.FlagPPH==false`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `Property=TempWorkPage.PolicyTreatyIn.TypeTax`
- **2** `Property-Remove`
  - syarat `pyWorkPage.PolicyTreatyIn.FlagPPH==false`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `Property=pyWorkPage.PolicyTreatyIn.TypeTax`

### 6 · `SaveJsonPolisTreatyIn_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 85.560 B · parameter: `isFOR`

- **1** `Call CopyToPolicy` · label `//` — Copy All To Policy Property
  - syarat `.Quotation.BusinessCode != "40"`: benar → lanjut; salah → lewati langkah [kode 2/3]
- **2** `Property-Set`
  - `.Policy.CaseID` = `.pzInsKey`
  - `.Policy.OperatorID` = `OperatorID.pyUserIdentifier`
  - `.Policy.Quotation.OperatorID` = `OperatorID.pyUserIdentifier`
- **3** `Property-Set` · halaman `InputData` · label `//`
  - syarat `@equals(param.isFOR,"EDM")`: benar → ; salah → lewati langkah [kode -/3]
  - `InputData.CARI1` = `pyWorkPage.pzInsKey`
  - `InputData.CARI2` = `pyWorkPage.PolicyNumber`
- **4** `RDB-List`
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-POLISTREATYIN`; RequestType=`GenerateNoPolicy`
- **5** `Property-Set` · halaman `InputData`
  - syarat `@equals(param.isFOR,"POLICY")`: benar → ; salah → lewati langkah [kode -/3] *(When langkah tidak dicentang)*
  - `InputData.CARI1` = `pyWorkPage.pzInsKey`
  - `InputData.CARI2` = `OutputParam.ERRMSG2`
- **6** `Property-Set` · halaman `pyWorkPage.PolicyTreatyIn`
  - `InputData.CARI3` = `@ASM.GetPageJSONString()`
- **7** `RDB-List` — insert ke dalam json_polis
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-POLISTREATYIN`; RequestType=`SavePolisTreatyIn_SQL`
- **8** `Property-Set`
  - `OutputParam.ERRMSG` = `OutputParam.IDPEGAOUT`

### 6 · `SaveViewSuggest`

Kelas `ASM-FW-GISFW-Work` · 128.682 B

- **1** `Property-Set` · halaman `pyWorkPage.OfferFacIn.ViewSuggest` — UNTUK FACIN
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `pyWorkPage.Quotation.BusinessFac=="F"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `Local.Idx` = `.pxListSubscript`
  - **1.1** `Property-Remove` — HAPUS JIKA PIC SUGEST KOSONG
    - syarat `.PICSuggest=="" && .DateSuggest = "" && .CommentSuggest = ""`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `Property=pyWorkPage.OfferFacIn.ViewSuggest(Local.Idx)`
  - **1.2** `(kosong)`
    - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
    - syarat `.PICSuggest=="" && .DateSuggest = "" && .CommentSuggest = ""`: benar → lewati langkah; salah → lanjut [kode 3/2]
    - syarat `.IsSave==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - **1.2.1** `Page-New` · halaman `InputData`
    - **1.2.2** `Property-Set` — set param
      - `InputData.CARI1` = `@replaceAll(pyWorkPage.pyWorkIDPrefix,"-","")`
      - `InputData.CARI2` = `.pxListSubscript`
      - `InputData.CARI3` = `.IsCedingConfirm`
      - `InputData.CARI4` = `.PICSuggest`
      - `InputData.CARI5` = `@FormatDateTime(.DateSuggest,"dd/MM/yyyy hh:mm:ss","Asia/Jakarta","in_ID")`
      - `InputData.CARI6` = `OperatorID.pyOrgDivision`
      - `InputData.CARI7` = `pyWorkPage.Quotation.BusinessFac`
      - `InputData.CARI8` = `@if(.IsCedingConfirm="Offer","1",@if(.IsCedingConfirm="Binding"\|\|.IsCedingConfirm="Accepted"\|\|.IsCedingConfirm="Policy","2",""))`
      - `InputData.CARI9` = `@if(.Approval="1","Accept",@if(.Approval="2","Reject",@if(.Approval="3","Ask",@if(.Approval="4","Banding",@if(.Approval="5","Reject Ceding",@if(.Approval="6","Ask Ceding",@if(.Approval="7","Decline",@if(.Approval="9","Revise",""))))))))`
      - `InputData.CARI10` = `@substring(.CommentSuggest,0,3990)`
    - **1.2.3** `RDB-List`
      - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Work`; Access=`RNM`; RequestType=`InsertViewSuggest_SQL`
    - **1.2.4** `Property-Set`
      - `.IsSave` = `"Yes"`
- **2** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.SuggestList` — UNTUK TREATY
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `pyWorkPage.Quotation.BusinessFac=="F"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **2.1** `(kosong)`
    - ulang: pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
    - syarat `.IsSave==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - **2.1.1** `Page-New` · halaman `InputData`
    - **2.1.2** `Property-Set` — set param
      - `InputData.CARI1` = `@replaceAll(pyWorkPage.pyWorkIDPrefix,"-","")`
      - `InputData.CARI2` = `.pxListSubscript`
      - `InputData.CARI3` = `"Policy"`
      - `InputData.CARI4` = `.OperatorName`
      - `InputData.CARI5` = `@FormatDateTime(.Date,"dd/MM/yyyy hh:mm:ss","Asia/Jakarta","in_ID")`
      - `InputData.CARI6` = `OperatorID.pyOrgDivision`
      - `InputData.CARI7` = `pyWorkPage.Quotation.BusinessFac`
      - `InputData.CARI8` = `2`
      - `InputData.CARI9` = `@if(.IsApproved="1","Accept", @if(.IsApproved="0","Reject",""))`
      - `InputData.CARI10` = `@substring(.Suggest,0,3990)`
    - **2.1.3** `RDB-List`
      - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Work`; RequestType=`InsertViewSuggest_SQL`
    - **2.1.4** `Property-Set`
      - `.IsSave` = `"Yes"`

### 6 · `SearchHierarkiSourceBizAgentTreatyIn_Act`

Kelas `ASM-FW-GISFW-Data-Agent` · 63.862 B

- **1** `RDB-List` · halaman `ClientEmail`
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-CLIENT`; RequestType=`BrowseClientEmail_SQL`
- **2** `Property-Set`
  - `pyWorkPage.OfferTreatyIn.QuotationData.SourceOfBusiness` = `@if(.ChildCount > 0, "", .ID)`
  - `pyWorkPage.OfferTreatyIn.QuotationData.SobName` = `@if(.ChildCount > 0, "", .ClientName)`
  - `pyWorkPage.OfferTreatyIn.QuotationData.SobLeader0` = `@if(.ChildCount > 0, "", .Leader0)`
  - `pyWorkPage.OfferTreatyIn.QuotationData.SobLeader1` = `@if(.ChildCount > 0, "", .Leader1)`
  - `Local.ClientID` = `@if(.ChildCount > 0, "", .ClientID)`
  - `Local.IsEmailAddress` = `false`
- **3** `(kosong)` · halaman `ClientEmail.pxResults`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **3.1** `Property-Set`
    - syarat `Local.ClientID==.ID`: benar → ; salah → lewati langkah [kode -/3]
    - `pyWorkPage.OfferTreatyIn.QuotationData.Email` = `.EMAILADDRESS`
    - `Local.IsEmailAddress` = `true`
- **4** `Property-Set`
  - syarat `Local.IsEmailAddress==false`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.OfferTreatyIn.QuotationData.Email` = `"The company has no email"`

### 6 · `serviceInsertArasapas_act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 18.013 B

- **1** `Call serviceInsertArasapas_act` · halaman `pyWorkPage`

### 6 · `SetCategoryAttach`

Kelas `ASM-FW-GISFW-Data-OfferFacIn` · 87.436 B · parameter: `pzInsKey`, `category`

- **1** `RDB-List`
  - syarat `pyWorkPage.Quotation.BusinessFac == "F" \|\| pyWorkPage.Quotation.BusinessFac == "T"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; ClassName=`ASM-FW-GISFW-Int-OFFERJSON`; Access=`ASM`; BrowsePage=`Attachment`; RequestType=`CategoryAttach_SQL`
- **2** `RDB-List`
  - syarat `pyWorkPage.Quotation.BusinessFac == "F" \|\| pyWorkPage.Quotation.BusinessFac == "T"`: benar → lewati langkah; salah → lanjut [kode 3/2]
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Int-OFFERJSON`; Access=`ASM`; BrowsePage=`Attachment`; RequestType=`AttachmentLife`
- **3** `Apply-DataTransform` · halaman `TempLIfe`
  - syarat `pyWorkPage.Quotation.BusinessFac == "F" \|\| pyWorkPage.Quotation.BusinessFac == "T"`: benar → lewati langkah; salah → lanjut [kode 3/2]
  - parameter: DataTransform=`setCategoryAttachment_DT`; PassParameterPage=`false`
- **4** `(kosong)` · halaman `Attachment.pxResults`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **4.1** `Property-Set`
    - `Param.GCNMCategory` = `.NOTE`
    - `Param.pzInsKey` = `pyWorkPage.pzInsKey`
    - `Param.category` = `"Reas"`
  - **4.2** `call  ASM-FW-GISFW-Int-OFFERJSON.InputParamUploadReas_act`
    - parameter: GCNMCategory=`Param.GCNMCategory`; category=`"Reas"`
  - **4.3** `Property-Set`
    - `.CountAttach` = `@Utilities.LengthOfPageList(AttachShowList.pxResults)`
  - **4.4** `Page-New` · halaman `AttachShowList`

### 6 · `SetCurrency_act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 39.111 B · parameter: `CURR`

- **1** `Property-Set`
  - `InputData.CARI1` = `Param.CURR`
- **2** `RDB-List`
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-CURRENCY`; BrowsePage=`Currency`; RequestType=`GetCurrency`
- **3** `Property-Set`
  - `.Currency` = `Currency.pxResults(1).HASIL1`

### 6 · `SetDueTo_act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 27.597 B

- **1** `Property-Set` — set due to us
  - syarat `.BalanceDueTo>=0`: benar → ; salah → lewati langkah [kode -/3]
  - `.DueTo` = `1`
- **2** `Property-Set`
  - syarat `.BalanceDueTo<0`: benar → ; salah → lewati langkah [kode -/3]
  - `.DueTo` = `0`

### 6 · `SetPPNPPH`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 63.772 B

- **1** `Property-Set` — Parameter for calling RD
  - `Param.pyReportName` = `"BrowseClientName_RD"`
  - `Param.pyReportClass` = `ASM-FW-GISFW-Int-AGENT`
  - `param.pyPageName` = `ListAgent`
  - `Param.ID` = `pyWorkPage.PolicyTreatyIn.QuotationData.SourceOfBusiness`
- **2** `Property-Set`
  - `.PPHValue` = `0`
  - `.PPNValue` = `0`
- **3** `Call Rule-Obj-Report-Definition.pxRetrieveReportData`
  - parameter: pyUseAlternateDB=`false`; pySkipSummaryProcessing=`false`
- **4** `Property-Set` — Set PPN & PPH
  - syarat `.FlagPPH=="true"`: benar → lewati syarat berikut; salah → lanjut [kode 5/2]
  - syarat `ListAgent.pxResults(1).STS_PKP == 1`: benar → ; salah → lewati langkah [kode -/3]
  - `.BrokerageFee` = `@divide(2.5,100,8) * (.PremiOgp + .PremiOnp)`
  - `.BrokerageFeeSebenarnya` = `@if(.TypeTax=="Inclusive",@divide(.Deduction1,@divide(102.2,100,8),8),.Deduction1)`
  - `.PPHValue` = `.BrokerageFeeSebenarnya* @divide(2,100,8)`
  - `.PPNValue` = `.BrokerageFeeSebenarnya* @divide(2.2,100,8)`
- **5** `Page-Remove` — Remove ListAgent
  - `Page=ListAgent`

### 6 · `SetReinstatementPct`

Kelas `ASM-FW-GISFW-Data-TreatyInLimits` · 66.961 B

- **1** `Property-Remove` — init
  - `Property=.Reinstatement_List`
- **2** `Property-Set` — init
  - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefLimit=`Current.ReinstatementValue`, pyStepsRepeatDefStart=`0`
  - `Local.maxidx` = `.ReinstatementValue`
  - `Local.runningidx` = `0`
  - `Local.LimitsSubscript` = `.pxListSubscript`
- **3** `(kosong)` — Loop as much as how much reinstatement
  - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`Local.maxidx`
  - **3.1** `Property-Set`
    - `.Reinstatement_List(<APPEND>).ReinstatementPct` = `"100"`
    - `.Reinstatement_List(<LAST>).ReinstatementValue` = `Local.runningidx + 1`
    - `.Reinstatement_List(<LAST>).ReinstatementNote` = `.ReinstatementNote`
    - `.Reinstatement_List(<LAST>).ReinstatementAmount1` = `.Limit`
    - `.Reinstatement_List(<LAST>).ReinstatementAmount2` = `.Limit2`
    - `.Reinstatement_List(<LAST>).ID` = `Local.LimitsSubscript`
    - `.Reinstatement_List(<LAST>).AdditionalPct` = `"100"`
    - `Local.runningidx` = `Local.runningidx + 1`
  - **3.2** `Property-Set` — @SizeOfPropertyList(.MDPList)==1
    - syarat `@SizeOfPropertyList(.MDPList)==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Reinstatement_List(<LAST>).AdditionalAmount1` = `@if(.MDPList(1).Currency=="IDR",.MDPList(1).Value,0)`
    - `.Reinstatement_List(<LAST>).AdditionalAmount2` = `@if(.MDPList(1).Currency=="USD",.MDPList(1).Value,0)`
  - **3.3** `Property-Set` — @SizeOfPropertyList(.MDPList)>2
    - syarat `@SizeOfPropertyList(.MDPList)>=2`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - `.Reinstatement_List(<LAST>).AdditionalAmount1` = `@if(.MDPList(1).Currency=="IDR",.MDPList(1).Value,0)`
    - `.Reinstatement_List(<LAST>).AdditionalAmount2` = `@if(.MDPList(2).Currency=="USD",.MDPList(2).Value,0)`

### 6 · `SetSurveyReport_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 22.275 B

- **1** `Property-Set`
  - `pyWorkPage.Quotation` = `pyWorkPage.PolicyTreatyIn.QuotationData`

### 6 · `SetTreatyCurrencyID`

Kelas `ASM-FW-GISFW-Work` · 40.377 B · parameter: `Currency`

- **1** `Property-Set` — Init SQL
  - `InputData.CARI1` = `pyWorkPage.PolicyTreatyIn.Currency`
- **2** `RDB-List`
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Int-TREATY_IN`; Access=`ASM`; BrowsePage=`OutputData`; RequestType=`GetCurrencyIDByName`
- **3** `Property-Set`
  - `pyWorkPage.PolicyTreatyIn.IDCurrency` = `OutputData.pxResults(1).CARI1`

### 6 · `SetTreatyIn_Act`

Kelas `Data-Portal` · 159.737 B · parameter: `ID`, `viewstate`, `revisionstate`

- **1** `Property-Set`
  - `FlagExcel.CARI1` = `0`
- **2** `call TreatyInInputVis`
  - parameter: Add=`1`; bypasscondition=`1`
- **3** `Property-Set` — Show layout input & reset error
  - `OutputParam.ERRMSG` = `""`
  - `TreatyIn.ID` = `Param.ID`
  - `TreatyWarning.CARI1` = `""`
- **4** `RDB-List` · halaman `TempResult` — Call RDB List
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Int-TREATY_IN`; Access=`ASM`; BrowsePage=`TempResult`; RequestType=`BrowseTreatyIn`
- **5** `Java` — Map oracle column to clipboard
  - Java:

```java
try{
  ClipboardPage tempPage2 = tools.findPage("TreatyIn");
  ClipboardPage DataJSON = tools.findPage("TempResult");
  ClipboardProperty DataJSONList  = DataJSON.getProperty(".pxResults");
  java.util.Iterator DataJSONListIter = DataJSONList.iterator();
  while (DataJSONListIter.hasNext())
  {
    ClipboardProperty DataJSONData = (ClipboardProperty)DataJSONListIter.next();
    ClipboardPage DataJSONPage = DataJSONData.getPageValue();
    String IsiDataJson = DataJSONPage.getString("CLASSOFBUSINESS");
    tempPage2.adoptJSONObject(IsiDataJson);
  }
} catch(InvalidStreamError e){
  oLog.error("ReloadSection:Invalid JSON Stream for data page params : "+e.getMessage());
} catch(Exception e){
  oLog.error("ReloadSection:Expection : "+e.getMessage());
}
```

- **6** `Property-Set`
  - syarat `param.viewstate==1`: benar → ; salah → lewati langkah [kode -/3]
  - `TreatyIn.ViewState` = `1`
  - `TreatyIn.IsEditData` = `1`
- **7** `Property-Set` — RevisionState
  - syarat `param.revisionstate==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `TreatyIn.RevisionState` = `1`
  - `TreatyIn.StatusAkseptasi` = `""`
  - `TreatyIn.PositionUsername` = `OperatorID.pxInsName`
  - `Local.update` = `1`
- **8** `RDB-List` — get date
  - syarat `param.revisionstate==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`RNM`; ClassName=`ASM-FW-GISFW-Int-TREATY_IN`; BrowsePage=`time`; RequestType=`GetCurrentDate`
- **9** `Property-Set` — add coment revisi
  - syarat `param.revisionstate==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `TreatyIn.CommentList(<APPEND>).Date` = `@toDateTime(time.pxResults(1).CARI1)`
  - `TreatyIn.CommentList(<LAST>).OperatorName` = `OperatorID.pxInsName`
  - `TreatyIn.CommentList(<LAST>).Suggest` = `"Create Revision"`
- **10** `Property-Set` · halaman `TreatyIn` — Insert to Json
  - syarat `param.revisionstate==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `InputParam.DATAPEGA` = `@ASM.GetPageJSONString()`
  - `InputParam.IDPEGA` = `TreatyIn.ID`
- **11** `RDB-List`
  - syarat `param.revisionstate==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; ClassName=`ASM-FW-GISFW-Int-TREATY_IN`; Access=`ASM`; RequestType=`SaveTreatyIn`
- **12** `Call ConvertHistoryDate`
- **13** `Call TreatySetReinstatement`
  - syarat `TreatyIn.ProportionType=="NonProportional"`: benar → ; salah → lewati langkah [kode -/3]
  - `TreatyIn.ProportionType`
- **14** `Call CheckDuplicateOffer`
- **15** `(kosong)` · halaman `TreatyIn.Limits`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  -  = `TreatyIn.Limits.Detail.`
  - **15.1** `(kosong)` · halaman `.Detail`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **15.1.1** `Property-Set`
      - syarat `@LengthOfPageList(.AchievementLists)>0`: benar → lanjut; salah → lewati langkah [kode 2/3]
      - `FlagExcel.CARI1` = `1`

### 6 · `SetValidateInstallment_Act`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · 56.786 B

- **1** `Page-Clear-Messages` · halaman `pyWorkPage.PolicyTreatyIn` — Init remove error message
- **2** `property-Set` — Init Property-Set-Message
  - `local.err` = `"Installment percentage is more than 100"`
- **3** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.ListInstallment` — Loop all installment, hitung ulang value nya, dan juga hitung total percent
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **3.1** `Property-Set` — Hitung Ulang Value
    - `.Premium` = `pyWorkPage.PolicyTreatyIn.BalanceDueTo* @divide(.InstallmentPercentage,100,4)`
    - `.PaymentTotal` = `pyWorkPage.PolicyTreatyIn.BalanceDueTo* @divide(.InstallmentPercentage,100,4)`
  - **3.2** `Property-Set` — Hitung Total Percentage
    - `local.pcttotal` = `local.pcttotal + .InstallmentPercentage`
- **4** `Page-Set-Messages` — Set Error Message when [local.pcttotal>100]
  - syarat `local.pcttotal>100`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Message=`local.err`; ClassOfPage=`ASM-FW-GISFW-Data-PolicyTreatyIn`; Page=`pyWorkPage.PolicyTreatyIn`

### 6 · `SetValue_Act`

Kelas `ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM` · 38.718 B · parameter: `SOB`, `CLASSOFBUSINESSID`, `TREATYYEAR`, `LIMITCURRENCY`, `CEDING`, `TREATYGROUP`, `TREATYTYPE`, `ID`

- **1** `call ASM-FW-GISFW-Work-NB.InputPolicyTreatyInDetail_preACT`
- **2** `Obj-Save` · halaman `pyWorkPage`
  - parameter: WithErrors=`true`; OnlyIfNew=`false`; WriteNow=`true`

### 6 · `SetValueRetro_Act`

Kelas `ASM-FW-GISFW-Int-TREATYOUTDETAIL` · 38.590 B · parameter: `SOB`, `CLASSOFBUSINESSID`, `TREATYYEAR`, `LIMITCURRENCY`, `CEDING`, `TREATYGROUP`, `TREATYTYPE`, `ID`

- **1** `call ASM-FW-GISFW-Work-NB.InputPolicyTreatyOutDetail_preACT` — call ASM-FW-GISFW-Work-NB.InputPolicyTreatyOutDetail_preACT
- **2** `Obj-Save` · halaman `pyWorkPage`
  - parameter: WithErrors=`true`; OnlyIfNew=`false`; WriteNow=`true`

### 6 · `TreatyInInputVis`

Kelas `Data-Portal` · 40.110 B · parameter: `Add`, `bypasscondition`

- **1** `Property-Set` — When Add=1, exit after step run
  - syarat `param.Add==1`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - transisi `1==1`: benar → keluar activity; salah → lanjut [kode 6/2]
  - `OutputParam.DATASHOW` = `1`
  - `OutputParam.ERRMSG` = `""`
  - `TreatyIn.ID` = `"UnknownId"`
- **2** `Property-Set` — When Add=0
  - syarat `param.Add==0`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `OutputParam.DATASHOW` = `0`
  - `TreatyIn.ProportionType` = `""`
  - `TreatyIn.ID` = `""`
- **3** `Page-Remove`
  - `Page=TreatyIn`

### 6 · `TreatyInNonSetTotal`

Kelas `Data-Portal` · 99.012 B

- **1** `Property-Remove` — remove value old total property
  - `Property=TreatyIn.TotalRetentionAmount`
  - `Property=TreatyIn.TotalEgnpiAmount`
  - `Property=TreatyIn.TotalEgnpiProportion`
  - `Property=TreatyIn.TotalLimitsAdjPct`
  - `Property=TreatyIn.TotalLimitsDeductible`
  - `Property=TreatyIn.TotalLimitsIOOLimit`
  - `Property=TreatyIn.TotalLimitsMdp`
  - `Property=TreatyIn.TotalLimitsPremiumEarned`
  - `Property=TreatyIn.TotalShareGross`
  - `Property=TreatyIn.TotalShareNet`
  - `Property=TreatyIn.TotalShareRnmLimit`
  - `Property=TreatyIn.TotalInstallmentAmount`
  - `Property=TreatyIn.TotalInstallmentPct`
- **2** `Property-Set`
  - `TreatyIn.Currency` = `TreatyIn.Retention(1).Currency`
  - `TreatyIn.CurrencyEgnpiAmount` = `"IDR"`
- **3** `Property-Set` · halaman `TreatyIn.Retention` — retention
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `TreatyIn.TotalRetentionAmount` = `TreatyIn.TotalRetentionAmount + .Amount`
- **4** `Property-Set` · halaman `TreatyIn.EGNPI` — egnpi total amount
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `TreatyIn.TotalEgnpiAmount` = `TreatyIn.TotalEgnpiAmount + .AmountIDR`
- **5** `Property-Set` · halaman `TreatyIn.EGNPI` — Egnpi Proportion
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `.Proportion` = `(.AmountIDR / TreatyIn.TotalEgnpiAmount) * 100`
  - `TreatyIn.TotalEgnpiProportion` = `TreatyIn.TotalEgnpiProportion + .Proportion`
- **6** `Property-Set` · halaman `TreatyIn.Limits` — limits
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `TreatyIn.TotalLimitsIOOLimit` = `TreatyIn.TotalLimitsIOOLimit + .Limit`
  - `TreatyIn.TotalLimitsAdjPct` = `TreatyIn.TotalLimitsAdjPct + .AdjRate`
  - `TreatyIn.TotalLimitsDeductible` = `TreatyIn.TotalLimitsDeductible + .Deductible`
  - `TreatyIn.TotalLimitsPremiumEarned` = `TreatyIn.TotalLimitsPremiumEarned + .PremiumEarned`
  - `TreatyIn.TotalLimitsMdp` = `TreatyIn.TotalLimitsMdp + .MDP`
  - `TreatyIn.TotalLimitsROL` = `(TreatyIn.TotalLimitsPremiumEarned / TreatyIn.TotalLimitsIOOLimit) * 100`
- **7** `Property-Set` · halaman `TreatyIn.Share` — share
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `TreatyIn.TotalShareRnmLimit` = `TreatyIn.TotalShareRnmLimit + .NusareLimit`
  - `TreatyIn.TotalShareGross` = `TreatyIn.TotalShareGross + .GrossPremium`
  - `TreatyIn.TotalShareNet` = `TreatyIn.TotalShareNet + .NetPremium`
- **8** `Property-Set` · halaman `TreatyIn.Installment` — installment
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - `TreatyIn.TotalInstallmentAmount` = `TreatyIn.TotalInstallmentAmount + .Amount`
  - `TreatyIn.TotalInstallmentPct` = `TreatyIn.TotalInstallmentPct + .InstallmentPct`

### 6 · `TreatyInputPctCommSpreading`

Kelas `ASM-FW-GISFW-Work` · 64.328 B

- **1** `Call FetchMasterTreatyIn` · halaman `pyWorkPage` — Set Master Treaty In to PolicyTreatyIn.TreatyIn page (Call master data for setting commision OGP & Spreading)
- **2** `(kosong)` · halaman `pyWorkPage.TreatyIn.Limits` — Set Comm% and Spreading
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **2.1** `(kosong)`
    - ulang: pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefLimit=`1`, pyStepsRepeatDefStart=`1`
    - syarat `pyWorkPage.PolicyTreatyIn.TreatyType==.TreatyType`: benar → lanjut; salah → lewati langkah [kode 2/3]
    - **2.1.1** `(kosong)` · halaman `.Detail`
      - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
      - **2.1.1.1** `Property-Set`
        - syarat `.TreatyGroup==pyWorkPage.PolicyTreatyIn.TreatyGroupName`: benar → lanjut; salah → lewati langkah [kode 2/3]
        - `pyWorkPage.PolicyTreatyIn.RiCommOgp` = `@replaceAll(.RIOGR,",",".")`
        - `pyWorkPage.PolicyTreatyIn.RiCommOgp` = `@replaceAll(.RIONR,",",".")`
        - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).SharePercentage` = `.SpreadingTotalPct`
        - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).TreatyType` = `.SpreadingTypeID`
        - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).TreatyName` = `.SpreadingType`
- **3** `Call BreakDownSpreading_Act` · halaman `pyWorkPage`
- **4** `Call CountSpreading_Act` · halaman `pyWorkPage.SpreadingRiskList(1)CountSpreading_Act` · label `//`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`

### 6 · `TreatyNonPropOutSetSpreading`

Kelas `ASM-FW-GISFW-Work` · 133.501 B · parameter: `spreading`

- **1** `Page-Clear-Messages` · halaman `pyWorkPage`
- **2** `Property-Set` — Set Spreading based on Master Data (When No Retro List)
  - syarat `.TreatyIn.FacultativeShare==0`: benar → ; salah → lewati langkah [kode -/3] *(When langkah tidak dicentang)*
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).TreatyType` = `Param.spreading`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).Currency` = `pyWorkPage.PolicyTreatyIn.Currency`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).CurrencyID` = `pyWorkPage.PolicyTreatyIn.IDCurrency`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).PremiumSpreaded` = `local.netpremi`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).SharePercentage` = `"100"`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
  - `local.sprderror` = `"Spreading in Master data is incomplete"`
- **3** `Property-Set` — Set Spreading based on Master Data (When No Retro List)
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).TreatyType` = `"ORS"`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).Currency` = `pyWorkPage.PolicyTreatyIn.Currency`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).CurrencyID` = `pyWorkPage.PolicyTreatyIn.IDCurrency`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).PremiumSpreaded` = `local.netpremi`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).SharePercentage` = `"100"`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
  - `local.sprderror` = `"Spreading in Master data is incomplete"`
- **4** `Page-Set-Messages` — Set Error when spreading ID is empty (When No Retro List)
  - syarat `Param.spreading==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - syarat `.TreatyIn.FacultativeShare==0`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Category=`pyMessageLabel`; Message=`local.sprderror`; ClassOfPage=`ASM-FW-GISFW-Work-NB`; Page=`pyWorkPage`
- **5** `Property-Remove` · label `//` — Reset .PolicyTreatyIn.SpreadingRiskList
  - syarat `.TreatyIn.FacultativeShare<1`: benar → lompat ke label `jmp`; salah → keluar activity [kode 1/6]
  - `Property=.PolicyTreatyIn.SpreadingRiskList`
- **6** `(kosong)` · halaman `Primary.TreatyIn.RetroList` · label `//` — Loop per Retro List
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **6.1** `Property-Set` — Set to Local
    - `local.RetroTypeID` = `.RetroTypeID`
    - `local.RetroPct` = `@divide(.RetroPct, Primary.TreatyIn.RNMShare,20)*100`
  - **6.2** `(kosong)` · halaman `.TotalNet`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **6.2.1** `Property-Set`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<APPEND>).TreatyType` = `local.RetroTypeID`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).Currency` = `.Currency`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).CurrencyID` = `.CurrencyID`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).SharePercentage` = `local.RetroPct`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).PremiumSpreaded` = `.Value`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).ClaimSpreaded` = `0`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).ClaimPercentage` = `0`
    - **6.2.2** `(kosong)` — When CurrencyID is empty, find from database
      - ulang: pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
      - **6.2.2.1** `Property-Set`
        - `InputData.CARI1` = `.Currency`
      - **6.2.2.2** `RDB-LIST`
        - parameter: RunInParallel=`false`; ApplyDeclaratives=`false`; ClassName=`ASM-FW-GISFW-Int-TREATY_IN`; Access=`ASM`; BrowsePage=`TempResult`; RequestType=`GetCurrencyIDByName`
      - **6.2.2.3** `Property-Set`
        - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).CurrencyID` = `TempResult.pxResults(1).CARI1`
- **7** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.SpreadingRiskList` · label `jmp` — Set Parameter for each spreading in [SpreadingRiskList]
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **7.1** `Property-Set` — Set Parameter for Calculate Spreading Value
    - `Param.SharePct` = `.SharePercentage`
    - `Param.Index` = `.pxListSubscript`
    - `Param.Type` = `"Premium"`
  - **7.2** `Call CountSpreading_Act` · halaman `pyWorkPage.PolicyTreatyIn` — SpreadingRiskList

### 6 · `TreatyNonPropSetSpreading`

Kelas `ASM-FW-GISFW-Work` · 177.927 B

- **1** `Property-Remove`
  - `Property=pyWorkPage.PolicyTreatyIn.SpreadingRiskList`
- **2** `Property-Set`
  - `InputData.CARI44` = `.TreatyIn.Share(1).SpreadingListXOL(1).ReinsTypeName`
- **3** `Property-Set` — Set Spreading based on Master Data (When No Retro List)
  - syarat `pyWorkPage.TreatyIn.Share(1).SpreadingTypeXOL!=""`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).TreatyType` = `pyWorkPage.TreatyIn.Share(1).SpreadingTypeIDXOL`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).Currency` = `pyWorkPage.PolicyTreatyIn.Currency`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).CurrencyID` = `pyWorkPage.PolicyTreatyIn.IDCurrency`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).PremiumSpreaded` = `local.netpremi`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).SharePercentage` = `"100"`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
  - `local.sprderror` = `"Spreading in Master data is incomplete"`
- **4** `(kosong)` · halaman `.TreatyIn.Share(1).SpreadingListXOL` — Set Spreading based on Master Data (When No Retro List) spreading baru
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - syarat `pyWorkPage.TreatyIn.Share(1).SpreadingTypeXOL==""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - **4.1** `Property-Set` — pyWorkpage.TreatyIn.Share(1).SpreadingListXOL
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<APPEND>).TreatyType` = `.ReinsTypeID`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).Currency` = `pyWorkPage.PolicyTreatyIn.Currency`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).CurrencyID` = `pyWorkPage.PolicyTreatyIn.IDCurrency`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).PremiumSpreaded` = `local.netpremi`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).SplitRNMSharePct` = `.Pct`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).SharePercentage` = `@divide(.Pct,pyWorkPage.TreatyIn.RNMShare,20)*100`
    - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(<LAST>).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
    - `local.sprderror` = `"Spreading in Master data is incomplete"`
- **5** `Property-Set` — set spreading jika treatyout non prop(xol retro)
  - syarat `pyWorkPage.PolicyTreatyIn.ClaimType=="XOL Retro"`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).TreatyType` = `pyWorkPage.TreatyIn.ShareReins(1).ReinsuranceListTONP(1).SpreadingTypeIDXOL`
- **6** `Property-Set` — Set Spreading based on Master Data (When No Retro List)
  - syarat `pyWorkPage.PolicyTreatyIn.FlagRetroTreaty==true`: benar → ; salah → lewati langkah [kode -/3]
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).TreatyType` = `"ORS"`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).Currency` = `pyWorkPage.PolicyTreatyIn.Currency`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).CurrencyID` = `pyWorkPage.PolicyTreatyIn.IDCurrency`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).PremiumSpreaded` = `local.netpremi`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).SharePercentage` = `"100"`
  - `pyWorkPage.PolicyTreatyIn.SpreadingRiskList(1).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
  - `local.sprderror` = `"Spreading in Master data is incomplete"`
- **7** `Page-Set-Messages` — Set Error when spreading ID is empty (When No Retro List)
  - syarat `pyWorkPage.TreatyIn.Share(1).SpreadingTypeIDXOL = ""`: benar → lanjut; salah → lewati langkah [kode 2/3]
  - syarat `.TreatyIn.FacultativeShare==0`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Category=`pyMessageLabel`; ClassOfPage=`ASM-FW-GISFW-Work-NB`; Message=`local.sprderror`; ContainingClassOfProperty=`@baseclass`; Page=`pyWorkPage`
- **8** `Property-Remove` · label `//` — Reset .PolicyTreatyIn.SpreadingRiskList
  - syarat `.TreatyIn.FacultativeShare<1`: benar → lompat ke label `jmp`; salah → keluar activity [kode 1/6]
  - `Property=.PolicyTreatyIn.SpreadingRiskList`
- **9** `(kosong)` · halaman `Primary.TreatyIn.RetroList` · label `//` — Loop per Retro List
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **9.1** `Property-Set` — Set to Local
    - `local.RetroTypeID` = `.RetroTypeID`
    - `local.RetroPct` = `@divide(.RetroPct, Primary.TreatyIn.RNMShare,20)*100`
  - **9.2** `(kosong)` · halaman `.TotalNet`
    - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
    - **9.2.1** `Property-Set`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<APPEND>).TreatyType` = `local.RetroTypeID`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).Currency` = `.Currency`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).CurrencyID` = `.CurrencyID`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).SharePercentage` = `local.RetroPct`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).pxObjClass` = `"ASM-FW-GISFW-Data-SpreadingRisk"`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).PremiumSpreaded` = `.Value`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).ClaimSpreaded` = `0`
      - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).ClaimPercentage` = `0`
    - **9.2.2** `(kosong)` — When CurrencyID is empty, find from database
      - ulang: pyStepsRepeatDefIteration=`1`, pyStepsRepeatDefHasRepeat=`REPEAT`, pyStepsRepeatDefStart=`1`, pyStepsRepeatDefLimit=`1`
      - **9.2.2.1** `Property-Set`
        - `InputData.CARI1` = `.Currency`
      - **9.2.2.2** `RDB-LIST`
        - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; Access=`ASM`; ClassName=`ASM-FW-GISFW-Int-TREATY_IN`; BrowsePage=`TempResult`; RequestType=`GetCurrencyIDByName`
      - **9.2.2.3** `Property-Set`
        - `Primary.PolicyTreatyIn.SpreadingRiskList(<LAST>).CurrencyID` = `TempResult.pxResults(1).CARI1`
- **10** `(kosong)` · halaman `pyWorkPage.PolicyTreatyIn.SpreadingRiskList` · label `jmp` — Set Parameter for each spreading in [SpreadingRiskList]
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **10.1** `Property-Set` — Set Parameter for Calculate Spreading Value
    - `Param.SharePct` = `.SharePercentage`
    - `Param.Index` = `.pxListSubscript`
    - `Param.Type` = `"Premium"`
  - **10.2** `Call CountSpreading_Act` · halaman `pyWorkPage.PolicyTreatyIn`

### 6 · `TreatyRealizationCheckDuplicate`

Kelas `ASM-FW-GISFW-Work` · 74.653 B

- **1** `Page-Clear-Messages` · halaman `pyWorkPage`
- **2** `Property-Set` — Init RDB
  - `InputData.CARI1` = `.PolicyTreatyIn.NoOffer`
  - `InputData.CARI2` = `.PolicyTreatyIn.StartDate`
  - `InputData.CARI3` = `.PolicyTreatyIn.EndDate`
  - `InputData.CARI4` = `.PolicyTreatyIn.CedingCoName`
  - `InputData.CARI5` = `.PolicyTreatyIn.TreatyYear`
  - `InputData.CARI6` = `.PolicyTreatyIn.Currency`
  - `InputData.Totaltsi` = `pyWorkPage.PolicyTreatyIn.BalanceDueTo`
  - `InputData.CARI7` = `.PolicyTreatyIn.BalanceDueTo`
- **3** `RDB-List` — Call RDB,
  - parameter: ApplyDeclaratives=`false`; RunInParallel=`false`; ClassName=`ASM-FW-GISFW-Work`; Access=`RNM`; BrowsePage=`ResultData`; RequestType=`TreatyRealizationCheckDuplicate`
- **4** `Property-Set`
  - `local.msg` = `"Protect Duplicate Policy; data is similar to "`
- **5** `(kosong)` · halaman `ResultData.pxResults`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **5.1** `Property-Set`
    - `local.msg` = `local.msg + .CARI1 + " "`
- **6** `Page-Set-Messages`
  - syarat `@SizeOfPropertyList(ResultData.pxResults) > 0 && .PolicyTreatyIn.ClaimType != "XOL"`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: Category=`pyMessageLabel`; ClassOfPage=`ASM-FW-GISFW-Work`; Message=`local.msg`; Page=`pyWorkPage`; ContainingClassOfProperty=`@baseclass`

### 6 · `TreatyRealizationCheckXOLList`

Kelas `ASM-FW-GISFW-Work` · 64.717 B

- **1** `Page-Clear-Messages` · halaman `pyWorkPage`
- **2** `(kosong)` — Check if TreatyXOLList is empty, if true continue else exit activity
  - syarat `@SizeOfPropertyList(pyWorkPage.PolicyTreatyIn.TreatyXOLList)<1`: benar → ; salah → keluar activity [kode -/6]
- **3** `Property-Set` — Init Call Activity & errormessage
  - `Param.ID` = `pyWorkPage.PolicyTreatyIn.NoOffer`
  - `local.err` = `"Error fetching XolList"`
- **4** `Call SetTreatyIn_Act` · halaman `DataPortal` — Call Treaty In
  - parameter: viewstate=`1`; ID=`3000004`
- **5** `Page-Copy` · halaman `TreatyIn`
  - parameter: CopyFrom=`TreatyIn`; CopyInto=`pyWorkPage.TreatyIn`
- **6** `Call InsertToTreatyXOLList` — Set XOL List (NewWorkPage)
- **7** `Property-Set` — Set XOL List (OldWorkPage)
  - `pyWorkPage.PolicyTreatyIn.OldData.TreatyXOLList` = `pyWorkPage.PolicyTreatyIn.TreatyXOLList`
- **8** `Page-Set-Messages` — When fetching failed, set message to page.
  - syarat `@SizeOfPropertyList(pyWorkPage.PolicyTreatyIn.TreatyXOLList)<1`: benar → ; salah → lewati langkah [kode -/3]
  - parameter: ClassOfPage=`ASM-FW-GISFW-Work`; Message=`local.err`; Page=`pyWorkPage`

### 6 · `TreatySetReinstatement`

Kelas `Data-Portal` · 23.971 B

- **1** `(kosong)` · halaman `TreatyIn.Limits`
  - ulang: pyStepsRepeatDefHasRepeat=`EMBEDDED`
  - **1.1** `Call SetReinstatementPct`
    - syarat `.Reinstatement_List(1).ReinstatementValue == ""`: benar → ; salah → lewati langkah [kode -/3]

## 7 · Data Transform

### 7 · `AddToListCommentsPolicyTreatyIn_DT`

Kelas `ASM-FW-GISFW-Work`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | APPEND_AND_MAP_TO | `.PolicyTreatyIn.SuggestList` |  | pyRelationNameAppend=NEW_PAGE |
| 1.1 | SET | `.Suggest` | `Param.Comment` |  |
| 1.2 | WHEN | `Param.Approved != ""` |  |  |
| 1.2.1 | SET | `.IsApproved` | `Param.Approved` |  |
| 1.3 | SET | `.Date` | `Param.Date` |  |
| 1.4 | SET | `.OperatorName` | `Param.Operator` |  |

### 7 · `btnCedingCO_DT`

Kelas `ASM-FW-GISFW-Data-Quotation`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | SET | `pyWorkPage.Quotation.btnQuotation` | `CedingCo` |  |
| 2 | SET | `SearchSOB.CARI1` | `""` |  |

### 7 · `btnSOB_DT`

Kelas `ASM-FW-GISFW-Data-Quotation`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | SET | `pyWorkPage.Quotation.btnQuotation` | `SOB` |  |
| 2 | SET | `SearchSOB.CARI1` | `""` |  |

### 7 · `DeptHeadTreatyIn_UW_postDT`

Kelas `ASM-FW-GISFW-Work`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | APPLY_MODEL | `AddToListCommentsPolicyTreatyIn_DT` |  | pyClassName=ASM-FW-GISFW-Work; pyModelName=AddToListCommentsPolicyTreatyIn_DT; pyRelationNameAppend=NEW_PAGE; pyPassCurrentParameterPage=false; param=Comment=.PolicyTreatyIn.Suggest, Operator=.PolicyTreatyIn.OperatorName, Approved=.PolicyTreatyIn.IsApproved, ApprovedtoDeptHead=.PolicyTreatyIn.isApprovedtoDeptHead, Date=.PolicyTreatyIn.SuggestDate |
| 2 | WHEN | `pyWorkPage.PositionNote=="ReasTreatyInAdmin" && pyWorkPage.PolicyTreatyIn.IsApproved == "0"` |  |  |
| 2.1 | SET | `.NBStatus` | `"NB WAS DECLINED BY  "+ @toUpperCase(OperatorID.pyUserName)` |  |
| 3 | WHEN | `pyWorkPage.PositionNote=="ReasTreatyInDirector" && pyWorkPage.PolicyTreatyIn.IsApproved == "0"` |  |  |
| 3.1 | SET | `.NBStatus` | `"NB IS IN "+@toUpperCase(pyWorkPage.pxCreateOpName)+"'S INBOX"` |  |
| 4 | WHEN | `pyWorkPage.PositionNote=="ReasTreatyInGroupLeader" && pyWorkPage.PolicyTreatyIn.IsApproved == "0"` |  | pyDisabled=false |
| 4.1 | SET | `.NBStatus` | `"NB IS IN "+@toUpperCase(pyWorkPage.pxCreateOpName)+"'S INBOX"` | pyDisabled=false |
| 5 | WHEN | `pyWorkPage.PositionNote=="ReasTreatyInDeptHead" && pyWorkPage.PolicyTreatyIn.IsApproved == "0"` |  | pyDisabled=false |
| 5.1 | SET | `.NBStatus` | `"NB IS IN "+@toUpperCase(pyWorkPage.pxCreateOpName)+"'S INBOX"` | pyDisabled=false |
| 6 | WHEN | `pyWorkPage.PositionNote=="ReasTreatyInGroupLeader" && pyWorkPage.PolicyTreatyIn.IsApproved == "1"` |  | pyDisabled=false |
| 6.1 | SET | `.NBStatus` | `"NB IS IN <nama-4>'S INBOX"` | pyDisabled=false |
| 7 | WHEN | `pyWorkPage.PositionNote=="ReasTreatyInDirector" && pyWorkPage.PolicyTreatyIn.IsApproved == "1"` |  |  |
| 7.1 | SET | `.NBStatus` | `"NB POLICY NO : "+pyWorkPage.PolicyTreatyIn.PolicyNo` |  |

### 7 · `DeptHeadTreatyInUW_preDT`

Kelas `ASM-FW-GISFW-Work`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | SET | `.PolicyTreatyIn.MasterID` | `""` |  |
| 2 | SET | `.PolicyTreatyIn.OperatorName` | `OperatorID.pyUserIdentifier` |  |
| 3 | SET | `.PolicyTreatyIn.SuggestDate` | `@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","")` |  |
| 4 | SET | `.PolicyTreatyIn.IsApproved` | `""` |  |
| 5 | SET | `.PolicyTreatyIn.isApprovedtoDeptHead` | `""` |  |
| 6 | SET | `.PolicyTreatyIn.Suggest` | `""` |  |
| 7 | SET | `pyWorkPage.PolicyTreatyIn.QuotationData` | `pyWorkPage.Quotation` |  |

### 7 · `InboxPolicyTreatyIn_postDT`

Kelas `ASM-FW-GISFW-Work`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | APPLY_MODEL | `AddToListCommentsPolicyTreatyIn_DT` |  | pyClassName=ASM-FW-GISFW-Work; pyModelName=AddToListCommentsPolicyTreatyIn_DT; pyPassCurrentParameterPage=false; param=Comment=.PolicyTreatyIn.Suggest, Operator=.PolicyTreatyIn.OperatorName, Approved=.PolicyTreatyIn.IsApproved, ApprovedtoDeptHead=.PolicyTreatyIn.isApprovedtoDeptHead, Date=.PolicyTreatyIn.SuggestDate |
| 2 | SET | `.PolicyTreatyIn.HasFacOut` | `"0"` |  |
| 3 | APPLY_MODEL | `TestTreatyToFacStatus` |  |  |
| 4 | WHEN | `pyWorkPage.PositionNote=="ReasTreatyInAdmin" && pyWorkPage.PolicyTreatyIn.IsApproved == "0"` |  |  |
| 4.1 | SET | `.NBStatus` | `"NB WAS DECLINED BY  "+ @toUpperCase(OperatorID.pyUserName)` |  |

### 7 · `InputPolicyTreatyIn_preDT`

Kelas `ASM-FW-GISFW-Work`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | WHEN | `.PolicyTreatyIn.StartDate==""` |  |  |
| 1.1 | SET | `.PolicyTreatyIn.StartDate` | `@DateTime.CurrentDate("dd/MM/yyyy","")` |  |
| 2 | WHEN | `.PolicyTreatyIn.EndDate==""` |  |  |
| 2.1 | SET | `.PolicyTreatyIn.EndDate` | `@DateTime.CurrentDate("dd/MM/yyyy","")` |  |
| 3 | WHEN | `.PolicyTreatyIn.StatementDate==""` |  |  |
| 3.1 | SET | `.PolicyTreatyIn.StatementDate` | `@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","")` |  |
| 4 | WHEN | `.Quotation.ProportionalType=="NonProportional"` |  |  |
| 4.1 | SET | `.PolicyTreatyIn.IsNewPolicyNonProp` | `"1"` |  |
| 5 | WHEN | `.PolicyTreatyIn.IsNewPolicyNonProp != "1"` |  |  |
| 5.1 | SET | `.PolicyTreatyIn.IsNewPolicyNonProp` | `"0"` |  |
| 6 | SET | `.PolicyTreatyIn.MasterID` | `""` |  |
| 7 | SET | `.PolicyTreatyIn.MarketingOfficer` | `pyWorkPage.Quotation.MarketingName` |  |
| 8 | SET | `.PolicyTreatyIn.Suggest` | `""` |  |
| 9 | SET | `.PolicyTreatyIn.SuggestDate` | `@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","")` |  |
| 10 | SET | `.PolicyTreatyIn.OperatorName` | `OperatorID.pyUserName` |  |
| 11 | SET | `pyWorkPage.FlagOnGoingPolicy` | `"1"` |  |
| 12 | SET | `TempEmail.CARI28` | `OperatorID.pyWorkBasketList(2).pyWorkBasketName` |  |
| 13 | WHEN | `pyWorkPage.PositionNote==""` |  |  |
| 13.1 | SET | `pyWorkPage.PositionNote` | `OperatorID.pyWorkBasketList(2).pyWorkBasketName` |  |
| 14 | SET | `pyWorkPage.PolicyTreatyIn.QuotationData` | `pyWorkPage.Quotation` |  |

### 7 · `SearchHierarkiSourceBizAgent_PostDT`

Kelas `ASM-FW-GISFW-Data-Agent`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | WHEN | `pyWorkPage.Quotation.btnQuotation=="SOB"` |  |  |
| 1.1 | SET | `pyWorkPage.Quotation.SourceOfBusiness` | `@if(.ChildCount > 0, "", .ID)` |  |
| 1.2 | SET | `pyWorkPage.Quotation.SobName` | `@if(.ChildCount > 0, "", .ClientName)` |  |
| 2 | OTHERWISE_WHEN | `pyWorkPage.Quotation.btnQuotation=="CedingCo"` |  |  |
| 2.1 | SET | `pyWorkPage.Quotation.CedingCo` | `@if(.ChildCount > 0, "", .ID)` |  |
| 2.2 | SET | `pyWorkPage.Quotation.CedingCoName` | `@if(.ChildCount > 0, "", .ClientName)` |  |
| 3 | OTHERWISE_WHEN | `pyWorkPage.Quotation.btnQuotation=="CedingCedant"` |  |  |
| 3.1 | SET | `pyWorkPage.OfferFacIn.CedingCedantID` | `@if(.ChildCount > 0, "", .ID)` |  |
| 3.2 | SET | `pyWorkPage.OfferFacIn.CedingCedant` | `@if(.ChildCount > 0, "", .ClientName)` |  |
| 4 | SET | `pyWorkPage.Quotation.SobLeader0` | `@if(.ChildCount > 0, "", .Leader0)` |  |
| 5 | SET | `pyWorkPage.Quotation.SobLeader1` | `@if(.ChildCount > 0, "", .Leader1)` |  |

### 7 · `setCategoryAttachment_DT`

Kelas `ASM-FW-GISFW-Work-LIFE`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | SET | `Attachment.pxResults(1).NOTE` | `"LIFE"` |  |

### 7 · `SystemSetOneYear_DT`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | SET | `.EndDate` | `@DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)` |  |

### 7 · `TestTreatyToFacStatus`

Kelas `ASM-FW-GISFW-Work`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | FOR_EACH_PAGE_IN | `.PolicyTreatyIn.SpreadingRiskList` |  | pyUpdateSourceContext=false |
| 1.1 | WHEN | `.TreatyType=="10015"` |  |  |
| 1.1.1 | SET | `Primary.PolicyTreatyIn.HasFacOut` | `"1"` |  |
| 1.2 | WHEN | `.TreatyType=="10218"` |  |  |
| 1.2.1 | SET | `Primary.PolicyTreatyIn.HasFacOut` | `"1"` |  |

### 7 · `TreatyEnableDisableInput`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`

| No | Aksi | Sasaran | Sumber | Lain |
| --- | --- | --- | --- | --- |
| 1 | SET | `.QuotationData.ProportionalType` | `"NonProportional"` |  |
| 2 | WHEN | `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp==""` |  |  |
| 2.1 | SET | `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp` | `"1"` |  |
| 3 | WHEN | `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp=="0"` |  |  |
| 3.1 | SET | `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp` | `"1"` |  |
| 4 | WHEN | `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp=="1"` |  |  |
| 4.1 | SET | `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp` | `"0"` |  |

## 8 · RDBList — naskah SQL

### 8 · `AttachmentLife`

Kelas `ASM-FW-GISFW-Int-OFFERJSON` · akses `ASM` · tabel: `CATEGORY_ATTACH_REAS` · dipakai: `SetCategoryAttach`

`pyBrowseSQL`:

```sql
select * from CATEGORY_ATTACH_REAS where note in ('OFFER','QUOTATION/PLACING SLIP','EMAIL','SOA','R/I SLIP','CLAUSES','EXTENDWPC','NOTA/INVOICE','OBJECTLIST','OTHERS','PHOTO','SURVEYREPORT')
```

### 8 · `BrowseClientEmail_SQL`

Kelas `ASM-FW-GISFW-Int-CLIENT` · akses `ASM` · tabel: `CLIENTEMAIL` · dipakai: `SearchHierarkiSourceBizAgentTreatyIn_Act`

`pyBrowseSQL`:

```sql
select id, emailaddress from clientemail
```

### 8 · `BrowseTreatyIn`

Kelas `ASM-FW-GISFW-Int-TREATY_IN` · akses `ASM` · tabel: `POOLDATA.M_TREATY_IN`, `POOLDATA.M_TREATY_IN_EDM` · dipakai: `SetTreatyIn_Act`

`pyBrowseSQL`:

```sql
select * from (
select JSONDATA as ClassofBusiness from pooldata.M_TREATY_IN where ID={TreatyIn.ID}
union all
select JSONDATA as ClassofBusiness from pooldata.M_TREATY_IN_edm where ID={TreatyIn.ID}
)
```

### 8 · `BrowseTreatyInDetailJoinEDM`

Kelas `ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM` · akses `RNM` · tabel: `POOLDATA.M_TREATY_IN_DETAIL_EDM` · dipakai: `InputPolicyTreatyInDetail_preACT`

`pyBrowseSQL`:

```sql
select JSONDATA as CLASSOFBUSINESS from pooldata.M_TREATY_IN_DETAIL_EDM where ID={TreatyIn.ID}
```

### 8 · `BrowseTreatyInJoinEDM`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · akses `RNM` · tabel: `POOLDATA.M_TREATY_IN`, `POOLDATA.M_TREATY_IN_EDM` · dipakai: `FetchMasterTreatyIn`, `InputPolicyTreatyInDetail_NonProp`

`pyBrowseSQL`:

```sql
select a.id, a.JSONDATA as CLASSOFBUSINESS from pooldata.m_treaty_in a
where a.ID = {pyWorkPage.PolicyTreatyIn.NoOffer}
union all
select b.id, b.JSONDATA as CLASSOFBUSINESS from pooldata.m_treaty_in_edm b
where b.ID = {pyWorkPage.PolicyTreatyIn.NoOffer}
```

### 8 · `BrowseTreatyOut`

Kelas `ASM-FW-GISFW-Data-PolicyTreatyIn` · akses `RNM` · tabel: `POOLDATA.M_TREATY_OUT` · dipakai: `InputPolicyTreatyOutDetail_NonProp`

`pyBrowseSQL`:

```sql
select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}
```

### 8 · `BrowseTreatyOutDetail`

Kelas `ASM-FW-GISFW-Int-TREATYOUTDETAIL` · akses `RNM` · tabel: `POOLDATA.M_TREATY_OUT` · dipakai: `InputPolicyTreatyOutDetail_preACT`

`pyBrowseSQL`:

```sql
select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}
```

### 8 · `CategoryAttach_SQL`

Kelas `ASM-FW-GISFW-Int-OFFERJSON` · akses `ASM` · tabel: `CATEGORY_ATTACH_REAS` · dipakai: `SetCategoryAttach`

`pyBrowseSQL`:

```sql
select * from CATEGORY_ATTACH_REAS order by note
```

### 8 · `CheckZipCode_SQL` ⛔

Kelas `ASM-FW-GISFW-Int-RW` · akses `RNM` · tabel: `RW` · dipakai: `ProtectCoverage_Act`, `ProtectFIREMBUPA_Act`

`pyBrowseSQL`:

```sql
select count(*) as CARI1 from rw where zipcode ={ParamInProdCase.CARI16}
```

### 8 · `FetchTreatyGroupOLDID`

Kelas `ASM-FW-GISFW-Int-TREATYGROUP` · akses `RNM` · tabel: `POOLDATA.TREATYGROUP` · dipakai: `FetchTreatyGroupOldID`

`pyBrowseSQL`:

```sql
select oldid as CARI1, treatygroupname as CARI2 from pooldata.treatygroup where ID = {InputData.CARI1}
```

### 8 · `GenerateNoPolicy`

Kelas `ASM-FW-GISFW-Int-POLISTREATYIN` · akses `ASM` · tabel: `DUAL` · dipakai: `SaveJsonPolisTreatyIn_Act`

`pyBrowseSQL`:

```sql
SELECT 'RNM-' || {InputData.CARI20} || '.T' || {pyWorkPage.PolicyTreatyIn.QuotationData.BusinessOldId} || '.' || to_char(sysdate,'MM.yyyy') || '.' || LPAD (POOLDATA.JSON_POLIS_TREATYIN_SEQ.NEXTVAL, 5, '0') AS "HASIL" FROM dual
```

### 8 · `GetAllCurrency` ⛔

Kelas `ASM-FW-GISFW-Int-CURRENCY` · akses `ASM` · tabel: `CURRENCY` · dipakai: `GetCurrencyMaster`

`pyBrowseSQL`:

```sql
select Currency as CURR, OldID as OldID from currency
```

### 8 · `GetBreakDownSpread_SQL`

Kelas `ASM-FW-GISFW-Work` · akses `RNM` · tabel: `POOLDATA.PROPORTIONALARRG` · dipakai: `BreakDownSpreading_Act`

`pyBrowseSQL`:

```sql
select DISTINCT Reinstypeid AS "TreatyType",PCT AS "SharePercentage" from pooldata.proportionalarrg a
where a.treatygroupid = {ParamData.CARI10}
and  a.PARENTREINSTYPEID ={ParamData.CARI11}
and TreatyDescID ='10001'
```

### 8 · `GetCategoryOccupation` ⛔

Kelas `ASM-FW-GISFW-Int-OCCUPATION` · akses `ASM` · tabel: `OCCUPATION` · dipakai: `ProtectFIREMBUPA_Act`, `SetFlagOccupation_ACT`

`pyBrowseSQL`:

```sql
select distinct CATEGORY AS CARI1 from OCCUPATION where OCCUPATIONCODE = {TempOccupation.CARI1}
```

### 8 · `GetClientID_SQL`

Kelas `ASM-FW-GISFW-Int-CLIENT` · akses `ASM` · tabel: `CLIENT` · dipakai: `InputPolicyTreatyInDetail_preACT`, `InputPolicyTreatyOutDetail_preACT`

`pyBrowseSQL`:

```sql
Select ID as CARI1 from Client where replace(name,' ','') = replace({SearchClient.CARI1},' ','')
```

### 8 · `GetCountClaim`

Kelas `Assign-Worklist` · akses `RNM` · tabel: `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` · dipakai: `CheckDuplicateOffer`

`pyBrowseSQL`:

```sql
select count(1) as "CARI1"  from DATAPEGA.PC_ASM_FW_GCNMFW_WORK where masterid={TreatyWarning.CAIREINSFACIN}
```

### 8 · `GetCurrency`

Kelas `ASM-FW-GISFW-Int-CURRENCY` · akses `ASM` · tabel: `CURRENCY` · dipakai: `ProtectCurrencyTSI_Act`, `SetCurrency_act`

`pyBrowseSQL`:

```sql
select CURRENCY as HASIL1 from currency where id = {InputData.CARI1}
```

### 8 · `GetCurrencyIDByName`

Kelas `ASM-FW-GISFW-Int-TREATY_IN` · akses `ASM` · tabel: `POOLDATA.CURRENCY` · dipakai: `SetTreatyCurrencyID`, `TreatyNonPropOutSetSpreading`, `TreatyNonPropSetSpreading`

`pyBrowseSQL`:

```sql
select ID AS CARI1,OLDID AS CARI2,CURRENCY AS CARI3,CURRENCYSYMBOL AS CARI4 from POOLDATA.CURRENCY where CURRENCY ={InputData.CARI1}
```

### 8 · `GetCurrentDate`

Kelas `ASM-FW-GISFW-Int-TREATY_IN` · akses `RNM` · tabel: `DUAL` · dipakai: `SetTreatyIn_Act`

`pyBrowseSQL`:

```sql
select sysdate as CARI1 from dual
```

### 8 · `GetDataAgentByNameNonLife_SQL` ⛔

Kelas `ASM-FW-GISFW-Int-AGENT` · akses `ASM` · tabel: `AGENT` · dipakai: `ProtectCedingCo`

`pyBrowseSQL`:

```sql
select ID,CLIENTID as "ClientID",CLIENTNAME as "ClientName",LEADER0 as "Leader0" from AGENT WHERE LEADER0 NOT LIKE 'L%' AND CLIENTNAME = {InputDataCredit.CARI5}
```

### 8 · `GetDataByID_SQL` ⛔

Kelas `ASM-FW-GISFW-Int-REINSURANCETYPE` · akses `ASM` · tabel: `M_REINSURANCETYPE` · dipakai: `cekSpreadingFactIn`

`pyBrowseSQL`:

```sql
select distinct a.jsondata.Note as CARI1 ,RP as CARI2 from m_reinsurancetype a, proportionalarrg b where a.jsondata.ID=b.REINSTYPEID and a.jsondata.ID={InputData.ReinstypeID} and TREATYDESCNAME='TREATY LIMIT'
```

### 8 · `GetDataCurrencyByName_SQL`

Kelas `ASM-FW-GISFW-Int-CURRENCY` · akses `ASM` · tabel: `CURRENCY` · dipakai: `InputPolicyTreatyEDMDetail_NP`, `InputPolicyTreatyInDetail_NonProp`, `InputPolicyTreatyOutDetail_NonProp`, `InsertToTreatyOutXOLList`, `InsertToTreatyXOLList`, `InsertToTreatyXOLListRetroShare`, `SumTSIPremiSpreadedRNM_ANEKA_Act`, `SumTSIPremiSpreadedRNM_FIRE_Act`, `SumTSIPremiSpreadedRNM_GOLF_Act`, `SumTSIPremiSpreadedRNM_MARINECARGO_Act`, `SumTSIPremiSpreadedRNM_MBU_Act`, `SumTSIPremiSpreadedRNM_PA_Act`

`pyBrowseSQL`:

```sql
select ID,OLDID,CURRENCY as "CURR",CURRENCYSYMBOL as "CurrencySymbol" from CURRENCY where CURRENCY ={InputDataCredit.CARI8}
```

### 8 · `GetDataDoubleCase` ⛔

Kelas `ASM-FW-GISFW-Int-policyjson` · akses `ASM` · tabel: `FACINPRODUCTION` · dipakai: `ProtectCoverage_Act`, `ProtectFIREMBUPA_Act`

`pyBrowseSQL`:

```sql
select distinct IDPEGA,NOPOLIS from facinproduction 
    where NOENDORS is null 
        AND replace(replace(INSUREDNAME,'.',''),' ','') ={ParamInProdCase.CARI10}
        AND SOB ={ParamInProdCase.CARI11}
        AND BUSINESSCODE ={ParamInProdCase.CARI12}
        AND RISLIPCEDING ={ParamInProdCase.CARI13}
        AND to_char(BEGINDATE,'dd/MM/yyyy')={ParamInProdCase.CARI14} AND to_char(ENDDATE,'dd/MM/yyyy')={ParamInProdCase.CARI15}
```

### 8 · `GetFlagReject_SQL` ⛔

Kelas `ASM-FW-GISFW-Int-policyjson` · akses `ASM` · tabel: `HISTORYAKSEPTASIPEGA` · dipakai: `Protection_Act`

`pyBrowseSQL`:

```sql
SELECT COUNT(*) AS CARI1 FROM HISTORYAKSEPTASIPEGA WHERE ID_PEGA = {DataSearch.CARI20} AND WORKBASKET != 'ReasFacInMarketing' AND WORKBASKET != 'ReasFacInTeamLeader' AND STATUS ='REJECT' ORDER BY TGL_TRANSFER DESC
```

### 8 · `GetKodeProdNonLife_SQL`

Kelas `ASM-FW-GISFW-Int-policyjson` · akses `RNM` · tabel: `POOLDATA.KODE_PRODUKSI` · dipakai: `GeneratePolicyNoTreaty_Act`

`pyBrowseSQL`:

```sql
SELECT KODE AS "ParamSeq.HASIL3" FROM POOLDATA.KODE_PRODUKSI WHERE TYPE ='NONLIFE'
```

### 8 · `GetKursLimitSpreading_SQL` ⛔

Kelas `ASM-FW-GISFW-Int-CURRENCY` · akses `ASM` · tabel: `TREATYEXCHANGEYEARLY` · dipakai: `SumTSIPremiSpreadedRNM_Act`

`pyBrowseSQL`:

```sql
select 
case  
    when ((select COUNT(rownum) from treatyexchangeyearly where IDCURRENCY = {SearchCurrencyValueIn.CARI1} and {SearchCurrencyValueIn.CARI10} between(startdate) and (enddate))<1) then
        (select toidr from (select * from treatyexchangeyearly where IDCURRENCY = {SearchCurrencyValueIn.CARI1} and  {SearchCurrencyValueIn.CARI10} between(startdate) and (enddate)order by startdate desc) where rownum = 1)
    when ((select COUNT(rownum) from treatyexchangeyearly where IDCURRENCY = {SearchCurrencyValueIn.CARI1} and {SearchCurrencyValueIn.CARI10} between(startdate) and (enddate))>=1) then
        (select toidr from (select * from treatyexchangeyearly where IDCURRENCY = {SearchCurrencyValueIn.CARI1} and  {SearchCurrencyValueIn.CARI10} between(startdate) and (enddate) order by startdate desc) where rownum = 1)
    end as CARI12
from treatyexchangeyearly where rownum = 1
```

### 8 · `GetObjectItembyName_SQL` ⛔

Kelas `ASM-FW-GISFW-Int-V_JN_OBJ_ITEM` · akses `ASM` · tabel: `V_JN_OBJ_ITEM` · dipakai: `ProtectFIREMBUPA_Act`

`pyBrowseSQL`:

```sql
SELECT MJOI_KODE AS CARI1 ,JN_OBJ_ITEM AS CARI2 FROM V_JN_OBJ_ITEM WHERE REPLACE(upper(JN_OBJ_ITEM),' ','') =REPLACE(upper({TempObjItemIn.CARI10}),' ','') AND ISACTIVE ='1'
```

### 8 · `GetOldIDBusiness_SQL`

Kelas `ASM-FW-GISFW-Int-OFFERJSON` · akses `ASM` · tabel: `BUSINESS` · dipakai: `InputPolicyTreatyInDetail_preACT`, `InputPolicyTreatyInPre_Act`, `InputPolicyTreatyOutDetail_preACT`

`pyBrowseSQL`:

```sql
select oldid as CARI1, GROUPPANEL as CARI2, ID as CARI3 from business where (ID ={OldID.CARI2} OR NOTE ={OldID.CARI2}) AND GROUPPANEL IS NOT NULL
```

### 8 · `GetSequenceNumber_SQL`

Kelas `ASM-FW-GISFW-Int-policyjson` · akses `RNM` · tabel: — · prosedur: `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` · dipakai: `GeneratePolicyNoTreaty_Act`

`pyBrowseSQL`:

```sql
BEGIN 
  
  POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER({ParamSeq.CARI1},{ParamSeq.CARI2},TO_DATE({ParamSeq.CARI3},'DD/MM/YYYY'), {ParamSeq.HASIL1 out},{ParamSeq.HASIL2 out});
  COMMIT; 
  
END;
```

### 8 · `GetSQLDate`

Kelas `ASM-FW-GISFW-Work` · akses `RNM` · tabel: `DUAL` · dipakai: `GeneratePolicyNoTreaty_Act`, `InputPolicyTreatyInPre_Act`

`pyBrowseSQL`:

```sql
select sysdate as CARI1 from dual
```

### 8 · `GETTanggalClosing_SQL`

Kelas `ASM-FW-GISFW-Int-policyjson` · akses `RNM` · tabel: `POOLDATA.TANGGAL_CLOSING` · dipakai: `GeneratePolicyNoTreaty_Act`

`pyBrowseSQL`:

```sql
SELECT * FROM POOLDATA.TANGGAL_CLOSING
```

### 8 · `GetTgl_InputJsonPolis_SQL` ⛔

Kelas `ASM-FW-GISFW-Int-policyjson` · akses `ASM` · tabel: `JSON_POLIS` · dipakai: `CheckSpreadingProtect_ACT`

`pyBrowseSQL`:

```sql
select tgl_input  as CARIDATETIME  from json_polis where nopolis ={TempNopolis.CARI1} order by tgl_input desc
```

### 8 · `GetTreatyName` ⛔

Kelas `ASM-FW-GISFW-Int-PROPORTIONALARRG` · akses `ASM` · tabel: `REINSURANCETYPE` · dipakai: —

`pyBrowseSQL`:

```sql
SELECT ID AS CARI1, NOTE AS CARI2 FROM REINSURANCETYPE WHERE ID = {InputData.CARI3}
```

### 8 · `GetTreatyName_SQL` ⛔

Kelas `ASM-FW-GISFW-Int-PROPORTIONALARRG` · akses `ASM` · tabel: `POOLDATA.REINSURANCETYPE`, `PROPORTIONALARRG`, `TREATYBUSINESS`, `TREATYCONTRACT` · dipakai: `GetTreatyName`

`pyBrowseSQL`:

```sql
select (select c.nourut from pooldata.reinsurancetype c where c.id = a.reinstypeid) as CARI5, a.REINSTYPEID as CARI1, a.REINSTYPENAME as CARI2, a.RP as CARI3, a.treatyyearid as CARI4 from PROPORTIONALARRG a
where a.TREATYDESCNAME = 'TREATY LIMIT'
        and a.treatygroupid in (select treatygroupid from treatybusiness 
                            where bizcode =  {pyWorkPage.OfferFacIn.QuotationData.BusinessCode} and ISACTIVE='1')
        and a.reinstypeid in (select reinstypeID from treatybusiness 
                            where  bizcode =  {pyWorkPage.OfferFacIn.QuotationData.BusinessCode} and ISACTIVE='1')
        and a.TREATYYEARID in (SELECT IDTREATYYEAR FROM treatycontract WHERE To_date({Track.CARI33},'DD/MM/RRRR') BETWEEN TREATYSTARTDATE AND TREATYENDDATE)
        and a.reinstypeid in (SELECT REINSTYPEID FROM treatycontract WHERE To_date({Track.CARI33},'DD/MM/RRRR') BETWEEN TREATYSTARTDATE AND TREATYENDDATE)
ORDER BY CARI5 ASC
```

### 8 · `GetTreatyName_SQL2` ⛔

Kelas `ASM-FW-GISFW-Int-PROPORTIONALARRG` · akses `ASM` · tabel: `POOLDATA.REINSURANCETYPE`, `PROPORTIONALARRG`, `TREATYBUSINESS` · dipakai: `GetTreatyName`

`pyBrowseSQL`:

```sql
select (select c.nourut from pooldata.reinsurancetype c where c.id = a.reinstypeid) as CARI5, a.REINSTYPEID as CARI1, a.REINSTYPENAME as CARI2, a.RP as CARI3, a.treatyyearid as CARI4 from PROPORTIONALARRG a
where a.TREATYDESCNAME = 'TREATY LIMIT' and REINSTYPEID ={Track.CARI44}
        and a.treatygroupid in (select treatygroupid from treatybusiness 
                            where bizcode =  {pyWorkPage.OfferFacIn.QuotationData.BusinessCode} and ISACTIVE='1')
ORDER BY CARI5 ASC
```

### 8 · `InsertHistoryAkseptasiPega_Sql`

Kelas `ASM-FW-GISFW-int-policyjson` · akses `ASM` · tabel: `HISTORYAKSEPTASIPEGA` · prosedur: `INSERT` · dipakai: `InsertHistoryAkseptasiPega`

`pyBrowseSQL`:

```sql
BEGIN

INSERT INTO HISTORYAKSEPTASIPEGA
(ID_PEGA, Tgl_Transfer, Status, Username, Workbasket,ID_KOMITE) 
VALUES
( {InsertHistory.CARI1},
sysdate,
{InsertHistory.CARI5},
{InsertHistory.CARI4},
{InsertHistory.CARI2},
{InsertHistory.CARI6}
);

COMMIT;

END;
```

### 8 · `InsertViewSuggest_SQL`

Kelas `ASM-FW-GISFW-Work` · akses `RNM` · tabel: `POOLDATA.HISTORYAKSEPTASIPRODUCTION` · prosedur: `INSERT` · dipakai: `SaveViewSuggest`

`pyBrowseSQL`:

```sql
BEGIN 
    INSERT INTO  POOLDATA.historyakseptasiproduction
               (IDPEGA      
                  ,TYPE_POLIS 
                  ,NOURUT     
                  ,POSISI     
                  ,PIC        
                  ,TGL_INP    
                  ,DIV        
                  ,TYPE       
                  ,PUTARAN    
                  ,APPROVAL   
                  ,KETERANGAN 
                  ,AKSES_LOGIN
                  ,B2B
                  ,BUSINESS_CODE
                  ,PERCENT_RNM)
     VALUES ({pyWorkPage.pzInsKey}, 
                  {InputData.CARI1}, 
                  {InputData.CARI2}, 
                  {InputData.CARI3}, 
                  {InputData.CARI4}, 
                  To_date({InputData.CARI5}, 'DD/MM/YYYY HH24:MI:SS'), 
                  {InputData.CARI6}, 
                  {InputData.CARI7}, 
                  {InputData.CARI8}, 
                  {InputData.CARI9},
                  substr({InputData.CARI10},0,3990),
                  {OperatorID.pyUserIdentifier},
                  {pyWorkPage.OfferFacIn.IsB2B},
                  {pyWorkPage.Quotation.BusinessCode},
                  {pyWorkPage.OfferFacIn.PercentShare}
                  );
    COMMIT; 
END;
```

### 8 · `SavePolisTreatyIn_SQL`

Kelas `ASM-FW-GISFW-Int-POLISTREATYIN` · akses `ASM` · tabel: — · prosedur: `POOLDATA.PEGA_JSON_POLIS_TREATYIN` · dipakai: `SaveJsonPolisTreatyIn_Act`

`pyBrowseSQL`:

```sql
BEGIN
  POOLDATA.PEGA_JSON_POLIS_TREATYIN(
      {pyWorkPage.pzInsKey},
      {pyWorkPage.PolicyTreatyIn.PolicyNo},      
      NULL,
      '0',
      {InputData.CARI21},             
      {OperatorID.pyUserIdentifier},
      {InputData.CARI3}, 
      {OutputData.HASIL1 out} 
  );
  COMMIT;
END;
```

### 8 · `SaveTreatyIn`

Kelas `ASM-FW-GISFW-Int-TREATY_IN` · akses `ASM` · tabel: — · prosedur: `POOLDATA.PEGA_TREATY_IN` · dipakai: `SetTreatyIn_Act`

`pyBrowseSQL`:

```sql
BEGIN 
    POOLDATA.PEGA_TREATY_IN({InputParam.DATAPEGA},
            {TreatyIn.ID},
            {TreatyIn.ProportionType},
            {TreatyIn.TreatyContractName},
            {TreatyIn.TeritorialScope},
            {TreatyIn.Commencement},
            {TreatyIn.Termination},
            {TreatyIn.ClassofBusiness},
            {TreatyIn.LeadingReinsSource},
            {TreatyIn.LeadingReinsSourceID},
            {TreatyIn.Ceding},
            {TreatyIn.CedingID},
            {TreatyIn.LeadingReinsID},
            {TreatyIn.NusareSharePct},
            {TreatyIn.BrokeragePct},
            {TreatyIn.Information},
            {TreatyIn.PositionUsername},
            {TreatyIn.Position},
            {TreatyIn.StatusAkseptasi},
            {TreatyIn.ChooseStatusAkseptasi},
            {TreatyIn.TreatyYear},
            {OutputParam.ERRMSG out }, {OutputParam.IDPEGAOUT out}, {OutputParam.STSSAVE out});
 COMMIT; 
END;
```

### 8 · `SelectSpreadingTreatyInProduction`

Kelas `ASM-FW-GISFW-Int-REINSURANCETYPE` · akses `RNM` · tabel: `REINSURANCETYPE` · dipakai: `InputPolicyTreatyInPre_Act`

`pyBrowseSQL`:

```sql
select ID as CARI1, NOTE as CARI2 from reinsuranceType where note like '%TRT%' or note = 'ORS' or note = 'FAC-OUT' or note ='QS (OR)' or note ='QS (R/I)'
```

### 8 · `TreatyRealizationCheckDuplicate`

Kelas `ASM-FW-GISFW-Work` · akses `RNM` · tabel: `POOLDATA.TREATYINPRODUCTION` · dipakai: `TreatyRealizationCheckDuplicate`

`pyBrowseSQL`:

```sql
select nopolis as CARI1 from pooldata.treatyinproduction where 
NOOFFER = {InputData.CARI1} and 
BEGINDATE = TO_DATE({InputData.CARI2}, 'YYYYMMDD') and
ENDDATE = TO_DATE({InputData.CARI3}, 'YYYYMMDD')and
INSUREDNAME=  {InputData.CARI4}and 
UW_YEAR = {InputData.CARI5}and
CURR_ID = {InputData.CARI6}and
BALANCE_DUE_TO = REPLACE({InputData.Totaltsi},',','.')
```

## 9 · Report Definition

| Nama | Kelas | Kolom | Filter | Dipakai |
| --- | --- | --- | --- | --- |
| `BrowseAgentHierarkiList_RD` | ASM-FW-GISFW-Int-AGENT | .ID, .ClientName, .Leader0, .ChildCount, .ClientID | A: .Leader0 = Param.Leader; B: .StatusActive IS NULL | AgentSourceBizTreatyIn_Act |
| `BrowseAgentNusaRe_RD` | ASM-FW-GISFW-Int-AGENT | .ID, .ClientID, .ClientName, .CountryID, .CountryName, .BU_ID, .Leader0, .ChildCount, .StatusActive, .UpdateUser, .UpdateTime, .Cause, .OtherCause, .Remarks, .RemarkStatus, .AgentType2, .STS_PKP | A: .StatusActive = 1; B: .AgentType2 != "LIFE INSURANCE"; D: .ChildCount = Param.ChildCount; E: .ClientName Contains Param.ClientName; G: .ClientID IS NOT NULL | SourceHierarki |
| `BrowseCedingCo_RD` | ASM-FW-GISFW-Int-AGENT | .ID, .ClientName, .Leader0, .ChildCount | A: .Leader0 = Param.CedingCoLeader; B: .ID = Param.ID; C: .StatusActive IS NULL | AgentSourceBizTreatyIn_Act |
| `BrowseClientName_RD` | ASM-FW-GISFW-Int-AGENT | .ID, .ClientName, .Leader0, .ChildCount, .CountryID, .AgentType2, .STS_PKP | A: .ClientName = Param.ClientName; B: .CountryID = Param.CountryID; C: .AgentType2 = Param.AgentType2; D: .ID = Param.ID; E: .StatusActive IS NULL | SetPPNPPH |
| `BrowseCurrency_RD` ⛔ | ASM-FW-GISFW-Int-CURRENCY | .ID, .CountryID, .Currency, .CurrencySymbol, .CountryName, .ISOSymbol, .Note, .OLDID | A: .Currency = Param.Currency; B: .ID = Param.ID; C: .Currency != "ITL" | GetCurrencyMaster |
| `BrowseCurrencyTreatyIn_RD` | ASM-FW-GISFW-Int-CURRENCY | .ID, .CountryID, .Currency, .CurrencySymbol, .CountryName, .ISOSymbol, .Note, .OLDID | A: .Currency = Param.Currency; B: .Currency != "ITL" | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `BrowseMarketingOfficer_RD` | ASM-FW-GISFW-Int-marketingofficer | .ID, .ClientID, .BranchDetailID, .BranchDetailName, .MOLeader, .MOStatus, .ClientName, .ClientID2, .BranchStatus, .TeamGroup, .Tanggal, .UserUpdate | A: .MOStatus = 1 | DetailPolicyTreatyIn |
| `BrowseReinsuranceType_RD` | ASM-FW-GISFW-Int-REINSURANCETYPE | .ID, .Note, .Type, .SOANote, .Code, .UserID, .Flag | A: .ID = Param.ID; B: .Note = Param.Note; C: .Flag = Param.Flag; D: .Type = Param.Type | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `BrowseTREATY_IN` | ASM-FW-GISFW-Int-TREATY_IN | .ID, .ProportionType, .TeritorialScope, .Commencement, .Termination, .ClassofBusiness, .Ceding, .CedingID, .TreatyContractName, .Information, .Position, .PositionUsername, .StatusAkseptasi, .ChooseStatusAkseptasi, .LeadingReinsSourceID, .LeadingReinsSource | A: .ID Contains Param.ID; H: .ProportionType = Param.ProportionType; B: .CedingID = Param.CedingID; F: .Ceding Contains Param.Ceding; C: .LeadingReinsSourceID = Param.LeadingReinsSourceID; G: .LeadingReinsSource Contains Param.LeadingReinsSource; D: .Commencement = Param.Commencement; E: .Terminatio … | CheckDuplicateOffer |
| `BrowseTreatyGroup_RD` | ASM-FW-GISFW-Int-TREATYGROUP | .ID, .TreatyGroupName, .TreatyGroupSOAName, .OJKBusinessID, .OJKBusinessName, .OJKBusinessNameIDN | A: .ID = Param.ID | FetchTreatyGroupOJK |
| `BrowseTreatyInDetail` | ASM-FW-GISFW-Int-TREATYINDETAIL | .ID, .TREATYID, .TREATYCONTRACTNAME, .PROPORTIONTYPE, .TREATYTYPE, .TREATYGROUP, .TREATYGROUPID, .CLASSOFBUSINESSID, .CLASSOFBUSINESS, .LIMITCURRENCY, .LIMITVALUE, .RETENTIONCURRENCY, .RETENTIONVALUE, .EPICURRENCY, .EPIVALUE, .MDPCURRENCY, .MDPVALUE, .NETPREMICURRENCY, .NETPREMIVALUE, .SOBID, .SOB,  … | A: .SOBID = Param.SOB; B: .CLASSOFBUSINESSID = Param.CLASSOFBUSINESSID; C: .TREATYYEAR = Param.TREATYYEAR; D: .CEDING = Param.CEDING; K: .CEDINGID = Param.CEDINGID; L: .CEDINGID = Param.CEDINGID1; E: .TREATYGROUP = Param.TREATYGROUP; F: .TREATYTYPE = Param.TREATYTYPE; G: .ID = Param.ID; I: .TREATYID … | BusinessAndSOBList, BusinessAndSOBListRetro |
| `BrowseTreatyJoinEDM` | ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM | .ID, .TREATYID, .TREATYCONTRACTNAME, .PROPORTIONTYPE, .TREATYTYPE, .TREATYGROUP, .TREATYGROUPID, .CLASSOFBUSINESSID, .CLASSOFBUSINESS, .LIMITCURRENCY, .LIMITVALUE, .RETENTIONCURRENCY, .RETENTIONVALUE, .EPICURRENCY, .EPIVALUE, .MDPCURRENCY, .MDPVALUE, .NETPREMICURRENCY, .NETPREMIVALUE, .SOBID, .SOB,  … | A: .SOBID = Param.SOB; B: .CLASSOFBUSINESSID = Param.CLASSOFBUSINESSID; C: .TREATYYEAR = Param.TREATYYEAR; D: .CEDING = Param.CEDING; K: .CEDINGID = Param.CEDINGID; L: .CEDINGID = Param.CEDINGID1; E: .TREATYGROUP = Param.TREATYGROUP; F: .TREATYTYPE = Param.TREATYTYPE; G: .ID = Param.ID; I: .TREATYID … | BusinessAndSOBList, InputPolicyTreatyInDetail_preACT |
| `BrowseTreatyOutDetail` | ASM-FW-GISFW-Int-TREATYOUTDETAIL | .ID, .TREATYID, .TREATYCONTRACTNAME, .PROPORTIONTYPE, .TREATYTYPE, .TREATYGROUP, .TREATYGROUPID, .CLASSOFBUSINESSID, .CLASSOFBUSINESS, .LIMITCURRENCY, .LIMITVALUE, .RETENTIONCURRENCY, .RETENTIONVALUE, .EPICURRENCY, .EPIVALUE, .MDPCURRENCY, .MDPVALUE, .NETPREMICURRENCY, .NETPREMIVALUE, .SOBID, .SOB,  … | A: .SOBID = Param.SOB; B: .CLASSOFBUSINESSID = Param.CLASSOFBUSINESSID; C: .TREATYYEAR = Param.TREATYYEAR; D: .CEDING = Param.CEDING; K: .CEDINGID = Param.CEDINGID; L: .CEDINGID = Param.CEDINGID1; E: .TREATYGROUP = Param.TREATYGROUP; F: .TREATYTYPE = Param.TREATYTYPE; G: .ID = Param.ID; I: .TREATYID … | BusinessAndSOBListRetro, InputPolicyTreatyOutDetail_preACT |
| `crmOpportunitiesList` | ASM-FW-SFAGISFW-Work-Opportunity | pyID, Name, Acc.Name, OpportunityStage, OpportunityAmount, CloseDate, OpportunityMustWin, pyOwnerUserID, TerritoryID, Org.Name, Acc.OrganizationID, pxCreateDateTime, pxCreateOpName, pxUpdateDateTime, pxUpdateOpName, Acc.TerritoryID, Org.TerritoryID, Acc.pzInsKey, Org.pzInsKey, Party.pxPartyRole, pzI … | A: .AccountID = Acc.pzInsKey; A: Acc.OrganizationID = Org.pzInsKey; A: Oper.pyUserIdentifier = .pyOwnerUserID; A: .SellingMode = Param.SellingMode; B: Oper.pyUserName = OperatorID.pyUserName; C: .TextNoQuotation Contains Param.NB | SFAPortal_OpportunitiesList |
| `GetListOpportunity` | ASM-FW-SFAGISFW-Work-Opportunity | pyID, Name, Acc.Name, OpportunityStage, OpportunityAmount, CloseDate, OpportunityMustWin, pyOwnerUserID, TerritoryID, Org.Name, Acc.OrganizationID, pxCreateDateTime, pxCreateOpName, pxUpdateDateTime, pxUpdateOpName, Acc.TerritoryID, Org.TerritoryID, Acc.pzInsKey, Org.pzInsKey, Party.pxPartyRole, pzI … | A: A.pzInsKey = .NBHandle; B: .TextNoQuotation Contains "NB-"; D: A.pyStatusWork != "Resolved-Completed"; E: A.Quotation.BusinessFac = T; A: A.pxCreateOperator = Param.UserIdentifier; F: A.pyStatusWork != "Resolved-Rejected"; C: .Name Contains Param.Search; G: .TextNoQuotation Contains Param.Search; … | SFAPortal_OpportunitiesList |

## 10 · Decision Table

### 10 · `BusinessType_DeT`

Kelas `ASM-FW-GISFW-Work` · bawaan `"UNKNOWN"` (RETURN) · evaluasi semua baris: `no` (baris pertama yang cocok menang) · kolom: `.Quotation.GroupPanel` (text, `=`), `.Quotation.BusinessOldId` (text, `=`)

| Baris | `.Quotation.GroupPanel` | `.Quotation.BusinessOldId` | Hasil |
| ---: | --- | --- | --- |
| 1 | `"001"` | *(apa saja)* | `"Medicare"` |
| 2 | `"002"` | `"03"` | `"PA"` |
| 3 | `"003"` | *(apa saja)* | `"Bonding"` |
| 4 | `"003"` | *(apa saja)* | `"BondingKBG"` |
| 5 | `"003"` | `"20"` | `"HE"` |
| 6 | `"003"` | *(apa saja)* | `"MarineHull"` |
| 7 | `"003"` | `"21"` | `"Glass"` |
| 8 | `"003"` | *(apa saja)* | `"Liability"` |
| 9 | `"003"` | `"16"` | `"AllRisk"` |
| 10 | `"003"` | `"06"` | `"AviationHull"` |
| 11 | `"003"` | `"14"` | `"Burglary"` |
| 12 | `"003"` | `"07"` | `"Car"` |
| 13 | `"003"` | `"08"` | `"Ear"` |
| 14 | `"003"` | `"09"` | `"ElectronicEquipment"` |
| 15 | `"003"` | `"19"` | `"Fidelity"` |
| 16 | `"003"` | `"27"` | `"GolfInsurance"` |
| 17 | `"003"` | `"17"` | `"CIT"` |
| 18 | `"003"` | `"18"` | `"CIS"` |
| 19 | `"003"` | `"11"` | `"Boiler"` |
| 20 | `"003"` | `"93"` | `"Workmen"` |
| 21 | `"003"` | `"10"` | `"MBD"` |
| 22 | `"003"` | `"12"` | `"ContractorsPM"` |
| 23 | `"003"` | *(apa saja)* | `"CustomBond"` |
| 24 | `"003"` | `"28"` | `"FireStyle1"` |
| 25 | `"003"` | `"B3"` | `"LandRig"` |
| 26 | `"003"` | *(apa saja)* | `"Aneka"` |
| 27 | `"004"` | *(apa saja)* | `"MarineCargo"` |
| 28 | `"005"` | *(apa saja)* | `"Travel"` |
| 29 | `"006"` | `"78"` | `"OilGas"` |
| 30 | `"006"` | *(apa saja)* | `"FireStyle1"` |
| 31 | `"006"` | *(apa saja)* | `"FireStyle2"` |
| 32 | `"006"` | *(apa saja)* | `"Fire"` |
| 33 | `"007"` | `"02"` | `"MBUCar"` |
| 34 | `"007"` | `"60"` | `"MBUMotorCycle"` |
| 35 | `"009"` | *(apa saja)* | `"Life"` |
| 36 | `"008"` | *(apa saja)* | `"Medicare"` |

*Sel bertanda kutip ditiru apa adanya — termasuk spasi di dalam kutip.* `pyRowNum` daftar-ATAU berbasis nol, dipetakan ke baris = `pyRowNum + 1`.

### 10 · `isApproved`

Kelas `ASM-FW-GISFW-Work` · bawaan `YES` (RETURN) · evaluasi semua baris: `no` (baris pertama yang cocok menang) · kolom: `pyWorkPage.PolicyTreatyIn.IsApproved` (text, `=`)

| Baris | `pyWorkPage.PolicyTreatyIn.IsApproved` | Hasil |
| ---: | --- | --- |
| 1 | `0` | `No` |

*Sel bertanda kutip ditiru apa adanya — termasuk spasi di dalam kutip.* `pyRowNum` daftar-ATAU berbasis nol, dipetakan ke baris = `pyRowNum + 1`.

## 11 · When — syarat yang DIJALANKAN

Yang dijalankan adalah `pyLogic` + `pyCondition` (P23). Kolom *penampil* hanya untuk perbandingan; bila berbeda, yang dijalankan benar.

| Nama | Kelas | Logika dijalankan | Kondisi | Penampil (tidak dijalankan) | Dipakai |
| --- | --- | --- | --- | --- | --- |
| `crmCreateOpportunity` | PegaCRM-Work- | `(A0 Or A1) And A2` | A0: true = @(Pega-RULES:ExpressionEvaluators).evaluateWhen("crmBypassOperatorAccessChecks"); A1: "true" equals Declare_crmOperatorAccess.canCreateOpportunity; A2: true = @(Pega-RULES:ExpressionEvaluators).evaluateWhen("crmIsOpen") | : [Double click to add condition] | SFAPortalOpportunitiesHeader |
| `isAllRisk` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "AllRisk" | A: pyWorkPage.Quotation.BusinessType = "AllRisk" | IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsAneka` ⛔ | ASM-FW-GISFW-Data | `A OR B OR C OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U OR V OR W OR X OR Y OR Z OR AA` | A: Rule IsLiability evaluates to true; B: Rule IsMarineHull evaluates to true; C: Rule isGrowingTrees evaluates to true; E: Rule IsExclusion evaluates to true; F: Rule isCAR evaluates to true; G: Rule IsMaintenance evaluates to true; H: Rule isElectronicEquipment evaluates to true; I: Rule isAviationHull evaluates to true; J: Rule isEAR evaluates to true; K: Rule isAllRisk evaluates to true; L: Rule IsGlass evaluates to true; M: Rule isFidelity evaluates to true; N: Rule isBillboardNeonSyariah evaluates to true; O: Rule isBurglary evaluates to true; P: Rule IsCIT evaluates to true; Q: Rule IsC … | A OR B OR C OR D OR E OR F: [Double click to add condition] | CheckSpreadingProtectAnekaGolf_ACT, ProtectCoverage_Act, SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_Act |
| `isApproved` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.PolicyTreatyIn.IsApproved = 1 | A: pyWorkPage.PolicyTreatyIn.IsApproved = 1 |  |
| `isAviationHull` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "AviationHull" | A: pyWorkPage.Quotation.BusinessType = "AviationHull" | IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsBillboardNeon` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.Quotation.BusinessName = "BILLBOARD/NEON SIGN" | A: pyWorkPage.Quotation.BusinessName = "BILLBOARD/NEON SIGN" | IsLiability |
| `isBillboardNeonSyariah` ⛔ | ASM-FW-GISFW-Data | `A OR B` | A: pyWorkPage.Quotation.BusinessType = "BillboardNeonSyariah"; B: pyWorkPage.Quotation.BusinessCode = "10106" | A: pyWorkPage.Quotation.BusinessType = "BillboardNeonSyariah" | IsAneka |
| `IsBoiler` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "Boiler" | A: pyWorkPage.Quotation.BusinessType = "Boiler" | IsAneka |
| `IsBonding` ⛔ | ASM-FW-GISFW-Work | `A OR B OR C` | A: .Quotation.BusinessType = "Bonding"; B: Rule IsBondingKBG evaluates to true; C: pyWorkPage.Quotation.BusinessType = "Bonding" | A OR B OR C OR D: [Double click to add condition] | CheckSpreadingProtect_ACT |
| `IsBondingAndCustomBonds` ⛔ | ASM-FW-GISFW-Data | `A OR B OR C` | A: pyWorkPage.Quotation.BusinessType = "Bonding"; B: Rule IsBondingKBG evaluates to true; C: Rule IsCustomBonds evaluates to true | A OR B OR C OR D: Kode Bisnis = "02"; Kode Bisnis = "58"; Kode Bisnis = "SB"; Kode Bisnis = "SG" | CheckSpreadingProtectAnekaGolf_ACT, IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_Act |
| `IsBondingKBG` ⛔ | ASM-FW-GISFW-Work | `A OR B` | A: .Quotation.BusinessType = "BondingKBG"; B: pyWorkPage.Quotation.BusinessType = "BondingKBG" | A: BusinessType = "BondingKBG" | IsBonding, IsBondingAndCustomBonds |
| `IsBuilderRisk` ⛔ | ASM-FW-GISFW-Data | `A OR B` | A: pyWorkPage.OfferFacIn.QuotationData.BusinessCode = "10085"; B: pyWorkPage.OfferFacIn.QuotationData.BusinessName = "BUILDER RISK" | A: [Double click to add condition] | SumTSIPremiSpreadedRNM_ANEKA_Act |
| `isBurglary` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "Burglary" | A: pyWorkPage.Quotation.BusinessType = "Burglary" | IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsCar` ⛔ | ASM-FW-GISFW-Work | `A` | A: .OfferFacIn.QuotationData.BusinessType = "Car" | A: BusinessType = "Car" | IsAneka |
| `IsCIS` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "CIS" | A: pyWorkPage.Quotation.BusinessType = "CIS" | IsAneka |
| `IsCIT` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "CIT" | A: pyWorkPage.Quotation.BusinessType = "CIT" | IsAneka |
| `IsClaim` | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.pyWorkIDPrefix = "CLM-" | A: pyWorkPage.pyWorkIDPrefix = "CLM-" | DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `isClaimTreaty` | ASM-FW-GISFW-Work | `G AND (A OR B OR C OR D OR E OR F)` | A: .PolicyTreatyIn.ClaimType = "XOL"; B: .PolicyTreatyIn.Claim != 0; C: .PolicyTreatyIn.SalvageValue != 0; D: .PolicyTreatyIn.ExcessLoss != 0; E: .PolicyTreatyIn.ClaimPaymentType = "Salvage"; F: .PolicyTreatyIn.ClaimPaymentType = "Claim"; G: .Quotation.ProportionalType != "Proportional" | A: [Double click to add condition] | InputRealizationTreatyIn |
| `IsCMI` ⛔ | ASM-FW-GISFW-Work | `A` | A: .OfferFacIn.QuotationData.BusinessName = "COMPREHENSIVE MACHINERIES INSURANCE" | A: Nama Bisnis = "COMPREHENSIVE MACHINERIES INSURANCE" | SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsContractorsPlantMachinery` ⛔ | ASM-FW-GISFW-Work | `A` | A: .OfferFacIn.QuotationData.BusinessType = "ContractorsPM" | A: BusinessType = "ContractorsPM" | IsAneka |
| `IsCrime` ⛔ | ASM-FW-GISFW-Data | `A AND B` | A: pyWorkPage.Quotation.BusinessCode = 10168; B: pyWorkPage.Quotation.BusinessType = "Aneka" | A: [Double click to add condition] | IsAneka |
| `IsCustomBonds` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "CustomBond" | A: pyWorkPage.Quotation.BusinessType = "CustomBond" | IsBondingAndCustomBonds |
| `IsEar` ⛔ | @baseclass | `A` | A: pyWorkPage.Quotation.BusinessType = "Ear" | A: pyWorkPage.Quotation.BusinessType = "Ear" | IsAneka |
| `IsEDM` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3 | A: pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3 | CheckSpreadingProtect_ACT, GetTreatyName, ProtectCoverage_Act, ProtectFIREMBUPA_Act, ProtectPremiPolicy_Act, ProtectShipData_Act, Protection_Act, cekSpreadingFactIn |
| `IsEdmAdjShareCedant` ⛔ | ASM-FW-GISFW-Work | `A AND B AND C` | A: pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3; B: pyWorkPage.OfferFacIn.QuotationData.EdmType = 4; C: pyWorkPage.OfferFacIn.QuotationData.Type = 11 | A: [Double click to add condition] | ProtectShareCedant_Act |
| `IsEdmAdjSpreading` ⛔ | ASM-FW-GISFW-Work | `A AND B AND C` | A: pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3; B: pyWorkPage.OfferFacIn.QuotationData.EdmType = 4; C: pyWorkPage.OfferFacIn.QuotationData.Type = 4 | A: [Double click to add condition] | SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_FIRE_Act, SumTSIPremiSpreadedRNM_GOLF_Act, SumTSIPremiSpreadedRNM_MBU_Act |
| `IsEDMRiSlip` ⛔ | ASM-FW-GISFW-Work | `A OR B` | A: pyWorkPage.OfferFacIn.QuotationData.EdmTypeNew = 4; B: pyWorkPage.OfferFacIn.QuotationData.Type = 0 | A: pyWorkPage.OfferFacIn.QuotationData.EdmTypeNew = 4 | Protection_Act |
| `isElectronicEquipment` ⛔ | ASM-FW-GISFW-Work | `A` | A: .OfferFacIn.QuotationData.BusinessType = " ElectronicEquipment " | A: BusinessType = " ElectronicEquipment " | IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsEngineering` ⛔ | ASM-FW-GISFW-Work | `A OR B OR C OR D OR E` | A: pyWorkPage.Quotation.BusinessType = "Car"; B: pyWorkPage.Quotation.BusinessType = "Ear"; C: pyWorkPage.Quotation.BusinessType = "MBD"; D: pyWorkPage.Quotation.BusinessCode = "10166"; E: pyWorkPage.Quotation.BusinessCode = "10032" | A: [Double click to add condition] | SumTSIPremiSpreadedRNM_Act |
| `IsEnvironmental` ⛔ | ASM-FW-GISFW-Data | `A AND B` | A: pyWorkPage.Quotation.BusinessCode = 10169; B: pyWorkPage.Quotation.BusinessType = "Aneka" | A: [Double click to add condition] | IsAneka |
| `IsErrorSpreading` ⛔ | @baseclass | `A OR B` | A: pyWorkPage.Quotation.OldPolicyNo = "<nomor-polis-1>"; B: pyWorkPage.Quotation.OldPolicyNo = "<nomor-polis-2>" | A: pyWorkPage.Quotation.OldPolicyNo = "<nomor-polis-1>" | GetTreatyName |
| `IsExclusion` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "Exclusion" | A: pyWorkPage.Quotation.BusinessType = "Exclusion" | IsAneka |
| `isFidelity` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "Fidelity" | A: pyWorkPage.Quotation.BusinessType = "Fidelity" | IsAneka |
| `IsFire` ⛔ | ASM-FW-GISFW-Work | `A OR B OR C OR D OR E` | A: Rule IsKPR evaluates to true; B: Rule IsOilGas evaluates to true; C: Rule IsFireStyle1 evaluates to true; D: Rule IsFireStyle2 evaluates to true; E: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Fire" | A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U: [Double click to add condition] | CheckSpreadingProtectFire_ACT, GetTreatyName, ProtectCurrencyTSI_Act, ProtectFIREMBUPA_Act, Protection_Act, SetFlagOccupation_ACT, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_FIRE_Act |
| `IsFireStyle1` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "FireStyle1" | A: [Double click to add condition] | IsFire |
| `IsFireStyle2` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "FireStyle2" | A: [Double click to add condition] | IsFire |
| `IsGlass` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "Glass" | A: pyWorkPage.Quotation.BusinessType = "Glass" | IsAneka |
| `IsGolfInsurance` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "GolfInsurance" | A: pyWorkPage.Quotation.BusinessType = "GolfInsurance" | CheckSpreadingProtectAnekaGolf_ACT, CheckSpreadingProtect_ACT, ProtectCoverage_Act, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_GOLF_Act |
| `IsGrowingTrees` ⛔ | @baseclass | `A` | A: pyWorkPage.Quotation.BusinessType = "GrowingTrees" | A: pyWorkPage.Quotation.BusinessType = "GrowingTrees" | IsAneka |
| `IsHE` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "HE" | A OR B OR C OR D OR E OR F: Kode Bisnis = "24"; Kode Bisnis = "18"; Kode Bisnis = "17"; Kode Bisnis = "10"; Kode Bisnis = "07"; Kode Bisnis = "06" | IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsKPR` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "KPR" | A OR B: [Double click to add condition] | IsFire |
| `IsLandRig` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "LandRig" | A: pyWorkPage.Quotation.BusinessType = "LandRig" | IsAneka |
| `IsLiability` ⛔ | ASM-FW-GISFW-Work | `A OR B OR C OR D OR E OR F OR G OR H` | A: .OfferFacIn.QuotationData.BusinessType = "Liability"; B: Rule IsProductsLiability evaluates to true; C: Rule IsProfessionalLiability evaluates to true; D: Rule IsWorkmenCompensation evaluates to true; E: pyWorkPage.Quotation.BusinessCode = "10048"; F: pyWorkPage.Quotation.BusinessCode = "10184"; G: Rule IsBillboardNeon evaluates to true; H: pyWorkPage.Quotation.BusinessCode = "10253" | A: BusinessType = "Liability" | IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsLife` ⛔ | ASM-FW-GISFW-Data | `A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P` | A: pyWorkPage.Quotation.BusinessOldId = "L1"; B: pyWorkPage.Quotation.BusinessOldId = "L2"; C: pyWorkPage.Quotation.BusinessOldId = "L3"; D: pyWorkPage.Quotation.BusinessOldId = "L4"; E: pyWorkPage.Quotation.BusinessOldId = "L5"; F: pyWorkPage.Quotation.BusinessOldId = "L6"; G: pyWorkPage.Quotation.BusinessOldId = "L7"; H: pyWorkPage.Quotation.BusinessOldId = "L8"; I: pyWorkPage.Quotation.BusinessOldId = "L9"; J: pyWorkPage.Quotation.BusinessOldId = "L10"; K: pyWorkPage.Quotation.BusinessOldId = "L11"; L: pyWorkPage.Quotation.BusinessOldId = "L12"; M: pyWorkPage.Quotation.BusinessOldId = "L13" … | A: pyWorkPage.OfferFacIn.QuotationData.BusinessCode = "10164" | ProtectPremiPolicy_Act, ProtectShareCedant_Act |
| `IsMaintenance` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "Maintenance" | A: pyWorkPage.Quotation.BusinessType = "Maintenance" | IsAneka |
| `IsMarineCargo` ⛔ | ASM-FW-GISFW-Work | `G OR H` | G: .OfferFacIn.QuotationData.BusinessType = "MarineCargo"; H: pyWorkPage.Quotation.BusinessType = "MarineCargo" | A OR B OR C OR D OR E OR F: Kode Bisnis = "24"; Kode Bisnis = "18"; Kode Bisnis = "17"; Kode Bisnis = "10"; Kode Bisnis = "07"; Kode Bisnis = "06" | CheckSpreadingProtectMCargoMBU_ACT, GetTreatyName, ProtectCoverage_Act, Protection_Act, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_MARINECARGO_Act |
| `IsMarineHull` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "MarineHull" | A: pyWorkPage.Quotation.BusinessType = "MarineHull" | IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsMarineHullOffshore` ⛔ | ASM-FW-GISFW-Data | `A AND (B OR C)` | A: pyWorkPage.Quotation.BusinessType = "Aneka"; B: pyWorkPage.Quotation.BusinessCode = "10157"; C: pyWorkPage.Quotation.BusinessName = "MARINE HULL OFFSHORE" | A OR B OR C OR D OR E OR F: Kode Bisnis = "24"; Kode Bisnis = "18"; Kode Bisnis = "17"; Kode Bisnis = "10"; Kode Bisnis = "07"; Kode Bisnis = "06" | SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsMBD` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "MBD" | A: pyWorkPage.Quotation.BusinessType = "MBD" | IsAneka, SumTSIPremiSpreadedRNM_ANEKA_Act, SumTSIPremiSpreadedRNM_Act |
| `IsMBU` ⛔ | ASM-FW-GISFW-Data | `A OR B` | A: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "MBUCar"; B: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "MBUMotorCycle" | A OR B OR C OR D: pyWorkPage.Quotation.BusinessCode = "10138"; pyWorkPage.Quotation.BusinessCode = "10009"; pyWorkPage.Quotation.BusinessCode = "10092"; pyWorkPage.Quotation.BusinessCode = "10096" | CheckSpreadingProtectMCargoMBU_ACT, ProtectFIREMBUPA_Act, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_MBU_Act |
| `IsMBUCar` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.Quotation.BusinessType = "MBUCar" | A: pyWorkPage.Quotation.BusinessType = "MBUCar" | CheckSpreadingProtect_ACT |
| `IsNotAdmin` | @baseclass | `A` | A: OperatorID.pyPosition != "Admin" | A: OperatorID.pyPosition != "Admin" | SFAPortalOpportunitiesHeader |
| `IsObjectWithQuantityYear` ⛔ | ASM-FW-GISFW-Data | `A OR B OR C OR D` | A: pyWorkPage.Quotation.BusinessType = "MBD"; B: pyWorkPage.Quotation.BusinessType = "Boiler"; C: pyWorkPage.Quotation.BusinessType = "ContractorsPM"; D: pyWorkPage.Quotation.BusinessType = "LandRig" | A OR B: pyWorkPage.Quotation.BusinessCode = "10"; pyWorkPage.Quotation.BusinessCode = "11" | SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsOilGas` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "OilGas" | A: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "OilGas" | IsFire |
| `IsOperatorLife` | Data-Portal | `A` | A: OperatorID.pyWorkGroup = "ReasLife" | A: OperatorID.pyAccessGroup = "GISFW:AuctionUsers" | SFAPortal_OpportunitiesList |
| `IsPA` ⛔ | Data-Party-Person | `A` | A: pyWorkPage.Quotation.BusinessType = "PA" | A: pyWorkPage.Quotation.BusinessType = "PA" | CheckSpreadingProtectPATravel_ACT, ProtectFIREMBUPA_Act, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_PA_Act |
| `IsPortRisk` ⛔ | ASM-FW-GISFW-Data | `A` | A: pyWorkPage.OfferFacIn.QuotationData.BusinessCode = "10188" | A: [Double click to add condition] | SumTSIPremiSpreadedRNM_ANEKA_Act |
| `IsProductsLiability` ⛔ | ASM-FW-GISFW-Work | `A AND B` | A: pyWorkPage.Quotation.BusinessType = "Aneka"; B: pyWorkPage.Quotation.BusinessCode = "10021" | A OR B OR C OR D OR E OR F: Kode Bisnis = "24"; Kode Bisnis = "18"; Kode Bisnis = "17"; Kode Bisnis = "10"; Kode Bisnis = "07"; Kode Bisnis = "06" | IsLiability |
| `IsProfessionalLiability` ⛔ | ASM-FW-GISFW-Work | `A AND B` | A: pyWorkPage.Quotation.BusinessType = "Aneka"; B: pyWorkPage.Quotation.BusinessCode = "10023" | A OR B OR C OR D OR E OR F: Kode Bisnis = "24"; Kode Bisnis = "18"; Kode Bisnis = "17"; Kode Bisnis = "10"; Kode Bisnis = "07"; Kode Bisnis = "06" | IsLiability |
| `IsRenewal` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.Quotation.StatusBusiness = 2 | A: pyWorkPage.Quotation.StatusBusiness = 2 | ProtectCoverage_Act, ProtectFIREMBUPA_Act, ProtectRenewal_Act |
| `isSellingModeB2B` | @baseclass | `A` | A: @(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") EQUALS "B2B" | A: @(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") EQUALS "B2B" | SFAPortalOpportunitiesHeader |
| `isSellingModeB2BB2C` | @baseclass | `A` | A: @(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") equals "B2B_B2C" | A: @(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") equals "B2B_B2C" | SFAPortalOpportunitiesHeader |
| `isSellingModeB2C` | @baseclass | `A` | A: @(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") equals "B2C" | A: @(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") equals "B2C" | SFAPortalOpportunitiesHeader |
| `IsSPVCreate` | ASM-FW-GISFW-Work | `A OR B` | A: pyWorkPage.pxCreateOperator = "<ID-operator-5>"; B: pyWorkPage.pxCreateOperator = "<ID-operator-6>" | A: [Double click to add condition] | InputRealizationTreatyIn |
| `IsSPVTreaty1` | ASM-FW-GISFW-Work | `A` | A: OperatorID.pyTelephone = "SPVTREATY1" | A: [Double click to add condition] | InputRealizationTreatyIn |
| `IsTBonding` ⛔ | ASM-FW-GISFW-Work | `A OR B OR C OR D` | A: .OfferFacIn.QuotationData.TeamGroup = 5; B: .OfferFacIn.QuotationData.MarketingName = "<nama-6>"; C: .OfferFacIn.QuotationData.MarketingName = "<nama-7>"; D: .OfferFacIn.QuotationData.MarketingName = "<nama-8>" | A: [Double click to add condition] | ProtectCoverage_Act |
| `IsTravel` ⛔ | ASM-FW-GISFW-Work | `A OR B` | A: .Quotation.BusinessType = "Travel"; B: pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Travel" | A: Kode Bisnis = "77" | CheckSpreadingProtectPATravel_ACT, SumTSIPremiSpreadedRNM_Act, SumTSIPremiSpreadedRNM_TRAVEL_Act |
| `IsTreaty1` | ASM-FW-GISFW-Work | `A` | A: OperatorID.pyTelephone = "TREATY1" | A: [Double click to add condition] | InputRealizationTreatyIn |
| `IsUW` | ASM-FW-GISFW-Work | `A OR B OR C` | A: pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInDirector; B: pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInGroupLeader; C: pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInUnderwriting | A: pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInGroupLeader | DetailPolicyTreatyIn |
| `IsWorkmenCompensation` ⛔ | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.Quotation.BusinessType = "Workmen" | A: pyWorkPage.Quotation.BusinessType = "Workmen" | IsLiability |
| `IsYieldShortfall` ⛔ | ASM-FW-GISFW-Data | `A AND B` | A: pyWorkPage.Quotation.BusinessType = "Aneka"; B: pyWorkPage.Quotation.BusinessCode = 10167 | A: pyWorkPage.Quotation.BusinessType = "GrowingTrees" | IsAneka |
| `NopolisEmpty` | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.PolicyTreatyIn.PolicyNo = "" | A: pyWorkPage.PolicyTreatyIn.PolicyNo = "" | InputRealizationTreatyIn |
| `pyIsIpadOrDesktop` | @baseclass | `A OR B` | A: Rule pyIsIPad evaluates to true; B: pxRequestor.pxDeviceType = desktop | : [Double click to add condition] | SFAPortalOpportunitiesHeader |
| `ToTREATYDEPTHEAD` | ASM-FW-GISFW-Work | `A` | A: pyWorkPage.LetterNo = "TREATYINDEPTHEAD" | A: pyWorkPage.LetterNo = "TREATYINDEPTHEAD" | InputRealizationTreatyIn |
| `TreatyMasterInEDM` | ASM-FW-GISFW-Data-PolicyTreatyIn | `A OR B OR C` | A: pyWorkPage.TreatyIn.EDMState = "1"; B: pyWorkPage.TreatyIn.EDMState = "2"; C: pyWorkPage.TreatyIn.EDMState = "3" | A: [Double click to add condition] | DetailPoliciesNonProportional, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp |

## 12 · Rujukan ke rule yang tidak ada di korpus, dan selisih kelas

Hanya dari rule yang terjangkau. Rule bawaan Pega (`py*`, `pz*`, `px*`, `Always`, `Never`) tidak dicantumkan.

| Dari | Jenis | Nama dirujuk | Lewat | Kelas konteks |
| --- | --- | --- | --- | --- |
| `GeneratePolicyNoTreaty_Act` | RDBList | `GenerateNoPolicyOJKID` | langkah 22 RDB-List |  |
| `GeneratePolicyNoTreaty_Act` | RDBList | `GenerateNoPolicyTreaty` | langkah 20 RDB-List |  |
| `GeneratePolicyNoTreaty_Act` | RDBList | `GenerateNoPolicyTreatyGroup` | langkah 21 RDB-List |  |
| `InputPolicyTreatyInDetail_preACT` | DecisionTable | `false` | langkah 14.8 Property-Map-DecisionTable |  |
| `InputPolicyTreatyOutDetail_preACT` | DecisionTable | `false` | langkah 15.8 Property-Map-DecisionTable |  |
| `SaveJsonPolisTreatyIn_Act` | Activity | `CopyToPolicy` | langkah 1 Call | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `InputRealizationTreatyIn` | Activity | `ToWorkbasket` | pxRuleReferences | ASM-FW-GISFW-Work |
| `InputRealizationTreatyIn` | Activity | `WorkBasket` | pxRuleReferences | ASM-FW-GISFW-Work |
| `InputRealizationTreatyIn` | Harness | `Perform` | pxRuleReferences | ASM-FW-GISFW-Work |
| `SFAPortalOpportunities` | When | `recordEvent` | pxRuleReferences | PegaCRM-Portal |
| `ListSuggest` | FlowAction | `TreatyInSuggest_FlowAct` | pyEditAction |  |
| `ListSuggest` | Section | `DetailPolicyTreatyInAddendum` | aksi refresh medan .IsApproved |  |
| `SFAPortalOpportunitiesHeader` | Activity | `SFAPopulateListHeader` | aksi refresh medan .pyTemplateInputBox | Data-Portal |
| `SFAPortalOpportunitiesHeader` | DataTransform | `SetStageViewParams` | pyPreDataTransform |  |
| `crmCreateOpportunity` | When | `crmBypassOperatorAccessChecks` | pxRuleReferences | PegaCRM-Work- |
| `crmCreateOpportunity` | When | `crmIsOpen` | pxRuleReferences | PegaCRM-Work- |

**Selisih kelas** (nama ada di korpus, kelas berkasnya bukan leluhur kelas konteks rujukan):

| Dari | Ke | Lewat | Kelas konteks | Kelas berkas |
| --- | --- | --- | --- | --- |
| `InputPolicyTreatyInDetail_NonProp` | `When/TreatyMasterInEDM` | pxRuleReferences | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `InputPolicyTreatyOutDetail_NonProp` | `When/TreatyMasterInEDM` | pxRuleReferences | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `TreatyInputPctCommSpreading` | `Activity/CountSpreading_Act` | langkah 4 Call | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `TreatyNonPropOutSetSpreading` | `Activity/CountSpreading_Act` | langkah 7.2 Call | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `TreatyNonPropSetSpreading` | `Activity/CountSpreading_Act` | langkah 10.2 Call | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `serviceInsertArasapas_act` | `Activity/serviceInsertArasapas_act` | langkah 1 Call | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `serviceInsertArasapas_act` | `Activity/serviceInsertArasapas_act` | pxRuleReferences | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `InputRealizationTreatyIn` | `Activity/SaveJsonPolisTreatyIn_Act` | pxRuleReferences | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `InputRealizationTreatyIn` | `Activity/SaveJsonPolisTreatyIn_Act` | shape Utility1 | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `InputRealizationTreatyIn` | `Activity/serviceInsertArasapas_act` | pxRuleReferences | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `InputRealizationTreatyIn` | `Activity/serviceInsertArasapas_act` | shape Utility2 | ASM-FW-GISFW-Work | ASM-FW-GISFW-Data-PolicyTreatyIn |
| `SFAPortalOpportunities` | `Section/SFAPortal_Opportunities` | pxRuleReferences | PegaCRM-Portal | Data-Portal |
| `BusinessAndSOBList` | `Activity/SetValue_Act` | aksi runActivity medan .pyTemplateInputBox | ASM-FW-GISFW-Int-TREATYINDETAIL | ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM |
| `BusinessAndSOBListRetro` | `Activity/SetValue_Act` | aksi runActivity medan .pyTemplateInputBox | ASM-FW-GISFW-Int-TREATYINDETAIL | ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM |
| `DetailPolicyTreatyIn` | `When/IsUW` | pxRuleReferences | ASM-FW-GISFW-Data-Installment | ASM-FW-GISFW-Work |
| `DetailPolicyTreatyIn` | `When/IsUW` | pxRuleReferences | ASM-FW-GISFW-Data-PolicyTreatyIn | ASM-FW-GISFW-Work |
| `ListSuggest` | `Activity/Protection_Act` | aksi refresh medan .IsApproved | ASM-FW-GISFW-Data-PolicyTreatyIn | ASM-FW-GISFW-Work |
| `ListSuggest` | `Activity/Protection_Act` | aksi refresh medan .Suggest | ASM-FW-GISFW-Data-PolicyTreatyIn | ASM-FW-GISFW-Work |
| `SFAPortalOpportunitiesHeader` | `When/crmCreateOpportunity` | pxRuleReferences | PegaCRM-Portal | PegaCRM-Work- |

## 13 · Padanan tombol → activity

Padanan setiap tombol / aksi klik-ubah layar NB Treaty In dengan rule XML yang dipanggilnya (prompt putaran 2 bab 4 akhir). Diisi tangan dari action set sel Section (pyActionSets: event -> aksi berurutan) yang dibaca ulang 2026-10-03; dibaca inventaris.py bab 13. Kolom: tombol (label / medan sel), section, tampil (pyVisible sel + wadah pyContainerVisibleWhen; nonaktif = pyDisabled), aksi_xml (urutan aksi action set), rule (rule yang dipanggil), padanan (kode / rute), status.

| Tombol / medan | Section | Tampil | Aksi XML (berurutan) | Rule dipanggil | Padanan kode / rute | Status |
| --- | --- | --- | --- | --- | --- | --- |
| Create opportunity | `SFAPortalOpportunitiesHeader` | crmCreateOpportunity && isSellingModeB2BB2C && IsNotAdmin (+ dua varian B2B / B2C) | click -> createWork | `Flow/InputRealizationTreatyIn` | pages/PortalNBTreatyIn.tsx tombol Create -> POST /api/nb-treaty-in/kasus (services.BuatKasus) | dibangun |
| Stage view / List view | `SFAPortalOpportunitiesHeader` | selalu (nonaktif menurut .CurrentMode) | click -> refresh SFAPopulateListHeader (pra-DT SetStageViewParams - tidak ada di korpus) -> refresh SFAPortal_Opportunities; click -> setValue .CurrentMode="List" -> refresh | `Activity/SFAPopulateListHeader` ⚠️ | - | tidak dibangun - milik portal CRM, bukan alur realisasi (status.json Section/SFAPortalOpportunitiesHeader) |
| Filter / .FilterTermForOpportunity | `SFAPortal_OpportunitiesList` | selalu | Filter click -> refresh; isian enter -> refresh; isian esc -> setValue "" -> refresh; ikon -> setValue .FilterTermForOpportunity="" -> postValue | `ReportDefinition/GetListOpportunity` | pages/PortalNBTreatyIn.tsx form saring (submit/enter = Filter; Esc = kosongkan + muat ulang; ikon pengosong C[1.2] pxIcon pyiconclearfield = tombol aria-label "Clear field": kosongkan isian TANPA muat ulang) -> GET /api/nb-treaty-in/kasus?cari= -> Param.Search = filter G RD `.TextNoQuotation Contains` (pengenal kasus, tanpa beda huruf; models.CocokCariPortal), pyMaxRecords 500 (models.BatasDaftarPortal) | RALAT R6 audit silang P3 (04-10-2026): dibangun - ikon pengosong kini tombol tersendiri (setValue "" -> postValue, tanpa refresh; bunyi lama: "ikon pengosong tanpa refresh tidak dibuat tombol tersendiri"); filter C `.Name` (kelas CRM, ditulis nol rule) tidak dibangun - P9 |
| .Name (tautan berkas) | `SFAPortal_OpportunitiesList` | selalu; nonaktif bila bukan marketing/pembuat/TeamLeader (identitas - tiket 04/05) | click -> openWorkByHandle | `FlowAction/InboxPolicyTreatyIn`, `FlowAction/DeptHeadTreatyIn_UW` | pages/PortalNBTreatyIn.tsx TabelPortal: tautan di sel "Offer No" (`.TextNoQuotation`) -> GET /api/nb-treaty-in/kasus/{id} (services.BukaKasus: hanya-baca bila bukan anggota antrean); kolom "Name" tidak dirender | RALAT R6 audit silang P3 (04-10-2026): dibangun dengan penyimpangan sadar - `.Name` ditulis NOL rule korpus dan tak berkolom (tautan selalu tanpa teks), aksi openWorkByHandle (.pzInsKey, kunci berkas yang sama) dipindah ke sel Offer No (tiket 11 RALAT W6; bunyi lama status: "dibangun") |
| Select Source Of Business | `DetailPolicyTreatyIn` | .ClaimType = 'XOL Retro' | click -> showHarness (InputQuotation_PreAct Acton=SOB) -> Harness SOB | `Activity/InputQuotation_PreAct`, `DataTransform/btnSOB_DT`, `Harness/SOB`, `Section/SourceHierarki`, `Activity/AgentSourceBizTreatyIn_Act` | pages/LayarKasus.tsx (admin, tampilTombolSOB: ClaimType XOL Retro) -> components/PilihSumberBisnis.tsx -> GET /api/nb-treaty-in/sumber-bisnis (services.DaftarSumberBisnis); XML tanpa refresh sesudah pilih (hanya showHarness) - layar pun tidak (F4) | dibangun paket P3 (RALAT tiket 11) - CATATAN-PUTARAN-2 §2c |
| Choose (pohon SOB) | `SourceHierarki` | .ChildCount = 0 | click -> runActivity SearchHierarkiSourceBizAgentTreatyIn_Act -> runScript opener.location.reload -> runScript window.close | `Activity/SearchHierarkiSourceBizAgentTreatyIn_Act`, `FlowAction/AgentSourceBizDetails`, `DataTransform/SearchHierarkiSourceBizAgent_PostDT` | components/PilihSumberBisnis.tsx klik baris / Choose -> POST /api/nb-treaty-in/kasus/{id}/pilih-sumber-bisnis (services.PilihSumberBisnis: SearchHierarkiSourceBizAgent_PostDT TANPA simpan) -> hasil dipegang state LayarKasus (pegangSumberBisnis) -> ikut Save/Submit/refresh (services.terimaSumberBisnis: cocok dengan BrowseAgentHierarkiList_RD yang dijalankan ulang, selain itu 422) - F4 | dibangun paket P3 - aktivitas tombol Choose (menulis OfferTreatyIn.QuotationData.*, dibaca nol rule NB) tidak (RALAT tiket 11); F4 putaran 3 (R2): klik tidak menyimpan, disimpan bersama Save/Submit |
| Choose Business | `DetailPolicyTreatyIn` | wadah .ClaimType != 'XOL Retro'; InputParam.CARI12 != 'treaty' && != 'TreatyPolicy' (InputParam kosong di NB) | click -> showHarness BusinessAndSOBList (popup, pyWindowName 'Business And SOB List', pySubmitData Yes) -> refresh | `Harness/BusinessAndSOBList`, `Section/BusinessAndSOBList`, `ReportDefinition/BrowseTreatyJoinEDM` | pages/LayarKasus.tsx (admin, bukanXOLRetro) -> components/PilihBisnis.tsx -> POST /api/nb-treaty-in/kasus/{id}/bisnis dengan isian layar (services.DaftarBisnis: kerjakan tanpa simpan, 409 bila bukan admin / ClaimType 'XOL Retro'; repository.DaftarBisnis: view TREATYINDETAILJOINEDM, filter H PROPORTIONTYPE = PolicyTreatyIn.QuotationData.ProportionalType - kosong diabaikan, ORDER BY TREATYID, <= 500) | dibangun - RALAT audit silang P3 W1: bunyi lama rule 'ReportDefinition/BrowseTreatyInDetail' dan padanan 'GET /api/nb-treaty-in/bisnis' (tabel TREATYINDETAIL tanpa saringan) keliru - grid itu wadah S11 pyContainerVisibleWhen 1=2 |
| Choose (popup bisnis) | `BusinessAndSOBList` | selalu; SATU grid tampil - S16/S17 pyGridProps/pyRDName BrowseTreatyJoinEDM (grid lama S11 RD BrowseTreatyInDetail wadah pyContainerVisibleWhen 1=2; LABEL 'For Treatyindetail join EDM' pyCondition 1=2). RALAT audit silang P3 W1: bunyi lama 'selalu (dua grid; ...)' | click -> runActivity SetValue_Act(ID=.ID) -> runScript opener.location.reload -> runScript window.close -> refresh | `Activity/SetValue_Act`, `Activity/InputPolicyTreatyInDetail_preACT` | components/PilihBisnis.tsx Choose -> POST /api/nb-treaty-in/kasus/{id}/pilih-bisnis (services.PilihBisnis; Obj-Save ditiru) -> layar dimuat ulang dari jawaban | dibangun |
| Choose Business R | `DetailPolicyTreatyIn` | wadah .ClaimType = 'XOL Retro'; InputParam.CARI12 != 'treaty' && != 'TreatyPolicy' | click -> showHarness BusinessAndSOBListRetro (popup) -> refresh | `Harness/BusinessAndSOBListRetro`, `Section/BusinessAndSOBListRetro`, `Activity/SetValueRetro_Act` | - | tidak dibangun - treaty KELUAR, [keputusan work owner] K8 butir 4: SetValueRetro_Act -> InputPolicyTreatyOutDetail_preACT membaca JSON M_TREATY_OUT (RDBList\BrowseTreatyOutDetail: `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`) |
| Survey Report | `DetailPolicyTreatyIn` | pyWorkPage.Quotation.ProportionalType != 'NonProportional'; nonaktif .QuotationData.IsSurveyReport=='No' \|\| '' | click -> showHarness HistoricalSurveyReport (popup) | `Harness/HistoricalSurveyReport`, `FlowAction/InputHistoricalSurveyReport`, `Activity/SetSurveyReport_Act` | - | tidak dibangun - K7 (keputusan WO 03-10-2026): tidak ada tabel di diagram grilling |
| Enable / Disable Input Type | `DetailPolicyTreatyIn` | .TreatyType='XOL' | click -> refresh thisSection dengan pra-DT TreatyEnableDisableInput | `DataTransform/TreatyEnableDisableInput` | pages/LayarKasus.tsx -> POST /kasus/{id}/hitung {aksi: TreatyEnableDisableInput} -> models.TreatyEnableDisableInput (tanpa simpan; hasil dipegang layar); Save/Submit: models.terimaEnableDisable menerima IsNewPolicyNonProp / QuotationData.ProportionalType hanya bila TreatyType XOL dan sama dengan hasil DT atas halaman server (bukan XOL diabaikan, lain 422) | RALAT R6 audit silang P3 (04-10-2026): dibangun - W3 (pola F4) |
| Save (bagian uang) / Save (kaki layar) | `DetailPolicyTreatyIn` | bagian uang: InputParam.CARI12 != 'treaty' && != 'TreatyPolicy' di wadah uang; kaki: selalu | click -> save | — | pages/LayarKasus.tsx satu tombol Save (admin) -> PUT /api/nb-treaty-in/kasus/{id} (services.SimpanDraf; medan wajib ditegakkan - AC 48); membawa pilihan Source Of Business yang dipegang layar (F4, services.terimaSumberBisnis) | dibangun - dua tombol XML beraksi sama, dibangun satu |
| Add / Delete (grid spreading) | `DetailPolicyTreatyIn` | InputParam.CARI12 != 'treaty' && != 'TreatyPolicy' di wadah uang | click -> addRow / deleteRow (.SpreadingRiskList) | — | pages/LayarKasus.tsx (admin) mengubah daftar di klien; diterima server lewat models.GabungMasukanLayar -> models.SpreadingDariLayar (admin: wadah S19) saat Save/Submit/hitung; kolom PremiumSpreaded/ClaimSpreaded tidak diterima - CountSpreading_Act 4.1 di server bila terpicu (models.TerapkanPemicu) | RALAT R6 audit silang P3 (04-10-2026): dibangun - W5 kolom hanya-baca dari server |
| Submit (IsApproved==1) | `DetailPolicyTreatyIn` | .IsApproved==1; nonaktif TempEmail.CARI28 != pyWorkPage.PositionNote | click -> runActivity SetDueTo_act -> finishAssignment (flow action InboxPolicyTreatyIn: post-DT InboxPolicyTreatyIn_postDT, InputPolicyTreatyInPost_Act) | `Activity/SetDueTo_act`, `FlowAction/InboxPolicyTreatyIn`, `DataTransform/InboxPolicyTreatyIn_postDT`, `Activity/InputPolicyTreatyInPost_Act`, `Activity/InsertHistoryAkseptasiPega`, `Activity/TreatyRealizationCheckDuplicate` | pages/LayarKasus.tsx tombol kirim (models.TombolUntuk = kirim; bolehKerja = anggota antrean) -> POST /kasus/{id}/kirim -> services.Kirim: SetDueTo -> validasi wajib/pesan -> PascaAdmin -> cekDuplikat -> riwayat -> connector; membawa pilihan Source Of Business yang dipegang layar (F4, services.terimaSumberBisnis) | dibangun |
| Submit (IsApproved==0) | `DetailPolicyTreatyIn` | .IsApproved==0; nonaktif TempEmail.CARI28 != pyWorkPage.PositionNote | click -> localAction PolicyTreatyInDeclineConfirm (modal; local action activity SetDueTo_act) | `FlowAction/PolicyTreatyInDeclineConfirm`, `Section/PolicyTreatyInDeclineConfirm`, `Activity/SetDueTo_act` | pages/LayarKasus.tsx (models.TombolUntuk = konfirmasi-tolak) -> Modal KONFIRMASI_TOLAK | dibangun |
| Yes / No (konfirmasi tolak) | `PolicyTreatyInDeclineConfirm` | selalu | Yes: click -> finishAssignment; No: click -> closeContainer | `FlowAction/InboxPolicyTreatyIn` | Modal LayarKasus.tsx: Yes -> POST /kasus/{id}/kirim (services.Kirim menjalankan SetDueTo lebih dulu); No -> tutup modal | dibangun |
| Submit (IsApproved = '0') | `DetailDeptHeadTreatyIn_UW` | .IsApproved = '0'; wadah InputParam.CARI12 != 'TreatyPolicy' | click -> finishAssignment (flow action DeptHeadTreatyIn_UW: post-DT DeptHeadTreatyIn_UW_postDT, InsertHistoryAkseptasiPega) | `FlowAction/DeptHeadTreatyIn_UW`, `DataTransform/DeptHeadTreatyIn_UW_postDT`, `Activity/InsertHistoryAkseptasiPega` | pages/LayarKasus.tsx tombol kirim -> POST /kasus/{id}/kirim -> services.Kirim (PascaAtasan -> riwayat -> connector) | dibangun |
| Submit (IsApproved == 1, tiga varian) | `DetailDeptHeadTreatyIn_UW` | .IsApproved == 1 && identitas operator / pyWorkPage.LetterNo ('TREATYINDEPTHEAD' -> finishAssignment; '' atau identitas tertentu -> nomor polis) | LetterNo 'TREATYINDEPTHEAD': click -> finishAssignment; selain itu: click -> runActivity GeneratePolicyNoTreaty_Act -> localAction ShowPolicyNoTreaty (modal) -> refresh GeneralDeptHeadTreatyIn_UW | `Activity/GeneratePolicyNoTreaty_Act`, `FlowAction/ShowPolicyNoTreaty`, `Section/ShowPolicyNoTreaty_SC`, `Activity/CekLimitTreatyAcc_Act`, `Activity/FetchTreatyGroupOldID` | models.TombolUntuk menurut POSISI: Sec Head -> kirim; Dept Head -> POST /kasus/{id}/nomor-polis (services.TerbitkanNomor) -> Modal nomor polis; GeneratePolicyNoTreaty_Act 11 -> FetchTreatyGroupOldID 3: TreatyGroupID kosong -> pesan "Cannot fetch Treaty Group ID, Contact IT" menahan OK/submit (services.validasiKirim, models.GrupTreatyTakTerbaca) | dibangun - identitas diganti posisi (tiket 05; ketiga syaratnya terdaftar sebagai tempat berperan models.TempatDetailDHSubmit* / TempatGeneralDHSubmit*, tidak dibaca layanan); LetterNo dari CekLimitTreatyAcc_Act tidak dibangun (K2: Sec Head selalu naik); RALAT R6 audit silang P3 (04-10-2026): FetchTreatyGroupOldID langkah 3 dibangun (7.4) |
| OK (nomor polis) | `ShowPolicyNoTreaty_SC` | selalu | click -> closeContainer -> finishAssignment | `FlowAction/DeptHeadTreatyIn_UW` | Modal LayarKasus.tsx OK -> tutup -> POST /kasus/{id}/kirim | dibangun |
| Survey Report (atasan) | `DetailDeptHeadTreatyIn_UW` | pyWorkPage.Quotation.ProportionalType != 'NonProportional'; nonaktif IsSurveyReport 'No' / '' | click -> showHarness HistoricalSurveyReportUW (popup) | `Harness/HistoricalSurveyReportUW`, `FlowAction/InputHistoricalSurveyReportUW` | - | tidak dibangun - K7 |
| .IsApproved (Approval) | `ListSuggest` | selalu (kedua layar) | change/click -> runActivity SetDueTo_act -> refresh CekLimitTreatyAcc_Act -> refresh Protection_Act -> refresh DetailPolicyTreatyIn / DetailDeptHeadTreatyIn_UW / DetailPolicyTreatyInAddendum | `Activity/SetDueTo_act`, `Activity/CekLimitTreatyAcc_Act`, `Activity/Protection_Act` | pages/LayarKasus.tsx radio -> POST /kasus/{id}/hitung {aksi: SetDueTo} -> models.SetDueTo | dibangun sebagian - CekLimitTreatyAcc_Act tidak (K2); Protection_Act varian Data-PolicyTreatyIn tidak ada di korpus (INVENTARIS 1.2) |
| .Suggest | `ListSuggest` | selalu (kedua layar) | change -> refresh Protection_Act -> refresh DetailPolicyTreatyIn / DetailDeptHeadTreatyIn_UW | `Activity/Protection_Act` | isian klien saja (tanpa panggilan backend) | Protection_Act tidak dibangun - varian tidak ada di korpus (INVENTARIS 1.2) |
| .StartDate | `DetailPolicyTreatyIn` | selalu | change -> refresh thisSection dengan pyPreDataTransform SystemSetOneYear_DT (di pyModes behaviors sel) | `DataTransform/SystemSetOneYear_DT` | medan.ts aksi [SystemSetOneYear] -> POST /hitung -> models.SystemSetOneYear (dibangun paket P4) | dibangun (paket P4) |
| .IsSurveyReport / .QuotationData.ProportionalType | `DetailPolicyTreatyIn` | lihat medan.ts | change -> (postValue ->) refresh thisSection TANPA activity | — | isian klien; .QuotationData.IsSurveyReport diterima server hanya bila sel tampil (Quotation bukan NonProportional), pyDefaultValue "No" (models.TerapkanNilaiBawaanSel); .QuotationData.ProportionalType sel pyReadOnly - hanya hasil tombol Enable / Disable (W3) | RALAT R6 audit silang P3 (04-10-2026): dibangun (tanpa panggilan backend) - W3, W4, W5 |
| .StatementType / .QuotationData.NoOfferSlip / .FlagRetroTreaty / .TypeTax | `DetailPolicyTreatyIn` | lihat medan.ts | change -> postValue | — | isian klien; diterima server lewat models.GabungMasukanLayar (daftar izin admin BERURUTAN dengan syarat tampil: FlagRetroTreaty `.ClaimType != XOL Retro`; TypeTax wadah S7 + `.FlagPPH = true`, pyDefaultValue "Inclusive") | RALAT R6 audit silang P3 (04-10-2026): dibangun - W4 sel tersembunyi tidak diterima, W5 nilai bawaan sel |
| .EndDate | `DetailPolicyTreatyIn` | selalu | change -> refresh ProtectDate | `Activity/ProtectDate` | medan.ts aksi [ProtectDate] -> POST /hitung -> models.ProtectDate | dibangun |
| .FlagPPH | `DetailPolicyTreatyIn` | wadah .ClaimType != 'XOL Retro' | change/click -> postValue -> runActivity RemoveTypeTax_ACT | `Activity/RemoveTypeTax_ACT` | medan.ts aksi [RemoveTypeTax] -> POST /hitung -> models.RemoveTypeTax; Save/Submit: FlagPPH yang berubah menjalankan RemoveTypeTax di server (models.GabungMasukanLayar), diterima hanya bila wadah S7 tampil | RALAT R6 audit silang P3 (04-10-2026): dibangun - W4 |
| .QuotationData.MOID | `DetailPolicyTreatyIn` | selalu | change -> postValue -> refresh CheckDataMkt | `Activity/CheckDataMkt` | medan.ts aksi [CheckDataMkt] -> POST /hitung -> models.TerapkanMO; halaman disimpan (Obj-Save langkah 5) | dibangun |
| .IDCurrency | `DetailPolicyTreatyIn` | .IsNewPolicyNonProp != 1 | change -> refresh SetCurrency_act(CURR=.IDCurrency) | `Activity/SetCurrency_act` | medan.ts aksi [SetCurrency] -> POST /hitung -> repository.NamaMataUang + models.SetelNamaMataUang | dibangun |
| .GrossPremium / .GrossClaim | `DetailPolicyTreatyIn` | wadah uang admin | change -> refresh CalculatePremi_Act(Action="PREMIUM" / "CLAIM") | `Activity/CalculatePremi_Act` | medan.ts aksi [CalculatePremi PREMIUM/CLAIM] -> POST /hitung -> models.CalculatePremi | dibangun |
| .PremiOgp / .PremiOnp / .Claim / .SalvageValue / .ExcessLoss / .Deduction1 / .Deduction2 | `DetailPolicyTreatyIn` | wadah uang admin | change / enter -> refresh CountOGPONP_Act | `Activity/CountOGPONP_Act` | medan.ts aksi [CountOGPONP] -> POST /hitung -> models.CountOGPONP | dibangun |
| .RiCommOgp / .ResultOgp1 / .OveriddingCommOgp / .ResultOgp2 / .RiCommOnp / .ResultOnp1 / .OveriddingCommOnp / .ResultOnp2 | `DetailPolicyTreatyIn` | wadah uang admin | change / enter -> refresh CountResult1_Act \| CountResult2Ogp_act \| CountResult1Onp_Act \| CountResult2Onp_act (Data="Pct" untuk persen, "Amount" untuk hasil) -> refresh CountOGPONP_Act | `Activity/CountResult1_Act`, `Activity/CountResult2Ogp_act`, `Activity/CountResult1Onp_Act`, `Activity/CountResult2Onp_act`, `Activity/CountOGPONP_Act` | medan.ts aksi [Count*(Data), CountOGPONP] -> POST /hitung {urutan} -> services.Hitung menjalankan keduanya berurutan atas halaman yang sama | dibangun - DIPERBAIKI 2026-10-03 (semula hanya langkah pertama) |
| .SharePercentage / .ClaimPercentage (grid spreading) | `DetailPolicyTreatyIn` | wadah uang admin | change -> refresh CountSpreading_Act(Index=.pxListSubscript) | `Activity/CountSpreading_Act` | LayarKasus.tsx InputAngka onBlur -> POST /hitung {aksi: CountSpreading, indeks} -> models.CountSpreading (langkah 4-5; 1-3 berlabel // tidak diport); aksi terbuka hanya bila models.SpreadingDariLayar; Save/Submit: %Share berubah -> CountSpreading_Act 4.1 di server | RALAT R6 audit silang P3 (04-10-2026): dibangun - W5, RALAT CountSpreading langkah 1-3 |
| .Installment | `DetailPolicyTreatyIn` | wadah uang admin | change -> refresh FillPaymentInstallment(Installment=.Installment) | `Activity/FillPaymentInstallment` | LayarKasus.tsx -> POST /hitung {aksi: FillPaymentInstallment} -> models.FillPaymentInstallment; Save/Submit: Installment berubah -> FillPaymentInstallment di server (models.TerapkanPemicu); ListInstallment tidak pernah diterima dari layar | RALAT R6 audit silang P3 (04-10-2026): dibangun - W5 |
| .InstallmentPercentage / .Premium (grid angsuran) | `DetailPolicyTreatyIn` | grid S45 `.ListInstallment` pyEditingMode / pyRowEditing / pyNextGenRowEditing readOnly (kedua layar); sel terkunci IsUW \|\| IsClaim \|\| IsShowOpenPolicyMarine \|\| IsOldData = 1 | change -> refresh SetValidateInstallment_Act / CountPctInstallment_Act(idx=.InstallmentNo) | `Activity/SetValidateInstallment_Act`, `Activity/CountPctInstallment_Act` | LayarKasus.tsx grid hanya-baca; POST /hitung SetValidateInstallment / CountPctInstallment ditolak 409 (services/aksiposisi.go); SetValidateInstallment_Act tetap berjalan sebagai CountOGPONP_Act langkah 10 (models.TerapkanPemicu saat sel uang berubah) | RALAT R6 audit silang P3 (04-10-2026): tidak terpicu - grid readOnly (audit silang P3 P2; bunyi lama status: "dibangun", padanan "LayarKasus.tsx InputAngka onBlur -> POST /hitung ...") |
| sel ber-aksi layar atasan (.StartDate, .ShareValue, .PremiOgp, .RiCommOgp, ..., grid spreading/angsuran) | `DetailDeptHeadTreatyIn_UW` | seluruhnya pyReadOnly (1==1) atau pyDisabled always | change -> refresh CountResult1_Act / CountRiCommOgp_act / CountNetPremi_act / SetCurrency_act / CountSpreading_Act / ... (sel terkunci: event tidak pernah terpicu) | — | medan.ts MEDAN_ATASAN_*: jenis tampil, tanpa aksi; server menolak isian (models.GabungMasukanLayar); KECUALI grid SpreadingRiskList subsection NonProp (baris 39-40: terbuka bagi atasan bila FacultativeShare 0/'' - W2) | RALAT R6 audit silang P3 (04-10-2026): tidak dibangun sebagai aksi - tidak terpicu (AC 49-52); grid Proporsional S96 tetap hanya-baca |
| subsection DetailPolicyTreatyOutNonProportional | `DetailPolicyTreatyIn / DetailDeptHeadTreatyIn_UW` | wadah .IsNewPolicyNonProp = 1 && .ClaimType = 'XOL Retro' | lihat INVENTARIS bab 5 Section/DetailPolicyTreatyOutNonProportional | `Section/DetailPolicyTreatyOutNonProportional` | - | tidak dibangun - treaty KELUAR, [keputusan work owner] K8 butir 4 (JSON M_TREATY_OUT: `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`) |
| subsection DetailPoliciesNonProportional (tampil) | `DetailPolicyTreatyIn / DetailDeptHeadTreatyIn_UW` | wadah .IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'; di dalamnya DetailPolicyTreatyInNonProportional bila !TreatyMasterInEDM \|\| IsEDMInputOnNB == true (selalu benar sesudah NonProp langkah 10) | (tampil saja; master TreatyIn dibaca saat pilih bisnis / pra-proses) | `Section/DetailPoliciesNonProportional`, `Section/DetailPolicyTreatyInNonProportional`, `When/TreatyMasterInEDM`, `Activity/InputPolicyTreatyInDetail_NonProp`, `Activity/TreatyRealizationCheckXOLList` | pages/LayarKasus.tsx tampilNonProp -> components/DetailNonProp.tsx; master dari services.siapkanNonProp (PembacaMasterTreaty = repository.MasterXOLDariJSON, K8) | dibangun (K8) |
| .pyTemplateButton Button (refresh TreatyInNonSetTotal) | `DetailPolicyTreatyInNonProportional / ...EDM` | OTHER 1=2 | click -> refresh TreatyInNonSetTotal | `Activity/TreatyInNonSetTotal` | - | tidak dibangun - (a) tampil 1=2 |
| Add / Delete (section SpreadingRiskList NonProp) | `SpreadingRiskList (di DetailPolicyTreatyInNonProportional; layar admin S17 DAN atasan S88)` | InputParam.CARI12 != 'treaty' && != 'TreatyPolicy' && (TreatyIn.FacultativeShare = 0 \|\| '') | click -> addRow / deleteRow (.SpreadingRiskList) | — | components/DetailNonProp.tsx (sunting = pelaku boleh bekerja, admin MAUPUN atasan; spreadingTerbuka) -> diterima server models.SpreadingDariLayar (kedua posisi: wadah NonProp tampil dan FacultativeShare 0/''); bunyi lama padanan: "components/DetailNonProp.tsx (admin, spreadingTerbuka) -> diterima server models.DaftarDariLayar (NonProp: hanya bila FacultativeShare 0/'')" | RALAT R6 audit silang P3 (04-10-2026): dibangun (K8) - W2 atasan juga |
| .SharePercentage / .ClaimPercentage (section SpreadingRiskList NonProp) | `SpreadingRiskList (di DetailPolicyTreatyInNonProportional; layar admin S17 DAN atasan S88)` | terkunci bila TreatyIn.FacultativeShare > 0 | change -> refresh CountSpreading_Act(SharePct, Index=.pxListSubscript, Type) | `Activity/CountSpreading_Act` | DetailNonProp.tsx onBlur -> POST /hitung {aksi: CountSpreading, indeks} -> models.CountSpreading; terbuka di kedua posisi (services.aksiTerbuka -> models.SpreadingDariLayar) | RALAT R6 audit silang P3 (04-10-2026): dibangun (K8) - W2 atasan juga |
| grid ListInstallment (masterDetail InstallmentList) | `DetailPolicyTreatyInNonProportional` | selalu (dalam subsection) | pyRowEditing masterDetail -> FlowAction InstallmentList (Section InstallmentList pyReadOnly) | `FlowAction/InstallmentList`, `Section/InstallmentList` | DetailNonProp.tsx (hanya-baca; rincian per mata uang ditampilkan terbuka); server tidak menerima ListInstallment NonProp dari layar | dibangun (K8) |

⚠️ = rule tidak ada di korpus.
