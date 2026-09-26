# PROMPT — sesi implementasi berikutnya: Claim Life tiket 14 **ronde 4** (penutupan) + persiapan tiket 02 sesudah keputusan o

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya. Brief ronde 3
> **`PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14-RONDE-3.md`** tetap berlaku untuk hal yang tidak
> diubah di sini *(§1-2 query katalog DBA, §2 k–l, §6, §7)*.
>
> ⛔ **BACA §1 DULU.** Sesi ini **tidak dimulai** bila gerbangnya belum terbuka. Tanpa Oracle dan
> tanpa satu pun keputusan baru, yang tersisa untuk executor hanya tiga sunting kecil (§3), dan itu
> tidak sepadan dengan satu sesi.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> tiket 14 bab `## Implementasi — ronde 3, 26 September 2026` → tiket 02 bab `## Implementasi` →
> `docs\adr\0006-…md` dan `0016-…md` **seluruhnya** → §2–§3 berkas ini.
>
> **SESI INI:** claim-life · tiket **14 ronde 4** + tiket **01 (penutupan)** + **persiapan tiket 02**
> *(ADR pengganti dan teks AC baru — hanya bila §2 o1–o3 `[DIPUTUSKAN]`)*

---

## 0. KEADAAN AWAL — 26 September 2026

| | Keadaan |
| --- | --- |
| `HEAD` | `7b832e6` — tiket 14 *claimed* 38/53 AC; tiket 01 *claimed* 4/7; tiket 02 *ready-for-agent* dengan empat blocker baru di bab `## Implementasi`-nya. Working tree bersih kecuali `APP_RNM\PANDUAN-MENJALANKAN.txt` *(panduan menjalankan untuk manusia, belum di-commit; commit bersama brief ini di Langkah 0)* |
| Uji yang lulus | `go vet`, `go vet -tags=db`, `gofmt` nol, `go build`; **77** test Go PASS; **19** test bertag `db` SKIP dengan pesan; `tsc --noEmit`; **5** test JS; `vite build` 87 modul |
| Migrasi | 16 berkas `.sql`, LF; **28** kolom `NUMBER(38,8)` + `AGE NUMBER(5)` + `DOCUMENT_CLAIM.ID NUMBER(19)`, **nol `NUMBER` polos** *(dikunci test)*; 19 pernyataan `CREATE` semuanya terbaca namanya; **belum pernah dijalankan di Oracle mana pun** — **empat** sesi berturut-turut di jalur B |
| Penjaga migrasi | `sudahAda` hanya `ORA-00955`; `CREATE` yang dilewati dibuktikan lewat `SYS.ALL_OBJECTS`; tabrakan nama constraint **menggagalkan** migrasi *(test db, SKIP tanpa Oracle)* |
| Keputusan diterapkan | **c** tuntas (ronde 3), **i**, **j** (ronde 2); **o** dicatat di tiket 02 |
| Keputusan menunggu | **d** *(empat sesi)*, **e′**, **m**, **n**; baru: **p**, **t**, **o1–o3** |

**Verifikasi independen 26 September 2026 atas commit `7b832e6`:** seluruh angka laporan ronde 3
tereproduksi *(77 · 19 · 5 · 38/53 · 11 berkas +516/−23 · 28/30/19/33 · CR = 0 di HEAD, index,
working tree)*; tujuh temuan ronde 2 memang tertutup, satu terhalang data DBA, satu catatan; kutipan
ADR-U-0016 Akibat 2 dan ADR-U-0006 benar kata demi kata. Yang **belum tepat** ada di §3, semuanya
kecil dan tentang **kata**, bukan kode.

---

## 1. GERBANG SESI — salah satu harus terbuka, kalau tidak jangan mulai

| Gerbang | Yang membukanya | Pemilik |
| --- | --- | --- |
| **G1 Oracle** | `ORACLE_DSN` + `ORACLE_SCHEMA` skema uji, hak `CREATE TABLE / SEQUENCE / INDEX / DROP`, disetel di terminal. ⛔ Skema uji = skema **kosong yang dibuat khusus** di instance **DEV**, **bukan** `POOLDATA` DEV dan bukan skema apa pun yang memuat tabel warisan sungguhan — test db **menghapus** `OS_AKSEPTASI_KLAIM_LIFE` di skema yang ditunjuknya *(§3-5)*. Sebelum §3-5 selesai, `go test -tags=db` hanya boleh dijalankan sesudah `ORACLE_SCHEMA` dibaca ulang dengan mata sendiri | DBA |
| **G2 Keputusan tiket 14** | minimal **d** disahkan *(satu-satunya yang menahan `resolved`)*; e′, m, n, p, t ikut bila ada | work owner |
| **G3 Keputusan tiket 02** | **o1, o2, o3** disahkan **dan** empat objek DBA §5 sudah diserahkan | work owner + DBA |

