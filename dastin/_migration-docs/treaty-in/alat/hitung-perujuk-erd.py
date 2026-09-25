# -*- coding: utf-8 -*-
r"""
hitung-perujuk-erd.py -- SEBELUM memindahkan berkas apa pun, hitung siapa yang merujuknya.

    diperiksa : 4-erd-dan-tabel-datar/*  dan  alat/*
    korpus    : seluruh _migration-docs/  (treaty-in, treaty-in-adjustment, claim-non-prop)
    keluaran  : alat/perujuk-erd.json  + laporan

KENAPA ADA
    Rujukan yang putus di berkas markdown TIDAK MENIMBULKAN GALAT, dan tidak ada
    yang menjalankannya untuk mengetahui. Memindahkan berkas yang masih dirujuk
    tujuh tempat merusak tujuh berkas tanpa satu pesan pun.

    Aturan yang sudah berdiri di _arsip/ISI-FOLDER.md: BASI BUKAN BERARTI TIDAK
    DIPAKAI. Yang diarsipkan hanya yang TERGANTIKAN UTUH oleh berkas lain yang
    menjawab pertanyaan yang sama dengan lebih baik.

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa sebuah nama berkas DISEBUT atau TIDAK PERNAH DISEBUT di korpus, dan
    oleh berkas mana saja. Ia juga memisahkan penyebutan oleh berkas yang SUDAH
    DIARSIPKAN -- sebab perujuk yang sendirinya sudah di dalam arsip bukan alasan
    menahan sesuatu di luar arsip.

APA YANG TIDAK
    Penyebutan BUKAN pembacaan, dan ketiadaan penyebutan BUKAN bukti tidak
    dipakai. Sebuah berkas dapat dipakai orang tanpa berkas lain menyebutnya.
    Perkakas ini karena itu MENGUSULKAN, ia tidak memindahkan. Pemindahannya
    perintah pemilik proses.

    Ia juga tidak membaca xlsx: rujukan dari dalam sel Excel tidak terlihat.
    Cacah berkas yang tidak terbaca dilaporkan.

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/hitung-perujuk-erd.py
"""
import io, os, re, json, sys
from collections import Counter, OrderedDict

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
DOCS = os.path.abspath(os.path.join(AKAR, '..'))
PERIKSA = [os.path.join(AKAR, '4-erd-dan-tabel-datar'), os.path.join(AKAR, 'alat')]
BINER = ('xlsx', 'xls', 'png', 'jpg', 'jpeg', 'pdf', 'zip', 'docx')
# Keluaran perkakas ini sendiri memuat SETIAP nama berkas, sehingga bila ikut
# disapu ia menaikkan cacah semuanya satu -- pemeriksa yang mencemari
# pengukurannya sendiri.
KELUARAN_SENDIRI = ('perujuk-erd.json',)

tolak = Counter()

# ── 1. Daftar berkas yang DIPERIKSA ────────────────────────────────────────
sasaran = OrderedDict()
for folder in PERIKSA:
    for f in sorted(os.listdir(folder)):
        p = os.path.join(folder, f)
        if os.path.isdir(p):
            tolak['folder di dalam folder yang diperiksa (tidak ikut dihitung)'] += 1
            continue
        if f.startswith('~$'):
            tolak['berkas kunci Excel ~$ (bukan artefak)'] += 1
            continue
        rel = os.path.relpath(p, AKAR).replace('\\', '/')
        sasaran[rel] = dict(nama=f, bita=os.path.getsize(p),
                            folder=os.path.basename(folder), perujuk=[], arsip=[])

# ── 2. Korpus ──────────────────────────────────────────────────────────────
korpus, tak_terbaca = [], 0
for dp, dn, fn in os.walk(DOCS):
    for f in fn:
        full = os.path.join(dp, f)
        if f in KELUARAN_SENDIRI:
            tolak['keluaran perkakas ini sendiri (memuat setiap nama berkas)'] += 1
            continue
        if f.rsplit('.', 1)[-1].lower() in BINER:
            tak_terbaca += 1
            tolak['berkas biner yang isinya TIDAK dapat disapu (xlsx dll)'] += 1
            continue
        try:
            korpus.append((os.path.abspath(full),
                           io.open(full, encoding='utf-8', errors='ignore').read()))
        except Exception:
            tak_terbaca += 1
            tolak['berkas yang tidak dapat dibuka'] += 1

# ── 3. Hitung ──────────────────────────────────────────────────────────────
# Berkas KATALOG menyebut hampir setiap berkas di folder yang didaftarnya. Ia
# menaikkan cacah SEMUA berkas satu, sehingga "disebut 1 kali" oleh katalog sama
# artinya dengan "tidak disebut siapa pun". Karena itu ia dipisahkan, bukan
# dibuang -- perujuk katalog tetap harus diperbaiki bila berkasnya berpindah.
KATALOG = ('ISI-FOLDER.md', 'INVENTARIS-BERKAS-INDUK.md')

