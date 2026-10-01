# 02: Register klaim Life + penomoran

**Status:** sebagian — baris adjustment pertama dan INSERT datar `OS_AKSEPTASI_KLAIM_LIFE` kini terjadi saat daftar (GILIRAN-14 butir bp; bukti Oracle menunggu uji `db`), penomoran serentak belum teruji, AC procedure menunggu teks baru work owner; pembulatan peserta 7.7 ditiru sejak GILIRAN-15 (OQ-N10 ditutup)

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

- [x] Klaim Life baru dapat didaftarkan lewat API dan muncul sebagai klaim berstatus awal. — bukti: `APP_RNM/internal/handlers/register.go:daftarKlaim`, `APP_RNM/internal/services/pendaftaran.go:Pendaftaran.Daftar` (tahap awal `Outstanding Claim`); uji `TestPendaftaranMengisiTahapDanWaktuBuat`
- [ ] Nomor klaim **diperoleh dari `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`**, bukan dihitung di
      aplikasi. — belum: `[keputusan work owner o1]` ADR-U-0043 menggantikannya — procedure tidak dipanggil, `penomorCounter.NomorBerikut` menulis ulangnya di Go; teks AC menunggu work owner
- [ ] Aplikasi **tidak** memuat logika pembentukan format nomor apa pun. *(ADR-0006)* — belum: bertentangan dengan ADR-U-0043 — format dirakit di `APP_RNM/internal/services/pendaftaran.go:RakitNomorKlaim`; teks AC menunggu work owner
- [x] Dua pendaftaran berurutan menghasilkan dua nomor berbeda. — bukti: `APP_RNM/internal/repository/penomor.go:Penomor.UrutNomorBerikut` (`SELECT … FOR UPDATE`, `NO_SEQ + 1`, lalu `UPDATE`)
- [x] Rule penomoran lama **tidak** dimigrasikan: tidak ada padanan `Generate_NoKlaim_Life` maupun
      `Generate_NoKlaim_LifeRetro` di kode. `[terverifikasi]` keduanya tidak terindeks sebagai
      rujukan aktif di `SaveOutStandingLife_Act`. — bukti: `APP_RNM/internal/services/pendaftaran.go:penomorCounter.NomorBerikut` satu-satunya penomor; `Generate_NoKlaim` nol kemunculan di `APP_RNM/`
- [x] Halaman React Register dapat mengirim pendaftaran dan menampilkan nomor yang diterima. — bukti: `APP_RNM/frontend/src/pages/claimlife/RegisterKlaim.tsx:RegisterKlaim` (`kirim` → `daftarKlaimLife`, lalu menampilkan `nomorKlaim`)
- [x] Nomor yang dihasilkan berbentuk `<prefix>K<kode bisnis>.MM.YYYY.<5 digit>` — contoh
      `RNML-KL1.08.2026.00936`. — bukti: `APP_RNM/internal/services/pendaftaran.go:RakitNomorKlaim`; uji `TestBentukNomorKlaim`
- [x] Prefix diperoleh lewat **lookup** ke `POOLDATA.KODE_PRODUKSI` (`TYPE='LIFE'`), **tidak**
      ditanam sebagai konstanta di kode. — bukti: `APP_RNM/internal/repository/penomor.go:Penomor.AwalanProduksi`; uji `TestNolAwalanNomorKlaimSebagaiLiteral`
- [x] Periode nomor mengikuti `POOLDATA.TANGGAL_CLOSING`, termasuk aturan cutover
      `TRUNC(now) <= 02/01/2026` → `12.2025`. — bukti: `APP_RNM/internal/repository/penomor.go:HitungPeriodeNomor`; uji `TestPeriodeNomorMengikutiHariTutupBuku`, `TestPeriodeCutoverDipertahankan`
- [x] Batas transaksi dipegang **Go**: commit terjadi segera setelah nomor terbentuk, sehingga lock
      `SELECT … FOR UPDATE` pada `GENERATE_SEQUENCE_NUMBER` tidak menahan pendaftar lain. — bukti: `APP_RNM/internal/services/pendaftaran.go:Pendaftaran.Daftar` (nomor di transaksi pertama); uji `TestPenomoranDiTransaksiSendiri`
- [ ] Dua pendaftaran serentak tidak pernah memperoleh nomor yang sama. — belum: hanya bentuk dua-transaksi yang dijaga statik (`TestPenomoranDiTransaksiSendiri`); nol uji dua sambungan serentak, dan cabang baris-kunci-pertama (`INSERT` tanpa kunci di `UrutNomorBerikut`) dapat berlomba

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §2b

- [ ] ⚠️ Klaim baru ditulis ke **`T_GENERAL_CLAIM`**, **dan** di-`INSERT` flat ke
      `OS_AKSEPTASI_KLAIM_LIFE` — hilir masih membaca dari sana. Test yang **menolak** penulisan
      `OS_AKSEPTASI_KLAIM_LIFE` justru **gagal**. Keduanya dalam **satu transaksi**.
      **REVISI 2026-09-18:** ⛔ pendaftaran **tidak** menulis tabel polis maupun marketing — keduanya
      **dihapus**. Test yang menemukan penulisan ke tabel polis atau marketing milik klaim
      **gagal**. *(AC 32 spec — koreksi 2026-09-16; `[keputusan work owner]`)* — **kode ada sejak GILIRAN-14 butir bp**: peserta terpilih kini lahir BERBARIS (`SavePesertaClaim` 7.8), jadi `BarisLamaDari` menghasilkan satu baris datar per peserta dan `OS_AKSEPTASI_KLAIM_LIFE` tertulis di transaksi pendaftaran. ⚠️ Bukti Oracle-nya `TestDaftarMenulisTigaTempatDanBarisDatar` (bertag `db`, **melewati** tanpa `ORACLE_DSN`) — AC ini tetap `[ ]` sampai uji itu dijalankan terhadap skema uji
- [x] ⚠️ **Tidak ada blob JSON** sebagai penyimpan isi klaim. *(AC 31 spec; penyimpangan sadar 1)* — bukti: uji `TestKolomUangDesimalDanNolJSON` (nol `JSON`/`CLOB`/`BLOB` di DDL klaim)
- [ ] Header memuat keempat field `PremiumListSummary` — `CLAIM_NO`, `PL_NUMBER`, `RISLIPRNM`,
      `BUSINESS_NAME` — beserta `CASEID` dan `CLAIM_RETRO`. *(AC 36 spec)* — belum: `PohonKlaim.Simpan` hanya mengisi `CLAIM_NO`, `POLICY_NO`, `BUSINESS_NAME` (dan `Daftar` tidak mengisi `NamaBisnis`); `RI_SLIP_RNM` dan jumlah `CLAIM_RETRO` tidak ditulis; `CASEID` pindah ke `T_WORK_CLAIM`
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS, MENUNGGU JAWABAN.**
      `T_CLAIM_POLICY` **dihapus**, dan kata **snapshot** **DICABUT**: data polis kini **dibaca
      hidup** ⚠️ **penyimpangan sadar** — polis yang berubah sesudah klaim dibuat **akan** mengubah
      tampilan klaim lama. Empat dari tujuh hal yang didaftar AC ini — ketiga **tanggal** dan
      **team group** — ⚠️ `[terbuka]` **kehilangan rumah**: bukan pindah, tetapi belum punya tempat
      (tiket 14 §Blocker). **Jangan tebak.** *(AC 37 spec)* — belum: `[terbuka]` sejak 2026-09-18 — kosong artinya; rumah ketiga tanggal dan team group menunggu work owner
- [ ] ⚠️ `[terbuka]` **AC INI KOSONG ARTINYA sejak 2026-09-18 — TIDAK DIHAPUS.** Atribut polis
      **tidak lagi disimpan di klaim sama sekali**, jadi "sekali per klaim" tidak punya yang
      dihitung. Yang menggantikan: ia **dibaca hidup** dari tabel polis lewat penunjuk di
      `T_GENERAL_CLAIM` ⚠️ **penyimpangan sadar**. *(AC 35 spec)* — belum: `[terbuka]` sejak 2026-09-18 — kosong artinya; atribut polis dibaca hidup, tidak ada yang dihitung
