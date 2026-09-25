# CONTEXT — Treaty In

Glosarium dan batas konteks modul Treaty In. Berkas ini **bukan** spesifikasi dan bukan tempat
keputusan implementasi. Keputusan ada di `docs/adr/`; temuan ada di `PENGETAHUAN.md` dan
`RINGKASAN-GRILLING.md`.

Disusun 23 September 2026 dari sesi penggalian kebutuhan atas ekspor Pega `D:\XML_NURE\Treaty In`.

---

## 1. Batas lingkup yang eksplisit

### 1.1 Treaty In Adjustment — masuk lingkup migrasi

Folder `D:\XML_NURE\Treaty In Adjustment` **masuk lingkup migrasi per 23 September 2026**; dasar
dan akibatnya dicatat di `../treaty-in-adjustment/KEPUTUSAN-GRILLING-ADJUSTMENT.md` **GRL-01**.
Embargo yang berlaku sebelumnya tidak pernah ditetapkan oleh ADR, sehingga berakhirnya tidak
menuntut revisi ADR.

Pemahaman AS-IS modulnya ada di `../treaty-in-adjustment/PENGETAHUAN.md`, dan status ADR-0048
**sedang direkonsiliasi** di grilling Adjustment (cabang A).

### 1.2 Yang dikeluarkan dari lingkup migrasi

| Dikeluarkan | Alasan |
|---|---|
| Tombol pengembang (`TestCopyDifference`, sejenisnya) | alat bantu pengembangan, bukan fungsi bisnis |
| Perlakuan khusus kontrak `ID == "1000951"` | diputuskan pemilik proses 22 Sep 2026: tidak dibawa |
| Penulisan selisih fac share ke `ActualValue.*` | diputuskan pemilik proses 22 Sep 2026: tidak dibawa |
| Area C — perilaku layar dinamis | direkonstruksi dari Section saat menulis spec, bukan lewat wawancara |
| Area F sisi masuk (penyimpanan objek, pembacaan master) | tertutup oleh keputusan layanan dokumen bersama dan keputusan adapter |

---

## 2. Aturan kerja — cara membaca sistem lama

### 2.0 Prinsip di atas semuanya: STRUKTUR YANG TERLIHAT BUKAN STRUKTUR YANG BERLAKU

Empat aturan di bawah lahir satu per satu dari empat kesalahan yang berbeda. Keempatnya ternyata
**satu bentuk**, dan pemilik proses meminta bentuk itu ditulis di atasnya supaya yang kelima
dikenali lebih cepat oleh orang yang memegang prinsipnya daripada oleh orang yang menghafal empat
kasus.

> **Di ekspor Pega, segala sesuatu yang tampak otoritatif belum tentu mengikat. Yang mengikat adalah
> apa yang benar-benar dieksekusi.**

| Yang tampak otoritatif | Ternyata | Aturan |
|---|---|---|
| sebuah langkah ada di daftar | bloknya bisa mati (`//`) dan tidak pernah berjalan | §2.1 |
| label dan deskripsi langkah | bisa bertentangan dengan kondisinya, dan kondisinya yang menang | §2.2 |
| kondisi yang tersimpan di sebelah sebuah kontrol | bisa diabaikan seluruhnya oleh `pyVisible = ALWAYS` | §2.2 |
| nama berkas tempat properti ditemukan | berkas Section membundel aturan lain di dalamnya | §2.8 |

**Cara memakainya:** setiap kali sebuah kesimpulan bersandar pada *keberadaan* sesuatu — sebuah
langkah, sebuah label, sebuah kondisi, sebuah nama berkas — tanyakan apa yang membuat benda itu
**berlaku**, dan periksa itu lebih dulu. Keberadaan bukan keberlakuan.

### 2.0-f Turunan KEENAM: kode yang tampak umum tetapi BERKONDISI PADA SATU BARIS DATA

Kelima turunan §2.0 sebelumnya tentang kode yang **mati**, **tak terjangkau**, atau **berlabel
salah**. Yang keenam berbeda, dan **cara menemukannya juga berbeda** — itu sebabnya ia diberi
tempatnya sendiri.

> **Kode yang tampak umum tetapi berkondisi pada satu baris data tertentu.**
> Ia **tidak mati**, **tidak tersembunyi**, dan **labelnya bahkan jujur**. Yang membuatnya tidak
> terlihat adalah bahwa ia hanya berjalan pada **dua dari sekian ribu kontrak** — sehingga setiap
> pengamatan atas kontrak lain **benar, dan tetap tidak lengkap**.

Bentuknya di Treaty In, ditemukan 24 September 2026:

```
pyStepsPreCondParamsWhen = TreatyIn.ID=="1000951"     14 kemunculan
pyStepsPreCondParamsWhen = TreatyIn.ID=="1000069"      2 kemunculan
pyStepsBlockName         = (kosong — hidup)
```

**Cara menemukannya bukan memeriksa blok dan bukan menyapu pemanggil, melainkan MENYAPU TETAPAN** —
nilai literal yang ditulis ke properti bisnis, dan nilai literal di dalam kondisi.
Perkakasnya `alat/sapu-tetapan-di-kode.py`, dan ia satu keluarga dengan pencari "cacat di balik
nilai bawaan" (§2.5): keduanya mencari **sesuatu yang benar untuk hampir semua hal, dan karena itu
tidak pernah diuji pada sisanya**.

**Akibatnya menjalar ke uji data**, dan itu harus ditangani saat ujinya dirancang, bukan saat
hasilnya dibaca: setiap uji yang menghitung atas seluruh baris **mengeluarkan baris istimewa itu
dari hitungan pokoknya dan melaporkannya terpisah** — lihat kepala
`PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql`.

### 2.0-g Turunan KETUJUH: nilai keadaan yang TIDAK PUNYA PENULIS

Ditemukan 24 September 2026, dan **cara menemukannya berbeda dari keenam sebelumnya** — itu sebabnya
ia diberi tempatnya sendiri.

