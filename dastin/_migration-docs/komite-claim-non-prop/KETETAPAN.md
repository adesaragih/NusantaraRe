> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : penilai
> Masukan: transkrip sesi `5696a074-4b6f-422e-a2aa-86e9a9c4b1ff.jsonl` (2026-09-20) · `PUTUSAN-01.md` §9 · `PENGETAHUAN.md` · `INVENTARIS-BUKTI.md` · `FAKTA-A1B.md` · `GRILL-05/06-PUTUSAN.md` · `GRILL-06/01-PEMBACAAN.md` · penyuntingan penutup ronde 6
> Status : TERBUKA
> Sifat  : HIDUP

# KETETAPAN — Komite Claim Non Prop

**Ini satu-satunya sumber ketetapan modul ini.** Rujuk dengan nomor; jangan menulis ulang
bunyinya di berkas lain. Ketetapan yang tidak ada di sini tidak mengikat siapa pun,
termasuk yang diyakini pernah diputuskan.

## Riwayat pemutakhiran

| Ronde | Tanggal | Apa yang masuk |
|---|---|---|
| 05 | 2026-09-20 | Berkas ini dibuat. Isi bagian 1–10 dipulihkan **dari transkrip sesi di disk**, bukan dari ingatan; bunyi disalin apa adanya. Ditambahkan bagian 11 (ketetapan ronde 5) dan peta penomoran di bawah |
| 06 | 2026-09-20 | `K5-7` · bunyi penuh `H-4` · jawaban `QF-1`, `QF-2`, `QF-3` sebagai deret `K6` · empat aturan kerja baru · akibat pembacaan ronde 6 terhadap ketetapan lama · status enam bunyi `H` |
| 06 · penutup | 2026-09-20 | `K6-4` menggantikan `H-2` yang bunyinya hilang · dua aturan kerja (kepala berkas hidup, aturan berhenti grilling) · kepala berkas hidup dimutakhirkan ke Ronde 06 |

## Peta penomoran

Deret ketetapan ronde 3 ditulis `G-1 … G-4` di seluruh badan berkas ini karena begitulah ia
lahir. Sejak 2026-09-20 nomor resminya **`J-1 … J-4`**, untuk mengakhiri tabrakan dengan
`G-01 … G-20` yang dipakai `PUTUSAN-01.md` bagi temuan interogasi (lihat bagian 10.2).
Judul bagian 3 dan keempat barisnya sudah memakai nomor baru; rujukan silang di dalam
badan ketetapan lain sengaja **tidak disunting**, agar bunyi aslinya tetap dapat diperiksa.

| Nomor lama | Nomor resmi | Muncul juga di |
|---|---|---|
| `G-1` | **`J-1`** | E-1 Status; `FAKTA-A1B.md` |
| `G-2` | **`J-2`** | bagian 9 baris 7; catatan F-12 |
| `G-3` | **`J-3`** | D-3 Status |
| `G-4` | **`J-4`** | — |

Deret yang tidak pernah dipakai ulang: **D** (ronde 1) · **E** (ronde 2) · **J** (ronde 3,
semula G) · **H** (ronde 4) · **K5** (ronde 5) · **C** (koreksi penilai) · **F** (fakta).

---

## 1. Ketetapan ronde 1 — D-1 … D-5

```
D-1 · Paritas atau perbaikan, dibelah per kelas risiko
  Bunyi   : Angka yang sudah terbit (nomor akseptasi lama, saldo OS, nilai terkirim):
            paritas mutlak, tidak dihitung ulang. Kebenaran catatan & wewenang (K-01,
            K-04, K-07, K-08, G-17): perbaikan wajib, paritas ditolak. Cacat berakibat
            data (G-08, G-15): perbaiki ke depan, tidak surut. Selebihnya: paritas
            default. Prinsip pengikat: di mana grilling menghasilkan cacat ber-EVIDENCED,
            perbaikan menjadi default dan paritas yang harus dibela. Setiap perbaikan
            masuk daftar DEVIASI DIHARAPKAN dengan uji paritas yang memang dirancang
            gagal. Ditandai sebagai paritas-dengan-pertanyaan-bisnis, bukan cacat:
            penolakan satu adjustment membuat seluruh klaim CLAIM REJECTED
            (KomitePostAdjustment 20).
  Alasan  : Cacat yang memalsukan riwayat orang tidak boleh diwarisi diam-diam, sementara
            angka uang wajib identik agar shadow-run tetap berdaya sebagai alat bukti.
  Asal    : ronde 1 (ketetapan D-1)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
D-2 · Batas modul
  Bunyi   : Komite adalah modul di dalam konteks Klaim, dengan agregat sendiri. Bukan
            konteks berjaringan. Alasannya bukan KOMITECOUNT/FLAGONGOINGCOMMITTE menempel
            di PC_ASM_FW_GCNMFW_WORK — itu artefak cara Pega menulis properti pyWorkPage
            menjadi kolom, bukan pernyataan tentang batas domain. Alasan yang mengikat:
            (a) kontrak tulis §7.2 menuntut keputusan komite dan status adjustment induk
            berubah bersama atau tidak sama sekali; (b) titik commit keuangan ada di modul
            ini. Ikut diputuskan: KOMITECOUNT dan FLAGONGOINGCOMMITTE tidak dimigrasikan
            sebagai kolom — turunan dari jenjang, dihitung saat dibaca.
  Alasan  : Memisahkan menjadi konteks berjaringan berarti memilih transaksi terdistribusi
            untuk proses yang di sistem lama pun satu commit.
  Asal    : ronde 1 (ketetapan D-2)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
D-3 · Satu konsep sirkulasi, tiga jenis
  Bunyi   : Satu konsep sirkulasi, tiga jenis. HASIL selalu berarti putusan komite atas
            usulan yang diajukan, tidak pernah nasib klaim. Nasib klaim adalah fungsi dari
            (jenis, hasil), ditulis sekali di satu tempat:
            ADJUSTMENT/DISETUJUI → adjustment diakseptasi
            ADJUSTMENT/DITOLAK → klaim ditolak
            TUTUP_KLAIM/DISETUJUI → klaim ditutup
            TUTUP_KLAIM/DITOLAK → klaim tetap berjalan
            TOLAK_KLAIM/DISETUJUI → klaim ditolak
            TOLAK_KLAIM/DITOLAK → klaim tetap berjalan
  Alasan  : Mekanismenya identik sampai ke tingkat langkah; yang berbeda hanya akibatnya.
  Asal    : ronde 1 (ketetapan D-3)
  Tanggal : 2026-09-20
  Status  : BERLAKU — diperluas satu dimensi oleh G-3
```

```
D-4 · Penerbit nomor akseptasi
  Bunyi   : Nomor akseptasi diterbitkan modul Komite, di dalam transaksi keputusan,
            memakai sumber urutan Oracle yang ada, dengan pembungkus pencatat di KLAIMNP.
            Periode buku ditentukan di satu lapis: Oracle. Sisi Go tidak menggeser ulang
            dan tidak menambal. Ambang 25 di rule dan @replaceAll 14.14 tidak dipindahkan.
            Penerbitan di dalam transaksi yang sama dengan pencatatan keputusan final,
            bukan outbox — nomor bukan efek luar, ia fakta internal yang menjadi syarat
            efek luar. Idempotensi ditegakkan constraint, bukan baca-lalu-tulis. Fase 1
            memakai prosedur lama apa adanya (hilir memeriksa panjang 23/24 karakter).
  Alasan  : Nomor dipakai sistem lain yang tidak ikut ditulis ulang, sehingga
            memindahkannya sekarang mempertaruhkan pencocokan pembayaran.
  Asal    : ronde 1 (ketetapan D-4)
  Tanggal : 2026-09-20
  Status  : BERLAKU — dasar faktualnya (isi badan PROC_GENERATE_SEQUENCE_NUMBER)
            bertanda TAFSIR (menunggu C-01)
```

