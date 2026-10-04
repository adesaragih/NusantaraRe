# 11: Layar realisasi — medan wajib, medan terkunci, dan bagian yang tidak dibangun

**Status:** selesai — tertahan hanya pihak luar: AC 45 `ProductionDate` lewat tempat berperan (K12); `DateofSurvey` ⛔ (b) K7; empat wadah 104 medan (spec §9.2 butir 17, Product & Underwriting) *(putaran 2, konsolidasi P10 04-10-2026 — rincian `docs/HASIL-IMPLEMENTASI.md` bab 9; semula: sebagian, implementasi 2026-10-03; awalnya ready-for-agent)*
**Blocked by:** 02
**Menutup:** AC 45 · 46 · 49 · 50 · 51 · 53 · 54 · 55 · 56 · 77 *(10 AC)* — US 27 · 28 · 29 · 31 · 32

## Hasil & nilai pengguna

Hari ini layar realisasi menuntut sejumlah medan diisi, mengunci sebagian lainnya, dan
menyembunyikan bagian yang tidak berlaku. ⚠️ `[terverifikasi]` **Delapan puluh** bagian layar
diberi syarat tampil yang **tidak mungkin pernah benar** — dan **tidak satu pun** dari yang 80 itu
menyembunyikan sebuah medan; seluruhnya label dan hiasan.

Sesudah tiket ini, pengguna melihat layar yang **menandai medan wajibnya dengan jelas**, mengunci
yang tidak boleh ia ubah, dan ⭐ **tidak memuat 80 bagian mati** yang selama ini ada tanpa pernah
tampil.

## Area codebase

- Lapisan handler: validasi medan wajib
- Antarmuka: penandaan medan wajib dan terkunci
- Antarmuka: daftar pilihan mata uang

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Medan wajib | **27** medan berbeda di **6** layar; sebarannya di `spec.md` §5.11 |
| Medan terkunci | **38**, seluruhnya lewat syarat hanya-baca yang **selalu benar**; setiap kunci mengenai **tepat satu medan** |
| Bagian mati | **91**; ⭐ **80** menempel pada sel tunggal, **11** pada wadah yang memuat medan |
| Pengecualian mata uang | `ReportDefinition\BrowseCurrency_RD.xml` · `BrowseCurrencyTreatyIn_RD.xml` |
| Pembersihan pesan galat | **12** tempat yang menghapus **sebelum** pesan dipasang |

## ADR terkait

- **ADR-0002** — RBAC memakai peran yang sudah ada

## Acceptance criteria

- [ ] 🟡 **AC 45** — layar menuntut medan wajibnya; **27** medan berbeda di **6** layar
- [x] **AC 46** — layar jenjang ketiga **tidak** mewajibkan enam medan yang wajib di layar admin
- [x] **AC 49** — **38** medan terkunci permanen
- [x] **AC 50** — **36** dari 38 di layar jenjang ketiga; **2** di layar biasa
- [x] **AC 51** — setiap kunci mengenai **tepat satu medan**, ⛔ bukan bagian atau tab
- [x] **AC 53** — **80** bagian mati **tidak dibangun**
- [x] **AC 54** — daftar pilihan mata uang **tidak memuat** kode yang dikecualikan
- [x] **AC 55** — kode jenis kontrak ditampilkan **apa adanya** bila keterangannya belum tersedia
- [x] **AC 56** — nilai kode di luar daftar yang dikenal **tetap diterima dan disimpan**
- [x] **AC 77** — pembersihan pesan galat di awal diterima apa adanya — **12** tempat

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **17** | **empat wadah** berisi **104** medan — masih dipakai atau ditinggalkan | ⚠️ menahan **keempatnya saja**, bukan tiketnya |
| **14** | apakah ada mata uang mati lain yang tidak disembunyikan | tidak menahan |

## Perintah verifikasi

