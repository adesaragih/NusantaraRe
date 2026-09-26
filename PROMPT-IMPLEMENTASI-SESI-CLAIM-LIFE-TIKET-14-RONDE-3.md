# PROMPT — sesi implementasi berikutnya: Claim Life tiket 14 **ronde 3** (penuntasan keputusan + dua cacat migrasi) + penutupan tiket 01

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya. Brief ronde 2
> **`PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14-RONDE-2.md`** tetap berlaku untuk hal yang tidak
> diubah di sini *(§2 k–l, §5 tabel objek DBA, §6, §7)*.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> tiket 14 bab `## Implementasi — ronde 2, 26 September 2026` **seluruhnya** → §3 berkas ini →
> `.scratch\claim-life\STRUKTUR-TABEL-CLAIM-LIFE.md` bab `T_CLAIMLF_ADJUSTMENT` sampai
> `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` *(tiga tabel yang keputusan §2 c belum sampai ke sana)*.
>
> **SESI INI:** claim-life · tiket **14 ronde 3** + tiket **01 (penutupan sisa AC, bila Oracle ada)**

---

## 0. KEADAAN AWAL — 26 September 2026

| | Keadaan |
| --- | --- |
| `HEAD` | `8f5453b` — tiket 14 *claimed* 38/53 AC; tiket 01 *claimed* 4/7. Working tree bersih kecuali berkas brief ini |
| Uji yang lulus | `go vet`, `go vet -tags=db`, `gofmt` nol, `go build`; **73** test Go PASS; **18** test bertag `db` SKIP dengan pesan; `tsc --noEmit`; **5** test JS; `vite build` 87 modul |
| Migrasi | 16 berkas `.sql`, LF, tanpa BOM, nol pengenal > 30 byte; 7 tabel · 5 sequence · 7 index · 7 PK · 6 FK · 4 cascade — **belum pernah dijalankan di Oracle mana pun**, tiga sesi berturut-turut di jalur B |
| Keputusan yang sudah diterapkan ronde 2 | **i** (`ACCEPT_STATUS` dibuang), **j** (tabel retro → `T_CLAIMLF_ADJ_SPREADING_RETRO`, dua kolom valuasi → nama korpus), **c hanya untuk `003`** *(lihat §3-1)* |
| Keputusan yang masih menunggu | **d** *(tiga sesi)*, **e** *(usulan lama dibalik oleh temuan korpus ronde 2, lihat §2 e′)* |
| Rantai alat | `go` dan `node` ada di mesin, tidak di PATH shell; `make` tidak ada; `.env` tidak dibaca kode *(env var disetel di terminal)* |

**Verifikasi independen 26 September 2026 atas commit `8f5453b`:** seluruh angka laporan ronde 2
tereproduksi *(73 · 18 · 5 · 38/53 · LF · nol nama > 30 byte)*; sebelas temuan ronde 2 memang
tertutup seperti dilaporkan; temuan korpus atas ketiga `TANGGAL_*` sahih. Yang **belum tepat** ada di
§3 — dua di antaranya cacat sungguhan, dan satu dari keduanya **berasal dari brief ronde 2 sendiri**.

---

## 1. PRASYARAT MANUSIA — sebelum sesi dimulai

1. **Oracle pengembangan** — DSN skema uji, hak `CREATE TABLE / SEQUENCE / INDEX / DROP`; env var
   di terminal, tidak pernah di berkas atau chat. Tanpa ini sesi keempat pun berakhir `claimed`.
2. ⭐ **Minta ke DBA daftar tipe kolom `OS_AKSEPTASI_KLAIM_LIFE` dari instance pengembangan** —
   hasil `SELECT COLUMN_NAME, DATA_TYPE, DATA_LENGTH, DATA_PRECISION, DATA_SCALE, NULLABLE FROM
   ALL_TAB_COLUMNS WHERE TABLE_NAME = 'OS_AKSEPTASI_KLAIM_LIFE' ORDER BY COLUMN_ID` *(55 baris
   katalog, bukan isi tabel)*. Ini satu-satunya yang dapat menutup **uji yang mengonfirmasi dirinya
   sendiri** *(tiket 14 ronde 2, kelemahan terbuka 1; §3-7 di sini)*. Simpan jawabannya sebagai
   `.scratch\claim-life\TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md` bertanda `[data DBA]`.