```
D-5 · Pemilik roster komite
  Bunyi   : Roster komite dimiliki sistem baru, disemai sekali dari EMAILKOMITE. Pilihan
            membaca tabel lama gugur: V_MST_USER_TEKNIS menarik nama lewat database link
            lintas sistem ke HR, sehingga memilihnya berarti melanggar ADR-0016 pada hari
            pertama. Tidak ada pembacaan runtime ke tabel lama. Identitas pengguna datang
            dari lapis aplikasi. Keputusan beku no. 4 (peran disalin ke keputusan saat
            diambil) membuat roster hanya dibutuhkan saat seleksi. Tidak ikut terbuka:
            FilterEmailKomiteWithLimit dan CreateChildKomiteCNP_Act. Yang diputuskan
            adalah pemilik rumahnya, bukan siapa yang berhak masuk.
  Alasan  : Roster adalah data organisasi yang berubah tiap mutasi jabatan; menaruhnya di
            tabel milik sistem lain membuat sistem baru tidak pernah dapat menegakkan
            aturan wewenangnya sendiri.
  Asal    : ronde 1 (ketetapan D-5)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

---

## 2. Ketetapan ronde 2 — E-1 … E-5

```
E-1 · Menyunting usulan sambil memutus
  Bunyi   : Menyatu, terekam dua peristiwa. Hanya jenjang 1 yang boleh menyunting penanda
            usulan (paritas, F-2). Keputusan jenjang 1 tercatat sebagai `usulan diubah`
            (membawa nilai sebelum & sesudah) dan `jenjang 1 memutus`, satu transaksi.
            Tiap keputusan mencatat versi usulan yang diputusnya; versi naik hanya saat
            penanda berubah. Turunan yang sudah tertetapkan: karena penyuntingan hanya
            mungkin di jenjang 1 dan terjadi tepat saat ia memutus, jenjang 2..n melihat
            keadaan identik satu sama lain — maka aturan "perubahan usulan membatalkan
            persetujuan sebelumnya" tidak ditulis; ia invarian yang dijamin bentuk, bukan
            aturan yang ditegakkan. Jenjang 1 menolak → suntingan tetap tercatat sebagai
            peristiwa, tanpa akibat hilir.
  Alasan  : Mempertahankan paritas perilaku sambil berhenti menyembunyikan bahwa objek
            berubah saat diputus.
  Asal    : ronde 2 (ketetapan E-1)
  Tanggal : 2026-09-20
  Status  : DIREVISI oleh G-1
```

```
E-2 · Kegagalan penerbitan nomor
  Bunyi   : Gagal menerbitkan nomor → seluruh transaksi keputusan batal. Pilihan
            "keputusan tersimpan, nomor menyusul" dan "keadaan menunggu nomor" keduanya
            ditolak. Empat ikatan: (1) tidak ada keputusan tersimpan, tidak ada nomor;
            (2) galat gangguan sistem wajib dapat dibedakan pengguna dari penolakan
            wewenang; (3) aksi keputusan membawa kunci idempotensi — percobaan ulang
            dengan kunci sama mengembalikan hasil yang sama, bukan nomor kedua (bahaya
            sesungguhnya adalah transaksi berhasil tetapi tanggapannya hilang);
            (4) masukan approver tidak hilang saat gagal.
  Alasan  : Hilir memperlakukan nomor sebagai syarat keberadaan akseptasi, sehingga
            akseptasi tanpa nomor mengekspor keadaan baru ke sistem yang tidak ikut
            ditulis ulang.
  Asal    : ronde 2 (ketetapan E-2)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
E-3 · Jenjang sebagai kedudukan
  Bunyi   : Diputuskan (A-1a): jenjang adalah kedudukan yang diisi seseorang, bukan
            seseorang. Pemegang adalah atribut kedudukan; perpindahan pemegang berwujud
            peristiwa tersendiri dengan pelaku, alasan, waktu. Keputusan ini berdiri
            karena menopang keputusan beku no. 4 dan keadaan TIDAK_SAMPAI (K-04), bukan
            karena pengalihan. Dipagari (A-1b): siapa berhak menjadi pengganti dan siapa
            berhak mengalihkan — bergantung aturan seleksi yang beku. Aksi pengalihan
            tidak diekspos sampai A-1b cair; modelnya menampung, permukaannya tidak ada.
            Didaftarkan: risiko pinjam akun sebagai risiko operasional terbuka, dimiliki
            pemilik proses.
  Alasan  : Menopang keputusan beku no. 4 dan keadaan TIDAK_SAMPAI.
  Asal    : ronde 2 (ketetapan E-3)
  Tanggal : 2026-09-20
  Status  : BERLAKU — bagian yang dipagari dijawab H-3 dan H-4
```

```
E-4 · Medan induk sebagai penghalang
  Bunyi   : Prinsip: keputusan komite tidak pernah diblokir oleh data yang tidak dapat
            disunting approver di layar ini; medan wajib yang pelakunya tak dapat mengisi
            bukan validasi, melainkan jebakan. Bila kelimanya pyReadOnly=true → tidak
            memblokir apa pun, dan itu paritas (aturan wajibnya memang tidak pernah
            menyala), serta G-20 dinilai ulang. Bila dapat disunting → memblokir hanya
            SETUJU.
  Alasan  : Medan wajib yang tidak dapat diisi pelakunya adalah jebakan, bukan validasi.
  Asal    : ronde 2 (ketetapan E-4)
  Tanggal : 2026-09-20
  Status  : BERLAKU — cabang pertama yang berlaku, ditetapkan oleh F-7
```

```
E-5 · Umur sirkulasi
  Bunyi   : Umur bukan keadaan. Tidak ada KEDALUWARSA, tidak ada pekerja latar, tidak ada
            kolom tulis baru. Umur = atribut turunan dari cap waktu yang sudah disimpan
            (usia sirkulasi sejak dibuat, usia jenjang aktif sejak menjadi aktif),
            dihitung saat dibaca.
  Alasan  : Menghasilkan data yang dibutuhkan untuk memutuskan SLA nanti dengan bukti,
            bukan dengan angka karangan.
  Asal    : ronde 2 (ketetapan E-5)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

---

## 3. Ketetapan ronde 3 — J-1 … J-4 (semula G-1 … G-4)

```
J-1 · Revisi E-1: penjaga di batas tulis
  Bunyi   : E-1 direvisi. Invarian tetap dijamin bentuk, tetapi ditegakkan juga sebagai
            penjaga murah di batas tulis: keputusan ditolak bila versi usulan berubah oleh
            jenjang selain yang pertama. Alasannya bukan ketidakpercayaan pada bentuk,
            melainkan bahwa bentuk yang dijamin 296 kendali layar akan berubah tanpa ada
            yang memberi tahu.
  Alasan  : Bentuk yang dijamin oleh 296 kendali layar akan berubah tanpa ada yang
            memberi tahu.
  Asal    : ronde 3 (ketetapan G-1)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
J-2 · .CNPFlagXOL dan .CurrencyID
  Bunyi   : .CNPFlagXOL dikunci ke jenjang 1 (penanda usulan menurut definisi, alasan
            seluruhnya A-1a), masuk versi usulan. .CurrencyID baca-saja seluruhnya —
            bukan dikunci ke jenjang 1: mengganti label mata uang sementara nilainya beku
            (F-9) bukan kemampuan, itu kejadian keempat pola panel-salinan. Keduanya
            DEVIASI DIHARAPKAN.
  Alasan  : Mengganti label mata uang sementara nilainya beku bukan kemampuan.
  Asal    : ronde 3 (ketetapan G-2)
  Tanggal : 2026-09-20
  Status  : BERLAKU — melarutkan F-12
```

