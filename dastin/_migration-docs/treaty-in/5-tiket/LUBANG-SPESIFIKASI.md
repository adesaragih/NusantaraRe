# Lubang spesifikasi — Treaty In

**Tanggal:** 24 September 2026

> Setiap tempat yang membuat saya ingin kembali ke ekspor Pega untuk menulis sebuah tiket.
> **Dilaporkan, tidak ditambal.**

> ### BENTUK SETIAP BUTIR — ditetapkan 24 September 2026
>
> **Setiap lubang membawa dua hal: SIAPA YANG MENUTUPNYA, dan APA YANG MEMBUATNYA DITAGIH.**
> Lubang tanpa pemilik dan tanpa saat penagihan bukan lubang yang tercatat — ia lubang yang
> dilupakan dengan cara yang lebih rapi. `CONTEXT.md` §2.9b.
>
> Baris **Yang menagih** ditambahkan ke **kedelapan butir, termasuk yang sudah ditutup**, supaya
> bentuknya seragam ketika daftar ini ditiru modul berikutnya. Menambalnya dari ekspor akan menjadi keputusan tak bertanda yang
> tidak pernah ditinjau siapa pun.

**Sepuluh lubang, EMPAT sudah ditutup (L-1, L-2, L-6, L-7), satu DITUTUP SEBAGIAN (L-8), dan satu diberi pemilik (L-4).** Diperbarui 24 September 2026, langkah 4 sesi to-spec. **L-3 tidak ada di sini** — ia lubang **lingkungan**, bukan lubang spesifikasi; lihat
§8.

---

## L-1 — Folder `ddl/` tidak ada — **DITUTUP 24 September 2026**

> ### DITUTUP — `2-to-spec/ddl-usulan/` berdiri, 31 berkas
>
> | | |
> |---|---|
> | berkas tabel | **29** — satu per entitas penyerahan pertama |
> | `00_SKEMA_DAN_AKUN.sql` | skema `TREATY_MASUK`, akun `TREATY_MASUK_APP` |
> | `Z00_KUNCI_ALAMI.sql` | **20** constraint `UNIQUE`, masing-masing menyebut nomor invariannya |
> | kunci asing | **28** |
> | pengenal melebihi 30 bita (§16) | **0** — diperiksa, bukan diandaikan |
>
> **Dibangkitkan, bukan diketik.** `alat/buat-kamus-dan-ddl.py` mengurai `SPEC-MODEL-DATA.md` §10
> menjadi satu definisi, dan `alat/tulis-kamus-dan-ddl.py` menulis **`KAMUS-KOLOM.md` dan
> `ddl-usulan/` dari definisi yang sama** — keduanya tidak dapat berbeda.
>
> **Nol dijalankan**, dan itu tertulis di kepala setiap berkas, bukan disembunyikan.
>
> **Butir 3 urutan wewenang — belum dipulihkan.** `KEPUTUSAN-PEMBAGIAN-TIKET.md` §0.2 mencoretnya
> selama folder DDL tidak ada. Folder itu kini ada, **tetapi isinya USULAN yang belum ditinjau
> siapa pun**, dan anak tangga wewenang yang berisi usulan lebih berbahaya daripada anak tangga
> yang hilang. **Pemulihannya keputusan pemilik proses**, bukan akibat otomatis dari berdirinya
> folder ini.
>
> **Tiga hal yang sengaja TIDAK ada di dalamnya**, masing-masing sebagai pernyataan keputusan:
> bentuk fisik **paket uang** (16 atribut, 9 entitas), constraint tingkat pencatatan untuk **enam
> paket** yang menunggu `T-1`…`T-4`, dan kunci asing untuk **dua induk polimorfik**.


