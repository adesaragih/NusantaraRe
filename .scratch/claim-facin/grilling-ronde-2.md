# Grilling ronde 2 — Claim Fac In

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Fac In\` — **482 berkas `.xml`**, READ-ONLY.
**Berkas ini MENAMBAH.** ⛔ `grilling-ronde-1.md` dan `grilling-ronde-1-ulang-docs.md` tidak disentuh
satu byte pun.

> **Cara membaca angka di berkas ini.** Setiap sensus menyebut **jendelanya** — apa yang dibaca, apa
> yang dikecualikan, dan penyebutnya. Angka yang saya dapat **berbeda** dari ronde sebelumnya
> **dilaporkan berdampingan**, tidak ditimpa. Bukti selalu **path berkas + nama rule + nomor langkah
> Pega**; ⛔ **nol nomor baris XML**.

| Tanda | Artinya |
| --- | --- |
| `[terverifikasi]` | dibaca sendiri dari korpus dengan pengurai XML; buktinya disebut |
| `[terbuka]` | **belum terjawab** — jangan ditebak |
| ⚠️ | jebakan atau risiko yang mudah terlewat |
| ⭐ | temuan yang mengubah gambaran |
| **USULAN** | belum disahkan siapa pun |

---

## LANGKAH 0

```
0a  isi .scratch\claim-facin\           2 berkas   ✅ grilling-ronde-2.md belum ada
0b  berkas .xml di korpus               482        ✅
0c  berkas seluruhnya di korpus         483        ✅  selisihnya .xlsx — §A2
0d  OUTPUT_HASIL_RNM\ATURAN-BACA-KORPUS-PEGA.md    TIDAK ADA   ✅ sesuai harapan — §A1
0e  docs\adr\*.md                       15         ✅
```

---

## §A — Dasar baca, dan dua angka yang diluruskan

### A1 — Aturan baca yang benar-benar dipakai

⛔ **`OUTPUT_HASIL_RNM\ATURAN-BACA-KORPUS-PEGA.md` TIDAK ADA** — dicari di seluruh
`D:\XML\RNM_BRD\`, nol hasil. Kalimat ronde 1 §A1 *"Aturan A–J dipakai sejak langkah pertama"*
karena itu **menunjuk berkas yang tidak ada**.

**Yang benar-benar saya pakai, disebut apa adanya:**

| Dasar | Sumber | Status |
| --- | --- | --- |
| Sembilan aturan baca ekspor Pega | bab `## Lampiran — Aturan baca ekspor Pega` di `komite-claim-prop\spec.md` | dipakai |
| Aturan **DataPage** | ronde 1 §A2 | **USULAN** |
| Empat butir tambahan | ronde ulang §C3·C4 | **USULAN** |

⛔ **Usulan tetap usulan.** Saya **tidak mengesahkannya**. Yang berhak mengesahkan adalah **work
owner**; sampai itu terjadi, setiap pernyataan yang bersandar padanya saya tandai **USULAN**.

⚠️ **Akibat yang perlu diketahui:** ronde 1 dan ronde ulang menyebut *"aturan A–J"* dengan huruf,
sedangkan lampiran yang ada memberi **sembilan aturan bernomor**, bukan berhuruf. **Pemetaan
huruf → nomor tidak ada di mana pun.** `[terbuka]`

### A2 — Angka korpus: **482 SAH**, dan berkas ke-483 bukan rule

`[terverifikasi]` Folder korpus berisi **483 berkas**: **482 `.xml`** ditambah **satu**
`Struktur_Register_Flow.xlsx`.

⭐ **Angka 482 milik ronde 1 SAH**, karena ia menghitung **rule**. Yang keliru hanyalah
**penyebutnya**: ronde 1 menulis *"482 berkas"*, seharusnya **"482 berkas `.xml`"**.

⛔ Berkas `.xlsx` itu **artefak turunan buatan tim sendiri untuk memetakan XML, bukan sumber
kebenaran** — ⛔ tidak boleh dikutip `[terverifikasi]`, dan **tidak saya buka**. Pola yang sama ada
di `Claim Prop` *(`Struktur_Flow_TreatyIn.xlsx`)*, `Komite Claim Prop`, dan `Master Contract Retro
Life`.

### A3 — ✅ Kelima pertanyaan lama: **KELIMANYA TERJAWAB** — `[keputusan work owner]` 2026-09-19

| # | Asal | Pertanyaan | Jawaban |
| --- | --- | --- | --- |
| 1 | ronde 1 §F5 | **Empat rule beda versi** — yang mana yang berlaku di produksi? | ✅ **A — versi Fac In yang berlaku** |
| 2 | ronde 1 §F5 | `pyStepsPreCondition = 0` dan `pyStepsRepeatDefHasRepeat = PROPERTYLIST` | ✅ **A — setara dengan nilai yang sudah dikenal** |
| 3 | ronde 1 §F5 | `pyWorkPage.Quotation` — ada atau tidak? | ✅ **ADA** — dengan satu ketetapan migrasi, lihat §A4 |
| 4 | ronde ulang §H4 | Enam ADR ditentang korpus Fac In | ✅ **A — ADR berlaku lintas-modul** |
| 5 | ronde ulang §H4 | 60 rule `When` identik, vonisnya berlawanan | ✅ **A — satu rule, satu nasib** |

⭐ **Uraian lengkap keenam jawaban — termasuk yang keenam dari §B4 — ada di §A4.**

⭐ **Pertanyaan 2 kini punya angka** `[terverifikasi]` — sensus atas **179 berkas `Activity`**:

```
pyStepsPreCondition        true 942 · false 130 · (kosong) 1431 · ⭐ "0" 6
pyStepsRepeatDefHasRepeat  (kosong) 1988 · EMBEDDED 409 · REPEAT 111 · ⭐ PROPERTYLIST 1
```

⛔ **Angka bukan jawaban.** Apa **arti** `0` dan `PROPERTYLIST` tetap tidak tertulis di mana pun.
Pertanyaan 2 **tetap terbuka**, dan work owner belum mengonfirmasi apa pun.

---

## §A4 — ✅ Enam jawaban work owner, lengkap dengan akibatnya

`[keputusan work owner]` **2026-09-19**. ⛔ Kalimat aslinya dikutip; akibat yang saya catat di
bawahnya adalah **konsekuensi yang terbaca**, bukan keputusan tambahan.

### A4-1 — Versi Fac In yang berlaku

> **"a"** — *A: versi Fac In yang berlaku.*

Keempat rule bersama dibaca menurut **versi Claim Fac In**, yang **lebih baru** pada keempatnya.

⚠️ **Akibat yang menyentuh modul lain:** kesimpulan **Claim Prop** tentang **kurs standar** dibuat
di atas rule yang **4 tahun 3 bulan lebih tua**, dan kini **perlu dibaca ulang**. Begitu pula
kesimpulan Claim Prop tentang **penyimpanan berkas**. ⛔ **Modul itu tidak disentuh dari sini** —
dicatat sebagai pekerjaan yang lahir dari jawaban ini.

### A4-2 — `0` setara kosong, `PROPERTYLIST` setara `EMBEDDED`

> **"A — setara dengan yang sudah dikenal"**

`pyStepsPreCondition = 0` dibaca **sama dengan kosong** *(tidak ada gerbang)*, dan
`pyStepsRepeatDefHasRepeat = PROPERTYLIST` dibaca **sama dengan `EMBEDDED`**.

✅ **Akibatnya melegakan:** sensus gerbang dan perulangan ronde 1 dan ronde 2 **tidak perlu
diulang** — keenam langkah ber-`0` memang sudah dihitung sebagai tanpa gerbang, dan satu langkah
`PROPERTYLIST` sudah dihitung sebagai berulang.

⚠️ **Tetapi catatannya belum punya rumah.** Aturan baca kanonnya **tidak ada berkasnya** *(§A1)*.
Sampai ada, ketetapan ini **hanya tercatat di sini**, dan modul berikutnya yang disensus **tidak
akan menemukannya**. `[terbuka]` — di mana catatan aturan baca disimpan.

### A4-3 — `Quotation` ADA, tetapi sumbernya berubah saat migrasi

> **"ada, tapi saat migrasi quotation langsung select dati table polis tidak di
> pyWorkPage.Quotation"**

**Dua hal sekaligus, dan keduanya mengikat:**

1. ✅ **`pyWorkPage.Quotation` ADA** pada objek kerja Claim Fac In — jadi **ke-38 rule `When` yang
   mengujinya HIDUP**, bukan sisa impor. ⛔ **Kebalikan** dari Claim Prop.
2. ⭐ **Di sistem baru, `Quotation` diambil LANGSUNG dari tabel polis — bukan dari halaman
   `pyWorkPage.Quotation`.**

⚠️ **Butir 2 adalah penyimpangan sadar**, karena sumber datanya berubah. ⛔ Ia **wajib masuk daftar
penyimpangan sadar** — lihat §H1b, tempatnya kini **dua**, bukan satu.

⛔ **Yang TIDAK diputuskan:** nama tabel polisnya, kolomnya, dan apakah ke-38 rule `When` itu ikut
dipindahkan atau digantikan pembacaan langsung. Ketiganya **belum ditetapkan**. `[terbuka]`

### A4-4 — ADR berlaku lintas-modul

> **"A — ADR berlaku lintas-modul"**

Kelima belas ADR **mengikat seluruh modul**, bukan hanya Claim Life.

⚠️ **Akibat di Fac In:** keenam pertentangan — **0003 · 0005 · 0008 · 0011 · 0013 · 0015** —
menjadi **penyimpangan sadar yang harus dicatat di spec Fac In kelak**, dan sebagian **menyentuh
uang**.

⚠️ **Akibat yang jauh lebih besar, dan menyentuh modul yang sudah selesai:** **Claim Prop** dan
**Komite Claim Prop** sudah **tuntas spec dan tiketnya tanpa pernah sekali pun diadu dengan 15 ADR**
— ronde ulang §H5 butir 6 sudah menandainya dan tidak mengerjakannya. Dengan jawaban ini,
mengadu keduanya **berubah dari "sebaiknya" menjadi "harus"**. ⛔ **Keduanya tidak disentuh dari
sini.**

### A4-5 — Satu rule, satu nasib

> **"A — satu rule, satu nasib"**

Rule yang **identitas empat bagiannya sama** diperlakukan **sama di semua modul**. ⛔ **Nol salinan
ganda logika klasifikasi** di sistem baru.

⚠️ **Akibat langsung, dan ia bertabrakan dengan keputusan yang sudah ada:** ke-60 rule `When` itu
**identik** dengan Claim Prop. Karena **A4-3** menyatakan halaman yang diujinya **ADA** di Fac In,
rule itu **hidup** — sehingga keputusan Claim Prop *"52 dari 61 tidak dimigrasikan karena halamannya
tidak ada"* **perlu ditinjau ulang**.

