# FAKTA A-1B — Pembentukan Sirkulasi

> Ronde 4, 20 September 2026. Pembacaan empat rule yang selama tiga ronde memagari A-1b, dan yang ternyata sudah ada di repo (`INVENTARIS-BUKTI.md` §2.4).
> Berkas ini memuat **fakta**, bukan rancangan. Tanda: `EVIDENCED` (berkas·langkah) · `DIPULIHKAN` (pengikatan lewat pencocokan nama halaman bind — tidak pernah naik menjadi `EVIDENCED`) · `TAFSIR`.
> Sumber: `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` (ver 01-01-28, 2026-08-04) · `CreateChildKomiteCloseNP_Act.xml` (01-01-24, 2025-10-01) · `ProteksiSendKomiteCNP_Act.xml` (01-01-24, 2025-12-30) · `ReportDefinition/FilterEmailKomiteWithLimit.xml` (01-01-17, 2024-08-19).

---

## F-13 · Jumlah jenjang lahir dari cacah baris, bukan dari aturan

`CreateChildKomiteCNP_Act` langkah 9 menyetel `ChildWorkPage.KomiteLoop = 1`, lalu langkah 26.13 **dan** 26.14 menimpanya dengan nilai yang sama:

```
ChildWorkPage.KomiteLoop = @Utilities.SizeOfPropertyList(pyReportContentPage.pxResults)
```

Langkah 26.15 menimpanya sekali lagi menjadi `1` bila `pyWorkPage.pxCreateOperator == "VINCENTVERNANDO_1"`. `EVIDENCED`.

Tidak ada tempat mana pun yang menyimpan "berapa jenjang seharusnya". Jumlah jenjang **adalah** banyaknya baris roster yang lolos saring, dihitung saat sirkulasi dibuat.

**Akibat yang ikut tertarik:** `KomiteLoop` bukan aturan yang perlu dimigrasikan, ia turunan — sejalan D-2 yang sudah menolak `KOMITECOUNT`/`FLAGONGOINGCOMMITTE` sebagai kolom. Yang harus dimigrasikan adalah **penyaringnya**, bukan bilangannya.

---

## F-14 · Penyaring roster: tiga syarat, satu urutan, satu kolom yang dipilih tapi tak pernah dipakai

`FilterEmailKomiteWithLimit`, class `ASM-FW-GCNMFW-Int-EMAILKOMITE`. `pyFilterLogic = A AND C AND B`. `EVIDENCED`:

| Label | Syarat |
|---|---|
| A | `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM` |
| B | `.STS_AKTIF = "1"` (nilai tetap di dalam rule) |
| C | `.STS_KLAIM = Param.STS_KLAIM` |

Parameter: `LIMIT_BOTTOM` (Double), `STS_KLAIM` (Text). Urutan hasil: **`.DEGREE` menaik** (`pySortType=ASC`, `pySortOrder=1`) — satu-satunya kolom berurutan. Delapan belas kolom diambil, termasuk `.LIMIT_TOP`, `.JABATAN`, `.TYPE_KOMITE`, `.STS_REJECT`, `.STS_ADJ`, `.STS_REG`, `.STS_SURVEY`, `.STS_SALVAGE`, `.STS_ADJUSTER`.

**`.LIMIT_TOP` diambil tetapi tidak pernah menyaring.** Begitu pula `.TYPE_KOMITE` dan enam kolom `STS_*` lainnya. `EVIDENCED`.

**Akibat yang ikut tertarik:** urutan jenjang persetujuan = urutan `DEGREE` pada roster, jadi urutan itu **data, bukan kode** — sejalan D-5. Dan karena hanya batas bawah yang menyaring, tidak ada satu pun anggota yang pernah dikecualikan karena klaimnya *terlalu kecil* baginya: penyaring menghasilkan **semua orang yang berwenang sampai nilai itu**, bukan orang yang tepat untuk nilai itu.

---

## F-15 · Tidak ada nilai klaim yang dibandingkan dengan ambang — pembandingnya dimatikan

`Param.LIMIT_BOTTOM` tidak pernah diisi dari nilai apa pun. Ia diisi konstanta (`CreateChildKomiteCNP_Act` langkah 26.3–26.6), `EVIDENCED`:

