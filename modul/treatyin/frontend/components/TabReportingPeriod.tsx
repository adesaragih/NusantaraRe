// Tab **Reporting Period** — `Section/TreatyInTabsProportional.xml` @26514.
//
// ⭐ 6 Oktober 2026 — tombol `Apply` MENGHITUNG, seperti di Pega.
//
// Tombolnya (@254615) menjalankan `Activity/TreatyInSetReport.xml`
// (`startdate = TreatyIn.ReportingStart`, `autocalculate = true`). Rumusnya
// hidup di services (`periode_pelaporan.go`) dan diukur ulang terhadap 4.548
// baris yang Pega simpan (99,5% cocok). Layar ini hanya mengirim isian dan
// menampilkan jawabannya — nol rumus di sini.
//
// ⚠️ Medan kepala (Start, End, Period, …) KOSONG saat tab dibuka: nilai
// tersimpannya hanya ada di `JSONDATA`, yang dilarang dibaca layar Treaty In
// (keputusan pemilik proses 6 Oktober 2026), dan `T_TREATY_REVISION` tidak
// punya kolom `Reporting*`. Lihat `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §17.

import { useState } from 'react'

import { Field, Kosong, Panel, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { hitungPeriodePelaporan, type BarisPeriodeWarisan, type OpsiPilihan } from '../api'
import { REPORTING_PERIOD } from '../labels'
import type { ModeForm } from '../mode'
import TanggalRedup from './TanggalRedup'

/** Pesan satu medan dari Activity (`Property-Set-Messages`), apa adanya. */
function Pesan({ teks }: { teks: string | undefined }) {
  if (teks === undefined || teks === '') return null
  return (
    <span className="trin__galat" role="alert">
      {teks}
    </span>
  )
}

export default function TabReportingPeriod({
  baris: barisWarisan,
  opsiPeriode = [],
  mode = 'lihat',
}: {
  baris: readonly BarisPeriodeWarisan[]
  opsiPeriode?: readonly OpsiPilihan[]
  mode?: ModeForm
}) {
  const [mulai, setMulai] = useState('')
  const [akhir, setAkhir] = useState('')
  const [periode, setPeriode] = useState('')
  const [interval, setSelang] = useState('')
  const [penyerahan, setPenyerahan] = useState('')
  const [konfirmasi, setKonfirmasi] = useState('')
  const [pelunasan, setPelunasan] = useState('')
  const [galat, setGalat] = useState<Record<string, string>>({})
  const [gagal, setGagal] = useState('')
  const [baris, setBaris] = useState<BarisPeriodeWarisan[]>(() => [...barisWarisan])
  const bisaUbah = mode === 'ubah'
  // `TreatyIn.ReportingPeriod='other'` @144723 — Interval dan `, and Interval`.
  const lain = periode === 'other'

  function terapkan() {
    setGagal('')
    hitungPeriodePelaporan({ mulai, akhir, periode, interval, penyerahan, konfirmasi, pelunasan })
      .then((h) => {
        setGalat(h.galat)
        // ⛔ Langkah 2 Activity: daftar DIBUANG lalu disusun ulang — tetapi
        // isian kosong KELUAR di langkah 4–7 SEBELUM ada baris baru. Daftar
        // kosong hanya bila jawabannya nol pesan.
        if (Object.keys(h.galat).length === 0) setBaris(h.baris)
        else setBaris([])
      })
      .catch((e: unknown) => {
        setGagal(e instanceof Error ? e.message : String(e))
      })
  }

  return (
    <Panel judul={REPORTING_PERIOD.judul}>
      <div className="form-grid">
        <div>
          <TanggalRedup label={REPORTING_PERIOD.mulai} value={mulai} onChange={setMulai} />
          <Pesan teks={galat.mulai} />
        </div>
        <div>
          <TanggalRedup label={REPORTING_PERIOD.akhir} value={akhir} onChange={setAkhir} />
          <Pesan teks={galat.akhir} />
        </div>
        <Pilih
          label={REPORTING_PERIOD.periode}
          value={periode}
          onChange={setPeriode}
          opsi={opsiPeriode.map((o) => ({ value: o.value, label: o.label }))}
        />
        {lain && <Field label={REPORTING_PERIOD.interval} value={interval} onChange={setSelang} error={galat.interval} />}
        {/* Label tanpa "(Days)" — layar Pega tidak menuliskannya. */}
        <Field label={REPORTING_PERIOD.penyerahan} value={penyerahan} onChange={setPenyerahan} error={galat.penyerahan} />
        <Field label={REPORTING_PERIOD.konfirmasi} value={konfirmasi} onChange={setKonfirmasi} error={galat.konfirmasi} />
        <Field label={REPORTING_PERIOD.pelunasan} value={pelunasan} onChange={setPelunasan} error={galat.pelunasan} />
      </div>

      <div className="trin__aksi">
        <button type="button" className="btn" onClick={terapkan} disabled={!bisaUbah}>
          {REPORTING_PERIOD.terapkan}
        </button>
        {/* Sel label @270162/@274479/@279526 — tampil SELALU di samping tombol,
            seperti di layar Pega; `, and Interval` hanya bila Period `other`. */}
        <span>{lain ? REPORTING_PERIOD.galatKosongInterval : REPORTING_PERIOD.galatKosong}</span>
      </div>
      {gagal !== '' && (
        <p className="trin__galat" role="alert">
          {gagal}
        </p>
      )}

      {baris.length === 0 ? (
        <Kosong pesan={REPORTING_PERIOD.tanpaBaris} />
      ) : (
        <div className="table-wrap">
          <table className="trin__tabel">
            <thead>
              <tr>
                <th scope="col">{REPORTING_PERIOD.kolomPeriode}</th>
                {/* `Auto Calculate` bersyarat `TreatyIn.ViewState != 1` @313578. */}
                {bisaUbah && <th scope="col">{REPORTING_PERIOD.kolomOtomatis}</th>}
                <th scope="col">{REPORTING_PERIOD.kolomTanggalAwal}</th>
                <th scope="col">{REPORTING_PERIOD.kolomPenyerahan}</th>
                <th scope="col">{REPORTING_PERIOD.kolomKonfirmasi}</th>
                <th scope="col">{REPORTING_PERIOD.kolomPelunasan}</th>
              </tr>
            </thead>
            <tbody>
              {baris.map((b, i) => (
                <tr key={i}>
                  <td>{b.periode}</td>
                  {bisaUbah && (
                    <td>
                      <input
                        type="checkbox"
                        aria-label={REPORTING_PERIOD.kolomOtomatis}
                        checked={b.hitungOtomatis === 'true'}
                        onChange={(e) => {
                          setBaris(baris.map((x, j) => (j === i ? { ...x, hitungOtomatis: e.target.checked ? 'true' : 'false' } : x)))
                        }}
                      />
                    </td>
                  )}
                  {(
                    [
                      ['tanggalAwal', 'tanggalAwalAsli', REPORTING_PERIOD.kolomTanggalAwal],
                      ['jatuhTempoKirim', 'jatuhTempoKirimAsli', REPORTING_PERIOD.kolomPenyerahan],
                      ['jatuhTempoKonfirmasi', 'jatuhTempoKonfirmasiAsli', REPORTING_PERIOD.kolomKonfirmasi],
                      ['jatuhTempoBayar', 'jatuhTempoBayarAsli', REPORTING_PERIOD.kolomPelunasan],
                    ] as const
                  ).map(([tampil, asli, label]) => (
                    <td key={tampil}>
                      {bisaUbah ? (
                        // ⛔ Kotak tanggal diisi nilai TERSIMPAN (`…Asli`), bukan
                        // terjemahan tampil — yang tampil membuat kotaknya kosong.
                        <TanggalRedup
                          label={label}
                          value={b[asli]}
                          onChange={(v) => {
                            setBaris(baris.map((x, j) => (j === i ? { ...x, [asli]: v } : x)))
                          }}
                        />
                      ) : (
                        b[tampil]
                      )}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Panel>
  )
}
