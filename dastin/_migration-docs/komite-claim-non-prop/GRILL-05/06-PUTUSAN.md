> Modul  : Komite Claim Non Prop · Ronde 05 · 2026-09-20
> Peran  : penilai
> Masukan: `00-LINGKUP.md` · `01-TEMUAN.md` · `02-SIDANG.md` · `03-LUBANG.md` · `04-PARITAS.md` · `05-GERBANG.md` — seluruhnya ronde 5 · bukti mentah pada 338 XML
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

# 06 · PUTUSAN

Ditulis dengan peran berbeda. Dakwaan ronde ini dicurigai sekeras dakwaan itu mencurigai
`PENGETAHUAN.md`. Lima temuan berubah, satu di antaranya berubah seluruhnya; empat temuan
turunan ditambahkan; dua koreksi metodologis dijatuhkan atas cara interogasi itu sendiri.

---

## 1. Koreksi metodologis

### M5-01 · Satu pasang transisi dibaca dari jendela teks, bukan dari langkahnya

Temuan N-01 versi pertama menyatakan router memberi giliran kepada baris **terakhir** yang
belum memutuskan. Dasarnya sepasang nilai `pyStepsTransParamsWhenTrue=2` dan
`…WhenFalse=2` yang terbaca dalam potongan 160 baris XML. Potongan itu memuat **dua**
langkah — langkah 6 beserta sub-langkahnya dan langkah 7 — dan pasangan `2/2` itu milik
langkah 7. Transisi sub-langkah 6.1 yang sebenarnya adalah `6/6`, yaitu keluar iterasi,
sehingga yang dipilih adalah kecocokan **pertama**.

Kekeliruan itu tidak tertangkap oleh pembaca langkah, karena `dump_act.py` menampilkan
sandi `6` sebagai "iterasi-ulang" pada kedua cabang — kata yang dapat dibaca sebagai
"ulangi iterasi" maupun "sudahi iterasi".

**Perbaikan yang berlaku mulai sekarang.** Pertama, tiap klaim yang bergantung pada sandi
transisi wajib dibaca dengan penambat langkahnya, bukan dengan jendela baris. Kedua,
`dump_act.py` menampilkan sandi transisi apa adanya di samping terjemahannya. Ketiga —
dan ini yang menyelamatkan perkara ini — uji kenyataan dijalankan **sebelum** temuan
dikirim, bukan sesudahnya: pembalikan urutan wewenang pada jalur utama akan terlihat oleh
penggunanya sejak hari pertama, dan itulah yang memicu pembacaan ulang.

### M5-02 · Ketiadaan diumumkan dari satu bentuk pencarian saja

Tiga temuan berdiri di atas hasil nol: N-13 (nol pemanggil lain), N-14 (nol penaik ticket),
dan pernyataan bahwa `IsFlagError` tidak dibaca siapa pun (N-08). Ketiganya dicari dengan
**satu** bentuk: pencarian teks atas nama. Di Pega, sebuah activity dapat dipanggil lewat
`pyActivityName` yang dirakit dari nilai properti, sebuah ticket dapat dinaikkan dari
Service Level atau dari `pyFlowActionPreProcessing`, dan sebuah properti dapat dibaca lewat
ekspresi yang merangkai namanya.

Ini **tidak** membatalkan ketiganya — bentuk perakitan nama tidak terlihat di mana pun pada
338 berkas — tetapi mengubah kekuatannya. N-13 dan N-14 sudah bertanda `TAK TERVERIFIKASI`
dan `sebagian`; pernyataan `IsFlagError` pada N-08 belum, dan diturunkan di sini menjadi
"tidak dibaca oleh satu pun rule dalam ekspor", bukan "tidak dibaca siapa pun".

**Perbaikan yang berlaku mulai sekarang.** Pernyataan ketiadaan menyebut bentuk
pencariannya, bukan hanya cakupannya.

### M5-03 · Pagar dipatuhi pada kesimpulan, hampir tidak pada alasan

