# ADR-0055 — Daftar keadaan siklus hidup dan perpindahan yang sah

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In
**Melaksanakan:** ADR-0046 (satu keadaan), ADR-0052 (satu rantai persetujuan), ADR-0054 (warisan)

## Konteks

ADR-0046 menetapkan **bahwa** ada satu keadaan bertipe tegas dengan himpunan tertutup. Ia tidak
menetapkan **apa saja isinya**. ADR ini mengisinya, dan mengisinya dalam bentuk daftar — karena
keadaan yang tidak punya perpindahan masuk maupun keluar adalah keadaan yang tidak pernah terjadi,
dan itu hanya ketahuan dari daftar, tidak dari uraian.

Dua hal ditemukan justru karena daftarnya disusun lebih dulu, dan keduanya argumen dari data:

**`Position` kosong menandai dua keadaan yang berlawanan.** Di sistem lama, `Resolve Complete`
(disetujui) dan `Decline` (ditolak) sama-sama menyetel `Position` ke nilai kosong. Hanya
pasangannya dengan `StatusAkseptasi` yang membedakan kontrak yang disetujui dari kontrak yang
ditolak. Itu bukan alasan selera untuk menyatukan kelima bendera menjadi satu keadaan — itu alasan
dari data.

**Yang menggantikan jalan pintas revisi bukan aturan baru, melainkan bentuk.** Di sistem lama,
`RevisionState == 1` membuat persetujuan Sec Head langsung menghasilkan `Resolve Complete`,
memotong dua tingkat. ADR-0052 menghapus jalan pintas itu. Yang menggantikannya bukan aturan
pengganti: perubahan atas versi yang sudah disetujui **melahirkan versi baru yang mulai dari
`DRAFT`**. Versi baru itu melewati keempat tingkat karena ia memang versi baru, bukan karena ada
aturan yang memaksanya.

## Keputusan

### Daftar keadaan — enam, ditambah satu keadaan warisan

| Keadaan | Padanan lama (`Position` + `StatusAkseptasi`) | Terminal |
|---|---|---|
| `DRAFT` | kosong atau `ReasTreatyInAdmin`, status kosong **atau** `Reject` | tidak |
| `MENUNGGU_SEC_HEAD` | `ReasTreatyInSecHead` + `Accept` | tidak |
| `MENUNGGU_DEPT_HEAD` | `ReasTreatyInDeptHead` + `Accept` | tidak |
| `MENUNGGU_DIREKTUR` | `ReasTreatyInDirector` + `Accept` | tidak |
| `DISETUJUI` | kosong + `Resolve Complete` | **ya** |
| `DITOLAK` | kosong + `Decline` | **ya** |
| `WARISAN_TAK_TERPETAKAN` | nilai apa pun di luar daftar di atas (ADR-0054) | tidak |

### Perpindahan yang sah — dua belas  → **TIGA BELAS sejak perubahan 24 September 2026**

> Tabel di bawah adalah daftar sebagaimana diputuskan 23 September 2026. Satu perpindahan
> ditambahkan kemudian: **`DRAFT` — `BATALKAN` → `DIBATALKAN`**. Lihat bagian
> "Perubahan 24 September 2026" di kaki berkas ini.

| Dari | Peristiwa | Ke |
|---|---|---|
| (versi baru lahir) | `LAHIR` | `DRAFT` |
| `DRAFT` | `AJUKAN` | `MENUNGGU_SEC_HEAD` |
| `MENUNGGU_SEC_HEAD` | `SETUJUI` | `MENUNGGU_DEPT_HEAD` |
| `MENUNGGU_SEC_HEAD` | `KEMBALIKAN` | `DRAFT` |
| `MENUNGGU_SEC_HEAD` | `TOLAK` | `DITOLAK` |
| `MENUNGGU_DEPT_HEAD` | `SETUJUI` | `MENUNGGU_DIREKTUR` |
| `MENUNGGU_DEPT_HEAD` | `KEMBALIKAN` | `DRAFT` |
| `MENUNGGU_DEPT_HEAD` | `TOLAK` | `DITOLAK` |
| `MENUNGGU_DIREKTUR` | `SETUJUI` | `DISETUJUI` |
| `MENUNGGU_DIREKTUR` | `KEMBALIKAN` | `DRAFT` |
| `MENUNGGU_DIREKTUR` | `TOLAK` | `DITOLAK` |
| `WARISAN_TAK_TERPETAKAN` | `PERBAIKAN_WARISAN` | keadaan sah yang dipilih eksplisit (ADR-0054) |

