// Chart daftar Bordereaux (section Pega `PieChartBordereaux` tidak ikut diekspor; dirancang sendiri - work owner
// 04-10-2026: "group by bisnis type, lalu diklik membuka ceding nya"): batang horizontal bertumpuk JUMLAH BERKAS per
// Business, warna per Type (PREMIUM / CLAIM / SUBROGATION); klik Business = batang per Ceding Business itu. Jumlah
// berkas, bukan nominal: premi dan klaim berbeda mata uang dan tidak dapat dijumlahkan. Tanpa pustaka grafik.

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilChart, type IrisanChart } from '../api'
import { SERI_TYPE, tampilanChart, teksPersen, type BatangChart, type JalurChart } from '../aturan'
import { BDX } from '../labels'

function IsiBatang({ b, bisaBuka }: { b: BatangChart; bisaBuka: boolean }) {
  return (
    <>
      <span className="bordereaux__batang-nama" title={b.label}>
        {b.label === '' ? '—' : b.label}
      </span>
      <span className="bordereaux__lintasan">
        {b.bagian.map((s) => (
          <span
            key={s.kunci}
            className={`bordereaux__ruas bordereaux__ruas--${s.seri % 4}`}
            style={{ width: `${s.lebar}%` }}
            title={`${s.kunci}: ${s.jumlah}`}
          />
        ))}
      </span>
      <span className="bordereaux__batang-nilai">{BDX.berkas(b.jumlah)}</span>
      <span className="bordereaux__batang-persen">{teksPersen(b.persen)}</span>
      <span className="bordereaux__chevron" aria-hidden="true">
        {bisaBuka ? '›' : ''}
      </span>
    </>
  )
}

/** Tampilan chart tanpa pengambilan data (dirender langsung oleh uji). */
export function PanelChart({
  irisan,
  galat,
  jalur,
  onJalur,
}: {
  irisan: IrisanChart[] | null
  galat: unknown
  jalur: JalurChart
  onJalur: (j: JalurChart) => void
}) {
  const v = irisan === null ? null : tampilanChart(irisan, jalur)
  return (
    <section className="bordereaux__kartu">
      <header className="bordereaux__chart-kepala">
        <h3 className="bordereaux__judul-kartu">{BDX.chart}</h3>
        <nav className="bordereaux__remah" aria-label={BDX.jejak}>
          <button type="button" className="bordereaux__remah-tombol" disabled={jalur.business === undefined} onClick={() => onJalur({})}>
            {BDX.semuaBusiness}
          </button>
          {jalur.business !== undefined && (
            <>
              <span className="bordereaux__remah-pisah" aria-hidden="true">
                ›
              </span>
              <button type="button" className="bordereaux__remah-tombol" disabled>
                {jalur.business === '' ? '—' : jalur.business}
              </button>
            </>
          )}
        </nav>
      </header>

      {galat !== null && <Gagal galat={galat} />}
      {irisan === null && galat === null && <Memuat pesan={BDX.memuatChart} />}
      {v !== null && (
        <>
          <div className="bordereaux__kpi">
            <div className="bordereaux__kpi-item">
              <span className="bordereaux__kpi-label">{BDX.totalBerkas}</span>
              <strong className="bordereaux__kpi-nilai">{v.total}</strong>
            </div>
            <div className="bordereaux__kpi-item">
              <span className="bordereaux__kpi-label">{v.tingkat === 'business' ? BDX.jumlahBusiness : BDX.jumlahCeding}</span>
              <strong className="bordereaux__kpi-nilai">{v.batang.length}</strong>
            </div>
          </div>
          {v.batang.length === 0 && <p className="muted">{BDX.chartKosong}</p>}
          {v.batang.length > 0 && (
            <>
              <div className="bordereaux__chart-sub">
                <ul className="bordereaux__legenda">
                  {SERI_TYPE.map((t, n) => (
                    <li key={t}>
                      <span className={`bordereaux__tanda bordereaux__ruas--${n}`} aria-hidden="true" />
                      {t}
                    </li>
                  ))}
                </ul>
                {v.bisaBuka && <span className="muted">{BDX.petunjukBuka}</span>}
              </div>
              <ol className="bordereaux__batang-daftar">
                {v.batang.map((b) => (
                  <li key={b.kunci}>
                    {v.bisaBuka ? (
                      <button type="button" className="bordereaux__batang bordereaux__batang--buka" onClick={() => onJalur({ business: b.kunci })}>
                        <IsiBatang b={b} bisaBuka />
                      </button>
                    ) : (
                      <div className="bordereaux__batang">
                        <IsiBatang b={b} bisaBuka={false} />
                      </div>
                    )}
                  </li>
                ))}
              </ol>
            </>
          )}
        </>
      )}
    </section>
  )
}

export default function ChartBordereaux({ segar }: { segar: number }) {
  const [irisan, setIrisan] = useState<IrisanChart[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [jalur, setJalur] = useState<JalurChart>({})

  useEffect(() => {
    let hidup = true
    ambilChart().then(
      (r) => {
        if (!hidup) return
        setIrisan(r.irisan)
        setGalat(null)
      },
      (g: unknown) => {
        if (hidup) setGalat(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [segar])

  return <PanelChart irisan={irisan} galat={galat} jalur={jalur} onJalur={setJalur} />
}