```
J-3 · Akseptasi bersyarat sebagai atribut usulan
  Bunyi   : Akseptasi bersyarat = atribut usulan, bukan hasil ketiga. HASIL tetap
            DISETUJUI/DITOLAK. Tabel D-3 diperluas satu dimensi, bukan satu baris:
            akibat = fungsi (jenis, bersyarat, hasil), tetap di satu tempat. Model baca
            menyajikan tiga nilai tampil dari fungsi yang sama; laporan tidak pernah
            membaca HASIL telanjang. Syarat terbit nomor: DISETUJUI pada jenjang terakhir
            DAN usulan tidak bersyarat.
  Alasan  : Akibat tetap dihitung di satu tempat, dan laporan tidak pernah membaca HASIL
            telanjang.
  Asal    : ronde 3 (ketetapan G-3)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
J-4 · Hak baca sirkulasi
  Bunyi   : Hak baca sirkulasi diturunkan dari hak baca klaim (baca sirkulasi ⊆ baca
            klaim). Ini paritas: komentar komite sudah masuk berkas klaim lewat
            ComiteeClaim (§5.2 langkah 11, §7.2). Pembatasannya adalah perubahan sah,
            didaftarkan sebagai pertanyaan terbuka ke pemilik proses — bukan keputusan
            spesifikasi, bukan pertanyaan ronde 4.
  Alasan  : Komentar komite sudah menjadi bagian berkas klaim di sistem lama.
  Asal    : ronde 3 (ketetapan G-4)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

---

## 4. Ketetapan ronde 4 — H-1 … H-7

> Sumber bunyi: enumerasi pada instruksi `KETETAPAN-KOMITE`, bagian 4. Bunyi yang tersedia berbentuk enumerasi ringkas, bukan kalimat normatif penuh — dicatat di bagian 10.

```
H-1 · Ambang berbasis nilai
  Bunyi   : Dua kelas dimigrasikan, aturan seleksi sebagai data, cabang mati dan lubang
            di atas 50 juta ditutup, pemulihan naik ke pemilik proses.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : ronde 4, jawaban atas FAKTA-A1B Q-1 (rekomendasi penulis adalah (a); yang
            diputuskan adalah (b) dengan tiga tambahan)
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
H-2 · Jalur tutup dan tolak klaim
  Bunyi   : Jalur tutup/tolak satu jenjang dengan pemegang dari roster berdasarkan peran,
            dan jenis dinyatakan pemanggil.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : ronde 4, jawaban atas FAKTA-A1B Q-2 (sejalan rekomendasi (a))
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
H-3 · Delegasi
  Bunyi   : Delegasi tetap sebagai atribut roster.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : ronde 4, jawaban atas FAKTA-A1B Q-3
  Tanggal : 2026-09-20
  Status  : BERLAKU — menjawab bagian E-3 yang dipagari
```

```
H-4 · Pengalihan sesaat
  Bunyi   : Pengalihan sesaat tetap tidak diekspos.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : ronde 4, jawaban atas FAKTA-A1B Q-3
  Tanggal : 2026-09-20
  Status  : DIGANTIKAN ronde 6 oleh bunyi penuh di bawah; baris ini tidak dihapus
```

> Bunyi ringkas di atas tidak pernah dicabut; ia dilengkapi. Kehilangan yang dicatat bagian
> 10.1 untuk `H-4` dengan demikian **tertutup**.

```
H-4 · Pengalihan sesaat  (bunyi penuh, dipulihkan ronde 6)
  Bunyi   : Pengalihan sesaat tetap tidak diekspos — bukan lagi karena A-1b beku, melainkan
            karena kebutuhan nyata yang selama ini ditambal satu baris kode (F-19) sudah
            terpenuhi oleh delegasi tetap sebagai atribut roster (H-3). Pengalihan sesaat
            menjadi kemampuan yang belum pernah ada yang meminta. Modelnya menampung (E-3),
            permukaannya tidak ada. Peran yang berwenang adalah keputusan pemilik proses,
            dengan default peran administratif tersendiri di luar keanggotaan komite.
  Alasan  : (tersirat pada bunyi; tidak disertakan terpisah pada sumber)
  Asal    : ronde 4, jawaban atas FAKTA-A1B Q-3; bunyi penuh dipulihkan ronde 6
  Tanggal : 2026-09-20
  Status  : BERLAKU — menggantikan bunyi ringkas ronde 4; dasar PG-02 dengan demikian
            dinyatakan, dan PG-02 berhenti menjadi pagar
```

```
H-5 · Validasi pembuatan
  Bunyi   : Validasi sebagai bagian kasus-guna pembuatan yang menolak, bukan menandai.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : ronde 4, adjudikasi F-21
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
H-6 · Sirkulasi tanpa jenjang
  Bunyi   : Sirkulasi tanpa jenjang ditolak saat dibuat, sebagai deviasi yang diuji.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : ronde 4, adjudikasi F-22; menguatkan keputusan beku no. 8
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
H-7 · Jalur VINCENTVERNANDO_1
  Bunyi   : Jalur VINCENTVERNANDO_1 tidak dimigrasikan.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : ronde 4, adjudikasi F-23
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

---

## 5. Koreksi penilai — C-01 … C-05

```
C-01 · Provenans DDL
  Bunyi   : 49 berkas DDL adalah reverse-engineer Toad dari sebuah Excel, bukan ekspor
            langsung — yang dibandingkan bukan hanya LAST_DDL_TIME melainkan isi badan
            prosedur PROC_GENERATE_SEQUENCE_NUMBER. Sampai itu dijalankan, lima penutupan
            yang bersandar padanya bertanda TAFSIR (menunggu C-01).
  Alasan  : Berkas hasil reverse-engineer tidak membuktikan isi objek yang berjalan hari
            ini.
  Asal    : instruksi KETETAPAN-KOMITE bagian 5; stempel asal berkas DDL ("SUMBER:
            DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)")
  Tanggal : 2026-09-20
  Status  : MENUNGGU perbandingan isi badan prosedur terhadap basis data berjalan
```

```
C-02 · PENGETAHUAN.md §0 salah menyatakan batas buktinya
  Bunyi   : PENGETAHUAN.md §0 menyatakan DDL tabel lama tidak ada, padahal sudah di repo
            sejak 18 September.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : instruksi ronde 4, catatan N-3; instruksi KETETAPAN-KOMITE bagian 5
  Tanggal : 2026-09-20
  Status  : BERLAKU
```

```
C-03 · Tanda DIPULIHKAN
  Bunyi   : DIPULIHKAN (pengikatan hasil pencocokan nama halaman bind) tidak boleh naik
            menjadi EVIDENCED.
  Alasan  : (tidak disertakan pada sumber)
  Asal    : instruksi ronde 4, aturan pelaporan; instruksi KETETAPAN-KOMITE bagian 5
  Tanggal : 2026-09-20
  Status  : BERLAKU — diulang sebagai aturan kerja no. 1
```

```
C-04 · Dasar penutupan G-07
  Bunyi   : (judul saja: "dasar penutupan G-07 separuh cacat")
  Alasan  : (tidak disertakan pada sumber)
  Asal    : instruksi KETETAPAN-KOMITE bagian 5 — hanya judul
  Tanggal : 2026-09-20
  Status  : MENUNGGU bunyi — lihat bagian 10, BUNYI HILANG
```

```
C-05 · A-1b dibelah dari A-1
  Bunyi   : (judul saja: "A-1b dibelah dari A-1")
  Alasan  : (tidak disertakan pada sumber)
  Asal    : instruksi KETETAPAN-KOMITE bagian 5; pembelahan terlihat pada tabel status
            aliran ronde 2 dan seterusnya (A-1a dan A-1b terdaftar terpisah)
  Tanggal : 2026-09-20
  Status  : BERLAKU — bunyi normatifnya tidak tercatat, lihat bagian 10
```

---

## 6. Sembilan keputusan beku

Disalin utuh dari `PUTUSAN-01.md` bagian 9. Ketetapan asal disebut dalam kurung.

