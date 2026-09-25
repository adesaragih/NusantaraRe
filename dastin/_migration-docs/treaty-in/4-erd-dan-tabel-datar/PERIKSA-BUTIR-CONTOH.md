# Pemeriksaan kesebelas butir dari contoh JSON

**Tanggal:** 24 September 2026

> **Contoh MEMATAHKAN, tidak MENGESAHKAN.** Setiap butir diperiksa **dari ekspor**, dan contohnya
> hanya menunjukkan ke mana harus melihat. Ketidakcocokan diberi nomor.

**Dikerjakan:** 7.10 (**selesai — §7.10-Z**), 7.2 sebagian. **Berikutnya:** 7.1, 7.9, lalu sisanya.

## URUTAN GAMBAR — apa yang ditahan butir mana

| ERD | Ditahan | Sebab |
|---|---|---|
| `ERD-PROPORSIONAL` · `ERD-NON-PROPORSIONAL` | **tidak ada** | 7.10 ditutup §7.10-Z. `CATATAN_PERSETUJUAN` **kembali ke bentuk §10.20**, dan yang berubah adalah **satu entitas baru** `PERISTIWA_KONTRAK` — kotak tambahan, bukan kotak yang digambar ulang |
| **`ERD-ADJUSTMENT-PROP` · `ERD-ADJUSTMENT-NON-PROP`** | **7.1** | `ID_VERSI_KONTRAK_DASAR` dipasang dengan anggapan yang ditunjuk adalah **versi pendahulu**. Bila sebagian `OLDID` menunjuk kontrak yang **DISALIN** dan bukan versi yang **DISESUAIKAN**, penunjuk itu menunjuk **dua benda berbeda** — dan kedua ERD Adjustment digambar di atasnya |

> Ditulis di sini supaya bila 7.1 ternyata rumit dan tergoda ditunda, **yang tertunda bersamanya
> terlihat**.

---

## 7.10 — RANTAI PERSETUJUAN: **LIMA JABATAN, DAN SATU DI ANTARANYA TIDAK PERNAH DAPAT DIMASUKI**

**Dilaporkan sebelum ERD digambar**, sesuai permintaan.

### Apa yang ada di ekspor

Sapuan seluruh nilai `TreatyIn.Position` di 708 berkas:

| Jabatan | Muncul sebagai KONDISI | Muncul sebagai NILAI YANG DISETEL |
|---|---:|---:|
| `ReasTreatyInAdmin` | 10 | **ya** |
| `ReasTreatyInSecHead` | 4 | **ya** |
| `ReasTreatyInDeptHead` | 2 | **ya** |
| `ReasTreatyInDirector` | 6 | **ya** |
| **`ReasTreatyInGroupLeader`** | **8** | **NOL** |

Rantai utamanya terbaca dari `DataTransform/Akseptasi_DT.xml`:

```
Position "" atau ReasTreatyInAdmin   ->  ReasTreatyInSecHead
Position ReasTreatyInSecHead         ->  ReasTreatyInDeptHead     (tolak -> Admin)
Position ReasTreatyInDeptHead        ->  ReasTreatyInDirector     (tolak -> Admin)
Position ReasTreatyInDirector        ->  ReasTreatyInAdmin        (selesai / tolak)
```

**Tiga persetujuan** — SecHead, DeptHead, Direktur. **Persis ADR-0055.**

### Tetapi `ReasTreatyInGroupLeader` ada, dan bentuknya aneh

| Yang diperiksa | Hasil |
|---|---|
| berapa kali diuji sebagai kondisi | **8**, di `Akseptasi_DT`, `TreatyInSetToDirector`, `TreatyInSetValue`, `TreatyInAddNew` |
| berapa kali **disetel** — di seluruh korpus, kedua modul | **NOL** |
| cabangnya di `Akseptasi_DT` | **`pyDisabled = true`** — dinonaktifkan secara tersurat |
| adakah aturan alur yang menyetelnya | **tidak ada folder `Flow` sama sekali** di ekspor Treaty In |

> **Ia keadaan yang diuji tetapi tidak pernah dapat dimasuki** — setidaknya tidak dari ekspor ini.
> Ini bentuk yang belum pernah kita temui: bukan langkah mati, bukan kondisi `NEVER`, melainkan
> **nilai keadaan yang tidak punya penulis**.