Bila **hanya G2** terbuka tanpa G1: sesi tetap sah, tetapi berakhir `claimed` lagi — tulis begitu
di awal laporan, bukan di akhir.

Prasyarat lain yang **tidak** membuka gerbang tetapi tetap diminta: daftar tipe kolom
`OS_AKSEPTASI_KLAIM_LIFE` *(brief ronde 3 §1-2; menutup §3-7 ronde 3)* dan versi Oracle +
`COMPATIBLE`.

---

## 2. KEPUTUSAN WORK OWNER

| | Keputusan | Keadaan |
| ---: | --- | --- |
| d | `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)`, nullable. **Tanpa ini tiket 14 tidak pernah `resolved`.** Bila jawabannya "tidak dipasang, integritas di Go", itu **juga** keputusan yang sah dan menutup AC 50 — yang tidak sah hanya membiarkannya `[USULAN]` untuk kali kelima | `[USULAN]` — empat sesi |
| e′ | Tiga `TANGGAL_*` bukan milik Claim Life; dibaca hidup lewat rantai polis → offer; AC 8 tiket 14 diberi `[terbuka — PremiumList Life]` *(brief ronde 3 §2 e′, bukti korpus diverifikasi dua kali)* | `[USULAN]` |
| m | Blok ralat bertanggal di STRUKTUR bab `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`: nama fisik `T_CLAIMLF_ADJ_SPREADING_RETRO` | `[USULAN]` |
| n | `git mv` berkas `006` (+`_down`) menjadi `006_t_claimlf_adj_spreading_retro.sql` — **hanya sebelum** migrasi sungguhan pertama | `[USULAN]` |
| p | ⭐ **Keputusan c BUKAN penyimpangan dari ADR-U-0016.** Akibat 2 ADR itu berbunyi *"kolom persen tidak termasuk — ia bukan uang, dan **tetap mengikuti ketetapan modulnya**"*. Keputusan c **adalah** ketetapan modul Claim Life, persis yang Akibat 2 serahkan. Ronde 3 menulis "PENYIMPANGAN SADAR" di kepala `004`/`005`/`006` dan di tiket 14 — label itu keliru dan menciptakan `[terbuka]` amandemen ADR yang **tidak perlu**. Usulan: **(p1)** c berlaku **untuk modul Claim Life** → label diganti *"ketetapan modul Claim Life atas kolom persen/rate, sesuai ADR-U-0016 Akibat 2"*, nol amandemen ADR; **(p2)** bila Anda ingin c berlaku **seluruh sistem**, barulah ADR-U-0016 Akibat 2 diamandemen — satu blok bertanggal, bukan menulis ulang | `[USULAN]` — rekomendasi **p1** |
| t | ⭐ **AC 12 tiket 14.** Klausa *"desimal presisi arbitrer"* di AC itu dibaca ronde 3 sebagai "`NUMBER` tanpa presisi", lalu dinyatakan dibalik oleh c. Bacaan lain, yang sejalan dengan rujukan AC-nya sendiri *(ADR-U-0003, ditulis "ADR-0003" di tiket: bukan float; test yang menemukan uang bertipe teks atau lewat float gagal)*: **"desimal presisi arbitrer" = jenis data desimal non-float, lawan dari `float`** — dan `NUMBER(38,8)` ↔ `apd.Decimal` memenuhinya. Usulan: satu baris dari Anda *"AC 12 terpenuhi oleh NUMBER(38,8); klausa sequence tercakup penyimpangan sadar AC 34"* → executor mencentang AC 12, daftar terbuka menjadi 14 | `[USULAN]` — rekomendasi **centang** |
| o1 | **ADR pengganti ADR-U-0006.** ADR-U-0006 sendiri menutup dengan *"bila kelak penomoran dipindah ke aplikasi, itu keputusan baru yang menggantikan ADR ini"*. Keputusan o *(26-09: "jangan ada lagi pemanggilan procedure, segala procedure hardcode dalam skrip")* adalah keputusan itu. Usulan: berkas baru `docs\adr\0043-penomoran-klaim-ditulis-di-aplikasi.md` *(status accepted, tanggal 2026-09-26, sumber: keputusan work owner, **menggantikan: ADR-U-0006**)*; ADR-U-0006 mendapat satu baris `status: superseded oleh ADR-U-0043` di frontmatter, isinya tidak diubah. Isi 0043: keputusan; alasan *(arahan work owner)*; akibat — sumber procedure **wajib** dari DBA *(ADR-U-0006 sudah menyatakan replikasi tanpa sumber = menebak)*, ruang nomor tetap **satu** dengan Pega selama koeksistensi *(o3)*, OQ-002 dan OQ-013 tetap terbuka | `[USULAN]` — executor menulis **hanya** bila `[DIPUTUSKAN]` |
| o2 | **Teks AC tiket 02 yang bertentangan dengan o** — bukan hanya AC 2 dan 3 seperti dicatat ronde 3: **AC 10** juga menyebut lock `SELECT … FOR UPDATE` *"pada `GENERATE_SEQUENCE_NUMBER`"* yang hari ini terjadi **di dalam** procedure. Teks AC hanya boleh diubah work owner. Usulan teks pengganti, untuk Anda tempel: **AC 2** → *"Nomor klaim dibentuk di `services`, dengan logika yang disalin dari sumber `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` yang diserahkan DBA. Nol pemanggilan procedure atau blok PL/SQL dari kode; test yang menemukan `BEGIN`/`CALL`/`EXEC` di teks SQL gagal."* **AC 3** → *"Logika penomoran berada di satu tempat dan diuji murni tanpa Oracle terhadap contoh nomor AC 7; pembacaan dan penulisan counter `GENERATE_SEQUENCE_NUMBER` hanya di `repository`."* **AC 10** → *"Counter `GENERATE_SEQUENCE_NUMBER` dikunci lewat `SELECT … FOR UPDATE` yang dikirim Go di dalam transaksi Go; commit segera setelah nomor terbentuk; kegagalan sesudahnya tidak membatalkan nomor yang sudah terbentuk (lubang nomor dicatat, bukan disembunyikan)."* AC 5, 7, 8, 9, 11 **tetap** — kini menjadi spesifikasi logika Go | `[USULAN]` — pemilik teks: work owner |
| o3 | **Ruang nomor selama koeksistensi.** Pega produksi masih memanggil procedure yang sama. Usulan: aplikasi baru membaca-menulis **tabel counter yang sama** (`GENERATE_SEQUENCE_NUMBER`) dengan kunci baris yang sama, sehingga dua pembangkit **berbagi satu urutan** dan tidak bentrok. Ini hanya benar bila logika Go **identik** dengan procedure — karena itu sumber procedure bukan pilihan, melainkan syarat | `[USULAN]` — work owner + DBA |
| k, l | tetap sebagaimana brief ronde 2 | berlaku |

