# 08: Kontrak hilir — rekam premium yang dikonsumsi Claim Life

**Status:** sebagian — uji kontrak HTTP terhadap Oracle dan index `PL_NUMBER` di `T_PREMIUM_LIST_DETAIL` belum ada; `M_LIFE_PREMIUM_SUMMARY` **ditulis** berkunci `PL_NUMBER` (OQ-PL-09 ditutup GILIRAN-18)

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 05a (rekam summary), 05b (alur simpan polis — penulisan detail NB menumpang di sana)

## Hasil & nilai pengguna

Sebagai **admin klaim Life**, saya ingin menemukan premium list dan peserta yang benar untuk klaim
yang sedang saya proses — berkunci `PL_NUMBER` — **segera setelah** polis
tersimpan, tanpa menunggu proses terjadwal apa pun, supaya klaim tidak pernah tertahan hanya karena
baris pesertanya belum sempat ditulis. *(User story 39–40 di spec; **ADR-0001**)*

## Area codebase

`internal/repository` (kueri pencarian premium summary + detail; penulisan detail),
`internal/services` (pemanggilan **inline** penulisan detail di alur simpan polis; kontrak yang
dipanggil konteks Claim Life), `internal/handlers` (endpoint pencarian), penandaan **kontrak lintas
konteks** di tempat bentuk rekam didefinisikan.

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InsertLifePremiumDetail_act` | `ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertLifePremiumDetail_act.xml` (286.827 byte, `pxUpdateDateTime` `20260211T064342.645 GMT`, ruleset `01-01-91`) | **penulis detail NB** |
| `SaveMasterLPDet` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/SaveMasterLPDet.xml` **dan** `Endorsement Life/RDBList/SaveMasterLPDet.xml` | `INSERT INTO POOLDATA.M_LIFE_PREMIUM_DETAIL` + `COMMIT;` (baris 252) |
| `InsertJsonPolisLife_Act` (EDM) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` | pemicu **inline** detail EDM (step **11.6**) |
| `GetPesertaClaim_sql1` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/GetPesertaClaim_sql1.xml` | **konsumen** |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/` **dan** `Endorsement Life/RDBList/` | penulis summary |

`[terverifikasi]` **`SaveMasterLPDet` adalah satu rule yang sama** untuk NB dan EDM, bukan dua rule
serupa: identitas identik, `pxUpdateDateTime` identik (`20260211T064412.717 GMT`), ukuran identik
(15.537 byte), dan diff atas dua ekspor yang dinormalisasi menyisakan **8 baris** — seluruhnya cap
waktu ekspor. **Penulisnya sudah seragam; yang berbeda hanyalah pemicunya.**

`[terverifikasi]` Kueri konsumen `GetPesertaClaim_sql1`: `SELECT * FROM POOLDATA.M_LIFE_PREMIUM_DETAIL`
dengan `WHERE PL_NUMBER = {pyWorkPage.PolicyDataLife.PremiumListSummary.PL_NUMBER}`, ditambah
pencocokan sebagian (`LIKE`) atas `CERTIFICATE_NO` dan `UPPER(NAME_OF_INSURED)`.

`[terverifikasi]` Kolom `INSERT` `SaveMasterLPDet` memuat **`PL_NUMBER` dan `PL_NUMBER_EDM`**,
`CERTIFICATE_NO`, `NAME_OF_INSURED`, `POLICY_NO`, `CURRENCY`, seluruh kolom uang gross/`*_REFUND`/
`*_RETRO`, `RATE`, `PRORATETYPE`, `SUM_AT_RISK_GROSS`, `SUM_AT_RISK_RETRO`, `RETROCEDED_SHARE`; PK
dari `M_LIFE_PREMIUM_DETAIL_SEQ.nextval`.

`[terverifikasi]` Class integrasi `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` adalah **tulang punggung
domain Life**: Claim Life 46 rule, Endorsement Life 25, PremiumList Life 19, Komite Claim Life 11.

### Peta langkah `InsertLifePremiumDetail_act` `[terverifikasi]`

Penomoran dari `<pyStepPageReference>` — akarnya **`RH_2`**, bukan `RH_1`:

| Step | Langkah | Catatan |
| --- | --- | --- |
| 1 | `Property-Set` | set `Param.pyReportName = "SelectNoJsonPolis_RD"`, `Param.pyReportClass = "ASM-FW-GISFW-Work-LIFE"` |
| **2** | `Call Rule-Obj-Report-Definition.pxRetrieveReportData` | ⚠️ **pemicu batch** — mengambil **daftar** kasus |
| 3 | *(loop `hasil.pxResults`)* | |
| 3.1 | `Property-Set` | |
| 3.2 | `Obj-Open-By-Handle` | buka work object per baris hasil |
| 3.3 | "Insert to table detail" | |
| 3.3.1 | `Page-Remove` | |
| 3.3.2 | `Property-Set` "get ceding co name" | |
| 3.3.3 | `Property-Set` "insert nilai dari data-batch → int" | precondition `TempError.CARIDESC==1` (baris 3634) |
| **3.3.4** | `RDB-List` → **`SaveMasterLPDet`** (baris 3730) | precondition **`hasilDetail.pxResults(1).PL_NUMBER==""`** (baris 3820) — **penjaga idempotensi** |
| 3.4 | `Property-Set` "Pega to jsondata" | |
| 3.5 | `Obj-Save` | |
| **3.6** | `Commit` | **AKTIF** |

`[terverifikasi]` **Nol `<pyStepsBlockName>` di berkas ini** — tidak ada langkah ter-remark.

`[terverifikasi]` Report Definition `SelectNoJsonPolis_RD` (class `ASM-FW-GISFW-Work-LIFE`) **dirujuk
tetapi tidak ada berkasnya di korpus** — asimetri rujukan; pemilih batch itu sendiri tidak terekspor.

## ⚠️ Penajaman kontrak hilir — **satu tabel, dua sudut pandang** `[keputusan work owner]`

Ditambahkan 2026-09-15 dari **verdict V14** grilling Endorsement Life.

Konteks Endorsement Life menulis **baris bernilai negatif** (jurnal balik) ke tabel yang **sama**
dengan yang ditulis dan dibaca di sini. Aturannya:

| Pembaca | Melihat |
| --- | --- |
| **Akuntansi / ringkasan premium** | **seluruh** baris — positif **dan** negatif; nettonya dari penjumlahan |
| **Klaim** | **hanya peserta hidup** — yang belum dibatalkan dan belum dihapus |

`[terverifikasi]` Penandanya adalah kolom **`EDMSTATUS`** pada `M_LIFE_PREMIUM_DETAIL`, terbaca dari
daftar `INSERT` di `SaveMasterLPDet`: ia diisi dari `TempValue.EDMStatus` dengan nilai
`Old` / `New` / `Delete` / `Batal`. Kolom `STATUS` **bukan** penandanya — ia berisi `0` untuk
`QR`/`QP` dan `1` untuk `TP`/`TR` (jenis transaksi).

⚠️ `[terverifikasi]` **Jalur new business — yakni tiket ini — tidak mengisi `EDMSTATUS` sama
sekali.** Sensus `InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` /
`INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY`): **nol** kemunculan. Baris NB karena itu masuk
dengan `EDMSTATUS` kosong/NULL — dan **harus tetap terlihat** oleh klaim.

`[terverifikasi]` Kueri klaim hari ini **belum menegakkan** kontrak ini: `GetPesertaClaim_sql1`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL`) menyaring hanya
dengan `PL_NUMBER`, `CERTIFICATE_NO`, dan `NAME_OF_INSURED` — **tanpa** penyaring status. ⚠️ Apakah
Pega menyaring di lapisan lain **tidak terbukti dari korpus**; ini **kontrak yang wajib ditegakkan
sistem baru**.

