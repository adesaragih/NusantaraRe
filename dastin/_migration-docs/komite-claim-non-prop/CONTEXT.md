> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : juru catat
> Masukan: KETETAPAN.md (2026-09-20, termasuk K6-1…K6-4) · PENGETAHUAN.md · FAKTA-A1B.md · GRILL-05/* · GRILL-06/01-PEMBACAAN.md · penyuntingan penutup ronde 6
> Status : TERBUKA
> Sifat  : HIDUP

# CONTEXT — glosarium Komite Claim Non Prop

Berkas ini **hanya** glosarium. Tidak ada rancangan, tidak ada nama tabel, tidak ada aturan.
Ia menentukan **akar kata** yang dipakai seluruh modul; nama tabel dan kolom tunduk
`SPEC-MODEL-DATA §16` (Indonesia, `UPPER_SNAKE_CASE`, ≤ 30 byte) dan dibentuk dari akar kata
di sini.

Istilah yang **berbeda artinya** di sistem lama dan baru diberi baris terpisah. Peleburan
dua arti ke dalam satu istilah adalah cara tercepat kehilangan keduanya.

---

## 1. Sirkulasi dan bentuknya

| Istilah | Artinya di modul ini | Padanan sistem lama | **Bukan** ini |
|---|---|---|---|
| **Sirkulasi** | Satu berkas yang mengedarkan **satu usulan** kepada sejumlah jenjang untuk diputuskan | case anak `ASM-FW-GCNMFW-Work-KomiteTreatyNonProp`, dibuat lewat `pxAddChildWork` | Bukan klaim. Bukan komite sebagai lembaga. Bukan daftar orang |
| **Jenis sirkulasi** | Sifat yang menentukan bagaimana jenjangnya dibentuk: **usulan pembayaran** atau **usulan tutup/tolak** | `TransferType`; bernilai `'2'` untuk jalur banyak-jenjang, dan **tidak pernah diisi** pada jalur tutup/tolak | Bukan `.Type` mentah dari adjustment — itu sumbernya, bukan istilahnya |
| **Jenjang** | Satu kedudukan di dalam sirkulasi yang berhak memutus, dengan derajat dan pemegangnya | satu baris `KomiteList` pada sirkulasi, dan satu baris `ComiteeClaim` pada adjustment induk | Bukan orang. Seorang dapat memegang jenjang pada banyak sirkulasi |
| **Derajat** | Urutan giliran sebuah jenjang di dalam sirkulasi, diambil dari roster saat pembentukan | `.DEGREE`, satu-satunya kolom berurutan pada `FilterEmailKomiteWithLimit` (F-14) | Bukan pangkat, bukan kelas kewenangan, dan bukan nomor urut yang dihitung ulang setelah pembentukan |
| **Jenjang aktif** | Jenjang berderajat terendah yang belum memutuskan (`K5-1`) | dua penentu yang tak saling memeriksa: `KomiteCount` dan baris pertama ber-`KomiteAproval==0` | Bukan "jenjang ke-n". Nomor urut adalah turunan, bukan penentu |
| **Pemegang** | Orang yang saat ini berhak memutus atas sebuah jenjang | `KomiteID`, berisi nama pengguna | Bukan pemilik jenjang selamanya — pemegang dapat berganti lewat delegasi tetap |
| **Roster** | Daftar acuan berisi siapa menduduki jenjang apa, dengan derajat dan kelas kewenangannya | tabel `EMAILKOMITE`, disaring `FilterEmailKomiteWithLimit` | Bukan daftar jenjang sebuah sirkulasi — itu salinan yang diambil dari roster saat pembentukan |

## 2. Yang diedarkan

| Istilah | Artinya di modul ini | Padanan sistem lama | **Bukan** ini |
|---|---|---|---|
| **Usulan** | Apa yang diminta diputuskan: satu penyesuaian nilai klaim, atau satu permintaan menutup/menolak klaim | satu baris `ClaimData.AdjustmentList(n)`, atau `TempCommiteClaim` pada jalur tutup/tolak | Bukan klaim. Satu klaim dapat melahirkan banyak usulan |
| **Versi usulan** | Keadaan usulan pada satu titik waktu, tersimpan sehingga dapat dibandingkan | **tidak ada padanan** — suntingan menimpa nilai sebelumnya tanpa jejak | Bukan revisi bernomor yang dapat dipilih pemutus |

## 3. Memutus, dan apa yang mengikutinya

Tiga istilah ini adalah tempat kesalahpahaman paling mahal di modul ini. Dipisah tegas atas
perintah `D-3`.

| Istilah | Artinya di modul ini | Padanan sistem lama | **Bukan** ini |
|---|---|---|---|
| **Keputusan** | Tindakan satu jenjang atas satu usulan, beserta komentarnya dan waktunya | `KomiteAproval` + `KomiteComment` + `DateApproval`/`DateApprove` pada baris jenjang | Bukan hasil sirkulasi — satu sirkulasi memuat banyak keputusan |
| **Hasil** | Putusan komite **atas usulan**: usulan diterima atau usulan ditolak | `AcceptStatus`, `"1"` terima dan `"2"` tolak | **TIDAK PERNAH** berarti nasib klaim. Ini larangan `D-3`, dan ia ada karena sistem lama mencampurnya |
| **Akibat** | Apa yang terjadi pada klaim karena sebuah hasil — ditentukan oleh jenis sirkulasi, bukan oleh hasilnya sendiri | tidak punya nama; tersebar di jalur A dan jalur B `KomitePostAdjustment` | Bukan sinonim hasil. Pada usulan tutup/tolak, **hasil "diterima" berarti klaim ditutup atau ditolak** |
| **Fungsi akibat** | Pemetaan (jenis, bersyarat, hasil) → akibat pada klaim, ditulis sekali di satu tempat | tidak ada padanan — tersebar di jalur A dan jalur B `KomitePostAdjustment` | Bukan hasil. Lihat *hasil* dan *akibat* |
| **Maksud** | Apa yang **diajukan** pengaju sebelum komite memutus: "diajukan untuk ditutup", "diajukan untuk ditolak" (`K5-7`) | `IsCloseFile` dan `IsReject`, ditulis ke klaim pada saat pengajuan | Bukan akibat. Maksud tercatat meski sirkulasi gagal lahir |
| **`TIDAK_SAMPAI`** | Keadaan sebuah jenjang yang tidak pernah mendapat giliran karena sirkulasi berakhir lebih dulu | **tidak ada padanan** — sistem lama menuliskan penolakan atas nama jenjang itu, lengkap dengan komentar dan jam milik penolak | Bukan "menolak". Bukan "belum memutuskan" — sirkulasi sudah selesai |

## 4. Akseptasi dan pembalikan

| Istilah | Artinya di modul ini | Padanan sistem lama | **Bukan** ini |
|---|---|---|---|
| **Akseptasi** | Pencatatan resmi bahwa sebuah usulan yang disetujui menjadi kewajiban bernomor | baris pada `OS_AKSEPTASI_KLAIM`, dengan `AcceptedNo` | Bukan keputusan. Bukan hasil. Akseptasi lahir **sesudah** hasil, dan hanya untuk hasil "diterima" |
| **Akseptasi bersyarat** | Akseptasi yang terbit dengan syarat yang belum terpenuhi | `IsSubjectivity`, disimpan terpisah di `OS_AKSEPTASI_SUBJECTIVITY` | Bukan hasil ketiga. `J-3` menetapkannya **atribut usulan**, bukan jenis putusan |
| **Pembalikan** | Pencatatan yang meniadakan akseptasi yang sudah terbit, tanpa menghapusnya | baris bernilai negatif, `STS_REJECT` | Bukan pembatalan. Yang dibalik tetap ada dan tetap terbaca |
| **Nomor akseptasi** | Pengenal resmi sebuah akseptasi, terbit sekali dan tidak pernah dipakai ulang | `AcceptedNo`, diterbitkan `PROC_GENERATE_SEQUENCE_NUMBER` | Bukan nomor sirkulasi. Bukan nomor klaim |
| **Periode buku** | Bulan pembukuan yang menentukan deret nomor akseptasi | `TANGGAL_CLOSING`, dibaca prosedur penerbit; di layar muncul sebagai `TglProd` | Bukan tanggal keputusan. Bukan tanggal hari ini |

## 5. Wewenang

| Istilah | Artinya di modul ini | Padanan sistem lama | **Bukan** ini |
|---|---|---|---|
| **Bagian treaty** | Persentase bagian perusahaan atas sebuah treaty, salah satu masukan tabel seleksi kelas kewenangan | `TreatyInMaster.RNMShare`, dibaca `CreateChildKomiteCNP_Act` langkah 10 ke `Local.CekLimitPersen` (EVIDENCED) | Bukan nilai uang. Bukan kelas kewenangan — ia masukan bagi kelas, bukan kelasnya |
| **Tabel seleksi** | Data yang menentukan kelas kewenangan dari tiga masukan: nilai usulan, bagian treaty, dan bersyarat | tidak ada padanan — lahir dari `K5-5`, `K6-1`, `K6-3`; di sistem lama tersebar sebagai `Flagkomite` dan empat konstanta di dalam langkah | Bukan roster. Tabel seleksi memilih **kelas**; roster memuat **orang** |
| **Aturan seleksi** | Satu baris tabel seleksi: satu kombinasi masukan, satu kelas kewenangan | tidak ada padanan | Bukan penyaring roster — itu yang menerapkan hasilnya |
| **Galat konfigurasi** | Keadaan ketika sebuah kombinasi masukan tidak tercakup aturan seleksi mana pun | tidak ada padanan — sistem lama memanggil penyaring dengan parameter tak terisi | **Bukan roster kosong**, dan bukan sirkulasi tanpa jenjang. Ini cacat data acuan, bukan hasil pembentukan |
| **Kelas kewenangan** | Kelompok jenjang mana yang harus memutus sebuah usulan, ditentukan oleh masukan yang dinyatakan pada tabel seleksi | `Flagkomite` bernilai 0/1/2, diterjemahkan menjadi `Param.LIMIT_BOTTOM` | Bukan jabatan. Bukan pangkat. Satu orang dapat berada di banyak kelas |
| **Delegasi tetap** | Pencatatan pada roster bahwa seorang memegang jenjang milik orang lain sampai dicabut (`H-3`) | **tidak ada padanan** — ditambal satu baris kode yang menukar satu nama dengan nama lain | Bukan pengalihan sesaat |
| **Pengalihan sesaat** | Memindahkan satu giliran yang sedang berjalan kepada orang lain, sekali pakai | **tidak ada padanan, dan tidak diekspos** (`H-4`) | Bukan delegasi tetap. Modelnya menampung, permukaannya tidak ada |

## 6. Dua nama yang nyaris sama — dibedakan tegas

| Nama sistem lama | Menempel pada | Isinya | Jangan tertukar karena |
|---|---|---|---|
| `ComiteeClaim` | **adjustment** (`ClaimData.AdjustmentList(n).ComiteeClaim`) | daftar jenjang **per usulan**; inilah yang dibaca jalur bersyarat | ejaannya salah dengan cara yang berbeda dari yang satunya |
| `ClaimComitee` | **klaim** (`ClaimData.ClaimComitee`) | salinan daftar jenjang **per klaim**, ditulis sekali saat pembentukan | keduanya hanya berbeda urutan dua kata, dan keduanya salah eja |

Pada sistem baru keduanya **tidak** dibawa sebagai nama. Keduanya adalah *jenjang* yang
menempel pada *sirkulasi*; yang kedua adalah turunan tampilan, bukan penyimpan.

## 7. Kata yang dilarang dipakai sebagai istilah

| Kata | Sebab dilarang |
|---|---|
| "approve" / "reject" tanpa objek | tidak menyatakan apakah yang dimaksud usulan atau klaim; `D-3` ada persis untuk ini |
| "komite" sebagai nama tabel | ambigu antara lembaga, roster, sirkulasi, dan jenjang |
| "loop" / "count" | `KomiteLoop` dan `KomiteCount` adalah dua penyimpan yang `K5-1` gabungkan menjadi satu turunan; membawa namanya membawa kembali perselisihannya |
| "flag" | `Flagkomite`, `CNPFlagXOL`, `FlagProrate`, `IsFlagError` — empat hal tak berhubungan dengan satu awalan |
