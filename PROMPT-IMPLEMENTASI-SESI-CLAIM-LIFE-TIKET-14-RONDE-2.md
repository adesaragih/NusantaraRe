# PROMPT — sesi implementasi berikutnya: Claim Life tiket 14 **ronde 2** (perbaikan hasil verifikasi) + penutupan tiket 01

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya; brief sesi sebelumnya
> **`PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14.md`** tetap berlaku untuk hal yang tidak diubah di
> sini *(§2 f–h: letak migrasi, skema uji dari migrasi, tabel tiruan warisan)*.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> tiket 14 bab `## Implementasi — 25 September 2026` → `.scratch\claim-life\STRUKTUR-TABEL-CLAIM-LIFE.md`
> **seluruhnya** *(inilah sumber nama kolom; ronde 1 melewatkan lima kolom darinya)* → §3 berkas ini.
>
> **SESI INI:** claim-life · tiket **14 ronde 2** + tiket **01 (penutupan sisa AC, bila Oracle ada)**

---

## 0. KEADAAN AWAL — 26 September 2026

| | Keadaan |
| --- | --- |
| `HEAD` | `76dcda4` — tiket 14 *claimed* 25/53 AC; tiket 01 *claimed* 4/7. Working tree bersih |
| Uji yang lulus | `go vet`, `go vet -tags=db`, `gofmt` nol, `go build`; **34** test Go; **17** test bertag `db` SKIP dengan pesan; `tsc --noEmit`; **5** test JS; `vite build` |
| Migrasi | 16 berkas `.sql` (8 maju + 8 mundur), LF, tanpa BOM; 7 tabel · 5 sequence · 7 index · 7 PK · 6 FK · 4 cascade — **belum pernah dijalankan di Oracle mana pun** |
| Oracle pengembangan | **belum ada** — dua sesi berturut-turut berjalan di jalur B |
| Rantai alat | `go` dan `node` ada di mesin, tidak di PATH shell; `make` tidak ada; `.env` tidak dibaca kode *(env var disetel di terminal)* |

**Verifikasi independen 26 September 2026 atas commit `76dcda4`** menemukan hal-hal di §3. Ronde 2
ada untuk menutupnya — bukan untuk mengulang ronde 1.

---

## 1. PRASYARAT MANUSIA — sebelum sesi dimulai

1. **Oracle pengembangan** — sama seperti brief sesi sebelumnya §1: DSN, skema uji, hak
   `CREATE TABLE / SEQUENCE / INDEX / DROP`; env var di terminal, tidak pernah di berkas atau chat.
2. ⭐ **Tanyakan ke DBA versi Oracle dan parameter `COMPATIBLE` instance sasaran** — pengembangan
   *dan* produksi. Di bawah **12.2**, nama objek lebih dari **30 byte** ditolak (`ORA-00972`).
   Rancangan saat ini memuat tiga nama seperti itu *(§2 j)*. Jawabannya menentukan nama tabel, bukan
   hanya nama di skema uji.
3. **Putuskan §2.** Baris `[USULAN]` disahkan dengan mengganti kata itu menjadi `[DIPUTUSKAN]` atau
   ditulis ulang. Butir **d** sudah dua sesi tertunda dan sendirian menahan tiket 14 di `claimed`.

---

## 2. KEPUTUSAN WORK OWNER — yang lama masih menunggu, yang baru dari verifikasi

