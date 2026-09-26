# PROMPT — sesi implementasi: Claim Life **lima tiket sekaligus** — 02, 03, 06, 04, 05

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya **kecuali §6 kalimat "satu
> sesi, satu tiket"**, yang diganti bab 1 berkas ini atas permintaan work owner 26 September 2026.
> Brief tiket 14 ronde 5–7 tetap menjadi rujukan skema, pagar, dan katalog.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> kelima tiket **utuh** *(`.scratch\claim-life\issues\02…06`)* → `STRUKTUR-TABEL-CLAIM-LIFE.md`
> **seluruhnya** → `KATALOG-TABEL-PESERTA-DAN-TREATY.md`, `TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md`,
> `SUMBER-PENOMORAN-DBA.md` → `CONTEXT.md` bab BusinessCode dan peran → §2–§4 berkas ini.
>
> **SESI INI:** claim-life · tiket **02, 03, 06, 04, 05** berurutan · tiket 14 dan 01 **tidak
> disentuh** kecuali dinyatakan di §4.

---

## 0. KEADAAN AWAL — 26 September 2026 malam

| | Keadaan |
| --- | --- |
| `HEAD` | `f754ad1` *(brief ini, brief ronde 7, dan katalog sudah ter-commit)*; tiket 14 *claimed* 40/53; tiket 01 *claimed* 4/7; tiket 02–06 *ready-for-agent*, nol AC dicentang. Working tree **bersih** — satu-satunya yang belum ter-commit adalah suntingan terakhir berkas brief ini |
| Uji | **99** test Go PASS; **20** test bertag `db` SKIP dengan pesan; 5 test JS; `vite build` 87 modul |
| Yang sudah ada dan **dipakai**, bukan ditulis ulang | skema 7 tabel + `T_MIGRASI` *(migrasi tertanam, pra-terbang bentuk, `-migrate`/`-migrate-down`)*; `repository.PohonKlaim.Simpan` *(satu transaksi: 6 tingkat + baris datar warisan + pagar `PeriksaNilaiWarisan`)*; `BongkarBarisLama`/`BarisLamaDari`; `AmbilHeader`, `AmbilPeserta`, `AmbilBaris`, `AmbilSpreading`; `models.Money`/`Ratio`/`Klaim`/`Peserta`/`BarisAdjustment`/`PohonKlaim`, `StatusBarisDariKode`; `services.Pelaku`, `WajibPeran`, `DalamTransaksi`; rute `GET /healthz`, `GET /api/klaim-life/{id}`; frontend satu halaman + proxy dev |
| Oracle | **G1 masih tertutup** — tidak ada user kosong. Seluruh test db ditulis dan **SKIP dengan pesan**; **nol** `-migrate` |
| Kredensial pengembangan | ada di `APP_RNM\.env` milik work owner; menunjuk `POOLDATA` — **hanya** untuk `go run ./cmd/api` *(baca)*, tidak pernah untuk test db atau migrasi |

---

## 1. LIMA TIKET DALAM SATU SESI — aturannya

1. **Urutan tetap: 02 → 03 → 06 → 04 → 05.** Itu urutan ketergantungannya *(03 butuh klaim dari
   02; 06 adalah gerbang simpan di 03; 04 menurunkan status dari baris 03; 05 memakai mesin 04)*.
   Tidak boleh dibolak-balik.
2. **Tiap tiket tetap satu unit penuh** brief induk §6 butir 3–8: `Status: claimed` → test dulu →
   kode → bab `## Implementasi — <tanggal>` → `/code-review` atas **titik tetap tiket itu** *(SHA
   sesudah commit tiket sebelumnya)* → **satu commit per tiket** dengan pesan
   `claim-life: tiket NN — <judul>` + baris `Tiket:`. Nol commit gabungan.
3. **Blocker di dalam tiket tidak menghentikan sesi.** AC yang menunggu keputusan *(§2)* atau
   Oracle diberi tanda di bab Implementasi dan tiket tetap `claimed`; sesi **lanjut** ke tiket
   berikutnya. Yang dilarang brief induk §6-2 adalah mengerjakan tiket lain **diam-diam** —
   di sini lanjutnya **dinyatakan**, per tiket, dengan alasannya.
