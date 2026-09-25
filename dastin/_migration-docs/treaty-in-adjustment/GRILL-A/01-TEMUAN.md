> Modul  : Treaty In Adjustment · Ronde A · 2026-09-24
> Peran  : interogator
> Masukan: 379 XML `Treaty In Adjustment/` + 329 XML `Treaty In/` · `PENGETAHUAN.md` · perkakas `tools/`
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 01 · TEMUAN RONDE A

## 1. Tabel temuan

| # | Temuan | Sisi | Mengubah |
|---|---|---|---|
| **NA-01** | Jenis addendum **dan** materialitasnya dipilih pengguna lewat dua radio group di picker Revisi — tanpa aturan, tanpa jejak, tanpa gerbang | 56-KHAS | TDA-13, §3.2, §4.2, §4.3 |
| **NA-02** | Penomoran revisi patah pada revisi kesepuluh; sebabnya **offset yang meleset satu**, bukan dua pengurai | 56-KHAS | TDA-12, §4.4, UA-1 |
| **NA-03** | Irisan kedua ekspor **317**, bukan 323; enam jalur memuat versi aturan berbeda | — | §1.2, `METODE` §8.1 |
| **NA-04** | **Sepuluh dari lima belas TDA berada di irisan** — ia cacat Treaty In yang berjalan hari ini | — | Lacak TDA |
| **NA-05** | `ValueDifference` punya **33 akar / 165 jalur daun**; 21 akar di antaranya agregat, dan hanya tujuh punya penulis khas addendum | IRISAN | §2, §5.6, §10.4 |
| **NA-06** | **TDA-07 dicabut**: membuka kontrak dari picker Adjustment **tidak** mengubah baris kontrak | IRISAN | TDA-07, TDA-08, §4.5, §7.2 |
| **NA-07** | Tiga dari lima field lapisan beku ADR-0040 dapat disunting di layar addendum, **tanpa syarat apa pun** | 56-KHAS + IRISAN | GRL-05, `UA-10` |
| **NA-08** | Dua dari tiga kontrol `Ceding` **tidak menulis `CedingID`** — nama dan pengenal cedant dapat tidak sinkron, dan ketidaksinkronan itu **tersimpan ke tabel datar** | 56-KHAS + IRISAN | §4.6, `UA-10`, titipan I/F/orang |
| **NA-09** | Dua hitungan sebelumnya ikut menghitung **salinan terbungkus**; satu di antaranya membalik tiga baris teratas sebuah tabel | — | §4.3, §4.5, `MA-04` |
| **NA-10** | Daftar **tertutup** setiap penulis `TreatyIn.Position`: lima nilai, dan `ReasTreatyInGroupLeader` bukan salah satunya. **Bukan temuan baru** — ADR-0055 §4.1 sudah menetapkannya 23 Sep; yang baru adalah **jejak `"BERNARD"`** | IRISAN | GRL-07, ADR-0052, ADR-0055, `UA-12` |
| **NA-11** | Pintu samping ADR-0055 §4 **juga ada di layar addendum**, dan dua di antaranya berperilaku lain dari yang tercatat: `Force Resolve Complete(dev)` **menyimpan** tanpa membereskan `RevisionState`; `TreatyInSetToDirector` **seluruh langkahnya mati** | IRISAN | §7.5, ADR-0055 §4, `UA-9`, cabang D |
| **NA-12** | Daur hidup `RevisionState` terbaca lengkap; **penolakan tidak mengosongkannya**, dan pewarisannya ke addendum terbukti ujung ke ujung — termasuk akibat kedua, addendum lahir **terkunci** | IRISAN | §7.4, TDA-08, `UA-9` |
| **NA-13** | **TDA-11 dikuatkan dengan sumber yang benar**: grid picker bersumber daftar halaman tanpa `WHERE`; konfigurasi RD saingan yang menyaring `Resolve Complete` tidak dapat hidup | 56-KHAS | §7.6, TDA-11, cabang B |
| **NA-14** | **Audit `TA-04`**: lima jenis penulis dikalibrasi lalu disapu ulang. Tiga klaim negatif **berdiri**, tiga angka **keliru**, satu daftar "tertutup" **bocor satu penulis** | — | §4.3, §4.6, §6.2, §7.4, `MA-05` |
| **NA-15** | `Force Resolve Complete(dev)` di layar addendum memanggil prosedur **kontrak**; `UPDATE` tidak mengenai baris, **tidak ada galat**, dan layar melapor *"Data Sudah Disimpan"* | aturan IRISAN, keterjangkauan **56-KHAS** | §6.4, `METODE` §2.2 |
| **NA-16** | Keterlihatan tombol "(dev)" = **wadah × sel**: dua tombol untuk **dua orang bernama di divisi IT**, satu (`Force Edit`) untuk **setiap** operator divisi IT | IRISAN | §7.5, `MA-07`, ADR-0055 |
| **NA-17a** | `CedingID` dapat diketik **teks bebas tanpa validasi** di `ShowSummary` — satu-satunya penyuntingnya di seluruh korpus; **tetapi layarnya dijaga satu nama operator** | IRISAN | §4.6, `IND-1`, `UA-10` |
| **NA-17b** | Tiga field lapisan beku dapat disunting **tanpa syarat keadaan apa pun** di `InputTreatyInOffer` dan `InputTreatyInAdjustment` — termasuk sesudah kontrak yang sudah disetujui dibuka kembali | IRISAN | GRL-05, `UA-10`, `IND-1` |
| **NA-18** | Keempat kontrol `SetTreatyIn_Act` **dibatasi peran inputor** di tingkat sel; hanya **`Revision`** mengirim `revisionstate=1`, dan hanya ia bekerja atas kontrak ber-`Resolve Complete` | IRISAN | §4.5, `MA-08`, diff eskalasi butir 1 |
| **NA-19** | `Save EDM(dev)` menyimpan addendum **tanpa syarat status apa pun**, dan bersama `Force Edit (dev)` memberi operator divisi IT jalan menyunting-lalu-menyimpan addendum yang **sudah disetujui** | IRISAN | GRL-02, `REV-1`, `UA-10` |
| **NA-21** | `Revision` mengirim **dua** parameter: ia memendekkan rantai **dan** mengunci sebagian kontrol lewat `ViewState = 1` | IRISAN | §7.7, `DB-10`, GRL-07 |
| **NA-22** | Kunci itu **tidak konsisten**: dicocokkan ke 33 akar §5.6, **tidak satu pun dari dua belas akar uang terkunci seluruhnya** — dan `Layers` menjaga `.Deductible` sedangkan kembarannya `LayersEDM` tidak | IRISAN | §7.7, `DB-10`, GRL-07 |
| **NA-20** | Pemicu penggeser `Termination` otomatis **tidak terjangkau**: seksinya ada di layar addendum, tetapi sel pemicunya `pyReadOnly = true`, dan sel tanggal yang dapat disunting **tidak punya peristiwa** | IRISAN | GRL-05, §4.6 |

## 2. Rincian

### NA-01 · Jenis dan materialitas dipilih di layar, bukan ditetapkan aturan