⛔ **Saya tidak meninjaunya di sini, dan tidak menyatakan keputusan Claim Prop batal.** Yang saya
nyatakan: **kedua keputusan tidak dapat berdiri bersama tanpa penjelasan**, dan penjelasan itu
belum ada. `[terbuka]`

### A4-6 — Baris Travel lama dibiarkan apa adanya

> **"A — dibiarkan apa adanya"**

Baris **Travel** yang sudah tersimpan **dipindahkan seperti aslinya**. ⛔ **Tidak dihitung ulang,
tidak ditambal, tidak ditolak.**

> ⚠️ **Risiko yang diterima sadar, ditulis tegas:** sejak perbaikan logika berjalan *(§B4)*,
> **baris Travel baru akan berbeda isinya dari baris Travel lama**. Laporan yang menjumlahkan
> keduanya **mencampur dua perilaku tanpa satu pun penanda yang membedakannya**. Selisih itu
> **bukan cacat migrasi** — ia akibat langsung dari keputusan yang benar, dan **tidak ada yang
> memberitahu pembaca laporan**.

⛔ **Yang TIDAK diputuskan:** apakah perlu ada penanda yang membedakan baris sebelum dan sesudah
perbaikan. `[terbuka]`

---

## §A5 — ✅ Tujuh keputusan rancangan — 2026-09-19

`[keputusan work owner — atas rekomendasi asisten]` **2026-09-19**, kecuali yang ditandai lain.
⛔ Tanda *"atas rekomendasi asisten"* dipakai **sengaja**: pilihannya datang dari rekomendasi saya,
dan **pencabutannya harus murah**. ⛔ **Nol di antaranya mengubah korpus** — seluruhnya tentang
sistem baru.

| # | Keputusan | Golongan |
| --- | --- | --- |
| **1** | **Daftar penyetuju kosong → penyerahan GAGAL TERANG-TERANGAN**, bukan tangga nol tingkat yang senyap | ⭐ **PERBAIKI** |
| **2** | **`.KomiteNo` diambil dari pengenal kasus komitenya**, bukan dari potongan teks mulai huruf ke-19 | ⭐ **PERBAIKI** |
| **3** | **`FlagOnGoingCommitte` tidak disimpan** — diturunkan dari kasus komitenya | ⭐ **PERBAIKI** |
| **4** | **`Quotation` dibaca SEKALI saat kasus dimuat** sebagai satu nilai turunan, dan **ke-38 rule `When` diganti SATU fungsi klasifikasi lini** | ⭐ **PERBAIKI** |
| **5** | **Alias SQL diberi nama sesuai isinya di batas pembacaan**, disertai **tabel pemetaan nama-lama → nama-benar** | ⭐ **PERBAIKI** |
| **6** | **Rule dipisahkan dari KEHIDUPANNYA** — satu rule dipindahkan sekali; hidup-matinya bergantung pada ada-tidaknya halaman yang diujinya pada objek kerja modul itu | ⭐ **PERBAIKI** |
| **7** | **Panggilan arasapas TIDAK MATI** — dan pemanggilnya **empat**, nol bergerbang hidup | ⚠️ **fakta**, `[keputusan work owner]` |

### A5-1 — Daftar penyetuju kosong: gagal terang-terangan

`CreateKMTNo_Act` langkah **9** mengisi `KomiteLoop` dari ukuran daftar dan `KomiteCount := 1`
**tanpa gerbang**. Bila daftarnya kosong, Pega melahirkan **kasus komite tanpa satu pun penyetuju,
tanpa pesan apa pun**.

⭐ **Di sistem baru, penyerahan itu DITOLAK dengan galat yang terlihat.**

⚠️ **Penyimpangan sadar.** ✅ **Sejalan dengan preseden yang sudah ada**: Komite Claim Life
**AC 4** menetapkan hal yang sama untuk bentuk yang sama.

### A5-2 — `.KomiteNo` dari pengenal kasus, bukan potongan huruf ke-19

Pega menulis `.KomiteNo` **dua kali**: langkah **13.1** dari `childPageKomite.pyID`, lalu
**13.1.1** dari `@substring(pyWorkPage.pxCoveredInsKeys(<last>),19)`.

⭐ **Di sistem baru dipakai yang pertama — pengenal kasus komitenya.** Angka **19** hanyalah panjang
awalan kelas; begitu nama kelas berubah satu huruf, potongan itu **diam-diam salah**.

⚠️ **Yang tetap harus dicari, dan belum ditemukan:** **kenapa 13.1 ditimpa 13.1.1.** Dua penulisan
berturut ke properti yang sama menandakan yang pertama dianggap kurang. ⛔ Alasannya **tidak ada di
korpus**. `[terbuka]`

### A5-3 — `FlagOnGoingCommitte` diturunkan, tidak disimpan

Pega mengisinya dengan **teks `"Send Commite"`** di `CreateKMTNo_Act` langkah **3**, tanpa gerbang.

⭐ **Di sistem baru kolom itu tidak dibuat** — keadaan *"komite sedang berjalan"* **diturunkan dari
kasus komitenya**.

✅ **Sejalan dengan keputusan 27 modul Komite Claim Prop**, yang membuang kolom bernama sama dengan
alasan yang sama.
⛔ **Catatan jujur, dan ia membatasi keputusan ini:** saya **belum membuktikan** kolom Fac In dan
kolom Claim Prop adalah kolom yang sama — **nama sama bukan bukti**. Bila ternyata berbeda,
keputusan ini **perlu ditinjau**. `[terbuka]`

### A5-4 — `Quotation` dibaca sekali, dan 38 `When` menjadi satu fungsi

⭐ **`Quotation` diambil LANGSUNG dari tabel polis, dibaca SEKALI saat kasus dimuat, disimpan
sebagai satu nilai turunan.** ⭐ **Ke-38 rule `When` yang mengujinya diganti SATU fungsi klasifikasi
lini.**

**Alasannya:** ke-38 rule itu menguji **satu halaman yang sama** untuk memilah lini; memindahkannya
satu per satu berarti membawa **38 potong logika yang berisi satu keputusan**.

✅ **Nama tabel dan kolomnya: DITUNDA ke waktu migrasi** — `[keputusan work owner]` 2026-09-19:
**"tabel & kolom polis sumber Quotation dibuat nanti saat migrasi, lewatkan itu, jangan jadiin
permasalahan."**

⛔ **Ia bukan pemblokir.** Keputusan A5-4 **berdiri penuh** tanpa nama tabelnya: `Quotation` dibaca
**sekali** saat kasus dimuat sebagai satu nilai turunan, dan ke-38 rule `When` diganti **satu fungsi
klasifikasi lini**. Yang ditunda hanyalah **dari kolom mana nilainya diambil**.

⚠️ **Satu catatan untuk yang kelak menulis spec, bukan keberatan:** karena tabel dan kolomnya
**dibuat saat migrasi**, kalimat *"diambil langsung dari tabel polis"* berarti **bagian dari rancangan
skema baru**, bukan pembacaan sesuatu yang sudah ada. ⛔ Cukup dicatat supaya tidak dikira sudah
tersedia.

> ⛔ **Kalimat lama, dikutip dan tidak dihapus:** *"⛔ **Yang BELUM ditetapkan, dan harus datang dari
> work owner atau DBA:** **nama tabel polisnya** dan **nama kolomnya**. ⛔ **Saya tidak menebaknya.**
> `[terbuka]`"*

### A5-5 — Alias SQL dibetulkan di batas pembacaan

`CariHistoryClaim_SQL` **5 dari 5** aliasnya berbohong; `BrowseHistoryClaim` **2 dari 3** — dan
⭐ **keduanya berbohong dengan pemetaan yang BERBEDA**.

⭐ **Di sistem baru, kolom diberi nama sesuai isinya di batas pembacaan**, disertai **tabel pemetaan
nama-lama → nama-benar**.

⚠️ **Tabel pemetaan itu WAJIB ada.** Tanpanya, orang yang membandingkan keluaran lama dan baru akan
mengira **datanya berubah**, padahal hanya **namanya** yang dibetulkan.
✅ **Sejalan dengan AC 106 Claim Prop.**

### A5-6 — Rule dipisahkan dari kehidupannya

Jawaban *"satu rule, satu nasib"* bertabrakan dengan keputusan Claim Prop *"52 dari 61 rule `When`
tidak dimigrasikan karena halamannya tidak ada"* — sebab di Fac In halaman itu **ada**.

⭐ **Penyelesaiannya: pisahkan RULE dari KEHIDUPANNYA.** Rule-nya **satu**, dipindahkan **sekali**,
disimpan **satu salinan**. Yang berbeda per modul **bukan rulenya**, melainkan **apakah halaman yang
diujinya ada pada objek kerja modul itu**.

✅ **Dengan begitu kedua keputusan berdiri bersama tanpa saling membatalkan:** Claim Prop tetap benar
*(di sana ia tidak pernah bernilai benar)*, Fac In tetap benar *(di sini ia hidup)*, dan sistem baru
tetap punya **satu salinan logika klasifikasi**.

⚠️ **Bila A5-4 dijalankan, butir ini ikut selesai sendiri** — satu fungsi klasifikasi meniadakan
pertanyaan *"rule mana yang hidup di modul mana"*.

> ✅ **DIKONFIRMASI** — `[keputusan work owner]` 2026-09-19: **"IKUTI REKOMENDASI."**
>
> ⛔ **Akibat yang menyentuh modul yang sudah selesai:** keputusan Claim Prop *"52 dari 61 rule
> `When` tidak dimigrasikan"* **tidak dibatalkan** — ia **dibaca ulang** sebagai pernyataan tentang
> **kehidupan rule di modul itu**, bukan tentang **dipindahkan atau tidaknya rule itu**. ⛔ Kalimat
> di modul Claim Prop **belum disesuaikan**, dan **tidak disentuh dari sini**. Lihat §H1 butir 12.

### A5-7 — Arasapas tidak mati — `[keputusan work owner]`

> **"SAYA CEK DI FACIN TIDAK MATI"**

✅ **Benar, dan berkas ini diralat mengikutinya.** Uraian lengkap berikut kesalahan saya ada di
**§D3**. Ringkasnya: `pyStepsPreCondition = false` mematikan **gerbangnya**, bukan **langkahnya** —
sehingga panggilan arasapas justru berjalan **tanpa saringan**, dari **empat** pemanggil.

---

## §B — `GetAllData_Act` dibaca langkah demi langkah

**Cara baca** `[terverifikasi]`: pengurai XML sungguhan *(`xml.etree`)*, menelusuri `pySteps` →
`rowdata` secara rekursif, sehingga langkah bersarang terbaca dengan nomor bertitik *(`16.1.1.1.2`)*
dan jebakan **`pyStepsPreCondParams` induk yang terbit sesudah `</pySteps>` anaknya** tidak menipu.

