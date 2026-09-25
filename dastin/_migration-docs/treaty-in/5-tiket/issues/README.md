# Papan tiket — Treaty In

**Tiket: aktif 46 · tertahan 5 · selesai 0 · mati 0** — **51 tiket**, bernomor `14`…`64`

**Yang sebenarnya dapat dimulai: 2** — `14`, `15`. *Dihitung ulang atas **seluruh papan**, bukan per batch.*

> ## ANGKA LAMA *"dapat dimulai: 4"* SALAH, DAN SEBABNYA LAYAK DIBACA
>
> Papan ini pernah mencetak **"dapat dimulai: 4 — `01`, `02`, `03`, `05`"**. **Keempatnya tidak dapat
> dimulai.** Ketiganya menambahkan **kolom** pada `VERSI_KONTRAK` dan satu menambahkan tabel yang
> merujuknya — dan **tidak ada satu pun tiket di papan yang membuat tabelnya**.
>
> **Sebab kekeliruannya, dan inilah yang akan terulang:**
>
> > Angka itu dihitung atas **penghalang di dalam papan**, sementara penghalangnya ada **di luar** —
> > di lingkup yang ronde itu belum digarap.
>
> Selama sebuah batch hanya memuat sebagian kemampuan sebuah modul, **"tidak ada penghalang di papan"
> tidak berarti "tidak ada penghalang"**. Papan yang tidak lengkap **selalu** melaporkan angka yang
> terlalu besar, dan ia melaporkannya dengan tenang.
>
> **Cara memeriksanya, sebelum mencetak angka itu di batch berikutnya:** untuk tiap tiket yang mengaku
> dapat dimulai, sebutkan **entitas yang disentuhnya** dan tunjukkan **tiket mana yang membuatnya**.
> Bila jawabannya "tidak ada tiket, tabelnya sudah ada", itu **andaian** — dan andaian itulah yang
> keliru di sini.

> **Lingkup papan ini: kemampuan Treaty In saja** — batch 1 (`14`…`44`) dan batch 2 (`45`…`64`).
>
> ### PEMISAHAN PAPAN — 25 September 2026
>
> Tiket **`01`…`13`** (kemampuan Adjustment) **dipindahkan** ke
> [`../../../treaty-in-adjustment/5-tiket/issues/`](../../../treaty-in-adjustment/5-tiket/issues/)
> atas perintah pemilik proses. **Ini membalik `GRL-01`** — *"satu model, satu spesifikasi, satu
> papan"* — dan pembalikannya dicatat, bukan disamarkan.
>
> **Yang dibalik hanya papannya.** Satu model data, satu daftar invarian, satu
> `DAFTAR-PEKERJAAN.md`, dan satu `ASUMSI-CLEAR.md` **tetap**. `P-60`…`P-66` **tidak ikut pindah**;
> ia tetap di `DAFTAR-PEKERJAAN.md` papan ini.
>
> **Penomoran tidak disusun ulang.** `14`…`64` tetap, `01`…`13` tetap — tidak ada tabrakan nomor
> antar-papan, dan setiap rujukan yang sudah ada tetap sah.
> Dari **66 kemampuan** di [`../DAFTAR-PEKERJAAN.md`](../DAFTAR-PEKERJAAN.md). **Batch 2 — 22 butir
> yang belum ditiketkan** (sapuan mekanis atas medan `*Asal:*`, bukan hitungan tangan).
>
> ### RALAT 25 September 2026 — penahan batch 2 BUKAN `REV-3`
>
> Baris ini pernah berbunyi: *"20 butir yang menyentuh daftar keadaan — belum ditiketkan;
> **penahannya tanggapan `REV-3`**."* **Itu keliru, dan kekeliruannya menyebar ke tiga berkas.**
>
> `REV-3` menyatakan sebaliknya di barisnya sendiri: keputusan ADR-0055 — *"tidak ada satu pun jalan
> menyetel keadaan selain melalui perpindahan di daftar"* — **tetap, seluruhnya**. Yang diusulkan
> berubah hanya **tabel §4**, dan keempat butirnya **koreksi atas pemerian sistem lama**. Tidak satu
> pun mengubah daftar keadaan, kolom keadaan, maupun perpindahannya. `KTV-4` sudah mengadilinya:
> ***"nol dari enam menentukan letak kolom."***
>
> | Salah | Benar |
> |---|---|
> | *"batch 2 ditunda, menunggu `REV-3`"* | *"batch 2 dikerjakan sesi berikutnya, sebab 51 kemampuan tidak muat satu sesi"* |
>
> Yang pertama membuat orang **menunggu jawaban yang tidak akan mengubah apa pun**; yang kedua
> membuat mereka **menjadwalkan sesi**. **Pemisahan batchnya tidak dicabut** — ia benar sebagai
> ukuran sesi. `REV` tetap diserahkan: ia **utang dokumentasi ADR induk**, bukan utang tiket.

