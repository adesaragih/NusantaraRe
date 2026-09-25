# FINDING-002 — Alur bercabang berdasarkan identitas orang

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

**Jenis**: laporan kondisi sistem lama (bukan keputusan migrasi)
**Status**: TERBUKTI DARI XML — tidak menunggu data Oracle
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

Dokumen ini memaparkan mekanismenya apa adanya. Ia **tidak** menilai kinerja, motif, atau kepatutan siapa pun. Nama-nama di bawah muncul karena tertulis di dalam kode, bukan karena dipersoalkan.

---

## 1. Tiga bentuk percabangan

### 1.1 Berdasarkan akun pembuat case

```
pyWorkPage.pxCreateOperator == "VINCENTVERNANDO_1"
```

| Rule | Jumlah langkah | Dibuat |
|---|---|---|
| `Activity\CreateChildKomiteCNP_Act.xml` baris 10720, 10967, 11195, 12582 | 4 | 2018-01-16 oleh `YOSUAAMBIKA` |
| `Activity\SaveDataToOSAksep_Act.xml` baris 8731 | 1 | 2022-01-03 oleh `GABRIELAMILITIA` |
| `Activity\SaveToOS.xml` baris 4511 | 1 | 2022-01-03 oleh `GABRIELAMILITIA` |
| `Activity\SendEmailKlaimRejectClose.xml` baris 693, 884 | 2 | 2023-11-20 oleh `ArlexyVarian` |

Case yang dibuat akun ini mengambil jalur berbeda dari case lain, di empat rule yang berbeda pula.

### 1.2 Berdasarkan nilai `.KomiteID` dan `.PICSuggest`

```
.KomiteID   == "Himawan" | "NANDINA" | "CHRISTINEANGELINA" | "CHRISTOPMARHASAK"
.PICSuggest == "Himawan" | "Nandina C" | "Christine Angelina Hutagalung"
```

| Rule | Jumlah langkah | Dibuat |
|---|---|---|
| `Activity\CreateChildKomiteCNP_Act.xml` baris 11926 dst. | 4 | 2020-02-18 oleh `MESDISILITONGA` |
| `Activity\SethistoryKlaimTreaty.xml` baris 592, 680, 1168, 1316 | 3 | 2022-12-01 oleh `AnanSosmita` |
| `DataTransform\InsertChronology_DT.xml` baris 426 | 1 | — |

Perhatikan `"Himawan"` muncul dalam dua ejaan berbeda: `KomiteID=="Himawan"` dan `PICSuggest=="Himawan"`, sementara rekannya ditulis `NANDINA` di satu tempat dan `Nandina C` di tempat lain, `CHRISTINEANGELINA` dan `Christine Angelina Hutagalung` di tempat lain lagi. Tidak ada satu bentuk kanonik.

### 1.3 Berdasarkan isi kolom komentar bebas — bentuk paling rapuh

```
@contains(.CommentSuggest, "Accepted by Himawan")
@contains(.CommentSuggest, "Rejected by Himawan")
@contains(.CommentSuggest, "Himawan")
```
`Activity\SethistoryKlaimTreaty.xml` baris 1256, 1268, 1404, 1416.

Keputusan alur diambil dengan mencari potongan teks di dalam kolom komentar yang diisi bebas oleh pengguna.

## 2. Apa yang terjadi bila string tidak cocok

Ketiga bentuk di atas dipakai sebagai `pyStepsPreCondParamsWhen`. Dalam Pega, precondition yang bernilai salah menyebabkan langkah itu **dilewati**, bukan menimbulkan galat.

Konsekuensinya sama untuk ketiganya, dan tidak bergantung pada tafsir:

| Pemicu | Akibat |
|---|---|
| Akun `VINCENTVERNANDO_1` diganti nama atau dinonaktifkan | 8 langkah berhenti berjalan. Tidak ada pesan galat. |
| Pemegang peran berganti orang | Cabang untuk nama lama tidak pernah terpicu lagi; nama baru tidak punya cabang |
| Komentar diketik `"Accepted By Himawan"` (huruf besar B) | `@contains` gagal, langkah dilewati |
| Komentar diketik `"Accepted by Pak Himawan"` | `@contains(.CommentSuggest,"Himawan")` tetap cocok, tetapi `"Accepted by Himawan"` tidak — dua kondisi bersaudara memberi hasil berbeda atas satu masukan |
| Pengguna mana pun mengetik nama itu di kolom komentar | Kondisi terpenuhi, langkah berjalan |