> Sebuah **nilai keadaan** diuji berkali-kali sebagai kondisi, tetapi **tidak pernah disetel** oleh
> aturan mana pun. Ia **bukan** langkah mati dan **bukan** kondisi yang tidak pernah benar — ia
> keadaan yang **tidak punya jalan masuk**.

Bentuknya di Treaty In: jabatan `ReasTreatyInGroupLeader` diuji **delapan kali**, disetel **nol
kali** di kedua modul, dan cabangnya di `DataTransform/Akseptasi_DT.xml` ber-`pyDisabled = true`.

**Cara menemukannya: sapu nilai yang DIBACA, lalu cari penulisnya — bukan sebaliknya.** Keenam
turunan sebelumnya berangkat dari **kode** (blok mati, pemanggil, label, tetapan); yang ini berangkat
dari **nilai**. Kumpulkan setiap literal yang dibandingkan terhadap sebuah properti keadaan, lalu
periksa mana yang tidak pernah muncul sebagai nilai yang **ditulis**.

**Dan jangan berhenti pada temuan.** Nilai tanpa penulis punya dua bacaan: ia **sisa mati**, atau
**penulisnya ada di tempat yang tidak ikut terekspor** — lihat §2.0-h. Yang memisahkan biasanya
**data**, dan pertanyaannya bukan *apakah* nilai itu pernah ada melainkan **kapan**.

### 2.0-h "TIDAK ADA DI EKSPOR" BUKAN "TIDAK ADA DI SISTEM"

> **ADA di ekspor → ada di sistem. Sah.**
> **TIDAK ADA di ekspor → tidak ada di sistem. TIDAK sah.**

Ini L-8 satu tingkat lebih tinggi. Di sana sebuah **perkakas** tidak pernah menunjukkan 414 properti.
Di sini **satu jenis aturan penuh** tidak muncul sama sekali — lebih mudah diperhatikan, **tetapi
hanya bila ada yang bertanya.**

**Langkah yang ditambahkan, dan ia tidak menunggu siapa pun:** bandingkan jenis aturan yang muncul
di ekspor terhadap jenis yang diketahui dipakai Pega. Hasil sapuan 24 September 2026 —
**21 jenis nol kemunculan** — ada di `4-erd-dan-tabel-datar/JENIS-ATURAN-TAK-TEREKSPOR.md`, dan lima
di antaranya mengubah klaim yang sudah kita buat.

Dua jenis yang paling berbahaya bila nol, dan sebabnya khusus:

> **`Rule-Declare-Expressions` dan `Rule-Declare-Trigger` MENULIS NILAI tanpa dipanggil siapa pun.**
> Sapuan penulis yang mencari `Property-Set` di aktivitas dan *data transform* **tidak akan pernah
> melihatnya.** Selama keduanya nol di ekspor, setiap kalimat *"properti ini tidak pernah ditulis"*
> punya lubang yang belum ditutup.

### 2.0-i Sapuan "siapa menulis X" mendaftar BENTUK penulisnya lebih dulu

Ditemukan 24 September 2026, ketika klaim *"kedua properti ini tidak pernah ditulis satu aturan
pun"* terbukti salah — dan sebabnya **bukan** jenis aturan yang tidak terekspor (§2.0-h), melainkan
**bentuk penulisan yang tidak pernah dicari**.

Sapuan kami hanya mengenali `<PropertiesName>`, bentuk `Property-Set`. Yang menulis kedua properti
itu berbentuk lain:

```xml
<pyDisplayProperty>.StatusActive</pyDisplayProperty>
<pyPropertyTarget>TreatyIn.CedingStatusActive</pyPropertyTarget>
<pySetValueOnSelect>true</pySetValueOnSelect>
```

**Pemetaan hasil pencarian di layar** — memilih satu baris master menyalin kolomnya ke properti
sasaran.

> **Sebelum menyapu "siapa menulis X", daftarkan dulu BENTUK-BENTUK penulisan yang mungkin:**
> penugasan langsung (`Property-Set`), **pemetaan hasil pencarian** (`pyPropertyTarget`), pemetaan
> koneksi basis data, pemetaan layanan, **ekspresi terdeklarasi**, **pemicu terdeklarasi**.
>
> Sapuan satu bentuk menjawab pertanyaan yang lebih sempit daripada yang ditanyakan, **dan diamnya
> terbaca sebagai ketiadaan.**

Di Treaty In bentuk ini ternyata **sempit** — empat sasaran berawalan `TreatyIn.` di seluruh korpus —
sehingga ia tidak meruntuhkan adjudikasi masukan-versus-turunan. Tetapi **sempitnya baru diketahui
setelah disapu**, dan sebelum itu ia lubang yang tidak berbatas.

### 2.0-j Sebuah sapuan tidak boleh bersandar pada SATU NAMA TAG — dan daftarnya disusun mesin

§2.0-i menyuruh **mendaftar** bentuk penulisnya. Menjalankannya, pada hari yang sama, memperlihatkan
bahwa daftar yang disusun **dari yang terpikir** tetap kurang. Daftar yang benar disusun **mekanis**:

> **Keluarkan seluruh nama elemen XML yang muncul di korpus, beserta jumlahnya DAN jenis aturan
> tempat ia muncul.** Lalu saring daftar itu dengan pertanyaan yang sedang ditanyakan.

Di korpus ini: **1.823 nama elemen, 174 berbentuk penugasan, empat yang benar-benar menulis properti
kontrak** — dan **hanya satu yang pernah kami sapu**. Yang terbesar terlewat karena beda satu
awalan: `<PropertiesName>` di aktivitas, `<pyPropertiesName>` di *data transform*.

> **ATURAN: setiap sapuan mencatat NAMA TAG YANG DIPAKAINYA di dalam dirinya sendiri, dan nama itu
> diperiksa terhadap daftar nama elemen sebelum hasilnya dipercaya.**

