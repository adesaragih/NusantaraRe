# -*- coding: utf-8 -*-
"""
tulis-kamus-dan-ddl.py — menulis KEDUA keluaran dari SATU definisi.

    masukan : alat/definisi-skema-treaty-masuk.json   <- dibuat buat-kamus-dan-ddl.py
    keluaran: 2-to-spec/KAMUS-KOLOM.md
              2-to-spec/ddl-usulan/*.sql

Keduanya membaca berkas definisi yang SAMA, sehingga keduanya TIDAK DAPAT BERBEDA.
Bila salah satu menyimpang dari SPEC-MODEL-DATA.md §10, ALATNYA yang salah.

Jalankan berurutan:
    PYTHONIOENCODING=utf-8 python alat/buat-kamus-dan-ddl.py
    PYTHONIOENCODING=utf-8 python alat/tulis-kamus-dan-ddl.py
"""
import io, os, json

AKAR   = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
KELUAR = os.path.join(AKAR, '2-to-spec')
DDL    = os.path.join(KELUAR, 'ddl-usulan')
os.makedirs(DDL, exist_ok=True)

D = json.load(io.open(os.path.join(AKAR, 'alat', 'definisi-skema-treaty-masuk.json'), encoding='utf-8'))
urut, kol, selisih = D['urut'], D['kolom'], D['selisih']

# dedupe: POTONGAN sempat masuk dua kali (judul §10.9 kosong, lalu isi dari §14.2)
seen = set(); u = []
for s_, e_ in urut:
    if e_ in seen: continue
    seen.add(e_); u.append((s_, e_))
urut = u

# Tabel tipe TIDAK lagi diketik ulang di sini. Sejak 24 September 2026 setiap
# kolom membawa tipe Oracle-nya sendiri di dalam berkas definisi ('oracle'),
# dihitung SEKALI oleh buat-kamus-dan-ddl.py. Dua salinan tabel tipe adalah dua
# tempat yang dapat berbeda -- dan salah satunya sudah pernah berbeda: berkas ini
# masih mencetak kalimat "enam kolom memakai CLOB" berbulan-bulan sesudah
# keputusan CLOB itu dicabut.
TIPE = dict((k, tuple(v)) for k, v in D['tipe'].items())
KHUSUS = D.get('khusus', {})
# hasil langkah 1: 2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md
TINGKAT_BELUM = {
    'BATAS_MAKSIMUM_KELOMPOK': 'T-1', 'BATAS_MAKSIMUM_NON_KELOMPOK': 'T-1',
    'BATAS_PILIHAN': 'T-1', 'NILAI_BATAS': 'T-2', 'LIMIT_AGREGAT': 'T-3', 'NILAI_EGNPI': 'T-4',
}
# Golongan C dipakai penulis untuk pernyataan keputusan; definisinya ada di
# buat-kamus-dan-ddl.py, dan dibaca dari berkas definisi.
PAKET_C = set(tuple(x) for x in D.get('paket_c', []))

SKEMA, AKUN = 'TREATY_MASUK', 'TREATY_MASUK_APP'
ENT = set(kol)


def oracle(e, c):
    """Tipe sudah dihitung di hulu dan ikut di dalam definisi. Bila sebuah kolom
    datang tanpa 'oracle', itu tanda berkas definisi berasal dari jalan lama --
    dan perkakas ini BERHENTI, bukan menebak."""
    t = c.get('oracle')
    if not t:
        raise SystemExit("BERHENTI -- kolom %s.%s tanpa tipe Oracle di berkas definisi. "
                         "Jalankan buat-kamus-dan-ddl.py lebih dulu." % (e, c['kolom']))
    return t


def nul(c):
    return 'N' if c['kosong'].lower().startswith('tidak') else 'Y'


def induk_dari(e):
    out = []
    for c in kol[e]:
        n = c['kolom']
        if n == 'ID_' + e:
            continue
        if n.startswith('ID_') and n[3:] in ENT:
            out.append((n, n[3:]))
    return out


