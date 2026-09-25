# `TDA-17` — tiga kolom tanpa rumah: putusan

**Tanggal:** 24 September 2026 · **Putaran penutup to-spec**
**Menutup:** §11.1 `SISA-DAN-SERAH-TERIMA.md`
**Sifat:** putusan to-spec. Tidak ada tiket, struct Golang, komponen React, maupun endpoint.

> ## ⚠ SEBAGIAN DIBANTAH 24 September 2026 — baca `F-15` LEBIH DULU
>
> `2-to-spec/F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` menemukan **rumus yang menghitung** dua dari
> ketiga kolom, di dalam aturan yang **hidup**, dan dua artefak induk yang **sudah** mengadili
> `PremiumEarnedList` sebagai **DITURUNKAN** dengan dasar `EVIDENCED`.
>
> | Kolom | Keadaan sesudah `F-15` |
> |---|---|
> | `DEDUCTIBLE_KEDUA` | **bertahan** — sapuan penulisnya kembali **bersih**; ia yang mengkalibrasi temuan itu |
> | `PERSEN_ROL` | **dipertikaikan** — `DetailCalculationROL` menghitungnya `premi ÷ limit × 100`. Kolomnya berdiri, **ditandai, tidak dicabut diam-diam** |
> | `PREMIUM_EARNED` | **TIDAK dibuat** — `DetailCalculation` langkah 25 menghitungnya `EGNPI × tarif penyesuaian` |
>
> §2 dan §3.3 di bawah **dipertahankan apa adanya sebagai penalaran putaran sebelumnya**, bukan
> sebagai keadaan sekarang. Yang berubah bukan bentuknya melainkan **nasibnya**.

> ## RINGKAS — putaran 24 September 2026 pagi
>
> **Ketiganya DISIMPAN.** Bukan dibuang.
>
> Dan pembingkaian yang dibawa ke putaran ini — *"muncul hanya di pohon cermin"* dan *"nol pembaca,
> jadi ia ditulis tetapi tidak pernah dibaca"* — **keduanya tidak bertahan sesudah diperiksa.**
> Rinciannya §1 dan §2.

---

## 1. Koreksi pertama — ketiganya ADA di pohon utama

§11.1 menyatakan: *"Mereka muncul di `PETA-TELUSUR-JSON.md` **hanya di dalam pohon cermin**, tidak
pernah di pohon utama."*

Disapu atas `4-erd-dan-tabel-datar/datar-treatyin-lama.csv` — sumber yang **mandiri** dari peta
telusur, dibangkitkan `buat-pohon-treatyin.py` dari 329 berkas ekspor:

| Ruas | Pohon **utama** | Pohon cermin |
|---|---:|---:|
| `Deductible2` | **3** | 6 |
| `PremiumEarned` | **4** | 8 |
| `ROL` / `RolPct` | **5** | 7 |

Jalur pohon utamanya, apa adanya:

```
TreatyIn.Limits.Deductible2                 Skalar
TreatyIn.Limits.ROLPct                      Skalar
TreatyIn.Limits.PremiumEarnedList           Page List
TreatyIn.Limits.PremiumEarnedList.Currency  Skalar
TreatyIn.Limits.PremiumEarnedList.Value     Skalar
—— dan dua agregat, yang bukan atribut ——
TreatyIn.TotalLimitsPremiumEarned           Skalar
TreatyIn.TotalLimitsROL                     Skalar
```

**`TreatyIn.Limits[]` adalah `LAYER`.** Keempat saudaranya yang sudah punya rumah — `DEDUCTIBLE`,
`MDP`, `MDP_PCT`, `ADJ_RATE` — berasal dari jalur yang sama persis.

### Kenapa §11.1 dapat benar dan tetap menyesatkan

Keduanya menyebut "pohon utama" untuk **semesta yang berbeda**:

| Sumber | Semestanya |
|---|---|
| `PETA-TELUSUR-JSON.md` | **667 jalur teradjudikasi**, dan **272 dikecualikan sebagai cermin**. Berlabel *"semestanya kurang"* sejak awal |
| `datar-treatyin-lama.csv` | **985 simpul** dari sapuan ekspor, tanpa adjudikasi |

**"Tidak ada di pohon utama peta telusur" bukan "tidak ada di pohon utama."** Yang pertama
pernyataan tentang sebuah berkas yang sudah mengaku bolong; yang kedua pernyataan tentang sistem
lama. §11.1 menulis yang pertama dan membacanya sebagai yang kedua.

Ini **instans `L-8` yang ketiga** — dan ia menguatkan, bukan melemahkan, kalimat §11.1 yang paling
penting: *"§10 tidak melewatkannya — §10 tidak pernah diperlihatkan kepadanya."*

