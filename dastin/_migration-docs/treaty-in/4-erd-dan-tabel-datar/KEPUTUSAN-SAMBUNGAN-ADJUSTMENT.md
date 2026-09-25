# Keputusan sambungan Treaty In ⟷ Treaty In Adjustment

**Tanggal:** 23 September 2026
**Keadaan:** hasil gerbang G1, G2, G3. **Ditulis sebelum satu kotak pun digambar.**
**Batas:** embargo Adjustment dibuka **sebagian** — hanya untuk mengetahui entitas, sambungan, dan
bentuk tabel selisih. Perilaku Adjustment **tidak diadili**; temuan tentangnya ada di
`TEMUAN-ADJUSTMENT-DITUNDA.md`.

---

> **KOREKSI, 23 September 2026 — embargo dicabut, §G2 dikoreksi.**
> Modul Treaty In Adjustment sudah dibedah; hasilnya di
> `../treaty-in-adjustment/PENGETAHUAN.md`. Satu kesimpulan berkas ini **tidak bertahan**:
> §G2 menyatakan nilai selisih historis tidak dapat direproduksi karena `OLDID` menunjuk baris
> yang bergerak. Ternyata mesin selisih tidak pernah membaca lewat `OLDID` — ia membaca
> `TreatyIn.OLDDATA`, potret penuh kontrak yang dibekukan ke dalam `JSONDATA` addendum itu sendiri
> (`Activity/TreatyInSetAddendumToHistory.xml`). Nilai lama **dapat** direproduksi.
> Selebihnya — G1, G1b, G3, dan G4 — bertahan.


## 0. Temuan yang mendahului ketiga gerbang, dan mengubah cara membacanya

**Kedua ekspor memuat aturan yang SAMA — byte-identik, bukan mirip.**

| Aturan | Ruleset · versi | `pzInsKey` |
|---|---|---|
| `Akseptasi_DT` | GISFW 01-01-94 | `…AKSEPTASI_DT #20260610T095554.165 GMT` |
| `SaveTreatyIn_Act` | GISFW 01-01-96 | `…SAVETREATYIN_ACT #20260901T081410.307 GMT` |
| `TreatyInForceEdit` | GISFW 01-01-55 | `…TREATYINFORCEEDIT #20190903T024752.625 GMT` |

Kunci instansnya identik. Maka **"Treaty In Adjustment" bukan aplikasi terpisah dengan modelnya
sendiri** — ia ekspor dari **ruleset yang sama**, diambil dari titik masuk yang berbeda.

| | Jumlah berkas |
|---|---|
| Treaty In | 329 |
| Treaty In Adjustment | 379 |
| **Irisan** | **323** |
| Hanya di Adjustment | **56** |
| Hanya di Treaty In | 6 |

**Permukaan Adjustment yang sebenarnya adalah 56 berkas itu**, dan bentuknya berbicara sendiri: 14
seksi/flow-action ber-akhiran `OldData`, tujuh aturan ber-`EDM`, dua *picker* master, dan SQL
`TreatyLoadMasterJoinEdm`.

> **Ini penerapan CONTEXT.md §2.0.** Membaca 379 berkas sebagai "model Adjustment" akan menghasilkan
> model yang 85 % adalah Treaty In. Yang terlihat sebagai dua modul adalah satu ruleset dengan dua
> pintu.

Satu kelas yang **benar-benar hanya ada di Adjustment**: `Int-treaty_in_edm`, 19 properti —
`ID`, `OLDID`, `EDMDate`, `EDMState`, `EDMMaterialType`, `Position`, `PositionUsername`,
`StatusAkseptasi`, `ChooseStatusAkseptasi`, `Ceding`, `CedingID`, `LeadingReinsSource`,
`LeadingReinsSourceID`, `ProportionType`, `Commencement`, `Termination`, `ClassOfBusiness`,
`TreatyContractName`, `Information`.

---

## G1 — Adjustment itu VERSI, bukan dokumen tersendiri

### Uji yang diminta, dan kenapa ia tidak memisahkan di model ini

