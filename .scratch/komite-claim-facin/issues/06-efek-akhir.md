# 06: Efek akhir — terbit sekali, oleh jenjang terakhir

**Status:** ready-for-agent
**Blocked by:** 05
**Menutup:** AC 43 · 44 · 45 · 46 · 47 · 48 · 49 · 50 · 51 *(9 AC)* — US 28 · 29 · 32 · 33

## Hasil & nilai pengguna

Hari ini Ringkasan akseptasi, nomor akseptasi, dokumen, dan surel **belum terbit sama sekali** — dan ⚠️ di sistem lama, langkah yang menyusunnya terpisah dari langkah yang menyimpannya, sehingga pernah dikira **penerimaan tidak tersimpan**.

Sesudah tiket ini, ⭐ **Efek akhir terbit HANYA ketika jenjang terakhir menyetujui**, dan **tepat satu kali**: ringkasan akseptasi disusun lalu **disimpan**, nomor terbit, dokumen tercetak, surel terkirim.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penyusun ringkasan | 51 langkah; ⛔ langkah penyimpan **di dalamnya** ber-remark |
| ⭐ Penyimpan sesungguhnya | ⭐ langkah penyimpan **diangkat ke rule pemanggil**, bergerbang *keputusan terima* **dan** *jenjang terakhir* |
| ⚠️ Jalur kegagalan | ⛔ **NOL jalur kegagalan pada 25 langkah penyimpanan** di sistem lama |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K1** — ⭐ Efek akhir menempel pada **jenjang terakhir menyetujui**

## Yang harus diuji

- [ ] ⭐ Efek akhir terbit **hanya** ketika **jenjang terakhir menyetujui**
- [ ] ⭐ Ringkasan akseptasi **disusun lalu disimpan** — penyimpanan terjadi **sekali**
- [ ] ⛔ **RALAT yang wajib dibaca:** vonis lama *“penolakan tersimpan, penerimaan tidak”* **BATAL** — ⭐ keduanya tersimpan; yang berbeda hanya **di rule mana** langkah penyimpannya duduk
- [ ] Urutannya: ringkasan → nomor → dokumen → surel → kasir
- [ ] Surel terkirim **hanya di lingkungan produksi**
- [ ] ⚠️ Kegagalan salah satu efek **tidak boleh** membatalkan efek yang sudah terbit tanpa jejak — ⭐ sistem baru **wajib** punya jalur kegagalan

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

## Seam & verifikasi

**Seam:** lapisan layanan komite — penyelesaian kasus.
1. Jalankan tangga tiga jenjang sampai selesai ⇒ ⭐ efek akhir terbit **tepat satu kali**.
2. Periksa jenjang pertama dan kedua ⇒ ⛔ **nol efek akhir** di keduanya.
3. Paksa gagal di tengah rangkaian efek ⇒ ⭐ kegagalannya **tercatat**, tidak hilang diam-diam.
4. Jalankan di lingkungan bukan produksi ⇒ ⛔ surel **tidak terkirim**.
