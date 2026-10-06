# Struktur tabel — R/I Rate Life

Modul `riratelife` (perintah work owner 05-10-2026). **RALAT R4 (06-10-2026, `MODUL.md`)**: ringkasan kini disimpan di
**tabel flat `RATE_LIFE_SUMMARY`** (migrasi inti `926`) yang menggantikan VIEW warisan bernama sama - keputusan work
owner K-F1/K-F2. Rincian rate tetap di tabel JSON warisan `M_RATE_LIFE` + view `RATE_LIFE` (tidak berubah). Nama
objek dan kunci JSON ditulis SEKALI di `backend/repository/rirl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `RATE_LIFE_SUMMARY` | **tabel flat** (926), PK `ID`, 6 kolom | Add, Edit, Delete, Simpan Upload, Rate Detail Save (OPERATORID/MODIFIEDDATE), alat pindahflat | grid, nama kembar, ID terpakai | `mastercontractretrolife`, `masterproductnamelife` (`SELECT ID, USEDBY` - tidak diubah) |
| `M_RATE_LIFE_SUMMARY` | tabel JSON warisan (347 baris DEV 06-10-2026, masih berubah) | **tidak pernah** (K-F1) | alat pindahflat; ID terpakai / ID tertinggi | prosedur Pega `PEGA_M_RATE_LIFE_SUMMARY` |
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
wajib-isi ditegakkan Go. Lebar keenam kolom cukup untuk data DEV (panjang maksimum isi dicek WO 06-10-2026, di bawah).
Alat pindah tetap melaporkan panjang maksimum (byte) tiap kolom di sumber SAAT berjalan dan MENOLAK `-jalankan`
bila ada nilai yang tidak muat - nol pemotongan.

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(10) NOT NULL, PK | **[terverifikasi data DEV 06-10-2026]** panjang maksimum isi 7 byte (mis. 1000043). `RATE_LIFE_SUMMARY.ID` VARCHAR2(10) (`modul/masterproductnamelife/docs/STRUKTUR-TABEL-MASTER-PRODUCT-NAME-LIFE.md` b180); `.ID` RIRate.xml b1089; `SEQ_M_RATE_LIFE_SUMMARY` |
| `USEDBY` | teks | VARCHAR2(500) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 99 byte; `UPPER(TRIM(USEDBY))` kembar = 0. "R/I RATE NAME" `.USEDBY` RIRate.xml b1273 (tanpa pyMax). Preseden yang menyalin nilai ini: `M_PRODUCTNAME_LIFE_PLAN.RIRATE` VARCHAR2(500) (STRUKTUR MPNL b181). Isian baru dibatasi Go `BatasNama` 200 byte |
| `TYPE` | teks | VARCHAR2(100) | **[terverifikasi data DEV 06-10-2026]** kosong di SEMUA baris. Tidak ada di XML `InboxSummaryRIRate`; RD `BrowseRateLifeSummary` memuat `.TYPE` sebagai kolom (catatan inventaris `Backup/jefri/OUTPUT FIX/03-celah/02-reportdefinition-nb.md` b246) tanpa pengisi yang diketahui. Modul tidak menulisnya: baris baru NULL, Edit tidak menimpa, pindahan apa adanya |
| `MODIFIEDDATE` | teks | VARCHAR2(50) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 23 byte (bentuk Pega `YYYYMMDDTHHMMSS.mmm GMT`, terbaru `20261006T040628.169 GMT`). "MODIFY DATE" `.MODIFIEDDATE` RIRate.xml b10879; dasar aturan delta alat pindah |
| `OPERATORID` | teks | VARCHAR2(200) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 17 byte. "MODIFY OPERATOR" `.OPERATORID` RIRate.xml b10709; Go menulis akun login (dipotong `BatasNama` 200 byte) |
| `FLAG` | teks | VARCHAR2(100) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 2 byte; isi AP 230, PM 106, PY 1, kosong 10. Pengisi dan artinya TIDAK ditemukan (MODUL.md, keputusan FLAG) → **[penyimpangan sadar - menunggu WO]**: ringkasan baru FLAG kosong, Edit tidak menimpa FLAG lama, alat pindah menyalin apa adanya |

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
