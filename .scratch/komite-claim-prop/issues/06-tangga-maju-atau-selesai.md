# 06: Tangga maju ke penyetuju berikutnya, atau selesai

**Status:** ready-for-agent

**Blocked by:** **04 (keputusan penyetuju tersimpan)**

## Hasil & nilai pengguna

Sebagai **pengelola proses**, kasus **otomatis berpindah** ke penyetuju berikutnya setelah satu
penyetuju memutuskan, dan **selesai** ketika penyetuju terakhir menyetujui. Tidak ada langkah manual
dan tidak ada tahap menggantung.

*(User story 19, 20 di spec)*

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Flow/KomiteTreaty_Flow.xml` | `[terverifikasi]` **satu jalur balik**: dari gerbang keputusan **kembali ke assignment yang sama**, bersyarat *"masih ada penyetuju berikutnya"*. Bila tidak, kasus **selesai** dengan status akhir *selesai-tuntas* |
| idem | `[terverifikasi]` **empat bentuk saja**: mulai → assignment → gerbang keputusan → selesai. Aksi pengguna pada layar komitelah yang memindahkan kasus dari assignment ke gerbang |
| `Komite Claim Prop/Activity/KomiteRouter.xml` | `[terverifikasi]` **langkah 5** menaikkan pencacah penyetuju |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 40** juga menaikkan pencacah, tetapi **gerbangnya dimatikan** |

⭐ `[terverifikasi]` **Tangga berputar DI DALAM satu tahap.** Setiap penyetuju adalah **kunjungan
baru ke tahap yang sama**, bukan tahap baru. Yang berubah tiap putaran hanya **kepada siapa kasus
dirutekan** dan **nilai pencacah**.

⚠️ `[terverifikasi]` Pembacaan riwayat pada langkah 8 dan 9 terjadi **jauh sebelum** penambahan
pencacah di langkah 40 — jadi tidak ada pembacaan indeks di luar daftar.

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Pencacah jumlah penyetuju **disimpan** dan **tidak dihitung ulang** tiap dibaca |

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 5 · 6 · 7 · 12 · 15 · 17

- [ ] Sesudah satu penyetuju memutuskan, kasus **berpindah** ke penyetuju berikutnya.
- [ ] Penyetuju terakhir yang **menyetujui** **menyelesaikan** kasus.
- [ ] Kasus yang selesai **tidak muncul lagi** di kotak kerja siapa pun.
- [ ] Pencacah penyetuju naik **tepat satu** per putaran; test yang menemukannya melompat **gagal**.
- [ ] Jumlah putaran **tidak melebihi** jumlah penyetuju yang dibutuhkan.

## Butir `[terbuka]` yang menyentuh tiket ini

- **Arti dua jenis perulangan Pega** dan **halaman apa yang diulang** belum terbaca.
- **Berapa kali perulangan berputar** belum terbaca — menyentuh berapa kali langkah di dalamnya
  berjalan.
- **Urutan pemeriksaan bila dua keluarga gerbang sama-sama terisi** tidak terbaca dari struktur.

⚠️ Ketiganya **tidak menghambat** perilaku tangga di atas, yang terbaca dari berkas alur, bukan dari
perulangan langkah.

## Seam & verifikasi

Memakai ulang seam Claim Prop.
