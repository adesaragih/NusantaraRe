# Audit silang putaran 3 — NB Treaty In (`nbtreatyin`)

Auditor independen (tidak ikut membangun modul ini), prompt putaran 3 bab 3.2. Cabang
`modul/nbtreatyin/p3r-audit`. Audit dimulai dari integrasi `71e6b91a`, lalu digabung dua kali:
`e4d847c0` (R2 F4, R3 K18/F2/F5/F8) dan `97fcb333` (R1 F3: kolom `T_GENERAL_POLIS.EDM_TYPE`,
medan dokumen lama, SuggestList lama). Layar dan pemuat diperiksa pada keadaan sesudah
kedua penggabungan itu (bab 6.5).

**Patokan: korpus XML READ-ONLY** `D:\NUSARE DEV\NusantaraRe\NB Treaty In (Done)\`. XML mentah
dibaca dengan dua lapis escape dibuka. Salinan section tertanam `pyIncludedRuleXML` dibuang,
sama dengan cara kerja `docs/alat/korpus.py`.

Bukti diambil dengan beberapa alat, semuanya hanya membaca:

- **Rantai pemanggilan** dari titik masuk nyata dicari dengan `docs/alat/graf.py`
  (`TITIK_MASUK` = `Flow/InputRealizationTreatyIn` dan `Harness/SFAPortalOpportunities`). Di
  atasnya ada skrip BFS sementara di scratchpad yang mengikuti sambungan graf yang sama.
- **Pembaca properti** dicari dengan `docs/alat/pemakai.py`.
- **Langkah activity** dibaca dari `korpus.json`. Untuk setiap langkah dicatat label, syarat,
  kode transisi, dan `pyPassCurrentParameterPage`.
- **Syarat wadah layar** dibaca dari XML mentah. Tiga pemeriksa hanya-baca dipakai untuk
  admin, atasan, serta NonProp/popup/portal. Setiap temuan besar mereka saya periksa ulang
  sendiri ke XML.

Penamaan dalam dokumen ini:

- Kode syarat/transisi Pega: `2` lanjut, `3` lewati langkah, `6` keluar activity.
- Langkah berlabel `//` = dinonaktifkan.
- Identitas operator yang tertulis di XML disamarkan sebagai `<ID-operator>`.

## Isi

