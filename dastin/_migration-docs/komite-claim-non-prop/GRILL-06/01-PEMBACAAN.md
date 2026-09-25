> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : pembaca
> Masukan: HitServiceToKasirKMT_Act.xml · InsertXOLKlaimCNP.xml · KomitePostAdjustment.xml · PUTUSAN-01.md · PENGETAHUAN.md §8.4 dan §11
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

## 0. Perbaikan alat, dilaporkan lebih dulu

Alat baca diperbaiki **sebelum** satu berkas pun dibaca. Dua perubahan:

1. Tiap pra-syarat dan transisi dicetak dengan penambat langkahnya (`@10.3`, `@6.1`,
   `@14.2.2.1.1.3`) dan sandi mentahnya di samping terjemahan; transisi tidak lagi
   disembunyikan ketika kedua cabangnya bernilai `2`.
2. Kata "iterasi-ulang" diganti `keluar(iterasi bila langkah ada di dalam loop, activity
   bila tidak)`, karena sandi `6` berarti dua hal bergantung kedudukan langkahnya.

Penambat langkah itu sendiri yang menutup `M5-01`: keluaran kini berbunyi
`@6.1 TRANSISI [true] ya:6=keluar-iterasi tidak:6=keluar-iterasi` dan
`@7 TRANSISI [(tanpa syarat)] ya:2=lanjut tidak:2=lanjut` — dua baris yang sebelumnya
tercampur dalam satu jendela teks.

---

## P6-1 · `G-03` — **DIKUKUHKAN**

```
Yang dibaca ulang : HitServiceToKasirKMT_Act.xml langkah 10.3
                    @10.3 PRA-SYARAT [.TreatyName=="UR"] ya:2=lanjut tidak:3=lewati-step
Pemilik pasangan  : langkah 10.3, dipastikan tiga cara — (a) penambat @10.3 dicetak dari
                    nomor langkah yang sama yang mencetak isinya, bukan dari kedekatan
                    baris; (b) langkah 10.4 punya pasangannya sendiri
                    (@10.4 [IsPEGASyariah] ya:2 tidak:3), sehingga keduanya tidak dapat
                    tertukar; (c) blok pembungkusnya, langkah 10, punya dua pra-syarat
                    terpisah yang tercetak sebagai @10, bukan @10.3
Bacaan lama       : PENGETAHUAN.md §8.4 — baris `UR` DILEWATI pada jalur Kasir
Bacaan baru       : blok muatan berjalan HANYA ketika `.TreatyName=="UR"`; baris selain UR
                    melewati langkah 10.3 seluruhnya. Bacaan §8.4 terbalik — persis seperti
                    yang G-03 catat
Uji kenyataan     : lolos. Bacaan ini tidak menyiratkan sistem mustahil berjalan: ia hanya
                    menyatakan baris mana yang mengisi halaman muatan, bukan apakah
                    kiriman pernah terjadi
Akibat            : G-03 tetap SAHIH. PENGETAHUAN.md §8.4 tetap dicabut
Tanda             : EVIDENCED HitServiceToKasirKMT_Act.xml langkah 10.3
```

**Instruksi ronde ini mengandaikan `G-03` akan runtuh. Ia tidak runtuh.**

Andaian itu juga bersandar pada premis yang keliru. Instruksi menyatakan `G-03` diturunkan
dari `SALAH` menjadi `RAGU`. `PUTUSAN-01.md` baris 57 justru mencantumkan `G-03` pada
daftar **"Diterima tanpa perubahan"**. Yang diturunkan ke `RAGU` adalah **bacaan gabungan
`G-03`+`G-04`** — "muatan pembayaran mungkin tidak pernah tersusun sampai JSON-nya cacat"
(baris 31) — bukan `G-03` sendiri. `G-03` tidak pernah punya status `RAGU` untuk dipulihkan.

**`PG-04` karena itu tidak menyusut.** Satu fakta struktural yang terbaca ronde ini justru
membuatnya lebih perlu, bukan kurang: `Connect-REST` berada di langkah **10.7**, di dalam
loop yang sama, dan pra-syaratnya hanya `IsPEGAPROD` — bukan `.TreatyName=="UR"`. Penulisan
log pun, langkah 10.9, tidak berpra-syarat sama sekali. Berapa kiriman yang benar-benar
terjadi per nomor akseptasi, dan apa isinya, tetap perkara `A-4` dan tetap dijawab oleh satu
`SELECT` atas log kiriman. Fakta strukturnya dicatat di sini; **artinya tidak disimpulkan**.

