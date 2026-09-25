> Modul  : Komite Claim Non Prop · Ronde 05 · 2026-09-20
> Peran  : interogator
> Masukan: `01-TEMUAN.md` · `02-SIDANG.md` · `KETETAPAN.md` (D-1, E-*, J-*, H-*, sembilan keputusan beku)
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

# 04 · SKENARIO PARITAS

Kolom **keluaran yang diamati** berisi perilaku sistem lama sebagaimana **terbaca dari
rule**, bukan sebagaimana terlihat dari data produksi — nol baris data produksi dipegang
(`INVENTARIS-BUKTI.md` §2.5). Di mana perilakunya menuntut data untuk dipastikan, itu
dinyatakan pada kolom terakhir.

Uji yang menguji deviasi **dirancang gagal**. Itu bukan cacat rancangan uji; itu cara
deviasi dibuktikan disengaja.

---

| ID | Skenario | Masukan | Keluaran yang diamati (sistem lama) | Diperiksa di | Sifat |
|---|---|---|---|---|---|
| P5-01 | Giliran pertama pada sirkulasi tiga jenjang | roster 3 baris, urut `.DEGREE ASC`, semua `KomiteAproval=0`, `TransferType='2'` | penugasan jatuh ke **baris ke-1** | penugasan aktif sesudah pembentukan | PARITAS |
| P5-02 | Giliran sesudah satu persetujuan | keadaan P5-01, baris ke-1 menyetujui | penugasan jatuh ke **baris ke-2** | penugasan aktif | PARITAS |
| P5-02b | Persetujuan di luar urutan | roster 3 baris; baris ke-2 memutuskan lebih dulu | penugasan jatuh ke baris ke-1, sedangkan `.KomiteCount` menunjuk jenjang lain | penugasan aktif vs hitungan jenjang | **DEVIASI DIHARAPKAN** |
| P5-03 | Sirkulasi berjenjang lima | roster 5 baris, `TransferType='2'`, semua belum memutuskan | giliran tetap berjalan lewat baris roster; empat nama keranjang tidak terpakai | penugasan aktif; daftar keranjang kosong | PARITAS |
| P5-04 | Jalur tutup klaim | `TypeComentAnalysis != "5"` | satu jenjang, pemegang tunggal, `IsCloseFile` terisi, `IsReject` kosong | roster sirkulasi; penanda pada klaim | PARITAS |
| P5-05 | Jalur tolak klaim | `TypeComentAnalysis == "5"` | satu jenjang; pada sistem lama anak membawa `IsCloseFile="0"` dan induk membawa kosong | penanda pada klaim dan pada sirkulasi | **DEVIASI DIHARAPKAN** |
| P5-06 | Klaim bersyarat diproses dua kali | klaim `IsSubjectivity=true`, dibentuk ulang | roster memuat anggota ganda; jumlah jenjang lebih kecil dari jumlah anggota | jumlah baris roster vs jumlah jenjang | **DEVIASI DIHARAPKAN** |
| P5-07 | Pembentukan gagal karena pesan validasi | keadaan yang memunculkan pesan pada halaman anak | sirkulasi tidak lahir; klaim tetap menerima nomor komite dari sirkulasi sebelumnya; klaim tersimpan; surat terkirim | nomor komite pada klaim; ada/tidaknya sirkulasi | **DEVIASI DIHARAPKAN** |
| P5-08 | Adjustment tanpa `SpreadingRisk` | `SpreadingRisk` kosong | tidak ada sirkulasi, tidak ada pesan, klaim tersimpan dengan penanda yang tak dibaca siapa pun | keluaran layar; ada/tidaknya sirkulasi | **DEVIASI DIHARAPKAN** |
| P5-09 | Bagian treaty di atas 30 persen | `RNMShare > 30`, nilai berapa pun | roster kosong; sirkulasi lahir tanpa jenjang | jumlah jenjang pada sirkulasi baru | **DEVIASI DIHARAPKAN** |
| P5-10 | Nilai adjustment antara 30 dan 50 juta, bagian ≤ 30 persen | `RNMShare ≤ 30`, nilai 40 juta | roster kelas kedua (ambang bawah 25.000.001) | isi roster | PARITAS |
| P5-11 | Nilai adjustment di atas 50 juta, bagian ≤ 30 persen | `RNMShare ≤ 30`, nilai 60 juta | roster kosong; sirkulasi lahir tanpa jenjang | jumlah jenjang | **DEVIASI DIHARAPKAN** |
| P5-12 | Klaim bersyarat bernilai di atas 25 juta | `IsSubjectivity=true`, nilai 40 juta | roster terluas (ambang bawah 0), bukan roster kelas kedua | isi roster | PARITAS — tertunda Q5-1 |
| P5-13 | Dua adjustment dalam satu klaim | `AdjustmentList` 2 baris, nilai 10 juta dan 40 juta | kelas komite ditentukan oleh nilai baris **terakhir**, bukan totalnya | isi roster | **DEVIASI DIHARAPKAN** |
| P5-14 | Jabatan anggota komite di luar empat nama yang dipaku | anggota roster dengan `OPERATOR_ID` lain | jabatan berasal dari kolom `.ID` basis data | jabatan pada roster sirkulasi | PARITAS |
| P5-15 | Jabatan anggota komite untuk salah satu dari empat nama | `OPERATOR_ID = CHRISTOPMARHASAK` | jabatan ditimpa menjadi "Technic Div. Head"/"CM" | jabatan pada roster sirkulasi | **DEVIASI DIHARAPKAN** |
| P5-16 | Hak membuka sirkulasi | pengguna mana pun yang dapat membuka klaimnya | tidak ada hak yang diperiksa — nol `pyPrivilegeName` berisi di 59 rule | keputusan izin | **DEVIASI DIHARAPKAN** |

