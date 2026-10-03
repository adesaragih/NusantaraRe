# Spesifikasi — Lima Modul yang Sudah Terverifikasi Penuh

> **Lingkup:** siklus **New Business** (`NB FacIn\`, K-005). Hanya bagian yang **terverifikasi penuh**
> dari korpus dan/atau sudah ditutup keputusan work owner (**K-001…K-021**, **ADR-0001…0006**).
> Bagian yang masih menunggu jawaban pihak lain ada di **Out of Scope**, lengkap dengan apa yang ditunggu.
>
> **Seam pengujian disetujui work owner 16 September 2026** — tiga seam runtime + satu pemeriksaan
> waktu-kompilasi. Tidak ada seam lain yang boleh ditambahkan tanpa keputusan baru.
>
> Kosakata mengikuti `../steering/GLOSARIUM.md`. Label mengikuti `CLAUDE.md` §3.

---

## Problem Statement

Aplikasi Facultative Inward berjalan di atas Pega yang akan dimatikan. Ketika itu terjadi, sebagian
pengetahuan tentang cara sistem ini menghitung dan mengambil keputusan **musnah bersamanya** — bukan
"sulit diambil", melainkan tidak ada sumber penggantinya.

Tetapi pekerjaan tidak dapat menunggu sampai semuanya diketahui. `[terverifikasi]` Sebagian besar
sistem lama **belum dapat dispesifikasikan sekarang**: isi 35 stored procedure yang dilewati seluruh
tulisan produksi tidak ada di korpus, tidak ada satu pun DDL sehingga tipe setiap kolom tidak
diketahui, 604 dropdown tidak punya daftar nilai, dan baris tabel keputusan yang menentukan arah
setiap keputusan underwriting tidak ikut terekspor.

Menunggu semua itu berarti **berhenti berbulan-bulan**. Memulai semuanya sekaligus berarti **menebak** —
dan menebak dilarang, karena setiap tebakan menghasilkan selisih angka yang tidak dapat dijelaskan saat
kedua sistem dijalankan berdampingan, lalu menyembunyikan salah-port yang sesungguhnya.

Yang dibutuhkan work owner adalah **potongan pekerjaan yang dapat dimulai hari ini tanpa satu pun
tebakan** — bagian yang buktinya sudah lengkap di korpus, keputusannya sudah diambil, dan hasilnya
dapat dibuktikan benar terhadap sistem lama.

---

## Solution

Lima modul dibangun lebih dulu, dipilih **bukan** karena mudah, melainkan karena **tidak satu pun
bagiannya bergantung pada jawaban yang belum masuk**:

| Modul | Mengapa sudah aman dibangun |
| --- | --- |
| **`pkg/money`** | Kebutuhan mata uang wajib terbukti berlapis; keadaan `Unknown` sudah diputuskan (K-012) |
| **`pkg/ratio`** | Skala rasio per lini bisnis **terkunci** (K-018); aturan penguraian pembagi komposit terbukti dari struktur ekspresi |
| **`internal/rules`** | Kondisi **nol dari 601** rule `When` tidak terbaca; sikap untuk yang belum terbukti sudah ditetapkan |
| **`services/acceptance`** | Logikanya terbaca penuh; yang hilang hanya **datanya**, dan itu digantikan fixture yang sekaligus jadi kontrak ke DBA (K-009) |
| **`services/premium`** | Rumus terbaca per baris; satuan tiap faktor terkunci; presisi pembulatan terbaca per langkah |

Ketiganya diuji lewat **tiga seam** saja, masing-masing di titik tertinggi modulnya, ditambah satu
pemeriksaan yang dijalankan **kompilator**, bukan test.

Ukuran keberhasilannya tunggal dan tidak dinegosiasikan: **rekonsiliasi eksak — nol selisih sampai
digit terakhir** (ADR-0001). Bukan toleransi. Bila `decimal` dipakai dan urutan operasi serta presisi
per-langkah direproduksi apa adanya, hasilnya **harus** identik; selisih sekecil apa pun berarti ada
salah-port, dan toleransi hanya menyembunyikannya.

---

## User Stories

### `pkg/money` — uang tidak pernah `float`

1. Sebagai pengembang, saya ingin nilai uang selalu berupa `decimal`, agar nilai pertanggungan puluhan
   miliar tidak kehilangan presisi diam-diam seperti yang terjadi pada `float64`.
2. Sebagai pengembang, saya ingin setiap nilai uang **membawa mata uangnya**, agar dua nilai dari mata
   uang berbeda tidak pernah dijumlahkan tanpa disadari.
3. Sebagai penguji rekonsiliasi, saya ingin pembulatan uang bersifat deterministik, agar selisih yang
   muncul pasti berarti salah-port dan bukan derau aritmetika.
4. Sebagai pengembang, saya ingin mata uang boleh berkeadaan **`Unknown` yang eksplisit**, agar 112
   layar yang di sistem lama berjalan mulus tanpa field mata uang tidak berhenti di sistem baru.
5. Sebagai work owner, saya ingin `Unknown` **tidak** diam-diam menjadi IDR, agar tebakan tidak masuk
   ke tempat yang paling mahal — penempatan valas.
6. Sebagai pengembang, saya ingin aritmetika **lintas mata uang** dengan `Unknown` gagal keras, agar
   celah pengetahuan muncul di titik yang benar-benar memerlukan jawabannya.
7. Sebagai tim migrasi, saya ingin implementasi ini **menghasilkan hitungan** berapa banyak nilai
   produksi yang tiba tanpa mata uang, agar pertanyaan ke DBA berubah dari dugaan menjadi ukuran.
8. Sebagai pengembang, saya ingin parser masukan menerima **string berkoma desimal**, agar 1.057 titik
   konversi di sistem lama tidak menjadi 1.057 cacat di sistem baru.

### `pkg/ratio` — skala melekat pada nilai

9. Sebagai pengembang, saya ingin rate dan persen memakai tipe **yang berbeda dari uang**, agar
   `Uang + Rasio` ditolak kompilator alih-alih menghasilkan angka yang salah diam-diam.
10. Sebagai pengembang, saya ingin satu-satunya jembatan antara keduanya adalah operasi eksplisit
    `Uang × Rasio → Uang`, agar setiap perkalian terlihat di kode dan dapat ditelusuri.
11. Sebagai aktuaris, saya ingin setiap rasio **membawa skalanya sendiri** (‰ atau %), agar properti
    rate yang sama tidak tertukar satuannya antar lini bisnis.
12. Sebagai work owner, saya ingin skala diisi **satu resolver terpusat**, agar tidak ada tempat kedua
    yang bisa menyimpang tanpa ketahuan.
13. Sebagai work owner, saya ingin resolver itu mengikuti **rumus**, bukan label layar, karena
    `[terverifikasi]` tiga label layar bertentangan dengan rumusnya.
14. Sebagai pengembang, saya ingin lini bisnis yang **belum ada di peta skala** memicu `panic`, bukan
    memakai skala default, karena menebak skala adalah persis cacat 10× yang tipe ini dirancang untuk
    mencegah.
15. Sebagai aktuaris, saya ingin **aturan penguraian pembagi komposit** tertulis sebagai aturan, agar
    pembagi `100000` dibaca sebagai `1.000 × 100` dan **menegaskan** rate ber-‰, bukan membantahnya.
16. Sebagai penguji rekonsiliasi, saya ingin peta skala terkunci dapat dibaca sebagai satu tabel, agar
    ketika angka tidak cocok saya tahu persis apa yang diasumsikan.
17. Sebagai work owner, saya ingin skala **tidak** disimpan di tabel konfigurasi yang dapat diubah
    tanpa deployment, agar perilaku sistem tidak berubah tanpa jejak version control.

### `internal/rules` — registry predikat `When`

18. Sebagai pengembang, saya ingin seluruh predikat `When` hidup di **satu registry**, agar tidak ada
    salinan logika gerbang yang tersebar dan menyimpang.
19. Sebagai pengembang, saya ingin setiap predikat menyebut **rule Pega asalnya** dalam komentar, agar
    rekonsiliasi paralel run dapat menunjuk sumber setiap perbedaan.
20. Sebagai auditor, saya ingin kondisi dibaca dari **kedua tag** kondisi, agar rule tidak tampak
    "tanpa kondisi" padahal kondisinya tersimpan di tag yang satunya.
21. Sebagai work owner, saya ingin predikat yang kondisinya **terbaca tetapi belum terbukti dieksekusi**
    memicu `panic`, agar tidak ada gerbang yang diam-diam selalu benar atau selalu salah.
22. Sebagai pengembang, saya ingin nama predikat yang **tidak dikenal** memicu `panic`, agar salah ketik
    tidak berubah menjadi gerbang yang selalu tertutup.
23. Sebagai work owner, saya ingin predikat yang hilang dari satu folder tetapi terbaca di folder lain
    **tetap diimplementasikan** dengan mencatat asal salinannya, karena "hilang dari ekspor" bukan
    berarti "usang".
24. Sebagai pengembang, saya ingin gerbang masuk siklus New Business memakai **ekspresi tersimpan**,
    bukan teks tampilannya, sesuai keputusan yang sudah diambil.
25. Sebagai penguji rekonsiliasi, saya ingin predikat yang perilakunya menyimpang dari namanya tetap
    diport **apa adanya** dengan namanya dipertahankan, agar ketertelusuran ke rule Pega asal tidak putus.

### `services/acceptance` — tangga akseptasi

26. Sebagai underwriter, saya ingin satu keputusan saya menghasilkan **satu transisi** pada tangga
    akseptasi, agar perilakunya sama persis dengan yang saya kenal hari ini.
27. Sebagai pengembang, saya ingin tangga akseptasi **bukan** loop yang menghitung seluruh rantai
    approver sekaligus, karena bentuk itu tidak ada di sistem lama dan berperilaku berbeda ketika rantai
    terputus di tengah.
28. Sebagai pengembang, saya ingin **antrean**, **kode jabatan tujuan**, dan **hasil keputusan** menjadi
    tiga hal terpisah, agar tidak disatukan menjadi satu "status" yang menghapus perbedaannya.
29. Sebagai pengembang, saya ingin kode jabatan tujuan bertipe berbeda dari token antrean, agar
    kompilator menangkap pertukaran keduanya.
30. Sebagai underwriter, saya ingin keadaan "tidak ada jabatan tujuan yang cocok" diperlakukan sebagai
    **tangga selesai** — penyelesaian normal — bukan sebagai galat.
31. Sebagai underwriter, saya ingin **Reject** dan **Decline** tetap dibedakan, karena yang satu masih
    dapat dibanding dan yang lain penolakan final.
32. Sebagai underwriter, saya ingin efek samping `Reject` — mematikan konfirmasi binding dan penerimaan
    R/I slip — tetap terjadi, karena itu perilaku yang dikehendaki dan menjadi syarat munculnya jalur
    banding.
33. Sebagai work owner, saya ingin hasil keputusan dibatasi pada **enam nilai**, agar nilai asing
    tertangkap alih-alih mengalir ke kolom yang menggerakkan alur.
34. Sebagai penguji rekonsiliasi, saya ingin nilai hasil keputusan di luar keenamnya memicu **`panic`
    disertai catatan** selama paralel run, bukan ditolak diam-diam, agar premis kita yang salah
    ketahuan di fase yang memang dirancang untuk menemukannya.
35. Sebagai work owner, saya ingin perbedaan perilaku antara paralel run dan produksi dikendalikan
    **satu flag fase eksplisit**, bukan dua basis kode yang bisa menyimpang satu sama lain.
36. Sebagai DBA, saya ingin bentuk tabel limit yang diharapkan sistem baru **dikirim lebih dulu sebagai
    fixture**, agar ketidakcocokan bentuk ketahuan saat permintaan dikirim dan bukan di akhir proyek.
37. Sebagai work owner, saya ingin ejaan nilai kolom jabatan ditandai sebagai **asumsi yang harus
    dikonfirmasi**, bukan sebagai fakta, karena bila ejaannya tidak cocok maka tidak ada approver yang
    pernah ditemukan dan seluruh tangga macet.

### `services/premium` — rumus premi

38. Sebagai aktuaris, saya ingin premi dihitung dengan **urutan operasi yang persis sama** seperti
    sistem lama, agar angkanya cocok sampai digit terakhir.
39. Sebagai pengembang, saya ingin setiap pembulatan menuliskan presisinya **di tempatnya**, karena satu
    rule tidak punya satu presisi — satu pembagi saja muncul dengan sembilan presisi berbeda.
40. Sebagai pengembang, saya ingin setiap pembulatan menyebut **rule dan langkah asalnya** dalam
    komentar, agar presisi yang tampak aneh dapat diperiksa ke sumbernya alih-alih "dirapikan".
41. Sebagai penguji rekonsiliasi, saya ingin uji rekonsiliasi memuat kasus pembulatan **di dalam loop
    akumulasi**, karena galat per-iterasi menumpuk dan port yang benar untuk satu nilai masih bisa
    meleset untuk daftar panjang.
42. Sebagai aktuaris, saya ingin satuan rate ditentukan **lini bisnis**, bukan satu satuan seragam,
    karena menyeragamkannya menggeser premi satu ordo besaran pada separuh portofolio.
43. Sebagai work owner, saya ingin percabangan lini bisnis terjadi **di dalam satu pintu masuk**, agar
    ada satu tempat yang dapat diaudit dan satu tempat yang dapat direkonsiliasi.
44. Sebagai penguji rekonsiliasi, saya ingin fixture per lini bisnis **sekaligus menjadi berkas
    rekonsiliasi**, agar tidak ada dua sumber kebenaran yang bisa menyimpang.
45. Sebagai aktuaris, saya ingin perhitungan menolak lini bisnis yang belum ada di peta skala, agar
    lini baru tidak ikut dihitung dengan asumsi diam-diam.
46. Sebagai work owner, saya ingin hal yang tampak keliru di sistem lama — perbandingan angka sebagai
    string, tautologi pembanding limit, gerbang yang menihilkan premi — **direproduksi apa adanya** dan
    ditandai sebagai kandidat perbaikan.
47. Sebagai work owner, saya ingin perbaikan label satuan pada layar dilakukan, karena `[terverifikasi]`
    **nol angka berubah** dan karena itu rekonsiliasi tidak terganggu.
48. Sebagai pengembang, saya ingin nilai uang masuk dan keluar dari perhitungan sebagai `Uang`, bukan
    angka telanjang, agar mata uang tidak hilang di tengah rantai perhitungan.
49. Sebagai aktuaris, saya ingin rate dan pro-rata diperlakukan sebagai **rasio bersatuan**, agar
    pembagi gabungan dalam satu operasi tidak salah dibaca sebagai satu satuan tunggal.
50. Sebagai penguji rekonsiliasi, saya ingin setiap baris rumus yang diport dapat ditunjuk ke **baris
    sumbernya di korpus**, agar perbedaan dapat didiagnosis tanpa membuka Pega.

### Lintas modul

51. Sebagai work owner, saya ingin arah dependensi **satu arah** dan tidak memotong lapisan, agar
    perubahan di lapisan bawah tidak merembet ke atas tanpa terlihat.
52. Sebagai work owner, saya ingin modul perhitungan dan modul tangga akseptasi **tidak saling
    memanggil**, karena di sistem lama keduanya memang tidak saling memanggil.
53. Sebagai pengembang baru di tim, saya ingin kosakata kode mengikuti glosarium proyek, agar nama
    warisan yang menyesatkan tidak tertanam ulang di sistem baru.
54. Sebagai auditor, saya ingin setiap keputusan rancangan dapat ditelusuri ke ADR atau keputusan work
    owner bernomor, agar alasannya tetap terbaca setelah orangnya berganti.
55. Sebagai work owner, saya ingin bagian yang **belum dapat dispesifikasikan** tertulis eksplisit
    beserta apa yang ditunggu, agar kekurangan itu terlihat alih-alih tersembunyi.

---

## Implementation Decisions

### Arsitektur umum

**Arah dependensi searah: `handlers → services → repository`** (`CLAUDE.md` §4.2). Tidak boleh terbalik,
tidak boleh memotong lapisan. `pkg/money` dan `pkg/ratio` adalah pustaka murni tanpa dependensi ke
lapisan mana pun.

**`services/premium` dan `services/acceptance` tidak saling memanggil.** `[terverifikasi]` Di sistem
lama, spreading tidak pernah mengubah antrean maupun kode jabatan tujuan (nol `Property-Set` di 74
activity terkait), dan tangga akseptasi membandingkan **nilai dasar akseptasi**, bukan premi. Menaruh
keduanya di balik satu façade berarti mengarang dependensi yang tidak ada — dilarang `CLAUDE.md` §1.

**Ketertelusuran (`CLAUDE.md` §4.6).** Setiap predikat, transisi, rumus, dan pembulatan menyebut rule
Pega asalnya dalam komentar, dengan nomor baris bila baris itu yang membuktikannya. Ini yang membuat
rekonsiliasi paralel run mungkin.

---

### Modul 1 — `pkg/money`

**Seam: tidak ada seam sendiri.** Teruji lewat Seam 1 (`services/premium`). Memberinya seam terpisah
berarti menguji pustaka `decimal` milik orang lain, bukan keputusan kita.

**Keputusan — bentuk tipe** (ADR-0006, K-012; `CLAUDE.md` §4.1):

```go
type Money struct {
    Amount   decimal.Decimal
    Currency Currency   // WAJIB menyertai — boleh berkeadaan Unknown yang eksplisit
}
```

**Keadaan `Unknown` adalah keadaan sah, bukan kegagalan.** Dasarnya `[terverifikasi]`: 112 Section
menampilkan nilai uang tanpa field mata uang mana pun, sementara hanya ada 93 pengikatan properti mata
uang di seluruh NB. Sistem lama juga **tidak pernah** menetapkan mata uang default di kode — ia selalu
datang dari data.

| Operasi | `Unknown` |
| --- | --- |
| Baca, tampilkan, simpan di memori | diizinkan |
| Aritmetika dalam mata uang yang sama-sama `Unknown` | diizinkan |
| Aritmetika **lintas** mata uang | **`panic`** |
| Tulis ke Oracle | **`panic`** — lihat Out of Scope butir 5 |

Default ke IDR **ditolak**: itu menebak, dan tebakannya akan salah persis pada kasus paling mahal.
`panic` tanpa kecuali juga **ditolak**: akan menghentikan 112 layar yang hari ini berjalan mulus.

**Parser masukan menerima string berkoma desimal.** `[terverifikasi]` 1.057 titik konversi di korpus.
Konversi terjadi **di batas input**, bukan tersebar di dalam perhitungan.

**Efek samping yang disengaja:** implementasi ini menghasilkan hitungan berapa banyak nilai produksi
tiba tanpa mata uang — angka itu menjawab Pengukuran D1 sebagai **ukuran**, bukan dugaan.

---

### Modul 2 — `pkg/ratio`

> ⚠️ **Disinkronkan 1 Oktober 2026:** paragraf di bawah sudah DIGANTI keputusan work owner — resolver
> lini bisnis → satuan kini **Seam 2**, `premium.SatuanRate(lini)`, beserta `SatuanProRata()`
> (`KEPUTUSAN-30-09-2026.md` butir 26–27). Teks asli dibiarkan untuk jejak.

**Seam: tidak ada seam sendiri.** Teruji lewat Seam 1. Nilai modul ini justru pada **skala yang
menempel**, dan itu hanya terbukti saat dipakai menghitung premi.

**Keputusan — bentuk tipe** (ADR-0004, K-010):

```go
type Scale uint8          // PerMille | Percent — tidak ada nilai default
type Ratio struct {
    Value decimal.Decimal
    Scale Scale           // skala MELEKAT pada nilai, bukan pada pemanggil
}

// Satu-satunya jembatan ke uang. Money + Ratio tidak dapat dikompilasi.
func (m Money) Times(r Ratio) Money
```

**Peta skala per lini bisnis — terkunci K-018:**

| Lini bisnis (COB) | Skala rasio | Pembagi yang menempel pada rate |
| --- | :-: | ---: |
| PA · Layering · FIRE | **‰** | 1.000 |
| MBU · ANEKA · BONDING · GOLF · MARINE CARGO | **%** | 100 |

**Aturan penguraian pembagi komposit — mengikat.** Pembagi gabungan pada satu operasi pembagian adalah
**hasil kali** sumbangan tiap faktor bersatuan, bukan satu satuan tunggal. Satuan rate dibaca dari
**faktor yang menempel padanya saja**, tidak pernah dari pembagi total:

```
@Math.divide((.TSI * .Rate * ProRatePercent), 100000, 4)
        100000  =  1000 (rate ber-‰)  ×  100 (ProRatePercent ber-%)
```

`[terverifikasi]` Asal: `GenerateLayerList_ACT.xml` L909 — dan ini **satu-satunya** ekspresi ber-rate di
seluruh 15 berkas bernuansa Layer di korpus NB, tanpa varian yang bertentangan.

⚠️ **Aturan ini bukan tafsir, dan buktinya struktural.** `FillPremiMBU_FacIn.xml` L1144 memakai bentuk
**bersarang**: pembagi **dalam** menempel pada rate, pembagi **luar** menempel pada pro-rata. Keterikatan
faktor→pembagi karena itu tertulis eksplisit dalam struktur ekspresinya.

**Resolver lini bisnis → skala mengikuti RUMUS, bukan label layar.** `[terverifikasi]` Tiga label layar
bertentangan dengan rumusnya; K-018 menetapkan rumus yang menang di ketiganya. **Lini bisnis yang belum
ada di peta → `panic`**, bukan skala default.

**Skala tidak disimpan di tabel konfigurasi.** Ditolak secara eksplisit di ADR-0004: skala adalah fakta
struktural terverifikasi dari korpus, bukan parameter bisnis. Menaruhnya di tabel yang dapat diubah tanpa
deployment mengulang kerentanan yang sudah ditemukan pada tabel prompt model AI — perilaku sistem berubah
tanpa jejak version control.

---

### Modul 3 — `internal/rules`

**Seam 1 dari 3 — registry predikat:**

```go
// Satu seam untuk 226 predikat. Predikat individual BUKAN seam.
func (r *Registry) Eval(name string, c *Case) bool
```

**Kondisi wajib dibaca dari KEDUA tag** (`CLAUDE.md` §4.5). `[terverifikasi]` **170 berkas (28,3 %)**
menyembunyikan kondisinya dari tag tampilan — 13 tagnya kosong, 157 berisi teks placeholder — sementara
seluruh 601 berkas punya tag ekspresi tersimpan terisi. Membaca hanya tag pertama membuat rule tampak
"tanpa kondisi" padahal kondisinya ada.

`[terverifikasi]` **Tidak ada satu pun rule `When` yang kondisinya tidak terbaca — 0 dari 601.**

**Sikap per rule yang sudah ditetapkan:**

| Rule | Sikap | Dasar |
| --- | --- | --- |
| `ToUW` · `ToJUW_A` | Implementasikan sesuai kondisi terbaca | `[terverifikasi]` |
| `LetterNoNull` · `IsEdmInternalRetro` | Implementasikan sesuai kondisi terbaca | `[dugaan]`, label ter-resolve |
| `IsPKSASM` | **`panic`** | Kondisi terbaca tetapi label belum ter-resolve dan bertanda sementara — belum terbukti sebagai yang dieksekusi |
| nama tidak dikenal | **`panic`** | Salah ketik tidak boleh berubah jadi gerbang tertutup |
| `IsOfferFacIn` | Ekspresi tersimpan yang berlaku | K-002 — teks tampilan menyebut kondisi lain, dicatat sebagai kandidat perbaikan |
| `IsSpreadingDepan` | Implementasikan, catat asal salinannya | K-003 — hilang dari folder NB, terbaca penuh di folder lain |
| `IsFacout` | Diport apa adanya, **nama dipertahankan** | K-019 — fitur usang secara bisnis, kode masih aktif |

⚠️ **`IsFacout` adalah contoh "nama bukan bukti" yang paling tajam.** Namanya menyiratkan fac out,
`[terverifikasi]` isinya menguji **hasil keputusan = Banding**. Namanya dipertahankan demi ketertelusuran
ke rule Pega asal (§4.6), **tetapi namanya bukan dokumentasi artinya** — komentar wajib menyebutkannya.

⚠️ **Caveat `IsSpreadingDepan`:** salinan yang terbaca berasal dari folder siklus lain pada versi ruleset
yang lebih lama. `[terverifikasi]` 95 rule di korpus ini berbeda versi antar folder, jadi tidak dapat
dipastikan salinan NB identik. Implementasikan kondisinya, catat asalnya.

> ### ⛔ Daftar nilai `IsB2B` yang baru masuk **TIDAK** mencabut `panic` pada `IsPKSASM`
>
> `[terverifikasi]` 17 September 2026, `D:\migrasi\RNM\DDL\B2B.xml` memuat empat nilai sah:
> `ASM` · `KBRU` · `BDX` · `SRB`. Godaannya adalah menyimpulkan bahwa `IsPKSASM` kini aman
> diimplementasikan. **Itu keliru.**
>
> **Penyebab `panic`-nya bukan nilai yang tidak dikenal.** `CLAUDE.md` §4.5 menyebut tiga hal, dan
> tidak satu pun disentuh oleh daftar nilai itu: `<pyConditionString>` masih berisi placeholder
> `[Double click to add condition]`, **labelnya belum ter-resolve**, dan `<pyTempText>` bernilai
> `true` — di ketiga folder korpus.
>
> Mengetahui **nilai apa saja yang sah** tidak membuktikan **kondisi itu benar-benar dieksekusi**.
> Itu dua pertanyaan berbeda, dan hanya yang kedua yang menyebabkan `panic`.
>
> **`IsPKSASM` tetap `panic`, tidak diubah.** Ia hanya dicabut bila kondisinya terbukti aktif —
> misalnya dari ekspor produksi tunggal yang labelnya sudah ter-resolve.

---

### Modul 4 — `services/acceptance`

**Seam 2 dari 3 — satu transisi tangga:**

```go
func Next(state LadderState, d Decision, limits LimitTable) (LadderState, error)
```

**Keputusan — bentuk keadaan** (ADR-0003, K-009):

```go
// Satu keputusan manusia = satu transisi. BUKAN loop yang menghitung
// seluruh rantai approver sekaligus — bentuk itu tidak ada di sistem lama
// dan berperilaku berbeda ketika rantai terputus di tengah.
type LadderState struct {
    Antrean                Antrean         // PositionNote — ruang nama ANTREAN
    NextApproverPosition   KodeJabatan     // LetterNo     — ruang nama KODE JABATAN
    HasilKeputusan         Hasil           // ProposalAcceptStatus
    TanggaSelesai          bool            // penyelesaian NORMAL, bukan galat
}
```

**Tiga field state terpisah, dan dua di antaranya tipe berbeda.** `[terverifikasi]` Token antrean dan
kode jabatan hidup di **ruang nama yang berbeda**. Tipe Go yang berbeda mencegah tertukar saat kompilasi —
salah satu tempat di mana sistem tipe menangkap cacat warisan.

⚠️ **`next_approver_position` disimpan di properti warisan bernama `LetterNo`, yang tidak berisi nomor
surat.** Kosakata ini mengikat (`GLOSARIUM.md`); nama warisan tidak boleh ditanam ulang di sistem baru
kecuali sebagai komentar asal.

**Tangga selesai adalah penyelesaian normal.** `[terverifikasi]` Ketika tidak ada jabatan tujuan yang
cocok, sistem lama mengambil cabang `Else` — wewenang sudah cukup. Merancangnya sebagai galat mengubah
perilaku.

**Domain hasil keputusan dikunci ke enam nilai** (K-021): Accept · Reject · Ask · **Banding** · Decline ·
Revise. `[terverifikasi]` Dekodernya literal dan hadir identik di ketiga folder korpus. Dua nilai lain
milik ranah ceding dan hidup di properti terpisah; satu nilai lagi nol jejak di seluruh 6.071 berkas.

⚠️ **Validasi domain ini adalah perilaku BARU.** `[terverifikasi]` Sistem lama tidak punya validasi
domain pada kolom ini — tidak punya validasi apa pun selain wajib-isi. Karena itu:

| Fase | Nilai di luar domain | Alasan |
| --- | --- | --- |
| Paralel run | **`panic` + catatan** | Nilai tak terduga berarti premis kita salah, dan paralel run memang ada untuk menemukannya |
| Produksi | **gagal keras di jalur tulis + catatan** | Menulis nilai tak dikenal ke kolom yang menggerakkan alur lebih berbahaya daripada berhenti |

**Reject dan Decline tidak disatukan** (K-014). Decline = penolakan final tanpa hak banding; Reject =
masih dapat dibanding. Efek samping Reject yang terverifikasi — mematikan konfirmasi binding dan
penerimaan R/I slip, lalu menjadi syarat munculnya jalur banding lewat hitungan riwayat berstatus reject —
adalah **perilaku yang dikehendaki**, direproduksi apa adanya.

**Nilai dasar akseptasi = nilai pertanggungan penuh** untuk New Business dan Renewal (K-015, disengaja).
Siklus endorsement memakai **selisih** terhadap before-image — di luar lingkup spec ini.

**Flag fase tunggal, bukan dua basis kode** (ADR-0002, K-008):

| Situasi | Paralel run | Produksi |
| --- | --- | --- |
| (a) hasil keputusan **dikenali**, baris keputusannya belum terverifikasi | `panic` | `decline` + catatan |
| (b) kondisi **tidak dikenali sama sekali** | `panic` | `panic` |

`[terverifikasi]` `decline` dipilih di produksi karena ia **hasil default yang terekam** pada rule
keputusan underwriting — memilihnya adalah reproduksi, bukan tebakan. Transisi jalur (a) dari `panic` ke
`decline` hanya boleh setelah jalur itu diverifikasi terhadap ekspor produksi.

**Tabel limit disuntikkan sebagai fixture, dan fixture itu adalah kontrak ke DBA** (ADR-0003, K-009).
Yang hilang dari tangga akseptasi adalah **datanya**, bukan logikanya: query-nya terbaca penuh, urutan
eskalasinya terverifikasi, cara berhentinya terverifikasi, dan ketiga folder berperilaku identik. Fixture
memuat kolom yang dibaca query untuk ketujuh tabel limit per lini bisnis.

⚠️ **Ejaan nilai kolom jabatan ditandai asumsi, bukan fakta.** Token routing di dalam rule ditulis
**tanpa spasi**; ejaan di Oracle belum diketahui. Bila keduanya tidak cocok, tidak ada approver yang
pernah ditemukan dan seluruh tangga macet total.

⚠️ **15 tautologi pembanding limit diport apa adanya.** `[terverifikasi]` Rule alur tangga akseptasi
eksklusif NB memuat 15 perbandingan yang selalu benar, plus empat nomor polis literal sebagai gerbang
alur. Sebelum bisnis memutuskan, port **mereproduksi** keduanya dan menandainya sebagai kandidat
perbaikan (`CLAUDE.md` §1).

---

### Modul 5 — `services/premium`

**Seam 3 dari 3 — satu pintu masuk perhitungan:**

```go
// Percabangan lini bisnis terjadi DI DALAM. Rumus per-COB adalah fungsi internal,
// bukan seam. Ini titik rekonsiliasi tahap 1 (ADR-0001).
func Calculate(in Input) (money.Money, error)
```

**Urutan operasi dipertahankan persis.** `round(a,4) * b`, bukan `round(a*b,4)`. Ini bukan preferensi
gaya — ia menentukan digit terakhir, dan digit terakhir adalah ukuran keberhasilan.

**Presisi pembulatan literal di tiap langkah** (ADR-0005, K-011), disertai komentar yang menyebut rule,
langkah, dan presisinya. Registry terpusat **ditolak** karena mengandaikan satu rule punya satu presisi,
dan `[terverifikasi]` itu tidak benar: aritmetika yang sama dibulatkan 4 desimal di satu rule dan 20
desimal di rule lain, dan satu pembagi muncul dengan **sembilan** presisi berbeda.

**Rumus yang diport dan baris sumbernya** `[terverifikasi]`:

| Lini | Rule Pega asal | Baris |
| --- | --- | --- |
| PA | `CalculatePremiPA_FacIn.xml` | L713 · L858 · L1003 · L1146 |
| MBU | `FillPremiMBU_FacIn.xml` | L1144 · L1626 |
| Layering | `GenerateLayerList_ACT.xml` | L909 |

Keempat baris PA sepakat pada ‰: dua memberi pembagi 1.000 langsung, dua memberi 100.000 yang terurai
menjadi 1.000 × 100. Penguatnya di lapisan layar: layar input PA mencetak tanda per mille **secara
harfiah**.

**Perbaikan label satuan dilakukan; nol angka berubah.** Label layar MBU diperbaiki ke `(%)` (K-016);
label tampilan PA dan daftar Layer diperbaiki ke `(‰)` (K-018). `[terverifikasi]` Ketiganya perbaikan
**label saja**, sehingga rekonsiliasi tidak terganggu sama sekali.

📌 **Arah perbaikannya tidak seragam — ini yang paling mudah salah baca.** "Rumus yang menang" **tidak**
berarti "semuanya persen": MBU bergerak ke %, sementara PA dan Layer justru bergerak ke ‰. Salah membaca
ini persis yang dulu melahirkan konflik yang ditutup K-018.

**Lini bisnis yang belum ada di peta skala → `panic`.** Bukan skala default.

---

## Testing Decisions

### Apa yang membuat sebuah test baik di proyek ini

**Hanya perilaku eksternal lewat seam yang disetujui.** Test tidak boleh menyentuh fungsi internal, dan
tidak boleh mengetahui bagaimana percabangan lini bisnis diimplementasikan. Bila sebuah test rusak karena
refactor yang tidak mengubah angka, test itu menguji hal yang salah.

**Rekonsiliasi eksak — nol selisih sampai digit terakhir** (ADR-0001, K-007). Toleransi **ditolak**: bila
urutan operasi dan presisi per-langkah direproduksi apa adanya, `decimal` bersifat deterministik dan
hasilnya harus identik. Selisih sekecil apa pun berarti salah-port, dan toleransi menyembunyikannya.

**Fixture adalah kontrak, bukan sekadar data uji.** Fixture tabel limit dikirim ke DBA sebagai pernyataan
bentuk data yang diharapkan; fixture per lini bisnis di perhitungan premi **adalah** berkas rekonsiliasi
tahap 1. Tidak ada dua sumber kebenaran yang bisa menyimpang.

### Yang diuji di tiap seam

**Seam `services/premium.Calculate`** — table-driven per lini bisnis. Mencakup transitif: aritmetika uang,
skala rasio yang melekat, resolver lini bisnis → skala, penguraian pembagi komposit, dan presisi pembulatan
per-langkah. Kasus wajib:

1. Setiap lini bisnis pada peta skala terkunci, dengan angka acuan dari korpus.
2. ⚠️ **Pembulatan DI DALAM loop akumulasi** (ADR-0005) — bukan hanya nilai tunggal. Galat per-iterasi
   menumpuk, sehingga port yang benar untuk satu nilai masih bisa meleset untuk daftar panjang. Ini
   bentuk kesalahan yang paling sulit terlihat, dan karena itu **wajib punya kasus sendiri**.
3. Bentuk pembagi komposit — memastikan pembagi 100.000 dibaca sebagai 1.000 × 100 dan menghasilkan ‰.
4. Bentuk pembagi **bersarang** (MBU) — memastikan pembagi dalam terikat ke rate dan pembagi luar ke
   pro-rata.
5. Lini bisnis tak dikenal → `panic`.
6. Nilai uang `Unknown` dalam aritmetika satu mata uang → diizinkan; lintas mata uang → `panic`.

**Seam `services/acceptance.Next`** — satu transisi per kasus uji, dengan tabel limit fixture. Kasus wajib:

1. Satu keputusan menghasilkan **tepat satu** transisi; tidak ada kasus uji yang mengharapkan rantai.
2. Tangga selesai — tidak ada jabatan tujuan yang cocok → **penyelesaian normal**, bukan galat.
3. Ketiga field state bergerak independen; antrean dan kode jabatan tidak pernah tertukar.
4. Reject dan Decline menghasilkan efek samping yang berbeda.
5. Hasil keputusan di luar domain enam nilai → `panic` + catatan pada fase paralel run.
6. Flag fase: jalur (a) `panic` saat paralel run dan `decline` saat produksi; jalur (b) `panic` di kedua
   fase.

**Seam `internal/rules.Eval`** — table-driven atas seluruh registry. Kasus wajib:

1. Rule yang kondisinya hanya ada di tag ekspresi tersimpan tetap terbaca benar.
2. Rule bertanda belum terbukti → `panic`.
3. Nama tidak dikenal → `panic`.
4. Gerbang masuk siklus memakai ekspresi tersimpan, bukan teks tampilan.

**Pemeriksaan waktu-kompilasi** — berkas uji yang **wajib gagal dikompilasi** ketika uang dijumlahkan
dengan rasio (ADR-0004). Ini jaminan tipe, bukan perilaku runtime, sehingga **bukan** unit test.

### Prior art

**Tidak ada.** Belum ada satu baris kode pun di repo ini — seluruh struktur folder target masih berupa
rencana. Karena itu ketiga seam di atas **menjadi** prior art bagi modul berikutnya, dan bentuknya harus
ditetapkan dengan sadar sekarang: table-driven, fixture sebagai kontrak, dan tidak ada test yang menyentuh
fungsi internal.

---

## Out of Scope

Setiap butir menyebut **apa yang ditunggu** — tidak ada yang dikeluarkan tanpa alasan dan tanpa jalan
keluar.

| # | Di luar lingkup | Menunggu |
| ---: | --- | --- |
| 1 | **Baris keputusan pemetaan hasil keputusan → konektor alur.** Kosakatanya pasti, logika pemilihnya tidak — baris tabel keputusannya tidak ikut terekspor | **Ekspor produksi tunggal** (permintaan butir 2 ke IT). Sementara itu ditangani flag fase K-008 |
| 2 | **Isi tabel limit yang sebenarnya.** Bentuknya diasumsikan lewat fixture, termasuk ejaan nilai kolom jabatan | **DBA** — isi + DDL tujuh tabel limit (permintaan butir 4) |
| 3 | **Seluruh jalur tulis produksi**: 35 stored procedure, pasangan nilai-sesudah/selisih, penomoran versi polis | **`ALL_SOURCE`** dari DBA (permintaan butir 3). Yang dapat dispesifikasikan sekarang hanya pemanggilnya |
| 4 | **Struktur tabel penyimpanan (flat vs JSON)** | **DITUNDA atas keputusan work owner** — bukan blocker teknis |
| 5 | **Kontrak `panic` saat menulis mata uang `Unknown` ke Oracle** (ADR-0006). Separuh kontrak ini tidak dapat diuji karena batasnya belum ada | **Seam repository**, yang sendirinya menunggu butir 3. Dicatat sebagai **kewajiban yang belum punya tempat**, bukan sebagai kontrak yang dilupakan. Repository tiruan **sengaja tidak ditarik sekarang** — itu mengarang bentuk yang belum diketahui |
| 6 | **Arti enumerasi yang belum dijawab Product** — jenis penawaran, jenis endorsement, jenis deductible, kelompok tim, jenis pro-rata, penanda kerja sama | **Product** — lihat daftar bertanda ☐ di glosarium |
| 7 | **604 dropdown tanpa daftar nilai.** Tidak ada layar isian yang dapat dibangun lengkap | **Product** |
| 8 | **Siklus endorsement dan hal-hal khas renewal.** Fokus spec ini New Business | Fase endorsement. Tiga pertanyaan matang sudah tersimpan utuh dan **tidak boleh dihapus** |
| 9 | **Predikat fac out beserta sepuluh perujuknya** — diport apa adanya, tidak dihapus | Tidak menunggu apa pun; **penghapusannya** adalah perubahan terpisah yang butuh keputusan tersendiri (K-019) |
| 10 | **Pemindahan alamat email ke konfigurasi** | Tidak menunggu apa pun — **pengecualian eksplisit tercatat** terhadap aturan endpoint (K-020). Dapat ditinjau ulang bila daftar penerima berubah |
| 11 | **Enam cabang alur yang memanggil rule yang hilang dari ekspor** — termasuk jalur Special Acceptance dan tangga akseptasi putaran kedua | **Ekspor produksi tunggal** (K-006). ⚠️ **Ditangguhkan, bukan dibuang** — "hilang dari ekspor" bukan berarti "usang" |
| 12 | **Mesin spreading terbesar** — gerbang parameternya belum dijelaskan, menentukan apakah logika sebesar 1,3 MB perlu diport sama sekali | **IT + Underwriting** |
| 13 | **Tipe setiap field** — nol rule properti, 150 ekspresi diikat kontrol yang bertentangan, korpus nol DDL | **DBA** |

---

## Further Notes

### Empat jebakan membaca korpus — keempatnya sudah pernah menghasilkan angka salah

Siapa pun yang mengaudit ulang angka di spec ini wajib tahu keempatnya:

1. Kondisi rule `When` ada di **dua tag**; membaca satu saja membuat rule tampak tanpa kondisi.
2. Perbandingan berkas antar folder dengan hash mentah **selalu gagal** — metadata ekspor berbeda di
   setiap berkas; buang baris bertanda metadata lebih dulu.
3. Urutan tag ekspor **tidak deterministik**; pengurutan wajib memakai pembanding ordinal, karena ada
   nama rule yang hanya berbeda kapitalisasi.
4. Membaca hanya langkah tingkat atas **kehilangan 63 % logika** — 7.021 dari 11.090 langkah bersarang,
   sampai kedalaman sembilan. Wajib parser rekursif.

### Nama bukan bukti

`CLAUDE.md` §3 butir 3 terbukti berulang kali di korpus ini, tiap kali dengan bentuk berbeda: predikat
bernama fac out yang menguji banding, rule bernama *Insert* yang menghapus lalu menyisipkan ulang, seluruh
tulisan Oracle yang memakai metode **baca**, dan Section bernama Treaty yang berkelas Fac In. **Arah
operasi, arti nilai, dan lingkup tidak pernah disimpulkan dari nama.**

### Perbaikan dipisahkan dari migrasi

Spec ini menemukan banyak hal yang tampak keliru — tautologi pembanding limit, label satuan yang
bertentangan dengan rumusnya, perbandingan angka sebagai string, predikat yang namanya menyesatkan.
**Semuanya direproduksi apa adanya**, diberi komentar sebagai kandidat perbaikan, dan didaftarkan.

Alasannya bukan kehati-hatian abstrak: rekonsiliasi menuntut nol selisih sampai digit terakhir. Satu
"perbaikan" diam-diam menghasilkan selisih yang tidak dapat dijelaskan, dan menyembunyikan salah-port
yang sesungguhnya.

### Tentang path di dokumen ini

Path berkas yang disebut seluruhnya adalah **rule Pega di sistem lama**, dan itu diwajibkan `CLAUDE.md`
§4.6 sebagai ketertelusuran. **Tidak ada path berkas sistem baru** yang disebut — path target berubah
cepat dan akan segera basi.

### Yang membuat spec ini dapat dikerjakan sekarang

Tidak satu pun dari lima modul ini menunggu jawaban yang belum masuk. Itu bukan kebetulan — itu kriteria
pemilihannya. Setiap kali sebuah modul ternyata menyentuh sesuatu yang belum diketahui, bagian itu
dipindahkan ke Out of Scope dengan menyebut apa yang ditunggu, alih-alih diisi dengan asumsi.

---

*Tidak ada nama orang, alamat email, kredensial, maupun data pelanggan di dokumen ini.*
