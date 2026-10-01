# Panduan cabang git — cabang mengatur waktu, CODEOWNERS mengatur wilayah

Untuk siapa: empat pengembang fullstack yang masing-masing memegang modul berbeda, dan tim inti yang
merawat kode bersama. Panduan ini menjawab pertanyaan yang selalu muncul di hari pertama — "cabang
`main` isinya folder apa saja, dan `dev` apa saja?" — lalu menetapkan nama cabang, arah merge, dan
aturan yang membuat empat orang mendorong ke repositori yang sama tanpa konflik.

Pembagian folder dan pemiliknya: `PANDUAN-TIM-PER-MODUL.md` (bersebelahan dengan berkas ini).
Deploy sebagian modul: `..\..\APP_RNM\PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`.

## 1. Salah paham yang harus dibereskan lebih dulu

**Cabang bukan sebagian folder. Cabang adalah seluruh repositori pada waktu yang berbeda.**

`main` bukan "hanya inti". `dev` bukan "modul-modul saja". Keduanya memuat **seluruh pohon repositori** —
`APP_RNM/` dengan dua puluh modulnya, `inti/`, `docs/`, `discovery/`, `dastin/`, `jefri/`, dan dokumen
akar. Yang membedakan bukan *folder mana*, melainkan *versi kapan*.

Keadaan repositori pada 1 Oktober 2026 membuktikannya:

| | `main` | `dev` |
| --- | --- | --- |
| Jumlah berkas | 2.129 | 2.137 |
| Folder tingkat-1 | `.github`, `APP_RNM`, `dastin`, `discovery`, `docs`, `jefri`, `keluaran-docx`, `scripts` | **sama persis** |
| `APP_RNM/modul/` | dua puluh modul | dua puluh modul |
| `APP_RNM/inti/` | ada | ada |

Selisih delapan berkas itu bukan folder yang ditinggal — itu berkas baru yang lahir di `dev` dan belum
naik ke `main`. Selisih commit-nya tujuh belas, dan arahnya selalu satu: `dev` di depan, `main` di
belakang, tidak pernah sebaliknya.

Perumpamaannya: `main` dan `dev` adalah dua naskah **lengkap** dari buku yang sama — bukan "bab 1-5 di
satu map, bab 6-10 di map lain". `dev` naskah draf terbaru, `main` naskah yang sudah dicetak.

Akibat praktisnya: ketika pemegang `nbtreatyin` menjalankan `git checkout modul/nbtreatyin`, folder di
laptopnya tetap berisi **semua** modul dan seluruh `inti/`. Memang harus begitu — kodenya meng-import
`inti/backend/kontrak` dan `inti/frontend/components`; tanpa itu aplikasi tidak bisa di-build. Yang
membatasi dia bukan cabang, melainkan `..\..\.github\CODEOWNERS`.

> **Cabang mengatur waktu — kapan pekerjaan naik ke rilis.**
> **CODEOWNERS mengatur wilayah — siapa boleh mengubah apa.**
> Dua hal berbeda, dan keduanya perlu.

## 2. Isi setiap cabang

### `main` — potret yang dinyatakan stabil

Seluruh pohon repositori pada titik terakhir yang ujinya hijau dan sudah disetujui. Patokannya: kalau
hari ini ada demo ke pengguna atau pemasangan ke server, yang diambil `main`. Tidak pernah di-push
langsung; satu-satunya jalan masuk adalah pull request dari `dev`.

### `dev` — potret yang sama, versi lebih baru

Juga seluruh pohon repositori — memuat pekerjaan yang sudah selesai dan sudah ditinjau, tapi belum
dinyatakan rilis. Inilah cabang yang ditarik tim setiap pagi, dan dasar bagi semua cabang modul.
`origin/HEAD` menunjuk ke sini, jadi `dev` adalah cabang bawaan repositori.

Tujuh belas commit yang sekarang ada di `dev` dan belum di `main` menyentuh `inti/`,
`modul/premiumlistlife/`, `modul/treatycontractout/`, `uji/`, `docs/bersama/`, sampai dokumen akar —
tujuh puluh delapan berkas di semua lapisan. Itu menegaskan sekali lagi: `dev` bukan soal folder
tertentu, melainkan kemajuan waktu.

### `modul/<nama>` — cabang tetap milik satu orang untuk satu modul