| | |
|---|---|
| **Apa** | sesi DDL ditangguhkan; folder `ddl/` tidak pernah lahir |
| **Akibat** | nama constraint, angka presisi INV-45, bentuk fisik `POTONGAN` berinduk dua, dan pilihan menurunkan INV-35 menjadi constraint — seluruhnya belum ada |
| **Perlakuan** | butir 3 urutan wewenang **DICORET** selama folder itu belum ada (`KEPUTUSAN-PEMBAGIAN-TIKET.md` §0.2, §4.3) |
| **Siapa menutup** | sesi DDL, **setelah** §10 selesai untuk 23 entitas |
| **Yang menagih** | **§10 sudah selesai 24 Sep 2026** — prasyaratnya jatuh, dan sesi DDL kini dapat dimulai. Penagihnya: gerbang pembuka sesi DDL, yang tidak dapat dijawab tanpa folder ini |

---

## L-2 — `SPEC-MODEL-DATA.md` §10 baru meliputi 4 entitas — **DITUTUP 24 September 2026**

| | |
|---|---|
| **Apa** | §10 memuat nama, tipe, dan keterisian untuk `KONTRAK`, `VERSI_KONTRAK`, `LAYER`, `DETAIL_PROPORSIONAL` saja. §11 menyatakannya sendiri |
| **Akibat** | kriteria selesai yang bisa gagal hampir selalu menyebut nama field; untuk entitas tanpa §10 ia tidak dapat ditulis |
| **Besarnya** | **23 entitas** menunggu §10 — 17 milik modul + 6 tabel acuan. Angka ini baru diketahui **setelah** `DAFTAR-PEKERJAAN.md` disusun; lihat §5.3 di sana |
| **Perlakuan** | jalan (iii) — daftar pekerjaan dulu, lalu §10 **hanya untuk entitas yang terpakai** |
| **Siapa menutup** | sesi to-spec |
| **Yang menagih** | **sudah ditagih dan ditutup** — setiap tiket menuntut kriteria selesai yang menyebut nama field, dan nama field hanya ada di §10 |

> **Ini lubang terbesar di seluruh alur**, dan satu-satunya yang menahan seluruh pembentukan tiket.

> ### DITUTUP — §10 kini meliputi 27 entitas
>
> §10.5–§10.22 ditulis 24 September 2026 untuk **23 entitas**, dan §10.3–§10.4 **diselesaikan ulang
> atas sumber yang lengkap** sesudah L-8. Kedua kewajiban §5.0b berlaku pada seluruhnya: setiap
> entitas menyatakan kunci alaminya dan lingkupnya, dan `SUMBU_REKONSILIASI`-nya — **termasuk ketika
> ia tidak punya**, karena itu keputusan dan bukan kekosongan.
>
> **Yang tersisa bukan lagi entitas tanpa §10.** `SPEC-MODEL-DATA.md` §11.2 mendaftar sisanya:
> empat keputusan yang mengubah jumlah entitas (§10.23c), tingkat pencatatan 14 paket uang,
> penomoran lima kunci alami, dan empat kepanjangan singkatan yang tinggal ditanyakan ke orang.

---

## L-4 — Tidak ada spesifikasi layar di mana pun

| | |
|---|---|
| **Apa** | `SPEC-INVARIAN.md` **nol** menyebut layar; `4-erd-dan-tabel-datar/ERD.md` **nol** |
| **Akibat** | bentuk layar tidak dapat menjadi kriteria selesai, dan pembagian tiket **per layar** mustahil — ia menuntut mengarang layarnya lebih dulu |
| **Perlakuan** | kriteria selesai dinyatakan sebagai **perilaku yang dapat gagal**, bukan susunan tampilan. Pembagian tiket memakai **kemampuan**, bukan layar (G1) |
| **Siapa menutup** | **DITETAPKAN 24 September 2026 (K-6): batch LAPISAN APLIKASI.** Bukan "kemungkinan tidak perlu ditutup" |
| **Pernyataan lingkup** | spesifikasi layar **tidak dibuat di batch lapisan data**, dan itu keputusan, bukan kelalaian. Rumusan lama — *"kemungkinan tidak perlu ditutup sama sekali"* — **dicabut**: ia membuat lubang yang punya pemilik terbaca seperti lubang yang boleh dilupakan |
| **Yang menagih** | tiket pertama yang kriteria selesainya **tidak dapat ditulis tanpa menyebut bentuk layar**. Selama tidak ada satu pun, lubang ini memang tidak perlu ditutup |

---

