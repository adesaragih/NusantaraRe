# Panduan deploy dan git per modul — APP_RNM bentuk B

Ditulis 30 September 2026, sesudah refactor bentuk B (`..\PROMPT-REFACTOR-BENTUK-B-MODUL.md`,
paket 1–8). Untuk siapa: pengembang yang memegang satu modul, dan siapa pun yang menyalakan
aplikasi dengan sebagian modul saja.

## 1. Bentuknya dalam satu layar

```
APP_RNM/
  cmd/api/            memasang modul dari daftar; MODUL_AKTIF; /healthz dan /api/modul-aktif
  inti/               SATU-SATUNYA kode bersama — tidak pernah mengimpor modul
    kontrak/          antarmuka lintas modul, TANPA implementasi (PembacaPolis, KlaimKomite)
    penjaga/          uji penjaga yang berlaku untuk SELURUH aplikasi
  modul/
    daftar.go         daftar modul: merakit Service tiap modul dan menyambung inti/kontrak
    claimlife/        models/ repository/ services/ handlers/ migrations/ modul.go
    premiumlist/      …
    komite/           …
    treaty/           … (tanpa migrations/: tco4, memakai tabel warisan)
  uji/                penunjang uji netral: skemauji (skema Oracle tiruan), lintasmodul
  frontend/src/
    inti/             Shell, KelompokMenu, PaletMenu, ui/dasar, klien.ts, lib/, hooks/, store/, labels.ts
    modul/daftar.ts   daftar modul frontend: merakit menu dan rute
    modul/<nama>/     pages/ components/ labels.ts api.ts menu.ts rute.tsx
    App.tsx           memasang modul yang AKTIF; Beranda.tsx layar awal aplikasi
```

**Aturan impor** — ditegakkan uji, bukan kesepakatan:

| Dari | Boleh mengimpor | Penjaga |
| --- | --- | --- |
| `modul/X/...` (Go) | `inti/...`, `modul/X/...`; berkas uji juga `uji/skemauji` ¹ | `inti/penjaga/impor_lintas_modul_test.go` |
| `inti/...` (Go) | `inti/...` saja | idem |
| `cmd/...` (Go) | `inti/...` dan daftar `nusantarare/modul` — tidak pernah satu modul langsung | idem |
| `frontend/src/modul/X/**` | `inti/**`, `modul/X/**` | `frontend/src/inti/lapisan.guard.test.ts` |
| `frontend/src/inti/**` | `inti/**` saja | idem |
| `modul/daftar.go`, `frontend/src/modul/daftar.ts`, `App.tsx`, `Beranda.tsx`, `uji/` | apa pun — merekalah tempat yang mengenal semua modul | — |

¹ `uji/skemauji` SENGAJA mengenal semua modul: ia membangun skema uji **utuh** (migrasi semua modul,
fixture Claim Life dan PremiumList). Karena itu ia satu-satunya jalur lintas modul yang disahkan —
hanya untuk berkas uji, dan hanya lewat paket itu.

Modul yang membutuhkan modul lain **tidak mengimpornya**: ia meminta antarmuka di `inti/kontrak`,
dan `modul/daftar.go` menyambungkannya. Dua sambungan yang ada hari ini:

| Antarmuka | Disediakan | Dipakai | Butir |
| --- | --- | --- | --- |
| `kontrak.PembacaPolis` | PremiumList (`services.PembacaPolis`) | Claim Life | pl4/av |
| `kontrak.KlaimKomite` | Claim Life (`services.KlaimUntukKomite`) | Komite | km3 |

## 2. Deploy memilih modul — `MODUL_AKTIF`

Satu biner memuat keempat modul; `MODUL_AKTIF` memilih yang **dipasang** saat menyala.

```powershell
$env:MODUL_AKTIF = 'claimlife,komite'   # dipisah koma; spasi dan huruf besar diabaikan
go run ./cmd/api                        # atau .\bin\api.exe
```

