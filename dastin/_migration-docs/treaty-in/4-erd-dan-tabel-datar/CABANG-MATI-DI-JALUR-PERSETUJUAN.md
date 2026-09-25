# Cabang dan langkah mati di jalur persetujuan — didaftar, bukan diadili

**Tanggal:** 24 September 2026
**Sebab:** *"membaca cabang yang hidup lalu menyimpulkan tentang seluruh cabang"* adalah bentuk yang
sama dengan kekeliruan kemarin, dan ia mengintai di mesin keadaan persetujuan.

---

## 0. KOREKSI LEBIH DULU — **tidak ada dua puluh lima cabang yang belum dibuka**

Tabel saya di `SAPUAN-DAN-NAMA-TAGNYA.md` §4 berbunyi:

> `pyDisabled == "true"` · **`Rule-Obj-Model` 26** · seluruhnya di `Akseptasi_DT.xml`

Kalimat itu **benar tetapi menyesatkan**, dan saya yang menyesatkan dengan menuliskannya begitu.
Yang sebenarnya:

| | |
|---|---:|
| `Akseptasi_DT.xml` di **`Treaty In`** | **13** baris mati |
| `Akseptasi_DT.xml` di **`Treaty In Adjustment`** | **13** baris mati |
| **berkas berbeda** | **satu** — dua salinan dari aturan yang sama |

Dan kedua salinan itu **identik secara logika**: 92 baris masing-masing, **nol baris yang hanya ada
di salah satunya**, dan ketiga-belas baris matinya **sama persis** — seluruhnya cabang **1.4
`GroupLeader`**, yang **sudah dibuka** untuk 7.10. *(Byte-nya berbeda — panjangnya sama, md5-nya
tidak — yang wajar untuk dua ekspor aturan yang sama.)*

Ketiga belasnya juga bukan tiga belas *cabang*: ia **empat cabang** (`1.4`, `1.4.1`, `1.4.2`,
`1.4.3`) dan **sembilan penugasan** di dalamnya.

> **Jadi: nol cabang mesin keadaan yang belum dilihat.** Saya melaporkan jumlah **baris** sebagai
> jumlah **tempat**, dan tidak menyebut bahwa dua modul membawa salinan aturan yang sama. Keduanya
> kesalahan yang sama bentuknya: **angka tanpa satuan, dan semesta tanpa batas.**

**Tetapi pertanyaannya tetap sah, hanya tempatnya salah.** Mesin keadaannya bersih; **aktivitas yang
menjalankannya tidak.** Di situlah kode mati jalur persetujuan berada, dan di situ ia **belum pernah
didaftar**.

---

## 1. Yang benar-benar ada: **20 langkah mati di keluarga persetujuan**

Keluarga didefinisikan mekanis: aktivitas yang **namanya** menyebut persetujuan/pengajuan/penolakan,
**atau** yang isinya menyentuh `StatusAkseptasi` atau `.Position`.

**Dan ke-20 itu identik di kedua modul** — berkas yang sama, langkah yang sama, jumlah yang sama.
Tidak ada satu pun yang mati di satu modul dan hidup di modul lain.

| Aktivitas | mati / langkah puncak | Akibatnya |
|---|---:|---|
| **`TreatyInSetToDirector`** | **4 / 4** | **seluruh aktivitas mati** — §2.1 |
| **`TreatyInDeclineConfirmation_postactEDM`** | **4 / 8** | yang mati **justru pencatatan penolakannya** — §2.2 |
| **`TreatyInSubmitEDM`** | **3 / 9** | **seluruh pemeriksaan sebelum kirim** — §2.3 |
| `SaveTreatyIn_EDM_Act` | 3 / 17 | penulisan ke tabel dasar — §2.4 |
| `TreatyInConvertCallData_act` | 3 / 15 | sisa uji coba konversi — tidak berakibat |
| `TreatyInSubmit` | 1 / 8 | hanya `TreatyInCheckID`; `TreatyInCheckError` **tetap hidup** — §2.3 |
| `SaveTreatyIn_Act` | 1 / 11 | penulisan `TREATYINOFFER` — §2.5 |
| `AddCommentList_Act` | 1 / 3 | satu `RDB-List` tanpa keterangan; kedua `Property-Set` **hidup** |

---

## 2. Isi tiap langkah mati yang berakibat — syarat, yang disetel, dan ke mana menyalurkan

