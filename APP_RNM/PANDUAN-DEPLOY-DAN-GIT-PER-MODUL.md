# Panduan deploy dan git per modul — APP_RNM bentuk B

Ditulis 30 September 2026, sesudah refactor bentuk B (`..\PROMPT-REFACTOR-BENTUK-B-MODUL.md`,
paket 1–8). Untuk siapa: pengembang yang memegang satu modul, dan siapa pun yang menyalakan
aplikasi dengan sebagian modul saja. Diperbarui 30 September 2026: menu dari tabel `M_NAV_MENU`
(`..\PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`) — bab 2, 3, 5, dan bab 6 baru.

## 1. Bentuknya dalam satu layar

```
APP_RNM/
  cmd/api/            memasang modul dari daftar; MODUL_AKTIF; /healthz, /api/modul-aktif, /api/menu
  inti/               SATU-SATUNYA kode bersama — tidak pernah mengimpor modul
    kontrak/          antarmuka lintas modul, TANPA implementasi (PembacaPolis, KlaimKomite)
    menu/             GET /api/menu: pembaca M_NAV_MENU, pohon GROUPMENU, SaringMenuUntukPelaku
    migrations/       tabel lintas modul, rentang 900–949: M_NAV_MENU (900)
    penjaga/          uji penjaga yang berlaku untuk SELURUH aplikasi
  modul/
    daftar.go         daftar modul: merakit Service tiap modul dan menyambung inti/kontrak
    claimlife/        models/ repository/ services/ handlers/ migrations/ modul.go
    premiumlistlife/  …
    komiteclaimlife/  …
    treatycontractout/ … (tanpa migrations/: tco4, memakai tabel warisan)
  uji/                penunjang uji netral: skemauji (skema Oracle tiruan), lintasmodul
  frontend/src/
    inti/             Shell, KelompokMenu, PaletMenu, ui/dasar, klien.ts, lib/, hooks/, store/, labels.ts
    modul/daftar.ts   daftar modul frontend: merakit menu dan rute
    modul/<nama>/     pages/ components/ labels.ts api.ts menu.ts rute.tsx
    App.tsx           memasang modul yang AKTIF; Beranda.tsx layar awal aplikasi
```

**Tabel nama modul** — satu-satunya sumber, tanpa singkatan *(keputusan work owner 30-09-2026, `PROMPT-REFACTOR-NAMA-MODUL.md`)*:

| Modul korpus | Backend Go `APP_RNM/modul/…` *(tanpa tanda hubung)* | Frontend `frontend/src/modul/…` dan `.scratch/…` | Nilai `MODUL_AKTIF` |
| --- | --- | --- | --- |
| Claim Life | `claimlife` *(tetap)* | `claim-life` | `claimlife` |
| PremiumList Life | `premiumlist` → **`premiumlistlife`** | `premiumlist` → **`premiumlist-life`** | `premiumlistlife` |
| Komite Claim Life | `komite` → **`komiteclaimlife`** | `komite` → **`komite-claim-life`** | `komiteclaimlife` |
| Treaty Contract Out | `treaty` → **`treatycontractout`** | `treaty` → **`treaty-contract-out`** | `treatycontractout` |
| *(modul berikutnya, mis. NB FacIn)* | `nbfacin` | `nb-facin` | `nbfacin` |

