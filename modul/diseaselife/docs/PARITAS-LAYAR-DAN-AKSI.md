# Paritas layar dan aksi — Disease Life

Setiap medan / tombol / aktivitas di XML Pega `D:\NUSARE DEV\Menu Disease\InboxDisease.xml` (section `InboxDisease`,
kelas `ASM-FW-GISFW-Int-DISEASE_LIFE`; satu-satunya XML) → endpoint / komponen, termasuk yang SENGAJA tidak dibuat.
Nomor `bNNNN` = baris XML (dibaca utuh, 6.451 baris). Isi Activity / Data Transform / Report Definition TIDAK ada di XML
(hanya dirujuk); aturan yang bergantung padanya ditandai.

## Section `InboxDisease`

| XML | Di Pega | Implementasi |
| --- | --- | --- |
| judul "DISEASE" b371 (`pyCaption DISEASE` b6223) | tampil | `DiseaseLife.tsx` `DSL.judul` |
| `InputParam.ERRMSG` b881 (pxDisplayText b884, read-only b829, tampil bila NOTBLANK b966) | pesan galat simpan | galat API ditampilkan `Gagal` di atas form (422 / 403 / 409 / 503 berkalimat) |
| "Number / ID" b1021 / b1044 (`DISEASE_LIFE.Number` b1050, pxTextInput b1053, `pyDisabled` true b1070, `pyDisabledNew` always b1069, lebar 105px b1063) | tampil, TIDAK dapat diisi | input `readOnly disabled`; kosong = dibentuk server `TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL)` (D1.1); badan berisi `id` = 400 |
| "ICD Code" b1200 / b1223 (`.ICD_Code` b1229, pxTextInput b1232, `pyRequired` **false** b1195 / b1244, `pyRequiredNew` false b1237, lebar 105px b1239) | input teks | `<input type="text" required>`; wajib = **keputusan WO D3, bukan XML** (lihat "Di luar XML") |
| perubahan ICD Code → refresh + `SetUpperCase_DT` (behaviors b1246-b1283, `pyName` b1265; ActionSets b1353-b1406, b1386) | huruf besar | `hurufBesar` saat medan ditinggalkan (`onBlur`) + `models.NormalICDCode` di server |
| "Disease" b1470 / b1492 (`.Disease` b1499, pxTextArea b1502, `pyRequired` true b1465 / b1515, `pyRequiredNew` always b1511, tinggi minimal 3 b1510) | textarea wajib | `<textarea rows={3} required>`; `models.NormalDisease` (wajib → 422 "Disease is required") |
| perubahan Disease → refresh + `SetUpperCase_DT` (behaviors b1519-b1556, b1538; ActionSets b1622-b1675, b1654) | huruf besar | `hurufBesar` `onBlur` + server. Isi data transform tidak ada di XML - arti nama diterapkan (preseden benefitlife T1) |
| Save b2050 / b2102 (pxButton b2059, refresh `AddToList_Act` b2121 / b2202) | simpan | `POST /api/disease-life` (201, ID baru) bila form tanpa ID; `PUT /api/disease-life/{id}` (200) sesudah Edit. Satu transaksi. View only 403 |
| Cancel b2313 / b2365 (`NewData_DT` b2388 / b2492), terlihat hanya bila `InputParam.DATASHOW = 'IsEdit'` b2522-b2525 | batal edit | tombol `DSL.batal` HANYA saat Edit; mengosongkan form (murni layar, nol endpoint) |
| grid `BrowseDiseaseLife_RD` b5094 (`pgRepPgSubSectionInboxDiseaseBBBB.pxResults` b3588, `pyRDAppliesTo` b5146, sumber Report Definition b5161) | grid | `GET /api/disease-life`; `SqlDaftar` atas `DISEASE_LIFE` |
| kolom "ID" b3758 (sel salin-tempel dari `BrowseBenefitLife_RD` b3738) = `.Number` b4321 (sel b4301 juga `BrowseBenefitLife_RD`) | tampil | kolom `DSL.id` = `ID` (`.Number` = kolom ID; prosedur Pega menulis ID) |
| kolom "ICD Code" b3894 = `.ICD_Code` b4484 | tampil | kolom `DSL.icdCode` = `ICD_CODE` |
| kolom "Disease" b4033 = `.Disease` b4629 (pxDisplayText b4632, wrap b4564) | tampil | kolom `DSL.disease` = `DISEASE`, dibungkus |
| kolom keempat tanpa judul b4154 berisi Edit b4777 / b4829 (`EditList_DT` b4852 dengan Number b4856, Disease b4862, ICD_Code b4868) | isi form | tombol `DSL.edit` per baris mengisi form (ID + ICD Code + Disease) → mode Edit |
| urut: ID `pySortType` DESC b5012 / `pySortOrder` 1 b5017 / `pyColumnSorting` true b5014; ICD Code `pyColumnSorting` true b5036 (`pySortType` NONE b5034); Disease false b5058; Edit false b5080; `pyGridSorting` true b5138 | bawaan ID menurun; ID / ICD Code dapat diurutkan | bawaan ID ANGKA menurun; kepala ID dan ICD Code = tombol urut (`pilihUrut`), Disease bukan tombol; `urut` / `arah` di API |
| `pyGridFiltering` true b5148 (`pyColumnFiltering` kosong per kolom b3735 / b3871 / b4007) | saring grid | dua kotak saring: ICD Code dan Disease ("memuat", tanpa beda huruf) → `icd` / `disease` di API (D3) |
| `pyPageSize` 10 b5169, `pyPageMode` Numeric b5140, `pyGridPaginator` b3489 / b3511 | halaman 10 | 10 baris per halaman di SERVER (OFFSET / FETCH NEXT), Previous / Next + "Page n of m" |
| `pyRowEditing` / `pyEditingMode` readOnly b5102 / b5119 | grid baca-saja | grid baca-saja; ubah lewat form |
| Delete (`pyGridDeleteActivityExists` false b460 / b754 / b1929 / b3262 / b3619 / b5297) | **tidak ada** | **tidak dibuat** - nol rute DELETE (uji `TestRute`), nol SQL DELETE (`TestNolJSONDanObjekLama`) |
| Upload CSV / View Upload / ringkasan / detail | **tidak ada** di XML | **tidak dibuat** |
| `pyGridNoResultsMessage` b5166 / b6091 | pesan grid kosong | `Kosong` "No disease yet." / "No disease matches the filter." |

