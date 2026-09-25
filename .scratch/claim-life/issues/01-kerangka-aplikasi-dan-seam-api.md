# 01: Kerangka aplikasi + seam API + tracer "buka satu klaim Life"

**Status:** claimed

**Blocked by:** None (can start immediately)

## Hasil & nilai pengguna

Seorang pengguna dapat membuka satu klaim Life dan melihat baris-baris `AdjustmentList`-nya di
layar. Itu perilaku yang tipis, tetapi ia menembus **seluruh lapisan** — React → HTTP → handler →
service → repository → Oracle — sehingga sekaligus **menegakkan seam** yang dipakai dua belas tiket
sesudahnya.

Nilainya bukan fiturnya, melainkan **jalurnya**: setelah tiket ini selesai, setiap tiket berikutnya
menambah perilaku di jalur yang sudah terbukti hidup, bukan membangun jalur baru.

## Area codebase

Seluruhnya **di dalam `OUTPUT_HASIL_RNM/`** — hari ini belum ada satu pun berkas kode di sana.

| Bagian | Yang dibuat |
| --- | --- |
| akar | `go.mod`, `Makefile`, `.env.example` |
| `cmd/api` | entry point HTTP |
| `internal/config` | pembacaan env var + koneksi Oracle |
| `internal/models` | agregat klaim Life + koleksi baris adjustment |
| `internal/repository` | pembacaan klaim + baris dari `POOLDATA` |
| `internal/services` | perakitan agregat |
| `internal/handlers` | satu endpoint REST baca-saja |
| `frontend/` | Vite + React; satu halaman yang menampilkan klaim dan barisnya |

Arah dependency **`handlers` → `services` → `repository`** (`CLAUDE.md` §5). Tidak boleh terbalik,
tidak boleh memotong lapisan.

## Rule Pega sumber

Tiket ini **tidak meniru perilaku bisnis** — ia hanya membaca. Bentuk data diambil dari kolom yang
sudah terbaca:

| Rule | Identitas | Dipakai untuk |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` 55 nama kolom `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` terbaca langsung dari blok PL/SQL-nya |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` bentuk baris `AdjustmentList` dan kolom yang ditampilkan |

## ADR terkait

**ADR-0003** (uang non-float — presisi desimal eksak; kolom Oracle `NUMBER` tanpa presisi),
**ADR-0013** (koneksi database dari env var; **alamat layanan keluar TIDAK** — ia di-lookup runtime
dari `M_LINK_SERVICE`), **ADR-0011** (agregat memuat **koleksi baris**, bukan satu status).

## Acceptance criteria

- [ ] `GET` satu klaim Life mengembalikan klaim beserta **seluruh** baris `AdjustmentList`-nya,
      bukan hanya baris terakhir.
- [ ] Halaman React menampilkan klaim dan daftar barisnya, masing-masing dengan statusnya.
- [ ] Status ditampilkan sebagai kata (Outstanding / Aksep / Ditolak), **bukan** sebagai nama field
      `STS_REJECT` dan bukan sebagai angka. *(AC 26 spec)*
- [ ] Tidak ada nilai uang yang direpresentasikan sebagai *binary floating point* di lapisan mana
      pun maupun di JSON respons. *(AC 22 spec)*
- [ ] Tidak ada host, endpoint, atau kredensial sebagai literal di kode — semuanya env var.
- [ ] Ada satu test ujung-ke-ujung yang menggerakkan sistem lewat HTTP dan memeriksa hasilnya lewat
      HTTP, terhadap skema uji Oracle — bukan mock repository.
- [ ] `handlers` tidak memanggil `repository` langsung; `repository` tidak memuat business logic.

## Skema uji

`[terbuka]` **OQ-001** (pemilik **DBA**) — tidak ada DDL di korpus, sehingga **tipe dan presisi**
kolom tidak diketahui. Untuk tiket ini dipakai **skema uji provisional**: nama kolom dari rule di
atas, tipe numerik berpresisi longgar, tipe teks berpanjang longgar.

Ini **sah untuk skema uji** dan **tidak sah untuk DDL produksi maupun migrasi** — keduanya tetap
menunggu OQ-001 dan ditangani di tiket **13**.