Butir terakhir adalah sifat mekanismenya: **kolom komentar bebas adalah masukan yang menentukan alur**, dan siapa pun yang boleh mengisi komentar dapat memenuhinya.

## 3. Batas klaim — apa yang TIDAK dapat saya simpulkan dari XML

1. **Apa yang sebenarnya dilakukan tiap cabang secara bisnis** tidak dapat dibaca dari kondisi itu sendiri; yang terbaca hanya langkah yang dijaganya.
2. **Apakah orang-orang itu masih menjabat** adalah pertanyaan organisasi → `_selesai/OPEN-QUESTIONS.md` A10, A11.
3. **Apakah cabang-cabang itu pernah benar-benar terpicu di produksi** adalah pertanyaan data, bukan pertanyaan kode.
4. Sebagian kondisi berada di `CreateChildKomiteCNP_Act`, yang sisi seberangnya ada di modul Komite → sebagian akibatnya `DEFERRED-TO-KOMITE-SESSION`.

## 4. Hubungan dengan sistem baru

Tidak ada cabang berbasis identitas yang boleh diwarisi. Ini bertemu dengan ADR-0006 (`rbac-dirancang-dari-nol`): karena `pyPrivilegeName` kosong di seluruh 279 berkas, satu-satunya model kewenangan yang benar-benar berjalan di sistem lama justru **kondisi-kondisi ini** — kewenangan yang ditulis sebagai nama orang di dalam kode, bukan sebagai peran.

Pemetaan tiap cabang ke **peran** (bukan ke orang) adalah prasyarat desain RBAC, dan itu memerlukan jawaban A10/A11.

## 5. Perilaku tiap cabang, dan kewenangan yang tersirat — EVIDENCED

Diminta pada Round 3: sebelum bertanya ke siapa pun, baca sendiri apa yang berubah saat kondisinya terpenuhi. Hasilnya menjawab sebagian besar pertanyaan tanpa perlu keluar.

| Nama di kode | Rule & langkah | Apa yang berubah saat kondisi terpenuhi | Kewenangan yang tersirat |
|---|---|---|---|
| `CHRISTINEANGELINA` | `CreateChildKomiteCNP_Act` | `.IDKomite="Claim Dept. Head"`, `.KomitePost="Department Head Claim"`, `.Initial="CA"` | **bukan kewenangan** — pemetaan nama ke jabatan |
| `CHRISTOPMARHASAK` | idem | `.IDKomite="Technic Div. Head"`, `.Initial="CM"` | **bukan kewenangan** — pemetaan nama ke jabatan |
| `Himawan` | idem | `.IDKomite="Operational Director"`, `.Initial="HY"` | **bukan kewenangan** — pemetaan nama ke jabatan |
| `NANDINA` | idem | `.IDKomite="Technical Director"`, `.Initial="NC"` | **bukan kewenangan** — pemetaan nama ke jabatan |
| `DARTO` | `CreateChildKomiteCNP_Act` baris 10668 | `.OPERATOR_ID` dipetakan ulang menjadi `"CHRISTINEANGELINA"` | **bukan kewenangan** — penggantian identitas |
| `Christine Angelina Hutagalung` / `Himawan` / `Nandina C` | `SethistoryKlaimTreaty` | `.IsCedingConfirm` diisi nama jabatan yang sama | **bukan kewenangan** — pemetaan nama ke jabatan |
| `VINCENTVERNANDO_1` | `CreateChildKomiteCNP_Act` (4 langkah) | daftar Komite hasil `FilterEmailKomiteWithLimit` **diganti satu anggota tetap**: `KomiteID="VINCENTVERNANDO_1"`, `KomiteEmail="klaim5@nusantarare.com"`, `IDKomite="1"`, `KomitePost="VC"`; dan **`ChildWorkPage.KomiteLoop = 1`** | **KEWENANGAN — mengesampingkan penjenjangan** |
| `VINCENTVERNANDO_1` | `SendEmailKlaimRejectClose` (2 langkah) | `Local.EmailCC` diganti daftar tetap | distribusi pemberitahuan |
| `VINCENTVERNANDO_1` | `SaveToOS` langkah 6.3.10 | melewati **`Call KonversiKlaim_Act`** — langkah berdeskripsi *"Untuk HIT ke arasapas"* | **KEWENANGAN — melewati integrasi hilir** |
| `VINCENTVERNANDO_1` | `SaveDataToOSAksep_Act` langkah 15.8.8 | mengubah cabang di sekitar **`Connect-REST`**, juga berdeskripsi *"Untuk HIT ke arasapas"* | **KEWENANGAN — melewati integrasi hilir** |

