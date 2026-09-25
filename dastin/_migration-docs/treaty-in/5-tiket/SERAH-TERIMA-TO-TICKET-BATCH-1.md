# Serah terima — ronde to-ticket batch 1 Treaty In

**Tanggal:** 25 September 2026 · **Ronde:** tiket `14`…`44`
**Sifat:** serah terima. **Bukan** pernyataan bahwa batch 1 selesai — ia pernyataan tentang **apa
yang ada, apa yang menahan, dan apa yang keliru.**

---

## 1. Gerbang selesai — dijawab satu per satu

| Gerbang | Keadaan |
|---|---|
| adjudikasi 51 kemampuan terhadap papan yang sudah ada | **YA** — tiga golongan berhitungan; enam koreksi terhadap sapuan kata, masing-masing bersebab |
| pemisahan batch 1 / batch 2 beralasan | **YA** — tetapi **sebabnya keliru**, dan sudah diralat 25 Sep. Lihat §4 `K-1` |
| dua kemampuan terbelah antar batch, garis pecahnya tertulis | **YA** — `P-01` dan `P-50`, garisnya di `USULAN-IRISAN-TREATY-IN-BATCH-1.md` §0.2 |
| 31 tiket ditulis, sembilan medan tanpa kecuali | **YA** — 31 ditulis, 0 ditolak perkakas |
| penahan luar papan bernama asli | **YA** — `KTV-A`, `T-2`…`T-6`, `F-15`, `F-16`, `L-3`, `Uji AD` |
| "menahan" dipisah dari "menunggui" | **YA** — lima menahan, delapan baris menunggui |
| **"yang sebenarnya dapat dimulai" dihitung** | **YA, dan angkanya pernah SALAH.** Lihat §4 `K-2` — ini kesalahan terpenting ronde ini |
| papan diperbarui | **YA** — 44 tiket, pencacah berlabel bendanya |

---

## 2. Yang dihasilkan

| Keluaran | Isi |
|---|---|
| `issues/14-…` … `issues/44-…` | **31 berkas tiket**, sembilan medan, urut ketergantungan |
| `USULAN-IRISAN-TREATY-IN-BATCH-1.md` | adjudikasi, pemisahan batch, 31 irisan beserta sisi penghalangnya |
| `issues/README.md` | papan 44 tiket, dua batch, "dapat dimulai", "menunggui" dipisah |
| `ASUMSI-CLEAR.md` §1a | **sembilan kode asumsi baru**, kolom tiketnya **diisi dari sapuan** |
| `USULAN-ATURAN-TA-12.md` | aturan baru yang lahir dari ronde ini |

**Angka yang dapat diperiksa ulang**, bukan diingat:

| | |
|---|---:|
| tiket di papan | **44** |
| kemampuan di `DAFTAR-PEKERJAAN.md` | 65 didaftar · 3 dicoret · **62 aktif** |
| kemampuan **belum disebut tiket mana pun** | **22** |
| kode `DIASUMSIKAN-CLEAR` yang hidup | **13** |

> Keempat angka di atas **dari sapuan mekanis atas medan `*Asal:*` dan `DIASUMSIKAN-CLEAR(`**, bukan
> dari hitungan tangan. Perintahnya dapat dijalankan ulang kapan saja.

---

## 3. Yang MENAHAN, dipisah dari yang MENGUKUR KERUSAKAN

### 3.1 Yang menahan — tidak boleh jalan sebelum dijawab

| # | Penahan | Tiket | Siapa mencabutnya |
|---|---|---|---|
| 1 | `DB-16a` · `DB-16b` · `DB-20` | `04`, `07`, `08` | **bisnis** — dan **ketiganya belum dikirim** |
| 2 | `Uji AD` | `38` | pemilik data; menentukan **golongan** tiketnya |
| 3 | `L-3` — tidak ada instans Oracle terjangkau | `38`, `43` | **kantor** — lubang lingkungan, bukan lubang spesifikasi |
| 4 | `KTV-A` — presisi kolom | **19 tiket**, dan **bertenggat** | gerbang sesi DDL, **sebelum `44` memuat data pertama** |

