# -*- coding: utf-8 -*-
"""
buat-kamus-dan-ddl.py — SATU definisi, DUA keluaran.

    2-to-spec/KAMUS-KOLOM.md
    2-to-spec/ddl-usulan/*.sql

Meniru syarat yang contoh claim-non-prop nyatakan di kepalanya sendiri:
"Dibangkitkan dari definisi yang sama dengan ddl-usulan/. Keduanya tidak dapat berbeda."
Definisinya BUKAN diketik ulang: ia DIURAI dari SPEC-MODEL-DATA.md §10, sehingga
kamus dan DDL tidak dapat menyimpang dari spesifikasi maupun dari satu sama lain.

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa setiap kolom yang dikeluarkannya ADA di §10, beserta tipe, keterisian,
    dan asal-usulnya.
APA YANG TIDAK
    Ia TIDAK dapat menyatakan sebuah entitas LENGKAP. Bila §10 menulis atribut
    dalam bentuk yang tidak dikenalinya, atribut itu hilang tanpa suara -- karena
    itu perkakas ini MEMBANDINGKAN jumlah terurai dengan jumlah yang DITULIS di
    judul bagiannya ("-- N atribut") dan MELAPORKAN setiap selisih.
    Selisih yang dilaporkan bukan kesalahan spesifikasi; ia batas perkakas ini.

PENOLAKAN dilaporkan satu baris per sebab.
"""
import io, os, re, json, sys

AKAR   = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
SPEC   = os.path.join(AKAR, 'SPEC-MODEL-DATA.md')
KELUAR = os.path.join(AKAR, '2-to-spec')
DDL    = os.path.join(KELUAR, 'ddl-usulan')

# ── tipe konseptual -> tipe Oracle ────────────────────────────────────────────
# Kelompoknya DIWARISI dari ADR-0003 lewat SPEC-MODEL-DATA.md §10 pendahuluan,
# yang menyatakan: "Tipe konseptual memakai kelompok ADR-0003 yang diwarisi dari
# modul Claim Non Prop" dan "Angka presisinya TIDAK ditetapkan di sini -- ia
# keputusan gerbang sesi DDL."
# Maka angka di bawah adalah USULAN yang mewarisi modul tetangga, BUKAN keputusan
# sesi ini. Setiap berkas DDL menyatakannya sebagai pernyataan keputusan.
TIPE = {
    'U1': ('NUMBER(38,20)', 'nilai uang'),
    'P1': ('NUMBER(38,20)', 'persentase dan porsi'),
    'P2': ('NUMBER(38,20)', 'tarif brokerage dan pajak'),
    'K1': ('NUMBER(38,20)', 'kurs'),
    'L1': ('NUMBER(38,20)', 'limit layer dan premi deposit'),
    'C1': ('NUMBER(9)',     'cacah dan nomor urut'),
    'T':  ('VARCHAR2(1000 CHAR)', 'teks'),
    'D':  ('DATE',          'tanggal'),
    'E':  ('VARCHAR2(40 CHAR)',  'himpunan tertutup'),
    'R':  ('NUMBER(19)',    'kunci asing'),
}
# KTV-A, 24 September 2026 -- KEPUTUSAN-TANPA-VERIFIKASI.md paragraf 7.
# Yang berubah dari jalan sebelumnya, masing-masing bersebab di sana:
#   P2  NUMBER(11,8)        -> NUMBER(38,20)        9989998 tidak muat di tiga
#                                                   angka depan koma, dan sistem
#                                                   lama MENULISKANNYA ke ROLPct.
#   T   VARCHAR2(255 CHAR)  -> VARCHAR2(1000 CHAR)  255 LEBIH SEMPIT daripada
#                                                   TREATYINDETAIL, yang memakai
#                                                   VARCHAR2(1000) pada 52 kolom.
#   ID_* bertipe C1         -> NUMBER(19)           kunci asing R sudah NUMBER(19);
#                                                   dua ujung satu relasi tidak
#                                                   boleh berbeda lebar.
ID_LEBAR = 'NUMBER(19)'

