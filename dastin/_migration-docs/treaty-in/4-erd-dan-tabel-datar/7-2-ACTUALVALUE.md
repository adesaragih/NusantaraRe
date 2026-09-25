# 7.2 — `ActualValue` memotret seluruhnya atau sebagian?

**Tanggal:** 24 September 2026
**Jawaban:** **SELURUHNYA** — dan untuk **dua dari tiga jenis addendum**. Jenis ketiga
(Addendum Premium) **tidak dipotret sama sekali**, karena langkah penyalinnya **dilompati**.

> ### KOREKSI 24 September 2026 — arah prasyaratnya terbalik, dan kalimat di bawah ikut terbalik
>
> Berkas ini semula berbunyi *"hanya untuk satu dari tiga jenis; dua jenis lainnya mendapat potret
> kosong"*. **Kebalikannya yang benar.** Uraian lama **dipertahankan di tempatnya** sebagai alasan,
> dengan penanda di masing-masing bagian — sebab yang keliru bukan pengamatannya, melainkan **arah
> pembacaan satu prasyarat**.
>
> Langkah 2 `SaveTreatyIn_EDM_Act` **hidup**, dan kondisinya memang `TreatyIn.EDMState=="3"`.
> Yang tidak dibaca adalah **kode tindakannya**:
>
> ```xml
> <pyStepsPreCondParamsWhen>TreatyIn.EDMState=="3"</pyStepsPreCondParamsWhen>
> <pyStepsPreCondParamsWhenTrue>1</pyStepsPreCondParamsWhenTrue>
> <pyStepsPreCondParamsWhenTruePrms>jmp</pyStepsPreCondParamsWhenTruePrms>
> ```
>
> dan **langkah 5 ber-`pyStepsBlockName = jmp`**.
>
> **Arti kode `1` dikalibrasi, bukan diandaikan** — disapu atas kedua ekspor, seluruh
> `pyStepsPreCondParamsWhenTrue`/`WhenFalse`:
>
> | Kode | Jumlah | Membawa nama blok? |
> |---:|---:|---|
> | **1** | **66** | **66 dari 66 — selalu** |
> | 2 | 18.412 | tidak pernah |
> | 3 | 4.890 | tidak pernah |
> | 4 | 412 | tidak pernah |
> | 5 | 16 | tidak pernah |
>
> Kode `2` muncul pada langkah **tanpa kondisi sama sekali** — yaitu "lanjut". **Maka kode `1` =
> lompat ke blok**, dan pada 46 dari 66 pemakaiannya nama blok itu ada di aktivitas yang sama.
>
> **Akibatnya: untuk `EDMState == "3"`, langkah 2, 3, DAN 4 dilompati seluruhnya** — lompatannya
> mendarat di langkah 5. Rinciannya `_migration-docs/treaty-in/2-to-spec/TEMUAN-0-5-ACTUALVALUE-TERBALIK.md`.
>
> **20 dari 66 pemakaian kode `1` menunjuk nama blok yang tidak ditemukan di aktivitas yang sama.**
> Sebabnya belum diperiksa, dan **tidak** dihitung sebagai cacat di sini.

> **Saya menjawab "sebagian" lebih dulu, dan itu salah.** Sapuan penugasan ruas demi ruas
> menghasilkan 8 cabang pokok, dan angka itu **benar untuk bentuk yang disapunya** — tetapi
> potretnya tidak dibuat ruas demi ruas. Ia **salin halaman**, bentuk kelima yang baru ditemukan
> hari ini (§7.1 berkas `7-1-OLDID.md` §5). Urutan koreksinya ditulis di §4.

---

## 1. Mesinnya, dibaca dari `SaveTreatyIn_EDM_Act.xml` langkah 2–4

```
2  Page-Copy     TreatyIn      -> TreatyInTemp          when TreatyIn.EDMState == "3"
                 label: "when edmtype= 2 copy data actual to treatyin.a…"     <- LABEL BOHONG
3  Page-Remove   TreatyInTemp.OLDDATA
                 TreatyInTemp.ValueDifference           TANPA PRASYARAT
4  Page-Copy     TreatyInTemp  -> TreatyIn.ActualValue  TANPA PRASYARAT
```

**Tiga hal terbaca sekaligus, dan ketiganya berakibat.**

### 1.1 Potretnya UTUH, bukan pilihan ruas

