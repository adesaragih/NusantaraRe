# 37: Add alamat risiko — popup InputRiskAddress, saran Zip Code dari RW, simpan ke RISKADDRESS

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026:
> *"ini tampilan menu add nya … kolomnya otomatis ke isi saat zipcode dimasukin, data nya di ambil dari tabel rw"*,
> disertai dua tangkapan layar. Tiga keputusan dijawab work owner lewat AskUserQuestion 03-10-2026 (lihat bawah).

**What to build:** tombol **Add** di popup Choose Risk Location (tiket 36) membuka popup isian alamat baru:
Country · Zip Code · (Province · City · District · Territory, tampil berantai) · Title · Address · Save / Close. Memilih
saran Zip Code mengisi enam medan dari tabel RW. Save menyimpan alamat ke master RISKADDRESS, lalu langsung mengisi
Risk Address objek.

**Blocked by:** tiket 36.

**Status:** frontend selesai 03-10-2026 (uji hijau); backend selesai 03-10-2026 (sesi c3, port Go butir 81).

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\`

- `NB FacIn\Section\ChooseRiskAddress.xml` sel 74 `Add`: DT `setOutputValue_DT` (`pyPortal.IsEdit=True`) →
  `ClearPageRiskAddress` (Page-New `InputRW`, `InputRiskAddress`) → showHarness `ChooseRiskLocation` (popup, tanpa
  WindowName).
- `NB FacIn\Section\InputRiskAddress.xml` (halaman `InputRiskAddress`, kelas `ASM-FW-GISFW-Int-RISKADDRESS`):
  sel 3 Country · 5 Province (visible `NationName != ''`) · 7 City (`ProvinceName != ''`) · 9 District (`CityName != ''`)
  · 11 Territory (`DistrictName != ''`) · 13 Zip Code · 20 Title (dropdown, `pyHasNoSelection=false`) · 21 Address
  (textarea) · 28 Save · 29 Close. Sel 18 Code, 19 Type, 22 Status Trans Pusat tersembunyi (`never`). Tidak ada medan
  wajib.
- Sel 13 Zip Code: autocomplete RD `BrowseRW_RD` (kelas `ASM-FW-GISFW-Int-RW` = `pooldata.rw`, `RDBList\BrowseRW_SQL.xml`).
  Memilih baris mengisi `.ZipCode → PostalCode`, `.NATIONNAME → NationName`, `.PROVINCENAME → ProvinceName`,
  `.CITYNAME → CityName`, `.DistrictName → DistrictName`, `.Note → TerritoryName`. RD menyaring `.STS_AKTIF = "1"`.
- Sel 28 Save: `SaveRiskAddress_Act` → `SaveAccumulationByRiskAddress_Act` → refresh `GridRiskAddress`.
- `Activity\SaveRiskAddress_Act.xml`: ID = `"UNKNOWNID"` bila kosong; RDB-List `UpdateMasterRiskAddress_SQL` =
  `BEGIN pooldata.InsertUpdateRISKADDRESS(ID, NationName, ProvinceName, DistrictName, CityName, TerritoryName, Title,
  Address, PostalCode, HASIL1 out, HASIL2 out); COMMIT; END;`. Lalu menulis balik ke
  `pyWorkPage.OfferFacIn.LocationList(ObjectNo).Property`: RoadName = Address, RoadType = Title, ASMAddress dirangkai,
  ASMZipCode, Country, Province, ASMCity, ASMDistrict, ASMRW. ID baru TIDAK ditulis balik.
- Cacat Pega: langkah penutup form mensyaratkan `HASIL1` mengandung "Sukses", padahal teks sukses prosedur adalah
  "Data sudah di simpan dengan ID :…", jadi form tidak pernah tertutup. ID di halaman tetap "UNKNOWNID", sehingga Save
  kedua menambah baris baru lagi.
- `DDL\INSERTUPDATERISKADDRESS.txt`: ID ada → UPDATE delapan medan; tidak ada → `vID := getcurrentsite ||
  lpad(RISKADDRESS_SEQ.nextval,12,'0')` lalu INSERT. Sukses: `ErrMsg := 'Data sudah di simpan dengan ID :' || vID`,
  `StsSimpan := 1`. Gagal: `StsSimpan := 0`, ErrMsg berisi `sqlerrm`. Fungsi `getcurrentsite` tidak ada di korpus.
- Title: REGEXP `(DESA|DUSUN|GANG|GEDUNG|JL\.|KOMPLEK|PERUMAHAN|OTHERS)` di
  `RDBList\SearchAccumulationbypersetase_SQL.xml`, sama dengan tangkapan layar.

## Keputusan work owner (03-10-2026)

- ~~**W-1** Save memanggil prosedur Pega yang sama, `pooldata.InsertUpdateRISKADDRESS`; ID dibuat prosedur.~~ →
  ⛔ **Diganti butir 81** (work owner, AskUserQuestion sesi c3 03-10-2026, *"Port to Go, keep ADR-0043"*): W-1 bertentangan
  dengan ADR-0043 (*"jangan ada lagi pemanggilan procedure"*). Logika prosedur di-port ke Go — lihat bab Backend.
- **W-2** `SaveAccumulationByRiskAddress_Act` (master akumulasi, `RDBMASTERACCUMULATION`) **ditunda**.
- **W-3** Risk Address ID objek diisi ID baru dari prosedur (Pega membiarkannya kosong).

## Keputusan agent

- **I-1** Country / Province / City / District / Territory = isian teks bebas. Pega juga mengizinkan isian bebas;
  saran berantai `BrowseNation_RD` … `BrowseTeritory_RD` = tahap berikut.
- **I-2** Saran Zip Code muncul setelah 3 karakter, jeda ketik 400 ms; disaring Zip Code + STS_AKTIF saja.
- **I-3** Setelah Save, kedua popup ditutup dan objek terisi (Pega: form tetap terbuka karena cacat "Sukses").
- **I-4** Zip Code dan Address wajib (Pega tanpa validasi), supaya tidak ada baris master kosong.
- **I-5** Judul modal = "Add" (label tombol pembuka; harness tanpa WindowName).

## Kontrak

- `GET /api/nbfacin/rw?zipCode=` → `{ baris: [{ zipCode, territoryName, districtName, cityName, provinceName,
  nationName }] }`.
- `POST /api/nbfacin/risk-address` badan `{ nationName, provinceName, districtName, cityName, territoryName, title,
  address, postalCode }` → `{ id }`.

## Acceptance criteria

- [x] Add membuka popup; medan tampil berantai; Title DESA; label diuji ke korpus.
- [x] Memilih saran Zip Code mengisi enam medan; Save → objek terisi + ID baru (W-3).
- [x] Backend `GET /api/nbfacin/rw` dan `POST /api/nbfacin/risk-address` (~~prosedur, W-1~~ → port Go, butir 81).
- [ ] Akumulasi (W-2, tiket tersendiri).

## Backend (sesi c3, 03-10-2026) — disusun agent

- `GET /api/nbfacin/rw?zipCode=` → `{"baris":[{"zipCode","territoryName","districtName","cityName","provinceName","nationName"}]}`;
  400 bila `zipCode` (dipangkas) < 3 atau > 255 karakter; 503; 500. Tanpa identitas (pola lookup). Sumber `POOLDATA.RW`,
  `STS_AKTIF = '1'`, DISTINCT, ZIPCODE **diawali** kata, urut ZIPCODE, ≤ 50 baris.
- `POST /api/nbfacin/risk-address` badan `{nationName, provinceName, districtName, cityName, territoryName, title, address,
  postalCode}` → **201** `{"id"}`. 401 tanpa identitas (menulis data master); 400 `postalCode`/`address` kosong (I-4), medan >
  4000 bita, JSON rusak; 503; 500. Title tidak dicocokkan ke daftar.
- **Port Go prosedur `InsertUpdateRISKADDRESS` (butir 81)**, cabang INSERT (Pega selalu `P_ID = "UNKNOWNID"` untuk alamat
  baru): dalam SATU transaksi aplikasi — `SELECT GETCURRENTSITE || LPAD(TO_CHAR(RISKADDRESS_SEQ.NEXTVAL),12,'0') FROM DUAL`
  (ekspresi ID prosedur **persis**; isi fungsi `getcurrentsite` tidak ada di korpus, jadi dipanggil apa adanya, bukan
  dikarang) lalu `INSERT` sembilan kolom dengan urutan parameter prosedur. Nol `CALL`, nol `COMMIT`/`ROLLBACK` di SQL
  (ADR-0043, ADR-U-0029). Pemeriksaan prosedur `[terverifikasi]` DDL: **tidak** ber-`COMMIT` pada sukses; `ROLLBACK` pada
  galat (yang akan ikut membatalkan transaksi pemanggil — salah satu alasan port). Pesan `ErrMsg` "Data sudah di simpan dengan
  ID :…" tidak dipakai (tidak ada prosedur yang dipanggil), ID dibaca langsung. Cabang UPDATE prosedur tidak diport (tidak
  ada layar ubah). Akumulasi ditunda (W-2).

**Bukti saran Zip Code `[terverifikasi]`:** `Section\InputRiskAddress.xml` sel 13 = pxAutoComplete, RD `BrowseRW_RD`, medan
cari `.ZipCode` (`pyUseForSearch`), `pyUseParameterForSearch` false — **tidak ada** pengaturan "diawali / mengandung" di
section (dicari: nol `pyMatch*`); RD: filter `ZipCode/PROVINCENAME/CITYNAME/DistrictName/Note = Param.*` + `STS_AKTIF = "1"`,
DISTINCT, `pyMaxRecords` 10000. `.NATIONNAME` tidak ada di DDL `RW.txt` (kolomnya `NATION`); alias `NATION as "NATIONNAME"`
ada di `RDBList\BrowseRW2_SQL.xml` (kelas RW yang sama) → `nationName` = `RW.NATION`.

**Keputusan agent — DISETUJUI work owner 03-10-2026** (butir 82, diteruskan sesi `nusantarare-0f`: *"setuju keputusan
agent"*):

| # | Keputusan | Dasar |
| --- | --- | --- |
| A123 | Zip Code **diawali** kata (`LIKE 'kata%'`) | mode cocok autocomplete tidak tertulis di korpus (bawaan Pega `belum terverifikasi`); kode pos diketik dari depan |
| A124 | Paling banyak **50** saran, urut ZIPCODE (lalu kolom lain) | usul sesi 0f; RD 10000 |
| A125 | `nationName` = `RW.NATION` | alias `BrowseRW2_SQL.xml`; pemetaan RD `.NATIONNAME` → kolom `[dugaan]` |
| A126 | POST menjawab **201** `{"id"}` | pola POST pembuat (tiket 29) |
| A127 | Medan > 4000 bita → 400; `zipCode` saran < 3 / > 255 karakter → 400 | lebar kolom RISKADDRESS; I-2 |
| A128 | ID baru > 15 karakter atau kosong → galat (500), tidak dipotong | prosedur `vID varchar2(15)` gagal sama; isi `getcurrentsite` `[dugaan]` ≤ 3 karakter |
