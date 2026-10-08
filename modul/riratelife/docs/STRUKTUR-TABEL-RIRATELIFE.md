# Struktur tabel — R/I Rate Life

Modul `riratelife` (perintah work owner 05-10-2026). **RALAT R6 (07-10-2026, `MODUL.md`, menggantikan R4)**:
ringkasan disimpan di **SATU tabel `M_RATE_LIFE_SUMMARY`** berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID
(migrasi inti `927` + `928`); JSONDATA ringkasan dan tabel flat `RATE_LIFE_SUMMARY` (926) dibuang. Rincian rate tetap di
tabel JSON warisan `M_RATE_LIFE` + view `RATE_LIFE` (tidak berubah). Nama objek dan kunci JSON ditulis SEKALI di
`backend/repository/rirl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `M_RATE_LIFE_SUMMARY` | tabel, PK `ID`, 5 kolom (sesudah 928) | Add, Edit, Delete, Simpan Upload, Rate Detail Save (OPERATORID/MODIFIEDDATE) | grid, nama kembar, ID terpakai / ID tertinggi | `mastercontractretrolife`, `masterproductnamelife` (`SELECT ID, USEDBY`, baca-saja) |
| `M_RATE_LIFE` | tabel JSON warisan, 96.038-98.305 baris DEV | Edit (salinan nama), Delete, Simpan Upload, Rate Detail tambah / ubah | ID terpakai / ID tertinggi | — |
| `RATE_LIFE` | view, 8 kolom | — | Rate Detail, jumlah rate, kembar upload | `mastercontractretrolife`, `masterproductnamelife`, Pega `GetRateRetro` (IDUSEDBY) |
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