| Langkah | Syarat | `Param.LIMIT_BOTTOM` |
|---|---|---|
| 26.3 | (tanpa syarat, dasar) | **tidak diisi sama sekali** — hanya `STS_KLAIM="NONPROP"` |
| 26.4 | `Flagkomite==1` | `0` |
| 26.5 | `Flagkomite==2` | `25000001` |
| 26.6 | `Subjectivity==true` | `0` |

Deskripsi langkah 26.3 dan 26.6 berbunyi persis: *"Set Parameter RD -- Untuk sementara dihilangkan (Param.LIMIT_BOTTOM==Local.TotalValueAdjust)"*. Jadi pembanding berbasis nilai memang pernah ada dan **sengaja dimatikan**, dan kalimat "untuk sementara" itu bertahan sampai versi 01-01-28.

**Jawaban langsung atas pertanyaan 3 ronde ini:** ambang **dapat** dievaluasi di dalam agregat Komite tanpa masukan nilai dari hulu — karena hari ini tidak ada nilai yang dievaluasi. Yang butuh masukan hulu hanyalah `Flagkomite` (F-16).

**Akibat yang ikut tertarik:** seluruh pertanyaan "per mata uang atau setelah konversi, kotor atau bersih" **tidak punya jawaban di sistem lama** — bukan karena tidak terbaca, melainkan karena tidak ada perbandingannya. Pertanyaan itu hidup lagi hanya bila ambang berbasis nilai dipulihkan (Q-1).

---

## F-16 · `Flagkomite`: empat konstanta, satu cabang mati, satu jalur tanpa parameter

`CreateChildKomiteCNP_Act` langkah 10 menetapkan: `LimitMax = 30.000.000,00` · `LimitMaxDivHead = 50.000.000,00` · `LimitPersenMax = 30,00` · `LimitPersenMaxDivHead = 30,00` · `Flagkomite = 0`. Pembandingnya `CekLimitPersen = pyWorkPage.TreatyInMaster.RNMShare`. `EVIDENCED`.

| Langkah | Syarat | Hasil |
|---|---|---|
| 12 | `RNMShare <= 30` **dan** `nilai <= 30jt` | `Flagkomite = 1` |
| 13 | `RNMShare <= 30` **dan** `30jt < nilai <= 50jt` | `Flagkomite = 2` |
| 14 | `RNMShare > 30` **dan** `RNMShare <= 30` **dan** `nilai <= 50jt` | `Flagkomite = 2` |

**Langkah 14 tidak pernah dapat menyala**: kedua konstanta pembatasnya bernilai `30,00`, sehingga syaratnya adalah himpunan kosong. Deskripsi langkahnya sendiri berbunyi *"LimitPersenMaxDivHead, CekLimitPersen >15 && <=30"* — angka 15 yang disebut deskripsi tidak ada di kode. `EVIDENCED`.

**Akibat yang paling keras:** bila `nilai > 50jt`, atau `RNMShare > 30`, tak satu pun cabang menyala → `Flagkomite = 0` → langkah 26.4/26.5 dilewati → **`Param.LIMIT_BOTTOM` tidak pernah diisi**, dan penyaring A dijalankan dengan parameter kosong. Perilaku runtime-nya (nol baris, seluruh baris, atau galat) **tidak dapat ditentukan dari XML** — itu perilaku Pega dan Oracle atas parameter Double yang tak berisi, bukan sesuatu yang tertulis di rule. `TAFSIR` berhenti di sini.

Ini berarti adjustment **terbesar** justru jatuh ke jalur yang paling tidak tertentu.

**Akibat yang ikut tertarik:** `RNMShare` dan nilai adjustment adalah dua masukan hulu yang tetap dibutuhkan agregat Komite saat sirkulasi dibuat, sekalipun ambangnya sendiri konstanta.

---

## F-17 · Yang dinamai "TotalValueAdjust" bukan total

`CreateChildKomiteCNP_Act` langkah 11 mengulang `pyWorkPage.ClaimData.AdjustmentList` dan di setiap putaran menyetel:

```
Local.TotalValueAdjust = .ValueAdjustment
```

