# 01: Sumber data realisasi treaty — dibaca dari view relasional, gagal baca menghentikan proses

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** —
**Menutup:** AC 15 · 16 · 17 · 36 · 37 · 38 · 57 · 58 · 89 *(9 AC)* — US 21 · 23 · 24 · 37

## Hasil & nilai pengguna

Hari ini data kontrak dibaca dengan **membongkar satu dokumen teks** menjadi properti saat
halaman dibuka. `[terverifikasi]` Enam langkah Java identik melakukannya, dan bila pembongkarannya
**gagal**, galatnya **hanya ditulis ke log** — aktivitas **tetap lanjut** dengan halaman kosong
atau separuh terisi. ⛔ Pengguna tidak diberi tahu apa pun.

Sesudah tiket ini, data kontrak dibaca dari **sumber relasional yang bentuknya dapat dinyatakan**,
dan ⭐ bila pembacaan gagal, pengguna **diberi tahu** dan prosesnya **berhenti** — tidak ada lagi
berkas yang tersimpan dari pembacaan yang gagal.

## Area codebase

- Lapisan repository: pembacaan data realisasi treaty
- Lapisan repository: pembacaan penempatan keluar *(hanya baca)*
- Lapisan service: penanganan kegagalan pembacaan

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pembongkaran dokumen — **tidak dimigrasi** | 6 langkah Java di `Activity\FetchMasterTreatyIn`, `SetTreatyIn_Act`, `InputPolicyTreatyInDetail_*`, `InputPolicyTreatyOutDetail_*` |
| Medan yang dipakai laporan | `ReportDefinition\BrowseTreatyInDetail.xml` — **33 medan** |
| Penempatan keluar, hanya baca | `RDBList\BrowseTreatyOut.xml` · `RDBList\BrowseTreatyOutDetail.xml` — keduanya `SELECT` |

## ADR terkait

- **ADR-0009** — migrasi penuh, tidak ada koeksistensi dua penulis

## Acceptance criteria

- [x] **AC 15** — data dibaca dari sumber relasional, bukan dari dokumen
- [x] **AC 16** — sistem baru **tidak menulis** dokumen
- [ ] 🟡 **AC 17** — ke-**33** medan yang dipakai laporan tersedia
- [ ] 🟡 **AC 89** — nol medan yang dipakai tetapi tidak tersedia
- [x] **AC 36** — kegagalan pembacaan **menghentikan** proses
- [x] **AC 37** — kegagalan pembacaan **menampilkan galat kepada pengguna**
- [x] **AC 38** — nol kasus tersimpan dari pembacaan yang gagal
- [ ] ⛔ **AC 57** — penempatan keluar **dapat dibaca** dari konteks ini
- [x] **AC 58** — penempatan keluar **tidak pernah ditulis** dari sini

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **P29** | contoh isi dokumen — hanya untuk **migrasi data lama**, bukan untuk pembacaan baru | ⭐ **tidak menahan** |
| **19** | satu aturan **bernama** menulis penempatan keluar tetapi tidak ada perintah tulisnya | tidak menahan |

## Perintah verifikasi

1. Buka sebuah realisasi treaty — ⭐ seluruh **33** medan laporan terisi.
2. Putus sumber data di tengah pembacaan — ⭐ pengguna **melihat galat**, ⭐ dan **nol** berkas
   tersimpan.
3. Coba menulis penempatan keluar dari konteks ini — ⭐ **ditolak**.

## Catatan

⚠️ `[penyimpangan sadar]` Menghentikan proses saat gagal baca **berbeda** dari perilaku Pega.
Alasannya tertulis di `spec.md` §5.9: **kegagalan yang terlihat lebih murah daripada yang
tersembunyi**.

## ⛔ RALAT implementasi 2026-10-03