`03-LUBANG.md` mencatat L5-3 — baris akseptasi mungkin menyimpan urutan jenjang — dan
menahan diri. Itu benar. Tetapi `01-TEMUAN.md` N-04 menyebut isi ekspresi surat
(`@if(childPageKomite.TransferType = "3","Reject","Close")`) di dalam badan temuan sebelum
memindahkannya ke L5-1. Menyebut ekspresinya bukan pelanggaran — itu rule, bukan isi
dokumen — tetapi jaraknya terhadap pagar A-6 tipis, dan ronde berikutnya tidak akan
seberuntung ini.

**Perbaikan yang berlaku mulai sekarang.** Ekspresi yang bermuara ke isi dokumen dicatat di
daftar lubang bertahan sejak awal, bukan dipindahkan setelah ditulis.

---

## 2. Adjudikasi

Hanya temuan yang **berubah** putusan atau tingkatnya.

### N-01 · SALAH P1 → **SAHIH P2**, dengan klaim yang diuji diganti

Klaim asli ("giliran jatuh pada baris terakhir") gugur oleh M5-01. Yang tersisa dan
terbukti adalah perkara lain di rule yang sama: dua penentu jenjang aktif hidup
berdampingan — `.KomiteCount` pada langkah 1–4 dan "baris pertama yang belum memutuskan"
pada langkah 6 — tanpa satu pun memeriksa yang lain. Itu sahih, dan tingkatnya P2 karena
keduanya sejalan selama roster tidak berubah dan keputusan tidak melompat urutan.

Satu akibat penting ikut berubah arah: perlakuan K-02 pada `PUTUSAN-01.md` §8 **bukan**
salah kaprah melainkan benar, dan kini berbukti langkah.

### N-02 · SAHIH P1 → **SAHIH P2**

Dakwaan aslinya: empat nama keranjang membatasi jenjang pada empat. Pemeriksaan keadaan
menunjukkan nilai dari langkah 1–4 hanya bertahan bila langkah 6 dan langkah 7 keduanya
tidak menulis `param.AssignTo`. Langkah 7 **selalu** menulis pada jalur `!='2'`. Langkah 6
menulis kapan pun ada baris belum memutuskan. Maka keempat nama itu hanya terpakai pada
keadaan "jalur banyak-jenjang, semua baris sudah memutuskan" — keadaan yang gerbang
`IsKomiteLoop` sudah akhiri.

Yang tersisa nyata: empat nama keranjang tertulis di dalam rule. Itu keluarga K-09, sudah
diputus oleh keputusan beku no. 5, dan tingkatnya P2.

### N-10 · KEPUTUSAN BISNIS P2 → **KEPUTUSAN BISNIS P1**

Menimpa ambang untuk klaim bersyarat bukan perkara tampilan; ia menentukan **siapa yang
berwenang menyetujui** klaim bersyarat di atas 25 juta. D-1 menempatkan kebenaran wewenang
pada kelas "perbaikan wajib, paritas ditolak". Sebuah keputusan bisnis yang menentukan
wewenang tidak boleh menumpang di P2 bersama ketidakjelasan.

Putusannya tetap `KEPUTUSAN BISNIS` — bentuknya memang bentuk aturan — tetapi tingkatnya
naik, dan ia berpindah ke daftar pemilik proses sebagai perkara yang memblokir tabel
seleksi H-1, bukan sebagai catatan kaki.

### N-11 · RAGU P1 → **DITUTUP, tingkat dicabut**

Interogator sudah menuliskan "LARUT — H-7" pada kolom cara menutup, lalu tetap
mencantumkan tingkat P1. Sebuah perkara yang larut tidak punya tingkat. H-7 menetapkan
jalur `VINCENTVERNANDO_1` tidak dimigrasikan; perbandingan `""` terhadap `0` pada jalur itu
karena itu tidak pernah terjadi di sistem baru. Ditutup tanpa perlakuan.

Yang **tidak** ikut tertutup, dan harus dinyatakan agar tidak hilang: nilai kosong dan nol
tidak boleh berbagi satu kolom. Itu naik menjadi bagian K5-6.

### N-13 · TAK TERVERIFIKASI P2 → **P3, hentikan penyelidikan**

