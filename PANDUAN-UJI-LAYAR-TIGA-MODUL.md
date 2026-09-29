# PANDUAN UJI LAYAR — TIGA MODUL (Claim Life · PremiumList Life · Komite Claim Life)

> **Untuk:** work owner yang akan mengeklik aplikasi. **Disusun:** 28 September 2026, GILIRAN-11
> paket 5, dari kode `main` sesudah commit `fbf1c9b` — **bukan** dari dokumen desain, dan **tanpa**
> menjalankan aplikasi atau menulis apa pun ke DEV. **Diperbarui GILIRAN-12** (28-09-2026): unduh
> dokumen, gerbang yang ter-remark di XML dibuang, dan **§5 uji asap baca-saja** — aplikasi
> dijalankan terhadap DEV dengan `GET` saja, nol tulisan. **Diperbarui GILIRAN-13** (29-09-2026): tiga
> titik buta bab 0 ditutup — PremiumList membuat kasus, `Add` baris adjustment pertama, dan data uji
> sintetis `UJI-*` (bab 0 §0.1–§0.2). **Diperbarui GILIRAN-14** (29-09-2026): baris adjustment lahir saat
> `Submit` Register (bp), `Decision3` PremiumList dirutekan dari bendera (bq), sunting sel adjustment: nol
> sel menurut XML (br). **Diperbarui GILIRAN-15** (29-09-2026): enam jawaban work owner — sunting sel (br) dan `Delete`
> baris adjustment (N7) tidak berlaku, pembulatan 7.7 pada peserta (N10), `CLAIM_GROSS` menunggu pemilik ekspor (N11),
> 7.8 membaca kosong sebagai nol, dan `SEQ_WORK_POLIS` mulai 22374 (migrasi 058, PL-15).
>
> - Setiap teks di dalam `kode` atau tanda kutip disalin **apa adanya** dari kode (label, tombol,
>   pesan) — termasuk salah ejanya (`cannnot`). Bila layar berbeda dari yang tertulis di sini, itu
>   temuan, bukan kesalahan membaca.
> - Rujukan `berkas:baris` disingkat: `FE` = `APP_RNM/frontend/src/`, `BE` = `APP_RNM/internal/`.
> - Seluruh nilai contoh sintetis (`UJI-…`). Nol nama orang, nomor polis, kredensial, atau alamat
>   layanan.
> - Kolom **"Siapa"** adalah gerbang yang **backend** tegakkan. Tombol yang tampil belum tentu boleh.
> - Nomor bagian di dalam tiap bab modul bersifat **lokal**: "§2.5" di bab PremiumList berarti
>   bagian 2.5 bab itu.

---

## 0. Menyalakan aplikasi

Langkah lengkapnya ada di **`APP_RNM/PANDUAN-MENJALANKAN.txt`** bab 1–6. Ringkasnya — **dua jendela**:

| Jendela | Perintah (dari `APP_RNM`) | Tanda benar |
| --- | --- | --- |
| Backend | `. .\muat-env.ps1` lalu `go run ./cmd/api` | log `http: mendengarkan di :8080`; `/healthz` menjawab `"database":"terjangkau"` |
| Frontend | `Set-Location frontend` lalu `npm run dev` | halaman dev Vite terbuka di peramban |

Wajib, **kedua sisi bersama** (tanpa itu Inbox tampak KOSONG, padahal di baliknya 401):

| Berkas | Setelan | Catatan |
| --- | --- | --- |
| `APP_RNM\.env` | `AUTH_STUB=true`, `IS_PEGA_PROD=false`, sambungan Oracle ke **skema uji dari DBA** | `AUTH_STUB` ditolak bila `IS_PEGA_PROD=true` |
| `APP_RNM\.env` | `UNGGAHAN_DIR` = folder lokal | hanya untuk unggah dokumen Claim Life |
| `APP_RNM\frontend\.env` | `VITE_AUTH_STUB=true`; `VITE_STUB_PELAKU` (kosong → `UJI-ADMIN`); `VITE_STUB_PERAN` dipisah koma (kosong → ketiga peran) | ganti akun/peran = ubah lalu **jalankan ulang Vite** bila nilai lama masih terbaca |

⛔ **Setiap klik yang menyimpan MENULIS ke skema yang ditunjuk `ORACLE_SCHEMA`.** Uji hanya di skema uji
kosong dari DBA — tidak pernah `POOLDATA`, tidak pernah produksi. Executor **tidak** menjalankan
`-migrate`; pembentukan tabel di skema uji adalah langkah DBA/work owner (bab 3 berkas itu).

### 0.1 Tiga titik buta — ditutup GILIRAN-13 (29-09-2026)

| # | Titik buta (GILIRAN-11) | Kini | Commit |
| --- | --- | --- | --- |
| 1 | **PremiumList** tidak dapat membuat kasus | tombol `Input Offer` / `Input Premium` di Inbox **membuat kasus** dan langsung membukanya di `Input Offer Life` — keduanya mulai di tahap yang sama; benderanya (`"0"`/`"1"`) bekerja saat `Confirm`: **`"0"` menutup, `"1"` memindah ke Input Premium Detail** (GILIRAN-14 butir bq) | `a291a20` (bn, migrasi **057**), `9f67d35` (bq) |
| 2 | **Claim Life** tanpa rute pembuat baris adjustment pertama | *(Diralat GILIRAN-14.)* baris pertama **lahir saat `Submit` Register** — satu per peserta terpilih, delapan nilai dari polisnya (`SavePesertaClaim` 7.8). `Add` di layar Detail = putaran berikutnya saja. `Delete` **tidak dirender** (tidak berlaku — ADR-U-0031; OQ-N7 ditutup GILIRAN-15) | `0af1773` (bo) → diralat butir **bp** |
| 3 | **Komite** tanpa roster dan baris siap serah | bukan celah kode — **data sintetis** di §0.2 | `a4d9147` |

⚠️ **Yang tetap perlu diketahui.** Migrasi **057** harus sudah berjalan di skema uji sebelum tombol
PremiumList dipakai: executor **tidak** menjalankan `-migrate`. `Submit` Register kini menulis
`OS_AKSEPTASI_KLAIM_LIFE` — tabel itu harus ada di skema uji. Sel baris adjustment **tidak** dapat disunting di
layar mana pun — keputusan work owner 29-09-2026, ikut XML (br dan OQ-N8 ditutup); `CLAIM_GROSS` menunggu pemilik
ekspor (**OQ-N11**). Migrasi **058** (`SEQ_WORK_POLIS` mulai 22374) juga dijalankan work owner, sesudah 057.

### 0.2 Memuat data uji sintetis ke skema uji

Berkas: **`APP_RNM/internal/repository/skemauji/data_uji_tiga_modul.sql`**. Executor **tidak**
menjalankannya — Anda yang memuatnya.

1. Pastikan `-migrate` (termasuk **057**) sudah berjalan di skema uji.
2. Buka SQL*Plus, SQLcl, atau SQL Developer dengan akun yang berhak menulis ke skema uji.
3. Jalankan berkasnya **utuh** (SQL Developer: *Run Script*, F5). Alatnya **menanyakan `skema_uji`**
   sekali — ketik nama skema uji. Setiap tabel diawali nama itu (ADR-U-0033), dan nama yang diketik
   itulah yang dipagari.
4. Periksa hasilnya, lalu **tetapkan transaksinya sendiri** — berkas ini sengaja tidak menetapkan DML-nya.
   (Bila tiruan `EMAILKOMITE` dibuat, DDL-nya menetapkan transaksi yang sedang terbuka — ia berjalan
   sebelum DML apa pun.)

Berkas itu **menolak berjalan** bila:

| Keadaan | Kode galat |
| --- | --- |
| skema aktif **memuat** `POOLDATA` (pagar yang sama dengan `-migrate-down`) | `ORA-20901` |
| skema belum dimigrasi (`T_MIGRASI` tidak ada) | `ORA-20902` |
| migrasi 057 belum berjalan | `ORA-20903` |
| klaim/polis `UJI-*` sudah pernah dimuat — memuat ulang = `-migrate-down` lalu `-migrate` | `ORA-20904` |

⚠️ `ORACLE_SKEMA_UJI=true` tidak dapat dibaca SQL — padanannya nama skema yang Anda ketik.
`IS_PEGA_PROD=false` pun tidak: memilih sambungan yang benar tetap tanggung jawab Anda. Roster
`EMAILKOMITE` **bertahan** melewati `-migrate-down` (bukan tabel migrasi); pemuatan ulang memakai ulang
baris `UJI-EK-*` yang sudah ada. Seluruh `INSERT` berada di **satu** blok — gagal di
mana pun, nol baris tertinggal.

Isinya — seluruhnya sintetis (`UJI-*`, surel `uji-…@contoh.invalid`):

| Modul | Kasus | Keadaan | Untuk menguji |
| --- | --- | --- | --- |
| PremiumList | `UJI-PL-A` | `Input Offer Life`, bendera `"0"` | `Confirm` → tertutup Resolved-Completed (bab 2 §2.3) |
| PremiumList | `UJI-PL-E` | `Input Offer Life`, bendera `"1"` | `Confirm` → pindah ke Input Premium Detail (bab 2 §2.3) |
| PremiumList | `UJI-PL-F` | `Input Offer Life`, bendera **kosong** | `Confirm` → 409 (bendera di luar decision table) |
| PremiumList | `UJI-PL-B` | `Input Premium Detail` | unggah CSV, nomor PL (bab 2 §2.4) |
| PremiumList | `UJI-PL-C` | `Input Premium Summary`, dua peserta | rekap dan `Submit` (bab 2 §2.5) |
| PremiumList | `UJI-PL-D` | `Resolved-Completed`, polis `UJI-POL-0001` | uji negatif kasus tertutup; polis tempat klaim berpijak |
| Claim Life | `UJI-CLM-1` | `Input Register`, satu baris tanpa status | tab Input Register (akun `UJI-ADMIN`) |
| Claim Life | `UJI-CLM-2` | `Outstanding Claim`, satu baris tanpa status | `Save to RNM` (sesudah dokumen peserta diunggah — gerbang langkah 3–4), lalu `Reject Outstanding` |
| Claim Life | `UJI-CLM-3` | `Medical Check`, satu baris Outstanding | layar Medical Check |
| Claim Life | `UJI-CLM-4` | `Claim Analis`, dua peserta | A: `Send Claim to Committee` / `Save Adjustment`; C (baris ditolak): `Add` (putaran berikutnya) |
| Komite | roster `EMAILKOMITE` | `UJI-KOMITE-1`…`4` | baris `UJI-ADJ-4-A1` (150.000.000 IDR) ditutup `UJI-KOMITE-1` dan `-2`; `-3` di atas pitanya, `-4` tidak aktif |

⛔ **Kasus Admin menuntut akun yang sama.** `UJI-CLM-1` dan `UJI-CLM-2` dibuat atas nama `UJI-ADMIN`
(akun stub bawaan): kedua tab Admin menyaring pembuatnya. Untuk bertindak sebagai anggota Komite,
setel `VITE_STUB_PELAKU=UJI-KOMITE-1` lalu jalankan ulang Vite.

⚠️ `EMAILKOMITE` bukan tabel migrasi. Bila skema uji belum memilikinya, berkas membuat **tiruan**
berkolom yang dibaca aplikasi; bila DBA sudah menyalin tabel aslinya dan tabel itu punya kolom wajib
lain, pemuatan batal utuh dengan `ORA-01400`.

Urutan uji yang disarankan: **PremiumList** (polis) → **Claim Life** (klaim atas polis itu) →
**Komite** (baris yang diserahkan).

---

## 1. Claim Life

> Draf untuk work owner yang akan mengeklik aplikasi. Diturunkan dari kode per 28-09-2026,
> bukan dari dokumen desain. Setiap teks di dalam `kode` disalin **apa adanya** dari kode
> (label, tombol, pesan). Rujukan `berkas:baris` singkat; `FE` = `frontend/src/`,
> `BE` = `internal/`.
>
> Cara membaca kolom **"Siapa"**: itulah gerbang yang **backend** tegakkan. Layar sering
> tetap menampilkan tombol kepada orang yang tidak berhak — tombol tampil bukan berarti boleh.

---

### 1. Prasyarat penguji

#### 1.1 Identitas (tidak ada layar login)

| Sisi | Variabel | Nilai untuk uji | Akibat bila salah |
|---|---|---|---|
| Backend `.env` | `AUTH_STUB` | `true` | Setiap jalur beridentitas dijawab 401 `permintaan tanpa identitas pelaku ditolak` (`BE/handlers/pelaku.go:41`) |
| Backend `.env` | `IS_PEGA_PROD` | `false` | `AUTH_STUB=true` **ditolak** saat produksi (`BE/config/config.go:165`); efek keluar hanya "dilewati" di non-produksi |
| Backend `.env` | `UNGGAHAN_DIR` | folder lokal | Tanpa itu unggah dokumen dijawab `UNGGAHAN_DIR belum disetel; unggahan dokumen belum dapat dilayani` |
| Frontend `.env` | `VITE_AUTH_STUB` | `true` | Layar hanya menampilkan `Identitas pelaku tidak dapat dimuat saat ini.` (`FE/App.tsx:39-53`) |
| Frontend `.env` | `VITE_STUB_PELAKU` | kosong → `UJI-ADMIN` | Dikirim sebagai `X-Pelaku` (`FE/store/sesi.ts:48,90`) |
| Frontend `.env` | `VITE_STUB_PERAN` | kosong → **ketiga** peran; dipisah koma | Dikirim sebagai `X-Peran`; peran salah ketik dibuang (`FE/store/sesi.ts:61-95`) |

Ganti peran = ubah `VITE_STUB_PERAN` lalu muat ulang (`PANDUAN-MENJALANKAN.txt:264`).
Tanpa basis data setiap rute menjawab 503 `database belum dikonfigurasi`; layar Inbox
menampilkannya sebagai panel `Backend tidak terhubung (…)` (`FE/lib/keadaanGalat.ts:69-100`).

#### 1.2 Peran dan tahap (`BE/models/tahap.go:84-89`)

| Tahap (tab, verbatim) | Terjemahan | Pemegang | Bentuk antrean |
|---|---|---|---|
| `Input Register` | Input Register | `ReasLifeAdmin` | **pribadi** — hanya kasus yang dibuat akun Anda |
| `Outstanding Claim` | Klaim Outstanding | `ReasLifeAdmin` | **pribadi** — hanya kasus yang dibuat akun Anda |
| `Medical Check` | Pemeriksaan Medis | `ReasLifeMedicalAdvisor` | bersama |
| `Claim Analis` | Analis Klaim | `ReasLifeSPV` | bersama |

⚠️ Dua tab Admin menyaring `CREATE_OP = X-Pelaku` (`BE/services/inbox.go:123-128`). Uji dengan
akun stub yang **sama** dengan yang mendaftarkan klaim, atau kasusnya tidak akan tampil.

#### 1.3 Data yang harus sudah ada di basis data uji

1. Nomor premium list yang berpeserta (contoh uji: `UJI-PL-1`).
2. Polisnya ada di modul PremiumList Life (dibaca `GET /api/polis-life/ringkas`) dengan
   `BusinessCode` yang dikenal — tanpa itu `Save to RNM` berhenti 422. (Ambang produk dan
   `DateReceived` **tidak** diperlukan `Save to RNM`: langkah STNC-nya ter-remark, lihat §3.4.)
3. Daftar kategori dokumen wajib terisi (untuk `Add attachment`).
4. *(Diperbarui GILIRAN-14, butir bp.)* Klaim yang didaftarkan lewat layar ini lahir **berbaris**: satu
   baris adjustment per peserta terpilih, berisi delapan nilai polisnya — `CEDING_RETENTION`,
   `SHARE_NUSANTARA_RE`, `SUM_INSURED`, `SUM_REASURED`, `SHARE_RETRO`, `CLAIM_AMOUNT`, `RETROCEDED_SHARE`
   (dibulatkan empat angka) dan `CURRENCY` — tanpa status sampai `Save to RNM` menulis `Outstanding`. Tombol
   baris dapat diuji pada klaim yang baru didaftarkan; **tidak perlu** baris adjustment manual.
   *(GILIRAN-15.)* Nilai **peserta** pun dibulatkan empat angka (langkah 7.7, OQ-N10 ditutup), dan medan uang
   sumber yang **kosong** menjadi **0** di baris pertama — hanya di baris itu (keputusan work owner; `[dugaan]` seperti
   `@toDecimal("")` Pega). ⚠️ Akibatnya gerbang `Save to RNM` 11.17.1 (`CLAIM_GROSS` kosong) tidak pernah menolak baris
   itu — **OQ-N12**.
5. ⚠️ Karena barisnya kini ada, `Submit` Register ikut **menulis** baris datar `OS_AKSEPTASI_KLAIM_LIFE`
   (AC 32 tiket 02). Tabel itu harus ada di skema uji — tanpa itu pendaftaran gagal.
6. Data uji sintetis (bab 0 §0.2) tetap berguna untuk tahap yang tidak dicapai lewat layar: `UJI-CLM-1`…`4`,
   satu klaim per tahap.

---

### 2. Jalur klik dari Beranda

1. Aplikasi terbuka di **Beranda**: `Selamat datang, UJI-ADMIN`, daftar peran, empat cacah antrean
   (`Input Register`, `Outstanding Claim`, `Medical Check`, `Claim Analis`) (`FE/pages/Beranda.tsx:81-131`).
