// Chart Aggregate (section Pega `PieChartAggregate`; bentuknya dirancang sendiri - work owner 04-10-2026: "jangan
// model pie, pilih model yang modern", "grouping berdasarkan ceding, lalu ceding bisa di buka tampilin treaty type,
// lalu bisa dibuka tampilin coverage"): batang horizontal bertumpuk RNM Value (USD), seluruh As At dijumlah (As At
// dibuang dari chart, keputusan work owner 04-10-2026). Tingkat Ceding ditumpuk per Treaty Type, tingkat Treaty Type
// per Coverage, tingkat Coverage batang tunggal. Tanpa pustaka grafik.

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilRingkasan, type Ringkasan } from '../api'
import { formatAngka, formatBulat, formatPersen, tampilanRingkasan, type BatangChart, type JalurChart } from '../aturan'
import { AG } from '../labels'

/** Jumlah warna seri di aggregate.css (`aggregate__ruas--0..7`). */
const WARNA = 8

function IsiBatang({ b, bisaBuka }: { b: BatangChart; bisaBuka: boolean }) {
  return (
    <>
      <span className="aggregate__batang-nama" title={b.label}>
        {b.label === '' ? '—' : b.label}
      </span>
      <span className="aggregate__lintasan">
        {b.bagian.map((s) => (
          <span
            key={s.kunci}
            className={`aggregate__ruas aggregate__ruas--${s.seri % WARNA}`}
            style={{ width: `${s.lebar}%` }}
            title={`${s.kunci === '' ? '—' : s.kunci}: ${formatAngka(s.nilai)}`}
          />
        ))}
      </span>
      <span className="aggregate__batang-nilai" title={b.nilai}>
        {formatAngka(b.nilai)}
      </span>
      <span className="aggregate__batang-persen">{formatPersen(b.persen)}</span>
      <span className="aggregate__chevron" aria-hidden="true">
        {bisaBuka ? '›' : ''}
      </span>
    </>
  )
}

export default function ChartAggregate({ segar }: { segar: number }) {
  const [jalur, setJalur] = useState<JalurChart>({})
  const [data, setData] = useState<Ringkasan | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    ambilRingkasan().then(
      (r) => {
        if (!hidup) return
        setData(r)
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

  return (
    <PanelChart
      data={data}
      galat={galat}
      jalur={jalur}
      onJalur={setJalur}
    />
  )
}

/** Tampilan chart tanpa pengambilan data (dirender langsung oleh `ChartAggregate.test.tsx`). */
export function PanelChart({
  data,
  galat,
  jalur,
  onJalur,
}: {
  data: Ringkasan | null
  galat: unknown
  jalur: JalurChart
  onJalur: (j: JalurChart) => void
}) {
  const v = data === null ? null : tampilanRingkasan(data.irisan, jalur)

  const buka = (kunci: string) => {
    if (jalur.ceding === undefined) onJalur({ ceding: kunci })
    else if (jalur.treatyType === undefined) onJalur({ ceding: jalur.ceding, treatyType: kunci })
  }

  return (
    <section className="aggregate__kartu">
      <header className="aggregate__chart-kepala">
        <div>
          <h3 className="aggregate__judul-kartu">{AG.chart}</h3>
          <nav className="aggregate__remah" aria-label={AG.jejak}>
            <button
              type="button"
              className="aggregate__remah-tombol"
              disabled={jalur.ceding === undefined}
              onClick={() => {
                onJalur({})
              }}
            >
              {AG.semuaCeding}
            </button>
            {jalur.ceding !== undefined && (
              <>
                <span className="aggregate__remah-pisah" aria-hidden="true">
                  ›
                </span>
                <button
                  type="button"
                  className="aggregate__remah-tombol"
                  disabled={jalur.treatyType === undefined}
                  onClick={() => {
                    onJalur({ ceding: jalur.ceding })
                  }}
                >
                  {v?.namaCeding ?? jalur.ceding}
                </button>
              </>
            )}
            {jalur.treatyType !== undefined && (
              <>
                <span className="aggregate__remah-pisah" aria-hidden="true">
                  ›
                </span>
                <button type="button" className="aggregate__remah-tombol" disabled>
                  {jalur.treatyType === '' ? '—' : jalur.treatyType}
                </button>
              </>
            )}
          </nav>
        </div>
      </header>

      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={AG.memuatChart} />}

      {v !== null && (
        <>
          <div className="aggregate__kpi">
            <div className="aggregate__kpi-item">
              <span className="aggregate__kpi-label">{AG.totalNilai}</span>
              <strong className="aggregate__kpi-nilai" title={v.total}>
                {formatAngka(v.total)}
              </strong>
            </div>
            <div className="aggregate__kpi-item">
              <span className="aggregate__kpi-label">{AG.labelTingkat[v.tingkat]}</span>
              <strong className="aggregate__kpi-nilai">{formatBulat(v.batang.length)}</strong>
            </div>
          </div>

          {v.batang.length === 0 && <p className="muted">{AG.chartKosong}</p>}
          {v.batang.length > 0 && (
            <>
              <div className="aggregate__chart-sub">
                {v.seri.length > 0 && (
                  <ul className="aggregate__legenda">
                    {v.seri.map((s, n) => (
                      <li key={s}>
                        <span className={`aggregate__tanda aggregate__ruas--${n % WARNA}`} aria-hidden="true" />
                        {s === '' ? '—' : s}
                      </li>
                    ))}
                  </ul>
                )}
                {v.bisaBuka && <span className="muted">{AG.petunjukBuka[v.tingkat === 'ceding' ? 'ceding' : 'treatyType']}</span>}
              </div>
              <ol className="aggregate__batang-daftar">
                {v.batang.map((b) => (
                  <li key={b.kunci}>
                    {v.bisaBuka ? (
                      <button
                        type="button"
                        className="aggregate__batang aggregate__batang--buka"
                        onClick={() => {
                          buka(b.kunci)
                        }}
                      >
                        <IsiBatang b={b} bisaBuka />
                      </button>
                    ) : (
                      <div className="aggregate__batang">
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