| | Keputusan | Keadaan |
| ---: | --- | --- |
| c | Kolom share / persen / rate bertipe `NUMBER(38,8)`; di Go `Ratio`, bukan `Money` | `[USULAN]` — belum dijawab sejak sesi lalu |
| d | `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)`, nullable. **Tanpa keputusan ini tiket 14 tidak pernah `resolved`** *(tiket 14 §Blocker)* | `[USULAN]` — belum dijawab sejak sesi lalu |
| e | Rumah lima kolom: tiga `TANGGAL_*` → `T_GENERAL_CLAIM`; `TEAM_GROUP` dibaca hidup lewat `MO_ID`; `BUSINESS_ID` tetap terbuka | `[USULAN]` — belum dijawab sejak sesi lalu |
| i | **`ACCEPT_STATUS` pada `T_WORK_CLAIM`.** `STRUKTUR-TABEL-CLAIM-LIFE.md` §"ACCEPT_STATUS — DIBUANG" *(keputusan work owner 2026-09-18, dengan bukti korpus)* membuangnya; diagram di tiket 14 dan spec §2b masih mencantumkannya; ronde 1 mengikuti diagram dan **membuat** kolom itu. Usulan: **ikuti STRUKTUR — kolom dibuang** dari `001_t_work_claim.sql`, dan diagram tiket 14 diberi catatan ralat | `[USULAN]` |
| j | **Nama lebih dari 30 byte** — tiga nama: tabel `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` (36) dan dua kolom STRUKTUR `RETROCESSION_VALUATION_BEGIN_DATE` (33) / `RETROCESSION_VALUATION_EXPIRED_DATE` (35). Usulan: kedua kolom memakai **nama korpus** `RETRO_VALUATION_BEGIN_DATE` / `RETRO_VALUATION_EXPIRED_DATE` *(STRUKTUR sendiri mencatat asalnya)* apa pun versinya; nama tabel **tetap** bila `COMPATIBLE ≥ 12.2`, dan menjadi `T_CLAIMLF_ADJ_SPREADING_RETRO` (29) bila di bawahnya | `[USULAN]` — jawaban §1 butir 2 menentukan |
| k | **Sumber `SpreadingList` / `RetroLifeList` data lama** (AC 46). Tabel datar warisan tidak memuatnya; sumbernya blob work object Pega. Ini **bukan** milik executor: perlu ekspor dari pemilik Pega *(lihat `DUMP-PROMPT-SELESAI\PERMINTAAN-KE-PEMILIK-EXPORT-PEGA.md` untuk bentuk permintaannya)*. Sampai ada, AC 46 `[terbuka]` | pemilik export Pega |
| l | **Aturan sunting berkas migrasi:** selama `T_MIGRASI` belum pernah ada di instance mana pun, berkas `001`–`008` **boleh disunting langsung**. Begitu satu instance pernah menjalankannya, setiap perubahan bentuk = **berkas baru bernomor lanjut**, tidak pernah menyunting yang lama | berlaku |

---

## 3. TEMUAN VERIFIKASI — wajib ditutup ronde ini, urut dari yang terberat

