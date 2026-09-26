# PROMPT — sesi implementasi berikutnya: Claim Life tiket 14 **ronde 6** — dua keputusan tipe warisan, sisa nama, dan penutupan bila Oracle ada

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya. Brief ronde 5
> **`PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14-RONDE-5.md`** tetap berlaku untuk hal yang tidak
> diubah di sini *(§4 Langkah A dan D, §5, §6, §7, §8)*, dan brief ronde 4 **§9** tetap menjadi
> rujukan katalog.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> tiket 14 bab `## Implementasi — ronde 5, 26 September 2026` → §2–§3 berkas ini. Bila §2 s atau
> w2 `[DIPUTUSKAN]`: `.scratch\claim-life\TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md` dan tiket 04
> baris `[data DBA] STS_REJECT`.
>
> **SESI INI:** claim-life · tiket **14 ronde 6** + tiket **01 (penutupan, bila G1)** + **persiapan
> tiket 02 (bila G3)**. Tanpa G1 dan tanpa keputusan baru, sesi ini **pendek**: §3-1 dan §3-2 saja.

---

## 0. KEADAAN AWAL — 26 September 2026 malam

| | Keadaan |
| --- | --- |
| `HEAD` | `f857352` — tiket 14 *claimed* **40/53**, daftar terbuka 13 = kotak `[ ]`; tiket 01 *claimed* 4/7; tiket 02 *ready-for-agent*, blocker o (AC 2, 3, 10; AC 3 harus mencakup perakitan format). Working tree bersih |
| Uji yang lulus | `go vet`, `go vet -tags=db`, `gofmt` nol, `go build`; **91** test Go PASS; **20** test bertag `db` SKIP dengan pesan; `tsc --noEmit`; **5** test JS; `vite build` 87 modul |
| Diterapkan ronde 5 | **x** pra-terbang bentuk *(hanya langkah yang belum tercatat; `CREATE TABLE` yang tak terurai menolak)*; **w** `WPC DATE` di `003`; **v1** `T_CLAIMLF_DOCUMENT` *(`007`, `SEQ_`, kode, test, STRUKTUR)*; peta tipe warisan **62 kolom dari katalog**, tabel tiruan 62; flag **`-migrate-down`** dengan pagar di `internal/config` *(satu pagar untuk test db dan flag)*; `KolomCreateTable` diangkat dari test ke kode |
| Migrasi | 16 berkas, 7 tabel + `T_MIGRASI`, 8 FK, 5 sequence — **belum pernah dijalankan di Oracle mana pun**, enam sesi |
| Ditahan executor, menunggu Anda | **`STS_REJECT`** *(ADR-U-0022 lawan katalog)*, **arti `CLAIM_RETRO`** *(w2)*, sisa nama `DOCUMENT_CLAIM` di `spec.md`, tiket 03, tiket 15, `revisi-penyimpanan-json-dibuang.md` |

**Verifikasi independen 26 September 2026 malam atas `e841f70` + `f857352`:** seluruh angka laporan
ronde 5 tereproduksi *(91 · 20 · 5 · 40/53 · 13 · 27 berkas +1.293/−99 · 2 baru · 2 ganti nama)*;
pra-terbang, pagar `-migrate-down`, peta 62 kolom, dan rename v1 dibaca di diff; nol cacat kode. Yang
belum tepat ada di §3, semuanya kata. Dua keputusan yang ditahan executor **diberi bukti** di §2
supaya tidak diputuskan dengan menebak.

---

## 1. GERBANG

| Gerbang | Yang membukanya | Keadaan |
| --- | --- | --- |
| **G0 tanpa syarat** | — | terbuka: §3-1, §3-2 |
| **G1 Oracle** | user **kosong** baru + `ORACLE_DSN` + `ORACLE_SCHEMA` + `ORACLE_SKEMA_UJI=true` | tertutup — DBA belum membuat user |
| **G2 keputusan tiket 14** | **s**, **w2** *(keduanya bisa diputuskan sekarang dengan bukti §2)*; y, j bila ada | tertutup |
| **G3 keputusan tiket 02** | **o1, o2, o3** | tertutup — hanya oleh keputusan; objek DBA sudah terbaca |

---

## 2. KEPUTUSAN WORK OWNER — dengan bukti