4. **Berhenti total** hanya bila: tiket sebelumnya gagal `go test ./...` dan tidak dapat
   dipulihkan; `/code-review` menemukan cacat AC yang tidak dapat ditutup di sesi itu; atau ada
   hal di §7. Bila berhenti di tengah, tiket yang belum disentuh tetap `ready-for-agent`.
5. **Verifikasi penuh** *(vet, vet db, gofmt, build, test, test db, typecheck, test JS, build JS)*
   dijalankan **sekali di akhir tiap tiket** sebelum commit-nya, bukan hanya di akhir sesi.
6. **Laporan akhir** memuat bab per tiket *(AC ditutup / total, SHA)* dan satu telemetri sesi.

---

## 2. KEPUTUSAN WORK OWNER YANG MENYENTUH KELIMA TIKET

| | Keputusan | Bukti / keadaan | Pemakai |
| ---: | --- | --- | --- |
| o1–o3 | **Penomoran klaim** *(brief ronde 4 §2, ronde 5)*. Tiket 02 AC 2, 3, 7–11 menyebut `PROC_GENERATE_SEQUENCE_NUMBER`; keputusan o melarang memanggil procedure. Sampai AC-nya ditulis ulang work owner, executor **memisahkan** penomoran di balik antarmuka `services.Penomor` *(§4 tiket 02)* dan **tidak** menulis logikanya. Bila `[DIPUTUSKAN]` sebelum sesi, executor menulis `PenomorCounter` dari `SUMBER-PENOMORAN-DBA.md` **termasuk perakitan format** `<prefix>K<kode>.MM.YYYY.<5 digit>` yang di Pega dirakit pemanggil | sumber procedure sudah terbaca — `[USULAN]` | 02 |
| aa | **Generator `CLM-xxxxxx` / `KMT-xxxxxx`** — tiket 14 AC 35 `[terbuka]`, tetapi tiket 02 **tidak dapat membuat satu baris `T_WORK_CLAIM` pun** tanpa itu. Usulan: sequence `SEQ_WORK_CLAIM` *(langkah migrasi baru `009`)* + `LPAD(6,'0')`, awalan menurut jenis baris, **tanpa reset tahunan**, dibaca lewat `nomorBerikut` yang sudah ada | fixture memakai `CLM-UJI001`; nol bukti korpus soal reset — `[USULAN]` | 02 |
| z | Mata uang header *(brief ronde 7)* — usulan **z1** kolom `CURRENCY` di `002` | agregat: 0 klaim campur — `[USULAN]` | 02 |
| ab | **Siapa pelaku permintaan HTTP.** `services.WajibPeran` sudah ada; sumber peran adalah satu tabel *(ADR-U-0030)* yang belum ada, dan tiket 07 yang menegakkannya. Tiket 05 butuh `ReasLifeAdmin` **sekarang**. Usulan: handler membaca pelaku dari header `X-Pelaku` *(akun)* dan `X-Peran` *(daftar peran)* **hanya bila** env `AUTH_STUB=true`, ditolak keras saat `IS_PEGA_PROD=true`, dan seluruh test wewenang berjalan di seam `services` dengan `Pelaku` langsung. Ini **bukan** autentikasi; ia penunda sampai tiket 07/IAM | `[USULAN]` — perlu persetujuan IAM/work owner | 05, 03 |
| ac | **Baris negatif jurnal balik** *(tiket 02 AC 28)*: kolom mana yang bertanda negatif belum ditulis di tiket. Executor membaca verdict V14 grilling Endorsement Life dan menulis kolomnya di tiket 02 bab Implementasi sebagai `[dugaan]` bila spec tidak tegas; penyaringnya tetap dipasang | spec Endorsement Life — `[USULAN]` | 02 |
| ad | **Kolom isi `T_CLAIMLF_DOCUMENT`** *(tiket 03 dokumen per peserta)*: diturunkan dari sensus `.DocumentList` di `SaveOutStandingLife_Act` — executor melakukan sensus itu dan menulis hasilnya sebagai langkah migrasi **`010`** *(tabel `007` hanya punya PK + FK)*, dengan bukti path XML per kolom | korpus — `[USULAN]` | 03 |
| — | `PREMIUM_SPREADED_NET` dua cabang *(tiket 03)*: **tidak** ditebak; kolom diisi `NULL` dan dilaporkan `[terbuka — Product+UW]` | korpus | 03 |
| — | `@addCalendar(...,0,0,0,1,0,0,0)` *(tiket 06)*: direplikasi apa adanya, satuannya `[dugaan]` **hari**; ditulis di test dan tiket | korpus | 06 |
| k, l | tetap | | semua |