# ══ KAMUS-KOLOM.md ═══════════════════════════════════════════════════════════
K = []
w = K.append
w('# KAMUS KOLOM — modul Treaty In (skema `%s`)' % SKEMA)
w('')
w('<!-- BERKAS INI DIBANGKITKAN. Jangan disunting dengan tangan. -->')
w('> **Dasar bukti:** `SPEC-MODEL-DATA.md` §10 — **diurai**, bukan diketik ulang.')
w('> Dibangkitkan `alat/buat-kamus-dan-ddl.py` + `alat/tulis-kamus-dan-ddl.py` dari definisi yang')
w('> **sama** dengan `ddl-usulan/`. **Keduanya tidak dapat berbeda.** Bila salah satu menyimpang')
w('> dari §10, **alatnya** yang salah.')
w('')
w('> ### BATAS BERKAS INI, dinyatakan di muka')
w('>')
w('> Ia dapat menyatakan sebuah kolom **ADA di §10**. Ia **tidak** dapat menyatakan sebuah entitas')
w('> **LENGKAP** — bila §10 menulis atribut dalam bentuk yang pengurainya tidak kenal, atribut itu')
w('> hilang tanpa suara. Karena itu §0 **membandingkan** jumlah terurai dengan jumlah yang §10 tulis')
w('> di judulnya sendiri, dan melaporkan setiap selisih. **"Nol selisih" bukan "lengkap"** — ia')
w('> hanya berarti judul dan isinya sepakat.')
w('')
w('**Cacah berkas ini:** **%d** entitas, **%d** kolom.' % (len(urut), sum(len(kol[e]) for _, e in urut)))
w('')
w('---')
w('')
w('## 0. Pemeriksaan diri — terurai versus yang DITULIS §10')
w('')
if selisih:
    w('| Entitas | Terurai dari \u00a710 | Judul \u00a710 menulis | Selisih | Kolom sisipan |')
    w('|---|---:|---:|---:|---:|')
    for e, a_, b_, sis in selisih:
        w('| `%s` | %d | **%d** | **%+d** | %d |' % (e, a_, b_, a_ - b_, sis))
    w('')
    w('> **Selisih ini SISI SPESIFIKASI, bukan sisi perkakas** \u2014 kolom yang perkakas ini')
    w('> sisipkan sendiri (mata uang golongan B dan C, pemekaran induk polimorfik) sudah')
    w('> **dikeluarkan** dari kolom "Terurai", dan jumlahnya dicetak di kolom terakhir supaya')
    w('> dapat diperiksa. Selisih yang tersisa berarti judul \u00a710 dan tabelnya benar-benar')
    w('> tidak sepakat.')
    w('>')
    w('> **TIDAK ditambal di sini.** Mengarang atribut yang hilang membuat kamus ini berbohong')
    w('> dengan cara yang paling sulit ditemukan: pembacanya tidak punya cara tahu mana yang')
    w('> berasal dari \u00a710 dan mana yang ditambahkan.')
else:
    w('**Tidak ada selisih** antara jumlah terurai dan jumlah yang ditulis di judul \u00a710 \u2014')
    w('seluruh 35 entitas. Tiga lubang cacah yang pernah berdiri sudah ditutup 24 September 2026:')
    w('`F-2` (\u00a710.2), `F-3` (\u00a710.21), dan `F-17` (\u00a710.19a, tabel kedua yang terhisap ke')
    w('entitas judulnya).')
    w('')
    w('> **"Nol selisih" tetap BUKAN "lengkap".** Ia berarti judul dan isinya sepakat. Semesta')
    w('> yang dipakai menyusun \u00a710 masih kurang **340 properti titik buta** (`L-8`, ditagih')
    w('> `M-4`), dan tidak ada pemeriksaan di berkas ini yang dapat melihatnya.')
