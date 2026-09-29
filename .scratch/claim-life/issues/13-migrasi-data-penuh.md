# 13: Migrasi data penuh Claim — Life

**Status:** sebagian — migrasi data belum punya pelari dan belum pernah dijalankan; posisi tahap dan `LAYER_*` tidak dibawa

**Blocked by:** 11 (kontrak Komite — jalur balik), **14 (skema relasional klaim — PREFACTOR)**

> ⚠️ **Dipersempit 2026-09-16 — revisi penyimpanan.** **Pemindahan bentuk data klaim
> (JSON + tabel flat warisan → enam tabel relasional, beserta normalisasi atribut polis) berpindah
> ke tiket 14**, yang menjadi PREFACTOR. Tiket ini **tinggal** menangani hal yang tidak menyentuh
> bentuk: keutuhan riwayat setelah cutover, penomoran, pembacaan master, dan kolom yang tidak
> ditafsirkan. AC yang berpindah ditandai di bawah.

## Hasil & nilai pengguna

Sebagai **work owner**, saya ingin seluruh data klaim Life dipindahkan — beserta **semua barisnya** —
sehingga tidak ada pekerjaan tertinggal di Pega dan riwayat putaran Komite tidak hilang.
*(User story 40, 41, 42 di spec)*

## Area codebase

`internal/repository` (pembacaan sumber + penulisan sasaran), skrip migrasi + DDL produksi di dalam
`OUTPUT_HASIL_RNM/`. Tidak menyentuh `handlers`/`frontend`.

## Rule Pega sumber