Setiap keadaan punya sedikitnya satu perpindahan masuk dan, kecuali yang terminal, sedikitnya satu
perpindahan keluar. Tidak ada keadaan yatim.

## Empat keputusan yang tertanam di daftar itu

### 1. `ReasTreatyInGroupLeader` tidak ada, dan itu dibuktikan bukan diduga

Di sistem lama, tingkat Group Leader punya **tiga perpindahan keluar dan nol perpindahan masuk**.
Sapuan menyeluruh atas ekspor menunjukkan **tidak ada satu pun aturan** yang pernah menyetel
`Position` ke nilai itu. Ia keadaan yang tidak pernah terjadi.

ADR-0052 sudah memutuskan tidak membawanya atas dasar ketiadaan bukti. Sekarang dasarnya lebih
kuat: bukan tidak ditemukan pemakaiannya, melainkan **tidak ada jalan masuknya sama sekali**.

### 2. `DIKEMBALIKAN` bukan keadaan — ia `DRAFT` yang punya riwayat penolakan

Sistem lama punya nilai `Reject` yang berdiri di samping `ReasTreatyInAdmin`. Pertanyaannya bukan
apakah nilai itu ada, melainkan apakah ia berbeda dalam hal apa pun **yang bukan sejarah**. Tiga
hal diperiksa, dan ketiganya dijawab oleh sapuan menyeluruh atas ekspor:

| Yang diperiksa | Hasil |
|---|---|
| apakah aturan kelengkapan yang berlaku berbeda | tidak — tidak ada aturan kelengkapan yang menyebut `Reject` |
| apakah siapa yang boleh menyuntingnya berbeda | tidak — keduanya `ReasTreatyInAdmin` |
| apakah apa yang boleh diubah berbeda | tidak — keterbukaan field dikendalikan `ViewState`, yang tidak pernah membaca `Reject` |

Kata `Reject` hanya muncul di **dua berkas** pada seluruh ekspor: aturan persetujuan
`Akseptasi_DT` dan satu kondisi tombol di `TreatyInActionButtons` — dan pada kondisi tombol itu ia
diperlakukan **identik** dengan `Accept`.

Maka `DIKEMBALIKAN` tidak disimpan sebagai keadaan. Riwayat penolakannya sudah tersimpan di
`CATATAN_PERSETUJUAN`; menyimpannya lagi sebagai keadaan berarti menyimpan satu fakta di dua
tempat, yang ADR-0041 larang. Ini bentuk yang sama dengan `ViewState`: sesuatu yang bisa diturunkan,
disimpan sebagai penanda tersendiri — dan prinsip yang sama membuangnya.

**Yang tetap bisa dilakukan:** menyaring daftar kerja untuk "kontrak yang dikembalikan kepada saya"
tetap mungkin. Ia **turunan** dari keputusan terakhir di `CATATAN_PERSETUJUAN`, dihitung saat
ditanya (ADR-0037), bukan bendera yang disimpan.

**Apa yang membatalkan keputusan ini:** bila bisnis menyatakan kontrak yang dikembalikan menanggung
kewajiban yang tidak ditanggung draft baru — misalnya wajib menyertakan tanggapan atas alasan
pengembalian sebelum boleh diajukan ulang — maka perbedaan itu nyata, `DIKEMBALIKAN` menjadi
keadaan tersendiri, dan **kewajiban itulah yang ditulis sebagai alasannya**. Yang tidak boleh:
dua keadaan bertahan hanya karena sistem lama punya dua nilai.

### 3. `DITOLAK` terminal bagi versinya

**Ini keputusan perancang, bukan warisan sistem lama.** Di sistem lama `Decline` menyetel `Position`
ke kosong, sehingga pengajuan berikutnya masuk lagi dari Sec Head: penawaran yang sudah ditolak
bisa dihidupkan kembali tanpa jejak bahwa ia pernah ditolak.

Membuatnya terminal memaksa kebangkitan itu menjadi **versi baru yang terlihat**.