---

## P6-2 · `G-08` — **DIKUKUHKAN, dengan sandinya diperjelas**

```
Yang dibaca ulang : InsertXOLKlaimCNP.xml langkah 1
                    @1 PRA-SYARAT [pyWorkPage.IsSubjectivity==true] ya:6 tidak:2=lanjut
                    dan langkah 4.1
                    @4.1 PRA-SYARAT [.TreatyName=="UR"] ya:3=lewati-step tidak:2=lanjut
Pemilik pasangan  : langkah 1 — langkah tingkat atas, bukan sub-langkah; dipastikan dari
                    penambat @1 dan dari ketiadaan blok `-- sub-steps --` di bawahnya.
                    Karena langkah 1 TIDAK berada di dalam loop, sandi 6 di sini berarti
                    KELUAR ACTIVITY, bukan keluar iterasi. Pemastian ketiga datang dari
                    rule itu sendiri: memo versinya berbunyi
                    "Exit Activity If IsSubjectivity==true"
Bacaan lama       : rincian layer tidak pernah tertulis untuk akseptasi bersyarat
Bacaan baru       : sama. Activity berhenti di langkah 1 ketika `IsSubjectivity==true`,
                    sehingga langkah 2 sampai akhir tidak berjalan
Uji kenyataan     : lolos. Ini penghentian yang disengaja dan diberi nama oleh pembuatnya
                    di memo versi; sistem berjalan bertahun-tahun karena jalur bersyarat
                    punya tujuan simpan tersendiri
Akibat            : G-08 tetap SAHIH. Tidak ada ketetapan, pagar, atau deviasi yang berubah
Tanda             : EVIDENCED InsertXOLKlaimCNP.xml langkah 1
```

Satu catatan yang berhenti di batas pagar: langkah 4.1 memakai pra-syarat yang sama dengan
langkah 10.3 pada P6-1 — `.TreatyName=="UR"` — dengan **polaritas terbalik**
(`ya:3=lewati-step` di sini, `ya:2=lanjut` di sana). Kedua rule itu milik `A-5` dan `A-4`.
Apakah keduanya dua sisi dari satu pembagian yang disengaja adalah pertanyaan tentang isi
kedua aliran beku itu, **dan karena itu tidak dijawab di sini, juga tidak dipakai sebagai
alasan**. Ia dicatat sebagai bahan pembuka `PG-04`/`PG-05`, bukan sebagai kesimpulan.

---

## P6-3 · `K-01` — **DIKUKUHKAN, dan diperjelas satu tingkat**

```
Yang dibaca ulang : KomitePostAdjustment.xml langkah 3 dan 4
                    @3 PRA-SYARAT [@contains(@toUpperCase(pyWorkPage.KomiteList(Local.
                       IdxKomite).KomiteID),@toUpperCase(OperatorID.pyUserIdentifier))]
                       ya:3=lewati-step tidak:2=lanjut
                    @4 PRA-SYARAT [ekspresi yang sama] ya:3=lewati-step tidak:2=lanjut
Pemilik pasangan  : langkah 3 dan langkah 4 masing-masing punya pasangannya sendiri, dengan
                    ekspresi identik; dipastikan dari dua penambat terpisah @3 dan @4 yang
                    dicetak bersama isi langkahnya. Langkah 4 adalah Page-Set-Messages dan
                    tidak punya `set` apa pun, sehingga tidak dapat tertukar dengan langkah 3
Bacaan lama       : kontrol wewenang memakai pencocokan sub-string
Bacaan baru       : sama, dan arahnya kini tegas. `@contains(KomiteID, pyUserIdentifier)` —
                    yang dicari adalah identitas pengguna DI DALAM pengenal komite. Bila
                    cocok, langkah 3 dan 4 DILEWATI, yaitu tidak ada galat. Bila tidak
                    cocok, keduanya berjalan dan memasang pesan "Invalid User Acceptance!"
Uji kenyataan     : lolos, dan justru menjelaskan mengapa sistem berjalan bertahun-tahun —
                    kontrol ini tidak pernah menghentikan apa pun
Akibat            : K-01 tetap SAHIH. Keputusan beku no. 1 tidak berubah. Deviasi 1 pada
                    REGISTER-DEVIASI.md berdiri lebih kokoh: uji "permintaan yang gagal
                    DITOLAK" memang dirancang gagal, sebab sistem lama tidak menolak
Tanda             : EVIDENCED KomitePostAdjustment.xml langkah 3 dan 4
```

