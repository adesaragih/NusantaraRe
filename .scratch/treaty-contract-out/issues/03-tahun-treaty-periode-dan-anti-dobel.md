# 03: Tahun treaty — CRUD, gerbang periode, dan anti-dobel

**Status:** selesai (28-09-2026)

**Blocked by:** 01 (skema + sequence harus ada)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin membuat dan mengubah **tahun treaty** beserta grup
treaty, underwriting year, proporsi, dan masa berlakunya — **tanpa pernah mengetik nomor
identitas**, **tanpa** dapat membuat periode yang berakhir sebelum dimulai, dan **tanpa** dapat
membuat tahun yang sama dua kali. *(User story 1–3, 5–6 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas tahun treaty |
| `internal/repository` | Upsert dikunci `ID`; identitas dari sequence |
| `internal/services` | Gerbang periode; gerbang anti-dobel |
| `internal/handlers` | Endpoint tahun treaty |
| `frontend/` | Layar daftar + form tahun treaty |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyYear_Act` | `@BASECLASS` / `SAVETREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyYear_Act.xml` | orkestrator simpan |
| `NewInputTreatyYear_Act` | `@BASECLASS` / `NEWINPUTTREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewInputTreatyYear_Act.xml` | baris baru |
| `SetTreatyYear_Act` | `@BASECLASS` / `SETTREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetTreatyYear_Act.xml` | isi form dari baris terpilih |
| `CheckYear` | `@BASECLASS` / `CHECKYEAR` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/CheckYear.xml` | ⚠️ **satu-satunya** validasi: `isNumber(TreatyYear)` |
| `SaveMasterTreatyYear_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyYear_SQL.xml` | → `POOLDATA.PEGA_TREATYYEAR` |
| `BrowseTreatyYear_RD` | `ASM-FW-GISFW-INT-TREATYYEAR` / `BROWSETREATYYEAR_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseTreatyYear_RD.xml` | daftar tahun |

`[data DBA]` `POOLDATA.PEGA_TREATYYEAR` — **upsert dikunci `ID`**; identitas baru dibuat di basis
data sebagai `'1' || lpad(TreatyYear_seq.nextval, 6, '0')`; `TGLUPDATE` diisi `SYSDATE`;
**tidak `COMMIT` sendiri**; keluaran `StsSimpan` **1 = sukses / 0 = gagal**.

## ADR terkait

**ADR-0006** (identitas lewat sequence basis data), **ADR-0007** (jejak audit),
**ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] Tahun treaty dapat dibuat dengan **grup treaty, underwriting year, proporsi, dan masa
      berlakunya**. *(AC 4 spec; User story 1)*
- [ ] Identitas tahun treaty **tidak pernah diketik pengguna** — dibuat dari sequence.
      *(AC 5 spec; **ADR-0006**)*
- [ ] Identitas berbentuk `'1'` diikuti nomor urut ber-*padding* nol **6 digit**.
      *(AC 6 spec; `[data DBA]`)*
- [ ] Menyimpan tahun treaty yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec — sisi tahun; `[data DBA]` upsert dikunci `ID`)*
- [ ] Masa berlaku yang **berakhir sebelum dimulai** **ditolak**, dengan pesan yang menyebut
      field-nya. *(AC 9 spec — sisi tahun)*
- [ ] **Tahun treaty yang sama tidak dapat dibuat dua kali**: kombinasi **(`STARTDATE`, `ENDDATE`,
      `TREATYGROUPID`)** yang **sudah ada** ditolak, dengan pesan yang menyebut tahun treaty mana
      yang sudah memakainya. *(AC 73 spec; `[keputusan work owner]` — **aturan baru**)*
- [ ] Gerbang anti-dobel **tidak** menghalangi pembaruan baris itu sendiri — memperbarui tahun
      treaty yang sudah ada dengan periode & grup yang tidak berubah tetap **berhasil**.
