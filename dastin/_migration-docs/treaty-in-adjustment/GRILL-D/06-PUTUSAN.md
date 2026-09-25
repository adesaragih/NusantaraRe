> Modul  : Treaty In Adjustment · Ronde D · 2026-09-24
> Sifat  : PUTUSAN. Tidak satu pun ditandai terverifikasi.

# 06 · PUTUSAN RONDE D

| # | Putusan | Keadaan |
|---|---|---|
| **GRL-04** | ~~Penyesuaian bukan entitas~~ | **🔴 BATAL** — dicoret, **tidak dihapus** |
| **GRL-12** | ~~Materialitas diturunkan dari baris selisih tersimpan~~ | **🔴 BATAL** — dicoret, **tidak dihapus** |
| **GRL-19** | **`DOKUMEN_ADDENDUM` adalah entitas tersendiri** | **terkunci** |
| **GRL-20** | **Materialitas adalah MASUKAN, dan ia sakelar dua arah** | **terkunci** |

---

## GRL-04 — **BATAL**

**Barisnya dicoret, bukan dihapus.** Keputusan yang pernah ada adalah bukti premisnya pernah
diperiksa; menghapusnya menghapus bukti itu.

**Pembatal yang tertulis di putusannya sendiri, dikutip apa adanya:**

> *"satu dokumen addendum mengubah **lebih dari satu** kontrak atau versi, **atau** satu revisi
> menggabungkan **lebih dari satu** dokumen addendum."*

| Cabang | Butir | Jawaban | Menyala? |
|---|---|---|---|
| pertama | `DB-3` — *"satu dokumen selalu mengubah tepat satu kontrak"* | **TIDAK** | **ya** |
| kedua | `DB-4` — *"dua addendum tidak pernah digabung jadi satu revisi"* | **TIDAK** | **ya** |

**Keduanya, bukan salah satu.** Pembatal berbentuk *"A **atau** B"* sudah menyala oleh satu cabang;
di sini dua-duanya.

**Yang TIDAK batal bersamanya, dan ini penting:** dasar GRL-04 bukan hanya kardinalitas. Ia juga
bertahan *"sebagai gerbang"* — `METODE` §8.3 menelusuri 19 properti satu per satu dan menemukan
**nol fakta tersisa** untuk `PENYESUAIAN`. **Penelusuran itu tetap berlaku**, dan karena itu
`PENYESUAIAN` **tetap tidak kembali**. Yang lahir adalah benda yang **berbeda**: bukan penyesuaian
di bawah versi, melainkan **dokumen di atas versi**.

---

## GRL-19 — `DOKUMEN_ADDENDUM` adalah entitas tersendiri, di ATAS versi

* **Cabang** — D · **Menggantikan** — GRL-04 · **Mengoreksi** — GRL-01 (bertambah satu entitas)

### Label per bagian

| Bagian | Label | Berubah dari |
|---|---|---|
| dokumen addendum menjadi entitas | **BARU** | sistem lama **tidak punya tempatnya sama sekali** — `TD-01` |
| persetujuan tetap melekat pada versi per kontrak | **PELESTARIAN** | jawaban A: *"1 kontrak di aksep 1 per 1"* |
| `PENYESUAIAN` tidak kembali | **PELESTARIAN** | penelusuran 19 properti `METODE` §8.3 berdiri |

### Keputusan

1. **`DOKUMEN_ADDENDUM` berdiri sendiri**, bukan di bawah `KONTRAK` maupun `VERSI_KONTRAK`.
2. **`VERSI_KONTRAK` memperoleh rujukan ke dokumen, dan rujukan itu BOLEH KOSONG.**
   Satu dokumen → banyak versi, lintas kontrak. Ini yang membuat `DB-3` dan `DB-4` terpenuhi
   **tanpa** mengubah letak persetujuan.
3. **Persetujuan tetap per versi per kontrak.** Tidak ada persetujuan di tingkat dokumen.
4. **Namanya `DOKUMEN_ADDENDUM`** — **bukan** `DOKUMEN_KONTRAK`, yang sudah dipakai modul induk
   untuk lampiran dan berstatus *"dirujuk, tidak dimiliki"* (§10.19). Dua benda berbeda tidak
   berbagi nama.

### Empat soal yang diminta, dijawab dengan ongkos salah dua arah