Dicabang dari `dev`, umurnya sepanjang migrasi modul itu. Contoh: `modul/nbtreatyin`,
`modul/claimprop`. Isinya tetap seluruh repositori; yang berubah hanya folder modul pemiliknya.

### `modul/<nama>/<topik>` — cabang tugas, umur pendek

Dicabang dari `modul/<nama>`, satu per tiket, di-merge balik lalu dihapus. Contoh:
`modul/nbtreatyin/spec-penyimpanan`, `modul/claimprop/tiket-07-komite`. Satu modul butuh belasan
tiket — `claimlife` menghabiskan dua belas giliran — jadi tingkat ini yang menjaga pull request tetap
kecil dan bisa ditinjau.

### `inti/<topik>` — perubahan berkas bersama

Dicabang dari `dev`, langsung ke `dev`, sekecil mungkin, dan **tidak pernah dibonceng di pull request
modul**. Semua yang disebut bab 4 masuk lewat sini.

### `perbaikan/<topik>` — hotfix

Dicabang dari `main` ketika yang sudah dipasang rusak. Di-merge ke `main` **dan** ke `dev` — kalau lupa
yang kedua, perbaikannya hilang pada rilis berikutnya.

## 3. Arah merge — satu jalur, tidak pernah berbalik

```
modul/<nama>/<topik>  ->  modul/<nama>  ->  dev  ->  main
inti/<topik>          ->  dev
perbaikan/<topik>     ->  main  dan  dev
```

Cabang modul **tidak pernah** langsung ke `main`.

| Pola | Dicabang dari | Di-merge ke | Umur |
| --- | --- | --- | --- |
| `main` | — | — | tetap |
| `dev` | `main` | `main` | tetap |
| `modul/<nama>` | `dev` | `dev` | sepanjang migrasi modul |
| `modul/<nama>/<topik>` | `modul/<nama>` | `modul/<nama>` | satu sampai tiga hari |
| `inti/<topik>` | `dev` | `dev` | satu hari |
| `perbaikan/<topik>` | `main` | `main` + `dev` | sependek mungkin |

## 4. Berkas bersama — satu-satunya sumber konflik yang tersisa

Isolasi satu folder per modul sudah menyelesaikan hampir seluruh masalah. Yang tersisa adalah berkas
yang dipakai semua modul. **Jangan pernah menyentuhnya di cabang modul:**

`go.mod`, `go.sum`, `package.json`, `package-lock.json`, `vite.config.ts`, `tsconfig.json`, `Makefile`,
`APP_RNM/inti/**` (kecuali `daftar/modul_<nama>_gen.go` milik Anda), `APP_RNM/frontend/**`,
`APP_RNM/cmd/**`, `APP_RNM/uji/**`, `.github/**`, `docs/bersama/**`.

Perlu mengubahnya? Cabang `inti/<topik>` -> pull request kecil ke `dev` -> di-merge lebih dulu -> semua
orang menarik `dev`. Menumpang di pull request modul membuat empat orang menunggu satu tinjauan.

**`package-lock.json` jangan pernah di-merge manual.** Bila bentrok:

```bash
git checkout --theirs APP_RNM/package-lock.json
npm install --prefix APP_RNM
git add APP_RNM/package-lock.json
```

Lockfile yang berlaku hanya satu: `APP_RNM/package-lock.json`. Aplikasi hidup di `APP_RNM/`, jadi
`npm install` selalu dijalankan dari sana (atau dengan `--prefix APP_RNM`). Lockfile yang muncul di akar
repositori adalah hasil `npm install` di folder salah — hapus, jangan di-commit.

**Nomor ADR adalah sumber konflik tersembunyi.** Empat orang menulis ADR berikutnya serentak
menghasilkan empat berkas `0044-*.md` dan empat suntingan `adr/00-INDEKS.md`. Cadangkan blok nomor per
orang, atau ambil nomornya lewat pull request satu baris ke `dev` sebelum menulis isinya.

**Rentang migrasi dan slot menu** sudah dipartisi di `MODUL.md` tiap modul dan ditegakkan
`inti/backend/penjaga`. Jangan keluar dari rentang Anda; mengubah rentang adalah pull request ke tim
inti.

## 5. Kebiasaan harian

**Tarik `dev` setiap pagi.** Satu kebiasaan ini yang paling menurunkan konflik.