### Dan contohnya menunjukkan EMPAT persetujuan, bukan tiga

Contoh `ID 1001701`: satu *"Copied from"*, lalu **`IsApproved: "Accept"` empat kali oleh empat orang
berbeda**, lalu `"Resolve Complete"` oleh orang kelima.

**Dua bacaan, dan ekspor tidak dapat memilih di antaranya:**

| | Bacaan | Akibatnya pada model |
|---|---|---|
| **A** | `GroupLeader` adalah **tingkat keempat yang nyata di produksi**, dan penyetelnya **tidak ada di ekspor ini** | ADR-0055 **kurang satu keadaan** `MENUNGGU_GROUP_LEADER`; K2 kurang satu syarat; P-31…P-33 kurang satu kemampuan |
| **B** | `GroupLeader` **sisa mati**, dan `Accept` keempat adalah **pengajuan pembuatnya sendiri** yang juga tercatat sebagai `Accept` | model tetap; yang perlu diperbaiki hanya pemahaman bahwa `CommentList` mencatat **pengajuan** dan **persetujuan** dengan penanda yang sama |

**Bacaan B punya satu dukungan yang tidak dimiliki A:** contohnya berisi enam entri untuk satu
kontrak, dan yang pertama *"Copied from ID 1001378"* — jadi `CommentList` memang mencatat peristiwa
yang **bukan** persetujuan. Bila ia mencatat penyalinan, ia dapat mencatat pengajuan.

**Bacaan A punya satu dukungan yang tidak dimiliki B:** cabangnya di `Akseptasi_DT` **dinonaktifkan**
— sesuatu dinonaktifkan karena pernah menyala. Dan ini menyambung **L-9**: bila ekspor ini dari QA,
ruleset produksi boleh punya cabang itu **menyala**.

### Akibatnya pada ERD — dan ini yang menentukan apakah menggambar boleh dilanjutkan

> **Bentuk `CATATAN_PERSETUJUAN` TIDAK berubah pada kedua bacaan.**
>
> Entitas itu tidak punya kunci alami, atributnya `WAKTU_KEPUTUSAN`, `NAMA_PEMUTUS`, `DISETUJUI`,
> `ALASAN` (§10.20), dan **tidak satu pun bergantung pada berapa tingkat persetujuan yang ada**.
> Tingkat keempat menambah **baris**, bukan **kolom**.

**Yang berubah bila bacaan A benar** bukan entitas di tengah gambar melainkan:

| Artefak | Perubahannya |
|---|---|
| ADR-0055 | satu keadaan baru `MENUNGGU_GROUP_LEADER`, dan dua perpindahan |
| `SPEC-INVARIAN.md` §3 daftar K2 | satu syarat persetujuan tambahan |
| INV-20 | *"salah satu dari delapan"* menjadi sembilan |
| `DAFTAR-PEKERJAAN.md` | satu kemampuan persetujuan tambahan |
| `4-erd-dan-tabel-datar/` | **nol** |

> ~~**Maka ERD boleh dilanjutkan.** Yang tertahan mesin keadaan, bukan gambar.~~

### ⚠ DICABUT 24 September 2026 — bentuknya MEMANG berubah, dan buktinya ditemukan sesudahnya

Sapuan bentuk penulis (`BENTUK-PENULIS-PROPERTI.md`) membuka **73 berkas *data transform*** yang
tidak pernah tersapu. Salah satunya menulis ke `CommentList`:

```
TreatyIn.CommentList(<APPEND>).IsApproved   = ""          ← KOSONG
TreatyIn.CommentList(<LAST>).OperatorName   = OperatorID.pxInsName
TreatyIn.CommentList(<LAST>).Date           = @CurrentDateTime()
TreatyIn.CommentList(<LAST>).Suggest        = … + "Had Created Internal Edit"
                                            … + "Had Created External Addendum"
                                            … + "Had Created Addendum Premium"
```
*(`DataTransform/TreatyInSetEditPre.xml`, bercabang pada `EDMState` 1, 2, 3)*

**Dua akibat, dan keduanya mengubah sesuatu yang sudah ditetapkan.**

#### a. ~~Bacaan A menguat, bacaan B melemah~~ — **DICABUT LAGI, lihat §7.10-Z**

