# 60: Underwriter mengubah Accumulation coverage

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 08-10-2026 (gambar layar form
> Coverage: Accumulation Code / Address, Choose Accumulation Code, Copy Accumulation) "di uw juga bisa ubah acumulation yang
> ada di tab spreading". Survei korpus oleh subagen (pengurai XML expat, atribut dipasangkan per `rowdata`).

**Status:** frontend selesai 08-10-2026 (uji hijau); backend diserahkan ke sesi c3.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- Di section UW, akumulasi diubah dari **tab Coverage**, bukan tab Spreading: `Section\InputInwardFacultativeDtl_IsUW.xml`
  tab Coverage (L13905-19237) → `CoverageList_IsUW` (L14181) → `FlowAction\CoverageListFire_FlowAction_IsUW.xml` →
  `Section\PropertyItemListCoverage_IsUW.xml` (pyEditAction `PropertyItemCoverageFacIn_FlowAction_IsUW`, L2687 / L9586) →
  `Section\InputCoverageFire_IsUW.xml` → `Section\CoverageItem_IsUW.xml`.
- Tab Spreading (L26286-38463) tidak memuat field akumulasi yang dapat diubah; hanya tombol "Total Accumulation"
  (`TotalAccumulation_Act` + harness `TotalAccumulation_FacIn` `pyReadOnly=Yes`) — ringkasan (K48 S-3, belum dibangun).
- `CoverageItem_IsUW`: `.AccumulationCode` / `.AccumulationDescription` Read-only (L6775 / L6943), diubah hanya lewat tombol;
  kontainer tombol Choose / Copy Accumulation visible bila PositionNote ∈ {ReasFacInUnderwriting, ReasFacInSeniorUnderwriting,
  ReasFacInJuniorUnderwriting, ReasFacInJuniorUnderwritingA} (L7247); seluruh field data coverage lain Read-only.