**Konsekuensi yang dinyatakan sadar, bukan akibat samping:** sebuah kontrak yang seluruh versinya
`DITOLAK` **tetap ada sebagai kontrak**. Ia punya identitas, punya kunci alami, dan **peringatan
duplikat akan menyebutnya** ketika penawaran serupa masuk lagi. Itu memang yang dikehendaki —
penawaran yang pernah ditolak untuk cedant, asal bisnis, periode, dan sifat proporsi yang sama
justru hal yang paling layak diperingatkan kepada orang yang memasukkannya.

**Apa yang membatalkan keputusan ini:** bila bisnis menyatakan penolakan boleh dicabut pada versi
yang sama, perpindahan `DITOLAK` — `AJUKAN ULANG` — `MENUNGGU_SEC_HEAD` ditambahkan. Satu baris,
bukan perubahan bentuk.

### 4. Tidak ada pintu yang membatalkan keadaan terminal

Di sistem lama, `DISETUJUI` terminal di dalam mesin persetujuannya — kedua cabang luar `Akseptasi_DT`
menuntut `StatusAkseptasi != 'Resolve Complete'` — tetapi ada **empat aturan di luar mesin** yang
membatalkannya dari samping:

| Aturan | Yang dilakukannya | Terlihat oleh |
|---|---|---|
| `TreatyInForceEdit` | membuka kunci kontrak yang sudah disetujui | **setiap pengguna layar penawaran** |
| `TreatyInForceResolveComplete` | menyetel selesai disetujui tanpa approver | dua nama pengembang |
| `TreatyInReturntoInputor` | mengembalikan ke admin, status dikosongkan | dua nama pengembang |
| `TreatyInSetToDirector` | melompati dua tingkat | lewat `TreatyInTestAgent` di layar penawaran |

Di model baru tidak ada satu pun jalan menyetel keadaan selain melalui perpindahan di daftar di
atas. Keadaan bukan kolom yang bisa ditulis; ia hasil dari perpindahan yang sah.

## Konsekuensi

- Daftar ini adalah sumber bagi kelengkapan per-perpindahan di `SPEC-INVARIAN.md`: setiap
  perpindahan punya daftar syarat yang harus terpenuhi sebelum ia boleh terjadi.
- Keadaan disimpan pada **versi kontrak**, bukan pada kontrak. Kontrak tidak punya keadaan; ia
  punya versi-versi yang masing-masing punya keadaan (ADR-0040).
- `ViewState`, `IsEditData`, `RevisionState`, `Position` dan `StatusAkseptasi` tidak ada di model
  baru, tidak dinamai ulang, dan tidak direkonsiliasi.

---

# Perubahan 24 September 2026 — keadaan kedelapan: `DIBATALKAN`

**Status:** diterima, 24 September 2026. Diputuskan pemilik proses.
**Sebab:** lubang L-7, ditemukan saat menyusun `5-tiket/DAFTAR-PEKERJAAN.md`.

## Lubang yang ditutup

Daftar dua belas perpindahan di atas memberi `DRAFT` **tepat satu** jalan keluar: `AJUKAN`.
Digabung INV-25 — paling banyak satu versi tak-terminal per kontrak — akibatnya:

> Pengisi kontrak yang membuat versi karena salah pencet, atau memulai addendum yang ternyata tidak
> jadi, **tidak punya jalan keluar**. Kontraknya terkunci: versi itu tidak dapat dibuang, dan versi
> kedua tidak boleh dibuat.

Satu-satunya jalan yang tersisa adalah **mengajukannya supaya ada yang menolaknya** — memakai jalur
persetujuan sebagai tempat sampah.

**Sistem lama juga tidak punya jalan keluar yang bersih**, dan yang dipakainya merusak:
`TreatyInDeclineConfirmation_postactEDM` langkah 6 dan 7 — keduanya **hidup** — **menghapus baris**
di `M_TREATY_IN_EDM` dan `TREATY_IN_EDM` setelah addendum ditolak. Jalur kontrak biasa
(`TreatyInDeclineConfirmation_postact`, empat langkah hidup) tidak menghapus apa pun.

## Keputusan

**Keadaan kedelapan `DIBATALKAN`, terminal, dengan satu perpindahan masuk.**

| Dari | Peristiwa | Ke |
|---|---|---|
| `DRAFT` | `BATALKAN` | `DIBATALKAN` |

Daftar keadaan menjadi **tujuh keadaan sah + satu keadaan warisan**; perpindahan menjadi
**tiga belas**.

### Kenapa bukan memakai `DITOLAK` yang sudah ada

