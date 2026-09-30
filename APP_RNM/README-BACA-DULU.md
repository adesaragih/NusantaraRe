# APP_RNM — baca dulu sebelum menyentuh kode

Folder ini memuat **seluruh aplikasi** hasil migrasi Pega → Go + React + Oracle. Dokumen
(spec, tiket, ADR, discovery) tetap berada di folder induk `OUTPUT_HASIL_RNM\`. Ditulis
25 September 2026 untuk pembaca yang **baru mengenal Go dan React**.

## 1. Peta folder — satu kalimat per folder

| Folder / berkas | Isinya, dalam bahasa sehari-hari |
| --- | --- |
| `cmd/api/` | Pintu masuk backend. Membaca pengaturan, menyambung Oracle, memasang modul yang AKTIF (`MODUL_AKTIF`, `rakit.go`), lalu menunggu permintaan |
| `inti/` | **Kode bersama semua modul** (refactor bentuk B, 30-09-2026): koneksi dan transaksi (`inti/db`), outbox efek keluar, resolver layanan, jejak audit, penomor, pelari migrasi, jawaban galat HTTP, gerbang unggahan, uang. Tidak pernah mengimpor modul |
| `inti/kontrak/` | Antarmuka lintas modul, TANPA implementasi (`PembacaPolis`, `KlaimKomite`) — modul tidak saling mengimpor |
| `inti/config/` | Membaca **env var** (variabel lingkungan) seperti `ORACLE_DSN`. Tidak ada alamat atau kata sandi yang ditulis di kode |
| `inti/penjaga/` | Test penjaga yang berlaku untuk SELURUH aplikasi: impor lintas modul, higiene migrasi, alamat layanan, nama orang, nama tabel telanjang |
| `modul/daftar.go` | Daftar modul: merakit `Service` tiap modul dan menyambung `inti/kontrak` |
| `modul/<nama>/` | Satu modul — `claimlife`, `premiumlistlife`, `komiteclaimlife`, `treatycontractout` (tabel nama di bawah) — dengan lapisannya sendiri: |
| `modul/<nama>/handlers/` | Penerima permintaan HTTP. Tugasnya sempit: baca permintaan, panggil *service*, tulis jawaban JSON |
| `modul/<nama>/services/` | Aturan dagang dan perakitan data. Di sinilah "klaim punya peserta, peserta punya baris" disusun |
| `modul/<nama>/repository/` | Satu-satunya lapisan yang berbicara ke Oracle. Seluruh SQL modul itu ada di sini, dan **hanya** di sini |
| `modul/<nama>/migrations/` | Berkas `.sql` bernomor yang membentuk tabel modul itu (rentang nomor per modul). Satu berkas = satu langkah, dan tiap langkah punya pasangan `_down.sql` untuk membatalkannya. Ditanam ke biner, jadi tidak perlu dicari di disk saat program jalan |
| `modul/claimlife/repository/barislamakolom.go` | Daftar **62 kolom** tabel datar warisan beserta tipenya dari katalog Oracle, ditulis SEKALI; **55** di antaranya yang ditulis rule Pega. Pembaca, penulis fixture, dan tabel tiruan mengambil daftar yang sama, sehingga urutan `SELECT` dan urutan `Scan` tidak mungkin berselisih |
| `modul/<nama>/models/` | Bentuk data (struct): `Klaim`, `Peserta`, `BarisAdjustment`. `Money` dan `Ratio` kini di `inti/uang/` |
| `uji/skemauji/` | Menyiapkan skema uji Oracle untuk test bertag `db`: menjalankan migrasi yang sama dengan aplikasi, mengisi fixture buatan, lalu membongkarnya |
| `inti/utils/` | Alat bantu umum: konversi teks ↔ desimal, format tanggal |
| `frontend/` | Tampilan React (TypeScript, berkas `.tsx`). Dijalankan Vite. `src/inti/` bersama, `src/modul/<nama>/` per modul, `src/App.tsx` merakit |
| `Makefile` | Daftar perintah: jalankan, bangun, uji. Setiap target adalah satu-dua perintah biasa |
| `bin/` | Hasil `go build` — tidak masuk git |

**Tabel nama modul** — satu-satunya sumber, tanpa singkatan *(keputusan work owner 30-09-2026, `PROMPT-REFACTOR-NAMA-MODUL.md`)*:

| Modul korpus | Backend Go `APP_RNM/modul/…` *(tanpa tanda hubung)* | Frontend `frontend/src/modul/…` dan `.scratch/…` | Nilai `MODUL_AKTIF` |
| --- | --- | --- | --- |
| Claim Life | `claimlife` *(tetap)* | `claim-life` | `claimlife` |
| PremiumList Life | `premiumlist` → **`premiumlistlife`** | `premiumlist` → **`premiumlist-life`** | `premiumlistlife` |
| Komite Claim Life | `komite` → **`komiteclaimlife`** | `komite` → **`komite-claim-life`** | `komiteclaimlife` |
| Treaty Contract Out | `treaty` → **`treatycontractout`** | `treaty` → **`treaty-contract-out`** | `treatycontractout` |
| *(modul berikutnya, mis. NB FacIn)* | `nbfacin` | `nb-facin` | `nbfacin` |

Aturan: **nama backend = nama dokumen `.scratch` tanpa tanda hubung**; frontend memakai nama `.scratch` persis. Nama paket Go = nama folder.

**Arah ketergantungan** — selalu satu arah, tidak pernah memotong:

```
handlers  →  services  →  repository  →  Oracle
```

`handlers` tidak boleh menyentuh `repository` langsung. Kalau Anda menulis SQL di luar
`repository`, atau memanggil `repository` dari `handlers`, itu salah tempat.

**Antarmodul** — `modul/X` tidak pernah mengimpor `modul/Y`; ia memakai antarmuka `inti/kontrak`
yang disambung `modul/daftar.go`. Cara deploy sebagian modul, commit per folder, dan menambah
modul: `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`.

## 2. Satu permintaan, dari layar sampai Oracle dan kembali

Contoh: pengguna mengetik `UJI-KLAIM-1` lalu menekan **Buka**.

| Langkah | Berkas | Yang terjadi |
| ---: | --- | --- |
| 1 | `frontend/src/modul/claim-life/pages/KlaimLife.tsx` | Fungsi `cari` dipanggil; halaman masuk keadaan "memuat" |
| 2 | `frontend/src/modul/claim-life/api.ts` | `ambilKlaimLife(id)` mengirim `GET /api/klaim-life/UJI-KLAIM-1` lewat `minta` (`inti/klien.ts`) |
| 3 | `modul/claimlife/handlers/klaimlife.go` | Backend menerima, mengambil `{id}` dari alamat, memanggil *service* |
| 4 | `modul/claimlife/services/klaimlife.go` | `Ambil` meminta header, peserta, dan baris ke *repository*, lalu merakitnya |
| 5 | `modul/claimlife/repository/klaimlife.go` | Tiga query SQL ke Oracle. Uang diminta sebagai **teks** lewat `TO_CHAR` |
| 6 | `modul/claimlife/models/klaimlife.go` | `MarshalJSON` menentukan bentuk JSON yang dikirim balik; kode status diterjemahkan menjadi kata |
| 7 | `KlaimLife.tsx` | Jawaban disimpan lewat `setKlaim`; React menggambar ulang tabel |

Nama medan JSON di langkah 6 (`nomorKlaim`, `jumlahKlaim`, …) **harus sama persis** dengan
`interface` di `api.ts`. Kalau salah satu diubah, yang lain ikut diubah.

## 3. Urutan membaca yang disarankan

**Go** — dari luar ke dalam: `cmd/api/main.go` → `cmd/api/rakit.go` → `inti/config/config.go` →
`modul/daftar.go` → `modul/claimlife/modul.go` → `modul/claimlife/handlers/handlers.go` →
`modul/claimlife/handlers/klaimlife.go` → `modul/claimlife/services/klaimlife.go` →
`modul/claimlife/repository/klaimlife.go` → `modul/claimlife/models/klaimlife.go` → `inti/uang/`.

**React** — dari titik masuk: `frontend/src/main.tsx` → `App.tsx` → `modul/daftar.ts` →
`modul/claim-life/rute.tsx` → `modul/claim-life/pages/KlaimLife.tsx` → `modul/claim-life/api.ts` →
`inti/klien.ts` → `inti/store/index.ts`.

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
| Bentuk tabel ke Oracle uji | `go run ./cmd/api -migrate` | menjalankan berkas di `modul/*/migrations/` SEMUA modul (tidak ikut `MODUL_AKTIF`) sekali masing-masing; aman diulang, dan **menolak** berjalan bila `IS_PEGA_PROD=true` |
| Uji backend **dengan** Oracle uji | `go test -tags=db ./...` | perlu `ORACLE_DSN` + `ORACLE_SCHEMA` + `ORACLE_SKEMA_UJI=true`; tanpa `ORACLE_DSN` test **melewati** dengan pesan, bukan lulus diam-diam |
| Jalankan backend | `go run ./cmd/api` | log `http: mendengarkan di :8080`; `MODUL_AKTIF=claimlife,komiteclaimlife` memasang sebagian modul |
| Periksa tipe frontend | `cd frontend` lalu `npm run typecheck` | tidak mencetak galat |
| Uji frontend | `npm test` | seluruh berkas uji `passed` |
| Jalankan frontend | `npm run dev` | buka `http://localhost:5173/` |
| Bangun frontend | `npm run build` | menjalankan `tsc` dulu, lalu Vite membuat `dist/` |

⛔ **`go test -tags=db` MENGHAPUS tabel di skema yang ditunjuk `ORACLE_SCHEMA`** — termasuk
`OS_AKSEPTASI_KLAIM_LIFE`, tabel datar warisan. Karena itu ia menolak berjalan kecuali
`ORACLE_SKEMA_UJI=true` **dan** `ORACLE_SCHEMA` bukan `POOLDATA`; penolakannya berupa **galat**,
bukan lewat. Tunjuk hanya skema uji kosong yang dibuat DBA khusus untuk itu, tidak pernah skema
mana pun yang memuat tabel warisan sungguhan — sekalipun di instance pengembangan. `-migrate`
tidak ikut dipagari: ia hanya `CREATE`, tidak pernah `DROP`.

Dengan `make`: `make test` · `make typecheck` · `make check` · `make test-db` · `make build`.

## 6. Tiga aturan yang paling sering ditegur reviewer

1. **Uang tidak pernah `float`, dan tidak pernah `Number` di JavaScript.** Oracle → Go → JSON →
   React membawanya sebagai **teks** desimal. Float membulatkan diam-diam. *(ADR-U-0003, ADR-U-0016)*
2. **Kode dan nomor tetap teks.** `"006"` bukan `6`. Mengubahnya menjadi bilangan menghilangkan
   nol di depan dan memecahkan penggolong. *(ADR-U-0022)*
3. **Akhiran baris berkas `.sql` dan `.go` selalu LF, tanpa BOM.** Alat Windows menulis CRLF dan BOM diam-diam; keduanya lolos ke teks SQL dan ditolak Oracle, sementara seluruh test tanpa basis data tetap hijau. Dijaga `.gitattributes` di akar dan dua test penjaga.
4. **Nol alamat host, kata sandi, atau nomor polis sungguhan di kode dan test.** Alamat dari
   env var; data uji berawalan `UJI-`. *(ADR-U-0004)*

Rujukan `ADR-U-nnnn` menunjuk `..\docs\adr\`; `ADR-D-<modul>-nnnn` ke `..\dastin\...\docs\adr\`;
`ADR-F-nnnn` ke `..\jefri\OUTPUT FIX\adr\`.

## 7. Dari mana pekerjaan datang

Setiap baris kode menunjuk satu **tiket**. Tiket Claim Life ada di
`..\.scratch\claim-life\issues\`; spec-nya di `..\.scratch\claim-life\spec.md`. Brief kerja
untuk sesi implementasi: `..\PROMPT-IMPLEMENTASI-GO-REACT.md`.
