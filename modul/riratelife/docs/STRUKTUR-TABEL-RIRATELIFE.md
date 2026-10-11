# Struktur tabel — R/I Rate Life

Modul `riratelife` (perintah work owner 05-10-2026). **RALAT R6 (07-10-2026, `MODUL.md`, menggantikan R4)**:
ringkasan disimpan di **SATU tabel `M_RATE_LIFE_SUMMARY`** berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID
(migrasi inti `927` + `928`); JSONDATA ringkasan dan tabel flat `RATE_LIFE_SUMMARY` (926) dibuang. Rincian rate tetap di
tabel JSON warisan `M_RATE_LIFE` + view `RATE_LIFE` (tidak berubah). Nama objek dan kolom ditulis SEKALI di
`backend/repository/rirl_tabel.go`.

**RALAT R7 (10-10-2026, `MODUL.md`)** - keputusan work owner *"view rate_life sudah tidak digunakan, sekarang gunakan tabel
M_RATE_LIFE"*: view `RATE_LIFE` **dibuang di DEV** (ORA-00942) dan `M_RATE_LIFE` kini **tabel FLAT 8 kolom tanpa JSONDATA**
(ALL_TAB_COLUMNS DEV dikirim work owner 10-10-2026; diubah DBA, tanpa migrasi). Rincian dibaca DAN ditulis di kolom
`M_RATE_LIFE` (INSERT / UPDATE kolom biasa); nol JSON.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `M_RATE_LIFE_SUMMARY` | tabel, PK `ID`, 5 kolom (sesudah 928) | Add, Edit, Delete, Simpan Upload, Rate Detail Save (OPERATORID/MODIFIEDDATE) | grid, nama kembar, ID terpakai / ID tertinggi | `mastercontractretrolife`, `masterproductnamelife` (`SELECT ID, USEDBY`, baca-saja) |
| `M_RATE_LIFE` | tabel warisan FLAT 8 kolom (RALAT R7), 96.038-98.305 baris DEV | Edit (UPDATE USEDBY), Delete, Simpan Upload, Rate Detail tambah / ubah (INSERT / UPDATE kolom) | Rate Detail, jumlah rate, kembar upload, ID terpakai / ID tertinggi | `mastercontractretrolife`, `masterproductnamelife`, `premiumlistlife`, `claimlife` (baca-saja, IDUSEDBY) |
| `RATE_LIFE` | view - **DIBUANG di DEV 10-10-2026, tidak dipakai lagi** (RALAT R7) | — | — | — |
| `SEQ_M_RATE_LIFE_SUMMARY` | sequence (923) | NEXTVAL | — | — |
| `SEQ_M_RATE_LIFE` | sequence (923) | NEXTVAL | — | — |
| `RATE_LIFE_SUMMARY` | tabel flat 926 - **DIBUANG 928** (riwayat di bawah) | — | — | — |

**ID ringkasan baru**: `SEQ_M_RATE_LIFE_SUMMARY` (923); `MaksID`/`AdaID` memeriksa `M_RATE_LIFE_SUMMARY` saja (satu
tabel). Prosedur Pega `PEGA_M_RATE_LIFE_SUMMARY` / `PEGA_M_PLAN_LIFE_SUMMARY` INVALID sesudah 928 (diterima WO): Pega
tidak lagi menulis ringkasan.

## M_RATE_LIFE_SUMMARY

