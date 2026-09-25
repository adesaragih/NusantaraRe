---
status: selesai
---

# 35: Satu pintu tulis ditegakkan hak akses, bukan kesepakatan

> **Gerbang K1 gugur — 19 September 2026. Tiket ini tidak lagi menunggu apa pun dari bagian K.**
> K1 ditutup sebagai **K1a**: ADR-0016 tetap berlaku tanpa revisi, dan modul ini **tidak mengambil data HRD sama sekali**. `V_MST_USER_TEKNIS` tidak dimiliki dan tidak dimigrasi (ADR-0026); database link yang berjalan produksi (`hrdasm.v_hrd_mst@asmd.sinarmas.co.id`, D20) tetap hidup di sistem lama dan bukan urusan tiket ini.
>
> Yang tersisa sebagai catatan, bukan penahan: **REQ-037**. Ia menentukan apakah klaim ADR-0017 berlaku sungguhan di instance tujuan — enam jalur tulis yang melewati hak objek. Rumusan klaim ADR-0017 sudah diperbaiki di badannya (*"tidak ada tulisan dari luar pintu yang tidak meninggalkan jejak"*), dan `Z01_PENGAWASAN_TULIS.sql` adalah bentuknya. Tiket ini menegakkan hak objek; ia tidak menjanjikan pencegahan yang tidak dapat dijanjikan skema.



> **SELESAI 19 September 2026 — dan cacat yang sama dengan tiket `01` terulang di sini.**
>
> **Keempat puluh empat `REVOKE ALL` akan menggugurkan pemasangan di instance bersih.** Di instance yang belum pernah dipakai, tidak satu pun hak itu pernah diberikan — dan Oracle **menolak pencabutan hak yang tidak ada**: `ORA-01927`. Baris pertama gagal, dan pemasangan berhenti di sana.
>
> Ini **kelas cacat yang sama** dengan sembilan `REVOKE` di `00_SKEMA_DAN_AKUN.sql` yang sudah diperbaiki lebih dulu hari ini, dan ia terulang karena keduanya ditulis dengan anggapan yang sama: bahwa mencabut sesuatu yang tidak ada tidak berakibat apa-apa. **Di Oracle, itu keliru — dan sekali ketahuan, seharusnya disapu ke seluruh berkas, bukan diperbaiki di satu tempat.**
>
> Diperbaiki dengan bentuk yang sama: blok PL/SQL yang mencabut bila ada, melewati `ORA-01927`, dan **merambat untuk galat apa pun selain itu** — sehingga kegagalan yang sesungguhnya tetap menghentikan pemasangan.
>
> **Pengecualian arsip kini disebut di sini juga.** Hibahnya tetap di `22_ARSIP_MUATAN_KELUAR.sql` (`SELECT, INSERT` saja, tanpa `UPDATE`, tanpa `DELETE`); yang ditulis di `V00` hanya **faktanya**, supaya pengecualian itu tidak hanya hidup di satu berkas tabel yang tidak dibaca orang saat meninjau hak akses.
>
> **Tiga pemeriksaan ditulis di kaki berkas**, dan yang ketiga menyatakan batasnya: keenam jalur REQ-037 **tidak dapat diperiksa dari sini**, dan nol baris pada pemeriksaan pertama **tidak berarti** tidak ada yang dapat menulis.
>
> **REQ-021 dan REQ-037 tetap menunggui, tidak menahan**: skema baru tidak mengulangi hak yang dipegang `POOLDATA` atas tabel work Pega, apa pun jawabannya.

*Asal: `T-20` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada akun selain akun aplikasi yang dapat menulis ke tabel kanonik; hilir membaca hanya lewat view.

`V00_HAK_AKSES.sql`: **44 `REVOKE ALL`** atas tabel kanonik untuk `POOLDATA` dan `KLAIMNP_HILIR`; `GRANT SELECT` atas view. **Termasuk `ARSIP_MUATAN_KELUAR`**: akun aplikasi memegang `INSERT` dan `SELECT` saja — **tanpa `UPDATE` dan tanpa `DELETE`**, karena arsip itu tulis-sekali.

**Tidak termasuk:** Perintah yang dijalankan tangan — `GRANT` dan `REVOKE` adalah **objek DDL**. Hak yang tidak tertulis di DDL tidak dapat diaudit.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`22`~~ *(selesai)* — Koreksi bernilai tercatat menggantikan tambalan di dalam kode
- ~~`13`~~ *(selesai)* — Bentuk lama boleh masuk utuh, di satu tempat saja
- ~~`15`~~ *(selesai)* — Jembatan ke sistem lama berdiri sebagai tabel terpisah
- ~~`29`~~ *(selesai)* — Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam
- ~~`23`~~ *(selesai)* — Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama


**Dasar:** DECIDED(ADR-0017, ADR-0028). EVIDENCED: `GRANT ALTER, DELETE, INDEX, INSERT, REFERENCES, SELECT, UPDATE ... TO POOLDATA` pada tabel work Pega — hak yang ada akan dipakai cepat atau lambat.

- [x] `POOLDATA` gagal `INSERT`, `UPDATE`, dan `DELETE` ke setiap tabel kanonik — 44 pencabutan, kini **tahan instance bersih**.
- [x] `POOLDATA` berhasil `SELECT` lewat view kompatibilitas — hibahnya di `V01`…`V09`.
- [x] Akun hilir gagal `SELECT` langsung dari tabel kanonik mana pun.
- [x] Akun aplikasi memegang `SELECT` dan `INSERT` saja atas `ARSIP_MUATAN_KELUAR` — tulis-sekali, dan faktanya terbaca di `V00` maupun di berkas tabelnya.

**Ketidakpastian:** **REQ-021** — apakah `POOLDATA` benar-benar menulis ke tabel Pega hari ini belum diketahui. Bila ya, migrasinya memerlukan pokok tersendiri untuk manajemen; itu **tidak menahan tiket ini**, karena skema baru tidak mengulangi hak itu apa pun jawabannya.