`Section/PickerTreatyInMasterRevisi.xml` memuat dua kontrol `pxRadioButtons` yang terikat langsung
ke `TreatyIn.EDMState` dan `TreatyIn.EDMMaterialType`, keduanya `pyReadOnly = false` dan
`pyDisabledNew = false`. Kondisi tampil diperiksa sampai akar seksi — sel, baris, tabel,
`Embed-Harness-SectionBody`, `Embed-Harness-Section`: **nol** `pyVisible`, `pyCondition`,
`pyVisibleWhen`. `PickerTreatyInMaster` tidak memuat keduanya.

Maka `TreatyCreateEDM` hanya **menyemai**; nilai yang terkirim ke `TreatyInEDMSetValue` adalah
nilai radio saat baris dipilih.

**Kadar.** Mekanismenya **terbaca**. Himpunan nilai yang ditawarkan **tidak terbaca**
(`pyListSource = associated`; aturan Field Value tidak ter-ekspor). Maka `EDMState = "2"`
berstatus **"masih bisa terjadi bila Field Value menawarkan 2"** — bukan terbukti. Dan "penyesuaian
premi selalu material" **terbaca hanya untuk jalur tombol Penyesuaian**.

### NA-02 · Penomoran revisi: offset meleset satu, dua kemungkinan akhir

Cacatnya pada `Activity/TreatyInRevisi_post.xml` **baris 848**, yang membaca mulai **indeks 10** —
digit satuan saja — pada pengenal yang digit puluhannya di indeks 9.

`RDBList/GetTreatyRevisionID.xml` memakai `SUBSTR(ID,10,2)` yang justru membaca kedua digit dengan
benar, tetapi **nilainya dibuang**: `HASIL1` hanya diuji kosong atau tidak. Jadi menyebut "dua
pengurai berbeda" sebagai **sebab** adalah salah, dan dicabut.

Pengenal `1234567/Rnn` panjangnya 11 karakter, sehingga `@substring(ID,10,12)` selalu meminta batas
akhir di luar teks — perilaku yang **tidak terbaca**. Dua akhir yang mungkin:

| Bila `@substring` … | Akibat | Sifat |
|---|---|---|
| memotong longgar | `/R09` -> `/R010` -> `/R11` -> **`/R02`**, menimpa revisi kedua lewat cabang `UPDATE` yang melapor berhasil | **diam** |
| gagal keras | revisi-dari-revisi selalu galat | **berisik** |

Yang memutuskan: **UA-1**, diperluas menjadi tiga kueri.

### NA-03 · Irisan 317, dan enam versi aturan yang berbeda

Diperiksa per `pzInsKey`, bukan per nama berkas (`METODE` §2.0). Dari 323 jalur bersama, **317**
ber-`pzInsKey` sama. Enam yang berbeda:

| Jalur | Adjustment | Treaty In |
|---|---|---|
| `Activity/AddSpreadingXOL.xml` | 01-01-96 | 01-01-95 |
| `Activity/LoadAttachmentData.xml` | 01-01-54 | 01-01-62 |
| `Activity/SetCurrName_Act.xml` | 01-01-56 | 01-01-94 |
| `Activity/TreatyInSubmit.xml` | 01-01-90 | 01-01-96 |
| `Activity/AddDeduction.xml` | 01-01-54 | 01-01-54 (kunci beda, versi sama) |
| `FlowAction/CoBList.xml` | 01-01-55 | 01-01-55 (idem) |

**Tidak satu pun** dari keenamnya menjadi bukti TDA-01 sampai TDA-15, sehingga tidak ada temuan
yang bergantung pada pilihan versi. `PENGETAHUAN.md` membaca versi dari ekspor Adjustment.

### NA-04 · Sepuluh dari lima belas TDA ada di irisan

Konsekuensinya diatur `METODE` §8.2: temuan di irisan adalah temuan **Treaty In**, diadili sekarang,
dan — sesuai GRL-01 butir 5 — dicatat sebagai usulan untuk langkah 8-10 induk. Sesi ini memutuskan
nasibnya **untuk jalur addendum saja**. Daftarnya ada di kolom sisi `KEPUTUSAN-GRILLING-ADJUSTMENT.md`.

### NA-05 · 33 akar `ValueDifference`, 165 jalur daun

Dihitung ulang atas ekspor Adjustment; daftarnya di `PENGETAHUAN.md` §5.6 (`METODE` §3.9 — hitungan
bukan daftar). Dua hal yang hanya terbaca dari daftarnya: **21 dari 33 akar adalah agregat
`Total*`**, dan hanya **tujuh** akar yang punya penulis ber-awalan `TreatyEDM*`.

### NA-06 · TDA-07 dicabut — sesudah penyisiran tertutup

Versi pertama menyatakan memilih kontrak di picker Adjustment menulis kembali ke baris kontrak.
**Salah.** Pencabutan adalah arah **membuang**, jadi ia dituntut penyisiran tertutup (`METODE`
§3.2); hasilnya ada di `PENGETAHUAN.md` §4.5 dan diringkas di sini.

**Syaratnya aktif, bukan sisa teks.** Keenam langkah penulis (6 s.d. 11 `SetTreatyIn_Act`, baris
1188/1392/1580/1767/1933/2117) ber-`pyStepsPreCondition = 'true'`, tanpa blok mati, dengan
**bila salah = 3 (LEWATI)**. Jebakan §6.2 — teks syarat yang tertulis tetapi nonaktif — **tidak**
berlaku di sini; diperiksa nilainya satu per satu.

**Penyisirannya tertutup.** Setiap kemunculan nama di kedua ekspor digolongkan menurut tag
pembawanya: `pyActivity` 16/16 (action set layar), `pyStepsActivityName` 2/3 (langkah `Call`),
`pxStepDefaultDescription` 2/3 (teks bawaan langkah yang sama), sisanya metadata aturan itu
sendiri. **Nol** pada pra/pasca-proses Flow Action maupun jalur lain.

**Hanya kontrol 7 dan 8 `InputTreatyInOffer` yang mengirim `revisionstate=1`**, dan keduanya
**tidak dapat dicapai dari layar addendum**: `Section/InputTreatyInAdjustment.xml` dan harness-nya
tidak menyertakan `InputTreatyInOffer` sama sekali.

**Yang tersisa**, dan ia temuan Treaty In: kedua kontrol itu `pxButton`, `pyReadOnly = false`,
**tanpa satu pun** `pyVisible`, `pyCondition`, `pyVisibleWhen`, `pyDisabledWhen`, atau
`pyPrivilege` sampai akar seksi. Mereka mengosongkan status akseptasi **dan menyimpan** dalam satu
tindakan. Lihat `07-AUDIT` §3 — ia mengoreksi kalimat peredam pada butir 1 daftar eskalasi induk.

**Ikutannya:** klaim §7.2 lama ikut dicabut, dan TDA-08 ditulis ulang.

### NA-07 · Lapisan beku terbuka di layar addendum

