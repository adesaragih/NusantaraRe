# Usulan irisan — Treaty In, batch 1

**Tanggal:** 24 September 2026 · **Langkah 2 to-ticket, ronde Treaty In**
**Status: DISETUJUI 24 September 2026, dengan lima jawaban.** Tiket `14`…`44` **sudah ditulis** ke
[`issues/`](issues/).

> ## LIMA JAWABAN PEMILIK PROSES, DAN AKIBATNYA
>
> **Daftar di bawah dipertahankan apa adanya sebagai bentuk yang disetujui.** Yang berubah dicatat di
> sini, bukan disunting diam-diam ke dalam tabelnya.
>
> | # | Pertanyaan | Jawaban | Yang berubah pada daftar di bawah |
> |---:|---|---|---|
> | 1 | irisan `14` sebagai **PEMBUAT PERTAMA** | **YA** — *"kolom tidak dapat ditambahkan ke tabel yang belum ada"* | tidak ada. **Tiket `01`…`05` memperoleh `Blocked by: 14`**, tiap suntingan bertanggal dan bersebab |
> | 2 | `P-01` dan `P-50` dibelah antar batch | **YA**, dengan **syarat**: belahannya ditulis **di baris `P-01` dan `P-50` sendiri** di `DAFTAR-PEKERJAAN.md` | syarat itu dipenuhi — **§2.4** berkas itu |
> | 3 | ketiga `G2` tegas tidak masuk | **YA**, tetapi **penandanya harus diperbaiki**: `TERTAHAN` dan "tegas tidak masuk" **sumbu yang berbeda** | **§4.1** `DAFTAR-PEKERJAAN.md`, dan pola ini naik sebagai **usulan `TA-12`** |
> | 4 | granularitas 31 irisan | **PERTAHANKAN 31** — keseragaman `22`…`30` justru tanda daftar anak versi memang seragam, dan menggabungkannya gagal syarat *"muat dalam satu jendela konteks"* | tidak ada |
> | 5 | irisan `39` **MENAHAN** atas `F-16` | **penahannya diubah** — yang menahan bukan ketiadaan entitas peran melainkan **bentuk `PERAN_PELAKU`**, dan itu **diputuskan**: ia **POTRET**, `KTV-D` | **`39` turun menjadi *menunggui***. `MENAHAN` menjadi **tiga**: `38`, `44`, dan — di papan penuh — `04`, `07`, `08` |
>
> ### Satu hal yang jawaban 5 tambahkan, dan ia lebih tajam daripada pertanyaannya
>
> Bedanya **permanen**, sebab jejak perubahan **tambah-saja** dan `ADR-0045` menyebutnya fakta mesin:
>
> > **Baris pertama yang ditulis dengan bentuk salah tidak dapat diperbaiki tanpa menulis ulang
> > sejarah** — hal yang jejak ini ada untuk mencegahnya.
>
> Maka `KTV-D` bertenggat pada **baris jejak pertama**, bukan pada pemuatan data pertama — lebih awal
> daripada `KTV-A` dan `KTV-2`. Tercatat di `ASUMSI-CLEAR.md` §3a.
**Nomor tiket:** mulai **14** — penomoran tidak disusun ulang.

> **Usulan yang menunggu persetujuan ADALAH keluaran.** Ia ditulis ke disk karena usulan yang hanya
> hidup di balasan **hilang begitu sesinya tutup**, dan yang hilang akan disusun ulang dari awal oleh
> orang yang tidak tahu ia pernah ada.