### 2.1 `TreatyInSetToDirector` — **seluruh aktivitas mati**

```
1  RDB-List      Getting Oracle time                                    MATI
2  Property-Set  Local.CurrTime = Time.pxResults(1).CARI1               MATI
3  RDB-List                                                             MATI
4  ?                                                                    MATI
```

Namanya menyatakan maksudnya: **menyetel kontrak langsung ke Direktur.** Ia **dirujuk** oleh aturan
lain — ia salah satu dari empat tempat `ReasTreatyInGroupLeader` diuji sebagai kondisi — tetapi
seluruh isinya dimatikan.

> **Ini jalur pintas persetujuan yang pernah ada lalu dimatikan seluruhnya.** Ia **tidak**
> membuktikan adanya tingkat persetujuan tambahan; ia membuktikan adanya **jalan pintas** yang
> melompati tingkat. Golongannya berbeda dari bacaan A, dan lebih dekat ke **P-40 — tidak ada jalan
> pintas** dan **eskalasi butir 1**.

### 2.2 `TreatyInDeclineConfirmation_postactEDM` — **yang mati adalah pencatatannya; yang hidup adalah penghapusannya**

```
1  Property-Set   Param.Info / Param.Comment                            MATI
2  Call AddCommentList_Act                                              MATI
3  Property-Set   CommentList(<LAST>).IsApproved = "Decline"            MATI
                  StatusAkseptasi               = "Decline"             MATI
4  Call SaveTreatyIn_Act                                                MATI
------------------------------------------------------------------ hidup:
5  Property-Set   InputData.CARI1 = TreatyIn.ID
6  RDB-List       RDB remove EDM
7  RDB-List       RDB remove EDM
8  call TreatyInInputVis
```

Bandingkan dengan saudaranya untuk kontrak dasar, **`TreatyInDeclineConfirmation_postact`**:
**keempat langkahnya hidup**, dan ia memang mencatat `Decline` lalu menyimpan.

