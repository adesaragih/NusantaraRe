# Struktur tabel — R/I Rate Life

Modul `riratelife` (perintah work owner 05-10-2026). **RALAT R6 (07-10-2026, `MODUL.md`, menggantikan R4)**:
ringkasan disimpan di **SATU tabel `M_RATE_LIFE_SUMMARY`** berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID
(migrasi inti `927` + `928`); JSONDATA ringkasan dan tabel flat `RATE_LIFE_SUMMARY` (926) dibuang. **RALAT R7 (07-10-2026)**: rincian rate
disimpan di **SATU tabel flat `M_RATE_LIFE`** berkolom ID, IDUSEDBY, USEDBY, TYPE, GENDER, CONTRACT, AGE, RATE (teks;
migrasi inti `929` + `930`); JSONDATA rincian dan view `RATE_LIFE` dibuang. Nama objek dan kolom ditulis SEKALI di
`backend/repository/rirl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `M_RATE_LIFE_SUMMARY` | tabel, PK `ID`, 5 kolom (sesudah 928) | Add, Edit, Delete, Simpan Upload, Rate Detail Save (OPERATORID/MODIFIEDDATE) | grid, nama kembar, ID terpakai / ID tertinggi | `mastercontractretrolife`, `masterproductnamelife` (`SELECT ID, USEDBY`, baca-saja) |
| `M_RATE_LIFE` | tabel, PK `ID`, 8 kolom (sesudah 930), 98.306 baris DEV 07-10-2026 | Edit (salinan nama), Delete, Simpan Upload, Rate Detail tambah / ubah | Rate Detail, jumlah rate, kembar upload, ID terpakai / ID tertinggi | `claimlife`, `premiumlistlife`, `mastercontractretrolife`, `masterproductnamelife` (baca-saja); Pega `GetRateRetro` (IDUSEDBY) |
| `RATE_LIFE` | view 8 kolom - **DIBUANG 930** (riwayat di bawah) | — | — | — |
| `SEQ_M_RATE_LIFE_SUMMARY` | sequence (923) | NEXTVAL | — | — |
| `SEQ_M_RATE_LIFE` | sequence (923) | NEXTVAL | — | — |
| `RATE_LIFE_SUMMARY` | tabel flat 926 - **DIBUANG 928** (riwayat di bawah) | — | — | — |

**ID ringkasan baru**: `SEQ_M_RATE_LIFE_SUMMARY` (923); `MaksID`/`AdaID` memeriksa `M_RATE_LIFE_SUMMARY` saja (satu
tabel). Prosedur Pega `PEGA_M_RATE_LIFE_SUMMARY` / `PEGA_M_PLAN_LIFE_SUMMARY` INVALID sesudah 928 (diterima WO): Pega
tidak lagi menulis ringkasan. Prosedur `PEGA_M_RATE_LIFE` INVALID sesudah 930 (diterima WO 07-10-2026).

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

Tabel Pega yang BENTUKNYA diubah migrasi inti 929/930 (RALAT R7, keputusan work owner 07-10-2026; karena itu tidak lagi
dinyatakan "Tabel warisan", preseden `adjusterconsultant` 870). Tabel di bab ini = kolom yang DIBUAT migrasi
(`929_m_rate_life_kolom.sql`, `TestKolomDDLCocokDenganStruktur`), nama dan urutan = view `RATE_LIFE` lama. Diisi 930 dari
JSONDATA APA ADANYA (98.306 baris DEV, termasuk 539 yatim); semua teks, semua NULLABLE; wajib-isi ditegakkan Go. Lebar
= panjang maksimum DEV dengan ruang; nilai lebih panjang menghentikan 930 (ORA-12899), tidak memotong. Indeks
**`IX_M_RATE_LIFE_IDUSEDBY`** (930): Rate Detail, jumlah rate, Delete, kembar upload, dan pembaca lain menyaring IDUSEDBY.

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `IDUSEDBY` | teks | VARCHAR2(10) | **[terverifikasi data DEV 07-10-2026]** panjang maksimum 7 byte. ID ringkasan `M_RATE_LIFE_SUMMARY.ID` (VARCHAR2(10)); 539 baris tanpa ringkasan (item terbuka WO, `MODUL.md` R7) |
| `USEDBY` | teks | VARCHAR2(500) | **[terverifikasi data DEV 07-10-2026]** panjang maksimum 95 byte. Salinan nama ringkasan (lebar = `M_RATE_LIFE_SUMMARY.USEDBY`); diganti Edit nama (A6) |
| `TYPE` | teks | VARCHAR2(100) | **[terverifikasi data DEV 07-10-2026]** kosong di SEMUA baris (panjang maksimum 0). Lebar = `M_RATE_LIFE_SUMMARY.TYPE`. Tidak ditulis modul |
| `GENDER` | teks | VARCHAR2(10) | **[terverifikasi data DEV 07-10-2026]** panjang maksimum 1 byte (U/M/F, A15) |
| `CONTRACT` | teks | VARCHAR2(10) | **[terverifikasi data DEV 07-10-2026]** panjang maksimum 3 byte; Go: bulat 0-120 |
| `AGE` | teks | VARCHAR2(10) | **[terverifikasi data DEV 07-10-2026]** panjang maksimum 3 byte; Go: bulat 0-120 |
| `RATE` | teks | VARCHAR2(50) | **[terverifikasi data DEV 07-10-2026]** panjang maksimum 20 byte; 90.437 baris berkoma desimal, 3.226 bertitik desimal - dipindah APA ADANYA, tetap teks (keputusan WO 07-10-2026). Isian baru disimpan berkoma desimal (A8) |

### Kolom warisan M_RATE_LIFE (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(10) NOT NULL, PK | **[terverifikasi data DEV 07-10-2026]** panjang maksimum isi 7 byte. `SEQ_M_RATE_LIFE` (923; terakhir 1117016 di DEV 07-10-2026), nomor terpakai dilewati |
| `JSONDATA` | CLOB, constraint `ENSURE_M_RATE_LIFE_JSON` (IS JSON) | **DIBUANG 930** (bersama constraint-nya). Ketujuh kunci view sudah di kolom; `FLAG` (tidak digunakan, RALAT R5) hanya di cadangan CSV `ID + JSONDATA` (`LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md` (a)) |

## RATE_LIFE

RIWAYAT - view warisan Pega `SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER,
a.JSONDATA.CONTRACT, a.JSONDATA.AGE, a.JSONDATA.RATE FROM M_RATE_LIFE a`, DIBUANG `930_m_rate_life_satu_tabel.sql`
(RALAT R7); dibangun ulang hanya oleh jalur mundur 930. Kolomnya kini kolom `M_RATE_LIFE` di atas.

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