- [ ] Tanggal mulai dan tanggal akhir tersimpan sebagai **tanggal**, bukan teks. *(AC 53 spec)*
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**; nilai status selain `1`
      **selalu** dibaca sebagai kegagalan. *(AC 38, 39 spec; **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit** — siapa dan kapan. *(AC 41 spec; **ADR-0007**)*
- [ ] ⚠️ **Fitur salin tahun treaty tidak dibangun.** Tidak ada jalur — layar, endpoint, maupun
      pekerjaan latar — yang menyalin isi satu tahun treaty ke tahun lain. Test yang menemukan
      jalur semacam itu **gagal**. *(AC 72 spec; `[fakta bisnis — work owner]` — penyimpangan
      sadar 3)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Mengapa fitur salin dibuang.** `[fakta bisnis — work owner]` `POOLDATA.PROSESCOPY` menyalin
isi satu tahun treaty ke tahun lain, dan dalam praktiknya membawa **nilai tahun lalu** — misalnya
batas QS tahun 2025 — ke tahun yang semestinya berbeda. Itu sumber kesalahan, bukan penghemat
waktu. Yang **tidak** dimigrasikan: `RDBList/SaveMasterCopyData_SQL.xml`,
`Activity/BrowseCopyData.xml`, `Activity/SaveTreatyYearMultiple_Act.xml`,
`Activity/NewInputTreatyYearMultiple_Act.xml`, dan bagian salin
`Activity/BrowseDeleteRowTreatyInContract.xml`.

⚠️ **Validasi existing sangat tipis.** `[terverifikasi]` `Activity/CheckYear.xml` hanya memeriksa
`@Default.isNumber(InputTreatyYear.TreatyYear)`. Tidak ada gerbang periode, tidak ada anti-dobel.
Kedua gerbang di tiket ini adalah **tambahan sadar** `[keputusan work owner]`, bukan tiruan.

⚠️ `[terverifikasi]` `RDBList/SaveMasterCopyData_SQL.xml` menulis
`dbms_output.put_line(errmsg)` dengan variabel lokal `errmsg` yang **dideklarasikan tetapi tidak
pernah diisi** — selalu mencetak NULL. Dicatat sebagai jejak; tidak dimigrasikan.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — gerbang anti-dobel adalah pertanyaan
keunikan di basis data; memalsukannya berarti tidak mengujinya.

```
go test ./internal/...
cd frontend && npm test
make check
```


---

## Pembacaan ulang XML — 28-09-2026 (sesi modul)

Nomor baris = nomor baris mentah berkas korpus (satu tag per baris). Rantai: `Harness/InboxTreatyContract.xml`
b151 `<pyLabel>InboxTreatyContract</pyLabel>` → `Section/GridTreatyContract.xml` (judul `TREATY CONTRACT OUT`
b1057) → `Section/InputTreatyContract.xml` (b2482) yang menyertakan `Section/InputDtlTreatyContact.xml` (b1206).

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| grid tahun | `InputTreatyContract.xml` `pyGridRDName BrowseTreatyYear_RD` b17355; judul kolom `Underwriting Year` b17378, `Transaction Year` b17531, `StartDate` b17684, `EndDate` b17837, `Treaty Group` b17990, `Reinsurance Type` b18136; sel `.UnderwritingYear` b18964, `.TreatyYear` b19125, `.StartDate` b19286, `.EndDate` b19446, `.TreatyGroupName` b19605, `.Proportion` b19724 (dipasangkan menurut urutan dalam satu grid `[dugaan kuat]`) | `GET /api/treaty-contract-out/tahun` → `InboxTreatyContract.tsx` enam kolom VERBATIM |
| `ReportDefinition/BrowseTreatyYear_RD.xml` | kelas `ASM-FW-GISFW-Int-TREATYYEAR` b40; 10 field b594–b735; sort `.ID DESC` b672; `pyMaxRecords` 500 b757 | `ORDER BY ID DESC`, halaman 20 (maks 500) |
| tombol `Add` | b16387 → `Activity/NewInputTreatyYear_Act.xml` mengosongkan `InputTreatyYear.ID/TreatyYear/TreatyGroupID` b469–b512, `ERRMSG` b422 | `formKosong()` |
| tombol `Edit` | b19939 → `Activity/SetTreatyYear_Act.xml` b20030: memuat `BrowseTreatyYear_RD` (b449, kelas b502) lalu menyalin Param → `InputTreatyYear.*` b1110–b1353 (ID, TreatyYear, TreatyGroupName, TreatyGroupID, Proportion, TglUpdate, UserID, StartDate, EndDate, UnderwritingYear, ReinsuranceType) | `formDari(baris)` |
| tombol `ReinsType` b20778 (membawa `.TreatyGroupName/.TreatyYear/.StartDate/.EndDate/.UnderwritingYear` b21310–b21346), `List Description` b22196 (`.TreatyYear` b22743) | membuka harness 2 dan 3 | berdiri, `disabled`, menyebut tiket 04 / 08 |
| tombol `Copy` b20459 (`Copy` b20409) + form `From`/`To` b2374/b3818 + `Proces` b5104 → `BrowseCopyData` b5128/b5211 | fitur salin | ➖ TIDAK dibawa (AC 72; `SaveMasterCopyData_SQL`, `SaveTreatyYearMultiple_Act`, `NewInputTreatyYearMultiple_Act`, `BrowseTreatyYearMultiple_RD` mati) |
| form lengkap `Input New Data` (`InputDtlTreatyContact.xml` b5437) | `ID` b6379 (`InputTreatyYear.ID`), `Treaty Group` b6560 (`.TreatyGroupName` b6628; pemilih `pyListSource reportdefinition` b6617 → `BrowseTreatyGroup_RD` b6624), `Reinsurance Type` b6800 → **`InputTreatyYear.Proportion`** b6829 (`.ID` pilihan b6857), `Start Date` b7532 (`.StartDate` b7561), `End Date` b7816 (`.EndDate` b7845), `Underwriting Year` b8004 → **`.TreatyYear`** b8035 (+`CheckYear` b8075/b8195), `Transaction Year` b8284 → **`.UnderwritingYear`** b8315 (+`CheckYear` b8350/b8475), `Modified Date` b9097 (`.TglUpdate` b9126), `Username` b9282 (`.UserID` b9313), `Save` b10332 → `SaveTreatyYear_Act` b10356/b10446, `Cancel` b10622 → `CancelActivityTreatyContract` b10640 | form layar; label VERBATIM per medan |
| form kedua di `InputTreatyContract.xml` sendiri (`Input New Data` b6351: `ID` b7282, `Treaty Year` b7455 + `CheckYear` b7520/b7646, `Treaty Group` b7732 (`.TreatyGroupID` b7758, `BrowseTreatyGroup_RD` b7792), `Reinsurance Type` b7925 → `InputTreatyYear.ReinsuranceType` b7954, `Modified Date` b8661, `By User` b9021, `Save` b10059 → `SaveTreatyYear_Act` b10082) | TANPA Start/End Date dan Underwriting Year; `ReinsuranceType` bukan parameter procedure | `[terbuka]` mana yang tampil di Pega; yang dibangun form LENGKAP (satu-satunya yang dapat mengisi seluruh parameter `PEGA_TREATYYEAR`) |
| `Activity/SaveTreatyYear_Act.xml` | `UserID←OperatorID.pyUserName` b280, `TglUpdate←@getCurrentTimeStamp()` b327; prasyarat `TreatyGroupID==""` b388/b622/b791, `TreatyYear==""` b411/b645/b814, `@Default.isNumber(TreatyYear)` b434/b668/b837; RDB-List `SaveMasterTreatyYear_SQL` b1100/b1160 → `POOLDATA.PEGA_TREATYYEAR` (10 param) | `models.PeriksaTahunTreaty` (grup wajib, tahun wajib, tahun angka) + `services.TahunTreatyTCO.Simpan` (UserID pelaku, TglUpdate sekarang) |
| `Activity/CheckYear.xml` | `@Default.isNumber(InputTreatyYear.TreatyYear)` b335; `Property-Set-Messages` b226 `pyMessageLabel CheckYearly` b268 — teks pesan TIDAK diekspor | `ErrTahunTreatyBukanAngka` (kosakata kami; teks Pega `[terbuka]`) |
| `ReportDefinition/BrowseTreatyGroup_RD.xml` | kelas `ASM-FW-GISFW-Int-TREATYGROUP` b40; `.ID` sort DESC b587, `.TreatyGroupName` b602, `.TreatyGroupSOAName` b618, `.OJKBusinessID` b632, `.OJKBusinessName` b647, `.OJKBusinessNameIDN` b663; saringan `.ID = Param.ID` b546–b555; maks 500 b688. Tabel fisik `POOLDATA.TREATYGROUP` (`ID`, `TREATYGROUPNAME`, `OLDID` — `NB Treaty In/RDBList/FetchTreatyGroupOLDID.xml`) | `GET /api/treaty-contract-out/grup-treaty` (`ID`, `TREATYGROUPNAME`, `ORDER BY ID DESC`), dibaca saja; kosong → 503 |
| `PEGA_TREATYYEAR` `[data DBA]` | upsert dikunci ID; `'1'‖lpad(TreatyYear_seq,6)`; TGLUPDATE SYSDATE; tidak COMMIT | `Sisip` (ID dari `SEQ_T_TREATYYEAR`) / `Perbarui` (seluruh kolom); satu transaksi + jejak `T_TREATYCO_JEJAK` |

### Gerbang tambahan sadar `[keputusan work owner]`

- **AC 9** — `models.PeriksaPeriodeTCO`: `EndDate` mendahului `StartDate` → 422, pesan menyebut kedua medan dan nilainya.
- **AC 73** — `repository.MasterTahunTreaty.CariDobel` di DALAM transaksi: `TREATYGROUPID = :1 AND TRUNC(STARTDATE) = TRUNC(:2)
  AND TRUNC(ENDDATE) = TRUNC(:3) AND (:4 IS NULL OR ID <> :4)` → 409 `GalatTahunTreatyDobel` menyebut ID baris lain;
  memperbarui baris itu sendiri tetap lolos. Bukan unique index: data warisan boleh sudah berduplikat.

### Ralat bertanggal 28-09-2026

1. **Brief §2** (*"kelompok Treaty Contract Out sudah ada di sidebar"*) — belum ada; ditambahkan ADITIF di tiket
   ini: `labels.ts` `MODUL.treatyContractOut`, `Shell.tsx` kelompok ke-18, `daftarMenu.ts` butir `tco-tahun`
   berlabel VERBATIM `InboxTreatyContract` (b151). Uji cacah bersama diperbarui (17→18 kelompok, 4→5 butir;
   14 kelompok tanpa butir tetap). Dua folder korpus lain (`Treaty In`, `Treaty In Adjustment`) tidak disentuh.
2. **Tiket ini, bab Area codebase** menyebut "Layar daftar + form tahun treaty" tanpa pemilih grup — XML: `Treaty
   Group` adalah pemilih dari `BrowseTreatyGroup_RD` (b6624/b7792) atas master `TREATYGROUP`; dibangun sebagai
   pembaca read-only kecil (`tco_gruptreaty.go`), bukan tiket baru.
3. **OQ-TCO-04 diperkuat**: kolom `PROPORTION` MENYIMPAN pilihan "Reinsurance Type" (form b6829, grid b19724);
   `InputTreatyYear.ReinsuranceType` (form kedua b7954) bukan parameter `PEGA_TREATYYEAR`. Layar memakai
   `PilihJenisReasuransi` (tiket 02) untuk medan itu — sumber daftar pilihan Pega (`pyListSource associated`
   b6857/b7976) tidak diekspor `[dugaan]`.
4. **OQ-TCO-05** ditampilkan di layar: label tahun bersilang antara grid dan form; keduanya VERBATIM.
5. **`daftarMenu.sinkron.test.ts`** dipersempit ke maksudnya: pencocokan palet berbasis potongan kata membuat
   kueri `NB Treaty In` menemukan `InboxTreatyContract` (`nb`/`in` di "inbox", `treaty` di kelompok) — butir yang
   sah; yang dijaga kini "tidak satu pun hasil BERKELOMPOK modul yang belum dimigrasi".

### Yang dibangun

`models/tco_tahun.go` (gerbang), `repository/tco_tahun.go` (daftar/ambil/sisip/perbarui/dobel), `tco_jejak.go`,
`tco_gruptreaty.go`, `services/tco_tahun.go` (`TahunTreatyTCO`: gudang + transaksi disuntik), `tco_gruptreaty.go`,
`handlers/rute_treaty_contract_out.go` (+4 rute), `handlers/tco_db_test.go` (lingkaran penuh HTTP: POST → ID
'1'+6 digit, daftar ID DESC, PUT menimpa tanpa baris baru, 409 dobel menyebut baris lain, 422 periode terbalik,
jejak 3 catatan, 404), frontend `pages/treaty-contract-out/InboxTreatyContract.tsx` (+uji), label `TAHUN_TCO`
(29 label diuji terhadap korpus), `api.ts` (+5 fungsi/antarmuka), menu + rute App.

**Status:** selesai 28-09-2026 — commit `treaty-contract-out: tiket 03 — tahun treaty, periode, anti-dobel`.

### Ralat bertanggal 29-09-2026 — tco5 menu satu butir

**Keputusan work owner tco5**: *"untuk menu hanya Treaty Contract Out; Treaty Contract ReinsType dan Treaty Contract
Description dihapus"*. Kelompok sidebar **Treaty Contract Out** memuat SATU butir berlabel **`Treaty Contract Out`**
(`MENU_TCO.treatyContractOut`; asalnya harness portal `InboxTreatyContract` b151). Layar tahun treaty ini pintu
masuk tunggal modul; kontrak dan klausul dibuka dari tombol `ReinsType` b20778 dan `List Description` b22196 di
barisnya, seperti Pega.
