# -*- coding: utf-8 -*-
r"""
lengkapi-erd-adjustment.py -- menyatukan ENAM sumber struktur menjadi satu bahan ERD.

    masukan : 4-erd-dan-tabel-datar/erd-skema-baru.json     35 tabel, dari ddl-usulan/
              alat/himpunan-struktur.json                   hasil cocok-enam-sumber
              4-erd-dan-tabel-datar/peta-nama-tabel-treatyin.tsv   nama tabel datar lama
              4-erd-dan-tabel-datar/datar-treatyin-lama.csv  985 simpul pohon Pega
    keluaran: 4-erd-dan-tabel-datar/erd-skema-terpadu.json  37 tabel + asal Pega + nama lama

KENAPA ADA
    erd-skema-baru.json dibangkitkan dari ddl-usulan/, dan ddl-usulan/ berisi 35 tabel.
    STRUKTUR-DATA.md -- yang MENGIKAT soal daftar entitas (urutan wewenang butir 4) --
    menyebut 37 di dalam gelombang ini. Dua yang hilang, NILAI_SELISIH dan
    BESARAN_DAPAT_DISESUAIKAN, adalah SATU-SATUNYA dua entitas milik modul Adjustment.

    Sebabnya terbaca di SPEC-MODEL-DATA.md sec 11.3: "modul Adjustment, embargo".
    Embargo itu sudah lewat. Penagihnya tidak dicabut, dan gambar yang dibangkitkan
    dari himpunan yang bolong TERLIHAT LENGKAP.

APA YANG PERKAKAS INI TAMBAHKAN, dan dari mana
    Tiap atribut kedua entitas itu MEMBAWA KUTIPAN SUMBERNYA di kolom 'dasar'.
    Tidak satu pun ditulis tanpa sumber.

APA YANG PERKAKAS INI TIDAK LAKUKAN, dan ini yang terpenting
    Ia TIDAK melengkapi apa yang sumbernya diam. Kedua entitas ini BELUM pernah
    melewati sec 10 SPEC-MODEL-DATA -- penamaan, tipe, dan keterisian tingkat atribut --
    sebagaimana ke-35 lainnya. Yang belum diputuskan dicetak sebagai LUBANG di dalam
    kartunya sendiri, bukan ditambal dengan tebakan.

    Ia juga TIDAK menulis DDL. Menggambar sebuah tabel bukan membangunnya.

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/lengkapi-erd-adjustment.py
"""
import io, os, re, csv, json, sys
from collections import Counter, OrderedDict

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
ERD = os.path.join(AKAR, '4-erd-dan-tabel-datar')
SUMBER = os.path.join(ERD, 'erd-skema-baru.json')
HIMP = os.path.join(AKAR, 'alat', 'himpunan-struktur.json')
PETA = os.path.join(ERD, 'peta-nama-tabel-treatyin.tsv')
POHON = os.path.join(ERD, 'datar-treatyin-lama.csv')
KELUAR = os.path.join(ERD, 'erd-skema-terpadu.json')

tolak = Counter()
d = json.load(io.open(SUMBER, encoding='utf-8'))
h = json.load(io.open(HIMP, encoding='utf-8'))

# ══ 1. DUA ENTITAS ADJUSTMENT -- tiap atribut membawa dasarnya ══════════════
# Kolom yang TIDAK ada di sini adalah kolom yang sumbernya DIAM, bukan kolom
# yang saya putuskan tidak ada.
ADJ = OrderedDict()

