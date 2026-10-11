# Paritas layar dan aksi — Cause Of Loss Life

Setiap medan / tombol / aktivitas di XML Pega `D:\NUSARE DEV\Menu Cause Of Loss\InboxCauseofLossLife.xml` (section
`InboxCauseofLossLife`, kelas `ASM-FW-GISFW-Int-CAUSEOFLOSS_LIFE`; satu-satunya XML) → endpoint / komponen, termasuk yang
SENGAJA tidak dibuat. Nomor `bNNNN` = baris XML (dibaca utuh, 6.568 baris). Isi Activity / Data Transform / Report
Definition TIDAK ada di XML (hanya dirujuk); aturan yang bergantung padanya ditandai.

## Section `InboxCauseofLossLife`

| XML | Di Pega | Implementasi |
| --- | --- | --- |
| judul "CAUSE OF LOSS" b343 (`pyCaption CAUSE OF LOSS` b6533) | tampil | `CauseOfLossLife.tsx` `COL.judul` |
| "ID" b823 / b835 (`.ID` b852, pxTextInput b855, `pyRequired` true b817 / b870, `pyDisabled` true b872, `pyDisabledNew` always b871, lebar 200px b865) | tampil, TIDAK dapat diisi | input `readOnly disabled`; kosong = dibentuk server `'1' \|\| LPAD(M_CAUSEOFLOSS_LIFE_SEQ, 5)` (K3); badan berisi `id` = 400 |
| "Cause of Loss" b1001 / b1013 (`.CauseofLoss` b1027, pxTextInput b1029, `pyRequired` true b996 / b1044, `pyRequiredNew` always b1037, lebar 200px b1039) | input teks wajib | `<input type="text" required>`; `models.NormalCauseOfLoss` (wajib → 422 "Cause of Loss is required") |
| pengubah huruf (data transform / refresh pada medan Cause of Loss) | **tidak ada** (nol `pyBehaviors` / `SetUpperCase*` di b967-b1133) | huruf TIDAK diubah (K4) - dibedakan dari Benefit yang punya `SetUpperCase_DT` |
| `InputParam.ERRMSG` b2018 (pxDisplayText b2021, read-only b1968, tampil bila NOTBLANK b2105) | pesan galat simpan | galat API ditampilkan `Gagal` di atas form (422 / 403 / 409 / 503 berkalimat) |
| Save b5314 / b5366 (pxButton, refresh `AddToList_Act` b5385 / b5464) | simpan | `POST /api/cause-of-loss-life` (201, ID baru) bila form tanpa ID; `PUT /api/cause-of-loss-life/{id}` (200) sesudah Edit. Satu transaksi. View only 403 |
| Cancel b5574 / b5624 (`NewData_DT` b5647 / b5754), terlihat hanya bila `InputParam.DATASHOW = 'IsEdit'` b5784-b5787 | batal edit | tombol `COL.batal` HANYA saat Edit; mengosongkan form (murni layar, nol endpoint) |
| grid `BrowseCauseofLossLife_RD` b2968 / b3981 (`pgRepPgSubSectionInboxCauseofLossLifeBBBBB.pxResults` b2816, `pyRDAppliesTo` b4034) | grid | `GET /api/cause-of-loss-life`; `SqlDaftar` atas `CAUSEOFLOSS_LIFE` |
| kolom "ID" b2989 = `.ID` b3415 | tampil | kolom `COL.id` = `ID` |
| kolom "Cause of Loss" b3127 = `.CauseofLoss` b3562 | tampil | kolom `COL.causeOfLoss` = `CAUSEOFLOSS` (100001 tampil kosong) |
| kolom ketiga tanpa judul b3245 berisi Edit b3705 / b3747 (`EditList_DT` b3770 dengan ID b3774-b3776 dan CauseofLoss b3780-b3781) | isi form | tombol `COL.edit` per baris mengisi form (ID + Cause of Loss) → mode Edit |
| sort ID `pySortType` ASC b3923, `pySortOrder` 1 b3929; ketiga kolom `pyColumnSorting` false b3925 / b3947 / b3969 | ID menaik, tanpa urut pilihan | `ORDER BY` ID ANGKA menaik tetap; kepala kolom bukan tombol |
| `pyGridFiltering` false b4036 | tanpa saring | **tidak dibuat** - nol kotak saring, nol parameter saring di API |
| `pyPageSize` 10 b4057, `pyPageMode` Numeric b4028, `pyGridPaginator` b2724 / b2746 | halaman 10 | 10 baris per halaman, tombol Previous / Next + "Page n of m" |
| `pyRowEditing` / `pyEditingMode` readOnly b3989 / b4007 | grid baca-saja | grid baca-saja; ubah lewat form |
| Delete (`pyGridDeleteActivityExists` false b431 / b725 / b1605 / b1890 / b2598 / b2847 / b4185 / b5192) | **tidak ada** | **tidak dibuat** - nol rute DELETE (uji `TestRute`), nol SQL DELETE (`TestNolJSONDanNamaLama`) |
| Upload CSV / View Upload / ringkasan / detail | **tidak ada** di XML | **tidak dibuat** |
| `pyGridNoResultsMessage` b4054 / b6293 | pesan grid kosong | `Kosong` "No cause of loss yet." |

## Di luar XML, mengikuti pola modul MASTER TREATY (WO boleh menolaknya - K4)

| Aturan | Alasan | Bila ditolak WO |
| --- | --- | --- |
| Cause of Loss dipangkas spasi tepinya sebelum disimpan | riratelife / ricommlife / ririsklife / benefitlife / planlife memangkas nama; isi `AddToList_Act` tidak ada | buang `strings.TrimSpace` di `models.NormalCauseOfLoss` |
| Cause of Loss TIDAK BOLEH KEMBAR tanpa beda huruf dan spasi tepi (422 "… is already used by ID …"); Edit boleh menyimpan namanya sendiri; baris kosong (100001) tidak pernah dianggap kembar | pola planlife K5 / ririsklife; data DEV tanpa kembar (fakta WO) | buang pemeriksaan `PemakaiNama` di `services.Simpan` |
| Cause of Loss paling panjang 200 byte (422 berkalimat) | lebar kolom `CAUSEOFLOSS` VARCHAR2(200) (K1); tanpa ini Oracle menjawab ORA-12899 = 500 mentah | tetap perlu selama kolom 200 - hanya pesannya |
| Halaman Previous / Next | wujud web paginator Pega Numeric | - |
| Menu View only: form dan Edit disembunyikan, tulis 403 | hak menu Full / View only aplikasi (pola semua modul MASTER TREATY) | - |
| ID dari sequence yang sudah ada = 409 berkalimat (bukan lompat) | keputusan work owner K3 | - |

**TIDAK diterapkan**: huruf besar (XML tanpa pengubah huruf - K4), saring dan urut pilihan (XML mematikannya).

## Pertanyaan terbuka

- **T1 - Save untuk Add dan Edit**: satu tombol Save → `AddToList_Act` untuk keduanya; isi aktivitas tidak ada.
  Diterapkan: form tanpa ID = Add, sesudah Edit = ubah baris itu. **Rekomendasi: pertahankan** (satu-satunya bacaan
  yang membuat Edit + Cancel `IsEdit` bermakna; sama dengan benefitlife T2).
- **T2 - baris 100001 (kosong)**: tidak dipakai produk (fakta WO); migrasi tidak mengisi / menghapusnya. Di layar ia
  tampil kosong dan dapat diisi lewat Edit. **Rekomendasi: biarkan**; penghapusan (bila diinginkan) adalah keputusan WO
  terpisah - XML tanpa Delete.
