# 40: Sub-tab Occupation — grid okupasi, Choose Occupation, Choose Class of Construction

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026 (urutan
> sub-tab objek, lihat tiket 39). Belum ada tangkapan layar; tampilan diturunkan dari XML.

**What to build:** sub-tab **Occupation** di baris objek FIRE. Isinya grid `.Property.OccupationList` (Occupation
ID · Occupation Name · Class Of Construction) dengan Tambah / Hapus. Baris yang dibuka berisi tombol Choose
Occupation, ID dan Name baca-saja, tombol Choose Class of Construction, dan Class of Construction baca-saja bertanda
wajib. Datanya ikut Save tab Object.

**Blocked by:** — (DDL TABLEOFLIMIT belum ada; lihat Kontrak).

**Status:** frontend selesai 03-10-2026 (uji hijau); backend → sesi c3 (sesudah tiket 39).

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\OccupationList.xml`: grid `.Property.OccupationList` (kelas `Data-Occupation`, master-detail, flow action
  `OccupationItemFacIn_FlowAction`, tanpa paging), kepala sel 14–16. Add = ikon → addRow + `SetCodeRiskExposure_DT`.
  Delete = ikon → deleteRow + `SetCodeRiskExposure_DT` / `GetLowestPctLimit_DT`.
- `Section\OccupationItemFacIn_Section.xml`: sel 3 "Choose Occupation", 4 Occupation ID, 5 Occupation Name, 8 "Choose
  Class of Construction", 9 Class of Construction (wajib). Medan baca-saja.
- `Section\ChooseOccupation.xml`: "Search Name/ID" (sel 1); grid RD `BrowseOccupationFacInFIRE_RD` (TYPE FIRE;
  DESCRIPTION / OLDID = kata cari), kolom ID / Name, 20 per halaman; Choose → `SetDataOccupation` (OldID, Name,
  KDRiskExposure).
- `Activity\SetDataOccupation.xml`: OccupationId = OldID, OccupationName = Name, TableOfLimit.Category = KDRiskExposure
  03 → III, 02 → II, 01 → I, lainnya "".
- `Section\ChooseClassofContraction.xml`: RD `BrowseTableOfLimit_RD` (Tahun = `pyWorkPage.OfferFacIn.CurrentYear`,
  Bizcode = `QuotationData.BusinessCode`, Category = baris); kolom Description (Limit tersembunyi); Choose →
  `SetDataClassofConstraction` (Description, PctLimit) + `GetLowestPctLimit_ACT`.
- `Activity\GetLowestPctLimit_ACT.xml`: PctLimit terendah lintas objek → `OfferFacIn.Parameters.LowestPctLimit` dkk.
  (dipakai `CountFormulaRNM_Act` untuk MaxTreatyCapacity).
- Data contoh: OccupationList 332 baris; TableOfLimit.PctLimit berkoma desimal.

## Keputusan agent

- **L-1** `GetLowestPctLimit_ACT` dihitung di backend / tahap kapasitas treaty, bukan di layar.
- **L-2** Tanda wajib Class of Construction hanya penanda. Save Pega (`SaveFacIn_Act`) tidak memvalidasi; penolakan
  = tahap Submit.
- **L-3** Popup Occupation memuat seluruh okupasi FIRE saat dibuka dan menyaring saat mengetik (20 per halaman di
  layar).
- **L-4** Popup tampil bergantian, tidak bertumpuk (pola E-6). Choose Class of Construction sebelum Occupation dipilih
  ditolak dengan pesan.
- **L-5** PctLimit disimpan sebagai TEKS apa adanya (Pega berkoma desimal; bukan uang).

## Kontrak

- `ObjekFire.occupations: OkupasiObjek[]` (`occupationId`, `occupationName`, `category`, `constructionClass`,
  `pctLimit`), ikut `GET`/`PUT …/objek`.
- `GET /api/nbfacin/occupation?cari=` (tiket 38) → baris tambah `kdRiskExposure`; `cari` kosong = semua okupasi FIRE
  (≤ 500).
- `GET /api/nbfacin/kasus/{caseId}/table-of-limit?category=` → `{ baris: [{ description, pctLimit }] }` — sumber
  TABLEOFLIMIT. DDL belum ada; contoh data `DDL\TABLEOFLIMIT.xml` (ID, BIZCODE, NOTE, TAHUN, CATEGORY, DESCRIPTION,
  PCTLIMIT).

## Acceptance criteria

- [x] Grid + Tambah / Hapus + form detail + dua popup; label diuji ke korpus.
- [x] Choose Occupation mengisi ID, Name, Category; Choose Class of Construction mengisi Description + PctLimit.
- [ ] Backend: simpan / baca okupasi (T_OCCUPATIONLIST rancangan?), endpoint table-of-limit, `kdRiskExposure`.