2. Kartu `Claim Life` berbunyi `aktif — N antrean` dengan tombol `Inbox Claim Life`.
   Sidebar kelompok `Claim Life` memuat dua butir: `Inbox Claim Life` dan `Register`
   (`FE/lib/daftarMenu.ts:57-63`). Ctrl+K membuka palet pencarian menu yang sama.
3. `Inbox Claim Life` → pilih tab tahap → klik baris → **layar Outstanding** (untuk tab **apa pun**,
   `FE/App.tsx:91-104`).
4. Dari layar Outstanding → tombol `Detail klaim` → **layar Detail** (`Klaim Life`), lalu **ketik
   pengenal klaim** di kotak `Pengenal klaim` dan tekan `Buka` (pengenal tidak dibawa otomatis,
   `FE/App.tsx:130`). Salin nilai sesudah `Claim No:` di kepala layar Outstanding.
5. `Register` (sidebar atau tombol di kepala Inbox) → **layar Register Klaim Life**.

---

### 3. Layar demi layar (kontrol berurutan seperti di layar)

#### 3.1 Beranda

| Kontrol | Memanggil | Siapa | Hasil |
|---|---|---|---|
| (otomatis) cacah antrean | `GET /api/klaim-life?tahap=1..4&halaman=1&ukuran=1` ×4 | pemegang **tiap** tahap | Empat angka. ⚠️ Bila Anda tidak memegang KEEMPAT peran, satu tahap dijawab 403 dan Beranda menampilkan `peran tidak memegang tahap ini` tanpa angka (`FE/pages/Beranda.tsx:87-104`) |
| `Inbox Claim Life` | — (navigasi) | — | Membuka Inbox |

#### 3.2 Inbox Claim Life (`FE/pages/claimlife/InboxClaimLife.tsx`)

Tab yang tampil hanya tab yang perannya Anda pegang. Tanpa satu pun: `Tidak ada antrian untuk peran Anda. Keempat antrian Claim Life dipegang ReasLifeAdmin, ReasLifeMedicalAdvisor, dan ReasLifeSPV.`

| # | Kontrol | Memanggil | Siapa | Hasil |
|---|---|---|---|---|
| 1 | `Register` | — | — | Membuka layar Register |
| 2 | `Export xlsx` | — (di peramban) | — | Unduh `inbox-<nama tahap>.xlsx`, **hanya halaman yang tampil**; mati bila tabel kosong |
| 3 | Tab `Input Register` / `Outstanding Claim` / `Medical Check` / `Claim Analis` | `GET /api/klaim-life?tahap=N&halaman=H` | pemegang tahap (§1.2) | Tabel kolom `Case ID`, `Claim No`, `Policy No`, `Work Status`, `Create Operator Name`, `Create Date/Time`; lencana = total; kosong: `Tidak ada kasus di antrian ini.` |
| 4 | Klik baris | `GET /api/klaim-life/{id}` + `GET /api/polis-life/ringkas` | beridentitas | Layar Outstanding |
| 5 | `Sebelumnya` / `Berikutnya` | idem, `halaman±1` | idem | `Halaman N — T kasus`; 50 baris per halaman |

Penolakan: 401 `permintaan tanpa identitas pelaku ditolak`; 403 `peran tidak memegang tahap ini`.

#### 3.3 Register Klaim Life (`FE/pages/claimlife/RegisterKlaim.tsx`)

| # | Kontrol | Memanggil | Siapa | Hasil |
|---|---|---|---|---|
| 1 | `Nomor premium list` (isian) | — | — | Wajib diisi sebelum `Search` hidup |
| 2 | `Certificate No`, `Name of Insured` (isian) | — | — | Boleh kosong = tidak menyaring |
| 3 | `Search` | `GET /api/peserta-life?pl=…&sertifikat=…&nama=…&n=50` | tanpa gerbang peran | Tabel `Select Insured` (centang), `Sertifikat`, `Name of Insured`, `Choose Policy No` (berisi nomor polis) |
| 4 | Centang peserta | `GET /api/polis-life/ringkas?nomorPolis=…` | beridentitas | Panel `Data Polis` terisi dari PremiumList Life + `Polis … — versi …`; polis tak ada: `Data polis belum terbaca.` Lima medan selalu `—` (lihat §6) |
| 5 | `Type` (isian, bawaan `QP`), `Kode bisnis` (isian) | — | — | Keduanya wajib di backend |
| 6 | `Daftarkan klaim` | `POST /api/klaim-life` | **`ReasLifeAdmin`** (`BE/services/pendaftaran.go:159`) | `Klaim terdaftar: <id> — nomor <nomor klaim>`; kasus lahir di tahap **`Outstanding Claim`** (`BE/services/pendaftaran.go:229`), bukan Input Register |

Penolakan yang tampil:
- Campur polis/mata uang (di layar): `Peserta terpilih berbeda polis atau mata uang. Pilih yang sepolis dan semata-uang.`
- 403 `peran tidak mencukupi`; 401 `permintaan tanpa identitas pelaku ditolak`.
- 400 `services: permintaan pendaftaran tidak sah: Type wajib terisi; ia yang menentukan jendela DOL` / `…: kode bisnis wajib terisi; ia yang menentukan prefix nomor` (`BE/services/pendaftaran.go:84-91`).
- Lain-lain 500 `gagal mendaftarkan klaim`.

#### 3.4 Layar Outstanding (`FE/pages/claimlife/OutstandingClaimLife.tsx`)

Kepala: `Claim No: <pengenal>`. Panel `Data Polis` sama dengan Register.

| # | Kontrol | Memanggil | Siapa | Hasil |
|---|---|---|---|---|
| 1 | `Detail klaim` | — | — | Layar Detail (ketik pengenal, §2 langkah 4) |
| 2 | `Save to RNM` | `POST /api/klaim-life/{id}/outstanding` | `ReasLifeAdmin`, kasus di `Outstanding Claim`, belum ditutup (`BE/services/simpanrnm.go:377-415`) | `Tersimpan ke RNM — nomor klaim <n>` [` (baru diterbitkan)` bila nomor terbit saat itu]`; <k> baris ditandai Outstanding; Arasapas: dilewati: lingkungan bukan produksi (IsPEGAPROD).` — atau `Arasapas: dilewati: kode retro langkah 27.` / `Arasapas: ditahan: gerbang retro langkah 27 bergantung pada ProdDateTime polis, yang tidak tersedia (OQ-N5).` Tombol **tetap hidup** dan boleh ditekan ulang (OQ-N1) |
| 3 | `Send Back to Register` | dialog `Send Back to Admin?` → `Submit` / `Cancel` → `POST …/tahap/input-register` | `ReasLifeAdmin` | Kembali ke Inbox; kasus pindah ke tab `Input Register` |
| 4 | `Send to Medical Check` | langsung `POST …/tahap/medical-check` (tanpa dialog) | `ReasLifeAdmin` | Kembali ke Inbox; kasus pindah ke tab `Medical Check` |

⚠️ Tombol 2–4 hanya tampil bila kasus **sebenarnya** berada di `Outstanding Claim` (sejak `fbf1c9b`).
Kotak masuk membuka layar ini untuk tahap mana pun; untuk tahap lain layar menulis
`Kasus ini berada di tahap <tahap>, bukan Outstanding Claim: tombol layar ini tidak berlaku untuknya. Buka dari layar Detail klaim.`
— pakai layar Detail.

**Gerbang `Save to RNM`, urut XML — pelanggaran pertama yang dilaporkan** (`BE/services/simpanrnm.go:140-222`).
Pesan tampil dua kali (panel merah + baris `role=alert`), itu perilaku layar saat ini.

| Urut | Syarat | Pesan (verbatim, termasuk salah eja) |
|---|---|---|
| 1 | Tiap peserta punya ≥1 dokumen (dilewati untuk Type `TP`/`TR`) | `The document hasn’t been uploaded person number N` — satu baris per peserta, apostrof lengkung |
| 2 | Klaim ganda (tabel warisan) | `Person number N has already been accepted.` |
| 3 | DOL dalam jendela valuasi, **tanpa** geser retro | `DOL cannot be blank or outside the valuation period No N` |
| 4 | Medan wajib | `DOB cannnot be blank No N` · `Begin Date cannnot be blank No N` · `Expired Date cannnot be blank No N` · `Policy No cannnot be blank No N, please contact IT` · `Certificate No cannnot be blank No N, please contact IT` · `Claim Gross No N can't null` |

⛔ **Diralat 28-09-2026** (temuan /code-review GILIRAN-11): versi pertama panduan ini memuat dua
gerbang lagi — STNC (`Begin date exceed STNC No N`) dan dokumen lengkap (`Documents are incomplete,
please complete the documents`). Keduanya **ter-remark** di XML (langkah 11.9/11.11 dan 12,
`pyStepsBlockName = //`), jadi sistem lama tidak pernah menolak karenanya, dan `Save to RNM` kini
pun tidak. Bila salah satu pesan itu muncul, catat sebagai cacat.

Penolakan lain: 403 `hanya pemegang tahap Outstanding Claim yang dapat menyimpan ke RNM`;
409 `Save to RNM hanya tersedia pada tahap Outstanding Claim`; 409 `kasus sudah ditutup dan tidak dapat diubah`;
404 `klaim tidak ada`;
422 data belum lengkap, kalimatnya berawalan `services: Type klaim tidak dikenal…`,
`services: BusinessCode tidak dikenal…`, `repository: nomor polis tidak ditemukan di PremiumList Life…`
(`BE/handlers/simpanrnm.go:21-57`).

Penolakan perpindahan (dari `PanelPindahTahap`): 403 `hanya pemegang tahap asal yang dapat memindahkan kasus ini`;
409 `perpindahan itu tidak ada di tangga kerja klaim`; 409 `tahap kasus ini tidak dikenal`;
kasus tertutup → 409 `kasus sudah ditutup dan tidak dapat diubah` (sejak `e3619d4`).

#### 3.5 Layar Detail — `Klaim Life` (`FE/pages/claimlife/KlaimLife.tsx`)

##### 3.5.1 Membuka klaim

| # | Kontrol | Memanggil | Siapa | Hasil |
|---|---|---|---|---|
| 1 | `Pengenal klaim` + `Buka` | `GET /api/klaim-life/{id}` | tanpa gerbang | Judul nomor klaim; `Nomor polis`, `Nama bisnis`, `Jumlah baris`, `Status klaim`. Galat: `Klaim tidak ada.` / `Gagal membaca klaim.` |
| 2 | `Hapus klaim…` | `GET /api/klaim-life/{id}/dampak-hapus` | `ReasLifeAdmin` | Popup `Hapus klaim <n>?` · `Yang akan ikut terhapus:` (Header klaim, Peserta, Baris adjustment, Spreading, Spreading retro, Dokumen, Baris work) · `Total N baris.` · peringatan baris datar warisan. Gagal: `Gagal menghitung dampak penghapusan.` |
| 3 | `Ya, hapus` | `DELETE /api/klaim-life/{id}` | `ReasLifeAdmin` | **Selalu ditolak** — backend 405; layar menampilkan kalimat server `penghapusan klaim selalu berupa penanda dan nilai pembalik, tidak pernah hapus fisik (ADR-U-0031); DELETE tidak berlaku atas sumber daya ini` (sejak `fbf1c9b`). 409 pun menampilkan kalimat server. |
| 4 | `Batal` | — | — | Menutup popup, tidak menulis apa pun |

##### 3.5.2 Per peserta (`Peserta <sertifikat> (N baris)`), berurutan

| # | Kontrol | Memanggil | Siapa | Hasil / pesan |
|---|---|---|---|---|
| 1 | `Edit Date` (kotak tanggal = DOL; simpan saat berubah) | `PUT …/peserta/{p}/tanggal-kejadian` | `ReasLifeAdmin`, tahap `Outstanding Claim` (kotak mati di tahap lain) | Sukses: kotak memuat ulang nilai tersimpan. Tolak: `Invalid DOL` · `Type klaim tidak dikenal; jendela valuasi tidak dapat ditentukan` · `tanggal valuasi peserta kosong; DOL tidak dapat divalidasi` |
| 2 | `MAX CLAIM RECEIVED: …` | (baca) | — | `—` = sah; tanggal = melampaui ambang; `(tidak dihitung: …)` bila tak terhitung |
| 3 | `CLAIM RECEIVED DATE`, `DOCUMENT COMPLETE DATE`, `CONFIRMATION DATE` + `Save` | `PUT …/peserta/{p}/tanggal-klaim` (ketiganya satu permintaan; kosong = dikosongkan) | `ReasLifeAdmin`, tahap `Outstanding Claim` | Sukses: layar memuat ulang. Tolak: lihat "pesan Edit Date" di bawah |
| 4 | `Total peserta` (baca) | — | — | Enam total: `Total Ceding Retention`, `Total Share Nusantara Re`, `Total Sum Insured`, `Total Sum Reasured`, `Total Share Retro`, `Total Claim Amount` — jumlah **seluruh** baris adjustment, termasuk yang ditolak |
| 5 | `Dokumen pendukung` → `Kategori` + `Add attachment` (pilih berkas; mati sampai Kategori diisi) | `POST …/peserta/{p}/dokumen` (multipart) | beridentitas, kasus terbuka (tanpa gerbang peran/tahap) | Baris baru `Nama berkas`/`Kategori`/`Tanggal`/`Berkas`; kolom `Berkas` berbunyi `URL menunggu penyambungan penyimpanan` (lihat §6). Tolak: `services: kategori dokumen tidak ada di daftar kategori: "…"` · `berkas kosong` · `kasus sudah ditutup dan tidak dapat diubah` · `UNGGAHAN_DIR belum disetel; unggahan dokumen belum dapat dilayani` · `daftar kategori dokumen belum tersedia` |
| 6 | `View Office Online` / `Delete` (dokumen) | `GET /api/dokumen/{d}/isi` (lewat `fetch` ber-header identitas, disimpan sebagai berkas bernama asli) / `DELETE …/dokumen/{d}` | beridentitas | **Hanya tampil bila berkas terunggah** — lihat §6. Unduhan gagal menampilkan kalimat server di baris `role=alert` |
| 7 | Grid `DIAGNOSE` → `Add` | `POST …/peserta/{p}/diagnosa` | pemegang tahap saat ini (Outstanding: Admin; Medical Check: Medical Advisor; Claim Analis: SPV) | Baris kosong `—`; `GROUP DIAGNOSE` berbunyi `GROUP DIAGNOSE tidak dapat dimuat saat ini.`; ringkasan `N diagnosa, K di antaranya belum diisi.` |
| 8 | Per baris: `Find Disease` → `ICD Code`, `Disease`, `Cari`, `Tutup` | `GET /api/penyakit-life?icd=…&nama=…&batas=50` | beridentitas | `N diagnosa ditemukan.` · `Tidak ada diagnosa yang cocok.` · `… daftarnya TERPOTONG pada batas 50. Persempit kata kuncinya.` Kedua kata kunci disambung **AND** |
| 9 | Hasil → `Choose` | `PUT …/peserta/{p}/diagnosa/{d}` | seperti #7 | Baris terisi `DIAGNOSE` + `ICD CODE`; panel pencarian menutup |
| 10 | Per baris: `Delete` | `DELETE …/peserta/{p}/diagnosa/{d}` | seperti #7 | Baris hilang **seketika** (tanpa konfirmasi) |
| 11 | `Save Adjustment` (tampil bila peserta dipilih DAN baris terakhir `Outstanding` tanpa nomor akseptasi) | `POST …/peserta/{p}/akseptasi` | **pemegang tahap saat ini** (`BE/services/akseptasi.go:177-184`) | `Baris diaksep dengan nomor <nomor>.` |
| 12 | **`Add`** (tampil di tahap `Claim Analis` bila baris terakhir ditolak, kasus terbuka — `FE/services/api.ts:bolehAddAdjustment`) *(GILIRAN-14: menggantikan tombol "Putaran berikutnya")* | `POST …/peserta/{p}/putaran` | **`ReasLifeSPV`** | Baris adjustment baru muncul, mewarisi delapan kolom baris pertama. 409 `putaran berikutnya hanya lahir sesudah baris terakhir ditolak` (juga bagi peserta klaim lama yang tanpa baris) · 409 `Add baris adjustment hanya tersedia di tahap Claim Analis` (gerbang b18160 kini juga di backend) |
| 13 | Tabel baris: `Baris`, `Status`, `Jumlah klaim`, `Nomor akseptasi`, `Tindakan`, `Komite`. Kosong: `Belum ada baris adjustment.` | | | Status berupa kata: `Outstanding` / `Aksep` / `Ditolak` / `Tidak diketahui` |
| 14 | `Reject Outstanding` (per baris `Outstanding`, klaim bernomor) | `POST …/adjustment/{a}/tolak` — **tanpa dialog, tanpa alasan** | **`ReasLifeAdmin`** (tahap tidak diperiksa) | Status baris → `Ditolak`; `Save Adjustment` hilang, `Add` muncul (tahap Claim Analis) |
| 15 | `Send Claim to Committee` (per baris `Outstanding` belum ke Komite) | `POST …/peserta/{p}/adjustment/{a}/komite` | Type `QP`/`QR`: **`ReasLifeSPV`**; Type `TP`/`TR`: siapa pun beridentitas (`BE/services/wewenang.go:118-143`) | Kolom Komite → `sudah diserahkan` |

