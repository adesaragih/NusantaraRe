# Tiket - Claim Life

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Claim Life | **15** | 15 | 0 | 0 | 0 | 0 |
| **Jumlah** | **15** | **15** | **0** | **0** | **0** | **0** |

---

# Claim Life

Jumlah tiket: **15**

## Claim Life - 01 - Kerangka aplikasi + seam API + tracer "buka satu klaim Life"

**Status:** ready-for-agent

**Blocked by:** None (can start immediately)

#### Hasil & nilai pengguna

Seorang pengguna dapat membuka satu klaim Life dan melihat baris-baris `AdjustmentList`-nya di
layar. Itu perilaku yang tipis, tetapi ia menembus **seluruh lapisan** — React → HTTP → handler →
service → repository → Oracle — sehingga sekaligus **menegakkan seam** yang dipakai dua belas tiket
sesudahnya.

Nilainya bukan fiturnya, melainkan **jalurnya**: setelah tiket ini selesai, setiap tiket berikutnya
menambah perilaku di jalur yang sudah terbukti hidup, bukan membangun jalur baru.

#### Area codebase

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

#### Rule Pega sumber

Tiket ini **tidak meniru perilaku bisnis** — ia hanya membaca. Bentuk data diambil dari kolom yang
sudah terbaca:

| Rule | Identitas | Dipakai untuk |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` 55 nama kolom `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` terbaca langsung dari blok PL/SQL-nya |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` bentuk baris `AdjustmentList` dan kolom yang ditampilkan |

#### ADR terkait

**ADR-U-0003** (uang non-float — presisi desimal eksak; kolom Oracle `NUMBER` tanpa presisi),
**ADR-U-0013** (koneksi database dari env var; **alamat layanan keluar TIDAK** — ia di-lookup runtime
dari `M_LINK_SERVICE`), **ADR-U-0011** (agregat memuat **koleksi baris**, bukan satu status).

#### Acceptance criteria

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

#### Skema uji

`[terbuka]` **OQ-001** (pemilik **DBA**) — tidak ada DDL di korpus, sehingga **tipe dan presisi**
kolom tidak diketahui. Untuk tiket ini dipakai **skema uji provisional**: nama kolom dari rule di
atas, tipe numerik berpresisi longgar, tipe teks berpanjang longgar.

Ini **sah untuk skema uji** dan **tidak sah untuk DDL produksi maupun migrasi** — keduanya tetap
menunggu OQ-001 dan ditangani di tiket **13**.

#### Perintah verifikasi

Toolchain belum ada; tiket ini yang membuatnya. Setelah selesai, ketiga perintah berikut **harus
berjalan hijau** dan dicatat di `Makefile`:

```
go build ./...
go test ./internal/...
cd frontend && npm test
```

Tambahkan target gabungan `make check` yang menjalankan ketiganya.

## Claim Life - 02 - Register klaim Life + penomoran

**Status:** ready-for-agent

**Blocked by:** 01 (kerangka aplikasi + seam API), **14 (skema relasional klaim — PREFACTOR)**

> ⚠️ **Diselaraskan 2026-09-16, DIRALAT 2026-09-18 — revisi penyimpanan.** Klaim baru ditulis ke
> **`T_GENERAL_CLAIM`** (spec §2b) — **dan tetap** di-`INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE`.
> Yang **dibuang hanya JSON**. Lihat blok AC "Penyimpanan relasional".
>
> ⚠️ **RALAT 2026-09-18** `[keputusan work owner]` — **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING`
> DIHAPUS**, bukan diganti nama. Pendaftaran klaim **tidak lagi menulis keduanya**; data polis dan
> marketing **dibaca hidup** dari tabel polis (`T_PREMIUM_LIST` dkk, modul PremiumList Life) lewat
> penunjuk di `T_GENERAL_CLAIM`. ⛔ Jangan membuat `T_CLAIMLF_POLICY`/`T_CLAIMLF_MARKETING`.
> Pohon yang berlaku: `spec.md` §2b RALAT D.

#### Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat mendaftarkan klaim Life baru dan klaim itu langsung memperoleh
nomornya sendiri — sehingga saya tidak perlu menomori manual dan nomor tidak pernah bentrok antar
petugas. *(User story 1 dan 2 di spec)*

#### Area codebase

`internal/models` (entitas klaim), `internal/repository` (pemanggilan stored procedure + penulisan
rekam klaim), `internal/services` (orkestrasi pendaftaran), `internal/handlers` (endpoint pendaftaran),
`frontend/` (formulir Register).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Flow/Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` shape `Assignment2` "Input Register" — tahap pertama siklus |
| `Claim Life/RDBList/GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` blok PL/SQL memanggil `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` dengan **3 parameter masuk dan 2 keluar**, lalu `COMMIT` |

Tanda tangan yang terbaca `[terverifikasi]`:

```sql
BEGIN
  POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER(
    {ParamSeq.CARI1}, {ParamSeq.CARI2}, TO_DATE({ParamSeq.CARI3},'DD/MM/YYYY'),
    {ParamSeq.HASIL1 out}, {ParamSeq.HASIL2 out});
  COMMIT;
END;
```

#### ADR terkait

**ADR-U-0006** (penomoran lewat stored procedure — **jangan replikasi logikanya**),
**ADR-U-0003** (uang non-float). Batas transaksi: **OQ-013 tertutup** — Go memegang transaksi.

#### Acceptance criteria

- [ ] Klaim Life baru dapat didaftarkan lewat API dan muncul sebagai klaim berstatus awal.
- [ ] Nomor klaim **diperoleh dari `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`**, bukan dihitung di
      aplikasi.
- [ ] Aplikasi **tidak** memuat logika pembentukan format nomor apa pun. *(ADR-U-0006)*
- [ ] Dua pendaftaran berurutan menghasilkan dua nomor berbeda.
- [ ] Rule penomoran lama **tidak** dimigrasikan: tidak ada padanan `Generate_NoKlaim_Life` maupun
      `Generate_NoKlaim_LifeRetro` di kode. `[terverifikasi]` keduanya tidak terindeks sebagai
      rujukan aktif di `SaveOutStandingLife_Act`.
- [ ] Halaman React Register dapat mengirim pendaftaran dan menampilkan nomor yang diterima.
- [ ] Nomor yang dihasilkan berbentuk `<prefix>K<kode bisnis>.MM.YYYY.<5 digit>` — contoh
      `RNML-KL1.08.2026.00936`.
- [ ] Prefix diperoleh lewat **lookup** ke `POOLDATA.KODE_PRODUKSI` (`TYPE='LIFE'`), **tidak**
      ditanam sebagai konstanta di kode.
- [ ] Periode nomor mengikuti `POOLDATA.TANGGAL_CLOSING`, termasuk aturan cutover
      `TRUNC(now) <= 02/01/2026` → `12.2025`.
- [ ] Batas transaksi dipegang **Go**: commit terjadi segera setelah nomor terbentuk, sehingga lock
      `SELECT … FOR UPDATE` pada `GENERATE_SEQUENCE_NUMBER` tidak menahan pendaftar lain.
- [ ] Dua pendaftaran serentak tidak pernah memperoleh nomor yang sama.

##### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §2b

- [ ] ⚠️ Klaim baru ditulis ke **`T_GENERAL_CLAIM`**, **dan** di-`INSERT` flat ke
      `OS_AKSEPTASI_KLAIM_LIFE` — hilir masih membaca dari sana. Test yang **menolak** penulisan
      `OS_AKSEPTASI_KLAIM_LIFE` justru **gagal**. Keduanya dalam **satu transaksi**.
      **REVISI 2026-09-18:** ⛔ pendaftaran **tidak** menulis tabel polis maupun marketing — keduanya
      **dihapus**. Test yang menemukan penulisan ke tabel polis atau marketing milik klaim
      **gagal**. *(AC 32 spec — koreksi 2026-09-16; `[keputusan work owner]`)*
- [ ] ⚠️ **Tidak ada blob JSON** sebagai penyimpan isi klaim. *(AC 31 spec; penyimpangan sadar 1)*
- [ ] Header memuat keempat field `PremiumListSummary` — `CLAIM_NO`, `PL_NUMBER`, `RISLIPRNM`,
      `BUSINESS_NAME` — beserta `CASEID` dan `CLAIM_RETRO`. *(AC 36 spec)*
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS, MENUNGGU JAWABAN.**
      `T_CLAIM_POLICY` **dihapus**, dan kata **snapshot** **DICABUT**: data polis kini **dibaca
      hidup** ⚠️ **penyimpangan sadar** — polis yang berubah sesudah klaim dibuat **akan** mengubah
      tampilan klaim lama. Empat dari tujuh hal yang didaftar AC ini — ketiga **tanggal** dan
      **team group** — ⚠️ `[terbuka]` **kehilangan rumah**: bukan pindah, tetapi belum punya tempat
      (tiket 14 §Blocker). **Jangan tebak.** *(AC 37 spec)*
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS.** Atribut polis
      **tidak lagi disimpan di klaim sama sekali**, jadi "sekali per klaim" tidak punya yang
      dihitung. Yang menggantikan: ia **dibaca hidup** dari tabel polis lewat penunjuk di
      `T_GENERAL_CLAIM` ⚠️ **penyimpangan sadar**. *(AC 35 spec)*
- [ ] ⚠️ **REVISI 2026-09-18 — berpindah pemilik, tidak batal.** `PRODUCT_NAME` /
      `PRODUCT_NAME_ID` **bukan lagi kolom klaim**; keduanya termasuk 27 kolom yang dibaca dari
      `T_PREMIUM_LIST`. Larangan membaca `m_product_life.JSONDATA` kini mengikat **modul PremiumList
      Life**. Yang mengikat di sini: pendaftaran klaim **tidak** menyimpan nama/id produk.
      *(AC 38 spec)*
- [ ] ⚠️ Keadaan tangga klaim (posisi, status akseptasi, `Type`) ditulis ke **`T_WORK_CLAIM`**,
      **bukan** sebagai kolom `T_GENERAL_CLAIM`. *(AC 46 spec; penyimpangan sadar 6)*
- [ ] Pendaftaran klaim berjalan dalam **satu transaksi** — penulisan `T_GENERAL_CLAIM`, baris
      `T_WORK_CLAIM`-nya, dan `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE` **utuh atau tidak sama
      sekali**. **REVISI 2026-09-18:** frasa "ketiga tabelnya" **dicabut** — tabel polis dan
      marketing dihapus. *(AC 49 spec)*

##### Pencarian peserta — **hanya peserta hidup** `[keputusan work owner]`

⚠️ **Penajaman kontrak hilir (2026-09-15, verdict V14 grilling Endorsement Life).** Layar Register
memuat pencarian peserta; konteks Endorsement Life menulis **baris bernilai negatif** (jurnal balik)
ke tabel yang sama. **Peserta yang sudah dibatalkan atau di-soft-delete tidak boleh muncul.**

`[terverifikasi]` Rantai pencarian peserta di korpus:

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `InputRegisterClaimLife` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `INPUTREGISTERCLAIMLIFE` / `RULE-OBJ-HTML-SECTION` | `Claim Life/Section/InputRegisterClaimLife.xml` |
| `LoadDataPesertaSpesifik_Act` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `LOADDATAPESERTASPESIFIK_ACT` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/LoadDataPesertaSpesifik_Act.xml` (67.369 byte) — merujuk SQL di bawah pada baris 545 |
| `GetPesertaClaim_sql1` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/GetPesertaClaim_sql1.xml` |

`[terverifikasi]` Kuerinya **tidak memuat penyaring status apa pun** — ia mengambil seluruh baris
`M_LIFE_PREMIUM_DETAIL` untuk `PL_NUMBER` itu, termasuk baris negatif. Sensus korpus: **nol**
kemunculan `EDMSTATUS` di seluruh modul `Claim Life/`. ⚠️ Apakah Pega menyaring di lapisan lain
**tidak terbukti dari korpus** — perlakukan sebagai **kontrak yang wajib ditegakkan sistem baru**.

`[terverifikasi]` **Kolom penandanya terbaca dari korpus** — `M_LIFE_PREMIUM_DETAIL` punya tiga
kolom berakhiran status (daftar `INSERT` di `SaveMasterLPDet`, `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL`
/ `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL`):

| Kolom | Diisi dari | Penanda hidup/mati? |
| --- | --- | --- |
| **`EDMSTATUS`** | `TempValue.EDMStatus` → `Old` / `New` / `Delete` / `Batal` | ✅ **ya** |
| `STATUS` | `0` untuk `QR`/`QP`, `1` untuk `TP`/`TR` | ❌ penanda **jenis transaksi** |
| `STATUSOLD` | `1`/`0` | ❌ |

⚠️ `[terverifikasi]` **Jalur new business tidak mengisi `EDMSTATUS` sama sekali** — sensus
`InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT`): **nol**
kemunculan. Baris NB masuk dengan `EDMSTATUS` kosong/NULL.

- [ ] Pencarian peserta **tidak menampilkan** peserta ber-`EDMSTATUS` `'Batal'` maupun `'Delete'`.
      *(AC 25 spec)*
- [ ] Peserta **new business** (`EDMSTATUS` kosong/NULL) **tetap muncul**. ⚠️ Penyaring naif
      `EDMSTATUS NOT IN ('Delete','Batal')` membuang seluruh peserta NB di Oracle — test wajib
      memuat kasus ini dan **harus gagal** bila penyaringnya naif. *(AC 26 spec)*
- [ ] Peserta ber-`EDMSTATUS` `'Old'` dan `'New'` **tetap muncul**. *(AC 27 spec)*
- [ ] Baris bernilai **negatif** hasil jurnal balik endorsement **tidak pernah** sampai ke layar
      Register maupun ke perhitungan klaim. *(AC 28 spec)*
- [ ] Penyaringan terjadi di **satu tempat** — repository pembaca peserta. Test yang menemukan jalur
      baca peserta tanpa penyaring **gagal**. *(AC 29 spec)*
- [ ] `STATUS` dan `STATUSOLD` **tidak** dipakai sebagai penanda hidup/mati. *(AC 30 spec)*
- [ ] Jalur baca **akuntansi/ringkasan premium** — bila ada di konteks ini — **tidak** menyaring:
      ia melihat seluruh baris positif dan negatif. **Satu tabel, dua sudut pandang.**

**Blocker parsial:** `[terbuka]` tipe dan nullability `EDMSTATUS` belum terbaca — apakah `NULL` atau
string kosong untuk baris NB. Masuk **OQ-001 (sisa)**, pemilik **DBA**. **Tidak memblokir tiket** —
implementasi menangani **keduanya** (`IS NULL` *atau* `= ''`) sampai DBA memastikan.

#### Catatan penutupan (2026-09-14)

**OQ-002 TERTUTUP** untuk Claim — Life `[data DBA]`. Kontrak penomoran kini diketahui penuh:

```
<prefix> + "K" + <kode bisnis L1/L2/…> + "." + MM.YYYY + "." + LPAD(seq,5)
contoh:  RNML-KL1.08.2026.00936
```

| Aspek | Isi |
| --- | --- |
| Prefix `RNML-` | `[terverifikasi]` **di-lookup**, bukan hardcode — `Claim Life/RDBList/GetKodeProdLife_SQL.xml`: `SELECT KODE FROM POOLDATA.KODE_PRODUKSI WHERE TYPE='LIFE'` |
| Tanda tangan proc | `PROC_GENERATE_SEQUENCE_NUMBER(p_class, p_jenis, p_proddate, OUT p_bulan, OUT p_seq_number)` |
| Cakupan sequence | per **`(class, jenis, tahun)`** di tabel `GENERATE_SEQUENCE_NUMBER` |
| Penguncian | `SELECT … FOR UPDATE` |
| `p_jenis` | membedakan **retro / non-retro** |
| Periode | digulir lewat `POOLDATA.TANGGAL_CLOSING` (`TANGGAL VARCHAR2(10)` — **string**) |
| Aturan cutover | `TRUNC(now) <= 02/01/2026` → periode **`12.2025`** — direplikasi apa adanya |
| Keluaran | `p_seq_number = LPAD(seq, 5, '0')` |

**OQ-013 TERTUTUP** untuk Claim — Life. `[data DBA]` Procedure **tidak** commit sendiri; `COMMIT`
yang terlihat di korpus milik **pembungkus Pega**. `[keputusan work owner]` **Go memegang batas
transaksi** — buka transaksi, panggil proc, commit/rollback di Go, dan **commit segera setelah nomor
terbentuk** agar lock `FOR UPDATE` lekas lepas.

AC tambahan yang mengikuti penutupan ini sudah masuk daftar di atas.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 03 - Baris `AdjustmentList` + Save ke Outstanding

**Status:** ready-for-agent

**Blocked by:** 02 (register klaim + penomoran), **14 (skema relasional klaim — PREFACTOR)**

> ⚠️ **Diselaraskan 2026-09-16 — revisi penyimpanan.** Dua perubahan mengikat:
> **(1)** baris adjustment melekat pada **PESERTA** (`T_CLAIMLF_PREMIUMLIST_DETAIL`), bukan pada
> klaim; **(2)** dokumen **per peserta** menjadi **gerbang simpan** ke Outstanding.
> Lihat blok AC "Penyimpanan relasional" dan "Dokumen per peserta" di bawah.

#### Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat menginput baris **AdjustmentList** pertama pada sebuah klaim
dan menyimpannya ke Outstanding — sehingga klaim itu tercatat sebagai sedang berjalan dan menjadi
antrean kerja yang terlihat. *(User story 3 di spec)*

Ini tiket yang **memperkenalkan unit keputusan** sistem ini: barisnya, bukan klaimnya.

#### Area codebase