Yang bertambah dari pembacaan dengan penambat: langkah 3 dan 4 hanya **memasang pesan**.
Tidak ada `Exit-Activity`, tidak ada lompatan, dan langkah 5 (`Obj-Open-By-Handle`) berjalan
tanpa memeriksa apa pun. Pemeriksaan wewenang di sistem lama bukan gerbang, melainkan
catatan. Ini penguatan `K-01`, bukan temuan terpisah.

---

## P6-4 · `K-07` — **DIKUKUHKAN, dan sasarannya kini tepat**

```
Yang dibaca ulang : KomitePostAdjustment.xml langkah 6 dan 7
                    @6 PRA-SYARAT [TempMainWork.ClaimData.AdjustmentList(Local.
                       IdxAdjustment).IsSubjectivity == true] ya:3=lewati-step tidak:2=lanjut
                    @7 PRA-SYARAT [ekspresi yang sama]        ya:2=lanjut tidak:3=lewati-step
Pemilik pasangan  : langkah 6 dan 7, masing-masing tercetak dengan penambatnya sendiri dan
                    dengan empat baris `set` yang berbeda sasarannya — langkah 6 menulis ke
                    pyWorkPage.KomiteList(Local.IdxKomite), langkah 7 menulis ke
                    TempMainWork.ClaimData.AdjustmentList(...).ComiteeClaim(<LAST>).
                    Perbedaan sasaran itu sendiri memastikan tidak ada pasangan yang
                    tertukar
Bacaan lama       : pada akseptasi bersyarat, keputusan mendarat di baris yang salah
Bacaan baru       : sama, dan kini lengkap. Kedua langkah SALING MENIADAKAN: langkah 6
                    berjalan ketika BUKAN bersyarat, langkah 7 ketika bersyarat. Pada jalur
                    bersyarat, keputusan ditulis ke baris TERAKHIR `ComiteeClaim` milik
                    adjustment induk, dan roster sirkulasi sendiri
                    (`pyWorkPage.KomiteList`) TIDAK MENERIMA APA-APA
Uji kenyataan     : lolos dengan satu ketegangan, dicatat di 07-AUDIT bagian akhir
Akibat            : K-07 tetap SAHIH. Deviasi 4 pada REGISTER-DEVIASI.md diperjelas:
                    yang salah bukan hanya indeks barisnya, melainkan juga halamannya
Tanda             : EVIDENCED KomitePostAdjustment.xml langkah 6 dan 7
```

---

## 5. Ringkasan

| Temuan | Hasil | Yang berubah |
|---|---|---|
| `G-03` | **DIKUKUHKAN** | tidak ada ketetapan berubah; premis instruksi dikoreksi; `PG-04` **tidak** menyusut |
| `G-08` | **DIKUKUHKAN** | tidak ada |
| `K-01` | **DIKUKUHKAN** | mekanismenya diperjelas: memasang pesan, tidak menghentikan |
| `K-07` | **DIKUKUHKAN** | sasarannya diperjelas: halaman berbeda, bukan hanya baris berbeda |

**Nol dari empat runtuh.** Misatribusi `M5-01` tidak berulang pada satu pun dari keempatnya.
Yang membedakan `M5-01` dari keempat ini: pada `M5-01` pasangan transisi dibaca dari
jendela `sed` sepanjang 160 baris yang memuat dua langkah; keempat pembacaan di atas dibaca
dari keluaran yang menambatkan tiap pasangan pada nomor langkahnya.

Satu koreksi yang ikut lahir, dan bukan tentang XML: status `G-03` selama ini dikutip
sebagai `RAGU`. `PUTUSAN-01.md` baris 57 menyatakan `G-03` "diterima tanpa perubahan".
Yang berstatus `RAGU` adalah bacaan gabungan `G-03`+`G-04`. Kutipan yang keliru itu masuk
ke instruksi ronde ini sebagai premis. Ia dicatat di `INVENTARIS-BUKTI.md` §3 sebagai pola
ketujuh, bukan disembunyikan.