### 5.1 Empat dari lima nama bukan kewenangan sama sekali

Keempatnya hanya menerjemahkan **nama orang menjadi nama jabatan**. Itu tabel referensi yang ditulis sebagai kode:

```
CHRISTINEANGELINA  ->  Claim Dept. Head
CHRISTOPMARHASAK   ->  Technic Div. Head
Himawan            ->  Operational Director
NANDINA            ->  Technical Director
DARTO              ->  (diganti menjadi CHRISTINEANGELINA)
```

Mengikuti kriteria Round 3 — *"kalau ia hanya mengisi nilai default berbeda, itu bukan kewenangan, itu preferensi, dan langsung NOT-MIGRATED tanpa perlu bertanya ke siapa pun"* — kelimanya **`CANDIDATE-NOT-MIGRATED`**, dan tidak perlu dibawa ke pemilik proses. Penggantinya di sistem baru adalah tabel `pengguna → jabatan` biasa.

Dikuatkan oleh `MEMORI_PEMAHAMAN.MD` baris 1167–1174, yang memuat pemetaan yang sama persis dan sudah menandainya sebagai risiko: *"perubahan personel memerlukan perubahan kode."* Label naik dari DERIVED menjadi **EVIDENCED**.

### 5.2 `VINCENTVERNANDO_1` adalah satu-satunya yang benar-benar kewenangan

Alurnya normal: `FilterEmailKomiteWithLimit` menghasilkan daftar anggota Komite, lalu `ChildWorkPage.KomiteLoop = @SizeOfPropertyList(pyReportContentPage.pxResults)` menetapkan berapa jenjang persetujuan yang dilalui.

Untuk case yang dibuat akun ini, keduanya dilewati: daftarnya diganti **satu anggota tetap**, dan `KomiteLoop` dipaksa **`1`**.

Artinya: **case yang dibuat akun ini hanya melalui satu persetujuan, siapa pun anggotanya seharusnya, berapa pun nilai klaimnya.** Ini wewenang mengesampingkan penjenjangan Komite — kategori paling berat dari tiga kategori di Round 3.

`MEMORI_PEMAHAMAN.MD` baris 1174 menyebut jalur ini juga mem-bypass `KonversiKlaim_Act` dan menyetel `InputData.CARI20 = 1` pada `InsertOSKlaimCNP` — cakupannya lebih luas dari yang terbaca di folder ini saja.

**Hubungan dengan FINDING-001:** keduanya menyentuh mekanisme penjenjangan yang sama. FINDING-001 mempersoalkan ambang yang dibandingkan tanpa konversi mata uang; temuan ini menunjukkan ada jalur yang **melewati ambang itu sama sekali**. Keduanya berdiri sendiri dan tidak saling menggugurkan.

**Satu-satunya pertanyaan yang tersisa untuk dibawa keluar**, karena tidak terjawab XML maupun memori — dan `MEMORI_PEMAHAMAN.MD` sendiri mencatatnya sebagai pertanyaan terbuka nomor 16: **apakah `VINCENTVERNANDO_1` akun uji atau akun bisnis nyata, dan apakah masih aktif.**

---

## 6. Sapuan pola pencocokan teks — seluruh folder

Diminta pada Round 3: yang satu itu hampir pasti bukan satu-satunya. Benar.

**Enam pasang**, bukan dua. Seluruhnya di `SethistoryKlaimTreaty`, dan semuanya melakukan hal yang sama: **menerjemahkan nama orang di dalam teks komentar menjadi nama jabatan**.