## Kebutuhan teknis (bukan fitur baru) — 97.586 baris

Pega memuat grid lewat Report Definition berhalaman (`pyPageSize` 10, Numeric). Di aplikasi ini **saring, urut, dan
halaman SELALU di server**: `SELECT … WHERE <saring> ORDER BY … OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY` + `COUNT(*)`
bersaring. Tidak ada jalan yang memuat seluruh `DISEASE_LIFE` (uji `TestSqlDaftarSelaluBerhalaman`, vitest "97.586
baris"). Urut ID memakai `TO_NUMBER(REGEXP_SUBSTR(ID, …))` (ID teks) - 97.586 baris diurut per permintaan; cukup untuk
satu layar master, dicatat sebagai pertanyaan bila WO melihatnya lambat.

## Di luar XML, mengikuti pola modul MASTER TREATY / keputusan WO (WO boleh menolaknya)

| Aturan | Alasan | Bila ditolak WO |
| --- | --- | --- |
| **ICD Code WAJIB** (422 "ICD Code is required") | keputusan WO D3 ("ICD Code dan Disease wajib (XML)") - **temuan: XML menandai ICD Code TIDAK wajib** (`pyRequired` false b1195 / b1244, `pyRequiredNew` false b1237); dipertahankan karena ICD Code dipakai pencarian diagnosa claimlife dan kunci keunikan | buang cabang kosong di `models.NormalICDCode` dan `periksaIsian` (keunikan tetap untuk ICD terisi) |
| ICD Code TIDAK BOLEH KEMBAR tanpa beda huruf dan spasi tepi (422 "ICD Code … is already used by ID …"); Edit boleh menyimpan ICD-nya sendiri | D3; data DEV 97.586 ICD unik (fakta WO) | buang pemeriksaan `PemakaiICD` di `services.Simpan` |
| ICD Code / Disease dipangkas spasi tepinya | pola modul MASTER TREATY; isi `AddToList_Act` tidak ada | buang `strings.TrimSpace` di `models.Normal*` |
| panjang paling banyak lebar kolom (ICD 100 byte, Disease 1000 byte) | lebar kolom warisan; tanpa ini Oracle menjawab ORA-12899 = 500 | tetap perlu selama kolom selebar itu - hanya pesannya |
| ID baru = `TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL)`, BUKAN rumus Pega `'1' \|\| LPAD(M_DISEASE_LIFE_SEQ, 5, '0')` | D1.1: rumus Pega menghasilkan 102051, 102052, … yang SUDAH dipakai data impor (asal ID kembar 102051) | - |
| ID dari sequence yang sudah ada = 409 berkalimat (bukan lompat) | D1.1 ("aplikasi tetap memeriksa ID belum terpakai"); prosedur Pega masih VALID dan dapat menulis | - |
| Halaman Previous / Next; menu View only: form dan Edit disembunyikan, tulis 403 | wujud web paginator Numeric; hak menu aplikasi | - |

## Pertanyaan terbuka

- **T1 - Save untuk Add dan Edit**: satu tombol Save → `AddToList_Act` untuk keduanya; isi aktivitas tidak ada.
  Diterapkan: form tanpa ID = Add, sesudah Edit = ubah baris itu. **Rekomendasi: pertahankan** (sama dengan benefitlife
  T2 / causeoflosslife T1).
- **T2 - `SetUpperCase_DT`**: isinya tidak ada di XML; diterapkan pada KEDUA medan yang memanggilnya. Baris uji
  `TEST123 / Sakit` (huruf kecil) dibuat lewat prosedur, bukan layar - tidak membuktikan sebaliknya. **Rekomendasi:
  pertahankan** (data DEV seluruhnya huruf besar).
- **T3 - ICD Code wajib** (lihat "Di luar XML"): XML tidak mewajibkannya. **Rekomendasi: pertahankan wajib** (keputusan
  D3; ICD Code kosong tidak dapat dicari claimlife dan tidak dapat dijaga keunikannya).