## L-5 — Huruf `G` bermakna tiga hal

| | |
|---|---|
| **Apa** | `G2`/`G3` berarti **gelombang** di `SPEC-MODEL-DATA.md` §2.3, **gerbang sambungan Adjustment** di `4-erd-dan-tabel-datar/KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md`, dan **gerbang pembagian tiket** di `KEPUTUSAN-PEMBAGIAN-TIKET.md` |
| **Akibat** | salah baca menghasilkan ruang lingkup yang salah |
| **Perlakuan** | **SUDAH DITUTUP.** Di seluruh `5-tiket/` gelombang ditulis `GEL-2`/`GEL-3`, dan satu baris koreksi ditulis ke `SPEC-MODEL-DATA.md` §2.3 **sendiri** — bukan hanya di berkas gerbang, karena pembaca §2.3 tidak akan membuka berkas gerbang untuk tahu |
| **Siapa menutup** | **SUDAH DITUTUP** — sesi tiket, 24 September 2026 |
| **Yang menagih** | koreksinya ditulis **di dalam §2.3 sendiri**, jadi penagihnya adalah pembacaan §2.3 itu sendiri. Itu sebabnya ia tidak cukup ditulis di berkas gerbang |

---

## L-6 — Daftar entitas §2.3 tidak lengkap

| | |
|---|---|
| **Apa** | `SPEC-MODEL-DATA.md` §2.3 tidak memuat `BATAS_PER_BAHAYA` (diperkenalkan §10.0b, ditegakkan INV-14) maupun **enam tabel acuan** — `MATA_UANG`, `JENIS_POTONGAN`, `JENIS_REASURANSI`, `BAHAYA`, `KELOMPOK_TREATY`, `KELAS_BISNIS` |
| **Akibat** | siapa pun yang menghitung ruang lingkup dari §2.3 **kekurangan tujuh entitas** |
| **Perlakuan** | butir 4 ditambahkan ke urutan wewenang: soal **daftar entitas**, `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` **mengikat**; §2.3 tidak dipakai untuk menghitung |
| **Sapuan** | seluruh berkas `.md` modul disisir untuk klaim jumlah entitas. **Tidak ada satu pun** selain angka "20" yang saya tulis sendiri di berkas gerbang pada 23 September dan sudah dikoreksi. §11 menyebut "entitas selebihnya" tanpa angka. **Kekurangan tujuh itu tidak pernah merambat** |
| **Siapa menutup** | **SUDAH DITUTUP 24 September 2026** — §2.3 dilengkapi dengan ketujuh entitas plus `PEMULIHAN_LIMIT` (kedelapan). **Butir 4 urutan wewenang TETAP BERLAKU**, dipertahankan di samping perbaikannya karena §2.3 dapat tertinggal lagi dan butir 4 tidak |
| **Sisa yang lahir dari penutupan ini** | `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` — sumber yang **mengikat** menurut butir 4 — **belum memuat `PEMULIHAN_LIMIT`**, sehingga wewenang tertinggi kini menyebut 27 sementara §10.23c memutuskan 28. Dilaporkan di `SPEC-MODEL-DATA.md` §1a, **tidak ditambal**; pemiliknya alur B |
| **Yang menagih** | setiap penghitungan ruang lingkup. Butir 4 urutan wewenang **memaksa** pembacanya memakai `STRUKTUR-DATA.md`, jadi penagihnya sudah terpasang di jalan yang wajib dilewati |

---

## L-7 — ADR-0055 tidak punya cara membuang draf — **DITUTUP 24 September 2026**

> **Ditemukan, dilaporkan, dan diputuskan pemilik proses pada hari yang sama.** ADR-0055 kini punya
> keadaan kedelapan `DIBATALKAN`, terminal, dengan perpindahan `DRAFT → BATALKAN → DIBATALKAN` dan
> tiga syarat mengikat. Kemampuannya **P-57**. Uraian di bawah dipertahankan sebagai alasannya.

**Lubang paling berakibat di daftar ini**, dan satu-satunya yang mengunci pengguna.

