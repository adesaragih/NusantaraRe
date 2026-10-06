# 32: Daftar case NB di portal Opportunity — `GET /api/nbfacin/opportunity`

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah: permintaan sesi `nusantarare-0f`
> 03-10-2026 yang meneruskan permintaan work owner: *"nb yang sudah di create, muncul disini [portal Opportunity]. …
> harus ada case id nya, group business, insured name, status."* Frontend portal sudah memakai kontrak ini (sesi 0f).

**What to build:** endpoint berhalaman yang mendaftar case NB (Fac In) untuk grid portal Opportunity.

**Blocked by:** migrasi **182/183** (tiket 31) dijalankan work owner — kolom Marketing dibaca lewat `T_QUOTATIONDATA.MOID`
(183 ber-FK ke 182). Sebelum keduanya dijalankan, rute ini menjawab 500 (tabel tidak ada).

**Status:** ready-for-human — dibangun 03-10-2026; uji tanpa Oracle hijau; uji terhadap Oracle **tidak** dijalankan

## Kontrak

`GET /api/nbfacin/opportunity?cari=<teks>&halaman=<n>` (halaman mulai 1; kosong = 1) →
200 `{"baris":[{"caseId","name","groupBusiness","insuredName","marketing","status"}],"total","halaman","ukuran"}`;
400 halaman bukan bilangan bulat ≥ 1 (atau melebihi batas offset int32, ±143 juta) / `cari` > 255 karakter; 401 tanpa
identitas (diperiksa lebih dulu); 503 tanpa DB; 500 tanpa rincian. Cacah dan halaman = dua kueri (tanpa snapshot bersama).
`POST` rute yang sama tetap pembuat case (tiket 29).

| Medan | Sumber |
| --- | --- |
| `caseId` | `T_WORK_POLIS.ID` — hanya `LINI = 'FAC'` (case Life tidak tampil) |
| `name` | `T_NB_OPPORTUNITY.BUSINESS_PROSPECT_NAME` |
| `groupBusiness` | `T_NB_OPPORTUNITY.GROUP_BUSINESS` |
| `insuredName` | `T_M_ACCOUNT.INSUREDNAME` lewat `ACCOUNT_ID` (subkueri `MAX`; akun boleh kosong) |
| `marketing` | `MARKETINGOFFICER.CLIENTNAME` lewat `T_QUOTATIONDATA.MOID` (tiket 31); kosong bila belum disimpan |
| `status` | `T_WORK_POLIS.STATUS_WORK` (permintaan work owner) |

Pencarian "mengandung", **tidak peka huruf** (`UPPER(kolom) LIKE` pola huruf besar, `ESCAPE`, `PolaCari` tiket 27) atas
`T_WORK_POLIS.ID` dan `BUSINESS_PROSPECT_NAME`; kosong = semua. Urutan `TGL_CREATE DESC, ID DESC`.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\SFAPortal_OpportunitiesList.xml` kotak cari: `pyPlaceholder` **"NB-1234 or Name"** (L1655).
- `ReportDefinition\GetListOpportunityF.xml` (grid Pega, kelas `ASM-FW-SFAGISFW-Work-Opportunity` ⋈ `ASM-FW-GISFW-Work-NB`
  pada `A.pzInsKey = .NBHandle`): logika `A AND B AND C AND E AND D AND F AND (G OR H) AND F1` — A/B `A.pyStatusWork !=
  "Resolved-Completed"` / `"Resolved-Rejected"`, C `A.Quotation.BusinessFac = F`, D `A.Quotation.TeamGroup = Param.TeamGroup`,
  E `.TextNoQuotation Contains "NB-"`, F `A.pxCreateOperator = Param.UserIdentifier`, G `.Name Contains Param.Search`,
  H `.TextNoQuotation = Param.Search` (G dan H `pyCaseInsensitive` true), F1 `.pyID IS NOT NULL`; urut `.pxCreateDateTime`
  DESC; `pyMaxRecords` 500.

## Keputusan agent — menunggu konfirmasi

| # | Keputusan | Dasar |
| --- | --- | --- |
| A95 | Ukuran halaman **15** | pola keputusan work owner tiket 27 (butir 73.2) |
| A96 | `cari` > 255 karakter → 400 | kolom tercari terlebar `BUSINESS_PROSPECT_NAME` VARCHAR2(255) (pola A73) |
| A97 | **Selisih dengan grid Pega, dicatat:** (1) tanpa saringan pembuat (F), team group (D), status Resolved (A/B); C `BusinessFac = F` digantikan `LINI = 'FAC'` `[dugaan]`; E (`.TextNoQuotation Contains "NB-"`) dan F1 (`.pyID IS NOT NULL`) tidak diterapkan — ID case NB selalu berawalan `NB-` dan tidak NULL; (2) Pega H `=` atas **`.TextNoQuotation`** (label grid "No Penawaran" / "Text No Quotation"), di sini **mengandung** atas `T_WORK_POLIS.ID` — apakah keduanya bernilai sama `belum terverifikasi`; (3) Pega G `.Name` dipetakan ke `BUSINESS_PROSPECT_NAME` `[dugaan]` (form tidak ada di korpus, tiket 26); (4) `status` = `STATUS_WORK` (Pega: `A.NBStatusNew`); (5) urut `TGL_CREATE` + `ID` (Pega `pxCreateDateTime`); (6) tanpa batas 500 baris — berhalaman | permintaan sesi 0f / work owner; (1) `belum terverifikasi` apakah perlu |
| A98 | Daftar meminta identitas (401) | permintaan sesi 0f |

## Comments

*Tanpa nama orang, tanpa data pelanggan — uji memakai data sintetis berawalan `UJI-`.*