```
1 · Wewenang memutus diperiksa dengan kesetaraan identitas penuh, ditegakkan di backend,
    dan permintaan yang gagal ditolak, bukan ditandai. (K-01)
2 · Keputusan hanya dicatat untuk pihak yang benar-benar bertindak; jenjang yang tidak
    sempat memutus punya keadaan tersendiri yang bukan "tolak". (K-04)
3 · Satu waktu keputusan per jenjang — satu fakta, satu kolom. (K-08)
4 · Peran dan jabatan diambil dari roster sebagai data, lalu disalin ke keputusan pada
    saat diambil sehingga riwayat tidak berubah ketika orangnya pindah jabatan. (K-03, T-06)
5 · Nilai lingkungan, kode akuntansi, alamat layanan, dan daftar penerima surat adalah
    konfigurasi; tidak ada cabang berdasarkan identitas orang di jalur mana pun. (K-09)
6 · Niat memanggil layanan luar dicatat sebelum panggilan, dengan kunci idempotensi per
    efek. (K-10)
7 · Nomor akseptasi diterbitkan sekali, dari nilai yang seluruhnya berpenulis, dan tidak
    pernah disunting sesudah terbit. (K-05, G-05, T-03)
8 · Sirkulasi tanpa jenjang ditolak saat dibuat; tidak ada jalur yang menutup sirkulasi
    tanpa keputusan. (G-17)
9 · Keadaan sirkulasi dan akibatnya pada klaim disimpan terpisah, sehingga "komite setuju"
    pada usulan tutup/tolak tidak terbaca sebagai "klaim disetujui". (G-19, §6)
```

Asal: `PUTUSAN-01.md` bagian 9 · 2026-09-20 · Status: BERLAKU.
Catatan pemakaian: keputusan no. 5 melarutkan pertanyaan terbuka pada F-23; keputusan no. 8 dikuatkan H-6.

---

## 7. Aturan yang mengikat cara kerja

```
A · Empat tanda bukti
  Bunyi   : EVIDENCED (berkas·langkah) · DIPULIHKAN (pengikatan hasil pencocokan nama
            halaman bind — tidak boleh naik jadi EVIDENCED) · TAFSIR · DIPUTUSKAN <nomor>.
  Asal    : instruksi ronde 4 (tiga tanda pertama) + instruksi KETETAPAN-KOMITE bagian 7
            (tanda keempat)
  Status  : BERLAKU
```

```
B · Aturan pagar
  Bunyi   : Setiap pagar dan setiap permintaan pengambilan bukti wajib menunjuk baris
            INVENTARIS-BUKTI.md yang membuktikan buktinya memang belum dipegang — nomor
            bagian dan nama objeknya, bukan kalimat "belum ada". Pagar atau permintaan
            yang tidak menunjuk baris inventaris tidak sah, dan siapa pun yang membacanya
            berhak mengabaikannya sampai rujukannya dilengkapi. Bila rujukan itu ternyata
            menunjuk baris yang isinya ada di repo, pagar itu gugur dengan sendirinya —
            tanpa perlu ronde, tanpa perlu keputusan.
  Asal    : INVENTARIS-BUKTI.md bagian 4
  Status  : BERLAKU
```

```
C · Alasan tunduk pada pagar
  Bunyi   : Pagar berlaku untuk alasan, bukan hanya untuk kesimpulan. Sebuah pertanyaan
            gugur bila salah satu cabangnya hanya dapat dinilai setelah A-1b, A-4, A-5,
            atau isi A-6 cair — termasuk ketika cabang itu tampak masuk akal.
  Asal    : instruksi ronde 3, uji V-6
  Status  : BERLAKU
```

```
D · Paritas terhadap perilaku yang tidak dapat ditentukan
  Bunyi   : Paritas terhadap perilaku yang tidak dapat ditentukan adalah kalimat kosong —
            jalur semacam itu wajib diberi perilaku yang tegas.
  Asal    : instruksi KETETAPAN-KOMITE bagian 7
  Status  : BERLAKU
```

```
E · Pemutakhiran inventaris
  Bunyi   : Inventaris dimutakhirkan setiap kali satu berkas bukti masuk atau satu lubang
            tertutup. Berkas yang tidak dimutakhirkan berhenti menjadi dasar yang sah bagi
            pagar mana pun.
  Asal    : INVENTARIS-BUKTI.md bagian 4
  Status  : BERLAKU
```

```
F · Ketetapan luar yang tetap mengikat
  Bunyi   : ADR-0003 (nilai uang memakai desimal, tidak pernah float64) · ADR-0006 (RBAC
            dirancang dari nol) · ADR-0016 (tanpa database link ke sistem lain) ·
            SPEC-MODEL-DATA §16 (penamaan Indonesia, UPPER_SNAKE_CASE, ≤ 30 byte,
            constraint <peran>_<tabel>[_n]).
  Asal    : konteks tetap ronde 1–4; ddl-usulan dan docs/adr sisi Claim
  Status  : BERLAKU
```

```
G · Zona waktu
  Bunyi   : Asia/Jakarta.
  Asal    : instruksi KETETAPAN-KOMITE bagian 7; perilaku sistem lama tercatat di
            PENGETAHUAN.md §8.2 (@CurrentDate("dd","Asia/Jakarta")) dan
            PROC_GENERATE_SEQUENCE_NUMBER (SYSTIMESTAMP AT TIME ZONE 'Asia/Jakarta')
  Status  : BERLAKU
```

```
H · Pembalikan makna jalur B
  Bunyi   : Di jalur B, AcceptStatus="1" (komite setuju) berarti klaim ditolak/ditutup —
            karena yang disetujui adalah usulan penolakan/penutupan. Salah membaca ini
            akan membalik seluruh logika di sistem baru.
  Asal    : PENGETAHUAN.md §6; ditegakkan D-3 dan keputusan beku no. 9
  Status  : BERLAKU
```

### Tambahan ronde 6

```
Aturan · Pernyataan ketiadaan menyebut bentuk pencariannya, bukan hanya cakupannya.
  Asal    : ronde 5, M5-02
  Tanggal : 2026-09-20
```

```
Aturan · Jawaban atas pertanyaan QF-x dicatat sebagai ketetapan bernomor pada ronde ia
         dijawab, sebelum ronde berikutnya dibuka. Pertanyaan yang jawabannya tidak
         berberkas dianggap belum dijawab.
  Asal    : ronde 6, pola ketujuh (INVENTARIS-BUKTI.md §3)
  Tanggal : 2026-09-20
```

```
Aturan · Status sebuah perkara dikutip dari berkas yang menetapkannya, dengan nomor
         barisnya. Berkas yang hanya menyebutnya tidak sah sebagai sumber status.
  Asal    : ronde 6, pola kedelapan
  Tanggal : 2026-09-20
```

```
Aturan · Sebuah ronde ditutup hanya setelah keluarannya diperiksa ada di disk. Keluaran
         yang masih berbentuk teks di chat berstatus belum terbit, dan tidak sah sebagai
         prasyarat, rujukan, maupun sumber istilah. Pemeriksaannya dilakukan terhadap
         disk, bukan terhadap transkrip.
  Asal    : ronde 6, pola kesembilan
  Tanggal : 2026-09-20
```

```
Aturan · Pada berkas hidup, baris `Ronde` dan `Masukan` di kepala menyatakan pemutakhiran
         terakhir, dan diperbarui tiap pemutakhiran. Itu bukan penyuntingan baris lama
         melainkan bagian dari pemutakhiran itu sendiri. Yang tambah-saja adalah isi,
         bukan kepala.
  Asal    : ronde 6, penerbitan
  Tanggal : 2026-09-20
```

```
Aturan · Grilling dihentikan ketika satu ronde penuh berlalu tanpa satu pun ketetapan
         berubah, tanpa satu pun pagar bergeser, dan tanpa satu pun pertanyaan baru naik
         ke pemilik proses. Ronde berikutnya dibuka hanya oleh bukti yang masuk — sebuah
         kueri pulang, sebuah rule diekspor, sebuah pagar dibuka — bukan oleh kehendak
         memeriksa ulang. Pemeriksaan ulang tanpa bukti baru menghasilkan temuan tentang
         pemeriksaan, bukan tentang sistem.
  Asal    : ronde 6, adjudikasi penutup
  Tanggal : 2026-09-20
```

---

## 8. Daftar fakta F-1 … F-24