Daftarnya, diperiksa sampai akar seksi, ada di `06-PUTUSAN` GRL-05 butir (i). Ringkasnya: `Ceding`,
`Commencement`, dan `Termination` dapat disunting pada pohon layar addendum **tanpa satu pun**
`pyVisible`, `pyCondition`, `pyVisibleWhen`, atau `pyDisabledWhen`. Karena tidak ada syarat sama
sekali, kedua nilai `EDMMaterialType` terjawab sekaligus.

`ProportionType` **dinyatakan tidak dapat dicapai** dari layar addendum: satu-satunya kontrol yang
dapat disuntingnya ada di `Section/ShowSummary.xml`, dan `ShowSummary` dirujuk hanya oleh
`InputTreatyInOffer` beserta harness dan flow action-nya.

### NA-08 · Menyunting nama cedant tidak selalu memperbarui ID-nya

Sel autocomplete `Ceding` di `TreatyInNONProportional` memuat `pyPropertyTarget = TreatyIn.CedingID`
— ia menulis keduanya. Dua kontrol `Ceding` di `InputTreatyInAdjustment` **nol** rujukan `CedingID`
di seluruh subpohon selnya. Satu-satunya aturan penulis `CedingID`, `TreatyInMappingDataconvert`,
dipanggil hanya oleh `TreatyInConvertCallData_act` dan bukan jalur penyuntingan layar.

**Ketiadaan kontrol untuk `CedingID` karena itu bukan bukti ia tidak pernah menyimpang** — ia justru
bentuk penyimpangannya (`METODE` §2.0). Diukur `UA-10`, yang kini memisahkan nama dari ID.

### NA-09 · Dua hitungan ikut menghitung salinan terbungkus

Ekspor Pega membungkus salinan aturan lain di dalam `pyIncludedRuleXML`, dan daftar versi lama
aturan yang sama di dalam `pyRuleVersionsList` (§1.2). Dua hitungan ronde ini ikut menghitungnya.

| Hitungan | Dilaporkan | Badan aturan | Salinan |
|---|---|---|---|
| `pyActivity = SetTreatyIn_Act` di `Section/InputTreatyInOffer.xml` | 8 | **8** | 0 |
| idem di `Harness/InputTreatyInOffer.xml` | 8 | **0** | 8 |
| `pyDisabledWhen` ber-`EDMMaterialType` di `Section/TreatyInNONProportional.xml` | 70 | **13** | 63 |
| idem di `Section/InputTreatyInOffer.xml` | 70 | **0** | 76 |
| idem di `Harness/InputTreatyInOffer.xml` | 70 | **0** | 76 |

**Kesimpulan kedua temuan tidak berubah**, dan yang kedua justru menguat: dihitung dari badan saja,
kondisi materialitas berjumlah **220 di 18 seksi**, bukan "ratusan" yang sebagian palsu. Yang
berubah adalah **di mana** kondisi itu berada: `InputTreatyInOffer` ternyata **tidak memuat satu
pun**.

**Bukti temuan lain diperiksa ulang, dan semuanya bersih** — nol salinan:

| Temuan | Bukti | Badan | Salinan |
|---|---|---|---|
| NA-01 radio picker | `pyValue = TreatyIn.EDMState` / `…EDMMaterialType` di `PickerTreatyInMasterRevisi` | 3 + 3 | 0 |
| GRL-05 lapisan beku | kontrol `Ceding`, `Commencement`, `Termination` di `InputTreatyInAdjustment` | 2 + 2 + 2 | 0 |
| GRL-05 lapisan beku | kontrol `Ceding` di `TreatyInNONProportional` | 2 | 0 |
| NA-08 | `pyPropertyTarget = TreatyIn.CedingID` di `TreatyInNONProportional` | 1 | 0 |
| NA-06 kontrol 7-8 | `pyActivity` di badan `Section/InputTreatyInOffer.xml` | 8 | 0 |

**Tidak ada temuan yang dikoreksi diam-diam.**

### NA-10 · Peran kelompok tidak dapat dimasuki — daftar penulis yang tertutup

> **Atribusi, dan ia harus disebut lebih dulu.** Ini **bukan** temuan baru. **ADR-0055 §4 butir 1**,
> diterima 23 September 2026, sudah menyatakannya dengan kalimat yang hampir sama: *"Sapuan
> menyeluruh atas ekspor menunjukkan tidak ada satu pun aturan yang pernah menyetel `Position` ke
> nilai itu."* Saya menyusun daftar ini tanpa membaca ADR-0055 lebih dulu, padahal ADR itu ada di
> daftar rekonsiliasi cabang A saya sendiri (`MA-06`).
>
> Yang tetap bernilai: daftarnya **dituangkan** — hitungan bukan daftar — sehingga dapat
> diverifikasi ulang; dan **satu hal yang benar-benar baru** muncul di ujungnya, yaitu jejak nama
> `"BERNARD"` yang membuat `UA-12` dapat mengukur sejarah, bukan hanya keadaan sekarang.

Pertanyaannya bukan "adakah yang menyebut `ReasTreatyInGroupLeader`" — menyebut bukan menulis
(`METODE` §2.0a). Pertanyaannya: **nilai apa saja yang dapat masuk ke `TreatyIn.Position`.**
Didaftar, bukan dicari:

| Aturan | Langkah | Nilai yang ditulis |
|---|---|---|
| `Akseptasi_DT` | 1.1.1.1 – 1.1.4.1 | `"ReasTreatyInSecHead"` |
| `Akseptasi_DT` | 1.2.1.1 | `"ReasTreatyInDeptHead"` |
| `Akseptasi_DT` | 1.3.1.1 | `"ReasTreatyInDirector"` |
| `Akseptasi_DT` | 1.2.2.1 / 1.3.2.1 / 1.5.2.1 / 2.2.2.1 | `"ReasTreatyInAdmin"` |
| `Akseptasi_DT` | 1.2.3.1 / 1.3.3.1 / 1.5.1.1 / 1.5.3.1 / 2.2.1.1 / 2.2.3.1 | `""` |
| `Akseptasi_DT` | 2.1.1 | `"ReasTreatyInSecHead"` |
| `Akseptasi_DT` | **1.4.1.1 / 1.4.2.1 / 1.4.3.1** | `"ReasTreatyInDirector"` / `"ReasTreatyInAdmin"` / `""` — **seluruh cabang 1.4 `pyDisabled = true`** |
| `TreatyInAddNew` | 1 | `"ReasTreatyInAdmin"` |
| `TreatyInSetEdit` | 3 | `"ReasTreatyInAdmin"` |
| `TreatyInReturntoInputor` | 2 | `"ReasTreatyInAdmin"` |
| `TreatyInForceResolveComplete` | 2 | `""` |

**Lima nilai yang mungkin: `ReasTreatyInAdmin`, `ReasTreatyInSecHead`, `ReasTreatyInDeptHead`,
`ReasTreatyInDirector`, dan kosong.** Tidak ada aturan di kedua ekspor yang menulis
`ReasTreatyInGroupLeader`. Jadi keadaan itu **tidak dapat dimasuki** — bukan karena tidak ada yang
menyebutnya, melainkan karena tidak ada nilai yang mengantarnya.

**Enam tempat yang menangani keadaan itu, semuanya pembaca:**