## Papan

### Tiket `01`…`13` TIDAK lagi di papan ini

Ketiga belasnya pindah ke papan Adjustment. **Satu tepi tetap melintas ke sini dan tidak boleh
hilang:**

| Tiket papan ini | Diblokir | Yang diantarkan tiket itu |
|---|---|---|
| **`40`** — nilai versi sebelumnya berdampingan, di-SELECT bukan disalin | **`01`** ⟦papan Adjustment⟧ | rujukan versi dasar yang disimpan eksplisit |

Dan **lima tepi ke arah sebaliknya**: `01`, `02`, `03`, `04`, `05` di papan Adjustment seluruhnya
diblokir **`14`** di papan ini. Itu sebab papan Adjustment mencetak **"dapat dimulai: 0"**.

### Batch 1 Treaty In — `14`…`44`

| # | Tiket | Ditahan oleh |
|---|---|---|
| `14` | **Kontrak dan versi pertamanya berdiri**, dan dapat ditemukan kembali dengan pengenalnya | — |
| `15` | **Himpunan acuan bertambah** tanpa mengubah arti apa pun yang sudah tercatat | — |
| `16` | Peringatan kunci alami ganda tanpa menghalangi simpan | `14` |
| `17` | Kontrak ditemukan kembali lewat nomor warisan | `14` |
| `18` | Kunci alami kontrak beku sesudah kontraknya lahir | `14` |
| `19` | Lapisan beku tidak berubah antar versi | `14` |
| `20` | Lebih dari satu mata uang berlaku, masing-masing berkurs dan berperiode | `14` · `15` |
| `21` | Kegagalan hitung menghasilkan keterangan, bukan nol | `20` |
| `22` | Retensi cedant per kelompok treaty per mata uang | `14` · `15` |
| `23` | EGNPI per kelompok treaty per mata uang | `14` · `15` |
| `24` | Portofolio masuk dan keluar yang menyertai kontrak | `14` |
| `25` | Periode pelaporan beserta jatuh temponya | `14` |
| `26` | Periode akumulasi, ditolak bila di luar periode kontrak | `14` |
| `27` | Termin pembayaran premi bernomor urut, per mata uang | `14` · `15` |
| `28` | Skala ko-asuransi sebagai beberapa baris | `14` |
| `29` | Batas tanggungan untuk bahaya apa pun di daftar | `14` · `15` |
| `30` | Dokumen dilampirkan, dirujuk bukan disalin | `14` |
| `31` | Layer beserta limit, deductible, MDP, dan pemulihan limitnya | `14` · `15` |
| `32` | Syarat berbeda tiap pemulihan limit | `31` |
| `33` | Bagian NuRe atas layer, hanya pada non-proporsional | `31` |
| `34` | Ketentuan proporsional per kelompok treaty di dalam layer | `31` · `15` |
| `35` | Tepat satu dari persen quota share atau jumlah lines surplus | `34` |
| `36` | Baris surplus tanpa quota share menerima kegagalan yang menyebut apa | `34` |
| `37` | Potongan atas premi, aturan sama di kedua pelekatannya | `33` · `34` · `15` |
| `38` | Penyebaran bagian NuRe per jenis reasuransi per mata uang | `33` · `34` · `15` · **`Uji AD`** · **`L-3`** |
| `39` | Jejak perubahan: siapa mengubah fakta apa dan kapan | `14` |
| `40` | Nilai versi sebelumnya berdampingan, di-SELECT bukan disalin | `14` · **`01`** ⟦papan Adjustment⟧ |
| `41` | Konsumen hilir membaca identitas kontrak dalam bentuk lama | `14` |
| `42` | Arsip JSON sistem lama disimpan dan dapat ditunjukkan | `14` |
| `43` | Sakelar penegakan trigger selama pemindahan | `14` |
| `44` | Isi tabel acuan dipindahkan dari sistem lama | `15` · `43` · **`KTV-A`** |

### Dapat dimulai hari pertama

**`14` dan `15` — dua, dan keduanya tidak bergantung pada tiket mana pun.**