w('')
w('---')
w('')
w('## 1. Kelompok tipe — **DIPUTUSKAN `KTV-A`, tanpa verifikasi**')
w('')
w('§10 pendahuluan menyatakan tipe konseptualnya **mewarisi kelompok ADR-0003** dari modul Claim')
w('Non Prop, dan bahwa **angka presisinya diserahkan ke gerbang sesi DDL**. Angka itu **sudah')
w('diputuskan 24 September 2026**, butir **`KTV-A`** — bukan lagi warisan yang menunggu.')
w('')
w('**Terverifikasi: tidak.** Dasarnya satu kalimat: terlalu lebar di Oracle **murah**, terlalu')
w('sempit **memotong data** dan potongannya baru ketahuan sesudah data masuk — ongkosnya tidak')
w('setangkup, jadi sisi murahnya yang diambil. **Syarat pembalikannya: sesi DDL boleh')
w('mempersempit, dan HANYA sebelum data dimuat.** Alasan tiap angka ada di')
w('`KEPUTUSAN-TANPA-VERIFIKASI.md` §7.')
w('')
w('| Kel. | Usulan tipe Oracle | Isi |')
w('|---|---|---|')
for k in ('U1', 'P1', 'P2', 'K1', 'L1', 'C1', 'T', 'D', 'E', 'R'):
    w('| `%s` | `%s` | %s |' % (k, TIPE[k][0], TIPE[k][1]))
w('')
w('**Tidak ada satu pun kolom `CLOB` maupun `BLOB`.** Kalimat yang berdiri di sini sebelumnya')
w('\u2014 *"enam kolom memakai `CLOB`"* \u2014 **sudah tidak benar sejak keputusan itu dicabut**, dan')
w('ia bertahan karena berkas ini dulu menyimpan salinan tabel tipenya sendiri. Salinan itu')
w('dihapus: tipe kini dihitung **sekali** di `alat/buat-kamus-dan-ddl.py` dan dibawa di dalam')
w('berkas definisi.')
w('')
w('**Dua penyimpangan dari tabel di atas, keduanya `KTV-A`:**')
w('')
w('| Yang menyimpang | Jadi | Kenapa |')
w('|---|---|---|')
w('| kolom bernama `ID_*` yang \u00a710 tulis `C1` | `NUMBER(19)` | kunci asing `R` sudah `NUMBER(19)`; **dua ujung satu relasi tidak boleh berbeda lebar** |')
w('| tujuh kolom teks bebas | `VARCHAR2(4000 CHAR)` | `KETERANGAN`, `CATATAN`, `PENGECUALIAN`, `KETENTUAN_KHUSUS`, `CATATAN_BORDEREAUX` pada `VERSI_KONTRAK`; `NILAI_SEBELUM` dan `NILAI_SESUDAH` pada `JEJAK_PERUBAHAN`. Lebarnya mengikuti lebar yang sistem lama pakai untuk teks bebasnya sendiri |')
w('')
w('> Alasan tiap angka di tabel kelompok tipe ada di `KEPUTUSAN-TANPA-VERIFIKASI.md` \u00a77,')
w('> butir **`KTV-A`** \u2014 beserta **syarat pembalikannya**: sesi DDL boleh mempersempit,')
w('> dan hanya **sebelum data dimuat**.')
w('')
w('---')
w('')
for s_, e in urut:
    w('## `%s`' % e)
    w('')
    w('*§%s · %d kolom*' % (s_, len(kol[e])))
    w('')
    w('| Kolom | Tipe | Null | Asal di sistem lama | Catatan |')
    w('|---|---|---|---|---|')
    for c in kol[e]:
        asal = c['asal'] or '—'
        if asal.strip('*') in ('baru',):
            asal = '**baru** — tidak ada di sistem lama'
        cat = (c['catatan'] or '').replace('|', r'\|')
        w('| `%s` | `%s` | %s | %s | %s |' % (c['kolom'], oracle(e, c), nul(c), asal, cat))
    w('')
io.open(os.path.join(KELUAR, 'KAMUS-KOLOM.md'), 'w', encoding='utf-8').write('\n'.join(K))

