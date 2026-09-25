# PROMPT — implementasi Go + React + Oracle dari tiket yang sudah ada

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia *(`disable-model-invocation:
> true`, `CLAUDE.md` §8)*, lalu tempel berkas ini **utuh** sebagai briefnya, ditutup satu baris:
> `SESI INI: <bounded-context> · tiket <NN[, NN…]>`. Skill itu sendiri pendek — ia memanggil
> `/tdd` pada seam yang **sudah disepakati** *(§5 brief ini adalah kesepakatannya)*, menjalankan
> test berkala dan penuh di akhir, lalu `/code-review`, lalu **commit**. Satu sesi mengerjakan
> satu tiket atau satu rantai pendek yang saling bergantung — bukan satu modul sekaligus.
>
> Ini **migrasi, bukan greenfield**. Kebenaran ada di korpus Pega dan di spec + tiket yang sudah
> ditulis darinya. Kode yang tidak dapat menunjuk tiketnya **tidak ditulis**.

---

## 0. KEADAAN AWAL — 25 September 2026

| | Keadaan |
| --- | --- |
| Kode | **nol** — 0 `.go`, 0 `go.mod`, 0 `package.json`. Seluruh struktur di §1 adalah target |
| Bahan | 20 dari 20 modul punya spec dan tiket: **365 tiket** di tiga pohon *(§3)*; **105 ADR** di tiga seri *(§4)* |
| Oracle pengembangan | **belum ada** — 8 tiket sudah berstatus `menunggu-instance` |
| Git | `OUTPUT_HASIL_RNM\` bukan repository git |

---

## 1. STRUKTUR — tidak boleh menyimpang

Akar aplikasi adalah **`OUTPUT_HASIL_RNM\` itu sendiri** *(`CLAUDE.md` §5: "semuanya menjadi anak
dari folder ini, tidak ada repo lain")*. Nama `my-web-app/` di bawah ini berarti folder itu.

```
my-web-app/                       = OUTPUT_HASIL_RNM\
├── cmd/api/main.go               entry point backend
├── internal/
│   ├── config/                   env var, koneksi Oracle, flag lingkungan
│   ├── handlers/                 HTTP controllers
│   ├── models/                   struct domain, tipe Money dan Ratio
│   ├── repository/               seluruh SQL — satu-satunya lapisan yang menyentuh Oracle
│   └── services/                 aturan dagang, transaksi, wewenang
├── pkg/utils/                    helper publik tanpa ketergantungan ke internal/
├── frontend/                     React via Vite, JSX
│   ├── public/
│   ├── src/{assets,components,hooks,pages,services,store}/
│   ├── src/App.jsx · src/main.jsx
│   ├── package.json · vite.config.js
├── go.mod · go.sum
└── Makefile                      run, build, test, db — kedua sisi
```

**Arah ketergantungan: `handlers → services → repository`.** Tidak terbalik, tidak memotong.
`pkg/` tidak boleh mengimpor `internal/`. Folder yang sudah ada — `.scratch\`, `discovery\`,
`docs\`, `dastin\`, `jefri\`, `scripts\` — **tidak disentuh** dan tidak dipindah.

---

## 2. FASE 0 — scaffold, satu sesi, sekali saja

Tiket: **tidak ada** — ini pekerjaan pendahuluan. Dikerjakan sekali; sesi berikutnya menemukan
semuanya sudah ada.

| Bagian | Isi |
| --- | --- |
| `go.mod` | modul `nusantarare`; Go ≥ 1.22 |
| `cmd/api/main.go` | baca config → buka koneksi → daftar handler → dengarkan. Nol aturan dagang di sini |
| `internal/config` | **hanya env var** *(ADR-U-0004)*: `ORACLE_DSN` · `ORACLE_SCHEMA` · `HTTP_ADDR` · `IS_PEGA_PROD` *(ADR-U-0005)* · alamat layanan luar. Nol literal host di kode mana pun |
| `internal/models` | `Money{Amount, Currency}` dan `Ratio{Value, Scale}` sebagai **tipe berbeda yang tidak dapat dijumlahkan** *(ADR-F-0004)*; desimal memakai **`cockroachdb/apd`** dengan konteks presisi **38** dinyatakan di satu tempat *(keputusan DECIDED-TEKNIS 18-09, `dastin\...\claim-non-prop\3-to-tickets\TICKETS.md` §T-keputusan; `shopspring/decimal` ditolak)*. ⛔ **`float64` tidak boleh muncul di jalur uang mana pun** — termasuk JSON keluar dan React |
| `internal/repository` | satu antarmuka per agregat; **nama skema eksplisit** di setiap query *(ADR-U-0033)*; transaksi dibuka-ditutup di `services`, **nol `COMMIT` di teks SQL** *(ADR-U-0029)* |
| `pkg/utils` | konversi teks↔desimal **satu fungsi** untuk seluruh batas procedure *(ADR-U-0034)*; format tanggal dua bentuk *(ADR-U-0022)* |
| `frontend/` | Vite + React; `services/api.js` satu klien HTTP *(axios)*; `store/` Zustand `[usulan]`; nol angka uang di-parse sebagai `Number` — dibawa sebagai string desimal |
| `Makefile` | `run-api` · `run-web` · `build` · `test` · `test-db` · `db-up` · `db-down` · `migrate` |

**Pustaka yang belum pernah diputuskan — dipakai sebagai `[usulan]`, dicatat di `go.mod` dengan
komentar, dan boleh diganti lewat keputusan tertulis:** driver Oracle `sijms/go-ora/v2` *(murni Go,
tanpa Instant Client)*; router `net/http` bawaan Go 1.22 *(nol ketergantungan)*.

**Oracle pengembangan:** `make db-up` menjalankan kontainer `gvenzl/oracle-free` `[usulan]` dan
membuat skema uji `POOLDATA`. Bila DBA menyediakan instance pengembangan, `ORACLE_DSN` menunjuk ke
sana dan `db-up` tidak dipakai. ⛔ **Tidak pernah** menunjuk instance produksi.

**`git init` di `OUTPUT_HASIL_RNM\` adalah bagian Fase 0, bukan usulan.** `/implement` menutup
setiap sesi dengan **commit**, dan `/code-review` bekerja atas `git diff <titik-tetap>...HEAD` —
keduanya tidak jalan tanpa repository. `.gitignore`: `node_modules/`, `dist/`, `*.env`,
`frontend/.vite/`. Commit pertama = scaffold kosong, pesan `fase-0: scaffold`. Korpus di folder
induk `D:\XML\RNM_BRD\` **tidak** ikut — repo hanya di `OUTPUT_HASIL_RNM\`.

**Tanda selesai Fase 0:** `make build` lulus, `make test` lulus dengan nol test *(kerangka ada)*,
`make run-api` menjawab `GET /healthz`, `make run-web` menampilkan satu halaman kosong. Nol aturan
dagang.

---

## 3. URUTAN KERJA — dua fase, dan satu gerbang di antaranya

### Fase 1 — modul yang **tidak menyentuh tabel akar bersama**. Mulai sekarang.

| Urutan | Bounded context | Tiket | Sumber |
| ---: | --- | --- | --- |
| 1 | Claim Life | `.scratch\claim-life\issues\` | paling terpisah — `CLAUDE.md` §9 |
| 2 | Life Master | `.scratch\master-product-name-life\` · `master-contract-retro-life\` | |
| 3 | Treaty Arrangement | `.scratch\treaty-contract-out\` | |
| 4 | Claim Non-Life | `.scratch\claim-prop\` · `claim-facin\` · `dastin\...\claim-non-prop\3-to-tickets\` | dastin: 29 tiket `selesai` = **DDL usulan tertulis, belum pernah dijalankan** — menjalankannya adalah pekerjaan Fase 1 |
| 5 | Komite | `.scratch\komite-claim-*\` · `dastin\...\komite-claim-non-prop\` | penegakan peran di `services` *(ADR-U-0030)* |
| 6 | Treaty Master | `dastin\...\treaty-in\5-tiket\` · `treaty-in-adjustment\5-tiket\` | |

Di dalam tiap konteks: ikuti baris **`Blocked by`** tiap tiket; nomor terkecil yang tidak
terblokir dikerjakan lebih dulu.

### ⛔ Gerbang — sebelum Fase 2

Tiga pihak membuat `T_WORK_POLIS` dan `T_GENERAL_POLIS` dengan rancangan kolom berbeda, dan tidak
saling tahu: PremiumList Life tiket **00**, NB Treaty In tiket **16**, dan Fac In *(spec 11 +
DDL draf, rekonsiliasi dinyatakan gugur — K-065)*. Kalau dua tiket pertama jalan duluan, tabel
lahir dengan bentuk Treaty In dan Fac In datang belakangan ke tabel yang sudah ada.

**Keputusan yang harus tertulis sebelum satu pun dari ketiganya dikerjakan** — milik work owner:

| | Keputusan | Usulan default, boleh disahkan dengan satu kata |
| ---: | --- | --- |
| a | Rancangan kolom `T_GENERAL_POLIS` dan `T_WORK_POLIS` | **gabungan**: kolom Treaty In *(79 + `REMARK`)* sebagai dasar, kolom Fac In yang belum ada ditambahkan, seluruhnya nullable *(ADR-U-0027)*; satu berkas migrasi, satu tiket pemilik |
| b | Kolom pembeda lini pada tabel akar | **ada dan wajib terisi** *(ADR-U-0026)*; K-069 yang mencabutnya berlaku untuk lingkup proyek Fac In, bukan untuk tabel bersama |
| c | Nama tabel anak untuk properti Pega yang sama | mengikuti nama yang **lebih dulu dibangun**; yang lain menyesuaikan saat gilirannya |

Keputusan ditulis ke `KEPUTUSAN-TABEL-AKAR-BERSAMA.md` di akar; tanpa berkas itu, Fase 2 tidak
dimulai. ⭐ Fase 1 cukup untuk berminggu-minggu — gerbang ini **tidak** menghentikan siapa pun.

### Fase 2 — sesudah gerbang

| Urutan | Bounded context | Tiket |
| ---: | --- | --- |
| 7 | Life Offer | PremiumList Life **00** dulu, lalu sisanya · Endorsement Life |
| 8 | Treaty Realisasi | NB Treaty In **16 → 17 → 18 → 19 → 20**, lalu sisanya · EDM Treaty In *(tiap tiket menyebut gate NB-nya)* |
| 9 | Facultative Inward | `jefri\OUTPUT FIX\05-tickets\` — NB 16 → RNW 8 → EDM 22 → Fac Out 14. **36 tiket `blocked`**; sebagian oleh tiket lain di rantainya, sebagian oleh bahan yang belum ada *(35 stored procedure, DDL, 604 dropdown)*. Yang tertahan bahan **tetap tertahan** sampai bahannya datang |

⚠️ Tidak ada tiket yang membuat **78 tabel flat** Fac In. DDL-nya ada sebagai draf di
`jefri\OUTPUT FIX\08-flat\`. Sebelum dijalankan, ia perlu **satu tiket pemilik** — dibuat di
`jefri\OUTPUT FIX\05-tickets\` mengikuti bentuk tiket yang ada, dan disetujui work owner.

---

## 4. ATURAN YANG MENGIKAT SETIAP TIKET

Sebelum menulis satu baris: baca `CLAUDE.md` §2–§4 dan §7–§10, `CONTEXT.md` *(glosarium)*,
spec modulnya, tiketnya, dan setiap ADR yang disebut tiket itu.

**Tiga seri ADR bernomor sama.** Di kode, komentar, dan pesan commit, sebut dengan awalan seri:
`ADR-U-nnnn` untuk `docs\adr\` · `ADR-D-<modul>-nnnn` untuk `dastin\...\docs\adr\` ·
`ADR-F-nnnn` untuk `jefri\OUTPUT FIX\adr\`. Nomor telanjang **tidak diterima**.

| Aturan | Sumber |
| --- | --- |
| Uang tidak pernah `float`; kolom uang `NUMBER(38,8)` | ADR-U-0003 · **ADR-U-0016** |
| `Money` dan `Ratio` tipe berbeda; skala rasio dari satu resolver | ADR-F-0004 |
| Konversi tipe **sekali saat masuk**; kode dan penanda tetap teks; uji nilai kolom langsung, bukan pulang-pergi saja | ADR-U-0022 |
| Pemecah dokumen punya **penampung medan tak dikenal**, wajib kosong sebelum selesai | ADR-U-0023 · U-0017 · U-0037 |
| Seluruh kolom nullable; wajib-isi ditegakkan di `services` | ADR-U-0027 |
| Transaksi dibuka-ditutup aplikasi; nol `COMMIT` di SQL; procedure ber-`ROLLBACK` dipanggil **terakhir** dan penandanya diperiksa | ADR-U-0029 |
| Wewenang diperiksa di `services`, sumber peran satu tabel, **nol nama orang** di kode | ADR-U-0030 |
| Hapus = penanda + nilai balik; pembatalan = generasi nol; nol `DELETE` di jalur pengguna | ADR-U-0031 · U-0024 |
| Nama skema eksplisit di setiap query | ADR-U-0033 |
| Hilir membaca tabel; penyusun JSON keluar **tidak dibangun** | ADR-U-0028 |
| Tetapan operasional dari tabel, bukan ditanam | ADR-U-0035 |
| Generasi = rantai baris; `NOURUT` memasang baris antar generasi; rumus selisih di `services`, nilai lama tidak dihitung ulang; tabel proyeksi bukan sumber | ADR-U-0018 · 0019 · 0020 · 0021 |

**Setiap klaim perilaku di kode menunjuk buktinya** — komentar satu baris: nomor tiket + AC, atau
`class / nama / tipe` rule Pega dari `<pxInsName>`. Nama rule saja **bukan** bukti.

---

## 5. UJI — apa yang membuat tiket selesai

`/tdd` menolak menulis test pada seam yang belum disepakati pengguna. **Tiga seam di bawah ini
adalah kesepakatan itu** — disetujui work owner lewat brief ini; executor tidak menambah seam
lain tanpa keputusan tertulis. Pemeriksaan **nilai kolom langsung** yang diwajibkan ADR-U-0022
dilakukan lewat **pembacaan publik seam `repository`** — bukan lewat SQL sampingan di dalam test
HTTP, yang oleh `/tdd` dianggap *side channel*.

| Seam | Dipakai untuk |
| --- | --- |
| **`repository`** terhadap **skema uji Oracle nyata** | setiap tiket penyimpanan; tabel dibuat oleh migrasi, diisi fixture, dibaca kembali |
| **HTTP** terhadap skema uji | tiket yang punya layar atau endpoint |
| **`services` murni** | rumus uang, tangga persetujuan, predikat — tanpa Oracle |

Wajib pada setiap tiket uang: nilai dibandingkan terhadap **angka dokumen produksi yang dikutip
spec** *(ADR-U-0021)*, bukan terhadap hasil kode sendiri. Wajib pada setiap tiket wewenang: **jalur
yang ditolak** diuji, bukan hanya yang berhasil. Pulang-pergi saja **tidak cukup** — kode
berawalan nol `"006"` yang kembali sebagai `"6"` lolos pulang-pergi dan memecahkan penggolong.

⛔ Fixture **tidak pernah** memuat nama orang, nomor polis nyata, atau potongan dokumen produksi.
Tiket `15-de-identifikasi-berkas-kasus` di Fac In menetapkan caranya; ikuti untuk semua modul.

---

## 6. SATU SESI, SATU TIKET — urutannya

1. Baca §4. Baca tiketnya. Baca spec modulnya di bagian yang dirujuk tiket.
2. Bila tiket berstatus `blocked` dan yang menahannya belum `resolved` — **berhenti**, laporkan,
   jangan mengerjakan tiket lain diam-diam.
3. Set `Status: claimed` pada berkas tiket.
4. Tulis test dulu dari acceptance criteria tiket, lalu kodenya, di lapisan yang benar.
5. Sebelum baris kode pertama: `git rev-parse HEAD` dicatat sebagai **titik tetap** sesi ini.
   Selama bekerja: `go vet ./...` dan test berkas tunggal berkala; `make test` dan `make test-db`
   penuh **sekali di akhir**; `make build` lulus.
6. Tambahkan di bawah tiket: `## Implementasi — <tanggal>`: berkas yang dibuat, test yang lulus,
   AC yang ditutup, dan **AC yang belum** beserta alasannya. Set `Status: resolved` hanya bila
   seluruh AC tertutup; bila tidak, tetap `claimed` dengan catatan.