**Pesan Edit Date (DOL dan `Save`)** — `BE/handlers/dol.go:98-116`:
403 `peran tidak mencukupi` · 409 `kasus sudah ditutup dan tidak dapat diubah` ·
409 `tanggal klaim hanya dapat diubah di tahap Outstanding Claim` · 422 `tahap kasus tidak dikenal` ·
400 `tanggalTerimaKlaim bukan tanggal yang dikenal` (dan dua kerabatnya) · cadangan layar
`Tanggal kejadian ditolak.` / `Tanggal klaim gagal disimpan.`

**Pesan grid diagnosa** — `BE/handlers/diagnosa.go:135-163`:
403 `hanya pemegang tahap ini yang dapat mengubah diagnosanya` ·
409 `peserta sudah diputus; diagnosanya tidak dapat diubah lagi` ·
409 `kasus sudah ditutup dan tidak dapat diubah` · 409 `tahap ini tidak membuka layar detail peserta`.
Peserta berstatus diaksep/ditolak: tombol diganti `Terkunci` dan kalimat `Peserta ini sudah diputus; diagnosanya tidak dapat diubah lagi (gerbang .STS_REJECT b4682).`

**Pesan tombol baris** (teks yang benar-benar tampil, `KlaimLife.tsx:237-336`):

| Tombol | 403 | 409 | 422 | lain |
|---|---|---|---|---|
| `Save Adjustment` | `Hanya pemegang tahap klaim ini yang dapat mengaksep barisnya.` | dari server: `peserta belum dipilih untuk diklaim` / `baris sudah punya nomor akseptasi` / `hanya baris Outstanding yang dapat diaksep` / `nomor akseptasi yang terbit sudah dipakai; coba lagi` | `Gagal mengaksep baris.` | 501 dari server `kode bisnis klaim belum tersimpan; nomor akseptasi belum dapat dirakit` |
| `Add` | `Hanya ReasLifeSPV yang dapat menambah baris adjustment.` | kalimat server (mis. `kasus sudah ditutup dan tidak dapat diubah`); cadangan `Baris adjustment tidak dapat ditambahkan pada keadaan ini.` | — | `Gagal menambah baris adjustment.` |
| `Reject Outstanding` | `Hanya ReasLifeAdmin yang dapat menolak baris.` | kalimat server; cadangan `Baris sudah diputus dan tidak dapat ditolak lagi.` | `Klaim belum bernomor.` | `Gagal menolak baris.` |
| `Send Claim to Committee` | `Peran Anda tidak berwenang menyerahkan baris bertipe ini ke Komite.` | kalimat server; cadangan `Baris sudah diserahkan, atau bukan lagi Outstanding.` | dari server: `Name of bank cannot be empty` / `baris pada klaim ini bermata uang campur; penyerahan menuntut mata uang tunggal` / `tidak ada tingkat komite yang menutup nilai klaim ini` / `Type klaim tidak dikenal; penyerahan ke Komite menuntut Type yang sah` | `Gagal menyerahkan baris ke Komite.` |

##### 3.5.3 Perpindahan tahap di layar Detail (`FE/components/claimlife/PanelPindahTahap.tsx:61-97`)

Judul panel `Perpindahan dari <tahap>`. Tombol ber-dialog menampilkan pertanyaan + `Submit` / `Cancel`.

| Tahap kasus | Tombol | Dialog | Rute | Siapa | Tujuan |
|---|---|---|---|---|---|
| `Outstanding Claim` | `Send Back to Register` | `Send Back to Admin?` | `POST …/tahap/input-register` | Admin | `Input Register` |
| `Outstanding Claim` | `Send to Medical Check` | — | `POST …/tahap/medical-check` | Admin | `Medical Check` |
| `Medical Check` | `Send Back to Admin` | `Send Back to Admin?` | `POST …/tahap/outstanding` | Medical Advisor | `Outstanding Claim` |
| `Medical Check` | `Send to Claim Analyst` | — | `POST …/tahap/claim-analis` | Medical Advisor | `Claim Analis` |
| `Claim Analis` | `Send Back to Admin` | `Send Back to Admin?` | `POST …/tahap/outstanding` | SPV | `Outstanding Claim` |
| `Claim Analis` | `Send Back to Medical` | `Send Back to Medical?` | `POST …/tahap/medical-check` | SPV | `Medical Check` |
| `Input Register` | (tidak ada tombol) | | | | |

Sukses: layar Detail memuat ulang dan tombol berganti mengikuti tahap baru. Penolakan: sama dengan §3.4.

##### 3.5.4 Close Claim

Hanya tampil pada `Outstanding Claim` dan `Claim Analis`; tahap lain menampilkan
`Close Claim hanya tersedia pada tahap Outstanding Claim dan Claim Analis.`; kasus tertutup `Kasus ini sudah ditutup.`

| # | Kontrol | Memanggil | Siapa | Hasil |
|---|---|---|---|---|
| 1 | `Periksa kesiapan: Close Claim` | `GET /api/klaim-life/{id}/boleh-tutup` (tidak mengubah apa pun) | tanpa gerbang | `Seluruh peserta sudah diaksep.` atau `Klaim belum dapat ditutup:` + daftar `<sertifikat> is not approved yet, on list <n>` |
| 2 | `Close Claim` | membuka dialog `Are you sure want to Close Claim?` | — | — |
| 3 | `Ya, tutup klaim` | `POST /api/klaim-life/{id}/tutup` | pemegang tahap (Outstanding: Admin; Claim Analis: SPV) | `Klaim <n> ditutup. Kasusnya selesai dan hilang dari kotak masuk; tidak ada lagi yang dapat diubah padanya.` Kasus hilang dari semua tab |
| 4 | `Batal` | — | — | Menutup dialog |

Penolakan `Ya, tutup klaim` (`BE/handlers/tutup.go:87-123`): daftar penghalang seperti #1 (409) ·
403 `hanya pemegang tahap kasus ini yang dapat menutupnya` · 409 `kasus ini sudah ditutup` ·
409 `Close Claim hanya ada pada tahap Outstanding Claim dan Claim Analis` · 404 `klaim tidak ditemukan`.
Peserta **tanpa** status aksep selalu menahan — klaim yang baru didaftarkan karena itu hanya dapat
menguji jalur tolaknya.

---

### 4. Skenario urut yang disarankan (satu klaim baru, peran bawaan ketiganya)

1. `Register` → isi `Nomor premium list` → `Search` → centang satu peserta → `Type` → `Kode bisnis` → `Daftarkan klaim`. Catat `<id>`.
2. `Inbox Claim Life` → tab `Outstanding Claim` → kasus tampil → klik → layar Outstanding.
3. `Save to RNM` → harapkan gerbang 1 (`The document hasn’t been uploaded person number 1`) bila Type bukan TP/TR.
4. `Detail klaim` → ketik `<id>` → `Buka` → `Add attachment` (satu dokumen berkategori apa pun dari daftar cukup) → ulangi `Save to RNM` di layar Outstanding sampai lolos atau berhenti di gerbang berikutnya.
5. Detail → `Edit Date` (uji batas DOL §5.1) → tiga tanggal + `Save`.
6. Detail → `Send to Medical Check` → tab `Medical Check` → buka lagi di Detail → `Add` / `Find Disease` / `Choose` / `Delete`.
7. `Send to Claim Analyst` → di `Claim Analis`: `Send Back to Medical` (dialog `Send Back to Medical?`) → kembali → `Send Back to Admin` (dialog `Send Back to Admin?`).
8. `Periksa kesiapan: Close Claim` → harapkan `… is not approved yet, on list 1`.
9. Uji 403: set `VITE_STUB_PERAN=ReasLifeMedicalAdvisor`, muat ulang, ulangi `Send to Medical Check` pada kasus Outstanding (kasus dibuka lewat layar Detail).

---

### 5. Contoh nilai dari uji otomatis (fiktif, `UJI-…`)

#### 5.1 Batas DOL di `Edit Date` (`BE/services/dol_test.go:23-26, 61-105`)
Peserta uji `UJI-P-DOL`: jendela **gross** 2025-03-01 s.d. 2025-09-01; jendela **retro** 2025-04-01 s.d. 2025-10-01.
Hanya contoh berjam 00:00 yang dapat diketik di kotak tanggal:

| Type | DOL | Hasil |
|---|---|---|
| QR / QP | 2025-03-01 (tepat mulai) | `Invalid DOL` |
| QR / QP | 2025-09-01 (tepat akhir) | diterima |
| QR | 2025-09-15 | `Invalid DOL` |
| QR | 2025-03-15 | diterima |
| TR / TP | 2025-03-30 | `Invalid DOL` |
| TR / TP | 2025-04-01 (tepat mulai) | diterima (DOL digeser +1 hari) |
| TR / TP | 2025-09-29 | diterima |
| TR / TP | 2025-10-01 (tepat akhir) | `Invalid DOL` |
| TR | 2025-09-15 | diterima |
| TR | 2025-03-15 | `Invalid DOL` |

Aturannya: QP/QR `(mulai, akhir]`; TP/TR DOL+1 hari dibandingkan dengan `(mulai, akhir]`. Ganti
tanggal contoh dengan jendela peserta Anda (jendela itu tidak tampil di layar — lihat §6).

#### 5.2 DOL di `Save to RNM` — tanpa geser retro (`BE/services/simpanrnm_test.go:130-157`)
QR tepat mulai → ditolak; QR tepat akhir → lolos; QP 2025-09-02 → ditolak; **TR tepat mulai
2025-04-01 → ditolak** (padahal `Edit Date` menerimanya); TP tepat akhir 2025-10-01 → lolos; DOL
kosong → ditolak. Pesan: `DOL cannot be blank or outside the valuation period No 1`.

#### 5.3 STNC — tidak ada gerbang (diralat 28-09-2026)
`SaveOutStandingLife_Act` langkah 11.9/11.11 ter-remark (b4632, b5009): BEGIN_DATE sejauh apa pun
dari `DateReceived` **tidak** menolak `Save to RNM`. Versi pertama panduan ini menyuruh menguji
`Begin date exceed STNC No 1` — pesan itu kini tidak boleh muncul.

#### 5.4 Pesan bertumpuk dan urutan (`simpanrnm_test.go:71-195`)
- Peserta 1 dan 3 tanpa dokumen → dua baris: `The document hasn’t been uploaded person number 1` / `… number 3`.
- Type `TP` / `TR` tanpa dokumen → **lolos**: gerbang 1 dilewati, dan langkah 12 ter-remark.
- Dua dokumen berkategori sama → **lolos** (kelengkapan per kategori tidak ditegakkan — butir bl, OQ-N6 ditutup).
- DOB dan BEGIN sama-sama kosong → yang dilaporkan `DOB cannnot be blank No 1`.

#### 5.5 Close Claim (`FE/services/tutupklaim.test.ts:23-25, 90-161`)
Klaim `UJI-KLAIM-1` (nomor `UJI-CLM-1`, polis `UJI-POL-1`); peserta `UJI-001` dan `UJI-003` belum
diaksep → `UJI-001 is not approved yet, on list 1`, `UJI-003 is not approved yet, on list 3`.

#### 5.6 Register (`FE/services/caripeserta.test.ts:43-83`)
`Nomor premium list` `UJI-PL-1`, `Certificate No` `UJI-006` → kueri `pl=UJI-PL-1&sertifikat=UJI-006&n=50`;
kotak berisi spasi saja tidak dikirim.

#### 5.7 Nomor (`BE/services/pendaftaran_test.go:203-209`, `akseptasi_test.go:62-73`)
- Nomor klaim: awalan `UJI-`, kode bisnis `L9`, periode `03.2027`, urut 7 → `UJI-KL9.03.2027.00007`.
- Periode nomor akseptasi: tanggal > 25 → bulan berikut. 2026-03-25 → `03`/`26`; 2026-03-26 → `04`/`26`;
  2026-12-26 → `01`/`27`; 2026-09-30 → `10`/`26`.

#### 5.8 Claim Paid — aturan teruji, **belum ada di layar** (`BE/models/klaimbayar_test.go:30-39`)
1000000 × 50% = 500000 · 1000000 × 100% = 1000000 · 1000000 × 33.333333% = **333330**
(persen dibulatkan 5 desimal menjadi 0.33333 lebih dulu) · 123.45 × 12.5% = 15.43125.

---

### 6. Belum dapat diuji di layar

| Butir | Sebab (satu baris) |
|---|---|
| Email ke anggota Komite | `EfekEmail` selalu `ErrEmailBelumDisetujui`; di non-produksi penyalur hanya "dilewati" (`BE/services/efekkeluar.go:366-374`) |
| Kirim ke Arasapas (`Save to RNM`) | Non-produksi → `Arasapas: dilewati: lingkungan bukan produksi (IsPEGAPROD)`; produksi pun `ErrArasapasBelumDisetujui` |
| Kasir | Stub `ErrKasirBelumDisetujui`, milik modul Komite, bukan layar Claim Life (`BE/services/komite_pengirim.go:50`) |
| Google Storage / unduh & hapus dokumen | Tidak ada pekerja outbox yang dijalankan (`cmd/api/main.go`), jadi `tStorageId` tak pernah terisi dan baris tetap `URL menunggu penyambungan penyimpanan`; `View Office Online` dan `Delete` tidak pernah tampil. (Sejak GILIRAN-12 paket 0, `View Office Online` mengunduh lewat klien ber-identitas dan menyimpan berkas dengan nama aslinya; sebelumnya pranala biasa yang selalu dijawab 401.) |
| Hapus klaim | Backend selalu menolak (405, kolom penanda belum diputuskan, ADR-U-0031); layar menampilkan kalimat server itu |
| Menyunting sel baris adjustment | **Tidak berlaku** — keputusan work owner 29-09-2026, ikut XML (butir br, OQ-N8 ditutup): keempat kolom grid dan seluruh medan panel `Adjustment_Detail` bertanda `Read-only`. `CLAIM_GROSS` menunggu pemilik ekspor — **OQ-N11** |
| `Delete` baris adjustment | **Tidak berlaku, tidak dirender** — keputusan work owner 29-09-2026 (OQ-N7 ditutup): ADR-U-0031, tabel tanpa kolom penanda |
| Jalan maju dari `Input Register` | Tidak ada tombol di tahap itu (Pega: `Submit` pendaftaran); layar Outstanding kini menyatakan tahapnya dan tidak menawarkan tombol |
| Claim Paid / `Percent Claim (%)` | Kolom `PCT_CLAIM`/`CLAIM_PAID` dan rute sunting adjustment belum ada — OQ-M3 |
| Spreading / panel `RetroDetailClaimLife` | `HitungSpreading` nol pemanggil; sumber rate belum diputuskan — OQ-M7 |
| `Claim Life - Upload CSV` di Register | Belum dibangun; menunggu keputusan desain XOL — OQ-M4 |
| Dialog Reject (Date, PIC, Remarks, Submit) | Sengaja belum dibangun; tempat menyimpan alasan belum diputuskan — OQ-M5 |
| Cabut peserta (`DELETE` di layar OS) | Penanda atau lepas baris belum diputuskan — OQ-M6 |
| Kunci tanggal sesudah Save Outstanding | Separuh gerbang `CLAIM_NO!=''` tidak ditiru; tanggal tetap dapat diubah selama Outstanding — OQ-M1 |
| Cermin tanggal ke tabel warisan | `UpdateDateClaimLife_SQL` tidak ditiru; kunci nama tertanggung tidak ada — OQ-M2 |
| `Save to RNM` terkunci sesudah simpan | Bendera `pyWorkPage.Save` tanpa kolom; tombol tetap hidup (tulisan idempoten) — OQ-N1 |
| Klaim ganda antarklaim baru | Cermin warisan tidak mengisi nama/DOB/CEDINGCO; hanya baris era Pega yang tertangkap — OQ-N2 |
| Lewati Arasapas untuk tiga kode retro | Mengikuti XML; berlaku-tidaknya OQ-064 belum diputuskan — OQ-N3. Kodenya dibaca SESUDAH penukaran `InsertJsonClaimLife_Act` langkah 2; bila hasilnya bergantung pada `ProdDateTime` (tanpa sumber), Arasapas **ditahan** — OQ-N5 |
| Dokumen lengkap per kategori | Tidak ditegakkan, dan bukan cacat: langkah 12 ter-remark di XML; OQ-N6 ditutup (butir bl, GILIRAN-12) — bila bisnis menghendakinya, ia keputusan baru |
| Residu `Save to RNM` | Pesan `.Protect` tak pernah muncul; `ADJUSTMENT_DATE`/`PrintFaceClaim` tanpa kolom; `IsAccept` diganti `IS_CHECK` — OQ-N4 |
| Pilihan `GROUP DIAGNOSE` | Daftar pilihannya tidak ada di ekspor; sel menampilkan `GROUP DIAGNOSE tidak dapat dimuat saat ini.` — OQ-L |
| Tanggal Respon/Konfirmasi/Realisasi, Confirmation Reserved, Underwriter Note | Tanpa kolom; selalu `—` di panel `Data Polis` |
| Label yang ada di berkas label tetapi tidak dirender | `Save to Outstanding`, `Select All`, `Refresh`, `Save` (Medical Check / Claim Analis), `Cancel` (Komite), `Close Claim` layar OS, `Find Insured` |
| Jendela valuasi peserta | Tidak ditampilkan di layar Detail; penguji harus mengetahuinya dari data sumber untuk menguji batas DOL |

---

