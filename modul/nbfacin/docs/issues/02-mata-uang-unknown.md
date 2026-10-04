# 02: Mata uang yang tidak diketahui sebagai keadaan eksplisit

**What to build:** Sistem baru dapat membaca dan menampilkan nilai uang yang **tidak membawa mata
uang**, sebagaimana sistem lama melakukannya hari ini — tanpa menebak, dan tanpa berhenti.

Dasarnya terukur: korpus memang kehilangan informasi ini (112 layar menampilkan uang tanpa field mata
uang mana pun), dan pengukuran produksi D1 menemukan **4 dari 691.925 baris** tiba tanpa mata uang.
Langka, tetapi nyata — dan justru kelangkaannya yang berbahaya: cacat yang muncul 4 kali dari 691.925
**tidak akan tertangkap pengujian sampel**.

Karena itu keadaan tidak-diketahui dibuat **eksplisit dan terlihat**, lalu gagal keras **hanya di
titik yang benar-benar memerlukan jawabannya**. Default diam-diam ke satu mata uang ditolak: itu
menebak, dan tebakannya akan salah persis pada kasus paling mahal.

**Blocked by:** 01

**Status:** wontfix — ditutup work owner 01-10-2026 (butir 34, A9)

- [ ] Keadaan mata uang tidak diketahui adalah **nilai sah yang eksplisit**, bukan kosong dan bukan kegagalan
- [ ] Baca, tampilkan, dan simpan di memori: **diizinkan**
- [ ] Aritmetika antar nilai yang sama-sama tidak diketahui mata uangnya: **diizinkan**
- [ ] Aritmetika **lintas mata uang** yang melibatkannya: **gagal keras**
- [ ] **Tidak ada** default diam-diam ke mata uang mana pun di kode
- [ ] Implementasi menghasilkan **hitungan** berapa banyak nilai tiba tanpa mata uang — sebagai ukuran, bukan dugaan
- [ ] Kewajiban **"gagal keras saat menulis ke Oracle"** dicatat sebagai TODO yang terikat pada tiket jalur tulis produksi; **tidak** diimplementasikan sekarang karena seam-nya belum ada, dan **tidak** dihapus dari daftar kewajiban

## Comments

### 2026-10-01 — DIUSULKAN ditutup oleh keputusan, sisanya dicatat (agent)

⚠️ Penutupan tiket wewenang work owner (keputusan agent A9, menunggu konfirmasi). Usulannya: tiket ini
**digantikan keputusan work owner butir 3** (`../KEPUTUSAN-30-09-2026.md`): ikuti
`inti/backend/uang` apa adanya — tidak ada keadaan `Unknown`, lintas mata uang →
`uang.ErrMataUangBerbeda`, bukan `panic`. Tidak ada kode baru di tiket ini.

| Kriteria | Keadaan |
| --- | --- |
| Keadaan Unknown eksplisit | ⛔ digantikan butir 3: mata uang kosong ikut apa adanya (`TestPremiPAMataUangKosongTidakDiberiBawaan`) |
| Aritmetika lintas mata uang gagal keras | ✅ lewat `inti/backend/uang` (galat, bukan panic) |
| Tanpa default diam-diam | ✅ |
| Hitungan nilai tanpa mata uang | ⏸ belum ada pemuat data; dibuat saat pemuat `repository` NB lahir |
| TODO gagal keras saat tulis Oracle | ⏸ tetap kewajiban jalur tulis produksi; **tidak dihapus** |

⚠️ Catatan dari data: berkas kasus PA dan MBU nyata menyimpan `Currency` kosong pada coverage.

### 2026-10-01 — DITUTUP (keputusan work owner, butir 34)

Usulan A9 dikonfirmasi: tiket digantikan butir 3 (`../KEPUTUSAN-30-09-2026.md`). Kewajiban "hitungan nilai tanpa mata uang"
pindah ke `17-hitungan-uang-tanpa-mata-uang.md`, terikat ke pemuat `repository`. Kewajiban "gagal keras
saat tulis Oracle" tetap milik jalur tulis produksi — **tidak dihapus**, dicatat di tiket 17.