# Tujuh kolom TEKS BEBAS. Ini BUKAN pengulangan keputusan CLOB yang dicabut:
# CLOB adalah TIPE LAIN, diputuskan perkakas, tidak pernah ditinjau siapa pun,
# dan satu entrinya menunjuk kolom yang tidak ada. Yang di bawah tetap VARCHAR2
# -- yang berubah hanya LEBARNYA, mengikuti lebar yang sistem lama pakai untuk
# teks bebasnya sendiri (ACHIEVEMENT dan T_STORAGE_IMAGE: VARCHAR2(4000)).
KHUSUS = {
    ('VERSI_KONTRAK', 'KETERANGAN'):        'VARCHAR2(4000 CHAR)',
    ('VERSI_KONTRAK', 'CATATAN'):           'VARCHAR2(4000 CHAR)',
    ('VERSI_KONTRAK', 'PENGECUALIAN'):      'VARCHAR2(4000 CHAR)',
    ('VERSI_KONTRAK', 'KETENTUAN_KHUSUS'):  'VARCHAR2(4000 CHAR)',
    ('VERSI_KONTRAK', 'CATATAN_BORDEREAUX'):'VARCHAR2(4000 CHAR)',
    ('JEJAK_PERUBAHAN', 'NILAI_SEBELUM'):   'VARCHAR2(4000 CHAR)',
    ('JEJAK_PERUBAHAN', 'NILAI_SESUDAH'):   'VARCHAR2(4000 CHAR)',
}

tolak = []   # (sebab, berapa)

# ── 1. urai §10 ───────────────────────────────────────────────────────────────
baris = io.open(SPEC, encoding='utf-8').read().split('\n')
JUDUL = re.compile(r'^#{3,4} (10\.\d+[a-z]?)\s+`([A-Z_]{2,})`(.*)$')
NAMA  = re.compile(r'^[A-Z][A-Z0-9_]*$')
PINDAH = re.compile(r'^Dan pada \*\*`([A-Z_]{2,})`\*\* ditambahkan')
pindah_terpakai = []

# Subbagian yang judulnya menyebut sebuah ATRIBUT, bukan entitas.
# §10.2a `SIFAT_MATERIAL_ADDENDUM` adalah atribut ke-48 VERSI_KONTRAK, bukan tabel.
INDUK_SUBBAGIAN = {'10.2a': 'VERSI_KONTRAK'}

ent = None; urut = []; kol = {}; klaim = {}
n_tolak_baris = 0
for l in baris:
    m = JUDUL.match(l)
    if m:
        if m.group(1) in INDUK_SUBBAGIAN:
            ent = INDUK_SUBBAGIAN[m.group(1)]
            continue
        ent = m.group(2)
        if ent not in kol:
            kol[ent] = []; urut.append((m.group(1), ent))
        d = re.search(r'(\d+)\s+atribut', m.group(3) or '')
        if d and ent not in klaim: klaim[ent] = int(d.group(1))
        continue
    if l.startswith('### ') and not JUDUL.match(l): ent = None; continue
    # Sebuah subbagian dapat memuat tabel KEDUA yang milik entitas LAIN, dan sec 10
    # mengumumkannya dengan satu kalimat: "Dan pada **`X`** ditambahkan ...".
    # Tanpa aturan ini, tabel kedua terhisap ke entitas judulnya -- yang sudah
    # terjadi: sec 10.19a menambahkan ID_DOKUMEN_ADDENDUM ke VERSI_KONTRAK, dan
    # kolom itu MENDARAT DI TABEL YANG SALAH sekaligus HILANG dari tabel yang benar.
    m2 = PINDAH.match(l)
    if m2:
        pindah_terpakai.append((m2.group(1), ent))
        ent = m2.group(1)
        continue
    if ent and l.startswith('| ') and l.count('|') >= 5:
        sel = [c.strip() for c in l.strip().strip('|').split('|')]
        nm = sel[0].strip('`* ')
        if not NAMA.match(nm): n_tolak_baris += 1; continue
        # bentuk lima kolom: Nama | Asal | Tipe | Boleh kosong | Catatan
        # bentuk empat kolom: Nama | Asal | Tipe | Boleh kosong
        asal, tipe, kosong = sel[1].strip('`* '), sel[2].strip('* '), sel[3].strip('* ')
        cat = sel[4] if len(sel) > 4 else ''
        kol[ent].append(dict(kolom=nm, asal=asal, tipe=tipe, kosong=kosong, catatan=cat))