`ActualValue` menerima **seluruh halaman `TreatyIn`**, dikurangi dua cabang yang dibuang lebih dulu.
Maka **tidak ada daftar "besaran yang dipotret"** — yang ada daftar **yang dikecualikan**, dan
isinya dua: `OLDDATA` dan `ValueDifference`.

> **Akibatnya pada butir 5.6 (bentuk tabel nilai selisih): tabel selisih dapat menampung selisih
> atas besaran APA PUN**, karena sisi "sekarang" tidak pernah membatasi. Yang membatasi adalah sisi
> **penghitung selisih**, bukan sisi potret — dan itu pertanyaan yang berbeda, dijawab §3.

> ### DIFF `D-7` SENGAJA BELUM DITERAPKAN — pernyataan keputusan, bukan `TODO`
>
> | | |
> |---|---|
> | **Apa yang sengaja belum dikerjakan** | kalimat untuk **pembaca arsip `JSONDATA`** belum ditambahkan ke berkas ini — bahwa `ActualValue` **berlapis**, dan pelapisannya tumbuh pada jenis **1 dan 2**, bukan pada addendum premi |
> | **Kenapa** | syarat 4 penerapan diff menetapkan `D-7` diterapkan **bersamaan dengan koreksi `F-1`**, tidak sebelumnya — sebab menambah kalimat yang benar ke berkas yang §1.2-nya terbalik menghasilkan berkas yang **bertentangan dengan dirinya sendiri** |
> | **Keadaan sebenarnya, 24 September 2026** | **`F-1` SUDAH DITERAPKAN** — §1.2, §1.3, dan §1.4 berkas ini sudah dikoreksi di bawah butir `P-1`. **Penahan `D-7` sudah lewat.** Dilaporkan, dan **parkirnya tidak dicabut sendiri**: pencabutan itu wewenang pemilik proses, sebab perintah penerapan berbunyi *"lima sekarang, dua diparkir"* |
> | **Akibat selama diparkir** | pembaca arsip `JSONDATA` lama **tidak diberi tahu** bahwa `ActualValue` berlapis, dan akan membaca potret yang salah — potret terluar bukan potret versi itu |
> | **Di mana terlihat** | §1.4 berkas ini menjelaskan pelapisannya, tetapi **tidak disimpan bersama arsipnya**; pembaca arsip tidak membuka berkas ini |
> | **Ditagih oleh** | satu kalimat izin dari pemilik proses. **Tidak ada pekerjaan lain yang tersisa** |

### 1.2 Labelnya berbohong — `edmtype 2` versus `EDMState == "3"`

Instans kelima dari *"label tidak dipercaya"* di modul ini. Dan bukan sekadar salah angka:
**`EDMState` dan `EDMMaterialType` adalah dua properti yang berbeda**, dan label itu menyebut nama
yang satu dengan nilai yang lain.

Pemetaan `EDMState` yang sebenarnya, dari `DataTransform/TreatyInSetEditPre.xml`:

| `EDMState` | Peristiwa yang dicatatnya | Dipotret? ~~semula~~ **koreksi 24 Sep** |
|---|---|---|
| `1` | *"Had Created **Internal Edit**"* — juga menyetel `EDMEffective = Commencement` | ~~TIDAK~~ **YA** |
| `2` | *"Had Created **External Addendum**"* | ~~TIDAK~~ **YA** |
| `3` | *"Had Created **Addendum Premium**"* — juga `ActualValue.EGNPI = EGNPI`, `AddendumPremi = "1"` | ~~YA~~ **TIDAK — dilompati** |

> **Kolom ketiga terbalik; dua kolom pertama benar** dan dikonfirmasi pemilik proses 24 September:
> `1` = Internal, `2` = External, `3` = Adjustment Premium. Yang keliru bukan arti nilainya.

### 1.3 Dan langkah 3–4 berjalan TANPA prasyarat — ini cacatnya

Prasyarat `EDMState == "3"` hanya melekat pada **langkah 2**. Langkah 4 menyalin `TreatyInTemp` ke
`ActualValue` **selalu**.