`internal/models` (baris adjustment sebagai entitas dengan statusnya sendiri), `internal/repository`
(penulisan baris), `internal/services` (aturan penyimpanan ke Outstanding), `internal/handlers`
(endpoint input baris), `frontend/` (formulir dan daftar baris adjustment).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` **satu-satunya penulis nilai `STS_REJECT = 0`** di seluruh korpus Claim Life |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` tombol "Save to Outstanding", dan bentuk baris adjustment |
| `Claim Life/Activity/SetIndexAdjustmentList.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` baris baru mewarisi **8 kolom** dari `AdjustmentList(1)`: `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`, `SUM_REASURED`, `SUM_INSURED`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `CURRENCYID`, `CURRENCY` — **tanpa** `STS_REJECT` |

#### ADR terkait

**ADR-U-0011** (unit status = baris `AdjustmentList`), **ADR-U-0003** (8 kolom uang + `EM_PERCENT`
persen; nama kolom tidak dapat dipakai menebak sifatnya).

#### Acceptance criteria

- [ ] Baris `AdjustmentList` yang baru disimpan ke Outstanding berstatus **Outstanding** (`0`).
      *(AC 1 spec)*
- [ ] Sebuah klaim dapat memuat **beberapa** baris adjustment sekaligus, masing-masing dengan
      statusnya sendiri. *(AC 25 spec)*
- [ ] Baris yang ditambahkan setelah baris pertama **mewarisi delapan kolom** di atas dari baris
      pertama, dan **tidak** mewarisi status. *(AC 5 spec)*
- [ ] Kedelapan kolom uang diperlakukan sebagai uang; `EM_PERCENT` **tidak** diperlakukan sebagai
      uang. *(AC 23 spec)*
- [ ] Tidak ada nilai uang sebagai *binary floating point* di lapisan mana pun maupun di kontrak API.
      *(AC 22 spec)*

##### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §2b

- [ ] ⚠️ Setiap baris adjustment menunjuk **satu peserta** lewat `PREMIUM_LIST_DETAIL_ID`. Test yang
      menemukan FK adjustment menunjuk **header klaim** **gagal**. *(AC 33 spec; penyimpangan
      sadar 2 — perbaikan relasi, bukan peniruan)*
- [ ] ⚠️ Dua peserta dengan masing-masing dua putaran adjustment menghasilkan **empat baris yang
      seluruhnya dapat ditelusuri ke peserta yang benar**. *(AC 34 spec)*
- [ ] ⚠️ Peserta menyimpan **penanda dipilih-untuk-diklaim** (`IS_CHECK`) — inilah penyimpan aturan
      "hanya peserta yang diklaim". *(AC 39 spec; penyimpangan sadar 3)*
- [ ] ⚠️ Tanggal **diterima**, **konfirmasi**, dan **penyelesaian** tersimpan **per peserta**, bukan
      di header. *(AC 40 spec)*
- [ ] Peserta menyimpan `STATUS`, `RECOMMENDATION`, `SOURCE_ID`, dan `CEDING_RETENTION`.
      *(AC 41 spec)*
- [ ] ⚠️ `STS_REJECT` baris adjustment diisi **nilai sebenarnya menurut aksi** — Admin insert
      Outstanding → `0`; SPV tambah Outstanding → `0`. Test yang menemukan nilai di-hardcode
      **gagal**. *(AC 42 spec; penyimpangan sadar 4)*
- [ ] ⚠️ `ACCEPTATION_DATE` **tidak** distempel saat insert; ia diisi **tanggal akseptasi
      sebenarnya** saat baris benar-benar diaksep. *(AC 43 spec; penyimpangan sadar 4)*
- [ ] Peserta beserta seluruh baris adjustment-nya tersimpan dalam **satu transaksi**. *(AC 49 spec)*
- [ ] Baris adjustment menyimpan **nama bank**, **id bank**, dan **nomor rekening**, dan ketiganya
      dapat diisi dari layar rincian adjustment. *(AC 56 spec; `[terverifikasi]` class
      `ASM-FW-GISFW-Data-AdjustmentLife`, tampil di `Claim Life/Section/AdjustmentDetail_Section.xml`)*
- [ ] Ketiga field bank **boleh kosong saat Save ke Outstanding** — ia baru menjadi gerbang pada
      **penyerahan ke Komite** (tiket 10). *(AC 57 spec)*

##### Dokumen per peserta — **gerbang simpan** ⚠️ BARU 2026-09-16

`[terverifikasi]` `Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) beriterasi atas
`.DocumentList` **per peserta** dan menolak simpan dengan
`"The document hasn't been uploaded person number "+<nomor>` serta
`"Documents are incomplete, please complete the documents"`.

- [ ] ⚠️ Dokumen tersimpan **per peserta** di `DOCUMENT_CLAIM` dan dapat dibaca dengan `SELECT`
      biasa — **bukan** lewat mekanisme lampiran bawaan. *(AC 44 spec; penyimpangan sadar 5)*
- [ ] ⚠️ Menyimpan ke Outstanding **ditolak** bila ada peserta yang dokumennya belum lengkap, dengan
      pesan yang **menyebut peserta mana**. *(AC 45 spec; penyimpangan sadar 5)*
- [ ] Kolom isian `DOCUMENT_CLAIM` **diturunkan dari sensus `.DocumentList`** pada activity di atas,
      dan **keputusannya dicatat** — **jangan tebak dari nama tabel**. *(tiket 14 §Catatan)*
- [ ] Halaman React menampilkan daftar baris adjustment dengan status masing-masing sebagai kata,
      bukan angka. *(AC 26 spec)*

##### Spreading adjustment ⚠️ BARU 2026-09-16

`[terverifikasi]` `Claim Life/Activity/SpreadingClaimLife_Act.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `SPREADINGCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY`) menghitung dan
mengisi `.SpreadingList` (per treaty-year, class `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`) beserta
`.RetroLifeList` (per reinsurer, class `ASM-FW-GISFW-Int-RETROCESSIONLIFE`) pada **baris
adjustment**.

- [ ] ⚠️ Menyimpan baris adjustment **menghitung dan menyimpan** spreading-nya: satu baris per
      treaty-year, dan di bawahnya satu baris per reinsurer. *(AC 59 spec; tiket 14; penyimpangan sadar —
      spreading dibekukan)*
- [ ] ⚠️ Setiap baris spreading menunjuk **satu baris adjustment**; setiap baris spreading retro
      menunjuk **satu baris spreading**. Test yang menemukan keduanya menggantung pada peserta atau
      pada header klaim **gagal**. *(AC 58 spec; tiket 14)*
- [ ] Nilai turunan tersimpan sesuai perhitungan yang terbukti: `AMOUNT = RetrocadedShare ×
      PERCENT_SHARE ÷ 100`, dan `PREMIUM_SPREADED_GROSS = RATE × (1 + EM_PERCENT) × AMOUNT`.
      *(`[terverifikasi]` `SpreadingClaimLife_Act`)*
- [ ] ⚠️ `[terbuka]` **`PREMIUM_SPREADED_NET` tidak dinyatakan selesai** sebelum Product + UW
      menetapkan rumusnya. `[terverifikasi]` korpus memuat **dua** cabang di rule yang sama —
      `GROSS − Comm` dan `GROSS − Discount − Comm` — dan **mana yang berlaku tidak terbaca**.
      **Jangan tebak.** *(tiket 14 §Blocker)*
- [ ] ⚠️ Nilai spreading **dibekukan**: perubahan master treaty sesudahnya **tidak mengubah** angka
      yang sudah tersimpan pada adjustment itu. *(AC 59 spec)*
- [ ] Adjustment tanpa retrosesi tersimpan dengan **nol baris** spreading — **bukan** kegagalan.
- [ ] Baris adjustment beserta seluruh spreading dan spreading retro-nya tersimpan dalam **satu
      transaksi**. *(AC 49 spec)*

#### Catatan

`[terverifikasi]` Empat kolom bernama `SHARE_*` / `*_SHARE` ternyata **uang**, bukan rasio — nama
kolom di korpus ini terbukti menipu (bandingkan `STS_REJECT`, yang nilai `1`-nya berarti *diaksep*).
Klasifikasi mengikat ada di **ADR-U-0003**; jangan menyimpulkan dari nama.

`[terbuka]` **OQ-060** (pemilik **Product+UW**) — apakah seluruh baris satu klaim wajib bermata uang
sama. `[terverifikasi]` dalam praktiknya seragam karena `CURRENCY` disalin dari baris 1, tetapi
strukturnya membolehkan campur. **Tidak memblokir tiket ini**; memblokir bentuk akhir tipe uang.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 04 - Mesin status per baris + aturan turunan "klaim selesai"

**Status:** ready-for-agent

**Blocked by:** 03 (baris `AdjustmentList` + Save ke Outstanding)

#### Hasil & nilai pengguna

Sebagai **pengguna mana pun**, saya melihat status tiap baris adjustment secara terpisah dan tahu
persis apa yang sudah diputuskan dan apa yang belum — dan sebuah klaim yang seluruh barisnya sudah
diputus tampak "selesai" sehingga antrean kerja saya bersih. *(User story 25 dan 28 di spec)*

Ini **inti spesifikasi**: tiket yang menetapkan bahwa yang diputuskan adalah baris, bukan klaim.

#### Area codebase

`internal/models` (status baris sebagai nilai tertutup), `internal/services` (transisi + aturan
turunan), `internal/handlers` (endpoint transisi), `frontend/` (tampilan status klaim dan baris).

Status klaim **dihitung**, bukan disimpan sebagai kolom mandiri — meskipun sistem lama
menyimpannya sebagai kolom.

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | penulis `0` |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` penulis `1` (6 `Property-Set`) dan `2` (2 `Property-Set`) |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `REJECTOSCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` penulis `2` dari sisi Admin |
| `Claim Life/Activity/SetSTS_Reject.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETSTS_REJECT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` penurunan status klaim → baris (`.STS_REJECT = Primary.STS_REJECT`) |

`[terverifikasi]` **Sensus penulis lengkap** ada di register **OQ-061**: tidak ada satu pun rule di
`Claim Life` maupun `Komite Claim Life` yang menulis `0` setelah `1` atau `2`.

#### ADR terkait

**ADR-U-0011** (unit = baris; kefinalan; arti tunggal nilai `2`), **ADR-U-0001** (jalur balik Komite
bekerja pada tingkat baris).

#### Acceptance criteria

- [ ] Baris yang sudah bernilai Aksep atau Ditolak **tidak dapat berubah lagi** melalui jalur mana
      pun. *(AC 2 spec)*
- [ ] Menolak sebuah baris **tidak** menutup klaim; klaim tetap dapat menerima baris baru.
      *(AC 4 spec)*
- [ ] Status "Ditolak" pada sebuah baris **selalu** berarti baris itu ditolak — **tidak pernah**
      berarti klaim selesai, apa pun sumber penolakannya.
- [ ] ⚠️ **Diselaraskan 2026-09-16:** pencerminan terjadi pada kolom `STS_REJECT` /
      `ACCEPTED_NO` di **`T_CLAIMLF_PREMIUMLIST_DETAIL`** dan **`T_GENERAL_CLAIM`** (spec §2b) — unit
      keputusannya tetap baris `T_CLAIMLF_ADJUSTMENT` (**ADR-U-0011**), yang kini menggantung pada
      **peserta**. *(AC 33 spec; penyimpangan sadar 2)*
- [ ] `PremiumListDetail` dan header klaim **selalu mencerminkan** baris adjustment terakhir, dan
      **tidak** ditulis sebagai status mandiri. *(AC 7 spec)*
- [ ] Klaim dilaporkan "selesai" **hanya** bila tidak ada baris berstatus Outstanding **dan** ada
      sekurangnya satu baris berstatus Aksep. *(AC 8 spec)*
- [ ] Bila tidak ada baris Outstanding dan tidak ada pula yang Aksep, klaim berada dalam keadaan
      **ditolak seluruhnya** — dan tetap dapat dilanjutkan dengan baris baru.
- [ ] Riwayat lengkap seluruh baris pada satu klaim dapat dilihat, sehingga putaran Komite terbaca.
      *(User story 27 spec)*

#### Catatan penutupan (2026-09-14)

**Aturan turunan "klaim selesai" disetujui work owner** `[keputusan work owner]` dan tercatat di
`CONTEXT.md` (Lampiran, butir 1):

- Status klaim adalah **turunan**, **bukan** kolom tersimpan.
- Unit keputusan = **baris `AdjustmentList`** (**ADR-U-0011**).
- **Header klaim mengikuti adjustment terakhir.**

`[terverifikasi 2026-09-14]` **Di Pega**, `AdjustmentList` **tidak punya tabel fisik** (kelas
`ASM-FW-GISFW-Data-AdjustmentLife`, berawalan `Data-` = embedded), dan barisnya di-persist ke
`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` yang kolomnya identik. Bukti:
`Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` /
`SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) — `pyStepsObjectName = .AdjustmentList`,
`pyStepsClassName = ASM-FW-GISFW-Data-AdjustmentLife`.

⚠️ **Diselaraskan 2026-09-16 — itu keadaan Pega, bukan keadaan sistem baru.** `[keputusan work
owner]` Di sistem baru baris adjustment punya **tabelnya sendiri**, `T_CLAIMLF_ADJUSTMENT`, yang
menggantung pada **peserta** (`T_CLAIMLF_PREMIUMLIST_DETAIL`) — bukan pada header klaim, dan **bukan**
pada `OS_AKSEPTASI_KLAIM_LIFE` (spec §2b, AC 33).

⚠️ **Koreksi 2026-09-16** `[keputusan work owner]`: `OS_AKSEPTASI_KLAIM_LIFE` **tetap di-`INSERT`
flat**, berdampingan dengan tabel relasional — hilir masih membaca dari sana. Yang **dibuang hanya
JSON**. Jadi baris adjustment tersimpan di **`T_CLAIMLF_ADJUSTMENT`** *dan* ikut terbawa ke rekam flat
itu; yang berubah adalah **di mana unit keputusan tinggal**, bukan hilangnya tabel akseptasi.
*(AC 32 spec)*

**Ini tidak mengubah ADR-U-0011** — unit keputusan tetap **baris**. Yang berubah hanya di mana baris
itu tinggal: dari tabel akseptasi bersama menjadi tabel adjungan per peserta.

`[data DBA]` `STS_REJECT` pada tabel **warisan** bertipe `NUMBER(38)`. ⚠️ Di skema baru tipenya
ditetapkan sendiri di tiket **14**; nilainya tetap `0`/`1`/`2` dan diisi **menurut aksi**, bukan
di-hardcode (AC 42 spec).

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 05 - Reject Outstanding oleh Admin — tanpa Komite

**Status:** ready-for-agent

**Blocked by:** 04 (mesin status per baris)

#### Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat menolak baris adjustment yang saya input sendiri tanpa melalui
Komite — sehingga kesalahan input dapat saya batalkan cepat — dan penolakan itu **hanya membatalkan
baris itu**, klaimnya tetap hidup dan saya dapat menginput baris pengganti. *(User story 6 dan 7 di
spec)*

Ini jalur penolakan **kedua** di sistem ini. Jalur pertama (lewat Komite) datang di tiket 11.

#### Area codebase

