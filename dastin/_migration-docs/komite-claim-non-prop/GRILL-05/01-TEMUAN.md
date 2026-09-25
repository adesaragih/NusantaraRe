> Modul  : Komite Claim Non Prop · Ronde 05 · 2026-09-20
> Peran  : interogator
> Masukan: `KomiteRouter.xml` · `KomitePostAdjustment.xml` · `IsKomiteLoop.xml` · `KomiteTreaty_Flow.xml` · `ReinstatementPremiumDetails.xml` · `CreateChildKomiteCNP_Act.xml` · `CreateChildKomiteCloseNP_Act.xml` — seluruhnya ekspor Pega 2026-09-08/09
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

# 01 · TEMUAN RONDE 5

## 1. Tabel temuan

| ID | Klaim yang diuji | Bukti berkas·langkah | Putusan | Tkt | Dampak | Cara menutup |
|---|---|---|---|---|---|---|
| N-01 | Satu rule memakai satu gagasan "jenjang aktif" | `KomiteRouter.xml` langkah 1–4 (`.KomiteCount`) vs langkah 6.1 (`.KomiteAproval==0`, transisi `6/6`) | **SAHIH** | P2 | dua penentu jenjang aktif hidup berdampingan dan dapat berselisih | KETETAPAN — satu penentu tegas |
| N-02 | Jumlah jenjang bebas mengikuti roster | `KomiteRouter.xml` langkah 1–4 | **SAHIH** | P1 | empat nama keranjang ditulis di dalam rule; jenjang ke-5 tak punya sasaran cadangan | KETETAPAN — seleksi sebagai data (H-1) |
| N-03 | Cabang mati pada router (K-02) | `KomiteRouter.xml` langkah 7 vs 1–4 | **KODE MATI** | P2 | pada jalur bukan-'2', langkah 1–4 seluruhnya tertimpa | LARUT ke N-02 |
| N-04 | `TransferType` diisi kedua jalur pembuatan | `CreateChildKomiteCNP_Act.xml` langkah 9 · `CreateChildKomiteCloseNP_Act.xml` nol kemunculan | **SAHIH** | P1 | jalur tutup/tolak tidak pernah mengisinya; perutean satu-orang bukan kebetulan | KETETAPAN — sudah dijawab H-2 |
| N-05 | Satu penanda, satu ejaan | `CreateChildKomiteCloseNP_Act.xml` langkah 6 | **SAHIH** | P2 | `IsCloseFile` ditulis angka `1/0` pada anak dan teks `"1"/""` pada induk, dalam satu langkah | KETETAPAN — tipe tegas di model sasaran |
| N-06 | Roster dibangun ulang tiap kali | `CreateChildKomiteCNP_Act.xml` langkah 26.1 & 26.8.1 | **SAHIH** | P1 | pada jalur bersyarat roster tidak dibersihkan lalu di-`APPEND` | KETETAPAN — pembentukan roster idempoten |
| N-07 | Pesan galat menghentikan pembuatan sirkulasi | `CreateChildKomiteCNP_Act.xml` langkah 29→31 · `CreateChildKomiteCloseNP_Act.xml` langkah 10→12 | **SAHIH** | **P0** | penjaga hanya melewati satu langkah; langkah sesudahnya menulis nomor komite dari case lain lalu menyimpan | KETETAPAN — H-5 diperluas |
| N-08 | Kegagalan pembentukan terlihat oleh pengguna | `CreateChildKomiteCNP_Act.xml` langkah 28 → 36; `IsFlagError` nol pembaca di 338 berkas | **SAHIH** | P1 | bila `SpreadingRisk` kosong, activity keluar diam-diam dengan penanda yang tak dibaca siapa pun | KETETAPAN — kegagalan yang terlihat |
| N-09 | Ambang `Flagkomite` punya tiga kelas yang hidup | `CreateChildKomiteCNP_Act.xml` langkah 10, 14, 26.3–26.6 | **SAHIH** | P1 | `>30 && <=30` mustahil benar; `Flagkomite=0` tidak menetapkan ambang sama sekali | KETETAPAN — H-1; batas bawah menunggu pemilik proses |
| N-10 | Klaim bersyarat memakai ambang kelasnya | `CreateChildKomiteCNP_Act.xml` langkah 26.5 lalu 26.6 | **KEPUTUSAN BISNIS** | P2 | bersyarat selalu memakai roster terluas, menimpa kelas 25 juta | pemilik proses — digabung ke pertanyaan H-1 |
| N-11 | `KomiteAproval` selalu numerik | `CreateChildKomiteCNP_Act.xml` langkah 26.9/26.10 vs `KomiteRouter.xml` langkah 6.1 | **RAGU** | P1 | `""` ditulis, `==0` diuji | LARUT — H-7 |
| N-12 | Deskripsi langkah menyatakan apa yang dilakukan langkah | `KomiteRouter.xml` langkah 5 | **SAHIH** | P3 | label menjanjikan penambahan hitungan yang tidak ada di dalamnya | LARUT — aturan "deskripsi bukan aturan" |
| N-13 | `KomitePostAdjustmentCWP` punya pintu masuk lain (Q-6) | `KomitePostAdjustment.xml` langkah 31; satu-satunya `Call` di 338 berkas | **TAK TERVERIFIKASI** | P2 | nol pemanggil lain **di dalam ekspor**; ruleset lain tidak diekspor | BUKTI — sudah diambil, batas dinyatakan |
| N-14 | `komiteAccept_ticket` tidak terhubung ke flow (Q-7) | `KomiteTreaty_Flow.xml`, `Decision1.pyTicketShapes` | **SALAH** | P2 | ticket menempel pada gerbang `KomiteLoop`; yang tidak ada adalah penaiknya | BUKTI — sebagian; penaik di luar ekspor |
| N-15 | Empat nilai `FlagProrate` punya arti yang tertulis (Q-3) | `ReinstatementPremiumDetails.xml` baris 609, 2540, 4471, 6687 | **TAK TERVERIFIKASI** | P2 | nilai 0/1/2/3 memilih satu dari empat blok tampilan; artinya tidak tertulis di ekspor | BUKTI — batas dinyatakan (Field Value tak diekspor) |

