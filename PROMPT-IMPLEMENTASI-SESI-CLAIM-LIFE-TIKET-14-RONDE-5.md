# PROMPT — sesi implementasi berikutnya: Claim Life tiket 14 **ronde 5** — pra-terbang bentuk tabel, tipe kolom warisan dari katalog, dan penutupan bila Oracle ada

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya. Brief ronde 4
> **`PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14-RONDE-4.md`** tetap berlaku untuk hal yang tidak
> diubah di sini *(§2 k–l, §4 Langkah A dan D, §6, §7, dan seluruh §9)*.
>
> **Beda dari ronde 4:** sesi ini **punya pekerjaan nyata tanpa Oracle dan tanpa keputusan baru**
> — §3-1 dan §3-2. Gerbang §1 hanya menentukan seberapa jauh sesi boleh melangkah, bukan apakah ia
> boleh dimulai.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> tiket 14 bab `## Implementasi — ronde 4, 26 September 2026` → brief ronde 4 **§9** →
> `.scratch\claim-life\TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md` dan `SUMBER-PENOMORAN-DBA.md`
> **seluruhnya** → §2–§3 berkas ini.
>
> **SESI INI:** claim-life · tiket **14 ronde 5** + tiket **01 (penutupan, bila G1)** + **persiapan
> tiket 02 (bila G3)**

---

## 0. KEADAAN AWAL — 26 September 2026 sore

| | Keadaan |
| --- | --- |
| `HEAD` | `ec843e2` — tiket 14 *claimed* **40/53** AC, daftar terbuka 13 = kotak `[ ]`; tiket 01 *claimed* 4/7; tiket 02 *ready-for-agent* dengan blocker o (AC 2, 3, 10) |
| Uji yang lulus | `go vet`, `go vet -tags=db`, `gofmt` nol, `go build`; **83** test Go PASS; **19** test bertag `db` SKIP dengan pesan; `tsc --noEmit`; **5** test JS; `vite build` 87 modul |
| Keputusan diterapkan ronde 4 | **d** (`FK_WORK_COVER_KEY`, `FK_ADJ_KOMITE`, nullable, tanpa `ON DELETE`), **p1**, **t** (AC 12 dicentang), **e′** (catatan AC 8), **m** (blok ralat STRUKTUR), **n** (`006_t_claimlf_adj_spreading_retro.sql`) |
| Pagar skema uji | `skemauji` menolak dengan **galat** kecuali `ORACLE_SKEMA_UJI=true` **dan** skema tidak memuat `POOLDATA`; keenam pemanggil `Buka()` hanya SKIP untuk `ErrTanpaOracle`; bawaan `POOLDATA` dicabut dari `Makefile` dan `.env.example` |
| Migrasi | 16 berkas, 8 FK *(6 kepemilikan + 2 penunjuk ke atas)*, **belum pernah dijalankan di Oracle mana pun** — lima sesi |
| Belum di-commit di working tree *(bukan milik sesi executor mana pun)* | `frontend/vite.config.ts` + `frontend/.env.example` *(proxy dev `DEV_PROXY_TARGET`, diuji ujung ke ujung)*; `APP_RNM/muat-env.ps1` + `muat-env.cmd` *(pemuat `.env`, nol rahasia)*; `.scratch/claim-life/TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md` + `SUMBER-PENOMORAN-DBA.md` *(`[data DBA — dibaca sendiri]`)*; brief ronde 4 **§9**; `PANDUAN-MENJALANKAN.txt` bab 3–6 *(diperbarui untuk pemuat dan proxy)* |
| Katalog instance pengembangan | sudah dibaca *(brief ronde 4 §9)*: Oracle 12.2.0.1; `POOLDATA` 760 tabel; **`DOCUMENT_CLAIM` sudah ada, 14 kolom, 295 baris**; `OS_AKSEPTASI_KLAIM_LIFE` **62** kolom; sumber procedure penomoran terbaca utuh |

