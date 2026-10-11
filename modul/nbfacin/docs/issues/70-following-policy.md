# 70: Following Policy Number — Search (popup), Copy Nopolis, Copy All, View Details (sesuai sistem baru)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 08-10-2026 (gambar layar tab
> General) "Following Policy Number, button search nya dibuatin berfungsi dong". Riset korpus: subagen sesi 0f (klaim inti
> diverifikasi ulang: `RDBList\BrowseDataViewPolis1.xml` L84–93, `Activity\CopyNopolisFollowing_Act.xml` L493/L540,
> `Activity\CopyLocationFollowing_Act.xml` L410).

**Status:** kontrak backend dikirim ke sesi c3 08-10-2026; frontend sesi 0f.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\Periode.xml`: label "Following Policy Number" (L19094) + tombol "Search" (L19833; `pyDisabledWhen .ConfirmBinding = 1`
  L19826; container `!IsLife`) → popup harness `ViewFollowingNB` → section `ViewOpenFollowingPolicy`. "Old Policy Number" =
  `.Following` read-only (L22283).
- Popup mode daftar: dropdown "View Option" (`pyPortal.FilterForViewRetro`, wajib, "Choose"), "Filter" (placeholder
  "Policy No / ID", Enter = cari), tombol "Refresh" (kosongkan filter), judul "Policy List"; grid `GridViewFollowingNB` 20 per
  halaman, kolom ID · Policy Number · Reff. Number · Insured Name · Source of Business · QQ Name · Start Period · End Period ·
  link "View Details" · "Copy All" · "Copy Nopolis" (View Details / Copy Nopolis disabled bila Policy Number kosong).
- `D:\migrasi\RNM\DDL\FilterForViewRetro.xml` (`DATA-PORTAL!FILTERFORVIEWRETRO`): 1 ID Pega · 2 Policy Number · 3 No Slip ·
  4 Reff. Number · 5 Insured Name · 6 Source Of Business · 7 QQ Name.
- Cari `Activity\GetListViewFollowing_act.xml` → `RDBList\BrowseDataViewPolis1.xml` atas `POOLDATA.FACINPRODUCTION` (DISTINCT,
  `NOENDORS is null`): 1 ID Pega (`IDPEGA = 'ASM-FW-GISFW-WORK '+filter`, bila filter memuat "NB-"/"RNW-"; "EDM-" →
  NOPOLIS dari baris IDPEGA itu), 2 `NOPOLIS =`, 3 `ARASAPAS.DETAIL_INVOICE`, 4 `RISLIPCEDING LIKE`, 5 `INSUREDNAME LIKE`,
  6 `SOB LIKE`, 7 `QQNAME LIKE`; filter memuat "RNM-E" → NOPOLIS dari baris `NOENDORS = filter` (menimpa opsi). Filter kosong →
  tidak mencari. Klausa disisipkan `{ASIS:}` (rawan injeksi). Hasil kosong → cadangan `json_following` (JSON).
- Copy Nopolis (`CopyNopolisFollowing_Act`): `OfferFacIn.Following = PolicyNo`, `FollowingNB = CaseID`, `IDFollowingNB =
  "ASM-FW-GISFW-WORK "+CaseID`, Obj-Save.
- Copy All (`CopyLocationFollowing_Act`): ketiga field + `PercentShare` + `Policy.Payment.RICommision`; `LocationList` DIGANTI
  lokasi polis sumber (dari `JSON_POLIS.DATA_JSON`); bisnis berbeda → lokasi tanpa FEAList / BuildingConstruction /
  ListCauseOfLoss / OccupationList / PropertyItemList / SurroundingRisk / TotalTSIList; Obj-Save.
- View Details: mode detail popup (Periode_IsUW + ViewInwardFacultativeDtl_IsUW read-only atas data JSON_POLIS, riwayat, Attachment).

## Keputusan work owner 08-10-2026

- **K70-1** "disistem baru kan tidak pakai json_data lagi, buatin sesuai sistem baru dong": seluruh fitur (Search, Copy
  Nopolis, Copy All, View Details) dibangun di atas **tabel sistem baru**, bukan JSON (ADR-0028). Cadangan `json_following`
  tidak diport.
- **K70-2** Teks View Option = `DDL\FilterForViewRetro.xml` (diserahkan work owner).
- **K70-3** Opsi 3 No Slip (ARASAPAS.DETAIL_INVOICE) **ditunda** (bersama 69b), tidak ditampilkan.

## Keputusan agent (sesi 0f)

- **A70-1** Klausa cari memakai bind parameter (tanpa `{ASIS:}`); aturan pola NB-/RNW-/EDM-/RNM-E ikut Pega.
- **A70-2** Copy All menyalin lokasi / item / coverage / deductible / spreading dari **tabel case sumber** sistem baru
  (bisnis sama = utuh; berbeda = hanya kepala lokasi, seperti daftar Property-Remove Pega); case sumber tak ada di tabel →
  galat berpesan. Konfirmasi di layar karena mengganti seluruh lokasi.
- **A70-3** View Details = popup baca-saja atas case sumber memakai endpoint GET yang sudah ada (General, Object/Coverage,
  riwayat akseptasi, Attachment); tanpa `CopypyWorkPage` (risiko merusak case terbuka di Pega tidak terbawa).
- **A70-4** Search / Copy hanya saat case boleh diubah dan `ConfirmBinding != 1` (Pega disabled).

## Backend sesi c3 (08-10-2026)

Berkas: `services/following.go`, `repository/following.go`, `models/following.go`, `handlers/following.go`; `services/objek.go`
(`opsiGantiObjek.spreading`), `handlers/kasus.go` (`bolehFollowing`).

| Endpoint | Jawaban |
| --- | --- |
| `GET /api/nbfacin/kasus/{caseId}/following?opsi=&filter=&halaman=` | `{baris, total, halaman, ukuran: 20}`; baris = api.ts `BarisFollowing` (tanggal DD-MM-YYYY) |
| `POST …/following/nopolis` `{caseIdSumber, policyNo}` | KasusNB (seperti GET kasus) |
| `POST …/following/salin-semua` `{caseIdSumber, policyNo}` | KasusNB |

- **K70-B1 Klausa cari** `[terverifikasi]` GetListViewFollowing_act 4.1-4.9 BERURUTAN (langkah sesudahnya menimpa):
  - opsi 2 NOPOLIS;
  - opsi 1 + "NB-" / "RNW-" IDPEGA = "ASM-FW-GISFW-WORK " + filter (4.2: when2 NB- = kode 5 Skip Whens);
  - filter memuat "RNM-E" -> NOPOLIS dari baris NOENDORS = filter (4.4);
  - opsi 1 + "EDM-" -> NOPOLIS dari baris IDPEGA itu (4.5);
  - opsi 4 RISLIPCEDING LIKE / 5 INSUREDNAME / 6 SOB / 7 QQNAME LIKE `%filter%`.

  ⚠️ **Ralat kontrak**: "RNM-E menimpa opsi apa pun" hanya berlaku untuk opsi 1 / 2. Langkah 4.6-4.9 (opsi 4-7) ada SESUDAH 4.4
  dan menimpanya balik. Opsi 1 tanpa NB- / RNW- / EDM- / RNM-E -> klausa kosong -> tanpa kueri (4.10). Filter kosong -> `[]`
  (langkah 3; filter dipangkas spasinya). Semua nilai lewat BIND (A70-1); `%` / `_` di filter tetap wildcard seperti Pega.
  Urutan `[keputusan agent]`: NOPOLIS, CASEID, lalu kolom lain. DISTINCT atas tanggal mentah, diformat sesudahnya. 400: opsi 3 /
  di luar 1-7, halaman < 1 / bukan bilangan. GET tanpa penjaga tahap (baca saja).
- **K70-B2 Penjaga POST** (A70-4): `CONFIRM_BINDING = "1"` -> 409 "Following Policy tidak dapat diubah setelah Binding.", diperiksa
  lebih dulu sehingga TL / UW tahap Policy (ConfirmBinding 1) juga 409. Selain itu aturan tab (`alasanUbah`: tahap + workbasket,
  case selesai). !IsLife selalu benar (nbfacin = FIRE). `bolehFollowing` = penjaga yang sama (lunak).
- **K70-B3 Polis diperiksa** `[keputusan agent]`: `caseIdSumber` + `policyNo` wajib ada di FACINPRODUCTION (IDPEGA, NOPOLIS,
  NOENDORS kosong), selain itu 404 "Polis Following tidak ditemukan.". Pega menulis nilai grid tanpa memeriksa.
- **K70-B4 Copy Nopolis**: `T_GENERAL_POLIS.FOLLOWING` = policyNo (baris General belum ada -> disisipkan).
- **K70-B5 Copy All** (A70-2): satu transaksi. FOLLOWING + PERCENT_SHARE (spreading sumber) + PAY_RI_COMMISION (deduction
  sumber) ditulis lebih dulu (`opsiGantiObjek.sebelum`). Seluruh lokasi diganti objek case sumber lewat `gantiObjek` (jalur Save tab
  Object). Spreading case sumber dipasang lewat `pertahankanSpreading` (posisi + oldId sama), dihitung ulang dengan % Share sumber.
  - Tabel tersentuh (dicek ke `ObjekOracle.tabel()`): T_GENERAL_POLIS (FOLLOWING / PERCENT_SHARE / PAY_RI_COMMISION); T_WORK_POLIS
    (sentuh); hapus + sisip T_LOCATIONLIST, T_PROPERTY, T_RISKLOCATION, T_BUILDINGCONSTRUCTION, T_SURROUNDINGRISK,
    T_PROPERTYITEMLIST, T_OCCUPATIONLIST, T_TABLEOFLIMIT, T_FEALIST, T_LISTCAUSEOFLOSS, T_COINSDATA, T_COVERAGELIST,
    T_DEDUCTIBLELIST, T_SPREADINGLIST (+ TSI / premi NR coverage lewat penyimpan spreading).
  - ⚠️ Penyimpangan `[keputusan agent]`: premi coverage dihitung ulang dengan periode case penerima. Ini aturan Save tab Object;
    Pega Page-Copy menyalin nilai apa adanya sampai Save berikutnya. CurrencyList / Payment TIDAK disentuh (Pega pun tidak).
  - Case sumber tidak ada di tabel sistem baru (polis lama Pega) -> 404 "Case sumber tidak ditemukan di sistem baru; Copy All
    hanya untuk case yang dibuat di sistem baru.".
  - Bisnis selalu sama (FIRE), jadi cabang "beda bisnis" (langkah 15) tidak berlaku. Isian objek sumber diperiksa ulang
    (`periksaObjek`, 400).
- **K70-B6 View Details** (A70-3): GET kasus / objek / spreading / lampiran tidak punya penjaga workbasket. Case lain (Resolved /
  bukan workbasket akun) terbaca tanpa 403. Hanya 401 tanpa identitas.

### Pertanyaan terbuka untuk work owner (dari sesi c3)

- **[pertanyaan terbuka] FollowingNB / IDFollowingNB tanpa kolom.** `CopyNopolisFollowing_Act` / `CopyLocationFollowing_Act`
  langkah 9 menulis `.FollowingNB` (CaseID sumber) dan `.IDFollowingNB` ("ASM-FW-GISFW-WORK " + CaseID). `Section\Periode.xml`
  L21181-21243 menampilkan isian "Following NB" yang, saat Refresh, MEMANGGIL `CopyLocationFollowing_Act` dengan `NB =
  .FollowingNB`. Sistem baru hanya punya `T_GENERAL_POLIS.FOLLOWING`. Menyimpan keduanya butuh migrasi baru (mis. `FOLLOWING_NB`,
  `ID_FOLLOWING_NB`) — menunggu persetujuan; tidak dibuat.
- `[catatan]` Copy All pada polis yang lahir di Pega (tanpa tabel sistem baru) tidak dapat dilakukan (404) sampai data lama
  dimigrasikan.

**Uji**:
- services: `TestJenisCariFollowing` (15 kasus urutan langkah + 4 opsi tak sah), `TestCariFollowing`,
  `TestSalinNopolisFollowing` (8 penjaga), `TestSalinSemuaFollowing` (objek sumber, spreading sumber terpasang, PS / RI, 404);
- repository: `TestSQLFollowing`;
- handlers: `TestFollowingHandler` (kunci HalamanFollowing / BarisFollowing = api.ts dua arah + urutan, 400 / 401 / 404 / 409,
  bolehFollowing);
- diralat: jumlah kunci KasusNB 15 -> 16 (`TestKunciKasusSamaDenganApiTs`, `TestBacaKasus`), akhiran JSON (`TestAkumulasiUWHandler`,
  `TestKunciPeriodeHandler`).

**Mutasi 19 / 19 mati** (tiga ditulis ulang supaya terkompilasi).

## Keputusan work owner 09-10-2026 (lewat sesi 9d)

- **K70-4** FollowingNB / IDFollowingNB **TIDAK perlu kolom dulu**. Hanya `T_GENERAL_POLIS.FOLLOWING` yang ditulis (Copy Nopolis /
  Copy All); tanpa migrasi. Isian "Following NB" `Section\Periode.xml` (Refresh -> CopyLocationFollowing_Act) tidak diport.