| # | Bunyi satu baris | Tanda | Berkas·langkah | Ronde |
|---|---|---|---|---|
| F-1 | Layar komite tidak pernah menyunting rekening penerima; tidak ada penulisan balik rekening di kedua activity pasca | EVIDENCED | `ShowTransfer.xml`; `KomitePostAdjustment.xml`, `KomitePostAdjustmentCWP.xml` | 3 |
| F-2 | Enam medan dapat disunting; empat di antaranya jenjang 1 saja | EVIDENCED | `ShowTransfer.xml` (`pyDisabledWhen = .KomiteCount != 1`) | 3 |
| F-3 | Medan bersufiks `2` bukan urusan Komite (tampilan baca-saja atas data hulu) | EVIDENCED | `ShowTransfer.xml`; nol pembaca di Activity modul Komite | 3 |
| F-4 | Nol SLA, timer, eskalasi, reassign, delegasi di seluruh modul | EVIDENCED | `Flow/KomiteTreaty_Flow.xml` (hanya `pyRouteTo=Custom`) | 3 |
| F-5 | `IsProposeClose` anak mendarat di induk hanya pada persetujuan terakhir non-bersyarat | EVIDENCED | `KomitePostAdjustment.xml` 14.12 | 3 |
| F-6 | 16 kemunculan medan rekening, seluruhnya `pyReadOnly=true` | EVIDENCED | `ShowTransfer.xml` | 3 |
| F-7 | Lima medan induk `pyReadOnly=true` **dan** `pyRequired=true` → G-20 separuh gugur | EVIDENCED | `ShowTransfer.xml` | 3 |
| F-8 | Delapan medan dapat-sunting, bukan enam: tambahan `.CNPFlagXOL` dan satu render `.CurrencyID`, keduanya tanpa kunci jenjang | EVIDENCED | `ShowTransfer.xml` | 3 |
| F-9 | Komite tidak dapat mengubah satu pun angka uang | EVIDENCED | `ShowTransfer.xml` (16 properti nilai, seluruhnya baca-saja) | 3 |
| F-10 | Empat rule pemagar A-1b sudah ada di repo | EVIDENCED | `Claim Non Prop/Activity/` dan `/ReportDefinition/` | 3 |
| F-11 | Sirkulasi berulang atas satu adjustment adalah jalur yang dirancang; ikatan tambahan: tiap sirkulasi wajib punya nomor urut per adjustment, bukan hanya cap waktu | EVIDENCED | `KomitePostAdjustment.xml` 16; `AdjustmentDetailNP.xml` (7×); `ProteksiSendKomiteCNP_Act.xml` (0 kemunculan `IsKomite`) | 3 |
| F-12 | Kendali sasaran `IsOutstanding`/`CNPFlagOuts` tidak dapat ditentukan dari XML (kedekatan posisi bukan bukti) | batas | `ShowTransfer.xml` | 3 |
| F-13 | Jumlah jenjang lahir dari cacah baris Report Definition, bukan dari aturan | EVIDENCED | `CreateChildKomiteCNP_Act.xml` 9, 26.13, 26.14, 26.15 | 4 |
| F-14 | Penyaring roster: `A AND C AND B`; urutan `.DEGREE` menaik; `.LIMIT_TOP` diambil tetapi tidak pernah menyaring | EVIDENCED | `FilterEmailKomiteWithLimit.xml` | 4 |
| F-15 | Tidak ada nilai klaim yang dibandingkan dengan ambang — pembandingnya dimatikan dengan catatan "untuk sementara" | EVIDENCED | `CreateChildKomiteCNP_Act.xml` 26.3–26.6 | 4 |
| F-16 | `Flagkomite`: empat konstanta, satu cabang mati (langkah 14), dan jalur nilai > 50 juta meninggalkan `Param.LIMIT_BOTTOM` tidak terisi | EVIDENCED + TAFSIR (perilaku runtime parameter kosong) | `CreateChildKomiteCNP_Act.xml` 10, 12, 13, 14 | 4 |
| F-17 | Yang dinamai `TotalValueAdjust` bukan total — ia nilai baris terakhir daftar | EVIDENCED | `CreateChildKomiteCNP_Act.xml` 11 | 4 |
| F-18 | Dua daftar diisi berbeda, lalu `.JABATAN` yang benar ditimpa daftar induk berisi `.ID` | EVIDENCED | `CreateChildKomiteCNP_Act.xml` 26.8, 26.13 | 4 |
| F-19 | Satu substitusi orang ditulis di dalam kode: `.OPERATOR_ID=="DARTO"` → `KomiteID="CHRISTINEANGELINA"` | EVIDENCED | `CreateChildKomiteCNP_Act.xml` 26.8.3 | 4 |
| F-20 | Jalur tutup/tolak lahir dari rule lain, tanpa roster dan tanpa ambang; jenis dibedakan `TypeComentAnalysis=="5"` | EVIDENCED | `CreateChildKomiteCloseNP_Act.xml` 6, 8, 9, 12 | 4 |
| F-21 | Gerbang `ProteksiSendKomiteCNP_Act` menandai galat, tidak menghentikan; penghentian terjadi di `@hasMessages` pada langkah pembuatan | EVIDENCED | `ProteksiSendKomiteCNP_Act.xml`; `CreateChildKomiteCNP_Act.xml` 29 | 4 |
| F-22 | Sirkulasi dapat lahir tanpa jenjang; satu-satunya pembatalan adalah alokasi layer kosong | EVIDENCED | `CreateChildKomiteCNP_Act.xml` 28, 29 | 4 |
| F-23 | Perlakuan khusus satu operator menyentuh empat tempat dan menghasilkan `KomiteAproval=""`, bukan `0` | EVIDENCED + batas (kesetaraan `""` dan `0` tidak dapat ditentukan dari XML) | `CreateChildKomiteCNP_Act.xml` 26.8, 26.9, 26.10, 26.15 | 4 |
| F-24 | Syarat tampil pemicu terbaca; peran pemicu tidak — keempat rule A-1b tidak memuat satu pun `pyPrivilege` | EVIDENCED (cacah) + DIPULIHKAN (makna) + batas (peran) | `AdjustmentDetailNP.xml` (7×); empat rule A-1b | 4 |

**Yang melarutkan fakta lain**
- **G-2 melarutkan F-12** — sasaran `pyDisabledWhen` tidak perlu lagi ditentukan; tangkapan layar turun menjadi opsional.
- **Keputusan beku no. 5 melarutkan pertanyaan terbuka pada F-23** — cabang berdasarkan identitas orang tidak ada di jalur mana pun, sehingga kesetaraan `""` terhadap `0` berhenti menjadi perkara.

**Yang mengoreksi fakta lain**
- **F-8 melengkapi F-2** — daftar enam medan menjadi delapan.
- **F-19 dikoreksi pembacaannya di ronde 4** — yang terjadi adalah **substitusi saat pembentukan**, bukan pengalihan sesaat.
- **F-7 mengoreksi G-20** — bagian G-20 tentang aturan lintas agregat gugur; bagian `.Comment` wajib berdiri.
- **F-1 mengoreksi `PENGETAHUAN.md` §10.1** — kalimat "rekening & nilai dikunci, hanya approver pertama boleh menyunting" salah.
- **F-18 menjelaskan K-03** — pemetaan nama orang di `SetKomiteList_Act` adalah gejala dari penimpaan daftar, bukan sebabnya.

---

## 9. Register deviasi diharapkan