> ## SATU TEMUAN YANG MENGUBAH PAPAN YANG SUDAH ADA — dibaca sebelum daftar irisan
>
> **Tidak ada satu pun dari tiga belas tiket di papan yang membuat `VERSI_KONTRAK`.**
>
> Tiket `01` menambahkan **kolom** `ID_VERSI_KONTRAK_DASAR` *"pada `VERSI_KONTRAK`"*; tiket `02`
> menambahkan kolom `SIFAT_MATERIAL_ADDENDUM`; tiket `05` menambahkan kolom nomor urut. Ketiganya
> **mengandaikan tabelnya sudah ada**, dan **PEMBUAT PERTAMA**-nya tidak pernah ditiketkan — sebab
> ia milik `P-01`, kemampuan Treaty In, yang ronde lalu memang di luar lingkup.
>
> | | |
> |---|---|
> | **Akibatnya pada angka papan** | `issues/README.md` mencetak *"Yang sebenarnya dapat dimulai: 4 — `01`, `02`, `03`, `05`"*. **Keempatnya tidak dapat dimulai** sampai `KONTRAK` dan `VERSI_KONTRAK` berdiri |
> | **Kenapa ini tidak terlihat ronde lalu** | angka itu dihitung atas **penghalang di dalam papan**, dan penghalangnya memang tidak ada di dalam papan — ia di luar, di lingkup yang belum digarap |
> | **Usulan** | irisan **14** menjadi **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`, dan tiket `01`, `02`, `03`, `04`, `05` memperoleh sisi penghalang **ke 14**. Tiket `06`…`13` sudah terhalang lewat keempatnya |
> | **Yang TIDAK saya lakukan** | menyunting kelima berkas tiket itu. Ia perubahan pada tiket yang sudah disetujui, dan usulan ini yang memintanya |
>
> **Sesudah koreksi ini, "yang sebenarnya dapat dimulai" untuk SELURUH papan adalah DUA** — `14` dan
> `15` — bukan empat.

---

## Perkiraan ongkos ronde ini

| | |
|---|---:|
| token terpakai sampai usulan ini ditulis | ± 150 rb |
| ~~perkiraan menulis 31 berkas tiket bila disetujui~~ | ~~± 260 rb~~ |
| **terpakai sebenarnya** untuk menulis 31 tiket + memperbarui enam berkas keadaan | **± 130 rb** |

*Perkiraannya **dua kali lipat** dari yang terpakai. Sebabnya dapat disebut: taksirannya dibuat atas
"31 berkas" tanpa memisahkan berkas ber-**PEMBUAT PERTAMA** dari tetangganya — persis kekeliruan yang
gerbang `G1` peringatkan untuk taksiran tiket, dan ternyata berlaku juga untuk taksiran token.*

---

# GERBANG 0

## 0.1 Saringan `G2` — berapa masuk, berapa tegas tidak

Dibangkitkan `alat/saring-kemampuan-ke-tiket.py`, yang **melaporkan apa yang ditolaknya**:

```
baris kemampuan terbaca : 65
DITOLAK / DIKELUARKAN -- satu baris per sebab:
   DICORET sebagai SIFAT (U-6)                                    3   (P-24, P-25, P-26)
   kemampuan Adjustment P-60..P-66 -- sudah ditiketkan ronde lalu  7
   baris P-nn berkolom kurang dari enam                           5   (tabel penahan §4, bukan baris kemampuan)
   SISA -- lingkup ronde ini                                     55