- [x] ⚠️ **REVISI 2026-09-18 — berpindah pemilik, tidak batal.** `PRODUCT_NAME` /
      `PRODUCT_NAME_ID` **bukan lagi kolom klaim**; keduanya termasuk 27 kolom yang dibaca dari
      `T_PREMIUM_LIST`. Larangan membaca `m_product_life.JSONDATA` kini mengikat **modul PremiumList
      Life**. Yang mengikat di sini: pendaftaran klaim **tidak** menyimpan nama/id produk.
      *(AC 38 spec)* — bukti: `APP_RNM/internal/repository/pohonklaim.go:PohonKlaim.Simpan` — nol kolom `PRODUCT_NAME`/`PRODUCTNAME` di INSERT header, peserta (`kolompeserta.go:insertPeserta`), maupun baris datar
- [ ] ⚠️ Keadaan tangga klaim (posisi, status akseptasi, `Type`) ditulis ke **`T_WORK_CLAIM`**,
      **bukan** sebagai kolom `T_GENERAL_CLAIM`. *(AC 46 spec; penyimpangan sadar 6)* — belum: posisi dan `Type` ditulis ke `T_WORK_CLAIM` (`PohonKlaim.Simpan`), tetapi `T_WORK_CLAIM` tak punya kolom status akseptasi — cermin `STS_REJECT`/`ACCEPTED_NO` tinggal di `T_GENERAL_CLAIM`
- [ ] Pendaftaran klaim berjalan dalam **satu transaksi** — penulisan `T_GENERAL_CLAIM`, baris
      `T_WORK_CLAIM`-nya, dan `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE` **utuh atau tidak sama
      sekali**. **REVISI 2026-09-18:** frasa "ketiga tabelnya" **dicabut** — tabel polis dan
      marketing dihapus. *(AC 49 spec)* — belum: transaksi kedua `Pendaftaran.Daftar` memuat work + header + peserta, tetapi INSERT datar bernol baris saat daftar (lihat AC penyimpanan di atas)

### Pencarian peserta — **hanya peserta hidup** `[keputusan work owner]`

> ⛔ **Dipertajam 01-10-2026 (OQ-N14, keputusan work owner):** "hidup" kini dinilai pada baris **versi polis
> terakhir** tiap sertifikat, bukan pada sembarang baris. Rinciannya di bab bertanggal
> *"Keputusan bertanggal — 1 Oktober 2026 (OQ-N14)"* di ekor tiket ini.

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

- [x] Pencarian peserta **tidak menampilkan** peserta ber-`EDMSTATUS` `'Batal'` maupun `'Delete'`.
      *(AC 25 spec)* — bukti: `APP_RNM/internal/repository/pesertapolis.go:PesertaHidup`; uji `TestPesertaHidupMenyaringBatalDanDelete`
- [x] Peserta **new business** (`EDMSTATUS` kosong/NULL) **tetap muncul**. ⚠️ Penyaring naif
      `EDMSTATUS NOT IN ('Delete','Batal')` membuang seluruh peserta NB di Oracle — test wajib
      memuat kasus ini dan **harus gagal** bila penyaringnya naif. *(AC 26 spec)* — bukti: `APP_RNM/internal/repository/pesertapolis.go:PesertaHidup`; uji `TestPesertaHidupMenyaringBatalDanDelete`, `TestPesertaNBTetapHidupDiJalurBacaKlaim`
- [x] Peserta ber-`EDMSTATUS` `'Old'` dan `'New'` **tetap muncul**. *(AC 27 spec)* — bukti: uji `TestPesertaHidupMenyaringBatalDanDelete` (kasus `Old`/`New`)
- [ ] Baris bernilai **negatif** hasil jurnal balik endorsement **tidak pernah** sampai ke layar
      Register maupun ke perhitungan klaim. *(AC 28 spec)* — belum: tidak ada penyaring baris bernilai negatif — hanya penyaring `EDMSTATUS` (`PesertaHidup`)
- [x] Penyaringan terjadi di **satu tempat** — repository pembaca peserta. Test yang menemukan jalur
      baca peserta tanpa penyaring **gagal**. *(AC 29 spec)* — bukti: uji `TestPenyaringPesertaHanyaSatuTempat`, `TestSQLCariPesertaSelaluBerpagar`
- [x] `STATUS` dan `STATUSOLD` **tidak** dipakai sebagai penanda hidup/mati. *(AC 30 spec)* — bukti: uji `TestStatusDanStatusOldBukanPenandaHidup`
- [ ] Jalur baca **akuntansi/ringkasan premium** — bila ada di konteks ini — **tidak** menyaring:
      ia melihat seluruh baris positif dan negatif. **Satu tabel, dua sudut pandang.** — belum: tidak ada jalur baca akuntansi/ringkasan premium di Claim Life untuk diperiksa

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

---

## Implementasi — 26 September 2026 malam (sesi batch 02–06)

**Keputusan §2 yang dipakai:** **aa** `[DIPUTUSKAN]` (sequence `SEQ_WORK_CLAIM`, migrasi `009`),
**z1** `[DIPUTUSKAN]` (kolom `CURRENCY` di `002`), **ab** `[DIPUTUSKAN]` (stub pelaku berpagar).
**Yang ditunggu:** **o1–o3** — dan itulah yang menahan sebagian besar AC tiket ini.

**Status: `claimed`** — **7 dari 26 AC tertutup**. Test Go **99 → 111**, nol FAIL; test bertag
`db` **20 → 23**, seluruhnya SKIP; frontend 87 → 88 modul. Angka diambil **sesudah** perbaikan
`/code-review`.

### Apa yang sekarang berjalan

`POST /api/klaim-life` dan `GET /api/peserta-life?pl=…&n=…` ada; `services.Pendaftaran.Daftar`
menulis baris work object, header, peserta terpilih, **dan** baris datar warisan dalam **satu
transaksi**, memakai `PohonKlaim.Simpan` yang sudah ada — nol `INSERT` baru untuk tabel yang sudah
punya penulis. Halaman React `Register` mencari peserta, memilih, dan mengirim.

### ⛔ Kenapa nomor klaim belum dapat dibentuk

AC 2, 3, dan 7–11 menuntut nomor datang dari `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`. Keputusan
work owner **o** (26-09-2026) melarang memanggil procedure mana pun. Keduanya tidak dapat benar
sekaligus, dan teks AC hanya boleh diubah work owner.

Yang dilakukan executor: **memisahkan tempatnya**, bukan memilih pihak. `services.Penomor` adalah
antarmuka; implementasi bawaannya `PenomorBelumDiputuskan` mengembalikan **galat terang** yang
menyebut butir o dan nomor AC yang bertentangan. Pintu HTTP menjawab **501**, bukan 500 — ini bukan
kerusakan melainkan keputusan yang belum diambil. ⛔ Nomor karangan yang tampak benar jauh lebih
berbahaya daripada galat: nomor klaim dibaca manusia dan dipakai di luar sistem ini.

⭐ Satu test mengunci bahwa **nomor tidak diambil bila permintaan ditolak**. Nomor yang sudah
terbentuk tidak dapat dikembalikan ke urutannya, jadi mengambilnya sebelum validasi berarti membuang
satu nomor setiap kali borang salah isi — dan lubang nomor itu terlihat oleh orang di luar sistem.

### Butir aa — identitas work object

`SEQ_WORK_CLAIM` (langkah `009`) memberi **angka urutannya**; awalan `CLM-`/`KMT-` dan `LPAD(6)`
dirakit di Go, sebab awalan bergantung jenis baris dan itu aturan dagang, bukan DDL. ⛔ **Tanpa
reset tahunan** — korpus tidak memuat satu pun bukti bahwa Pega me-reset urutan ini, dan mengarang
reset membuat dua baris bernomor sama pada tahun berbeda.

### Butir ac — kolom negatif, dijawab korpus

Brief meminta executor membaca verdict **V14** grilling Endorsement Life. Hasilnya: V14 **tidak
menyebut kolom bertanda negatif** sama sekali. Yang ditetapkannya adalah **mekanismenya** — baris
negatif hasil jurnal balik ditulis ke tabel yang **sama** dan hidup berdampingan dengan baris
positifnya; akuntansi melihat keduanya, dan kontrak Claim Life berbunyi *"peserta yang sudah EDM
Batal atau soft-delete TIDAK BOLEH MUNCUL"*.

Jadi AC 23 dipenuhi lewat penyaring `EDMSTATUS`, bukan lewat penyaring tanda. ⛔ Executor **tidak
mengarang** kolom bertanda negatif yang tidak ada di mana pun.