`internal/services` (aturan penolakan + gerbang status baris), `internal/handlers` (endpoint reject),
`frontend/` (kontrol "Reject Outstanding" pada baris adjustment).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `REJECTOSCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` menulis `2` ke **dua tingkat berbarengan**: `.STS_REJECT` dan `…PremiumListDetail(idx).STS_REJECT` |
| `Claim Life/Section/RejectOSClaimLife_Sec.xml` | layar pemicu — `[terverifikasi]` merujuk activity di atas 2× sebagai `<pyActivity>` | |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` gerbang tombol "Reject Outstanding": `pyWorkPage.pyPosition =='ReasLifeAdmin' && pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO !='' && .STS_REJECT == 0` |

Gerbang itu menguji **`.STS_REJECT` tingkat baris** — bukti bahwa wewenang pun diukur per baris.

#### ADR terkait

**ADR-U-0011** (dua sumber penolakan dengan makna setara pada tingkat baris; `RejectOSClaimLife_Act`
**bukan** dead rule), **ADR-U-0002** (peran).

#### Acceptance criteria

- [ ] Penolakan oleh `ReasLifeAdmin` membuat **hanya baris itu** berstatus Ditolak; baris lain pada
      klaim yang sama tidak berubah. *(AC 3 spec)*
- [ ] Klaim **tidak** tertutup oleh penolakan itu, dan baris adjustment baru dapat diinput
      sesudahnya.
- [ ] Penolakan hanya mungkin pada baris yang **masih Outstanding** — baris yang sudah Aksep atau
      Ditolak menolak upaya itu.
- [ ] Penolakan hanya mungkin bila klaim sudah punya nomor (padanan `CLAIM_NO !=''`).
- [ ] Nilai status hasil penolakan Admin **sama** dengan hasil penolakan Komite; tidak ada nilai
      khusus yang membedakan keduanya.
- [ ] Pencerminan ke tingkat `PremiumListDetail` terjadi bersamaan, bukan menyusul. *(AC 7 spec)*
- [ ] ⚠️ **Diselaraskan 2026-09-16:** penolakan langsung oleh `ReasLifeAdmin` menulis `STS_REJECT`
      = **`2`** sebagai **nilai sebenarnya menurut aksi** — bukan nilai yang di-hardcode seperti di
      Pega. Test yang menemukan nilai di-hardcode **gagal**. *(AC 42 spec; penyimpangan sadar 4)*

#### Catatan

`[keputusan work owner 2026-09-14]` Bahwa penolakan Admin membatalkan **baris saja** — bukan klaim —
berasal dari work owner. `[terverifikasi]` Korpus **tidak membedakan** kedua sumber penolakan: kedua
rule menulis `2` dalam bentuk yang sama persis. Perbedaan itu memang tidak seharusnya ada.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 06 - Validasi Date of Loss per `Type` + `ContentNote` dari `BusinessCode`

**Status:** ready-for-agent

**Blocked by:** 03 (baris `AdjustmentList` + Save ke Outstanding)

#### Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, sistem menolak *Date of Loss* yang berada di luar jendela valuasi polis —
sehingga klaim yang tidak tertanggung tidak pernah masuk ke siklus — dan jenis klaim terisi otomatis
dari kode produk sehingga saya tidak salah memilih. *(User story 4 dan 5 di spec)*

#### Area codebase

`internal/services` (aturan validasi + penurunan jenis klaim), `internal/models` (data acuan
`BusinessCode` → `ContentNote`), `internal/handlers` (pesan kesalahan yang dapat dibaca),
`frontend/` (penampilan pesan pada formulir).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/ValidasiDOL_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `VALIDASIDOL_ACT` / `RULE-OBJ-ACTIVITY`, 59.747 byte | dua cabang jendela valuasi per `Type` |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | penurunan `ContentNote` dari `BusinessCode` |

`[terverifikasi]` Dua cabang validasi DOL:

| Cabang | Jendela | Argumen `@addCalendar` |
| --- | --- | --- |
| `Type=="QR" \|\| Type=="QP"` | `GROSS_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `(.DATE_OF_LOSS,0,0,0,0,0,0,0)` |
| `Type=="TR" \|\| Type=="TP"` | `RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `(.DATE_OF_LOSS,0,0,0,1,0,0,0)` |

Gagal → `local.errmsg = "Invalid DOL"`, dengan `local.Begin==false || local.Expired==true`.

#### ADR terkait

**ADR-U-0012** (`Type` menggerbangi dua hal: wewenang **dan** jendela validasi; `TP` = Payable,
`TR` = Receivable — `[keputusan work owner]`, tidak ada di korpus).

#### Acceptance criteria

- [ ] Klaim ber-`Type` `QP`/`QR` dengan *Date of Loss* di luar `GROSS_VALUATION_BEGIN_DATE` …
      `_EXPIRED_DATE` ditolak dengan pesan yang setara `"Invalid DOL"`. *(AC 13 spec)*
- [ ] Klaim ber-`Type` `TP`/`TR` diuji terhadap `RETROCESSION_VALUATION_*`, dengan pergeseran
      tanggal yang sama seperti Pega. *(AC 14 spec)*
- [ ] `ContentNote` terisi sesuai tabel `BusinessCode` `L1`–`L21` di `CONTEXT.md`. *(AC 15 spec)*
- [ ] Pemetaan `BusinessCode` → `ContentNote` diwujudkan sebagai **data acuan**, bukan rangkaian
      `if` bercabang.
- [ ] Validasi dan penurunan jenis klaim membaca **satu** field `Type` yang sama — bukan dua salinan
      seperti di Pega.

#### Catatan penutupan (2026-09-14)

**DOL dibawa sebagai paritas** `[keputusan work owner]` — validasi *Date of Loss* per `Type`
dipertahankan persis seperti perilaku Pega. Tercatat di `CONTEXT.md` (Lampiran, butir 2).

`[terbuka — catatan kecil, NON-PEMBLOKIR]` Satuan pergeseran tanggal pada cabang `TP`/`TR`
(`@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)`) belum dikonfirmasi **Product+UW**; definisi
`@addCalendar` tidak ada di korpus. **Ini bukan blocker**: pergeseran direplikasi apa adanya dari
argumen yang terbaca. Bila kelak Product+UW menyatakan satuannya berbeda, itu perubahan satu baris
disertai testnya.

`[terbuka]` **OQ-020** — arti `QP`/`QR` belum dijawab. **Tidak memblokir**: perilaku kedua cabang
sudah terbaca penuh; hanya namanya yang belum.

`[data DBA]` Kolom terkait di `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`: `TYPE VARCHAR2(10)`.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 07 - Penegakan peran di lapisan layanan + wewenang kirim-Komite per `Type`

**Status:** ready-for-agent

**Blocked by:** 05 (reject Outstanding oleh Admin) — gerbang perlu tindakan nyata untuk dijaga

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin setiap tindakan pada klaim hanya dapat dilakukan peran yang
berhak — dan penegakannya berada di **lapisan layanan**, bukan sekadar di visibilitas layar —
sehingga wewenang tidak dapat dilewati dengan memanggil API langsung. *(User story 15 di spec)*

#### Area codebase

`internal/services` (penegakan wewenang sebelum setiap transisi), `internal/handlers` (identitas
pemanggil), `frontend/` (kontrol yang tampil/tersembunyi mengikuti peran — **sebagai kenyamanan,
bukan sebagai penegakan**).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | tiga gerbang `<pyCondition>` di bawah |
| `Claim Life/Flow/Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` 15 kemunculan `pyPosition` — penegakan sesungguhnya di sistem lama bertumpu pada penugasan tahap di sini |

`[terverifikasi]` Tiga gerbang pada `AdjustmentDetail_Section`:

| Kontrol | `<pyCondition>` |
| --- | --- |
| "Reject Outstanding" | `pyWorkPage.pyPosition =='ReasLifeAdmin' && …CLAIM_NO !='' && .STS_REJECT == 0` |
| Jalur Komite (`GetListKomiteLife`) | `pyWorkPage.pyPosition =='ReasLifeSPV' \|\| pyWorkPage.Type = 'TP' \|\| pyWorkPage.Type = 'TR'` |
| "Save to Outstanding" | `pyWorkPage.pyPosition =='ReasLifeSPV'` |

`[terverifikasi]` `Type` dibaca dari **dua salinan** di Pega — `pyWorkPage.Type` (gerbang wewenang)
dan `pyWorkPage.PolicyDataLife.Type` (validasi DOL) — disalin di
`Claim Life/Activity/LoadDataPeserta_Act.xml` (`ASM-FW-GCNMFW-WORK` / `LOADDATAPESERTA_ACT` /
`RULE-OBJ-ACTIVITY`): `pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type`. **Sumber otoritatif =
`PolicyDataLife.Type`.**

#### ADR terkait

**ADR-U-0002** (tiga peran; rangkap peran tidak diperbolehkan), **ADR-U-0012** (wewenang kirim-Komite
bergantung `Type`; **dibawa apa adanya sebagai paritas**, risiko RBAC diterima dan ditinjau saat
konteks Komite/IAM digarap).

#### Acceptance criteria

- [ ] `ReasLifeMedicalAdvisor` **tidak dapat** mengubah status akseptasi baris mana pun. *(AC 9 spec)*
- [ ] Untuk klaim ber-`Type` `QP` atau `QR`, **hanya `ReasLifeSPV`** yang dapat mengirim ke Komite;
      upaya oleh peran lain ditolak. *(AC 10 spec)*
- [ ] Untuk klaim ber-`Type` `TP` atau `TR`, `ReasLifeAdmin` **dapat** mengirim ke Komite.
      *(AC 11 spec)*
- [ ] Penolakan wewenang terjadi **di lapisan layanan**, dan tetap terjadi meskipun kontrol UI-nya
      ditampilkan. *(AC 12 spec)*
- [ ] Wewenang dan validasi membaca **satu** field `Type` yang sama — tidak ada dua salinan yang
      dapat berbeda.
- [ ] Tidak ada nama orang ter-hardcode di lapisan mana pun; wewenang diukur dari peran akun.

#### Catatan

⚠️ `[terverifikasi]` Penegakan di sistem lama **tidak seragam**: `pyPosition` tidak ada sama sekali
di `Section/InputOSClaimLife.xml`, `Section/RejectOSClaimLife_Sec.xml`, `Section/ClaimComite.xml`,
maupun `Harness/Committe_Life.xml`. **Ketidakseragaman itu tidak ditiru** — sistem baru menegakkan
di lapisan layanan untuk semua tindakan.

⚠️ Kelonggaran `TP`/`TR` **adalah risiko yang diterima sadar** (**ADR-U-0012**): wewenang bergantung
pada atribut data, bukan peran. Jangan "memperbaikinya" dalam tiket ini — itu keputusan terpisah.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 08 - Tahap Medical Check & Claim Analis + jalur balik

**Status:** ready-for-agent

**Blocked by:** 07 (penegakan peran) — tiap tahap milik peran tertentu

#### Hasil & nilai pengguna

Sebagai **ReasLifeMedicalAdvisor**, saya menerima klaim yang menunggu telaah medis dan dapat
mencatat hasilnya; sebagai **ReasLifeSPV**, saya menerima klaim yang siap dianalisis. Keduanya dapat
**mengembalikan** kasus ke admin bila data kurang, dan SPV dapat mengembalikan ke medis bila telaah
perlu diulang — sehingga tidak ada keputusan yang diambil di atas data tidak lengkap.
*(User story 8, 11–14, 23, 24 di spec)*

#### Area codebase

`internal/models` (posisi tahap pada klaim), `internal/services` (transisi antar tahap + jalur
balik), `internal/handlers` (endpoint submit dan kembalikan), `frontend/` (layar Medical Check dan
Claim Analis, kontrol kembalikan).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Flow/Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` empat tahap: `Assignment2` "Input Register", `Assignment1` "Outstanding Claim", `Assignment3` "Medical Check", `Assignment4` "Claim Analis" |
| `Claim Life/When/IsSendtoAdmin.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN` | `[terverifikasi]` `pyWorkPage.SendtoAdmin = 1`; menggerbangi pengembalian dari **tiga titik** di `Register_Flow` |
| `Claim Life/When/IsSendtoMedical.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN` | pengembalian SPV → medis; `[keputusan work owner]` artinya — kondisinya tidak terbaca dari tag |
| `Claim Life/Section/MedicalCheckClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `MEDICALCHECKCLAIMLIFE` / `RULE-OBJ-HTML-SECTION` | layar telaah medis; memuat gerbang `pyPosition` |
| FlowAction `MEDICALCHECK`, `AKSEPTASICLAIMLIFE` | — | `[terverifikasi]` tindakan pada kedua tahap |

#### ADR terkait

**ADR-U-0002** (peran per tahap), **ADR-U-0011** (tahap 2 dan 3 **tidak mengubah** status baris —
baris tetap Outstanding sepanjang Medical Check dan sampai keputusan Komite).

#### Acceptance criteria

- [ ] Klaim dapat berpindah Register → Outstanding → Medical Check → Claim Analis, masing-masing
      hanya oleh peran yang berhak.
- [ ] Status baris adjustment **tetap Outstanding** sepanjang perpindahan tahap — perpindahan tahap
      bukan keputusan akseptasi.
- [ ] `ReasLifeMedicalAdvisor` dapat mengembalikan kasus ke `ReasLifeAdmin`.
- [ ] `ReasLifeSPV` dapat mengembalikan kasus ke `ReasLifeAdmin`.
- [ ] `ReasLifeSPV` dapat mengembalikan kasus ke `ReasLifeMedicalAdvisor`.
- [ ] Kasus yang dikembalikan muncul kembali di antrean peran tujuan.
- [ ] Pengembalian **tidak** mengubah status baris adjustment mana pun.

#### Catatan

`[terverifikasi]` Pengembalian bukan jalur langka: `SendtoAdmin` menggerbangi pengembalian dari
**tiga titik** di `Register_Flow`. Perekaman pelakunya ditangani tiket **09** — di sistem lama kedua
penanda hanya menyimpan nilai `1`, tanpa pelaku dan tanpa waktu.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 09 - Jejak audit — siapa + kapan untuk setiap transisi dan setiap jalur balik

**Status:** ready-for-agent

**Blocked by:** 08 (tahap + jalur balik) — seluruh transisi harus ada dulu untuk dapat direkam

#### Hasil & nilai pengguna

Sebagai **auditor**, saya dapat mengetahui **siapa** dan **kapan** untuk setiap transisi status dan
setiap pengembalian kasus — sehingga setiap keputusan dapat dipertanggungjawabkan, dan pengembalian
kasus dapat ditelusuri. Sebagai **ReasLifeAdmin**, saya melihat siapa yang mengembalikan kasus
kepada saya dan kapan, sehingga saya tahu apa yang diminta.
*(User story 10, 29, 30 di spec)*

**Ini penyimpangan sadar dari sistem lama — sebuah perbaikan, bukan paritas.**

#### Area codebase

`internal/models` (entri jejak audit), `internal/repository` (penyimpanan jejak),
`internal/services` (perekaman pada setiap transisi — satu tempat, bukan tersebar),
`internal/handlers` (identitas pelaku), `frontend/` (tampilan riwayat pada klaim).

Jejak direkam **per baris `AdjustmentList`**, bukan per klaim — karena unit statusnya baris
(**ADR-U-0011**).

#### Rule Pega sumber

| Rule | Identitas | Keadaan sekarang |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` hanya `CREATEOPNAME` + empat kolom tanggal: `ACCEPTATION_DATE`, `CONFIRMATION_DATE`, `CLAIM_RECEIVED_DATE`, `COMPLETE_DATE` |
| `Claim Life/When/IsSendtoAdmin.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN` | `[terverifikasi]` hanya menyimpan nilai `1` — **tanpa pelaku, tanpa waktu** |
| `Claim Life/When/IsSendtoMedical.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN` | idem |
| `Claim Life/RDBList/InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` / `RNM!INSERTLOGSERVICECLAIM` / `RULE-CONNECT-SQL` | `[terverifikasi]` `INSERT INTO pooldata.monitoring_klaim_log` — **log layanan**, bukan jejak keputusan |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml`, `SendEmailKlaimLF.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `…` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` memakai `OperatorID.pyUserIdentifier` / `pyUserName` sebagai data |

#### ADR terkait

**ADR-U-0007** (jejak audit setiap transisi dan setiap jalur balik — penyimpangan sadar),
**ADR-U-0002** (pelaku diidentifikasi lewat akun berperan, bukan nama ter-hardcode),
**ADR-U-0011** (jejak per baris).

#### Acceptance criteria

- [ ] Setiap transisi status baris menghasilkan catatan berisi **pelaku dan waktu**. *(AC 16 spec)*
- [ ] Setiap pengembalian (`SendtoAdmin`, `SendtoMedical`) menghasilkan catatan berisi **pelaku dan
      waktu**. *(AC 17 spec)*
- [ ] Perubahan nilai `Type` menghasilkan catatan berisi pelaku dan waktu. *(AC 18 spec)* — `Type`
      menyentuh keamanan, bukan sekadar data (**ADR-U-0012**).
- [ ] Jejak melekat pada **baris** yang bersangkutan, dan riwayat satu klaim dapat dibaca utuh
      lintas seluruh barisnya.
- [ ] Rekam akseptasi lama tetap ditulis sebagaimana adanya — kontrak dengan Komite tidak berubah
      karena tiket ini.
- [ ] Pelaku dicatat sebagai identitas akun; **tidak ada nama orang ter-hardcode**.

#### Catatan

`[keputusan work owner 2026-09-14]` Jejak lama **tidak dapat direkonstruksi ke belakang** — data
sebelum cutover hanya punya `CREATEOPNAME` + empat tanggal. Riwayat transisi lengkap hanya ada untuk
kejadian **setelah** cutover. Konsekuensi ini disadari dan diterima.

`[terbuka]` **OQ-013** (pemilik **DBA**) — `COMMIT` berada di dalam blok PL/SQL, sehingga
atomisitas "tulis akseptasi + tulis jejak audit" belum dapat dipastikan. **Tidak memblokir** tiket
ini; memblokir jaminan atomisitasnya.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 10 - Kontrak Komite — penyerahan kasus

**Status:** ready-for-agent

**Blocked by:** 07 (penegakan peran + wewenang per `Type`)

#### Hasil & nilai pengguna

Sebagai **ReasLifeSPV**, saya dapat menyerahkan baris adjustment ke **Komite Life** untuk diputuskan
— dan Komite menerima nilai klaim beserta mata uangnya, sehingga mereka tidak perlu membaca balik ke
sistem ini. *(User story 17, 18, 19 di spec)*

Komite adalah **sistem luar**; tiket ini membangun **batasnya**, bukan isinya.

#### Area codebase

`internal/models` (muatan penyerahan), `internal/repository` (pembuatan rekam kasus Komite),
`internal/services` (aturan penyerahan + pemilihan roster), `internal/handlers` (endpoint serahkan),
`frontend/` (kontrol "Send ke Komite" pada baris adjustment).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/CreateKMTLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 121.652 byte | `[terverifikasi]` sepuluh langkah: `Property-Set` → `Call pxRetrieveReportData` → `Property-Set` ×3 → `Call pxAddChildWork` → `Obj-Refresh-And-Lock` → `Property-Set` → `Obj-Save` → `Call SendEmailKlaimLF` |
| `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml` | `ASM-FW-GCNMFW-INT-EMAILKOMITE` / `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`, 67.827 byte | `[terverifikasi]` pemilihan roster; menerima `Param.LIMIT_BOTTOM` dan `Param.STS_KLAIM` |
| `Claim Life/Section/ClaimComite.xml`, `Claim Life/Harness/Committe_Life.xml` | — | `[terverifikasi]` pemicu dari UI (`<pyActivity>CreateKMTLife_Act</pyActivity>`, 2× masing-masing) — **bukan** shape flow |

`[terverifikasi]` Kelas kasus anak: `ASM-FW-GCNMFW-Work-KomiteLife`. Muatan yang menyeberang
sekarang:

```
childPageKomite.CLMNO            childPageKomite.KomiteCount
childPageKomite.KomiteLoop       childPageKomite.IndexAdjustment
childPageKomite.IndexPremiumList childPageKomite.KomiteList(<APPEND>).KomiteID
childPageKomite.KomiteList(<LAST>).IDKomite / .KomiteAproval / .KomiteEmail
sisi induk: .IsKomite  .KomiteNo  .TotalKomite
```

#### ADR terkait

**ADR-U-0001** (tiga kontrak batas; muatan **diperluas** atas keputusan work owner),
**ADR-U-0003** (nilai uang yang menyeberang memakai representasi yang sama di kedua sisi),
**ADR-U-0012** (siapa yang boleh menyerahkan, bergantung `Type`).

#### Acceptance criteria

- [ ] Penyerahan membuat kasus anak berkelas Komite Life, membawa penunjuk **baris** yang diserahkan.
- [ ] Muatan penyerahan memuat **nilai klaim**, **`CURRENCY`**, dan **status baris saat penyerahan**
      — tiga hal yang **tidak** ada di sistem lama. *(AC 24 spec; `[keputusan work owner]`)*
- [ ] Nilai uang yang menyeberang memakai representasi yang sama dengan di dalam sistem — tidak
      dikonversi menjadi *floating point* di batas. *(AC 22 spec)*
- [ ] Penyerahan hanya mungkin pada baris yang **masih Outstanding**.
- [ ] **Gerbang rekening pembayaran** `[terverifikasi]`: penyerahan **ditolak** bila salah satu dari
      **nama bank**, **id bank**, atau **nomor rekening** pada baris yang diserahkan **kosong**,
      dengan pesan yang setara `"Name of bank cannot be empty"`. Ini **paritas perilaku existing**,
      bukan penyimpangan — sumbernya `Claim Life/Activity/GetListKomiteLife.xml`
      (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`), prasyarat
      `.NameOfBank=="" || .NoAccount=="" || .IDOfBank==""`. *(AC 57 spec)*
