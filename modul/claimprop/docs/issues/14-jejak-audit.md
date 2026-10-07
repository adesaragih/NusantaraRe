# 14: Jejak audit klaim — empat cacat yang diperbaiki

**Status:** sebagian 07-10-2026 — log layanan `MONITORING_KLAIM_LOG` dibangun; jejak Claim History tidak disimpan (OQ-CP-17) — RALAT 07-10-2026 (semula `ready-for-agent`)
**AC yang dipegang (RALAT 07-10-2026, prompt §7 butir 4):** AC 133 (wewenang hapus kronologi lewat peran) — layar Claim Prop tidak punya aksi hapus kronologi yang berbukti; berlaku bila Claim History disimpan (OQ-CP-17).
**Blocked by:** 00 (PREFACTOR) · 10 (wewenang)
**Menutup:** AC 77 · 78 · 79 · 80 · 81 · 82 · 83 *(7 AC)* — US 59–62

## Hasil & nilai pengguna

Auditor dapat merekonstruksi klaim dari awal sampai akhir: setiap aksi meninggalkan entri berisi
aksi, pelaku, waktu, dan tingkat wewenang; riwayat tampil urut waktu; **semua** peran terekam tanpa
kecuali; dan memo penutupan klaim benar-benar sampai ke jejak.

⚠️ Keempat cacat di bawah **membalik** perilaku Pega. Itu disengaja dan disetujui — tanpa perbaikan
ini jejak audit tidak dapat dipercaya sebagai bukti.

## Area codebase

Entitas jejak audit · enum tingkat wewenang · migrasi entri warisan · pengurutan riwayat.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `DataTransform/InsertChronology_DT.xml` | 1 | ⚠️ digerbangi jabatan pelaku, dan **seluruh penulisan bersarang di bawahnya** → aksi peran itu **nol entri** |
| `DataTransform/InsertChronology_DT.xml` | 1.1.1 | membaca slot memo pada halaman kronologi |
| `DataTransform/InsertChronology_DT.xml` | 1.1.4–1.1.7 | tingkat wewenang ditulis ke properti berawalan `Is…` yang **bukan boolean** |
| `Activity/SethistoryKlaimTreaty.xml` | — | penulis kedua; **tidak** punya pengecualian jabatan, dan mengeja tingkat terbawah **berbeda** |
| `Activity/CloseClaimProp.xml` | 5 | ⚠️ memo penutupan ditulis ke slot yang **tidak pernah dibaca** |
| `Activity/MakeLowercase_Act.xml` | 1 | ⭐ **BARU 2026-09-19 (ronde 6)** — ⚠️ **lokasi kerugian dan uraian laporan ditimpa versi huruf kecilnya**; **nilai asli tidak disimpan di mana pun**. Teks bebas yang diketik pengguna hilang bentuk aslinya sebelum sempat tercatat |
| `Activity/RemoveLossAlloction_act.xml` | 2 | ⭐ **BARU** — menulis kronologi *"Delete Loss Allocation"*; satu jenis peristiwa jejak audit yang belum terdaftar |
| `Activity/AddInterest_act.xml` | — | ⭐ **BARU** — catatan pengembangnya *"Untuk History Claim"*, penulis kronologi ketiga di luar kedua penulis yang sudah tercatat |

⚠️ `[terverifikasi]` Sensus **329 berkas berlingkup halaman kronologi**: slot yang dibaca muncul
**34 kali di 27 berkas**; slot yang ditulis step 5 muncul **satu kali di satu berkas** — yaitu langkah
yang salah itu. **Lingkup halaman wajib disebut:** grep nama telanjang memberi 315 dan 57, karena
slot itu juga dipakai sebagai parameter umum di rule SQL lain.

## ADR terkait

**ADR-0014** (tingkat wewenang — sumbernya dari tiket 10).

## Acceptance criteria

- [ ] `[terverifikasi]` Setiap aksi meninggalkan entri berisi **aksi, pelaku, waktu, dan tingkat wewenang** *(AC 77 spec)*
- [ ] ⚠️ Tingkat wewenang disimpan sebagai **enum tertutup**, dan **nama pelaku di kolom terpisah**. **Alasan menyimpang:** di Pega satu properti berawalan `Is…` menyimpan tingkat wewenang, bukan boolean — namanya berbohong dan isinya tidak dapat disaring *(AC 78 spec)*
- [ ] ⚠️ Migrasi memetakan **kedua ejaan** tingkat terbawah ke satu nilai enum, dan memindahkan nama pelaku ke kolom pelaku. **Alasan menyimpang:** dua penulis mengeja tingkat terbawah berbeda, dan korpus sendiri mengakalinya dengan pencocokan substring *(AC 79 spec)*
- [ ] ⚠️ **Semua peran terekam tanpa kecuali**; test yang menemukan satu peran tidak berjejak **gagal**. **Alasan menyimpang:** di Pega satu jabatan dikecualikan dan seluruh penulisan bersarang di bawah gerbang itu, sehingga aksinya menghasilkan nol entri *(AC 80 spec)*
- [ ] ⚠️ **Memo penutupan klaim tercatat di riwayat.** **Alasan menyimpang:** di Pega memo ditulis ke slot yang tidak pernah dibaca, sehingga alasan penutupan hilang seluruhnya *(AC 81 spec)*
- [ ] ⚠️ `[data DBA]` Riwayat ditampilkan **urut dari kolom tanggal**, bukan urutan baris. **Alasan menyimpang:** Pega menampilkan **urutan baris tersimpan** — grid jejak audit tidak punya konfigurasi pengurutan sama sekali, baik pada `Section/InputAcceptation.xml` maupun `Section/OutstandingClaim.xml`, sehingga mengurutkan menurut tanggal **mengubah apa yang dilihat pengguna**. Ini **penyimpangan sadar**. **Bukti:** satu baris contoh memuat entri terakhir bertanggal lebih lambat dari sebelas lainnya, jadi urutan tersimpan memang bukan kronologis *(AC 82 spec)*
- [ ] `[terverifikasi]` Jejak audit mencatat **jenis aksi**, bukan nilai sebelum/sesudah — batas yang diwarisi dan **tidak** diperluas dalam spec ini *(AC 83 spec)*

## Perintah verifikasi

```
jalankan test "aksi oleh tiap peran -> semuanya menghasilkan entri jejak"
jalankan test "tutup klaim dengan memo -> memo tercatat dan terbaca kembali"
jalankan test "entri bertanggal tidak urut -> riwayat tampil urut tanggal"
jalankan test "tingkat wewenang tersimpan sebagai enum tertutup, bukan teks bebas"
jalankan test "migrasi: dua ejaan tingkat terbawah -> satu nilai enum, nama ke kolom pelaku"
cari nama pelaku yang tercampur ke kolom tingkat          -> nihil
```