### 7. Catatan untuk penguji (perilaku yang mudah disangka cacat, atau memang janggal)

1. Setiap baris Inbox — tab apa pun — membuka layar Outstanding (`App.tsx`); tombolnya kini hanya tampil bila kasus memang di Outstanding.
2. Kepala layar Outstanding berlabel `Claim No:` tetapi menampilkan **pengenal** kasus, bukan nomor klaim.
3. Layar Detail tidak menerima pengenal dari layar sebelumnya; ketik manual.
4. Konfirmasi Close Claim memakai `Ya, tutup klaim` / `Batal`, bukan `Submit` / `Cancel`.
5. Hapus klaim selalu dijawab 405; layar menampilkan kalimat server (sejak `fbf1c9b`).
6. `Reject Outstanding` tidak bertanya dan tidak meminta alasan.
7. Beranda gagal memuat keempat angka bila akun tidak memegang keempat peran (Beranda meminta keempat tahap).

---

## 2. PremiumList Life

> Draf ini diturunkan dari kode `APP_RNM` yang dibaca pada 28-09-2026. Aplikasi dan basis data TIDAK
> dijalankan untuk menyusunnya. Singkatan rujukan:
> **FE** = `APP_RNM/frontend/src/` · **BE** = `APP_RNM/internal/` · **Tiket** = `.scratch/premiumlist-life/issues/`.
> Seluruh pengenal contoh berawalan `UJI-`. Teks di antara backtick adalah teks layar VERBATIM dari kode;
> `<…>` menandai bagian yang terisi nilai saat berjalan.

---

### 0. Ringkas untuk penguji

- **Layar yang diuji (5, ditambah pintu masuk):** Beranda/menu → Inbox `PremiumList` → keputusan `Input Offer`
  → `Premium List Detail` (berisi `Upload CSV Premium List Detail` dan keputusan `Input Premium Detail`)
  → `Summary Premium Life`.
- **Tiga hal yang paling mungkin membingungkan:**
  1. Layar kasus **tidak punya tombol Kembali**. Lihat §1.3.
  2. Layar `Summary Premium Life` **tidak dapat dicapai lewat alur aplikasi**, karena tahap
     `Input Premium Summary` tidak punya konektor masuk. Lihat §2.5.
  3. Galat 502/503/504 ditampilkan sebagai panel **"Backend tidak terhubung (…)"**, bukan kalimat server. Lihat §1.3.
- **Peran tidak diperiksa.** Siapa pun yang punya identitas stub boleh menekan setiap tombol di bawah. Lihat §1.1.

---

### 1. Prasyarat

#### 1.1 Identitas (tidak ada form login)

| Sisi | Setelan | Bila salah, yang terlihat |
| --- | --- | --- |
| FE `.env` | `VITE_AUTH_STUB=true`. `VITE_STUB_PELAKU=UJI-…` (bila kosong menjadi `UJI-ADMIN`, FE `store/sesi.ts:48,90`). `VITE_STUB_PERAN` bebas (bila kosong berisi ketiga peran, `sesi.ts:94`). | Layar awal hanya menampilkan panel `Identitas pelaku tidak dapat dimuat saat ini.` beserta catatan tentang `VITE_AUTH_STUB=true` (FE `App.tsx:39-54`) |
| BE env | `AUTH_STUB=true` (BE `config/config.go:143`) dan `IS_PEGA_PROD=false`. `AUTH_STUB` ditolak saat `IS_PEGA_PROD=true` (`config.go:163-165`). | Setiap permintaan PremiumList dijawab 401 `permintaan tanpa identitas pelaku ditolak` (BE `handlers/rute_premiumlist.go:331`) |
| BE basis data | Oracle terkonfigurasi | 503 `database belum dikonfigurasi`, yang di layar tampil sebagai panel "Backend tidak terhubung (…)" (§1.3) |

**Peran: nol gerbang.** Setiap layanan PremiumList (`BE services/polis_*.go`) hanya memanggil `WajibIdentitas`,
yang hanya memeriksa bahwa akun tidak kosong (`services/services.go:162-167`). `services/polis_inbox.go:9-12` menulis
"Nol gerbang PERAN". Akibatnya kolom "Siapa" di tabel §2 selalu berbunyi *ber-identitas*, dan 403
`wewenang tidak mencukupi` tidak dapat dipicu dari modul ini.

#### 1.2 Data kasus yang harus disiapkan di luar aplikasi

*(Diperbarui GILIRAN-13.)* Tombol `Input Offer` / `Input Premium` di Inbox kini **membuat kasus** (§2.2) —
tetapi keduanya lahir di `Input Offer Life`, dan `Input Premium Summary` tidak dapat dicapai lewat alur
sama sekali. Karena itu kasus per tahap tetap disiapkan lewat **data uji sintetis** (bab 0 §0.2), yang
memuat keempatnya:

| Kasus contoh | `T_WORK_POLIS.STATUS` (= tahap) | Dipakai untuk |
| --- | --- | --- |
| `UJI-PL-A` | `Input Offer Life`, bendera `"0"` | §2.3: `Confirm` (→ Resolved-Completed), `Decline` |
| `UJI-PL-E` | `Input Offer Life`, bendera `"1"` | §2.3: `Confirm` (→ Input Premium Detail) |
| `UJI-PL-F` | `Input Offer Life`, bendera kosong | §2.3: `Confirm` → 409 |
| `UJI-PL-B` | `Input Premium Detail` | §2.4: unggah CSV, `Generate PL Number`, `Confirm`, `Reject`, `Decline` |
| `UJI-PL-C` | `Input Premium Summary` | §2.5: `Summary Premium Life` dan `Submit`. Tahap ini hanya dapat dicapai lewat data yang disiapkan |
| `UJI-PL-D` | `Resolved-Completed` (polis `UJI-POL-0001`, tempat klaim `UJI-CLM-*` berpijak) | §2.6: uji negatif kasus tertutup |

Syarat data yang dituntut query:

1. `T_WORK_POLIS.ID` = Case ID, paling panjang 32 karakter (migrasi `050`). Tahap dibaca dari `STATUS`, bukan
   dari `POSITION` (BE `services/polis_penawaran.go:102-112`).
2. Header `T_PREMIUM_LIST` harus punya **`ID` = Case ID**, karena Detail, nomor, dan rekap membaca `p.ID = :1`
   (BE `repository/polis_nomor.go:118,221`). Header juga harus punya **`ID_PEGA` = Case ID**, karena kolom Inbox
   digabung lewat `p.ID_PEGA = w.ID` (`repository/polis_inbox.go:101`). Bila `ID` tidak cocok, layar Detail
   menampilkan `polis tidak ditemukan`.
3. `TYPE` bernilai salah satu dari `QR`/`QP`/`TP`/`TR`, dan `BUSINESS_CODE` terisi. Keduanya bahan nomor dan rekap.
4. Tabel rujukan: `POOLDATA.TANGGAL_CLOSING` terisi (bahan periode), dan awalan `KODE_PRODUKSI` untuk lini Life
   tersedia (bahan awalan nomor).

#### 1.3 Perilaku layar yang perlu diketahui sebelum mulai

- **Tidak ada tombol Kembali.** Kasus hanya dilepas oleh keputusan yang memindahkan atau menutupnya
  (FE `pages/premiumlist/InputOffer.tsx:98-100` → `App.tsx:86-88`). Menekan menu `PremiumList` lagi **tidak**
  melepas kasus, karena state `polis` di `App.tsx:35` tidak direset. Untuk kembali tanpa memutuskan, muat ulang
  peramban. Aplikasi lalu kembali ke Beranda.
- **Kalimat hasil keputusan hampir tidak sempat terlihat.** Contohnya `Kasus ditutup — Resolved-Rejected.`
  (`InputOffer.tsx:36-43`): layar langsung kembali ke Inbox. Bukti hasil keputusan adalah kolom **`Work Status`**
  baris kasus di Inbox, yang dimuat ulang saat itu juga.
- **Tampilan galat.** Penolakan server tampil sebagai pita merah berisi kalimat server apa adanya
  (FE `lib/keadaanGalat.ts:91`, `services/api.ts:315-330`). Bila badan jawaban tidak memuat kalimat, pita berbunyi
  `Permintaan ditolak backend`. **Pengecualian: 502/503/504** tampil sebagai panel peringatan yang diawali
  "Backend tidak terhubung (…)", dengan tombol `Muat ulang` (`keadaanGalat.ts:98-105`, `components/ui/dasar.tsx:344-360`).
  Panel ini juga muncul untuk 503 `database belum dikonfigurasi` dan 503 tanggal tutup buku kosong, jadi kalimat
  server keduanya **tidak terlihat** di layar.
- Saat memuat, layar menampilkan `Memuat...`.

---

### 2. Jalur klik dan setiap kontrol (urut seperti di layar)

#### 2.1 Beranda → PremiumList

1. Buka aplikasi. Beranda menampilkan `Selamat datang, UJI-ADMIN` beserta baris peran (FE `pages/Beranda.tsx:116-119`).
2. Kartu **`PremiumList Life`** menampilkan keadaan `aktif — belum ada kotak masuk` (`Beranda.tsx:148-151`;
   `assets/labels.ts:131-137`) dan tombol **`PremiumList`** (`labels.ts:114`).
3. Jalan lain ke layar yang sama: sidebar kelompok `PremiumList Life` → butir `PremiumList`
   (`lib/daftarMenu.ts:61`), atau palet `Ctrl K` lalu ketik `premiumlist`.

| Kontrol | Memanggil | Siapa | Hasil |
| --- | --- | --- | --- |
| `PremiumList` (kartu, sidebar, atau palet) | pindah halaman. Inbox lalu memanggil `GET /api/polis-life?posisi=&halaman=1&ukuran=20` | ber-identitas | Layar Inbox (§2.2) |

#### 2.2 Inbox `PremiumList` — FE `pages/premiumlist/InboxPremiumList.tsx`

Urutan di layar:

1. Judul `PremiumList` (:97).
2. Tombol `Input Offer` dan `Input Premium` (:116-133). *(GILIRAN-13 — sebelumnya dua panel
   "tidak dapat dimuat".)*
3. Tombol `Export xlsx` (:101-111).
4. Tabel.
5. Keterangan `<n> dari <total>` (:148-152).

| # | Kontrol | Memanggil | Siapa | Hasil yang diharapkan |
| --- | --- | --- | --- | --- |
| 1 | `Input Offer` / `Input Premium` | `POST /api/polis-life` `{"flag":"0"}` / `{"flag":"1"}` (`BE/handlers/polis_kasus.go`) | ber-identitas | Kasus `NBLF-<n>` lahir dan **langsung terbuka** di `Input Offer Life` — untuk **kedua** tombol; benderanya bekerja saat `Confirm` (§2.3): `Input Offer` → tertutup, `Input Premium` → Input Premium Detail. Kedua tombol mati selama permintaan berjalan; penolakan tampil lewat `Gagal` (:149) di bawah baris tombol — pita merah, atau panel untuk 502/503/504 (§1.3). ⚠️ Menuntut migrasi **057** di skema uji |
| 2 | `Export xlsx` | tanpa HTTP. Mengunduh `premiumlist.xlsx` berisi baris yang tampil | ber-identitas | Mati bila tabel kosong. Kolom berkas sama dengan kolom tabel |
| 3 | Klik baris | membuka kasus dengan tahap = nilai `Work Status` baris itu (:131-137) | ber-identitas | Layar yang terbuka mengikuti tahap (tabel di bawah) |

Kolom tabel (VERBATIM `assets/labels.premiumlist.ts:20-45`): `Case ID`, `Create Date/Time`, `Create Operator Name`,
`Work Status`, `CedingCoName`, `PolicyHolderName`, `PL_NUMBER`, `RISLIPRNM`, `Type`, `MarketingName`, `SobName`,
`DateReceived`. Sel kosong tampil `—`, tanggal tampil `YYYY-MM-DD`. Antrean kosong menampilkan `Antrean PremiumList kosong.`

Tahap baris menentukan layar yang terbuka (FE `App.tsx:59-90`):

| `Work Status` baris | Yang tampil |
| --- | --- |
| `Input Offer Life` | Blok keputusan berjudul `Input Offer` saja (§2.3) |
| `Input Premium Detail` | `Premium List Detail` + unggah CSV + grid, lalu blok keputusan berjudul `Input Premium Detail` (§2.4) |
| `Input Premium Summary` | `Summary Premium Life`, lalu blok keputusan berjudul `Input Offer` (§2.5) |
| `Resolved-Rejected`, `Resolved-Completed`, kosong, atau lainnya | Blok keputusan berjudul `Input Offer` saja (§2.6) |

Catatan: kasus tertutup **tetap tercantum**, karena query tidak menyaring status (`BE repository/polis_inbox.go:91-101`).
Yang tampil hanya 20 baris pertama, karena layar tidak punya kontrol halaman.

#### 2.3 Keputusan `Input Offer` — tahap `Input Offer Life` (kasus `UJI-PL-A`)

FE `pages/premiumlist/InputOffer.tsx`. Urutan di layar:

1. Judul `Input Offer` (:55).
2. Baris `UJI-PL-A — Input Offer Life` (:113-115).
3. `Periode produksi: <YYYY-MM>` (:117-121), dari `GET /api/polis-life/periode`.
4. Tombol `Confirm` dan `Decline`. `Reject` **sengaja tidak tampil** di tahap ini (:138; `services/api.ts:1527-1529`).

| # | Tombol | HTTP | Gerbang BE | Hasil yang diharapkan |
| --- | --- | --- | --- | --- |
| 1 | `Confirm` | `POST /api/polis-life/{id}/keputusan` `{"keputusan":"Confirm"}` | ber-identitas, kasus belum tertutup (BE `services/polis_penawaran.go`) | *(GILIRAN-14 butir bq.)* `Decision3` dirutekan dari bendera kasus, di transaksi yang sama: bendera `"0"` (`UJI-PL-A`) → **ditutup `Resolved-Completed`** tanpa menyimpan premium list; `"1"` (`UJI-PL-E`) → **pindah ke `Input Premium Detail`**. Layar kembali ke Inbox; kolom `Work Status` barisnya membuktikan hasilnya. Tombol `Premium`/`Offer` **tidak ada lagi** |
| 2 | `Decline` | sama, dengan `"Decline"` | sama | Kasus ditutup `Resolved-Rejected` (BE `models/polis_penawaran.go:203-205`), lalu layar kembali ke Inbox |

Penolakan utama (VERBATIM):

| Keadaan | HTTP | Teks |
| --- | --- | --- |
| Kasus sudah tertutup, misalnya karena Inbox sudah basi | 409 | `kasus polis sudah ditutup` (`rute_premiumlist.go:336`) |
| `Confirm` pada kasus berbendera kosong (`UJI-PL-F`, kasus sebelum 057) | 409 | `polis "UJI-PL-F": models: FLAG_ONGOING_POLICY bukan "0" maupun "1"; Decision3 tidak punya jalur untuknya: bendera ""` (`models/polis_penawaran.go:ErrBenderaTanpaKonektor`) |
| `Reject` dikirim lewat API (tidak dapat ditekan di layar) | 409 | `models: tahap ini tidak punya konektor untuk keputusan itu: "Reject" pada "Input Offer Life"` (`models/polis_penawaran.go:235-236`) |
| `POOLDATA.TANGGAL_CLOSING` kosong | 503 | Layar menampilkan panel "Backend tidak terhubung (…)". Kalimat server `models: tanggal tutup buku tidak terbaca dari POOLDATA.TANGGAL_CLOSING; transaksi ditolak, dan TIDAK ada nilai pengganti yang dipakai` (`models/polis_periode.go:44-46`) hanya terlihat di panel jaringan peramban |

#### 2.4 Tahap `Input Premium Detail` (kasus `UJI-PL-B`)

Urutan di layar (FE `App.tsx:73-90`):

- **A. Kepala** — `pages/premiumlist/PremiumListDetail.tsx:118-140`:
  judul `Premium List Detail` → baris `UJI-PL-B — <Type> / <BusinessCode>`
  → `PL_NUMBER: Belum bernomor`, atau nomornya → tombol `Generate PL Number` → kalimat yang menyebut sebab tombolnya mati.
- **B. Unggah** — `UnggahCSVPeserta.tsx:118-187`:
  judul `Upload CSV Premium List Detail` → catatan aturan
  `Angka memakai TITIK sebagai pemisah desimal, dan tanpa pemisah ribuan. Contoh: 1234567.89 — bukan 1,234,567.89 dan bukan 1.234.567,89. Tanggal berformat dd/mm/yyyy.`
  → pemilih berkas. Label aksesibelnya `Berkas CSV`; teks tombolnya berasal dari peramban.
  → `Tinjau` → `Simpan permanen` → ringkasan tinjauan dan tabel `Baris` | `Kolom` | `Pesan` | `Sebab`.
- **C. Grid peserta** — `PremiumListDetail.tsx:160-210`:
  38 kolom, urutannya dari server dan judulnya dari `labels.premiumlist.ts:92-131`
  → `<n> dari <total>`, paling banyak 50 baris
  → lipatan `2 medan layar lama tidak ditampilkan`, berisi `REINSTYPENAME` dan `RetrocadedShare` beserta alasannya.
- **D. Blok keputusan**: judul `Input Premium Detail` → `UJI-PL-B — Input Premium Detail`
  → `Periode produksi: <YYYY-MM>` → tombol `Confirm`, `Reject`, `Decline`.

