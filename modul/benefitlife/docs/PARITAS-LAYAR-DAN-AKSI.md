# Paritas layar dan aksi — Benefit

Setiap medan / tombol / aktivitas di XML Pega `D:\NUSARE DEV\Menu Benefit\InboxBenefit.xml` (section `InboxBenefit`,
kelas `ASM-FW-GISFW-Int-BENEFIT_LIFE`; satu-satunya XML) → endpoint / komponen, termasuk yang SENGAJA tidak dibuat.
Nomor `bNNN` = baris XML (dibaca utuh, 5.783 baris). Isi Activity / Data Transform / Report Definition TIDAK ada di XML
(hanya dirujuk); aturan yang bergantung padanya ditandai.

## Section `InboxBenefit`

| XML | Di Pega | Implementasi |
| --- | --- | --- |
| judul "INSURANCE BENEFIT" b313 (`pyCaption INSURANCE BENEFIT` b5616) | tampil | `BenefitLife.tsx` `BN.judul` |
| `InputParam.ERRMSG` b823 (pxDisplayText b826, read-only b781, tampil bila NOTBLANK b908) | pesan galat simpan | galat API ditampilkan `Gagal` di atas form (422 / 403 / 409 / 503 berkalimat) |
| "Number / ID" b964 / b987 (`.Number` b993, pxTextInput b996, `pyDisabled` true b1013, `pyDisabledNew` always b1012, lebar 105px b1006) | tampil, TIDAK dapat diisi | input `readOnly disabled`; kosong = dibentuk server `'1' \|\| LPAD(M_BENEFIT_LIFE_SEQ, 5)` (K3); badan berisi `id` = 400 |
| `.Number` = ID | - | **keputusan work owner**: tidak ada kolom NUMBER (Number di 5 baris JSON, selalu = ID) |
| "Benefit" b1143 / b1166 (`.Benefit` b1172, pxTextArea b1175, `pyRequired` true b1137 / b1188, `pyRequiredNew` always b1185) | textarea wajib | `<textarea required>`; `models.NormalBenefit` (wajib → 422 "Benefit is required") |
| perubahan Benefit → refresh dengan `SetUpperCase_DT` b1192-b1211 (ActionSets b1295-b1327) | nilai jadi huruf besar | frontend: huruf besar saat medan ditinggalkan (`onBlur`, peristiwa `change` b1199); backend: `strings.ToUpper` sebelum simpan. Isi DT tidak ada di XML - **T1** |
| Save b1723 / b1772 (pxButton, refresh `AddToList_Act` b1791 / b1875) | simpan | `POST /api/benefit-life` (201, ID baru) bila form tanpa ID; `PUT /api/benefit-life/{id}` (200) sesudah Edit. Satu transaksi. View only 403 |
| Cancel b1985 / b2035 (`NewData_DT` b2058 / b2165), terlihat hanya bila `InputParam.DATASHOW = 'IsEdit'` b2195-b2198 | batal edit | tombol `BN.batal` HANYA saat Edit; mengosongkan form (murni layar, nol endpoint) |
| grid `BrowseBenefitLife_RD` b3411 / b4451 (`pgRepPgSubSectionInboxBenefitBBBB.pxResults` b3261, `pyRDAppliesTo` b4503) | grid | `GET /api/benefit-life`; `SqlDaftar` atas `BENEFIT_LIFE` |
| kolom "ID" b3431 = `.Number` b3858 (pxDisplayText) | tampil | kolom `BN.id` = `ID` |
| kolom "Benefit" b3570 = `.Benefit` b4020 (pxDisplayText, wrap b3956) | tampil | kolom `BN.benefit` = `BENEFIT`, teks dibungkus |
| kolom ketiga tanpa judul b3708 berisi Edit b4168 / b4220 (`EditList_DT` b4243 dengan Number b4247 dan Benefit b4253) | isi form | tombol `BN.edit` per baris mengisi form (ID + Benefit) → mode Edit |
| sort ID `pySortType` DESC b4391, `pySortOrder` 1 b4397; `pyGridSorting` true b4493; kolom Benefit `pyColumnSorting` false b4415 | ID menurun | bawaan ID MENURUN (angka); kepala ID membalik arah (`arah=asc`); Benefit tidak dapat diurutkan |
| `pyGridFiltering` true b4505 | saring kolom | dua kotak saring (ID, Benefit) "memuat" tanpa beda huruf - bentuk saring = pola ririsklife |
| `pyPageSize` 10 b4526, `pyPageMode` Numeric b4498, `pyGridPaginator` b3162 | halaman 10 | 10 baris per halaman, tombol Previous / Next + "Page n of m" |
| `pyRowEditing` / `pyEditingMode` readOnly b4457 / b4476 | grid baca-saja | grid baca-saja; ubah lewat form |
| Delete (`pyGridDeleteActivityExists` false b402 / b696 / b1602 / b2935 / b3292, `pyDeleteActivity` kosong) | **tidak ada** | **tidak dibuat** - nol rute DELETE (uji `TestRute`), nol SQL DELETE (`TestNolJSONDanNamaLama`) |
| Upload CSV / View Upload / ringkasan / detail | **tidak ada** di XML | **tidak dibuat** |
| `pyGridNoResultsMessage` b5448 | pesan grid kosong | `Kosong` "No benefit yet." / "No benefit matches the filter." |

## Di luar XML, mengikuti pola modul MASTER TREATY (WO boleh menolaknya - K6)

| Aturan | Alasan | Bila ditolak WO |
| --- | --- | --- |
| Benefit dipangkas spasi tepinya sebelum disimpan | riratelife / ricommlife / ririsklife memangkas nama; isi `AddToList_Act` tidak ada | buang `strings.TrimSpace` di `models.NormalBenefit` |
| Benefit paling panjang 200 byte (422 berkalimat) | lebar kolom `BENEFIT` VARCHAR2(200) (K1); tanpa ini Oracle menjawab ORA-12899 = 500 mentah | tetap perlu selama kolom 200 - hanya pesannya |
| Saring ID / Benefit "memuat" + halaman Previous / Next | wujud web `pyGridFiltering` / paginator Pega (pola ririsklife) | - |
| Menu View only: form dan Edit disembunyikan, tulis 403 | hak menu Full / View only aplikasi (pola semua modul MASTER TREATY) | - |
| ID dari sequence yang sudah ada = 409 berkalimat (bukan lompat) | keputusan work owner K3 | - |

**TIDAK diterapkan** (walau ada di pola ririsklife): penolakan Benefit kembar - tidak ada di XML (K6).

## Pertanyaan terbuka

- **T1 - `SetUpperCase_DT`**: XML hanya menyebut namanya pada perubahan textarea Benefit (b1211); isinya tidak ada.
  Diterapkan sebagai "Benefit menjadi huruf besar" (frontend saat medan ditinggalkan, backend sebelum simpan). Data
  lama TIDAK diubah migrasi (944 memindahkan apa adanya). `benefit_bukti_k2.sql` kueri G menghitung Benefit DEV yang
  belum huruf besar. **Rekomendasi: pertahankan** (nama DT jelas); bila kueri G > 0 dan WO ingin data lama ikut,
  itu UPDATE terpisah atas keputusan WO - tidak dikerjakan di sini.
- **T2 - Save untuk Add dan Edit**: satu tombol Save → `AddToList_Act` untuk keduanya; isi aktivitas tidak ada.
  Diterapkan: form tanpa ID = Add, sesudah Edit = ubah baris itu. **Rekomendasi: pertahankan** (satu-satunya bacaan
  yang membuat Edit + Cancel `IsEdit` bermakna).