Meski itu tidak menambah keadaan: **`DITOLAK` berarti seseorang yang berwenang menolak.** Pengisi
yang membuang drafnya sendiri tidak ditolak siapa pun. Memakai satu keadaan untuk keduanya membuat
**setiap hitungan penolakan tercemar draf yang dibuang**, dan tidak ada cara memisahkannya kemudian.

Keadaan menyatakan apa yang terjadi; ia harus menyatakan yang benar.

### Kenapa bukan melonggarkan INV-25

Melonggarkannya berarti **membuang invarian untuk menyelesaikan masalah alur kerja**. Hasilnya draf
terlantar yang menumpuk selamanya tanpa ada yang tahu mana yang sungguhan.

## Tiga syarat yang mengikat

| # | Syarat | Catatan |
|---|---|---|
| **a** | **Baris tidak dihapus.** Versi yang dibatalkan tetap tersimpan, tetap membawa nomornya, tetap membawa jejaknya | jawaban langsung atas apa yang dilakukan sistem lama |
| **b** | **Nomor revisinya tidak dipakai ulang.** Lompatan penomoran itu jujur; nomor yang dipakai dua kali tidak | **tidak menuntut mekanisme baru** — lihat di bawah |
| **c** | **Hanya dari `DRAFT`, dan hanya oleh pembuatnya** | membatalkan sesudah diajukan adalah hal yang **berbeda** |

### Syarat (b) jatuh sendiri dari (a)

**INV-04 sudah menetapkan `NOMOR_URUT_VERSI` unik di dalam satu `KONTRAK`.** Selama barisnya tidak
dihapus — syarat (a) — nomor yang sudah terpakai tetap menempati tempatnya, dan constraint yang
sudah ada menolak pemakaian ulangnya.

Jadi (b) **bukan aturan tambahan**; ia akibat (a) di bawah invarian yang sudah berlaku. Tidak ada
trigger baru, tidak ada kolom baru. Yang muncul hanyalah **lompatan** dalam deret nomor, dan
lompatan itu memang yang dikehendaki.

### Apa yang TIDAK diputuskan di sini

**Penarikan sesudah diajukan** — membatalkan versi yang sudah berada di antrian persetujuan —
**bukan bagian keputusan ini.** Ia kemampuan tersendiri: pelakunya berbeda, dan akibatnya menyentuh
orang lain yang sudah mulai menilai. Bila ia diperlukan, ia pertanyaan tersendiri dengan jawabannya
sendiri.

Menggabungkannya sekarang adalah larangan *"satu tiket memuat hal yang pasti dan hal yang belum
diputuskan"*, diterapkan pada keputusan alih-alih pada tiket.

## Konsekuensi yang harus dikerjakan sesi to-spec

| Berkas | Yang berubah |
|---|---|
| `SPEC-INVARIAN.md` **INV-20** | himpunan keadaan menjadi **delapan** nilai, bukan tujuh |
| `SPEC-INVARIAN.md` **INV-23** | tidak ada perpindahan keluar dari `DISETUJUI`, `DITOLAK`, **maupun `DIBATALKAN`** |
| `SPEC-INVARIAN.md` **INV-25** | `DIBATALKAN` terhitung **terminal**, sehingga versi yang dibatalkan tidak lagi menahan pembuatan versi berikutnya |
| `SPEC-INVARIAN.md` §3 | baris kelengkapan baru: `DRAFT` → `DIBATALKAN` — pelakunya **pembuat versi itu sendiri** |
| `SPEC-MODEL-DATA.md` §10.2 | `KEADAAN_SIKLUS_HIDUP` bertambah satu nilai |
| `5-tiket/DAFTAR-PEKERJAAN.md` | satu kemampuan baru, bergolongan **BARU** |

## Satu akibat yang BELUM diputuskan, dan sengaja dibiarkan terbuka

**Kontrak yang seluruh versinya `DIBATALKAN`.** Bila versi pertama sebuah kontrak dibatalkan dan
tidak pernah ada versi kedua, kontrak itu ada tetapi tidak pernah punya isi yang berlaku.

Pertanyaannya — apakah ia muncul di pencarian (`5-tiket/DAFTAR-PEKERJAAN.md` P-56), apakah ia
terhitung sebagai kontrak dalam laporan apa pun — **tidak dijawab di sini**, karena ia pertanyaan
bisnis, bukan pertanyaan bentuk. Dicatat supaya tidak ditemukan belakangan sebagai kejutan.