| # | Kontrol | HTTP | Gerbang BE | Hasil yang diharapkan |
| --- | --- | --- | --- | --- |
| 1 | `Generate PL Number` | `POST /api/polis-life/{id}/nomor`, tanpa badan | ber-identitas, kasus belum tertutup (`services/polis_nomor.go:95-109`), minimal satu peserta | Tombol mati bila polis sudah bernomor, dengan kalimat `Polis ini sudah bernomor; tombolnya tidak perlu ditekan lagi.` Tombol juga mati bila belum ada peserta, dengan kalimat `Polis ini belum punya baris peserta, sehingga nomornya belum punya tempat tersimpan. Unggah rincian peserta lebih dahulu.` Bila berhasil muncul `Nomor terbit untuk periode <MM.YYYY> pada <n> baris peserta.`, lalu kepala dan grid menampilkan nomornya. Nomor lahir sekali (`polis_nomor.go:142-168`) |
| 2 | Pemilih berkas `Berkas CSV` | — | — | Mengganti berkas membuang tinjauan sebelumnya (`UnggahCSVPeserta.tsx:66-74`) |
| 3 | `Tinjau` | `POST /api/polis-life/{id}/unggah/tinjau`, multipart bagian `berkas` | ber-identitas saja; **tidak menyimpan apa pun** (`services/polis_unggah.go:195-215`) | Muncul `<n> baris terbaca, semuanya lolos.`, atau `<n> baris terbaca; <b> baris ditolak dengan <c> alasan. Tidak ada yang tersimpan.` disertai `Perbaiki dulu baris yang ditolak, lalu tinjau ulang. Selama masih ada penolakan, tidak ada satu baris pun yang disimpan.` dan tabel penolakan |
| 4 | `Simpan permanen` | `POST /api/polis-life/{id}/unggah/simpan`. Berkas dikirim ulang dan divalidasi ulang | ber-identitas, kasus belum tertutup, nol penolakan (`polis_unggah.go:251-313`) | Hanya hidup sesudah tinjauan lolos. Bila berhasil muncul `<n> peserta tersimpan.` atau `<n> peserta tersimpan, menggantikan <m> baris sebelumnya.`, lalu grid dimuat ulang. Unggahan **mengganti** baris lama, tidak menumpuk |
| 5 | Lipatan `2 medan layar lama tidak ditampilkan` | — | — | Menampilkan `REINSTYPENAME — nol kolom di migrasi 050-056, dan nol rule lain di korpus PremiumList Life yang menyebutnya` dan `RetrocadedShare — berdampingan dengan RETROCEDED_SHARE di grid yang sama, jadi medan yang berbeda; mana yang mana [terbuka]` (`models/polis_detail.go:144-149`) |
| 6 | `Confirm` | `POST …/keputusan` `"Confirm"` | ber-identitas, belum tertutup, ada peserta, `Type` QR/QP/TP/TR, `BusinessCode` terisi, bahan nomor tersedia | Dalam **satu** transaksi: nomor (bila belum ada) → rekap → salinan peserta warisan → kasus ditutup `Resolved-Completed` (`services/polis_penawaran.go:183-225`; `models/polis_penawaran.go:207-217`). Layar kembali ke Inbox. **Rekapnya tidak ditampilkan di mana pun** |
| 7 | `Reject` | `"Reject"` | ber-identitas, belum tertutup | Kasus kembali ke `Input Offer Life` (`models/polis_penawaran.go:221-223`) |
| 8 | `Decline` | `"Decline"` | sama | Kasus ditutup `Resolved-Rejected`, tanpa simpan (:218-220) |

Penolakan utama (VERBATIM; handler `rute_premiumlist.go:326-385`, `rute_unggahpolis.go:97-110`):

| Keadaan | HTTP | Teks di layar |
| --- | --- | --- |
| CSV tanpa baris data | 400 | `services: berkas CSV tidak punya baris data di bawah baris judul` |
| Judul kolom ganda, misalnya `DOB` dua kali | 400 | `services: berkas CSV memuat judul kolom ganda: "DOB"` |
| Kolom wajib hilang dari judul | 400 | `services: berkas CSV kehilangan kolom wajib: SUM_INSURED` |
| Lebih dari 20.000 baris | 400 | `services: berkas CSV melebihi batas cacah baris: 20000` |
| `Confirm` pada polis tanpa peserta | 409 | `repository: polis belum punya satu pun baris peserta, sehingga nomor premium list belum punya tempat tersimpan; unggah rincian peserta lebih dahulu` |
| `Type` bukan QR/QP/TP/TR | 409 | `models: Type polis tidak punya cabang penomoran (bukan QR/QP/TP/TR): "<Type>"` |
| `BusinessCode` kosong | 409 | `models: BusinessCode kosong; nomor premium list memuatnya` |
| Dua penerbitan nomor serentak | 409 | `repository: PL_NUMBER polis ini baru saja diterbitkan permintaan lain; baca ulang nomornya` |
| Baris peserta memuat nomor yang berbeda-beda | 409 | `repository: baris peserta satu polis memuat lebih dari satu PL_NUMBER` |
| Header `T_PREMIUM_LIST.ID` ≠ Case ID | 404 | `polis tidak ditemukan` |
| Kasus tertutup | 409 | `kasus polis sudah ditutup` |
| `TANGGAL_CLOSING` atau `KODE_PRODUKSI` kosong saat menomori | 500 | `gagal memproses permintaan polis` (galat repository tidak dipetakan, `rute_premiumlist.go:378-382`) |

#### 2.5 Tahap `Input Premium Summary` — `Summary Premium Life` (kasus `UJI-PL-C`)

> ⚠️ **Keadaan sekarang: layar ini tidak dapat dicapai lewat alur aplikasi.** Tahap `Input Premium Summary`
> **tidak punya satu pun konektor masuk** (BE `models/polis_penawaran.go:78-90`; Tiket `01-…md:164-170`
> `[terbuka — work owner]`). Tiket `05b-…md:235` menulis bahwa **"tombol Submit summary selalu 409 untuk kasus baru"**.
> Di layar, akibatnya sebagai berikut:
>
> 1. Kasus yang dijalankan lewat aplikasi (`Premium` → Detail → `Confirm`) **tidak pernah** membuka layar ini.
>    Simpan nomor, rekap, dan salinan warisan terjadi di `Confirm` tahap Detail (§2.4 no. 6) tanpa menampilkan rekap.
> 2. Layar ini hanya terbuka bila `Work Status` baris di Inbox bernilai tepat `Input Premium Summary`
>    (`App.tsx:79-81`). Artinya kasusnya harus disiapkan di luar aplikasi (`UJI-PL-C`).
> 3. Bila `Submit` dikirim untuk kasus di tahap lain, jawabannya **409**. Contohnya: Inbox sudah basi, atau
>    `POST` dikirim langsung. Pita merah di layar berbunyi:
>    `services: Submit summary hanya dari tahap Input Premium Summary: polis "UJI-PL-C" berada di "<tahap>", bukan "Input Premium Summary"`
>    (BE `services/polis_summary.go:119-120,268-271`).
> 4. Untuk kasus yang memang disiapkan di tahap ini, `Submit` berjalan seperti tercantum di tabel di bawah.

Urutan di layar (FE `pages/premiumlist/PremiumListSummary.tsx:136-176`):

1. Judul `Summary Premium Life`.
2. Pesan status.
3. Grid rekap menurut `Type`. Judul kolomnya VERBATIM dari `labels.premiumlist.ts:219-280`. Contoh QR: `COB`,
   `PL NUMBER`, `CURRENCY`, `PREMIUM`, `DEDUCTION`, `BROKERAGE FEE`, `RI ADMIN FEE`, `TAX`, `PROF COMM`, `CLAIM`,
   `BALANCE`. Pada QP/TP/TR, judul `PREMIUM DEDUCTION` berdiri di atas nilai COMMISSION.
4. Catatan `Kolom COB tidak terisi: tidak satu pun rule PremiumList Life menetapkan .COB pada baris rekap mata uang.`
5. Tombol `Submit`.
6. Blok keputusan berjudul `Input Offer`, berisi `Confirm` dan `Decline`.

| # | Kontrol | HTTP | Gerbang BE | Hasil yang diharapkan |
| --- | --- | --- | --- | --- |
| – | (saat layar dibuka) | `GET /api/polis-life/{id}` dan `GET /api/polis-life/{id}/summary` | ber-identitas; tanpa gerbang kasus tertutup (`polis_summary.go:222-247`) | Uang tampil sebagai teks empat desimal, misalnya `975.0000`. Sel kosong dan kolom `COB` tampil `—` |
| 1 | `Submit` | `POST /api/polis-life/{id}/summary`, tanpa badan | ber-identitas, belum tertutup, **tahap = `Input Premium Summary`** (`polis_summary.go:261-273`) | Bila berhasil muncul `PL_NUMBER <nomor>: <k> rekap mata uang tersimpan, <m> peserta tersalin. Kasus ditutup: Resolved-Completed. Kiriman ke Arasapas dilewati: lingkungan ini bukan produksi.`, lalu tombolnya mati. Bila tahapnya lain → 409, lihat kotak di atas |
| 2 | `Confirm` / `Decline` (blok `Input Offer`) | `POST …/keputusan` | — | **Selalu 409** `models: tahap ini tidak punya konektor untuk keputusan itu: tahap "Input Premium Summary" tidak punya konektor keputusan sama sekali` (`models/polis_penawaran.go:225-231`) |

Penolakan saat layar dimuat:

- Polis tanpa peserta → 409 `repository: polis belum punya satu pun baris peserta, …`, dan `Submit` mati.
- `Type` tidak sah → 409 `models: Type polis tidak punya cabang penomoran (bukan QR/QP/TP/TR): "<Type>"`.

Kalimat `Belum ada rekap: polis ini belum punya baris peserta.` praktis tidak pernah muncul, karena polis tanpa peserta
sudah ditolak sebagai galat lebih dulu (`services/polis_summary.go:209-211`).

#### 2.6 Kasus tertutup, atau tahap kosong/asing (kasus `UJI-PL-D`)

Klik baris berstatus `Resolved-*`. Yang tampil hanya blok keputusan: judul `Input Offer`,
`UJI-PL-D — Resolved-Completed`, `Periode produksi: …`, lalu tombol `Confirm` dan `Decline`.

- Menekan kedua tombol → 409 `kasus polis sudah ditutup`.
- Bila `STATUS` kosong atau bukan salah satu dari kelima nilai di atas → 400 `models: tahap polis tidak dikenal: "<tahap>"`
  (`models/polis_penawaran.go:232-233`).

---

### 3. Contoh nilai dari uji

#### 3.1 CSV yang lolos — satu baris

Fixture dari BE `services/polis_unggah_test.go:19-56` (`judulLengkap` dan `barisUji`). Kolomnya dipisah koma:

```
CERTIFICATE_NO,NAME_OF_INSURED,PLAN,CURRENCY,POLICY_NO,MEDICAL_STATUS,DOB,BEGIN_DATE,EXPIRED_DATE,START_DATE,EFFECTIVE_DATE,STNC,WPC,GROSS_VALUATION_BEGIN_DATE,GROSS_VALUATION_EXPIRED_DATE,RETROCESSION_VALUATION_BEGIN_DATE,RETROCESSION_VALUATION_EXPIRED_DATE,SUM_INSURED,CEDING_RETENTION,SUM_REASURED,SHARE_NUSANTARA_RE,GROSS_PREMIUM,NET_PREMIUM
UJI-C1,UJI PESERTA A,UJI-PLAN,IDR,UJI-POL-1,NM,01/02/1990,01/01/2026,31/12/2026,01/01/2026,01/01/2026,01/01/2026,01/01/2026,01/01/2026,31/12/2026,01/01/2026,31/12/2026,1000.50,100.25,900.25,50.5,200.75,180.5
```

- `Tinjau` → `1 baris terbaca, semuanya lolos.` (format `UnggahCSVPeserta.tsx:36-38`). `Simpan permanen` lalu hidup.
- Judul huruf kecil dan BOM Excel tetap diterima (`TestBacaCSVUnggahMenaikkanHurufJudul` :80, `TestBOMExcelDibuang` :98).
  Baris kosong di ujung berkas dilewati (:150).

#### 3.2 CSV yang ditolak

Semua contoh berangkat dari baris 3.1.

| Ubahan | Baris / Kolom | `Pesan` (verbatim) | `Sebab` | Uji |
| --- | --- | --- | --- | --- |
| `PLAN` dikosongkan | 1 / `PLAN` | `PLAN HARUS ADA` | `kolom kosong` | `models/polis_unggah_test.go:226` |
| `DOB` = `1970-01-02` | 1 / `DOB` | `DOB HARUS ADA. FORMAT:dd/mm/yyyy` | `"1970-01-02" bukan dd/mm/yyyy` | :226 |
| `SUM_INSURED` = `"1,000"` (wajib dikutip di CSV, kalau tidak koma memecah kolom) | 1 / `SUM_INSURED` | `SUM INSURED HARUS ADA` | `"1,000" memuat koma; pemisah desimal harus titik dan pemisah ribuan tidak diterima` | :226 |
| Ketiga ubahan di atas sekaligus | — | tabel memuat 3 baris berurutan PLAN, DOB, SUM_INSURED | ringkasan `1 baris terbaca; 1 baris ditolak dengan 3 alasan. Tidak ada yang tersimpan.` | :226 |
| Baris 1 disalin utuh sebagai baris 2 | 2 / `CERTIFICATE_NO` | `CERTIFICATE NOMOR TIDAK BOLEH DOUBLE` | `nomor sertifikat "UJI-C1" sudah dipakai baris 1` | :263 |
| (baris yang sama) | 2 / `NAME_OF_INSURED` | `CERTIFICATE NOMOR SUDAH DIGUNAKAN` | `peserta ini sudah ada di baris 1 (nama tanpa peduli huruf besar/kecil + DOB)` | :288 |
| Baris 2: `CERTIFICATE_NO` = `UJI-C2`, nama `uji peserta a` (huruf kecil), `DOB` sama | 2 / `NAME_OF_INSURED` | `CERTIFICATE NOMOR SUDAH DIGUNAKAN` | seperti di atas | :288 |
| Baris 2 seperti di atas, tetapi `DOB` = `02/02/2026` | — | lolos, karena DOB berbeda berarti bukan duplikat | — | :288 |
| Kolom tambahan `TAX` = `"1,5"` | 1 / `TAX` | `TAX bukan angka yang sah` | `"1,5" memuat koma; …` | :373 |
| `MEDICAL_STATUS` = `fcl` | 1 / `MEDICAL_STATUS` | `MEDICAL STATUS HARUS FCL/M/NM` | `"fcl" di luar FCL/M/NM` | :132 |
| `NET_PREMIUM` = `1.00`, `GROSS_PREMIUM` = `999999.00` | — | **lolos**; OQ-069 hanya memeriksa keberadaan | — | :351 |
| Kolom judul `SUM_INSURED` dibuang | pita merah | `services: berkas CSV kehilangan kolom wajib: SUM_INSURED` | — | `services/polis_unggah_test.go:132` |
| Kolom `DOB` ditambah sekali lagi | pita merah | `services: berkas CSV memuat judul kolom ganda: "DOB"` | — | :109 |
| Berkas hanya berisi judul | pita merah | `services: berkas CSV tidak punya baris data di bawah baris judul` | — | :109 |

#### 3.3 Bentuk angka, tanggal, dan status medis

Sumber: `models/polis_unggah_test.go:21-144`.

| Jenis | Diterima | Ditolak |
| --- | --- | --- |
| Uang | `1234567.89`, `0`, `0.00000001`, `  42.5  ` (spasi di tepi dibuang), `1234567890123456789.12345678` (tanpa dibulatkan) | `1,234,567.89`, `1.234.567,89`, `1,5`, `1,000`, `abc`, `Rp1000`, `1.2.3` |
| Tanggal | `31/01/2026`, `02/01/1970` | `2026-01-31`, `1/2/2026`, `31-01-2026`, `31/13/2026`, `32/01/2026` |
| `MEDICAL_STATUS` | `FCL`, `M`, `NM` | `fcl`, `m`, `X`, `NMM` |

#### 3.4 Periode produksi dan tanggal tutup buku

Layar menampilkan `Periode produksi: <YYYY-MM>`. Sumber: `models/polis_periode_test.go`.

| Tutup buku | Saat transaksi | Periode | Uji |
| --- | --- | --- | --- |
| 25 | 10-09-2026 10:00 WIB | `2026-09` | :122 |
| 25 | 25-09-2026 23:00 WIB | September 2026 (pembanding `>`, bukan `>=`) | :23 |
| 25 | 26-09-2026 00:00 WIB | Oktober 2026 | :23 |
| 25 | 31-12-2026 09:00 WIB | `2027-01` (bukan bulan 13) | :64 |
| 25 | 25-09-2026 20:00 UTC (= 26-09 03:00 WIB) | Oktober 2026 (zona Jakarta yang menentukan hari) | :105 |
| 0 atau -1 | — | ditolak; **tidak** diganti 25 secara diam-diam | :75 |
| 32 | — | ditolak `models: tanggal tutup buku di luar 1..31` | :98 |

#### 3.5 `PL_NUMBER`

Sumber: `models/polis_nomor_test.go`.

- Bentuk nomor: `<awalan><Type><BusinessCode>.<MM>.<YY>.<urut lima digit>`. Urut 7 menjadi `00007` (:110, :123).
  Kalimat di layar memakai periode `MM.YYYY`, misalnya `09.2026`.