# ══ ddl-usulan/ ══════════════════════════════════════════════════════════════
KEPALA = '''-- =====================================================================
-- %s
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tidak ada instans Oracle yang terjangkau (L-3). Dinyatakan di sini, bukan
-- disembunyikan: setiap baris di bawah adalah bacaan, bukan hasil uji.
--
-- DIBANGKITKAN dari SPEC-MODEL-DATA.md §10, dari definisi yang SAMA dengan
-- 2-to-spec/KAMUS-KOLOM.md. Jangan disunting dengan tangan -- suntingan
-- tangan membuat keduanya dapat berbeda.
--
-- PRESISI -- DIPUTUSKAN 24 September 2026, butir KTV-A.
--   Angka presisi di bawah bukan lagi warisan yang menunggu diputuskan. Ia
--   DIPUTUSKAN, TANPA VERIFIKASI, dengan dasar dan syarat pembalikan tertulis
--   di KEPUTUSAN-TANPA-VERIFIKASI.md sec 7.
--   Kaidahnya: terlalu lebar di Oracle MURAH -- NUMBER dan VARCHAR2 disimpan
--   panjang-berubah; terlalu sempit MEMOTONG DATA, dan potongannya baru
--   ketahuan sesudah data masuk. Ongkosnya tidak setangkup, jadi sisi murahnya
--   yang diambil.
--   SYARAT PEMBALIKAN: sesi DDL boleh MEMPERSEMPIT, dan HANYA SEBELUM DATA
--   DIMUAT. Sesudah itu tidak -- baris yang tidak muat tidak punya tempat pergi.
--   Yang memberi angkanya: Uji AP (panjang teks sebenarnya) dan Uji AQ.
--
-- SATU BATAS YANG DINYATAKAN, bukan disembunyikan:
--   VARCHAR2(4000 CHAR) sah dideklarasikan, tetapi pada basis data AL32UTF8
--   batas BITA-nya tetap 4000. Teks 4.000 aksara yang memuat aksara berbita
--   ganda tetap dapat ditolak saat disisipkan, kecuali MAX_STRING_SIZE=EXTENDED.
--   Tidak ada instans yang dapat ditanyai (L-3), jadi ini dicatat, bukan diuji.
--
-- TIDAK ADA CREATE PROCEDURE, CREATE FUNCTION, maupun trigger pembawa aturan
-- bisnis di seluruh folder ini -- ADR-0056 (K-4).
-- TIDAK ADA DML.
-- =====================================================================

'''


def kedalaman(e, lihat=frozenset()):
    if e in lihat:
        return 0
    p = induk_dari(e)
    return 0 if not p else 1 + max(kedalaman(x[1], lihat | {e}) for x in p)


# Bersihkan keluaran lama SEBELUM menulis. Tanpa ini, berkas dari jalan sebelumnya
# tertinggal dengan nomor urut lama dan folder memuat DUA tabel yang sama bernomor beda.
# Yang dihapus DILAPORKAN -- penghapusan yang diam sama berbahayanya dengan penolakan yang diam.
_lama = sorted(f for f in os.listdir(DDL) if f.lower().endswith('.sql'))
for f in _lama:
    os.remove(os.path.join(DDL, f))

