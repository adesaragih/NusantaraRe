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

## Keputusan work owner 07-10-2026 - daftar hanya case buatan sendiri

"di menu nb facin, list yang muncul hanya list yg dia buat sndri" - `GET /api/nbfacin/opportunity` disaring
`T_WORK_POLIS.CREATE_OP` = akun login (diminta ke c3; frontend tidak berubah).

### Backend (sesi c3, 07-10-2026)

- `repository/portal.go`: saringan `w.LINI = :1 AND w.CREATE_OP = :2` (bind, tanpa literal) di kueri halaman DAN kueri cacah
  - total / jumlah halaman dihitung sesudah saringan. Cari `:3/:4`, OFFSET / FETCH bergeser `:3/:4` atau `:5/:6`.
  `services/portal.go` meneruskan `pelaku.AkunID` (`inti.PelakuDari`); tanpa identitas tetap 401 (sudah ada sebelumnya).
- Uji: teks SQL kedua kueri (repository), akun A tidak melihat case akun B dan total = milik sendiri (services, pengganti
  repository yang meniru saringan), 401 tanpa sesi (handler, sudah ada). Mutasi: saringan dibuang dari SQL / AkunID tidak
  diteruskan - keduanya menggagalkan uji.
- ⚠️ **Case lama / hasil muat data** `[terverifikasi]` pada kode loader: `T_WORK_POLIS.CREATE_OP` diisi dari
  `POOLDATA.JSON_POLIS.USERNAME` (aturan V-48, `services/loader/aturan.go`) - pengenal OPERATOR Pega, bukan LOGIN_ID akun
  sistem baru. Bentuknya di sampel `DDL\CONTOH` (`pxCreateOperator`, 114 berkas): huruf besar tanpa spasi / '@', 6-16
  karakter (3 berkas huruf campuran). Case seperti itu hanya tampil bila LOGIN_ID akun baru SAMA PERSIS dengan pengenal
  operator Pega (peka huruf) - kesamaan itu `belum terverifikasi`; bila tidak sama, case tersebut tidak terlihat oleh siapa
  pun di portal ini. Apakah baris FAC hasil muat sudah ada di DEV juga belum terverifikasi (agent tidak membaca DEV). Baris
  T_WORK_POLIS milik premiumlistlife sebelum 059 ber-CREATE_OP kosong, tetapi ber-LINI bukan FAC (tidak tampil di portal ini
  sejak awal). Case yang dibuat aplikasi baru (tiket 29) selalu ber-CREATE_OP = LOGIN_ID pembuat.
  `[pertanyaan terbuka]` work owner / tim login: pemetaan operator Pega -> LOGIN_ID untuk case lama.

### Terjawab 07-10-2026 - case lama tidak masuk portal ini

Work owner: "case lama yang sudah masuk ke json_polis tidak akan masuk ke menu nb facin ini, menu ini hanya untuk yg lgi proses
input." Pertanyaan terbuka c3 (case hasil muat ber-CREATE_OP operator Pega tak terlihat siapa pun) DITUTUP: portal NB FacIn hanya
untuk case yang sedang diinput di aplikasi baru (CREATE_OP = LOGIN_ID sejak tiket 29); tidak perlu pemetaan operator Pega →
LOGIN_ID untuk portal ini.

## Kolom NB Status (work owner 07-10-2026: "disini ditambahin NB status dong, biar tau nb nya lgi dimana")

`[terverifikasi]` Grid `GetListOpportunityF` sel 95 = `A.NBStatus` dengan judul KOSONG. Kini (sesi 0f) berjudul "NB Status"
(permintaan work owner) dan berisi `BarisCaseNB.nbStatus` = `T_GENERAL_POLIS.NB_STATUS` (tiket 57); belum pernah di-Submit =
"NEW NB" (`Flow\InputInwardFacultativeOffer.xml` awal flow). Backend diminta ke c3; GET portal kini butuh migrasi 208.

