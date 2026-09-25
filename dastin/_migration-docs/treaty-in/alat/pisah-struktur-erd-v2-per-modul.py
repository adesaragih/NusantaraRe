# -*- coding: utf-8 -*-
r"""
pisah-struktur-erd-v2-per-modul.py -- memecah STRUKTUR-DATA-ERD-V2.md menjadi DUA,
satu per modul, masing-masing disimpan di folder modulnya sendiri.

    masukan : alat/erd-v2.json                 <- SUMBER YANG SAMA dengan
                                                  bangkitkan-dari-erd-v2.py
    keluaran: ../treaty-in/STRUKTUR-DATA-TREATY-IN.md
              ../treaty-in-adjustment/STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md

KENAPA DARI SUMBERNYA, BUKAN DARI BERKAS GABUNGANNYA
    STRUKTUR-DATA-ERD-V2.md berbunyi di kepalanya sendiri: "tidak ditulis tangan,
    dan tidak boleh disunting -- sunting sumbernya, lalu jalankan ulang alatnya."
    Memotong keluarannya menjadi dua berarti melahirkan dua berkas yang TIDAK DAPAT
    dibangkitkan ulang, dan yang akan basi pada perubahan sumber berikutnya.
    Perkakas ini membacanya dari erd-v2.json, persis seperti yang gabungannya.

SUMBU PEMISAHNYA, dan ia mekanis
    Nama lembar yang memuat "EDM" -> modul Treaty In Adjustment.
    Selain itu -> modul Treaty In.
    Itu SATU-SATUNYA aturan, dan ia dicetak di dalam kedua keluarannya.

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa tiap kotak, tiap tabel unik, dan tiap relasi di sumbernya MENDARAT di
    sedikitnya satu keluaran -- dan bahwa cacah keduanya menutup ke cacah sumber.

APA YANG TIDAK
    Ia TIDAK memisahkan MODEL; ia memisahkan GAMBAR. Ke-34 tabel yang dipakai
    kedua modul muncul di KEDUA berkas, dan itu BUKAN penggandaan melainkan
    kenyataan sistem lama: satu tabel yang sama dipakai jalur kontrak maupun
    jalur addendum. Penandanya ikut dicetak per baris.

    Dan ia TIDAK mengubah satu huruf pun isi: nama tabel, kardinalitas, kunci,
    jalur Pega, padanan model baru, dan tingkat bukti seluruhnya apa adanya.

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/pisah-struktur-erd-v2-per-modul.py
"""
import io, os, json
from collections import Counter, OrderedDict, defaultdict

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
SUMBER = os.path.join(AKAR, 'alat', 'erd-v2.json')
KELUAR = {
    'Treaty In': os.path.join(AKAR, 'STRUKTUR-DATA-TREATY-IN.md'),
    'Treaty In Adjustment': os.path.join(AKAR, '..', 'treaty-in-adjustment',
                                         'STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md'),
}
TANGGAL = '25 September 2026'

tolak = Counter()
d = json.load(io.open(SUMBER, encoding='utf-8'))
lembar = d['lembar']

# ── 1. bagi lembar menurut sumbu EDM ────────────────────────────────────────
MODUL = OrderedDict([('Treaty In', []), ('Treaty In Adjustment', [])])
for nama in lembar:
    MODUL['Treaty In Adjustment' if 'EDM' in nama else 'Treaty In'].append(nama)

# ── 2. relasi: dicocokkan lewat kolom "Muncul di lembar" ────────────────────
def pendek(nama):
    return nama.replace('Treaty In ', '')


relasi_modul = defaultdict(list)
relasi_takterklasifikasi = []
for r in d['relasi']:
    if not r.get('Anak') or not r.get('Induk'):
        tolak['baris di Daftar Relasi yang bukan relasi (catatan kaki)'] += 1
        continue
    ml = [x.strip() for x in (r.get('Muncul di lembar') or '').split(',')]
    kena = False
    for mod, lbs in MODUL.items():
        if any(pendek(l) in ml for l in lbs):
            relasi_modul[mod].append(r)
            kena = True
    if not kena:
        relasi_takterklasifikasi.append(r)
        tolak['relasi yang tidak muncul di lembar mana pun -- dibawa ke KEDUA berkas'] += 1

