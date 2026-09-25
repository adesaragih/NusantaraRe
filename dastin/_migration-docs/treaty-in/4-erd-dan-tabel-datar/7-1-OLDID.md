# 7.1 — `OLDID` punya berapa arti?

**Tanggal:** 24 September 2026
**Jawaban singkat:** **TIGA.** Itu cabang yang menahan menurut aturan F.1, jadi berkas ini berhenti
pada laporan dan **tidak** memilih bentuk modelnya sendiri.

---

## 1. Cara menyapunya

Seluruh penulis `OLDID` dicari atas **kelima** bentuk penulis — keempat bentuk `BENTUK-PENULIS-PROPERTI.md`
**ditambah bentuk kelima yang baru ditemukan hari ini**, `<CopyFrom>`/`<CopyInto>` pada
`pyStepsCallParams` (salin halaman). Lihat §5.

| | |
|---|---:|
| penulis `TreatyIn.OLDID` | **4** |
| pembaca `TreatyIn.OLDID` | **3** |
| di antaranya penulis yang **mati** | **0** — keempatnya hidup |

---

## 2. Keempat penulis, dan apa yang ditulis masing-masing

### Arti 1 — **versi pendahulu** (addendum / revisi)

`Treaty In Adjustment/Activity/TreatyInEDMSetValue.xml`:

```
TreatyIn.OLDID = TreatyIn.ID          <- ID kontrak/versi yang sedang dibuka
TreatyIn.EDMState = Param.InternalType
...
Call TreatyInRevisi_post              <- membangkitkan ID baru
```

`Treaty In Adjustment/Activity/TreatyInRevisi_post.xml`:

```
langkah 4  (revisi pertama)   : ID = ID + "/R01"
langkah 5  (sudah ada revisi) : OLDID = ID                     <- ID revisi SEBELUMNYA
                                ID    = <nomor pokok> + "/R" + (urut+1)
```

> **Rantainya berbutir satu tingkat:** `OLDID` menunjuk **pendahulu langsung** — kontrak dasar untuk
> `R01`, dan revisi sebelumnya untuk `R02` ke atas. Format ID-nya sendiri membawa nomor urut
> (`1001378/R02`), jadi urutan **tidak** bergantung pada `OLDID`.

### Arti 2 — **kontrak asal salinan**

`Treaty In/Activity/TreatyInCopy.xml`:

```
CommentList(<LAST>).Suggest = "Copied from ID " + TreatyIn.ID
TreatyIn.OLDID              = TreatyIn.ID        <- kontrak yang DISALIN
TreatyIn.ID                 = "UnknownId"        <- baris baru, belum bernomor
TreatyIn.Position           = "ReasTreatyInAdmin"
```

> Di sini `OLDID` menunjuk **kontrak lain yang berdiri sendiri**, bukan pendahulu dalam satu rantai
> versi. Kontrak hasil salinan adalah kontrak **baru** dengan riwayat persetujuannya sendiri —
> `CommentList` dikosongkan lebih dulu.

### Arti 3 — **nomor penawaran warisan** *(yang membuat jawabannya tiga, bukan dua)*

`Treaty In/Activity/TreatyInMappingDataconvert.xml`, dibuat **23 Juli 2019** oleh `ALDO SAPUTRA`:

```
TreatyIn.Information = TreatyIn.Information + " No Offer: " + TempResult.pxResults(1).Ceding
TreatyIn.OLDID       = TempResult.pxResults(1).Ceding
```

**Dua baris berurutan, nilai yang sama.** Baris pertama menamai isinya sendiri: **nomor penawaran**.
*(Nama ruas `Ceding` di situ tidak berarti apa-apa — halamannya diisi langkah `Java` berjudul
"Map oracle column to clipboard", jadi namanya nama kolom menurut kedudukan, bukan menurut isi.
Ini jebakan "label tidak dipercaya" dalam bentuk nama ruas.)*

Pemanggilnya `TreatyInConvertCallData_act.xml` — **konversi sekali jalan dari sistem lama**: ia
menyemai `TreatyIn.ID = "30001"`, mengulang baris warisan, memanggil pemeta ini, lalu menyetel
`StatusAkseptasi = "Resolve Complete"` dan menyimpan.

**Dan pembacanya mengukuhkannya**, `Treaty In/Activity/SaveTreatyInOffer_Act.xml`:

```
SaveData1.IDPEGA  = TreatyIn.ID
SaveData1.NOOFFER = TreatyIn.OLDID        <- TANPA SYARAT
```

> Kolom tujuannya bernama **`NOOFFER`**. Pembaca itu memperlakukan `OLDID` sebagai **nomor
> penawaran** untuk **setiap** kontrak — termasuk yang `OLDID`-nya sebenarnya berarti versi
> pendahulu atau kontrak asal salinan. **Itu bukan dugaan tentang niat; itu yang tertulis.**

### Ketiga pembacanya, untuk kelengkapan

| Pembaca | Membacanya sebagai |
|---|---|
| `FetchTreatyExistingProduction.xml` — `@if(EDMState="", ID, OLDID)` | **arti 1** — untuk addendum, cari kontrak produksinya lewat pendahulunya |
| `SaveTreatyInOffer_Act.xml` — `NOOFFER = OLDID` | **arti 3** |
| `SaveTreatyIn_EDM_Act.xml` — `InputParam.OLDIDPEGA = OLDID` | diteruskan ke prosedur simpan; artinya mengikuti pemanggilnya |

---

## 3. KENAPA INI BERHENTI, dan apa yang saya TIDAK putuskan

Aturan putusan yang ditetapkan: satu arti → seam tetap; dua arti → dua penunjuk; **lebih dari dua →
berhenti dan laporkan.** Jawabannya tiga, jadi berkas ini berhenti.

