---
status: accepted
---

# Mata uang yang tidak diketahui adalah keadaan eksplisit, bukan default

`Money.Currency` mengizinkan keadaan **`Unknown` yang eksplisit**. Nilai ber-mata-uang-tidak-diketahui
boleh dibaca dan ditampilkan, tetapi **`panic`** ketika dipakai dalam aritmetika lintas mata uang
atau ditulis ke Oracle. Tidak ada default diam-diam.

## Mengapa

`[terverifikasi]` Korpus memang kehilangan informasi ini: **112 Section menampilkan nilai uang tanpa
field mata uang mana pun**, dan hanya ada 93 pengikatan properti mata uang di seluruh NB. Sekaligus
`[terverifikasi]` sistem lama **tidak pernah** menetapkan mata uang default di kode — ia selalu
datang dari data (`CurrencyList` / `CurrencyMaster`), dan kurs selalu dari hasil query.

## Considered Options

- **Default ke IDR** — ditolak. Menebak (`CLAUDE.md` §3 butir 4), dan tebakan itu akan salah persis
  pada kasus yang paling mahal: penempatan valas.
- **`panic` tanpa kecuali bila mata uang tidak ada** — ditolak. Akan menghentikan aplikasi pada 112
  layar yang di sistem lama berjalan mulus — kesalahan yang sama bentuknya dengan mengganti
  `decline` terekam menjadi `panic` (lihat ADR-0002).

## Consequences

Prinsipnya sama dengan ADR-0002: celah pengetahuan dibuat **eksplisit dan terlihat**, lalu gagal
keras **hanya** di titik yang benar-benar memerlukan jawabannya.

| Operasi | Mata uang `Unknown` |
| --- | --- |
| Baca, tampilkan, simpan di memori | diizinkan |
| Aritmetika dalam satu mata uang yang sama-sama `Unknown` | diizinkan |
| Aritmetika lintas mata uang | **`panic`** |
| Tulis ke Oracle | **`panic`** |

Efek samping yang disengaja: implementasi menghasilkan **hitungan berapa banyak nilai produksi yang
tiba tanpa mata uang**. Angka itu masuk kuesioner — ia mengubah pertanyaan dari dugaan menjadi
ukuran.
