> Modul  : Treaty In Adjustment · ronde C
> Dibuat : 2026-09-24
> Sifat  : lubang **dilaporkan, tidak ditambal** (`METODE` §6.2). Setiap butir membawa **pemilik**
>          dan **saat penagihan**; tidak satu pun memblokir rancangan (`METODE` §7.2).
> Status : terbuka

# Lubang ronde C

## 1. Kueri data baru — seri `UA`

Nomor melanjutkan `UA-18`. **Tidak satu pun memblokir rancangan.** Dua di antaranya menyatakan batas
pembuktiannya sendiri, dan itu bagian dari bunyinya, bukan catatan kaki.

### `UA-19` — sebaran waktu `ProRatePercent ≠ 100`

**Bentuknya KAPAN, bukan APAKAH.** Untuk baris yang memuat nilai selain 100, ambil sebaran waktunya.

| Cabang hasil | Bacaan |
|---|---|
| berhenti pada suatu tanggal | kemampuannya pernah hidup lalu dimatikan; tanggal berhentinya dicocokkan dengan jejak perubahan aturan |
| masih berjalan sampai sekarang | ada **penulis aktif** di luar ekspor yang kita pegang — jauh lebih kuat daripada sekadar keberadaan |
| tidak pernah ada | bacaan NC-03 menguat |

**Batas pembuktian:** ia dapat **mematahkan** NC-03, tidak dapat **mengesahkannya**. Sebuah nilai
dapat masuk ke data lewat jalan yang bukan aturan — perbaikan manual, prosedur, pemuatan awal, atau
aturan yang sejak itu dihapus.

**Menagih:** `GRL-15`. **Pemilik:** DBA, sesudah izin kueri baca-saja turun.

### `UA-20` — nilai `ProRatePercent` pada baris ber-`IsProRate` benar

Cacah per nilai: berapa **0**, berapa **100**, berapa lainnya.

Sebabnya satu ujung yang **tidak terbaca dari ekspor**: `TreatyCalculateProratePct` langkah 3
menyetel `ProRatePercent = 0`, dan langkah 4 hanya berjalan bila `IsProRate` benar. Penghitungnya
dirujuk **hanya dari aturan layar** (NC-03), sedangkan mesin selisih memanggil pengalinya. Maka bila
`IsProRate` dicentang tanpa peristiwa layar yang menjalankan penghitungnya, nilai yang dikalikan
mesin selisih **tidak diketahui** — dan 0 di antara kemungkinannya.

> Ini **bukan** kesimpulan. Ekspor tidak dapat mengatakan nilai apa yang tersisa di halaman clipboard
> pada saat pengajuan. Yang ekspor katakan hanya: **tidak ada aturan non-layar yang menghitungnya.**

**Menagih:** `GRL-15`, dan bila hasilnya memuat 0 dalam jumlah berarti, ia naik menjadi **temuan
uang** — pengali nol memusnahkan seluruh selisih yang diproratakan. **Pemilik:** DBA.

### `UA-21` — luas kerusakan penomoran `/Rnn`

Dua angka, dipisah: berapa **kontrak** memuat nomor `/Rnn` ganda; dan berapa **addendum** yang
tertimpa (baris yang `EDMDATE`-nya berubah tanpa pertambahan nomor).

**Hasilnya tidak mengubah `GRL-17`** — bila nol, aturannya tetap ditulis, karena ia menutup kelas
persoalan dan bukan satu nilai. Yang berubah hanya perkiraan berapa banyak pengecualian migrasi yang
menunggu. **Menagih:** `GRL-17`, `TDA-01`, `TDA-12`. **Pemilik:** DBA.

---

## 2. Pertanyaan ke orang — seri `DB`

Nomor melanjutkan `DB-17`. Bahasa bisnis, **tanpa nama tabel maupun nama aturan** (`METODE` §4.7).
Bentuknya **pernyataan untuk dibantah**, bukan pertanyaan terbuka.

### `DB-18` — share fakultatif

> *"Ketika sebuah addendum mengubah bagian fakultatif, perubahan itu memang **tidak** dihitung
> sebagai selisih — dan itu disengaja, bukan kelalaian."*

**Bila dibantah** (yaitu bila bisnis memang mengharapkannya dihitung): tidak ada yang berubah,
`GRL-16` sudah membawanya. **Bila dibenarkan**: himpunan besaran yang dipadankan menyusut satu, dan
`GRL-16` menerima satu penyaring — bukan pembongkaran.

**Menagih:** `GRL-16`. **Pemilik:** pengisi kontrak dan bagian yang membukukan retro.

### `DB-19` — penyesuaian premi yang tidak mengubah angka

> *"Penyesuaian premi yang ternyata tidak mengubah satu angka pun **tidak pernah dibuat**; dan
> seandainya dibuat, ia tidak perlu diperlakukan sebagai perubahan yang material."*

**Bila dibantah**, yang berubah **bukan** `GRL-14` melainkan **`GRL-12`**, ke arah pilihan (c) yang
sudah disediakan di sana — materialitas sebagai pernyataan, bukan turunan.

**Menagih:** batas kasus `GRL-14`. **Pemilik:** pengisi kontrak dan penyetuju.

> Kedua butir di atas bergabung dengan `DB-1`…`DB-17` yang sudah siap kirim di
> `GRILL-A/07-AUDIT.md` §3. **Dikirim sebagai satu paket**, bukan satu per satu.

---

## 3. Usulan revisi ADR — `REV-6`

`REV-1`…`REV-5` sudah ada di [`../USULAN-REVISI-ADR.md`](../USULAN-REVISI-ADR.md). Artefak induk
**tidak disunting**; ini draf yang diserahkan sekaligus di akhir grilling.