| | |
|---|---|
| **Apa** | dua belas perpindahan ADR-0055 disisir: **`DRAFT` punya tepat satu perpindahan keluar — `AJUKAN`.** Tidak ada `BATALKAN`, `BUANG`, `TARIK`, maupun `HAPUS` |
| **Akibat** | digabung INV-25 (satu versi belum-selesai per kontrak) dan INV-22 (perpindahan di luar daftar ditolak): **pengisi kontrak yang salah pencet tidak punya jalan keluar, dan kontraknya terkunci.** Satu-satunya jalan tersisa adalah mengajukannya supaya ada yang menolaknya |
| **Sistem lama** | **juga tidak punya "buang draf".** `TreatyInDeclineConfirmation_postact` (4 langkah, seluruhnya **hidup**) hanya menyetel keadaan. `TreatyInDeclineConfirmation_postactEDM` (8 langkah, 4 hidup) — langkah 6 dan 7 *"RDB remove EDM"* **HIDUP**, **menghapus baris** di `M_TREATY_IN_EDM` dan `TREATY_IN_EDM`. Jadi orang memang memakai jalur penolakan sebagai tempat sampah, dan untuk addendum buktinya lenyap |
| **Kenapa mendesak** | di sistem lama jalan keluarnya ada meski buruk. Di rancangan baru jalan itu **hilang** dan tidak ada penggantinya: `DITOLAK` terminal, ADR-0045 melarang jejak hilang, INV-25 melarang versi kedua |
| **Perlakuan** | dilaporkan sebagai lubang, **tidak ditambal sendiri**; diputuskan pemilik proses. Kemampuannya kini **P-57**, bergolongan **BARU**. Uraian lengkapnya di `DAFTAR-PEKERJAAN.md` §2.2 |
| **Siapa menutup** | **SUDAH DITUTUP** — pemilik proses, ADR-0055 perubahan 24 September 2026 |
| **Yang menagih** | penanda di `SPEC-INVARIAN.md` INV-20, INV-23, INV-25, §3 dan `SPEC-MODEL-DATA.md` §10.2 — **keadaan adalah atribut yang §10.2-nya sudah ditulis**, sehingga sesi to-spec tidak dapat melewatinya |

Tiga bentuk yang ditimbang, dan yang dipilih:

| Pilihan | Yang perlu ditimbang |
|---|---|
| **✓ DIPILIH** — perpindahan `DRAFT → DIBATALKAN`, keadaan terminal baru | menambah keadaan kedelapan; jejaknya utuh, ADR-0045 aman |
| perpindahan `DRAFT → DITOLAK` langsung oleh pengisinya | tidak menambah keadaan, tetapi `DITOLAK` lalu berarti dua hal berbeda |
| longgarkan INV-25 | mengubah invarian, bukan menambah perpindahan |

> **Diketahui sebelum sesi to-spec melanjutkan §10, bukan sesudahnya** — karena pilihan yang diambil
> menambah keadaan, dan keadaan adalah atribut `VERSI_KONTRAK` yang §10.2-nya **sudah** ditulis.
> Penandanya sudah dipasang di `SPEC-INVARIAN.md` (INV-20, INV-23, INV-25, §3) dan
> `SPEC-MODEL-DATA.md` §10.2, supaya sesi to-spec tidak dapat melewatkannya.
>
> **Satu syarat ternyata gratis.** "Nomor revisi tidak dipakai ulang" tidak menuntut mekanisme
> baru: INV-04 sudah membuat `NOMOR_URUT_VERSI` unik per kontrak, jadi selama barisnya tidak
> dihapus, constraint yang sudah ada menolak pemakaian ulangnya sendiri.

---

## L-8 — Pohon 985 simpul punya titik buta 414 properti — **DITUTUP SEBAGIAN 24 September 2026**