1. Kosongkan satu medan wajib, simpan — ⭐ **ditolak**, dengan pesan yang menyebut medannya.
2. Coba sunting medan terkunci — ⭐ **tidak bisa**, dan tampil sebagai terkunci.
3. Cari kode mata uang yang dikecualikan di daftar pilihan — ⭐ **tidak ada**.
4. Kirim kode jenis kontrak yang tidak dikenal — ⭐ **tetap tersimpan**, ditampilkan apa adanya.

## Catatan

⭐ **Delapan puluh bagian mati dibuang tanpa ditanyakan kepada siapa pun** `[keputusan work owner]`
P44 — `[terverifikasi]` sebabnya struktural: membangunnya berbiaya nol, membuangnya juga.
⛔ **Empat wadah berisi 104 medan** menunggu Product & Underwriting dan **tidak dibangun** sebelum
dijawab.

## Hasil implementasi 2026-10-03

- Medan wajib per layar dari `pyRequired` (`models.DaftarMedanWajib`): 25 medan di dua layar
  realisasi + `ProductionDate` (tempat berperan, tiket 05) + `DateofSurvey` (layar survei historis,
  tidak dibangun) = 27 (AC 45). Ditegakkan pada Submit DAN Save (AC 48).
- Pilihan radio/dropdown bersumber "associated values" rule Property — **tidak ada di korpus**
  (TypeTax, IsSurveyReport, StatementType, ClaimType, ClaimPaymentType, DueTo) ⇒ isian teks apa adanya
  (AC 56), tidak dikarang. Approval memakai nilai `1`/`0` dengan teks dari `SaveViewSuggest`
  (Accept/Reject).
- Tidak dibangun: tombol/layar Survey Report (penyimpanan survei tidak dirancang, `[terbuka]`),
  pemilih SOB (hanya untuk XOL Retro), Choose Business R (treaty keluar, JSON).