1. [Ringkasan](#1--ringkasan)
2. [Angka status.json](#2--angka-statusjson)
3. [Bagian 1 — rule dibangun: uji berharapan XML](#3--bagian-1--rule-dibangun-uji-berharapan-xml)
4. [Bagian 2 — 30 rule kategori (a)](#4--bagian-2--30-rule-kategori-a)
5. [Bagian 3 — 2 rule kategori (c)](#5--bagian-3--2-rule-kategori-c)
6. [Bagian 4 — setiap layar lawan Section/Harness/FlowAction](#6--bagian-4--setiap-layar-lawan-sectionharnessflowaction)
7. [Temuan WAJIB-BANGUN dan WAJIB-PERBAIKI](#7--temuan-wajib-bangun-dan-wajib-perbaiki)
8. [RALAT dokumen yang diperlukan (untuk orkestrator)](#8--ralat-dokumen-yang-diperlukan-untuk-orkestrator)
9. [Commit dan verifikasi](#9--commit-dan-verifikasi)

---

## 1 · Ringkasan

| Bagian | Hasil |
| --- | --- |
| **1. Rule dibangun** | 48 rule diperiksa: 44 Activity, 2 When, 2 DecisionTable. **17 rule tanpa uji berharapan XML, kini berujian** (19 uji Go + 2 uji vitest). Ada 28 tombol/aksi XML berstatus dibangun; semuanya kini berujian. **5 activity ternyata tidak pernah terpicu** dari layar (7.4), meski status.json menyebutnya "dibangun". |
| **2. Kategori (a)** | 30 rule (bukan 29, lihat bab 2). **Ke-30-nya terbukti tetap tidak dibangun.** Nol WAJIB-BANGUN. Dua bunyi bukti di status.json perlu RALAT kata (bab 8), tetapi kesimpulannya tetap. |
| **3. Kategori (c)** | Bukti WO dirujuk. Medan yang dibutuhkan tidak ada di 39 kolom view dan tidak tercakup F1, jadi **keduanya tetap tidak dibangun**. |
| **4. Layar** | Admin, atasan + ListSuggest, NonProp, popup bisnis, popup SOB, dan portal diperiksa sel per sel, sekitar 500 elemen. **10 selisih kecil diperbaiki dengan uji.** Ada **1 selisih perilaku server** (aksi hitung atasan dapat menyimpan halaman); ukurannya kecil, jadi langsung diperbaiki. **6 temuan WAJIB-PERBAIKI** dengan bukti (bab 7). **Nol menu/halaman di luar Pega.** |

Perbaikan yang sudah dibuat di cabang ini:

1. **`POST /kasus/{id}/hitung` hanya menjalankan aksi sel yang TERBUKA di layar posisi berkas**
   (`services/aksiposisi.go`). Sebelumnya atasan dapat memanggil `CheckDataMkt`, yang
   ber-`Obj-Save`, dan dengan begitu menyimpan halaman tanpa tombol Save. Ini terbukti merah
   dulu oleh `TestHitungHanyaAksiSelTerbukaDiLayarPosisi`.
2. **Grid angsuran `.ListInstallment` kini hanya-baca** (`pyEditingMode=readOnly`). Akibatnya
   `SetValidateInstallment` dan `CountPctInstallment` dikeluarkan dari aksi admin.
3. **Label dan teks tampil disamakan dengan LABEL/pyLabel XML**: grid spreading, grid angsuran,
   ListSuggest, modal tolak, modal nomor polis, portal, caption kotak centang, label
   pendamping, judul grid total NonProp, dan `pyNoSelectionText`.
4. **Dropdown TreatyType** grid NonProp dan grid atasan kini memakai
   `BrowseReinsuranceType_RD` (`.ID` → `.Note`).
5. **Paginasi ListSuggest** lima baris per halaman.
6. **LABEL "EDM"** yang wadahnya mustahil tampil tidak lagi dirender.

---

## 2 · Angka status.json

`docs/alat/status.json`, 176 rule terjangkau, dihitung ulang pada keadaan sesudah merge `e4d847c0`:

| Kategori | status.json (huruf pertama di baris) | HASIL-IMPLEMENTASI bab 6 |
| --- | --- | --- |
| dibangun | **109** (107 + 1 "dibangun sebagian" + 1 "dibangun (K8)") | 109 |
| (a) | **30** | 29 |
| (b) | **35** | 36 |
| (c) | **2** | 2 |

Selisihnya satu baris, yaitu `ReportDefinition/crmOpportunitiesList`. status.json menulisnya
"**(a) + (b)** P44/AC 53", sedangkan HASIL menaruhnya di (b).
`RDBList/SaveTreatyIn` berlaku sebaliknya: status.json menulis "(b) … juga (a)", dan HASIL
juga menaruhnya di (b).

Audit ini memeriksa **ke-30** rule berhuruf pertama (a), termasuk `crmOpportunitiesList`.
Orkestrator perlu memilih satu cara hitung dan menyamakan HASIL dengan status.json.

---

## 3 · Bagian 1 — rule dibangun: uji berharapan XML

Ukuran lulus: ada uji yang **nilai harapannya dibaca atau dihitung tangan dari langkah XML**,
bukan dihitung ulang dengan rumus kode. Kolom "Sumber harapan" mengutip langkah XML-nya. Tanda
**+P3** berarti uji itu ditambahkan audit ini.

Lokasi berkas uji: `M/` = `backend/models/`, `H/` = `backend/handlers/`,
`R/` = `backend/repository/`, dan `F/` = `frontend/`. Uji `R/*_db_test.go` ditulis tetapi
belum dijalankan karena penahan K11.

### 3.1 Activity (44)

| Rule | Uji | Sumber harapan |
| --- | --- | --- |
| `AgentSourceBizTreatyIn_Act` | `H/sumberbisnis_test.go` TestDaftarSumberBisnis | RD `BrowseAgentHierarkiList_RD` filter A (LEADER0 IS NULL), urut ClientName |
| `CalculatePremi_Act` | `M/hitung_test.go` TestCalculatePremiDariGrossPremium; **+P3** `M/hitung_onp_test.go` TestCalculatePremiClaim | langkah 1 (PREMIUM) dan langkah 2-3 (CLAIM): `@divide(.GrossClaim*RNMShareP,100,4)` = 200, lalu CountNetPremi langkah 2 GrossClaim = 500 |
| `CheckDataMkt` | **+P3** `M/pilihbisnis_test.go` TestTerapkanMOMenurutCheckDataMkt; **+P3** `H/aksisel_test.go` TestCheckDataMktMengisiMODanMenyimpan | langkah 2-4 (11 medan, MOID kosong → kosong); langkah 5 `Obj-Save` (tersimpan) |
| `CountNetPremi_act` | `M/hitung_test.go` TestRantaiNetPremiDanBalance, TestMasterTidakTersediaLangkahnyaDilewati; **+P3** TestCountNetPremiNonPropMemakaiRNMShare | langkah 4-6 (NetPremium 830, Balance*); langkah 3 NonProp `@divide(.Claim,RNMShare,4)*100` = 75, bukan 60 |
| `CountOGPONP_Act` | `M/hitung_test.go` TestPremiKosongMenolkanLimaMedan, TestKlaimMengisiClaimTypeSOA, TestCountOGPONPTanpaMaster; `H/alur_test.go` TestHitungUrutanActionSet | langkah 1-2, 8 (`SOA`/`Claim`), 9-10 |
| `CountOverridingCommOgp_Act` | **+P3** `M/hitung_onp_test.go` TestAktivitasLayarDeptHeadMembulatkanEmpatDesimal | langkah 3 `@Math.divide(.ResultOgp2,.PremiOgp,4)*100` = 6,67; langkah 5 NetPremium 2799,9 |
| `CountOverridingCommOnp_Act` | **+P3** sama | langkah 3-4 (ONP) |
| `CountPctInstallment_Act` | **+P3** `M/hitung_onp_test.go` TestCountPctInstallmentBarisIdx | langkah 1 `@Math.divide(@toDecimal(.Premium),@toDecimal(.BalanceDueTo),4)*100` = 16,67; PaymentTotal = Premium |
| `CountResult1Onp_Act` | **+P3** TestCountResult1OnpPctDanAmount, TestCountResult1OnpKeluarTanpaMenghitung | langkah 1-5: Pct 200 / NetPremium 1800; Amount 15; keluar langkah 1/2/3 tanpa CountNetPremi; pesan `ErrorMsg1` > 100 |
| `CountResult1_Act` | `M/hitung_test.go` TestRantaiNetPremiDanBalance, TestGrossPremiumMemakaiRNMShareP, TestRiCommOgpSeratusTidakDihitung, TestRiCommOgpLebihSeratusMemasangPesan, TestPembagiNolGagalTerang; `H/alur_test.go` TestHitungUrutanActionSet (Amount) | langkah 1-7 |
| `CountResult2Ogp_act` | `M/hitung_test.go` TestRantaiNetPremiDanBalance (Pct); **+P3** TestCountResult2OgpAmount | langkah 3 Pct (ResultOgp2 50); langkah 4 Amount (OveriddingCommOgp 3, NetPremium 970) |
| `CountResult2Onp_act` | **+P3** TestCountResult2OnpPctAmountDanKosong | langkah 2 (kosong → keluar), 3 Pct 100, 4 Amount 2,5 |
| `CountRiCommOgp_act` | **+P3** TestAktivitasLayarDeptHeadMembulatkanEmpatDesimal | langkah 3 (33,33, dibulatkan 4 desimal sebelum ×100); langkah 1, 2, 4 `//` |
| `CountRiCommOnp_act` | **+P3** sama | langkah 2 (33,33); langkah 3 CountNetPremi |
| `CountSpreading_Act` | `M/hitung_test.go` TestSpreadingDuaBarisRata; `M/layar_test.go` TestSpreadingBagiRataPresisiSepuluh; `H/nonprop_test.go` TestNonPropPilihBisnisHitungSimpanBacaKembali | langkah 4-5 |
| `FetchTreatyGroupOJK` | `H/alur_test.go` TestTanggaPenuhDanNomorPolisSekali; **+P3** `H/aksisel_test.go` TestNomorPolisMemakaiOJKTerisiDanMengisiGrupLama | GeneratePolicyNoTreaty_Act langkah 13 (hanya bila `OJKBusinessID==""`) dan 28 (`".T"+OJKBusinessID`) |
| `FetchTreatyGroupOldID` | **+P3** sama | GeneratePolicyNoTreaty_Act langkah 11 (bila kosong), FetchTreatyGroupOldID langkah 6 |
| `FillPaymentInstallment` | `M/hitung_test.go` TestFillPaymentInstallmentTigaKali | langkah 3.1-3.5 (33,3333 / 33,3334) |
| `GeneratePolicyNoTreaty_Act` | `M/tangga_test.go` TestRakitNomorPolis, TestTipeNomorPolis, TestTanggalProduksiNomor; `M/tutupbuku_test.go`; `H/alur_test.go` TestTanggaPenuhDanNomorPolisSekali; `R/nomorpolis_db_test.go` (K11) | langkah 5.3, 7-9, 25-28 |
| `InputPolicyTreatyInDetail_NonProp` | `M/nonprop_detail_test.go` TestInputDetailNonProp, TestInputDetailNonPropRetro; `H/nonprop_test.go` (2 uji) | langkah 6-23 atas master fiktif, hitung tangan |
| `InputPolicyTreatyInDetail_preACT` | `H/alur_test.go` TestPilihBisnis; `H/logika_test.go` TestPilihBisnisKomisiOgpDariRIONRView; `M/layar_test.go` TestKunciCariBisnis; `M/nonprop_detail_test.go` TestPPNPPHLapisanXOL; `R/kontrak_db_test.go` (K11) | langkah 3-6, 8, 14, 17, 18 |
| `InputPolicyTreatyInPost_Act` | `H/alur_test.go` TestAdminMenolakDiselesaikanDitolak, TestTanggaPenuhDanNomorPolisSekali, TestNourutUsulanSamaDenganSubskripSuggestList; `H/logika_test.go` TestCekPolisSerupaHanyaSaatAdminMenyetujui | langkah 2-4 |
| `InputPolicyTreatyInPre_Act` | `M/tangga_test.go` TestPraprosesAdmin*; `M/tutupbuku_test.go`; `H/logika_test.go` TestTanggalProduksiDariTanggalClosing; `H/nonprop_test.go` TestNonPropCekDaftarXOLSaatDibuka; `M/nonprop_master_test.go` TestPerluCekDaftarXOL | langkah 2-4, 9, 10 |
| `InputQuotation_PreAct` | `M/sumberbisnis_test.go` (TombolSOB / PostDT) | langkah 1 (`Acton=="SOB"` → btnSOB_DT) |
| `InsertHistoryAkseptasiPega` | `H/alur_test.go` (riwayat per submit); `R/riwayat_db_test.go` (K11) | RDB `InsertHistoryAkseptasiPega_Sql` |
| `InsertToTreatyXOLList` | `M/nonprop_test.go` TestInsertToTreatyXOLList, …TanpaPPH | langkah 3.x, hitung tangan |
| `InsertToTreatyXOLListRetroShare` | `M/nonprop_test.go` TestRetroShare* (3) | langkah 2.x |
| `ProtectDate` | **+P3** `M/pilihbisnis_test.go` TestProtectDatePesanBilaAkhirSebelumMulai; **+P3** `F/medan.test.ts` (aksi sel `.EndDate`) | langkah 1-2: pesan verbatim "End Date Cannot be less than Start Date" bila `@CompareDates(.StartDate,.EndDate)` |
| `RemoveTypeTax_ACT` | **+P3** TestRemoveTypeTaxHanyaBilaFlagPPHTidakBenar; **+P3** `F/medan.test.ts` (`.FlagPPH`) | langkah 2 `FlagPPH==false` → Property-Remove `.TypeTax` |
| `SaveViewSuggest` | `M/usulan_test.go` (4 uji); `H/alur_test.go` TestNourutUsulanSamaDenganSubskripSuggestList | langkah 2 "UNTUK TREATY" |
| `SetCurrency_act` | **+P3** `H/aksisel_test.go` TestSetCurrencyMengisiNamaMataUangDariID; **+P3** `F/medan.test.ts` (`.IDCurrency`) | langkah 1-3 `.Currency = Currency.pxResults(1).HASIL1`; tanpa Obj-Save |
| `SetDueTo_act` | `M/hitung_test.go` TestBalanceNegatifDueToNol, TestRantaiNetPremiDanBalance | ≥ 0 → 1, < 0 → 0 |
| `SetPPNPPH` | `M/setppnpph_test.go`; `M/hitung_test.go` TestSetPPNPPHBersyaratFlagAtauPKP; `H/logika_test.go` TestSetPPNPPHMembacaStatusPKPAgen | langkah 2, 4 |
| `SetReinstatementPct` | `M/nonprop_master_test.go` TestSetReinstatementPctSatuMDP | langkah 2-3 |
| `SetTreatyCurrencyID` | `H/alur_test.go` TestPilihBisnis (IDCurrency dari nama) | preACT langkah 4 → langkah 1-3 |
| `SetTreatyIn_Act` | `R/masterxol_test.go` (4); `R/masterxol_penjaga_test.go` (3); `M/nonprop_master_test.go` TestTreatySetReinstatement; `R/masterxol_db_test.go` (K11) | langkah 3-5, 13 (F1) |
| `SetValidateInstallment_Act` | `M/hitung_test.go` TestAngsuranLebihSeratusMemasangPesan | langkah 2-4 (pesan verbatim) |
| `SetValue_Act` | `H/alur_test.go` TestPilihBisnis | langkah 1 → preACT |
| `TreatyInputPctCommSpreading` | `M/komisi_test.go` (6); `H/logika_test.go` TestPilihBisnisKomisiOgpDariRIONRView | langkah 2-2.1.1.1 (RIONR menimpa RIOGR) |
| `TreatyNonPropSetSpreading` | `M/nonprop_test.go` TestTreatyNonPropSetSpreading* (3) | langkah 2-7 |
| `TreatyRealizationCheckDuplicate` | `H/alur_test.go` TestPolisSerupaMenahan; `H/logika_test.go` TestCekPolisSerupaHanyaSaatAdminMenyetujui | langkah 1-6, pesan verbatim |
| `TreatyRealizationCheckXOLList` | `M/nonprop_master_test.go` TestCekDaftarXOL; `H/nonprop_test.go` TestNonPropCekDaftarXOLSaatDibuka | langkah 2-8 |
| `TreatySetReinstatement` | `M/nonprop_master_test.go` TestTreatySetReinstatement | langkah 1-1.1 |
| `serviceInsertArasapas_act` | `H/alur_test.go` TestKonversiSesudahSelesai | CARI1/2/3/21 (KEPUTUSAN-RONDE-12 butir 7) |

### 3.2 When (2) dan DecisionTable (2)

| Rule | Uji | Sumber harapan |
| --- | --- | --- |
| `When/NopolisEmpty` | `M/tangga_test.go` TestTanggaDeptHeadMenyetujuiSelesai; `H/alur_test.go` TestTanggaPenuhDanNomorPolisSekali | `PolicyNo = ""` → Decision8 → Assignment3 |
| `When/TreatyMasterInEDM` | **+P3** `M/pilihbisnis_test.go` TestTreatyMasterInEDMTigaNilai; `M/nonprop_detail_test.go` (EDMState 2) | `EDMState` = "1" OR "2" OR "3" |
| `DecisionTable/BusinessType_DeT` | `M/penggolong_test.go` TestKe128KodeSamaDenganSistemLama, …BerhentiDiBarisPertama, …BawaanUnknown, …MembandingkanTeks | 128 kode dibangkitkan dari baris XML oleh skrip terpisah (independen dari kode) |
| `DecisionTable/isApproved` | `M/tangga_test.go` TestIsApproved* (3) | kolom `IsApproved = "0"` → No, bawaan YES |

### 3.3 Tombol dan aksi sel (`docs/alat/tombol.json`, entri dibangun)

| Tombol / aksi | Uji |
| --- | --- |
| Create opportunity | `H/alur_test.go` TestBuatKasusDiAntreanAdmin |
| Filter / `.FilterTermForOpportunity` | `H/portal_test.go` TestDaftarPortalSesuaiGetListOpportunity, TestGerbangDaftarPortal |
| `.Name` (tautan berkas) | `H/alur_test.go` TestBukanAnggotaAntreanDitolak (buka hanya-baca) |
| Select Source Of Business / Choose (pohon SOB) | `H/sumberbisnis_test.go` (paket F4); `F/components/pilihSumberBisnis.test.ts` |
| Choose Business / Choose (popup) | `H/alur_test.go` TestPilihBisnis; `H/nonprop_test.go` |
| Enable / Disable Input Type | `M/layar_test.go` TestEnableDisableSelaluNol |
| Save | `H/alur_test.go` TestSimpanDrafHanyaAdmin, TestMedanWajibMenahanKirimDanSimpan |
| Add / Delete grid spreading | **+P3** `M/layar_daftar_test.go` TestGabungMasukanDaftarMenurutGridXML; `F/nonprop.test.ts` |
| Submit (IsApproved 1 / 0), Yes / No | `H/alur_test.go` TestTanggaAdminMenyetujuiNaikKeSecHead, TestAdminMenolakDiselesaikanDitolak; `M/layar_test.go` TestTombolPerPosisi |
| Submit atasan (tiga varian), OK nomor polis | `H/alur_test.go` TestTanggaPenuhDanNomorPolisSekali; `M/layar_test.go` TestTombolPerPosisi |
| `.IsApproved` (Approval, runActivity SetDueTo_act) | `M/hitung_test.go` TestBalanceNegatifDueToNol; **+P3** `H/aksisel_test.go` TestHitungHanyaAksiSelTerbukaDiLayarPosisi (SetDueTo atasan 200) |
| `.StartDate` (SystemSetOneYear_DT) | `H/tanggal_test.go`; `F/tanggalMulai.test.ts` |
| `.EndDate`, `.FlagPPH`, `.QuotationData.MOID`, `.IDCurrency` | **+P3** `F/medan.test.ts` (aksi sel) + uji model/HTTP di 3.1 |
| `.GrossPremium` / `.GrossClaim`, sel uang OGP/ONP | `F/medan.test.ts` (urutan action set) + 3.1 |
| Grid spreading `.SharePercentage` / `.ClaimPercentage` | `M/hitung_test.go` TestSpreadingDuaBarisRata |
| `.Installment` | `M/hitung_test.go` TestFillPaymentInstallmentTigaKali |
| Grid angsuran `.InstallmentPercentage` / `.Premium` | **Tidak terpicu**, karena grid `readOnly` (7.4). Ditolak server: **+P3** TestHitungHanyaAksiSelTerbukaDiLayarPosisi |
| Subsection NonProp, grid NonProp | `H/nonprop_test.go`; `M/nonprop_*_test.go`; `F/nonprop.test.ts` |

### 3.4 Uji yang ditambahkan audit ini

Go:

- `backend/models/hitung_onp_test.go`: 8 uji.
- `backend/models/pilihbisnis_test.go`: 4 uji.
- `backend/models/layar_daftar_test.go`: 1 uji.
- `backend/handlers/aksisel_test.go`: 4 uji, termasuk TestHitungHanyaAksiSelTerbukaDiLayarPosisi.

Vitest:

- `frontend/medan.test.ts`: 2 uji baru (aksi sel; teks tampil XML).
- `frontend/labels.test.ts`: 7 uji, baru.
- `frontend/nonprop.test.ts`: 2 uji (judul grid total).
- `frontend/paginasi.test.ts`: 2 uji, baru.

Seluruh harapan diambil dari langkah/LABEL XML yang dikutip di komentar uji.

---

## 4 · Bagian 2 — 30 rule kategori (a)

Untuk setiap rule ditulis rantai terpendek dari titik masuk nyata menurut `graf.py`, lalu
alasan XML bahwa rule itu tidak pernah berjalan atau bahwa efeknya dibaca nol rule NB.

Singkatan rantai yang dipakai di tabel:

| Singkatan | Rantai lengkap |
| --- | --- |
| **ATS** | `Flow/InputRealizationTreatyIn` →(connector Assignment1→Decision1) `FlowAction/DeptHeadTreatyIn_UW` → `Activity/InputPolicyTreatyInPre_Act` |
| **ADM** | `Flow/InputRealizationTreatyIn` →(Assignment2→Decision3) `FlowAction/InboxPolicyTreatyIn` → `Section/GeneralPolicyTreatyIn` → `Section/DetailPolicyTreatyIn` |
| **XOL** | ATS → langkah 10 `TreatyRealizationCheckXOLList` → langkah 4 `SetTreatyIn_Act` |
| **SOB** | ADM →(showHarness `.pyTemplateInputBox`) `Harness/SOB` → `Section/SourceHierarki` |
| **PRT** | `Harness/SFAPortalOpportunities` → `Section/SFAPortalOpportunitiesHeader` |

| # | Rule | Rantai | Bukti tidak terjangkau / nol pembaca (XML) | Putusan |
| --- | --- | --- | --- | --- |
| 1 | `Activity/CheckDuplicateOffer` | XOL → langkah 14 | Dipanggil `SetTreatyIn_Act` langkah 14 dengan `pyPassCurrentParameterPage=false` tanpa parameter, sehingga `Param.ID` kosong. Langkah 1–4 berlabel `//`. Langkah 5 `TreatyWarning.CAIREINSFACIN = Param.ID` = "". Langkah 6 GetCountClaim `masterid={…}` menghasilkan 0, jadi syarat langkah 7–8 `HasilClaim.pxResults(1).CARI1>0` selalu salah. | tetap (a) |
| 2 | `Activity/ConvertHistoryDate` | XOL → langkah 12 | Hanya mengalang `TreatyIn.CommentList`. `pemakai.py "CommentList"` → 2 rule: dirinya sendiri dan `SetTreatyIn_Act` langkah 9, yang berprasyarat `revisionstate==1` (lihat #18). Nol Section/Activity/When NB yang membacanya. | tetap (a) |
| 3 | `Activity/InputParamUploadReas_act` | ATS → langkah 1 `SetCategoryAttach` → 4.2 | Dipanggil dengan `pass=true` dan mengisi `AttachShowList`. Hasilnya hanya dipakai `SetCategoryAttach` 4.3 (`.CountAttach`). `pemakai.py "CountAttach"`: pembaca lain hanya `CheckRISLIP`, `CheckRISLIP_EDM`, `CheckSpreadingProtect_ACT`, ketiganya tak-terjangkau (`Protection_Act`, sambungan PUTUS di graf.py). | tetap (a) |
| 4 | `Activity/InputPolicyTreatyEDMDetail_NP` | ADM → popup bisnis → `SetValue_Act` → preACT 16 → NonProp langkah 11 | Syarat varian EDM mustahil (sudah terbukti WO): `TreatyMasterInEDM` dan `IsEDMInputOnNB != true`, padahal langkah 10 mengisi `IsEDMInputOnNB=true` tepat saat `TreatyMasterInEDM`. | tetap (a) |
| 5 | `Activity/InsertToTreatyOutXOLList` | via #4 atau `InputPolicyTreatyOutDetail_NonProp` 27 | Pemanggilnya hanya #4 (mustahil) dan jalur treaty keluar (K8 butir 4, (b)). | tetap (a) |
| 6 | `Activity/SearchHierarkiSourceBizAgentTreatyIn_Act` | SOB → runActivity Choose | Langkah 1–4 hanya menulis `pyWorkPage.OfferTreatyIn.QuotationData.{SourceOfBusiness, SobName, SobLeader0, SobLeader1, Email}`. `pemakai.py "OfferTreatyIn"`: selain dirinya hanya `AgentSourceBizTreatyIn_Act` (membaca `btnQuotation`, tidak ditulis rule ini), `Harness/SOB`, dan `Section/SourceHierarki` (nama kelas, bukan properti). `QuotationData.Email` dibaca nol rule. | tetap (a) |
| 7 | `Activity/SetCategoryAttach` | ATS → langkah 1 | Untuk `BusinessFac` "T": langkah 1 (CategoryAttach_SQL → `Attachment`) dan 4 (`.CountAttach`) berjalan. Pembacanya sama dengan #3, yaitu hanya rule `Protection_Act` yang tak-terjangkau. | tetap (a) |
| 8 | `Activity/TreatyInInputVis` | XOL → langkah 2 | Dipanggil `Add=1`. Langkah 1 menulis `OutputParam.DATASHOW/ERRMSG` dan `TreatyIn.ID="UnknownId"`, lalu **transisi `1==1` → 6 (keluar)**, sehingga langkah 3 `Page-Remove` tidak berjalan (RALAT bunyi status.json, bab 8). `SetTreatyIn_Act` langkah 3 menimpa `TreatyIn.ID` dan `ERRMSG`. `DATASHOW` dibaca nol rule. | tetap (a) |
| 9 | `Activity/TreatyInNonSetTotal` | ATS → `GeneralDeptHeadTreatyIn_UW` → NonProp → tombol | Satu-satunya pemicu: `.pyTemplateButton` di `DetailPolicyTreatyInNonProportional` / …EDM, `pyVisible` `1=2`. | tetap (a) |
| 10 | `DataTransform/btnCedingCO_DT` | ADM → `InputQuotation_PreAct` langkah 2 | Prasyarat `Param.Acton=="Ceding"`. `pemakai.py "Acton"`: satu-satunya pengirim (`Section/DetailPolicyTreatyIn`) mengirim `Acton=SOB`. | tetap (a) |
| 11 | `DataTransform/setCategoryAttachment_DT` | ATS → `SetCategoryAttach` 3 | Langkah 3 prasyarat `BusinessFac=="F"\|\|=="T"` dengan T→3 (lewati). Kasus NB = "T". | tetap (a) |
| 12 | `FlowAction/Installments_ReadOnly` | ATS → NonProp → …EDM | Hanya `pyEditAction` di `DetailPolicyTreatyInNonProportionalEDM` (#22, mustahil). | tetap (a) |
| 13 | `RDBList/AttachmentLife` | ATS → `SetCategoryAttach` 2 | Langkah 2 prasyarat T→3. Kasus "T". | tetap (a) |
| 14 | `RDBList/BrowseClientEmail_SQL` | SOB → #6 langkah 1 | Hasilnya hanya menulis `OfferTreatyIn.QuotationData.Email` (#6), yang dibaca nol rule. | tetap (a) |
| 15 | `RDBList/CategoryAttach_SQL` | ATS → `SetCategoryAttach` 1 | Berjalan, tetapi halaman `Attachment` dibaca hanya rule tak-terjangkau (#3). | tetap (a) |
| 16 | `RDBList/GenerateNoPolicy` | `Flow/InputRealizationTreatyIn` (Utility1) → `SaveJsonPolisTreatyIn_Act` langkah 4 | RDB-List tanpa BrowsePage, jadi hasil masuk `PolicyTreatyIn.pxResults`. `pemakai.py "PolicyTreatyIn\.pxResults"`: satu-satunya pembaca adalah `GeneratePolicyNoTreaty_Act` langkah 24, berlabel `//`. Pemanggilnya sendiri (b) AC 16. | tetap (a) |
| 17 | `RDBList/GetCountClaim` | #1 langkah 6 | `masterid = ""` menghasilkan 0 (#1). | tetap (a) |
| 18 | `RDBList/GetCurrentDate` | XOL → langkah 8 | Prasyarat `param.revisionstate==1`. `TreatyRealizationCheckXOLList` dipanggil `pass=false` dari Pre_Act 10, tanpa parameter. Langkah 4-nya memanggil `SetTreatyIn_Act` dengan **`pyPassCurrentParameterPage=true`**, sehingga halaman parameternya hanya berisi `Param.ID` (= NoOffer, langkah 3), dan `revisionstate` kosong (RALAT bunyi, bab 8). | tetap (a) |
| 19 | `ReportDefinition/BrowseAgentNusaRe_RD` | SOB → sumber autocomplete `SearchSOB.CARI1` | `pemakai.py "SearchSOB"`: hanya dikosongkan `btnSOB_DT` / `btnCedingCO_DT` dan dipakai sel itu sendiri. Tanpa action set. TreeGrid tidak memakainya. | tetap (a) |
| 20 | `ReportDefinition/BrowseCedingCo_RD` | SOB → `AgentSourceBizTreatyIn_Act` langkah 2 | Prasyarat `pyWorkPage.OfferTreatyIn.QuotationData.btnQuotation=="CedingCo"`, ditulis nol rule (`btnCedingCO_DT` menulis `pyWorkPage.Quotation.btnQuotation`). | tetap (a) |
| 21 | `ReportDefinition/BrowseTREATY_IN` | #1 langkah 1 | Langkah 1–4 `CheckDuplicateOffer` berlabel `//`. | tetap (a) |
| 22 | `ReportDefinition/crmOpportunitiesList` | PRT → `SFAPortal_OpportunitiesList` `pyGridRDName` | Hanya metadata sel label kolom. Grid memakai `pyRDName=GetListOpportunity`. Rujukan lain adalah parameter tombol Stage view (wadah `1=2`). Juga (b) P44/AC 53. | tetap (a)+(b) |
| 23 | `Section/DetailPolicyTreatyInNonProportionalEDM` | ATS → `DetailPoliciesNonProportional` | Wadah S2 `TreatyMasterInEDM && IsEDMInputOnNB != true`, mustahil (#4). | tetap (a) |
| 24 | `Section/Installments_ReadOnly` | via #12 | Hanya lewat #12. | tetap (a) |
| 25 | `When/IsClaim` | ATS → `DetailDeptHeadTreatyIn_UW` | `pyWorkPage.pyWorkIDPrefix = "CLM-"`, sedangkan kasus NB berawalan `NB-`. Pemakainya `pyReadOnlyCondition` sel `.DueDate` / `.InstallmentPercentage` di grid yang toh `readOnly` (7.4). | tetap (a) |
| 26 | `When/isClaimTreaty` | Flow connector Decision5→Assignment5 | Nol connector ke Decision5. Daftar connector Flow dicek: hanya `dari: Decision5`. | tetap (a) |
| 27 | `When/isSellingModeB2B` | PRT → syarat tombol Create | `getDataSystemSetting("PegaCRM-","SellingMode")` (nilai tidak di korpus). Ketiga When saling meniadakan, masing-masing hanya memilih satu dari tiga tombol "Create opportunity" yang semuanya `createWork`. | tetap (a) |
| 28 | `When/isSellingModeB2BB2C` | sama | sama | tetap (a) |
| 29 | `When/isSellingModeB2C` | sama | sama. Varian B2C ber-`WorkClass_OpportunityInd`, `D_crmAppExtPage` tidak di korpus. | tetap (a) |
| 30 | `When/pyIsIpadOrDesktop` | PRT → `pyContainerVisibleWhen` wadah Create | `pyIsIPad` (tidak di korpus) \|\| `pxRequestor.pxDeviceType = desktop`, benar di peramban desktop. Wadahnya dibangun tampil-selalu. | tetap (a) |

**Hasil bagian 2:** nol rule (a) yang terjangkau dan efeknya dibaca. **Nol WAJIB-BANGUN.**

---

## 5 · Bagian 3 — 2 rule kategori (c)

Keabsahan (c) sudah diputuskan WO. Di sini hanya dicek bahwa medan yang dibutuhkan memang tidak
tersedia.

| Rule | Medan yang dibaca rule terjangkau dari halaman hasilnya | Di 39 kolom view `TREATYINDETAILJOINEDM`? | Di F1 (`models/masterxol.go`)? |
| --- | --- | --- | --- |
| `Activity/FetchMasterTreatyIn` (`TreatyInputPctCommSpreading` langkah 1, hanya jalur bukan-NonProp) | `Limits().Detail().TreatyType/TreatyGroup/RIOGR/RIONR` (langkah 2.1–2.1.1.1, **dibangun dari kolom view**); `SpreadingTotalPct`, `SpreadingTypeID`, `SpreadingType` (2.1.1.1); `TreatyIn.RNMShareP`, `BrokeragePercentP` (CalculatePremi 1-2, CountNetPremi 2, CountResult1 3 di jalur Proporsional) | **Tidak.** Ke-39 kolom adalah 33 kolom RD ditambah BROKERAGE, COMMENCEMENT, TERMINATION, RIOGR, RIONR, RNM_SHARE. Nama medan Spreading* dan RNMShareP/BrokeragePercentP tidak ada. `RNM_SHARE`/`BROKERAGE` adalah kolom lain yang maknanya masih `[terbuka]` (PERTANYAAN-untuk-DBA, AC 26). | **Tidak.** F1 hanya memuat medan jalur NonProp (`RNMShare`, bukan `RNMShareP`; tanpa Spreading*; tanpa BrokeragePercentP). |
| `RDBList/BrowseTreatyInDetailJoinEDM` (preACT langkah 9 → 10 `adoptJSONObject` → 13) | `TreatyIn.INSTALLMENT(n).Currency` dan `.InstallmentList(m).{Installment, PaymentDate, InstallmentPct, AcceptStatus}` (preACT 13–13.1.1) | **Tidak.** View hanya punya `INSTALLMENTNO` (cacah), tanpa daftar angsuran. | **Tidak.** F1 `Installment.InstallmentList.*` adalah properti **`Installment`** dari JSON `M_TREATY_IN` / `_EDM`. Rule ini membaca **`INSTALLMENT`** (nama properti Pega peka huruf) dari JSON **`M_TREATY_IN_DETAIL_EDM`**, tabel di luar dua tabel pengecualian K8. |

**Kedua rule (c) tetap tidak dibangun.**

---

## 6 · Bagian 4 — setiap layar lawan Section/Harness/FlowAction

Diperiksa sel per sel: jalur, label, wajib (`pyRequired*`), terkunci/nonaktif
(`pyReadOnly*` / `pyDisabled*`), syarat tampil sel DAN setiap wadah (`pyContainerVisibleWhen`),
dan setiap tombol/aksi (event → aksi → rule → parameter, berurutan).

| Layar | Diperiksa | Cocok | Selisih | Sah tidak dibangun |
| --- | --- | --- | --- | --- |
| Admin `DetailPolicyTreatyIn` + `GeneralPolicyTreatyIn` + FA `InboxPolicyTreatyIn` | 145 | 82 | 26 | 37 |
| ListSuggest (kedua layar) | 12 | 10 | 2 | – |
| `SpreadingRiskList`, `InstallmentList`, `PolicyTreatyInDeclineConfirm` | 40 | 26 | 9 | 5 |
| Atasan `DetailDeptHeadTreatyIn_UW` + `GeneralDeptHeadTreatyIn_UW` + FA + `ShowPolicyNoTreaty_SC` | 145 | 98 | 30 (+S1) | 31 |
| Subsection NonProp (`DetailPoliciesNonProportional` → `DetailPolicyTreatyInNonProportional`) | 103 | 52 | 43 | 8 |
| Popup pilih bisnis (`Harness/Section BusinessAndSOBList`) | 34 | 3 | 28 | 3 |
| Popup SOB (`Harness/SOB`, `SourceHierarki`, FA `AgentSourceBizDetails`) | 13 | 10 | (F4, kini sudah digabung) | 3 nol efek |
| Portal (`SFAPortalOpportunities`, `SFAPortal_Opportunities`, `…Header`, `…List`) + menu | 30 | 11 | 9 | 10 |

### 6.1 Yang cocok (ringkas)

Bagian berikut sama dengan XML:

- **Medan wajib.** `models/layar.go` `medanWajibAdmin` / `medanWajibAtasan` sama dengan
  `pyRequired`/`pyRequiredWhen`, termasuk syarat wadah S19 dan S14 admin serta S90 atasan
  (`ResultOnp1` wajib di atasan).
- **Syarat tampil sel dan wadah** di kedua layar: S7, S8, S14, S15, S17, S19, S49 admin; S78,
  S85, S86, S88, S90 atasan. Juga wadah NonProp `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL
  Retro'` dan Share Facultative `FacultativeShare>0`.
- **Urutan dan parameter action set refresh admin.**
- **Tombol:**
  - Submit, ditambah modal tolak Yes/No.
  - Submit atasan: varian identitas diganti posisi (K5/K12), LetterNo (K2).
  - Save, Choose Business, Enable/Disable, Add/Delete, dan OK nomor polis.
- **Menu.** Migrasi 968 hanya `UPDATE M_NAV_MENU … WHERE KODE='nbtreatyin'`; `menu.ts` satu
  halaman; `rute.tsx` hanya portal + layar kasus. **Nol menu/halaman di luar Pega.**

### 6.2 Selisih yang DIPERBAIKI di cabang ini (dengan uji)

| # | Elemen | Bukti XML | Perbaikan |
| --- | --- | --- | --- |
| P1 | `POST /hitung` menjalankan aksi apa pun di posisi mana pun. `CheckDataMkt` (Obj-Save) membuat atasan dapat menyimpan halaman. | Semua sel ber-refresh `DetailDeptHeadTreatyIn_UW` ber-`pyReadOnly`. Yang terbuka hanya radio `ListSuggest .IsApproved` (runActivity `SetDueTo_act`). | `services/aksiposisi.go` menjawab 409 untuk aksi di luar layar posisi. Uji `H/aksisel_test.go` TestHitungHanyaAksiSelTerbukaDiLayarPosisi (merah dulu: Suggest atasan tersimpan lewat CheckDataMkt). |
| P2 | Grid angsuran admin tersunting (DueDate, %, Premium) | `DetailPolicyTreatyIn` S45 `.ListInstallment`: `pyEditingMode`, `pyRowEditing`, dan `pyNextGenRowEditing` = `readOnly`. Bandingkan grid spreading `readWrite`. | `LayarKasus.tsx` grid hanya-baca. `SetValidateInstallment`/`CountPctInstallment` keluar dari aksi admin. |
| P3 | Judul kolom grid spreading | LABEL "Type Treaty", "% Share", "Premium", "% Share", "Claim"; kaki "Total" (`SpreadingRiskList` S2, admin, atasan) | `labels.ts` + `F/labels.test.ts` |
| P4 | Judul grid angsuran + judul wadah | "No", "Due Date", "% Installment", "Premium Nusantara Re", "Total Payment"; S45 `pyTitle` "Installment Data Information" | sama |
| P5 | Kolom ListSuggest; paginasi | LABEL "PIC"; `pyPageMode=Numeric`, `pyPageSizeOther=5`, `pyGridPaginator` | `labels.ts`; `paginasi.ts` + `F/paginasi.test.ts` |
| P6 | Modal tolak dan nomor polis | FA `pyLabel` "Confirm Decline NB" + LABEL "Are you sure you want to decline this NB"; FA "Show PolicyNo" + `pyID` LABEL "telah diaksep menjadi" `PolicyNo` | `labels.ts`, `LayarKasus.tsx` |
| P7 | Portal | LABEL "Opportunity"; kolom "Offer No", "Group Business", "Insured Name", "Marketing", "Status"; `pyPlaceholder` "NB-1234 or Name" | `labels.ts`, `PortalNBTreatyIn.tsx` |
| P8 | Caption kotak centang; label pendamping | `.FlagRetroTreaty` "Overiding Commision", `.FlagPPH` "With Tax" (admin) / "Include Tax" (atasan); LABEL "Q", "/", "Of" di depan `.Quartal`, `.YearOfQuartal`, `.LayerPartType` | `medan.ts` + `F/medan.test.ts` |
| P9 | Dropdown TreatyType NonProp dan grid atasan | `SpreadingRiskList` dan `DetailDeptHeadTreatyIn_UW` `.TreatyType`: `pyListSource=reportdefinition` `BrowseReinsuranceType_RD` (`.ID` / `.Note`), `pyNoSelectionText` "Choose". Grid admin Proporsional `pyListSource=pageList` `ListSpreading.pxResults` sudah benar. | `LayarKasus.tsx`, `DetailNonProp.tsx` (`acuan.jenisReas`) |
| P10 | Grid total NonProp; LABEL "EDM" | Judul di atas `.Value`, `.Currency` tanpa judul, kecuali Total Spreaded; `InstallmentList` kolom ke-3 tanpa judul; LABEL "EDM" berwadah S2 (mustahil) | `nonprop.ts` `kolomTotal` + `F/nonprop.test.ts`; `DetailNonProp.tsx` |

### 6.3 Sah tidak dibangun (syarat dicek ke XML)

- **Elemen mati (`1=2` / `NEVER`):** `.BizName`, `.DueTo` admin, label pengisi "s"/"spc"/"space",
  dan LABEL "Section for …". Juga total berlabel admin S39 dan `SpreadingRiskList` S6, tombol
  `TreatyInNonSetTotal`, Stage/List view S14, "Create Opportunity" showMenu, dan
  `SFAPortal_OpportunitiesList_Header`.
- **Keputusan tertulis:**
  - K7: Survey Report (kedua layar).
  - K8 butir 4: "Choose Business R", `DetailPolicyTreatyOutNonProportional`.
  - K9: grid Breakdown Spreading S34 / S100.
  - K2: `CekLimitTreatyAcc_Act` dan varian Submit LetterNo.
  - K5/K12: identitas operator di syarat Submit, LABEL "NON EDM".
- **`Protection_Act`** varian kelas Data-PolicyTreatyIn tidak ada di korpus.
- **Pemilih SOB, nol efek** (bab 4 #6, #19, #20): autocomplete `SearchSOB.CARI1`, perluasan
  simpul pohon (`btnQuotation` OfferTreatyIn ditulis nol rule), dan aktivitas Choose.

### 6.4 Perlu cek (bukan selisih, belum dapat diputuskan dari korpus)

- **Default sel.** Apakah Pega menerapkan `pyDefaultValue` sel saat render: `.QuotationData.IsSurveyReport`
  "No", `.TypeTax` "Inclusive". Lihat W5.
- **Tombol bawaan FlowAction.** `pyShowFAButtons` / `pySaveButtonVisibility=Always`; harness
  Perform tidak ada di korpus.
- **Nilai kosong.** Apakah `FacultativeShare` kosong = 0 di ekspresi `pyCondition`, padahal
  `SpreadingRiskList` sendiri memeriksa `= 0 || = ''`.
- **Parameter RD dropdown.** `ID=.TreatyType` pada dropdown RD: apakah menyempitkan daftar.
- **Teks opsi.** Radio `.DueTo`, Approval, serta pilihan IsSurveyReport/StatementType/TypeTax/ClaimType:
  pilihannya ada di rule Property yang tidak terekspor, sehingga dirender sebagai teks bebas.

### 6.5 Keadaan sesudah R1/R2/R3 (`97fcb333`)

**Pemilih SOB (R2, F4).**

- Hasil `SearchHierarkiSourceBizAgent_PostDT` kini dipegang layar dan dikirim pada Save/Submit.
  `POST /pilih-sumber-bisnis` menjadi pencarian tanpa simpan. Ini sesuai XML: FA
  `AgentSourceBizDetails` `pyActionTransformRule` hanya mengubah clipboard, dan tombol Choose
  `SourceHierarki` (`runActivity` → `opener.location.reload` → `window.close`) tidak berisi
  `Obj-Save`.
- Satu uji paket itu (`hitungNetPremi`) disesuaikan ke aksi sel admin `CountOGPONP` (7.4).
  Tidak ada selisih layar baru.

**Kolom `T_GENERAL_POLIS.EDM_TYPE` (R1, F3).**

- Kolom ini dibaca syarat kedua `InputPolicyTreatyInPre_Act` langkah 10
  (`[.PolicyTreatyIn.EDMType=="3"] T->3`). Cocok dengan `PerluCekDaftarXOL`.
- Fungsi itu juga menambah syarat ketiga, `QuotationData.ProportionalType != "Proportional"`,
  yang **tidak ada di XML**. Penyimpangan ini ditandai `[penyimpangan sadar]` (penjaga AC 33),
  tetapi tanpa keputusan WO yang dikutip. Lihat catatan 7.4.
- Bukan medan layar.

**Pemuat dokumen lama (R1, F3).**

- Pemetaan `SuggestList` lama (`models/usulanlama.go`) sama dengan `SaveViewSuggest` langkah
  2.1.2: CARI3 "Policy", CARI4 `.OperatorName`, CARI5 `.Date`, CARI8 2, CARI9 Accept/Reject,
  dan CARI10 `@substring(.Suggest,0,3990)`.
- Uji petik `pemakai.py` atas enam medan "dibuang" di `medan_abaikan_lama.json` bagian `pola`
  cocok dengan bukti yang ditulis R1: `QuotationData.BusinessType`, `.SobName`,
  `.CedingCoName`, `.TeamGroup`, `.BranchCode`, `IsOJKNopolis`.
  - Pembaca terjangkau `.BusinessType` / `.SobName` hanya halaman `pyWorkPage.Quotation`, atau
    langkah berlabel `//` (`CheckDataMkt` 6–7).
  - `.CedingCoName` yang dibaca `TreatyRealizationCheckDuplicate` adalah
    `.PolicyTreatyIn.CedingCoName`, yang berkolom.
  - `.TeamGroup` di `BrowseMarketingOfficer_RD` adalah kolom tabel MO, bukan medan polis.

---

## 7 · Temuan WAJIB-BANGUN dan WAJIB-PERBAIKI

**WAJIB-BANGUN: nol** (bab 4).

> ⭐ **Keadaan sesudah putaran 3 (konsolidasi P3K, 04-10-2026): W1–W6 dan 7.4 seluruhnya diperbaiki** di
> paket R5 (W1) dan R6 (W2–W6, 7.4), digabung ke integrasi `312f3b1c` dan `724ef8eb`. Penanda per butir di bawah;
> isi temuan tidak diubah.

### W1 · WAJIB-PERBAIKI (besar) — popup pilih bisnis membaca sumber yang salah

> ✅ **Diperbaiki putaran 3** — paket **R5**: `258bf1af` (server: `repository.DaftarBisnis` RD `BrowseTreatyJoinEDM`, filter `PROPORTIONTYPE`, `POST /kasus/{id}/bisnis`; uji `handlers/daftarbisnis_test.go`, db `TestDaftarBisnisRDBrowseTreatyJoinEDM`), `f06b1790` (layar: judul dan 25 kolom XML, saring/urut per kolom, 50/halaman), `312f3b1c` (RALAT spec §5.1, tiket 01/11, status.json, tombol.json).

- **Bukti XML.** `Section/BusinessAndSOBList`:
  - Grid lama S11 (`pgRepPgSubSectionBusinessAndSOBListBB`, kelas `TREATYINDETAIL`, RD
    `BrowseTreatyInDetail`) ber-`pyContainerVisibleWhen=1=2`. Memo rule: "Hidden the old one,
    now use treatyindetail join edm".
  - Grid aktif S16/S17 `pyGridProps/pyRDName=BrowseTreatyJoinEDM`, kelas
    `ASM-FW-GISFW-Int-TREATYINDETAILJOINEDM`, `pyRDParams` `PROPORTIONALTYPE =
    .QuotationData.ProportionalType`. Ini filter H RD `.PROPORTIONTYPE = Param.PROPORTIONALTYPE`,
    yang diabaikan bila kosong.
  - Grid ber-`pyGridFiltering`/`pyGridSorting=true`, `pyPageSize=50`, tanpa kotak cari.
  - Judul 25 kolom: Treaty Offer ID, Contract Name, Class of Business, Source Of Business, Insured
    Name (`.CEDING`), Proportion Type, Treaty Type, Treaty Group, Treaty Year, Currency, Limit,
    Currency, Retention, Currency, EPI, Layer, (kosong), Part of, (kosong), Currency, MDP,
    Currency, Net Premium, Currency, Share RNM Value.
  - Judul jendela `pyWindowName` "Business And SOB List".
- **Bukti kode:**
  - `repository/acuan.go` `DaftarDetailKontrak` membaca tabel `TREATYINDETAIL` (`tabelDetail`)
    tanpa filter PROPORTIONTYPE. Komentarnya "Popup tidak mengirim parameter" keliru.
  - `services/layanan.go:347` `DaftarBisnis` tidak tahu konteks kasus (`GET /bisnis` tanpa id).
  - `PilihBisnis.tsx` memasang kotak "TREATYID" + Filter, dan judul kolom memakai nama kolom view.
- **Akibat.** Kontrak EDM tidak pernah tampil di popup. Bila sudah ada pilihan sebelumnya, saringan
  jenis proporsi tidak berlaku.
- **Usulan.**
  - Baca view `TREATYINDETAILJOINEDM` dengan filter `PROPORTIONTYPE` dari
    `PolicyTreatyIn.QuotationData.ProportionalType` kasus. Rute membawa id kasus.
  - Ganti judul kolom dan judul jendela dengan teks XML.
  - Saring dan urut per kolom di klien atas ≤ 500 baris, paging 50, tanpa tombol Filter.
  - Uji `R/` ber-tag db.
  - RALAT INVENTARIS 174/179/5124/5125, tombol.json, dan komentar `acuan.go`.

### W2 · WAJIB-PERBAIKI (besar, butuh RALAT spec) — grid spreading NonProp tersunting di layar atasan

> ✅ **Diperbaiki putaran 3** — paket **R6**: `172ddb44` (server: cabang atasan menerima `SpreadingRiskList` NonProp bila `FacultativeShare` 0/'', `CountSpreading` di aksi atasan), `5b8c929d` (layar), `724ef8eb` (RALAT spec AC 52, US 11, §5.11; AC 52 tetap ✅).

- **Bukti XML:**
  - `DetailDeptHeadTreatyIn_UW` menyertakan `DetailPoliciesNonProportional` dengan sel
    SUB_SECTION `pyReadOnly=false`, `pyEditOptions=Auto`. Begitu pula
    `DetailPolicyTreatyInNonProportional` → `SpreadingRiskList` (`pyReadOnly=false`, `Auto`).
  - Di dalam `SpreadingRiskList`, Add/Delete tampil bila `FacultativeShare` 0/''.
    `.TreatyType`/`.SharePercentage`/`.ClaimPercentage` hanya `pyReadOnlyCondition
    pyWorkPage.TreatyIn.FacultativeShare >0`, dengan action set refresh `CountSpreading_Act`.
  - Grid Proporsional atasan, sebaliknya, `pyReadOnly=true` tanpa syarat.
- **Bukti kode:**
  - `LayarKasus.tsx` `sunting={admin && boleh}`.
  - `models/layar.go` `GabungMasukanLayar`: cabang atasan tidak menerima daftar.
  - `services/aksiposisi.go`: atasan hanya `SetDueTo`.
  - Spec AC 52 menyatakan atasan hanya mengisi ListSuggest.
- **Usulan.** RALAT AC 52 dan INVENTARIS:5153, atau minta keputusan WO tertulis. Bila XML diikuti:
  - Buka sunting NonProp bagi atasan dengan syarat yang sama.
  - Terima `SpreadingRiskList` di cabang atasan bila `PolisNonPropBaru` dan FacultativeShare 0/''.
  - Tambahkan `CountSpreading` ke `aksiAtasan` dengan syarat yang sama.

### W3 · WAJIB-PERBAIKI (sedang) — server menerima medan terkunci hasil tombol Enable/Disable

> ✅ **Diperbaiki putaran 3** — paket **R6** `eaa49344`: diterima hanya bila TreatyType XOL dan sama dengan `TreatyEnableDisableInput` atas halaman server; lainnya 422 (`TestEnableDisableHanyaDariTombolXOL`).

- **Bukti XML.** Sel `.QuotationData.ProportionalType` `pyReadOnly=true`
  (`pyEditOptions=Read-only`). `IsNewPolicyNonProp` bukan sel layar. Keduanya hanya diubah DT
  `TreatyEnableDisableInput` (tombol bertampil `.TreatyType='XOL'`, refresh tanpa simpan).
- **Bukti kode.** `models/layar.go` `medanAdmin` menerima `IsNewPolicyNonProp` dan
  `QuotationData.ProportionalType` dari layar apa adanya.
- **Usulan (pola F4).** Terima hanya bila `TreatyType=='XOL'` dan nilainya sama dengan hasil
  `TreatyEnableDisableInput` atas halaman tersimpan. Hasil DT itu selalu "0", lihat
  `TestEnableDisableSelaluNol`.

### W4 · WAJIB-PERBAIKI (sedang) — server menerima medan di wadah/sel yang TERSEMBUNYI

> ✅ **Diperbaiki putaran 3** — paket **R6** `eaa49344`: daftar izin admin + syarat tampil sel/wadah (S19, S7, S14, `IDCurrency`, `IsSurveyReport`); FlagPPH berubah → `RemoveTypeTax` (`TestMedanAdminTersembunyiTidakDiterima`, `TestKirimanLayarTidakMenimpaUangMasterNonProp`).

- **Bukti XML:**
  - S19 `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1` membungkus seluruh medan uang
    dan `.Installment`.
  - `.IDCurrency` `pyCondition .IsNewPolicyNonProp != 1`.
  - S7 `.ClaimType != 'XOL Retro'` membungkus FlagPPH/TypeTax.
  - S14 `.QuotationData.ProportionalType = 'Proportional'` membungkus Quartal/YearOfQuartal.
  - `IsSurveyReport` `pyCondition pyWorkPage.Quotation.ProportionalType != 'NonProportional'`.
- **Bukti kode.** `models/layar.go` `GabungMasukanLayar` admin menyalin `medanAdmin` tanpa
  predikat tampil. Bagi polis NonProp baru, `PremiOgp`/`Deduction1`/`Deduction2` adalah milik
  server (`nonprop_detail.go`), dan `services.turunkan` tidak menghitungnya ulang. Akibatnya
  kiriman layar dapat menimpa nilai hasil master XOL.
- **Usulan.** Saring medan admin dengan predikat tampil yang sudah ada (`wadahUangAdmin`,
  `bukanXOLRetro`, `proporsionalQD`, `bukanNonProp`, `bukanNonPropBaru`), lalu tambahkan uji HTTP.

### W5 · WAJIB-PERBAIKI (sedang, butuh keputusan) — kolom hanya-baca di daftar dari layar dan nilai bawaan sel

> ✅ **Diperbaiki putaran 3** — paket **R6** `174d5557` (kolom hanya-baca spreading dihitung server, `CountSpreading_Act` 4.1) dan `eaa49344` (`ListInstallment` tidak diterima dari layar; `pyDefaultValue` `IsSurveyReport` "No", `TypeTax` "Inclusive"). Dua `[penyesuaian sadar]` tercatat di HASIL bab 5b butir 32.

- **Kolom hanya-baca diterima dari layar.**
  - Grid spreading `.PremiumSpreaded`/`.ClaimSpreaded` `Read-only`, dan grid angsuran seluruhnya
    `readOnly`. Meski begitu, `DaftarDariLayar` menerima kedua daftar utuh. `turunkan` hanya
    menjumlah total (`HitungTotalSpreading`) dan tidak menghitung ulang `CountSpreading_Act`
    langkah 4.1 per baris.
  - Usulan: hitung ulang 4.1 di server dan simpan `ListInstallment` hanya dari
    `FillPaymentInstallment` server, atau validasi. Perlu keputusan, sejalan dengan F4.
- **Nilai bawaan sel.** `.QuotationData.IsSurveyReport` `pyDefaultValue="No"` dan `.TypeTax`
  `pyDefaultValue="Inclusive"` tidak diterapkan. Perilaku runtime Pega perlu cek.

### W6 · WAJIB-PERBAIKI (kecil–sedang, butuh keputusan WO, bab 0 butir 7) — unsur layar tanpa dasar Section XML

> ✅ **Diperbaiki putaran 3** — paket **R6** `5b8c929d` (layar) dan `724ef8eb` (RALAT tiket 11, tombol.json): spanduk NBStatus, kolom portal Position/No Polis, judul panel buatan dibuang; ikon pengosong saring, format grid master NonProp, kolom kosong pertama grid total dibangun; tautan `.Name` → "Offer No" `[penyimpangan sadar]` (HASIL bab 5b butir 33); panel History, tombol Kembali, Cancel modal dipertahankan dengan dasar tertulis.

Unsur yang ada di kode tetapi tidak di Section XML:

- tombol "Kembali" (`LayarKasus.tsx`);
- spanduk NBStatus layar kasus (XML hanya di grid portal);
- panel "History" + `GET /kasus/{id}/riwayat` (dasarnya AC 72 spec, bukan Section);
- judul panel buatan "General", "Premium & Claim", "Spreading Risk", "Suggest", dan judul panel
  NonProp. XML: wadah NOHEADER, sedangkan S24/S25 berlabel "OGP"/"ONP";
- kolom portal "Position" dan "No Polis" (grid `GetListOpportunity` hanya 6 kolom);
- tombol "Cancel" modal nomor polis (bawaan `Modal` inti);
- teks daftar kosong portal;
- kolom "Name" (`.Name`, ditulis nol rule): tautan dipindah ke kolom "Offer No". Ini perlu
  RALAT eksplisit.

Tiga hal kecil yang juga belum dibangun:

- ikon pengosong saring (`setValue ""` + `postValue`, tanpa refresh);
- format angka per sel grid master NonProp (beberapa tanpa `pyDecimalPlaces`);
- kolom kosong pertama grid TotalLimitIOONP/TotalFacShareRnmNP.

### 7.4 Catatan untuk status.json (bukan WAJIB-PERBAIKI kode)

> ✅ **Diperbaiki putaran 3** — paket **R6** `eaa49344` (`FetchTreatyGroupOldID` langkah 3 menahan Submit Dept Head, `TestGrupTreatyKosongMenahanSubmitDeptHead`), `97e4ddfc` (langkah 1–3 `//` `CountSpreading` tidak lagi diport), `724ef8eb` (`PerluCekDaftarXOL` berdasar spec-penyimpanan AC 33, tiket 19). Kelima aktivitas tak terpicu **ditetapkan kategori (a)** di status.json oleh konsolidasi P3K `7052dcb4` (HASIL bab 6: dibangun 103, (a) 36, (b) 35, (c) 2).

Lima activity berstatus "dibangun" ternyata **tidak pernah terpicu** dari layar. Fungsinya ada
di `models` dan kini berujian, tetapi aksi hitungnya ditolak server:

- `CountRiCommOgp_act`, `CountRiCommOnp_act`, `CountOverridingCommOgp_Act`,
  `CountOverridingCommOnp_Act`. Satu-satunya pemicu: sel `.ResultOgp1`/`.ResultOnp1`/
  `.ResultOgp2`/`.ResultOnp2` `DetailDeptHeadTreatyIn_UW`, semuanya `pyReadOnly=true`.
- `CountPctInstallment_Act`. Satu-satunya pemicu: sel `.Premium` grid S45 `readOnly`, di kedua
  layar.
- `SetValidateInstallment_Act` tetap terjangkau sebagai `CountOGPONP_Act` langkah 10.
- `CountNetPremi_act` terjangkau sebagai sub-panggilan.

**Usulan:** status.json menambahkan catatan "(a) tidak terpicu: sel pemicu ber-pyReadOnly /
grid readOnly" pada kelima baris itu, atau memindahkannya ke (a). Keputusan ada pada orkestrator.

Tiga catatan lain:

- **`FetchTreatyGroupOldID`** langkah 3 (Page-Set-Messages "Cannot fetch Treaty Group ID,
  Contact IT" bila `TreatyGroupID==""`) tidak dibangun dan tidak disebut status.json. Perlu
  keputusan: pesan itu menahan Submit Dept Head bila kontrak tanpa grup.
- **`CountSpreading`** (`models/angsuran.go`) memport langkah 1–3 yang berlabel `//`. Tidak ada
  beda perilaku, karena langkah 4 menimpa dengan rumus yang sama, tetapi komentar kodenya perlu
  RALAT.
- **Uji paket F4** `H/sumberbisnis_test.go` `hitungNetPremi` disesuaikan memakai aksi sel admin
  `CountOGPONP` di commit merge, karena `CountNetPremi` tidak lagi diterima.
- **`PerluCekDaftarXOL`** (`models/nonprop_detail.go`) menambah syarat
  `QuotationData.ProportionalType != "Proportional"` di luar XML `InputPolicyTreatyInPre_Act`
  langkah 10. Penyimpangan ini ditandai sadar (penjaga AC 33), tetapi tanpa keputusan WO.
  Perlu dicatat di HASIL bab 5, atau minta konfirmasi WO.

---

## 8 · RALAT dokumen yang diperlukan (untuk orkestrator)

Audit ini tidak menyunting HASIL, PERMINTAAN, status.json, maupun INVENTARIS.

> ✅ **Seluruh RALAT di tabel ini sudah diterapkan putaran 3**: status.json `TreatyInInputVis`, `GetCurrentDate`, lima
> baris 7.4, dan INVENTARIS dibangkitkan ulang (P3K `7052dcb4`); HASIL bab 6 disamakan (P3K); tombol.json, INVENTARIS
> popup, `labels.ts` (R5 `312f3b1c`, `f06b1790`); spec AC 52 (R6 `724ef8eb` — W2 dibangun, bukan menunggu keputusan).

| Dokumen | Bunyi lama | Bunyi baru | Bukti |
| --- | --- | --- | --- |
| status.json `Activity/TreatyInInputVis` | "langkah 3 `Page-Remove TreatyIn`" (seolah berjalan) | langkah 1 bertransisi `1==1` → 6 (keluar), sehingga langkah 3 tidak berjalan. Kesimpulan nol efek tetap. | `Activity/TreatyInInputVis.xml` langkah 1 `pyStepsTransParams` |
| status.json `RDBList/GetCurrentDate` | "`TreatyRealizationCheckXOLList` langkah 4 dengan parameter `viewstate` dan `ID` saja" | langkah 4 `pyPassCurrentParameterPage=true`, sehingga isian `viewstate=1, ID=3000004` diabaikan dan halaman parameter hanya membawa `Param.ID` = NoOffer. `revisionstate` tetap kosong. | `Activity/TreatyRealizationCheckXOLList.xml` langkah 4 |
| HASIL bab 6 | (a) 29 / (b) 36 | samakan dengan status.json (bab 2) | — |
| status.json 5 baris (7.4) | "dibangun" | tambahkan "tidak terpicu" | 7.4 |
| `docs/alat/tombol.json` | Stage/List view "tampil selalu"; popup bisnis "dua grid", RD `BrowseTreatyInDetail` | wadah S14 `1=2`; satu grid aktif `BrowseTreatyJoinEDM` (W1) | `SFAPortalOpportunitiesHeader.xml`, `BusinessAndSOBList.xml` |
| INVENTARIS-XML 174/179/5124/5125/5153 | popup bisnis `BrowseTreatyInDetail`; atasan tanpa sunting NonProp | W1, W2 | idem |
| spec AC 52 | "medan lain hanya-baca" (atasan) | menunggu keputusan W2 | W2 |
| `labels.ts` (komentar `KOLOM_BISNIS`) | "nama kolom view, apa adanya" | W1 | W1 |

---

## 9 · Commit dan verifikasi

| Commit | Isi |
| --- | --- |
| `54b33597` | uji berharapan XML untuk 17 rule dibangun tanpa uji (bab 3.4) |
| `d9944fa0` | P1 — hitung hanya aksi sel terbuka di layar posisi (`services/aksiposisi.go`) |
| `6a898d57` | P2 (server) — grid angsuran hanya-baca: dua aksi keluar dari daftar admin |
| `3b85a5d0` | P2–P9 (layar) — label, modal, portal, caption, grid angsuran hanya-baca, dropdown RD |
| `82c1e4a1` | P10 — grid total NonProp, LABEL EDM |
| `ecb8d771` | P5 — paginasi ListSuggest |
| `3ae24c8a` | gabung implementasi `e4d847c0` (F4, K18, F2, F5, F8) + penyesuaian helper uji F4 |
| `6cca634b` | dokumen ini (draf) |
| `3eeeca5a` | gabung implementasi `97fcb333` (F3 dokumen lama, EDM_TYPE) |

Verifikasi dijalankan sesudah penggabungan `97fcb333`:

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` | bersih |
| `go vet -tags db ./modul/nbtreatyin/...` | bersih |
| `go test ./...` | `nbtreatyin` hijau. Merah hanya baseline claimlife (`models`, `repository`, `services`) dan `inti/backend/penjaga` `TestNolAlamatLayananDiKode`, yang merah karena worktree tanpa `.env`. |
| `npm run typecheck` | 0 galat |
| `npx vitest run modul/nbtreatyin frontend/daftar` | 17 berkas, 111 uji lulus |
