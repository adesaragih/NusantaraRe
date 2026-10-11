# 04: Keputusan penyetuju tersimpan

**Status:** dibangun *(RALAT 08-10-2026; status lama: `ready-for-agent`)*

> **RALAT 08-10-2026** — implementasi satu modul (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7). Kalimat lama tetap di bawah, dikutip di sini:
>
> - AC 81-82 ditempel di sini: simpan keputusan hanya oleh pemilik `KOMITE_OPERATORID` tingkat berjalan, ditegakkan di layanan (403 terbaca); uji `TestBukanPemilikDitolak`.


**Blocked by:** **02 (kotak kerja dan rute giliran)** · **03 (layar komite)**

## Hasil & nilai pengguna

Sebagai **penyetuju komite**, saya menyatakan **setuju atau tolak** dan menulis **catatan**, lalu
menekan kirim. Keputusan saya tersimpan pada baris saya, lengkap dengan **tanggal dan identitas
saya** — tanpa saya perlu mengisinya.

*(User story 13, 14, 18 di spec)*

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Section/ShowTransfer.xml` | `[terverifikasi]` **dua isian selalu dapat disunting dan wajib isi**: keputusan setuju/tolak, dan catatan. Keduanya **tanpa gerbang** |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 6** menulis keputusan, catatan, dan tanggal ke baris penyesuaian induk · **langkah 27** menulis catatan · tanggal diisi **waktu saat itu** |
| idem | `[terverifikasi]` **langkah 4** membuka kasus klaim induk **dengan penguncian**, dan kuncinya **dilepas saat penyimpanan** — bukan dipegang selama komite menimbang |

⭐ `[terverifikasi]` **Penguncian hanya menyelimuti proses simpan**, beberapa detik. Penyetuju bisa
membuka layar berjam-jam tanpa mengunci apa pun.

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | **Satu kolom tanggal persetujuan saja** — properti tanggal kedua di Pega **sengaja tidak dijadikan kolom**; jangan ditambahkan kembali karena "ada di korpus" |
| 2026-09-18 | Pada jalur **tolak** dan **tutup**, penyetuju memang hanya mengisi **dua** hal ini. **Ditiru apa adanya** |

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 23 · 25 · 29 · 53 · 54

- [ ] Keputusan dan catatan tersimpan pada **baris penyetuju yang bersangkutan**, bukan di header.
- [ ] Tanggal dan identitas penyetuju **terisi otomatis**; penyetuju tidak bisa mengubahnya.
- [ ] Keduanya **wajib isi** — pengiriman tanpa salah satunya **ditolak**.
- [ ] Pada jalur **tolak** dan **tutup**, hanya kedua isian ini yang tersedia.
- [ ] Penguncian kasus klaim induk **hanya berlangsung selama proses simpan**, dan **dilepas**
      sesudahnya. Test yang menemukan kunci tertahan sesudah simpan **gagal**.
- [ ] Hanya **satu** kolom tanggal persetujuan yang ada. Test yang menemukan kolom tanggal kedua
      **gagal**.

## Butir `[terbuka]` yang menyentuh tiket ini

Tidak ada.

## Seam & verifikasi

Memakai ulang seam Claim Prop.
