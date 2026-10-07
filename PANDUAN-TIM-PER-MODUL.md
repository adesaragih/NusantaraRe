# Panduan tim per modul — satu folder, satu pemilik

Untuk siapa: pengembang **fullstack** yang memegang satu modul, dan tim inti yang merawat kode bersama.
Keputusan work owner 30 September 2026 (`..\..\PROMPT-STRUKTUR-TIM-SATU-FOLDER-PER-MODUL.md`): repositori
dibagi ke beberapa orang, **satu orang mengerjakan satu modul**, dan saat di-push balik **tidak perlu
merge manual dan tidak ada konflik**. Panduan ini menjelaskan bentuk yang membuat itu mungkin dan cara
bekerja di dalamnya. Menjalankan aplikasi: `APP_RNM\PANDUAN-MENJALANKAN.txt`; deploy sebagian modul dan
menu: `APP_RNM\PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`; isi cabang `main` dan `dev`, nama cabang, arah
merge, dan pembagian empat orang: `PANDUAN-CABANG-GIT.md`.

## 1. Bentuknya

Semua milik satu modul tinggal di satu folder:

```
APP_RNM/modul/<nama>/
  MODUL.md      pemilik, rentang migrasi, slot menu, prefix rute, kontrak, pernyataan untuk penjaga
  backend/      models/ repository/ services/ handlers/ migrations/ modul.go
  frontend/     pages/ components/ labels.ts api.ts menu.ts rute.tsx *.test.ts
  docs/         spec, tiket (issues/), grilling, catatan
```

Dua puluh folder berdiri sejak 30-09-2026 — satu per folder korpus. Empat sudah dimigrasi (`claimlife`,
`premiumlistlife`, `komiteclaimlife`, `treatycontractout`); enam belas lainnya **kerangka**: hanya
`MODUL.md` (dengan rentang migrasi dan slot menunya) dan `docs/` bila dokumen `.scratch`-nya ada.
Kerangka tidak terdaftar di mana pun sampai ia mendapat `backend/modul.go`.

Nama modul = nama folder korpus tanpa spasi, huruf kecil (`NB Treaty In` → `nbtreatyin`) — satu nama
untuk folder, `const Nama` Go, `MODUL_AKTIF`, dan `KODE` kelompok `M_NAV_MENU`.

## 2. Siapa memegang apa

`..\..\.github\CODEOWNERS` menyatakannya; nama akun di sana penanda yang diisi work owner.

| Jalur | Pemilik | Yang terjadi bila disunting |
| --- | --- | --- |
| `APP_RNM/modul/<nama>/` | pemilik modul | pull request ditinjau pemilik modul |
| `APP_RNM/modul/<nama>/MODUL.md`, `APP_RNM/modul/<nama>/backend/migrations/9*` *(slot menu)* | tim inti — pemilik modul yang menulisnya | rentang, slot, pernyataan penjaga, `Status`, dan butir menu adalah keputusan yang ditinjau tim inti; tetap nol berkas bersama disunting |
| `APP_RNM/inti/backend/daftar/modul_<nama>_gen.go` | pemilik modul | berkas bangkitan — lahir sekali, saat modul dimulai (bab 4) |
| `APP_RNM/inti/` (selain di atas), `cmd/`, `frontend/`, `uji/` | tim inti | menyentuh **semua** modul: seluruh uji + tinjauan tim inti |
| `go.mod`, `go.sum`, `package.json`, `package-lock.json`, `vite.config.ts`, `tsconfig.json`, `Makefile` | tim inti | pustaka baru = pull request ke tim inti |
| `APP_RNM/modul/_templat/`, `APP_RNM/inti/docs/`, `APP_RNM/PANDUAN-*.md`, `.github/` | tim inti | — |

Yang **selalu** lewat tim inti, karena ia bersama menurut sifatnya:

- kontrak lintas modul baru (antarmuka di `inti/backend/kontrak`, bab 7);
- rentang migrasi atau slot menu tambahan (dari cadangan 760–899 dan 990–999);
- menyalakan menu modul Anda: berkas slot menu `9*` di folder migrasi modul Anda — SATU
  `UPDATE … SET DIMIGRASI = '1' WHERE KODE = '<nama>'`, nol `INSERT` (menu datar, keputusan work owner
  30-09-2026: satu modul satu menu) — ditinjau tim inti lewat `CODEOWNERS`. Dulu lewat angka kunci di
  `frontend/Shell.test.ts`, yang membuat dua modul berkonflik di baris yang sama;