**Penegakannya milik konteks Claim Life** (spec §16, AC 25–30; tiket 02 Claim Life). Tiket ini
mengikatnya sebagai **syarat kontrak** yang harus terbukti lewat test kontrak lintas konteks.

## ⚠️ Penyimpangan sadar — tanpa job, penulisan INLINE `[keputusan work owner]`

| | Pega existing | Sistem baru |
| --- | --- | --- |
| **Logika penulisan NB** | `InsertLifePremiumDetail_act` | **TETAP** `InsertLifePremiumDetail_act` — tidak diganti, tidak ditulis ulang |
| **Pemicu NB** | **job/batch** (step 2 `pxRetrieveReportData` + loop step 3) | **INLINE saat proses insert/simpan polis** |
| **Pemicu EDM** | inline di alur simpan (step 11.6) | inline di alur simpan — **tidak berubah** |

**Yang disamakan adalah TIMING, bukan logikanya.** Step 2 dan loop step 3 **tidak dimigrasikan** —
keduanya semata mesin batch. Yang dimigrasikan adalah **badan per-kasus** (3.1–3.6), dipanggil sekali
untuk kasus yang sedang disimpan.

**Konsekuensi positif:** data peserta langsung tersedia untuk klaim **tanpa jeda job**, dan NB
seragam dengan EDM.