- [ ] Penolakan gerbang rekening terjadi di **lapisan layanan**, dan tetap terjadi meskipun kontrol
      UI-nya ditampilkan. *(sejalan AC 12 spec)*
- [ ] Penyerahan berlaku untuk **semua** adjustment, bukan hanya yang di atas ambang nilai tertentu.
      *(`[keputusan work owner]`, langkah 4 mesin status)*
- [ ] Wewenang penyerahan mengikuti aturan tiket 07: `QP`/`QR` hanya SPV; `TP`/`TR` bebas peran.
- [ ] Perubahan pada bentuk muatan diperlakukan sebagai **perubahan kontrak lintas konteks**, dan
      ditandai demikian di kode.
- [ ] **Jumlah tingkat komite = COUNT baris roster `EMAILKOMITE` yang aktif (`STS_AKTIF = "1"`) dan
      ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`** — dihitung saat penyerahan, **tidak** dibaca dari
      konstanta mana pun.
- [ ] Ambang yang dipakai mencari roster adalah **nilai mutlak** klaim: klaim bernilai negatif
      dicari dengan tandanya dihilangkan — `@if(IsADj<0, IsADj*-1, IsADj)`.
- [ ] ⚠️ `[keputusan work owner]` Bila tidak ada baris roster yang cocok, penyerahan **gagal
      terang-terangan** — bukan diam-diam membuat tangga nol tingkat. **Penjaga defensif**, bukan
      alur normal: roster dijamin **≥ 1** secara bisnis (limit berjenjang selalu menutup nilai
      klaim). `[terverifikasi]` Pega **tidak** punya gerbang ini — `CreateKMTLife_Act` men-set
      `KomiteCount = 1` meski `KomiteLoop = 0`.
- [ ] Semua baris adjustment pada satu klaim bermata uang sama sebelum penyerahan (invariant
      **OQ-060**).

##### Rujukan Komite ⚠️ BARU 2026-09-16

- [ ] ⚠️ Penyerahan ke Komite menyimpan **`KOMITE_ID`** = **identitas kasus komite**, yaitu
      `T_WORK_CLAIM.ID` baris komite yang baru lahir; baris yang belum pernah dikirim ber-`KOMITE_ID`
      **`NULL`**. *(AC 61 spec; tiket 14; penyimpangan sadar — rujukan, bukan salinan)* — ⚠️
      `[keputusan work owner]` REVISI 2026-09-17; bentuk penautan ini **tidak ada di korpus Pega**,
      lihat catatan korpus di tiket 14 §`KOMITE_ID`. `[terbuka]` tipenya mengikuti tipe
      `T_WORK_CLAIM.ID` yang belum ditetapkan.
- [ ] ⚠️ **Roster dan keputusan komite per anggota TIDAK disimpan di konteks ini.** Test yang
      menemukan tabel/kolom penyimpan `KomiteAproval`, `KomiteComment`, atau `DateApprove` di Claim
      Life **gagal**. Keduanya milik **Komite Claim Life**. *(AC 61 spec; tiket 14; penyimpangan sadar)*
- [ ] ⚠️ Keputusan komite **dibaca lewat join**, bukan disalin ke adjustment. Rantainya **tiga
      lompatan**, bukan satu: `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` → `T_WORK_CLAIM` (`ID` = `KOMITE_ID`,
      `COVER_KEY` = `ID` baris klaim) → `T_GENERAL_KOMITE` (**`ID` = `T_WORK_CLAIM.ID`, shared
      primary key — ⛔ tidak ada kolom `WORK_CLAIM_ID`, REVISI 2026-09-18**) →
      `T_KOMITE_KOMITELIST` (keputusan per anggota, diurut `KOMITE_URUT`). *(tiket 14 §`KOMITE_ID`;
      tiket 00 Komite Claim Life)* — ⚠️ `[keputusan work owner]` REVISI 2026-09-17; rantai ini
      **bentuk baru, tidak ada di korpus Pega**. Di Pega penautannya lewat `pxAddChildWork` +
      `CLMNO` + indeks posisi; lihat catatan korpus di tiket 14 §`KOMITE_ID`.
- [ ] ⚠️ Rujukan memakai **ID stabil**, bukan indeks posisi. Test yang menemukan padanan
      `IndexPremiumList` / `IndexAdjustment` sebagai kunci rujukan **gagal**. *(`[terverifikasi]`
      `Claim Life/Activity/CreateKMTLife_Act.xml` memakai `.pxListSubscript`; penyimpangan sadar; **AC 62 spec**)*
- [ ] ✅ **TERTUTUP 2026-09-17** — tabel komite **sudah** menampung keputusan per baris adjustment
      lewat **`T_GENERAL_KOMITE.ADJUSTMENT_ID`** (FK → `T_CLAIMLF_ADJUSTMENT.ID`), ditetapkan di
      **tiket 00 Komite Claim Life** (§Tabel, plus AC "`ADJUSTMENT_ID` berada di `T_GENERAL_KOMITE`,
      **BUKAN** di `T_WORK_CLAIM`"). Penunjuk dua arah — `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dan
      `T_GENERAL_KOMITE.ADJUSTMENT_ID` — diisi dalam **satu transaksi** saat kirim komite.

#### Catatan penutupan (2026-09-14)

**OQ-032 TERTUTUP** `[terverifikasi]` — tangga komite **sepenuhnya data-driven**, tanpa konstanta:

1. `Claim Life/Activity/GetListKomiteLife.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
   `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) set `Local.IsADj = .CLAIM_AMOUNT`, panggil report
   `FilterEmailKomiteWithLimit`, hasil ke `.KomiteList`.
2. `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml` (`ASM-FW-GCNMFW-INT-EMAILKOMITE` /
   `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`) filter
   `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM && .STS_KLAIM = Param.STS_KLAIM && .STS_AKTIF = "1"`.
3. `Claim Life/Activity/CreateKMTLife_Act.xml` set
   `childPageKomite.KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)`.

**OQ-037 TERTUTUP** — ambang adalah **data di tabel** `POOLDATA.EMAILKOMITE`
(`LIMIT_BOTTOM INTEGER`, `LIMIT_TOP INTEGER`, `STS_AKTIF`, `STS_KLAIM`), **bukan** hardcode seperti
di Komite Claim FacIn. ⚠️ `[data DBA]` roster **tanpa kolom mata uang** → pita berlaku atas satu
mata uang implisit.

**OQ-060 TERTUTUP** — bentuk uang `(amount, currency)` per baris, **invariant: satu klaim satu mata
uang**.

`[terbuka]` **OQ-035** — kepemilikan `serviceInsertArasapasClaimLife_act` (satu salinan, dua
konteks). **Tidak memblokir** tiket ini; menyentuh tiket 12.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 11 - Terima & tampilkan hasil keputusan Komite