---

## 2. Rincian

### N-01 · Dua gagasan "jenjang aktif" dalam satu rule

`KomiteRouter.xml` memakai dua penentu yang berbeda untuk perkara yang sama.

**Penentu pertama — hitungan.** Langkah 1–4 memilih sasaran dari `.KomiteCount`: nilai
`1`, `2`, `3`, `4` masing-masing memberi nama keranjang tetap.

**Penentu kedua — keadaan baris.** Langkah 6 beriterasi atas `.KomiteList` dengan pra-syarat
`Primary.TransferType=='2'`; sub-langkah tunggalnya berpra-syarat `.KomiteAproval==0`
(`pyStepsPreCondParamsWhenFalse=3`, yakni lewati langkah) dan mengerjakan
`param.AssignTo = .KomiteID`. Transisi sesudahnya dibaca mentah dari XML:
`pyStepsTransParamsWhen=true`, `pyStepsTransParamsWhenTrue=6`, `pyStepsTransParamsWhenFalse=6`.
Sandi `6` berarti **keluar iterasi**. Karena baris yang tidak cocok dilewati tanpa menjalankan
transisinya, iterasi berhenti pada **kecocokan pertama** — baris pertama yang belum
memutuskan.

Kedua penentu tidak pernah saling memeriksa. `.KomiteCount` bertambah satu setiap kali
`KomitePostAdjustment.xml` langkah 29 dijalankan, sedangkan `.KomiteAproval` ditulis per
baris. Bila satu jenjang pernah memutuskan di luar urutan — atau bila roster berubah di
tengah jalan (N-06) — keduanya menunjuk jenjang yang berlainan, dan yang menang adalah
langkah 6 karena berjalan belakangan.

> **Akibat yang ikut tertarik:** perlakuan yang ditulis `PUTUSAN-01.md` §8 untuk K-02 —
> "jenjang aktif = baris pertama yang belum memutuskan" — **berbukti benar**, dan berhenti
> menjadi tafsir. Yang tersisa bukan urutannya, melainkan bahwa `.KomiteCount` tetap hidup
> sebagai penentu kedua; model sasaran harus punya satu penentu saja, dan hitungan jenjang
> menjadi turunannya, bukan saingannya.