- Register `_oq011` (korpus `RNM_BRD`): `COVERAGEITEM_ISUW`, `INPUTCOVERAGEFIRE_ISUW`, `PROPERTYITEMCOVERAGEFACIN_FLOWACTION_ISUW`
  tidak tercatat; `PROPERTYITEMLISTCOVERAGE_ISUW` (#354) dan `INPUTINWARDFACULTATIVEDTL_ISUW` (#260) tercatat berkonflik — hanya
  varian NB FacIn yang dipakai.

## Keputusan

- **K60-1** `[keputusan agent]` Lokasi = tab Coverage (sesuai korpus), walau pemicu menyebut "tab spreading"; gambar layar
  pemicu memang form Coverage. Konfirmasi work owner bila yang dimaksud lain.
- **K60-2** `[keputusan work owner 08-10-2026]` Posisi = enam posisi When `IsUWSpread` (sama dengan tiket 59; backend memakai
  `rules.Eval("ISUWSPREAD")`): "DepHead dan UnderwritingFinancial juga bisa dong ubah accumulation". ⚠️ RALAT rancangan awal
  (empat posisi kontainer tombol `CoverageItem_IsUW` L7247 `[terverifikasi]`, tanpa DepHead / UnderwritingFinancial) — kini
  penyimpangan sadar dari L7247. `ReasFacInJuniorUnderwriting` dibaca `ReasFacInJuniorUnderwritingB` (ADR-0032).
- **K60-3** `[keputusan agent]` Save tab Coverage di UW memakai endpoint sempit `PUT …/akumulasi` yang hanya menulis
  Accumulation Code / Description; field lain diabaikan (di Pega Read-only).

## Frontend (sesi 0f)

- `components/BacaSajaKasus.tsx`: prop `kecuali` + `PILIH_UBAH_AKUMULASI` (`[data-ubah-akumulasi]`).
- `components/FormCoverage.tsx`: tombol Choose / Copy Accumulation dan pembungkus `PopupAkumulasi` bertanda `data-ubah-akumulasi`.
- `components/TabCoverage.tsx`: prop `hanyaAkumulasi`; Save → `simpanAkumulasi`; tombol Save bertanda saat `hanyaAkumulasi`.
- `pages/InwardFacultative.tsx`: `akumulasiTerbuka` (`KasusNB.ubahAkumulasi`), `kecualiDetail` untuk tab Coverage; banner
  `TEKS_INWARD.bacaSajaKecualiSpreadingAkumulasi`.
- `api.ts`: `KasusNB.ubahAkumulasi`, `simpanAkumulasi`.
- Uji: `BlokEmailType.test.tsx` (terapkanBacaSaja + kecuali), `InwardFacultative.test.tsx`, `TabCoverage.test.tsx`.

## Backend (diminta ke sesi c3, 08-10-2026)

1. Tahap `uwAkumulasi` (K60-2) di `services/bacasaja.go`; workbasket = POSITION.
2. GET kasus `ubahAkumulasi`.
3. `PUT /api/nbfacin/kasus/{caseId}/akumulasi` (badan = PUT objek; bentuk + OLDID sama; hanya dua kolom akumulasi ditulis).
4. PUT objek tetap 409 di uw; endpoint pencarian / Add New akumulasi tidak dijaga tahap.

## Terbuka

- "Total Accumulation" di tab Spreading (K48 S-3) belum dibangun.
- Isi `SearchAccumAct`, `CopyAccumulationCode_Act`, `SetFilterAccum_DT`, `ChooseAccumulation_FacIn` tidak ditelusur ulang untuk
  varian UW (sama dengan marketing menurut atribut aksi).

## Backend sesi c3 — selesai 08-10-2026

- **K60-B1 Tahap:** `services/akumulasiuw.go` `alasanUbahAkumulasi`.
  - Stage = `tahapUW`, yaitu `rules.Eval("ISUWSPREAD")`, aturan yang sama dengan tab Spreading. Workbasket wajib = POSITION (403).
  - Tahap lain (marketing / TL / Decline / nama lama ReasFacInJuniorUnderwriting) -> 409 `ErrAkumulasiHanyaUW`. Marketing / TL
    tetap memakai PUT …/objek.
  - > **RALAT K60-2** — penyimpangan sadar `[keputusan work owner 08-10-2026 "DepHead dan UnderwritingFinancial juga bisa dong
    > ubah accumulation"]`. Rancangan awal memakai empat posisi ekspresi visibilitas `Section\CoverageItem_IsUW.xml` L7247
    > (Underwriting, Senior, Junior → B, Junior A; tanpa DepHead / UnderwritingFinancial). Kini enam posisi IsUWSpread; daftar
    > tulisan tangan dihapus.
- **K60-B2 GET kasus:** `ubahAkumulasi` (bool, selalu ada, sesudah `ubahBarisSpreading`) = uw + pemegang workbasket. Lunak:
  galat baca tahap -> false; tanpa pembaca tahap -> false (berbeda dari `bolehUbahSpreading` / `ubahBarisSpreading` yang true
  saat penjaga tidak aktif — tanpa pembaca tahap tidak ada tahap uw).
- **K60-B3 `PUT /api/nbfacin/kasus/{caseId}/akumulasi`.** Badan `{baris: ObjekFire[]}` diurai dengan pengurai yang sama dengan
  PUT …/objek (`uraiBadanObjek`). Jawaban = daftar objek baca ulang.
  - Bentuk: jumlah lokasi / item / coverage per posisi + oldId dibandingkan dengan pohon coverage tersimpan (`BacaSpreading` di
    transaksi yang sama: ID baris coverage, urut SEQ_NO = posisi GET …/objek, seperti K48-7). Beda -> 409 `ErrObjekBerubah`
    "Objek berubah sejak tab Coverage dimuat - muat ulang tab", tanpa tulisan.
  - Yang ditulis: hanya `T_COVERAGELIST.ACCUMULATION_CODE` / `ACCUMULATION_DESCRIPTION` per ID
    (`repository.ObjekOracle.UbahAkumulasiCoverage`, teks kosong -> NULL seperti PUT objek) dan TGL_UPDATE case. Medan lain badan
    diabaikan: tanpa ganti objek, tanpa hitung ulang premi / coverage / spreading.
  - Validasi = aturan PUT …/objek untuk dua medan itu (`lebarCoverage`: 50 / 500 byte -> 400). ⚠️ PUT objek TIDAK memeriksa kode
    akumulasi terhadap master; Pega memang hanya mengisinya lewat popup Choose / Copy. Tidak ada aturan baru.
  - Kolom sudah ada (migrasi 195). Tanpa migrasi.
- **K60-B4 Tidak dijaga tahap** (diperiksa, tidak diubah): GET `/akumulasi`, `/akumulasi/saran/{jenis}`, `/akumulasi/czone`,
  `/akumulasi/zipcode` (tanpa pelaku) dan POST `/akumulasi` (Add New master: hanya identitas, tanpa caseId). PUT …/objek tetap
  409 di semua tahap uw.
- **K60-B5 ADR-0032:** pemetaan workbasket dipindah ke paket daun `services/rules/peta` (`peta.Workbasket`, `peta.Petakan`) dan
  dipakai generator `bangkit` (`PetaWorkbasket = peta.Workbasket`). Registry bangkitan tidak berubah
  (`TestRegistrySamaDenganKeluaranBangkit` dengan korpus lulus).
- **Uji:**
  - `services/akumulasiuw_test.go`: matriks tahap × workbasket — enam posisi uw lolos, nama lama / marketing / TL / Decline 409,
    non-pemegang 403; hanya dua kolom ditulis menurut posisi; bentuk / oldId beda 409; lebar 400; tanpa pembaca tahap.
  - `handlers/spreadinguw_test.go`: `TestAkumulasiUWHandler` (HTTP, PUT objek 409 di uw, GET `ubahAkumulasi`, 401, 400, 409);
    `TestKunciKasusSamaDenganApiTs` (14 kunci, dua arah + urutan).
  - `repository/akumulasicoverage_test.go`.
  - Mutasi: 10, semua tertangkap.
