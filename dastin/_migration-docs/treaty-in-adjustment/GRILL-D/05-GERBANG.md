> Modul  : Treaty In Adjustment · Ronde D · 2026-09-24

# 05 · GERBANG — akibat berantai, diperiksa satu per satu

**Untuk setiap keputusan yang TIDAK berubah, dikatakan tidak berubah.** Diam bukan jawaban.

## 1. Keputusan yang diperiksa

| # | Isi ringkas | Tersentuh? | Alasannya |
|---|---|---|---|
| **GRL-01** | addendum adalah `VERSI_KONTRAK` | **DIKOREKSI, bukan dibongkar** | bertambah **satu entitas** di atas versi. Pernyataan *"addendum adalah versi"* **tetap benar** — dokumen bukan addendum, ia **pembungkus** beberapa addendum |
| **GRL-07** | ADR-0052, satu rantai persetujuan | **TIDAK** | jawaban A: *"1 kontrak di aksep 1 per 1"*. Persetujuan tetap per versi per kontrak; rantainya tidak disentuh dokumen |
| **GRL-08** | ADR-0055, daftar keadaan dan perpindahan | **TIDAK** | keadaan melekat pada **versi**. Dokumen tidak punya keadaan, tidak punya perpindahan, dan tidak punya persetujuan |
| **GRL-09** | pengenal tampilan `…/Rnn` | **TERSENTUH — dan pertanyaannya baru** | lihat §2 |
| **GRL-10** | rantai versi, dasar = versi berlaku terakhir | **TIDAK** | rantai dibentuk `ID_VERSI_KONTRAK_DASAR` antar-versi **dalam satu kontrak**. Dokumen menyentuh beberapa kontrak, tetapi **tidak menggabungkan rantainya** |
| **GRL-11** | versi berlaku adalah turunan | **TIDAK** | turunan dihitung per kontrak dari keadaan versinya; dokumen tidak masuk perhitungan |
| **GRL-13** | `JENIS_ADDENDUM` bernilai dua | **TERSENTUH — dan ekspor TIDAK dapat memutuskannya** | lihat §3 |
| **GRL-14** | `ActualValue` dibuang | **TIDAK** | ia tentang pohon nilai, bukan tentang dokumen maupun materialitas |
| **GRL-18** | jenis adalah masukan, beku sejak diajukan | **SEBAGIAN** | lihat §4 |

## 2. `GRL-09` — nomor dokumen adalah PENGENAL atau ATRIBUT?

Jawaban B: *"punya nomor sendiri, tapi tetap ada key merujuk ke ID treaty addendum"*.

| Bacaan | Artinya |
|---|---|
| **pengenal entitas** | nomor dokumen **adalah** kunci alaminya; `…/Rnn` tetap milik versi dan tidak bertabrakan |
| atribut biasa | dokumen punya pengenal sintetis, dan nomornya sekadar teks |

> **Diputuskan: PENGENAL — yaitu kunci alami `DOKUMEN_ADDENDUM`, bukan kunci utamanya.**
>
> Alasannya bentuk: nomor itu **beredar di luar sistem** — ia yang ditulis di kertas dan disebut
> orang. `GRL-09` sudah menetapkan pola yang sama untuk pengenal tampilan versi: **pengenal yang
> dipakai manusia dipelihara sebagai kunci alami, bukan sebagai kunci utama**.
>
> **`…/Rnn` tidak tersentuh.** Ia pengenal **versi**; nomor dokumen pengenal **dokumen**. Keduanya
> hidup di tingkat yang berbeda, dan satu dokumen dapat memayungi beberapa `…/Rnn` dari kontrak
> yang berbeda-beda.

**Arah dampak bila salah:** bila nomor dokumen ternyata berulang dan bukan pengenal, ia turun
menjadi atribut biasa — **satu constraint dicabut**, tidak ada data yang hilang. `UA-19` yang
mengukurnya (`LD-2`).

## 3. `GRL-13` — jenis melekat di DOKUMEN atau di VERSI?

> **Ekspor TIDAK DAPAT memutuskan ini, dan itu dinyatakan alih-alih disamarkan.**
>
> Sistem lama menulis `EDMState` **per baris addendum**, dan ia **tidak punya dokumen sebagai
> benda**. Maka tidak ada satu pun baris di ekspor yang dapat memisahkan *"jenis milik dokumen"*
> dari *"jenis milik versi"* — keduanya menghasilkan data yang sama persis.

**Diputuskan dengan ongkos salah dua arah: jenis tetap melekat pada VERSI.**

| Arah salah | Ongkos |
|---|---|
| melekat pada versi, ternyata milik dokumen | seluruh versi di bawah satu dokumen membawa nilai yang sama — **berulang, tetapi benar**. Menaikkannya kelak: satu kolom dipindahkan |
| melekat pada dokumen, ternyata milik versi | **satu dokumen tidak dapat memayungi versi berjenis berbeda** — dan bila kenyataannya bisa, dokumennya harus dipecah, yaitu **mengubah data, bukan hanya skema** |

**Tidak setangkup.** Dan bacaan "melekat pada versi" adalah satu-satunya yang **tidak kehilangan
kemampuan** apa pun.