### N-02 · Empat nama keranjang di dalam rule

`KomiteRouter.xml` langkah 1–4 masing-masing berpra-syarat `.KomiteCount==1`, `==2`,
`==3`, `==4` dan menetapkan `param.AssignTo` ke tetapan teks `"komitepnc"`, `"komitepnc2"`,
`"komitepnc3"`, `"komitepnc4"`.

`KomiteCount` lahir bernilai `1` pada kedua pembuat anak, bertambah satu di
`KomitePostAdjustment.xml` langkah 29, dan dibandingkan dengan `KomiteLoop` oleh
`When/IsKomiteLoop.xml`, yang berbunyi `.AcceptStatus = "1"` **dan** `.KomiteCount <= .KomiteLoop`.
`KomiteLoop` sendiri adalah jumlah baris laporan roster (`FAKTA-A1B.md` F-13), yang tidak
dibatasi empat.

> **Akibat yang ikut tertarik:** pada sirkulasi dengan lebih dari empat jenjang, langkah
> 1–4 tidak memberi sasaran cadangan apa pun. Tabel roster pada model sasaran karena itu
> harus menyimpan sasaran penugasan sebagai **data per jenjang**, bukan mengandalkan
> daftar tetap di dalam kode.

### N-03 · Yang mati bukan cabangnya, melainkan langkah 1–4

Langkah 7 berpra-syarat `Primary.TransferType!='2'` dengan `pyStepsPreCondition=true`, dan
isinya `param.AssignTo = .Komite.KomiteID` — tanpa syarat di dalam. Langkah 6 dan 7 karena
itu saling melengkapi: satu di antaranya **selalu** berjalan.

Pada jalur `!='2'`, langkah 7 selalu menimpa apa pun yang ditulis langkah 1–4. Pada jalur
`=='2'`, nilai langkah 1–4 hanya bertahan bila tak satu pun baris roster ber-`KomiteAproval==0`.

> **Akibat yang ikut tertarik:** K-02 pada `PENGETAHUAN.md` §11 berbunyi "routing bercabang
> yang cabangnya mati"; yang mati sebenarnya adalah empat langkah keranjang. Itu mengubah
> perlakuannya — bukan menghidupkan cabang, melainkan menghapus empat langkah.

### N-04 · Satu penulis `TransferType`, satu jalur tanpa penulis

`Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` langkah 9 menetapkan
`ChildWorkPage.TransferType = .Type` — satu-satunya penulis di seluruh 338 berkas.
`CreateChildKomiteCloseNP_Act.xml` tidak menyebut `TransferType` sama sekali.

Maka sirkulasi yang lahir dari jalur tutup/tolak **selalu** memenuhi `TransferType!='2'`
dan selalu diarahkan lewat `KomiteRouter` langkah 7, yaitu ke satu halaman `.Komite`
tunggal. Itu sejalan dengan roster satu orang yang dipaku di `CreateChildKomiteCloseNP_Act.xml`
langkah 8, dan sejalan dengan ketetapan H-2.

Nilai `.Type` yang membuat jalur banyak-jenjang hidup adalah `'2'`. Nilai lain yang
terbukti dipakai di tempat lain adalah `'3'` (`Claim Non Prop/Harness/KomiteCNP.xml` baris
3577 dan 3963). Arti keduanya tidak tertulis di ekspor — `Rule-Obj-FieldValue` tidak
diekspor (`INVENTARIS-BUKTI.md` §1).

> **Akibat yang ikut tertarik:** satu akibat lagi mendarat di isi surat dan **tertahan
> pagar A-6** — dicatat sebagai lubang bertahan `L5-1` di `03-LUBANG.md`, tanpa perlakuan
> dan tanpa alasan yang menembus pagar.

### N-05 · Satu penanda, dua ejaan, satu langkah

`CreateChildKomiteCloseNP_Act.xml` langkah 6 memuat, berurutan:

