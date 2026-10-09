# Paritas layar dan aksi — R/I Risk

Setiap tombol / medan / aktivitas di kedua XML Pega (`D:\NUSARE DEV\Menu RI Risk\InboxSummaryRIRisk.xml` = **S**,
`InboxRIRisk.xml` = **D**) → endpoint / komponen, termasuk yang SENGAJA tidak dibuat. Isi Activity / Data Transform
tidak ada di XML (hanya dirujuk); aturannya ASUMSI di `MODUL.md` bab Asumsi. Nomor = baris XML.

## Ringkasan — section `InboxSummaryRIRisk` (kelas `ASM-FW-GISFW-Int-RIRISK_LIFE_SUMMARY`)

| XML | Di Pega | Implementasi |
| --- | --- | --- |
| judul "R/I RISK SUMMARY" S375 | tampil | `RIRiskLife.tsx` `RK.judul` |
| `OutputParam.ERRMSG` S903 (pxDisplayText, NOTBLANK) | pesan galat | galat API ditampilkan `Gagal` |
| `.ID` S1073 (label "ID" S1042, disabled) | tampil, tidak dapat diisi | form `RIRiskLife.tsx` (kosong = dibentuk server) |
| `.USEDBY` S1256 (label "R/I RISK NAME" S1225, wajib S1274) | input | form `RIRiskLife.tsx`; wajib, dipangkas, tidak kembar tanpa beda huruf (`IX_RIRISK_LIFE_SUMMARY_NAMA`, 937) |
| Save S1796 → `AddToListSummary_Act` S1819 | **TERSEMBUNYI**: wadah `pyContainerVisibleWhen` `1=2` S1511, `pyIsVisibilityOption` CONDITION S1546 | **DITAMPILKAN atas keputusan work owner 08-10-2026** (sama dengan ricommlife / riratelife): `POST /api/ri-risk-life` (201, ID = site \|\| LPAD(M_RIRISK_LIFE_SUMMARY_SEQ, 6)), `PUT /api/ri-risk-life/{id}` (200, nama ikut ke rincian - satu transaksi); View only 403 |
| Cancel S2063 → `NewDataSummary_DT` S2091 (`DATASHOW = 'IsEdit'` S2231) | **TERSEMBUNYI** (wadah yang sama) | **DITAMPILKAN atas keputusan work owner 08-10-2026**, hanya saat Edit |
| Upload CSV S3156 (local action) | terlihat (`1=1` S3549) | `UnggahCSV.tsx` (pilih berkas, baca di peramban) |
| View Upload S3676 → harness `ViewCSVResult_RIRisk` S3695 | terlihat | `POST /api/ri-risk-life/unggah/pratinjau`, `HasilUnggah.tsx` |
| Simpan Upload S4690 → `SubmitRIRisk_Act` S4714 | terlihat | `POST /api/ri-risk-life/unggah` (satu transaksi; nama baru = ringkasan baru) |
| label `Format excel : USEDBY, CONTRACT, YEAR, MONTH, RISK` S5495 | **TERSEMBUNYI** (`1=2` S5263) | tidak dirender; urutan kolomnya = kepala CSV `models.KolomCSV` |
| grid `BrowseRIRiskSummary` S7127: ID S7001 · R/I RISK NAME S7147 · MODIFY OPERATOR S7295 · MODIFY DATE S7396; sort ID ASC S9806; `pyPageSize` 50 S10030 | grid | `GET /api/ri-risk-life` (filter ID / nama, urut per kolom, 50 baris) |
| Edit S8589 → `EditListSummary_DT` S8616 (ID, UsedBy) | isi form | `RIRiskLife.tsx` mengisi form; Save menyimpan (`PUT`) |
| Detail S8869 → `setIDUsedBy_Act` S8886 + harness `InboxRIRisk` S8923 (`pyWindowName` "Ri Comm" S8917 - salin-tempel Pega) | popup | `RincianDetail.tsx`, `GET /{id}/detail` |
| Delete S9624 → `DeleteSummaryDetail` S9648 (DeleteID=.ID S9664) | hapus ringkasan + rincian | `DELETE /api/ri-risk-life/{id}` satu transaksi, `KonfirmasiHapus.tsx` |
| grid `BrowseRIRiskLife_RD` S6978 di section ringkasan | tidak berlabel / tidak dipakai layar | tidak dibuat (grid rincian ada di D) |

