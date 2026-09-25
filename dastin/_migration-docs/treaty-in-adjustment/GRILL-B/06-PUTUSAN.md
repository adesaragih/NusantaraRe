> Modul  : Treaty In Adjustment · Ronde B · 2026-09-25
> Peran  : penilai
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 06 · PUTUSAN RONDE B

## GRL-09 — Pengenal tampilan: diturunkan untuk versi baru, dilestarikan apa adanya untuk warisan

* **Cabang** — B (identitas dan penomoran), `QB-1`.
* **Label — dipisah, karena dua bagiannya punya risiko berbeda** (`METODE` §6.5):

  | Bagian | Label | Berubah dari |
  |---|---|---|
  | pengenal warisan dibawa apa adanya, beserta asal-usulnya | **PELESTARIAN** | — |
  | bentuk turunan **tanpa lebar tetap dan tanpa anggapan panjang pengenal induk** | **PERUBAHAN** | offset tetap `@substring(ID,10,12)`, anggapan pengenal kontrak selalu tujuh karakter, dan lebar nomor dua digit yang patah melewati `R09` (**TDA-12**, NA-02) |

* **Keputusan**

  1. **Satu sumber kebenaran.** Untuk versi **baru**, pengenal tampilan **diturunkan saat
     ditampilkan** dan **tidak disimpan**. Yang disimpan hanyalah **nilai warisan**, sebagai
     atribut warisan beserta asal-usulnya (ADR-0042).

     Bila to-spec kelak memerlukan penyimpanannya untuk versi baru — misalnya demi pencarian atau
     antarmuka ke luar — itu **harus beralasan** dan dijaga invarian dengan **uji negatif yang
     dapat gagal** (`METODE` §2.2, §5.3). Menyimpan turunan tanpa penjaga adalah bentuk "satu fakta
     dua tempat" yang ADR-0041 larang.
  2. **Pengenal warisan dilestarikan apa adanya** — termasuk yang berbentuk `/R010`, yang melompat,
     dan yang berulang. Ia **tidak diturunkan ulang**.

     Preseden **GRL-06** berlaku di sini, dan alasannya disebut (`METODE` §3.6): GRL-06 melarang
     menghitung ulang nilai warisan yang merupakan **fakta yang dilihat orang**. Pengenal
     `‹kontrak›/Rnn` sudah beredar di luar sistem; menurunkannya ulang berarti mengganti nomor yang
     sudah dipegang orang lain dengan nomor yang "seharusnya" — rapi, konsisten, dan tidak dapat
     dipulihkan. Ini **berbeda** dari kelima penanda lama (GRL-08), yang **mekanisme** dan tidak
     pernah dilihat siapa pun.
  3. **Awalan pengenal diambil dari pengenal kontrak sebagaimana ditetapkan ADR-0040 butir 1** —
     "addendum atas" adalah **hubungan versi**, sehingga awalannya adalah identitas kontrak
     induknya, apa pun panjangnya. **Tidak ada anggapan tujuh karakter.**
  4. **Melestarikan nomor tidak memulihkan isinya**, dan itu perkara tersendiri. Karena **TDA-01**,
     sebuah dokumen yang beredar dengan nomor tertentu dapat merujuk isi yang **sudah tertimpa**.
     Dititipkan ke **cabang I**.

* **Bahan to-spec yang lahir bersama putusan ini**

  | # | Bunyi | Berubah dari |
  |---|---|---|
  | **B-1** | Menyimpan versi dengan **nomor urut** yang sudah dipakai di dalam kontrak yang sama **ditolak**, dan pesannya **menyebut nomor yang bentrok** | tabrakan masuk cabang `UPDATE`, menimpa baris yang ada, lalu melapor berhasil (TDA-01) |
  | **B-2** | Versi baru yang **pengenal tampilannya** sama dengan pengenal warisan mana pun di kontrak yang sama **ditolak** | tidak ada padanannya — di sistem lama tabrakan pengenal justru menjadi `UPDATE` |

  B-2 tidak tergantikan oleh B-1: B-1 menjaga **nomor urut**, B-2 menjaga **pengenal tampilan**.
  Keduanya dapat berbeda justru karena pengenal warisan tidak selalu cocok dengan urutannya.

* **Titipan ke cabang I**

  | Isi |
  |---|
  | Bagaimana `NOMOR_URUT_VERSI` diberikan kepada baris warisan, dan **dari nomor berapa** penomoran baru dimulai — mengingat deret warisan dapat bolong, berulang, atau berbentuk `/R010` |
  | Nasib dokumen yang beredar dengan nomor yang isinya sudah tertimpa (TDA-01) |