Seluruh baris `[USULAN]` disahkan dengan mengganti kata itu menjadi `[DIPUTUSKAN]` **di berkas ini**
sebelum sesi. Yang masih `[USULAN]` saat sesi berjalan diperlakukan sebagaimana kolom "keadaan"
menyatakan — **tidak** ditebak.

---

## 3. FAKTA KATALOG YANG MENGUBAH CARA MENULIS KODE

| Fakta *(katalog pengembangan, agregat saja)* | Akibat |
| --- | --- |
| `M_LIFE_PREMIUM_DETAIL` **66,8 juta baris**, 85 kolom; index pada `PL_NUMBER`, `CERTIFICATE_NO`, `POLICY_NO`, **tidak** pada `EDMSTATUS` | pencarian peserta **wajib** berawalan `PL_NUMBER` *(atau `CERTIFICATE_NO`)* dan berbatas hasil *(`FETCH FIRST :n ROWS ONLY`)*; nol query tanpa penyaring ber-index; test statik menolak `SELECT` atas tabel ini tanpa `PL_NUMBER`/`CERTIFICATE_NO` di `WHERE` |
| `EDMSTATUS`: `NULL` 59,1 juta · `Batal` 7,7 juta · `Old` 661 · `New` 148 · `Delete` 142 · `''` **0** | baris NB = `NULL`; penyaring `EDMSTATUS IS NULL OR EDMSTATUS NOT IN ('Batal','Delete')`; test AC 26 memuat kasus `NULL` **dan** `''` |
| Keempat tanggal valuasi + `WPC` ada di `M_LIFE_PREMIUM_DETAIL` sebagai `DATE` | tiket 06 membacanya dari peserta yang **disalin** ke `T_CLAIMLF_PREMIUMLIST_DETAIL` saat dipilih *(kolomnya sudah ada di `003`)*, bukan query ulang ke tabel 66 juta baris |
| `RETROCESSIONLIFE` adalah **VIEW** ber-13 kolom `VARCHAR2(4000)` semua; `TREATYYEAR_LIFE` tabel 7 kolom; `RATE_LIFE` sebagian terbaca | angka spreading tiba sebagai **teks**: urai lewat `utils.ParseDecimal`, laporkan yang gagal; pemisah desimalnya belum diketahui — periksa dengan agregat sebelum menulis pembaca |
| `GetJsonProductLife` → `m_product_life.JSONDATA` | ⛔ **tidak ditiru** *(tiket 02 AC 38)*; nama produk dari view `PRODUCT_LIFE` |
| `STS_REJECT` di DEV bernilai 0/1/2/**4** | mesin status 04 memetakan `4` ke *tidak diketahui* dan **tidak** pernah menulis 4; artinya `[terbuka — tiket 04/work owner]` |

---

## 4. RENCANA PER TIKET — seam, pintu, dan yang dikunci test

Seam tetap tiga *(brief induk §5)*: `repository` ↔ skema uji *(db, SKIP tanpa Oracle)*, HTTP ↔
skema uji, `services` murni. Uang: `Money`/`Ratio` dari teks; nol float; nilai dibandingkan dengan
angka contoh spec bila ada.

### Tiket 02 — Register klaim + penomoran

| Bagian | Isi |
| --- | --- |
| Pintu | `POST /api/klaim-life` *(daftar: `PL_NUMBER`, penunjuk polis, `TYPE`, `BUSINESS_CODE`, daftar peserta terpilih)* → 201 `{id, nomorKlaim}`; `GET /api/peserta-life?pl=<PL_NUMBER>&n=<batas>` |
| `services` | `Pendaftaran.Daftar(ctx, pelaku, permintaan)`: satu transaksi *(`DalamTransaksi`)*: `T_WORK_CLAIM` *(ID dari **aa**, `LINI=LIFE`, `TYPE`, `CASE_ID`)* → `T_GENERAL_CLAIM` *(shared PK; `CLAIM_NO` dari `Penomor`)* → peserta terpilih disalin ke `T_CLAIMLF_PREMIUMLIST_DETAIL` *(termasuk keempat tanggal valuasi, `WPC`, `IS_CHECK`)* → baris datar warisan lewat `BarisLamaDari` + `PeriksaNilaiWarisan`. `Penomor` = antarmuka `NomorBerikut(ctx, tx, jenis, tanggal) (nomor string, err)`; implementasi **sementara** `PenomorBelumDiputuskan` mengembalikan galat terang *(AC 2, 3, 7–11 `[terbuka — o]`)*; test memakai `PenomorPalsu` |
| `repository` | `PesertaPolis.Cari(ctx, plNumber, batas)` atas `M_LIFE_PREMIUM_DETAIL` dengan penyaring §3 dan **ac**; `Pendaftaran.Simpan` **memakai** `PohonKlaim.Simpan` yang ada — tidak menulis INSERT baru untuk tabel yang sudah punya penulis |
| Test murni | penyaring `EDMSTATUS` *(NULL, '', Old, New, Batal, Delete)*; baris negatif; `Penomor` dipanggil **sekali** per pendaftaran dan **sesudah** validasi; kegagalan `Penomor` membatalkan seluruh transaksi |
| Test db *(SKIP)* | pendaftaran menulis 3 tempat + baris datar dalam satu transaksi; dua pendaftaran = dua `ID` berbeda; `T_WORK_CLAIM.LINI = 'LIFE'`; nol tulisan ke tabel polis/marketing *(negatif: test yang menemukan tabel itu **gagal**)* |
| Frontend | halaman `Register`: cari peserta per `PL_NUMBER`, pilih, kirim; menampilkan nomor **atau** pesan "nomor belum dapat dibentuk" *(o belum diputuskan)* — uang tetap teks |
| Tetap terbuka | AC 15, 16 *(kosong artinya)*, 17 *(pemilik PremiumList Life)*; AC 2, 3, 7–11 sampai **o** |

### Tiket 03 — Baris `AdjustmentList` + Save ke Outstanding

| Bagian | Isi |
| --- | --- |
| Pintu | `POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment` *(tambah baris)*; `POST /api/klaim-life/{id}/simpan-outstanding` *(gerbang dokumen + DOL tiket 06, lalu status `0`)*; `GET /api/klaim-life/{id}` yang ada sudah menampilkan baris |
| `services` | `Adjustment.Tambah`: baris ke-2 dst **mewarisi 8 kolom** dari baris pertama peserta itu, **tidak** mewarisi status; `STS_REJECT` diisi `0` **oleh aksi Save ke Outstanding**, bukan saat insert *(test menolak hardcode: baris yang belum disimpan ke Outstanding **tidak** punya status)*; `ACCEPTATION_DATE` **tidak** disentuh; bank boleh kosong; `Spreading.Hitung(baris, treaty)` murni: `AMOUNT = RetrocadedShare × PERCENT_SHARE ÷ 100`, `PREMIUM_SPREADED_GROSS = RATE × (1 + EM_PERCENT) × AMOUNT`, `PREMIUM_SPREADED_NET = NULL` `[terbuka]`; tanpa retro → nol baris spreading, **bukan** galat |
| `repository` | pembaca treaty *(`RETROCESSIONLIFE`, `TREATYYEAR_LIFE`, `RATE_LIFE`)* — teks → `ParseDecimal`, galat dilaporkan; penulisan lewat `PohonKlaim.Simpan`; dokumen: langkah migrasi **`010`** kolom isi *(**ad**)* + `Dokumen.Daftar/Simpan` |
| Test murni | pewarisan 8 kolom; rumus spreading dibandingkan **angka contoh spec/korpus** bila ada, kalau tidak dengan angka bulat yang ditulis di test beserta turunannya; gerbang dokumen menyebut **peserta mana**; satu transaksi *(kegagalan spreading membatalkan baris)* |
| Test db *(SKIP)* | dua peserta × dua putaran = empat baris tertelusur; spreading dibekukan *(ubah treaty sesudah simpan → angka tetap)*; dokumen `SELECT` biasa |
| Frontend | daftar baris per peserta dengan status **kata**; formulir baris + tiga field bank; tombol Save ke Outstanding dengan pesan gerbang |

### Tiket 06 — Validasi DOL per `Type` + `ContentNote`

| Bagian | Isi |
| --- | --- |
| `models` | `BusinessCode` → `ContentNote` sebagai **tabel data** dari `CONTEXT.md` L1–L21 *(satu `map`/slice, nol `if` bercabang)*; nilai `ContentNote` tertutup: `DEATH`, `HEALTH`, `CI`, `TPD`, `TI` |
| `services` | `ValidasiDOL(tipe, dol, peserta)`: `QP`/`QR` → jendela `GROSS_VALUATION_*`; `TP`/`TR` → `RETRO_VALUATION_*` **dengan pergeseran** `+1` satuan `[dugaan: hari]`; gagal → galat setara `"Invalid DOL"` yang menyebut peserta; dipanggil dari Save ke Outstanding tiket 03; **satu** field `Type` |
| Test murni | tiap `Type` × *(di dalam, tepat di batas, di luar, tanggal kosong)*; pemetaan 21 kode; kode di luar daftar → galat, bukan `DEATH` diam-diam |
| Tetap terbuka | satuan `@addCalendar` *(Product+UW)*; arti `QP`/`QR` *(OQ-020)* |

### Tiket 04 — Mesin status per baris

| Bagian | Isi |
| --- | --- |
| `models` | `StatusBaris` sudah ada; tambah `Klaim.StatusTurunan()`: *selesai* ⇔ nol Outstanding **dan** ≥1 Aksep; *ditolak seluruhnya* ⇔ nol Outstanding dan nol Aksep; selain itu *berjalan*. **Tidak** disimpan sebagai kolom |
| `services` | `Transisi(baris, ke)`: dari Outstanding ke Aksep/Ditolak saja; dari Aksep/Ditolak ke apa pun → galat *(kefinalan, ADR-U-0011)*; kode `4` → tidak diketahui, tidak pernah ditulis; pencerminan `STS_REJECT`/`ACCEPTED_NO` ke `T_CLAIMLF_PREMIUMLIST_DETAIL` dan `T_GENERAL_CLAIM` **dalam transaksi yang sama** dengan penulisan baris *(satu fungsi repository `PohonKlaim.PerbaruiStatusBaris`)* |
| Pintu | `GET /api/klaim-life/{id}` menambah `statusKlaim` *(kata)* dan tiap baris membawa `status` *(sudah ada)*; riwayat = seluruh baris, bukan yang terakhir |
| Test murni | tabel transisi lengkap *(termasuk yang **ditolak**)*; aturan turunan pada 6 kombinasi; menolak satu baris tidak menutup klaim |
| Test db *(SKIP)* | pencerminan dua tingkat terjadi bersamaan; baris final tidak berubah lewat jalur mana pun |

### Tiket 05 — Reject Outstanding oleh Admin

| Bagian | Isi |
| --- | --- |
| Pintu | `POST /api/klaim-life/{id}/adjustment/{adjId}/tolak` → 200 baris terbaru; 403 bila bukan `ReasLifeAdmin`; 409 bila baris bukan Outstanding; 422 bila klaim belum bernomor |
| `services` | `Tolak(ctx, pelaku, klaimID, adjID)`: `WajibPeran(pelaku, "ReasLifeAdmin")` → klaim punya `CLAIM_NO` → baris Outstanding → `Transisi(baris, Ditolak)` *(tiket 04)*; nilai `2` **menurut aksi**, sama dengan jalur Komite; pencerminan bersamaan |
| Pelaku | sesuai **ab**; tanpa `AUTH_STUB=true` pintu HTTP menjawab 401 dan test HTTP-nya SKIP dengan pesan, test `services` tetap jalan |
| Test murni | jalur **ditolak** diuji: peran salah, baris Aksep, baris sudah Ditolak, klaim tanpa nomor; jalur berhasil: hanya baris itu berubah, klaim tetap menerima baris baru |

---

## 5. URUTAN SESI

**Langkah 0.** `git add PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-02-06-BATCH.md` lalu commit
`docs: ralat Langkah 0 dan catatan penjaga, brief batch tiket 02–06` *(hanya berkas ini; tiga dokumen
lain sudah ter-commit di `f754ad1`)*. `git status --porcelain` kosong; SHA dicatat sebagai titik
tetap **tiket 02**; uji tanpa Oracle hijau *(99 test, 20 SKIP db, 5 JS)*.

**Tiga penjaga yang PASTI tersentuh sesi ini — perbarui angkanya, jangan dilonggarkan:**

| Penjaga | Sekarang | Sesudah | Sebab |
| --- | ---: | ---: | --- |
| `TestSeluruhCreateDapatDibacaNamanya` *(cacah `CREATE`)* | 19 | **20** | `009` menambah `SEQ_WORK_CLAIM` *(aa)* |
| `TestKolomDDLCocokDenganStruktur` | — | STRUKTUR ikut diralat | `002` `CURRENCY` *(z1)*; blok ralat bertanggal di bab `T_GENERAL_CLAIM`, satu baris kolom baru di tabelnya |
| Cacah objek `-migrate` *(Langkah A, bila G1)* | 8 · 5 · 15 | **8 · 6 · 15** | sequence baru; `010` hanya `ALTER TABLE … ADD`, nol objek baru |

**Sensus dokumen (ad) — hasil awal, dilanjutkan executor:** kelas Pega `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`
memetakan ke tabel warisan `DOCUMENT_CLAIM` *(14 kolom di brief ronde 4 §9: `NAMAFILE`, `MIME`,
`KATEGORI_1/2`, `TANGGAL`, `T_STORAGE_ID`, `NOAKSEP`, `NOPREKAS`, `PAYMENTDATE`, `INSKEY_*`,
`IDPEGA`, `PXCREATEOPERATOR`)*; `InsertDocument_Act` merujuk `InsertGoogleStorage_Act` dan rule
`Insert_T_Storage_SQL` / `GetLinkStorage_SQL` *(penyimpanan berkas di Google Storage, ADR-U-0010;
`T_STORAGE_ID` adalah penunjuknya)*. Yang belum dibaca: `Claim Life/Section/DocumentLife.xml`
*(properti yang diisi manusia)* dan aturan **"lengkap"** di `SaveOutStandingLife_Act` *(kategori mana
yang wajib per peserta)* — keduanya dibaca executor sebelum menulis `010`, dengan path XML per kolom.

**Per tiket, berurutan 02 → 03 → 06 → 04 → 05:** `Status: claimed` → test dulu pada seam §4 →
kode → verifikasi penuh → bab Implementasi *(AC ditutup / total; AC terbuka dengan pemiliknya;
keputusan §2 yang dipakai atau ditunggu)* → `/code-review` atas titik tetap tiket itu → perbaiki
temuan AC → commit → SHA baru = titik tetap tiket berikutnya.

**Migrasi baru** *(`009` sequence work claim bila **aa**; `010` kolom dokumen bila **ad**; `002`
`CURRENCY` bila **z1**)*: berkas lanjutan **bernomor baru**, kecuali `002` yang masih boleh
disunting langsung *(§2 l — `T_MIGRASI` belum ada di mana pun; nyatakan itu di tiket)*.
`TestKolomDDLCocokDenganStruktur` menuntut STRUKTUR ikut diralat *(blok bertanggal)*.

**Tiket 14 dan 01** tidak disentuh, **kecuali** satu baris di tiket 14 bab Implementasi ronde 7
*(ditulis sesi ini)*: langkah migrasi mana yang ditambah sesi ini dan kenapa.

**Laporan akhir sesi**: satu bab per tiket + telemetri §8.

---

## 6. GAYA KODE — tetap mengikat

Brief ronde 5 §6. Tambahan untuk batch: tiap paket/berkas baru menyebut **tiket** pemiliknya di
komentar kepala; nama fungsi layanan mengikuti kata kerja tiket *(`Daftar`, `Tambah`,
`SimpanOutstanding`, `ValidasiDOL`, `Transisi`, `Tolak`)*; frontend tetap `.tsx`, uang teks.

---

## 7. YANG MEMERLUKAN PERSETUJUAN MANUSIA — berhenti dan tanya

Brief ronde 5 §7, ditambah: **menyimpan baris apa pun** dari `M_LIFE_PREMIUM_DETAIL` ke artefak
*(termasuk fixture — buat sendiri `UJI-*`; kolom `KTP` **tidak pernah** masuk fixture)* · membaca
`m_product_life.JSONDATA` · menulis logika penomoran tanpa **o** `[DIPUTUSKAN]` · membaca peran dari
mana pun selain **ab** · menambah langkah migrasi di luar `009`/`010`/`002` · menyentuh tiket 07–15.

---

## 8. TELEMETRI EKSEKUSI — bab wajib di laporan akhir

| Besaran | Cara ukur |
| --- | --- |
| Per tiket: AC ditutup / total, AC terbuka dengan pemilik, SHA titik tetap dan commit | tiket + `git log` |
| Per tiket: test murni · test db *(SKIP/PASS)* · test JS ditambah | `go test -v` sebelum/sesudah |
| Keputusan §2 yang `[DIPUTUSKAN]` vs ditunggu, dan AC mana yang terkena | tulis di baris pertama tiap bab |
| Query ke `M_LIFE_PREMIUM_DETAIL`: seluruhnya ber-`PL_NUMBER`/`CERTIFICATE_NO` dan berbatas | test statik |
| Berkas dibuat / diubah per tiket, baris berisi | `git diff --stat <titik tetap tiket>..<commit tiket>` |
| Sub-agen review per tiket: token, panggilan | laporan harness |
| Token sesi utama · biaya · jam dinding | **⛔ tidak diukur** — nyatakan |

---

## 9. SESUDAH TIKET 02 — verifikasi independen `753cef2`, dan yang harus terjadi sebelum tiket 03

**Yang tereproduksi:** 26 berkas +1.598/−27; 23 test db SKIP dengan pesan; 88 modul; 7/26 AC
*(5, 13, 20, 21, 22, 24, 25)*; keputusan aa, z1, ab diterapkan; sepuluh perbaikan tinjauan ada di
kode *(pagar pelaku, `\b`, urutan nomor→pengenal, `TRIM` di SQL + `PesertaHidup`, `CURRENCY` ditulis
dan dibaca, stub sebagai argumen, 501 tanpa nama objek, kunci polis+sertifikat)*.

⛔ **Yang TIDAK benar: `753cef2` memuat satu test MERAH.** `go test ./...` = **110 PASS, 1 FAIL**,
bukan "111, nol FAIL" seperti ditulis log dan bab Implementasi. `TestSetiapPemanggilBukaMemeriksaBolehDilewati`
mengunci cacah pemanggil `skemauji.Buka()` = 7, sedangkan `services/pendaftaran_db_test.go` — yang
lahir dari perbaikan tinjauan nomor 7 — menjadi pemanggil **ke-8**. Verifikasi penuh dijalankan
**sebelum** perbaikan tinjauan, lalu tidak diulang. Aturan §1-5 karena itu diperjelas: **verifikasi
penuh dijalankan lagi sesudah perbaikan `/code-review`, tepat sebelum commit.**

**Langkah 02-L — commit lanjutan tiket 02, sebelum apa pun yang lain**
*(pesan `claim-life: tiket 02 lanjutan — penjaga hijau, peserta disalin lengkap`)*:

| # | Yang dikerjakan |
| ---: | --- |
| 1 | `TestSetiapPemanggilBukaMemeriksaBolehDilewati`: `mau` 7 → **8**, komentar menyebut `pendaftaran_db_test.go`; ralat satu baris di bab Implementasi tiket 02: *"110 PASS + 1 FAIL saat commit, diperbaiki di lanjutan"* — bukan menghapus klaim lama |
| 2 | **Peserta disalin lengkap** *(§4 tiket 02, belum dikerjakan; prasyarat tiket 06)*. `models.Peserta` bertambah: `SumberID` *(kolom `ID` M_LIFE_PREMIUM_DETAIL → `SOURCE_ID`)*, `IsCheck`, `TanggalValuasiGrossMulai/Selesai`, `TanggalValuasiRetroMulai/Selesai`, `WPC`, `TanggalMulai/Efektif/Berakhir/Lapse`, `STNC`, dan uang polis `SumInsured`, `SumReasured`, `GrossPremium`, `NetPremium`, `CedingRetention`, `ShareNusantaraRe`, `ShareRetro`, `RetrocededShare` *(Money)*, `EMPercent` *(Ratio)*. ⛔ **Nol nama orang di model**: `NAME_OF_INSURED`/`POLICY_HOLDER` tetap hanya di `CalonPeserta` untuk layar, tidak disalin ke `T_CLAIMLF_PREMIUMLIST_DETAIL` *(kolomnya ada di DDL, tetap NULL sampai keputusan work owner — catat)* |
| 3 | **Server membaca ulang peserta terpilih sendiri** lewat `(PL_NUMBER, CERTIFICATE_NO)` *(keduanya ber-index)* di dalam `Daftar`, lalu menyalinnya — klien hanya mengirim nomor sertifikat. Nilai polis **tidak** dipercaya dari badan HTTP. Query-nya ikut aturan penjaga: ber-index dan berbatas |
| 4 | `PohonKlaim.Simpan` menulis kolom peserta baru; `AmbilPeserta` membacanya *(TO_CHAR untuk angka dan tanggal)*; test murni pulang-pergi model ↔ `BarisLama` bila tersentuh; test db *(SKIP)*: keempat tanggal valuasi terbaca kembali sama persis |
| 5 | `CalonPeserta` diberi tag JSON camelCase *(`nomorSertifikat`, …)* seperti kontrak lain; `api.ts` mengikuti |
| 6 | Daftar "AC terbuka" tiket 02 berjumlah **19**, bukan 18 *(7 + 1 + 6 + 2 + 2 + 1)*; angkanya diralat |

**Keputusan baru untuk work owner — `[USULAN]`:**

| | Keputusan | Keadaan |
| ---: | --- | --- |
| ae | **`CASEID` klaim baru.** `Daftar` membiarkannya kosong, sehingga baris datar warisan klaim baru **tidak dapat dikelompokkan** oleh hilir yang membaca `OS_AKSEPTASI_KLAIM_LIFE` per `CASEID` *(Arasapas)*, dan `Hapus`/`CacahBarisLama` memakai sumbu itu. Usulan **ae1**: `CASEID` = pengenal work object *(`CLM-xxxxxx`)* — di Pega pun `CASEID` adalah pengenal work object-nya, jadi ini paritas, bukan satu nilai dua arti; dicatat sebagai penyimpangan sadar di tiket 02. **ae2**: tetap kosong sampai work owner menetapkan sumber lain | `[USULAN]` — rekomendasi **ae1** |

**Sesudah 02-L hijau dan ter-commit:** lanjut **03 → 06 → 04 → 05** persis §4–§5, tiap tiket satu
unit. Butir **o** masih menahan AC 2, 3, 7–11 tiket 02; tidak ada tiket berikutnya yang bergantung
padanya, jadi ia tidak menghalangi rantai.

---

*Disusun 26 September 2026 malam atas permintaan work owner "proses 5 tiket sekaligus": kelima
tiket dibaca utuh, rantai ketergantungan 02→03→06→04→05 diturunkan dari kolom "Blocked by", katalog
tabel peserta dan treaty dibaca dari instance pengembangan (agregat saja).*
