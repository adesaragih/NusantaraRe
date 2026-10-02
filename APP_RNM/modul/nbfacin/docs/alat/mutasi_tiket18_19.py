"""Uji mutasi tiket 18 (rumus premi lini lain), 19 (pembayaran MARINE), dan penjaga
rekonsiliasi EDM (A44) - nbfacin.

Jalankan dari mana saja (Windows, pemanggil `py`, CLAUDE.md §4a):

    PYTHONIOENCODING=utf-8 py modul/nbfacin/docs/alat/mutasi_tiket18_19.py

Cara kerja: setiap mutan mengganti SATU jangkar teks di satu berkas Go, lalu
`go vet` + `go test` paket terdampak; mutan "tertangkap" bila test gagal. Berkas
asli selalu dipulihkan (`finally`). Hasil ditulis ke `mutasi_tiket18_19.log` di
folder ini.

Pengaman instrumen (aturan "uji instrumennya dulu"):
- baseline: seluruh paket terdampak wajib hijau SEBELUM mutan pertama;
- setiap jangkar wajib muncul TEPAT SEKALI di kode sekarang - bila tidak, skrip
  berhenti dengan galat (jangkar basi tidak boleh dilewati diam-diam);
- mutan yang tidak terkompilasi dilaporkan terpisah, tidak dihitung.
"""
import datetime
import pathlib
import subprocess
import sys

APP = pathlib.Path(__file__).resolve().parents[4]  # .../APP_RNM
SRV = 'modul/nbfacin/backend/services/'
LL, PM, TG = SRV + 'premium/lini_lain.go', SRV + 'premium/premium.go', SRV + 'acceptance/tangga.go'
RK, BY = SRV + 'rekonsiliasi/rekonsiliasi.go', SRV + 'pembayaran/marine.go'
TT = SRV + 'premium/total.go'
PK = {
    LL: ['./' + SRV + 'premium/', './' + SRV + 'rekonsiliasi/'],
    PM: ['./' + SRV + 'premium/', './' + SRV + 'rekonsiliasi/'],
    TT: ['./' + SRV + 'premium/', './' + SRV + 'rekonsiliasi/'],
    TG: ['./' + SRV + 'acceptance/'],
    RK: ['./' + SRV + 'rekonsiliasi/'],
    BY: ['./' + SRV + 'pembayaran/', './' + SRV + 'rekonsiliasi/'],
}