> ### PERNYATAAN LINGKUP — K-2, arahan pemilik proses
>
> **Beres untuk 27 entitas gelombang 1. Tiga ratus empat puluh properti sisanya TETAP TERBUKA, dan
> pemiliknya gelombang 2.**
>
> | | Jumlah |
> |---|---:|
> | properti tidak pernah terlihat perkakas mana pun | **414** |
> | dibatasi pada 14 kelas penyandang entitas penyerahan pertama | **74** |
> | di antaranya yang menuntut keputusan, dan **sudah diadili di §10** | **41** |
> | **sisanya — belum pernah diadili siapa pun** | **340** |
>
> **Lubang ini TIDAK berbunyi "ditutup" tanpa kata "sebagian".** Yang beres adalah bagian yang
> menyentuh penyerahan pertama; 340 properti selebihnya milik kelas yang entitasnya belum
> dimodelkan, dan **tidak ada satu pun langkah yang pernah membukanya**.
>
> | | |
> |---|---|
> | **Siapa menutup sisanya** | sesi to-spec **gelombang 2** |
> | **Yang menagih** | `buat-pohon-treatyin.py` mencetak *"DITOLAK penjaga primary_ok: 584 rujukan jalur"* setiap kali dijalankan — penagih yang tidak bergantung pada ingatan siapa pun |
> | **Label yang DILARANG selama ini terbuka** | **"tidak dipakai"** (K-3). Selama 340 properti, `L-9`, dan `L-10` masih terbuka, satu-satunya label yang sah adalah **`DIPAKAI`**, **`NOL DI EKSPOR INI`**, dan **`BELUM TERPERIKSA`** |


| | |
|---|---|
| **Apa** | penjaga `primary_ok` di `alat/buat-pohon-treatyin.py` **membuang** properti yang ditulis aturan ber-*applies-to* **kelas anak**, alih-alih menggantungkannya di bawah kelas anaknya. Aturan semacam itu tidak menyumbang satu simpul pun |
| **Bagaimana terlihat** | `BreakDownSprdList` — sumber `RINCIAN_PENYEBARAN` menurut `SPEC-MODEL-DATA.md` §3.5 — muncul **54 kali** di ekspor dan **nol kali** di pohon |
| **Besarnya** | **414 properti** tidak pernah terlihat perkakas mana pun. Dibatasi pada 14 kelas penyandang 27 entitas penyerahan pertama: **74**, yang **41** di antaranya menuntut keputusan |
| **Akibat** | `SPEC-MODEL-DATA.md` **§10.3 (`LAYER`) dan §10.4 (`DETAIL_PROPORSIONAL`) dibuka kembali** — keduanya sudah dinyatakan selesai. `peta-nama-tabel-treatyin.tsv` dan `ERD-STRUKTUR-TREATYIN.html` tidak punya tabel untuk `BreakDownSprdList` |
| **Perlakuan** | perkakas mandiri `alat/sapu-properti-per-kelas.py` ditulis dan dijalankan; hasilnya terukur, bukan dugaan. Uraian lengkap `4-erd-dan-tabel-datar/TITIK-BUTA-POHON.md`. Ke-41 properti diadili **di §10 entitasnya masing-masing** |
| **Siapa menutup** | sesi to-spec, **di dalam pekerjaan §10 yang sedang berjalan** — bukan pekerjaan terpisah |
| **Yang menagih** | **perkakasnya sendiri**: `buat-pohon-treatyin.py` kini mencetak *"DITOLAK penjaga primary_ok: 584 rujukan jalur"* setiap kali dijalankan. Penagih yang tidak bergantung pada ingatan siapa pun |

> **Penjaganya tidak dicabut.** Mencabutnya mengembalikan simpul palsu ke tingkat akar, dan simpul
> palsu **terbaca sebagai fakta** sementara simpul hilang setidaknya tidak mengarang apa pun. Yang
> diperbaiki adalah perlakuan atas yang ditolak penjaga, lewat sumber kedua yang mandiri.

---

## L-9 — Ekspor ini diambil dari lingkungan apa

