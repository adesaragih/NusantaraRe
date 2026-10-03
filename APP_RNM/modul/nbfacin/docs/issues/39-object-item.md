# 39: Sub-tab Object Item — grid item objek, form detail, total TSI per mata uang

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026: *"bisa
> lanjut ke smua ga? atau satu2 berurutan, karena itu semua harus dibuat"* (urutan: Object Item, Occupation, FEA,
> Loss Record, Loss Record Internal). Belum ada tangkapan layar; tampilan diturunkan dari XML.

**What to build:** sub-tab **Object Item** di baris objek FIRE. Isinya grid item (Object Item Type · Condition · Year ·
Unit(s) · Currency · TSI Object Item) dengan Tambah / Hapus. Baris yang dibuka menampilkan form detail. Di bawahnya
grid Total TSI per mata uang. Datanya ikut Save tab Object.

**Blocked by:** daftar Condition dan Adjustment Pct. (PctAdjust2) — aturan properti Pega tidak ada di korpus.

**Status:** frontend selesai 03-10-2026 (uji hijau); backend selesai 03-10-2026 (sesi c3; migrasi 188 ditulis, BELUM dijalankan).

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\Property.xml`: tab "Object Item" meng-include `PropertyItemList`.
- `Section\PropertyItemList.xml`:
  - Layout NB `!IsEDM`. Grid `.Property.PropertyItemList` (kelas `Data-PropertyItem`, master-detail, flow action
    `PropertyItemFacIn_FlowAction`); kolom sel 19–24. Add = ikon → addRow + `SumTSIObjItem_act` + `SetPropertyItemNo`.
    Delete = ikon → deleteRow + `UpdatePropertyItemNo`.
  - Grid Total `.Property.TotalTSIList` (Currency · Total TSI, sel 59–60, 4 desimal).
- `Section\PropertyItemFacIn_Section.xml`: Object Item Type (RD `BrowseV_JN_OBJ_ITEM`, nilai MJOI_KODE, "Choose") ·
  Object Item Note · Year · Unit(s) · Condition (associated) · Currency (RD `BrowseCurrency_RD`, Currency != "ITL") ·
  TSI Object Item (All Unit) · Year of Planting · No of Trees · Area ( Hectar ) · Remark of `<ItemType>` · centang
  Adjustable · Adjustment Pct. % (PctAdjust2 dropdown bila `!IsAdjustableFlag`; PctAdjustOther angka bila
  `IsAdjustableFlag`).
- `Activity\SetObjItemType_Act.xml`: RD `GetObjectItem` (MJOI_KODE) → `.ItemType = JN_OBJ_ITEM`,
  `.PropertiItemNote = KETERANGAN`.
- `Activity\AddPaymentCurency_ACT.xml`: `TotalTSIList(Name = CURR, TSI = Σ TSIObjectItem where Currency = CURR)` bila
  > 0; tanpa konversi kurs, tanpa pembulatan.
- `Activity\ValidateAdjustPct.xml`: "%Adjustment can't be less than 60% or more than 100%".
- `Activity\ResetPct_Adjustment.xml`: Adjustable "false" → `PctAdjustOther = 100`.
- `SetErrorMessageUnit_Act` (Unit <= 0) dan `SetErrorMessageTSIObjectItem_Act` (TSI < 0) memakai field value
  `ErrorMessageUnit` / `TSIObjectItemErrorMessage`, yang tidak ada di korpus.
- Rancangan flat: `T_PROPERTYITEMLIST` (induk T_PROPERTY) — CURRENCY, ITEM_TYPE(_ID), PROPERTI_ITEM_NOTE, REMARK,
  IS_ADJUSTABLE_FLAG, PCT_ADJUST2, PCT_ADJUST_OTHER, PROPERTY_ITEM_NO, TSI_OBJECT_ITEM (uang), dan lain-lain.

## Keputusan agent

- **K-1** Hitung ulang premi / spreading saat TSI berubah (`CountPremi_ACT`, `cekSpreadingFactIn`) = tab Coverage.
- **K-2** Urutan baris Total = urutan munculnya mata uang di daftar item (Pega: urutan tabel CURRENCY).
- **K-3** Pesan Unit <= 0 dan TSI minus = teks sistem baru (field value Pega tidak ada di korpus). TSI harus desimal
  bertitik.
- **K-4** Daftar Condition dan PctAdjust2 menunggu aturan properti dari work owner.
- **K-5** TSI = uang, teks desimal sampai ke layar. Total dihitung `jumlahDesimal` (BigInt, ADR-0003/0016);
  tampil 4 desimal gaya Indonesia (`formatNumber`).
- **K-6** Tambah / Hapus = teks seperti grid objek (di Pega berupa ikon). PropertyItemNo = urutan baris.

## Kontrak

- `ObjekFire.items: ItemObjek[]` (lihat `frontend/api.ts`), ikut `GET`/`PUT …/objek`.
- `GET /api/nbfacin/jenis-item-objek` → `{ baris: [{ kode, nama, keterangan }] }`.
- `GET /api/nbfacin/mata-uang` → `{ baris: string[] }`.

## Acceptance criteria

- [x] Grid item + Tambah / Hapus + form detail; label diuji ke korpus.
- [x] Object Item Type mengisi nama + Note; Adjustable mengatur Adjustment Pct.; galat menahan Save.
- [x] Total TSI per mata uang eksak (uji tanpa float).
- [x] Backend: simpan / baca item (T_PROPERTYITEMLIST, migrasi 188), dua endpoint pilihan.
- [ ] Daftar Condition dan PctAdjust2.

## Backend (sesi c3, 03-10-2026) — disusun agent

- `GET`/`PUT …/objek`: tiap baris objek membawa `items` (kunci persis `ItemObjek`; selalu larik). Disimpan ke
  `T_PROPERTYITEMLIST` (banyak per property, `SEQ_NO` = `PROPERTY_ITEM_NO` = 1..n), dihapus-sisip ulang bersama objek
  (item sebelum property). Dibaca dengan kueri kedua per case, dipasangkan lewat `T_PROPERTY.ID`.
- **Uang/persen (ADR-0003/0016):** `tsi`, `pctAdjust2`, `pctAdjustOther` melintas JSON sebagai teks desimal bertitik
  (angka JSON → 400). services memeriksa bentuknya dan membandingkan batas dengan `apd` (nol float). Repository menulis
  `TO_NUMBER(:n, 'FM9…9D99999999', <NLS titik>)` ke `NUMBER(38,8)` dan membaca `TO_CHAR TM9` ber-NLS
  titik, diurai `apd` lalu ditulis ulang kanonik (Oracle TM9 memberi `.5` → `0.5`).
- **Validasi PUT** (400, `baris[n].items[m].<medan>`):
  - `currency` WAJIB dan harus ada di `CURRENCY` tanpa ITL (A133);
  - `tsi` kosong atau `^[0-9]{1,30}(\.[0-9]{1,8})?$`;
  - `unit` kosong atau bulat > 0 (`SetErrorMessageUnit_Act` `.Unit<=0` `[terverifikasi]`);
  - `pctAdjust2` / `pctAdjustOther` kosong atau desimal yang sama; bila `isAdjustable` → `pctAdjustOther` wajib dan 60..100
    (`ValidateAdjustPct` `Local.adjpct<60||Local.adjpct>100` `[terverifikasi]`);
  - lebar kolom; Condition / PctAdjust2 tanpa enumerasi.
- `GET /api/nbfacin/jenis-item-objek` → `{"baris":[{"kode","nama","keterangan"}]}`;
  `GET /api/nbfacin/mata-uang` → `{"baris":["…"]}`. Tanpa identitas; 503 tanpa basis data.
- Migrasi **188** (`188_t_propertyitemlist.sql` + `_down`); amandemen loader `amandemenItem` (skema 79 tabel / **1.418**
  kolom). ⛔ Ditulis, **tidak dijalankan** agent.

**Bukti `[terverifikasi]`:**
- `BrowseV_JN_OBJ_ITEM`: logika `A AND B AND C`; A `.JN_OBJ_ITEM = Param`, B `.ISACTIVE = 1`, C `.KELOMPOK = Param` — dropdown
  sel mengirim parameter kosong (`pyReportDefParams` MJOI_KODE/JN_OBJ_ITEM tanpa nilai), filter berparameter kosong dibuang
  (pola tiket 38); DISTINCT; urut JN_OBJ_ITEM ASC; maks 10000.
- Asal KETERANGAN: `SetObjItemType_Act` memanggil RD `GetObjectItem` (kelas V_JN_OBJ_ITEM yang sama) dengan
  `Param.JN_OBJ_ITEM = Param.OBJ_ITEM`, filter `.MJOI_KODE StartsWith`, lalu `.PropertiItemNote =
  DataObjItem.pxResults(1).KETERANGAN`, `.ItemType = pxResults(1).JN_OBJ_ITEM` — jadi KETERANGAN adalah **kolom view
  V_JN_OBJ_ITEM**, bukan rumus.
- `BrowseCurrency_RD`: `.Currency = Param` (dibuang bila kosong), `.ID = Param` (dibuang), `.Currency != "ITL"`; maks 500;
  tanpa urutan.
- Kontrol: PropertyYear / Unit / TSIObjectItem / PctAdjustOther = pxNumber; Year / NoOfTree / AreaHectar = pxTextInput;
  Note / Remark = pxTextArea; Condition = pxDropdown.

**Keputusan agent (menunggu konfirmasi):**

| # | Keputusan | Dasar |
| --- | --- | --- |
| A132 | Enam kolom baru `T_PROPERTYITEMLIST`: `PROPERTY_YEAR`, `UNIT`, `YEAR`, `NO_OF_TREE`, `AREA_HECTAR` `VARCHAR2(50)`, `CONDITION` `VARCHAR2(500)`; angka disimpan teks | medan ada di layar, tidak di rancangan (pola A110); nama/tipe = medan bernama sama di tabel rancangan lain (`YEAR`/`UNIT` 50, `CONDITION` 500); angka teks pola A129 |
| A133 | `currency` **wajib** per item dan harus ada di `CURRENCY` (≠ ITL) | rancangan `CURRENCY … DEFAULT 'UNKNOWN' NOT NULL` (K-069) + K-012 (ADR modul 0006 `Unknown` eksplisit — bukan `docs/bersama` ADR-0006; aplikasi tidak menulis `UNKNOWN`); Pega tidak mewajibkan (pyRequired false). ⚠️ Baris lama bermata uang `UNKNOWN` akan ditolak saat disimpan ulang sampai mata uangnya dipilih |
| A134 | `keterangan` = KETERANGAN **baris yang sama** di daftar jenis item | Pega mengambil `pxResults(1)` dari `MJOI_KODE StartsWith` (urut JN_OBJ_ITEM, tanpa ISACTIVE) — sama hasilnya bila tidak ada kode yang menjadi awalan kode lain `[dugaan]`; data view tidak ada di korpus |
| A135 | Mata uang urut `CURRENCY` | RD tanpa urutan (urutan Pega tidak terdefinisi) |
| A136 | `tsi`/pct > 8 desimal → **400** (bukan dibulatkan) | ADR-0016 membulatkan saat MEMUAT data lama; isian layar ditolak supaya nilai yang disimpan = yang diketik |
| A137 | Tidak Adjustable → `pctAdjustOther` tidak dibatasi 60..100; `pyMin`/`pyMax` kontrol (PropertyYear 4, PctAdjustOther 2/3) tidak diperlakukan sebagai aturan nilai | `ValidateAdjustPct` hanya untuk Adjustable; pyMin/pyMax di blok `Embed-Control-Mode` diduga panjang tampilan `[dugaan]` |

Tidak diport: `.Property.TotalTSIList` (dihitung frontend, K-5; rancangan tidak punya tabelnya), `.PctAdjust1` (ada di
section, tidak di kontrak), `.SubLimit` (dropdown `FacInKurs.pxResults`, kelas GrandTSI — ada di section baris 9108/9388,
tidak di kontrak dan tidak di rancangan T_PROPERTYITEMLIST), hitung ulang premi (K-1).