> ~~`CommentList` mencatat peristiwa bukan-persetujuan dengan **`IsApproved` KOSONG**, bukan dengan
> `"Accept"`. Maka keempat `"Accept"` di contoh **keempatnya persetujuan sungguhan**.~~
>
> **Salah, dan sebabnya sama dengan sebelumnya: saya membaca SATU penulis lalu menyimpulkan tentang
> SELURUH penulis.** `TreatyInSetEditPre` memang menulis `IsApproved = ""`, tetapi ia bukan satu-
> satunya penulis. Ketika seluruh penulis `CommentList` dibuka — §7.10-Z — hasilnya **membalik
> kembali ke bacaan B**, dan kali ini dengan mesin yang terbaca, bukan dengan dugaan.

#### b. `DISETUJUI` TIDAK BOLEH "tidak boleh kosong" — dan itu KOLOM, bukan baris

§10.20 menetapkan `DISETUJUI` bertipe **E**, **boleh kosong: tidak**. Data lama membuktikan
sebaliknya: entri peristiwa punya `IsApproved` **kosong**.

| Pilihan | Bentuknya |
|---|---|
| `DISETUJUI` **boleh kosong** | kosong berarti *"baris ini peristiwa, bukan keputusan"* — dan **kosong yang berarti sesuatu** adalah bentuk yang proyek ini hindari |
| **nilai enumerasi ketiga** — `PERISTIWA` di samping `DISETUJUI` / `DITOLAK` | tidak ada kosong yang bermakna; **usul saya** |
| entitas terpisah untuk peristiwa | menambah entitas, dan peristiwanya sudah punya rumah: `JEJAK_PERUBAHAN` |

> **Saya salah kemarin.** Saya menulis *"tingkat keempat menambah baris, bukan kolom"* — dan itu
> benar untuk tingkat keempatnya. Yang **tidak** saya lihat: `CommentList` memuat **dua jenis baris**,
> dan jenis kedua **tidak muat** di kolom yang sudah ditetapkan. **Kotak `CATATAN_PERSETUJUAN` di
> ERD berubah.**

**Maka ERD `CATATAN_PERSETUJUAN` menunggu satu keputusan** — dan keputusannya diambil di §7.10-Z:
**bukan nilai enumerasi ketiga, melainkan baris peristiwa keluar dari tabelnya.**

---

## 7.10-Z — **DIPUTUSKAN: bacaan B.** Mesin persetujuannya dibaca utuh

Sapuan penulis yang lengkap membuka `DataTransform/Akseptasi_DT.xml` — **mesin keadaan
persetujuannya sendiri**, yang tidak pernah dibaca utuh karena ia *data transform* dan penyisir kami
hanya mengenali aktivitas. Isinya menjawab 7.10 tanpa uji data dan tanpa wawancara.

### Z.1 Rantainya, dibaca dari penulisnya

```
1   RevisionState != 1 && StatusAkseptasi != "Resolve Complete"      (alur biasa)
  1.1  Position=="" || =="ReasTreatyInAdmin"     -> SecHead   StatusAkseptasi "Accept"   << PENGAJUAN
  1.2  Position=="ReasTreatyInSecHead"  + Accept -> DeptHead  StatusAkseptasi "Accept"
  1.3  Position=="ReasTreatyInDeptHead" + Accept -> Director  StatusAkseptasi "Accept"
  1.4  Position=="ReasTreatyInGroupLeader"       -> Director                          *** MATI ***
  1.5  Position=="ReasTreatyInDirector" + Accept -> ""        StatusAkseptasi "Resolve Complete"
2   RevisionState == 1                                               (alur revisi)
  2.1  Position=="" || =="ReasTreatyInAdmin"     -> SecHead   StatusAkseptasi "Accept"   << PENGAJUAN
  2.2  Position=="ReasTreatyInSecHead"  + Accept -> ""        StatusAkseptasi "Resolve Complete"
```

**Dan urutan pemanggilannya menutup rantainya** — `Activity/TreatyInAkseptasi_Act.xml`:

```
1. Apply-DataTransform  (Akseptasi_DT)      <- StatusAkseptasi disetel di sini
2. Call AddCommentList_Act                  <- baris CommentList ditambahkan SESUDAHNYA
3. Call SaveTreatyIn_Act
```