1. **Nilai master `TreatyIn.*` tidak ada di view.** Bunyi lama (Hasil): *"data dibaca dari sumber
   relasional"*. Benar untuk 33 kolom RD — tetapi `TreatyIn.RNMShareP`, `RNMShare`,
   `BrokeragePercentP`, `CurrencyList`, `INSTALLMENT`, `Limits/Share` milik JSON master (P29) dan
   **tidak punya kolom padanan**; `RNM_SHARE` view belum boleh dipakai dalam perhitungan
   (PERTANYAAN-untuk-DBA). Langkah rantai uang yang membaginya **dilewati selama nilainya kosong**
   (`models.MasterTersedia`) — penyimpangan sadar, dicatat.
2. **`InputPolicyTreatyInDetail_preACT` dibangun sebagian** — langkah 3-8, 11, 14, 15 (seluruhnya
   membaca view dan tabel acuan); langkah 9-10, 13, 16-18 (JSON master) tidak.
   `pxResults(1).CURRENCYID` bukan kolom RD maupun view → ID mata uang selalu dicari menurut nama
   (langkah 4).
3. **AC 57 tidak dapat dipenuhi seperti tertulis.** Seluruh bagian treaty keluar
   (`InputPolicyTreatyOutDetail_*`, `BusinessAndSOBListRetro`, `DetailPolicyTreatyOutNonProportional`)
   membaca JSON `M_TREATY_OUT` (P29). Data treaty keluar tidak dapat ditampilkan tanpa sumber
   relasional baru — `[terbuka]`.

## ⛔ RALAT putaran 2 (P4, 03-10-2026) — `TreatyInputPctCommSpreading`

1. **RALAT atas RALAT 2 di atas.** Bunyi lama: *"langkah 9-10, 13, 16-18 (JSON master) tidak."* Bunyi
   baru: langkah **17** (`Activity\TreatyInputPctCommSpreading.xml`, syarat preACT
   `pyWorkPage.Quotation.ProportionalType=="NonProportional"` benar → lewati) **dibangun sebagian dari
   view**. Bukti XML: langkah 1 `FetchMasterTreatyIn` memuat JSON master `where ID =
   PolicyTreatyIn.NoOffer` (`RDBList\BrowseTreatyInJoinEDM`); langkah 2 kalang `TreatyIn.Limits`
   (2.1 syarat `PolicyTreatyIn.TreatyType==.TreatyType`), 2.1.1 kalang `.Detail` (2.1.1.1 syarat
   `.TreatyGroup==PolicyTreatyIn.TreatyGroupName`): `RiCommOgp = @replaceAll(.RIOGR,",",".")` lalu
   **ditimpa** `RiCommOgp = @replaceAll(.RIONR,",",".")`. Keempat medan itu kolom view
   `TREATYINDETAILJOINEDM` (`TREATYTYPE`, `TREATYGROUP`, `RIOGR`, `RIONR`) ⇒ dibaca
   `repository.KomisiKontrak` (baris ber-`TREATYID = NoOffer`, lewat `pilihKolom`), diterapkan
   `models.TreatyInputPctCommSpreading`, dipanggil `services.PilihBisnis`. Hasil akhir = RIONR,
   RiCommOnp tidak disentuh — ditiru apa adanya. Uji: `models/komisi_test.go`,
   `handlers/logika_test.go` TestPilihBisnisKomisiOgpDariRIONRView.
2. **Tetap tidak dibangun, alasan (c):** tiga baris `SpreadingRiskList(1).SharePercentage/TreatyType/
   TreatyName = .SpreadingTotalPct/.SpreadingTypeID/.SpreadingType` — bukan kolom view (39 kolom =
   33 kolom RD + `BROKERAGE` `COMMENCEMENT` `TERMINATION` `RIOGR` `RIONR` `RNM_SHARE`), hanya di JSON
   master; pengecualian baca-JSON K8 terbatas pada medan jalur NonProp/XOL. Langkah 3
   `BreakDownSpreading_Act` (K9) dan 4 (berlabel `//`) tidak.
3. ⚠️ Urutan `Limits/Detail` dokumen tidak ada di view: baris dibaca `ORDER BY ID`, baris cocok
   terakhir menang (sama dengan kalang Pega bila hanya satu baris cocok).

