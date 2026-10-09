# Paritas layar dan aksi — Plan

Setiap medan / tombol / aktivitas / sumber pilihan di XML Pega `D:\NUSARE DEV\Menu Plan\InboxProductType.xml` (section
`InboxProductType`, kelas `ASM-FW-GISFW-Int-PRODUCT_TYPE_LIFE`; satu-satunya XML, dibaca utuh 7.676 baris) → endpoint /
komponen, termasuk yang SENGAJA tidak dibuat. `bNNNN` = baris XML. Isi Activity / Report Definition TIDAK ada di XML
(hanya dirujuk).

## Section `InboxProductType`

| XML | Di Pega | Implementasi |
| --- | --- | --- |
| judul "Plan" b331 (`pyCaption Plan` b7605); judul tersembunyi layout "Product Type" b644 | tampil | `PlanLife.tsx` `PL.judul` |
| `OutputParam.ERRMSG` b5415 (pxDisplayText b5418, read-only b5373, tampil NOTBLANK b5502) | pesan galat | galat API di atas form (`Gagal`) |
| "Plan Name" b812 / b834 (`.CoverName` b841, pxTextInput b844, `pyRequired` false b859); perubahan → refresh b864-b875 tanpa DT | teks | input teks; **K5** wajib + unik (di luar XML) |
| "Business" b1079 (`.Business` b1105, pxAutoComplete b1107, `pyAllowFreeFormInput` true b1125); perubahan → refresh b1132-b1143 | autocomplete | input + daftar pilihan; `GET /api/plan-life/pilihan/business` |
| sumber Business: RD `BrowseBusinessLife_RD` b1183 kelas `ASM-FW-GISFW-Int-BUSINESS` b1189, nilai `.Note` b1187; kolom `.OLDID` b1200, `.Note` b1231 (cari b1230), `.ID` b1265 → `.BusinessID` b1268; parameter `ID` b1296, `Note` b1304, `Group` b1310 (nilai tidak tertulis) | daftar | tabel `BUSINESS`, kolom OLDID / Note / ID; **K4 `GROUPPANEL = '009'` - DITURUNKAN DARI DATA, BUKAN DARI XML** (isi RD dan nilai `Group` tidak ada di XML; ke-21 business yang dipakai plan DEV semuanya di grup ini, grup ini persis 21 baris LIFE INSURANCE / OLDID L1-L21; selaras RD yang sama di XML Retro Life: `.OLDID StartsWith "L"`, `mastercontractretrolife`) |
| "Benefit" b1432 / b1455 (`.Benefit` b1461, pxAutoComplete b1464, `pyAllowFreeFormInput` true b1478); refresh b1485-b1497 | autocomplete | input + daftar pilihan; `GET /api/plan-life/pilihan/benefit` |
| sumber Benefit: RD `BrowseBenefitLife_RD` b1537 kelas `BENEFIT_LIFE` b1543, nilai `.Benefit` b1541 / b1552 (cari b1550), `.Number` → `.BenefitID` b1563-b1564 (tidak tampil b1562) | daftar | tabel `BENEFIT_LIFE` (modul benefitlife), kolom BENEFIT / ID |
| teks bebas (`pyAllowFreeFormInput`) | diterima apa adanya | **K4 keputusan WO**: server MENCOCOKKAN ULANG teks ke master (sama persis tanpa beda huruf / spasi tepi); cocok = nama + ID dari master; tidak cocok / cocok >1 = **422** berkalimat - pasangan nama / ID tidak pernah basi |
| Save b6354 / b6403 (`SaveProductTypeLife_Act` b6422 / b6504) | simpan | `POST /api/plan-life` (201, ID baru) / `PUT /api/plan-life/{id}` (200) sesudah Edit; satu transaksi; View only 403 |
| New b6615 / b6667 (`NewProductTypeLife_act` b6686 / b6790), tampil HANYA bila `InputParam.DATASHOW = 'IsEdit'` b6824-b6827 | kosongkan form | tombol `PL.baru` HANYA saat Edit (murni layar) - RALAT R2 |
| grid `BrowseProductTypeLife_RD` b4325 (`pgRepPgSubSectionInboxProductTypeBBBB.pxResults` b2808, `pyRDAppliesTo` b4378) | grid | `GET /api/plan-life` atas `PRODUCT_TYPE_LIFE` |
| kolom "Plan Name" b2978 (`.CoverName` b3543), "Business" b3117 (`.Business` b3692), "Benefit" b3255 (`.Benefit` b3839), kolom tanpa judul b3371 berisi Edit | tampil | empat kolom; TANPA kolom ID |
| Edit b3975 / b4027 (`EditProductTypeLife_Act` b4045 / b4146 dengan ID b4060, Business b4066, BusinessID b4072, Benefit b4078, BenefitID b4084, CoverName b4090) | isi form | tombol `PL.edit` per baris mengisi form → mode Edit |
| urut: `pyGridSorting` true b4368; Plan Name `pyColumnSorting` true b4247, Business false b4269, Benefit true b4291; `pySortType` NONE, `pyMaxSortOrder` 0 b2890 | urut per kolom, tanpa bawaan | kepala Plan Name / Benefit dapat diklik; bawaan backend ID menaik (stabil - **penyimpangan sadar**, XML tanpa urutan bawaan) |
| `pyGridFiltering` false b4380 | tanpa saring | **tidak dibuat** kotak saring |
| `pyPageSize` 10 b4401, `pyPageMode` Numeric b4373, `pyGridPaginator` b2709 | halaman 10 | 10 baris, Previous / Next + "Page n of m" |
| medan `.ID` di form | **tidak ada** (hanya parameter Edit b4060) | tidak ditampilkan - RALAT R1 (petunjuk prompt menyebut `.ID` tampil) |
| Delete (`pyGridDeleteActivityExists` false b419 / b713 / b2187 / b2839), Upload, ringkasan / detail | **tidak ada** | **tidak dibuat** - nol rute DELETE (uji `TestRute`), nol SQL DELETE |

