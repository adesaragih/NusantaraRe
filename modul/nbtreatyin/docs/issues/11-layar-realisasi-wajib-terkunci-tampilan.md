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