- `ChildWorkPage.IsCloseFile = @if(TempCommiteClaim.TypeComentAnalysis!="5",1,0)`
- `ChildWorkPage.IsReject   = @if(TempCommiteClaim.TypeComentAnalysis=="5","1","")`
- `.IsCloseFile             = @if(TempCommiteClaim.TypeComentAnalysis!="5","1","")`
- `.IsReject                = @if(TempCommiteClaim.TypeComentAnalysis=="5","1","")`

Anak menerima `0` pada penolakan; induk menerima kosong. Gerbang jalur B di
`KomitePostAdjustment.xml` langkah 1 menguji `IsCloseFile=="1" || IsReject=="1"`, sehingga
penolakan tetap lolos lewat `IsReject`.

> **Akibat yang ikut tertarik:** setiap penguji yang berbunyi "tidak kosong" akan membaca
> `"0"` sebagai benar. Kolom pendaratan `IsCloseFile` dan `IsReject` yang disebut G-12
> karena itu tidak boleh dipetakan sebagai teks bebas; keduanya satu boolean.

### N-06 · Roster tidak dibersihkan pada jalur bersyarat

`CreateChildKomiteCNP_Act.xml` langkah 26 sub-langkah 1 adalah `Property-Remove` atas
`.ComiteeClaim`, berpra-syarat `.IsSubjectivity==true` dengan aksi **lewati-step**. Artinya
pembersihan berjalan pada klaim biasa dan **tidak** berjalan pada klaim bersyarat.

Sub-langkah 26.8.1 kemudian menulis `Primary.ComiteeClaim(<APPEND>).KomiteID = .OPERATOR_ID`
untuk tiap baris hasil laporan. Pada klaim bersyarat yang diproses lebih dari sekali, baris
lama bertahan dan baris baru menumpuk di belakangnya.

Sub-langkah 26.13 dan 26.14 menetapkan
`ChildWorkPage.KomiteLoop = @Utilities.SizeOfPropertyList(pyReportContentPage.pxResults)`
— jumlah baris **laporan**, bukan jumlah baris roster yang terbentuk.

> **Akibat yang ikut tertarik:** sejak penumpukan pertama, `KomiteLoop` lebih kecil daripada
> jumlah anggota `KomiteList`. Gerbang `IsKomiteLoop` (`.KomiteCount <= .KomiteLoop`) karena
> itu berhenti sebelum seluruh anggota memutuskan, dan sisa anggota tidak pernah menerima
> giliran. Pembentukan roster pada sistem baru harus idempoten: menyusun ulang, bukan
> menambahkan.

### N-07 · Penjaga pesan galat hanya melewati satu langkah

Pada `CreateChildKomiteCNP_Act.xml`, langkah 29 `Call pxAddChildWork` berpra-syarat
`@hasMessages(myStepPage)` dengan aksi **lewati-step**. Langkah 30 (`Obj-Refresh-And-Lock`)
dan langkah 31 tetap berjalan, dan langkah 31 menulis:

```
.KomiteNo        = @substring(pyWorkPage.pxCoveredInsKeys(<LAST>),19,30)
Param.NoKomite   = @substring(pyWorkPage.pxCoveredInsKeys(<LAST>),19,30)
```

Pola yang sama ada di `CreateChildKomiteCloseNP_Act.xml`: langkah 10 dijaga
`@hasMessages`, langkah 11–12 tetap berjalan dan menulis `.KomiteNo` dari
`pxCoveredInsKeys(<LAST>)`, langkah 13 `Obj-Save pyWorkPage`, langkah 14 mengirim surat.

`<LAST>` menunjuk case tercakup terakhir. Bila pembuatan dilewati, yang terbaca adalah
sirkulasi **sebelumnya** — atau, bila belum pernah ada, kosong.

Uji kenyataan: ini jalur galat, sehingga sistem tetap dapat berjalan bertahun-tahun tanpa
seorang pun menemuinya. Yang tidak dapat disimpulkan dari berkas adalah seberapa sering
`@hasMessages` benar di titik itu — dan itu tidak diperlukan untuk menetapkan perlakuannya.