- pernyataan untuk penjaga di `MODUL.md` (bab 6) — pengecualian penjaga adalah keputusan tim inti;
- pustaka npm/Go baru.

## 3. Alur kerja: cabang per modul → commit di folder sendiri → pull request → CI → merge

```powershell
git switch -c modul/nbtreatyin/tiket-01          # satu cabang per pekerjaan, berawalan nama modul
# ... bekerja HANYA di APP_RNM/modul/nbtreatyin/ ...
git add -A -- APP_RNM/modul/nbtreatyin
git commit -m "nbtreatyin: tiket 01 - ..."
git push -u origin modul/nbtreatyin/tiket-01
# buka pull request -> CODEOWNERS memanggil peninjau -> CI menjalankan gerbang SEMUA modul (bab 8) -> merge
```

Kenapa tidak ada merge manual: pekerjaan sehari-hari satu modul hanya menulis di foldernya sendiri, dan
setiap berkas yang dulu harus disunting tiap modul kini ditutup:

| Dulu disunting setiap modul | Kini |
| --- | --- |
| `modul/daftar.go` — daftar modul dan sambungan kontrak tangan | `inti/backend/daftar` **dibangkitkan**: satu berkas per modul; kontrak dinyatakan di `Pendaftaran()` modulnya dan disambung perakit |
| `frontend/src/modul/daftar.ts` — satu baris per modul | `frontend/daftar.ts` memakai `import.meta.glob` atas `modul/*/frontend/menu.ts` dan `rute.tsx` — nol baris per modul |
| `inti/labels.ts` `MODUL` — nama kelompok | nama modul di `menu.ts` modul (`KELOMPOK_<NAMA>`, nama folder korpus); label tombol dari `M_NAV_MENU.LABEL`; katalog 20 folder korpus di perakit, isinya tetap |
| Nomor migrasi "berikutnya" | rentang **40 nomor per modul**, ditetapkan di `MODUL.md` sejak awal |
| Menu baru di `inti/migrations/901…` | **slot menu** di folder migrasi modul sendiri (950–989): satu `UPDATE DIMIGRASI`, nol `INSERT` (901 kini meratakan tabel, milik inti) |
| Daftar per modul di penjaga (`letakStruktur`, peta penyuntikan, kaskade, …) | penjaga membaca `modul/*` secara umum; yang khusus modul dinyatakan di `MODUL.md`-nya (bab 6) |

Konflik masih mungkin — dan dikenali — hanya bila dua orang menyunting **modul yang sama** (satu pemilik
per modul mencegahnya) atau **berkas milik tim inti** (lewat tinjauannya).

## 4. Memulai modul kerangka

Contoh `nbtreatyin` (rentang migrasi 320–359, slot menu 968–969 — dari `MODUL.md`-nya).

**4.1 Backend** — `modul/nbtreatyin/backend/{models,repository,services,handlers}/`, masing-masing
mengimpor hanya `inti/backend/...` dan dirinya sendiri (arah `handlers → services → repository`). Lalu
`modul/nbtreatyin/backend/modul.go`.

⚠️ `//go:embed migrations/*.sql` tidak terkompilasi selama folder `migrations/` belum memuat satu
berkas `.sql` pun: tulis migrasi pertama (4.3) sebelum `go build`, atau — bila modul tidak bermigrasi,
seperti `treatycontractout` — buang baris `go:embed`, `berkasMigrasi`, dan `Migrasi:` sekaligus.

```go
// Package backend merakit modul NB Treaty In (`nbtreatyin`).
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbtreatyin/backend/handlers"
	"nusantarare/modul/nbtreatyin/backend/services"
)

//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "nbtreatyin"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`).
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi, // buang baris ini (dan go:embed) bila modul tidak bermigrasi
		// Membutuhkan: []inti.Kontrak{inti.KontrakDari[kontrak.X]()},  lalu inti.Ambil[kontrak.X](p)
		// Menyediakan: []inti.Kontrak{inti.KontrakDari[kontrak.Y]()},  lalu inti.Sediakan[kontrak.Y](p, ...)
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return Modul{svc: services.DariDasar(p.Dasar()), stubPelaku: p.Config().AuthStub}, nil
		},
	}
}

