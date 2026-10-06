# Struktur tabel — R/I Rate Life

Modul `riratelife` (perintah work owner 05-10-2026) menulis dua tabel fisik warisan Pega dan membaca dua view di
atasnya; nol tabel baru, nol DDL pada objek warisan, nol migrasi modul. Satu-satunya objek baru = dua sequence ID
(migrasi inti 923, K2). Nama objek dan kunci JSON ditulis SEKALI di `backend/repository/rirl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `M_RATE_LIFE_SUMMARY` | tabel (JSON), 339 baris DEV | Add, Edit, Delete, Simpan Upload | ID terpakai / ID tertinggi | — |
| `M_RATE_LIFE` | tabel (JSON), 96.038-98.305 baris DEV | Edit (salinan nama), Delete, Simpan Upload, Rate Detail tambah / ubah | ID terpakai / ID tertinggi | — |
| `RATE_LIFE_SUMMARY` | view, 6 kolom | — | grid, nama kembar | `mastercontractretrolife`, `masterproductnamelife` (ID, USEDBY) |
| `RATE_LIFE` | view, 8 kolom | — | Rate Detail, jumlah rate, kembar upload | `mastercontractretrolife`, `masterproductnamelife`, Pega `GetRateRetro` (IDUSEDBY) |
| `SEQ_M_RATE_LIFE_SUMMARY` | sequence (923) | NEXTVAL | — | — |
| `SEQ_M_RATE_LIFE` | sequence (923) | NEXTVAL | — | — |

Bukti: `RATE_LIFE` = `SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER,
a.JSONDATA.CONTRACT, a.JSONDATA.AGE, a.JSONDATA.RATE FROM M_RATE_LIFE a` dan agregatnya
(`modul/claimlife/docs/KATALOG-TABEL-PESERTA-DAN-TREATY.md` b118-b140); `RATE_LIFE_SUMMARY.ID` VARCHAR2(10)
(`modul/masterproductnamelife/docs/STRUKTUR-TABEL-MASTER-PRODUCT-NAME-LIFE.md` b180); objek ada di DEV
(`modul/masterproductnamelife/docs/OQ-MASTER-PRODUCT-NAME-LIFE.md` b51).

✅ **RALAT R1** (MODUL.md, 06-10-2026): definisi view `RATE_LIFE_SUMMARY` dibaca WO dari `ALL_VIEWS` (owner POOLDATA):
`SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG
FROM M_RATE_LIFE_SUMMARY a`. Asumsi A1 terbukti.

⚠️ View atas CLOB tanpa indeks: setiap kueri `RATE_LIFE ... WHERE IDUSEDBY = :n` mengurai seluruh JSON. Modul ini
membaca kunci kembar upload SEKALI per unggah (daftar IN), dan Delete/Edit memilih baris lewat view
(`ID IN (SELECT ID FROM RATE_LIFE WHERE IDUSEDBY = :n)`) supaya maknanya sama dengan pembaca lain.

## M_RATE_LIFE_SUMMARY

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | ? | | "ID" | `SEQ_M_RATE_LIFE_SUMMARY` (nomor terpakai dilewati) |
| `JSONDATA` | CLOB | ? | | JSON | sisip `JSON_OBJECT`, ubah `JSON_MERGEPATCH` (A2) |

Kunci JSON yang ditulis (terbukti, RALAT R1): `USEDBY` ("R/I RATE NAME", wajib, tidak kembar), `OPERATORID` (akun login),
`MODIFIEDDATE` (format Pega `YYYYMMDDTHHMMSS.mmm GMT`). Kunci lain dibiarkan; `TYPE` dan `FLAG` tidak ditulis
(tidak ada di XML).

## M_RATE_LIFE

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | tidak | unik (data) | rate ID | `SEQ_M_RATE_LIFE` (nomor terpakai dilewati) |
| `JSONDATA` | CLOB | ? | | JSON | sisip `JSON_OBJECT(... ABSENT ON NULL)`; Edit nama `JSON_MERGEPATCH` kunci `USEDBY`; Rate Detail Edit `JSON_MERGEPATCH` kunci `GENDER`, `CONTRACT`, `AGE`, `RATE` (kosong = kunci dibuang) |

Kunci JSON yang ditulis (terbukti dari view): `IDUSEDBY` (ID ringkasan), `USEDBY`, `GENDER` (U/M/F), `CONTRACT`
(0-120 atau tidak ditulis), `AGE` (0-120), `RATE` (teks berkoma desimal). `TYPE` tidak pernah ditulis (NULL di
seluruh baris DEV). Rate Detail (`InboxRIRate`): `GENDER` dan `AGE` boleh kosong (kunci tidak ditulis), `CONTRACT` wajib.

## RATE_LIFE_SUMMARY

View warisan (dinyatakan di `MODUL.md`), dibaca saja.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | ? | | "ID" | `M_RATE_LIFE_SUMMARY.ID` |
| `USEDBY` | VARCHAR2(4000) | ? | | "R/I RATE NAME" | `JSONDATA.USEDBY` (terbukti, R1) |
| `OPERATORID` | VARCHAR2(4000) | ? | | "MODIFY OPERATOR" | `JSONDATA.OPERATORID` (terbukti, R1) |
| `MODIFIEDDATE` | VARCHAR2(4000) | ? | | "MODIFY DATE" | `JSONDATA.MODIFIEDDATE` (terbukti, R1) |
| `TYPE` | VARCHAR2 | ? | | tidak dipakai | `JSONDATA.TYPE` — tidak dirujuk XML (R1) |
| `FLAG` | VARCHAR2 | ? | | tidak dipakai | `JSONDATA.FLAG` — tidak dirujuk XML (R1) |

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