| Berkas | Perannya | Keadaan |
|---|---|---|
| `Akseptasi_DT` cabang 1.4 | perpindahan **dari** keadaan itu | **MATI** — `pyDisabled = true` |
| `TreatyInSetToDirector` langkah 4.5 | precondition atas keadaan itu | hidup; langkah 1-4 induknya mati |
| `TreatyInSetValue` langkah 2.6 | precondition; menulis `PositionUsername = "BERNARD"` | hidup, kondisinya tak pernah benar |
| `Section/InputTreatyInOffer` | `pyCondition = 1==2 && (…)` | **tidak pernah benar** (`METODE` §2.0) |
| `Section/TreatyInActionButtons` | **kondisi TAMPIL tombol kolom "Actions"**, membaca `TreatyIn.Position` | hidup — **pembaca, bukan penulis** |
| `When/IsTreatyUser` | membaca **workbasket operator**, bukan `TreatyIn.Position` | hidup |

Dua hal yang mudah tertukar dan harus dipisah: **workbasket** `ReasTreatyInGroupLeader` ada, dan
orang boleh menjadi anggotanya; **perkara** tidak pernah diarahkan ke sana, karena pengarahannya
lewat `Position`. `TreatyInActionButtons` menuntut `OperatorID.pyWorkBasketList(2).pyWorkBasketName
= TreatyIn.Position` — kecocokan yang, untuk nilai ini, tidak akan pernah terjadi.

Yang tersisa bukan jalan masuk, melainkan **jejak**: `TreatyInSetValue` langkah 2.6 menanam nama
orang, `"BERNARD"`, persis seperti `"IRVANDY"`, `"YOHANESKRISTIAWAN"`, dan `"NANDINA"` pada tingkat
lain. Nama itu satu-satunya cara mengenali baris yang pernah melewati tingkat kelompok — dan itulah
sebabnya `UA-12` diperluas ke sejarah.

### NA-11 · Pintu samping ADR-0055 §4 — dua koreksi dan satu perluasan

> **Atribusi.** Keempat pintu samping sudah didaftar **ADR-0055 §4**, termasuk catatan bahwa
> `TreatyInForceEdit` *"terlihat oleh setiap pengguna layar penawaran"*. Itu bukan temuan ronde ini
> (`MA-06`). Yang berikut ini **menambah atau mengoreksi** daftar tersebut.

**Perluasan — pintunya ada di layar addendum juga.** `Section/TreatyInActionButtons` disertakan
oleh `Section/InputTreatyInOffer` **dan** `Section/InputTreatyInAdjustment`, masing-masing empat
kemunculan di badan aturan. ADR-0055 menyebut "layar penawaran" saja.

**Koreksi 1 — `Force Resolve Complete(dev)` menyimpan, dan meninggalkan `RevisionState`.** Selnya
membawa `pyActivity = SaveTreatyIn_Act`; `TreatyInForceResolveComplete` menulis `StatusAkseptasi =
"Resolve Complete"`, `Position = ""`, `PositionUsername = ""` — dan **tidak menyentuh
`RevisionState`**. Ia satu-satunya aturan yang dapat melahirkan keadaan beku §7.4 butir 4, dan
karena itu satu-satunya cara sebuah kontrak dapat tampak selesai sambil tetap bertanda sedang
direvisi.

**Koreksi 2 — pintu keempat terkunci dari dalam.** ADR-0055 mencatat `TreatyInSetToDirector`
"melompati dua tingkat", terjangkau lewat `TreatyInTestAgent` di layar penawaran. Jalan menujunya
memang hidup: `TreatyInTestAgent` dirujuk **dua kali di badan** `Section/InputTreatyInOffer`, dan
satu-satunya langkahnya memanggil `TreatyInSetToDirector`. Tetapi aktivitas itu sendiri:

| Langkah | `pyStepsBlockName` |
|---|---|
| 1, 2, 3, 4 — **seluruh langkah tingkat atas** | **`//`** |
| 4.1 – 4.7 | bersarang **di dalam** langkah 4 yang mati |

**Tidak ada satu langkah hidup pun.** Dan seandainya hidup, langkah 4.5 menuntut
`TreatyIn.Position == "ReasTreatyInGroupLeader"` — keadaan yang tidak dapat dimasuki (NA-10). Dua
lapis mati sekaligus.

**Penjaga ketiga tombol "(dev)" adalah nama orang sebagai teks**, bukan peran, bukan privilese. Dua
nama itu tertulis di dalam aturan; bila nama tampilan orangnya berubah, penjaganya berubah arti
tanpa ada yang menyentuh kode.

### NA-12 · `RevisionState` — penolakan tidak mengosongkannya, dan warisannya membawa dua akibat

Rinciannya di §7.4. Tiga hal yang membalik dugaan:

* **Reject mempertahankan `"1"`** (cabang `2.2.2.4` menulis ulang nilainya). Yang mengosongkan
  hanya Accept dan Decline.
* Addendum yang mewarisi `1` bukan hanya berantai dua tingkat; `TreatyInSetEdit` langkah 2
  membacanya dan menyetel `ViewState = 1` — addendum itu **lahir terkunci**.
* Pasangan `Resolve Complete` + `RevisionState = 1` adalah keadaan yang **tidak ditangani cabang
  mana pun** di `Akseptasi_DT`.

### NA-13 · TDA-11 dikuatkan, dan sumbernya diperbaiki

Rinciannya di §7.6. Grid picker bersumber `TempMasterList.pxResults`, diisi `TreatyLoadMasterJoinEdm`
— `UNION` dua tabel **tanpa `WHERE`**. Konfigurasi kedua pada grid yang sama
(`pyRDName = BrowseTREATY_IN`, `StatusAkseptasi = "Resolve Complete"`) tidak dapat hidup, sebab
kolom yang ditampilkan adalah alias `.CARI1`-`.CARI7` yang tidak dikembalikan RD itu.

**Dua tanda yang bertentangan di satu kendali** — dan yang menang ditentukan oleh kolom yang
tampil, bukan oleh yang tertulis lebih rapi (`METODE` §5.1).

### NA-14 · Audit `TA-04` — lima jenis penulis dikalibrasi, lalu semua klaim negatif disapu ulang

`MA-05` menunjukkan sebuah sapuan penulis dapat mengembalikan nol semata-mata karena tag yang salah.
Maka setiap klaim "nol" dan "satu-satunya" diperiksa ulang dengan perkakas baru
[`tools/tulis.py`](../tools/tulis.py), yang menyapu **lima jenis penulis** dan **mencetak apa yang
ditolaknya**.

**Kalibrasi dijalankan lebih dulu** — setiap penyapu diuji atas kasus positif yang sudah diketahui,
dan hasilnya dilaporkan sebelum sapuan mana pun dipercaya:

| Jenis | Tag sasaran | Kasus positif | Hasil |
|---|---|---|---|
| `PS` Property-Set di Activity | `PropertiesName` / `PropertiesValue` | `SetTreatyIn_Act` 7 menulis `RevisionState` | **4** — LULUS |
| `DT` Data Transform | `pyPropertiesName` / `pyPropertiesValue` (+ `pyDisabled`) | `Akseptasi_DT` menulis `Position` | **48** — LULUS |
| `SEC` kontrol Section | `pyPropertyTarget` | autocomplete `Ceding` -> `CedingID` | **2** — LULUS |
| `BIND` pengikatan kontrol | `pyValue` pada sel ber-`pyReadOnly != true` | radio group `EDMState` di picker | **2** — LULUS |
| `RD` Report Definition | `pyTargetProperty` | `BrowseTREATY_IN` memetakan `.ID` | **29** — LULUS |

**Titik buta yang dilaporkan perkakasnya sendiri, bukan disembunyikan:** **32** langkah `Java` dan
**191** langkah `RDB-List` / `Connect-REST` / `Obj-Browse` tidak terurai. Setiap klaim negatif di
bawah karena itu berlaku untuk **penulisan yang terbaca dari aturan**, bukan untuk seluruh
kemungkinan.

**Hasil sapuan ulang:**

| # | Klaim | Jenis yang disapu | Hasil |
|---|---|---|---|
| 1 | "Tidak ada aturan yang menetapkan `EDMState = 2`" (§4.2, TDA-13) | kelima | **BERDIRI**. Aturan menulis `1` (`TreatyCreateEDM` 2.1) dan `3` (3.1); `TreatyInEDMSetValue` menulis dari `Param.InternalType`. Nilai `2` **hanya** dari radio picker |
| 2 | "`EDMMaterialType` hanya ditulis satu aturan" (§4.3) | kelima | **KELIRU sebagai kalimat**. Penulisnya **tiga**: satu Property-Set + dua pengikatan kontrol (satu dapat disunting, satu hanya-baca). "Satu-satunya **aturan**" benar; "satu-satunya **penulis**" tidak |
| 3 | "Penulis `CedingID` ada dua" (NA-08) | kelima | **KELIRU**. Penulisnya **empat** — tambahan: `TreatyInSetReinsured` 1.2 (DT) dan `Section/ShowSummary` (`pxTextInput` dapat disunting). Tiga dari empat menulis nama **dan** pengenal bersama; kesimpulan pokok NA-08 tetap berdiri |
| 4 | "Kontrol `Ceding` di `InputTreatyInAdjustment` ada dua" (§4.6) | `BIND` | **KELIRU**. Ada **satu** |
| 5 | "Halaman `PoductName` tidak pernah diisi" (§6.2) | kelima **+ Java** | **BERDIRI**. Delapan kemunculan di badan, semuanya teks precondition; nol sasaran tulis; nol di 32 sumber Java |
| 5a | "`Ceding` 2 + `Commencement` 2 + `Termination` 2 dapat disunting di layar addendum" (GRL-05) | `BIND` | **KELIRU angkanya — `1 + 1 + 1`.** Kesimpulannya **tidak berubah**: ketiganya tetap dapat disunting tanpa syarat apa pun. Dan sapuan yang sama memunculkan **NA-17** |
| 6 | "Daftar tertutup penulis `RevisionState`" (§7.4) | kelima | **BOCOR SATU**. `TreatyInCopy` langkah 1 mengosongkannya. Ia aktivitas "salin kontrak", **tidak ada di jalur addendum** — kesimpulan §7.4 tidak berubah |

#### Titik buta ditutup 25 September 2026 — jenis penulis keenam, ketujuh, dan kedelapan

Tiga jenis yang semula hanya dilaporkan sebagai titik buta kini **disapu**, sehingga kadar klaimnya
naik dari "berdiri untuk lima jenis penulis" menjadi berdiri untuk **delapan**.

**6. Langkah `Java` — 32 sumber, disapu sebagai teks.**

| Nama properti dicari | Kemunculan di 32 `pyStepsJavaSource` |
|---|---|
| `EDMState`, `EDMMaterialType`, `RevisionState`, `PoductName`, `CedingID`, `Ceding`, `Commencement`, `Termination`, `Position`, `StatusAkseptasi` | **nol, seluruhnya** |

Satu-satunya penulisan yang keluar dari langkah Java adalah `putString` ke **halaman parameter**:
`FileBase64`, `Ext`, `LinkRefFrom`, `pyIsDataPageRefreshed`. Tidak satu pun properti bisnis.

**Tetapi ada satu penulis generik, dan ia harus dinamai:** delapan langkah berjudul *"Map oracle
column to clipboard"* menjalankan

```java
ClipboardPage tempPage2 = tools.findPage("TreatyIn");
String IsiDataJson = DataJSONPage.getString("CLASSOFBUSINESS");
tempPage2.adoptJSONObject(IsiDataJson);
```

`adoptJSONObject` menulis **setiap properti yang ada di dalam JSON** ke halaman `TreatyIn`. Jadi ia
**dapat** menulis `EDMState`, `RevisionState`, dan seluruh nama di tabel atas — tetapi hanya
**nilai yang sudah tersimpan di baris itu**. Ia **memulihkan**, tidak **mencipta**. Klaim seperti
"tidak ada aturan yang menetapkan `EDMState = 2`" karena itu tidak goyah: `adoptJSONObject` akan
memulangkan `2` bila `2` pernah tersimpan, dan `2` hanya dapat lahir dari radio picker.

Dan ia **tidak dapat** menjelaskan `PoductName`: `adoptJSONObject` bekerja pada halaman `TreatyIn`,
sedangkan `PoductName` adalah **halaman tingkat atas**, bukan properti `TreatyIn`.

**7. Connect-REST — pemetaan responsnya diperiksa.** Aturan Connect-REST di korpus ini **satu per
ekspor**, bukan tiga: `ServiceGoogle`. Responsnya dipetakan `pyMapToKey = UploadDoc.Response`,
`pyMapTo = JSON` — ke halaman `UploadDoc`, bukan ke `TreatyIn`. Nol sentuhan properti bisnis. (Angka
"tiga" pada perintah kerja agaknya dari **enam langkah** `Connect-REST`, yang semuanya memanggil
aturan yang satu itu.)

**8. RDB-List, Connect-REST, dan Obj-Browse — 191 langkah. Rumusan "hanya memuat, tidak mencipta"
TIDAK dapat dinyatakan, dan itu harus dikatakan.** Dari 42 aturan RDB-List per ekspor:

| Jenis SQL | Jumlah (kedua ekspor) |
|---|---|
| `SELECT` | 48 |
| pemanggilan prosedur | 14 |
| `DELETE` | 12 |
| `INSERT` | 6 |
| `UPDATE` | 4 |

**Delapan belas dari empat puluh dua menulis ke basis data**, dan **empat mencipta nilai baru**:
`GenerateImageID_SQL`, `GetCurrentDate`, `InsertToLogAchievement_SQL`, dan `InsertAttachment2_Sql`.
Contohnya yang paling penting justru di jalur utama: `SaveTreatyIn` memanggil `PEGA_TREATY_IN`, yang
di cabang sisip menjalankan **`M_TREATY_IN_SEQ.nextval`** dan mengembalikannya sebagai
`IDPegaOut` — sebuah **pengenal kontrak yang benar-benar baru**.