---

## Catatan per sifat

**Tujuh belas skenario, sembilan di antaranya deviasi.** Itu bukan tanda ronde ini terlalu
berani; itu akibat langsung D-1, yang menetapkan bahwa di mana grilling menghasilkan cacat
ber-`EVIDENCED`, perbaikan menjadi default dan paritas yang harus dibela. Sembilan deviasi
di atas seluruhnya berdiri di atas langkah yang dibaca, bukan di atas dugaan.

**P5-01, P5-02, dan P5-02b adalah satu berkas uji.** Dua yang pertama menegaskan urutan
giliran yang memang dipertahankan — baris pertama yang belum memutuskan, urut `.DEGREE ASC`.
Yang ketiga memasang keadaan yang membuat kedua penentu jenjang aktif berselisih (N-01) dan
karena itu dirancang gagal: sistem baru hanya punya satu penentu, sehingga tidak dapat
menirukan selisihnya.

**P5-09 dan P5-11 tampak sama tetapi tidak.** Keduanya menghasilkan roster kosong lewat
`Flagkomite=0`, tetapi lewat pintu berbeda: P5-09 lewat bagian treaty, P5-11 lewat nilai.
Sebuah perbaikan yang hanya menutup salah satunya akan meloloskan satu uji dan menggagalkan
satu lagi — dan itu persis yang ingin diketahui.

**P5-12 belum berstatus tetap.** Ia ditandai PARITAS sekarang karena belum ada ketetapan
yang membalikkannya; bila Q5-1 dijawab (b), ia berpindah ke deviasi. Tidak ditandai
"tertunda" pada kolom sifat agar tabel tetap dapat dijalankan apa adanya; penundaannya
tercatat di kolom yang sama.

**P5-16 diuji meski jawabannya sudah pasti.** ADR-0006 menetapkan RBAC dibangun dari nol,
sehingga sistem baru pasti memeriksa hak. Uji ini tetap ada agar tercatat bahwa yang
dibandingkan bukan "hak lama vs hak baru" melainkan "tanpa hak vs berhak" — dan agar tidak
ada yang kelak mencari hak lama untuk ditiru.

**Nol skenario menyentuh A-4, A-5, dan isi A-6.** Skenario pembayaran, pembalikan, dan isi
dokumen tidak dirancang ronde ini, bukan karena tidak menarik, melainkan karena
pembandingnya berada di balik pagar dan uji tanpa pembanding adalah kalimat kosong.