### ⛔ Tabel peserta berisi 66,8 juta baris — dan itu mengubah cara kodenya ditulis

Katalog instance pengembangan: index pada `PL_NUMBER`, `CERTIFICATE_NO`, `POLICY_NO`, dan **tidak
ada** index berawalan `EDMSTATUS`. Query tanpa penyaring ber-index bukan "agak lambat" melainkan
pemindaian penuh yang menahan basis data. Dua aturan dipasang dan **dijaga test statik**: setiap
query menyaring dengan `PL_NUMBER`/`CERTIFICATE_NO`, dan setiap query berbatas hasil.

⚠️ Penyaring naif `EDMSTATUS = ''` **keliru**, dan bukan sedikit: agregat menghitung **59,1 juta**
baris `NULL` dan **nol** baris teks kosong — peserta new business justru yang `NULL`, sehingga
penyaring naif membuang hampir seluruh tabel tanpa satu pun galat.

### Butir ab — stub pelaku, dan apa yang ia BUKAN

⛔ **Ini bukan autentikasi.** Header `X-Pelaku`/`X-Peran` dapat ditulis siapa saja yang dapat
mengirim permintaan. Tiga pagar dipasang: mati secara bawaan; **ditolak saat memuat konfigurasi**
bila `IS_PEGA_PROD=true` (bukan saat permintaan pertama — proses yang menyala akan melayani
permintaan sebelum ada yang sempat menyadarinya); dan seluruh uji wewenang berjalan di seam
`services` dengan `Pelaku` langsung, sehingga aturannya tidak bergantung pada jalur ini sama sekali.

### AC yang masih terbuka — 18

| Sebab | Nomor |
| --- | --- |
| **Menunggu butir o** (teks AC bertentangan dengan keputusan work owner) | 2, 3, 7, 8, 9, 10, 11 |
| ⛔ **Penyaring tanda tidak dipasang** — dicabut centangnya sesudah tinjauan | 23 |
| Menunggu Oracle (G1) — jalurnya ada, buktinya belum | 1, 4, 12, 14, 18, 19 |
| `[terbuka]` sejak 2026-09-18, kosong artinya | 15, 16 |
| Pemilik lain — PremiumList Life · layar Register bernomor | 17 · 6 |
| Tidak berlaku di konteks ini (jalur baca akuntansi) | 26 |

### Penjaga yang tersentuh, dan angkanya diperbarui — bukan dilonggarkan

`TestSeluruhCreateDapatDibacaNamanya` **19 → 20** (sequence baru); `TestKolomDDLCocokDenganStruktur`
menuntut STRUKTUR diralat untuk `CURRENCY`, dan blok ralat bertanggal ditulis di sana.

⚠️ **Dua penjaga baru gagal membuktikan dirinya pada percobaan pertama**, dan keduanya diperbaiki
sebelum dipakai: penjaga query mencari `PL_NUMBER` **di mana pun** — termasuk di daftar kolom
`SELECT` dan di `ORDER BY` — sehingga meluluskan query yang penyaringnya dicabut; dan penjaga AC 25
mencocokkan `STATUS` di dalam `EDMSTATUS`, menuduh penyaring yang justru benar.


### ⛔ Yang ditemukan `/code-review`, dan diperbaiki sebelum commit

| # | Temuan | Perbaikan |
| ---: | --- | --- |
| 1 | ⛔ **Jalur tulis tidak berpagar wewenang.** `Daftar` tidak pernah memeriksa pelaku, sehingga permintaan tanpa identitas — keadaan **bawaan** saat stub mati, yaitu keadaan produksi — tetap menulis klaim dengan `CREATE_OP_NAME` kosong. Bukan gagal-tertutup, melainkan **gagal-anonim** | pendaftaran menolak pelaku tanpa identitas; ⚠️ peran yang sebenarnya tetap `[terbuka — tiket 07]` |
| 2 | ⛔ **Penjaga AC 25 MATI.** Polanya ditulis `\b` di dalam raw string, dan di sana dua backslash berarti backslash **harfiah** — ia tidak pernah cocok dengan apa pun. Bab ini sempat membanggakannya sebagai penjaga yang diperbaiki | → ``; dibuktikan dengan menjalankan kedua pola berdampingan atas `"WHERE STATUS = 1"` |
| 3 | ⛔ **AC 23 dicentang tanpa dasar.** Penyaring tanda memang tidak dipasang, dan bab ini mengakuinya di paragraf yang sama | centangnya **dicabut** |
| 4 | ⛔ **Sequence dibakar sebelum penomoran yang pasti gagal.** `SEQ_WORK_CLAIM` non-transaksional, jadi setiap `POST` membuang satu pengenal lalu menjawab 501 — melanggar prinsip yang berkas itu sendiri tulis | urutan ditukar: nomor dulu, baru pengenal |
| 5 | ⛔ **AC 29 justru dilanggar**: aturan penyaring hidup ada **dua kali** dan berbeda — SQL tanpa `TRIM`, Go dengan `TRIM`, dan `PesertaHidup` nol pemanggil produksi. `"Batal "` berspasi lolos SQL | `TRIM` di SQL, dan hasil disaring ulang lewat `PesertaHidup` sehingga aturannya punya satu rumah |
| 6 | ⛔ **Butir z1 separuh**: kolom `CURRENCY` ada di DDL tetapi tidak ditulis maupun dibaca — daftar klaim USD, baca lagi, mata uangnya hilang | ditulis di `INSERT` header dan dibaca `AmbilHeader` |
| 7 | ⚠️ **Test db tiket 02 hilang seluruhnya** — satu baris utuh §4 brief tidak dikerjakan | tiga test db ditulis: tiga tempat + baris datar, dua pengenal berbeda, dan penomoran gagal membatalkan seluruh transaksi |
| 8 | ⚠️ `stubPelakuAktif` variabel paket yang diubah `Router` — dua Router dalam satu proses berbagi satu saklar | dibawa sebagai argumen |
| 9 | ⚠️ Pesan 501 membocorkan nama objek basis data ke badan HTTP | rinciannya tinggal di log |
| 10 | ⚠️ Frontend: kunci pilihan memakai nomor sertifikat saja, sehingga sertifikat kembar antar polis tercentang berbarengan; polis dan mata uang diambil dari peserta pertama tanpa memeriksa campuran | kunci menjadi polis+sertifikat; campuran polis/mata uang ditolak di layar |

⚠️ **Pelajaran yang saya catat sendiri:** saat memeriksa temuan 2, alat ukur pertama saya **tertelan
perangkap yang sama** — heredoc memakan backslash gandanya, sehingga kedua pola menjadi identik dan
hasilnya "tidak ada masalah". Saya nyaris menolak tuduhan yang benar. Yang menyelamatkannya hanya
satu: angka yang identik untuk pola yang seharusnya berbeda terlihat mencurigakan, dan diukur ulang
dengan alat yang ditulis tanpa heredoc.

### Yang MASIH belum dikerjakan dari §4 brief

⚠️ **Peserta belum disalin lengkap.** Brief menuntut peserta terpilih disalin *"termasuk keempat
tanggal valuasi, `WPC`, `IS_CHECK`"*; `models.Peserta` belum punya medan itu dan `INSERT` peserta
belum menulisnya. ⛔ **Akibatnya tiket 06 terblokir**: §3 brief menyandarkan validasi DOL pada kolom
itu justru supaya tidak perlu query ulang ke tabel 66 juta baris.

### Penjaga ketiga §5 — dicatat di sini karena belum dapat dijalankan

Cacah objek sesudah `-migrate` berubah **8 · 5 · 15 → 8 · 6 · 15**: satu sequence baru
(`SEQ_WORK_CLAIM`, langkah `009`), nol tabel dan nol index baru. ⛔ Belum dapat dibuktikan — G1
tertutup dan `-migrate` belum pernah berjalan di Oracle mana pun.


---

## Implementasi — lanjutan tiket 02, 26 September 2026 malam

⛔ **Ralat atas bab di atas: commit `753cef2` MEMUAT SATU TEST MERAH.** Bab itu dan pesan commit-nya
menulis *"111 test, nol FAIL"*; yang sebenarnya **110 PASS + 1 FAIL**.
`TestSetiapPemanggilBukaMemeriksaBolehDilewati` mengunci cacah pemanggil `skemauji.Buka()` = 7,
sedangkan `services/pendaftaran_db_test.go` — yang lahir dari perbaikan tinjauan nomor 7 — menjadi
pemanggil **kedelapan**.

