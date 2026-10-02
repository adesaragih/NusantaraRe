"""Uji mutasi layar tiket 21 (CoverageCargo.tsx). Bita asli dipulihkan persis.

    PYTHONIOENCODING=utf-8 py modul/nbfacin/docs/alat/mutasi_layar21.py
"""
import pathlib
import subprocess
import sys

APP = pathlib.Path(__file__).resolve().parents[4]
F = APP / 'modul/nbfacin/frontend/pages/CoverageCargo.tsx'
UJI = ['npx', 'vitest', 'run', 'modul/nbfacin/frontend/pages/CoverageCargo.test.ts']
M = [
    ('tukar urutan Rate dan LimitofLiability',
     '        <Field label={MEDAN.rate.label} value={rate} onChange={ubah(setRate)} />\n        <Field label={MEDAN.limitOfLiability.label} value="" onChange={tanpaUbah} readOnly />',
     '        <Field label={MEDAN.limitOfLiability.label} value="" onChange={tanpaUbah} readOnly />\n        <Field label={MEDAN.rate.label} value={rate} onChange={ubah(setRate)} />'),
    ('Min Premi dapat diisi',
     '<Field label={MEDAN.minPremi.label} value="" onChange={tanpaUbah} readOnly />',
     '<Field label={MEDAN.minPremi.label} value={tsi} onChange={ubah(setTsi)} />'),
    ('tombol Pega aktif',
     '<button type="button" className="btn btn--ghost btn--sm" disabled aria-describedby="nbfacin-belum-diport">',
     '<button type="button" className="btn btn--ghost btn--sm" aria-describedby="nbfacin-belum-diport">'),
    ('keterangan belum diport hanya tooltip', '<small>({TEKS.belumDiport})</small>', ''),
    ('kelompok layout digabung', '      </div>\n      <div className="form-grid">\n        <Field label={MEDAN.tsi.label}', '        <Field label={MEDAN.tsi.label}'),
]
bita = F.read_bytes()
crlf = b'\r\n' in bita
asli = bita.decode('utf-8').replace('\r\n', '\n')
if subprocess.run(UJI, cwd=APP, capture_output=True, shell=True).returncode:
    sys.exit('BASELINE MERAH')
n = 0
try:
    for nama, a, b in M:
        if asli.count(a) != 1:
            sys.exit(f'JANGKAR BASI: {nama}')
        t = asli.replace(a, b)
        F.write_bytes((t.replace('\n', '\r\n') if crlf else t).encode('utf-8'))
        r = subprocess.run(UJI, cwd=APP, capture_output=True, shell=True)
        n += r.returncode != 0
        print('TERTANGKAP' if r.returncode else 'LOLOS     ', nama)
finally:
    F.write_bytes(bita)
print(f'HASIL: {n}/{len(M)} tertangkap; berkas asli pulih: {"ya" if F.read_bytes() == bita else "TIDAK"}')
