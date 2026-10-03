# 01: Terima penyerahan dari Claim Prop — kasus komite lahir

**Status:** ready-for-agent

**Blocked by:** **00 (skema — PREFACTOR)** · **kontrak muatan penyerahan Claim Prop** — bukan
penyelesaian modul Claim Prop.

## Hasil & nilai pengguna

Sebagai **petugas klaim**, ketika saya menyerahkan satu baris penyesuaian ke komite, **kasus komite
terbentuk** lengkap dengan daftar penyetuju, pencacah tangga, dan penunjuk ke baris penyesuaian
saya. Saya tahu penyerahan saya diterima.

*(User story 5, 19 di spec)*

Ini **tracer bullet** pertama yang terlihat: menembus penerimaan penyerahan → penyimpanan →
pembacaan kembali.

## ⚠️ Penyerahan membekukan baris penyesuaian induknya — **PERILAKU BARU**

`[keputusan work owner]` 2026-09-19 — **saat kasus komite lahir, baris penyesuaian induknya
ditandai BEKU di sisi Claim Prop:** tidak dapat diubah dan tidak dapat dihapus, **selamanya**.
**Bekunya tidak mencair meskipun komite menolak** — perbaikan dilakukan dengan **baris penyesuaian
baru**, bukan dengan menyunting baris lama.

⛔ **PENEGAKANNYA MILIK CLAIM PROP, BUKAN TIKET INI.** Daftar penyesuaian adalah milik modul Claim
Prop, jadi penolakan sunting dan penolakan hapus dibangun di sana — pada tiket **08** *(baris
adjustment)* dan tiket **11** *(penyerahan ke Komite)* modul itu. **Tiket ini hanya pemicunya:**
ia mencatat penyerahan, dan pencatatan itulah yang menjadikan baris induknya beku.

⚠️ **Ini perilaku BARU, bukan paritas.** `[terverifikasi]` Di seluruh ekspor modul ini: **nol
pemeriksaan, nol pesan kesalahan, nol penanganan** untuk baris penyesuaian yang berubah atau hilang
sesudah diserahkan. **Aturan ini dibangun, bukan dimigrasikan.**

**Lingkup:** Claim Prop dan Claim Non Prop; `[keputusan work owner]` **Claim Non Prop menyusul —
fokus sekarang Claim Prop saja. Claim — Life tidak termasuk.**

⛔ **Tidak berlaku surut** — kasus komite lama hasil migrasi tidak terlindungi; lihat butir **14**
register `[terbuka]` dan tiket **13**.

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` | `[terverifikasi]` **langkah 13** menetapkan jenis pengajuan dan penunjuk posisi baris · **langkah 14** menyalin isi baris penyesuaian · **langkah 15** menyalin catatan komite · **langkah 22.1** membangun daftar penyetuju, tiap baris ber-keputusan awal `0` |
| idem | `[terverifikasi]` pencacah **jumlah penyetuju** diisi dari **jumlah baris daftar penyetuju**, sekali saat kasus dibuat; pada jalur tutup/tolak yang memakai satu penyetuju tetap nilainya `1` |

⭐ `[terverifikasi]` **Daftar penyetuju ditetapkan saat kasus dibuat dan tidak berubah di tengah
jalan** — nol penulis di modul komite ini.

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Pencacah jumlah penyetuju **disimpan**, diisi ulang di titik yang sama dengan Pega — **bukan** dihitung ulang tiap dibaca |
| 2026-09-18 | **Aturan awalan**: data berawalan titik dan halaman kasus = **komite**; berawalan halaman induk = **klaim**, **dibaca saja** |

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 1 · 8 · 9 · 16 · **76** · **77** · **78**

- [ ] Penyerahan satu baris penyesuaian **melahirkan** baris komite di tabel lintas-lini **dan**
      header kasus komite ber-kunci sama.
- [ ] Penyerahan **membuat satu baris penyetuju per anggota daftar**, masing-masing ber-keputusan
      awal kosong dan ber-urutan sesuai jenjangnya.
- [ ] Pencacah **jumlah penyetuju** terisi dari jumlah baris daftar; pencacah **penyetuju keberapa**
      dimulai di penyetuju pertama.
- [ ] Penunjuk ke baris penyesuaian induk terisi, dan penunjuk balik dari sisi klaim terisi
      **dalam satu transaksi yang sama**. Test yang menemukan salah satunya kosong **gagal**.
- [ ] Muatan penyerahan tanpa daftar penyetuju, atau dengan jumlah penyetuju kurang dari satu,
      **ditolak** — kasus komite **tidak dibuat**.
- [ ] ⚠️ Sesudah kasus komite lahir, baris penyesuaian induknya **tidak dapat disunting** dan
      **tidak dapat dihapus** — termasuk lewat kaskade hapus klaim. Ditegakkan **di sisi Claim
      Prop**; tiket ini menguji bahwa **penandanya benar-benar terpasang** saat penyerahan tercatat.
- [ ] ⚠️ Bekunya **tidak mencair** sesudah komite menolak; test yang **berhasil menyunting** baris
      sesudah penolakan **gagal**.

## Butir `[terbuka]` yang menyentuh tiket ini

- **Berapa kali perulangan berputar** di Pega belum terbaca — menyentuh berapa kali nilai
  penyesuaian ditulis saat kasus dibuat.

## Seam & verifikasi

Memakai ulang seam Claim Prop. Muatan penyerahan **dipalsukan di seam** sampai sisi Claim Prop
nyata, supaya kedua modul bisa dikerjakan paralel.
