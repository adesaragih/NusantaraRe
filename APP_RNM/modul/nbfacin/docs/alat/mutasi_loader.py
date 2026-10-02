"""Uji mutasi tiket 22 (loader.Flatten + generator). Tiap mutasi merusak SATU aturan;
`go test` paket loader (dan bangkit) wajib gagal. Bita asli dipulihkan persis.

    PYTHONIOENCODING=utf-8 py modul/nbfacin/docs/alat/mutasi_loader.py
"""
import pathlib
import subprocess
import sys

APP = pathlib.Path(__file__).resolve().parents[4]
L = APP / 'modul/nbfacin/backend/services/loader'
UJI = ['go', 'test', './modul/nbfacin/backend/services/loader/...']
F, A, B = L / 'flatten.go', L / 'aturan.go', L / 'bangkit/main.go'
AM = L / 'amandemen.go'
D = L / 'dokumen.go'
M = [
    ('V-28b sufiks Old tidak dibuang', A, '\tcase strings.HasSuffix(medan, "Old"):\n\t\treturn alasanSufiksOld\n', ''),
    ('V-19 OldData tidak dibuang', A, '\t"OldData":                "V-19', '\t"OldDataX":               "V-19'),
    ('V-28 Total* tidak dibuang', A, '\t"TotalTSIList":           "V-28', '\t"TotalTSIListX":          "V-28'),
    ('V-16 Aneka mendahului FIRE', F,
     '\tcase adaDiBawah(lokasi, func(s *simpul) bool {\n\t\treturn berisiLarik(s.medan("Property").medan("PropertyItemList"))\n\t}):\n\t\treturn "FIRE", nil\n\tcase adaKunciDiBawah(lokasi, "AnekaList"):\n\t\treturn "Aneka", nil\n',
     '\tcase adaKunciDiBawah(lokasi, "AnekaList"):\n\t\treturn "Aneka", nil\n\tcase adaDiBawah(lokasi, func(s *simpul) bool {\n\t\treturn berisiLarik(s.medan("Property").medan("PropertyItemList"))\n\t}):\n\t\treturn "FIRE", nil\n'),
    ('V-16 larik kosong dihitung ada', F, '&& len(s.anak) > 0 }', '&& len(s.anak) >= 0 }'),
    ('V-16 AnekaList dicari dari akar (ikut Fac Retro)', F, 'case adaKunciDiBawah(lokasi, "AnekaList"):', 'case adaKunciDiBawah(akar, "AnekaList"):'),
    ('V-16 langkah 5 menerima BusinessType apa pun yang dikenal', F, '\tif bt.teks != "Life" && bt.teks != "PA" {', '\tif false {'),
    ('K-069 warisan UNKNOWN tidak terhitung', F,
     '\t\t\t\t\tif v.Teks == kodeTakDiketahui {\n\t\t\t\t\t\tp.hasil.Diagnostik.MataUangTakDiketahui[ba.tabel+"."+k.nama]++\n\t\t\t\t\t}\n', ''),
    ('K-063 (c) pewarisan dimatikan', F, '\t\t\tif sumber, waris := warisMataUang[ba.tabel]; waris && k.nama == kolomKodeUang {',
     '\t\t\tif sumber, waris := warisMataUang[ba.tabel]; false && waris && k.nama == kolomKodeUang {'),
    ('K-069 sentinel tidak dipasang', F, '\t\t\tba.b.Kolom[k.nama] = Nilai{Teks: kodeTakDiketahui}\n', ''),
    ('BAHAN §7 koma desimal tidak dinormalkan', F, '\t\ts = strings.Replace(s, ",", ".", 1)\n', ''),
    ('BAHAN §7 koma desimal tidak terhitung', F, '\t\tp.hasil.Diagnostik.KomaDesimal[tabel+"."+k.nama]++\n', ''),
    ('angka > 38 digit tidak terhitung', F, '\tif digitBermakna(d) > 38 {', '\tif digitBermakna(d) > 99 {'),
    ('V-22 medan Policy sendiri ikut dilipat', A, '\t{"T_CURRENCYLIST", "Policy"}:      {dalam: map[string]string{"Payment": "Pay"}},',
     '\t{"T_CURRENCYLIST", "Policy"}:      {dalam: map[string]string{"Payment": "Pay"}, medanSendiri: true},'),
    ('V-22b awalan Pay hilang', A, '\t{"T_GENERAL_POLIS", "PolicyData"}: {dalam: map[string]string{"Payment": "Pay"}, medanSendiri: true},',
     '\t{"T_GENERAL_POLIS", "PolicyData"}: {dalam: map[string]string{"Payment": ""}, medanSendiri: true},'),
    ('V-24b CargoList/PolicyData tidak dilipat', A, '\t{"T_CARGOLIST", "PolicyData"}:     {medanSendiri: true},\n', ''),
    ('BAHAN §0 CURRENCY_CODE dari Name dimatikan', A, '\t"T_CURRENCYLIST": {"Name": "CURRENCY_CODE"},', '\t"T_CURRENCYLIST": {"NameX": "CURRENCY_CODE"},'),
    ('V-49 SOURCE_ID dimatikan', A, '\t"T_GENERAL_POLIS": "SOURCE_ID",', '\t"T_GENERAL_POLISX": "SOURCE_ID",'),
    ('V-27 cabang kosong tetap lahir', F, '\tif p.hasil.Diagnostik.Terpetakan > awalTerpetakan {', '\tif true {'),
    ('V-41 SEQ_NO mulai nol', F, 'seq, pos = i+1, posisiLarik(posisi, i)', 'seq, pos = i, posisiLarik(posisi, i)'),
    ('PARENT_TABLE berisi tabel sendiri', F, '\t\tb.Kolom[kolomTabelInduk] = Nilai{Teks: induk.tabel}', '\t\tb.Kolom[kolomTabelInduk] = Nilai{Teks: t}'),
    ('SRC_PATH kosong', F, '\t\tb.Kolom[kolomJalurSumber] = Nilai{Teks: jalur}', '\t\tb.Kolom[kolomJalurSumber] = Nilai{Teks: ""}'),
    ('medan tak terpetakan tidak terhitung', F, '\t\tp.hasil.Diagnostik.TakTerpetakan[kunci]++\n', ''),
    ('butir 68.1 teks menyimpang diurai sebagai angka', F, '\tif !k.angka() || kolomTeksMenyimpang[tabel+"."+k.nama] != "" {', '\tif !k.angka() {'),
    ('butir 68.3 kode dari Currency/Name dimatikan', A, '"T_COVERAGELIST": "T_CURRENCY", "T_ANEKALIST": "T_CURRENCY",', '"T_COVERAGELISTX": "T_CURRENCY", "T_ANEKALIST": "T_CURRENCY",'),
    # Butir 69: kembaran FR IKUT 68.3. Mutasi "perluas ke FR" putaran sebelumnya kini perilaku
    # sah; penggantinya mematikan kembaran FR. (Catatan putaran 68: menambah T_FR_COVERAGELIST ke
    # peta bool SAJA dulu LOLOS karena anak Currency FR adalah baris T_FR_CURRENCY.)
    ('butir 69 kembaran FR dimatikan', A, '"T_FR_COVERAGELIST": "T_FR_CURRENCY",', '"T_FR_COVERAGELISTX": "T_FR_CURRENCY",'),
    # Mutasi "anak tabel apa pun" (kodeDariCurrency[...] != "") DIKELUARKAN, putaran 69 LOLOS:
    # [terverifikasi] 02-10-2026 di skema_gen.go, anak keempat tabel kodeDariCurrency yang punya
    # kolom NAME hanyalah baris halaman Currency-nya (16 jalur, semuanya .../Currency) - tidak ada
    # dokumen sah yang membedakan. Kode tetap memeriksa tabel anak bila skema kelak berubah.
    ('ralat 69 >38 digit dihitung dari koefisien mentah', F, '\tif digitBermakna(d) > 38 {', '\tif d.NumDigits() > 38 {'),
    ('K-069 7b IsCedingConfirm tidak diselamatkan', F, 'selamat := p.selamatkan(ba, jalurAnakDari(jalur, k), v, medanDiselamatkan[k])', 'selamat := 0'),
    ('K-069 7b medan diselamatkan terhitung dua kali', F, 'if n := hitungDaun(v) - selamat; n > 0 {', 'if n := hitungDaun(v) + 0*selamat; n > 0 {'),
    ('butir 68.3 T_ANEKALIST dimatikan', A, '"T_ANEKALIST": "T_CURRENCY",', '"T_ANEKALISTX": "T_CURRENCY",'),
    ('butir 69 T_FR_ANEKALIST dimatikan', A, '"T_FR_ANEKALIST": "T_FR_CURRENCY",', '"T_FR_ANEKALISTX": "T_FR_CURRENCY",'),
    ('butir 68.3 kode berbeda menimpa diam-diam', F, '\t\t\tif v.Teks != kode {', '\t\t\tif false {'),
    ('ADR-0023 medan tak dikenal tidak masuk penampung', F,
     '\tp.hasil.Penampung = append(p.hasil.Penampung, MedanTakDikenal{Kunci: ba.b.Kunci, Jalur: jalurAtauAkar(jalur),', '\t_ = append(p.hasil.Penampung, MedanTakDikenal{Kunci: ba.b.Kunci, Jalur: jalurAtauAkar(jalur),'),
    ('ADR-0023 subpohon tidak masuk penampung', F, '\t\tp.tampungSubpohon(ba.b.Kunci, jalur, v)\n', ''),
    ('ADR-0023 catatTak tanpa nilai', F, 'Medan: medan, Nilai: teks, Penunjuk: penunjuk})', 'Medan: medan, Penunjuk: penunjuk})'),
    ('ADR-0023 penanda Penunjuk hilang', F, 'Medan: medan, Nilai: teks, Penunjuk: penunjuk})', 'Medan: medan, Nilai: teks})'),
    ('V-27 entri penampung baris batal tidak pindah', F, '\t\tp.hasil.Penampung[i].Kunci = induk.b.Kunci\n', ''),
    ('butir 68.3 Currency ganda saling timpa', F, 'if lama, sudah := anakCurrency[ba.induk.b]; sudah && lama != v.Teks {',
     'if lama, sudah := anakCurrency[ba.induk.b]; false && sudah && lama != v.Teks {'),
    ('ADR-0023 penampung tanpa nilai', F, 'MedanTakDikenal{Kunci: kunci, Jalur: jalur, Medan: nama, Nilai: a.teks}', 'MedanTakDikenal{Kunci: kunci, Jalur: jalur, Medan: nama}'),
    ('subpohon tak terpetakan tidak terhitung', F, '\t\tp.hasil.Diagnostik.TakTerpetakan[jalur+" :: *"] += n\n', ''),
    ('V-16 langkah 2 ikut cabang dibuang', F, '\t\t\tif _, dibuang := cabangDibuang[s.kunci[i]]; dibuang {',
     '\t\t\tif _, dibuang := cabangDibuang[s.kunci[i]]; false && dibuang {'),
    ('V-30 SEQ_NO opsi tanpa celah', F, 'posisi+"#"+strconv.Itoa(n), n)', 'posisi+"#"+strconv.Itoa(n), 1)'),
    # Cabang cadangan sebabDokumen ("dokumen tidak dapat diurai") TIDAK dimutasi: galat
    # json.Decoder.Token yang diketahui hanya SyntaxError / EOF / ErrUnexpectedEOF, dan
    # galat bentuk dokumen sendiri bertipe galatBentukDokumen - tak satu masukan pun yang
    # dapat disusun mencapainya. Putaran 02-10-2026 mencatatnya LOLOS; dikeluarkan dengan
    # alasan ini, bukan dilabeli "ekuivalen".
    ('pesan galat sintaks menyalin isi', F, '\t\treturn fmt.Sprintf("sintaks JSON tidak sah di bita %d", sintaks.Offset)',
     '\t\treturn err.Error()'),
    ('unsur larik bernilai tidak terhitung', F, '\t\tif s.jenis == simpulLarik || !metadata(s.kunci[i]) {', '\t\tif s.jenis != simpulLarik && !metadata(s.kunci[i]) {'),
    ('V-30 opsi tidak terhitung Terpetakan', F, '\tif ok {\n\t\tp.hasil.Diagnostik.Terpetakan++\n\t}\n\treturn err\n}', '\treturn err\n}'),
    ('ROW_UID tanpa IDPEGA', F, '"loader.Flatten\\x00" + idpega + "\\x00"', '"loader.Flatten\\x00" + "\\x00"'),
    ('V-48 medan kerja tidak ke T_WORK_POLIS', F, '\tif len(ks) == 0 && sasaran == "T_GENERAL_POLIS" {', '\tif len(ks) == 0 && sasaran == "T_GENERAL_POLISX" {'),
    ('V-17b ganti nama dimatikan', A, '\t{"VehicleList", "Occupation"}: "OccupationList",', '\t{"VehicleList", "OccupationX"}: "OccupationList",'),
    ('V-30 faktor tak dikenali', F, '\t\tif k == medanTandaCek || polaOpsi', '\t\tif k == "X" || polaOpsi'),
    ('kunci ganda diterima', D, '\t\t\t\tif ada[k] {', '\t\t\t\tif false {'),
    ('JENIS_WORK apa pun diterima', F, '\tif !ok || !jenisWork[j] || nomor == "" {', '\tif !ok || nomor == "" {'),
    ('kolom ganda menimpa', F, '\tif _, sudah := b.Kolom[k.nama]; sudah {', '\tif _, sudah := b.Kolom[k.nama]; false && sudah {'),
    ('asal kolom V-47 dianggap medan', A, '\t\tcase kolomFKV47[k.nama]:\n\t\t\treturn asalKosong, alasanFKV47\n', ''),
    ('P1 COVERAGE_INITIAL dicabut', AM, '"T_COVERAGELIST":    {{nama: "COVERAGE_INITIAL", tipe: "VARCHAR2(500)", medan: "CoverageInitial"}},', ''),
    ('P2-P3 jalur T_ADDITIONALSHIP dicabut', AM, '{jalur: "CargoList/Ship/AdditionalShip", tabel: "T_ADDITIONALSHIP", induk: "T_SHIP"},', ''),
    ('P2-P3 ADDITIONAL_SHIP_REF_ID dicabut', A, '"T_ADDITIONALSHIP": "ADDITIONAL_SHIP_REF_ID",', '"T_ADDITIONALSHIPX": "ADDITIONAL_SHIP_REF_ID",'),
    ('P5 CURRENCY_REF_ID dicabut', A, '"T_CURRENCYLIST":   "CURRENCY_REF_ID",', '"T_CURRENCYLISTX":   "CURRENCY_REF_ID",'),
    ('P6 POLICY_TSI dicabut', A, ', khusus: map[string]string{"TSI": "POLICY_TSI"}}', '}'),
    # Butir 71 mencabut P7: tiga mutasi P7 lama (V-50 / induk langsung / leluhur) tak berjangkar lagi.
    # Penggantinya menagih bahwa penunjuk TIDAK dibuang diam-diam (A67, A68):
    ('butir 71 penunjuk dibuang diam-diam', F, '\t\tp.hasil.Diagnostik.PenunjukBelumDikonversi[kunci]++', '\t\tp.hasil.Diagnostik.Dibuang[kunci]++'),
    ('butir 72 kolom INDEX_CARGO dicabut', AM, '{tabel: "T_COVERAGELIST", nama: "INDEX_CARGO", medan: "IndexCargo"},', ''),
    ('butir 72 penunjuk diurai sebagai angka', AM, 'const tipePenunjuk = "VARCHAR2(50)"', 'const tipePenunjuk = "NUMBER"'),
    ('pembaca xlsx menggeser sel kosong', B, '\t\t\ti := indeksKolom(c.R)\n', '\t\t\ti := len(row.C) - len(row.C) + func() int { n := 0; for _, x := range r { if x != "" { n++ } }; return n }()\n'),
]
asli = {}
for f in {F, A, B, D, AM}:
    asli[f] = f.read_bytes()
if subprocess.run(UJI, cwd=APP, capture_output=True).returncode:
    sys.exit('BASELINE MERAH')
n = 0
try:
    for nama, f, a, b in M:
        crlf = b'\r\n' in asli[f]
        teks = asli[f].decode('utf-8').replace('\r\n', '\n')
        if teks.count(a) != 1:
            sys.exit(f'JANGKAR BASI ({teks.count(a)}): {nama}')
        t = teks.replace(a, b)
        f.write_bytes((t.replace('\n', '\r\n') if crlf else t).encode('utf-8'))
        r = subprocess.run(UJI, cwd=APP, capture_output=True)
        n += r.returncode != 0
        print('TERTANGKAP' if r.returncode else 'LOLOS     ', nama)
        f.write_bytes(asli[f])
finally:
    for f, b in asli.items():
        f.write_bytes(b)
pulih = all(f.read_bytes() == b for f, b in asli.items())
print(f'HASIL: {n}/{len(M)} tertangkap; berkas asli pulih: {"ya" if pulih else "TIDAK"}')