| # | Temuan | Letak | Yang dikerjakan |
| ---: | --- | --- | --- |
| 1 | **Lima kolom STRUKTUR tidak ada di DDL** `T_CLAIMLF_PREMIUMLIST_DETAIL`: `GROSS_VALUATION_BEGIN_DATE`, `GROSS_VALUATION_EXPIRED_DATE`, dua `RETRO*_VALUATION_*` *(nama per §2 j)*, `STNC_TREATY`. STRUKTUR menandainya "keputusan tiket 14"; ronde 1 melewatkannya tanpa catatan | `003_…sql` | tambahkan; test yang membandingkan **seluruh** daftar kolom DDL dengan tabel markdown STRUKTUR, tabel demi tabel — supaya selisih semacam ini tidak lolos lagi |
| 2 | **`ACCEPT_STATUS` ada di `T_WORK_CLAIM`** padahal STRUKTUR membuangnya | `001_…sql`, `models/pohonklaim.go` | ikuti §2 i |
| 3 | **Pembaca tabel warisan dari Oracle tidak ada.** `BongkarBarisLama` murni, tetapi tidak ada fungsi yang membaca `OS_AKSEPTASI_KLAIM_LIFE` menjadi `[]BarisLama`; `TestBongkarDataLamaDariTabelTiruan` menulis ke tabel tiruan lalu **mengabaikannya** dan membongkar slice di memori | `repository/pohonklaim.go`, `pohonklaim_db_test.go` | tulis `AmbilBarisLama(ctx, caseID)` — 55 kolom, **tanggal dan angka lewat `TO_CHAR` ber-argumen NLS** seperti `fmtDesimal` di `klaimlife.go`, sebab `BarisLama` bertipe teks; test membaca **kembali dari tabel**, baru membongkar |
| 4 | **Pembacaan sampai cicit hanya sebagian:** `AmbilSpreading` membaca `IDR`/`USD` saja — bukan `RETROCADED_SHARE`, `RATE`; `ambilRetro` membaca `AMOUNT` saja — bukan `PERCENT_SHARE`, `RATE`, `PREMIUM_SPREADED_GROSS/NET`, `COMMISION`, `OVR_COMM`. Dan **galat parse ditelan** (`if d, err := …; err == nil`) — bertentangan dengan "dilaporkan, bukan didiamkan" dan dengan `klaimlife.go` yang mengembalikan galat | `pohonklaim.go` baris 424–433, 483–487 | baca seluruh kolom uang/rasio lewat `TO_CHAR`; galat parse **dikembalikan**; test db memeriksa tiap kolom pulang-pergi digit demi digit |
| 5 | **Idempoten hanya untuk jalur mulus.** DDL Oracle menutup transaksinya sendiri; bila langkah gagal di pernyataan ke-2, pernyataan ke-1 sudah jadi tanpa catatan, dan `JalankanMigrasi` berikutnya mati di `ORA-00955`. Per-pernyataan tidak ada toleransi "sudah ada" | `migrasi.go` `JalankanMigrasi` | untuk `CREATE TABLE/INDEX/SEQUENCE` yang menjawab `ORA-00955` *(nama sudah dipakai)* atau `ORA-02264` *(nama constraint sudah dipakai)*: catat sebagai dilewati, lanjutkan, laporkan di `LaporanMigrasi`; test db: jalankan pernyataan pertama langkah 001 secara manual, lalu `JalankanMigrasi` harus **lulus** dan melaporkannya |
| 6 | **AC 11 belum dikerjakan:** kolom bank ada di tabel, `BarisAdjustment` tidak punya medannya, pembongkar tidak mengisinya | `models/klaimlife.go`, `migrasidata.go`, `pohonklaim.go` | tambah `NamaBank`, `IDBank`, `NomorRekening` ke `BarisAdjustment`; pembongkar memetakan `NAME_OF_BANK`, `IDBANK` → `ID_BANK`, `ACCOUNTNO` → `ACCOUNT_NO`; `Simpan`/`BarisLamaDari` membawanya; test murni + test db |
| 7 | **14 AC "tertulis tetapi belum punya test"** — nomor 5, 7, 9, 10, 23, 28, 31, 33, 39, 40, 41, 42, 45, 48 *(penomoran = urutan kotak `- [ ]` di tiket)* | `migrasi_test.go` | satu pernyataan test per AC, atas teks SQL yang ditanam; test #1 di atas menutup sebagian besar sekaligus |
| 8 | **Kotak AC di tiket tidak dicentang** — 25 AC dilaporkan tertutup, `[x]` = 0. Kepala bab Implementasi menulis "32 test", badan "34" | tiket 14 | centang `[x]` pada AC yang tertutup **di daftar aslinya** *(itu satu-satunya sunting yang diizinkan di badan tiket)*; ralat cacah test |
| 9 | Baris datar warisan ditulis **15 dari 55 kolom** oleh `Simpan`; 40 sisanya NULL. Hilir *(Arasapas)* membaca tabel itu | `pohonklaim.go` `Simpan` | **bukan** diperbaiki di sini — 14 atribut polis dibaca hidup dari tabel polis dan itu milik tiket 02/03. Tulis catatan `[terbuka — tiket 02/03]` di kode dan tiket: kolom mana yang masih NULL |
| 10 | Penjaga akhiran baris belum ada — BOM sudah dijaga, CR belum | `migrasi_test.go`, akar repo | test: nol byte CR di berkas migrasi; `.gitattributes` dengan `*.sql text eol=lf` dan `*.go text eol=lf` *(hanya dua pola itu — jangan `text=auto` global, ia mengubah ribuan berkas dokumen)* |
| 11 | `Hapus` di repository adalah `DELETE` fisik. Sah untuk **membuktikan kaskade** (AC 38), **tidak** untuk jalur pengguna — ADR-U-0031: hapus = penanda + nilai balik. Belum ada pemanggil di services/handlers, dan harus tetap begitu sampai tiket 15 | `pohonklaim.go` | komentar kepala `Hapus`: "⛔ hanya untuk test kaskade; jalur pengguna tiket 15 mengikuti ADR-U-0031"; test statik: nol pemanggil `Hapus` di luar `_test.go` |

