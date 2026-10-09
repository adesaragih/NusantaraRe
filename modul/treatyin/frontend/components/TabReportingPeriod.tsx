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

import { Field, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { PilihCari as Pilih } from '../../../../inti/frontend/components/ui/pilihSaring'
import { hitungPeriodePelaporan, type BarisPeriodeWarisan, type OpsiPilihan } from '../api'
import { useProperti } from '../halaman'
import { REPORTING_PERIOD } from '../labels'
import type { ModeForm } from '../mode'
import { KotakTanggalKetik } from './TanggalKetik'
import TanggalRedup from './TanggalRedup'
import { saringAngka } from './saringAngka'
import { SelKosongPega, TataPegaBlok } from './tataPega'
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
      {/* ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Section/TreatyInTabsProportional.xml`,
          badan `Default` L1231 (bertumpuk):

            1. `Inline grid double` L2282 — SATU sel, separuh kanan kosong →
               `Stacked with labels left` L2581: [`Inline grid double` L2880:
               Start Date (`Inline labels left` L3178) | End Date (L3907)],
               Period, Interval (`ReportingPeriod='other'`).
            2. `Inline grid triple` L5749 — Submission | Confirmation |
               Settlement, masing-masing `Inline labels left` (L6047, L6780,
               L7513).
            3. `Inline labels left` L8487 — tombol Apply + kalimat wajib isi.
            4. Grid `ReportingPeriodList` L10306.

          Satu medan `Inline labels left` digambar `kiri` — label di KIRI
          medan, sama dengan Pega.

          ⛔ Sel `g2`/`g3` di bawah memakai `div.trin__tata--kiri` telanjang,
          BUKAN `TataPegaBlok`: aturan `.trin__tata-wadah + .trin__tata-wadah`
          (jarak antarblok bertumpuk) ikut mengenai sel kedua dan ketiga grid,
          sehingga End Date, Confirmation, dan Settlement turun 12px dari
          teman sebarisnya. Di gambar Pega 01 ketiganya SEGARIS. */}
      {/* `trin__rp` — pengait jarak rapi KHUSUS tab ini (`treatyin.css` "TAB
          REPORTING PERIOD — RAPI"); `Panel` tidak menerima `className`. */}
      <div className="trin__rp">
      <TataPegaBlok tata="tumpuk">
      <TataPegaBlok tata="g2">
        <TataPegaBlok tata="kiri">
          <TataPegaBlok tata="g2">
            <div className="trin__tata trin__tata--kiri">
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
            <div className="trin__tata trin__tata--kiri">
              <TanggalRedup
                label={REPORTING_PERIOD.akhir}
                value={akhir}
                onChange={(v) => {
                  setAkhir(keSimpan(v))
                }}
              />
              <Pesan teks={galat.akhir} />
            </div>
          </TataPegaBlok>
          <Pilih
            label={REPORTING_PERIOD.periode}
            value={periode}
            onChange={setPeriode}
            opsi={opsiPeriode.map((o) => ({ value: o.value, label: o.label }))}
          />
          {lain && <Field label={REPORTING_PERIOD.interval} value={interval} onChange={(v) => { setSelang(saringAngka(v, true)) }} error={galat.interval} />}
        </TataPegaBlok>
        <SelKosongPega />
      </TataPegaBlok>
      <TataPegaBlok tata="g3">
        {/* Label tanpa "(Days)" — layar Pega tidak menuliskannya. */}
        <div className="trin__tata trin__tata--kiri">
          <Field label={REPORTING_PERIOD.penyerahan} value={penyerahan} onChange={(v) => { setPenyerahan(saringAngka(v, true)) }} error={galat.penyerahan} />
        </div>
        <div className="trin__tata trin__tata--kiri">
          <Field label={REPORTING_PERIOD.konfirmasi} value={konfirmasi} onChange={(v) => { setKonfirmasi(saringAngka(v, true)) }} error={galat.konfirmasi} />
        </div>
        <div className="trin__tata trin__tata--kiri">
          <Field label={REPORTING_PERIOD.pelunasan} value={pelunasan} onChange={(v) => { setPelunasan(saringAngka(v, true)) }} error={galat.pelunasan} />
        </div>
      </TataPegaBlok>
      </TataPegaBlok>

      {/* `Inline labels left` L8487 — tombol dan kalimatnya mengalir sebaris. */}
      <div className="trin__aksi">
        <button
          type="button"
          className="btn btn--sm"
          onClick={() => {
            terapkan()
          }}
          disabled={!bisaUbah}
        >
          {REPORTING_PERIOD.terapkan}
        </button>
        {/* Sel label @270162/@274479/@279526 — tampil SELALU di samping tombol,
            seperti di layar Pega; `, and Interval` hanya bila Period `other`.
            Huruf KECIL `.trin__catatan` — di gambar Pega 01 kalimat ini seukuran
            label medan, bukan seukuran isi. */}
        <span className="trin__catatan">{lain ? REPORTING_PERIOD.galatKosongInterval : REPORTING_PERIOD.galatKosong}</span>
      </div>
      {gagal !== '' && (
        <p className="trin__galat" role="alert">
          {gagal}
        </p>
      )}

      {/* ⭐ Bentuk Pega (8 Oktober 2026): grid tetap tampil tanpa baris —
          kepala kolom + satu sel "No items". */}
      <div className="table-wrap">
          <table className="trin__tabel trin__tabel--pega">
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
              {baris.length === 0 && (
                <tr>
                  <td colSpan={bisaUbah ? 6 : 5} className="trin__kosong-pega">
                    {REPORTING_PERIOD.tanpaBaris}
                  </td>
                </tr>
              )}
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
                        // ⛔ Kotak TANPA label tampak: judul kolom sudah di kepala
                        // grid (gambar Pega 01); `label` hanya nama pembaca layar.
                        <KotakTanggalKetik
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
      </div>
    </Panel>
  )
}