## Rincian — section `InboxRIRisk` (kelas `ASM-FW-GISFW-Int-RI_RISK_LIFE`)

| XML | Di Pega | Implementasi |
| --- | --- | --- |
| judul "R/I RISK DETAIL" D382 | tampil | `RK.judulRincian` |
| Clear Field D1125 → `clearInputFieldRIRisk_act` D1148 | kosongkan form | tombol `RK.bersihkan` |
| `.ID` D1347 disabled | tampil | form, kosong = dibuat server (`'1' || LPAD(seq, 5)`) |
| `TempIDUsedBy.ID` D1524 (`1=2` D1643) | tersembunyi | IDUSEDBY dari ringkasan di server |
| `TempIDUsedBy.USEDBY` D1727 (label "R/I RISK NAME" D1696, disabled, wajib) | tampil | nama ringkasan, bukan isian (badan berisi `usedby` = 400) |
| `.CONTRACT` D1920 (pxTextInput, wajib D1933, placeholder 0, `ChangeDotToPoint_DT` D1963) | input wajib | `models.NormalContract` (bulat, ≤ 10 angka) |
| `.YEAR` D2202 (TIDAK wajib D2216, pyMax 4 D2220) | input | `models.NormalYear` (kosong / bulat ≤ 4 angka) |
| `.MONTH` D2389 (TIDAK wajib, pyMax 4 D2407; label D2358 bertuliskan "YEAR", preview "MONTH" D2381) | input | `models.NormalMonth`; label "MONTH" (= grid D8277) - **keputusan work owner 08-10-2026: label tetap "MONTH"** |
| `.RISK` D2579 (pxNumber D2582, wajib D2595, label "RISK (PERMIL)" D2548) | input wajib | `models.NormalRisk` (tidak negatif, koma / titik, ≤ 38 angka) |
| Save D3132 → `AddToList_Act` D3155 | tambah / ubah | `POST /{id}/detail`, `PUT /{id}/detail/{detailId}` |
| Cancel D3399 → `NewRIRiskLife_Act` D3423 (`DATASHOW = 'IsEdit'` D3565) | batal edit | tombol `RK.batal` hanya saat EDIT |
| Upload CSV / View Upload / Simpan Upload D4494 / D5014 / D6028 | **TERSEMBUNYI** (wadah `1=2` D4210) | tidak dibuat di popup (ada di halaman ringkasan) |
| grid `BrowseRIRiskLife_RD` D7612 (idusedby D7508): ID D7635 · R/I RISK NAME D7789 · CONTRACT D7956 · YEAR D8123 · MONTH D8277 · RISK (PERMIL) D8423; sort ID ASC D7607; `pyPageSize` Other = 200 D10169/D10243 | grid | `GET /{id}/detail` 200 baris; RISK koma desimal (aturan COMM ricommlife) |
| EDIT D9784 → `EditRIRiskLife_Act` D9808 (ID, Contract, Year, risk, UsedTo, Month) | isi form | tombol `RK.editRincian` mengisi form |
| Delete per baris | **tidak ada** di D | tidak dibuat (baris hilang hanya lewat Delete ringkasan) |

## Perbedaan dari R/I Comm Life (yang wajib tampak di kode)

- Ringkasan: di XML Save / Cancel dan label format TERSEMBUNYI; Save / Cancel DITAMPILKAN atas keputusan work owner
  08-10-2026 (paritas perilaku dengan Comm), label format tetap tidak.
- Rincian: medan dan kolom MONTH; YEAR tidak wajib (Comm: wajib tepat 4 angka); CONTRACT teks (Comm: NUMBER(5));
  label "RISK (PERMIL)"; tombol "EDIT" (Comm "Edit"); Cancel = `NewRIRiskLife_Act` (Comm `NewData_DT`); EDIT =
  `EditRIRiskLife_Act` (Comm `EditList_DT`); grid 200 baris (Comm 50).
- ID rincian `'1' || LPAD(seq, 5, '0')` lebar 6 (Comm: site || LPAD 6).
- Kembar rincian (CONTRACT, YEAR, MONTH) (Comm: CONTRACT, YEAR).