> **Akibat yang ikut tertarik:** klaim dapat menyimpan nomor komite milik sirkulasi lain,
> lalu surat dikirim atas nomor itu. H-5 sudah menetapkan validasi yang **menolak**, bukan
> menandai; temuan ini memperluasnya: yang harus menolak bukan hanya pembuatan, melainkan
> seluruh urutan yang bersandar pada keberhasilannya.

### N-08 · Keluar diam-diam ketika `SpreadingRisk` kosong

`CreateChildKomiteCNP_Act.xml` langkah 28 berpra-syarat
`@LengthOfPageList(.SpreadingRisk)=0 || @LengthOfPageList(ChildWorkPage.Adjustment.SpreadingRisk)=0`
dengan aksi **lompat ke label `END`**. Label `END` adalah langkah 36, yang isinya
`pyWorkPage.ClaimData.IsFlagError = "true"`, disusul langkah 37 `Obj-Save pyWorkPage`.

`IsFlagError` muncul di **satu** berkas saja dari 338 — berkas ini sendiri. Tidak ada
Section, Harness, When, atau Activity lain yang membacanya.

> **Akibat yang ikut tertarik:** pengguna menekan tombol, tidak ada sirkulasi yang lahir,
> tidak ada pesan, dan klaim tetap tersimpan. Setiap jalur yang berakhir tanpa hasil pada
> sistem baru wajib berakhir dengan kegagalan yang terlihat — ini kasus uji paritas yang
> dirancang gagal.

### N-09 · Cabang ambang yang mustahil, dan kelas tanpa ambang

`CreateChildKomiteCNP_Act.xml` langkah 10 menetapkan empat tetapan:

```
Local.LimitMax               = 30000000.00
Local.LimitMaxDivHead        = 50000000.00
Local.LimitPersenMax         = 30.00
Local.LimitPersenMaxDivHead  = 30.00
Local.Flagkomite             = 0
```

Langkah 12 menetapkan `Flagkomite=1` bila `CekLimitPersen<=30.00` **dan** `TotalValueAdjust<=30 juta`.
Langkah 13 menetapkan `Flagkomite=2` bila `CekLimitPersen<=30.00` **dan** nilai berada di
antara 30 juta dan 50 juta. Langkah 14 menetapkan `Flagkomite=2` bila
`CekLimitPersen>Local.LimitPersenMax && CekLimitPersen<=Local.LimitPersenMaxDivHead` —
karena kedua tetapan bernilai `30.00`, syarat itu berbunyi `>30 && <=30` dan **tidak pernah
benar**. Deskripsi langkah itu berbunyi "LimitPersenMaxDivHead, CekLimitPersen >15 && <=30";
angka 15 tidak ada di kode mana pun.

`CekLimitPersen` diisi dari `pyWorkPage.TreatyInMaster.RNMShare` (langkah 10), yaitu
persentase bagian — bukan nilai uang.

Ambangnya sendiri dipakai di langkah 26: 26.3 menetapkan parameter laporan **tanpa**
`LIMIT_BOTTOM` dan berkomentar "Untuk sementara dihilangkan"; 26.4 (`Flagkomite==1`)
menetapkan `LIMIT_BOTTOM = 0`; 26.5 (`Flagkomite==2`) menetapkan `25000001`; 26.6
(`Subjectivity==true`) menetapkan `0`.

> **Akibat yang ikut tertarik:** ketika `RNMShare > 30`, `Flagkomite` tetap `0`, tidak satu
> pun dari 26.4/26.5/26.6 berjalan, dan `Param.LIMIT_BOTTOM` **tidak pernah diisi** sebelum
> `pxRetrieveReportData` dipanggil di 26.7. Penyaring `FilterEmailKomiteWithLimit` (F-14)
> membandingkan `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM` terhadap parameter kosong. Di situlah
> sirkulasi tanpa jenjang (F-22, G-17, H-6) lahir — asal-usulnya kini berbukti langkah,
> bukan dugaan.

### N-10 · Bersyarat menimpa kelas 25 juta

Sub-langkah 26.5 dan 26.6 berjalan berurutan dan keduanya menulis `Param.LIMIT_BOTTOM`.
Karena 26.6 berada sesudahnya, klaim bersyarat bernilai di atas 25 juta berakhir dengan
`LIMIT_BOTTOM = 0`, yaitu roster **terluas**, bukan roster kelasnya.