// Modul memenuhi inti.Modul.
type Modul struct {
	svc        *services.Service
	stubPelaku bool
}

func (Modul) Nama() string                                  { return Nama }
func (m Modul) DaftarkanRute(mux *http.ServeMux)            { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
```

**4.2 Daftar** — dari `APP_RNM`: `go generate ./inti/backend/daftar` (atau `make generate`). Ia menulis
SATU berkas baru, `inti/backend/daftar/modul_nbtreatyin_gen.go` — commit bersama modulnya; di
`CODEOWNERS` berkas itu milik Anda. Lupa menjalankannya? `inti/backend/daftar/bangkit` merah dan menyebut
perintah ini.

**4.3 Migrasi** — `backend/migrations/320_<nama>.sql` + `320_<nama>_down.sql`, nomor di rentang
`MODUL.md` Anda, tiga digit, LF tanpa BOM. Setiap tabel yang Anda buat harus tercatat di
`docs/STRUKTUR-TABEL-*.md` modul Anda — dokumen itu **mulai mengikat** sejak modul punya
`backend/modul.go` (`TestKolomDDLCocokDenganStruktur`).

**Tabel kerja `T_WORK_<…>`** `[keputusan work owner 01-10-2026]` — setiap tabel kerja baru (satu baris per
work object, padanan kasus Pega) memuat sembilan kolom ini, dengan **nama dan tipe seperti `T_WORK_CLAIM`**:

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | `VARCHAR2(32) NOT NULL`, PK | pengenal kasus (`pyID`, tanpa awalan kelas) |
| `COVER_KEY` | `VARCHAR2(32)`, FK ke `ID` tabel itu sendiri **tanpa `ON DELETE`**, ber-index | kasus induk; kosong bila XML tidak memakai penunjuk induk |
| `LINI` | `VARCHAR2(16)` | lini |
| `POSITION` | `VARCHAR2(64)` | posisi kasus menurut properti Pega yang disimpannya — `pyPosition` (peran) di klaim, `Position` (layar) di polis |
| `STATUS_WORK` | `VARCHAR2(32)` | `pyWorkStatus` VERBATIM |
| `CREATE_OP` | `VARCHAR2(64)` | akun pembuat (`pxCreateOperator`) |
| `CREATE_OP_NAME` | `VARCHAR2(128)` | nama pembuat — akun pelaku sampai login menyediakan nama tampilan |
| `TGL_CREATE` | `DATE` | waktu lahir kasus, `SYSDATE` |
| `TGL_UPDATE` | `DATE` | waktu ubah terakhir, `SYSDATE` di **setiap** pernyataan yang mengubah baris kasus |

Kolom posisi selalu bernama **`POSITION`** `[keputusan work owner 01-10-2026, migrasi Claim Life 023]` —
isinya mengikuti properti Pega modul itu, jadi **jangan membandingkan `POSITION` antartabel kerja**:
klaim menyimpan `pyPosition` (nama peran), polis menyimpan `Position` (layar `Offer`/`Premium`).
**Bukan kolom tabel kerja:** `TYPE` milik header kasus (`T_GENERAL_CLAIM.TYPE` sejak 023), dan `CASE_ID`
tidak ada — `ID` adalah pengenal kasus (migrasi data lama memakai CASEID warisan sebagai `ID`). Kolom
khusus modul tetap boleh (`TAHAP`, `SENDTO_*` klaim; `FLAG_ONGOING_POLICY` polis). Contoh:
`modul/claimlife/backend/migrations/001_t_work_claim.sql` (+ 016, 017, 023) dan
`modul/premiumlistlife/backend/migrations/059_seragam_kolom_t_work_polis.sql`; `LINI`, `POSITION`, dan
`STATUS_WORK` `T_WORK_POLIS` tetap `VARCHAR2(255)` karena datanya sudah ada (memperkecil dapat gagal).

**4.4 Frontend** — `modul/nbtreatyin/frontend/menu.ts`:

```ts
import type { MenuModul } from '../../../inti/frontend/modul'

export const NAMA_NBTREATYIN = 'nbtreatyin'
/** Nama modul - nama folder korpus VERBATIM (= M_NAV_MENU.LABEL barisnya). */
export const KELOMPOK_NBTREATYIN = 'NB Treaty In'
export const HALAMAN_NBTREATYIN = ['nbtreatyin-inbox'] as const
export type HalamanNbTreatyIn = (typeof HALAMAN_NBTREATYIN)[number]

/** Halaman yang dibuka tombol "NB Treaty In" di sidebar. */
export const HALAMAN_AWAL_NBTREATYIN: HalamanNbTreatyIn = 'nbtreatyin-inbox'

export const PENDAFTARAN_MENU: MenuModul<HalamanNbTreatyIn> = {
  nama: NAMA_NBTREATYIN,
  kelompok: KELOMPOK_NBTREATYIN,
  halaman: HALAMAN_NBTREATYIN,
  halamanAwal: HALAMAN_AWAL_NBTREATYIN,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    nbtreatyin: HalamanNbTreatyIn
  }
}
```

dan `rute.tsx` yang mengekspor `export const RUTE_MODUL: RuteModul<HalamanNbTreatyIn> = RuteNbTreatyIn`.
Nama ekspor `PENDAFTARAN_MENU` dan `RUTE_MODUL` SAMA di setiap modul: `frontend/daftar.ts` mengumpulkannya
dari folder, dan menolak folder yang hanya punya salah satunya — juga `halamanAwal` yang bukan salah satu
halaman modul itu. Menu DATAR (keputusan work owner 30-09-2026): satu modul satu tombol; label tombol
datang dari `M_NAV_MENU.LABEL`, halaman lain modul dibuka dari dalam halaman awalnya.

**4.5 Menu** — slot menu Anda: `backend/migrations/968_menu_nbtreatyin.sql` (+ `_down.sql`) yang HANYA
menyalakan `DIMIGRASI` baris modul Anda — satu `UPDATE`, nol `INSERT`; bentuk SQL persis di
`APP_RNM\PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6. Barisnya sudah ada sejak 900; `LABEL`-nya = nama
modul Anda di `menu.ts` VERBATIM. Berkas slot `9*` ditinjau tim inti (`CODEOWNERS`, bab 2); tidak ada
angka kunci di berkas bersama.