ADJ['BESARAN_DAPAT_DISESUAIKAN'] = dict(
    arti='tabel acuan: besaran apa yang boleh disesuaikan, dan satuannya',
    induk=None, sekat='B -- dipakai bersama kedua modul',
    lembar='Penyesuaian',
    pk='ID_BESARAN',
    kunci_alami='KODE',
    isi_awal='31 baris sebagai AWAL, bukan sebagai batas (G3)',
    kolom=[
        dict(kolom='ID_BESARAN', tipe='C1', kosong='tidak',
             dasar='STRUKTUR-DATA.md sec 2 -- kunci utama'),
        dict(kolom='KODE', tipe='T', kosong='tidak',
             dasar='SPEC-MODEL-DATA sec 10.22 -- bentuk seragam tabel acuan; '
                   'kunci alami, unik seluruh tabel'),
        dict(kolom='NAMA', tipe='T', kosong='tidak',
             dasar='SPEC-MODEL-DATA sec 10.22'),
        dict(kolom='SATUAN_BESARAN', tipe='E', kosong='tidak',
             dasar='KEPUTUSAN-SAMBUNGAN G3 butir 2 -- UANG / PERSENTASE, '
                   'himpunan TERTUTUP sejak awal'),
        dict(kolom='AKTIF', tipe='E', kosong='tidak',
             dasar='SPEC-MODEL-DATA sec 10.22 -- nilai lama tidak dihapus, ia dimatikan'),
    ],
    lubang=[
        'Kunci alaminya DISEBUT DUA KALI DENGAN NAMA BERBEDA: STRUKTUR-DATA sec 2 '
        'berbunyi "nama besaran", SPEC-MODEL-DATA sec 10.22 berbunyi KODE untuk '
        'keenam tabel acuan. Yang dipakai di sini KODE, sebab sec 10.22 menyebut '
        'alasannya -- nama berubah ejaannya, kode tidak. BELUM DIPUTUSKAN pemilik proses.',
    ])

ADJ['NILAI_SELISIH'] = dict(
    arti='satu besaran yang berubah pada sebuah versi terhadap versi dasarnya',
    induk='VERSI_KONTRAK', sekat='ADJ -- milik modul Treaty In Adjustment',
    lembar='Penyesuaian',
    pk='ID_NILAI_SELISIH',
    kunci_alami='ID_BESARAN + KUNCI_PADANAN, di dalam versinya',
    isi_awal=None,
    kolom=[
        dict(kolom='ID_NILAI_SELISIH', tipe='C1', kosong='tidak',
             dasar='STRUKTUR-DATA.md sec 2 -- kunci utama'),
        dict(kolom='ID_VERSI_KONTRAK', tipe='R', kosong='tidak',
             dasar='STRUKTUR-DATA.md sec 2 -- induknya VERSI_KONTRAK'),
        dict(kolom='ID_BESARAN', tipe='R', kosong='tidak',
             dasar='KEPUTUSAN-SAMBUNGAN G3 butir 1 -- merujuk tabel acuan, '
                   'BUKAN teks bebas (ADR-0038)'),
        dict(kolom='JENIS_INDUK', tipe='E', kosong='tidak',
             dasar='SEAM-ADJUSTMENT.md sec 3 -- "Letak besaran: JENIS_INDUK, KUNCI_PADANAN"'),
        dict(kolom='KUNCI_PADANAN', tipe='T', kosong='tidak',
             dasar='SEAM-ADJUSTMENT.md sec 3 + tiket 06 -- identitas bisnis yang STABIL '
                   'lintas versi; MAJEMUK untuk Bentuk B; mata uang ADA DI DALAMNYA (E-1)'),
        dict(kolom='ID_MATA_UANG', tipe='R', kosong='ya',
             dasar='KEPUTUSAN-SAMBUNGAN G3 butir 2 -- WAJIB terisi bila satuannya UANG, '
                   'WAJIB kosong bila PERSENTASE; ditegakkan CHECK'),
        dict(kolom='NILAI', tipe='U1', kosong='tidak',
             dasar='KEPUTUSAN-SAMBUNGAN G3 butir 3 -- bilangan eksak, tidak pernah teks '
                   '(ADR-0003)'),
    ],
    lubang=[
        'PAKET UANGNYA BELUM LENGKAP. G3 butir 2 berbunyi baris ber-satuan UANG '
        '"wajib membawa mata uang DAN PAKET UANGNYA". Paket uang menurut ADR-0007 / '
        '0029 / 0039 lebih dari dua kolom -- SEAM sec 3 menyebut TINGKAT_PENCATATAN, '
        'PERSEN_BAGIAN_DIPAKAI, KURS, TANGGAL_KURS, SUMBER_KURS, NILAI_IDR pada bentuk '
        'BACA. Mana di antaranya ikut TERSIMPAN di sini belum diputuskan.',
        'BERAPA PENUNJUK VERSI YANG DISIMPAN -- DUA SUMBER BERTENTANGAN. '
        'KEPUTUSAN-SAMBUNGAN G2 berbunyi "Dua penunjuk disimpan, bukan satu": '
        'ID_VERSI_KONTRAK_LAMA dan ID_VERSI_KONTRAK_BARU. STRUKTUR-DATA sec 1.1 '
        'berbunyi ID_VERSI_KONTRAK_DASAR "adalah satu-satunya penunjuk yang disimpan, '
        'dan ia tidak punya pasangan". Gambar ini memakai SATU -- induk VERSI_KONTRAK -- '
        'sebab butir 4 urutan wewenang menjadikan STRUKTUR-DATA mengikat soal daftar '
        'entitas dan induknya. Ini PILIHAN WEWENANG, bukan adjudikasi.',
        'TIPE ORACLE-nya BELUM melewati KTV-A. Kelompok tipe di KAMUS-KOLOM sec 1 '
        'diputuskan untuk 35 tabel; kedua tabel ini tidak ikut di dalamnya.',
    ])