Pertanyaan "apakah `KomitePostAdjustmentCWP` punya pintu masuk lain" tidak mengubah apa pun
pada model sasaran. Jalur tutup/tolak adalah cabang di dalam kasus-guna keputusan baik ia
punya satu pemanggil maupun tiga. Menaikkan biaya untuk membuktikan ketiadaan di ruleset
yang tidak diekspor adalah biaya tanpa hasil.

Diturunkan ke P3 dan penyelidikannya dihentikan. Q-6 ditutup.

---

## 3. Temuan turunan

Yang interogator lewatkan, ditemukan saat memeriksa dakwaannya.

### T5-01 · Penutupan case bergantung pada hitungan, bukan pada keputusan

`KomitePostAdjustment.xml` langkah 20 berlabel `EXIT` berpra-syarat
`pyWorkPage.AcceptStatus==2` dan mengerjakan `pyWorkPage.KomiteCount = pyWorkPage.KomiteLoop`.
Langkah 28 memanggil `ASMForceCaseClose` dengan pra-syarat
`pyWorkPage.KomiteCount>=pyWorkPage.KomiteLoop`. Langkah 29 menambah hitungan satu.

Artinya penolakan tidak menutup case secara langsung; ia **memalsukan hitungan** agar
gerbang penutupan terpenuhi. Interogator membaca langkah 29 untuk N-02 dan melewati langkah
20 dan 28 seluruhnya.

Akibat yang ikut tertarik: setiap keadaan yang membuat `KomiteLoop` keliru — dan N-06
membuatnya keliru pada klaim bersyarat — mengenai **dua** hal sekaligus, yaitu kapan
lingkar berhenti dan apakah case tertutup. Model sasaran harus menutup sirkulasi dari
keadaan keputusan, bukan dari perbandingan dua bilangan.

### T5-02 · `.AcceptStatus` dikosongkan tanpa syarat oleh router

`KomiteRouter.xml` langkah 5 mengerjakan `.AcceptStatus = ""` dan **tidak** berpra-syarat.
Gerbang `When/IsKomiteLoop.xml` menuntut `.AcceptStatus = "1"`.

Maka keadaan keputusan hanya bernilai "1" pada jendela antara keputusan ditulis dan router
berjalan untuk giliran berikutnya. Interogator mencatat bahwa **label** langkah 5 keliru
(N-12) dan melewatkan bahwa **isi** langkah 5 adalah bagian yang menentukan dari mekanika
lingkar.

Akibat yang ikut tertarik: keadaan keputusan pada sistem lama bersifat sesaat dan tidak
dapat dibaca ulang. Sirkulasi pada model sasaran harus menyimpan keputusan per jenjang
sebagai catatan tetap, bukan sebagai satu medan yang dipakai bergantian.

### T5-03 · Klaim ditandai tutup/tolak pada saat pengajuan, bukan pada saat keputusan

`CreateChildKomiteCloseNP_Act.xml` langkah 6 — satu langkah yang sama yang membuat halaman
anak — juga menulis ke **induk**: `pyWorkPage.ClaimData.Remark_Close`, `.NamePIC`,
`.IsCloseFile`, dan `.IsReject`. Langkah 12 menambahkan
`pyWorkPage.CNPStatusCase = "COMITEE ACCEPTANCE (DEPT. HEAD)"`, dan langkah 13 menyimpan.

Artinya klaim sudah membawa penanda tutup atau tolak sebelum komite memutuskan apa pun.
Inversi arti jalur B sudah dikenal sejak `PENGETAHUAN.md`; yang baru adalah **titik
waktunya**, dan itu mengubah rancangan: penanda pada klaim bukan hasil keputusan komite,
melainkan pernyataan maksud pengaju.

Akibat yang ikut tertarik: bila sirkulasi kemudian gagal lahir (N-07) atau ditolak, klaim
tetap membawa penanda itu. Model sasaran harus memisahkan "diajukan untuk ditutup" dari
"ditutup".

### T5-04 · Pola `Obj-Refresh-And-Lock` sesudah penyuntingan muncul lagi, di berkas lain