* **Apakah penimpaan meninggalkan jejak yang dapat diukur — diperiksa, dan jawabannya sebagian.**

  | Kandidat jejak | Dapat dipakai? |
  |---|---|
  | `RevisionDate` | **ya, sebagian** — `TreatyInEDMSetValue` langkah 3 menulisnya `@CurrentDateTime()` pada setiap pembuatan. Baris yang `RevisionDate`-nya **lebih baru daripada** `RevisionDate` baris bernomor lebih tinggi di kontrak yang sama adalah tanda penimpaan |
  | `CommentList` | **ya, sebagian** — langkah 3 yang sama **mengosongkannya** (`CommentList = ""`). Addendum yang sudah disetujui seharusnya membawa jejak persetujuan; yang kosong padahal statusnya bukan draf adalah tanda |
  | `OLDID` | **tidak** — ia hanya menunjuk induk, dan penimpaan tidak mengubahnya |

  Maka ukurannya **ada**, dan ditulis sebagai **`UA-14`**. Yang **tidak** terukur: isi yang hilang
  itu sendiri — ia tertimpa tanpa salinan, dan itu dinyatakan apa adanya.

* **Ditolak** — **(b)** sepenuhnya diturunkan: pengenal warisan tidak selalu dapat diturunkan ulang
  (TDA-12, NA-02), dan menurunkannya berarti mengganti nomor yang sudah beredar. **(c)** sepenuhnya
  dilestarikan sebagai teks: membawa masuk kembali anggapan lebar tetap dan offset, yaitu cacat
  yang justru sedang diperbaiki.
* **Arah dampak** — memilih (a) padahal pengenal itu tak pernah dilihat siapa pun: satu kolom teks
  dibawa tanpa pembaca, murah dan dapat dibuang. Memilih (b) padahal ia dipakai di luar: nomor yang
  sudah dikutip tidak lagi dapat ditemukan, dan ketahuannya saat seseorang mencari dokumen lama —
  mahal dan tidak dapat dibatalkan. Ongkosnya tidak setangkup; sisi murah diambil tanpa menunggu
  `DB-11`.
* **Syarat pembalikan** — `DB-11`. Bila pembaca di luar sistem dibenarkan tidak ada, butir 2
  menyusut menjadi "simpan hanya untuk baris yang benar-benar menyimpang", dan model berkurang satu
  atribut.
* **Status** — **terkunci**, dengan `DB-11` dan `UA-14` terbuka.

---

## MA-09 · Bukti dibangun dari ingatan atas dua temuan sendiri, tanpa membacanya ulang

Bukti B2 saya memuat dua kekeliruan, dan **keduanya bertentangan dengan temuan saya sendiri yang
sudah tertulis**:

1. Saya menulis `HASIL1` dipakai sebagai "nomor tertinggi + 1". `PENGETAHUAN.md` §4.4 butir 1 dan
   **`MA-03`** sudah menetapkan bahwa nilainya **dibuang** dan hanya diuji kosong.
2. Saya menyimpulkan `HASIL1` kosong bila yang dipilih sebuah addendum — padahal baris itu memenuhi
   `LIKE`-nya sendiri, sebab prosedur simpan menulis ke `M_TREATY_IN_EDM` **dan** `TREATY_IN_EDM`.
   Deret `R09 -> R010` yang saya sendiri temukan (NA-02) hanya mungkin **karena** langkah 5 berjalan
   atas revisi yang dipilih — jadi temuan saya sudah memuat bantahannya.

Bentuknya bukan `MA-05` sampai `MA-08` (pemeriksaan lebih sempit daripada klaimnya), melainkan yang
baru: **membangun rantai sebab baru tanpa membaca ulang rantai yang sudah ditulis.** Keduanya ada
di berkas yang sama, dua bagian di atasnya.

**Yang ditambahkan ke daftar periksa, usulan `TA-08`:** setiap rantai sebab yang menyentuh mekanisme
yang **sudah punya temuan bernomor** dibaca ulang dari temuannya, bukan disusun ulang dari ingatan.
Bila hasilnya berbeda dari temuan itu, salah satunya salah — dan itu dinyatakan sebelum bukti baru
dipakai.

Dan satu lagi, yang ditemukan pada giliran yang sama: **`NB-03`** menunjukkan pertanyaan lanjutan B2
sudah dijawab `4-erd-dan-tabel-datar/STRUKTUR-DATA.md`. Itu **`MA-06` terulang** — sumber induk
tidak dibaca lebih dulu. `TA-05` sudah mengaturnya untuk **ADR**; ia diperluas: **seluruh artefak
induk di daftar masukan ronde**, bukan ADR saja.

---

## GRL-10 — Rantai versi: tidak ada percabangan, dan rujukan dasar disimpan tepat di satu tempat

* **Cabang** — B, `QB-2`.
* **Label per bagian:**

  | Bagian | Label | Berubah dari |
  |---|---|---|
  | versi dasar = **versi berlaku terakhir pada saat versi baru dibuat** | **PERUBAHAN** | dasar = **baris mana pun yang dipilih di picker**, termasuk versi lama dan termasuk baris kontrak (§4.4 butir 2, TDA-11) |
  | `OLDID` warisan dibawa apa adanya, beserta asal-usulnya | **PELESTARIAN** | — sejalan GRL-09 butir 2 |
  | rujukan dasar **disimpan eksplisit** sebagai `ID_VERSI_KONTRAK_DASAR` | **KONFIRMASI** | sudah diputuskan induk; lihat di bawah |