RELASI_ADJ = [
    dict(induk='VERSI_KONTRAK', anak='NILAI_SELISIH', kolom='ID_VERSI_KONTRAK',
         kolom_induk='ID_VERSI_KONTRAK', kard='1:N', on_delete='(belum diputuskan)',
         dasar='STRUKTUR-DATA.md sec 2'),
    dict(induk='BESARAN_DAPAT_DISESUAIKAN', anak='NILAI_SELISIH', kolom='ID_BESARAN',
         kolom_induk='ID_BESARAN', kard='1:N', on_delete='(belum diputuskan)',
         dasar='KEPUTUSAN-SAMBUNGAN G3 butir 1'),
    dict(induk='MATA_UANG', anak='NILAI_SELISIH', kolom='ID_MATA_UANG',
         kolom_induk='ID_MATA_UANG', kard='1:N', on_delete='(belum diputuskan)',
         dasar='KEPUTUSAN-SAMBUNGAN G3 butir 2'),
]

# ══ 2. ASAL PEGA per entitas, dari peta nama + pohon ════════════════════════
asal, nama_lama = {}, {}
with io.open(PETA, encoding='utf-8') as fh:
    for r in csv.DictReader(fh, delimiter='\t'):
        nt = (r.get('NAMA_T') or '').strip()
        if not nt or 'TIDAK DIBUAT' in nt:
            tolak['baris peta nama bertanda TIDAK DIBUAT'] += 1
            continue
        jl = (r.get('JALUR_PEGA') or '').strip()
        for cal in re.split(r'[+·]', r.get('PADANAN_DDL') or ''):
            cal = cal.strip().strip('`')
            if not re.match(r'^[A-Z][A-Z0-9_]+$', cal):
                continue
            nama_lama.setdefault(cal, []).append(nt)
            if jl:
                asal.setdefault(cal, []).append(jl)