Aturan: **nama backend = nama dokumen `.scratch` tanpa tanda hubung**; frontend memakai nama `.scratch` persis. Nama paket Go = nama folder.

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
$env:MODUL_AKTIF = 'claimlife,komiteclaimlife'   # dipisah koma; spasi dan huruf besar diabaikan
go run ./cmd/api                        # atau .\bin\api.exe
```

| Nama | Modul | Awalan rute |
| --- | --- | --- |
| `claimlife` | Claim Life | `/api/klaim-life`, `/api/peserta-life`, `/api/penyakit-life`, `/api/dokumen` |
| `premiumlistlife` | PremiumList Life | `/api/polis-life` |
| `komiteclaimlife` | Komite Claim Life | `/api/komite` |
| `treatycontractout` | Treaty Contract Out | `/api/treaty-contract-out` (+ pekerja latar lampiran) |
| *(aplikasi)* | selalu ada | `/healthz`, `/api/modul-aktif`, `/api/menu` |

- **Kosong (bawaan) = semua modul.**
- **Nama lama** (`premiumlist`, `komite`, `treaty` — sebelum 30-09-2026) **ditolak** saat menyala
  dengan pesan yang menyebut nama barunya; tidak ada dua nama untuk satu modul.
- ⚠️ **Deploy backend dan frontend bersama** untuk perubahan nama ini: `GET /api/modul-aktif` kini
  menjawab nama baru, dan bundel frontend lama (atau yang tertahan di cache peramban) tidak
  mengenalinya — menu PremiumList, Komite, dan Treaty akan hilang sampai bundelnya diganti.
- Modul yang tidak disebut: rutenya **tidak didaftarkan** — jawabannya 404 berbadan JSON
  `{"galat":"modul <nama> tidak aktif di proses ini (MODUL_AKTIF)"}`, supaya layar tidak menyangka
  backend mati — **pekerja latarnya tidak jalan**, dan **menunya tidak tampil**: `GET /api/menu` tidak
  mengirim butir modul itu (sidebar dan palet Ctrl+K, bab 6), dan rute serta kartu Beranda-nya disaring
  dari `GET /api/modul-aktif` (`{"modul":["claimlife","komiteclaimlife"]}`). **Tidak ada env Vite** untuk
  ini, jadi satu bangunan frontend melayani deploy mana pun.
- Nama yang salah ketik **menolak menyala**: backend berhenti dengan pesan yang menyebut nama itu dan
  nama-nama yang dikenal. `-migrate` / `-migrate-down` tidak membaca `MODUL_AKTIF` sama sekali.
- **Migrasi tidak ikut `MODUL_AKTIF`.** `-migrate` selalu menjalankan migrasi SEMUA modul terdaftar,
  supaya skema utuh (tabel satu modul dirujuk modul lain, dan data warisan tidak memilih modul).
- Bila `/api/modul-aktif` gagal dibaca, frontend memasang **semua** rute dan kartu Beranda — persis
  perilaku sebelum `MODUL_AKTIF` ada. Satu pembacaan yang gagal tidak mengosongkan aplikasi.
- Bila `/api/menu` gagal, **sidebar menampilkan galatnya** di tempat menu (Beranda tetap ada) — bukan
  menu kosong diam-diam. Sebelum migrasi 900 dijalankan jawabannya 503
  `{"galat":"tabel M_NAV_MENU belum ada - migrasi 900 belum dijalankan (-migrate, oleh work owner)"}`.

⚠️ **Ketergantungan yang tetap ada saat sebagian modul mati:**

1. Layar Register dan Outstanding **Claim Life** mengisi panel data polis dari
   `GET /api/polis-life/ringkas` — rute **PremiumList**. Tanpa `premiumlistlife`, panel itu menampilkan
   galat "modul premiumlistlife tidak aktif"; halaman lain Claim Life tetap jalan. (Di backend, Claim Life
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
| Claim Life | `APP_RNM/modul/claimlife/` | `APP_RNM/frontend/src/modul/claim-life/` |
| PremiumList Life | `APP_RNM/modul/premiumlistlife/` | `APP_RNM/frontend/src/modul/premiumlist-life/` |
| Komite Claim Life | `APP_RNM/modul/komiteclaimlife/` | `APP_RNM/frontend/src/modul/komite-claim-life/` |
| Treaty Contract Out | `APP_RNM/modul/treatycontractout/` | `APP_RNM/frontend/src/modul/treaty-contract-out/` |

```powershell
# Commit HANYA folder modul Anda. `-o` (--only) mengambil isi jalur itu dari pohon kerja dan
# mengabaikan apa pun yang sudah di-stage sesi lain — stage di pohon ini dipakai bersama.
git add -A -- APP_RNM/modul/claimlife APP_RNM/frontend/src/modul/claim-life
git commit -o -m "claimlife: ..." -- APP_RNM/modul/claimlife APP_RNM/frontend/src/modul/claim-life