Keduanya **PEMBUAT PERTAMA**, dan itu bukan kebetulan: `14` membuat `KONTRAK` dan `VERSI_KONTRAK`
yang **seluruh papan** menggantung padanya, `15` membuat keenam tabel acuan yang **setiap kunci
asing** menunggunya. Begitu keduanya mendarat, **yang dapat dimulai melompat ke tiga belas**.

> **Taksiran keduanya tidak dibandingkan dengan tiket mana pun**, sesuai aturan `PEMBUAT PERTAMA`.

### Menunggui, tidak menahan

Tiket ini **boleh jalan**. Jawaban yang ditunggu mengubah satu constraint, satu lebar kolom, atau satu
jalur migrasi — bukan rancangannya.

| # | Menunggu | Yang berubah bila jawabannya lain |
|---|---|---|
| `01` | `REV-3` | **keadaan mana yang dihitung "berlaku"** — bukan bahwa ada penunjuk dasar |
| `02` | `DB-20` | **titik beku** materialitas, dan titik beku milik `07` |
| `09` | `DB-16a` | **bentuk kolom** dokumennya, yang ikut `04` |
| `14`, `31` | `T-6` | tiga kolom mata uang kosong pada `VERSI_KONTRAK`, dua pada `LAYER`. **Dicabut**, bukan diubah |
| `23` | `T-4` · `29` `T-2` · `31` `T-3` | **tingkat pencatatan** paket uangnya — satu constraint, bukan bentuk tabel |
| `31` | **`F-15`** | `PERSEN_ROL` dapat **dicabut**, dan tabel anak premi diperoleh dapat **ditambahkan**. Bentuk `LAYER` tidak berubah |
| `39` | **`F-16`** | **dari mana** nilai `PERAN_PELAKU` diambil. **Bentuknya sudah diputuskan** `KTV-D` — potret, bukan rujukan — dan bentuk itu tidak berubah |
| `43` | `L-3` | **ujinya**, bukan bentuk sakelarnya |
| enam belas **PEMBUAT PERTAMA** | `KTV-A` | **lebar kolom**. Dapat dipersempit **hanya sebelum data dimuat** — dan yang memuat data pertama adalah `44` |

### Menahan sungguhan — lima

| # | Penahan | Kenapa ia menahan, bukan menunggui |
|---|---|---|
| `04` | `DB-16a` | bila dokumen ternyata dikirim ke luar, **jenis bernilai tiga** dan penanda kembali. Itu kolom, bukan penyetelan |
| `07` | `DB-20` | titik beku adalah **pokok** tiketnya |
| `08` | `DB-16b` | tanggalnya **adalah** pokok tiketnya, dan `KTV-2` menuntut kolomnya **dicabut sebelum data masuk** bila dibantah |
| `38` | `Uji AD` · `L-3` | `Uji AD` menentukan **golongannya**, dan bila **B** tiket ini menuntut medan `CARA MENYALAKANNYA` yang belum ada. `L-3` berarti uji negatif MV-nya **belum pernah dijalankan di mana pun** |
| `44` | `KTV-A` | **bertenggat**: ia yang **memuat data pertama**, dan sesudahnya presisi tidak dapat dipersempit lagi |

**Empat pertanyaan bisnis di antaranya — `DB-16a`, `DB-16b`, `DB-20` — belum dikirim.**

---

## UKURAN SELESAI DI PAPAN INI — dibaca sebelum menilai satu tiket pun

Skill `to-tickets` menuntut tiap tiket memotong tegak **skema · API · UI · uji**. **Di papan ini, dua
dari empat lapisan itu belum ada.**

`L-4` menyatakan: *"tidak ada spesifikasi layar di mana pun."* Lubang itu sudah punya pemilik —
**batch lapisan aplikasi** — dan **bukan batch ini**.

> ### JANGAN MENGARANG LAYAR UNTUK MEMENUHI BENTUK IRISAN TEGAK.
>
> Layar yang dikarang akan dipakai sebagai kriteria selesai oleh orang yang **tidak tahu ia
> karangan** — dan ia akan membangun yang salah dengan keyakinan penuh.

Yang berlaku di sini gerbang **`G1`** modul ini: **pembagian tiket memakai kemampuan, bukan layar.**
Irisan tegak di papan ini berarti tegak melalui **lapisan yang ada**:

```
skema  →  constraint dan invarian  →  jalur migrasi  →  uji negatif DAN positif
```

| | |
|---|---|
| **Sebuah tiket selesai bila** | perilakunya **dapat gagal** dan **dapat diperiksa** |
| **Bukan bila** | ia dapat diperagakan di layar |

### Dan satu kewajiban yang lahir dari ukuran ini