| Nama | Modul | Awalan rute |
| --- | --- | --- |
| `claimlife` | Claim Life | `/api/klaim-life`, `/api/peserta-life`, `/api/penyakit-life`, `/api/dokumen` |
| `premiumlist` | PremiumList Life | `/api/polis-life` |
| `komite` | Komite Claim Life | `/api/komite` |
| `treaty` | Treaty Contract Out | `/api/treaty-contract-out` (+ pekerja latar lampiran) |
| *(aplikasi)* | selalu ada | `/healthz`, `/api/modul-aktif` |

- **Kosong (bawaan) = semua modul.**
- Modul yang tidak disebut: rutenya **tidak didaftarkan** — jawabannya 404 berbadan JSON
  `{"galat":"modul <nama> tidak aktif di proses ini (MODUL_AKTIF)"}`, supaya layar tidak menyangka
  backend mati — **pekerja latarnya tidak jalan**, dan **menunya tidak tampil** — sidebar, palet Ctrl+K, dan kartu Beranda. Frontend membaca
  daftar modul aktif dari `GET /api/modul-aktif` (`{"modul":["claimlife","komite"]}`); **tidak ada env
  Vite** untuk ini, jadi satu bangunan frontend melayani deploy mana pun.
- Nama yang salah ketik **menolak menyala**: backend berhenti dengan pesan yang menyebut nama itu dan
  nama-nama yang dikenal. `-migrate` / `-migrate-down` tidak membaca `MODUL_AKTIF` sama sekali.
- **Migrasi tidak ikut `MODUL_AKTIF`.** `-migrate` selalu menjalankan migrasi SEMUA modul terdaftar,
  supaya skema utuh (tabel satu modul dirujuk modul lain, dan data warisan tidak memilih modul).
- Bila `/api/modul-aktif` gagal dibaca, frontend menampilkan **semua** menu — persis perilaku sebelum
  `MODUL_AKTIF` ada. Satu pembacaan yang gagal tidak mengosongkan aplikasi.

⚠️ **Ketergantungan yang tetap ada saat sebagian modul mati:**

1. Layar Register dan Outstanding **Claim Life** mengisi panel data polis dari
   `GET /api/polis-life/ringkas` — rute **PremiumList**. Tanpa `premiumlist`, panel itu menampilkan
   galat "modul premiumlist tidak aktif"; halaman lain Claim Life tetap jalan. (Di backend, Claim Life
   membaca polis lewat `kontrak.PembacaPolis` di dalam proses, dan itu tetap tersambung apa pun
   `MODUL_AKTIF`.) Ketergantungan ini lewat **HTTP**, jadi penjaga impor tidak melihatnya.
   Menghapusnya berarti rute Claim Life sendiri yang menyajikan polis lewat `kontrak.PembacaPolis` —
   rute baru, di luar lingkup refactor bentuk B (yang melarang perubahan rute). `[pertanyaan terbuka]`
2. **Komite** menuntaskan baris klaim Claim Life lewat `kontrak.KlaimKomite` di dalam proses — tetap
   tersambung apa pun `MODUL_AKTIF`. Menyerahkan kasus ke Komite adalah tombol di layar Claim Life.
3. Kartu antrean di Beranda adalah cacah kotak masuk Claim Life; tanpa `claimlife` kartu itu tidak
   tampil dan tidak diminta **lagi sesudah daftar modul aktif terbaca**. Pada muatan pertama Beranda
   dapat mengirim empat permintaan itu sebelum `/api/modul-aktif` menjawab; jawabannya diabaikan
   begitu daftarnya tiba.

Memecah modul ke **proses berbeda** di belakang reverse proxy (per awalan rute di atas) belum
didukung: frontend membaca `/api/modul-aktif` dari satu backend saja. `[pertanyaan terbuka]`

## 3. Git per folder

Seluruh aplikasi satu repositori; `git pull` selalu menarik semuanya. Yang dapat dipisah per modul
adalah **apa yang Anda commit, lihat, dan salin**.

| Modul | Folder backend | Folder frontend |
| --- | --- | --- |
| Claim Life | `APP_RNM/modul/claimlife/` | `APP_RNM/frontend/src/modul/claimlife/` |
| PremiumList Life | `APP_RNM/modul/premiumlist/` | `APP_RNM/frontend/src/modul/premiumlist/` |
| Komite Claim Life | `APP_RNM/modul/komite/` | `APP_RNM/frontend/src/modul/komite/` |
| Treaty Contract Out | `APP_RNM/modul/treaty/` | `APP_RNM/frontend/src/modul/treaty/` |