urut_file = sorted((e for _, e in urut), key=lambda e: (kedalaman(e), e))
manifest = []
for i, e in enumerate(urut_file, 1):
    nm = '%02d_%s.sql' % (i, e)
    L = [KEPALA % nm, 'CREATE TABLE %s.%s (' % (SKEMA, e)]
    lb = max(len(c['kolom']) for c in kol[e]) + 2
    body = ['    %-*s %-22s %s' % (lb, c['kolom'], oracle(e, c),
                                   'NOT NULL' if nul(c) == 'N' else '') for c in kol[e]]
    L.append(',\n'.join(x.rstrip() for x in body))
    L.append(');')
    L.append('')
    # §16: constraint berbentuk <peran>_<tabel>[_n] -- TIDAK mengeja kolom,
    # dan tidak mengeja tabel kedua. Batas 30 bita diperiksa di bawah.
    L.append('ALTER TABLE %s.%s ADD CONSTRAINT PK_%s PRIMARY KEY (%s);'
             % (SKEMA, e, e[:27], kol[e][0]['kolom']))
    for j, (n, par) in enumerate(induk_dari(e), 1):
        L.append('ALTER TABLE %s.%s ADD CONSTRAINT FK_%s_%d FOREIGN KEY (%s)'
                 % (SKEMA, e, e[:25], j, n))
        L.append('    REFERENCES %s.%s (ID_%s);' % (SKEMA, par, par))

    # induk POLIMORFIK -- KTV-B, 24 September 2026.
    # Dulu: satu kolom ID_INDUK_*, tanpa kunci asing, dengan pernyataan keputusan
    # "basis data TIDAK menolak baris yatim". Sekarang: dua kolom bernama, dua
    # kunci asing sungguhan (dipasang otomatis di atas oleh induk_dari), dan satu
    # CHECK yang menuntut TEPAT SATU terisi.
    ganda = [c['kolom'] for c in kol[e]
             if c['kolom'] in ('ID_BAGIAN', 'ID_DETAIL_PROPORSIONAL')]
    if len(ganda) == 2:
        L += ['',
              '-- ---------------------------------------------------------------------',
              '-- INDUK POLIMORFIK -- KTV-B: dua kolom bernama, TEPAT SATU terisi',
              '-- ---------------------------------------------------------------------',
              '--   Menggantikan kolom tunggal ID_INDUK_%s, yang tidak dapat' % e,
              '--   punya kunci asing karena sasarannya bergantung nilai kolom lain',
              '--   (INV-17). Kedua kolom di bawah MASING-MASING punya kunci asing,',
              '--   dipasang di atas -- baris yatim menjadi mustahil.',
              '--   Dasar dan syarat pembalikannya: KEPUTUSAN-TANPA-VERIFIKASI.md sec 7,',
              '--   butir KTV-B. Ditinjau ulang bila induknya bertambah menjadi TIGA.',
              'ALTER TABLE %s.%s ADD CONSTRAINT CK_%s_INDUK' % (SKEMA, e, e[:22]),
              '    CHECK ( (%s IS NOT NULL AND %s IS NULL)' % (ganda[0], ganda[1]),
              '         OR (%s IS NULL     AND %s IS NOT NULL) );' % (ganda[0], ganda[1]),
              '']

    # PAKET UANG -- satu atribut konseptual, LEBIH DARI SATU kolom fisik
    # sejak P-8 hanya GOLONGAN C yang masih tanpa mata uang; A dan B sudah berbentuk
    paket = [c['kolom'] for c in kol[e] if (e, c['kolom']) in PAKET_C]
    if paket:
        L += ['',
              '-- ---------------------------------------------------------------------',
              '-- PAKET UANG GOLONGAN C -- KOLOM MATA UANG DIPASANG, BOLEH KOSONG',
              '-- ---------------------------------------------------------------------',
              '--   APA        : %s' % ', '.join(paket),
              '--                memperoleh kolom mata uang BERNAMA di sebelahnya,',
              '--                nullable. Namanya menyebut paketnya (MATA_UANG_<paket>)',
              '--                dan BUKAN "KODE_MATA_UANG" polos, sebab baris ini memuat',
              '--                lebih dari satu paket uang atau sudah punya kolom mata',
              '--                uang lain -- satu nama untuk dua arti adalah cacat yang',
              '--                ditanam dengan tangan sendiri.',
              '--   KENAPA     : sistem lama TIDAK merekam mata uang untuk besaran ini',
              '--                sama sekali. Kolom nullable yang tak terpakai murah;',
              '--                menambahkannya SESUDAH data masuk mahal, dan tidak ada',
              '--                sumber untuk mengisi baris warisannya.',
              '--   AKIBAT     : kolomnya akan KOSONG pada seluruh baris hasil migrasi.',
              '--                Itu BUKAN kegagalan -- itu keadaan yang benar, dan',
              '--                pembacanya harus tahu bahwa kosong berarti "tidak',
              '--                pernah dicatat", bukan "belum diisi".',
              '--   DILIHAT DI : KEPUTUSAN-TANPA-VERIFIKASI.md sec 7, butir KTV-C',
              '--   DITAGIH    : jawaban T-6 dari teknik treaty. Ia MENYEMPITKAN, tidak',
              '--                lagi memblokir: bila T-6 menyatakan besaran ini selalu',
              '--                dalam mata uang kontrak, kolomnya DICABUT SEBELUM DATA',
              '--                DIMUAT dan bacaannya diserahkan ke KODE_MATA_UANG_KONTRAK.',
              '--']

    belum = [c['kolom'] for c in kol[e] if c['kolom'] in TINGKAT_BELUM]
    if belum:
        L += ['', '-- ---------------------------------------------------------------------',
              '-- CONSTRAINT YANG SENGAJA TIDAK DIPASANG -- ini PERNYATAAN KEPUTUSAN,',
              '-- bukan TODO. Jangan dipasang tanpa jawaban yang disebut di baris DITAGIH.',
              '-- ---------------------------------------------------------------------']
        for b in belum:
            L += ['-- %s' % b,
                  '--   APA        : tidak ada CHECK maupun NOT NULL yang memasangkan paket uang',
                  '--                ini dengan PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40).',
                  '--   KENAPA     : tingkat pencatatannya BELUM DITENTUKAN -- sapuan penulisnya',
                  '--                tidak memisahkan apa pun. Constraint yang dipasang sekarang',
                  '--                menuntut nilai yang migrasi tidak tahu cara mengisinya.',
                  '--   AKIBAT     : kolom ini menerima angka tanpa menyatakan angkanya untuk',
                  '--                SELURUH TREATY atau untuk BAGIAN NuRe. Pembaca laporan tidak',
                  '--                dapat mengetahuinya dari basis data.',
                  '--   DILIHAT DI : 2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md §3',
                  '--   DITAGIH    : ketika jawaban %s kembali dari teknik treaty.' % TINGKAT_BELUM[b],
                  '--']
    L += ['', 'GRANT SELECT, INSERT, UPDATE, DELETE ON %s.%s TO %s;' % (SKEMA, e, AKUN)]
    io.open(os.path.join(DDL, nm), 'w', encoding='utf-8').write('\n'.join(L) + '\n')
    manifest.append((nm, e, len(kol[e]), len(induk_dari(e)), len(belum)))