---

## 3. TEMUAN VERIFIKASI RONDE 3 — kecil, semuanya tentang kata

| # | Temuan | Letak | Yang dikerjakan |
| ---: | --- | --- | --- |
| 1 | ⚠️ **Label "PENYIMPANGAN SADAR dari ADR-U-0016 Akibat 2" keliru** *(lihat §2 p)*. Akibat 2 justru **menyerahkan** kolom persen ke modul; c adalah ketetapan modul itu. Label ini ada di kepala `004`, `005`, `006` *(paragraf yang sama disalin tiga kali)*, di tiket 14 bab ronde 3 *("Penyimpangan sadar yang ditemukan tinjauan")*, dan di `presisiSah` `strukturkolom_test.go` *(komentar nilai `NUMBER(38,8)`)* | tiga `.sql`, tiket 14, `strukturkolom_test.go` | sesudah **p** `[DIPUTUSKAN]`: p1 → ganti label di keempat tempat menjadi *"ketetapan modul Claim Life, sesuai ADR-U-0016 Akibat 2"* dan cabut `[terbuka]` amandemen; p2 → amandemen ADR + label *"sesuai ADR-U-0016 amandemen 2026-09-2x"*. Tanpa keputusan: **tidak disentuh** |
| 2 | ⚠️ **Konsekuensi o di tiket 02 kurang satu AC**: ronde 3 menulis "AC 2 dan 3 harus ditulis ulang"; **AC 10** juga *(§2 o2)* | tiket 02 bab `## Implementasi` | tambah satu baris ralat: *"ralat 2026-09-2x: AC 10 ikut, sebab lock `FOR UPDATE` yang disebutnya hari ini ada di dalam procedure"* |
| 3 | ⚠️ **AC 12 dibaca terlalu sempit** *(§2 t)* — "presisi arbitrer" ≠ "tanpa presisi" menurut rujukan AC-nya sendiri | tiket 14 | sesudah **t** `[DIPUTUSKAN]`: centang, pindahkan dari daftar terbuka, daftar menjadi 14 = kotak `[ ]` |
| 4 | ℹ️ `TestNolNamaTabelTelanjangDiQuery` memindai **komentar** juga *(regex `(?i)` atas seluruh teks berkas)*. Komentar berbahasa Inggris yang memuat "from X" atau "join Y" akan gagal palsu. Hari ini 33 rujukan, nol palsu | `batasanpemakaian_test.go` | catatan saja; **tidak diubah** kecuali gagal palsu benar-benar muncul — saat itu buang baris komentar (`//`) sebelum memindai |
| 5 | ⛔ **Skema uji tidak punya pagar.** `skemauji.Pasang` memanggil `Bongkar` lebih dulu, dan `Bongkar` mengirim `DROP TABLE <skema>.OS_AKSEPTASI_KLAIM_LIFE CASCADE CONSTRAINTS` **tanpa syarat** *(hanya `ORA-00942` yang ditoleransi)*, lalu `BongkarMigrasi`, lalu `DROP TABLE <skema>.T_MIGRASI`. Satu-satunya penjaga adalah `IS_PEGA_PROD`. Bila `ORACLE_SCHEMA` menunjuk skema **DEV** yang memuat tabel warisan sungguhan *(misalnya `POOLDATA` di instance pengembangan)* dan `IS_PEGA_PROD=false`, `go test -tags=db` **menghapus tabel warisan itu**. `Makefile` bahkan memberi bawaan `ORACLE_SCHEMA ?= POOLDATA`. Ditemukan 26-09 saat menjawab pertanyaan koneksi | `skemauji.go` `Buka`/`Bongkar`, `Makefile`, `README-BACA-DULU.md` | **wajib sebelum test db pertama:** `Buka()` menolak *(galat, bukan SKIP)* kecuali env `ORACLE_SKEMA_UJI=true` **dan** `ORACLE_SCHEMA` **bukan** `POOLDATA` *(dua syarat, keduanya)*; pesan galatnya menyebut alasannya; test murni untuk penjaga ini *(tiga kasus: tanpa env, skema `POOLDATA`, keduanya benar)*; bawaan `POOLDATA` di `Makefile` dicabut; README bab 5 memuat larangan yang sama dalam satu kalimat. **`-migrate` tidak ikut dipagari** — ia hanya `CREATE`, tidak pernah `DROP` |