tolak.append(('baris tabel di dalam §10 yang bukan baris atribut', n_tolak_baris))

# ── 2. koreksi terukur, masing-masing bersebab ───────────────────────────────
# 2a. PEMULIHAN_LIMIT terhisap ke ekor LAYER (§10.3b subbagian di bawah §10.3)
if 'LAYER' in kol:
    idx = [i for i, c in enumerate(kol['LAYER']) if c['kolom'] == 'ID_PEMULIHAN_LIMIT']
    if idx:
        i = idx[0]
        kol['PEMULIHAN_LIMIT'] = kol['LAYER'][i:]
        kol['LAYER'] = kol['LAYER'][:i]
        urut.append(('10.3b', 'PEMULIHAN_LIMIT')); klaim['PEMULIHAN_LIMIT'] = len(kol['PEMULIHAN_LIMIT'])

# 2b. POTONGAN -- §10.9 sengaja TIDAK mengulang atributnya; ia ada di §14.2
teks = '\n'.join(baris)

def iris(awal, akhir):
    """Iris segmen dokumen lebih dulu. Regex yang menyapu SELURUH dokumen pernah
    memungut 66 baris untuk tabel yang isinya 5 -- penolakan yang tidak terlihat
    karena hasilnya tetap 'berhasil'."""
    try:
        a = teks.index(awal); b = teks.index(akhir, a + 1)
        return teks[a:b]
    except ValueError:
        return ''

seg142 = iris('### 14.2', '### 14.3')
m = re.search(r'\| Nama \| Asal \| Tipe \| Boleh kosong \|\n\|[-|]+\|\n((?:\|.*\n)+)', seg142)
if m:
    kol['POTONGAN'] = []
    for l in m.group(1).strip().split('\n'):
        sel = [c.strip() for c in l.strip().strip('|').split('|')]
        nm = sel[0].strip('`* ')
        if NAMA.match(nm):
            kol['POTONGAN'].append(dict(kolom=nm, asal=sel[1].strip('`* '), tipe=sel[2].strip('* '),
                                        kosong=sel[3].strip('* '), catatan='§14.2'))
    urut.append(('10.9', 'POTONGAN'))
else:
    tolak.append(('POTONGAN tidak terbaca dari §14.2', 1))

# 2c. enam tabel acuan -- satu bentuk seragam di §10.22
ACUAN = ['MATA_UANG', 'JENIS_POTONGAN', 'JENIS_REASURANSI', 'BAHAYA', 'KELOMPOK_TREATY', 'KELAS_BISNIS']
seg1022 = iris('### 10.22', '### 10.23')
m = re.search(r'\| Nama \| Tipe \| Boleh kosong \| Catatan \|\n\|[-|]+\|\n((?:\|.*\n)+)', seg1022)
if m:
    pola = []
    for l in m.group(1).strip().split('\n'):
        sel = [c.strip() for c in l.strip().strip('|').split('|')]
        if len(sel) < 3: continue
        nm = sel[0].strip('`‹›* ')
        if not nm: continue
        pola.append((nm, sel[1].strip('* '), sel[2].strip('* '), sel[3] if len(sel) > 3 else ''))
    for a in ACUAN:
        kol[a] = []
        for nm, tp, ks, ct in pola:
            if nm.startswith('ID_') and 'entitas' in nm: nm = 'ID_' + a
            if nm == 'ID_INDUK' and a != 'JENIS_REASURANSI': continue   # §10.22: hanya JENIS_REASURANSI
            kol[a].append(dict(kolom=nm, asal='tabel acuan (ADR-0038)', tipe=tp, kosong=ks, catatan=ct))
        urut.append(('10.22', a)); klaim[a] = len(kol[a])