**Yang dapat dinyatakan, dan lebih tepat:** langkah-langkah ini **tidak pernah menulis langsung ke
`TreatyIn.*`**. Hasilnya mendarat di halaman hasilnya sendiri (`TempResult.pxResults`,
`TempMasterList.pxResults`, `Time.pxResults`), dan berpindah ke properti bisnis **hanya** lewat
Property-Set atau langkah Java — dua jenis yang sudah disapu. Nilai ciptaan di atas pun mengikuti
jalur itu: `IDPegaOut` menjadi `TreatyIn.ID` lewat `SaveTreatyIn_Act` **langkah 10**, sebuah
Property-Set yang penyapu `PS` memang menemukannya.

Ketujuh nama properti yang diaudit **tidak muncul** sebagai sasaran di SQL mana pun; yang muncul
(`CEDINGID`, `COMMENCEMENT`, `TERMINATION`, `EDMSTATE`, `EDMMATERIALTYPE`) seluruhnya **parameter
keluar** menuju prosedur simpan, atau alias kolom `SELECT` pada pemuat picker.

**Penanda mati juga disapu ulang.** Seluruh **73 Data Transform** di kedua ekspor diperiksa untuk
baris ber-`pyDisabled = true`. Hasilnya: **satu berkas saja** — `Akseptasi_DT`, tiga belas baris,
seluruhnya cabang 1.4. Tidak ada klaim hidup/mati lain yang terpengaruh. Mesin selisih seluruhnya
**Activity**, bukan Data Transform, sehingga penanda ini tidak berlaku di sana dan `pre.py` tetap
memadai.

### NA-15 · Tombol paksa di layar addendum menyimpan ke tabel kontrak, dan tidak menulis apa pun

Rinciannya di §6.4. Rantainya: `SaveTreatyIn_Act` langkah 4 -> `POOLDATA.PEGA_TREATY_IN` ->
`ELSE UPDATE … WHERE ID = IDPega` pada `M_TREATY_IN` dan `TREATY_IN`. Pengenal addendum
(`‹kontrak›/Rnn`) tidak ada di kedua tabel itu, sehingga `UPDATE` **tidak mengenai baris**.

Di Oracle itu bukan galat: tidak ada exception, `COMMIT` tetap jalan, `StsSave := 1`, dan pesannya
*"Data Sudah Disimpan Dengan ID : …"*.

* **Baris kontrak asalnya aman** — penjaganya `WHERE ID = IDPega`, dan pengenal addendum tidak
  pernah sama dengan pengenal kontrak.
* **Addendumnya tidak tersimpan** — perubahan status hilang bersama halaman, sementara layar
  melapor berhasil.

Kegagalan senyap (`METODE` §2.2). **Tidak dapat diukur UA mana pun**, justru karena tidak ada yang
tersimpan; itu dikatakan apa adanya, bukan dikarang jadi kueri.

### NA-16 · Ketiga tombol "(dev)" dijaga divisi, bukan terbuka untuk semua

Rantai syarat ditelusuri sampai akar seksi. Di atas ketiga sel itu ada layout `pySectionId = S3`
dengan `pyIsVisibilityOption = CONDITION` dan
`pyContainerVisibleWhen = OperatorID.pyOrgDivision = 'IT'`.

Jadi yang benar: tombol-tombol itu terlihat oleh **operator divisi `IT`**. Di dalam divisi itu,
`Force Resolve Complete(dev)` dan `ReturnToInputor(dev)` masih disaring dua nama; **`Force Edit
(dev)` tidak**, sebab `pyVisible = ALWAYS` membuang kondisi namanya.

`pyAssociatedPrivileges = false` pada **setiap** lapis — tidak ada privilese sama sekali. Pada
`Section/InputTreatyInAdjustment`, tiga dari empat penyertaan `TreatyInActionButtons` berada di
layout `S14` ber-`pyContainerVisibleWhen = OutputParam.DATASHOW=1`.

Ini mengoreksi **dua** pernyataan sekaligus: klaim saya sendiri di putaran sebelumnya, dan
**ADR-0055 §4** yang mencatat `TreatyInForceEdit` *"terlihat oleh setiap pengguna layar
penawaran"*.

### NA-17 · `ShowSummary` membuka kelima field lapisan beku sebagai teks bebas

Rinciannya di §4.6. Yang penting: dua field — **`CedingID`** dan **`ProportionType`** — tidak dapat
disunting di layar mana pun **kecuali** di sini, dan di sini keduanya `pxTextInput`, tanpa daftar
pilihan dan tanpa pemeriksaan.

`ShowSummary` dirujuk hanya oleh `InputTreatyInOffer` beserta harness dan flow action-nya, jadi ini
**layar penawaran** — cacat Treaty In yang berjalan hari ini, bukan cacat jalur addendum. Diperiksa
lebih dulu sesuai `MA-06`: **tidak satu pun dokumen di `_migration-docs/treaty-in/` menyebut
`ShowSummary`.** Ia titik buta induk.

Dicatat sebagai usulan untuk induk lewat **GRL-01 butir 5**.

### NA-18 · Keempat kontrol di layar penawaran dibatasi peran, dan hanya satu menyentuh kontrak yang sudah disetujui

Rantai tampil keempat sel ber-`pyActivity = SetTreatyIn_Act` ditelusuri sampai akar. Labelnya kini
**terbaca** — sesuatu yang daftar eskalasi induk menyatakan tidak dapat dibaca:

| # | Label | Baris | Parameter yang dikirim | `pyVisible` | Syarat sel |
|---|---|---|---|---|---|
| 1 | **Revision** | 269782 | `ID`, `viewstate=1`, **`revisionstate=1`** | `OTHER` | workbasket(2) = `ReasTreatyInAdmin` **&&** `Position = ''` **&&** `StatusAkseptasi = 'Resolve Complete'` |
| 2 | **Copy** | 269396 | `ID` saja | `OTHER` | idem |
| 3 | **View** | 268955 | `ID`, `viewstate=1` | `ALWAYS` | kondisinya **tidak dievaluasi** |
| 4 | **Edit** | 268415 | `ID`, `viewstate=` **kosong**, `revisionstate=` **kosong** | `OTHER` | workbasket(2) = `ReasTreatyInAdmin` && `Position` ∈ {`ReasTreatyInAdmin`, ``} && `StatusAkseptasi` ∉ {`Decline`, `Resolve Complete`} |

> **Ralat 25 September 2026.** Versi pertama tabel ini menulis `Edit` **mengirim
> `revisionstate=1`**. Ia mengirim **nama parameternya dengan nilai kosong**. Saya menghitung
> **keberadaan kata** `revisionstate` di dalam subpohon sel, bukan **nilainya** — bentuk kesalahan
> yang sama dengan `MA-05`. Akibatnya besar dan arahnya melegakan: precondition
> `param.revisionstate==1` pada `SetTreatyIn_Act` langkah 7 **tidak terpenuhi**, sehingga `Edit`
> **tidak** menulis `RevisionState`, **tidak** mengosongkan `StatusAkseptasi`, dan **tidak**
> memendekkan rantai persetujuan. Lihat `MA-08`.

Wadahnya `S168` ber-`pyContainerVisibleWhen = OutputParam.DATASHOW !=1`; `pyAssociatedPrivileges =
false` di setiap tingkat.

