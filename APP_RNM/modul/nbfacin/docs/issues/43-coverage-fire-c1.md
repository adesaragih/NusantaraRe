# 43: Tab Coverage FIRE — tahap C1 (objek → item → coverage, hitung premi basis 1–4)

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026: *"lanjut
> coverage"*. Belum ada tangkapan layar; tampilan diturunkan dari XML. Tahapan C1–C4 diusulkan sesi 0f.

**What to build:** tab **Coverage** (kasus FIRE), bertingkat:
- grid objek (No. · Object Name · Location);
- grid item (Object Item Type · Currency · TSI Object Item · Total Gross Premium · ‰ Total Net Rate) dan grid total
  per mata uang;
- grid coverage (Coverage · ‰ Standard Rate · Premi; Tambah = lima coverage FIRE otomatis bila kosong; Hapus);
- form coverage: Coverage Basis 1–4, Choose Coverage, rate, persen, Gross Premium.

Premi dihitung backend (`CountPremi_ACT`). Save = `PUT …/objek`.

**Blocked by:** — (DDL `COVERAGE_FACIN` / `COVERAGE` diberikan work owner 03-10-2026).

**Status:** frontend selesai 03-10-2026 (uji hijau); backend selesai 03-10-2026 sesi c3 (uji hijau; migrasi 193 ditulis,
**belum dijalankan**).