> ⚠️ **Cakupan diubah 2026-09-15.** Judul lama: *"Kontrak Komite — jalur balik hasil keputusan"*.
> Tiket ini **tidak lagi menulis `STS_REJECT`**. `[terverifikasi]` Penulisannya milik **Komite Claim
> Life**, bukan Claim — Life: `Komite Claim Life/Activity/KomitePostAdjustment.xml`
> (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`) yang menulis ke
> dua tingkat baris. Implementasinya ada di **tiket Komite 05**.
> Tiket ini kini **membaca dan menampilkan** hasil itu, lalu melanjutkan siklus klaim.

**Status:** ready-for-agent

**Blocked by:** 10 (kontrak Komite — penyerahan) · **Komite 05** (jalur balik — penulisan
`STS_REJECT`)

#### Hasil & nilai pengguna

Sebagai **ReasLifeSPV**, saya melihat hasil keputusan Komite pada baris adjustment saya — sehingga
saya tahu apakah baris itu diaksep atau ditolak — dan bila ditolak, saya dapat menambah baris
adjustment baru untuk mengajukan ulang dengan angka yang diperbaiki.
*(User story 20, 21, 22 di spec)*

Inilah yang membuat **klaim tidak terminal**: penolakan menghasilkan putaran berikutnya, bukan akhir.

#### Area codebase

`internal/services` (pembacaan hasil; pembuatan baris lanjutan; perhitungan status klaim turunan),
`internal/handlers` (endpoint status klaim + tambah baris), `frontend/` (tampilan hasil dan riwayat
putaran).

**Tidak** menulis `STS_REJECT` — itu milik Komite 05.

#### Rule Pega sumber

| Rule | Identitas | Peran di tiket ini |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` **penulis** hasil — 6 `Property-Set` bernilai `1`, 2 bernilai `2`, digerbangi `KomiteCount == KomiteLoop`. Tiket ini **membacanya**, tidak menjalankannya |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` **satu rule bersama** — hash ternormalisasi `c50bfd9a12` identik di kedua modul |
| `Claim Life/Activity/SetIndexAdjustmentList.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY` | pewarisan 8 kolom ke baris lanjutan |

#### ADR terkait

**ADR-U-0001** (Kontrak 2 — jalur balik; `AcceptStatus` dipetakan ke `STS_REJECT` **di batas**, tidak
disimpan sebagai status kedua), **ADR-U-0011** (terminal per baris; klaim tidak terminal; revisi =
baris baru), **ADR-U-0007** (jejak audit).

#### Acceptance criteria

- [ ] Hasil keputusan Komite **terbaca** pada baris `AdjustmentList` yang diserahkan — Aksep atau
      Ditolak. *(AC 4 spec Claim Life)*
- [ ] Klaim **tetap** dapat menerima baris adjustment baru setelah penolakan Komite.
- [ ] Baris baru yang ditambahkan setelah penolakan berstatus Outstanding dan **mewarisi delapan
      kolom** dari baris pertama **tanpa** mewarisi status. *(AC 5 spec Claim Life)*
- [ ] `PremiumListDetail` dan header klaim **mencerminkan** baris terakhir setelah hasil diterapkan.
      *(AC 7 spec Claim Life)*
- [ ] Status klaim "selesai" dihitung sebagai keadaan **turunan** dari kumpulan baris — bukan kolom
      tersimpan.
- [ ] `AcceptStatus` **tidak** disimpan sebagai status kedua di konteks ini. *(**ADR-U-0001**)*
- [ ] Hasil keputusan **hanya diterapkan** ketika putaran Komite sudah mencapai **tingkat terakhir**;
      hasil dari tingkat antara **tidak** mengubah status baris mana pun di konteks ini.
      *(AC 6 spec)* ⚠️ **Penegakannya milik Komite Claim Life tiket 05** — tiket ini hanya wajib
      **tidak menerapkan lebih awal**. Dicatat agar AC 6 punya jejak pemilik, bukan tampak terlewat.
- [ ] Tiket ini **tidak** menulis `STS_REJECT` — diverifikasi dengan tidak adanya jalur tulis status
      baris di konteks Claim — Life.

#### Catatan — mengapa cakupan diubah

`[terverifikasi]` Penulisan `STS_REJECT` dilakukan rule di modul **`Komite Claim Life`** dengan class
`ASM-FW-GCNMFW-Work-KomiteLife`. Menempatkannya di tiket Claim — Life akan membuat **dua tiket
`ready-for-agent` sama-sama mengklaim penulisan yang sama** — dua agent dapat mengimplementasikannya
berdua.

`[keputusan work owner 2026-09-15]` Kepemilikan ditetapkan: **Komite menulis, Claim Life membaca.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 12 - Efek keluar asinkron + antre-ulang + flag lingkungan

**Status:** ready-for-agent

**Blocked by:** 09 (jejak audit) — kegagalan efek keluar harus masuk jalur audit, bukan hanya log

#### Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat mengunggah dan mengunduh dokumen pendukung klaim; sebagai
**ReasLifeSPV**, anggota Komite menerima notifikasi saat kasus diserahkan. Dan sebagai **pengguna
mana pun**, alur klaim saya **tidak pernah tertahan** karena layanan luar sedang gagal — kegagalan
itu tercatat dan diantre ulang, bukan hilang diam-diam.
*(User story 32–36 di spec)*

#### Area codebase

`internal/services` (interface efek keluar + orkestrasi asinkron + antre-ulang),
`internal/repository` (implementasi pemanggil), `internal/config` (alamat layanan + flag
lingkungan), `internal/handlers` (endpoint unggah/unduh), `frontend/` (kontrol dokumen).

#### Rule Pega sumber

| Efek | Rule | Identitas |
| --- | --- | --- |
| Unggah berkas | `Claim Life/Activity/InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `INSERTGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY`, 160.027 byte |
| Ambil URL / hapus | `GetUrlGoogleStorage_Act.xml`, `DeleteGoogleStorage_Act.xml` | kelas sama |
| Token penyimpanan | `Claim Life/RDBList/GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETTOKENSTORAGE_SQL` / `RULE-CONNECT-SQL` |
| Email | `Claim Life/Activity/SendEmailKlaimLF.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `SENDEMAILKLAIMLF` / `RULE-OBJ-ACTIVITY` |
| Arasapas | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` | `[terverifikasi]` **satu-satunya salinan di korpus** — OQ-035 |
| ⚠️ ~~Konversi ke produksi lewat payload JSON~~ — **DIBUANG 2026-09-16** | ~~`Claim Life/ConnectREST/convertJsonNusareToProductionClaimLife.xml`~~ | `[keputusan work owner]` hilir **`SELECT` langsung dari tabel klaim** (spec §2b, §12) |
| Pencatatan layanan | `Claim Life/RDBList/InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` / `RNM!INSERTLOGSERVICECLAIM` / `RULE-CONNECT-SQL` → `INSERT INTO pooldata.monitoring_klaim_log` |

`[terverifikasi]` Token penyimpanan **tidak** diterbitkan aplikasi — ia datang dari Oracle:
`pooldata.GET_TOKEN_STORAGE({UploadDoc.App}, {OperatorID.pyUserIdentifier}, {UploadDoc.Kodestring OUT}, {DocAPI.ResponseMsg OUT})`.

#### ADR terkait

**ADR-U-0008** (asinkron, tidak memblokir, antre-ulang), **ADR-U-0010** (penyimpanan tetap Google
Storage), **ADR-U-0013** (**resolusi endpoint lewat `M_LINK_SERVICE`** — meralat **ADR-U-0004** untuk
alamat layanan keluar), **ADR-U-0005** (flag lingkungan), **ADR-U-0007** (kegagalan masuk jalur audit).

#### Acceptance criteria

- [ ] ⚠️ **Diselaraskan 2026-09-16:** efek keluar tinggal **tiga** — berkas, email, Arasapas.
      Konversi ke produksi **tidak lagi mengirim payload JSON**; hilir membaca **langsung dari
      tabel klaim**. Test yang menemukan payload JSON dikirim keluar **gagal**. *(AC 55 spec;
      penyimpangan sadar 1)*
- [ ] Seluruh efek keluar berada **di balik interface**, sehingga dapat diganti dalam test.
- [ ] Kegagalan efek keluar mana pun **tidak menahan** transisi status klaim. *(AC 19 spec)*
- [ ] Kegagalan tercatat di **jalur audit**, bukan hanya di log layanan, dan **dapat diantre ulang**.
      *(AC 20 spec)*
- [ ] Di lingkungan non-production, klaim **tetap tersimpan**; keempat efek keluar tidak berjalan.
      *(AC 21 spec)*
- [ ] Tidak ada host, endpoint, atau kredensial sebagai literal di kode.
- [ ] Kegagalan konfigurasi dapat dibedakan dari kegagalan jaringan, agar antre-ulang tidak berputar
      sia-sia.
- [ ] Dokumen yang sudah diunggah dapat diunduh kembali.
- [ ] Alamat endpoint keluar di-resolve lewat **runtime lookup** ke `M_LINK_SERVICE` dengan kunci
      `(KATEGORI_1, KATEGORI_2)` — untuk Arasapas: `("Klaim", "insertClaimLife")`.
- [ ] **Tidak ada URL** sebagai literal, konstanta, **maupun env var** di kode. Yang boleh menjadi
      konstanta hanyalah **kunci kategori**. *(**ADR-U-0013**)*
- [ ] Pemisahan dev–prod terjadi lewat **isi tabel per-database**, bukan lewat percabangan di kode.
- [ ] Bila kunci kategori tidak ditemukan di `M_LINK_SERVICE`, kegagalan **terang-terangan** dan
      masuk jalur audit — bukan diam-diam melewati efek keluar.

#### Catatan penutupan (2026-09-14)

**OQ-047 TERTUTUP** `[terverifikasi + data DBA]` — kontrak resolusi endpoint diketahui penuh:
`Claim Life/Activity/GetLinkService.xml` (`ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` /
`RULE-OBJ-ACTIVITY`) melakukan `Obj-Browse` atas `M_LINK_SERVICE` dengan
`.KATEGORI_1 = Param.Kategori_1 AND .KATEGORI_2 = Param.Kategori_2`, mengambil `.URL`, lalu
`Connect-REST`.

`[terverifikasi]` Kunci Claim — Life terbaca di
`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`:
`Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`.
`[data DBA]` Isi tabel 19 baris; endpoint Life →
`<konfigurasi alamat layanan>

**OQ-018 TERTUTUP untuk Claim — Life** `[terverifikasi]` — **nol** URL endpoint bisnis ter-hardcode
di modul ini; seluruh `http(s)://` yang ada hanyalah `pyHelpURI` ke `community.pega.com` (132+) dan
3 tautan penampil dokumen Office. Pembedaan lingkungan memakai
`Claim Life/When/IsPEGAPROD.xml` (`@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN`):
`pzProductionLevel = "5"`.

**Flag lingkungan — keputusan spesifikasi tetap berlaku:** flag hanya menggerbangi **efek keluar**,
tidak pernah **penyimpanan**. Ini penyimpangan sadar dari Pega (di sana `IsPEGAPROD` juga
menggerbangi simpan utama); tanpa itu lingkungan non-production tidak dapat dipakai menguji.

`[terbuka]` **OQ-035** — kepemilikan `serviceInsertArasapasClaimLife_act`. **Tidak memblokir.**

#### Catatan

`[terverifikasi]` "Tidak memblokir" adalah **paritas**, bukan perubahan: `Claim Life` hanya punya
**pencatatan** (`InsertLogServiceClaim`), bukan gerbang keberhasilan seperti `IsSuccessHitService`
di konteks facultative. Yang **ditambahkan** adalah antre-ulang.

⚠️ Konsekuensi yang diterima (**ADR-U-0008**): sebuah klaim dapat mencapai keadaan akhir sementara
efek keluarnya masih tertunda.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 13 - Migrasi data penuh Claim — Life

**Status:** ready-for-agent

**Blocked by:** 11 (kontrak Komite — jalur balik), **14 (skema relasional klaim — PREFACTOR)**

> ⚠️ **Dipersempit 2026-09-16 — revisi penyimpanan.** **Pemindahan bentuk data klaim
> (JSON + tabel flat warisan → enam tabel relasional, beserta normalisasi atribut polis) berpindah
> ke tiket 14**, yang menjadi PREFACTOR. Tiket ini **tinggal** menangani hal yang tidak menyentuh
> bentuk: keutuhan riwayat setelah cutover, penomoran, pembacaan master, dan kolom yang tidak
> ditafsirkan. AC yang berpindah ditandai di bawah.

#### Hasil & nilai pengguna

Sebagai **work owner**, saya ingin seluruh data klaim Life dipindahkan — beserta **semua barisnya** —
sehingga tidak ada pekerjaan tertinggal di Pega dan riwayat putaran Komite tidak hilang.
*(User story 40, 41, 42 di spec)*

#### Area codebase

`internal/repository` (pembacaan sumber + penulisan sasaran), skrip migrasi + DDL produksi di dalam
`OUTPUT_HASIL_RNM/`. Tidak menyentuh `handlers`/`frontend`.

#### Rule Pega sumber

| Rule | Identitas | Dipakai untuk |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` 55 nama kolom `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` terbaca langsung dari blok PL/SQL |
| `Claim Life/RDBList/GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL` | penomoran yang tidak boleh melompat/mengulang setelah migrasi |

`[terverifikasi]` Tabel akseptasi Life **terpisah** dari lini lain: `OS_AKSEPTASI_KLAIM_LIFE`
(7 kemunculan, hanya `Claim Life` + `Komite Claim Life`) versus `OS_AKSEPTASI_KLAIM` (191 kemunculan,
enam modul non-Life). Migrasi ini **tidak menyentuh** tabel non-Life.

#### ADR terkait

**ADR-U-0009** (migrasi penuh; koeksistensi **ditolak**; cutover memutus, bukan bertahap),
**ADR-U-0011** (seluruh **baris** ikut pindah, bukan hanya keadaan terakhir),
**ADR-U-0006** (penomoran), **ADR-U-0007** (jejak audit tidak dapat direkonstruksi ke belakang).

#### Acceptance criteria

⚠️ **Berpindah ke tiket 14** (jangan dikerjakan di sini): seluruh klaim + seluruh baris adjustment
terbawa; jumlah baris per klaim utuh; presisi uang; desimal presisi arbitrer; migrasi dapat
dijalankan ulang; normalisasi atribut polis. Rujukan silang: AC 51–54 spec.

- [ ] Klaim yang sedang berada di tengah siklus terbawa **beserta statusnya** dan posisi tahapnya.
      ⚠️ **Diselaraskan:** posisi tahap mendarat di **`T_WORK_CLAIM`** (spec §2b), bukan di header
      klaim. *(AC 46 spec; penyimpangan sadar 6)*
- [ ] Setelah migrasi, penomoran klaim **tidak melompat dan tidak mengulang**. *(AC 42 spec)*
- [ ] Tidak ada periode dua penulis ke rekam akseptasi Life dari sisi Claim — Life. *(ADR-U-0009)*
      ⚠️ **Diselaraskan:** setelah cutover, penulis satu-satunya adalah keenam tabel klaim baru;
      `OS_AKSEPTASI_KLAIM_LIFE` **tetap ditulis** — ⚠️ **koreksi 2026-09-16**
      `[keputusan work owner]`: setelah cutover, penulis klaim adalah **kedelapan tabel relasional
      baru** *dan* `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE`, karena hilir masih membaca dari sana.
      Yang **dibuang hanya JSON**. Test yang **menolak** penulisan `OS_AKSEPTASI_KLAIM_LIFE` justru
      **gagal**. *(AC 32 spec)*
- [ ] ⚠️ Master `RATE_LIFE`, `PRODUCTINWARD_LIFE`, `CURRENCY` diperlakukan sebagai **view atas
      `JSONDATA`**, bukan sebagai tabel relasional biasa. **Pengecualian 2026-09-16:** **produk**
      dibaca dari **`product_life` relasional** (hasil migrasi Master Product Name Life), bukan dari
      `m_product_life.JSONDATA`. *(AC 38 spec)*
- [ ] `NO_SEQ` terjaga **per kombinasi `(CLASS, JENIS, TAHUN)`**, bukan sebagai penghitung global.
- [ ] Kolom `LAYER_1`…`LAYER_4` dipindahkan apa adanya **tanpa ditafsirkan** — perannya belum
      terverifikasi.

#### Catatan penutupan (2026-09-14)

**OQ-001 TERTUTUP** untuk Claim — Life `[data DBA]` — DDL **12 tabel/view** diserahkan, mencakup
**seluruh** persistensi Claim Life. Daftar lengkap dan tipenya di `discovery/open-questions.md`
OQ-001. Yang mengikat tiket ini:

| Temuan | Konsekuensi untuk migrasi |
| --- | --- |
| Kolom uang = Oracle **`NUMBER` tanpa presisi** | Go **wajib** desimal presisi arbitrer; **`float64` dilarang** (**ADR-U-0003**) |
| `STS_REJECT NUMBER(38)` di tabel klaim | ⚠️ berbeda dari `STS_REJECT VARCHAR2(15)` di `EMAILKOMITE` — nama sama, tabel berbeda, tipe berbeda |
| `RATE_LIFE`, `PRODUCTINWARD_LIFE`, `CURRENCY` = **VIEW atas `JSONDATA`** | master sesungguhnya **JSON**; membacanya sebagai tabel relasional biasa akan menyesatkan |
| `TANGGAL_CLOSING.TANGGAL VARCHAR2(10)` | tanggal disimpan sebagai **string**, bukan `DATE` |
| `LAYER_1`…`LAYER_4 VARCHAR2(10)` | ⚠️ **peran belum terverifikasi** — ikut dipindah apa adanya, jangan ditafsirkan |

**`AdjustmentList` tidak punya tabel fisik** `[terverifikasi]` — kelasnya
`ASM-FW-GISFW-Data-AdjustmentLife` (berawalan `Data-` = embedded, tanpa tabel), induk
`ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL`. Bukti:
`Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` /
`SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`). **Baris adjustment di-persist ke
`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`** — jadi "memindahkan seluruh baris" berarti memindahkan seluruh
baris tabel itu, bukan mencari tabel adjustment yang tidak pernah ada.

**OQ-002 TERTUTUP** — sequence per `(class, jenis, tahun)` di `GENERATE_SEQUENCE_NUMBER`
(PK komposit, `TAHUN VARCHAR2(5)`, `NO_SEQ NUMBER`, `MM_YYYY VARCHAR2(10)`). Migrasi nomor kini
dapat direncanakan: yang harus dijaga adalah **`NO_SEQ` per kombinasi kunci**, bukan satu penghitung
global.

`[terbuka]` **OQ-018** — lingkungan `pega-nusre` belum dinyatakan (IT-infra). **Tidak memblokir
penulisan skrip migrasi**; memblokir **pemilihan sumber data** saat eksekusi.

`[terbuka]` **OQ-013 untuk modul non-Life** — tidak menyentuh migrasi Claim Life.

#### Catatan

`[keputusan work owner 2026-09-14]` Opsi koeksistensi — sistem baru hanya menerima klaim baru
sementara klaim berjalan diselesaikan di Pega — **ditolak**. Penolakan itu dicatat di **ADR-U-0009**
supaya tidak diusulkan ulang.

`[keputusan work owner]` Jejak audit lama **tidak dapat direkonstruksi**: data sebelum cutover hanya
punya `CREATEOPNAME` + empat kolom tanggal. Riwayat transisi lengkap hanya ada setelah cutover.

#### Perintah verifikasi

```
go test ./internal/...
make check
```

Tambah verifikasi khusus migrasi: hitungan klaim dan hitungan baris per klaim sebelum dan sesudah
pemindahan harus sama.

## Claim Life - 14 - Skema relasional klaim (6 tabel, berakar di `T_WORK_CLAIM`) + migrasi — **PREFACTOR**

**Status:** ready-for-agent

**Blocked by:** 01 (kerangka aplikasi + seam API)

⚠️ **Ini PREFACTOR dan harus dikerjakan LEBIH DULU daripada tiket 02–13.** Nomornya 14 karena
seri 01–13 sudah terbit dan tidak dinomori ulang — **urutan ditentukan tepi pemblokir, bukan
nomor**. Tiket **02, 03, 04, 05, 12, 13** kini memblokir pada tiket ini.

Alasannya: bentuk penyimpanan berubah total — dokumen JSON + satu tabel flat warisan → **enam tabel
relasional**. Tidak ada irisan lain yang dapat berdiri sebelum bentuk barunya ada.
*"Make the change easy, then make the easy change."*

#### Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin setiap atribut klaim menjadi **kolom bernama** dengan tipe yang
benar, dan seluruh klaim lama pindah **tanpa kehilangan satu nilai pun** — termasuk **seluruh baris
adjustment**, bukan hanya keadaan terakhir. Sebagai **organisasi**, saya ingin bentuk klaim dapat
diperiksa, dicari, dan divalidasi, bukan tersembunyi di dalam satu dokumen. *(Spec §2b, §14)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL enam tabel klaim + `T_WORK_CLAIM` + sequence; skrip migrasi & rekonsiliasi |
| `internal/models` | Agregat klaim: header → peserta → adjustment & dokumen |
| — | Skrip rekonsiliasi nilai uang, tanggal, dan jumlah baris |

#### Bentuk baru — spec §2b RALAT D/E (2026-09-18)

⚠️ **REVISI 2026-09-18** `[keputusan work owner]` — **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING`
DIHAPUS**, dan **`T_WORK_CLAIM` naik menjadi AKAR.** Skema klaim Life turun dari **8 tabel menjadi
6**, dan pohonnya kini **enam tingkat**. ⛔ Jangan membuat `T_CLAIMLF_POLICY` maupun
`T_CLAIMLF_MARKETING` — penghapusan ini **bukan** penggantian nama.

```
TINGKAT 1   T_WORK_CLAIM ─────────────────── akar · LINTAS-LINI · satu baris per work object
            PK  ID   TEKS BERFORMAT: CLM-xxxxxx (klaim) / KMT-xxxxxx (komite)
            FK  COVER_KEY → T_WORK_CLAIM.ID   (menunjuk dirinya sendiri; NULL bila tak punya induk)
            LINI · PY_POSITION · ACCEPT_STATUS · SENDTO_ADMIN · SENDTO_MEDICAL · TYPE
            CASEID · CREATE_OP · CREATE_OP_NAME · TGL_UPDATE  <- PINDAHAN dari header klaim
            |
   +--------+-------------------------------------------+
   | baris KLAIM  COVER_KEY = NULL                      | baris KOMITE  COVER_KEY = ID baris klaim
   |                                                    |
TINGKAT 2                                           TINGKAT 2
   +--1:1-- T_GENERAL_CLAIM                             +--1:1-- T_GENERAL_KOMITE
           PK ID = T_WORK_CLAIM.ID   <- SHARED PK               PK ID = T_WORK_CLAIM.ID baris komite
             (tidak ada kolom FK terpisah)                        <- SHARED PK, tidak ada kolom FK
           penunjuk polis (ke LUAR, bukan anak):                FK ADJUSTMENT_ID -> T_CLAIMLF_ADJUSTMENT.ID
             CASEID_POLICY   [terbuka]                            (tutup lingkar)
             POLICY_NO       [terbuka]                          KOMITE_LOOP · KOMITE_COUNT · ACCEPT_STATUS
             ENDORSMENT_NO   [terbuka]                          |
           |                                         TINGKAT 3  |
TINGKAT 3  |                                                    +--1:N-- T_KOMITE_KOMITELIST
           +--1:N-- T_CLAIMLF_PREMIUMLIST_DETAIL                        PK ID
           |        PK ID                                              FK DATA_KOMITE_ID
           |        FK CLAIM_ID -> T_GENERAL_CLAIM.ID  CASCADE             -> T_GENERAL_KOMITE.ID  CASCADE
           |        |                                                   KOMITE_URUT · KOMITE_ID · ID_KOMITE
TINGKAT 4  |        +--1:N-- T_CLAIMLF_ADJUSTMENT                       KOMITE_EMAIL · KOMITE_APROVAL
           |        |        PK ID                                      KOMITE_COMMENT · DATE_APPROVE
           |        |        FK PREMIUM_LIST_DETAIL_ID
           |        |           -> T_CLAIMLF_PREMIUMLIST_DETAIL.ID  CASCADE
           |        |        FK KOMITE_ID -> T_WORK_CLAIM.ID  <- penunjuk balik KE ATAS, nullable
           |        |        |
TINGKAT 5  |        |        +--1:N-- T_CLAIMLF_ADJUSTMENT_SPREADING
           |        |                 FK ADJUSTMENT_ID -> T_CLAIMLF_ADJUSTMENT.ID  CASCADE
           |        |                 |
TINGKAT 6  |        |                 +--1:N-- T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO
           |        |                          FK SPREADING_ID
           |        |                             -> T_CLAIMLF_ADJUSTMENT_SPREADING.ID  CASCADE
TINGKAT 4  |        +--1:N-- DOCUMENT_CLAIM                      <- LINTAS-LINI
           |                 FK PREMIUM_LIST_DETAIL_ID -> ...DETAIL.ID   (Life saja)
           |                 ON DELETE di Go -- induk beda tabel per lini
           |
           +-- DIBACA dari luar, BUKAN anak, TIDAK ikut cascade:
                 T_PREMIUM_LIST dkk (polis) · master marketing officer · EMAILKOMITE ·
                 master retro · security reinsurer · currency
```

**ENAM tabel klaim** — `T_GENERAL_CLAIM`, `T_CLAIMLF_PREMIUMLIST_DETAIL`, `T_CLAIMLF_ADJUSTMENT`,
`T_CLAIMLF_ADJUSTMENT_SPREADING`, `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`, `DOCUMENT_CLAIM` —
digantung pada akar `T_WORK_CLAIM`. Kedalaman **enam tingkat**.

✅ **Ejaan penaut induk-anak — DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.**
`[keputusan work owner]` **Ejaan final `COVER_KEY` (snake_case);** `CoverKey` warisan Pega
**tidak dipakai** sebagai nama kolom.
Alasan: seluruh kolom SQL baru proyek ini snake_case, dan Oracle melipat identifier tanpa kutip
menjadi huruf besar — `CoverKey` akan menjadi identifier **tanpa garis bawah**, berbeda dari
kolom-kolom sekitarnya. Kolomnya dibuat oleh **tiket 00 Komite Claim Life** lewat `ALTER`;
tiket ini hanya membuat tabel `T_WORK_CLAIM`.

##### Sebelas relasi — seluruh kunci tamu ber-index

| # | Induk | Anak | Kunci tamu | Kard | Hapus |
| --- | --- | --- | --- | --- | --- |
| 1 | `T_WORK_CLAIM` | `T_WORK_CLAIM` | `COVER_KEY` | 1:N | di Go |
| 2 | `T_WORK_CLAIM` | `T_GENERAL_CLAIM` | **tidak ada kolom terpisah** — `T_GENERAL_CLAIM.ID` = `T_WORK_CLAIM.ID` (**shared PK**) | 1:1 | — |
| 3 | `T_GENERAL_CLAIM` | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `CLAIM_ID` | 1:N | CASCADE |
| 4 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `T_CLAIMLF_ADJUSTMENT` | `PREMIUM_LIST_DETAIL_ID` | 1:N | CASCADE |
| 5 | `T_CLAIMLF_ADJUSTMENT` | `T_CLAIMLF_ADJUSTMENT_SPREADING` | `ADJUSTMENT_ID` | 1:N | CASCADE |
| 6 | `T_CLAIMLF_ADJUSTMENT_SPREADING` | `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` | `SPREADING_ID` | 1:N | CASCADE |
| 7 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `DOCUMENT_CLAIM` | `PREMIUM_LIST_DETAIL_ID` | 1:N | **di Go** |
| 8 | `T_WORK_CLAIM` | `T_GENERAL_KOMITE` | **tidak ada kolom terpisah** — `T_GENERAL_KOMITE.ID` = `T_WORK_CLAIM.ID` baris komite (**shared PK**) | 1:1 | di Go |
| 9 | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | 1:N | CASCADE |
| 10 | `T_CLAIMLF_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` | 1:1 | di Go |
| 11 | `T_CLAIMLF_ADJUSTMENT` | `T_WORK_CLAIM` | `KOMITE_ID` | N:1 | penunjuk |

Relasi **8 · 9 · 10** dan kolom `COVER_KEY` dibuat oleh **tiket 00 Komite Claim Life**, bukan di
sini. Tiket ini membuat relasi **1 · 2 · 3 · 4 · 5 · 6 · 7 · 11** dan **tabel** `T_WORK_CLAIM`.

✅ **REVISI 2026-09-18 — relasi 2 dan 8 memakai SHARED PRIMARY KEY.** `[keputusan work owner]`
Keduanya **tidak punya kunci tamu**: `T_GENERAL_CLAIM.ID` **sama persis** dengan `T_WORK_CLAIM.ID`
baris klaim, dan `T_GENERAL_KOMITE.ID` **sama persis** dengan `T_WORK_CLAIM.ID` baris komite.
⛔ **Kolom `WORK_CLAIM_ID` DIBUANG** — bukan diganti nama, **tidak ada**. Test yang menemukan
kolom `WORK_CLAIM_ID` **gagal**. Karena keduanya PK, indexnya sudah ada dengan sendirinya.

✅ **`DOCUMENT_CLAIM` — `ON DELETE` DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.**
`[keputusan work owner]` Penghapusannya **ditangani di Go**, **bukan** cascade basis data. Alasan:
ia **LINTAS-LINI** dan induknya **tabel yang berbeda per lini**, sehingga satu `ON DELETE CASCADE`
tidak dapat seragam. Untuk Life induknya tetap `T_CLAIMLF_PREMIUMLIST_DETAIL`.

##### Kolom yang bertambah dan yang pindah — 2026-09-18

| Tabel | Aksi | Kolom |
| --- | --- | --- |
| `T_GENERAL_CLAIM` | **TAMBAH** | `CASEID_POLICY` · `POLICY_NO` · `ENDORSMENT_NO` — ketiganya ⚠️ `[terbuka]` |
| `T_GENERAL_CLAIM` | **BUANG** | `PL_NUMBER` → diganti nama menjadi **`POLICY_NO`** |
| `T_GENERAL_CLAIM` | **BUANG** | `CREATE_OP` · `CREATE_OP_NAME` · `TGL_UPDATE` → **pindah** ke `T_WORK_CLAIM` |
| `T_GENERAL_CLAIM` | **BUANG** | `CASEID` → **pindah** ke `T_WORK_CLAIM` *(DIPUTUSKAN 2026-09-18)* |
| `T_WORK_CLAIM` | **TAMBAH** | `CREATE_OP` · `CREATE_OP_NAME` · `TGL_UPDATE` (pindahan) |
| `T_WORK_CLAIM` | **TAMBAH** | `CASEID` (pindahan) *(DIPUTUSKAN 2026-09-18)* |
| `T_WORK_CLAIM` | **TAMBAH** | `LINI` — kolom **ADA**; nilai Life = konstanta lini Life. `[terbuka — Non-Life]` hanya **daftar enum lintas-lini** |

⚠️ `[terbuka]` **Ketiga penunjuk polis belum menunjuk apa pun yang ada.** Ditulis apa adanya di atas
**bukan** karena sudah beres. **Jangan dijawab sendiri** — pemilik **work owner**:

| Penunjuk | Mengapa belum menunjuk |
| --- | --- |
| `CASEID_POLICY` | Menunjuk **kolom apa?** `T_PREMIUM_LIST` **sengaja membuang** `pyID`/`px*`/`py*` — catatannya berbunyi *"ID work Pega diganti ID sequence baru"*. **Case id Pega milik polis tidak disimpan di sisi sana.** |
| `POLICY_NO` | **Tidak ada di header polis.** Ia ada di `T_PREMIUM_LIST_DETAIL` — **per peserta**. Header polis punya `PL_NUMBER`. Join header-klaim → header-polis lewat `POLICY_NO` **tidak nyambung.** |
| `ENDORSMENT_NO` | **Belum menunjuk satu versi.** Di sisi polis namanya **`NOENDORS`** `[terverifikasi]` (`Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml`), dan **versi berjalan ditentukan `PRODKE` terbesar**. |

✅ **Nasib `CASEID` — DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.** `[keputusan work owner]`
`CASEID` **PINDAH** ke `T_WORK_CLAIM`. Alasannya **sama** dengan
`CREATE_OP`/`CREATE_OP_NAME`/`TGL_UPDATE`: ia identitas **work object**, bukan atribut klaim.
`T_GENERAL_CLAIM` **tidak lagi memuat** `CASEID`.

##### Mengapa kedua tabel boleh dihapus `[terverifikasi]`

| Tabel | Bukti |
| --- | --- |
| `T_CLAIM_POLICY` | **27 dari 32 kolomnya tersedia di `T_PREMIUM_LIST`** (32 = 23 kolom dasar + 9 kolom tambahan hasil audit; `ID`/`CLAIM_ID` tidak dihitung). Lima sisanya kehilangan rumah — §Blocker. |
| `T_CLAIM_MARKETING` | **Terbukti turunan, nol nilai asli milik klaim.** `Claim Life/Activity/SetMOClaim_Act.xml` mengisi `MarketingData` dari `PolicyDataLife` (`MOID` · `MarketingCode` · `MarketingName` · `TeamGroup` · `BranchCode` · `BranchName`) **atau** dari master marketing officer (`ASM-FW-GISFW-Int-marketingofficer`, lewat `Claim Life/ReportDefinition/BrowseMarketingOfficer_RD.xml`). |

⚠️ **Penyimpangan sadar — potret berubah menjadi baca hidup.** `[keputusan work owner]` Keputusan
lama menyebut `PolicyDataLife` sebagai **snapshot** — *"disalin, bukan baca live"* — dan itu **kini
dibalik**. Akibatnya nyata dan diterima: **polis yang berubah sesudah klaim dibuat akan mengubah
tampilan klaim lama**; untuk polis ber-endorsement itu **pasti terjadi**. Migrasi karena itu
**tidak** menyalin atribut polis ke mana pun.

✅ **Spec §2b selaras** — `.scratch/claim-life/spec.md` §2b memuat RALAT 2026-09-18 A–E dengan pohon
dan tabel relasi yang sama. Teks lama di spec **tidak dihapus**, hanya diralat.

⚠️ **Dua tabel spreading tetap ada** — keduanya **terlewat di audit awal** dan ditemukan pada
verifikasi 2026-09-16. Bukti di bawah.

##### Dua tabel spreading — `[terverifikasi]` induk = baris adjustment

`[terverifikasi]` `AdjustmentList` (class `ASM-FW-GISFW-Data-AdjustmentLife`) punya anak
`SpreadingList` → `RetroLifeList`. Bukti:

| Rule | Class / Nama / Tipe | Path | Isi |
| --- | --- | --- | --- |
| `SpreadingClaimLife_Act` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `SPREADINGCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/SpreadingClaimLife_Act.xml` | mengisi `.SpreadingList` **dan** `.RetroLifeList` beserta perhitungannya |
| `AdjustmentDetail_Section` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-HTML-SECTION` | `Claim Life/Section/AdjustmentDetail_Section.xml` | grid `.SpreadingList` |
| `RetroDetailClaimLife` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `RETRODETAILCLAIMLIFE` / `RULE-HTML-SECTION` | `Claim Life/Section/RetroDetailClaimLife.xml` | menampilkan `.REINSURERNAME`, `.PERCENTSHARE`, `.Amount` |

`[terverifikasi]` Kelas baris, dari deklarasi halaman di `SpreadingClaimLife_Act`:

```
.SpreadingList   = OutwardList.pxResults   → ASM-FW-GISFW-Int-TREATYYEAR_LIFE      (per treaty-year)
.RetroLifeList   = RetroLife.pxResults     → ASM-FW-GISFW-Int-RETROCESSIONLIFE     (per reinsurer)
```

**`T_CLAIMLF_ADJUSTMENT_SPREADING`** (per treaty-year) — `ID` (PK), **`ADJUSTMENT_ID`** (FK →
`T_CLAIMLF_ADJUSTMENT.ID`), `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`, `TREATY_YEAR_LIFE`,
`RETROCADED_SHARE` **NUMBER**, `RATE` **NUMBER**, `IDR` **NUMBER**, `USD` **NUMBER**, `CURRENCY`.

**`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`** (per reinsurer) — `ID` (PK), **`SPREADING_ID`** (FK →
`T_CLAIMLF_ADJUSTMENT_SPREADING.ID`), `REINSURER_NAME`, `PERCENT_SHARE` **NUMBER**, `AMOUNT`
**NUMBER**, `RATE` **NUMBER**, `PREMIUM_SPREADED_GROSS` **NUMBER**, `PREMIUM_SPREADED_NET`
**NUMBER**, `COMMISION` (sic) **NUMBER**, `OVR_COMM` **NUMBER**, `TREATY_TYPE_ID`,
`TREATY_TYPE_NAME`.

`[terverifikasi]` Perhitungan yang membekukan nilainya, dari `SpreadingClaimLife_Act`:

```
.RetrocadedShare        = local.ClaimNet   (atau .IDR / .USD menurut mata uang)
.Amount                 = local.AmountRetroShare * @divide(.PERCENTSHARE, 100, 4)
.RATE                   = @divide(local.Rate, 1, 4)          ← local.Rate sudah dibagi 1000 (per-mil)
.PREMIUM_SPREADED_GROSS = local.Rate * (1 + local.EMPercent) * .Amount
.PREMIUM_SPREADED_NET   = .PREMIUM_SPREADED_GROSS - local.Comm
.PREMIUM_SPREADED_NET   = .PREMIUM_SPREADED_GROSS - local.Discount - local.Comm      ← cabang kedua
```

⚠️ **`PREMIUM_SPREADED_NET` punya DUA rumus** `[terverifikasi]` — satu mengurangi komisi saja, satu
mengurangi diskon **dan** komisi. Keduanya ada di rule yang sama. Cabang mana yang berlaku
**tidak terbaca dari korpus** → **`[terbuka]`**, dicatat di §Blocker. **Jangan tebak.**

⚠️ **Penyimpangan sadar (baru) — spreading adjustment DIBEKUKAN dan DISIMPAN.**
`[keputusan work owner]` Nilai spreading pada saat adjustment dibuat tersimpan apa adanya; perubahan
master treaty sesudahnya **tidak** mengubahnya. Konsisten dengan spreading di PremiumList Life.

##### `KOMITE_ID` pada `T_CLAIMLF_ADJUSTMENT`

⚠️ **Penyimpangan sadar (baru) — roster & keputusan komite TIDAK disimpan di Claim Life.**
`[keputusan work owner]` `T_CLAIMLF_ADJUSTMENT` menyimpan **satu kolom rujukan**, `KOMITE_ID`
(**nullable** — `NULL` bila baris belum pernah dikirim ke Komite). Roster dan keputusan per anggota
milik konteks **Komite Claim Life**; Claim Life **melihat**nya lewat join. **Tidak ada
`T_CLAIMLF_ADJUSTMENT_KOMITE`.**

⚠️ **`KOMITE_ID` berisi identitas kasus komite — yaitu `T_WORK_CLAIM.ID` baris komite.**
`[keputusan work owner]` REVISI 2026-09-17. Ia **tidak** menunjuk `T_GENERAL_KOMITE.ID`. Kolomnya
**ber-index**. Alasan work owner: saat user melihat baris adjustment, ia langsung tahu baris itu
merujuk kasus komite yang mana, tanpa query terbalik.

✅ **Tipe `KOMITE_ID` — DIPUTUSKAN 2026-09-18, `[terbuka]` tipe DITUTUP.**
`[keputusan work owner]` Ia mengikuti `T_WORK_CLAIM.ID`, yang kini **teks berformat**: baris klaim
`CLM-xxxxxx` (contoh `CLM-123456`), baris komite `KMT-xxxxxx` (contoh `KMT-000789`). **Bukan angka
sequence.** `KOMITE_ID` karena itu berisi teks berformat `KMT-xxxxxx`.

⚠️ **Penyimpangan sadar dari ADR-U-0006 — dicatat, bukan dilanggar diam-diam.**
`[keputusan work owner]` **ADR-U-0006** menetapkan identitas dari **sequence**. `T_WORK_CLAIM` dan
kedua tabel ber-shared-PK memakai **nomor bisnis berformat** alih-alih sequence murni. ADR-U-0006
tetap berlaku untuk identitas tabel klaim lainnya (`T_CLAIMLF_*`, `DOCUMENT_CLAIM`).

⚠️ `[terbuka]` **Dua hal TETAP terbuka. Jangan tebak:** (a) **generator** nomor `CLM-`/`KMT-` —
siapa yang membuatnya, apakah ada sequence di belakang prefiks, apakah di-reset per tahun —
pemilik **DBA / work owner**; (b) apakah `KOMITE_ID` dan `COVER_KEY` dipasangi
`REFERENCES T_WORK_CLAIM(ID)` atau dibiarkan tanpa constraint — pemilik **DBA / work owner**.

Rantai penuh untuk membaca keputusan komite dari sebuah baris adjustment:

```
T_CLAIMLF_ADJUSTMENT.KOMITE_ID
  → T_WORK_CLAIM   (ID = KOMITE_ID ; COVER_KEY = ID baris klaim)
  → T_GENERAL_KOMITE  (ID = T_WORK_CLAIM.ID  <- SHARED PK ; ADJUSTMENT_ID = adjustment.ID)
  → T_KOMITE_KOMITELIST (keputusan per anggota, diurut KOMITE_URUT)
```

⚠️ **`T_WORK_CLAIM` TIDAK memuat `ADJUSTMENT_ID`** `[keputusan work owner]` 2026-09-17 — penutup
lingkar ke baris adjustment ada di **`T_GENERAL_KOMITE.ADJUSTMENT_ID`**. Kolom **`COVER_KEY`** pada
`T_WORK_CLAIM` **ditambahkan oleh tiket 00 Komite Claim Life** lewat `ALTER`, bukan oleh tiket ini —
tiket ini hanya membuat tabelnya.

⚠️ **REVISI 2026-09-17** `[keputusan work owner]` — **`KMT_NO` DIBUANG.** `T_WORK_CLAIM` adalah tabel
**satu baris per work object**: saat kirim komite, kasus komite mendapat **barisnya sendiri**
(`ID` = identitas kasus komite) dan menunjuk klaimnya lewat **`COVER_KEY`**. Tidak ada kolom nomor
terpisah. Test yang menemukan `KMT_NO` **gagal**.

###### ⚠️ Catatan korpus — bentuk penautan ini TIDAK ADA di Pega

`[keputusan work owner]` `COVER_KEY` dan **shared primary key** adalah **bentuk baru**, bukan
temuan korpus.
Jangan mencarinya di Pega. Begitu pula `KMT_NO` yang sempat dirancang lalu dibuang.

Sensus 2026-09-17 atas **81 berkas korpus** — seluruh modul `Komite Claim Life` (47 dari 53 berkas;
sisanya `.xlsx`) dan 34 berkas `Claim Life` (**seluruh** `RDBList/`, `Section/ClaimComite.xml`,
`Harness/Committe_Life.xml`, dan aktivitas kunci) menghasilkan:

| Yang dicari | Hasil |
| --- | --- |
| string `KMT-` | **0 berkas** |
| string `KMT_NO` | **0 berkas** — sempat dirancang, dibuang 2026-09-17 |
| `KMT` lainnya | hanya di **nama rule** — `CreateKMTLife_Act`, `Generate_NoAccept_KMT_Life`, `HitServiceToKasirKMTLife_Act` |

`[terverifikasi]` **Penaut nyata di Pega bukan nomor kasus.** `Claim Life/Activity/CreateKMTLife_Act.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) menautkan kasus
komite ke klaim lewat **`pxAddChildWork`** (3×, mekanisme parent-child bawaan Pega), `pzInsKey` (7×),
`pyID` (2×), ditambah muatan `childPageKomite.CLMNO`, `.IndexAdjustment`, `.IndexPremiumList`. **Nol**
nilai literal ber-`KMT`.