**`GRL-13` tidak batal; lingkupnya dipertegas.** Yang ditambahkan: pernyataan terang bahwa
pertanyaan dokumen-versus-versi **pernah diajukan dan diputuskan tanpa bukti ekspor**.

## 4. `GRL-18` — bagian mana yang berdiri dan bagian mana yang jatuh

| Bagian | Keadaan |
|---|---|
| 1 · jenis addendum adalah **masukan** | **BERDIRI** |
| 2 · dapat dibetulkan selama `DRAFT`, **beku sejak diajukan** | **BERDIRI untuk jenis** — tetapi **tidak otomatis berlaku untuk materialitas**; itu `DB-20` (`LD-3`) |
| 3 · tiap perubahan berjejak | **BERDIRI** |
| — · `EDMMaterialType` dihapus **tanpa pengganti** | **JATUH.** Ia bersandar pada `GRL-12`. `GRL-20` mengembalikannya sebagai atribut masukan |

> **`GRL-18` tidak batal.** Tiga bagiannya berdiri; yang jatuh adalah **akibat** dari `GRL-12`, bukan
> penalaran `GRL-18` sendiri.

## 5. `DB-16` dirumuskan ulang — sekarang ada bukti dokumen benda tersendiri

| | |
|---|---|
| **Bunyi lama** | *"Addendum eksternal selalu berupa dokumen yang disepakati dan ditandatangani cedant, dengan tanggal berlaku sendiri; revisi internal tidak pernah dikirim ke luar."* |
| **Kenapa dirumuskan ulang** | ia **empat klaim dalam satu kalimat** — persis kesalahan `MA-12`. Dan sebagiannya **sudah terjawab**: `DB-3`/`DB-4` membuktikan dokumen adalah benda tersendiri |

**Penggantinya, satu fakta satu butir:**

> **`DB-16a`.** *"Revisi internal tidak pernah disertai dokumen yang dikirim ke cedant."*
> **`DB-16b`.** *"Dokumen addendum selalu punya tanggal berlaku sendiri, yang dapat berbeda dari
> tanggal berlaku versi yang dipayunginya."*

`DB-16a` menjawab butir a `GRL-19` (boleh kosong). `DB-16b` menentukan apakah tanggal berlaku milik
dokumen atau milik versi — dan itu **belum diputuskan di mana pun**.

## 6. `CABANG-K-PEMETAAN-TO-SPEC.md`

**Ditandai perlu memuat `GRL-12` sampai `GRL-20`.** **Tetap pekerjaan sesi to-spec, bukan ronde
ini** — dan ia tidak dapat dikerjakan sebelum `DB-20`, `DB-16a`, dan `DB-16b` terjawab, sebab
ketiganya menentukan letak kolom.

---

## 7. Sisa ronde D — dipisah menurut pembacanya

### 7.1 Yang MENAHAN MODEL — tanpa jawabannya, bentuk datanya tidak dapat ditentukan

| # | Sisa | Siapa menjawab | Apa yang berubah bila jawabannya datang |
|---|---|---|---|
| **SD-1** | **`DB-20`** — apakah pilihan materialitas beku sesudah diajukan | **bisnis** | menentukan apakah materialitas dan jenis punya **titik beku yang sama**; bila berbeda, invariannya bercabang |
| **SD-2** | **`DB-16a`** — apakah revisi internal pernah disertai dokumen | **bisnis** | menentukan apakah rujukan dokumen **boleh kosong** (butir a `GRL-19`) |
| **SD-3** | **`DB-16b`** — apakah dokumen punya tanggal berlaku sendiri | **bisnis** | menentukan apakah **tanggal berlaku** milik dokumen atau milik versi — satu kolom berpindah tingkat |
| **SD-4** | paket **`REV-1` … `REV-6`** belum ditanggapi, dan kini **`REV-4` bertambah alasannya** | pemilik ADR induk | ADR-0049 dan ADR-0037 keduanya tersentuh `GRL-20` |

> **Keempatnya menentukan LETAK KOLOM**, bukan sekadar isinya. Karena itu `CABANG-K` tidak dapat
> diperbarui sebelum ketiga `DB` terjawab.

### 7.2 Yang hanya MENGUKUR KERUSAKAN — tidak menahan model

| # | Uji | Yang diukurnya |
|---|---|---|
| **MD-1** | **`UA-3`**, kepalanya ditulis ulang | berapa baris warisan **melanggar** `INV-69`/`INV-70` — Non Material berselisih uang, atau Material berubah teksnya |
| **MD-2** | **`UA-19`** (baru) | berapa nomor dokumen addendum **berulang lintas cedant** di arsip — menentukan apakah keunikan global (`GRL-19` butir c) menolak data yang sah |
| **MD-3** | ongkos pengisian ulang nomor dokumen warisan dari arsip kertas | **berapa banyak** versi warisan yang perlu diisi, dan apakah sepadan (`LD-1`) |
| **MD-4** | rekonsiliasi angka `GERBANG-0-TO-SPEC.md` | tidak mengubah kesimpulan apa pun (`LD-6`, `MA-13`) |

> **MD-1 tidak menahan model, tetapi ia menahan MIGRASI.** Invariannya sudah dapat ditulis; yang
> belum dapat diputuskan adalah **perlakuan atas yang melanggar**, dan itu keputusan pemilik proses
> sesudah angkanya ada.
