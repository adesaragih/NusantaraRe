# Urutan tiket — Komite Claim Prop

⛔ **Ini bukan tiket.** Berkas ini hanya menjelaskan **tiket mana harus selesai sebelum tiket mana**.
Isinya tidak menambah keputusan apa pun.

Sumber: `spec.md` dan `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`.

---

## Empat belas tiket

| # | Judul | Blocked by |
| --- | --- | --- |
| **00** | Skema penyimpanan komite — tiga tabel + dua kolom usul — **PREFACTOR** | skema klaim Claim Prop |
| **01** | Terima penyerahan dari Claim Prop — kasus komite lahir | 00 · kontrak muatan Claim Prop |
| **02** | Kotak kerja penyetuju dan rute giliran | 01 |
| **03** | Layar komite — 93 kolom, tiga wajah | 01 |
| **04** | Keputusan penyetuju tersimpan | 02 · 03 |
| **05** | Persetujuan bersyarat dan dua penanda usul | 04 |
| **06** | Tangga maju ke penyetuju berikutnya, atau selesai | 04 |
| **07** | Penolakan menghentikan tangga dan menolak sisa penyetuju | 06 |
| **08** | Nomor akseptasi — terbit sekali, di tingkat akhir | 06 |
| **09** | Angka uang — total penyesuaian per mata uang | 03 |
| **10** | Efek keluar ke basis data — akseptasi, riwayat, log | 08 |
| **11** | Efek keluar — dokumen PDF akseptasi | 08 |
| **12** | Efek keluar — Kasir, arasapas, dan email | 08 · 10 |
| **13** | Migrasi data lama | 00 |

---

## Urutan pengerjaan

```
                       00  PREFACTOR
                        |
              +---------+---------+
              |                   |
             01                  13  migrasi
              |
        +-----+-----+
        |           |
       02          03
        |           |
        +-----+-----+-----------+
              |                 |
             04                09  angka uang
              |
        +-----+-----+
        |           |
       05          06
                    |
              +-----+-----+
              |           |
             07          08
                          |
                    +-----+-----+
                    |     |     |
                   10    11    12
                          (12 juga menunggu 10)
```

**Frontier awal:** hanya **00**. Sesudah 00 selesai, **01** dan **13** terbuka bersamaan.

**Jalur terpanjang:** `00 → 01 → 02/03 → 04 → 06 → 08 → 10 → 12` — delapan tingkat.

**Yang bisa dikerjakan paralel:**

- **02** dan **03** sesudah 01
- **05**, **06**, dan **09** sesudah prasyaratnya masing-masing
- **07** dan **08** sesudah 06
- **10** dan **11** sesudah 08
- **13** berjalan sendiri sepanjang waktu, hanya menunggu 00

---

## ⚠️ Tiket yang menunggu modul Claim Prop

`[keputusan work owner]` 2026-09-18 — **seam memakai ulang milik Claim Prop, tidak menambah seam
baru.** Akibatnya menyentuh **seluruh empat belas tiket**: tidak satu pun bisa **diverifikasi
ujung-ke-ujung** sampai seam Claim Prop berdiri.

Di luar itu, **tiga tiket punya ketergantungan isi** pada modul Claim Prop:

| Tiket | Menunggu apa dari Claim Prop |
| --- | --- |
| **00** | **tabel klaim** — baris komite tinggal di tabel lintas-lini yang sama |
| **01** | **kontrak muatan penyerahan** — bentuk data yang diserahkan ke komite. **Bukan** penyelesaian modul Claim Prop; muatannya dapat dipalsukan di seam supaya kedua modul jalan paralel |
| **09** | **rumus angka klaim** — empat dari lima penghitung berkelas kasus klaim, jadi **pemiliknya Claim Prop**. Modul ini hanya menjumlahkan penyesuaian per mata uang |

⚠️ Tiket **12** menyentuh Claim Prop lewat syarat nomor akseptasi, tetapi **rule pemeriksanya sudah
ada di modul ini** — yang di luar hanyalah SQL-nya, dan itu **tidak ditelusuri**
`[keputusan work owner]`.

---

## Pengujian — berlaku untuk seluruh tiket

`[keputusan work owner]` 2026-09-18:

| Hal | Ketetapan |
| --- | --- |
| **Seam** | **memakai ulang seam Claim Prop**, tidak menambah |
| **Efek keluar** | diuji dengan **layanan sungguhan**, bukan pengganti tiruan, di **lingkungan uji terpisah** |
| **Urutan efek keluar** | pengujiannya **MENUNGGU** urutan barunya ditetapkan |

⚠️ `[terbuka]` **Baris mana di tabel alamat layanan yang menunjuk lingkungan uji** belum ditetapkan.
Menyentuh tiket **11** dan **12**.

---

## ⚠️ Dua titik yang SENGAJA DIUBAH — muncul di tiket mana saja

| Titik | Muncul di tiket |
| --- | --- |
| **Urutan efek keluar terhadap penyimpanan** | **10** · **11** · **12** |
| **Ketelitian angka** | **00** · **03** *(tampilan)* · **09** *(hitungan)* |

⛔ Keduanya **bukan peniruan**. Jangan diperlakukan sebagai cacat pemindahan bila hasilnya berbeda
dari Pega.