# (berkas, nama, jangkar, pengganti, catatan ekuivalen bila diketahui)
M = [
    # --- FIRE (CountPremi_ACT)
    (LL, 'FIRE NetRate "0" dipakai', 'if in.NetRate != "0" && in.NetRate != "" {', 'if in.NetRate != "" {', ''),
    (LL, 'FIRE NetRate diabaikan', '\t\trate = netRate', '\t\t_ = netRate', ''),
    (LL, 'FIRE NetRate 0.0 diterima', '\t\tif netRate.IsZero() {', '\t\tif false {', ''),
    (LL, 'FIRE LossLimit teks → angka', '\tif teks == "" || teks == "0" {', '\tif teks == "" || d.IsZero() {', ''),
    (LL, 'FIRE LossLimit 0.0 diterima', '\tif d.IsZero() {\n\t\treturn nil, fmt.Errorf("%w: LossLimit %q"', '\tif false {\n\t\treturn nil, fmt.Errorf("%w: LossLimit %q"', ''),
    (LL, 'FIRE LossLimit tanpa bawaan', '\tlossLimit, err := lossLimitFire(in.LossLimit, lossLimit)', '\t_, err := lossLimitFire(in.LossLimit, lossLimit)', ''),
    (LL, 'FIRE FirstScale diabaikan', 'faktor = append(faktor, rasioPersen(firstScale))', '_ = firstScale', ''),
    (LL, 'FIRE PctAdjustment diabaikan', 'faktor = append(faktor, rasioPersen(adj))', '_ = adj', ''),
    (LL, 'FIRE presisi 4', 'Currency: in.MataUang}, 20, faktor...)', 'Currency: in.MataUang}, 4, faktor...)', ''),
    (LL, 'FIRE basis 3/4 ditolak', 'case "1", "3", "4":', 'case "1":', ''),
    (LL, 'FIRE basis 5 diterima', '\tcase "5":\n\t\treturn uang.Money{}, fmt.Errorf("%w: %v", ErrBentukBelumDiport, errLayering)', '\tcase "5":', ''),
    (LL, 'FIRE basis tak dikenal = belum diport', 'fmt.Errorf("%w: %q", ErrCoverageBasisTakDikenal, in.CoverageBasis)', 'fmt.Errorf("%w: %q", ErrBentukBelumDiport, in.CoverageBasis)', ''),
    # --- ANEKA / GOLF (CountPremiCoverageAneka, FillPremiGolf)
    (LL, 'ANEKA metode 3 = pro-rata', 'p = rasioProRata(apd.New(100, 0))', 'p = rasioProRata(proRata)', ''),
    (LL, 'ANEKA metode 2 = pro-rata', 'p = rasioPeriodePendek(periode)', 'p = rasioProRata(proRata)', ''),
    (LL, 'ANEKA 14.1 berlaku untuk GOLF', 'if in.LiniBisnis != LiniGolf && (proRata == nil', 'if (proRata == nil', ''),
    (LL, 'ANEKA 14.1 hilang', 'if in.LiniBisnis != LiniGolf && (proRata == nil || proRata.IsZero()) {', 'if false {', ''),
    (LL, 'ANEKA indemnity tanpa MBD', 'if in.MBD && in.IndemnityPercentage != "" {', 'if in.IndemnityPercentage != "" {', ''),
    (LL, 'ANEKA MBD tanpa cek indemnity', 'if in.MBD && in.IndemnityPercentage != "" {', 'if in.MBD {', ''),
    (LL, 'ANEKA loading dibuang', 'return kaliLossLimitPersen(jumlah, lossLimitBawaan(lossLimit))', '_ = jumlah\n\treturn kaliLossLimitPersen(pokok, lossLimitBawaan(lossLimit))', ''),
    (LL, 'ANEKA pokok presisi 8', 'pokok, err := premiKomposit(tsi, 4, dasar...)', 'pokok, err := premiKomposit(tsi, 8, dasar...)', ''),
    (LL, 'ANEKA LossLimit diabaikan', 'return kaliLossLimitPersen(jumlah, lossLimitBawaan(lossLimit))', 'return kaliLossLimitPersen(jumlah, apd.New(100, 0))', ''),
    (LL, 'ANEKA LossLimit 0.0 bukan bawaan', '\tif d == nil || d.IsZero() {', '\tif d == nil {', ''),
    (LL, 'skala pembilang tidak dipulihkan', 'hasil.Exponent > ideal {', 'false {', ''),
    (LL, 'tanpa Reduce', '\thasil.Reduce(hasil)\n', '', ''),
    # --- MARINE (CountGPWMarinePAMbu_Act)
    (LL, 'MARINE master policy diabaikan', '\tif in.MasterPolicy {', '\tif false {', ''),
    (LL, 'MARINE presisi 2', 'Currency: in.MataUang}, 4, rasioRate(LiniMarineCargo, rate))', 'Currency: in.MataUang}, 2, rasioRate(LiniMarineCargo, rate))', ''),
    # --- pembulatan (A37)
    (PM, 'premi setengah-genap', '\tctx.Rounding = apd.RoundHalfUp\n\thasil := new(apd.Decimal)\n\tif _, err := ctx.Quantize(hasil, d, -desimal)', '\tctx.Rounding = apd.RoundHalfEven\n\thasil := new(apd.Decimal)\n\tif _, err := ctx.Quantize(hasil, d, -desimal)', ''),
    (TG, 'bulatNol setengah-genap', '\tctx.Rounding = apd.RoundHalfUp\n\thasil := new(apd.Decimal)\n\tif _, err := ctx.Quantize(hasil, x, 0)', '\tctx.Rounding = apd.RoundHalfEven\n\thasil := new(apd.Decimal)\n\tif _, err := ctx.Quantize(hasil, x, 0)', ''),
    # --- rekonsiliasi (tiket 16/18, A44, tiket 19)
    (RK, 'rekon: master policy tidak pernah', 'MasterPolicy: teksDi(akar, "QuotationData", "PolicyType") == "1" && teksDi(akar, "QuotationData", "IsMOP") == "MOP"}', 'MasterPolicy: false}', ''),
    (RK, 'rekon: LossLimit dari kunci lain', 'LossLimit: teksDi(c, "LostLimit")', 'LossLimit: teksDi(c, "LossLimit")', ''),
    # Dua mutan berikut dulu LOLOS karena data nyata tidak menjangkau cabangnya - itu
    # CELAH TES, bukan mutan ekuivalen (ralat 01-10-2026). Ditutup
    # TestBerkasKasusCabangBuatan (masukan buatan).
    (RK, 'rekon: IsMBD selalu false', 'if mbd, err = rules.Eval("IsMBD", k); err != nil {', 'if _, err = rules.Eval("IsMBD", k); err != nil {', ''),
    (RK, 'rekon: pro-rata dari coverage', 'ProRatePercent: teksDi(akar, "ProRatePercent")', 'ProRatePercent: teksDi(c, "ProRatePercent")', ''),
    (RK, 'rekon: GOLF memakai jalur ANEKA', 'premium.LiniGolf:        {"LocationList", "Property", "RiskLocation", "AnekaList", "CoverageList"},', 'premium.LiniGolf:        {"LocationList", "Property", "RiskLocation", "OccupationList", "AnekaList", "CoverageList"},', ''),
    (RK, 'rekon: ANEKA tanpa OccupationList', 'premium.LiniAneka:       {"LocationList", "Property", "RiskLocation", "OccupationList", "AnekaList", "CoverageList"},', 'premium.LiniAneka:       {"LocationList", "Property", "RiskLocation", "AnekaList", "CoverageList"},', ''),
    (RK, 'rekon: catatan EDM hilang', '\tif edm {\n\t\t// CountGrossPremi_Act langkah 5', '\tif false {\n\t\t// CountGrossPremi_Act langkah 5', ''),
    (RK, 'rekon: EDM lini lain dibandingkan', 'if edm && l.Lini != premium.LiniFire {', 'if false {', ''),
    (RK, 'rekon: beda EDM tidak dicek', '\t\t\tif edm {\n\t\t\t\tif alasan := bedaEDM(c)', '\t\t\tif false {\n\t\t\t\tif alasan := bedaEDM(c)', ''),
    (RK, 'rekon: tanda [dugaan] TSI item hilang', 'if sama, err := samaNilai(teksDi(c.m, "TSI"), teksDi(c.induk, "TSIObjectItem")); err != nil || !sama {', 'if false {', ''),
    (RK, 'rekon: TSI item jadi galat lagi', '\t\t\t\th.Catatan = strings.Trim(h.Catatan+" | "+catatanTSIItem(c), " |")', '\t\t\t\tif catatanTSIItem(c) != "" {\n\t\t\t\t\th.Status = Galat\n\t\t\t\t}', ''),
    (RK, 'rekon: FlagDelete item diabaikan', 'if teksDi(c.m, "FlagDelete") == "1" && teksDi(c.induk, "FlagDelete") != "1" {', 'if teksDi(c.m, "FlagDelete") == "1" {', ''),
    (RK, 'rekon: baris pembayaran hilang', '\thasil = append(hasil, barisPembayaran(nama, kasusJSON{akar})...)\n', '', ''),
    (RK, 'rekon: pembayaran dinyatakan cocok', 'Langkah: langkahPembayaran, Status: BelumTercakup,', 'Langkah: langkahPembayaran, Status: Cocok,', ''),
    # --- total FIRE per item/lokasi (tiket 07, SumTotalTSIPremiGross_Act langkah 1)
    (TT, 'total: Rate ×100', 'ctx.Mul(rate, keAtas, apd.New(1000, 0))', 'ctx.Mul(rate, keAtas, apd.New(100, 0))', ''),
    (TT, 'total: seri @divide diterima', '\t\tif genap.Cmp(keAtas) != 0 {', '\t\tif false {', ''),
    (TT, 'total: TSI 0 tidak dijaga', '\tif tsi.IsZero() {\n\t\treturn apd.New(0, 0), nil', '\tif false {\n\t\treturn apd.New(0, 0), nil', ''),
    (TT, 'total: entri baru tanpa premi', 'TotalMataUang{TSI: tsi, Premium: premi}', 'TotalMataUang{TSI: tsi, Premium: uang.Money{Amount: apd.New(0, 0), Currency: it.MataUang}}', ''),
    (TT, 'total: premi entri lama tak ditambah', 'if e.Premium, err = premi.Add(e.Premium); err != nil {', 'if _, err = premi.Add(e.Premium); err != nil {', ''),
    (TT, 'total: mata uang kosong diterima', '\t\tif it.MataUang == "" {', '\t\tif false {', ''),
    (RK, 'rekon: TSI total dibanding teks', '\t\t\tif sama, err := samaNilai(tsiLama, tsiBaru); err == nil && sama {', '\t\t\tif false {', ''),
    (RK, 'rekon: entri LewatDouble dibandingkan', '\t\t\tif e.LewatDouble {\n\t\t\t\thasil = append(hasil', '\t\t\tif false {\n\t\t\t\thasil = append(hasil', ''),
    (TT, 'total: LewatDouble tak ditandai', '\t\t\te.LewatDouble = true\n', '', ''),
    (RK, 'rekon: catatan TSI item dibanding teks', 'if sama, err := samaNilai(teksDi(c.m, "TSI"), teksDi(c.induk, "TSIObjectItem")); err != nil || !sama {', 'if teksDi(c.m, "TSI") != teksDi(c.induk, "TSIObjectItem") {', ''),
    # --- pembayaran.Marine (tiket 19)
    (BY, 'bayar: tanpa gerbang marine', 'if err != nil || !marine {', 'if _ = marine; err != nil {', ''),
    (BY, 'bayar: cabang EDM dihitung', '\tif edm {\n\t\treturn Hasil{}, ErrCabangEDM', '\tif _ = edm; false {\n\t\treturn Hasil{}, ErrCabangEDM', ''),
    (BY, 'bayar: tanpa cabang dihitung', '\tif !nb {\n\t\treturn Hasil{}, ErrTanpaCabang', '\tif _ = nb; false {\n\t\treturn Hasil{}, ErrTanpaCabang', ''),
    (BY, 'bayar: diskon tidak dijumlah', 'if hasil.Diskon, err = hasil.Diskon.Add(diskon); err != nil {', 'if _, err = hasil.Diskon.Add(diskon); err != nil {', ''),
    (BY, 'bayar: premi tidak dijumlah', 'if hasil.Premium, err = hasil.Premium.Add(premi); err != nil {', 'if _, err = hasil.Premium.Add(premi); err != nil {', ''),
    (BY, 'bayar: kosong diterima', '\tif teks == "" {\n\t\treturn uang.Money{}, fmt.Errorf("%w: kosong", ErrMasukan)', '\tif false {\n\t\treturn uang.Money{}, fmt.Errorf("%w: kosong", ErrMasukan)',
     'ekuivalen: utils.ParseDecimal("") juga gagal, masukan kosong tetap ditolak lewat ErrMasukan'),
    (BY, 'bayar: mata uang kosong diterima (A46)', '\tif mataUang == "" {', '\tif false {', ''),
    (BY, 'bayar: tanpa coverage diterima', '\tif !hasil.Berlaku {\n\t\treturn Hasil{}, fmt.Errorf("%w: tanpa coverage kargo"', '\tif false {\n\t\treturn Hasil{}, fmt.Errorf("%w: tanpa coverage kargo"', ''),
]


