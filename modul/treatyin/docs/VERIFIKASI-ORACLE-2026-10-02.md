# Verifikasi Oracle — 2 Oktober 2026

**`L-3` SUDAH USANG.** Ia berbunyi *"tidak ada instans Oracle yang terjangkau"*, dan atas dasar itu
seluruh pernyataan keputusan di modul ini menulis bahwa DDL-nya belum pernah dijalankan di mana pun.
**Itu tidak lagi benar.** Berkas ini mencatat apa yang terjadi, apa yang terbukti, dan apa yang
masih menghalangi.

---

## 1 · Yang ditemukan, bukan yang dikerjakan

| Hal | Keadaan |
| --- | --- |
| Docker / podman | **tidak ada** di mesin ini — `make db-up` tidak dapat dijalankan |
| `sqlplus` | **ada** |
| Instans Oracle | **terjangkau** — `DEV_NUSARE2`, TCP 1521 terbuka |
| `ORACLE_SCHEMA` di `.env` | **`POOLDATA`** — skema warisan Pega: 816 tabel, 428 sequence, 369 view, 198 procedure |
| `ORACLE_SKEMA_UJI` | **barisnya tidak ada** |

⛔ **Migrasi `400`–`419` dan `440`–`441` SUDAH DIJALANKAN terhadap POOLDATA**, pada 2 Oktober 2026
pukul 11:13–11:14, **bukan oleh pekerjaan ini**. Ke-22 langkah tercatat di `T_MIGRASI`; ke-29 tabel
dan ke-29 sequence berdiri; seluruhnya **kosong** (nol baris).

### Kenapa itu layak dicatat, bukan dilewati

`make migrate` **tidak berpagar**. `cmd/api/main.go` menjalankan `jalankanMigrasi(dasar)` tanpa
memeriksa `IS_PEGA_PROD` dan tanpa `PastikanSkemaUji()` — sementara `-migrate-down` **berpagar
keduanya**. Akibatnya asimetris, dan tidak dapat dibatalkan dengan perkakas aplikasi sendiri:

> Migrasi maju **boleh** menulis ke POOLDATA; migrasi mundur **menolak** membongkarnya dari sana,
> sebab POOLDATA bukan skema uji.

Ke-29 tabel Treaty In kini berdiri di skema warisan Pega, dan `-migrate-down` tidak akan
membuangnya. Membuangnya menuntut `DROP TABLE` manual oleh DBA.

**Pekerjaan ini TIDAK menjalankan `make migrate`,** dan tidak akan — menulis DDL ke skema warisan
bersama bukan keputusan pengembang modul. Yang dijalankan di sini hanya **baca** dan **`INSERT` yang
diakhiri `ROLLBACK`**.

---

## 2 · Yang TERBUKTI — nol `ORA-` pada DDL

Ke-22 langkah migrasi kedua modul diterima Oracle tanpa satu galat pun. Yang ikut terbukti sah, dan
sebelumnya hanya dapat diargumentasikan:

| Yang diragukan | Hasil |
| --- | --- |
| `LIMIT` sebagai nama kolom pada `LAYER` | **diterima** — ia kata kunci PL/SQL, bukan kata cadangan SQL |
| `CK_POTONGAN_INDUK` sebagai `ALTER TABLE` tersendiri | **diterima** |
| `ALTER TABLE … MODIFY NOMOR_URUT_VERSI NULL` (migrasi `441`) | **diterima** — kolomnya kini `NULLABLE=Y` |
| Kunci asing menunjuk tabelnya sendiri | **diterima** |
| Pengenal terpanjang | **diterima**, seluruhnya ≤ 30 bita |
| Presisi `NUMBER(38,8)` menggantikan `NUMBER(38,20)` | **diterima** |

Inventaris constraint yang berdiri: **29 PRIMARY KEY · 22 UNIQUE · 33 FOREIGN KEY · 1 CHECK bernama**.

---

## 3 · Yang DIUJI PERILAKUNYA — `INSERT` lalu `ROLLBACK`