## Perintah verifikasi

Toolchain belum ada; tiket ini yang membuatnya. Setelah selesai, ketiga perintah berikut **harus
berjalan hijau** dan dicatat di `Makefile`:

```
go build ./...
go test ./internal/...
cd frontend && npm test
```

Tambahkan target gabungan `make check` yang menjalankan ketiganya.

---

## Implementasi — 25 September 2026

**Status: `claimed`** — 4 dari 7 AC tertutup. Ketiga sisanya tertahan oleh satu hal yang sama:
belum ada instance Oracle yang dapat dibuati tabel.

### Berkas yang dibuat

Jendela cacah disebut: **total** = seluruh baris; **berisi** = baris bukan-kosong.

| Berkas | Total | Berisi |
| --- | ---: | ---: |
| `internal/models/klaimlife.go` | 211 | 196 |
| `internal/models/klaimlife_test.go` | 145 | 137 |
| `internal/repository/klaimlife.go` | 169 | 157 |
| `internal/repository/klaimlife_db_test.go` | 225 | 204 |
| `internal/repository/skemauji/skemauji.go` | 232 | 213 |
| `internal/services/klaimlife.go` | 68 | 59 |
| `internal/handlers/klaimlife.go` | 49 | 43 |
| `internal/handlers/klaimlife_db_test.go` | 170 | 154 |
| `frontend/src/pages/KlaimLife.jsx` | 104 | 95 |
| `frontend/src/services/api.test.js` | 36 | 30 |
| **jumlah** | **1.409** | **1.288** |

Diubah: `internal/handlers/handlers.go` (satu rute), `frontend/src/App.jsx`,
`frontend/src/services/api.js`, `frontend/package.json`, `Makefile` (target `check`).

⚠️ Tabel "Area codebase" tiket ini — `go.mod`, `Makefile`, `cmd/api`, kelima paket `internal/`,
`frontend/` — **sudah dikerjakan Fase 0** sebelum tiket ini dimulai, sesuai brief
`PROMPT-IMPLEMENTASI-GO-REACT.md` bab 2. Tiket ini menambahkan perilakunya, bukan kerangkanya.

### Test yang lulus

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` | lulus |
| `go vet -tags=db ./...` | lulus — test Oracle ikut kompilasi |
| `go test ./...` | lulus, **7 test** di `internal/models` |
| `npm test` | lulus, **5 test** di `src/services/api.test.js` |
| `go build ./...` · `go build -o bin/api` | lulus |
| `npm run build` | lulus, 87 modul |
| `gofmt -l` | nol berkas belum rapi |

`go test -tags=db` **melewati** dengan alasan `ORACLE_DSN belum dikonfigurasi` — bukan lulus
diam-diam.

### Code review — enam cacat diperbaiki sebelum commit

Dua poros (standar dan spec) atas titik tetap `2c48186`. Yang menyentuh AC diperbaiki:

1. ⛔ **Status tingkat klaim dikarang.** `Klaim.MarshalJSON` sempat menerjemahkan `STS_REJECT`
   header menjadi kata. spec.md, "Aturan yang mengikat": *"**`STS_REJECT` adalah status baris**,
   bukan status klaim"* dan *"Nilai `2` selalu berarti 'baris ini ditolak', tidak pernah 'klaim
   selesai'"*. Header hanyalah **cerminan** baris terakhir. Kata status tingkat klaim **dibuang**;
   kode mentahnya dibawa tanpa tafsir.
2. ⛔ **Angka mentah bocor ke layar.** `KlaimLife.jsx` sempat mencetak `(kode 9)` untuk status tak
   dikenal — AC-3 melarang menampilkan status sebagai angka. Dibuang dari layar; kode tetap ada di
   kontrak API bagi yang menelusuri.
3. **Relasi hilang dari skema uji.** Lisensi "provisional" melonggarkan tipe dan presisi, **bukan**
   relasi. Ditambahkan dua `FOREIGN KEY … ON DELETE CASCADE`, index `PREMIUM_LIST_DETAIL_ID`, dan
   index **UNIK** `KOMITE_ID` sebagaimana `STRUKTUR-TABEL-CLAIM-LIFE.md` menetapkannya.
4. **Jalur tulis fixture tidak dijaga NLS.** Jalur baca sudah aman lewat argumen NLS pada
   `TO_CHAR`; jalur tulis mem-bind teks desimal ke kolom `NUMBER` dan akan gagal `ORA-01722` pada
   sesi ber-pemisah koma. Ditambahkan `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,'`.