Uji yang ditetapkan: *apakah Adjustment punya daur hidup dan persetujuannya sendiri, terpisah dari
versi yang disesuaikannya?*

**Jawabannya YA**, dan buktinya tegas:

| Bukti | Berkas |
|---|---|
| `Int-treaty_in_edm` membawa `Position`, `StatusAkseptasi`, `ChooseStatusAkseptasi`, `PositionUsername` sendiri | inventaris kelas, ekspor Adjustment |
| Rantai persetujuan khusus EDM | `TreatyInSubmitEDM`, `TreatyInAkseptasiEDM_Act`, `TreatyInDeclineConfirmation_postactEDM` |
| Penyimpanan khusus EDM | `SaveTreatyIn_EDM_Act`, `SaveTreatyInDetailEdm_Act` |

**Tetapi uji itu tidak memisahkan di model kita**, dan itu harus dikatakan terang daripada dipakai
diam-diam: ADR-0055 sudah menetapkan bahwa **setiap versi kontrak menempuh persetujuannya sendiri**.
Punya daur hidup sendiri karena itu **ciri versi di model ini**, bukan tanda sesuatu yang bukan
versi. Uji itu memisahkan pada model yang versinya tidak punya daur hidup; model kita bukan itu.

### Pembeda yang memisahkan, dan jawabannya

Pertanyaan yang benar-benar memisahkan: **apakah baris Adjustment memuat DELTA, atau memuat
DOKUMEN KONTRAK YANG UTUH?**

- delta + penunjuk → dokumen tersendiri;
- dokumen utuh + identitas sendiri + penunjuk ke pendahulunya → **versi**.

Buktinya satu baris, dari `RDBList/SaveTreatyInEDM.xml`:

```
BEGIN POOLDATA.PEGA_M_TREATY_IN_EDM(
   {InputParam.DATAPEGA},        -- JSONDATA: SELURUH pohon kontrak
   {TreatyIn.ID}, {TreatyIn.OLDID}, {InputParam.STSCARI1},
   {TreatyIn.ProportionType}, {TreatyIn.TreatyContractName}, {TreatyIn.TeritorialScope},
   {TreatyIn.Commencement}, {TreatyIn.Termination}, {TreatyIn.ClassofBusiness},
   {TreatyIn.LeadingReinsSource}, {TreatyIn.LeadingReinsSourceID},
   {TreatyIn.Ceding}, {TreatyIn.CedingID}, {TreatyIn.LeadingR… )
```

Sebuah baris EDM menyimpan **`JSONDATA` penuh** — bentuk yang sama persis dengan `M_TREATY_IN` —
ditambah kolom kepala yang sama. Ia bukan catatan perubahan; ia **salinan lengkap kontrak pada satu
titik waktu, dengan identitasnya sendiri dan penunjuk ke pendahulunya.**

Itu definisi sebuah versi.

### Keputusan G1

> **Adjustment adalah VERSI_KONTRAK.** Ia bukan entitas baru.

Konsekuensinya, dan ini yang membuat keputusan ini murah: **tidak ada entitas Adjustment yang perlu
dibuat.** Yang perlu dibuat hanyalah entitas **selisih** (G3), dan sambungan yang membuat selisih
itu punya dasar (G2).

Ini sejalan dengan ADR-0040 §1 — *"addendum berhenti menjadi kolom dan menjadi hubungan versi"* —
dan sekarang ia **terbukti dari ekspor**, bukan hanya diputuskan.

---

## G1b — apa yang dibaca hilir setelah sebuah addendum disetujui

Uji ini dijalankan tiga kali, dan **dua kali pertama saya salah** karena membaca daftar langkah
tanpa mengurai sarangnya. Yang di bawah adalah hasil penguraian pohon `pySteps` yang benar.

### Kedua aktivitas simpan, diurai per langkah

