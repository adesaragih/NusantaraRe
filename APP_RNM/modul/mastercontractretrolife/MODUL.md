# Modul `mastercontractretrolife` — Master Contract Retro Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `mastercontractretrolife` |
| Folder korpus | `Master Contract Retro Life` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERCONTRACTRETROLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `100-139` |
| Slot menu | `958-959` |
| Prefix rute API | — (ditetapkan spec modul ini) |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` — paket Go `nusantarare/modul/mastercontractretrolife/backend/...` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `tampilan.ts` `mcrl.css` `menu.ts` `rute.tsx` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, RALAT, OQ, STRUKTUR — dulu `.scratch/master-contract-retro-life/` |

## Migrasi

Rentang `100-139` **tetap kosong**: K1 (`docs/RALAT-DEV-30-09-2026.md`, preseden tco4) — modul ini menulis
dan membaca lima tabel warisan `POOLDATA`, nol tabel baru, nol DDL (`TestMCRLNolMigrasiDiRentang`,
`TestMCRLNolDDL`). Peta tabelnya: `docs/STRUKTUR-TABEL-MASTER-CONTRACT-RETRO-LIFE.md`.

Slot menu `958-959`: `backend/migrations/958_menu_mastercontractretrolife.sql` (+ `_down`) — satu
`UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1'` baris modul ini, nol `INSERT` (menu datar 30-09-2026,
`APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6). ⛔ **`-migrate` dijalankan work owner**, bukan sesi
pengembang; sampai 958 dijalankan, tombol menu tetap "belum dimigrasi" di basis data yang sudah berjalan.

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Judul yang tidak ada berarti modul ini tidak menyatakan apa pun untuk jenis itu. Nilai di dalam
`` ` `` dibaca apa adanya.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang dokumen STRUKTUR modul ini gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris =
kepemilikan tabel berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `TREATYYEAR_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |
| `TREATYCONTRACT_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |
| `TREATYREINSURER_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |
| `TREATYSECURITYREINSURER_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |
| `TREATYBUSINESS_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/mastercontractretrolife/...
go test -tags db -p 1 ./modul/mastercontractretrolife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/mastercontractretrolife
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` di akar repo, bab 8).
