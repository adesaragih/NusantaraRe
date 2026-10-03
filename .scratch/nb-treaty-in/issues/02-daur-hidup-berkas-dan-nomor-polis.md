# 02: Daur hidup berkas realisasi dan nomor polis — satu nomor, sekali, tanpa bentrok

**Status:** ready-for-agent
**Blocked by:** 01
**Menutup:** AC 31 · 59 · 73 · 74 *(4 AC)* — US 1 · 4 · 5

## Hasil & nilai pengguna

Hari ini admin treaty membuka penawaran yang sudah disetujui dan melengkapinya menjadi realisasi.
`[terverifikasi]` Nomor polis dibentuk dari sebuah deret, dan ⚠️ salah satu bagiannya diambil dari
slot parameter generik yang **tidak terlihat diisi** di aktivitas mana pun.

Sesudah tiket ini, sebuah realisasi treaty **lahir dari penawaran yang disetujui**, mendapat
**satu nomor polis yang tidak bentrok**, dan ⭐ pengguna **diperingatkan** ketika ia hendak membuat
penawaran yang sudah pernah ada.

## Area codebase

- Lapisan service: daur hidup berkas realisasi
- Lapisan service: pembentukan nomor polis
- Lapisan handler: peringatan penawaran ganda

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Titik masuk alur | `Flow\InputRealizationTreatyIn.xml` — `pyStartActivity` = `Start1` |
| Pembentukan nomor polis | `RDBList\GenerateNoPolicy.xml` |
| Peringatan penawaran ganda | `Activity\CheckDuplicateOffer.xml` langkah 6 dan 8, lewat `RDBList\GetCountClaim.xml` |
| Rantai pemanggilnya | `TreatyRealizationCheckXOLList` → `SetTreatyIn_Act` → `CheckDuplicateOffer` |

## ADR terkait

- **ADR-0006** — penomoran lewat deret basis data

## Acceptance criteria

- [ ] **AC 73** — nomor polis memuat awalan tetap, penanda treaty, bulan-tahun, dan nomor urut
      berdigit tetap
- [ ] **AC 74** — nomor dibentuk **sekali** per berkas; tidak berubah pada penyimpanan berikutnya
- [ ] **AC 31** — dua berkas **tidak pernah** bernomor sama
- [ ] **AC 59** — peringatan muncul ketika jumlah berkas klaim terhubung **lebih dari nol**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **P7** | asal salah satu bagian nomor polis **tidak terlihat diisi** | ⚠️ tidak menahan — polanya diketahui |

## Perintah verifikasi

1. Realisasikan dua penawaran berbarengan — ⭐ nomor polisnya **berbeda**.
2. Simpan ulang salah satunya — ⭐ nomornya **tidak berubah**.
3. Buat penawaran yang sudah pernah ada — ⭐ **peringatan muncul**.
