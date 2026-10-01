# PROMPT — STRUKTUR TIM: **SATU FOLDER PER MODUL** *(backend + frontend + dokumen)* untuk pengembang fullstack *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang `main`)*

> Keputusan work owner 30-09-2026: repo dibagi ke beberapa orang, **satu orang fullstack mengerjakan satu modul**, dan saat di-push balik
> **tidak perlu merge manual dan tidak ada konflik**. Bentuk B *(`inti/` + `modul/<nama>/`)* dan tabel nama modul sudah ada; brief ini
> menyatukan frontend dan dokumen ke folder modul, dan menutup setiap berkas bersama yang masih harus disunting tiap modul.
> **Nol perubahan perilaku**; angka uji sama sebelum dan sesudah *(boleh bertambah hanya untuk penjaga baru, disebut)*. Pindah dengan `git mv`.
> **Hanya sesi ini** yang bekerja di `OUTPUT_HASIL_RNM` selama refactor.

## 0.1 RALAT 30-09-2026 *(asisten, sesudah mengukur `main` @ `420bc66`)* — **bab ini mengalahkan bab lain bila bertentangan**

Gambaran pohon lengkap *(akar repo sampai 20 folder modul)* yang disetujui work owner hanya terbentuk bila lima butir berikut ikut dikerjakan.

| # | Masalah di brief lama | Bukti | Ganti menjadi |
| ---: | --- | --- | --- |
| R1 | Hanya 4 folder modul yang dipindah; 16 modul lain baru muncul saat dimulai, jadi `CODEOWNERS` dan rentang migrasi tidak bisa ditetapkan sekali di awal | `APP_RNM/modul/` kini 4 folder | paket 6 **juga** membuat **16 folder kerangka** dari `_templat/`, satu per folder korpus yang belum dibangun, nama dari tabel nama *(`claimfacin`, `claimprop`, `claimnonprop`, `komiteclaimfacin`, `komiteclaimprop`, `komiteclaimnonprop`, `nbfacin`, `rnwfacin`, `endorsmentfacin`, `nbtreatyin`, `edmtreatyin`, `treatyin`, `treatyinadjustment`, `endorsementlife`, `mastercontractretrolife`, `masterproductnamelife`)*. Isi kerangka **hanya** `MODUL.md` *(pemilik penanda, rentang R2, slot menu R3, status "belum dimigrasi")* dan `docs/` dari `.scratch/<nama-panjang>/` bila ada. **Nol** `modul.go`, nol halaman, nol menu — kerangka tanpa `backend/modul.go` tidak terdaftar, sehingga `TestIsiAwalMenuDimigrasiSamaDenganModulBackend` tetap hijau |
| R2 | "Modul baru mengambil blok 100 berikutnya" **tidak muat** | pelari migrasi mengurutkan nama berkas **sebagai teks** di semua sumber sekaligus *(`inti/migrasi/migrasi.go:166` `sort.Strings`)*; nomor 4 digit akan salah urut *(`1000_` sebelum `101_`)*; ruang kosong 3 digit tinggal 100–299, 320–899, 950–999 = 7 blok ratusan untuk 16 modul | rentang **40** nomor, ditetapkan sekarang, urut **hulu ke hilir** supaya migrasi modul hilir yang merujuk tabel modul hulu selalu berjalan sesudahnya — tabel R2 di bawah. PremiumList **050–099** *(PANDUAN-DEPLOY masih menulis 050–079 — diperbarui)* |
| R3 | Isi menu modul baru wajib di `inti/migrations/901–949` | PANDUAN-DEPLOY bab 6 dan `TestMenuHanyaDiMigrasiInti` | itu folder bersama: dua pengembang menambah `901_…` bersamaan = bentrok, dan pengembang modul menyentuh folder inti. Ganti: tiap modul punya **dua slot menu** di rentang **950–999**, berkas tinggal di **folder migrasi modulnya sendiri** *(mis. `modul/nbfacin/backend/migrations/962_menu_nbfacin.sql`)*. Urutan teks menjamin slot 95x berjalan **sesudah** `900_m_nav_menu` di skema baru. Bentuk SQL persis bab 6 PANDUAN-DEPLOY. Penjaga diganti: baris `M_NAV_MENU` hanya di `900` dan di slot menu milik modul itu *(dibaca dari `MODUL.md`, nol nama modul di penjaga)*; 900–949 tetap milik `inti` saja |
| R4 | `go generate` menulis daftar modul dari folder | `modul/daftar.go:60-66` **menyambung kontrak dengan tangan** *(PembacaPolis PremiumList ke Claim Life, KlaimKomite Claim Life ke Komite)* — pembangkit tidak dapat menebak sambungan itu | tiap `backend/modul.go` menyatakan kontrak yang **disediakan** dan yang **dibutuhkan** lewat `inti/backend/kontrak`; perakit di `inti` menyambung menurut jenis antarmuka; kontrak dibutuhkan tanpa penyedia = galat saat mulai yang menyebut kontrak dan modulnya *(bukan nil diam-diam)*. Uji: sambungan hasil perakit = sambungan `daftar.go` sekarang. Nama paket `modul.go` = `backend` di semua modul, jadi berkas bangkitan **wajib memakai alias impor** = nama modul |
| R5 | Folder yang ada tetapi tidak disebut | `APP_RNM/uji/` *(skemauji, lintasmodul)*, `APP_RNM/pkg/utils`, `inti/menu/`, `inti/migrations/` | `uji/` dan `pkg/` **tetap di tempatnya**, milik tim inti di `CODEOWNERS`. `inti/menu/` dan `inti/migrations/` ikut ke `inti/backend/`. Dokumen tersegel *(`grilling-ronde-*`, `VERIFIKASI-*`, `KOREKSI-*`)* dipindah **hanya** dengan `git mv`; isinya byte-identik, dibuktikan hash sebelum dan sesudah di laporan |