### NB Status di daftar (permintaan work owner 07-10-2026 "disini ditambahin NB status dong, biar tau nb nya lgi dimana")

Backend (sesi c3): tiap baris `nbStatus` = `T_GENERAL_POLIS.NB_STATUS` (migrasi 208, ditulis Submit tiket 57) lewat
`LEFT JOIN … g ON g.ID = w.ID` (PK = ID kasus, satu baris per kasus - cacah tidak terpengaruh, kueri cacah tidak berubah);
kosong / NULL → "NEW NB" `[terverifikasi]` (`Flow\InputInwardFacultativeOffer.xml`: awal flow men-set `.NBStatus = "NEW NB"`
dan `.NBStatusNew = "NEW NB"`). Kolom `status` (STATUS_WORK) tetap. ⚠️ Migrasi 208 kini **WAJIB sebelum biner** - tanpa
kolom NB_STATUS kueri portal ORA-00904 (pembacaan tidak dibuat tahan: pembacaan bersyarat katalog menambah kueri per
permintaan; urutan deploy cukup). Uji: SQL (kolom + JOIN), NB_STATUS apa adanya / kosong → "NEW NB" (services), kunci JSON
`nbStatus` (handler); mutasi bawaan "NEW NB" → tertangkap.

### Kotak masuk Beranda (permintaan work owner 07-10-2026 "ini kenapa tidak ada ya, padahal ini login teamleader")