Ekspor **tidak dapat menjawab satu pun** — `TD-01` nol. Maka keempatnya diputuskan dengan
membandingkan ongkos salahnya, **dan tidak satu pun ditandai terverifikasi**.

#### a. Wajib punya dokumen, atau boleh kosong?

> **Boleh kosong.**

| Arah salah | Ongkos |
|---|---|
| boleh kosong, ternyata selalu ada | satu kolom yang tak pernah kosong — **tidak ada ongkos** |
| wajib, ternyata revisi internal tidak punya dokumen | **migrasi harus mengarang dokumen** untuk setiap revisi internal warisan, dan `EDMState = 1` berbunyi *"Internal Edit"* — revisi yang memang tidak ke luar |

**Syarat pembalikan:** bila bisnis menyatakan setiap versi — termasuk internal — selalu disertai
dokumen bernomor, kolomnya menjadi wajib untuk versi **baru** saja; warisan tetap boleh kosong.

#### b. Boleh menyentuh kontrak dari cedant yang BERBEDA?

> **Ya — dokumen berdiri sendiri, tidak berinduk pada cedant.**

| Arah salah | Ongkos |
|---|---|
| berdiri sendiri, ternyata selalu satu cedant | satu kolom cedant turunan yang selalu terisi sama — **murah** |
| berinduk cedant, ternyata boleh lintas cedant | **perubahan skema + migrasi ulang**, dan dokumen yang sah **ditolak** sampai itu selesai |

**Tidak setangkup.** Dan bentuk berinduk-cedant juga bertabrakan dengan butir 2: rujukan dari versi
ke dokumen sudah membawa kontraknya, sehingga cedant terbaca tanpa disimpan dua kali (ADR-0041).

#### c. Nomor unik GLOBAL atau per cedant?

> **Global**, dengan uji data sebelum migrasi.

Alasannya bukan selera melainkan bentuk: **dokumen berdiri sendiri** (butir b), sehingga tidak ada
induk untuk melingkupi keunikannya. Keunikan per cedant **tidak dapat dikompilasi** pada entitas
yang tidak berinduk cedant.

**Dan ini keputusan yang paling mungkin salah di seluruh GRL-19**, karena arah salahnya **menolak
data yang sah** — kekeliruan yang sudah tiga kali terjadi di korpus ini.

**Maka ia disertai ujinya, bukan keyakinan:** `UA-19` — berapa nomor dokumen addendum yang berulang
lintas cedant di arsip. Bila ada, keunikannya dilingkupi, bukan dicabut.

#### d. Nama

> **`DOKUMEN_ADDENDUM`.** Alasannya di butir 4 keputusan.

### Akibat yang harus diketahui sebelum entitas ini dibangun

> **Seluruh baris warisan akan punya kolom dokumen yang KOSONG, dan migrasi tidak punya sumber untuk
> mengisinya.**

`TD-01` nol berarti nomor itu **tidak pernah ada di sistem**. Pengisiannya **pekerjaan orang** —
mengetik ulang dari arsip kertas — bukan pekerjaan migrasi. Itu bukan cacat rancangan; itu **akibat
yang harus diketahui pemilik proses sebelum ia menyetujui entitas ini**, dan ia naik sebagai butir
sisa di `05-GERBANG.md`.

### Arah dampak bila GRL-19 salah seluruhnya

Bila ternyata dokumen **tidak** perlu menjadi entitas — misalnya nomornya cukup sebagai atribut teks
pada versi — ongkosnya **satu tabel dibuang dan kolomnya dipindahkan**, tanpa kehilangan data.
Sebaliknya, bila ia dibutuhkan dan tidak dibuat, satu dokumen yang menyentuh lima kontrak akan
tersimpan **lima kali** sebagai teks, dan tidak ada yang tahu kelimanya satu benda.

---

## GRL-12 — **BATAL**

**Pembatal yang tertulis di putusannya sendiri:**

> *"bila dibantah dengan **satu pemakaian yang harus diketahui SEBELUM perubahannya dibuat**,
> keputusan ini ditinjau ulang ke arah **(c)**."*

**Jawaban C:** materialitas *"dipilih untuk menentukan mana dan tidak boleh diubah"* — ia dipakai
**sebelum** perubahan dikerjakan, untuk menentukan field mana yang boleh disunting.