```

### Lingkupnya **55 berbaris, bukan 56** — dan yang keenam-puluh-enam memang tidak punya baris

Perintah ronde ini menyebut **56 kemampuan**. Angka itu **benar sebagai cacah kemampuan**, dan
perkakas hanya menemukan **55** — sebabnya sudah tercatat:

> **`P-59` tidak punya baris tabel sama sekali.** Ia diputuskan di `DAFTAR-PEKERJAAN.md` §2.2b
> — *"syarat berbeda tiap pemulihan limit"*, bergolongan **BARU** — dan tidak pernah didaftar di
> tabel kemampuan mana pun. Dilaporkan di `2-to-spec/SISA-DAN-SERAH-TERIMA.md` §12.4.

**`P-59` dibawa masuk lingkup dengan tangan**, bukan lewat perkakas, dan barisnya **tetap tidak
dikarang** — isi kolom Gol, Asal, dan Entitas-nya diambil apa adanya dari §2.2b. Maka
**lingkup = 56**.

### Tiga yang **TEGAS TIDAK MASUK**, dan sebab masing-masing

| Kemampuan | Tepi `G2` yang mengenainya | Sebab, dan kenapa ia bukan sekadar "tertahan" |
|---|---|---|
| **`P-20`** — mencatat cara pembukuan XOL, **hanya pada kontrak non-proporsional** | *"invarian atas `AccountingMode` — menunggu Uji X-2"* | Separuh keduanya **adalah** invariannya. Memotongnya menjadi *"kolomnya ada, aturannya menyusul"* gagal **U-3** — ia kalimat ber-"supaya nanti", dan kolom yang dicatat tanpa aturannya persis bentuk *"kontrak yang setengah dicatat"* yang §2 larang |
| **`P-41`** — wewenang persetujuan dibatasi **nilai kontrak** | *"apa pun yang mewujudkan butir eskalasi yang belum diputuskan"* — **butir 5** | Dokumen batas wewenangnya **belum ada**. `G3` sudah menguji pemecahannya dan mencatat kegagalannya: *"tidak dapat dipecah — yang belum diketahui duduk tepat di tengah kemampuannya"* |
| **`P-49`** — tabel acuan pembagian kapasitas punya **pemilik** | idem — **butir 7** | Kolom Entitas-nya berbunyi *(belum ditetapkan)*. Tidak ada yang dapat dibangun dari kemampuan yang entitasnya belum ada |

> ### Dua pernyataan yang bertentangan, dan yang mana yang menang
>
> `DAFTAR-PEKERJAAN.md` §4 menandai ketiganya **`TERTAHAN`** — yaitu **di dalam** penyerahan pertama,
> menunggu jawaban. `KEPUTUSAN-PEMBAGIAN-TIKET.md` §2 menandainya **tegas tidak masuk** penyerahan
> pertama.
>
> **Gerbang `G2` yang menang**, sebab ia **mengikat** dan berkas daftar pekerjaan tidak. Bedanya
> bukan kata: yang `TERTAHAN` **ditiketkan tanpa taksiran**; yang **tegas tidak masuk** tidak
> ditiketkan sama sekali. Dilaporkan, bukan diselaraskan diam-diam.

### `GEL-2` dan `GEL-3` — **nol kemampuan**, dan itu dinyatakan bukan didiamkan

Sapuan atas `RETRO_KELUAR` dan `PENCAPAIAN` di seluruh 56 baris kemampuan kembali **nol**. Keduanya
memang **entitas** yang ditunda, dan tidak ada satu pun `P-nn` yang membawanya — `DAFTAR-PEKERJAAN.md`
§3 sudah menyatakannya.

**Ditulis di sini justru karena nol.** Angka nol yang tidak dicetak tidak dapat dibedakan dari sapuan
yang tidak pernah dijalankan.

### Hitungan `G2`

| | |
|---|---:|
| kemampuan dalam lingkup ronde ini | **56** |
| **tegas TIDAK masuk** | **3** — `P-20`, `P-41`, `P-49` |
| dari `GEL-2` / `GEL-3` | **0** |
| **masuk, dan digolongkan ke batch** | **53** |

---

## 0.2 Batch 1 versus batch 2 — dan enam koreksi terhadap sapuan kata

Sapuan kata perkakas memberi **28 / 27**. Sapuan kata **tidak dapat membaca arti**, dan batasnya
dicetak di dalam keluarannya sendiri. Enam baris dikoreksi dengan mata; masing-masing disebut supaya
dapat dibantah.

| Kemampuan | Sapuan | Adjudikasi | Sebab |
|---|---|---|---|
| `P-10` · `P-11` | batch 2 | **batch 1** | katanya *"di**tolak** bila periodenya di luar periode kontrak"*. Itu **rentang tanggal**, bukan keadaan siklus hidup. `REV-3` tidak menyentuhnya |
| `P-22` | batch 2 | **batch 1** | *"mengisi keduanya di**tolak**"* — aturan pengisian kolom, bukan perpindahan keadaan |
| `P-18` | batch 2 | **batch 1** | kata pemicunya ada di kalimat penyebaran, bukan pada keadaan. Penahannya `Uji AD` dan `L-3`, bukan `REV-3` |
| `P-45` | batch 2 | **batch 1** | sama dengan `P-10`. **Dan ia sudah ditiketkan** — tiket `03` |
| `P-53` | batch 2 | **batch 1** | sakelar penegakan trigger saat pemindahan. Tidak ada keadaan siklus hidup di dalamnya |
| `P-30` | batch 1 | **batch 2** | ia **saudara kembar `P-29`** — keduanya memeriksa syarat pada perpindahan **pengajuan**. Memisahkan keduanya antar batch berarti separuh gerbang pengajuan dibangun tanpa separuhnya |

### Dua kemampuan yang **TERBELAH antar batch**, dan kenapa itu bukan jalan pintas

`G3` memerintahkan **DIPECAH tepat di garis tempat kepastian berhenti**. Dua kemampuan jatuh persis
di garis itu, dan memaksanya utuh ke salah satu batch merusak keduanya.

#### `P-01` — identitas berdiri di batch 1, keadaan lahirnya di batch 2

| | |
|---|---|
| **Bunyinya** | *"**PK** dapat membuat kontrak baru dan melihatnya kembali dengan pengenalnya; **versi pertamanya lahir langsung dalam keadaan `DRAFT`**"* |
| **Yang menyentuh `REV-3`** | hanya anak kalimat kedua — **nama keadaan lahirnya** |
| **Kalau dipaksa seluruhnya ke batch 2** | `KONTRAK` dan `VERSI_KONTRAK` **tidak pernah dibuat**, dan **setiap** irisan batch 1 menggantung pada tabel yang belum ada. Batch 1 menjadi nol |
| **Kalau dipaksa seluruhnya ke batch 1** | ia membawa `DIASUMSIKAN-CLEAR(REV-3)`, dan perintah ronde ini **melarangnya** untuk batch 1 |
| **Garis pecahnya** | **kolomnya** `KEADAAN_SIKLUS_HIDUP` berdiri di batch 1 — ia kolom, dan tabelnya tidak dapat berdiri tanpanya. **Himpunan nilai sahnya dan aturan perpindahannya** milik batch 2 |

> **Ini bukan penyelundupan.** Yang dibangun batch 1 adalah **tempat**; yang ditunda adalah
> **daftar nilai dan mesin perpindahannya** — dan `REV-3` merevisi persis yang kedua. Bila `REV-3`
> ditolak seluruhnya, irisan 14 **tidak berubah satu kriteria pun**.

#### `P-50` — hanya pecahan tabel acuannya yang batch 1

`P-50` adalah **payung**, dan `DAFTAR-PEKERJAAN.md` sendiri memerintahkan *"DIPECAH per kelompok
entitas **sebelum** penaksiran"*.

| Pecahan | Batch | Sebab |
|---|---|---|
| **isi tabel acuan** — mata uang, jenis potongan, jenis reasuransi, bahaya, kelompok treaty, kelas bisnis | **1** | **nol keadaan ditulis**. Ia murni pemindahan himpunan, dan seluruh kunci asing menunggunya |
| kepala kontrak · daftar anak versi · cabang non-prop · cabang prop · potongan dan penyebaran · catatan dan jejak | **2** | pemindahan kepala **menulis `KEADAAN_SIKLUS_HIDUP`**, termasuk `WARISAN_TAK_TERPETAKAN` (`P-51`). Seluruh pecahan lain menggantung padanya, sehingga tidak satu pun dapat dinyatakan selesai sebelum daftar keadaannya pasti |

> **Kenapa pecahan anak tidak dimasukkan batch 1 lalu ditandai tertahan:** ia akan berupa tiket yang
> penghalangnya **tidak ada di papan mana pun** — bukan tertahan, melainkan menggantung. Perintah
> ronde ini sudah menolak bentuk itu untuk batch 2.

### Hitungan batch — **kedua-duanya**

| | |
|---:|---|
| **53** | kemampuan masuk sesudah `G2` |
| **2** | terbelah antar batch — `P-01`, `P-50` |
| **30** | **batch 1**, utuh *(termasuk `P-45` yang sudah ditiketkan)* |
| **21** | **batch 2**, utuh |
| | |
| **31** | **irisan batch 1 yang diusulkan** — 29 utuh *(30 − `P-45`)* + 2 pecahan |
| **20** | **ditunda ke batch 2** — 18 utuh *(21 − `P-04`, `P-42`, `P-38` yang sudah tersentuh papan)* + 2 pecahan |

### Perbandingan dengan perkiraan perintah

Perintah memperkirakan batch 2 **± 19**; hitungan ini memberi **21 utuh**, atau **20** sesudah yang
sudah tersentuh papan dikeluarkan. **Selisihnya dua ke atas**, dan penyebabnya dapat disebut satu per
satu:

| Kemampuan | Kenapa ia masuk batch 2 padahal kalimatnya tidak berbunyi "keadaan" |
|---|---|
| **`P-27`** | *"kurs dibekukan pada versinya"*. **Titik bekunya adalah sebuah keadaan** — ADR-0036 membekukan **saat disetujui**. Bila `REV-3` mengubah daftar keadaan, titik beku itu ikut berpindah |
| **`P-30`** | saudara kembar `P-29`, lihat tabel koreksi di atas |
| **`P-56`** | pencarian **menurut keadaan siklus hidup** disebut terang di kalimatnya; nilai yang dapat dicari **adalah** daftar yang `REV-3` revisi |

---

# LANGKAH 1 — Adjudikasi batch 1 terhadap tiga belas tiket yang sudah ada

**Tiga golongan, tidak ada yang keempat.**

| Golongan | Berapa | Kemampuan |
|---|---:|---|
| **SUDAH DITIKETKAN** | **1** | `P-45` |
| **BERSINGGUNGAN** | **3** | `P-01` *(pecahan identitas)* · `P-18` · `P-44` |
| **BARU** | **27** | sisanya |
| **jumlah** | **31** | menutup ke jumlah irisan batch 1 |

## SUDAH DITIKETKAN — satu

| Kemampuan | Tiket | Kenapa ia sudah mencakup |
|---|---|---|
| **`P-45`** | **`03`** | Tiket `03` **adalah** `P-45` yang diubah: *"tanggal berlaku **pada versi** ditolak bila di luar periode kontrak"*, dengan pembawanya disebut supaya tidak bertabrakan dengan tanggal dokumen. **Tidak ada sisa** |

## BERSINGGUNGAN — tiga, dan bagian yang tersisa disebut

| Kemampuan | Tiket yang menyentuhnya | Bagian yang **tersisa**, dan itu yang menjadi irisan |
|---|---|---|
| **`P-01`** *(identitas)* | `01`, `02`, `05` menambahkan **kolom** pada `VERSI_KONTRAK` | **tabelnya sendiri** — `KONTRAK` dan `VERSI_KONTRAK`, kunci alaminya, kunci asingnya, dan pengenalnya. **Tidak satu pun tiket membuatnya** |
| **`P-18`** | `11` memasang **pemantau kebasian** untuk `INV-69`/`INV-70` | **pola pemantaunya sudah ada; penyebarannya belum.** Tersisa: `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN`, dan pemantau untuk **`INV-47`, `INV-50`, `INV-51`** — tiga *materialized view* induk yang **tidak satu pun punya pemantau** |
| **`P-44`** | `01` membangun `ID_VERSI_KONTRAK_DASAR` | **sisi bacanya.** `01` menyimpan rujukan dasar; `P-44` menampilkan nilai versi dasar **berdampingan** dengan yang sedang disusun, dan menegakkan bahwa ia **di-SELECT, tidak disalin** (`ADR-0048` butir 2) |

> **Tanpa langkah ini tiket kembar lahir.** Yang paling dekat: `P-44` dan tiket `01` sama-sama dapat
> ditulis sebagai *"versi dasar"*, dan dua tiket dengan kriteria selesai yang bertentangan
> **keduanya tampak sah di papan**.

---

# LANGKAH 2 — Daftar irisan batch 1

## Ukuran selesai di papan ini

Sama dengan ronde Adjustment, dan alasannya tidak berubah: `L-4` menyatakan **tidak ada spesifikasi
layar di mana pun**, dan lubang itu milik **batch lapisan aplikasi**.

```
skema  →  constraint dan invarian  →  jalur migrasi  →  uji negatif DAN positif
```

**Jangan mengarang layar.** Sebuah irisan selesai bila perilakunya **dapat gagal** dan **dapat
diperiksa**, bukan bila ia dapat diperagakan.

## Daftar

| # | Judul | Dari | Blocked by | Yang diantarkan | Penahan luar papan | Menahan / menunggui |
|---:|---|---|---|---|---|---|
| **14** | Kontrak dan versi pertamanya berdiri, dan dapat ditemukan kembali dengan pengenalnya | `P-01` *(pecahan identitas)* | None | `KONTRAK` dan `VERSI_KONTRAK` ada beserta kunci alami, kunci asing, dan pengenalnya; satu kontrak dapat dibuat dan dibaca kembali. Kolom `KEADAAN_SIKLUS_HIDUP` **ada**; **daftar nilai dan perpindahannya tegas bukan bagian irisan ini**. **PEMBUAT PERTAMA** | `KTV-A` · `T-6` | **menunggui** — `KTV-A` mengubah **lebar kolom**, `T-6` mencabut tiga kolom mata uang yang kosong. Tidak satu pun mengubah bentuk tabelnya |
| **15** | Himpunan acuan bertambah tanpa mengubah arti apa pun yang sudah tercatat | `P-48` | None | Keenam tabel acuan ada; menambah baris baru **tidak** menyentuh kontrak yang sudah tercatat; kode ganda ditolak (`INV-68`). **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **16** | Pengisi kontrak diperingatkan bahwa kunci alaminya sama dengan kontrak lain, tanpa dihalangi menyimpan | `P-02` | 14 | Kontrak berkunci alami kembar **tersimpan**, dan peringatannya muncul menyebut kontrak pembandingnya. **Memperingatkan, bukan menolak** (`ADR-0040` §2) | — | — |
| **17** | Kontrak ditemukan kembali lewat nomor warisan sistem lama, berdampingan dengan pengenal baru | `P-03` | 14 | Pencarian dengan nomor warisan menemukan barisnya; baris non-warisan bernomor warisan kosong **tidak** mengganggu pencarian | — | — |
| **18** | Kunci alami kontrak tidak dapat diubah sesudah kontraknya lahir, dan penolakannya menyebut apa yang diubah | `P-05` | 14 | Mengubah cedant, asal bisnis, sifat proporsi, atau tanggal periode pada kontrak yang sudah ada **ditolak** (`INV-19`), dan pesannya menyebut ruas yang dicoba | — | — |
| **19** | Lapisan beku tidak berubah antar versi | `P-43` | 14 | Versi baru pada kontrak yang sama **tidak dapat** membawa cedant, asal bisnis, sifat proporsi, atau periode yang berbeda; percobaannya ditolak | — | — |
| **20** | Lebih dari satu mata uang berlaku pada satu kontrak, masing-masing dengan kurs dan periodenya | `P-06` | 14 · 15 | `MATA_UANG_KONTRAK` ada; dua mata uang pada satu versi diterima; mata uang kembar pada versi yang sama ditolak (`INV-07`). **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **21** | Perhitungan yang gagal menghasilkan keterangan apa yang gagal, bukan angka nol dan bukan kurs satu | `P-28` | 20 | Konversi tanpa kurs **tidak** menghasilkan nol maupun 1; ia menghasilkan keterangan bernama (`ADR-0035`, `INV-38`, `INV-42`) | — | — |
| **22** | Retensi cedant dicatat per kelompok treaty per mata uang | `P-07` | 14 · 15 | `RETENSI_CEDANT` ada; kunci alaminya **memuat mata uang** (`INV-08`); baris kembar ditolak. **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **23** | EGNPI dicatat per kelompok treaty per mata uang | `P-08` | 14 · 15 | `EGNPI` ada, `INV-09` terpasang, baris kembar ditolak. **PEMBUAT PERTAMA** | `KTV-A` · `T-4` | **menunggui** — `T-4` menentukan **tingkat pencatatan** `NILAI_EGNPI`, bukan bentuk tabelnya |
| **24** | Portofolio masuk dan keluar yang menyertai kontrak dicatat | `P-09` | 14 | `PORTOFOLIO` ada; arah masuk dan keluar dapat dibedakan; kunci alaminya ditegakkan (`INV-66`). **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **25** | Periode pelaporan dicatat beserta tanggal jatuh temponya, dan ditolak bila di luar periode kontrak | `P-10` | 14 | `PERIODE_PELAPORAN` ada; periode di luar periode kontrak **ditolak** (`INV-55`); periode kembar ditolak (`INV-10`). **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **26** | Periode akumulasi dicatat, dan ditolak bila di luar periode kontrak | `P-11` | 14 | `PERIODE_AKUMULASI` ada; `INV-11` dan `INV-56` terpasang. **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **27** | Termin pembayaran premi dicatat bernomor urut, per mata uang | `P-12` | 14 · 15 | `TERMIN` ada; kunci alaminya **memuat mata uang** (`INV-12`, §10.16a). **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **28** | Skala ko-asuransi dicatat sebagai beberapa baris, bukan satu nilai | `P-13` | 14 | `SKALA_KOASURANSI` ada; dua baris berbeda persen limit pada satu versi diterima; persen kembar ditolak (`INV-13`). **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **29** | Batas tanggungan dicatat untuk bahaya apa pun yang ada di daftar bahaya, termasuk yang ditambahkan kemudian | `P-14` | 14 · 15 | `BATAS_PER_BAHAYA` ada; bahaya baru di tabel acuan **langsung dapat dipakai tanpa mengubah skema** (`ADR-0038`); `INV-14` terpasang. **PEMBUAT PERTAMA** | `KTV-A` · `T-2` | **menunggui** — `T-2` menentukan tingkat pencatatan `NILAI_BATAS` |
| **30** | Dokumen dilampirkan pada kontrak, dan dokumennya dirujuk bukan disalin | `P-15` | 14 | `DOKUMEN_KONTRAK` ada; yang tersimpan **rujukan**, bukan isi berkas (`ADR-0027`, `INV-59`); `INV-67` terpasang. **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **31** | Layer dicatat beserta limit, deductible, MDP, dan ketentuan pemulihan limitnya | `P-16` | 14 · 15 | `LAYER`, `PEMULIHAN_LIMIT`, `NILAI_MDP`, `NILAI_MDP_MINIMUM` ada; kunci alami layer `INV-05` terpasang; MDP tersimpan **per mata uang**, bukan sebagai satu angka. **PEMBUAT PERTAMA** | `KTV-A` · `T-3` · `T-6` · **`F-15`** | **menunggui** — `F-15` dapat **mencabut `PERSEN_ROL`** dan dapat **menambah tabel anak premi diperoleh**. Bentuk `LAYER` sendiri tidak berubah |
| **32** | Syarat berbeda ditetapkan untuk tiap pemulihan limit pada sebuah layer | `P-59` **BARU** | 31 | Pemulihan pertama dan kedua dapat berpersentase **berbeda**; sistem lama menuliskannya tetap `100` (`TETAPAN-DI-KODE.md` §3). `INV-49` terpasang. Wajib **CARA MENYALAKANNYA** | — | — |
| **33** | Bagian NuRe atas sebuah layer dicatat, dan bagian hanya dapat dibuat pada kontrak non-proporsional | `P-17` | 31 | `BAGIAN`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM` ada; membuat bagian pada kontrak proporsional **ditolak**; `INV-64` terpasang. **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **34** | Ketentuan proporsional dicatat per kelompok treaty di dalam sebuah layer | `P-21` | 31 · 15 | `DETAIL_PROPORSIONAL`, `NILAI_CADANGAN_PREMI` ada; `INV-06` terpasang. **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **35** | Tepat satu dari persen quota share atau jumlah lines surplus terisi, ditentukan jenis treaty-nya | `P-22` | 34 | Mengisi keduanya **ditolak**; mengisi yang tidak sesuai jenis treaty **ditolak**; mengisi yang benar diterima | — | — |
| **36** | Baris surplus tanpa baris quota share pada versi yang sama menghasilkan kegagalan yang menyebutkan apa yang kurang | `P-23` | 34 | Menyimpan surplus tanpa quota share **ditolak**, dan pesannya menyebut **apa** yang kurang — bukan galat tanpa isi | — | — |
| **37** | Potongan atas premi dicatat, dengan aturan yang sama di kedua tempat ia melekat | `P-19` | 33 · 34 · 15 | `POTONGAN` ada dengan **dua kolom induk bernama** dan `CHECK` tepat satu terisi (`KTV-B`); rumus potongannya **tertulis sekali** dan berlaku di kedua pelekatan (`INV-63`); `INV-15` terpasang sebagai **dua** `UNIQUE`. **PEMBUAT PERTAMA** | `KTV-A` | **menunggui** |
| **38** | Bagian NuRe disebarkan ke susunan retro internal per jenis reasuransi per mata uang | `P-18` | 33 · 34 · 15 | `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN` ada; `INV-16` sebagai dua `UNIQUE`; **`INV-47`, `INV-50`, dan `INV-51` terpasang BESERTA pemantau kebasiannya** — ketiganya *materialized view*, dan **tidak satu pun punya pemantau hari ini**. **PEMBUAT PERTAMA** | **`Uji AD`** · **`L-3`** · `KTV-A` | **MENAHAN** — `Uji AD` menentukan **golongannya**, `L` atau `B`; bila `B`, irisan ini **menuntut medan `CARA MENYALAKANNYA`** yang belum ada. `L-3` berarti uji negatif MV-nya **belum pernah dijalankan di Oracle mana pun** |
| **39** | Pemeriksa jejak melihat siapa mengubah fakta apa dan kapan, untuk kontrak mana pun | `P-47` | 14 | `JEJAK_PERUBAHAN` ada dan terisi pada setiap penulisan; nilai sebelum dan sesudah tersimpan; `PELAKU` dan `PERAN_PELAKU` **dibekukan sebagai teks**. **PEMBUAT PERTAMA** | **`F-16`** · `KTV-A` | **MENAHAN** — `ADR-0045` menuntut **peran yang berlaku saat itu**, dan `ADR-0044` menaruh entitas peran serta penugasan bertanggal di gelombang 1. **§10 tidak memuat keduanya**, sehingga `PERAN_PELAKU` akan terisi nilai yang tidak dapat diperiksa siapa pun |
| **40** | Pengisi kontrak melihat nilai versi sebelumnya berdampingan dengan versi yang sedang disusun | `P-44` | 14 · **`01`** | Nilai versi dasar terlihat berdampingan, dan ia **di-SELECT lewat `ID_VERSI_KONTRAK_DASAR`, tidak disalin** (`ADR-0048` butir 2); tidak ada kolom salinan yang dibangun | — | — |
| **41** | Konsumen hilir yang tidak dikenal tetap dapat membaca identitas kontrak dalam bentuk lama | `P-54` | 14 | Bentuk baca lama tersedia dan menghasilkan nilai yang sama dengan sistem lama untuk kontrak yang sama (`ADR-0051`) | — | — |
| **42** | Pelaksana migrasi menyimpan dokumen JSON sistem lama sebagai arsip, dan dapat menunjukkan arsipnya ada | `P-55` | 14 | Arsip tersimpan saat pemindahan dan keberadaannya dapat ditunjukkan per kontrak; **`INV-61` tetap berlaku** — arsip **tidak punya jalur baca dari aplikasi** (`ADR-0034`) | — | — |
| **43** | Pelaksana migrasi mematikan penegakan trigger selama pemindahan dan menyalakannya kembali, dan keadaan sakelarnya terlihat | `P-53` **BARU** | 14 | Sakelar ada; keadaannya **terbaca tanpa membuka basis data**; menyalakan kembali **memeriksa ulang** baris yang masuk selama ia mati. Wajib **CARA MENYALAKANNYA** | `L-3` | **menunggui** — bentuk sakelarnya tidak bergantung pada instans; **ujinya** bergantung |
| **44** | Isi keenam tabel acuan dipindahkan dari sistem lama | `P-50` *(pecahan tabel acuan)* | 15 · 43 | Keenam himpunan terisi dari sumber lama **apa adanya**, tanpa menghitung ulang apa pun (`ADR-0042`); baris yang tidak dapat dipetakan menjadi **pengecualian migrasi bernomor**, bukan baris yang didiamkan | `KTV-A` | **MENAHAN** — `KTV-A` **bertenggat**: presisi boleh dipersempit **hanya sebelum data dimuat**, dan irisan inilah yang memuat data pertama |