### 3.2 Yang hanya mengukur kerusakan — tidak menahan satu tiket pun

| Uji | Yang diukurnya |
|---|---|
| `UA-3` | berapa baris warisan melanggar `INV-69`/`INV-70` |
| `UA-18` | berapa daftar memuat dua baris berkunci padanan sama |
| `UA-19` | sebaran waktu `ProRatePercent` bukan 100 |
| `UA-21` | berapa kontrak bernomor `/Rnn` ganda, berapa addendum tertimpa |
| `Uji AN` | baris ganda `TREATYINDETAIL` — akibat prosedur tanpa cabang `ELSE` |

**Pemisahan ini yang menentukan jadwal.** Golongan pertama menahan pekerjaan; golongan kedua
menahan **angka kerugian historis**, dan itu pekerjaan yang berbeda dengan pembaca yang berbeda.

---

## 4. Kesalahan yang dicatat

### `K-1` — sebab pemisahan batch keliru, dan kekeliruannya menyebar ke tiga berkas

Ronde ini menulis: *"batch 2 ditunda, menunggu `REV-3`, sebab `REV-3` dapat mengubah kolom
keadaan."*

**`REV-3` menyatakan sebaliknya, di barisnya sendiri:** keputusan ADR-0055 — *"tidak ada satu pun
jalan menyetel keadaan selain melalui perpindahan di daftar"* — **tetap, seluruhnya**. Yang
diusulkan berubah hanya **tabel §4**, dan keempat butirnya **koreksi atas pemerian sistem lama**.
Dan `KTV-4` sudah mengadilinya lebih dulu: ***"nol dari enam menentukan letak kolom."***

| Salah | Benar |
|---|---|
| *"batch 2 ditunda, menunggu `REV-3`"* | *"batch 2 dikerjakan sesi berikutnya, sebab 51 kemampuan tidak muat satu sesi"* |

**Pemisahannya tidak dicabut** — ia benar sebagai ukuran sesi. Yang salah hanya sebabnya. Diralat
25 September di `issues/README.md`, `ASUMSI-CLEAR.md` baris 3, dan berkas ini.

> **Kenapa ini bukan soal kata-kata:** yang pertama membuat orang **menunggu jawaban yang tidak akan
> mengubah apa pun**; yang kedua membuat mereka **menjadwalkan sesi**.

### `K-2` — angka *"dapat dimulai"* dihitung atas tepi yang salah

Papan pernah mencetak **"dapat dimulai: 4 — `01`, `02`, `03`, `05`"**. **Keempatnya tidak dapat
dimulai.** Ketiganya menambahkan **kolom** pada `VERSI_KONTRAK` dan satu menambahkan tabel yang
merujuknya — dan **tidak ada satu pun tiket di papan waktu itu yang membuat tabelnya**.

> Angka itu dihitung atas **penghalang di dalam papan**, sementara penghalangnya ada **di luar** —
> di lingkup yang ronde itu belum digarap.

**Dan ia akan terulang.** Selama sebuah batch hanya memuat sebagian kemampuan sebuah modul,
*"tidak ada penghalang di papan"* **tidak berarti** *"tidak ada penghalang"*. Papan yang tidak
lengkap **selalu** melaporkan angka yang terlalu besar, dan ia melaporkannya **dengan tenang**.

**Cara memeriksanya, sebelum mencetak angka itu di batch mana pun:**

> Untuk tiap tiket yang mengaku dapat dimulai, sebutkan **entitas yang disentuhnya** dan tunjukkan
> **tiket mana yang membuatnya**. Bila jawabannya *"tidak ada tiket, tabelnya sudah ada"*, itu
> **andaian** — dan andaian itulah yang keliru.

### `K-3` — dua kemampuan jatuh di antara dua batch

Sapuan mekanis 25 September menemukan **22** kemampuan belum ditiketkan; `§0.2` memperkirakan
**20**. Ketiga selisihnya punya sebab berbeda, dan ketiganya layak dicatat:

