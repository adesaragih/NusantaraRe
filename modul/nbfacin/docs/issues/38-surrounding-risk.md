# 38: Sub-tab Surrounding Risk — empat sisi, Other Description, saran Occupation

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026: *"kamu
> lanjut ke tab surroundingrisk"*. Tangkapan layar tab dan lima aturan properti dropdown menyusul dari work owner
> pada hari yang sama (tangkapan tidak disalin).

**What to build:** sub-tab **Surrounding Risk** di baris objek FIRE (tab Object, tiket 35). Isinya empat blok sisi
(Front / Left / Back / Right) dan blok Other Description. Datanya ikut Save tab Object.

**Blocked by:** — (aturan properti dropdown sudah ditambahkan work owner 03-10-2026).

**Status:** frontend selesai 03-10-2026 (uji hijau; lima dropdown dari aturan properti); backend selesai 03-10-2026 (sesi c3; migrasi 187 ditulis, BELUM dijalankan).

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\Property.xml`: sub-tab "Surrounding Risk" meng-include `RiskAround`.
- `Section\RiskAround.xml`:
  - Blok `Front` / `Left` / `Back` / `Right` (pyTitle). Tiap blok berisi Occupation (sel 9/23/37/51), Construction
    (12/26/40/54), Distance (meter) (13/27/41/55), dan Note (14/28/42/56), atas
    `.Property.SurroundingRisk.{Front|Left|Back|Right}{Occupation|Construction|Distance|Note}`.
  - Other Description: 65 Ownership (`.Property.Ownership`) · 66 House keeping Status · 67 Flood Area Status · 68 Flood
    Area (visible `FloodAreaStatus==0`) · 69 Housekeeping Remark · 77–79 centang Production Process / Job With a Chance of
    Fire / Flammable Item (`.Property.Is{ProductionProcess|HotWorkProcess|FlammableItem}Flag`).
  - Tidak ada medan wajib.
- Occupation: autocomplete data page `D_BrowseOccupationFacInFIRE` → RD `BrowseOccupationFacInFIRE_RD` (TYPE "FIRE";
  urut Name, OldID; maks 500). Nilai = `.OldID`; memilih mengisi Note sisi itu dengan `.Name`. DDL `OCCUPATION.txt`.
- Construction, Ownership, House keeping Status, Flood Area Status, Flood Area: `pyListSource=associated`. Aturan
  propertinya TIDAK ada di korpus NB FacIn; work owner menambahkannya 03-10-2026 di `D:\migrasi\RNM\DDL\`
  (`FrontConstruction.xml`, `Ownership.xml`, `HousekeepingStatus.xml`, `FloodAreaStatus.xml`, `FloodArea.xml`;
  PromptList, pasangan `pyStandardValue`/`pyLocalizedValue` per rowdata, diuji `labels.test.ts`):
  Construction "Silahkan pilih" / I / II / III (nilai teks panjang, terpanjang 215 bita); Ownership 2 Not Informed /
  0 Own / 1 Rent; HousekeepingStatus 0 Not Informed / 1 Good / 2 Fair / 3 Poor; FloodAreaStatus 2 Not Informed /
  0 Yes / 1 No; FloodArea "Silahkan pilih" / 1 Low / 2 Medium / 3 High / 4 Very High. Data contoh (Ownership "2",
  HousekeepingStatus "0", FloodAreaStatus "2") = "Not Informed".
- `Activity\NegativeIsNotAllowed.xml`: Distance < 0 → "Jarak Resiko Sekitar tidak boleh Minus". Distance = pxNumber,
  min 0, maks 100000, 2 desimal.
- Rancangan flat: `T_SURROUNDINGRISK` (induk T_PROPERTY) hanya berisi FLOOD_AREA_STATUS dan HOUSEKEEPING_STATUS.
  `T_PROPERTY` di rancangan punya OWNERSHIP dan IS_*_FLAG, tetapi migrasi 186 membuatnya sebagian (A109).

## Keputusan agent

- **J-1** Paragraf `InputFireObject_Q1_RiskFactorOption` (isi tidak ada di korpus): teks diambil dari tangkapan layar
  work owner (`TEKS_SEKITAR.faktorRisiko`).
- **J-2** Saran Occupation muncul setelah 2 karakter, jeda 400 ms; isian bebas tetap boleh (Pega mengizinkan).
- **J-3** Ownership / House keeping / Flood Area Status tanpa baris kosong -> objek baru berawal pilihan pertama
  "Not Informed". Left / Back / Right Construction memakai daftar FrontConstruction `[dugaan]` (hanya aturan Front
  yang dikirim; tangkapan layar menampilkan "Silahkan pilih" di keempat sisi).
- **J-6** Tata letak = tangkapan layar: dua blok "Other Description" berdampingan (kiri isian, kanan paragraf + tiga
  centang).
- **J-4** Note = tampil-saja (Pega: baca-saja), diisi dari Occupation.
- **J-5** Distance minus menahan Save tab Object dan membuka sub-tab Surrounding Risk baris itu.

## Kontrak

- `ObjekFire` (lihat `frontend/api.ts`) bertambah: `ownership`, `isProductionProcess`, `isHotWorkProcess`,
  `isFlammableItem`, dan `surroundingRisk { front|left|back|right: { occupation, construction, distance, note },
  housekeepingStatus, floodAreaStatus, floodArea, housekeepingRemark }`. Ikut `GET`/`PUT …/objek`.
- `GET /api/nbfacin/occupation?cari=` → `{ baris: [{ oldId, name }] }`.

## Acceptance criteria

- [x] Empat sisi + Other Description tampil; Flood Area bersyarat; label diuji ke RiskAround.xml.
- [x] Memilih Occupation mengisi Note; Distance minus → pesan Pega, Save tertahan.
- [x] Backend: simpan/baca medan baru (migrasi 187), endpoint Occupation.
- [x] Daftar lima dropdown (aturan properti `DDL\*.xml`, diuji).

## Backend (sesi c3, 03-10-2026) — disusun agent

- `GET`/`PUT /api/nbfacin/kasus/{caseId}/objek` membawa `ownership`, `isProductionProcess`, `isHotWorkProcess`,
  `isFlammableItem`, `surroundingRisk` persis kontrak `api.ts` (pemetaan medan eksplisit di handler). `surroundingRisk`
  selalu objek utuh (sisi kosong = teks kosong).
- Simpan: `T_PROPERTY` + `OWNERSHIP` / `IS_*_FLAG` (boolean teks `true`/`false`, A115); satu baris `T_SURROUNDINGRISK` per
  property, dihapus-sisip ulang bersama baris objek (anak sebelum induk). Baca: `LEFT JOIN` (`UNIQUE PARENT_ID`).
- Validasi PUT (400, pesan menyebut indeks dan sisi): `baris[n].surroundingRisk.<sisi>.distance` kosong atau angka
  `^[0-9]+(\.[0-9]{1,2})?$` ≤ 100000 — `[terverifikasi]` keempat kontrol Distance `Section\RiskAround.xml`: pyMin 0,
  pyMax 100000, pyDecimalPlaces 2; lebar kolom (bita) tiap medan teks. Construction / Ownership / Housekeeping / Flood
  **tidak** dicocokkan ke daftar.
- `GET /api/nbfacin/occupation?cari=` → `{"baris":[{"oldId","name"}]}`; 400 bila `cari` (dipangkas) < 2 atau > 255
  karakter; 503; 500. Tanpa identitas (pola lookup). SQL: `TYPE = 'FIRE' AND (UPPER(OLDID) LIKE :p OR UPPER(NAME) LIKE :p)`,
  urut `NAME, OLDID`, ≤ 500.
- Migrasi **187** (`187_t_surroundingrisk.sql` + `_down`): `ALTER TABLE T_PROPERTY ADD` empat kolom rancangan; `CREATE TABLE
  T_SURROUNDINGRISK` utuh + 18 kolom baru; `SEQ_T_SURROUNDINGRISK`. Amandemen loader `amandemenSekitar` (skema 79 tabel /
  **1.412** kolom). ⛔ Ditulis, **tidak dijalankan** agent.

**Bukti RD Occupation `[terverifikasi]`:** `ReportDefinition\BrowseOccupationFacInFIRE_RD.xml` logika `A AND B AND D AND E AND
(F OR C)`; RD itu sendiri memuat peringatan Pega *"the filter will be dropped entirely"* untuk filter berparameter kosong —
sel 9 hanya mengirim `TYPE`, jadi RD memuat SEMUA okupasi FIRE (urut Name, OldID; maks 500) lalu autocomplete mencari
`.OldID`/`.Name` di dalamnya. Lebar Construction diukur dari `DDL\FrontConstruction.xml` (3 nilai standar, terpanjang
215 bita; dua cara: tag `pyStandardValue` dan cacah baris = 3).

**Keputusan agent — DISETUJUI work owner 03-10-2026** (butir 82, diteruskan sesi `nusantarare-0f`: *"setuju keputusan
agent"*):

| # | Keputusan | Dasar |
| --- | --- | --- |
| A129 | `*_DISTANCE` disimpan **teks** `VARCHAR2(50)`, bukan NUMBER | pola `NUMBER_OF_FLOOR` rancangan; menghindari konversi desimal NLS; bukan uang |
| A130 | 18 kolom baru `T_SURROUNDINGRISK` (bukan di rancangan): Distance/FloodArea `VARCHAR2(50)`, Construction/HousekeepingRemark `VARCHAR2(500)`, Occupation/Note `VARCHAR2(1000)` | medan ada di layar, tidak di rancangan (pola A110); Construction 215 bita terukur; Occupation/Note diisi dari `OCCUPATION.OLDID`/`NAME` VARCHAR2(1000) — lebar sumber (pola butir 80), panjang data nyata tidak ada di korpus |
| A131 | Occupation dicari **mengandung**, **tidak peka huruf**, di OLDID ATAU NAME, langsung di Oracle (≤ 500 hasil) | mode cocok autocomplete Pega tidak tertulis di korpus `belum terverifikasi`; Pega memotong ke 500 baris PERTAMA lalu mencari di dalamnya — port mencari di seluruh tabel (okupasi di luar 500 pertama tetap bisa ditemukan) |