**Verifikasi independen 26 September 2026 sore atas commit `ec843e2`:** seluruh angka laporan ronde 4
tereproduksi *(83 · 19 · 5 · 40/53 · 13 · 20 berkas +781/−42)*; kedua FK dan ketiadaan `ON DELETE`
dibaca di DDL; pagar dua syarat dan `BolehDilewati` di enam pemanggil dibaca di diff; klaim korpus
*"nama tabel Pega di Claim Life ≤ 23 byte"* dihitung ulang dan **benar** *(terpanjang
`OS_AKSEPTASI_KLAIM_LIFE`, 23)*. Yang belum tepat ada di §3, dan **dua di antaranya adalah pekerjaan
sungguhan** yang tidak menunggu siapa pun.

---

## 1. GERBANG — menentukan seberapa jauh, bukan apakah mulai

| Gerbang | Yang membukanya | Keadaan 26-09 sore |
| --- | --- | --- |
| **G0 tanpa syarat** | — | **terbuka**: §3-1, §3-2, §3-3, Langkah 0 |
| **G1 Oracle** | user **kosong** baru di instance pengembangan + `ORACLE_DSN` + `ORACLE_SCHEMA` + `ORACLE_SKEMA_UJI=true`, disetel di terminal atau lewat `APP_RNM\.env` + `muat-env` | tertutup — DBA belum membuat user; `POOLDATA` **bukan** skema uji |
| **G2 keputusan tiket 14** | **v** *(nama tabel dokumen)* dan **w** *(`WPC`, `CLAIM_RETRO`)* | tertutup — keduanya `[USULAN]` sejak brief ronde 4 §9 |
| **G3 keputusan tiket 02** | **o1, o2, o3** disahkan. ⭐ Objek DBA §5-4 sampai §5-6 brief ronde 4 **sudah terpenuhi isinya** lewat `SUMBER-PENOMORAN-DBA.md`; yang tersisa hanya konfirmasi DBA bahwa produksi sama dengan pengembangan | tertutup — hanya oleh keputusan |

---

## 2. KEPUTUSAN WORK OWNER

