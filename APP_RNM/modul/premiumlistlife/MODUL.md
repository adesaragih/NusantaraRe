# Modul `premiumlistlife` — PremiumList Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `premiumlistlife` |
| Folder korpus | `PremiumList Life` |
| GROUPMENU | `TREATY` |
| Pemilik | `@PEMILIK-PREMIUMLISTLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `050-099` |
| Slot menu | `954-955` |
| Prefix rute API | `/api/polis-life` |
| Kontrak disediakan | `kontrak.PembacaPolis` (dipakai `claimlife`) |
| Kontrak dipakai | — |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `migrations/` `modul.go` — paket Go `nusantarare/modul/premiumlistlife/backend/...` |
| `frontend/` | `pages/` `labels.ts` `api.ts` `menu.ts` `rute.tsx` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, LAPORAN, OQ — dulu `.scratch/premiumlist-life/` |

## Migrasi

Rentang `050-099` *(ralat 30-09-2026: dulu tertulis `050–079` di panduan deploy)*, terpakai
`050–058`. Slot menu `954-955` tidak terpakai: baris modul ini sudah
`DIMIGRASI = '1'` sejak 900, dan menu datar (30-09-2026) tidak punya butir — slot hanya menyalakan
`DIMIGRASI` (`APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6). Nama berkas migrasi yang sudah ada tidak pernah diubah: `T_MIGRASI` mencatat nama.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/premiumlistlife/...
go test -tags db ./modul/premiumlistlife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/premiumlistlife
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` di akar repo, bab 8).

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Penjaganya berlaku untuk setiap modul; yang KHUSUS modul ini dinyatakan di sini, supaya
mengubahnya tidak pernah menyunting berkas di luar folder ini. Judul yang tidak ada berarti modul
ini tidak menyatakan apa pun untuk jenis itu. Nilai di dalam `` ` `` dibaca apa adanya.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang dokumen STRUKTUR modul ini gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris =
kepemilikan tabel berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `M_TEMPUPLOADLIFE` | tabel warisan penampung unggahan CSV, ditulis `RDBList/InsertDataUploadLife.xml` di sistem lama; dibaca tiket 04, tidak dibuat |

### Pesan verbatim yang bukan nama orang

Konstanta teks yang cocok dengan pola nama orang tetapi BUKAN nama orang (`TestNolNamaOrangDiKode`).
Nilainya DIBACA dari sumber konstantanya, tidak diketik ulang.

| Paket | Konstanta | Alasan |
| --- | --- | --- |
| `backend/models` | `PesanNamaTertanggung` | ValidasiUploadPL_act `local.err3` - pesan kolom NAME_OF_INSURED. |
| `backend/models` | `PesanPolicyHolder` | ValidasiUploadPL_act `local.err17` - pesan rujukan master POLICY HOLDER. |

## Brief acuan

`PROMPT-IMPLEMENTASI-MODUL-PREMIUMLIST-LIFE.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` (folder `OUTPUT_HASIL_RNM\`).
