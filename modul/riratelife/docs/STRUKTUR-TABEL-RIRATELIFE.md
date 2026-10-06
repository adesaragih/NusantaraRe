# Struktur tabel — R/I Rate Life

Modul `riratelife` (perintah work owner 05-10-2026). **RALAT R4 (06-10-2026, `MODUL.md`)**: ringkasan kini disimpan di
**tabel flat `RATE_LIFE_SUMMARY`** (migrasi inti `926`) yang menggantikan VIEW warisan bernama sama - keputusan work
owner K-F1/K-F2. Rincian rate tetap di tabel JSON warisan `M_RATE_LIFE` + view `RATE_LIFE` (tidak berubah). Nama
objek dan kunci JSON ditulis SEKALI di `backend/repository/rirl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `RATE_LIFE_SUMMARY` | **tabel flat** (926), PK `ID`, 6 kolom | Add, Edit, Delete, Simpan Upload, Rate Detail Save (OPERATORID/MODIFIEDDATE), alat pindahflat | grid, nama kembar, ID terpakai | `mastercontractretrolife`, `masterproductnamelife` (`SELECT ID, USEDBY` - tidak diubah) |
| `M_RATE_LIFE_SUMMARY` | tabel JSON warisan, 339 baris DEV | **tidak pernah** (K-F1) | alat pindahflat; ID terpakai / ID tertinggi | prosedur Pega `PEGA_M_RATE_LIFE_SUMMARY` |
| `M_RATE_LIFE` | tabel JSON warisan, 96.038-98.305 baris DEV | Edit (salinan nama), Delete, Simpan Upload, Rate Detail tambah / ubah | ID terpakai / ID tertinggi | — |
| `RATE_LIFE` | view, 8 kolom | — | Rate Detail, jumlah rate, kembar upload | `mastercontractretrolife`, `masterproductnamelife`, Pega `GetRateRetro` (IDUSEDBY) |
| `SEQ_M_RATE_LIFE_SUMMARY` | sequence (923) | NEXTVAL | — | — |
| `SEQ_M_RATE_LIFE` | sequence (923) | NEXTVAL | — | — |

**ID ringkasan baru**: `SEQ_M_RATE_LIFE_SUMMARY` (923, sudah jalan di DEV). `MaksID`/`AdaID` memeriksa tabel flat
**DAN** `M_RATE_LIFE_SUMMARY`: ringkasan yang masih di JSON (belum dipindah, atau ditulis Pega sesudah cutover) tetap
dihitung terpakai, sehingga ID baru aplikasi tidak pernah sama dengan ID yang kelak dibawa alat pindah (yang akan
menolak putaran bila ID sama berisi beda). Membaca JSON warisan untuk ini aman: SELECT saja.

## RATE_LIFE_SUMMARY

Tabel flat (migrasi inti `926_rate_life_summary_flat.sql`). Kolom = PERSIS keenam kolom VIEW warisan
(`SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG
FROM M_RATE_LIFE_SUMMARY a`, ALL_VIEWS dibaca WO 06-10-2026 - semua VARCHAR2; K-F2). Semua NULLABLE kecuali PK;
wajib-isi ditegakkan Go. Alat pindah melaporkan panjang maksimum (byte) tiap kolom di sumber dan MENOLAK `-jalankan`
bila ada nilai yang tidak muat - nol pemotongan.

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(10) NOT NULL, PK | `RATE_LIFE_SUMMARY.ID` VARCHAR2(10) (`modul/masterproductnamelife/docs/STRUKTUR-TABEL-MASTER-PRODUCT-NAME-LIFE.md` b180); `.ID` RIRate.xml b1089; `SEQ_M_RATE_LIFE_SUMMARY` |
| `USEDBY` | teks | VARCHAR2(500) | "R/I RATE NAME" `.USEDBY` RIRate.xml b1273 (tanpa pyMax; lebar layar 200 px). Preseden yang menyalin nilai ini: `M_PRODUCTNAME_LIFE_PLAN.RIRATE` VARCHAR2(500) (STRUKTUR MPNL b181, `.RIRATE ← .USEDBY`). Isian baru dibatasi Go `BatasNama` 200 byte |
| `TYPE` | teks | VARCHAR2(100) | **[penyimpangan sadar]** tidak ada di XML, data DEV tidak terdokumentasi di repo → preseden K6 teks; tidak ditulis modul (baris baru NULL), baris pindahan apa adanya |
| `MODIFIEDDATE` | teks | VARCHAR2(50) | "MODIFY DATE" `.MODIFIEDDATE` RIRate.xml b10879; bentuk Pega `YYYYMMDDTHHMMSS.mmm GMT` = 23 byte. **[penyimpangan sadar]** lebar data lama tidak terdokumentasi → K6 teks 50 |
| `OPERATORID` | teks | VARCHAR2(200) | "MODIFY OPERATOR" `.OPERATORID` RIRate.xml b10709; Go menulis akun login dipotong `BatasNama` 200 byte. **[penyimpangan sadar]** lebar data lama tidak terdokumentasi |
| `FLAG` | teks | VARCHAR2(100) | **[penyimpangan sadar]** seperti `TYPE` |

**Indeks `IX_RATE_LIFE_SUMMARY_NAMA` = `UPPER(TRIM(USEDBY))`** (fungsi): dipakai SETIAP Save dan Simpan Upload
(`SqlPemakaiNama` - nama tidak kembar tanpa beda huruf, pencocokan upload). Sekaligus PENGAMAN migrasi: bila view lama
belum dibuang, CREATE TABLE 926 dilewati pra-terbang (ORA-00955), lalu `CREATE INDEX` atas view gagal ORA-01702 dan
`-migrate` berhenti keras tanpa mencatat 926 (`docs/LANGKAH-WO-RIRATELIFE-FLAT.md`).

## M_RATE_LIFE_SUMMARY

Tabel warisan (dinyatakan di `MODUL.md`), dibaca saja sejak RALAT R4.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | ? | | ID terpakai, sumber pindah | prosedur Pega / `SEQ_M_RATE_LIFE_SUMMARY` (sebelum R4) |
| `JSONDATA` | CLOB | ? | | sumber pindah | kunci `USEDBY`, `TYPE`, `MODIFIEDDATE`, `OPERATORID`, `FLAG` (RALAT R1) |

## M_RATE_LIFE

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | tidak | unik (data) | rate ID | `SEQ_M_RATE_LIFE` (nomor terpakai dilewati) |
| `JSONDATA` | CLOB | ? | | JSON | sisip `JSON_OBJECT(... ABSENT ON NULL)`; Edit nama kunci `USEDBY`; Rate Detail Edit kunci `GENDER`, `CONTRACT`, `AGE`, `RATE` (kosong = kunci dibuang) - baca `FOR UPDATE`, ganti di Go, tulis utuh (RALAT R3) |

Kunci JSON yang ditulis (terbukti dari view): `IDUSEDBY` (ID ringkasan), `USEDBY`, `GENDER` (U/M/F), `CONTRACT`
(0-120 atau tidak ditulis), `AGE` (0-120), `RATE` (teks berkoma desimal). `TYPE` tidak pernah ditulis (NULL di
seluruh baris DEV). Rate Detail (`InboxRIRate`): `GENDER` dan `AGE` boleh kosong (kunci tidak ditulis), `CONTRACT` wajib.

## RATE_LIFE

View warisan (dinyatakan di `MODUL.md`), dibaca saja.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | tidak | | rate ID | `M_RATE_LIFE.ID` |
| `IDUSEDBY` | VARCHAR2(4000) | ya | | kunci ringkasan | `JSONDATA.IDUSEDBY` |
| `USEDBY` | VARCHAR2(4000) | ya | | "USEDBY" | `JSONDATA.USEDBY` |
| `TYPE` | VARCHAR2(4000) | ya | | — | `JSONDATA.TYPE` (NULL seluruhnya) |
| `GENDER` | VARCHAR2(4000) | ya | | "GENDER" | `JSONDATA.GENDER` |
| `CONTRACT` | VARCHAR2(4000) | ya | | "CONTRACT" | `JSONDATA.CONTRACT` |
| `AGE` | VARCHAR2(4000) | ya | | "AGE" | `JSONDATA.AGE` |
| `RATE` | VARCHAR2(4000) | ya | | "RATE" | `JSONDATA.RATE` |
