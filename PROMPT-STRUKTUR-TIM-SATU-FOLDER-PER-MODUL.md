# PROMPT — STRUKTUR TIM: **SATU FOLDER PER MODUL** *(backend + frontend + dokumen)* untuk pengembang fullstack *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang `main`)*

> Keputusan work owner 30-09-2026: repo dibagi ke beberapa orang, **satu orang fullstack mengerjakan satu modul**, dan saat di-push balik
> **tidak perlu merge manual dan tidak ada konflik**. Bentuk B *(`inti/` + `modul/<nama>/`)* dan tabel nama modul sudah ada; brief ini
> menyatukan frontend dan dokumen ke folder modul, dan menutup setiap berkas bersama yang masih harus disunting tiap modul.
> **Nol perubahan perilaku**; angka uji sama sebelum dan sesudah *(boleh bertambah hanya untuk penjaga baru, disebut)*. Pindah dengan `git mv`.
> **Hanya sesi ini** yang bekerja di `OUTPUT_HASIL_RNM` selama refactor.

## 0. BENTUK TUJUAN

```
APP_RNM/
  go.mod  package.json  package-lock.json  vite.config.ts  tsconfig.json   ← milik tim inti
  cmd/api/
  inti/
    backend/ …           (isi inti/ Go sekarang)
    frontend/ …          (isi frontend/src/inti sekarang: Shell, ui, klien, lib, tema, labels bersama)
  modul/
    claimlife/
      backend/           models/ repository/ services/ handlers/ migrations/ modul.go
      frontend/          pages/ components/ labels.ts api.ts menu.ts rute.tsx *.test.ts
      docs/              spec, tiket, grilling, PARITAS, LAPORAN, OQ  (isi .scratch/claim-life/)
      MODUL.md           pemilik, rentang migrasi, prefix rute API, kontrak yang disediakan/dipakai
    premiumlistlife/  komiteclaimlife/  treatycontractout/   (bentuk sama)
  frontend/              index.html, main.tsx, App.tsx perakit — tanpa kode modul
docs/bersama/            ADR, CONTEXT, panduan lintas modul
.github/CODEOWNERS
```

## 1. TIGA MASALAH TEKNIS YANG HARUS DITUTUP *(diukur asisten 30-09-2026)*

| Masalah | Bukti | Solusi |
| --- | --- | --- |
| **Resolusi paket npm** dari folder di luar `frontend/` | `package.json` dan `node_modules` kini di `APP_RNM/frontend/`; berkas di `APP_RNM/modul/<nama>/frontend/` mencari `react` ke atas dari foldernya sendiri dan **tidak** menemukannya | pindahkan `package.json`, `package-lock.json`, `vite.config.ts`, `tsconfig.json` ke **`APP_RNM/`** *(root Vite = `APP_RNM/frontend`, `resolve` dan `server.fs.allow` mencakup `APP_RNM/modul` dan `APP_RNM/inti`)*; `tsconfig` `include` = `frontend`, `inti/frontend`, `modul/*/frontend`; `vitest` `include` sama; `npm install` sekali di `APP_RNM/` |
| **Go** membaca folder `frontend/` dan `docs/` di dalam modul | modul Go `nusantarare` di `APP_RNM/go.mod` | aman: Go mengabaikan folder tanpa `.go`; paket pindah ke `modul/<nama>/backend/...`; `go:embed migrations/*.sql` *(kini di `modul/<nama>/modul.go`)* ikut pindah bersama foldernya |
| **Uji yang membaca dokumen `.scratch`** | 5 berkas uji merujuk `.scratch/{claim-life,premiumlist-life,komite-claim-life,treaty-contract-out,inti}` | jalur diperbarui ke `modul/<nama>/docs/` dan `docs/bersama/`; `.scratch/` modul yang **belum** dibangun *(claim-facin, nb-treaty-in, …)* tetap di tempatnya sampai modulnya dimulai |

## 2. MENUTUP BERKAS BERSAMA — supaya menambah atau mengubah modul tidak menyentuh berkas milik orang lain

