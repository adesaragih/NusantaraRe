> Modul  : Treaty In Adjustment · ronde C
> Dibuat : 2026-09-24
> Sifat  : skenario paritas — **bentuk selisih ditulis SEBELUM run**, syarat ADR-0043 yang tidak
>          boleh dilewati. Ini bukan hasil; ini harapan yang dapat gagal.
> Status : terbuka sampai ada run paritas

# Paritas ronde C

ADR-0043 menuntut: *"tuliskan lebih dulu, per cacat yang sudah dikenali, **bentuk selisih apa** yang
seharusnya ia hasilkan — arahnya, besarannya, dan populasi terdampak. Tulis itu **sebelum** run
dijalankan. Tanpa itu, pembacanya akan merasionalisasi apa pun yang muncul."*

Tiga putusan ronde ini menghasilkan selisih yang **dapat diramalkan bentuknya**. Satu menghasilkan
selisih yang **harus dikecualikan dari kriteria identik**, dan itu yang paling mudah salah dibaca.

---

## PC-1 — `GRL-14`: selisih EGNPI pada addendum premi berubah dari NOL menjadi BUKAN NOL

| | |
|---|---|
| **Populasi** | seluruh versi warisan ber-`EDMState = 3` yang premi aktualnya berbeda dari kontrak dasarnya |
| **Arah** | sistem lama: **nol baris selisih EGNPI**. Kalkulator baru: **ada baris**, sebesar (premi aktual − premi dasar) |
| **Besaran** | persis selisih yang selama ini tersimpan di `ActualValue.EGNPI` terhadap `OLDDATA.EGNPI` — dan **belum pernah dibandingkan oleh siapa pun**, karena sistem lama tidak pernah menghitungnya |
| **Kelas cacat** | **TDA-16** |

> **Selisih ini WAJIB muncul.** Bila run paritas menghasilkan nol untuk populasi ini, yang salah
> bukan datanya melainkan **kalkulator barunya** — ia berarti premi aktual tidak ikut terbaca sebagai
> nilai versi, dan `GRL-14` tidak benar-benar terpasang.

Ini satu-satunya kelas di daftar ini yang **kegagalannya berbentuk ketiadaan selisih**, dan karena
itu ia paling mudah lolos: nol terbaca seperti "cocok".

---

## PC-2 — `GRL-16`: selisih share fakultatif muncul untuk pertama kalinya

| | |
|---|---|
| **Populasi** | versi warisan yang bagian fakultatifnya berbeda dari versi dasarnya |
| **Arah** | sistem lama: **tidak pernah dihitung** pada jalur pengajuan addendum (NC-04, dua sebab). Kalkulator baru: baris selisih seperti besaran lain |
| **Besaran** | tidak dapat diramalkan dari sistem lama — tidak ada angka pembanding yang pernah ada |
| **Kelas cacat** | **NC-04** |

**Batas yang harus tertulis di laporan run:** untuk kelas ini, *"selisih yang tidak dapat
dijelaskan"* **tidak berarti kalkulator baru salah**, karena tidak ada angka lama untuk
menjelaskannya terhadap apa. Yang diperiksa hanya **kewajaran besarannya**, dan bila `DB-18`
dibenarkan, seluruh kelas ini **dikeluarkan**, bukan diperbaiki.

---

## PC-3 — `GRL-15`: tidak ada selisih, dan itu ramalannya

| | |
|---|---|
| **Populasi** | seluruh versi |
| **Arah** | **nol selisih** |
| **Alasan** | sistem lama selalu menghasilkan `ProRatePercent = 100` (NC-01+NC-03), dan sistem baru tidak memprorata sama sekali. Dua jalan berbeda yang bertemu di angka yang sama |
| **Kelas cacat** | — |

**Dan di sinilah ramalan ini dapat gagal, yang justru membuatnya berguna:** bila run paritas
menemukan versi yang selisihnya **bukan** nol pada besaran yang diproratakan, maka
`ProRatePercent ≠ 100` pernah terjadi — dan itu **mematahkan NC-03** dari arah yang sama sekali
berbeda dengan `UA-19`.

> Ini uji yang **memisahkan**: dua bacaan tentang masa lalu menghasilkan hasil run yang berbeda.
> Ia bukan pengamatan yang kebetulan konsisten.

---

## PC-4 — `GRL-17`: `NOMOR_URUT_VERSI` DIKECUALIKAN dari kriteria identik

**Ini bukan kelas cacat. Ini pengecualian yang harus tertulis, atau run paritas akan menandai
keberhasilan sebagai kegagalan.**

ADR-0043 menetapkan *"setiap angka identik"*. `NOMOR_URUT_VERSI` **tidak punya padanan di sistem
lama** — sistem lama tidak menyimpan nomor urut, ia hanya menyimpan **pengenal** `‹kontrak›/Rnn`.
Maka:

| Yang tunduk kriteria identik | Yang tidak |
|---|---|
| **pengenal `/Rnn`**, dilestarikan apa adanya (GRL-09, GRL-17) | **`NOMOR_URUT_VERSI`**, yang **diberikan**, bukan dipindahkan |

**Yang justru wajib diperiksa run paritas untuk kelas ini** adalah dua hal lain, dan keduanya dapat
gagal:

1. **Setiap versi warisan punya tepat satu `NOMOR_URUT_VERSI`, dan tidak ada yang kembar di dalam
   satu kontrak.** Kembar berarti kronologinya gagal memisahkan dua baris, dan kontraknya masuk
   pengecualian migrasi.
2. **Versi berlaku turunan (GRL-11) atas data hasil migrasi sama dengan baris yang hari ini dibaca
   hilir.** Bila berbeda, yang berbeda itu **diperiksa satu per satu**, bukan disamakan —
   perbedaannya justru isi `DB-12` dan `UA-17` (arti hukum versus arti operasional).

---

## Apa yang belum dapat ditulis, dan kenapa

**Populasi terdampak dalam angka** — belum ada satu pun. Tidak ada instans basis data yang
terjangkau dan izin kueri baca-saja belum turun; `UA-19`, `UA-20`, dan `UA-21` yang mengisinya.

Bentuknya ditulis sekarang justru karena itu: **bentuk harus mendahului angka**, atau uji terima
berubah jadi latihan menjelaskan.