Itu bukan mirip dengan pembatalnya. **Itu pembatalnya.** Sesuatu yang dibutuhkan **sebelum**
akibatnya ada **tidak dapat** diturunkan dari akibatnya.

**Dan ia dikuatkan ekspor, bukan hanya keterangan:** `TD-02` — **660 sel** terkunci oleh kondisi
yang menyebut `EDMMaterialType`.

### Premis pendukung GRL-12 juga gugur

GRL-12 bersandar pada **ADR-0037 butir 3** dengan alasan *"materialitas tidak disepakati"*. Jawaban C
menyatakan ia **dipilih di muka** — dan butir yang sama lalu berkata sebaliknya:

> *"bila sebuah besaran boleh disepakati, ia dinaikkan menjadi **masukan**."*

**ADR-0037 tidak dilanggar oleh pembatalan ini; ia justru yang memerintahkannya.**

---

## GRL-20 — Materialitas adalah MASUKAN, dan ia sakelar DUA ARAH

* **Cabang** — D · **Menggantikan** — GRL-12 · **Menyentuh** — GRL-18 bagian 2

### Label per bagian

| Bagian | Label | Berubah dari |
|---|---|---|
| materialitas menjadi **atribut masukan** pada versi | **PELESTARIAN** | sistem lama memang menyimpannya: `TreatyIn.EDMMaterialType` |
| ia **menentukan field mana yang boleh disunting** | **PELESTARIAN** | 660 sel terkunci — mekanismenya nyata dan berfungsi |
| aturannya ditegakkan **saat simpan**, bukan hanya di layar | **PERUBAHAN** | *"materialitas hanya ditegakkan di layar"* — `TDA-10` |
| pembedaannya **dua arah** | **PERUBAHAN** | *"Non Material berarti uang tidak berubah"* — rumusan lama hanya menyebut satu arah |

### Keputusan

1. **`SIFAT_MATERIAL_ADDENDUM` kembali menjadi atribut masukan `VERSI_KONTRAK`**, dengan dua nilai.
2. **Ia dipilih sebelum perubahan dikerjakan**, dan menentukan himpunan field yang boleh disunting.
3. **Aturannya menjadi invarian yang ditegakkan saat simpan** — bukan hanya penguncian layar.
4. **Himpunannya saling lepas**, dan itu terbaca dari ekspor (`TD-02`).

### a. Apa persisnya yang dikunci — dibaca dari kondisinya, bukan dari namanya

| Pilihan | Yang TIDAK boleh berubah |
|---|---|
| **MATERIAL** | `Exclusions`, `SpecialConditions`, `TreatyContractName`, `ContractRefNo` — **teks kesepakatan dan identitas kontrak** |
| **NON MATERIAL** | **setiap besaran uang dan porsi** — `Value`, `Amount`, `Currency`, `RNMShare`, `FacultativeShare`, `FacultativeShareBrokerage`, `Deduction`, `DeductionPct`, keempat batas bahaya beserta mata uangnya, dan tanggal jatuh tempo pelaporan |

### Invarian yang diusulkan, dan ia lebih kuat daripada sistem lama

> **`INV-69`** — pada `VERSI_KONTRAK` bersifat **NON MATERIAL**, tidak boleh ada satu pun baris
> `NILAI_SELISIH` yang besarannya bertipe uang atau porsi.

| Uji | Jawaban |
|---|---|
| **Dapatkah ia gagal?** | **ya** — versi Non Material dengan satu baris selisih premi |
| **Dapatkah ia dikompilasi menjadi nama kolom?** | **ya** — `NILAI_SELISIH` dikelompokkan `ID_VERSI_KONTRAK`, disaring jenis besarannya, dibandingkan dengan `SIFAT_MATERIAL_ADDENDUM` pada versinya |

> **`INV-70`** — pada `VERSI_KONTRAK` bersifat **MATERIAL**, teks kesepakatan dan identitas kontrak
> **tidak berubah** terhadap versi dasarnya.

Keduanya lolos kedua uji. **Yang kedua adalah arah yang sistem lama punya dan tidak pernah
seorang pun tulis** — ia baru terbaca di `TD-02`.

### Uji negatif WAJIB disertai uji positif

