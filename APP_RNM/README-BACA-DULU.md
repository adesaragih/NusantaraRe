# APP_RNM — baca dulu sebelum menyentuh kode

Folder ini memuat **seluruh aplikasi** hasil migrasi Pega → Go + React + Oracle. Dokumen
(spec, tiket, ADR, discovery) tetap berada di folder induk `OUTPUT_HASIL_RNM\`. Ditulis
25 September 2026 untuk pembaca yang **baru mengenal Go dan React**.

## 1. Peta folder — satu kalimat per folder

**Struktur tim satu folder per modul** *(keputusan work owner 30-09-2026,
`..\PROMPT-STRUKTUR-TIM-SATU-FOLDER-PER-MODUL.md`)*: satu orang fullstack memegang satu modul, dan
segala milik modul itu — backend, frontend, dokumen — tinggal di **satu folder**
`modul/<nama>/`. Alur kerja tim, kepemilikan, dan cara memulai modul:
`PANDUAN-TIM-PER-MODUL.md`. Cabang git dan pembagian tim: `PANDUAN-CABANG-GIT.md`.

| Folder / berkas | Isinya, dalam bahasa sehari-hari |
| --- | --- |
| `cmd/api/` | Pintu masuk backend. Membaca pengaturan, menyambung Oracle, memasang modul yang AKTIF (`MODUL_AKTIF`, `rakit.go`), lalu menunggu permintaan |
| `inti/backend/` | **Kode Go bersama semua modul**: koneksi dan transaksi (`db/`), outbox efek keluar, resolver layanan, jejak audit, penomor, pelari migrasi, jawaban galat (`galat/`), unggahan, uang (`uang/`), alat bantu (`utils/`), menu (`menu/`). Paket akarnya bernama `backend` dan diimpor dengan alias `inti` — `inti.Dasar`, `inti.Modul` |
| `inti/backend/kontrak/` | Antarmuka lintas modul, TANPA implementasi (`PembacaPolis`, `KlaimKomite`) — modul tidak saling mengimpor |
| `inti/backend/perakit.go` | Perakit: menyambung kontrak menurut pernyataan `Menyediakan`/`Membutuhkan` setiap modul; kontrak yang dibutuhkan tanpa penyedia = backend menolak menyala |
| `inti/backend/daftar/` | Daftar modul Go — **dibangkitkan** `go generate ./inti/backend/daftar` dari folder `modul/`: satu berkas `modul_<nama>_gen.go` per modul yang punya `backend/modul.go`. Jangan disunting tangan |
| `inti/backend/config/` | Membaca **env var** (variabel lingkungan) seperti `ORACLE_DSN`. Tidak ada alamat atau kata sandi yang ditulis di kode |
| `inti/backend/penjaga/` | Test penjaga yang berlaku untuk SELURUH aplikasi: impor lintas modul, rentang migrasi dan slot menu, higiene migrasi, alamat layanan, nama orang, nama tabel telanjang. Membaca `modul/*` secara umum; yang khusus satu modul dinyatakan di `MODUL.md`-nya |
| `inti/backend/migrations/` | Berkas `.sql` tabel **lintas modul** milik `inti`, rentang 900–949. Hari ini satu: `M_NAV_MENU` (900), tabel menu yang dibaca `GET /api/menu` |
| `inti/frontend/` | Kerangka React bersama: Shell, `components/ui/dasar`, `klien.ts`, `lib/`, `hooks/`, `store/`, `labels.ts`, `styles.css`; `uji/sumber.ts` menyebut letak kode frontend bagi uji yang membaca sumber |
| `modul/<nama>/` | **Satu modul = satu folder = satu pemilik.** Dua puluh folder — satu per folder korpus; `claimlife`, `premiumlistlife`, `komiteclaimlife`, `treatycontractout` sudah dimigrasi, enam belas lainnya kerangka (hanya `MODUL.md` dan `docs/`) |
| `modul/<nama>/MODUL.md` | Pemilik, rentang migrasi dan slot menu (dibaca penjaga), prefix rute API, kontrak yang disediakan/dipakai, cara menguji modul itu saja, dan pernyataan khusus modul untuk penjaga |
| `modul/<nama>/backend/modul.go` | `Pendaftaran()`: nama modul, migrasinya, kontraknya, dan cara merakitnya |
| `modul/<nama>/backend/handlers/` | Penerima permintaan HTTP. Tugasnya sempit: baca permintaan, panggil *service*, tulis jawaban JSON |
| `modul/<nama>/backend/services/` | Aturan dagang dan perakitan data. Di sinilah "klaim punya peserta, peserta punya baris" disusun |
| `modul/<nama>/backend/repository/` | Satu-satunya lapisan yang berbicara ke Oracle. Seluruh SQL modul itu ada di sini, dan **hanya** di sini |
| `modul/<nama>/backend/migrations/` | Berkas `.sql` bernomor yang membentuk tabel modul itu — nomor di rentang modulnya (`MODUL.md`), menu di slot menunya. Satu berkas = satu langkah, dan tiap langkah punya pasangan `_down.sql` untuk membatalkannya |
| `modul/<nama>/backend/models/` | Bentuk data (struct): `Klaim`, `Peserta`, `BarisAdjustment`. `Money` dan `Ratio` di `inti/backend/uang/` |
| `modul/claimlife/backend/repository/barislamakolom.go` | Daftar **62 kolom** tabel datar warisan beserta tipenya dari katalog Oracle, ditulis SEKALI; **55** di antaranya yang ditulis rule Pega |
| `modul/<nama>/frontend/` | Tampilan React modul itu: `pages/`, `components/`, `labels.ts`, `api.ts`, `menu.ts` (`PENDAFTARAN_MENU`), `rute.tsx` (`RUTE_MODUL`) |
| `modul/<nama>/docs/` | Spec, tiket (`issues/`), grilling, catatan modul itu — dulu `..\.scratch\<nama-panjang>\` |
| `modul/_templat/` | Templat folder modul (hanya `MODUL.md` dan `docs/`) |
| `uji/skemauji/` | Menyiapkan skema uji Oracle untuk test bertag `db`: menjalankan migrasi yang sama dengan aplikasi, mengisi fixture buatan, lalu membongkarnya |
| `frontend/` | **Perakit** React (TypeScript, `.tsx`): `index.html`, `main.tsx`, `App.tsx`, `Beranda.tsx`, `daftar.ts` (daftar modul dari folder lewat `import.meta.glob`), `katalogKorpus.ts`. Juga `public/`, `dist/`, dan `.env` frontend |
| `package.json` · `vite.config.ts` · `tsconfig.json` · `node_modules/` | Di folder ini (APP_RNM), bukan di `frontend/`: perintah `npm` dijalankan dari sini |
| `Makefile` | Daftar perintah: jalankan, bangun, uji, bangkitkan daftar. Setiap target adalah satu-dua perintah biasa |
| `bin/` | Hasil `go build` — tidak masuk git |

**Tabel nama modul** — satu-satunya sumber, tanpa singkatan *(keputusan work owner 30-09-2026, `PROMPT-REFACTOR-NAMA-MODUL.md`)*:

| Modul korpus | Folder `APP_RNM/modul/…` = `MODUL_AKTIF` *(tanpa tanda hubung)* | Dokumen lama `.scratch/…` *(sebelum 30-09-2026)* |
| --- | --- | --- |
| Claim Life | `claimlife` | `claim-life` |
| PremiumList Life | `premiumlistlife` *(dulu `premiumlist`)* | `premiumlist-life` |
| Komite Claim Life | `komiteclaimlife` *(dulu `komite`)* | `komite-claim-life` |
| Treaty Contract Out | `treatycontractout` *(dulu `treaty`)* | `treaty-contract-out` |
| *(kerangka, mis. NB Treaty In)* | `nbtreatyin` | `nb-treaty-in` |

Aturan: **nama modul = nama folder korpus tanpa spasi, huruf kecil = nama `.scratch` lama tanpa tanda
hubung**. Satu nama untuk folder, `const Nama` Go, `MODUL_AKTIF`, dan `KODE` kelompok `M_NAV_MENU`.

**Arah ketergantungan** — selalu satu arah, tidak pernah memotong:

```
handlers  →  services  →  repository  →  Oracle
```

`handlers` tidak boleh menyentuh `repository` langsung. Kalau Anda menulis SQL di luar
`repository`, atau memanggil `repository` dari `handlers`, itu salah tempat.

**Antarmodul** — `modul/X` tidak pernah mengimpor `modul/Y` (Go maupun TypeScript); ia menyatakan
kontrak `inti/backend/kontrak` yang ia sediakan atau butuhkan di `Pendaftaran()`, dan perakit
menyambungnya. Cara deploy sebagian modul dan menambah menu: `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`.

## 2. Satu permintaan, dari layar sampai Oracle dan kembali

Contoh: pengguna mengetik `UJI-KLAIM-1` lalu menekan **Buka**.

| Langkah | Berkas | Yang terjadi |
| ---: | --- | --- |
| 1 | `modul/claimlife/frontend/pages/KlaimLife.tsx` | Fungsi `cari` dipanggil; halaman masuk keadaan "memuat" |
| 2 | `modul/claimlife/frontend/api.ts` | `ambilKlaimLife(id)` mengirim `GET /api/klaim-life/UJI-KLAIM-1` lewat `minta` (`inti/frontend/klien.ts`) |
| 3 | `modul/claimlife/backend/handlers/klaimlife.go` | Backend menerima, mengambil `{id}` dari alamat, memanggil *service* |
| 4 | `modul/claimlife/backend/services/klaimlife.go` | `Ambil` meminta header, peserta, dan baris ke *repository*, lalu merakitnya |
| 5 | `modul/claimlife/backend/repository/klaimlife.go` | Tiga query SQL ke Oracle. Uang diminta sebagai **teks** lewat `TO_CHAR` |
| 6 | `modul/claimlife/backend/models/klaimlife.go` | `MarshalJSON` menentukan bentuk JSON yang dikirim balik; kode status diterjemahkan menjadi kata |
| 7 | `KlaimLife.tsx` | Jawaban disimpan lewat `setKlaim`; React menggambar ulang tabel |

Nama medan JSON di langkah 6 (`nomorKlaim`, `jumlahKlaim`, …) **harus sama persis** dengan
`interface` di `api.ts`. Kalau salah satu diubah, yang lain ikut diubah.

## 3. Urutan membaca yang disarankan

**Go** — dari luar ke dalam: `cmd/api/main.go` → `cmd/api/rakit.go` → `inti/backend/config/config.go` →
`inti/backend/daftar/daftar.go` → `inti/backend/perakit.go` → `modul/claimlife/backend/modul.go` →
`modul/claimlife/backend/handlers/handlers.go` → `…/handlers/klaimlife.go` → `…/services/klaimlife.go` →
`…/repository/klaimlife.go` → `…/models/klaimlife.go` → `inti/backend/uang/`.

**React** — dari titik masuk: `frontend/main.tsx` → `frontend/App.tsx` → `frontend/daftar.ts` →
`modul/claimlife/frontend/menu.ts` → `…/rute.tsx` → `…/pages/KlaimLife.tsx` → `…/api.ts` →
`inti/frontend/klien.ts` → `inti/frontend/store/index.ts`.

Setiap berkas punya komentar kepala yang menjelaskan **untuk apa berkas ini** dan istilah
yang dipakainya.

## 4. Istilah yang sering muncul

| Istilah | Artinya di proyek ini |
| --- | --- |
| **env var** | Nilai pengaturan yang dibaca dari lingkungan proses, bukan dari kode. Contohnya di `.env.example` |
| **handler** | Fungsi Go yang menjawab satu alamat HTTP |
| **service** | Fungsi Go yang memuat aturan dagang; boleh memanggil banyak *repository* |
| **repository** | Fungsi Go yang menjalankan SQL. Tidak tahu apa-apa soal HTTP |
| **model** | Struct Go yang menggambarkan satu jenis data |
| **seam** | Titik tempat test menggerakkan sistem: lewat HTTP, lewat *repository*, atau lewat *service* murni |
| **fixture** | Data contoh buatan untuk test — selalu berawalan `UJI-`, tidak pernah data sungguhan |
| **skema uji** | Tabel sementara yang dibuat test lalu dibuang lagi (`uji/skemauji/`) |
| **test statik** | Test yang membaca kode sumber sebagai teks, bukan menjalankannya. Dipakai untuk aturan yang tidak dapat dijaga kompilator — misalnya "`Hapus` tidak boleh dipanggil di luar test" (`batasanpemakaian_test.go`) |
| **komponen** (React) | Fungsi yang mengembalikan tampilan |
| **`useState`** | Kotak penyimpan nilai di dalam komponen; mengganti isinya menggambar ulang tampilan |
| **`interface`** (TypeScript) | Bentuk data; hanya ada saat pemeriksaan tipe, tidak ikut ke browser |

## 5. Menjalankan — dari folder `APP_RNM\` ini

Bila `go`, `node`, atau `npm` tidak dikenal di terminal, tambahkan dulu ke PATH sesi
(rinciannya di `..\PANDUAN-RANTAI-ALAT.md` §1):

```powershell
$env:Path = 'C:\Program Files\Go\bin;C:\Program Files\nodejs;' + $env:Path
```

| Tujuan | Perintah | Yang diharapkan |
| --- | --- | --- |
| Uji backend tanpa Oracle | `go vet ./...` lalu `go test ./...` | `ok` di setiap paket `inti/...`, `modul/...`, `cmd/api` |
| Bentuk tabel ke Oracle uji | `go run ./cmd/api -migrate` | menjalankan berkas di `modul/*/backend/migrations/` SEMUA modul dan `inti/backend/migrations/` (tidak ikut `MODUL_AKTIF`) sekali masing-masing; aman diulang, dan **menolak** berjalan bila `IS_PEGA_PROD=true` |
| Uji backend **dengan** Oracle uji | `go test -tags=db ./...` | perlu `ORACLE_DSN` + `ORACLE_SCHEMA` + `ORACLE_SKEMA_UJI=true`; tanpa `ORACLE_DSN` test **melewati** dengan pesan, bukan lulus diam-diam |
| Jalankan backend | `go run ./cmd/api` | log `http: mendengarkan di :8080`; `MODUL_AKTIF=claimlife,komiteclaimlife` memasang sebagian modul |
| Periksa tipe frontend | `npm run typecheck` — dari `APP_RNM/`, bukan `frontend/` *(sejak 30-09-2026)* | tidak mencetak galat |
| Uji frontend | `npm test` | seluruh berkas uji `passed` |
| Jalankan frontend | `npm run dev` | buka `http://localhost:5174/` |
| Bangun frontend | `npm run build` | menjalankan `tsc` dulu, lalu Vite membuat `dist/` |