Karena "dapat diperagakan" tidak tersedia sebagai bukti, **uji menjadi satu-satunya bukti**. Maka tiap
tiket membawa **uji negatif dan uji positif**, bukan hanya yang pertama.

> Alasannya bukan kesetaraan: ketiga kekeliruan lingkup di modul ini **lolos uji negatif**. Yang
> menangkapnya hanya uji positif — uji yang menuntut **data sah diterima**. Constraint yang menolak
> terlalu banyak lulus setiap uji negatif yang pernah ditulis untuknya.

---

## Aturan papan

- **Status adalah medan di dalam berkas tiket** (`status: aktif | tertahan | selesai | mati`), bukan
  lokasi foldernya. Papan **dibangkitkan dari berkas**, dan **tidak pernah menghapus tiket**.
- **Empat golongan, tidak ada yang kelima.**
- **Penahan dari luar papan ditulis dengan nama aslinya** — `REV-3`, `DB-16b`, `Uji AD`, `L-3`,
  `KTV-A`, `F-15`, `F-16`, `T-2`…`T-6` — tidak diterjemahkan menjadi nomor tiket.
- **Penahan yang selesai dicoret, tidak dihapus**, supaya rantainya tetap terbaca:
  `~~05~~ *(selesai)*`.
- **Nomor yang lompat bukan kekeliruan.** Penomoran tidak disusun ulang.
- **"Menunggui, tidak menahan" dipisahkan dari penahan sungguhan**, supaya papan tidak berteriak
  serigala.
- **Pencacah diberi label bendanya.** Papan ini mencacah **tiket**; kemampuan `P-nn` punya
  pencacahnya sendiri di `DAFTAR-PEKERJAAN.md`, dan keduanya **tidak pernah dijumlahkan**.
- **Tiket kemampuan yang DIUBAH menyebut bunyi lamanya** di medan `Menggantikan:` — bukan hanya kode
  cacatnya.
- **Satu kode per `DIASUMSIKAN-CLEAR(...)`.** Penanda berisi dua kode **tidak dapat disapu**; tiket
  `07` sempat memuat `(DB-20, REV-3)` dan **sudah dipecah** menjadi dua baris, 24 September 2026.

## Satu angka yang menggambarkan keadaan, dan ia bukan "aktif"

Papan ini mencetak **"Yang sebenarnya dapat dimulai"** — tiket yang **tidak tertahan** dan **seluruh
penghalangnya, di dalam maupun di luar papan, sudah selesai**.

Angka itu, bukan cacah "aktif", yang menjawab *"apakah batch ini layak dimulai sekarang"*. Hari ini
ia **2 dari 44** — dan itu keadaan yang sehat, bukan buruk: keduanya fondasi yang selebihnya
menunggu.

## Asumsi yang dianggap clear

Sebagian tiket bersandar pada penahan yang **dianggap clear atas perintah pemilik proses**, bukan yang
benar-benar tertutup. Masing-masing ditulis di medan `Dasar:` sebagai `DIASUMSIKAN-CLEAR(<kode>)`, dan
daftarnya beserta **apa yang ditinjau ulang bila ia patah** ada di
[`../ASUMSI-CLEAR.md`](../ASUMSI-CLEAR.md).

Satu sapuan atas `DIASUMSIKAN-CLEAR(` di folder ini mengeluarkan daftarnya. **Tidak ada yang perlu
diingat.**

## Satu ketaksamaan bentuk yang tercatat, bukan ditambal

Ketiga belas tiket `01`…`13` **tidak punya medan `golongan:` di frontmatter**, sementara
`BENTUK-TIKET.md` §6.1 mewajibkannya dan tiket `14`…`44` memuatnya. Golongan sebuah tiket —
pelestarian, perubahan, atau baru — **menentukan cara mengujinya**, dan tanpa medan itu pembacanya
harus menyimpulkannya dari isi.

**Tidak ditambal dari sini:** mengisinya berarti **memutuskan golongan tiket yang sudah disetujui**,
dan itu wewenang pemilik proses. **Ditagih** pada putaran berikutnya yang menyentuh ketiga belas tiket
itu.


---

## Batch 2 Treaty In — `45`…`64` · daur hidup dan warisan