Jadi rantai `KOMITE_ID` → `T_WORK_CLAIM` (`ID`/`COVER_KEY`) → `T_GENERAL_KOMITE` → `T_KOMITE_KOMITELIST`
adalah **rancangan pengganti** mekanisme parent-child Pega — sah sebagai keputusan desain, tetapi
**bukan** `[terverifikasi]`. Siapa pun yang menuliskannya ulang wajib memakai tanda
`[keputusan work owner]`, bukan `[terverifikasi]`.

Catatan: `COVER_KEY` **lebih dekat** ke bentuk Pega daripada rancangan `KMT_NO` sebelumnya, karena
Pega memang menautkan induk-anak lewat mekanisme *cover* (`pxAddChildWork`, `pzInsKey`) dan bukan
lewat kolom nomor. Tetapi nama `COVER_KEY` sendiri tetap tidak muncul di korpus.

`[terbuka]` **Nasib `CLMNO`.** `[terverifikasi]` `CLMNO` adalah nomor klaim yang **benar-benar**
dibawa ke kasus komite di Pega — muncul di `Claim Life/Activity/CreateKMTLife_Act.xml` dan
`Komite Claim Life/Activity/SendEmailKlaimLife.xml` (2 berkas korpus), dan tercatat di tiket 10
baris 33 serta tiket 01 Komite sebagai bagian muatan penyerahan. Tetapi ia **tidak muncul sama
sekali** dalam rancangan rantai relasional di atas. Apakah `CLMNO` tetap dibawa (misalnya sebagai
kolom pada `T_GENERAL_KOMITE`), atau memang digantikan mekanisme `ID`/`COVER_KEY`, **belum pernah
dinyatakan**. **Jangan tebak** — keputusan work owner.