| | Keputusan | Bukti | Keadaan |
| ---: | --- | --- | --- |
| s | **`STS_REJECT` di jalur tulis tabel warisan.** Ronde 5 membingkainya sebagai *"ADR-U-0022 lawan katalog"*. Bingkainya lebih sempit dari itu: ADR-U-0022 melindungi kode **berawalan nol** seperti `"006"`; `STS_REJECT` adalah **angka satu digit menurut aksi** *(0 · 1 · 2, tiket 04)*, dan tiket 04 **sudah** mencatat `[data DBA] STS_REJECT warisan bertipe NUMBER(38)`. Mem-bind teks `"1"` ke kolom `NUMBER` sah di Oracle. Usulan: **(s1)** skema baru **tetap teks** *(ADR-U-0022 utuh)*; jalur tulis ke tabel warisan **menolak** nilai yang bukan digit *(pagar kompatibilitas, bukan konversi — galat terang lebih baik daripada `ORA-01722` di lapangan)*; test murni untuk pagar itu | agregat DEV *(hanya cacah, nol baris)*: `STS_REJECT` bernilai **0** (8.044) · **1** (5.227) · **2** (346) · **4** (77) · **NULL** (0). Nol nilai berawalan nol, nol nilai non-angka | `[USULAN]` — rekomendasi **s1** |
| s′ | ⭐ **Kode `4` ada di data DEV (77 baris) dan tidak dikenal siapa pun**: tiket 04 mengenumerasi 0/1/2, `StatusBarisDariKode` memetakan `4` → *"tidak diketahui"* *(sudah benar: tidak menebak)*. Artinya milik **tiket 04 / work owner**; korpus tidak memuat penetapan nilai 4 dalam pola yang dapat saya temukan | agregat DEV | `[terbuka — tiket 04]`, bukan ronde ini |
| w2 | **Arti `CLAIM_RETRO`.** Ronde 5 menahannya: tipe `NUMBER` diketahui, arti tidak. Bukti menunjuk **uang** *(bagian retro dari jumlah klaim)*, bukan perbandingan: nol nilai di rentang 0–1 atau 1–100, 4.778 nilai di atas 100, 475 berdesimal, dan pada 8.755 baris `CLAIM_RETRO < CLAIM_AMOUNT` *(782 sama, 220 lebih besar — anomali data, dicatat)*; label Pega di `AdjustmentDetail_Section` berbunyi *"Claim Retro"* di samping label uang lain. Spec sendiri menaruhnya di **header** *(`PremiumListSummary`)*. Usulan: `002` `CLAIM_RETRO` → **`NUMBER(38,8)`**, di Go **`Money`** dengan mata uang header; `BarisLamaDari` menulis nilai header ke tiap baris datar *(sekarang kosong)*; STRUKTUR diberi ralat bertanggal; `[terbuka]` di STRUKTUR dicabut | agregat DEV + label korpus | `[USULAN]` — rekomendasi **disahkan** |
| y, j, o1–o3, v2 | tetap sebagaimana brief ronde 5 §2 | | `[USULAN]` / `[terbuka]` |
| k, l | tetap | | berlaku |

---

## 3. TEMUAN VERIFIKASI RONDE 5 — semuanya kata

| # | Temuan | Letak | Yang dikerjakan |
| ---: | --- | --- | --- |
| 1 | ⚠️ Bab ronde 5 tabel *"Yang dijalankan"* menulis **89 test**; telemetri di bab yang sama dan hasil nyata **91** *(dua test paket `config` lahir sesudah tabel itu ditulis)* | tiket 14 bab ronde 5 | satu baris ralat di bab ronde 6 |
| 2 | ⚠️ **Nama lama tertinggal di komentar** untuk tabel **baru**: `007_t_claimlf_document.sql` baris 1 *(judul kepala masih `DOCUMENT_CLAIM`)*, `008_sequences.sql` baris 15 *(`T_CLAIMLF_*` dan `DOCUMENT_CLAIM`)*, `pohonklaim.go` baris 154, 403, 405. Yang **benar** tetap `DOCUMENT_CLAIM`: nama kelas Pega dan rujukan ke tabel warisan *(`models/pohonklaim.go` 111, `migrasi.go` 179, `007` baris 19)* | tiga berkas | ganti hanya rujukan ke tabel **baru**; test statik kecil: di `migrations/*.sql` dan `pohonklaim.go`, kata `DOCUMENT_CLAIM` hanya boleh muncul dalam frasa `Int-DOCUMENT_CLAIM` atau bersama kata `warisan` |
| 3 | ⚠️ Sisa nama lama sebagai tabel **baru** di `spec.md`, `issues/03`, `issues/15`, `revisi-penyimpanan-json-dibuang.md` | berkas milik work owner | **bukan executor**: tiket lain dan spec tidak disunting executor. Usulan untuk Anda: satu baris ralat bertanggal di kepala tiap berkas, bukan ganti nama massal |
| 4 | ℹ️ Pra-terbang membandingkan **nama** kolom, bukan tipe — sengaja, tertulis di komentarnya. Dicatat saja | `migrasi.go` | tidak diubah |