Pemeriksaan atas seluruh perkakas kami: `4-erd-dan-tabel-datar/SAPUAN-DAN-NAMA-TAGNYA.md`. Dua
terdampak — penyisir tetapan (aktivitas saja) dan **penyisir langkah hidup**, yang mengenali
`pyStepsBlockName` tetapi **tidak** `pyDisabled`; keduanya penanda mati, dan **irisannya nol**.

**Tiga jebakan yang lahir bersamanya, dan ketiganya berdiri sendiri:**

1. **Nama sasaran berawalan titik belum tentu menyasar halaman yang sama.** `pyTargetProperty` di
   *Report Definition* berbunyi `.ID`, `.TreatyYear` — relatif terhadap **kelas laporannya**, bukan
   halaman kontrak. **Periksa jenis aturannya sebelum menghitungnya sebagai penulis.**
2. **Nama tag dapat berbohong seperti label langkah.** `pyStepsActivityName` terdengar seperti nama
   aktivitas yang dipanggil; ia berisi `"Call AddCommentList_Act"` — metode **dan** sasarannya
   menyatu. Mencocokkannya secara penuh menghasilkan nol untuk semuanya.
3. **Indeks bawaan yang dihitung mesin punya batasnya sendiri.** `pxRuleReferences` terlihat seperti
   jawaban resmi atas *"siapa memanggil siapa"*; seluruh 63 aktivitas yang disebutnya nol-perujuk
   **ditemukan pencarian teks di berkas lain**. Ia dapat menyatakan sesuatu **dirujuk**, tidak
   pernah **tidak dirujuk**.

**Dan satu hal yang menenangkan, tetapi bukan karena kehati-hatian:** sapuan **pencarian teks** kebal
terhadap kegagalan ini **secara bentuk** — ia tidak mengenal tag sama sekali. Ia tetap tidak dapat
menyatakan **jenis** apa pun, jadi keduanya dipakai bersama, bukan salah satu.

#### 2.0-j.1 Dua bentuk kekeliruan pelaporan yang lahir dari sapuan, dan keduanya akan terulang

> **MELAPORKAN JUMLAH BARIS SEBAGAI JUMLAH TEMPAT.** *"26 `pyDisabled` di data transform"* terbaca
> sebagai 26 cabang; ia **13 baris × dua salinan berkas yang sama**, dan ketiga-belasnya **satu**
> cabang beserta isinya. Setiap angka dari sapuan menyebut **satuannya** — baris, langkah, cabang,
> berkas, atau tempat — dan menyebut bila **dua modul membawa salinan berkas yang sama**, karena
> korpus ini memang memuat keduanya.

> **NOL YANG MUSTAHIL ADALAH ALAT CARI, bukan hasil.** Daftar nama elemen menjawab *"nama mana lagi
> yang dipakai untuk gagasan ini"*. Ia **tidak** menjawab *"gagasan ini dapat berbentuk apa lagi"*.
> Yang menjawab yang kedua adalah hasil yang **tidak mungkin benar**: sebuah properti dengan empat
> layar sendiri dan **nol penulis**. Ketika sebuah sapuan mengembalikan nol untuk sesuatu yang jelas
> dipakai, **yang salah sapuannya**, dan bentuk yang belum terpikir ada di situ. Begitulah bentuk
> penulis **kelima** — salin halaman, `<CopyFrom>`/`<CopyInto>` di parameter langkah — ditemukan
> sesudah keempat bentuk lain didaftar mekanis.

### 2.0a Rujukan bukan panggilan. Panggilan bukan langkah hidup.

Ditetapkan 23 September 2026 setelah saya membuat kesalahan yang sama **tiga kali dalam satu
pemeriksaan**. Ia pasangan alami §2.1, dan kegagalannya berbentuk sama.

> **Nama yang muncul di sebuah berkas, langkah yang memanggilnya, dan langkah yang benar-benar
> berjalan adalah tiga hal berbeda. Hanya yang ketiga perilaku.**

| Tingkat | Yang dibuktikannya | Cara memeriksanya |
|---|---|---|
| **Rujukan** — nama muncul di berkas | tidak membuktikan apa pun | `grep` |
| **Panggilan** — ada langkah yang memanggilnya | benda itu dimaksudkan dipakai | daftar `pySteps` |
| **Langkah hidup** — langkah itu tidak ber-blok `//`, dan tidak ada induknya yang ber-blok `//` | **perilaku** | uraikan pohon `pySteps`, periksa blok induk maupun anak |

**Setiap pernyataan berbentuk "dipanggil dari N tempat" harus menyebut N-nya langkah hidup.**
Kalimat "dirujuk di delapan berkas" bukan pernyataan tentang perilaku dan tidak boleh dipakai
sebagai satu.

Tiga kesalahan yang melahirkan aturan ini, seluruhnya dalam satu pemeriksaan yang sama: menghitung
rujukan `SaveTreatyInOffer_Act` sebagai panggilan hidup; membaca daftar langkah datar tanpa
mengurai sarangnya sehingga langkah anak tertukar dengan langkah induk; dan menyimpulkan "tidak
terlihat hilir" dari satu aktivitas padahal 323 berkas dipakai bersama dua modul.

**Penyisiran mundur sekali** dijalankan atas seluruh klaim berbentuk hitungan pemanggil di dokumen
modul — hasilnya di `PENGETAHUAN.md` dan `INVENTARIS-STRUKTUR-DATA.md`, dan satu klaim **gugur**.


Aturan ini mengikat siapa pun yang membaca ekspor Pega modul ini. Ketiganya lahir dari kesalahan
nyata yang terjadi dan harus dicabut selama sesi penggalian.

### 2.1 Periksa dulu apakah langkahnya hidup, baru baca kondisinya

Langkah dengan `pyStepsBlockName = "//"` **dinonaktifkan** — isinya tidak berjalan. Membaca kondisi
dari blok mati sama menyesatkannya dengan mempercayai label. Aturan ini **mendahului** 2.2.