# ── 00_SKEMA_DAN_AKUN.sql ────────────────────────────────────────────────────
S = [KEPALA % '00_SKEMA_DAN_AKUN.sql',
     '-- Dijalankan PALING AWAL, sebelum seluruh berkas tabel.',
     '-- Nama skema: ADR-0028 dan CONTEXT.md §3.5 -- TREATY_MASUK, BUKAN TREATY_IN,',
     '-- karena POOLDATA.TREATY_IN dan POOLDATA.TREATY_IN_EDM SUDAH ADA sebagai tabel',
     '-- sistem lama dan skema baru berdampingan di instance yang sama.',
     '--',
     '-- Pengenal di berkas ini, dihitung (batas 30 bita, §16):',
     '--   TREATY_MASUK            12',
     '--   TREATY_MASUK_APP        16',
     '-- Di atas 30 bita: NOL.',
     '',
     '-- BAGIAN 0 -- yang harus diputuskan DBA sebelum berkas ini berarti',
     '--   tablespace data dan indeks, kuota, profil kata sandi, dan apakah akun',
     '--   aplikasi boleh membuat objek. Tidak satu pun diputuskan di sini.',
     '',
     'CREATE USER %s IDENTIFIED BY "<diisi DBA>";' % SKEMA,
     'CREATE USER %s IDENTIFIED BY "<diisi DBA>";' % AKUN,
     '',
     'GRANT CREATE SESSION TO %s;' % AKUN,
     '']
io.open(os.path.join(DDL, '00_SKEMA_DAN_AKUN.sql'), 'w', encoding='utf-8').write('\n'.join(S) + '\n')

