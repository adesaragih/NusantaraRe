# Panduan deploy dan git per modul — APP_RNM bentuk B

Ditulis 30 September 2026, sesudah refactor bentuk B (`..\PROMPT-REFACTOR-BENTUK-B-MODUL.md`,
paket 1–8). Untuk siapa: pengembang yang memegang satu modul, dan siapa pun yang menyalakan
aplikasi dengan sebagian modul saja. Diperbarui 30 September 2026: menu dari tabel `M_NAV_MENU`
(`..\PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`) — bab 2, 3, 5, dan bab 6 baru. Diperbarui lagi 30 September
2026: **struktur tim satu folder per modul** (`..\PROMPT-STRUKTUR-TIM-SATU-FOLDER-PER-MODUL.md`) — bab 1,
3, 4, 5, dan 6. Alur kerja tim dan cara memulai modul kerangka: `..\docs\bersama\PANDUAN-TIM-PER-MODUL.md`.

## 1. Bentuknya dalam satu layar

Satu modul = **satu folder** `modul/<nama>/` berisi backend, frontend, dokumen, dan `MODUL.md` — dipegang
satu orang fullstack (`..\.github\CODEOWNERS`).

```
APP_RNM/
  go.mod go.sum package.json package-lock.json vite.config.ts tsconfig.json   milik tim inti
  cmd/api/              memasang modul dari daftar; MODUL_AKTIF; /healthz, /api/modul-aktif, /api/menu
  inti/
    backend/            kode Go BERSAMA - tidak pernah mengimpor modul (kecuali daftar/)
      kontrak/          antarmuka lintas modul, TANPA implementasi (PembacaPolis, KlaimKomite)
      perakit.go        menyambung kontrak menurut Pendaftaran() setiap modul
      daftar/           daftar modul BANGKITAN (go generate): modul_<nama>_gen.go per modul
      menu/             GET /api/menu: pembaca M_NAV_MENU, pohon GROUPMENU, SaringMenuUntukPelaku
      migrations/       tabel lintas modul 900-949: M_NAV_MENU (900)
      penjaga/          uji penjaga SELURUH aplikasi - membaca modul/* dan MODUL.md, nol nama modul
    frontend/           kerangka React bersama: Shell, ui/dasar, klien.ts, lib/, hooks/, store/, labels.ts
  modul/
    <nama>/             20 folder, satu per folder korpus (4 dimigrasi, 16 kerangka)
      MODUL.md          pemilik, rentang migrasi, slot menu, prefix rute, kontrak, pernyataan penjaga
      backend/          models/ repository/ services/ handlers/ migrations/ modul.go
      frontend/         pages/ components/ labels.ts api.ts menu.ts rute.tsx
      docs/             spec, tiket, grilling (dulu ..\.scratch\<nama-panjang>\)
    _templat/           templat folder modul
  uji/                  penunjang uji netral: skemauji (skema Oracle tiruan), lintasmodul
  frontend/             perakit: index.html main.tsx App.tsx Beranda.tsx daftar.ts katalogKorpus.ts
..\docs\bersama\       ADR, CONTEXT, STRUKTUR-TABEL-INTI, PANDUAN-TIM-PER-MODUL
..\.github\CODEOWNERS   pemilik setiap folder
```

**Tabel nama modul** — satu-satunya sumber, tanpa singkatan *(keputusan work owner 30-09-2026, `PROMPT-REFACTOR-NAMA-MODUL.md`)*:

| Modul korpus | Folder `APP_RNM/modul/…` = `MODUL_AKTIF` | Dokumen lama `.scratch/…` |
| --- | --- | --- |
| Claim Life | `claimlife` | `claim-life` |
| PremiumList Life | `premiumlistlife` *(dulu `premiumlist`)* | `premiumlist-life` |
| Komite Claim Life | `komiteclaimlife` *(dulu `komite`)* | `komite-claim-life` |
| Treaty Contract Out | `treatycontractout` *(dulu `treaty`)* | `treaty-contract-out` |
| *(kerangka, mis. NB Treaty In)* | `nbtreatyin` | `nb-treaty-in` |