* **Keputusan**

  1. **Tidak ada percabangan di model baru.** Versi dasar sebuah penyesuaian adalah **versi berlaku
     terakhir pada saat versi baru itu dibuat** — bukan "versi *n* − 1". Rumusan "*n* − 1" salah,
     dan saya cabut: versi *n* − 1 dapat berkeadaan `DITOLAK` atau `DIBATALKAN`, dan versi terminal
     yang tidak disetujui **bukan dasar**.
  2. **`OLDID` warisan dibawa apa adanya** sebagai atribut warisan dengan asal-usulnya. Pada data
     lama ia **tidak** selalu menunjuk versi terakhir (NB-02), dan menurunkannya akan mengarang
     jawaban yang rapi untuk baris yang sebenarnya bercabang.
  3. **Bentuk simpannya bukan pertanyaan terbuka — ia KONFIRMASI.**
     `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` §1.1 sudah menetapkan
     **`ID_VERSI_KONTRAK_DASAR` disimpan eksplisit, tepat satu penunjuk, tanpa pasangan**, dengan
     alasan yang persis sama: picker lama meng-UNION kontrak dengan addendum sehingga dasar sebuah
     penyesuaian belum tentu versi tepat sebelumnya. `ERD.md` §2.6 menegaskan sisi lainnya —
     `NILAI_SELISIH` **tidak** punya relasi ke "versi lama"; sisi lama dibaca lewat induk barisnya.

     **Rekomendasi saya di B2 — menurunkan, jangan simpan — karena itu salah dan dicabut.** Dua
     alasannya runtuh: "satu fakta dua tempat" justru berbalik arah, sebab rujukan dasar **sudah**
     tersimpan untuk keperluan `NILAI_SELISIH` (`METODE` §8.3), sehingga menurunkannya di tempat
     lain akan **membuat** sumber kedua; dan INV-25 tidak cukup menopangnya, sebab ia membatasi
     versi tak-terminal, bukan menjamin dasar selalu versi sebelumnya.
  4. **Dua hal yang tetap perlu ditegaskan, dan keduanya sejalan dengan keputusan induk:**
     * **Versi non-material yang tidak menghasilkan satu pun baris selisih tetap butuh dasarnya
       diketahui.** Karena `ID_VERSI_KONTRAK_DASAR` menggantung pada `VERSI_KONTRAK` — bukan pada
       `NILAI_SELISIH` — ia terisi meski tabel selisihnya kosong. Bentuk induk sudah benar untuk
       kasus ini; yang perlu ditulis hanyalah **invariannya**.
     * **Rujukan dasar disimpan tepat di satu tempat**, yaitu `VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR`.
       Tidak ada salinan di `NILAI_SELISIH`, dan tidak ada turunan yang dihitung di tempat lain.

* **Bahan to-spec yang lahir bersama putusan ini**

  | # | Bunyi | Berubah dari |
  |---|---|---|
  | **B-3** | `ID_VERSI_KONTRAK_DASAR` **wajib terisi** pada setiap versi yang lahir dari penyesuaian, **termasuk yang tidak menghasilkan baris selisih** — dan wajib **kosong** pada versi pertama sebuah kontrak | tidak ada padanannya; sistem lama mengisi `OLDID` tanpa aturan tentang kapan ia harus kosong |
  | **B-4** | Versi yang ditunjuk `ID_VERSI_KONTRAK_DASAR` harus **versi berlaku terakhir** pada saat versi penunjuk dibuat, dan **tidak boleh** berkeadaan `DITOLAK` atau `DIBATALKAN` | dasar = baris mana pun yang dipilih di picker (TDA-11) |

  Keduanya **dapat gagal**, dan uji negatifnya jelas: coba simpan versi penyesuaian tanpa dasar,
  dan coba tunjuk versi yang sudah `DITOLAK`.

  **Cakupan B-3 dan B-4: versi NON-WARISAN saja.** Baris warisan melanggar keduanya sejak hari
  pertama, dan melanggarnya ke dua arah yang berlawanan:

  | Bila baris warisan diperlakukan begini | Yang dilanggar |
  |---|---|
  | `ID_VERSI_KONTRAK_DASAR` **diisi dari `OLDID`** | **B-4** — `OLDID` dapat menunjuk versi yang bukan terakhir (percabangan, `UA-1(c)`), bahkan yang sudah `DITOLAK` |
  | `ID_VERSI_KONTRAK_DASAR` **dibiarkan kosong** | **B-3** — ia versi penyesuaian, jadi dasarnya wajib terisi |

  Menyatakan cakupannya bukan pelonggaran: sebuah invarian yang **pasti** dilanggar data warisan
  akan dimatikan orang pada hari pertama migrasi, dan sesudah dimatikan ia tidak menjaga apa pun
  lagi (`METODE` §5.3). Lebih baik cakupannya disebut daripada penjaganya dicabut diam-diam.