| | Keputusan | Keadaan |
| ---: | --- | --- |
| v | **Nama tabel dokumen** — `DOCUMENT_CLAIM` sudah ada di `POOLDATA` (14 kolom: `ID VARCHAR2(100) NOT NULL`, `IDPEGA`, `TANGGAL`, `NAMAFILE`, `MIME`, `KATEGORI_1/2`, `NOAKSEP`, `NOPREKAS`, `PAYMENTDATE`, `INSKEY_LINK`, `INSKEY_DATA`, `T_STORAGE_ID`, `PXCREATEOPERATOR`). **(v1)** tabel baru diganti nama `T_CLAIMLF_DOCUMENT`, tabel warisan tidak disentuh; **(v2)** pakai tabel warisan apa adanya dengan pemetaan kolom, sesuai ADR-U-0042 *(bentuknya kini dapat direkayasa balik)* | `[USULAN]` — rekomendasi **v1** sekarang, v2 ditinjau saat tiket dokumen dikerjakan |
| w | **`WPC` dan `CLAIM_RETRO`** — tabel warisan: `WPC DATE`, `CLAIM_RETRO NUMBER`; DDL baru `003` menulis `WPC VARCHAR2(32)`, `002` `CLAIM_RETRO` *(tipe periksa)*; kode memperlakukan keduanya teks. Usulan: ikuti tabel warisan; `CLAIM_RETRO` sebagai `Ratio` atau `Money` **belum diketahui artinya** — tanya work owner; STRUKTUR diberi ralat bertanggal | `[USULAN]` |
| x | ⭐ **Pra-terbang bentuk sebelum `-migrate`.** Kelemahan "keberadaan diperiksa, bentuk tidak" (ronde 3) kini terbukti nyata lewat `DOCUMENT_CLAIM`. Usulan: sebelum menjalankan **langkah apa pun**, `JalankanMigrasi` memeriksa setiap `CREATE TABLE` yang tabelnya **sudah ada** di katalog: daftar kolom `ALL_TAB_COLUMNS` dibandingkan dengan daftar kolom DDL *(nama saja, urutan bebas)*. Berbeda → **menolak sebelum satu pernyataan pun dikirim**, dengan pesan menyebut tabel dan kolom yang berselisih. Sama → dilewati seperti sekarang. Ini bukan pagar `POOLDATA` *(produksi kelak memang menyasar skema warisan)*, melainkan pagar **bentuk** | `[USULAN]` — rekomendasi **disahkan**; §3-2 |
| y | **Semantik hapus dengan FK ke atas** — menghapus baris klaim yang masih ditunjuk baris komite akan ditolak `ORA-02292`. Dua jalan untuk tiket 15: lepaskan penunjuk dulu, atau tolak penghapusan selama komite ada | `[terbuka — tiket 15]`, bukan ronde ini |
| j (tinjau ulang) | Oracle pengembangan 12.2.0.1; bila DBA memastikan `COMPATIBLE ≥ 12.2` di **kedua** lingkungan, nama 36 byte sah. Executor sudah menurunkan dasarnya ke `[dugaan]` di STRUKTUR. Tetap 29 byte *(mengikuti sistem berjalan, ≤ 23 byte)* atau kembalikan ke 36? | `[USULAN]` — rekomendasi **tetap 29 byte**; bukan pekerjaan ronde ini |
| o1, o2, o3 | tetap sebagaimana brief ronde 4 §2, dengan **dua fakta baru** dari sumber procedure: **(a)** nol `COMMIT` di dalam procedure — OQ-013 punya bukti; **(b)** format lengkap `<prefix>K<kode>.MM.YYYY.<5 digit>` dirakit **pemanggil Pega**, procedure hanya mengeluarkan `LPAD(no_seq,5,'0')` dan `MM.YYYY` — jadi teks AC 3 tiket 02 harus mencakup **perakitan format**, bukan hanya counter | `[USULAN]` — G3 |
| k, l | tetap | berlaku |

---

## 3. TEMUAN VERIFIKASI RONDE 4 — dua pekerjaan, dua ralat kata