⭐ **Dua jangkar, dan angkanya BEDA** *(aturan I3/P4)*:

```
jangkar A  pxObjClass = Embed-ActivitySteps       50
jangkar B  pyStepsActivityNameUC                  46
pengurai XML rekursif                             50
```

⛔ **Keduanya dilaporkan, tidak dipilih salah satu.** Selisih **4** persis sama dengan jumlah langkah
yang **`pyStepsActivityName`-nya kosong** — yaitu langkah **berulang murni** yang tidak memanggil
metode apa pun. **Angka yang dipakai berkas ini: 50**, dengan selisihnya diterangkan.

### B1 — Kerangkanya

| Bagian | Langkah | Isinya |
| --- | --- | --- |
| **Pengambilan outstanding per lini** | **3 – 9** | tujuh `RDB-List`, masing-masing bergerbang satu lini: `IsFire` · `IsAneka` · `isGolfInsurance` · `IsMarineCargo` · `IsMBU` · `IsTravel` · `IsPA` |
| **Pemilah lini** | **11** | tiga lompatan ke tanda `MBU` · `TRAVEL` · `PA` |
| **Tiga cabang lini** | **12 · 13 · 14** | `TRAVEL` · `PA` · `MBU`, masing-masing membaca `OutputData.pxResults` |
| **Perulangan objek** | **16** | `pyWorkPage.ClaimData.ObjectList` → `.ObjectItemList` → **estimasi** *(16.1.1)* dan **adjustment** *(16.1.2)* |
| **Kurs** | **17** | `RDB-List` `SearchNilaiKursOutput` — *"ConvertValue for Currency"* |
| **Konversi nilai** | **18** | dua cabang: **18.2** *"outstanding ANEKA FIRE"* · **18.3** *"outstanding MBU"* |
| **Penutup** | **19** | *"set All Data to property"* |

### B2 — Langkah ber-remark: **SATU**

`[terverifikasi]` Hanya **langkah 10** — metode `Java`, catatan *"Remove Duplikat Data"*.
⚠️ Menariknya **langkah 15 `Java` *"remove currency yg sama"* HIDUP** — jadi pembersih duplikat mata
uang berjalan, sedangkan pembersih duplikat data **dimatikan**. ⛔ Alasannya tidak tertulis.
`[terbuka]`

### B3 — Flag prakondisi

```
true 29 · false 2 · (kosong) 19        jumlah 50
```

Dua yang **`false`** — jadi **gerbangnya tersimpan tetapi tidak berlaku** *(aturan C1)*:
**langkah 14.2** dan **langkah 16.1.1.1**. ⚠️ **16.1.1.1 adalah induk seluruh blok estimasi**;
gerbangnya mati, jadi blok itu **selalu masuk**.

### B4 — ⭐⭐ Kesembilan lompatan — tandanya ADA semua, tetapi **DUA PASANG TERTUKAR**

`[terverifikasi]` **Sembilan lompatan, seluruhnya arah kode 1 *(lompat ke langkah)*, seluruhnya
bergerbang HIDUP, dan kesembilan tandanya ADA — nol yang menggantung.**

| Dari | Gerbang | Melompat ke tanda | Tanda itu ada di | Gerbang langkah pendarat menguji |
| --- | --- | --- | --- | --- |
| **11** | `IsMBU` | `MBU` | **14** | `IsMBU` ✅ |
| **11** | `IsTravel` | `TRAVEL` | **12** | `IsTravel` ✅ |
| **11** | `IsPA` | `PA` | **13** | `IsPA` ✅ |
| **16.1.1.1.1** | `IsMBU` | ⚠️ **`TRAVEL2`** | **16.1.1.1.2** | `IsTravel` ❌ |
| **16.1.1.1.1** | `IsTravel` | ⚠️ **`MBU2`** | **16.1.1.1.3** | `IsMBU` ❌ |
| **16.1.1.1.1** | `IsPA` | `PA2` | **16.1.1.1.4** | `IsPA` ✅ |
| **16.1.2.1.1** | `IsMBU` | ⚠️ **`TRAVEL3`** | **16.1.2.1.2** | `IsTravel` ❌ |
| **16.1.2.1.1** | `IsTravel` | ⚠️ **`MBU3`** | **16.1.2.1.3** | `IsMBU` ❌ |
| **16.1.2.1.1** | `IsPA` | `PA3` | **16.1.2.1.4** | `IsPA` ✅ |

⭐ **Di langkah 11 pasangannya BENAR. Di kedua blok bersarang, MBU dan Travel TERTUKAR.**

**Diperiksa dua kali** `[terverifikasi]`: sekali lewat pengurai rekursif, sekali lagi dengan
membaca medan mentah `pyStepsPreCondParamsWhen` / `…WhenTrue` / `…WhenTruePrms` baris demi baris.
Kedua cara memberi hasil yang sama.

#### ⚠️ Akibatnya tidak simetris — dan hanya **satu lini** yang dirugikan

Urutan langkah di dalam blok adalah **`TRAVEL2` (.2) → `MBU2` (.3) → `PA2` (.4)**.

| Lini | Mendarat di | Yang terjadi |
| --- | --- | --- |
| **MBU** | `TRAVEL2` *(.2)* | gerbangnya `IsTravel` → salah → **dilewati**, lalu **jatuh ke `MBU2` (.3)** → `IsMBU` benar → **berjalan**. ✅ **Selamat karena urutan, bukan karena benar** |
| **PA** | `PA3`/`PA2` *(.4)* | ✅ benar |
| ⭐ **Travel** | `MBU2` *(.3)* | ia **melompati `TRAVEL2` (.2) — langkahnya sendiri**. `MBU2` → `IsMBU` salah → dilewati. `PA2` → `IsPA` salah → dilewati. ⛔ **Tidak satu pun dari ketiganya berjalan** |

⭐ **Kesimpulannya:** untuk lini **Travel**, langkah `TRAVEL2` dan `TRAVEL3` — yaitu pengisian nilai
pada **perulangan estimasi** *(16.1.1)* **dan** **perulangan adjustment** *(16.1.2)* — **tidak
pernah berjalan**. Di **kedua** perulangan.

⚠️ Yang **tidak** terpengaruh: pengambilan outstanding Travel di **langkah 12**, karena pemilah di
langkah 11 pasangannya benar.

⛔ **Saya tidak menyatakan ini cacat yang harus diperbaiki, dan tidak menyimpulkan apa pun untuk
aplikasi Go.** Yang saya nyatakan: **struktur ekspornya berbunyi demikian**, dibaca dua cara.
**Apakah itu memang dikehendaki adalah pertanyaan untuk work owner** — lihat §H pertanyaan baru.

#### ⭐ Seberapa luas dampaknya — pemanggilnya **SATU**

`[terverifikasi]` Penyisiran ke-**482** berkas `.xml` untuk nama `GetAllData_Act`, **tidak peka huruf
besar-kecil**: **satu pemanggil saja**, `Section\InputEstimasiDetail.xml`
*(`ASM-FW-GCNMFW-WORK-PNC!INPUTESTIMASIDETAIL`, ruleset `GCNMFW` versi `01-01-27`)*.

⭐ **Dan cara memanggilnya bukan panggilan activity biasa** — ia **membuka harness sebagai popup**:

```
pyActivity      GetAllData_Act          pyActivityClass  ASM-FW-GCNMFW-Data-Object
pyHarnessName   Outstanding             pyUsingPage      outstanding
pyTarget        popup                   pyReadOnly       Yes
pyWidth 1600 · pyHeight 800            muncul 10 kali di section yang sama
```

⚠️ **Artinya yang terpengaruh lompatan tertukar itu adalah isi layar popup *Outstanding*** yang
dibuka dari layar rincian estimasi — **bukan jalur simpan**.

⛔ **Tetapi ini BELUM membuktikan activity itu tidak menulis balik ke kasus.** Harness-nya
`pyReadOnly = Yes`, sedangkan **langkah 19** berbunyi *"set All Data to property"* dan **sasaran
penulisannya belum saya telusuri**. `[terbuka]` — lihat §H1.

**Identitas empat bagiannya**, supaya tidak tertukar dengan rule bernama sama di modul lain:

```
pxInsName  ASM-FW-GCNMFW-DATA-OBJECT!GETALLDATA_ACT
kelas      ASM-FW-GCNMFW-Data-Object
ruleset    GCNMFW        versi 01-01-07
```

#### ⚠️ Pemeriksaan work owner — dan satu nama yang TIDAK ADA di ekspor

`[data work owner]` 2026-09-19 — **"tercatat dari TRAVEL, TRAVEL1, TRAVEL2."**

⛔ **Ekspor tidak memuat `TRAVEL1`.** `[terverifikasi]` Penyisiran **seluruh 482 berkas `.xml`**
untuk nama blok berawalan `TRAVEL` / `MBU` / `PA` memberi **tepat sembilan**, seluruhnya di
`GetAllData_Act`:

```
TRAVEL   lgk 12            MBU   lgk 14            PA   lgk 13
TRAVEL2  lgk 16.1.1.1.2    MBU2  lgk 16.1.1.1.3    PA2  lgk 16.1.1.1.4
TRAVEL3  lgk 16.1.2.1.2    MBU3  lgk 16.1.2.1.3    PA3  lgk 16.1.2.1.4
```

⚠️ **Nol `TRAVEL1`, nol `MBU1`, nol `PA1`.**

> ### ✅ DITARIK — `[keputusan work owner]` 2026-09-19
>
> **"IYA SAYA SALAH GA ADA TRAVEL1."**
>
> ⛔ **Pernyataan `[data work owner]` di atas DITARIK oleh work owner sendiri.** Ia **tidak dihapus**
> — tetap dikutip sebagai jejak, supaya terbaca bahwa dugaan itu pernah ada dan sudah diselesaikan.
>
> ✅ **Yang berlaku: ekspor dan kenyataan SEJALAN** — tandanya **TRAVEL · TRAVEL2 · TRAVEL3**,
> dan **tidak ada `TRAVEL1`**. ✅ Kekhawatiran bahwa *"ekspor berbeda dari layar"* **tidak terbukti**,
> dan itu **menguatkan** seluruh pembacaan §B4.

⛔ **Tetapi pertanyaan aslinya TETAP TERBUKA.** Yang ditanyakan adalah **apakah `TRAVEL2`
DIJALANKAN**, bukan apakah tandanya **tercatat** — dan menarik dugaan `TRAVEL1` **tidak menjawabnya**.
Yang masih diperlukan: **menelusuri satu kasus Travel dan melihat langkah mana yang benar-benar
berjalan.** `[terbuka]`