Uji masuk akal bisnis: ini dapat memang disengaja — klaim bersyarat mungkin sengaja
diedarkan lebih luas. Karena itu putusannya `KEPUTUSAN BISNIS`, bukan cacat.

> **Akibat yang ikut tertarik:** aturan seleksi jenjang sebagai data (H-1) harus memuat
> dimensi "bersyarat" di samping dimensi nilai, atau keputusan ini hilang dalam migrasi
> tanpa seorang pun menyadarinya.

### N-11 · `""` ditulis, `==0` diuji

Sub-langkah 26.9 dan 26.10 menulis `KomiteAproval = ""` untuk jalur
`pyWorkPage.pxCreateOperator=="VINCENTVERNANDO_1"`, sedangkan `KomiteRouter` langkah 6.1
menyeleksi dengan `.KomiteAproval==0`.

Apakah Pega menyamakan `""` dengan `0` pada perbandingan itu tidak dapat diputuskan dari
berkas. Putusan diturunkan ke `RAGU` sesuai aturan berhenti nomor 3 (`00-LINGKUP.md` §6).

> **Akibat yang ikut tertarik:** tidak ada — H-7 sudah menetapkan jalur
> `VINCENTVERNANDO_1` tidak dimigrasikan, sehingga perkara ini larut tanpa perlakuan.
> Yang tersisa hanya pelajarannya: nilai kosong dan nol tidak boleh berbagi satu kolom
> pada model sasaran.

### N-12 · Label langkah menjanjikan yang tidak dikerjakannya

`KomiteRouter.xml` langkah 5 berlabel "KomiteCount Increment & AcceptStatus Reset". Isinya
satu baris: `set .AcceptStatus = ""`. Penambahan `KomiteCount` sebenarnya terjadi di
`KomitePostAdjustment.xml` langkah 29.

> **Akibat yang ikut tertarik:** ini kejadian kedua yang tercatat di mana deskripsi langkah
> menyatakan aturan yang tidak ada di kodenya — yang pertama adalah angka 15 pada N-09.
> Larangan "jangan menaikkan komentar atau deskripsi langkah menjadi aturan" karena itu
> bukan kehati-hatian, melainkan aturan yang sudah dua kali menyelamatkan pembacaan.

### N-13 · `KomitePostAdjustmentCWP` — satu pintu masuk di dalam ekspor (Q-6)

Pencarian teks atas seluruh 338 berkas menemukan `KomitePostAdjustmentCWP` sebagai:
satu `pyStepsActivityName` berbunyi `Call KomitePostAdjustmentCWP` di
`KomitePostAdjustment.xml` (langkah 31, label `CLS`), satu baris indeks rujukan
`Embed-Reference-Rule` di berkas yang sama, dan definisi rule-nya sendiri.

Tidak ada pemanggil lain — **di dalam ekspor**. Ruleset lain tidak diekspor, sehingga
putusannya `TAK TERVERIFIKASI`, bukan `SAHIH`.

> **Akibat yang ikut tertarik:** jalur B tidak perlu dimodelkan sebagai kasus-guna
> berdiri sendiri; ia satu cabang di dalam kasus-guna keputusan. Batasnya dinyatakan,
> bukan ditutupi: bila kelak ditemukan pemanggil di ruleset lain, ketetapan ini gugur.

### N-14 · Ticket itu terhubung — yang hilang adalah penaiknya (Q-7)

`PENGETAHUAN.md` §5 menyatakan `komiteAccept_ticket` "tidak dirujuk transisi mana pun di
flow ini". Pembacaan ulang `KomiteTreaty_Flow.xml` menunjukkan sesuatu yang lain: ticket
tidak pernah dirujuk oleh transisi karena ticket memang bukan transisi. Ia muncul di
`pyTicketShapes` milik shape `Decision1` — gerbang `KomiteLoop`, kelas
`Data-MO-Gateway-DataXOR`. Dua shape lain yang punya `pyTicketShapes` (`END52` dan
`ASSIGNMENT63`) memuat baris kosong tanpa `pyMOName`.