# ── Z00_KUNCI_ALAMI.sql ──────────────────────────────────────────────────────
# Sumber tiap baris: pernyataan "Kunci alami dan lingkupnya" di §10 entitasnya,
# dan nomor invariannya di SPEC-INVARIAN.md. TIDAK ada yang disimpulkan dari bentuk.
KUNCI = [
    ('INV-04', 'VERSI_KONTRAK',       ['ID_KONTRAK', 'NOMOR_URUT_VERSI']),
    ('INV-05', 'LAYER',               ['ID_VERSI_KONTRAK', 'NOMOR_LAYER', 'BAGIAN_LAYER']),
    ('INV-06', 'DETAIL_PROPORSIONAL', ['ID_LAYER', 'ID_KELOMPOK_TREATY']),
    ('INV-07', 'MATA_UANG_KONTRAK',   ['ID_VERSI_KONTRAK', 'KODE_MATA_UANG']),
    ('INV-08', 'RETENSI_CEDANT',      ['ID_VERSI_KONTRAK', 'ID_KELOMPOK_TREATY', 'KODE_MATA_UANG']),
    ('INV-09', 'EGNPI',                ['ID_VERSI_KONTRAK', 'ID_KELOMPOK_TREATY', 'KODE_MATA_UANG']),
    ('INV-10', 'PERIODE_PELAPORAN',   ['ID_VERSI_KONTRAK', 'PERIODE']),
    ('INV-12', 'TERMIN',               ['ID_VERSI_KONTRAK', 'NOMOR_TERMIN', 'KODE_MATA_UANG']),
    ('INV-11', 'PERIODE_AKUMULASI',   ['ID_VERSI_KONTRAK', 'PERIODE']),
    ('INV-13', 'SKALA_KOASURANSI',    ['ID_VERSI_KONTRAK', 'PERSEN_LIMIT']),
    ('INV-14', 'BATAS_PER_BAHAYA',    ['ID_VERSI_KONTRAK', 'ID_BAHAYA']),
    # KTV-B: induk polimorfik menjadi dua kolom bernama, jadi kunci alaminya
    # menjadi DUA constraint -- satu per pelekatan. Oracle melewatkan baris yang
    # salah satu kolom kuncinya NULL, sehingga masing-masing hanya menjaga
    # pelekatannya sendiri. Bentuk lama MENCAMPUR dua ruang pengenal.
    ('INV-15', 'POTONGAN',            ['ID_BAGIAN', 'ID_JENIS_POTONGAN']),
    ('INV-15', 'POTONGAN',            ['ID_DETAIL_PROPORSIONAL', 'ID_JENIS_POTONGAN']),
    ('INV-16', 'PENYEBARAN',          ['ID_BAGIAN', 'ID_JENIS_REASURANSI']),
    ('INV-16', 'PENYEBARAN',          ['ID_DETAIL_PROPORSIONAL', 'ID_JENIS_REASURANSI']),
    ('INV-64', 'BAGIAN',              ['ID_LAYER']),
    ('INV-65', 'RINCIAN_PENYEBARAN',  ['ID_PENYEBARAN', 'ID_JENIS_REASURANSI']),
    ('INV-66', 'PORTOFOLIO',          ['ID_VERSI_KONTRAK', 'ARAH_PORTOFOLIO', 'JENIS_PORTOFOLIO']),
    ('INV-67', 'DOKUMEN_KONTRAK',     ['ID_VERSI_KONTRAK', 'ID_DOKUMEN']),
]
for a in ('MATA_UANG', 'JENIS_POTONGAN', 'JENIS_REASURANSI', 'BAHAYA', 'KELOMPOK_TREATY', 'KELAS_BISNIS'):
    KUNCI.append(('INV-68', a, ['KODE']))

Z = [KEPALA % 'Z00_KUNCI_ALAMI.sql',
     '-- Dijalankan PALING AKHIR -- seluruh tabel harus sudah ada.',
     '-- Setiap baris menyebut NOMOR INVARIANNYA. Constraint tanpa invarian tidak',
     '-- ditulis di sini: ia tidak punya tempat untuk gagal.',
     '']
