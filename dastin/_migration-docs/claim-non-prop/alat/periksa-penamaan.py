# -*- coding: utf-8 -*-
"""Periksa seluruh pengenal di ddl-usulan/ terhadap aturan penamaan SPEC bagian 16.

HANYA MEMBACA. Berkas ini tidak pernah mengubah apa pun.

Aturan yang diperiksa:
  1. Tidak ada pengenal melewati 30 byte  -> keluar dengan kode 1 bila ada.
  2. Tidak ada sisa singkatan dari daftar yang dibatalkan 19 September 2026.
  3. Tidak ada nama constraint yang dipakai dua kali.

Jalankan:  python alat/periksa-penamaan.py
Keluar 0 bila bersih, 1 bila ada yang melanggar.
"""
from __future__ import print_function
import io, os, re, sys, glob, collections

BATAS = 30
AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '2-to-spec', 'ddl-usulan')

# Singkatan dari daftar yang DIBATALKAN. Tidak boleh muncul lagi sebagai kata
# utuh di dalam pengenal mana pun. Lihat SPEC bagian 16.1 dan tiket 02 (mati).
SINGKATAN_LAMA = ('KLM','LYR','MTU','AKS','ALK','ADJ','RTC','PRP','RKP','OBP',
                  'KRN','DOK','EST','TTB','TRF','KRK','NKM','DPT','MGK','MGP',
                  'MGN','PBK','PNY')

POLA = [
    ('tabel',       r'CREATE TABLE KLAIMNP\.(\w+)'),
    ('view',        r'CREATE OR REPLACE VIEW KLAIMNP\.(\w+)'),
    ('sequence',    r'CREATE SEQUENCE KLAIMNP\.(\w+)'),
    ('constraint',  r'ADD CONSTRAINT (\w+)'),
    ('index',       r'CREATE (?:UNIQUE )?INDEX KLAIMNP\.(\w+)'),
    ('akun/peran',  r'CREATE (?:USER|ROLE|TABLESPACE) (\w+)'),
    ('kebijakan',   r'CREATE AUDIT POLICY (\w+)'),
]

def tanpa_komentar(teks):
    """Komentar dibuang supaya prosa tidak ikut tercacah sebagai pengenal."""
    return u'\n'.join(baris.split('--')[0] for baris in teks.split('\n'))

def kumpulkan():
    per_jenis = collections.defaultdict(set)
    constraint_per_berkas = []
    for berkas in sorted(glob.glob(os.path.join(AKAR, '*.sql'))):
        kode = tanpa_komentar(io.open(berkas, encoding='utf-8').read())
        for jenis, pola in POLA:
            per_jenis[jenis] |= set(re.findall(pola, kode))
        for nama in re.findall(r'ADD CONSTRAINT (\w+)', kode):
            constraint_per_berkas.append((os.path.basename(berkas), nama))
        for blok in re.findall(r'CREATE TABLE KLAIMNP\.\w+\s*\((.*?)\n\);', kode, re.S):
            per_jenis['kolom'] |= set(re.findall(r'^\s{2,}([A-Z][A-Z_0-9]*)\s+\S', blok, re.M))
    return per_jenis, constraint_per_berkas

def main():
    per_jenis, constraints = kumpulkan()
    semua = set()
    for v in per_jenis.values():
        semua |= v
    semua = set(x for x in semua if x)

    print('Pengenal di %s' % os.path.normpath(AKAR))
    for jenis, _ in POLA + [('kolom', None)]:
        n = len(per_jenis.get(jenis, ()))
        if not n:
            continue
        panjang = max(per_jenis[jenis], key=lambda s: len(s.encode('utf-8')))
        print('  %-12s %4d   terpanjang: %-28s %2d byte'
              % (jenis, n, panjang, len(panjang.encode('utf-8'))))
    print('  %-12s %4d' % ('UNIK', len(semua)))

    gagal = []

    lewat = sorted(x for x in semua if len(x.encode('utf-8')) > BATAS)
    print('\n1. Di atas %d byte : %d' % (BATAS, len(lewat)))
    for x in lewat:
        print('     ! %s (%d byte)' % (x, len(x.encode('utf-8'))))
        gagal.append(x)

    pola_sing = re.compile(r'(?:^|_)(%s)(?:_|$)' % '|'.join(SINGKATAN_LAMA))
    sisa = sorted(x for x in semua if pola_sing.search(x))
    print('2. Sisa singkatan : %d' % len(sisa))
    for x in sisa:
        print('     ! %s' % x)
        gagal.append(x)

    cacah = collections.Counter(n for _, n in constraints)
    ganda = sorted(n for n, k in cacah.items() if k > 1)
    print('3. Nama constraint dipakai dua kali : %d' % len(ganda))
    for n in ganda:
        print('     ! %s di %s' % (n, ', '.join(f for f, x in constraints if x == n)))
        gagal.append(n)

    if gagal:
        print('\nGAGAL. %d pelanggaran. Pembangkitan DDL berhenti di sini.' % len(gagal))
        return 1
    print('\nBERSIH. %d pengenal unik, tidak satu pun melanggar.' % len(semua))
    return 0

if __name__ == '__main__':
    sys.exit(main())