Tanpa penjumlahan dan tanpa penyaring. Sesudah perulangan, nilainya adalah `ValueAdjustment` **baris terakhir** daftar — bukan total, dan belum tentu adjustment yang sedang dikirim ke komite. Nilai inilah yang dipakai langkah 12–14 untuk menentukan `Flagkomite`. `EVIDENCED`.

**Penerapan default D-1, saya jalankan tanpa bertanya:** ini cacat ber-`EVIDENCED`, jadi perbaikan menjadi default dan paritas yang harus dibela. Nilai yang menentukan kelas kewenangan adalah nilai **adjustment yang sedang diajukan**. Masuk **DEVIASI DIHARAPKAN**, dengan uji paritas yang dirancang gagal pada klaim ber-lebih-dari-satu adjustment.

**Akibat yang ikut tertarik:** pada klaim dengan satu adjustment, cacat ini tidak terlihat — maka uji paritasnya wajib memakai klaim dengan minimal dua adjustment bernilai beda kelas, kalau tidak deviasinya lolos tanpa terdeteksi.

---

## F-18 · Dua daftar diisi berbeda, lalu yang benar ditimpa yang salah

Di dalam perulangan hasil roster (`CreateChildKomiteCNP_Act` langkah 26.8), `EVIDENCED`:

| Sasaran | `IDKomite` diisi dari |
|---|---|
| `Primary.ComiteeClaim` (daftar di induk) | `.ID` — pengenal baris roster |
| `ChildWorkPage.KomiteList` (daftar di case Komite) | **`.JABATAN`** — jabatan dari roster |

Lalu langkah 26.13 menimpa seluruhnya:

```
ChildWorkPage.KomiteList = .ComiteeClaim
```

Sehingga `.JABATAN` yang sudah benar **tidak pernah sampai** ke case Komite; yang sampai adalah `.ID`. `EVIDENCED`.

**Akibat yang ikut tertarik:** inilah sebab `SetKomiteList_Act` di sisi Komite harus menebak jabatan dari nama orang (lima cabang `DARTO`/`CHRISTINEANGELINA`/`CHRISTOPMARHASAK`/`Himawan`/`NANDINA`), dan sebab kronologi hanya mengenal tiga dari lima nama. K-03 bukan kemalasan penulis rule Komite — ia **gejala** dari penimpaan ini. Penerapan default D-1: jabatan diambil dari roster, cacat ber-`EVIDENCED` ini diperbaiki, masuk DEVIASI DIHARAPKAN.

---

## F-19 · Satu substitusi orang ditulis di dalam kode

`CreateChildKomiteCNP_Act` langkah 26.8.3, `EVIDENCED`:

```
JIKA .OPERATOR_ID == "DARTO"  →  Primary.ComiteeClaim(<LAST>).KomiteID = "CHRISTINEANGELINA"
```

Baris roster tetap milik `DARTO` (email, ID, jabatan ikut dari barisnya), tetapi **yang berhak memutus diganti** menjadi orang lain. Karena langkah 26.13 menyalin daftar induk ke case Komite, substitusi ini ikut terbawa.

**Akibat yang ikut tertarik:** pengalihan pemegang jenjang **sudah terjadi hari ini**, hanya saja berwujud satu baris kode tanpa pelaku, tanpa alasan, dan tanpa waktu. E-3 sudah memutuskan bentuk yang benar (jenjang = kedudukan, perpindahan = peristiwa); fakta ini menunjukkan kebutuhannya nyata, bukan antisipatif. Yang belum terjawab hanyalah siapa yang berwenang melakukannya → Q-3.

---

## F-20 · Jalur tutup/tolak lahir dari rule lain, tanpa roster dan tanpa ambang

`CreateChildKomiteCloseNP_Act`, `EVIDENCED`:

- Satu jenjang, ditulis di kode (langkah 8): `KomiteList(1).KomiteID = "CHRISTINEANGELINA"`, `KomiteEmail = "christine_angelina@nusantarare.com"`, `IDKomite = "Claim Dept. Head"`, `Initial = "CA"`, `KomitePost = "Department Head Claim"`, `KomiteLoop = 1`.
- Tidak memanggil `FilterEmailKomiteWithLimit`, tidak menghitung `Flagkomite`, tidak menyentuh ambang apa pun.
- Jenis dibedakan satu nilai: `TempCommiteClaim.TypeComentAnalysis == "5"` → `IsReject`; selain itu → `IsCloseFile`. Nilai `"5"` adalah satu-satunya nilai yang pernah dibandingkan di seluruh ekspor Claim (5 kemunculan).
- Menyalin `ClaimData.{ListClaimAmount, CNPSpreadLoss, SpreadingRisk}` ke case anak, menulis `pyWorkPage.ClaimData.ClaimComitee = ChildWorkPage.KomiteList`, dan menyetel `CNPStatusCase = "COMITEE ACCEPTANCE (DEPT. HEAD)"`.

**Akibat yang ikut tertarik:** D-3 memperlakukan tiga jenis sebagai satu konsep, dan itu tetap benar untuk mekanika keputusan — tetapi **pembentukannya memang dua jalur berbeda**, dan hanya jalur ADJUSTMENT yang memakai roster. Dua dari tiga jenis sirkulasi hari ini tidak pernah menyentuh `EMAILKOMITE`, sehingga D-5 (roster dimiliki sistem baru) tidak dengan sendirinya mengatur mereka.

---

## F-21 · Gerbang `ProteksiSendKomiteCNP_Act`: menandai, tidak menghentikan

Rule ini menjalankan validasi lalu menandai galat lewat `Property-Set-Messages`; **tidak ada `Exit-Activity` di jalur galat mana pun**. `EVIDENCED`. Syarat yang diperiksa:

nomor polis · `Payable` · `PaymentType` · `NoAccount`/`NameOfBank`/`BranchOfBank` · `Currency` · `IDOfBank` · email profil operator · `Occupation` · kecocokan nomor rekening terhadap mata uang (dua arah, `Curr1`/`Curr2`, `NoAccount`/`NoAccount2`) · `NetClaim < 0` bila `DirectToKasir` · baris `TreatyName=="UR"` bila `PaymentType != "7"` · dan *"Adjustment value should not be more than estimation value"* per mata uang.

Penghentian sesungguhnya terjadi di tempat lain: `CreateChildKomiteCNP_Act` langkah 29 hanya memanggil `pxAddChildWork` bila `@hasMessages(myStepPage)` bernilai salah. `EVIDENCED`.

Rule ini juga memanggil `SetProtectionEstimation`, yang **tidak ada di repo** (`INVENTARIS-BUKTI.md` §2.3).

**Pagar dipatuhi:** sebagian syarat di atas menyentuh nilai yang berakhir di baris OS akseptasi dan muatan Kasir. Saya laporkan keberadaannya dan berhenti; tidak ada penilaian atas akibatnya di A-4 maupun A-5.

**Akibat yang ikut tertarik:** gerbang dan pembentukan adalah dua rule terpisah yang hanya terhubung lewat pesan pada halaman. Siapa pun yang memanggil `CreateChildKomiteCNP_Act` tanpa melewati gerbang akan membuat sirkulasi tanpa satu pun validasi.

---

## F-22 · Sirkulasi dapat lahir tanpa jenjang; satu-satunya pembatalan bukan soal jenjang

Bila penyaring roster tidak mengembalikan baris: `KomiteLoop = 0`, `KomiteList` kosong, dan `pxAddChildWork` **tetap dipanggil** — penjaganya hanya `@hasMessages`. Tidak ada pemeriksaan atas jumlah jenjang di titik mana pun. `EVIDENCED`.

Satu-satunya jalur yang membatalkan pembuatan adalah langkah 28:

```
JIKA @LengthOfPageList(.SpreadingRisk)=0 || @LengthOfPageList(ChildWorkPage.Adjustment.SpreadingRisk)=0  →  JUMP END
```

— yaitu ketika **alokasi layer kosong**, bukan ketika komitenya kosong. `EVIDENCED`.

Dilaporkan apa adanya, tanpa menilai keputusan beku no. 8.

**Akibat yang ikut tertarik:** G-17 kini punya jalur lahirnya, bukan hanya jalur matinya. Dan cacah case ber-`KOMITECOUNT=0` yang masih terdaftar sebagai permintaan kueri berubah maknanya: ia bukan lagi mencari tahu *apakah mungkin*, melainkan *seberapa sering*.