Tabel Pega yang BENTUKNYA diubah migrasi inti 927/928 (karena itu tidak lagi dinyatakan "Tabel warisan", preseden
`adjusterconsultant` 870). Tabel di bab ini = kolom yang DIBUAT migrasi (`927_m_rate_life_summary_kolom.sql`,
`TestKolomDDLCocokDenganStruktur`); lebar = 926. Diisi 928 dari tabel flat (sumber kebenaran sejak DEV 07-10-2026 09:09);
baris yang sudah dihapus aplikasi ikut dibuang. Semua NULLABLE; wajib-isi ditegakkan Go. Indeks
**`IX_M_RATE_LIFE_SUMMARY_NAMA` = `UPPER(TRIM(USEDBY))`** (928): dipakai setiap Save dan Simpan Upload (`SqlPemakaiNama`).

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `USEDBY` | teks | VARCHAR2(500) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 99 byte; `UPPER(TRIM(USEDBY))` kembar = 0. "R/I RATE NAME" `.USEDBY` RIRate.xml b1273 (tanpa pyMax). Preseden yang menyalin nilai ini: `M_PRODUCTNAME_LIFE_PLAN.RIRATE` VARCHAR2(500) (STRUKTUR MPNL b181). Isian baru dibatasi Go `BatasNama` 200 byte |
| `TYPE` | teks | VARCHAR2(100) | **[terverifikasi data DEV 06-10-2026]** kosong di SEMUA baris. Tidak ada di XML `InboxSummaryRIRate`; RD `BrowseRateLifeSummary` memuat `.TYPE` sebagai kolom (catatan inventaris `Backup/jefri/OUTPUT FIX/03-celah/02-reportdefinition-nb.md` b246) tanpa pengisi yang diketahui. Modul tidak menulisnya: baris baru NULL, Edit tidak menimpa, pindahan apa adanya |
| `MODIFIEDDATE` | teks | VARCHAR2(50) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 23 byte (bentuk Pega `YYYYMMDDTHHMMSS.mmm GMT`, terbaru `20261006T040628.169 GMT`). "MODIFY DATE" `.MODIFIEDDATE` RIRate.xml b10879 |
| `OPERATORID` | teks | VARCHAR2(200) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 17 byte. "MODIFY OPERATOR" `.OPERATORID` RIRate.xml b10709; Go menulis akun login (dipotong `BatasNama` 200 byte) |

### Kolom warisan M_RATE_LIFE_SUMMARY (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(10) NOT NULL, PK | **[terverifikasi data DEV 06-10-2026]** panjang maksimum isi 7 byte (mis. 1000043). `RATE_LIFE_SUMMARY.ID` VARCHAR2(10) (`modul/masterproductnamelife/docs/STRUKTUR-TABEL-MASTER-PRODUCT-NAME-LIFE.md` b180); `.ID` RIRate.xml b1089; `SEQ_M_RATE_LIFE_SUMMARY` |
| `JSONDATA` | CLOB, constraint `ENSURE_M_RATE_LIFE_SUMMARY_JSON` (IS JSON) | **DIBUANG 928** (bersama constraint-nya). Kunci `USEDBY`, `TYPE`, `MODIFIEDDATE`, `OPERATORID` sudah di kolom; `FLAG` (tidak digunakan, RALAT R5) hanya di cadangan CSV `ID + JSONDATA` (LANGKAH-WO (a)) |

## RATE_LIFE_SUMMARY