## Di luar XML (WO boleh menolaknya)

| Aturan | Asal | Bila ditolak WO |
| --- | --- | --- |
| Plan Name WAJIB | K5, pola modul MASTER TREATY | buang cek di `models.PeriksaIsian` |
| Plan Name UNIK tanpa beda huruf (422 "already used by ID …") | K5 (data DEV: 32 nama unik) | buang `PemakaiNama` di `services.Simpan` |
| Business dan Benefit WAJIB | K5: XML tidak menandai wajib (`pyRequired` true nol kali), tetapi XML / aktivitas TIDAK membuktikan boleh kosong - isi `SaveProductTypeLife_Act` tidak ada; data DEV keenam kunci terisi 32/32. Tanpa nilai, pencocokan K4 tidak punya ID untuk diisi | izinkan kosong = simpan NULL nama dan ID |
| Pencocokan ulang master (422 bila tidak cocok) | K4 keputusan WO | - |
| Pilihan Business `GROUPPANEL = '009'` | K4, diturunkan dari data | ganti nilai `models.GrupBusinessLife` |
| Spasi tepi dipangkas; panjang ≤ 200 byte | pola benefitlife; lebar kolom 947 (tanpa ini ORA-12899 = 500) | - |
| ID dari sequence yang sudah ada = 409 berkalimat | K3 | - |

## Pertanyaan terbuka

- **T1 - urutan bawaan grid**: XML tanpa urutan bawaan; dipakai ID menaik (stabil, urutan pembuatan). Rekomendasi:
  pertahankan.
- **T2 - Save untuk Add dan Edit**: satu Save → `SaveProductTypeLife_Act`; form tanpa ID = Add, sesudah Edit = ubah.
  Rekomendasi: pertahankan (pola benefitlife).