| # | Deviasi | Diperintahkan oleh | Uji paritas yang dirancang gagal |
|---|---|---|---|
| 1 | Wewenang memutus memakai kesetaraan identitas penuh dan permintaan yang gagal **ditolak** | D-1, keputusan beku no. 1 (K-01) | Pengguna ber-identitas yang mengandung `KomiteID` sebagai sub-string mengirim keputusan; sistem lama menerima, sistem baru menolak |
| 2 | Peran dan jabatan diambil dari roster, bukan dari nama orang di kode | D-1, keputusan beku no. 4 (K-03), F-18 | Anggota komite yang tidak termasuk lima nama di `SetKomiteList_Act`; sistem lama menghasilkan jabatan kosong/"Claim Admin", sistem baru menghasilkan jabatan roster |
| 3 | Jenjang yang tidak sempat memutus berkeadaan **TIDAK_SAMPAI**, bukan TOLAK | D-1, keputusan beku no. 2 (K-04) | Penolakan di jenjang 2 dari 4; sistem lama mencatat jenjang 3 dan 4 sebagai menolak dengan komentar dan jam milik penolak, sistem baru tidak |
| 4 | Keputusan bersyarat dicatat pada jenjang yang bertugas, bukan pada baris terakhir daftar | D-1 (K-07) | Akseptasi bersyarat berjenjang lebih dari satu; sistem lama menimpa baris `ComiteeClaim(<LAST>)`, sistem baru menulis pada jenjang yang memutus |
| 5 | Satu kolom waktu keputusan per jenjang | D-1, keputusan beku no. 3 (K-08) | Bandingkan dua kolom `DateApproval` dan `DateApprove` sistem lama terhadap satu kolom sistem baru |
| 6 | Sirkulasi tanpa jenjang ditolak saat dibuat | D-1 (G-17), keputusan beku no. 8, H-6, F-22 | Penyaring roster mengembalikan nol baris; sistem lama tetap membuat case dan menutupnya diam-diam, sistem baru menolak pembuatan |
| 7 | `.CNPFlagXOL` dikunci ke jenjang 1; `.CurrencyID` baca-saja seluruhnya | G-2 | Jenjang 2 mengubah `.CNPFlagXOL`; sistem lama menerima, sistem baru menolak |
| 8 | Nilai penentu kelas kewenangan adalah nilai adjustment yang **sedang diajukan** | D-1, F-17 | **Wajib memakai klaim dengan ≥ 2 adjustment beda kelas** — pada klaim ber-adjustment tunggal deviasi ini lolos tanpa terlihat |
| 9 | Jabatan diambil dari roster; `.JABATAN` tidak ditimpa `.ID` | D-1, F-18 | Bandingkan `IDKomite` pada case Komite: sistem lama berisi pengenal baris roster, sistem baru berisi jabatan |
| 10 | Jalur `VINCENTVERNANDO_1` tidak dimigrasikan | H-7, keputusan beku no. 5, F-23 | Sirkulasi yang dibuat operator tersebut; sistem lama menghasilkan satu jenjang ber-`KomiteAproval=""`, sistem baru mengikuti aturan umum |
| 11 | Cabang mati `Flagkomite` langkah 14 ditutup | H-1, D-1, F-16 | `RNMShare` di antara 15 dan 30 dengan nilai ≤ 50 juta; sistem lama tidak pernah menyalakan cabang, sistem baru menyalakannya |
| 12 | Lubang di atas 50 juta ditutup dengan kelas yang tegas | H-1, aturan kerja D, F-16 | Adjustment bernilai > 50 juta; sistem lama menjalankan penyaring dengan parameter tak terisi, sistem baru memakai kelas yang ditetapkan |
| 13 | Validasi pembuatan **menolak**, bukan menandai | H-5, F-21 | Pembuatan sirkulasi dengan rekening kosong; sistem lama membuat pesan lalu bergantung pada `@hasMessages`, sistem baru menolak kasus-guna |

---

## 10. Yang hilang dan yang bentrok

### 10.1 Hilang

| Hal | Keterangan |
|---|---|
| **C-04 — dasar penutupan G-07 separuh cacat** | Hanya judulnya yang tercatat. Bunyi normatifnya tidak ada di sumber mana pun. `BUNYI HILANG` |
| **C-05 — A-1b dibelah dari A-1** | Pembelahannya terlihat pada tabel status aliran sejak ronde 2, tetapi kalimat normatif yang memerintahkannya tidak pernah tercatat. `BUNYI HILANG` |
| **Alasan H-1 … H-7** | Enumerasi ronde 4 memuat *apa* yang diputuskan, tanpa *mengapa*. Tujuh baris kolom Alasan kosong. `BUNYI HILANG` |
| **Kalimat normatif penuh H-1 … H-7** | Yang tersedia adalah enumerasi ringkas dipisah titik-tengah, bukan kalimat normatif utuh seperti D, E, dan G. Dicatat apa adanya sesuai larangan parafrase |
| **G-01 … G-20 (temuan interogasi)** | `GRILL-01.md` dihapus 2026-09-20 atas permintaan. Bunyi penuh kedua puluh temuan hilang; yang tersisa hanya klaim yang dikutip ulang di `PUTUSAN-01.md` bagian 3 dan sidang bagian 8. Nomor G-01…G-20 karena itu **bentrok deret** dengan G-1…G-4 ronde 3 — lihat 10.2 |
| **K-01 … K-12 (temuan `PENGETAHUAN.md` §11)** | Masih utuh di `PENGETAHUAN.md`. Tidak hilang; dicatat di sini hanya karena dirujuk D-1 tanpa disalin |
| **`PENILAIAN-01.md`** | Tidak pernah ada. Seluruh isi yang dirujuk dengannya hidup hanya sebagai judul di instruksi |
| **Uji paritas untuk tiap deviasi** | Hanya uji F-17 yang pernah dirancang eksplisit ("klaim dengan ≥ 2 adjustment beda kelas"). Dua belas uji lain pada bagian 9 disusun juru catat dari bunyi deviasinya, bukan disalin dari sumber — **tandai sebagai belum diratifikasi** |
| **Jawaban ronde 1–3 sebagai berkas** | Tidak pernah tersimpan; dipulihkan ke dalam berkas ini dari transkrip. Bila transkrip hilang sebelum berkas ini disimpan, D, E, dan G hilang bersamanya |

### 10.2 Bentrok

| Hal | Sumber A | Sumber B |
|---|---|---|
| **Deret nomor `G`** | `PUTUSAN-01.md` memakai `G-01 … G-20` untuk temuan interogasi | Ronde 3 memakai `G-1 … G-4` untuk ketetapan. Dua deret berbeda dengan huruf yang sama; keduanya dirujuk D-1 dan F-… tanpa pembeda |
| **Cakupan F-2** | Ronde 3: "Enam medan dapat disunting" | F-8: delapan. Instruksi ronde 4 menyebut F-8 "melengkapi" F-2, tetapi bunyi F-2 tidak pernah dicabut |
| **Dasar faktual D-4** | D-4: "Fase 1 memakai prosedur lama apa adanya" (BERLAKU) | C-01: isi badan `PROC_GENERATE_SEQUENCE_NUMBER` belum diverifikasi terhadap basis data berjalan; lima penutupan yang bersandar padanya bertanda `TAFSIR (menunggu C-01)` |
| **Status A-1b** | Ronde 2 dan 3: A-1b **BEKU**, "dibuka oleh ekspor 2 rule" | F-10 dan `INVENTARIS-BUKTI.md` §2.4: keempat rule sudah di repo sebelum pagarnya dipasang. Ronde 4 menyatakan A-1b "TERBUKA SETELAH PEMBACAAN"; kalimat pencabutan pagar tidak pernah ditulis |
| **Pagar pengalihan** | E-3: "aksi pengalihan tidak diekspos **sampai A-1b cair**" | H-4: "pengalihan sesaat tetap tidak diekspos" — A-1b sudah cair, pagar tetap berlaku atas dasar lain yang tidak dinyatakan |
| **Sumber "zona Asia/Jakarta"** | Instruksi KETETAPAN-KOMITE bagian 7 mencantumkannya sebagai aturan kerja | Tidak ada ronde yang pernah memutuskannya; yang ada hanya perilaku sistem lama di `PENGETAHUAN.md` §8.2 dan di badan prosedur Oracle |

---

*Berkas ini tidak memutuskan apa pun. Ia memindahkan keputusan yang sudah diambil menjadi daftar bernomor yang dapat dirujuk. Setiap baris menyebut asalnya; yang tidak dapat menyebut asalnya ada di bagian 10, bukan di bagian 1–7.*
---