⛔ **`go test -tags=db` MENGHAPUS tabel di skema yang ditunjuk `ORACLE_SCHEMA`** — termasuk
`OS_AKSEPTASI_KLAIM_LIFE`, tabel datar warisan. Karena itu ia menolak berjalan kecuali
`ORACLE_SKEMA_UJI=true` **dan** `ORACLE_SCHEMA` bukan `POOLDATA`; penolakannya berupa **galat**,
bukan lewat. Tunjuk hanya skema uji kosong yang dibuat DBA khusus untuk itu, tidak pernah skema
mana pun yang memuat tabel warisan sungguhan — sekalipun di instance pengembangan. `-migrate`
tidak ikut dipagari: ia hanya `CREATE`, tidak pernah `DROP`.

Dengan `make`: `make test` · `make typecheck` · `make check` · `make test-db` · `make build` · `make generate`
(tulis ulang daftar modul Go sesudah menambah atau membuang `modul/<nama>/backend/modul.go`).

## 6. Tiga aturan yang paling sering ditegur reviewer

1. **Uang tidak pernah `float`, dan tidak pernah `Number` di JavaScript.** Oracle → Go → JSON →
   React membawanya sebagai **teks** desimal. Float membulatkan diam-diam. *(ADR-U-0003, ADR-U-0016)*