| Berkas bersama sekarang | Menjadi |
| --- | --- |
| `modul/daftar.go` *(daftar modul Go)* | **dibangkitkan**: `go generate` memindai `modul/*/backend/modul.go` dan menulis `inti/backend/daftar_gen.go`; berkas hasil **di-commit** dan dijaga uji "hasil bangkit = isi folder" *(tanpa `go generate`, uji merah menyebut perintahnya)* — tiap modul hanya punya `modul.go` sendiri |
| `frontend/src/modul/daftar.ts` | `import.meta.glob('../../modul/*/frontend/menu.ts', { eager: true })` dan padanannya untuk `rute.tsx` — daftar terbentuk dari folder, **nol** baris per modul di berkas bersama |
| `inti` `labels.ts` `MODUL` *(nama kelompok)* | nama kelompok pindah ke `menu.ts` modul masing-masing; `inti` hanya memuat label kerangka |
| Nomor migrasi | tabel rentang di **`MODUL.md`** tiap modul; penjaga uji: setiap berkas migrasi di `modul/<nama>/backend/migrations` berada di rentang modulnya, dan dua modul tidak berbagi rentang *(Claim Life 001–029, Komite 030–049, PremiumList 050–099, Treaty 300–319, inti 900–949; modul baru mengambil blok 100 berikutnya yang kosong)* |
| Penjaga statik lintas modul | di `inti/backend/penjaga`, memindai `modul/*` secara umum; **nol** daftar nama modul di dalam penjaga |
| `go.mod`, `go.sum`, `package.json`, `package-lock.json` | milik tim inti *(CODEOWNERS)*; tim modul mengajukan pustaka baru lewat pull request ke tim inti |

## 3. ATURAN KEPEMILIKAN

- **`.github/CODEOWNERS`**: `/APP_RNM/modul/<nama>/  @<pemilik-modul>` per modul; `/APP_RNM/inti/`, `/APP_RNM/cmd/`, `/APP_RNM/frontend/`, berkas konfigurasi
  root, dan `/docs/bersama/` → `@<tim-inti>`. Nama akun **dibiarkan sebagai penanda** `@PEMILIK-CLAIMLIFE` dst. — work owner mengisinya.
- **Penjaga impor** *(sudah ada, diperluas ke frontend)*: `modul/X` tidak boleh mengimpor `modul/Y` — baik Go maupun TypeScript; lintas modul hanya lewat
  `inti/.../kontrak`.
- **`MODUL.md`** tiap modul: pemilik, rentang migrasi, prefix rute API *(`/api/klaim-life`, `/api/polis-life`, `/api/komite`, `/api/treaty-contract-out`)*,
  kontrak yang disediakan dan dipakai, cara menjalankan uji modul itu saja.

## 4. URUTAN — satu commit per paket, uji hijau dan angka sama di setiap commit

| # | Paket |
| ---: | --- |
| 1 | konfigurasi npm/Vite/TS/Vitest naik ke `APP_RNM/`; `npm install` ulang; uji dan build sama |
| 2 | `inti/` dipecah `inti/backend` + `inti/frontend` |
| 3 | per modul *(Treaty, PremiumList, Komite, Claim Life)*: `backend/` + `frontend/` + `docs/` + `MODUL.md`; jalur uji yang membaca dokumen diperbarui |
| 4 | daftar modul dibangkitkan *(Go `go generate`, frontend `import.meta.glob`)*; label kelompok ke modul |
| 5 | penjaga rentang migrasi + penjaga impor lintas modul frontend |
| 6 | `CODEOWNERS` + `docs/bersama/` + `PANDUAN-TIM-PER-MODUL.md` *(alur: cabang per modul → commit di folder sendiri → pull request → CI semua modul → merge; cara menambah modul baru: salin templat folder, isi `MODUL.md`, ambil rentang migrasi)* + templat folder modul kosong `modul/_templat/` |

Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`)*, `gofmt`, `tsc`, `vitest`, `vite build`. `PANDUAN-MENJALANKAN` diperbarui
*(perintah `npm` kini dari `APP_RNM/`)*.

## 5. LAPORAN

Satu pesan: tabel **paket → commit → berkas dipindah → angka uji**; daftar berkas yang masih bersama beserta pemiliknya; bukti bahwa **menambah modul
templat baru tidak mengubah satu pun berkas di luar `modul/<nama>/`** *(uji coba dengan `modul/_contoh` yang dibuat lalu dihapus)*; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 30 September 2026 dari pengukuran `APP_RNM` (letak `package.json`/`vite.config.ts`/`tsconfig.json`, `go:embed` per modul, lima uji yang membaca
`.scratch`) dan keputusan work owner: fullstack, satu orang satu modul.*
