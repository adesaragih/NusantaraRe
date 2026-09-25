# Panduan rantai alat — menutup dan menjalankan Fase 0

Disusun 25 September 2026, sesudah audit langsung atas mesin ini.
Untuk siapa pun yang menyiapkan mesin bagi proyek ini: IT, pengembang berikutnya, atau
sesi agent berikutnya.

> ⚠️ **Revisi 25 September 2026 (sore):** aplikasi dipindahkan ke `OUTPUT_HASIL_RNM\APP_RNM\`.
> Setiap `Set-Location 'D:\XML\RNM_BRD\OUTPUT_HASIL_RNM'` di bawah dibaca sebagai
> `...\OUTPUT_HASIL_RNM\APP_RNM`; `bin\api.exe` kini `APP_RNM\bin\api.exe`; frontend memakai
> TypeScript, sehingga `npm run typecheck` ikut dijalankan sebelum `npm run build`. Isi lama di
> bawah tidak diubah. Peta baca untuk pemula: `APP_RNM\README-BACA-DULU.md`.

---

## 0. Keadaan mesin ini — terukur, bukan dikira

| Alat | Keadaan | Letak |
| --- | --- | --- |
| `go` | ✅ **go1.26.8** windows/amd64 | `C:\Program Files\Go\bin\go.exe` |
| `node` | ✅ **v24.21.0** | `C:\Program Files\nodejs\node.exe` |
| `npm` | ✅ **11.19.0** | `C:\Program Files\nodejs\npm.cmd` |
| `git` | ✅ **2.55.0** | di PATH |
| `sh` | ✅ ada (dibawa Git) | `C:\Program Files\Git\bin\sh.exe` |
| `make` | ⛔ **tidak ada**, di mana pun | — |
| `docker` | ⛔ **tidak ada** | — |

Lingkungan lain yang sudah terukur:

| | |
| --- | --- |
| OS | Windows 11 Pro, build 26100, AMD64 |
| Sesi PowerShell | **bukan** sesi elevasi (`admin: False`) |
| `GOPROXY` | `https://proxy.golang.org,direct` — **terjangkau (200)** |
| `GOSUMDB` | `sum.golang.org` — **terjangkau (200)** |
| registry npm | `registry.npmjs.org` — **terjangkau (200)** |
| `GOMODCACHE` | `C:\Users\Administrator\go\pkg\mod` |
| Sumber winget | ⛔ **gagal diperbarui** — katalognya tak terjangkau dari jaringan ini |
| WSL | ⛔ belum terpasang; `HypervisorPresent = False` |
| Ruang bebas | C: 138 GB · D: 868 GB |

---

## 1. ⚠️ Jebakan yang paling mahal: PATH, bukan pemasangan

Laporan Fase 0 versi pertama menyatakan Go, Node, dan npm **tidak terpasang**. Itu keliru.
Ketiganya terpasang dan terdaftar di `Path` **tingkat mesin**; yang tidak memuatnya hanya
PATH shell yang dipakai sesi agent.

⛔ **Jangan menyimpulkan "alat tidak ada" dari `command -v` di satu shell.** Periksa
mesinnya:

```powershell
# 1. apakah binernya ada di tempat lazimnya
Test-Path 'C:\Program Files\Go\bin\go.exe'
Test-Path 'C:\Program Files\nodejs\node.exe'

# 2. apakah terdaftar di PATH tingkat mesin dan tingkat pengguna
([Environment]::GetEnvironmentVariable('Path','Machine') -split ';') -match 'Go\\bin|nodejs'
([Environment]::GetEnvironmentVariable('Path','User')    -split ';') -match 'Go\\bin|nodejs'

# 3. baru simpulkan
Get-Command go, node, npm -ErrorAction SilentlyContinue
```

**Bila alatnya ada tetapi shell tidak melihatnya**, cukup sisipkan untuk sesi itu:

```powershell
$env:Path = 'C:\Program Files\Go\bin;C:\Program Files\nodejs;' + $env:Path
```

Untuk Git Bash:

```bash
export PATH="/c/Program Files/Go/bin:/c/Program Files/nodejs:$PATH"
```

⭐ **Agar permanen bagi pengguna ini** (tanpa hak administrator):