Ketiga kekeliruan lingkup di korpus ini **seluruhnya lolos uji negatif**; hanya uji positif yang
menangkapnya.

| # | Yang disisipkan | Harus |
|---|---|---|
| **N-7a** | versi Non Material + satu baris selisih premi | **GAGAL** |
| **N-7a+** | versi Non Material + satu baris selisih **tanggal pelaporan** | **BERHASIL** — bukan besaran uang |
| **N-7b** | versi Material + perubahan `PENGECUALIAN` | **GAGAL** |
| **N-7b+** | versi Material + perubahan limit **dan** pengecualian **tidak** disentuh | **BERHASIL** |
| **N-7c+** | versi **warisan** yang melanggar keduanya | **BERHASIL dimuat** — lihat butir c |

### b. Kapan pilihannya dikunci — **TIDAK DIPUTUSKAN DI SINI**

Keterangan *"tidak boleh diubah"* dapat berarti dua hal yang berbeda akibatnya:

| Bacaan | Artinya |
|---|---|
| (i) | **field** yang tidak boleh diubah — yaitu penguncian yang `TD-02` ukur |
| (ii) | **pilihan materialitasnya sendiri** tidak boleh diubah sesudah ditetapkan |

**Ekspor tidak memisahkan keduanya**, dan saya **tidak memilih sendiri**. `TD-02` membuktikan (i)
ada; ia **tidak** membuktikan (ii) tidak ada.

**Butir daftar-untuk-dibantah baru — dan ia ditulis sebagai SATU fakta, sesuai `MA-12`:**

> **`DB-20`.** *"Sesudah sebuah addendum diajukan, pilihan material / non material-nya tidak dapat
> diubah lagi oleh siapa pun."*

Bila **benar**, materialitas beku **sejak diajukan** — sama dengan jenis di GRL-18 bagian 2.
Bila **dibantah**, kedua sumbu punya **titik beku yang berbeda**, dan itu keputusan tersendiri.

### c. Nilai warisan — baris lama MUNGKIN melanggar

`TDA-10` menyatakan aturan ini **tidak pernah ditegakkan di sisi simpan** di sistem lama; ia hanya
penguncian layar. **Penguncian layar dapat dilewati** — lewat jalur simpan yang tidak melewati
layar itu, lewat tombol paksa, atau lewat perubahan yang dibuat sebelum materialitas dipilih.

> **Maka baris warisan tidak boleh diandaikan patuh**, dan `INV-69`/`INV-70` **tidak dapat langsung
> ditegakkan atas data lama** tanpa mengetahui berapa yang melanggar.

**Ujinya `UA-3`, yang artinya berubah — lihat butir d.** Perlakuan atas yang melanggar
(ditolak, ditandai, atau dinaikkan menjadi Material) adalah **keputusan migrasi**, bukan keputusan
ronde ini, dan ia tidak dapat diambil sebelum angkanya ada.

### d. `UA-3` ditulis ulang kepalanya

| | Sebelum | Sesudah |
|---|---|---|
| **Yang diukur** | ongkos beralih ke bentuk (c) | **berapa baris warisan yang MELANGGAR** `INV-69` dan `INV-70` |
| **Kenapa berubah** | bentuk (c) kini **menjadi keputusan**, bukan alternatif | yang perlu diketahui bukan lagi *"berapa mahal pindah"*, melainkan *"berapa banyak data lama yang tidak muat"* |

**Bunyinya:** *berapa baris `M_TREATY_IN_EDM` bersifat Non Material yang selisih besaran uangnya
bukan nol, dan berapa baris Material yang teks kesepakatannya berbeda dari versi dasarnya.*

### Arah dampak bila GRL-20 salah

Bila materialitas ternyata **memang** dapat diturunkan — misalnya bisnis menyatakan pilihan di muka
itu hanya kenyamanan layar — ongkosnya **satu kolom yang tidak dibaca siapa pun**, dan invariannya
dicabut. Murah.

Sebaliknya, bila ia **tidak** dijadikan masukan dan ternyata dibutuhkan sebelum perubahan ada,
**tidak ada cara menambalnya**: pada saat ia dibutuhkan, akibat yang menjadi sumbernya belum ada.
**Tidak setangkup** — dan itu yang membuat pembatalan GRL-12 tidak dapat ditunda.