| | |
|---|---|
| **Apa** | `SystemSettings/LinkService.xml` bernilai **`3`**, dan keterangannya sendiri memetakan **3 = Quality assurance** (1 Sandbox · 2 Development · 4 Staging · 5 Production) |
| **Kenapa ia besar** | bila ekspor ini ruleset **QA**, maka sebelas sesi membaca aturan yang **boleh berbeda** dari produksi — dan ruleset QA memang ada supaya boleh berbeda. Setiap langkah mati, setiap `NEVER &&`, tetapan **15% ORS**, dan kedua nomor kontrak berkondisi: seluruhnya mungkin hanya ada di QA |
| **Bacaan kedua** | ini ruleset **produksi**, dan *link service*-nya menunjuk QA — yang berarti produksi memanggil lingkungan QA, dan **itu cacat tersendiri** |
| **Uji yang dijalankan** | sapuan `pxCreateSystemID` atas 708 berkas. **Ia menyempitkan, tidak memutuskan** — penanda itu mencatat tempat aturan **dibuat**, bukan tempat ekspor diambil. Hasilnya di `4-erd-dan-tabel-datar/PRA-PEMISAHAN-ANGKA-DAN-TABRAKAN.md` §1b |
| **Yang TIDAK terancam** | **rancangan sistem baru.** Yang dirancang adalah sistem baru; yang terancam **penggolongan PELESTARIAN versus PERUBAHAN** dan **daftar temuan** |
| **Siapa menutup** | **kantor** — pertanyaan ke orang, bukan ke berkas |
| **Yang menagih** | **medan `DASAR:` bergolongan `EVIDENCED` pada setiap tiket**, dan ia dibuat **dapat disapu**: setiap rujukan bukti menyebut **ekspor mana** yang menjadi sumbernya — `EVIDENCED(SetSpreadName@ekspor-2026-09)`. Tanpa penanda itu, jawaban "QA" menuntut membaca ulang **setiap** tiket; dengan penanda itu, **satu sapuan** mengeluarkan daftar terdampak. Penagih yang tidak dapat dijalankan mekanis adalah penagih yang bergantung pada ingatan |

> **Pertanyaannya satu kalimat dan jawabannya tiga puluh detik: EKSPOR INI DIAMBIL DARI LINGKUNGAN
> MANA, DAN OLEH SIAPA.** Ia naik ke **puncak** daftar pertanyaan ke orang — di atas slip
> non-proporsional dan di atas batas wewenang.

---

## L-10 — Dua puluh satu JENIS aturan tidak ada di ekspor

| | |
|---|---|
| **Apa** | ekspor berfolder menurut jenis aturan. **Dua belas jenis muncul; dua puluh satu jenis yang diketahui dipakai Pega tidak muncul sama sekali** — di antaranya `Rule-Obj-Flow`, `Rule-Obj-CaseType`, `Rule-Declare-Expressions`, `Rule-Declare-Trigger`, `Rule-Access-Role-Obj` |
| **Kenapa baru ketahuan** | *"tidak ada di ekspor"* terbaca sebagai *"tidak ada di sistem"* selama sebelas sesi. Ketiadaan satu **jenis penuh** lebih mudah diperhatikan daripada ketiadaan 414 properti (L-8) — **tetapi hanya bila ada yang bertanya** |
| **Yang paling berat** | `Declare Expression` dan `Declare Trigger` **menulis nilai tanpa dipanggil siapa pun**. Sapuan penulis kita mencari `Property-Set`; bentuk itu **tidak akan pernah melihat keduanya** |
| **Klaim yang terdampak** | **lima**, didaftar satu per satu di `4-erd-dan-tabel-datar/JENIS-ATURAN-TAK-TEREKSPOR.md` §3.1. Yang paling patut dicurigai **§14.5 `SPEC-MODEL-DATA.md`**, karena di situ ketiadaan penulis dipakai sebagai **alasan membuang dua atribut** |
| **Akibat lanjutan** | ADR-0055 tiga belas perpindahan adalah perpindahan **yang terlihat dari sisi aktivitas**; dan **P-41 tertahan karena jenis aturannya memang tidak diminta**, bukan karena belum ditemukan |
| **Siapa menutup** | **kantor** — satu pertanyaan, ditumpangkan pada permintaan ekspor kedua |
| **Yang menagih** | kelima klaim di §3.1 masing-masing **menyebut L-10 di tempatnya sendiri**, sehingga pembacanya tertagih tanpa membuka berkas ini. Dan **P-41** tidak dapat lepas dari TERTAHAN sampai jawabannya datang |

