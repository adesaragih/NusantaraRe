> Modul  : Treaty In Adjustment · ronde C
> Dibuat : 2026-09-24
> Sifat  : temuan baru ronde C (`NC-xx`). Seluruhnya disapu ulang atas ekspor dengan perkakas
>          terkalibrasi; tidak satu pun diwarisi dari berkas penggrill.
> Kolom sisi: ditentukan `tools/cat56.py` (katalog 56 aturan khas), **bukan dari nama berkas**.

# Temuan ronde C

## Ringkasan

| # | Temuan | Sisi | Menyentuh |
|---|---|---|---|
| **NC-01** | `TreatyIn.EDMEffective` punya **tepat satu** penulis nilai, dan ia menyamakannya dengan `Commencement` | 56-KHAS (`TreatyInSetEditPre`) | E3a |
| **NC-02** | `IsProRate` dinyalakan **kotak centang layar**, bukan diturunkan dari tanggal | IRISAN | E3a, ADR-0037 |
| **NC-03** | Penghitung `ProRatePercent` dirujuk **hanya dari aturan layar**, tidak dari satu pun Activity | IRISAN | E3a |
| **NC-04** | Selisih share fakultatif **tidak pernah dihitung** di jalur pengajuan addendum — dua sebab yang berdiri sendiri | IRISAN | E3b |
| **NC-05** | `ActualValue` berarti **dua hal yang berbeda** menurut `EDMState` | IRISAN | C3 |
| **NC-06** | Payung I2 **sudah terjawab** ADR-0042 + ADR-0043 + ADR-0054; yang tersisa satu hal, dan ia bukan yang tertulis di daftar | — | I2 |
| **NC-07** | Seluruh mesin selisih addendum berada di **IRISAN** — ia ada di ekspor Treaty In juga | IRISAN | GRL-01 butir 5 |

---

## NC-01 — `EDMEffective` tidak dapat digeser siapa pun

`tools/tulis.py sapu .EDMEffective` — 5 penulisan, dan **tidak satu pun** membuatnya berbeda dari
tanggal mulai kontrak:

| Bentuk | Berkas | Isi |
|---|---|---|
| DT | `TreatyInSetEditPre` 1.5 | `TreatyIn.EDMEffective = TreatyIn.Commencement` |
| DT | `TreatyCalculateProratePct` 4.2 (dua ekspor) | sasarannya **`Param.EDMEffective`**, bukan `TreatyIn.EDMEffective` — ia menambahkan komponen jam, bukan menetapkan tanggal |
| BIND | `Section/TreatyInNONProportional` (dua ekspor) | `pxDateTime` **hanya-baca** |

> Penyapu mencocokkan akhiran `.EDMEffective`, sehingga `Param.EDMEffective` **ikut terjaring**.
> Itu dinyatakan di sini supaya angka 5 tidak dibaca sebagai lima penulis `TreatyIn.EDMEffective`.
> Penulis `TreatyIn.EDMEffective` yang sebenarnya: **satu**.

**Ditolak perkakas:** 0 simpul salinan terbungkus, 0 berkas gagal urai.
**Titik buta yang perkakas nyatakan sendiri:** 32 langkah Java, 191 langkah SQL/REST — keduanya
**tidak** dapat dinyatakan nol.

## NC-02 — pro rata dinyalakan orang, bukan diturunkan dari tanggal

`tools/tulis.py sapu .IsProRate` — **2 penulisan, keduanya BIND**, `pxCheckbox`, **dapat disunting**,
di `Section/TreatyInNONProportional.xml` pada **kedua** ekspor. Nol penulisan dari Activity, Data
Transform, kontrol Section, maupun Report Definition.

ADR-0037 induk menetapkan *"faktor prorata addendum = turunan dari tanggal, bukan centang"*. Sistem
lama adalah **centang**. Maka ADR-0037 berlabel **PERUBAHAN**, dan ia berubah dari sesuatu yang
terbaca, bukan dari dugaan.

## NC-03 — dan persennya dihitung di layar, bukan saat pengajuan

`tools/dt.py TreatyCalculateProratePct`:

```
1    SET  TreatyIn.ProRateDays        0
2    SET  TreatyIn.ProRateTotalDays   0
3    SET  TreatyIn.ProRatePercent     0
4    WHEN TreatyIn.IsProRate == true
4.4    SET ProRateDays       = @DateTimeDifference(Param.EDMEffective, Param.Termination, 'D')
4.5    SET ProRateTotalDays  = @DateTimeDifference(Param.Commencement, Param.Termination, 'D')
4.6    SET ProRatePercent    = @divide(ProRateDays, ProRateTotalDays)
4.7    SET ProRatePercent    = ProRatePercent * 100
```

