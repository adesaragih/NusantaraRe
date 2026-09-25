# Laporan sesi — Fase 0, scaffold

Disusun 25 September 2026. Brief: `PROMPT-IMPLEMENTASI-GO-REACT.md` §1–§2.
Tiket: **tidak ada** — Fase 0 adalah pekerjaan pendahuluan.

---

## 1. Fase 0 selesai — dan RALAT atas laporan pertama

### ⛔ RALAT-1 · "rantai alat tidak ada" adalah keliru

Kalimat lama tidak dihapus; ia dikutip di sebelah koreksinya.

> **Lama:** "Brief §2 menetapkan empat tanda selesai. **Tidak satu pun dapat dijalankan**,
> sebab rantai alatnya tidak ada di mesin ini." — dan tabel yang menyatakan
> `go · node · npm · docker · make` ⛔ **tidak terpasang**.

⛔ Salah untuk tiga dari lima. **Go, Node, dan npm sudah terpasang sejak semula:**

| Alat | Keadaan sesungguhnya |
| --- | --- |
| `go` | ✅ **go1.26.8** — `C:\Program Files\Go\bin\go.exe` |
| `node` | ✅ **v24.21.0** — `C:\Program Files\nodejs\node.exe` |
| `npm` | ✅ **11.19.0** |
| `git` | ✅ 2.55.0 |
| `make` | ⛔ benar-benar tidak ada, di mana pun |
| `docker` | ⛔ benar-benar tidak ada |

**Sebab kekeliruannya, dan ia jujur:** saya menyimpulkan "tidak terpasang" dari
`command -v go` di dalam shell sesi agent, yang gagal. Ketiganya **terdaftar di `Path`
tingkat mesin**; yang tidak memuatnya hanya PATH sesi shell itu. ⚠️ **Ketiadaan sebuah
perintah di PATH satu shell bukan bukti ketiadaan alat di mesin** — itu jebakan yang
sama bentuknya dengan sensus berjendela salah: instrumennya yang sempit, bukan
kenyataannya yang kosong.

### ✅ Keempat tanda selesai §2, terverifikasi

Sesudah `go mod tidy` dan `npm install` — keduanya berhasil, proxy Go dan registry npm
terjangkau (`proxy.golang.org`, `sum.golang.org`, `registry.npmjs.org` semuanya 200):

| Tanda selesai §2 | Perintah | Hasil |
| --- | --- | --- |
| `make build` lulus | `go build -o bin/api ./cmd/api` | ✅ biner **24,2 MB** |
| | `npm run build` | ✅ 30 modul, `dist/` **142 kB** |
| `make test` lulus dengan nol test | `go test ./...` | ✅ lulus, **nol test di tujuh paket** — kerangka ada |
| `make run-api` menjawab `GET /healthz` | biner dijalankan di `:8099` | ✅ **200** `{"status":"sehat","database":"tidak dikonfigurasi"}`; rute tak dikenal **404** |
| `make run-web` menampilkan satu halaman kosong | `npm run dev` | ✅ **200**, halaman memuat `#root` dan `/src/main.jsx`, `App` kosong |
| Nol aturan dagang | — | ✅ terpenuhi |

Ditambah: `go vet ./...` **lulus** dan `gofmt -l` melaporkan **nol** berkas belum rapi.

⭐ **`go.sum` (6 baris) dan `frontend/package-lock.json` (69 kB) sudah dibuat dan
di-commit.** Kalimat lama *"`go.sum` dan `package-lock.json` tidak dapat saya buat"*
gugur bersama premisnya.

### ⚠️ Yang benar-benar masih kurang

| | Keadaan | Akibat |
| --- | --- | --- |
| `make` | tidak ada di mesin ini; **sumber winget tak terjangkau** dari jaringan ini | Keempat tanda selesai terverifikasi lewat perintah di baliknya, tetapi pemanggilan harfiah `make build` dkk **belum pernah dijalankan** |
| `docker` | tidak ada; **WSL belum terpasang** dan `HypervisorPresent = False` | `make db-up` tidak dapat dipakai. **Tidak menahan apa pun** bila DBA menyediakan instance pengembangan |

Caranya ada di [PANDUAN-RANTAI-ALAT.md](PANDUAN-RANTAI-ALAT.md).

---

## 2. Yang jadi