**Tiga hal yang membalik pernyataan yang sudah beredar di daftar eskalasi induk:**

1. Keempatnya **dibatasi peran** di tingkat sel — `pyUserData/pyVisible = OTHER` dengan
   `pyCondition` yang menuntut `OperatorID.pyWorkBasketList(2).pyWorkBasketName =
   'ReasTreatyInAdmin'`, yaitu peran pengisi (baris 268820, 269268, 269652, 270086). Bukan "tidak
   dibatasi peran dan selalu terlihat oleh siapa pun". Wadahnya `S168` ber-`pyContainerVisibleWhen
   = OutputParam.DATASHOW != 1`; `pyAssociatedPrivileges = false` di setiap tingkat — **penjaganya
   peran lewat workbasket, bukan privilese**.
2. Yang mengirim `revisionstate=1` **hanya `Revision`**, dan hanya ia mensyaratkan
   `StatusAkseptasi = 'Resolve Complete'`. Jadi yang "mengosongkan persetujuan lalu menyimpan" ada
   **satu**, bukan dua — dan tidak ada jalan memendekkan rantai bagi kontrak yang belum pernah
   disetujui.
3. **Labelnya terbaca** dari `pyLabel` di dalam sel: `Revision`, `Copy`, `View`, `Edit`
   (`Section/InputTreatyInOffer.xml` baris 269782, 269396, 268955, 268415).

Yang tetap berdiri: `Revision` memang mengosongkan `StatusAkseptasi` sebuah kontrak yang sudah
disetujui dan menyimpannya dalam satu tindakan, meninggalkan komentar "Create Revision" (§4.5).

### NA-17a · `CedingID` teks bebas — satu-satunya penyuntingnya, di balik satu nama operator

`Section/ShowSummary.xml` memuat `pxTextInput` **dapat disunting** (`pyReadOnly = false`,
`pyDisabledNew = false`, tanpa `pyDisabledWhen`) untuk `TreatyIn.CedingID`. Itu **satu-satunya**
tempat di seluruh korpus di mana pengenal cedant dapat diketik langsung, terpisah dari namanya dan
tanpa daftar pilihan. Bentuk paling murni dari ketidaksinkronan yang diukur `UA-10`.

**Tetapi kadarnya jauh lebih rendah daripada yang saya laporkan, dan itu harus disebut lebih dulu.**
Tombol yang membuka `ShowSummary` — label *"dsp sum"*, `pyLocalAction = ShowSummary` di
`Section/InputTreatyInOffer` — berbunyi:

```
pyVisible   = OTHER
pyCondition = OperatorID.pyUserName= 'ALDO SAPUTRA1'
```

Satu nama operator, dan namanya **berakhiran `1`** — berbeda satu karakter dari `'ALDO SAPUTRA'`
yang dipakai ketiga tombol "(dev)" di layar yang sama.

> **KOREKSI 26 September 2026 — kadarnya saya nyatakan terlalu tinggi.** Saya menyebut pola ini
> "sekeluarga dengan `1=2`". **Ia tidak sekeluarga.** `1=2` salah **secara logika**: ia tidak pernah
> benar, apa pun datanya, dan itu terbaca dari ekspor. `pyUserName = 'ALDO SAPUTRA1'` hanya salah
> bila **tidak ada operator bernama persis itu** — dan itu **pertanyaan data, bukan ekspor**
> (`METODE` §3.1, §4.6). Akun uji dengan nama berakhiran angka lazim ada.
>
> **Kadar yang benar:** *tidak terbaca apakah operator itu ada; bila ada, `ShowSummary` terjangkau
> olehnya, dan `CedingID` dapat diketik bebas.* Diperiksa lewat **`UA-15`**.

Maka `ShowSummary` **tampaknya** layar pengembang — dan "tampaknya" itu menunggu `UA-15`. Di dalamnya bahkan ada satu sel
`Ceding` ber-`pyCondition = 1=2` — tidak pernah tampil. Yang tetap berlaku: **tidak ada satu pun
penyunting `CedingID` yang tervalidasi di mana pun**, dan penjaga satu-satunya berupa **nama orang
di dalam aturan** (keluarga TDA-09), bukan peran atau privilese.

### NA-17b · Lapisan beku dapat disunting pada kontrak yang dibuka kembali

Inilah bentuk NA-17 yang benar-benar menyentuh lapisan beku. Menyunting kelima field itu **saat
kontrak sedang diisi** adalah kerja biasa — lapisan beku ADR-0040 berlaku **antar versi**, bukan
saat kontrak dilahirkan (`METODE` §1.2). Yang menjadi cacat adalah menyuntingnya pada kontrak yang
**sudah disetujui lalu dibuka kembali**.

Kedua jalur pembukaan kembali diperiksa sampai akar seksi:

| Field | `InputTreatyInOffer` | `InputTreatyInAdjustment` | Syarat keadaan |
|---|---|---|---|
| `Ceding` | `pyReadOnly = false` | `pyReadOnly = false` | **tidak ada** |
| `Commencement` | `pyReadOnly = false` | `pyReadOnly = false` | **tidak ada** |
| `Termination` | `pyReadOnly = false` | `pyReadOnly = false` | **tidak ada** |
| `ProportionType` | `pyReadOnly = true` | `pyReadOnly = true` | — |
| `CedingID` | tidak ada kontrol | tidak ada kontrol | — |

Satu-satunya wadah di atasnya adalah `pyContainerVisibleWhen = OutputParam.DATASHOW != 1`, dan
`DATASHOW` ditulis `TreatyInInputVis` dari `param.Add` — ia membedakan **tambah-baru** dari
**muat-yang-ada**, bukan boleh-sunting dari tidak. `pyAssociatedPrivileges = false` di setiap
tingkat.

**Maka sesudah `Revision` atau `Force Edit (dev)` membuka kembali sebuah kontrak, ketiga field itu
dapat disunting — tanpa satu pun syarat yang membedakan kontrak yang sudah disetujui dari draf.**
Di model baru pembukaan seperti itu melahirkan versi baru, sehingga lapisan beku ADR-0040 berlaku
dan perubahan semacam ini ditolak. Label **PERUBAHAN** pada GRL-05 karena itu berdiri dengan dasar
yang lebih tepat daripada sebelumnya.

### NA-19 · `Save EDM(dev)` menyimpan addendum tanpa syarat status

Dua tombol simpan untuk addendum, di seksi yang sama:

| Tombol | `pyCondition` sel | Wadah | Syarat status |
|---|---|---|---|
| **Save** | `TreatyIn.ViewState != '1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'` | `S2: TreatyMasterInEDM` | **ada** — addendum yang sudah disetujui tidak dapat disimpan |
| **Save EDM(dev)** | `OperatorID.pyOrgDivision = 'IT' && TreatyMasterInEDM` | `S3: OperatorID.pyOrgDivision = 'IT'` | **tidak ada** |

Tombol pertama adalah peredamnya, dan ia bekerja: addendum ber-`Resolve Complete` tidak punya
tombol simpan biasa. Tombol kedua **melewatinya** — tidak menyebut `StatusAkseptasi` maupun
`ViewState` sama sekali.

Digabung dengan `Force Edit (dev)` (`pyVisible = ALWAYS`, terlihat setiap operator divisi IT, dan
menyetel `ViewState = 0`), hasilnya: **seorang operator divisi IT dapat membuka kunci addendum yang
sudah disetujui, menyuntingnya, lalu menyimpannya di tempat** — tanpa versi baru, tanpa persetujuan
ulang.

**Akibatnya untuk GRL-02.** Bagian **PELESTARIAN** di GRL-02 berbunyi *"niat membekukan angka dasar
persetujuan sudah dipenuhi sistem lama"*. Itu tetap benar **untuk mesinnya** — `OLDDATA` dan
`ValueDifference` memang beku di dalam `JSONDATA`. Tetapi ia **tidak terlindung dari penyuntingan
manual** lewat pasangan tombol di atas. Sumber PERUBAHAN ketiga karena itu ditambahkan ke GRL-02
dan ke `REV-1`: *baris addendum yang sudah disetujui dapat disunting dan ditimpa di tempat oleh
operator divisi IT*. Besarannya tidak dapat diukur langsung — penyuntingan di tempat tidak
meninggalkan jejak — tetapi `UA-10` menangkap sebagiannya lewat pelanggaran lapisan beku.

### NA-20 · Penggeser `Termination` otomatis ada, tetapi pemicunya tidak dapat ditekan

Pernyataan saya sebelumnya — *"`TreatyInSetTreatyYear` dirujuk `TreatyInNONProportional`, bukan
layar addendum"* — **bertentangan dengan dirinya sendiri** dan dicabut: `TreatyInNONProportional`
**adalah** kolom NILAI BARU pada layar addendum, disertakan `Section/InputTreatyInAdjustment`
sebanyak **tujuh kali di badan**, dan saya sendiri memakainya sebagai bukti di GRL-05.

Diperiksa sampai akar, dan hasilnya berbeda dari yang saya duga maupun yang dikhawatirkan:

| Sel | Berkas | Memicu? | Keadaan |
|---|---|---|---|
| `TreatyIn.Commencement`, `pxDateTime` | `TreatyInNONProportional` | **ya** — `pyEvent = change` -> `refresh`, `pyDataTransform = TreatyInSetTreatyYear` | **`pyReadOnly = true`, `pyDisabledNew = true`** |
| `TreatyIn.OLDDATA.Commencement` | `TreatyInNONProportionalOldData` | ya | `pyReadOnly = true` |
| `.Commencement`, `pxDateTime` **dapat disunting** | `InputTreatyInAdjustment` | **tidak** — `pyEvent` kosong, tanpa aksi | dapat disunting |
| `.Commencement`, `pxDateTime` **dapat disunting** | `InputTreatyInOffer` | **tidak** | dapat disunting |

Dan `TreatyInSetTreatyYear` **tidak dipanggil aktivitas mana pun** — hanya dua sel di atas.

**Maka penggeseran otomatis `Termination` tidak dapat dipicu pengguna**: sel yang memicunya
hanya-baca, dan sel yang dapat disunting tidak memicunya. Bukan karena seksinya tidak ada di layar
addendum — ia ada. Ini bentuk `METODE` §2.0 yang lain: **kemampuan yang terpasang pada kontrol yang
tidak dapat disentuh**.

Satu sisa yang tidak ditutup: sel pemicu itu juga membawa `pyDisabledWhen = TreatyIn.EDMMaterialType = 2`
di samping `pyReadOnly = true`. Syarat itu **tidak dievaluasi** selama selnya hanya-baca; ia sisa,
dan dicatat supaya tidak dibaca sebagai bukti oleh pembaca berikutnya.

### NA-21 · `Revision` juga mengunci — selektif, dan kuncinya bocor

Rinciannya di §7.7. Tiga angka yang mengubah pembacaan:

* `Revision` mengirim `viewstate=1` **dan** `revisionstate=1`; langkah 6 `SetTreatyIn_Act` menulis
  `ViewState = 1`, `IsEditData = 1`.
* Dari sel FIELD yang dapat disunting di seluruh `Section/`: **96 dijaga `ViewState`, 882 tidak.**
  Yang dijaga terpusat di `DetailLimits` (20), `Layers` (20), `TreatyInTabsNonProportional` (14),
  `TreatyInTabsProportional` (13) — batas per bahaya, deductible, mata uang, bagian, brokerase.
* **Ketiga field lapisan beku tidak dijaga sama sekali.**

**Akibatnya bukan memberatkan, melainkan mengubah bentuk pertanyaannya.** Kunci yang menyasar angka
uang membuat rantai dua tingkat terbaca sebagai **kebijakan yang mungkin dikehendaki** (`METODE`
§1.2): perubahan tanpa uang cukup ditinjau lebih pendek. Yang cacat karena itu bukan rantai
pendeknya, melainkan **kuncinya yang bocor** — cedant dan periode, yang justru menentukan kontrak
itu kontrak yang mana, tetap dapat diubah.

`DB-10` ditulis ulang mengikuti bentuk ini.

### NA-22 · Kunci `ViewState` tidak konsisten — dicocokkan ke daftar 33 akar

§5.6 mendaftar 33 akar `ValueDifference`; 21 di antaranya agregat `Total*`, sehingga **dua belas**
adalah besaran yang benar-benar dapat disunting. Dicocokkan ke sel yang dijaga `ViewState`:

| Seksi | Dijaga | Tanpa penjaga | Akar yang disentuh |
|---|---|---|---|
| `Installments` | **0** | 2 (`.AmountTotal`, `.PctTotal`) | #6 `Installment` |
| `DetailShareRetro` | **0** | 2 | #5 |
| `ShareRetro` | 1 | **27** | #12 |
| `Share` | 2 | **33** | #11, #12 |
| `LayersEDM` | 3 | **18** (`.Deductible`, `.AgregateLimit`, `.AdjRate`) | #10 |
| `DetailLimits` | 20 | 5 (**`.Value`**, **`.Currency`**) | #9, #10 |
| `Layers` | 20 | 3 | #10 |
| `TreatyInTabsNonProportional` | 14 | **92** | #1, #10 |

**Tidak satu pun dari dua belas akar terkunci seluruhnya**, dan pasangan yang paling telak:
`Layers.xml` menjaga `.Deductible`, `.AgregateLimit`, `.AdjRate`; **`LayersEDM.xml`, kembarannya,
tidak menjaga satu pun.**

**Akibatnya, dan ia membatalkan bacaan saya sendiri satu giliran sebelumnya.** Dugaan bahwa
`Revision` "mengunci angka uang" sehingga rantai dua tingkat terbaca sebagai kebijakan — **gugur**.
Sesudah `Revision`, deductible, batas agregat, nilai dan persentase bagian, serta total angsuran
**masih dapat diubah**, dan perubahan itu menempuh rantai dua tingkat.

`ViewState = 1` karena itu bukan penanda "revisi tanpa uang"; ia penjaga yang dipasang **sebagian,
tanpa pola yang dapat dipertanggungjawabkan**.