Pencarian teks pemanggil di **kedua** ekspor mengembalikan **hanya aturan layar** —
`Harness/InputTreatyInOffer`, `Section/InputTreatyInOffer`, `Section/TreatyInNONProportional` — dan
**nol Activity**. `TreatyEDMCalculateDifference` langkah 11 memanggil `TreatyEDMProRateCalculation`
(pengali), **bukan** penghitung persennya.

### Akibat yang ditarik dari NC-01 + NC-03, dan ia aritmetika, bukan tafsiran

Karena `EDMEffective` selalu sama dengan `Commencement` (NC-01):

```
ProRateDays      = selisih(Commencement, Termination)
ProRateTotalDays = selisih(Commencement, Termination)
ProRatePercent   = 1 × 100 = 100
```

> **Pro rata di sistem lama tidak pernah dapat menghasilkan apa pun selain 100 %.** Ia tidak mati,
> tidak tersembunyi, dan labelnya jujur — yang membuatnya tidak berarti adalah bahwa **satu-satunya
> masukan yang dapat membedakannya dipaku pada sebuah konstanta**.

Dan satu ujung yang **tidak** dapat dijawab ekspor: bila `IsProRate` dicentang tanpa peristiwa layar
yang menjalankan penghitungnya, `ProRatePercent` bernilai apa saat mesin selisih mengalikannya.
Langkah 3 menyetelnya **0**. Itu **`UA` baru**, bukan kesimpulan.

## NC-04 — share fakultatif: dua ketiadaan yang berdiri sendiri

| # | Ketiadaan | Bukti |
|---|---|---|
| 1 | `TreatyEDMDifferenceDeduction` langkah **4 dan 5 MATI** | `tools/pre.py`, label aslinya *"Facultative share not enable yet in this edm"* |
| 2 | `TreatyInDifferenceFacShare` **tidak dipanggil** dari rantai `TreatyEDMCalculateDifference` | pemanggilnya hanya dirinya sendiri, `…FacShareTotal`, dan `TreatyInActualUpdateValueShare` — yang di ekspor Adjustment dirujuk **hanya** `Section/TreatyInActualShare` |