# ── 3. pohon, disalin dari pembangkit gabungannya ──────────────────────────
def pohon(kotak):
    akar, tumpuk = [], []
    for k in kotak:
        n = dict(k, anak=[])
        while tumpuk and tumpuk[-1]['dalam'] >= n['dalam']:
            tumpuk.pop()
        if not tumpuk:
            if n['dalam'] != 0:
                tolak['kotak berkedalaman > 0 tanpa induk di atasnya'] += 1
            akar.append(n)
        else:
            tumpuk[-1]['anak'].append(n)
        tumpuk.append(n)
    return akar


def pohon_md(simpul, out, dalam=0):
    for n in simpul:
        f = n['fk'] or {}
        kunci = ('PK `%s`' % n['pk']) if n['pk'] else \
                ('FK `%s` → `%s`.`%s` · %s' % (f.get('kolom'), f.get('induk'),
                                               f.get('kolom_induk'), f.get('on_delete')))
        out.append('%s- **`%s`**%s — %s' % ('  ' * dalam, n['nama'],
                                            ' *%s*' % n['kard'] if n['kard'] else '', kunci))
        pohon_md(n['anak'], out, dalam + 1)


# ── 4. himpunan tabel per modul ────────────────────────────────────────────
tabel_modul = {}
kotak_modul = {}
for mod, lbs in MODUL.items():
    t = OrderedDict()
    n = 0
    for nama in lbs:
        for k in lembar[nama]['kotak']:
            n += 1
            if k['nama'] not in t:
                t[k['nama']] = k
    tabel_modul[mod] = t
    kotak_modul[mod] = n

A, B = list(MODUL)
bersama = [x for x in tabel_modul[A] if x in tabel_modul[B]]
hanya = {A: [x for x in tabel_modul[A] if x not in tabel_modul[B]],
         B: [x for x in tabel_modul[B] if x not in tabel_modul[A]]}
semua_tabel = set(tabel_modul[A]) | set(tabel_modul[B])

# tabel yang HANYA muncul di Daftar Relasi, tidak pernah sebagai kotak
# Kurung di bawah BUKAN kerapian. Tanpa kurungnya, Python menghitung
# {Anak} - semua_tabel LEBIH DULU lalu menggabungkannya dengan {Induk} -- dan
# seluruh induk lolos tanpa disaring. Ia sempat mencetak 15, bukan 4.
tabel_relasi_saja = sorted(({r['Induk'] for r in d['relasi'] if r.get('Induk')} |
                            {r['Anak'] for r in d['relasi'] if r.get('Anak')})
                           - semua_tabel)
for x in tabel_relasi_saja:
    tolak['tabel yang hanya muncul di Daftar Relasi, tidak pernah sebagai kotak'] += 1