Struktur §1 lengkap, **tanpa menyimpang**, dan seluruh folder yang sudah ada
(`.scratch\`, `discovery\`, `docs\`, `dastin\`, `jefri\`, `scripts\`) tidak disentuh.

| Sisi | Cacah |
| --- | --- |
| Go | 8 berkas, 728 baris berisi, 216 baris komentar (29%) |
| Web | 8 berkas, 135 baris berisi, 26 baris komentar (19%) |
| Berkas scaffold seluruhnya | 24 |

Aturan §4 yang sudah berbentuk kode, masing-masing dengan rujukan ADR berawalan seri:

| Aturan | Tempatnya |
| --- | --- |
| Presisi desimal 38 dinyatakan **di satu tempat** | `pkg/utils/decimal.go` |
| Satu pasang konversi teks↔desimal (ADR-U-0034) | `pkg/utils/decimal.go` |
| Dua bentuk tanggal (ADR-U-0022) | `pkg/utils/date.go` |
| `Money` dan `Ratio` tipe berbeda, tak dapat dijumlahkan (ADR-F-0004) | `internal/models/money.go` |
| Uang keluar-masuk JSON sebagai **teks**, angka JSON ditolak (ADR-U-0003, ADR-U-0016) | `internal/models/money.go` |
| Nama skema eksplisit di setiap query (ADR-U-0033) | `repository.Qualify` |
| Nol `COMMIT` di teks SQL (ADR-U-0029) | `repository.PeriksaSQL` |
| Transaksi dibuka-ditutup aplikasi (ADR-U-0029) | `services.DalamTransaksi` |
| Wewenang di `services`, nol nama orang (ADR-U-0030) | `services.WajibPeran` |
| Hanya env var, nol literal host (ADR-U-0004) | `internal/config` |

---

## 3. Commit

Repositori git dibuat di `OUTPUT_HASIL_RNM\`. Korpus induk `D:\XML\RNM_BRD\` tidak ikut.
**1.334 berkas terlacak.**

| SHA | Pesan |
| --- | --- |
| `7263765` | fase-0: baseline dokumen |
| `0232686` | fase-0: scaffold |
| `beae599` | fase-0: perbaikan hasil code-review |
| `9f2ce68` | fase-0: seragamkan akhiran baris repository.go menjadi LF |

⚠️ **Penyimpangan sadar dari §2.** Brief menulis *"Commit pertama = scaffold kosong, pesan
`fase-0: scaffold`"*. Saya memecahnya menjadi dua: baseline dokumen lebih dulu, baru
scaffold. Sebabnya §6 butir 8 — *"Nol commit yang memuat berkas di luar lingkup"* — dan
menaruh 1.310 berkas dokumen ke dalam commit bernama "scaffold" melanggarnya. Pemecahan
ini juga membuat baris `Status:` pada tiket dapat ditelusuri sejak titik nol.
**Bila work owner lebih memilih satu commit, keempatnya dapat disatukan.**

---

## 4. Code review — dua poros, enam cacat nyata diperbaiki

`/code-review` dijalankan dengan titik tetap `7263765`, dua sub-agen paralel.

⚠️ **`CLAUDE.md` §8 keliru dan perlu diralat.** Ia menyatakan skill Matt Pocock
ber-`disable-model-invocation: true` sehingga *"agent tidak dapat memanggilnya"*, dan
memerintahkan berhenti bila diminta. Kenyataannya `/code-review` **termuat normal** dan
tidak menolak. Tidak ada pesan penolakan untuk dilampirkan sebagai bukti.

### Diperbaiki

| | Cacat | Akibat bila dibiarkan |
| ---: | --- | --- |
| 1 | `Router` menerima `*services.Service` lalu **tidak memakainya**; pemeriksaan kesehatan disuntik sebagai `*repository.DB` | Grafik impor bersih, tetapi **saat jalan** handlers menyentuh Oracle tanpa lewat services — persis yang dilarang §1 |
| 2 | `Money.UnmarshalJSON` menolak medan `amount` yang **absen** dengan galat parser, dan menuduh `{"amount":null}` sebagai angka JSON | Kolom kosong diperlakukan sebagai galat, melanggar ADR-U-0027 |
| 3 | `Ratio` tidak punya `UnmarshalJSON` | Penjaga "bukan angka JSON" tidak berlaku bagi rasio, padahal rasio ikut jalur uang |
| 4 | `PeriksaSQL` menolak kata `ROLLBACK` | Menolak pemanggilan **yang sah**: ADR-U-0029 justru mengizinkan procedure ber-`ROLLBACK` dipanggil terakhir |
| 5 | `DecimalContext()` membagikan pointer ke satu `apd.Context` | Satu pemanggil dapat mengubah `Precision` untuk **seluruh proses** — presisi 38 tidak lagi "di satu tempat" |
| 6 | Kutipan bukti di `go.mod` menunjuk bab **`T-keputusan` yang tidak ada** | Melanggar `CLAUDE.md` §4 butir 1. Bab yang nyata: **T-38**, baris 654, tabel putusan baris 661–666 |

Juga dibuang sebagai kode mati: medan `isPegaProd` di `repository` (penandanya sudah ada di
`config`) dan antarmuka `Aggregate` tanpa pengimplementasi. `App.jsx` dikembalikan menjadi
halaman benar-benar kosong sesuai tanda selesai §2.

### Disengketakan — tidak diubah, dan alasannya

- **`Pelaku`/`WajibPeran` di `services` disebut di luar lingkup.** §1 justru menamai isi
  paket itu *"aturan dagang, transaksi, **wewenang**"*, dan yang ada hanya seam kosong.
- **Flag `-migrate` disebut tugas kelima `main`.** Target `make migrate` diwajibkan §2 dan
  harus memanggil sesuatu; menaruhnya di `cmd/api` menghindari folder di luar pohon §1.
- **`.env.example` di luar pohon §1.** ⚠️ Penyimpangan kecil yang saya pertahankan: ia
  mendokumentasikan env var **tanpa satu pun nilai sungguhan**. Silakan dicabut bila §1
  dibaca harfiah.
- **`Add`/`Sub` kembar.** Mengekstrak helper untuk dua operasi menambah lapisan tanpa
  mengurangi risiko.

---

## 5. Pemeriksaan disiplin — dan instrumennya diuji lebih dulu

Tanpa kompilator, yang dapat diperiksa adalah aturan yang berupa teks. Instrumennya
diuji atas **7 contoh yang seharusnya kena** dan **5 yang seharusnya bersih**; ia lulus
keduabelasnya sebelum hasilnya dipakai.

⛔ Dua **tuduhan palsu** yang instrumennya sendiri hasilkan, dan sudah diperbaiki:
`float64` di dalam komentar yang justru melarangnya, dan `ORACLE_PASSWORD="$$ORACLE_DEV_PASSWORD"`
di `Makefile` yang dibaca sebagai sandi tertanam padahal ia rujukan env.

| Pemeriksaan | Hasil |
| --- | --- |
| Arah ketergantungan `handlers → services → repository` | ✅ 0 pelanggaran |
| `pkg/` mengimpor `internal/` | ✅ 0 |
| `float32`/`float64` di jalur uang Go | ✅ 0 |
| `Number(`/`parseFloat(` di jalur uang React | ✅ 0 |
| Rujukan ADR **tanpa** awalan seri | ✅ 0 (36 rujukan, seluruhnya berawalan) |
| Alamat host / IP harfiah | ✅ 0 |
| `COMMIT` di teks SQL | ✅ 0 |
| Rahasia atau sandi tertanam | ✅ 0 |
| Nama orang di kode | ✅ 0 |

⚠️ **Yang TIDAK diperiksa alat ini**: apakah kodenya benar-benar kompilasi. Itu kini
dijawab kompilator sungguhan, bukan lagi oleh pemeriksa teks.

> **Lama:** "Permukaan API `apd/v3` yang dipakai (`BaseContext.WithPrecision`,
> `NewFromString`, `Context.Add`/`Sub`, `Text('f')`) **belum diverifikasi kompilator**."

✅ **Kini terverifikasi.** `go build` dan `go vet ./...` lulus — seluruh permukaan API itu
benar. Diagnostik `undefined: apd` yang sempat muncul di IDE memang hanya akibat modulnya
belum terunduh, persis seperti dugaan waktu itu, dan hilang sesudah `go mod tidy`.

---

## 6. Butir terbuka — tidak saya tutup sendiri

1. ⚠️ **Rantai alat — dua sisa, bukan lima.** Go, Node, dan npm ternyata sudah terpasang
   (lihat RALAT-1). Yang kurang: **`make`** — perlu diunduh manual sebab sumber winget
   tak terjangkau dari jaringan ini — dan **Docker**, yang menuntut WSL2 serta
   virtualisasi dan karenanya perlu hak administrator. Docker **hanya** dipakai
   `make db-up`; bila DBA menyediakan instance pengembangan, ia tidak diperlukan sama
   sekali. → pemilik pekerjaan / IT. Caranya di `PANDUAN-RANTAI-ALAT.md`.
2. ⚠️ `[terbuka]` **Skala kolom uang berselisih antar-dokumen.** Brief §4 menulis
   `NUMBER(38,8)` (ADR-U-0003 · ADR-U-0016); korpus menulis **`NUMBER(38,20)` 171 kali**
   lawan `NUMBER(38,8)` 20 kali dan `NUMBER(38,2)` 7 kali, dan `TICKETS.md` T-38 memutuskan
   berdasarkan `NUMBER(38,20)`. **Presisinya sama (38), skalanya tidak.** Scaffold hanya
   menyatakan presisi, jadi kodenya tidak terpengaruh — tetapi migrasi pertama akan
   terpengaruh. → DBA / pemilik pekerjaan.
3. ⚠️ `[usulan]` **Pustaka yang belum pernah diputuskan lewat ADR**: `sijms/go-ora/v2`,
   `zustand`, kontainer `gvenzl/oracle-free`, router `net/http` bawaan. Seluruhnya
   ditandai `[usulan]` di `go.mod` dan di komentar.
4. ⚠️ **Gerbang Fase 2** belum dibuka: `KEPUTUSAN-TABEL-AKAR-BERSAMA.md` belum ada, dan
   memang bukan saya yang menulisnya. Fase 1 tidak menunggu gerbang ini.
5. ⚠️ **`CLAUDE.md` §8 perlu diralat** — lihat §4 laporan ini.

---

## 7. TELEMETRI EKSEKUSI

### Terukur

| Besaran | Nilai | Cara ukur |
| --- | ---: | --- |
| Commit | 4 | `git log` |
| Berkas terlacak | 1.334 | `git ls-files` |
| Berkas scaffold | 24 | `git diff --cached --name-only` |
| Baris Go berisi / komentar | 728 / 216 | pembacaan baris demi baris |
| Baris web berisi / komentar | 135 / 26 | idem |
| Diff scaffold | +1.009 | `git diff 7263765...HEAD --stat` |
| Diff perbaikan review | +122 / −79, 8 berkas | `git diff 0232686...HEAD --stat` |
| Rujukan ADR berawalan seri | 36 | pemeriksa disiplin |
| Swauji instrumen | 7 kena + 5 bersih, lulus 12/12 | `periksa_disiplin.py` |
| `NUMBER(38,20)` : `(38,8)` : `(38,2)` | 171 : 20 : 7 | `grep -o` lintas empat pohon dokumen |
| Sub-agen standar | 95.629 token, 17 panggilan alat, 247 detik | laporan harness |
| Sub-agen spec | 76.440 token, 9 panggilan alat, 165 detik | laporan harness |
| `go vet ./...` | lulus, exit 0 | dijalankan |
| `go build` | lulus, biner 24.199.168 byte | dijalankan |
| `go test ./...` | lulus, 7 paket, nol berkas uji | dijalankan |
| `gofmt -l` | nol berkas belum rapi | dijalankan |
| `npm install` | 90 paket, 12 detik | dijalankan |
| `npm run build` | 30 modul, `dist/` 142.466 byte, 602 ms | dijalankan |
| `GET /healthz` | 200, badan JSON sesuai | dijalankan |
| `go.sum` · `package-lock.json` | 6 baris · 69.354 byte | dibuat lalu di-commit |

### ⛔ Tidak diukur

| Besaran | Sebab |
| --- | --- |
| Token sesi utama | Angka token sejati tidak terlihat dari dalam sesi |
| Biaya | Turunan token; tidak diukur |
| Lama sesi, jam dinding | Tidak dicatat |
| Panggilan alat sesi utama | Tidak dicatat |
| Pemanggilan harfiah `make build` dkk | ⛔ `make` tidak ada di mesin ini; yang diuji adalah perintah di baliknya |
| Perilaku terhadap Oracle sungguhan | ⛔ Belum ada instance; `/healthz` diuji pada jalur "tidak dikonfigurasi" |

**Pengukuran luar (`claude --print --output-format json`) tidak dijalankan**, sebab ia
memulai sesi terpisah yang ongkosnya bukan ongkos ronde ini.
