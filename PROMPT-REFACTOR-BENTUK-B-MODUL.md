# PROMPT — REFACTOR **BENTUK B**: satu `APP_RNM`, kode bersama di `inti/`, tiap modul di `modul/<nama>/` *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang `main`)*

> **RALAT NAMA 30-09-2026** *(keputusan work owner 30-09-2026, `PROMPT-REFACTOR-NAMA-MODUL.md`)*: nama folder modul di brief ini (`premiumlist`, `komite`,
> `treaty`) sudah diseragamkan menurut tabel berikut (`66d39dd` backend, `0d0495b` frontend). Isi brief di
> bawah dibiarkan sebagai catatan sejarah.

| Modul korpus | Backend Go `APP_RNM/modul/…` *(tanpa tanda hubung)* | Frontend `frontend/src/modul/…` dan `.scratch/…` | Nilai `MODUL_AKTIF` |
| --- | --- | --- | --- |
| Claim Life | `claimlife` *(tetap)* | `claim-life` | `claimlife` |
| PremiumList Life | `premiumlist` → **`premiumlistlife`** | `premiumlist` → **`premiumlist-life`** | `premiumlistlife` |
| Komite Claim Life | `komite` → **`komiteclaimlife`** | `komite` → **`komite-claim-life`** | `komiteclaimlife` |
| Treaty Contract Out | `treaty` → **`treatycontractout`** | `treaty` → **`treaty-contract-out`** | `treatycontractout` |
| *(modul berikutnya, mis. NB FacIn)* | `nbfacin` | `nb-facin` | `nbfacin` |

Aturan: **nama backend = nama dokumen `.scratch` tanpa tanda hubung**; frontend memakai nama `.scratch` persis. Nama paket Go = nama folder.

> Keputusan work owner 29-09-2026: *"bisa menggunakan bentuk B"*. Tujuan: **deploy** dapat memilih modul yang aktif, dan **push/pull git**
> hampir tidak pernah konflik karena tiap modul hanya menyentuh foldernya sendiri. **Nol perubahan perilaku**: setiap rute, layar, pesan,
> migrasi, dan uji tetap sama; buktinya adalah angka uji yang **sama persis** sebelum dan sesudah. Commit dengan jalur eksplisit atau
> `git mv` supaya riwayat berkas terbawa.

## 0. SYARAT MULAI — jangan mulai sebelum ketiganya benar

1. Sesi **giliran 17** tiga modul dan sesi **Treaty lanjutan 6** sudah selesai dan melapor.
2. Suntingan frontend sesi tampilan yang belum di-commit di `main` *(`index.html`, `App.tsx`, `PagarGalat.tsx`, `labels.ts`, `styles.css`,
   `Shell.tsx`, `dasar.tsx`, `hooks/useTema.ts`, `lib/tema*.ts`)* sudah di-commit atau dibuang oleh pemiliknya — refactor memindahkan
   berkas-berkas itu.
3. **Hanya sesi ini** yang bekerja di `OUTPUT_HASIL_RNM` selama refactor. Sesi lain ditahan sampai laporan akhir.

Langkah 0: `git status --porcelain` kosong; catat **angka dasar** *(Go tingkat atas dan total, dengan dan tanpa tag `db`; vitest berkas dan uji;
build modul; `go vet`; `tsc`)*.

## 1. KEADAAN SEKARANG *(diukur asisten 29-09-2026)*

Modul Go `nusantarare` *(go 1.22)*, paket `internal/{config,handlers,models,repository,services}`. Berkas non-uji berawalan modul:

| Paket | `tco_` | `polis_` | `komite_` | tanpa awalan *(Claim Life + bersama)* |
| --- | ---: | ---: | ---: | ---: |
| `models` | 10 | 8 | 3 | 17 |
| `repository` | 16 | 10 | 4 | 29 |
| `services` | 14 | 11 | 7 | 27 |
| `handlers` | 8 | 1 | 0 | 22 |

Frontend: `pages/{claimlife,premiumlist,komite,treaty-contract-out}` sudah per modul; `services/api.ts` **satu berkas 2.845 baris** untuk
semua modul; `App.tsx`, `Shell.tsx`, `assets/labels.ts`, `lib/daftarMenu.ts` bersama. Ketergantungan lintas modul yang diketahui:
Claim Life membaca polis PremiumList *(`PolisRingkas`, keputusan pl4/av)*; Komite memanggil fungsi `services` Claim Life *(km3, penyerahan)*.

## 2. BENTUK TUJUAN

```
APP_RNM/
  cmd/api/                 memasang modul dari daftar; MODUL_AKTIF
  inti/                    satu-satunya kode bersama
    config/  db/  pelaku/  galat/  outbox/  layanan/  penomor/  migrasi/  jejak?/  utils/
    kontrak/               antarmuka lintas modul (mis. PembacaPolis, PenyerahanKomite) — TANPA implementasi
  modul/
    claimlife/    models/ repository/ services/ handlers/ migrations/  rute.go  modul.go
    premiumlist/  …
    komite/       …
    treaty/       …
  frontend/src/
    inti/         Shell, ui/dasar, KelompokMenu, PaletMenu, klien HTTP (minta, galat, keadaanGalat), lib bersama, tema
    modul/<nama>/ pages/ components/ labels.ts  api.ts  menu.ts  rute.tsx
    App.tsx       merakit rute dan menu dari modul yang aktif
```