`AddCommentList_Act` menulis **`CommentList(<LAST>).IsApproved = TreatyIn.StatusAkseptasi`** pada
kedua cabangnya. *(Ia juga memuat `IsApproved = "Accept"` tertulis tetap satu baris di atasnya, di
langkah yang sama — dan **ditimpa** oleh baris `StatusAkseptasi` berikutnya. Label tanpa akibat.)*

### Z.2 Kenapa contohnya memuat EMPAT `Accept` sementara tingkatnya TIGA

> **Langkah 1.1 — pengajuan oleh pembuatnya sendiri — menyetel `StatusAkseptasi = "Accept"`.**
> Lalu `AddCommentList_Act` menyalin nilai itu ke barisnya.
>
> **empat `Accept` = satu pengajuan + tiga persetujuan.** Persis ADR-0055.

**Bacaan B benar, dan sekarang ia punya mesinnya**, bukan hanya dukungan tidak langsung. Bacaan A
juga kehilangan penyangganya yang terakhir:

| Penyangga bacaan A | Keadaannya sesudah Z.1 |
|---|---|
| *"cabang `GroupLeader` dinonaktifkan, berarti pernah menyala"* | cabang 1.4 menyalurkan `GroupLeader -> Director`, **sejajar dengan 1.3** (`DeptHead -> Director`). Ia **jalur alternatif pada tingkat yang sama**, bukan tingkat keempat |
| *"ada yang menyetel `Position` ke `GroupLeader` di luar ekspor"* | **tidak satu pun cabang, termasuk yang mati, menyetelnya**. Ia keadaan yang hanya pernah **diuji**, tidak pernah **dituju** |

> **ADR-0055 tetap tiga persetujuan. `MENUNGGU_GROUP_LEADER` TIDAK ditambahkan.** INV-20 tetap
> delapan nilai. Uji AI **dicabut** — jawabannya sudah terbaca dari ekspor, dan uji yang jawabannya
> sudah terbaca adalah **biaya**, bukan ketelitian.

#### Z.2a LINGKUP PENUTUPAN INI — dan kenapa ia harus disebut

Bukti di atas datang dari cabang yang **hidup**. Ia menjawab *"berapa tingkat yang BERJALAN
SEKARANG"*. Ia **tidak** dengan sendirinya menjawab *"berapa tingkat yang PERNAH ADA"* — dan bacaan
A selalu tentang yang kedua. Menyimpulkan dari cabang hidup tentang seluruh cabang adalah bentuk
yang sama dengan kekeliruan Z.a: membaca satu penulis lalu menyimpulkan tentang seluruh penulis.

Karena itu seluruh kode **mati** di jalur persetujuan didaftar tersendiri —
`CABANG-MATI-DI-JALUR-PERSETUJUAN.md`. Hasil yang menyangga penutupan ini:

| Yang dicari di seluruh kode mati jalur persetujuan | Hasil |
|---|---:|
| cabang mati di `Akseptasi_DT` selain 1.4 `GroupLeader` | **nol** |
| langkah mati yang menyetel `Position` ke nilai di luar empat jabatan yang dikenal | **nol** |
| langkah mati yang menambah **tingkat** persetujuan | **nol** |
| beda antara kedua modul pada keluarga ini | **nol** |

> **Maka: tiga tingkat berlaku untuk jalur yang hidup, DAN tidak ada jejak tingkat lain di antara
> 20 langkah mati serta 4 cabang mati yang ada di ekspor.**
>
> **Yang tetap TIDAK dapat ditutup dari sini:** cabang yang **dihapus** — bukan dimatikan — tidak
> meninggalkan jejak apa pun, dan **L-10** berdiri untuk jenis aturan yang tidak terekspor. Kedua
> hal itu ditutup oleh **ekspor kedua**, bukan oleh pembacaan.

Dan kode mati itu memperlihatkan hal lain yang bukan tingkat melainkan **jalan pintas**:
`TreatyInSetToDirector` — seluruh empat langkahnya mati — adalah jalur *"langsung ke Direktur"* yang
pernah ada. Golongannya **P-40** dan **eskalasi butir 1**, bukan ADR-0055.

### Z.3 Temuan baru yang ikut terbuka: **alur revisi hanya punya SATU persetujuan**

Cabang 2 adalah alur *"revisi sesudah diterima"*. Di sana **SecHead menyetujui, lalu kontraknya
langsung `Resolve Complete`** — **DeptHead dan Direktur dilewati.**