**4.6 `MODUL.md`** — `Status` → `dimigrasi` bersamaan dengan berkas slot menu Anda: uji Beranda,
sidebar, dan palet (`frontend/Beranda.test.ts`, `Shell.test.ts`, `daftar.sinkron.test.ts`,
`daftar.modulAktif.test.ts`) membaca `Status` ini dan merah bila ia tidak sesuai dengan menunya. Isi
`Prefix rute API` dan kontrak; tambahkan pernyataan untuk penjaga bila perlu (bab 6).

**4.7 Gerbang** — bab 8, lalu pull request.

## 5. Modul di luar dua puluh folder korpus

Bukan bagian migrasi ini: keputusan work owner. Bila diputuskan: salin `modul/_templat/` menjadi
`modul/<nama>/`, ajukan rentang migrasi dan slot menu dari cadangan (760–899, 990–999) lewat pull
request tim inti, dan minta tim inti menambah kelompok `M_NAV_MENU`-nya (isi awal 900 hanya memuat dua
puluh folder korpus; slot menu modul tidak boleh membuat kelompok) serta barisnya di
`frontend/katalogKorpus.ts`. Atau — keputusan work owner 04-10-2026 ("Modul luar korpus tanpa slot", delapan modul
master) — TANPA slot: `Slot menu` — di `MODUL.md`, dan baris menunya lahir `DIMIGRASI '1'` di langkah inti
(`modulLuarKorpus`, `inti/backend/penjaga/menu_test.go`).

## 6. Pernyataan untuk penjaga di `MODUL.md`

Penjaga `inti/backend/penjaga` berlaku untuk setiap modul dan tidak memuat nama modul. Yang khusus satu
modul dinyatakan di bab `## Pernyataan untuk penjaga` `MODUL.md`-nya — satu judul `###` per jenis, satu
tabel per judul, nilai di dalam backtick dibaca apa adanya. Judul yang tidak ada = tidak menyatakan apa pun.