| # | Tiket | Ditahan oleh |
|---|---|---|
| `45` | **Daftar keadaan dan perpindahan sah berdiri**; perpindahan di luar daftar ditolak | `14` |
| `46` | Versi terminal beku, menjangkau **seluruh entitas anaknya** | `45` |
| `47` | Paling banyak satu versi tak-terminal per kontrak | `45` |
| `48` | Pengajuan menolak dua syarat, memperingatkan enam | `45` |
| `49` | **SH** menyetujui, versi pindah ke antrian DH | `48` · `54` |
| `50` | **DH** menyetujui, versi pindah ke antrian DR | `49` |
| `51` | **DR** menyetujui, versi menjadi `DISETUJUI` bila seluruh `K2` | `50` |
| `52` | Pengaju tidak menyetujui; tiap tingkat orang berbeda; tanpa jalan pintas | `51` |
| `53` | Pengembalian ke `DRAFT` dengan alasan wajib | `49` |
| `54` | **Tiap perpindahan meninggalkan catatan persetujuan** | `45` |
| `55` | Penolakan versi — kontrak maupun addendum — meninggalkan catatan | `49` · `54` |
| `56` | Pembatalan draf oleh pembuatnya sendiri | `45` · `54` |
| `57` | Angka rupiah pada versi disetujui tidak bergeser | `46` · `20` |
| `58` | Wewenang persetujuan dibatasi nilai kontrak | `51` · **dokumen batas wewenang** |
| `59` | Kepala kontrak warisan dipindahkan beserta keadaannya | `45` · `44` |
| `60` | Keadaan warisan tanpa padanan → `WARISAN_TAK_TERPETAKAN` | `59` |
| `61` | Baris warisan diperbaiki ke keadaan sah yang dipilih eksplisit | `60` |
| `62` | Kontrak dicari, termasuk menurut keadaan siklus hidupnya | `45` |
| `63` | Cara pembukuan XOL, hanya pada non-proporsional | `14` · **`Uji X-2`** |
| `64` | Tabel acuan pembagian kapasitas punya pemiliknya | `15` · **penetapan pemilik** |

### Tiga PEMBUAT PERTAMA baru di batch ini

| Tiket | Entitas |
|---|---|
| `54` | **`CATATAN_PERSETUJUAN`** — ditemukan lewat pemeriksaan entitas-per-irisan; tanpa pemeriksaan itu, `49` akan menulis catatan ke tabel yang belum ada |
| `45` | himpunan nilai `KEADAAN_SIKLUS_HIDUP` — **kolomnya** dibuat `14`, **artinya** dibuat di sini |
| `59` | isi kepala kontrak warisan — data, bukan tabel |

### Kenapa `49`, `50`, `51` tiga tiket dan bukan satu

`U-2` ditegakkan harfiah: satu tiket tidak melintasi lebih dari **satu** perpindahan, dan `GRL-08`
menulis `SETUJUI×3` — **tiga perpindahan berbeda**, bukan satu perpindahan berparameter.

> Satu tiket yang melintasi tiga tingkat **gagal sebagai satu blok**. Bila tingkat kedua yang salah,
> tiketnya *"belum selesai"* — dan **tidak ada yang tahu tingkat mana yang patah**.

Dan wewenangnya berbeda per tingkat: `ADR-0044` menuntut **setiap larangan punya tepat satu sebab**,
sehingga penolakan karena **peran** dan penolakan karena **keadaan** harus berpesan berbeda — tiga
kali, dengan jawaban yang berbeda tiap tingkat.

---

## Satu entitas yang TIDAK punya tiket, dan tidak punya kemampuan

Pemeriksaan entitas-per-irisan atas **seluruh 35 entitas** `KAMUS-KOLOM.md` terhadap **64 tiket**
menemukan satu yatim:

| Entitas | Keadaan |
|---|---|
| **`PERISTIWA_KONTRAK`** | **nol tiket menyebutnya, dan nol kemampuan `P-nn` menghasilkannya** |

Ia berdiri di `SPEC-MODEL-DATA.md` §10.21a dan di `ddl-usulan/`, lahir dari pemisahan `CommentList`
sistem lama menjadi **dua** entitas: persetujuan (`IsApproved` terisi) dan **peristiwa**
(`IsApproved` kosong — *"Copied from ID …"*, *"Had Created Internal Edit"*).

> **Tabelnya akan berdiri kosong selamanya**, dan tidak ada pencacah yang memperlihatkannya. Ia
> bentuk `K-3` satu tingkat lebih dalam: di sana **kemampuan** jatuh di antara dua batch; di sini
> **entitas** tidak punya kemampuan sama sekali.

| | |
|---|---|
| **Siapa menutup** | pemilik proses — satu kalimat: kemampuan `P-67` tersendiri, atau ia ikut tiket `59` sebagai bagian pemindahan warisan |
| **Yang menagih** | berkas ini, dan `ddl-usulan/` yang memuat tabelnya |
