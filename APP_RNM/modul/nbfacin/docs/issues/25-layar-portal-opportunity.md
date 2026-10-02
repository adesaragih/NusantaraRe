# 25: Halaman depan NB FacIn — portal Opportunity (port harness `SFAPortalOpportunities`)

> Disusun agent (sesi `nusantarare-0f`) dari XML korpus atas perintah work owner 02-10-2026 — **bukan
> hasil `/to-tickets`**. Work owner menyatakan halaman depan lama (blok coverage kargo, tiket 21) salah
> tempat dan menyerahkan pengerjaan halaman depan ke sesi ini.

**What to build:** halaman yang dibuka tombol **NB FacIn** di sidebar = portal Opportunity, meniru harness
`SFAPortalOpportunities` apa adanya: judul, tombol **Create opportunity**, kotak saring + tombol
**Filter**, dan grid daftar. Halaman coverage kargo (tiket 21) **tidak dihapus**, hanya tidak lagi
menjadi halaman awal.

## Bukti `[terverifikasi]` — korpus `D:\migrasi\RNM\NB FacIn\`

| Rule (`pxInsName`) | Ruleset · versi | Peran |
| --- | --- | --- |
| `PEGACRM-PORTAL!SFAPORTALOPPORTUNITIES` (`Harness\SFAPortalOpportunities.xml`) | PegaCRM-SFA 07-22-01 | titik masuk portal |
| `PEGACRM-PORTAL!SFAPORTALOPPORTUNITIESHEADER` (`Section\SFAPortalOpportunitiesHeader.xml`) | SFAGISFW 01-01-39 | judul + tombol |
| `DATA-PORTAL!SFAPORTAL_OPPORTUNITIES` (`Section\SFAPortal_Opportunities.xml`) | SFAGISFW 01-01-39 | pembungkus |
| `DATA-PORTAL!SFAPORTAL_OPPORTUNITIESLIST` (`Section\SFAPortal_OpportunitiesList.xml`) | SFAGIS 01-01-32 | saring + grid |

Work owner mengonfirmasi 02-10-2026: versi XML ini benar.

### Yang tampil (label VERBATIM, diuji `labels.test.ts` terhadap korpus)

| Unsur | Label | Sel · baris |
| --- | --- | --- |
| Judul | `Opportunity` | Header sel 57, `pyValue` L516 |
| Tombol | `Create opportunity` | Header sel 72, `pyLabel` L2486 — `createWork` |
| Kotak saring | placeholder `NB-1234 or Name` (label `Filter Term for Opportunity` tidak dirender, `pyIncludeLabel=false`) | List sel 9, L1655 / L1622 |
| Ikon hapus isian | — (`pxIcon`) | List sel 10 |
| Tombol | `Filter` | List sel 13, L2902 |
| Grid `GetListOpportunityF` (kolom 1–8) | `Offer No` · `Name` · *(kosong: tautan `View`)* · `Group Business` · `Insured Name` · `Marketing` · *(kosong: `NBStatus`)* · `Status` | List sel 89–96 (L14643–L15727) |

### Yang tersembunyi permanen di Pega — TIDAK diport

`Stage view` / `List view` (kontainer `1=2`, Header L768), `Create Opportunity` sel 71 (`1=2`),
`All` / `Individual` / `Corporate` (`1==2`, List L3782), `Phase :` / `Proposal` / `Closed` / `Export`×3 /
`Refresh` (kontainer `NEVER`), `SFAPortal_OpportunitiesList_Header` (`1=2`).

## Keputusan agent (menunggu konfirmasi work owner)

- **B-1 — grid yang ditampilkan = `GetListOpportunityF` (grid L19).** Pega memilih satu dari tiga grid
  menurut operator: L14 bila `pyWorkBasketList(2)=='ReasTreatyInAdmin'`, L24 bila `IsOperatorLife`, L19
  untuk sisanya (= operator NB FacIn non-Life). Model peran (OQ RBAC) belum diputuskan, jadi yang dipilih
  grid untuk operator NB FacIn biasa. Literal workgroup/workbasket **tidak** diport.
- **B-2 — daftar belum punya sumber data.** RD `GetListOpportunityF` (kelas
  `ASM-FW-SFAGISFW-Work-Opportunity`) membaca tabel kerja Pega; padanannya di sistem baru menunggu jalur
  pemuatan (tiket 22–24, sesi c3). Layar menampilkan `BelumTersedia` (bukan tabel kosong / "tidak ada
  data" — pola `inti`), dan kotak saring + **Filter** dinonaktifkan selama itu.
- **B-3 — `Create opportunity` membuka halaman form Opportunity (tiket 26)** di dalam modul. Di Pega
  tombol ini `createWork` dengan `D_crmAppExtPage.WorkClass_Opportunity` + `CreateWork_Flow`; data page
  itu **tidak ada di korpus**.
- **B-4 — nol catatan pengembang di layar**, meniru keputusan work owner 30-09-2026 untuk
  treatycontractout. Keterangan "belum terverifikasi" tinggal di kode dan tiket.

## Belum terverifikasi

- Isi `D_crmAppExtPage`, `crmOpportunitiesList`, `SFAPopulateListHeader`, `TipeBisnisAct`, enam Data
  Transform filter, When `HasOppIndFilter`/`HasOppBusFilter`/`NoOppTypeFilter`/`crmIsReview` — **tidak ada
  di korpus** (23 rule, dicek dua cara).
- Tautan `Name` (`openWorkByHandle .pzInsKey`, nonaktif oleh `pyDisabledWhen` marketing/pembuat/
  `ReasFacInTeamLeader`) dan tautan `View` (`showHarness ViewOutstandingCase`, activity `GetData_ACT`
  11.214 baris, belum diurai) — belum diport; baris data memang belum ada.
- Paging efektif (20/50/10 saling bertentangan) dan nilai filter `BusinessFac` `T`/`F` tanpa kutip.

Rincian bedah lengkap (tabel sel, aksi, kondisi tampil): hasil bedah sesi 0f, uji instrumen 15/15.

**Status:** ready-for-human — dibangun 02-10-2026, menunggu tinjauan work owner atas B-1…B-4

- [x] Tombol NB FacIn membuka portal (halaman awal modul)
- [x] Judul, tombol, kotak saring, dan judul kolom VERBATIM korpus, diuji terhadap XML
- [x] Unsur tersembunyi permanen tidak dirender
- [x] Daftar tanpa sumber data tampil sebagai `BelumTersedia`, bukan tabel kosong
- [x] `Create opportunity` membuka form Opportunity (tiket 26)
- [ ] Daftar terisi — menunggu jalur pemuatan (tiket 22–24) + endpoint daftar

## Comments