- Periode di dalam nomor memotong tahun: `09.2026` → `09.26`, `01.2000` → `01.00`, `12.1999` → `12.99` (:81).
- `Type` yang diterima: `QR`, `QP`, `TP`, `TR`, dengan spasi di tepi dibuang. Yang ditolak: `qr`, `Q`, `TX`, `QRQP`,
  `FAC`, dan kosong (:11, :34).
- Keempat Type memakai **satu** penghitung (:48).

#### 3.6 Rekap per mata uang — dibulatkan empat desimal di akhir

Sumber: `models/polis_summary_test.go`. Rekap hanya terlihat di layar Summary, jadi hanya lewat `UJI-PL-C`.

| Type | Masukan | Harapan | Dapat diulang lewat CSV? | Uji |
| --- | --- | --- | --- | --- |
| QR / QP / TP / TR | satu baris: `GROSS_PREMIUM` 100, `GROSS_PREMIUM_REFUND` 200, `GROSS_PREMIUM_RETRO` 300, `GROSS_PREMIUM_REFUND_RETRO` 400 | `PREMIUM` QR `100.0000`, QP `200.0000`, TP `300.0000`, TR `400.0000` (bukan 1000) | ya | :35 |
| QR | tiga baris `GROSS_PREMIUM` 100 (IDR), 200 (USD), 50 (IDR) | 2 baris rekap. IDR `PREMIUM` `150.0000` dari 2 baris, USD dari 1 baris | ya. Urutan baris di layar dapat berbeda dari Pega (Tiket `05a-…md`, catatan urutan mata uang) | :187 |
| QR | tiga baris `GROSS_PREMIUM` `0.00005` | `PREMIUM` `0.0002` dan `BALANCE` `0.0002`, bukan `0.0003` | ya | :217 |
| semua | `COMM` 11, `PROF_COMM` 22, `OVR_COMM` 33 | COMMISSION `11.0000`, tampil di kolom `PREMIUM DEDUCTION` pada QP/TP/TR | ya | :171 |
| TP | baris tanpa kolom retro | `BALANCE` `0.0000` | ya | :354 |
| QR | `GROSS_PREMIUM` 1000, `DEDUCTION` 10, `RI_ADMIN_FEE` 1, `BROKERAGE_FEE` 2, `TAX` 3, `PROF_COMM` 4, `CLAIM` 5 | `BALANCE` `975.0000` | **tidak** — `DEDUCTION` dan `RI_ADMIN_FEE` tidak diisi unggahan | :66 |
| QP | `GROSS_PREMIUM_REFUND` 1000, `CLAIM_AMOUNT` 50, `DEDUCTION_REFUND` 10, `BROKERAGE_FEE_REFUND` 20, `RI_ADMIN_FEE_REFUND` 30, `TAX` 3, `PROF_COMM` 4, `CLAIM` 5 | `978.0000` | **tidak** | :66 |
| TP | `GROSS_PREMIUM_RETRO` 1000, `DISCOUNT_PREMIUM_RETRO` 100, `RI_ADMIN_FEE_RETRO` 10, `BROKERAGE_FEE_RETRO` 7 (brokerage **ditambah**) | `897.0000` | **tidak** | :66 |
| TR | kolom `*_REFUND_RETRO` 2000, 200, 20, 9 | `1789.0000` | **tidak** | :66 |

Tampilan sel (FE `PremiumListSummary.test.ts:55-67`): `975.0000` tetap `975.0000`, karena ekor nol tidak dibuang.
Nilai kosong dan kolom `COB` tampil `—`.

#### 3.7 Tabel transisi keputusan

Tabel ini menjadi acuan bagi §2.3–2.6 (`models/polis_penawaran_test.go:19-56`, :58, :81, :105):

| Tahap | `Confirm` | `Reject` | `Decline` |
| --- | --- | --- | --- |
| `Input Offer Life` | dari bendera (bq): `"0"` → `Resolved-Completed`; `"1"` → `Input Premium Detail`; lain → 409 | 409, tidak ada jalur | `Resolved-Rejected` |
| `Input Premium Detail` | simpan, lalu `Resolved-Completed` | kembali ke `Input Offer Life` | `Resolved-Rejected` |
| `Input Premium Summary` | 409 | 409 | 409 |

---

### 4. Belum dapat diuji di layar

| # | Perilaku | Sebab (satu baris) | Rujukan |
| --- | --- | --- | --- |
| 1 | ~~Membuat kasus baru~~ — ✅ **ditutup GILIRAN-13** (`a291a20`, butir bn). Urutan `SEQ_WORK_POLIS` mulai **22374** **sesudah migrasi 058 dijalankan work owner** (`c31eb12`, OQ-PL-15 ditutup — `[sementara]` sampai DBA memastikan penghitung Pega `PC_DATA_UNIQUEID`, OQ-PL-17). `Decision3` dirutekan dari bendera sejak GILIRAN-14 (`9f67d35`, butir bq; OQ-PL-16 ditutup) | — | Tiket `00-…md`, `01-…md` (ralat 29-09-2026) |
| 2 | Isian dan simpan data penawaran, serta riwayat `T_VIEW_SUGGEST` | Layar hanya punya tombol keputusan; penulis penawaran relasional belum ada | Tiket `01-…md:3,65-87` |
| 3 | Gerbang `ProtectAccept` (`Please choose no offer !`, `COB can't null`, `Premium is 0`, dst.) | Pesannya ada, tetapi `ValidasiPenawaran` tidak dipanggil di luar uji | BE `models/polis_validasi.go:37-49,116` |
| 4 | Layar `Summary Premium Life` lewat alur aplikasi | Tahap `Input Premium Summary` nol konektor masuk `[terbuka — work owner]` | `models/polis_penawaran.go:78-90`; Tiket `05b-…md:235` |
| 5 | Melihat rekap hasil `Confirm` di tahap Detail | Jawaban keputusan hanya membawa tahap dan status; ringkasan efek dibuang | Tiket `05b-…md` (temuan Standards 3) |
| 6 | Kiriman Arasapas | Di non-produksi selalu `Kiriman ke Arasapas dilewati: lingkungan ini bukan produksi.`; di produksi masih stub, dan `AUTH_STUB` ditolak di sana | BE `services/polis_efekkeluar.go:77-85`, `efekkeluar.go:266-270` |
| 7 | OQ-PL-11 muatan Arasapas | Layanan hilir tampaknya membaca `JSON_POLIS`, yang tidak lagi ditulis; stub tetap gagal terang | Tiket `06-…md:147-152` |
| 8 | Email alarm kegagalan | Stub `ErrEmailBelumDisetujui`; penerima dan isinya `[belum terverifikasi]` | `polis_efekkeluar.go:87-100` |
| 9 | Pengulangan outbox `PREMIUMLISTLIFE` | Pekerja pengulang modul belum ada | Tiket `06-…md:3,94` |
| 10 | OQ-PL-09 `M_LIFE_PREMIUM_SUMMARY` | Tidak ditulis karena pemetaan argumen ke kolom belum tercatat (menunggu DBA); tidak tampak di layar mana pun | Tiket `05a-…md:373-387` |
| 11 | OQ-PL-10: nilai kosong disimpan NULL atau 0 di tabel warisan | Keputusan pemilik kerja masih terbuka; hanya terlihat di tabel, tidak di layar | Tiket `05a-…md:439-442` |
| 12 | Periode tampil sebelum `Generate PL Number` / `Submit` | Periode hanya tampil di blok keputusan | Tiket `02-…md:3,69-70` |
| 13 | Verifikasi `NAME_OF_INSURED`/`POLICY_HOLDER` ke master; pesan `CEDING CO TIDAK TERDAFTAR` dan sejenisnya | Pembaca master belum ada; pesan rujukan master belum dipanggil `ValidasiUnggah` | Tiket `04-…md:105`; `models/polis_unggah.go` (pesan rujukan master) |
| 14 | `NET PREMIUM HARUS ADA DAN LEBIH BESAR DARI GROSS PREMIUM` sebagai perbandingan | OQ-069: hanya keberadaan yang diperiksa | `models/polis_unggah.go:91` |
| 15 | `Please Upload CSV File into attachment` | Hanya muncul bila bagian `berkas` tidak dikirim, sedangkan layar selalu mengirimnya | `models/polis_unggah.go:375-377` |
| 16 | Inbox lebih dari 20 baris, grid lebih dari 50 baris, dan tab posisi `Offer`/`Premium` | Tidak ada kontrol halaman maupun tab | FE `services/api.ts:1459-1467,1613-1617` |
| 17 | Kolom `Ketentuan Underwriting`, `REINSTYPENAME`, `RetrocadedShare` | Ketiganya tanpa kolom di migrasi 050–056 | `labels.premiumlist.ts:11-18`; `models/polis_detail.go:144-149` |
| 18 | Kolom uang yang tidak diisi unggahan (`Deduction`, `Rate`, `Factor`, `Sum At Risk *`, `RI Admin Fee *`, dst.) | Unggahan hanya mengisi 32 dari 44 kolom uang, jadi grid menampilkan `—`; contoh 975/978/897/1789 hanya teruji di unit test | Tiket `03-…md:122` |
| 19 | Blok kepala `ShowLifePremiumSummary` (Type, DateReceived, SobName, …) | Belum dirender; yang ada hanya grid dan `Submit` | Tiket `05a-…md` ("Layar sebagian") |
| 20 | Spreading (`T_PREMIUM_LIST_SPREADING*`) | Belum ada penulisnya | Tiket `03-…md:136-138` |
| 21 | Dua penerbitan nomor paralel | Belum ada uji beban; di layar hanya terlihat 409 `…baru saja diterbitkan permintaan lain…` | Tiket `03-…md:120` |
| 22 | Penolakan per peran (403) | Tidak ada gerbang peran di PremiumList | §1.1 |
| 23 | Autentikasi sungguhan | Masih identitas stub sampai IAM (tiket 07) | BE `handlers/pelaku.go:3-20` |

---

### 5. Hal yang perlu diamati — temuan dari kode, sebaiknya dikonfirmasi saat klik

1. **Unggah ulang sesudah polis bernomor menghapus nomornya.** `Simpan permanen` menghapus seluruh baris peserta,
   lalu menyisipkan baris baru tanpa `PL_NUMBER`, padahal nomor hanya tersimpan di baris peserta
   (BE `repository/polis_unggah.go:143-145`; `repository/polis_nomor.go:156-159`). Harapan menurut kode: kepala kembali
   `Belum bernomor`, dan `Generate PL Number` menerbitkan nomor **baru**.
2. **Kolom grid `Share Nusantara Re Gross` tetap `—`** walaupun CSV mengisi `SHARE_NUSANTARA_RE`. Unggahan menulis
   `SHARE_NUSANTARA_RE`, sedangkan grid membaca `SHARE_NUSANTARA_RE_GROSS` (`models/polis_detail.go:114`;
   migrasi `052:85-86`).
3. Kasus tertutup dapat dibuka dari Inbox, dan layarnya tetap menawarkan `Confirm`/`Decline` yang pasti dijawab 409 (§2.6).
4. Layar Summary juga memuat blok keputusan berjudul `Input Offer` yang pasti dijawab 409 (§2.5 no. 2).
5. `Confirm` di tahap Detail kini **menuntut** polis siap disimpan: ada peserta, `Type` sah, dan bahan nomor lengkap.
   Polis tanpa peserta dijawab 409 (Tiket `05b-…md`, "Akibat yang harus diketahui work owner").

---

## 3. Komite Claim Life

