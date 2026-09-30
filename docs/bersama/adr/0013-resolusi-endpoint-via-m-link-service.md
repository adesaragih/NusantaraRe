---
status: accepted
tanggal: 2026-09-14
sumber: penutupan OQ-047 + OQ-018 (korpus + data DBA), keputusan work owner
menggantikan: ADR-0004
---

# Alamat layanan keluar di-resolve runtime dari `M_LINK_SERVICE` — bukan env var

Resolusi URL keluar **harus** berupa **lookup runtime** ke tabel `M_LINK_SERVICE` lewat kunci
`(KATEGORI_1, KATEGORI_2)`, **setiap kali** layanan dipanggil.

**DILARANG** menanam alamat sebagai literal, sebagai konstanta, **maupun sebagai env var.**
Yang boleh menjadi konstanta di kode hanyalah **kunci kategori** — identifier endpoint, bukan
alamatnya.

Pemisahan dev–prod ditangani **isi tabel per-database**: database dev memuat baris dev, database
production memuat baris production. Kode tidak tahu bedanya, dan tidak perlu tahu.

> ⚠️ **ADR ini menggantikan ADR-0004**, yang memutuskan hal sebaliknya (alamat menjadi env var,
> lookup tidak direplikasi). Lihat §"Mengapa ADR-0004 dibatalkan".

## Kontrak yang ditiru `[terverifikasi]`

`Claim Life/Activity/GetLinkService.xml`
(`ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY`) melakukan
`Obj-Browse` atas tabel `M_LINK_SERVICE` dengan

```
.KATEGORI_1 = Param.Kategori_1  AND  .KATEGORI_2 = Param.Kategori_2
```

mengambil kolom `.URL`, lalu memanggil `Connect-REST`.

`[terverifikasi]` Kunci untuk Claim — Life terbaca langsung di
`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`:
`Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`.

`[terverifikasi]` Lima berkas `Claim Life` memakai jalur ini: `GetLinkService.xml`,
`InsertGoogleStorage_Act.xml`, `GetUrlGoogleStorage_Act.xml`, `DeleteGoogleStorage_Act.xml`,
`serviceInsertArasapasClaimLife_act.xml`.

`[data DBA]` Isi tabel **19 baris**. Endpoint Life:
`http://10.100.10.75:7315/Nusare-Integration-WS/resources1/restws/NusareClaim/insertClaimLife`.

## Mengapa ADR-0004 dibatalkan

ADR-0004 ditulis **sebelum** isi `M_LINK_SERVICE` diketahui (OQ-047 masih terbuka). Alasannya waktu
itu: "mekanisme lookup tidak perlu direplikasi; env var lebih sederhana." Dua hal membatalkannya:

1. **`[data DBA]` Tabelnya nyata dan dipakai bersama** — 19 baris melayani banyak modul, bukan hanya
   Claim Life. Memindahkan alamat Claim Life ke env var akan membuat **satu endpoint punya dua
   sumber kebenaran** selama modul lain masih membacanya dari tabel.
2. **`[keputusan work owner]`** Pemisahan dev–prod di organisasi ini memang dikerjakan lewat isi
   database, bukan lewat konfigurasi aplikasi. Env var akan memindahkan tanggung jawab itu ke tempat
   yang tidak dikelola tim yang sama.

## Considered Options

- **Runtime lookup `M_LINK_SERVICE` setiap panggilan** — dipilih
- Env var per endpoint (**ADR-0004**) — ditolak, alasan di atas
- Lookup sekali saat start lalu di-cache selamanya — ditolak: perubahan alamat di tabel tidak akan
  terlihat tanpa restart, sehingga menghidupkan kembali persoalan dua sumber kebenaran.
  *(Cache ber-TTL pendek boleh dipakai sebagai optimisasi, asal kebenarannya tetap tabel.)*

## Consequences

- **Ketergantungan Oracle bertambah satu**: memanggil layanan luar kini menuntut database hidup.
  Ini ketergantungan ketiga di luar data, setelah penomoran (**ADR-0006**) dan token penyimpanan.
- `internal/config` **tidak** memuat URL. Ia memuat koneksi database; alamat datang dari
  `internal/repository`.
- **Kunci kategori menjadi bagian kontrak** — mengubah `("Klaim", "insertClaimLife")` adalah
  perubahan yang memutus, setara mengubah nama endpoint.
- Efek keluar (**ADR-0008**) harus membedakan **kunci tidak ditemukan** dari **jaringan gagal**:
  yang pertama kesalahan konfigurasi data dan tidak layak diantre ulang berkali-kali; yang kedua
  layak.
- **ADR-0005** (flag lingkungan) tetap berlaku untuk *apakah* efek keluar berjalan; ADR ini hanya
  mengatur *ke mana* ia pergi.

## Bukti pendukung: tidak ada alamat ter-hardcode untuk ditinggalkan `[terverifikasi]`

Penutupan **OQ-018** memastikan tidak ada yang perlu "dibersihkan" lebih dulu — modul `Claim Life`
memang tidak pernah menanam alamat bisnis:

| URL yang ditemukan di modul | Jumlah | Sifat |
| --- | ---: | --- |
| `https://community.pega.com/help_v88/…` | 132+ | `pyHelpURI`, tautan bantuan |
| `https://view.officeapps.live.com/op/view.aspx?src=` | 3 | penampil dokumen Office |

Pembedaan lingkungan di Pega memakai `Claim Life/When/IsPEGAPROD.xml`
(`@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN`): `pzProductionLevel = "5"` — sebuah **flag**, bukan
hostname tertanam.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-035** | `serviceInsertArasapasClaimLife_act` satu salinan dipakai dua konteks — kepemilikan kunci kategorinya belum ditetapkan |
| **OQ-018** (modul lain) | URL literal di 15 modul lain belum disapu; ADR ini menetapkan arah, penerapannya menyusul per modul |
