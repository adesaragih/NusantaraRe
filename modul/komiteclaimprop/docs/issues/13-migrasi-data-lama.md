# 13: Migrasi data lama

**Status:** alat sensus uji-kering dibangun; `-jalankan` ditolak (OQ tangga lama) *(RALAT 08-10-2026; status lama: `ready-for-agent`)*

> **RALAT 08-10-2026** — implementasi satu modul (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7). Kalimat lama tetap di bawah, dikutip di sini:
>
> - Sensus DEV 08-10-2026: 2.451 kasus CLMP lama, 419 ditunda pemuat Claim Prop (baris berlaku STS_REJECT 1): 2 punya work object KomiteTreaty ber-BLOB di DATAPEGA, 316 hanya riwayat akseptasi (tanpa tingkat / jabatan / operator), 101 tanpa jejak komite. Tangga tidak dikarang; AC 83 (seluruh data dipindahkan) ditempel di sini dan menunggu sumber tangga (OQ pemilik ekspor Pega).


**Blocked by:** **00 (skema — PREFACTOR)**

## Hasil & nilai pengguna

Sebagai **pengguna sistem baru**, kasus komite lama beserta daftar penyetuju dan keputusannya
**terbaca di sistem baru**, sehingga riwayat tidak hilang saat pindah.

## Yang dipindahkan

| Dari Pega | Ke sistem baru |
| --- | --- |
| kasus komite | baris komite di tabel lintas-lini + header kasus komite |
| daftar penyetuju | satu baris per penyetuju, beserta keputusan, catatan, dan tanggalnya |
| pencacah tangga | kedua pencacah apa adanya |
| penunjuk ke baris penyesuaian induk | penunjuk baris penyesuaian di header kasus |
| **dua penanda usul** | ⭐ **benar/salah → `'1'`/`'0'` · kosong → `'0'`** — lihat aturan konversi di bawah |

### ⭐ Aturan konversi dua penanda usul

`[keputusan work owner]` 2026-09-19 — kolom tujuannya **`CHAR(1)`, wajib isi, nilai sah hanya
`'1'` dan `'0'`** *(tiket 00)*. Karena itu:

| Nilai lama di Pega | Menjadi |
| --- | --- |
| bermakna **benar** *(dicentang)* | **`'1'`** |
| bermakna **salah** *(tidak dicentang)* | **`'0'`** |
| **kosong / tidak ada nilainya** | **`'0'`** |

⛔ **Konversinya otomatis dan tanpa pengecualian** — migrasi **tidak boleh** meninggalkan kolom
usul kosong, dan **tidak boleh** menolak baris hanya karena nilainya kosong.

⚠️ **Ini satu-satunya tempat migrasi memperbaiki nilai.** Ia **tidak** melonggarkan keputusan
2026-09-19 tentang penunjuk posisional di bawah — penunjuk tetap **dipindahkan apa adanya**.

## ⚠️ Kunci lama bersifat POSISIONAL

`[terverifikasi]` Di Pega, kasus komite menunjuk baris penyesuaian induknya lewat **nomor urut
baris**, bukan pengenal. Kuncinya dua bagian: kunci kasus klaim induk, ditambah **posisi** di dalam
daftar penyesuaian klaim itu.

⚠️ **Artinya: bila sebuah baris penyesuaian pernah disisipkan atau dihapus di klaim induk sesudah
kasus komitenya dibuat, kasus komite itu menunjuk baris yang KELIRU — dan tidak ada yang
menyadarinya.** `[terverifikasi]` **Tidak ada penanganan untuk hal itu di Pega** — tidak ada
pemeriksaan, tidak ada pesan.

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| **2026-09-19** | **Penunjuk posisional DIPINDAHKAN APA ADANYA, termasuk yang sudah salah alamat di Pega. Tidak diperiksa dulu, tidak diperbaiki.** Konsisten dengan preseden *"tiru apa adanya"* |

> ⚠️ **Risikonya, ditulis tegas:** **penunjuk yang salah di Pega akan TETAP SALAH sesudah migrasi.**
> Sistem baru **mewarisi kekeliruan itu apa adanya, tanpa penanda.** ⛔ Migrasi **tidak boleh**
> memperbaikinya diam-diam, dan **tidak boleh** menolak baris yang mencurigakan.

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 61 · 64 · 65 · 66 · 71 · **75**

- [ ] Kasus komite lama terbaca lengkap dengan daftar penyetuju dan keputusannya.
- [ ] Penunjuk posisional diterjemahkan menjadi penunjuk baris penyesuaian **apa adanya** — termasuk
      yang mengarah ke baris yang berbeda dari yang semestinya.
- [ ] ⛔ Migrasi **tidak menolak** dan **tidak memperbaiki** baris yang penunjuknya mencurigakan.
      Test yang menemukan baris ditolak atau diperbaiki **gagal**.
- [ ] Kedua pencacah tangga dipindahkan apa adanya.
- [ ] Kedua penanda usul terisi `'1'` atau `'0'` sesudah migrasi — **nol baris berkolom usul
      kosong**. Test yang menemukan kolom usul kosong sesudah migrasi **gagal**.
- [ ] Nilai uang lama yang **melewati batas digit** menyebabkan **kegagalan yang terlihat**, bukan
      pembulatan diam-diam.

## Butir `[terbuka]` yang menyentuh tiket ini

- ⭐ **DIPERSEMPIT 2026-09-19 — kini hanya untuk KASUS LAMA HASIL MIGRASI.** Kasus komite lama yang
  penunjuk baris penyesuaiannya **sudah salah alamat atau menunjuk baris yang tidak ada**: apa yang
  dilakukan aplikasi saat kasus itu dibuka — **menolak terang-terangan**, atau **diam seperti
  Pega**. Terdaftar sebagai butir **14** di register `spec.md`.
  > ⛔ **Kalimat lamanya dikutip, tidak dihapus:** *"**Apa yang terjadi bila baris penyesuaian
  > induk tidak ketemu** *(saat berjalan, bukan saat migrasi)* — tidak ada penanganannya di
  > korpus."*
  >
  > **Kenapa menyempit:** `[keputusan work owner]` 2026-09-19 membekukan baris penyesuaian sejak ia
  > diserahkan ke komite *(bab 11)*, sehingga **kasus baru tidak bisa lagi kehilangan baris
  > induknya**. ⚠️ **Kunci beku itu tidak berlaku surut**, jadi pertanyaannya **tidak hilang** —
  > ia hanya tinggal berlaku untuk **data lama hasil migrasi**. ⛔ **Butir ini TIDAK ditutup.**

## Seam & verifikasi

Migrasi diverifikasi terhadap **skema uji**, bukan produksi.