2. **Kode dan nomor tetap teks.** `"006"` bukan `6`. Mengubahnya menjadi bilangan menghilangkan
   nol di depan dan memecahkan penggolong. *(ADR-U-0022)*
3. **Akhiran baris berkas `.sql` dan `.go` selalu LF, tanpa BOM.** Alat Windows menulis CRLF dan BOM diam-diam; keduanya lolos ke teks SQL dan ditolak Oracle, sementara seluruh test tanpa basis data tetap hijau. Dijaga `.gitattributes` di akar dan dua test penjaga.
4. **Nol alamat host, kata sandi, atau nomor polis sungguhan di kode dan test.** Alamat dari
   env var; data uji berawalan `UJI-`. *(ADR-U-0004)*

Rujukan `ADR-U-nnnn`, `ADR-D-<modul>-nnnn`, dan `ADR-F-nnnn` menunjuk dokumen keputusan yang dulu
tinggal di luar `APP_RNM` (`docs\bersama\adr\`, `dastin\`, `jefri\`). Folder itu dipangkas dari repo
pada 1 Oktober 2026 supaya repositori hanya memuat aplikasi; isinya tetap tersimpan di riwayat git dan
dapat dipulihkan dengan `git checkout <commit-sebelum-pangkas> -- <jalur>`.

## 7. Dari mana pekerjaan datang

Setiap baris kode menunjuk satu **tiket**. Tiket Claim Life ada di
`modul\claimlife\docs\issues\`; spec-nya di `modul\claimlife\docs\spec.md` — dokumen setiap modul di
`modul\<nama>\docs\`. Brief kerja
untuk sesi implementasi: `..\PROMPT-IMPLEMENTASI-GO-REACT.md`.
