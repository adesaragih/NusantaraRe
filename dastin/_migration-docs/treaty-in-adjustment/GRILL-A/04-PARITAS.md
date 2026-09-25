> Modul  : Treaty In Adjustment · Ronde A · 2026-09-24
> Peran  : interogator
> Masukan: GRL-01 s.d. GRL-04 · `METODE` §6.5
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 04 · SKENARIO PARITAS

## Kenapa fase ini hampir kosong, dan itu bukan pengisi

Paritas menguji **pelestarian**: perilaku baru dicocokkan dengan perilaku lama atas data lama.
Keempat putusan ronde A berlabel **PERUBAHAN** atau **BARU** — dan menurut `METODE` §6.5, keduanya
justru **tidak dapat** diuji terhadap data lama: perubahan karena data lama memperlihatkan perilaku
yang dibuang, dan yang baru karena tidak punya pembanding.

Maka ronde ini melahirkan **dua** skenario paritas, bukan satu daftar panjang — dan keduanya bukan
paritas perilaku melainkan **paritas data**.

## P-A1 · Nilai lama hasil migrasi harus cocok dengan potret `OLDDATA`

| | |
|---|---|
| **Sifat** | paritas data, bukan paritas perilaku |
| **Dasar** | GRL-02; `PENGETAHUAN.md` §2.1, §10.1 |
| **Bentuk** | untuk setiap addendum warisan, nilai versi pendahulu hasil migrasi dibandingkan dengan sub-pohon `OLDDATA` di dalam `JSONDATA` addendum itu |
| **Yang diharapkan** | cocok |
| **Bila tidak cocok** | **alat deteksi mutu migrasi**, bukan kegagalan uji — ambangnya diputuskan di cabang I |
| **Kenapa ia mungkin** | karena sistem lama membekukan potret; inilah satu-satunya hal dari sistem lama yang dapat dipakai memeriksa hasil migrasi |

## P-A2 · Selisih tersimpan versus selisih terhitung

| | |
|---|---|
| **Sifat** | paritas data |
| **Dasar** | GRL-02, GRL-03; ADR-0048 (*"ketidakcocokan … ia alat deteksi"*) |
| **Bentuk** | selisih yang dibekukan pada persetujuan dibandingkan dengan selisih yang dihitung ulang dari kedua versi |
| **Yang diharapkan** | **berbeda** pada tiga keadaan yang sudah diketahui: baris disisipkan/dihapus di tengah daftar, mata uang berbeda antar versi, dan angka ringkasan bagian |
| **Kenapa perbedaannya bukan kegagalan** | yang berbeda adalah **angka yang dulu salah**; ketiga keadaan itu diukur `UA-5` dan `UA-6` |

## Yang TIDAK dapat diuji paritas, dan disebut supaya tidak dikira terlewat

| Putusan | Kenapa |
|---|---|
| GRL-01 — versi menggantikan pendahulunya | **kemampuan baru**: langkahnya ada di kode tetapi mati (`METODE` §3.8). Tidak ada perilaku lama untuk dibandingkan; rencana penyalaannya dititipkan ke cabang D |
| GRL-03 butir 2 — selisih di tabel tersendiri | pembandingnya ada tetapi cacat (sub-pohon `ValueDifference`), jadi ia bahan **uji migrasi**, bukan paritas |
| GRL-04 — penyesuaian bukan entitas | tidak menghasilkan perilaku yang dapat diamati |
