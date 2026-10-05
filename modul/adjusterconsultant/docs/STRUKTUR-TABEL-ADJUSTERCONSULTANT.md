# Struktur tabel — Adjuster Consultant

Modul `adjusterconsultant` (keputusan work owner 05-10-2026) menulis tabel warisan `POOLDATA.ADJUSTERCONSULTANT`
langsung; nol tabel baru. Katalog DEV 05-10-2026: 8 baris, PK `ID`, sequence `ADJUSTERCONSULTANT_SEQ`, nol trigger,
nol FK, nol objek bergantung. Pemakai di Pega: modul klaim (Claim Fac In, Claim Prop, Claim Non Prop) memilih
Appointed Adjuster / Consultant lewat `BrowseAdjusterConsultant` dan menyimpan ID serta salinan namanya di data kasus.

## ADJUSTERCONSULTANT

Tabel warisan Pega. Tabel di bab ini = kolom yang DIBUAT migrasi modul ini saja (`870_adjusterconsultant_aktif.sql`,
`TestKolomDDLCocokDenganStruktur`; keputusan work owner 05-10-2026: flag nonaktif pengganti hapus); kolom warisan di
bab berikutnya.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `IS_ACTIVE` | VARCHAR2(1) | tidak | | "Status" | `1` aktif / `0` nonaktif; `DEFAULT '1'` (baris lama aktif), `CK_ADJUSTERCONSULTANT_ACTIVE` |

### Kolom warisan ADJUSTERCONSULTANT (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(100), NOT NULL, PK | situs aktif `M_SITE_DATABASE.ID` + `ADJUSTERCONSULTANT_SEQ` 4 digit (`GetIDConsultanAdj_SQL`); DEV: `1` + 4 digit |
| `NAME` | VARCHAR2(500) | huruf besar (`SaveAdjusterConsultant_Act`); wajib dan unik tanpa beda huruf (keputusan work owner) |
| `ADDRESS` | VARCHAR2(1000) | huruf besar |
| `TELPNO` | VARCHAR2(50) | apa adanya (teks bebas, boleh beberapa nomor) |
| `USERNAME` | VARCHAR2(50) | akun login pengubah terakhir |
| `EDITDATE` | DATE | SYSDATE setiap simpan |

## Tabel yang dibaca saja

| Tabel | Kolom | Untuk |
| --- | --- | --- |
| `M_SITE_DATABASE` | `ID`, `CURRENT_SITE` | awalan ID baru (situs aktif `CURRENT_SITE` 1) |