## ADR terkait

**ADR-0001** (batas konteks Claim Life ↔ hulu; perubahan bentuk rekam = perubahan kontrak),
**ADR-0011** (bentuk `PremiumListSummary` / `PremiumListDetail` yang dipakai mesin status klaim),
**ADR-0003** (uang non-float menyeberang batas), **ADR-0015** (batas transaksi dipegang Go —
`SaveMasterLPDet` commit sendiri, jadi ia titik potong).

## Acceptance criteria

- [ ] Baris peserta **new business** ditulis ke `M_LIFE_PREMIUM_DETAIL` lewat logika
      `InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT`),
      dipanggil **INLINE sebagai bagian alur simpan polis** — **bukan** oleh job, cron, worker
      terjadwal, atau antrean tunda. Test yang menemukan penjadwal di jalur ini **gagal**.
      `[keputusan work owner]` — belum: inline di `services/polis_summary.go:simpanDalam`, tetapi uji yang gagal bila ada penjadwal belum ada
- [ ] **Setelah simpan polis NB berhasil**, baris `M_LIFE_PREMIUM_DETAIL` untuk `PL_NUMBER` itu
      **sudah ada** — dibuktikan dengan membacanya **segera** sesudah respons simpan, tanpa menunggu
      apa pun. — belum: struktural (satu commit), belum dibuktikan dengan membaca terhadap Oracle
- [x] Baris itu **dapat dibaca jalur baca klaim**: kueri bergaya `GetPesertaClaim_sql1` — berkunci
      `PL_NUMBER`, dengan pencocokan sebagian pada `CERTIFICATE_NO` dan `NAME_OF_INSURED` **tanpa
      peduli huruf besar/kecil** — mengembalikannya. — bukti: `repository/pesertapolis.go:sqlCariPeserta`; uji `TestKolomBacaClaimLifeDiisiPenulisPremiumList`, `TestPenulisWarisanMengisiKunciBacaKlaim`
- [x] `SaveMasterLPDet` dipakai **apa adanya** sebagai penulis bersama — tidak dibuatkan salinan,
      tidak divariasikan per jalur. Kolom `PL_NUMBER_EDM` tetap ada di skema dan **dibiarkan kosong**
      oleh jalur new business. *(pengisiannya milik konteks Endorsement Life)* — bukti: `repository/polis_warisan.go:kolomPesertaWarisan`; uji `TestKolomWarisanVERBATIMDariSaveMasterLPDet`
- [ ] Penjaga idempotensi dipertahankan: penulisan hanya terjadi bila baris untuk `PL_NUMBER` itu
      belum ada (setara precondition `PL_NUMBER==""` pada step 3.3.4). Menyimpan ulang polis yang
      sama **tidak** menggandakan baris — dibuktikan dengan menyimpan dua kali lalu menghitung baris. — belum: menyimpang sadar — salinan milik work yang sama diganti, bukan dilewati (`[terbuka — work owner]`); nol uji simpan-dua-kali terhadap Oracle
- [ ] `SaveMasterLPDet` diperlakukan sebagai **titik potong transaksi** (ia `COMMIT;` sendiri, baris
      252): dipanggil **setelah** transaksi penomoran + summary commit, konsisten dengan aturan
      urutan transaksi campuran. *(AC 20–21, 24 spec)* — belum: pl2 — salinan berada di dalam satu transaksi simpan; titik potong ditiadakan
- [ ] Kegagalan penulisan detail **tidak** membatalkan premium list yang sudah tersimpan; keadaannya
      **terdeteksi** dan pemanggilan ulang aman berkat penjaga idempotensi. *(AC 23 spec)* — belum: pl2 — kegagalan salinan membatalkan seluruh simpan (satu transaksi), bukan dibiarkan
- [ ] Rekam `M_LIFE_PREMIUM_SUMMARY` dan `M_LIFE_PREMIUM_DETAIL` jalur new business dapat ditemukan
      lewat `PL_NUMBER`. *(AC 28 spec)* *(pencarian lewat `PL_NUMBER_EDM` diuji di konteks Endorsement Life)* — belum: keduanya kini ditulis berkunci `PL_NUMBER` (`M_LIFE_PREMIUM_SUMMARY` sejak GILIRAN-18, `SummaryWarisan.Ganti`); pencariannya baru terbukti di uji `db` `TestSummaryWarisanDitulisSepertiProsedur`, yang SKIP tanpa Oracle