## 11. Ketetapan ronde 5 — K5-1 … K5-6

> Bunyi disalin apa adanya dari `GRILL-05/06-PUTUSAN.md` bagian 5. Dua baris
> terakhir tiap blok ditambahkan di sini, bukan di sumbernya.

```
K5-1 · Satu penentu jenjang aktif
  Bunyi  : Jenjang aktif sebuah sirkulasi ditentukan oleh satu hal saja, yaitu jenjang
           berderajat terendah yang belum memutuskan. Hitungan jenjang adalah turunan
           yang dihitung dari roster, bukan penyimpan tandingan, dan tidak boleh
           dipakai sebagai syarat penugasan maupun syarat penutupan.
  Alasan : Sistem lama memakai dua penentu yang tak pernah saling memeriksa (N-01), dan
           penutupan case bergantung pada perbandingan hitungan yang dipalsukan saat
           penolakan (T5-01).
  Asal   : ronde 5, N-01 dan T5-01
  Tanggal: 2026-09-20
  Status : BERLAKU
```

```
K5-2 · Pembuatan sirkulasi dan akibatnya adalah satu transaksi
  Bunyi  : Pembuatan sirkulasi, penulisan nomornya ke klaim, penyimpanan klaim, dan
           pengiriman pemberitahuan berhasil bersama atau gagal bersama. Tidak ada
           langkah sesudah pembuatan yang boleh berjalan ketika pembuatan tidak jadi.
  Alasan : Penjaga `@hasMessages` di sistem lama hanya melewati satu langkah; langkah
           sesudahnya menulis nomor komite milik sirkulasi lain lalu menyimpannya (N-07).
  Asal   : ronde 5, N-07 — memperluas H-5
  Tanggal: 2026-09-20
  Status : BERLAKU
```

```
K5-3 · Kegagalan pembentukan wajib terlihat
  Bunyi  : Setiap jalur pembentukan sirkulasi yang berakhir tanpa sirkulasi berakhir
           dengan kegagalan yang terlihat oleh pengguna dan tercatat. Penanda yang tidak
           dibaca siapa pun bukan penanganan galat.
  Alasan : Ketika daftar spreading kosong, sistem lama keluar diam-diam dan menulis
           `IsFlagError` yang tidak dibaca oleh satu pun rule dalam ekspor (N-08).
  Asal   : ronde 5, N-08
  Tanggal: 2026-09-20
  Status : BERLAKU
```

```
K5-4 · Pembentukan roster bersifat idempoten
  Bunyi  : Membentuk roster sebuah sirkulasi menyusun ulang isinya dari awal. Pembentukan
           kedua atas perkara yang sama menghasilkan roster yang sama, bukan roster yang
           bertambah panjang.
  Alasan : Pada jalur bersyarat, sistem lama melewati pembersihan lalu menambahkan,
           sehingga jumlah jenjang tidak lagi sama dengan jumlah anggota (N-06).
  Asal   : ronde 5, N-06
  Tanggal: 2026-09-20
  Status : BERLAKU
```

```
K5-5 · Tidak ada kombinasi tanpa aturan seleksi
  Bunyi  : Tabel seleksi jenjang menerima dua masukan — nilai adjustment dan bagian
           treaty — dan setiap kombinasi keduanya menghasilkan satu aturan. Kombinasi
           yang tidak tercakup adalah galat konfigurasi, bukan roster kosong.
  Alasan : Ketika bagian melebihi 30 persen atau nilai melebihi 50 juta, sistem lama
           tidak menetapkan ambang sama sekali dan memanggil laporan dengan parameter
           yang tidak pernah diisi (N-09).
  Asal   : ronde 5, N-09 — melengkapi H-1 dan H-6
  Tanggal: 2026-09-20
  Status : BERLAKU
```

```
K5-6 · Penanda dua-keadaan adalah satu boolean
  Bunyi  : Penanda yang hanya punya dua keadaan disimpan sebagai satu boolean dengan satu
           ejaan pada seluruh lapisan. Nilai kosong dan nilai nol tidak pernah berbagi
           satu kolom, dan tidak ada kolom yang artinya bergantung pada apakah ia kosong
           atau berisi "0".
  Alasan : `IsCloseFile` ditulis angka pada anak dan teks pada induk dalam satu langkah
           yang sama (N-05); `KomiteAproval` ditulis kosong sementara pembandingnya
           menguji nol (N-11).
  Asal   : ronde 5, N-05 dan N-11
  Tanggal: 2026-09-20
  Status : BERLAKU
```

```
K5-7 · Maksud dan akibat adalah dua penanda berbeda
  Bunyi  : Penanda "diajukan untuk ditutup" dan "diajukan untuk ditolak" dicatat pada saat
           pengajuan dan tidak pernah menjadi penanda "ditutup" atau "ditolak". Yang kedua
           hanya lahir dari keputusan komite, lewat fungsi akibat D-3. Sirkulasi yang gagal
           lahir atau ditolak meninggalkan maksudnya tercatat dan akibatnya tidak.
  Alasan : CreateChildKomiteCloseNP_Act langkah 6 menulis IsCloseFile dan IsReject ke induk
           pada langkah yang sama yang membuat halaman anak — sebelum komite memutuskan apa
           pun (T5-03) — sehingga klaim membawa penanda itu bahkan ketika sirkulasinya gagal
           lahir (N-07).
  Asal   : ronde 5, T5-03 dan N-07 — melengkapi D-3
  Tanggal: 2026-09-20
  Status : BERLAKU  (dicatat ronde 6)
```

### Akibat ronde 5 terhadap ketetapan lama

| Ketetapan | Akibat | Dasar |
|---|---|---|
| `H-5` | **DIPERLUAS** oleh `K5-2` — yang menolak bukan hanya langkah pembuatan, melainkan seluruh urutan yang bergantung padanya | GRILL-05 N-07, sidang F-21 |
| `H-1` | **DILENGKAPI** oleh `K5-5` — tabel seleksi menerima dua masukan, dan tidak ada kombinasi tanpa aturan | GRILL-05 N-09 |
| `H-6` | **DIKUATKAN** — sebab sirkulasi tanpa jenjang kini berbukti langkah, bukan dugaan | GRILL-05 N-09, sidang F-22 |
| `H-7` | **DIPAKAI MENUTUP** N-11 tanpa perlakuan | GRILL-05 06-PUTUSAN §2 |
| `D-5` | **DIPERLUAS** — pemetaan orang→jabatan yang harus pindah ke data bertambah empat baris yang ditulis di dalam rule | sidang F-19 |
| `D-2` | **DIPERLUAS** — batas transaksi menutupi kedua activity pembuat anak, bukan hanya activity keputusan | GRILL-05 T5-04 |
| `PUTUSAN-01.md` §8, perlakuan `K-02` | **NAIK DARI TAFSIR MENJADI BERBUKTI**, lalu diserap `K5-1` | GRILL-05 N-01 |

### Tambahan bagian 10.1 — yang hilang, per ronde 5

| Hal | Keterangan |
|---|---|
| **Arti nilai `.Type` dan `FlagProrate`** | Keduanya terbukti dipakai; artinya tidak tertulis di ekspor karena `Rule-Obj-FieldValue` tidak diekspor. `BUNYI HILANG` sampai B5-2 pulang |
| **Penaik `komiteAccept_ticket`** | Ticket menempel pada gerbang `KomiteLoop`; tak satu pun dari 338 berkas menaikkannya. Pemiliknya di luar ekspor |

### Tambahan bagian 10.2 — yang bentrok, per ronde 5

