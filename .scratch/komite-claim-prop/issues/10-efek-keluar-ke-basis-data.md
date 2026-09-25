# 10: Efek keluar ke basis data — akseptasi, riwayat, log

**Status:** ready-for-agent

**Blocked by:** **08 (nomor akseptasi)**

## Hasil & nilai pengguna

Sebagai **bagian akseptasi**, baris akseptasi tercatat di daftar akseptasi klaim sehingga angkanya
bisa ditelusuri. Sebagai **auditor**, setiap keputusan komite tercatat di riwayat akseptasi, dan
penolakan klaim tercatat tersendiri.

*(User story 25, 30, 31 di spec)*

## Perilaku Pega yang ditiru

`[terverifikasi]` Empat dari delapan efek keluar bersifat penulisan basis data, di
`Komite Claim Prop/Activity/KomitePostAdjustment.xml`:

| Langkah | Rule yang dipanggil | Yang ditulis |
| --- | --- | --- |
| **17** | `SaveAcceptation_Act` | baris akseptasi ke daftar akseptasi klaim |
| **28** | `InsertJsonClaimTreaty_act` | data klaim bentuk JSON |
| **31** | `InsertLogServiceClaim` | baris log pemantauan |
| **33** | `InsertHistoryAkseptasiPega_Sql` | **6 kolom** riwayat akseptasi |

`[terverifikasi]` Penulisan penolakan klaim ada di `KomitePost_Close` **langkah 13** dan
`KomitePost_Reject` **langkah 13**, keduanya **di dalam langkah berulang**.

⚠️ `[terverifikasi]` Kolom **kotak kerja** pada riwayat akseptasi diisi **literal tetap**, selalu,
tanpa cabang. `[terbuka]` di modul lain: baris komite **ikut terhitung** oleh pembaca riwayat di
modul Fac In / Treaty In. ⛔ **Bukan urusan tiket ini, dan jangan mengusulkan membuang kolomnya.**

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | **Sembilan rule yang menyimpan sendiri tetap menyimpan sendiri.** *"Ditiru apa adanya."* |

> ⚠️ **Risiko diterima sadar:** bila langkah sesudahnya gagal, **kesembilan tabel tetap terisi**,
> dan **tidak ada rule pembatal di korpus** — tidak ada penghapusan, pembatalan, maupun penanda
> batal. **Data separuh jadi adalah perilaku yang disengaja ditiru, bukan cacat yang terlewat.**

## ⚠️ TITIK YANG SENGAJA DIUBAH — urutan terhadap penyimpanan

`[keputusan work owner]` 2026-09-18. **Ini perubahan sadar, bukan peniruan.**

`[terverifikasi]` Di Pega **kedelapan efek keluar berjalan SEBELUM penyimpanan**. Penyimpanan ada di
**langkah 41**, sedangkan pengiriman ke Kasir di **34** dan email di **35**. **Nol yang sesudah.**

⚠️ Akibatnya di Pega: bila penyimpanan gagal, **uang sudah dikirim, email sudah sampai, dokumen
sudah dibuat, dan tiga tabel log sudah terisi** — sementara kasusnya sendiri tidak tersimpan.

> `[keputusan work owner]` — **URUTAN INI TIDAK DITIRU.** Fakta di atas dicatat sebagai **fakta**,
> **tidak mengikat rancangan**. Urutannya **akan disesuaikan**.
>
> ⛔ **Jangan mengunci urutan Pega sebagai syarat.** ⛔ Urutan penggantinya **belum diputuskan**,
> dan **tidak ditetapkan di tiket ini**.

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 45 · 56 · 57 · 58 · 59 · 62 · 63 · 69 · 70 · 73

- [ ] Nomor akseptasi terbit → baris akseptasi, riwayat, dan log **tercatat**.
- [ ] Riwayat akseptasi terisi **6 kolom**; kolom ketujuh yang tidak pernah diisi jalur komite
      **tetap kosong**.
- [ ] Penolakan klaim tercatat **tersendiri**.
- [ ] Keempat penulisan **menyimpan sendiri** dan **tidak dibatalkan** bila langkah sesudahnya gagal.
- [ ] Kolom kotak kerja terisi **literal tetap**, tanpa cabang.

## Butir `[terbuka]` yang menyentuh tiket ini

- **Kolom akun operator** pada riwayat akseptasi **tidak pernah diisi jalur komite**. Adakah penulis
  lain di luar modul ini, atau kolom itu memang selalu kosong. ⛔ Tidak ditebak, **tidak diusulkan
  dibuang**.
- **Urutan terhadap penyimpanan** belum ditetapkan — lihat blok di atas.

## Seam & verifikasi

`[keputusan work owner]` 2026-09-18 — **efek keluar diuji dengan layanan sungguhan, bukan pengganti
tiruan**, di **lingkungan uji terpisah, bukan produksi**.
⚠️ **Pengujian urutan efek keluar MENUNGGU** urutan barunya ditetapkan.