7. `/code-review` dengan titik tetap dari langkah 5. Sumber spec-nya adalah **path tiket** — sebut
   di argumen dan di pesan commit, sebab tracker ini local-markdown tanpa nomor `#123`. Temuan
   yang menyentuh AC diperbaiki sebelum commit; yang menyentuh gaya dicatat di tiket.
8. **Commit** ke branch berjalan, satu commit per tiket, pesan:
   `<konteks>: tiket NN — <judul singkat>` + baris `Tiket: .scratch/<konteks>/issues/NN-*.md`.
   Nol commit yang memuat berkas di luar lingkup tiket.
9. Laporan akhir sesi: cacah AC tertutup lawan total, berkas disentuh, SHA commit, dan bab
   telemetri yang memisahkan yang terukur dari yang ditaksir.

---

## 7. YANG MEMERLUKAN PERSETUJUAN MANUSIA — berhenti dan tanya

Menjalankan migrasi pada instance selain kontainer lokal · perubahan authn/authz · akses
credential atau data teratur · perubahan sistem luar · `git push` ke remote mana pun *(repo ini
lokal; remote adalah keputusan tersendiri)* · membuat tiket baru untuk 78 tabel flat · **apa pun**
yang menyentuh `T_WORK_POLIS` / `T_GENERAL_POLIS` sebelum `KEPUTUSAN-TABEL-AKAR-BERSAMA.md` ada.

---

## 8. DISIPLIN

Korpus `D:\XML\RNM_BRD\` **READ-ONLY** — dibaca sebagai bukti, tidak pernah ditulis.
`D:\XML\nusantara-re\` ⛔ **terlarang** — tidak dibaca, tidak dibandingkan, tidak ditiru.
Nol secret, token, alamat host, data pelanggan, dump produksi — di kode, test, fixture, commit.
Berkas tersegel tidak disunting: `grilling-ronde-*`, `VERIFIKASI-*`, `KOREKSI-*-DIJALANKAN`,
`PROMPT-*`, seluruh `dastin\` dan `jefri\` kecuali baris `Status:` dan bab `## Implementasi` pada
tiketnya sendiri. Butir `[terbuka]` **tidak ditutup** oleh executor.

---

*Disusun 25 September 2026, sesudah audit menemukan nol kode, 365 tiket di tiga pohon, dan tiga
pembuat tabel akar yang sama.*

**SESI INI:** *(isi di sini)*
