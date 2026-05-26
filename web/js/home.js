// ── Score color shift ────────────────────────────────
const SCORE_THEMES = [
  { min: 7.5, bg: '#0d1e18', border: '#1a3028', scoreColor: '#4ecfb0', labelColor: '#3a6a4a', streakColor: '#2a4a38' },
  { min: 4.5, bg: '#101828', border: '#1a2540', scoreColor: '#7aadff', labelColor: '#3a5a80', streakColor: '#2d4060' },
  { min: 2.5, bg: '#1a1018', border: '#301520', scoreColor: '#d98878', labelColor: '#6a3a3a', streakColor: '#4a2530' },
  { min: 0,   bg: '#1e0e10', border: '#3a1420', scoreColor: '#c44f4f', labelColor: '#6a2a2a', streakColor: '#4a1a20' },
];

function applyScoreTheme() {
  const card = document.getElementById('todayCard');
  if (!card) return;
  const score = parseFloat(document.getElementById('todayScore').textContent);
  const theme = SCORE_THEMES.find(t => score >= t.min);
  card.style.background = theme.bg;
  card.style.borderColor = theme.border;
  document.getElementById('todayScore').style.color = theme.scoreColor;
  document.getElementById('todayLabel').style.color = theme.labelColor;
  document.getElementById('todayTitle').style.color = theme.streakColor;
  document.querySelector('.today-streak').style.color = theme.streakColor;
}

applyScoreTheme();
document.addEventListener('htmx:afterSettle', applyScoreTheme);

// ── Radar chart ──────────────────────────────────────
const RADAR_METRICS = [
  { label: 'Energy',     value: 8, color: '#4ecfb0', angle: -90 },
  { label: 'Sleep',      value: 7, color: '#38b8a0', angle: -18 },
  { label: 'Happiness',  value: 6, color: '#7eb8f7', angle:  54 },
  { label: 'Pain',       value: 2, color: '#d98878', angle: 126 },
  { label: 'Depression', value: 3, color: '#a84560', angle: 198 },
];
const CX = 160, CY = 95, R = 60;

function radarPt(r, deg) {
  const rad = deg * Math.PI / 180;
  return [CX + r * Math.cos(rad), CY + r * Math.sin(rad)];
}

function buildRadar() {
  const svg = document.getElementById('radarChart');
  const ns = 'http://www.w3.org/2000/svg';
  svg.innerHTML = '';

  function mkEl(tag, attrs) {
    const el = document.createElementNS(ns, tag);
    Object.entries(attrs).forEach(([k, v]) => el.setAttribute(k, v));
    return el;
  }

  function polyPts(r) {
    return RADAR_METRICS.map(m => radarPt(r, m.angle).map(v => v.toFixed(1)).join(',')).join(' ');
  }

  [0.2, 0.4, 0.6, 0.8, 1.0].forEach(s => {
    svg.appendChild(mkEl('polygon', {
      points: polyPts(R * s), fill: 'none',
      stroke: '#1e2235', 'stroke-width': '1'
    }));
  });

  RADAR_METRICS.forEach(m => {
    const [x2, y2] = radarPt(R, m.angle);
    svg.appendChild(mkEl('line', {
      x1: CX, y1: CY, x2: x2.toFixed(1), y2: y2.toFixed(1),
      stroke: '#1e2235', 'stroke-width': '1'
    }));
  });

  const dataPts = RADAR_METRICS.map(m =>
    radarPt(m.value / 10 * R, m.angle).map(v => v.toFixed(1)).join(',')
  ).join(' ');
  svg.appendChild(mkEl('polygon', {
    points: dataPts,
    fill: 'rgba(122,173,255,0.09)',
    stroke: '#7aadff',
    'stroke-width': '1.5',
    'stroke-linejoin': 'round'
  }));

  RADAR_METRICS.forEach(m => {
    const [x, y] = radarPt(m.value / 10 * R, m.angle);
    svg.appendChild(mkEl('circle', {
      cx: x.toFixed(1), cy: y.toFixed(1),
      r: '3.5', fill: m.color
    }));
  });

  RADAR_METRICS.forEach(m => {
    const [x, y] = radarPt(R + 17, m.angle);
    const a = m.angle;
    const anchor = (a === -90 || a === 90) ? 'middle'
                 : (a > -90 && a < 90)     ? 'start' : 'end';
    const text = mkEl('text', {
      x: x.toFixed(1), y: (y + 4).toFixed(1),
      'font-size': '10.5',
      'font-family': '-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
      fill: m.color, opacity: '0.85',
      'text-anchor': anchor
    });
    text.textContent = m.label;
    svg.appendChild(text);
  });
}

buildRadar();

// ── Chart filters ────────────────────────────────────
function setTime(btn) {
  document.querySelectorAll('.time-btn').forEach(b => b.classList.remove('active'));
  btn.classList.add('active');
}

function setMetricPill(pill) {
  document.querySelectorAll('.pill').forEach(p => p.classList.remove('active'));
  pill.classList.add('active');
  const id = pill.dataset.metric;
  document.querySelectorAll('.metric-card').forEach(c => {
    c.classList.toggle('active', c.dataset.metric === id);
  });
  applyMetric(id);
}

function selectMetric(id) {
  document.querySelectorAll('.pill').forEach(p => {
    p.classList.toggle('active', p.dataset.metric === id);
  });
  document.querySelectorAll('.metric-card').forEach(c => {
    c.classList.toggle('active', c.dataset.metric === id);
  });
  applyMetric(id);
}

function applyMetric(id) {
  const legend   = document.getElementById('chartLegend');
  const lineWrap = document.getElementById('lineChartWrap');
  const radar    = document.getElementById('radarChart');
  const groups   = document.querySelectorAll('.mline');

  if (id === 'all') {
    legend.style.display   = 'flex';
    lineWrap.style.display = 'none';
    radar.style.display    = 'block';
  } else {
    legend.style.display   = 'none';
    lineWrap.style.display = 'block';
    radar.style.display    = 'none';
    groups.forEach(g => {
      const match = g.dataset.metric === id;
      g.style.opacity = match ? '1' : '0';
      g.querySelector('.mline-area').style.display = match ? '' : 'none';
    });
  }
}

// ── Modal ────────────────────────────────────────────
function openModal() {
  document.getElementById('modalOverlay').classList.add('open');
  document.body.style.overflow = 'hidden';
}

function closeModal() {
  const overlay = document.getElementById('modalOverlay');
  overlay.classList.remove('open');
  overlay.addEventListener('transitionend', () => {
    document.body.style.overflow = '';
  }, { once: true });
}

function handleOverlayClick(e) {
  if (e.target === document.getElementById('modalOverlay')) closeModal();
}

function updateVal(key, val) {
  document.getElementById('val-' + key).textContent = val;
}

applyMetric('energy');