DDL di Oracle auto-commit; **DML tidak**. Seluruh uji di bawah menyisipkan baris, membaca
penolakannya, lalu `ROLLBACK` — **nol perubahan menetap**.

### Negatif — constraint sungguh menolak

| Invarian | Yang dicoba | Hasil |
| --- | --- | --- |
| `INV-04` | dua versi bernomor urut sama pada satu kontrak | `ORA-00001 (UQ_VERSI_KONTRAK)` |
| kunci asing | versi yatim, kontrak induk tidak ada | `ORA-02291 (FK_VERSI_KONTRAK_1) parent key not found` |
| `INV-44` | mata uang kontrak tidak ada di tabel acuan | `ORA-02291 (FK_VERSI_KONTRAK_MATA_UANG)` |
| `INV-68` | `KODE` ganda di tabel acuan | `ORA-00001 (UQ_MATA_UANG)` |
| `INV-64` | dua `BAGIAN` pada satu `LAYER` | `ORA-00001 (UQ_BAGIAN)` |
| `INV-05` | dua layer tanpa bagian, nomor sama | `ORA-00001 (UQ_LAYER)` |
| `INV-05` | dua layer berbagian, nomor dan bagian sama | `ORA-00001 (UQ_LAYER)` |
| `KTV-B` | potongan dengan **kedua** induk terisi | `ORA-02290 (CK_POTONGAN_INDUK)` |
| `KTV-B` | potongan dengan **nol** induk terisi | `ORA-02290 (CK_POTONGAN_INDUK)` |
| `INV-15` | jenis potongan ganda pada pelekatan yang sama | `ORA-00001 (UQ_POTONGAN)` |

### Positif — data yang sah DITERIMA

Uji positif ada sebab **constraint yang menolak terlalu banyak lulus setiap uji negatif yang pernah
ditulis untuknya**.

| Tiket | Yang dibuktikan | Hasil |
| --- | --- | --- |
| `14` | satu kontrak dengan **tiga** versi berturut-turut | ketiganya diterima; `KONTRAK` tetap **satu** baris — lapisan beku tidak disalin |
| `37` | dua potongan, **jenis sama**, pelekatan berbeda | keduanya diterima — kedua `UNIQUE` tidak saling mencampur |

### Yang terbukti BELUM ditegakkan — dan itu memang yang tertulis

| Invarian | Yang dicoba | Hasil |
| --- | --- | --- |
| `INV-53` | kontrak bertanggal berakhir lebih awal daripada tanggal mulai | **DITERIMA** — sesuai pernyataan keputusan di migrasi `401` |

---

## 4 · Satu klaim saya yang TERBUKTI KELIRU, dan dicabut

Migrasi `413`, dokumen STRUKTUR, dan `MODUL.md` pernah menyatakan:

> ~~`UQ_LAYER` tidak menegakkan `INV-05` sepenuhnya: `BAGIAN_LAYER` boleh kosong, dan Oracle
> memperlakukan NULL sebagai tidak sama dengan NULL di kunci unik komposit, sehingga baris
> tanpa bagian diterima berulang.~~

**Keliru.** Aturan Oracle yang sebenarnya: sebuah entri dilewati indeks unik **hanya bila SELURUH
kolom kuncinya NULL**. Pada `UQ_LAYER`, `ID_VERSI_KONTRAK` dan `NOMOR_LAYER` selalu terisi, sehingga
barisnya terindeks dan duplikatnya ditolak.

Kekeliruan ini lolos **dua kali** — dari penulisnya dan dari peninjau sumbu Standards — dan yang
memperbaikinya bukan argumen, melainkan satu `INSERT`.

> **Pelajarannya, dan ia berlaku untuk seluruh modul:** jangan menulis pernyataan keputusan tentang
> **perilaku** basis data yang belum dijalankan. Pernyataan tentang **bentuk** boleh; tentang
> perilaku, tunggu Oracle.