## Hitungan

| | |
|---|---:|
| irisan batch 1 | **31** — bernomor `14`…`44`, tanpa lompat |
| kemampuan yang diwakilinya | **31**, seluruhnya berbeda |
| **PEMBUAT PERTAMA** | **18** |
| **MENAHAN** | **3** — `38`, `39`, `44` |
| menunggui, boleh jalan | **17** |
| tanpa penahan luar sama sekali | **11** |
| **Yang sebenarnya dapat dimulai** | **2** — `14` dan `15` |

*Keenam angka di atas dihitung ulang dari tabelnya sendiri, bukan ditulis tangan.*

> **Dua dari tiga puluh satu, dan itu angka yang menggambarkan keadaan.** Bukan kabar buruk: irisan
> `14` adalah **PEMBUAT PERTAMA** dua entitas yang **seluruh papan** menggantung padanya, dan `15`
> adalah tabel acuan yang setiap kunci asing menunggunya. Begitu keduanya mendarat, **yang dapat
> dimulai melompat ke tiga belas**.

## Sisi penghalang ke papan yang sudah ada

**Ini yang perintah minta diperiksa, dan hasilnya tidak nol.**

| Tiket yang sudah ada | Memperoleh penghalang baru | Sebab |
|---|---|---|
| `01`, `02`, `03`, `04`, `05` | **`14`** | kelimanya menambahkan kolom atau tabel anak pada `VERSI_KONTRAK`, yang `14` buat |
| `06`…`13` | — | sudah terhalang lewat kelima di atas |
| `11` | **`38`** *(usulan)* | `11` memasang pemantau kebasian untuk `INV-69`/`INV-70`; `38` memasangnya untuk `INV-47`/`INV-50`/`INV-51`. **Polanya harus satu**, dan yang menetapkan polanya sebaiknya yang lebih dulu mendarat — **usulan saya: `11` tetap lebih dulu, dan `38` mengikutinya**, bukan sebaliknya |