5. **Mata uang hilang saat jumlah kosong.** `CURRENCY` kini dibawa walau `CLAIM_AMOUNT` NULL —
   keduanya kolom terpisah.
6. **Kode mati dan ketidakseragaman.** `KlaimLife.Nama()` nol pemanggil → dibuang. `make check`
   sempat memakai `./internal/...` tanpa `go vet`, beda dari `make test` → disamakan ke `./...`
   plus `vet`.

Temuan yang **tidak** diubah, beserta alasannya, ada di bab Catatan butir 7.

### AC yang DITUTUP

- [x] **Status ditampilkan sebagai kata** (Outstanding / Aksep / Ditolak), bukan nama field dan
      bukan angka. Diuji termasuk jebakannya: nilai `"1"` berarti **diaksep** meski kolomnya
      bernama `STS_REJECT`, dan `"00"`/`"01"`/`"006"` **bukan** `"0"`/`"1"`/`"6"` — kode tetap
      teks (ADR-U-0022). Kode di luar ketiga nilai spec dilaporkan *Tidak diketahui*, tidak
      ditebak, dan tidak dicetak sebagai angka.
- [x] **Tidak ada nilai uang sebagai binary floating point** di lapisan mana pun maupun di
      kontrak API. Oracle menyerahkan `CLAIM_AMOUNT` sebagai **teks** lewat
      `TO_CHAR(…,'TM9','NLS_NUMERIC_CHARACTERS=''.,''')` — nilainya tidak pernah melewati
      `float64` di driver; JSON membawa `amount` sebagai teks; React menolak jumlah yang datang
      sebagai angka JSON.
- [x] **Tidak ada host, endpoint, atau kredensial sebagai literal.**
- [x] **`handlers` tidak memanggil `repository` langsung; `repository` tidak memuat business
      logic.** Perakitan agregat ada di `services`; seluruh SQL ada di `repository`.

### AC yang BELUM ditutup — dan sebabnya

- [ ] **`GET` satu klaim mengembalikan klaim beserta seluruh baris `AdjustmentList`-nya.**
      Kode dan test-nya ada dan kompilasi; **belum pernah dijalankan** terhadap Oracle.
- [ ] **Halaman React menampilkan klaim dan daftar barisnya.** Halaman ada dan `npm run build`
      lulus; **belum pernah dilihat** menampilkan data sungguhan, dan **belum punya test render**.
- [ ] **Satu test ujung-ke-ujung lewat HTTP terhadap skema uji Oracle, bukan mock.**
      `internal/handlers/klaimlife_db_test.go` ditulis persis begitu — sistem digerakkan lewat
      HTTP, hasilnya diperiksa lewat HTTP, nol mock, nol SQL sampingan. Ia **melewati** karena
      `ORACLE_DSN` kosong.

⛔ **Penghalang tunggal:** belum ada instance Oracle pengembangan yang dapat dibuati tabel.
Yang diperlukan: DSN, nama skema, dan hak `CREATE TABLE` di skema uji. Begitu ada,
`make test-db` menjalankan ketiga AC di atas tanpa perubahan kode.

### Catatan

1. ⚠️ **Skema uji provisional** di `internal/repository/skemauji/` memakai `NUMBER` tanpa presisi
   dan teks berpanjang longgar, sebagaimana bab "Skema uji" tiket ini mengizinkan. ⛔ Ia **bukan**
   DDL produksi dan **bukan** migrasi — itu milik tiket 14 — dan ia menolak berjalan bila
   `IS_PEGA_PROD` bernilai benar.
2. ⚠️ `[terbuka]` **Skala kolom uang berselisih di tiga tempat**: tiket ini menyebut *"kolom Oracle
   `NUMBER` tanpa presisi"*, brief bab 4 menyebut `NUMBER(38,8)`, korpus menulis `NUMBER(38,20)`
   **171 kali** lawan `NUMBER(38,8)` 20 kali. Presisinya sama (38), skalanya tidak. Tidak
   menghambat tiket ini; **mengikat di tiket 14**. → DBA / pemilik pekerjaan.
