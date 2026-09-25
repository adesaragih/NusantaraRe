# `F-13` — tiga *materialized view* induk menegakkan tanpa pemantau kebasian

**Tanggal:** 24 September 2026
**Ditemukan saat:** adjudikasi kemampuan Adjustment terhadap `P-01`…`P-59` (to-ticket)
**Sifat:** **temuan berpemilik.** Bukan tiket, bukan keputusan. **Tidak ditambal dari modul
Adjustment** — ia lubang induk, dan menambalnya dari sana berarti memutuskan sesuatu untuk daftar
induk di luar lingkup yang diperintahkan.

> ## TEMUAN
>
> **`INV-47`, `INV-50`, dan `INV-51` ditegakkan lewat *materialized view* ber-`REFRESH ON COMMIT`,
> dan tidak satu pun `P-nn` di `DAFTAR-PEKERJAAN.md` menyebut pemantau kebasiannya.**

## Kenapa ia bukan kerapian

Aturan modul ini sudah menyatakannya sendiri, dan menyatakannya sebagai bahaya:

> ***MV* yang gagal me-refresh berhenti menegakkan tanpa satu galat pun.**

Itu instans `§2.0` — *struktur yang terlihat bukan struktur yang berlaku* — yang **lahir dari
rancangan kita sendiri**, bukan dari sistem lama. Sebuah constraint yang berdiri di atas MV yang
basi **terlihat terpasang di setiap pemeriksaan yang pernah dijalankan orang**, dan tidak menolak
apa pun.

Dan ketiganya bukan invarian pinggiran:

| Invarian | Yang dijaganya |
|---|---|
| `INV-47` | jumlah nilai penyebaran per rincian sama dengan besaran induknya |
| `INV-50` | persentase penyebaran berjumlah seratus |
| `INV-51` | nilai penyebaran per mata uang menutup terhadap induknya |

**Ketiganya menjaga uang**, dan ketiganya menjaga bagian model yang paling sering disesuaikan.

## Apa yang sudah ada, dan apa yang tidak

`UJI-NEGATIF-INVARIAN.md` sudah menuntut **dua** hal untuk penegakan yang dapat berhenti diam-diam:

| Tuntutan | Keadaan |
|---|---|
| **uji negatif** yang dijalankan, bukan diargumentasikan | tertulis — §2 |
| **pemantau kebasian** — `REFRESH_MODE`, `STALENESS`, terjadwal | tertulis — §3 |

**Tuntutan kedua tidak pernah menjadi kemampuan.** Ia ada sebagai kalimat di berkas invarian, dan
tidak ada satu baris pun di `DAFTAR-PEKERJAAN.md` yang seorang pelaku dapat dinyatakan telah
menyelesaikannya.

> Sebuah kewajiban yang tidak punya `P-nn` **tidak akan pernah masuk papan**, tidak akan pernah
> ditaksir, dan tidak akan pernah dinyatakan selesai. Ia akan dibaca, disetujui, dan dilewati.

## Bagaimana ia ketahuan dari modul Adjustment

`INV-69` dan `INV-70` — materialitas — juga ditegakkan lewat MV. Saat memeriksa apakah "pemantau
kebasian" sudah punya kemampuan sendiri yang dapat dirujuk, jawabannya **tidak** — dan pemeriksaan
yang sama menunjukkan ketiadaan itu **sudah ada sejak induk**, bukan lahir dari Adjustment.

Irisan yang menegakkan `INV-69`/`INV-70` karenanya membawa pemantau kebasiannya **di dalam kriteria
selesainya sendiri**, sebagai jalan keluar sementara. Itu menutup dua invarian dan **membiarkan
tiga**.

## Dua bentuk yang mungkin, dan pilihannya bukan milik saya

| Bentuk | Akibatnya |
|---|---|
| **satu `P-nn` tersendiri** — *"**PJ** dapat melihat bahwa setiap penegakan ber-MV masih segar, dan diberi tahu bila tidak"* | satu tiket, satu pemilik, berlaku untuk seluruh MV sekarang dan nanti |
| **sifat** (`U-6`) — dibawa setiap tiket yang memasang MV | tidak ada tiket tunggal yang dapat dilewati, tetapi **dikerjakan berulang kali dengan tafsir berbeda** — persis nasib yang `U-6` peringatkan untuk sifat yang tersebar |

Saya condong ke **yang pertama**, dengan alasan yang dapat diperiksa: pemantau kebasian punya
**pelaku nyata** (`PJ`, pemeriksa jejak) dan **titik selesai yang jelas** (ia berbunyi kepada
seseorang) — dua syarat yang `U-6` pakai untuk memisahkan kemampuan dari sifat. Tetapi
**menambahkannya sebagai `P-67` adalah keputusan atas daftar induk**, dan itu bukan milik sesi
Adjustment.

| | |
|---|---|
| **Siapa menutup** | **pemilik proses** — satu kalimat memilih bentuknya |
| **Yang menagih** | berkas ini; dan kriteria selesai irisan `INV-69` yang menyebut pemantau kebasian **tanpa punya nomor kemampuan untuk dirujuk** |
| **Bila dibiarkan** | tiga invarian uang berdiri di atas mekanisme yang **dapat berhenti tanpa memberi tahu siapa pun**, dan tidak ada tiket yang pernah menyalakan pemberitahuannya |
