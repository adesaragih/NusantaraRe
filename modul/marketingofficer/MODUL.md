# Modul `marketingofficer` — Marketing Officer

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
03-10-2026: *"buatkan di modul baru dengan nama marketingofficer, gunanya untuk insert update table
marketingofficer; ingat jangan di inti, buat modul baru"*. Padanan Pega: form `InputMarketingOfficer` di NB FacIn,
RNW Fac In, dan Endorsment Fac In (popup `InputKotaMOName`), bukan menu portal; menu sendiri adalah keputusan work
owner. Baris menunya dibuat migrasi inti `906` (slot menu modul tidak boleh membuat kelompok), dan `Folder korpus`
di bawah = label menu `M_NAV_MENU.LABEL` (keputusan work owner: "Marketing Officer", kelompok MASTER).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — diambil dari
cadangan `760-899` dan `990-999`.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `marketingofficer` |
| Folder korpus | `Marketing Officer` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MARKETINGOFFICER` |
| Status | dimigrasi |
| Rentang migrasi | `760-799` |
| Slot menu | `990-990` |
| Prefix rute API | `/api/marketing-officer` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-MARKETINGOFFICER.md` — peta tabel warisan dan sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `marketingofficer.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas)

- Tulis langsung ke `POOLDATA.MARKETINGOFFICER`: nol tabel baru, nol DDL, prosedur `PEGA_MARKETINGOFFICER` tidak
  dipanggil. Nol hapus — nonaktif = `MOSTATUS` 2.
- `AKSES_LOGIN` = `M_LOGIN_GO.LOGIN_ID` (akun aktif, maks. 50 karakter). `CLIENTID` (Marketing Code) =
  `M_LOGIN_GO.CONTACT_ID`, atau `CLIENTID` lama bila akun itu sudah punya baris MO; tidak pernah berubah.
- Satu baris AKTIF per Marketing Code dan per akun. Leader harus baris `LEADER` aktif. Sub Branch dari `BRANCH`.
- Layar (permintaan work owner 03-10-2026): halaman depan = daftar leader; dari leader dibuka anggotanya; setiap MO
  punya log perubahan dari `MARKETINGOFFICER_LOG`. Tampilan sama dengan Kelola User (tema disalin ke CSS modul).
- Rinciannya: `docs/STRUKTUR-TABEL-MARKETINGOFFICER.md` dan dokumentasi paket `backend/services`.

## Migrasi

Rentang `760-799`: `760_marketingofficer_log.sql` memperbaiki tabel warisan `MARKETINGOFFICER_LOG` (izin work owner
03-10-2026) - tambah `LOG_TIME` dan `AKSES_LOGIN`, trigger warisan tidak disentuh; nol tabel baru. Karena 760
mengubahnya, `MARKETINGOFFICER_LOG` TIDAK terdaftar "Tabel warisan" di bawah, dan skema uji (`uji/skemauji`) membuat
tiruannya sebelum migrasi. Slot menu `990`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`. Barisnya sendiri
dibuat migrasi inti `906_m_nav_menu_marketingofficer.sql`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/marketingofficer/...
npx vitest run modul/marketingofficer
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `MARKETINGOFFICER` | tabel warisan POOLDATA milik Pega; modul ini menambah dan mengubah barisnya, tidak pernah membuat atau mengubah strukturnya (keputusan work owner 03-10-2026: nol tabel baru) |
| `BRANCH` | tabel warisan POOLDATA, sumber Branch dan Sub Branch; dibaca saja |
