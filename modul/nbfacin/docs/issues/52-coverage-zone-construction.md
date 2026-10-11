# 52: Coverage FIRE — Zone, 4.2 Construction, dan sisa isian yang ditunda tiket 43

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 06-10-2026 (gambar layar Pega
> coverage 4.2 PRGBI) "saat pilih 4.2 PRGBI kenapa beda dengan pega? kolom2 di coverage disesuaikan dong dengan pega".
> Survei korpus oleh subagen (pengurai XML); butir kunci diperiksa ulang sesi 0f.

**Status:** frontend selesai 06-10-2026 di atas kontrak USULAN (uji hijau). Backend selesai 06-10-2026 (sesi c3, bab
"Backend" di bawah); migrasi 203 DITULIS, belum dijalankan.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\Section\CoverageItem.xml` (`ASM-FW-GISFW-DATA-COVERAGE!COVERAGEITEM`)

Hanya empat syarat di section ini yang bergantung pada coverage (L4252, L4399, L5049, L6983):

| Sel | Isian | Tampil bila | Aksi |
|---|---|---|---|
| 6 | Zone (`.Zone`, autocomplete) | `.Coverage` ∈ {100835, 100821, 100829, 100816, 100817, 100828} dan Zone kosong | change → DT `SetPropertyCoverageDT_FacIn` + `CountPremi_ACT(CallAct="rate")` |
| 11 + 12 | Zone (teks baca-saja) + ikon "Switch to edit mode" | coverage sama dan Zone terisi | ikon → `SetShowZone_DT` (`.Zone=""`, `.ConstructionEarthquake=""`) |
| 13 | 4.2 Construction (`.ConstructionEarthquake`, autocomplete) | `.Coverage` ∈ {100821, 100829} | change → DT + `CountPremi_ACT` |

- Daftar Zone = RD `ASM-FW-GISFW-INT-ZONES!BROWSEZONES_RD`: tabel ZONES, saring `.OLDID = OLDID coverage` dan zona X/Y/Z
  (RSMD) atau 1–6 (gempa) dari `SetPropertyCoverageDT_FacIn` langkah 10–11; disimpan `.Description` (mis. "Daerah Zona 5"),
  `.ID` ke `InputParam.CARI28`.
- Daftar Construction = `RDBList\SearchClassConstructionReas_SQL.xml`: `select distinct construction from construction order
  by CONSTRUCTION`.
- Tarif 4.2 PRGBI (100829) saat Zone berubah = `RDBList\RateEarthquake_Sql.xml`: `SELECT RATE AS CARI1 FROM POOLDATA.ZONEOJK
  WHERE DESCRIPTION={Earthquake.CARI1}` (`CountPremiRate_ACT` langkah 8, Earthquake.CARI1 = `.Zone`).
- Tarif lain (`SetRatePolis_ACT`): gempa 100821 `m_eqs_rate`, RSMD `m_rsmd_rate`, FLEXAS `m_flexas_rate`, banjir
  `m_flood_rate`, terorisme `m_terorisme_rate` — **semua lewat db-link `@asmd.sinarmas.co.id`** (server luar, isi tidak di
  korpus; sejenis OQ db-link). `CountPremiRate_ACT`: 4.1A CC / OTHERS = 0.1; 4.3 = 0.5 (DKI JAKARTA / BANTEN / JAWA BARAT)
  atau 0.45.
- Sisa tiket 43 (P-1): Number of Layer + grid Layering (`.CoverageBasis==5`, `.LayerList`: Layer No · TSI · TSI Liability ·
  Policy Rate · Premium, detail `LayerListDtl`, `GenerateLayerList_ACT`); tombol "View Indemnity Table" (`ViewIndemnity_LA`,
  `fire.m_bi_indemnity@asmd…` — db-link); Indemnity Unit → `SetIndemnityRate_ACT` (`m_bi_indemnity`).
- ASM TSI / ASM Premium, tombol Spreading / Deductible: kontainer `1 = 2` → tersembunyi di Pega.
- Kode coverage: 100829 = 4.2 PRGBI, 100828 = 4.1A CC (contoh kasus); 100821, 100835, 100816, 100817 `[dugaan]` gempa /
  RSMD lain.

## Jawaban work owner 06-10-2026

1. DDL dikirim (`D:\migrasi\RNM\DDL\`): `ZONES.txt` (view `POOLDATA.ZONES` atas `m_zones`: ID, OLDID, ZONE, DESCRIPTION,
   RATE, DESCRIPTION2, ACTIVEDATE, TYPE), `ZONEOJK.txt` (tabel `POOLDATA.ZONEOJK`: ID, ZONE, DESCRIPTION, RATE NUMBER(20,3)),
   `CONSTRUCTION.txt` (view `POOLDATA.CONSTRUCTION` atas `m_construction`). Hanya dibaca — tanpa migrasi.
2. "setuju": tarif dari tabel db-link tidak diisi otomatis; ‰ Gross Rate tetap manual untuk coverage itu. Zone / 4.2
   Construction tetap tampil dan tersimpan.
3. "coverage basis 5 tidak usah, itu tidak dipakai lagi" — Layering tidak dibuat.

## Frontend (sesi 0f)

- `components/FormCoverage.tsx`: Zone (sel 6 autocomplete `IsianSaran`; sel 11/12 teks + ikon "Switch to edit mode" =
  `SetShowZone_DT`), 4.2 Construction (sel 13), Choose Coverage mengosongkan Zone / Construction lalu `callAct="rate"`;
  Zone dipilih → `callAct="rate"`. `labels.ts` (`ZONA_COVERAGE`, `COVERAGE_ZONA`, `COVERAGE_KONSTRUKSI`); `api.ts`
  (`CoverageObjek.zone` / `.constructionEarthquake`, `hitungCoverage.callAct`, `saranZona`, `saranKonstruksiGempa`);
  `nbfacin.css`; uji `FormCoverageZona.test.tsx`.

## Kontrak USULAN (menunggu sesi c3)

- `CoverageObjek` + `zone`, `constructionEarthquake` (teks; simpan di GET/PUT `…/objek`, kolom T_COVERAGELIST).
- `GET /api/nbfacin/zona?coverage=&oldId=&q=` → `{ baris: [{ id: ZONES.ID, label: ZONES.DESCRIPTION }] }`; saring
  `OLDID = oldId` (tanpa peka huruf) dan ZONE ∈ X/Y/Z bila coverage 100835/100828/100817/100816, 1–6 bila 100829/100821;
  `q` = mengandung, pada DESCRIPTION.
- `GET /api/nbfacin/konstruksi-gempa?q=` → `{ baris: [{ id, label }] }` = `select distinct construction from construction
  order by CONSTRUCTION`.
- `POST …/hitung-coverage` + `callAct: "rate"`: coverage 100829 dengan Zone terisi → `.Rate` = `ZONEOJK.RATE` WHERE
  `DESCRIPTION = zone` (`RateEarthquake_Sql`), lalu premi dihitung; coverage lain: tarif tidak diubah (butir 2).

## Backend (sesi c3, 06-10-2026)

- **K52-1** Kolom: `T_COVERAGELIST.ZONE` VARCHAR2(1000) dan `CONSTRUCTION_EARTHQUAKE` VARCHAR2(500), migrasi
  `203_coverage_zona.sql` (+ `_down`) — ALTER ADD biasa (pola 195). ZONE ada di rancangan sebagai VARCHAR2(50); dilebarkan ke
  lebar `ZONEOJK.DESCRIPTION` (amandemen loader `amandemenLebar`, pola butir 80) karena nilainya = `ZONES.DESCRIPTION` dan
  sampel `DDL\CONTOH` sudah sampai 111 byte (402 tag Zone; dihitung dua cara). CONSTRUCTION_EARTHQUAKE tidak di rancangan
  (0 sampel) → amandemen loader `amandemenKonstruksiGempa`; lebar = kolom konstruksi saudaranya (`*_CONSTRUCTION` 187),
  lebar sumber belum terverifikasi (view atas JSON). Lebih panjang → 400 ber-jalur
  (`baris[i].items[j].coverages[k].zone paling banyak 1000 byte`).
- **K52-2** `GET /api/nbfacin/zona?coverage=&oldId=&q=` (`repository/zona.go` `SaranZona`): RD `BrowseZones_RD` versi NB
  (sama dengan RNW; Endorsment beda isi, tidak dipakai). Saringan kosong dibuang; `UPPER` kedua sisi; ZONE ∈ daftar
  `SetPropertyCoverageDT_FacIn` menurut kode `coverage` (kode lain / kosong → tanpa saringan zona); `oldId` kosong →
  tanpa saringan OLDID (seperti RD). Urut DESCRIPTION, ID (RD tanpa urutan); batas 50 (RD 500, pola A124). Parameter
  dipangkas, > 255 karakter → 400; tanpa basis data → 503.
- **K52-3** `GET /api/nbfacin/konstruksi-gempa?q=`: DISTINCT CONSTRUCTION urut CONSTRUCTION, NULL dibuang, `q` Contains
  tidak peka huruf, batas 50.
- **K52-4** ⛔ *Diganti K52-7 (keputusan work owner lanjutan) — teks lama dipertahankan sebagai riwayat.*
  `POST …/hitung-coverage` + `callAct`: `"rate"` → sebelum `CountPremi`, bila `coverage.oldId` **persis**
  `"4.2 PRGBI"` (perbandingan `==` harfiah `CountPremiRate_ACT`, Okupasi.CARI2 = `.OLDID` — BUKAN kode 100829 seperti
  kontrak usulan; keduanya sama untuk data contoh) dan `zone` tidak kosong → Rate = `ZONEOJK.RATE WHERE DESCRIPTION = zone`
  (baris pertama urut ID bila ganda `[dugaan]`), lalu premi dihitung seperti biasa. ⚠️ **Menyimpang dari Pega, sesuai
  kontrak:** zona tanpa baris ZONEOJK (atau RATE NULL) → Rate DIBIARKAN (Pega mengosongkannya). ⚠️ **Juga menyimpang,
  sesuai kontrak ("dengan Zone terisi"):** Zone KOSONG pada 4.2 PRGBI → Rate tidak diubah, padahal Pega tetap menjalankan
  `RateEarthquake_Sql` (langkah 8–9 hanya bersyarat `Okupasi.CARI2=="4.2 PRGBI"`), tidak menemukan baris, dan
  mengosongkan `.Rate`. Akibatnya di layar: Choose Coverage 4.2 PRGBI (frontend mengosongkan Zone lalu `callAct="rate"`)
  di Pega mengosongkan Rate, di sini Rate coverage sebelumnya terbawa. Coverage lain → Rate tidak diubah dan basis data
  tidak dibaca. `callAct` lain → 400; absen = hitung biasa. Butuh ZONEOJK tetapi tanpa
  basis data → 503.
- **K52-5** ⛔ *Diganti K52-8 (cabang lain kini dibangun).* Cabang `CountPremiRate_ACT` lain **TIDAK ada** di backend (tiket 43 / 44 hanya memakai lima OLDID itu untuk
  flag net rate `CekNetRate_ACT`, `services/netrate.go`) dan tidak ditambahkan: FLEXAS = max RATELOWER okupasi,
  "4.1A CC" / "OTHERS" = 0.1, "4.3" = 0.5 (DKI JAKARTA / BANTEN / JAWA BARAT) atau 0.45. `CountPremi_ACT` juga memanggil
  `CekNetRate_ACT` tepat sebelum `CountPremiRate_ACT`; syarat langkahnya belum terverifikasi (penelaah spec membacanya
  `Param.CallAct=="rate"`). Hitung-coverage tidak menjalankannya — flag net rate tetap jalur tiket 44
  (`POST …/hitung-net-rate` dan simpan `…/objek`).
- **K52-6** Urutan deploy: 203 sebelum biner baru — biner baru membaca / menulis ZONE dan CONSTRUCTION_EARTHQUAKE di
  `…/objek`; tanpa 203 GET / PUT objek gagal ORA-00904.

## Sisa

- `CountPremiRate_ACT` lain pada `CallAct="rate"` (4.1A CC / OTHERS = 0.1; 4.3 = 0.5 / 0.45 menurut provinsi; FLEXAS =
  max RATELOWER okupasi) — sesi c3 diminta melaporkan apakah sudah ada di backend tiket 43 / 44.
- View Indemnity Table dan Indemnity Unit (`m_bi_indemnity`, db-link) tetap ditunda.

## Keputusan work owner 06-10-2026 (lanjutan)

1. "ikutin sesuai pega saja sesuai saran km": 4.2 PRGBI dengan Zone kosong / tak ada di ZONEOJK → ‰ Gross Rate dikosongkan
   (`.Rate = InputEarthquake.pxResults(1).CARI1`, hasil kosong = kosong), seperti Pega.
2. "iya silahkan ditambah, jangan halu ya": cabang `CountPremiRate_ACT` lain pada `CallAct="rate"` ditambahkan. Bukti
   `[terverifikasi]` `Activity\CountPremiRate_ACT.xml` (sesi 0f, semua langkah `pyStepsPreCondition=true`):
   - @427 `Local.PROVINSI = @if(…LocationList(param.IdxProperty).Property.Province=="DKI JAKARTA","0.5",@if(=="BANTEN","0.5",
     @if(=="JAWA BARAT","0.5","0.45")))`;
   - FLEXAS (@973–1651): per okupasi `.Property.OccupationList` (`.OccupationId`) `RDBList\RateFlexas_Sql`: `SELECT RATELOWER
     AS CARI1 FROM OCCUPATION WHERE OCCUPATIONCODE = {Okupasi.CARI1}`; nilai terbesar (`Param.tes > .CARI1`; RATELOWER
     NUMBER(20,6)) → `.Rate`;
   - "4.1A CC" / "OTHERS" (@1874) → `.Rate = 0.1`; "4.3" (@2019) → `.Rate = Local.PROVINSI`.
   - Layar mengirim `province` dan `occupationIds` objek bersama `callAct` (isian saat itu, belum tentu tersimpan).

## Backend lanjutan (sesi c3, 06-10-2026)

- **K52-7** `callAct` "rate" mengikuti `CountPremiRate_ACT` (`services/ratecoverage.go` `terapkanRate`), kunci = OLDID
  harfiah. 4.2 PRGBI: Rate = `ZONEOJK.RATE` menurut Zone; Zone kosong, zona tanpa baris, atau RATE NULL → Rate
  **DIKOSONGKAN** seperti Pega, lalu premi dihitung (Rate kosong → premi 0). Zone kosong tidak membaca basis data (Oracle
  `DESCRIPTION = ''` = NULL, tidak cocok baris mana pun); Zone berisi spasi dicari apa adanya. Selisih K52-4 (Rate dibiarkan)
  dicabut.
- **K52-8** Cabang lain:
  - FLEXAS — `RateFlexas_Sql` untuk setiap `occupationIds` (kosong / ganda dibuang), SATU kueri `OCCUPATIONCODE IN (...)`
    (`repository/occupation.go` `RateLowerOkupasi`), baris pertama per kode urut ID bila ganda `[dugaan]`; Rate = RATELOWER
    TERBESAR. Kode tanpa baris atau RATELOWER NULL dilewati; tanpa okupasi / tanpa nilai → Rate kosong.
    `[dugaan]` dua hal: (a) `Param.tes > .CARI1` dibaca NUMERIK (RATELOWER NUMBER(20,6); tipe properti CARI1 tidak ada di
    korpus — bila Pega membandingkan teks, mis. "10" < "9", hasil bisa beda); (b) di Pega okupasi PERTAMA tanpa hasil
    menambahkan nilai kosong ke temprate — perilaku pembandingan kosong tidak terverifikasi, di sini dilewati.
    ⚠️ `.OccupationId` yang disimpan popup Choose Occupation = `OCCUPATION.OLDID` (tiket 40), sedangkan SQL asal
    mencocokkan `OCCUPATIONCODE` — diikuti apa adanya; bila di data OLDID ≠ OCCUPATIONCODE, FLEXAS akan selalu kosong
    (sama dengan Pega). Isi data tidak ada di korpus — belum terverifikasi.
  - "4.1A CC" / "OTHERS" → Rate 0.1.
  - "4.3" → 0.5 bila `province` PERSIS "DKI JAKARTA" / "BANTEN" / "JAWA BARAT" (peka huruf, tanpa pangkas), selain itu 0.45
    (termasuk province kosong).
  - OLDID lain → Rate tidak diubah dan basis data tidak dibaca (tarif db-link, butir 2).
- **K52-9** Masukan: `province` (teks) dan `occupationIds` (larik teks) di badan hitung-coverage. Dikirim BERISI tanpa
  `callAct` → **400** (dipilih menolak, bukan mengabaikan: kesalahan pemanggil terlihat; nilai kosong / larik kosong
  tetap diterima). Batas: province ≤ 255 karakter, occupationIds ≤ 100 kode, tiap kode ≤ 1000 byte
  (= `T_OCCUPATIONLIST.OCCUPATION_ID`). Butuh OCCUPATION / ZONEOJK tanpa basis data → 503; cabang yang tidak membaca
  (4.3, 4.1A CC, Zone kosong, tanpa okupasi) tetap 200. Tanpa migrasi baru.
- **K52-10** `CekNetRate_ACT` (dipanggil `CountPremi_ACT` sebelum `CountPremiRate_ACT`, syarat belum terverifikasi) tetap
  tidak dijalankan hitung-coverage — jalur tiket 44.
- **K52-11** ✅ *Terjawab 06-10-2026 — pengecualian ADR-0035, tanpa tabel dan tanpa migrasi (bab "Keputusan work owner 06-10-2026 — pengecualian ADR-0035").* Ketegangan dengan **ADR-0035** ("Nol tanggal, ambang, atau nama tetap di dalam kode"): tarif 0.1 / 0.5 /
  0.45 dan tiga nama provinsi ditanam di `services/ratecoverage.go`, persis seperti Pega menanamnya di
  `CountPremiRate_ACT` (literal di langkah activity itu; tabel tetapan lain yang memuatnya belum dicari). Dibiarkan atas "ikutin sesuai pega" dan "tanpa
  migrasi baru"; ADR-0035 belum mencatat pengecualian ini. ⛔ **Pertanyaan terbuka untuk work owner:** pindahkan ke tabel
  tetapan (butuh migrasi baru + pemilik isinya) atau catat sebagai pengecualian ADR-0035.
- **K52-12** Catatan penelaah spec (06-10-2026), dicatat tanpa mengubah perilaku:
  - "Baris pertama urut ID" (FLEXAS ganda, juga ZONEOJK) = urutan TEKS — `OCCUPATION.ID` VARCHAR2(1000), `ZONEOJK.ID`
    VARCHAR2(20) ("10" sebelum "9"). Pega tanpa ORDER BY; urutan ini hanya supaya pasti, `[dugaan]`.
  - Halaman `InputFlexas` / `InputEarthquake` tidak pernah di-Page-Remove di activity. Bila RDB-List Pega tidak
    mengosongkan `pxResults` saat nol baris, okupasi / zona tanpa baris di Pega bisa memakai nilai basi panggilan
    sebelumnya; di sini dilewati / dikosongkan. Perilaku RDB-List itu belum terverifikasi.
  - "Rate kosong → premi 0" hanya bila Net Rate kosong: `HitungCoverage` mode percent memakai `NetRate` bila terisi
    (tiket 44, `services/coverage.go`). Pemanggil yang mengirim `mode:"amount"` bersama `callAct` menghitung Rate dari
    Premium sesudah cabang rate (layar selalu mengirim percent).
  - `province` absen bersama `callAct` pada 4.3 → 0.45 tanpa galat (sama dengan Province kosong di Pega).

## Keputusan work owner 06-10-2026 — pengecualian ADR-0035 (menjawab K52-11)

"a. brarti ini ikut sesuai saran km": tarif tetap 0.1 ("4.1A CC" / "OTHERS"), 0.5 / 0.45 ("4.3") dan tiga nama provinsi
("DKI JAKARTA", "BANTEN", "JAWA BARAT") **tetap ditanam di kode** (`backend/services/ratecoverage.go`) sebagai
**pengecualian ADR-0035** — nilainya literal di activity Pega (`Activity\CountPremiRate_ACT.xml` @427, @1796, @1941). Tidak
dipindah ke tabel tetapan; tanpa migrasi.
- **K52-13** (c3, 06-10-2026) Bug DEV: Save tab Object "`baris[0].items[0].coverages[2].zone paling banyak 1000 byte`" —
  nilai `ZONES.DESCRIPTION` (view atas JSON m_zones, lebar tidak tertulis) melampaui lebar `ZONEOJK.DESCRIPTION` yang
  dipakai 203. Perintah work owner "ini juga perbaiki, buat 4000": batas services `zone` 4000 byte, migrasi BARU
  `204_coverage_zona_4000.sql` `ALTER TABLE … MODIFY (ZONE VARCHAR2(4000))` (+ `_down` kembali 1000 — ⚠️ GAGAL ORA-01441
  bila ada Zone > 1000 byte tersimpan); 203 tidak diubah (kemungkinan sudah jalan di DEV). Amandemen loader
  `T_COVERAGELIST.ZONE` VARCHAR2(4000); penjaga rancangan memutar 182..204. ⛔ Sampai 204 jalan, biner baru menerima Zone
  1001–4000 byte tetapi INSERT gagal ORA-12899 — jalankan 204 sebelum / bersama biner.
  ⚠️ Akibat ikutan: Zone > 1000 byte tidak akan pernah cocok `ZONEOJK.DESCRIPTION` VARCHAR2(1000) → rate 4.2 PRGBI kosong
  (K52-7) untuk zona sepanjang itu — bila ZONES dan ZONEOJK memakai teks deskripsi yang berbeda, itu pertanyaan data.
  `CONSTRUCTION_EARTHQUAKE` VARCHAR2(500) TIDAK diubah: sumbernya view `CONSTRUCTION` atas JSON m_construction, lebar
  `CONSTRUCTION` tidak tertulis (ekspresi dot-notation JSON mengembalikan VARCHAR2(4000) secara bawaan `[dugaan]`), jadi
  bisa melampaui 500 — menunggu keputusan work owner (ukur dulu dengan kueri di bawah).

  Langkah DEV (work owner / DBA, BUKAN agent):
  1. Sebelum: `SELECT NAMA FROM <SKEMA>.T_MIGRASI WHERE NAMA LIKE '20%' ORDER BY NAMA;` (harap 200..203 tercatat, 204 belum)
     dan `SELECT COLUMN_NAME, DATA_TYPE, DATA_LENGTH FROM ALL_TAB_COLUMNS WHERE OWNER = '<SKEMA>' AND TABLE_NAME =
     'T_COVERAGELIST' AND COLUMN_NAME IN ('ZONE', 'CONSTRUCTION_EARTHQUAKE');` (harap VARCHAR2 1000 / VARCHAR2 500).
     Bila 203 BELUM tercatat, `-migrate` menjalankan 203 lalu 204 berurutan.
  2. `-migrate`.
  3. Sesudah: kueri ALL_TAB_COLUMNS di atas → ZONE `DATA_LENGTH` = 4000; `SELECT NAMA FROM <SKEMA>.T_MIGRASI WHERE NAMA =
     '204_coverage_zona_4000';`.
  Ukur panjang sebenarnya (opsional, baca saja):
  `SELECT MAX(LENGTHB(DESCRIPTION)) FROM POOLDATA.ZONES;` · `SELECT MAX(LENGTHB(DESCRIPTION)) FROM POOLDATA.ZONEOJK;` ·
  `SELECT MAX(LENGTHB(CONSTRUCTION)) FROM POOLDATA.CONSTRUCTION;` · `SELECT MAX(LENGTHB(ZONE)) FROM <SKEMA>.T_COVERAGELIST;`