Daftar periksa tiket `31` *"INV-05 terpasang atas NOMOR_LAYER + BAGIAN_LAYER di dalam satu versi"*
karena itu **terpenuhi**, bukan gagal.

---

## 5 · Yang masih menghalangi, dan siapa pemiliknya

`L-3` tertutup, tetapi **`make test-db` belum dapat dijalankan di sini**. `uji/skemauji` menolak
POOLDATA dengan **galat, bukan skip**, dan penolakan itu benar: `Pasang`/`Bongkar` menghapus tabel di
skema yang ditunjuk, dan POOLDATA memuat 816 tabel warisan sungguhan.

> ⚠️ **RALAT.** Tabel di sini pernah mendaftar empat butir dengan kolom **"Pemilik: DBA"** dan
> memperlakukannya sebagai penghalang di luar kendali. **Tidak ada DBA di proyek ini** — tim ini
> yang membuat seluruh tabel sendiri. Bingkai itu membuat pekerjaan sendiri terbaca sebagai
> permintaan yang sedang menunggu jawaban orang lain, dan itu salah.

| Yang dibutuhkan | Keadaan 2 Oktober 2026 |
| --- | --- |
| **skema uji tersendiri** di `DEV_NUSARE2` | ⛔ **TERHALANG HAK AKSES.** Akun aplikasi hanya punya `CREATE SESSION`, `CREATE VIEW`, `UNLIMITED TABLESPACE`, `CONNECT`, `RESOURCE` — **nol `CREATE USER`**. Membuat skema menuntut kredensial yang lebih tinggi, dan kredensial itu tidak ada di `.env` maupun di repo. Ini bukan tugas yang dilewati; ini tugas yang **tidak dapat dijalankan dengan kredensial yang tersedia**. |
| `ORACLE_SCHEMA` diarahkan ke sana, `ORACLE_SKEMA_UJI=true` | menunggu butir di atas |
| ke-29 tabel **dibuang dari POOLDATA** | ✅ **skripnya ditulis**: [`../alat/buang-tabel-treatyin-dari-pooldata.sql`](../alat/buang-tabel-treatyin-dari-pooldata.sql). **Sengaja belum dijalankan** — tabelnya dibuat orang lain hari itu juga, POOLDATA skema bersama, dan sesudah dibuang migrasi tidak dapat dijalankan ulang ke tempat yang benar selama butir pertama belum tertutup |
| pagar pada jalur `-migrate` maju | ✅ **selesai** — `cmd/api/main.go` kini menolak `IS_PEGA_PROD=true` dan `ORACLE_SCHEMA` yang memuat nama skema warisan. Penyimpangan §10 di `KEPUTUSAN-PENYELARASAN-REPO.md` |

**Yang perlu diputuskan pemilik proses, dan ia satu kalimat:** kredensial mana yang boleh dipakai
untuk `CREATE USER` di `DEV_NUSARE2` — atau, bila tidak ada, apakah modul ini boleh memakai satu
skema aplikasi yang sudah berdiri selain POOLDATA.

Begitu skema uji ada, `make test-db` menjalankan
`modul/treatyin/backend/repository/invarian_db_test.go`, yang **mengabadikan seluruh uji §3 sebagai
uji otomatis** — negatif dan positif, masing-masing menuntut constraint **yang disebut namanya**,
sebab penolakan karena alasan lain adalah uji yang lulus secara kebetulan.

---

## 6 · Apa artinya bagi kolom "selesai" di papan

Tetap **nol**, dan itu jujur. Yang berubah: penghalangnya bukan lagi *"tidak ada Oracle"*, melainkan
*"tidak ada skema uji"* — satu kalimat dari DBA, bukan satu lingkungan yang harus dibangun.

Sembilan constraint kini terbukti menolak dan dua jalur positif terbukti diterima, tetapi **tiket
tidak dinyatakan selesai dari berkas ini**: status adalah medan di dalam berkas tiket, dan yang
berwenang mengubahnya pemilik proses. Yang dapat dikatakan berkas ini: **buktinya sudah ada**.