`[terverifikasi]` Di Pega justru sebaliknya — `KomiteList` adalah **anak `AdjustmentList`**:
`Komite Claim Life/Activity/KomitePostAdjustment.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` /
`KOMITEPOSTADJUSTMENT`) menulis `…AdjustmentList(idx).KomiteList(k).KomiteAproval` /
`.KomiteComment` / `.DateApprove`; `Claim Life/Activity/CreateKMTLife_Act.xml` mengisi
`childPageKomite.KomiteList`, `IndexPremiumList`, dan `IndexAdjustment = .pxListSubscript`.

⚠️ **Penyimpangan sadar (baru) — rujukan memakai ID stabil, bukan indeks posisi.**
`[keputusan work owner]` Pega merujuk baris lewat **subscript posisi** (`IndexPremiumList`,
`IndexAdjustment`) — rusak begitu urutan bergeser. Sistem baru memakai **`ADJUSTMENT_ID`** dan
**`KOMITE_ID`**. Kelas bug yang sama sudah diperbaiki di Endorsement Life dengan `PARENT_ID`.

Daftar kolom lengkap tiap tabel ada di **spec §2b** dan
`.scratch/claim-life/revisi-penyimpanan-json-dibuang.md`.

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InsertJsonKlaimLife_sql` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!INSERTJSONKLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/InsertJsonKlaimLife_sql.xml` | ⚠️ **INSERT flat 51 kolom** ke `OS_AKSEPTASI_KLAIM_LIFE` — **bukan** JSON; sumber migrasi |
| `UpdateOsAkseptasiClaimLife_sql` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | 55 nama kolom terbaca langsung |
| `SavePesertaClaim` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEPESERTACLAIM` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/SavePesertaClaim.xml` | ⚠️ **bukti induk adjustment = peserta** |
| `SaveOutStandingLife_Act` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/SaveOutStandingLife_Act.xml` | ⚠️ bukti dokumen **per peserta** |

`[terverifikasi]` `OS_AKSEPTASI_KLAIM_LIFE` menyimpan **satu baris per peserta** dan **mengulang 14
atribut polis** pada setiap baris peserta.

#### ADR terkait

**ADR-U-0003** (uang non-float), **ADR-U-0006** (identitas lewat sequence), **ADR-U-0009** (migrasi
penuh; koeksistensi ditolak), **ADR-U-0011** (seluruh **baris** adjustment ikut pindah).

#### Acceptance criteria

- [ ] ⚠️ Skema klaim **relasional penuh**: setiap atribut menjadi **kolom bernama**. Test yang
      menemukan kolom JSON menyimpan atribut klaim **gagal**. *(AC 31 spec; penyimpangan sadar 1)*
- [ ] ⚠️ **Enam tabel klaim** ada dengan PK dan FK sesuai diagram — `T_GENERAL_CLAIM`,
      `T_CLAIMLF_PREMIUMLIST_DETAIL`, `T_CLAIMLF_ADJUSTMENT`, `T_CLAIMLF_ADJUSTMENT_SPREADING`,
      `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`, `DOCUMENT_CLAIM`. **REVISI 2026-09-18:** ⛔ **tidak
      ada** `T_CLAIMLF_POLICY` maupun `T_CLAIMLF_MARKETING` — keduanya **dihapus**, bukan diganti
      nama. Test yang menemukan salah satunya **gagal**. *(AC 48 spec; penyimpangan sadar 8)*
- [ ] ⚠️ **REVISI 2026-09-18** — **tidak** seluruh FK `ON DELETE CASCADE`. Yang **CASCADE**: relasi
      **3 · 4 · 5 · 6** (klaim → peserta → adjustment → spreading → spreading retro). Yang **di Go**:
      relasi **1** (`COVER_KEY`) dan **11** (`KOMITE_ID`, penunjuk). Yang ⚠️ `[terbuka]`: relasi
      **7** (`DOCUMENT_CLAIM`, kini lintas-lini — induk beda tabel per lini) dan relasi **2**
      (belum punya kolom). Lihat tabel sebelas relasi di atas. **Jangan menyeragamkan sendiri.**
- [ ] ⚠️ FK `T_CLAIMLF_ADJUSTMENT` menunjuk **`T_CLAIMLF_PREMIUMLIST_DETAIL.ID`**, **bukan**
      `T_GENERAL_CLAIM.ID`. Test yang menemukan adjustment menggantung pada header **gagal**.
      *(AC 33 spec; penyimpangan sadar 2)*
- [ ] ⚠️ FK `DOCUMENT_CLAIM` menunjuk **`T_CLAIMLF_PREMIUMLIST_DETAIL.ID`**. *(AC 44 spec;
      penyimpangan sadar 5)*
- [ ] ⚠️ **`T_WORK_CLAIM` ada sebagai tabel mandiri** — bukan anak `T_GENERAL_CLAIM`, dan keadaan tangga
      **tidak** menjadi kolom header klaim. *(AC 46 spec; penyimpangan sadar 6)*
- [ ] Header memuat keempat field `PremiumListSummary` (`CLAIM_NO`, `PL_NUMBER`, `RISLIPRNM`,
      `BUSINESS_NAME`) beserta `CASEID` dan `CLAIM_RETRO`. *(AC 36 spec)*
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS, MENUNGGU JAWABAN.**
      `T_CLAIM_POLICY` **dihapus**, jadi tidak ada tabel yang memuat kesembilan kolom itu. **27 dari
      32** kolomnya dibaca dari `T_PREMIUM_LIST`; **lima** ⚠️ `[terbuka]` **kehilangan rumah** —
      `TEAM_GROUP`, `BUSINESS_ID`, `TANGGAL_RESPON`, `TANGGAL_REALISASI`,
      `TANGGAL_KONFIRMASI_BALIK` (§Blocker). `PL_NUMBER` tetap tidak ada, tetapi karena **diganti
      nama menjadi `POLICY_NO`** di `T_GENERAL_CLAIM` — bukan karena tabel ini menolaknya.
      **Jangan tebak rumah kelima kolom itu.** *(AC 37 spec)*
- [ ] `T_CLAIMLF_PREMIUMLIST_DETAIL` memuat **kesembilan kolom tambahan** hasil audit — termasuk
      `IS_CHECK` dan ketiga tanggal per peserta. *(AC 39, 40, 41 spec)*
- [ ] `T_CLAIMLF_ADJUSTMENT` memuat `CLAIM_AMOUNT`, `STS_REJECT`, `ACCEPTEDNO`, `ACCEPTATION_DATE`.
      *(AC 42, 43 spec)*
- [ ] `T_CLAIMLF_ADJUSTMENT` memuat **ketiga kolom bank** — `NAME_OF_BANK`, `ID_BANK`, `ACCOUNT_NO` —
      dan migrasi mengisinya dari kolom warisan `NAME_OF_BANK`, `IDBANK`, `ACCOUNTNO`.
      *(AC 56 spec; `[terverifikasi]`)*
- [ ] ⚠️ Seluruh uang dan share bertipe **desimal presisi arbitrer**; seluruh tanggal **`DATE`**;
      seluruh kolom **nullable**; identitas dari **sequence**. Test yang menemukan kolom uang
      bertipe teks atau melewati `float` **gagal**. *(AC 50 spec; **ADR-U-0003**, **ADR-U-0006**;
      penyimpangan sadar 7)*
- [ ] Seluruh klaim Life terbawa **beserta seluruh baris adjustment**-nya — bukan hanya keadaan
      terakhir. Jumlah baris per klaim setelah migrasi **sama** dengan sebelumnya. *(**ADR-U-0011**)*
- [ ] Setiap baris adjustment hasil migrasi **menunjuk peserta yang benar**. *(AC 34 spec)*
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS.** Tidak ada
      `T_CLAIM_POLICY` untuk dinormalkan ke dalamnya; atribut polis **tidak dipindahkan ke mana
      pun** — ia **dibaca hidup** dari tabel polis ⚠️ **penyimpangan sadar**. Yang **tetap
      mengikat**: bila atribut polis **berbeda antar baris peserta** dalam satu klaim, migrasi
      **melaporkannya** dan **tidak** diam-diam memilih salah satu. *(AC 51 spec)*
- [ ] ⚠️ Migrasi **melaporkan** bahwa data lama ber-`ACCEPTATION_DATE` = waktu insert dan
      `STS_REJECT` = `0` karena **di-hardcode** di sumbernya, bukan karena nilainya sebenarnya.
      Migrasi **tidak mengarang** tanggal akseptasi. *(AC 52 spec; penyimpangan sadar 4)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
      tepat**, bukan dengan toleransi. *(AC 53 spec; **ADR-U-0003**)*
- [ ] Tanggal yang berupa teks menjadi `DATE` **tanpa pergeseran zona waktu**; yang **tidak dapat
      diurai dilaporkan**, bukan didiamkan. *(AC 54 spec)*
- [ ] ⚠️ Migrasi **tidak mengambil apa pun** dari `m_product_life.JSONDATA`; `PRODUCT_NAME` dan
      `PRODUCT_NAME_ID` berasal dari **`product_life` relasional**. *(AC 38 spec)*
- [ ] Penomoran klaim **tidak melompat dan tidak mengulang** setelah migrasi; sequence pindah dengan
      **nilai berjalan yang benar**. *(**ADR-U-0006**)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

##### Kolom yang bertambah dan yang pindah ⚠️ BARU 2026-09-18

- [ ] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `CREATE_OP`, `CREATE_OP_NAME`, `TGL_UPDATE` —
      ketiganya **pindah ke `T_WORK_CLAIM`**. Test yang menemukannya di header klaim **gagal**.
- [ ] ⚠️ `T_WORK_CLAIM` memuat `CREATE_OP`, `CREATE_OP_NAME`, `TGL_UPDATE`, **`CASEID`**
      (seluruhnya pindahan) dan `LINI`.
- [ ] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `CASEID` — pindah ke `T_WORK_CLAIM`. Test yang
      menemukannya di header klaim **gagal**. *(DIPUTUSKAN 2026-09-18; `[keputusan work owner]`)*
- [ ] ⚠️ `T_GENERAL_CLAIM` **tidak lagi memuat** `PL_NUMBER`; kolom itu **diganti nama** menjadi
      **`POLICY_NO`**. Test yang menemukan `PL_NUMBER` di header klaim **gagal**.
- [ ] ⚠️ `T_GENERAL_CLAIM` memuat ketiga **penunjuk polis** — `CASEID_POLICY`, `POLICY_NO`,
      `ENDORSMENT_NO` — dan atribut polis **tidak** disalin ke klaim. Test yang menemukan salinan
      atribut polis di tabel klaim **gagal**.
- [ ] ⚠️ `[terbuka]` **Isi ketiga penunjuk polis belum dapat ditulis** — tidak satu pun menunjuk
      kolom yang ada di sisi polis (§Blocker). Kolomnya dibuat; **pengisiannya** menunggu keputusan
      work owner. **Jangan tebak.**
- [ ] ⚠️ Kolom `LINI` **ADA** pada `T_WORK_CLAIM`, dan untuk Life isinya **konstanta lini Life**.
      `[terbuka — Non-Life]` hanya **daftar nilai enum lintas-lini**, yang ditetapkan saat konteks
      Non-Life digarap — **bukan** keberadaan kolomnya. *(DIPUTUSKAN 2026-09-18)*

##### Identitas: shared PK dan nomor bisnis berformat ⚠️ BARU 2026-09-18

- [ ] ⚠️ **`T_GENERAL_CLAIM.ID` sama persis dengan `T_WORK_CLAIM.ID` baris klaim** — **shared
      primary key**, 1:1, **tanpa kolom penyambung**. `T_GENERAL_CLAIM.ID` sekaligus PK **dan** FK
      ke `T_WORK_CLAIM.ID`. *(relasi 2; `[keputusan work owner]`)*
- [ ] ⚠️ ⛔ **Tidak ada kolom `WORK_CLAIM_ID`** di mana pun — pada `T_GENERAL_KOMITE` maupun
      tabel lain. Hubungan ke kasus komite juga **shared PK**: `T_GENERAL_KOMITE.ID` =
      `T_WORK_CLAIM.ID` baris komite. Test yang menemukan kolom `WORK_CLAIM_ID` **gagal**.
      *(relasi 8; `[keputusan work owner]`)*
- [ ] ⚠️ **`T_GENERAL_KOMITE.ADJUSTMENT_ID` TETAP ADA** — penutup lingkar ke baris adjustment; ia
      **bukan** bagian shared PK dan **tidak** ikut dibuang. *(relasi 10)*
- [ ] ⚠️ **`T_WORK_CLAIM.ID` bertipe teks berformat** — baris klaim `CLM-xxxxxx`, baris komite
      `KMT-xxxxxx`. **Bukan angka sequence.** Test yang menemukan tipe numerik **gagal**.
      *(`[keputusan work owner]`)*
- [ ] ⚠️ `T_WORK_CLAIM.COVER_KEY`, `T_GENERAL_CLAIM.ID`, `T_GENERAL_KOMITE.ID`, dan
      `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` **bertipe sama** dengan `T_WORK_CLAIM.ID` — teks berformat.
- [ ] ⚠️ **Penyimpangan sadar dari ADR-U-0006 dicatat di artefak**, tidak dilanggar diam-diam:
      identitas `T_WORK_CLAIM` dan kedua tabel ber-shared-PK adalah **nomor bisnis berformat**,
      bukan sequence. ADR-U-0006 **tetap berlaku** untuk `T_CLAIMLF_*` dan `DOCUMENT_CLAIM`.
- [ ] ⚠️ `[terbuka]` **Generator nomor `CLM-`/`KMT-` belum ditetapkan** — siapa yang membuatnya,
      apakah ada sequence di belakang prefiks, apakah di-reset per tahun. Pemilik **DBA / work
      owner**. Kolomnya dibuat; **pembangkitannya** menunggu jawaban. **Jangan tebak.**

##### Spreading adjustment + rujukan Komite ⚠️ BARU 2026-09-16

- [ ] ⚠️ **`T_CLAIMLF_ADJUSTMENT_SPREADING` ada**, dengan FK **`ADJUSTMENT_ID`** → `T_CLAIMLF_ADJUSTMENT.ID`
      dan **`ON DELETE CASCADE`**. Test yang menemukannya menggantung pada peserta atau pada header
      klaim **gagal**. *(`[terverifikasi]` `SpreadingClaimLife_Act` mengisi `.SpreadingList` pada
      baris adjustment; **AC 58 spec**)*
- [ ] ⚠️ **`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` ada**, dengan FK **`SPREADING_ID`** →
      `T_CLAIMLF_ADJUSTMENT_SPREADING.ID` dan **`ON DELETE CASCADE`**. *(AC 58 spec)*
- [ ] Pohon klaim berkedalaman **lima tingkat** — klaim → peserta → adjustment → spreading →
      spreading retro — dan menghapus klaim **mengkaskade sampai tingkat terdalam**. Test wajib
      memeriksa **cicit** (`_SPREADING_RETRO`) ikut hilang. *(AC 60 spec)*
- [ ] `T_CLAIMLF_ADJUSTMENT_SPREADING` memuat `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`,
      `TREATY_YEAR_LIFE`, `RETROCADED_SHARE`, `RATE`, `IDR`, `USD`, `CURRENCY`.
- [ ] `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` memuat `REINSURER_NAME`, `PERCENT_SHARE`, `AMOUNT`,
      `RATE`, `PREMIUM_SPREADED_GROSS`, `PREMIUM_SPREADED_NET`, `COMMISION`, `OVR_COMM`,
      `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`. *(`[terverifikasi]`
      `Claim Life/Section/RetroDetailClaimLife.xml` + `SpreadingClaimLife_Act`)*
- [ ] ⚠️ Seluruh kolom uang dan persen pada kedua tabel bertipe **desimal presisi arbitrer**;
      **tidak** melewati `float`. *(**ADR-U-0003**)*
- [ ] ⚠️ Nilai spreading **dibekukan**: perubahan master treaty setelah adjustment tersimpan
      **tidak mengubah** angka yang sudah ada. *(penyimpangan sadar — spreading disimpan)*
- [ ] ⚠️ **`T_CLAIMLF_ADJUSTMENT` memuat kolom `KOMITE_ID`**, **nullable** dan **ber-index** — `NULL`
      bila baris belum pernah dikirim ke Komite. *(AC 61 spec; penyimpangan sadar — rujukan, bukan salinan)*
- [ ] ⚠️ **Tidak ada tabel `T_CLAIMLF_ADJUSTMENT_KOMITE`.** Roster dan keputusan per anggota **tidak**
      disimpan di Claim Life. Test yang menemukan tabel itu **gagal**. *(AC 61 spec; penyimpangan sadar)*
- [ ] ⚠️ Rujukan ke Komite memakai **`KOMITE_ID`**, bukan **indeks posisi**. Test yang menemukan
      padanan `IndexPremiumList`/`IndexAdjustment` sebagai kunci rujukan **gagal**.
      *(AC 62 spec; `[terverifikasi]` `CreateKMTLife_Act` memakai `.pxListSubscript`; penyimpangan sadar)*
- [ ] Migrasi **membongkar** `SpreadingList` dan `RetroLifeList` dari data lama ke kedua tabel,
      dan setiap baris dapat ditelusuri ke **baris adjustment yang benar**.
- [ ] Setiap FK baru (`ADJUSTMENT_ID`, `SPREADING_ID`) **ber-index**.
- [ ] ⚠️ `KOMITE_ID` **ber-index** juga, dan berisi **identitas kasus komite** = `T_WORK_CLAIM.ID`
      baris komite. Test yang menemukannya menunjuk `T_GENERAL_KOMITE.ID` **gagal**.
      *(§`KOMITE_ID` pada `T_CLAIMLF_ADJUSTMENT`; `[keputusan work owner]` REVISI 2026-09-17)*
- [ ] ✅ Tipe `KOMITE_ID` **SUDAH DITETAPKAN 2026-09-18** — teks berformat `KMT-xxxxxx`, mengikuti
      `T_WORK_CLAIM.ID`. *(`[keputusan work owner]`)*
- [ ] ⚠️ `[terbuka]` **Apakah `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)`**
      atau dibiarkan tanpa constraint — **belum diputuskan**. Pemilik **DBA / work owner**. Tiket
      ini **tidak dinyatakan selesai** sebelum jawabannya ada. **Jangan tebak.**

##### Koreksi: `OS_AKSEPTASI_KLAIM_LIFE` TETAP ditulis ⚠️ 2026-09-16

⚠️ `[keputusan work owner]` **Keputusan sebelumnya dibalik.** Sistem baru menulis **dua tempat**:
tabel relasional baru **dan** `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE`, karena hilir (Arasapas /
produksi) masih membaca dari sana. **Yang dibuang hanya JSON** (`JSON_KLAIM` / serialisasi
`ClaimData`).

- [ ] ⚠️ Menyimpan klaim menulis **tabel relasional baru** *dan* `INSERT` flat ke
      `OS_AKSEPTASI_KLAIM_LIFE`. Test yang menuntut `OS_AKSEPTASI_KLAIM_LIFE` **tidak** ditulis
      adalah **keliru** dan harus dibalik.
- [ ] ⚠️ Yang **tetap dibuang** hanya **blob JSON** — `JSON_KLAIM` dan serialisasi `ClaimData`.
      *(AC 31, 34 spec tetap berlaku untuk JSON saja)*
- [ ] Kedua penulisan berada dalam **satu transaksi**; kegagalan pada salah satunya **membatalkan
      keduanya**. *(AC 49 spec)*

#### Blocker

✅ **PEMBLOKIR RELASI 2 DICABUT 2026-09-18.** `[keputusan work owner]` Relasi
`T_WORK_CLAIM` → `T_GENERAL_CLAIM` **dipecahkan lewat SHARED PRIMARY KEY**: `T_GENERAL_CLAIM.ID`
**sama persis** dengan `T_WORK_CLAIM.ID` baris klaim. Tidak perlu kolom penyambung, dan
**`WORK_CLAIM_ID` dibuang** dari seluruh rancangan. Relasi 8 dipecahkan dengan cara yang sama.

> Teks lama: *"⛔ PEMBLOKIR — relasi nomor 2 tidak punya kolom di sisi mana pun … Tiket ini tidak
> dinyatakan selesai sebelum dijawab."* — **dicabut.** Relasi nomor 2 **bukan lagi** alasan tiket
> ini belum selesai.

⚠️ **Tiket ini MASIH belum dapat dinyatakan selesai**, tetapi karena **satu** hal lain:
⚠️ `[terbuka]` **apakah `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES T_WORK_CLAIM(ID)`**
atau dibiarkan tanpa constraint — pemilik **DBA / work owner**. Tipe kolomnya **sudah** ditetapkan
(teks berformat); yang belum hanya **constraint**-nya.

(Yang tetap benar: **OQ-001 dan OQ-002 (bagian penyimpanan) ditutup 2026-09-16** — tabel klaim
dirancang sendiri, dan jalur tulis klaim tidak lagi lewat stored procedure.)

##### ⚠️ `[terbuka]` — lima kolom kehilangan rumah

**Status pemblokirnya BELUM ditetapkan work owner.** Keputusan 2026-09-18 menyatakan kelima kolom ini
`[terbuka]` dan **tidak** menyatakannya memblokir. Ini 5 dari 32 kolom `T_CLAIM_POLICY` yang
**tidak** tersedia di `T_PREMIUM_LIST`. **JANGAN DITEBAK:**

| Kolom | Keadaan |
| --- | --- |
| `TEAM_GROUP` | Tidak ada di `T_PREMIUM_LIST`; **ada di master marketing officer** — diambil lewat `MO_ID` |
| `BUSINESS_ID` | Tidak ada di `T_PREMIUM_LIST`; dipakai sebagai **parameter** `Claim Life/RDBList/Generate_NoAccept_Life.xml`. **Sumbernya wajib ditetapkan** |
| `TANGGAL_RESPON` · `TANGGAL_REALISASI` · `TANGGAL_KONFIRMASI_BALIK` | **Tanggal proses klaim, bukan atribut polis.** Tidak ada di `T_PREMIUM_LIST`, dan **tidak lagi ada di klaim**. Rumahnya `T_GENERAL_CLAIM` atau `T_WORK_CLAIM` — **belum diputuskan** |

⚠️ **Usulan — bukan keputusan.** Ketiga `TANGGAL_*` bukan perkara kosmetik: `[terverifikasi]`
ketiganya **tampil di empat section aktif** (`InputRegisterClaimLife`, `InputOSClaimLife`,
`InputAkseptasiClaimLife`, `MedicalCheckClaimLife`) **tanpa gerbang visibilitas** — layar akan
kehilangan isi bila rumahnya tidak ditetapkan. Atas dasar itu **saya usulkan** butir ini ikut
memblokir tiket 14, tetapi **penetapannya keputusan work owner** dan **belum diambil**.

##### ⚠️ `[terbuka]` tambahan — tidak memblokir pembuatan tabel, tetapi wajib dijawab

- **Isi ketiga penunjuk polis** (`CASEID_POLICY`, `POLICY_NO`, `ENDORSMENT_NO`) — tidak satu pun
  menunjuk kolom yang ada di sisi polis. **TETAP terbuka.**
- **Generator nomor `CLM-`/`KMT-`** — siapa yang membuatnya, sequence di belakang prefiks atau
  tidak, reset per tahun atau tidak. Pemilik **DBA / work owner**.

✅ **Ditutup 2026-09-18, tidak lagi terbuka:** relasi 2 tanpa kolom (shared PK) · relasi 8
memakai `WORK_CLAIM_ID` (shared PK, kolom dibuang) · nasib `CASEID` (pindah ke `T_WORK_CLAIM`) ·
`ON DELETE` `DOCUMENT_CLAIM` (di Go) · tipe `T_WORK_CLAIM.ID`/`KOMITE_ID` (teks berformat) ·
keberadaan kolom `LINI` (ada; hanya daftar enum lintas-lini yang menunggu Non-Life).

⚠️ `[terbuka — Non-Life]` **tidak memblokir:** relasi & cascade formal `T_WORK_CLAIM` ditetapkan
saat konteks Non-Life digarap. Yang mengikat di sini hanya: klaim life dihapus → baris work-nya ikut.

⚠️ `[terbuka]` **tidak memblokir skema, memblokir perilaku hitung:** **`PREMIUM_SPREADED_NET` punya
dua rumus** di rule yang sama — `GROSS − Comm` dan `GROSS − Discount − Comm`
(`Claim Life/Activity/SpreadingClaimLife_Act.xml`). **Cabang mana yang berlaku tidak terbaca dari
korpus.** Kolomnya tetap dibuat di tiket ini; **rumusnya** ditetapkan Product + UW sebelum tiket
**03** dinyatakan selesai. **Jangan tebak.**

⚠️ `[terbuka — selaraskan Komite Claim Life]` **tidak memblokir tiket ini:** tabel komite di konteks
**Komite Claim Life** perlu menampung keputusan **per baris adjustment** — kolom `ADJUSTMENT_ID`
(FK → `T_CLAIMLF_ADJUSTMENT.ID`) beserta `KOMITE_APROVAL`, `KOMITE_COMMENT`, `DATE_APPROVE`, dan
roster. Apakah tabel di sana sudah punya kolom itu **belum diperiksa**. Penyelarasannya pekerjaan
konteks Komite, bukan di sini.

#### Catatan

⚠️ **Sumber migrasi kedua bernama "Update" tetapi isinya INSERT.** `[terverifikasi]`
`Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL`) berisi
**`INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` dengan 55 kolom** — **tidak ada `UPDATE` di
dalamnya**. Dicatat supaya **tidak ada waktu terbuang mencari jalur `UPDATE` yang memang tidak
ada**. Ketiga kolom bank berasal dari rule ini (**OQ-066**).

⚠️ **Nama berbohong — sepasang, saling membalik.** `[terverifikasi]` `InsertJsonKlaimLife_sql`
bernama "Json" tetapi **INSERT flat 51 kolom**; sebaliknya `GetJsonProductLife`
(`ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!GETJSONPRODUCTLIFE`) bernama "Json" dan **memang** JSON
(`json_table` atas `m_product_life.JSONDATA`). Nama tidak memberi tahu apa-apa di kedua arah —
**baca kodenya** (**OQ-066**).

⚠️ **`DOCUMENT_CLAIM` kolomnya belum final.** Rinciannya (nama berkas / kategori / kelengkapan)
diturunkan dari sensus `.DocumentList` di `SaveOutStandingLife_Act` saat tiket **03** dikerjakan.
Tiket ini menyiapkan **tabel dan kuncinya**; kolom isian ditetapkan di sana — **jangan tebak dari
nama tabel**.

⚠️ **`T_WORK_CLAIM` kolomnya belum final** — bentuk lengkapnya menyusul saat penyelarasan
lintas-lini. Untuk Claim Life yang wajib ada: posisi tangga, status akseptasi, penanda pengembalian
ke Admin dan ke Medical, serta `Type`.

#### Seam & perintah verifikasi

**Seam: API HTTP Claim — Life** terhadap **skema uji Oracle nyata** — kaskade tiga tingkat, presisi
desimal, konversi tanggal, dan normalisasi atribut polis **hanya berperilaku benar pada basis data
sungguhan**; memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Claim Life - 15 - Hapus klaim — popup konfirmasi, kaskade tiga tingkat, dan baris work

**Status:** ready-for-agent

**Blocked by:** 14 (skema relasional klaim — PREFACTOR), 03 (peserta, adjustment, dan dokumen harus
ada agar dapat dihitung dan dihapus)

> Tiket ini menutup **AC 47** — yang sebelumnya tidak dirujuk tiket mana pun — dan **sisi perilaku**
> AC 48. Tiket 14 hanya membuat *constraint* `ON DELETE CASCADE`-nya; yang dilihat pengguna —
> peringatan berisi jumlah, dan pembatalan yang benar-benar tidak mengubah apa pun — belum ada
> pemiliknya.

#### Hasil & nilai pengguna

Sebagai **`ReasLifeAdmin`**, saya dapat menghapus sebuah klaim **beserta seluruh isinya** —
peserta, seluruh baris adjustment, seluruh spreading dan spreading retro, dan seluruh dokumen —
tetapi **tidak sebelum diberi peringatan berisi jumlah baris yang akan ikut terhapus**, supaya saya
dapat membatalkan ketika angkanya tidak sesuai dugaan. *(User story 54 di spec)*

⚠️ **RALAT 2026-09-18** — kata **"polis, marketing"** **DICABUT** dari daftar di atas: kedua
tabelnya **dihapus**. Data polis dan marketing **dibaca hidup** dari tabel polis dan **bukan milik
klaim**, jadi menghapus klaim **tidak menyentuhnya**.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Hapus berkaskade tiga tingkat + baris `T_WORK_CLAIM`, **dalam satu transaksi** |
| `internal/services` | Hitung jumlah baris terdampak sebelum menghapus; orkestrasi |
| `internal/handlers` | Endpoint pratinjau dampak + endpoint hapus |
| `frontend/` | Popup konfirmasi Ya/Batal dengan rincian jumlah per jenis |

#### Bentuk yang dihapus — spec §2b

⚠️ **RALAT 2026-09-18** `[keputusan work owner]` — **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING`
DIHAPUS**, bukan diganti nama, jadi keduanya **tidak ada lagi untuk ikut terhapus**. ⛔ Jangan
membuat `T_CLAIMLF_POLICY`/`T_CLAIMLF_MARKETING`. Data polis dan marketing **dibaca hidup** dari
tabel polis — menghapus klaim **tidak boleh menyentuhnya sama sekali**.