| Kondisi | Akibat |
|---|---|
| `@contains(.CommentSuggest,"Accepted by CHRISTINEANGELINA")` | `.CommentSuggest = "Accepted by Kepala Departemen Klaim"` |
| `@contains(.CommentSuggest,"Rejected by CHRISTINEANGELINA")` | `.CommentSuggest = "Rejected by Kepala Departemen Klaim"` |
| `@contains(.CommentSuggest,"Accepted by Himawan")` | `.CommentSuggest = "Accepted by Direktur Operational"` |
| `@contains(.CommentSuggest,"Rejected by Himawan")` | `.CommentSuggest = "Rejected by Direktur Operational"` |
| `@contains(.CommentSuggest,"Accepted by NANDINA")` | `.CommentSuggest = "Accepted by Direktur Teknik"` |
| `@contains(.CommentSuggest,"Rejected by NANDINA")` | `.CommentSuggest = "Rejected by Direktur Teknik"` |

Ditambah `@contains(.CommentSuggest,"NANDINA")` (4x), `"Himawan"` (4x), `"CHRISTINEANGELINA"` (2x) sebagai kondisi tanpa awalan.

### 6.1 Pola pencocokan teks lain di folder ini

Sapuan `@contains` / `@startsWith` / `@endsWith` / `@indexOf` / `@matches` menemukan 27 ekspresi. Yang **bukan** identitas orang:

| Ekspresi | Sifat |
|---|---|
| `@contains(.TreatyName,"R/I")` (4x) | penanda jenis treaty di dalam nama |
| `@contains(.IsCedingConfirm,"Admin")` (2x) | mencari jabatan di dalam field yang diisi nama jabatan |
| `@contains(SearchData.CARI2,"/R0")`, `@contains(pyWorkPage.ClaimData.NoPla,"/")` | penguraian nomor lewat teks |
| `@contains(InputSpreading.CARI1,"TH"/"ST"/"RD"/"ND")` | **menebak akhiran bilangan urut bahasa Inggris** dari teks |
| `@contains(pyWorkPage.ClaimData.PolicyData.PolicyNo,"RNM-Q")` | penanda jenis polis di dalam nomor |
| `@contains(Local.Emailto,"syariah")` | pemilahan unit bisnis dari alamat surel |
| `@contains(Local.CheckOld,"Previous")` | penanda status di dalam teks |
| `@startsWith(param.WorkStatus,"Resolved-")` | pola standar Pega, wajar |
| `@indexOf("fdf.fsf.sfsdf",".")` (2x) | **string uji coba tertinggal di kode produksi** |

Dua yang paling rapuh setelah identitas orang: `"syariah"` di alamat surel (unit bisnis ditentukan dari isi alamat email) dan `@indexOf("fdf.fsf.sfsdf",".")` yang jelas sisa percobaan.

### 6.2 Ejaan yang sama tidak konsisten — sebagian cabang tidak akan pernah terpicu

**Koreksi atas versi pertama.** Saya menulis *"orang yang sama ditulis dalam dua ejaan di dalam rule yang sama"*. **Itu keliru.** Penelusuran per baris menunjukkan dua ejaan itu berada di **dua rule yang berbeda**:

| Rule | Bentuk yang dicocokkan |
|---|---|
| `Activity\SethistoryKlaimTreaty.xml` baris 538, 680, 822 | `"Christine Angelina Hutagalung"`, `"Himawan"`, `"Nandina C"` — **bentuk panjang** |
| `DataTransform\InsertChronology_DT.xml` baris 367, 426, 486 | `"CHRISTINEANGELINA"`, `"Himawan"`, `"NANDINA"` — **bentuk pendek** |

Kesimpulannya tetap berlaku, bahkan lebih tajam: **dua rule berselisih tentang bentuk data pada field yang sama.** Karena `.PICSuggest` hanya dapat berisi satu bentuk pada satu waktu:

- **`Himawan`** dieja sama di kedua rule — kedua cabangnya bekerja.
- **`Christine` dan `Nandina`** dieja berbeda — untuk keduanya **tepat satu dari dua rule tidak akan pernah cocok.** Mana yang mati bergantung pada bentuk yang sebenarnya tersimpan di produksi, dan itu pertanyaan data, bukan pertanyaan kode.

Perhatikan pula: kondisi `.PICSuggest` memakai bentuk panjang (`"Christine Angelina Hutagalung"`), sedangkan kondisi `@contains(.CommentSuggest,...)` memakai bentuk pendek (`"CHRISTINEANGELINA"`). Dua mekanisme untuk orang yang sama, dengan ejaan yang berbeda.

---

## 7. Risiko migrasi data — dinyatakan sekarang, bukan ditemukan belakangan