- [ ] Nilai uang ditulis dan dibaca sebagai **desimal presisi arbitrer**; nilai yang ditulis hulu
      dibaca hilir **identik**, tanpa pembulatan di perbatasan. *(AC 16 spec; **ADR-0003**)* — belum: angka warisan dikirim sebagai teks yang bergantung NLS sesi; nol uji pulang-pergi
- [x] Bentuk kedua rekam ditandai di kode sebagai **kontrak lintas konteks**; mengubahnya memaksa
      pembaruan sadar di sisi Claim Life. *(AC 29 spec; **ADR-0001**)* — bukti: uji `TestKolomBacaClaimLifeDiisiPenulisPremiumList`, `TestPolicyDataLifeSamaDiKeduaSisi`
- [x] ⚠️ **Jalur baca klaim menyaring peserta batal/delete.** Pembacaan peserta bergaya
      `GetPesertaClaim_sql1` **tidak menampilkan** baris ber-`EDMSTATUS` `'Batal'` atau `'Delete'`
      — hanya **peserta hidup**. `[keputusan work owner]` *(verdict V14 grilling Endorsement Life;
      AC 25 spec Claim Life)* — bukti: `repository/pesertapolis.go:sqlCariPeserta` (`penyaringHidup`); uji `TestPesertaHidupMenyaringBatalDanDelete`
- [x] Peserta **new business** tetap muncul di jalur baca klaim meski jalur NB **tidak mengisi**
      `EDMSTATUS` (kosong/NULL). ⚠️ Penyaring naif `NOT IN ('Delete','Batal')` membuang seluruh
      peserta NB di Oracle — test wajib memuat kasus ini. *(AC 26 spec Claim Life)* — bukti: uji `TestPesertaNBTetapHidupDiJalurBacaKlaim`
- [x] **Jalur baca akuntansi/ringkasan premium TIDAK menyaring** — ia melihat **seluruh** baris,
      positif maupun negatif, karena nettonya diperoleh dari penjumlahan. **Satu tabel, dua sudut
      pandang**, dan itu disengaja. *(AC 49 spec Endorsement Life)* — bukti: `repository/polis_summary.go:sqlBarisUangPolis` (tanpa penyaring `EDMSTATUS`); uji `TestPenyaringPesertaHanyaSatuTempat`
- [ ] Ada test kontrak yang menembus dari simpan polis sampai pencarian bergaya Claim Life — **satu — belum: skema uji belum memasang 050–056; uji kontrak lewat HTTP belum ada

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

- [ ] ⚠️ Claim Life membaca peserta dari **`T_PREMIUM_LIST_DETAIL`**, bukan dari
      `M_LIFE_PREMIUM_DETAIL` maupun dari CLOB JSON. Kontrak bacanya tetap sama bentuknya.
      *(AC 33 spec; penyimpangan sadar 1)* — belum: dilampaui pl2 — Claim Life tetap membaca `M_LIFE_PREMIUM_DETAIL`
- [ ] ⚠️ Aturan **peserta hidup** (`EDMSTATUS` bukan `Delete`/`Batal`, NULL tetap muncul) berlaku
      **apa adanya** pada tabel baru — pindah tabel **tidak** mengubah aturannya. — belum: jalur baca klaim belum membaca tabel baru (pl2)
- [ ] ⚠️ Kolom `PL_NUMBER` pada tabel peserta **ber-index** sejak hari pertama — ia kunci baca Claim
      Life pada tabel berjutaan baris. *(AC 49 spec)*
      seam**, API HTTP, terhadap skema uji Oracle. — belum: migrasi `052` hanya meng-index `PREMIUM_LIST_ID` dan `PARENT_ID`; `PL_NUMBER` tanpa index

## Blocker

**Tidak ada.** **OQ-068 ditutup 2026-09-15** — `[terverifikasi]` + `[keputusan work owner]`.

Sebelumnya tiket ini `needs-info` karena korpus tidak memperlihatkan penulis `M_LIFE_PREMIUM_DETAIL`
di jalur new business — hanya `SaveMasterLPDet` di Endorsement. Work owner kemudian menambahkan
`InsertLifePremiumDetail_act` **dan** salinan `SaveMasterLPDet` ke `PremiumList Life/`, lalu
menetapkan bahwa logikanya dipakai apa adanya sementara **pemicunya** diubah dari job menjadi inline.
Jalur NB kini terbukti dan seragam dengan EDM.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## Implementasi — 28-09-2026 (giliran 10): kontrak dua sisi, dikunci dan dibuktikan menggigit