```powershell
$lama = [Environment]::GetEnvironmentVariable('Path','User')
[Environment]::SetEnvironmentVariable(
  'Path', 'C:\Program Files\Go\bin;C:\Program Files\nodejs;' + $lama, 'User')
```

Tutup lalu buka kembali terminalnya. ⚠️ Jangan menulis ke `Path` tingkat **Machine** tanpa
alasan — itu menuntut elevasi dan menyentuh seluruh pengguna.

---

## 2. Menutup Fase 0 **tanpa** `make` — sudah terbukti jalan

Setiap target `Makefile` adalah satu perintah. Seluruh baris di bawah ini **sudah
dijalankan di mesin ini dan lulus**; keluaran yang diharapkan ikut dicantumkan supaya
dapat dibandingkan.

```powershell
$env:Path = 'C:\Program Files\Go\bin;C:\Program Files\nodejs;' + $env:Path
Set-Location 'D:\XML\RNM_BRD\OUTPUT_HASIL_RNM'
```

| Target | Perintah | Keluaran yang diharapkan |
| --- | --- | --- |
| *(sekali saja)* | `go mod tidy` | membuat `go.sum`, 6 baris |
| *(sekali saja)* | `cd frontend; npm install` | `added 90 packages`, membuat `package-lock.json` |
| `build` | `go build -o bin/api.exe ./cmd/api` | biner ±24,2 MB |
| `build` | `cd frontend; npm run build` | `✓ 30 modules transformed`, `dist/` ±142 kB |
| `test` | `go vet ./...` lalu `go test ./...` | `[no test files]` di 7 paket, exit 0 |
| `run-api` | `go run ./cmd/api` | log `http: mendengarkan di :8080` |
| `run-web` | `cd frontend; npm run dev` | `VITE ready`, melayani di port 5173 |
| `migrate` | `go run ./cmd/api -migrate` | menuntut `ORACLE_DSN`; melaporkan belum ada migrasi |

**Memeriksa `/healthz`:**

```powershell
$env:HTTP_ADDR = ':8099'
Start-Process .\bin\api.exe -NoNewWindow
Start-Sleep 2
(Invoke-WebRequest 'http://127.0.0.1:8099/healthz' -UseBasicParsing).Content
```

Jawaban yang benar tanpa Oracle:

```json
{"status":"sehat","database":"tidak dikonfigurasi"}
```

⚠️ **Server Vite mengikat `localhost`, bukan `127.0.0.1`.** Memanggil `http://127.0.0.1:5173/`
dapat gagal walau servernya sehat — pakai `http://localhost:5173/`.

---

## 3. Memasang `make`

Sumber winget **tidak terjangkau** dari jaringan ini, jadi `winget install` bukan jalan.
Tiga pilihan, dari yang paling ringan:

### 3a. Scoop — tanpa hak administrator ⭐ disarankan

Scoop memasang ke profil pengguna, jadi tidak menuntut elevasi, dan mengambil dari GitHub
yang **terjangkau** dari mesin ini.

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
scoop install make
```

⚠️ Baris kedua **menjalankan skrip dari internet**. Bila kebijakan keamanan kantor
melarangnya, pakai 3b.

### 3b. Unduh biner GNU Make langsung — paling konservatif

Ambil paket **ezwinports** `make-*-without-guile-w32-bin.zip`, bongkar, lalu taruh
`make.exe` di satu folder — misalnya `C:\Users\Administrator\bin` — dan tambahkan folder
itu ke `Path` tingkat pengguna dengan cara di §1. Tidak perlu administrator, tidak ada
pemasang, tidak ada layanan yang berjalan.

### 3c. Chocolatey — perlu administrator

```powershell
choco install make
```

Hanya bila kantor memang memakai Chocolatey; ia menuntut elevasi dan jangkauan ke
`chocolatey.org`.

### Sesudah `make` ada — satu syarat tambahan

`Makefile` proyek ini memakai sintaks shell POSIX, misalnya:

```make
db-up:
	@test -n "$$ORACLE_DEV_PASSWORD" || (echo "setel ORACLE_DEV_PASSWORD dulu" && exit 1)