Dan satu arah sebaliknya:

| Irisan baru | Terhalang tiket lama | Sebab |
|---|---|---|
| `40` | **`01`** | `01` membangun `ID_VERSI_KONTRAK_DASAR`; `40` membacanya |

## `DIASUMSIKAN-CLEAR` — satu kode per penanda

Kewajiban 3 ronde ini: **satu kode per penanda**, supaya sapuannya bekerja. Bentuk yang dipakai:

```
Dasar:  EVIDENCED(...)
        DECIDED(KTV-B)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-6)
```

| Kode | Irisan yang bersandar padanya | Kenapa |
|---|---|---|
| `KTV-A` | kedelapan belas **PEMBUAT PERTAMA**, dan `44` | presisi kolom. **Bertenggat** — `44` memuat data pertama |
| `T-6` | `14`, `31` | kelima kolom mata uang golongan C duduk di `VERSI_KONTRAK` dan `LAYER` |
| `F-15` | `31` | nasib `PERSEN_ROL` dan premi diperoleh |
| `F-16` | `39` | ketiadaan entitas peran |
| `Uji AD` · `L-3` | `38` | golongan `P-18` dan uji negatif MV |
| `T-2` · `T-3` · `T-4` | `29`, `31`, `23` | tingkat pencatatan tiga paket uang |