Yang **sudah benar** dan tidak perlu disentuh: 28 + 1 + 1 kolom `NUMBER` semuanya berpresisi;
`sudahAda` hanya `ORA-00955` dan pembuktian keberadaan lewat `SYS.ALL_OBJECTS` dengan `UPPER` di
kedua sisi; `namaObjekDibuat` dipanggil atas pernyataan **mentah** `p` *(masih `{skema}`)*, sesuai
regex-nya; 19 `CREATE` dikunci cacahnya; test tabrakan constraint menuntut galat menyebut `001` dan
`ORA-02264` **dan** `T_MIGRASI` kosong; komentar `005` yang basi dicabut; daftar terbuka 15 = kotak
`[ ]`; ralat angka ronde 2 benar; kutipan ADR benar; `PANDUAN-MENJALANKAN.txt` utuh dan tidak
ter-commit; nol host/sandi/nama orang; fixture `UJI-*`.

---

## 4. URUTAN SESI

**Langkah 0 — titik tetap.** `git add PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14-RONDE-4.md
APP_RNM/PANDUAN-MENJALANKAN.txt` *(keduanya, dengan nama — bukan `git add -A`)*, commit
`docs: brief sesi tiket 14 ronde 4 + panduan menjalankan`. Lalu `git status --porcelain` kosong;
`git rev-parse HEAD` dicatat. Dari `APP_RNM\`: `go vet ./...`, `go test ./...`,
`cd frontend; npm run typecheck; npm test` — hijau sebelum apa pun diubah.

**Langkah B — sunting kecil §3** *(hanya yang keputusannya sudah `[DIPUTUSKAN]`; §3-2 tanpa
syarat)*. Urutan: 3-2 → 3-1 → 3-3. Nol perubahan perilaku kode; `go test ./...` tetap 77.

**Langkah C — keputusan §2 tiket 14 yang `[DIPUTUSKAN]`.** **d** → `REFERENCES` di `001`
(`COVER_KEY`) dan `004` (`KOMITE_ID`) *(atau, bila keputusannya "di Go": catatan di kedua berkas dan
di tiket, tanpa constraint)*; cabut `[terbuka]`-nya; centang AC 50. **e′** → catatan AC 8, `002`
tidak disentuh. **m** → satu blok ralat bertanggal di STRUKTUR. **n** → `git mv` dua berkas `006`
**sebelum** Langkah A-2; `TestSetiapLangkahPunyaJalurMundur` dan `TestSeluruhCreateDapatDibacaNamanya`
tetap lulus. Yang masih `[USULAN]` **tidak disentuh**.

**Langkah A — hanya bila G1 terbuka.**
1. `go test -tags=db ./internal/... -count=1 -v` — seluruh **19** SKIP harus menjadi PASS, termasuk
   `TestNamaConstraintBertabrakanMenggagalkanMigrasi` dan `TestLangkahGagalSeparuhJalanTetapSelesai`.
   Yang gagal diperbaiki **di ronde ini** dan dicatat apa adanya.
2. **Migrasi sungguhan**: `go run ./cmd/api -migrate` dua kali. Run kedua: **0 dijalankan, 8
   dilewati, `ObjekSudahAda` kosong**. Query `SELECT OBJECT_TYPE, COUNT(*) FROM SYS.ALL_OBJECTS WHERE
   OWNER = :1 GROUP BY OBJECT_TYPE` — diharapkan **8** TABLE · **5** SEQUENCE · **15** INDEX
   *(7 eksplisit + 8 milik `PRIMARY KEY`)*; selisih **dicatat**.
3. Tutup tiket 01 bila 11 test db-nya PASS: centang tiga AC sisanya, `Status: resolved`, commit
   terpisah `claim-life: tiket 01 — penutupan AC yang menunggu Oracle`.

**Langkah D — persiapan tiket 02, hanya bila G3 terbuka.**
1. Tulis `docs\adr\0043-penomoran-klaim-ditulis-di-aplikasi.md` sesuai §2 o1; tambahkan baris
   `status: superseded oleh ADR-U-0043` ke frontmatter ADR-U-0006 **tanpa mengubah isinya**.
2. Simpan sumber procedure dan DDL dari DBA sebagai `.scratch\claim-life\SUMBER-PENOMORAN-DBA.md`
   bertanda `[data DBA]` — **teks procedure boleh, isi tabel counter sungguhan tidak**.
3. Di tiket 02 bab `## Implementasi`: catat bahwa o1–o3 `[DIPUTUSKAN]`, AC 2/3/10 sudah ditulis ulang
   work owner *(kutip tanggalnya)*, dan blocker 1–4 ronde 3 mana yang tertutup. **Tidak** mulai
   mengerjakan AC tiket 02 di sesi ini — itu sesi tersendiri dengan brief tersendiri.