> ### ~~Uraian semula~~ — DICORET 24 September 2026, dan penggantinya di bawahnya
>
> ~~Untuk addendum jenis 1 dan 2, `TreatyInTemp` tidak pernah diisi pada permintaan itu — dan
> isinya tetap disalin menimpa `ActualValue`. Potretnya menjadi kosong atau basi, bukan
> "tidak dibuat".~~
>
> **Yang sebenarnya terjadi — arahnya terbalik.** Prasyarat pada langkah 2 **melompat** ketika
> `EDMState == "3"`, dan lompatannya mendarat di langkah 5 — sehingga **langkah 3 dan 4 ikut
> terlewat**. Maka:
>
> | Jenis | Langkah 2-3-4 | `ActualValue` pada penyimpanan |
> |---|---|---|
> | `1` Internal · `2` External | **berjalan** | **ditimpa salinan utuh pohon utama, setiap kali simpan** |
> | **`3` Addendum Premium** | **dilompati** | **tidak disentuh** — isinya tetap apa yang disunting pengguna |
>
> **Dan untuk jenis 3 itu bukan cacat, melainkan justru yang membuatnya bekerja:** `ActualValue`
> adalah pohon yang penggunanya sunting, disemai `TreatyInSetEditPre` langkah 3.5. Menimpanya pada
> setiap simpan akan menghapus masukan penggunanya sendiri.
>
> **Cacat yang sebenarnya ada di tempat lain, dan lebih berat:** mesin selisih premi
> `TreatyEDMDifferencePremium` **tidak pernah membaca `ActualValue`** — ia beriterasi atas
> `TreatyIn.EGNPI`. Maka pada addendum premi **selisih EGNPI selalu nol**. Dicatat sebagai
> **TDA-16** di `../../treaty-in-adjustment/PENGETAHUAN.md` §5.4.
>
> Bentuknya sudah dikenal di proyek ini: **langkah berlabel benar, berprasyarat benar, diikuti
> langkah tanpa prasyarat yang membatalkan maksudnya** — sekeluarga dengan langkah berjudul
> *"if empty exit activity"* yang tidak memuat kode keluar.

**Ini tidak dapat dipastikan dari ekspor saja** — apakah `TreatyInTemp` kosong atau membawa sisa
tergantung keadaan clipboard saat itu, dan clipboard tidak ada di ekspor. **Yang dapat dipastikan:
tidak ada penjaga.** Ujinya data, dan ia tegas: **Uji AM**.

### 1.4 Potretnya BERSARANG, dan ia tumbuh tiap addendum

Langkah 3 membuang `OLDDATA` dan `ValueDifference` dari salinannya — **tetapi tidak membuang
`ActualValue`**. Maka potret yang baru memuat potret yang lama di dalamnya:

```
ActualValue( ActualValue( ActualValue( … ) ) )
```

Dalamnya bertambah satu tiap kali addendum **jenis 1 atau 2** dibuat — **bukan jenis 3**, yang
langkah penyalinnya dilompati — dan seluruhnya tersimpan di `M_TREATY_IN.JSONDATA` sebagai satu
dokumen.

> **Dikoreksi 24 September 2026.** Semula tertulis *"tiap kali addendum jenis 3 dibuat"*. Ujinya
> **Uji AM-3** tetap sah; yang berubah **kontrak yang disasarnya** — yang diukur adalah kontrak
> dengan banyak revisi internal dan addendum eksternal, bukan kontrak yang sering di-addendum-premi.

> **Akibat yang dapat diukur:** ukuran `JSONDATA` tumbuh **berlipat**, bukan bertambah, pada kontrak
> yang sering di-addendum-premi. **Uji AM-3** menghitungnya, dan ia juga menjawab §1.3 sekaligus —
> kedalaman sarang nol pada kontrak yang pernah di-addendum berarti potretnya memang tidak pernah
> terisi.
>
> **Akibat pada model: tidak ada, dan itu disengaja.** ADR-0034 menetapkan `JSONDATA` menjadi
> **arsip**, dan model baru menyimpan potret sebagai **baris nilai**, bukan sebagai halaman
> bersarang. Sarang itu **tidak dibawa** — dan karena ia bukan pelestarian, ia tidak menuntut tiket.
> Yang menuntut perhatian hanyalah **pembacaan arsipnya**: siapa pun yang kelak membaca `JSONDATA`
> lama harus tahu bahwa `ActualValue` di dalamnya berlapis.

---

## 2. Empat halaman potret, dan bentuk pengisian masing-masing

