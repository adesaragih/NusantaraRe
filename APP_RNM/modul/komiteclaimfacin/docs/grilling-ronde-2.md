# Grilling Ronde 2 — Komite Claim Fac In · JANTUNG MODUL

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim FacIn\` — **114 berkas**, READ-ONLY
**Pembanding:** `Claim Fac In` *(482)* · `Komite Claim Prop` *(80)*

> ⛔ **Ronde ini MEMBACA.** Ia tidak memutuskan, kecuali **mencatat dua keputusan work owner** yang
> sudah diberikan *(§8)*. ⛔ Nol spec, nol tiket, nol kesimpulan untuk Go.
> ⛔ Nol berkas korpus dibuka untuk ditulis · nol nilai kredensial dibaca, dicetak, atau disalin.

**Sensus ronde ini:** ⭐ **RALAT terhadap ronde 1: dua** · butir `[terbuka]` ditutup **6** ·
`[terbuka]` baru **4** · keputusan work owner dicatat **2** · ⭐⭐ **temuan yang membalik vonis: satu**

---

## §1 — Empat rule `KomitePost_*`

### 1a — ⭐ `KomitePostAct` ternyata **PENYALUR**, bukan pelaksana

`[terverifikasi]` **6 langkah**, dan isinya **percabangan menurut jenis penyerahan**:

| Langkah | Gerbang | Memanggil | Keadaan |
| --- | --- | --- | --- |
| 1 | — | `Obj-Refresh-And-Lock` | hidup |
| 2 | `TransferType=="2"` | **`KomitePost_Adjustment`** | hidup |
| ⛔ 3 | `TransferType=="1"` | ~~`KomitePost_Survey`~~ | ⛔ **BER-REMARK — MATI** |
| 4 | `TransferType=="3"` | **`KomitePost_Reject`** | hidup |
| 5 | `TransferType=="4"` | **`KomitePost_CloseClaim`** | hidup |
| ⛔ 6 | — | `Property-Set` | ⛔ ber-remark |

⭐⭐ **Empat jenis penyerahan, dan yang jenis `1` — SURVEY — jalurnya MATI.**
⚠️ Rule `KomitePost_Survey` **tidak ada di daftar 23 rule khas**, dan ⛔ **belum saya cari apakah
berkasnya ada sama sekali**. `[terbuka]` **butir 14**

⚠️ **Dan ini menyambung dengan temuan ronde 1:** `KomitePost_CloseClaim` dan `KomitePost_Reject`
masing-masing punya gerbang **berflag `false`** pada langkah 1 yang berbunyi
`pyWorkPage.TransferType=="1" && …` — ⭐ **gerbang mati yang menguji jenis penyerahan yang jalurnya
juga mati**.

### 1b — ⭐⭐ Urutan: **SIMPAN DULU, efek keluar menyusul**

`[terverifikasi]` Pada keempat rule, **tanpa kecuali**:

| Rule | Langkah | `Obj-Save` | Efek keluar sesudahnya |
| --- | ---: | --- | --- |
| `KomitePost_Adjustment` | **72** | **7.2.1.10 · 7.2.1.11** | surel *(7.2.1.15)* · konversi *(17)* · Kasir *(25.2.1.1)* |
| `KomitePost_CloseClaim` | **38** | **7 · 8** | surel berlampiran *(11.10.8)* · konversi *(12.4)* |
| `KomitePost_Reject` | **47** | **8 · 9** | surel berlampiran *(13.10.8)* · konversi *(15)* |
| `KomitePostAct` | 6 | — *(hanya kunci)* | — |

> ### ⭐⭐ INI KEBALIKAN KOMITE CLAIM LIFE
>
> Di Komite Claim Life, **kedelapan efek keluar berjalan SEBELUM penyimpanan**, dan itu dicatat
> sebagai **penyimpangan sadar pertama** modul tersebut.
>
> ⭐ **Di sini penyimpanan lebih dulu.** ⛔ **Saya tidak menyimpulkan mana yang lebih benar** —
> yang dicatat: **polanya berbeda**, dan keputusan yang diambil untuk Life **tidak dapat dipinjam
> begitu saja** ke sini.

### 1c — ⭐ Sembilan pembangkit nomor akseptasi **per lini produk: SEMUANYA BER-REMARK**

`[terverifikasi]` `KomitePost_Adjustment` langkah **7.2.1.2.8** sampai **7.2.1.2.14** memanggil
`GenerateNoAccept` dan `GenerateNoAcceptNonFire`, **masing-masing bergerbang satu lini produk**:

```
7.2.1.2.8   GenerateNoAccept          gerbang OutputData.START_DATE==""   [REMARK]
7.2.1.2.9   GenerateNoAcceptNonFire   gerbang IsAneka                     [REMARK]
7.2.1.2.10  GenerateNoAcceptNonFire   gerbang isGolfInsurance             [REMARK]
7.2.1.2.11  GenerateNoAcceptNonFire   gerbang IsMarineCargo               [REMARK]
7.2.1.2.12  GenerateNoAcceptNonFire   gerbang IsMBU                       [REMARK]
7.2.1.2.13  GenerateNoAcceptNonFire   gerbang IsTravel                    [REMARK]
7.2.1.2.14  GenerateNoAcceptNonFire   gerbang IsPA                        [REMARK]
```

⭐⭐ **Ketujuhnya MATI.** ⚠️ Yang **hidup** adalah langkah **7.2.1.2.6** — `GetSequenceNumber_SQL`,
bergerbang `OutputData.START_DATE==""`.

⭐ **Artinya penomoran akseptasi per lini produk sudah ditinggalkan, diganti satu pembangkit
tunggal.** ✅ Sejalan **ADR-0006**. ⛔ **Kapan dan kenapa tidak tertulis di korpus.**

### 1d — ⭐⭐ `COMMIT;` jatuh di mana — dan **nol jalur kegagalan**

`[terverifikasi]` Letak rule ber-`COMMIT;` relatif terhadap `Obj-Save`:

| Rule | `COMMIT;` | Relatif `Obj-Save` |
| --- | --- | --- |
| `KomitePost_Adjustment` | `GetSequenceNumber_SQL` **7.2.1.2.6** | ⚠️ **SEBELUM** simpan |
| | `SaveOSClaim_SQL` **8** · `InsertHistoryAkseptasiPega_Sql` **21** | sesudah |
| `KomitePost_CloseClaim` | `SaveOSClaim_SQL` **12 · 12.3** | sesudah |
| `KomitePost_Reject` | `InsertClaimRejected_Sql` **14 · 14.2** | sesudah |

⚠️ **Satu `COMMIT;` jatuh SEBELUM `Obj-Save`** — pembangkit nomor urut. ⭐ **Itu justru yang
dikehendaki:** nomor yang sudah terbit tidak boleh ikut batal. ✅ Sejalan keputusan §8-Q2.

> ### ⭐⭐ ANGKA YANG PALING KERAS DI RONDE INI
>
> `[terverifikasi]` **Dari 25 langkah `Obj-*` di seluruh modul, dan dari 14 langkah `Obj-*` di
> keempat rule ini — NOL punya jalur kegagalan.**
>
> ⛔ Tidak satu pun `Obj-Save`, `Obj-Refresh-And-Lock`, atau `Obj-Open-By-Handle` memiliki baris
> transisi bersyarat. ⚠️ **Bila penyimpanan gagal, tidak ada yang menangkapnya** — dan efek keluar
> tetap berjalan sesudahnya.
>
> ✅ **Butir terbuka 5 ronde 1 DITUTUP** dengan angka ini.

### 1e — Tulis ke kasus induk dari keempat rule: **NOL**

`[terverifikasi]` ⛔ **Nol penugasan bersasaran `Primary.*` atau `pyWorkCover.*`** di keempat rule
`KomitePost_*`.

⭐ **Jadi ke-18 penulisan ke induk yang ronde 1 catat ada di rule LAIN** — terbaca di §6:
`HitServiceToKasirKMT_Act` *(`Primary.StatusKasir`)*, `PreSecurityReas_Act`, `PreShowRetro_Act`,
dan rule lampiran.

⚠️ `KomitePost_Adjustment` langkah **26** ber-`Obj-Save` dengan **flag `false`** — ⭐ **gerbangnya
mati, jadi langkahnya berjalan TANPA saringan**, padahal syarat tertulisnya
`pyWorkPage.KomiteCount==pyWorkPage.KomiteLoop` *(hanya di tingkat akhir)*.

---

## §2 — ⭐⭐⭐ `SetProteksiSubmiteKomite` — dan ia **BUKAN** yang saya duga

### 2a — Bentuknya: **tiga langkah**

`[terverifikasi]` Kelas `ASM-FW-GCNMFW-Work-Komite`.

| Langkah | Isi |
| --- | --- |
| 1 | `ProteksiKomite.CARI1 = 0` ⭐ **baku: TIDAK boleh** · `Local.USERNAME` diisi dari akun operator |
| 2 | *(pembungkus)* |
| ⭐ **2.1** | `ProteksiKomite.CARI1 = 1` ⭐ **boleh** — bergerbang **tiga syarat** |

### 2b — ⭐⭐ Rantai gerbang langkah 2.1, dibaca utuh

```
[5/2]  OperatorID.pyPosition == "IT Developer"      <== kode 5 pada sisi BENAR
[2/3]  .KomiteAproval == 0
[2/3]  Local.USERNAME == .KomiteID
```

> ### ⭐⭐⭐ TEMUAN YANG MEMBALIK VONIS ADR-0014
>
> ⭐ **Syarat ketiga — `Local.USERNAME == .KomiteID` — ADALAH penegakan pemilik giliran.**
> Ia memeriksa bahwa **pengguna yang sedang bekerja adalah pemilik `KomiteID`**, persis yang
> ADR-0014 tuntut.
>
> ⚠️⚠️ **Dan syarat PERTAMA melewatinya.** Kode **5** *(Skip Whens)* pada sisi benar berarti:
> begitu `pyPosition == "IT Developer"` terpenuhi, **kedua syarat sisanya TIDAK DIUJI LAGI**, dan
> `CARI1` tetap disetel **1**.
>
> ⭐ **Jadi modul ini PUNYA penegakan — dengan satu pintu belakang yang melewatinya.**

### 2c — ⭐ Ke mana `ProteksiKomite.CARI1` pergi

`[terverifikasi]` **Jendela: seluruh 114 berkas.** Ia dibaca **di satu tempat**:

**`Section\ShowTransfer.xml`**, medan `pyDisabledWhen`:

```
ProteksiKomite.CARI1 == 0  ||  pyWorkPage.Adjustment.AcceptedNo != ''
```

⭐⭐ **Kendali keputusan pada layar komite DINONAKTIFKAN bila `CARI1` bernilai 0** — yaitu bila
pengguna **bukan pemilik giliran** dan **bukan "IT Developer"**.

⚠️⚠️ **Tetapi penegakannya ada di LAYAR, bukan di lapisan layanan.** ⛔ ADR-0014 menyatakan tegas
bahwa mengandalkan penempatan atau tampilan **tidak cukup**: *"Pengguna lain ditolak di lapisan
layanan"*.

⭐ **Pemanggilnya satu:** `SetValueKomite` langkah **15** — rule yang juga menyalin nilai
klasifikasi dari kasus induk.

### 2d — ⭐ Apakah ini pintu belakang yang SAMA dengan dua modul lain? **BUKAN**

`[terverifikasi]` **Jendela: seluruh 114 berkas, SETIAP tag.** `OperatorID.pyPosition` muncul di
**tiga berkas**, dan ⭐ **keduanya berbeda urusan**:

| Berkas | Syarat | Yang diaturnya |
| --- | --- | --- |
| `DataTransform\ChronologyInsertion_DT` | `!= "IT Developer"` | ⭐ **penghapusan catatan kronologi** — **SAMA** dengan Claim Prop dan Komite Claim Prop |
| ⭐⭐ `Activity\SetProteksiSubmiteKomite` | `== "IT Developer"` | ⭐⭐ **wewenang MENYIMPAN KEPUTUSAN KOMITE** — **BARU, belum pernah ada** |
| ⭐ `Activity\ApprovalKomite_Act` | `== "SPV B"` | ⭐ lihat §3 — **wewenang berdasarkan pita nilai uang** |

⭐ **Jawaban tegas: yang di `SetProteksiSubmiteKomite` BUKAN pintu belakang yang sama.** Dua
kemunculan sebelumnya mengatur **catatan kronologi**; yang ini mengatur **siapa boleh memutuskan
uang**. ⚠️ **Wilayahnya jauh lebih berat.**

⭐ Catatan pengembang pada rule itu berbunyi *"Add kondisi when OperatorID.pyPosition=='IT
Developer'"* — ⭐ **ia ditambahkan belakangan, bukan bagian rancangan asli.**

### 2e — Pertanyaan, bukan keputusan

⛔ **Tidak diputuskan di sini** — lihat §7.2 **Q1**.

---

## §3 — `KomiteRouter` dan `ApprovalKomite_Act`

### 3a — ⭐ `KomiteRouter`: **tangga empat tingkat lama sudah MATI**

`[terverifikasi]` **8 langkah**, dan **lima di antaranya ber-remark**:

| Langkah | Isi | Keadaan |
| --- | --- | --- |
| ⛔ 1–4 | `param.AssignTo = "komitepnc"` / `"komitepnc2"` / `"komitepnc3"` / `"komitepnc4"`, bergerbang `.KomiteCount==1..4` | ⛔ **BER-REMARK — MATI** |
| ⛔ 5 | `.KomiteCount = .KomiteCount+1` · `.AcceptStatus = ""` | ⛔ **BER-REMARK** |
| 6 | pembungkus, gerbang `Primary.TransferType=='2'` | ⚠️ **flag `false` — gerbang mati** |
| ⭐ **6.1** | ⭐ **`param.AssignTo = .KomiteID`**, bergerbang `.KomiteAproval==0` | ✅ **HIDUP** |
| ⛔ 7 | `param.AssignTo = .Komite.KomiteID`, gerbang `TransferType!='2'` | ⛔ **BER-REMARK** |

⭐⭐ **Antrean berjenjang empat tingkat bernama `komitepnc1..4` sudah ditinggalkan**, diganti
**penempatan langsung ke akun operator penyetuju giliran**.

✅ **Pola yang sama persis dengan Komite Claim Prop**, yang mencatat *"antrean berjenjang lama —
empat antrean bertingkat yang keenam langkahnya di-remark — tidak dibuat"*.

⚠️ **Langkah 6 bergerbang mati** berarti langkah 6.1 **berjalan tanpa saringan `TransferType`** —
⭐ penempatan terjadi untuk **semua jenis penyerahan**, bukan hanya jenis `2`.

### 3b — ⭐⭐ `ApprovalKomite_Act`: **wewenang berdasarkan JABATAN dan PITA NILAI UANG**

`[terverifikasi]` **15 langkah.** Yang menonjol:

| Langkah | Isi |
| --- | --- |
| 1–2.1 | `Local.TotalAdj` **diakumulasi** dari `.ValueAdjustment` |
| 2.2 | `Local.Retro = 1` bila `.IsFacRetro==1` |
| ⚠️ 3 | **`Exit-Activity`** bila `Local.Retro==1` — ⭐ **jalur retro keluar lebih awal** |
| 4 | `Obj-Browse`, gerbang `pyWorkPage.KomiteCount==1` |
| ⭐⭐ **5** | **`OperatorID.pyPosition=="SPV B"`** **DAN** **`Local.TotalAdj > 30.000.000,00 && <= 57.750.000,00`** |
| 6.2 | `Property-Remove` bila `KomiteList(1).KomiteID == .OPERATOR_ID` |
| ⚠️ 7.1 | mengisi `KomiteList(<APPEND>)` — ⚠️ **flag `false`, gerbangnya mati** |
| 8 | `pyWorkPage.KomiteLoop = @LengthOfPageList(KomiteList)` |

> ### ⭐⭐ PITA NILAI UANG YANG DI-HARDCODE
>
> `[terverifikasi]` **`Local.TotalAdj > 30000000.00 && Local.TotalAdj <= 57750000.00`** —
> dua ambang **tertulis langsung di dalam gerbang**, dipasangkan dengan **jabatan `"SPV B"`**.
>
> ⚠️ **Ini bentuk wewenang yang BELUM PERNAH terlihat di dua modul komite lain**, tempat seluruh
> medan hak akses kosong. ⛔ **Saya tidak menyimpulkan apa pun** — angkanya dicatat apa adanya.
>
> ⚠️ **Dan angka kedua ganjil:** `57.750.000` bukan bilangan bulat rapi. ⛔ Dari mana ia berasal
> **tidak tertulis**. `[terbuka]` **butir 15**

### 3c — `AcceptStatus`: **satu penulis, tujuh pembaca**

`[terverifikasi]` **Jendela: seluruh 114 berkas.**

| | Berkas |
| --- | --- |
| ⭐ **MENULIS** *(`PropertiesName`)* | ⭐ **`KomitePostAct`** — **satu**, sekali |
| **MEMBACA** *(nilai atau gerbang)* | `KomitePost_Adjustment` **27** · `KomitePost_Reject` **10** · `KomitePost_CloseClaim` **9** · `SendEmailKlaim_KMT` **5** · `SetDataForInformation_Act` **2** · `InsertJsonClaimNonMBU_act` **1** · ⭐ `When\IsKomiteLoop` **1** |

⭐ **Satu penulis, tujuh pembaca** — bentuk yang rapi. ⚠️ Bandingkan Komite Claim Prop, tempat
`AcceptStatus` juga ditulis **tepat satu** penugasan.

### 3d — ⭐ Vonis **ADR-0014**: ⚠️ **MENENTANG SEBAGIAN** — dan bedanya tajam dari Komite Claim Prop

| | Komite Claim Prop | ⭐ Komite Claim Fac In |
| --- | --- | --- |
| penempatan tugas | ada — `KomiteRouter` | ada — `KomiteRouter` |
| ⭐ **pemeriksaan pemilik giliran** | ⛔ **NOL** | ⭐ **ADA** — `Local.USERNAME == .KomiteID` |
| letak pemeriksaannya | — | ⚠️ **LAYAR** *(`pyDisabledWhen`)*, bukan lapisan layanan |
| pintu belakang | — | ⚠️ **`"IT Developer"` melewatinya lewat kode 5** |
| **vonis** | ⛔ **MENENTANG** | ⚠️ **MENENTANG SEBAGIAN** |

⭐ **Mekanismenya ADA dan niatnya terbaca** — itu lebih jauh daripada Komite Claim Prop.
⛔ **Tetapi ADR-0014 menuntut penolakan di LAPISAN LAYANAN**, dan yang ada hanya penonaktifan
kendali layar — yang **dapat dilewati siapa pun yang memanggil layanannya langsung**.

⛔ **Nol kesimpulan untuk Go ditarik dari sini.**

---

## §4 — Daur hidup diadu dengan Komite Claim Prop

### 4a — ⭐⭐ **IDENTIK, sampai ke nama penghubungnya**

`[terverifikasi]` Dibaca berdampingan, wadah demi wadah:

| | Komite Claim FacIn | Komite Claim Prop |
| --- | --- | --- |
| kelas kerja | `…Work-Komite` | `…Work-KomiteTreaty` |
| `END52` | Event-End | Event-End |
| `ASSIGNMENT63` | **`KomiteRouter`** · rute **Custom** · impl **WorkList** | **`KomiteRouter`** · rute **Custom** · impl **WorkList** |
| `Decision1` | **`KomiteLoop`** · **DataXOR** | **`KomiteLoop`** · **DataXOR** |
| `Start1` | Event-Start | Event-Start |

**Penghubung — keempatnya, berdampingan:**

| Pengenal | Dari → Ke | Nama | Sama? |
| --- | --- | --- | --- |
| `TRANSITION54` | `ASSIGNMENT63` → `Decision1` | **`ViewTransferDtl`** | ✅ |
| `Transition3` | `Decision1` → `END52` | **`NoLoop`** | ✅ |
| `Transition2` | `Decision1` → `ASSIGNMENT63` | **`IsKomiteLoop`** | ✅ |
| `Transition1` | `Start1` → `ASSIGNMENT63` | **`[Always]`** | ✅ |

⭐⭐⭐ **VONIS: SAMA.** Bentuk, pengenal, nama gerbang, arah penghubung, jenis rute, jenis
keputusan — **tidak satu pun berbeda**. ⛔ **Yang berbeda hanya kelas kerjanya.**

⭐ **Kedua modul komite adalah SATU RANCANGAN ALUR yang disalin dan diubah kelasnya.**
✅ Ini menutup **butir terbuka 7** ronde 1.

### 4b — `KomiteLoop`: tangga maju

`[terverifikasi]` Gerbang `Decision1` bernama **`KomiteLoop`**, berjenis **DataXOR**, dan
penghubung baliknya bernama **`IsKomiteLoop`** — ⭐ rule `When` yang **khas modul ini**.

⭐ **Kondisinya terbaca** *(lihat §6)*: ia menguji **`.AcceptStatus`** dan **`.KomiteCount`**.

⚠️ **Pencacahnya dinaikkan di mana?** `KomiteRouter` langkah 5 yang menaikkan `.KomiteCount`
⛔ **BER-REMARK — mati**. ⭐ **Jadi kenaikan pencacah terjadi di tempat lain**, dan ⛔ **belum saya
temukan**. `[terbuka]` **butir 16**

---

## §5 — Empat berkas produksi dan ruleset `ADESAMUEL@`

### 5a — Keempat berkas `pegaprdnusare`

`[terverifikasi]`

| Berkas | Identitas | Ruleset |
| --- | --- | --- |
| ⭐ `Activity\SaveAccept_ACT` | `…DATA-OBJECTITEM!SAVEACCEPT_ACT` | `GCNMFW 01-01-22` |
| ⭐ `Activity\SaveReject_ACT_KMT` | `…DATA-OBJECTITEM!SAVEREJECT_ACT_KMT` | `GCNMFW 01-01-22` |
| `ReportDefinition\GetEmailUser_RD` | `DATA-ADMIN-OPERATOR-ID!GETEMAILUSER_RD` | `GISFW 01-01-83` |
| `When\IsBondingKBG` | `@BASECLASS!ISBONDINGKBG` | `SFAGIS 01-01-31` |

### 5b — Diadu isinya dengan salinan di modul lain

⛔ **`md5` berkas TIDAK dipakai** — ia terbukti bukan alat yang sah. Yang dipakai: **sidik isi
dengan 27 medan waktu/ekspor dibuang.**

| Berkas | Pembanding | Hasil |
| --- | --- | --- |
| `GetEmailUser_RD` | Komite Claim Prop | ✅ **IDENTIK** |
| `IsBondingKBG` | Claim Fac In | ✅ **IDENTIK** |
| ⭐ `SaveAccept_ACT` | ⛔ **tidak ada pembanding** | khas modul ini |
| ⭐ `SaveReject_ACT_KMT` | ⛔ **tidak ada pembanding** | khas modul ini |

⭐ **Dua yang punya pembanding TIDAK berbeda isinya.** ⚠️ Jadi keputusan *"salinan produksi yang
ditiru"* **tidak mengubah apa pun untuk kedua berkas itu**. ✅ **Butir terbuka 3 ronde 1 DITUTUP.**

⭐ **Dua sisanya khas modul ini** — `SaveAccept_ACT` dan `SaveReject_ACT_KMT` — jadi ⛔ tidak ada
yang dapat diadu, dan ⭐ **keduanya menyimpan keputusan aksep/tolak**. ⛔ Isinya belum dibaca.
`[terbuka]` **butir 17**

### 5c — ⭐⭐ Ruleset `ADESAMUEL@`: **satu berkas, dan ia yang paling penting**

`[terverifikasi]` **`Section\ShowTransfer.xml`** · `…WORK-KOMITE!SHOWTRANSFER` · versi `01-01-01`.

> ### ⭐⭐ DAN ITU LAYAR YANG MEMBAWA SATU-SATUNYA PENEGAKAN WEWENANG
>
> `ShowTransfer` adalah **layar keputusan komite**, dan §2c membuktikan **ia satu-satunya tempat
> `ProteksiKomite.CARI1` dibaca** — yaitu **satu-satunya penegakan pemilik giliran di modul ini**.
>
> ⚠️ **Layar itu tinggal di ruleset bernama akun perorangan**, terpisah dari **56 berkas `GISFW`**
> dan **52 berkas `GCNMFW`**.
>
> ⛔ **Dicatat sebagai temuan tentang BERKAS, bukan tentang orang.** ⛔ Nol penilaian terhadap siapa
> pun, nol usulan kepegawaian. ⚠️ Yang perlu diketahui: **siapa pun yang kelak mencari layar
> keputusan komite di ruleset resmi tidak akan menemukannya.**

✅ **Butir terbuka 2 ronde 1 DITUTUP.**

⚠️ **Bandingkan Claim Fac In:** di sana ruleset yang sama memuat **berkas alur** — daur hidup
seluruh modul. ⭐ **Polanya berulang: yang tinggal di ruleset perorangan justru yang paling
menentukan.**

---

## §6 — `IsKomiteLoop` dan tujuh properti tanpa penulis

### 6a — ⛔⛔ **RALAT: ronde 1 SALAH — kondisinya terbaca**

`[terverifikasi]` Kondisi `IsKomiteLoop` **ADA dan terbaca**:

| Medan | Isi |
| --- | --- |
| `pyConditionString` | ⭐ **`Status Penerimaan = 1`** |
| `pyDesignatedProperty` / `pyUnmodifiedPath` | ⭐ **`.AcceptStatus`** dan **`.KomiteCount`** |
| `pySimpleCondition` | `false` |
| `pyClassName` | `ASM-FW-GCNMFW-Work-Komite` |

> ⛔ **RALAT terhadap ronde 1 §4.** Kalimat lamanya **dikutip utuh, tidak dihapus**:
> *"kondisi **terbaca** 48 dari 49 · ⛔ **tidak terbaca di tujuh medan** 1 — ⭐ **`IsKomiteLoop`**"*
> dan *"⛔ **Bukan berarti kosong** — ia menuntut medan kedelapan yang belum saya kenal."*
>
> ⭐ **Yang benar: 49 dari 49 kondisi terbaca.** ⚠️ **Sebabnya bukan medan yang kurang, melainkan
> PENYARING SAYA:** ronde 1 mengekstraksi properti dengan pola yang **mewajibkan awalan halaman**
> *(`pyWorkPage.` · `pyWorkCover.` · `pxProcess.` · `OperatorID.`)*. ⛔ `IsKomiteLoop` menguji
> **`.AcceptStatus`** — **jalur RELATIF, tanpa awalan halaman** — sehingga polanya tidak menangkap.
>
> ⭐ **Ini jebakan sensus yang sama bentuknya dengan dua kegagalan sebelumnya:** angka yang salah
> bukan karena korpusnya, melainkan karena **jendelanya**.

✅ **Butir terbuka 8 ronde 1 DITUTUP.**

### 6b — ⭐ Ketujuh properti tanpa penulis: **TIDAK mati**

`[terverifikasi]` Ditelusuri satu per satu:

| Properti | Kenapa tidak punya penulis di sini |
| --- | --- |
| `pyWorkPage.Quotation.BusinessCode` · `BusinessName` | ⭐ dibaca dari **kasus induk**; hanya `BusinessType` dan `GroupPanel` yang **disalin** ke halaman komite |
| ⭐ `pyWorkPage.OfferFacIn.QuotationData.BusinessType` | ⭐ **diuji langsung pada halaman induk** lewat `pyWorkCover.*` — tidak perlu penulis lokal |
| `pyWorkPage.pyWorkIDPrefix` · `pyWorkCover.pyWorkIDPrefix` | ⭐ **bawaan Pega** — diisi mesin, bukan rule |
| `pxProcess.pzProductionLevel` · `pxProcess.pxSystemNodeID` | ⭐ **bawaan Pega** — keadaan proses berjalan |

⭐ `[terverifikasi]` **`SetValueKomite` langkah 12** menyalin **tiga** nilai dari kasus induk:

```
pyWorkPage.Quotation.BusinessType  =  pyWorkCover.OfferFacIn.QuotationData.BusinessType
pyWorkPage.Quotation.GroupPanel    =  pyWorkCover.OfferFacIn.QuotationData.GroupPanel
pyWorkPage.CLMNO                   =  pyWorkCover.pyID
```

⭐⭐ **Modul ini membaca kasus induk secara langsung lewat `pyWorkCover` di banyak tempat** —
`HitServiceToKasirKMT_Act` *(17 titik)* · `SendEmailKlaim_KMT` *(13)* · `SendErrorDirectKasir` *(13)*
· `SetDataForInformation_Act` *(9)*. ⭐ **Jadi ia tidak menyalin seluruh data klaim; ia
MEMBACANYA saat dibutuhkan.**

✅ **Butir terbuka 9 ronde 1 DITUTUP.**

---

## §7 — Penutup

### 7.1 — ⚠️ Kesimpulan ronde ini yang **PALING RAWAN SALAH**

⚠️⚠️ **§3d — vonis ADR-0014 "MENENTANG SEBAGIAN".**

**Kenapa rawan:** ia bersandar pada **satu pembacaan** — bahwa `ProteksiKomite.CARI1` hanya dibaca
di `pyDisabledWhen` sebuah Section, jadi penegakannya *"hanya di layar"*. ⛔ **Saya menyisir nama
`ProteksiKomite` di seluruh 114 berkas**, tetapi ⚠️ **halaman itu bisa saja dibaca lewat nama lain**
— alias halaman, parameter, atau pemetaan di sisi layanan. ⭐ **Pola "medan keenam belas" sudah
berkali-kali menggigit proyek ini.**

**Yang kokoh dan tidak rawan:** ⭐ **§4 daur hidup identik** — dibaca berdampingan, entri demi entri,
dan **tidak satu pun berbeda**; ⭐ **§2b rantai gerbang** — ketiga syarat dan ketiga kodenya terbaca
langsung dari medan, bukan disimpulkan.

**Rawan kedua:** §1c *"tujuh pembangkit nomor per lini semuanya mati"*. Saya membaca **penanda
remark**, dan ⛔ **belum memeriksa apakah ada pembangkit per lini di rule LAIN** yang menggantikannya.

### 7.2 — ⭐ Pertanyaan untuk work owner — **dua**

❓ **Q1** — **Pintu belakang `"IT Developer"` pada wewenang keputusan komite**

`[terverifikasi]` `SetProteksiSubmiteKomite` langkah 2.1 menyetel *"boleh memutuskan"* bila
**salah satu** dari tiga terpenuhi, dan syarat pertamanya **`pyPosition=="IT Developer"`** memakai
**kode 5** sehingga **memotong dua syarat sisanya** — termasuk ⭐ **`Local.USERNAME == .KomiteID`**,
yaitu pemeriksaan pemilik giliran itu sendiri.

⚠️ **Ini BUKAN pintu belakang yang sama** dengan dua kemunculan sebelumnya *(yang mengatur
penghapusan catatan kronologi)*. ⭐ Yang ini mengatur **siapa boleh memutuskan uang**. Catatan
pengembangnya menunjukkan ia **ditambahkan belakangan**.

**Bedanya kalau A atau B.** **A — tidak dibawa:** wewenang di sistem baru **hanya** pemilik giliran,
tanpa pengecualian; bila pengembang perlu menembus, itu lewat **jalur administratif yang tercatat**,
bukan gerbang tersembunyi. **B — dibawa apa adanya:** pintu belakang tetap ada, dan **siapa pun yang
jabatannya diketik `"IT Developer"` dapat memutuskan klaim siapa pun**.

➡️ **Rekomendasi: A — jangan dibawa.** ⭐ Alasannya: ⚠️ ia **melewati satu-satunya penegakan
wewenang yang dipunyai modul ini**, dan ⛔ ia bergantung pada **teks jabatan** — data kepegawaian
yang berubah tanpa sepengetahuan sistem. ✅ Sejalan keputusan Q4 Claim Prop *(izin eksplisit lewat
peran, bukan teks jabatan)*. ⚠️ **Penyimpangan sadar**, dan **wajib dicatat sebagai lubang yang
ditutup, bukan fitur yang dibuang**.

---

❓ **Q2** — **Pita nilai uang yang di-hardcode pada `ApprovalKomite_Act`**

`[terverifikasi]` Langkah 5 menggerbangi **`OperatorID.pyPosition=="SPV B"`** bersama
**`Local.TotalAdj > 30.000.000,00 && <= 57.750.000,00`** — dua ambang **tertulis langsung di dalam
gerbang**.

⚠️ Angka kedua **ganjil** — `57.750.000` — dan ⛔ **dari mana ia berasal tidak tertulis di korpus**.

**Bedanya kalau A atau B.** **A — pita nilai menjadi DATA**, bukan kode: disimpan sebagai baris
tabel wewenang yang dapat diubah tanpa menyentuh kode. **B — ditiru apa adanya** sebagai tetapan.

➡️ **Rekomendasi: A — jadikan data.** ⭐ Alasannya: ambang uang **berubah mengikuti kebijakan**, dan
⚠️ **tetapan di dalam kode berarti setiap perubahan kebijakan menuntut rilis**. ✅ Preseden ada:
Claim Prop menjadikan roster wewenang sebagai **baris tabel**, dengan catatan *"mengganti pemegang
jabatan cukup dengan mengubah baris tabel"*. ⛔ **Nilai awalnya tetap yang tercatat di korpus** —
30.000.000 dan 57.750.000 — ⚠️ **dan asal-usul angka kedua wajib ditanyakan ke work owner.**

### 7.3 — Register `[terbuka]` sesudah ronde 2

**✅ DITUTUP — enam:**

| # | Butir | Ditutup oleh |
| --- | --- | --- |
| **2** | Berkas mana di ruleset `ADESAMUEL@` | §5c — `Section\ShowTransfer` |
| **3** | Isolasi 4 berkas `pegaprdnusare` | §5a · §5b — dua identik, dua khas |
| **5** | Berapa `Obj-*` punya jalur kegagalan | §1d — **NOL** dari 25 |
| **7** | Apakah daur hidup kedua modul sama | §4a — **SAMA**, sampai nama penghubung |
| **8** | Kondisi `IsKomiteLoop` | §6a — terbaca; ⛔ **ronde 1 salah** |
| **9** | Apakah 7 properti diisi dari luar | §6b — dibaca dari induk atau bawaan Pega |

**⭐ BARU — empat:**

| # | Butir | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| **14** | Apakah `KomitePost_Survey` ada berkasnya — jalur `TransferType=="1"` mati | asisten | tidak |
| ⭐ **15** | **Asal-usul ambang `57.750.000`** | work owner | ⚠️ **ya** — Q2 |
| **16** | Di mana `KomiteCount` dinaikkan, sebab langkah penaiknya ber-remark | asisten | ⚠️ **ya** — tangga tak terbaca |
| **17** | Isi `SaveAccept_ACT` dan `SaveReject_ACT_KMT` | asisten | tidak |

**Sisa ronde 1 yang masih terbuka: 1 · 4 · 6 · 10 · 11 · 12 · 13.**

⚠️ **Butir 18 yang sempat saya siapkan — *"cacah pasti rule ber-`COMMIT` bersama"* — TIDAK JADI
LAHIR:** ia dijawab di §8b pada blok yang sama, **9 dari 9**. ⛔ Dicatat supaya nomor 18 tidak
dipakai ulang untuk butir lain.

⭐ **Total: 13 − 6 + 4 = 11 butir terbuka.** **Dihitung dua cara, dan keduanya sepakat:**
*(a)* mencacah baris ketiga tabel di atas ⇒ `7 sisa + 4 baru = 11`;
*(b)* delta ⇒ `13 − 6 ditutup + 4 baru = 11`.

⚠️ **Butir 4 ronde 1** *(isi `SetProteksiSubmiteKomite`)* **tidak ditutup melainkan NAIK menjadi
pertanyaan Q1** — isinya sudah dibaca, **keputusannya belum diambil**.

### 7.4 — ⛔ RALAT terhadap ronde 1 — **dua**

| # | Yang diralat | Kalimat LAMA, dikutip | Yang benar |
| --- | --- | --- | --- |
| ⭐ **1** | kondisi `IsKomiteLoop` | *"kondisi terbaca **48 dari 49** · tidak terbaca **1** — `IsKomiteLoop` … ia menuntut medan kedelapan yang belum saya kenal"* | ⭐ **49 dari 49 terbaca.** Sebabnya **penyaring saya** yang mewajibkan awalan halaman, bukan medan yang kurang |
| **2** | sifat pintu belakang | ronde 1 §2c menyiratkan `"IT Developer"` adalah **kemunculan ketiga pola yang sama** | ⚠️ **BUKAN sama.** Dua sebelumnya mengatur **catatan kronologi**; yang ini **wewenang keputusan komite**. ⭐ Pola kronologi **juga ada di sini**, di berkas terpisah |

### 7.5 — Yang paling menahan untuk ronde 3

| # | Yang dikerjakan | Kenapa menahan |
| --- | --- | --- |
| ⭐ **1** | **Baca `SaveAccept_ACT` dan `SaveReject_ACT_KMT`** *(butir 17)* | dua rule **khas** yang **menyimpan keputusan**, dan keduanya dari **salinan produksi** |
| ⭐ **2** | **Temukan di mana `KomiteCount` dinaikkan** *(butir 16)* | ⚠️ tanpa itu **tangga penyetuju tidak terbaca** |
| **3** | **Sisir `ProteksiKomite` lewat nama lain** *(kerawanan §7.1)* | vonis ADR-0014 bergantung padanya |
| **4** | **Periksa apakah ada pembangkit nomor per lini di rule lain** | kerawanan kedua §7.1 |
| **5** | **Baca 5 Section dan 5 FlowAction** | layar belum dibaca isinya, baru dicacah gerbangnya |
| **6** | **Adu 47 rule `When` bersama dengan Claim Fac In** | ⛔ sengaja dilewati ronde 1 dan 2 |

---

## §8 — ✅ Dua keputusan work owner — 2026-09-19

`[keputusan work owner — atas rekomendasi asisten]` **2026-09-19**. ⛔ Tanda *"atas rekomendasi
asisten"* dipakai **sengaja**: pilihannya datang dari rekomendasi saya, dan **pencabutannya harus
murah**.

### 8a — Q1 → butir terbuka **1** — ⚠️ **MENYEMPIT, tidak tertutup**

> ⭐ **`ASM-FW-GCNMFW-Work-Komite` diperlakukan sebagai kelas kerja KHUSUS Claim Fac In**, sampai
> terbukti sebaliknya.

⚠️ **Alasannya:** ekspor ini **hanya memuat Fac In**; menganggapnya generik berarti merancang kolom
untuk lini yang **belum pernah terlihat**.

⛔ **Butir 1 TIDAK ditutup** — satu pandangan DBA ke tabel produksi dapat membalikkannya, dan **itu
belum dilakukan**.

⚠️ **Dan ronde ini menambah satu alasan untuk berhati-hati:** §4a membuktikan alur kedua modul
komite **identik sampai ke nama penghubungnya**. ⭐ Rancangan yang disalin seperti itu **bisa saja
dimaksudkan generik sejak awal** — ⛔ **bukan bukti, hanya alasan untuk tidak menutup butirnya.**

### 8b — Q2 → keputusan baru — **satu aksi, satu transaksi**

> ⭐ **Kesembilan `COMMIT;` di dalam SQL tidak direplikasi.** ⚠️ **Kecuali pembangkit nomor urut**,
> yang nomornya **harus bertahan walau transaksi induknya batal** — kalau ikut dibatalkan, nomor
> yang sama **terpakai ulang**.

**Alasan yang mengikat.** ⛔ **Klaim brief bahwa "empat rule ber-`COMMIT` identitas empat bagiannya
sama di kedua modul" — SAYA VERIFIKASI ULANG, dan hasilnya LEBIH KUAT daripada klaimnya:**

`[terverifikasi]` **Jendela: medan `pyBrowseSQL` pada seluruh berkas `RDBList` ketiga modul,
dibandingkan lewat identitas EMPAT BAGIAN.**

| | Jumlah |
| --- | ---: |
| rule ber-`COMMIT;` Komite Claim FacIn | **9** |
| rule ber-`COMMIT;` Claim Fac In | **9** |
| rule ber-`COMMIT;` Komite Claim Prop | **10** |
| ⭐ **identik empat bagian dengan Komite Claim Prop** | ⭐ **9 dari 9** |
| identik empat bagian dengan Claim Fac In | **6 dari 9** |
| ⭐⭐ **identik dengan SALAH SATU modul saudara** | ⭐⭐ **9 dari 9** |

⭐⭐ **SELURUH sembilan rule ber-`COMMIT;` modul ini adalah rule bersama** — bukan empat.
**Kesembilan yang identik dengan Komite Claim Prop:** `GetSequenceNumber_SQL` · `GetTokenStorage_SQL`
· `InsertClaimPNC` · `InsertClaimRejected_Sql` · `InsertHistoryAkseptasiPega_Sql` ·
`InsertLOGDirectKasir_SQL` · `InsertLogServiceClaim` · `Insert_T_Storage_SQL` · `SaveOSClaim_SQL`.

⛔ **Karena itu "satu rule dua nasib" bukan risiko sebagian, melainkan MUTLAK:** memutuskan berbeda
untuk modul ini berarti **sembilan rule yang benar-benar satu diperlakukan dua cara**.

> ⛔ **RALAT terhadap brief blok ini.** Kalimatnya dikutip: *"⭐ **empat rule ber-`COMMIT` itu
> identitas EMPAT BAGIANNYA SAMA di kedua modul** — `GetSequenceNumber_SQL` · `SaveOSClaim_SQL` ·
> `InsertClaimPNC` · `GetTokenStorage_SQL`."* ⭐ **Yang benar: sembilan, bukan empat** — dan
> keempat yang disebut memang termasuk di dalamnya. ⚠️ **Angkanya terlalu kecil, bukan terlalu
> besar** — alasannya justru menguat.

⭐ **Alasan kedua yang berdiri sendiri:** **`COMMIT;` di tengah pekerjaan orang lain membuat
kegagalan separuh jalan meninggalkan data separuh tersimpan**, dan ⚠️ **§1d membuktikan modul ini
punya NOL jalur kegagalan pada 25 langkah `Obj-*`** — jadi **tidak ada apa pun yang akan
menangkapnya**.

✅ Sejalan keputusan **K1 Claim Fac In** 2026-09-19 dan **ADR-0006**.

⚠️ **Dan satu akibat yang wajib disebut:** modul ini **nol metode `Commit` Pega**, sehingga
menghapus kesembilan `COMMIT;` berarti **penutupan transaksi harus dibangun dari nol** di sistem
baru — bukan sekadar tidak menyalin sesuatu.

---

## Lampiran — bukti berkas lain tidak disentuh

```
korpus Komite Claim FacIn   114 berkas .xml   — nol dibuka untuk ditulis
korpus Claim Fac In         482 berkas .xml   — nol dibuka untuk ditulis
korpus Komite Claim Prop     80 berkas .xml   — nol dibuka untuk ditulis
docs\adr\                    15 berkas        — nol disunting
grilling-ronde-1.md                           — nol disunting
.scratch\claim-facin\         7 berkas        — nol disunting
.scratch\claim-prop\         27 berkas        — nol disunting
.scratch\komite-claim-prop\  27 berkas        — nol disunting
CLAUDE.md                                     — nol disunting
```

⛔ **Nol kode Go/React · nol `CREATE TABLE` · nol DDL · nol nomor baris XML dikutip · nol nilai
kredensial dibaca atau disalin · nol keputusan work owner baru di luar §8 · nol ADR direvisi ·
nol kesimpulan untuk aplikasi Go.**