## ⛔ RALAT K8 — master jalur NonProp/XOL (putaran 2, 03-10-2026)

`[keputusan work owner]` **K8** (PROMPT-NB-TREATY-IN-PUTARAN-2 bab 2; PESAN-KOREKSI-PUTARAN-2 bagian D).

1. **`[penyimpangan sadar]` atas P29.** Bunyi lama (RALAT 2026-10-03 butir 1): *"`TreatyIn.RNMShareP`,
   `RNMShare`, `BrokeragePercentP`, `CurrencyList`, `INSTALLMENT`, `Limits/Share` milik JSON master (P29) dan
   **tidak punya kolom padanan**"*. Bunyi baru: untuk jalur **NonProporsional / XOL saja**, master dibaca
   **BACA-SAJA** dari `JSONDATA` `M_TREATY_IN` / `M_TREATY_IN_EDM` — persis SQL `RDBList\BrowseTreatyIn`
   (`select JSONDATA ... from pooldata.M_TREATY_IN where ID={TreatyIn.ID} union all ... M_TREATY_IN_edm ...`)
   dan `RDBList\BrowseTreatyInJoinEDM` — di **SATU** fungsi `repository.MasterXOLDariJSON`, di balik
   `services.PembacaMasterTreaty` (kelak kontrak modul `treatyin`, PERMINTAAN-TIM-INTI bagian E). Hanya
   medan daftar `models.SkalarMasterXOL` / `models.DaftarMasterXOL` yang lolos (selebihnya dibuang di
   pengurai). Dasar: medan itu tidak ada di view, dan tabel master relasional `treatyin` (`KONTRAK`,
   `LAYER`, `BAGIAN`, `PEMULIHAN_LIMIT`, `TERMIN`, `POTONGAN`) nol baris (dicek 03-10-2026).
2. **Daftar medan K8 lawan XML.** K8 menyebut `TreatyIn.TreatyXOLList`, `Share().SpreadingListXOL`,
   `Installment`, `RetroList`, `FacultativeShare`, `FlagPPH`, `TypeTax`. Dari XML: `TreatyXOLList` adalah
   KELUARAN (`PolicyTreatyIn.TreatyXOLList`, bukan medan master); `FlagPPH`/`TypeTax` dibaca dari halaman
   POLIS (`pyWorkPage.PolicyTreatyIn.*`, isian layar admin), bukan master; `RetroList` hanya dibaca
   `TreatyNonPropSetSpreading` langkah 8-9 yang berlabel `//` — **tidak dibaca**. Selebihnya yang dibaca
   rule terjangkau (nomor langkah di `models/masterxol.go`): `Share()` beserta `GrossPremiumList`,
   `NetPremiumList`, `DeductionList`, `DeductionTotalList`, `RnmLimitList`, `SpreadingListXOL`;
   `Installment().InstallmentList`; `FacultativeShare`, `FacultativeShareList()`; `RNMShare`; `EDMState`;
   `ProportionType`; `Limits()` (pemulihan); ringkasan dan total yang ditampilkan subsection
   `DetailPolicyTreatyInNonProportional`.
3. **Kegagalan membaca master** saat pilih bisnis = **422, nol simpanan** (AC 36-38, sama dengan view).
   Saat pra-proses (`TreatyRealizationCheckXOLList`) mengikuti XML: master kosong, pesan
   `"Error fetching XolList"`.
4. **AC 57 tetap ⛔** (K8 butir 4): treaty keluar — `RDBList\BrowseTreatyOut`
   `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`.
   **AC 58 tetap ✅**: nol penulisan treaty keluar maupun master treaty masuk.
5. **RALAT butir 2 di atas:** `InputPolicyTreatyInDetail_preACT` langkah **16 dan 18 dibangun** (K8);
   langkah 17 sebagian dari view (RALAT P4 di atas); langkah 9-10 (`M_TREATY_IN_DETAIL_EDM` — bukan tabel
   K8) dan 13 tetap tidak.
