# 40: Sub-tab Occupation — grid okupasi, Choose Occupation, Choose Class of Construction

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026 (urutan
> sub-tab objek, lihat tiket 39). Belum ada tangkapan layar; tampilan diturunkan dari XML.

**What to build:** sub-tab **Occupation** di baris objek FIRE. Isinya grid `.Property.OccupationList` (Occupation
ID · Occupation Name · Class Of Construction) dengan Tambah / Hapus. Baris yang dibuka berisi tombol Choose
Occupation, ID dan Name baca-saja, tombol Choose Class of Construction, dan Class of Construction baca-saja bertanda
wajib. Datanya ikut Save tab Object.

**Blocked by:** — (DDL TABLEOFLIMIT belum ada; lihat Kontrak).

**Status:** frontend selesai 03-10-2026 (uji hijau); backend sebagian 03-10-2026 (sesi c3: okupasi + `kdRiskExposure` selesai, migrasi 189 ditulis BELUM dijalankan; endpoint table-of-limit menunggu DDL TABLEOFLIMIT dan asal BusinessCode).

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
- `Section\ChooseClassofContraction.xml`: RD `BrowseTableOfLimit_RD` (Tahun **KOSONG** — `<pyValue/>` / `<Tahun/>`, filter
  dibuang; Bizcode = `QuotationData.BusinessCode`, Category = baris); kolom Description (Limit tersembunyi); Choose →
  `SetDataClassofConstraction` (Description, PctLimit) + `GetLowestPctLimit_ACT`. ⚠️ Ralat 03-10-2026 (sesi c3,
  dikonfirmasi sesi 0f): tulisan awal "Tahun = `pyWorkPage.OfferFacIn.CurrentYear`" milik **autocomplete** Class of
  Construction di `Section\OccupationItemFacIn_Section.xml` (data page `D_BrowseTableOfLimit`), bukan popup ini.
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
- [x] Backend: simpan / baca okupasi (T_OCCUPATIONLIST + T_TABLEOFLIMIT rancangan, migrasi 189), `kdRiskExposure`, cari kosong.
- [x] Backend: endpoint table-of-limit (BIZCODE butir 89; TAHUN = tahun Begin date, A153).

## Backend (sesi c3, 03-10-2026) — disusun agent

- `GET`/`PUT …/objek`: tiap baris objek membawa `occupations` (kunci persis `OkupasiObjek`; selalu larik). Disimpan ke
  **tabel rancangan** `T_OCCUPATIONLIST` (jalur `LocationList/Property/OccupationList`, induk T_PROPERTY — bukan jalur
  RiskLocation) + `T_TABLEOFLIMIT` (satu per okupasi): `category` / `constructionClass` / `pctLimit` = `.TableOfLimit.
  Category / Description / PctLimit`. Urut `SEQ_NO`. Dihapus-sisip ulang bersama objek (TableOfLimit → okupasi → …).
- `pctLimit` **teks apa adanya** (keputusan work owner butir 68.1 — rancangan NUMBER, isi teks; koma dan spasi ujung
  tidak diubah). Validasi PUT: lebar kolom saja (400 `baris[n].occupations[m].<medan>`); tanpa enumerasi; Class of
  Construction tidak diwajibkan (L-2).
- `GET /api/nbfacin/occupation`: tiap baris membawa `kdRiskExposure` (`OCCUPATION.KDRISKEXPOSURE`, DDL VARCHAR2(1000));
  `cari` kosong → seluruh okupasi FIRE urut NAME, OLDID, ≤ 500 (popup Choose Occupation); 1 karakter tetap 400.
- Migrasi **189** (`189_t_occupationlist.sql` + `_down`). Tanpa kolom baru → jumlah kolom loader tetap **1.418**; dua
  pelebaran di `amandemenLebar`. ⛔ Ditulis, **tidak dijalankan** agent.

