# 07: Penolakan menghentikan tangga dan menolak sisa penyetuju

**Status:** dibangun *(RALAT 08-10-2026; status lama: `ready-for-agent`)*

> **RALAT 08-10-2026** — implementasi satu modul (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7). Kalimat lama tetap di bawah, dikutip di sini:
>
> - S25 (`KomiteCount := KomiteLoop`, `AcceptanceStatus := 2`, `AktifButton := 0`), S26.1 sisa tangga menunggu -> 2, S27 Notes; `FlagOnGoingCommitte` dibuang (keputusan 27).


**Blocked by:** **06 (tangga maju atau selesai)**

## Hasil & nilai pengguna

Sebagai **pengelola proses**, **satu penolakan menghentikan rangkaian** — penyetuju berikutnya tidak
diminta menilai sesuatu yang sudah ditolak. Sisa penyetuju **ditandai otomatis**, jadi daftarnya
tidak tertinggal setengah terisi.

*(User story 21, 22 di spec)*

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 12** — bila keputusan **tolak**, melompat ke tanda `EXT`, yaitu **langkah 25** |
| idem | `[terverifikasi]` **langkah 25** memaksa pencacah ke nilai akhir sehingga tangga **berhenti** |
| idem | `[terverifikasi]` **langkah 26** menyisir daftar penyetuju; anaknya **26.1** menandai sisa penyetuju **tertolak otomatis**. Langkah 26 adalah **langkah berulang** |
| idem | `[terverifikasi]` **langkah 23 · 25** menulis status penolakan ke baris penyesuaian induk |

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | **Lima perintah lompat menggantung dibuang** — kelimanya menunjuk tanda yang tidak ada, dan kelimanya bergerbang mati. ⛔ Tidak ditebak ke mana seharusnya melompat, tidak ditandai cacat, **dibuang** |

⚠️ **Lompatan yang HIDUP tetap dipindahkan.** Yang dibuang hanya yang **mati**. Lompatan ke tanda
`EXT` pada langkah 12 **hidup dan tandanya ada** — ia bagian dari tiket ini.

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 5 · 13 · 62 · 68

- [ ] Satu keputusan **tolak** menghentikan rangkaian; penyetuju berikutnya **tidak menerima** kasus.
- [ ] Seluruh sisa penyetuju **ditandai tertolak**; tidak ada baris penyetuju yang tertinggal
      berkeputusan kosong.
- [ ] Status penolakan sampai ke **baris penyesuaian induk**.
- [ ] Kasus **selesai** sesudah penolakan — tidak menggantung.
- [ ] ⛔ Tidak ada perintah lompat menggantung yang ikut dipindahkan. Test yang menemukan lompatan
      ke tanda yang tidak ada **gagal**.

## Butir `[terbuka]` yang menyentuh tiket ini

- **Berapa kali perulangan berputar** — menyentuh berapa kali penandaan sisa penyetuju dikerjakan.

## Seam & verifikasi

Memakai ulang seam Claim Prop.
