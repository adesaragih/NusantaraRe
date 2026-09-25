# -*- coding: utf-8 -*-
"""
panen-kolom-dari-penulis.py — 24 September 2026

KENAPA PERKAKAS INI ADA
    `TREATY_IN` dan `TREATY_IN_EDM` adalah dua tabel yang paling banyak dirujuk
    berkas permintaan DBA, dan justru keduanya yang DDL-nya tidak kita punya.
    Uji A-0 menutupnya -- tetapi hanya SESUDAH DBA menjalankannya, dan bila ia
    kembali dengan puluhan nama salah, berkasnya diputar lagi. Satu putaran
    diukur dalam hari.

    Ada kebenaran dasar KEDUA yang sudah ada dan belum dipakai: PENULISNYA.
    Badan prosedur di PROCEDURE/*.txt memuat INSERT dan UPDATE yang menyebutkan
    kolomnya SATU PER SATU. Kolom yang ditulis sistem lama setiap hari PASTI ADA.

BATAS KEMAMPUAN PERKAKAS INI, dinyatakan di dalam dirinya sendiri:
    Ia dapat menyatakan sebuah kolom ADA -- karena ada yang menulisinya.
    Ia TIDAK dapat menyatakan sebuah kolom TIDAK ADA: kolom yang tidak pernah
    ditulis prosedur mungkin tetap ada di tabel, diisi jalur lain atau tidak
    diisi sama sekali. Yang tidak cocok karena itu ditandai PALING PATUT
    DICURIGAI, bukan salah.
    Uji A-0 tetap yang menutup sisanya, dan ia dijalankan DBA.
"""
import collections
import io
import os
import re

AKAR = r'D:\XML_NURE\_migration-docs\treaty-in'
SQL = os.path.join(AKAR, 'PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql')
PROSEDUR = os.path.join(AKAR, 'PROCEDURE')
TABEL_DDL = os.path.join(AKAR, 'Table')

RE_INSERT = re.compile(
    r'INSERT\s+INTO\s+(?:POOLDATA\.)?([A-Z_0-9]+)\s*\((.*?)\)\s*VALUES', re.I | re.S)
RE_UPDATE = re.compile(
    r'UPDATE\s+(?:POOLDATA\.)?([A-Z_0-9]+)\s+SET\s+(.*?)(?:WHERE|\Z)', re.I | re.S)
RE_SELECT_FROM = re.compile(
    r'SELECT\s+(.*?)\s+FROM\s+(?:POOLDATA\.)?([A-Z_0-9]+)', re.I | re.S)


def panen():
    kol = collections.defaultdict(set)
    asal = collections.defaultdict(set)
    if not os.path.isdir(PROSEDUR):
        return kol, asal
    for f in sorted(os.listdir(PROSEDUR)):
        s = io.open(os.path.join(PROSEDUR, f), encoding='utf-8', errors='replace').read()
        for t, daftar in RE_INSERT.findall(s):
            for k in re.findall(r'\b([A-Z][A-Z_0-9]{1,})\b', daftar):
                kol[t.upper()].add(k)
                asal[t.upper()].add(f + ' (INSERT)')
        for t, daftar in RE_UPDATE.findall(s):
            for k in re.findall(r'\b([A-Z][A-Z_0-9]{1,})\s*=', daftar):
                kol[t.upper()].add(k)
                asal[t.upper()].add(f + ' (UPDATE)')
    return kol, asal


def dari_ddl():
    peta = {}
    if not os.path.isdir(TABEL_DDL):
        return peta
    for f in os.listdir(TABEL_DDL):
        s = io.open(os.path.join(TABEL_DDL, f), encoding='utf-8', errors='replace').read()
        m = re.search(r'CREATE TABLE\s+"([^"]+)"\."([^"]+)"', s)
        if m:
            peta[m.group(2).upper()] = set(re.findall(
                r'"\s*([A-Z_0-9]+)\s*"\s+(?:VARCHAR2|NUMBER|DATE|CLOB|TIMESTAMP)', s))
    return peta


