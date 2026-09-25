# PROMPT — sesi implementasi berikutnya: Claim Life tiket 14 (PREFACTOR skema relasional) + penutupan tiket 01

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya *(revisi 25 September 2026 sore:
> akar aplikasi `APP_RNM\`, frontend `.tsx`, kode mudah dibaca pemula)*. Berkas ini **hanya** menambah
> keadaan awal, prasyarat, dan keputusan untuk sesi ini.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> `.scratch\claim-life\issues\14-skema-relasional-klaim-dan-migrasi.md` *(seluruhnya, termasuk
> §Blocker)* → `.scratch\claim-life\spec.md` §2b dan bab *Acceptance Criteria* →
> `.scratch\claim-life\STRUKTUR-TABEL-CLAIM-LIFE.md` → `docs\adr\` 0003 · 0006 · 0009 · 0011 · 0016 ·
> 0022 · 0027 · 0029 · 0033 → tiket 01, bab `## Implementasi`.
>
> **SESI INI:** claim-life · tiket **01 (penutupan sisa AC)** + **14**

---

## 0. KEADAAN AWAL — 25 September 2026, sore

| | Keadaan | Akibat untuk sesi ini |
| --- | --- | --- |
| `HEAD` | `6136a9d` — tiket 01 *claimed*, 4 dari 7 AC | tiga AC sisanya menunggu Oracle |
| Working tree | **44 perubahan belum commit**: pemindahan ke `APP_RNM\`, konversi `.tsx`, README, ralat dokumen | langkah 0 di §4 meng-commit-nya **sebelum** titik tetap diambil |
| Oracle pengembangan | **belum ada** — `ORACLE_DSN` kosong. `tnsnames.ora` mesin ini memuat 8 alias DEV *(isinya tidak pernah dikutip ke mana pun)* | tanpa instance, sesi berjalan di **jalur B** (§4) |
| Rantai alat | `go` 1.26.8 dan `node` 24 ada di mesin, **tidak** di PATH shell; `make` **tidak ada** | jalankan perintah di balik target Makefile langsung; sisipkan PATH (§1) |
| `.env` | `.env.example` hanya dokumentasi — **tidak ada kode yang memuat `.env`**; `config.Load` membaca env var proses | env var disetel di terminal, bukan di berkas |
| Uji yang lulus saat ini | `go vet`, `go vet -tags=db`, `gofmt`, `go build`, 7 test Go, `tsc --noEmit`, 5 test JS, `vite build` | titik nol yang harus tetap hijau |

---

## 1. PRASYARAT MANUSIA — dikerjakan work owner sebelum sesi dimulai

1. **Minta ke DBA** satu instance pengembangan: DSN, nama skema uji, dan hak `CREATE TABLE`,
   `CREATE SEQUENCE`, `CREATE INDEX`, `DROP` pada skema itu. ⛔ Bukan alias berpenanda PROD.
   Uji sambungannya dulu dengan `C:\oracle12i\bin\sqlplus.exe <pengguna>@<alias DEV>`.
2. **Setel env var di terminal yang akan menjalankan sesi** — tidak pernah ditulis ke berkas repo,
   tidak pernah ditempel ke chat atau prompt:

   ```powershell
   $env:Path = 'C:\Program Files\Go\bin;C:\Program Files\nodejs;' + $env:Path
   $env:ORACLE_DSN    = 'oracle://<pengguna>:<sandi>@<host>:<port>/<service>'   # bentuk DSN go-ora
   $env:ORACLE_SCHEMA = '<skema uji>'
   $env:IS_PEGA_PROD  = 'false'
   Set-Location 'D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\APP_RNM'
   ```

   Host, port, dan service diambil manusia dari entri alias DEV di `tnsnames.ora`.
3. **Bila instance belum ada saat sesi dimulai, sesi tetap jalan** — jalur B di §4. Yang tidak
   dikerjakan adalah menutup AC, bukan menulis kodenya.
4. **Isi §2 butir c–e.** Tiap baris bertanda `[USULAN]` disahkan dengan mengganti kata itu menjadi
   `[DIPUTUSKAN]`, atau ditulis ulang. Baris yang masih `[USULAN]` saat sesi dimulai **tidak boleh
   dijadikan dasar keputusan** oleh executor — kolomnya tetap dibuat, keputusannya dicatat terbuka.

---

## 2. KEPUTUSAN WORK OWNER YANG BERLAKU DI SESI INI

| | Keputusan | Keadaan |
| ---: | --- | --- |
| a | **Tiket 01 dianggap cukup membuka tiket 14.** Seam-nya sudah kompilasi dan lulus vet; tiga AC sisanya hanya menunggu instance dan ditutup di langkah A. Baris `Blocked by: 01` pada tiket 14 **tidak** menahan sesi ini | disahkan dengan memakai prompt ini |
| b | **Kolom uang bertipe `NUMBER(38,8)`** — ADR-U-0016 *(memperluas ADR-U-0003 yang menulis "tanpa presisi")*. `NUMBER(38,20)` yang ditemukan sesi sebelumnya ada di dokumen `dastin\` dan ADR-D-CNP-0003 — modul Claim Non-Prop, **bukan** korpus Pega *(korpus memuat nol `NUMBER(38,x)`)*. Penyatuan lintas-seri adalah butir Steering 5d, bukan urusan executor | berlaku |
| c | Kolom **share / persen / rate** (`PERCENT_SHARE`, `RATE`, `RETROCADED_SHARE`) bertipe `NUMBER(38,8)` juga, dan di Go dibawa sebagai `Ratio`, bukan `Money` (ADR-F-0004) | `[USULAN]` |
| d | Constraint `KOMITE_ID` dan `COVER_KEY`: dipasangi `REFERENCES T_WORK_CLAIM(ID)`, keduanya nullable (ADR-U-0027). *Tiket 14 §Blocker menyatakan tiket tidak `resolved` sebelum butir ini diputuskan* | `[USULAN]` |
| e | Lima kolom kehilangan rumah (tiket 14 §Blocker): `TANGGAL_RESPON`, `TANGGAL_REALISASI`, `TANGGAL_KONFIRMASI_BALIK` → kolom `T_GENERAL_CLAIM` *(tanggal proses klaim, bukan atribut polis)*; `TEAM_GROUP` → tidak disimpan, dibaca hidup dari master marketing officer lewat `MO_ID`; `BUSINESS_ID` → **tetap terbuka**, sumbernya belum diketahui | `[USULAN]` |
| f | **Letak DDL dan migrasi:** `APP_RNM\internal\repository\migrations\` — berkas `.sql` bernomor urut, ditanam lewat `//go:embed`, dijalankan `repository.Migrate` yang dipanggil `go run ./cmd/api -migrate` *(target `make migrate` yang sudah ada)*. Ini memenuhi baris `migrations/` pada tiket 14 **tanpa** menambah folder di luar struktur §1 brief induk. Jalur mundur: berkas `*_down.sql` berpasangan + flag `-migrate-down` `[usulan]` | berlaku |
| g | **Skema uji = hasil migrasi + fixture** *(brief induk §5)*, bukan DDL tangan. `internal\repository\skemauji\` dari tiket 01 diubah: `Pasang` memanggil `Migrate`, bukan DDL sendiri; `Bongkar` memanggil jalur mundur. Fixture tetap `UJI-*` | berlaku |
| h | **Data lama untuk menguji migrasi** = fixture sintetis di tabel tiruan `OS_AKSEPTASI_KLAIM_LIFE` *(55 kolom, `[terverifikasi]` dari `UpdateOsAkseptasiClaimLife_sql`)* yang dibuat oleh skema uji. Menjalankan migrasi terhadap data sungguhan adalah tiket **13** dan tunduk pada brief induk §7 | berlaku |

---

## 3. YANG TETAP `[terbuka]` — executor TIDAK menutupnya

Kolomnya dibuat; **isi, pembangkit, atau constraint-nya** menunggu pemilik. Ditulis apa adanya di
tiket dan di komentar kode dengan tanda `[terbuka]`, bukan ditebak:

| Butir | Pemilik | Sumber |
| --- | --- | --- |
| Constraint `REFERENCES` untuk `KOMITE_ID` / `COVER_KEY` — bila §2 d masih `[USULAN]` | DBA / work owner | tiket 14 §Blocker |
| Rumah `BUSINESS_ID`; dan keempat lainnya bila §2 e masih `[USULAN]` | work owner | tiket 14 §Blocker |
| Isi ketiga penunjuk polis `CASEID_POLICY`, `POLICY_NO`, `ENDORSMENT_NO` | work owner | tiket 14 |
| Generator nomor `CLM-xxxxxx` / `KMT-xxxxxx` | DBA / work owner | tiket 14 |
| Rumus `PREMIUM_SPREADED_NET` — dua cabang di rule yang sama | Product + UW, sebelum tiket 03 | tiket 14 |
| Nasib `CLMNO` pada kasus komite | work owner | tiket 14 |
| Daftar enum `LINI` lintas-lini | saat konteks Non-Life digarap | tiket 14 |
| Tipe/nullability `EDMSTATUS` (`NULL` atau `''`) | DBA — OQ-001 sisa | tiket 02 |

Akibat yang harus ditulis jujur: **selama butir pertama masih terbuka, tiket 14 berakhir `claimed`**
walau seluruh kodenya jadi dan hijau.

---

## 4. URUTAN SESI

**Langkah 0 — commit restrukturisasi, lalu titik tetap.** Bila `git status --porcelain` tidak
kosong, commit dulu — ini milik work owner, bukan milik tiket, dan pengecualian sadar atas aturan
"nol commit di luar lingkup tiket":

```
git add -A
git commit -m "fase-0: pindah aplikasi ke APP_RNM dan konversi frontend ke TypeScript (.tsx)" -m "Tiket: .scratch/claim-life/issues/01-kerangka-aplikasi-dan-seam-api.md"
git rev-parse HEAD        # <- titik tetap sesi ini
```

Lalu dari `APP_RNM\`: `go vet ./...`, `go test ./...`, `cd frontend; npm run typecheck; npm test` —
semua harus hijau **sebelum** apa pun diubah.

**Langkah A — hanya bila `ORACLE_DSN` terisi: tutup tiket 01.**
`go test -tags=db ./internal/... -count=1 -v`. Bila 11 test yang tadinya *SKIP* kini *PASS*: centang
tiga AC sisanya di tiket 01, `Status: resolved`, tambah baris tanggal di bab `## Implementasi`,
commit `claim-life: tiket 01 — penutupan AC yang menunggu Oracle`. Bila ada yang gagal: perbaiki
dulu, jangan lanjut ke 14 dengan seam yang merah. Tanpa instance: lewati langkah ini, catat.

**Langkah B — tiket 14.** `Status: claimed`. Kerjakan berurutan; tiap butir punya test-nya lebih
dulu (`/tdd`, pada seam brief induk §5):

1. **DDL** tujuh tabel — enam tabel klaim + `T_WORK_CLAIM` — persis pohon tiket 14 §"Bentuk baru":
   PK, shared PK (relasi 2), FK relasi 3·4·5·6 `ON DELETE CASCADE`, relasi 1·7·11 tanpa cascade
   *(hapus di Go)*, index pada setiap FK dan `KOMITE_ID`, sequence untuk `T_CLAIMLF_*` dan
   `DOCUMENT_CLAIM` (ADR-U-0006), `T_WORK_CLAIM.ID` teks berformat. Tipe: §2 b–c. Seluruh kolom
   nullable (ADR-U-0027). ⛔ Nol `T_CLAIMLF_POLICY`, `T_CLAIMLF_MARKETING`, `WORK_CLAIM_ID`,
   `KMT_NO`, `T_CLAIMLF_ADJUSTMENT_KOMITE` — test yang menemukannya **gagal**.
2. **Pelari migrasi** `repository.Migrate`: idempoten *(aman dijalankan ulang)*, berurutan, mencatat
   versi yang sudah jalan di satu tabel kecil `[usulan]`, punya jalur mundur yang **diuji**.
   ⛔ Menolak berjalan bila `IS_PEGA_PROD=true`. Nol `COMMIT` di teks SQL (ADR-U-0029).
3. **Skema uji** memakai migrasi itu (§2 g) + tabel tiruan `OS_AKSEPTASI_KLAIM_LIFE` (§2 h).
   Test tiket 01 harus tetap hijau di atas skema hasil migrasi.
4. **Models** untuk bentuk baru — `T_WORK_CLAIM`, spreading, spreading retro, dokumen — tanpa
   aturan dagang; `Money`/`Ratio` untuk setiap kolom uang/share; kode dan penanda tetap teks.
5. **Repository** yang diminta AC 14: menulis satu pohon klaim dalam **satu transaksi** — tabel
   relasional **dan** `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE` — nol JSON; membaca kembali sampai
   cicit (`_SPREADING_RETRO`); kaskade hapus sampai tingkat terdalam **diuji**.
6. **Migrasi data lama**: pembongkar baris `OS_AKSEPTASI_KLAIM_LIFE` *(satu baris per peserta, 14
   atribut polis berulang — tidak disalin ke klaim)* ke pohon baru; **seluruh** baris adjustment,
   bukan hanya keadaan terakhir (ADR-U-0011); uang **tanpa berubah satu digit** dan dibandingkan
   **tepat** (ADR-U-0021 / AC 53); tanggal teks → `DATE` tanpa geser zona; yang tak terurai
   **dilaporkan**; atribut polis yang berbeda antar baris peserta **dilaporkan**; `ACCEPTATION_DATE`
   hasil hardcode lama **dilaporkan** (AC 52). Laporan rekonsiliasi = keluaran fungsi, diuji.
7. **Verifikasi penuh sekali di akhir:** `go vet ./...`, `go vet -tags=db ./...`, `gofmt -l`,
   `go test ./...`, `go test -tags=db ./internal/...`, lalu di `frontend`: `npm run typecheck`,
   `npm test`, `npm run build`. Frontend tidak berubah di tiket ini — tetap harus hijau.

**Jalur B — tanpa instance Oracle:** butir 1, 2, 4, 5, 6 ditulis lengkap beserta test-nya; test
bertag `db` **melewati dengan pesan**, bukan lulus diam-diam. Tiket 14 berakhir `claimed` dengan
daftar AC yang menunggu instance. Jangan mengarang hasil eksekusi.

**Penutup — brief induk §6 butir 6–9:** bab `## Implementasi — <tanggal>` di tiket 14 *(AC ditutup
lawan total, butir §3 yang tetap terbuka, berkas)*, `/code-review` atas titik tetap langkah 0 dengan
path tiket, perbaiki temuan yang menyentuh AC, commit
`claim-life: tiket 14 — skema relasional klaim dan migrasi` + baris
`Tiket: .scratch/claim-life/issues/14-skema-relasional-klaim-dan-migrasi.md`. Satu commit per
tiket; tiket 01 dan 14 **tidak** digabung.

---

## 5. GAYA KODE — permintaan work owner, mengikat

Pembacanya **baru mengenal Go dan React**. Setiap berkas baru dibuka komentar kepala: *untuk apa
berkas ini*, *dibaca sesudah apa*, *aturan mana yang dijaganya*. Istilah — migrasi, idempoten,
sequence, shared PK, cascade, `embed`, transaksi — dijelaskan sekali dalam bahasa sederhana saat
pertama muncul. Satu konsep per komentar. Nama fungsi berbahasa Indonesia yang terbaca sebagai
kalimat (`JalankanMigrasi`, `BongkarBarisLama`). Berkas SQL diberi komentar `--` per tabel: tabel ini
untuk apa, siapa induknya. Perbarui `APP_RNM\README-BACA-DULU.md` bab 1 dan 5 bila ada folder atau
perintah baru — dalam bahasa yang sama.

---

## 6. YANG MEMERLUKAN PERSETUJUAN MANUSIA — berhenti dan tanya

Brief induk §7, ditambah untuk sesi ini: menjalankan `Migrate` terhadap skema **selain** skema uji
yang diberikan DBA · membaca tabel produksi mana pun *(termasuk `OS_AKSEPTASI_KLAIM_LIFE`
sungguhan)* · mengubah `T_WORK_CLAIM` di luar kolom yang disebut tiket 14 · menambah folder di luar
§1 brief induk · `git push`.

---

## 7. TELEMETRI EKSEKUSI — bab wajib di laporan akhir sesi

Pisahkan **terukur** dari **ditaksir**; angka token sejati tidak terlihat dari dalam sesi.

| Besaran | Cara ukur |
| --- | --- |
| AC ditutup / total — tiket 01 dan tiket 14 terpisah | hitung centang di berkas tiket |
| Butir §3 yang tetap terbuka | daftar |
| Objek DB dibuat oleh migrasi: tabel · sequence · index · constraint | dari log `Migrate` |
| Test: lulus · gagal · SKIP, per tag | keluaran `go test -v` dan vitest |
| Lama `go test -tags=db` | keluaran `ok … <detik>` |
| Berkas dibuat / diubah, baris berisi | `git diff --stat`, hitung baris bukan-kosong |
| SHA titik tetap dan SHA commit | `git rev-parse` |
| Sub-agen review: token dan panggilan alat | laporan harness |
| Token sesi utama · biaya · jam dinding | **⛔ tidak diukur** — nyatakan |

---

*Disusun 25 September 2026 sore, sesudah verifikasi log tiket 01, pemindahan ke `APP_RNM\`, dan
konversi frontend ke TypeScript. Brief induk: `PROMPT-IMPLEMENTASI-GO-REACT.md`.*