else:
    tolak.append(('enam tabel acuan tidak terbaca dari §10.22', 6))

# hapus entitas kosong
kosong_ent = [e for e in kol if not kol[e]]
for e in kosong_ent: del kol[e]
urut = [(s, e) for s, e in urut if e in kol]
if kosong_ent: tolak.append(('entitas berjudul §10.x tanpa satu pun baris atribut: ' + ', '.join(kosong_ent), len(kosong_ent)))

# ── PAKET UANG: bentuk fisiknya, diputuskan 24 September 2026 (P-8) ──────────
# Tiga golongan, dibaca dari kolom *Asal* §10 -- bukan satu bentuk untuk keenam belasnya.
#   A  asalnya DAFTAR per mata uang  -> BARIS ANAK, satu per mata uang
#   B  asalnya Amount+Currency pada satu baris, dan induknya SUDAH per mata uang
#      -> kolom KODE_MATA_UANG pada baris itu sendiri. Ini yang membebaskan
#         INV-08, INV-09, INV-12, yang kunci alaminya memang menyebut mata uang.
#   C  asalnya SKALAR POLOS, tanpa mata uang apa pun di sistem lama
#      -> DITAHAN. Menetapkan mata uang untuknya = mengarang fakta. Menunggu T-6.
PAKET_A = {
    ('LAYER', 'MDP'), ('LAYER', 'MDP_MINIMUM'),
    ('BAGIAN', 'PREMI_BRUTO'), ('BAGIAN', 'PREMI_BRUTO_MINIMUM'),
    ('DETAIL_PROPORSIONAL', 'CADANGAN_PREMI'),
}
PAKET_B = {
    ('EGNPI', 'NILAI_EGNPI'), ('RETENSI_CEDANT', 'NILAI_RETENSI'),
    ('TERMIN', 'NILAI_TERMIN'), ('NILAI_PENYEBARAN', 'NILAI'), ('LAYER', 'LIMIT'),
    # KTV-C koreksi 1: NILAI_BATAS BUKAN golongan C. Sec 10.18 menulis asalnya
    # sendiri 'Earthquake / FloodJab / FloodNation / RSMDLimit + MATA UANGNYA',
    # dan keempat kolom Currency* itu tercatat 'melebur ke paket uangnya'.
    ('BATAS_PER_BAHAYA', 'NILAI_BATAS'),
}
PAKET_C = {
    ('VERSI_KONTRAK', 'BATAS_MAKSIMUM_KELOMPOK'),
    ('VERSI_KONTRAK', 'BATAS_MAKSIMUM_NON_KELOMPOK'),
    ('VERSI_KONTRAK', 'BATAS_PILIHAN'),
    ('LAYER', 'DEDUCTIBLE'), ('LAYER', 'LIMIT_AGREGAT'),
}