| # | Temuan | Letak | Yang dikerjakan |
| ---: | --- | --- | --- |
| 1 | ⭐ **Ronde 3 §3-7 kini dapat dikerjakan tanpa Oracle**: tipe kolom `OS_AKSEPTASI_KLAIM_LIFE` sudah ada di `TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md` *(62 kolom; di dalam 55: `STS_REJECT NUMBER(38,0)`, `CLAIM_RETRO NUMBER`, `WPC DATE` — ketiganya ditebak teks)* | `barislamakolom.go`, `skemauji.go`, test baru | `kolomAngkaLama` += `STS_REJECT`, `CLAIM_RETRO`; `kolomTanggalLama` += `WPC`; `TipeKolomBarisLama` mengikuti dokumen **termasuk panjang `VARCHAR2`**; tabel tiruan **62** kolom *(tujuh yang tidak ditulis rule ikut dibuat, supaya bentuk tiruan = bentuk sungguhan)*; **test murni** yang membaca tabel markdown dokumen itu dan membandingkan ketiganya kolom demi kolom *(pola `strukturkolom_test.go`)*; `AmbilBarisLama` tetap 55 kolom. ⚠️ `Simpan` menulis `CLAIM_RETRO` teks ke kolom `NUMBER`: `BarisLamaDari` harus menjamin isinya angka atau NULL, dan test murni mengunci itu. **Tidak** menyentuh `002`/`003` — itu **w** |
| 2 | ⭐ **Pra-terbang bentuk** *(§2 x; dikerjakan bila `[DIPUTUSKAN]`, dan **disiapkan** bagian murninya bila belum)* | `migrasi.go`, `migrasi_test.go`, `pohonklaim_db_test.go` | pemecah kolom DDL diangkat dari test (`kolomMenurutDDL`) ke kode non-test; `bentukSama(ctx, tabel, kolomDDL)` membaca `SYS.ALL_TAB_COLUMNS` *(berawalan skema, ADR-U-0033)*; `JalankanMigrasi` memanggilnya untuk tiap `CREATE TABLE` yang `objekAda` **sebelum** loop eksekusi, menolak dengan pesan berisi selisih; test murni untuk pembanding *(sama · kolom kurang · kolom lebih · beda huruf besar-kecil)*; test db: buat `DOCUMENT_CLAIM` berbentuk lain di skema uji → `-migrate` **gagal** menyebut `DOCUMENT_CLAIM` dan kolomnya, `T_MIGRASI` tetap kosong; **`ObjekSudahAda` hanya berisi tabel yang bentuknya sama** |
| 3 | ⚠️ **Tiket 02 menulis `[dugaan]` bahwa lock `FOR UPDATE` ada di dalam procedure** *(ralat ronde 4)*. Sudah **terlihat**: `SUMBER-PENOMORAN-DBA.md` baris 56–63 procedure memuat `SELECT … FOR UPDATE` atas `(CLASS, JENIS, TAHUN)`. Penanda harus naik ke `[terverifikasi]` | tiket 02 bab `## Implementasi` | satu baris ralat: *"dugaan ini terverifikasi 26-09 lewat `SUMBER-PENOMORAN-DBA.md`; lock memang di dalam procedure, dan procedure tidak `COMMIT`"* |
| 4 | ℹ️ `PANDUAN-MENJALANKAN.txt` bab 5 masih menyebut celah proxy | `PANDUAN-MENJALANKAN.txt` | **sudah diralat** oleh sesi verifikasi sore ini *(belum di-commit)*; executor hanya memastikan isinya cocok dengan `vite.config.ts` |

Yang **sudah benar** dan tidak perlu disentuh: FK d tanpa `ON DELETE` + fixture komite; pagar dua
syarat dan `BolehDilewati`; `rapatkanSpasi`/`klausaConstraint`; label p1 di tiga `.sql` dan
`presisiSah`; blok ralat STRUKTUR *(termasuk paragraf `[dugaan]` j — dibiarkan)*; `006` diganti
nama; AC 12 dan 50 dicentang dengan alasannya; daftar terbuka 13 = kotak `[ ]`; tolakan atas
tuduhan "rename memutus ledger" benar *(katalog: `T_MIGRASI` tidak ada di mana pun)*.

---

## 4. URUTAN SESI

**Langkah 0 — titik tetap.** Commit berkas yang tertunda **dengan nama**, dua commit:
`docs: brief ronde 5, katalog DEV, panduan pemuat env` *(brief ini, dua dokumen `[data DBA]`,
`PANDUAN-MENJALANKAN.txt`, brief ronde 4)* dan `tooling: proxy dev Vite + pemuat .env`
*(`frontend/vite.config.ts`, `frontend/.env.example`, `muat-env.ps1`, `muat-env.cmd`)*. ⛔ Bukan
`git add -A`: `APP_RNM/.env` dan `frontend/.env` memang diabaikan git, tetapi berkas lain yang tidak
disebut di sini tidak boleh ikut. Lalu `git status --porcelain` kosong, `git rev-parse HEAD`
dicatat, uji tanpa Oracle hijau.

**Langkah B — G0.** Urutan: §3-3 → §3-1 → §3-2 *(bagian murni selalu; bagian `JalankanMigrasi`
hanya bila x `[DIPUTUSKAN]`)*. Tiap nomor: test dulu, diuji **gagal** pada kasus buruknya, lalu
kode. `go test ./...` harus bertambah, bukan berkurang.