3. **Versi Oracle dan `COMPATIBLE`** — masih belum dijawab sejak brief ronde 2 §1-2. Keputusan **j**
   sudah memilih nama 29 byte tanpa menunggu jawaban itu, jadi kini pertanyaan ini tinggal informasi;
   tetap tanyakan, sebab produksi bisa berbeda dari pengembangan.
4. **Putuskan §2.** Baris `[USULAN]` disahkan dengan mengganti kata itu menjadi `[DIPUTUSKAN]` atau
   ditulis ulang. Butir **d** sudah **tiga sesi** tertunda dan sendirian menahan tiket 14 di `claimed`.

---

## 2. KEPUTUSAN WORK OWNER

| | Keputusan | Keadaan |
| ---: | --- | --- |
| c | Kolom share / persen / rate bertipe `NUMBER(38,8)`; di Go `Ratio`, bukan `Money` | `[DIPUTUSKAN]` 26-09 — **penerapannya belum tuntas**: baru `003`; delapan kolom di `004`/`005`/`006` masih `NUMBER` polos *(§3-1; dikerjakan ronde ini tanpa keputusan baru)* |
| d | `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)`, nullable. **Tanpa keputusan ini tiket 14 tidak pernah `resolved`** *(tiket 14 §Blocker)* | `[USULAN]` — belum dijawab tiga sesi |
| e′ | **Ketiga `TANGGAL_*` bukan milik Claim Life.** Temuan korpus ronde 2 *(diverifikasi ulang: sumbernya `PremiumList Life/RDBList/GetOfferLife_sql.xml` atas `POOLDATA.JSON_OFFER_LIFE`, dipasang `setNoOffer_Act.xml` kelas `ASM-FW-GISFW-Work-LIFE`; di Claim Life hanya tampil di 4 berkas `Section\`, nol `Activity\`, nol di `UpdateOsAkseptasiClaimLife_sql.xml`)*. Usulan **menggantikan** usulan lama: ketiganya **tidak dibuat** di `T_GENERAL_CLAIM` maupun `T_WORK_CLAIM`; Claim Life membacanya **hidup** lewat rantai penunjuk polis → offer, bentuk relasionalnya milik modul **PremiumList Life**. AC 8 tiket 14 diberi `[terbuka — PremiumList Life]` untuk tiga kolom itu; `TEAM_GROUP` tetap lewat `MO_ID`; `BUSINESS_ID` tetap terbuka | `[USULAN]` |
| m | **Ralat nama fisik di STRUKTUR.** Sejak **j**, DDL memakai `T_CLAIMLF_ADJ_SPREADING_RETRO`, sedangkan `STRUKTUR-TABEL-CLAIM-LIFE.md`, `spec.md`, tiket 14 *(AC 2, 38, 40 dan diagram)*, tiket 15, dan `revisi-penyimpanan-json-dibuang.md` masih menulis nama 36 byte. Usulan: **satu blok ralat bertanggal** di kepala bab STRUKTUR `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` *("nama fisik sejak 26-09-2026: `T_CLAIMLF_ADJ_SPREADING_RETRO`, keputusan work owner j; nama logis di dokumen tidak diubah")*, dan **tidak menyentuh** spec, tiket 15, maupun teks AC tiket 14 | `[USULAN]` |
| n | **Nama berkas migrasi 006** masih `006_t_claimlf_adjustment_spreading_retro.sql` *(+ `_down`)*. Nama berkas = kunci `T_MIGRASI`, jadi **hanya boleh diganti selama belum satu instance pun menjalankannya** *(§2 l)* — yaitu ronde ini, **sebelum** Langkah A-2. Usulan: ganti menjadi `006_t_claimlf_adj_spreading_retro.sql` dan `_down`-nya lewat `git mv` | `[USULAN]` |
| o | ⭐ **Penomoran tiket 02 lawan arahan "tidak ada lagi pemanggilan procedure".** Tiket 02 AC 2 menuntut nomor dari `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` *(ADR-U-0006: jangan replikasi logikanya)*, dan blok PL/SQL warisannya memuat `COMMIT` yang `PeriksaSQL` tolak. Dua jalan, keduanya melanggar sesuatu: **(o1)** Go memanggil procedure lewat `BEGIN … END;` **tanpa** `COMMIT` di teks, Go yang commit — melanggar arahan "tidak ada lagi pemanggilan procedure"; **(o2)** DBA/work owner menyerahkan logika penomoran untuk ditulis ulang di Go — melanggar ADR-U-0006. **Bukan** untuk dikerjakan ronde ini; diputuskan sekarang supaya sesi tiket 02 tidak berhenti di baris pertamanya | `[USULAN]` — pemilik work owner + DBA |
| k, l | tetap sebagaimana brief ronde 2 | berlaku |

---

## 3. TEMUAN VERIFIKASI RONDE 2 — wajib ditutup ronde ini, urut dari yang terberat

| # | Temuan | Letak | Yang dikerjakan |
| ---: | --- | --- | --- |
| 1 | ⛔ **Keputusan c diterapkan sebagian.** Delapan kolom share/rate masih `NUMBER` polos: `004` `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`, `SHARE_RETRO`, `RETROCEDED_SHARE`; `005` `RETROCADED_SHARE`, `RATE`; `006` `PERCENT_SHARE`, `RATE`. Laporan ronde 2 menulis "kolom share/persen/rate kini `NUMBER(38,8)`" — benar hanya untuk `003`; tiket sendiri menulis "lima kolom `003`". `CEDING_RETENTION` bertipe `NUMBER(38,8)` di `003` tetapi `NUMBER` di `004`. **Tidak ada test yang menangkapnya**: `TestGolonganTipeDDLCocokDenganStruktur` sengaja kasar, `TestAC41…` hanya menuntut kata `NUMBER` | `004`, `005`, `006` | ubah kedelapannya menjadi `NUMBER(38,8)`; komentar kepala ketiga berkas menyebut keputusan c; `TestAC41…` diperketat menjadi `NUMBER(38,8)`; test baru **nol `NUMBER` tanpa presisi** di seluruh berkas maju *(yang sah hanya `NUMBER(38,8)`, `AGE NUMBER(5)`, `DOCUMENT_CLAIM.ID NUMBER(19)`)* — diuji gagal dulu dengan satu kolom dikembalikan ke `NUMBER` |
| 2 | ⛔ **Toleransi `ORA-02264` pada `CREATE` tidak aman — dan itu usulan brief ronde 2 §3-5 yang keliru.** `ORA-02264` berarti *nama constraint sudah dipakai objek lain*, dan Oracle baru memeriksanya bila **tabelnya belum ada** *(bila ada, jawabannya `ORA-00955` lebih dulu)*. Kode kini melewati `CREATE TABLE` itu, mencatatnya di `ObjekSudahAda`, lalu mencatat langkahnya **sukses** di `T_MIGRASI` — padahal tabelnya **tidak pernah dibuat**. Skenario nyata sesudah **j**: instance yang pernah memuat tabel lama dengan `PK_T_CLAIMLF_SPR_RETRO` → `CREATE TABLE T_CLAIMLF_ADJ_SPREADING_RETRO` gagal `ORA-02264` → dilewati diam-diam. `TestPenggolongGalatObjekSudahAda` justru **mengunci** perilaku ini *(`ORA-02264` di daftar `harusYa`)* | `migrasi.go` `sudahAda`, `JalankanMigrasi`; `migrasi_test.go`; `pohonklaim_db_test.go` | **(a)** `sudahAda` hanya `ORA-00955`; `ORA-02264` dipindah ke `harusTidak`. **(b)** Sesudah sebuah `CREATE` dilewati, **buktikan objeknya ada**: `SELECT COUNT(*) FROM ALL_OBJECTS WHERE OWNER = :1 AND OBJECT_NAME = :2` *(nama diambil dari `{skema}.<nama>` pernyataan itu)*; nol → kembalikan galat *"dilaporkan sudah ada, tidak ditemukan"*, langkah **tidak** dicatat. Ini sekaligus menutup separuh kelemahan terbuka 2 ronde 2 *(bentuk masih tidak diperiksa, keberadaan kini diperiksa — tulis begitu)*. **(c)** Test db baru: sesudah `BongkarMigrasi`, buat `UJI_TABRAKAN (ID VARCHAR2(1), CONSTRAINT PK_T_WORK_CLAIM PRIMARY KEY (ID))`, lalu `JalankanMigrasi` **harus gagal** dengan galat menyebut `001` dan `ORA-02264`, dan `T_MIGRASI` **tidak** memuat `001`; bersihkan `UJI_TABRAKAN` di `defer` |
| 3 | ⚠️ **AC 12 tidak dicentang dan tidak ada di daftar "AC yang masih terbuka — 15"** — daftar itu memuat 14 nomor *(20 21 38 51 53 · 8 15 27 35 50 · 31 33 · 19 46)*. AC 12 = *"seluruh uang dan share desimal presisi arbitrer; tanggal `DATE`; nullable; identitas dari sequence"*. Alasannya tidak tertulis: mungkin §3-1, mungkin bagian "identitas dari sequence" yang memang penyimpangan sadar *(AC 34)* | tiket 14 | sesudah §3-1: centang **bila** bagian sequence dinyatakan tercakup penyimpangan sadar AC 34 — atau tulis alasannya di daftar terbuka. Daftar terbuka harus **berjumlah sama** dengan kotak `[ ]` |
| 4 | ⚠️ **Kelemahan terbuka 3 ronde 2 keliru** — "`UX_ADJ_KOMITE_ID` UNIQUE, konsekuensinya belum pernah dibahas". STRUKTUR bab `T_CLAIMLF_ADJUSTMENT` memuat `[keputusan work owner]`: *"satu baris `AdjustmentList` = TEPAT satu kasus komite … `KOMITE_ID` ber-index UNIK meski nullable"*. Sudah diputuskan, dan index unik Oracle memang mengizinkan banyak `NULL` | tiket 14 bab Implementasi ronde 3 | tutup butir itu dengan menunjuk bab STRUKTUR-nya; **tidak** ada perubahan DDL |
| 5 | ⚠️ **AC 40 dicentang `[x]` padahal teks AC-nya menyebut nama tabel lama** `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`; test `TestAC40…` memeriksa nama baru. Sah secara isi, tidak sah secara jejak — pembaca tiket tidak tahu ada ganti nama | tiket 14, STRUKTUR *(§2 m)*, berkas `006` *(§2 n)* | catat di bab Implementasi ronde 3 satu baris *"AC 2, 38, 40: nama fisik = `T_CLAIMLF_ADJ_SPREADING_RETRO` per keputusan j"*; lalu §2 m dan n **hanya bila `[DIPUTUSKAN]`** |
| 6 | ⚠️ **Angka di bab Implementasi ronde 2 basi** — ditulis sebelum perbaikan `/code-review`: "24 berkas +2.247/−136" *(nyata: 76dcda4→8f5453b **+2.293**/−136 atas 24 berkas, termasuk brief; 8b18b00→8f5453b +2.123/−136 atas 23)*; "`strukturkolom_test.go` 231/213" *(nyata 368/344)*; "jumlah berkas baru 1.279/1.201" *(nyata tujuh berkas 1.206 berisi)* | tiket 14 | satu baris ralat di bab ronde 3; angka ronde 3 diambil **sesudah** commit terakhir, bukan sebelum review |
| 7 | ⚠️ **Uji tabel warisan mengonfirmasi dirinya sendiri** *(diakui ronde 2)*: `kolomAngkaLama` / `kolomTanggalLama` dan `TipeKolomBarisLama` sama-sama tebakan atas nama | `barislamakolom.go`, `skemauji.go` | **hanya bila §1-2 sudah ada:** kedua peta dan `TipeKolomBarisLama` diturunkan dari berkas `[data DBA]` itu; test murni membandingkan 55 kolom satu per satu dengan dokumen tersebut; tabel tiruan memakai tipe DBA. Tanpa §1-2: **tidak disentuh**, tetap `[terbuka]` |
| 8 | ⚠️ Kalimat *"nama `TANGGAL_RESPON` nol kemunculan … itu nama karangan spec"* terlalu jauh: yang nol adalah ejaan `TANGGAL_RESPON`; properti Pega **`TanggalRespon`** ada *(`setNoOffer_Act.xml`, 2 kemunculan; 4 `Section\` Claim Life)*. Spec menurunkan nama kolom dari properti sungguhan, bukan mengarangnya. Kesimpulan pokoknya *(tanggal penawaran, milik Offer)* **tetap sahih** | tiket 14 | satu baris koreksi di bab ronde 3 |
| 9 | ℹ️ `TestKolomDDLCocokDenganStruktur` membaca `../../../.scratch/claim-life/STRUKTUR-…md` — berkas **di luar** `APP_RNM\`. Sah selama satu repositori; gagal terang bila `APP_RNM` dipindah *(memang `t.Fatalf`)* | — | catatan saja, **tidak diubah** |

Yang **sudah benar** dan tidak perlu disentuh: `ACCEPT_STATUS` hilang dari DDL, model, `INSERT`;
47 kolom `003` = STRUKTUR; `AmbilBarisLama` + test membaca kembali dari tabel; sebelas kolom
uang/rasio cicit dibaca dan galat dikembalikan; `PernyataanLangkah` untuk test separuh jalan;
`TestAC45Dan48` memeriksa `CREATE INDEX` sungguhan; kolom bank pulang-pergi; 18/55 dikunci test;
CR + `.gitattributes`; `Hapus` dijaga test statik; nol host/sandi/nama orang; fixture `UJI-*`.

---

## 4. URUTAN SESI

**Langkah 0 — titik tetap.** `git add PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14-RONDE-3.md`
dan commit `docs: brief sesi tiket 14 ronde 3 dari verifikasi independen 8f5453b`. Lalu
`git status --porcelain` kosong; `git rev-parse HEAD` dicatat. Dari `APP_RNM\`: `go vet ./...`,
`go test ./...`, `cd frontend; npm run typecheck; npm test` — hijau sebelum apa pun diubah.

**Langkah B dulu, A sesudahnya** — urutan ini **dibalik** dari ronde 2, sebab §3-2 harus diperbaiki
dan §2 n *(bila disahkan)* harus terjadi **sebelum** migrasi sungguhan pertama.

**Langkah B — tiket 14 ronde 3.** Status tetap `claimed` sampai Langkah A. Urutan: §3 nomor
**1 → 2 → 5** *(ralat baris + §2 n bila `[DIPUTUSKAN]`)* **→ 3 → 4 → 6 → 8 → 7** *(hanya dengan
§1-2)*. Tiap nomor: test dulu *(`/tdd` pada seam brief induk §5)*, diuji **gagal** pada kasus
buruknya, lalu kode, lalu `go vet` dan test berkas itu. Tanpa Oracle: test db ditulis dan **SKIP
dengan pesan**; jangan mengarang hasil.

**Langkah C — keputusan §2 yang `[DIPUTUSKAN]`.** Terapkan **hanya** yang disahkan:
**d** → `REFERENCES` di `001` *(`COVER_KEY`)* dan `004` *(`KOMITE_ID`)*, cabut `[terbuka]`-nya,
centang AC 50; **e′** → AC 8: catatan `[terbuka — PremiumList Life]` untuk tiga tanggal, `002` tetap
tidak disentuh; **m** → blok ralat di STRUKTUR *(satu blok, bertanggal, tidak mengubah teks lama)*;
**n** → `git mv` dua berkas `006`, `TestSetiapLangkahPunyaJalurMundur` harus tetap lulus; **o** →
**tidak dikerjakan** di sini, hanya dicatat di tiket 02 bab `## Implementasi` sebagai keputusan
yang diterima. Yang masih `[USULAN]` **tidak disentuh** dan tetap tercatat terbuka.

**Langkah A — hanya bila `ORACLE_DSN` terisi.**
1. `go test -tags=db ./internal/... -count=1 -v` — seluruh 19+ SKIP harus menjadi PASS *(termasuk
   test §3-2c yang baru)*. Yang gagal diperbaiki **di ronde ini** dan dicatat.
2. **Migrasi sungguhan** ke skema uji: `go run ./cmd/api -migrate`. Lalu `-migrate` **kedua**
   harus melaporkan **0 dijalankan, 8 dilewati, `ObjekSudahAda` kosong**. Catat log keduanya
   *(tanpa DSN)* di tiket. Query sekali: `SELECT OBJECT_TYPE, COUNT(*) FROM ALL_OBJECTS WHERE
   OWNER = :1 GROUP BY OBJECT_TYPE` — yang diharapkan: **8** TABLE *(7 + `T_MIGRASI`)* · **5**
   SEQUENCE · **15** INDEX *(7 eksplisit + 8 yang Oracle buat sendiri untuk tiap `PRIMARY KEY`,
   termasuk `PK_T_MIGRASI`)*. Selisih dari angka itu **dicatat**, bukan disesuaikan diam-diam.
3. Tutup tiket 01: bila 11 test db-nya PASS, centang tiga AC sisanya, `Status: resolved`, commit
   `claim-life: tiket 01 — penutupan AC yang menunggu Oracle`. **Terpisah** dari commit tiket 14.

**Verifikasi penuh sekali di akhir:** `go vet ./...`, `go vet -tags=db ./...`, `gofmt -l`,
`go test ./...`, `go test -tags=db ./internal/...`, `cd frontend; npm run typecheck; npm test; npm run build`.

**Penutup — brief induk §6 butir 6–9:** bab `## Implementasi — ronde 3, <tanggal>` di tiket 14
*(temuan §3 ditutup lawan 9; AC ditutup lawan 53 dengan daftar terbuka **berjumlah sama** dengan
kotak `[ ]`; keputusan §2 yang diterapkan lawan yang ditahan; ralat §3-6 dan §3-8)*, `/code-review`
atas titik tetap Langkah 0 dengan path tiket, **angka bab ditulis sesudah perbaikan review**, commit
`claim-life: tiket 14 ronde 3 — penuntasan keputusan c, penjaga ORA-02264` + baris
`Tiket: .scratch/claim-life/issues/14-…md`.

**Tanda `resolved` untuk tiket 14** — dan hanya bila **semuanya**: §2 d `[DIPUTUSKAN]` dan
diterapkan; Langkah A-2 selesai dengan run kedua 0/8/kosong; seluruh test db PASS; tiap AC
`[terbuka]` tinggal yang pemiliknya memang di luar executor *(8, 15, 27, 35 work owner; 31, 33 konteks
Komite; 46 pemilik export Pega; 19 modul PremiumList Life)*. Bila satu saja tidak terpenuhi:
`claimed`, dengan alasan yang mana.

---

## 5. SESUDAH INI — tiket 02, prasyaratnya

Tabel objek DBA di brief ronde 2 §5 **tetap berlaku** dan belum satu pun terpenuhi:
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` + tabel `GENERATE_SEQUENCE_NUMBER`, `POOLDATA.KODE_PRODUKSI`
(`TYPE='LIFE'`), `POOLDATA.TANGGAL_CLOSING`, `M_LIFE_PREMIUM_DETAIL`. Ditambah ronde ini:
**keputusan §2 o** — tanpa itu sesi tiket 02 berhenti di AC keduanya.

---

## 6. GAYA KODE — tetap mengikat

Pembaca **baru mengenal Go dan React**: komentar kepala tiap berkas *(untuk apa, dibaca sesudah apa,
aturan mana yang dijaga)*; istilah dijelaskan sekali; satu konsep per komentar; nama berbahasa
Indonesia yang terbaca sebagai kalimat. Perbarui `APP_RNM\README-BACA-DULU.md` bila ada perintah,
berkas, atau nama berkas migrasi yang berubah.

---

## 7. YANG MEMERLUKAN PERSETUJUAN MANUSIA — berhenti dan tanya

Brief induk §7, ditambah: `-migrate` terhadap skema **selain** skema uji dari DBA · membaca tabel
produksi mana pun · mengganti nama tabel atau berkas migrasi di luar §2 j/n · menyunting STRUKTUR
di luar blok ralat §2 m · menyentuh `dastin\` / `jefri\` · `git push`.

---

## 8. TELEMETRI EKSEKUSI — bab wajib di laporan akhir

| Besaran | Cara ukur |
| --- | --- |
| Temuan §3 ditutup / 9; AC tiket 14 ditutup / 53 *(daftar terbuka = kotak `[ ]`)*; AC tiket 01 ditutup / 7 | hitung di tiket |
| Kolom `NUMBER` tanpa presisi tersisa *(harus 0)* | test §3-1 |
| `-migrate` sungguhan: dijalankan · dilewati · pernyataan · `ObjekSudahAda`, run pertama dan kedua | log `jalankanMigrasi` |
| Objek DB nyata per `OBJECT_TYPE` sesudah migrasi | query Langkah A-2, dicatat |
| Test: PASS · FAIL · SKIP per tag; lama `go test -tags=db` — sebut apakah SKIP atau PASS | keluaran `go test -v` |
| Berkas dibuat / diubah / diganti nama, baris berisi; SHA titik tetap dan commit — **diukur sesudah commit terakhir** | `git diff --stat <titik tetap>..HEAD`, `git rev-parse` |
| Sub-agen review: token, panggilan alat | laporan harness |
| Token sesi utama · biaya · jam dinding | **⛔ tidak diukur** — nyatakan |

---

*Disusun 26 September 2026 dari verifikasi independen commit `8f5453b`: 73 test dijalankan ulang,
4 berkas migrasi yang berubah dan 6 berkas Go baru dibaca utuh, diff 8 berkas Go lama dibaca,
tiket 14 bab ronde 2 dibaca utuh, klaim korpus atas `TANGGAL_*` direproduksi pada 9.461 berkas.*