> ⛔⛔ **RALAT.** 2026-10-03 (putaran 2, paket P3) — bunyi lama dikutip utuh lalu ditarik untuk butir
> pemilih SOB:
> > *"Tidak dibangun: tombol/layar Survey Report (penyimpanan survei tidak dirancang, `[terbuka]`),
> > pemilih SOB (hanya untuk XOL Retro), Choose Business R (treaty keluar, JSON)."*
> > *(status.json lama: "pemilih hierarki sumber bisnis/ceding hanya tampil bila ClaimType 'XOL
> > Retro' (jalur retro, P29) — tidak dibangun")*
>
> ⭐ **Pemilih SOB DIBANGUN.** "Hanya tampil bila XOL Retro" bukan bukti tak terjangkau, dan P29
> (data kontrak dari view, nol JSON) tidak menyentuh jalur ini: seluruh datanya relasional (tabel
> `AGENT`). Bukti XML:
> - `Section/DetailPolicyTreatyIn` (dan salinannya di `GeneralPolicyTreatyIn`, layar admin flow
>   action `InboxPolicyTreatyIn`): tombol `Select Source Of Business`, pyVisible
>   `.ClaimType = 'XOL Retro'`, click → showHarness `SOB` (pyUsingPage `pyWorkPage.Quotation`,
>   pySubmitData Yes) + aktivitas `InputQuotation_PreAct(Acton=SOB)` → langkah 1 `btnSOB_DT`
>   (`Quotation.btnQuotation = "SOB"`). Layar atasan tidak memuat tombol ini.
> - `Harness/SOB` → `Section/SourceHierarki`: TreeGrid `TempBusinessSource.pxResults`,
>   pyDeferLoadActivity `AgentSourceBizTreatyIn_Act` (langkah 3 `pxShowReport`
>   `BrowseAgentHierarkiList_RD`), pyRowEditing masterDetail → `FlowAction/AgentSourceBizDetails`
>   pra-proses `SearchHierarkiSourceBizAgent_PostDT`: `Quotation.SourceOfBusiness/SobName/
>   SobLeader0/SobLeader1 = @if(.ChildCount > 0, "", ...)`.
> - Akibat pada data kasus: `Quotation.SourceOfBusiness` (medan diagram `T_POLIS_QUOTATION`) dibaca
>   `SetPPNPPH` langkah 1-3 (STS_PKP agen) → syarat PPH/PPN langkah 4.
>
> Dibangun: `models/sumberbisnis.go`, `repository/agen.go`, `services/sumberbisnis.go`, rute
> `GET /sumber-bisnis` dan `POST /kasus/{id}/pilih-sumber-bisnis`, `components/PilihSumberBisnis.tsx`.
> Tetap tidak dibangun, alasan (a) tak berpengaruh/tak terjangkau di NB: autocomplete `Search`
> (`SearchSOB.CARI1`, dibaca nol rule); perluasan simpul (`Param.Leader` diisi hanya bila
> `pyWorkPage.OfferTreatyIn.QuotationData.btnQuotation=="SOB"` — properti yang ditulis nol rule ⇒
> setiap simpul = daftar akar yang sama); aktivitas tombol `Choose`
> (`SearchHierarkiSourceBizAgentTreatyIn_Act` menulis `OfferTreatyIn.QuotationData.*`, dibaca nol
> rule NB); cabang Ceding (`btnCedingCO_DT`/`BrowseCedingCo_RD`, hanya lewat `Acton=Ceding` yang
> tak pernah dikirim).

> ⛔⛔ **RALAT** 2026-10-04 (putaran 3, paket R2 — F4, keputusan WO 04-10-2026: **ikuti XML**) — dua
> bunyi lama RALAT paket P3 di atas dikutip apa adanya:
> > *"Dibangun: `models/sumberbisnis.go`, `repository/agen.go`, `services/sumberbisnis.go`, rute
> > `GET /sumber-bisnis` dan `POST /kasus/{id}/pilih-sumber-bisnis`, `components/PilihSumberBisnis.tsx`."*
> > *"Akibat pada data kasus: `Quotation.SourceOfBusiness` (medan diagram `T_POLIS_QUOTATION`) dibaca
> > `SetPPNPPH` langkah 1-3 (STS_PKP agen) → syarat PPH/PPN langkah 4."*
> >
> > (kode P3, `services/sumberbisnis.go`: *"`[penyesuaian sadar]` Pega memegang hasil PostDT di
> > clipboard sampai Save/Submit layar utama ... jadi hasilnya DISIMPAN saat klik"*)
>
> ⭐ **Bunyi baru.** Klik baris **tidak menyimpan**. `POST /kasus/{id}/pilih-sumber-bisnis` menjadi
> **pencarian tanpa simpan**: menjalankan `btnSOB_DT` + `SearchHierarkiSourceBizAgent_PostDT` atas
> baris RD yang dibaca ulang, lalu menjawab keempat nilai `Quotation.SourceOfBusiness/SobName/
> SobLeader0/SobLeader1`. Layar **memegang** nilai itu (`pegangSumberBisnis`) dan mengirimnya
> bersama **Save / Submit** (dan refresh). Server (`services.terimaSumberBisnis`, dari `kerjakan`
> layar admin) menerimanya **hanya** bila:
> 1. `Quotation.SourceOfBusiness` kiriman berbeda dari nilai server (selain itu tak ada pilihan baru;
>    tiga medan lain tanpa kolom dan tanpa pembaca NB tidak diambil dari layar);
> 2. ClaimType isian layar **`XOL Retro`** — selain itu tombolnya tak tampil, medan **terkunci**:
>    kiriman diabaikan, nilai tersimpan dipakai (pola AC 49–51);
> 3. sama persis dengan hasil PostDT salah satu baris RD `BrowseAgentHierarkiList_RD` yang
>    **dijalankan ulang** di server (filter efektif NB `LEADER0 IS NULL AND STATUSACTIVE IS NULL`);
>    kosong-semua sah hanya bila RD memuat simpul beranak (`ChildCount > 0`). Tidak cocok → **422**
>    *"Source Of Business "…" tidak cocok dengan hasil pencarian hierarki sumber bisnis
>    (BrowseAgentHierarkiList_RD) - pilih ulang lewat tombol Select Source Of Business"*.
>
> Yang tersimpan hanya medan berkolom: `SOURCE_OF_BUSINESS` (`T_POLIS_QUOTATION`); `SobName`,
> `SobLeader0`, `SobLeader1` tetap tanpa kolom (katalog: dibuang). Ralat kecil bunyi P3 kedua:
> `SetPPNPPH` langkah 1 membaca **`pyWorkPage.PolicyTreatyIn.QuotationData.SourceOfBusiness`**
> (`Param.ID`), bukan `Quotation.SourceOfBusiness`.
>
> **Bukti XML** (dibaca ulang 2026-10-04):
> - `DataTransform/SearchHierarkiSourceBizAgent_PostDT.xml`: langkah 1 `WHEN
>   pyWorkPage.Quotation.btnQuotation=="SOB"` → 1.1 `pyWorkPage.Quotation.SourceOfBusiness =
>   @if(.ChildCount > 0, "", .ID)`, 1.2 `.SobName = @if(.ChildCount > 0, "", .ClientName)`; langkah 4–5
>   `SobLeader0/SobLeader1 = @if(.ChildCount > 0, "", .Leader0/.Leader1)` — **nol Obj-Save**.
>   `FlowAction/AgentSourceBizDetails.xml` hanya `pyPreProcessingTransformRule` (dan
>   `pyActionTransformRule`) = PostDT itu.
> - `Section/SourceHierarki.xml` tombol `Choose` (pyVisible `.ChildCount = 0`): `runActivity
>   SearchHierarkiSourceBizAgentTreatyIn_Act` → `runScript opener.location.reload` → `runScript
>   window.close`; aktivitas itu (langkah 1–4) hanya menulis `pyWorkPage.OfferTreatyIn.QuotationData.*`
>   dan `Local.*` — **nol Obj-Save**. `Harness/SOB.xml` tak memuat tombol lain.
> - `Section/DetailPolicyTreatyIn.xml` tombol `Select Source Of Business`: satu-satunya perilaku
>   `click → showHarness SOB` (`pySubmitData=Yes`, `InputQuotation_PreAct`) — **tanpa `refresh`**
>   (bandingkan `Choose Business`: `showHarness` lalu `refresh`). Jadi XML **tidak menjalankan refresh
>   berhitung** sesudah klik/Choose; layar baru pun tidak — PPN/PPH berubah pada refresh berikutnya
>   yang dipicu pengguna, atau pada Save/Submit (turunan dihitung ulang server, AC 49).
> - `Activity/SetPPNPPH.xml` langkah 1 `Param.ID = pyWorkPage.PolicyTreatyIn.QuotationData.
>   SourceOfBusiness`; penyalinnya hanya `InputPolicyTreatyIn_preDT` langkah 14
>   (`PolicyTreatyIn.QuotationData = pyWorkPage.Quotation`), preACT 14.9, `GeneratePolicyNoTreaty_Act`
>   langkah 10. ⚠️ `[penyesuaian sadar]`: apakah `opener.location.reload` membuat mesin Pega
>   menjalankan ulang pra-proses itu **tidak ada di korpus**. Di sini pra-proses dijalankan ulang di
>   setiap permintaan atas data tersimpan, maka pilihan yang dipegang layar disalin ke `QuotationData`
>   sesudah diterima (`models.SalinKeQuotationData`) supaya `SetPPNPPH` memakai nilai layar saat
>   refresh/Save/Submit — bukan nilai tersimpan — dan supaya pilihan tidak tertimpa salinan lama
>   (`NilaiQuotation`, ID-23).
>
> Uji: `handlers/sumberbisnis_test.go` (`TestKlikBarisSumberBisnisTidakMenyimpan`,
> `TestKlikSimpulBeranakMengosongkan`, `TestKlikSumberBisnisDitolak`,
> `TestSaveMenyimpanSumberBisnisYangCocok`, `TestNilaiSumberBisnisPalsuDitolak`,
> `TestSumberBisnisTerkunciBilaBukanXOLRetro`, `TestSaveSimpulBeranakMenghapusSumberBisnis`,
> `TestPilihanDipegangMenggerakkanPPNPPH`), `models/sumberbisnis_test.go` (`TestSumberBisnisKiriman`,
> `TestCocokHasilPostDT`), `frontend/components/pilihSumberBisnis.test.ts` ("F4").

> ⛔ **RALAT** 2026-10-03 (putaran 2, paket P2) — dua butir lain kalimat yang sama: Survey Report →
> **K7** (tidak ada tabel di diagram grilling; tetap tidak dibangun); `Choose Business R` dan subsection
> NonProp → **paket P5** (K8). Pemilih SOB: RALAT paket P3 di atas.

## Hasil implementasi putaran 2 — 2026-10-03 (paket P2, cabang `modul/nbtreatyin/p2-layar`)

Setiap sel `DetailPolicyTreatyIn` / `DetailDeptHeadTreatyIn_UW` dicocokkan ulang dengan XML,
**termasuk syarat wadah** (`pyContainerVisibleWhen`) yang putaran 1 lewatkan. Selisih yang diperbaiki:

| Sel | Bunyi XML | Sebelumnya |
| --- | --- | --- |
| admin `.FlagRetroTreaty` | `pyVisible` `.ClaimType != 'XOL Retro'` | dibalik (`=== 'XOL Retro'`) |
| admin `.FlagPPH`, `.TypeTax`, tombol Choose Business | wadah `.ClaimType != 'XOL Retro'` | selalu |
| `.Quartal`, `.YearOfQuartal`, `.TreatyYear` kedua ("U/Y") | wadah `.QuotationData.ProportionalType = 'Proportional'` | selalu; sel U/Y tidak ada |
| `.LayerType` `.Layer` `.LayerPartType` `.LayerPart` | wadah `.QuotationData.ProportionalType = 'NonProportional'` | selalu |
| bagian uang, grid spreading, angsuran | admin `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1`; atasan `.IsNewPolicyNonProp != 1` | selalu |
| atasan `.DueTo` `.FlagPPH` `.QuotationData.NoOfferSlip` | `pyDisabled=always` | FlagPPH dan No Offer Slip dapat diisi |
| atasan `.CedingCoName` `.MarketingOfficer` | `NOTBLANK` | selalu |
| atasan `.PolicyNo`, `.ProductionDate` (General) | wadah `pyWorkPage.FlagViewPolicy = 1` (nol pengisi terjangkau) | selalu |
| atasan `.TotalSharePercentagePremium` dst. berlabel | tampil (2 desimal); di admin wadah `1=2` | tidak ada |
| admin `.RiCommOgp` … `.ResultOnp2` (8 sel) | action set: `Count*_Act(Data)` LALU `CountOGPONP_Act` | hanya langkah pertama |
| `.Deduction1/2` | `pxCurrency` (K3) | tanpa kode mata uang ("persen P29") |

- **AC 45 / 48** — medan wajib di wadah tersembunyi tidak berlaku (Pega tidak me-render selnya):
  `TypeTax` (`FlagPPH` DAN bukan XOL Retro), medan uang admin/atasan (wadah IsNewPolicyNonProp) —
  `models/layar.go` `wajibUangAdmin`, `wajibUangAtasan`; uji `TestMedanWajibIkutWadahTampil`.
- Padanan setiap tombol/aksi dengan rule XML: `docs/alat/tombol.json` (INVENTARIS bab 13).
- Format penyajian (K14) — tiket 07, AC 86.

## ⛔ RALAT putaran 3 — audit silang W1 (04-10-2026): popup `Harness/BusinessAndSOBList`

Popup tombol `Choose Business` dicocokkan ulang dengan `Section/BusinessAndSOBList` (grid AKTIF S16/S17, RD
`BrowseTreatyJoinEDM`; grid S11 RD `BrowseTreatyInDetail` dan LABEL *"For Treatyindetail join EDM"* ber-`1=2`) dan
showHarness `Section/DetailPolicyTreatyIn` (`pyWindowName`, `pySubmitData=Yes`). Selisih yang diperbaiki:

| Sel / setelan | Bunyi XML | Sebelumnya |
| --- | --- | --- |
| judul jendela | `pyWindowName` *"Business And SOB List"* | *"Choose Business"* (label tombol) |
| judul 25 kolom | LABEL baris 1: Treaty Offer ID, Contract Name, Class of Business, Source Of Business, Insured Name, Proportion Type, Treaty Type, Treaty Group, Treaty Year, Currency, Limit, Currency, Retention, Currency, EPI, Layer, *(kosong)*, Part of, *(kosong)*, Currency, MDP, Currency, Net Premium, Currency, Share RNM Value | nama kolom view (`TREATYID` …) |
| saring | `pyGridFiltering true`, kolom 2-26 `pyColumnFilteringDropDown true` (kolom tombol `false`); nol kotak cari | kotak teks *"TREATYID"* + tombol *Filter* (buatan) |
| urut | `pyGridSorting true`, kolom 2-26 `pyColumnSorting true`; kolom `.TREATYID` `pySortType ASC` | tidak dapat diurut |
| paging | `pyPageMode Numeric`, `pyPageSize 50` | seluruh baris satu halaman |
| isi | RD `BrowseTreatyJoinEDM` (view, filter H `ProportionalType` kasus, ≤ 500) | tabel TREATYINDETAIL tanpa saringan (tiket 01 RALAT W1) |
| syarat buka | wadah tombol `.ClaimType != 'XOL Retro'`, hanya layar admin | dijaga layar saja; kini juga server (409) untuk isi popup dan Choose |

Bunyi lama komentar `frontend/labels.ts` `KOLOM_BISNIS`: *"Kolom grid popup `BusinessAndSOBList` (nama kolom view, apa
adanya)."* → bunyi baru: judul VERBATIM LABEL baris 1, `kolom` = properti sel baris 2 di bawahnya. Kode:
`components/PilihBisnis.tsx`, `gridbisnis.ts`, `labels.ts`; uji `labels.test.ts` ("popup pilih bisnis"),
`gridbisnis.test.ts`, `backend/handlers/daftarbisnis_test.go`. AC tiket ini tidak berubah status; AC 53 (bagian mati
tidak dibangun) kini juga berlaku untuk grid S11 yang sebelumnya dibangun.

## ⛔ RALAT putaran 3 (paket R6, 04-10-2026) — audit silang P3 W3–W6: masukan layar dan unsur layar

Seluruhnya dibaca ulang ke XML (`Section/DetailPolicyTreatyIn`, `DetailDeptHeadTreatyIn_UW`,
`SpreadingRiskList`, `DetailPolicyTreatyInNonProportional`, `SFAPortal_OpportunitiesList`).

### W4 — medan dari sel / wadah TERSEMBUNYI tidak diterima

Bunyi lama (komentar `models.medanAdmin`, dikutip): *"DAFTAR IZIN layar admin ... isian medan yang
dapat diketik/dipilih di section"* — daftar tanpa syarat tampil: setiap medan di daftar diterima
walau selnya tersembunyi. Bunyi baru: `medanAdmin` daftar izin **berurutan dengan syarat tampil sel dan
wadahnya**, dinilai atas halaman yang sedang digabung:

| Medan | Syarat tampil (XML) |
| --- | --- |
| medan uang, `.Installment` | wadah S19 `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1` |
| `.IDCurrency` | sel `.IsNewPolicyNonProp != 1` |
| `.FlagPPH`, `.TypeTax` | wadah S7 `.ClaimType != 'XOL Retro'`; `.TypeTax` juga sel `.FlagPPH = true` |
| `.FlagRetroTreaty` | sel `.ClaimType != 'XOL Retro'` |
| `.Quartal`, `.YearOfQuartal` | wadah S14 `.QuotationData.ProportionalType = 'Proportional'` |
| `.QuotationData.IsSurveyReport` | sel `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` |

Polis NonProp baru: `PremiOgp` / `Deduction1` / `Deduction2` milik `InputPolicyTreatyInDetail_NonProp`
18–19 tidak lagi dapat ditimpa kiriman layar. `.FlagPPH` yang berubah menjalankan `RemoveTypeTax_ACT`
(change → runActivity) di server, supaya TypeTax tersembunyi tidak tertinggal.
Uji: `models/layar_masukan_test.go` `TestMedanAdminTersembunyiTidakDiterima`,
`TestFlagPPHDilepasMenghapusTypeTax`; `handlers/masukanlayar_test.go`
`TestKirimanLayarTidakMenimpaUangMasterNonProp`.

### W3 — hasil tombol Enable / Disable Input Type (pola F4, AC 49–51)

Bunyi lama (`models.medanAdmin`, dikutip): *"// hasil tombol "Enable / Disable Input Type"
(TreatyEnableDisableInput) "IsNewPolicyNonProp", "QuotationData.ProportionalType""* — diterima apa
adanya. Bunyi baru: kedua medan (`.QuotationData.ProportionalType` sel `pyReadOnly=true`;
`.IsNewPolicyNonProp` bukan sel) diterima **hanya bila `.TreatyType='XOL'`** (pyVisible tombol) **dan
sama dengan hasil `TreatyEnableDisableInput` atas halaman server** (langkah 1 "NonProportional",
langkah 2–4 selalu "0"). TreatyType bukan XOL: medan terkunci, kiriman diabaikan; nilai lain: **422**
*"Nilai … bukan hasil tombol Enable / Disable Input Type (TreatyEnableDisableInput) - medan ini
terkunci di layar"*. Uji: `TestEnableDisableHanyaHasilTombol`, `TestEnableDisableHanyaDariTombolXOL`.
Ikutan: uji R5 `TestDaftarBisnisMenyaringJenisProporsiKasus` kini memakai nilai server / hasil tombol.

### W5 — kolom hanya-baca daftar dari server; nilai bawaan sel

Bunyi lama (komentar `GabungMasukanLayar`, dikutip): *"Admin `medanAdmin` (daftar izin), beserta
daftar `DaftarDariLayar` (SpreadingRiskList, ListInstallment - baris boleh ditambah/dihapus, tombol
Add/Delete layar admin)"*. Bunyi baru:

- `.ListInstallment` (grid S45/S107 `pyEditingMode`/`pyRowEditing` `readOnly`) **tidak pernah**
  diterima dari layar. Barisnya ditulis action set server (`models.TerapkanPemicu`, sesudah
  `services.turunkan`): `.Installment` berubah → `FillPaymentInstallment`; sel uang ber-action set
  `CountOGPONP_Act` berubah → langkah 10 `SetValidateInstallment_Act` (dan langkah 9 CountSpreading).
- Grid spreading: `.PremiumSpreaded` / `.ClaimSpreaded` (`Read-only`) tidak diterima; ditulis
  `CountSpreading_Act` langkah 4.1 bila sel %Share berubah / baris dihapus / baris baru ber-%Share,
  selain itu nilai server (menurut urutan) bertahan. `[penyesuaian sadar]` baris dihapus tanpa ubahan
  %Share dihitung ulang (baris tak berkunci); urutan beberapa refresh dalam satu kiriman tak terbaca —
  `FillPaymentInstallment` memakai `BalanceDueTo` akhir.
- `pyDefaultValue` sel terbuka (dua satu-satunya di layar NB selain label mati): `.QuotationData.
  IsSurveyReport` "No", `.TypeTax` "Inclusive" — dipakai bila kosong dan sel tampil, saat layar admin
  dirender dan sesudah kiriman digabung (`models.TerapkanNilaiBawaanSel`). `[tafsiran]` (bab 6.4
  audit, "perlu cek"): `pyDefaultValue` sel Section = nilai kontrol bila properti kosong saat sel
  dirender, lalu ikut terkirim bersama form; pemicunya hanya render sel - nol Activity/DataTransform
  korpus yang MENGISI kedua medan itu (`RemoveTypeTax_ACT` hanya menghapus `.TypeTax`; rule lain
  hanya membaca `.TypeTax=="Inclusive"` / `IsSurveyReport=="Yes"`).
Uji: `TestKolomHanyaBacaSpreadingDariServer`, `TestAngsuranDariServer`, `TestNilaiBawaanSel`
(models); `TestKolomHanyaBacaSpreadingTidakDariLayar`, `TestAngsuranDariActionSetServer`,
`TestNilaiBawaanSelLayarAdmin` (handlers).

### W6 — unsur layar tanpa dasar Section XML (bab 0 butir 7)

**Dibuang** (nol sel XML): spanduk `NBStatus` layar kasus (`NBStatus` hanya kolom grid portal);
kolom portal "Position" dan "No Polis" (grid `GetListOpportunity` tidak memuatnya); judul panel
buatan "General", "Premium & Claim", "Spreading Risk", "Suggest", dan judul panel NonProp
"Spreading Risk"/"Installment" — wadahnya `NOHEADER` / `pyIncludeHeader=false` (S2 berjudul
"General" pun tanpa header). Yang tampil di XML dan dipertahankan: "Installment Data Information"
(S45/S107 `pyIncludeHeader=true`), "Limits", "Share", "Total Spreaded", "Share Facultative"
(NonProp, `pyIncludeHeader=true`), dan **LABEL Heading 4 "OGP" / "ONP"** (`pyIncludeLabel=true`
wadah S24/S25 admin, S93/S94 atasan) — bagian uang kini dikelompokkan per wadah XML.

**Dibangun (butir kecil):** ikon pengosong saring portal C[1.2] (pxIcon `pyiconclearfield`: click →
setValue `""` → postValue, **tanpa** refresh); format sel grid master NonProp (`LimitShareSummaryList
.Limit` pxNumber tanpa `pyDecimalPlaces`; `LimitFacShareSummaryList` seluruh sel tanpa Format → apa
adanya); kolom pertama FIELD kosong grid `TotalLimitIOONP` (S7) dan `TotalFacShareRnmNP` (S57).

**⛔ RALAT tautan `.Name`.** XML: kolom grid ke-2 LABEL "Name", sel `.Name` pxLink (`LabelPreview
.Name`) click → `openWorkByHandle(.pzInsKey)`. `.Name` **ditulis nol rule** di korpus NB (sisir
`PropertiesName` seluruh korpus: nol penulis properti kasus `.Name`) dan tak berkolom di delapan tabel
— tautannya selalu tanpa teks, tidak dapat diklik. Bunyi lama (`tombol.json`, dikutip): *"`.Name`
(tautan berkas) … dibangun … pages/PortalNBTreatyIn.tsx tautan"*. Bunyi baru: kolom "Name" **tidak
dirender**; aksi `openWorkByHandle` (kunci berkas yang sama) dipasang di sel "Offer No"
(`.TextNoQuotation` = pengenal kasus `NB-<n>`). `[penyimpangan sadar]` — tanpa ini berkas tak dapat
dibuka dari portal.

**Dipertahankan, dengan dasar bukan Section** (tidak ada padanan sel, tetapi dibutuhkan): panel
"History" + `GET /kasus/{id}/riwayat` — spec **AC 72** (*"Riwayat dapat dibaca berurutan waktu"*);
tombol "Kembali" — modul satu halaman (bab 0 butir 7): `openWorkByHandle` membuka berkas di area
kerja portal Pega, padanannya kembali ke daftar; tombol "Cancel" modal nomor polis — tombol bawaan
komponen `Modal` inti (menutup tanpa submit, sama dengan menutup jendela lokal ShowPolicyNoTreaty);
teks daftar kosong portal — komponen `Kosong` inti; peringatan hanya-baca — layar bagi pelaku bukan
anggota antrean (AC 11, 92).

Uji: `frontend/pages/PortalNBTreatyIn.test.tsx`, `frontend/components/DetailNonProp.test.tsx`,
`frontend/labels.test.ts` (BAGIAN), `frontend/medan.test.ts` (kelompok wadah uang).
