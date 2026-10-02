# Struktur Tabel — Master Product Name Life: PETA TABEL WARISAN yang ditulis

> Disusun 01-10-2026 (paket 10, brief `PROMPT-IMPLEMENTASI-MODUL-MASTER-PRODUCT-NAME-LIFE.md`).

**P1 (`RALAT-DEV-30-09-2026.md`)** — modul ini **tidak membuat satu tabel, sequence, maupun constraint pun**. Produk
tetap JSON di dua tabel lama seperti Pega (OQ-MPNL-01, bawaan "JSON seperti Pega"); lampiran direkam di tabel lama
(P5). Rentang migrasi `140`–`179` sengaja kosong (`TestMPNLNolMigrasiDiRentang`); satu-satunya migrasi modul ini
adalah slot menu `960` (`UPDATE DIMIGRASI`).

Berkas ini **peta**, bukan DDL: tabel → kolom VERBATIM → tipe katalog → penulis Pega → cara kolom itu ditulis modul
ini. Penjaga yang memakainya: `TestKolomDDLCocokDenganStruktur` / `TestTabelBukanMilikKitaTidakDibuat`
(`inti/backend/penjaga/strukturkolom_test.go`) — ketiga tabel di bawah dinyatakan "Tabel warisan" di `MODUL.md`, jadi
migrasi mana pun yang **membuatnya** merah.

Sumber tipe: `docs/dba-procedures-and-ddl.md` `[data DBA]`. ⛔ Procedure **tidak dipanggil**
(`PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE`): isinya ditiru di Go — upsert dikunci `ID`, ID baru
`'1' || LPAD(M_PRODUCT_LIFE_SEQ.NEXTVAL, 5, '0')`, kedua tabel di SATU transaksi (P4), nol `COMMIT` di teks SQL.

## M_PRODUCT_LIFE

Sisi umum produk — halaman Pega `ProductName`. Penulis Pega `SaveProductName_Act` 9 b1833 (`SaveProductNameLIfe` →
`PEGA_M_PRODUCT_LIFE`) dan 10 b2021 (`SaveProductNameLIfeFlat`). Dibaca view `PRODUCT_LIFE` dan `DOCUMENTCLAIM_LIFE`.

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(6) | `'1' || LPAD(M_PRODUCT_LIFE_SEQ.NEXTVAL, 5, '0')` — dari server, tidak pernah dari klien; > 5 digit atau ID terpakai di salah satu tabel = gagal terang (P6) |
| `JSONDATA` | CLOB, `IS JSON` | halaman `ProductName` berkunci Pega (`POLICYHODER`, `UnderwritingLimitList`, `OutwardList`, …); kunci yang dibaca view dijamin ada; kunci lama yang tidak dikelola layar dipertahankan |
| `RIRISKID` | VARCHAR2(10) | = JSON `RIRISKID` (`SaveProductNameLIfeFlat` b84) |
| `RIRISK` | VARCHAR2(100) | = JSON `RIRISK` |

⭐ **Ralat 01-10-2026 (lanjutan 1 L1):** dua baris lama dikutip — *"`PRODUCTNAME` | VARCHAR2(1000) | = JSON `PRODUCTNAME` (OQ-MPNL-08)"*,
*"`BEGIN_DATE` | DATE | = `BEGIN` inward `dd/MM/yyyy` → `DATE`; NULL bila kosong (OQ-MPNL-08)"* — dicabut: katalog DEV `ALL_TAB_COLUMNS`
`M_PRODUCT_LIFE` hanya empat kolom di atas (`testdata/katalog-dev.json`, uji `TestKolomDitulisAdaDiKatalogDEV`).

## M_PRODUCTINWARD_LIFE

Sisi inward — halaman Pega `ProductNameInward`. Penulis Pega `SaveProductName_Act` 15 b2864
(`SaveProductNameInwardLIfe` → `PEGA_M_PRODUCT_INWARD_LIFE`). Dibaca view `PRODUCTINWARD_LIFE` (Claim Life).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(6) | = ID produk untuk baris baru (R14, OQ-MPNL-02); baris lama dipertahankan ID-nya |
| `JSONDATA` | CLOB, `IS JSON` | halaman `ProductNameInward`; `PRODUCTID` = ID produk; tanggal `dd/MM/yyyy` |

## M_ATTACHMENTPRODUCTNAME

Rekam lampiran produk. Penulis Pega `InsertAttachProdName_Sql` b84, penghapus `DeleteAttachProdName_Sql` b85, pembaca
`GetAttachmentProdName_Sql` b84. Tipe kolom belum ada di katalog modul ini `[data DBA]`.

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | — | `TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3')` (dicoba ulang bila bentrok) |
| `TREATYID` | — | ID produk (`ProductName.ID`) |
| `CATEGORY` | — | `File` (`ProductNameSaveAttachment` 2.1 b475) |
| `FILENAME` | — | nama berkas asli |
| `FILEMIMETYPE` | — | ekstensi huruf kecil (`InsertGoogleStorage_Act` 3 b566) |
| `DATA_JSON` | — | NULL (`""` di Pega) |
| `USERNAME` | — | akun pelaku |
| `T_STORAGE_ID` | — | `ImageID` objek penyimpanan |

## Tabel bersama dan master yang disentuh (digambarkan dokumen pemiliknya)

Tabel di bawah TIDAK digambarkan di sini — pemiliknya yang menggambarkannya, dan dua gambaran yang berbeda akan
saling membatalkan (`TestDokumenSTRUKTURSepakatAtasTabelBersama`).

| Tabel | Sentuhan modul ini |
| --- | --- |
| `T_STORAGE_IMAGE` | ditulis pelaksana stub lampiran (`Insert_T_Storage_SQL` b85: `IMAGEID`, `URLPUBLIC` NULL, `APPFOLDER`, `EXPDATE`, `FILENAME`, `APPNAME`, `STORAGE = 'standard'`), dihapus `DeleteStorage_SQL` b85 |
| `T_LOG_SERVICE_RNM` | outbox bersama (Claim Life 015): `MODUL = 'MASTERPRODUCTNAMELIFE'`, `JENIS_EFEK = 'unggah-lampiran'` |
| `T_FOLDER_IMAGE` | dibaca `APPNAME` (`GetAppName_SQL` b58) |
| `AGENT`, `CLIENT`, `CURRENCY`, `RIRISK_LIFE_SUMMARY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE` | dibaca saja — pemilih master dan `PLAN LIST` (OQ-MPNL-04) |
| `TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE` | dibaca saja — `On Retention` → `BrowseReinstypeOR_SQL` b84 (milik Master Contract Retro Life) |