⛔ **Dan pertanyaan aslinya belum terjawab.** Yang ditanyakan adalah **apakah `TRAVEL2`
DIJALANKAN**, bukan apakah tandanya **tercatat**. Ketiga tanda itu memang ada — itu sudah saya
baca sendiri. Yang belum: **menelusuri satu kasus Travel dan melihat langkah mana yang benar-benar
berjalan.** `[terbuka]`

#### ⛔ DICABUT — `[keputusan work owner]` 2026-09-19 · **`GetAllData_Act` DITIRU APA ADANYA**

> **"GetAllData_Act untuk ini ikuti apa adanya, jangan jadikan permasalahan."**

⭐ **Lompatan yang tertukar DITIRU APA ADANYA. Logikanya TIDAK diperbaiki di sistem baru.**
⛔ **Ia bukan penyimpangan sadar**, dan **bukan pemblokir apa pun.**

⚠️ **Yang tetap berdiri sebagai temuan:** pasangan gerbang→tanda memang **tertukar** di dua blok
bersarang — itu dibaca dua kali dengan dua cara, dan tidak dicabut. Yang dicabut adalah
**perlakuannya**, bukan **temuannya**.

> ### ⛔ Yang DICABUT, dikutip utuh dan tidak dihapus
>
> Keputusan sebelumnya pada hari yang sama berbunyi:
>
> > **"sepertinya itu kesalahan developer sebelumnya, perbaiki logic nya di go."**
> >
> > *⭐ **Lompatan yang tertukar dinyatakan SALAH TULIS, dan logikanya DIPERBAIKI di sistem baru —
> > bukan ditiru.*** ⚠️ **Ini penyimpangan sadar**, dan sejauh ronde 2 ia **yang pertama dan
> > satu-satunya** yang tercatat untuk Claim Fac In … **Pasangan yang diperbaiki** — tiap lini
> > menuju tandanya sendiri: `IsMBU` → `MBU2`/`MBU3` · `IsTravel` → `TRAVEL2`/`TRAVEL3` ·
> > `IsPA` → `PA2`/`PA3`."*
>
> **Sebab pencabutannya, dan ia bersandar pada bukti yang baru terbaca:** lihat blok di bawah —
> **`GetAllData_Act` tidak menulis apa pun ke kasus**, sehingga cacat itu **hanya menyentuh tampilan
> satu popup**, bukan data yang tersimpan.

#### ⭐ Buktinya: **NOL penugasan ke `pyWorkPage`** — activity ini tidak menyimpan apa pun

`[terverifikasi]` Sensus **seluruh 61 penugasan properti** di ke-50 langkah, dikelompokkan menurut
halaman sasarannya:

| Halaman sasaran | Jumlah |
| --- | --- |
| `Local` *(variabel kerja)* | **33** |
| **relatif** — seluruhnya jatuh ke `OutputData.pxResults` | **14** |
| `MataUang` | **9** |
| `TempOutstanding` | **2** *(langkah 19)* |
| `InputData` · `SearchNilaiKursInput` · `DataView` | **1** masing-masing |
| ⭐ **`pyWorkPage`** | ⭐ **NOL** |

⚠️ **Ke-14 penugasan relatif itu hampir menipu:** dua di antaranya bersarang di bawah
**16: `pyWorkPage.ClaimData.ObjectList`**. Tetapi ditelusuri ke atas, **halaman terdekat pada
langkah yang menugaskan selalu `OutputData.pxResults`** — kesepuluh-empatnya diperiksa satu per
satu, bukan disimpulkan dari dua contoh.

⭐ **Kesimpulan:** cacat lompatan itu **hanya menyentuh isi popup `Outstanding`**, dan
**tidak pernah mengotori data kasus yang tersimpan**. Itulah yang membuat *"jangan jadikan
permasalahan"* berdiri di atas bukti, bukan sekadar pilihan.

⛔ **Yang BELUM dibuktikan:** apakah `TempOutstanding` disimpan oleh rule **lain**. Namanya
*"Temp"* dan harness-nya hanya-baca — itu petunjuk, **bukan bukti**. `[terbuka]`

✅ **Dan pemeriksaan layar Pega tidak lagi memblokir apa pun.** Ia tetap berguna untuk memastikan
cara membaca lompatan di modul lain, tetapi **untuk Claim Fac In ia sudah tidak menahan keputusan
mana pun.**

<!-- teks lama di bawah ini SUDAH TIDAK BERLAKU; dibiarkan sebagai jejak -->

⚠️ **Ini penyimpangan sadar**, dan sejauh ronde 2 ia **yang pertama dan satu-satunya** yang
tercatat untuk Claim Fac In. ⛔ Ia **wajib ikut terbawa** ke bab *titik yang sengaja diubah* ketika
spec modul ini kelak ditulis — bersama alasannya, supaya orang berikutnya **tidak menganggapnya
cacat pemindahan**.

**Pasangan yang diperbaiki** — tiap lini menuju tandanya sendiri:

| Gerbang | Di Pega melompat ke | Di sistem baru |
| --- | --- | --- |
| `IsMBU` | ~~`TRAVEL2` / `TRAVEL3`~~ | **`MBU2` / `MBU3`** |
| `IsTravel` | ~~`MBU2` / `MBU3`~~ | **`TRAVEL2` / `TRAVEL3`** |
| `IsPA` | `PA2` / `PA3` | **`PA2` / `PA3`** — sudah benar, tidak berubah |

⚠️ **Bentuk perbaikan yang benar TIDAK ditetapkan di sini.** Tabel di atas menyebut **pasangan
yang benar**, bukan cara mewujudkannya. Yang mana pun dipilih — membetulkan sasaran lompatannya,
atau membuang lompatannya sama sekali karena **tiap langkah pendarat sudah menguji gerbangnya
sendiri** — hasilnya sama. ⛔ Itu keputusan rancangan, dan **bukan milik berkas grilling**.

⚠️ **Dasar "pasangan yang benar" adalah langkah 11**, yang pasangannya memang benar, **bukan
dokumen mana pun**. Tidak ada satu berkas pun di korpus yang menyatakan niat aslinya. ⛔ Bila kelak
ternyata langkah 11 yang justru keliru, perbaikan ini ikut keliru.

⚠️ **Akibat yang TIDAK ikut diputuskan — data lama.** Bila `TRAVEL2` dan `TRAVEL3` memang tidak
pernah berjalan, maka baris **Travel** yang sudah tersimpan selama ini **kekurangan isi yang
seharusnya ada**. Memperbaiki logika di sistem baru **tidak** memperbaiki data lama, dan sejak
perbaikan itu berjalan **baris Travel baru akan berbeda dari baris Travel lama**. ⛔ Apa yang harus
dilakukan terhadap data lama **belum diputuskan** — lihat §H pertanyaan baru. `[terbuka]`

### B5 — Sebaran arah gerbang di rule ini

```
lanjut-when 174 · lewati-langkah 30 · lompat-ke-langkah 9 · (kosong) 3
kode 4 (keluar iterasi) 0 · kode 5 (lewati when) 0 · kode 6 (keluar activity) 0
```

⭐ **Nol kode 5** — jadi **seluruh rantai gerbang di rule ini bersifat DAN**, tidak ada yang berubah
menjadi ATAU *(aturan C3)*.

### B6 — Halaman yang dibaca dan ditulis

`[terverifikasi]` `pyStepsObjectName` yang terisi: `OutputData` *(langkah 3–9)* ·
`OutputData.pxResults` *(12 · 13 · 14 · 16.1.1.1 · 16.1.2.1 · 18.2 · 18.3)* ·
`pyWorkPage.ClaimData.ObjectList` *(16)* · `.ObjectItemList` *(16.1)* · `.EstimationList` *(16.1.1)* ·
`.Adjustment` *(16.1.2)* · `MataUang.pxResults` *(17 · 18)* · `SearchNilaiKursOutput` *(17.2)*.

⚠️ **`.CARI7` · `.CARI12` · `.CARI13` · `.CARI15` muncul sebagai gerbang** *(14.1 · 11.1 · 12.1 ·
13.1 · 18.2.1 · 18.3.1)*. Penampung bernama `CARI<n>` itu **pola yang sama** dengan yang ditemukan
di Master Contract Retro Life, tempat parameter procedure dikirim lewat penampung bernomor alih-alih
nama kolom. ⛔ Apa isi tiap `CARI<n>` di modul ini **tidak terbaca dari rule ini**. `[terbuka]`

---

## §C — `CreateKMTNo_Act` — pembuat kasus komite

⭐ **Dua jangkar, BEDA lagi** *(aturan I3/P4)*: **jangkar A 42** · **jangkar B 35** · **pengurai XML
42**. Selisih **7** sama dengan jumlah langkah tanpa metode. **Angka yang dipakai: 42.**

### C1 — Bentuknya

`[terverifikasi]` **42 langkah · NOL ber-remark · NOL lompatan.**
Flag prakondisi: **true 9 · false 1 · (kosong) 32**.
Arah: `lanjut-when` 154 · `lewati-langkah` 12 · `keluar-activity` **1** · (kosong) 7.

⭐ **Satu-satunya penjaga di seluruh rule ini** adalah **langkah 10**:
`Call pxAddChildWork`, bergerbang `@hasMessages(pyWorkPage)` → **keluar-activity**. Artinya: **bila
ada pesan kesalahan pada kasus, pembuatan kasus komite dibatalkan.** Di luar itu, **nol penjaga**.

⚠️ **Satu gerbang MATI** *(flag `false`)*: **langkah 6.8.1**, `Local.SizeADj>0` — pemeriksaan
*"ada tidaknya baris adjustment"* **tersimpan tetapi tidak berlaku**.

### C2 — ⭐ Penunjuk posisional: **LIMA, bukan empat**

`[terverifikasi]` Seluruhnya di **langkah 6.9**, dan tidak satu pun bergerbang:

| Yang diisi | Diisi dari |
| --- | --- |
| `childPageKomite.Adjustment.IDObject` | `.IndexObject` |
| `childPageKomite.Adjustment.CoverageSubscript` | `.IndexCoverage` |
| `childPageKomite.IndexObjectItem` | `Local.IndexObjectItem` |
| `childPageKomite.IndexObject` | `.ObjectIndex` |
| ⭐ **`childPageKomite.IndexAdjustment`** | **`local.IndexAdjust`** |

⭐ **RALAT terhadap ronde ulang §B5:** ia menulis **empat** penunjuk — `IDObject` ·
`CoverageSubscript` · `IndexObjectItem` · `IndexObject`. Yang **kelima**, `IndexAdjustment`,
**terlewat**. Kalimat lamanya dikutip di §H3.

⚠️ **Dua kejanggalan pada kelimanya:**

1. **`IndexObject` diisi dari `.ObjectIndex`, sedangkan `Adjustment.IDObject` diisi dari
   `.IndexObject`.** Dua nama sumber yang berbeda untuk dua penunjuk yang namanya nyaris sama.
   ⛔ Mana yang mana **tidak terbaca dari rule ini**. `[terbuka]`