```powershell
# Commit HANYA folder modul Anda. `-o` (--only) mengambil isi jalur itu dari pohon kerja dan
# mengabaikan apa pun yang sudah di-stage sesi lain — stage di pohon ini dipakai bersama.
git add -A -- APP_RNM/modul/claimlife APP_RNM/frontend/src/modul/claimlife
git commit -o -m "claimlife: ..." -- APP_RNM/modul/claimlife APP_RNM/frontend/src/modul/claimlife

# Riwayat dan beda satu modul saja
git log --oneline -- APP_RNM/modul/claimlife APP_RNM/frontend/src/modul/claimlife
git diff HEAD~1 -- APP_RNM/modul/komite

# Menarik perubahan: pull utuh, lalu lihat apa yang berubah di luar modul Anda
git pull
git diff --stat ORIG_HEAD -- APP_RNM/inti APP_RNM/frontend/src/inti APP_RNM/modul/daftar.go APP_RNM/frontend/src/modul/daftar.ts

# Opsional: salinan kerja yang hanya memuat inti + satu modul (git sparse-checkout, mode cone).
# Mode cone menerima FOLDER; berkas yang langsung berada di folder induknya (modul/daftar.go,
# APP_RNM/*.md) ikut dengan sendirinya. Folder dokumen di OUTPUT_HASIL_RNM\ ikut tersembunyi.
git sparse-checkout set APP_RNM/cmd APP_RNM/inti APP_RNM/uji APP_RNM/modul/claimlife APP_RNM/frontend
git sparse-checkout disable   # kembali ke salinan utuh
```

⚠️ Sparse-checkout satu modul **tidak dapat membangun** aplikasi: `modul/daftar.go` dan
`frontend/src/modul/daftar.ts` mengimpor keempat modul. Pakai untuk membaca dan menyunting; bangun
dan uji selalu dari salinan utuh.

**Aturan yang berlaku untuk setiap commit** (sama dengan gerbang refactor bentuk B):

```powershell
go build ./... ; go build -tags db ./... ; go vet ./... ; go vet -tags db ./... ; gofmt -l cmd inti modul uji
go test ./... ; go test -tags db ./...          # tanpa ORACLE_DSN: test db SKIP dengan pesan
cd frontend ; npx tsc --noEmit ; npx vitest run ; npx vite build
```

- Menyentuh `inti/`, `modul/daftar.go`, `frontend/src/inti/`, `frontend/src/modul/daftar.ts`, atau
  `App.tsx` berarti menyentuh **semua modul**: jalankan seluruh uji, dan minta tinjauan pemilik setiap
  modul (bab 4).
- Menambah kebutuhan lintas modul = menambah antarmuka di `inti/kontrak` (tanpa implementasi) dan satu
  sambungan di `modul/daftar.go` — tidak pernah impor langsung.
- **Nama berkas migrasi tidak pernah diubah**: `T_MIGRASI` mencatat nama, dan nama baru membuat
  `-migrate` menjalankannya ulang. Rentang nomor per modul:

| Modul | Rentang | Terpakai |
| --- | --- | --- |
| Claim Life | 001–029 | 001–022 |
| Komite Claim Life | 030–049 | 030 |
| PremiumList Life | 050–079 | 050–058 |
| Treaty Contract Out | 300–319 | kosong — tco4: nol tabel baru (`TestTCONolTabelBaru`) |

## 4. Pemilik folder

Pemilik adalah **sesi/tim yang memegang brief modul itu**; nama orangnya ditetapkan work owner dan
sengaja tidak ditulis di repositori.