* **Titipan ke `I2`** — dua pertanyaan, dan keduanya milik cabang I:

  1. **Bagaimana `ID_VERSI_KONTRAK_DASAR` diisi untuk baris warisan?** Pilihannya bukan hanya
     "dari `OLDID`" atau "kosong": ia dapat diisi **versi berlaku terakhir menurut urutan**, dengan
     `OLDID` mentah tersimpan di sampingnya sebagai keterangan. Keputusannya di I2, bersama
     keputusan umum tentang aturan baru mana yang boleh dijalankan atas baris warisan.
  2. **Kenapa `OLDID` mentah tetap disimpan di sampingnya** — dan alasannya **preseden GRL-06**
     (`METODE` §3.6), bukan kebiasaan. `OLDID` adalah **fakta yang tercatat pada barisnya**, bukan
     mekanisme; dan ia dapat menunjuk **baris yang isinya sudah tertimpa** (TDA-01, `UA-14`).
     Menurunkannya ulang menjadi "versi terakhir menurut urutan" akan **menghapus satu-satunya
     saksi** bahwa baris itu pernah menunjuk ke tempat lain — rapi, konsisten, dan tidak dapat
     dipulihkan. Itu persis kerusakan yang GRL-06 cegah.

* **Titipan ke cabang I** — baris warisan yang `OLDID`-nya menunjuk versi **bukan-terakhir** melanggar
  B-4 sejak hari pertama. Cabang I memutuskan apakah ia dipetakan apa adanya (dan B-4 hanya berlaku
  untuk versi baru) atau diperbaiki. Besarannya `UA-1` butir (c).
* **Arah dampak** — mengunci "tidak ada percabangan" padahal bisnis memang sengaja merevisi dari
  versi lama: ketahuan saat `UA-1(c)` tidak nol atau saat `DB-2` dibantah, dan perbaikannya
  **menambah satu keadaan sah pada B-4** — satu baris. Membiarkan percabangan: setiap pembaca
  selisih harus menelusuri pohon, dan tidak ada satu pun "versi berlaku" yang tunggal. Tidak
  setangkup; sisi murah diambil.
* **Status** — **terkunci**. `DB-2` (ditulis ulang) dan `UA-1(c)` terbuka, keduanya mempersempit.

---

## MA-10 · `TA-04` menemukan ulang daftar yang sudah ada di folder induk

`PETA-SUMBER-INDUK.md` baru saja saya susun, dan pada berkas pertama yang saya buka ternyata
tersimpan jawaban yang saya bangun ulang dari nol dua giliran lalu.

`4-erd-dan-tabel-datar/BENTUK-PENULIS-PROPERTI.md`, **24 September 2026** — hari yang sama dengan
`MA-05` — menyusun daftar bentuk penulis properti **secara mekanis**: 1.823 nama elemen di korpus,
disaring menjadi 174 kandidat, lalu **empat bentuk** yang benar-benar menulis properti kontrak,
lengkap dengan kolom *"Pernah disapu?"*. Tiga di antaranya ditandai **TIDAK**.

| Bentuk menurut induk | Padanan di `tools/tulis.py` |
|---|---|
| `<PropertiesName>` — Activity `Property-Set` | `PS` |
| `<pyPropertiesName>` — Data Transform | `DT` |
| `<pyPropertyTarget>` — Harness / Section | `SEC` |
| `<pyTargetProperty>` — Report Definition | `RD` — induk menandainya *"tidak perlu, bukan penulis halaman kontrak"*, sejalan dengan catatan saya |

Dan §2.1 berkas itu **sudah mendaftar** penulis Data Transform yang saya laporkan sebagai temuan
pada 25-26 September: `Position`, `ViewState`, `RevisionState` di `Akseptasi_DT`; `Ceding` dan
`CedingID` di `TreatyInSetReinsured`; `Termination` di `TreatyInSetTreatyYear`; `EDMState` di
`TreatyCreateEDM`.

Berkas kedua lebih telak lagi. `SAPUAN-DAN-NAMA-TAGNYA.md`, tanggal yang sama, sudah menetapkan
**aturannya**:

> **SEBUAH SAPUAN ATAS KORPUS YANG BERISI BANYAK JENIS ATURAN TIDAK BOLEH BERSANDAR PADA SATU NAMA
> TAG.** Setiap sapuan mencatat nama tag yang dipakainya di dalam dirinya sendiri, dan nama itu
> diperiksa terhadap `alat/datar-nama-elemen.csv` sebelum hasilnya dipercaya.

Itu `TA-04` saya, ditulis dua hari lebih dulu, dengan dasar yang **lebih kuat** — daftar nama elemen
yang disusun mekanis, bukan lima jenis yang saya pilih sendiri.

**Ini `MA-06` untuk ketiga kalinya**, dan ongkosnya kali ini paling besar: satu putaran penuh
dipakai membangun perkakas dan kalibrasi yang padanannya sudah ada. `TA-05` — baca sumber induk
lebih dulu — sudah saya tulis sesudah `MA-06`, dan **tidak saya jalankan** karena saya
memperlakukannya sebagai aturan untuk **ADR**, bukan untuk seluruh folder. `PETA-SUMBER-INDUK.md`
ada supaya kesalahan ini tidak dapat terjadi keempat kalinya.

