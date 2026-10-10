# Paritas layar dan aksi — Cover Life

Setiap medan / tombol / aktivitas di XML Pega `D:\NUSARE DEV\Menu Cover\InboxCoverLife.xml` (section `InboxCoverLife`,
kelas `ASM-FW-GISFW-Int-COVER_LIFE`; satu-satunya XML) → endpoint / komponen, termasuk yang SENGAJA tidak dibuat. Nomor
`bNNNN` = baris XML (dibaca utuh, 6.685 baris). Isi Activity / Data Transform / Report Definition TIDAK ada di XML (hanya
dirujuk); aturan yang bergantung padanya ditandai.

## Section `InboxCoverLife`

| XML | Di Pega | Implementasi |
| --- | --- | --- |
| judul "COVER" b368 | tampil | `CoverLife.tsx` `CVL.judul` |
| medan ID di form | **tidak ada** (form S2 hanya dua sel: Cover b848, Note b1026) | **tidak dibuat** - ID dibentuk server `'1' \|\| LPAD(M_COVER_LIFE_SEQ, 5)` (C2); saat Edit ID tampil sebagai teks status "Editing ID …" (teks aplikasi, bukan medan); badan berisi `id` = 400 |
| "Cover" b848 / b871 (`.Cover` b877, pxTextInput b880, `pyRequired` true b843 / b895, `pyRequiredNew` always b888, lebar 300px b890) | input teks wajib | `<input type="text" required>`; `models.NormalCover` (wajib → 422 "Cover is required") |
| "Note" b1026 / b1048 (`.Note` b1055, pxTextArea b1058, `pyRequired` false b1020 / b1070, `pyRequiredNew` false b1063, lebar 300px b1065) | textarea opsional | `<textarea>` tanpa `required`; `models.NormalNote` (hanya spasi = kosong / NULL) |
| pengubah huruf (data transform / refresh pada Cover / Note) | **tidak ada** (`pyBehaviors` kosong b897 / b1072; nol `SetUpperCase*`) | huruf TIDAK diubah (C3) - dibedakan dari Disease yang punya `SetUpperCase_DT` |
| `InputParam.ERRMSG` b2052 (pxDisplayText b2055, read-only, tampil bila NOTBLANK b2138) | pesan galat simpan | galat API ditampilkan `Gagal` di atas form (422 / 403 / 409 / 503 berkalimat) |
| Save b5357 / b5407 (pxButton b5367, refresh `AddToList_Act` b5426 / b5508) | simpan | `POST /api/cover-life` (201, ID baru) bila form tanpa ID; `PUT /api/cover-life/{id}` (200) sesudah Edit. Satu transaksi. View only 403 |
| Cancel b5618 / b5671 (`NewData_DT` b5694 / b5798), terlihat hanya bila `InputParam.DATASHOW = 'IsEdit'` b5829-b5831 | batal edit | tombol `CVL.batal` HANYA saat Edit; mengosongkan form (murni layar, nol endpoint) |
| grid `BrowseCoverLife_RD` b4026 (`pgRepPgSubSectionInboxCoverLifeBBBBB.pxResults` b2850, `pyRDAppliesTo` b4078) | grid | `GET /api/cover-life`; `SqlDaftar` atas `M_COVER_LIFE` |
| kolom "ID" b3023 (sel salin-tempel dari `BrowseCauseofLossLife_RD` b3002) = `.ID` b3449 | tampil | kolom `CVL.id` = `ID` |
| kolom "Cover" b3161 = `.Cover` b3596 | tampil | kolom `CVL.cover` = `COVER` |
| kolom ketiga tanpa judul b3277 berisi Edit b3732 / b3784 (`EditList_DT` b3807 dengan ID b3811, Cover b3817, Note b3823) | isi form | tombol `CVL.edit` per baris mengisi form (ID + Cover + Note) → mode Edit; Note dikirim API walau tidak tampil di grid |
| Note di grid | **tidak ada** (hanya tiga kolom: ID, Cover, Edit) | **tidak ditampilkan** |
| sort ID `pySortType` ASC b3967, `pySortOrder` 1 b3972; ketiga kolom `pyColumnSorting` false b3969 / b3991 / b4013 | ID menaik, tanpa urut pilihan | `ORDER BY` ID ANGKA menaik tetap; kepala kolom bukan tombol |
| `pyGridFiltering` false b4080 | tanpa saring | **tidak dibuat** - nol kotak saring, nol parameter saring di API |
| `pyPageSize` **50** b4101, `pyPageMode` Numeric b4072, `pyGridPaginator` b2758 / b2780 | halaman 50 | 50 baris per halaman, tombol Previous / Next + "Page n of m" |
| `pyRowEditing` / `pyEditingMode` readOnly b4034 / b4051 | grid baca-saja | grid baca-saja; ubah lewat form |
| Delete (`pyGridDeleteActivityExists` false b457 / b749 / b1638 / b1925 / b2633 / b2882 / b4229 / b5236) | **tidak ada** | **tidak dibuat** - nol rute DELETE (uji `TestRute`), nol SQL DELETE (`TestKodeTidakMenyebutViewLama`) |
| Upload CSV / View Upload / ringkasan / detail | **tidak ada** di XML | **tidak dibuat** |
| `pyGridNoResultsMessage` b4098 / b6385 | pesan grid kosong | `Kosong` "No cover yet." |

## Di luar XML, mengikuti pola modul MASTER TREATY (WO boleh menolaknya - C3)

| Aturan | Alasan | Bila ditolak WO |
| --- | --- | --- |
| Cover dipangkas spasi tepinya sebelum disimpan | pola modul MASTER TREATY; isi `AddToList_Act` tidak ada | buang `strings.TrimSpace` di `models.NormalCover` |
| Cover TIDAK BOLEH KEMBAR tanpa beda huruf dan spasi tepi (422 "… is already used by ID …"); Edit boleh menyimpan namanya sendiri | C3; data DEV tanpa kembar (fakta WO) | buang pemeriksaan `PemakaiNama` di `services.Simpan` |
| Cover paling panjang 200 byte, Note 1000 byte (422 berkalimat) | lebar kolom 085 (C1.1); tanpa ini Oracle menjawab ORA-12899 = 500 | tetap perlu selama kolom selebar itu - hanya pesannya |
| Note yang hanya spasi disimpan kosong (NULL) | 4/4 baris DEV NULL; menghindari Note "kosong" yang bukan NULL | simpan apa adanya |
| Halaman Previous / Next; menu View only: form dan Edit disembunyikan, tulis 403 | wujud web paginator Numeric; hak menu aplikasi | - |
| ID dari sequence yang sudah ada = 409 berkalimat (bukan lompat) | keputusan work owner C2 | - |

**TIDAK diterapkan**: huruf besar (XML tanpa pengubah huruf - C3), saring dan urut pilihan (XML mematikannya), medan ID
di form (XML tidak memuatnya).

## Pertanyaan terbuka

- **T1 - Save untuk Add dan Edit**: satu tombol Save → `AddToList_Act` untuk keduanya; isi aktivitas tidak ada.
  Diterapkan: form tanpa ID = Add, sesudah Edit = ubah baris itu. **Rekomendasi: pertahankan** (sama dengan
  causeoflosslife T1).
- **T2 - teks "Editing ID …"**: form Pega tidak menampilkan ID saat Edit; aplikasi menampilkan teks status kecil supaya
  pengguna tahu baris mana yang diubah (bukan medan, tidak dapat diisi). **Rekomendasi: pertahankan**; dibuang bila WO
  ingin paritas visual mutlak.