**Sebabnya satu, dan bukan kebetulan:** verifikasi penuh saya jalankan **sebelum** perbaikan
`/code-review`, lalu tidak diulang sesudahnya. Klaim lama **tidak dihapus** dari bab di atas; ia
diralat di sini, supaya jejaknya tetap terbaca.

### Yang dikerjakan lanjutan ini

| # | Isi |
| ---: | --- |
| 1 | Cacah pemanggil **7 → 8**; `go test ./...` kini **111 PASS, 0 FAIL** |
| 2 | ⭐ **Peserta disalin lengkap** — `models.Peserta` bertambah `SumberID`, `IsCheck`, keempat tanggal valuasi, `WPC`, empat tanggal polis, `STNC`, delapan medan uang, dan `EMPercent` (`Ratio`). Inilah **prasyarat tiket 06**: validasi DOL membaca jendela valuasi dari sini, bukan bertanya ulang ke tabel 66,8 juta baris |
| 3 | ⛔ **Server membaca ulang peserta sendiri** lewat `(PL_NUMBER, CERTIFICATE_NO)` — keduanya ber-index. Klien kini hanya mengirim **nomor sertifikat**; nilai polis tidak pernah dipercaya dari badan HTTP, sebab nilai itu menentukan angka klaim dan jendela DOL |
| 4 | Satu daftar kolom (`kolompeserta.go`) dipakai **tulis dan baca**. Menulisnya dua kali adalah cara paling mudah membuat urutan bind berselisih, dan selisih itu tidak terlihat sampai ada nilai yang mendarat di kolom yang salah |
| 5 | `CalonPeserta` diberi tag JSON camelCase; `api.ts` dan halaman Register mengikuti |
| 6 | ⭐ **butir ae1** — `CASEID` = pengenal work object |

### ⭐ Butir ae1 — dan ralat atas penolakan saya sendiri

Ronde tinjauan lalu saya menolak menyamakan `CASEID` dengan pengenal work object, dengan alasan
*"satu nilai dua arti"*. **Itu keliru.** Di Pega pun `CASEID` **adalah** pengenal work object-nya,
jadi menyamakannya adalah **paritas** dengan sistem berjalan. Yang justru merusak adalah
membiarkannya kosong: baris datar warisan klaim baru tidak dapat dikelompokkan hilir yang membaca
`OS_AKSEPTASI_KLAIM_LIFE` per `CASEID`, dan `Hapus` serta `CacahBarisLama` memakai sumbu itu.
`[keputusan work owner butir ae1]`.

### ⛔ Nol nama orang disalin

`NAME_OF_INSURED` dan `POLICY_HOLDER` punya kolomnya di DDL, tetapi **tidak disalin**. Nama
tertanggung hanya diperlukan **layar** saat memilih, dan itu dilayani `CalonPeserta`. Menyalinnya ke
tabel klaim berarti menduplikasi data pribadi tanpa satu pun AC yang memintanya. Kolomnya tetap
`NULL`. ⛔ Kolom `KTP` tidak pernah dibaca sama sekali. `[terbuka — work owner]`

### Ralat cacah AC terbuka

Daftar "AC yang masih terbuka" di bab sebelumnya menulis **18**; yang benar **19**
(7 butir o + 1 penyaring tanda + 6 menunggu Oracle + 2 kosong artinya + 2 pemilik lain + 1 tidak
berlaku). Sensusnya sendiri benar: **7 `[x]` + 19 `[ ]` = 26**.

---

### Ralat menurut keputusan work owner — 27 September 2026 (A2, butir o1/o2/o3)