**Dua yang pertama sudah punya rumah di model, dan keduanya tidak perlu dibahas ulang:**

| Arti | Sudah ada sebagai |
|---|---|
| 1 — versi pendahulu | `ID_VERSI_KONTRAK_DASAR` (seam Adjustment) |
| 2 — kontrak asal salinan | **`ID_KONTRAK_DISALIN_DARI`** — sudah ada di `SPEC-MODEL-DATA.md` §10.1, bersumber *"`OLDID` pecahan 1"* |

*(Nama yang diusulkan `ID_KONTRAK_ASAL_SALINAN` **tidak saya pakai**: atribut dengan arti itu sudah
ada dan sudah dirujuk berkas lain. Mengganti namanya sekarang adalah perubahan nama, bukan
keputusan model.)*

**Yang belum punya rumah adalah arti ketiga**, dan di situ letak pertanyaannya — karena §10.1 sudah
memuat satu atribut yang mendekatinya:

```
NOMOR_KONTRAK_WARISAN   asal: `ID`   boleh kosong: ya
                        "hanya terisi pada baris hasil migrasi (ADR-0042)"
```

> **Dua bacaan, dan keduanya mengubah hal yang berbeda:**
>
> **(i) `NOMOR_KONTRAK_WARISAN` memang benda yang sama, dan ASALNYA yang salah tercatat.** Bila
> begitu, asalnya bukan `ID` melainkan `OLDID`, dan §10.1 diperbaiki — satu baris.
>
> **(ii) Keduanya benda yang berbeda**: satu nomor kontrak di sistem lama, satu **nomor penawaran**
> yang mendahului kontraknya. Bila begitu model kekurangan satu atribut, dan namanya menyebut
> penawaran.
>
> **Ekspor tidak dapat memilih di antaranya** — ia memperlihatkan `OLDID` diisi dari kolom warisan
> dan diberi label *"No Offer"*, tetapi tidak memperlihatkan apakah nomor itu sekaligus nomor
> kontrak lamanya. **Yang memisahkan adalah data**, dan ujinya murah: bentuk nilai `OLDID` pada
> baris hasil migrasi versus bentuk `ID` di `TREATYINOFFER.NOOFFER`.

**Ditambahkan sebagai Uji AL**, dan pertanyaannya **bukan** *"apakah berbeda"* melainkan **bentuk
apa** yang muncul di masing-masing.

---

## 4. AKIBAT YANG DITARIK SEKARANG, apa pun putusannya

### 4.1 Uji AC menghitung dua benda sebagai satu — kepalanya diperbaiki

Uji AC menghitung **lompatan nomor revisi** untuk mengukur addendum yang dihapus. Karena `OLDID`
juga dipakai untuk salinan dan untuk nomor penawaran, penelusuran rantai yang bersandar padanya akan
mencampur ketiganya.

**Yang menyelamatkannya:** nomor urut revisi **tidak** dibaca dari `OLDID`, melainkan dari **format
`ID` itu sendiri** (`/R01`, `/R02`) — `TreatyInRevisi_post` membangkitkannya dengan
`@substring(ID,10,12)`. Jadi Uji AC tetap sah **asalkan ia menyaring pada pola `ID`, bukan pada
keterisian `OLDID`**. Kepalanya sudah ditulis ulang dengan peringatan itu.

### 4.2 `TREATYINOFFER` — satu sebab lagi kenapa isinya tidak dapat dipercaya

Pembacanya menulis `NOOFFER = OLDID` tanpa syarat, sedangkan `OLDID` berarti tiga hal. Maka pada
kontrak hasil salinan dan pada addendum, kolom `NOOFFER` berisi **bukan nomor penawaran**.

Ini bertemu dengan temuan `CABANG-MATI-DI-JALUR-PERSETUJUAN.md` §2.5: langkah yang memanggil
`SaveTreatyInOffer_Act` **mati**. Jadi kolom itu **berhenti diisi**, dan **isi yang sempat masuk
sebelum berhenti tidak seragam artinya.**

> **Naik ke butir eskalasi 3** — *"tabel yang dibaca sistem lain tidak punya pengisi yang hidup"* —
> sebagai kalimat kedua: **bukan hanya berhenti diisi, isinya pun bercampur tiga arti.** Siapa pun
> yang membacanya perlu tahu keduanya.

---

## 5. Dan sapuan ini menemukan **BENTUK PENULIS KELIMA**

Sapuan pertama atas `OLDDATA` mengembalikan **nol penulis**, dan itu **salah**. Penulisnya berbentuk
**salin halaman**, yang sasarannya duduk di parameter langkah, bukan di elemen penugasan:

```xml
<pyStepsCallParams>
  <CopyFrom>TreatyIntemp</CopyFrom>
  <CopyInto>TreatyIn.OLDDATA</CopyInto>
</pyStepsCallParams>
```

| | |
|---|---:|
| langkah salin-halaman di korpus | **64** |
| di antaranya menyasar halaman potret kontrak | **4** |

> **Ini instans ketiga dari aturan `CONTEXT.md` §2.0-j dalam satu hari**, dan ia yang paling tajam:
> dua yang pertama adalah **nama tag yang berbeda untuk gagasan yang sama**; yang ini **bentuk yang
> berbeda sama sekali** — penugasan lewat parameter metode, bukan lewat pasangan nama/nilai.
>
> `alat/sapu-penulis-properti.py` sudah diperluas ke bentuk kelima, dan
> `BENTUK-PENULIS-PROPERTI.md` diperbarui dari "empat bentuk" menjadi **lima**.
