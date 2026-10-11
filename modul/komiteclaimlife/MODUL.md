# Modul `komiteclaimlife` — Komite Claim Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `komiteclaimlife` |
| Folder korpus | `Komite Claim Life` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-KOMITECLAIMLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `030-049` |
| Slot menu | `—` |
| Prefix rute API | `/api/komite` |
| Kontrak disediakan | — |
| Kontrak dipakai | `kontrak.KlaimKomite` (disediakan `claimlife`) |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `migrations/` `modul.go` — paket Go `nusantarare/modul/komiteclaimlife/backend/...` |
| `frontend/` | `pages/KasusKomite.tsx` `labels.ts` `api.ts` `nilai.ts` `layar.ts` (modul TANPA menu) `rute.tsx` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, LAPORAN, OQ — dulu `.scratch/komite-claim-life/` |

## Migrasi

Rentang `030-049`, terpakai `030`. Tabel tangga Komite sendiri lahir di migrasi Claim Life `013`
(sebelum modul ini berdiri) dan tetap di sana: `T_MIGRASI` mencatat nama, bukan letak.

## Menu

**Tidak ada** (perintah work owner 09-10-2026: "kode menu nya di hapus dari repo, anggap menu itu tidak pernah ada,
karena digabung di menu klaim nya masing-masing"). Baris `M_NAV_MENU` (isi awal 900) dan hak akunnya dibuang migrasi
inti 949; slot menu `952-953` dilepas. Inbox Komite dipindah menjadi tabel komite di bawah inbox Claim Life
(`modul/claimlife/frontend/components/TabelKomite.tsx`); modul ini dipasang bagi pemegang menu `claimlife`
(`MODUL_DIPINJAM` frontend/App.tsx, frontend `layar.ts`) dan rutenya dipinjam (`ruteDipinjam` cmd/api/rakit.go).

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/komiteclaimlife/...
go test -tags db ./modul/komiteclaimlife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/komiteclaimlife
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` di akar repo, bab 8).

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Penjaganya berlaku untuk setiap modul; yang KHUSUS modul ini dinyatakan di sini, supaya
mengubahnya tidak pernah menyunting berkas di luar folder ini. Judul yang tidak ada berarti modul
ini tidak menyatakan apa pun untuk jenis itu. Nilai di dalam `` ` `` dibaca apa adanya.

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah (`TestKaskadeHanyaPadaRelasiTerdaftar`).

| Awalan berkas | Relasi |
| --- | --- |
| `030_` | relasi 9: roster komite. 013 (migrasi Claim Life) membuatnya TANPA kaskade (cacat), 030 memasangnya lewat ALTER. |

## Brief acuan

`PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-LIFE.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` (folder `OUTPUT_HASIL_RNM\`).
