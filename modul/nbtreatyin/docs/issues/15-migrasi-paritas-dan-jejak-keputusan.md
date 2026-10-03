# 15: Migrasi, paritas, dan jejak keputusan — tidak ada butir terbuka yang ditutup diam-diam

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** 01 · 08 · 09 · 10
**Menutup:** AC 68 · 70 · 93 · 94 · 95 · 96 *(6 AC)* — US 43 · 47 · 48

## Hasil & nilai pengguna

Hari ini riwayat kontrak treaty ada di sistem lama, dan ⛔ sebagian nilainya **tidak dapat dibaca
tanpa membongkar dokumen teks**. ⚠️ Sebagian tanggalnya **ambigu secara mutlak**.

Sesudah tiket ini, data lama terbaca di sistem baru, ⭐ setiap penyimpangan dari perilaku lama
**tercatat beserta alasannya**, dan ⭐ setiap butir yang masih terbuka **tertulis di satu tempat** —
sehingga tidak ada yang terlupakan saat go-live.

## Area codebase

- Alat migrasi: pemindahan data lama
- Uji paritas: perbandingan perilaku lama dan baru
- Uji dokumen: kelengkapan penanda dan rujukan

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Bentuk data lama | dokumen teks hasil potret halaman kerja; ⛔ **bukan skema** |
| Tanggal ambigu | dua susunan tercampur, disimpan sebagai teks |
| Enam penyimpangan sadar | `spec.md` §10.1 |

## ADR terkait

- ⭐ **ADR-0009** — seluruh data dipindahkan; ⛔ tidak ada koeksistensi dua penulis
- **ADR-0003** — uang tidak `float`

## Acceptance criteria

- [ ] 🟡 **AC 68** — data lama **terbaca** di sistem baru; berkas lama dapat dibuka *(putaran 2: pemuat tiket 22 dibangun — dokumen lama ditulis lewat `SimpanHalaman` dan dibaca `BacaHalaman`, `TestPemuatLamaMenulisLewatAntarmukaSama` bertag `db`, **belum dijalankan**, K11)*
- [ ] ⛔ **AC 70** — setiap penyimpangan dari perilaku lama **tercatat beserta alasannya**
- [ ] ⛔ **AC 93** — setiap acceptance criterion merujuk **bab asalnya**
- [ ] ⛔ **AC 94** — setiap acceptance criterion membawa **penanda**
- [ ] ⛔ **AC 95** — ⛔ butir terbuka **tidak ditutup**; tidak ada yang dinyatakan selesai tanpa
      pemiliknya
- [x] **AC 96** — ⛔ nilai berupa **nama orang** tidak muncul di artefak mana pun

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⚠️ **15** | data lama bertanggal **ambigu mutlak** — satu nilai, dua arti | ⛔ **MENAHAN migrasi tanggal**, bukan tiketnya |
| ✅ ~~**P29**~~ | ~~contoh isi dokumen — dibutuhkan untuk membaca data lama~~ | ✅ **TERJAWAB 2026-09-22** — **tiga** contoh diterima |
| ⛔⛔ **galat uang tersimpan** | ekor `2,76 x 10^-7` pada `NetPremium`, sudah permanen di produksi | ⛔ **MENENTUKAN AMBANG PARITAS** — lihat catatan di bawah |

## Perintah verifikasi

1. Pindahkan satu berkas lama, buka di sistem baru — ⭐ seluruh medan yang dipakai **terisi**.
2. Sisir spec dan tiket — ⭐ **nol** acceptance criterion tanpa penanda, ⭐ **nol** tanpa rujukan bab.
3. Sisir seluruh artefak — ⭐ **nol** nama orang, **nol** nomor polis apa adanya.
4. Bandingkan hasil penggolongan dan potongan pajak lama lawan baru — ⭐ **sama**.

## Catatan

⚠️ **Paritas dapat diuji sekarang untuk yang sudah dibangun** — penggolongan, pajak, tanggal,
riwayat. ⛔ **Paritas rantai perhitungan uang menunggu tiket 13.**

---

## ⛔⛔ Ambang paritas uang — tidak dapat ditetapkan sebelum satu keputusan turun

`[terbuka]` `[terverifikasi]` 2026-09-22 — ⚠️ **dicatat di sini karena tiket inilah yang memiliki
uji paritas**; butir aslinya ada di **P29**, dan **tidak ditutup di sini**.

Data produksi memuat **galat angka uang yang sudah tersimpan permanen**:

```
premium angsuran  148157378.220000069   x 4  =  592629512.880000276
NetPremium        592629512.880000276        <- sama persis
nilai bulat 2 desimal                        =  592629512.88
selisih                                      =  2,76 x 10^-7
```

⛔⛔ **Akibatnya langsung bagi tiket ini:** sistem baru memakai aritmetika desimal yang benar
*(ADR-0003, nol `float`)*, sehingga ia **tidak akan menghasilkan ekor itu**. ⛔ **Uji paritas yang
menuntut kesamaan persis akan GAGAL — justru karena sistem barunya benar.**

| Bila keputusan P29 jatuh ke | Ambang paritas yang berlaku |
| --- | --- |
| **a** ikuti apa adanya | ⛔ paritas **persis**; galat lama sengaja ditiru |
| **b** bulatkan saat migrasi | ⚠️ paritas **tidak berlaku** untuk angka historis — angkanya memang berubah |
| **c** simpan apa adanya, bulatkan saat dihitung | ⭐ paritas dengan **toleransi**, dan besarnya toleransi itu yang harus ditetapkan |

⛔ **Ambangnya tidak ditetapkan di sini.** Ia mengikuti keputusan `[work owner]` dan `[Finance]`
pada **P29**. ⚠️ **Yang tidak boleh terjadi:** uji paritas ditulis dengan ambang yang dikarang,
lalu selisihnya diam-diam dianggap wajar.

> ⚠️ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat di atas semula berakhir dengan
> ⛔ *"… yang menunggu P18."* **P18 ditarik.** Tiket 13 kini menunggu **P30** dan **P8**.
