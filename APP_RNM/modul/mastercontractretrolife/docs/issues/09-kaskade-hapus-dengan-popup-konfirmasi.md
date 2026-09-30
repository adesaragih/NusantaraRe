# 09: Kaskade hapus induk → anak, dengan popup konfirmasi sebelum apa pun terhapus

**Status:** ready-for-agent

**Blocked by:** 02 (kontrak), 05 (reinsurer), 06 (security reinsurer), 07 (business) — keempat
entitas harus ada agar kaskade dapat diuji utuh

## Hasil & nilai pengguna

Sebagai **admin master**, saya menghapus sebuah induk dan **seluruh anaknya ikut terhapus** —
sehingga tidak ada baris yatim yang tertinggal — tetapi **hanya setelah** sistem memberi tahu apa dan
berapa yang akan ikut hilang, dan saya menyetujuinya. *(User story 27–29 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Penghitungan anak per induk; penghapusan |
| `internal/services` | **Aturan kaskade + pratinjau dampak**; larangan hapus tahun treaty |
| `internal/handlers` | Endpoint pratinjau hapus; endpoint eksekusi hapus |
| `frontend/` | **Dialog konfirmasi Ya/Batal** yang menyebut apa & berapa |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Sasaran |
| --- | --- | --- | --- |
| `DeleteTreatyLimit_SQL` | `ASM-FW-GISFW-INT-RETROCESSIONLIFE` / `ASM!DELETETREATYLIMIT_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/DeleteTreatyLimit_SQL.xml` | `treatycontract_life` |
| `DeleteSecurityReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!DELETESECURITYREINSURER_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/DeleteSecurityReinsurer_SQL.xml` | `TREATYREINSURER_LIFE` |
| `DeleteSecurityReinsurerLife_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!DELETESECURITYREINSURERLIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/DeleteSecurityReinsurerLife_SQL.xml` | `TREATYSECURITYREINSURER_LIFE` |
| `DeleteRowBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!DELETEROWBUSINESSLIST` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/DeleteRowBusinessList.xml` | `treatybusiness_life` |
| `DeleteTreatyLimit_Act`, `DeleteSecurityLife_Act`, `DeleteSecurityReinsurerLife_Act`, `DeleteRowBusiness` | `ASM-FW-GISFW-…` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/` | orkestrator hapus |

`[terverifikasi]` **Keempatnya `DELETE … WHERE ID = …` datar, tanpa kaskade** — menghapus induk
meninggalkan anak yatim. `[terverifikasi]` **Tidak ada penghapus untuk `TREATYYEAR_LIFE`**.

`[data DBA]` **Keempat FK sudah terpasang** dengan mode **`ON DELETE CASCADE`**:

| FK | Kolom anak | → induk |
| --- | --- | --- |
| 1 | `TREATYCONTRACT_LIFE.IDTREATYYEAR` | `TREATYYEAR_LIFE.ID` |
| 2 | `TREATYREINSURER_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |
| 3 | `TREATYSECURITYREINSURER_LIFE.TREATYREINSURERID` | `TREATYREINSURER_LIFE.ID` |
| 4 | `TREATYBUSINESS_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |

## ADR terkait

**ADR-0007** (jejak audit penghapusan), **ADR-0001** (master ini dirujuk Claim Life, Komite Claim
Life, dan Master Product Name Life — penghapusan berdampak lintas konteks).

## Acceptance criteria

- [ ] ⚠️ Menghapus **kontrak** memunculkan **konfirmasi lebih dulu** yang menyebut **apa dan berapa**
      yang akan ikut terhapus — reinsurer, security reinsurer, dan business di bawahnya. *(AC 32
      spec; `[keputusan work owner]` — penyimpangan sadar 3)*
- [ ] ⚠️ Menghapus **reinsurer** memunculkan konfirmasi yang menyebut berapa **security reinsurer**
      akan ikut terhapus. *(AC 33 spec)*
- [ ] ⚠️ Menekan **Ya** menghapus induk **beserta seluruh sub-pohonnya**; **tidak ada baris yatim**
      yang tertinggal — dibuktikan dengan menghitung baris anak sesudahnya. *(AC 34 spec)*
- [ ] ⚠️ Menekan **Batal** membuat **tidak ada satu pun** baris terhapus — induk maupun anak.
      Dibuktikan dengan menghitung baris sebelum dan sesudah. *(AC 35 spec)*
- [ ] Menghapus baris **tanpa anak** tetap memerlukan konfirmasi, dengan pesan yang menyatakan
      **tidak ada anak** yang terpengaruh. *(AC 36 spec)*
- [ ] Menghapus **security reinsurer** atau **business** (baris daun) menghapus **hanya baris itu**.
      *(AC 37 spec)*
- [ ] ⚠️ **Tahun treaty tetap tidak dapat dihapus** — kaskade **tidak** membuka jalur hapus untuknya.
      *(AC 2 spec; `[fakta bisnis — work owner]`)*
- [ ] Pratinjau dampak dihitung **sesaat sebelum** dialog tampil, bukan dari data lama yang mungkin
      sudah berubah.
- [ ] Penghapusan tercatat di **jejak audit**: siapa, kapan, dan berapa baris ikut terhapus.
- [ ] Menghapus baris yang **sudah tidak ada** memberi pesan yang jelas, bukan galat mentah.

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Basis data adalah lapis kedua, bukan pengganti.** `[data DBA]` FK `ON DELETE CASCADE` akan ikut
menghapus anak — tetapi **FK tidak dapat menjelaskan apa pun kepada pengguna**. Pesan yang menyebut
apa dan berapa anak **tetap tanggung jawab Go**, dan pratinjau dihitung sebelum penghapusan dimulai.

⚠️ **Perubahan arah keputusan.** Versi awal grilling mencatat "tolak hapus bila punya anak".
`[keputusan work owner]` **Direvisi menjadi kaskade + konfirmasi** — hapus induk memang dimaksudkan
menghapus seluruh sub-pohonnya. Tiket ini mengikuti keputusan **final**.

⚠️ **Kejanggalan penamaan** `[terverifikasi]`: `DeleteTreatyLimit_SQL` ber-class
`ASM-FW-GISFW-INT-RETROCESSIONLIFE` tetapi menghapus `treatycontract_life` — class tidak sejalan
dengan tabel sasarannya. Jangan tiru penamaannya.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — kaskade FK hanya berperilaku benar pada
basis data sungguhan; memalsukannya berarti tidak menguji apa pun yang penting.

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
| *"`[data DBA]` **Keempat FK sudah terpasang** dengan mode **`ON DELETE CASCADE`**"*; *"FK `ON DELETE CASCADE` akan ikut menghapus anak"* | ⛔ DEV **nol FK** (K1). Kaskade **di Go**, satu transaksi, anak lebih dulu: (1) security, (2) reinsurer dan business, (3) kontrak — **sesudah** popup konfirmasi (K2) |
| TAMBAHAN-TIKET: *"Test kaskade **tidak boleh** mengandaikan aplikasi yang menghapus anak"* | justru aplikasi yang menghapus anak; uji membuktikan urutan dan hasil akhirnya |
| *"Penghapusan tercatat di **jejak audit**"* | nol tabel jejak (K6): satu baris log server berisi cacah baris terhapus per tabel, tanpa nama orang |
| — | pesan sukses VERBATIM `"Data Berhasil di Hapus"` (kontrak/reinsurer/security) dan `"Data Dengan ID <id> Berhasil di Hapus"` (business); konfirmasi yang jumlahnya tidak lagi cocok dengan data ditolak 409 tanpa satu baris pun terhapus |