| `SaveTreatyIn_Act` (kontrak) — 11 langkah | | `SaveTreatyIn_EDM_Act` (addendum) — 17 langkah | |
|---|---|---|---|
| 8 `call SaveTreatyInDetail_Act` | **hidup** | 12 RDB-List *"insert into m treaty in"* | **MATI** |
| 9 `call SaveTreatyInOffer_Act` | **MATI** | 13 `call SaveTreatyInDetail_Act` | **MATI** |
| | | 14 `call SaveTreatyInDetailEdm_Act` | **hidup** |

### Uji salin-tempel — hasilnya: BUKAN sisa salin-tempel

Pemilik proses mengajukan tafsir tandingan: langkah mati itu mungkin sisa salinan dari aktivitas
kontrak, dimatikan **justru karena tidak berlaku** untuk addendum — yang berarti maksudnya
**penolakan**, bukan penggantian.

Dua pemeriksaan memisahkannya, dan keduanya menjawab hal yang sama:

| Pemeriksaan | Hasil |
|---|---|
| **(a) Adakah padanannya di aktivitas kontrak?** | **Tidak.** Aktivitas kontrak **tidak punya** langkah berdeskripsi *"when Resolve Complete insert into m treaty in"*. Penyisipan master kontraknya adalah langkah 4, tanpa deskripsi dan tanpa syarat keadaan. Langkah 12 di aktivitas EDM **tidak punya asal-usul untuk disalin.** |
| **(b) Apakah polanya salin-lalu-tambah?** | **Tidak.** Langkah khas EDM (14) berada **di tengah**, bukan di belakang; ekornya (15–17) sejajar dengan ekor aktivitas kontrak (10–11). Urutannya disusun, bukan ditempel. |

**Maka kalimat "bukti maksud" bertahan** — tetapi sekarang ia teruji, bukan diasumsikan. Langkah 12
dan 13 ditulis **untuk** alur addendum, dengan syarat keadaan `Resolve Complete` yang hanya masuk
akal di alur itu, lalu dimatikan.

### Dua klaim saya yang harus dicabut

**Cabutan 1 — `TREATYINOFFER` tidak ditulis oleh jalur kontrak juga.** Saya menulis "jalur EDM nol
panggilan ke `SaveTreatyInOffer_Act`, jalur kontrak dua". Rujukannya memang dua, tetapi
**langkahnya MATI** (langkah 9). Saya menghitung rujukan, bukan langkah hidup — kesalahan yang sama
untuk ketiga kalinya dalam satu pemeriksaan.

Penulis `TREATYINOFFER` yang benar-benar hidup hanya satu: `TreatyInPopulateDetail` langkah 2.4,
dan **ia tidak dipanggil aktivitas mana pun** — hanya oleh kontrol di layar `InputTreatyInOffer`.

**Cabutan 2 — "bukan sebagian" tidak punya dasar yang cukup.** Saya menyimpulkan setiap addendum
tidak terlihat hilir dari **satu** aktivitas. Dengan 323 berkas bersama, layar addendum sangat
mungkin mencapai penulis lain yang belum saya telusuri. Diamnya satu aturan bukan diamnya sistem —
kalimat yang sudah dipakai tiga kali di sesi ini, dan saya melanggarnya sendiri.

### Yang tersisa sebagai klaim yang berdiri

> `TREATYINDETAIL` — tabel datar rincian kontrak — menerima baris dari jalur simpan **kontrak**
> (langkah 8 hidup) dan **tidak** dari jalur simpan **addendum** (langkah 13 mati). Addendum
> menulis ke `TREATYINDETAILEDM` sebagai gantinya.

Itu terbukti dari blok, bukan dari rujukan. Selebihnya — apakah ada jalur lain yang membawa
addendum ke hilir — **dijawab data, bukan kode**: **Uji Z** di berkas permintaan DBA.

---

## G1c — pemisahan fakta dari rancangan, dan kenapa ADR tidak memutuskan fakta

Saya menyerahkan "tensi" ini kepada pemilik proses. Koreksinya diterima dan ia penting melampaui
kasus ini:

> **ADR adalah keputusan yang kita buat; ia menetapkan apa yang akan DIBANGUN. Ia tidak dapat
> menetapkan apa yang sistem lama LAKUKAN, karena ia lahir belakangan.** Urutan wewenang
> menyelesaikan pertikaian tentang apa yang dibangun; ia tidak pernah dimaksudkan membungkam bukti
> tentang apa yang ada.

Memakainya begitu akan membuat ADR **kebal terhadap temuan**, dan itu persis yang tidak dikehendaki.
Dipisahkan, keduanya terjawab sekaligus dan tensinya hilang:

| Pertanyaan | Jawaban | Yang berwenang |
|---|---|---|
| **FAKTA** — apakah sistem lama mengganti atau berdampingan? | **BERDAMPINGAN.** Baris kontrak tidak pernah diperbarui oleh jalur addendum. Disengaja atau terlantar tidak mengubah faktanya. | bukti ekspor |
| **RANCANGAN** — apakah sistem baru mengganti? | **YA**, penuh. | ADR-0040 |

**G1 berdiri — tetapi karena ia pertanyaan rancangan, bukan karena bukti sistem lama mendukungnya.**
Bukti sistem lama justru melawannya, dan kita tetap memilih mengganti.

### Akibatnya wajib: ini PERUBAHAN, bukan pelestarian

Ditandai di `SPEC-MODEL-DATA.md` dan di ADR-0040 sendiri:

> Di sistem lama, addendum yang disetujui **tidak menggantikan apa pun** — ia tersimpan di tabelnya
> sendiri dan baris kontrak tetap seperti semula. Di sistem baru, **versi yang disetujui
> menggantikan pendahulunya**, dan hilir membaca versi yang berlaku.

Tanpa penandaan itu, seseorang setahun lagi akan menemukan perilaku lama, mengira spesifikasinya
keliru, dan "memperbaikinya" kembali ke keadaan rusak.

---

## G2 — `OLDID` menunjuk SEBUAH BARIS, dan barisnya bisa berpindah. Ini temuan.

### Apa yang benar-benar disimpan

`OLDID` ditulis di tiga tempat, dan ketiganya sama:

```
TreatyInCopy          :  TreatyIn.OLDID  <-  TreatyIn.ID
TreatyInEDMSetValue   :  TreatyIn.OLDID  <-  TreatyIn.ID
TreatyInRevisi_post   :  TreatyIn.OLDID  <-  TreatyIn.ID
```

`TreatyIn.ID` adalah **ID baris yang sedang dibuka**. Baris apa? Jawabannya ada di *picker*-nya,
`RDBList/TreatyLoadMasterJoinEdm.xml`:

```sql
select a.ID … from pooldata.treaty_in a
UNION
select b.ID … from pooldata.treaty_in_edm b
```

Yang dapat dipilih sebagai "yang disesuaikan" adalah **baris kontrak ATAU baris addendum**, dan
keduanya disajikan dalam satu daftar tanpa pembeda. Maka `OLDID` memuat **ID baris mana pun yang
dipilih**.

### Kenapa itu bukan penunjuk versi

Di sistem lama **tidak ada identitas versi sama sekali**. `M_TREATY_IN` punya **satu baris per
kontrak**, dan baris itu memuat **keadaan kontrak saat ini** — ia disunting, bukan ditambah.

Maka ketika `OLDID` menunjuk sebuah baris `treaty_in`, ia menunjuk **sasaran yang bergerak**:

| Saat | Yang dibaca sebagai "nilai lama" |
|---|---|
| hari selisih dihitung | nilai kontrak hari itu |
| setahun kemudian | nilai kontrak **setelah** disunting entah berapa kali |

### Akibatnya, dinyatakan dan tidak diperbaiki diam-diam

> **Nilai selisih historis tidak dapat direproduksi.** Menghitung ulang selisih sebuah Adjustment
> lama hari ini akan menghasilkan angka yang berbeda dari yang dibukukan dulu, **tanpa ada seorang
> pun mengubah Adjustment itu.**