**Verifikasi penuh sekali di akhir:** `go vet ./...`, `go vet -tags=db ./...`, `gofmt -l`,
`go test ./...`, `go test -tags=db ./internal/...`, `cd frontend; npm run typecheck; npm test; npm run build`.

**Penutup — brief induk §6 butir 6–9:** bab `## Implementasi — ronde 4, <tanggal>` di tiket 14
*(gerbang mana yang terbuka; §3 ditutup lawan 4; AC ditutup lawan 53 dengan daftar terbuka = kotak
`[ ]`; keputusan diterapkan lawan ditahan; hasil Langkah A bila ada)*, `/code-review` atas titik
tetap Langkah 0, **angka ditulis sesudah perbaikan review**, commit
`claim-life: tiket 14 ronde 4 — <isi sebenarnya>` + baris `Tiket: .scratch/claim-life/issues/14-…md`.
Perubahan ADR dan tiket 02 masuk commit **terpisah** `docs: ADR-U-0043 penomoran di aplikasi`.

**Tanda `resolved` untuk tiket 14** — hanya bila **semuanya**: **d** `[DIPUTUSKAN]` dan
diterapkan *(apa pun isinya)*; Langkah A-2 selesai 0/8/kosong; seluruh test db PASS; tiap AC
`[terbuka]` tinggal yang pemiliknya di luar executor *(8, 15, 27, 35 work owner; 31, 33 konteks
Komite; 46 pemilik export Pega; 19 modul PremiumList Life; 12 bila t tidak disahkan)*. Bila satu
saja tidak terpenuhi: `claimed`, dengan alasan yang mana.

---