⛔ **AC yang menyebut procedure diralat.** Teks lama menuntut memanggil
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`; `[keputusan work owner]` *"jangan ada lagi pemanggilan
procedure, segala procedure hardcode dalam skrip"* melarangnya. **ADR-U-0043** ditulis, dan ia
**meng-supersede ADR-U-0006** pada bagian penomoran bisnis.

| Butir | Teks lama | Teks baru |
| --- | --- | --- |
| penerbit nomor klaim | panggil `PROC_GENERATE_SEQUENCE_NUMBER` | logikanya **ditulis ulang di Go** (`PenomorCounterOracle`); yang tetap di Oracle hanya `SELECT … FOR UPDATE` |
| periode | *(tidak disebut)* | hari tutup buku dari `TANGGAL_CLOSING`; lewat hari itu → geser satu bulan **beserta tahunnya** |
| cabang cutover | *(tidak disebut)* | `TRUNC(now) <= 02/01/2026` → `12.2025`/`2025`; **dipertahankan** meski sudah lewat, sebab pengurai nomor lama tiket 13 memerlukannya |
| bentuk nomor | *(tidak disebut)* | `RNML-K<kode bisnis>.<MM.YYYY>.<5 digit>` |

⚠️ **Tulisan ke tabel warisan yang disengaja:** `GENERATE_SEQUENCE_NUMBER` di-`UPDATE`/`INSERT`
aplikasi. Inheren pada keputusan ini — penghitung yang tidak disimpan bukan penghitung. Dicatat,
bukan disembunyikan.

⚠️ **Label tidak dinaikkan.** Bentuk procedure-nya berlabel `[data DBA — belum dikonfirmasi DBA]` di
sumbernya; kode yang menirunya memakai label yang sama.

`PenomorBelumDiputuskan` **tetap menjadi bawaan** — gagal tertutup. `register.go` yang memasang
`PenomorCounterOracle`, sehingga jalur pendaftaran berhenti menjawab 501.

## Implementasi — 27 September 2026 (A2, penutupan stub o1)

**Empat centang bergeser, dan sebabnya:**

| AC | Sebab bergeser |
| --- | --- |
| bentuk `<prefix>K<kode bisnis>.MM.YYYY.<5 digit>` | `RakitNomorKlaim` diuji dengan contoh AC-nya sendiri (`RNML-KL1.08.2026.00936`), termasuk urut > 5 digit yang **tidak** dipotong. Sebelum ini fungsinya **tanpa uji sama sekali** |
| prefix lewat lookup `KODE_PRODUKSI` | konstanta `AwalanNomorKlaim = "RNML-"` **dibuang**; `AwalanProduksi(ctx, tx, "LIFE")` membacanya saat jalan. `[terverifikasi]` `RDBList/GetKodeProdLife_SQL.xml` baris 85. Penjaga arah-balik: nol literal `"RNML-"` di kode services |
| periode mengikuti `TANGGAL_CLOSING` + cutover | `HariClosing` + `HitungPeriodeNomor`; cutover `TRUNC(now) <= 02/01/2026 → 12.2025` diuji `TestPeriodeCutoverDipertahankan` |
| batas transaksi dipegang Go | penomoran dipindah ke **transaksi sendiri**, di-commit sebelum transaksi pendaftaran. Sebelumnya kuncian `SELECT … FOR UPDATE` dipegang sampai seluruh pendaftaran selesai |

**Yang TIDAK bergeser, dan sebabnya:**

| AC | Sebab tetap terbuka |
| --- | --- |
| nomor dari `PROC_GENERATE_SEQUENCE_NUMBER` | ⛔ **DIGANTIKAN ADR-U-0043** `[keputusan work owner, butir o1]`: procedure tidak dipanggil. Teks AC ini perlu ditulis ulang work owner — executor tidak menutup `[terbuka]` |
| aplikasi tidak memuat logika format nomor *(ADR-0006)* | ⛔ **DIGANTIKAN ADR-U-0043** untuk nomor bisnis. Sama: teks AC menunggu work owner |
| didaftarkan lewat API dan muncul berstatus awal | menuntut jalan Oracle sungguhan; belum dijalankan |
| dua pendaftaran berurutan → dua nomor berbeda | sama — menuntut Oracle |
| dua pendaftaran **serentak** tidak pernah senomor | sama, dan menuntut **dua sambungan** serentak. Bentuknya dijaga statik (`TestPenomoranDiTransaksiSendiri`), perilakunya belum |
| halaman React Register | kelompok A3 |

**Konsekuensi yang diterima** *(dinyatakan, bukan disembunyikan)*: dua transaksi berarti nomor
dapat **terbakar** bila transaksi kedua gagal — urutannya berlubang. Lubang tidak merusak apa pun;
kuncian global merusak setiap pendaftaran serentak.

**Terbuka baru:** awalan **akseptasi** (`RNML-A`, `RNML-AR`) masih dirakit di kode — keduanya tidak
ada di `KODE_PRODUKSI`. Dari mana huruf `A`/`AR` datang adalah `[terbuka — work owner]`.

## Ralat menurut XML — 27 September 2026 (butir av)

**Yang diralat:** ronde A3 Register pertama membangun `Ceding`, `Class of Business`, dan
`Marketing Officer` sebagai **isian ber-autocomplete**, lengkap dengan rute
`GET /api/rujukan/{jenis}` dan tiga pembaca tabel. **KELIRU.**

**Bukti, path + baris** *(`Section\InputRegisterClaimLife.xml`)*:

| Medan | `pyReadOnly` | `pyEditOptions` | `pyLabelFor` | terikat |
| --- | --- | --- | --- | --- |
| `Type` | `true` b9058 | `Read-only` b9068 | `Type` b9093 | `.PolicyDataLife.Type` |
| `Marketing Officer` | `true` b9845 | `Read-only` b9855 | `MarketingName` b9879 | `.PolicyDataLife.MarketingName` |
| `Ceding` | `true` b11493 | `Read-only` b11503 | `CedingCoName` b11529 | `.PolicyDataLife.CedingCoName` |
| `Class of Business` | `true` b12127 | `Read-only` b12135 | `BusinessName` b12159 | `.PolicyDataLife.BusinessName` |

Keempatnya **read-only** dan terisi dari halaman polis — bukan diketik.

⚠️ **Sebab salah bacanya layak dicatat**, sebab ia akan terulang: blok kontrol berdiri **SEBELUM**
labelnya di DOM section, dan pembacaan ronde pertama memakai jendela **ke depan** dari label. Itu
grep dengan langkah tambahan, bukan pembacaan pohon.

**Keputusan av** `[DIPUTUSKAN work owner 27-09-2026]`: sumber kesebelas medan `PolicyDataLife`
adalah modul **PremiumList Life**. Claim Life **tidak** membaca cermin JSON-nya
*(`JSON_POLIS`, `POLICYJSONLIFE`, `SEARCH_POLIS`, `pc_*`)*; medannya diisi kelak dari tabel
relasional modul itu — sejalan dengan ralat 2026-09-18 di kepala tiket ini *("dibaca hidup dari
tabel polis, `T_PREMIUM_LIST` dkk")*.

**Yang dibuang:** `repository/rujukan.go`, `services/rujukan.go` *(+ujinya)*, `handlers/rujukan.go`,
rute `GET /api/rujukan/{jenis}`, dan `api.ts` `cariRujukan`/`JENIS_RUJUKAN`/`BarisRujukan`. Kode mati
di tiga lapis; **nol** berkas `.tsx` pernah memakainya. Tidak dipakai ulang untuk `PreCaimLife_Act`:
pencarian `BUSINESS` ber-`OLDID` pun bergantung `PolicyDataLife.BusinessCode`, jadi ia **ikut
menunggu modul** — menyimpan kodenya berarti menyimpan kode mati kedua.

**Di layar sekarang:** `PanelDataPolis` menampilkan **seluruh** medan VERBATIM dengan penanda
"menunggu modul PremiumList Life"; ujinya mengunci cacahnya pada **11**.

**AC:** tidak ada AC tiket ini yang berubah centangnya — yang diralat adalah cara membacanya.

## Ralat menurut XML — 27 September 2026 (paket Register(2), lanjutan 10 §1)

### 1. `Find Insured` menyaring TIGA kriteria, bukan dua

`RDBList/GetPesertaClaim_sql1.xml:85` *(dipanggil `LoadDataPesertaSpesifik_Act` langkah
`RDB-List` b485)*:

> `SELECT * FROM POOLDATA.M_LIFE_PREMIUM_DETAIL`
> `WHERE PL_NUMBER = {pyWorkPage.PolicyDataLife.PremiumListSummary.PL_NUMBER}`
> `  AND CERTIFICATE_NO LIKE '%'||{SearchPolicyHolder.CARI2}||'%'`
> `  AND UPPER(NAME_OF_INSURED) LIKE '%'||{SearchPolicyHolder.CARI3}||'%'`

| Kriteria | Perlakuan | Bukti |
| --- | --- | --- |
| `PL_NUMBER` | **wajib**, sama dengan | b85 |
| `CERTIFICATE_NO` | `LIKE` berpagar `%`, **tanpa** `UPPER` | b85 |
| `NAME_OF_INSURED` | `LIKE` berpagar `%`, **dengan** `UPPER` di kedua sisi | b85 + `@toUpperCase` b405 |

Rutenya karena itu `GET /api/peserta-life?pl=&sertifikat=&nama=&n=`. Meng-`UPPER` sertifikat pun
akan membuat pencariannya berhenti memakai index-nya — dan tabel itu 66,8 juta baris.

⚠️ **Penyimpangan sadar, dilaporkan OQ-E**: Pega memasang kedua `LIKE` **tanpa syarat**, dan di
Oracle `X LIKE '%'` bernilai FALSE ketika `X` NULL. Kotak kosong di Pega karena itu membuang
peserta ber-nama NULL. Kita memasang `LIKE` hanya untuk kotak yang terisi.

### 2. `+7 jam` di `LoadDataPesertaSpesifik_Act` TIDAK ditiru

Langkah b662 mengulang hasil dan menambah **7 jam** ke delapan medan tanggal *(b738 `DOB`, b784
`BEGIN_DATE`, b804 `EXPIRED_DATE`, b824 `EFFECTIVE_DATE`, b844/b864 valuasi gross, b884/b904
valuasi retro)* lewat `@addCalendar(...,0,0,0,0,7,0,0)`.

Itu tambalan zona waktu JDBC Pega *(UTC → WIB)*, bukan aturan bisnis. Pembaca kita mengambil
tanggal lewat `TO_CHAR(..,'YYYY-MM-DD HH24:MI:SS')`, sehingga tidak ada pergeseran yang perlu
ditambal. **Menirunya justru akan menggeser tanggal tujuh jam ke depan.**

### 3. `Select Insured` — dua medan adalah PILIHAN, bukan salinan

`SaveInsuredClaim_Act` memindahkan 38 medan; **36 di antaranya salinan lurus** yang backend sudah
salin sendiri saat `POST` *(`kolomSalin`)*. Dua sisanya aturan:

| Medan | Aturan | Baris |
| --- | --- | --- |
| `AGE` | `AGE` → `ENTRY_AGE` → `CURRENT_AGE`; hanya **kosong** yang menjatuhkan | b661 |
| `SHARE_NUSANTARA_RE` | `"0"` **atau** kosong → `SHARE_NUSANTARA_RE_GROSS` | b601, b1412 |

⛔ Keduanya **sengaja asimetris** dalam memperlakukan `"0"`, dan ujinya mengunci perbedaan itu.
Menyeragamkannya mengubah angka uang: umur bayi nol tahun akan naik menjadi `ENTRY_AGE`.

Uang di Pega dinormalkan `@divide(@toDecimal(@replaceAll(x,",",".")),1,4)` — perbaikan koma
desimal. Tidak ditiru dan tidak perlu: `TO_CHAR(..,'TM9','NLS_NUMERIC_CHARACTERS=''.,''')` sudah
memberi titik desimal kanonik langsung dari kolom `NUMBER`.

**AC:** tidak ada AC tiket ini yang berubah centangnya. Kolom `AGE` sudah disediakan migrasi 003
sejak awal — ia hanya tidak pernah terisi; **nol migrasi baru**.

## Ralat 28-09-2026 — pergeseran bulan `HitungPeriodeNomor` menjepit ke akhir bulan seperti `ADD_MONTHS`

⛔ **Ini mengubah keluaran penomoran Claim Life yang sudah jalan.** Dicatat di sini, bukan hanya di
tiket PremiumList tempat ia ketahuan, sebab fungsi yang diperbaiki adalah **milik penomoran klaim**.

**Ditemukan** saat tiket 03 PremiumList hendak memakai ulang `repository.HitungPeriodeNomor` —
bukan lewat uji yang merah, melainkan lewat pembacaan sebelum memakai.

### Cacatnya

`HitungPeriodeNomor` menggeser periode satu bulan dengan `saat.AddDate(0, 1, 0)` ketika hari
transaksi melewati hari tutup buku. **Go melimpahkan tanggal yang tidak ada**:

```
31 Januari 2026 + 1 bulan  ->  3 Maret 2026     (Go AddDate)
31 Januari 2026 + 1 bulan  ->  29 Februari 2026 (Oracle ADD_MONTHS, menjepit)
```

`[data DBA]` `SUMBER-PENOMORAN-DBA.md` menyebut `ADD_MONTHS(+1)`, dan Oracle `ADD_MONTHS`
**menjepit** ke akhir bulan tujuan — ia tidak pernah melimpah. Akibatnya periode nomor klaim
menjadi `03.2026` padahal seharusnya `02.2026`: **Februari terlewat sama sekali**.

### Kapan ia menyala

Setiap penerbitan nomor pada tanggal **29, 30, atau 31** dari bulan yang penggantinya lebih pendek.
Dengan hari tutup buku yang lazim (25), seluruh tanggal itu melewati ambang, jadi seluruhnya
bergeser — dan seluruhnya bergeser **ke bulan yang salah**. Tidak ada galat, tidak ada peringatan:
nomornya terbentuk, tersimpan, dan baru terlihat keliru saat rekonsiliasi tutup buku.

### Tujuh kasus yang mengunci perbaikannya

`TestPeriodeNomorTidakMelompatiBulanPendek`, hari tutup buku 25:

| Saat | Sebelum ralat | Sesudah ralat | Sebab |
| --- | --- | --- | --- |
| 2026-01-29 | `03.2026` | **`02.2026`** | Februari 2026 hanya 28 hari |
| 2026-01-30 | `03.2026` | **`02.2026`** | idem |
| 2026-01-31 | `03.2026` | **`02.2026`** | kasus yang paling jauh melimpah |
| 2026-03-31 | `05.2026` | **`04.2026`** | April 30 hari |
| 2026-05-31 | `07.2026` | **`06.2026`** | Juni 30 hari |
| 2026-08-31 | `10.2026` | **`09.2026`** | September 30 hari |
| 2026-10-31 | `12.2026` | **`11.2026`** | November 30 hari |

⚠️ Penjaganya **dibuktikan merah lebih dahulu** atas implementasi lama (ketujuh kasus gagal), lalu
hijau sesudah perbaikan.

### Bagaimana diperbaiki

Pergeserannya **didelegasikan** ke `models.PeriodeProduksi` — aturan periode tiket 02 PremiumList,
yang menaikkan **nomor bulan** dengan pergantian tahun, bukan menambah tiga puluh hari. Karena yang
dibutuhkan hanya bulan dan tahun, menaikkan nomor bulan setara persis dengan penjepitan
`ADD_MONTHS`.

Akibat sampingan yang disengaja: aturan periode kini **satu**, dipakai penomoran klaim maupun
premium list. Dua salinan aturan periode adalah dua periode yang suatu hari berselisih diam-diam.

### Perubahan kedua: zona pembacaan hari

Hari transaksi kini dibaca di zona **Asia/Jakarta**, sama dengan `SYSDATE` server yang dibaca
`TRUNC(v_now)` di procedure. Sebelumnya ia dibaca di zona `saat` sendiri — sehingga cabang cutover
dan pergeseran bulan dapat berselisih sehari di sekitar tengah malam, dan sehari di sini berarti
satu bulan buku. Dikunci `TestBatasCutoverDibacaDiJakarta`.

### AC

Tidak ada AC tiket ini yang berubah centangnya. Yang berubah adalah **keluaran** penomoran pada
tujuh tanggal di atas — ke arah yang benar. **Nol migrasi baru.**

Kodenya: `APP_RNM/internal/repository/penomor.go`, `HitungPeriodeNomor` (kini mengembalikan galat,
sebab `models.PeriodeProduksi` **menolak** hari tutup buku kosong atau di luar 1..31 alih-alih
menebaknya). Commit `8f69682`.

## ⛔ Catatan bertanggal — 29 September 2026 (GILIRAN-13 paket 2)

Baris adjustment pertama kini dapat lahir lewat tombol `Add` di layar Detail (tahap Claim Analis; tiket 03,
butir bo). **Pendaftaran tidak diubah** — peserta tetap lahir tanpa baris (AC 32 di atas tetap terbuka).
Pega melahirkan baris pertama saat pendaftaran (`SavePesertaClaim` 7.8 b3671, hidup, WHEN b3919); selisih
urutan kerjanya dicatat sebagai **OQ-N9** (tiket 03).

## ⛔ Ralat bertanggal — 29 September 2026 (GILIRAN-14 paket 1, butir **bp**: baris adjustment lahir saat Submit Register)

`[DIPUTUSKAN; veto work owner]` butir **bp**. Catatan GILIRAN-13 di atas ("pendaftaran tidak diubah";
**OQ-N9** "apakah pendaftaran semestinya melahirkan baris pertama") **dicabut**: itu bukan pertanyaan,
melainkan bunyi XML yang terlewat.

**Bukti** (pohon, `pyStepsBlockName` dicetak): `Activity/SavePesertaClaim.xml`, dipanggil tombol `Submit`
(`Section/InputRegisterClaimLife.xml` b27369 → b27393). Langkah 7 berulang atas `TempDetail.pxResults`
(b1461); **7.7** b2744 menambahkan peserta (`PremiumListDetail(<APPEND>)`) dan **7.8** b3671 menulis
`PremiumListDetail(<LAST>).AdjustmentList(<LAST>)` — keduanya hidup, keduanya ber-WHEN `.IsCheck=="true"`
(b3631 untuk 7.7, **b3919** untuk 7.8; True=2 lanjut, False=3 lewati). Delapan medan 7.8: `CEDING_RETENTION`
b3696, `SHARE_NUSANTARA_RE` b3742 (`@if` yang sama dengan peserta b2845), `SUM_INSURED` b3768, `SUM_REASURED`
b3788, `SHARE_RETRO` b3808, `CLAIM_AMOUNT` b3828, `RETROCEDED_SHARE` b3848 — ketujuhnya
`@divide(@toDecimal(@replaceAll(.X,",",".")),1,4)` — dan `CURRENCY` b3868.

**Yang dibangun:**
- `services.BarisPendaftaran` (murni) + `LahirkanBarisPendaftaran`, dipanggil `Pendaftaran.Daftar` di
  transaksi kedua, sebelum `Simpan`. Pembulatan empat angka lewat `utils.DecimalContext` (setengah ke atas),
  sama dengan rekap PremiumList. Status **tidak** ditulis (7.8 tidak menyentuhnya; `Save to RNM` 22.1.3.2
  yang menulis `0`). Kolom sumber kosong tetap kosong (`[dugaan]` Pega menjadikannya nol; ADR-U-0027).
- `CLAIM_AMOUNT` kini **dibaca** dari `M_LIFE_PREMIUM_DETAIL` (posisi 28 `kolomSalin`, di ekor supaya nol
  posisi lain bergeser) dan **disimpan** di peserta (`kolomPeserta`, 7.7 b3280) — sebelumnya tidak dibaca
  sama sekali. Penulis PremiumList sudah mengisinya (`polis_warisan.go`), jadi kontrak dua sisi tetap utuh.
- ⛔ **Temuan sampingan:** tiruan `M_LIFE_PREMIUM_DETAIL` di skema uji tertinggal empat kolom yang
  `kolomSalin` baca sejak lanjutan 10 (`SHARE_NUSANTARA_RE_GROSS`, `AGE`, `ENTRY_AGE`, `CURRENT_AGE`). Setiap
  uji `db` yang mendaftarkan klaim akan gagal ORA-00904 — tidak terlihat karena uji itu selalu melewati.
  Dilengkapi, dan kini dijaga `skemauji.TestTiruanPesertaPolisMemuatSetiapKolomSalin`.

**Uji:** `TestBarisPendaftaranMenyalinDelapanMedan7_8` (contoh literal, pembulatan terlihat),
`TestBarisPendaftaranHanyaBagiPesertaTerpilih` (`"false"`/kosong → nol baris),
`TestBarisPendaftaranMembiarkanKosongTetapKosong`, `TestLahirkanBarisPendaftaranSatuPerPeserta` (baris ini
yang diwarisi putaran), `TestPesertaMenyimpanJumlahKlaim`, `TestPesertaMembacaKembaliJumlahKlaim`,
`TestUrutanKolomSalinDikunci` (29 posisi); `db`: `TestDaftarMenulisTigaTempatDanBarisDatar` (baris + satu
baris datar).

⚠️ **Akibat bagi skema uji:** pendaftaran kini MENULIS `OS_AKSEPTASI_KLAIM_LIFE` (AC 32). Tabel itu harus ada
di skema uji sebelum `Submit` Register diuji di layar.

**Pertanyaan terbuka baru:**
- **OQ-N10** — 7.7 membulatkan nilai **peserta** ke empat angka (`@divide(…,1,4)` b2770–b3561); peserta Go
  disimpan apa adanya dari `NUMBER(38,8)`. Baris pertama (7.8) sudah dibulatkan, jadi pada sumber berdesimal
  lebih dari empat, nilai peserta dan baris pertamanya berbeda di angka kelima. Tiru juga pembulatan
  peserta?

## ⛔ Keputusan work owner bertanggal — 29 September 2026 (GILIRAN-15 paket 2: N10 dan "kosong = nol" di 7.8)

**OQ-N10 — DITUTUP: pembulatan 7.7 DITIRU.** `services.BulatkanPesertaPendaftaran`, dipanggil `Pendaftaran.Daftar`
sebelum baris 7.8 lahir, membulatkan sepuluh medan peserta ke empat angka setengah ke atas dengan `bulat` bersama —
`[terverifikasi]` `SavePesertaClaim.xml` 7.7 b2744: GROSS_PREMIUM b2770, NET_PREMIUM b2824, SHARE_NUSANTARA_RE b2845,
SUM_INSURED b2872, CEDING_RETENTION b2893, SUM_REASURED b3107, EM_PERCENT b3134, CLAIM_AMOUNT b3281, SHARE_RETRO
b3401, RETROCEDED_SHARE b3561. Kosong **tetap kosong** di peserta. Uji: `TestPesertaDibulatkanSeperti7_7` (contoh
literal).

**⛔ Penyimpangan bertanggal terhadap ADR-U-0027 — HANYA di langkah 7.8.** Keputusan work owner: ikut Pega. Ketujuh
medan uang baris adjustment yang lahir saat Submit (`@divide(@toDecimal(…),1,4)` b3697–b3849) membaca sumber KOSONG
sebagai **0** — `[dugaan]` itulah yang `@toDecimal("")` Pega hasilkan; perilaku fungsi Pega itu tidak ada di korpus, yang
pasti keputusannya. `CURRENCY` b3868 bukan `@toDecimal` dan tetap teks apa adanya. Satu fungsi
konversi yang memuatnya: `services.kosongJadiNol78` (komentarnya menyatakan penyimpangan ini); ADR-U-0027 tetap berlaku di
seluruh tempat lain. Uji: `TestBarisPendaftaranKosongDibacaNol`.

### Tambahan tinjauan — 29 September 2026 (GILIRAN-15, `/code-review`)

- **Ralat rentang brief.** Brief GILIRAN-15 menyebut langkah 7.7 "b3600–b3671". Yang benar **b2744–b3631** (WHEN b3631);
  b3671 adalah awal 7.8. Kesepuluh medan 7.7 di atas dibaca dari rentang yang benar.
- **"Delapan medan" 7.8 = tujuh `@toDecimal` + `CURRENCY`.** Kosong = nol hanya berlaku bagi ketujuh medan uang;
  `CURRENCY` b3868 disalin sebagai teks (diuji dengan `USD`).
- Fungsi konversi dinamai ulang `kosongJadiNol78` (terikat langkahnya) dan dijaga satu pemanggil
  (`TestKosongJadiNolHanyaSatuPemanggil`); uji literal 7.7 kini membawa lebih dari empat desimal di SETIAP medan.
- **OQ-N12** *(untuk work owner — lahir dari tinjauan GILIRAN-15)* — dua jawaban 29-09-2026 bertemu dengan catatan 7
tiket 03. (1) N11: "`CLAIM_GROSS` sementara **kosong**". (2) "Kosong = nol di 7.8": `CLAIM_AMOUNT` baris yang lahir saat
Submit tidak pernah kosong (paling kecil `0`). (3) Catatan 7 (tiket 14): `CLAIM_GROSS` = `CLAIM_AMOUNT`, sehingga gerbang
`Save to RNM` 11.17.1 (`simpanrnm.go`, b6043 `.CLAIM_GROSS==""`) membaca `CLAIM_AMOUNT`. Akibatnya gerbang 11.17.1 **tidak
pernah menolak** baris yang lahir saat Register (dan putaran yang mewarisinya), sedangkan di Pega — `CLAIM_GROSS` tanpa
penulis — gerbang itu menolak. Perilaku TIDAK diubah executor. Pilih: (a) pertahankan catatan 7 (gerbang praktis mati untuk
baris 7.8); (b) pisahkan `CLAIM_GROSS` dari `CLAIM_AMOUNT` (kolom baru — migrasi — dan gerbang 11.17.1 akan menolak setiap
baris sampai penulisnya ada); (c) tunggu jawaban OQ-N11 dari pemilik ekspor.

## Keputusan bertanggal — 29 September 2026 (GILIRAN-17 paket 1: OQ-M4 ditutup) `[keputusan work owner 29-09-2026 — lembar keputusan, "rekomendasi"]`

Klaim XOL **tidak** membawa peserta dari berkas. Tombol `Claim Life - Upload CSV` di Register
(`FlowAction/UploadCSV_ClaimLife.xml:144`, `UploadCSVClaimLife_Act` b490–b539; pemetaan `SetClaimXOL_Act`
`InputRegisterClaimLife` b6763, 28 kolom b478–b1065, bendera XOL b1169, `LoadDataPeserta_Act` b282/b327) **tidak
dibangun**. Kolom yang dibawa berkas itu — uang, jendela valuasi, nama tertanggung — tetap dibaca ulang dari
PremiumList saat pendaftaran. Ini keputusan tercatat, bukan celah paritas.

## Keputusan bertanggal — 1 Oktober 2026 (OQ-N14): pencarian peserta memakai versi polis terakhir `[keputusan work owner 01-10-2026, "ikuti rekomendasi semua"]`

**Sebabnya.** Temuan B1 Endorsement Life (`bf45753`, `9160a8a`): Endorsement kini menulis versi baru polis ke
`M_LIFE_PREMIUM_DETAIL` — `PL_NUMBER` tetap nomor polis, `PL_NUMBER_EDM` = `<polis>/01`, `/02`, … *(rumus
`Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml` b84 `NOPOLIS||'/'||CARI14`)*, `EDMSTATUS`
`Old`/`New`/`Delete`/`Batal` — dan baris versi lama **tidak diubah**. Penyaring lama (`penyaringHidup` atas setiap
baris) meloloskan baris new business yang masih ber-`EDMSTATUS` NULL sesudah sertifikatnya di-`Delete` di `/01`: peserta
yang sudah dihapus endorsement tetap tampil dan dapat diklaim, dan `AmbilUntukKlaim` (tanpa `ORDER BY`) menyalin
baris versi mana saja.

**Aturannya** — satu tempat, `repository/pesertapolis.go` (`sqlVersiPeserta`, `peringkatVersi`, `pilihVersiHidup`):

1. Versi baris = angka sesudah `<PL_NUMBER>/` di `PL_NUMBER_EDM`; kosong — atau tidak berawalan nomor polisnya — = 0
   (new business).
2. Per `CERTIFICATE_NO` di satu `PL_NUMBER` hanya satu baris: `ROW_NUMBER() OVER (PARTITION BY CERTIFICATE_NO ORDER
   BY versi DESC, TGL_INPUT DESC NULLS LAST, LENGTH(ID) DESC, ID DESC) = 1`.
3. Baris terpilih ber-`TRIM(EDMSTATUS)` `Delete`/`Batal` → sertifikat tidak tampil di `Find Insured` dan ditolak
   pendaftaran dengan pesan yang sama dengan "tidak ada" (`… tidak ada atau sudah batal/delete`).
4. Jendela peringkat selalu di subkueri yang `WHERE`-nya `PL_NUMBER = :1` (`Cari`) atau `PL_NUMBER = :3 AND
   CERTIFICATE_NO = :4` (`AmbilUntukKlaim`) — nol pemindaian penuh atas 66,8 juta baris.

⚠️ **Dua penyimpangan dari rumus brief, disengaja.** (a) Rumus versi `REGEXP_SUBSTR(PL_NUMBER_EDM, '[0-9]+$')`
diperketat: ia membaca angka di akhir **nomor polis sendiri** sebagai versi bila baris new business ber-`PL_NUMBER_EDM`
= `PL_NUMBER` (polis berakhiran `-2024` → versi 2024, menang atas `/01`; kasus uji h). (b) `TGL_INPUT DESC` ditambah
`NULLS LAST` (Oracle menaruh NULL **pertama** pada `DESC`) dan `ID` dibandingkan secara **angka** (`LENGTH` dulu —
`ID` teks berisi `TO_CHAR(urutan.NEXTVAL)`, dan `'99' > '100'` secara teks). ⚠️ Tidak satu pun penulis korpus
(`SaveMasterLPDet`, PremiumList, Endorsement) mengisi `TGL_INPUT`; pada baris mereka seri versi praktis diputus `ID`.

**Inventori pembaca `M_LIFE_PREMIUM_DETAIL` di modul ini:**

| Berkas | Pembaca | Aturan baru? | Alasan |
| --- | --- | --- | --- |
| `repository/pesertapolis.go` | `Cari` (`sqlCariPeserta`) | **wajib** — diterapkan | daftar `Find Insured`: menentukan sertifikat yang dapat dipilih |
| `repository/pesertapolis.go` | `AmbilUntukKlaim` (`sqlAmbilPesertaKlaim`) | **wajib** — diterapkan | baris yang disalin ke klaim; server membaca ulang, klien hanya menyebut sertifikat |
| `repository/gandawarisan.go` | `DOBSumberKosong`, `StatusWarisanTerakhir`, `AdaWarisanSamaDOL`, `IsiTertanggungCermin` | tidak | berkunci `ID` baris sumber (`SOURCE_ID`) yang **sudah dipilih** `AmbilUntukKlaim` saat pendaftaran; membaca riwayat baris itu (Save to RNM), bukan memutuskan dapat-tidaknya diklaim |
| `repository/kolompeserta.go`, `models/klaimlife.go`, `services/simpanrnm.go` | — | tidak | hanya komentar; nol SQL ke tabel itu |

**SQL lama** (`Cari`; `AmbilUntukKlaim` sama bentuknya dengan `CERTIFICATE_NO = :2`, tanpa `ORDER BY`):

```sql
SELECT PL_NUMBER, POLICY_NO, CERTIFICATE_NO, NAME_OF_INSURED, CURRENCY, EDMSTATUS
  FROM POOLDATA.M_LIFE_PREMIUM_DETAIL
 WHERE PL_NUMBER = :1 AND (EDMSTATUS IS NULL OR TRIM(EDMSTATUS) NOT IN ('Batal','Delete'))
   [AND CERTIFICATE_NO LIKE '%'||:2||'%'] [AND UPPER(NAME_OF_INSURED) LIKE '%'||:3||'%']
 ORDER BY CERTIFICATE_NO FETCH FIRST n ROWS ONLY
