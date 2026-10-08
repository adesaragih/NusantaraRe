# 05: Persetujuan bersyarat dan dua penanda usul

**Status:** dibangun SEBAGIAN — subjectivity bertingkat tertahan OQ-KCP-01 *(RALAT 08-10-2026; status lama: `ready-for-agent`)*

> **RALAT 08-10-2026** — implementasi satu modul (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7). Kalimat lama tetap di bawah, dikutip di sini:
>
> - Persetujuan bersyarat di tangga >1 tingkat: isian tingkat 1 dibaca lagi di tingkat akhir (S16/S17/S21/S23/S24/S34) tetapi header kasus komite tidak punya kolom untuk menyimpannya — **OQ-KCP-01**; sementara ditolak validasi. Tangga satu tingkat dibangun utuh.


**Blocked by:** **04 (keputusan penyetuju tersimpan)**

## Hasil & nilai pengguna

Sebagai **penyetuju komite**, saya bisa menandai bahwa persetujuan saya **bersyarat**, lalu menulis
**isi syaratnya**. Saya juga bisa **mengusulkan penutupan** atau **pencadangan** klaim. Keempatnya
hanya muncul ketika memang relevan, jadi layar tidak penuh isian yang tidak berlaku.

*(User story 15, 16, 17 di spec)*

## Rantai gerbangnya

`[terverifikasi]` Keempat isian ini **bersyarat**, dan rantainya **bertingkat**:

| Isian | Muncul bila |
| --- | --- |
| penanda persetujuan bersyarat | penyetuju **menyetujui** **dan** jalur **penyesuaian** |
| isi syarat | penanda bersyarat **dicentang** |
| usul tutup klaim | jalur **penyesuaian** |
| usul cadangkan klaim | jalur **penyesuaian** |

⛔ **RALAT ke catatan lama:** keempatnya pernah disebut *"selalu dapat disunting"*. **Itu keliru** —
yang selalu dapat disunting hanya **dua** isian di tiket 04.

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Section/ShowTransfer.xml` | `[terverifikasi]` keempatnya bergerbang **famili sel**; kedua penanda usul berupa **kotak centang**, dapat disunting, **tidak wajib isi** |
| idem | `[terverifikasi]` **nol penulis di activity mana pun** — keduanya ditulis **langsung dari layar** |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 24** menulis penanda bersyarat dan isi syaratnya ke baris penyesuaian induk |

## ⭐ Kedua penanda usul punya DUA tempat di Pega

`[terverifikasi]` **KomitePostAdjustment langkah 11** — *"add propose close and propose reserve"*,
**tanpa gerbang**, jadi selalu berjalan — menyalin keduanya ke **kasus klaim induk dengan nama
berbeda**:

```
kasus komite . usul tutup klaim       ->  kasus klaim . penanda tutup berkas
kasus komite . usul cadangkan klaim   ->  kasus klaim . penanda klaim dicadangkan
```

⚠️ **RALAT ke catatan lama:** pernyataan *"tidak ada tulis-menulis lintas modul"* **hanya benar
untuk layar**. Lewat activity, komite **memang menulis ke kasus klaim induk**.

⛔ Apakah sisi klaim ikut disimpan adalah **urusan modul Claim Prop**, bukan tiket ini.

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Dua penanda usul **disimpan di tabel komite**, karena penyetuju yang mengisinya |
| **2026-09-19** | Keduanya disimpan di **header kasus komite** — *"komite menyimpan catatan usulnya sendiri, terpisah dari nilai akhir yang tercatat di kasus klaim"* |

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 24 · 55 · **74**

- [ ] Penanda bersyarat **hanya muncul** bila penyetuju menyetujui **dan** jalurnya penyesuaian.
- [ ] Isi syarat **hanya muncul** setelah penanda bersyarat dicentang.
- [ ] Kedua penanda usul **hanya muncul** pada jalur penyesuaian.
- [ ] Keempatnya **tidak wajib isi** — pengiriman tanpa mengisinya tetap diterima.
- [ ] Kedua penanda usul tersimpan di **header kasus komite**.
- [ ] Nilai yang sama juga **sampai ke kasus klaim induk**, sesuai perilaku Pega.
- [ ] Kedua penanda usul tersimpan **hanya** sebagai `'1'` atau `'0'`; kotak yang **tidak disentuh**
      tersimpan `'0'`, bukan kosong, dan nilai lain **ditolak**.

## Butir `[terbuka]` yang menyentuh tiket ini

- ✅ **TERJAWAB** `[keputusan work owner]` 2026-09-19 — **kedua penanda usul disimpan sebagai teks
  satu huruf: `'1'` bila diusulkan, `'0'` bila tidak. Kolomnya `CHAR(1)` dan wajib isi — kotak
  centang yang tidak pernah disentuh tersimpan `'0'`, bukan kosong.** ⚠️ Migrasi data lama: nilai
  yang di Pega bermakna benar/salah dikonversi otomatis menjadi `'1'`/`'0'`, dan baris lama yang
  nilainya kosong menjadi `'0'`.
  > ⛔ **Kalimat lamanya dikutip, tidak dihapus:** *"**Nilai apa yang tersimpan** untuk kedua
  > penanda usul — di Pega ia kotak centang, tetapi nilai tersimpannya **tidak terbaca dari
  > korpus**. ⛔ Jangan ditebak."*
  >
  > ⚠️ Butir ini **tidak pernah masuk register 13 butir** di `spec.md`, jadi menjawabnya
  > **tidak mengubah jumlah register**. ⛔ **Nol butir register ditutup.**

- **Tidak ada butir `[terbuka]` lain yang menyentuh tiket ini.**

## Seam & verifikasi

Memakai ulang seam Claim Prop.