Yang **sudah benar**: pra-terbang hanya untuk langkah yang belum tercatat *(kolom audit DBA pada
tabel lama tidak menggagalkan migrasi berikutnya)*; `CREATE TABLE` yang tak terurai ditolak; pagar
tunggal di `internal/config` dipakai test db **dan** `-migrate-down`; `TestPetaTipeCocokDenganKatalog`
membaca dokumen `[data DBA]` baris demi baris dan mengunci 62; fixture `CLAIM_RETRO` angka; klaim
kosong soal "mengunci kolom angka" sudah diralat sendiri; tiket 02 `[dugaan]` naik ke
`[terverifikasi]`; STRUKTUR bab `T_CLAIMLF_DOCUMENT` dan ralat tiga tipe; tolakan atas *"003 tanpa
jalur naik"* benar.

---

## 4. URUTAN SESI

**Langkah 0 — titik tetap.** `git add PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14-RONDE-6.md`,
commit `docs: brief sesi tiket 14 ronde 6`. `git status --porcelain` kosong; `git rev-parse HEAD`
dicatat; uji tanpa Oracle hijau.

**Langkah B — G0.** §3-1 → §3-2. Nol perubahan perilaku; `go test ./...` bertambah satu *(test
statik §3-2)*.

**Langkah C — G2 bila terbuka.** **s1** → pagar digit di jalur tulis tabel warisan *(`BarisLamaDari`
atau `NilaiBarisLama`, pilih tempat yang paling dekat dengan bind)*, galat menyebut kolom dan
nilainya, test murni; `models` **tidak disentuh** *(kode tetap teks)*. **w2** → `002` `CLAIM_RETRO
NUMBER(38,8)`, `models.Klaim` mendapat `ClaimRetro models.Money`, `AmbilHeader` membacanya lewat
`fmtDesimal`, `BarisLamaDari` menulis nilai header ke `CLAIM_RETRO` tiap baris datar, STRUKTUR
ralat bertanggal, `TestNilaiKolomAngkaSelaluAngkaAtauKosong` kini **benar-benar** menguji
`CLAIM_RETRO` *(kepala test diperbarui: klaim yang tadinya kosong kini terisi)*. Keduanya **hanya**
selama `T_MIGRASI` belum ada di instance mana pun *(§2 l)*.

**Langkah A — G1 bila terbuka.** Persis brief ronde 5 §4 Langkah A: `ALL_OBJECTS` skema harus
**0** dulu; test db → `-migrate` dua kali → cacah objek *(8 · 5 · 15)* → `-migrate-down` sekali
*(bukti pintu masuk manusia, lalu `-migrate` lagi)* → tiket 01.

**Langkah D — G3 bila terbuka.** Persis brief ronde 5 §4 Langkah D.

**Verifikasi penuh, penutup, tanda `resolved`** — persis brief ronde 5 §4. Tanpa G1: `claimed`,
ditulis di awal laporan.

---

## 5. YANG DIMINTA KE DBA

Tetap brief ronde 5 §5. Yang paling menentukan masih **butir 1**: user kosong. Enam sesi telah
berakhir `claimed` karena itu saja.

---

## 6–8. GAYA KODE · PERSETUJUAN MANUSIA · TELEMETRI

Persis brief ronde 5 §6–§8. Tambahan §7: **membaca agregat tabel warisan** *(cacah, min, max)*
boleh untuk keputusan, **baris tidak pernah**; tulis di laporan agregat mana yang dibaca.

---

*Disusun 26 September 2026 malam dari verifikasi independen `e841f70` + `f857352`: 91 test
dijalankan ulang, diff 27 berkas dibaca utuh, agregat `STS_REJECT` dan `CLAIM_RETRO` dibaca dari
instance pengembangan (nol baris), label `CLAIM_RETRO` dicari di korpus Pega.*