Diminta pada Round 3. Ini konsekuensi langsung dari bagian 6.

**Case yang alurnya pernah ditentukan oleh pencocokan teks tidak menyimpan keputusan itu di tempat lain mana pun selain di dalam teks komentar.**

Saat migrasi, keadaan case seperti itu harus direkonstruksi dari isi `.CommentSuggest`. Yang menggagalkannya:

1. **Teks sudah ditimpa.** Langkah-langkah di bagian 6 **menulis ulang** `.CommentSuggest` — dari `"Accepted by Himawan"` menjadi `"Accepted by Direktur Operational"`. Jadi teks yang tersimpan di produksi sebagian sudah bentuk terjemahan, sebagian masih bentuk asli, tergantung apakah langkah itu sempat berjalan. **Dua populasi bercampur di satu kolom.**
2. **Ejaan berbeda tidak terbaca.** `"Accepted By Himawan"`, `"Accepted by Pak Himawan"`, atau nama yang diketik dengan spasi berlebih tidak akan cocok dengan pola mana pun, dan keputusannya hilang.
3. **Nama di luar kelima itu tidak punya cabang sama sekali.** Persetujuan oleh orang lain tidak pernah diterjemahkan, dan tidak ada jejak strukturnya.

**Akibatnya untuk rencana migrasi:** kolom "siapa menyetujui" pada sistem baru **tidak dapat diisi lengkap dari data lama secara otomatis.** Sebagian akan kosong, dan besarnya bagian yang kosong tidak dapat diketahui dari XML. Kuantifikasinya perlu satu profil data yang murah dan tanpa data nasabah — dimasukkan ke `pengetahuan/PULL-LIST.csv` sebagai **REQ-013**.

Ini **tidak** menghalangi keputusan Round 3 tentang `tindakan + pelaku + waktu` terstruktur. Ia hanya menetapkan bahwa **sebagian riwayat lama tidak akan terbawa**, dan itu harus disepakati sekarang, bukan ditemukan saat UAT.

---

## 8. Jalur yang melewati penjenjangan Komite — `VINCENTVERNANDO_1`

Bagian ini berdiri sendiri, sejajar dengan FINDING-001, karena sifatnya berbeda dari enam bagian di atas. Yang lain adalah pemetaan nama ke jabatan; yang ini mengubah berapa persetujuan yang harus dilalui sebuah klaim.

Netral dan faktual. Tanpa penilaian niat, tanpa menyebut siapa pun bersalah.

### 8.1 Mekanisme normal

`CreateChildKomiteCNP_Act` menyusun daftar anggota Komite dari hasil Report Definition `FilterEmailKomiteWithLimit`, yang disaring dengan `Param.LIMIT_BOTTOM` menurut nilai klaim. Jumlah jenjang persetujuan lalu ditetapkan dari banyaknya baris hasil:

```
ChildWorkPage.KomiteLoop = @SizeOfPropertyList(pyReportContentPage.pxResults)
```

### 8.2 Mekanisme untuk case yang dibuat akun ini

Bila `pyWorkPage.pxCreateOperator == "VINCENTVERNANDO_1"`, empat langkah menggantikan keduanya:

| Yang diganti | Nilai pengganti |
|---|---|
| `Primary.ComiteeClaim` | satu baris tetap: `KomiteID="VINCENTVERNANDO_1"`, `KomiteEmail="klaim5@nusantarare.com"`, `IDKomite="1"`, `KomitePost="VC"`, `Initial="VC"`, `KomiteAproval=""` |
| `ChildWorkPage.KomiteList` | satu baris tetap dengan isi yang sama |
| `ChildWorkPage.KomiteLoop` | **`1`** |

`Activity\CreateChildKomiteCNP_Act.xml` baris 10720, 10967, 11195, 12582.

### 8.3 Akibatnya

**Case yang dibuat akun ini melalui satu persetujuan, berapa pun nilai klaimnya, dan oleh satu alamat tetap.** Hasil `FilterEmailKomiteWithLimit` tidak ikut menentukan apa pun pada jalur ini, sehingga ambang `LIMIT_BOTTOM` — inti dari seluruh mekanisme penjenjangan — tidak berlaku.

Dua langkah lain pada rule berbeda memakai kondisi yang sama: `SendEmailKlaimRejectClose` mengganti `Local.EmailCC` dengan daftar tetap; `SaveDataToOSAksep_Act` dan `SaveToOS` mengubah cabang langkah tanpa `Property-Set` di blok itu sendiri (tujuannya belum ditelusuri).