> **Untuk ADDENDUM, menolak berarti MENGHAPUS.** Penolakan addendum tidak meninggalkan baris
> `CommentList`, tidak menyetel `StatusAkseptasi`, dan tidak menyimpan apa pun — ia menjalankan dua
> `RDB remove EDM`.
>
> **Ini menguatkan butir eskalasi yang sudah ada** (*"addendum yang ditolak dihapus, bukan
> ditandai"*) dengan hal yang belum kita punya: **kode penandaannya ADA dan DIMATIKAN.** Bukan
> kelalaian, melainkan **penggantian** — dan itu kalimat eskalasi yang berbeda.
>
> **Akibat pada model, dan ia menyentuh dua kotak yang sedang digambar:** `CATATAN_PERSETUJUAN` dan
> `PERISTIWA_KONTRAK` **tidak akan pernah memuat penolakan addendum** dari data lama, karena
> barisnya tidak pernah ditulis dan addendumnya sendiri tidak ada lagi. Itu **bukan** lubang model —
> ia **batas data** yang harus tertulis, supaya tidak ada yang menyimpulkan *"tidak ada penolakan
> addendum"* dari nol baris.

### 2.3 Pemeriksaan sebelum kirim — **hidup untuk kontrak dasar, mati untuk addendum**

| | `TreatyInSubmit` (kontrak) | `TreatyInSubmitEDM` (addendum) |
|---|---|---|
| `Call TreatyInCheckID` | **MATI** | **MATI** |
| `call TreatyInCheckError` | **hidup** | **MATI** |
| langkah keluar bergalat | **hidup** | **MATI** |

Isi langkah keluar yang mati di sisi addendum:

```
when OutputParam.ERRMSG == ""
     OutputParam.ERRMSG = "ID error: Found one or more empty field/ ID value"
```

> **Addendum dapat dikirim tanpa satu pun pemeriksaan isi; kontrak dasar tidak.** Keduanya lalu
> memanggil `Akseptasi_DT`, `AddCommentList_Act`, dan penyimpanan yang sama.
>
> **Belum diadili di sini**, sesuai perintah: didaftar. Tetapi ia **calon kuat** butir eskalasi, dan
> ia menyentuh embargo modul Adjustment — jadi ia dicatat, tidak dinaikkan sendiri.

### 2.4 `SaveTreatyIn_EDM_Act` — penulisan ke tabel **dasar** dimatikan, penulisan ke tabel **addendum** hidup

```
10  Property-Set  OutputParam.ERRMSG = "Terjadi kesalahan saat menyimpan data - " + …   MATI
                  when OutputParam.STSSAVE == 0
12  RDB-List      "Activity when Resolve Complete insert into m treaty in"              MATI
                  when TreatyIn.ID != "" || TreatyIn.ID != "UnknownId"
13  call SaveTreatyInDetail_Act    "when Resolve Complete, save treatyindetail"         MATI
------------------------------------------------------------------------------ hidup:
14  call SaveTreatyInDetailEdm_Act
```

> **Ini bukti langsung untuk butir eskalasi yang selama ini bersandar pada pengamatan data:**
> *"perubahan kontrak yang disetujui tidak ikut ke tabel rinciannya"*. Langkah yang akan
> menuliskannya **ada, dan dimatikan** — dan yang tetap hidup hanya penulisan ke tabel **addendum**.
>
> Langkah 10 juga layak dicatat tersendiri: **pesan galat penyimpanan dimatikan.** Bila penyimpanan
> gagal (`STSSAVE == 0`), tidak ada yang memberi tahu siapa pun. Ini `penegakan-yang-bisa-berhenti-diam-diam`
> dalam bentuk paling harfiah.
>
> Dan perhatikan syarat langkah 12: `ID != "" || ID != "UnknownId"` — **`||` yang selalu benar**.
> Langkahnya mati, jadi tidak berakibat; tetapi bila ada yang menghidupkannya kembali, syaratnya
> tidak menyaring apa pun. Dicatat supaya tidak dihidupkan tanpa dibaca.

### 2.5 `SaveTreatyIn_Act` — `TREATYINOFFER` berhenti diisi, dan inilah langkahnya

```
9  call Data-Portal.SaveTreatyInOffer_Act                               MATI
   when TreatyIn.StatusAkseptasi == "Resolve Complete"
```

> **Uji AB menanyakan KAPAN `TREATYINOFFER` berhenti diisi.** Pertanyaan *"apa yang menghentikannya"*
> sekarang **terjawab dari ekspor**: langkah yang mengisinya dimatikan. Uji AB **tetap dikirim** —
> ia menjawab *kapan*, yang tidak terbaca dari ekspor, dan **jangan tanyakan apakah, tanyakan
> kapan** berlaku persis di sini.

---

## 3. Yang TIDAK ditemukan, dan itu jawaban juga

| Yang dicari | Hasil |
|---|---:|
| cabang mati di `Akseptasi_DT` selain 1.4 | **nol** |
| langkah mati yang **menyetel `Position`** ke nilai di luar keempat jabatan yang dikenal | **nol** |
| langkah mati yang **menambah tingkat persetujuan** | **nol** |
| beda antara modul `Treaty In` dan `Treaty In Adjustment` pada keluarga ini | **nol** |

> **Tidak satu pun dari 20 langkah mati memperlihatkan tingkat persetujuan yang lain.** Yang mereka
> perlihatkan adalah **jalan pintas** (§2.1), **penggantian pencatatan dengan penghapusan** (§2.2),
> dan **pemeriksaan yang dimatikan** (§2.3, §2.4).
>
> **Batas yang tetap disebut:** ini hanya meliputi penanda mati yang **ada di ekspor**. Cabang yang
> dihapus — bukan dimatikan — tidak meninggalkan jejak apa pun, dan **L-10** tetap berdiri untuk
> jenis aturan yang tidak terekspor.

---

## 4. "TERJANGKAU DARI LAYAR" — diperiksa untuk ketiga aturan yang menyangga eskalasi, dan **ketiganya bertahan**

`pyDisabled` di `Section` dan `Harness` **persis** yang menentukan apakah sesuatu terjangkau dari
layar, dan kalimat *"hidup serta terjangkau dari layar"* dipakai untuk menaikkan tiga aturan ke
butir eskalasi tata kelola. Karena itu ketiganya diperiksa satu per satu.

### 4.1 Cara memeriksanya, dan jebakan yang harus dilewati lebih dulu

Panggilan dari layar duduk di `pyActivity`, di dalam sebuah **perilaku kendali**
(`Embed-Control-Mode-Behaviors`) milik sebuah **ragam kendali** (`Embed-Control-Mode`). Ragam itu
membawa **dua** ruas yang harus dibaca bersama:

```xml
<pyDisabledWhen>TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2</pyDisabledWhen>
<pyDisabled>true</pyDisabled>
```

> **`pyDisabled = true` SENDIRIAN bukan "selalu mati".** Ia berarti *"pematian dipasang"*; yang
> menentukan **kapan** ada di `pyDisabledWhen` di sebelahnya. Membaca yang pertama saja menghasilkan
> kesimpulan terbalik — bentuk yang sama dengan jebakan `pyVisible = ALWAYS`.
>
> Yang benar-benar **mati tanpa syarat** adalah `pyDisabled = true` dengan `pyDisabledWhen`
> **kosong**. Ketiga besaran itu dihitung terpisah di bawah.

### 4.2 Hasilnya, dihitung per kemunculan kendali di kedua modul

| Aturan | kendali **hidup** | mati **bersyarat** | **mati tanpa syarat** |
|---|---:|---:|---:|
| `AddSpreadingXOL` | **8** | 0 | **0** |
| `TreatyInNonAddItem` | **132** | 42 | **0** |
| `TreatyInSetBrokerage` | **35** | 28 | **7** |

> **Ketiganya tetap terjangkau dari layar. Tidak ada kalimat yang dicabut**, dan butir eskalasi tata
> kelola tetap menyebut **tiga** aturan logika transaksi, bukan nol.

Ketujuh kendali `TreatyInSetBrokerage` yang mati tanpa syarat seluruhnya ada di **layar cermin**:
`TreatyInActualShare` (`ActualValue.BrokeragePercent`), `TreatyInTabsNPValueDifferenceProp`
(`ValueDifference.Brokerage…`), `TreatyInTabsNPValueDifference_` (`ValueBeforeProrate.Broker…`), dan
`TreatyInTabsNonProportionalOldData` (`OLDDATA.BrokeragePercent`). **Itu justru benar**: layar potret
dan layar nilai lama memang tidak boleh disunting, dan kendali yang sama di layar penyuntingan
sebenarnya tetap hidup.

### 4.3 Dan pemeriksaan ini membuka satu hal yang BUKAN tentang keterjangkauan

Syarat pematian yang muncul berulang di ketiga aturan hanya dua:

```
TreatyIn.ViewState = 1            -> kontrak sedang dilihat, bukan disunting
TreatyIn.EDMMaterialType = 2      -> ???
```

**`EDMMaterialType` berstatus DITUNDA** di `SPEC-MODEL-DATA.md` §6, dan alasan penundaannya tertulis
sebagai *"menunggu konfirmasi bahwa tidak ada pemakaian lain"*.

> **Pemakaian lain itu ADA, dan baru terbaca sekarang: `EDMMaterialType = 2` MEMATIKAN ruas-ruas
> masukan di layar** — brokerase, bagian seragam, prorata, dan penambahan baris non-proporsional.
>
> Ini menyentuh **ADR-0049**, yang menetapkan materialitas **diturunkan** dari ada-tidaknya baris
> selisih. Bila materialitas diturunkan, sementara di sistem lama ia **menentukan apa yang boleh
> disunting**, maka nilai turunan menjadi penjaga penyuntingan yang menghasilkannya. **Itu tabrakan,
> dan ia dilaporkan, bukan diselesaikan di sini.**
>
> **Yang berubah sekarang juga:** butir §6 `EDMMaterialType` tidak lagi menunggu *"konfirmasi bahwa
> tidak ada pemakaian lain"* — jawabannya **ada pemakaian lain**. Yang ditunggu berganti menjadi
> keputusan atas tabrakan dengan ADR-0049.

### 4.4 Dan satu hal yang harus dikerjakan perkakasnya

Ketika `langkah-hidup.py` diperbaiki, ia melaporkan **tiga besaran terpisah**, tidak satu pun
dijumlahkan:

| Besaran | Penandanya | Artinya |
|---|---|---|
| blok mati di aktivitas | `pyStepsBlockName == "//"` | kode yang tidak dijalankan |
| cabang mati di *data transform* | `pyDisabled == "true"` pada `Rule-Obj-Model` | cabang logika yang dimatikan |
| kendali mati di layar | `pyDisabled` **+** `pyDisabledWhen` pada `Rule-HTML-*` | **dua** besaran lagi: mati tanpa syarat, dan mati bersyarat beserta syaratnya |

Yang terakhir **tidak boleh dilaporkan sebagai satu angka**: *"1.082 kendali mati"* akan salah,
karena sebagian besar mati **hanya pada keadaan tertentu**, dan keadaannya itulah isinya.
