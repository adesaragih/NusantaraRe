# Struktur tabel — Treaty Exchange Yearly

Modul `treatyexchangeyearly` (keputusan work owner 05-10-2026) menulis tabel warisan `POOLDATA.TREATYEXCHANGEYEARLY`
langsung; nol tabel baru, nol DDL, nol migrasi sendiri. Katalog DEV 05-10-2026: 140 baris, 25 mata uang, tahun
2017-2026; seluruh kolom VARCHAR2; tanpa PK, constraint, indeks, trigger. Sequence `TREATYEXCHANGE_SEQ` (nilai
terakhir 117, cache 0) dipakai prosedur `PEGA_TREATYEXCHANGE`, yang menulis `M_TREATYEXCHANGE` (JSON, 117 baris) dan
tabel ini sekaligus - 23 baris tabel ini tanpa pasangan JSON.

Temuan data DEV:

- ID tidak unik: 10114, 10115, 10116 masing-masing dua baris; dua ID di luar pola (13358, 13359). Edit memakai ROWID.
- Kombinasi kembar (tahun, quarter, mata uang): 2019 USD, 2020 SGD, 2025 USD - dibiarkan; baris baru ditolak kembar.
- 11 baris STARTDATE berbentuk rusak `20190801T00000.000 GMT`; baris 2024-2025 ber-jam (`T105000`). Dibiarkan bila
  tanggalnya tidak diubah.
- TOUSD tidak konsisten (`14500`, `1`, `0`, kosong); QURRENCYID selalu kosong.

Pembaca (membandingkan tanggal SEBAGAI TEKS - format Pega wajib dijaga): FacIn / NB Treaty In
`GetCurrencyToIDR_SQL`, `GetKursLimitSpreading_SQL`; Treaty In `BrowseTreatyExchangeYearly_RD` (baris pertama per
`TreatyYear`, `Currency`, `Quarter = "0"`); `UploadCSVAggregate_Act`; modul `treatycontractout` dan `aggregate`
(dibaca saja).

## TREATYEXCHANGEYEARLY

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(50) | ya | | "ID" | situs aktif + `TREATYEXCHANGE_SEQ` 4 digit (`10118` berikutnya di DEV) |
| `QURRENCYID` | VARCHAR2(255) | ya | | tidak | selalu kosong; tidak ditulis |
| `TREATYYEAR` | VARCHAR2(255) | ya | | "Treaty Year" | 4 angka |
| `STARTDATE` | VARCHAR2(255) | ya | | "Start Date" | `YYYYMMDDT000000.000 GMT` |
| `ENDDATE` | VARCHAR2(255) | ya | | "End Date" | `YYYYMMDDT000000.000 GMT`; tidak sebelum Start Date |
| `TOIDR` | VARCHAR2(255) | ya | | "To IDR" | angka bertitik desimal, wajib |
| `TOUSD` | VARCHAR2(255) | ya | | "To USD" | angka bertitik desimal, boleh kosong |
| `USERID` | VARCHAR2(255) | ya | | "Last Edited By" | akun login |
| `DATEIU` | VARCHAR2(255) | ya | | "Last Edited" (WIB) | waktu simpan, format Pega GMT |
| `IDCURRENCY` | VARCHAR2(255) | ya | | "Currency" | `CURRENCY.ID` |
| `CURRENCY` | VARCHAR2(255) | ya | | "Currency" | salinan kode `CURRENCY.CURRENCY` |
| `QUARTER` | VARCHAR2(255) | ya | | "Quarter" | `0` tahunan, `1`-`4` |
| `DATEIN` | VARCHAR2(255) | ya | | View | waktu Add, format Pega GMT |

## Tabel yang dibaca saja

| Tabel | Kolom | Untuk |
| --- | --- | --- |
| `CURRENCY` (view) | `ID`, `CURRENCY`, `NOTE` | pilihan Currency, kode yang disalin, nama tampil |
| `M_SITE_DATABASE` | `ID`, `CURRENT_SITE` | awalan ID baru (situs aktif `CURRENT_SITE` 1) |