2. **`local.IndexAdjust` ditulis huruf kecil**, sedangkan empat lainnya `Local.` / `.`.
   ⚠️ Pega **tidak peka huruf besar-kecil** untuk ini *(aturan I1)*, jadi kemungkinan besar tidak
   berakibat — tetapi ia **petunjuk bahwa baris ini ditulis belakangan**, terpisah dari empat lainnya.

⛔ **Nol pemeriksaan** bila yang ditunjuk hilang atau bergeser — sama seperti Claim Prop, tetapi
**dengan tiga tingkat lebih dalam** *(objek → item objek → adjustment)*, sehingga peluang bergesernya
lebih besar. ⛔ **Saya tidak menyimpulkan apa pun untuk aplikasi Go.**

### C3 — Daftar penyetuju

`[terverifikasi]` **Langkah 8** `Call SetListKomite_act`, lalu **langkah 9**, **tanpa gerbang**:

```
childPageKomite.KomiteList  := .ComiteeClaim
childPageKomite.KomiteLoop  := @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)
childPageKomite.KomiteCount := 1
```

⚠️ **Nol penjaga bila daftarnya kosong.** Bila `.ComiteeClaim` kosong, `KomiteLoop` menjadi **0**
sementara `KomiteCount` tetap **1** — **tangga nol tingkat, senyap**. Bentuk yang **sama persis**
dengan yang sudah tercatat di modul Komite Claim Life. ⛔ Di sini **belum pernah dibawa ke work
owner**. `[terbuka]`

### C4 — Yang ditulis ke kasus klaim induk

| Langkah | Gerbang | Yang ditulis |
| --- | --- | --- |
| **3** | — | `pyWorkPage.ClaimData.DaftarObjek := .IndexObject` |
| ⭐ **3** | — | **`pyWorkPage.ClaimData.FlagOnGoingCommitte := "Send Commite"`** |
| **4** | — | empat penulisan tanggal ke `pyWorkPage.…` |
| **6.8.1.1** | `true` | `.IsKomite := 1` |
| **12** | — | `.IsKomite := 1` *(kedua kalinya)* · dan sebuah kolom `ObjectList(...)` di-nol-kan |
| **13.1** | — | `.KomiteNo := childPageKomite.pyID` |
| **13.1.1** | `true` | `.KomiteNo := @substring(pyWorkPage.pxCoveredInsKeys(<last>),19)` |

⭐ **Tiga hal yang menonjol:**

1. **`FlagOnGoingCommitte` diisi teks `"Send Commite"`** — bukan `1`, bukan benar/salah. ⚠️ Kolom
   bernama sama ada di modul `Claim Prop`. ⛔ **Saya tidak menyimpulkan keduanya sama** — nama yang
   sama bukan bukti rule atau kolom yang sama. Bentuk nilainya di sini: **teks, dan teksnya
   mengandung salah ketik** *(`Commite`, bukan `Committee`)*.
2. **`.IsKomite := 1` ditulis DUA KALI** — langkah 6.8.1.1 *(bergerbang)* dan langkah 12 *(tanpa
   gerbang)*. Yang kedua membuat yang pertama **tidak berpengaruh pada hasil akhir**.
3. **`.KomiteNo` juga ditulis dua kali** — 13.1 dari `childPageKomite.pyID`, lalu 13.1.1 dari
   **potongan teks berawal huruf ke-19** dari kunci instance. ⚠️ Angka **19** ditanam langsung di
   rumus. ⛔ Alasannya tidak tertulis. `[terbuka]`

### C5 — ⚠️ Urutan efek keluar terhadap penyimpanan

`[terverifikasi]` **`Call SendEmailKlaim` di langkah 14**, **`Obj-Save` di langkah 19** — jadi
**email keluar SEBELUM kasus tersimpan**. Ditambah `Call InsertProgressClaim` **(16)** dan
`Apply-DataTransform` *"set Note history status claim"* **(17)**, keduanya juga sebelum simpan.

⛔ **Ini fakta tentang Pega, bukan usulan rancangan.**

---

## §D — Kandidat `pyMemo` diuji dengan membuka rule-nya

### D1 — ⚠️ Jendela saya BERBEDA dari ronde 1, dan keduanya dilaporkan

`[terverifikasi]` **459 catatan `pyMemo` di 449 berkas**, **457 pasangan (berkas, catatan) berbeda**
— **ketiga angka COCOK** dengan ronde 1 §F1.

Tetapi **penggolongannya tidak cocok**, karena ronde 1 **tidak mencantumkan daftar kata kuncinya**,
sehingga golongan (c)-nya **tidak dapat direproduksi**. Jendela saya, disebut penuh:

```
(d) kandidat LARANGAN   kata: hapus · jangan · jgn · ganti · buang
(c) kandidat MENGUBAH   kata: skip · lewati · sementara · belum · harusnya · seharusnya · fix ·
                              error · bug · salah · off · disable · nonaktif · dummy · test · coba ·
                              tidak dipakai · ga dipake · gak dipake · tdk dipakai · nanti · todo ·
                              sblm · blm · pending · revisi · ubah · tambah
```

| Golongan | Ronde 1 | Ronde 2 *(jendela saya)* |
| --- | --- | --- |
| (d) larangan | **12** | **32** |
| (c) mengubah pembacaan | **118** | **91** |
| (a)+(b) sisanya | **327** | **336** |

⛔ **Keduanya dilaporkan, tidak ada yang ditimpa.** Angka yang berbeda dua cara = **belum punya
data** *(aturan I3/P4)*. Yang **pasti** hanyalah **459 · 449 · 457**.

### D2 — **Diuji 12 dari 91** *(atau dari 118 menurut ronde 1)* — ⭐ **5 pindah golongan**

⛔ **Belum selesai. 79 kandidat lagi belum dibuka.**

| Rule | Catatan | Hasil pembukaan rule |
| --- | --- | --- |
| ⭐ `Activity/GetUrlGoogleStorage_Act.xml` | *"FIX ERROR HANDLING"* | ✅ **SUDAH DIPATUHI** — `pyOnException` terisi di langkah **6.1** dan **6.5** → **pindah (c) → (b)** |
| ⭐ `Activity/InsertGoogleStorage_Act.xml` | *"FIX ERROR HANDLING"* | ✅ **SUDAH DIPATUHI** — `pyOnException` terisi di langkah **7** dan **11** → **pindah (c) → (b)** |
| `Activity/SetInputParam_Act.xml` | *"tambah save"* | ✅ **SUDAH DIPATUHI** — `Obj-Save` di langkah **5**, langkah terakhir → **pindah (c) → (b)** |
| `Activity/ViewKomite_act.xml` | *"tambah save"* | ✅ **SUDAH DIPATUHI** — `Obj-Save` di langkah **4**, langkah terakhir → **pindah (c) → (b)** |
| `Activity/CheckEstimateValue.xml` | *"test"* | ✅ **catatannya menyesatkan** — **60 langkah**, **7** `Property-Set-Messages`, pesan di 8 titik. Ini rule validasi sungguhan, bukan uji coba → **pindah (c) → (b)** |
| ⚠️ `Activity/SaveAcceptation.xml` | *"ubah bagian konversi"* | ⚠️ **MASIH ADA, dan gerbangnya DIMATIKAN** — sehingga berjalan tanpa saringan; lihat D3 dan ralatnya |
| ⚠️ `Activity/SetDisable_ACT.xml` | *"disable button send komite"* | ⚠️ **belum dapat diputuskan** — 12 langkah, **3 ber-remark**, ada `Page-Set-Messages` di **7.1** |
| ⚠️ `Activity/CountSpreadingClaim_ACT.xml` | *"ubah di step 5"* | ⚠️ **belum dapat diputuskan** — 37 langkah, langkah **5** ada, tetapi **perubahan yang dimaksud tidak tertulis** |
| ⚠️ `Activity/SetRejectClaim_pre.xml` | *"tambah unutk set date reject"* | ⚠️ **belum dapat diputuskan** — hanya **2 langkah** |
| ⚠️ `Activity/DeleteTest.xml` | *"fix 2.2"* | ⚠️ **rule bernama "Test" dan HIDUP** — 8 langkah, nol remark, memanggil `ProtectionDate_Act` |
| `Activity/CheckTotalSpreadingPct_Act.xml` | *"test"* | ⚠️ 8 langkah, 1 `Property-Set-Messages` — **bukan uji coba**, tetapi golongan barunya belum pasti |
| `Activity/SetTempLocation_Act.xml` | *"test"* | ⚠️ 2 langkah, **1 ber-remark** — terlalu kecil untuk diputuskan |

⭐ **Arah lesetnya sama dengan ronde ulang:** penggolongan dari kata kunci **melebih-lebihkan**
golongan yang mengkhawatirkan. **5 dari 12 ternyata sudah beres.**

⚠️ **Dua catatan ini berbunyi sama persis dengan yang ada di modul Komite Claim Prop** — *"FIX ERROR
HANDLING"*. ⭐ **Tetapi hasilnya BERLAWANAN:** di sana penanganannya **tidak berjejak** dan tercatat
sebagai butir terbuka; **di sini `pyOnException` benar-benar terpasang di keempat titik.** ⛔ **Nama
catatan yang sama bukan bukti keadaan yang sama** — ini buktinya.

### D3 — ⭐⭐ `SaveAcceptation` — empat langkah mati, dan arasapas **berjalan TANPA SARINGAN**

> ⛔ **Judul bagian ini diralat 2026-09-19.** Bentuk lamanya, dikutip: *"### D3 — ⭐⭐
> `SaveAcceptation` — empat langkah mati, dan arasapas **DIMATIKAN**"*. Sebabnya di ralat bawah.

`[terverifikasi]` 20 langkah, **4 ber-remark**:

| Langkah | Keadaan | Isinya |
| --- | --- | --- |
| **9** · **9.1** | ⛔ **REMARK** | `Call SaveAccept_ACT` — *"activity save to db"* |
| **9.2** | **HIDUP** *(flag true)* | ⭐ **`Call HitServiceToKasir_Act`** — kiriman uang ke Kasir |
| **10** | ⛔ **REMARK** | `RDB-List` — *"Untuk Save OS_Akseptasi_klaim"* |
| **11** | ⛔ **REMARK** | `Call PrintPDFAcceptanceNote` — *"per 1 aksep"* |
| **12** | **HIDUP** | `Call PrintPDFAccep_MultiAksep` |
| **14** | **HIDUP** | **`Obj-Save`** |
| **15** | **HIDUP** | `Call InsertJsonClaimNonMBU_act` |
| ⭐ **16** | ⚠️ **BERJALAN — gerbangnya yang mati** *(flag `false`)* | **`Call KonversiKlaim_Act`** — *"Untuk Hit Service Arasapas"* |