Temuan yang **sudah benar** dan tidak perlu disentuh: BOM sudah dijaga; `PeriksaSQL` pada tiap
pernyataan; `{skema}` diganti saat jalan; `IS_PEGA_PROD` ditolak; jalur mundur membaca `T_MIGRASI`;
sequence punya pembaca (`nomorBerikut`); shared PK; kaskade hanya relasi 3·4·5·6; uang `NUMBER(38,8)`;
fixture `UJI-*`; nol nama orang.

---

## 4. URUTAN SESI

**Langkah 0 — titik tetap.** `git status --porcelain` harus kosong; `git rev-parse HEAD` dicatat.
Dari `APP_RNM\`: `go vet ./...`, `go test ./...`, `cd frontend; npm run typecheck; npm test` — hijau
sebelum apa pun diubah.

**Langkah A — hanya bila `ORACLE_DSN` terisi.**
1. Tutup tiket 01: `go test -tags=db ./internal/... -count=1 -v`; bila 17 SKIP menjadi PASS,
   centang tiga AC sisanya, `Status: resolved`, commit
   `claim-life: tiket 01 — penutupan AC yang menunggu Oracle`.
2. **Jalankan migrasi sungguhan** ke skema uji: `go run ./cmd/api -migrate`. Apa pun yang patah
   *(nama > 30 byte, tipe, sintaks)* diperbaiki **di ronde ini** dan dicatat di tiket. Lalu
   `-migrate` kedua kali harus melaporkan **0 dijalankan, 8 dilewati** — itu bukti idempoten.

**Langkah B — tiket 14 ronde 2.** Status tetap `claimed`. Urutan: §3 nomor 1 → 2 → 7 *(bentuk dulu,
test bentuknya sekalian)* → 5 → 3 → 4 → 6 → 10 → 11 → 8 → 9. Tiap nomor: test dulu (`/tdd` pada seam
brief induk §5), lalu kode, lalu `go vet` dan test berkas itu. Tanpa Oracle: test db ditulis dan
**SKIP dengan pesan**; jangan mengarang hasil.

**Langkah C — keputusan §2 yang sudah `[DIPUTUSKAN]`.** Terapkan **hanya** yang disahkan:
d → tambah constraint di `001`/`004` *(masih boleh disunting, §2 l)* dan cabut catatan `[terbuka]`-nya;
c → ubah tipe kolom share/rate; e → tambah tiga kolom tanggal ke `002`; i, j → sesuai isinya.
Yang masih `[USULAN]` **tidak disentuh** dan tetap tercatat terbuka.

**Verifikasi penuh sekali di akhir:** `go vet ./...`, `go vet -tags=db ./...`, `gofmt -l`,
`go test ./...`, `go test -tags=db ./internal/...`, `cd frontend; npm run typecheck; npm test; npm run build`.

**Penutup — brief induk §6 butir 6–9:** bab `## Implementasi — ronde 2, <tanggal>` di tiket 14
*(AC ditutup lawan total; temuan §3 yang ditutup lawan yang belum; butir terbuka)*, `/code-review`
atas titik tetap langkah 0 dengan path tiket, commit
`claim-life: tiket 14 ronde 2 — perbaikan hasil verifikasi` + baris `Tiket: .scratch/claim-life/issues/14-…md`.
Tiket 01 dan 14 **tidak** digabung dalam satu commit.

