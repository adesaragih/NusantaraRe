# 08: Penggolongan lini usaha di sisi komite

**Status:** ready-for-agent
**Blocked by:** 00 · `claim-facin\issues\02` *(klasifikasi lini)*
**Menutup:** AC 56 · 57 · 58 · 59 · 60 *(5 AC)* — US 34–36

## Hasil & nilai pengguna

Hari ini ⛔ `[terverifikasi]` Penggolongan lini usaha di sisi komite **membaca medan yang tidak pernah disalin** ke objek kerja komite. ⚠️ Akibatnya cabang **MBU** dan **Travel** **tidak pernah terbit** — diam-diam, tanpa galat.

Sesudah tiket ini, Setiap lini usaha **tergolong benar di sisi komite**, termasuk **MBU dan Travel** yang dulu tak pernah terbit.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| 13 penggolong hidup | ⭐ dari 49 rule penggolong, **13 dipakai hidup**; ⛔ **36 tidak dipakai sama sekali** |
| ⚠️ Cacat lama | ⛔ sembilan penggolong menguji medan yang **tidak disalin**; **dua dipakai hidup lima kali** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K2** — ⭐ **Data kutipan disalin UTUH** — ⚠️ `[penyimpangan sadar]`, **cacat yang diperbaiki**. ⛔ Bukan daftar medan bernama — **daftar bernama itulah yang melahirkan cacat ini**

## Yang harus diuji

- [ ] ⭐ **13 penggolong** dialihkan; ⛔ **36 tidak dibangun**
- [ ] ⭐ Penggolongan membaca **data kutipan lengkap**
- [ ] ⭐ Cabang **MBU** dan **Travel** **terbit** bila datanya memenuhi — ⚠️ inilah cacat lama yang ditutup
- [ ] ⛔ Penggolong berhubungan **ATAU** tetap berperilaku **ATAU**, ⚠️ bukan **DAN**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **29** | Cacah baris lama MBU/Travel terdampak — `[data DBA]` | tidak menahan |

## Seam & verifikasi

**Seam:** lapisan layanan komite — penyusunan ringkasan akseptasi.
1. Jalankan kasus komite lini **MBU** ⇒ ⭐ cabangnya **terbit**.
2. Jalankan lini **Travel** ⇒ ⭐ cabangnya **terbit**.
⚠️ Keduanya **tidak terbit di sistem lama** — ⭐ uji ini membuktikan cacatnya tertutup.
3. Cari pemanggilan salah satu dari 36 penggolong yang tidak dialihkan ⇒ ⛔ **nihil**.
