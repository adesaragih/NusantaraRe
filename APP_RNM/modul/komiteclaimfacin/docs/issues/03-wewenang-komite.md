# 03: Wewenang komite — ditegakkan di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** 00 · `claim-facin\issues\08` *(wewenang klaim)*
**Menutup:** AC 23 · 24 · 25 · 26 · 27 · 28 · 29 *(7 AC)* — US 19–21

## Hasil & nilai pengguna

Hari ini ⛔ `[terverifikasi]` Hanya **satu** pemeriksaan pemilik giliran ada di korpus, dan ia menempel pada **tombol Submit di layar**. ⚠️ Di modul saudaranya **tidak ada sama sekali**. ⛔ Siapa pun yang dapat memanggil lapisan layanan dapat menyimpan keputusan untuk jenjang mana pun.

Sesudah tiket ini, ⭐ **Hanya akun beku pada jenjang berjalan** yang dapat menyimpan keputusan, dan penolakannya terjadi **di lapisan layanan** — ⛔ bukan di layar.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pemeriksaan pemilik giliran | satu-satunya di korpus; menempel pada **tombol Submit** |
| ⚠️ Penukaran identitas | ⛔ **dua akun ditukar menjadi akun ketiga SEBELUM pemeriksaan** — ⭐ **tidak dibawa** |
| Modul saudara | ⛔ **nol pemeriksaan** — tidak di aktivitas, tidak di layar |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K10** — ⭐ **ADR-0014 BERDIRI** — penegakan pemilik giliran **tetap di lapisan layanan**
- **K3 · K12** — ⭐ Pembandingan memakai **identitas akun**, ⛔ bukan alias

## Yang harus diuji

- [ ] Hanya **akun beku pada jenjang berjalan** dapat menyimpan keputusan
- [ ] Percobaan oleh akun lain **ditolak dengan galat** — ⛔ bukan diabaikan diam-diam
- [ ] Percobaan yang ditolak **terekam** di jejak audit
- [ ] ⛔ Penukaran akun menjadi akun ketiga **tidak dibawa**
- [ ] ⭐ Layar boleh menyembunyikan tindakan tak berwenang — ⛔ **itu bukan penegakan**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | ⛔⛔ **Isi daftar jabatan dan susunan jenjang BELUM ADA** — korpus tidak memuat satu pun rule otorisasi; 13 medan privilese seluruhnya kosong | ⛔⛔ **MENAHAN PEMBANGUNAN tiket ini.** ⭐ Tiketnya tetap ditulis dan jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum work owner memberikan daftarnya |

## Seam & verifikasi

**Seam:** lapisan layanan komite — ⭐ **uji lewat sini, bukan lewat layar**.
2. Panggil sebagai **bukan pemegang giliran** ⇒ ⛔ **ditolak**.
3. Sembunyikan tombolnya di layar, lalu panggil layanan **langsung** ⇒ ⛔ **tetap ditolak**.
4. Periksa jejak audit ⇒ ⭐ percobaan yang ditolak **tercatat**.
5. Panggil sebagai akun yang dulu 'menyamar' ⇒ ⛔ **ditolak** — ⭐ penukaran tidak dibawa.