# KTV-C, 24 September 2026 -- golongan C memperoleh kolom mata uang BERNAMA,
# boleh kosong. Bukan 'KODE_MATA_UANG' polos: tiga di antaranya duduk di baris
# VERSI_KONTRAK yang SUDAH punya KODE_MATA_UANG_KONTRAK, dan dua di baris LAYER
# yang SUDAH punya KODE_MATA_UANG milik LIMIT (golongan B). Satu nama untuk dua
# arti adalah cacat yang ditanam dengan tangan sendiri.
MATA_UANG_C = {
    ('VERSI_KONTRAK', 'BATAS_MAKSIMUM_KELOMPOK'):     'MATA_UANG_BATAS_KELOMPOK',
    ('VERSI_KONTRAK', 'BATAS_MAKSIMUM_NON_KELOMPOK'): 'MATA_UANG_BATAS_NON_KELOMPOK',
    ('VERSI_KONTRAK', 'BATAS_PILIHAN'):               'MATA_UANG_BATAS_PILIHAN',
    ('LAYER', 'DEDUCTIBLE'):                          'MATA_UANG_DEDUCTIBLE',
    ('LAYER', 'LIMIT_AGREGAT'):                       'MATA_UANG_LIMIT_AGREGAT',
}

# KTV-B, 24 September 2026 -- induk polimorfik menjadi DUA kolom bernama.
# Sec 10 tetap menulis satu atribut konseptual; yang di bawah adalah BENTUK
# FISIKNYA, dan sec 10 sendiri menyerahkan bentuk fisik ke sesi DDL.
INDUK_GANDA = {
    ('POTONGAN', 'ID_INDUK_POTONGAN'):     ('BAGIAN', 'DETAIL_PROPORSIONAL'),
    ('PENYEBARAN', 'ID_INDUK_PENYEBARAN'): ('BAGIAN', 'DETAIL_PROPORSIONAL'),
}


# Golongan B memperoleh KODE_MATA_UANG pada baris yang sama. Disisipkan ke DEFINISI,
# sehingga KAMUS-KOLOM.md dan ddl-usulan/ memperolehnya dari sumber yang sama.
for _e, _v in list(kol.items()):
    _baru = []
    for _c in _v:
        _baru.append(_c)
        if (_e, _c['kolom']) in PAKET_B:
            # BATAS_PER_BAHAYA berbeda pada satu hal, dan sebabnya ditulis di KTV-C:
            # mata uangnya BUKAN bagian kunci alami mana pun (kuncinya ID_BAHAYA,
            # INV-14), jadi NOT NULL hanya akan menggagalkan migrasi baris warisan
            # yang mata uangnya kosong.
            _b_kosong = 'ya' if _e == 'BATAS_PER_BAHAYA' else 'tidak'
            _b_cat = ('**denominasi paket uang di sebelahnya** — bagian kunci alami '
                      '(INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B')
            if _e == 'BATAS_PER_BAHAYA':
                _b_cat = ('**denominasi paket uang di sebelahnya**. Golongan **B**, dipindahkan '
                          'dari golongan C oleh `KTV-C` — §10.18 menulis asalnya *"+ mata uangnya"*. '
                          '**Boleh kosong**: ia bukan bagian kunci alami (INV-14 memakai `ID_BAHAYA`)')
            _baru.append(dict(
                kolom='KODE_MATA_UANG',
                asal='`Currency` / `CurrencyID` pada baris yang sama',
                tipe='T',
                kosong=_b_kosong,
                catatan=_b_cat))
        if (_e, _c['kolom']) in MATA_UANG_C:
            _baru.append(dict(
                kolom=MATA_UANG_C[(_e, _c['kolom'])],
                asal='**tidak ada di sistem lama**',
                tipe='T',
                kosong='ya',
                catatan='**denominasi `%s`** — `KTV-C`. Sistem lama tidak merekam mata uang '
                        'untuk besaran ini sama sekali; kolomnya disediakan sekarang karena '
                        'menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat '
                        'bila `T-6` menyatakan ia selalu mata uang kontrak**' % _c['kolom']))
    kol[_e] = _baru