# Riwayat dan beda satu modul saja
git log --oneline -- APP_RNM/modul/claimlife APP_RNM/frontend/src/modul/claim-life
git diff HEAD~1 -- APP_RNM/modul/komiteclaimlife

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
npx tsc --noEmit ; npx vitest run ; npx vite build   # dari APP_RNM/ - package.json di sini sejak 30-09-2026
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
| `inti` *(tabel lintas modul, `inti/migrations/`)* | 900–949 | 900 (`M_NAV_MENU` + isi awal) — isi menu HANYA di sini (bab 6) |

## 4. Pemilik folder

Pemilik adalah **sesi/tim yang memegang brief modul itu**; nama orangnya ditetapkan work owner dan
sengaja tidak ditulis di repositori.

| Folder | Pemilik | Brief acuan (di `..\`) |
| --- | --- | --- |
| `modul/claimlife/`, `frontend/src/modul/claim-life/` | sesi modul Claim Life | `PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE*.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` |
| `modul/premiumlistlife/`, `frontend/src/modul/premiumlist-life/` | sesi modul PremiumList Life | `PROMPT-IMPLEMENTASI-MODUL-PREMIUMLIST-LIFE.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` |
| `modul/komiteclaimlife/`, `frontend/src/modul/komite-claim-life/` | sesi modul Komite Claim Life | `PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-LIFE.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` |
| `modul/treatycontractout/`, `frontend/src/modul/treaty-contract-out/` | sesi modul Treaty Contract Out | `PROMPT-IMPLEMENTASI-MODUL-TREATY-CONTRACT-OUT.md`, `PROMPT-LANJUTAN-TREATY-CONTRACT-OUT-*.md` |
| `inti/`, `uji/`, `cmd/`, `modul/daftar.go`, `frontend/src/inti/`, `frontend/src/modul/daftar.ts`, `App.tsx`, `Beranda.tsx` | **bersama** — perubahan disetujui work owner dan ditinjau pemilik setiap modul | `PROMPT-REFACTOR-BENTUK-B-MODUL.md` |

`.scratch/<modul>/` (spec, tiket, catatan) tetap di tempatnya dan dimiliki pemilik modulnya.

## 5. Menambah modul baru

Contoh: **NB FacIn** — menurut tabel nama, backend dan `MODUL_AKTIF` `nbfacin`, frontend `nb-facin`
(nama dokumen `.scratch/nb-facin/` persis).

**Backend**

1. `modul/nbfacin/{models,repository,services,handlers}/` — hanya mengimpor `inti/...` dan
   dirinya sendiri. `services.DariDasar(dasar *inti.Dasar)` membangun `Service` di atas akar bersama.
2. `handlers.DaftarkanRute(mux, svc, stubPelaku)` mendaftarkan seluruh rute modul dengan satu awalan
   (URL rute ditetapkan spec modulnya; nama folder tidak menentukan URL).
3. Bila bermigrasi: pilih rentang nomor yang belum terpakai (catat di tabel bab 3), taruh berkas di
   `modul/nbfacin/migrations/` beserta `_down.sql`, dan tanam dengan `//go:embed migrations/*.sql`.
4. `modul/nbfacin/modul.go` (`package nbfacin`): `const Nama = "nbfacin"`, `Baru(...)`, dan metode `Nama()`,
   `DaftarkanRute(mux)`, `JalankanPekerja(ctx)` (tanpa pekerja: `inti.TanpaPekerja()`), serta
   `SumberMigrasi()` bila bermigrasi — pola keempat modul yang ada.
5. `modul/daftar.go`: satu baris di `Rakit`, dan satu baris di `SumberMigrasi` bila bermigrasi. Butuh
   modul lain? Tambah antarmuka di `inti/kontrak` dan sambungkan di sini.

**Frontend**

6. `frontend/src/modul/nb-facin/`: `pages/`, `components/`, `labels.ts`, `api.ts` (memakai
   `minta`/`mintaFormulir` dari `inti/klien`), `menu.ts` (`NAMA_NBFACIN = 'nbfacin'`,
   `HALAMAN_NBFACIN`, `MENU_NBFACIN`), `rute.tsx` (`RuteNbFacin`, menyimpan keadaan
   kasusnya sendiri).
7. `frontend/src/modul/daftar.ts`: satu baris di `MODUL_FRONTEND` dan satu anggota di union `Halaman`.
   Nama kelompok sidebar-nya sudah ada di `MODUL` (`inti/labels.ts`) bila modulnya salah satu folder
   korpus.
8. **Menu**: butir `menu.ts` baru TAMPIL hanya bila ada barisnya di `M_NAV_MENU` — kelompok modulnya
   sudah ada di isi awal (20 folder korpus, `DIMIGRASI = '0'`); tambahkan butirnya menurut bab 6.

**Yang akan memeriksa Anda** — jalankan gerbang bab 3; yang biasanya berbunyi:

| Uji | Menangkap |
| --- | --- |
| `inti/penjaga/impor_lintas_modul_test.go`, `frontend/src/inti/lapisan.guard.test.ts` | impor lintas modul |
| `cmd/api/rakit_test.go` | rute modul nonaktif harus 404; migrasi di disk = migrasi di pelari |
| `frontend/src/modul/daftar.modulAktif.test.ts` | nama frontend ≠ `const Nama` di `modul/*/modul.go` |
| `frontend/src/modul/daftar.sinkron.test.ts`, `frontend/src/Shell.test.ts` | sidebar ↔ palet, butir menu berbukti korpus |
| `frontend/src/modul/daftar.menuTabel.test.ts`, `inti/penjaga/menu_test.go` | isi `M_NAV_MENU` ↔ `menu.ts` dua arah, bentuk SQL menu (bab 6) |
| `inti/penjaga/*` | higiene migrasi, kata cadangan Oracle, alamat layanan, nama orang, nama tabel telanjang |

`Shell.test.ts` mengunci jumlah butir menu (hari ini lima). Menambah butir menuntut bukti XML
korpus dan menyunting angka itu dengan alasan — itu disengaja.

## 6. Menambah menu — satu baris di `M_NAV_MENU` + satu butir di `menu.ts`

Sejak 30-09-2026 sidebar dan palet Ctrl+K dirakit dari tabel **`M_NAV_MENU`** lewat `GET /api/menu`
(permintaan work owner: dasar akses menu per akun sesudah login ada). Tabelnya dua tingkat:

| Baris | `PARENT_ID` | `KODE` | Contoh |
| --- | --- | --- | --- |
| kelompok modul | kosong | nama modul backend (tabel nama modul) | `claimlife`, `nbfacin` |
| butir menu | ID kelompoknya | kunci halaman frontend (`menu.ts`) | `inbox`, `tco-tahun` |

`GROUPMENU` (`TREATY`, `FACULTATIVE`, `KLAIM`, `MASTER`) menjadi kepala bagian sidebar, dalam urutan itu.
Isi awal (migrasi 900) memuat 20 kelompok — satu per folder korpus — dan 5 butir; daftarnya di
`..\.scratch\inti\STRUKTUR-TABEL-INTI.md`.

**Frontend memotong pohon tabel dengan rute yang terdaftar** (`inti/lib/daftarMenu.ts` `susunMenu`):

- baris tabel **tanpa** rute frontend tidak tampil, dan dicatat di konsol peramban
  (`menu: 1 butir M_NAV_MENU tanpa rute frontend, tidak tampil: …`);
- butir `menu.ts` **tanpa** baris tabel tidak tampil — di sidebar maupun palet;
- kelompok `DIMIGRASI = '0'` tampil terlipat "belum dimigrasi"; kelompok yang sudah dimigrasi tetapi
  tanpa butir (modulnya tidak ada di `MODUL_AKTIF`) tidak tampil.

Jadi menambah satu menu = **dua sisi, satu deploy**:

1. **Tabel** — migrasi BARU di `inti/migrations/` dengan nomor bebas berikutnya di 900–949
   (mis. `901_menu_nbfacin.sql` + `901_menu_nbfacin_down.sql`). ⛔ **Jangan menyunting 900**:
   `T_MIGRASI` mencatat nama, jadi isi 900 yang diubah tidak pernah dijalankan ulang. ⛔ **Jangan dari
   folder migrasi modul**: isi menu hanya di `inti` (`TestMenuHanyaDiMigrasiInti`). Bentuknya PERSIS
   bentuk 900 — kedua penjaga membacanya dengan pola itu:

   ```sql
   -- 901 - menu NB FacIn: modulnya mendapat layar pertamanya, lalu butirnya.
   UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE
   WHERE KODE = 'nbfacin' AND PARENT_ID IS NULL
   /
   INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
   SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, k.ID, 'nbfacin-inbox', 'Inbox NB FacIn', k.GROUPMENU, k.MODUL, 1, k.DIMIGRASI
   FROM {skema}.M_NAV_MENU k
   WHERE k.KODE = 'nbfacin' AND k.PARENT_ID IS NULL
   AND NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU b WHERE b.KODE = 'nbfacin-inbox')
   /
   ```

   `UPDATE` **sebelum** `INSERT`: butir mewarisi `GROUPMENU`, `MODUL`, dan `DIMIGRASI` induknya saat
   disisipkan. Keduanya idempoten, jadi langkah yang gagal separuh jalan aman diulang pelari. Jalur
   mundurnya (`_down.sql`) menghapus butirnya (`DELETE … WHERE KODE = 'nbfacin-inbox'`) lalu
   mengembalikan `DIMIGRASI` ke `'0'`. `UPDATE DIMIGRASI` hanya untuk modul yang mendapat layar
   **pertamanya**; `TestIsiAwalMenuDimigrasiSamaDenganModulBackend` menuntut `DIMIGRASI = '1'` tepat
   untuk modul yang punya `modul/<nama>/modul.go`.

2. **Kode** — satu butir di `frontend/src/modul/<nama>/menu.ts` dengan `modul` = `KODE` baris tabel dan
   `label` = `LABEL`-nya **VERBATIM** (dan halamannya di `rute.tsx`).

3. **Penjaga** — `frontend/src/modul/daftar.menuTabel.test.ts` (dua arah: `KODE` ↔ `menu.ts`, `LABEL`,
   induk = modul pemilik) dan `inti/penjaga/menu_test.go` (bentuk SQL, idempoten, `CHECK GROUPMENU`,
   20 kelompok = folder korpus, daftar butir yang dikunci — sunting angka/daftarnya dengan alasan, seperti
   `Shell.test.ts`).

4. **`-migrate` dijalankan work owner** — sampai itu, backend baru pun tetap menjawab dari baris lama,
   dan butir baru tidak tampil.

Menyembunyikan satu menu: `STATUS_AKTIF = '0'` pada barisnya (atau pada kelompoknya — seluruh butirnya
ikut hilang); pembaca menu hanya membaca baris `'1'`. Itu **perubahan data di Oracle** — tulis ke DB,
jadi dilakukan work owner/DBA dengan persetujuan, bukan executor. Penjaga membaca migrasi, jadi tidak
melihat perubahan data semacam itu; bila menu itu memang dibuang untuk seterusnya, tuliskan sebagai
migrasi `inti` berikutnya dan sunting kedua penjaga.

**Di luar lingkup hari ini** (dicatat, tidak dibangun): tabel akses per akun (mis. `M_NAV_MENU_AKSES`:
akun atau peran → `MENU_ID`) dan login. Titik sambungnya sudah ada: `inti/menu`
`SaringMenuUntukPelaku(pelaku, menu)` — hari ini meneruskan semua — dipanggil `GET /api/menu` untuk
setiap permintaan.