```

**SQL baru `Cari`** — saringan sertifikat di dalam jendela (menyaring seluruh kelompok satu sertifikat; peringkat di
dalamnya tidak berubah), status dan nama di luar, atas versi terakhir saja: nama dari versi yang sudah diganti tidak
lagi menemukan pesertanya (kasus uji j).

```sql
SELECT PL_NUMBER, POLICY_NO, CERTIFICATE_NO, NAME_OF_INSURED, CURRENCY, EDMSTATUS
  FROM (SELECT PL_NUMBER, POLICY_NO, CERTIFICATE_NO, NAME_OF_INSURED, CURRENCY, EDMSTATUS,
               ROW_NUMBER() OVER (PARTITION BY CERTIFICATE_NO ORDER BY
                 NVL(TO_NUMBER(CASE WHEN SUBSTR(TRIM(PL_NUMBER_EDM), 1, LENGTH(PL_NUMBER) + 1) = PL_NUMBER || '/'
                   THEN REGEXP_SUBSTR(SUBSTR(TRIM(PL_NUMBER_EDM), LENGTH(PL_NUMBER) + 2), '^[0-9]+$') END), 0) DESC,
                 TGL_INPUT DESC NULLS LAST, LENGTH(ID) DESC, ID DESC) AS RN_VERSI
          FROM POOLDATA.M_LIFE_PREMIUM_DETAIL
         WHERE PL_NUMBER = :1 [AND CERTIFICATE_NO LIKE '%'||:2||'%'])
 WHERE RN_VERSI = 1 AND (EDMSTATUS IS NULL OR TRIM(EDMSTATUS) NOT IN ('Batal','Delete'))
   [AND UPPER(NAME_OF_INSURED) LIKE '%'||:3||'%']
 ORDER BY CERTIFICATE_NO FETCH FIRST n ROWS ONLY
