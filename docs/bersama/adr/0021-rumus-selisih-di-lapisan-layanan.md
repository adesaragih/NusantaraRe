---
status: accepted
tanggal: 2026-09-23
sumber: spec penyimpanan EDM Treaty In - `.scratch/edm-treaty-in/`, keputusan work owner
---

# Rumus selisih hidup di lapisan layanan, dan nilai lama tidak pernah dihitung ulang

Selisih endorsemen dihitung di lapisan **layanan**, bukan di basis data. Lapisan `repository`
hanya mengambil **dua baris**: baris generasi ini dan baris yang ditunjuknya.

## Satu rumus untuk semua

```
selisih.X = baris_ini.X - baris_yang_ditunjuk.X
```

Tanpa memandang endorsemen pertama atau berlapis, proporsional atau non-proporsional.

`[terverifikasi]` Terbukti dari dokumen produksi pada enam medan sekaligus - nilai baru nol, nilai
lama positif, selisih negatif sebesar nilai lama.

Empat aturan turunannya:

| Golongan | Perlakuan |
| --- | --- |
| Uang | **dikurangi** |
| Persen | **disalin, tidak dikurangi** |
| Kunci | disalin |
| Penanda arah | diturunkan dari **tanda jumlah**, bukan dari satu baris |

## Nilai lama tidak pernah dihitung ulang

`[keputusan work owner]` Angka selisih yang berasal dari sistem lama **disalin apa adanya** saat
migrasi. Ia **tidak** dihitung ulang dengan rumus baru.

Sebabnya: sistem lama memuat **varian rumus** yang sengaja **tidak ditiru** - varian yang
mengurangi terhadap nilai yang ternyata kosong, sehingga aritmetiknya tampak ganjil. Menghitung
ulang berarti mengubah angka historis yang sudah dilaporkan.

`[penyimpangan sadar]` Perbedaan antara rumus baru dan varian lama itu **diterima dengan sadar**,
dan ditandai pada barisnya supaya terlihat, bukan disembunyikan.

## Akibat

1. Nol rumus selisih di basis data - nol view, nol kolom terhitung, nol trigger.
2. `repository` tidak menghitung apa pun; ia mengambil dua baris dan mengembalikannya.
3. Test membandingkan hasil lapisan layanan terhadap angka dokumen produksi, bukan terhadap tabel
   proyeksi - sebab tabel proyeksi adalah salinan, bukan sumber *(ADR-0020)*.