Aturan: **nama modul = nama folder korpus tanpa spasi, huruf kecil**. Satu nama untuk folder,
`const Nama` Go, `MODUL_AKTIF`, dan `KODE` kelompok `M_NAV_MENU`. Paket Go `backend/modul.go` bernama
`backend` di setiap modul; daftar bangkitan mengimpornya dengan alias nama modul.

**Aturan impor** — ditegakkan uji, bukan kesepakatan:

| Dari | Boleh mengimpor | Penjaga |
| --- | --- | --- |
| `modul/X/backend/...` (Go) | `inti/backend/...`, `modul/X/...`; berkas uji juga `uji/skemauji` ¹. **Tidak pernah** `inti/backend/daftar` | `inti/backend/penjaga/impor_lintas_modul_test.go` |
| `inti/backend/...` (Go) | `inti/...` saja — kecuali daftar bangkitan `inti/backend/daftar`, yang mengimpor PAKET AKAR `modul/<nama>/backend` | idem |
| `cmd/...` (Go) | `inti/...` saja — modul dipasang lewat `inti/backend/daftar` | idem |
| Letak kode Go modul | hanya `modul/<nama>/backend/` | idem |
| `modul/X/frontend/**` | `inti/frontend/**`, `modul/X/frontend/**` | `inti/frontend/lapisan.guard.test.ts` |
| `inti/frontend/**` | `inti/frontend/**` saja | idem |
| `frontend/**` (perakit), `uji/` | apa pun — merekalah tempat yang mengenal semua modul | — |

¹ `uji/skemauji` SENGAJA mengenal semua modul: ia membangun skema uji **utuh** (migrasi semua modul,
fixture Claim Life dan PremiumList). Karena itu ia satu-satunya jalur lintas modul yang disahkan —
hanya untuk berkas uji, dan hanya lewat paket itu.

Modul yang membutuhkan modul lain **tidak mengimpornya**: `Pendaftaran()` di `backend/modul.go`-nya
menyatakan `Membutuhkan` (antarmuka `inti/backend/kontrak`) dan membacanya dengan `inti.Ambil`;
penyedianya menyatakan `Menyediakan` dan menyerahkannya dengan `inti.Sediakan`. Perakit
(`inti/backend/perakit.go`) membangun penyedia lebih dulu, dan **menolak menyala** — dengan kalimat yang
menyebut kontrak dan modulnya — bila kontrak dibutuhkan tanpa penyedia, disediakan dua modul, melingkar,
atau dinyatakan tetapi tidak diserahkan. Dua sambungan yang ada hari ini:

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
  dari `GET /api/modul-aktif` (`{"modul":["claimlife","komiteclaimlife"]}`, urutan nama modul). **Tidak ada env Vite** untuk
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

Seluruh aplikasi satu repositori; `git pull` selalu menarik semuanya. Satu modul = **satu jalur**:
yang Anda commit, lihat, dan salin adalah `APP_RNM/modul/<nama>/`.

```powershell
# Commit HANYA folder modul Anda.
git add -A -- APP_RNM/modul/claimlife
git commit -m "claimlife: ..." -- APP_RNM/modul/claimlife

# Riwayat dan beda satu modul saja
git log --oneline -- APP_RNM/modul/claimlife
git diff HEAD~1 -- APP_RNM/modul/komiteclaimlife

# Menarik perubahan: pull utuh, lalu lihat apa yang berubah di luar folder modul
git pull
git diff --stat ORIG_HEAD -- APP_RNM/inti APP_RNM/frontend APP_RNM/cmd APP_RNM/uji
```