| | |
|---|---|
| **Apa** | perubahan atas kontrak yang **sudah disetujui Direktur** dapat disahkan kembali oleh **satu** tingkat |
| **Golongannya** | bukan lubang model — ia **perilaku sistem lama yang harus diputuskan dibawa atau tidak** |
| **Kenapa berat** | tidak ada ambang nilai di cabang itu. Revisi sebesar apa pun melewati jalur yang sama |
| **Ke mana** | `DAFTAR-ESKALASI-MANAJEMEN.md`, dan ia **menuntut jawaban bisnis** sebelum ADR-0055 memuat perpindahan revisinya |

### Z.4 `CommentList` memuat **dua jenis fakta**, dan pemisahnya mekanis

Seluruh penulis `CommentList` di korpus, dan apa yang masing-masing tulis:

| Penulis | `IsApproved` | `Suggest` | Jenis baris |
|---|---|---|---|
| `AddCommentList_Act` | `= StatusAkseptasi` | komentar orang, bebas | **keputusan** — *dan pengajuan, lihat Z.5* |
| `TreatyInDeclineConfirmation_postact(EDM)` | `= "Decline"` | — | **keputusan** |
| `TreatyInSetEditPre` | `= ""` | `"… Had Created Internal Edit / External Addendum / Addendum Premium"` | **peristiwa** |
| `TreatyInCopy` | tidak ditulis | `"Copied from ID " + ID` | **peristiwa** |
| `SetTreatyIn_Act` | tidak ditulis | `"Create Revision"` | **peristiwa** |
| `TreatyInMappingDataconvert` | disalin dari tabel lama | disalin | **warisan** |

> ### ATURAN MIGRASI — dapat dijalankan tanpa menafsirkan satu baris pun
>
> ```
> baris CommentList ber-IsApproved TERISI  ->  CATATAN_PERSETUJUAN
> baris CommentList ber-IsApproved KOSONG  ->  jejak peristiwa
> ```

**Ke mana baris peristiwa itu pergi?** Ujinya satu baris: **apakah ia menyebut ruas yang berubah?**

| Baris peristiwa | Menyebut ruas? |
|---|---|
| `"… Had Created Internal Edit"` | **tidak** |
| `"… Had Created External Addendum"` | **tidak** |
| `"… Had Created Addendum Premium"` | **tidak** |
| `"Copied from ID 1001378"` | **tidak** — ia menyebut **kontrak**, bukan ruas |
| `"Create Revision"` | **tidak** |

> **Tidak satu pun menyebut ruas. Maka `JEJAK_PERUBAHAN` (P-47) TIDAK MUAT** — entitas itu berbutir
> ruas, dan memaksakan peristiwa kasar ke dalamnya berarti membuat baris yang kolom
> `NAMA_RUAS`-nya kosong pada separuh isinya, yaitu bentuk yang baru saja ditolak.
>
> **Maka: entitas ke-28 (atau ke-29 bila `PEMULIHAN_LIMIT` diterima) — `PERISTIWA_KONTRAK`.**
> Dilaporkan sebagai **penambahan**, bukan diselipkan. Atributnya empat:
> `ID_VERSI_KONTRAK`, `WAKTU`, `NAMA_PELAKU`, `JENIS_PERISTIWA` (E) — ditambah
> `ID_KONTRAK_DISALIN_DARI` yang **sudah ada di §10.1**, sehingga peristiwa *"Copied from"* tidak
> menambah kolom teks bebas.
>
> **Dan `CATATAN_PERSETUJUAN` kembali ke bentuk semula:** `DISETUJUI` tetap **tidak boleh kosong**,
> tanpa nilai enumerasi ketiga. §10.20 **tidak berubah**, dan usul `PERISTIWA` **dicabut**.

### Z.5 TABRAKAN YANG DILAPORKAN, BUKAN DIPUTUSKAN — baris pengajuan terbawa ikut

Aturan mekanis di Z.4 memindahkan **setiap** baris ber-`IsApproved` terisi ke
`CATATAN_PERSETUJUAN`. Tetapi baris pertama — **pengajuan** — juga membawa `"Accept"`, karena
langkah 1.1 menyetel `StatusAkseptasi` sebelum barisnya ditulis.

