# 08: Wewenang — ditegakkan di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi) · 05 (penyesuaian)
**Menutup:** AC 72 · 73 · 74 *(3 AC)* — US 19 · 20 · 21

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

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | ⛔⛔ **Isi daftar jabatan dan susunan jenjang BELUM ADA** — korpus tidak memuatnya | ⛔⛔ **MENAHAN PEMBANGUNAN tiket ini.** ⭐ Tiketnya tetap ditulis, jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum daftar jabatan diberikan work owner |

## Perintah verifikasi

1. Panggil lapisan layanan sebagai pengguna **tanpa wewenang** — ⛔ **ditolak**.
2. Periksa jejak audit — ⭐ percobaan yang ditolak **tercatat**.
3. ⛔ Sembunyikan tombolnya di layar lalu panggil layanan langsung — ⛔ **tetap ditolak**.