| Halaman | Diisi dengan | Penulisnya |
|---|---|---|
| **`ActualValue`** | **salin halaman utuh** + 8 cabang disunting ruas demi ruas sesudahnya | `SaveTreatyIn_EDM_Act` (utuh) · 14 aktivitas/DT (ruas) |
| **`OLDDATA`** | **salin halaman utuh** | `TreatyInSetAddendumToHistory` ← dipanggil `TreatyInRevisi_post` langkah 1 |
| **`ValueDifference`** | **ruas demi ruas saja** — tidak ada salin halaman | `TreatyEDMDifference*` (4 berkas) + 18 lainnya |
| **`ValueBeforeProrate`** | **nol penulis** atas kelima bentuk; hanya muncul sebagai kendali layar **mati tanpa syarat** | — |

> **`TreatyInSetAddendumToHistory` ber-`pyUsage` = "Addendum will be set on the history property
> (unused)"** — dan ia **dipanggil langkah hidup**. Label berbohong lagi, kali ini pada ruas
> "kegunaan" aturannya sendiri.

**Dan `OLDDATA` memang tidak perlu disimpan ulang**, sejalan dengan aturan pemilik proses: sisi
"lama" di-`SELECT` dari versi sebelumnya. Yang berubah hanya **dasarnya** — sebelumnya keputusan
rancangan, sekarang juga terbaca bahwa sistem lama membuatnya sebagai **salinan utuh** yang memang
tidak menambah keterangan apa pun di luar versi pendahulunya.

---

## 3. `BESARAN_DAPAT_DISESUAIKAN` — dijawab dari sisi PENGHITUNG SELISIH, bukan dari potret

Karena potretnya utuh, daftar besaran yang dapat disesuaikan **tidak** dapat dibaca dari
`ActualValue`. Ia terbaca dari `ValueDifference`, yang diisi ruas demi ruas:

| Halaman | akar | **agregat** (`Total…`, `…SummaryList`) | **pokok** |
|---|---:|---:|---:|
| `ActualValue` — sisi ruas | 29 | 21 | **8** |
| `ValueDifference` | 31 | 22 | **9** |

**Kesembilan besaran pokok yang punya penghitung selisih:**

```
BrokeragePercent · EGNPI · FacultativeShare · FacultativeShareBrokerage
FacultativeShareList · Installment · Limits · RNMShare · Share
```

Perbandingan kedua daftar, dan selisihnya **satu**:

| | |
|---|---|
| hanya di `ValueDifference` | **`Installment`** — termin punya penghitung selisih, tetapi tidak punya penyuntingan ruas di sisi potret |
| hanya di `ActualValue` | **nol** |

> **Maka `BESARAN_DAPAT_DISESUAIKAN` bersumber pada daftar `ValueDifference`, sembilan besaran** —
> dan itu sumber yang lebih kuat daripada *"31 akar bercermin"*, karena 22 dari 31 akar itu
> **agregat** yang menurut ADR-0037 tidak disimpan sama sekali.
>
> **`Installment` ditandai**, bukan dirapikan: ia satu-satunya yang punya sisi selisih tanpa sisi
> potret bersuntingan. Sejalan dengan §10.16 yang menemukan tingkat termin terbalik, dan **ditagih
> saat tabel nilai selisih digambar.**

---

## 4. Urutan koreksi di berkas ini, supaya pembaca berikutnya tidak mengulang jalannya

| Langkah | Hasil | Keadaan |
|---|---|---|
| sapuan penugasan ruas demi ruas, empat bentuk | `ActualValue` = 8 cabang pokok → **"SEBAGIAN"** | **salah** — bentuk penulisnya bukan itu |
| memeriksa `OLDDATA` dengan sapuan yang sama | **nol penulis** | **salah**, dan justru inilah yang memaksa mencari bentuk lain |
| sapuan `Page-*` atas nama halaman | **nol** | **salah lagi** — sasarannya di `pyStepsCallParams`, bukan di `pyParamArray` |
| sapuan `<CopyFrom>`/`<CopyInto>` | 64 langkah, 4 menyasar halaman potret | **benar** |

> **Yang mengungkapnya bukan kecurigaan melainkan angka nol yang mustahil.** `OLDDATA` punya layar
> sendiri di empat berkas `…OldData.xml`; sesuatu yang ditampilkan di empat layar **tidak mungkin**
> tidak punya penulis. **Nol yang mustahil adalah alat cari**, dan ia layak berdiri di samping
> "tetapan yang ditulis langsung di kode".