**`KTV-B` TIDAK diberi kode `DIASUMSIKAN-CLEAR`**, dan itu keputusan: syarat pembalikannya bukan
jawaban yang ditunggu melainkan **keadaan yang mungkin terjadi kelak** — induk ketiga. Ia masuk
`Dasar:` sebagai `DECIDED(KTV-B)`. Memberi kode asumsi pada sesuatu yang tidak ditunggu siapa pun
membuat sapuan asumsi mengeluarkan baris yang tidak pernah dapat ditutup.

## Yang TIDAK menjadi irisan, dan sebabnya

| Hal | Kenapa bukan irisan |
|---|---|
| bentuk fisik **paket uang** | **sifat** (`U-6`) — dibawa tiap irisan sebagai invarian; bentuk fisiknya dibangun **PEMBUAT PERTAMA** |
| **presisi kolom** | sama — ia `KTV-A`, dibawa sebagai penanda, bukan sebagai pekerjaan |
| **pemantau kebasian** sebagai tiket tersendiri | gagal `U-1` — tidak ada pelaku di luar tiket. Ia **kriteria selesai** di dalam `38`, persis seperti di `11` |
| daftar nilai keadaan dan mesin perpindahannya | **batch 2** — ia yang `REV-3` revisi |
| **mesin pro rata** | `GRL-15` — sengaja tidak dibangun. Masuk **YANG TEGAS BUKAN BAGIAN TIKET INI** pada irisan `19`, sebagai pernyataan keputusan |
| pemindahan kepala kontrak dan anaknya | **batch 2** — ia menulis keadaan |