```
T_GENERAL_CLAIM                       ← yang dihapus pengguna
  └─ T_CLAIMLF_PREMIUMLIST_DETAIL  1:N   ikut
        ├─ T_CLAIMLF_ADJUSTMENT     1:N   ikut  ⬅ CUCU
        │     └─ T_CLAIMLF_ADJUSTMENT_SPREADING        1:N  ikut  ⬅ CICIT
        │            └─ T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO  1:N  ikut  ⬅ CICIT-CUCU
        └─ DOCUMENT_CLAIM         1:N   ikut  ⬅ CUCU

T_WORK_CLAIM                      ← baris work klaim life itu ikut terhapus
```

⚠️ **RALAT 2026-09-18 — klaim "tiga tingkat" di tiket ini SUDAH SALAH SEJAK 2026-09-16, bukan
karena ralat hari ini.** Dua tabel spreading ditemukan pada verifikasi 2026-09-16 dan **tidak pernah
masuk** ke pohon di tiket ini. Kaskade sebenarnya menyentuh **lima tingkat** — klaim → peserta →
adjustment → spreading → spreading retro — persis seperti **AC 60 spec** dan tiket **14**. Uji wajib
memeriksa **sampai cicit-cucu** (`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`), bukan hanya dua tingkat
cucu. Pohon yang berlaku: `spec.md` §2b RALAT D.

#### ADR terkait

**ADR-U-0007** (jejak audit), **ADR-U-0015** (kegagalan ditangani eksplisit, tidak ditelan).

#### Acceptance criteria

- [ ] ⚠️ Menghapus klaim **mengkaskade** ke peserta, **seluruh baris adjustment**, **seluruh
      spreading**, **seluruh spreading retro**, dan **seluruh dokumen**. **REVISI 2026-09-18:** kata
      **"polis, marketing"** **DICABUT** — kedua tabelnya dihapus. Test wajib memeriksa **setiap
      tingkat** secara terpisah, **sampai cicit-cucu**. *(AC 48 spec; penyimpangan sadar 8)*
- [ ] ⚠️ Menghapus klaim **tidak menyentuh tabel polis maupun marketing** — data itu **dibaca
      hidup** dari `T_PREMIUM_LIST` dkk dan **bukan milik klaim**. Test yang menemukan penghapusan
      menyentuh tabel polis **gagal**. *(REVISI 2026-09-18; `[keputusan work owner]`)*
- [ ] ⚠️ Penghapusan **didahului popup konfirmasi Ya/Batal** yang menyebut **jumlah baris tiap
      jenis** yang akan ikut terhapus. *(AC 48 spec)*
- [ ] ⚠️ Memilih **Batal** **tidak mengubah apa pun** — tidak ada baris terhapus, tidak ada status
      berubah, tidak ada jejak audit penghapusan. *(AC 48 spec)*
- [ ] Jumlah yang ditampilkan popup **sama persis** dengan jumlah yang benar-benar terhapus —
      dihitung dari data, bukan dari perkiraan.
- [ ] ⚠️ Menghapus klaim life **menghapus juga baris `T_WORK_CLAIM`**-nya; tidak ada keadaan tangga
      yang tertinggal tanpa klaim. *(AC 47 spec; penyimpangan sadar 6)*
- [ ] Seluruh penghapusan berjalan dalam **satu transaksi**: kegagalan di tingkat mana pun
      **membatalkan seluruhnya**, dan klaim tetap utuh. *(AC 49 spec)*
- [ ] Penghapusan mencatat **jejak audit** — siapa, kapan, dan berapa baris tiap jenis.
      *(**ADR-U-0007**)*
- [ ] Penghapusan yang gagal menghasilkan kegagalan **terang-terangan**, bukan sebagian terhapus
      diam-diam. *(**ADR-U-0015**)*
- [ ] Menghapus klaim **tidak menyentuh** `M_LIFE_PREMIUM_DETAIL` — ia hanya dibaca sebagai sumber
      snapshot peserta.
- [ ] ⚠️ **Perlakuan `OS_AKSEPTASI_KLAIM_LIFE` saat klaim dihapus** ditetapkan eksplisit. Ia **tetap
      ditulis** saat klaim disimpan (koreksi 2026-09-16, AC 32 spec), sehingga menghapus klaim
      menimbulkan pertanyaan: barisnya ikut dihapus, atau ditinggal karena hilir sudah membacanya?
      `[terbuka]` — **keputusan work owner**, **jangan tebak**. Tiket ini **tidak dinyatakan selesai**
      sebelum jawabannya ada.
- [ ] Menghapus klaim **tidak menghapus** peserta di premium list sumbernya — yang terhapus hanya
      **snapshot** milik klaim itu.

#### Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka — Non-Life]` **tidak memblokir:** relasi & cascade formal `T_WORK_CLAIM` ditetapkan saat
konteks Non-Life digarap. Untuk Claim Life yang mengikat hanya AC di atas: klaim life dihapus →
baris work-nya ikut.

#### Catatan

⚠️ **Mengapa kaskade tidak cukup diuji lewat constraint basis data.** Tiket 14 memasang
`ON DELETE CASCADE`; itu menjamin *baris* hilang, bukan bahwa **angka di popup benar** dan bahwa
**Batal benar-benar membatalkan**. Keduanya perilaku layanan dan layar — dan keduanya tempat bug
biasanya muncul.

⚠️ **Pola yang diikuti.** Kaskade + popup konfirmasi sudah ditetapkan di Master Contract Retro Life
(tiket 09) dan Treaty Contract Out (tiket 10). Bedanya di sini: **tiga tingkat**, dan ada **satu
tabel di luar pohon** (`T_WORK_CLAIM`) yang ikut.

#### Seam & perintah verifikasi

**Seam: API HTTP Claim — Life** terhadap **skema uji Oracle nyata** — kaskade tiga tingkat dan
keutuhan setelah Batal **hanya berperilaku benar pada basis data sungguhan**; memalsukannya berarti
tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```