| Folder | Pemilik | Brief acuan (di `..\`) |
| --- | --- | --- |
| `modul/claimlife/`, `frontend/src/modul/claimlife/` | sesi modul Claim Life | `PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE*.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` |
| `modul/premiumlist/`, `frontend/src/modul/premiumlist/` | sesi modul PremiumList Life | `PROMPT-IMPLEMENTASI-MODUL-PREMIUMLIST-LIFE.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` |
| `modul/komite/`, `frontend/src/modul/komite/` | sesi modul Komite Claim Life | `PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-LIFE.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` |
| `modul/treaty/`, `frontend/src/modul/treaty/` | sesi modul Treaty Contract Out | `PROMPT-IMPLEMENTASI-MODUL-TREATY-CONTRACT-OUT.md`, `PROMPT-LANJUTAN-TREATY-CONTRACT-OUT-*.md` |
| `inti/`, `uji/`, `cmd/`, `modul/daftar.go`, `frontend/src/inti/`, `frontend/src/modul/daftar.ts`, `App.tsx`, `Beranda.tsx` | **bersama** — perubahan disetujui work owner dan ditinjau pemilik setiap modul | `PROMPT-REFACTOR-BENTUK-B-MODUL.md` |

`.scratch/<modul>/` (spec, tiket, catatan) tetap di tempatnya dan dimiliki pemilik modulnya.

## 5. Menambah modul baru

Contoh nama: `endorsement` (sama di backend, frontend, dan `MODUL_AKTIF`).

**Backend**

1. `modul/endorsement/{models,repository,services,handlers}/` — hanya mengimpor `inti/...` dan
   dirinya sendiri. `services.DariDasar(dasar *inti.Dasar)` membangun `Service` di atas akar bersama.
2. `handlers.DaftarkanRute(mux, svc, stubPelaku)` mendaftarkan seluruh rute modul dengan satu awalan
   (`/api/endorsement-...`).
3. Bila bermigrasi: pilih rentang nomor yang belum terpakai (catat di tabel bab 3), taruh berkas di
   `modul/endorsement/migrations/` beserta `_down.sql`, dan tanam dengan `//go:embed migrations/*.sql`.
4. `modul/endorsement/modul.go`: `const Nama = "endorsement"`, `Baru(...)`, dan metode `Nama()`,
   `DaftarkanRute(mux)`, `JalankanPekerja(ctx)` (tanpa pekerja: `inti.TanpaPekerja()`), serta
   `SumberMigrasi()` bila bermigrasi — pola keempat modul yang ada.
5. `modul/daftar.go`: satu baris di `Rakit`, dan satu baris di `SumberMigrasi` bila bermigrasi. Butuh
   modul lain? Tambah antarmuka di `inti/kontrak` dan sambungkan di sini.

**Frontend**

6. `frontend/src/modul/endorsement/`: `pages/`, `components/`, `labels.ts`, `api.ts` (memakai
   `minta`/`mintaFormulir` dari `inti/klien`), `menu.ts` (`NAMA_ENDORSEMENT = 'endorsement'`,
   `HALAMAN_ENDORSEMENT`, `MENU_ENDORSEMENT`), `rute.tsx` (`RuteEndorsement`, menyimpan keadaan
   kasusnya sendiri).
7. `frontend/src/modul/daftar.ts`: satu baris di `MODUL_FRONTEND` dan satu anggota di union `Halaman`.
   Nama kelompok sidebar-nya sudah ada di `MODUL` (`inti/labels.ts`) bila modulnya salah satu folder
   korpus.

**Yang akan memeriksa Anda** — jalankan gerbang bab 3; yang biasanya berbunyi:

| Uji | Menangkap |
| --- | --- |
| `inti/penjaga/impor_lintas_modul_test.go`, `frontend/src/inti/lapisan.guard.test.ts` | impor lintas modul |
| `cmd/api/rakit_test.go` | rute modul nonaktif harus 404; migrasi di disk = migrasi di pelari |
| `frontend/src/modul/daftar.modulAktif.test.ts` | nama frontend ≠ `const Nama` di `modul/*/modul.go` |
| `frontend/src/modul/daftar.sinkron.test.ts`, `frontend/src/Shell.test.ts` | sidebar ↔ palet, butir menu berbukti korpus |
| `inti/penjaga/*` | higiene migrasi, kata cadangan Oracle, alamat layanan, nama orang, nama tabel telanjang |

`Shell.test.ts` mengunci jumlah butir menu (hari ini lima). Menambah butir menuntut bukti XML
korpus dan menyunting angka itu dengan alasan — itu disengaja.