Memulai sebuah modul kerangka menambah SATU berkas di luar foldernya: daftar bangkitan
`APP_RNM/inti/backend/daftar/modul_<nama>_gen.go` (`go generate ./inti/backend/daftar`), yang di
`CODEOWNERS` juga milik pemilik modul itu. Satu berkas per modul, jadi dua cabang yang memulai dua
modul tidak pernah berkonflik di sana.

⚠️ Sparse-checkout satu modul **tidak dapat membangun backend**: daftar bangkitan mengimpor setiap
modul terdaftar. Frontend dapat (daftar `import.meta.glob` hanya melihat folder yang ada). Pakai untuk
membaca dan menyunting; bangun dan uji selalu dari salinan utuh.

**Aturan yang berlaku untuk setiap commit** (gerbang lengkap, sama dengan CI):

```powershell
go build ./... ; go build -tags db ./... ; go vet ./... ; go vet -tags db ./... ; gofmt -l cmd inti modul uji
go test ./... ; go test -tags db ./...          # tanpa ORACLE_DSN: test db SKIP dengan pesan
npx tsc --noEmit ; npx vitest run ; npx vite build   # dari APP_RNM/ - package.json di sini sejak 30-09-2026
```

- Menyentuh berkas milik tim inti (`inti/`, `cmd/`, `frontend/`, `uji/`, konfigurasi npm/Go/Vite/TS)
  berarti menyentuh **semua modul**: jalankan seluruh uji, dan minta tinjauan tim inti (`CODEOWNERS`).
- Menambah kebutuhan lintas modul = antarmuka baru di `inti/backend/kontrak` (pull request tim inti,
  tanpa implementasi) dan pernyataan `Menyediakan` / `Membutuhkan` di `Pendaftaran()` kedua modul —
  tidak pernah impor langsung.
- **Nama berkas migrasi tidak pernah diubah**: `T_MIGRASI` mencatat nama, dan nama baru membuat
  `-migrate` menjalankannya ulang. **Nomor selalu tiga digit** (pelari mengurutkan nama sebagai teks:
  `1000_` akan berjalan sebelum `101_`). Rentang setiap modul dinyatakan `MODUL.md`-nya dan dijaga
  `inti/backend/penjaga/rentang_test.go` — urut hulu ke hilir, supaya migrasi modul hilir yang merujuk
  tabel modul hulu selalu berjalan sesudahnya:

| GROUPMENU | Modul | Rentang migrasi | Slot menu | Terpakai |
| --- | --- | --- | --- | --- |
| KLAIM | `claimlife` | 001–029 | 950–951 | 001–022 |
| KLAIM | `komiteclaimlife` | 030–049 | 952–953 | 030 |
| TREATY | `premiumlistlife` | 050–099 *(dulu tertulis 050–079)* | 954–955 | 050–058 |
| MASTER | `treatycontractout` | 300–319 | 956–957 | kosong — tco4: nol tabel baru (`TestTCONolTabelBaru`) |
| MASTER | `mastercontractretrolife` | 100–139 | 958–959 | kerangka |
| MASTER | `masterproductnamelife` | 140–179 | 960–961 | kerangka |
| FACULTATIVE | `nbfacin` | 180–219 | 962–963 | kerangka |
| FACULTATIVE | `rnwfacin` | 220–259 | 964–965 | kerangka |
| FACULTATIVE | `endorsmentfacin` | 260–299 | 966–967 | kerangka |
| TREATY | `nbtreatyin` | 320–359 | 968–969 | kerangka |
| TREATY | `edmtreatyin` | 360–399 | 970–971 | kerangka |
| TREATY | `treatyin` | 400–439 | 972–973 | kerangka |
| TREATY | `treatyinadjustment` | 440–479 | 974–975 | kerangka |
| TREATY | `endorsementlife` | 480–519 | 976–977 | kerangka |
| KLAIM | `claimfacin` | 520–559 | 978–979 | kerangka |
| KLAIM | `claimprop` | 560–599 | 980–981 | kerangka |
| KLAIM | `claimnonprop` | 600–639 | 982–983 | kerangka |
| KLAIM | `komiteclaimfacin` | 640–679 | 984–985 | kerangka |
| KLAIM | `komiteclaimprop` | 680–719 | 986–987 | kerangka |
| KLAIM | `komiteclaimnonprop` | 720–759 | 988–989 | kerangka |
| — | `inti` *(tabel lintas modul, `inti/backend/migrations/`)* | 900–949 | — | 900 (`M_NAV_MENU` + isi awal) |
| — | cadangan, dibagi tim inti lewat pull request | 760–899 dan 990–999 | | |