---

## 2. Koreksi kedua — "nol pembaca" benar, dan ia menjawab pertanyaan yang lain

Bukti yang dibawa: ketujuh kolom hanya muncul di `SaveTreatyInDetail.xml`, sebuah **penulis**;
kalibrasi lulus dengan `TREATYID`, 21 rujukan. **Itu tidak dibantah.**

Yang dibantah adalah kesimpulan yang ditarik darinya:

> Sapuan itu menyapu **kolom tabel datar** `TREATYINDETAIL`. Ia menyatakan **kolom datarnya** tidak
> pernah dibaca. Ia **tidak** menyatakan **besarannya** tidak dipakai.

Dan besarannya **dipakai**. `PENGETAHUAN.md` modul Adjustment §5.2, tentang bentuk rumus mesin
selisih:

> *"Satu-satunya pengecualian bentuk adalah **`Deductible2`**, yang dibungkus `@toDecimal(...)` **di
> kedua sisi** — bukti bahwa sebagian nilai memang tersimpan sebagai teks."*

Rumus selisih berbentuk `baru − OLDDATA`. Sebuah ruas yang **dikurangkan** pada kedua sisinya
**ada di pohon utama dan dibaca** — kalau tidak, tidak ada yang dapat dikurangkan.

> **Tabel datar adalah proyeksi yang diturunkan, bukan sumber.** Membuang sebuah besaran karena
> proyeksi datarnya tidak dibaca berarti menurunkan model dari ringkasan — persis yang aturan modul
> ini larang.

**Batas sapuan "nol pembaca", dan ia wajib ikut tertulis:**

| Batas | Akibatnya |
|---|---|
| disapu atas **aturan yang ter-ekspor saja** | `L-10` — **21 jenis aturan nol kemunculan** di ekspor. "Nol pembaca" berarti nol di antara jenis yang ada |
| `Declare Expression` **membaca tanpa dipanggil** | ia salah satu dari 21 jenis itu. Sebuah ekspresi deklaratif yang membaca `ROL_PCT` **tidak akan terlihat oleh sapuan pemanggil mana pun** |

---

## 3. Putusan per kolom

### 3.1 `DEDUCTIBLE2` → **DISIMPAN**, atribut `LAYER`

| | |
|---|---|
| **Asal** | `Limits[].Deductible2` |
| **Nama baru** | **`DEDUCTIBLE_KEDUA`** — 16 bita, mengikuti saudaranya `DEDUCTIBLE` |
| **Tipe** | U1 paket uang — sejajar `DEDUCTIBLE` |
| **Boleh kosong** | ya |
| **Label** | **PELESTARIAN** — besaran yang sudah ada, diberi rumah yang sudah seharusnya |
| **Arah dampak bila salah** | bila ia ternyata bukan deductible kedua melainkan besaran lain, satu kolom berganti nama. Murah |
| **Syarat pembalikan** | bila sapuan penulisnya menunjukkan ia hanya perancah layar tanpa arti bisnis, kolomnya dicabut **sebelum data dimuat** |

> **Satu hal yang ikut dibawa dan tidak boleh hilang:** ia **tersimpan sebagai teks** di sistem lama
> — itu sebab `@toDecimal` dipasang di kedua sisi rumusnya. Migrasinya **harus mengubah tipe**, dan
> baris yang tidak dapat dikonversi **tidak boleh disamarkan menjadi nol** (ADR-0035).

### 3.2 `ROL_PCT` → **DISIMPAN**, atribut `LAYER`

| | |
|---|---|
| **Asal** | `Limits[].ROLPct` |
| **Nama baru** | **`PERSEN_ROL`** — 10 bita, mengikuti `PERSEN_PENYESUAIAN` dan `PERSEN_MINIMUM_DEPOSIT` |
| **Tipe** | P2 persentase |
| **Boleh kosong** | ya |
| **Label** | **PELESTARIAN** |
| **Arah dampak bila salah** | ROL adalah **turunan** — premi dibagi limit. Bila bisnis menyatakan ia selalu dihitung dan tidak pernah disepakati, ia turun menjadi turunan dan kolomnya dicabut |
| **Syarat pembalikan** | jawaban bisnis atas *"apakah rate on line pernah disepakati sebagai angka, atau selalu dihitung"* — **butir baru untuk `PERMINTAAN-TEKNIK-TREATY.md`** |

> `TotalLimitsROL` **tidak** ikut disimpan — ia agregat, turunan (ADR-0037).

### 3.3 `PREMIUM_EARNED` → **DISIMPAN**, tetapi **bukan kolom `LAYER`**

Ini yang berbeda dari dua di atas, dan bedanya bukan detail:

```
TreatyIn.Limits.PremiumEarnedList            Page List
TreatyIn.Limits.PremiumEarnedList.Currency   Skalar
TreatyIn.Limits.PremiumEarnedList.Value      Skalar
```

**Ia daftar bernilai per mata uang** — bentuk yang persis sama dengan lima paket uang yang `P-8`
sudah golongkan **golongan A** dan jadikan **tabel anak per mata uang**: `NILAI_MDP`,
`NILAI_MDP_MINIMUM`, `NILAI_CADANGAN_PREMI`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM`.

| | |
|---|---|
| **Usulan bentuk** | tabel anak **`NILAI_PREMI_DIPEROLEH`** — 21 bita — berinduk `LAYER`, kunci alami kode mata uang, sejajar kelima saudaranya |
| **Label** | **PELESTARIAN** bentuknya; `P-8` sudah memutuskan polanya, ini memperluasnya |
| **Arah dampak bila salah** | bila ia ternyata satu nilai bermata uang tunggal, ia menyusut menjadi kolom `LAYER` + `KODE_MATA_UANG`. Menyusutkan tabel anak menjadi kolom **lebih murah** daripada memekarkan kolom menjadi tabel sesudah data dimuat |

> ### DAN INI MENGUBAH CACAH ENTITAS — 34 menjadi 35
>
> Itu **bukan keputusan to-spec biasa**. Setiap penambahan entitas sebelumnya — `PEMULIHAN_LIMIT`,
> `PERISTIWA_KONTRAK`, kelima tabel anak `P-8` — diputuskan **pemilik proses**, dan §10.23c
> menetapkan presedennya: tiga dari empat usul entitas saya **gugur** setelah diperiksa.
>
> **Maka bentuknya saya usulkan, tidak saya terapkan.** Yang saya terapkan hanya §3.1 dan §3.2 —
> dua kolom `LAYER`, yang tidak mengubah cacah entitas.
>
> | | |
> |---|---|
> | **Siapa memutuskan** | pemilik proses |
> | **Yang menagih** | berkas ini, dan ketiadaan `PREMIUM_EARNED` di `KAMUS-KOLOM.md` |
> | **Bila ditunda** | `LAYER` berdiri tanpa premi diperoleh, dan selisih atasnya **tidak dapat dihitung** — sama persis dengan keadaan sistem lama, jadi penundaannya **tidak memburukkan apa pun** |

---

## 4. Penagih §11.1 — **belum dicabut, dan sebabnya**

Kewajiban 3 putaran ini berbunyi: *"Penagih yang sudah dijawab dicabut."*

**Penagih §11.1 belum dapat dicabut seluruhnya**, karena ia menagih **tiga** kolom dan baru **dua**
yang mendarat:

| Kolom | Keadaan | Penagih |
|---|---|---|
| `DEDUCTIBLE2` | diputuskan, masuk §10.3 | **dicabut** sesudah `KAMUS-KOLOM.md` dibangkitkan ulang |
| `ROL_PCT` | ~~idem~~ **dipertikaikan `F-15`** — berdiri di §10.3 dengan penanda | **berbunyi kembali**, kini atas **nasibnya**, bukan atas rumahnya |
| `PREMIUM_EARNED` | ~~diusulkan~~ **TIDAK dibuat** — `F-15` | **tetap berbunyi**, dan bunyinya **berganti**: bukan lagi *"pemilik proses memutuskan bentuknya"* melainkan *"pemilik proses memutuskan **nasibnya**"* |

**Penagih yang dipersempit tetap penagih.** Yang dilarang adalah membiarkannya berbunyi untuk hal
yang sudah dijawab — bukan mempertahankannya untuk yang belum.

---

## 5. Yang berubah pada `L-8`

§11.1 menyatakan sebab ketiga kolom ini hilang adalah `L-8`. **Itu bertahan**, dan putaran ini
**menaikkan kadarnya**:

> `L-8` bukan lagi risiko bahwa titik buta memakan sesuatu. **Ia terbukti memakan tiga besaran
> `LAYER`** — dua di antaranya kini dipulihkan, satu menunggu keputusan.

Dan angkanya belum selesai: **340 properti titik buta belum diperiksa siapa pun** (`M-4`). Tiga yang
ditemukan berasal dari 74 yang **sudah** diperiksa. Perbandingan itu ditulis di sini supaya tidak ada
yang membaca pemulihan ini sebagai penutupan `L-8`.

| | |
|---|---|
| **Siapa menutup `L-8`** | gelombang 2 — 340 properti sisa |
| **Yang menagih** | `M-4` di `SISA-DAN-SERAH-TERIMA.md` §5.2 |