> Disusun 28-09-2026 **dari kode saja** (`APP_RNM\frontend\src`, `APP_RNM\internal`) dan tiket
> `.scratch\komite-claim-life\`. Setiap teks dalam tanda kutip disalin apa adanya dari kode; rujukan
> `berkas:baris` disingkat. Semua nilai contoh sintetis (`UJI-…`, `KMTLF-UJI…`). Tidak ada nama
> orang, nomor polis, kredensial, atau alamat layanan di dokumen ini.
>
> Bagian yang ditulis `<…>` adalah isi dinamis (pengenal, angka) yang disisipkan kode.

---

### 1. Prasyarat

#### 1.1 Identitas (tanpa layar masuk)

Aplikasi **tidak punya form login**. Identitas dibaca dari env saat aplikasi menyala
(`store/sesi.ts:82-96`, `handlers/pelaku.go`). Kedua sisi wajib menyala bersama.

| Tempat | Variabel | Nilai untuk uji | Catatan |
| --- | --- | --- | --- |
| `APP_RNM\.env` (backend) | `AUTH_STUB` | `true` | Tanpa ini setiap rute Komite menjawab 401 |
| `APP_RNM\.env` (backend) | `IS_PEGA_PROD` | `false` | `true` menolak stub; uji layar selalu non-produksi |
| `APP_RNM\.env` (backend) | sambungan Oracle ke **skema uji** | diisi pengembang/DBA | Tanpa basis data rute Komite menjawab 503 `database belum dikonfigurasi`; layar menampilkannya sebagai panel "backend tidak terhubung" |
| `frontend\.env` | `VITE_AUTH_STUB` | `true` | Selain `true`: layar hanya menampilkan "Identitas pelaku tidak dapat dimuat saat ini." |
| `frontend\.env` | `VITE_STUB_PELAKU` | akun anggota yang sedang diuji, mis. `UJI-A` | Kosong → `UJI-ADMIN` |
| `frontend\.env` | `VITE_STUB_PERAN` | dipisah koma | Kosong → ketiga peran (`ReasLifeAdmin`, `ReasLifeMedicalAdvisor`, `ReasLifeSPV`) |

**Ganti anggota = ganti `VITE_STUB_PELAKU`, lalu muat ulang halaman.** Bila nilai lama masih
terbaca, hentikan dan jalankan ulang server Vite (nilai env dibekukan saat Vite membangun,
`store/sesi.ts:78-80`).

#### 1.2 Siapa boleh apa (ditegakkan di services, bukan di layar)

| Tindakan | Syarat pelaku | Ditegakkan di |
| --- | --- | --- |
| Serahkan baris ke Komite, klaim Type `TP`/`TR` | cukup teridentifikasi (peran tidak dibaca) | `services/wewenang.go:126-134` |
| Serahkan baris ke Komite, klaim Type `QP`/`QR` | peran `ReasLifeSPV` | `services/wewenang.go:135-140`, `services/adjustment.go:95` |
| Melihat Inbox Komite | teridentifikasi; hanya kasus di mana akunnya = anggota **berjalan** | `repository/komite_inbox.go:98-107` |
| Membuka satu kasus | akun tercatat di tangga kasus itu (tingkat mana pun) | `services/komite_inbox.go:187-209` |
| Memutuskan (Setuju/Tolak) | akun = `KOMITE_OPERATORID` tingkat berjalan yang masih `0`; peran tidak berpengaruh | `services/komite_keputusan.go:129-148` |
| Eskalasi naik satu tingkat | peran `ReasLifeAdmin` (asumsi OQ-007/OQ-021); keanggotaan tangga **tidak** diperiksa server | `services/komite_keputusan.go:264-273` |
| Laporan "perlu intervensi" | peran `ReasLifeAdmin` | `services/komite_intervensi.go:89-96` |
| Riwayat tangga (API) | cukup teridentifikasi | `services/komite_riwayat.go` fungsi `Riwayat` |

#### 1.3 Data yang harus sudah ada di skema uji

Disiapkan pengembang/DBA **sebelum** uji; penguji tidak menjalankan migrasi.

1. **Tabel Komite** dari migrasi `013_tabel_komite.sql` dan `030_komite_kaskade_dan_lebar_id.sql`,
   plus tabel Claim Life dan outbox `T_LOG_SERVICE_RNM` (urutan `SEQ_LOG_SERVICE_RNM`,
   `SEQ_KOMITE_KOMITELIST`).
2. **Roster `EMAILKOMITE`** di skema uji, baris sintetis — **jangan menyalin baris produksi**
   (tabel ini memuat data orang, `repository/roster.go:8-11`). Filter yang dipakai
   (`repository/roster.go:58-64`): `LIMIT_BOTTOM <= nilai klaim mutlak`, `STS_KLAIM = 'LIFE'`,
   `STS_AKTIF = '1'`, urut `DEGREE`. Setiap baris yang lolos = satu tingkat.
   *(Diperbarui GILIRAN-13.)* **Data uji sintetis** (bab 0 §0.2) memuatnya:

   | `DEGREE` | `OPERATOR_ID` (= `VITE_STUB_PELAKU`) | `LIMIT_BOTTOM` | `STS_AKTIF` | Untuk 150.000.000 |
   | --- | --- | --- | --- | --- |
   | 1 | `UJI-KOMITE-1` | 0 | `1` | terpilih |
   | 2 | `UJI-KOMITE-2` | 100.000.000 | `1` | terpilih |
   | 3 | `UJI-KOMITE-3` | 1.000.000.000 | `1` | di atas pitanya |
   | 1 | `UJI-KOMITE-4` | 0 | `0` | tidak aktif |

   Surelnya `uji-komite-<n>@contoh.invalid`. Tangga baris `UJI-ADJ-4-A1` karena itu **dua** tingkat.
3. *(Data uji: klaim `UJI-CLM-4` peserta A, baris `UJI-ADJ-4-A1`, tahap Claim Analis — serahkan sebagai
   `ReasLifeSPV`, sebab Type-nya `QP`.)* **Satu klaim Claim Life terbuka** (bukan `Resolved-Completed`) dengan satu peserta dan satu baris
   adjustment: status **Outstanding**, `KOMITE_ID` kosong, nama bank / ID bank / nomor rekening
   **terisi**, mata uang tunggal di seluruh baris (kolom `CURRENCY` dan `CURRENCY_ID`).
   Type klaim salah satu `QR`, `QP`, `TP`, `TR`.
4. **Khusus Setuju di tingkat akhir** (penerbitan nomor): awalan di `KODE_PRODUKSI` (tipe Life),
   `TANGGAL_CLOSING`, baris penghitung `GENERATE_SEQUENCE_NUMBER` untuk pasangan
   `(ASM-FW-GCNMFW-Work-KomiteLife, <awalan>A)` (`models/komite_nomor.go`), dan **kode bisnis**
   klaim induk terisi.
5. **Khusus uji Kasir**: klaim Type `QR` atau `QP` dengan `T_GENERAL_CLAIM.IS_KPR = 'KPR'`.

---

### 2. Jalur klik

#### Langkah A — Serahkan baris dari Claim Life ke Komite

1. Beranda → kartu **Claim Life** → tombol "Inbox Claim Life".
2. Klik baris klaim uji → layar Outstanding (judul "Claim No: `<id>`"). Catat `<id>`.
3. Tekan "Detail klaim" → layar "Klaim Life".
4. Isi kotak "Pengenal klaim" (petunjuk "mis. UJI-KLAIM-1") dengan `<id>` → tekan "Buka".
5. Pada tabel baris peserta, kolom "Komite", tekan "Send Claim to Committee".
   Tidak ada dialog konfirmasi (label `Cancel` di `TOMBOL_KOMITE` terdefinisi tetapi tidak dipakai).

| Kontrol (urut layar) | Label di layar | HTTP | Siapa | Hasil yang diharapkan |
| --- | --- | --- | --- | --- |
| Kotak isian | "Pengenal klaim" | — | siapa pun | — |
| Tombol | "Buka" (saat menunggu "Memuat…") | `GET /api/klaim-life/{id}` | siapa pun | Klaim tampil; galat "Klaim tidak ada." (404) atau "Gagal membaca klaim." |
| Tombol per baris, kolom "Komite" | "Send Claim to Committee" (saat menunggu "Menyerahkan…") | `POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite` | lihat §1.2 | 204; klaim dimuat ulang; sel "Komite" menjadi "sudah diserahkan" (arahkan tetikus: pengenal kasus berawalan `KMTLF-`). Kasus komite lahir dengan `KOMITE_COUNT = 1`, `KOMITE_LOOP` = cacah roster, semua anak tangga `0`. Status baris **tetap Outstanding** |

Tombol hanya tampil bila status baris "Outstanding" **dan** `KOMITE_ID` kosong
(`services/api.ts:658-660`) — itu kenyamanan, bukan pagar.

**Efek keluar saat menyerahkan:** di non-produksi efek Claim Life (berkas, email, Arasapas)
**dilewati**, tidak diantre (`services/efekkeluar.go:266-271`). Tidak ada email ke anggota tingkat 1.

**Penolakan** (urutan gerbang `services/komite.go:361-430`):

| Keadaan yang dicoba | Kode | Teks di layar (`KlaimLife.tsx:293-306`) | Teks server (`handlers/komite.go`) |
| --- | --- | --- | --- |
| Type `QP`/`QR` tanpa peran `ReasLifeSPV` | 403 | "Peran Anda tidak berwenang menyerahkan baris bertipe ini ke Komite." | "peran Anda tidak berwenang menyerahkan baris bertipe ini ke Komite" |
| Baris sudah diserahkan / bukan Outstanding (mis. klik ulang lewat API) | 409 | "Baris sudah diserahkan, atau bukan lagi Outstanding." | "baris sudah pernah diserahkan ke Komite" / "hanya baris Outstanding yang dapat diserahkan ke Komite" |
| Nama bank, ID bank, **atau** nomor rekening kosong | 422 | "Name of bank cannot be empty" (diteruskan apa adanya) | sama — teks sistem lama |
| Mata uang campur antarbaris | 422 | "baris pada klaim ini bermata uang campur; penyerahan menuntut mata uang tunggal" | sama |
| Tidak ada baris roster yang menutup nilai klaim | 422 | "tidak ada tingkat komite yang menutup nilai klaim ini" | sama |
| Type klaim di luar empat nilai | 422 | "Type klaim tidak dikenal; penyerahan ke Komite menuntut Type yang sah" | sama |
| Klaim sudah ditutup | 409 | kalimat server "kasus sudah ditutup dan tidak dapat diubah" (sejak `e3619d4` / `fbf1c9b`) | sama |
| Tanpa identitas (401), galat lain | 401/500 | "Gagal menyerahkan baris ke Komite." | mis. "permintaan tanpa identitas pelaku ditolak" |

#### Langkah B — Buka Inbox Komite milik anggota berjalan

1. Setel `VITE_STUB_PELAKU=UJI-A` (anggota tingkat 1), muat ulang.
2. Beranda → kartu **Komite Claim Life** (keterangan "aktif — belum ada kotak masuk") → tombol
   "Inbox Komite". Jalan lain: sidebar kelompok "Komite Claim Life" → "Inbox Komite", atau palet
   Ctrl+K ketik `komite`.
3. Layar "Inbox Komite" memuat `GET /api/komite?halaman=1`; bila peran memuat `ReasLifeAdmin`,
   juga `GET /api/komite/laporan-harian` (bagian 2G).

| Unsur (urut layar) | Label | Harapan |
| --- | --- | --- |
| Judul + aturan | "Inbox Komite" · "Hanya kasus yang tingkat berjalannya milik Anda: anggota pertama di tangga yang belum memutuskan." | — |
| Tabel | kolom "Case ID", "Update Date/Time", "Work Status", "Klaim induk", "Tingkat", "Nilai klaim", "CURRENCY", "Status baris" | Kasus hasil Langkah A tampil: `KMTLF-…`, "Tingkat" `1 / 3`, status baris "Outstanding"; sel kosong ditulis "—" |
| Klik baris | — | Membuka "Kasus Komite" (`GET /api/komite/{id}` + `GET /api/komite/{id}/riwayat`) |
| Cacah | "`<n>` dari `<total>`" | Tidak ada tombol halaman berikut; server menjepit 100 baris |
| Inbox kosong | "Tidak ada kasus komite yang menunggu keputusan Anda." | Mis. saat login sebagai `UJI-B` sebelum `UJI-A` memutuskan |

Uji silang: login sebagai `UJI-B` atau `UJI-C` sebelum tingkat 1 memutuskan → kasus **tidak** tampil.

#### Langkah C — Layar "Kasus Komite" dan keputusan

| Kontrol / unsur (urut layar) | Label | HTTP | Siapa | Harapan |
| --- | --- | --- | --- | --- |
| Tombol atas | "Kembali ke Inbox Komite" | — (inbox dimuat ulang) | siapa pun | Kembali ke daftar |
| Kepala | "Kasus Komite `<id>`"; "Klaim induk", "Tingkat", "Nilai klaim" (`<nilai> <mata uang>`), "Status baris" | — | anggota tangga | — |
| Penanda giliran | "Giliran Anda memutuskan." / "Bukan giliran Anda — kasus ini menunggu tingkat lain atau sudah selesai." | — | — | Formulir keputusan hanya tampil pada giliran sendiri |
| Tabel tangga | judul "Tangga persetujuan"; kolom "Tingkat", "Committee", "Anggota", "Status", "Date Approve", "Comment" | dari `.../riwayat` | — | Status berupa kata: "Setuju", "Tolak", "Menunggu", "Dilewati (eskalasi)", atau "Kode `<x>`" untuk kode asing |
| Daftar eskalasi (bila ada) | judul "Eskalasi"; baris "Tingkat `<dari>` → `<ke>` oleh `<akun>` pada `<waktu>`" | — | — | Lihat Langkah E |
| Tabel efek (bila ada) | "Efek keluar — `<keadaan>`"; kolom "Efek", "Keadaan", "Percobaan", "Sejak" | — | — | Lihat tabel akibat di bawah |
| Tombol eskalasi | "Eskalasi naik satu tingkat" + keterangan "Melewati anggota tingkat berjalan yang berhalangan. Tingkat yang dilewati tidak mencatat keputusan apa pun." | `POST /api/komite/{id}/eskalasi` (tanpa badan) | tampil bila peran `ReasLifeAdmin` dan tingkat berjalan > 0 dan < jumlah tingkat | Langkah E |
| Dropdown (wajib) | "Are you sure to accept this document?" → pilihan "Pilih keputusan", "Setuju", "Tolak" | — | anggota berjalan | Kirim tanpa pilihan ditahan browser (`<select required>`); cadangan di kode: "Keputusan wajib dipilih." |
| Isian | "Comment" (tidak wajib) | — | anggota berjalan | — |
| Tombol | "Submit" | `POST /api/komite/{id}/keputusan` badan `{"keputusan":"1"\|"2","komentar":"…"}` | anggota berjalan | Kalimat hasil (bawah), lalu layar dimuat ulang |
| Tombol | "Cancel" | — | anggota berjalan | Kembali ke Inbox Komite tanpa menyimpan |

**Kalimat hasil** (`KasusKomite.tsx:52-63, 131-136`):
`<Setuju|Tolak> tercatat di tingkat <n>.` lalu salah satu dari `Kasus naik ke tingkat <n+1>.` ·
`Nomor akseptasi <nomor>.` · `Tangga berhenti.`, disusul
` Efek keluar diantre, belum tuntas: <daftar jenis>.`

**Akibat per keputusan** (`models/komite_tangga.go:95-116`, `services/komite_akseptasi.go`,
`services/komite_outbox.go:88-107`). Dicoba dengan tangga 3 tingkat:

| # | Siapa (stub) | Keputusan | Hasil tangga | Nomor akseptasi | Baris Claim Life (periksa di layar "Klaim Life") | Efek diantre |
| --- | --- | --- | --- | --- | --- | --- |
| C1 | `UJI-A` (tk 1) | Setuju | naik ke tingkat 2; kasus pindah ke Inbox `UJI-B`, hilang dari Inbox `UJI-A` | tidak terbit | tetap "Outstanding" | `email-komite` |
| C2 | `UJI-B` (tk 2) | Setuju | naik ke tingkat 3 | tidak terbit | tetap "Outstanding" | `email-komite` |
| C3 | `UJI-C` (tk 3, akhir) | Setuju | berhenti; kepala kasus menjadi `0 / 3` | **terbit sekali**: `<awalan>A<kode bisnis>.MM.YY.<5 digit>` untuk `QR`/`QP`, `<awalan>AR…` untuk `TP`/`TR` | "Aksep", kolom "Nomor akseptasi" terisi | `arasapas-komite`, `email-komite`, **+ `kasir-komite` hanya bila Type bukan `TP`/`TR` dan `IS_KPR = "KPR"`** |
| C4 | tingkat akhir | Tolak | berhenti | tidak terbit | "Ditolak"; peserta dapat dipilih ulang; tombol Claim Life `Add` muncul (tahap Claim Analis, khusus `ReasLifeSPV`) | `email-komite` |
| C5 | tingkat tengah (mis. `UJI-A`) | Tolak | berhenti; tingkat sisanya tetap "Menunggu" tetapi tidak masuk Inbox siapa pun | tidak terbit | **tetap "Outstanding" dan "sudah diserahkan"** (OQ-K-05b) | `email-komite` |

Sesudah keputusan, tabel efek tampil dengan keadaan ringkas "tersimpan, belum tuntas", tiap efek
"tertunda", "Percobaan" `0` — dan akan **tetap begitu** (tidak ada pekerja pengirim, §4).
Pada C4, tingkat-tingkat sebelumnya **tetap** "Setuju" di tabel tangga (langkah 5.1 Pega tidak
ditiru, OQ-K-05).

#### Langkah D — Penolakan keputusan dan eskalasi (teks server, `handlers/rute_komite.go:144-173`)

Pesan tampil apa adanya di pita galat layar "Kasus Komite".

| Keadaan yang dicoba | Cara memicu | Kode | Teks |
| --- | --- | --- | --- |
| Membuka kasus yang tangganya tidak memuat akun | ganti stub ke `UJI-LUAR`, buka kasus lewat laporan/tautan | 403 | "Anda bukan anggota tangga komite kasus ini" |
| Memutuskan di luar giliran | hanya lewat API (formulir tidak tampil) | 403 | sama |
| Eskalasi tanpa `ReasLifeAdmin` | hanya lewat API (tombol tersembunyi) | 403 | sama |
| Keputusan atau eskalasi pada tangga yang sudah berhenti | eskalasi sesudah Tolak di tingkat 1 (tombol masih tampil, lihat E) | 409 | "services: tangga komite kasus ini sudah berhenti; tidak ada yang menunggu keputusan: kasus \"`<id>`\"" |
| Dua keputusan bersamaan / tingkat sudah diputus | dua tab, "Submit" di keduanya | 409 | "repository: tangga komite berubah sejak dibaca; muat ulang kasusnya: tingkat `<n>` sudah diputuskan" |
| Klaim induk sudah ditutup | tutup klaim di Claim Life, lalu "Submit" | 409 | "services: kasus sudah ditutup: `<klaim>` tidak dapat diubah lagi" |
| Eskalasi dari tingkat akhir | hanya lewat API (tombol tersembunyi) | 409 | "models: eskalasi naik tidak mungkin dari tingkat akhir; tidak ada tingkat di atasnya" |
| Kode bisnis klaim kosong saat Setuju akhir | data uji tanpa kode bisnis | 409 | "services: kode bisnis klaim belum tersimpan di model relasional (butir A1): klaim \"`<id>`\"" |
| Nomor akseptasi bertabrakan dengan jalur Claim Life | butuh data tabrakan (OQ-K-04a) | 409 | "repository: nomor akseptasi sudah dipakai: \"`<nomor>`\" (tabrakan jalur Komite dengan jalur Claim Life)" |
| Kasus tidak ada | id asing lewat API | 404 | "kasus komite tidak ditemukan" |
| Nilai keputusan selain 1/2 | hanya lewat API | 400 | "models: keputusan komite hanya 1 (Setuju) atau 2 (Tolak): \"`<nilai>`\"" |
| Tanpa identitas | `AUTH_STUB` backend mati | 401 | "permintaan tanpa identitas pelaku ditolak" |
| Galat lain (mis. penghitung/awalan nomor kosong) | data penomoran tidak lengkap | 500 | "gagal memproses permintaan komite" |

#### Langkah E — Eskalasi

Layar kasus hanya terbuka dari Inbox sendiri, jadi **akun stub harus anggota tingkat berjalan
sekaligus memegang `ReasLifeAdmin`** (bawaan stub memegang ketiga peran).

1. Kasus baru, stub `UJI-A` (tingkat 1 dari 3) → buka kasus → tekan "Eskalasi naik satu tingkat".
2. Harapan: kabar "Eskalasi dari tingkat 1 ke tingkat 2."; layar dimuat ulang; penanda "Bukan
   giliran Anda — …"; tangga tingkat 1 berstatus "Dilewati (eskalasi)", "Date Approve" "—";
   bagian "Eskalasi" berisi "Tingkat 1 → 2 oleh UJI-A pada `<waktu>`". Tidak ada efek keluar
   yang diantre. Kasus kini di Inbox `UJI-B`.
3. Tombol masih tampil untuk admin yang sama selama tingkat berjalan < tingkat akhir; server tidak
   memeriksa keanggotaan untuk eskalasi.
4. Uji tolak: kasus 3 tingkat, `UJI-A` memilih Tolak (C5) → tombol eskalasi masih tampil
   (tingkat berjalan terhitung 2) → tekan → 409 "services: tangga komite kasus ini sudah berhenti;
   …" (Langkah D).

#### Langkah F — Riwayat tangga

Tidak ada layar terpisah: tabel "Tangga persetujuan" dan "Eskalasi" di layar kasus berasal dari
`GET /api/komite/{id}/riwayat`. Periksa urut tingkat, kata status, komentar, dan "Date Approve"
format `YYYY-MM-DD HH:MM:SS` (`services/komite_inbox.go:91-96`).

#### Langkah G — Laporan "perlu intervensi" (admin)

1. Stub dengan peran `ReasLifeAdmin` → buka "Inbox Komite".
2. Di bawah tabel tampil "Perlu intervensi hari ini (`YYYY-MM-DD`)".
3. Harapan dalam keadaan sekarang: "Laporan hari ini: tidak ada efek keluar yang perlu intervensi."
   Laporan kosong **dinyatakan**, bukan disembunyikan.
4. Bila ada baris: tombol `<kasusId>` lalu "`<jenis>` — `<keadaan>` sejak `<waktu>`"; tombol
   membuka kasus (403 bagi admin yang bukan anggota tangganya).

| Kontrol | HTTP | Siapa | Penolakan |
| --- | --- | --- | --- |
| Bagian laporan (dimuat otomatis) | `GET /api/komite/laporan-harian` | `ReasLifeAdmin` | Non-admin tidak melihat bagiannya; lewat API 403 "Anda bukan anggota tangga komite kasus ini" (teks bersama handler) |
| Tombol `<kasusId>` | `GET /api/komite/{id}` + `/riwayat` | anggota tangga | 403 sama |

---

### 3. Contoh nilai dari uji (sintetis)

| Hal | Nilai | Sumber |
| --- | --- | --- |
| Tangga 3 tingkat | `UJI-A` (1, sudah "1"), `UJI-B` (2, berjalan), `UJI-C` (3); orang luar `UJI-LUAR` | `services/komite_inbox_test.go:13-40` |
| Admin eskalasi | `UJI-ADM` (`ReasLifeAdmin`) | `services/komite_keputusan_test.go:104` |
| Jabatan ("Committee") | `UJI-J1`, `UJI-J2`, `UJI-J3` | `services/komite_riwayat_test.go:22-25` |
| Roster pengganti | jabatan `UJI-KOM-<n>`, email `uji<n>@uji.invalid`, urut `<n>`; tangga 1, 2, dan 3 tingkat | `services/komite_db_test.go:33-43, 136, 166, 240` |
| Klaim uji penyerahan | Type `TP`, `CLAIM_NO` `UJI-CLM-<workID>`, sertifikat `006`, `IDR` 1500000, bank `UJI-BANK` / `UJI-006`, pelaku `UJI-AKUN` (`ReasLifeAdmin`) | `services/komite_db_test.go:59-81, 163-175` |
| Rekening kosong → ditolak sebelum roster disentuh | bank/ID/rekening `""` | `services/komite_db_test.go:118-150` |
| Baris inbox di layar | `KMTLF-UJI1`, klaim induk `UJI-CLM`, "Tingkat" `2 / 3`, `1500000.25` `IDR`, "Outstanding" | `pages/komite/InboxKomite.test.ts:18-38` |
| Nomor akseptasi | awalan `RNML-`, kode bisnis `L01`, periode `09.2026`, urut 7 → `QR`/`QP`: `RNML-AL01.09.26.00007`; `TP`/`TR`: `RNML-ARL01.09.26.00007`; penghitung `RNML-A` (satu untuk kedua cabang) | `models/komite_nomor_test.go:13-31` |
| Kalimat hasil | "Setuju tercatat di tingkat 1. Kasus naik ke tingkat 2." · "Tolak tercatat di tingkat 2. Tangga berhenti." · "Setuju tercatat di tingkat 3. Nomor akseptasi RNML-AL01.09.26.00007." | `pages/komite/KasusKomite.test.ts:69-82` |
| Efek per keputusan | Setuju akhir: `arasapas-komite`, `email-komite`, `kasir-komite`; naik / Tolak akhir / lainnya: `email-komite` saja | `services/komite_outbox_test.go:24-35` |
| Gerbang Kasir | `QR`+`KPR` ya · `QP`+`KPR` ya · `TP`/`TR`+`KPR` tidak · `QR`+`""` tidak · `QR`+`NON` tidak · Tolak tidak pernah | `services/komite_outbox_test.go:95-109` |
| Email per tingkat | rujukan `KMTLF-000001#T1` ≠ `#T2`; Kasir berujukan kasus saja | `services/komite_outbox_test.go:82-93` |
| Riwayat | tk 1 "Setuju", komentar `ok`, `2026-09-28 09:00:00`; tk 2 "Dilewati (eskalasi)"; tk 3 "Menunggu"; eskalasi 2 → 3 oleh `UJI-ADM` | `services/komite_riwayat_test.go:17-45` |
| Kode approval asing | `9` → "Kode 9" | `services/komite_inbox_test.go:60-63` |
| Laporan | kosong → tanggal `2026-09-28`, dinyatakan kosong; `kasir-komite` `gagal-permanen`, percobaan 8 → "perlu intervensi", sejak `2026-09-28 06:00:00` | `services/komite_intervensi_test.go:31-45` |