# ── 5. tulis ───────────────────────────────────────────────────────────────
LAIN = {A: B, B: A}
hasil = []
for mod, lbs in MODUL.items():
    tabel = tabel_modul[mod]
    di_lembar = defaultdict(list)
    for nama in lbs:
        for k in lembar[nama]['kotak']:
            di_lembar[k['nama']].append(nama)

    m = []
    w = m.append
    w('# Struktur data sistem lama — **%s**' % mod)
    w('')
    w('> ### %s' % d['banner'][0])
    w('>')
    w('> %s' % d['banner'][1])
    w('')
    w('**Dibangkitkan:** [`alat/pisah-struktur-erd-v2-per-modul.py`](%salat/pisah-struktur-erd-v2-per-modul.py) '
      '· **Sumber tunggal:** `alat/erd-v2.json` — **sumber yang sama** dengan berkas gabungannya, '
      '`treaty-in/4-erd-dan-tabel-datar/STRUKTUR-DATA-ERD-V2.md`.'
      % ('' if mod == A else '../treaty-in/'))
    w('')
    w('> **Berkas ini tidak ditulis tangan, dan tidak boleh disunting.** Ia bagian **%s** dari '
      'potret `%s`. Sunting sumbernya, lalu jalankan ulang alatnya.' % (mod, d['sumber']))
    w('')
    w('---')
    w('')
    w('## 0. Bagaimana pemisahannya dilakukan, dan apa yang TIDAK dipisah')
    w('')
    w('**Sumbu pemisahnya satu, dan ia mekanis:** nama lembar yang memuat **`EDM`** masuk modul')
    w('**Treaty In Adjustment**; selain itu masuk modul **Treaty In**. Tidak ada aturan kedua.')
    w('')
    w('| Modul | Lembar |')
    w('|---|---|')
    for mm, ll in MODUL.items():
        w('| %s%s | %s |' % ('**' if mm == mod else '', mm + ('**' if mm == mod else ''),
                             ' · '.join('*%s*' % x for x in ll)))
    w('')
    w('> ### YANG DIPISAH ADALAH GAMBARNYA, BUKAN MODELNYA')
    w('>')
    w('> **%d dari %d tabel dipakai KEDUA modul**, dan karena itu muncul di **kedua** berkas.'
      % (len(bersama), len(semua_tabel)))
    w('> Itu **bukan penggandaan** melainkan kenyataan sistem lama: satu tabel yang sama dipakai')
    w('> jalur kontrak **dan** jalur addendum. Kolom **Bersama modul lain** menandainya per baris.')
    w('>')
    w('> | | |')
    w('> |---|---:|')
    w('> | tabel unik di **%s** | **%d** |' % (mod, len(tabel)))
    w('> | — di antaranya **juga** di %s | **%d** |' % (LAIN[mod], len(bersama)))
    w('> | — **hanya** di %s | **%d**%s |'
      % (mod, len(hanya[mod]),
         (' — ' + ', '.join('`%s`' % x for x in hanya[mod])) if hanya[mod] else ''))
    w('> | tabel unik di kedua modul digabung | **%d** |' % len(semua_tabel))
    w('')
    w('---')
    w('')
    w('## 1. Hitungan')
    w('')
    w('| | Jumlah |')
    w('|---|---:|')
    w('| Tabel unik di modul ini | **%d** |' % len(tabel))
    w('| Kotak tergambar (satu tabel dapat muncul di beberapa lembar) | **%d** |' % kotak_modul[mod])
    w('| Relasi yang menyentuh lembar modul ini | **%d** |' % len(relasi_modul[mod]))
    w('| Relasi **tidak terklasifikasi** — dibawa ke kedua berkas | **%d** |'
      % len(relasi_takterklasifikasi))
    for nama in lbs:
        w('| Lembar *%s* | %d kotak |' % (nama, len(lembar[nama]['kotak'])))
    w('')
    w('## 2. Entitas — satu baris per tabel unik')
    w('')
    w('| Tabel | Kard. | Kunci utama | Kunci asing | ON DELETE | Induk | Muncul di lembar | Bersama %s |'
      % LAIN[mod])
    w('|---|---|---|---|---|---|---|---|')
    for nm, k in tabel.items():
        f = k['fk'] or {}
        w('| `%s` | %s | %s | %s | %s | %s | %s | %s |'
          % (nm, k['kard'] or '—',
             '`%s`' % k['pk'] if k['pk'] else '*tidak dinyatakan*',
             '`%s`' % f['kolom'] if f else '—',
             f.get('on_delete', '—') or '—',
             '`%s`' % f['induk'] if f else '**akar**',
             ', '.join(pendek(x) for x in di_lembar[nm]),
             '**ya**' if nm in bersama else '**TIDAK — khas modul ini**'))
    w('')
    w('> **Kolom *Kunci utama* berbunyi *tidak dinyatakan* untuk %d dari %d tabel, dan itu salinan '
      'setia.** Berkas sumbernya mencetak **PK hanya pada kotak akar**; pada kotak anak baris '
      'keduanya dipakai untuk **FK**. Maka kekosongan ini berarti **sumbernya diam**, bukan '
      '**tabelnya tanpa kunci utama**.'
      % (sum(1 for k in tabel.values() if not k['pk']), len(tabel)))
    w('')
    w('## 3. Asal di sistem lama dan padanannya di model baru')
    w('')
    w('| Tabel | Jalur Pega | Padanan di model baru | Tingkat bukti |')
    w('|---|---|---|---|')
    for nm, k in tabel.items():
        w('| `%s` | %s | %s | %s |'
          % (nm, ('`%s`' % k['jalur']) if k['jalur'] else '*(tidak ada)*',
             ('**%s**' % k['padanan']) if k['padanan'] else '*(tidak ada)*',
             k['bukti'] or '—'))
    w('')
    w('> **Kolom *Padanan di model baru* adalah BACAAN, bukan pemetaan yang mengikat.** Yang '
      'mengikat soal nama kolom dan tipe: `treaty-in/2-to-spec/KAMUS-KOLOM.md`; soal daftar '
      'entitas: `treaty-in/4-erd-dan-tabel-datar/STRUKTUR-DATA.md`.')
    w('')
    w('## 4. Relasi yang menyentuh modul ini — %d' % len(relasi_modul[mod]))
    w('')
    w('| # | Induk | Kard. | Anak | Kunci tamu | ON DELETE | Muncul di lembar |')
    w('|---|---|---|---|---|---|---|')
    for r in relasi_modul[mod]:
        w('| %s | `%s` | %s | `%s` | `%s` | %s | %s |'
          % (r.get('#', ''), r['Induk'], r.get('Kard.', '—'), r['Anak'],
             r.get('Kunci tamu', '—'), r.get('ON DELETE', '—'),
             r.get('Muncul di lembar', '—')))
    w('')
    if relasi_takterklasifikasi:
        w('### 4a. Relasi yang TIDAK muncul di lembar mana pun — %d'
          % len(relasi_takterklasifikasi))
        w('')
        w('Ketiganya ada di *Daftar Relasi* sumbernya, dan kolom **Muncul di lembar**-nya berbunyi')
        w('*(tidak terklasifikasi)*. **Dibawa ke KEDUA berkas**, sebab membuangnya dari keduanya')
        w('akan menghilangkannya sama sekali — dan memilih salah satu berarti memutuskan sesuatu')
        w('yang sumbernya tidak nyatakan.')
        w('')
        w('| # | Induk | Kard. | Anak | Kunci tamu | ON DELETE |')
        w('|---|---|---|---|---|---|')
        for r in relasi_takterklasifikasi:
            w('| %s | `%s` | %s | `%s` | `%s` | %s |'
              % (r.get('#', ''), r['Induk'], r.get('Kard.', '—'), r['Anak'],
                 r.get('Kunci tamu', '—'), r.get('ON DELETE', '—')))
        w('')
        if tabel_relasi_saja:
            w('> **%d tabel di antaranya tidak pernah digambar sebagai kotak** di lembar mana pun: %s. '
              'Ia ada di daftar relasi dan tidak di gambar — dan itu keadaan sumbernya, bukan '
              'kekeliruan perkakas ini.'
              % (len(tabel_relasi_saja), ', '.join('`%s`' % x for x in tabel_relasi_saja)))
            w('')
    w('## 5. Pohon per lembar')
    w('')
    for nama in lbs:
        w('### %s — %d kotak' % (nama, len(lembar[nama]['kotak'])))
        w('')
        pohon_md(pohon(lembar[nama]['kotak']), m)
        w('')
    w('## 6. Catatan & Batas — disalin apa adanya dari sumbernya')
    w('')
    w('| | |')
    w('|---|---|')
    for a, b in d['catatan']:
        if a:
            w('| **%s** | %s |' % (a, b.replace('|', '\\|')))
    w('')
    w('## 7. Batas berkas ini')
    w('')
    w('| Ia TIDAK dapat menyatakan… | Sebabnya |')
    w('|---|---|')
    w('| bahwa himpunan tabelnya **lengkap** | ia salinan sebuah **gambar**. Gambar itu disusun dari '
      'pohon clipboard Pega, dan pohon itu **kurang 340 properti titik buta** (`L-8`, ditagih `M-4`) |')
    w('| bahwa pembagian dua modul ini **benar secara bisnis** | sumbu pemisahnya **nama lembar**, '
      'dan nama lembar adalah keputusan orang yang menggambar — bukan fakta yang dibaca dari ekspor |')
    w('| bahwa `ON DELETE` di sini **perilaku sistem lama** | Catatan & Batas sumbernya sendiri '
      'menyatakannya **usulan rancangan**, bukan perilaku yang terbaca |')
    w('')
    p = os.path.normpath(KELUAR[mod])
    os.makedirs(os.path.dirname(p), exist_ok=True)
    try:
        io.open(p, 'w', encoding='utf-8').write('\n'.join(m) + '\n')
        st = 'DITULIS'
    except Exception as ex:
        st = 'GAGAL — %s' % ex
        tolak['berkas keluaran tidak dapat ditulis'] += 1
    hasil.append((mod, p, len(tabel), kotak_modul[mod], len(relasi_modul[mod]), st))