---

## Tiga angka basi yang ditemukan saat menghitung, dan tidak ditambal dari sini

| Berkas | Bunyinya | Yang dihitung ulang |
|---|---|---|
| `DAFTAR-PEKERJAAN.md` §2 judul | *"60 didaftar, 57 aktif"* | **65 berbaris, 3 dicoret, 62 aktif** — dan `P-59` tanpa baris membuatnya **66 / 63** sebagai kemampuan |
| `DAFTAR-PEKERJAAN.md` §4 Hitungan | *"58 didaftar, 55 aktif"* | idem |
| `DAFTAR-PEKERJAAN.md` §5.1 / §5.3 | *"21 entitas milik Treaty In · 27 tersentuh"* | **29 + 6 acuan = 35**, cocok dengan `2-to-spec/KAMUS-KOLOM.md` sesudah `F-18` didaftarkan |

**Ditagih di gerbang pembuka ronde ini, dan dicatat di sini** — `2-to-spec/SISA-DAN-SERAH-TERIMA.md`
§12.4 sudah menunjuk ronde inilah pemiliknya. **Tidak disunting bersama usulan ini** karena
menyunting angka pada berkas yang menjadi **masukan** usulan, di dalam putaran yang sama, membuat
usulannya tidak dapat diperiksa terhadap apa pun.