3. ⚠️ Tracer ini sengaja **tidak membawa medan bernama orang** (`NAME_OF_INSURED`,
   `POLICY_HOLDER`). AC tiket ini tidak memerlukannya, dan menahannya menjauhkan fixture dari
   data pribadi. Tiket berikutnya yang memerlukannya membawanya dengan aturan de-identifikasi.
4. ⚠️ Tiket ini meminta `npm test` hijau, sehingga **vitest** ditambahkan sebagai `[usulan]` —
   belum pernah diputuskan lewat ADR, boleh diganti. Versi dipatok `^2.1.9` sebab vitest 3+
   menuntut Vite 6.
5. ⚠️ Target `make check` yang diminta tiket ini **menambah satu target** di luar delapan yang
   disebut brief bab 2. Perbedaan itu dicatat, bukan didiamkan.
6. ⚠️ Rujukan ADR di berkas tiket ini masih **telanjang** (`ADR-0003`, `ADR-0011`, `ADR-0013`);
   di kode ditulis berawalan seri sesuai brief bab 4 — `ADR-U-0003`, `ADR-U-0011`, `ADR-U-0013`.
   Seluruh berkas tiket `claim-life` memuat **68 rujukan telanjang**; membetulkannya di luar
   lingkup tiket ini.
7. ⚠️ **Temuan review yang TIDAK diubah, beserta alasannya.** (a) Test murni ditaruh di
   `internal/models`, bukan `internal/services` — brief bab 5 menamai seam murni sebagai
   "`services` murni"; yang diuji memang perilaku tanpa Oracle, tetapi paketnya berbeda dari yang
   disebut brief. (b) Persiapan test `db` berulang di dua paket; menyatukannya menuntut paket
   bantu tersendiri. (c) Nilai uji dibandingkan terhadap angka fixture, bukan angka dokumen
   produksi (ADR-U-0021) — tiket ini pembacaan, bukan perhitungan, dan spec tidak mengutip angka
   produksi untuk baris ini. (d) Berkas test HTTP mengimpor `repository` untuk merakit tumpukan
   nyata; itu peran composition root seperti `main`, dan paket produksi `handlers` tetap nol impor
   `repository`.
8. ⚠️ Pemeriksa disiplin yang dipakai sesi ini (arah ketergantungan, float di jalur uang, awalan
   seri ADR, alamat harfiah, rahasia tertanam) adalah **instrumen sesi di scratchpad**, bukan
   bagian repo — ia diuji lebih dulu atas 7 contoh yang harus kena dan 5 yang harus bersih.
   Menjadikannya alat repo memerlukan tiket tersendiri, sebab `scripts/` termasuk folder yang
   brief bab 1 nyatakan tidak disentuh.

### Pemindahan dan konversi — 25 September 2026 (sore), sesi asisten atas permintaan work owner

Seluruh berkas di tabel "Berkas yang dibuat" kini berada di **`APP_RNM\`** *(dipindah dengan
`git mv`; riwayat utuh)*, dan frontend dikonversi ke **TypeScript**: `KlaimLife.jsx` → `KlaimLife.tsx`,
`api.js` → `api.ts`, `api.test.js` → `api.test.ts`, `App.jsx` / `main.jsx` → `.tsx`,
`store/index.js` → `index.ts`; ditambah `tsconfig.json`, `src/vite-env.d.ts`, dan komentar
bahasa sederhana untuk pembaca yang baru mengenal Go dan React (`APP_RNM\README-BACA-DULU.md`).
Perilaku tidak berubah. Dijalankan ulang dari `APP_RNM\` dan lulus: `go vet ./...`,
`go vet -tags=db ./...`, `gofmt -l` nol, `go build ./...`, `go test ./...` (7 test), `npm run typecheck`
(`tsc --noEmit`), `npm test` (5 test), `npm run build` (87 modul). Status tetap **`claimed`** — ketiga AC
yang menunggu Oracle tidak berubah keadaannya.