# Sumber asal KEDUA, dan ia yang mutakhir: SPEC-MODEL-DATA sec 2.3 ("Sumber di
# JSON") dan sec 10.22 ("Isi awal dari"). Dipakai karena peta nama tabel datar
# memuat 17 padanan yang namanya sudah tidak ada di DDL -- pemeriksaan C.
SPEC = os.path.join(AKAR, 'SPEC-MODEL-DATA.md')
tspec = io.open(SPEC, encoding='utf-8').read()
blok = re.search(r'### 2\.3 Daftar entitas lengkap(.*?)\n## 3\.', tspec, re.S)
blok10 = re.search(r'### 10\.22 (.*?)\n### 10\.23', tspec, re.S)
n_spec = 0
for b in (blok, blok10):
    if not b:
        tolak['blok SPEC-MODEL-DATA yang tidak terbaca (sec 2.3 / sec 10.22)'] += 1
        continue
    for l in b.group(1).split('\n'):
        m = re.match(r'^\|\s*\*{0,2}`([A-Z][A-Z0-9_]+)`\*{0,2}\s*\|\s*([^|]+?)\s*\|', l)
        if not m:
            continue
        nm, sm = m.group(1), m.group(2).strip()
        if not sm or sm in ('—', '-'):
            tolak['baris sec 2.3 / 10.22 yang kolom sumbernya kosong'] += 1
            continue
        # Tanda tebal markdown dibuang -- keluarannya xlsx dan html, bukan markdown
        sm = sm.replace('`', '').replace('**', '').strip()
        asal.setdefault(nm, []).append(sm)
        n_spec += 1

# simpul pohon, untuk menyatakan cacah dan keberadaannya di ekspor Adjustment
pohon = {}
with io.open(POHON, encoding='utf-8-sig') as fh:
    for r in csv.DictReader(fh):
        pohon[r['JALUR']] = r

# ══ 3. SUSUN keluaran ══════════════════════════════════════════════════════
tabel = OrderedDict(d['tabel'])
relasi = list(d['relasi'])
lembar = OrderedDict(d['lembar']) if isinstance(d['lembar'], dict) else OrderedDict(d['lembar'])
di_lembar = dict(d.get('di_lembar', {}))

# Tipe Oracle disamakan dengan yang SUDAH DIPAKAI ddl-usulan/, bukan dikarang:
# kunci utama dan kunci asing NUMBER(19), teks VARCHAR2(1000 CHAR), himpunan
# tertutup VARCHAR2(40 CHAR), uang NUMBER(38,20). Lihat KAMUS-KOLOM sec 1.
ORA = {'C1': 'NUMBER(19)', 'R': 'NUMBER(19)', 'T': 'VARCHAR2(1000 CHAR)',
       'E': 'VARCHAR2(40 CHAR)', 'U1': 'NUMBER(38,20)', 'D': 'DATE'}

LEMBAR_ADJ = 'Penyesuaian — modul Adjustment'
for nm, t in ADJ.items():
    if nm in tabel:
        tolak['entitas Adjustment yang TERNYATA sudah ada di DDL (tidak ditimpa)'] += 1
        continue
    kol = []
    for k in t['kolom']:
        if k['tipe'] not in ORA:
            tolak['kolom bertipe konseptual yang tidak punya padanan Oracle'] += 1
            continue
        kol.append(dict(nama=k['kolom'], tipe=ORA[k['tipe']],
                        wajib=(k['kosong'] == 'tidak'), dasar=k['dasar']))
    tabel[nm] = dict(
        nama=nm, pk=t['pk'], kolom=kol,
        fk=[dict(kolom=r['kolom'], induk=r['induk'], kolom_induk=r['kolom_induk'],
                 on_delete=None, constraint='(BELUM ADA)')
            for r in RELASI_ADJ if r['anak'] == nm],
        cek=[], berkas='(BELUM ADA — tidak ada berkas di ddl-usulan/)',
        belum_ddl=True, arti=t['arti'], sekat=t['sekat'],
        kunci_alami_teks=t['kunci_alami'], isi_awal=t['isi_awal'], lubang=t['lubang'])

no = max([r.get('no', 0) for r in relasi] or [0])
for r in RELASI_ADJ:
    no += 1
    relasi.append(dict(no=no, induk=r['induk'], anak=r['anak'], kolom=r['kolom'],
                       kard=r['kard'], on_delete=r['on_delete'],
                       lembar_induk=next((L for L, isi in lembar.items()
                                          if r['induk'] in isi), LEMBAR_ADJ),
                       lembar_anak=LEMBAR_ADJ, dasar=r['dasar']))
