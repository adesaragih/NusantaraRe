# 12: Efek keluar — Kasir, arasapas, dan email

**Status:** dibangun SEBAGIAN — efek diantre outbox, pelaksana berhenti terang *(RALAT 08-10-2026; status lama: `ready-for-agent`)*

> **RALAT 08-10-2026** — implementasi satu modul (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7). Kalimat lama tetap di bawah, dikutip di sini:
>
> - Konversi Arasapas (S29), Kasir (S34), dan email (S35) diantre `inti/backend/outbox` hanya di produksi (`IsPEGAPROD`); pelaksana `PelaksanaKomiteClaimProp` me-resolve M_LINK_SERVICE lalu berhenti `…BelumDisetujui` sampai panggilan nyata disetujui. Pemeriksaan S3 `getStatusKonversi_Act` dipindah ke pelaksana (konversi kini asinkron). `DIRECTTOKASIR_LOG` / StatusKasir = hasil pelaksana (OQ-CP-03).


**Blocked by:** **08 (nomor akseptasi)** · **10 (efek keluar ke basis data)**

## Hasil & nilai pengguna

Sebagai **bagian keuangan**, data pembayaran **terkirim ke Kasir** setelah akseptasi, sehingga
pembayaran bisa diproses — dan hanya bila **nomor akseptasi sudah tercatat**. Sebagai **pihak
terkait**, saya menerima **email pemberitahuan** hasil akseptasi tanpa membuka sistem.

*(User story 27, 28, 29 di spec)*

## Perilaku Pega yang ditiru

| Langkah | Rule | Yang keluar |
| --- | --- | --- |
| **29** | `KonversiKlaim_Act` | panggilan layanan **arasapas** |
| **34** | `HitServiceToKasirKMT_Act` | data pembayaran ke **Kasir** |
| **35** | `SendEmailKlaim_KMT` | **email** ke penerima akseptasi |

`[terverifikasi]` Alamat tujuan diambil dari **tabel alamat layanan**, disaring **dua kunci**.
Pasangan kuncinya berbeda per tujuan: satu untuk Kasir, satu untuk klaim, dua untuk penyimpanan
berkas.

## Syarat sebelum mengirim ke Kasir

`[data work owner]` 2026-09-18 — pengiriman ke Kasir **hanya boleh jalan bila nomor akseptasi sudah
tercatat** di daftar akseptasi klaim. **Itu syarat resmi yang dipertahankan.**

`[terverifikasi]` Rule pemeriksanya mengambil nomor akseptasi, membuang titiknya, mencari di tabel
akseptasi, dan menyalakan penanda hanya bila ketemu. ⛔ SQL-nya di luar modul ini — **tidak
ditelusuri**.

## ⚠️ Sikap terhadap kegagalan — **melompat dan melanjutkan**

`[keputusan work owner]` 2026-09-18 — **kiriman ke Kasir gagal → email TETAP dikirim, kasus tetap
jalan, nol percobaan ulang. DITIRU APA ADANYA.**

`[terverifikasi]` Mekanismenya: bila pengiriman gagal, alur **melompat ke langkah pengiriman email**
lalu selesai. Pola serupa pada arasapas, yang melompat ke langkah pencatatan log. Di **dalam** rule
pengirim Kasir ada jalur galat tersendiri yang **mengirim email khusus kegagalan**.

> ⚠️ **Risiko diterima sadar:** kegagalan **tidak terlihat** sampai ada orang yang memeriksa
> belakangan. Tidak ada pemberitahuan bahwa transfer gagal; yang sampai ke penerima justru **email
> keberhasilan akseptasi**. **Selisih antara "email terkirim" dan "uang terkirim" hanya ketahuan
> dari pemeriksaan manual.**

⚠️ `[terverifikasi]` **Email sendiri tidak punya penanganan gagal sama sekali.**

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

**Diverifikasi oleh:** spec.md AC 45 · 46 · 47 · 49 · 50 · 52 · 69 · 70 · 73

- [ ] Nomor akseptasi tercatat → data pembayaran **terkirim ke Kasir**.
- [ ] Nomor akseptasi **belum** tercatat → pengiriman ke Kasir **tidak berjalan**.
- [ ] Kegagalan pengiriman ke Kasir **tidak menghentikan kasus**, dan **email tetap terkirim**.
- [ ] Tidak ada percobaan ulang otomatis.
- [ ] Kegagalan pengiriman **tercatat** di jalur galatnya sendiri.
- [ ] Alamat tujuan diambil dengan **pasangan kunci yang benar** untuk masing-masing tujuan.

## Butir `[terbuka]` yang menyentuh tiket ini

- **Unggah berkas dan email tanpa penanganan gagal** — belum dibawa ke work owner.
- **Urutan terhadap penyimpanan** belum ditetapkan.
- **Baris mana di tabel alamat layanan yang menunjuk lingkungan uji** belum ditetapkan.
- **Seberapa besar selisih** akibat pembulatan di batas penyimpanan — angka yang dikirim ke Kasir
  berasal dari hitungan di tiket 09.

## Seam & verifikasi

Efek keluar diuji dengan **layanan sungguhan di lingkungan uji terpisah**
`[keputusan work owner]` 2026-09-18. ⚠️ **Pengujian urutan MENUNGGU** urutan barunya ditetapkan.