| Rule | Identitas | Dipakai untuk |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` 55 nama kolom `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` terbaca langsung dari blok PL/SQL |
| `Claim Life/RDBList/GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL` | penomoran yang tidak boleh melompat/mengulang setelah migrasi |

`[terverifikasi]` Tabel akseptasi Life **terpisah** dari lini lain: `OS_AKSEPTASI_KLAIM_LIFE`
(7 kemunculan, hanya `Claim Life` + `Komite Claim Life`) versus `OS_AKSEPTASI_KLAIM` (191 kemunculan,
enam modul non-Life). Migrasi ini **tidak menyentuh** tabel non-Life.

## ADR terkait

**ADR-0009** (migrasi penuh; koeksistensi **ditolak**; cutover memutus, bukan bertahap),
**ADR-0011** (seluruh **baris** ikut pindah, bukan hanya keadaan terakhir),
**ADR-0006** (penomoran), **ADR-0007** (jejak audit tidak dapat direkonstruksi ke belakang).

## Acceptance criteria

⚠️ **Berpindah ke tiket 14** (jangan dikerjakan di sini): seluruh klaim + seluruh baris adjustment
terbawa; jumlah baris per klaim utuh; presisi uang; desimal presisi arbitrer; migrasi dapat
dijalankan ulang; normalisasi atribut polis. Rujukan silang: AC 51–54 spec.

- [ ] Klaim yang sedang berada di tengah siklus terbawa **beserta statusnya** dan posisi tahapnya.
      ⚠️ **Diselaraskan:** posisi tahap mendarat di **`T_WORK_CLAIM`** (spec §2b), bukan di header
      klaim. *(AC 46 spec; penyimpangan sadar 6)* — belum: posisi tahap tidak ada di sumber warisan — migrasi melaporkannya, tidak membawanya (`TestTahapSelaluDilaporkanTidakAdaDiSumber`); migrasi data belum pernah dijalankan
- [ ] Setelah migrasi, penomoran klaim **tidak melompat dan tidak mengulang**. *(AC 42 spec)* — belum: `PeriksaTandaAir` hanya menjaga "tidak mundur" (`TestPenomoranTidakMundurSesudahMigrasi`) dan nol pemanggil; migrasi data belum pernah dijalankan
- [ ] Tidak ada periode dua penulis ke rekam akseptasi Life dari sisi Claim — Life. *(ADR-0009)*
      ⚠️ **Diselaraskan:** setelah cutover, penulis satu-satunya adalah keenam tabel klaim baru;
      `OS_AKSEPTASI_KLAIM_LIFE` **tetap ditulis** — ⚠️ **koreksi 2026-09-16**
      `[keputusan work owner]`: setelah cutover, penulis klaim adalah **kedelapan tabel relasional
      baru** *dan* `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE`, karena hilir masih membaca dari sana.
      Yang **dibuang hanya JSON**. Test yang **menolak** penulisan `OS_AKSEPTASI_KLAIM_LIFE` justru
      **gagal**. *(AC 32 spec)* — belum: "nol periode dua penulis" adalah sifat cutover, tidak terbukti dari kode; yang terverifikasi hanya jalur tulis datar (`TestJalurTulisDatarWarisanWAJIBADA`)
- [ ] ⚠️ Master `RATE_LIFE`, `PRODUCTINWARD_LIFE`, `CURRENCY` diperlakukan sebagai **view atas
      `JSONDATA`**, bukan sebagai tabel relasional biasa. **Pengecualian 2026-09-16:** **produk**
      dibaca dari **`product_life` relasional** (hasil migrasi Master Product Name Life), bukan dari
      `m_product_life.JSONDATA`. *(AC 38 spec)* — belum: master tidak disentuh migrasi (`TestMasterViewTidakDisentuh`), tetapi pembaca `RATE_LIFE`/`PRODUCTINWARD_LIFE` belum ada dan produk tidak dibaca dari `product_life` relasional di mana pun (OQ-M7) — *ralat 29-09-2026 (GILIRAN-17): pembaca sempit keduanya kini ada (`ambangproduk.go` butir bh, `ratelife.go` OQ-M7); `OUTWARDRATEID` tetap terbuka — DBA*
- [x] `NO_SEQ` terjaga **per kombinasi `(CLASS, JENIS, TAHUN)`**, bukan sebagai penghitung global. — bukti: `APP_RNM/internal/repository/penomor.go:Penomor.UrutNomorBerikut` (kunci `CLASS`, `JENIS`, `TAHUN`), `APP_RNM/internal/repository/migrasinomor.go:TandaAirSequence`; uji `TestTandaAirDihitungPerKombinasiKunci`
- [ ] Kolom `LAYER_1`…`LAYER_4` dipindahkan apa adanya **tanpa ditafsirkan** — perannya belum
      terverifikasi. — belum: migrasi Claim Life sengaja TIDAK membawa `LAYER_1`..`LAYER_4` (`kolomTidakDibawa`, uji `TestKolomTakDibawaHanyaAdaDiKatalog`) — nilainya tinggal di tabel datar, tidak dipindahkan

## Catatan penutupan (2026-09-14)

**OQ-001 TERTUTUP** untuk Claim — Life `[data DBA]` — DDL **12 tabel/view** diserahkan, mencakup
**seluruh** persistensi Claim Life. Daftar lengkap dan tipenya di `discovery/open-questions.md`
OQ-001. Yang mengikat tiket ini:

| Temuan | Konsekuensi untuk migrasi |
| --- | --- |
| Kolom uang = Oracle **`NUMBER` tanpa presisi** | Go **wajib** desimal presisi arbitrer; **`float64` dilarang** (**ADR-0003**) |
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

## Catatan

`[keputusan work owner 2026-09-14]` Opsi koeksistensi — sistem baru hanya menerima klaim baru
sementara klaim berjalan diselesaikan di Pega — **ditolak**. Penolakan itu dicatat di **ADR-0009**
supaya tidak diusulkan ulang.

`[keputusan work owner]` Jejak audit lama **tidak dapat direkonstruksi**: data sebelum cutover hanya
punya `CREATEOPNAME` + empat kolom tanggal. Riwayat transisi lengkap hanya ada setelah cutover.

## Perintah verifikasi

```
go test ./internal/...
make check
```

Tambah verifikasi khusus migrasi: hitungan klaim dan hitungan baris per klaim sebelum dan sesudah
pemindahan harus sama.

---

## Pembacaan ulang XML dan katalog — 27 September 2026

**Sensus 1 — apakah posisi tahap ada di sumber migrasi? Dua cara, keduanya NOL.**

| Cara | Jendela | Hasil |
| --- | --- | --- |
| 1 | 62 nama kolom `kolomWarisan` yang berbau `POSITION`/`STAGE`/`TAHAP`/`ASSIGN`/`WORK`/`ROLE` | **0** |
| 2 | berkas `Claim Life/**/*.xml` yang menyebut `OS_AKSEPTASI_KLAIM_LIFE` **dan** `pyPosition` | **0** |

⛔ **Posisi tahap TIDAK ADA di tabel datar warisan.** Ia hidup di work object Pega, yang bukan
bagian dari 12 tabel/view OQ-001. Jadi AC *"terbawa beserta statusnya dan posisi tahapnya"*
terpenuhi separuh: **status** terbawa (`STS_REJECT`), **tahap tidak dapat** — bukan karena
terlewat, melainkan karena tidak ada di sumbernya. Menurunkannya dari `STS_REJECT` berarti
**merancang**, dan ADR-U-0002 sudah menyatakan model peran DIRANCANG bukan dimigrasikan. Itu
keputusan work owner, bukan kelengkapan skrip.

**Sensus 2 — kolom katalog lawan medan struct.**

⚠️ **Ralat penamaan:** ini BUKAN "dua cara" dalam arti bab 4a — dua wadah yang dikurangkan bukan
dua cara menghitung satu hal, dan menyebutnya begitu adalah trap 5. Yang benar-benar diperiksa dua
arah adalah **daftarnya**: tiap nama wajib ADA di katalog dan wajib TIDAK ADA di sumber produksi.

| Cara | Hasil |
| --- | --- |
| katalog `kolomWarisan` | **62** kolom |
| medan struct `BarisLama` | **55** medan |
| selisih | **7**: `STS_KONVERSI`, `TGL_KONVERSI`, `LAYER_1`…`LAYER_4`, `INDEXLIST` |

⭐ Ketujuhnya **sengaja tidak dibawa**: model relasional tidak punya rumah bagi mereka, dan
mengarang rumah berarti **menafsirkan** kolom yang perannya belum terverifikasi. Yang dilakukan
migrasi: tidak menyentuhnya. Tabel datarnya tetap ada (AC 3), jadi nilai lamanya utuh di tempatnya.

## Implementasi — 27 September 2026 (tiket 13)

**Status: `claimed`** — **2 dari 6 AC tertutup.** Titik tetap `b0afeb4`.

⛔ **Yang TIDAK dijalankan, dan tidak akan dijalankan tanpa Anda:** migrasi datanya sendiri.
Menjalankan migrasi data menuntut **persetujuan manusia**, begitu pula membaca tabel produksi mana
pun. Yang ada di tiket ini **aturan dan pagarnya**, seluruhnya dapat diuji tanpa Oracle. OQ-018
menyebutnya sama: *"Tidak memblokir penulisan skrip migrasi; memblokir pemilihan sumber data saat
eksekusi."*

| Angka | Nilai | Cara 1 | Cara 2 | Sepakat? | Label |
| --- | --- | --- | --- | --- | --- |
| AC tertutup | **2 dari 6** | `grep -c '^- \[x\]'` → 2 | `grep -cE '^- \[[x ]\]'` → 6 dikurangi `grep -c '^- \[ \]'` → 4 | ✅ | `[terverifikasi]` |
| test Go | **228 PASS · 0 FAIL · 34 SKIP** *(dari 222/34)* | cacah awalan `--- PASS`/`--- FAIL`/`--- SKIP` dari `go test -tags=db ./internal/... -v` → **262** | `grep -rhoE '^func Test[A-Za-z0-9_]+' internal/ --include=*_test.go \| wc -l` → **262** | ✅ | `[terverifikasi]` |

| Berkas | Isi |
| --- | --- |
| `repository/migrasinomor.go` *(baru)* | `KunciSequence`, `TandaAirSequence`, `PeriksaTandaAir` |
| `repository/migrasinomor_test.go` *(baru)* | 3 kasus |
| `repository/migrasibatas_test.go` *(baru)* | 4 penjaga, dua di antaranya **terbalik** |

**Penjaga, tiap-tiap dibuktikan gagal lalu dipulihkan:**

| Mutasi | Hasil |
| --- | --- |
| `LAYER_1` dibandingkan di kode | merah — *"menafsirkan kolom LAYER"* |
| `LAYER_1..4` ditambahkan ke struct `BarisLama` | merah — *"BarisLama kini membawa LAYER_1"* |
| `INSERT INTO POOLDATA.RATE_LIFE` disisipkan | merah — master itu VIEW atas `JSONDATA` |
| `SELECT JSONDATA FROM m_product_life` disisipkan | merah — produk dibaca dari `product_life` relasional |
| jalur tulis datar warisan dibuang dari `Simpan` | merah — penjaga **TERBALIK** |
| penjagaan kunci komposit dilucuti | merah — tanda air jadi global |

⛔ **Dua mutasi semula gagal KOMPILASI, dan itu bukan lulus.** `b.LAYER_1` tidak dapat dikompilasi
sebab medannya memang tidak ada — yang justru menyingkap sensus 2 di atas. Setelah medannya
ditambahkan, mutasinya kompilasi dan penjaganya merah. Diam karena kompilasi gagal terbaca persis
seperti diam karena lulus; bedanya hanya terlihat bila keluarannya dibaca, bukan di-`grep`.

⭐ **Dua penjaga TERBALIK — bentuk yang jarang dan memang diminta tiket.**
Yang satu menuntut jalur tulis `OS_AKSEPTASI_KLAIM_LIFE` **tetap ada** (*"Test yang menolak
penulisan … justru gagal"*), sebab hilir masih membaca dari sana dan yang dibuang hanya JSON. Yang
lain menuntut ketujuh kolom tak-terbawa **tetap** tak terbawa. Keduanya ada supaya orang yang kelak
"membersihkan" atau "melengkapi" — niat yang sangat wajar — menabrak kalimatnya lebih dulu, bukan
menabrak Arasapas yang berhenti menerima data.

### ⛔ AC yang belum tertutup — 4

⛔ **STANDAR GANDA SAYA, ditemukan review.** Saya mencabut centang AC *"penomoran tidak melompat"*
dengan alasan *"aturan tanpa pengurai adalah aturan tanpa masukan"* — lalu **mencentang** AC
*"`NO_SEQ` per kombinasi kunci"* yang bersandar pada **fungsi yang sama**, dengan **nol** pemanggil
produksi yang sama. Satu berkas, dua vonis berlawanan. Centang kedua dicabut.

| AC | Sebab |
| --- | --- |
| *"terbawa beserta statusnya dan **posisi tahapnya**"* | **status terbawa; tahap tidak ada di sumber.** Sensus dua arah di atas: nol dari 62 kolom, nol berkas korpus. Menurunkannya dari `STS_REJECT` berarti MERANCANG - `[terbuka — work owner]`: dari mana tahap klaim berjalan diambil saat cutover |
| *"penomoran tidak melompat dan tidak mengulang"* | aturannya ada dan diuji, tetapi **pengurai** nomor lama menjadi `(CLASS, JENIS, TAHUN, urut)` belum ada — bentuk nomornya butir **o1–o3**, masih terbuka |
| *"`NO_SEQ` terjaga per kombinasi kunci"* | ⛔ **dicabut sesudah review.** `TandaAirSequence`, `PeriksaTandaAir`, dan `NomorTerpakai` punya **nol** pemanggil produksi — masukannya belum ada, sama persis dengan AC di atasnya. Aturannya benar dan terbukti; yang belum, ia belum menjaga apa pun |
| *"master diperlakukan sebagai view atas `JSONDATA`; produk dari `product_life` relasional"* | ⛔ **dicabut sesudah review.** Yang ada hanya **larangannya**: nol kode membaca `RATE_LIFE`, `PRODUCTINWARD_LIFE`, `CURRENCY`, maupun `product_life`. Larangan tanpa pembaca menjaga ruang kosong |

### Hasil /code-review — titik tetap `b0afeb4`

**Spec — 3 temuan, semuanya sahih.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ **standar ganda**: AC *"`NO_SEQ` per kombinasi"* dicentang atas fungsi yang **sama** dan pemanggil **sama-sama nol** dengan AC yang saya cabut centangnya sendiri | **dicabut**; alasannya ditulis di atas |
| AC master-sebagai-view: yang ada hanya **larangannya**, nol pembaca | **dicabut** |
| ⛔ `polaJSONProduk` melarang `JSONDATA` **telanjang** se-modul — bertabrakan dengan AC-nya sendiri, yang justru menuntut ketiga master diperlakukan sebagai view **ATAS** `JSONDATA` | **dipersempit** ke pasangan `m_product_life` + `JSONDATA`; dibuktikan dua arah: JSON produk tertangkap, `JSONDATA` sebagai view **tidak** dituduh |
| ketiga sensus (tahap absen · 62/55/7 · tulis datar di dalam `Simpan`) | **diperiksa ulang oleh reviewer dan cocok seluruhnya** |
| keempat penjaga dimutasi reviewer | **keempatnya merah**; caveat dicatat: menghapus gelung tulis datar secara naif gagal KOMPILASI, dan penjaga tidak dapat membedakannya |

⛔ **Backslash termakan untuk kelima kalinya di sesi ini.** Regex `polaJSONProduk` mendarat memuat
**newline sungguhan** di dalam kelas karakternya — dan ia tetap **kompilasi dan hijau**, yang justru
bentuk kegagalan paling berbahaya: bukan merah palsu melainkan hijau palsu. Ditemukan hanya karena
bentuk yang mendarat diperiksa, bukan hasil testnya. Sejak itu backslash dibangun dari
`bytes([92])` dan hasilnya selalu dilihat.

⭐ **Nilai review di tiket ini bukan menemukan kode salah** — kodenya benar dan penjaganya
menggigit. Ia menemukan **penilaian saya yang tidak konsisten terhadap bukti yang sama**, dan satu
penjaga yang melarang persis hal yang AC-nya perintahkan.

**lanjut dari sini:** tiket 13 selesai. Empat AC menunggu keputusan Anda, bukan menunggu kode.

**Standards — 3 pelanggaran keras + 5 elakan TERBUKTI + 6 bau.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ **LABEL DINAIKKAN.** Saya menulis `[terverifikasi]` OQ-002 untuk fakta DDL penomoran; sumbernya `SUMBER-PENOMORAN-DBA.md` melabeli dirinya sendiri `[data DBA — belum dikonfirmasi DBA]`, dan nomor OQ-nya pun salah (DDL milik OQ-001, bukan OQ-002) | **diralat di dua berkas**, dengan path + baris sumbernya |
| prosa berkata penghitung wajib **"DI ATAS"** tanda air; kodenya menerima **kesamaan** | **prosa yang salah**, bukan kodenya: `NO_SEQ` = nomor terakhir yang DITERBITKAN (`SUMBER-PENOMORAN-DBA.md` 72-86). Prosa diperbaiki; prosa yang lebih ketat daripada kodenya mengundang orang "memperbaiki" kode yang benar |
| §4a trap 5: "sensus dua cara" untuk 62 − 55 = 7 | **diralat** — dua wadah yang dikurangkan bukan dua cara |
| ⛔ **B2 paling telak:** pola `INSERT INTO … RATE_LIFE` **tidak pernah mungkin menggigit** — `Qualify` wajib (ADR-U-0033), jadi nama tabel SELALU terpisah dari kata kerjanya. Penjaga itu hanya sanggup menangkap kode yang ADR-U-0033 sudah tolak | **diganti aturan atas NAMA** |
| B1 variabel perantara · B3 baris berbeda · B5 struct sematan | ketiganya **diganti** satu aturan atas nama, lebih sederhana dari pola yang digantikannya |
| B4 penjaga TERBALIK memeriksa **teks**, bukan **penulisan** | kini menuntut `INSERT INTO %s` + `datarT` **di dalam** `Simpan` |
| berkas `.sql` tidak pernah dibaca padahal di-`go:embed` dan dijalankan | **ikut dibaca** |
| bau: `contains` menulis ulang `strings.Contains` · pesan galat tanpa batas · `Urut int` tanpa kutipan rentang | ketiganya **diperbaiki**; batas 5 digit `LPAD` dikutip |

⛔ **Vonis yang paling pahit dan paling berguna:** *"Kedua penjaga terbalik jatuh. Tidak satu pun
menambah deteksi — B4 dan B5 hanya tertangkap penjaga yang SUDAH ADA; B1–B3 tidak tertangkap apa
pun."* Lima elakan, semuanya kompilasi, semuanya hijau. Penjaga yang tidak menambah deteksi lebih
buruk daripada tidak ada penjaga: ia membeli rasa aman dengan harga nol.

⛔ **Backslash termakan untuk keenam kalinya**, dan kali ini **tertangkap sebelum test dijalankan**
karena bentuk yang mendarat diperiksa lebih dulu. Itu satu-satunya cara jebakan ini dapat dilihat.

⭐ **Yang tersisa dinyatakan, bukan ditambal:** `CURRENCY` **tidak dapat dijaga lewat nama** — ia
juga nama kolom dan muncul puluhan kali secara sah. Yang menjaganya hanya kenyataan bahwa nol kode
menyentuh master mana pun hari ini; bila itu berubah, ia harus ditinjau dengan tangan. Dan
`buangKomentarSumber` hanya membuang komentar sebaris penuh — komentar blok tetap tinggal.


## Implementasi — A4, kode saja (28 September 2026)

⛔ **Tidak dijalankan terhadap DEV.** Brief melarangnya, dan `T_MIGRASI` DEV masih di `016`.
Yang dibangun kodenya; uji bertag `db` **melewati** selama `ORACLE_DSN` belum dikonfigurasi.

### Apa yang tabel datar TIDAK punya — dan karena itu tidak dikarang

Daftar 55 kolom `OS_AKSEPTASI_KLAIM_LIFE` dibaca VERBATIM dari
`UpdateOsAkseptasiClaimLife_sql.xml`. Ia **tidak memuat** satu pun kolom tahap, waktu lahir kasus,
maupun status kerja. Maka:

| Kolom baru | Diisi dari | Keadaan |
| --- | --- | --- |
| `TAHAP` *(butir at)* | — | **KOSONG**, dan dilaporkan tiap kasus. Tahap kosong tidak menawarkan tombol apa pun — gagal **tertutup**, yang benar untuk kasus yang tahapnya memang tidak diketahui |
| `TGL_CREATE` *(butir au)* | `CLAIM_RECEIVED_DATE` | **Penggantian**, dilaporkan. Ia **bukan** `pxCreateDateTime`: tanggal klaim *diterima* tidak sama dengan waktu kasusnya lahir di Pega |
| `STATUS_WORK` *(butir bb)* | `COMPLETE_DATE` terisi → `Resolved-Completed` | **KESIMPULAN**, bukan bacaan; ditandai `[terbuka — work owner]` pada setiap baris laporan |

⚠️ Alternatif untuk `STATUS_WORK` — mengosongkan seluruhnya — membuat setiap kasus warisan yang
sudah selesai tampak **masih berjalan** di kotak masuk. Itu selisih yang lebih besar, jadi
kesimpulannya diambil **dan dinyatakan**, bukan didiamkan ke salah satu arah.

### Dokumen warisan — dan kenapa OQ-J membatasi apa yang mungkin

`POOLDATA.DOCUMENT_CLAIM` `[katalog DEV]` punya `IDPEGA`, `NOAKSEP`, `NAMAFILE`, `MIME`,
`KATEGORI_1/2`, `T_STORAGE_ID`. Tautannya:

- `IDPEGA` = `pyWorkPage.pzInsKey` *(`InsertDocument_Act.xml` b185)* → menunjuk **klaim**, bukan peserta;
- `KATEGORI_1` = kunci kelompok `DL-`+angka *(`SaveAttachLife.xml` b595)* → pembukuan unggahan, **bukan** identitas;
- **`NOAKSEP`** ↔ `OS_AKSEPTASI_KLAIM_LIFE.NO_ACCEPTATION` → **satu-satunya** tautan ke peserta yang sumbernya punya.

⛔ Saringan Pega yang sebenarnya `.KATEGORI_1 = .DOCUMENT` pada halaman **peserta**
*(`LoadDocumentLife_ACT.xml` b495)*, dan kolom `DOCUMENT` **tidak ada** di mana pun yang kami
terima — **OQ-J**. Pemetaan ini karena itu **tidak dapat** meniru saringan aslinya.

Yang dilakukan: gantung lewat `NOAKSEP`; setiap yang tidak cocok **dilaporkan dan tidak
dipindahkan**. ⚠️ Arahnya disengaja: dokumen milik orang lain yang menempel pada peserta yang
salah jauh lebih buruk daripada dokumen yang dilaporkan hilang. `IDPEGA` yang bertentangan dengan
`NOAKSEP` dilaporkan pula; `KATEGORI_1` berbentuk asing **dibawa apa adanya** dan dilaporkan —
kunci yang bentuknya aneh tetap menunjuk kelompok unggahan yang nyata.

⛔ `BASE64` **tidak** dibawa ke struktur mana pun: isi berkas tidak pernah masuk artefak. Migrasi
memindahkan **catatannya**; berkasnya sudah ada di penyimpanan dan dirujuk `T_STORAGE_ID`.

### Diagnosa warisan — tidak ada sumbernya, dan itu dinyatakan

Tabel datar punya `DISEASE` dan `ICD_CODE` **tunggal** — satu pasang per baris peserta —
sedangkan `.DiagnoseList` *(butir bd)* adalah **daftar**. Daftar itu hidup di halaman kerja Pega,
tersimpan sebagai **BLOB** milik mesin Pega, dan **tidak diekspor sebagai tabel**.

Akibatnya migrasi membawa paling banyak **satu** diagnosa per peserta; daftar yang lebih panjang
**tidak dapat dipulihkan**. Dilaporkan setiap jalan, `[terbuka — work owner]`, dan **nol** baris
`T_CLAIMLF_DIAGNOSE` dibuat dari tebakan.

### ⚠️ Satu hal yang A4 harus tahu dari paket sebelumnya

`Update_T_Storage_SQL.xml` menulis `EXPDATE` dengan `'DD/MM/YYYY'` dan `TANGGAL_UPLOAD` dengan
`'MM/DD/YYYY'` — **dua bentuk berbeda di satu pernyataan**. Nilai warisan pada kedua kolom itu
tidak dapat dibedakan untuk tanggal 1–12 tiap bulan.

### Penjaga

`TestA4TidakMenyentuhBasisData` membaca berkasnya sendiri dan menolak `db.sql`, `QueryContext`,
`ExecContext`, `tx.tx`, `Qualify(`, dan `time.Now()`. Sebabnya keras: fungsi migrasi yang diam-diam
membuka koneksi akan **menjalankan dirinya sendiri terhadap DEV** saat seseorang menjalankan uji,
dan brief melarang eksekusi A4 sampai work owner menyetujuinya. Dibuktikan merah dengan menanam
`time.Now()`.

**AC:** tidak ada AC tiket ini yang berubah centangnya — **kodenya** yang ada, bukan eksekusinya.