⭐ **Tiga temuan sekaligus:**

1. **Penyimpanan berpindah rumah** — `Call SaveAccept_ACT` *(9.1)* dan penyimpanan
   `OS_Akseptasi_klaim` *(10)* **dimatikan**, diganti `Obj-Save` polos di **14**.
2. **Kiriman ke Kasir *(9.2)* dan PDF *(12)* berjalan SEBELUM `Obj-Save` *(14)*; penulisan JSON
   *(15)* berjalan SESUDAHNYA.** ⚠️ Jadi urutannya **campuran** — berbeda dari Claim Prop dan
   Komite Claim Prop, tempat **kedelapan** efek keluar berjalan sebelum penyimpanan.
3. ⭐⭐ **Panggilan arasapas BERJALAN — dan justru TANPA SARINGAN.** Lihat ralat di bawah.

> ### ⛔ RALAT 2026-09-19 — **kesalahan saya, dikoreksi work owner**
>
> **Kalimat lama, dikutip apa adanya:**
>
> > *"⭐ **16** ⛔ **MATI** *(flag `false`)* — **`Call KonversiKlaim_Act`** — *"Untuk Hit Service
> > Arasapas"*"* · dan · *"⭐ **Panggilan arasapas MATI** — gerbangnya `false`. ⛔ Di modul Claim
> > Prop panggilan arasapas **hidup**."*
>
> ⛔ **Itu SALAH, dan salahnya pada penafsiran — bukan pada pembacaan medan.**
> `pyStepsPreCondition = false` berarti **gerbangnya tersimpan tetapi TIDAK DITEGAKKAN**. Ia
> **tidak** mematikan langkahnya. Yang mematikan langkah adalah `pyStepsBlockName = "//"`
> *(remark)* — dan langkah 16 **tidak** ber-remark.
>
> ⚠️ **Akibatnya BERKEBALIKAN dari yang saya tulis:** langkah itu **berjalan lebih sering**, bukan
> tidak pernah. Kedua baris gerbang yang tersimpan berbunyi:
>
> ```
> when pyWorkPage.AcceptStatus == "1"                  T=lanjut-when  F=lewati-langkah
> when pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop  T=lanjut-when  F=lewati-langkah
> ```
>
> Seandainya gerbangnya **hidup**, arasapas hanya dipanggil bila **status akseptasi = 1** DAN
> **tangga komite sudah tuntas**. Karena gerbangnya **mati**, kedua syarat itu **tidak berlaku** —
> panggilan berjalan **apa pun status akseptasinya, dan walau tangga komite belum selesai**.
>
> ✅ **Work owner memeriksanya langsung di Fac In dan menyatakan "TIDAK MATI"** —
> `[keputusan work owner]` 2026-09-19. Pemeriksaan itu **benar**, dan berkas ini mengikuti.

#### ⭐ Dan ternyata pemanggilnya EMPAT, nol di antaranya bergerbang hidup

`[terverifikasi]` Penyisiran ke-**482** berkas `.xml` untuk pemanggil `KonversiKlaim_Act`:

| Pemanggil | Langkah | Flag prakondisi | Artinya |
| --- | --- | --- | --- |
| `Activity/CLaimFaceSheet_Act.xml` | **47** | **(kosong)** | ⚠️ **nol gerbang sama sekali** — selalu jalan |
| `Activity/CloseClaim.xml` | **8.4** | **(kosong)** | ⚠️ **nol gerbang sama sekali** — selalu jalan |
| `Activity/ChooseDla_Act.xml` | **6** | **`false`** | gerbang **3 baris** tersimpan, **tidak ditegakkan** |
| `Activity/SaveAcceptation.xml` | **16** | **`false`** | gerbang **2 baris** tersimpan, **tidak ditegakkan** |

⭐ **Empat pemanggil, dan tidak satu pun benar-benar bergerbang.** Dua **memang tidak punya gerbang**,
dua **punya tetapi dimatikan**. ⛔ **Saya tidak menyimpulkan apa pun untuk aplikasi Go**, dan
⛔ **tidak menyatakan ini cacat** — yang saya nyatakan: **layanan arasapas dipanggil dari empat
tempat, tanpa satu pun saringan yang berlaku.** `[terbuka]` apakah itu memang dikehendaki.

### D4 — Yang menyebut nama orang, akun, atau alamat

`[terverifikasi]` **ADA**, dan ⛔ **nilainya tidak saya salin**:

- `ConnectREST/SendAcceptationToKasir.xml` — catatan **"tambah auth"**.
  ⭐ Inilah **satu dari 8 ConnectREST yang berautentikasi** *(ronde 1 §E3)*. Bentuk yang tercatat:
  **profil autentikasi yang tertanam di dalam rule**. ⛔ **Nilainya tidak dibaca, tidak dicetak,
  tidak disalin.**
  ⚠️ **Karena kredensial itu beredar di dalam berkas ekspor, ia SEBAIKNYA DIGANTI sesudah migrasi.**
  Itu urusan tim pemilik layanan Kasir — dicatat di sini semata supaya tidak terlewat, ⛔ **tidak
  ditindaklanjuti dari berkas ini.**
- Ruleset **`ADESAMUEL@`** — nama akun orang di dalam identitas rule. Lihat §G1.

---

## §E — ⭐ `pyXMLSignature` TERPECAHKAN dengan pola ketiga

### E1 — Polanya: **XML di dalam XML**

Ronde 1 §C8 dan ronde ulang §E1 gagal dengan dua pola. **Pola ketiga berhasil** `[terverifikasi]`:

> **Isi medan `pyXMLSignature` adalah dokumen XML tersendiri yang lengkap**, berawalan
> `<?xml version="1.0"?>`, berakar `<methodSignature>`, berisi `<pyStepsCallParams>`, dan **setiap
> anaknya adalah satu parameter** — **nama parameter menjadi nama elemen**, sedangkan sifatnya
> menjadi **atribut**: `TYPE` · `INOUT` · `REQUIRED` · `DESCRIPTION` · `SIZE` · `DEFAULTVALUE` ·
> `SMARTPROMPT*`.

**Kedua pola sebelumnya gagal** karena memperlakukan isinya sebagai **teks**, bukan sebagai
**dokumen XML yang harus diurai ulang**.

### E2 — Hasilnya

```
Activity                                179
pyXMLSignature terisi                   178     kosong 1  (InsertLogServiceClaim.xml)
berhasil diurai sebagai XML             178     gagal 0
TOTAL parameter                         169
activity yang PUNYA parameter            72
activity yang parameternya NOL          106
```

| Sifat | Sebaran |
| --- | --- |
| `INOUT` | **IN 162** · **OUT 7** |
| `REQUIRED` | **`0` 154** · **`-1` 15** |
| `TYPE` | STRING **147** · INTEGER **13** · Date **4** · DateTime **2** · TrueFalse **1** · Decimal **1** · Double **1** |

**Activity dengan parameter terbanyak:** `InsertDocument_Act` **9** *(`IDPEGA` · `NAMAFILE` · `MIME`
· `KATEGORI_1` · `KATEGORI_2` · `NOAKSEP` · `NOPREKAS` · `PAYMENTDATE` · `BASE64`)* ·
`GetCoverageAneka_Act` **7** · `InsertGoogleStorage_Act` **6**.

⚠️ **Dua hal yang TIDAK saya simpulkan:**

1. **Arti `REQUIRED = -1` versus `0`.** Dugaan yang wajar adalah *"-1 = wajib, 0 = tidak"*, tetapi
   **itu tidak tertulis di mana pun** — persis jenis nilai tak berdokumen yang sudah menjadi
   **Pertanyaan 2** ronde 1. ⛔ **Tidak ditebak.** `[terbuka]`
2. ⚠️ **Hanya 2 dari 169 parameter bertipe angka pecahan** *(1 `Decimal`, 1 `Double`)*, sedangkan
   **147 bertipe `STRING`**. ⛔ Apakah ada parameter uang yang lewat sebagai teks **belum saya
   periksa satu per satu** — itu pekerjaan ronde berikutnya, dan ia menyentuh **ADR-0003**.

---

## §F — Dua pernyataan "nol" yang diperiksa ulang

### F1 — ⭐ Gerbang aksi tombol: `pyActionConditions` memang kosong, **tetapi gerbang tombol ADA**

**Jendela saya, disebut penuh:** seluruh **106 berkas layar** — `Section` **57** · `Harness` **16** ·
`FlowAction` **33** — disisir untuk **SELURUH** nama medan yang mengandung
`cond|action|disab|visib|require|readonly|enable|click|button|behav`, **tidak peka huruf besar-kecil**
*(aturan I1)*. **Hasilnya 121 nama medan berbeda** — bukan 3, bukan 11.

`[terverifikasi]` **`pyActionConditions`: hadir 1 438 kali, terisi NOL.** ✅ **Pernyataan ronde 1 dan
ronde ulang BERTAHAN**, dan kini jendelanya **jauh lebih lebar**.

⭐ **Tetapi menyebutnya "gerbang aksi tombol tidak ketemu" TERLALU LUAS.** Sensus 121 medan menemukan
gerbang tombol yang **nyata**:

| Medan | Terisi | Nilainya | Nyata atau baku? |
| --- | --- | --- | --- |
| ⭐ **`pyDisableSubmit`** | **727** di 49 berkas | `false` 675 · ⭐ **`true` 52** | **52 NYATA** |
| ⭐ **`pyLocalAction`** | **28** di 10 berkas | `SetPassWordSP` 6 · ⭐ **`PreventRejectClaim` 4** · `ShowRetro` 4 · `MessageBeforeDeleteTreatyGroup` 4 | **NYATA — aksi bernama** |
| `pyShowFAButtons` | 33 | `true` 32 · `false` 1 | 1 nyata |
| `pyIsClientDisableWhen` | 203 di 94 berkas | **seluruhnya `false`** | ⛔ baku, bukan gerbang |
| `pyDisqualifyAction` | 33 | seluruhnya `true` | ⛔ baku |
| `pySaveButtonVisibility` | 33 | seluruhnya `Always` | ⛔ baku |
| `pyActionsType` | 23 | seluruhnya `standard` | ⛔ baku |

⭐ **Kesimpulan yang benar:** *"`pyActionConditions` kosong di seluruh 106 berkas layar"* — **SAH**.
*"Gerbang aksi tombol tidak ada"* — ⛔ **TIDAK SAH**: ada **52 titik `pyDisableSubmit = true`** dan
**4 aksi lokal bernama**, satu di antaranya berbunyi **`PreventRejectClaim`**.

### F2 — ⚠️ Alias SQL: angka saya **BERBEDA** dari ronde ulang