---

## F-23 · Perlakuan khusus satu operator, menyentuh empat tempat

Bila `pyWorkPage.pxCreateOperator == "VINCENTVERNANDO_1"`, `EVIDENCED`:

| Langkah | Akibat |
|---|---|
| 26.8 | seluruh perulangan roster **dilewati** |
| 26.9 / 26.10 | daftar diisi satu baris tetap: `KomiteID="VINCENTVERNANDO_1"`, email `klaim5@nusantarare.com`, `KomitePost="VC"`, `Initial="VC"`, dan `KomiteAproval = ""` — **untai kosong, bukan `0`** seperti anggota biasa |
| 26.15 | `KomiteLoop = 1` |

**Akibat yang ikut tertarik:** `KomiteRouter` memilih jenjang aktif dengan uji `.KomiteAproval == 0`; jalur ini menghasilkan `""`, bukan `0`. Apakah keduanya setara bagi pembanding Pega **tidak dapat ditentukan dari XML** — sama seperti F-12, batasnya saya nyatakan, tidak saya tafsirkan. Fakta yang berdiri: satu identitas orang mengubah bentuk data, bukan hanya alur.

---

## F-24 · Pemicu: syarat tampil terbaca, peran tidak

`.IsKomite = 1` disetel di `CreateChildKomiteCNP_Act` langkah 8, sebelum apa pun dikerjakan, dan disetel ulang di langkah 27. Di sisi Claim, `Section/AdjustmentDetailNP.xml` memakai `.IsKomite != 1` sebanyak **7 kali** sebagai syarat — `EVIDENCED` untuk cacahnya, `DIPULIHKAN` untuk maknanya sebagai syarat tampil kendali pengirim, karena pengikatan kendali ke kondisi tidak terbaca di ekspor (sebab yang sama dengan F-12).

**Tidak dapat ditentukan:** peran atau hak apa yang boleh memicu pembentukan sirkulasi. Keempat rule A-1b **tidak memuat satu pun `pyPrivilege`**. Wewenangnya, bila ada, hidup di harness atau di akses layar yang tidak ikut diekspor.

**Akibat yang ikut tertarik:** ADR-0006 (RBAC dirancang dari nol) berlaku penuh di sini — tidak ada warisan yang bisa ditiru, bahkan tidak ada yang bisa dibaca.

---

## Pertanyaan

Tiga, dan ketiganya hal yang benar-benar tidak terbaca di keempat rule. Dua kandidat lain saya bunuh sendiri: arti `.LIMIT_TOP` (tidak pernah menyaring — F-14, tidak berkonsekuensi) dan siapa yang berhak **menjadi** pengganti (terjawab F-14: siapa pun baris roster yang lolos penyaring yang sama).

```
Q-1 · Ambang berbasis nilai: dipulihkan atau tetap dua kelas tetap
  Aliran           : A-1b
  Pertanyaan       : Apakah sistem baru memulihkan ambang kewenangan berbasis nilai
                     adjustment, atau mempertahankan dua kelas tetap (0 dan 25.000.001)
                     yang berlaku hari ini?
  Mengapa perlu    : F-15 menunjukkan pembandingnya dimatikan dengan catatan "untuk
                     sementara" yang bertahan sampai versi terakhir. Tanpa jawaban ini,
                     aturan seleksi jenjang tidak dapat ditulis sama sekali — dan F-16
                     menunjukkan jalur nilai besar hari ini berakhir pada parameter
                     yang tidak terisi.
  Pilihan          : (a) Pulihkan — ambang dibandingkan terhadap nilai adjustment yang
                         diajukan. Konsekuensi: `.LIMIT_BOTTOM` dan `.LIMIT_TOP` pada
                         roster menjadi hidup; jumlah jenjang berubah mengikuti nilai;
                         seluruh sirkulasi berjalan hari ini memakai aturan berbeda.
                     (b) Pertahankan dua kelas tetap. Konsekuensi: `Flagkomite` ikut
                         dimigrasikan beserta empat konstantanya; cabang mati pada
                         langkah 14 diperbaiki (cacat ber-EVIDENCED, D-1), dan jalur
                         nilai besar harus diberi kelas yang tegas — hari ini ia tidak
                         punya.
  Uji V-6          : (a) tidak bersandar aliran beku — nilai adjustment dan RNMShare
                     tersedia di agregat Klaim, bukan di A-4/A-5. (b) tidak bersandar
                     aliran beku.
  Rekomendasi      : (a) — dua kelas tetap bukan kebijakan yang pernah diputuskan,
                     melainkan sisa penonaktifan sementara; memigrasikannya berarti
                     membekukan keadaan darurat orang lain menjadi aturan kita.
  Bila tak dijawab : (b) dengan jalur nilai besar diberi kelas tertinggi secara tegas —
                     paritas yang dapat dijalankan lebih aman daripada aturan baru yang
                     belum ada yang memilikinya.
  Bergantung pada  : —
```