Ini bukan hipotesis: ia akibat langsung dari tidak adanya identitas versi, ditambah `M_TREATY_IN`
yang disunting di tempat. Ia juga persis alasan ADR-0036 disebut **prasyarat teknis**, bukan pilihan
rancangan — dan sekarang alasannya terbukti, bukan diperkirakan.

Masuk daftar temuan Treaty In. **Tidak diperbaiki diam-diam.**

### Keputusan G2 — ditandai sebagai PERUBAHAN, bukan pelestarian

> **Di sistem baru, penunjuknya adalah `ID_VERSI_KONTRAK`, bukan `ID_KONTRAK`.**

| | Sistem lama | Sistem baru |
|---|---|---|
| Yang ditunjuk | baris kontrak **atau** baris addendum, tanpa pembeda | **satu versi kontrak tertentu** |
| Sifat sasaran | dapat berubah nilainya | **beku** — INV-24, ADR-0036 |
| Reproduksi selisih | tidak mungkin | mungkin, bertahun-tahun kemudian |

Dua penunjuk disimpan, bukan satu:

| Nama | Bita | Artinya |
|---|---|---|
| `ID_VERSI_KONTRAK_LAMA` | 21 | versi yang menjadi **dasar** selisih |
| `ID_VERSI_KONTRAK_BARU` | 21 | versi yang **dihasilkan** penyesuaian |

`SELISIH = NILAI_BARU − NILAI_LAMA`, dan keduanya dibaca lewat bentuk baca
`NILAI_VERSI_KONTRAK` (`SEAM-ADJUSTMENT.md` §3) — **di-SELECT, tidak disalin**.

Karena kedua versi beku, selisih yang tersimpan **dapat direkonsiliasi ulang kapan saja**, dan
ketidakcocokan menjadi **alat deteksi**, bukan masalah.

---

---

## G4 — Apakah `treaty_in` dan `treaty_in_edm` berbagi ruang pengenal?

Pertanyaan ini naik dari pemeriksaan sampingan menjadi **prasyarat**: bila nilai `OLDID` dapat cocok
di kedua tabel sekaligus, rantai addendum tidak dapat ditelusuri sama sekali, dan migrasi tidak
dapat menyusun ulang keadaan kontrak yang berlaku.

**Jawabannya: tidak berbagi sequence, tetapi tabrakan nilai MUSTAHIL** — karena pengenal addendum
**diturunkan dari pengenal kontrak dengan penyambungan teks**.

Bukti, dari `Activity/TreatyInRevisi_post.xml`:

```
TreatyIn.ID <- TreatyIn.ID + "/R01"
TreatyIn.ID <- @If(@substring(TreatyIn.ID,10,12) < 10,
                   @substring(TreatyIn.ID,0,7) + "/R0" + (@toInt(@substring(…
```

Sebuah pengenal addendum selalu berbentuk `‹pengenal kontrak›/Rnn`. Karena akhiran itu selalu ada,
**tidak ada nilai yang dapat muncul di kedua tabel**, dan `OLDID` tidak ambigu secara nilai — hanya
secara tipe, dan tipenya bahkan terbaca dari bentuk nilainya sendiri.

**Akibatnya untuk migrasi: rantai addendum DAPAT ditelusuri.** Itu kabar baik, dan ia prasyarat yang
terpenuhi.