Dua kali selama pembedahan modul ini, kesimpulan ditarik dari langkah mati dan harus dicabut:
sekali dari langkah penghapus baris pada aktivitas penyimpanan, sekali dari blok penjumlahan limit
pada perhitungan *Rate on Line*. Yang kedua melahirkan laporan penggandaan penyebut yang tidak ada.

Perkakas peringkas apa pun wajib ikut mencetak penanda itu, bukan menyembunyikannya.

#### Perluasan 24 September 2026 — HIDUP belum cukup; kode tindakannya ikut dibaca

Aturan §2.1 di atas benar dan **tidak cukup**. Sebuah kekeliruan lolos justru karena mematuhinya:
langkahnya diperiksa, ternyata **hidup**; kondisinya dibaca, ternyata **benar** — dan kesimpulannya
tetap **terbalik**.

Sebabnya: pada Pega, `pyStepsPreCondParamsWhen` hanya menyatakan **apa yang diuji**. Yang menyatakan
**apa yang terjadi bila uji itu benar** adalah pasangan terpisah:

```
pyStepsPreCondParamsWhenTrue      kode tindakan
pyStepsPreCondParamsWhenTruePrms  parameternya, bila kode itu membutuhkannya
pyStepsPreCondParamsWhenFalse     kode tindakan untuk sisi sebaliknya
```

**Arti kodenya dikalibrasi dari sebarannya, bukan dari kamus mana pun** — disapu atas seluruh
`Treaty In` dan `Treaty In Adjustment`:

| Kode | Jumlah | Membawa nama blok (`…Prms`)? | Artinya |
|---:|---:|---|---|
| **1** | **66** | **66 dari 66 — selalu** | **lompat ke blok bernama itu** |
| **2** | 18.412 | tidak pernah | **lanjut** — ia muncul pada langkah yang tidak berkondisi sama sekali |
| 3 | 4.890 | tidak pernah | belum dikalibrasi |
| 4 | 412 | tidak pernah | belum dikalibrasi |
| 5 | 16 | tidak pernah | belum dikalibrasi |

Kasus yang melahirkannya: sebuah langkah berkondisi `TreatyIn.EDMState=="3"` dengan kode `1` dan
parameter `jmp`, sementara langkah beberapa baris di bawahnya ber-`pyStepsBlockName = jmp`.
Dibaca tanpa kode tindakannya, ia terbaca *"berjalan untuk jenis 3"*. Sebenarnya ia
**melompati jenis 3** — **dan membawa dua langkah berikutnya ikut terlewat**, karena lompatannya
mendarat sesudah keduanya.

> **Maka: kondisi tanpa kode tindakannya hanya separuh kalimat, dan separuh yang tersisa dapat
> membalik artinya.** Setiap klaim tentang "langkah ini berjalan ketika X" menyebut **ketiganya** —
> kondisinya, kode tindakannya, dan ke mana lompatannya mendarat.

Ini **bukan** turunan §2.0: strukturnya berlaku sepenuhnya, tidak mati, tidak tersembunyi, dan
labelnya tidak berbohong. Yang terbalik **pembacaan kita**.

**Sisa yang belum diperiksa:** 20 dari 66 pemakaian kode `1` menunjuk nama blok yang tidak ditemukan
di aktivitas yang sama. Mungkin bloknya terdefinisi di tempat lain, mungkin lompatan ke tempat yang
tidak ada. **Tidak dihitung sebagai cacat** sampai diperiksa.

### 2.2 Label langkah tidak dipercaya; hanya kondisinya yang dibaca

Dalam modul ini saja, deskripsi langkah bertentangan dengan preconditionnya **lima kali**:

| Label | Kenyataan |
|---|---|
| "when edmtype = 2" | berjalan untuk dua keadaan; memo "for edmtype != 3" yang benar |
| "Spreading Lama" / "Spreading Baru" | bukan dua generasi — jalur otomatis dan jalur manual |
| `@SizeOfPropertyList(.MDPList)>2` | kondisinya `>=2`; memo aturannya sendiri: "perbaiki when 3.3" |
| "if empty exit activity" | memasang pesan galat, **tidak ada kode keluar**; aktivitas lanjut |
| "when RNMShareAcrossTheBoard != true" | kondisinya `TreatyIn.OptionLimit == 1` |

### 2.2c LINGKUP DITULIS DARI PENULISNYA, BUKAN DARI BENTUK YANG TERLIHAT

Ditetapkan 24 September 2026, sesudah invarian **ketiga** ditemukan salah lingkup dalam satu hari.

> **LINGKUP DITULIS DARI BENTUK YANG TERLIHAT, BUKAN DARI PENULISNYA.**

Ketiganya keliru dengan bentuk yang sama. Sebuah daftar **tampak** menggantung langsung pada induk
besarnya, lalu keunikannya ditulis *"unik di dalam induk besar itu"* — padahal aktivitas yang
menulisnya menunjukkan ada **satu tingkat pengelompokan di antaranya**.

| Invarian | Yang ditulis | Yang sebenarnya |
|---|---|---|
| INV-47, INV-50 | dijumlahkan "per sumbu **pihak**" | sumbunya **jenis reasuransi**; tidak ada sumbu pihak di keluarga penyebaran |
| INV-12 | nomor termin unik **di dalam versi** | termin disusun **per mata uang**; kontrak bermata uang dua punya dua rangkaian 1..N |

Contoh yang paling bersih adalah yang terakhir: daftar bernama `Installment` ternyata baris **per
mata uang**, dan termin yang sebenarnya hidup satu tingkat di bawahnya di `InstallmentList`.
**Constraint yang ditulis dari bentuk yang terlihat akan menolak data yang sah.**

**Cara menerapkannya:** sebelum menulis lingkup sebuah keunikan, buka aktivitas yang **menambah
baris** ke daftar itu dan lihat apa yang dilewati gelungnya. Labelnya kadang justru jujur — yang
satu ini berbunyi *"Set currency list (also total's currency)"* — tetapi hanya terbaca oleh yang
membuka penulisnya.