**Jendela saya, disebut penuh:** ke-**63** berkas `RDBList`, medan **`pyBrowseSQL`**, pola
**`<sumber> as "<alias>"`** dengan alias **berkutip ganda**; pembanding memakai **potongan terakhir
sesudah titik**, tidak peka huruf besar-kecil.

| | Ronde ulang §E2 | Ronde 2 *(jendela saya)* |
| --- | --- | --- |
| pasangan alias | **50** | **42** |
| **TIDAK COCOK** | **38** | **39** |
| cocok | 12 | **3** |

⛔ **Kedua hasil dilaporkan; tidak ada yang ditimpa.** Angka yang berbeda dua cara = **belum punya
data** *(aturan I3/P4)*. Sebabnya hampir pasti **alias tanpa kutip** yang jendela saya buang, dan
**3 positif palsu** yang ronde ulang §H5 butir 3 catat tetapi tidak bersihkan. ⛔ **Saya belum dapat
membersihkan ketiganya**, karena ronde ulang **tidak menyebut yang mana**. `[terbuka]`

⚠️ **Yang PASTI, karena dibaca langsung:** di `CariHistoryClaim_SQL.xml` modul ini, **5 alias dan
5-lima-nya berbohong** — pola yang **sama persis** dengan yang tercatat di Claim Prop:

```
a.DATA_JSON.DateOfLoss   -> "START_DATE"      tanggal kejadian  dinamai  tanggal mulai
IDPEGA                   -> "BRANCH_CODE"     id Pega           dinamai  kode cabang
a.DATA_JSON.ClaimNo      -> "BRANCH_NAME"     nomor klaim       dinamai  nama cabang
a.DATA_JSON.CauseOfLoss  -> "BUSINESS_CODE"   penyebab kerugian dinamai  kode bisnis
NOPOLIS                  -> "POLICY_NO"       ✅ satu-satunya yang jujur
```

---

## §G — Tiga hal yang belum pernah tercatat di modul mana pun

### G1 — ⭐ Ruleset `ADESAMUEL@` — dan **flow satu-satunya modul ini ada di dalamnya**

`[terverifikasi]` Sensus ruleset atas **482 berkas `.xml`**:

```
GCNMFW      352      GISFW      115      ⭐ ADESAMUEL@   5
GCNMFWInt     5      SFAGIS       5
```

**Kelima rule di `ADESAMUEL@`, seluruhnya versi `01-01-01`:**

| Berkas | `pxInsName` | Kelas |
| --- | --- | --- |
| ⭐ **`Flow\Register_Flow.xml`** | `ASM-FW-GCNMFW-WORK-PNC!REGISTER_FLOW` | `…Work-PNC` |
| `Section\ClaimSurvey.xml` | `…WORK-PNC!CLAIMSURVEY` | `…Work-PNC` |
| `Section\InputInwardFacultativeDtl.xml` | `ASM-FW-GCNMFW-WORK!INPUTINWARDFACULTATIVEDTL` | `…Work` |
| `Section\ViewHistoryClaim.xml` | `…WORK-PNC!VIEWHISTORYCLAIM` | `…Work-PNC` |
| `Section\ViewPolis.xml` | `…WORK-PNC!VIEWPOLIS` | `…Work-PNC` |

⭐ **Modul ini punya TEPAT SATU berkas `Flow` — dan ia ada di ruleset bernama akun orang.**
`Register_Flow.xml` adalah **daur hidup kasus seluruh modul**, dan ia **tidak berada di `GCNMFW`
maupun `GISFW`**.

**Menyentuh uang atau efek keluar?** `[terverifikasi]` **Tidak langsung** — kelimanya **Flow** dan
**Section**, nol `Activity`, nol `RDBList`, nol `ConnectREST`. ⚠️ **Tetapi `Flow` menentukan tahap
mana yang berjalan**, jadi pengaruhnya tidak kecil.

⚠️ **Dicatat sebagai temuan tentang BERKAS, bukan tentang orang.** ⛔ Saya tidak menilai siapa pun
dan **tidak mengusulkan tindakan kepegawaian apa pun**. Yang perlu diketahui hanyalah: **ruleset
bernama akun perorangan biasanya ruleset kerja pribadi, dan biasanya tidak dimaksudkan naik ke
produksi.** Apakah di sini ia memang dipakai di produksi **tidak terbaca dari korpus**.

> ### ✅ TERJAWAB — `[keputusan work owner]` 2026-09-19
>
> **"ikuti aja begitu, itu yang akan dipake."**
>
> ⭐ **`Flow\Register_Flow.xml` di ruleset `ADESAMUEL@` versi `01-01-01` ADALAH daur hidup kasus
> yang berlaku**, dan **itulah yang dipindahkan**. ⛔ **Tidak dicari penggantinya di ruleset resmi**,
> dan ⛔ **tidak diperlakukan sebagai ekspor yang kurang lengkap**.
>
> ⚠️ **Yang tetap harus dicatat, karena ia fakta dan bukan keberatan:** daur hidup seluruh modul
> ini tinggal di ruleset **bernama akun perorangan**, terpisah dari `GCNMFW` *(352 rule)* dan
> `GISFW` *(115 rule)* tempat semua rule lain tinggal. Siapa pun yang kelak mencari
> `Register_Flow` di ruleset resmi **tidak akan menemukannya** — dan itu bukan tanda ekspor rusak.
>
> ⛔ **Keempat Section lain di ruleset yang sama belum diputuskan** — jawaban ini menyebut yang
> **akan dipakai**, dan yang ditanyakan adalah **flow**-nya. Isi keempat Section itu juga belum
> dibaca *(§H4 butir 5)*. `[terbuka]`

### G2 — ⭐ `BrowseHistoryClaim` — kembaran yang aliasnya juga berbohong

`[terverifikasi]` **Ada, satu berkas**, `RDBList\BrowseHistoryClaim.xml`:

```
pxInsName    ASM-FW-GCNMFW-INT-V_POLIS!GCNM!BROWSEHISTORYCLAIM
ruleset      GCNMFWInt        versi 01-01-06        jenis RULE-CONNECT-SQL
```

**SQL-nya, dikutip apa adanya:**

```
select a.DATA_JSON.DateOfLoss as "BRANCH_CODE",
       NOPOLIS AS "POLICY_NO",
       a.DATA_JSON.ClaimNo AS "BUSINESS_CODE"
from   json_klaim a
where  nopolis = {TempPolis.PolicyNo}
order by TGL_INPUT ASC
```

⭐ **Dua dari tiga aliasnya berbohong:**

| Sumber | Alias | Kenyataannya |
| --- | --- | --- |
| `a.DATA_JSON.DateOfLoss` | `"BRANCH_CODE"` | ❌ **tanggal kejadian** dinamai **kode cabang** |
| `NOPOLIS` | `"POLICY_NO"` | ✅ jujur |
| `a.DATA_JSON.ClaimNo` | `"BUSINESS_CODE"` | ❌ **nomor klaim** dinamai **kode bisnis** |

⭐ **Dan ia BUKAN kembaran yang sama:** dibanding `CariHistoryClaim_SQL`, **pemetaan bohongnya
berbeda** — di sana `DateOfLoss → START_DATE` dan `ClaimNo → BRANCH_NAME`; di sini
`DateOfLoss → BRANCH_CODE` dan `ClaimNo → BUSINESS_CODE`. ⚠️ **Dua rule yang mirip, berbohong dengan
cara yang berbeda.** Siapa pun yang menyamakan keduanya akan salah.

⭐ **Jejak Save-As yang terbaca** `[terverifikasi]`: `pzOriginalInstanceKey` menyebut kelas
**`ASM-FW-GCNMFW-INT-V_STS_CLAIM`** *(bertanggal 2017)*, sedangkan `pxInsName` kini berkelas
**`…INT-V_POLIS`** *(bertanggal 2022)*. **Rule ini disalin dari kelas lain dan pindah kelas.**

⚠️ **Siapa pemanggilnya belum saya sisir** — itu menuntut menyisir seluruh 482 berkas untuk rujukan
nama, dan **tidak sempat di ronde ini**. `[terbuka]`

### G3 — ADR-0012 — ⛔ **TETAP tidak dapat diputuskan**

**Vonis: ⛔ tidak dapat diputuskan.** Tidak berubah dari ronde ulang §B10.

**Sebabnya, ditulis tegas:** ADR-0012 menyangkut **wewenang mengirim ke komite**. Yang ditemukan di
layar ronde ini:

- `[terverifikasi]` `Activity/SetDisable_ACT.xml` berjudul catatan **"disable button send komite"** —
  12 langkah, **3 ber-remark**, satu `Page-Set-Messages` di langkah **7.1**.
- `[terverifikasi]` **52 titik `pyDisableSubmit = true`** di berkas layar *(§F1)*.
- `[terverifikasi]` aksi lokal bernama **`PreventRejectClaim`**.

⛔ **Ketiganya menunjukkan ADA pengendalian tombol, tetapi tidak satu pun menunjukkan pengendalian
itu berdasarkan WEWENANG.** Tanpa itu, ADR-0012 tidak dapat dinyatakan dikuatkan maupun ditentang.
⛔ **Dilarang mengusulkan revisi ADR sebagai keputusan** — ADR hanya dicabut work owner.

---

## §H — Penutup

### H1 — Yang paling menahan untuk ronde 3

