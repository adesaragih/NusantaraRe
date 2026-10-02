# Modul `facout` — Fac Out (retrosesi keluar)

⛔ **DIPARKIR — bukan di `APP_RNM/modul/`.** Folder ini disiapkan untuk menjadi `APP_RNM/modul/facout/`, tetapi
modul di luar folder korpus baru lolos gerbang setelah pull request tim inti. Alasan dan daftar
permintaannya: `APP_RNM/modul/nbfacin/docs/KEPUTUSAN-30-09-2026.md`. Pindahkan folder ini utuh ke sana
setelah disetujui.

⚠️ **Kerangka — belum dimigrasi.** Modul di luar dua puluh folder korpus
(`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 5), dibuat atas keputusan work owner 30-09-2026: Fac Out
menjadi **modul baru**, bukan bagian `nbfacin`. Belum ada kode: tanpa `backend/modul.go` modul ini tidak
terdaftar.

⚠️ **Rentang migrasi `760-799` dan slot menu `990-991` di bawah adalah USULAN** dari rentang cadangan
(`760-899`, `990-999`) — belum disetujui tim inti. Yang juga masih harus diajukan ke tim inti sebelum
modul ini mendapat layar (bab 5): kelompok `M_NAV_MENU`-nya (isi awal 900 hanya memuat dua puluh folder
korpus) dan barisnya di `frontend/katalogKorpus.ts`.

📌 **Tidak ada folder korpus `Fac Out` sendiri.** Rule Fac Out tinggal di dalam folder korpus siklus
Fac In (`NB FacIn`, dan jalur endorsement di `Endorsment Fac In`) — lihat
`docs/09-facout/01-temuan-dan-rancangan-facout.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `facout` |
| Folder korpus | — (tidak ada; rule Fac Out di dalam `NB FacIn` dan `Endorsment Fac In`) |
| GROUPMENU | `FACULTATIVE` |
| Pemilik | `@PEMILIK-FACOUT` |
| Status | belum dimigrasi |
| Rentang migrasi | `760-799` |
| Slot menu | `990-991` |
| Prefix rute API | — (ditetapkan spec modul ini) |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | 14 tiket + indeks (`issues/`), temuan & rancangan (`09-facout/`), 8 audit sumber angka tiket (`10-audit/`) — dipindah dari `jefri/OUTPUT FIX/` 30-09-2026, isi tidak diubah. Asal dan tujuan tautan lama: `docs/PETA-ASAL.md` |
| `backend/` | `models/` `repository/` `services/` `handlers/` `migrations/` `modul.go` — lahir saat modul dimulai |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `menu.ts` `rute.tsx` — lahir saat modul dimulai |

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/facout/...
npx vitest run modul/facout
```

## Pernyataan untuk penjaga

Tambahkan bab `###` di sini bila modul ini membutuhkannya — jenis dan bentuk tabelnya di
`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 6 (contoh lengkap: `modul/claimlife/MODUL.md`).