# KTV-B -- induk polimorfik dimekarkan menjadi DUA kolom bernama + CHECK.
# Dilaporkan, tidak didiamkan: perkakas mencetak berapa kolom yang dimekarkannya.
n_mekar = 0
for _e, _v in list(kol.items()):
    _baru = []
    for _c in _v:
        kunci = (_e, _c['kolom'])
        if kunci in INDUK_GANDA:
            n_mekar += 1
            for _par in INDUK_GANDA[kunci]:
                _baru.append(dict(
                    kolom='ID_' + _par,
                    asal=_c['asal'],
                    tipe='R',
                    kosong='ya',
                    catatan='`KTV-B` — salah satu dari **dua** pelekatan. **Tepat satu** dari '
                            'kedua kolom terisi, dijaga `CHECK`. Menggantikan `%s`, yang tidak '
                            'dapat punya kunci asing karena sasarannya bergantung nilai kolom '
                            'lain (INV-17)' % _c['kolom']))
        else:
            _baru.append(_c)
    kol[_e] = _baru

# Golongan A menjadi tabel anak, satu baris per mata uang.
for _e, _k in sorted(PAKET_A):
    _anak = _e[:12] + '_' + _k[:16]
    _anak = ('NILAI_' + _k)[:30]
    if _anak in kol: continue
    kol[_anak] = [
        dict(kolom='ID_' + _anak, asal='*baru*', tipe='C1', kosong='tidak', catatan=''),
        dict(kolom='ID_' + _e, asal='*baru*', tipe='R', kosong='tidak',
             catatan='induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama'),
        dict(kolom='KODE_MATA_UANG', asal='`Currency` pada baris daftarnya', tipe='T', kosong='tidak',
             catatan='**kunci alami** di dalam induknya — belum bernomor'),
        dict(kolom='NILAI', asal='`Value` pada baris daftarnya', tipe='U1', kosong='tidak', catatan=''),
    ]
    urut.append(('10.x', _anak))

# dedupe urut -- POTONGAN masuk dua kali: sekali dari judul sec 10.9 yang
# sengaja kosong, sekali lagi saat isinya diambil dari sec 14.2.
_lihat = set(); _u = []
for _s, _e in urut:
    if _e in _lihat: continue
    _lihat.add(_e); _u.append((_s, _e))
urut = _u

# ── 3. pemeriksaan diri: terurai versus yang DITULIS di judul ────────────
# Kolom yang DISISIPKAN perkakas ini -- P-8 golongan B, KTV-B, KTV-C -- TIDAK ada
# di sec 10 dan tidak boleh terbaca sebagai selisih spesifikasi. Ia dihitung
# terpisah, sebab selisih yang bercampur dua sebab tidak dapat ditindaklanjuti.
DISISIPKAN = set()
for _e, _v in kol.items():
    for _c in _v:
        if 'P-8 golongan B' in (_c['catatan'] or '') or '`KTV-B`' in (_c['catatan'] or '') \
           or '`KTV-C`' in (_c['catatan'] or '') or 'KTV-C`' in (_c['catatan'] or ''):
            DISISIPKAN.add((_e, _c['kolom']))

selisih = []
for s, e in urut:
    if e not in klaim:
        continue
    sisip = sum(1 for c in kol[e] if (e, c['kolom']) in DISISIPKAN)
    # KTV-B memekarkan SATU atribut sec 10 menjadi DUA kolom fisik: yang bertambah
    # hanya satu per pemekaran, jadi yang dihitung bukan jumlah kolom barunya.
    mekar = sum(1 for (ee, kk) in INDUK_GANDA if ee == e)
    terurai_sec10 = len(kol[e]) - sisip + mekar
    if klaim[e] != terurai_sec10:
        selisih.append((e, terurai_sec10, klaim[e], sisip + mekar))

# ── 3b. batas 30 bita -- perkakas BERHENTI, tidak memotong diam-diam ────────
kepanjangan = []
for e in kol:
    if len(e) > 30:
        kepanjangan.append(('tabel', e, len(e)))
    for c in kol[e]:
        if len(c['kolom']) > 30:
            kepanjangan.append(('kolom %s' % e, c['kolom'], len(c['kolom'])))