## 4. Pemilik folder

Pemilik ditetapkan di **`..\.github\CODEOWNERS`** — nama akun di sana PENANDA (`@PEMILIK-CLAIMLIFE`,
`@TIM-INTI`, …) yang diisi work owner. Aturannya:

| Jalur | Pemilik |
| --- | --- |
| `APP_RNM/modul/<nama>/` (backend, frontend, docs) | pemilik modul itu |
| `APP_RNM/modul/<nama>/MODUL.md`, `APP_RNM/modul/<nama>/backend/migrations/9*` (slot menu) | ditulis pemilik modul, disetujui tim inti — rentang, slot, `Status`, pernyataan penjaga, butir menu |
| `APP_RNM/inti/backend/daftar/modul_<nama>_gen.go` (bangkitan) | pemilik modul itu |
| `APP_RNM/inti/`, `cmd/`, `frontend/`, `uji/`, `pkg/`, `modul/_templat/` | tim inti |
| `APP_RNM/go.mod`, `go.sum`, `package.json`, `package-lock.json`, `vite.config.ts`, `tsconfig.json`, `Makefile` | tim inti — pustaka baru lewat pull request |
| `docs/bersama/`, `docs/agents/`, `.github/`, dokumen akar | tim inti |

Brief acuan setiap modul tercatat di `MODUL.md`-nya.

## 5. Memulai modul kerangka, atau menambah modul

Keenam belas folder korpus yang belum dimigrasi SUDAH punya folder kerangka (`MODUL.md` dengan rentang
migrasi dan slot menunya, `docs/` bila ada). Memulai satu = menambah `backend/` dan `frontend/` di
foldernya sendiri — langkah lengkapnya, beserta kerangka `modul.go`, `menu.ts`, dan `rute.tsx`:
`..\docs\bersama\PANDUAN-TIM-PER-MODUL.md` bab 4. Modul di luar dua puluh folder korpus bukan bagian
migrasi ini — keputusan work owner, lewat `modul/_templat/` dan pull request tim inti (bab 5 panduan itu).

**Yang akan memeriksa Anda** — jalankan gerbang bab 3; yang biasanya berbunyi:

| Uji | Menangkap |
| --- | --- |
| `inti/backend/penjaga/impor_lintas_modul_test.go`, `inti/frontend/lapisan.guard.test.ts` | impor lintas modul, kode Go di luar `backend/` |
| `inti/backend/daftar/bangkit/main_test.go` | daftar modul basi — jalankan `go generate ./inti/backend/daftar` |
| `inti/backend/daftar/daftar_test.go`, `inti/backend/perakit_test.go` | kontrak tanpa penyedia, sambungan kontrak |
| `inti/backend/penjaga/rentang_test.go` | migrasi di luar rentang atau slot, nomor bukan tiga digit, rentang bertumpuk, menu di luar slot |
| `cmd/api/rakit_test.go` | rute modul nonaktif harus 404; migrasi di disk = migrasi di pelari |
| `frontend/daftar.modulAktif.test.ts`, `frontend/daftar.rakit.test.ts` | nama frontend ≠ nama folder / `const Nama`; `menu.ts` tanpa `rute.tsx` |
| `frontend/daftar.sinkron.test.ts`, `frontend/Shell.test.ts` | sidebar ↔ palet, butir menu berbukti korpus |
| `frontend/daftar.menuTabel.test.ts`, `inti/backend/penjaga/menu_test.go` | isi `M_NAV_MENU` ↔ `menu.ts` dua arah, bentuk SQL menu (bab 6), kelompok modul = `LABEL` tabel |
| `inti/backend/penjaga/*` | higiene migrasi, dokumen STRUKTUR, kata cadangan Oracle, alamat layanan, nama orang, nama tabel telanjang |