6. `IsEDMInputOnNB` (NonProp langkah 10) disimpan di kolom `T_GENERAL_POLIS.IS_EDM_INPUT_ON_NB` (medan
   `PolicyTreatyIn` diagram, katalog paket penyimpanan); syarat tampil `DetailPoliciesNonProportional`
   (`!TreatyMasterInEDM || IsEDMInputOnNB == true`) selalu benar sesudah langkah 10, sehingga varian EDM
   subsection tidak terjangkau.

## ⭐ Putaran 2 — P9 (04-10-2026): SATU ukuran K8 `[menunggu konfirmasi WO]`

Dasar: tinjauan spec P9 (temuan 6–7): K8 sempat dipakai dengan dua ukuran — `models/masterxol.go` membaca medan
di luar daftar harfiah K8 (butir 2 RALAT K8 di atas: `RNMShare`, `RnmShareDeducted`, `EDMState`, `ProportionType`,
`Limits()` beserta `MDPList`/`Reinstatement_List`, tiga `*SummaryList`, tiga belas `Total*NP`), sedangkan status
`ConvertHistoryDate` menolak dengan alasan *"di luar daftar medan pengecualian K8"*. Salah satunya keliru.

**Ukuran yang dipakai di semua tempat:** K8 = **medan master yang DIBACA rule terjangkau jalur NB NonProp** —
setiap medan di `backend/models/masterxol.go` berkutip langkah XML yang membacanya; sumbernya tetap hanya
`JSONDATA` `M_TREATY_IN` / `M_TREATY_IN_EDM` (K8 butir 2), nol penulisan JSON (K8 butir 3), treaty keluar tetap ⛔
(K8 butir 4). ⚠️ `[menunggu konfirmasi WO]` — tafsiran ini melampaui bunyi harfiah daftar K8 (`TreatyXOLList`,
`SpreadingListXOL`, `Installment`, `RetroList`, `FacultativeShare`, `FlagPPH`, `TypeTax`); dicatat di
PERMINTAAN-TIM-INTI bagian F1.

Alasan di `docs/alat/status.json` dinilai ulang dengan ukuran yang sama:

| Rule | Semula | Kini | Bukti XML |
| --- | --- | --- | --- |
| `Activity/ConvertHistoryDate` | (c) *"di luar daftar medan pengecualian K8"* | **(a)** efeknya dibaca nol rule NB | terjangkau di jalur NonProp (`SetTreatyIn_Act` 12 ← `TreatyRealizationCheckXOLList` 4), tetapi satu-satunya efeknya (`TreatyIn.CommentList().Date` +7 jam) dibaca NOL rule: `CommentList` hanya ada di `Activity\SetTreatyIn_Act.xml` (langkah 9, `revisionstate==1`) dan `Activity\ConvertHistoryDate.xml` |
| `FetchMasterTreatyIn`, tiga baris `SpreadingRiskList(1)` `TreatyInputPctCommSpreading`, `CalculatePremi_Act` 1-2 | (c) | (c) tetap — jalur **Proporsional** (preACT 17 hanya bila BUKAN NonProportional; sel di wadah `.IsNewPolicyNonProp != 1`), di luar ukuran K8 | `InputPolicyTreatyInDetail_preACT` langkah 16/17 saling meniadakan |
| `InputPolicyTreatyInDetail_preACT` 9-10, 13; `RDBList/BrowseTreatyInDetailJoinEDM` | *"P29"* | (c) — JSON tabel `M_TREATY_IN_DETAIL_EDM`, di luar dua tabel K8 butir 2 | `select JSONDATA as CLASSOFBUSINESS from pooldata.M_TREATY_IN_DETAIL_EDM where ID={TreatyIn.ID}` |
| `SetTreatyIn_Act` 8-11 (`GetCurrentDate`, `SaveTreatyIn`) | (b) | (b) tetap — penulisan JSON master, dilarang K8 butir 3 | prasyarat `param.revisionstate==1`, tidak pernah dari NB |

Medan yang dibaca `models/masterxol.go` tidak berubah (setiap medan sudah berkutip langkah XML). AC tidak berubah
status.