if kepanjangan:
    print('BERHENTI -- pengenal melewati batas 30 bita Oracle:')
    for jenis, nm, n in kepanjangan:
        print('   %-24s %-34s %d' % (jenis, nm, n))
    raise SystemExit(1)

# ── 4. keluaran ───────────────────────────────────────────
def oracle(e, c):
    """Tipe Oracle satu kolom. Hasilnya DISIMPAN ke dalam berkas definisi, supaya
    penulis keluaran tidak perlu mengulang tabel tipe -- dua salinan tabel tipe
    adalah dua tempat yang dapat berbeda, dan sudah pernah berbeda."""
    if (e, c['kolom']) in KHUSUS:
        return KHUSUS[(e, c['kolom'])]
    # KTV-A: pengenal dan rujukan HARUS sama lebar di kedua ujung relasi.
    if c['kolom'].startswith('ID_'):
        return ID_LEBAR
    t = c['tipe'].split()[0].strip('`*')
    return TIPE.get(t, ('VARCHAR2(1000 CHAR)', '?'))[0]

def nullable(c):
    return 'N' if c['kosong'].lower().startswith('tidak') else 'Y'

n_id_dilebarkan = 0
for e in kol:
    for c in kol[e]:
        c['oracle'] = oracle(e, c)
        if c['kolom'].startswith('ID_') and c['tipe'].split()[0].strip('`*') == 'C1':
            n_id_dilebarkan += 1

os.makedirs(DDL, exist_ok=True)
json.dump({'urut': urut, 'kolom': kol, 'klaim': klaim, 'selisih': selisih,
           'paket_c': sorted(list(x) for x in PAKET_C),
           'induk_ganda': dict(('%s.%s' % k, list(v)) for k, v in INDUK_GANDA.items()),
           'mata_uang_c': dict(('%s.%s' % k, v) for k, v in MATA_UANG_C.items()),
           'tipe': dict((k, list(v)) for k, v in TIPE.items()),
           'khusus': dict(('%s.%s' % k, v) for k, v in KHUSUS.items())},
          io.open(os.path.join(AKAR, 'alat', 'definisi-skema-treaty-masuk.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)

print('=== buat-kamus-dan-ddl.py ===')
print('entitas         : %d' % len(urut))
print('kolom seluruhnya: %d' % sum(len(v) for v in kol.values()))
print()
print('DITOLAK -- satu baris per sebab:')
for s, n in tolak: print('   %-70s %d' % (s, n))
print()
print('TABEL KEDUA di dalam satu subbagian, dipindahkan ke entitas yang benar:')
for _tgt, _asal in pindah_terpakai:
    print('   %-24s <- tabel kedua di dalam subbagian %s' % (_tgt, _asal))
if not pindah_terpakai:
    print('   (tidak ada)')
print()
print('DISISIPKAN perkakas ini -- BUKAN dari sec 10, dan dihitung terpisah:')
print('   %-58s %d' % ('kolom mata uang golongan B (P-8) dan C (KTV-C)',
                       sum(1 for x in DISISIPKAN)))
print('   %-58s %d' % ('atribut induk polimorfik dimekarkan jadi dua (KTV-B)', n_mekar))
print('   %-58s %d' % ('kolom ID_* dilebarkan C1 -> NUMBER(19) (KTV-A)', n_id_dilebarkan))
print()
print('PEMERIKSAAN DIRI -- terurai versus jumlah yang DITULIS di judul sec 10')
print('(sesudah kolom sisipan dan pemekaran dikeluarkan dari hitungan):')
if not selisih:
    print('   tidak ada selisih.')
for e, a, b, sis in selisih:
    print('   %-24s terurai %2d, judulnya menulis %2d  -> SELISIH %+d  (sisipan %d)'
          % (e, a, b, a - b, sis))
print()
for s, e in urut: print('   %-8s %-26s %2d kolom' % (s, e, len(kol[e])))