**Satu hal yang tetap merupakan tambahan, dan diusulkan balik ke induk sebagai `IND-4`:** daftar
induk memuat empat bentuk **di sisi aturan**. `tools/tulis.py` menambahkan bentuk kelima di **sisi
antarmuka** — pengikatan sel (`pyValue` pada sel ber-`pyReadOnly != true`), yaitu penulisan yang
terjadi ketika pengguna mengetik. Bentuk itulah yang memunculkan `NA-01` (radio picker), `NA-17b`
(lapisan beku terbuka), dan seluruh analisis `ViewState` di `NA-21`/`NA-22`. Ia tidak ada di daftar
induk.

**Dan satu peringatan yang saya bawa pulang dari berkas itu:** `SAPUAN-DAN-NAMA-TAGNYA.md` §1
menandai `langkah-hidup.py` — penyisir `pyStepsBlockName` milik induk — sebagai **terdampak, "dan
ini yang terberat"**. Setiap klaim hidup/mati saya bersandar pada penanda yang sama. Bagian itu
**dibaca sebelum klaim hidup/mati berikutnya**, dan hasilnya dilaporkan apa adanya.

---

## GRL-11 — Versi berlaku adalah turunan; `KONTRAK` tidak membawa penunjuk

* **Cabang** — I (butir `I1`), diajukan lebih dulu karena mandiri terhadap `EXP-1`.
* **Label** — **PERUBAHAN**, berubah dari: *"kontrak tidak pernah digantikan; baris `TREATY_IN`
  selalu yang berlaku"* (§6.3 — dua langkah pengganti ber-blok `//`).
* **Keputusan**

  1. **Versi berlaku dihitung, tidak disimpan.** `KONTRAK` tidak membawa `ID_VERSI_BERLAKU`.
     Dasarnya ADR-0037 butir 2 — turunan tidak disunting dan tidak menulis balik — dan preseden
     terdekatnya baris terakhir tabel ADR-0037: *"'Boleh diubah' = turunan dari keadaan siklus
     hidup"*. Menyimpannya melanggar ADR-0041.
  2. **Rumusnya dinyatakan tepat, dan dicocokkan ke daftar keadaan ADR-0055.** Delapan keadaan:
     `DRAFT`, `MENUNGGU_SEC_HEAD`, `MENUNGGU_DEPT_HEAD`, `MENUNGGU_DIREKTUR`, `DISETUJUI`,
     `DITOLAK`, `DIBATALKAN`, `WARISAN_TAK_TERPETAKAN`. **Tidak ada keadaan "digantikan"**, dan
     `DISETUJUI` adalah **satu-satunya** keadaan yang berarti "pernah berlaku" — `DITOLAK` dan
     `DIBATALKAN` terminal tanpa pernah berlaku, sisanya belum selesai. Maka rumusnya:

     > **Versi berlaku sebuah kontrak = versi ber-`KEADAAN_SIKLUS_HIDUP = DISETUJUI` dengan
     > `NOMOR_URUT_VERSI` tertinggi. Kontrak tanpa versi `DISETUJUI` tidak punya versi berlaku, dan
     > itu keadaan yang sah.**

     Versi yang digantikan **tetap** `DISETUJUI`; ia berhenti berlaku bukan karena keadaannya
     berubah, melainkan karena ada yang bernomor lebih tinggi. Itu konsisten dengan ADR-0055 yang
     menjadikan `DISETUJUI` **terminal**.
  3. **Syarat yang membuat rumus itu benar, dinyatakan eksplisit:** *urutan `NOMOR_URUT_VERSI`
     harus sama dengan urutan waktu persetujuan.* Untuk versi baru, **INV-25** menjaminnya — paling
     banyak satu versi tak-terminal per kontrak, sehingga versi tidak dapat disetujui di luar
     urutan. **Untuk baris warisan, jaminan itu tidak ada**, dan menyediakannya adalah tugas `I2`.

* **Bahan to-spec**

  | # | Bunyi | Kasus yang dapat gagal |
  |---|---|---|
  | **I-1** | Versi berlaku = versi `DISETUJUI` ber-`NOMOR_URUT_VERSI` tertinggi; kontrak tanpa versi `DISETUJUI` tidak punya versi berlaku | kontrak dengan versi 3 `DISETUJUI` dan versi 4 `DITOLAK` -> versi berlaku **3**; kontrak yang hanya punya versi pertama `DISETUJUI` -> versi berlaku **versi pertama**; kontrak yang seluruh versinya `DIBATALKAN` -> **tidak ada** versi berlaku, dan itu bukan galat |
  | **I-2** | `NOMOR_URUT_VERSI` sebuah kontrak harus **searah** dengan urutan waktu persetujuan versinya | dua versi `DISETUJUI` yang nomornya terbalik terhadap tanggal persetujuannya -> ditolak saat migrasi |

