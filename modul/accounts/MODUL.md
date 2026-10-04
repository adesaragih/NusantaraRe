# Modul `accounts` — Accounts

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
04-10-2026: *"buat modul/menu baru, jangan sentuh modul lain, nama modulnya accounts, tabelnya T_M_ACCOUNT"*. Padanan
Pega: layar Account aplikasi SFAGIS (kunci `ASM-SFAGIS-WORK-ACCOUNT ACC-n`) — **tidak ada di korpus XML**; acuannya
tangkapan layar Pega dari work owner dan katalog DEV. Baris menunya dibuat migrasi inti `908` (slot menu modul tidak
boleh membuat kelompok), dan `Folder korpus` di bawah = label menu `M_NAV_MENU.LABEL` (keputusan work owner:
"Accounts", kelompok MASTER).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — diambil dari
cadangan `760-899` dan `990-999`.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `accounts` |
| Folder korpus | `Accounts` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-ACCOUNTS` |
| Status | dimigrasi |
| Rentang migrasi | `840-879` |
| Slot menu | `994-995` |
| Prefix rute API | `/api/accounts` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-ACCOUNTS.md` — tabel yang diubah dan dibaca, sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `accounts.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas) — keputusan work owner 04-10-2026

- Tulis langsung ke `POOLDATA.T_M_ACCOUNT`: **tambah saja** - akun hanya diinput sekali, nol ubah, nol hapus, nol
  layar View. `CLIENT` (organisasi) dan
  `BUSINESSGROUP` dibaca saja.
- **Insured Name** \* = organisasi `CLIENT` (`FLAG` `Org`): `INSUREDID` = `CLIENT.ID`, `INSUREDNAME` = `CLIENT.NAME`;
  **Org ID** = `ORG-n`. **Group Business** \* = `BUSINESSGROUP` (`ID` → `GROUPBUSINESSID`, `NOTE` → `GROUPBUSINESS`).
  Wajib isi: *"Value cannot be blank"* (tangkapan layar Pega).
- **Owner** = `CREATEOP` (akun pelaku saat dibuat), **Create Date** = `CREATEDATE`; keduanya diisi sekali saat dibuat.
  **Description** bebas. **Territory** dihapus.
- `ID` = `ASM-SFAGIS-WORK-ACCOUNT ACC-<SEQ_T_M_ACCOUNT>`; di layar `ACC-n`. Nomor yang sudah terpakai dilewati.
- Pasangan Insured + Group Business yang sudah dimiliki akun yang ada ditolak (409).
- Tampilan sama dengan Kelola User (permintaan work owner 04-10-2026): tema disalin ke token `--acc-*` di
  `accounts.css`, tidak menumpang ke kelas inti.

## Migrasi

Rentang `840-879`:

| Berkas | Isi |
| --- | --- |
| `840_t_m_account_kolom` | `T_M_ACCOUNT` + `CREATEDATE DATE`, `CREATEOP VARCHAR2(100)`, `DESCRIPTION VARCHAR2(4000)` — tanpa default, baris lama kosong |
| `841_t_m_account_pk` | `PK_T_M_ACCOUNT` atas `ID` (DEV 04-10-2026: 18.658 ID, semua berbeda, nol kosong) |
| `842_seq_t_m_account` | `SEQ_T_M_ACCOUNT` mulai dari nomor ACC terbesar + 1, dihitung saat migrasi berjalan (blok sequence-dari-kueri) |

Karena 840 dan 841 mengubahnya, `T_M_ACCOUNT` TIDAK terdaftar "Tabel warisan" di bawah, dan skema uji
(`uji/skemauji`) membuat tiruannya sebelum migrasi. Slot menu `994`: satu `UPDATE DIMIGRASI` baris modul ini, nol
`INSERT`. Barisnya sendiri dibuat migrasi inti `908_m_nav_menu_accounts.sql`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/accounts/...
npx vitest run modul/accounts
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `BUSINESSGROUP` | tabel warisan POOLDATA; sumber Group Business (`ID`, `NOTE`); dibaca saja |

`CLIENT` (sumber Insured Name dan Org ID) juga dibaca saja, tetapi TIDAK terdaftar di atas: migrasi modul Company
Detail mengubahnya, dan skema uji membuat tiruannya.