| # | Yang dikerjakan | Kenapa menahan |
| --- | --- | --- |
| **1** | ⭐ **79 kandidat `pyMemo` sisanya** dibuka satu per satu | ronde ini baru menguji **12**; dan 5 dari 12 **pindah golongan**, jadi sisanya pasti masih menyimpan kejutan |
| **2** | ⭐ **Sisir pemanggil `BrowseHistoryClaim` dan `CariHistoryClaim_SQL`** | dua rule berbohong dengan cara **berbeda**; siapa yang memakai yang mana menentukan data siapa yang salah baca |
| **3** | ✅ ~~Periksa apakah `Register_Flow` di ruleset `ADESAMUEL@` benar-benar yang berjalan di produksi~~ — **TERJAWAB** `[keputusan work owner]` 2026-09-19: **ia yang dipakai**. Yang **tersisa**: baca **isi** kelima rule `ADESAMUEL@`, bukan identitasnya | keempat Section-nya belum dibaca sama sekali |
| **4** | **Baca ke-52 titik `pyDisableSubmit = true`** dan keempat aksi lokal | satu-satunya jalan menuntaskan **ADR-0012** |
| **5** | **Periksa 147 parameter bertipe `STRING`** — mana yang sebenarnya uang | menyentuh **ADR-0003** *(uang non-float)* |
| **6** | **Bersihkan selisih alias 50 vs 42** dengan satu jendela yang disepakati | dua angka berbeda = belum punya data |
| **7** | **Isolasi satu lompatan kode 4** *(ronde ulang §F3)* | masih gagal sejak ronde ulang |
| ~~**8**~~ | ✅ **SELESAI 2026-09-19** — sasarannya `TempOutstanding` dan `DataView.CARI25`; **nol penugasan ke `pyWorkPage`** di seluruh 61 penugasan. Teks lamanya: ~~Telusuri sasaran penulisan `GetAllData_Act` langkah 19~~ | menentukan apakah lompatan tertukar di §B4 hanya mengotori **tampilan popup**, atau ikut mengotori **data kasus** — ⚠️ **naik peringkat** sesudah keputusan memperbaiki logika: ia yang menentukan apakah **data Travel lama** ikut cacat |
| ~~**9**~~ | ⛔ **GUGUR 2026-09-19** — tidak ada baris tersimpan yang terdampak: activity itu **tidak menulis ke kasus**, dan logikanya **ditiru apa adanya**. Teks lamanya: ~~Ukur berapa banyak baris Travel yang terdampak~~ | ~~perbaikan logika tidak memperbaiki data lama~~ |
| **10** | ⭐⭐ **Adu Claim Prop dan Komite Claim Prop dengan 15 ADR** | **A4-4** mengubahnya dari *"sebaiknya"* menjadi **"harus"**; kedua modul itu **tuntas tanpa pernah diadu dengan dokumen** |
| **11** | ⭐ **Baca ulang kesimpulan kurs standar dan penyimpanan berkas di Claim Prop** | **A4-1**: rule yang dipakai di sana **4 tahun 3 bulan lebih tua** dari yang dinyatakan berlaku |
| **12** | ⭐ **Tinjau ulang keputusan Claim Prop *"52 dari 61 rule `When` tidak dimigrasikan"*** | **A4-3 + A4-5** membuatnya **tidak dapat berdiri bersama** keputusan Fac In tanpa penjelasan |
| ~~**13**~~ | ✅ **DITUNDA 2026-09-19** — *"dibuat nanti saat migrasi, jangan jadiin permasalahan"*. Bukan pemblokir; A5-4 berdiri penuh tanpanya. Teks lamanya: ~~Tetapkan tabel dan kolom polis sumber `Quotation`~~ | ~~A4-3 memutuskan sumbernya berubah, tanpa menyebut tabelnya~~ |

### H1b — ⚠️ Penyimpangan sadar yang tercatat ronde ini: **SATU**

| # | Titik | Di Pega | Di sistem baru | Bagian |
| --- | --- | --- | --- | --- |
| ~~**1**~~ | ⛔ **DICABUT 2026-09-19** — *"ikuti apa adanya, jangan jadikan permasalahan"*. Teksnya dibiarkan sebagai jejak: ~~**Pasangan lompatan lini pada `GetAllData_Act`**~~ | ~~`IsMBU` melompat ke tanda **Travel** dan `IsTravel` ke tanda **MBU**~~ | ~~**tiap lini menuju tandanya sendiri**; logikanya **diperbaiki**~~ → ✅ **DITIRU APA ADANYA** | §B4 |
| **2** | **Sumber `Quotation`** | dibaca dari halaman **`pyWorkPage.Quotation`**, yang **ADA** dan diuji **38 rule `When`** | ⭐ **diambil LANGSUNG dari tabel polis**, dibaca **sekali** saat kasus dimuat; ke-38 rule diganti **satu fungsi klasifikasi** | §A4-3 · §A5-4 |
| **3** | **Daftar penyetuju kosong** | kasus komite lahir **tanpa penyetuju, tanpa pesan** | **penyerahan DITOLAK** dengan galat yang terlihat | §A5-1 |
| **4** | **Sumber `.KomiteNo`** | potongan teks mulai **huruf ke-19** dari kunci instance | **pengenal kasus komitenya** | §A5-2 |
| **5** | **`FlagOnGoingCommitte`** | disimpan sebagai teks **`"Send Commite"`** | **tidak disimpan** — diturunkan | §A5-3 |
| **6** | **Nama kolom hasil SQL** | **5 dari 5** dan **2 dari 3** alias berbohong, dengan pemetaan berbeda | **dinamai sesuai isinya** + tabel pemetaan nama-lama → nama-benar | §A5-5 |

⭐ **Jumlah titik yang terurai kini LIMA.**

> ⛔ **RALAT 2026-09-19** — kalimat lamanya dikutip: *"⭐ **Jumlah titik yang terurai kini ENAM**,
> naik dari dua."* Titik **1** dicabut pada hari yang sama, jadi **enam → lima**.

⚠️ **Ditambah enam lagi yang lahir dari A4-4** — keenam pertentangan ADR **0003 · 0005 · 0008 ·
0011 · 0013 · 0015** kini **resmi menjadi penyimpangan sadar**, karena ADR dinyatakan berlaku
lintas-modul. ⛔ **Belum saya daftarkan satu per satu di sini** — rinciannya ada di ronde ulang
§B10, dan memindahkannya menuntut membaca ulang keenamnya. **Jumlah penyimpangan sadar modul ini
karena itu: 2 yang terurai + 6 yang tercatat di ronde ulang = 8.**

⛔ **Daftar ini wajib ikut terbawa ke spec** ketika modul ini kelak punya spec. Sampai itu terjadi,
**berkas inilah satu-satunya tempat penyimpangan sadar Claim Fac In tercatat.**

### H2 — Kesimpulan ronde 2 yang **PALING RAWAN salah**

⚠️ **Yang paling rawan adalah §B4 — "dua pasang lompatan tertukar, dan lini Travel dirugikan".**

**Kenapa rawan:** ia bersandar pada **satu asumsi yang belum diuji di layar Pega** — bahwa langkah
yang dilompati **benar-benar dilewati**, dan bahwa eksekusi **jatuh berurutan** ke langkah
berikutnya sesudah gerbang menolak. Kedua hal itu saya baca dari **struktur ekspor**, bukan dari
menjalankan sistemnya. Bila Pega ternyata memperlakukan *"lompat ke tanda"* berbeda — misalnya
melanjutkan dari induk, bukan dari tanda — **seluruh kesimpulan Travel runtuh**.

Yang **tidak** rawan dan tetap berdiri: **pasangan gerbang→tanda memang tertukar**, dan itu dibaca
**dua kali dengan dua cara**.

⚠️ **Dan kerawanan itu kini menanggung beban lebih berat**, karena `[keputusan work owner]`
2026-09-19 memutuskan **memperbaiki logikanya di sistem baru**. Bila asumsi *"langkah yang dilompati
benar-benar dilewati, lalu eksekusi jatuh berurutan"* ternyata **keliru**, maka yang diperbaiki
adalah sesuatu yang **tidak pernah rusak** — dan sistem baru justru **menyimpang tanpa sebab**.
⛔ **Menguji satu kasus Travel di layar Pega** akan menutup kerawanan ini sekaligus; itu **satu
pemeriksaan**, bukan satu ronde.

**Kedua paling rawan:** §D1, penggolongan `pyMemo`. Angka **91** milik jendela saya, **118** milik
ronde 1, dan **tidak satu pun dari keduanya terbukti**.

### H3 — RALAT terhadap ronde 1 dan ronde ulang

| # | Yang diralat | Angka/kalimat LAMA, dikutip | Yang benar |
| --- | --- | --- | --- |
| **1** | Penunjuk posisional `CreateKMTNo_Act` | ronde ulang §B5: *"`CreateKMTNo_Act` langkah **6.9** membawa **empat** penunjuk posisional ke kasus komite: `IDObject` · `CoverageSubscript` · `IndexObjectItem` · `IndexObject`"* | ⭐ **LIMA** — yang kelima **`IndexAdjustment` := `local.IndexAdjust`** |
| **2** | Penyebut angka korpus | ronde 1: *"482 berkas"* | **482 berkas `.xml`**; folder berisi **483 berkas** |
| **3** | Pernyataan gerbang aksi tombol | ronde ulang §E3: *"Gerbang aksi tombol: tidak ketemu di 11 medan yang disisir"* | **`pyActionConditions` kosong di 121 medan yang disisir** — tetapi **gerbang tombol ADA**: `pyDisableSubmit = true` **52 titik** |
| **4** | Pasangan alias SQL | ronde ulang §E2: *"**38 dari 50** pasangan TIDAK COCOK"* | jendela saya memberi **39 dari 42**; ⛔ **keduanya berdiri**, tidak ada yang ditimpa |
| **5** | `pyXMLSignature` | ronde ulang §E1: *"gagal dibaca, dan itu dilaporkan"* | ⭐ **BERHASIL** dengan pola ketiga — **169 parameter di 72 activity** |
| **6** | Golongan `pyMemo` | ronde 1 §F1: *"(c) kandidat **MENGUBAH** pembacaan **118**"* | jendela saya memberi **91**; ⛔ **keduanya dilaporkan** |
| **7** | Berkas aturan baca | ronde 1 §A1: *"Aturan A–J dipakai sejak langkah pertama"* | berkasnya **TIDAK ADA**; dasar yang dipakai disebut di §A1 |

### H4 — ⭐ Yang seharusnya dikerjakan tetapi **TIDAK diperintahkan** blok ini

⛔ **Disebutkan, tidak dikerjakan.**

| # | Butir |
| --- | --- |
| **1** | **Menyisir `pxUpdateSystemID` di 20 modul korpus lain** — masih menggantung sejak ronde ulang §H5 |
| **2** | **Memeriksa apakah properti yang diuji tiap `When` benar-benar ditulis** |
| **3** | **Mengadu Claim Prop dan Komite Claim Prop dengan 15 ADR** — keduanya sudah selesai tanpa pernah diadu dengan dokumen |
| **4** | **Membaca `SetListKomite_act`** — pembangun daftar penyetuju, yang §C3 hanya sebut namanya |
| **5** | **Memeriksa kelima rule `ADESAMUEL@` isinya**, bukan hanya identitasnya |
| **6** | **Menghitung ulang 60 rule `When`** dengan identitas empat bagian terhadap **seluruh** 20 modul, bukan hanya Claim Prop |
| **7** | **Membaca `HitServiceToKasir_Act`** — kiriman uang yang di `SaveAcceptation` berjalan **sebelum** `Obj-Save` |

---

## Lampiran — bukti berkas lain tidak disentuh

Sidik jari MD5 diambil **sebelum** ronde ini dan dibandingkan **sesudahnya**:

```
korpus Claim Fac In        482 berkas .xml
berkas lama claim-facin      2 berkas  (712 + 735 baris)
docs\adr\                   15 berkas
modul lain di .scratch\    177 berkas
```

Hasil perbandingan ditulis di laporan ronde ini. ⛔ **Nol berkas korpus dibuka untuk ditulis**;
seluruh pembacaan memakai pengurai XML **hanya-baca**. ⛔ **Nol nilai rahasia disalin.**