Akar `FacultativeShare` (#3) dan `FacultativeShareBrokerage` (#4) pada daftar 33 akar
(`PENGETAHUAN.md` §5.6) berpenulis **tunggal** `TreatyInDifferenceFacShare`. Digabung dengan butir 2
di atas: **kedua akar itu tidak pernah terisi lewat jalur pengajuan addendum.**

Yang kedua lebih penting daripada yang pertama, dan bedanya layak disebut: butir 1 adalah kemampuan
yang **dimatikan dengan sengaja**; butir 2 adalah kemampuan yang **tidak pernah disambungkan**.
Menyalakan blok mati tidak menutup butir 2.

## NC-05 — `ActualValue` berarti dua hal, dan pembedanya `EDMState`

Dari `PENGETAHUAN.md` §2.2 dan §5.4, keduanya dikutip apa adanya:

| `EDMState` | Apa yang terjadi pada `ActualValue` | Maka artinya |
|---|---|---|
| **1, 2** | `SaveTreatyIn_EDM_Act` langkah 2–4 **menimpanya** dengan salinan seluruh halaman dikurangi `OLDDATA` dan `ValueDifference` | **potret sesudah-simpan** — keluaran, bukan masukan |
| **3** | langkah 2 **dilompati**; pohon itu justru yang **disunting pengguna**, disemai `TreatyInSetEditPre` 3.5 (`ActualValue.EGNPI = EGNPI`) | **masukan pengguna** — nilai baru yang sedang diajukan |

Ini bentuk yang sudah dilarang dibawa ke sistem baru oleh aturan induk: **satu nama untuk dua arti.**

Dan ia bertemu dengan **TDA-16** (`PENGETAHUAN.md` §5.4, temuan 27 September): mesin selisih premi
beriterasi `TreatyIn.EGNPI`, sedangkan pengguna menyunting `TreatyIn.ActualValue.EGNPI`. Maka pada
addendum premi **selisih EGNPI selalu nol** — angka yang menjadi alasan jenis addendum itu ada tidak
pernah muncul sebagai selisih.

> Di bawah **GRL-12** (material = versi punya sedikitnya satu baris `NILAI_SELISIH` tersimpan),
> sebuah `PENYESUAIAN_PREMI` yang **hanya** mengubah premi aktual terbaca **non-material** —
> bertentangan dengan **GRL-13** yang mencatat *"penyesuaian premi selalu material"*.
>
> Itu bukan dua keputusan yang bertabrakan. Itu **satu keputusan yang bertabrakan dengan satu
> cacat**, dan C3 yang menentukan mana yang mengalah.

## NC-06 — payung I2 sudah terjawab; sisanya satu hal, dan ia tidak ada di daftar

Ketiga ADR dibaca (prasyarat `MA-10`), dan **payung I2** — *"aturan baru mana yang boleh dijalankan
atas baris warisan, dan asal-usulnya dicatat bagaimana"* — **sudah dijawab**:

| Sumber | Kutipan yang menjawab |
|---|---|
| **ADR-0042** butir 2 | *"Invarian ditegakkan pada penulisan, bukan pada baris. Baris warisan ditulis sekali oleh migrasi dan tidak bisa disunting, jadi ia tidak pernah melanggar apa pun."* |
| **ADR-0042** butir 3 | aturan **sentuh-perbaiki**: begitu kontrak warisan disentuh, ia memenuhi invarian saat itu juga dan penanda warisannya dicabut |
| **ADR-0042** §materialitas | kolom membawa **asal-usulnya** — diturunkan aturan, atau diimpor sebagai warisan |
| **ADR-0043** | migrasi **memindahkan**; menghitung ulang adalah **peristiwa bisnis tersendiri**, dan kriterianya berlawanan |
| **ADR-0054** | nilai warisan tanpa padanan mendarat di `WARISAN_TAK_TERPETAKAN`, nilai aslinya tersimpan |

**Maka payung I2 turun menjadi KONFIRMASI**, dan anggaran GRILL turun **4 → 3 butir + satu residu**.

**Residu yang ketiga ADR TIDAK jawab**, dan ia satu-satunya yang tersisa:

> ADR-0042 dan ADR-0043 menjawab *"apa yang boleh dilakukan atas nilai warisan"*. Keduanya **tidak**
> menjawab *"berapa `NOMOR_URUT_VERSI` sebuah baris warisan"* — dan itu bukan nilai warisan, ia
> **mekanisme sistem baru** yang harus diberikan kepada baris lama.

Dan itu menyentuh **GRL-11** langsung: versi berlaku = versi `DISETUJUI` dengan `NOMOR_URUT_VERSI`
**tertinggi**. Syaratnya *"urutan nomor = urutan persetujuan"*. Untuk baris warisan syarat itu
**tidak dapat dipenuhi dari nomor lamanya**, karena `PENGETAHUAN.md` §4.4 membuktikan penomoran
`/Rnn` **meleset dan dapat dipakai ulang** — memilih `…/R01` ketika R02 sudah ada menghasilkan R02
lagi, dan tabrakannya berakhir sebagai `UPDATE` yang menimpa (TDA-01).

> Nomor yang dapat dipakai ulang bukan urutan. Menurunkan "versi berlaku" darinya berarti menurunkan
> keadaan hukum sebuah kontrak dari sebuah penghitung yang sudah terbukti salah.

## NC-07 — mesin selisih addendum ada di IRISAN, seluruhnya

`tools/cat56.py` mengembalikan **56 aturan** khas Adjustment. Tidak satu pun dari berikut ini ada di
dalamnya:

```
TreatyEDMCalculateDifference   TreatyEDMDifferencePremium   TreatyEDMDifferenceShare
TreatyEDMDifferenceDeduction   TreatyEDMProRateCalculation  TreatyCalculateProratePct
TreatyInDifferenceFacShare     TreatyInNONProportional      SaveTreatyIn_EDM_Act
```

Hanya `TreatyInSetEditPre` yang 56-KHAS.

**Akibatnya pada ke mana temuan dikirim:** NC-02, NC-03, NC-04, NC-05, dan NC-07 sendiri adalah
**temuan Treaty In** — cacat yang berjalan di modul induk hari ini. Ronde ini memutuskan nasibnya
**untuk jalur addendum saja**, dan mencatatnya sebagai usulan untuk langkah 8–10 induk (GRL-01
butir 5).

Dan satu kalimat yang pantas dicatat karena ia mengubah cara membaca seluruh modul: **"mesin selisih
addendum" bukan milik modul Adjustment.** Ia dikirim bersama ekspor induk. Permukaan khas Adjustment
ternyata lebih tipis lagi daripada yang sudah dicatat §1.3.