Backend sesi c3, pola nbtreatyin (keputusan work owner 06-10-2026 "beranda menunjukkan berapa banyak case yang masuk di akun
dia, mengikuti workbasket"):
- `GET /api/nbfacin/kotak-masuk` → `[{workbasket, nama, jumlah}]`: satu baris per workbasket ANTREAN NB Fac In yang DIPEGANG
  akun (`inti.Pelaku.Peran`), urut antrean - ReasFacInMarketing, ReasFacInTeamLeader, lalu kesembilan tujuan Banding
  SetBanding_ACT (tiket 57; `TestAntreanKasusNB` menjaga setiap tujuan Submit ada di antrean). Workbasket di luar antrean
  (modul lain) tidak ditampilkan. Jumlah 0 tetap dikirim (Beranda menyembunyikannya). Nama = `M_WORKBASKET.NAME`, kosong /
  tak terdaftar → kode.
- `GET /api/nbfacin/kotak-masuk/kasus?workbasket=` → `[baris portal]` (bentuk sama dengan `GET /api/nbfacin/opportunity`, NB
  status kosong → "NEW NB"); workbasket kosong = semua yang dipegang; workbasket tidak dipegang / bukan antrean → 403; urut
  `TGL_UPDATE` terbaru dulu (NULLS LAST), paling banyak 500.
- Saringan (`repository/kotakmasuk.go`): `LINI = 'FAC'`, `STATUS_WORK` bukan `Resolved-%` (Decline tiket 57 →
  "Resolved-Decline"), lalu ATAU: (a) workbasket marketing → POSITION "Offer" / "ReasFacInMarketing" DAN `CREATE_OP` = akun
  (selaras portal "hanya yang dia buat sendiri"; keduanya dijumlah ke ReasFacInMarketing); (b) workbasket lain → POSITION =
  workbasket (seluruh case, siapa pun pembuatnya). Tanpa sesi 401; tanpa basis data 503. Tanpa migrasi.
- Uji: `TestKotakMasuk` / `TestDaftarMenunggu` (TL melihat seluruh case Team Leader; selesai tidak; marketing hanya buatannya;
  akun tanpa workbasket antrean → []; 403 workbasket asing; 401), `TestSQLKotakMasuk` (bind berurutan), `TestKotakMasukHandler`
  (bentuk JSON larik); mutasi 5 titik → tertangkap.

### Kolom inbox "New Business" (work owner 07-10-2026 "di tampilan inbox, keluarin kolom2 ini", gambar inbox Pega)

`GET /api/nbfacin/kotak-masuk/kasus` kini `[BarisKotakMasukNbFacIn]` (bukan baris portal); cacah kotak masuk tidak berubah.
Report inbox Pega tidak ada di korpus - kolom dari gambar (pola nbtreatyin). Sumber per kunci (backend sesi c3):

| Kunci | Sumber | Catatan |
| --- | --- | --- |
| `caseId` | T_WORK_POLIS.ID | |
| `insuredName` | T_M_ACCOUNT.INSUREDNAME lewat T_NB_OPPORTUNITY.ACCOUNT_ID (MAX, pola portal) | |
| `businessName` | T_NB_OPPORTUNITY.CLASS_OF_BUSINESS | nama Class Of Business case |
| `sourceOfBusiness` | T_QUOTATIONDATA.SOB_NAME | |
| `cedingCompany` | T_QUOTATIONDATA.CEDING_CO_NAME | gabungan ";" seperti General |
| `inceptionDate` | T_GENERAL_POLIS.START_DATE_TIME → DD-MM-YYYY WIB | Begin date |
| `shareRnm` | T_GENERAL_POLIS.PERCENT_SHARE (teks desimal TM9) | sumber yang sama dengan % Share RNM tab Spreading / Cedant |
| `marketing` | MARKETINGOFFICER.CLIENTNAME lewat MOID | |
| `status` | padanan `.IsCedingConfirm` dari T_WORK_POLIS.FLAG_ONGOING_POLICY: "0" → "Offer", "2" → "Binding", lainnya "" | keputusan work owner 07-10-2026 ("setuju , buatin jangan halu"; semula STATUS_WORK). `[terverifikasi]` `Activity\InputOfferFacInEngineer_preACT.xml` langkah 3 (IsCedingConfirm "Offer" bila FlagOnGoingPolicy 0, tiap layar dibuka); konektor `Flow\InputInwardFacultativeOffer.xml` men-set "Binding" bersama FlagOnGoingPolicy 2 (Decision29 / Decision35). "Policy" (SetBanding_ACT jalur EDM) dan "Accepted" (hanya baris riwayat SetValidateDate 3.2) TIDAK ditebak - belum ada padanan tersimpan `[pertanyaan terbuka]` tahap binding / EDM |
| `createdDate` | TO_CHAR(T_WORK_POLIS.TGL_CREATE, 'DD-MM-YYYY') | |
| `offeringDate` | T_GENERAL_POLIS.OFFERING_DATE (Pega YYYYMMDD) → DD-MM-YYYY | |
| `validityDate` | "" | `.OfferFacIn.DateValidity` (GetDateValidity_ACT, dipanggil pemeriksa spreading tahap UW) tidak dihitung / disimpan; kolom rancangan DATE_VALIDITY belum dibuat `belum terverifikasi` |
| `statusInbox` | T_GENERAL_POLIS.NB_STATUS; kosong → "NEW NB" | |
| `aksiTerakhir` | APPROVAL + " by " + PIC baris HISTORYAKSEPTASIPRODUCTION TERAKHIR (KEEP DENSE_RANK LAST atas TGL_INP, NOURUT; IDPEGA = pzInsKey) | `Activity\SetNBStatus_Act` langkah 6 (StsApproval "Accept by " … + ViewSuggest(<LAST>).PICSuggest); label APPROVAL tersimpan dipakai apa adanya; tanpa riwayat → "" |
| `tglUpdate` | TO_CHAR(T_WORK_POLIS.TGL_UPDATE, 'YYYY-MM-DD HH24:MI:SS') | SYSDATE server basis data, dianggap WIB `[dugaan]` |

Uji: `TestSQLKotakMasuk` (kolom berurutan = pemindai, bind :1 / :2 awalan pzInsKey lalu saringan), `TestKeBarisInbox` (pemetaan),
`TestKotakMasukHandler` (kunci JSON persis).