| | |
|---|---|
| **Akibatnya** | `CATATAN_PERSETUJUAN` akan memuat **satu baris pengajuan per kontrak** yang bukan persetujuan. **INV-28** — *"setiap perpindahan meninggalkan satu baris"* — akan **terpenuhi kelebihan satu**, dan setiap hitungan "berapa kali kontrak ini disetujui" **lebih satu** |
| **Dapatkah dipisahkan mekanis?** | **Tidak dari isi barisnya.** Baris pengajuan dan baris persetujuan punya bentuk yang identik: `OperatorName`, `Date`, `IsApproved="Accept"`, `Suggest` = komentar. Yang membedakan hanya **kedudukannya** — ia baris ber-`IsApproved` terisi yang **paling awal** pada satu kontrak |
| **Kenapa tidak saya putuskan** | aturan *"baris paling awal adalah pengajuan"* **rapuh**: kontrak yang pernah **ditolak lalu diajukan ulang** akan punya lebih dari satu pengajuan, dan contoh yang saya pegang tidak memuat kasus itu. Memutuskannya dari satu contoh adalah **menyusun dari contoh** |
| **Yang menutupnya** | **uji data**: sebaran jumlah baris ber-`IsApproved="Accept"` per kontrak, **dipisah** antara kontrak yang pernah ber-`Reject` dan yang tidak. Ditambahkan sebagai **Uji AK** |

### Z.6 Dua bentuk baris yang juga harus punya rumah

| Bentuk | Rumahnya |
|---|---|
| `IsApproved = "Accept"` **tanpa** `Suggest` | **`CATATAN_PERSETUJUAN` dengan `ALASAN` kosong.** `AddCommentList_Act` menulis `Suggest = TreatyIn.Comment`; bila penyetuju tidak mengetik komentar, ia kosong. Sah, dan §10.20 sudah mengizinkan |
| `IsApproved` kosong **dan** `Suggest` kosong | **belum ditemukan satu pun penulisnya.** Bila muncul di data, ia baris tanpa isi — dan itu **temuan migrasi**, bukan bentuk yang dimodelkan. Uji AK menghitungnya sekalian |

### Yang memutuskannya ~~— dicabut, ketiganya sudah terjawab~~

| | |
|---|---|
| **pertanyaan bisnis, satu kalimat** | *"Berapa tingkat persetujuan yang harus dilewati sebuah kontrak treaty masuk, dan apakah ada peran Group Leader di antaranya?"* — jawabannya tiga puluh detik |
| **uji data** | pada `M_TREATY_IN.JSONDATA`: berapa entri `CommentList` ber-`IsApproved = "Accept"` per kontrak. **Modusnya 3** mendukung bacaan B; **modusnya 4** mendukung A. Ditambahkan sebagai **Uji AI** |
| **ekspor kedua** | bila `Akseptasi_DT` produksi memuat cabang `GroupLeader` **menyala**, bacaan A terbukti — dan berkas itu **ditambahkan ke daftar lima belas** |

**Ketiganya dicatat. Tidak satu pun ditebak.**

---

## 7.2 — `ActualValue` TERJAWAB SEBAGIAN: ia potret nilai SEKARANG saat addendum dibuat

Ditemukan lewat sapuan yang sama, di berkas yang sama:

```
TreatyIn.ActualValue.EGNPI = TreatyIn.EGNPI        (DataTransform/TreatyInSetEditPre.xml)
TreatyIn.AddendumPremi     = "1"
```

> **`ActualValue` diisi dari nilai berjalan pada saat addendum dibuat.** Ia bukan cermin kosong dan
> bukan salinan lama — ia **potret "sekarang" yang dibekukan di awal penyesuaian**.

Itu melengkapi mesin selisih, dan **sejalan dengan aturan pemilik proses**:

| Sisi | Dari mana |
|---|---|
| **nilai sekarang** | `ActualValue`, dipotret saat addendum dibuat |
| **nilai lama** | **di-SELECT dari versi sebelumnya** — `OLDDATA` tidak disimpan kembali |
| **selisih** | sekarang − lama |

**Yang masih terbuka:** hanya `EGNPI` yang terlihat dipotret di berkas ini. Apakah seluruh cabang
`ActualValue` dipotret di tempat lain, atau hanya sebagian, **belum disapu** — dan itu menentukan
bentuk tabel nilai selisih (butir 5.6 arahan). Disapu bersama 7.1.
