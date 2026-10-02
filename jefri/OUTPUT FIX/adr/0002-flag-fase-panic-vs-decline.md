---
status: accepted
---

# Satu flag fase: `panic` saat paralel run, `decline`+log saat produksi

Baris keputusan `DecisionTable/IsUWAccepted` tidak ikut terekspor, sehingga pemetaan **status →
konektor flow** tidak terbaca meskipun arti tiap nilai status sudah terverifikasi. Untuk jalur yang
belum terverifikasi, sistem berperilaku berbeda menurut fase: **`panic` selama paralel run** (agar
setiap jalur yang belum terbukti muncul ke permukaan — itu memang tujuan paralel run), dan
**`decline` + catatan log saat produksi** (karena `decline` adalah perilaku terekam sistem lama,
dan menjatuhkan case produksi tanpa alasan bisnis lebih berbahaya). Keduanya dikendalikan **satu
flag fase eksplisit**, bukan dua basis kode.

## Considered Options

- **`panic` di kedua fase** — ditolak untuk jalur tipe (a) di bawah. Paralel run akan berhenti pada
  kasus yang di sistem lama berjalan mulus sebagai `decline`, dan di produksi ia menjatuhkan case
  yang seharusnya lanjut.
- **`decline` di kedua fase** — ditolak. Menyembunyikan celah pengetahuan kita justru pada fase yang
  seharusnya menemukannya.

## Consequences

Dua jenis celah dibedakan tegas, dan **hanya yang pertama** mengikuti flag fase:

| | Situasi | Paralel run | Produksi |
| --- | --- | --- | --- |
| **(a)** | Status **dikenali** artinya, baris keputusannya belum terverifikasi | `panic` | `decline` + log |
| **(b)** | Kondisi **tidak dikenali sama sekali** | `panic` | **`panic`** |

Transisi sebuah jalur (a) dari `panic` menjadi `decline`+log **hanya** boleh dilakukan setelah jalur
itu diverifikasi terhadap ekspor produksi tunggal — bukan karena ia sering muncul dan mengganggu.

`[terverifikasi]` `decline` adalah `<pyDefaultResult>` rule `IsUWAccepted`, jadi memilihnya di
produksi adalah reproduksi, bukan tebakan.

Mesin keadaannya sendiri dibangun **dari sisi penulis status**, yang seluruhnya terverifikasi berikut
efek sampingnya (mis. rule penulis `2` mematikan `ConfirmBinding` dan `ReceivedRiSlip`; penulis `4`
menyalakan `IsBanding` dan menyetel tujuan banding). Tabel keputusan hanya menentukan pemetaan
status → konektor, dan itu diverifikasi belakangan tanpa menahan pekerjaan.