### Dua kontrak, masing-masing diuji di kedua sisi

| Kontrak | Sisi penulis | Sisi pembaca | Uji |
| --- | --- | --- | --- |
| **pl4** `GET /api/polis-life/ringkas` | `services.PolicyDataLife` (tag JSON, dibaca lewat `json.Marshal`) | `frontend/src/services/api.ts` `interface PolicyDataLife` | Go membaca antarmuka TS (`TestPolicyDataLifeSamaDiKeduaSisi`); TS membaca tag Go (`policydatalife.kontrak.test.ts`), dengan daftar kunci bertipe `Record<keyof PolicyDataLife, true>` sehingga `tsc` menolak daftar yang menyimpang dari antarmukanya |
| **pl2** `M_LIFE_PREMIUM_DETAIL` | `repository/polis_warisan.go` (PremiumList) | `repository/pesertapolis.go` `kolomSalin` + `sqlCariPeserta` (Claim Life) | `TestKolomBacaClaimLifeDiisiPenulisPremiumList`: setiap kolom yang dibaca Claim Life harus diisi penulis |

⚠️ **Instrumen dibuktikan menggigit** — mutasi sementara, lalu dipulihkan (`git checkout`, diff kosong):

- tag `prodKe` → `prodke` di Go: uji Go **merah** ("kunci PolicyDataLife berselisih"), uji TS **merah**;
- baris `CURRENCY` dicabut dari `kolomPesertaWarisan`: uji kontrak **merah** — *"membaca [CURRENCY] … penulis
  PremiumList tidak mengisinya"*.

Pembaca daftar pilih diuji atas jawaban yang sudah diketahui (ujung `ID…CURRENT_AGE` dan
`PL_NUMBER…EDMSTATUS`), supaya kontrak tidak lulus hampa karena pembacanya rusak.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| penulisan detail NB INLINE, bukan job | ✅ `simpanDalam` di transaksi simpan/tutup (tiket 05b); nol penjadwal |
| sesudah simpan berhasil, baris untuk `PL_NUMBER` itu sudah ada | ✅ struktural — satu commit · ⚠️ belum dibuktikan terhadap Oracle |
| dapat dibaca jalur baca klaim | ✅ kontrak kolom + kunci (`PL_NUMBER`, `CERTIFICATE_NO`) |
| `SaveMasterLPDet` apa adanya; `PL_NUMBER_EDM` kosong di NB | ✅ pemetaan 80 kolom VERBATIM; `PL_NUMBER_EDM ← p.PL_NUMBER_EDM` (kosong di NB) |
| penjaga idempotensi | ⚠️ **menyimpang, sadar**: Pega *melewati* bila baris untuk `PL_NUMBER` sudah ada (`PL_NUMBER==""` b3820); kami *mengganti* salinan milik work yang sama (`DELETE … PL_NUMBER AND IDPEGA`). Keduanya tidak menggandakan; yang kami pilih membuat unggah ulang sesudah simpan ikut tercermin. `[terbuka — work owner]` bila "lewati" yang dikehendaki |
| titik potong `SaveMasterLPDet` | ➖ lenyap — satu transaksi (pl2) |
| peserta NB (EDMSTATUS NULL) tetap hidup di jalur baca klaim | ✅ `TestPesertaNBTetapHidupDiJalurBacaKlaim` |
| bentuk rekam ditandai kontrak lintas konteks | ✅ kedua uji di atas; mengubah salah satu sisi memerahkan uji |
| `M_LIFE_PREMIUM_SUMMARY` dapat ditemukan lewat `PL_NUMBER` | ⚠️ ditulis sejak GILIRAN-18 (OQ-PL-09 ditutup, tiket 05a); bukti Oracle = uji `db` SKIP |
| ⚠️ AC relasional "Claim Life membaca dari `T_PREMIUM_LIST_DETAIL`" (2026-09-16) | ⚠️ **dilampaui pl2** (28-09-2026): Claim Life tetap membaca `M_LIFE_PREMIUM_DETAIL`, maka tabel itu ditulis |
| test kontrak menembus satu seam lewat HTTP terhadap Oracle | ⚠️ belum — skema uji belum memasang 050–056 |

### Angka

Go **544 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **349** · tsc bersih.