> **Permintaannya tidak meminta isinya, hanya DAFTARNYA** — dan itu yang membuatnya murah untuk
> dipenuhi. Jawaban **nol** membuat kesimpulan lama berdiri **di atas jawaban alih-alih di atas
> ketiadaan**; jawaban **bukan nol** memberi daftar yang cakupannya langsung terukur.

---

## 9. Keempat penahan LUAR — menjadi **tiket tertahan**, bukan penghenti

**Ditetapkan 24 September 2026.** Empat hal menunggu kantor dan **tidak akan terjawab di dalam sesi
mana pun**. Selama ini keempatnya berperilaku seperti penghenti; mulai sekarang keempatnya
**tiket tertahan bernomor menurut G3**.

| Penahan | Nomor lubang | Siapa dapat menjawab |
|---|---|---|
| instans Oracle yang terjangkau | **L-3** | kantor — infrastruktur |
| izin baca-saja ke data produksi | — | kantor — pemilik data |
| ekspor kedua: 15 berkas **dan** daftar jenis aturan | **L-9, L-10** | kantor — administrator Pega |
| dokumen batas wewenang · pemilik tabel kapasitas | eskalasi butir 5, 7 | manajemen |

> **Aturan yang berlaku untuk seluruh pekerjaan berikutnya:** bila sebuah langkah menyentuh salah
> satu dari keempatnya, **langkah itu tetap jalan**, dan bagian yang menunggu **dipecah keluar**
> menjadi tiket tertahan tersendiri — bernomor, muncul di peta liputan, **tanpa taksiran**, dengan
> medan `PENGHALANG` yang menyebut **apa · siapa · apa yang berubah bila jawabannya datang**.
>
> Yang dilarang adalah membiarkan salah satunya **menghentikan** sebuah blok. Penahan yang
> menghentikan pekerjaan berperilaku sama dengan pekerjaan yang tidak pernah dimulai — dan ia tidak
> terlihat di papan mana pun.

---

## 8. Yang BUKAN lubang spesifikasi

### L-3 — empat invarian belum terbukti: LUBANG LINGKUNGAN

INV-47, INV-50, INV-51, dan bentuk lintas baris INV-31 ditegakkan lewat *materialized view*
ber-`REFRESH ON COMMIT` dengan `CHECK` padanya. Uji negatifnya **belum pernah dijalankan**, karena
**tidak ada instans Oracle yang terjangkau**.

**Tidak ada sesi mana pun yang dapat menutupnya** — bukan to-spec, bukan DDL, bukan sesi tiket.
Karena itu ia tidak ditaruh di sini: menaruhnya di daftar lubang spesifikasi membuatnya tampak
seakan ada sesi yang bisa mengerjakannya.

**Tempatnya: daftar yang harus disediakan kantor.**

| Yang diminta | Untuk |
|---|---|
| satu instans Oracle yang terjangkau, boleh salinan | menjalankan `UJI-NEGATIF-INVARIAN.md` N-0…N-5f |
| **EKSPOR KEDUA — 15 berkas aturan dari lingkungan PRODUKSI** | menutup **L-9**. Daftarnya di bawah, dan ia **dua persen dari 708** |
| izin kueri **baca-saja** ke basis data produksi | Uji A…**AD**, termasuk Uji X-2 yang menahan P-20, Uji AC yang mengukur addendum yang hilang, dan **Uji AD** yang memutuskan sumbu mata uang penyebaran |

| | |
|---|---|
| **Siapa menutup** | **kantor** — tidak ada sesi mana pun yang dapat |
| **Yang menagih** | **empat invarian tetap berstatus BELUM DIBUKTIKAN di tabel §1 `SPEC-INVARIAN.md`**, yaitu tabel yang dibaca orang yang tidak membaca sisanya. Selama baris itu ada, penagihnya berbunyi setiap kali seseorang membuka berkas invarian |

Selama keduanya belum ada, **P-18 tetap TERTAHAN** dan keempat invarian itu tetap **klaim, bukan
fakta**.

### Permintaan ketiga: ekspor kedua, dan ia sengaja TIDAK LENGKAP