---

## Yang ditanyakan sebelum satu berkas tiket pun ditulis

1. **Irisan `14` sebagai PEMBUAT PERTAMA `KONTRAK` dan `VERSI_KONTRAK`, dan kelima tiket lama
   memperoleh penghalang ke sana** — disetujui? Ini mengubah angka papan dari 4 menjadi 2.
2. **`P-01` dan `P-50` dibelah antar batch** — disetujui? Bila tidak, batch 1 menjadi **nol irisan**,
   sebab tidak ada yang membuat `VERSI_KONTRAK`.
3. **Ketiga `G2` — `P-20`, `P-41`, `P-49` — tegas tidak masuk**, meski `DAFTAR-PEKERJAAN.md`
   menandainya `TERTAHAN`. Gerbang `G2` yang menang; dikonfirmasi?
4. **Granularitasnya** — 31 irisan untuk 31 kemampuan. Terlalu halus di bagian daftar anak versi
   (`22`…`30`, sembilan irisan yang bentuknya hampir seragam)? Menggabungkannya melanggar `G1`, tetapi
   `G1` milik pemilik proses untuk diubah.
5. **`39` bertanda MENAHAN atas `F-16`** — atau diturunkan menjadi *menunggui* dengan `PERAN_PELAKU`
   berupa teks bebas sementara?