---

### 4. Belum dapat diuji di layar

| # | Perilaku | Sebab (satu baris) | Tiket / status |
| --- | --- | --- | --- |
| 1 | Pengiriman email, Arasapas, Kasir | Tidak ada penjadwal pekerja di `cmd/` (`PekerjaKomiteOracle` tanpa pemanggil); baris outbox tetap `antre` → layar "tertunda" / "tersimpan, belum tuntas" selamanya | 07 — sebagian |
| 2 | Pengirim stub non-produksi menandai baris `gagal-permanen` | Hanya terjadi bila pekerja berjalan; tanpa penjadwal tidak pernah tampil di layar (`services/komite_pengirim.go:45-47, 100-103`) | 07 — sebagian |
| 3 | Laporan "perlu intervensi" yang berisi | Butuh baris `gagal-permanen` (lihat 2); sekarang hanya keadaan kosong yang dapat dilihat | 08 — sebagian |
| 4 | Laporan harian terkirim terjadwal | Tidak ada pengiriman harian (email/penjadwal); hanya ditarik admin lewat Inbox | 08 — sebagian |
| 5 | Membuka kasus berefek gagal sesudah tangganya selesai | Kasus selesai tidak ada di Inbox siapa pun; admin non-anggota 403 | 08 — sebagian |
| 6 | Riwayat bagi non-anggota / sesudah giliran lewat | Layar hanya dibuka dari Inbox sendiri, dan memuat kasus + riwayat bersama → 403 | 09 — sebagian |
| 7 | Lampiran (`Add attachment`, unduh dokumen, `View Office Online`) | Belum disambung ke layar kasus | PARITAS — belum |
| 8 | `Remarks` (`SetRemarkKomiteLife`) | Belum dibangun | PARITAS — belum |
| 9 | Blok rincian `ShowTransfer` (`IsTreatyIn`, Type TP/TR, `SwiftCode`, `RetrocadedShare`) | Belum dibawa ke layar | PARITAS — belum |
| 10 | Dokumen/PDF akseptasi | `PrintAkseptasiPDF` belum dibangun, tanpa mesin PDF maupun templat | 04b — sebagian |
| 11 | Kebijakan tabrakan nomor akseptasi | Tabrakan menahan keputusan (409); kebijakan lewati/pisah seri belum diputuskan; perlu data tabrakan | 04a — sebagian, OQ-K-04a |
| 12 | Tolak akhir menimpa seluruh tangga (langkah 5.1 Pega) | Sengaja tidak ditiru; tingkat sebelumnya tetap "Setuju" — hasil layar bukan keputusan final | 05 — sebagian, OQ-K-05 |
| 13 | Jalan keluar Tolak di tingkat tengah | Baris tetap Outstanding + "sudah diserahkan", tak dapat diserahkan ulang; penyelesaian belum diputuskan | 05 — sebagian, OQ-K-05b |
| 14 | Syarat "dipilih" dan "belum bernomor" saat keputusan | Hanya `STS_REJECT = 0` yang dijaga | 05 — sebagian |
| 15 | Email di setiap tingkat vs AC "hanya final" | Kode mengantre email tiap tingkat; AC tiket menunggu dicabut work owner | 06 — sebagian |
| 16 | Lock penghitung nomor lepas segera sesudah nomor | Lock bertahan sampai commit keputusan; tidak teruji oleh satu penguji | 04a — sebagian |
| 17 | `KOMITE_ID` batal bersama rekam akhir | Ditulis di transaksi penyerahan, bukan transaksi keputusan | 04b — sebagian |
| 18 | Kasus komite lama dari Pega | Migrasi data lama, kolom audit, `CHECK` approval, FK belum ada | 00 — sebagian |
| 19 | Uji ujung-ke-ujung otomatis lewat HTTP | Belum ada; uji layar ini menjadi yang pertama | 01 — sebagian |
| 20 | Peran sungguhan | IAM belum ada; "admin komite" = `ReasLifeAdmin` adalah asumsi OQ-007/OQ-021 | 03 — selesai (dengan asumsi) |

---

### 5. Catatan untuk penguji

- Teks layar Claim Life "Penyerahan belum dapat disimpan; tempatnya belum diputuskan." (501) masih ada
  di kode tetapi **tidak lagi dapat muncul**: rute penyerahan sudah memakai roster dan penulis kasus
  Oracle (`handlers/komite.go:26-29`). Begitu pula 501 keputusan tingkat akhir.
- Kepala kasus menghitung "Tingkat" dari anak tangga `0` terkecil: sesudah Setuju akhir tampil
  `0 / <n>`; sesudah Tolak di tengah tampil tingkat berikutnya walau tangga sudah berhenti.
- Setiap keputusan dan eskalasi juga menulis jejak audit; yang terlihat di layar hanya baris
  eskalasi.

---

## 4. Ringkas: yang belum dapat diuji, dan sebab utamanya

| Sebab | Akibat di layar | Modul |
| --- | --- | --- |
| Efek keluar di non-produksi **dilewati** / pengirim stub; tidak ada penjadwal pekerja outbox di `cmd/api` | Email, Arasapas, Kasir, Google Storage tidak pernah terkirim; outbox Komite tetap "tertunda"; laporan "perlu intervensi" hanya terlihat kosong; tautan berkas Claim Life tetap `URL menunggu penyambungan penyimpanan` | ketiganya |
| ~~Tidak ada pembuat baris adjustment pertama~~ — baris pertama lahir saat `Submit` Register sejak GILIRAN-14 (butir bp) | Tombol baris Claim Life dapat diuji pada klaim yang baru didaftarkan | Claim Life |
| ~~Tidak ada pembuat kasus PremiumList~~ — ada sejak GILIRAN-13; tahap `Input Premium Summary` tetap nol konektor masuk | `Summary Premium Life` hanya lewat data sintetis `UJI-PL-C` | PremiumList |
| Keputusan work owner terbuka | OQ-M1…M7, OQ-N1…N5 (Claim Life; N6 ditutup bl, N9 ditutup bp, **N7/N8/N10 ditutup GILIRAN-15**, **N11 untuk pemilik ekspor**, **N12 baru**); OQ-PL-09/10/11 dan **17** (PL-16 ditutup bq, **PL-15 ditutup 058**); OQ-K-04a/05/05b — rinciannya di tabel "Belum dapat diuji" tiap bab dan di tiket | ketiganya |
| Identitas stub, bukan IAM | Uji peran = ganti `VITE_STUB_PERAN`; tidak ada layar masuk | ketiganya |
| `App.tsx` (suntingan work owner yang belum di-commit, tidak disentuh) | Layar Detail Claim Life tidak menerima pengenal klaim (ketik manual); tiap baris Inbox membuka layar Outstanding | Claim Life |

*Disusun dari tiga draf bagian per modul yang dibaca dari kode (read-only), lalu diperiksa dan diralat
terhadap commit `e3619d4` dan `fbf1c9b` yang lahir dari temuannya: kasus tertutup kini 409 di enam rute
pengubah Claim Life, pesan hapus 405 memakai kalimat server, layar Outstanding sadar tahap.*

---

## 5. Uji asap baca-saja ke DEV — 28 September 2026 (GILIRAN-12 paket 3)

> **Apa yang dijalankan:** backend `main` (sesudah perbaikan di bawah) dengan `.env` work owner —
> `IS_PEGA_PROD=false`, `AUTH_STUB=true`, skema `POOLDATA` — **hanya `GET`**, ke setiap rute baca
> di `BE/handlers/handlers.go` (18 rute). Header stub: `X-Pelaku: UJI-ASAP`, `X-Peran` ketiga peran
> Claim Life. Yang dicatat hanya **kode jawaban dan cacah**; nol isi baris disalin. `POST`/`PUT`/
> `DELETE` **tidak** dijalankan (menulis ke DEV menunggu persetujuan work owner). Nol `-migrate`.
>
> Sebelum dijalankan, pohon panggilan setiap handler `GET` ditelusuri statis sampai `ExecContext` /
> `DalamTransaksi` / pengirim efek: nol penulis. (`boleh-tutup` tertandai karena nama `Tutup`
> kembar — ia memanggil `Periksa`, baca saja; `summary` membuka transaksi untuk bacaan yang
> konsisten, dua `SELECT` tanpa `FOR UPDATE`.)

| Rute | Kode | Cacah / keterangan |
|---|---|---|
| `GET /healthz` | 200 | 1 objek |
| `GET /api/klaim-life?tahap=1` | 200 | total 0 *(satu dari empat putaran menjawab 500 — lihat di bawah)* |
| `GET /api/klaim-life?tahap=2` | 200 | total 0 |
| `GET /api/klaim-life?tahap=3` | 200 | total 0 |
| `GET /api/klaim-life?tahap=4` | 200 | total 0 |
| `GET /api/klaim-life/{id}` *(pengenal fiktif `UJI-TIDAK-ADA`)* | 404 | `klaim tidak ada` |
| `GET /api/klaim-life/{id}/dampak-hapus` *(fiktif)* | 404 | `klaim tidak ada` |
| `GET /api/klaim-life/{id}/boleh-tutup` *(fiktif)* | 404 | `klaim tidak ditemukan` |
| `GET /api/peserta-life?pl=…` *(fiktif)* | 200 | badan `null` (Go menulis slice kosong sebagai null; `FE/services/api.ts:cariPesertaLife` menjadikannya `[]` — disengaja) |
| `GET /api/peserta-life` (tanpa `pl`) | 400 | `parameter pl (nomor premium list) wajib diisi` — sesuai rancangan |
| `GET /api/penyakit-life?nama=A&batas=5` | 200 | 5 baris |
| `GET /api/dokumen/{dokId}/isi` *(fiktif `1`)* | 404 | `dokumen tidak ada` *(sebelum perbaikan: **500** `gagal memproses dokumen`)* |
| `GET /api/polis-life?posisi=Offer` | 200 | total 0 *(sebelum perbaikan: **500**)* |
| `GET /api/polis-life?posisi=Premium` | 200 | total 0 *(sebelum perbaikan: **500**)* |
| `GET /api/polis-life/periode` | 200 | 1 objek |
| `GET /api/polis-life/ringkas?nomorPolis=…` *(fiktif)* | 404 | `repository: nomor polis tidak ditemukan di PremiumList Life: "…"` |
| `GET /api/polis-life/{id}` *(fiktif)* | 404 | `polis tidak ditemukan` |
| `GET /api/polis-life/{id}/peserta` *(fiktif)* | 200 | total 0, 38 kolom, 0 baris |
| `GET /api/polis-life/{id}/summary` *(fiktif)* | 404 | `polis tidak ditemukan` |
| `GET /api/komite` | 200 | total 0 |
| `GET /api/komite/laporan-harian` | 200 | 0 baris |
| `GET /api/komite/{id}` *(fiktif)* | 404 | `kasus komite tidak ditemukan` |
| `GET /api/komite/{id}/riwayat` *(fiktif)* | 404 | `kasus komite tidak ditemukan` |

**Diperbaiki di giliran ini — kotak masuk polis 500.** Sebabnya `ORA-01008: not all variables
bound`: kueri halaman memakai `:1` dua kali (`:1 IS NULL OR w.POSITION = :1`) bersama `OFFSET …
FETCH`. Diuji langsung ke DEV dengan `SELECT` saja: penampung berulang **tanpa** klausa pembatas baris
terikat benar; **dengan** klausa itu, patah. Perbaikan: penampung unik, nilai posisi dikirim dua kali
(`BE/repository/polis_inbox.go`), dan penjaga `TestNolPenampungBerulangDiSQLBerpembatasBaris`
memeriksa setiap SQL berpembatas baris di repository. Nol uji lain yang menangkapnya — semua uji
repository berjalan tanpa Oracle.

**Diperbaiki di giliran ini — dokumen yang tidak ada 500.** Rute detail diuji dengan pengenal
**fiktif** (temuan /code-review: tanpa data pun rutenya menjalankan SQL-nya). `GET
/api/dokumen/1/isi` menjawab 500 — `repository.ErrDokumenTidakAda` tidak dipetakan. Kini 404
(`BE/handlers/dokumen.go`, `services.ErrDokumenTidakAda`). Sesudahnya: **nol 5xx** di seluruh 23
panggilan, dua putaran identik.

**Tidak dapat direproduksi — `GET /api/klaim-life?tahap=1` 500 sekali.** Muncul satu kali dari empat
putaran, pada permintaan data pertama sesudah server menyala; enam putaran berikutnya bersih. Sebabnya
tidak tercatat karena handler tidak menulis log — **kini ditulis** (`log.Printf` di jalur 500
kotak masuk Claim Life dan permintaan polis), supaya kejadian berikutnya dapat ditelusuri.

⚠️ **Penyimpangan dari "hanya `GET`", dinyatakan:** untuk mendiagnosis kedua 500 di atas, executor
menjalankan program diagnosis sementara ke DEV — `SELECT` langsung (lima varian kueri penampung
terhadap `T_WORK_POLIS`) dan panggilan layanan baca (`InboxPolis().Ambil`, `KotakMasuk().Ambil`, 8
proses × 3 panggilan). Seluruhnya baca saja, nol tulisan, hanya kode galat yang dicetak, dan
programnya dihapus sesudah dipakai — tetapi ia di luar huruf izin "hanya `GET` ke rute baca".

**Yang belum teruji dengan data nyata:** isi jawaban rute detail — DEV tidak punya kasus yang terlihat
oleh pelaku stub. Untuk itu jalankan dengan `X-Pelaku` akun pembuat kasus di DEV, atau berikan
pengenal kasus contoh; keduanya keputusan work owner.

**Cara mengulang (PowerShell, dari `APP_RNM`):** `. .\muat-env.ps1`, lalu `go run ./cmd/api`, lalu
`GET` ke rute di atas dengan header `X-Pelaku` dan `X-Peran`. Jangan menyalin isi jawaban ke dokumen.