**Belum dibangun — `GET …/table-of-limit`** (perintah sesi 0f 03-10-2026: *"Silakan lanjut tiket 40 tanpa endpoint
table-of-limit dulu"*; `cariTableOfLimit` di frontend mendapat 404 sampai endpoint ini ada):
- `[terverifikasi]` RD `BrowseTableOfLimit_RD`: kelas `ASM-FW-GISFW-Int-TABLEOFLIMIT`; filter `.Tahun = Param.Tahun`,
  `.Bizcode = Param.Bizcode`, `.Category = Param.Category`, `.ID = Param.ID` (A AND B AND C AND D); DISTINCT; maks 500;
  urut Category lalu Description.
- `[terverifikasi]` `SetValidateDate_Act`: `.CurrentYear = @DateTime.FormatDateTime(.PolicyData.StartDateTime,"yyyy",
  "Indonesia/Jakarta","in_ID")` — dapat diturunkan dari `T_GENERAL_POLIS.START_DATE_TIME`.
- `[terverifikasi]` `OfferFacIn.QuotationData` diisi **satu halaman utuh**: `Activity\SetCedingCo_Act.xml` (baris 458–459)
  `pyWorkPage.OfferFacIn.QuotationData = pyWorkPage.Quotation` — jadi BusinessCode = `pyWorkPage.Quotation.BusinessCode`
  `[dugaan]` (bergantung pada perilaku salin-halaman Pega). ⚠️ **Ralat** (temuan code review sumbu spec): tulisan awal
  bab ini menyebut "tidak ada yang mengisinya" — dua cara pencarian saya sama-sama hanya mencari jalur berakhiran
  `.BusinessCode`, sehingga salin-halaman itu luput.
- `[pertanyaan terbuka]` Asal `Quotation.BusinessCode`: setter daun hanya `SetBusinessType_Act` (mengosongkan, `""`) dan
  `CopyToPolicyListPASSG`; nol section NB FacIn yang mengikat `.BusinessCode`. Kolom `T_QUOTATIONDATA.BUSINESS_CODE` juga
  belum dibuat (183 sebagian).
- `[pertanyaan terbuka]` Nama tabel `TABLEOFLIMIT` hanya dari kelas + berkas contoh `[dugaan]` (tidak ada SQL korpus);
  tipe TAHUN / BIZCODE / PCTLIMIT tidak diketahui (contoh "70,000 ") — DDL diminta lewat sesi 0f.

**Keputusan agent — DISETUJUI work owner 03-10-2026** (butir 82, diteruskan sesi `nusantarare-0f`: *"setuju keputusan
agent"*):

| # | Keputusan | Dasar |
| --- | --- | --- |
| A138 | `OCCUPATION_ID` / `OCCUPATION_NAME` `VARCHAR2(1000)` (rancangan 50 / 500) | diisi dari `OCCUPATION.OLDID` / `NAME` VARCHAR2(1000) — pola butir 80 "Widen joined columns", sama dengan tiket 38 |
| A139 | `T_TABLEOFLIMIT.PCT_LIMIT` `VARCHAR2(50)` | butir 68.1 menetapkan teks tetapi menahan lebarnya (tiket 23); contoh 7 bita; melebarkan nanti tidak merusak |
| A140 | `T_OCCUPATIONLIST.PARENT_ID` tanpa FK; indeks (PARENT_TABLE, PARENT_ID) | rancangan berinduk jamak (tiga jalur) — FK ke satu induk akan menolak dua jalur lainnya |
| A141 | `cari` 1 karakter tetap 400; kosong = semua | kontrak tiket 40 hanya membuka "kosong"; saran Surrounding Risk tetap ≥ 2 |

⚠️ **Risiko tercatat:** simpan objek menghapus `T_RISKLOCATION` dan okupasi berinduk `T_PROPERTY` saja. Okupasi berinduk
`T_RISKLOCATION` (jalur `RiskLocation/OccupationList`, hanya ditulis loader) akan **yatim** tanpa galat karena tanpa FK
(A140). Belum terjadi: loader belum menulis ke Oracle; ditangani bersama tiket 23 / pemuat.

**Keputusan work owner W-5 (03-10-2026, diteruskan sesi 0f; butir 86):** BusinessCode table-of-limit = **kode Group
Business** (dipilih di Create opportunity). Kolom mana yang cocok dengan `TABLEOFLIMIT.BIZCODE` dibuktikan dari data contoh
saat endpoint dibangun — belum dikerjakan. DDL `TABLEOFLIMIT.txt` (03-10-2026): seluruh kolom VARCHAR2(4000 BYTE) —
PCTLIMIT dan TAHUN teks.

## Backend table-of-limit (sesi c3, 03-10-2026) — disusun agent

- `GET /api/nbfacin/kasus/{caseId}/table-of-limit?category=` → `{"baris":[{"description","pctLimit"}]}`. Tanpa identitas;
  404 case tidak ada; 400 `category` > 50 bita; **409** bila BIZCODE tidak dapat ditentukan; 503.
- BIZCODE (butir 89): `BUSINESS.ID` ber-`NOTE` = `T_NB_OPPORTUNITY.CLASS_OF_BUSINESS` dan `BUSINESSGROUPID` =
  `GROUP_BUSINESS_ID` case. Class of Business / Group Business kosong, 0 baris, atau > 1 baris → 409 dengan pesan jelas.
- TABLEOFLIMIT (DDL `TABLEOFLIMIT.txt`, semua VARCHAR2(4000)): `BIZCODE = kode AND TAHUN = tahun Begin date (WIB)`
  [+ `CATEGORY = category` bila diisi]; Begin date kosong → 409;
  DISTINCT atas kolom laporan RD (Bizcode, Category, Description, PctLimit, Note); urut Category, Description; ≤ 500.
  `pctLimit` teks apa adanya.

**`Tahun` — dua jalan di korpus `[terverifikasi]`:** layar Occupation punya DUA jalan ke RD yang sama:
(1) **autocomplete** Class of Construction di `Section\OccupationItemFacIn_Section.xml` (data page `D_BrowseTableOfLimit`)
mengirim `Tahun = pyWorkPage.OfferFacIn.CurrentYear` (2 kemunculan; `_IsUW` 2) — CurrentYear = tahun Begin date
(`SetValidateDate_Act`); (2) **tombol** Choose Class of Construction (`FlowAction\ChooseClassofContraction` →
`Section\ChooseClassofContraction.xml`) mengirim Tahun KOSONG (`<pyValue/>`, `<Tahun/>`) → filter dibuang (peringatan Pega
"the filter will be dropped entirely", `BrowseOccupationFacInFIRE_RD.xml` baris 1867). ⚠️ Ralat: tulisan awal bab ini
menyebut CurrentYear "hanya di `Section\TableOfLimit.xml`" — salah; pencarian pertama hanya mencocokkan bentuk tag
`<Tahun>` dan melewatkan bentuk `pyName`/`pyValue` (temuan code review sumbu spec).

**Keputusan agent (menunggu konfirmasi):**

| # | Keputusan | Dasar |
| --- | --- | --- |
| A153 | Endpoint **menyaring TAHUN = tahun Begin date** (WIB) — jalan autocomplete; Begin kosong → 409 | ✅ **DISETUJUI work owner 03-10-2026** (butir 92: *SARING tahun Begin date*); korpus punya dua jalan (autocomplete bertahun, tombol tanpa tahun) |
| A154 | NOTE dicocokkan persis; 0 / > 1 baris BUSINESS → 409 | ✅ **DISETUJUI work owner 03-10-2026** (butir 93: *"setuju A154"*); arahan sesi 0f; nama tersimpan = BUSINESS.NOTE pilihan datalist (cobSah) |