`MEMORI_PEMAHAMAN.MD` baris 1174 mencatat jalur ini juga mem-bypass `KonversiKlaim_Act` dan menyetel `InputData.CARI20 = 1` pada `InsertOSKlaimCNP` — jadi cakupannya melampaui apa yang terbaca di folder ini.

### 8.4 Beda derajat kepastian dengan FINDING-001

| | FINDING-001 | Bagian ini |
|---|---|---|
| Yang terbaca | ambang dibandingkan tanpa konversi mata uang | daftar Komite diganti satu orang tetap, `KomiteLoop` dipaksa 1 |
| Penjelasan alternatif yang masuk akal | **ada** — `.ValueAdjustment` mungkin sudah IDR sejak hulu | **tidak ada** — mekanismenya eksplisit dan tidak ambigu |
| Yang belum diketahui | apakah ada Adjustment non-IDR di produksi | apakah akun itu dipakai di produksi, dan sejak kapan |

Keduanya menyentuh mekanisme penjenjangan yang sama dan tidak saling menggugurkan.

### 8.5 Yang menutup bagian ini

**REQ-016**, prioritas **BLOCKER**, tanpa data nasabah:

| Pertanyaan | Menentukan |
|---|---|
| Berapa case dibuat oleh `VINCENTVERNANDO_1` | bila nol, jalur ini mati dan bagian ini gugur |
| Rentang tanggal case pertama sampai terakhir | akun uji yang tertinggal, atau jalur yang masih hidup |
| Berapa di antaranya bernilai di atas ambang `30.000.000` | besaran paparan |
| Berapa yang tercatat hanya punya satu persetujuan | apakah mekanismenya benar-benar berjalan seperti terbaca |

Pertanyaan "akun uji atau akun bisnis nyata" tetap EXTERNAL (`_selesai/OPEN-QUESTIONS.md` A14) dan sudah tercatat sebagai pertanyaan terbuka nomor 16 di `MEMORI_PEMAHAMAN.MD`.

---

## 9. `DARTO` — penggantian identitas, bukan pemetaan jabatan

Versi pertama menempatkan `DARTO` dalam tabel yang sama dengan empat nama lainnya. **Itu salah penempatan**: jenisnya berbeda.

Empat nama lain adalah pencarian **nama ke jabatan** (`CHRISTINEANGELINA` menghasilkan `"Claim Dept. Head"`). `DARTO` melakukan hal lain:

```
.OPERATOR_ID == "DARTO"   ->   .OPERATOR_ID diganti menjadi "CHRISTINEANGELINA"
```
`Activity\CreateChildKomiteCNP_Act.xml` baris 10580 (deskripsi) dan 10668 (kondisi).

**Yang diganti adalah identitas pelakunya, bukan atributnya.** Tindakan yang dilakukan satu akun tercatat atas nama akun lain.

### 9.1 Konsekuensinya bukan RBAC, melainkan keutuhan jejak audit

Untuk case yang melewati langkah ini, riwayat di sistem lama **tidak menunjukkan siapa yang sebenarnya bertindak**. Yang tercatat adalah identitas pengganti.

### 9.2 Konsekuensi migrasi — dinyatakan sekarang

> **Riwayat yang dimigrasi membawa identitas yang mungkin tidak benar, dan tidak ada cara memulihkannya.**

Nilai aslinya tidak disimpan di mana pun sebelum diganti — penggantiannya berupa `Property-Set` langsung, bukan penyalinan ke field cadangan. Setelah `.OPERATOR_ID` bernilai `"CHRISTINEANGELINA"`, tidak ada jejak bahwa sebelumnya ia `"DARTO"`.

Tiga akibat yang harus disepakati, bukan ditemukan belakangan:

1. **Baris riwayat yang terdampak tidak dapat dikoreksi saat migrasi.** Kita tidak tahu baris mana yang terdampak, karena hasilnya tidak dapat dibedakan dari baris yang memang dibuat `CHRISTINEANGELINA`.
2. **Jumlahnya tidak dapat diukur** — bahkan dengan query. Tidak ada penanda.
3. **Sistem baru tidak boleh menyediakan mekanisme serupa.** Perwakilan atau pendelegasian, bila memang dibutuhkan, dicatat sebagai *"X bertindak untuk Y"* dengan kedua identitas tersimpan — bukan dengan menimpa salah satunya.