### 2.3 Pisahkan "tertulis di aturan" dari "benar-benar berjalan"

Setiap dokumen yang menyebut aturan, pemeriksaan, atau penjaga di sistem lama **wajib** memisahkan
keduanya secara terlihat, dan menyebutkan sebab tiap aturan yang tidak berjalan (langkah mati,
pemanggil mati, hanya kontrol layar).

Ini jenis kesalahan yang paling mahal: ia membuat orang percaya ada penjaga yang tidak ada, lalu
perkiraan pekerjaan menganggap pemulihannya gratis. Sesuatu yang tidak berjalan hari ini adalah
**kemampuan baru**, bukan pelestarian.

### 2.4 Model datanya ada di JSONDATA, bukan di tabel relasional

`M_TREATY_IN` menyimpan seluruh halaman clipboard sebagai satu kolom `JSONDATA`. `TREATY_IN` hanya
punya 20 kolom dan `TREATYINDETAIL` adalah proyeksi datar. Struktur data disusun dari `JSONDATA`.

Tabel relasional dipakai untuk hal lain: ia memberi tahu field mana yang dulu dianggap layak
dilaporkan. Itu masukan tentang **prioritas**, bukan tentang **kelengkapan**.

### 2.4b Kata yang berarti lain di konteks lain tidak boleh muncul telanjang di keduanya

Aturan glosarium, ditetapkan 23 September 2026 setelah **tabrakan nama ketiga** dalam dua sesi:

| Kata | Arti di satu tempat | Arti di tempat lain |
|---|---|---|
| **penyebaran** | pembagian bagian NuRe ke susunan retro internal | penempatan risiko ke pihak luar |
| **keadaan** (`ViewState`) | bendera tampilan layar | kedudukan kontrak dalam siklus hidupnya |
| **fakultatif** | risiko yang **diterima** NuRe per risiko — sumber klaim, milik modul Klaim | bagian NuRe sendiri yang **ditempatkan keluar** per risiko — milik Treaty In |

Setiap kata semacam itu **wajib membawa pembedanya di dalam namanya** — arah, pemilik, atau
lingkupnya — di kedua konteks, bukan hanya di salah satunya. Ini penerapan ADR-0021: pengenal
dinamai menurut isinya dan pemiliknya.

Untuk fakultatif, pembedanya **arah**, dan arah itu mutlak: yang satu uang masuk, yang lain uang
keluar. Karena itu ADR-0020 — *"sistem ini hanya memiliki Treaty; fakultatif dibaca, tidak
dimodelkan"* — dinyatakan **tidak berlaku** di Treaty In: ia bicara tentang **fakultatif masuk**,
sementara yang dimiliki Treaty In adalah **fakultatif keluar**. Alasannya bukan "subjeknya
berbeda", melainkan **arahnya berlawanan**.

### 2.5 Cacat bersembunyi di balik nilai bawaan

Pada setiap rumus, tanyakan: *nilai bawaan apa yang membuat rumus ini tampak benar, dan apa yang
terjadi kalau nilai itu bukan bawaannya?* Pola ini muncul empat kali di modul ini — jepitan nol
yang identik dengan lantai premi selama deposit sama dengan minimum; jalur otomatis yang tidak
butuh validasi selama persentase master berjumlah 100; tarif premi reinstatement yang tertimpa
tanpa akibat selama keduanya 100; dan Rate on Line yang benar selama layer bermata uang tunggal.

---

### 2.6 Keluaran ditulis saat LANGKAHNYA selesai, bukan saat SESINYA selesai

Ditetapkan pemilik proses pada 23 September 2026, berlaku untuk **sisa proyek** dan untuk **setiap
modul**, bukan hanya Treaty In.

> **Sebuah sesi dapat berakhir kapan saja. Apa yang hanya ada di konteks kerja tidak ada.**

Setiap berkas keluaran ditulis ke disk begitu langkah yang menghasilkannya selesai — bukan
dikumpulkan di akhir sesi. Berlaku untuk spesifikasi, peta telusur, ADR, catatan, maupun skrip DDL.

**Sebabnya tercatat supaya tidak diulang.** Sesi to-spec Treaty In menjalankan tujuh dari sepuluh
langkah dan menghasilkan hampir nol artefak yang bertahan: keluarannya tersusun di prompt sebagai
daftar di bagian terpisah, sehingga terbaca sebagai pekerjaan akhir. Ketika konteks kerjanya
hilang, tujuh langkah kerja hilang bersamanya dan sesi berikutnya menemukan empat dari sembilan
masukannya tidak ada.

Pemilik proses mencatat bahwa aturan yang benar sebenarnya sudah pernah ditulis — untuk peta
telusur JSON secara khusus: *"jangan dikerjakan di akhir sebagai formalitas, isi sambil berjalan"* —
tetapi tidak digeneralisasikan ke seluruh keluaran. Aturan ini adalah generalisasi itu.

**Penerapannya ke belakang:** prompt sesi mana pun yang mendaftar keluarannya di satu bagian dan
urutan kerjanya di bagian lain dibaca dengan aturan ini, bukan menurut susunannya. Setiap berkas
ditulis saat langkah yang menghasilkannya selesai.

### 2.7 Perkakas analisis disimpan sebagai BERKAS, bukan perintah sekali pakai

Ditetapkan 23 September 2026, dan ia pasangan §2.6.

Penyisir dan pengklasifikasi yang dipakai membedah ekspor ditulis sebagai berkas skrip di
scratchpad, bukan dijalankan sebagai perintah sekali tempel. Ketika konteks kerja sesi to-spec
hilang, **skripnya bertahan** — sehingga adjudikasi 667 jalur JSON dapat **dijalankan ulang atas
sumbernya** dan hasilnya terverifikasi, bukan diingat.

Perbedaannya menentukan: angka yang diingat adalah klaim; angka yang dapat dijalankan ulang adalah
bukti. Perkakas yang bertahan mengubah pemulihan dari rekonstruksi menjadi reproduksi.