**Tabel R2 — rentang migrasi dan slot menu semua modul** *(ditulis juga ke `MODUL.md` masing-masing; nama berkas migrasi yang sudah ada tidak berubah)*

| GROUPMENU | Modul backend | Rentang migrasi | Slot menu |
| --- | --- | --- | --- |
| KLAIM | `claimlife` | 001–029 | 950–951 |
| KLAIM | `komiteclaimlife` | 030–049 | 952–953 |
| TREATY | `premiumlistlife` | 050–099 | 954–955 |
| MASTER | `treatycontractout` | 300–319 *(tetap kosong, tco4)* | 956–957 |
| MASTER | `mastercontractretrolife` | 100–139 | 958–959 |
| MASTER | `masterproductnamelife` | 140–179 | 960–961 |
| FACULTATIVE | `nbfacin` | 180–219 | 962–963 |
| FACULTATIVE | `rnwfacin` | 220–259 | 964–965 |
| FACULTATIVE | `endorsmentfacin` | 260–299 | 966–967 |
| TREATY | `nbtreatyin` | 320–359 | 968–969 |
| TREATY | `edmtreatyin` | 360–399 | 970–971 |
| TREATY | `treatyin` | 400–439 | 972–973 |
| TREATY | `treatyinadjustment` | 440–479 | 974–975 |
| TREATY | `endorsementlife` | 480–519 | 976–977 |
| KLAIM | `claimfacin` | ~~520–559~~ **560–599** *(ralat 01-10-2026)* | 978–979 |
| KLAIM | `claimprop` | ~~560–599~~ **520–559** *(ralat 01-10-2026: skema Claim Fac In bergantung pada tabel lini PROP, jadi Prop harus berjalan lebih dulu)* | 980–981 |
| KLAIM | `claimnonprop` | 600–639 | 982–983 |
| KLAIM | `komiteclaimfacin` | 640–679 | 984–985 |
| KLAIM | `komiteclaimprop` | 680–719 | 986–987 |
| KLAIM | `komiteclaimnonprop` | 720–759 | 988–989 |
| — | `inti` | 900–949 | — |
| — | cadangan, dibagi tim inti lewat pull request | 760–899 dan 990–999 | — |

Penjaga R2: setiap berkas migrasi berada di rentang atau slot menu modulnya; dua `MODUL.md` tidak berbagi nomor; nomor selalu 3 digit.

**Paket di bab 4 bertambah**: paket 4 memuat R4; paket 5 memuat penjaga R2 dan R3; paket 6 memuat R1. Laporan menambah bukti: skema uji dari nol menjalankan `900` sebelum slot `95x`, dan uji `_contoh` memakai satu rentang cadangan lalu dihapus.

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