> ### SYARAT yang ditambahkan 24 September 2026 — rantainya dapat kehilangan mata rantai
>
> Kesimpulan di atas ditulis **sebelum** diketahui bahwa barisnya dapat **dihapus**.
>
> `TreatyInDeclineConfirmation_postactEDM` langkah 6 dan 7 — keduanya **hidup** — menghapus baris di
> `M_TREATY_IN_EDM` dan `TREATY_IN_EDM` setelah sebuah addendum ditolak. Jalur kontrak biasa tidak
> melakukannya.
>
> **Rantai yang mata rantainya dapat hilang bukan rantai yang utuh.** Kesimpulan G4 karena itu
> **tidak dibatalkan, melainkan diberi syarat:**
>
> | | |
> |---|---|
> | **Berlaku** | rantai dapat ditelusuri **di antara addendum yang masih ada** |
> | **Tidak dijamin** | bahwa yang masih ada adalah **seluruh** addendum yang pernah dibuat |
> | **Diuji** | **Uji AC** — lompatan dan pengulangan nomor revisi |
>
> Dua akibat yang menunggu hasil Uji AC:
>
> 1. **Nomor revisi mungkin dipakai ulang.** Pengenal berbentuk `‹kontrak›/Rnn` dan nomornya
>    dihitung dari baris yang **ada**. Bila `/R02` dihapus, addendum berikutnya dapat memperoleh
>    `/R02` lagi — sehingga **dua addendum berbeda pernah memakai pengenal yang sama pada waktu yang
>    berbeda**, dan setiap dokumen di luar sistem yang menyebut `‹kontrak›/R02` menjadi tidak tertentu.
> 2. **Kueri penomorannya sendiri cacat** — `ROWNUM = 1` mendahului `ORDER BY`, yang di Oracle
>    mengembalikan baris **sembarang**, bukan tertinggi. Dua cacat yang bergandengan.
>
> Di model baru keduanya tidak terbawa: pengenal versi dibangkitkan `SEQUENCE` (INV-02), dan
> `NOMOR_URUT_VERSI` unik per kontrak (INV-04) atas baris yang **tidak pernah dihapus**
> (ADR-0055 perubahan 24 Sep 2026).

**Tetapi caranya membawa dua hal yang harus dicatat**, dan keduanya **tidak diadili di sini**:

| Yang terlihat | Bukti |
|---|---|
| Pengenal dibentuk dengan **menyunting teks** — pola yang dilarang di model baru | `TreatyIn.ID + "/R01"` |
| Penguraiannya memakai **offset tetap** — `@substring(ID,0,7)` dan `SUBSTR(ID,10,2)` — yang mengandaikan pengenal kontrak selalu selebar tujuh karakter | `TreatyInRevisi_post`, `GetTreatyRevisionID` |

Keduanya masuk `TEMUAN-ADJUSTMENT-DITUNDA.md`. Di model baru pengenal versi dibangkitkan
`SEQUENCE` dan tidak pernah dari teks (INV-02), sehingga bentuk ini tidak terbawa.

Satu hal yang **belum** dipastikan dan tidak saya tebak: untuk revisi kedua dan seterusnya, apakah
`OLDID` menunjuk **kontrak** atau **revisi sebelumnya**. Kedua bentuk dapat dihasilkan urutan
langkah di `TreatyInRevisi_post`, dan memastikannya menuntut membaca urutan langkah aktivitas itu
sampai habis — yaitu mengadili perilaku Adjustment. **Ia ditunda**, dan ia **tidak menghalangi**
rancangan: model baru menyimpan `ID_VERSI_KONTRAK_DASAR` secara eksplisit, sehingga bentuk rantainya
tidak perlu disimpulkan.

## G3 — Tabel nilai selisih: SEMPIT, dengan kosakata besaran TERTUTUP

### Apa yang ada di sistem lama

**Tidak ada tabel selisih.** Sapuan SQL atas kedua ekspor tidak menemukan satu pun tabel yang
menyimpannya. Selisih hidup **di dalam `JSONDATA`**, sebagai sub-pohon cermin `ValueDifference.*`.

Jadi tabel selisih adalah **benda baru**. Arahan pemilik proses menetapkannya, dan tidak ada bentuk
lama yang membatasinya.

### Yang menentukan: apakah daftar besaran yang boleh disesuaikan TERTUTUP

Diukur dari cermin `ValueDifference.*` — 161 jalur, di bawah **31 akar**. Sisanya, **sekitar 100
akar pohon utama, tidak punya cermin sama sekali**: seluruh field kepala yang bersifat uraian,
periode, mata uang, portofolio, akumulasi, pelaporan.