def jalan(args):
    return subprocess.run(args, cwd=APP, capture_output=True, text=True, encoding='utf-8', errors='replace')


CRLF, LF = '\r\n', '\n'


def main():
    baris = [f'# Uji mutasi tiket 18/19 + A44 - {datetime.datetime.now():%Y-%m-%d %H:%M}', f'# {len(M)} mutan']
    semua = sorted({p for ps in PK.values() for p in ps})
    dasar = jalan(['go', 'test', '-count=1'] + semua)
    if dasar.returncode:
        print(dasar.stdout[-2000:], dasar.stderr[-2000:])
        sys.exit('BASELINE MERAH - instrumen tidak sah, tidak ada mutan dijalankan')
    baris.append('baseline: hijau (' + ' '.join(semua) + ')')
    # Bita asli disimpan persis (CRLF/LF tidak diubah); jangkar dicocokkan pada teks
    # berakhir-baris LF, lalu mutan ditulis dengan akhir baris asli berkas itu.
    bita = {f: (APP / f).read_bytes() for f in PK}
    asli = {f: b.decode('utf-8').replace(CRLF, LF) for f, b in bita.items()}
    crlf = {f: CRLF.encode() in b for f, b in bita.items()}
    for f, nama, a, _, _ in M:
        if asli[f].count(a) != 1:
            sys.exit(f'JANGKAR BASI: {nama} - muncul {asli[f].count(a)} kali di {f}')
    tangkap, lolos, takKompil = 0, [], []
    for f, nama, a, b, ekuivalen in M:
        p = APP / f
        try:
            mutan = asli[f].replace(a, b)
            p.write_bytes((mutan.replace(LF, CRLF) if crlf[f] else mutan).encode('utf-8'))
            v = jalan(['go', 'vet'] + PK[f])
            if v.returncode:
                takKompil.append(nama)
                baris.append(f'TAK TERKOMPILASI  {nama}: {v.stderr.strip().splitlines()[-1][:160]}')
                continue
            r = jalan(['go', 'test', '-count=1'] + PK[f])
            if r.returncode:
                tangkap += 1
                baris.append(f'TERTANGKAP        {nama}')
            else:
                lolos.append(nama)
                baris.append(f'LOLOS             {nama}' + (f'  [{ekuivalen}]' if ekuivalen else ''))
        finally:
            p.write_bytes(bita[f])
    dihitung = len(M) - len(takKompil)
    baris.append(f'HASIL: {tangkap}/{dihitung} tertangkap; lolos {len(lolos)}: {lolos}; tak terkompilasi {len(takKompil)}: {takKompil}')
    pulih = all((APP / f).read_bytes() == bita[f] for f in PK)  # bita persis, bukan teks
    baris.append('berkas asli pulih: ' + ('ya' if pulih else 'TIDAK'))
    teks = '\n'.join(baris) + '\n'
    (pathlib.Path(__file__).with_suffix('.log')).write_text(teks, encoding='utf-8', newline='\n')
    print(teks)


if __name__ == '__main__':
    main()