### `REV-6` — ADR-0037, faktor prorata

**Yang diusulkan:** ADR-0037 menetapkan *"faktor prorata addendum = turunan dari tanggal, bukan
centang"*. Keputusannya **tidak diusulkan berubah**. Yang diusulkan ditambahkan dua hal:

1. **Labelnya PERUBAHAN, dan sekarang ada buktinya.** Sistem lama memakai **centang** — `IsProRate`
   ditulis hanya oleh `pxCheckbox` layar, dua penulisan, di kedua ekspor (NC-02).
2. **Satu kalimat konteks yang mengubah bobot keputusannya:** prorata di sistem lama **tidak pernah
   dapat menghasilkan apa pun selain 100 %**, karena `EDMEffective` punya satu penulis yang
   menyamakannya dengan `Commencement` dan kontrol layarnya hanya-baca (NC-01, NC-03). Maka ADR-0037
   bukan memperbaiki rumus yang berjalan keliru; ia **memberi arti pada sesuatu yang selama ini
   tidak berarti**.

**Kenapa ini bukan suntingan diam-diam:** butir 2 mengubah cara pembaca menilai prioritas ADR itu.
Tanpa itu, pembaca berikutnya akan mengira ada rumus berjalan yang perlu dipindahkan.

---

## 4. Titipan ke ronde D — tujuh TDA tanpa nasib

**Ditemukan dengan mencocokkan dua tabel di dalam `../KEPUTUSAN-GRILLING-ADJUSTMENT.md` sendiri:**
**Lacak TDA** terhadap **Ronde**. `METODE` §6.3 menuntut setiap TDA membawa nasib — diperbaiki,
dilestarikan, atau ditunda. Tujuh tidak punya satu pun, dan menunjuk cabang yang sudah tidak
bersidang:

| TDA | Pokok | Ditugaskan ke | Keadaan cabang itu |
|---|---|---|---|
| `TDA-01` | penjaga duplikat mati; tabrakan jadi `UPDATE` yang menimpa | cabang B | **DITUTUP 26 Sep** |
| `TDA-07` | dua kontrol layar mengosongkan status akseptasi kontrak | cabang D | **ditutup tanpa ronde** |
| `TDA-09` | peran dari `pyWorkBasketList(2)` / `pyTelephone` / nama tersemat | cabang G | **ditutup tanpa ronde** |
| `TDA-11` | picker menyatukan kontrak dan addendum tanpa pembeda | cabang H | **ditutup tanpa ronde** |
| `TDA-12` | offset pengurai nomor revisi meleset satu | cabang B | **DITUTUP** |
| `TDA-13` | jenis dan materialitas dipilih bebas di radio picker | menunggu `EXP-1` | **`EXP-1` DITUTUP 26 Sep** |
| `TDA-15` | penggandaan pohon di dalam satu `JSONDATA` | cabang F | **ditutup tanpa ronde** |

Ditambah tiga yang tertutup oleh ronde ini dan **perlu ditandai**, bukan diadili ulang:

| TDA | Nasibnya sekarang | Oleh |
|---|---|---|
| `TDA-04` | pemadanan per posisi → **diperbaiki** | `E1` KONFIRMASI — `KUNCI_PADANAN` |
| `TDA-05` | selisih tanpa memeriksa mata uang → **diperbaiki** | `E-1a` — mata uang masuk kunci (ADR-0053) |
| `TDA-06` | ringkasan tak tersimpan → **diperbaiki** | **`GRL-14`** — tidak ada `ActualValue` yang dapat ditimpa |
| `TDA-16` | selisih EGNPI selalu nol → **diperbaiki** | **`GRL-14`** |

### Kenapa ini lubang, bukan kerapian

`METODE` §7.2 menyatakan grilling selesai bila **semua bisa di-spec**. TDA tanpa nasib **tidak bisa
di-spec**: ia sampai ke sesi to-spec sebagai cacat yang tercatat, ditugaskan kepada cabang yang tidak
akan pernah bersidang, dan **tidak ada yang akan tahu bahwa ia menunggu**.

Bentuknya sama dengan yang sudah ditemukan di artefak induk pada hari yang sama: **sebuah keadaan
diperbarui di satu tabel dan tidak di tabel sebelahnya, di dalam berkas yang sama.** Di sana ia
angka kemampuan; di sini ia nasib TDA.

### Perkiraan bentuk ronde D — dan ia perkiraan, bukan hasil

Sebagian besar akan jatuh ke **KONFIRMASI** atau **"diperbaiki oleh GRL-nn"** dengan kutipan. Saya
memperkirakan **tiga** yang menjadi pertanyaan sungguhan: `TDA-07`, `TDA-11`, dan `TDA-13` — yang
terakhir dalam bentuk yang berubah, karena sesudah GRL-12 dan GRL-13 pertanyaannya bukan lagi *"nilai
apa yang sah"* melainkan *"apakah pengisi masih memilih apa pun"*.

Yang memastikannya adalah sapuannya sendiri, bukan perkiraan ini.

---

## 5. Yang BUKAN lubang ronde ini

| Hal | Ke mana |
|---|---|
| paket `REV-1`…`REV-6` belum ditanggapi pemilik ADR | prasyarat cabang **K**; **tidak dapat ditutup sesi mana pun** |
| `PP-1` — siapa mengerjakan to-spec induk | hanya memengaruhi arah dampak GRL-01 |
| tujuh temuan IRISAN ronde C | usulan untuk langkah 8–10 **induk** (GRL-01 butir 5); artefak induk tidak disunting |
| bentuk fisik `NILAI_SELISIH` | sesi DDL induk |