Berlaku juga saat tidak ada yang hilang: skrip yang tersimpan dapat dijalankan ulang ketika ekspor
diperbarui, sehingga temuan dapat disegarkan tanpa menyusun ulang metodenya.

### 2.8 Cabang sebuah properti ditentukan `pzIndexOwnerKey`, bukan nama berkas

Berkas Section di ekspor **membundel aturan lain di dalamnya**.
`Section/TreatyInNONProportional.xml` memuat entri indeks milik empat aturan berbeda, termasuk 193
entri milik `TREATYINTABSPROPORTIONAL`.

Akibatnya, menyimpulkan cabang sebuah properti dari nama berkas tempat ia ditemukan dapat salah
total — sebuah properti proporsional akan tampak dipakai di cabang non-proporsional semata karena
nama berkas pembungkusnya. Bacalah `pzIndexOwnerKey` pada entri rujukannya.

Ini bentuk yang sama dengan §2.2: nama berkas adalah label, dan label tidak dipercaya.

### 2.9 Nama orang: pemegang sekarang adalah turunan, pelaku dahulu adalah fakta

Ditetapkan 23 September 2026 setelah keputusan atas satu atribut ternyata berlaku umum. Ia akan
dipakai lagi di modul lain, karena setiap modul Pega NuRe menyimpan nama orang di beberapa tempat.

> **Nama yang menyatakan SIAPA MEMEGANG SESUATU SEKARANG adalah turunan yang akan usang.
> Nama yang menyatakan SIAPA MELAKUKAN SESUATU DAHULU adalah fakta beku.**

| Bentuk | Contoh di Treaty In | Perlakuan |
|---|---|---|
| pemegang sekarang | `PositionUsername` — nama orang yang memegang kontrak | **dibuang**; penugasan yang berlaku dibaca dari peran bertanggal (ADR-0044) |
| pelaku dahulu | `CommentList[].OperatorName` — siapa menyetujui, kapan | **disimpan**; membekukan namanya pada saat keputusan justru yang benar (ADR-0045) |

Pembedanya bukan apakah isinya nama orang, melainkan **apakah nilainya bisa berubah tanpa ada
peristiwa baru**. Pemegang berubah ketika orangnya pindah jabatan, tanpa apa pun terjadi pada
kontraknya — itu tanda ia turunan. Pelaku tidak pernah berubah, karena peristiwanya sudah lewat.

Aturan yang sama berlaku pada **nama perusahaan**: nama cedant, reasuradur pemimpin, dan broker
adalah milik data master dan **dirujuk**, tidak disalin (ADR-0023, ADR-0041). Yang disalin hanya
bila ia bagian dari peristiwa yang dibukukan.

### 2.9b Lubang yang dicatat membawa PEMILIK dan SAAT PENAGIHAN

Ditetapkan 24 September 2026, dari kalimat yang lahir saat sebuah properti ditemukan hilang untuk
**kedua** kalinya:

> **Lubang yang dicatat tetapi tidak ditindaklanjuti berperilaku persis seperti lubang yang tidak
> diketahui.**

`SpreadingTypeID` disebut **dua kali** di `PETA-TELUSUR-JSON.md` §5 sebagai properti yang diketahui
hilang dari penyisiran. Ia tidak pernah berjalan sampai menjadi atribut, dan baru masuk model
sembilan hari kemudian lewat jalan yang sama sekali lain. **Mencatatnya tidak menolong sama sekali.**

Maka bentuk yang mengikat:

> **SETIAP LUBANG YANG DICATAT MEMBAWA DUA HAL: SIAPA YANG MENUTUPNYA, DAN APA YANG MEMBUATNYA
> DITAGIH.**
>
> Lubang tanpa pemilik dan tanpa saat penagihan bukan lubang yang tercatat — ia **lubang yang
> dilupakan dengan cara yang lebih rapi**.

"Siapa menutupnya" saja **tidak cukup**, dan itu justru yang membuat aturan ini perlu: daftar lubang
yang sudah punya kolom pemilik tetap gagal, karena tidak ada satu pun langkah yang **wajib
membacanya**.

**Saat penagihan** adalah peristiwa yang **pasti terjadi** dan yang pembacanya **tidak dapat
melewatinya** — *"saat §10 entitas itu ditulis"*, *"saat gerbang berikutnya dijawab"*, *"saat DDL
menulis constraint-nya"*. Ia **bukan tanggal** dan **bukan niat**.

### 2.10 Kepanjangan sebuah singkatan ada di pola penamaan SAUDARANYA, bukan di kamus industri

Ditetapkan 24 September 2026 setelah `PLA` terpecahkan — singkatan yang sudah berstatus *"tidak ada
kepanjangannya di mana pun"* selama tujuh sesi.

Ia tidak dipecahkan dengan menebak dari pengetahuan pasar, dan tidak dengan mencari kamus. Ia
dipecahkan dengan **membaca tetangganya di kelas yang sama**. Aktivitas konversi
`TreatyInMappingDataconvertProp.xml` memasangkan setiap daftar `Detail[]` dengan ruas asalnya pada
kelas `Data-TreatyInBusinessLimit`, dan kelas itu menamai ruasnya **panjang-panjang**:

```
PreliminaryLossAdvice  ->  PLAList
ClaimCooperation       ->  ClaimCoopList
CashLossLimit          ->  CashLossList
AmountEPI              ->  EPIList
```

**Caranya, untuk dipakai lagi:**

1. Cari **kelas mana** yang memuat singkatan itu.
2. Cari aktivitas **konversi atau pemetaan** yang menulis ke sana — bentuk sumbernya sering memakai
   nama panjang, karena ia ditulis untuk dibaca manusia atau sistem lain.
3. Baca **pola penamaan saudara-saudaranya** di kelas itu. Bila saudaranya dieja penuh, singkatan
   itu punya pasangan penuh di suatu tempat.
4. Sapu **seluruh korpus**, bukan satu modul.