# Nama berkas dicari dengan BATAS KIRI. Tanpa itu "TABEL-DATAR.md" cocok di
# dalam "RINGKASAN-TABEL-DATAR.md" -- dan cacahnya membengkak oleh berkas yang
# sama sekali bukan dia. Batas kanan tidak perlu: ekstensinya sudah mengunci.
pola = {rel: re.compile(r'(?<![0-9A-Za-z_.-])' + re.escape(s['nama']))
        for rel, s in sasaran.items()}

for rel, s in sasaran.items():
    diri = os.path.abspath(os.path.join(AKAR, rel))
    s['katalog'] = []
    for jalur, teks in korpus:
        if jalur == diri:
            continue
        if pola[rel].search(teks):
            r = os.path.relpath(jalur, DOCS).replace('\\', '/')
            if '_arsip/' in r:
                s['arsip'].append(r)
            elif os.path.basename(jalur) in KATALOG:
                s['katalog'].append(r)
            else:
                s['perujuk'].append(r)

json.dump(sasaran, io.open(os.path.join(AKAR, 'alat', 'perujuk-erd.json'), 'w',
                           encoding='utf-8'), ensure_ascii=False, indent=1)

# ══════════════════════════ CETAK ══════════════════════════════════════════
P = print
P('=== hitung-perujuk-erd.py ===')
P()
P('KORPUS  : %d berkas terbaca, %d tidak terbaca' % (len(korpus), tak_terbaca))
P('SASARAN : %d berkas di 4-erd-dan-tabel-datar/ dan alat/' % len(sasaran))
P()
nol = [(r, s) for r, s in sasaran.items() if not s['perujuk']]
P('TIDAK DISEBUT SATU BERKAS KERJA PUN  (katalog dan arsip tidak dihitung) : %d' % len(nol))
P('   -- "katalog" = ISI-FOLDER.md / INVENTARIS-BERKAS-INDUK.md, yang menyebut')
P('      hampir SETIAP berkas dan karena itu menaikkan cacah semuanya satu.')
P()
for r, s in sorted(nol, key=lambda x: (x[1]['folder'], -x[1]['bita'])):
    tanda = []
    if s['katalog']:
        tanda.append('katalog %d' % len(s['katalog']))
    if s['arsip']:
        tanda.append('_arsip %d' % len(s['arsip']))
    P('   %-52s %8d bita   %s' % (r, s['bita'], ' · '.join(tanda) or 'NOL di mana pun'))
P()
P('DISEBUT 1-2 BERKAS SAJA -- calon terkuat sesudah yang nol:')
for r, s in sorted(sasaran.items(), key=lambda x: len(x[1]['perujuk'])):
    n = len(s['perujuk'])
    if 1 <= n <= 2:
        P('   %-52s %d  <- %s' % (r, n, ', '.join(s['perujuk'])))
P()
P('DISEBUT BANYAK -- JANGAN DIPINDAHKAN tanpa memperbaiki perujuknya:')
for r, s in sorted(sasaran.items(), key=lambda x: -len(x[1]['perujuk'])):
    n = len(s['perujuk'])
    if n >= 3:
        P('   %-52s %d' % (r, n))
P()
# ── Nama yang BERTABRAKAN: ada di lebih dari satu tempat di korpus ─────────
# Sapuan nama telanjang memungut perujuk milik berkas bernama sama di modul
# lain. Cacahnya diukur, bukan diandaikan nol.
semua_nama = Counter()
for dp, dn, fn in os.walk(DOCS):
    for f in fn:
        semua_nama[f] += 1
tabrakan = [(r, s, semua_nama[s['nama']]) for r, s in sasaran.items()
            if semua_nama[s['nama']] > 1]
P('NAMA YANG BERTABRAKAN -- ada di lebih dari satu tempat di korpus : %d' % len(tabrakan))
P('   Untuk berkas ini, cacah perujuk di atas TERLALU BESAR: sapuan nama telanjang')
P('   memungut penyebutan yang sebenarnya menunjuk berkas bernama sama di modul lain.')
for r, s, n in sorted(tabrakan, key=lambda x: -len(x[1]['perujuk'])):
    P('   %-52s %d salinan · %d perujuk terhitung' % (r, n, len(s['perujuk'])))
P()
P('DITOLAK / TIDAK DIHITUNG -- satu baris per sebab:')
for k, v in tolak.most_common():
    P('   %-64s %d' % (k, v))
P()
P('-- BATAS PERKAKAS INI, dan ia menentukan cara membaca angka di atas --')
P('   1. PENYEBUTAN BUKAN PEMBACAAN. Sebuah berkas yang disebut tujuh kali')
P('      bisa saja tidak pernah dibuka siapa pun.')
P('   2. NOL PENYEBUTAN BUKAN BUKTI TIDAK DIPAKAI. Orang membuka berkas tanpa')
P('      berkas lain menyebutnya -- terutama xlsx dan html, yang dibuka langsung.')
P('   3. Isi %d berkas BINER tidak disapu. Rujukan dari dalam sel Excel TIDAK' % tak_terbaca)
P('      TERLIHAT di sini.')
P('   4. Ia MENGUSULKAN, ia tidak memindahkan.')