Butir menu baru menuntut bukti XML korpus, dan setiap butir menu baru lewat tinjauan tim inti — itu
disengaja. Sejak 30-09-2026 tinjauannya lewat `CODEOWNERS` atas berkas slot menu `9*` di folder
migrasi modul, bukan lewat angka kunci di `Shell.test.ts` (dulu lima): angka di berkas bersama
membuat dua modul yang menambah butir bersamaan berkonflik di baris yang sama. Kelompok yang belum
dimigrasi dibaca dari `Status` di `MODUL.md` setiap modul.

## 6. Menambah menu — satu `UPDATE DIMIGRASI` di slot modul + `HALAMAN_AWAL` di `menu.ts`

Sejak 30-09-2026 sidebar dan palet Ctrl+K dirakit dari tabel **`M_NAV_MENU`** lewat `GET /api/menu`
(permintaan work owner: dasar akses menu per akun sesudah login ada). Sejak migrasi **901**
(keputusan work owner 30-09-2026: *"menu jangan ada model seperti child … 1 modul 1 menu"*) tabelnya
**datar**: 20 baris, satu per folder modul korpus; `KODE` = `MODUL` = nama modul (tabel nama modul);
`URUTAN` = urutan di dalam `GROUPMENU`. Tidak ada `PARENT_ID`, tidak ada butir anak. Daftarnya di
`..\docs\bersama\STRUKTUR-TABEL-INTI.md`.

Sidebar: kepala `GROUPMENU` (`TREATY`, `FACULTATIVE`, `KLAIM`, `MASTER`, dalam urutan itu), di bawahnya
**satu tombol per modul** berlabel `LABEL` (nama folder korpus VERBATIM). Klik tombol membuka **halaman
awal** modul (`HALAMAN_AWAL_<X>`); halaman lain modul itu dibuka dari dalamnya — Register dari tombol
Register di Inbox Claim Life, Outstanding dan Detail dari baris kasus. Tidak ada kelompok yang dilipat.

`GET /api/menu` — satu tingkat di bawah golongan:

```json
{"golongan":[
  {"kode":"TREATY","modul":[
    {"kode":"nbtreatyin","label":"NB Treaty In","modul":"nbtreatyin","urutan":1,"dimigrasi":false},
    {"kode":"premiumlistlife","label":"PremiumList Life","modul":"premiumlistlife","urutan":5,"dimigrasi":true}]},
  {"kode":"KLAIM","modul":[ … ]}]}
```

**Frontend memotong menu tabel dengan modul frontend yang terdaftar** (`inti/frontend/lib/daftarMenu.ts`
`susunMenu`):

- baris `DIMIGRASI = '1'` **tanpa** modul frontend tidak tampil, dan dicatat di konsol peramban
  (`menu: 1 modul M_NAV_MENU DIMIGRASI='1' tanpa modul frontend, tidak tampil: …`);
- modul frontend **tanpa** baris tabel tidak tampil — di sidebar maupun palet;
- baris `DIMIGRASI = '0'` tampil sebagai tombol **nonaktif** "belum dimigrasi" (`aria-disabled`);
- modul yang sudah dimigrasi tetapi tidak ada di `MODUL_AKTIF` tidak dikirim backend, jadi tidak tampil.

Jadi "menambah menu" = modul mendapat layar pertamanya — **dua sisi, satu deploy, satu folder**:

1. **Tabel** — migrasi BARU di folder migrasi modul Anda sendiri,
   `modul/<nama>/backend/migrations/`, dengan nomor di **slot menu** modul itu (`MODUL.md`, mis.
   `nbfacin` 962–963): `962_menu_nbfacin.sql` + `962_menu_nbfacin_down.sql`. Isinya SATU pernyataan —
   menyalakan `DIMIGRASI` baris modul Anda:

   ```sql
   -- 962 - menu NB FacIn: modul ini mendapat layar pertamanya.
   UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE
   WHERE KODE = 'nbfacin'
   /
   ```

   Jalur mundurnya sama dengan `'0'`. ⛔ **Nol `INSERT`**: baris modul Anda sudah ada sejak 900, dan
   `INSERT INTO M_NAV_MENU` di slot menu MERAH (`TestSlotMenuHanyaMenyentuhMenuModulnya`). ⛔ Bentuk
   lama `… AND PARENT_ID IS NULL` mati di ORA-00904 sesudah 901. ⛔ **Jangan menyunting 900 maupun
   901**: `T_MIGRASI` mencatat nama, jadi isinya yang diubah tidak pernah dijalankan ulang. ⛔ **Jangan di
   `inti/backend/migrations/`**: 902–949 milik inti (`TestMenuHanyaDi900DanSlotMenuModulnya`). Urutan nama
   berkas sebagai teks menjamin slot 95x berjalan **sesudah** 900 dan 901 (`TestSlotMenuBerjalanSesudah900`).
   `TestMenuDimigrasiSamaDenganModulBackend` menuntut `DIMIGRASI = '1'` — sesudah semua slot — tepat untuk
   modul yang punya `modul/<nama>/backend/modul.go`.

2. **Kode** — `modul/<nama>/frontend/menu.ts`: `HALAMAN_AWAL_<X>` (salah satu `HALAMAN_<X>`) sebagai
   `PENDAFTARAN_MENU.halamanAwal`, dan `kelompok` = nama folder korpus = `LABEL` barisnya **VERBATIM**
   (halamannya di `rute.tsx`). Label tombol datang dari tabel, bukan dari modul.

3. **Penjaga** — `frontend/daftar.menuTabel.test.ts` (dua arah atas hasil bersih 900 + 901 + slot setiap
   modul ↔ `MODUL_FRONTEND`: baris dimigrasi ↔ modul dengan `HALAMAN_AWAL`, `LABEL` = nama modul; uji
   gigit kedua arah), `inti/backend/penjaga/menu_test.go` (skema tiruan 900 + 901 + slot: 20 baris = folder
   korpus, nol butir, nol `PARENT_ID`, `CHECK GROUPMENU`, blok 901 berpelindung katalog), dan
   `rentang_test.go` (slot hanya `UPDATE DIMIGRASI` modulnya). Berkas slot ditinjau tim inti lewat
   `CODEOWNERS` (bab 5).

4. **`-migrate` dijalankan work owner** — termasuk **901**. Backend baru sudah benar **sebelum** 901
   (pembacanya menyaring `KODE = MODUL` tanpa menyebut `PARENT_ID`, jadi lima butir lama tersaring) dan
   **sesudahnya**.

Menyembunyikan satu menu: `STATUS_AKTIF = '0'` pada baris modulnya; pembaca menu hanya membaca baris
`'1'`. Itu **perubahan data di Oracle** — tulis ke DB,
jadi dilakukan work owner/DBA dengan persetujuan, bukan executor. Penjaga membaca migrasi, jadi tidak
melihat perubahan data semacam itu; bila menu itu memang dibuang untuk seterusnya, tuliskan sebagai
migrasi di slot menu modul pemiliknya.

**Di luar lingkup hari ini** (dicatat, tidak dibangun): tabel akses per akun (mis. `M_NAV_MENU_AKSES`:
akun atau peran → `MENU_ID`) dan login. Titik sambungnya sudah ada: `inti/backend/menu`
`SaringMenuUntukPelaku(pelaku, menu)` — hari ini meneruskan semua — dipanggil `GET /api/menu` untuk
setiap permintaan.