* **Titipan ke `I2`, dengan ambiguitasnya disebut terang.** "Versi yang benar-benar berlaku hari
  ini" untuk kontrak warisan punya **dua arti yang dapat berbeda**:

  | Arti | Isinya | Di mana ia tersimpan |
  |---|---|---|
  | **hukum** | addendum terakhir yang disetujui | `TREATY_IN_EDM` |
  | **operasional** | baris yang selama ini dibaca hilir | `TREATY_IN` |

  Di sistem lama baris `TREATY_IN` **tidak pernah digantikan** addendum (§6.3), tetapi ia **dapat
  disunting di tempat** lewat `Revision` — termasuk **sesudah** addendum terakhir dibuat. Maka
  kedua arti itu dapat menunjuk isi yang berbeda, dan **urutan di antara keduanya hanya dapat
  direkonstruksi dari cap waktu**: komentar `"Create Revision"` pada baris kontrak, dan
  `RevisionDate` pada addendum. Diukur **`UA-17`**.

  Akibatnya ditarik sampai ke orang (`METODE` §4.5): memilih arti hukum berarti angka yang keluar
  ke akuntansi dan ke pihak lawan **berubah** pada hari peralihan — itu **`DB-12`**, dan ia
  menggantung pada keputusan ini.

  Ditambah dua titipan yang sudah ada: pengisian `ID_VERSI_KONTRAK_DASAR` untuk baris warisan, dan
  alasan menyimpan `OLDID` mentah di sampingnya (GRL-10).

* **Ditolak** — **(b)** penunjuk tersimpan: kolom yang harus dijaga konsisten pada setiap
  persetujuan, pembatalan, dan migrasi; bila ia menyimpang, satu-satunya cara mengetahuinya adalah
  menghitung ulang turunannya — yang berarti (a) tetap harus ada. **(c)** penunjuk khusus warisan:
  membawa cabang permanen di dalam model untuk masalah berumur satu hari, dan cabang itu tidak akan
  pernah dibuang; keputusan migrasinya dapat dituliskan seluruhnya ke dalam `NOMOR_URUT_VERSI` dan
  `KEADAAN`.
* **Arah dampak** — (a) salah: ongkosnya **menambah satu view atau kolom terhitung**, dapat
  dilakukan kapan saja tanpa mengubah kebenaran data. (b) salah: kolom yang menyimpang diam-diam,
  dan tidak ada cara mengetahuinya. Tidak setangkup.
* **Status** — **terkunci**. `UA-17` dan `DB-12` terbuka; keduanya mempersempit arti "berlaku",
  bukan membatalkan bentuknya.

---

## GRL-12 — Materialitas diturunkan dari baris selisih yang tersimpan; jenis tidak menurunkannya

* **Cabang** — C (butir `C1`), gerbang bagi C2, E1, dan E3.
* **Label** — **PERUBAHAN**, berubah dari: *"dipilih pengguna lewat radio; Non Material mengunci
  field uang di layar; tanpa penjaga di sisi simpan"* (NA-01, TDA-10, TDA-13).
* **Keputusan**

  1. **Untuk versi baru: material = versi itu memiliki sedikitnya satu baris `NILAI_SELISIH`.**
     Tidak ada atribut materialitas pada versi baru. `JENIS_ADDENDUM` tetap sumbu terpisah dan
     **tidak** menurunkan materialitas.
  2. **Diturunkan dari baris selisih yang TERSIMPAN, bukan dari hitung ulang.** Materialitas sebuah
     versi dibaca dari baris `NILAI_SELISIH` miliknya, yang **beku sesudah persetujuan** (GRL-02,
     ADR-0036). Akibatnya penting dan disengaja: **perbaikan rumus `E1` di kemudian hari tidak
     mengubah materialitas versi yang sudah disetujui.** Bila materialitas dihitung ulang saat
     dibaca, satu perubahan rumus akan mengklasifikasi ulang sejarah — bentuk kerusakan yang sama
     yang GRL-06 cegah.
  3. **Nilai warisan tetap disimpan.** Rumusan butir 1 berbunyi *"tidak ada atribut untuk versi
     **baru**"*. `EDMMATERIALTYPE` warisan **dibawa apa adanya** sebagai nilai warisan beserta
     asal-usulnya (GRL-06, ADR-0042). Pembanding yang dihitung `UA-3` **disimpan terpisah** dan
     tidak pernah menggantikannya.
  4. **`TDA-10` dinyatakan diperbaiki.** Materialitas turunan **tidak dapat berbohong**: ia bukan
     pernyataan yang perlu ditegakkan, melainkan pembacaan atas apa yang tersimpan. Cacat "nol
     penegakan di sisi simpan" hilang bersama hal yang perlu ditegakkan.

* **Alasan menolak (b) — turunan dari jenis — memakai bukti kombinasi, bukan tanggal.**
  Di radio picker, `EDMState` menawarkan `1` dan `2`, dan `EDMMaterialType` menawarkan `1` dan `2`,
  **keduanya bebas dipasangkan** (`EXP-1`, NA-01). Jadi di sistem lama jenis 1 dan jenis 2
  masing-masing dapat bermaterialitas material **maupun** non-material — materialitas **tidak dapat**
  diturunkan dari jenis, sebab satu jenis memuat dua-duanya.

  Tanggal pembuatan properti — `EDMState` Juni 2020, `EDMMaterialType` Oktober 2020 — tetap dicatat
  sebagai **fakta sejarah, bukan bukti niat**, dan **tidak dipakai** sebagai argumen arah.