| Judul `###` | Kolom tabel | Dibaca |
| --- | --- | --- |
| `Tabel warisan: dibaca, tidak dibuat` | Tabel · Alasan | `TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat` |
| `Nama STRUKTUR yang berbeda di DDL` | Jenis (`tabel`/`kolom`) · STRUKTUR · DDL · Alasan | `TestKolomDDLCocokDenganStruktur`, `TestGolonganTipeDDLCocokDenganStruktur` |
| `Nama terlarang di migrasi` | Nama · Sebab | `TestNamaYangDibuangTidakAda` (migrasi SEMUA modul) |
| `Kaskade ON DELETE CASCADE` | Awalan berkas · Relasi | `TestKaskadeHanyaPadaRelasiTerdaftar` — judul ada = kebijakan kaskade modul ini ditegakkan |
| `Penyuntikan wajib di handler` | Berkas · Penyusun · Wajib · Catatan | `TestHandlerMenyuntikkanImplementasiNyata`, `TestNolCabang501StubDiHandler` |
| `Pesan verbatim yang bukan nama orang` | Paket · Konstanta · Alasan | `TestNolNamaOrangDiKode` |

Contoh lengkap: `APP_RNM\modul\claimlife\MODUL.md` dan `APP_RNM\modul\treatycontractout\MODUL.md`.
Pernyataan yang berbentuk salah menggagalkan penjaga dengan kalimat yang menyebut berkas dan judulnya —
termasuk judul `###` yang bukan salah satu dari enam di atas (salah ketik satu huruf dulu membuat
pernyataan itu diam-diam tidak dibaca).

## 7. Kontrak lintas modul

Modul tidak pernah mengimpor modul lain. Bila modul A butuh sesuatu dari modul B:

1. tim inti menambah antarmukanya di `inti/backend/kontrak` (tanpa implementasi) — pull request;
2. `Pendaftaran()` B menyatakan `Menyediakan: []inti.Kontrak{inti.KontrakDari[kontrak.X]()}` dan
   `Bangun`-nya memanggil `inti.Sediakan[kontrak.X](p, implementasi)`;
3. `Pendaftaran()` A menyatakan `Membutuhkan` yang sama dan membacanya dengan `inti.Ambil[kontrak.X](p)`.

Perakit (`inti/backend/perakit.go`) membangun B sebelum A. Kontrak yang dibutuhkan tanpa penyedia, dua
penyedia, ketergantungan melingkar, atau pernyataan yang tidak ditepati membuat backend **menolak
menyala** dengan kalimat yang menyebut kontrak dan modulnya — tidak pernah nil diam-diam.

## 8. Gerbang — setiap commit, dan CI setiap pull request

Dari `APP_RNM`:

```powershell
go generate ./inti/backend/daftar ; git status --short    # daftar bangkitan tidak boleh berubah
go build ./... ; go build -tags db ./... ; go vet ./... ; go vet -tags db ./... ; gofmt -l cmd inti modul uji
go test ./... ; go test -tags db ./...                     # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx tsc --noEmit ; npx vitest run ; npx vite build
```

Uji satu modul saja selama bekerja: `go test ./modul/<nama>/...` dan `npx vitest run modul/<nama>`.
Pull request tetap menjalankan gerbang **semua** modul: modul Anda dapat mematahkan penjaga lintas
aplikasi (impor, rentang, menu, alamat layanan).

## 9. Yang masih bersama

Dicatat supaya tidak ada yang mengira semuanya sudah per modul:

| Berkas | Pemilik | Kapan pengembang modul menyentuhnya |
| --- | --- | --- |
| `frontend/Beranda.tsx` — kartu antrean Claim Life | tim inti | bila kartu Beranda Claim Life berubah |
| `inti/frontend/labels.ts` — `MENU` (label kedua butir Claim Life) | tim inti | bila label butir Claim Life berubah — Shell memakai `MENU.inbox` sebagai judul cadangan, jadi memindahkannya mengubah perilaku |
| `frontend/Shell.test.ts` — kelompok Treaty Contract Out tepat satu butir | tim inti | bila keputusan tco5 berubah (keputusan work owner) |
| `inti/backend/penjaga/lintasaplikasi_test.go` — daftar klien HTTP keluar yang disetujui | tim inti | bila modul memanggil layanan luar (butuh persetujuan manusia) |
| `inti/backend/penjaga/lintasaplikasi_test.go` — dua uji tco4 Treaty Contract Out | tim inti | bila keputusan tco4 berubah (keputusan work owner) |
| `inti/backend/penjaga/menu_test.go` — isi awal 900 | tim inti | tidak pernah: 900 tidak disunting |
| `inti/backend/kontrak/` | tim inti | kontrak baru (bab 7) |