```
Q-2 · Jalur tutup dan tolak klaim: satu jenjang tetap atau memakai roster
  Aliran           : A-1b
  Pertanyaan       : Apakah sirkulasi TUTUP_KLAIM dan TOLAK_KLAIM tetap satu jenjang,
                     atau memakai roster dan urutan DEGREE seperti jalur ADJUSTMENT?
  Mengapa perlu    : F-20 menunjukkan keduanya lahir dari rule lain yang tidak menyentuh
                     roster maupun ambang. D-5 menetapkan pemilik roster, bukan siapa
                     yang memakainya. Tanpa jawaban ini, dua dari tiga jenis sirkulasi
                     tidak punya aturan pembentukan.
  Pilihan          : (a) Tetap satu jenjang. Konsekuensi: jabatan penentunya diambil
                         dari roster, bukan dari nama orang (perbaikan cacat D-1 yang
                         sudah saya terapkan), tetapi kedalamannya tetap satu.
                     (b) Memakai roster dan urutan DEGREE seperti jalur ADJUSTMENT.
                         Konsekuensi: menutup klaim menjadi sama beratnya dengan
                         menyetujui pembayaran; jumlah jenjang mengikuti data roster.
  Uji V-6          : (a) dan (b) sama-sama tidak bersandar aliran beku — keduanya soal
                     pembentukan, bukan soal akibat pada OS, Kasir, atau isi cetak.
  Rekomendasi      : (a) — menutup klaim adalah keputusan yang membatalkan, bukan yang
                     mengeluarkan uang; menyamakan kedalamannya dengan jalur pembayaran
                     menambah beban tanpa ada yang pernah memintanya.
  Bila tak dijawab : (a).
  Bergantung pada  : —
```

```
Q-3 · Siapa berwenang mengalihkan pemegang jenjang
  Aliran           : A-1b
  Pertanyaan       : Peran apa yang berwenang memindahkan sebuah jenjang dari satu
                     pemegang ke pemegang lain pada sirkulasi yang sedang berjalan?
  Mengapa perlu    : E-3 menetapkan bentuknya (kedudukan, bukan orang; perpindahan =
                     peristiwa) dan menahan permukaannya sampai A-1b cair — A-1b kini
                     cair. F-19 menunjukkan pengalihan sudah terjadi hari ini sebagai
                     satu baris kode tanpa pelaku. Siapa yang berhak tidak terbaca di
                     rule mana pun, dan F-24 menunjukkan keempat rule A-1b tidak memuat
                     satu pun `pyPrivilege`.
  Pilihan          : (a) Peran administratif tersendiri, di luar keanggotaan komite.
                     (b) Pemegang jenjang itu sendiri, melimpahkan ke anggota roster
                         lain yang lolos penyaring yang sama.
                     (c) Pemilik proses klaim (peran yang membentuk sirkulasi).
  Uji V-6          : ketiganya tidak bersandar aliran beku — syarat kelayakan pengganti
                     sudah terjawab F-14, sehingga tidak satu pun cabang menunggu
                     aturan seleksi yang belum diketahui.
  Rekomendasi      : (a) — (b) membiarkan orang memilih penggantinya sendiri pada
                     keputusan yang justru dirancang berjenjang, dan (c) memberi
                     pembentuk sirkulasi kuasa atas siapa yang menilainya.
  Bila tak dijawab : (a), dengan aksi pengalihan tetap tidak diekspos sampai perannya
                     ada.
  Bergantung pada  : —
```
