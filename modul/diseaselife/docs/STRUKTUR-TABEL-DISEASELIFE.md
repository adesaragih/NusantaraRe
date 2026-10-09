# Struktur tabel — Disease Life

Modul `diseaselife` (keputusan work owner 08-10-2026 K0, D1-D4, K5, `MODUL.md`). SATU tabel: `DISEASE_LIFE`, tabel
Pega yang **SUDAH flat** - TIDAK di-RENAME, kolomnya TIDAK diubah (dibaca claimlife, ditulis prosedur Pega
`PEGA_DISEASE_LIFE` yang masih VALID). Migrasi MODUL 080-081 hanya menambah sequence aplikasi dan PK; nol tabel baru,
nol kolom baru. Nama objek dan kolom ditulis SEKALI di `backend/repository/dsl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca / penulis lain |
| --- | --- | --- | --- | --- |
| `DISEASE_LIFE` | tabel flat, 3 kolom, **97.586 baris DEV**, sesudah 081: PK `PK_DISEASE_LIFE (ID)` | Save Add / Edit | grid (berhalaman + bersaring di server), ID terpakai, ICD Code kembar | `claimlife` (pencarian diagnosa, kueri tidak berubah); prosedur `PEGA_DISEASE_LIFE` (INSERT / UPDATE langsung - masih VALID) |
| `SEQ_DISEASE_LIFE` | sequence aplikasi BARU (080) | NEXTVAL | — | — |
| `M_DISEASE_LIFE` (3 baris JSON lama), `PEGA_M_DISEASE_LIFE`, `M_DISEASE_LIFE_SEQ` (last_number 2052) | objek Pega lama | — | — | TIDAK disentuh (D2, uji `TestMigrasiTanpaTabelLainDanTanpaTulisData`, `TestPeriksaTulis`); pertanyaan terbuka |
| `T_CLAIMLF_DIAGNOSE`, `T_CLAIMLF_PREMIUMLIST_DETAIL` | tabel claimlife | — | — | menyimpan TEKS `ICD_CODE` / `DISEASE`, bukan ID (fakta WO) |

**ID baru** (D1.1): `TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL)`; sequence mulai ID angka tertinggi + 1 dihitung di basis data
tempat migrasi berjalan (DEV: 197586). Aplikasi tetap memeriksa ID belum terpakai sebelum INSERT (`AdaID`); NEXTVAL yang
menabrak ID = 409 berkalimat. `models.BentukID`, uji `TestBentukIDDariSequence`.

**Mengapa BUKAN rumus Pega** `'1' || LPAD(M_DISEASE_LIFE_SEQ.NEXTVAL, 5, '0')` (prosedur `PEGA_DISEASE_LIFE`): data
`DISEASE_LIFE` adalah katalog ICD hasil impor dengan ID `1nnnnn` sampai 197585, sedangkan `M_DISEASE_LIFE_SEQ` di 2052.
Sequence itu menghasilkan 102051 yang SUDAH dipakai data impor - itulah asal ID kembar 102051 (`C718 …` impor dan
`TEST123 / Sakit` uji) karena tabel tanpa PK menerimanya - dan 102052, 102053, 102054, … juga sudah terpakai. Rumus
Pega dipakai apa adanya = setiap Add bertabrakan.

### Kolom warisan DISEASE_LIFE (tidak dibuat dan tidak diubah migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(100), NULLABLE → NOT NULL sesudah 081 (PK) | `1nnnnn`, tertinggi 197585 DEV; satu ID kembar 102051 (dihapus WO D2 sebelum 081); grid "ID" b3758 = `.Number` b4321, form "Number / ID" b1021 disabled b1070 |
| `ICD_CODE` | VARCHAR2(100) | 97.586 unik, maks 7 byte, huruf besar (kecuali baris uji); form "ICD Code" b1200 (`.ICD_Code` b1229) |
| `DISEASE` | VARCHAR2(1000) | maks 290 byte, huruf besar (kecuali baris uji `Sakit`); form "Disease" b1470 (`.Disease` b1499, pxTextArea b1502) |

Bab di atas sengaja berjudul tingkat tiga (bukan `## DISEASE_LIFE`): penjaga `TestKolomDDLCocokDenganStruktur`
membandingkan tabel berjudul `##` dengan kolom yang DIBUAT DDL - DDL modul ini tidak membuat kolom apa pun.

### Objek yang dibuat migrasi modul

| Objek | Migrasi | Bentuk |
| --- | --- | --- |
| `SEQ_DISEASE_LIFE` | 080 | START WITH NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0) + 1, INCREMENT BY 1 NOCACHE NOCYCLE |
| `PK_DISEASE_LIFE` | 081 | PRIMARY KEY (ID) + indeks unik bernama sama; NOT NULL implisit pada ID |