| Hal | Sumber A | Sumber B |
|---|---|---|
| **Urutan giliran jenjang** | `PUTUSAN-01.md` §8: "baris pertama yang belum memutuskan", ditulis sebagai perlakuan tanpa membuka rule-nya | `KomiteRouter.xml` langkah 6.1 (`WhenTrue=6`/`WhenFalse=6`): sama persis. **Bentrok selesai** — sumber B membuktikan sumber A |
| **Penentu jenjang aktif** | `KomiteRouter.xml` langkah 1–4: `.KomiteCount` | `KomiteRouter.xml` langkah 6: baris pertama ber-`KomiteAproval==0`. Keduanya hidup di satu rule; diselesaikan oleh `K5-1` |

---

## 12. Ketetapan ronde 6 — K6-1 … K6-3

```
K6-1 · Kelas kewenangan: dua kelas tetap, tanpa angka dari deskripsi
  Bunyi  : Dua kelas kewenangan dipertahankan; ambang berbasis nilai tidak dipulihkan
           sekarang. Angka 15 tidak dipakai dan tidak diajukan lagi — menaikkan komentar
           atau deskripsi langkah menjadi ambang kewenangan adalah bentuk lain dari
           menaikkan DIPULIHKAN menjadi EVIDENCED. Yang naik ke pemilik proses tinggal satu
           pertanyaan, dan K5-5 menyelesaikan sisanya tanpa angka.
  Alasan : Jawaban atas QF-1. Cabang langkah 14 CreateChildKomiteCNP_Act berbunyi
           `>30 && <=30` dan tidak pernah benar; angka 15 hanya ada pada deskripsi langkah,
           tidak pada kode mana pun.
  Asal   : ronde 6, jawaban QF-1 (rekomendasi penulis adalah (b) dengan batas bawah 15;
           yang diputuskan adalah (b) tanpa angka, ditopang K5-5)
  Tanggal: 2026-09-20
  Status : BERLAKU
```

```
K6-2 · Hak baca sirkulasi diturunkan dari hak baca klaim
  Bunyi  : Hak baca sirkulasi diturunkan dari hak baca klaim; ia bukan peran tersendiri di
           RBAC. Paritas, J-4 dikukuhkan.
  Alasan : Jawaban atas QF-2. Komentar komite sudah menjadi bagian berkas klaim di sistem
           lama; membatasinya menghapus konteks dari orang yang hari ini sudah membacanya,
           dan menjadikannya peran tersendiri menambah permukaan RBAC tanpa kebutuhan.
  Asal   : ronde 6, jawaban QF-2 (sejalan rekomendasi (a))
  Tanggal: 2026-09-20
  Status : BERLAKU
```

```
K6-3 · Bersyarat adalah masukan ketiga pada tabel seleksi
  Bunyi  : QF-3 dilebur ke QF-1. Bersyarat menjadi masukan ketiga pada tabel seleksi K5-5,
           di samping nilai dan bagian treaty. Sebabnya: langkah 26.6 membawa deskripsi yang
           sama persis dengan langkah 26.3 — ambang bersyarat adalah sisa penonaktifan yang
           sama, bukan kebijakan terpisah. Isinya hari ini mereproduksi perilaku sekarang.
  Alasan : Jawaban atas QF-3. Kesamaan deskripsi kedua sub-langkah menunjukkan keduanya
           lahir dari satu penonaktifan, sehingga penimpaan ambang bersyarat bukan aturan
           yang berdiri sendiri.
  Asal   : ronde 6, jawaban QF-3 (rekomendasi penulis adalah (a) sebagai kebijakan; yang
           diputuskan adalah melebur ke K5-5 sebagai masukan ketiga, dengan isi awal yang
           mereproduksi perilaku sekarang)
  Tanggal: 2026-09-20
  Status : BERLAKU — menutup N-10, yang GRILL-05/07-AUDIT tandai TERBUKA
```

```
K6-4 · Jalur tutup dan tolak klaim
  Bunyi  : Sirkulasi TUTUP_KLAIM dan TOLAK_KLAIM tetap satu jenjang. Pemegangnya
           diselesaikan dari roster berdasarkan peran, bukan dari nama yang ditulis di
           dalam aturan. Bila tidak ada baris roster aktif yang memegang peran itu,
           sirkulasi ditolak saat dibuat. Jenis sirkulasi dinyatakan oleh pemanggil,
           tidak diturunkan dari nilai kolom analisis.
  Alasan : Menutup klaim adalah keputusan yang membatalkan, bukan yang mengeluarkan uang,
           sehingga kedalamannya tidak disamakan dengan jalur pembayaran. Nama orang,
           surel, inisial, dan dua sebutan jabatan yang berbeda untuk satu kedudukan
           ditulis di dalam rule (F-20) — melanggar keputusan beku no. 5.
  Asal   : ronde 4 Q-2, bunyi hilang; ditetapkan ulang ronde 6 dari F-20
  Tanggal: 2026-09-20
  Status : BERLAKU — menggantikan H-2, yang tetap tercatat BUNYI HILANG
```

### Akibat ronde 6 terhadap ketetapan lama

| Ketetapan | Akibat | Dasar |
|---|---|---|
| keputusan beku no. 1 (`K-01`) | **DIKUKUHKAN**, mekanismenya diperjelas: pemeriksaan wewenang memasang pesan dan tidak menghentikan apa pun | `GRILL-06/01-PEMBACAAN.md` P6-3 |
| `D-1` (`K-07`) | **DIKUKUHKAN**, sasarannya diperjelas: jalur bersyarat menulis ke halaman yang berbeda, bukan hanya ke baris yang berbeda | P6-4 |
| `G-03`, `G-08` | **DIKUKUHKAN**; tidak ada ketetapan yang berubah | P6-1, P6-2 |
| `N-10` | **DITUTUP** oleh `K6-3` | jawaban QF-3 |
| `QF-1`, `QF-2`, `QF-3` | seluruhnya **DIJAWAB** dan bernomor | `K6-1`, `K6-2`, `K6-3` |

### Tambahan bagian 10.1 — status enam bunyi `H`, per ronde 6

Baris "Kalimat normatif penuh `H-1 … H-7`" pada bagian 10.1 **tidak dihapus**; ia diberi
status per nomor. `H-4` sudah pulih penuh di bagian 4.

| # | Status baru | Diserap oleh |
|---|---|---|
| `H-1` | **RINGKAS — MENGIKAT APA ADANYA, diserap** | `K6-1` (dua kelas, pemulihan tidak dilakukan) dan `K5-5` (aturan seleksi sebagai data, cabang mati dan lubang di atas 50 juta) |
| `H-2` | tetap **`BUNYI HILANG`** — digantikan `K6-4`, bukan ditemukan | tidak terserap di mana pun. `CONTEXT.md` bagian 1 hanya membakukan istilah *jenis sirkulasi*; ia tidak memuat aturan "satu jenjang dengan pemegang dari roster berdasarkan peran". Ketiga ADR tidak menyebutnya |
| `H-3` | **RINGKAS — MENGIKAT APA ADANYA, diserap** | `CONTEXT.md` bagian 5, istilah *delegasi tetap*, yang merujuk `H-3`; dan bunyi penuh `H-4` yang bersandar padanya |
| `H-5` | **RINGKAS — MENGIKAT APA ADANYA, diserap** | `K5-2` dan `ADR-0031` |
| `H-6` | **RINGKAS — MENGIKAT APA ADANYA, diserap** | keputusan beku no. 8 dan `K5-5` (kombinasi tak tercakup adalah galat konfigurasi); bentuk penolakannya oleh `K5-3` |
| `H-7` | **RINGKAS — MENGIKAT APA ADANYA, diserap** | bagian 9 baris 10 dan keputusan beku no. 5 |

### Tambahan bagian 10.2 — yang bentrok, per ronde 6

| Hal | Sumber A | Sumber B |
|---|---|---|
| **Status `G-03`** | Dikutip berulang sebagai `RAGU`, "diturunkan atas dasar kemustahilan bisnis" | `PUTUSAN-01.md` baris 57: `G-03` ada pada daftar "Diterima tanpa perubahan". Yang berstatus `RAGU` adalah bacaan gabungan `G-03`+`G-04` (baris 31). **Bentrok selesai** — sumber B menang, dan kekeliruan kutipan dicatat sebagai pola kedelapan |