**Langkah C — G2 bila terbuka.** **v1** → ganti nama di `007` (+`_down`), `pohonklaim.go`,
`models`, test, STRUKTUR *(blok ralat bertanggal, seperti m)*; `tabelDikecualikan` di
`strukturkolom_test.go` ikut. **w** → `003` `WPC DATE`, `002` `CLAIM_RETRO` sesuai arti yang
diputuskan, `BarisLama`/`Simpan`/`BarisLamaDari` mengikuti, STRUKTUR ralat. Keduanya **hanya**
selama `T_MIGRASI` belum ada di instance mana pun *(§2 l — masih benar)*.

**Langkah A — G1 bila terbuka.** Persis brief ronde 4 §4 Langkah A, dengan tambahan pertama:
`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER(:1)` harus **0** sebelum apa pun; bila
tidak nol, berhenti dan tanya. Lalu test db → `-migrate` dua kali → cacah objek *(8 · 5 · 15)* →
tiket 01.

**Langkah D — G3 bila terbuka.** Persis brief ronde 4 §4 Langkah D; `SUMBER-PENOMORAN-DBA.md`
sudah ada, jadi D-2 menjadi *"tandai dokumen itu dikonfirmasi DBA"* bila konfirmasinya ada.

**Verifikasi penuh sekali di akhir**, **penutup** *(bab `## Implementasi — ronde 5`, `/code-review`,
angka sesudah review, commit `claim-life: tiket 14 ronde 5 — <isi sebenarnya>`)*, dan **tanda
`resolved`** — semuanya persis brief ronde 4 §4. Tanpa G1, tiket 14 berakhir `claimed` lagi dan
kalimat itu ditulis di **awal** laporan.

---

## 5. YANG DIMINTA KE DBA — diperbarui

| # | Objek | Keadaan |
| ---: | --- | --- |
| 1 | **User kosong** di instance pengembangan, hak `CREATE TABLE / SEQUENCE / INDEX` + kuota tablespace, DSN | ⛔ belum — membuka G1 |
| 2 | Konfirmasi dua dokumen `[data DBA — dibaca sendiri]` sama dengan **produksi** | belum |
| 3 | `COMPATIBLE` pengembangan **dan** produksi | belum — menentukan tinjauan j |
| 4–6 | sumber procedure, `GENERATE_SEQUENCE_NUMBER`, `KODE_PRODUKSI`, `TANGGAL_CLOSING` | ✅ isinya ada *(SUMBER-PENOMORAN-DBA.md)*; tinggal konfirmasi |
| 7 | `M_LIFE_PREMIUM_DETAIL` DDL + OQ-001 sisa; **di mana `T_PREMIUM_LIST`** *(tidak ada di `POOLDATA`)* | belum — tiket 02 AC 20–22, bukan ronde ini |

---

## 6. GAYA KODE — tetap mengikat

Sama dengan brief ronde 4 §6. Test murni yang membaca dokumen markdown *(§3-1)* mengikuti pola
`strukturkolom_test.go`: gagal terang bila dokumennya hilang, tolak pemetaan yang tidak terpakai.

---

## 7. YANG MEMERLUKAN PERSETUJUAN MANUSIA — berhenti dan tanya

Brief ronde 4 §7, ditambah: `-migrate` ke skema **apa pun** yang `ALL_OBJECTS`-nya tidak nol ·
membaca **isi** tabel warisan mana pun *(katalog boleh, baris tidak)* · mengubah `007` di luar
keputusan v.

---

## 8. TELEMETRI EKSEKUSI — bab wajib di laporan akhir

Sama dengan brief ronde 4 §8, ditambah satu baris: **kolom tabel tiruan warisan** *(harus 62)* dan
**selisih peta tipe lawan dokumen katalog** *(harus 0)*, dari test §3-1.

---

*Disusun 26 September 2026 sore dari verifikasi independen commit `ec843e2`: 83 test dijalankan
ulang, diff 20 berkas dibaca utuh, klaim korpus ≤ 23 byte dihitung ulang, katalog instance
pengembangan dibaca siang harinya.*