* **Alasan menolak (c) — masukan yang diperiksa — diganti.** Alasan saya sebelumnya, *"tidak ada
  pihak yang pernah menyatakan niat dengan sadar"*, **dicabut**: itu tidak terbaca dari ekspor
  (`METODE` §4.3, §4.6), dan yang terbaca justru sebaliknya — memilih Non Material **mengunci 220
  kondisi field uang** di layar, jadi pilihan itu punya akibat fungsional. Dua alasan yang sah:

  | # | Alasan |
  |---|---|
  | 1 | **`METODE` §3.8** — (c) memasang **penolakan baru** pada hari peralihan: pekerjaan yang hari ini lolos akan ditolak karena niat dan akibat tidak cocok. Seberapa sering itu terjadi diukur **`UA-3`**, dan angkanya belum ada |
  | 2 | **ADR-0037 butir 3** — *"Bila sebuah besaran boleh disepakati, ia dinaikkan menjadi masukan."* Materialitas **tidak disepakati**; ia sifat dari apa yang berubah. Menaikkannya menjadi masukan tanpa kesepakatan berarti membalik aturan itu |

* **Bahan to-spec**

  | # | Bunyi | Kasus yang dapat gagal |
  |---|---|---|
  | **C-1** | Materialitas sebuah versi = benar bila versi itu punya sedikitnya satu baris `NILAI_SELISIH` tersimpan; dibaca dari baris tersimpan, tidak dihitung ulang | versi yang hanya mengubah keterangan -> **non-material**; versi yang mengubah satu limit -> **material**; versi yang sudah disetujui **tetap bermaterialitas sama** sesudah rumus selisih diubah |
  | **C-2** | `JENIS_ADDENDUM` tidak menentukan materialitas, dan materialitas tidak menentukan jenis | versi berjenis "revisi internal" yang mengubah limit -> **material**, dan jenisnya tetap "revisi internal" |

* **Pembatal — `DB-15`**, ditulis dalam bahasa bisnis:

  > *"Pembedaan material dan non-material tidak dipakai di luar layar pengisian — ia tidak menentukan
  > jalur persetujuan, tidak menentukan jenis dokumen yang dikirim ke cedant, tidak menentukan
  > pencatatan akuntansi, dan tidak muncul di laporan mana pun. Benar?"*

  Bila **dibantah dengan satu pemakaian yang harus diketahui SEBELUM perubahannya dibuat**,
  keputusan ini ditinjau ulang ke arah **(c)** — sebab hanya niat yang dinyatakan di muka yang dapat
  dipakai sebelum akibatnya ada.

* **Titipan ke cabang H — akibat bagi orang di hari pertama** (`METODE` §4.5). Mode "Non Material
  mengunci field uang" **hilang** dari layar. Tanpa atribut, materialitas hanya diketahui **sesudah**
  perubahan dibuat. Dua hal yang mengikutinya, keduanya bukan keputusan grilling:

  | # | Isi |
  |---|---|
  | **R-H1** | perlu tidaknya **bantuan layar yang tidak disimpan** — misalnya peringatan *"perubahan ini mengubah angka"* sebelum pengajuan. Ia bantuan, bukan atribut; tidak disimpan, tidak menjadi invarian |
  | pemberitahuan | pengisi kontrak perlu diberi tahu sebelum peralihan bahwa mode pengunci itu hilang |

* **`REV-4`** diajukan untuk ADR-0049 — **bukan** keputusannya bahwa materialitas tidak menjadi
  kolom, melainkan **sumber turunannya**: dari jenis menjadi dari baris selisih. Aturan paket REV
  sama: ADR tidak disunting, diserahkan sekaligus di akhir grilling.
* **Arah dampak** — (a) salah, ternyata bisnis memutuskan materialitas menurut pertimbangan yang
  tidak tercermin di angka: ongkosnya **menaikkan materialitas menjadi masukan**, persis prosedur
  ADR-0037 butir 3, dan nilai warisannya sudah tersimpan apa adanya sehingga tidak ada yang hilang.
  (b) salah: mengunci definisi yang membuat pemeriksaan mustahil. (c) salah: memasang kewajiban baru
  yang mungkin tidak ada yang meminta. Keduanya lebih mahal dibatalkan.
* **Status** — **terkunci**. `DB-15` terbuka sebagai pembatal; `UA-3` mengukur ongkos (c) bila kelak
  ditinjau ulang.

---

## GRL-13 — `JENIS_ADDENDUM` bernilai dua untuk versi baru

* **Cabang** — C (butir `C2`).
* **Label per bagian** (`METODE` §6.5):

  | Bagian | Label | Berubah dari |
  |---|---|---|
  | `PENYESUAIAN_PREMI` dibawa sebagai nilai jenis | **PELESTARIAN** | — perilaku jenis 3 dibawa: ia menyemai pohon nilai aktual, menyalakan penandanya, dan membelokkan jalur simpan |
  | pembedaan Internal/External **tidak** dibawa untuk versi baru | **PERUBAHAN** | *"tiga nilai `EDMState`; Internal dan External dipilih pengguna di radio tanpa akibat selain teks komentar"* |