**Dan hasil negatif ikut dicatat, beserta jebakannya.** Teknik ini dijalankan atas `RSMD` pada hari
yang sama dan **tidak menemukan kepanjangannya**: 44 kemunculan di Treaty In dan Treaty In
Adjustment, nol di modul lain, dan saudaranya di layar memang dieja penuh — *"Earthquake Limit"*,
*"Flood Limit (Jabodetabek)"*, *"Flood Limit (Nationwide)"* — sementara ia sendiri tetap
*"RSMD Limit"*. Yang diperoleh hanya **golongannya**: ia **nama bahaya**, sekeluarga dengan gempa
dan banjir, yang mengukuhkan `BAHAYA` sebagai tabel acuan.

> **Jebakan yang sudah memakan satu pencarian:** menyapu `rsmd` tanpa memperhatikan huruf besar
> memungut **`KursMDP`**. Satu-satunya kemunculan "RSMD" di modul Claim Non Prop adalah itu, dan ia
> bukan RSMD sama sekali. Dicatat supaya tidak ada yang mengulanginya dan mengira ia menemukan
> sesuatu.

`RSMD` **tetap butir wawancara.** Yang berubah: sekarang tercatat bahwa korpusnya sudah disapu
habis, sehingga pencarian berikutnya tidak diulang melainkan ditanyakan ke orang.

---

## 3. Glosarium

Istilah reasuransi dibiarkan dalam bahasa aslinya bila itu yang dipakai di pasar.

### 3.1 Benda pokok

**Treaty In** — kontrak reasuransi masuk: NuRe menerima sesi risiko dari *cedant*. Satu kontrak
berlaku untuk satu periode dan satu *source of business*.

**Cedant** (di sistem lama: `Ceding`) — perusahaan asuransi yang menyerahkan risiko kepada NuRe.

**Source of Business (SoB)** (`LeadingReinsSource`) — asal bisnis; bersama cedant, periode, dan
`ProportionType` membentuk **kunci alami** sebuah kontrak.

**Leading Reinsurer** — reasuradur pemimpin yang menetapkan syarat; NuRe mengikutinya.

**Addendum** (di sistem lama: EDM) — perubahan atas kontrak yang sudah disetujui, dicatat sebagai
versi tersendiri, bukan sebagai penimpaan. Lihat ADR-0040.

**Rangkaian treaty** — dugaan bahwa beberapa kontrak antar-tahun membentuk satu kesinambungan.
**Belum dipastikan ada**; bergantung pada apakah riwayat lintas tahun pernah diminta. Sampai
dipastikan, ia bukan bagian dari model.

### 3.2 Dua cabang

**Proporsional** — NuRe menerima persentase dari setiap risiko.

- **Quota Share (QS)** — satu persentase untuk seluruh portofolio. Yang disepakati **satu
  bilangan**: persentasenya. Retensi adalah definisinya, bukan kesepakatan terpisah.
- **Surplus** — cedant menahan **retensi** sebagai *jumlah uang*, dan menyerahkan kelebihannya
  sampai sekian **lines**. Yang disepakati **dua hal**: retensi dan jumlah lines. Kapasitas adalah
  hasil kali keduanya. Persentase sesi **tidak pernah menjadi syarat kontrak** di Surplus, karena
  ia berbeda untuk setiap risiko.

**Non-proporsional (XOL)** — NuRe menanggung kerugian di atas sebuah *deductible* sampai sebuah
*limit*, disusun dalam **layer**.

### 3.3 Besaran

**Limit** — batas tanggungan sebuah layer. Besaran **tingkat 100%**.

**Deductible** — batas bawah tanggungan layer.

**EPI** (*Estimated Premium Income*) — perkiraan premi. **Belum dipastikan** apakah dicatat pada
tingkat 100% atau bagian NuRe; rumus pencapaian mengandaikan yang pertama. Lihat daftar wawancara.

**EGNPI** (*Estimated Gross Net Premium Income*) — dasar perhitungan premi XOL.

**MDP** (*Minimum and Deposit Premium*) — dua besaran, bukan satu: **deposit** yang dibayar di muka
dan **minimum** yang menjadi lantai. Keduanya **tidak wajib sama**.

**Rate on Line (ROL)** — premi layer dibagi limitnya. Kebalikannya adalah *payback period*.

**Reinstatement** — pemulihan limit setelah tergerus klaim. Dikutip **berpasangan**: porsi limit
yang dipulihkan, dan tarif premi yang ditagih untuk pemulihan itu. Keduanya dinegosiasikan
**terpisah** — "pulihkan penuh, gratis" adalah ketentuan yang sah. Hubungannya **perkalian**,
bukan persamaan.

**Deduction / potongan** — dibedakan menurut **dasar perhitungannya**, bukan menurut tetap atau
proporsional:

- terhadap **premi bruto**: brokerage, ceding commission, overriding, pajak, premium reserve;
- terhadap **laba**: *profit commission* — tidak bisa dihitung sebelum periode selesai;
- **jumlah yang disepakati**: tidak bergerak mengikuti premi.

**Spreading / penyebaran** — pembagian bagian NuRe ke susunan retro internal NuRe.

**Tingkat pencatatan** — apakah sebuah nilai uang dinyatakan pada **100% treaty** atau pada
**bagian NuRe**. Tingkat adalah **sifat besarannya**, bukan konvensi yang dipilih. Lihat ADR-0039.

### 3.3b Aturan penamaan istilah

**Istilah pasar dipertahankan dalam bentuk aslinya, dan setiap istilah yang dipertahankan WAJIB
punya entri di glosarium ini beserta kepanjangan dan artinya. Selebihnya Bahasa Indonesia.**

Yang berbahaya bukan istilah asing; yang berbahaya adalah **singkatan yang kepanjangannya tidak ada
di mana pun**. Dengan aturan itu `treaty`, `EGNPI`, `MDP`, `PLA`, `cash loss`, dan
`claim cooperation` semuanya boleh — tidak satu pun boleh muncul tanpa entri.