## 5. YANG DIMINTA KE DBA — satu daftar, supaya diminta sekali

| # | Objek | Untuk | Bentuk yang diminta |
| ---: | --- | --- | --- |
| 1 | Skema uji + `ORACLE_DSN` | G1 — seluruh test db, migrasi sungguhan | DSN di terminal, hak `CREATE TABLE / SEQUENCE / INDEX / DROP` |
| 2 | `ALL_TAB_COLUMNS` untuk `OS_AKSEPTASI_KLAIM_LIFE` | menutup uji swa-konfirmasi (ronde 3 §3-7) | 55 baris katalog: nama, tipe, panjang, presisi, skala, nullable |
| 3 | Versi Oracle + `COMPATIBLE`, pengembangan **dan** produksi | batas 30 byte pengenal | dua angka |
| 4 | **Sumber `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`** + DDL `GENERATE_SEQUENCE_NUMBER` | G3 — tanpa ini o berarti menebak *(ADR-U-0006 Considered Options)* | teks procedure dan DDL; **bukan** isi counter |
| 5 | `POOLDATA.KODE_PRODUKSI` baris `TYPE='LIFE'` | prefix nomor (AC 8 tiket 02) | DDL + baris LIFE |
| 6 | `POOLDATA.TANGGAL_CLOSING` + aturan cutover | periode nomor (AC 9 tiket 02) | DDL + contoh isi |
| 7 | `M_LIFE_PREMIUM_DETAIL` | pencarian peserta ber-`EDMSTATUS` (AC 20–22 tiket 02) | DDL, dan jawaban OQ-001 sisa: NB `NULL` atau `''` |

---

## 6. GAYA KODE — tetap mengikat

Pembaca **baru mengenal Go dan React**: komentar kepala tiap berkas; istilah dijelaskan sekali; satu
konsep per komentar; nama berbahasa Indonesia yang terbaca sebagai kalimat. Perbarui
`APP_RNM\README-BACA-DULU.md` bila ada perintah, berkas, atau nama berkas migrasi yang berubah.
ADR baru mengikuti bentuk ADR yang ada: frontmatter `status / tanggal / sumber`, bab Kenapa, Akibat,
Yang TIDAK diputuskan.

---

## 7. YANG MEMERLUKAN PERSETUJUAN MANUSIA — berhenti dan tanya

Brief induk §7, ditambah: menulis atau mengubah berkas di `docs\adr\` *(hanya o1 `[DIPUTUSKAN]`)* ·
mengubah teks AC tiket mana pun *(tidak pernah; itu milik work owner)* · `-migrate` terhadap skema
selain skema uji · mengganti nama tabel/berkas migrasi di luar §2 j/n · menyunting STRUKTUR di luar
blok ralat m · menyimpan isi tabel produksi apa pun *(sumber procedure boleh, isi counter tidak)* ·
menyentuh `dastin\` / `jefri\` · `git push`.

---

## 8. TELEMETRI EKSEKUSI — bab wajib di laporan akhir

| Besaran | Cara ukur |
| --- | --- |
| Gerbang yang terbuka (G1 / G2 / G3) dan keputusan §2 yang diterapkan | tulis di baris pertama laporan |
| Temuan §3 ditutup / 4; AC tiket 14 ditutup / 53 *(daftar terbuka = kotak `[ ]`)*; AC tiket 01 ditutup / 7 | hitung di tiket |
| `-migrate` sungguhan: dijalankan · dilewati · pernyataan · `ObjekSudahAda`, run pertama dan kedua | log `jalankanMigrasi` |
| Objek DB nyata per `OBJECT_TYPE` sesudah migrasi | query Langkah A-2 |
| Test: PASS · FAIL · SKIP per tag — sebut apakah SKIP atau PASS | keluaran `go test -v` |
| Berkas dibuat / diubah / diganti nama, baris berisi; SHA titik tetap dan tiap commit — **diukur sesudah commit terakhir** | `git diff --stat <titik tetap>..HEAD` |
| Sub-agen review: token, panggilan alat | laporan harness |
| Token sesi utama · biaya · jam dinding | **⛔ tidak diukur** — nyatakan |

---

*Disusun 26 September 2026 dari verifikasi independen commit `7b832e6`: 77 test dijalankan ulang,
diff 11 berkas dibaca utuh, ADR-U-0006 dan ADR-U-0016 dibaca utuh, cacah kolom `NUMBER`, pernyataan
`CREATE`, dan byte CR dihitung ulang dengan alat sendiri.*
