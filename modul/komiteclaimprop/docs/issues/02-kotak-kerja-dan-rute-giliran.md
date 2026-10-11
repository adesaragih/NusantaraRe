# 02: Kotak kerja penyetuju dan rute giliran

**Status:** dibangun *(RALAT 08-10-2026; status lama: `ready-for-agent`)*

> **RALAT 08-10-2026** — implementasi satu modul (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7). Kalimat lama tetap di bawah, dikutip di sini:
>
> - Daftar kerja = KomiteRouter S6.1 (baris tangga PERTAMA ber-keputusan 0), saringan `w.LINI = 'PROP'` ketat + awalan `TKMT-`. AC 81-82 (wewenang simpan keputusan, galat terbaca) ditempel ke tiket **04**.


**Blocked by:** **01 (kasus komite lahir)**

## Hasil & nilai pengguna

Sebagai **penyetuju komite**, kasus yang menunggu keputusan saya muncul di **kotak kerja saya** —
dan **hanya saat giliran saya**. Saya tidak melihat pekerjaan orang lain, dan saya tidak diminta
menilai sesuatu yang belum dilihat penyetuju sebelum saya.

*(User story 1, 2, 3 di spec)*

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Flow/KomiteTreaty_Flow.xml` · `ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITETREATY_FLOW` | `[terverifikasi]` satu **assignment** bernama `KomiteRouter`, dirutekan **khusus** ke satu pengguna, berbentuk antrean per-pengguna |
| `Komite Claim Prop/Activity/KomiteRouter.xml` | `[terverifikasi]` 8 langkah, **6 di-remark**. Yang hidup **langkah 6** dan anaknya **6.1**: bergerbang *"keputusan penyetuju masih kosong"*, lalu menetapkan tujuan rute dari **akun operator** penyetuju yang sedang giliran. Pencacah dinaikkan **langkah 5** |

⭐ `[terverifikasi]` **Tangga persetujuan yang lama — antrean berjenjang `komitepnc1..4` — sudah
mati**; keenam langkah itulah yang di-remark. ⛔ **Jangan dipindahkan.**

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | **Hak akses per activity diabaikan** — daftar hak akses Pega menyebut kelas **tanpa menyebut hak**, jadi ia tidak menegakkan apa pun |

⚠️ Konsekuensinya: **siapa boleh melihat apa ditentukan aturan peran sistem baru**, bukan disalin
dari Pega. Tiket ini hanya menegakkan **giliran**, bukan wewenang.

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 3 · 4 · 10 · 11 · 14 · 67

- [ ] Kotak kerja seorang penyetuju **hanya** memuat kasus yang sedang **gilirannya**.
- [ ] Penyetuju yang belum tiba gilirannya **tidak melihat** kasus itu.
- [ ] Penyetuju yang sudah memutuskan **tidak melihatnya lagi** di kotak kerja.
- [ ] Tujuan rute diambil dari **akun operator** penyetuju, bukan jabatan dan bukan nama antrean.

## Butir `[terbuka]` yang menyentuh tiket ini

- **Pewarisan kelas kerja** tidak ada di ekspor Pega — menyentuh bagaimana properti kasus komite
  diwarisi. Tidak menghambat tiket ini.

## Seam & verifikasi

Memakai ulang seam Claim Prop.