| Akar bercermin | Jalur |
|---|---|
| `Share` | 50 |
| `Limits` | 22 |
| `FacultativeShareList` | 21 |
| `LimitSummaryList` | 9 |
| `Installment` | 8 |
| `EGNPI` | 6 |
| 21 akar `Total*` (agregat — turunan, tidak disimpan) | 2 masing-masing |
| `BrokeragePercent`, `RNMShare`, `FacultativeShare`, `FacultativeShareBrokerage` | 1 masing-masing |

**Tertutup: ya.** Dan ia tertutup pada bentuk yang lebih tajam daripada sekadar "pendek":

> **Setiap besaran yang dapat disesuaikan adalah UANG atau PERSENTASE. Tidak ada satu pun tanggal,
> dan tidak ada satu pun teks.**

### Kenapa putusan KUNCI_PADANAN tidak otomatis berlaku di sini

Kami pernah menolak bentuk sempit — `SEAM-ADJUSTMENT.md` §3, rujukan polimorfik `JENIS_INDUK` +
`KUNCI_PADANAN` — dan menerimanya hanya dengan tiga syarat karena ia polimorfisme.

Putusan itu **tidak dibawa ke sini**, dan alasannya bukan kenyamanan:

| | KUNCI_PADANAN | Tabel selisih |
|---|---|---|
| Kedudukan | **data pokok** — bentuk baca atas nilai kontrak yang kanonik | **catatan** tentang dua nilai yang kanoniknya ada di tempat lain |
| Bila isinya salah | nilai kontrak terbaca salah | catatan salah; **nilai kontraknya tetap benar dan dapat dihitung ulang** |
| Ragam tipe yang harus ditampung | uang, tanggal, teks, persentase, enumerasi | **dua saja**: uang dan persentase |

Keberatan pokok terhadap bentuk sempit adalah **tipe tidak dapat ditegakkan, karena satu kolom
nilai harus memuat uang, tanggal, dan persentase sekaligus.** Keberatan itu **tidak berlaku di
sini**, karena himpunan yang harus ditampung hanya dua, dan keduanya bilangan eksak.

### "Tertutup" adalah fakta sistem lama, bukan putusan bisnis — penandaan

Kesimpulan "daftarnya tertutup" ditarik dari 31 akar yang punya cermin `ValueDifference.*`
sementara ~100 akar tidak. **Tetapi 31 itu mencerminkan apa yang pernah dibangun orang, bukan apa
yang boleh disesuaikan menurut bisnis.**

Seratus akar tanpa cermin punya dua tafsir yang sangat berbeda, dan ekspor tidak membedakannya:
ia memang tidak boleh disesuaikan, **atau** ia boleh disesuaikan dan sistem lama tidak pernah
mencatat selisihnya. **Diamnya sistem lama bukan bukti** — kalimat yang sudah dipakai berkali-kali
di sesi ini.

**Bentuk sempit justru yang tahan terhadap ketidakpastian ini**, dan itu memperkuat putusan G3,
bukan melemahkannya: daftar bertambah berarti **baris** bertambah di tabel acuan, bukan **kolom**
di tabel selisih. Tidak ada yang perlu dirancang ulang.

Dua hal yang dijaga sebagai gantinya:

1. **`BESARAN_DAPAT_DISESUAIKAN` diisi 31 sebagai AWAL, bukan sebagai batas.** Ketiadaan sebuah
   besaran di sana **tidak** diperlakukan sebagai larangan sampai dikonfirmasi bisnis.
2. **Ke-31 itu masuk daftar pertanyaan bisnis sebagai DAFTAR UNTUK DIBANTAH**, bukan sebagai
   pertanyaan terbuka. Bukan *"apa saja yang boleh disesuaikan?"* melainkan *"ini yang selama ini
   tercatat sebagai dapat disesuaikan — adakah yang kurang?"*

> Daftar untuk dibantah selalu mendapat jawaban lebih baik daripada pertanyaan terbuka.

### Keputusan G3

> **SEMPIT — satu baris per besaran yang berubah — dengan kosakata besaran TERTUTUP di tabel acuan.**

Dan tipenya dijaga meskipun kolom nilainya satu keluarga, lewat **tiga hal**:

1. **`ID_BESARAN` merujuk tabel acuan `BESARAN_DAPAT_DISESUAIKAN`** — bukan teks bebas. Sebuah
   besaran yang tidak ada di tabel acuan tidak dapat ditulis. Ini ADR-0038: himpunan yang bertambah
   tanpa mengubah arti disimpan sebagai data.
2. **Tabel acuan itu membawa `SATUAN_BESARAN`** — `UANG` atau `PERSENTASE`. Satu constraint
   memastikan baris selisih cocok dengan satuan besarannya: baris ber-satuan `UANG` **wajib**
   membawa mata uang dan paket uangnya; baris ber-satuan `PERSENTASE` **wajib** tidak membawanya.
3. **Kolom nilainya bilangan eksak**, bukan teks. Tidak ada titik di mana angka berubah menjadi
   teks dan kembali lagi — penyakit yang di sistem lama membuat presisi tidak dapat dipulihkan
   (ADR-0003).

Dengan ketiganya, bentuk sempit **tidak menjadi polimorfisme**: ia satu tipe nilai dengan satu
pembeda satuan yang berasal dari tabel acuan, bukan satu kolom yang berganti arti.

**Yang membatalkannya:** bila kelak ada besaran yang boleh disesuaikan dan ia **bukan** uang maupun
persentase — misalnya tanggal berlaku, atau kelas bisnis — maka keberatan polimorfisme kembali
berlaku dan bentuknya ditinjau ulang. Karena itu `SATUAN_BESARAN` ditulis sebagai himpunan tertutup
sejak awal: penambahan satuan ketiga **harus** memaksa perubahan yang terlihat, bukan diam-diam
memuat tanggal di kolom angka.

---

## Ringkasan ketiga gerbang

| Gerbang | Keputusan | Dasar |
|---|---|---|
| **G1** | Adjustment adalah **VERSI_KONTRAK**, bukan entitas baru; `PENYESUAIAN` **dibuang** | baris EDM memuat `JSONDATA` kontrak **utuh** + identitas sendiri + penunjuk pendahulu + persetujuan sendiri; nol fakta tersisa di luar model yang ada |
| **G1b** | Di sistem lama addendum **tidak menggantikan** — tetapi karena dua langkah penggantinya **dimatikan**, bukan karena dirancang begitu | `SaveTreatyIn_EDM_Act`: dua langkah ber-blok `//`, dan nol panggilan ke `SaveTreatyInOffer_Act` |
| **G2** | Penunjuknya **`ID_VERSI_KONTRAK`**, bukan `ID_KONTRAK` — **PERUBAHAN, bukan pelestarian** | `OLDID ← TreatyIn.ID`, dan *picker*-nya meng-UNION kontrak dengan addendum; tidak ada identitas versi di sistem lama |
| **G3** | Tabel selisih **SEMPIT**, kosakata besaran **tertutup** di tabel acuan, satuan `UANG`/`PERSENTASE` | 31 akar bercermin dari ~130; ~100 akar tanpa cermin; **nol** besaran bertipe tanggal atau teks |

**Tidak ada gerbang yang tertinggal terbuka.** Menggambar boleh dimulai setelah pemilik proses
memeriksa G1b — karena bila ia menilai temuan itu membalik G1, bentuk pokok ERD-nya berubah.

---

## Entitas baru yang lahir dari ketiga gerbang — hanya tiga

| Nama | Bita | Milik | Artinya |
|---|---|---|---|
| `NILAI_SELISIH` | 13 | Adjustment | satu baris per besaran yang berubah pada sebuah versi |
| `BESARAN_DAPAT_DISESUAIKAN` | 25 | **bersama** | tabel acuan: besaran apa yang boleh disesuaikan, dan satuannya |

**Dua, bukan tiga.** `PENYESUAIAN` dibuang di §G1c: seluruh 19 properti kelas addendum sudah
tertampung `KONTRAK` dan `VERSI_KONTRAK`, sehingga ia tidak membawa satu fakta pun.

Seluruh nama ≤ 30 bita.
