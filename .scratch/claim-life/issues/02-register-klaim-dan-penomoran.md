# 02: Register klaim Life + penomoran

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

## Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat mendaftarkan klaim Life baru dan klaim itu langsung memperoleh
nomornya sendiri — sehingga saya tidak perlu menomori manual dan nomor tidak pernah bentrok antar
petugas. *(User story 1 dan 2 di spec)*

## Area codebase

`internal/models` (entitas klaim), `internal/repository` (pemanggilan stored procedure + penulisan
rekam klaim), `internal/services` (orkestrasi pendaftaran), `internal/handlers` (endpoint pendaftaran),
`frontend/` (formulir Register).

## Rule Pega sumber

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

## ADR terkait

**ADR-0006** (penomoran lewat stored procedure — **jangan replikasi logikanya**),
**ADR-0003** (uang non-float). Batas transaksi: **OQ-013 tertutup** — Go memegang transaksi.

## Acceptance criteria

- [ ] Klaim Life baru dapat didaftarkan lewat API dan muncul sebagai klaim berstatus awal.
- [ ] Nomor klaim **diperoleh dari `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`**, bukan dihitung di
      aplikasi.
- [ ] Aplikasi **tidak** memuat logika pembentukan format nomor apa pun. *(ADR-0006)*
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

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §2b

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

### Pencarian peserta — **hanya peserta hidup** `[keputusan work owner]`

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

## Catatan penutupan (2026-09-14)

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

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Implementasi

### Keputusan work owner 26 September 2026 — penomoran klaim

⭐ **`[DIPUTUSKAN]`** — *"jangan ada lagi pemanggilan procedure, segala procedure hardcode dalam
skrip"*.

Artinya untuk tiket ini: nomor klaim **tidak** diambil dengan memanggil
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`. Logika penomorannya ditulis di aplikasi.

⚠️ **Konsekuensi yang harus disadari sebelum sesi tiket 02 dimulai:**

1. ⛔ **Keputusan ini menyimpang dari ADR-U-0006**, yang menetapkan penomoran lewat stored procedure
   dan melarang mereplikasi logikanya. Penyimpangan sadar ini perlu dicatat di ADR — atau ADR-U-0006
   perlu dicabut/direvisi — supaya tidak terbaca sebagai pelanggaran diam-diam.
2. ⛔ **AC nomor 2 dan 3 tiket ini bertentangan dengan keputusan tersebut.** AC 2 menuntut nomor
   *"diperoleh dari `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`, bukan dihitung di aplikasi"*; AC 3
   menuntut aplikasi *"tidak memuat logika pembentukan format nomor apa pun"*. Keduanya harus
   ditulis ulang sebelum tiket ini dikerjakan; executor tidak mengubah teks AC sendiri.
3. ⚠️ **Produksi masih memakai procedure yang sama.** Selama dua sistem berjalan berdampingan,
   dua pembangkit nomor atas satu ruang nomor dapat menghasilkan **nomor bentrok**. Siapa yang
   memegang urutan selama masa itu belum ditetapkan.
4. `[data DBA]` Yang tetap diperlukan agar logikanya dapat ditulis ulang dengan benar: **sumber
   `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`** dan DDL tabel `GENERATE_SEQUENCE_NUMBER`, isi
   `POOLDATA.KODE_PRODUKSI` untuk `TYPE='LIFE'`, serta `POOLDATA.TANGGAL_CLOSING` beserta aturan
   cutover-nya. Tanpa keempatnya, menulis ulang logika penomoran berarti menebak.

**Ralat 26 September 2026 (ronde 4) atas butir 2 — AC-nya tiga, bukan dua.** **AC 10** ikut
bertentangan. Teksnya: *"Batas transaksi dipegang **Go**: commit terjadi segera setelah nomor
terbentuk, sehingga lock `SELECT … FOR UPDATE` pada `GENERATE_SEQUENCE_NUMBER` tidak menahan
pendaftar lain."* AC itu memerikan pembagian kerja procedure-dipanggil-di-dalam-transaksi-Go: Go
memegang batas transaksi, dan lock baris counter terjadi di pihak yang membentuk nomor. Begitu
keputusan **o** mencabut pemanggilan procedure, `SELECT … FOR UPDATE` harus dikirim Go sendiri,
dan kalimat AC-nya tidak lagi memerikan apa yang dibangun.

~~⚠️ `[dugaan]` Bahwa lock itu **hari ini** berada di dalam badan
`PROC_GENERATE_SEQUENCE_NUMBER` adalah kesimpulan dari ADR-U-0006 dan dari bunyi AC-nya, **bukan
hal yang sudah dilihat**: sumber procedure-nya justru butir 4 di atas, dan belum diserahkan.~~

✅ **`[terverifikasi]` 26 September 2026 sore — dugaan di atas terbukti, penandanya naik.** Sumber
procedure sudah dibaca dan disimpan di `.scratch\claim-life\SUMBER-PENOMORAN-DBA.md`. Badan
procedure memuat `SELECT no_seq … FROM POOLDATA.GENERATE_SEQUENCE_NUMBER WHERE class = :p_class AND
jenis = :p_jenis AND tahun = :v_tahun FOR UPDATE`. Jadi lock itu memang **di dalam** procedure, atas
`(CLASS, JENIS, TAHUN)`.

⭐ **Dua fakta yang ikut terbaca, dan keduanya mengubah teks AC yang harus ditulis work owner:**
**(a)** procedure **tidak memuat satu pun `COMMIT`** — batas transaksi memang sudah dipegang
pemanggil, sehingga AC 10 lebih mudah dipenuhi daripada dikira. **(b)** format lengkap
`<prefix>K<kode>.MM.YYYY.<5 digit>` **dirakit pemanggil Pega**, bukan procedure; procedure hanya
mengeluarkan `LPAD(no_seq,5,'0')` dan `MM.YYYY`. Karena itu teks pengganti **AC 3** harus mencakup
**perakitan format**, bukan hanya pembacaan counter — usulan teks di brief ronde 4 §2 o2 belum
menyebut itu.

Jadi yang harus ditulis ulang work owner sebelum tiket ini dikerjakan: **AC 2, AC 3, dan AC 10**.
Executor tidak mengubah teks AC sendiri.

Status tiket tidak diubah: **`ready-for-agent`**, dan butir 1–4 di atas adalah blocker-nya.
