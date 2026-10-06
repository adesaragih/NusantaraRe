# Struktur tabel — Reinsurance Type

Modul `reinsurancetype` (keputusan work owner 05-10-2026) menulis tabel warisan `POOLDATA.REINSURANCETYPE` langsung;
nol tabel baru, nol DDL, nol migrasi sendiri. Katalog DEV 05-10-2026: 130 baris, seluruh kolom VARCHAR2(100), tanpa
PK; indeks unik `REINSURANCETYPE_INDEX1` (`ID`, `NOTE`). Sequence `M_REINSURANCETYPE_SEQ` (nilai terakhir 263) dipakai
prosedur `PEGA_REINSURANCETYPE`; `REINSURANCETYPE_SEQ` (642) tidak dipakai. `M_REINSURANCETYPE` (JSON Pega, 102 baris,
102 cocok ID) tidak disentuh. Data masih ditulis dari luar aplikasi ini (TGLUPDATE terakhir 05-10-2026).

Temuan data DEV:

- Type: 1 (8 baris, Own Retention), 2 (76, Treaty Out), 3 (14, Facultative), 4 (31, Treaty In), kosong 1.
- Flag: `active` 84, `inactive` 12, `1` 5 (disaring `BrowseReinsuranceTypeLimit_RD` Contract Retro Life), kosong 29.
- Nama kembar: OR, SURPLUS, 2ND SURPLUS, XL (dua baris masing-masing) - dibiarkan; nama baru ditolak kembar.
- CODE `00` kecuali Type 4 (kode 2 digit 10-51); NOURUT hanya 23 baris (1-7); GROUPTYPE 14 baris (OR, QS, RI, SPL),
  tidak ditulis prosedur Pega.

Pembaca: Pega `BrowseReinsuranceType_RD` (saringan ID, Note, Flag, Type), `GetReinsuranceTypeBYName_SQL`
(`type ='4' and FLAG ='active' and note = …`), `GetReinstypeIDbyName_SQL` (`where note = …`), view `LIMITTREATYIN`;
modul `treatycontractout` (Flag `active`, Type 1-3) dan `mastercontractretrolife` (Flag `1`) - dibaca saja.

## REINSURANCETYPE

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(100) | ya | UNIQUE (ID, NOTE) | "ID" | `1` + `M_REINSURANCETYPE_SEQ` 4 digit (`10263` berikutnya di DEV) |
| `NOTE` | VARCHAR2(100) | ya | UNIQUE (ID, NOTE) | "Name" | wajib, huruf besar, tidak kembar |
| `TYPE` | VARCHAR2(100) | ya | | "Type" | `1` Own Retention, `2` Treaty Out, `3` Facultative, `4` Treaty In |
| `SOANOTE` | VARCHAR2(100) | ya | | "SOA Name" | huruf besar, boleh kosong |
| `CODE` | VARCHAR2(100) | ya | | "Code" | angka, kosong = `00` |
| `FLAG` | VARCHAR2(100) | ya | | "Flag" | `active` / `inactive`; warisan `1` / kosong dibiarkan |
| `USERID` | VARCHAR2(100) | ya | | "Last Edited By" | akun login |
| `TGLUPDATE` | VARCHAR2(100) | ya | | "Last Edited" (WIB) | waktu simpan, format Pega GMT |
| `NOURUT` | VARCHAR2(100) | ya | | "No Urut" | angka, boleh kosong |
| `GROUPTYPE` | VARCHAR2(100) | ya | | "Group Type" | OR / QS / RI / SPL, boleh kosong |