---

## 10. Cabang yang tidak mungkin terpicu — otomatis tidak dimigrasi

Diminta pada Round 4. Hasil sapuan seluruh `pyStepsPreCondParamsWhen` di 279 berkas.

### 10.1 Yang terbukti

| Cabang | Kenapa mati | Rule |
|---|---|---|
| `.PICSuggest=="CHRISTINEANGELINA"` **atau** `"Christine Angelina Hutagalung"` — salah satu | Dua rule mencocokkan bentuk berbeda atas field yang sama | `InsertChronology_DT` vs `SethistoryKlaimTreaty` |
| `.PICSuggest=="NANDINA"` **atau** `"Nandina C"` — salah satu | idem | idem |

Mana dari pasangan itu yang mati **belum dapat ditentukan dari XML** — ia bergantung bentuk yang tersimpan di produksi. Yang pasti: untuk setiap pasangan, **tepat satu sisi tidak pernah cocok.**

### 10.2 Yang saya cari tetapi TIDAK ditemukan

Supaya tidak dibaca sebagai daftar yang lebih besar dari sebenarnya, ini yang saya uji dan hasilnya nihil:

| Pola yang dicari | Hasil |
|---|---|
| Kondisi membandingkan literal dengan literal (`"a"=="b"`) | **TIDAK DITEMUKAN DI XML** |
| Kondisi self-contradictory (`X=="a" && X=="b"`) | **TIDAK DITEMUKAN DI XML** |
| String uji coba di dalam **kondisi** (`fdf`, `asdf`, `test123`) | **TIDAK DITEMUKAN DI XML** |

Catatan atas yang terakhir: `@indexOf("fdf.fsf.sfsdf",".")` memang ada, tetapi **bukan** sebagai precondition — ia berada di dalam ekspresi nilai. Jadi ia sisa percobaan yang tetap dieksekusi, bukan cabang mati. Saya sempat menggolongkannya bersama cabang mati pada ringkasan Round 4; itu penggolongan yang salah.

Satu kandidat lagi yang saya periksa dan **gugur**: `HitServiceToKasir_Act` memuat `(pyWorkIDPrefix=="CLMP-" && ...) || (pyWorkIDPrefix=="CLM-" && ...)`. Terlihat seperti kontradiksi, tetapi itu dua kelompok yang di-OR — keduanya dapat terpenuhi. Bukan cabang mati.

### 10.3 Berapa besar penghematannya

Dua pasang cabang, tiga langkah. **Lebih kecil dari dugaan awal.** Penghematan lingkup kerjanya nyata tetapi kecil, dan saya lebih baik menyampaikan angka yang benar daripada daftar yang panjang.

### 8.6 Cakupannya lebih luas dari penjenjangan Komite

Telusuran langkah tujuan pada `SaveToOS` dan `SaveDataToOSAksep_Act` selesai, dan hasilnya memperluas temuan ini.

| Rule | Langkah | Yang dilewati |
|---|---|---|
| `SaveToOS` | `pySteps(6).pySteps(3).pySteps(10)` | `Call KonversiKlaim_Act` — *"Untuk HIT ke arasapas"* |
| `SaveDataToOSAksep_Act` | `pySteps(15).pySteps(8).pySteps(8)` | cabang di sekitar `Connect-REST` — *"Untuk HIT ke arasapas"* |

Jadi untuk case yang dibuat akun ini, bukan hanya penjenjangan Komite yang diganti — **pengiriman ke Arasapas, sistem hilir, juga tidak berjalan sebagaimana case lain.**

Ini mengonfirmasi catatan `MEMORI_PEMAHAMAN.MD` baris 1174 yang menyebut jalur ini mem-bypass `KonversiKlaim_Act`, dan menaikkan statusnya dari **memori** menjadi **EVIDENCED**.

**Akibatnya untuk REQ-016**: pertanyaannya bertambah satu. Bukan hanya berapa case yang dibuat akun ini dan berapa persetujuan yang dilaluinya, tetapi juga **berapa di antaranya tidak pernah sampai ke Arasapas**. Bila ada case produksi yang tidak terkirim ke sistem hilir, itu selisih data antar sistem yang belum pernah masuk peta mana pun.
