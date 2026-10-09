# Pertanyaan terbuka — penyebut `CacahTerukur` sesudah korpus KEDUA ikut dimuat

**Untuk pemilik proses. Diajukan 6 Oktober 2026.**
**Nol pekerjaan terblokir oleh ini. Yang rusak adalah ALAT UKURNYA, bukan datanya.**

---

## Pertanyaannya

> **Sesudah `M_TREATY_IN_EDM` ikut didaratkan ke tabel yang sama, `CacahTerukur`
> harus menyebut korpus yang mana?**

| Jawab | Akibatnya |
| --- | --- |
| **"`M_TREATY_IN` saja"** (keadaan hari ini) | `pemuat -cocokkan` menandai ⚠️ pada **setiap** baris selamanya, sebab cacah nyata selalu lebih besar. Tandanya menjadi derau, dan derau yang selalu menyala sama saja dengan tidak ada penjaga. |
| **"jumlah keduanya"** | satu angka, tetapi ia tidak lagi dapat menjawab *"korpus mana yang kurang?"* bila selisih muncul. |
| **"dua angka"** | `Pendaratan` memperoleh ruas kedua (`CacahTerukurEDM`), dan `cetakCacah` mengadu `nyata` lawan `jumlah keduanya` sambil tetap dapat menyebut sisi mana yang meleset. ⚠️ Menuntut perubahan struktur yang juga menyentuh salinan peta di modul Adjustment. |

---

## Keadaan terukur, supaya jawabannya tidak perlu diukur ulang

Sapuan 6 Oktober 2026 atas **seluruh** kedua korpus, tanpa `ROWNUM`:
**1.855** dokumen `M_TREATY_IN` + **280** dokumen `M_TREATY_IN_EDM` (kedua sisi
`New`/`OLDDATA` ikut), **2.135 dokumen terurai, nol gagal urai**.

| | baris |
| --- | ---: |
| dari `M_TREATY_IN` | **115.563** |
| dari `M_TREATY_IN_EDM` (dua sisi) | **60.407** |
| **jumlah yang akan berdiri di tabel** | **175.970** |
| yang `CacahTerukur` hari ini sebut | **115.563** |

⭐ **Yang MENETAPKAN arti `CacahTerukur` hari ini:** sapuan yang sama mengukur ulang
ke-26 entri yang angkanya sudah terisi sebelum ronde ini, dan **ke-26-nya cocok persis**
dengan kolom `M_TREATY_IN` (702 · 2.298 · 4.210 · 11.475 · 15.746 · …). Jadi penyebutnya
bukan soal selera — ia sudah tertetapkan oleh saudara-saudaranya, dan ketiga angka yang
ronde ini isi mengikutinya.

---

## Mengapa ini bukan keputusan yang boleh diambil sendiri

Pilihan **"dua angka"** mengubah `repository.Pendaratan`, dan peta itu **disalin** di
modul Adjustment (modul dilarang saling mengimpor). `uji/lintasmodul/peta_pendaratan_test.go`
membuktikan salinannya tidak basi, jadi perubahan strukturnya menjalar ke dua modul
sekaligus.

Pilihan **"jumlah keduanya"** membuang kemampuan menunjuk sisi yang meleset — dan itu
kemampuan yang baru saja terbukti berharga: ronde ini menemukan `CacahLarik` **buta**
terhadap sembilan tabel migrasi 438/439, dan yang memperlihatkannya justru pembanding
per-korpus.

---

## Yang TIDAK dilakukan sambil menunggu

- ⛔ Struktur `Pendaratan` **tidak diubah**.
- ⛔ Ketiga angka baru **tidak dijumlahkan** diam-diam; angka EDM-nya ditulis sebagai
  komentar di `pendaratan_peta.go`, di sebelah angkanya, supaya jawaban mana pun dapat
  dipasang tanpa mengukur lagi.
- ⚠️ `pemuat -cocokkan` **akan** menandai ⚠️ pada setiap baris mulai hari muatan EDM
  masuk. Itu dinyatakan di kepala peta, bukan dibiarkan mengejutkan pembaca berikutnya.