def kolom_dipakai_sql():
    """Kolom nyata yang dirujuk berkas permintaan, per tabel."""
    s = io.open(SQL, encoding='utf-8').read()
    pakai = collections.defaultdict(set)
    RE_TA = re.compile(r'\b(?:FROM|JOIN)\s+POOLDATA\.([A-Z_0-9]+)\s+([a-zA-Z]\w*)\b', re.I)
    RE_TP = re.compile(r'\b(?:FROM|JOIN)\s+POOLDATA\.([A-Z_0-9]+)\s*\n', re.I)
    RE_Q = re.compile(r'\b([a-zA-Z]\w*)\.([A-Z][A-Z_0-9]{1,})\b')
    RE_MAYA = re.compile(r'\b([A-Z][A-Z_0-9]{1,})\s+(?:VARCHAR2|NUMBER|DATE|CLOB)')
    RE_JSONALIAS = re.compile(r'\)\s*\)\s*([a-zA-Z]\w*)\b')
    KATA = {'SELECT', 'FROM', 'WHERE', 'AND', 'OR', 'NOT', 'IN', 'IS', 'NULL', 'AS',
            'ON', 'CASE', 'WHEN', 'THEN', 'ELSE', 'END', 'COUNT', 'SUM', 'MIN',
            'MAX', 'AVG', 'GROUP', 'ORDER', 'BY', 'PATH', 'COLUMNS', 'NESTED'}
    for bagian in s.split(';'):
        if 'SELECT' not in bagian.upper():
            continue
        alias = {a: t.upper() for t, a in RE_TA.findall(bagian)}
        maya_a = set(RE_JSONALIAS.findall(bagian))
        maya_k = set(RE_MAYA.findall(bagian))
        polos = [t.upper() for t in RE_TP.findall(bagian)]
        for a, k in RE_Q.findall(bagian):
            if k in KATA or k in maya_k or a in maya_a or a not in alias:
                continue
            pakai[alias[a]].add(k)
        # tabel tanpa alias: kolom telanjang tidak dapat dipetakan dengan aman
        # kecuali hanya ada SATU tabel nyata di pernyataan itu
        if len(polos) == 1 and not alias:
            for k in re.findall(r'\b([A-Z][A-Z_0-9]{2,})\b', bagian):
                if k in KATA or k in maya_k:
                    continue
                pakai[polos[0]].add(k)
    return pakai


def main():
    kol, asal = panen()
    ddl = dari_ddl()
    pakai = kolom_dipakai_sql()

    print('=== KOLOM YANG DIPANEN DARI PENULISNYA (prosedur) ===')
    for t in sorted(kol):
        print('   %-22s %3d kolom   sumber: %s'
              % (t, len(kol[t]), ', '.join(sorted(asal[t]))))
    print()
    for t in ('TREATY_IN', 'TREATY_IN_EDM'):
        if t in kol:
            print('   %s:\n      %s' % (t, ', '.join(sorted(kol[t]))))
    print()
    print('=== PEMERIKSAAN BERKAS PERMINTAAN ===')
    total_ok = total_curiga = 0
    for t in sorted(pakai):
        dikenal = set(ddl.get(t, set())) | set(kol.get(t, set()))
        if not dikenal:
            print('   %-22s tidak ada sumber apa pun -> seluruh %d rujukan BELUM TERPERIKSA'
                  % (t, len(pakai[t])))
            continue
        ok = pakai[t] & dikenal
        curiga = pakai[t] - dikenal
        total_ok += len(ok)
        total_curiga += len(curiga)
        print('   %-22s terverifikasi %2d   PATUT DICURIGAI %2d  %s'
              % (t, len(ok), len(curiga), ', '.join(sorted(curiga)) if curiga else ''))
    print()
    print('   terverifikasi tanpa DBA : %d' % total_ok)
    print('   patut dicurigai         : %d' % total_curiga)
    print()
    print('   CATATAN BATAS: "patut dicurigai" BUKAN "salah". Kolom yang tidak')
    print('   pernah ditulis prosedur mungkin tetap ada. Uji A-0 yang menutupnya.')


if __name__ == '__main__':
    main()