| Kemampuan | Apa yang terjadi |
|---|---|
| `P-38` | disebut *"sudah tersentuh papan"*. **Tersentuh, bukan tercakup** — tiket `14` menyebut `INV-24`, tetapi medan `*Asal:*`-nya tidak menyebut `P-38`, dan isinya tidak menegakkannya |
| `P-20` · `P-49` | **tidak masuk batch mana pun.** Keduanya `TERTAHAN` jawaban luar (`Uji X-2`, pemilik tabel acuan), dan keduanya bukan kemampuan daur hidup — sehingga tidak tersaring ke batch 2 maupun tersisa di batch 1 |

**Kemampuan yang tidak masuk batch mana pun tidak akan pernah ditiketkan**, dan tidak ada pencacah
yang memperlihatkannya. Ketiganya dimasukkan ke batch 2.

### `K-4` — dua penanda asumsi tidak dapat disapu

`ASUMSI-CLEAR.md` mencatat tiket `38` bergantung pada **`L-3`** dan **`Uji AD`**. Tiket `38` memang
menyebut keduanya — di medan `PENGHALANG` dan `Blocked by`, **bukan sebagai `DIASUMSIKAN-CLEAR(...)`**.

> Sebuah asumsi yang tidak membawa penandanya **tidak akan muncul di sapuan**, dan **sapuan satu
> perintah adalah seluruh alasan penanda itu ada.**

Diperbaiki 25 September: tiket `38` kini membawa `DIASUMSIKAN-CLEAR(Uji AD)` dan
`DIASUMSIKAN-CLEAR(L-3)` di samping `KTV-A`.

### `K-5` — batas perkakas pemeriksa saya sendiri, dicatat supaya tidak jadi temuan palsu

Pemeriksa dua arah yang saya jalankan melaporkan `KTV-A` **tidak cocok** — sapuan menemukan 19
tiket, kolomnya mencatat 11. **Kolomnya benar**; ia menulis rentang `22…31`, dan pemeriksanya tidak
dapat memekarkan rentang.

> Perkakas pemeriksa yang tidak mengenali notasi ringkas akan **melaporkan berkas yang benar sebagai
> salah** — dan temuan palsu memakan waktu yang sama dengan temuan sungguhan.

---

## 5. Yang ditagih SEBELUM implementasi

| # | Tagihan | Kenapa ia tidak dapat menunggu |
|---|---|---|
| 1 | **`14` dan `15` dikerjakan lebih dulu** | keduanya **PEMBUAT PERTAMA**; `14` membuat `KONTRAK` dan `VERSI_KONTRAK` yang **seluruh papan** menggantung padanya, `15` membuat keenam tabel acuan yang **setiap kunci asing** menunggunya. Begitu keduanya mendarat, yang dapat dimulai melompat ke **tiga belas** |
| 2 | **`KTV-A` dipersempit sebelum `44` berjalan** | `44` **memuat data pertama**, dan sesudah data masuk presisi **tidak dapat dipersempit lagi**. Ini satu-satunya tagihan bertenggat di ronde ini |
| 3 | **Tiga pertanyaan bisnis dikirim** — `DB-16a`, `DB-16b`, `DB-20` | ketiganya menahan tiket, dan **belum satu pun dikirim**. Mengirim paket bertiga sekaligus lebih murah daripada tiga putaran |
| 4 | **Taksiran `14` dan `15` tidak dibandingkan dengan tiket lain** | aturan `PEMBUAT PERTAMA`: irisan tegak membuat tiket pertama sebuah entitas 3–5× lebih besar. Membandingkannya dengan tiket biasa membuat seluruh perkiraan melenceng |

---

## 6. Apa yang ronde ini TIDAK kerjakan

| Hal | Kenapa |
|---|---|
| **22 kemampuan batch 2** | ukuran sesi — 51 kemampuan tidak muat satu jendela konteks. **Bukan** karena `REV-3` (`K-1`) |
| berkas tiket batch 2 | menunggu daftar irisannya disetujui |
| implementasi apa pun | hanya atas perintah pemilik proses |
| penaksiran | belum diminta, dan `PEMBUAT PERTAMA` harus ditandai lebih dulu agar taksirannya tidak dibandingkan salah |