`CreateChildKomiteCNP_Act.xml` menjalankan `Obj-Refresh-And-Lock pyWorkPage` pada langkah 3
**dan** pada langkah 30, sementara di antaranya `Obj-Save ChildWorkPage` (langkah 17) dan
`Obj-Save pyWorkPage` (sub-langkah 26.17) sudah berjalan. `CreateChildKomiteCloseNP_Act.xml`
memakai pola yang sama pada langkah 1 dan 11.

Ini keluarga K-11, yang `PUTUSAN-01.md` §8 sudah tutup sebagai "bukan temuan". Perkara itu
**tidak dibuka kembali**; yang dicatat hanyalah bahwa cakupannya lebih luas dari satu
berkas, sehingga ketetapan batas transaksi (D-2) harus menutupi kedua pembuat anak, bukan
hanya activity keputusan.

---

## 4. Konsolidasi akar

| Akar | Bunyi | Menghimpun |
|---|---|---|
| R5-1 | Satu fakta disimpan di dua tempat yang tak pernah saling memeriksa | N-01, N-05, T5-01, F-18 |
| R5-2 | Penjaga menjaga satu langkah, bukan akibat yang bergantung padanya | N-07, N-08, T5-02 |
| R5-3 | Aturan wewenang ditulis sebagai tetapan di dalam langkah | N-02, N-09, N-10, F-19, F-20 |
| R5-4 | Pembentukan menambah alih-alih menyusun ulang | N-06, F-13, T5-03 |
| R5-5 | Arti sebuah nilai hidup di luar ekspor | N-04, N-13, N-14, N-15 |

---

## 5. Ketetapan baru yang lahir

Siap dipindahkan ke `KETETAPAN.md` deret `K5`. Bunyi di bawah ini adalah bunyi yang
dipindahkan, bukan ringkasannya.

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
```

```
K5-2 · Pembuatan sirkulasi dan akibatnya adalah satu transaksi
  Bunyi  : Pembuatan sirkulasi, penulisan nomornya ke klaim, penyimpanan klaim, dan
           pengiriman pemberitahuan berhasil bersama atau gagal bersama. Tidak ada
           langkah sesudah pembuatan yang boleh berjalan ketika pembuatan tidak jadi.
  Alasan : Penjaga `@hasMessages` di sistem lama hanya melewati satu langkah; langkah
           sesudahnya menulis nomor komite milik sirkulasi lain lalu menyimpannya (N-07).
  Asal   : ronde 5, N-07 — memperluas H-5
```

```
K5-3 · Kegagalan pembentukan wajib terlihat
  Bunyi  : Setiap jalur pembentukan sirkulasi yang berakhir tanpa sirkulasi berakhir
           dengan kegagalan yang terlihat oleh pengguna dan tercatat. Penanda yang tidak
           dibaca siapa pun bukan penanganan galat.
  Alasan : Ketika daftar spreading kosong, sistem lama keluar diam-diam dan menulis
           `IsFlagError` yang tidak dibaca oleh satu pun rule dalam ekspor (N-08).
  Asal   : ronde 5, N-08
```

```
K5-4 · Pembentukan roster bersifat idempoten
  Bunyi  : Membentuk roster sebuah sirkulasi menyusun ulang isinya dari awal. Pembentukan
           kedua atas perkara yang sama menghasilkan roster yang sama, bukan roster yang
           bertambah panjang.
  Alasan : Pada jalur bersyarat, sistem lama melewati pembersihan lalu menambahkan,
           sehingga jumlah jenjang tidak lagi sama dengan jumlah anggota (N-06).
  Asal   : ronde 5, N-06
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
```

---

## 6. Apa yang tidak saya setujui, dan mengapa itu penting

Menyetujui semuanya berarti meneruskan, bukan mengadili. Yang ditolak ronde ini: satu
temuan P1 yang salah seluruhnya (N-01 versi pertama), satu penaikan tingkat yang tidak
tertanggung (N-02), satu tingkat yang melekat pada perkara yang sudah larut (N-11), dan
satu penyelidikan yang tidak akan mengubah apa pun (N-13).

Satu di antaranya — N-01 — akan menjadi ketetapan yang membalik urutan wewenang seluruh
modul seandainya lolos. Itu ukuran yang sebenarnya dari ronde ini.
