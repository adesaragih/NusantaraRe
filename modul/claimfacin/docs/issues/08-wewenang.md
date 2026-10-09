# 08: Wewenang — ditegakkan di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi) · 05 (penyesuaian)
**Menutup:** AC 72 · 73 · 74 *(3 AC)* — US 19 · 20 · 21

> ⛔ **RALAT 10-10-2026.** Status lamanya dikutip utuh, tidak dihapus: *"**Status:** ready-for-agent"* → status itu
> **bertentangan** dengan tiket ini sendiri dan `urutan-tiket.md`, yang menyatakan tiket ini **MENAHAN PEMBANGUNAN**
> sampai butir 6 dijawab — jadi ia sebenarnya **tertahan**. ⭐ Penahannya **kini terjawab**: ADR-0030
> (`docs/bersama/adr/0030-aturan-peran-ditetapkan-sekali-lintas-modul.md`) + keanggotaan workbasket, dibangun
> 10-10-2026 di lapisan layanan (`backend/services/layanan.go` `Pemegang`, `services/aksi.go`):
>
> - **Input Register** dan **Input Estimasi** = pembuat kasus (`CREATE_OP`).
> - **Choose Surveyor** = anggota workbasket **`ReasKlaimTeknik`** — `ReasPNCTeknik` yang ditulis XML tidak ada di DEV
>   (`PINDAI.md` §2).
> - Server menolak tulisan ke **kasus tertutup** (`ErrKasusTertutup`) dan ke **adjustment yang sedang di komite**
>   (`ErrAdjustmentDiKomite`); bukan pemegang assignment = **403** (`ErrBukanPemegang`). Aksi yang tidak tampil / aktif
>   di tata layar kasus ditolak (`ErrAksiTertutup`) — layar tidak menjadi penjaga satu-satunya.
>
> ⚠️ *"13 medan privilese"* di tiket ini berbeda dari spec Bab 15 (**15** nama medan); yang mengikat spec.

## Hasil & nilai pengguna

Hari ini ⛔ `[terverifikasi]` Korpus **tidak memuat satu pun rule otorisasi**, dan 13 medan privilese seluruhnya kosong. Siapa pun yang dapat memanggil lapisan layanan dapat mengerjakan apa pun.

Sesudah tiket ini, Wewenang **ditegakkan di lapisan layanan** — ⛔ bukan di layar. Percobaan yang tidak berwenang **ditolak dan terekam**, bukan diabaikan diam-diam.

## Area codebase

- Lapisan layanan: penegakan wewenang
- Jejak audit atas percobaan yang ditolak

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Wewenang | ⛔ **nihil di korpus** — 13 medan privilese kosong; ADR-0014 mencatatnya ABSENT |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0014** — wewenang ditegakkan di lapisan layanan
- **ADR-0002** — peran ditegakkan di lapisan layanan, bukan visibilitas layar
- **ADR-0012** — wewenang eksplisit

## Acceptance criteria

- [ ] **AC 72 · 73 · 74** — wewenang
- [ ] ⭐ Penolakan terjadi **di lapisan layanan**
- [ ] ⭐ Percobaan yang ditolak **terekam**

> ⛔ **RALAT 10-10-2026 — `[penyimpangan sadar]` DIBUANG: tidak dibangun.** Butir lamanya dikutip utuh, tidak dihapus:
> *"⭐ Percobaan yang ditolak **terekam**"* → janji ini **di luar XML** (spec AC 72: sistem lama nol gerbang wewenang,
> jadi tidak ada perekaman penolakan untuk ditiru). Yang dibangun: penolakan dijawab **403 / 409** oleh handler
> (`backend/handlers/rute.go`), **tanpa baris jejak** di tabel mana pun. Perintah verifikasi 2 tidak berlaku.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | ⛔⛔ **Isi daftar jabatan dan susunan jenjang BELUM ADA** — korpus tidak memuatnya | ⛔⛔ **MENAHAN PEMBANGUNAN tiket ini.** ⭐ Tiketnya tetap ditulis, jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum daftar jabatan diberikan work owner |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"⛔⛔ **MENAHAN PEMBANGUNAN tiket ini.**"* →
> **tidak menahan lagi** (butir 6 = register **spec**): wewenang memakai ADR-0030 + keanggotaan workbasket, bukan daftar
> jabatan — lihat RALAT di kepala tiket. Susunan jenjang **komite** menurut jabatan (roster `EMAILKOMITE` FACIN) beralih
> ke workbasket pola Komite Claim Prop di **tahap 2** (`komiteclaimfacin`, OQ-CFI-04).

## Perintah verifikasi

1. Panggil lapisan layanan sebagai pengguna **tanpa wewenang** — ⛔ **ditolak**.
2. Periksa jejak audit — ⭐ percobaan yang ditolak **tercatat**.
3. ⛔ Sembunyikan tombolnya di layar lalu panggil layanan langsung — ⛔ **tetap ditolak**.