```

GNU Make di Windows memakai `cmd.exe` bila tidak menemukan `sh`. **Jalankan `make` dari
Git Bash**, atau paksa shell-nya:

```powershell
make SHELL="C:/Program Files/Git/bin/sh.exe" build
```

Sesudah itu, jalankan keempat tanda selesai §2 brief dalam bentuk harfiahnya:
`make build` · `make test` · `make run-api` · `make run-web`.

---

## 4. Docker dan Oracle pengembangan

### ⭐ Periksa dulu: instance pengembangan mungkin sudah ada

Audit 25 September 2026 menemukan mesin ini **sudah memuat konfigurasi klien Oracle
lama**, dan di dalamnya daftar layanan yang sudah dikenal:

| Temuan | Nilai |
| --- | --- |
| Klien Oracle 12.2, **32-bit** | `C:\oracle12i` — ada `sqlplus.exe` |
| Klien Oracle 9i | `C:\oracle9i` — ada `sqlplus.exe` |
| `tnsnames.ora` | `C:\oracle12i\network\admin\` — **18 entri** |
| Alias berpenanda **DEV** | **8** |
| Alias berpenanda **PROD** | 8 |
| Alias tanpa penanda lingkungan | 2 |

⛔ **Isi berkas itu tidak dikutip ke mana pun** — ia memuat alamat layanan internal.
Yang dicatat hanya jumlahnya.

⭐ **Artinya: delapan layanan pengembangan sudah terdaftar di mesin ini.** Sebelum
memasang Docker sama sekali, tanyakan ke DBA apakah salah satunya dapat dipakai sebagai
instance pengembangan proyek ini, dan minta kredensial beserta hak `CREATE TABLE` pada
skema ujinya. Uji sambungannya dengan perkakas yang sudah ada:

```
C:\oracle12i\bin\sqlplus.exe <pengguna>@<alias DEV>
```

⛔ **Jangan** memakai alias berpenanda PROD. Brief §7: menjalankan migrasi pada instance
selain kontainer lokal menuntut persetujuan manusia.

### Apakah Docker benar-benar perlu?

⛔ **Tidak, bila DBA menyediakan instance pengembangan.** Brief §2 menyatakannya terang:
bila instance tersedia, `ORACLE_DSN` menunjuk ke sana dan `db-up` tidak dipakai. Docker
hanya melayani `make db-up`, yaitu jalan darurat ketika tidak ada instance.

### Oracle Instant Client 23.26 — terpasang, dan kapan ia terpakai

`C:\oracle\instantclient_23_0` (varian **basic**, 36 berkas, 387 MB, `oci.dll` 23.26.1.0),
terdaftar di `Path` tingkat pengguna.

⚠️ **Backend proyek ini tidak memerlukannya.** Driver `sijms/go-ora/v2` murni Go, dan
brief §2 memilihnya justru supaya Instant Client tidak diperlukan. Ia berguna untuk
perkakas DBA 64-bit dan bila kelak driver diganti ke yang berbasis OCI.

**Varian `basic` dipilih, bukan `basiclite`**, sebab charset Oracle sasaran **belum
diketahui** — nol DDL di korpus, OQ-001 masih terbuka. `basiclite` hanya membawa
`oraociicus.dll` (US English, charset terbatas); `basic` membawa `oraociei.dll` 298 MB
dengan seluruh charset dan pesan multi-bahasa, plus `ojdbc11/17.jar` dan `ucp11/17.jar`.
Memilih lite berarti mengandaikan charset yang tidak dapat diverifikasi.

⚠️ **Jebakan PATH yang tersisa.** `Path` tingkat **mesin** memuat `C:\oracle12i\bin` dan
`C:\oracle9i\bin`, dan PATH mesin selalu dibaca **sebelum** PATH pengguna. Perkakas baris
perintah karenanya tetap memakai klien 12.2. Itu **sengaja dibiarkan**: mendahulukan
Instant Client dapat mematahkan aplikasi 32-bit yang bergantung pada `C:\oracle12i`.
Buktinya `genezi -v` melaporkan *"Client Shared Library 32-bit - 12.2.0.1.0"*, bukan
23.26. Bila kelak ada program **64-bit** yang memuat `oci.dll`, ia dapat gagal karena
menemukan `oci.dll` 32-bit lebih dulu; perbaikannya bukan mengubah PATH mesin, melainkan
menunjuk direktori DLL secara eksplisit dari program itu.

`sqlplus` **tidak** ikut dalam paket `basic` — ia paket terpisah. Untuk sekarang pakai
`sqlplus` milik `C:\oracle12i\bin`.

### Bila tetap dipasang

Prasyaratnya **belum terpenuhi** di mesin ini dan seluruhnya menuntut administrator:

| Prasyarat | Keadaan | Cara |
| --- | --- | --- |
| Virtualisasi CPU | `HypervisorPresent = False` | aktifkan di BIOS/UEFI |
| `VirtualMachinePlatform` | belum terbaca (perlu elevasi) | `dism /online /enable-feature /featurename:VirtualMachinePlatform /all` |
| WSL2 | ⛔ belum terpasang | `wsl --install` lalu **reboot** |
| Docker Desktop | ⛔ belum terpasang | pemasang resmi; winget di sini tak terjangkau |

Sesudah itu:

```powershell
$env:ORACLE_DEV_PASSWORD = '<sandi lokal, jangan pernah sandi produksi>'
make db-up      # menjalankan gvenzl/oracle-free, membuat skema POOLDATA
```

### Env var yang dibaca aplikasi

Salin `.env.example` menjadi `.env` lalu isi. `.env` **tidak pernah** masuk git.

| Env var | Arti |
| --- | --- |
| `ORACLE_DSN` | Kosong berarti jalan tanpa database — sah di Fase 0 |
| `ORACLE_SCHEMA` | **Wajib** bila `ORACLE_DSN` terisi (ADR-U-0033) |
| `HTTP_ADDR` | Bawaan `:8080` |
| `IS_PEGA_PROD` | Penanda lingkungan (ADR-U-0005) |
| `SERVICE_<NAMA>` | Alamat layanan luar |
| `VITE_API_BASE_URL` | Alamat backend bagi frontend; kosong berarti sama-asal |

⛔ **`ORACLE_DSN` tidak pernah menunjuk instance produksi.** Menjalankan migrasi pada
instance selain kontainer lokal menuntut persetujuan manusia (brief §7).

---

## 5. Daftar periksa akhir

```powershell
$env:Path = 'C:\Program Files\Go\bin;C:\Program Files\nodejs;' + $env:Path
Set-Location 'D:\XML\RNM_BRD\OUTPUT_HASIL_RNM'