```

**SQL baru `AmbilUntukKlaim`** (`<kolomSalin>` = 29 kolom, tidak berubah; bind `pl, sertifikat, pl, sertifikat`):

```sql
SELECT <kolomSalin> FROM POOLDATA.M_LIFE_PREMIUM_DETAIL
 WHERE PL_NUMBER = :1 AND CERTIFICATE_NO = :2
   AND ID = (SELECT ID FROM (SELECT ID, EDMSTATUS, <peringkat versi yang sama> AS RN_VERSI
                               FROM POOLDATA.M_LIFE_PREMIUM_DETAIL
                              WHERE PL_NUMBER = :3 AND CERTIFICATE_NO = :4)
              WHERE RN_VERSI = 1 AND (EDMSTATUS IS NULL OR TRIM(EDMSTATUS) NOT IN ('Batal','Delete')))
 FETCH FIRST 1 ROWS ONLY
```

**Kasus yang mengunci** (`repository/versiterakhir_test.go`; tabel yang sama dijalankan uji `db` atas Oracle,
`versiterakhir_db_test.go`, **MELEWATI** tanpa `ORACLE_DSN`):

| | Baris satu sertifikat | Hasil |
| --- | --- | --- |
| a | NB saja | tampil |
| b | NB + `/01 Old` | tampil **sekali**, nilai dari `/01` |
| c | NB + `/01 Delete` | tidak tampil, pendaftaran ditolak |
| d | NB + `/01 Batal` | tidak tampil, pendaftaran ditolak |
| e | NB + `/01 Delete` + `/02 New` | tampil dari `/02` |
| f | NB + `/01 " Delete "` (berspasi) | tidak tampil |
| g | dua baris `/01` (`Old` lebih baru, `Delete` lebih lama) | `TGL_INPUT` terbaru menang → tampil |
| h | polis berakhiran angka, NB ber-`PL_NUMBER_EDM` = `PL_NUMBER`, + `/01 Delete` | tidak tampil (versi NB = 0) |
| i | dua baris `/01` tanpa `TGL_INPUT`, `ID` 99 `Delete` dan 100 `Old` | `ID` terbesar **secara angka** → tampil |
| j | nama lama di NB, nama baru di `/01 Old`, cari nama lama | tidak tampil |

Uji gigit: pemeriksa yang sama dijalankan atas SQL termutasi — tanpa `ROW_NUMBER` kasus (c) merah; tanpa `PL_NUMBER`
di jendela penjaga "nol pemindaian penuh" merah. Mutasi kode produksi (tanpa `ROW_NUMBER`, tanpa `PL_NUMBER` di jendela
`Cari`/`AmbilUntukKlaim`, rumus versi naif, tanpa `NULLS LAST`, status di dalam jendela) — keenamnya merah.

⚠️ **Tiruan bersama.** `uji/skemauji` (di luar folder modul ini) belum memuat `PL_NUMBER_EDM` dan `TGL_INPUT`; uji `db`
Claim Life menambahkannya ke **tiruan skema uji** (`ALTER TABLE <skema uji>.M_LIFE_PREMIUM_DETAIL ADD …`, ORA-01430
ditoleransi) — bukan DDL tabel warisan. Pemilik `uji/skemauji` sebaiknya memasukkan keduanya ke `ddlTiruanPesertaPolis`.

**Di luar cakupan keputusan ini** *(dicatat, tidak diubah)*: klaim yang **sudah terdaftar** sebelum endorsement
menghapus sertifikatnya tetap memakai baris sumber (`SOURCE_ID`) yang dipilih saat pendaftaran; nol pemeriksaan ulang
di Save to RNM. Bila bisnis menghendakinya, itu keputusan baru. Tiket 11 Endorsement Life (di luar folder ini)
**tidak** disunting — temuan B1 tertutup dari sisi pembaca di sini.

