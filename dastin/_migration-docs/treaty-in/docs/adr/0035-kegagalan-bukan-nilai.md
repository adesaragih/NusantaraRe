# ADR-0035 — Kegagalan tidak pernah disamarkan menjadi nilai

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In, dan diusulkan untuk seluruh modul

## Konteks

Lima temuan terpisah di modul ini ternyata satu penyakit:

| Keadaan | Yang dilakukan sistem lama |
|---|---|
| Limit nol saat menghitung Rate on Line | `ROLPct` diisi `9989998`, tersimpan sebagai `998999800` |
| Mata uang di luar slot IDR/USD pada reinstatement | hasil diisi `0` |
| Kurs valas tidak ditemukan | pengali tinggal kosong, nilai IDR menjadi `0` |
| Prosedur Oracle gagal | `ErrMsg` dan `StsSave` tidak diisi; Pega menerima kosong |
| Persentase tidak terurai dari `VARCHAR2(1000)` | nilai lolos apa adanya |

Nilai samarannya selalu tampak masuk akal — nol, satu, atau angka besar — sehingga hilir tidak
punya cara membedakannya dari hasil sungguhan.

## Keputusan

1. Besaran yang **tidak dapat dihitung** tidak pernah diganti dengan nilai pengganti. Ia bernilai
   tidak-ada, dan tipe datanya harus **memaksa** pembaca menangani ketidakadaan itu.
2. Ketidakadaan itu **membawa sebabnya**: limit nol, kurs tidak ditemukan, mata uang tidak dikenal,
   nilai tidak terurai. Sebab yang berbeda menuntut tindakan yang berbeda.
3. Nilai pengganti yang netral — nol, satu, string kosong — **dilarang secara khusus**, karena
   justru nilai netral yang paling mudah lolos tanpa disadari.
4. Di perbatasan tempat masukan tidak bisa ditolak, kegagalan **diangkat ke pemanggil**, tidak
   ditelan.

## Catatan penerapan

Aturan ini **tidak** berarti setiap ketidakmampuan menghitung menghentikan pekerjaan orang. Kontrak
yang limitnya belum diisi memang belum bisa punya Rate on Line, dan itu normal selama penyusunan.

Yang dilarang adalah **menyimpan angka palsu** untuk keadaan itu. Bedakan:

- *belum bisa dihitung, dan itu wajar di tahap ini* — bukan angka, tidak menghalangi akseptasi;
- *gagal dihitung padahal seharusnya bisa* — bukan angka, **menghalangi** akseptasi.

## Instans yang tertutup oleh ADR ini

Selain kelima di atas: pengambilan `pxResults(1)` tanpa memeriksa ada berapa hasil; pembagian
dengan bagian NuRe tanpa memeriksa ia terisi; dan penjaga rasio kerugian yang memeriksa premi
bruto sementara pembaginya premi neto. Semuanya dicatat sebagai instans, bukan temuan baru.