⛔ **Migrasi 193 WAJIB dijalankan di DEV SEBELUM backend baru dipasang** — urutan … → 192 → **193**. `repository/item.go`
kini membaca / menulis `T_PROPERTYITEMLIST.TOTAL_GROSS_PREMI` / `TOTAL_NET_RATE` dan `T_COVERAGELIST`; tanpa 193, `GET` /
`PUT …/objek` gagal ORA-00904 dan **seluruh tab Object ikut rusak**, bukan hanya tab Coverage.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\InputInwardFacultativeDtl.xml`: tab "Coverage" (visible IsFire) → `CoverageList` (objek, IsFire) /
  `InputDtlCoverage_FacIn` (!IsFire) / `SummaryCoverage_Section`; tombol Save (`SaveFacIn_Act`).
- `Section\CoverageList.xml` (sel 37–39) → flow action `CoverageListFire_FlowAction` → `PropertyItemListCoverage.xml` (grid
  item sel 17–21; grid `.Property.TotalTSIPremiGrossList` sel 51–54) → flow action `PropertyItemCoverageFacIn_FlowAction` →
  `InputCoverageFire.xml` (grid `.CoverageList` sel 16–18) → flow action `CoverageItem` → `CoverageItem.xml`.
- `CoverageItem.xml`: label sel 4–47 diuji `labels.test.ts`; ‰ Gross Rate wajib; syarat tampil per `.CoverageBasis`
  (2 = First Loss/First Scale/Indemnity/TSILiability; 3 = EML/PML; 4/5 = Sub Limit; 5 = Layering grid).
- `DDL\CoverageBasis.xml`: 1 Sum Insured / 2 First Loss / 3 EML / PML / 4 Sub Limit / 5 Layering Basis.
- `Activity\AddCoverageAutoFire.xml`: lima coverage, urutan korpus 100815, 100828, 100829, 100825, 100840 (basis 1,
  %Indemnity 100) lalu `CountPremi_ACT`.
- `Activity\CountPremi_ACT.xml` (57 langkah): Prorate dari periode polis; `.TSI` = TSIObjectItem item; Sum Insured percent:
  `.Premium=@Math.divide((.TSI*.Rate*Local.Prorate*.IndemnityPercentage*Local.LossLimit),1000000000,20)`; amount:
  `.Rate=@Math.divide((.Premium * 1000000000),(.TSI *Local.Prorate* .IndemnityPercentage * Local.LossLimit), 20)`; varian
  adjustable (÷1e11), basis 2 (×FirstScale ÷1e11; TSILiability = TSI×FirstLoss/100), basis 3 (EmlPml), basis 4/5
  (SubLimit); diskon; NetRate menimpa Rate bila terisi; total `PropertyItemList.TotalGrossPremi = Σ Premium`.
- `Activity\SumTotalTSIPremiGross_Act.xml` (cabang FIRE): `TotalTSIPremiGrossList` per mata uang, Rate = Premium/TSI×1000.
- Rancangan flat: `LocationList/Property/PropertyItemList/CoverageList` → `T_COVERAGELIST`.

## Keputusan agent

- **P-1** Layering (basis 5), Zone / 4.2 Construction, Accumulation, Indemnity Unit, View Indemnity Table, Deductible,
  lookup rate standar (`SetRatePolis_ACT`), Copy Coverage / Deductible, Summary = tahap C2–C4.
- **P-2** Pilihan Days = 365 / 366 `[dugaan]` (aturan `.Day` tidak ada; data contoh 365).
- **P-3** ‰ Gross Rate bertanda wajib tanpa menahan Save (pola L-2).
- **P-4** Hitung premi = `POST …/hitung-coverage` dengan jeda 500 ms; jawaban lama dibuang; premi tidak dihitung di layar
  (tanpa float).
- **P-5** Tab Coverage dan tab Object disimpan masing-masing lewat tombol Save tab itu.
- **P-6** "Total Gross Rate" ditampilkan tanpa akhiran "%" (nilainya permil; Pega memberi "%").

## Kontrak

- `ItemObjek.coverages?: CoverageObjek[]`, `totalGrossPremi?` dan `totalNetRate?` (server); `ObjekFire.totalPerCurrency?:
  TotalCoverage[]` (server) — ikut `GET`/`PUT …/objek`.
- `GET /api/nbfacin/coverage?cari=` → `{ baris: [{ id, oldId, nama }] }` (COVERAGE_FACIN, Type FIRE; ~~aktif~~ — view
  tanpa ACTIVESTATUS, A160).
- `GET /api/nbfacin/coverage-otomatis` → lima coverage `AddCoverageAutoFire`.
- `POST /api/nbfacin/kasus/{caseId}/hitung-coverage` badan `{ coverage, tsi, mode: "percent" | "amount", modeDiskon?,
  isAdjustable?, pctAdjustOther? }` → `CoverageObjek` terhitung.

## Acceptance criteria

- [x] Grid berantai + Tambah (otomatis) / Hapus + form basis 1–4; label diuji ke korpus; tanpa float di layar.
- [x] Backend: hitung `CountPremi_ACT` (basis 1–4, prorate, diskon, net rate), simpan coverage, total item / mata uang.
- [x] DDL COVERAGE_FACIN (03-10-2026).
- [ ] Migrasi 193 dijalankan di DEV (work owner).

## Backend (sesi c3, 03-10-2026)

**Rute** (`handlers/coverage.go`, `handlers/objek.go`):
- `POST …/kasus/{caseId}/hitung-coverage` — tanpa menyimpan, tanpa identitas (pola lookup). Medan server `tsi`,
  `tsiLiability`, `proRatePercent` dari badan diabaikan. 400: desimal tak sah (jalur `coverage.<medan>` / `tsi` /
  `pctAdjustOther`), basis 5 ("Layering belum didukung"), basis bukan 1–4, `mode` / `modeDiskon` di luar
  `percent` / `amount`; 404 case tidak ada; **409** Begin / End date case kosong (Prorate tidak dapat dihitung); 503.
- `GET …/objek` / `PUT …/objek`: `items[].coverages[]` (selalu larik ke luar), `items[].totalGrossPremi` (BACA-SAJA, dihitung
  ulang saat PUT), `items[].totalNetRate` (disimpan apa adanya), `totalPerCurrency[]` per objek (BACA-SAJA, dihitung saat
  GET). Saat PUT seluruh coverage DIHITUNG ULANG server; 400 ber-indeks `baris[n].items[m].coverages[k].<medan>`; item
  ber-coverage dengan mata uang > 10 byte → 400 (`CURRENCY_CODE` VARCHAR2(10)); Begin / End kosong dan ada coverage → 409.
- `GET /api/nbfacin/coverage?cari=` → `{baris:[{id,oldId,nama}]}`; `cari` dipangkas, > 255 karakter → 400, kosong = semua.
- `GET /api/nbfacin/coverage-otomatis` → lima baris urut korpus; kode tanpa baris di `COVERAGE` → `{id, oldId:"", nama:""}`.

**Rumus** (`services/coverage.go`) — port `Activity\CountPremi_ACT.xml` (57 langkah dibaca utuh lewat pengurai langkah):
Prorate langkah 3–8; PctAdjustment 11–15; LossLimit / SubLimit kosong → 100 (17–18); rumus per basis × mode 19–50 dijalankan
BERURUTAN, yang kemudian menimpa (PctAdjustment ≠ 0 → pembagi ×100; NetRate terisi → menggantikan Rate); diskon 51–52;
TotalGrossPremi 55–56. Total per mata uang = `Activity\SumTotalTSIPremiGross_Act.xml` cabang FIRE (Rate = Premium / TSI
× 1000; TSI 0 → 0). **Tidak diport:** langkah 9 (jalur EDM), 53–54 (min / max rate `SetRatePolis`, tahap C4), basis 5.
`[dugaan]` DiffYear = tahun penuh Start → End (DateTimeDifference "Y" dibulatkan ke bawah); arti kode aksi syarat langkah
Pega (2 lanjut / 3 lewati / 5 jalankan) diturunkan dari konsistensi logika langkah 3, 8, 9, 14 — tidak tertulis di korpus.

**Tabel:** migrasi **193** (ditulis, belum dijalankan; urutan DEV … → 192 → **193**): `T_PROPERTYITEMLIST` `ALTER ADD`
`TOTAL_GROSS_PREMI` / `TOTAL_NET_RATE` NUMBER(38,8); `T_COVERAGELIST` sebagian (34 kolom, rancangan berinduk jamak,
`PARENT_TABLE` `T_PROPERTYITEMLIST`, tanpa FK, indeks `(PARENT_TABLE, PARENT_ID)`), `CURRENCY_CODE` = mata uang item
(disetujui sesi 0f). Loader: skema 199 kolom diperiksa. Peta kolom: `STRUKTUR-TABEL-NB-FACIN.md`.

### Keputusan agent (menunggu konfirmasi)

| # | Keputusan | Dasar |
| --- | --- | --- |
| A155 | `@Math.divide(…, 20)` dihitung dengan presisi lebar lalu dibulatkan **setengah-ke-atas** ke 20 desimal; nilai tersimpan dibulatkan setengah-ke-atas ke 8 desimal (NUMBER(38,8), pola A148) | mode pembulatan `@Math.divide` Pega tidak tertulis di korpus `[dugaan]` |
| A156 | `item.PctAdjust1` dianggap kosong / 0: PctAdjustment = PctAdjustOther bila ≥ 60, selain itu 0; langkah 16 tidak berlaku | `.PctAdjust1` tidak ada di model tiket 39; arahan sesi 0f ("empty/0, catat sebagai keputusan agent") |
| A157 | `modeDiskon` kosong diturunkan: DiscountPercentage terisi → percent; kalau tidak, Discount terisi → amount; keduanya kosong → tanpa diskon | `Param.DiscountStatus` dikirim layar Pega, tidak tersimpan di data |
| A158 | Saat PUT mode hitung tiap coverage diturunkan: Rate terisi → percent; Rate kosong dan Premium terisi → amount; keduanya kosong → hanya TSI / TSILiability / Prorate, premi dibiarkan kosong (diperhalus A162) | `Param.PremiStatus` tidak tersimpan; heuristik disetujui sesi 0f |
| A162 | Rate DAN Premium terisi: percent, KECUALI percent tidak menghasilkan kembali premi yang dikirim sedangkan amount menghasilkan kembali rate yang dikirim → amount (premi pengguna dipertahankan) | `[terverifikasi]` temuan review sumbu spec: rate dibalik disimpan 8 desimal, sehingga A158 polos menggeser premi ketikan saat Save (TSI 123456789012, premi 1000000 → 999999.9909972); Pega menyimpan rate 20 desimal di halaman. Uji `TestHitungSimpanPremiTetap` (+ mutasi tertangkap) |
| A163 | `coverageBasis` kosong atau di luar 1–5 → 400 (POST dan PUT; basis 5 → 400 sesuai permintaan); hasil hitung yang melebihi NUMBER(38,8) (30 digit bulat) → 400 | Pega tidak menolak basis kosong — langkah rumus 19–50 tidak cocok dan premi tidak dihitung `[terverifikasi]` langkah 19–50. ⚠️ Data lama hasil loader dengan basis kosong atau 5 membuat Save tab Object ditolak 400 sampai tahap C2–C4 `[pertanyaan terbuka]` |
| A164 | `Day` case (EDMDay) terisi tetapi bukan bilangan / 0 → 409; LossLimit / SubLimit bernilai nol dalam bentuk apa pun ("0", "0.00") → 100 | Pega membandingkan teks `""` / `"0"` pada langkah 17 `[dugaan]` — "0.00" mungkin tidak didefault di Pega; layar mengirim desimal terurai |
| A159 | `TotalGrossPremi` dihitung ulang saat PUT (Σ Premium, item tanpa coverage → kosong); ⚠️ bagian TotalNetRate diganti tiket 44 (A166) — `TotalNetRate` disimpan apa adanya dari PUT; `PCT_ADJUSTMENT` disimpan tetapi tidak dikirim JSON | rumus TotalNetRate tahap C2; kontrak frontend tanpa pctAdjustment |
| A160 | Popup Choose Coverage **tanpa** saringan ACTIVESTATUS (`.ACTIVESTATUS = 1 OR IS NULL`, syarat D/E RD) | `[terverifikasi]` RD `ReportDefinition\BrowseCoverageFacIn_RD.xml` `pyFilterLogic` "A AND B AND C AND (D OR E)" dengan D/E atas ACTIVESTATUS; DDL `DDL\COVERAGE_FACIN.txt` = VIEW lima kolom (ID, BIZCODE, NAMACOVERAGE, TYPE, OLDID) — 0 kemunculan ACTIVESTATUS (dihitung dua cara: daftar kolom view dan `grep -c`). Kunci JSON `M_COVERAGE` untuk status aktif tidak ditebak; DDL `M_COVERAGE` tidak ada. `[pertanyaan terbuka]` DBA: tambah ACTIVESTATUS ke view, atau semua baris M_COVERAGE memang aktif? |
| — | Coverage otomatis dibaca dari tabel `COVERAGE` (bukan `COVERAGE_FACIN`), urutan korpus 100815, 100828, 100829, 100825, 100840; baris pertama per ID | `[terverifikasi]` `Activity\AddCoverageAutoFire.xml` langkah 3 + Obj-Browse kelas `ASM-FW-GISFW-Int-COVERAGE` (`RDBList\SearchCoverageIDSQL.xml` "FROM COVERAGE"); disetujui sesi 0f ("ikuti korpus") |
