"""Uji mutasi tiket R05 (daftar kasus renewal) - rnwfacin.

    PYTHONIOENCODING=utf-8 py modul/rnwfacin/docs/alat/mutasi_r05.py

Satu mutan = satu jangkar teks diganti di `backend/services/daftar.go`, lalu
`go vet` + `go test` paket layanan. Baseline wajib hijau, setiap jangkar wajib
muncul tepat sekali, berkas asli selalu dipulihkan. Log: `mutasi_r05.log`.
"""
import datetime
import pathlib
import subprocess
import sys

APP = pathlib.Path(__file__).resolve().parents[4]  # .../APP_RNM
DF = 'modul/rnwfacin/backend/services/daftar.go'
PK = ['./modul/rnwfacin/backend/services/']
M = [
    ('tanpa filter selesai', 'b.StatusKerja != models.StatusKerjaSelesai && b.PembuatID', 'b.PembuatID', ''),
    ('tanpa filter ditolak', ' && b.StatusKerja != models.StatusKerjaDitolak {', ' {', ''),
    ('tanpa filter pembuat', 'b.PembuatID == sr.PembuatID &&', '', ''),
    ('pembuat tanpa peka huruf', 'b.PembuatID == sr.PembuatID &&', 'strings.EqualFold(b.PembuatID, sr.PembuatID) &&', ''),
    ('tanpa filter team', '\t\t\tb.TeamGroup == sr.TeamGroup && ', '\t\t\t', ''),
    ('urut naik', 'return lolos[i].DibuatPada.After(lolos[j].DibuatPada)', 'return lolos[i].DibuatPada.Before(lolos[j].DibuatPada)', ''),
    ('pyID naik', 'return lolos[i].IDKasus > lolos[j].IDKasus', 'return lolos[i].IDKasus < lolos[j].IDKasus', ''),
    ('batas 499', 'const batasBaris = 500', 'const batasBaris = 499', ''),
    ('terpotong >=', 'if len(lolos) > batasBaris {', 'if len(lolos) >= batasBaris {', ''),
    ('pembuat kosong diterima', '\tif sr.PembuatID == "" {', '\tif false {', ''),
    ('team kosong diterima', '\tif sr.TeamGroup == "" {', '\tif false {', ''),
    ('status kosong diterima', '\t\tif b.StatusKerja == "" {', '\t\tif false {', ''),
    ('pesan pembuat memuat login', 'fmt.Errorf("%w: pembuat", ErrSaringanKosong)', 'fmt.Errorf("%w: pembuat %s", ErrSaringanKosong, sr.PembuatID)',
     'ekuivalen: di cabang ini PembuatID kosong, tidak ada login yang bocor'),
    ('pesan team memuat login', 'fmt.Errorf("%w: team group", ErrSaringanKosong)', 'fmt.Errorf("%w: team group %s", ErrSaringanKosong, sr.PembuatID)', ''),
]


def jalan(args):
    return subprocess.run(args, cwd=APP, capture_output=True, text=True, encoding='utf-8', errors='replace')


def main():
    baris = [f'# Uji mutasi R05 - {datetime.datetime.now():%Y-%m-%d %H:%M}', f'# {len(M)} mutan']
    if jalan(['go', 'test', '-count=1'] + PK).returncode:
        sys.exit('BASELINE MERAH - instrumen tidak sah')
    baris.append('baseline: hijau')
    p = APP / DF
    asli = p.read_text(encoding='utf-8')
    for nama, a, _, _ in M:
        if asli.count(a) != 1:
            sys.exit(f'JANGKAR BASI: {nama} - muncul {asli.count(a)} kali')
    tangkap, lolos, takKompil = 0, [], []
    for nama, a, b, ekuivalen in M:
        try:
            baru = asli.replace(a, b)
            if 'strings.' in b:
                baru = baru.replace('"sort"', '"sort"\n\t"strings"', 1)
            p.write_text(baru, encoding='utf-8', newline='')
            if jalan(['go', 'vet'] + PK).returncode:
                takKompil.append(nama)
                baris.append(f'TAK TERKOMPILASI  {nama}')
                continue
            if jalan(['go', 'test', '-count=1'] + PK).returncode:
                tangkap += 1
                baris.append(f'TERTANGKAP        {nama}')
            else:
                lolos.append(nama)
                baris.append(f'LOLOS             {nama}' + (f'  [{ekuivalen}]' if ekuivalen else ''))
        finally:
            p.write_text(asli, encoding='utf-8', newline='')
    baris.append(f'HASIL: {tangkap}/{len(M) - len(takKompil)} tertangkap; lolos {lolos}; tak terkompilasi {takKompil}')
    baris.append('berkas asli pulih: ' + ('ya' if p.read_text(encoding='utf-8') == asli else 'TIDAK'))
    teks = '\n'.join(baris) + '\n'
    pathlib.Path(__file__).with_suffix('.log').write_text(teks, encoding='utf-8', newline='\n')
    print(teks)


if __name__ == '__main__':
    main()