```bash
git checkout modul/nbtreatyin
git fetch origin
git merge origin/dev
```

Gunakan `merge`, bukan `rebase`, untuk cabang yang sudah dipush — menulis ulang riwayat yang sudah
dilihat orang lain memaksa mereka memperbaiki salinannya. Untuk cabang topik yang **belum** dipush,
`git rebase origin/dev` lebih rapi.

**Pull request kecil dan sering** — sekali per tiket, bukan sekali per modul. Cabang modul yang hidup
tiga minggu tanpa merge akan bentrok hebat di `go.mod`, `inti/backend/daftar/daftar.go`, dan
`APP_RNM/frontend/daftar.ts`.

**Pindah modul tanpa stash** memakai worktree (`.worktrees/` sudah diabaikan `.gitignore`):

```bash
git worktree add .worktrees/nbtreatyin modul/nbtreatyin
```

## 6. Pesan commit

Riwayat adalah satu-satunya cara menelusuri "perubahan ini kenapa" setelah pemiliknya lupa. Formatnya:

```
<lingkup>: <apa yang berubah>
```

`<lingkup>` adalah nama modul, `inti`, atau `docs`. Contoh yang sudah ada di repositori ini dan layak
ditiru:

```
inti: sidebar satu tombol per modul, dikelompokkan GROUPMENU
treaty-contract-out: ReinsType anak Treaty Limit - porsi + induknya, autocomplete ikut XML
docs: brief implementasi Master Contract Retro Life - ralat PK/FK dari katalog DEV, nol DDL
```

Pesan seperti `KOMIT`, `HEHE`, atau `LALALLALLA` — ada empat di riwayat `dev` dari masa repositori ini
dipegang sendiri — tidak lagi diterima sejak tim berisi empat orang.

## 7. Proteksi cabang yang harus dinyalakan di GitHub

| | `main` | `dev` |
| --- | --- | --- |
| Push langsung | ditolak | ditolak |
| Pull request wajib | ya | ya |
| Persetujuan | 1, dari orang lain | 1, dari orang lain |
| Tinjauan CODEOWNERS | wajib | wajib |
| Gerbang uji (`go test`, `vitest`) hijau | wajib | wajib |
| Riwayat linear | ya | — |
| Hapus cabang sesudah merge | ya | ya |

"Require approval from someone other than the last pusher" harus menyala: ADR 0040 melarang menyetujui
pekerjaan sendiri.

Agar tinjauan CODEOWNERS berfungsi, penanda di `..\..\.github\CODEOWNERS` harus diganti akun GitHub
sebenarnya — dua puluh `@PEMILIK-*` dan `@TIM-INTI`. Selama masih penanda, aturannya tertulis tapi tidak
menegakkan apa pun.

## 8. Pembagian empat orang

Dasarnya: **modul yang berbagi tabel atau kontrak harus satu pemilik.** ADR 0041 menyatakan endorsemen
berbagi tabel dengan new business, dan tiap `komiteclaim*` memakai kontrak `claim*` pasangannya.
Memisahkan pasangan itu ke dua orang berarti menjadwalkan konflik.

| Orang | Lingkup | Modul | Rentang migrasi | Slot menu |
| --- | --- | --- | --- | --- |
| A | Fakultatif | `nbfacin`, `rnwfacin`, `endorsmentfacin`, `claimfacin`, `komiteclaimfacin` | 180-299, 520-559, 640-679 | 962-967, 978-979, 984-985 |
| B | Treaty In | `nbtreatyin`, `edmtreatyin`, `treatyin`, `treatyinadjustment` | 320-479 | 968-975 |
| C | Master & Life | `mastercontractretrolife`, `masterproductnamelife`, `endorsementlife`, rawat `premiumlistlife` dan `treatycontractout` | 100-179, 480-519, (050-099, 300-319) | 958-961, 976-977, (954-957) |
| D | Klaim Prop/Non-Prop | `claimprop`, `komiteclaimprop`, `claimnonprop`, `komiteclaimnonprop`, rawat `claimlife` dan `komiteclaimlife` | 560-639, 680-759, (001-049) | 980-983, 986-989, (950-953) |

Rentang migrasi keempatnya **tidak bersinggungan** — itulah yang membuat `backend/migrations/` tidak
pernah bentrok. A memegang lima modul, tetapi tiga di antaranya berbagi tabel, jadi bebannya setara B.