**Dua singkatan yang statusnya berubah 24 September 2026:**

| Singkatan | Keadaan | Dari mana |
|---|---|---|
| **`PLA`** | **TERPECAHKAN — *Preliminary Loss Advice*.** Ambang yang mewajibkan cedant memberi pemberitahuan awal kerugian | ruas `PreliminaryLossAdvice` pada `Data-TreatyInBusinessLimit`, yang mengisi `PLAList`. Nama sementara `TBD_PLA_ARTI_BELUM_DIKETAHUI` **dicabut** |
| **`RSMD`** | **MASIH TEBAKAN**, dan kini diketahui **tidak ada di korpus** — 44 kemunculan, nol kepanjangan. Yang pasti hanya golongannya: **nama bahaya**, sekeluarga dengan gempa dan banjir | disapu §2.10; tetap **butir wawancara** |

**Tingkat pencatatan `PLA` tetap terbuka.** Kepanjangan menjawab *apa*, bukan *pada tingkat mana
angkanya dicatat* — ia tetap satu dari lima di `SPEC-MODEL-DATA.md` §7.

Dua penerapan yang sudah diputuskan:

| Nama lama | Nama baru | Alasan |
|---|---|---|
| `CessionList`, `CessionPct` | `NILAI_PENYERAHAN`, `PERSEN_PENYERAHAN` | **"sesi" dilarang** — aturan 2.4b. Dalam proyek ini "sesi" berarti sesi kerja dan dipakai puluhan kali sehari; sebagai istilah reasuransi ia berarti penyerahan risiko. Yang lebih mudah diubah adalah yang di model |
| `IOOLimitList` | `KAPASITAS_SURPLUS` | kepanjangan "IOO" **tidak ada di mana pun** di ekspor. Dinamai menurut apa yang dihitungnya: retensi × jumlah *lines*. Nama lamanya dicatat di peta telusur |

### 3.4 Istilah sistem lama yang TIDAK dibawa

| Istilah lama | Nasib |
|---|---|
| `ViewState`, `IsEditData`, `RevisionState`, `Position`, `StatusAkseptasi` | dilebur jadi **satu keadaan siklus hidup**; lihat ADR-0046 |
| `CessionPct` | pecah: di QS ia persentase sesi, di Surplus ia jumlah lines x 100 — dua benda berbeda |
| `SpreadingTypeID` sebagai pemilih mode | mode tidak ditentukan oleh keterisian sebuah field |
| `autocalculate` | tidak ada sakelar yang mematikan perhitungan |
| `9989998` / `998999800` | kegagalan bukan nilai; lihat ADR-0035 |

---

## 3.5 Nama skema

| | |
|---|---|
| Skema | **`TREATY_MASUK`** |
| Akun aplikasi | **`TREATY_MASUK_APP`** |

**Bukan `TREATY_IN`.** `POOLDATA.TREATY_IN` dan `POOLDATA.TREATY_IN_EDM` **sudah ada** sebagai tabel
sistem lama, dan ADR-0028 menempatkan skema baru bersebelahan dengan `POOLDATA` di instance yang
sama. Selama masa berdampingan, "cek TREATY_IN" akan berarti dua hal berbeda dalam satu percakapan.

"treaty" dipertahankan sebagai istilah pasar; "masuk" adalah kata Indonesia utuh, bukan singkatan
buatan — dan justru bagian **arah** inilah yang berbahasa Indonesia, karena ia keterangan, bukan
istilah pasar. Ia juga menyiapkan `TREATY_KELUAR` bila konteks Treaty Out kelak dibangun.

Preseden `KLAIMNP` **tidak dijadikan patokan**: "NP" adalah singkatan buatan yang aturan penamaan
justru larang, jadi mengikuti bentuknya berarti mewarisi penyimpangannya.

## 4. Batas konteks

### 4.1 Di dalam Treaty In

Kontrak, addendum, layer, detail proporsional, penyebaran, reinstatement, potongan, installment,
jadwal pelaporan, dan persetujuan atas semuanya.

### 4.2 Di luar, dan pemiliknya harus ditetapkan

| Konteks | Isinya | Status |
|---|---|---|
| **Program retro NuRe sendiri** | `PROPORTIONALARRG` — susunan baku NuRe. Tidak ada satu pun aturan Pega yang menulisnya. Diduga milik konteks *Treaty Out* yang belum bernama | pemilik belum ada; Treaty In menjadi **pembaca** lewat adapter yang gagal keras |
| **Penulis produksi dan pencapaian** | `TREATYINPRODUCTION`, `ACHIEVEMENT` | penulisnya di luar ekspor ini; harus ditetapkan |
| **Layanan dokumen bersama** | token, `T_STORAGE_IMAGE`, `M_LINK_SERVICE` | Treaty In **dan** Klaim sama-sama pemakai; pemiliknya belum ada |
| **Klaim** | gerusan limit yang mungkin menentukan porsi reinstatement | bersinggungan; belum ditelusuri |
| **Buku besar / keuangan** | kurs pembukuan dan penilaian ulang posisi valas | harus dipastikan; bila berbeda dari kurs modul ini, itu temuan tersendiri |
| **Treaty In Adjustment** | selisih nilai kontrak | **bukan konteks tetangga** — bagian modul ini sejak GRL-01 |

### 4.3 Batas yang ditulis sebelum ada yang membangun

Bila proses latar pengawas jatuh tempo dibangun kelak, ia **mencatat peristiwa tentang kontrak**
— misalnya "laporan periode Q2 terlambat sekian hari", sebagai catatan tersendiri — dan
**tidak pernah mengubah kontraknya**. Proses latar yang mengubah keadaan kontrak yang sudah
disetujui termasuk kelas yang sama dengan suntingan yang mengubah penyebaran secara surut.

Batas ini ditulis sekarang, sebelum ada yang membangunnya, karena batas yang ditulis sesudahnya
jauh lebih mahal.