go version                    # go1.26.8 atau lebih baru (go.mod menuntut >= 1.22)
node --version ; npm --version
make --version                # sesudah bab 3
go vet ./...                  # exit 0
go test ./...                 # [no test files] di 7 paket, exit 0
gofmt -l cmd internal pkg     # tidak mencetak apa pun
go build -o bin/api.exe ./cmd/api
cd frontend ; npm ci ; npm run build ; cd ..
git status --porcelain        # kosong
```

⚠️ `npm ci` menuntut `package-lock.json` — berkas itu sudah di-commit, jangan dihapus.

---

## 6. Untuk sesi agent berikutnya

1. **Periksa mesin, bukan shell**, sebelum melaporkan alat tidak ada — §1.
2. Rujukan ADR di kode, komentar, dan pesan commit **wajib berawalan seri**
   (`ADR-U-nnnn`, `ADR-D-<modul>-nnnn`, `ADR-F-nnnn`). Nomor telanjang ditolak.
3. Titik tetap `/code-review` diambil **sebelum** baris kode pertama: `git rev-parse HEAD`.
4. Tulis skrip Python dengan `py` dan `PYTHONIOENCODING=utf-8`; jalur korpus memuat spasi.
5. ⚠️ Skrip Python yang menulis berkas di Windows **mengubah LF menjadi CRLF** bila dibuka
   dalam mode teks. Pakai mode biner (`io.open(p, 'wb')`) agar akhiran baris tidak berubah
   — satu commit sesi ini terpaksa dipakai hanya untuk membereskannya.
6. Fase 1 dimulai dari **Claim Life**, `.scratch\claim-life\issues\`, nomor terkecil yang
   tidak terblokir. Gerbang tabel akar bersama **tidak** menahan Fase 1.