lembar[LEMBAR_ADJ] = list(ADJ.keys())
for nm in ADJ:
    di_lembar.setdefault(nm, []).append(LEMBAR_ADJ)

# asal Pega dan nama lama ditempelkan ke SETIAP tabel yang punya
tempel = 0
for nm in tabel:
    a = sorted(set(asal.get(nm, [])))
    nl = sorted(set(nama_lama.get(nm, [])))
    if a or nl:
        tabel[nm]['asal_pega'] = a
        tabel[nm]['nama_lama'] = nl
        tempel += 1

d2 = dict(d)
d2['tabel'] = tabel
d2['relasi'] = relasi
d2['lembar'] = lembar
d2['di_lembar'] = di_lembar
d2['adj_lubang'] = {k: v['lubang'] for k, v in ADJ.items()}
d2['peta_belum'] = h.get('peta_basi', [])
d2['tanpa_nasib'] = len(h.get('tanpa_nasib', []))
d2['cacah_simpul'] = h.get('cacah', {}).get('simpul', 0)
d2['sumber_rantai'] = [
    'ddl-usulan/*.sql -> bangkitkan-erd-skema-baru.py -> erd-skema-baru.json  (35 tabel)',
    'erd-skema-baru.json + STRUKTUR-DATA.md + KEPUTUSAN-SAMBUNGAN + SEAM-ADJUSTMENT'
    ' -> lengkapi-erd-adjustment.py -> erd-skema-terpadu.json  (37 tabel)',
    'peta-nama-tabel-treatyin.tsv -> kolom "nama lama" tiap kartu',
    'datar-treatyin-lama.csv (985 simpul) -> kolom "asal Pega" tiap kartu',
]
json.dump(d2, io.open(KELUAR, 'w', encoding='utf-8'), ensure_ascii=False, indent=1)

# ══════════════════════════ CETAK ══════════════════════════════════════════
P = print
P('=== lengkapi-erd-adjustment.py ===')
P()
P('MASUK dari ddl-usulan/            : %d tabel' % len(d['tabel']))
P('DITAMBAHKAN dari to-spec Adjustment: %d tabel' % len(ADJ))
for nm, t in ADJ.items():
    P('   %-28s %d kolom berdasar · %d lubang dicetak · induk %s'
      % (nm, len(t['kolom']), len(t['lubang']), t['induk'] or '(berdiri sendiri)'))
P('KELUAR                            : %d tabel · %d relasi · %d lembar'
  % (len(tabel), len(relasi), len(lembar)))
P()
P('ASAL PEGA dan NAMA LAMA ditempelkan: %d dari %d tabel' % (tempel, len(tabel)))
P('   -- dari peta nama tabel datar   : %d padanan' % len(nama_lama))
P('   -- dari SPEC-MODEL-DATA 2.3+10.22: %d baris' % n_spec)
P('   tabel TANPA asal Pega           : %d' % (len(tabel) - tempel))
for nm in tabel:
    if 'asal_pega' not in tabel[nm]:
        P('        - %s' % nm)
P()
P('DITOLAK / TIDAK DIKERJAKAN -- satu baris per sebab:')
for k, v in tolak.most_common():
    P('   %-62s %d' % (k, v))
if not tolak:
    P('   (tidak ada)')
P()
P('-- BATAS PERKAKAS INI, dan ia bagian dari keluarannya --')
P('   1. Ia MENGGAMBAR kedua tabel; ia TIDAK MEMBANGUNNYA. Nol berkas di')
P('      ddl-usulan/, nol pasal di KAMUS-KOLOM.md. Gambar bukan DDL.')
P('   2. Kolom yang TIDAK ada di sini bukan kolom yang saya putuskan tidak ada --')
P('      ia kolom yang sumbernya DIAM. Keduanya belum melewati sec 10.')
P('   3. Satu pertentangan sumber DIPILIH, bukan diadili: berapa penunjuk versi')
P('      yang disimpan. Lihat lubang di kartu NILAI_SELISIH.')