| Aturan | Isi |
| --- | --- |
| **Impor** | `modul/X/...` boleh mengimpor `inti/...` dan `inti/kontrak`; **tidak boleh** mengimpor `modul/Y/...` langsung. Ketergantungan lintas modul lewat antarmuka di `inti/kontrak`, disambungkan di `cmd/api` *(PremiumList menyediakan `PembacaPolis`, Claim Life memakainya; Claim Life menyediakan layanan penyerahan/jalur balik, Komite memakainya)*. Ditegakkan penjaga statik baru |
| **Pendaftaran** | tiap modul punya `modul.go` yang mendaftarkan rute, pekerja latar, migrasi *(embed foldernya sendiri)*, dan penjaga; `cmd/api` hanya memanggil daftar modul |
| **Migrasi** | berkas pindah ke `modul/<nama>/migrations/`, **nama berkas tidak berubah** *(`T_MIGRASI` mencatat nama — mengubahnya membuat `-migrate` menjalankan ulang)*; penjalan migrasi `inti/migrasi` mengumpulkan dari **semua** modul terdaftar, **tidak** bergantung `MODUL_AKTIF` *(skema selalu utuh)*; rentang nomor per modul tetap |
| **`MODUL_AKTIF`** | env baru, daftar dipisah koma, kosong = semua; modul nonaktif: rute tidak didaftarkan, pekerja latar tidak jalan, menu tidak tampil *(frontend membaca daftar modul aktif dari satu rute `GET /api/modul-aktif`, bukan env Vite kedua)*; `.env.example` dan PANDUAN diperbarui |
| **Frontend** | `api.ts` dipecah menjadi `inti/klien.ts` *(minta, galat, identitas)* + `modul/<nama>/api.ts`; label per modul pindah ke foldernya; `daftarMenu` dan rute dirakit dari `menu.ts`/`rute.tsx` tiap modul; uji sinkron sidebar ↔ palet tetap |
| **Penjaga yang ada** | penjaga lintas modul *(kata cadangan Oracle, nama orang, alamat layanan, kaskade, nol tabel baru Treaty, dst.)* dipindah ke `inti` dan memindai **seluruh** `modul/*`; penjaga milik satu modul ikut modulnya; penghitung migrasi dihitung **per modul** sehingga modul lain tidak pernah memaksa menyunting berkas uji modul ini |

## 3. URUTAN — satu commit per paket, uji hijau dan angka **sama** di setiap commit

| # | Paket | Catatan |
| ---: | --- | --- |
| 1 | `inti/` backend: config, db/Tx, pelaku, galat, outbox + pekerja, resolver layanan, penomor, penjalan migrasi, utils; `inti/kontrak` | berkas tanpa awalan dipilah: bersama → `inti`, Claim Life → paket 5; catat pemilahan di laporan |
| 2 | `modul/treaty` *(berkas `tco_*`)* | paling terisolasi, nol migrasi |
| 3 | `modul/premiumlist` *(`polis_*`, migrasi 050–079)* | menyediakan `PembacaPolis` |
| 4 | `modul/komite` *(`komite_*`, migrasi 030–049)* | memakai kontrak Claim Life |
| 5 | `modul/claimlife` *(sisanya, migrasi 001–029)* | memakai `PembacaPolis`; menyediakan kontrak untuk Komite |
| 6 | `cmd/api` + `MODUL_AKTIF` + `GET /api/modul-aktif` | uji: modul nonaktif → rutenya 404 dan menunya hilang; migrasi tetap lengkap |
| 7 | frontend `inti/` + `modul/<nama>/` | `api.ts` dipecah; `App.tsx` merakit dari modul |
| 8 | penjaga statik impor lintas modul + dokumen | `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`: cara deploy memilih modul, cara commit/pull per folder, siapa pemilik folder mana, cara menambah modul baru |

Setiap commit: `go build ./...`, `go vet ./...` *(dengan dan tanpa tag `db`)*, `go test ./...` *(dengan dan tanpa tag `db`)*, `gofmt -l`, `npx tsc
--noEmit`, `npx vitest run`, `npx vite build` — **angka uji sama dengan dasar** *(boleh bertambah hanya untuk uji penjaga baru paket 6 dan 8, disebut)*.

## 4. LARANGAN

Nol perubahan perilaku, nol perubahan nama berkas migrasi, nol perubahan nama rute atau teks layar, nol tabel baru, nol `-migrate`, nol
panggilan ke DEV. Dokumen `.scratch/<modul>/` **tidak** dipindah. Bila satu paket tidak dapat hijau tanpa mengubah perilaku, **berhenti** di paket
itu dan laporkan.

## 5. LAPORAN

Satu pesan: tabel **paket → commit → berkas dipindah → angka uji** *(dasar dan sesudah)*; daftar pemilahan berkas tanpa awalan *(inti atau Claim
Life)*; antarmuka `inti/kontrak` yang lahir; cara memakai `MODUL_AKTIF`; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 29 September 2026 dari pengukuran struktur `APP_RNM` (modul Go, paket, cacah berkas per awalan, susunan frontend, ukuran `api.ts`) dan
ketergantungan lintas modul yang tercatat di keputusan pl4/av dan km3.*