# ── 6. laporan ─────────────────────────────────────────────────────────────
print('=== pisah-struktur-erd-v2-per-modul.py ===')
print()
print('SUMBER : %s  (%d lembar, %d relasi)' % (d['sumber'], len(lembar), len(d['relasi'])))
print()
for mod, p, nt, nk, nr, st in hasil:
    print('%-22s %s' % (mod, st))
    print('   %-20s %s' % ('berkas', p))
    print('   %-20s %d tabel unik · %d kotak · %d relasi' % ('isi', nt, nk, nr))
print()
print('PENUTUPAN CACAH -- dua arah:')
print('   kotak: %d + %d = %d   sumber: %d   %s'
      % (kotak_modul[A], kotak_modul[B], kotak_modul[A] + kotak_modul[B],
         sum(len(L['kotak']) for L in lembar.values()),
         'MENUTUP' if kotak_modul[A] + kotak_modul[B] == sum(len(L['kotak']) for L in lembar.values()) else 'TIDAK'))
print('   tabel unik gabungan: %d   sumber: %d   %s'
      % (len(semua_tabel), len({k['nama'] for L in lembar.values() for k in L['kotak']}),
         'MENUTUP' if len(semua_tabel) == len({k['nama'] for L in lembar.values() for k in L['kotak']}) else 'TIDAK'))
