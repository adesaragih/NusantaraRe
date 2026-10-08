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
//
// ⭐ 7 Oktober 2026 — PENAMPUNG HALAMAN (`../halaman.tsx`). Medan kepala dan
// `ReportingPeriodList` hidup di halaman `TreatyIn` form, berejaan Pega, dan
// tanggalnya bentuk simpan `YYYYMMDD`. Isian bertahan saat pindah tab, dan
// tab Accumulation MEMBACA `ReportingStart`/`ReportingEnd` dari sini
// (`TreatyInSetAccountReport`).

import { useState } from 'react'

import { Field, Kosong, Panel, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { hitungPeriodePelaporan, type BarisPeriodeWarisan, type OpsiPilihan } from '../api'
import { useProperti } from '../halaman'
import { REPORTING_PERIOD } from '../labels'
import type { ModeForm } from '../mode'
import TanggalRedup from './TanggalRedup'
import { keSimpan, tanggalTampil } from './tanggalIso'

/** Satu baris `TreatyIn.ReportingPeriodList` — ejaan Pega, tanggal `YYYYMMDD`. */
export interface BarisReportingPeriod {
  Period: string
  AutoCalculate: string
  InitialDate: string
  SubmissionDue: string
  ConfirmationDue: string
  SettlementDue: string
}

/** Baris jawaban rute / kontrak → baris halaman (bentuk TERSIMPAN). */
export const barisReportingPeriod = (b: BarisPeriodeWarisan): BarisReportingPeriod => ({
  Period: b.periode,
  AutoCalculate: b.hitungOtomatis,
  InitialDate: b.tanggalAwalAsli,
  SubmissionDue: b.jatuhTempoKirimAsli,
  ConfirmationDue: b.jatuhTempoKonfirmasiAsli,
  SettlementDue: b.jatuhTempoBayarAsli,
})

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
  // ⭐ Properti halaman `TreatyIn` — ejaan ekspor (`TreatyIn.ReportingStart` …).
  const [mulai, setMulai] = useProperti('ReportingStart', '')
  const [akhir, setAkhir] = useProperti('ReportingEnd', '')
  const [periode, setPeriode] = useProperti('ReportingPeriod', '')
  const [interval, setSelang] = useProperti('ReportingInterval', '')
  const [penyerahan, setPenyerahan] = useProperti('ReportingSubmission', '')
  const [konfirmasi, setKonfirmasi] = useProperti('ReportingConfirmation', '')
  const [pelunasan, setPelunasan] = useProperti('ReportingSettlement', '')
  // Pesan Activity — milik tab, bukan halaman.
  const [galat, setGalat] = useState<Record<string, string>>({})
  const [gagal, setGagal] = useState('')
  const [baris, setBaris] = useProperti<BarisReportingPeriod[]>('ReportingPeriodList', () => barisWarisan.map(barisReportingPeriod))
  const bisaUbah = mode === 'ubah'
  // `TreatyIn.ReportingPeriod='other'` @144723 — Interval dan `, and Interval`.
  const lain = periode === 'other'

  /** `TreatyInSetReport` — `awal` = `param.startdate` (kosong: ReportingStart). */
  function terapkan(awal = '') {
    setGagal('')
    hitungPeriodePelaporan({ mulai, akhir, periode, interval, penyerahan, konfirmasi, pelunasan, awal })
      .then((h) => {
        setGalat(h.galat)
        // ⛔ Langkah 2 Activity: daftar DIBUANG lalu disusun ulang — tetapi
        // isian kosong KELUAR di langkah 4–7 SEBELUM ada baris baru. Daftar
        // kosong hanya bila jawabannya nol pesan.
        if (Object.keys(h.galat).length === 0) setBaris(h.baris.map(barisReportingPeriod))
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
          {/* Kotak tanggal menerima bentuk simpan; jawabannya bentuk kabel →
              disimpan kembali `YYYYMMDD`. */}
          <TanggalRedup
            label={REPORTING_PERIOD.mulai}
            value={mulai}
            onChange={(v) => {
              setMulai(keSimpan(v))
            }}
          />
          <Pesan teks={galat.mulai} />
        </div>
        <div>
          <TanggalRedup
            label={REPORTING_PERIOD.akhir}
            value={akhir}
            onChange={(v) => {
              setAkhir(keSimpan(v))
            }}
          />
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
        <button
          type="button"
          className="btn"
          onClick={() => {
            terapkan()
          }}
          disabled={!bisaUbah}
        >
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
                  <td>{b.Period}</td>
                  {bisaUbah && (
                    <td>
                      <input
                        type="checkbox"
                        aria-label={REPORTING_PERIOD.kolomOtomatis}
                        checked={b.AutoCalculate === 'true'}
                        onChange={(e) => {
                          setBaris(baris.map((x, j) => (j === i ? { ...x, AutoCalculate: e.target.checked ? 'true' : 'false' } : x)))
                        }}
                      />
                    </td>
                  )}
                  {(
                    [
                      ['InitialDate', REPORTING_PERIOD.kolomTanggalAwal],
                      ['SubmissionDue', REPORTING_PERIOD.kolomPenyerahan],
                      ['ConfirmationDue', REPORTING_PERIOD.kolomKonfirmasi],
                      ['SettlementDue', REPORTING_PERIOD.kolomPelunasan],
                    ] as const
                  ).map(([kunci, label]) => (
                    <td key={kunci}>
                      {bisaUbah ? (
                        // ⛔ Kotak tanggal diisi nilai TERSIMPAN, bukan terjemahan
                        // tampil — yang tampil membuat kotaknya kosong.
                        <TanggalRedup
                          label={label}
                          value={b[kunci]}
                          onChange={(v) => {
                            setBaris(baris.map((x, j) => (j === i ? { ...x, [kunci]: keSimpan(v) } : x)))
                            // ⭐ `change` sel Initial Date → `TreatyInSetReport(startdate =
                            // .InitialDate, autocalculate = .AutoCalculate)`. Langkah 1
                            // KELUAR bila Auto Calculate baris itu tidak dicentang.
                            if (kunci === 'InitialDate' && b.AutoCalculate === 'true') terapkan(keSimpan(v))
                          }}
                        />
                      ) : (
                        tanggalTampil(b[kunci])
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