RIWAYAT - tabel flat migrasi inti `926_rate_life_summary_flat.sql` (RALAT R4, sudah jalan di DEV), DIBUANG
`928_m_rate_life_summary_satu_tabel.sql` (RALAT R6); dibangun ulang hanya oleh jalur mundur 928. Kolom = kolom VIEW
warisan KECUALI FLAG (`SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID,
a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`, ALL_VIEWS dibaca WO 06-10-2026). Indeks `IX_RATE_LIFE_SUMMARY_NAMA`
(`UPPER(TRIM(USEDBY))`) ikut terbuang bersama tabelnya.

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(10) NOT NULL, PK | **[terverifikasi data DEV 06-10-2026]** panjang maksimum isi 7 byte (mis. 1000043). `RATE_LIFE_SUMMARY.ID` VARCHAR2(10) (`modul/masterproductnamelife/docs/STRUKTUR-TABEL-MASTER-PRODUCT-NAME-LIFE.md` b180); `.ID` RIRate.xml b1089; `SEQ_M_RATE_LIFE_SUMMARY` |
| `USEDBY` | teks | VARCHAR2(500) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 99 byte; `UPPER(TRIM(USEDBY))` kembar = 0. "R/I RATE NAME" `.USEDBY` RIRate.xml b1273 (tanpa pyMax). Preseden yang menyalin nilai ini: `M_PRODUCTNAME_LIFE_PLAN.RIRATE` VARCHAR2(500) (STRUKTUR MPNL b181). Isian baru dibatasi Go `BatasNama` 200 byte |
| `TYPE` | teks | VARCHAR2(100) | **[terverifikasi data DEV 06-10-2026]** kosong di SEMUA baris. Tidak ada di XML `InboxSummaryRIRate`; RD `BrowseRateLifeSummary` memuat `.TYPE` sebagai kolom (catatan inventaris `Backup/jefri/OUTPUT FIX/03-celah/02-reportdefinition-nb.md` b246) tanpa pengisi yang diketahui. Modul tidak menulisnya: baris baru NULL, Edit tidak menimpa, pindahan apa adanya |
| `MODIFIEDDATE` | teks | VARCHAR2(50) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 23 byte (bentuk Pega `YYYYMMDDTHHMMSS.mmm GMT`, terbaru `20261006T040628.169 GMT`). "MODIFY DATE" `.MODIFIEDDATE` RIRate.xml b10879 |
| `OPERATORID` | teks | VARCHAR2(200) | **[terverifikasi data DEV 06-10-2026]** panjang maksimum 17 byte. "MODIFY OPERATOR" `.OPERATORID` RIRate.xml b10709; Go menulis akun login (dipotong `BatasNama` 200 byte) |

## M_RATE_LIFE

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi). **Sejak RALAT R7 (10-10-2026) FLAT, tanpa JSONDATA** -
kolom per ALL_TAB_COLUMNS DEV dikirim work owner 10-10-2026 (sama dengan kolom view `RATE_LIFE` lama).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | tidak | unik (data) | rate ID | `SEQ_M_RATE_LIFE` (nomor terpakai dilewati) |
| `IDUSEDBY` | VARCHAR2(10) | ya | | kunci ringkasan | ID ringkasan `M_RATE_LIFE_SUMMARY` |
| `USEDBY` | VARCHAR2(500) | ya | | "R/I RATE NAME" | salinan nama ringkasan; Edit nama = `UPDATE M_RATE_LIFE SET USEDBY = :nama WHERE IDUSEDBY = :id` (Go `BatasNama` 200) |
| `TYPE` | VARCHAR2(100) | ya | | — | tidak pernah ditulis (NULL di seluruh baris DEV) |
| `GENDER` | VARCHAR2(10) | ya | | "GENDER" | U / M / F; kosong = NULL |
| `CONTRACT` | VARCHAR2(10) | ya | | "CONTRACT" | teks bulat 0-120; kosong = NULL |
| `AGE` | VARCHAR2(10) | ya | | "AGE" | teks bulat 0-120; kosong = NULL |
| `RATE` | VARCHAR2(50) | ya | | "RATE" | teks berkoma desimal (`0.50` → `0,50`, A8); Go menolak > 50 byte (`models.LebarRate`) |

Sisip (Simpan Upload, Rate Detail Save) = `INSERT (ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE)`; Rate Detail Edit =
`UPDATE … SET GENDER, CONTRACT, AGE, RATE WHERE ID = :id AND TRIM(IDUSEDBY) = :ringkasan` (dulu "kunci JSON dibuang", kini
NULL); Delete = `DELETE … WHERE IDUSEDBY = :id`. Rate Detail (`InboxRIRate`): `GENDER` dan `AGE` boleh kosong, `CONTRACT` wajib.

## Riwayat — view RATE_LIFE (dibuang di DEV 10-10-2026, RALAT R7)

Tidak dipakai lagi oleh modul mana pun (penjaga `TestTidakAdaViewRateLifeDiKodeProduksi`). Definisi lama, untuk sejarah:

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
