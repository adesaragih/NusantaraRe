# Struktur Tabel — Accounts: tabel warisan yang diubah dan dibaca

**Keputusan work owner 04-10-2026:** modul `accounts` mengelola (**tambah saja** - akun hanya diinput sekali; nol ubah, nol hapus) tabel warisan
`POOLDATA.T_M_ACCOUNT` — *"buat modul/menu baru, nama modulnya accounts, tabelnya T_M_ACCOUNT"*. Padanan Pega: layar
Account aplikasi SFAGIS (kunci `ASM-SFAGIS-WORK-ACCOUNT ACC-n`), **tidak ada di korpus XML**; acuannya tangkapan
layar Pega dan katalog DEV (`ALL_TAB_COLUMNS`, agregat, 04-10-2026). Modul ini **tidak membuat tabel baru**.

Pembaca lain: popup `ChooseAccount` modul NB FacIn (`Search Group Business`) membaca `T_M_ACCOUNT`;
`T_NB_OPPORTUNITY.ACCOUNT_ID` merujuk `T_M_ACCOUNT.ID` dalam bentuk PENUH (DEV: 2 dari 2 baris) — karena itu ID akun
baru tetap berbentuk penuh.

## T_M_ACCOUNT

Tabel warisan Pega (DEV 04-10-2026: 18.658 baris; sebelum migrasi modul ini tanpa PK, unique, FK, maupun trigger).
Tabel di bab ini = kolom yang DIBUAT migrasi modul ini saja (`840_t_m_account_kolom.sql`,
`TestKolomDDLCocokDenganStruktur`); kolom warisan di bab berikutnya.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `CREATEDATE` | DATE | ya | | "Create Date" | `SYSDATE` saat akun dibuat lewat aplikasi; baris lama kosong |
| `CREATEOP` | VARCHAR2(100) | ya | | "Owner" | akun pelaku (`LOGIN_ID`) saat akun dibuat; baris lama kosong |
| `DESCRIPTION` | VARCHAR2(4000) | ya | | "Description" | isian bebas form; kosong = NULL |

### Kolom warisan T_M_ACCOUNT (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(1020), NOT NULL | `ASM-SFAGIS-WORK-ACCOUNT ACC-n`; **PK `PK_T_M_ACCOUNT`** (migrasi 841). DEV: semua berbeda, nol kosong; satu baris lama ber-ID `ASM-SFAGIS-WORK-ACCOUNT` tanpa nomor (data Pega, dibiarkan; tidak dapat dibuka dari layar) |
| `GROUPBUSINESSID` | VARCHAR2(128) | `BUSINESSGROUP.ID` (DEV: 18.658 dari 18.658 cocok) |
| `GROUPBUSINESS` | VARCHAR2(256) | salinan `BUSINESSGROUP.NOTE` saat disimpan (DEV: 371 baris lama memuat nama lama, dibiarkan) |
| `INSUREDID` | VARCHAR2(1020), NOT NULL | `CLIENT.ID` organisasi `ASM-SFAGIS-WORK-ORG ORG-n` (DEV: 18.655 ber-FLAG `Org`) |
| `INSUREDNAME` | VARCHAR2(256) | salinan `CLIENT.NAME` saat disimpan |

## Sequence dan kunci

- `SEQ_T_M_ACCOUNT` (migrasi 842): nomor `n` akun baru; mulai dari nomor ACC terbesar + 1, dihitung di basis data
  tempat migrasi berjalan (DEV 04-10-2026: terbesar 4916561). `ACCOUNT_SEQ` warisan TIDAK dipakai — nilainya
  tertinggal ±4 juta dari nomor terbesar. Nomor yang sudah terpakai dilewati aplikasi.
- `PK_T_M_ACCOUNT` (migrasi 841) atas `ID`.

## Tabel yang dibaca saja

| Tabel | Kolom dibaca | Untuk |
| --- | --- | --- |
| `CLIENT` | `ID`, `IDVIEW`, `NAME`, `FLAG` (`Org`) | pilihan "Insured Name" dan "Org ID" (dikelola modul Company Detail) |
| `BUSINESSGROUP` | `ID`, `NOTE` | pilihan "Group Business" (DEV: 45 baris, ID unik, NOTE terisi) |