Artinya: menaikkan `komiteAccept_ticket` memindahkan case ke gerbang lingkar jenjang.
Tidak satu pun dari 338 berkas menaikkannya — pencarian `Obj-Set-Tickets`, `SetTicket`, dan
`pyTicketName` tidak menghasilkan apa pun.

> **Akibat yang ikut tertarik:** ada jalan masuk ke gerbang lingkar yang melewati
> penugasan, dan pemiliknya di luar ekspor. Model sasaran karena itu tidak boleh
> mengandaikan bahwa satu-satunya cara maju ke jenjang berikutnya adalah lewat keputusan
> di layar.

### N-15 · Empat nilai `FlagProrate`, empat blok tampilan (Q-3)

`Section/ReinstatementPremiumDetails.xml` memuat empat wadah dengan
`pyContainerVisibleWhen` berbunyi `.FlagProrate = 0`, `= 1`, `= 2`, `= 3` (baris 609,
2540, 4471, 6687), masing-masing berjudul tingkat `h2`. Berkas yang sama ada di kedua
folder.

Yang terbukti: **empat** nilai dipakai, dan keempatnya memilih tampilan rumus
reinstatement yang berbeda. Yang tidak terbukti: apa arti bisnis tiap nilai, dan siapa
yang menulisnya — `Rule-Obj-FieldValue` tidak ikut diekspor.

> **Akibat yang ikut tertarik:** `FlagProrate` bukan penanda dua-nilai; kolom sasarannya
> harus menampung empat nilai berkode, dan rumus reinstatement sisi Claim
> (`FINDING-00*` sisi Claim) punya empat varian tampilan yang harus dicocokkan sebelum
> rumusnya dipindahkan.

---

## 3. Tiga uji wajib atas temuan sendiri

**Uji lingkar.** Satu-satunya temuan yang berdiri di atas ketiadaan adalah N-13
(nol pemanggil lain) dan N-14 (nol penaik ticket). Keduanya menyangkut **pemanggilan
antar-rule**, yang justru termasuk jenis data yang G-01 nyatakan tidak terekspor untuk
parameter metode — tetapi nama activity yang dipanggil **memang** terekspor
(`pyStepsActivityName`), sehingga ketiadaan di sini bermakna. Meski begitu keduanya
ditandai `TAK TERVERIFIKASI` dan `SALAH sebagian`, bukan `SAHIH`, dan batasnya dinyatakan:
nol **di dalam ekspor**.

**Uji kenyataan.** Tiga temuan menyiratkan kerusakan yang seharusnya terlihat: N-06
(penumpukan roster), N-07 (nomor komite dari case lain), N-08 (keluar diam-diam). Ketiganya
hanya bekerja pada jalur bersyarat atau jalur galat, yang jarang; tak satu pun menyiratkan
jalur utama rusak. Karena itu tak satu pun diturunkan.

Uji ini juga menjatuhkan satu pembacaan awal. Bacaan pertama atas `KomiteRouter` langkah 6
menyimpulkan bahwa giliran jatuh pada baris **terakhir** yang belum memutuskan — yang
berarti urutan jenjang terbalik terhadap derajat pada setiap sirkulasi yang pernah berjalan.
Itu tidak lolos uji kenyataan: pembalikan urutan wewenang pada jalur utama akan terlihat
oleh penggunanya sejak hari pertama. Pemeriksaan ulang atas XML mentah menunjukkan
bacaan itu keliru — transisi `2/2` yang terbaca ternyata milik langkah 7, sedangkan langkah
6.1 bertransisi `6/6`. Rinciannya di `06-PUTUSAN.md` M5-01.

**Uji masuk akal bisnis.** N-10 semula terbaca sebagai cacat urutan langkah; pembacaan itu
akan menyiratkan perusahaan tidak sengaja memperluas komite untuk klaim bersyarat selama
bertahun-tahun. Lebih masuk akal bahwa itu memang kebijakan. Putusannya diubah menjadi
`KEPUTUSAN BISNIS` sebelum dikirim. N-09 tidak lolos uji yang sama dengan hasil berbeda:
`>30 && <=30` tidak dapat dibaca sebagai kebijakan apa pun — itu salah tulis.