tak_terkompilasi = []
_n_uq = {}
for inv, e, cols in KUNCI:
    ada = {c['kolom'] for c in kol.get(e, [])}
    hilang = [c for c in cols if c not in ada]
    if hilang:
        tak_terkompilasi.append((inv, e, hilang)); continue
    _n_uq[e] = _n_uq.get(e, 0) + 1
    # Satu entitas dapat punya lebih dari satu kunci alami sejak KTV-B; namanya
    # diberi nomor urut, tetap di bawah 30 bita (sec 16).
    _nm = ('UQ_%s' % e[:27]) if _n_uq[e] == 1 else ('UQ_%s_%d' % (e[:25], _n_uq[e]))
    Z.append('-- %s' % inv)
    Z.append('ALTER TABLE %s.%s ADD CONSTRAINT %s UNIQUE (%s);'
             % (SKEMA, e, _nm, ', '.join(cols)))
Z += ['',
      '-- ---------------------------------------------------------------------',
      '-- INV-08, INV-09, INV-12 -- DIBEBASKAN 24 September 2026',
      '-- ---------------------------------------------------------------------',
      '--   Ketiganya sempat TIDAK dapat dikompilasi: kunci alaminya menyebut',
      '--   "kode mata uang" dan kolomnya tidak ada, karena mata uang melebur ke',
      '--   dalam paket uang. P-8 golongan B memutuskan mata uang tinggal pada',
      '--   BARIS YANG SAMA -- dan itu memang bentuknya di sistem lama.',
      '--   Sejak itu ketiganya tertulis sebagai UNIQUE di atas.',
      '--',
      '--   Catatan yang TIDAK boleh hilang: menulis UNIQUE ini TANPA mata uang',
      '--   bukan jalan tengah -- ia mengubah artinya menjadi "dilarang dua baris',
      '--   bermata uang berbeda", dan itu MENOLAK DATA YANG SAH.',
      '--',
      '-- KUNCI ALAMI YANG SENGAJA BUKAN CONSTRAINT -- bukan kelalaian:',
      '--   KONTRAK              : cedant + SoB + periode + sifat proporsi.',
      '--                          ADR-0040 §2 MEMPERINGATKAN, tidak melarang. §10.1',
      '--                          menyatakannya, dan ketiadaan nomor INV-nya disengaja.',
      '--   NILAI_PENYEBARAN     : kode mata uang di dalam rincian -- BELUM BERNOMOR;',
      '--                          kedudukan entitasnya sendiri ditangguhkan Uji AD',
      '--                          (§10.23c butir 3). SPEC-INVARIAN.md §7.3.',
      '--   PEMULIHAN_LIMIT      : belum bernomor; entitasnya baru diterima §10.23c.',
      '--   CATATAN_PERSETUJUAN  : TIDAK ADA, dan itu keputusan -- peristiwa yang sama',
      '--   JEJAK_PERUBAHAN        dapat terjadi dua kali pada versi yang sama.',
      '--   PERISTIWA_KONTRAK      §10.20, §10.21, §10.21a. SPEC-INVARIAN.md §7.4.',
      '']
io.open(os.path.join(DDL, 'Z00_KUNCI_ALAMI.sql'), 'w', encoding='utf-8').write('\n'.join(Z) + '\n')

print('DIBERSIHKAN lebih dulu: %d berkas .sql lama dihapus' % len(_lama))
print()
print('DITULIS')
print('   2-to-spec/KAMUS-KOLOM.md       %d entitas, %d kolom'
      % (len(urut), sum(len(kol[e]) for _, e in urut)))
print('   2-to-spec/ddl-usulan/          %d berkas tabel' % len(manifest))
print('   kunci asing                    %d' % sum(m[3] for m in manifest))
print('   constraint sengaja TIDAK dipasang, bertuliskan pernyataan keputusan : %d'
      % sum(m[4] for m in manifest))
print()
print('TIDAK DIBUATKAN TABEL, dan itu disengaja:')
print('   NILAI_SELISIH                  menunggu grilling Adjustment (§11.3)')
print('   BESARAN_DAPAT_DISESUAIKAN      menunggu grilling Adjustment (§11.3)')
print('   RETRO_KELUAR, PENCAPAIAN       GEL-2 / GEL-3, di luar gelombang ini')