* **Keputusan**

  1. **`JENIS_ADDENDUM` bernilai dua untuk versi baru:** `PENYESUAIAN_PREMI` dan satu nilai untuk
     selebihnya.
  2. **`DB-16` adalah pintu ke (b+).** Bila dibenarkan — addendum eksternal berupa dokumen yang
     disepakati cedant, bertanggal sendiri — penandanya ditambahkan. **Tempatnya diputuskan sesudah
     `DB-16` terjawab**, sebab bila "eksternal" berarti dokumen, tempatnya mungkin
     `DOKUMEN_KONTRAK`, bukan enum jenis.
  3. **`EDMState` warisan dibawa apa adanya** sebagai nilai warisan beserta asal-usulnya — pola
     GRL-06, GRL-09, GRL-12. Tidak ada data lama yang hilang; yang dipertaruhkan hanya versi **baru**
     yang lahir sebelum `DB-16` terjawab.
  4. **`REV-5` diajukan** untuk himpunan jenis ADR-0049, yang runtuh sendiri sesudah GRL-12:
     "administratif" sama dengan non-material turunan, sehingga tersisa dua — dan dua itu persis
     bentuk ini.

* **Nama nilai kedua: bahan to-spec, dengan satu petunjuk domain dan satu peringatan.**

  `EGNPI` lazimnya *Estimated Gross Net Premium Income*, dan jenis 3 menyalin `EGNPI` ke
  `ActualValue.EGNPI` — penyesuaian dari **estimasi** ke **aktual**. Jadi kata "estimasi" di
  ADR-0049 punya **jangkar domain**, walaupun tidak tertulis di ekspor.

  **Tetapi "perubahan estimasi" kemungkinan terlalu sempit sebagai nama nilai kedua**: kombinasi
  (1,1) dan (2,1) juga mengubah limit, share, dan tarif — bukan hanya estimasi. Nama finalnya
  **diputuskan di to-spec**; yang dikunci di sini hanya **jumlah nilainya**, bukan namanya.

  Butir daftar untuk dibantah: **`DB-17`** — *"`EGNPI` pada kontrak treaty berarti perkiraan
  pendapatan premi bruto neto untuk periode itu, dan 'nilai aktual' berarti realisasinya. Benar?"*
  Dikaitkan ke **C3**.

* **Syarat yang WAJIB dijawab `C3`, dan ia sudah punya bukti yang memberatkan.**

  Di sistem lama, *"penyesuaian premi selalu material"* **terbaca** (`EXP-1`: radio tidak menawarkan
  jenis 3, sehingga ia tidak dapat dipasangkan dengan Non Material). Di model baru, materialitas
  diturunkan dari baris `NILAI_SELISIH` (GRL-12). Maka bila penyesuaian premi hanya mengubah nilai
  aktual, dan nilai aktual tidak menghasilkan baris selisih, versi `PENYESUAIAN_PREMI` akan terbaca
  **non-material** — bertentangan dengan perilaku lama.

  **Dan itu bukan kekhawatiran teoretis: `TDA-16` menunjukkan sistem lama memang menghasilkan
  selisih EGNPI nol pada addendum premi** (NB-08). Jadi menyalin perilaku lama apa adanya **tidak**
  menyelamatkan keadaan ini; ia justru sumbernya.

  | # | Bunyi | Kasus yang dapat gagal |
  |---|---|---|
  | **C-3** | Versi `PENYESUAIAN_PREMI` yang mengubah nilai aktual **menghasilkan baris `NILAI_SELISIH`**, sehingga terbaca material | versi `PENYESUAIAN_PREMI` yang hanya mengubah premi aktual -> **material**, bukan non-material |

* **Titipan ke cabang H** — radio Internal/External **hilang** dari layar. **`R-H2`**: perlu tidaknya
  penggantinya, dan dalam bentuk apa. Masuk pula ke daftar pemberitahuan untuk pengisi kontrak,
  bersama `R-H1` (mode pengunci field uang yang juga hilang).
* **Ditolak** — **(c)** tanpa sumbu jenis: menghapus satu-satunya pembedaan yang terbaca dari
  perilaku, dan mengubah boolean menjadi enum kelak lebih mahal daripada menambah satu nilai enum.
  **(b+) sekarang**: memilih tempat penanda sebelum tahu benda apa yang ditaruh.
* **Arah dampak** — (b) salah, Internal/External ternyata bermakna: ongkosnya satu nilai enum
  ditambahkan atau satu atribut di tempat yang lebih tepat, dan versi baru yang lahir sementara itu
  kehilangan klasifikasinya — sebanyak jarak sampai `DB-16` terjawab. (c) salah: bentuknya berubah
  dari boolean ke enum, menyentuh setiap pembaca. Tidak setangkup.
* **Status** — **terkunci**, dengan `DB-16` sebagai pintu, `DB-17` menopang penamaan, dan **`C-3`
  sebagai syarat yang dibawa ke C3**.
