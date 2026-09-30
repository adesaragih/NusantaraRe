# 08: Terapkan ke semua — pratinjau, konfirmasi, jejak audit, dan laporan sebagian-gagal

**Status:** ready-for-agent

**Blocked by:** 07 (fitur ini bekerja atas baris business)

## Hasil & nilai pengguna

Sebagai **admin master**, saya menerapkan satu business ke **seluruh kontrak berjenis reasuransi
sama** sekaligus — tetapi hanya setelah sistem memberi tahu **berapa baris akan terpengaruh** dan
saya menyetujuinya, dan hasilnya tercatat sehingga dapat ditelusuri.
*(User story 24–26 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Kueri baris terdampak; penerapan per baris |
| `internal/services` | Pratinjau; orkestrasi penerapan; **penghitungan berhasil/gagal** |
| `internal/handlers` | Endpoint pratinjau; endpoint eksekusi; respons memuat rekap |
| `frontend/` | Dialog konfirmasi dengan jumlah baris; ringkasan hasil |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveBusinessToAllLife_Act` | `ASM-FW-GISFW-…` / `SAVEBUSINESSTOALLLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveBusinessToAllLife_Act.xml` | orkestrator penerapan massal |
| `SaveTreatyBusinessAll_Life_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!SAVETREATYBUSINESSALL_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveTreatyBusinessAll_Life_SQL.xml` | penulisan massal |

`[terverifikasi]` Gerbang pemilihan baris: **`.REINSTYPEID == Param.REINSTYPEID`** — penerapan
menyasar seluruh baris berjenis reasuransi sama.

`[terverifikasi]` Di Pega **tidak ada konfirmasi dan tidak ada jejak audit** — satu klik dapat
mengubah puluhan baris tanpa peringatan dan tanpa jejak.

⚠️ `[data DBA]` **`COMMIT` berada di dalam tiap procedure** — penerapan massal karena itu **tidak
atomik**. Bila gagal di tengah, sebagian baris **sudah tersimpan** dan tidak dapat di-rollback.

## ADR terkait

**ADR-0007** (jejak audit — inti tiket ini), **ADR-0003**, **ADR-0015** (batas transaksi dipegang
aplikasi; di sini **tidak ada** transaksi menyeluruh yang mungkin).

## Acceptance criteria

- [ ] ⚠️ **Terapkan ke semua** menampilkan **jumlah baris yang akan terpengaruh** dan **menunggu
      konfirmasi** sebelum dijalankan. *(AC 28 spec; `[keputusan work owner]` — penyimpangan sadar 5)*
- [ ] Pratinjau menyebut **jenis reasuransi** yang menjadi dasar pemilihan, sehingga pengguna tahu
      cakupannya.
- [ ] ⚠️ Hasil penerapan massal **tercatat di jejak audit**, termasuk **berapa baris berubah** dan
      **siapa** pelakunya. *(AC 29 spec; `[keputusan work owner]`)*
- [ ] Penerapan massal yang **dibatalkan** pada dialog konfirmasi **tidak mengubah apa pun**.
      *(AC 30 spec)*
- [ ] ⚠️ Bila penerapan gagal di tengah, sistem melaporkan **berapa berhasil dan berapa gagal** —
      **tidak** menampilkan sukses tunggal yang menyesatkan. *(AC 31 spec; `[data DBA]` tiap
      procedure commit sendiri)*
- [ ] Baris yang gagal **disebutkan** — pengguna tahu mana yang perlu diulang.
- [ ] `o_message` diperiksa **per baris**, bukan hanya pada baris terakhir. *(AC 46 spec)*
- [ ] Pratinjau yang menunjukkan **nol baris** memberi tahu pengguna, dan tidak menampilkan dialog
      konfirmasi yang sia-sia.

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Sifat tidak-atomik harus terlihat, bukan disembunyikan.** `[data DBA]` Karena `COMMIT` ada di
dalam tiap procedure, tidak ada cara membungkus penerapan massal dalam satu transaksi. Menyajikannya
seolah "semua atau tidak sama sekali" akan **berbohong kepada pengguna**. Yang benar: laporkan
hasilnya apa adanya.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — sifat commit per baris tidak dapat difake
dengan jujur, dan justru itulah yang diuji.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 30-09-2026 — sesi implementasi (paket 0)

> Sumber: `RALAT-DEV-30-09-2026.md` (K1–K8 katalog DEV, R1–R12 pembacaan ulang XML) dan `PARITAS-LAYAR-DAN-AKSI.md`. Kalimat di atas **tidak dihapus**; yang berlaku adalah ralat ini.

| Kalimat lama | Ralat |
| --- | --- |
| *"menerapkan satu business ke **seluruh kontrak berjenis reasuransi sama**"*; *"Gerbang pemilihan baris: `.REINSTYPEID == Param.REINSTYPEID`"* | ⭐ **terbalik**: prakondisi 3.1/3.2 `WhenTrue 3 = lewati` → sasaran = kontrak **lain di tahun treaty yang sama** yang jenisnya **BERBEDA** (pesan Pega `"Copied to all reins types."`) — R2 |
| *"penerapan massal karena itu **tidak atomik**"*; AC *"laporkan berapa berhasil dan berapa gagal"* | procedure tidak dipanggil (R4): satu transaksi — semua sasaran atau tidak sama sekali; AC sebagian-gagal gugur |
| *"tercatat di jejak audit"* | nol tabel jejak (K6): `USERID`/`TGLUPDATE` tiap baris + satu baris log server berisi cacah |
| — | tiap sasaran mendapat baris **baru** (`INSERT`, `TREATYBUSINESS_LIFE_SEQ`), tanpa penjaga dobel seperti Pega — OQ-MCRL-06 |