n_rel_kena = len({r['#'] for r in relasi_modul[A]} | {r['#'] for r in relasi_modul[B]}
                 | {r['#'] for r in relasi_takterklasifikasi})
n_rel_asli = len([r for r in d['relasi'] if r.get('Induk') and r.get('Anak')])
print('   relasi: %d mendarat   dari %d relasi sejati   %s'
      % (n_rel_kena, n_rel_asli, 'MENUTUP' if n_rel_kena == n_rel_asli else 'TIDAK'))
print()
print('BERSAMA kedua modul : %d tabel' % len(bersama))
for mod in (A, B):
    print('HANYA %-16s : %d  %s' % (mod, len(hanya[mod]),
                                    ', '.join(hanya[mod]) if hanya[mod] else '(tidak ada)'))
print()
print('DITOLAK -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-64s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
if tabel_relasi_saja:
    print('   tabel yang hanya ada di Daftar Relasi, apa adanya: %s' % ', '.join(tabel_relasi_saja))
print()
print('-- BATAS PERKAKAS INI --')
print('   Ia memisahkan GAMBAR, bukan MODEL. %d dari %d tabel dipakai KEDUA modul dan'
      % (len(bersama), len(semua_tabel)))
print('   karena itu muncul di kedua berkas -- itu kenyataan sistem lama, bukan penggandaan.')
print('   Dan sumbu pemisahnya NAMA LEMBAR, yang merupakan keputusan orang yang menggambar.')