**L-9 tidak dapat dijawab dari ekspor yang kita punya** — `pxCreateSystemID` mencatat tempat aturan
**dibuat**, bukan tempat ekspor **diambil**. Maka pencariannya dihentikan di sana, dan yang diminta
sumber kedua.

> **Seluruh temuan besar modul ini bersandar pada lima belas berkas.** Bila kelimabelasnya
> **byte-identik** antara ekspor yang kita punya dan ekspor baru dari produksi, maka pertanyaan
> lingkungan **tidak lagi mengancam temuan mana pun** — terlepas dari jawabannya.

**Lima belas dari 708 adalah dua persen.** Permintaan dua persen jauh lebih mungkin dipenuhi
daripada permintaan seluruhnya, dan itulah sebabnya ia berbentuk begini.

| Berkas | Temuan yang bersandar padanya |
|---|---|
| `Activity/SaveTreatyIn_Act.xml` · `Activity/SaveTreatyIn_EDM_Act.xml` | langkah mati; addendum tidak sampai ke hilir |
| `Activity/TreatyInDeclineConfirmation_postact.xml` · `…postactEDM.xml` | penghapusan baris addendum — L-7, Uji AC, eskalasi butir 4 |
| `Activity/TreatyInCheckID.xml` | daftar K1, enam syarat yang mati |
| `Activity/FetchQSfromMasterXOL.xml` | **tetapan 15% ORS**, dan **dua nomor kontrak berkondisi** |
| `Activity/SetSpreadName.xml` | `ReinsID` diisi dari `ReinsTypeID`; INV-47 dan INV-50 |
| `Activity/SetSpreadingXOL.xml` | `BreakDownSprdListXOL` — L-8 |
| `Activity/TreatyInSetValueInstallment.xml` | `Installment[]` per mata uang; INV-12 |
| `Activity/SetReinstatementPct.xml` | tetapan `100`/`100`; entitas `PEMULIHAN_LIMIT`; P-59 |
| `DataTransform/TreatyInForceResolveComplete.xml` · `TreatyInForceEdit.xml` | persetujuan tanpa penyetuju |
| `Activity/SaveTreatyInDetail_Act.xml` | penulisan nol ke tabel datar |
| `Activity/TreatyInNonAddItem.xml` · `Activity/TreatyInSetBrokerage.xml` | bukti 5a — cabang penulis `Limits[]` dan `Share[]` |
| `SystemSettings/LinkService.xml` | **penanda lingkungan itu sendiri** |

**Perbandingannya mekanis:** hash per berkas, lalu beda pada yang berbeda. **Ketiga hasilnya
berguna:**

| Hasil | Artinya |
|---|---|
| seluruhnya identik | seluruh temuan besar **aman**, dan L-9 turun dari **ancaman** menjadi **catatan** |
| sebagian berbeda | yang berbeda itulah **daftar temuan yang harus dibaca ulang**, dan cakupannya **terukur** |
| ekspor kedua tidak dapat diperoleh | itu **jawaban juga**, dan ia naik ke eskalasi sebagai **ketidakmampuan memverifikasi** — bukan sebagai pertanyaan yang belum sempat ditanyakan |

> **Permintaannya bukan mustahil, ia belum diajukan.** Contoh data yang disediakan pemilik proses
> membawa `pxCreateSystemID: "pegaprdnusare"` pada baris tahun 2023 dan 2024 — **sistem produksinya
> ada dan terjangkau orang.**

### Dan satu pertanyaan ditumpangkan pada permintaan yang sama — **L-10**

> **ADAKAH ATURAN BERJENIS `Rule-Obj-Flow`, `Rule-Obj-CaseType`, `Rule-Declare-Expressions`,
> `Rule-Declare-Trigger`, DAN `Rule-Access-Role-Obj` UNTUK APLIKASI INI? BILA ADA, BERAPA BANYAK
> DAN NAMANYA APA?**
>
> **Tidak perlu isinya. Cukup daftarnya** — dan itu membuatnya jauh lebih murah daripada lima belas
> berkas di atas. Ia dapat dijawab dari layar pencarian aturan dalam satu menit.