**Tanda `resolved` untuk tiket 14** — dan hanya bila **semuanya**: §2 d `[DIPUTUSKAN]` dan
diterapkan; migrasi pernah berjalan sungguhan dan idempoten (langkah A-2); seluruh test db PASS;
tiap AC bertanda `[terbuka]` di tiket tinggal yang pemiliknya memang di luar executor. Bila satu saja
tidak terpenuhi: `claimed`, dengan alasan yang mana.

---

## 5. SESUDAH INI — tiket 02, prasyaratnya disiapkan sekarang

Tiket berikutnya menurut tepi pemblokir: **02 Register klaim + penomoran** *(blocked by 01 dan 14)*.
Ia menuntut objek yang **belum ada di skema uji** dan hanya DBA yang punya bentuk aslinya:

| Objek | Dipakai untuk | Yang diminta ke DBA |
| --- | --- | --- |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` + tabel `GENERATE_SEQUENCE_NUMBER` | nomor klaim — ADR-U-0006: **jangan replikasi logikanya** | sumber procedure dan DDL tabelnya, atau skema uji yang sudah memuat keduanya |
| `POOLDATA.KODE_PRODUKSI` (`TYPE='LIFE'`) | prefix nomor lewat lookup | DDL + isi baris `LIFE` |
| `POOLDATA.TANGGAL_CLOSING` | periode nomor, aturan cutover | DDL + contoh isi |
| `M_LIFE_PREMIUM_DETAIL` *(modul PremiumList Life)* | pencarian peserta ber-`EDMSTATUS` | DDL, dan jawaban OQ-001 sisa: `EDMSTATUS` NB `NULL` atau `''` |

Minta sekarang, supaya sesi tiket 02 tidak lagi berjalan di jalur B.

---

## 6. GAYA KODE — tetap mengikat

Pembaca **baru mengenal Go dan React**: komentar kepala tiap berkas *(untuk apa, dibaca sesudah apa,
aturan mana yang dijaga)*; istilah dijelaskan sekali; satu konsep per komentar; nama berbahasa
Indonesia yang terbaca sebagai kalimat. Perbarui `APP_RNM\README-BACA-DULU.md` bila ada perintah atau
berkas baru.

---

## 7. YANG MEMERLUKAN PERSETUJUAN MANUSIA — berhenti dan tanya

Brief induk §7, ditambah: `-migrate` terhadap skema **selain** skema uji dari DBA · membaca tabel
produksi mana pun · mengganti nama tabel di luar §2 j · menambah folder di luar §1 brief induk ·
menyentuh `dastin\` / `jefri\` · `git push`.

---

## 8. TELEMETRI EKSEKUSI — bab wajib di laporan akhir

| Besaran | Cara ukur |
| --- | --- |
| Temuan §3 ditutup / 11; AC tiket 14 ditutup / 53; AC tiket 01 ditutup / 7 | hitung di tiket |
| `-migrate` sungguhan: dijalankan · dilewati · pernyataan, run pertama dan kedua | log `jalankanMigrasi` |
| Objek DB nyata sesudah migrasi: `SELECT COUNT(*) FROM ALL_OBJECTS WHERE OWNER=:1` per `OBJECT_TYPE` | query sekali, dicatat |
| Test: PASS · FAIL · SKIP per tag; lama `go test -tags=db` | keluaran `go test -v` |
| Berkas dibuat / diubah, baris berisi; SHA titik tetap dan commit | `git diff --stat`, `git rev-parse` |
| Sub-agen review: token, panggilan alat | laporan harness |
| Token sesi utama · biaya · jam dinding | **⛔ tidak diukur** — nyatakan |

---

*Disusun 26 September 2026 dari verifikasi independen commit `76dcda4`: 34 test dijalankan ulang,
16 berkas migrasi dibaca utuh, kolom DDL dibandingkan tabel demi tabel dengan
`STRUKTUR-TABEL-CLAIM-LIFE.md`.*
